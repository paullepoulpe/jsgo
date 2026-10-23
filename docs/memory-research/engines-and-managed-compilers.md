# Cluster report: JS engine front ends and compilers on managed runtimes

Scope: V8 (parser, preparser, Zone, scopes, Turbofan/Turboshaft/Maglev), SpiderMonkey (LifoAlloc,
ParseNode, ParserAtoms, Stencil), Hermes (ESTree arena, lazy compilation), QuickJS (atoms), Closure
Compiler on the JVM (per-object sizes as the target), Roslyn, and briefly Scala 3, Kotlin K2
and typescript-go (a JS-family compiler that is itself written in Go).

Method: read the projects' current source (raw files from v8/v8 main, mozilla-firefox/firefox main,
facebook/hermes main, bellard/quickjs, google/closure-compiler tag v20261015,
dotnet/roslyn main, scala/scala3 main, microsoft/typescript-go main) and V8's blog sources
(v8/v8.dev). Nothing was built or run. jsgo was read, not modified. All sizes for jsgo are from the
struct definitions; JVM sizes assume 64-bit HotSpot with compressed oops and compressed class pointers
(12 B header, 4 B references, 8 B alignment, subclass fields can fill superclass gaps).

---

## 0. Findings that matter most

1. **jsgo's AST is not larger than Java's. The difference is how the GC treats it.** A Java
   `Node` is 56 B, and every boolean or int property adds a 24 B `IntPropListItem`. jsgo's Node is
   64 B with 40 boolean properties stored inline. For a typical node with 1-2 properties, Java uses
   80-104 B and jsgo 64 B. Java pulls ahead on large inputs because G1 is generational: the
   long-lived AST is promoted and is not traced again on every young collection. Go's
   non-generational GC traces all 40 pointer bytes of every live node on every cycle. **The goal for
   the AST is pointer-free (noscan) memory, not fewer bytes.** For types and JSDoc the opposite holds:
   jsgo's type objects are 2-3x their Java size (table in §1.5), so bytes matter there.
2. **Every engine studied keeps positions as 32-bit source offsets and computes line and column
   only when asked**, from a per-file table of line starts plus a "last line looked up" cursor
   (V8 `Scanner::Location`, SpiderMonkey `TokenPos{uint32 begin,end}` + `SourceCoords`, Hermes
   `SMLoc`, Scala 3 `Span`, typescript-go `TextRange{int32,int32}`). Tokens
   are value objects in a ring of 4 (V8 `token_storage_[4]`, SpiderMonkey `tokens[ntokens=4]`) or
   structs (Roslyn's red `SyntaxToken`). jsgo's parser keeps **one TokenBase (32 B), one SourceRange
   (16 B) and 1-2 SourcePositions (32 B each) per token until parsing finishes**: 80-140 B of
   pointer-bearing memory per token. That is the main parse-time peak lever.
3. **V8's pointer compression is the model for "software compressed oops" in Go:** 32-bit offsets
   from a base (V8 reports up to 43% less heap). Turboshaft applies the same idea inside the compiler
   with `OpIndex`, a uint32 offset into one operation buffer, and keeps per-op analysis data in side
   tables. In Go, an index scheme only pays off if the indexed memory holds no Go pointers. Then it
   costs the GC almost nothing, but anything it refers to by ID has to be kept alive by an explicit
   table. V8 removed its compressed *zone* pointers (`graph-zone-traits.h` now uses plain
   pointers): it compresses long-lived heap data, not short-lived compiler temporaries. Do the same
   in jsgo: index the long-lived AST, not per-pass scratch data.
4. **Arenas in C++ engines exist to free memory in bulk.** In Go, an arena of pointer-bearing
   structs only saves malloc time. It is still scanned, one interior pointer keeps the whole chunk
   alive, and keeping chunks across passes turns garbage into live, scanned memory at the peak. V8
   Zones, SpiderMonkey LifoAlloc and Hermes's allocator are always discarded at the end of a
   compilation unit. Roslyn's pools have a hard size cap ("the pool is not storage"; builders with
   capacity ≥128 are dropped).
5. **Lazy techniques (lazy function parsing, red nodes, stencils) mostly do not fit Closure**,
   because a whole-program optimizer reads and changes every function. The parts that do fit:
   offsets instead of materialized positions, transient tokens, lazily parsed and later discarded
   JSDoc type expressions (only a modest gain, see R9), and pre-serialized runtime libraries, which
   already exist as `runtime_libs.typedast`.
6. **Outside this cluster but important:** jsgo pins Go 1.24 (`go.mod`). Go 1.26 turns on the
   "Green Tea" GC by default, which the release notes say cuts GC overhead by 10-40% for programs
   that allocate many small objects (another ~10% with AVX-512). That is a no-code experiment and
   should be measured before or alongside any layout work.

---

## 1. Per project

### 1.1 V8 (front end and optimizing compilers)

- **Zone allocation (`src/zone/zone.h`, `zone-segment.h`).** A bump allocator (`position_`,
  `limit_`, `segment_head_`) over segments from an `AccountingAllocator`. It has no per-object
  free: the whole Zone is discarded after a parse or compile job. A `ZoneSnapshot` allows rollback
  (`Reset`/`Snapshot`). Allocations are tagged per type for `--trace-zone-stats`, which reports
  per-zone and per-phase high-water marks. The parser and every Turbofan, Turboshaft and Maglev
  phase use zones with a defined lifetime.
- **Compact AST nodes (`src/ast/ast.h`).** `AstNode` = `int position_` (a source *offset*) +
  `uint32_t bit_field_` (6-bit node type plus per-subclass flags through `NextBitField`). That is
  8 B with no vtable and **no parent or sibling pointers**. Children are typed fields, and lists
  are zone vectors. `VariableProxy` stores `union { const AstRawString* raw_name_; Variable* var_; }`
  plus an intrusive `next_unresolved_`: the name before resolution and the variable after it share
  one slot. Typical nodes are 16-40 B. The AST is used once and then discarded after bytecode
  generation, which is why it can be this lean.
- **String interning (`src/ast/ast-value-factory.h`).** `AstRawString` = literal bytes in the zone
  + `raw_hash_field_` + `is_one_byte_` + a union of `next_` and the internalized handle. Strings
  are deduplicated in `string_table_`, so equality is pointer comparison. **Strings and AST have
  separate lifetimes**: `ast_raw_string_zone_` outlives `single_parse_zone_`.
  `AstStringConstants` is a shared, prebuilt table of well-known strings.
- **Scopes (`src/ast/scopes.h`, `variables.h`).** `Variable` = scope, name, `local_if_not_shadowed_`,
  intrusive `next_` (ThreadedList), `int index_`, `int initializer_position_`, and two `uint16_t`
  bitfields (mode:4, kind:3, location:3, flags), about 48 B. `Scope` = outer/inner/sibling links
  (first-child/next-sibling tree), `VariableMap` (a zone hash map keyed by the interned string
  pointer), intrusive `locals_`/`unresolved_list_`/`decls_` lists, start/end positions as ints, and a
  `uint32 flags_`.
- **Preparser and lazy parsing (`src/parsing/preparser.h`, `preparse-data.h`; v8.dev/blog/preparser).**
  The preparser builds no AST. `PreParserExpression` is a 32-bit coded value (`TypeField`,
  `ExpressionTypeField`, ...). It still tracks variables, so early errors and scope resolution stay
  correct. "Skipping inner functions": scope-allocation results are stored "as a dense array of flags
  per variable" (is_used, maybe_assigned, forced context allocation, inner_scope_calls_eval). The data
  is varint-encoded, kept only for functions that contain inner functions, and freed after the lazy
  compile. Each function is preparsed at most once.
- **Positions.** In the AST a position is an int offset. Line and column come from the script's
  line-ends array only when needed. Generated code uses a delta- and VLQ-encoded
  `SourcePositionTable`, and `SourcePosition` packs everything into one `uint64_t`
  (`ScriptOffsetField` 30 bits, `InliningIdField` 16 bits, ...).
  **"Lazy source positions"** (v8.dev/blog/v8-lite): bytecode is compiled *without* position tables,
  which are regenerated by reparsing when a stack trace needs them. This required making bytecode
  generation deterministic, with a stress mode checking that eager and lazy compilation produce
  identical output. Lazy feedback allocation and bytecode flushing are the same pattern: allocate
  optional data on first use and drop it when cold. Together they cut heap by 18% on average.
- **Tokens (`src/parsing/scanner.h`).** `TokenDesc token_storage_[4]` holds
  current/next/next_next/next_next_next. Each `TokenDesc` has `Location{beg_pos,end_pos}`, the token
  value and literal buffers. Tokens never become heap objects.
- **Turbofan node (`src/compiler/node.h`).** Variable-size allocation: the `Use` records come
  *before* the `Node` in memory and the input pointers come *after* it ("Saving space for big graphs
  is important"). Pointer arithmetic recovers Node⇄Use, the id/inline-count/inline-capacity fields are
  packed into one `bit_field_`, and inputs overflow to an out-of-line block.
  `ZoneNodePtr = GraphZoneTraits::Ptr<Node>` was once compressible (`kCompressGraphZone`, 2020), but
  today `graph-zone-traits.h` maps it to a plain pointer (`FullZonePtr`).
- **Turboshaft (`src/compiler/turboshaft/index.h`, `graph.h`, `operations.h`, `sidetable.h`,
  `snapshot-table.h`; v8.dev/blog/leaving-the-sea-of-nodes).** All operations of a graph live in one
  append-only `OperationBuffer` of 8 B slots. `OpIndex` is "an offset from the beginning of the
  operations buffer. Compared to `Operation*`, it is more memory efficient (32bit) and stable when the
  operations buffer is re-allocated." The operation header is `opcode` (uint8),
  `saturated_use_count` (uint8) and `input_count` (uint16), followed by inline `OpIndex` inputs.
  Analysis data is never stored in operations: it goes into `GrowingOpIndexSidetable<T>` (dense
  vectors indexed by OpIndex). Phases copy the input graph into a new output graph, which also
  restores locality. `SnapshotTable` gives key→value state with O(1) Get/Set and snapshots costing
  O(changes since the common ancestor), used by load elimination instead of copying state per block.
  Reported results: compile time halved versus Sea of Nodes, Sea of Nodes had about 3x more L1d
  misses, and the new load elimination is up to 190x faster with "way less memory".
- **Maglev (`src/maglev/maglev-ir.h`).** `NodeBase::bitfield_` is one `uint64_t`
  `[opcode:16][input_count:16][properties:16][extras:16]`. Inputs are allocated in the zone directly
  before the node (`input_base() - index`).

### 1.2 SpiderMonkey

- **LifoAlloc (`js/src/ds/LifoAlloc.h`).** A bump allocator over a chain of `BumpChunk`s. Chunk size
  doubles up to 1 MB and then grows geometrically. Oversize allocations get their own chunks
  (`oversize_` list) so normal chunks are not left half-empty. `LifoAllocScope`/`Mark`/`release`
  free a whole scope at once, and released chunks are reused by the next similar task.
- **ParseNode (`js/src/frontend/ParseNode.h`, `Token.h`).** `pn_type` (uint16 kind), 1-bit flags,
  `TokenPos pn_pos` = `{uint32 begin, end}` (code-unit offsets only), and `ParseNode* pn_next` (the
  link inside the parent's list). There are no parent or previous pointers. Subclasses differ in
  size: `NameNode` adds a 32-bit `TaggedParserAtomIndex atom_`; `BinaryNode` adds `left_`/`right_`;
  `ListNode` adds `head_`, `tail_` (a pointer to the last `pn_next`), `count_` and `xflags`.
  `SyntaxParseHandler` (syntax-only parsing for lazy functions) uses an **enum as its "node"**
  (`NodeGeneric`, `NodeName`, ...), so no tree is allocated at all.
- **Positions (`TokenStream.h`).** `SourceCoords::lineStartOffsets_` is "one element per line … plus
  one sentinel". Converting an offset to a line is "O(log n) in the worst case … but in practice
  they're heavily clustered and we do better than that by using the previous lookup's result
  (lastIndex_) as a starting point". Tokens live in a ring, `Token tokens[ntokens = 4]`.
- **ParserAtoms (`ParserAtom.h`).** A `ParserAtom` is a 12-byte header (hash, length, flags such as
  `HasTwoByteCharsFlag` and `UsedByStencilFlag`) followed by inline characters, allocated in a
  LifoAlloc and indexed through `entries_`. References are a **32-bit tagged index**:
  `0x1YYY_YYYY` is a 28-bit atom index, `0x2000_YYYY` a well-known atom, `0x2001_YYYY` a static
  length-1 string, `0x2002_YYYY` a static length-2 string matching `[A-Za-z0-9$_]{2}`, and
  `0x2003_YYYY` the numbers "100"-"255". Short names and well-known names never touch the table.
- **Bindings (`vm/Scope.h`).** `AbstractBindingName<TaggedParserAtomIndex>` is a single
  `uint32_t bits_`: the atom index plus two flag bits (`ClosedOverFlag`, `TopLevelFunctionFlag`).
  Scope data is a header (slot info such as `nextFrameSlot`) with a trailing array of these 4-byte
  names grouped by kind. Variables are not objects after parsing.
- **Stencil (`Stencil.h` SMDOC, `CompilationStencil.h`).** The output of parsing and bytecode
  emission is "designed for fast serialization to and from disk by preferring indices over pointers
  and vectors instead of graphs to allow bulk operations … a series of vectors/spans segregated by
  data type (a.k.a. Data Oriented Design)". `ScriptStencil` refers to its GC things as a
  `[gcThingsOffset, +gcThingsLength)` range into a shared `gcThingData` vector, using 32-bit tagged
  indices (`TaggedScriptThingIndex`). An immutable `CompilationStencil` can *borrow* its spans from
  an XDR buffer or an mmapped file. Lazily parsed functions are later "delazified" into a new Stencil
  and merged in (`CompilationStencilMerger`). Moving the frontend off the GC required replacing GC
  atoms with ParserAtoms first. Mozilla's newsletters report lower memory use, mainly for
  browser-internal JS, without a single headline number.

### 1.3 Hermes and QuickJS

- **Hermes ESTree (`include/hermes/AST/ESTree.h`, `Context.h`).** `Node : llvh::ilist_node<Node>`
  (intrusive prev/next, like Closure's sibling list) + `NodeKind kind_` + `unsigned parens_` +
  `SMRange sourceRange_` + `SMLoc debugLoc_`. `SMLoc` is a raw pointer into the source buffer and
  the line is resolved on demand by the SourceErrorManager. There is no parent pointer. `NodeList` is
  a `simple_ilist`. Nodes come from `Context::allocator_`, a BumpPtrAllocator that "may be rolled
  back to parse subtrees during pre-parsing (for lazy parsing)", and `delete` is disabled.
  Identifiers live in a **separate** `identifierAllocator_` "because we don't want to revert the
  strings when we revert subtrees", deduplicated by `StringTable` (`UniqueString*`, compared by
  pointer).
- **Hermes lazy compilation (`Parser/PreParser.h`).** For each preparsed function it stores a map
  from body start `SMLoc` to `{end SMLoc, strictMode, directives}`. The real parse skips those bodies
  and compiles them on first call.
- **Hermes IR (`IR/IR.h`).** LLVM style: `Value` has a users list, `Instruction` has
  `SmallVector<Use,2> Operands`. Nothing unusual about its layout.
- **QuickJS (`quickjs.c`).** It has **no AST**: the parser emits bytecode directly. `JSString` is a
  12-byte header (`len:31`, `is_wide_char:1`, `hash:30`, `atom_type:2`, `hash_next` as a uint32
  *index*) followed by 8-bit or 16-bit characters. `JSAtom` is a `uint32_t`, and
  `JS_ATOM_TAG_INT` (bit 31) marks integer atoms that have no table entry. The atom table is
  `atom_array` plus a hash chain of indices and a free list.

### 1.4 Closure Compiler on the JVM (the target to beat)

Instance fields read from `rhino/Node.java`, `rhino/jstype/*.java`, `jscomp/AbstractVar.java`,
`AbstractScope.java`, `parsing/parser/util/*.java` at v20261015:

| Object | JVM size (compressed oops) | jsgo now | Notes |
|---|---|---|---|
| `Node` (base) | **56 B**: 12 header + 8 refs (token, next, previous, first, parent, jstypeOrColor, originalName, propListHead) + 2 ints (linenoCharno, length) | **64 B** | jsgo: 4 node pointers + typeOrColor unsafe.Pointer = 40 B of pointer words (ptrdata) |
| `StringNode` | 56 B (the `str` ref fills the padding) | 64 B (payload) | |
| `NumberNode` | 64 B (double) | 64 B | |
| each boolean/int prop | **+24 B** `IntPropListItem` (12 + next + byte + int); the chain is copied on change; the SOURCE_FILE item is a shared tail | 0 B (40 inline bits) | Java's real per-node cost is often 80-104 B |
| `JSDocInfo` | 32 B (2 longs + Object) | 32 B (2 int64 + `any`) | equal |
| `JSType` fields | 32 B (resolveResult, templateTypeMap, Boolean, registry, 3 booleans) | `JSTypeBase` 64 B (`self` interface 16 + `resolveResult` interface 16 + 2 ptrs + serials + bools) | 2.0x |
| `PrototypeObjectType` | 64 B + `PropertyMap` 40 B + `TreeMap` 40 B + **72 B per property** (TreeMap.Entry 40 + Property 32) | ? | |
| `InstanceObjectType` | 64 B | 152 B | 2.4x |
| `FunctionType` | ≈112 B (+ArrowType, maps) | 224 B | 2.0x |
| `UnionType` | 40 B (+ImmutableList) | 128 B | 3.2x |
| `NamedType` | 88 B | 208 B | 2.4x |
| `Var` | 40 B | 40 B (`AbstractVarBase`) | equal |
| `Scope` | 32 B + LinkedHashMap (56 B + table + 40 B/entry) | ≈112 B (`AbstractScopeBase` + parent, lookupParent, arena, 2 slices, bloom, count) | |
| parser `SourcePosition` | 32 B | 32 B (`*SourceFile` + 3 Go `int`) | |
| parser `SourceRange` / `Token` / `IdentifierToken` | 24 / 24 / 24 B | 16 / 32 / 56 B | |

With JDK 25's optional compact headers (JEP 519, 8-byte headers) the Java numbers drop by 4-8 B per
object. Java also gets scalar replacement and TLAB bump allocation for short-lived temporaries, and
young collections that trace only roots, survivors and dirty cards. **Bytes per AST node are already
comparable. The gap is that Go traces the long-lived AST on every cycle and that types are about 2x
larger.**

Closure already has two "stencil-like" features that jsgo ports. (a) Runtime libraries are a
serialized TypedAST (`runtime_libs.typedast`, 3.4 MB) and are deserialized lazily per file
(`TypedAstDeserializer.typedAstFilesystem`, a `Supplier<Node>` per file). (b) Staged builds decode
`LazyAst` byte blobs per file. jsgo decodes those bytes directly (`ScriptNodeDeserializer` over
`scriptBytes`), and `Compiler.GetTypedAstDeserializer` removes the entry after first use, so the bytes
are released. `ExternsSummary` (extern property names) is a compact summary used by later stages.

### 1.6 Roslyn (C#/VB)

- **Green nodes (`Syntax/GreenNode.cs`, `GreenNode.NodeFlagsAndSlotCount.cs`).** `ushort _kind` +
  a `ushort` packing 12 flag bits and a 4-bit slot count + `int _fullWidth` = 8 B of payload on top of
  the CLR object header. Nodes store **widths, not absolute positions**, and have no parent, which
  is why they can be shared. Rare data (diagnostics, annotations) lives in
  `ConditionalWeakTable<GreenNode, …>` **side tables**, signalled by a flag bit in the node.
- **Hash-consing (`InternalSyntax/SyntaxNodeCache.cs`).** Green nonterminals with up to 3 children
  are looked up by (kind, flags, child identities), so repeated small subtrees are stored once. Nodes
  with diagnostics or annotations are excluded. Keyword and punctuation tokens without trivia are
  static singletons (`s_tokensWithNoTrivia`, `s_tokensWithSingleTrailingSpace`, ...).
- **Red nodes (`Syntax/SyntaxNode.cs`).** A red node is `_parent`, `_syntaxTree`, green and
  absolute position. It is **created lazily** when someone asks for a child, published with
  `Interlocked.CompareExchange` into a field of the parent. Very long lists keep their red children
  through `WeakReference` (`GetWeakRedElement`) so they can be collected. Red tokens are **structs**
  (`SyntaxToken`, `SyntaxNodeOrToken` = parent, green token, position, index) and never allocate.
- **Pooling (`PooledObjects/ObjectPool\`1.cs`, `ArrayBuilder.cs`).** The pool has a fixed size, and
  its comments say: "it is not the goal to keep all returned objects. Pool is not meant for storage",
  and "Keeping checked out objects for long durations … reduces usefulness of pooling". `ArrayBuilder`
  is returned to the pool only if its capacity is below `PooledArrayLengthLimitExclusive = 128`. In
  their measurements most builders hold 0-2 items, and keeping a 50K-element builder alive "is not
  worth it".

### 1.7 Scala 3, Kotlin K2, typescript-go

- **Scala 3 (`util/Spans.scala`, `core/Names.scala`, `tasty/TreeUnpickler.scala`, `ast/Trees.scala`).**
  `Span` is a value class over one `Long`: start (26 bits), end (26 bits), point delta (12 bits), so
  it is never allocated. Names are interned in a global table. TASTy (serialized typed trees) is read
  **lazily**: `ValDef`/`DefDef` hold `preRhs: LazyTree = Tree | Lazy[Tree]`, which is materialized by
  `readLater` on first `.rhs`.
- **Kotlin K2.** The command-line frontend builds FIR from IntelliJ's LighterAST, a flat list of
  markers with flyweight nodes, instead of materializing PSI objects. JetBrains reports K2 analysis up
  to 3.76x faster, but I found no number for the LightTree change alone.
- **typescript-go (`internal/ast/ast.go`, `ast_generated.go`, `core/arena.go`).** `Node` =
  `Kind` (int16) + `Flags` + `Loc core.TextRange{pos, end int32}` + `id atomic.Uint64` +
  `Parent *Node` + `data nodeData` (an interface pointing back to the enclosing typed struct, the same
  "self" pattern jsgo uses). `NodeFactory` keeps one `core.Arena[T]` per node kind (a slice that
  doubles up to 256 elements per chunk). Identifier text is a Go `string`. So they use typed slabs
  with pointers, much like jsgo today. They avoid the long-lived-heap problem mainly because tsc does
  not run ~1,000 rewriting passes over a heap that stays live.

---

## 2. Recommendations for jsgo, ranked by expected impact on peak live heap / pointer count

Labels: **[in flight: X]** means an implementation track already covers the area and this item adds
external evidence or a specific design. Gains are estimates from the context numbers and
`docs/gclab-2026-10-18.md`. None of them are measured here.

### R1. Make the long-lived AST GC-free ("software pointer compression") and then cash it in on RSS [in flight: rhino.Node side tables / pointer-free AST]

- **Evidence.** V8 pointer compression (base + uint32 offset, up to 43% less heap), Turboshaft
  `OpIndex`, SpiderMonkey Stencil, and the table above: Java's AST bytes are not smaller, so jsgo's
  extra cost is the scanning.
- **Applies to.** `internal/rhino/node.go` (`Node`, `nodeExtBox`, `Node_PropListItem`,
  `nodeExtraShards`), node slabs and clone paths, `NodeTraversal`.
- **Mechanism, Go-adapted.**
  1. Remove every Go pointer from `Node`. Replace `typeOrColor` with a uint32 ID into a
     compilation-owned table (`[]JSType` by serial; colors already have pool IDs). Replace JSDoc
     and the ext box with one uint32 "aux" index into a `[]nodeAux` table. Most nodes have neither,
     so the type ID stays inline for typed builds and JSDoc/ext go to aux.
  2. Choose the link representation:
     (a) Keep the four `*Node` links in slabs typed as `[]byte`/`[]uint64`. This needs unsafe, keeps
     64 B nodes and keeps the API.
     (b) Use uint32 node indices with a process-wide two-level chunk table (`chunks[i>>k][i&mask]`)
     and the node's own index stored in it. `*Node` stays the API type because each accessor
     converts index→pointer. Node becomes about 56 B: 5×4 links + type ID 4 + aux 4 +
     linenoCharno 4 + length 2 + token 1 + kind 1 + payload 8 + propFlags 8. Because `Node` then
     has no pointer fields, `make([]Node, n)` is automatically noscan with no unsafe slab trick.
  3. Once node indices are dense, keep per-node side data in dense slices (Turboshaft
     `GrowingSidetable`; see R4) instead of maps keyed by address.
- **The point the context undersells: noscan reduces GC CPU but not RSS by itself.** Go's heap
  goal is `marked × (1 + GOGC/100)` and counts noscan bytes. RSS falls only through (i) fewer bytes
  (variant b: −12% of AST bytes; −25% if `previous` or the payload is moved out), or (ii) using the
  now-cheap GC cycles to lower GOGC or set GOMEMLIMIT. With the AST noscan, a cycle on bundle-mega
  marks about 20% of today's pointer-bearing bytes. GOGC=50 then doubles the number of cycles but
  still cuts mark work by about 60%, and the heap goal drops from 2.0x to 1.5x of live. On
  bundle-mega (700 MB live, RSS 1.3 GB) that could plausibly bring RSS to roughly 0.9-1.0 GB, before
  any byte savings.
  Alternative: an **off-heap cage** (one `syscall.Mmap` reservation, uint32 offsets exactly as in
  V8). It is excluded from the heap goal entirely, so RSS ≈ AST + 2×(the rest). It is not portable
  without platform-specific reservation code, and memory is reclaimed only when the whole compile ends.
- **Rough gain.** Mark work −45-60% on typed builds and about −85% on bundle (gclab shares). Live
  bytes −12-25% of the AST: typed-full −20-44 MB (AST about 177 MB), bundle-mega −65-155 MB. The RSS
  gain comes mainly from the GOGC lever.
- **Liveness cost (the Go version of "interior pointers keep things alive").** A type or JSDoc
  object referenced only by ID must be kept alive by its table. Tables therefore only grow until
  their lifetime ends: types until `typesToColors`, which matches what happens today, and JSDoc for
  the whole compile. Dead nodes are never reclaimed within a compile (gclab: about 90 MiB allocated
  against 35 MiB live on typed-app). That is acceptable, and it is what V8 does with Zones.
- **Pass API:** unchanged with variant (a) or (b). **Output risk:** low; it is a pure substrate
  change and the oracle tests cover it. **Effort:** high (substrate only). **Mode:** default.

### R2. Parser: offsets-only positions, transient tokens, inline ranges [in flight: parser memory]

- **Evidence.** V8 (`TokenDesc token_storage_[4]`, `Location{beg_pos,end_pos}`), SpiderMonkey
  (`TokenPos{uint32,uint32}`, `tokens[4]`, `SourceCoords` with `lastIndex_`), Hermes (`SMLoc`),
  Roslyn (struct tokens, singleton punctuation tokens),
  Kotlin K2 (LightTree).
- **Applies to.** `internal/parser/scanner.go` (`newToken`, `GetTokenRange`),
  `line_number_scanner.go` (`getTokenSourceRange` allocates 1-2 `SourcePosition` and one
  `SourceRange` per token), `token.go`, `identifier_token.go`, `literal_token.go`,
  `util_source_position.go`, `util_source_range.go`, `trees_parse_tree.go` (`Location *SourceRange`),
  `parse_arena.go`, `internal/parsing/ir_factory.go`.
- **Today.** Each token costs `TokenBase` 32 B (`self` interface 16 + `*SourceRange` 8 + `TokenType`
  int32 + padding) + `SourceRange` 16 B (2 pointers) + 1-2 × `SourcePosition` 32 B (pointer + three
  8-byte ints). That is 80-112 B, plus the value for identifiers (56 B) and literals. All of it holds
  pointers and stays in the ParseArena until IRFactory finishes. Parse trees add one more range and
  position per tree.
- **Proposal.**
  1. `type Range struct{ Start, End uint32 }` stored *inline* in `ParseTreeBase` and in token values
     (8 B, no pointers).
  2. Make tokens values in a small lookahead ring, sized to the parser's deepest `PeekTokenAt`.
     Store only what trees keep (identifier and literal: `{type uint8; Range; valueID uint32}` = 16 B)
     inside the tree. Punctuation tokens and keywords are never retained.
  3. Give each `SourceFile` a `lineStarts []uint32` (Java's `SourceFile` already has line offsets)
     with a "last line" cursor. IRFactory's monotone walk is then O(1) amortized. Compute
     `(line, column)` only in IRFactory's `setSourceInfo`, error reporting and `GetStart()`-style
     slow paths, which can still build a `*SourcePosition` for the few error messages.
  4. Drop the `Source *SourceFile` field from positions: one arena belongs to one file.
- **Gain.** Removes 80-140 B of pointer-bearing memory per token for the whole parse. For an 8.9 MB
  file with about 1.5-2.5 M tokens that is very roughly 150-350 MB at the parse peak; measure it with
  `arena.tokens/positions/ranges.used` counters first. Trees save another 8-16 B each and lose a
  pointer.
- **Pass API:** untouched; only the parser and IRFactory change. **Output risk:** low. Columns must
  stay in UTF-16 code units, and line terminators must match Java's `LineNumberScanner` (CR, LF, CRLF,
  LS, PS). The offset→line/col function is deterministic. **Effort:** medium-high.
  **Mode:** default.
- **Maximal variant.** Represent the parse tree as a flat, pointer-free array (Kotlin LightTree or
  Turboshaft style). It needs a parser and IRFactory rewrite, so do it only if (1)-(4) are not enough.

### R3. Arena and pool lifetime discipline (Zone/LifoAlloc semantics adapted to Go; Roslyn caps) [in flight: removing cross-pass pools/arenas]

- **Rules to apply across `internal/jscomp`, `internal/jutil` and parser arenas.**
  1. An arena of pointer-bearing objects is not a Zone in Go: the GC still scans it, and one
     reachable element pins its whole chunk and everything that chunk points to. Keep an arena only
     if (a) its element type is noscan, or (b) it is created and dropped with one phase: per pass,
     never retained across passes (V8 discards zones at the end of each compile). Arenas kept across
     passes (Round 7 "inter-pass arena retention", Round 10 "scopeArena chunk retention") turn
     transient garbage into live, scanned heap at the peak. Undo them unless they are noscan.
  2. Separate lifetimes: V8 has a string zone and a parse zone, Hermes an identifier allocator and an
     AST allocator. Never put a compilation-lifetime object (a `Var` that a pass stores in a
     compilation-wide map, a JSDoc node) in a per-pass slab, or it pins the slab.
  3. Pools follow Roslyn: bounded (about 2×GOMAXPROCS), short checkouts, objects above a size cap are
     not returned (ArrayBuilder: capacity ≥128), and "the pool is not storage". Go's `sync.Pool` is
     cleared every two GC cycles; hand-rolled retained pools are not.
  4. Add per-arena accounting printed per pass under a flag, like V8's `--trace-zone-stats`, so the
     swarm can see which arena is live at the peak.
- **Gain.** It depends on how much retained arena memory sits at the peak; likely tens of MB on
  mega builds, plus large simplification (goal 2). **Risk:** none. **Effort:** low-medium.
  **Mode:** default.

### R4. Dense side tables instead of address-keyed maps (Turboshaft side tables, Roslyn ConditionalWeakTable + flag bit) [in flight: rhino.Node side tables]

- `nodeExtraShards [64]{RWMutex; map[uintptr]nodeExtraEntry}` (node.go) and similar per-node maps
  should become dense slices indexed by node index or aux index once R1 provides indices. That is
  noscan, needs no hashing or locking, and removes the hazard of a slab address being reused with a
  stale entry (today only the `kindFlagHasExtra` bit protects against it). Keep Roslyn's discipline:
  a flag bit in the node says "has side data".
- **Gain:** small bytes, simpler code. **Risk:** low. **Effort:** low once R1 exists.
  **Mode:** default.

### R5. Data-flow state: snapshots instead of copies (Turboshaft `SnapshotTable`) [in flight: CFG/data flow]

- **Problem.** gclab: `MustBeReachingVariableDef` copies an `ObjectLinkedHashMap` per CFG node,
  357 MiB at one site and 21% of bundle churn.
- **Option A** (in flight): number variables per function and use dense `[]int32` or bitset
  lattices. This is quadratic in (CFG nodes × vars) for huge functions: 10k nodes × 2k vars × 4 B
  = 80 MB.
- **Option B** (V8): a `SnapshotTable` keyed by var index. Each CFG node keeps a snapshot handle
  (4-8 B). Set and Get are O(1), and a merge costs the number of changes since the common
  ancestor. jsgo's existing `hamt_p_map.go` is the persistent-map alternative (O(log n), structural
  sharing).
- **Caveat.** Closure passes later query the state at arbitrary nodes (`getFlowStateAt(n).getIn()`).
  A SnapshotTable only has one snapshot open at a time, so random queries must reopen snapshots;
  persistent maps answer them directly. Keep Java's iteration order wherever output depends on
  iterating a lattice.
- **Gain:** large churn reduction on bundle and less peak in pathological functions.
  **Risk:** medium. **Effort:** medium. **Mode:** default.

### R6. Scopes and vars: what engines do after resolution [in flight: scope/var/reference simplification]

- jsgo's `Var` is already 40 B (`nameNode`, `input`, `scope`, `nameID`, `indexAndFlags`,
  `nameHash`), the same as Java. `Scope` is about 112 B against Java's 32 B plus a map. Engine
  patterns:
  1. Key lookups by interned ID (V8's `VariableMap` keys on the interned pointer; jsgo's `nameID`
     is the analog).
  2. Store a scope's vars contiguously in the creator's arena and keep only a `[lo,hi)` range in the
     Scope. Closure's `SyntacticScopeCreator` declares a scope's vars all at once, before creating
     inner scopes, so they are contiguous. That removes `varList` (24 B) per scope.
     entries pushed on scope entry and popped on exit. A lookup from the current scope is then one
     probe with no per-scope hash tables. Queries on non-ancestor scopes fall back to a linear scan of
     the range.
  4. For anything that persists across passes (the global scope, TypedScope in typed builds), think
     about SpiderMonkey's 4-byte `BindingName` arrays (ID + flag bits).
- **Gain:** mostly churn (scopes/vars are 8.5-11% of allocation), and 1-3% of typed peak through
  TypedVar/TypedScope. **Pass API:** Scope/Var methods unchanged. **Risk:** low-medium (var
  iteration order must be preserved). **Mode:** default.

### R7. Types: use Java's sizes as the budget [in flight: JSType struct shrinking]

- The table in §1.4 shows jsgo types at 2.0-3.2x Java. Contributors, in order: the two 16 B interface
  fields in `JSTypeBase` (`self`, `resolveResult`); 8-byte Go `int` fields where Java uses 4-byte
  `int`; separate `bool`s that Java packs; and per-type `*JSTypeRegistry` pointers.
- Engine patterns that apply: pack every boolean into one flags word;
  V8 and SpiderMonkey bitfields; SpiderMonkey keeping the rare `resolveResult`-style data out of
  line; Roslyn's flag bit plus side table.
- Specific options:
  1. `resolveResult` becomes a concrete pointer (8 B) or moves to a side table once the type is
     resolved; most types resolve to themselves, so a flag can mean "self".
  2. Move ints to int32 and bools to a flag word.
  3. Keep `PropertyMap` entries as values in a slice-backed ordered map. Java's
     `TreeMap.Entry` + `Property` costs 72 B per property, so a slice of 32 B values plus an index
     beats Java.
- **Gain:** types are 15% of typed peak; halving them gives about −7% of typed-full peak (−40 MB).
  The type-check peak is where Java wins today. **Risk:** medium (equality and hash semantics).
  **Mode:** default.

### R8. Store source text one byte per character when possible; tag short strings (V8 `is_one_byte`, SpiderMonkey `HasTwoByteChars` and static length-1/2 atoms, QuickJS `is_wide_char`)

- **Applies to:** `ParseArena.utf16Buf` and the scanner in `internal/parser`, and
  `rhino_string_pool.go`.
- **Source text.** If every UTF-16 code unit is below 0x100, index *i* in a Latin-1 `[]byte` is the
  same as UTF-16 index *i*, so Java's column semantics hold exactly. For pure-ASCII files, scan the
  original UTF-8 bytes with no decode at all. Fall back to `[]uint16` otherwise. Implement it with a
  generic scanner over `uint8 | uint16`. **Gain:** about 9 MB on the 8.9 MB bundle parse (−50% of the
  text buffer), noscan memory only. **Risk:** low. **Effort:** medium.
- **String IDs.** Encode 1- and 2-character identifiers (most names after renaming) and small
  integers directly in the 32-bit ID, with a tag bit and no pool entry, and pre-assign IDs for
  well-known strings ("prototype", "goog", "exports", ...). The memory gain is tiny. The real gain
  is fewer pool lookups and `GetStringID() == wkPrototype` comparisons. **Risk:** low, provided ID
  allocation order does not affect any output ordering; check whether anything iterates by string ID.

### R9. JSDoc: drop strings, materialize type expressions lazily, flush them late [in flight: JSDoc compaction]

- **Cheap part.** `JSTypeExpression{root *Node; sourceName string}` (`js_type_expression.go`) becomes
  `{root; sourceFileID uint32}`, saving 12 B and one pointer per type expression.
- **Lazy part** (V8 lazy source positions and bytecode flushing, TASTy lazy trees).
  At IRFactory time, run JsDocInfoParser in a "validate only" mode: parse each type expression into a
  scratch arena that is reset per comment, report exactly the same warnings, and keep only
  `(fileID, offset, len)`. Materialize the `*Node` tree on the first `GetRoot()`, positions included.
  V8 had to make its eager and lazy paths produce identical results, with a stress mode, and the same
  applies here: an oracle-test stress mode that materializes everything eagerly and compares.
  Optionally flush the trees after the last pass that reads them (after `typesToColors` in ADVANCED).
  Needs an audit of `GetType`/`GetParameterType`/`GetReturnType`/`GetBaseType`/`GetTypeNodes`...
  callers by pass order.
- **Honest gain.** The typed-full peak is *during* type checking, when every type expression has been
  materialized, so the peak barely moves. The benefit is in untyped bundle builds (if no optimization
  pass reads type expressions) and in the heap after checking. **Effort:** medium-high. **Risk:**
  moderate. Rank: last among the default-mode items.

### Lazy ideas: what applies to Closure's pipeline

| Idea | Verdict for jsgo |
|---|---|
| Lazy parsing of function bodies (V8/SM/Hermes) | **No** in default and fast mode. Every function is checked, normalized and optimized in the first passes, early errors and JSDoc warnings must be complete, and nothing is "never called". The only real preparser is dependency pruning, which already exists (`LazyParsedDependencyInfo`). |
| Lazily materialized red nodes (Roslyn) | **No** for the AST: passes need `GetParent()` everywhere and change the tree in place. The underlying principle, store the minimum and derive the rest, applies to positions (R2) and per-node derived data (R4 side tables). |
| Positions derived on demand (V8 lazy source positions) | **Yes** in the parser (R2). Not for `rhino.Node`: `linenoCharno` int32 + `length16` is already 6 B and is copied around by passes. |
| Compact stencil for runtime libraries | **Already done** (`runtime_libs.typedast`, lazily injected per file; keep it embedded so it never lands on the Go heap). |
| Compact stencil for default externs (~2 MB source, 128 files in `externs.zip`) | **Startup time only, not memory.** Passes read the externs AST late (RenameVars, PureFunctionIdentifier, property passes), so it must be materialized. TypedAST's `OptimizationJsdoc` is lossy, so typed builds cannot use it. Could be a `--fast_mode` startup option for untyped builds. |
| Lazy extern types | **No.** It changes registry contents and order, "property defined anywhere" checks and (dis)ambiguation inputs. |
| Lazy JSDoc type expressions | **Yes, modest** (R9). |
| Stage-2/3 TypedAST inputs decoded on demand | **Already done** (streaming decode from `LazyAst` bytes, entry removed after use). |

---

## 3. Fast-mode-only ideas, and ideas considered and rejected

### Possible in `--fast_mode` (relaxed but deterministic)

- **Skip JSDoc type-expression parsing in untyped builds** except for tags optimizations read
  (`@define` with its type, `@const`, `@nosideeffects`, `@export`, `@nocollapse`, `@noinline`,
  `@enum`, `@constructor`/`@interface`/`@record`, `@license`/`@preserve`). Diagnostics change
  (malformed-type warnings are lost), not code. This is the V8 Lite-mode pattern: don't allocate
  optional data. It needs the same audit as R9.
- **Externs from an embedded pre-serialized AST** for untyped builds (startup time only).
- **Per-node `length` and `originalName` only when a source map is requested.** About 2 B per node
  plus side-table entries; they matter only for some diagnostics and source maps. Low value.

### Rejected

- **Red/green persistent trees (Roslyn).** Closure's ~1,000 passes change the tree in place and rely
  on node identity and parent links. Persistence costs O(depth) allocation per edit and would mean
  rewriting every pass.
- **Hash-consing subtrees (Roslyn `SyntaxNodeCache`).** AST nodes have identity, parents and per-node
  positions. JSDoc type trees look like candidates, but `ClosureRewriteModule` and module renaming
  change type-name nodes in place, and positions differ for every occurrence.
- **Compressed pointers inside ordinary Go heap objects (V8's 2020 `kCompressGraphZone`).** A 32-bit
  offset is invisible to Go's GC, so it is only safe when it points into noscan or off-heap memory
  (R1). V8 itself went back to full pointers for zone graphs.
- **Weak references for recomputable data (Roslyn weak red children, recoverable trees).** Go has
  `weak.Pointer` since 1.24, but a batch compiler where every pass needs the data gains nothing, and
  collection timing becomes a source of variance.
- **Streaming parse→IR per top-level statement.** UMD bundles are one top-level statement, and
  IRFactory works top-down, so it doesn't help the case that hurts.
- **AST-less compilation (QuickJS)** and **per-compile lazy delazification (SpiderMonkey)** do not
  apply to a whole-program optimizer.
- **Off-heap mmap cage as the first step.** Workable, and it gives the best RSS because it is
  excluded from the heap goal, but it needs platform-specific reservation and manual lifetimes. Try
  on-heap noscan slabs plus a lower GOGC first (R1).

### Outside this cluster, worth measuring first

- **Go 1.26 Green Tea GC** (default since 1.26; 10-40% less GC overhead for programs that allocate
  many small objects). jsgo is on Go 1.24.7. Measuring it is a toolchain bump.
- **GOMEMLIMIT plus a lower GOGC once the AST is noscan** (the cost and benefit are argued in R1).

---

## Sources (primary)

- V8: [zone.h](https://github.com/v8/v8/blob/main/src/zone/zone.h), [ast.h](https://github.com/v8/v8/blob/main/src/ast/ast.h), [ast-value-factory.h](https://github.com/v8/v8/blob/main/src/ast/ast-value-factory.h), [scopes.h](https://github.com/v8/v8/blob/main/src/ast/scopes.h), [variables.h](https://github.com/v8/v8/blob/main/src/ast/variables.h), [preparser.h](https://github.com/v8/v8/blob/main/src/parsing/preparser.h), [preparse-data.h](https://github.com/v8/v8/blob/main/src/parsing/preparse-data.h), [scanner.h](https://github.com/v8/v8/blob/main/src/parsing/scanner.h), [source-position.h](https://github.com/v8/v8/blob/main/src/codegen/source-position.h), [compiler/node.h](https://github.com/v8/v8/blob/main/src/compiler/node.h), [graph-zone-traits.h](https://github.com/v8/v8/blob/main/src/compiler/graph-zone-traits.h), [turboshaft/index.h](https://github.com/v8/v8/blob/main/src/compiler/turboshaft/index.h), [graph.h](https://github.com/v8/v8/blob/main/src/compiler/turboshaft/graph.h), [operations.h](https://github.com/v8/v8/blob/main/src/compiler/turboshaft/operations.h), [snapshot-table.h](https://github.com/v8/v8/blob/main/src/compiler/turboshaft/snapshot-table.h), [sidetable.h](https://github.com/v8/v8/blob/main/src/compiler/turboshaft/sidetable.h), [maglev-ir.h](https://github.com/v8/v8/blob/main/src/maglev/maglev-ir.h); blog: [preparser](https://v8.dev/blog/preparser), [v8-lite](https://v8.dev/blog/v8-lite), [pointer-compression](https://v8.dev/blog/pointer-compression), [leaving-the-sea-of-nodes](https://v8.dev/blog/leaving-the-sea-of-nodes), [background-compilation](https://v8.dev/blog/background-compilation)
- SpiderMonkey: [Stencil.h](https://github.com/mozilla-firefox/firefox/blob/main/js/src/frontend/Stencil.h), [CompilationStencil.h](https://github.com/mozilla-firefox/firefox/blob/main/js/src/frontend/CompilationStencil.h), [ParserAtom.h](https://github.com/mozilla-firefox/firefox/blob/main/js/src/frontend/ParserAtom.h), [ParseNode.h](https://github.com/mozilla-firefox/firefox/blob/main/js/src/frontend/ParseNode.h), [Token.h](https://github.com/mozilla-firefox/firefox/blob/main/js/src/frontend/Token.h), [TokenStream.h](https://github.com/mozilla-firefox/firefox/blob/main/js/src/frontend/TokenStream.h), [SyntaxParseHandler.h](https://github.com/mozilla-firefox/firefox/blob/main/js/src/frontend/SyntaxParseHandler.h), [LifoAlloc.h](https://github.com/mozilla-firefox/firefox/blob/main/js/src/ds/LifoAlloc.h), [vm/Scope.h](https://github.com/mozilla-firefox/firefox/blob/main/js/src/vm/Scope.h); [SpiderMonkey newsletter 96/97](https://spidermonkey.dev/blog/2022/01/14/newsletter-firefox-96-97.html)
- Hermes: [ESTree.h](https://github.com/facebook/hermes/blob/main/include/hermes/AST/ESTree.h), [Context.h](https://github.com/facebook/hermes/blob/main/include/hermes/AST/Context.h), [PreParser.h](https://github.com/facebook/hermes/blob/main/include/hermes/Parser/PreParser.h), [IR.h](https://github.com/facebook/hermes/blob/main/include/hermes/IR/IR.h); QuickJS: [quickjs.c](https://github.com/bellard/quickjs/blob/master/quickjs.c)
- Closure (v20261015): [Node.java](https://github.com/google/closure-compiler/blob/v20261015/src/com/google/javascript/rhino/Node.java), [jstype/](https://github.com/google/closure-compiler/tree/v20261015/src/com/google/javascript/rhino/jstype), [AbstractVar.java](https://github.com/google/closure-compiler/blob/v20261015/src/com/google/javascript/jscomp/AbstractVar.java), [AbstractScope.java](https://github.com/google/closure-compiler/blob/v20261015/src/com/google/javascript/jscomp/AbstractScope.java), [TypedAstDeserializer.java](https://github.com/google/closure-compiler/blob/v20261015/src/com/google/javascript/jscomp/serialization/TypedAstDeserializer.java)
- Roslyn: [GreenNode.cs](https://github.com/dotnet/roslyn/blob/main/src/Compilers/Core/Portable/Syntax/GreenNode.cs), [GreenNode.NodeFlagsAndSlotCount.cs](https://github.com/dotnet/roslyn/blob/main/src/Compilers/Core/Portable/Syntax/GreenNode.NodeFlagsAndSlotCount.cs), [SyntaxNode.cs](https://github.com/dotnet/roslyn/blob/main/src/Compilers/Core/Portable/Syntax/SyntaxNode.cs), [SyntaxNodeCache.cs](https://github.com/dotnet/roslyn/blob/main/src/Compilers/Core/Portable/Syntax/InternalSyntax/SyntaxNodeCache.cs), [C# InternalSyntax/SyntaxToken.cs](https://github.com/dotnet/roslyn/blob/main/src/Compilers/CSharp/Portable/Syntax/InternalSyntax/SyntaxToken.cs), [ObjectPool`1.cs](https://github.com/dotnet/roslyn/blob/main/src/Dependencies/PooledObjects/ObjectPool%601.cs), [ArrayBuilder.cs](https://github.com/dotnet/roslyn/blob/main/src/Dependencies/PooledObjects/ArrayBuilder.cs)
- Scala 3: [Spans.scala](https://github.com/scala/scala3/blob/main/compiler/src/dotty/tools/dotc/util/Spans.scala), [Trees.scala](https://github.com/scala/scala3/blob/main/compiler/src/dotty/tools/dotc/ast/Trees.scala), [TreeUnpickler.scala](https://github.com/scala/scala3/blob/main/compiler/src/dotty/tools/dotc/core/tasty/TreeUnpickler.scala); Kotlin K2: [JetBrains K2 benchmarks post](https://blog.jetbrains.com/kotlin/2024/04/k2-compiler-performance-benchmarks-and-how-to-measure-them-on-your-projects/)
- typescript-go: [ast.go](https://github.com/microsoft/typescript-go/blob/main/internal/ast/ast.go), [core/arena.go](https://github.com/microsoft/typescript-go/blob/main/internal/core/arena.go)
- Go: [Go 1.26 release notes (Green Tea GC)](https://go.dev/doc/go1.26)
