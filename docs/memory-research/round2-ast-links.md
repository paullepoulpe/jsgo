# Round 2: int32 node links, chunked node storage and copying compaction

Status: implemented (2026-10-23) on the `ast-links` worktree; measurements
below. Companion to `round2-design.md` (design) and `node-holder-census.md`
(the holder list the forwarding registry implements). The type/color handle
(`typeOrColor` → 32-bit `tc` + `aux`) is the concurrent `ast-handles` work;
until the two are merged `Node` still carries one Go pointer, so the chunks
are still scanned (see "What is and is not pointer-free yet").

## 1. Node storage and the id scheme (`internal/rhino/node_store.go`)

- `Node.next/first/parent/previous` are `nodeID` (uint32); a new `self`
  field holds the node's own id. `Node` is 56 bytes (was 64): 4 × 4 B links
  + 4 B self + 8 B typeOrColor + payload 8 + propFlags 8 + lineno 4 +
  length 2 + token 1 + kindFlags 1.
- Chunks: `nodeChunk = { header (one node wide); nodes [4095]Node }`,
  4096 × 56 B = 224 KiB (28 pages, no size-class waste). Slot 0 is the
  header (`owner *nodeGen`, `index`, bump counter), so **id = chunk<<12 |
  slot with slot in 1..4095 and id 0 = nil**. 2^20 chunks of id space
  (4 G nodes per process; the last chunk index is reserved for the
  forwarding marker).
- Directory: one process-wide table `nodeChunkTable [2^20]uintptr` (8 MiB
  of BSS, demand-paged, **not traced by the GC**: the slots are integers).
  `nodeAt(id)` = `(*Node)(chunkTable[id>>12] + (id&4095)*56)`: an atomic
  load of the slot and no bounds check. A cleared slot panics
  ("stale node id") instead of dereferencing.
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
  lifetime grouping: ids are global and links across generations are valid
  while both are alive. Limit (documented, tests only): a node that outlives
  the Compiler that made it and links into a generation that was since
  released panics on traversal; a window of a few ms exists between a
  generation's collection and its cleanup in which such a (buggy) traversal
  reads freed memory.
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
- `Finish` clears the old chunks' table slots immediately and drops the
  generations' chunk lists (a stale holder keeps only its own chunk alive).

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
Enabled with `--fast_mode`; `JSGO_AST_COMPACT=0|1` overrides (1 = also in
default mode), `JSGO_AST_COMPACT_LOG=1` logs, `JSGO_AST_POISON=1` is the
verifier.

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
  nested shadow compiler compacting a tree the outer compilation holds).

## 3. Measurements

(filled in below)

## 4. What is and is not pointer-free yet

`Node` still holds `typeOrColor unsafe.Pointer`, and the chunk header holds
`owner *nodeGen`, so a chunk is a pointerful 224 KiB object and the GC scans
all of it (the links themselves are no longer pointers: 1 pointer word per
node instead of 5). Once the handles work lands, `Node` has no pointer and
the chunk type's pointer prefix is the 8-byte `owner` only: Go's allocator
accounts and scans objects up to `PtrBytes`, so the AST then costs 8 bytes
of scan work per chunk. The chunk header is laid out for that (the pointer
is at offset 0).

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
