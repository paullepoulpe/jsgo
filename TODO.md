# jsgo port TODO

Living list of what's left to port, with rough effort estimates, in
order of expected pass-rate lift for the `testdata/fixtures/extracted`
benchmark. Updated whenever a milestone commits.

**Current state:** 722/1444 = **50.00%** extracted pass rate. Tier-0
through Tier-3 hand-authored corpora all 100%.

**Goal:** 80%+ = ~1155 extracted passes. Gap: **+433**.

Since the last revision, the CFG, dataflow, externs, P1-2
(JSC_USELESS_CODE) and P1-4 (JSC_HIDDEN_SUPERCLASS / INTERFACE
PROPERTY MISMATCH) ports have all landed, plus a
nominal-nullability default matching Java's legacy mode.

**P0 sweep — completed (with the listed scope caveats):**
- P0-1 ✅ Parser precedence-climbing + expression-level TypeInference
  (binary / unary / ternary / array / object literals); `+` follows
  the JS string-bias rule; functions without @return default to
  Void; compound-assigns parse as ASSIGN_*. Scope: rest args /
  spread still skipped — testtemplatetype4/5 patterns stay
  unmatched.
- P0-2 ✅ VarCheck (JSC_UNDEFINED_VARIABLE), reports once per file,
  narrow extern-name stub. Scope: real DefaultExterns.js loader is
  still P2-1 — names not in the stub will FP if referenced bare.
- P0-3 ✅ VariableReferenceCheck (JSC_VAR_MULTIPLY_DECLARED_ERROR)
  for intra-file duplicates; honours `@suppress {duplicate}`;
  global frame seeded with extern-predeclared names so
  `Disposable`-style collisions match the oracle.
- P0-4 ✅ TemplateType + TemplateTypeMap + call-site unification;
  TemplatizedType (Array<T>, Map<K,V>, …) participates in subtype
  / property-lookup / unification / substitution. Scope: bounded
  generics (`@template T extends Base`) still produce
  JSC_UNSUPPORTED_BOUNDED_GENERIC_TYPES on the Java side and
  remain unported (Java itself defers — see P3-1).
- P0-5 ✅ **Real CFG + forward DataFlowAnalysis + TypeInferencePass**
  ported end-to-end:
  * `ControlFlowGraph` (controlflowgraph.go) with synthetic
    ENTRY/EXIT and branch-labeled edges.
  * `ControlFlowAnalysis` (controlflowanalysis.go) builds the CFG
    per FUNCTION / SCRIPT — proper IF with BranchOnTrue / OnFalse
    arms, RETURN edges straight to EXIT, reachability-aware joins.
  * `FlowScope` (flowscope.go) — persistent name→JSType map with
    Bind / BindAll / Equal and the lattice-join used at confluence
    points.
  * `DataFlowAnalysis` (dataflow.go) — forward-worklist runner that
    supports per-branch transfer outputs.
  * `TypeInferencePass` (typeinferencepass.go) — concrete transfer
    function: VAR/LET/CONST/ASSIGN refine the LHS NAME in the OUT
    scope; IF emits branch-specific scopes via `buildRefinements`
    (which itself recognises truthy, null-check, typeof, instanceof,
    `!cond`, `&&` / `||` falsy-branch). Expression inference in
    TypeCheck consults the IN FlowScope of the enclosing statement
    via `pushEnclosingFlow`.
  * `MissingReturnCheck` is now CFG-driven: walks `EXIT.Predecessors`
    restricted to reachable-from-ENTRY nodes and fires
    `JSC_MISSING_RETURN_STATEMENT` unless every one is a
    `return EXPR`. Replaces the AST-walk `definitelyReturns`
    approximation.

  Scope caveats still open: the CFG doesn't yet model loops (no
  BranchLoop / back-edge), try/catch/finally flow, labeled
  break/continue, or assignment-in-expression refinement (e.g.
  `if ((x = f())) { … }`). Those match the Java pass as the next
  wave of CFG edges.

## Methodology note

The buckets below are derived from grouping the 845 currently-failing
fixtures by the Java oracle's first diagnostic code (see
`scripts/gap-analysis` if you want to re-run the count). Each bucket
lists the Java files that port would touch. Item counts are upper
bounds — many fixtures fail for multiple reasons, so individual fixes
may move more.

---

## P0: biggest levers

### [P0-1] Expression-level type inference (`TypeInference.java`, narrow)
**Impact:** ~120–180 fixtures. **Effort:** 1 week.

Closure's `TypeInference` computes a JSType for every expression
node. jsgo currently handles only literals, NAME, THIS, GETPROP, NEW,
and direct CALL. Missing:

- Binary arithmetic: `+` (number+number → number, string+any → string),
  `-`/`*`/`/`/`%` → number.
- Comparison: `<`/`>`/`<=`/`>=`/`==`/`!=`/`===`/`!==` → boolean.
- Logical: `&&`/`||`/`??` → left | right union (with narrowing).
- Bitwise: `&`/`|`/`^`/`~`/`<<`/`>>`/`>>>` → number.
- Unary: `!` → boolean, `-`/`+`/`~` → number, `typeof` → string.
- Ternary `? :` → left | right union.
- Compound assignment `+=`/`-=`/`*=`/etc.
- Array literal `[…]` → `Array`.
- Object literal `{a: …}` → anonymous structural `ObjectType`.
- Function expression → `FunctionType` from @param/@return.

Ports: `TypeInference.java` (the non-dataflow part) plus
`TypeValidator.java`'s ADD / BITWISE / COMPARISON helpers.

### [P0-2] Undeclared-identifier check (`VarCheck.java`)
**Impact:** ~53 fixtures (all expecting `JSC_UNDEFINED_VARIABLE`).
**Effort:** 2 days.

Port of `VarCheck.java`: walk every NAME, confirm it's declared in
some scope (var/let/const, @param, function name, class name, imported
global). Straightforward; the hard part is populating the default
globals (`undefined`, `Infinity`, `NaN`, `globalThis`, etc.) from
externs.

### [P0-3] Duplicate-declaration check (`VarCheck.java` → `VariableReferenceCheck.java`)
**Impact:** ~26 fixtures (`JSC_VAR_MULTIPLY_DECLARED_ERROR`).
**Effort:** 1 day.

Walk declarations; flag a second declaration in the same scope or
across a scope-chain conflict. Small diagnostic port with large
dedicated Java test coverage.

### [P0-4] Real `@template` with call-site unification
**Impact:** ~30–60 fixtures. **Effort:** 1 week.

Current stand-in binds `@template T` names to `Unknown` so we don't
false-positive. A real port:

- `TemplateType` in `jstype` (port of
  `com/google/javascript/rhino/jstype/TemplateType.java`).
- `TemplateTypeMap` (substitutions).
- Unification at CALL / NEW sites:
  `visitCall` infers arg types, unifies against param `TemplateType`s,
  builds a substitution, applies to `@return`.
- `JSTypeResolver.resolveTemplates` port.

### [P0-5] Dataflow narrowing (real port of `TypeInference.java`'s control flow)
**Impact:** ~80–120 fixtures (many currently-suppressed
"found is union" cases). **Effort:** 2 weeks.

The big one. Java's `TypeInference` is a proper forward dataflow
analysis over the CFG, with `FlowScope` + `JoinOp`. Needed for:

- `if (x === null) { … }` narrowing `x` to non-null in the else branch.
- `typeof x === 'number'` narrowing.
- `instanceof Foo` narrowing.
- `if (!x)` narrowing unions.
- `JSC_MISSING_RETURN_STATEMENT` (exit-CFG reachability).
- `JSC_POSSIBLE_INEXISTENT_PROPERTY`.

Requires a control-flow-graph port (`ControlFlowGraph.java`,
`ControlFlowAnalysis.java`, `DataFlowAnalysis.java`). Biggest single
architectural piece still missing.

---

## P1: structural ports that enable checks

### [P1-1] Real `Parser.java` token stream
**Impact:** ~50 fixtures (parse-error + constructs we skip).
**Effort:** 1–2 weeks.

jsgo's parser is char-level; Java's is token-driven. Missing JS
constructs:

- Destructuring: `var {a, b} = obj`, `var [x, y] = arr`.
- Rest/spread: `f(...args)`, `function f(...args) {}`.
- Object literals at the expression level (only in `@type {…}` today).
- Template literals (`` `hello ${x}` ``).
- Regex literals `/.../flags`.
- `for`/`for-in`/`for-of`/`while`/`do`/`switch`/`try`/`throw`.
- `async`/`await`, generators.
- Optional chaining `a?.b`.
- Nullish coalescing already parses.
- ES modules (`import`/`export`).

Ports: `parser/TokenType.java`, `parser/Token.java`,
`parser/Scanner.java` (full), `parser/Parser.java` (full).

### [P1-2] Expression-statement side-effect check (`JSC_USELESS_CODE`)
**Impact:** ~57 fixtures. **Effort:** 3 days.

Port of `CheckSideEffects.java`. Detects expression statements with
no observable effect (`a + b;`, `x.y;`, `foo ? bar : baz;` in
statement position). Requires side-effect classification on CALL
(which needs an @nosideeffects scan).

### [P1-3] Nominal `@implements` check
**Impact:** ~15 fixtures (`JSC_INTERFACE_METHOD_NOT_IMPLEMENTED` +
subset of `JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH`).
**Effort:** 2 days.

When a class carries `@implements {I}`, walk I's structural member
set and verify the class has each member with a matching type. Ports
the interface-verification part of `TypeCheck.java:visitClass`.

### [P1-4] `@constructor` subclass property-mismatch check
**Impact:** ~16 fixtures (`JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH`).
**Effort:** 2 days.

When a subclass overrides a property inherited via `@extends`, check
the new type is a subtype of the inherited one. Ports the corresponding
branch of `TypeCheck.java:visitClass`.

### [P1-5] Property-access permission (`JSC_ILLEGAL_PROPERTY_ACCESS` + `ILLEGAL_PROPERTY_CREATION`)
**Impact:** ~17 fixtures. **Effort:** 3 days.

Ports `AccessControls.java` + the `@struct`/`@dict` / `@private` /
`@protected` check branches in `TypeCheck.java`.

### [P1-6] Function arity check (`JSC_WRONG_ARGUMENT_COUNT`)
**Impact:** ~7 fixtures. **Effort:** 1 day.

Requires parsing `@param` modifiers (`{=T}` optional, `{...T}` rest,
`{T=}` default) and comparing against the call-site arg count. Today
jsgo skips the arity check entirely to avoid false positives when a
param is optional.

### [P1-7] `@record` duck-typing over externs + inline records
**Impact:** ~10–20 fixtures. **Effort:** 3 days.

Currently we only match structural types against other named object
types. Extend to include: (a) anonymous object literal values, (b)
`Object` externs with extra ad-hoc properties. Requires literal-type
inference (P0-1) + extern population (P2-1).

---

## P2: infrastructure + correctness

### [P2-1] Real externs loader
**Impact:** ~40–80 fixtures (tests that assume `document.*`,
`window.*`, `Array.prototype.*` beyond our tiny stub).
**Effort:** 1 week.

Closure ships `DefaultExterns.js` + lots of additional extern files.
A real port:

- Wire `tools/extract_fixtures` to ALSO support `addExterns` tests
  (currently skipped — 270 tests).
- A `Compiler.compile` variant that accepts both externs and source.
- Extern type resolution: walk `@constructor` / `@interface` /
  prototype assignments in externs the same way we do in source.
- A minimal `Object.prototype` / `Array.prototype` / `String.prototype`
  built-in wired into the global scope.

### [P2-2] `ScopedCallback` + `enterScope` / `exitScope`
**Impact:** small. **Effort:** 2 days.

Port of the `ScopedCallback` sub-interface on `NodeTraversal`.
Enables cleaner scope handling in checkers that care (not strictly
blocking P0s but prerequisite for a cleaner scope-creator port).

### [P2-3] Inner `ScopeBuilder` classes in `TypedScopeCreator`
**Impact:** small. **Effort:** 1 week.

Java uses `AbstractScopeBuilder` / `FunctionScopeBuilder` /
`ClassScopeBuilder` inner classes that walk the tree via
`ScopedCallback`. jsgo collapses them into free methods on
`TypedScopeCreator`. Porting the inner-builder structure makes
closure-captured and per-block hoisting semantics easier to match.

### [P2-4] `Compiler.java` pass pipeline
**Impact:** 0 directly; unlocks multi-pass work.
**Effort:** 3 days.

Current `internal/jscomp/compiler.go` is a 30-line glue. A real
`Compiler` port gives us:

- `PhaseOptimizer` — pass registry + ordering.
- `CompilerInput` / `SourceFile` — so we can load multiple files and
  real externs.
- `Result` — a structured compilation report.
- `CompilerOptions` — parameterise checker behaviour.

### [P2-5] `Msg.java` + `DiagnosticType` registry
**Impact:** small. **Effort:** 2 days.

Replace the flat `const` diagnostic codes in `diag/diag.go` with a
registry that mirrors Java's `DiagnosticType.create(id, msgTemplate,
defaultLevel)`. Required for honouring `--jscomp_warning=<group>`
flags and lets us emit diagnostic messages that are
byte-for-byte identical to Java's.

### [P2-6] Intrusive linked-list children in `Node`
**Impact:** 0 (cosmetic). **Effort:** 2 days.

Replace the slice-backed `Node.children` with a `first` / `child.next`
intrusive list to match `Node.java` exactly. Zero behavioural impact;
matters for faithful port only.

---

## P3: lower-priority or out-of-scope buckets

### [P3-1] `JSC_UNSUPPORTED_BOUNDED_GENERIC_TYPES`
**Impact:** 38 fixtures. **Effort:** N/A.

Java itself says "currently still in development". Skip — matches
Java's own TODO.

### [P3-2] `JSC_NON_STRINGIFIABLE_OBJECT_KEY`
**Impact:** 15 fixtures. **Effort:** 2 days.

Specific to computed-property access on maps. Narrow.

### [P3-3] `JSC_DETERMINISTIC_TEST`
**Impact:** 9 fixtures. **Effort:** 3 days.

Dataflow-dependent lint — probably falls out once P0-5 lands.

### [P3-4] `JSC_USED_GLOBAL_THIS`
**Impact:** 5 fixtures. **Effort:** 1 day.

Scoped `this` lint; narrow.

### [P3-5] `JSC_TYPE_PARSE_ERROR` / `JSC_UNRECOGNIZED_TYPE_ERROR`
**Impact:** ~33 fixtures. **Effort:** 3 days.

Need error reporting on bad `@type` / `@param` annotations — port of
`JsDocInfoParser.java`'s warning paths (we currently silently drop
bad type expressions).

---

## Putting it together

Estimated fastest path to 80%+:

1. P0-1 (expression inference) → +120–180 fixtures
2. P0-2 + P0-3 (undefined / duplicate-decl) → +79
3. P1-2 (`JSC_USELESS_CODE`) → +57
4. P0-4 (@template) → +30–60
5. P1-4 (superclass property mismatch) → +16
6. P0-5 (dataflow narrowing) → +80–120

Cumulative: ~380–510 fixtures added → ~75%–85%. The dataflow port
(P0-5) is the single biggest lift and the last piece needed to
reliably clear 80%.
