# Cluster report: JS/TS toolchains written in Go (go-js-toolchains)

Sources read: microsoft/typescript-go at 89d5d5b2 (20 Aug 2026; the repo is now archived and the code lives in microsoft/TypeScript, PR #4919), evanw/esbuild at f6058f83 (9 Aug 2026), the Go tree (cmd/compile/internal/ir, cmd/internal/src), the Go 1.24.7 runtime in the local GOROOT (malloc.go, msize.go, sizeclasses.go), and tsgo/esbuild/gopls PRs, changelogs and blog posts. Everything was read only: nothing was built or run, and /home/user/jsgo was not modified.

## TL;DR

1. **New finding, cheap and safe: size-class waste in jsgo's slabs.** Since Go 1.22, a pointerful object larger than 512 B gets an 8-byte malloc header, and that header is added *before* the size class is chosen. So `make([]rhino.Node, 128)` asks for 8192+8 B and lands in the **9472 B** class: 74 B per node instead of 64, about 13% of every IRFactory-slab node byte. The 16/32/64-node chunks waste 10–16%. typescript-go avoids this by growing arenas with `slices.Grow` and using the full rounded capacity (`internal/core/arena.go`). Fixing jutil.Slab and the other ~40 `make([]T, 2^k)` chunk sites is a few lines and cannot change output. Estimated saving: about 10–20 MB on typed-full and 30–60 MB on bundle-mega (2–8% of peak). The heap profile counts the full slot, so the waste is already inside the measured "AST 160 MB".
2. **Measure slab-retained dead nodes.** tsgo found 3.2% of nodes "abandoned" in its *immutable* AST (PR #4785) and added `ReleaseArenas` because arenas pinned dead nodes (PR #4024). jsgo's mutable AST, with ~1,000 passes deleting code, can only be worse. A 64-B node allocated on its own costs exactly 64 B in Go (heap bits sit in the span, no header). Slabs therefore save only malloc CPU and give up per-node reclamation. Under Go 1.26's Green Tea GC, only objects of 512 B or less get the new span-batched marking, which 8 KB slab chunks do not.
3. **Type and symbol references.** Follow tsgo's `*Type` header pattern and esbuild's index-typed references: store `*JSTypeBase` (8 B) instead of the `JSType` interface (16 B, 2 scanned words) in fields and slices (union alternates, template type lists, Property.type_, TypedVar.type_…). Use int32 indexes or uint32 string IDs, not pointers and strings, in TypedVar/Var/Reference.
4. **Parser target design** (in flight, so given as a reference): no token objects (tsgo `ScannerState`, esbuild `Lexer`), int32 start/end on parse trees (esbuild `Loc`, Go `XPos`), line and column computed later from a per-file line-start table (`[]int32`), and the source scanned as UTF-8 with a sparse table for non-ASCII characters (tsgo `PositionMap`, esbuild `LineOffsetTable`) instead of a full `[]uint16` copy. Closure computes line and column as a pure function of the offset, so this stays byte-identical.
5. **Side tables keyed by dense IDs** (tsgo `PagedLinkStore`: 256-entry `*[256]V` pages, links split per concern, a flag bit gating the lookup). This is the natural home for rhino's cold properties once nodes have dense int32 indices.
6. **Runtime levers.** Moving from Go 1.24.7 to 1.26 enables the Green Tea GC (10–40% less GC CPU per the release notes). GOGC/GOMEMLIMIT tuning trades CPU for RSS headroom. esbuild's "GC off" is the opposite of jsgo's goal.

---

## 1. Per-project techniques

### 1.1 microsoft/typescript-go ("tsgo", a function-by-function port of tsc to Go)

This is the closest analogue to jsgo. Its stated method is "a direct, file-by-file and function-by-function translation" ([10x faster TypeScript announcement and coverage](https://2ality-com.onrender.com/2025/03/typescript-in-go.html)), and its early memory was about half of tsc's. PR #581 measured `src/compiler` single-threaded: tsc 517,882K versus tsgo 267,819K. The team kept checker logic 1:1 and changed only the substrate. They did *not* attempt pointer-free ASTs or types, and their full VS Code check still holds about 3.5–4.7 GB live. The lessons below are the substrate changes that got them there.

**a) The common header is embedded in a kind-specific struct, with one interface word for downcasts.**
- `ast.Node` = `{Kind int16; Flags uint32; Loc TextRange{pos,end int32}; id atomic.Uint64; Parent *Node; data nodeData}` = 48 B, with 3 pointer words (Parent plus the 2-word interface). Source: [internal/ast/ast.go L179](https://github.com/microsoft/typescript-go/blob/main/internal/ast/ast.go). Each kind struct embeds `Node` first (`NodeDefault{Node}`), and `data` points back at the enclosing struct. This is the same trick as jsgo's `JSTypeBase.self`.
- Children are typed fields: `BinaryExpression{Left, Type, OperatorToken, Right *Node; …}` is 104 B. Lists are `NodeList{Loc; Nodes []*Node}`. There are no sibling pointers, only Parent. tsgo's per-node bytes are **not** smaller than rhino's uniform 64 B: an Identifier is 72 B and a BinaryExpression plus its operator-token node is 152 B. The win is fewer *structural* pointers per node (Parent plus the parent's child slot, about 2, against rhino's 4: first/next/prev/parent).
- `checker.Type` = `{flags TypeFlags; objectFlags ObjectFlags; id TypeId; symbol *Symbol; alias *TypeAlias; checker *Checker; data TypeData}` (types.go L681). **Every reference to a type is a `*Type` (8 B)**. Interface values exist only inside the header (`data`), and kind tests are flag tests. jsgo instead stores `JSType` interfaces (16 B) in ~60 fields and slices of `internal/rhino/jstype_*.go`.
- Discipline on struct size: PR #4744 removed duplicate embedded bases, and PR #4788 moved zero-size marker bases first, because Go pads a trailing zero-size field. That shrank 45 node structs by 8 B. The Go compiler enforces sizes with `cmd/compile/internal/ir/sizeof_test.go`. jsgo already does this for Node (`_ [64-unsafe.Sizeof(Node{})]byte`).

**b) Per-file, per-kind arenas that live and die with the AST, sized to Go size classes.**
- `core.Arena[T]{data []T}` grows 1→2→…→256 elements using `slices.Grow[[]T](nil, nextSize)`. The source comment reads "Grow rounds up to the next size class", and the arena then uses `len==cap` (see [internal/core/arena.go](https://github.com/microsoft/typescript-go/blob/main/internal/core/arena.go)). So every chunk fills its size class, including the 8-byte malloc header slot.
- `NodeFactory` holds about 48 arenas, one per *hot* node kind (ast_generated.go L20). Cold kinds are allocated individually. The parser is pooled with `sync.Pool`, but `putParser` zeroes the factory (`*p = Parser{scanner: …}`). Node arenas are therefore **never reused across files**: they belong to that file's AST. Scratch state (scanner caches) is reused.
- Small slices also come from arenas: `nodeSliceArena core.Arena[*ast.Node]` for child lists, and the binder's `singleDeclarationsArena` for the very common single-element `Symbol.Declarations` slice.
- **Lessons on interior pointers.** PR #4024 added `NodeFactory.ReleaseArenas()` "to break the GC link between nodes and the arena they came from": a cached node builder kept its current chunk alive. PR #4785 removed a speculative `async` lookahead whose discarded nodes "are never freed; they keep occupying slots in the source file's node arenas for as long as the AST lives". tsgo now reports "Nodes created / reachable / abandoned" (VS Code: 339,698 abandoned = 3.2%, down to 1.35% after the fix; −19 MB).

**c) "Links": per-concern side tables instead of fat nodes and symbols.**
- TS-JS keeps one big `NodeLinks`/`SymbolLinks` per entity. tsgo splits them into about 27 typed stores (`nodeLinks`, `signatureLinks`, `typeNodeLinks`, `enumMemberLinks`, `switchStatementLinks`, `valueSymbolLinks`, `aliasSymbolLinks`, …; checker.go around L674–700), so an entity pays only for the categories it actually uses. The emitter does the same: TS-JS's `node.original` and `node.emitNode` fields moved into `printer.EmitContext` maps (`original map[*Node]*Node`, `emitNodes LinkStore`), which live only for the emit (see [internal/printer/emitcontext.go](https://github.com/microsoft/typescript-go/blob/main/internal/printer/emitcontext.go)).
- The store type evolved with measurements:
  - `LinkStore[K,V]` = `map[K]*V` plus an arena of V.
  - PR #4226 made the two hottest stores **ID-indexed paged arrays**: −14% wall, −4.3% memory single-threaded. Sparse stores stayed as maps so memory would not grow.
  - PR #4329 made pages `*[256]V` (the page table shrank 66%), stored small V inline in pages and large V as arena pointers, and falls back to a map for high page indexes (`PagedLinkStore` in [internal/core/linkstore.go](https://github.com/microsoft/typescript-go/blob/main/internal/core/linkstore.go)).
  - PR #4741 allocates IDs in per-checker blocks so pages stay dense: −118 MB RSS at 4 checkers, −1 GB at 16.
  - PR #4703 stores links inline on checker-created (transient) symbols: one pointer chase and no separate allocation.
- **IDs are lazy.** `Node.id` and `Symbol.id` are assigned on first use (`ast.GetNodeId`), so only entities that ever get links consume ID space.
- **Flag bit plus side map.** JSDoc lives in `SourceFile.jsdocCache map[*Node][]*Node`, gated by `NodeFlagsHasJSDoc` on the node, so there is no JSDoc field on Node.

**d) Laziness: compute and keep only what is read.**
- **Lazy JSDoc in .ts files** (PR #2801). The comment is parsed on the first `(*Node).JSDoc()` call. Only tags needed for checking (`@see`, `@link`, `@deprecated`) are parsed eagerly, and JSDoc diagnostics from a lazy parse are thrown away. Server-style parse: −23% bytes on dom.generated.d.ts, −33% on a TSX test.
- **Line map only on demand.** PR #4787 removed the last check-path use of the line map, so a file's line table is now built only for diagnostics or LSP requests.
- **UTF-16 positions.** Source text stays a Go UTF-8 `string`, and identifier text is a zero-copy substring (`s.text[s.tokenStart:s.pos]`). UTF-16 offsets for the LSP come from `PositionMap`: an `asciiOnly` flag plus one `{utf8Pos, delta}` entry *per non-ASCII character*, with binary search (see [internal/ast/positionmap.go](https://github.com/microsoft/typescript-go/blob/main/internal/ast/positionmap.go)).
- **No token objects.** The scanner exposes one current token through `ScannerState{pos, fullStartPos, tokenStart int; token Kind; tokenValue string; tokenFlags; …}`. Lookahead and backtracking copy that struct by value (`Mark`/`Rewind`, [internal/scanner/scanner.go L192](https://github.com/microsoft/typescript-go/blob/main/internal/scanner/scanner.go)). Positions are stored as `TextRange{pos,end int32}` = 8 B per node.

**e) Pointer-free fixed-size cache keys.** TS-JS builds string keys for relation and instantiation caches, which is cheap in V8 thanks to ropes. tsgo first ported them as Go strings, where the profile was dominated by `WriteByte`. PR #2288 switched to 128-bit xxh3 hashes (`CacheHashKey`, 16 B, no pointers): VS Code allocations went from 34.6M to 25.6M (−26%) with memory −1%. PR #4784 builds keys in a 192-byte stack buffer and hashes once (average key 11.5 B). Relation caches are `map[CacheHashKey]RelationComparisonResult`, so both key and value are pointer-free and the map is noscan.

**f) Determinism without ordered maps.** tsgo uses plain Go maps (`SymbolTable map[string]*Symbol`). Wherever order is observable it sorts, for example `sortSymbols` by (file index, declaration position, name, id) in `checker/utilities.go L362`. It accepted baseline differences from TS-JS rather than emulate JS `Map` insertion order. `collections.OrderedMap` exists but is used only once in the checker.

**g) Memory hygiene PRs worth copying as patterns:**
- #259: pool Parser and Scanner, preallocate; parse allocations on checker.ts −81%.
- #1366: dom.d.ts allocations −77%.
- #4798: avoid temporary composite mappers; allocations −21%, peak RSS −5.5%.
- #2991: dedupe and replace a per-export `Set[string]` with bit flags; live heap 1.25 GB → 0.75 GB.
- #1745 and #2604: break retention cycles (a buildinfo↔program cycle, a writer pool that retained the program).

### 1.2 evanw/esbuild

esbuild optimizes for speed, not live heap. The CLI calls `debug.SetGCPercent(-1)` for one-shot builds (`cmd/esbuild/main.go` around L310–323), so its RSS equals everything it ever allocated. Its data *layout* is still the most pointer-frugal of the cluster:

- **Value-typed AST edges with no parent or sibling pointers.** `Expr{Data E; Loc logger.Loc}` and `Stmt{Data S; Loc}` are 24 B values held *inline* in parents. `E` and `S` are marker interfaces; `EBinary{Left, Right Expr; Op; IsParenthesized}` is 56 B. `a + b` costs about 88 B in 3 objects (rhino: 192 B in 3 nodes), and the only pointers are interface data words. Leaves without state are shared singletons (`ENullShared`, `EThisShared`, `EMissingShared`, `SEmptyShared`… [internal/js_ast/js_ast.go L570](https://github.com/evanw/esbuild/blob/main/internal/js_ast/js_ast.go)).
- **Symbols as flat per-file value arrays referenced by `Ref{SourceIndex, InnerIndex uint32}`** (8 B, pointer-free; [internal/ast/ast.go L385](https://github.com/evanw/esbuild/blob/main/internal/ast/ast.go)). The bundle-wide `SymbolMap{SymbolsForSource [][]Symbol}` is just an array of the per-file arrays, so merging is free, and cloning symbols for the linker is a slice copy with no rewiring. Identifiers carry a `Ref`, not a name. `Symbol.Link` gives union-find merging. `Index32` stores `^index`, so the zero value means "invalid", which makes zeroed memory a valid "nil".
- **Maps keyed by pointer-free keys.** `Part.SymbolUses map[Ref]SymbolUse`, `TopLevelSymbolToPartsFromParser map[Ref][]uint32`, and `ExprComments map[logger.Loc][]string` (a side table for rare data keyed by offset). Key and value types without pointers make the map noscan.
- **Positions as int32 byte offsets** (`logger.Loc{Start int32}`, `Range{Loc; Len int32}`). Line and column are computed only for messages and source maps. `LineOffsetTable{byteOffsetToStartOfLine; byteOffsetToFirstNonASCII int32; columnsForNonASCII []int32}` builds a UTF-16 column table **only for lines that contain non-ASCII** ("this mapping takes up a lot of memory and generates a lot of garbage… we avoid creating a table for ASCII-only lines", [internal/sourcemap/sourcemap.go L496](https://github.com/evanw/esbuild/blob/main/internal/sourcemap/sourcemap.go)).
- **Lexer without token objects.** One current token, and `Identifier MaybeSubstring{String; Start Index32}` aliasing the source. String literal escapes are decoded lazily into `[]uint16` only when needed (`decodedStringLiteralOrNil`). The source stays UTF-8.
- **Few full-AST passes and short-lived scope trees.** The scope tree is "not attached to the AST"; pass 2 replays `scopesInOrder` (see [docs/architecture.md](https://github.com/evanw/esbuild/blob/main/docs/architecture.md)).
- **History.** Release 0.6.19 cut memory 30–40% "by avoiding unnecessary per-file symbol maps" and by preallocating source-map arrays (CHANGELOG-2020.md, #304). esbuild moved to Go 1.26 and its Green Tea GC (CHANGELOG.md).

### 1.3 Go compiler (cmd/compile, cmd/internal/src)

- **Compact, pointer-free positions.** `src.XPos{index int32; lico uint32}` is 8 B. `lico` packs line (20 bits), column (8 bits, saturating), and is_stmt/prologue bits; the `*PosBase` pointer is replaced by an index into a table ([cmd/internal/src/xpos.go](https://github.com/golang/go/blob/master/src/cmd/internal/src/xpos.go), pos.go L319–343). This is the same idea as rhino's `linenoCharno int32` (20+12 bits), but the parser's `syntax.Pos{base *PosBase; line, col uint32}` already avoids per-position objects. jsgo's parser `SourcePosition{*SourceFile; Offset, Line, Column int}` is a separate 32-B object holding a pointer.
- **Small common header.** `ir.miniNode{pos XPos; op Op; bits bitset8; esc uint16}` is 12 B, embedded in kind structs. Sizes are pinned by `ir/sizeof_test.go` (`miniNode` 12 B, `Name` 160 B, `Func` 320 B).

### 1.4 gopls and go/parser

- **gopls v0.12 (2023) made typed syntax trees ephemeral.** It used to hold typed ASTs for all packages; afterwards it keeps compact export data and indexes in a file cache and keeps full ASTs only for open packages. Steady-state memory fell 497→232 MB on gopls itself, 3090→832 MB on kubernetes, and 5039→863 MB on google-cloud-go: about −75% on average across 28 repositories ([Go blog: gopls scalability](https://go.dev/blog/gopls-scalability), [v0.12 notes](https://golang.google.cn/s/gopls-v0.12)). The price was some operations about 50% slower until v0.13.
- **`parser.SkipObjectResolution`.** The parser stops allocating `ast.Object`/`Scope` bindings that the type checker recomputes anyway ([go/parser docs](https://pkg.go.dev/go/parser@master), [golang/go#52463](https://golang.google.cn/issue/52463)). The general lesson: do not build parse-time structures that a later pass rebuilds.
- *Skipped:* CUE and x/tools/go/ssa. I found no layout or GC write-ups there that bear on jsgo.

### 1.5 Go runtime facts that decide what pays off (verified in the local Go 1.24.7 source)

- **Small objects of 512 B or less** keep heap bits inside the span, with no per-object header. A standalone `new(rhino.Node)` costs exactly 64 B (size class 6, 128 objects per 8 KB span).
- **Pointerful objects of 513 B to 32,760 B get an 8-byte malloc header added before the size class is chosen** (`mallocgcSmallScanHeader`: `size += mallocHeaderSize`, malloc.go L1452). `roundupsize` knows this, so `slices.Grow` and `append` use the slack. **`make([]T, n)` with power-of-two byte sizes always spills into the next class.** The heap profile records the whole slot (`profilealloc(mp, x, size)` after rounding, L1511), so this waste is part of jsgo's measured AST bytes.
- **Objects above 32,760 B** are page-rounded, carry no header, and have no waste at exact page multiples (for example `[]any` of 2048 = 32 KB).
- **Green Tea GC** (opt-in in Go 1.25, default in 1.26): span-batched marking for small-object spans (objects of 512 B or less), "10–40% reduction in GC overhead", about 10% more with AVX-512 ([Go 1.26 release notes](https://go.dev/doc/go1.26), [golang/go#73581](https://laisky.notion.site/runtime-green-tea-garbage-collector-Issue-73581-golang-go-277ba4011a8681189344e40a749fa05f)). Larger objects, such as 8 KB slab chunks, keep the old per-object path.

---

## 2. Recommendations for jsgo, ranked by expected impact on peak live heap and pointer count

| # | What | Peak gain (est.) | Keeps *rhino.Node/pass APIs | Output risk | Effort | Mode |
|---|---|---|---|---|---|---|
| R1 | Slab chunks sized to size classes | typed-full 10–20 MB; bundle-mega 30–60 MB | yes | none | hours | default |
| R2 | Measure and bound slab-retained dead nodes; per-node allocation for 64-B nodes | 0–10% of AST (measure) | yes | none | days | default |
| R3 | `*JSTypeBase` instead of the `JSType` interface in fields and slices; drop `self` | 5–15 MB typed; fewer scanned words | yes | none | ~1–2 wk | default |
| R4 | Parser: no token objects, int32 ranges, lazy line/col, UTF-8 plus sparse UTF-16 delta | bundle parse peak −100–250 MB (est.) | yes (IRFactory reads offsets) | low if line/col derivation matches | in flight | default |
| R5 | Paged dense-ID side tables for cold node and symbol data | depends on cold-property prevalence | yes | none | medium (needs dense IDs) | default |
| R6 | Lean TypedVar/Var layout (esbuild symbol arrays) | maybe 10–40 MB typed (measure TypedVar count) | yes | none | medium (scope work in flight) | default |
| R7 | Pointer-free composite cache keys | allocation volume, not peak | yes | none if equality stays identity-only | small | default |
| R8 | Runtime: Go 1.26 Green Tea; GOGC/GOMEMLIMIT policy | GC CPU −10–40%; RSS/live 1.5–1.9 → about 1.3–1.5 | yes | none | small | default |

### R1. Size-class-aware chunking for every slab (from tsgo `core.Arena`) — do first

- **Files:**
  - `internal/jutil/slab.go` (`Slab.grow`: chunks of 16/32/64/128 elements).
  - `internal/rhino/node.go` `nodeCloneSlab.alloc` (exact subtree count, cap 512).
  - `internal/parser/parse_arena.go` (`arenaSlab.grow`, 128…8192 elements).
  - `internal/jscomp/scope.go` `scopeArena`.
  - Registry slabs in `internal/rhino/jstype_js_type_registry.go` (`unionTypeSlab`, `templateTypeMapSlab`, `templatizedTypeSlab`).
  - `rhino_string_pool.go` `ptrSlab` (1024 strings = 16 KB + 8 → 18,432 B class, 11% waste).
  - About 30 per-pass `make([]X, chunkSize)` sites (inline_functions.go, optimize_calls.go, inline_simple_methods.go, pure_function_identifier.go, serialization_*.go, infer_consts.go, remove_unused_code.go …).
- **Mechanism:** allocate chunks with `c := slices.Grow([]T(nil), n); c = c[:cap(c)]` (tsgo's trick), or pick `n` so that `n*sizeof(T)+8` fits a class (for Node: 127 nodes = 8136 B → 8192 class). A generic `jutil.NewChunk[T](n)` keeps call sites tidy.
- **Numbers for `[]rhino.Node`** (64 B, pointerful):

  | Request | Bytes | Class | Waste | Cost per node |
  |---|---|---|---|---|
  | 16 nodes | 1032 | 1152 | 10% | 72 B |
  | 32 nodes | 2056 | 2304 | 11% | 72 B |
  | 64 nodes | 4104 | 4864 | 16% | 76 B |
  | 128 nodes | 8200 | 9472 | 13% | 74 B |

  With `Grow`, the same classes hold 17/35/75/147 nodes, about 64.4 B per node.
- **Estimate:** if most of typed-full's 160 MB AST is IRFactory and JSDoc slab nodes, about 13% → **~15–20 MB (3–4% of the 552 MB peak)**. In untyped bundle-mega the AST is most of the 700 MB, so **~30–60 MB**. RSS falls correspondingly. Struct slabs with odd element sizes vary from 0–15%, so audit each one with `unsafe.Sizeof`.
- **Caveat:** clone slabs sized to an exact subtree count cannot use slack unless the slab persists across clones. Either leave them (about 5–10% average waste) or carry the remainder in a per-pass clone slab, which couples retention (see R2).
- **APIs, output and risk:** unchanged; zero risk.

### R2. Dead nodes retained by slabs: measure, then decide

- **Why:** a Go interior pointer to any node keeps its whole 8–9 KB chunk alive (tsgo #4024, #4785). jsgo's IRFactory allocates parsed nodes in 128-node chunks (`internal/parsing/ir_factory.go`, `jutil.Slab[rhino.Node]`). JSDoc nodes come from `jsdocNodeSlab`, and clones from `cloneTreeWithSlab`. Optimization passes then delete large amounts of code: early closure-primitive and module rewriting before the typed peak, and RemoveUnusedCode and inlining before the late peaks in ADVANCED bundles. tsgo saw 3.2% abandoned nodes with an *immutable* AST; jsgo's mutable AST can only be worse. Mitigating factor: parse order is source order, so deleting whole functions frees whole chunks.
- **Measure:** behind JSGO_MEMSTATS, at the peak phase and after the big deletion passes, count nodes reachable from the externs and js roots (including JSDoc type trees and nodes referenced from JSTypes) and compare with the nodes handed out by slabs, as tsgo's "Nodes created / reachable / abandoned" diagnostic does.
- **Mitigation if retention is above ~10%:**
  - (a) For node kinds or phases with high churn, allocate individually. A 64-B node has the same density as in a slab (no header, heap bits in span). The cost is about 15–25 ns of malloc per node, which is acceptable under "wall-time regressions are OK". Under Go 1.26 these objects also get Green Tea's batched marking, which 8 KB chunks do not.
  - (b) Smaller size-class-fit chunks (31 or 63 nodes).
  - (c) Never keep a long-lived slab or arena field in a long-lived owner (tsgo's `ReleaseArenas` lesson). This overlaps the in-flight "remove cross-pass pools/arenas" work.
- **APIs, output and effort:** unchanged; no output risk; measurement takes a day, mitigation is small.

### R3. Store `*JSTypeBase`, not the `JSType` interface (tsgo `*Type` header; esbuild and Go index references)

- **Files:** `internal/rhino/jstype_*.go` fields. Union `alternates/finalAlternates []JSType`, `templateTypes []JSType` (Named/Templatized/NoResolved), `TemplateTypeMap.templateValues/resolvedTemplateValues`, `Property.type_`/`owner`, `FunctionTypeBase.returnType/typeOfThis/implementedInterfaces/extendedInterfaces`, `JSTypeBase.resolveResult`, `EnumType.elementType`, `PrototypeObjectType.implicitPrototype`, `ProxyObjectType.referencedType`. Also `internal/jscomp/typed_var.go` `TypedVar.type_`. rhino.Node already does this (`typeOrColor unsafe.Pointer` to `*JSTypeBase`).
- **Mechanism:** fields hold `*JSTypeBase` (8 B, 1 scan word, against 16 B and 2 words). Getters return `b.self` (or a kind-switch cast, since jsgo already asserts `ObjectTypeBase` at offset 0 in jstype_object_type.go L165). A further step: replace the 16-B `self` interface with a kind byte plus a cast, which saves 16 B on every type. tsgo keeps its `data` interface, so this step goes beyond tsgo.
- **Estimate:** types are about 15% of typed peak (~83 MB on typed-full). Interface-valued fields and slices are a sizeable share, so expect roughly **5–15 MB plus about 16 B per type if `self` goes**, and fewer pointer words scanned on every GC cycle. Passes keep using the `JSType` interface (accessors convert), so they stay 1:1 with no output risk. Fold this into the in-flight JSType shrinking.

### R4. Parser and positions (tsgo scanner, esbuild lexer, Go XPos) — reference design for the in-flight parser work

- **Tokens:** `Token` becomes a ~16-B value struct `{type uint8; flags uint8; start, end int32; value rhinoStringID/literal index}` held in the scanner's lookahead ring (`[8]Token` by value; no `self` interface, no `*SourceRange`). Parse trees embed the values they need (IdentifierToken and literal tokens).
- **Ranges:** `ParseTreeBase.Location *SourceRange` (8 B pointer + 16 B range + 2×32 B positions) becomes inline `start, end int32`.
- **Line and column:** computed in IRFactory from a per-file `lineStarts []int32` (4 B per line), with a monotone cursor because nodes are produced roughly in source order. Closure's `SourceFile.getSourcePosition(offset)` already derives line and column from offset plus the line-offset table, so this is **byte-identical by construction**. The line terminator set and the column units (UTF-16 chars) must stay identical.
- **Text:** keep `Contents string` (UTF-8). Scan bytes, and turn byte offsets into Java char offsets with a tsgo-style `PositionMap` (an ASCII-only fast path plus one `{utf8Pos, delta}` entry per non-ASCII character) or esbuild's per-line `columnsForNonASCII`. That removes the `ContentsUTF16` copy (2 B per char; 17.8 MB for the 8.9 MB bundle). Invalid UTF-8 must map exactly as jsgo's current `StringToUTF16` decoding does, with U+FFFD counting as one unit (encode it in the delta table) — this is the main risk.
- **Estimate:** at roughly one token per 4–6 source bytes and about 112 B of token, range and position objects per token held in the per-file arena, the 8.9 MB bundle carries an estimated 100–250 MB of position and token objects at its parse peak. This design reduces that to almost nothing. It does not help the typed-full and typed-mega peaks, which occur during type checking.

### R5. Dense-ID paged side tables for cold and rare data (tsgo PagedLinkStore and per-concern links)

- **Targets:**
  - `nodeExtBox` plus the `Node_PropListItem` linked list (`internal/rhino/node.go` L596; each cold property is a boxed list item with an `any` value).
  - The global `nodeExtraShards [64]map[uintptr]nodeExtraEntry`.
  - Long-lived node-keyed maps that live across the type-checking peak (for example the `LinkedHashMap[*rhino.Node,*TypedScope]` memo in the typed scope creator, and `scopedNameTable LinkedHashMap[*Node, …]` in the registry).
- **Mechanism:**
  - One typed store per concern (`[]*[256]V` pages, created on first touch, falling back to a map for sparse high pages as in #4329).
  - Keyed by a dense node index and gated by a flag bit on the node (jsgo already has `hasExtBox` and the JSDoc bit).
  - Pointer-free V (ints, enums, string IDs) gives **noscan pages**.
  - Iteration over these stores is never observable, so Java order is not an issue as long as the stores replace lookups, not iterated collections.
- **Prerequisite:** dense IDs. These come free if the in-flight int32-indexed AST lands (index = ID; tsgo's #4741 shows IDs must be clustered). Without it, adding a field would grow Node from 64 to 72 B (80-B class), a net loss. A slab-relative ID (8 KB-aligned chunks of size class 8192, plus a per-64 MB-arena page table) is possible but depends on runtime internals. I would not do it.
- **Estimate:** proportional to the number of nodes carrying ext boxes, which is unknown, so measure the count. Typically 32–80 B per cold node drops to 4–8 B per property-page slot.

### R6. TypedVar/Var/Reference as lean value records (esbuild symbol arrays)

- **Mechanism:** Vars as values in per-scope or per-file slices (jsgo's `scopeArena` already slabs them). Name as `rhinoStringID` (4 B, not a 16-B string), `index int32`, `input` as an int32 input ID, type as `*JSTypeBase` (R3), scope as an int32 index into the scope slab.
- **Estimate:** about 88 B → about 40–48 B per TypedVar. Closure declares every qualified name (`a.b.c`) as a TypedVar in the global typed scope, so typed-full may have several hundred thousand to a million of them. Count them before committing. Keeps APIs, no output risk, coordinate with the in-flight scope work.

### R7. Pointer-free cache keys (tsgo #2288 and #4784) — allocation volume

- **Targets:** `EqualityChecker_CacheKey` and `SubtypeChecker_CacheKey` (heap-allocated `{left, right JSType; hashCode}` objects inside an ObjectLinkedHashMap) and similar per-lookup key objects. Use `map[[2]uint32]result` keyed by type serial IDs, so there is no key allocation and the map is noscan.
- **Semantic catch:** Closure's SubtypeChecker CacheKey falls back to non-identity equality for recursive templatized types ("Identity is insufficient in the case of recursive parameterized types"). Keep that fallback path, or the cache hits change. Apply the swap only where equality is pure identity.
- **Gain:** fewer allocations and less GC churn. These caches are transient, so the peak barely moves.

### R8. Runtime and GC policy

- **Go 1.26 (Green Tea).** jsgo's toolchain is go1.24.7. Expect 10–40% less GC CPU on small-object-heavy heaps (esbuild already moved). Benchmark it together with R2(a), because slab chunks larger than 512 B do not get Green Tea's span marking.
- **RSS headroom.** At GOGC=100 the heap goal is about twice the marked live heap. Measured RSS/live is 879/552 = 1.59 (typed-full), 1.53 (typed-mega) and 1.86 (bundle-mega). If RSS matters more than CPU, lower GOGC (for example 50) or set a GOMEMLIMIT just above the expected peak during parse and type-check. Every extra cycle re-marks the whole live heap, so pair this with the pointer reductions above. **Do not** copy esbuild's `SetGCPercent(-1)`: it makes RSS equal total allocation.

---

## 3. Fast-mode-only ideas, and ideas considered and rejected

### Fast mode (deterministic but not Java-identical)

- **Plain Go maps plus deterministic sorting where order is observed** (tsgo's `sortSymbols` by (file index, position, name, id)) instead of Java HashMap and LinkedHashMap order emulation. Collections fall to a Go map (or slice) with no slot array, no index and no order emulation. Collections are about 7% of the peak and the top allocation site. This changes output order, so fast mode only.
- **Lazy JSDoc type-expression parsing** (tsgo #2801). Keep the comment range plus a flag bit and build the rhino type-expression trees on first access. JSDoc is about 10% of the typed peak and likely more in UMD bundles full of library JSDoc. In default mode this would move or suppress JSDoc parse warnings and their order, so fast mode only unless warnings are proven absent or replayed in order.

### Considered and rejected

- **esbuild's GC-off CLI mode:** the opposite of the RSS goal.
- **esbuild's parent-less, value-typed Expr/Stmt trees and shared leaf singletons:** about 2× denser and far fewer pointers, but `*rhino.Node` passes need getParent/getNext/getPrevious, per-node line numbers and in-place mutation. Every pass would need rewriting.
- **tsgo per-kind node structs:** common kinds would be *larger* than rhino's uniform 64 B (Identifier 72 B, BinaryExpression 104 B plus a token node), and the change would break the API.
- **Dropping rhino's `previous` pointer** (tsgo and esbuild have no sibling links): O(1) last-child access and removal depend on it, and SCRIPT/BLOCK nodes with 10^5 children would make removal quadratic. Leave it to the int32-AST evaluation.
- **tsgo's lazily assigned `atomic.Uint64` id on every node:** pushes Node into the 80-B class. Only worth it if the ID comes free (int32 AST).
- **128-bit hashed keys for collections whose iteration order is observable:** they cannot reproduce Java order.
- **gopls-style serialize-and-drop of ASTs:** whole-program ADVANCED optimization needs the AST resident. The gopls analogue (compact summaries) is Colors plus "COLORS_AND_SIMPLIFIED_JSDOC", which jsgo already ports. A stage-1/stage-2 TypedAST split does not lower the type-checking peak.
- **tsgo's parallel checker pool:** each checker duplicates the type universe (+20% memory at 4 checkers in tsgo's own numbers). That is the wrong direction here.

### What Go lacks, and how these projects coped

- **No compressed oops:** answered with int32 or uint32 indices where references are dense (esbuild `Ref`/`Loc`, Go `XPos`, tsgo `TypeId`/NodeId-keyed pages), not with pointers.
- **No generational GC:** esbuild turns the GC off; tsgo accepts parallel marking on many cores (jsgo runs at t=1, so it cannot); gopls shortens object lifetimes.
- **Interior pointers retain whole chunks:** tsgo `ReleaseArenas`, removed speculation, and arenas scoped to the file's AST.
