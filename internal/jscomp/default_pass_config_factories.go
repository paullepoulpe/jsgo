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
	false, "requires assumePropertiesAreStaticallyAnalyzable to be enabled")

func defaultPassConfig_requirePropertiesAreStaticallyAnalyzable(options *CompilerOptions) PassFactory_PreconditionResult {
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

func passFactoryOf[P CompilerPass](name string, fn func(AbstractCompiler) P) *PassFactory {
	return makePassFactory(name, func(c AbstractCompiler) CompilerPass { return fn(c) })
}

func loopPassFactoryOf[P CompilerPass](name string, fn func(AbstractCompiler) P) *PassFactory {
	return makeLoopPassFactory(name, func(c AbstractCompiler) CompilerPass { return fn(c) })
}

func checkedPassFactoryOf[P CompilerPass](
	name string, precond func(*CompilerOptions) PassFactory_PreconditionResult, loop bool, fn func(AbstractCompiler) P,
) *PassFactory {
	return makeCheckedPassFactory(name, precond, loop, func(c AbstractCompiler) CompilerPass { return fn(c) })
}

func actionPassFactory(name string, fn func(AbstractCompiler, rhino.Node, rhino.Node)) *PassFactory {
	return makePassFactory(name, func(c AbstractCompiler) CompilerPass {
		return CompilerPassFunc(func(externs, root rhino.Node) { fn(c, externs, root) })
	})
}

func newReplaceMessagesForOptions(c AbstractCompiler) *ReplaceMessages {
	opts := c.GetOptions()
	return NewReplaceMessages(c, opts.GetMessageBundle(), opts.GetStrictMessageReplacement())
}

// Stateless PassFactory singletons shared across all DefaultPassConfig instances.
var (
	defaultPass_checkSideEffects = makePassFactory("checkSideEffects", func(c AbstractCompiler) CompilerPass {
		opts := c.GetOptions()
		return NewCheckSideEffects(c, opts.GetCheckSuspiciousCode(), opts.ShouldProtectHiddenSideEffects())
	})
	defaultPass_stripSideEffectProtection = passFactoryOf(PassNames_STRIP_SIDE_EFFECT_PROTECTION, NewCheckSideEffects_StripProtection)
	defaultPass_suspiciousCode            = makePassFactory("suspiciousCode", func(c AbstractCompiler) CompilerPass {
		opts := c.GetOptions()
		var callbacks []NodeTraversal_Callback
		if opts.GetCheckSuspiciousCode() {
			callbacks = append(callbacks, NewCheckSuspiciousCode(), NewCheckDuplicateCase(c))
		}
		if opts.Enables(DiagnosticGroups_GLOBAL_THIS) {
			callbacks = append(callbacks, NewCheckGlobalThis(c))
		}
		if opts.Enables(DiagnosticGroups_DEBUGGER_STATEMENT_PRESENT) {
			callbacks = append(callbacks, NewCheckDebuggerStatement(c))
		}
		return defaultPassConfig_combineChecks(c, callbacks)
	})
	defaultPass_extraRequires = makePassFactory("checkExtraRequires", func(c AbstractCompiler) CompilerPass {
		return NewCheckExtraRequires(c, c.GetOptions().GetUnusedImportsToRemove())
	})
	defaultPass_checkMissingRequires = makePassFactory("checkMissingRequires", func(c AbstractCompiler) CompilerPass {
		return NewCheckMissingRequires(c, c.GetModuleMetadataMap())
	})
	defaultPass_checkJsDocAndEs6Modules = makePassFactory("checkJsDocAndEs6Modules", func(c AbstractCompiler) CompilerPass {
		return defaultPassConfig_combineChecks(c, []NodeTraversal_Callback{NewCheckJSDoc(c), NewEs6CheckModule(c)})
	})
	defaultPass_generateExports = makePassFactory(PassNames_GENERATE_EXPORTS, func(c AbstractCompiler) CompilerPass {
		conv := c.GetCodingConvention()
		pass := NewGenerateExports(c, c.GetOptions().ShouldExportLocalPropertyDefinitions(), conv.GetExportSymbolFunction(), conv.GetExportPropertyFunction())
		return CompilerPassFunc(func(externs, root rhino.Node) {
			pass.Process(externs, root)
			c.AddExportedNames(pass.GetExportedVariableNames())
		})
	})
	defaultPass_generateIjs         = passFactoryOf("generateIjs", NewConvertToTypedInterface)
	defaultPass_removeExtraRequires = passFactoryOf("removeExtraRequires", NewExtraRequireRemover)
	defaultPass_exportTestFunctions = makePassFactory(PassNames_EXPORT_TEST_FUNCTIONS, func(c AbstractCompiler) CompilerPass {
		if conv := c.GetCodingConvention(); conv.GetExportSymbolFunction() != nil {
			return NewExportTestFunctions(c, conv.GetExportSymbolFunction(), conv.GetExportPropertyFunction())
		}
		return NewErrorPass(c, defaultPassConfig_GENERATE_EXPORTS_ERROR)
	})
	defaultPass_gatherRawExports = makePassFactory(PassNames_GATHER_RAW_EXPORTS, func(c AbstractCompiler) CompilerPass {
		pass := NewGatherRawExports(c)
		return CompilerPassFunc(func(externs, root rhino.Node) {
			pass.Process(externs, root)
			c.AddExportedNames(pass.GetExportedVariableNames())
		})
	})
	defaultPass_angularPass              = passFactoryOf(PassNames_ANGULAR_PASS, NewAngularPass)
	defaultPass_replaceMessagesForChrome = makePassFactory(PassNames_REPLACE_MESSAGES, func(c AbstractCompiler) CompilerPass {
		return NewReplaceMessagesForChrome(c, NewGoogleJsMessageIdGenerator(c.GetOptions().GetTcProjectId()))
	})
	defaultPass_fullReplaceMessages = makePassFactory(PassNames_REPLACE_MESSAGES, func(c AbstractCompiler) CompilerPass {
		return newReplaceMessagesForOptions(c).GetFullReplacementPass()
	})
	defaultPass_protectMessages = makePassFactory("protectMessages", func(c AbstractCompiler) CompilerPass {
		return newReplaceMessagesForOptions(c).GetMsgProtectionPass()
	})
	defaultPass_replaceProtectedMessages = makePassFactory("replaceProtectedMessages", func(c AbstractCompiler) CompilerPass {
		return newReplaceMessagesForOptions(c).GetReplacementCompletionPass()
	})
	defaultPass_closureGoogScopeAliasesForIjs = makePassFactory("closureGoogScopeAliasesForIjs", func(c AbstractCompiler) CompilerPass {
		return ScopedAliases_Builder_(c).Build()
	})
	defaultPass_injectRuntimeLibraries = makePassFactory("InjectRuntimeLibraries", func(c AbstractCompiler) CompilerPass {
		return NewInjectRuntimeLibraries(c, jutil.NewLinkedHashSet(c.GetOptions().GetForceLibraryInjectionList()...))
	})
	defaultPass_removeWeakSources             = passFactoryOf("removeWeakSources", NewRemoveWeakSources)
	defaultPass_declaredGlobalExternsOnWindow = passFactoryOf(PassNames_DECLARED_GLOBAL_EXTERNS_ON_WINDOW, NewDeclaredGlobalExternsOnWindow)
	defaultPass_checkTypeImportCodeReferences = passFactoryOf("checkTypeImportCodeReferences", NewCheckTypeImportCodeReferences)
	defaultPass_closureCheckModule            = makePassFactory("closureCheckModule", func(c AbstractCompiler) CompilerPass {
		return NewClosureCheckModule(c, c.GetModuleMetadataMap())
	})
	defaultPass_checkClosureImports = makePassFactory("checkGoogRequires", func(c AbstractCompiler) CompilerPass {
		return NewCheckClosureImports(c, c.GetModuleMetadataMap())
	})
	defaultPass_rewriteGoogJsImports = makePassFactory("rewriteGoogJsImports", func(c AbstractCompiler) CompilerPass {
		return NewRewriteGoogJsImports(c, RewriteGoogJsImports_Mode_LINT_AND_REWRITE, c.GetModuleMap())
	})
	defaultPass_closureReplaceGetCssName = actionPassFactory("closureReplaceGetCssName", func(c AbstractCompiler, externs, jsRoot rhino.Node) {
		opts := c.GetOptions()
		var cssNames *jutil.LinkedHashSet[string]
		if opts.ShouldGatherCssNames() {
			cssNames = jutil.NewLinkedHashSet[string]()
		}
		NewReplaceCssNames(c, opts.GetCssRenamingMap(), ReplaceCssNames_CssNameCollectorFunc(func(name string) {
			if cssNames != nil {
				cssNames.Add(name)
			}
		}), opts.GetCssRenamingSkiplist()).Process(externs, jsRoot)
		c.SetCssNames(cssNames)
	})
	defaultPass_createSyntheticBlocks = makePassFactory("createSyntheticBlocks", func(c AbstractCompiler) CompilerPass {
		opts := c.GetOptions()
		return NewCreateSyntheticBlocks(c, *opts.GetSyntheticBlockStartMarker(), *opts.GetSyntheticBlockEndMarker())
	})
	defaultPass_earlyPeepholeOptimizations = makePassFactory("earlyPeepholeOptimizations", func(c AbstractCompiler) CompilerPass {
		opts := c.GetOptions()
		peepholes := []AbstractPeepholeOptimization{NewPeepholeRemoveDeadCode()}
		if opts.GetJ2clPass().ShouldAddJ2clPasses() {
			peepholes = append(peepholes, NewJ2clEqualitySameRewriterPass(opts.ShouldUseTypesForLocalOptimization()))
		}
		return NewPeepholeOptimizationsPassOfList(c, "earlyPeepholeOptimizations", peepholes)
	})
	defaultPass_earlyInlineVariables = makePassFactory("earlyInlineVariables", func(c AbstractCompiler) CompilerPass {
		return NewInlineVariables(c, inlineVariablesMode(c.GetOptions()))
	})
	defaultPass_peepholeOptimizations = makeLoopPassFactory(PassNames_PEEPHOLE_OPTIMIZATIONS, func(c AbstractCompiler) CompilerPass {
		return defaultPassConfig_createPeepholeOptimizationsPass(c, PassNames_PEEPHOLE_OPTIMIZATIONS, true)
	})
	defaultPass_peepholeOptimizationsOnceNormalized = makePassFactory(PassNames_PEEPHOLE_OPTIMIZATIONS, func(c AbstractCompiler) CompilerPass {
		return defaultPassConfig_createPeepholeOptimizationsPass(c, PassNames_PEEPHOLE_OPTIMIZATIONS, true)
	})
	defaultPass_peepholeOptimizationsOnceNonNormalized = makePassFactory(PassNames_PEEPHOLE_OPTIMIZATIONS, func(c AbstractCompiler) CompilerPass {
		return defaultPassConfig_createPeepholeOptimizationsPass(c, PassNames_PEEPHOLE_OPTIMIZATIONS, false)
	})
	defaultPass_latePeepholeOptimizations = makePassFactory("latePeepholeOptimizations", func(c AbstractCompiler) CompilerPass {
		useTypes := c.GetOptions().ShouldUseTypesForLocalOptimization()
		return NewPeepholeOptimizationsPass(
			c, "latePeepholeOptimizations",
			NewStatementFusion(), NewPeepholeRemoveDeadCode(),
			NewPeepholeMinimizeConditions(true), NewPeepholeSubstituteAlternateSyntax(true),
			NewPeepholeReplaceKnownMethods(true, useTypes), NewPeepholeFoldConstants(true, useTypes))
	})
	defaultPass_checkVars                   = passFactoryOf(PassNames_CHECK_VARS, NewVarCheck)
	defaultPass_inferConsts                 = passFactoryOf(PassNames_INFER_CONSTS, NewInferConsts)
	defaultPass_checkRegExp                 = makeCheckRegExpFactory(PassNames_CHECK_REG_EXP, true)
	defaultPass_checkRegExpForOptimizations = makeCheckRegExpFactory("checkRegExpForOptimizations", false)
	defaultPass_checkVariableReferences     = passFactoryOf(PassNames_CHECK_VARIABLE_REFERENCES, NewVariableReferenceCheck)
	defaultPass_checkSuper                  = passFactoryOf("checkSuper", NewCheckSuper)
	defaultPass_clearTypedScopeCreatorPass  = actionPassFactory("clearTypedScopeCreatorPass", func(c AbstractCompiler, _, _ rhino.Node) { c.ClearTypedScopeCreator() })
	defaultPass_clearTopTypedScopePass      = actionPassFactory("clearTopTypedScopePass", func(c AbstractCompiler, _, _ rhino.Node) {
		c.SetTopScope(nil)
		for _, input := range c.GetInputsInOrder() {
			input.SetTypedScope(nil)
		}
	})
	defaultPass_inferTypes = actionPassFactory(PassNames_INFER_TYPES, func(c AbstractCompiler, _, srcRoot rhino.Node) {
		pass := NewTypeInferencePass(c, c.GetReverseAbstractInterpreter(), c.GetTypedScopeCreator().(*TypedScopeCreator))
		c.SetTypeCheckingHasRun(true)
		c.SetTopScope(pass.InferAllScopes(srcRoot.GetParent()))
	})
	defaultPass_inferJsDocInfo = passFactoryOf("inferJsDocInfo", NewInferJSDocInfo)
	defaultPass_checkTypes     = actionPassFactory(PassNames_CHECK_TYPES, func(c AbstractCompiler, externs, root rhino.Node) {
		opts := c.GetOptions()
		check := NewTypeCheckWithScope(c, c.GetReverseAbstractInterpreter(), c.GetTypeRegistry(), c.GetTopScope(), c.GetTypedScopeCreator().(*TypedScopeCreator)).
			ReportUnknownTypes(opts.Enables(DiagnosticGroup_ForType(TypeCheck_UNKNOWN_EXPR_TYPE))).
			ReportMissingProperties(!opts.Disables(DiagnosticGroup_ForType(TypeCheck_INEXISTENT_PROPERTY)))
		check.Process(externs, root)
		c.GetErrorManager().SetTypedPercent(check.GetTypedPercent())
	})
	defaultPass_checkControlFlow = makePassFactory("checkControlFlow", func(c AbstractCompiler) CompilerPass {
		opts := c.GetOptions()
		var callbacks []NodeTraversal_Callback
		var unreachable *CheckUnreachableCode
		if !opts.Disables(DiagnosticGroups_CHECK_USELESS_CODE) {
			unreachable = NewCheckUnreachableCode(c)
			callbacks = append(callbacks, unreachable)
		}
		if !opts.Disables(DiagnosticGroups_MISSING_RETURN) {
			mr := NewCheckMissingReturn(c)
			mr.sharedCfgSource = unreachable
			callbacks = append(callbacks, mr)
		}
		return defaultPassConfig_combineChecks(c, callbacks)
	})
	defaultPass_checkAccessControls = passFactoryOf("checkAccessControls", NewCheckAccessControls)
	defaultPass_lintChecks          = makePassFactory(PassNames_LINT_CHECKS, func(c AbstractCompiler) CompilerPass {
		return defaultPassConfig_combineChecks(c, []NodeTraversal_Callback{
			NewCheckConstPrivateProperties(c), NewCheckConstantCaseNames(c), NewCheckEmptyStatements(c), NewCheckEnums(c),
			NewCheckEs6ModuleFileStructure(c), NewCheckEs6Modules(c), NewCheckNoMutatedEs6Exports(c), NewCheckGoogModuleTypeScriptName(c),
			NewCheckInterfaces(c), NewCheckJSDocStyle(c), NewCheckMissingSemicolon(c), NewCheckNullabilityModifiers(c),
			NewCheckPrimitiveAsObject(c), NewCheckPrototypeProperties(c), NewCheckUnusedPrivateProperties(c), NewCheckUnusedLabels(c),
			NewCheckUselessBlocks(c), NewCheckVar(c),
		})
	})
	defaultPass_analyzerChecks = makePassFactory(PassNames_ANALYZER_CHECKS, func(c AbstractCompiler) CompilerPass {
		return defaultPassConfig_combineChecks(c, []NodeTraversal_Callback{
			NewCheckArrayWithGoogObject(c), NewImplicitNullabilityCheck(c), NewCheckNestedNames(c),
		})
	})
	defaultPass_checkRequiresAndProvidesSorted = makePassFactory("checkRequiresAndProvidesSorted", func(c AbstractCompiler) CompilerPass {
		return defaultPassConfig_combineChecks(c, []NodeTraversal_Callback{
			NewCheckProvidesSorted(CheckProvidesSorted_Mode_COLLECT_AND_REPORT),
			NewCheckRequiresSorted(CheckRequiresSorted_Mode_COLLECT_AND_REPORT),
		})
	})
	defaultPass_checkStrictMode = makePassFactory("checkStrictMode", func(c AbstractCompiler) CompilerPass {
		if c.GetOptions().ExpectStrictModeInput() {
			return NewStrictModeCheck(c, CheckLevel_ERROR)
		}
		return NewStrictModeCheck(c, CheckLevel_OFF)
	})
	defaultPass_processTweaks = makePassFactory("processTweaks", func(c AbstractCompiler) CompilerPass {
		return NewProcessTweaks(c, c.GetOptions().GetTweakProcessing().ShouldStrip())
	})
	defaultPass_processDefinesCheck    = createProcessDefines(ProcessDefines_Mode_CHECK)
	defaultPass_processDefinesOptimize = createProcessDefines(ProcessDefines_Mode_OPTIMIZE)
	defaultPass_stripCode              = makePassFactory("stripCode", func(c AbstractCompiler) CompilerPass {
		opts := c.GetOptions()
		return NewStripCode(c, opts.GetStripTypes(), opts.GetStripNameSuffixes(), opts.GetStripNamePrefixes(), opts.GetTweakProcessing().ShouldStrip())
	})
	defaultPass_checkConsts               = passFactoryOf("checkConsts", NewConstCheck)
	defaultPass_rewriteCallerCodeLocation = passFactoryOf("rewriteCallerCodeLocation", NewRewriteCallerCodeLocation)
	defaultPass_replaceToggles            = passFactoryOf("replaceToggles", NewReplaceToggles)
	defaultPass_replaceIdGenerators       = actionPassFactory(PassNames_REPLACE_ID_GENERATORS, func(c AbstractCompiler, externs, root rhino.Node) {
		opts := c.GetOptions()
		pass := NewReplaceIdGenerators(c, opts.NeedsTranspilationOf(parser.FeatureSet_Feature_TEMPLATE_LITERALS),
			opts.GetIdGenerators(), opts.ShouldGeneratePseudoNames(), opts.GetIdGeneratorsMapSerialized(), opts.GetXidHashFunction())
		pass.Process(externs, root)
		c.SetIdGeneratorMap(pass.GetSerializedIdMappings())
	})
	defaultPass_replaceStrings = actionPassFactory("replaceStrings", func(c AbstractCompiler, externs, root rhino.Node) {
		opts := c.GetOptions()
		pass := NewReplaceStrings(c, opts.GetReplaceStringsPlaceholderToken(), opts.GetReplaceStringsFunctionDescriptions())
		pass.Process(externs, root)
		c.SetStringMap(pass.GetStringMap())
	})
	defaultPass_closureCodeRemoval = makePassFactory("closureCodeRemoval", func(c AbstractCompiler) CompilerPass {
		opts := c.GetOptions()
		return NewClosureCodeRemoval(c, opts.ShouldRemoveAbstractMethods(), opts.ShouldRemoveClosureAsserts())
	})
	defaultPass_closureOptimizePrimitives = makePassFactory("closureOptimizePrimitives", func(c AbstractCompiler) CompilerPass {
		return NewClosureOptimizePrimitives(c, c.GetOptions().GetOutputFeatureSet().Contains(parser.FeatureSet_ES2015))
	})
	defaultPass_rescopeGlobalSymbols = makePassFactory("rescopeGlobalSymbols", func(c AbstractCompiler) CompilerPass {
		opts := c.GetOptions()
		return NewRescopeGlobalSymbols(c, *opts.GetRenamePrefixNamespace(), opts.AssumeCrossChunkNamesForRenamePrefixNamespace(), opts.GetOptimizeLocalAccessForGlobalSymbolNamespace())
	})
	defaultPass_convertChunksToESModules    = passFactoryOf("convertChunksToESModules", NewConvertChunksToESModules)
	defaultPass_inlineAndCollapseProperties = makeCheckedPassFactory("inlineAndCollapseProperties", inlineAndCollapsePrecondition, false, func(c AbstractCompiler) CompilerPass {
		opts := c.GetOptions()
		return InlineAndCollapseProperties_Builder_(c).
			SetPropertyCollapseLevel(opts.GetPropertyCollapseLevel()).
			SetChunkOutputType(opts.GetChunkOutputType()).
			SetHaveModulesBeenRewritten(opts.GetProcessCommonJSModules()).
			SetModuleResolutionMode(opts.GetModuleResolutionMode()).Build()
	})
	defaultPass_collapseObjectLiterals = makeLoopPassFactory(PassNames_COLLAPSE_OBJECT_LITERALS, func(c AbstractCompiler) CompilerPass {
		return NewInlineObjectLiterals(c, c.GetUniqueNameIdSupplier())
	})
	defaultPass_disambiguateProperties = makeCheckedPassFactory(PassNames_DISAMBIGUATE_PROPERTIES, defaultPassConfig_requirePropertiesAreStaticallyAnalyzable, false, func(c AbstractCompiler) CompilerPass {
		return NewDisambiguateProperties(c, c.GetOptions().GetPropertiesThatMustDisambiguate())
	})
	defaultPass_devirtualizeMethods = makeCheckedPassFactory(PassNames_DEVIRTUALIZE_METHODS, defaultPassConfig_requirePropertiesAreStaticallyAnalyzable, false, func(c AbstractCompiler) CompilerPass {
		return OptimizeCalls_Builder_().SetCompiler(c).SetConsiderExterns(false).AddPass(NewDevirtualizeMethods(c)).Build()
	})
	defaultPass_optimizeCalls = makeCheckedPassFactory(PassNames_OPTIMIZE_CALLS, defaultPassConfig_requirePropertiesAreStaticallyAnalyzable, true, func(c AbstractCompiler) CompilerPass {
		return OptimizeCalls_Builder_().SetCompiler(c).SetConsiderExterns(false).
			AddPass(NewOptimizeReturns(c).AsCallGraphCompilerPass()).
			AddPass(NewOptimizeParameters(c).AsCallGraphCompilerPass()).Build()
	})
	defaultPass_optimizeConstructors = makePassFactory("optimizeConstructors", func(c AbstractCompiler) CompilerPass {
		return OptimizeCalls_Builder_().SetCompiler(c).SetConsiderExterns(false).AddPass(NewOptimizeConstructors(c).AsCallGraphCompilerPass()).Build()
	})
	defaultPass_markPureFunctions = checkedPassFactoryOf("markPureFunctions", defaultPassConfig_requirePropertiesAreStaticallyAnalyzable, false, NewPureFunctionIdentifier_Driver)
	defaultPass_inlineVariables   = makeLoopPassFactory(PassNames_INLINE_VARIABLES, func(c AbstractCompiler) CompilerPass {
		return NewInlineVariables(c, inlineVariablesMode(c.GetOptions()))
	})
	defaultPass_inlineConstants = makeLoopPassFactory("inlineConstants", func(c AbstractCompiler) CompilerPass {
		return NewInlineVariables(c, InlineVariables_Mode_CONSTANTS_ONLY)
	})
	defaultPass_inlineSimpleMethods               = checkedPassFactoryOf("inlineSimpleMethods", defaultPassConfig_requirePropertiesAreStaticallyAnalyzable, true, NewInlineSimpleMethods)
	defaultPass_deadAssignmentsElimination        = loopPassFactoryOf(PassNames_DEAD_ASSIGNMENT_ELIMINATION, NewDeadAssignmentsElimination)
	defaultPass_deadPropertyAssignmentElimination = checkedPassFactoryOf("deadPropertyAssignmentElimination", defaultPassConfig_requirePropertiesAreStaticallyAnalyzable, true, NewDeadPropertyAssignmentElimination)
	defaultPass_inlineFunctions                   = makeLoopPassFactory(PassNames_INLINE_FUNCTIONS, func(c AbstractCompiler) CompilerPass {
		opts := c.GetOptions()
		return NewInlineFunctions(c, c.GetUniqueNameIdSupplier(), opts.GetInlineFunctionsLevel(), opts.AssumeStrictThis() || opts.ExpectStrictModeInput(), opts.GetMaxFunctionSizeAfterInlining())
	})
	defaultPass_inlineProperties     = loopPassFactoryOf(PassNames_INLINE_PROPERTIES, NewInlineProperties)
	defaultPass_isolatePolyfills     = passFactoryOf("IsolatePolyfills", NewIsolatePolyfills)
	defaultPass_removeUnusedCode     = makeLoopPassFactory(PassNames_REMOVE_UNUSED_CODE, createRemoveUnusedCodePass)
	defaultPass_removeUnusedCodeOnce = makePassFactory(PassNames_REMOVE_UNUSED_CODE, createRemoveUnusedCodePass)
	defaultPass_crossChunkCodeMotion = makeLoopPassFactory(PassNames_CROSS_CHUNK_CODE_MOTION, func(c AbstractCompiler) CompilerPass {
		return NewCrossChunkCodeMotion(c, c.GetChunkGraph(), c.GetOptions().GetParentChunkCanSeeSymbolsDeclaredInChildren())
	})
	defaultPass_crossChunkMethodMotion = makeLoopPassFactory(PassNames_CROSS_CHUNK_METHOD_MOTION, func(c AbstractCompiler) CompilerPass {
		return NewCrossChunkMethodMotion(c, c.GetCrossChunkIdGenerator(), false, c.GetOptions().GetCrossChunkCodeMotionNoStubMethods())
	})
	defaultPass_flowSensitiveInlineVariables = passFactoryOf(PassNames_FLOW_SENSITIVE_INLINE_VARIABLES, NewFlowSensitiveInlineVariables)
	defaultPass_coalesceVariableNames        = makePassFactory(PassNames_COALESCE_VARIABLE_NAMES, func(c AbstractCompiler) CompilerPass {
		return NewCoalesceVariableNames(c, c.GetOptions().ShouldGeneratePseudoNames())
	})
	defaultPass_postNormalizePeephole = makePassFactory(PassNames_POST_NORMALIZE_PEEPHOLE, func(c AbstractCompiler) CompilerPass {
		opts := c.GetOptions()
		var optsList []AbstractPeepholeOptimization
		if opts.ShouldCollapseVariableDeclarations() {
			optsList = append(optsList, NewExploitAssigns())
		}
		if opts.ShouldOptimizeLetAndConst() {
			optsList = append(optsList, NewOptimizeLetAndConstPeephole(opts.ShouldTreatGlobalScopeAsIsolated()))
		}
		return NewPeepholeOptimizationsPassOfList(c, PassNames_POST_NORMALIZE_PEEPHOLE, optsList)
	})
	defaultPass_collapseVariableDeclarations       = passFactoryOf(PassNames_COLLAPSE_VARIABLE_DECLARATIONS, NewCollapseVariableDeclarations)
	defaultPass_extractPrototypeMemberDeclarations = makeCheckedPassFactory(PassNames_EXTRACT_PROTOTYPE_MEMBER_DECLARATIONS, defaultPassConfig_requirePropertiesAreStaticallyAnalyzable, false, func(c AbstractCompiler) CompilerPass {
		switch c.GetOptions().GetExtractPrototypeMemberDeclarationsMode() {
		case CompilerOptions_ExtractPrototypeMemberDeclarationsMode_USE_GLOBAL_TEMP:
			return NewExtractPrototypeMemberDeclarations(c, ExtractPrototypeMemberDeclarations_Pattern_USE_GLOBAL_TEMP)
		case CompilerOptions_ExtractPrototypeMemberDeclarationsMode_USE_CHUNK_TEMP:
			return NewExtractPrototypeMemberDeclarations(c, ExtractPrototypeMemberDeclarations_Pattern_USE_CHUNK_TEMP)
		case CompilerOptions_ExtractPrototypeMemberDeclarationsMode_USE_IIFE:
			return NewExtractPrototypeMemberDeclarations(c, ExtractPrototypeMemberDeclarations_Pattern_USE_IIFE)
		default:
			panic(jutil.NewIllegalStateExceptionMsg("unexpected"))
		}
	})
	defaultPass_rewriteFunctionExpressions                   = passFactoryOf(PassNames_REWRITE_FUNCTION_EXPRESSIONS, NewFunctionRewriter)
	defaultPass_collapseAnonymousFunctions                   = passFactoryOf(PassNames_COLLAPSE_ANONYMOUS_FUNCTIONS, NewCollapseAnonymousFunctions)
	defaultPass_rewriteGlobalDeclarationsForTryCatchWrapping = passFactoryOf("rewriteGlobalDeclarationsForTryCatchWrapping", NewRewriteGlobalDeclarationsForTryCatchWrapping)
	defaultPass_aliasStrings                                 = makePassFactory("aliasStrings", func(c AbstractCompiler) CompilerPass {
		opts := c.GetOptions()
		return NewAliasStrings(c, c.GetChunkGraph(), opts.ShouldOutputJsStringUsage(), opts.GetAliasStringsMode())
	})
	defaultPass_ambiguateProperties = makeCheckedPassFactory(PassNames_AMBIGUATE_PROPERTIES, defaultPassConfig_requirePropertiesAreStaticallyAnalyzable, false, func(c AbstractCompiler) CompilerPass {
		opts := c.GetOptions()
		return NewAmbiguateProperties(c, opts.GetPropertyReservedNamingFirstChars(), opts.GetPropertyReservedNamingNonFirstChars(), c.GetExternProperties())
	})
	defaultPass_markUnnormalized = actionPassFactory("markUnnormalized", func(c AbstractCompiler, _, _ rhino.Node) {
		c.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_RAW)
	})
	defaultPass_normalize     = passFactoryOf(PassNames_NORMALIZE, Normalize_CreateNormalizeForOptimizations)
	defaultPass_externExports = passFactoryOf(PassNames_EXTERN_EXPORTS, NewExternExportsPass)
	defaultPass_denormalize   = makePassFactory("denormalize", func(c AbstractCompiler) CompilerPass {
		return NewDenormalize(c, c.GetOptions().GetOutputFeatureSet())
	})
	defaultPass_invertContextualRenaming = passFactoryOf("invertContextualRenaming", MakeDeclaredNamesUnique_GetContextualRenameInverter)
	defaultPass_renameProperties         = makeCheckedPassFactory("renameProperties", defaultPassConfig_requirePropertiesAreStaticallyAnalyzable, false, func(c AbstractCompiler) CompilerPass {
		opts := c.GetOptions()
		jutil.CheckState(opts.GetPropertyRenaming() == PropertyRenamingPolicy_ALL_UNQUOTED)
		prevPropertyMap := opts.GetInputPropertyMap()
		return CompilerPassFunc(func(externs, root rhino.Node) {
			rprop := NewRenameProperties(c, opts.ShouldGeneratePseudoNames(), prevPropertyMap,
				opts.GetPropertyReservedNamingFirstChars(), opts.GetPropertyReservedNamingNonFirstChars(), opts.GetNameGenerator())
			rprop.Process(externs, root)
			c.SetPropertyMap(rprop.GetPropertyMap())
		})
	})
	defaultPass_removePropertyRenamingCalls = passFactoryOf("removePropertyRenamingCalls", NewRemovePropertyRenamingCalls)
	defaultPass_renameVars                  = makePassFactory("renameVars", func(c AbstractCompiler) CompilerPass {
		prevVariableMap := c.GetOptions().GetInputVariableMap()
		return CompilerPassFunc(func(externs, root rhino.Node) {
			c.SetVariableMap(runVariableRenaming(c, prevVariableMap, externs, root))
		})
	})
	defaultPass_renameLabels              = passFactoryOf("renameLabels", NewRenameLabels)
	defaultPass_convertToDottedProperties = passFactoryOf(PassNames_CONVERT_TO_DOTTED_PROPERTIES, NewConvertToDottedProperties)
	defaultPass_checkAstValidity          = passFactoryOf("checkAstValidity", NewAstValidator)
	defaultPass_varCheckValidity          = makePassFactory("varCheckValidity", func(c AbstractCompiler) CompilerPass { return NewVarCheckWithValidityCheck(c, true) })
	defaultPass_instrumentForCodeCoverage = makePassFactory("instrumentForCodeCoverage", func(c AbstractCompiler) CompilerPass {
		opts := c.GetOptions()
		return NewCoverageInstrumentationPassWithInstrumentOption(c, CoverageInstrumentationPass_CoverageReach_CONDITIONAL, opts.GetInstrumentForCoverageOption(), opts.GetProductionInstrumentationArrayName())
	})
	defaultPass_gatherExternPropertiesCheck = makePassFactory("gatherExternProperties", func(c AbstractCompiler) CompilerPass {
		return NewGatherExternProperties(c, GatherExternProperties_Mode_CHECK)
	})
	defaultPass_gatherExternPropertiesOptimize = makePassFactory("gatherExternProperties", func(c AbstractCompiler) CompilerPass {
		return NewGatherExternProperties(c, GatherExternProperties_Mode_OPTIMIZE)
	})
	defaultPass_polymerPass             = passFactoryOf("polymerPass", NewPolymerPass)
	defaultPass_chromePass              = passFactoryOf("chromePass", NewChromePass)
	defaultPass_j2clConstantHoisterPass = loopPassFactoryOf("j2clConstantHoisterPass", NewJ2clConstantHoisterPass)
	defaultPass_j2clClinitPass          = makeLoopPassFactory("j2clClinitPass", func(c AbstractCompiler) CompilerPass {
		return NewJ2clClinitPrunerPass(c, c.GetChangeTracker().GetChangedScopeNodesForPass("j2clClinitPass"))
	})
	defaultPass_j2clPropertyInlinerPass       = passFactoryOf("j2clPropertyInlinerPass", NewJ2clPropertyInlinerPass)
	defaultPass_j2clPass                      = passFactoryOf("j2clPass", NewJ2clPass)
	defaultPass_j2clUtilGetDefineRewriterPass = passFactoryOf("j2clUtilGetDefineRewriterPass", NewJ2clUtilGetDefineRewriterPass)
	defaultPass_j2clAssertRemovalPass         = passFactoryOf("j2clAssertRemovalPass", NewJ2clAssertRemovalPass)
	defaultPass_j2clSourceFileChecker         = passFactoryOf("j2clSourceFileChecker", NewJ2clSourceFileChecker)
	defaultPass_j2clChecksPass                = passFactoryOf("j2clChecksPass", NewJ2clChecksPass)
	defaultPass_checkConformance              = makePassFactory(PassNames_CHECK_CONFORMANCE, func(c AbstractCompiler) CompilerPass {
		opts := c.GetOptions()
		mode := opts.GetConformanceReportingMode()
		return NewCheckConformanceWithConfigs(c, opts.GetConformanceConfigs(), &mode)
	})
	defaultPass_removeCastNodes = passFactoryOf("removeCastNodes", NewRemoveCastNodes)
	defaultPass_typesToColors   = actionPassFactory("typesToColors", func(c AbstractCompiler, externs, js rhino.Node) {
		NewConvertTypesToColors(c, SerializationOptions_Builder_().SetIncludeDebugInfo(c.GetOptions().ShouldSerializeExtraDebugInfo()).Build()).Process(externs, js)
		c.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_COLORS_AND_SIMPLIFIED_JSDOC)
	})
	defaultPass_serializeTypedAst = makePassFactory("serializeTypedAst", func(c AbstractCompiler) CompilerPass {
		opts := c.GetOptions()
		return SerializeTypedAstPass_CreateFromPath(c, opts.GetTypedAstOutputFile(),
			SerializationOptions_Builder_().SetIncludeDebugInfo(opts.ShouldSerializeExtraDebugInfo()).
				SetRuntimeLibraries(c.GetRuntimeJsLibManager().GetInjectedLibraries()).Build())
	})
	defaultPass_removeUnnecessarySyntheticExterns = passFactoryOf("removeUnnecessarySyntheticExterns", NewRemoveUnnecessarySyntheticExterns)
	defaultPass_optimizeToEs6                     = passFactoryOf("optimizeToEs6", NewSubstituteEs6Syntax)
	defaultPass_whitespaceWrapGoogModules         = passFactoryOf("whitespaceWrapGoogModules", NewWhitespaceWrapGoogModules)
	defaultPass_rewriteCommonJsModules            = passFactoryOf(PassNames_REWRITE_COMMON_JS_MODULES, NewProcessCommonJSModules)
	defaultPass_rewriteScriptsToEs6Modules        = passFactoryOf(PassNames_REWRITE_SCRIPTS_TO_ES6_MODULES, NewEs6RewriteScriptsToModules)
	defaultPass_gatherModuleMetadataPass          = makePassFactory(PassNames_GATHER_MODULE_METADATA, func(c AbstractCompiler) CompilerPass {
		c.GetSynthesizedExternsInput()
		opts := c.GetOptions()
		return NewGatherModuleMetadata(c, opts.GetProcessCommonJSModules(), opts.GetModuleResolutionMode())
	})
	defaultPass_createModuleMapPass = makePassFactory(PassNames_CREATE_MODULE_MAP, func(c AbstractCompiler) CompilerPass {
		return NewModuleMapCreator(c, c.GetModuleMetadataMap())
	})
	defaultPass_gatherGettersAndSetters = passFactoryOf(PassNames_GATHER_GETTERS_AND_SETTERS, NewGatherGetterAndSetterProperties)
	defaultPass_addSyntheticScript      = actionPassFactory("ADD_SYNTHETIC_SCRIPT", func(c AbstractCompiler, _, _ rhino.Node) { c.InitializeSyntheticCodeInput() })
	defaultPass_removeSyntheticScript   = actionPassFactory("REMOVE_SYNTHETIC_SCRIPT", func(c AbstractCompiler, _, _ rhino.Node) { c.RemoveSyntheticCodeInput() })
	defaultPass_mergeSyntheticScript    = actionPassFactory("MERGE_SYNTHETIC_SCRIPT", func(c AbstractCompiler, _, _ rhino.Node) { c.MergeSyntheticCodeInput() })
	defaultPass_protectLocaleData       = passFactoryOf("protectLocaleData", NewLocaleDataPasses_ProtectGoogLocale)
	defaultPass_substituteLocaleData    = makePassFactory("SubstituteLocaleData", func(c AbstractCompiler) CompilerPass {
		return NewLocaleDataPasses_LocaleSubstitutions(c, c.GetOptions().GetLocale())
	})
	defaultPass_transpileOnlyClosureUnaware = makePassFactory("TranspileOnlyClosureUnaware", func(c AbstractCompiler) CompilerPass {
		return NewTranspileAndOptimizeClosureUnaware(c, NestedCompilerRunner_Mode_TRANSPILE_ONLY)
	})
	defaultPass_transpileAndOptimizeClosureUnaware = makePassFactory("TranspileAndOptimizeClosureUnaware", func(c AbstractCompiler) CompilerPass {
		return NewTranspileAndOptimizeClosureUnaware(c, NestedCompilerRunner_Mode_TRANSPILE_AND_OPTIMIZE)
	})
	defaultPass_customPassesBeforeOptimizationLoop = makePassFactory("runCustomPasses", func(c AbstractCompiler) CompilerPass {
		return defaultPassConfig_runInSerial(c.GetOptions().GetCustomPassesAt(CustomPassExecutionTime_BEFORE_OPTIMIZATION_LOOP))
	})
	defaultPass_customPassesAfterOptimizationLoop = makePassFactory("runCustomPasses", func(c AbstractCompiler) CompilerPass {
		return defaultPassConfig_runInSerial(c.GetOptions().GetCustomPassesAt(CustomPassExecutionTime_AFTER_OPTIMIZATION_LOOP))
	})

	defaultPassConfig_forbidDynamicImportUsage = passFactoryOf("FORBID_DYNAMIC_IMPORT", NewForbidDynamicImportUsage)
	defaultPassConfig_rewriteDynamicImports    = makePassFactory("REWRITE_DYNAMIC_IMPORT", func(c AbstractCompiler) CompilerPass {
		opts := c.GetOptions()
		return NewRewriteDynamicImports(c, opts.GetDynamicImportAlias(), opts.GetChunkOutputType())
	})
)

func makeCheckRegExpFactory(name string, reportErrors bool) *PassFactory {
	return makePassFactory(name, func(c AbstractCompiler) CompilerPass {
		pass := NewCheckRegExp(c, c.GetOptions().GetAssumePropertiesAreStaticallyAnalyzable(), reportErrors)
		return CompilerPassFunc(func(externs, root rhino.Node) {
			pass.Process(externs, root)
			c.SetHasRegExpGlobalReferences(pass.IsGlobalRegExpPropertiesUsed())
		})
	})
}

func (c *DefaultPassConfig) makeClosurePrimitives() *PassFactory {
	return makePassFactory("closurePrimitives", func(compiler AbstractCompiler) CompilerPass {
		c.preprocessorSymbolTableFactory.MaybeInitialize(compiler)
		return NewProcessClosurePrimitives(compiler)
	})
}

func (c *DefaultPassConfig) makeClosureProvidesRequires() *PassFactory {
	return makePassFactory("closureProvidesRequires", func(compiler AbstractCompiler) CompilerPass {
		c.preprocessorSymbolTableFactory.MaybeInitialize(compiler)
		pass := NewProcessClosureProvidesAndRequires(compiler, compiler.GetOptions().ShouldPreservesGoogProvidesAndRequires())
		return CompilerPassFunc(func(externs, root rhino.Node) {
			pass.Process(externs, root)
			compiler.AddExportedNames(pass.GetExportedVariableNames())
		})
	})
}

func (c *DefaultPassConfig) makeClosureGoogScopeAliases() *PassFactory {
	return makePassFactory("closureGoogScopeAliases", func(compiler AbstractCompiler) CompilerPass {
		c.preprocessorSymbolTableFactory.MaybeInitialize(compiler)
		return ScopedAliases_Builder_(compiler).
			SetPreprocessorSymbolTable(c.preprocessorSymbolTableFactory.GetInstanceOrNull()).
			SetModuleMetadataMap(compiler.GetModuleMetadataMap()).
			SetInvalidModuleGetHandling(ScopedAliases_InvalidModuleGetHandling_GIVE_UNIQUE_NAME).Build()
	})
}

func (c *DefaultPassConfig) makeClosureRewriteModule() *PassFactory {
	return makePassFactory("closureRewriteModule", func(compiler AbstractCompiler) CompilerPass {
		c.preprocessorSymbolTableFactory.MaybeInitialize(compiler)
		return NewClosureRewriteModule(compiler, c.preprocessorSymbolTableFactory.GetInstanceOrNull(), compiler.GetTopScope())
	})
}

func createRemoveUnusedCodePass(c AbstractCompiler) CompilerPass {
	opts := c.GetOptions()
	return NewRemoveUnusedCode_Builder(c).
		RemoveLocalVars(opts.ShouldRemoveUnusedLocalVariables()).RemoveGlobals(opts.ShouldRemoveUnusedVariables()).
		PreserveFunctionExpressionNames(false).RemoveUnusedPrototypeProperties(opts.ShouldRemoveUnusedPrototypeProperties()).
		RemoveUnusedThisProperties(opts.IsRemoveUnusedClassProperties()).
		RemoveUnusedObjectDefinePropertiesDefinitions(opts.IsRemoveUnusedClassProperties()).
		RemoveUnusedPolyfills(len(opts.GetForceLibraryInjectionList()) == 0 && opts.GetInjectPolyfillsNewerThan() == nil).
		AssumeGettersArePure(opts.GetAssumeGettersArePure()).Build()
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

func defaultPassConfig_createPeepholeOptimizationsPass(c AbstractCompiler, passName string, expectAstIsNormalized bool) CompilerPass {
	jutil.CheckArgumentMsg1(expectAstIsNormalized == c.GetLifeCycleStage().IsNormalized(), "%s", c.GetLifeCycleStage())
	opts := c.GetOptions()
	useTypes := opts.ShouldUseTypesForLocalOptimization()
	optimizations := make([]AbstractPeepholeOptimization, 0, 10)
	if expectAstIsNormalized {
		optimizations = append(optimizations, NewMinimizeExitPoints())
	}
	optimizations = append(optimizations, NewPeepholeMinimizeConditions(false), NewPeepholeSubstituteAlternateSyntax(false),
		NewPeepholeReplaceKnownMethods(false, useTypes), NewPeepholeRemoveDeadCode())
	if opts.GetJ2clPass().ShouldAddJ2clPasses() {
		optimizations = append(optimizations, NewJ2clEqualitySameRewriterPass(useTypes), NewJ2clStringValueOfRewriterPass(), NewJ2clUndefinedChecksRewriterPass())
	}
	optimizations = append(optimizations, NewPeepholeFoldConstants(false, useTypes), NewPeepholeCollectPropertyAssignments())
	return NewPeepholeOptimizationsPassOfList(c, passName, optimizations)
}

func defaultPassConfig_combineChecks(c AbstractCompiler, callbacks []NodeTraversal_Callback) CompilerPass {
	jutil.CheckArgument(len(callbacks) != 0)
	return NewCombinedCompilerPassOfList(c, callbacks)
}

func createProcessDefines(mode ProcessDefines_Mode) *PassFactory {
	return makePassFactory("processDefines_"+mode.String(), func(c AbstractCompiler) CompilerPass {
		opts := c.GetOptions()
		return NewProcessDefines_Builder(c).
			PutReplacements(defaultPassConfig_getAdditionalReplacements(opts)).PutReplacements(opts.GetDefineReplacements()).
			SetMode(mode).SetRecognizeClosureDefines(opts.GetClosurePass()).
			SetEnableZonesDefineName(opts.GetEnableZonesDefineName()).SetZoneInputPattern(opts.GetZoneInputPattern()).
			SetUnknownDefinesToIgnore(opts.GetUnknownDefinesToIgnore()).Build()
	})
}

func runVariableRenaming(c AbstractCompiler, prevVariableMap *VariableMap, externs, root rhino.Node) *VariableMap {
	opts := c.GetOptions()
	reservedNames := jutil.NewLinkedHashSet[string]()
	if ns := opts.GetRenamePrefixNamespace(); ns != nil {
		reservedNames.Add(*ns)
	}
	for name := range c.GetExportedNames().All() {
		reservedNames.Add(name)
	}
	for name := range parsing.ParserRunner_GetReservedVars().All() {
		reservedNames.Add(name)
	}
	var renamePrefix string
	if p := opts.GetRenamePrefix(); p != nil {
		renamePrefix = *p
	}
	var reservedChars ASCIICharSet
	rn := NewRenameVars(c, renamePrefix, opts.GetVariableRenaming() == VariableRenamingPolicy_LOCAL,
		opts.ShouldGeneratePseudoNames(), opts.ShouldPreferStableNames(),
		prevVariableMap, reservedChars, reservedNames, opts.GetNameGenerator())
	rn.Process(externs, root)
	return rn.GetVariableMap()
}

func defaultPassConfig_runInSerial(passes []CompilerPass) CompilerPass {
	return CompilerPassFunc(func(externs, root rhino.Node) {
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
