# Census of `*rhino.Node` holders across pass boundaries

Input to round-2 step 1 (`round2-design.md`, "The hard part"): a pointer-free,
compacting AST moves nodes, so every `*rhino.Node` that something other than the
tree holds across a compaction point has to be found and forwarded. This file
lists those holders (static census, section 1), measures how many nodes they hold
and how many of those are already detached from the tree at each candidate
compaction point (dynamic census, section 2), and says where compaction is safe
(section 3).

Method. Static: grep over `internal/` (tests excluded) for `*rhino.Node`,
`[]*rhino.Node`, `map[*rhino.Node]`, `LinkedHashMap[*rhino.Node, …]`, closures
returning nodes, and the struct bodies of every long-lived type (Compiler and
everything reachable from it, rhino types, colors, serialization). Dynamic: a
temporary hook (`.gocache/holder-census.patch`, not committed; enabled with
`JSGO_HOLDER_CENSUS=1`, `JSGO_HOLDER_CENSUS_EVERY=N`,
`JSGO_HOLDER_CENSUS_PASSES=a,b`) that, right after `typesToColors`, at the end
of every optimization-loop iteration, every 10 passes and at stage ends, walks
the tree from `Compiler.externAndJsRoot` into a visited set (plus the side trees
reachable from live nodes: JSDoc type expressions, `DECLARED_TYPE_EXPR`,
`CLOSURE_UNAWARE_SHADOW`), then enumerates every holder category below and
counts refs, distinct nodes, nodes not in the visited set ("detached"), the
distinct roots of those detached subtrees and their total size ("retained").
Raw outputs: `.gocache/holder-census-typed-{app,full}.txt`; runner:
`.gocache/holder-census-run.py`; hooked binary: `.bin/jsgo-census`.

## 1. Static census

"Java" = the Java Closure Compiler keeps the same reference (the port is 1:1 for
everything named after a Java field; Go-only holders are marked).

### 1.1 Roots and per-input roots (always attached; trivial to forward)

| holder | what | lifetime | Java |
|---|---|---|---|
| `Compiler.externsRoot/jsRoot/externAndJsRoot` (`internal/jscomp/compiler.go:91-93`) | the 3 ROOT nodes | whole compilation | yes |
| `Compiler.scriptNodeByFilename` (`compiler.go:106`, `LinkedHashMap[string,*Node]`) | one SCRIPT per input (+ synthetic externs) | whole compilation; cleared/refilled on re-parse (`compiler.go:2456, 5236`) | yes |
| `CompilerInput_JsAst.root` (`compiler_input.go:801`) via `Compiler.inputsById`, `JSChunkGraph` → `JSChunk.inputs`, `Compiler.externs` | SCRIPT per input | whole compilation | yes |
| `PhaseOptimizer.jsRoot` (`phase_optimizer.go:45`) | jsRoot | per stage | yes |
| `PhaseOptimizer_Loop.Process` locals `scriptSizes map[*Node]int`, `changedScopes` (`phase_optimizer.go:418-421, 487-516`) | SCRIPT nodes keyed for the AST-size heuristic | **across loop iterations** (not pass-local) | yes |
| `RuntimeJsLibManager.nodeForCodeInsertion` (`js_runtime_js_lib_manager.go:38`, closure set at `compiler.go:1030-1036`) | captures the Compiler, not a node | whole compilation | yes |

### 1.2 Change tracking

| holder | what | lifetime | Java |
|---|---|---|---|
| `ChangeTracker.changeTimeline` (`change_tracker.go:35`; `Timeline[*Node]`: `eventsByValue map[*Node]uint32` + `chunks` of events, `timeline.go:40-53`) | every change-scope root (FUNCTION/SCRIPT) ever reported, in order; removed only by `ReportFunctionDeleted` (`change_tracker.go:108`) | whole compilation; 5-25k entries | yes (`Compiler.changeTimeline`) |
| `ChangeVerifier.clonesByCurrent` (`change_verifier.go:33`) | clone map | validity checks only (debug) | yes |

### 1.3 Modules (set in checks, never cleared)

| holder | what | lifetime | Java |
|---|---|---|---|
| `ModuleMetadataMap_ModuleMetadata.rootNode` (`modules_module_metadata_map.go:143`; `Compiler.moduleMetadataMap`, `compiler.go:264`) | SCRIPT/MODULE_BODY per module (+ nested) | whole compilation | yes |
| `Binding.sourceNode` (`modules_binding.go:48`), `Export.exportNode/nameNode` (`modules_export.go:58-59`), `Import.importNode/nameNode` (`modules_import.go:46-47`) via `ModuleMap` → `Module.namespace/boundNames/localNameToLocalExport` (`modules_module.go:79-81`; `Compiler.moduleMap`, `compiler.go:266`) | declaration nodes inside module bodies | whole compilation; readers after checks: `GlobalNamespace.SetupModuleMetadata` (`global_namespace.go:1146`, used by inlineAndCollapseProperties/polymer), `polymer_pass.go:73` | yes |

### 1.4 Types (alive from `inferTypes` until `typesToColors`)

Dropped wholesale by `ConvertTypesToColors.Process` → `Compiler.ClearJSTypeRegistry`
(`serialization_convert_types_to_colors.go:170,192`; `compiler.go:2191-2197` nils
`typeRegistry`, `typeValidator`, `abstractInterpreter`).

| holder | what | Java |
|---|---|---|
| `JSTypeRegistry.scopedNameTable` keys (`internal/rhino/jstype_js_type_registry.go:171`) | scope root nodes (1.2k / 3.0k) | yes |
| `JSTypeRegistry.nameTableGlobalRoot` (`:179`) | synthetic global root (always detached) | yes |
| `JSTypeRegistry.nonNullableTypeNames` keys (`:200`) | scope roots | yes |
| `JSTypeRegistry_ClosureNamespace.definitionNode` (`:3006`) | goog.provide/module call (detached after module rewriting) | yes |
| `FunctionTypeBase.source` (`jstype_function_type.go:146`) | FUNCTION nodes (8.9k / 26.9k distinct) | yes |
| `Property.propertyNode` (`jstype_property.go:73`) | declaration nodes (28.6k / 225k distinct) | yes |
| `EnumType.source` (`jstype_enum_type.go:53`), `TemplateType.typeTransformation` (`jstype_template_type.go:50`) | | yes |
| `TypeMismatch.location` (`type_mismatch.go:42`) via `TypeValidator.mismatches` (`type_validator.go:47`) | | yes |
| `Node.TYPEDEF_TYPE` prop (`node.go:4818`) holds a JSType, `Node.typeOrColor`/ext box (`node.go:3104, 601-603`) | per node, leaves with the type column | yes |
| `Property.docInfo` (`jstype_property.go:76`), `JSTypeRegistry.objectDocInfo` (`:250`) | JSDocInfo (owns type-expression trees, 1.6) | yes |

### 1.5 Typed scopes (alive from `inferTypes` until the clear passes)

Cleared by `clearTypedScopeCreatorPass` / `clearTopTypedScopePass`
(`default_pass_config.go:742-761`, scheduled at `:2057-2058` and `:2209-2225`)
**unless `options.PreservesDetailedSourceInfo()`** (IDE mode, not the CLI).

| holder | what | Java |
|---|---|---|
| `Compiler.typedScopeCreator` → `TypedScopeCreator.memoized` (`typed_scope_creator.go:60`, `LinkedHashMap[*Node,*TypedScope]`), `reservedNamesForScope` (`:64`), `functionsWithNonEmptyReturns` (`:69`), `unresolvedNodes` (`:79`), `undeclaredNamesForClosure` (`:85`), `providedNamesFromCall` (`:88`), `weakImports[].moduleLocalNode` (`:162`) | 23k / 93k scope roots as keys | yes |
| `AbstractScopeBase.rootNode` (`abstract_scope.go:262`), `AbstractVarBase.nameNode` (`abstract_var.go:161`) through `TypedScope.vars`, reached from `memoized`, `Compiler.topScope` (`compiler.go:214`), `CompilerInput.typedScope` (`compiler_input.go:62`) | 31k / 238k var name nodes in memoized scopes; 13k / 51k in topScope; 0.8k / 98k in per-input module scopes | yes |

### 1.6 Node-owned side trees (not holders of tree nodes; trees the compaction must copy too)

| owner | what | lifetime | Java |
|---|---|---|---|
| `JSTypeExpression.root` (`js_type_expression.go:49`) reached from `JSDocInfo` (`js_doc_info.go:79` `propertyValues`; `TypeExprIter`) and `JSDocInfo.typeTransformations` (`:597`) | detached type-expression trees (62k / 159k nodes at parse; 2 / 24 after colors) | JSDoc is on nodes (`node.go:602`), on Properties and in the registry; type trees are dropped at `typesToColors` | yes |
| `JSDocInfo_Marker` positions (`js_doc_info.go:699-762`, `SourcePositionBase[*Node]`) | only with documentation parsing (not the CLI) | | yes |
| `Node.DECLARED_TYPE_EXPR` (`node.go:2649`, `*Node`), `Node.CLOSURE_UNAWARE_SHADOW` (`node.go:1410`, `*Node`) | detached subtrees hung off a node prop | per node; 0 in both tasks | yes |

### 1.7 Diagnostics, injection, transpilation namespace

| holder | what | lifetime | Java |
|---|---|---|---|
| `JSError.node` (`js_error.go:44`) in `SortingErrorManagerBase.messages` (`sorting_error_manager.go:42`) | the node a warning was reported on; read at report time (`js_error.go:367`, `suppress_doc_warnings_guard.go:83`, `closure_unaware_code_warnings_guard.go:57`) | whole compilation (4 / 31 entries) | yes |
| `RuntimeJsLibManager.lastInjectedLibrary` (`js_runtime_js_lib_manager.go:42`) | last injected `$jscomp` library statement; **read** by `injectLibrary` (`:494-499`) to place the next one | from first injection to end | yes |
| `Compiler.transpilationNamespace` (`compiler.go:110`, built lazily at `:2377-2394`): `GlobalNamespace.root/externsRoot/globalRoot` (`global_namespace.go:267-269`), `spreadSiblingCache` (`:270`), `defineCandidateScripts` (`:301`), `nameMap` → `GlobalNamespace_Name.refsSlice/refsIndex/declaration` (`:2245-2246`) → `GlobalNamespace_Ref.node` (`:3708`), `Ref.scope` → `Scope.rootNode` | names of externs + the first (runtime-library) script; 4k names / 20k refs in typed-full | first transpilation pass → end (only ES5 output / transpiling builds) | yes |
| `globalNamespace_globalRoot` (`global_namespace.go:325`) and ~15 other package-level sentinel nodes (`type_mismatch.go:96`, `replace_css_names.go:104,220`, `pure_function_identifier.go:147`, `function_argument_injector.go:47`, `chrome_pass.go:48`, `ijs_jsdoc_util.go:35`, `closure_coding_convention.go:458`, `check_closure_imports.go:183-187`, `locale_data_passes.go:96`) | process-global template/sentinel nodes, never attached | process | yes (static finals) |

### 1.8 Symbol tables (not built by the CLI)

`SymbolTable` (`symbol_table.go:76-97, 1772, 1946`; `BuildKnownSymbolTable` has no
callers outside tests) and `PreprocessorSymbolTable` (`preprocessor_symbol_table.go:35-42`,
built only when `PreservesDetailedSourceInfo()`, `:170-176`). Both Java-identical
holders of many nodes; relevant only for IDE-style embedding.

### 1.9 Go-only holders

| holder | what | notes |
|---|---|---|
| `cachedCCAArena` (`cross_chunk_reference_collector.go:54`): `TopLevelStatement.statementNode`, `Reference.nameNode`, `Var` slabs | arena kept between the two cross-chunk passes; evicted at stage ends (`compiler.go:3033`) | stale, never read; evict at compaction |
| `SourceFile.sourceInfoTemplateVal` (`source_file.go:81`, a `rhino.Node` **by value**) | source-info template copied into new nodes | not in any chunk; must stay representable (or become a plain struct) |
| `scriptNodeDeserializer_Runner.templateNodeVal` (`serialization_script_node_deserializer.go:101`, by value, pooled, zeroed) | same idea | |
| `nodeTraversalPool` (`node_traversal.go:878`, zeroed on release `:895-903`), `devirtStrSlabPool`, `costEstimatorPool`, `subclassRelSlabPool` | pooled scratch | no node pointers survive release (checked) |
| `Compiler.clearedTypeRegistry` (`compiler.go:183`, `weak.Pointer`) | diagnostics only | |
| `TypedAst` suppliers (`Compiler.typedAstFilesystem/runtimeLibraryTypedAsts`, `compiler.go:174,203`) | `func() *Node` that build a fresh node per call (`ScriptNodeDeserializer.DeserializeNew`, `typedAstDeserializer_createStubWeakScriptNode`) | no cached nodes |

### 1.10 Checked and clean

No `*Node` anywhere in `internal/colors` (ColorRegistry, Color), `AccessorSummary`,
`VariableMap`, `JsMessage`, `CompilerOptions` (`GetDefineReplacements` builds a
fresh map per call), `commentsPerFile` (`parser.Comment`), `AstAnalyzer`,
`JSChunk`/`JSChunkGraph`, `DefaultPassConfig` (factories capture the compiler),
`PassFactory`. Everything else with Node-typed fields (≈70 pass files:
`es6_normalize_classes.go` 40 fields, `replace_messages.go` 29, …,
`NodeTraversal`, `MemoizedScopeCreator` (`coalesce_variable_names.go:135`,
`node_traversal.go:1241`), `ReferenceCollector`, per-pass `GlobalNamespace`
instances, `Scope`/`Var` from `SyntacticScopeCreator`) is pass-local and dies with
its pass, so it needs no forwarding as long as compaction only runs between
passes.

## 2. Dynamic census

typed-app (971 Closure Library files, PRUNE, 408 inputs) and typed-full (1,818
files, SORT_ONLY, ES5 output, 1,944 inputs). "live" = tree + side trees
reachable from live nodes; "detached/retained" = held nodes not reachable from the
roots / the size of the subtrees they pin. Both runs byte-identical (exit 0; the
hook does not touch the AST).

### 2.1 Live nodes at the candidate points

| point | typed-app live (tree + side) | typed-full live (tree + side) |
|---|---|---|
| after pass 10 (parsed, pre-typecheck) | 380,382 (318,489 + 61,893 JSDoc type trees; 23,387 JSDocInfos) | 2,145,068 (1,986,513 + 158,555; 62,672 JSDocInfos) |
| after `checkTypes` | 381,234 (+ 40,271 JSTypes) | 2,134,654 (+ 169,757 JSTypes) |
| end of checks (= before `removeCastNodes`/`typesToColors`, the first two optimization passes) | 375,973 | 2,146,324 |
| before `typesToColors` | 374,263 (313,525 + 60,738; 54,297 of the side nodes reachable only through types) | 2,139,146 (1,985,296 + 153,850; 134,415 only through types) |
| **after `typesToColors`** | **313,527** (313,525 + 2; 8,144 simplified JSDocInfos) | **1,985,298** (+ 2; 20,240 JSDocInfos) |
| loop 1 (4 passes) iter 2 | 181,758 | 1,644,896 |
| loop 3 (12 passes) iter 3 | 167,092 | 1,579,962 |
| end of optimizations | 166,947 | 1,579,112 |
| end of finalizations | 163,841 | 1,490,591 |

`typesToColors` alone sheds 16 % (typed-app) / 7 % (typed-full) of the live
nodes (the JSDoc type-expression trees) on top of the whole type graph.

### 2.2 Holders at each point (distinct nodes; "det" = detached, "ret" = nodes pinned)

typed-app:

| holder | after checkTypes | before typesToColors | after typesToColors | loop 3 iter 3 | end of finalizations |
|---|---|---|---|---|---|
| scriptNodeByFilename / CompilerInput.ast.root | 408/409 attached | 408 | 408 | 408 | 408 |
| changeTracker.timeline | 1,955 attached | 2,323 (det 1) | 2,323 (det 6, ret 48) | 5,540 (det 6, ret 48) | 5,576 (det 6, ret 48) |
| moduleMetadata.rootNode | 409 | 409 (det 1) | 409 (det 1) | 409 (det 1) | 409 (det 1) |
| moduleMap Binding.sourceNode | 387 attached | 387 (det 254, ret 617) | 387 (det 254, ret 617) | 387 (det 387, ret 1,147) | same |
| moduleMap Export.exportNode | 265 attached | 265 (det 132, ret 194) | same | 265 (det 265, ret 724) | same |
| errorManager JSError.node | 4 | 4 (det 1) | 4 (det 1) | 4 (det 2, ret 257) | same |
| typedScopeCreator.memoized keys + scope roots / var nameNodes | 23,251 / 31,165 attached | nil | nil | nil | nil |
| topScope var nameNodes / CompilerInput.typedScope roots + names | 12,941 / 40 + 782 attached | nil | nil | nil | nil |
| registry scopedNameTable / nonNullable / closureNamespace keys | 1,190 / 13 / 6 (det 1 synthetic) | 1,190 (det 36, ret 636) / 13 (det 13) / 6 (det 6) | nil | nil | nil |
| FunctionType.source / Property.propertyNode / EnumType.source / TemplateType | 8,932 / 28,588 / 91 / 31 all attached | 8,928 (det 103, ret 5,004) / 28,587 (det 298, ret 5,295) / 91 / 31 | nil | nil | nil |
| **total detached roots / retained nodes** | 1 / 1 | 299 / 5,753 (1.5 % of live) | **152 / 666 (0.21 %)** | 286 / 1,452 (0.87 %) | 286 / 1,452 (0.89 %) |

typed-full:

| holder | after checkTypes | before typesToColors | after typesToColors | loop 1 iter 2 | loop 3 iter 3 | end of finalizations |
|---|---|---|---|---|---|---|
| scriptNodeByFilename / CompilerInput.ast.root | 1,944/1,945 | 1,944 | 1,944 | 1,944 | 1,944 | 1,944 |
| changeTracker.timeline | 5,097 attached | 6,725 (det 2, ret 16) | 6,725 (det 76, ret 629) | 25,679 (det 76) | 23,423 (det 76) | 24,246 (det 76) |
| moduleMetadata.rootNode | 1,945 | 1,945 (det 1) | same | same | same | same |
| moduleMap Binding.sourceNode | 8,139 attached | 8,139 (det 3,404, ret 21,084) | same | 8,139 (det 7,638, ret 38,004) | same | same |
| moduleMap Export.exportNode / nameNode | 5,658 / 37 attached | 5,658 (det 960, ret 3,733) / 37 | same | 5,658 (det 5,157, ret 20,506) / 37 (det 37, ret 147) | same | same |
| errorManager JSError.node | 21 | 21 (det 2) | 21 (det 2) | 31 (det 5, ret 90) | 31 (det 6, ret 385) | same |
| runtimeJsLibManager.lastInjectedLibrary | – | – | – | 1 (det 1, ret 77) | same | same |
| transpilationNamespace: Name.declaration / Ref.node / roots | – | – | – | 3,991 (det 116, ret 3,634) / 19,650 (det 584, ret 5,956) / 3 (det 1: the global sentinel) | 3,991 (det 116) / 19,650 (det 616, ret 6,252) / 3 | same |
| typedScopeCreator.memoized keys / var nameNodes | 92,900 / 238,188 attached | nil | nil | nil | nil | nil |
| topScope var nameNodes / CompilerInput.typedScope roots + names | 51,252 / 476 + 97,833 attached | nil | nil | nil | nil | nil |
| registry scopedNameTable / nonNullable / closureNamespace keys | 2,965 / 448 / 385 (det 1) | 2,965 (det 395) / 448 (det 183) / 385 (det 385, ret 552) | nil | nil | nil | nil |
| FunctionType.source / Property.propertyNode / EnumType.source / TemplateType | 26,945 / 225,314 / 408 / 32 attached | 26,945 / 225,312 (det 1,211, ret 4,097) / 408 / 32 | nil | nil | nil | nil |
| **total detached roots / retained nodes** | 1 / 1 | 3,915 / 22,027 (1.0 %) | **3,274 / 21,752 (1.1 %)** | 7,725 / 44,717 (2.7 %) | 7,752 / 45,308 (2.9 %) | 7,752 / 45,308 (3.0 %) |

Side trees after colors: 2 (typed-app) / 24 (typed-full; JSDoc written by
transpilation passes) type-expression nodes; `DECLARED_TYPE_EXPR` and
`CLOSURE_UNAWARE_SHADOW` are 0 in both tasks.

### 2.3 Walk cost (one core, M-series)

A plain pre-order walk (`GetFirstChild`/`GetNext`, count only) of the attached
tree costs 7-9 ns/node right after parsing (typed-full 2.0 M nodes: 13-18 ms)
and 11-30 ns/node later in the optimization loop when the live nodes are spread
over a churned heap (typed-full 1.64 M nodes: 50 ms at loop 1 iter 3). The
census walk with a `map[*Node]` visited set costs 100-300 ns/node (0.3-0.8 s on
typed-full), which is the reason the hooked runs are slower; it is not a
compaction cost. A copying compaction does the plain walk plus a 48-B copy and
four link rewrites per node, so a few tens of ms per compaction on typed-full is
the right expectation (the design's "~10-20 ns per live node" is the walk alone).

## 3. Recommendations

### 3.1 Safe compaction points, and the forwarding list at each

Post-colors the holder set is small and closed. At every point from "after
`typesToColors`" to the end of finalizations the complete list of `*Node`
holders outside the tree is:

1. `Compiler.externsRoot/jsRoot/externAndJsRoot`, `scriptNodeByFilename` values,
   `CompilerInput_JsAst.root` for every input (`inputsById`, chunk inputs,
   `externs`), `PhaseOptimizer.jsRoot` — always attached; N+3 pointers.
2. `ChangeTracker.changeTimeline` — 2-26k entries (`eventsByValue` keys and the
   event chunk values; the map must be rebuilt with the forwarded keys). Up to 76
   entries point at detached functions that were never `ReportFunctionDeleted`
   (629 nodes pinned in typed-full).
3. `ModuleMetadataMap` root nodes (N, 1 detached) and `ModuleMap`
   `Binding.sourceNode` / `Export.exportNode/nameNode` (8.1k + 5.7k in typed-full,
   of which 3.4k-7.6k are detached after module rewriting, pinning 21-38k
   nodes).
4. `JSError.node` in the error manager (≤ 31; ≤ 6 detached).
5. `RuntimeJsLibManager.lastInjectedLibrary` (1; detached from loop 1 on in
   typed-full, 77 nodes; it is read again by the next injection).
6. `Compiler.transpilationNamespace` when built (transpiling builds): 3 roots,
   `Name.declaration` (4k), `Ref.node` (20k; also `Name.refsIndex` map keys),
   `spreadSiblingCache` and `defineCandidateScripts` keys, and `Ref.scope` →
   `Scope.rootNode`; 0.6-0.7k refs detached (6k nodes pinned).
7. `PhaseOptimizer_Loop.Process` locals `scriptSizes` (map keyed by SCRIPT
   nodes) when compacting inside a loop; rebuild the map or key it by input.
8. Node-owned side trees (JSDoc type expressions: 2-24 nodes after colors;
   `DECLARED_TYPE_EXPR`, `CLOSURE_UNAWARE_SHADOW`: 0) — copied as part of the
   owning node, not forwarded.
9. Go-only: `cachedCCAArena` (evict, as `evictPassArenaCaches` already does at
   stage ends), the package-level sentinel nodes (keep them in a process-global,
   never-compacted chunk; they are never attached), `SourceFile.sourceInfoTemplateVal`
   (a by-value Node outside any chunk: give it a reserved id or make it a plain
   source-position struct).

So the safe points are: **right after `typesToColors`** (the largest payoff:
16 % / 7 % of live nodes plus all checks-stage garbage become reclaimable at
once; the forwarding list is items 1-5 only, since the transpilation namespace is
not built yet and no loop is running), **the end of every optimization-loop
iteration or every N passes between loops** (items 1-7), and **the end of the
optimization stage** (items 1-6; the finalization stage frees too little to be
worth a copy). Item counts: ≈ 2.5k (typed-app) to ≈ 60k (typed-full) pointers to
rewrite per compaction, all enumerable from the Compiler, versus 1.6-2.0 M nodes
copied.

The checks → optimizations boundary is **not** a useful point by itself: the
first two optimization passes are `removeCastNodes` and `typesToColors`
(`default_pass_config.go:2682, 2982-2984`), so "end of checks" has the full type
graph alive (section 3.3).

### 3.2 Holders to nil before compacting, and Java parity

- `typedScopeCreator`, `topScope`, `CompilerInput.typedScope`: already nil at
  every post-colors point because `clearTypedScopeCreatorPass` and
  `clearTopTypedScopePass` run in checks (Java does the same). The exception is
  `options.PreservesDetailedSourceInfo()` (IDE embedding; Java keeps the scopes
  too): in that mode either disable compaction or add the 93k-scope / 238k-name
  tables to the forwarding list.
- `typeRegistry`, `typeValidator` (mismatch locations), `abstractInterpreter`:
  nil from `typesToColors` on (Java same).
- `transpilationNamespace`: Java never clears it. Nilling it before a compaction
  would make `GetTranspilationNamespace` rebuild it from the current AST, which
  can differ from the stale one Java consults (its names come from externs and
  the first script only, so the difference is small but not provably empty);
  forward it instead (item 6) to stay byte-identical.
- `cachedCCAArena`: Go-only, stale; evict.
- Nothing else needs nilling; every other holder is either forwarded (3.1) or
  pass-local.

### 3.3 Hard holders

- **Detached nodes that Java also keeps**: `ModuleMap` bindings/exports
  (3.4k-7.6k detached, 21-38k nodes), `changeTimeline` entries (≤ 76, 629
  nodes), `JSError.node` (≤ 6), `lastInjectedLibrary` (1, and read later),
  `transpilationNamespace` refs (≈ 0.7k, 6k nodes). Forwarding them means the
  compaction must also copy their detached subtrees (they are extra roots: the
  distinct detached roots are 152-286 / 3.3k-7.8k, 0.2-3 % of live nodes), or
  the holder must be proven unread. `ModuleMap` is read in the optimization
  stage (`GlobalNamespace.SetupModuleMetadata`, polymer) but only for MODULE_BODY
  roots, which no longer exist after module rewriting, so dropping it after the
  last checks-stage reader would be semantically safe in Java terms; the
  byte-identical-first policy says copy the 0.2-3 % instead and decide later.
  `lastInjectedLibrary` must be copied: a later `$jscomp` injection walks from it.
- **Pre-colors compaction is impractical**: before `typesToColors` the type graph
  holds 27k function sources + 225k property nodes + 93k scope roots + 238k var
  name nodes (typed-full), i.e. ≈ 600k pointers scattered through
  `FunctionType`, `Property`, `TypedScope`/`TypedVar` and three registry maps,
  plus 134k JSDoc type-tree nodes reachable only through types. That is the
  case the design already avoids by dropping the type column wholesale at
  colors; nothing before `typesToColors` should compact.
- **IDE mode** (`PreservesDetailedSourceInfo`, `SymbolTable`,
  `PreprocessorSymbolTable`): large Java-identical holders (hundreds of
  thousands of node references); out of scope for the CLI, but compaction must
  be off or fully forwarded there.
- **`Timeline` and `GlobalNamespace_Name.refsIndex` are maps keyed by `*Node`**:
  forwarding is a rebuild, not an in-place rewrite (hash by pointer). With
  node ids instead of pointers as keys (round-2 "dense node-id-indexed side
  tables") the rebuild disappears, because compaction can assign ids in
  pre-order and remap through the forwarding table.
