# jsgo port TODO

Exhaustive list of remaining work to reach **100 % parity with
closure-compiler v20260415** on the entire extracted
`TypeCheckTest.java` corpus (1444 fixtures) and on any other corpus
we point jsgo at, under both the default and `strictCheckTypes`
diagnostic sets. Nothing is out of scope.

## Porting discipline

Every remaining item below must be delivered as a faithful 1:1 port
of its Java counterpart, not a fresh implementation. Divergences
accumulate into subtle bugs that are painful to unpick at scale;
mirroring Java's structure makes behavioural mismatches surface at
the same spot in the Go code as in the Java source and keeps the
oracle usable as a debugger.

Concretely:

- Preserve class, method, and field names recognisable to a Java
  reader. `TypeCheck.visitCall`, `TypeInference.traverseGetProp`,
  etc.
- Preserve the Java dispatch order and control-flow structure — do
  not flatten a three-method call chain into one function, do not
  reorder visitor branches for aesthetics, do not swap a fallthrough
  for early-return.
- Preserve data-structure choices when they are observable
  (insertion-ordered maps, deterministic iteration, per-node
  property side tables, lattice shapes, visitor dispatch tables).
- When Java itself has a `TODO` or flags something as in-development
  (e.g. `JSC_UNSUPPORTED_BOUNDED_GENERIC_TYPES`), port the
  in-development path faithfully — do not invent our own fix.
- Every non-trivial `internal/` Go file carries a `// Port of
  <java-file> @ closure-compiler v20260415` header (enforced by
  `scripts/ci-guards.sh`). Narrowing qualifier is removed once the
  port matches Java 1:1.
- Skeletons under `internal/jscomp/skel/` are the authoritative
  starting points for every un-ported Java file; read them in
  parallel with the upstream source at `.cache/closure-src/…` when
  porting.

## Current state

964 / 1444 = **66.76 %** on default, 883 / 1444 = **61.15 %** on
`strictCheckTypes`. Tier-0..Tier-3 hand-authored corpora are all at
100 %. Gap to 100 % on default: **480** fixtures across **108**
distinct diagnostic-tuple buckets.

## Triage snapshot (current failing buckets on default)

```
193  want=[JSC_TYPE_MISMATCH]                      got=[]
 37  want=[JSC_UNSUPPORTED_BOUNDED_GENERIC_TYPES]  got=[]   [Java itself marks in-dev]
 17  want=[]                                       got=[JSC_TYPE_MISMATCH]                (our FPs)
 15  want=[JSC_NON_STRINGIFIABLE_OBJECT_KEY]       got=[]
 11  want=[JSC_UNRECOGNIZED_TYPE_ERROR]            got=[]
 11  want=[JSC_INEXISTENT_PROPERTY]                got=[]
  7  want=[JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH] got=[]
  6  want=[JSC_WRONG_ARGUMENT_COUNT]               got=[]
  6  want=[JSC_TYPE_MISMATCH]                      got=[JSC_TYPE_MISMATCH]                (wrong position)
  5  want=[JSC_INVALID_CAST]                       got=[]
  5  want=[JSC_TYPE_PARSE_ERROR]                   got=[]
  5  want=[JSC_UNRECOGNIZED_TYPEOF_ERROR]          got=[]
  5  want=[JSC_NOT_A_CONSTRUCTOR]                  got=[]
  5  want=[JSC_UNDEFINED_VARIABLE]                 got=[]
  5  want=[JSC_TYPE_MISMATCH,JSC_USELESS_CODE]     got=[JSC_USELESS_CODE]
  4  want=[JSC_REFERENCE_BEFORE_DECLARE,JSC_TYPE_MISMATCH] got=[]
  4  want=[JSC_TYPE_MISMATCH,JSC_USED_GLOBAL_THIS] got=[]
  4  want=[JSC_POSSIBLE_INEXISTENT_PROPERTY]       got=[]
  4  want=[JSC_DETERMINISTIC_TEST]                 got=[]
  4  want=[JSC_USED_GLOBAL_THIS]                   got=[]
  4  want=[JSC_INCOMPATIBLE_ALIAS_ANNOTATION]      got=[]
  4  want=[JSC_VAR_MULTIPLY_DECLARED_ERROR]        got=[]
  4  want=[JSC_INCOMPATIBLE_EXTENDED_PROPERTY_TYPE] got=[JSC_INCOMPATIBLE_EXTENDED_PROPERTY_TYPE]  (wrong position)
  3  want=[JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH] got=[]
  3  want=[JSC_INVALID_MODULE_PATH]                got=[]
  3  want=[JSC_PARSE_ERROR]                        got=[JSC_UNDEFINED_VARIABLE]
  3  want=[JSC_PROPERTY_ASSIGNMENT_TO_READONLY_VALUE] got=[]
  3  want=[]                                       got=[JSC_INEXISTENT_PROPERTY]          (our FPs)
  3  want=[JSC_DETERMINISTIC_TEST,JSC_USELESS_CODE] got=[JSC_DETERMINISTIC_TEST,JSC_USELESS_CODE]  (wrong position)
```

Full list: `go run ./tools/fbucket b`. 108 distinct buckets total.

---

## Work list, organized by Java source file

The organising principle is: every item names the Java source it
ports. Anything already partially ported is listed under the Go
file it lives in (all such files carry a `// Port of …` header
naming the Java source); anything un-ported is listed under its
skeleton path in `internal/jscomp/skel/`.

### Parser — `parser.*` (`internal/parser/`)

Current ports are narrow: `scanner.go` is character-level, `parser.go`
covers a narrow JS subset.

Gaps vs. `parser/Scanner.java`:

- Full `TokenType` enum (port `parser/TokenType.java`).
- String-literal tokenizer (escapes, template-literal pieces).
- Number-literal tokenizer (hex, octal, binary, `_`-separators, `n`
  bigint suffix, decimal edge cases).
- Regex-literal tokenizer.
- Comment/trivia preservation through `JsDocTokenStream`.

Gaps vs. `parser/Parser.java`:

- Array and object literals at expression position.
- Destructuring patterns (array and object, nested, default values,
  rest).
- Spread/rest arguments and parameters.
- Template literals (`` `text ${x}` ``), tagged templates.
- Regex literals as expressions.
- `GETELEM` (`x[0]`) as a full expression form.
- Optional chaining (`?.`, `?.(`, `?.[`).
- `async` / `await`.
- Generator functions (`function*`, `yield`, `yield*`).
- Full computed property syntax (`{[k]: v}`, `class { [k]() {} }`).
- Cast expressions `(/** @type {T} */ (e))` as first-class grammar.
- Class-body grammar: class fields, static fields, accessors,
  private `#` members, static blocks.
- `for`/`for-in`/`for-of` with destructured binding.
- `import` / `export` declarations (named, default, namespace, re-export).
- `super.x`, `super(…)` in all positions Java accepts.
- Labeled statements and label-target resolution.

Skeletons under `internal/jscomp/skel/jscomp/parsing/parser/`:
`scanner.go`, `parser.go`, `token.go`, `tokentype.go`,
`identifiertoken.go`, `literaltoken.go`, `templateliteraltoken.go`.

### JSDoc parsing — `jscomp.parsing.*` (`internal/jscomp/parsing/`)

Gaps vs. `JsDocInfoParser.java`:

- `@enum` tag (type + values).
- `@const` with optional type.
- `@dict` / `@struct`.
- `@nosideeffects` + `@noside`.
- `@abstract`.
- `@public` / `@protected` / `@private`.
- `@param` modifiers: `{=T}`, `{...T}`, `{T=}`, default values.
- Callable-expression types: `function(A, B): C`,
  `function(this:T, …)`, `function(new:T, …)`.
- Multi-`@extends` accumulator (for interfaces).
- Multi-`@implements`.
- Position-tracked errors inside type expressions (needed for
  `JSC_UNRECOGNIZED_TYPE_ERROR`, `JSC_TYPE_PARSE_ERROR`).
- `@template` with bounds.
- `@lends`.
- `@typedef` across record / union / templatized forms.
- `@ngInject`, `@idGenerator`, `@wizaction`, `@polymerBehavior`,
  remaining long-tail tags recognised by Java.

Gaps vs. `IRFactory.java`: every construct newly introduced in the
parser needs the matching IR materialisation (port the whole file,
not just what the current corpus forces).

Gaps vs. `TypeTransformationParser.java`: required for
`@template T := … =:` transformations.

Skeletons: `skel/jscomp/parsing/{irfactory,jsdocinfoparser,jsdoctoken,jsdoctokenstream,parserrunner,typetransformationparser}.go`.

### Rhino AST & type expressions — `rhino.*` (`internal/rhino/`)

Gaps vs. `Node.java` (`internal/rhino/node.go`):

- `SideEffectFlags` cache.
- Full ~50-entry `Prop` enum (currently a handful).
- `cloneNode`, `cloneTree`.
- Intrusive linked-list children (Java uses `first`/`last`/`next`).
- `QualifiedName` extraction helper.
- `JSDocInfo` attachment on every node shape that can carry it.
- `getLineno`/`getCharno`/source-file tracking on every node.

Gaps vs. `JSDocInfo.java` (`internal/rhino/jsdoc.go`):

- Multi-`@extends` list (currently single).
- `@implements` list.
- `@template` bounds.
- `@typedef` + visibility bits.
- Full builder pattern (`JSDocInfoBuilder`).

Gaps vs. `Token.java` (`internal/rhino/token.go`): bring to parity
with Java's enum once parser breadth lands.

`JSTypeExpression.java` (`internal/rhino/jstype_expr.go`) is currently
a stub — port the full class.

Other `rhino/` Java files not yet present in `internal/rhino/`
(skeletons exist under `skel/rhino/`):

- `ClosurePrimitive.java`
- `ErrorReporter.java`
- `InputId.java`
- `IR.java`
- `NominalTypeBuilder.java`
- `NonJSDocComment.java`
- `Outcome.java`
- `QualifiedName.java`
- `SimpleSourceFile.java`
- `SourcePosition.java`
- `StaticRef.java`, `StaticScope.java`, `StaticSlot.java`,
  `StaticSourceFile.java`, `StaticSymbolTable.java`
- `TokenStream.java`, `TokenUtil.java`
- `TypeDeclarationsIR.java`

### JSType hierarchy — `rhino.jstype.*` (`internal/rhino/jstype/`)

Currently ported (all as subsets): `JSType`, `FunctionType`,
`ObjectType`, `UnionType`, `EnumType`, `TemplateType`,
`TemplatizedType`, `JSTypeRegistry`.

Gaps in already-ported files:

- `jstype.go` — widen to full `JSType.java` (visitor dispatch,
  `canAssignTo`, `matchesObjectContext`, cache hooks).
- `object_type.go` — `@protected` / `@private` / `@public` access
  bits, `@struct` vs `@dict` semantics, sealed-for-inference flag,
  split `PrototypeObjectType` vs `InstanceObjectType`.
- `function_type.go` — optional / rest / default-value param
  tracking, `THIS` binding type, closure-capture environment,
  `minArity` / `maxArity`, `isOrdinaryFunction` / `isConstructor` /
  `isInterface`.
- `templatized_type.go` — inheritance through templatized
  supertypes, recursive substitution at nested generic positions.
- `union_type.go` — union-over-union flattening, deterministic
  diagnostic ordering, `UnionTypeBuilder`.
- `enum_type.go` — full `EnumType` + separate `EnumElementType`.
- `template_type.go` — bounded `@template` constraints.
- `registry.go` — full `JSTypeRegistry.java`: native-type table,
  late-resolution worklist, property-type side table,
  scope-registered forward declarations.

Java JSType files not yet present in `internal/rhino/jstype/`
(skeletons exist under `skel/rhino/jstype/`):

- `AllType.java`
- `ArrowType.java`
- `BigIntType.java`
- `BooleanLiteralSet.java`
- `BooleanType.java`
- `EnumElementType.java` (separate from `EnumType`)
- `InstanceObjectType.java`
- `KnownSymbolType.java`
- `NamedType.java` (forward-ref nominal)
- `NoObjectType.java`, `NoResolvedType.java`, `NoType.java`
- `NullType.java`
- `NumberType.java`
- `Property.java`, `PropertyMap.java`
- `PrototypeObjectType.java`
- `ProxyObjectType.java`
- `RecordType.java` (distinct from anon `ObjectType`),
  `RecordTypeBuilder.java`
- `StringType.java`
- `SymbolType.java`
- `UnknownType.java`
- `ValueType.java`
- `VoidType.java`
- `FunctionParamBuilder.java`
- `TemplateTypeMap.java`, `TemplateTypeReplacer.java`
- `SubtypeChecker.java`, `EqualityChecker.java`,
  `RelationshipVisitor.java`
- `CanCastToVisitor.java`, `ContainsUpperBoundSuperTypeVisitor.java`,
  `JSTypeIterations.java`, `Visitor.java`
- `TypeStringBuilder.java`
- `JSTypeResolver.java`
- `JSTypeClass.java`, `JSTypeNative.java`
- `EquivalenceMethod.java`

### Typed-scope creation — `TypedScopeCreator.java` (`internal/jscomp/scope.go`)

The current 2186-line port is a subset that collapses Java's three
inner-builder classes into free functions. Gaps:

- `AbstractScopeBuilder` / `ClassScopeBuilder` /
  `FunctionScopeBuilder` inner-class separation, with the same
  two-pass `createScopeIn` interaction that Java uses.
- Destructuring-pattern declarations (var, let, const, parameter).
- Per-block `let` / `const` hoisting rules.
- `@enum` slot type + `EnumElementType` registration.
- `Object.defineProperty` / `Object.defineProperties` property
  registration.
- ES6 module import bindings.
- `goog.require` / `goog.provide` / `goog.module` interaction
  (coordinates with `CheckClosureImports`).
- `recordDecl` vs `recordFunction` vs `recordClass` vs
  `recordTypedef` vs `recordPrototypeAssign` dispatch — preserve
  Java's exact ordering.
- Separate `AbstractScope` / `Scope` / `TypedScope` / `TypedVar` /
  `Var` / `AbstractVar` split (skeletons exist).

### Type inference — `TypeInference.java` + `SemanticReverseAbstractInterpreter.java`

Files: `internal/jscomp/{inference.go, typeinferencepass.go, refinement.go, flowscope.go}`.

Gaps vs. `TypeInference.java`:

- Binary-op inference: `ADD`, `SUB`, `MUL`, `DIV`, `MOD`, bitwise,
  comparison — currently narrow, widen to Java's full dispatch
  including string-vs-number `ADD` fallout.
- `GETELEM` inference.
- Compound-assign result-type sharpening.
- `GETPROP` refinement (`if (x.foo) { … x.foo … }`).
- Assignment-with-union-LHS.
- Falsy-on-union-with-literal-false.
- Nested AND/OR sub-expression refinement.
- Loop-aware dataflow: `while (!x) { x = u(); }` must refine `x`
  after the loop body.
- Missing literal cases: BigInt, Symbol, regex.
- `FlowScope` port as a proper lattice (currently a narrow map).

Gaps vs. `SemanticReverseAbstractInterpreter.java`: full walk of
all condition shapes (`in`, `delete`, `instanceof`, `typeof`,
equality with `==` vs `===`, literal narrowing, tuple-destructure
discrimination).

Skeletons under `skel/jscomp/typeinfer/`:
`chainablereverseabstractinterpreter.go`,
`closurereverseabstractinterpreter.go`, `flowscope.go`,
`reverseabstractinterpreter.go`,
`semanticreverseabstractinterpreter.go`.

### Type checking — `TypeCheck.java` (`internal/jscomp/typecheck.go`)

Skeleton target: `skel/jscomp/typecheck.go` (the un-narrowed port).
Current file is 2121 lines covering a subset; remaining gaps:

- `visitArrayLit` / `visitObjectLit` (full inference + contextual
  typing + key checks).
- `visitGetElem`.
- `visitFunction` body check (current: signature only).
- `visitRegExpLit`, `visitTemplateLit`.
- `visitCast` / ASSIGN-with-cast mismatch.
- `checkPropertyAccess` non-structural broadening (full
  `JSC_INEXISTENT_PROPERTY` + `JSC_POSSIBLE_INEXISTENT_PROPERTY`).
- `validateNonStringifiableKey` → `JSC_NON_STRINGIFIABLE_OBJECT_KEY`.
- Call/NEW arity and "not a constructor": `JSC_WRONG_ARGUMENT_COUNT`,
  `JSC_NOT_A_CONSTRUCTOR`, `JSC_CONSTRUCTOR_NOT_CALLABLE`,
  `JSC_INSTANTIATE_ABSTRACT_CLASS`.
- `JSC_INEXISTENT_PARAM`, `JSC_INVALID_PARAM`.
- `JSC_PROPERTY_ASSIGNMENT_TO_READONLY_VALUE`,
  `JSC_CONFLICTING_GETTER_SETTER_TYPE`.
- `JSC_ILLEGAL_PROPERTY_CREATION`,
  `JSC_ILLEGAL_PROPERTY_CREATION_ON_UNION_TYPE`.
- `JSC_INVALID_CAST`, `JSC_ILLEGAL_IMPLICIT_CAST`.
- `JSC_EXPECTED_THIS_TYPE`, `JSC_NOT_FUNCTION_TYPE`.
- `JSC_CYCLIC_INHERITANCE_ERROR`, `JSC_EXTENDS_NON_OBJECT`,
  `JSC_EXTENDS_WITHOUT_TYPEDEF`, `JSC_IMPLEMENTS_NON_INTERFACE`.
- `JSC_INTERFACE_METHOD_NOT_EMPTY`,
  `JSC_INTERFACE_METHOD_NOT_IMPLEMENTED`.
- `JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH`,
  `JSC_INCOMPATIBLE_EXTENDED_PROPERTY_TYPE` (class-style).
- `JSC_TEMPLATE_TRANSFORMATION_ON_CLASS`,
  `JSC_UNSUPPORTED_BOUNDED_GENERIC_TYPES` (port Java's in-dev
  branch faithfully — including the TODO it carries).
- `JSC_UNKNOWN_LENDS`, `JSC_LENDS_ON_NON_OBJECT`,
  `JSC_INCOMPATIBLE_ALIAS_ANNOTATION`.
- `JSC_ENUM_INITIALIZER_NOT_ENUM`, `JSC_ILLEGAL_CLASS_KEY`,
  `JSC_ILLEGAL_OBJLIT_KEY`, `JSC_BAD_TYPE_FOR_BIT_OPERATION`.
- `JSC_MULTIPLE_VAR_DEF`, `JSC_RESOLVED_TAG_EMPTY`,
  `JSC_DICT_EXTEND_STRUCT_TYPE`, `JSC_UNKNOWN_OVERRIDE`,
  `JSC_ILLEGAL_PROPERTY_ACCESS`.
- `JSC_UNRECOGNIZED_TYPEOF_ERROR`, `JSC_UNKNOWN_TYPEOF_VALUE`.

### Type validation + mismatch tracking — `TypeValidator.java`, `TypeMismatch.java`

Skeletons: `skel/jscomp/typevalidator.go`, `skel/jscomp/typemismatch.go`.
No Go port yet; currently `typecheck.go` does validation inline.
Port as separate files to match Java layout — unblocks consistent
diagnostic formatting + handles the JSC_TYPE_MISMATCH false
positives driven by subtype-check disagreements.

### Control-flow analysis — `ControlFlowAnalysis.java` (`internal/jscomp/controlflowanalysis.go`)

Gaps:

- Labeled `break` / `continue` to **nested** loops.
- Exception-edge fan-out: each call / throw to nearest catch.
- `with` statement modelling.
- Generator suspend edges (`yield`).
- `async` / `await` suspend edges.
- Assignment-in-condition edge handling.
- `finally` re-raise edges.

### Dataflow analysis — `DataFlowAnalysis.java` (`internal/jscomp/dataflow.go`)

Gaps:

- Backward direction.
- Lattice-annotation side table (currently stored directly on
  `CfgNode`).
- Fixed-point iteration limit heuristic.
- `JoinOp` abstraction as a separate interface.

### Compiler driver — `Compiler.java` (`internal/jscomp/compiler.go`)

Currently a ~65-line shim. Gaps (skeletons under `skel/jscomp/`):

- `AbstractCompiler.java`.
- `Compiler.java` full: `PhaseOptimizer`, pass registry, `Result`,
  source-map context, module graph.
- `CompilerInput.java`, `CompilerInputProvider.java`.
- `CompilerOptions.java` beyond current flag subset, plus
  `CompilerOptionsPreprocessor.java`.
- `DefaultPassConfig.java` with the real ordered pass list.
- `PassConfig.java`, `PassFactory.java`, `PassListBuilder.java`,
  `PassNames.java`, `CompilerPass.java`.
- `NodeTraversal.java` — widen past current subset (full
  pre-order / post-order / scope-callback variants).
- `JSError.java`, `DiagnosticType.java`, `DiagnosticGroup.java`,
  `DiagnosticGroups.java`, `DiagnosticGroupWarningsGuard.java`,
  `DiagnosticGroupPathSuppressingWarningsGuard.java`,
  `DiagnosticToSuppressionMapper.java`,
  `AllowlistWarningsGuard.java`, `CheckLevel.java`,
  `AbstractMessageFormatter.java`.

### Java passes not yet started (each maps to a skeleton)

Skeletons live under `internal/jscomp/skel/jscomp/`; porting = move
the file out of `skel/`, drop the `//go:build skeletons` tag, fill
in bodies, wire into `compiler.go`, add the pass's `JSC_*` codes to
`internal/diag/diag.go`.

- `CheckAccessControls.java` → `checkaccesscontrols.go`
  — unlocks `JSC_ILLEGAL_PROPERTY_ACCESS`.
- `CheckClosureImports.java` → `checkclosureimports.go`
  — unlocks `goog.require` / `goog.provide` / `JSC_INVALID_MODULE_PATH`.
- `CheckConformance.java` → `checkconformance.go`.
- `CheckDebuggerStatement.java` → `checkdebuggerstatement.go`.
- `CheckGlobalThis.java` → `checkglobalthis.go`
  — unlocks `JSC_USED_GLOBAL_THIS`.
- `CheckJSDoc.java` → `checkjsdoc.go`
  — unlocks bad-`@typedef` / bad-`@enum` diagnostics,
  `JSC_INVALID_INTERFACE_MEMBER_DECLARATION`.
- `CheckMissingRequires.java` → `checkmissingrequires.go`.
- `CheckMissingReturn.java` → `checkmissingreturn.go`
  — unlocks the missing-path `JSC_MISSING_RETURN_STATEMENT`.
- `CheckRegExp.java` → `checkregexp.go`.
- `CheckSuper.java` → `checksuper.go`.
- `CheckSuspiciousCode.java` → `checksuspiciouscode.go`
  — unlocks `JSC_SUSPICIOUS_IN`, `JSC_SUSPICIOUS_SEMICOLON`,
  `JSC_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR`.
- `CheckTypeImportCodeReferences.java` → `checktypeimportcodereferences.go`.
- `CheckUnreachableCode.java` → `checkunreachablecode.go`.
- `InferConsts.java` → `inferconsts.go`.
- `InferJsDocInfo.java` → `inferjsdocinfo.go`.
- `NodeIterators.java`, `NodeNameExtractor.java`, `NodeUtil.java`
  — utility ports used by most of the above.
- `AstAnalyzer.java`, `AstFactory.java`, `AstManipulations.java`,
  `AstValidator.java`.
- `BasicBlock.java`.
- `SymbolTable.java` (+ `ScopedName.java`).
- `TypeMatchingStrategy.java`, `TypeMismatch.java`,
  `TypeTransformation.java`, `TypeValidator.java`.
- There is no skeleton for `CheckDeterministicTest.java` in the
  current generation pass — add it to `scripts/gen-skeletons.sh`
  and port to unlock `JSC_DETERMINISTIC_TEST` cleanly.

### Currently-ported checks that need widening

- `internal/jscomp/checkdupvardecl.go` — port of part of
  `VariableReferenceCheck.java`. Add the type-mismatch arm
  (`JSC_DUP_VAR_DECLARATION_TYPE_MISMATCH` full semantics), the
  `JSC_FUNCTION_MASKS_VARIABLE` arm, the full
  `JSC_VAR_MULTIPLY_DECLARED_ERROR` path.
- `internal/jscomp/checksideeffects.go` — `@nosideeffects`-aware
  CALL classification, sub-expression `a[b].c;` reporting.
- `internal/jscomp/superclasscheck.go` — constructor-arg
  compatibility, multi-`@extends` interface diamond check,
  `validateImplements` loop
  (`JSC_INTERFACE_METHOD_NOT_IMPLEMENTED`).
- `internal/jscomp/varcheck.go` — per-scope unreached-NAME check,
  `JSC_REFERENCE_BEFORE_DECLARE` for var hoisting.
- `internal/jscomp/varrefcheck.go` — full port of
  `VariableReferenceCheck.java`.
- `internal/jscomp/checkjsdoctypes.go` — extend when `CheckJSDoc`
  lands.
- `internal/jscomp/checkduplicatemember.go` — widen to class-body
  member checks.
- `internal/jscomp/controlflowgraph.go` — widen alongside
  `controlflowanalysis.go`.
- `internal/jscomp/flowscope.go` — replace with the full port of
  `jscomp/type/FlowScope.java` (skeleton at
  `skel/jscomp/typeinfer/flowscope.go`).
- `internal/jscomp/nodetraversal.go` — widen to full `NodeTraversal`.

### Externs

- `internal/jscomp/default_externs.go` — currently an embedded
  subset. Port Java's full `DefaultExterns.js` bundle (`Window.*`,
  `HTMLElement.*`, typed-array methods, `String.prototype.*`,
  `Array.prototype.*`, `Object.prototype.*`, `Map`, `Set`,
  `Promise`, `Symbol`, `Iterator`, `goog.*`). Unlocks the remaining
  `JSC_UNDEFINED_VARIABLE` misses and reduces TYPE_MISMATCH false
  positives.

---

## Remaining diagnostic-tuple gaps (ledger)

Every failing fixture shows up in one of these three categories.
Listed so nothing is missed — each maps back to one of the porting
areas above. Re-run `go run ./tools/fbucket b` after every landing.

### A. Pure misses (want=[X] got=[])

| Count | Diagnostic | Fix-site |
|-------|-----------|----------|
| 193 | `JSC_TYPE_MISMATCH` | `inference.go`, `typecheck.go`, `typevalidator.go` (new port) |
| 37  | `JSC_UNSUPPORTED_BOUNDED_GENERIC_TYPES` | `typecheck.go`'s `@template` branch — port Java's in-dev path |
| 15  | `JSC_NON_STRINGIFIABLE_OBJECT_KEY` | `typecheck.go` `visitGetElem` / computed-prop |
| 11  | `JSC_UNRECOGNIZED_TYPE_ERROR` | `jsdocinfoparser.go` position-tracked, `scope.go` resolve-time |
| 11  | `JSC_INEXISTENT_PROPERTY` | `typecheck.go` non-structural branch |
|  7  | `JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH` (class-style) | `superclasscheck.go` |
|  6  | `JSC_WRONG_ARGUMENT_COUNT` | `typecheck.go` call arity |
|  5  | `JSC_INVALID_CAST` | `typecheck.go` visitCast |
|  5  | `JSC_TYPE_PARSE_ERROR` | `jsdocinfoparser.go` |
|  5  | `JSC_UNRECOGNIZED_TYPEOF_ERROR` | `checkjsdoc.go` / `typecheck.go` |
|  5  | `JSC_NOT_A_CONSTRUCTOR` | `typecheck.go` visitNew |
|  5  | `JSC_UNDEFINED_VARIABLE` | externs + `varcheck.go` |
|  4  | `JSC_DETERMINISTIC_TEST` | port `CheckDeterministicTest` |
|  4  | `JSC_USED_GLOBAL_THIS` | port `CheckGlobalThis` |
|  4  | `JSC_INCOMPATIBLE_ALIAS_ANNOTATION` | `typecheck.go` alias branch |
|  4  | `JSC_VAR_MULTIPLY_DECLARED_ERROR` | `varrefcheck.go` |
|  4  | `JSC_POSSIBLE_INEXISTENT_PROPERTY` | `typecheck.go` + diagnostic level |
|  4  | `JSC_REFERENCE_BEFORE_DECLARE` | `varrefcheck.go` |
|  3  | `JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH` | `typecheck.go` + `superclasscheck.go` |
|  3  | `JSC_INVALID_MODULE_PATH` | port `CheckClosureImports` |
|  3  | `JSC_PROPERTY_ASSIGNMENT_TO_READONLY_VALUE` | `typecheck.go` |
|  2  | `JSC_EXTENDS_NON_OBJECT` | `typecheck.go` |
|  2  | `JSC_ILLEGAL_CLASS_KEY` | `typecheck.go` |
|  2  | `JSC_SUSPICIOUS_SEMICOLON` | port `CheckSuspiciousCode` |
|  2  | `JSC_TEMPLATE_TRANSFORMATION_ON_CLASS` | `typecheck.go` |
|  2  | `JSC_ILLEGAL_OBJLIT_KEY` | `typecheck.go` |
|  2  | `JSC_FUNCTION_MASKS_VARIABLE` | `varrefcheck.go` |
|  2  | `JSC_EXPECTED_THIS_TYPE` | `typecheck.go` |
|  2  | `JSC_INTERFACE_METHOD_NOT_IMPLEMENTED` | `superclasscheck.go` |
|  2  | `JSC_IMPLEMENTS_NON_INTERFACE` | `typecheck.go` |
|  2  | `JSC_DUP_VAR_DECLARATION_TYPE_MISMATCH` | `checkdupvardecl.go` widening |
|  1 each | `JSC_UNKNOWN_LENDS`, `JSC_CONSTRUCTOR_NOT_CALLABLE`, `JSC_UNKNOWN_OVERRIDE`, `JSC_INVALID_PARAM`, `JSC_INSTANTIATE_ABSTRACT_CLASS`, `JSC_CONFLICTING_GETTER_SETTER_TYPE`, `JSC_INCOMPATIBLE_EXTENDED_PROPERTY_TYPE`, `JSC_BAD_TYPE_FOR_BIT_OPERATION`, `JSC_MULTIPLE_VAR_DEF`, `JSC_INEXISTENT_PARAM`, `JSC_CYCLIC_INHERITANCE_ERROR`, `JSC_ENUM_INITIALIZER_NOT_ENUM`, `JSC_NOT_FUNCTION_TYPE`, `JSC_INTERFACE_METHOD_NOT_EMPTY`, `JSC_UNKNOWN_TYPEOF_VALUE`, `JSC_RESOLVED_TAG_EMPTY`, `JSC_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR`, `JSC_DICT_EXTEND_STRUCT_TYPE`, `JSC_SUSPICIOUS_IN`, `JSC_ILLEGAL_PROPERTY_CREATION_ON_UNION_TYPE`, `JSC_USELESS_CODE` (edge case), `JSC_ILLEGAL_PROPERTY_ACCESS`, `JSC_EXTENDS_WITHOUT_TYPEDEF`, `JSC_LENDS_ON_NON_OBJECT` | see the per-pass sections above |

### B. False positives (want=[] got=[X])

| Count | Spurious diagnostic | Fix-site |
|-------|---------------------|----------|
| 17 | `JSC_TYPE_MISMATCH` | narrow `inference.go` / loop-aware refinement / `typevalidator.go` port |
|  3 | `JSC_INEXISTENT_PROPERTY` | broaden "unknown-structure" receiver recognition |
|  1 | `JSC_DETERMINISTIC_TEST` | our heuristic — port `CheckDeterministicTest` precisely |
|  1 | `JSC_MISSING_RETURN_STATEMENT` | CFG defect; revisit once `CheckMissingReturn` lands |
|  1 | `JSC_DUP_VAR_DECLARATION_TYPE_MISMATCH` | narrow duplicate-decl detection |

### C. Position mismatches (same code, wrong line/col)

Same diagnostic reported at the wrong location — usually because
we emit on a parent node while Java emits on a child (or vice
versa). Fix each by consulting Java's exact report-site.

- 6× `JSC_TYPE_MISMATCH`
- 4× `JSC_INCOMPATIBLE_EXTENDED_PROPERTY_TYPE`
- 3× `JSC_DETERMINISTIC_TEST,JSC_USELESS_CODE`
- ~15 further single-count combinations (see `fbucket b`)

### D. Mixed tuples (partial match)

Some fixtures match one of an n-tuple but miss another. These drop
out automatically as the underlying pure-miss / FP entries above
are resolved; tracked in `fbucket` output for regression watch.

---

## Supporting work

- Populate the hold-out corpus slot (SPEC guardrail #5) and keep
  its pass-rate locked to the visible corpus.
- Performance: stay within `[0.1×, 2×]` of Java on Tier-1, verified
  at each major port landing.
- Oracle regeneration: re-run `scripts/gen-oracles.sh` whenever
  fixture layout or Java pinned version changes.
- Triage tooling: `tools/fbucket` keeps the bucket table; re-run
  after every significant change, paste the snapshot back into
  this file.
- CI floor in `harness/harness_test.go` raised alongside progress
  so regressions show up immediately.
- Skeleton refresh: re-run `scripts/gen-skeletons.sh` when the Java
  tag is bumped, to keep `internal/jscomp/skel/` in sync.
