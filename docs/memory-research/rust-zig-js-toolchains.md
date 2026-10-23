# Cluster: rust-zig-js-toolchains (oxc, swc, Bun/esbuild, Biome/rowan)

Sources read (shallow clones, read only, under `scratchpad/research/src/`):
oxc `0e6498f7` (2026-10-23), swc `d2cc0ad7` (2026-10-22), Bun Zig tag `bun-v1.3.14`
(plus Bun `main` `5c6d4e24`, which is now a Rust port that keeps the same design:
`src/ast/new_store.rs`, `src/ast/*.rs`), biome `main`, rowan `main`. I also used the oxc
performance doc (oxc-project.github.io `src/docs/learn/performance.md`) and esbuild's
`docs/architecture.md`, which Bun's parser is derived from. jsgo facts come from
`internal/rhino/node.go`, `internal/jscomp/{scope,abstract_var,reference,typed_var}.go`,
`internal/parser/*` and `docs/gclab-2026-10-18.md`.

---------------------------------------------------------------------------------------------

## 1. Per project

### 1.1 oxc (Rust)

1. **Bump arena that owns the whole AST** (`crates/oxc_allocator/src/allocator.rs`, `arena/mod.rs`,
   derived from bumpalo). Chunks at least double in size. `Box<'a,T>` is one pointer and never
   runs `Drop`. `Vec<'a,T>` (`vec2/raw_vec.rs`) is `ptr + len:u32 + cap:u32 + &alloc` = 24 B.
   `Allocator::reset()` keeps only the last (largest) chunk and rewinds it; the docs tell users
   to reuse one allocator per file ("If workload is completely uniform, it reaches stable state
   on the 3rd round"). **Nothing is freed node by node.** The minifier rewrites the AST in place
   (`TakeIn`/`ReplaceWith` put a dummy into the arena), so replaced subtrees stay in the arena as
   garbage until it is dropped. The perf doc credits the arena with about 20 %.
2. **AST layout is asserted at compile time** (`crates/oxc_ast/src/generated/assert_layouts.rs`).
   The `Expression` and `Statement` enums are 16 B (tag + Box), because every variant is boxed;
   the perf doc says boxing gave about 10 %. Every struct starts with `span:{u32,u32}` (8 B) and
   `node_id:u32`. Sizes: IdentifierReference 32 B (span, node_id, reference_id, name 16 B),
   BindingIdentifier 32, BinaryExpression 48, ComputedMember 48, CallExpression 64 (includes a
   24 B argument Vec), StaticMemberExpression 64 (the property name is stored inline, not as a
   node), NumericLiteral 40, StringLiteral 48, Function 96, Program 144. **The AST has no
   parent or sibling pointers**, only owning child pointers. A node costs about 32–64 B plus a
   16 B enum slot in its parent. Moving spans from usize to u32 gave up to 5 % on large files.
3. **Strings are not interned globally.** `Ident<'a>` (`crates/oxc_str/src/ident.rs`) is
   `ptr + (len:u32 | hash:u32)` = 16 B. It points into the source text, or into the arena for
   the rare escaped identifiers, and carries a precomputed hash. `CompactStr` is a 24 B
   inline-or-heap string for lifetime-free uses. The lexer `Token` is a `u128` (16 B: start,
   end, kind, flags). Unescaped string values go to a side map keyed by token start
   (`lexer/mod.rs: escaped_strings: FxHashMap<u32, JSStr>`), so tokens carry no payload. The
   perf doc says removing a globally locked string cache made parallel parsing about 30 %
   faster.
4. **The semantic model is struct-of-arrays, indexed by u32 ids** (`oxc_syntax`: NodeId,
   ScopeId, SymbolId and ReferenceId are all `NonMaxU32`, so `Option<Id>` is still 4 B). See
   `crates/oxc_semantic/src/scoping.rs`:
   - `ScopeTable{parent_ids, node_ids, flags}` and
     `SymbolTable{spans, flags, scope_ids, declarations}` are each one allocation with one
     len/cap (`multi_index_vec!`, "modeled on Zig's MultiArrayList").
   - The 2-D parts (symbol names, per-scope `bindings: ArenaHashMap<Ident,SymbolId>`, per-symbol
     `resolved_references: ArenaVec<ReferenceId>`) live in a private arena
     (`ScopingCell`, via `self_cell`) that is dropped in one go.
   - `Reference` is 16 B with no pointers: `{node_id, symbol_id:Option, scope_id, flags}`.
   - Redeclarations are in a side map keyed by SymbolId instead of a `Vec` per symbol (the
     comment: 32 B per symbol becomes 4 B).
   - `AstNodes` (`node/nodes.rs`) is `IndexVec<NodeId, AstNode{kind:16 B, scope_id}>` plus
     side vectors `parent_ids` and `flags`. **Parent links live in a side table**, built only on
     request (`with_build_nodes(true)`).

   `Stats{nodes, scopes, symbols, references}` pre-sizes every table exactly when semantic is
   rebuilt.

   Rough totals: a symbol is about 20 B of SoA + 16 B name + 24 B reference-list header + one
   binding map entry, about 80–90 B. A reference is 16 B + 4 B in its symbol's list. A scope is
   about 12 B of SoA plus its bindings map.
5. **Scope data is not rebuilt per pass** (`crates/oxc_minifier/src/compressor.rs`, `lib.rs`).
   `Scoping` is built once before compression. The fixed-point peephole loop then updates it
   incrementally ("it only prunes references for nodes it drops, and no longer rebuilds
   liveness from scratch each pass"). It is rebuilt exactly once more, before mangling, and
   pre-sized from the first build's `Stats`.

   The mangler (`crates/oxc_mangler/src/lib.rs`, esbuild's slot algorithm) uses symbol-indexed
   `Vec<u32>` slots, frequency vectors, `BitSet`s and `InlineString` names, all in a temporary
   allocator that the caller resets between files. It writes new names into
   `Scoping.symbol_names` (a side table by SymbolId), and codegen reads names by symbol id, so
   **the AST is never rewritten for renaming**.

### 1.2 swc (Rust)

1. **`hstr::Atom` is 8 B** (`crates/hstr/src/lib.rs`, `dynamic.rs`). Strings of 7 bytes or
   fewer are stored inline in the word. Longer ones are a tagged pointer to a
   `ThinArc<{hash:u64}, [u8]>` entry with a precomputed hash. Stores are thread-local
   (`thread_local! GLOBAL_DATA`), with no global mutex ("parsers are embarrassingly parallel").
   Atoms are refcounted, and the global store is never shrunk unless `global_atom_store_gc()` is
   called.
2. **Hygiene uses marks instead of per-pass scope objects**
   (`swc_common/src/syntax_pos/hygiene.rs`, `swc_ecma_transforms_base/src/resolver/mod.rs`).
   `Ident{span:8 B, ctxt:SyntaxContext(u32), sym:Atom(8 B), optional}`. The resolver runs once
   and stamps every identifier with the SyntaxContext of its binding scope. After that,
   `Id = (Atom, SyntaxContext)` identifies a binding ("Yes, (Atom, SyntaxContext) always uniquely
   identifies a variable binding"), passes key their maps by `Id`, and a final `hygiene` pass
   renames clashes. Marks and contexts are u32 indices into global tables
   (`MarkData{parent}`, `SyntaxContextData{outer_mark, prev_ctxt}`).
3. **The minifier's usage analysis is recomputed each iteration**
   (`swc_ecma_minifier/src/compress/mod.rs`: `analyze()` inside `optimize_unit_repeatedly`).
   `ProgramData{vars: FxHashMap<Id, Box<VarUsageInfo>>, scopes: Vec<ScopeData>}`. VarUsageInfo
   packs about 30 booleans into one `u32` bitflags field plus u32 counters (`program_data.rs`).
4. **The AST is plain malloc'd `Box<Expr>` trees.** `Expr` is an inline enum, about 80 B
   according to a commented-out assertion in `expr.rs:174`. `Span` is 8 B (`BytePos(u32)` × 2,
   tested). `swc_allocator` offers an opt-in "scoped" arena. **swc is the least compact of the
   cluster; its useful lessons are the Atom and the resolve-once binding ids.**

### 1.3 Bun's JS parser (Zig, derived from esbuild)

1. **`Expr` is a 32 B value**: `{loc:i32, data:24 B tagged union}`
   (`src/js_parser/ast/Expr.zig:2130`, `comptime assert_eql(@sizeOf(Data), 24) // Do not
   increase the size of Expr`). Leaves live inside the union with no allocation: E.Identifier
   (Ref + 3 bools), booleans, numbers (f64), null, undefined, this, missing. Other nodes are
   pointers into the Store; E.Binary is two inline Exprs plus an op, about 72 B. `Stmt` is at
   most 24 B (compile-time check in `Stmt.zig:295`). There are no parent pointers.
2. **The "Store" is a typed block allocator** (`ast/NewStore.zig`). There is one block type of
   `largest_size * 512 * 2` bytes; blocks form a linked list and are allocated by bumping.
   `reset()` rewinds to the first block and reuses the chain, with no free and no zeroing.
   Instances are `threadlocal`. `Expr.Data.Store.begin()` means "create, or reset" and is called
   per file by the transpiler. The bundler instead redirects allocation (`ASTMemoryAllocator`) to
   per-thread mimalloc heaps and keeps ASTs until linking.
3. **Identifiers carry a `Ref`, not a name**: `packed struct(u64){inner_index:u31, tag:u2,
   source_index:u31}` (`ast/base.zig:97`). `Symbol` is asserted to be 88 B
   (`ast/Symbol.zig:142`). It holds `original_name` (a slice into the source), `link` (for
   merging), `use_count_estimate`, `nested_scope_slot` (precomputed for the minifier) and flags.
   Symbols are stored in a flat array per file. `Symbol.Map` is an array of arrays
   `[source][inner]`: "a two-level array was more efficient in profiles … doesn't involve a hash
   … trivial to merge symbol maps from multiple files".
4. **Scopes are built once, in the parse pass.** `Scope{id, kind, parent:*Scope, children,
   members: StringHashMap(Member{ref, loc}), generated: []Ref, flags}`
   (`ast/Scope.zig`). The visit pass replays them in order (esbuild: `scopesInOrder`; "the scope
   tree is not attached to the AST"). The parser does only two full-AST passes. Afterwards
   everything works on Refs, and the renamer uses the precomputed slots and character
   frequencies.

### 1.4 Biome and rowan (Rust, red-green trees)

1. **Green trees are immutable and position-independent.** `GreenNodeHead{kind, text_len:u32}`
   is 8 B in Biome (`biome_rowan/src/green/node.rs`, static assert). Children are 16 B slots
   (`rel_offset:u32` + a pointer), stored in the same ThinArc allocation as the header. Tokens
   own their text and trivia.
2. **Hash-consing** (`rowan/src/green/node_cache.rs`). Tokens, and nodes with at most 3
   children, are deduplicated through a hash built from the children's hashes and compared by
   child identity. The source says "deduping saves 17% of the memory for green nodes".
3. **Red nodes are created on demand** (`rowan/src/cursor.rs: NodeData{rc:u32, parent, index,
   green, offset}`). They are refcounted and freed as soon as no cursor holds them. Parent
   pointers and absolute offsets exist only for nodes currently being looked at. Edits build new
   green nodes that share all untouched subtrees.
4. **Biome's semantic model** (`biome_js_semantic/src/semantic_model/model.rs`) stores
   `Vec<Binding>` and `Vec<Scope>` indexed by u32 ids. Its lookup tables are keyed by source
   offsets (`FxHashMap<TextSize, BindingId>`), which only works because the trees are immutable.
5. Boa's interner uses `Sym(NonZeroUsize)` with keywords pre-interned as statics. This teaches
   nothing beyond what jsgo's 32-bit StringID already does.

### 1.5 Summary table

| | AST node | symbol / scope | index vs pointer | freeing | strings | scope reuse |
|---|---|---|---|---|---|---|
| oxc | 32–64 B + 16 B parent slot; no parent or sibling links | symbol ≈ 80–90 B (SoA), reference 16 B, scope ≈ 12 B + map | AST: arena pointers; semantic: u32 ids | whole arena at once | 16 B slice + hash, no global interning | built once, updated incrementally |
| swc | `Expr` ≈ 80 B inline enum, boxed children | `VarUsageInfo` per Id in a HashMap | pointers; bindings via u32 ctxt | malloc and Drop | 8 B Atom (inline ≤ 7 B) | resolver once; usage data rebuilt per iteration |
| Bun / esbuild | 32 B Expr value (leaves inline) + pointee | Symbol 88 B in per-file arrays | Ref = u64 (source, inner) | Store reset per file | slices into the source | scopes built in parse pass, replayed |
| rowan / Biome | 8 B header + 16 B per child, deduplicated | Vec by u32 id, maps by offset | green: Arc; red: refcounted cursors | refcounts | tokens deduplicated | immutable trees, rebuilt per edit |

---------------------------------------------------------------------------------------------

## 2. Recommendations for jsgo, ranked by impact on peak live heap and pointer count

Baseline numbers. `rhino.Node` is 64 B and contains **5 pointer words** (next, first, parent,
previous, typeOrColor). typed-full's AST is 160 MB, about 2.5 M nodes and about 12.5 M pointer
words. gclab (`docs/gclab-2026-10-18.md`) measured the AST at 37–89 % of the live heap and
45–85 % of GC scan work. It also measured that **typed-app allocates about 90 MiB of nodes
against about 35 MiB live**, and bundle allocates about 3× its live AST. That ratio drives
everything in R1.

Two Go facts matter throughout:
- **Noscan memory removes GC mark CPU but not RSS.** The heap goal is
  live × (1 + GOGC/100) and includes noscan bytes, so peak RSS only drops when bytes drop.
- **Go has no arena reset**, and an interior pointer keeps its whole allocation alive.

### R1. Pointer-free node records in noscan chunks, keeping `*rhino.Node` as the handle
This is the oxc/Bun arena, adapted to Go.

- **Applies to:** `internal/rhino/node.go` (Node, link methods, clone slab), `jutil.Slab[rhino.Node]`
  users (16 sites), and about 34 direct `rhino.Node{}` allocations.
- **Mechanism.** A record with `self, next, first, parent, previous int32` (20 B), `payload`
  (8), `propFlags` (8), `linenoCharno` (4), `length16 + token + kindFlags` (4), and
  `typeOrColor` as a `uint32` handle (4) is **48 B with no Go pointers**, so `[]Node` chunks are
  allocated noscan.
  - `GetNext()` becomes `chunks[i>>12][i&4095]`. It still returns a `*Node`: an interior
    pointer into a noscan chunk is legal and keeps that chunk alive.
  - The index space is global int32 (the 35 M-node googleapis build fits). Chunks of 1–4 K
    records are reserved per goroutine from an atomic counter, like Bun's `threadlocal` Store.
  - The handle points into id tables: JSType serial id → registry table, Color id → color
    table, and an append-only `[]unsafe.Pointer` page table for JSDocInfo and the extension
    box. With this layout, "drop all JSTypes after colors" is simply clearing a table.
- **Gain.** −16 B/node (−25 %): about −40 MB on typed-full and roughly −100 MB or more on
  bundle-mega's AST. The AST's 5 pointer words per node, about 12.5 M on typed-full, drop out
  of marking entirely; gclab estimates −50–80 % of mark work.
- **Passes stay 1:1:** pointer identity, `map[*Node]` keys and every accessor are unchanged.
  Byte-identical output is safe, since nothing depends on addresses.
- **The catch, which oxc and Bun accept and jsgo may not.** With int32 links the GC cannot see
  references between nodes, so the chunk table must keep every chunk alive: **nodes are never
  freed.** oxc and Bun accept this because their arena dies per file or per process.
  - Today's 128-node slabs already pin dead nodes, but a slab can die once all its nodes are
    dead.
  - With gclab's 2.6–3× allocated-to-live ratio, late-compile AST bytes would reach
    48 × 2.6–3 / 64 ≈ **2×** today's live AST.
  - That is probably fine for typed-full, whose peak is at checkTypes before most garbage is
    produced. It is **dangerous for untyped bundle-mega** if its peak falls in the
    optimization loop.
- **Before adopting it,** count allocated nodes against nodes reachable from the root (plus
  JSDoc roots) at each phase boundary, for typed-full and bundle-mega.
- **Option for chunk-level reclamation** (my Go-specific adaptation, not taken from these
  projects):
  - Give each chunk a small scanned `refs []*chunk` remembered set. Every link store (all of
    them are in node.go) records the target's chunk when it differs from the source's.
  - The global chunk table becomes `uintptr` plus a `weak.Pointer` (Go 1.24+) used only to
    detect freed slots.
  - A chunk then dies when no `*Node` handle and no referencing chunk reaches it. That gives
    today's slab-level reclamation, the same retention as 128-record chunks, with
    near-zero scan cost.
  - Sets are rebuilt at pass boundaries.
  - Risk: one missed barrier is a use-after-free, so it needs a debug verifier that walks
    links and asserts the target chunk is in the source's refs set.
- **Effort:** large inside rhino; zero in passes. **Mode:** default, as long as the garbage
  question is solved.

### R2. Dense side tables by node id instead of `map[*Node]` and `map[uintptr]`
oxc's `IndexVec<NodeId, …>` and Bun's arrays.

- **Applies to:** `nodeExtraShards` (64 RWMutex shards of `map[uintptr]nodeExtraEntry`,
  node.go:3108), the sourceFile property chunks, and 57 `map[*rhino.Node]` maps in jscomp.
- **Mechanism.** R1's `self` gives each node a dense id for free. A `NodeTable[V]` is a paged
  `[]*[4096]V` with pages allocated on first touch; it is noscan when V has no pointers, needs
  no hashing and no locks.
  - For maps that are only looked up (Java's `HashMap<Node,…>` can't have a meaningful order,
    since Node uses identity hashes), swap the type and keep the call sites.
  - For maps that are iterated, keep the insertion-order slot list with a node-id index (the
    in-flight int32 LinkedHashMap, keyed by dense id).
- **Gain.** A pointer-keyed map entry costs about 24–48 B and is scanned; a dense entry costs
  `sizeof(V)` and is unscanned. Peak impact is moderate (only long-lived tables count); churn
  and GC impact are larger.
- **1:1:** yes. **Risk:** none for lookup-only maps. **Mode:** default.

### R3. Pointer-light Reference, Var and Scope records (oxc `Reference` is 16 B; Bun's symbol arrays)
- **Applies to:** `reference.go` (`Reference{nameNode, basicBlock, scope, inputId}` = 32 B with
  4 pointers, each allocated separately and listed in `ArrayList[*Reference]`),
  `abstract_var.go` (Var is about 40 B with 3 pointers), `scope.go` (Scope is about 112 B with
  7 pointer words), `typed_var.go` and `typed_scope.go`.
- **Mechanism.**
  - Reference becomes `{nameNode *Node; scope, bb, input int32}` = 24 B with 1 pointer, stored
    by value in chunked slabs or per-variable value slices. `*Reference` stays an interior
    pointer, so the getters are unchanged. With R1 the `nameNode` becomes int32 and the whole
    record is noscan.
  - Var can drop `input *CompilerInput`, which is derivable from the node's source-file id,
    leaving 2 pointers.
  - The scopes of one creator go into one `[]Scope` with int32 parent indices.
  - Pre-size from the previous run's counts, as oxc does with `Stats`.
- **Gain.** gclab: scopes, vars and references are 8.5–11 % of churn, plus most of the 17–19 %
  pass-local share. At the typed peak, TypedScope and TypedVar are about 6 % of the heap
  (≈ 33 MB), of which roughly 30–40 % is removable.
- **1:1:** yes. **Mode:** default. This overlaps the in-flight scope/var work; the layout above
  is offered to it.

### R4. Parser: 16 B tokens and u32 spans with lazy line and column (oxc)
- **Applies to:** `internal/parser` (`SourcePosition{*SourceFile; Offset, Line, Column int}` =
  32 B; `SourceRange{Start, End *SourcePosition}`; `ParseTreeBase{self ParseTree (16 B
  interface), Type, Location}`; Token structs).
- **Mechanism.**
  - Ranges become `{start, end uint32}` UTF-16 offsets, with line and column computed by binary
    search in a per-file `[]uint32` table of line starts (Java's UTF-16 columns stay exact).
  - Tokens become 16 B values: offsets, kind and flags.
  - Decoded values go into side maps keyed by start offset, and only when they differ from the
    raw text.
  - Parse-tree references become typed 8 B pointers or indices instead of 16 B interfaces.
- **Gain.** About 80 B and 3 pointers of position data per parse tree become 8 B with no
  pointers. This hits the single-large-file parse peak (bundle).
- **1:1:** parser-internal. **Risk:** low with a position oracle. Already in flight; this is
  the concrete target layout.

### R5. JSDoc type expressions stored as green trees (rowan)
- **Applies to:** `rhino/js_type_expression.go` (`JSTypeExpression{root *Node, sourceName}`),
  `parsing/js_doc_info_parser.go`.
- **Mechanism.**
  - Encode each type-expression tree once into a pointer-free blob of (kind, StringID, relative
    offset, child count).
  - Hash-cons blobs with at most 3 children: `{string}`, `{!Array<string>}`,
    `{function(number): boolean}` repeat thousands of times in typed Closure code.
  - `JSTypeExpression` keeps the blob id plus the root's absolute line, column and source file.
  - `GetRoot()` (54 call sites in jscomp) materializes a "red" `rhino.Node` tree on demand with
    exact absolute positions, cached per JSTypeExpression so identity stays stable while used.
  - Mutations become copy-on-write: re-encode and repoint. The known mutation sites are
    `es6_rewrite_modules.go:767`, `process_common_js_modules.go:2164–2229`,
    `es6_rename_type_references.go:128`, `check_js_doc.go:856`.
- **Gain.** JSDoc is about 10 % of typed peak (≈ 55 MB), much of it type-expression nodes at
  64 B each; rowan measured 17 % from dedup alone. Gains are only real if the materialized
  trees are transient: dropped after TypedScopeCreator has resolved them, or rebuilt each time.
- **Risk:** medium (identity, in-place mutation). **Effort:** high. **Mode:** default if
  positions stay exact. This is for the JSDoc compaction team to evaluate.

### R6. Long-lived symbol tables keyed by dense ids (Bun's `[source][inner]` array of arrays)
- **Applies to:** GlobalNamespace `Name`/`Ref`, the type registry's name maps, and
  ReferenceCollector's `ObjectLinkedHashMap[*Var, …]`.
- **Mechanism.** `map[string]*Name` becomes `map[StringID]int32` plus a `[]nameRec` SoA.
  References become 16 B records in per-name slices. Vars get a dense ordinal per scope-creator
  run, so collector maps become slices indexed by ordinal, with a separate insertion-order list
  to keep LinkedHashMap order.
- **Gain:** large on churn, moderate on peak wherever GlobalNamespace is retained.
- **1:1:** yes behind getters. **Mode:** default.

### R7. Allocator hygiene (oxc `reset` keeps only the largest chunk; oxc `Stats` pre-sizing)
- In Go, a retained arena chunk is live heap. Keep cross-pass pools only for pointer-free
  buffers, keep at most one chunk, and drop everything at phase boundaries (before
  checkTypes, and after typesToColors).
- Pre-size per-pass tables from the previous run's counts, as oxc's `with_stats` does.
- Clip long-lived slices: append's 2×/1.25× growth leaves up to 50 % slack. `slices.Clip`, or
  an exact copy, at phase ends (JSDoc lists, registry slices).
- **Gain:** a few percent of peak. Simple; default mode.

### R8. Reusing scopes across passes (oxc incremental Scoping, esbuild bind-once, swc resolve-once)
- Memoize `SyntacticScopeCreator` results per scope root, invalidated by Closure's per-function
  change stamps. Verify in debug builds by rebuilding and comparing.
- **Gain:** allocation volume (scopes are 8.5–11 % of churn), not peak, because syntactic scopes
  are transient.
- **Risk:** any pass that edits the AST without reporting the change gets stale Vars and wrong
  output.
- **Mode:** fast_mode first; default only after the verifier is clean on the corpora.

### R9. Strings: keep the current design
jsgo's 4 B StringID in noscan pages is denser than oxc's 16 B `Ident`, swc's 8 B Atom and Bun's
16 B slices. Two minor ideas:
- Encode strings of 3 bytes or fewer directly in the ID, as swc does for 7 bytes or fewer.
  This saves pool lookups for mangled names.
- Add a per-parser local intern cache in front of the RWMutex shards (swc's thread-local
  stores; oxc saw 30 % on parallel parsing). That is a wall-time gain only.

Ranking by impact on peak heap and pointer count:
**R1 > R4 (parse peak) > R5 > R3 > R2 > R6 > R7 > R8 (allocation volume only) > R9.**

---------------------------------------------------------------------------------------------

## 3. Fast-mode ideas and rejected ideas

**Fast mode only:**
- **R8 scope memoization without the verifier.**
- **Hash-consing type expressions with positions dropped.** This is R5 without exact
  positions: it changes diagnostic positions, but stays deterministic.

**Rejected:**
- **No parent or sibling links, with parents in an on-demand side table (oxc AstNodes, rowan
  red nodes).** rhino's API needs O(1) `getParent`, `getNext` and `getPrevious` for every node,
  at any time.
- **Inline leaf nodes (Bun's identifier, number and boolean inside the parent's 24 B
  union).** Every leaf must be an addressable `*Node` with a parent.
- **swc's (sym, ctxt) binding ids, or esbuild's Refs, on NAME nodes, with passes keyed by
  them.** Every pass using Scope/Var would have to be rewritten. (Closure's
  MakeDeclaredNamesUnique is the moral equivalent, and passes don't exploit it either.)
- **Source slices instead of interning (oxc, Bun).** They would be larger than 4 B IDs and
  would keep source text alive. Likewise refcounted atoms (swc), which have no place under a
  GC.
- **A pure arena with reset per pass for pass-local data (Bun's Store).** In Go, the GC already
  reclaims unreachable pass data. Retained arenas raise peak and were the 22-round problem.
  Arenas pay off only for pointer-free scratch.
- **Interning whole JSDocInfo objects** (position-free `@const`, `@private`, `@override`).
  This needs JSDocInfo to be immutable and never compared by identity; the gain is small and
  the priority low.
- **Offset-keyed semantic tables (Biome).** They need immutable trees; jsgo nodes are mutated
  and synthesized nodes share positions.
