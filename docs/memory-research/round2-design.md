# Round 2 design: a pointer-free AST, per-phase columns and AST rebuilds

Status: proposal (2026-10-23). Inputs: the six reports in this directory, the
round-1 agents' measurements, and the fragmentation probe below.

## Why the AST, and why not "just smaller"

- The AST is the largest long-lived structure: 32 % of typed-full's peak heap,
  76–89 % of untyped builds' heap. Each `rhino.Node` is 64 B with 5 pointer
  words (next, first, parent, previous, typeOrColor).
- It is **not** bigger than Java's: a Closure `Node` is ~56 B under compressed
  oops and every non-inline property adds a 24 B `PropListItem` (jsgo inlines
  40 boolean props). What Java does not pay is re-tracing it: G1 is
  generational and leaves the old AST alone, while Go's collector re-marks
  every node pointer every cycle (80–90 % of jsgo's peak heap is
  pointer-bearing; `peak_scan_bytes` in memgate).
- So the goal for the AST is "invisible to the GC", and fewer bytes second.
  (Types are the opposite case: 2–3× Java's size; shrinking them is a round-1
  item.)

## Measured: fragmentation is small, slab pinning is not

A probe forced a full GC after every pass and read
`/memory/classes/heap/unused` (holes in partially used spans):

| Task | Variant | Max in-use spans | Live | Holes |
|---|---|---:|---:|---:|
| bundle-mega | slabs (before round 1) | 697 MB | 674 MB | 23 MB (3 %) |
| bundle-mega | `Slab.New` → `new(T)` | 640 MB | 604 MB | 36 MB (6 %) |
| typed-full | slabs | 510 MB | 497 MB | 13 MB (3 %) |
| typed-full | `new(T)` | 471 MB | 458 MB | 13 MB (3 %) |

Holes stay at 2–7 % through the whole pipeline: rebuilding the AST through
Go's allocator to "defragment" would recover at most that much, and a copy does
not come out contiguous either (the 2026-10-18 GC lab measured that
re-allocating a churned AST in pre-order recovers little of the 2.5× locality
gap). The real waste was **dead objects pinned inside live slab chunks**
(40–70 MB more "live" heap with slabs), which round 1 removes.

A rebuild is worth it only as the reclamation step of a pointer-free AST.

## Counter-evidence from round 1 (rhino.Node agent, measured)

Nodes allocated over a run vs live at the peak: typed-full 2.99 M / 2.15 M
(ratio 1.39), typed-mega 1.42, typed-app 1.49, bundle 1.96 by the end. A
48-B node in a never-freed arena beats 64-B GC-managed nodes on bytes only
below a ratio of 64/48 = 1.33, so **without reclamation the byte gain is a
wash** (typed-full: arena ≤ 144 MB vs 138 MB of scanned nodes today); the
gain is GC mark work (−140 to −240 MB of scannable heap per cycle, about 30 %
of marking). Also: types and JSDoc reached only through ids must be held by
per-compilation tables, which re-creates the after-colors retention that
per-node allocation just fixed, unless the type column is dropped at colors.

So the order below changes: first the cheap, safe steps (Go 1.27; a 32-bit
type/color handle, which removes 1 of 5 pointer words per node and lets the
type column be dropped wholesale), and the pointer-free arena only together
with compaction, which requires an exhaustive census of `*Node` holders that
live across a compaction point (pass-local maps die with their pass; the risk
is Compiler fields and other long-lived holders).

## Design

### 1. Node core: 48 B, no Go pointers

```go
type Node struct {
    next, first, parent, previous nodeID // int32 links (0 = nil)
    self    nodeID   // own id: index into per-phase columns
    aux     uint32   // handle into the cold-data table (JSDoc, rare props, overflow length/original name)
    payload uint64   // double bits or two 32-bit string ids (as today)
    propFlags uint64 // 40 bool props + small enums + source file id (as today)
    linenoCharno int32
    length16 uint16
    token    Token
    kindFlags uint8
}
```

- Nodes live in chunks `[4096]Node` (192 KiB, allocated as noscan because the
  struct has no pointers) referenced from a per-compilation chunk directory.
  Parallel parsing reserves whole chunks per goroutine.
- The public API is unchanged: passes keep `*rhino.Node` (an interior pointer
  into a chunk) and the accessors translate ids (`n.GetNext()` =
  `dir.at(n.next)`). The link fields are private to package rhino (245 uses
  in node.go; 35 allocation sites repo-wide), so the change stays in rhino plus
  the allocation sites.
- `typeOrColor` and the ext box leave the node:
  - **type column** (`[]*JSTypeBase` by node id) exists only from type
    inference to typesToColors, then is dropped wholesale — which also fixes
    "a dead typed node pins the whole registry" by construction;
  - **color column** (`[]uint32` color ids, pointer-free) after typesToColors;
  - **cold table** behind `aux` for JSDoc and rare properties.

### 2. Reclamation: copying compaction at safe points

The GC no longer sees node-to-node references, so dead nodes are never freed by
it (jsgo allocates ~2.6× more nodes than are live). Reclaim by copying the live
tree into fresh chunks in pre-order and dropping the old chunks:

- **When:** right after typesToColors (typed builds shed types, JSDoc type
  trees and checks-only garbage), every N iterations of the optimization loop,
  and at stage boundaries (CHECKS → OPTIMIZATIONS → FINALIZATIONS, where a
  TypedAST round trip already rebuilds the AST today).
- **Cost:** ~10–20 ns per live node; a few tens of ms for millions of nodes,
  against seconds per loop iteration.
- **Bonus:** pre-order layout (walks ~2.5× faster than a churned heap per the
  GC lab) and a natural point to switch per-phase columns.

### 3. The hard part: cross-pass `*Node` holders

Every `*Node` held across a compaction point must be forwarded (old node stores
its new id; holders are rewritten). Holders to audit: compiler roots and
`CompilerInput` AST roots, JSDoc type-expression roots, function types'
`source` and properties' `propertyNode` (gone after colors), change-tracking
state, cached analyses on the Compiler (GlobalNamespace caches,
ModuleMetadata, symbol tables), side maps keyed by `*Node`.

Safety net: a debug mode that keeps old chunks alive but **poisons** them
(token = POISON; accessors assert) so any stale holder fails loudly, and runs
the full parity corpus (40k Java builds) in that mode before the default flips.

### 4. Migration steps (each measured with memgate, byte-identical)

1. **Measure:** dead/live node counts per pass (done for whole runs, above),
   and a census of `*Node` holders alive at each candidate compaction point.
2. **Type/color handle:** `typeOrColor` becomes a 32-bit id into
   registry-owned tables (5 → 4 pointer words per node, −20 % AST marking; the
   JSType table is dropped wholesale at colors). Safe, no reclamation needed.
3. **Columns for cold data, int32 links:** the node becomes pointer-free
   (48 B, noscan); AST scan work → 0. Only together with step 4.
4. **Compaction after typesToColors and in the loop**, behind `--fast_mode`,
   with the poison verifier in tests.
5. **Default** once the parity corpus is clean in poison mode.
6. With the AST unscanned, GC cycles are cheap: lower the default GOGC (or add
   a `-Xmx`-like soft limit via `debug.SetMemoryLimit`) to turn the saved GC
   CPU into lower RSS.

### Expected effect (estimates)

- AST bytes −25 % (64 → 48 B); AST mark work → 0 (50–85 % of all GC scan
  work today).
- Typed builds: types/JSDoc trees released wholesale at colors (the lifetimes
  agent measured ~100 MB retained after colors on typed-full before round 1).
- Untyped builds: dead nodes reclaimed at every compaction; RSS bounded by
  ~2× live AST at the copy point.
- With GOGC lowered to spend the saved CPU: peak RSS −20 % or more on large
  builds.

## Other round-2 candidates (from the research, not done in round 1)

- Colors as compact, interned id lists (`ownProperties` as ordered string-id
  slices; union colors interned).
- Hash-consing nullable/optional unions and templatized types (identity audit
  first; fast mode if identity matters).
- Positions as offsets + per-file line tables end-to-end (parser and node).
- Dense node-id-indexed side tables replacing `map[*rhino.Node]...`.
- Upgrade the toolchain to Go 1.27 (Green Tea GC): measured −3 to −9 % RSS and
  up to −35 % GC CPU with no code change.
