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
- **No Go-only code paths.** If a helper, struct, or check doesn't
  exist in Java, it doesn't belong here. When a divergence looks
  unavoidable, first map the Java call chain to the Go equivalent
  before adding anything new.

---

## 1. Current state

**1101 / 1444 = 76.25 %** on default, **1017 / 1444 = 70.43 %** on
`strictCheckTypes` (baseline commit `fc26e86` at the time of the
last session refresh; TODO was re-audited at `9816618` so §2
snapshot is captured there). Tier-0..Tier-3 hand-authored corpora
are all at 100 %. Gap to 100 % on default: **343** fixtures across
the buckets below; strict lags by **84** more.

Diff vs. the previous TODO snapshot: default +137 (964 → 1101),
strict +134 (883 → 1017).

## 2. Triage snapshot (current failing buckets on default)

Re-derived from `go run ./tools/fbucket b` at `9816618`. Rows with
three or more fixtures are shown verbatim; single-count tuples are
summarised at the bottom.

```
 179 want=[JSC_TYPE_MISMATCH]                          got=[]
  24 want=[]                                           got=[JSC_TYPE_MISMATCH]                    (our FPs)
  11 want=[JSC_NON_STRINGIFIABLE_OBJECT_KEY]           got=[]
  11 want=[JSC_TYPE_MISMATCH]                          got=[JSC_TYPE_MISMATCH]                    (wrong position)
   8 want=[JSC_INEXISTENT_PROPERTY]                    got=[]
   7 want=[JSC_UNRECOGNIZED_TYPE_ERROR]                got=[]
   6 want=[JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH]    got=[]
   5 want=[JSC_INVALID_CAST]                           got=[]
   4 want=[]                                           got=[JSC_INEXISTENT_PROPERTY]              (our FPs)
   4 want=[JSC_TYPE_MISMATCH,JSC_USELESS_CODE]         got=[JSC_USELESS_CODE]
   4 want=[]                                           got=[JSC_POSSIBLE_INEXISTENT_PROPERTY]     (our FPs)
   4 want=[JSC_DETERMINISTIC_TEST]                     got=[]
   3 want=[JSC_INVALID_MODULE_PATH]                    got=[]
   3 want=[JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH]     got=[]
   3 want=[JSC_TYPE_MISMATCH,JSC_USED_GLOBAL_THIS]     got=[JSC_USED_GLOBAL_THIS]
   3 want=[JSC_DETERMINISTIC_TEST,JSC_USELESS_CODE]    got=[JSC_DETERMINISTIC_TEST,JSC_USELESS_CODE] (wrong position)
   3 want=[JSC_REFERENCE_BEFORE_DECLARE,JSC_TYPE_MISMATCH] got=[]
   3 want=[JSC_PROPERTY_ASSIGNMENT_TO_READONLY_VALUE]  got=[]
   3 want=[JSC_UNDEFINED_VARIABLE]                     got=[JSC_TYPE_MISMATCH,JSC_UNDEFINED_VARIABLE] (over-reporting)
   2 want=[JSC_POSSIBLE_INEXISTENT_PROPERTY,JSC_USELESS_CODE] got=[JSC_POSSIBLE_INEXISTENT_PROPERTY]
   2 want=[JSC_PARSE_ERROR]                            got=[JSC_POSSIBLE_INEXISTENT_PROPERTY,JSC_UNDEFINED_VARIABLE]
   2 want=[JSC_TYPE_MISMATCH]                          got=[JSC_INEXISTENT_PROPERTY]              (wrong code)
   2 want=[JSC_TYPE_MISMATCH]                          got=[JSC_POSSIBLE_INEXISTENT_PROPERTY]     (wrong code)
   2 want=[JSC_UNDEFINED_VARIABLE]                     got=[JSC_POSSIBLE_INEXISTENT_PROPERTY,JSC_UNDEFINED_VARIABLE]
   2 want=[JSC_INEXISTENT_PROPERTY]                    got=[JSC_INEXISTENT_PROPERTY,JSC_POSSIBLE_INEXISTENT_PROPERTY]
   2 want=[JSC_UNDEFINED_VARIABLE]                     got=[JSC_UNDEFINED_VARIABLE,JSC_USELESS_CODE]
   2 want=[JSC_UNDEFINED_VARIABLE]                     got=[]
   2 want=[JSC_TEMPLATE_TRANSFORMATION_ON_CLASS,JSC_UNRECOGNIZED_TYPE_ERROR] got=[]
   2 want=[JSC_EXPECTED_THIS_TYPE]                     got=[]
   2 want=[JSC_MISSING_RETURN_STATEMENT,JSC_TYPE_MISMATCH] got=[]
```

Single-count residuals: `JSC_CYCLIC_INHERITANCE_ERROR`,
`JSC_INCOMPATIBLE_EXTENDED_PROPERTY_TYPE`,
`JSC_ILLEGAL_PROPERTY_CREATION_ON_UNION_TYPE`,
`JSC_VAR_MULTIPLY_DECLARED_ERROR`, `JSC_UNKNOWN_TYPEOF_VALUE`,
`JSC_NOT_FUNCTION_TYPE`, `JSC_INEXISTENT_PARAM`,
`JSC_UNKNOWN_OVERRIDE`, `JSC_RESOLVED_TAG_EMPTY`,
`JSC_EXTENDS_WITHOUT_TYPEDEF`, `JSC_ILLEGAL_PROPERTY_CREATION`,
`JSC_ILLEGAL_IMPLICIT_CAST`, `JSC_INSTANTIATE_ABSTRACT_CLASS`,
`JSC_UNKNOWN_LENDS`, `JSC_LENDS_ON_NON_OBJECT`,
`JSC_ENUM_INITIALIZER_NOT_ENUM`, `JSC_INVALID_INTERFACE_MEMBER_DECLARATION`,
`JSC_INVALID_PARAM`, `JSC_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR`,
`JSC_SUSPICIOUS_IN`, plus ~15 tuple/position combinations. Full list:
`go run ./tools/fbucket b`.

Net triage: the 179-count `want=[JSC_TYPE_MISMATCH]/got=[]` bucket
plus the 24-count FP bucket still dominate. Every remaining 3+
bucket maps directly to a file in section 3; single-count tuples
map to section 5.

---

## 3. Per-file port status

Every file under `internal/` with a `// Port of …` header, the
Java source it mirrors, one-line status, and a short list of Java
entry points / fields / branches missing on the Go side. Method
names are Java-side unless noted.

### 3.1 `internal/rhino/`

#### `internal/rhino/node.go` — `rhino/Node.java` (3944 LoC)
**Status: subset.** 389 LoC, ~41 exported methods vs. ~464 Java
public/protected members.
- Missing: `SideEffectFlags` (cache + getter/setter), `cloneNode`
  distinct from `CloneTree`, intrusive `first`/`last`/`next`
  linked-list children (Go uses a slice), ~40 `Prop` enum values
  (only the frequently-touched handful are defined), `extractName`,
  `isValidQualifiedName`, `getSecondChild`/`getLastChild` alongside
  existing helpers, `srcref`/`srcrefIfMissing`, `getLineno`
  vs. `Lineno` convention mismatch on a few callers, token-level
  classifiers (`isOptionalChainCall`, `isPureDelete`,
  `matchesName`/`matchesQualifiedName`), `JSType` storage +
  `getJSType`/`setJSType`, `getGrandparent`, full
  `iterChildrenInReverse` helpers.
- Go-only: the slice-based child model and `Children()` iterator
  — not in Java. Must be replaced by first/next linked-list fields
  before `Node` is considered 1:1.

#### `internal/rhino/jsdoc.go` + `jsdoc_info_builder.go` — `rhino/JSDocInfo.java` (2922 LoC)
**Status: subset.** ~404 LoC each. Covers the bit flags the current
corpus exercises.
- Missing: `@throws` / `@preserveTry` / `@closurePrimitive`
  accessors; full visibility lattice
  (`getOriginalVisibility`/`setOriginalVisibility`); bit-precise
  `BitField` enum; `getMarkers` (source-position-tagged marker
  list needed for `JSC_UNRECOGNIZED_TYPE_ERROR` positions);
  `cloneAndReset` + `cloneOn`; `getSourceComment`; `@modifies`
  slots; `@ngInject`, `@idGenerator`, `@wizaction`,
  `@polymerBehavior`, `@customElement` (parsed + surfaced); multi-
  `@extends` is present, multi-`@implements` list is present.
  `JSDocInfoBuilder` omits `recordThrowType`, `recordSuppressions`
  merge semantics, `recordTypeTransformations`.

#### `jstype_expr.go` + `qualified_name.go` + `simple_source_file.go` + `input_id.go`
**Status: near-complete.** `jstype_expr.go` missing the unified
`evaluate(StaticTypedScope, JSTypeRegistry)` entry (callers
resolve through the scope creator directly). Verify
`QualifiedName.getOwner`/`matchesRoot`/`isSimple` against Java —
a few callers reproduce logic inline.

#### `static_*.go` (`StaticRef`, `StaticScope`, `StaticSlot`, `StaticSourceFile`, `StaticSymbolTable`)
**Status: stub interfaces** (17–67 LoC each). Surface only so the
`Static*Typed*` subinterfaces compile. Matches Java layout.

#### `internal/rhino/token.go` — `rhino/Token.java`
**Status: subset of the enum.** 380 LoC, covers ~95 of the ~140
Java tokens.
- Missing: TS-only tokens (`INTERFACE`, `IMPLEMENTS_MEMBERS`,
  `DECLARE`, `TYPE_ALIAS`); some AST-only tokens
  (`PLACEHOLDER1`/`PLACEHOLDER2`, `ROOT`, `INPUT_ID`,
  `IMPORT_STAR`, `MODULE_BODY`, `DYNAMIC_IMPORT`, `IMPORT_META`,
  `ES6_SUPER`). Add when parser breadth lands.

### 3.2 `internal/rhino/jstype/`

Currently ported (all as subsets unless noted):

#### `jstype.go` — `rhino/jstype/JSType.java` (1705 LoC)
**Status: thin dispatch shim.** 124 LoC, 7 Go-side methods. Java
exposes 100+ public methods. Most `JSType` predicates (`isFoo`,
`toMaybeFooType`, `visit`, `autobox`, `autoboxesTo`,
`restrictByTypeof`, `getRestrictedTypeGivenOutcome`,
`getTypesUnderEquality`) are implemented inline on the concrete
types or not yet present.
- Missing: the whole `Visitor<T>` double-dispatch path apart from
  `can_cast_to_visitor`/`contains_upper_bound_super_type_visitor`.
- Missing: `meet`, `isStrictSubtypeOf`, `matchesObjectContext`,
  `getOwnerFunction`, `collapseUnion`,
  `getTypeOfThis`.

#### `object_type.go` + `prototype_object_type.go` + `instance_object_type.go` + `object_type_ext.go` — `ObjectType.java` + `PrototypeObjectType.java` + `InstanceObjectType.java` (929+545+~200 LoC)
**Status: subset.** Split across four Go files tracking the Java
hierarchy.
- Missing: `@private`/`@protected`/`@public` access bits as part
  of `Property.Info` (visibility lattice not yet consulted in
  `CheckAccessControls`); deprecation markers; `@struct` vs
  `@dict` enforcement on definition site; sealed-for-inference
  flag; `hasOwnDeclaredProperty`, `defineDeclaredProperty`,
  `setPropertyJSDocInfo` preserving source positions; `isRecordType`
  vs anonymous-object distinction.

#### `function_type.go` — `FunctionType.java` (1649 LoC)
**Status: subset** (432 LoC, 29 methods).
- Missing: the full `Builder` chain (`withReturnType`,
  `withTemplateKeys`, `withClosurePrimitiveId`, `withKind`);
  `setImplicitMatch`; `setImplementedInterfaces` merged with
  `@implements`; `getOwnerFunction`; `isOrdinaryFunction`
  distinction from `isArrow`; `getInstanceType`/`setInstanceType`
  across constructor vs. interface; arity fields (`minArity`,
  `maxArity`) partially present but not wired through every call
  site; TypeScript-style `readonly` modifier on params.

#### `union_type.go` — `UnionType.java` (939 LoC)
**Status: subset** (94 LoC, 4 methods) — delegates most work to
`registry.go` when flattening.
- Missing: the full `UnionTypeBuilder` (add/union/flatten state
  machine); deterministic alternate ordering by token-stream
  index; union-of-union collapse with cycle detection;
  `restrictByNotNullOrUndefined`; `getGreatestSubtype`
  dispatched per-alternate; `meetWith`.

#### `enum_type.go` + `enum_element_type.go` — `EnumType.java` + `EnumElementType.java`
**Status: subset.** `enum_type.go` 129 LoC, `enum_element_type.go`
101 LoC.
- Missing: `@enum` attribute inheritance across typedef-chained
  aliases; `JSC_ENUM_INITIALIZER_NOT_ENUM` is driven by
  `scope.go` not here.

#### `templatized_type.go` — `TemplatizedType.java` (239 LoC)
**Status: subset** (61 LoC, 3 methods). Core substitution works
through `template_type_replacer.go`.
- Missing: `getTemplateArgs` keyed by template key, not position;
  inheritance through templatized supertypes when the supertype
  itself is templatized again; `isInstantiatedGenericType`.

#### `template_type.go` + `template_type_map.go` + `template_type_replacer.go` — matching Java files
**Status: subset.** `template_type_map.go` is 454 LoC, fairly
complete.
- Missing: bounded template constraint storage
  (`JSC_UNSUPPORTED_BOUNDED_GENERIC_TYPES` depends on this —
  currently the bound is parsed but not enforced through the
  replacer); `TemplateTypeReplacer.KeyTracker` visitation for
  cycle detection; `TemplateTypeMap.getDecoupledTemplateKeys`.

#### `named_type.go` — `NamedType.java` (565 LoC)
**Status: subset** (272 LoC, 22 methods).
- Missing: forward-declared type resolution through the module
  graph; `resolveViaClosureNamespace` / `resolveViaRegistry`
  split that backs `JSC_UNRECOGNIZED_TYPE_ERROR` reporting.

#### `registry.go` — `rhino/jstype/JSTypeRegistry.java` (2545 LoC)
**Status: stub.** 196 LoC, 11 methods — barely scratches Java.
**Biggest lever in the whole port.**
- Missing: native-type table and eager construction
  (`createNativeTypes`, `getNativeType`); late-resolution
  worklist (`resolveTypes`, `registerTypeResolver`); property-
  type side table (`getPropertyTypes`, `declarePropertyOnInstance`,
  `findTypeForProperty`); scope-registered forward declarations
  (`forwardDeclareType`); closure-namespace registry
  (`registerClosureNamespace`, `getClosureNamespace`); template-
  key registration; `ScopedName` keyed type table.
- Go-only: an ad-hoc `namespacePropNode` map on `TypedScope`
  substitutes for Java's property-type side table.

#### `subtype_checker.go` + `equality_checker.go` + `can_cast_to_visitor.go` + `contains_upper_bound_super_type_visitor.go` + `visitor.go` + `arrow_type.go` + `function_param_builder.go`
**Status: near-complete.** Structure matches Java 1:1 for these
files (449 + 351 + 261 + 182 + 165 + 152 + 134 LoC). Verify
`SubtypeChecker.isStructuralSubtypeHelper` ordering of property
checks vs. cycle detection.

#### `property.go` + `property_map.go` — `Property.java` + `PropertyMap.java`
**Status: near-complete structure**, subset of semantics (293 +
367 LoC). Missing `isSystemProp`, `getOwnerInstantiation`,
snapshot-in-time cloning used for interface-conflict diagnostics.

#### `proxy_object_type.go` — `ProxyObjectType.java`
**Status: subset** (454 LoC). Base class for `NamedType` and
`TemplatizedType`; verify every Java override is present.

#### `record_type.go` + `record_type_builder.go` — matching Java
**Status: subset** (278 + 107 LoC). Missing: nominal-vs-structural
duality in `isRecordType`; record-type printer simplified.

#### Singleton + primitive types (`value_type`, `boolean_type`, `number_type`, `string_type`, `null_type`, `void_type`, `big_int_type`, `boolean_literal_set`, `symbol_type`, `known_symbol_type`, `unknown_type`, `all_type`, `no_type`, `no_object_type`, `no_resolved_type`)
**Status: near-complete.** One-file-per-Java-class, 10-80 LoC each.
Structure matches Java.

#### `jstype_class.go` + `jstype_native.go` + `jstype_iterations.go` + `jstype_resolver.go` + `type_string_builder.go`
**Status: near-complete** for the thin utility classes.
`JSTypeResolver` (241 Go LoC vs. ~500 Java LoC) is missing late-
resolution ordering and the `Phase` enum's effect on reporting
suppression — shows up in several position-mismatch bugs.

### 3.3 `internal/jscomp/` — type checker core

#### `internal/jscomp/typecheck.go` — `jscomp/TypeCheck.java` (3470 LoC)
**Status: broad subset.** 2826 LoC, 86 Go-side funcs covering
~80 % of Java's visitor surface.
- Missing visitors: `visitArrayLit`, full `visitObjectLit`
  (contextual typing + `JSC_ILLEGAL_OBJLIT_KEY` full coverage),
  `visitRegExpLit`, `visitTemplateLit`, `visitCast` +
  `checkForImplicitCast`, `checkForOfTypes`, `visitObjectPattern`,
  `checkAbstractMethodInConcreteClass`,
  `checkAbstractPropertyAccess`, `visitInterfacePropertyAssignment`
  (full branch), `checkDestructuringAssignment`,
  `checkPropertyInheritanceOnAssignment` sibling branches,
  `checkDeclaredPropertyAgainstPrototypalInheritance`.
- Missing diagnostics / branches: `JSC_INVALID_CAST`,
  `JSC_ILLEGAL_IMPLICIT_CAST`,
  `JSC_UNSUPPORTED_BOUNDED_GENERIC_TYPES` (Java marks in-dev —
  port the in-dev path as-is), `JSC_TEMPLATE_TRANSFORMATION_ON_CLASS`,
  `JSC_INCOMPATIBLE_ALIAS_ANNOTATION`,
  `JSC_PROPERTY_ASSIGNMENT_TO_READONLY_VALUE`,
  `JSC_CONFLICTING_GETTER_SETTER_TYPE`,
  `JSC_ILLEGAL_PROPERTY_CREATION_ON_UNION_TYPE`,
  `JSC_EXPECTED_THIS_TYPE`, `JSC_NOT_FUNCTION_TYPE`,
  `JSC_CYCLIC_INHERITANCE_ERROR`, `JSC_EXTENDS_WITHOUT_TYPEDEF`,
  `JSC_INTERFACE_METHOD_NOT_EMPTY` (partial),
  `JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH` (partial),
  `JSC_UNKNOWN_LENDS`, `JSC_LENDS_ON_NON_OBJECT`,
  `JSC_ILLEGAL_CLASS_KEY`, `JSC_BAD_TYPE_FOR_BIT_OPERATION`,
  `JSC_MULTIPLE_VAR_DEF`, `JSC_RESOLVED_TAG_EMPTY`,
  `JSC_DICT_EXTEND_STRUCT_TYPE`, `JSC_UNKNOWN_TYPEOF_VALUE`,
  `JSC_INEXISTENT_PARAM`, `JSC_INVALID_PARAM`,
  `JSC_CONSTRUCTOR_NOT_CALLABLE`,
  `JSC_INSTANTIATE_ABSTRACT_CLASS`, `JSC_UNKNOWN_OVERRIDE`.
- Verify: `reportMissingProperty` exact position (the 24 FP
  `JSC_TYPE_MISMATCH` bucket and the 11 wrong-position bucket
  both land here).
- Go-only: `isSafeCalleeForArity`, `isSafeInstanceMethodForArity`,
  `isDefinitelyNonConstructor`, `isInGuardedBranch`, `containsPropTest`,
  `otherSideIsNullOrUndefined`, `dereferencedObject`,
  `isKnownSurfaceBuiltin`, `defineOnOpenObject`,
  `buildExternTypeSet`, `buildLoosePropertyIndex`. These were
  introduced to work around missing pieces of `TypeValidator`,
  `NodeUtil`, and the externs bundle. Each needs to be migrated
  into its Java home and removed once the upstream surface lands.

#### `internal/jscomp/superclasscheck.go` — `TypeCheck.java` (inheritance-check subset)
**Status: subset.** 1162 LoC.
- Missing: constructor-arg supertype compatibility check;
  full `JSC_INTERFACE_METHOD_NOT_IMPLEMENTED` across multi-
  extends diamond; `JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH`
  full property-descriptor signature comparison (6-count
  bucket); `JSC_INCOMPATIBLE_EXTENDED_PROPERTY_TYPE` exact
  position (4-wrong-position bucket); `JSC_CYCLIC_INHERITANCE_ERROR`.
- Go-only: the `scanConstructorJSDoc` helper is a local pre-
  pass — Java does this inline in `TypeCheck.checkConstructor`.

#### `internal/jscomp/typevalidator.go` — `jscomp/TypeValidator.java` (1303 LoC)
**Status: subset.** 794 LoC, 29 `expect*` methods present.
- Missing: `expectInterfaceProperty`, `expectInterfaceProperties`
  (wire through `SuperclassCheck`); `expectAllInterfaceProperties`;
  `caseEnumElementType`, `caseFunctionType`, `caseNamedType`,
  `caseUnionType`, `caseTemplatizedType`, `caseTemplateType`
  visitor hooks — currently duplicated in
  `SubtypeChecker`/`typecheck.go`; `getReadableTypeName`;
  `mismatch` registration through `TypeMismatch` accumulator.
- Go-only: `tvMismatch`/`tvRegisterMismatchAndReportWithDiag` are
  the Go re-spelling of the overloaded Java `mismatch`/
  `registerMismatchAndReport`. Keep but rename to match Java's
  exact method names before declaring the port complete.

#### `internal/jscomp/typemismatch.go` — `jscomp/TypeMismatch.java`
**Status: near-complete** (200 LoC). Minor: `MismatchType` enum
missing a few late-added entries — audit against Java when
`TypeValidator` port lands.

#### `internal/jscomp/scope.go` — `jscomp/TypedScopeCreator.java` (3755 LoC)
**Status: broad subset.** 3054 LoC, 80 Go-side funcs vs. ~89 Java
methods — coverage is on par in count but Java's three builders
(`AbstractScopeBuilder`, `ClassScopeBuilder`, `FunctionScopeBuilder`)
are collapsed into free functions.
- Missing: full `AbstractScopeBuilder` / `ClassScopeBuilder` /
  `FunctionScopeBuilder` separation — requires the
  `Scope`/`TypedScope`/`TypedVar`/`Var`/`AbstractVar` split
  (skeletons exist). Destructuring parameter declarations
  (`declareDestructuringParameter`, `declareNamesInPositionalParameter`);
  ES6 module imports (`defineModuleImport`,
  `declareExportsInGoogModuleScope`,
  `finishDeclaringGoogModule`, `markGoogModuleExportsAsConst`);
  `Object.defineProperty`/`defineProperties` property
  registration (`defineObjectLiteral` helpers);
  `declareNativeFunctionType`/`declareNativeType`/
  `declareNativeValueType` (tied to `JSTypeRegistry` port);
  `declareTemplatesInFunctionBody`; `inferTypeForDestructuredTarget`;
  `getLValueRootScope`; `findOverriddenProperty` (shared with
  `SuperclassCheck`); `determineOwnerTypeForClassMember`;
  `isTypedConstantDeclaration`; `isValidTypedefDeclaration`;
  `identifyEnumOrTypedefDeclaration`.
- Go-only: `resolveObjectLiteralType`, `PushRefinement`/
  `PopRefinement` sit in `TypedScope` but Java puts the
  equivalent in `FlowScope`. Move when `FlowScope` port lands.

#### `internal/jscomp/inference.go` + `typeinferencepass.go` + `refinement.go` + `flowscope.go` — `jscomp/TypeInference.java` (2849 LoC) + `type/FlowScope.java`
**Status: subset.** 1087 + 198 + 287 + 135 LoC vs. Java's 2849-line
`TypeInference` alone.
- Missing traversers (Java has 41 `traverse*`, Go has 11 `infer*`):
  `traverseAnd`, `traverseOr`, `traverseNullishCoalesce` (Go
  collapses all three through `inferUnder`); `traverseArrayLiteral`;
  `traverseAwait`, `traverseYieldAll`; `traverseAssignOp`,
  `traverseShortCircuitingBinOpAssignment`,
  `traverseAssignUnsignedRightShift`, `traverseUnsignedRightShift`;
  `traverseBigIntCompatibleBinaryOperator` / `*UnaryOperator*`;
  `traverseCatch`; `traverseClass`, `traverseClassMemberRhs`;
  `traverseDeclaration`, `traverseDeclarationChild`;
  `traverseDestructuringPattern` (+ helper);
  `traverseDynamicImport`; `traverseImportSpec`;
  `traverseInstantiation`; `traverseNewTarget`; `traverseOptChain`;
  `traverseReturn`; `traverseSuper`; `traverseTaggedTemplateLit`;
  `traverseUnaryPlus`; `traverseWithinShortCircuitingBinOp`.
- Missing: loop-aware dataflow (`while (!x) x = u();` must refine
  `x` after the body); BigInt/Symbol/regex literals;
  compound-assign result-type sharpening; assignment-with-union-LHS;
  falsy-on-union-with-literal-false; nested AND/OR sub-expression
  refinement.
- `flowscope.go` is a thin map (Bind/BindAll/Equal/Get). Java's
  `FlowScope` is a lattice with
  `inferSlotType`/`inferQualifiedSlot`/`createChildFlowScope`
  and a proper `JoinOp`. Full port blocks several `JSC_TYPE_MISMATCH`
  misses.
- Missing semantic reverse-abstract-interpreter: `in`, `delete`,
  `instanceof` discrimination on union alternates, `==` vs.
  `===` refinement over literals, tuple-destructure discrimination.
  Skeletons at `skel/jscomp/typeinfer/`.
- Go-only: `InferType`, `inferUnder`, `inferWithNarrowing` are
  local simplifications of Java's `getTypeOfNode`/`traverse`/
  `dataflow` entry points. Rename to the Java spellings as part
  of the widening.

#### `internal/jscomp/parsing/jsdocinfoparser.go` — `parsing/JsDocInfoParser.java` (2895 LoC)
**Status: subset.** 869 LoC, 35 funcs vs. Java's 104 methods.
- Missing tags: `@enum` fully (type + values); `@const` with
  optional type; `@dict` / `@struct` (flag carried but not
  checked here); `@nosideeffects` + `@noside`; `@abstract`;
  `@public` / `@protected` / `@private` (parsed but not threaded
  to visibility); `@param` modifiers `{=T}`, `{...T}`, `{T=}`,
  default values; callable types `function(A,B):C`,
  `function(this:T,…)`, `function(new:T,…)`; multi-`@extends`
  accumulator (interfaces); multi-`@implements`; `@template`
  with bounds (parser reads them, no downstream enforcement);
  `@lends`; `@typedef` over record/union/templatized forms;
  `@ngInject`, `@idGenerator`, `@wizaction`, `@polymerBehavior`.
- Missing: position-tracked errors inside type expressions —
  drives the 7-count `JSC_UNRECOGNIZED_TYPE_ERROR` bucket and
  the 1-count wrong-position case.

#### `internal/jscomp/parsing/jsdoctoken.go` + `jsdoctokenstream.go`
**Status: near-complete** for the tags currently parsed (80 + 298
LoC). Add the trailing tag set when `JsDocInfoParser` widens.

### 3.4 `internal/jscomp/` — passes

#### `internal/jscomp/checkaccesscontrols.go` — `jscomp/CheckAccessControls.java` (1466 LoC)
**Status: broad subset.** 793 LoC, 19 Go funcs vs. ~37 Java
methods. Wired into the pass list.
- Missing: `canAccessDeprecatedTypes`; `canOverrideVisibility`;
  `checkPackagePropertyVisibility`, `checkPrivatePropertyVisibility`,
  `checkProtectedPropertyVisibility`, `checkPropertyOverrideVisibility`,
  `checkPropertyOverrideVisibilityIsSame`; `enterAccessControlScope`
  / `exitAccessControlScope` scope stack (currently one
  ambient scope); `getEffectiveVisibilityForNonOverriddenProperty`;
  `isPackageAccessAllowed`, `isPrivateAccessAllowed`;
  `isDocumentedDeclaration`; `isExtendsTarget`;
  `shouldEmitDeprecationWarning` precise conditions; `isDeletion`
  (for deprecation on delete). Drives the single-count
  `JSC_ILLEGAL_PROPERTY_ACCESS` bucket.
- Go-only: the `cacWalk` + inline traversal is a thin post-order
  substitute for Java's `ShouldTraverse`/`Visit` pair. Align
  naming when `NodeTraversal` widens.

#### `internal/jscomp/checkjsdoc.go` — `jscomp/CheckJSDoc.java` (799 LoC)
**Status: subset.** 734 LoC, 17 funcs vs. 33 Java methods.
- Missing: `checkJsDocInBlockComments`,
  `containsOnlySuppressionFor`, `validateClosurePrimitive`,
  `validateDefaultValue`, `validateImplicitCast`,
  `validateIsUsedViaDotConstructor`, `validateJsDocTypeNames`,
  `validateNoCollapse`, `validateRestParameter`,
  `validateTsType`. Drives single-count
  `JSC_INVALID_INTERFACE_MEMBER_DECLARATION`.

#### `internal/jscomp/checkjsdoctypes.go` — part of `TypeValidator.java`
**Status: subset** (1096 LoC). Reads `@type` expressions from
JSDoc and checks them against resolved types. Once
`TypeValidator.expect*` coverage is complete, collapse into
`typevalidator.go` to match Java layout.

#### `checkmissingreturn.go` + `checksuspiciouscode.go` + `checksuper.go` + `checkglobalthis.go`
**Status: near-complete.** Sizes vs. Java: 472/203, 476/281,
556/479, 244/208. Residuals: one FP `JSC_MISSING_RETURN_STATEMENT`
on `??` nullable (CFG defect — revisit when `ControlFlowAnalysis`
widens); single-count `JSC_SUSPICIOUS_IN` and
`JSC_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR` residuals
(verify against Java branch); 3-count
`JSC_TYPE_MISMATCH,JSC_USED_GLOBAL_THIS` where the missing
`TYPE_MISMATCH` is a `typecheck.go` gap, not a `checkglobalthis`
gap.

#### `internal/jscomp/checkmissingrequires.go` — `CheckMissingRequires.java` (839 LoC)
**Status: broad subset** (1083 LoC, covers the main branches;
module metadata is stubbed).
- Missing: `ModuleMetadataMap` wired through `CheckClosureImports`
  + the real module graph (currently a local stub); `isIIFE`,
  `isGoogScopeBody`, `isTypeScriptSource` are present; verify
  `swapRootOfQualifiedName` against Java.

#### `internal/jscomp/checkclosureimports.go` — `CheckClosureImports.java` (551 LoC)
**Status: subset** (628 LoC). Handles the common `goog.require`
patterns; stubs the module metadata map.
- Missing: `checkValidImportCodeReference`; `enterModule`
  tracking; `isValidDestructuringImport` (narrow). Drives the
  3-count `JSC_INVALID_MODULE_PATH` bucket and the 2-count
  `JSC_PARSE_ERROR` misroute.

#### `checkregexp.go` + `checkdebuggerstatement.go` + `checkunreachablecode.go` + `checktypeimportcodereferences.go`
**Status: near-complete.** Straightforward passes — 232/~40,
77/small, 340/95, 193/small LoC. `CheckConformance.java`
(`checkconformance.go`, 112 LoC) is a **stub** with only the
entry point; low priority — no TypeCheckTest fixture exercises
conformance.

#### `internal/jscomp/checksideeffects.go` — `CheckSideEffects.java`
**Status: subset** (211 LoC).
- Missing: `@nosideeffects`-aware CALL classification;
  sub-expression `a[b].c;` reporting; the ModifiedStateChecker
  sub-pass.

#### `internal/jscomp/checkdupvardecl.go` — a subset of `TypeCheck.java` + `VariableReferenceCheck.java`
**Status: subset** (397 LoC). Carries the bulk of
`JSC_DUP_VAR_DECLARATION_TYPE_MISMATCH` detection.
- Missing: the `JSC_FUNCTION_MASKS_VARIABLE` arm (2-count bucket);
  `JSC_VAR_MULTIPLY_DECLARED_ERROR` cross-file variant (1-count
  residual).

#### `internal/jscomp/checkduplicatemember.go` — `CheckSuspiciousCode.java` (duplicate-member branch)
**Status: subset** (124 LoC).
- Missing: class-body member checks (class fields + methods
  conflict detection).

#### `internal/jscomp/checkdeterministictest.go` — `JSC_DETERMINISTIC_TEST` branch of `TypeCheck.java`
**Status: subset** (579 LoC). Pure-miss count is 4, position-
mismatch count is 3, FP count is 1 — so both broader and
narrower than Java in different spots.
- Verify: exact Java dispatch for symbol/number cross-type
  comparison; exact position of report when the comparison sits
  inside a `|| void 0` chain.
- Add a skeleton for this when `gen-skeletons.sh` runs next —
  there isn't one in the current tree because Java keeps the
  check inline in `TypeCheck.java`.

#### `internal/jscomp/varcheck.go` — `VarCheck.java`
**Status: subset** (197 LoC).
- Missing: per-scope unreached-NAME reporting (drives the 2-
  count `want=[JSC_UNDEFINED_VARIABLE] got=[]` bucket);
  `JSC_REFERENCE_BEFORE_DECLARE` for `var` hoisting edge cases
  (3-count `JSC_REFERENCE_BEFORE_DECLARE,JSC_TYPE_MISMATCH`
  bucket).

#### `internal/jscomp/varrefcheck.go` — `VariableReferenceCheck.java`
**Status: subset** (500 LoC).
- Missing: `checkEarlyReference` across hoisting of
  function-scoped `var`; `JSC_VAR_MULTIPLY_DECLARED_ERROR` per-
  script aggregation; `checkFunctionMasksVariable` class-scope
  overload.

#### `inferjsdocinfo.go` + `inferconsts.go` + `typetransformation.go` + `ast_analyzer.go`
**Status: subset.** 347, 284, 987, 203 LoC. `InferJsDocInfo` misses
prototype-assignment branches. `TypeTransformation` handles the
standard TTL operators; `JSC_TEMPLATE_TRANSFORMATION_ON_CLASS`
(2-count bucket) is reported from `typecheck.go` only when the
TTL transform is detected. `AstAnalyzer` covers side-effect +
mutation predicates used by `CheckSideEffects`.

#### `symboltable.go` + `scoped_name.go` — `SymbolTable.java` + `ScopedName.java`
**Status: subset** (1774 + 100 LoC). Currently unused by the
corpus; kept to satisfy `TypedScopeCreator` references.

#### `internal/jscomp/basicblock.go` + `cfg.go` + `controlflowgraph.go` + `controlflowanalysis.go` + `dataflow.go`
**Status: CFG/DFG subset.** `controlflowanalysis.go` 735 LoC vs.
1050 Java LoC; `dataflow.go` 146 vs. 499.
- Missing: labeled `break`/`continue` to nested loops; exception
  fan-out to nearest catch per call/throw; `with` modelling
  (low-value but needed for conformance parity); generator
  suspend edges (`yield`); async/await suspend edges;
  assignment-in-condition edges; `finally` re-raise edges.
- Missing (DFG): backward direction; lattice side-table
  (currently annotated directly on `CfgNode`); iteration-limit
  heuristic; `JoinOp` as a proper interface.

#### `internal/jscomp/nodetraversal.go` — `NodeTraversal.java`
**Status: subset** (201 LoC). Pre/post order + scope callbacks
present. Missing: `traverseFunctionOutOfBand`,
`traverseWithScope`, the various `traverseRoots` helpers, the
`ScopeCreator` two-pass protocol's post-order phase.

#### `internal/jscomp/nodeutil.go` — `NodeUtil.java` (5851 LoC!)
**Status: thin subset.** 633 LoC, 28 funcs vs. ~127 Java methods.
**This is a high-leverage gap** — every pass reaches for
`NodeUtil` helpers and currently re-implements them locally.
- Missing (selection): `getAssignedValue`, `isBleedingFunctionName`,
  `isConstantDeclaration`, `isExpressionResultUsed`,
  `isFunctionObjectLit`, `isGetPropOnThis`,
  `isInFunction`, `isInSyntheticScript`, `isNameNodeInsideExtends`,
  `isNonEs6ClassExpression`, `isObjectLitKey`,
  `isPrototypeMethod`, `isPrototypeProperty`,
  `mayBeObjectLitKey`, `mayBeString`, `mayEffectMutableState`,
  `newQualifiedNameNode`, `removeChildren`, `isLhsByDestructuring`,
  `tryGetStringFromBoolean`, `getFunctionParameters`,
  `getBestLValue`, `getBestLValueOwner`, `getBestJSDocInfo`.

#### `internal/jscomp/nodeiterators.go` + `nodenameextractor.go`
**Status: near-complete** (343 + 202 LoC).

#### `internal/jscomp/compiler.go` + `compileroptions.go`
**Status: thin shim.** `compiler.go` 112 LoC, `compileroptions.go`
65 LoC.
- Missing: `AbstractCompiler.java`; `PhaseOptimizer`;
  `Result`/`JSError` table; module graph; source-map
  context; `DefaultPassConfig` with a real ordered pass list;
  `PassConfig`/`PassFactory`/`PassListBuilder`/`PassNames`/
  `CompilerPass`; `CompilerInput`/`CompilerInputProvider`;
  `CompilerOptionsPreprocessor`. The current driver wires
  passes imperatively in `Check`.

#### `internal/jscomp/default_externs.go` — externs bundle
**Status: subset** (625 LoC, ES5 + a sliver of DOM).
- Missing: the full ES2022 surface (Iterator, AsyncIterator,
  WeakRef, FinalizationRegistry); Web IDL externs (Window.*,
  HTMLElement.*, typed arrays, Intl.*); `Promise` static and
  instance methods; `Symbol.*` well-known keys as declared
  properties; `goog.*` runtime surface. Drives the 2-count
  `want=[JSC_UNDEFINED_VARIABLE] got=[]` bucket and several
  `JSC_TYPE_MISMATCH` FPs.

### 3.5 `internal/parser/`

#### `internal/parser/scanner.go` — `parsing/parser/Scanner.java` (1375 LoC)
**Status: narrow.** 632 LoC character-level scanner.
- Missing: full `TokenType` enum; string-literal escapes
  + template pieces; numeric literals
  (hex/octal/binary/`_`-separators/`n` bigint suffix/decimal
  edge cases); regex literals; comment/trivia preservation
  through `JsDocTokenStream`; source-position stream.

#### `internal/parser/parser.go` — `parsing/parser/Parser.java` (3694 LoC)
**Status: narrow.** 2561 LoC covering a JS subset tailored to the
corpus fixtures.
- Missing: array and object literals at expression position;
  destructuring patterns (array+object, nested, default, rest);
  spread/rest args and params; template literals + tagged templates;
  regex literals as expressions; `GETELEM` expression form;
  optional chaining (`?.`, `?.(`, `?.[`); `async`/`await`;
  generator functions (`function*`, `yield`, `yield*`);
  computed property syntax (`{[k]:v}`, `class { [k]() {} }`);
  cast `(/** @type {T} */ (e))` as first-class grammar;
  class-body: fields + static fields + accessors + private `#`
  + static blocks; `for`/`for-in`/`for-of` with destructured
  binding; `import`/`export` declarations (named, default,
  namespace, re-export); `super.x`, `super(…)`; labeled
  statements. ~86 tree node types under Java's
  `parsing/parser/trees/` — Go side uses the rhino `Node` model
  directly, so each syntactic form needs a dedicated materialiser
  inside `parser.go` when porting `IRFactory.java`.

### 3.6 `internal/diag/diag.go`
**Status: subset.** 197 LoC. Mirrors the cut of
`DiagnosticType.java` + `JSError.java` that jsgo emits. Widen as
new passes lift their codes.

### 3.7 Test scaffolding

`internal/jscomp/strictcheck_test.go` (67 LoC) tracks
`TypeCheckStrictTest.java`. Intentionally narrow.

---

## 4. Unported Java files (skeletons in `internal/jscomp/skel/`)

Each Java file listed below has a companion skeleton under
`internal/jscomp/skel/`. Porting = move the file out of `skel/`,
drop the `//go:build skeletons` tag, fill in the body, wire into
`compiler.go`, add any new codes to `internal/diag/diag.go`.

### 4.1 `rhino.*` core (no Go counterpart yet)

`ClosurePrimitive`, `ErrorReporter`, `IR`, `NominalTypeBuilder`,
`NonJSDocComment`, `Outcome`, `SourcePosition`, `TokenStream`,
`TokenUtil`, `TypeDeclarationsIR` — all in `skel/rhino/`.
`HamtPMap`, `PMap`, `RhinoStringPool`, `PropTranslator`, `Msg`,
`JSIdentifier` have no skeleton (utility classes the Go side
collapses into stdlib equivalents — low priority).

### 4.2 `rhino.jstype.*` (no Go counterpart)

`EquivalenceMethod`, `RelationshipVisitor` — skeletons under
`skel/rhino/jstype/`. `SimpleReference`, `SimpleSlot`,
`StaticTypedRef`, `StaticTypedScope`, `StaticTypedSlot` —
Java-side helper interfaces; fold into `internal/rhino/` when
the `Scope`/`TypedScope`/`TypedVar` split lands.
`AutoValue_FunctionType_Parameter` / `AutoValue_JSTypeRegistry_ClosureNamespace`
are generated and inlined on the Go side.

### 4.3 `jscomp.*` infrastructure

Driver / options / diagnostics (all in `skel/jscomp/`):
`AbstractCompiler`, `Compiler` (full rewrite over the current
shim), `CompilerInput`, `CompilerInputProvider`, `CompilerOptions`
(full), `CompilerOptionsPreprocessor`, `CompilerPass`,
`DefaultPassConfig`, `PassConfig`, `PassFactory`,
`PassListBuilder`, `PassNames`, `AbstractMessageFormatter`,
`AllowlistWarningsGuard`, `CheckLevel`, `DiagnosticType`,
`DiagnosticGroup` (full), `DiagnosticGroups`,
`DiagnosticGroupWarningsGuard`,
`DiagnosticGroupPathSuppressingWarningsGuard`,
`DiagnosticToSuppressionMapper`, `JSError` (full),
`DefaultExterns`.

### 4.4 `jscomp.*` scope hierarchy

`AbstractScope`, `AbstractVar`, `Scope`, `ScopeCreator`,
`TypedScope`, `TypedScopeCreator` (full), `TypedVar`, `Var` — all
in `skel/jscomp/`. Required prerequisite for splitting the
current monolithic `internal/jscomp/scope.go` into the Java
builder classes.

### 4.5 `jscomp.*` type-check infrastructure

Full 1:1 rewrites pending for: `TypeCheck`, `TypeValidator`,
`TypeMismatch`, `TypeInference`, `TypeInferencePass`,
`TypeMatchingStrategy`, `TypeTransformation`. Current Go ports
are subsets (see section 3) — skeleton body is the un-narrowed
target.

### 4.6 `jscomp.type.*` (inference plumbing)

`FlowScope`, `ReverseAbstractInterpreter`,
`ChainableReverseAbstractInterpreter`,
`ClosureReverseAbstractInterpreter`,
`SemanticReverseAbstractInterpreter` — all in
`skel/jscomp/typeinfer/`. Required for full refinement coverage
(several `JSC_TYPE_MISMATCH` buckets).

### 4.7 `jscomp.*` utilities

`AstFactory`, `AstManipulations`, `AstValidator`, `NodeUtil`
(full — compare §3.4 entry), `NodeTraversal` (full),
`NodeIterators`, `NodeNameExtractor`, `ControlFlowAnalysis`
(full), `ControlFlowGraph` (full), `DataFlowAnalysis` (full),
`BasicBlock` (full), `SymbolTable` (full), `ScopedName` (full) —
all in `skel/jscomp/`.

### 4.8 `jscomp.parsing.*`

Full rewrites: `IRFactory` (3900+ LoC Java; gates on parser
widening), `JsDocInfoParser`, `JsDocToken`, `JsDocTokenStream`,
`ParserRunner`, `TypeTransformationParser`. No skeleton for
`Config`, `AutoValue_Config`, `Annotation`, `FeatureCollector`,
`ParserConfiguration`, `ParsingUtil` — generate on next
`scripts/gen-skeletons.sh` run.

### 4.9 `jscomp.parsing.parser.*` (the scanner/parser + 86 trees)

Skeletons exist for `Scanner`, `Parser`, `Token`, `TokenType`,
`IdentifierToken`, `LiteralToken`, `TemplateLiteralToken` under
`skel/jscomp/parsing/parser/`. No skeletons for the 86 tree
classes under `parsing/parser/trees/*.java` nor for
`ClassOrObjectElementInfo`, `FeatureSet`, `Identifiers`,
`Keywords`, `LineNumberScanner`, `PredefinedName`, `SourceFile`,
`StringLiteralToken`, `util/*`. Go port currently materialises a
trimmed tree directly on `rhino.Node` — either port every tree
class 1:1 or keep the current approach and document the
divergence.

### 4.10 Notes on coverage

- `CheckExtraRequiresPassConfig.java` — no skeleton; generate on
  the next `gen-skeletons.sh` run.
- `CheckDeterministicTest` is not a Java file (inlined in
  `TypeCheck.java`); Go's `checkdeterministictest.go` split it out
  for locality. Merge back when `TypeCheck.java` is ported 1:1.

---

## 5. Diagnostic-tuple gaps (ledger)

Re-derived from `go run ./tools/fbucket b` at `9816618`. Listed so
nothing is missed. Re-run after every landing.

### 5.1 Pure misses (`want=[X] got=[]`)

| Count | Diagnostic | Fix-site |
|------:|-----------|----------|
| 179 | `JSC_TYPE_MISMATCH` | `inference.go`, `typecheck.go`, `typevalidator.go` + externs widening |
|  11 | `JSC_NON_STRINGIFIABLE_OBJECT_KEY` | `typecheck.go` `visitGetElem` / computed-prop + `findObjectWithNonStringifiableKey` |
|   8 | `JSC_INEXISTENT_PROPERTY` | `typecheck.go` non-structural branch + `TypeValidator.caseObjectType` |
|   7 | `JSC_UNRECOGNIZED_TYPE_ERROR` | `jsdocinfoparser.go` position-tracked errors + `scope.go` resolve-time |
|   6 | `JSC_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH` | `superclasscheck.go` class-style branch |
|   5 | `JSC_INVALID_CAST` | `typecheck.go` `visitCast` |
|   4 | `JSC_DETERMINISTIC_TEST` | `checkdeterministictest.go` dispatch on symbol/number cross-type |
|   3 | `JSC_INVALID_MODULE_PATH` | `checkclosureimports.go` validate-path branch |
|   3 | `JSC_HIDDEN_INTERFACE_PROPERTY_MISMATCH` | `typecheck.go` + `superclasscheck.go` |
|   3 | `JSC_PROPERTY_ASSIGNMENT_TO_READONLY_VALUE` | `typecheck.go` new branch |
|   2 | `JSC_UNDEFINED_VARIABLE` | `varcheck.go` unreached-NAME; externs widening |
|   2 | `JSC_TEMPLATE_TRANSFORMATION_ON_CLASS,JSC_UNRECOGNIZED_TYPE_ERROR` | `typecheck.go` @template-TTL branch |
|   2 | `JSC_EXPECTED_THIS_TYPE` | `typecheck.go` call-with-@this branch |
|   2 | `JSC_MISSING_RETURN_STATEMENT,JSC_TYPE_MISMATCH` | `checkmissingreturn.go` + `typevalidator.go` |
| 1 each | `JSC_CYCLIC_INHERITANCE_ERROR`, `JSC_INCOMPATIBLE_EXTENDED_PROPERTY_TYPE`, `JSC_ILLEGAL_PROPERTY_CREATION_ON_UNION_TYPE`, `JSC_UNKNOWN_TYPEOF_VALUE`, `JSC_NOT_FUNCTION_TYPE`, `JSC_INEXISTENT_PARAM`, `JSC_UNKNOWN_OVERRIDE`, `JSC_RESOLVED_TAG_EMPTY`, `JSC_EXTENDS_WITHOUT_TYPEDEF`, `JSC_ILLEGAL_PROPERTY_CREATION`, `JSC_ILLEGAL_IMPLICIT_CAST`, `JSC_INSTANTIATE_ABSTRACT_CLASS`, `JSC_UNKNOWN_LENDS`, `JSC_LENDS_ON_NON_OBJECT`, `JSC_ENUM_INITIALIZER_NOT_ENUM`, `JSC_INVALID_INTERFACE_MEMBER_DECLARATION`, `JSC_INVALID_PARAM`, `JSC_TYPE_PARSE_ERROR`, `JSC_MISSING_RETURN_STATEMENT` (single) | see section 3 |

### 5.2 False positives (`want=[] got=[X]`)

| Count | Spurious diagnostic | Fix-site |
|------:|---------------------|----------|
| 24 | `JSC_TYPE_MISMATCH` | `inference.go` loop-aware refinement + `typevalidator.go` port + externs widening |
|  4 | `JSC_INEXISTENT_PROPERTY` | broaden "unknown-structure" receiver recognition in `typecheck.go:checkPropertyAccess` |
|  4 | `JSC_POSSIBLE_INEXISTENT_PROPERTY` | same receiver check, narrow side |
|  1 | `JSC_NON_STRINGIFIABLE_OBJECT_KEY` (`testenumofsymbol1.js`) | `checkTypeContainsObjectWithBadKey` on Symbol-keyed enum |
|  1 | `JSC_MISSING_RETURN_STATEMENT` | CFG defect on `??` nullable; revisit after `ControlFlowAnalysis` widening |

### 5.3 Position mismatches (same code, wrong line/col)

Same diagnostic reported at the wrong location — usually because
we emit on a parent node while Java emits on a child (or vice
versa). Fix each by consulting Java's exact report-site.

- 11× `JSC_TYPE_MISMATCH`
- 3× `JSC_DETERMINISTIC_TEST,JSC_USELESS_CODE`
- 1× `JSC_UNRECOGNIZED_TYPE_ERROR`
- ~10 further single-count combinations.

### 5.4 Mixed tuples (partial match, mostly over-reporting)

| Count | Shape |
|------:|-------|
| 4 | `want=[JSC_TYPE_MISMATCH,JSC_USELESS_CODE]` got only `[JSC_USELESS_CODE]` |
| 3 | `want=[JSC_TYPE_MISMATCH,JSC_USED_GLOBAL_THIS]` got only `[JSC_USED_GLOBAL_THIS]` |
| 3 | `want=[JSC_UNDEFINED_VARIABLE]` got `[JSC_TYPE_MISMATCH,JSC_UNDEFINED_VARIABLE]` (we over-report TYPE_MISMATCH) |
| 3 | `want=[JSC_REFERENCE_BEFORE_DECLARE,JSC_TYPE_MISMATCH]` got `[]` |
| 2 | `want=[JSC_POSSIBLE_INEXISTENT_PROPERTY,JSC_USELESS_CODE]` got only `[JSC_POSSIBLE_INEXISTENT_PROPERTY]` |
| 2 | `want=[JSC_UNDEFINED_VARIABLE]` got `[JSC_POSSIBLE_INEXISTENT_PROPERTY,JSC_UNDEFINED_VARIABLE]` |
| 2 | `want=[JSC_UNDEFINED_VARIABLE]` got `[JSC_UNDEFINED_VARIABLE,JSC_USELESS_CODE]` |
| 2 | `want=[JSC_INEXISTENT_PROPERTY]` got `[JSC_INEXISTENT_PROPERTY,JSC_POSSIBLE_INEXISTENT_PROPERTY]` |
| 2 | `want=[JSC_TYPE_MISMATCH]` got `[JSC_INEXISTENT_PROPERTY]` (wrong code) |
| 2 | `want=[JSC_TYPE_MISMATCH]` got `[JSC_POSSIBLE_INEXISTENT_PROPERTY]` (wrong code) |
| 2 | `want=[JSC_PARSE_ERROR]` got `[JSC_POSSIBLE_INEXISTENT_PROPERTY,JSC_UNDEFINED_VARIABLE]` (module-path misroute) |

These drop out automatically as the underlying pure-miss / FP
entries resolve.

---

## 6. Supporting work

- Populate the hold-out corpus slot (SPEC guardrail #5) and keep
  its pass-rate locked to the visible corpus.
- Performance: stay within `[0.1×, 2×]` of Java on Tier-1,
  verified at each major port landing.
- Oracle regeneration: re-run `scripts/gen-oracles.sh` whenever
  fixture layout or Java pinned version changes.
- Triage tooling: `tools/fbucket` keeps the bucket table; re-run
  after every significant change, paste the snapshot back into
  this file (sections 1 + 2 + 5).
- `tools/fbucket` itself covers default mode only. For strict
  mode, invoke `jscomp.CheckSourceWithOptions` with
  `GroupStrictCheckTypes=LevelWarning` over `testdata/oracle-strict/…`.
  Consider extending `fbucket` with a `-strict` flag.
- CI floor in `harness/harness_test.go` raised alongside progress
  so regressions show up immediately.
- Skeleton refresh: re-run `scripts/gen-skeletons.sh` when the
  Java tag is bumped, to keep `internal/jscomp/skel/` in sync.
  Add missing generators for `CheckExtraRequiresPassConfig`,
  `Config`/`ParserConfiguration`, `parsing/parser/trees/*`.
