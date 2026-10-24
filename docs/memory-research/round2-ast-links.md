# Round 2: int32 node links, chunked node storage and copying compaction

Status: implemented (2026-10-23/24) on the `ast-links` worktree, merged
with the `ast-handles` work (`typeOrColor` → 32-bit `tc` + `aux` handles,
`node_tables.go`); measurements below. Companion to `round2-design.md`
(design) and `node-holder-census.md` (the holder list the forwarding
registry implements). With both, `Node` holds no Go pointer at all
(`TestNodeHasNoPointers`).

## 1. Node storage and the id scheme (`internal/rhino/node_store.go`)

- `Node.next/first/parent/previous` are `nodeID` (uint32); a new `self`
  field holds the node's own id. `Node` is 56 bytes (was 64): 4 × 4 B links
  + 4 B self + payload 8 + propFlags 8 + lineno 4 + tc 4 + aux 4 +
  length 2 + token 1 + kindFlags 1 (52, padded to 56; the 48-byte layout of
  the design would need `tc` as a per-chunk column indexed by `self`, one
  more dependent load per type/color access — not done).
- Chunks: `nodeChunk = { header (one node wide); nodes [479]Node }`,
  480 × 56 B = 26,880 B. Slot 0 is the header (`owner *nodeGen`, `index`,
  bump counter), so **id = chunk<<9 | slot with slot in 1..479 and id 0 =
  nil**. 2^23 chunks of id space (4 G nodes per process; the last chunk
  index is reserved for the forwarding marker). The size is chosen for
  the GC: a chunk is a *small* object (≤ 32 KiB, size class 27,264 with the
  8-byte malloc header, 1.4 % slack), and Go accounts and scans a small
  object with pointers only up to its last pointer (`PtrBytes`), i.e. the
  header's `owner` word; a large object is accounted whole
  (`heapSetTypeLarge` returns the element size). So the AST costs the GC
  8 bytes per chunk: `peak_scan_bytes` no longer contains it (section 3).
- Directory: one process-wide table `nodeChunkTable [2^23]uintptr` (64 MiB
  of BSS, demand-paged, 8 B per chunk ever allocated, **not traced by the
  GC**: the slots are integers; they are *written* as integers too — a
  pointer store would carry a write barrier that shades the slot's old
  value, and a cleanup clears the slot of a chunk the GC already freed).
  `nodeAt(id)` = `(*Node)(chunkTable[id>>9] + (id&511)*56)`: an atomic load
  of the slot and no bounds check. A cleared slot panics ("stale node id")
  instead of dereferencing.
- `*Node` stays the handle passes hold: an interior pointer into a chunk;
  `GetNext()` etc. translate ids. `idOf(n)` panics for a by-value node
  (`self == 0`): the by-value template nodes (`SourceFile`, `IRFactory`,
  the TypedAST deserializer) are only read as templates and keep working;
  `FunctionToBlockMutator`'s scratch ROOT/SCRIPT, which get linked, became
  chunk nodes.

### Lifetime: generations, and how a node finds its directory

The table is global, so a node needs no directory pointer. What makes the
global table safe across compilations (tests run thousands in one process)
is the ownership scheme:

- A **generation** (`nodeGen`) owns a list of chunks; every chunk header
  points back at its generation. `NewCompiler` starts a generation and makes
  it current (`rhino.NewNodeGeneration`); compaction starts another.
- Any held `*Node` keeps its chunk alive, hence its generation, hence
  every chunk of the generation, hence everything its links can reach.
  The table holds untraced addresses, so a generation that nothing
  references is collectable; its `runtime.AddCleanup` (one per chunk)
  clears the table slots. `TestNodeGenerationsAreReleased` checks that 20
  finished compilations leave no live chunk behind.
- Chunk indices are **never reused**, so a stale id always hits an empty
  slot and panics; it can never alias another chunk.
- Nodes are allocated in whichever generation is current, which is only a
  lifetime grouping: ids are global. A link written between nodes of two
  generations (`Node.ref`, every link write goes through it) **ties** the
  two generations (each keeps the other in a `peers` set), because an id
  keeps nothing alive by itself; in one compilation this never happens (a
  compaction's successor never links back), in the test suite it happens
  constantly (the compiler under test keeps allocating into the generation
  of a later helper compiler). Without the tie the first merged test runs
  crashed in the GC ("found pointer to free object").
  Limit: a stale holder (a bug) that traverses into a released generation
  panics on the cleared slot; between a generation's collection and its
  cleanup (a few ms) such a traversal would read freed memory.
- The initial generation (package-level sentinel nodes, allocated during
  `init`) is never released and never compacted; `Forward` leaves its nodes
  in place (the JSDoc serializer's placeholder type expression is one).

### Concurrency

Allocation is lock-free: an atomic add on the current chunk's bump counter;
when the chunk is full, a mutex-protected slow path allocates and publishes
the next chunk (atomic slot store). Parallel parsing and the parallel passes
(coalesce, dead assignments, flow-sensitive inlining, object literals)
allocate concurrently without any plumbing; ids are stable for the process.
Chunk switches are the only contention point (once per 4095 nodes).

## 2. Compaction (`internal/rhino/node_compact.go`, `internal/jscomp/compiler_compact.go`)

`rhino.Compactor` copies the live tree in pre-order into a fresh generation
(56-B copy + 4 link writes per node, ~35 ns/node measured), forwards holders
and releases the old generations:

- An old node remembers its copy (`old.next = newID`,
  `old.previous = nodeForwardMarker`), so `Forward` is idempotent and holders
  can be forwarded in any order. A holder of a node not reachable from the
  roots gets the **whole detached subtree it belongs to** copied as an extra
  root (the holder may point at an interior node, e.g. a module binding's
  NAME inside a removed module body).
- Side trees owned by nodes are copied with them: JSDoc type-expression
  roots and type transformations (once per `JSDocInfo`), and the node-valued
  props `DECLARED_TYPE_EXPR` / `CLOSURE_UNAWARE_SHADOW`.
- The aux table (`node_tables.go`) is rebuilt: each copied node's aux
  entry (its `*JSDocInfo` or `*nodeExtBox`) moves to a fresh block and
  `Finish` swaps the compilation's aux allocator to those blocks and
  releases the old ones (`NodeTables.rotateAux`), so the cold data of dead
  nodes is dropped too. Type and color handles are object-owned and need
  nothing.
- `Finish` clears the old chunks' table slots immediately and drops the
  generations' chunk lists (a stale holder keeps only its own chunk alive).
- Before the compaction that follows `typesToColors` the driver runs
  `runtime.GC()`: the type graph was just dropped, and without the
  collection the copy allocated a second AST on top of the not-yet-reclaimed
  types (typed-full peak RSS 731 MB instead of 535 MB in the first runs).

### Forwarding registry (`Compiler.compactAST`)

The explicit, checked list, from the census (section 3.1):

1. `externsRoot`, `jsRoot`, `externAndJsRoot` (the tree copied first), the
   caller's local `externs`/`root` (PhaseOptimizer and Loop pass them by
   pointer), `scriptNodeByFilename`, `CompilerInput.ast.root` for
   `inputsById`, the chunk graph's inputs and `externs`,
   `phaseOptimizer.jsRoot`.
2. `ChangeTracker.changeTimeline` (`Timeline.ForwardValues`: values
   rewritten, the pointer-keyed index rebuilt).
3. `ModuleMetadataMap` root nodes (nested included); `ModuleMap`
   `Binding.sourceNode`, `Export.exportNode/nameNode` through `namespace`,
   `boundNames`, `localNameToLocalExport`.
4. `JSError.node` of the error manager's errors and warnings.
5. `RuntimeJsLibManager.lastInjectedLibrary`.
6. `transpilationNamespace` (`GlobalNamespace.forwardNodes`): roots,
   `spreadSiblingCache` and `defineCandidateScripts` (rebuilt), every
   `Name.declaration` and `Ref.node` (with `refsIndex` rebuilt), and the
   refs' `Scope` chains (`rootNode`, `Var.nameNode`).
7. The optimization loop's `scriptSizes` (rebuilt through the `extra`
   callback).
8. Go-only: `cachedCCAArena` evicted. Asserted nil: `typedScopeCreator`,
   `topScope`, `CompilerInput.typedScope` (compaction only runs after
   `typesToColors`, and never with `PreservesDetailedSourceInfo`).
9. Never compacted: a `NestedCompilerRunner`'s shadow compiler (its tree is
   spliced into, and held by, the outer compilation).

### Policy

Hooks after every pass of `PhaseOptimizer.Process` and at the end of every
`Loop` iteration, plus the end of the optimization stage; compaction is
forced right after `typesToColors` and otherwise runs when the nodes
allocated since the last compaction (chunk count × 4095, free to read)
exceed a quarter of the live nodes at the last compaction (and 64 k).
On by default (after the poison runs below came back clean);
`JSGO_AST_COMPACT=0|1` overrides, `JSGO_AST_COMPACT_LOG=1` logs,
`JSGO_AST_POISON=1` is the verifier.

### Poison mode results

Old chunks are kept alive, their slots cleared (any id into them panics in
`nodeAt`) and every old node overwritten with an unresolvable id and token
0xFF, so an unforwarded holder fails on first use.

- memgate `--fast` with `JSGO_AST_POISON=1`: typed-app, bundle, checks-only,
  typed-full, typed-mega, bundle-mega all MATCH. Compactions after
  `typesToColors`: typed-app 314,196 nodes copied (153 extra roots / 676
  nodes; the census predicted 152 / 666), 16 ms; bundle 1.20 M (1 / 1),
  29 ms; typed-full 2.01 M (3,290 / 21,792; census 3,274 / 21,752), 73 ms;
  typed-mega 3.68 M (3,843 / 23,912), 149 ms; bundle-mega 4.28 M (1 / 1),
  136 ms.
- `JSGO_AST_COMPACT=1 JSGO_AST_POISON=1 go test ./internal/jscomp`: green
  after two fixes it found (the shared JSDoc placeholder sentinel, and the
  nested shadow compiler compacting a tree the outer compilation holds);
  green again on the merged tree (compaction on by default, 177 s).
- Known pre-existing flake (handles-only baseline 0f0b036e, 1 run in 8):
  `TestTranspileAndOptimizeClosureUnawareTest` loses a
  `CLOSURE_UNAWARE_SHADOW` prop ("reportCodeChange should have been
  called") — an aux handle released under the nested compiler; not from
  this work.

## 3. Measurements

memgate on the six tasks, one binary per step, all MATCH (byte-identical
with Java in default mode, fast mode and poison mode). "base" is the branch
with the handles work only (0f0b036e: `tc`/`aux` handles, pointer links);
"ids" is the merged tree with `JSGO_AST_COMPACT=0` (step 1 alone: ids, no
reclamation); "compact" is the merged tree as shipped (compaction on). The
machine ran three other agents' builds throughout, so wall and CPU are
±10 % noisy; live heap, objects and GC CPU are the robust columns.

### Default GOGC (RSS, wall, CPU, GC CPU)

| task | step | wall s | cpu s | RSS MB | live MB | scan MB | M obj | GCs | GC cpu s |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| typed-app | base | 2.59 | 1.72 | 181 | 69 | 33 | 1.8 | 13 | 0.30 |
| | ids | 2.24 | 1.70 | 184 | 72 | 45 | 1.3 | 13 | 0.13 |
| | compact | 1.60 | 1.70 | 171 | 71 | 44 | 1.3 | 14 | 0.12 |
| bundle | base | 4.64 | 5.52 | 302 | 127 | 78 | 3.5 | 19 | 1.03 |
| | ids | 4.98 | 5.15 | 310 | 138 | 123 | 1.8 | 17 | 0.15 |
| | compact | 5.40 | 5.58 | 280 | 124 | 99 | 1.9 | 22 | 0.19 |
| checks-only | base | 2.30 | 2.88 | 324 | 157 | 89 | 4.8 | 15 | 0.60 |
| | ids | 2.45 | 2.56 | 303 | 153 | 129 | 3.3 | 14 | 0.11 |
| | compact | 2.65 | 2.76 | 303 | 156 | 127 | 3.3 | 14 | 0.11 |
| typed-full | base | 11.42 | 12.77 | 575 | 261 | 161 | 9.9 | 23 | 2.31 |
| | ids | 13.85 | 14.43 | 583 | 255 | 217 | 7.2 | 23 | 0.83 |
| | compact | 12.68 | 13.16 | 475 | 230 | 196 | 7.2 | 26 | 0.53 |
| typed-mega | base | 18.20 | 22.02 | 980 | 430 | 261 | 15.5 | 26 | 4.17 |
| | ids | 22.80 | 23.84 | 944 | 454 | 368 | 10.7 | 25 | 1.24 |
| | compact | 20.05 | 20.81 | 787 | 382 | 335 | 10.8 | 28 | 0.78 |
| bundle-mega | base | 23.33 | 28.17 | 919 | 435 | 242 | 15.5 | 29 | 4.90 |
| | ids | 30.63 | 31.37 | 956 | 449 | 413 | 10.2 | 29 | 0.81 |
| | compact | 27.49 | 28.29 | 877 | 414 | 378 | 10.2 | 29 | 0.77 |

- Step 1 (ids, no reclamation): bytes are a wash, as predicted by
  `round2-design.md` (56-byte nodes in never-freed chunks vs 64-byte
  GC-managed ones at an allocated/live ratio of 1.4–2); objects −27…−48 %;
  **GC CPU −57…−84 %** (the collector no longer traces node links; a
  chunk costs it one pointer); CPU +0…+13 % from the id translation on
  every link read (a dependent table load) and the lost per-node freeing.
- Step 3 (compaction): RSS −5…−20 % and peak live −3…−12 % vs base on the
  optimizing tasks (dead nodes and the cold data of dead nodes are
  reclaimed, and the pre-order copy restores locality); CPU back to base
  (−6…+4 %), GC CPU −60…−84 %. checks-only has no optimization stage and
  so no compaction. Compactions on the shipped binary: typed-app 314 k
  nodes copied in 9 ms (152 extra roots / 666 nodes); bundle 1.20 M in
  21 ms, then one loop compaction (1.18 M, 20 ms); typed-full 2.01 M in
  49 ms (3,274 / 21,752); typed-mega 3.68 M in 87 ms (3,826 / 23,866);
  bundle-mega 4.28 M in 73 ms (1 / 1). The copy runs at 23–27 ns per node.
- "scan MB" does **not** drop, and rises with ids: Go accounts every
  object above 512 bytes that contains a pointer at its full size
  (`heapSetTypeSmallHeader`/`heapSetTypeLarge` return `span.elemsize`),
  and a chunk holds one pointer (its generation). The mark loop still
  stops at that pointer — the GC CPU column is the real cost — but the
  pacer's `heapScan` and memgate's `peak_scan_bytes` count the whole AST.
  Making chunks noscan would mean giving up the GC-driven lifetime
  (section 4) for a Compiler-death release, or 512-byte chunks of 8 nodes
  (+11 % AST bytes); neither is done here.

### `--gogc 10` (stable peak live)

| task | base live / scan / GC cpu | ids | compact |
|---|---|---|---|
| typed-app | 84 / 45 / 2.62 | 84 / 56 / 1.31 | 82 / 50 / 0.86 |
| bundle | 164 / 107 / 10.09 | 154 / 133 / 1.82 | 159 / 144 / 1.74 |
| checks-only | 188 / 118 / 5.01 | 175 / 150 / 1.34 | 177 / 151 / 1.26 |
| typed-full | 270 / 167 / 24.91 | 268 / 229 / 7.49 | 290 / 258 / 4.47 |
| typed-mega | 441 / 264 / 41.87 | 434 / 351 / 11.60 | 529 / 490 / 6.67 |
| bundle-mega | 447 / 272 / 67.85 | 454 / 418 / 7.61 | 494 / 468 / 7.96 |

At GOGC 10 the collector samples the live heap during the copy itself, so
the compaction's transient (old and new AST both alive for 50–90 ms) shows
up as a higher "peak live" (+7…+20 % on the mega tasks); at default GOGC it
is absorbed by the heap goal and peak RSS still falls. A chunk-by-chunk
release during the copy would remove it (nodes of one script are
chunk-contiguous right after parsing, but not after optimization passes).

### Step 4: GOGC trade-off with compaction (default GOGC = 100)

| task | GOGC 100 RSS / cpu / GC cpu | GOGC 75 | GOGC 50 | base (GOGC 100) |
|---|---|---|---|---|
| typed-app | 171 / 1.70 / 0.12 | 155 / 1.89 / 0.14 | 140 / 1.66 / 0.18 | 181 / 1.72 / 0.30 |
| bundle | 280 / 5.58 / 0.19 | 304 / 5.37 / 0.27 | 241 / 5.35 / 0.34 | 302 / 5.52 / 1.03 |
| checks-only | 303 / 2.76 / 0.11 | 288 / 2.61 / 0.16 | 294 / 3.11 / 0.26 | 324 / 2.88 / 0.60 |
| typed-full | 475 / 13.16 / 0.53 | 450 / 12.17 / 0.60 | 466 / 12.51 / 0.94 | 575 / 12.77 / 2.31 |
| typed-mega | 787 / 20.81 / 0.78 | 765 / 20.68 / 0.86 | 781 / 21.14 / 1.45 | 980 / 22.02 / 4.17 |
| bundle-mega | 877 / 28.29 / 0.77 | 853 / 28.70 / 1.08 | 689 / 29.38 / 1.40 | 919 / 28.17 / 4.90 |

GOGC 50 with compaction costs 0.2–1.5 s of GC CPU — still below base at
GOGC 100 on every task — for RSS −9…−25 % vs base (bundle-mega 689 vs
919 MB, typed-mega 781 vs 980, typed-full 466 vs 575, bundle 241 vs 302,
typed-app 140 vs 181). The default GOGC is left unchanged here (the runner
can set `GOGC`); lowering it to 50 is the recommended follow-up once the
other agents' allocation-volume work is in, since GC CPU scales with the
remaining pointerful heap (types, scopes, strings), not with the AST.

## 4. Pointer-free, not noscan

`Node` has no pointer; the only pointer in a chunk is the header's `owner`
(offset 0), kept on purpose: it is what lets a held `*Node` keep its
generation (and everything its links reach) alive, so chunks are released
by the GC's own reachability instead of by a compiler-lifetime rule that a
test holding a tree after dropping its compiler would break. Because the
chunk is a small object, that one pointer is all the GC accounts and
scans (section 1); a genuinely noscan chunk would save nothing more.

## 5. Address-order audit

Compaction changes every node's address. The only address-derived values in
the compiler are hash bucket placements, never orderings: the two
`uintptr(unsafe.Pointer(n))` hashes in `type_inference.go` (type inference
runs before `typesToColors`, i.e. before any compaction) and the type-pair
hashes in `jstype_js_type_registry.go` (the registry is dropped at colors).
Pointer-keyed Go maps (`Timeline`, `GlobalNamespace.refsIndex`, the loop's
`scriptSizes`, the two GlobalNamespace caches) are rebuilt on forwarding;
their iteration order is random in Go anyway, so nothing byte-identical
could depend on it. The six benchmark tasks match Java byte for byte with
compaction on, in fast and default mode, and in poison mode.
