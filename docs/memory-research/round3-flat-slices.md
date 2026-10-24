# Round 3 spike: nodes in slices instead of chunks (2026-10-24)

Question: round 2 stores nodes in 480-node chunks behind a process-wide
chunk table so that `*Node` stays a valid handle. Would plain slices
(`[]Node`, links are indexes) be faster, and should there be one slice per
file (cheap, local, parallel-friendly) that becomes one unified slice only
when cross-file optimizations need it?

Spike code: `internal/rhino/node_flat_spike.go` (copies a parsed tree into
slices, walk kernels) and `internal/rhino/flatspike/` (benchmarks on real
inputs). Not used by the compiler.

Machine: Linux amd64, 4 cores, 15 GB, Go 1.27.3. Corpus: the typed-full
inputs (1,820 files, 1,898,797 nodes), parsed by the real parser.

## 1. Lookup cost: slices vs chunks

Two kernels, ns per node, median of 6 runs (20 iterations each):

- **walk**: recursive first-child/next-sibling walk reading every node's
  token (what `NodeTraversal` does).
- **parent chains**: from every node, follow `parent` to the root (what
  `getParent()`-heavy checks do).

| Layout | walk | parent chains |
|---|---|---|
| chunks, `*Node` API, as parsed (today, before compaction) | 20.3 | 41.6 |
| chunks, `*Node` API, after compaction (today, pre-order) | 16.5 | 37.4 |
| chunks, inlined table lookup, as parsed | 18.1 | 24.7 |
| chunks, inlined table lookup, after compaction | 14.1 | 22.4 |
| per-file slices, parse order | 15.9 | 21.3 |
| per-file slices, pre-order | 12.5 | — |
| per-file slices addressed by a global (file, index) id | 19.1 | — |
| one unified slice, parse order | 15.3 | 20.6 |
| one unified slice, pre-order | **11.2** | **18.2** |

Reading:

- A slice in tree order is the fastest layout: 32 % faster than today's
  compacted chunks on walks, 51 % on parent chains.
- About half of today's overhead is a missed inline: `nodeAt` costs 100
  against Go's budget of 80, so every `GetNext`/`GetFirstChild`/`GetParent`
  is a real call (7 % of typed-full CPU is `nodeAt` itself, 0.77 s of it on
  the function entry). See section 4 for the fix.
- Layout matters as much as the lookup: pre-order beats parse order by
  20–30 % in every store. A rebuild after a pass gets this for free.
- **Per-file and unified slices cost the same** when the code has the
  file's slice in hand (12.5 vs 11.2). Per-file slices addressed through a
  global (file, index) id lose the whole gain (19.1, no better than today):
  the second level is the chunk table again.

## 2. How much do passes move nodes across files?

Instrumented run (`JSGO_AST_COMPACT=0`): after every pass, every node's
SCRIPT is compared with its SCRIPT after the previous pass.

| Task | Passes run | Passes that move nodes to another file | Nodes moved |
|---|---|---|---|
| typed-full (ADVANCED, 1 chunk) | 207 | 3: `optimizeCalls` (2,018 nodes, 75 files), `inlineVariables` (599), `earlyInlineVariables` (362) | ~3,000 of 1.9 M |
| bundle (SIMPLE) | 128 | none | 0 |

Passes that copy code across files (`inlineFunctions`) clone it, so the
clone is a new node, not a moved one. Neither task has several output
chunks, so `crossChunkCodeMotion`/`crossChunkMethodMotion` did not run;
they move whole statements between files and belong to the same class as
the three above.

Nodes created after parsing: typed-full 484 k (inlineFunctions 138 k,
latePeephole 57 k, the ES6 rewrites ~140 k), bundle 253 k (inlineFunctions
157 k), against 1.9 M / 1.1 M parsed.

## 3. Merging per-file slices into one

Concatenating the 1,820 per-file slices into one and rebasing every link:
67–164 ms (unoptimized: fresh 106 MB allocation each time, single-threaded).
Negligible once per compilation; too much to do around every pass.

## 4. Quick win independent of the storage decision

`nodeAt` with the stale check as an inline `panic` (no call) and a direct
atomic load inlines at cost 34, keeping every check:

| Task | before (5 runs, median) | after (median) | change |
|---|---|---|---|
| typed-full | 27.11 s (25.87–27.51) | 25.26 s (24.70–26.99) | −6.8 % |
| bundle | 11.40 s (10.90–12.17) | 10.70 s (10.26–11.07) | −6.1 % |

Interleaved runs, outputs byte-identical. The stale-id panic message no
longer carries the id (a `fmt` call would not inline).

## 5. What this means for the design

The handle migration is the same work in every variant: `NodeID` (uint32)
plus a `*NodeStore` holding `[]Node`, accessors on the store
(`s.Next(id)`, `s.Token(id)`), creation through the store
(`s.NewName(...)`), `map[*Node]` → `map[NodeID]` (56 sites), about 9,700
`*Node` uses in 392 files plus 1,500 in tests. A pass is handed the store
it works on; it does not need to know whether that store is one file or
the whole program.

That makes per-file vs unified a layering decision on top, not a fork:

- **Per-file stores** fit parsing and script-local passes: parallel workers
  (main's Round 3A already runs peephole, InlineFunctions discovery and
  several checks per script) allocate into their own file's slice without
  locks, growth copies one file, and a file can be rebuilt (pre-order,
  garbage dropped) on its own.
- **The unified store** is required by the passes that move nodes between
  files (section 2) and by global analyses that follow handles from one
  file into another (type checking, CollapseProperties, RemoveUnusedCode
  read across files through scopes and references). Resolving a
  (file, index) handle through a file table measured no faster than
  today's chunks, so global passes should run on the unified store, not on
  per-file stores behind a directory.
- **Switching** is a rebuild: copy the trees into the new store(s) in
  pre-order and forward the holders that live across pass boundaries. That
  list already exists (`node-holder-census.md`, the compaction forwarding
  registry), so the switch reuses compaction, with ids changing at the
  switch only. At ~0.1 s per switch at typed-full scale, switches have to
  be grouped (parse + script-local checks per file → unified for type
  checking and the global optimization loop → per-file for script-local
  late passes and output), not done around every pass.

Open points:

1. **Allocation target.** With per-file stores a new node must be created
   in the file it will be inserted into. Script-local passes know it;
   a global pass creating a node for another file runs on the unified
   store, so this holds as long as cross-file moves stay global.
2. **Parallel allocation in the unified store.** Appending from several
   goroutines needs either per-worker ranges reserved before a parallel
   phase (with growth only between phases) or running such phases on
   per-file stores. The second is the per-file mode itself.
3. **Detached subtrees held across passes** (change timeline, side maps)
   must be copied at a switch like the trees; the census lists them.

Recommendation: do the handle migration on a unified store first (one
mode, the simplest correct design, and what the cross-file passes need),
keeping per-file slices only as parse buffers merged once after parsing.
Add the per-file mode for script-local phases afterwards, when main's
parallel passes are merged and allocation contention or per-file rebuilds
show up in profiles; the pass code does not change for it. Before the
migration, ship the `nodeAt` inline fix (section 4).

## 6. Goal: type-safe, plain Go (2026-10-24)

Direction from the project owner: move towards code that is type safe and
unsurprising, so that a junior developer can read and explore it. Speed
and memory still matter, but an unusual construct has to earn its place.

Rule of thumb for new and existing code:

- No `unsafe`, `uintptr` addresses, `runtime.AddCleanup`/finalizers or
  weak pointers unless a benchmark shows a clear win (several percent on a
  benchmark task) that a plain version cannot get, and then only behind a
  small typed API in one file, with a comment saying why.
- Prefer distinct named types over raw integers (`type NodeID uint32`,
  not `uint32`), so the compiler catches mixed-up handles.
- Prefer plain slices and maps owned by one object (a Compiler, a
  NodeStore) over process-wide tables, so lifetime follows ordinary
  ownership and the GC frees things without help.
- One mode over two: add a second mode (per-file stores) only when a
  profile asks for it.

How the slice design serves this. The NodeStore migration removes the
most unusual code on the branch: the chunk table of untraced addresses,
generations with `AddCleanup` and cross-generation ties, the `unsafe`
handle tables for types, colors and JSDoc, and the GC-lifetime rules in
`node_store.go`/`node_tables.go`. A node's type, color and JSDoc become
typed per-store columns indexed by `NodeID` (`types []*JSTypeBase`,
cleared wholesale at typesToColors; `colors []*colors.Color`;
`jsdoc []*JSDocInfo`), and the store is freed like any other value. The
per-file mode is deferred under the one-mode rule.

Inventory of the remaining unusual code (non-test, at 2d9eaf68), in the
order to clean it up:

| Area | Where | What | Plain replacement |
|---|---|---|---|
| AST storage | `rhino/node_store.go`, `node_tables.go`, `node_compact.go`, `node.go` (aux) | chunk table of `uintptr`, `unsafe.Pointer` handle tables, `AddCleanup` generations, weak reporter in `jscomp/compiler.go` | NodeStore with typed columns (above) |
| JSType casts | `rhino/jstype_*.go` (~50 uses) | interface value → embedded `*JSTypeBase` through `unsafe` offsets | a `base() *JSTypeBase` method on the interface; measure the type-checking cost |
| Scope/var casts | `jscomp/abstract_scope.go`, `abstract_var.go` | `unsafe` cast from base struct to outer struct | generics or an interface method |
| String tricks | `jscomp/scope.go` (hand-written hash), `rhino_string_pool.go`, `jutil/jstring.go`, `code_printer.go`, `compiler.go`, `default_name_generator.go`, `sourcemap/` | zero-copy `unsafe.String`, pointer-equality fast paths, custom hashing | `strings.Builder`, `map[string]`, `hash/maphash`; keep only what a benchmark defends |
| Pools | `sync.Pool` in 10 files | object reuse | delete where the allocation rate no longer shows in profiles (round 1 removed most) |

Each row is its own change with a benchmark before/after and the
byte-identical gate, so a row whose plain version costs too much can be
kept with a comment instead.
