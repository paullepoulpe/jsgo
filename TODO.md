# jsgo port TODO

Living list of remaining work, grounded in (a) per-file audit of jsgo
against the corresponding Java sources, and (b) bucket-by-oracle-code
triage of the 698 failing fixtures in the extracted corpus.

**Current state:** 1014 / 1444 = **70.22 %** extracted pass rate under
the default diagnostic set, 932 / 1444 = **64.54 %** under
`strictCheckTypes`. Tier-0 through Tier-3 hand-authored corpora all
100 %.

**Goal:** 80 %+ ≈ 1155 extracted passes. Gap: **+141** on default
(≈ 10 % of the corpus).

**Batch-merge state (Apr 2026):** 8 parallel ports landed (union-flatten,
templatized inheritance, NamedType position-tracked resolver,
`@nosideeffects`, `@private`/`@protected` visibility, optional/rest
param modifiers, `DUP_VAR_DECLARATION_TYPE_MISMATCH`, arity +
`NOT_A_CONSTRUCTOR`). Five of the new check passes — `JSTypeResolver`,
`CheckImplements`, `CheckAccessControls`, `CheckArity`,
`CheckDupVarDeclMismatch` — compile and have full fixture coverage but
are **not yet wired into the pipeline** (see `compiler.go` deferral
comments) pending false-positive cleanup on transitively-inherited
methods, varargs arity edge cases, and `@dict`/`@struct` bracket-access
patterns.

## Triage snapshot (current failing buckets)

Regenerate with `go run ./tools/fbucket`. Snapshot below taken against
the default diagnostic set at the 964-pass baseline.

```
193  want=[JSC_TYPE_MISMATCH]                     got=[]
 37  want=[JSC_UNSUPPORTED_BOUNDED_GENERIC_TYPES] got=[]   [Java's TODO]
 17  want=[]                                      got=[JSC_TYPE_MISMATCH]   [our FPs]
 15  want=[JSC_NON_STRINGIFIABLE_OBJECT_KEY]      got=[]
 11  want=[JSC_UNRECOGNIZED_TYPE_ERROR]           got=[]
 11  want=[JSC_INEXISTENT_PROPERTY]               got=[]
  7  want=[JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH] got=[]
  6  want=[JSC_WRONG_ARGUMENT_COUNT]              got=[]
  6  want=[JSC_TYPE_MISMATCH]                     got=[JSC_TYPE_MISMATCH]   [wrong line/col]
  5  want=[JSC_NOT_A_CONSTRUCTOR]                 got=[]
  5  want=[JSC_UNRECOGNIZED_TYPEOF_ERROR]         got=[]
  5  want=[JSC_UNDEFINED_VARIABLE]                got=[]
  5  want=[JSC_TYPE_MISMATCH,JSC_USELESS_CODE]    got=[JSC_USELESS_CODE]
  5  want=[JSC_INVALID_CAST]                      got=[]
  5  want=[JSC_TYPE_PARSE_ERROR]                  got=[]
  4  want=[JSC_TYPE_MISMATCH,JSC_USED_GLOBAL_THIS] got=[]
  4  want=[JSC_INCOMPATIBLE_EXTENDED_PROPERTY_TYPE] got=[JSC_INCOMPATIBLE_EXTENDED_PROPERTY_TYPE]  [wrong pos]
  4  want=[JSC_POSSIBLE_INEXISTENT_PROPERTY]      got=[]
  4  want=[JSC_VAR_MULTIPLY_DECLARED_ERROR]       got=[]
  4  want=[JSC_DETERMINISTIC_TEST]                got=[]
```

Since the last full refresh (833-pass baseline), the following buckets
have materially shrunk — thanks to the flow-narrowing push,
`checkdupvardecl.go`, and `checkjsdoctypes.go`:

- `JSC_TYPE_MISMATCH` (miss): 274 → 193  (−81)
- `JSC_UNRECOGNIZED_TYPE_ERROR`: 25 → 11  (−14)
- `JSC_INEXISTENT_PROPERTY`: 20 → 11  (−9)
- `JSC_INCOMPATIBLE_EXTENDED_PROPERTY_TYPE`: 9 → 0 missed (4 now emit but at wrong position)
- `JSC_DUP_VAR_DECLARATION_TYPE_MISMATCH`: 8 → out of top-20 (mostly fixed)
- `JSC_DETERMINISTIC_TEST`: 8 → 4
- `JSC_USELESS_CODE` standalone: 14 → out of top-20

## Audit snapshot (per-file completeness)

| File | Header | Gap | Effort | Impact |
|------|--------|-----|--------|--------|
| `internal/parser/parser.go` | narrow grammar | array/object literals at expression level, destructuring, spread/rest, template literals, regex, GETELEM, optional chaining, async/await, generators, full computed-property | 2–3 w | High |
| `internal/parser/scanner.go` | subset | Char-level only — no full TokenType enum, no string-literal/number-literal/regex tokenizers | 1–2 w | Medium (blocks parser) |
| `internal/jscomp/parsing/jsdocinfoparser.go` | scope-limited | @enum, @const, @dict/@struct, @nosideeffects, @abstract, @protected/@private, @param modifiers (`{=T}`/`{...T}`/`{T=}`), `function(A,B):C` callable expressions, type-error reporting | 3 d | Medium–High |
| `internal/jscomp/inference.go` | narrow | binary ops (ADD/SUB/MUL/DIV/MOD/BITWISE/COMPARISON), GETELEM, compound-assign result-type sharpening, missing literal cases | 1–2 w | High (~80–120 fixtures) |
| `internal/jscomp/typecheck.go` | subset | visitArrayLit / visitObjectLit, visitGetElem, visitFunctionExpression body, regex / template literal visit, ASSIGN-with-cast | 1–2 w | High |
| `internal/jscomp/typeinferencepass.go` | narrow | GETPROP refinement (`if (x.foo) { … x.foo … }`), assignment-with-union-LHS, falsy-on-union-with-literal-false | 3–5 d | Medium |
| `internal/jscomp/refinement.go` | narrow | nested AND/OR sub-expression refinement, GETPROP narrowing | 2 d | Medium |
| `internal/jscomp/scope.go` | subset | destructuring patterns, per-block hoisting, @enum slot type, Object.defineProperty, module imports | 1–2 w (blocked on parser) | Low–Medium |
| `internal/jscomp/controlflowanalysis.go` | narrow | labeled break/continue to **nested** loops, exception-edge fan-out from each call/throw to nearest catch, `with`, generator suspend, assignment-in-condition | 1 w | Low–Medium |
| `internal/jscomp/dataflow.go` | narrow forward-only | backward direction, lattice-annotation side table, fixed-point heuristic | N/A | None for current scope |
| `internal/jscomp/compiler.go` | shim | PhaseOptimizer, CompilerInput / SourceFile, multi-file Result, CompilerOptions | 1 w | Low |
| `internal/jscomp/varcheck.go` | narrow | per-scope unreached-NAME check, EARLY_REFERENCE for var hoisting | 3 d | Low |
| `internal/jscomp/checksideeffects.go` | narrow | @nosideeffects-aware CALL classification, sub-expression `a[b].c;` reporting | 2–3 d | Medium |
| `internal/jscomp/superclasscheck.go` | narrow | constructor-arg compatibility, multi-`@extends` interface diamond check (P0-A4) | 3–5 d | Medium |
| `internal/rhino/jstype/jstype.go` | subset | EnumType, NoResolvedType, NamedType (forward refs), distinct RecordType vs anon ObjectType | 1–2 w | Medium |
| `internal/rhino/jstype/object_type.go` | subset | @protected/@private/@public access bits, @struct vs @dict semantics, sealed-for-inference flag | 3 d | Low–Medium |
| `internal/rhino/jstype/function_type.go` | subset | optional/rest/default-value param tracking, THIS binding type, closure-capture environment | 3 d | Low–Medium |
| `internal/rhino/jstype/templatized_type.go` | narrow | inheritance through templatized supertypes, recursive substitution at nested generic positions | 3 d | Medium (template tests) |
| `internal/rhino/jstype/union_type.go` | narrow | union-over-union flattening optimisation, deterministic diagnostic ordering | 1 d | Low |
| `internal/rhino/node.go` | subset | `SideEffectFlags` cache, full ~50-entry Prop enum, cloneNode, intrusive linked-list children | 2–3 d | Cosmetic |
| **Java passes not started at all** | — | CheckMissingOverrides, CheckAccessControls (@private/@protected), CheckJSDoc (bad @typedef / @enum), SymbolTable, ConformancePass | weeks each | Mixed |

---

## Prioritised work list

Each item gives the relevant Java source(s), the audit gaps it
closes, the triage buckets it unblocks, an effort guess, and an
estimated fixture-rate lift on the extracted corpus. Items are listed
in expected ROI order — finishing the top half should clear the path
to ~75 %.

### P0-A — Expression inference completeness  *(highest ROI)*
**Java:** `TypeInference.java` (binary-op / GETELEM / compound-assign
branches), `TypeValidator.java` (ADD / BITWISE / COMPARISON helpers).
**Files:** `internal/jscomp/inference.go`, `internal/jscomp/typecheck.go`.
**Audit gaps closed:** binary ops in `inference.go`, GETELEM,
compound-assign sharpening, `visitArrayLit` / `visitObjectLit`,
ASSIGN-with-cast.
**Triage buckets unblocked:** ~120 of the 274 `JSC_TYPE_MISMATCH`
misses (the ones blocked on missing operator inference) plus most of
the 8 `[TYPE_MISMATCH,USELESS_CODE]` and 6 `[TYPE_MISMATCH][USELESS_CODE]`
mixed buckets.
**Effort:** 1–2 weeks. **Impact:** **+80 to +120 fixtures**.

### P0-A2 — `@nosideeffects` + better `CheckSideEffects`
**Java:** `CheckSideEffects.java` (full), `JsDocInfoParser.java`'s
`@nosideeffects` branch.
**Files:** `internal/jscomp/parsing/jsdocinfoparser.go`,
`internal/jscomp/checksideeffects.go`,
`internal/jscomp/superclasscheck.go` (call-effect classification).
**Audit gaps closed:** parser-level `@nosideeffects` + `@noside`,
sub-expression useless reporting (`a[b].c;`), call effect awareness.
**Triage buckets unblocked:** 14 `[USELESS_CODE]`-only +
8 `[TYPE_MISMATCH,USELESS_CODE]` + 7 `[DETERMINISTIC_TEST,USELESS_CODE]`.
**Effort:** 2–3 days. **Impact:** **+20 to +30 fixtures**.

### P0-A3 — `@enum` support
**Java:** `JSDocInfoParser.java` (@enum tag), `JSTypeRegistry.java`
(EnumType + EnumElementType creation), `TypeCheck.java` (enum-element
mismatch).
**Files:** `internal/jscomp/parsing/jsdocinfoparser.go`,
`internal/rhino/jstype/jstype.go` (EnumType + EnumElementType),
`internal/jscomp/scope.go` (recordEnum).
**Audit gaps closed:** `@enum`, EnumType class, enum-element subtype.
**Triage buckets unblocked:** ~25 of the `JSC_TYPE_MISMATCH` pool that
involve `@enum` declarations + the
`testaliasedenum*` / `testenumof*` clusters.
**Effort:** 4–5 days. **Impact:** **+25 to +40 fixtures**.

### P0-A4 — Multi-`@extends` on interfaces
**Java:** `JSDocInfoParser.java` (`@extends` accumulator),
`TypeCheck.java` (`JSC_INCOMPATIBLE_EXTENDED_PROPERTY_TYPE`).
**Files:** `internal/rhino/jsdoc.go` (extendsExpr → list),
`internal/jscomp/parsing/jsdocinfoparser.go` (parse multiple),
`internal/jscomp/superclasscheck.go` (diamond-check loop).
**Triage buckets unblocked:** 9 `[INCOMPATIBLE_EXTENDED_PROPERTY_TYPE]`.
**Effort:** 1 day. **Impact:** **+9 fixtures**.

### P0-A5 — `JSC_INEXISTENT_PROPERTY` broadening (without FP regression)
**Java:** `TypeCheck.java`'s `visitGetProp` non-structural branch.
**Files:** `internal/jscomp/typecheck.go`.
**Audit gaps closed:** non-structural property-existence check; needs
guarded broadening (the naive enable cost us 43 fixtures last attempt
— the win is conditional on an opt-out for builtins-with-empty-member
maps).
**Triage buckets unblocked:** 20 `[INEXISTENT_PROPERTY]`.
**Effort:** 3–5 days (most of the work is FP avoidance). **Impact:**
**+10 to +18 fixtures** (after FP balancing).

### P0-A6 — TYPE_MISMATCH FP reduction
**Files:** `internal/jscomp/inference.go`, `internal/jscomp/typecheck.go`,
`internal/jscomp/scope.go`.
**Audit gaps closed:** loop-aware dataflow (currently `while (!x) { x =
u(); }` doesn't refine `x` after the loop), GETPROP-refinement, narrow
assignment-with-union, missing function expression THIS-binding.
**Triage buckets unblocked:** 17 `want=[] got=[JSC_TYPE_MISMATCH]`
plus 8 of the 8 `[TYPE_MISMATCH,USELESS_CODE]` mixed bucket where we
also emit the right `USELESS_CODE`.
**Effort:** 4–5 days. **Impact:** **+15 to +25 fixtures**.

### P1-B — Cast expressions + parser breadth
**Java:** `Parser.java` (cast / array / spread / template grammar),
`Scanner.java` (full TokenType enum), `IRFactory.java`.
**Files:** `internal/parser/parser.go`, `internal/parser/scanner.go`.
**Audit gaps closed:** cast `(/** @type {T} */ (e))`, GETELEM
(`x[0]`), template literals (`` `text ${x}` ``), regex literals,
spread / rest, optional chaining.
**Triage buckets unblocked:** ~30 of the remaining
`JSC_TYPE_MISMATCH` pool (cast-blocked), 5 `[INVALID_CAST]`, parser
silent-failures across the corpus.
**Effort:** 2–3 weeks. **Impact:** **+40 to +60 fixtures**.

### P1-C — `JSC_UNRECOGNIZED_TYPE_ERROR` / `JSC_TYPE_PARSE_ERROR`
**Java:** `JsDocInfoParser.java` (warning paths in
`parseTypeExpressionAnnotation`), `TypeRegistry.java` (NamedType
late-resolve diagnostic).
**Files:** `internal/jscomp/parsing/jsdocinfoparser.go` (need
position-tracked warnings), `internal/rhino/jstype_expr.go`,
`internal/jscomp/scope.go` (resolve-time error emission).
**Audit gaps closed:** position-tracked errors inside JSDoc comments;
NamedType (forward-reference unresolved nominal).
**Triage buckets unblocked:** 25 `[UNRECOGNIZED_TYPE_ERROR]` +
5 `[TYPE_PARSE_ERROR]`.
**Effort:** 4–5 days. **Impact:** **+25 to +30 fixtures**.

### P1-D — Function arity + overload checks
**Java:** `TypeCheck.java`'s `JSC_WRONG_ARGUMENT_COUNT`, `JSC_NOT_A_CONSTRUCTOR`,
optional/rest/default-value `@param` modifiers in
`JsDocInfoParser.java`.
**Files:** `internal/jscomp/typecheck.go`,
`internal/rhino/jstype/function_type.go` (optional / rest / minArity),
`internal/jscomp/parsing/jsdocinfoparser.go`.
**Triage buckets unblocked:** 7 `[WRONG_ARGUMENT_COUNT]` +
7 `[NOT_A_CONSTRUCTOR]`.
**Effort:** 3–4 days. **Impact:** **+12 to +15 fixtures**.

### P1-E — `JSC_DUP_VAR_DECLARATION_TYPE_MISMATCH`
**Java:** `TypeCheck.java`'s `checkDuplicateVarDeclaration`.
**Files:** `internal/jscomp/varrefcheck.go` (extend with type-mismatch
arm) or new pass.
**Triage buckets unblocked:** 8 `[DUP_VAR_DECLARATION_TYPE_MISMATCH]` +
6 `[DUP_VAR_DECLARATION_TYPE_MISMATCH][UNDEFINED_VARIABLE]`.
**Effort:** 1 day. **Impact:** **+10 fixtures**.

### P2-A — `CheckAccessControls` (`@private` / `@protected`)
**Java:** `AccessControls.java` (full), `JSDocInfoParser.java` (visibility
tags), `TypeCheck.java`'s `@struct`/`@dict` branches.
**Files:** new `internal/jscomp/accesscontrols.go`,
`internal/rhino/jstype/object_type.go` (visibility bits),
`internal/jscomp/parsing/jsdocinfoparser.go`.
**Triage buckets unblocked:** 7 `[ILLEGAL_PROPERTY_ACCESS]`.
**Effort:** 1 week. **Impact:** **+8 to +12 fixtures**.

### P2-B — `JSC_INTERFACE_METHOD_NOT_IMPLEMENTED`
**Java:** `TypeCheck.java`'s `validateImplements` loop.
**Files:** `internal/jscomp/superclasscheck.go` (extend), or new
`internal/jscomp/checkimplements.go`.
**Triage buckets unblocked:** 5 `[INTERFACE_METHOD_NOT_IMPLEMENTED]` +
the rest of the 9 `[HIDDEN_SUPERCLASS_PROPERTY_MISMATCH]` we still
miss (class-style fixtures).
**Effort:** 3 days. **Impact:** **+8 to +12 fixtures**.

### P2-C — Real externs loader (`DefaultExterns.js` superset)
**Java:** `Compiler.java`'s extern bundle wiring,
`com.google.javascript.jscomp.resources.DefaultExterns.js`.
**Files:** `internal/jscomp/default_externs.go` (extend embedded
source), or load from `testdata/externs/`.
**Audit gaps closed:** `goog.*`, `Window.*`, `HTMLElement.*` deeply,
typed-array prototype methods, `String.prototype.*` complete signatures.
**Triage buckets unblocked:** 6 `[UNDEFINED_VARIABLE]` (`goog`, etc.) +
indirectly improves the false-positive rate on TYPE_MISMATCH.
**Effort:** 1 week (mostly content). **Impact:** **+10 to +20 fixtures**.

### P2-D — `JSC_DETERMINISTIC_TEST`
**Java:** `CheckDeterministicTest.java`.
**Files:** new `internal/jscomp/checkdeterministictest.go`.
**Triage buckets unblocked:** 8 `[DETERMINISTIC_TEST]` +
7 `[DETERMINISTIC_TEST,USELESS_CODE]`.
**Effort:** 3 days. **Impact:** **+15 fixtures**.

### P2-E — `JSC_NON_STRINGIFIABLE_OBJECT_KEY`
**Java:** `TypeCheck.java`'s `validateNonStringifiableKey` branch.
**Files:** `internal/jscomp/typecheck.go` (visit GETELEM /
COMPUTED_PROP).
**Triage buckets unblocked:** 14 `[NON_STRINGIFIABLE_OBJECT_KEY]`.
**Effort:** 2–3 days (gated on GETELEM parser support — P1-B).
**Impact:** **+14 fixtures**.

### P3 — Lower-leverage / out-of-scope items

* `JSC_UNSUPPORTED_BOUNDED_GENERIC_TYPES` — 37 fixtures, **skip**:
  Java itself flags this as still-in-development.
* `with`-statement modeling — rare; deferred.
* Generator / async suspend CFG edges — rare in checker corpus.
* SymbolTable / module graph — needs multi-file support first.
* PhaseOptimizer + CompilerOptions / Result — architectural debt;
  no direct fixture impact.
* Source-position tracking inside JSDoc type expressions — needed by
  P1-C; would also benefit suggestion-style diagnostics.

---

## Putting it together — fastest path to 80 %

Baseline is 964 / 1444. Cumulative columns show the pass total after
each item lands. Impact per-item is capped by the bucket size, which
has already shrunk since the previous revision of this list.

| # | Item | Effort | Item impact | Cumulative |
|---|------|--------|-------------|------------|
| 1 | P0-A expression inference completeness | 1–2 w | +60 to +90 (cap 193) | 1024–1054 |
| 2 | P0-A6 TYPE_MISMATCH FP reduction | 4–5 d | +12 to +17 (cap 17) | 1036–1071 |
| 3 | P2-E `NON_STRINGIFIABLE_OBJECT_KEY` (after P1-B) | 2–3 d | +12 to +15 (cap 15) | 1048–1086 |
| 4 | P1-C `JSC_UNRECOGNIZED_TYPE_ERROR` | 3–4 d | +8 to +11 (cap 11) | 1056–1097 |
| 5 | P0-A5 `JSC_INEXISTENT_PROPERTY` broadening | 3–5 d | +8 to +11 (cap 11) | 1064–1108 |
| 6 | P2-B `INTERFACE_METHOD_NOT_IMPLEMENTED` + hidden-superclass fixes | 3 d | +5 to +7 (cap 7) | 1069–1115 |
| 7 | P1-D arity / not-a-constructor | 3–4 d | +8 to +11 (cap 11 across 2 buckets) | 1077–1126 |
| 8 | P1-B parser breadth (cast / GETELEM / template) | 2–3 w | +30 to +50 (unblocks #3 + #7) | 1107–1176 |
| 9 | P2-D `DETERMINISTIC_TEST` | 3 d | +3 to +4 (cap 4) | 1110–1180 |
| 10 | P2-C real externs loader | 1 w | +5 to +15 (UNDEFINED_VARIABLE + TYPE_MISMATCH) | 1115–1195 |
| 11 | P2-A `CheckAccessControls` | 1 w | +3 to +5 | 1118–1200 |

Doing items 1–5 puts us at roughly **1064–1108 / 1444 ≈ 74–77 %**.
Items 1–8 cross the 80 % line (1155), with item 8 (parser breadth)
the single largest individual lift and a prereq for items 3 and 7.
Items 9–11 are gravy once the leverage from 1–8 is captured.

The remaining ~30–60 fixtures past that point need destructuring,
async/await, modules, and the longer-tail JSDoc tags — multi-week
ports each.

The remaining ~80–100 fixtures past that point need destructuring,
async/await, modules, and the longer-tail JSDoc tags — multi-week
ports each.
