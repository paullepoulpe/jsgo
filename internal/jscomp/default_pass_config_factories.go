// Port of com/google/javascript/jscomp/DefaultPassConfig.java @ closure-compiler v20261021

/*
 * Copyright 2009 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

var errPropsNotStaticallyAnalyzable = NewPassFactory_PreconditionResult(
	false,
	"requires assumePropertiesAreStaticallyAnalyzable to be enabled",
)

func defaultPassConfig_requirePropertiesAreStaticallyAnalyzable(
	options *CompilerOptions,
) PassFactory_PreconditionResult {
	if options.GetAssumePropertiesAreStaticallyAnalyzable() {
		return PassFactory_PreconditionResult_SUCCESS
	}
	return errPropsNotStaticallyAnalyzable
}

func inlineAndCollapsePrecondition(options *CompilerOptions) PassFactory_PreconditionResult {
	if options.GetPropertyCollapseLevel() == CompilerOptions_PropertyCollapseLevel_NONE ||
		options.GetAssumePropertiesAreStaticallyAnalyzable() {
		return PassFactory_PreconditionResult_SUCCESS
	}
	return errPropsNotStaticallyAnalyzable
}

// Stateless PassFactory singletons shared across all DefaultPassConfig instances.
var (
	defaultPass_stripSideEffectProtection = makePassFactory(PassNames_STRIP_SIDE_EFFECT_PROTECTION, func(compiler AbstractCompiler) CompilerPass {
		return NewCheckSideEffects_StripProtection(compiler)
	})
	defaultPass_checkMissingRequires = makePassFactory("checkMissingRequires", func(compiler AbstractCompiler) CompilerPass {
		return NewCheckMissingRequires(compiler, compiler.GetModuleMetadataMap())
	})
	defaultPass_checkJsDocAndEs6Modules = makePassFactory("checkJsDocAndEs6Modules", func(compiler AbstractCompiler) CompilerPass {
		return defaultPassConfig_combineChecks(compiler, []NodeTraversal_Callback{
			NewCheckJSDoc(compiler),
			NewEs6CheckModule(compiler),
		})
	})
	defaultPass_generateIjs = makePassFactory("generateIjs", func(compiler AbstractCompiler) CompilerPass {
		return NewConvertToTypedInterface(compiler)
	})
	defaultPass_removeExtraRequires = makePassFactory("removeExtraRequires", func(compiler AbstractCompiler) CompilerPass {
		return NewExtraRequireRemover(compiler)
	})
	defaultPass_exportTestFunctions = makePassFactory(PassNames_EXPORT_TEST_FUNCTIONS, func(compiler AbstractCompiler) CompilerPass {
		convention := compiler.GetCodingConvention()
		if convention.GetExportSymbolFunction() != nil {
			return NewExportTestFunctions(
				compiler,
				convention.GetExportSymbolFunction(),
				convention.GetExportPropertyFunction())
		}
		return NewErrorPass(compiler, defaultPassConfig_GENERATE_EXPORTS_ERROR)
	})
	defaultPass_gatherRawExports = makePassFactory(PassNames_GATHER_RAW_EXPORTS, func(compiler AbstractCompiler) CompilerPass {
		pass := NewGatherRawExports(compiler)
		return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
			pass.Process(externs, root)
			compiler.AddExportedNames(pass.GetExportedVariableNames())
		})
	})
	defaultPass_angularPass = makePassFactory(PassNames_ANGULAR_PASS, func(compiler AbstractCompiler) CompilerPass {
		return NewAngularPass(compiler)
	})
	defaultPass_closureGoogScopeAliasesForIjs = makePassFactory("closureGoogScopeAliasesForIjs", func(compiler AbstractCompiler) CompilerPass {
		return ScopedAliases_Builder_(compiler).Build()
	})
	defaultPass_injectRuntimeLibraries = makePassFactory("InjectRuntimeLibraries", func(compiler AbstractCompiler) CompilerPass {
		return NewInjectRuntimeLibraries(
			compiler,
			jutil.NewLinkedHashSet(compiler.GetOptions().GetForceLibraryInjectionList()...))
	})
	defaultPass_removeWeakSources = makePassFactory("removeWeakSources", func(compiler AbstractCompiler) CompilerPass {
		return NewRemoveWeakSources(compiler)
	})
	defaultPass_declaredGlobalExternsOnWindow = makePassFactory(PassNames_DECLARED_GLOBAL_EXTERNS_ON_WINDOW, func(compiler AbstractCompiler) CompilerPass {
		return NewDeclaredGlobalExternsOnWindow(compiler)
	})
	defaultPass_checkTypeImportCodeReferences = makePassFactory("checkTypeImportCodeReferences", func(compiler AbstractCompiler) CompilerPass {
		return NewCheckTypeImportCodeReferences(compiler)
	})
	defaultPass_closureCheckModule = makePassFactory("closureCheckModule", func(compiler AbstractCompiler) CompilerPass {
		return NewClosureCheckModule(compiler, compiler.GetModuleMetadataMap())
	})
	defaultPass_checkClosureImports = makePassFactory("checkGoogRequires", func(compiler AbstractCompiler) CompilerPass {
		return NewCheckClosureImports(compiler, compiler.GetModuleMetadataMap())
	})
	defaultPass_rewriteGoogJsImports = makePassFactory("rewriteGoogJsImports", func(compiler AbstractCompiler) CompilerPass {
		return NewRewriteGoogJsImports(compiler, RewriteGoogJsImports_Mode_LINT_AND_REWRITE, compiler.GetModuleMap())
	})
	defaultPass_earlyPeepholeOptimizations = makePassFactory("earlyPeepholeOptimizations", func(compiler AbstractCompiler) CompilerPass {
		opts := compiler.GetOptions()
		useTypes := opts.ShouldUseTypesForLocalOptimization()
		peepholes := []AbstractPeepholeOptimization{NewPeepholeRemoveDeadCode()}
		if opts.GetJ2clPass().ShouldAddJ2clPasses() {
			peepholes = append(peepholes, NewJ2clEqualitySameRewriterPass(useTypes))
		}
		return NewPeepholeOptimizationsPassOfList(compiler, "earlyPeepholeOptimizations", peepholes)
	})
	defaultPass_peepholeOptimizations = makeLoopPassFactory(PassNames_PEEPHOLE_OPTIMIZATIONS, func(compiler AbstractCompiler) CompilerPass {
		return defaultPassConfig_createPeepholeOptimizationsPass(compiler, PassNames_PEEPHOLE_OPTIMIZATIONS, true)
	})
	defaultPass_peepholeOptimizationsOnceNormalized = makePassFactory(PassNames_PEEPHOLE_OPTIMIZATIONS, func(compiler AbstractCompiler) CompilerPass {
		return defaultPassConfig_createPeepholeOptimizationsPass(compiler, PassNames_PEEPHOLE_OPTIMIZATIONS, true)
	})
	defaultPass_peepholeOptimizationsOnceNonNormalized = makePassFactory(PassNames_PEEPHOLE_OPTIMIZATIONS, func(compiler AbstractCompiler) CompilerPass {
		return defaultPassConfig_createPeepholeOptimizationsPass(compiler, PassNames_PEEPHOLE_OPTIMIZATIONS, false)
	})
	defaultPass_checkVars = makePassFactory(PassNames_CHECK_VARS, func(compiler AbstractCompiler) CompilerPass {
		return NewVarCheck(compiler)
	})
	defaultPass_inferConsts = makePassFactory(PassNames_INFER_CONSTS, func(compiler AbstractCompiler) CompilerPass {
		return NewInferConsts(compiler)
	})
	defaultPass_checkVariableReferences = makePassFactory(PassNames_CHECK_VARIABLE_REFERENCES, func(compiler AbstractCompiler) CompilerPass {
		return NewVariableReferenceCheck(compiler)
	})
	defaultPass_checkSuper = makePassFactory("checkSuper", func(compiler AbstractCompiler) CompilerPass {
		return NewCheckSuper(compiler)
	})
	defaultPass_clearTypedScopeCreatorPass = makePassFactory("clearTypedScopeCreatorPass", func(compiler AbstractCompiler) CompilerPass {
		return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
			compiler.ClearTypedScopeCreator()
		})
	})
	defaultPass_clearTopTypedScopePass = makePassFactory("clearTopTypedScopePass", func(compiler AbstractCompiler) CompilerPass {
		return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
			compiler.SetTopScope(nil)
			for _, compilerInput := range compiler.GetInputsInOrder() {
				compilerInput.SetTypedScope(nil)
			}
		})
	})
	defaultPass_inferTypes = makePassFactory(PassNames_INFER_TYPES, func(compiler AbstractCompiler) CompilerPass {
		return CompilerPassFunc(func(unused rhino.Node, srcRoot rhino.Node) {
			globalRoot := srcRoot.GetParent()
			inferencePass := NewTypeInferencePass(
				compiler,
				compiler.GetReverseAbstractInterpreter(),
				compiler.GetTypedScopeCreator().(*TypedScopeCreator))
			compiler.SetTypeCheckingHasRun(true)
			compiler.SetTopScope(inferencePass.InferAllScopes(globalRoot))
		})
	})
	defaultPass_inferJsDocInfo = makePassFactory("inferJsDocInfo", func(compiler AbstractCompiler) CompilerPass {
		return NewInferJSDocInfo(compiler)
	})
	defaultPass_checkAccessControls = makePassFactory("checkAccessControls", func(compiler AbstractCompiler) CompilerPass {
		return NewCheckAccessControls(compiler)
	})
	defaultPass_lintChecks = makePassFactory(PassNames_LINT_CHECKS, func(compiler AbstractCompiler) CompilerPass {
		return defaultPassConfig_combineChecks(compiler, []NodeTraversal_Callback{
			NewCheckConstPrivateProperties(compiler),
			NewCheckConstantCaseNames(compiler),
			NewCheckEmptyStatements(compiler),
			NewCheckEnums(compiler),
			NewCheckEs6ModuleFileStructure(compiler),
			NewCheckEs6Modules(compiler),
			NewCheckNoMutatedEs6Exports(compiler),
			NewCheckGoogModuleTypeScriptName(compiler),
			NewCheckInterfaces(compiler),
			NewCheckJSDocStyle(compiler),
			NewCheckMissingSemicolon(compiler),
			NewCheckNullabilityModifiers(compiler),
			NewCheckPrimitiveAsObject(compiler),
			NewCheckPrototypeProperties(compiler),
			NewCheckUnusedPrivateProperties(compiler),
			NewCheckUnusedLabels(compiler),
			NewCheckUselessBlocks(compiler),
			NewCheckVar(compiler),
		})
	})
	defaultPass_analyzerChecks = makePassFactory(PassNames_ANALYZER_CHECKS, func(compiler AbstractCompiler) CompilerPass {
		return defaultPassConfig_combineChecks(compiler, []NodeTraversal_Callback{
			NewCheckArrayWithGoogObject(compiler),
			NewImplicitNullabilityCheck(compiler),
			NewCheckNestedNames(compiler),
		})
	})
	defaultPass_checkRequiresAndProvidesSorted = makePassFactory("checkRequiresAndProvidesSorted", func(compiler AbstractCompiler) CompilerPass {
		return defaultPassConfig_combineChecks(compiler, []NodeTraversal_Callback{
			NewCheckProvidesSorted(CheckProvidesSorted_Mode_COLLECT_AND_REPORT),
			NewCheckRequiresSorted(CheckRequiresSorted_Mode_COLLECT_AND_REPORT),
		})
	})
	defaultPass_stripCode = makePassFactory("stripCode", func(compiler AbstractCompiler) CompilerPass {
		opts := compiler.GetOptions()
		return NewStripCode(
			compiler,
			opts.GetStripTypes(),
			opts.GetStripNameSuffixes(),
			opts.GetStripNamePrefixes(),
			opts.GetTweakProcessing().ShouldStrip())
	})
	defaultPass_checkConsts = makePassFactory("checkConsts", func(compiler AbstractCompiler) CompilerPass {
		return NewConstCheck(compiler)
	})
	defaultPass_rewriteCallerCodeLocation = makePassFactory("rewriteCallerCodeLocation", func(compiler AbstractCompiler) CompilerPass {
		return NewRewriteCallerCodeLocation(compiler)
	})
	defaultPass_replaceToggles = makePassFactory("replaceToggles", func(compiler AbstractCompiler) CompilerPass {
		return NewReplaceToggles(compiler)
	})
	defaultPass_closureOptimizePrimitives = makePassFactory("closureOptimizePrimitives", func(compiler AbstractCompiler) CompilerPass {
		return NewClosureOptimizePrimitives(
			compiler, compiler.GetOptions().GetOutputFeatureSet().Contains(parser.FeatureSet_ES2015))
	})
	defaultPass_convertChunksToESModules = makePassFactory("convertChunksToESModules", func(compiler AbstractCompiler) CompilerPass {
		return NewConvertChunksToESModules(compiler)
	})
	defaultPass_collapseObjectLiterals = makeLoopPassFactory(PassNames_COLLAPSE_OBJECT_LITERALS, func(compiler AbstractCompiler) CompilerPass {
		return NewInlineObjectLiterals(compiler, compiler.GetUniqueNameIdSupplier())
	})
	defaultPass_devirtualizeMethods = makeCheckedPassFactory(
		PassNames_DEVIRTUALIZE_METHODS,
		defaultPassConfig_requirePropertiesAreStaticallyAnalyzable,
		false,
		func(compiler AbstractCompiler) CompilerPass {
			return OptimizeCalls_Builder_().
				SetCompiler(compiler).
				SetConsiderExterns(false).
				AddPass(NewDevirtualizeMethods(compiler)).
				Build()
		},
	)
	defaultPass_optimizeCalls = makeCheckedPassFactory(
		PassNames_OPTIMIZE_CALLS,
		defaultPassConfig_requirePropertiesAreStaticallyAnalyzable,
		true,
		func(compiler AbstractCompiler) CompilerPass {
			return OptimizeCalls_Builder_().
				SetCompiler(compiler).
				SetConsiderExterns(false).
				AddPass(NewOptimizeReturns(compiler).AsCallGraphCompilerPass()).
				AddPass(NewOptimizeParameters(compiler).AsCallGraphCompilerPass()).
				Build()
		},
	)
	defaultPass_optimizeConstructors = makePassFactory("optimizeConstructors", func(compiler AbstractCompiler) CompilerPass {
		return OptimizeCalls_Builder_().
			SetCompiler(compiler).
			SetConsiderExterns(false).
			AddPass(NewOptimizeConstructors(compiler).AsCallGraphCompilerPass()).
			Build()
	})
	defaultPass_markPureFunctions = makeCheckedPassFactory(
		"markPureFunctions",
		defaultPassConfig_requirePropertiesAreStaticallyAnalyzable,
		false,
		func(compiler AbstractCompiler) CompilerPass {
			return NewPureFunctionIdentifier_Driver(compiler)
		},
	)
	defaultPass_inlineConstants = makeLoopPassFactory("inlineConstants", func(compiler AbstractCompiler) CompilerPass {
		return NewInlineVariables(compiler, InlineVariables_Mode_CONSTANTS_ONLY)
	})
	defaultPass_inlineSimpleMethods = makeCheckedPassFactory(
		"inlineSimpleMethods",
		defaultPassConfig_requirePropertiesAreStaticallyAnalyzable,
		true,
		func(compiler AbstractCompiler) CompilerPass {
			return NewInlineSimpleMethods(compiler)
		},
	)
	defaultPass_deadAssignmentsElimination = makeLoopPassFactory(PassNames_DEAD_ASSIGNMENT_ELIMINATION, func(compiler AbstractCompiler) CompilerPass {
		return NewDeadAssignmentsElimination(compiler)
	})
	defaultPass_deadPropertyAssignmentElimination = makeCheckedPassFactory(
		"deadPropertyAssignmentElimination",
		defaultPassConfig_requirePropertiesAreStaticallyAnalyzable,
		true,
		func(compiler AbstractCompiler) CompilerPass {
			return NewDeadPropertyAssignmentElimination(compiler)
		},
	)
	defaultPass_inlineProperties = makeLoopPassFactory(PassNames_INLINE_PROPERTIES, func(compiler AbstractCompiler) CompilerPass {
		return NewInlineProperties(compiler)
	})
	defaultPass_isolatePolyfills = makePassFactory("IsolatePolyfills", func(compiler AbstractCompiler) CompilerPass {
		return NewIsolatePolyfills(compiler)
	})
	defaultPass_flowSensitiveInlineVariables = makePassFactory(PassNames_FLOW_SENSITIVE_INLINE_VARIABLES, func(compiler AbstractCompiler) CompilerPass {
		return NewFlowSensitiveInlineVariables(compiler)
	})
	defaultPass_collapseVariableDeclarations = makePassFactory(PassNames_COLLAPSE_VARIABLE_DECLARATIONS, func(compiler AbstractCompiler) CompilerPass {
		return NewCollapseVariableDeclarations(compiler)
	})
	defaultPass_rewriteFunctionExpressions = makePassFactory(PassNames_REWRITE_FUNCTION_EXPRESSIONS, func(compiler AbstractCompiler) CompilerPass {
		return NewFunctionRewriter(compiler)
	})
	defaultPass_collapseAnonymousFunctions = makePassFactory(PassNames_COLLAPSE_ANONYMOUS_FUNCTIONS, func(compiler AbstractCompiler) CompilerPass {
		return NewCollapseAnonymousFunctions(compiler)
	})
	defaultPass_rewriteGlobalDeclarationsForTryCatchWrapping = makePassFactory("rewriteGlobalDeclarationsForTryCatchWrapping", func(compiler AbstractCompiler) CompilerPass {
		return NewRewriteGlobalDeclarationsForTryCatchWrapping(compiler)
	})
	defaultPass_markUnnormalized = makePassFactory("markUnnormalized", func(compiler AbstractCompiler) CompilerPass {
		return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
			compiler.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_RAW)
		})
	})
	defaultPass_normalize = makePassFactory(PassNames_NORMALIZE, func(compiler AbstractCompiler) CompilerPass {
		return Normalize_CreateNormalizeForOptimizations(compiler)
	})
	defaultPass_externExports = makePassFactory(PassNames_EXTERN_EXPORTS, func(compiler AbstractCompiler) CompilerPass {
		return NewExternExportsPass(compiler)
	})
	defaultPass_invertContextualRenaming = makePassFactory("invertContextualRenaming", func(compiler AbstractCompiler) CompilerPass {
		return MakeDeclaredNamesUnique_GetContextualRenameInverter(compiler)
	})
	defaultPass_removePropertyRenamingCalls = makePassFactory("removePropertyRenamingCalls", func(compiler AbstractCompiler) CompilerPass {
		return NewRemovePropertyRenamingCalls(compiler)
	})
	defaultPass_renameLabels = makePassFactory("renameLabels", func(compiler AbstractCompiler) CompilerPass {
		return NewRenameLabels(compiler)
	})
	defaultPass_convertToDottedProperties = makePassFactory(PassNames_CONVERT_TO_DOTTED_PROPERTIES, func(compiler AbstractCompiler) CompilerPass {
		return NewConvertToDottedProperties(compiler)
	})
	defaultPass_checkAstValidity = makePassFactory("checkAstValidity", func(compiler AbstractCompiler) CompilerPass {
		return NewAstValidator(compiler)
	})
	defaultPass_varCheckValidity = makePassFactory("varCheckValidity", func(compiler AbstractCompiler) CompilerPass {
		return NewVarCheckWithValidityCheck(compiler, true)
	})
	defaultPass_gatherExternPropertiesCheck = makePassFactory("gatherExternProperties", func(compiler AbstractCompiler) CompilerPass {
		return NewGatherExternProperties(compiler, GatherExternProperties_Mode_CHECK)
	})
	defaultPass_gatherExternPropertiesOptimize = makePassFactory("gatherExternProperties", func(compiler AbstractCompiler) CompilerPass {
		return NewGatherExternProperties(compiler, GatherExternProperties_Mode_OPTIMIZE)
	})
	defaultPass_polymerPass = makePassFactory("polymerPass", func(compiler AbstractCompiler) CompilerPass {
		return NewPolymerPass(compiler)
	})
	defaultPass_chromePass = makePassFactory("chromePass", func(compiler AbstractCompiler) CompilerPass {
		return NewChromePass(compiler)
	})
	defaultPass_j2clConstantHoisterPass = makeLoopPassFactory("j2clConstantHoisterPass", func(compiler AbstractCompiler) CompilerPass {
		return NewJ2clConstantHoisterPass(compiler)
	})
	defaultPass_j2clClinitPass = makeLoopPassFactory("j2clClinitPass", func(compiler AbstractCompiler) CompilerPass {
		return NewJ2clClinitPrunerPass(
			compiler,
			compiler.GetChangeTracker().GetChangedScopeNodesForPass("j2clClinitPass"))
	})
	defaultPass_j2clPropertyInlinerPass = makePassFactory("j2clPropertyInlinerPass", func(compiler AbstractCompiler) CompilerPass {
		return NewJ2clPropertyInlinerPass(compiler)
	})
	defaultPass_j2clPass = makePassFactory("j2clPass", func(compiler AbstractCompiler) CompilerPass {
		return NewJ2clPass(compiler)
	})
	defaultPass_j2clUtilGetDefineRewriterPass = makePassFactory("j2clUtilGetDefineRewriterPass", func(compiler AbstractCompiler) CompilerPass {
		return NewJ2clUtilGetDefineRewriterPass(compiler)
	})
	defaultPass_j2clAssertRemovalPass = makePassFactory("j2clAssertRemovalPass", func(compiler AbstractCompiler) CompilerPass {
		return NewJ2clAssertRemovalPass(compiler)
	})
	defaultPass_j2clSourceFileChecker = makePassFactory("j2clSourceFileChecker", func(compiler AbstractCompiler) CompilerPass {
		return NewJ2clSourceFileChecker(compiler)
	})
	defaultPass_j2clChecksPass = makePassFactory("j2clChecksPass", func(compiler AbstractCompiler) CompilerPass {
		return NewJ2clChecksPass(compiler)
	})
	defaultPass_removeCastNodes = makePassFactory("removeCastNodes", func(compiler AbstractCompiler) CompilerPass {
		return NewRemoveCastNodes(compiler)
	})
	defaultPass_typesToColors = makePassFactory("typesToColors", func(compiler AbstractCompiler) CompilerPass {
		return CompilerPassFunc(func(externs rhino.Node, js rhino.Node) {
			NewConvertTypesToColors(
				compiler,
				SerializationOptions_Builder_().
					SetIncludeDebugInfo(compiler.GetOptions().ShouldSerializeExtraDebugInfo()).
					Build()).
				Process(externs, js)
			compiler.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_COLORS_AND_SIMPLIFIED_JSDOC)
		})
	})
	defaultPass_removeUnnecessarySyntheticExterns = makePassFactory("removeUnnecessarySyntheticExterns", func(compiler AbstractCompiler) CompilerPass {
		return NewRemoveUnnecessarySyntheticExterns(compiler)
	})
	defaultPass_optimizeToEs6 = makePassFactory("optimizeToEs6", func(compiler AbstractCompiler) CompilerPass {
		return NewSubstituteEs6Syntax(compiler)
	})
	defaultPass_whitespaceWrapGoogModules = makePassFactory("whitespaceWrapGoogModules", func(compiler AbstractCompiler) CompilerPass {
		return NewWhitespaceWrapGoogModules(compiler)
	})
	defaultPass_rewriteCommonJsModules = makePassFactory(PassNames_REWRITE_COMMON_JS_MODULES, func(compiler AbstractCompiler) CompilerPass {
		return NewProcessCommonJSModules(compiler)
	})
	defaultPass_rewriteScriptsToEs6Modules = makePassFactory(PassNames_REWRITE_SCRIPTS_TO_ES6_MODULES, func(compiler AbstractCompiler) CompilerPass {
		return NewEs6RewriteScriptsToModules(compiler)
	})
	defaultPass_createModuleMapPass = makePassFactory(PassNames_CREATE_MODULE_MAP, func(compiler AbstractCompiler) CompilerPass {
		return NewModuleMapCreator(compiler, compiler.GetModuleMetadataMap())
	})
	defaultPass_gatherGettersAndSetters = makePassFactory(PassNames_GATHER_GETTERS_AND_SETTERS, func(compiler AbstractCompiler) CompilerPass {
		return NewGatherGetterAndSetterProperties(compiler)
	})
	defaultPass_addSyntheticScript = makePassFactory("ADD_SYNTHETIC_SCRIPT", func(compiler AbstractCompiler) CompilerPass {
		return CompilerPassFunc(func(externs rhino.Node, js rhino.Node) {
			compiler.InitializeSyntheticCodeInput()
		})
	})
	defaultPass_removeSyntheticScript = makePassFactory("REMOVE_SYNTHETIC_SCRIPT", func(compiler AbstractCompiler) CompilerPass {
		return CompilerPassFunc(func(externs rhino.Node, js rhino.Node) {
			compiler.RemoveSyntheticCodeInput()
		})
	})
	defaultPass_mergeSyntheticScript = makePassFactory("MERGE_SYNTHETIC_SCRIPT", func(compiler AbstractCompiler) CompilerPass {
		return CompilerPassFunc(func(externs rhino.Node, js rhino.Node) {
			compiler.MergeSyntheticCodeInput()
		})
	})
	defaultPass_protectLocaleData = makePassFactory("protectLocaleData", func(compiler AbstractCompiler) CompilerPass {
		return NewLocaleDataPasses_ProtectGoogLocale(compiler)
	})
	defaultPass_substituteLocaleData = makePassFactory("SubstituteLocaleData", func(compiler AbstractCompiler) CompilerPass {
		return NewLocaleDataPasses_LocaleSubstitutions(compiler, compiler.GetOptions().GetLocale())
	})
	defaultPass_transpileOnlyClosureUnaware = makePassFactory("TranspileOnlyClosureUnaware", func(compiler AbstractCompiler) CompilerPass {
		return NewTranspileAndOptimizeClosureUnaware(compiler, NestedCompilerRunner_Mode_TRANSPILE_ONLY)
	})
	defaultPass_transpileAndOptimizeClosureUnaware = makePassFactory("TranspileAndOptimizeClosureUnaware", func(compiler AbstractCompiler) CompilerPass {
		return NewTranspileAndOptimizeClosureUnaware(compiler, NestedCompilerRunner_Mode_TRANSPILE_AND_OPTIMIZE)
	})

	defaultPassConfig_forbidDynamicImportUsage = makePassFactory("FORBID_DYNAMIC_IMPORT", func(compiler AbstractCompiler) CompilerPass {
		return NewForbidDynamicImportUsage(compiler)
	})
	defaultPassConfig_rewriteDynamicImports = makePassFactory("REWRITE_DYNAMIC_IMPORT", func(compiler AbstractCompiler) CompilerPass {
		opts := compiler.GetOptions()
		return NewRewriteDynamicImports(compiler, opts.GetDynamicImportAlias(), opts.GetChunkOutputType())
	})
)

func (c *DefaultPassConfig) initFields() {
	options := c.options

	c.preprocessorSymbolTableFactory = NewPreprocessorSymbolTable_CachedInstanceFactory()

	c.checkSideEffects = makePassFactory("checkSideEffects", func(compiler AbstractCompiler) CompilerPass {
		return NewCheckSideEffects(
			compiler,
			options.GetCheckSuspiciousCode(),
			options.ShouldProtectHiddenSideEffects())
	})
	c.stripSideEffectProtection = defaultPass_stripSideEffectProtection

	c.suspiciousCode = makePassFactory("suspiciousCode", func(compiler AbstractCompiler) CompilerPass {
		var callbacks []NodeTraversal_Callback
		if options.GetCheckSuspiciousCode() {
			callbacks = append(callbacks, NewCheckSuspiciousCode(), NewCheckDuplicateCase(compiler))
		}
		if options.Enables(DiagnosticGroups_GLOBAL_THIS) {
			callbacks = append(callbacks, NewCheckGlobalThis(compiler))
		}
		if options.Enables(DiagnosticGroups_DEBUGGER_STATEMENT_PRESENT) {
			callbacks = append(callbacks, NewCheckDebuggerStatement(compiler))
		}
		return defaultPassConfig_combineChecks(compiler, callbacks)
	})

	c.extraRequires = makePassFactory("checkExtraRequires", func(compiler AbstractCompiler) CompilerPass {
		return NewCheckExtraRequires(compiler, options.GetUnusedImportsToRemove())
	})
	c.checkMissingRequires = defaultPass_checkMissingRequires
	c.checkJsDocAndEs6Modules = defaultPass_checkJsDocAndEs6Modules

	c.generateExports = makePassFactory(PassNames_GENERATE_EXPORTS, func(compiler AbstractCompiler) CompilerPass {
		convention := compiler.GetCodingConvention()
		pass := NewGenerateExports(
			compiler,
			options.ShouldExportLocalPropertyDefinitions(),
			convention.GetExportSymbolFunction(),
			convention.GetExportPropertyFunction())
		return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
			pass.Process(externs, root)
			compiler.AddExportedNames(pass.GetExportedVariableNames())
		})
	})

	c.generateIjs = defaultPass_generateIjs
	c.removeExtraRequires = defaultPass_removeExtraRequires
	c.exportTestFunctions = defaultPass_exportTestFunctions
	c.gatherRawExports = defaultPass_gatherRawExports

	c.closurePrimitives = makePassFactory("closurePrimitives", func(compiler AbstractCompiler) CompilerPass {
		c.preprocessorSymbolTableFactory.MaybeInitialize(compiler)
		return NewProcessClosurePrimitives(compiler)
	})

	c.closureProvidesRequires = makePassFactory("closureProvidesRequires", func(compiler AbstractCompiler) CompilerPass {
		c.preprocessorSymbolTableFactory.MaybeInitialize(compiler)
		pass := NewProcessClosureProvidesAndRequires(compiler, options.ShouldPreservesGoogProvidesAndRequires())
		return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
			pass.Process(externs, root)
			compiler.AddExportedNames(pass.GetExportedVariableNames())
		})
	})

	c.angularPass = defaultPass_angularPass

	c.replaceMessagesForChrome = makePassFactory(PassNames_REPLACE_MESSAGES, func(compiler AbstractCompiler) CompilerPass {
		return NewReplaceMessagesForChrome(compiler, NewGoogleJsMessageIdGenerator(options.GetTcProjectId()))
	})

	c.closureGoogScopeAliasesForIjs = defaultPass_closureGoogScopeAliasesForIjs

	c.closureGoogScopeAliases = makePassFactory("closureGoogScopeAliases", func(compiler AbstractCompiler) CompilerPass {
		c.preprocessorSymbolTableFactory.MaybeInitialize(compiler)
		return ScopedAliases_Builder_(compiler).
			SetPreprocessorSymbolTable(c.preprocessorSymbolTableFactory.GetInstanceOrNull()).
			SetModuleMetadataMap(compiler.GetModuleMetadataMap()).
			SetInvalidModuleGetHandling(ScopedAliases_InvalidModuleGetHandling_GIVE_UNIQUE_NAME).
			Build()
	})

	c.injectRuntimeLibraries = defaultPass_injectRuntimeLibraries
	c.removeWeakSources = defaultPass_removeWeakSources
	c.declaredGlobalExternsOnWindow = defaultPass_declaredGlobalExternsOnWindow
	c.checkTypeImportCodeReferences = defaultPass_checkTypeImportCodeReferences
	c.closureCheckModule = defaultPass_closureCheckModule

	c.closureRewriteModule = makePassFactory("closureRewriteModule", func(compiler AbstractCompiler) CompilerPass {
		c.preprocessorSymbolTableFactory.MaybeInitialize(compiler)
		return NewClosureRewriteModule(
			compiler,
			c.preprocessorSymbolTableFactory.GetInstanceOrNull(),
			compiler.GetTopScope())
	})

	c.checkClosureImports = defaultPass_checkClosureImports
	c.rewriteGoogJsImports = defaultPass_rewriteGoogJsImports

	c.closureReplaceGetCssName = makePassFactory("closureReplaceGetCssName", func(compiler AbstractCompiler) CompilerPass {
		return CompilerPassFunc(func(externs rhino.Node, jsRoot rhino.Node) {
			var cssNames *jutil.LinkedHashSet[string]
			if options.ShouldGatherCssNames() {
				cssNames = jutil.NewLinkedHashSet[string]()
			}
			pass := NewReplaceCssNames(
				compiler,
				options.GetCssRenamingMap(),
				ReplaceCssNames_CssNameCollectorFunc(func(cssName string) {
					if cssNames != nil {
						cssNames.Add(cssName)
					}
				}),
				options.GetCssRenamingSkiplist())
			pass.Process(externs, jsRoot)
			compiler.SetCssNames(cssNames)
		})
	})

	c.createSyntheticBlocks = makePassFactory("createSyntheticBlocks", func(compiler AbstractCompiler) CompilerPass {
		return NewCreateSyntheticBlocks(
			compiler,
			*options.GetSyntheticBlockStartMarker(),
			*options.GetSyntheticBlockEndMarker())
	})

	c.earlyPeepholeOptimizations = defaultPass_earlyPeepholeOptimizations

	c.earlyInlineVariables = makePassFactory("earlyInlineVariables", func(compiler AbstractCompiler) CompilerPass {
		return NewInlineVariables(compiler, inlineVariablesMode(options))
	})

	c.peepholeOptimizations = defaultPass_peepholeOptimizations
	c.peepholeOptimizationsOnceNormalized = defaultPass_peepholeOptimizationsOnceNormalized
	c.peepholeOptimizationsOnceNonNormalized = defaultPass_peepholeOptimizationsOnceNonNormalized

	c.latePeepholeOptimizations = makePassFactory("latePeepholeOptimizations", func(compiler AbstractCompiler) CompilerPass {
		const late = true
		useTypes := options.ShouldUseTypesForLocalOptimization()
		return NewPeepholeOptimizationsPass(
			compiler,
			"latePeepholeOptimizations",
			NewStatementFusion(),
			NewPeepholeRemoveDeadCode(),
			NewPeepholeMinimizeConditions(late),
			NewPeepholeSubstituteAlternateSyntax(late),
			NewPeepholeReplaceKnownMethods(late, useTypes),
			NewPeepholeFoldConstants(late, useTypes))
	})

	c.checkVars = defaultPass_checkVars
	c.inferConsts = defaultPass_inferConsts

	c.checkRegExp = makePassFactory(PassNames_CHECK_REG_EXP, func(compiler AbstractCompiler) CompilerPass {
		pass := NewCheckRegExp(compiler, options.GetAssumePropertiesAreStaticallyAnalyzable(), true)
		return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
			pass.Process(externs, root)
			compiler.SetHasRegExpGlobalReferences(pass.IsGlobalRegExpPropertiesUsed())
		})
	})

	c.checkRegExpForOptimizations = makePassFactory("checkRegExpForOptimizations", func(compiler AbstractCompiler) CompilerPass {
		pass := NewCheckRegExp(compiler, options.GetAssumePropertiesAreStaticallyAnalyzable(), false)
		return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
			pass.Process(externs, root)
			compiler.SetHasRegExpGlobalReferences(pass.IsGlobalRegExpPropertiesUsed())
		})
	})

	c.checkVariableReferences = defaultPass_checkVariableReferences
	c.checkSuper = defaultPass_checkSuper
	c.clearTypedScopeCreatorPass = defaultPass_clearTypedScopeCreatorPass
	c.clearTopTypedScopePass = defaultPass_clearTopTypedScopePass
	c.inferTypes = defaultPass_inferTypes
	c.inferJsDocInfo = defaultPass_inferJsDocInfo

	c.checkTypes = makePassFactory(PassNames_CHECK_TYPES, func(compiler AbstractCompiler) CompilerPass {
		return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
			check := NewTypeCheckWithScope(
				compiler,
				compiler.GetReverseAbstractInterpreter(),
				compiler.GetTypeRegistry(),
				compiler.GetTopScope(),
				compiler.GetTypedScopeCreator().(*TypedScopeCreator)).
				ReportUnknownTypes(options.Enables(DiagnosticGroup_ForType(TypeCheck_UNKNOWN_EXPR_TYPE))).
				ReportMissingProperties(!options.Disables(DiagnosticGroup_ForType(TypeCheck_INEXISTENT_PROPERTY)))
			check.Process(externs, root)
			compiler.GetErrorManager().SetTypedPercent(check.GetTypedPercent())
		})
	})

	c.checkControlFlow = makePassFactory("checkControlFlow", func(compiler AbstractCompiler) CompilerPass {
		var callbacks []NodeTraversal_Callback
		var unreachable *CheckUnreachableCode
		if !options.Disables(DiagnosticGroups_CHECK_USELESS_CODE) {
			unreachable = NewCheckUnreachableCode(compiler)
			callbacks = append(callbacks, unreachable)
		}
		if !options.Disables(DiagnosticGroups_MISSING_RETURN) {
			mr := NewCheckMissingReturn(compiler)
			mr.sharedCfgSource = unreachable
			callbacks = append(callbacks, mr)
		}
		return defaultPassConfig_combineChecks(compiler, callbacks)
	})

	c.checkAccessControls = defaultPass_checkAccessControls
	c.lintChecks = defaultPass_lintChecks
	c.analyzerChecks = defaultPass_analyzerChecks
	c.checkRequiresAndProvidesSorted = defaultPass_checkRequiresAndProvidesSorted

	c.checkStrictMode = makePassFactory("checkStrictMode", func(compiler AbstractCompiler) CompilerPass {
		defaultLevel := CheckLevel_OFF
		if options.ExpectStrictModeInput() {
			defaultLevel = CheckLevel_ERROR
		}
		return NewStrictModeCheck(compiler, defaultLevel)
	})

	c.processTweaks = makePassFactory("processTweaks", func(compiler AbstractCompiler) CompilerPass {
		return NewProcessTweaks(compiler, options.GetTweakProcessing().ShouldStrip())
	})

	c.processDefinesCheck = c.createProcessDefines(ProcessDefines_Mode_CHECK)
	c.processDefinesOptimize = c.createProcessDefines(ProcessDefines_Mode_OPTIMIZE)

	c.stripCode = defaultPass_stripCode
	c.checkConsts = defaultPass_checkConsts
	c.rewriteCallerCodeLocation = defaultPass_rewriteCallerCodeLocation
	c.replaceToggles = defaultPass_replaceToggles

	c.replaceIdGenerators = makePassFactory(PassNames_REPLACE_ID_GENERATORS, func(compiler AbstractCompiler) CompilerPass {
		return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
			pass := NewReplaceIdGenerators(
				compiler,
				options.NeedsTranspilationOf(parser.FeatureSet_Feature_TEMPLATE_LITERALS),
				options.GetIdGenerators(),
				options.ShouldGeneratePseudoNames(),
				options.GetIdGeneratorsMapSerialized(),
				options.GetXidHashFunction())
			pass.Process(externs, root)
			compiler.SetIdGeneratorMap(pass.GetSerializedIdMappings())
		})
	})

	c.replaceStrings = makePassFactory("replaceStrings", func(compiler AbstractCompiler) CompilerPass {
		return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
			pass := NewReplaceStrings(
				compiler,
				options.GetReplaceStringsPlaceholderToken(),
				options.GetReplaceStringsFunctionDescriptions())
			pass.Process(externs, root)
			compiler.SetStringMap(pass.GetStringMap())
		})
	})

	c.closureCodeRemoval = makePassFactory("closureCodeRemoval", func(compiler AbstractCompiler) CompilerPass {
		return NewClosureCodeRemoval(
			compiler,
			options.ShouldRemoveAbstractMethods(),
			options.ShouldRemoveClosureAsserts())
	})

	c.closureOptimizePrimitives = defaultPass_closureOptimizePrimitives

	c.rescopeGlobalSymbols = makePassFactory("rescopeGlobalSymbols", func(compiler AbstractCompiler) CompilerPass {
		return NewRescopeGlobalSymbols(
			compiler,
			*options.GetRenamePrefixNamespace(),
			options.AssumeCrossChunkNamesForRenamePrefixNamespace(),
			options.GetOptimizeLocalAccessForGlobalSymbolNamespace())
	})

	c.convertChunksToESModules = defaultPass_convertChunksToESModules

	c.inlineAndCollapseProperties = makeCheckedPassFactory(
		"inlineAndCollapseProperties",
		inlineAndCollapsePrecondition,
		false,
		func(compiler AbstractCompiler) CompilerPass {
			return InlineAndCollapseProperties_Builder_(compiler).
				SetPropertyCollapseLevel(options.GetPropertyCollapseLevel()).
				SetChunkOutputType(options.GetChunkOutputType()).
				SetHaveModulesBeenRewritten(options.GetProcessCommonJSModules()).
				SetModuleResolutionMode(options.GetModuleResolutionMode()).
				Build()
		},
	)

	c.collapseObjectLiterals = defaultPass_collapseObjectLiterals

	c.disambiguateProperties = makeCheckedPassFactory(
		PassNames_DISAMBIGUATE_PROPERTIES,
		defaultPassConfig_requirePropertiesAreStaticallyAnalyzable,
		false,
		func(compiler AbstractCompiler) CompilerPass {
			return NewDisambiguateProperties(compiler, options.GetPropertiesThatMustDisambiguate())
		},
	)

	c.devirtualizeMethods = defaultPass_devirtualizeMethods
	c.optimizeCalls = defaultPass_optimizeCalls
	c.optimizeConstructors = defaultPass_optimizeConstructors
	c.markPureFunctions = defaultPass_markPureFunctions

	c.inlineVariables = makeLoopPassFactory(PassNames_INLINE_VARIABLES, func(compiler AbstractCompiler) CompilerPass {
		return NewInlineVariables(compiler, inlineVariablesMode(options))
	})

	c.inlineConstants = defaultPass_inlineConstants
	c.inlineSimpleMethods = defaultPass_inlineSimpleMethods
	c.deadAssignmentsElimination = defaultPass_deadAssignmentsElimination
	c.deadPropertyAssignmentElimination = defaultPass_deadPropertyAssignmentElimination

	c.inlineFunctions = makeLoopPassFactory(PassNames_INLINE_FUNCTIONS, func(compiler AbstractCompiler) CompilerPass {
		return NewInlineFunctions(
			compiler,
			compiler.GetUniqueNameIdSupplier(),
			options.GetInlineFunctionsLevel(),
			options.AssumeStrictThis() || options.ExpectStrictModeInput(),
			options.GetMaxFunctionSizeAfterInlining())
	})

	c.inlineProperties = defaultPass_inlineProperties
	c.isolatePolyfills = defaultPass_isolatePolyfills

	c.removeUnusedCode = makeLoopPassFactory(PassNames_REMOVE_UNUSED_CODE, func(compiler AbstractCompiler) CompilerPass {
		return NewRemoveUnusedCode_Builder(compiler).
			RemoveLocalVars(options.ShouldRemoveUnusedLocalVariables()).
			RemoveGlobals(options.ShouldRemoveUnusedVariables()).
			PreserveFunctionExpressionNames(false).
			RemoveUnusedPrototypeProperties(options.ShouldRemoveUnusedPrototypeProperties()).
			RemoveUnusedThisProperties(options.IsRemoveUnusedClassProperties()).
			RemoveUnusedObjectDefinePropertiesDefinitions(options.IsRemoveUnusedClassProperties()).
			RemoveUnusedPolyfills(
				len(options.GetForceLibraryInjectionList()) == 0 && options.GetInjectPolyfillsNewerThan() == nil).
			AssumeGettersArePure(options.GetAssumeGettersArePure()).
			Build()
	})
	c.removeUnusedCodeOnce = makePassFactory(PassNames_REMOVE_UNUSED_CODE, c.removeUnusedCode.internalFactory)

	c.crossChunkCodeMotion = makeLoopPassFactory(PassNames_CROSS_CHUNK_CODE_MOTION, func(compiler AbstractCompiler) CompilerPass {
		return NewCrossChunkCodeMotion(
			compiler,
			compiler.GetChunkGraph(),
			options.GetParentChunkCanSeeSymbolsDeclaredInChildren())
	})

	c.crossChunkMethodMotion = makeLoopPassFactory(PassNames_CROSS_CHUNK_METHOD_MOTION, func(compiler AbstractCompiler) CompilerPass {
		return NewCrossChunkMethodMotion(
			compiler,
			compiler.GetCrossChunkIdGenerator(),
			false,
			options.GetCrossChunkCodeMotionNoStubMethods())
	})

	c.flowSensitiveInlineVariables = defaultPass_flowSensitiveInlineVariables

	c.coalesceVariableNames = makePassFactory(PassNames_COALESCE_VARIABLE_NAMES, func(compiler AbstractCompiler) CompilerPass {
		return NewCoalesceVariableNames(compiler, options.ShouldGeneratePseudoNames())
	})

	c.postNormalizePeephole = makePassFactory(PassNames_POST_NORMALIZE_PEEPHOLE, func(compiler AbstractCompiler) CompilerPass {
		var optimizations []AbstractPeepholeOptimization
		if options.ShouldCollapseVariableDeclarations() {
			optimizations = append(optimizations, NewExploitAssigns())
		}
		if options.ShouldOptimizeLetAndConst() {
			optimizations = append(optimizations,
				NewOptimizeLetAndConstPeephole(options.ShouldTreatGlobalScopeAsIsolated()))
		}
		return NewPeepholeOptimizationsPassOfList(compiler, PassNames_POST_NORMALIZE_PEEPHOLE, optimizations)
	})

	c.collapseVariableDeclarations = defaultPass_collapseVariableDeclarations

	c.extractPrototypeMemberDeclarations = makeCheckedPassFactory(
		PassNames_EXTRACT_PROTOTYPE_MEMBER_DECLARATIONS,
		defaultPassConfig_requirePropertiesAreStaticallyAnalyzable,
		false,
		func(compiler AbstractCompiler) CompilerPass {
			var pattern ExtractPrototypeMemberDeclarations_Pattern
			switch options.GetExtractPrototypeMemberDeclarationsMode() {
			case CompilerOptions_ExtractPrototypeMemberDeclarationsMode_USE_GLOBAL_TEMP:
				pattern = ExtractPrototypeMemberDeclarations_Pattern_USE_GLOBAL_TEMP
			case CompilerOptions_ExtractPrototypeMemberDeclarationsMode_USE_CHUNK_TEMP:
				pattern = ExtractPrototypeMemberDeclarations_Pattern_USE_CHUNK_TEMP
			case CompilerOptions_ExtractPrototypeMemberDeclarationsMode_USE_IIFE:
				pattern = ExtractPrototypeMemberDeclarations_Pattern_USE_IIFE
			default:
				panic(jutil.NewIllegalStateExceptionMsg("unexpected"))
			}
			return NewExtractPrototypeMemberDeclarations(compiler, pattern)
		},
	)

	c.rewriteFunctionExpressions = defaultPass_rewriteFunctionExpressions
	c.collapseAnonymousFunctions = defaultPass_collapseAnonymousFunctions
	c.rewriteGlobalDeclarationsForTryCatchWrapping = defaultPass_rewriteGlobalDeclarationsForTryCatchWrapping

	c.aliasStrings = makePassFactory("aliasStrings", func(compiler AbstractCompiler) CompilerPass {
		return NewAliasStrings(
			compiler,
			compiler.GetChunkGraph(),
			options.ShouldOutputJsStringUsage(),
			options.GetAliasStringsMode())
	})

	c.ambiguateProperties = makeCheckedPassFactory(
		PassNames_AMBIGUATE_PROPERTIES,
		defaultPassConfig_requirePropertiesAreStaticallyAnalyzable,
		false,
		func(compiler AbstractCompiler) CompilerPass {
			return NewAmbiguateProperties(
				compiler,
				options.GetPropertyReservedNamingFirstChars(),
				options.GetPropertyReservedNamingNonFirstChars(),
				compiler.GetExternProperties())
		},
	)

	c.markUnnormalized = defaultPass_markUnnormalized
	c.normalize = defaultPass_normalize
	c.externExports = defaultPass_externExports

	c.denormalize = makePassFactory("denormalize", func(compiler AbstractCompiler) CompilerPass {
		return NewDenormalize(compiler, options.GetOutputFeatureSet())
	})

	c.invertContextualRenaming = defaultPass_invertContextualRenaming

	c.renameProperties = makeCheckedPassFactory(
		"renameProperties",
		defaultPassConfig_requirePropertiesAreStaticallyAnalyzable,
		false,
		func(compiler AbstractCompiler) CompilerPass {
			jutil.CheckState(options.GetPropertyRenaming() == PropertyRenamingPolicy_ALL_UNQUOTED)
			prevPropertyMap := options.GetInputPropertyMap()
			return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
				rprop := NewRenameProperties(
					compiler,
					options.ShouldGeneratePseudoNames(),
					prevPropertyMap,
					options.GetPropertyReservedNamingFirstChars(),
					options.GetPropertyReservedNamingNonFirstChars(),
					options.GetNameGenerator())
				rprop.Process(externs, root)
				compiler.SetPropertyMap(rprop.GetPropertyMap())
			})
		},
	)

	c.removePropertyRenamingCalls = defaultPass_removePropertyRenamingCalls

	c.renameVars = makePassFactory("renameVars", func(compiler AbstractCompiler) CompilerPass {
		prevVariableMap := options.GetInputVariableMap()
		return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
			compiler.SetVariableMap(c.runVariableRenaming(compiler, prevVariableMap, externs, root))
		})
	})

	c.renameLabels = defaultPass_renameLabels
	c.convertToDottedProperties = defaultPass_convertToDottedProperties
	c.checkAstValidity = defaultPass_checkAstValidity
	c.varCheckValidity = defaultPass_varCheckValidity

	c.instrumentForCodeCoverage = makePassFactory("instrumentForCodeCoverage", func(compiler AbstractCompiler) CompilerPass {
		return NewCoverageInstrumentationPassWithInstrumentOption(
			compiler,
			CoverageInstrumentationPass_CoverageReach_CONDITIONAL,
			options.GetInstrumentForCoverageOption(),
			options.GetProductionInstrumentationArrayName())
	})

	c.gatherExternPropertiesCheck = defaultPass_gatherExternPropertiesCheck
	c.gatherExternPropertiesOptimize = defaultPass_gatherExternPropertiesOptimize
	c.polymerPass = defaultPass_polymerPass
	c.chromePass = defaultPass_chromePass
	c.j2clConstantHoisterPass = defaultPass_j2clConstantHoisterPass
	c.j2clClinitPass = defaultPass_j2clClinitPass
	c.j2clPropertyInlinerPass = defaultPass_j2clPropertyInlinerPass
	c.j2clPass = defaultPass_j2clPass
	c.j2clUtilGetDefineRewriterPass = defaultPass_j2clUtilGetDefineRewriterPass
	c.j2clAssertRemovalPass = defaultPass_j2clAssertRemovalPass
	c.j2clSourceFileChecker = defaultPass_j2clSourceFileChecker
	c.j2clChecksPass = defaultPass_j2clChecksPass

	c.checkConformance = makePassFactory(PassNames_CHECK_CONFORMANCE, func(compiler AbstractCompiler) CompilerPass {
		reportingMode := options.GetConformanceReportingMode()
		return NewCheckConformanceWithConfigs(compiler, options.GetConformanceConfigs(), &reportingMode)
	})

	c.removeCastNodes = defaultPass_removeCastNodes
	c.typesToColors = defaultPass_typesToColors

	c.serializeTypedAst = makePassFactory("serializeTypedAst", func(compiler AbstractCompiler) CompilerPass {
		return SerializeTypedAstPass_CreateFromPath(
			compiler,
			options.GetTypedAstOutputFile(),
			SerializationOptions_Builder_().
				SetIncludeDebugInfo(compiler.GetOptions().ShouldSerializeExtraDebugInfo()).
				SetRuntimeLibraries(compiler.GetRuntimeJsLibManager().GetInjectedLibraries()).
				Build())
	})

	c.removeUnnecessarySyntheticExterns = defaultPass_removeUnnecessarySyntheticExterns
	c.optimizeToEs6 = defaultPass_optimizeToEs6
	c.whitespaceWrapGoogModules = defaultPass_whitespaceWrapGoogModules
	c.rewriteCommonJsModules = defaultPass_rewriteCommonJsModules
	c.rewriteScriptsToEs6Modules = defaultPass_rewriteScriptsToEs6Modules

	c.gatherModuleMetadataPass = makePassFactory(PassNames_GATHER_MODULE_METADATA, func(compiler AbstractCompiler) CompilerPass {
		// Ensure the synthetic externs input exists so metadata is created for it.
		compiler.GetSynthesizedExternsInput()
		return NewGatherModuleMetadata(
			compiler,
			options.GetProcessCommonJSModules(),
			options.GetModuleResolutionMode())
	})

	c.createModuleMapPass = defaultPass_createModuleMapPass
	c.gatherGettersAndSetters = defaultPass_gatherGettersAndSetters
	c.addSyntheticScript = defaultPass_addSyntheticScript
	c.removeSyntheticScript = defaultPass_removeSyntheticScript
	c.mergeSyntheticScript = defaultPass_mergeSyntheticScript
	c.protectLocaleData = defaultPass_protectLocaleData
	c.substituteLocaleData = defaultPass_substituteLocaleData
	c.transpileOnlyClosureUnaware = defaultPass_transpileOnlyClosureUnaware
	c.transpileAndOptimizeClosureUnaware = defaultPass_transpileAndOptimizeClosureUnaware
}

func inlineVariablesMode(options *CompilerOptions) InlineVariables_Mode {
	if options.ShouldInlineVariables() {
		return InlineVariables_Mode_ALL
	}
	if options.ShouldInlineLocalVariables() {
		return InlineVariables_Mode_LOCALS_ONLY
	}
	panic(jutil.NewIllegalStateExceptionMsg("No variable inlining option set."))
}

func (c *DefaultPassConfig) getFullReplaceMessagesPass() *PassFactory {
	options := c.options
	return makePassFactory(PassNames_REPLACE_MESSAGES, func(compiler AbstractCompiler) CompilerPass {
		return NewReplaceMessages(
			compiler,
			options.GetMessageBundle(),
			options.GetStrictMessageReplacement()).
			GetFullReplacementPass()
	})
}

func (c *DefaultPassConfig) getProtectMessagesPass() *PassFactory {
	options := c.options
	return makePassFactory("protectMessages", func(compiler AbstractCompiler) CompilerPass {
		return NewReplaceMessages(
			compiler,
			options.GetMessageBundle(),
			options.GetStrictMessageReplacement()).
			GetMsgProtectionPass()
	})
}

func (c *DefaultPassConfig) getReplaceProtectedMessagesPass() *PassFactory {
	options := c.options
	return makePassFactory("replaceProtectedMessages", func(compiler AbstractCompiler) CompilerPass {
		return NewReplaceMessages(
			compiler,
			options.GetMessageBundle(),
			options.GetStrictMessageReplacement()).
			GetReplacementCompletionPass()
	})
}

func defaultPassConfig_createPeepholeOptimizationsPass(
	compiler AbstractCompiler, passName string, expectAstIsNormalized bool,
) CompilerPass {
	jutil.CheckArgumentMsg1(
		expectAstIsNormalized == compiler.GetLifeCycleStage().IsNormalized(),
		"%s",
		compiler.GetLifeCycleStage())
	const late = false
	opts := compiler.GetOptions()
	useTypes := opts.ShouldUseTypesForLocalOptimization()
	optimizations := make([]AbstractPeepholeOptimization, 0, 10)
	if expectAstIsNormalized {
		optimizations = append(optimizations, NewMinimizeExitPoints())
	}
	optimizations = append(
		optimizations,
		NewPeepholeMinimizeConditions(late),
		NewPeepholeSubstituteAlternateSyntax(late),
		NewPeepholeReplaceKnownMethods(late, useTypes),
		NewPeepholeRemoveDeadCode(),
	)
	if opts.GetJ2clPass().ShouldAddJ2clPasses() {
		optimizations = append(
			optimizations,
			NewJ2clEqualitySameRewriterPass(useTypes),
			NewJ2clStringValueOfRewriterPass(),
			NewJ2clUndefinedChecksRewriterPass(),
		)
	}
	optimizations = append(
		optimizations,
		NewPeepholeFoldConstants(late, useTypes),
		NewPeepholeCollectPropertyAssignments(),
	)
	return NewPeepholeOptimizationsPassOfList(compiler, passName, optimizations)
}

func defaultPassConfig_combineChecks(
	compiler AbstractCompiler, callbacks []NodeTraversal_Callback,
) CompilerPass {
	jutil.CheckArgument(len(callbacks) != 0)
	return NewCombinedCompilerPassOfList(compiler, callbacks)
}

func (c *DefaultPassConfig) createProcessDefines(mode ProcessDefines_Mode) *PassFactory {
	options := c.options
	return makePassFactory("processDefines_"+mode.String(), func(compiler AbstractCompiler) CompilerPass {
		return NewProcessDefines_Builder(compiler).
			PutReplacements(defaultPassConfig_getAdditionalReplacements(options)).
			PutReplacements(options.GetDefineReplacements()).
			SetMode(mode).
			SetRecognizeClosureDefines(compiler.GetOptions().GetClosurePass()).
			SetEnableZonesDefineName(options.GetEnableZonesDefineName()).
			SetZoneInputPattern(options.GetZoneInputPattern()).
			SetUnknownDefinesToIgnore(options.GetUnknownDefinesToIgnore()).
			Build()
	})
}

func (c *DefaultPassConfig) runVariableRenaming(
	compiler AbstractCompiler, prevVariableMap *VariableMap, externs rhino.Node, root rhino.Node,
) *VariableMap {
	options := c.options
	var reservedChars ASCIICharSet
	reservedNames := jutil.NewLinkedHashSet[string]()
	if ns := options.GetRenamePrefixNamespace(); ns != nil {
		reservedNames.Add(*ns)
	}
	for name := range compiler.GetExportedNames().All() {
		reservedNames.Add(name)
	}
	for name := range parsing.ParserRunner_GetReservedVars().All() {
		reservedNames.Add(name)
	}
	var renamePrefix string
	if p := options.GetRenamePrefix(); p != nil {
		renamePrefix = *p
	}
	rn := NewRenameVars(
		compiler,
		renamePrefix,
		options.GetVariableRenaming() == VariableRenamingPolicy_LOCAL,
		options.ShouldGeneratePseudoNames(),
		options.ShouldPreferStableNames(),
		prevVariableMap,
		reservedChars,
		reservedNames,
		options.GetNameGenerator())
	rn.Process(externs, root)
	return rn.GetVariableMap()
}

func (c *DefaultPassConfig) getCustomPasses(executionTime CustomPassExecutionTime) *PassFactory {
	options := c.options
	return makePassFactory("runCustomPasses", func(compiler AbstractCompiler) CompilerPass {
		return defaultPassConfig_runInSerial(options.GetCustomPassesAt(executionTime))
	})
}

func defaultPassConfig_runInSerial(passes []CompilerPass) CompilerPass {
	return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
		for _, pass := range passes {
			pass.Process(externs, root)
		}
	})
}

func defaultPassConfig_getAdditionalReplacements(options *CompilerOptions) *jutil.LinkedHashMap[string, rhino.Node] {
	additionalReplacements := jutil.NewLinkedHashMap[string, rhino.Node]()
	if options.ShouldMarkAsCompiled() || options.GetClosurePass() {
		additionalReplacements.Put(defaultPassConfig_COMPILED_CONSTANT_NAME, rhino.IR_TrueNode())
	}
	if options.GetClosurePass() && options.GetLocale() != nil && !options.DoLateLocalization() {
		additionalReplacements.Put(defaultPassConfig_CLOSURE_LOCALE_CONSTANT_NAME, rhino.IR_String(*options.GetLocale()))
	}
	return additionalReplacements
}
