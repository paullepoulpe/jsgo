# Data-oriented compilers → jsgo memory layout (cluster: data-oriented-compilers)

Projects studied (source read from shallow sparse clones under `research/src/{zig,carbon,rust,go}`,
plus PR threads): Zig self-hosted compiler, Carbon toolchain, rustc, Go `cmd/compile` (+ LLVM
idioms only where they translate to Go). jsgo code read: `internal/rhino/node.go`,
`rhino_string_pool.go`, `jstype_*.go`, `internal/jutil/{slab,collections}.go`,
`internal/jscomp/{scope,var,graph_linked_directed_graph,compiler}.go`, `docs/gclab-2026-10-18.md`,
`docs/performance-roadmap.md`. Nothing in /home/user/jsgo was modified, built or run.

## 0. Summary (top items)

1. **Pointer API, integer guts (the Go-SSA hybrid).** All four compilers use 32-bit indices
   instead of pointers, but the three that use indices end to end (Zig, Carbon, rustc) were
   designed that way. Go's own SSA backend is the precedent for a *mutable* IR: `*Value`
   pointers in the API, plus dense `int32` IDs so that side tables are slices. For jsgo:
   keep `*rhino.Node` as the pass API, but make the `Node` struct pointer-free (int32
   links, a self-id, a 32-bit type/color/ext handle) and store nodes in fixed-size chunks
   behind a global chunk directory. The node shrinks from 64 B (40 B of pointers) to 48 B
   with no pointers. Estimate on typed-full: −40 MB live, about −65 MB RSS, and about 12.5 M
   fewer pointers to trace on every GC cycle.
2. **Index-based mutable IRs need explicit reclamation, and none of the studied compilers
   gets it from a GC.** Zig and Carbon IRs are append-only and freed whole (per file or per
   function). rustc compacts and renumbers MIR inside a closed `Body`. Go SSA frees a value
   only when `Uses==0`. jsgo needs a mark/sweep of the node store at pass boundaries (with a
   root registry and a poisoning debug mode). Without it, a never-freed store grows with
   allocation volume: gclab measured 2.6× more node bytes allocated than live on typed-app,
   so R1 alone would *raise* peak heap on churn-heavy builds.
3. **Go-specific arena rule.** A pointerful slab chunk is kept alive by any interior pointer.
   Its *dead* elements are still scanned, so they keep their referents alive (transitive
   retention). Arenas only help peak heap when elements are pointer-free or die together.
   Audit `Slab[rhino.Node]` (parse), `scopeArena`, graph chunks and the JSDoc slabs on this
   basis, and measure "reachable from root" against "retained".
4. **Flat immutable arrays for immutable trees** (Carbon post-order tree, Zig ZIR). These do
   not fit the main AST, which Closure mutates constantly. They do fit JSDoc type expressions
   (lazily materialized behind `JSTypeExpression.GetRoot()`), prebuilt externs, and the
   parser's tokens and positions (8 B tokens; `SourceRange` as two `uint32` offsets; line and
   column computed lazily from a line table).
5. **Cheap guards that lock in gains.** Go and rustc both pin hot-type sizes in tests. Go's
   `freePackage` also uses a finalizer to verify that a released structure really is
   collected. jsgo can do the same for types after `typesToColors`, and assert "no pointers"
   for every noscan element type.

---

## 1. Per project: the techniques that matter

### 1.1 Zig (ziglang/zig master, `lib/std/zig/Ast.zig`, `Zir.zig`, `src/Air.zig`, `src/InternPool.zig`, `src/Air/Liveness.zig`, `src/Air/Legalize.zig`)

**(a) Struct-of-arrays tokens and nodes with u32 indices**
([Ast.zig](https://github.com/ziglang/zig/blob/master/lib/std/zig/Ast.zig))
- `TokenList = MultiArrayList(struct{ tag: Token.Tag (u8), start: u32 })` = **5 B/token**.
  There is no end offset (it is re-lexed when needed) and no line or column (computed from
  the source).
- `Node = { tag: u8, main_token: u32, data: Data (8 B) }` = **13 B/node** in SoA, with no
  padding. There are no parent, sibling or position fields. `Data` is an untagged union of
  pairs: `node_and_node`, `opt_node_and_opt_node`, `extra_range`, `node_and_extra`, and so
  on. The node tag decides which variant applies, and a comptime assert keeps `Data` at 8 B.
- `MultiArrayList` puts all columns in one allocation, sorted by alignment.
- Capacity is estimated up front: "8:1 ratio of source bytes to token count" and "2:1 ratio
  of tokens to AST nodes" (`Ast.parse`). One reservation avoids the growth copies.
- Measured in PR [#7920](https://github.com/ziglang/zig/pull/7920) (`zig fmt` over std): peak
  memory 84.7 → 68.4 MiB (−19.3%), wall −11.4%. Parsing alone was −22% wall and −28%
  instructions. Andrew Kelley's "Practical DOD" talk (secondhand summary on
  [forums.swift.org](https://forums.swift.org/t/what-are-the-benefits-of-minimizing-the-memorylayout-size-stride-of-a-struct/74205/3))
  reports the follow-up ZIR rework at −40% peak memory and −39% wall.

**(b) "extra" side array for variable-sized payloads, plus specialized tags for common shapes**
- A node with more than two operands stores an `ExtraIndex` into one `extra_data: []u32`. A
  `SubRange{start,end}` points at a run of child indices there. There is no per-node slice
  allocation and no hash table.
- Common shapes get dedicated tags so they fit inline: `block_two` (at most two statements,
  stored inline) versus `block` (an `extra_range`); likewise `call_one` and similar.
- The parser collects children on a reusable `scratch` stack, then copies them into
  `extra_data` in one go.
- ZIR and AIR follow the same scheme: `instructions: MultiArrayList(Inst{tag u8, data 8 B})`,
  `extra: []u32`, and `string_bytes: []u8` with `NullTerminatedString = enum(u32)`.
  `extraData(T, idx)` decodes a struct of u32 fields by reflection.

**(c) One index space for constants and instructions (`Ref` encoding)**
- In ZIR, a `Ref` below `InternPool.static_len` is a well-known interned type or value
  (`u8_type`, `bool_true`, and so on). Anything above that is an instruction.
- In AIR, the most significant bit of a `Ref` says "instruction index" versus "InternPool
  index".
- "None" is a sentinel (`maxInt`). Optional indices therefore cost no extra byte.

**(d) InternPool: every type and value, hash-consed into one pool with 32-bit ids**
- `Item = { tag: u8, data: u32 }` (SoA), plus `extra: []u32`, `limbs`, `strings`, and
  `string_bytes`.
- The hash map stores **only `{value: u32, hash: u32}` = 8 B entries**. Keys are rebuilt from
  `items`/`extra` for comparison (adapted hashing), so the map holds no copies of keys.
- Multi-threaded: one `Local` per thread (append-only lists plus an arena) and the thread id
  in the top bits of each index, so writers never contend on a shared append point. The maps
  are sharded and readers acquire them lock-free.
- Measured in PR [#15569](https://github.com/ziglang/zig/pull/15569) ("Use InternPool for all
  types and constant values"): ghostty maxrss 1130 → 463 MiB. Other projects moved −3% to
  +4% in memory and ±10% in wall time. **Interning wins only where duplication is high.**

**(e) Lifetimes tied to compilation units; mutation by in-place slot overwrite**
- ZIR is "the minimum amount of information needed": once AstGen is done, the AST, tokens
  and source are not needed. ZIR is cached on disk as its raw pointer-free arrays.
- AIR is per function, and so is its `Liveness`. The `Liveness` tomb bits are 4 bits per
  instruction ("unreferenced" plus 3 operand-death bits), with a sparse `special` hash map
  for the rare instructions that need more.
- `Air/Legalize.zig` rewrites by **overwriting an instruction's slot in place**
  (`replaceInst`), so every existing reference to that index stays valid. New instructions
  are appended, and sequences are spliced in as `block` bodies stored in `extra`. Superseded
  `extra` data is left behind as garbage until the whole function's AIR is freed.

### 1.2 Carbon (carbon-language/carbon-lang trunk, `toolchain/{lex,parse,sem_ir,base,docs}`)

**(a) 8-byte tokens and a line table**
([token_info.h](https://github.com/carbon-language/carbon-lang/blob/trunk/toolchain/lex/token_info.h),
[tokenized_buffer.h](https://github.com/carbon-language/carbon-lang/blob/trunk/toolchain/lex/tokenized_buffer.h))
- `TokenInfo` = 8 B: an 8-bit kind, a 1-bit "has leading space", a 23-bit payload
  (identifier id, int id, or the matching bracket's token index), and an `int32` byte
  offset. The byte offset is "hot while lexing but then generally cold".
- `TokenIndex::Bits = 23`, so a file is limited to 8 M tokens.
- `LineInfo{start, indent}` per line. Line and column are derived by searching the line
  table.
- Out-of-band booleans: `recovery_tokens_` is a `BitVector`.
- The lexer keeps an upper bound on the number of parse nodes per token kind
  (`expected_max_parse_tree_size`), so the parser reserves its storage exactly once.

**(b) Flat post-order parse tree**
([docs/parse.md](https://github.com/carbon-language/carbon-lang/blob/trunk/toolchain/docs/parse.md),
[tree.h](https://github.com/carbon-language/carbon-lang/blob/trunk/toolchain/parse/tree.h))
- `NodeImpl` = 8 B: kind, a `has_error` bit, and the token index (`static_assert(sizeof == 8)`).
  The tree is a single `SmallVector<NodeImpl>` in depth-first post-order.
- Structure is implied by "bracketing" nodes (the opening child) or by a fixed child count
  per node kind.
- Subtree sizes **were dropped from the node**. They are computed on demand
  (`TreeAndSubtrees`), and the doc says it may "eliminate it if it turns out to be
  unnecessary and a meaningful cost".
- Checking consumes the tree in post-order with an explicit node stack, with no recursion.
  The design rules in docs/README.md are "Distinct steps", "Vectorized storage" (flyweight
  int32 handles instead of pointers) and "Iterative processing".

**(c) SemIR value stores**
([value_store.h](https://github.com/carbon-language/carbon-lang/blob/trunk/toolchain/base/value_store.h),
[inst.h](https://github.com/carbon-language/carbon-lang/blob/trunk/toolchain/sem_ir/inst.h))
- `Inst` = 16 B `{kind, type_id, arg0, arg1}`, all `int32`. Locations live in a **parallel
  vector** (`loc_ids_`, struct-of-arrays), as do constant values.
- `ValueStore<IdT, T>` is chunked: **4 KiB chunks**, id = `chunk << bits | pos`. References
  stay stable and growth never copies old data. The code notes it is "up to 5%" slower than
  one contiguous vector.
- `CanonicalValueStore` provides deduplication with ids (ints, floats, names, types).
- `BlockValueStore` slab-allocates id arrays ("blocks") and can canonicalize them.
- `IdTag` uses spare high bits of an id to tag which store it belongs to, checked in debug
  builds.

**(d) Mutation**
- The IR is append-only. A rewrite replaces a placeholder instruction with a `splice_block`
  that points at a new block, or holds the instructions aside until their final shape is
  known (docs/check/README.md "Rewrites").
- Nothing is ever deleted. Everything is freed with the file.

### 1.3 rustc (rust-lang/rust master)

**(a) Hash-consed types in a session arena**
([ty/context.rs](https://github.com/rust-lang/rust/blob/master/compiler/rustc_middle/src/ty/context.rs),
[intern.rs](https://github.com/rust-lang/rust/blob/master/compiler/rustc_data_structures/src/intern.rs))
- `CtxtInterners` has one `ShardedHashMap<InternedInSet<T>, ()>` per kind (types, args, type
  lists, predicates, consts, layouts, and so on).
- `intern_ty` allocates `WithCachedTypeInfo{internee, flags, outer_exclusive_binder}` in the
  arena. The flags are **precomputed summaries**, so "has params" style queries never walk
  the type.
- `Ty<'tcx>` is a pointer-sized handle compared by address (`Interned<T>`).
- `List<T>` is a length-prefixed, thin-pointer slice, also interned.
- `TypedArena` and `DroplessArena` grow their chunks from 4 KiB to 2 MiB and are freed with
  the session. Nothing is freed per object.

**(b) Index newtypes and owner-local ids**
- `newtype_index!` defines a `u32` with a reserved niche, so `Option<Idx>` is still 4 B.
  `IndexVec<I, T>` is a vector indexed by that type.
- `HirId = { owner: OwnerId(u32), local_id: ItemLocalId(u32) }`.
- `OwnerNodes.nodes: IndexVec<ItemLocalId, ParentedNode{ parent: ItemLocalId, node }>`. **Parent
  links are dense local ids in a side vector**, not pointers stored in the nodes.
- `TypeckResults` is one table per owner (function): `ItemLocalMap<Ty>` for node types,
  adjustments, field indices, and so on. Side tables are keyed by small local ids and are
  dropped with the owner's results.

**(c) Compact spans** ([span_encoding.rs](https://github.com/rust-lang/rust/blob/master/compiler/rustc_span/src/span_encoding.rs))
- An 8-byte `Span` uses one of four inline bit layouts and falls back to an interner index
  for rare cases. Quote: "An earlier version … used only 4 bytes for Span, but that was
  slower because many fewer spans could be stored inline."

**(d) Data-flow sets**
([bit_set.rs](https://github.com/rust-lang/rust/blob/master/compiler/rustc_index/src/bit_set.rs))
- `MixedBitSet` is `DenseBitSet` for small domains and `ChunkedBitSet` for large ones.
- `ChunkedBitSet` chunks are `Zeros`, `Ones`, or `Mixed(Rc<[Word]>)`. Mixed chunks are
  **shared copy-on-write** "because it's surprisingly common to duplicate an entire chunk".

**(e) Mutation of index IRs**
([patch.rs](https://github.com/rust-lang/rust/blob/master/compiler/rustc_mir_transform/src/patch.rs),
[simplify.rs](https://github.com/rust-lang/rust/blob/master/compiler/rustc_mir_transform/src/simplify.rs))
- MIR is mutated in place. `MirPatch` queues new statements, blocks and terminator
  replacements and applies them all at once.
- Deleted statements become `Nop` tombstones.
- `remove_dead_blocks` and `SimplifyLocals` **compact and renumber**: a replacement map plus a
  `MutVisitor` that rewrites every index. This is safe because every holder of a block or
  local index lives inside the `Body`, a closed world.

**(f) Size discipline**
- `static_assert_size!` pins the sizes of hot AST and HIR types.
- PR [#37445](https://github.com/rust-lang/rust/pull/37445) shrank `Expr_` from 144 to 64 B,
  about −1% peak RSS on one workload. Later PRs switched AST vectors to `ThinVec` (8 B
  instead of 24 B).

### 1.4 Go cmd/compile (go master) and LLVM idioms

**(a) SSA values: pointer API plus dense ids**
([ssa/value.go](https://github.com/golang/go/blob/master/src/cmd/compile/internal/ssa/value.go),
[func.go](https://github.com/golang/go/blob/master/src/cmd/compile/internal/ssa/func.go),
[cache.go](https://github.com/golang/go/blob/master/src/cmd/compile/internal/ssa/cache.go))
- `Value` is 112 B and `Block` is 312 B on 64-bit (pinned by `sizeof_test.go`).
- `ID int32` is allocated densely, so analyses use `[]T` indexed by ID instead of maps.
- `Argstorage [3]*Value` is inline small-vector storage that `Args` points into.
- `FreeValue` requires `Uses == 0`. The freed value keeps its ID, goes on a free list linked
  through `Argstorage[0]`, and is reused by `NewValue`.
- `Cache` holds preallocated `[2000]Value`, `[200]Block` and `[2000]Location`, reused across
  functions. `Reset()` clears only the used prefix, found by binary search on `ID != 0`.
- `allocators.go` (generated) keeps a `sync.Pool` per power-of-two size class for slices.
- `SparseSet` and `SparseMap` (Briggs–Torczon) are pooled per function.

**(b) Compact positions and kind-specific node structs**
- `src.XPos` is 8 B `{index int32 into a PosBase table; lico uint32 = 20-bit line + column +
  statement bits}`.
- `ir.miniNode` is 12 B. Each op has its own struct, sized for that op.
- `types.Type` is 96 B, with `extra any` for kind-specific data. Derived types are memoized
  on the element type, so the pointer type and slice type are created once.

**(c) Release as soon as possible**
- After machine code is generated: `fn.Body = nil; fn.Dcl = nil; fn.ClosureVars = nil`
  ([gc/compile.go](https://github.com/golang/go/blob/master/src/cmd/compile/internal/gc/compile.go)).
- The unified IR noder writes the type-checked syntax to export data, reads the IR back
  lazily (`todoBodies`), and sets `fn.Body = nil // free memory` for imported inline
  candidates.
- `freePackage` zeroes the `types2.Package`. Under `-d=gccheck=1` it **sets a finalizer and
  forces up to 10 GCs to verify the package is actually collected**
  ([noder/unified.go](https://github.com/golang/go/blob/master/src/cmd/compile/internal/noder/unified.go)).

**(d) LLVM idioms and how they translate to Go**
- `BumpPtrAllocator` becomes typed slabs only. Go cannot bump-allocate pointerful data
  untyped, and the `arena` GOEXPERIMENT is on hold.
- `SmallVector<T,N>` becomes an inline `[N]T` plus a slice into it (jsgo's
  `UnionType.inlineAlternates` and graph-node buffers already do this). Caveat: the struct
  must not be copied by value.
- `DenseMap` (one flat open-addressed array with empty/tombstone keys) becomes an `int32`
  open-addressing index (already in flight for `LinkedHashMap`).
- `PointerIntPair`, tagged pointers and hung-off `Use` lists do not translate (see §4).

---

## 2. Go constraints that shape every recommendation

- A `make([]T, n)`, map or chunk is **noscan exactly when T has no pointers**. No `unsafe` is
  needed. This holds for maps too: when both key and value are pointer-free, the map's
  internal arrays are noscan (both the Swiss and classic implementations).
- Noscan bytes still count toward the heap goal, so RSS ≈ live × (1 + GOGC/100) at peak.
  Fewer **bytes** lower RSS; fewer **pointers** lower GC CPU. On typed-full, RSS/live =
  879/552 ≈ 1.6.
- A pointer to one element keeps the whole chunk alive. Worse, if the element type has
  pointers, **dead elements of a live chunk are still scanned and keep their referents
  alive**. A detached subtree whose nodes share parse slabs with live code keeps its
  children, JSDoc and types reachable. The arenas in Zig, Carbon and rustc are free of this
  only because they have no tracing GC and their lifetimes are phase- or unit-scoped.
- Go has no thread-local storage. A zero-argument accessor like `n.GetNext()` can only
  resolve an id through a **global** structure, such as a chunk directory with chunks
  installed atomically. Parallel allocation should reserve whole chunks per goroutine
  (Zig's per-thread `Local`, with a tid in the high bits).
- The GC is non-generational, so long-lived pointers (AST, types, cached scopes) are
  re-traced on every cycle. Pass-local churn costs mostly `mallocgc` (12–15% of CPU per
  gclab).
- Interior pointers and fixed-size arrays let `nodeChunks[id>>12][id&4095]` compile without
  bounds checks when the directory is a fixed-size array and the indices are masked.

---

## 3. Recommendations for jsgo, ranked by expected peak-heap / pointer impact

Numbers use the context's measurements: typed-full peak 552 MB live and 879 MB RSS, with the
AST at about 160 MB, roughly 2.5 M nodes at 64 B. Untyped builds have 76–89% of live heap in
the AST (gclab, bundle). All gains are estimates.

### R1. Pointer-free 48 B `rhino.Node` in a chunked global node store; `*Node` API unchanged
*Sources: Go SSA (pointer API plus dense ids), Carbon `ValueStore` chunks, Zig and rustc u32
indices, Zig per-thread locals.*

**Layout.** 48 B with zero pointers:

```go
type NodeID int32 // 0 = nil
type Node struct {
    id                         NodeID // self, so a *Node can be stored as a link (4 B)
    next, first, parent, previous NodeID // 16 B; previous keeps Closure's circular first.previous = last
    typeOrColor                uint32 // color id | JSType id | ext-box id, discriminated by kindFlags as today
    payload                    uint64
    propFlags                  uint64
    linenoCharno               int32
    length16                   uint16
    token                      Token
    kindFlags                  uint8
} // 48 B, no pointers -> []Node chunks are noscan
```

**Store.** `var nodeChunks [1<<14]atomic.Pointer[[4096]Node]` gives 67 M nodes, more than
the 35 M-node googleapis case. Each chunk is 192 KiB.

```go
func at(id NodeID) *Node {
    return &nodeChunks[uint32(id)>>12&(1<<14-1)].Load()[id&4095]
}
```

`GetNext()` becomes `if n.next==0 {return nil}; return at(n.next)`.
- Chunk numbers come from a global allocator, so concurrent compilations (tests) share the
  directory without a per-node table pointer.
- Each allocating goroutine reserves a whole chunk at a time. Slots within it are plain bumps.
- When a compilation ends, its chunk numbers are recycled.
- jsgo already uses this exact pattern for `sourceFilePropChunks`.

**The `typeOrColor` slot.** Two variants:
- **(a) All ids.** Colors already live in an interned pool, so they get dense ids. JSTypes get
  ids from a per-compilation `[]JSType` table, assigned in the registry's alloc slabs
  (`SerialGen` is already there). JSDoc and ext boxes get ids into tables.
  - Caveat: the type table keeps every type ever assigned to a node alive until it is dropped
    at `typesToColors` or `removeTypes`.
- **(b) Hybrid.** Colors and ext boxes go by id. JSTypes go in a separate per-chunk
  `[4096]unsafe.Pointer` column, allocated lazily (pointerful, 8 B per node, only in typed
  builds). The whole column is dropped at `typesToColors`.
  - Types then keep their natural Go-GC lifetime, and untyped builds pay nothing.

Option (b) carries less risk.

**What moves.**
- `nodeExtraShards` (a global map keyed by `uintptr(unsafe.Pointer(n))`) becomes a paged
  column indexed by `NodeID`.
- `NewNode` and the 16 files that build Node values or use `Slab[rhino.Node]` (the parse,
  JSDoc and deserializer paths, plus a few ES6 passes) switch to `store.alloc()`.
- `cloneTreeWithSlab` allocates contiguously, so clones are laid out in pre-order.

**Gain.**
- AST −16 B/node: typed-full ≈ 160 → 120 MB (−40 MB, about −7% of peak, about −65 MB RSS).
  Untyped builds −19–22% of live.
- About 100 MB of pointer-bearing AST bytes (12.5 M pointers) leave GC marking. Per gclab,
  that is the AST's 45–60% (typed) and about 85% (untyped) share of scan work.
- Side tables keyed by `NodeID` become noscan (R6).

**Passes.** 1:1. Only rhino substrate and allocation sites change. **Byte identity:**
unaffected (pure substrate). **Effort:** L. **Mode:** default.

**Hard prerequisite: R2.** gclab measured about 90 MiB of nodes allocated against 35 MiB
live on typed-app. A store that never frees ends up at roughly 2.6 × 48/64 ≈ 2× today's AST
bytes on churn-heavy (untyped, optimization-loop-peak) builds.

### R2. Explicit node reclamation: mark/sweep at pass boundaries; pre-order compaction later
*Sources: Go SSA `FreeValue`/free list with `Uses` invariants; rustc closed-world compaction
and renumbering (`remove_dead_blocks`, `SimplifyLocals`); Zig per-function lifetime; ASan-style
quarantine.*

**Mark.** Run at `PhaseOptimizer` boundaries, triggered like GOGC: when nodes allocated since
the last sweep exceed X% of the last live count. Roots:
- `jsRoot` and `externsRoot`, plus `scriptNodeByFilename` and the CompilerInput `JsAst` roots.
- Node-valued properties, such as the closure-unaware shadow.
- Every `JSTypeExpression` root reachable from the JSDoc of marked nodes.
- While types are alive: `FunctionType.source`, `Property.propertyNode`, and TypedVar name
  nodes. Simplest policy: sweep only after `typesToColors` in typed builds, and always in
  untyped builds. Typed-full's peak is during type checking anyway, when churn is low.
- A small registry for the few other cross-pass holders.

**Sweep.** Push unmarked slots onto a free list. The same mark also frees ext-box and JSDoc
table entries (and type ids in variant R1a).

**Safety.** The only hazard is a `*Node` to a *dead* node held across a boundary. Live nodes
never move, so holders of live nodes are unaffected.
- In a `jsgo_debug` build, freed slots are poisoned (`token = FREED`), held in quarantine for
  K sweeps, and every accessor asserts on them.
- Run the full parity corpus in that mode.

**Cost.** One AST walk, roughly 10–30 ns/node, so 25–75 ms per sweep on typed-full.
**Gain.** Reclaims dead nodes precisely. That includes garbage that today's 128-node
pointerful parse slabs retain transitively, which current heap profiles count as "live".

**Optional second step.** Evacuate live nodes into fresh chunks in pre-order. gclab measured
walks 2.5× faster on contiguous pre-order layouts. This moves live nodes, so *every*
cross-boundary holder must be remapped, which is rustc's closed-world requirement. Do it only
for the optimization loop, after types are gone, and only once the R2 root registry has
proven complete.

**Passes.** 1:1. **Effort:** M. **Byte identity:** safe if the root set is complete; the
debug mode plus the oracle catch omissions. **Mode:** default (mark/sweep); compaction
behind a flag until proven.

### R3. Arena and slab hygiene: Go-specific rules, plus one measurement first
*Source: the contrast between Zig, Carbon and rustc arenas (unit-scoped, no tracing GC) and
Go's interior-pointer and dead-element retention.*

**Rules.**
1. Use a pointerful `Slab[T]` only for objects that die together, and `clear()` an element
   when it is logically freed.
2. Long-lived or mixed-lifetime slabs must hold pointer-free `T`.
3. Never keep a slab across passes (in flight: removal of cross-pass pools).

**Audit:**
- `Slab[rhino.Node]` (parse and JSDoc; mixed lifetimes after parse)
- `scopeArena` (`[]Scope`, `[]Var`, `[]*Var` chunks)
- `LinkedDirectedGraph` node and edge chunks
- `JSDocInfo` `InfoSlab`, `TypeExprSlab`, `AnySliceSlab`
- the registry type slabs

**Measurement (cheap, decisive for R2).** At the peak point and mid-loop, compare:
- nodes reachable from roots × 64 B, against
- heap-profile in-use bytes for node allocation sites.

The difference is retained garbage.

**Effort:** S–M. **Mode:** default.

### R4. JSDoc type expressions as flat, immutable pre-order arrays, materialized lazily
*Sources: Carbon flat post-order tree, Zig `extra` and `string_bytes`.*

**Encoding.**
- Store each `JSTypeExpression` as one run per expression in a per-file `[]uint32` "extra"
  array: `{token u8 | childCount u8 | flags u16}`, `stringID u32`, and a packed
  `linenoCharno` for each node. That is about 12 B/node instead of 64 B, with no pointers.
- `GetRoot()` (51 call sites in 25 files) builds the `*Node` tree on first use and caches it.
  Module rewriting mutates JSDoc type nodes in place, so a materialized tree must stay
  materialized.
- Type checking reads most expressions once. After checks, most are never touched again.

**Gain.** JSDoc is 10% of typed-full's peak (about 55 MB). If half of that is type-expression
nodes, the gain is about −20 MB and about 430 k fewer pointerful nodes. Externs are heavily
typed, so the effect also shows on untyped builds that still parse externs' JSDoc.

**Passes.** 1:1 (behind `JSTypeExpression`). **Byte identity:** materialization must
reproduce positions and properties exactly; the oracle covers this. **Effort:** M.
Coordinate with the JSDoc-compaction swarm. **Mode:** default.

### R5. Parser tokens and positions in the Zig/Carbon model (coordinate with the parser swarm)
- **Tokens.** Carbon-style 8 B `{kind u8, flags, 23-bit payload (StringID or literal-table
  index), u32 offset}`, or Zig's 5 B `{tag, start}` in SoA.
- **Locations.** `ParseTreeBase.Location` becomes a value `{start, end uint32}` (UTF-16
  offsets, 8 B, no pointers). Today it is a `*SourceRange` (16 B) pointing to two
  `*SourcePosition` (32 B each).
- **Line and column.** Computed from a per-file `[]uint32` line-start table. IRFactory
  creates nodes in near-source order, so a monotone cursor gives O(1) amortized
  `linenoCharno`.
- **Reserve once.** Size node and token chunks from the source length (Zig: 8 bytes per
  token, 2 tokens per node; or Carbon's per-token upper bound).
- **Source text.** Store it as Latin-1 `[]byte` when every code unit is below 0x100,
  otherwise UTF-16. This is what Java's own compact strings do. UTF-16 indices equal byte
  indices, so offsets and charno are unchanged.

**Gain.** For the single-file case (8.9 MB, parse-time peak): text 17.8 → 8.9 MB, and
locations about 64–80 B → 8 B per parse tree. **Passes:** parser-internal; Rhino-facing APIs
unchanged. **Byte identity:** offsets preserved by construction. **Effort:** M. **Mode:**
default.

### R6. Dense ids as collection keys: noscan side tables for pass-local maps
*Sources: rustc `IndexVec` and `ItemLocalMap`, Go SSA ID-indexed slices and sparse sets,
Zig index-only hash entries.*

- **Node-keyed maps.** Add `rhino.NodeMap[V]`, `NodeSet` and `NodeLinkedHashMap[V]`, keyed by
  `NodeID` internally, with the same method sets as the `jutil` maps.
  - Replace the 64 `map[*rhino.Node]` and 122 `jutil.*[*rhino.Node, …]` declarations in
    jscomp mechanically.
  - With scalar values the storage becomes noscan, and keys shrink from 8 B to 4 B.
  - A per-pass dense table can be a paged array (one page per 4096-node chunk, allocated on
    touch).
- **Insertion order.** A Briggs–Torczon sparse map iterates in insertion order only if
  nothing is removed. Removal swaps the last element into the hole, so use tombstones to keep
  `LinkedHashMap` semantics.
- **String-keyed maps.** For maps keyed by names, offer `StringID` keys where the pass already
  has the node (`n.GetStringID()`). Keep a pool-side column of Java `String.hashCode` per
  `StringID` so `JavaHashMap` ordering is computed without touching bytes. This also trims
  the ~5% `jutil.StringHashCode` CPU.

**Gain.** Mostly churn and scan of pass-local structures (31–40% of allocated bytes are
collection entries), with a smaller effect on peak. **Passes:** declarations change; logic
stays 1:1. **Effort:** M (mechanical). **Mode:** default.

### R7. Data-flow and per-function scratch (refines the in-flight CFG work)
- **Lattices.** Use rustc's `ChunkedBitSet`: Zeros/Ones/Mixed chunks with **copy-on-write
  shared `[]uint64` chunks**. The Go GC handles the sharing, so no `Rc` is needed.
  `MustBeReachingVariableDef`-style per-node copies then share unchanged chunks.
  - Dense-but-small variable sets use one flat bitset.
  - This is the 357 MiB site in `flowSensitiveInlineVariables`.
- **Per-function scratch.** Model it on Go's `ssa.Cache`: preallocated pointer-free arrays,
  reset per function by clearing only the used prefix.
  - Per-function analysis data (CFG, lattices, liveness) then lives only as long as that
    function. That is Zig's per-function AIR and `Liveness` lifetime, and Go's
    `fn.Body = nil` after codegen.
  - Liveness-style facts as packed bits (Zig: 4 tomb bits per instruction plus a sparse map
    for exceptions).

**Passes:** lattice internals only. **Byte identity:** safe. **Effort:** M. **Mode:** default.

### R8. Memoize derived union types; intern only what is duplicated
*Sources: Go `types.Type` memo fields (the pointer type cached on its element), rustc
interners, Zig InternPool numbers.*

- `CreateNullableType` and `CreateOptionalType` build a fresh `UnionType` (128 B) on every
  call, and `?T` / `T|undefined` are ubiquitous. A memo slot on the base type (or a small
  registry map) for `T|null` and `T|undefined` would remove the duplicates.
- **First measure** how many unions are structurally equal at the typed-full peak. Zig's
  InternPool cut ghostty from 1130 to 463 MiB but did nothing for low-duplication projects.
- **Identity audit needed.** Java creates distinct objects, and `JSTypeBase` holds per-object
  traversal state (`visited`, `hashCodeInProgress`, `templateCheckResult`). Sharing can change
  cycle-detection paths.
- Equality and hash are structural, so `HashMap` order is unaffected.
- **Mode:** default if the audit is clean, otherwise `--fast_mode`. **Effort:** S–M.

### R9. Size and pointer regression guards
*Sources: Go `sizeof_test.go`, rustc `static_assert_size!`, Go `freePackage` with gccheck.*

- Extend `node.go`'s compile-time 64 B assert to all hot types (JSType structs, Var, Scope,
  Reference, `LhmSlot` instances, graph nodes, store element types).
- Add a reflection test asserting "**no pointer fields**" for every type meant to be noscan:
  Node in R1, R4 encodings, lattice chunks.
- For "release early" work (types after colors, parse arenas, source text), add a debug flag
  that zeroes the root object, attaches `runtime.AddCleanup`, forces up to 10 GCs, and fails
  if the object was not collected. This catches a stale cache or pool reference that silently
  keeps a whole structure alive.

**Effort:** S. **Mode:** tests and debug only.

### R10. (After R1 and R2) An off-heap node store
- Put node chunks in `mmap`ed memory. Pointers into non-heap memory are legal and ignored by
  the GC.
- The AST then no longer inflates the heap goal: RSS ≈ AST + (1 + GOGC/100) × rest, instead
  of (1 + GOGC/100) × (AST + rest). At GOGC=100 on typed-full that is about −120 MB RSS.
- Requires R2's explicit lifetimes and an `munmap` at compilation end. A stale `*Node` after
  unmapping faults instead of corrupting memory.
- Needs `unsafe`, and `-race`/checkptr compatibility must be checked.
- **Mode:** behind a flag first. **Effort:** M.

### R11. Prebuilt default externs as a flat snapshot (low memory impact)
- Precedent: Zig ZIR is cached as raw pointer-free arrays.
- Embed the default externs as a pre-order array (token, StringIDs, flags, positions,
  JSDoc runs) and materialize them by a linear copy into the node store.
- Removes the externs parse churn (positions, parse trees) for every compilation. Peak heap
  after load is unchanged.
- **Byte identity:** must equal a parse (checked by a test comparing trees). **Mode:**
  default. **Effort:** M.

---

## 4. Fast-mode-only ideas, and ideas considered and rejected

**Possible in `--fast_mode` (deterministic but not Java-identical):**
- **Union and record interning without the identity audit** (R8 taken further). Hash-cons all
  structural types after resolution, as Zig's InternPool and rustc do.
- **Dropping JSDoc type-expression trees right after type checking** in ADVANCED. Only valid
  if no later pass reads them; that needs an audit. Without the audit, it is fast-mode only.
- **Go maps keyed by ids with Go iteration order** where the default mode emulates Java
  `HashMap` order. This saves the emulation's slot arrays.

**Rejected:**
- **Full SoA (`MultiArrayList`) for `rhino.Node`.** It breaks the `*Node` struct API. Zig's
  gain is padding (13 B vs 16 B), and jsgo's Node already has none. Sparse optional columns
  (R1b types, `nodeExtra`) capture the only real gain.
- **Flat post-order immutable tree (Carbon) for the main AST.** Closure detaches, replaces and
  inserts constantly, and needs O(1) `getParent`, `getPrevious` and `getLastChild`. It is used
  only for immutable data (R4, R5, R11).
- **Dropping `parent` and `previous` links** (Zig and Carbon have none). That would turn
  O(1) Closure APIs into O(n).
- **23-bit (Carbon) or relative 16-bit indices.** jsgo's builds reach 35 M nodes, so ids must
  be 32-bit.
- **Tagged pointers, `PointerIntPair`, low-bit stealing.** They stay pointers, so the GC still
  scans them. checkptr flags them, and they save nothing for GC. jsgo's discriminator bits in
  `kindFlags` are the right Go form.
- **The Go `arena` experiment, or untyped bump arenas for pointerful data.** The experiment is
  frozen, and Go needs per-object type information to scan.
- **Per-node `weak.Pointer` or finalizers** to learn when nodes die. The per-object cost is far
  too high. Chunk-level weak pointers do not work either: int32 links are invisible to the
  GC, so a live chunk reachable only through ids would be collected.
- **Wholesale InternPool-style rewrite of JSType** into `{tag u8, data u32} + extra`. It would
  rewrite a ~150-method interface hierarchy and every pass that touches types. Types are
  15–21% of typed peak and die at `typesToColors`, so the return does not justify it.
- **In-place slot overwrite on replace** (Zig `replaceInst`, Carbon splice). Closure's
  `replaceWith` leaves the old node a valid detached object that callers keep using, so its
  identity cannot be reused.
- **Hash-consing AST subtrees.** Nodes have parent links, positions and mutable properties.

## Sources
- Zig: [Ast.zig](https://github.com/ziglang/zig/blob/master/lib/std/zig/Ast.zig), [Zir.zig](https://github.com/ziglang/zig/blob/master/lib/std/zig/Zir.zig), [Air.zig](https://github.com/ziglang/zig/blob/master/src/Air.zig), [InternPool.zig](https://github.com/ziglang/zig/blob/master/src/InternPool.zig), [Liveness.zig](https://github.com/ziglang/zig/blob/master/src/Air/Liveness.zig), [Legalize.zig](https://github.com/ziglang/zig/blob/master/src/Air/Legalize.zig), [PR #7920](https://github.com/ziglang/zig/pull/7920), [PR #15569](https://github.com/ziglang/zig/pull/15569), [talk summary (secondhand)](https://forums.swift.org/t/what-are-the-benefits-of-minimizing-the-memorylayout-size-stride-of-a-struct/74205/3)
- Carbon: [toolchain/docs/README.md](https://github.com/carbon-language/carbon-lang/blob/trunk/toolchain/docs/README.md), [parse.md](https://github.com/carbon-language/carbon-lang/blob/trunk/toolchain/docs/parse.md), [idioms.md](https://github.com/carbon-language/carbon-lang/blob/trunk/toolchain/docs/idioms.md), [check/README.md](https://github.com/carbon-language/carbon-lang/blob/trunk/toolchain/docs/check/README.md), [token_info.h](https://github.com/carbon-language/carbon-lang/blob/trunk/toolchain/lex/token_info.h), [tree.h](https://github.com/carbon-language/carbon-lang/blob/trunk/toolchain/parse/tree.h), [value_store.h](https://github.com/carbon-language/carbon-lang/blob/trunk/toolchain/base/value_store.h), [inst.h](https://github.com/carbon-language/carbon-lang/blob/trunk/toolchain/sem_ir/inst.h)
- rustc: [context.rs](https://github.com/rust-lang/rust/blob/master/compiler/rustc_middle/src/ty/context.rs), [intern.rs](https://github.com/rust-lang/rust/blob/master/compiler/rustc_data_structures/src/intern.rs), [rustc_arena](https://github.com/rust-lang/rust/blob/master/compiler/rustc_arena/src/lib.rs), [span_encoding.rs](https://github.com/rust-lang/rust/blob/master/compiler/rustc_span/src/span_encoding.rs), [bit_set.rs](https://github.com/rust-lang/rust/blob/master/compiler/rustc_index/src/bit_set.rs), [hir.rs OwnerNodes](https://github.com/rust-lang/rust/blob/master/compiler/rustc_hir/src/hir.rs), [typeck_results.rs](https://github.com/rust-lang/rust/blob/master/compiler/rustc_middle/src/ty/typeck_results.rs), [patch.rs](https://github.com/rust-lang/rust/blob/master/compiler/rustc_mir_transform/src/patch.rs), [simplify.rs](https://github.com/rust-lang/rust/blob/master/compiler/rustc_mir_transform/src/simplify.rs), [PR #37445](https://github.com/rust-lang/rust/pull/37445), [PR #80014](https://github.com/rust-lang/rust/pull/80014)
- Go: [ssa/value.go](https://github.com/golang/go/blob/master/src/cmd/compile/internal/ssa/value.go), [ssa/func.go](https://github.com/golang/go/blob/master/src/cmd/compile/internal/ssa/func.go), [ssa/cache.go](https://github.com/golang/go/blob/master/src/cmd/compile/internal/ssa/cache.go), [ssa/allocators.go](https://github.com/golang/go/blob/master/src/cmd/compile/internal/ssa/allocators.go), [ir/mini.go](https://github.com/golang/go/blob/master/src/cmd/compile/internal/ir/mini.go), [src/xpos.go](https://github.com/golang/go/blob/master/src/cmd/internal/src/xpos.go), [gc/compile.go](https://github.com/golang/go/blob/master/src/cmd/compile/internal/gc/compile.go), [noder/unified.go](https://github.com/golang/go/blob/master/src/cmd/compile/internal/noder/unified.go)
