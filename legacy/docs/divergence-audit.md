# Divergence audit (2026-04-23, commit `1cf7fea`)

This document audits every significant divergence between the jsgo Go
port and its Java reference (closure-compiler `v20260415`). It was
produced from a static read of `internal/` against
`.cache/closure-src/com/google/javascript/...`, cross-referenced with
the self-declared `// Port of …` headers, skeletons under
`internal/jscomp/skel/`, and the porting-discipline rules in
`TODO.md § "Porting discipline"` (lines 9-44).

The discipline rules we are auditing against:

- Preserve class/method/field names recognisable to a Java reader.
- Preserve Java dispatch order and control-flow structure.
- Preserve data-structure choices when they are observable.
- No Go-only code paths: if a helper/struct/check doesn't exist in
  Java, it doesn't belong.
- Every non-trivial `internal/` file carries a `// Port of
  <java-file> @ closure-compiler v20260415` header.

## Summary

| Category | Count |
|---|---|
| Go-only helpers with no Java counterpart (in `typecheck.go` alone) | 8 (confirmed); ~40 across `scope.go` / `inference.go` / `superclasscheck.go` |
| Renamed methods in prod code (excl. `tvMismatch`/already fixed) | 3 confirmed (`buildRefinements`, `collectRefinements`, `inferType`) |
| Collapsed Java inner classes | 4 in `TypedScopeCreator` + 2 in `TypeInference` |
| Parser grammar shortcuts | Entire IRFactory layer (4352 Java LoC) skipped; Scanner stays char-level; ~21 tree categories unsupported |
| Un-ported Java skeletons staged in `skel/` | 95 (72 jscomp + 23 rhino) |
| Files with narrowed port qualifier (`narrow`/`subset`/`fragment`) | 21 of 208 `// Port of` headers |
| Go-only passes (no Java file at all) | 3 (`superclasscheck.go`, `default_externs.go`, `checkdeterministictest.go`) |

### Method-count sanity check

| Java source | Java LoC | Go file | Go LoC | Go LoC / Java LoC |
|---|---|---|---|---|
| `TypeCheck.java` | 3470 | `typecheck.go` | 3251 | 94% |
| `TypeInference.java` | 2849 | `inference.go` + `typeinferencepass.go` + `refinement.go` + `flowscope.go` + `dataflow.go` | 3600 | 126% |
| `TypedScopeCreator.java` | 3755 | `scope.go` | 3379 | 90% |
| `TypeValidator.java` | 1303 | `typevalidator.go` | 1004 | 77% |
| `parser/Parser.java` | 3694 | `parser/parser.go` | 3138 | 85% |
| `parser/Scanner.java` | 1375 | `parser/scanner.go` | 815 | 59% |
| `rhino/Node.java` | 3944 | `rhino/node.go` | 770 | 20% |
| `IRFactory.java` | 4352 | (none) | 0 | 0% |

Note: `typecheck.go` at 94 % of Java LoC still carries multiple
Go-only helpers (see §1). The low LoC ratio does not mean "faithful
port"; it means we replaced large methods with short helpers of our
own invention.

---

## Top 20 divergences (ranked by impact on 1:1 parity)

### 1. `rhino/Node`: slice-backed children vs Java's intrusive linked list

- **Where**: `internal/rhino/node.go:60-80` (Go), `.cache/closure-src/com/google/javascript/rhino/Node.java:1678-1681`
- **Java equivalent**: `Node.first` / `Node.next` / `Node.previous` —
  an intrusive circular doubly-linked list where every `Node` owns its
  own sibling links.
- **Why we diverged**: the header says "Java uses Node.first
  linked-list; slice view is easier to read." Self-flagged on
  `node.go:8-10`: *"The Java side uses an intrusive singly-linked
  list (Node.first + child.next), but the slice + explicit
  parent/next/previous maintenance below produces the same observable
  API."*
- **Behavioural risk**: Java's `Node.detach()` is O(1) (unlink from
  siblings); our `Node.Detach()` does a linear search in a slice.
  Traversal order is the same, but passes that detach/reattach
  children in hot loops (ES-rewriters, Es6ExtractClasses, etc.) will
  be O(n) per detach.
- **Migration sketch**: Drop the `children []*Node` slice, add
  `first/next/previous *Node` fields, rewrite `AddChildToBack /
  AddChildToFront / Detach / RemoveChild / Children()`, provide an
  iterator helper. ~300 LoC change, but blast radius is small
  because all children access goes through `n.Children()` or
  `FirstChild()/Next()`.

### 2. `parser/parser.go`: char-level scanner, no `TokenType` stream

- **Where**: `internal/parser/scanner.go:1-14`, `internal/parser/parser.go:10-13`
- **Java equivalent**: `parsing/parser/Scanner.java` (1375 LoC)
  emits `Token` records bearing a `TokenType` enum
  (`parsing/parser/TokenType.java`). `Parser.java` (3694 LoC) consumes
  the token stream via `Peek(TokenType)` / `Consume(TokenType)`.
- **Why we diverged**: self-flagged on `scanner.go:7-11`: *"This file
  carries the character-level scanning primitives the narrow jsgo
  grammar currently uses […] but does not yet emit a full Token type.
  A future expansion adds a TokenType enum + Peek/Consume(TokenType)
  interface so the Parser stops peeking at raw chars."*
  `parser.go:13-14` makes the same confession.
- **Behavioural risk**: every keyword lookahead is a `Peek()` over
  runes; whitespace skipping is hand-rolled (`SkipTrivia`), ASI rules
  are emulated via `skipTriviaTracked()` returning "seen newline"
  booleans rather than by checking a token's `precededByLineBreak`
  flag. Bug surface: identifier vs keyword disambiguation in
  `async` / `await` / `yield` / `of` / `static` / `set` / `get`
  contexts; Unicode identifier escapes (`Adentifier`); regex vs
  division ambiguity (currently resolved by `tryPrimary`'s context).
- **Migration sketch**: Port `parsing/parser/Token.java`,
  `TokenType.java`, `IdentifierToken.java`, `LiteralToken.java`,
  `TemplateLiteralToken.java`. Split Scanner into an `advance() Token`
  loop + a `peek() TokenType` cache. Rewrite Parser's 55 `tryX` / 
  `finishX` methods to consume tokens instead of bytes. Estimate:
  ~2500 Go LoC net change.

### 3. `parser/`: no `IRFactory` — direct `rhino.Node` construction

- **Where**: `internal/parser/parser.go` builds `rhino.Node` trees
  directly. No `trees/*Tree.java` ports, no `IRFactory.java` port.
- **Java equivalent**: 86 `*Tree` classes under
  `.cache/closure-src/com/google/javascript/jscomp/parsing/parser/trees/`
  plus `parsing/IRFactory.java` (4352 LoC) which translates the
  untyped parse tree into `rhino.Node`.
- **Why we diverged**: not explicitly flagged; the parser was a
  greenfield file, not generated from the skeleton
  (`internal/jscomp/skel/jscomp/parsing/parser/parser.go` exists but is
  un-wired).
- **Behavioural risk**: jsgo's parser does double duty as parser + IR
  lowering, which means every grammar production has been re-invented
  in the lowered shape. Bugs like "ITER_SPREAD emitted at wrong level
  in call arg list" can only be fixed by hand in the parser, not by
  tweaking IRFactory. Also blocks adoption of tree-level passes
  (feature detection, syntactic fixes) that Java runs on the parse
  tree before IRFactory.
- **Migration sketch**: Port all 86 trees as `type XxxTree struct {
  ... }` plus IRFactory visitor. Large (~5000 LoC) but mechanical.

### 4. `internal/jscomp/superclasscheck.go` — Go-only pass, not a
   separate file in Java

- **Where**: `internal/jscomp/superclasscheck.go` (1419 LoC)
- **Java equivalent**: none. Its diagnostics
  (`JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH` /
  `JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH`) are fired by
  `TypeCheck.checkDeclaredPropertyAgainstNominalInheritance`
  (line 1566) and `checkDeclaredPropertyAgainstPrototypalInheritance`
  (line 1665) — both inline in `TypeCheck.java`.
- **Why we diverged**: header on `superclasscheck.go:1-3` is honest:
  *"Port of com/google/javascript/jscomp/TypeCheck.java's
  checkPropertyVisibility / validateMethodOverride branches"* — but
  the port lives in a separate file with its own driver (`walk` +
  `checkClass` + `checkPrototypeAssign` + `checkInterfaceExtendsConflicts`
  etc.). The 15+ methods on `SuperclassCheck` have no direct Java
  counterparts.
- **Behavioural risk**: the super-check runs in its own AST walk
  instead of inline from `TypeCheck.visitClass` / `visitAssign`. Any
  mutation `TypeCheck` does to the scope during its walk won't be
  visible to `SuperclassCheck` (and vice versa), and interface
  validation runs twice (`TypeValidator.expectInterfaceProperties`
  during visitClass + `SuperclassCheck.validateImplements`).
- **Migration sketch**: Fold `SuperclassCheck` back into the relevant
  `TypeCheck.visit*` methods, matching Java's dispatch.

### 5. `TypeInference` is split into three files with non-Java method names

- **Where**: `internal/jscomp/inference.go` (`InferType`,
  `inferGetProp`, `inferCall`, …), `typeinferencepass.go` (`Transfer`,
  `transferDecl`, `transferAssign`, …), `refinement.go`
  (`buildRefinements`, `collectRefinements`).
- **Java equivalent**: single class `TypeInference.java` with a
  `traverse(Node, FlowScope)` dispatcher (line 552) and per-kind
  methods `traverseCall` / `traverseGetProp` / `traverseAssign` /
  `traverseReturn` / `traverseYieldAll` / `traverseCatch` /
  `traverseNew` / `traverseDeclaration`. **Every** Java method
  returns `FlowScope`; **our** `inferX` methods return `jstype.JSType`.
- **Why we diverged**: self-flagged on `inference.go:2-12`:
  *"Java's TypeInference is a ~2500-line control-flow-sensitive
  dataflow analysis […] This file ports the expression-level half:
  a flat node-kind → JSType switch […] this module stays monotonic
  and never consults a FlowScope."*
- **Behavioural risk**: expression-level narrowing (typeof/instanceof
  inside a conditional expression) can't propagate through
  sub-expressions the way Java's `FlowScope`-threading does.
- **Migration sketch**: rename `inferX` to `traverseX`, return
  `*FlowScope`, thread the flow scope through every call. Move
  `buildRefinements` into a proper
  `SemanticReverseAbstractInterpreter` port (skeleton exists:
  `internal/jscomp/skel/jscomp/typeinfer/semanticreverseabstractinterpreter.go`).

### 6. `TypedScopeCreator`: four inner classes collapsed into free
   functions

- **Where**: `internal/jscomp/scope.go:1-16` (self-flagged:
  *"collapses the inner-builder hierarchy into free functions on
  TypedScopeCreator"*)
- **Java equivalent**: `TypedScopeCreator.java` line 758
  (`abstract class AbstractScopeBuilder`), 3127
  (`NormalScopeBuilder`), 3232 (`FunctionScopeBuilder`), 3497
  (`ClassScopeBuilder`), 3664 (`FirstOrderFunctionAnalyzer`).
  Each is a `NodeTraversal.Callback` implementation with its own
  `enterScope` / `visit` / `shouldTraverse` contract.
- **Why we diverged**: Go doesn't have inheritance, and the port
  pre-dates the recursive-descent shape. Functions like
  `createScopeIn` / `buildMethodScope` / `collectThisProperties`
  cover what Java's three inner classes do.
- **Behavioural risk**: per-block hoisting, closure-captured scope
  nuances, and the two-pass NormalScopeBuilder traversal (declare
  first, fill second) are only approximated. Fixtures needing
  declaration-before-use semantics or `let` TDZ checks will rely on
  VarCheck/VarRefCheck rather than scope builder ordering.
- **Migration sketch**: Port the three builders as Go structs with
  matching methods. Keep a `ScopeBuilder` interface if callers want
  polymorphism.

### 7. `TypeCheck` Go-only helpers with no Java counterpart

All live in `internal/jscomp/typecheck.go`:

| Go helper | file:line | Java equivalent |
|---|---|---|
| `buildGlobalPropNameSet` | 1819 | none — Java consults `JSTypeRegistry` per-access |
| `buildExternTypeSet` | 1873 | none — Java handles externs via `isFromExterns()` on `JSType` |
| `buildLoosePropertyIndex` | 1901 | none — Java's `ObjectType.canPropertyBeDefined()` does this on demand |
| `dereferencedObject` | 2305 | sort-of: `JSType.toMaybeObjectType()` but simpler; ours is a standalone helper that only handles `TemplatizedType` / `ObjectType` |
| `isKnownSurfaceBuiltin` | 2264 | none — hardcoded switch for "Array" only |
| `checkPossibleInexistent` | 1754 | Java emits `POSSIBLE_INEXISTENT_PROPERTY` from `checkPropertyAccessHelper` (line 2113) |
| `defineOnOpenObject` | 2277 | Java calls `ObjectType.defineInferredProperty` directly |
| `isSafeCalleeForArity` / `isSafeInstanceMethodForArity` / `isDefinitelyNonConstructor` | 1136, 1181, 1209 | none — ad-hoc arity-check helpers |
| `insideOwnConstructor` | 1931 | none |

TODO.md § 3.3 already flags the first five. Total Go-only helpers
in `typecheck.go`: ~15. None of these survive a strict 1:1 audit.

- **Why we diverged**: all added opportunistically to paper over
  spurious diagnostics that Java avoids through its richer
  `JSTypeRegistry` + `isFromExterns()` + `canPropertyBeDefined()`
  plumbing. Rather than porting that plumbing, we added side indexes.
- **Migration sketch**: port `JSTypeRegistry` (skeleton at
  `internal/rhino/jstype/registry.go` is already present at ~600
  LoC), wire `isFromExterns()` + `canPropertyBeDefined()` onto
  `ObjectType`, then delete all six helpers from `typecheck.go`.

### 8. `TypedScope` carries a refinements stack; Java's scope is
   immutable, refinements live on `FlowScope`

- **Where**: `internal/jscomp/scope.go:64-73`
  (`refinements []map[string]jstype.JSType`)
- **Java equivalent**: `AbstractScope` / `TypedScope` have **no**
  refinement stack. Narrowing lives on `FlowScope` /
  `LinkedFlowScope`.
- **Why we diverged**: self-flagged on `scope.go:64-72`: *"Mirrors
  the FlowScope refinement shim Java's TypeInference threads through
  its dataflow."*
- **Behavioural risk**: our `TypedScope.PushRefinement` /
  `PopRefinement` pattern is effectively a dynamic-scoped overlay;
  any code that looks up a name while a push is active sees the
  overlay. Java's FlowScope is a value (immutable, explicitly passed)
  — much safer, no "who pushed when" debugging.
- **Migration sketch**: Eliminate `TypedScope.refinements`. Thread
  `*FlowScope` through every `InferType` call (see #5).

### 9. `parser/parser.go`: no `async` / `await` / `yield` / generator
   grammar

- **Where**: `internal/parser/parser.go` — `grep "async\|await\|yield"`
  returns 0 hits in parser.go and 1 in scanner.go (just a comment
  listing them as "contextual keywords" but not implemented).
- **Java equivalent**: `Parser.java` has full coverage —
  `parseAsyncArrowFunction`, `parseAwaitExpression`,
  `parseYieldExpression`, `parseGeneratorDeclaration`, etc.
- **Why we diverged**: no header note; the gap is apparent from
  grep. Fixtures that use `async function` / `await expr` /
  `function* gen()` / `yield expr` presumably fail to parse or
  degrade to an unrelated token.
- **Behavioural risk**: high — any fixture with ES2017+ async
  primitives silently mis-parses.
- **Migration sketch**: Add `tryAsyncFunction`, `tryAwaitExpr`,
  `tryYieldExpr`, `tryGeneratorFunction` on Parser; mirror Java's
  `allowYield` / `allowAwait` contextual state bits.

### 10. `parser/parser.go`: no proper destructuring-default /
   rest-element parse for non-param contexts

- **Where**: `internal/parser/parser.go:461` (`parsePatternOrIdent`),
  `480` (`parseArrayPattern`), `546` (`parseObjectPattern`). These
  fire on param lists and `for (let {x} of arr)` but there's no
  full-LHS destructuring in assignment position.
- **Java equivalent**: `Parser.parseAssignmentPattern` handles
  arbitrary `{a: b = 1, ...rest} = expr`.
- **Why we diverged**: similar — the parser grew feature-by-feature.
- **Migration sketch**: converge the three pattern-parsers into one
  `parsePattern(allowIn, allowAwait, allowYield)` matching Java.

### 11. `SymbolTable` Go stub

- **Where**: `internal/jscomp/symboltable.go:16-20` (self-flagged)
- **Java equivalent**: `SymbolTable.java` (2406 LoC) backed by
  `StaticSymbolTable<S, R>` + `StaticSlot` + `StaticRef` +
  `StaticScope` in `rhino/`.
- **Why we diverged**: Go generics weren't used; the interfaces are
  stubbed with TODOs. Header says *"Full port"* but body says
  otherwise.
- **Behavioural risk**: SymbolTable is not wired to any downstream
  indexer or IDE integration, so this is latent.
- **Migration sketch**: port `rhino/StaticScope.go` /
  `StaticSlot.go` / `StaticRef.go` as concrete interfaces (skeletons
  already under `internal/jscomp/skel/rhino/`).

### 12. `checkdeterministictest.go` / `checkdupvardecl.go` —
   TypeCheck branches extracted to separate files

- **Where**: `internal/jscomp/checkdeterministictest.go` (598 LoC),
  `internal/jscomp/checkdupvardecl.go`
- **Java equivalent**: both live inside `TypeCheck.java` — the
  former in `visitCompare` → `testEquality`, the latter in
  `checkDuplicateVarDeclaration`.
- **Why we diverged**: extracted for locality during initial
  bring-up. `TODO.md:788` explicitly flags the decision: *"for
  locality. Merge back when TypeCheck.java is ported 1:1."*
- **Behavioural risk**: these pass sites run on a fresh AST walk;
  they don't see the flow scope `TypeCheck` is maintaining. Fixtures
  that need refined types for determinism fire might under-report.
- **Migration sketch**: inline both into `TypeCheck.visit*`.

### 13. `default_externs.go` — hand-curated, not a port

- **Where**: `internal/jscomp/default_externs.go` (783 LoC, hand
  written)
- **Java equivalent**: `CommandLineRunner` bundles a `.jar` of
  externs (`es6.js`, `html5.js`, `w3c_dom1.js`, …, ~15 files, tens of
  thousands of lines of declaration). `DefaultExterns.java` exposes
  them.
- **Why we diverged**: self-flagged on `default_externs.go:3`:
  *"(narrow; hand-curated subset)"*. Porting the real externs bundle
  requires a full JSDoc type-language parser + ES grammar to consume
  the externs `.js` files.
- **Behavioural risk**: any fixture that expects globals we haven't
  declared (`IArrayLike`, DOM types, ES6 builtins) gets
  `JSC_UNDEFINED_VARIABLE` or missing-property FPs.
- **Migration sketch**: ship the real externs files as embedded
  assets once the parser is capable of consuming them. Blocked on #2
  and #3.

### 14. `SemanticReverseAbstractInterpreter` skeleton, refinement.go
   approximation

- **Where**: `internal/jscomp/refinement.go:1-28` (header + entry
  points `buildRefinements` / `collectRefinements`)
- **Java equivalent**: `type/SemanticReverseAbstractInterpreter.java`
  (663 LoC) with `caseEquality`, `caseInstanceOf`, `caseTypeOf`,
  `caseAndOrNotShortCircuiting`, `caseAndOrMaybeShortCircuiting`,
  `caseIn`, `caseDelete`.
- **Why we diverged**: self-flagged — *"Our port turns the same
  condition shapes into a flat `name → JSType` map"*. Our approach
  discards the `FlowScope` return value and just mutates a map.
- **Behavioural risk**: transitive qualified-name narrowing
  (`x.y.z != null`) and nested equality (`(a == b) == true`) are not
  handled. Also: our map-based approach can't express partial-union
  narrowing that Java's SRAI produces via `restrictByNotNullOrUndefined`.
- **Migration sketch**: port the full SRAI (skeleton present:
  `skel/jscomp/typeinfer/semanticreverseabstractinterpreter.go`).

### 15. `TypeValidator` subset — missing methods

- **Where**: `internal/jscomp/typevalidator.go:10-15`
  (*"The subset ported here covers the public API surface Java
  exposes […]"*)
- **Java equivalent**: `TypeValidator.java` (1303 LoC, ~60 methods)
- **Missing methods**: `expectYieldType`, `expectSpreadIterable`,
  `expectAutoboxesToIterable`, `expectValidTemplateType`,
  `expectTypesMatchForLHS`, etc. (~20 expect* methods absent)
- **Why we diverged**: only methods callers need have been ported.
- **Migration sketch**: port each method as callers land.

### 16. `parser/parser.go` uses `Scanner.Save / Scanner.Restore` for
   speculative parse; Java uses a token-level lookahead buffer

- **Where**: `internal/parser/scanner.go:116-125` (`Save` / `Restore`)
- **Java equivalent**: `Parser.java` uses `peekType(int depth)` /
  `peekExpression()` — bounded token lookahead, no rollback.
- **Why we diverged**: self-flagged on `parser.go:5-7`: *"speculative
  parse paths implemented via Scanner.Save/Restore."*
- **Behavioural risk**: rolling back the scanner means any
  side-effect (line-count update, JSDoc capture, comment buffering)
  has to be reversed too — currently the scanner only rolls back
  position + line, so any captured side-effect is lost. Bugs in
  JSDoc association on backtracked paths are a known risk.
- **Migration sketch**: replace with bounded token lookahead once the
  Scanner emits a Token stream (part of #2).

### 17. `Scanner.go` keyword detection is `strings.HasPrefix` based,
   not trie-based

- **Where**: `internal/parser/scanner.go:204` (`AtKeyword`), 233
  (`ConsumeKeyword`)
- **Java equivalent**: `Keywords.java` precomputes a perfect-hash
  table from identifier strings to `TokenType`. `Identifiers.java`
  tracks reserved-word sets per strict/non-strict mode.
- **Why we diverged**: no header note; simplification.
- **Behavioural risk**: `strict`-mode reserved words
  (`implements`, `interface`, `let`, `package`, `private`, etc.) are
  not differentiated from non-strict. `await`/`yield`/`async` get
  special-cased ad hoc.
- **Migration sketch**: port `Keywords.java` as a Go `map[string]Token`.

### 18. `NodeTraversal` is a widened subset

- **Where**: `internal/jscomp/nodetraversal.go:2` (*"widened subset"*)
- **Java equivalent**: `NodeTraversal.java` (~1500 LoC) with
  `AbstractPostOrderCallback`, `AbstractPreOrderCallback`,
  `AbstractScopedCallback`, `AbstractShallowCallback`, etc.
- **Why we diverged**: the traversal abstractions are hand-inlined
  per caller (TypeCheck has its own `Visit` / `ShouldTraverse`).
- **Migration sketch**: port `Callback` interface + the four
  abstract-callback helpers once multiple passes share them.

### 19. `Node.jstypeOrColor any` (type erasure) vs Java's
   `Node.typei` field

- **Where**: `internal/rhino/node.go:73-77` (stored as `any`)
- **Java equivalent**: `Node.typei` is typed `@Nullable TypeI` —
  an interface common to `JSType` and `Color` (the colour-only
  late-pass representation).
- **Why we diverged**: *"Stored as `any` to avoid an import cycle
  (jstype imports rhino)."* — real, genuine Go limitation.
- **Behavioural risk**: every `SetJSType` / `GetJSType` site casts;
  type confusion would not be caught by the compiler.
- **Migration sketch**: introduce a `rhino.TypeI` interface in the
  `rhino` package, have both `JSType` and (future) `Color` satisfy
  it. Requires splitting `jstype` to depend on an interface file
  that sits in `rhino/` or in a neutral `typei` package.

### 20. `checkjsdoc.go` omits Java tokens not present in the Go AST

- **Where**: `internal/jscomp/checkjsdoc.go:1-5` — *"Java-only tokens
  MEMBER_FUNCTION_DEF / GETTER_DEF / SETTER_DEF / COMPUTED_PROP /
  CAST / EXPORT / WITH / REST / DEFAULT_VALUE are absent from the Go
  parser, so those branches are omitted and noted with TODO
  comments."*
- **Java equivalent**: `CheckJSDoc.java` full
- **Why we diverged**: these tokens don't exist in our parser (see
  #2, #3, #10). So the JSDoc checks that would fire on them can't
  fire.
- **Migration sketch**: add the missing tokens to
  `internal/rhino/token.go` and to `parser.go`, then uncomment the
  check branches.

---

## Parser divergences (dedicated section)

### Approach

**How we scoped down.** Three distinct shortcuts stack:

1. **No `TokenType` layer.** Scanner exposes rune-level primitives
   (`Peek`, `Advance`, `ConsumeByte`, `ConsumeKeyword`). The Parser
   `Peek()`s characters and runs `strings.HasPrefix`-style keyword
   matches. Java emits `Token` records bearing `TokenType` enum + lexeme.

2. **No `*Tree` layer.** The Parser emits `*rhino.Node` directly. Java
   emits 86 distinct `*Tree` classes; `IRFactory.java` (4352 LoC)
   lowers them to `Node`.

3. **Speculative parse via `Scanner.Save/Restore`.** `tryX` methods
   return `(n, ok bool)` and rollback the scanner on failure. Java's
   Parser uses bounded `peekType(depth)` and commits once a production
   is disambiguated — no rollback.

### Features genuinely not implemented

Evidenced by absence from `parser.go`:

- **`async` / `await` / `yield`** — 0 hits in parser.go (only
  scanner.go:193 comment listing them as "contextual keywords");
- **Generator functions (`function*`)** — no `tryGenerator` method;
- **Import / export statements** — no `tryImport` / `tryExport`;
- **Computed property keys on class / object literals** —
  `parseClassMember` (line 795) handles plain identifier / string
  keys but not `[expr]: …`;
- **Tagged-template literal in member / call position** — partial
  (there's a `2010-2014` comment about adding it);
- **Numeric separators** (`1_000_000`), **BigInt literals** (`1n`) —
  `ScanNumericLiteral` (scanner.go:501) doesn't handle them;
- **Private class fields** (`#foo`) — no handling in
  `parseClassMember`;
- **Decorators** (`@decorator`) — not in grammar;
- **Optional catch binding** (`try{}catch{}`) — `finishTry` requires
  a catch-name.

### Features half-implemented

- **Destructuring** — parameter / `for-of` LHS only, not arbitrary
  assignment LHS (#10);
- **Arrow functions** — basic form only; no `async` arrow;
- **Template literals** — plain and tagged but not template expression
  positions in every context;
- **Optional chaining** — parsed but optional-chain propagation rules
  (e.g. `a?.b.c` vs `a?.b?.c`) are not fully honoured in all contexts;
- **Strict-mode reserved words** — not differentiated from
  non-strict (#17).

### Going 1:1 with Java's `Parser.java` would require

1. Port `parsing/parser/TokenType.java` (~90 enum values) to
   `internal/parser/tokentype.go`.
2. Port `parsing/parser/Token.java` + `IdentifierToken.java` +
   `LiteralToken.java` + `TemplateLiteralToken.java` (~4 Go files,
   ~200 LoC each).
3. Port `parsing/parser/Keywords.java` as a Go `map[string]TokenType`.
4. Rewrite `Scanner` to emit `Token` records in a `next()/peek()` loop
   (~1500 Go LoC).
5. Port all 86 `*Tree` classes under `trees/` to Go structs (~2000
   LoC, mechanical).
6. Port `IRFactory.java` (~4000 Go LoC).
7. Rewrite `parser.go` as a `Token`-stream consumer with bounded
   lookahead (~4000 Go LoC, replacing the current 3138).

Total: ~12-13k LoC of new code, replacing ~4000 LoC of current code.
Net +8k LoC but unlocks full ES coverage + tree-level passes.

---

## Class-collapse inventory

| Java inner class | Java source | Go collapsed into |
|---|---|---|
| `TypedScopeCreator.AbstractScopeBuilder` | `TypedScopeCreator.java:758` | free functions on `*TypedScopeCreator` in `scope.go` |
| `TypedScopeCreator.NormalScopeBuilder` | `TypedScopeCreator.java:3127` | `scope.go:createScopeIn` |
| `TypedScopeCreator.FunctionScopeBuilder` | `TypedScopeCreator.java:3232` | `scope.go:buildExpressionFunctionScope` + `buildMethodScope` |
| `TypedScopeCreator.ClassScopeBuilder` | `TypedScopeCreator.java:3497` | `scope.go:recordClass` + `collectThisProperties` |
| `TypedScopeCreator.FirstOrderFunctionAnalyzer` | `TypedScopeCreator.java:3664` | not ported |
| `TypedScopeCreator.IdentifyEnumsAndTypedefsAsNonNullable` | `TypedScopeCreator.java:704` | not ported; nullability is hand-rolled per-fixture |
| `TypeInference` class | `TypeInference.java` | split across `inference.go` + `typeinferencepass.go` + `refinement.go` |
| `LinkedFlowScope` | `LinkedFlowScope.java` | merged into `flowscope.go` |
| `TypeCheck.TypeCheckOptions` (enum) | inner | collapsed into `CompilerOptions` |

---

## Method-rename inventory

| Go name | file:line | Java original |
|---|---|---|
| `buildRefinements` | refinement.go:33 | `SemanticReverseAbstractInterpreter.getPreciserScopeKnowingConditionOutcome` |
| `collectRefinements` | refinement.go:39 | (implicit helpers inside SRAI) |
| `inferType` | inference.go:32 | `TypeInference.traverse` (but signature differs: returns JSType not FlowScope) |
| `inferCall` / `inferGetProp` / `inferAssign` | inference.go:335, 255 | `TypeInference.traverseCall` / `traverseGetProp` / `traverseAssign` — rename `infer` → `traverse` |
| `inferExpr` | typeinferencepass.go:262 | no Java equivalent; helper |
| `transferDecl` / `transferAssign` / `transferIf` / `transferWhile` / `transferDoWhile` / `transferFor` | typeinferencepass.go:70-224 | `TypeInference.traverseDeclaration` / `traverseAssign` / …; the `transferX` convention is ours. Java names them `traverseX`. |
| `Check` (on several pass types) | misc | Java uses `process(Node externs, Node root)` |

The `tvMismatch` → `mismatch` rename is done (#121). The
`inferX`/`traverseX` split is still pending.

---

## Un-ported Java skeletons under `internal/jscomp/skel/`

Total: **158** skeleton files staged but un-ported (per
`internal/jscomp/skel/README.md:8`).

Categories of interest:

- **Scopes**: `abstractscope.go`, `abstractvar.go`, `scope.go`,
  `scopecreator.go`, `scopedname.go`, `typedscope.go`,
  `typedscopecreator.go`, `typedvar.go`, `var.go` — the full
  `*Scope*` hierarchy is a skeleton, while our `scope.go` collapses
  it (see #6).
- **Type system**: `typecheck.go` (696 LoC skeleton, the "full" port
  target), `typeinference.go`, `typeinferencepass.go`,
  `typematchingstrategy.go`, `typemismatch.go`, `typetransformation.go`,
  `typevalidator.go`. Our current files are narrow subsets.
- **SRAI**: `typeinfer/chainablereverseabstractinterpreter.go`,
  `reverseabstractinterpreter.go`,
  `semanticreverseabstractinterpreter.go`,
  `closurereverseabstractinterpreter.go` — none ported; we have
  `refinement.go` as an approximation.
- **Parser**: `parsing/parser/{parser, scanner, token, tokentype,
  identifiertoken, literaltoken, templateliteraltoken}.go` — all
  skeletons; our current parser doesn't use them.
- **Passes**: `checkglobalthis.go`, `checksuper.go`, `passconfig.go`,
  `defaultpassconfig.go`, `passfactory.go`, `passlistbuilder.go`,
  `nodetraversal.go`, `scopecreator.go`, etc.

---

## Divergent token / structural semantics

| Divergence | Go | Java |
|---|---|---|
| Child list | `children []*Node` slice | intrusive circular linked list (`first`, `next`, `previous`) |
| Type storage on Node | `any` (type-erased) | `TypeI` interface |
| Prop bag | `map[Prop]any` | `IntObjectHashMap` (or a bit-flag int for prop bits) |
| Source file reference | filename `string` | `StaticSourceFile` reference |
| JSDoc attachment | direct `*JSDocInfo` field | side-table on `Node.propertyValue` map |
| Position | `lineno`, `charno` ints | packed `int` with line in high bits + char in low bits |
| Token enum size | ~80 tokens (`rhino/token.go`) | ~140 tokens |
| `TokenType` for parser | not separate (shares `rhino.Token`) | separate `TokenType` enum |

---

## Self-flagged divergences (comment census)

`grep -rn "TODO(\|divergence\|shortcut\|Go-only" internal --include=*.go`:

- **57** lines matching `TODO(...)` outside `skel/` (of which ~40 are
  our own, the rest Java-inherited `johnlenz`/`nickreid`/`b/...` TODOs).
- **0** lines matching `"Go-only"` (but the concept is all over:
  `go-port`, `port:`, `tier-N`, `later-tier`, `proxy:`, `Tier-5`,
  `registry-integration` prefixes).
- Most TODOs cluster in `internal/rhino/jstype/` (28 occurrences in
  `proxy_object_type.go` alone) and `internal/rhino/jstype/
  template_type_replacer.go` (5).

---

## Narrowed-port qualifier census

`grep -rn "// Port of.*(narrow\|subset\|fragment\|partial\|simplified)"`:

- **5** file-top hits (checkclosureimports, refinement, nodeutil,
  jsdoc_info_builder, jstype.go RestrictByTypeof region).
- `grep "(subset)" internal/*.go`: 30+ body-level comments. The
  distribution concentrates in checkjsdoc, checkaccesscontrols,
  nodetraversal, typevalidator, checkmissingrequires.
- Files with no qualifier but body-level "subset" confession:
  `typecheck.go` (header says "subset"), `scope.go`, `inference.go`,
  `superclasscheck.go`, `checkjsdoc.go` — i.e. all five of our
  largest type-check files admit to being subsets in the body, and
  three flag it in the header.

---

## Go-only code paths (passes with no Java source file at all)

| Go file | LoC | Java mapping |
|---|---|---|
| `internal/jscomp/superclasscheck.go` | 1419 | subset of `TypeCheck.java` extracted to its own file (#4) |
| `internal/jscomp/checkdeterministictest.go` | 598 | subset of `TypeCheck.java` extracted (#12) |
| `internal/jscomp/checkdupvardecl.go` | ~80 | subset of `TypeCheck.java` extracted (#12) |
| `internal/jscomp/checkduplicatemember.go` | ~200 | subset of `StrictModeCheck.java` extracted |
| `internal/jscomp/default_externs.go` | 783 | hand-curated (#13) |
| `internal/jscomp/astmanipulations.go` (part) | some | `AstManipulations.java` - ported, but our API shape is Go-idiomatic |

All are compiled into the default pass chain in
`internal/jscomp/compiler.go`.

---

## Recommendations

### Priority 1 — blocks all further 1:1 parity work

1. **Port `rhino.Node` to an intrusive linked list (#1)**. Every
   pass touches it; every structural divergence cascades from here.
   ~300 LoC, high-value.
2. **Replace expression-level `inferType` with a
   `traverse(Node, *FlowScope) *FlowScope` engine (#5, #8)**.
   Enables proper flow-sensitive narrowing. ~1500 LoC rewrite,
   touches inference.go + typeinferencepass.go + refinement.go +
   scope.go.
3. **Port `JSTypeRegistry` beyond the current stub (#7)**. Unblocks
   deletion of `buildGlobalPropNameSet`, `buildExternTypeSet`,
   `buildLoosePropertyIndex`, `isKnownSurfaceBuiltin`, and several
   sibling helpers in typecheck.go.

### Priority 2 — quality / maintenance

4. **Fold `superclasscheck.go` and `checkdeterministictest.go` back
   into `TypeCheck` (#4, #12)**. Mechanical. Keeps the dispatch tree
   in one file and frees `SuperclassCheck` for deletion.
5. **Port `SemanticReverseAbstractInterpreter` proper (#14)**.
   Skeleton exists; estimated 700 Go LoC.
6. **Port the four Java inner classes of `TypedScopeCreator` (#6)**.
   Start with `ClassScopeBuilder`. ~1200 Go LoC per builder.

### Priority 3 — larger investments, unblock strict-parity work

7. **Port the parser properly (#2, #3)**: TokenType enum, Token
   types, 86 Tree classes, IRFactory. Unblocks full ES coverage.
   ~12k Go LoC.
8. **Replace the hand-curated externs with the closure-compiler
   externs bundle (#13)**. Blocked on #7.
9. **Port `AbstractScope` / `StaticSlot` / `StaticRef` interfaces
   (#11)**. Unblocks `SymbolTable` being usable.

### Priority 4 — hygiene

10. **Adopt consistent narrowing qualifiers.** Only 5 `// Port of`
    headers declare `narrow`/`subset` explicitly, yet all five
    largest files are subsets. Either every narrow port flags itself
    in the header (and ci-guards.sh checks for it), or none do and
    the convention is dropped.
11. **Audit every `// TODO(later-tier)` / `// TODO(Tier-N)` against
    a tier plan.** 40+ such markers; half can probably be resolved
    now.

---

## Appendix: discipline breaches catalogued by CLAUDE.md / TODO.md rule

| Rule (TODO.md:9-44) | Breaches |
|---|---|
| *"Preserve class, method, and field names recognisable to a Java reader."* | `inferType` ≠ `traverse`; `transferDecl` ≠ `traverseDeclaration`; `buildRefinements` / `collectRefinements` aren't Java names at all; `Check` ≠ `process`. (§ Method-rename inventory) |
| *"Preserve the Java dispatch order and control-flow structure."* | `superclasscheck.go` runs as a separate traversal, not from inside `TypeCheck.visitClass` (#4). `checkdeterministictest.go` likewise (#12). |
| *"Preserve data-structure choices when they are observable."* | `Node.children []*Node` (#1). `TypedScope.refinements []map[]` (#8). `Prop` map instead of packed int flags. |
| *"No Go-only code paths."* | ~15 helpers in `typecheck.go` alone (#7). Entire `superclasscheck.go` pass (#4). `default_externs.go` hand-curated (#13). |
| *"Every non-trivial internal/ Go file carries a `// Port of <java-file>` header."* | Universally satisfied — 208 headers counted. |
| *"Narrowing qualifier is removed once the port matches Java 1:1."* | Only 5 files declare narrow/subset explicitly; a dozen large files silently elide the qualifier despite being subsets in the body. |
| *"No Go-only checks."* | `checkdeterministictest.go`, `superclasscheck.go`, `checkdupvardecl.go` — all three are TypeCheck subsets that grew into independent passes. |
