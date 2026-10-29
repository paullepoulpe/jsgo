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
)

// DefaultPassConfig holds the pass factories and pipeline schedule for native
// JSCompiler passes across whitespace-only, transpile-only, checks,
// optimizations, and finalizations stages.
type DefaultPassConfig struct {
	options *CompilerOptions

	// Symbol table for registering references removed during preprocessing.
	preprocessorSymbolTableFactory *PreprocessorSymbolTable_CachedInstanceFactory

	inferTypes *PassFactory
}

const (
	defaultPassConfig_COMPILED_CONSTANT_NAME       = "COMPILED"
	defaultPassConfig_CLOSURE_LOCALE_CONSTANT_NAME = "goog.LOCALE"
)

var (
	defaultPassConfig_GLOBAL_SYMBOL_NAMESPACE_PATTERN = jutil.Pattern_Compile("^[a-zA-Z0-9$_]+$")
	defaultPassConfig_GENERATE_EXPORTS_ERROR          = DiagnosticType_Error(
		"JSC_GENERATE_EXPORTS_ERROR",
		"Exports can only be generated if export symbol/property functions are set.",
	)

	// Shared immutable marker passes used across pipeline construction.
	emptyBeforeStandardChecks                     = PassFactory_CreateEmptyPass("beforeStandardChecks")
	emptyAfterStandardChecks                      = PassFactory_CreateEmptyPass(PassNames_AFTER_STANDARD_CHECKS)
	emptyBeforeTypeChecking                       = PassFactory_CreateEmptyPass(PassNames_BEFORE_TYPE_CHECKING)
	emptyBeforeSerialization                      = PassFactory_CreateEmptyPass(PassNames_BEFORE_SERIALIZATION)
	emptyBeforeEarlyOptimizationsTranspilation    = PassFactory_CreateEmptyPass(PassNames_BEFORE_EARLY_OPTIMIZATIONS_TRANSPILATION)
	emptyBeforeStandardOptimizations              = PassFactory_CreateEmptyPass(PassNames_BEFORE_STANDARD_OPTIMIZATIONS)
	emptyObfuscationPassMarker                    = PassFactory_CreateEmptyPass(PassNames_OBFUSCATION_PASS_MARKER)
	emptyOptimizationsHalfwayPoint                = PassFactory_CreateEmptyPass(PassNames_OPTIMIZATIONS_HALFWAY_POINT)
	emptyBeforeEarlyOptimizationLoop              = PassFactory_CreateEmptyPass(PassNames_BEFORE_EARLY_OPTIMIZATION_LOOP)
	emptyAfterEarlyOptimizationLoop               = PassFactory_CreateEmptyPass(PassNames_AFTER_EARLY_OPTIMIZATION_LOOP)
	emptyBeforeMainOptimizations                  = PassFactory_CreateEmptyPass(PassNames_BEFORE_MAIN_OPTIMIZATIONS)
	emptyAfterMainOptimizations                   = PassFactory_CreateEmptyPass(PassNames_AFTER_MAIN_OPTIMIZATIONS)
	emptyBeforeChunkMotion                        = PassFactory_CreateEmptyPass("beforeChunkMotion")
	emptyAfterChunkMotion                         = PassFactory_CreateEmptyPass("afterChunkMotion")
	emptyBeforeExtractPrototypeMemberDeclarations = PassFactory_CreateEmptyPass(PassNames_BEFORE_EXTRACT_PROTOTYPE_MEMBER_DECLARATIONS)
	emptyBeforeRenameProperties                   = PassFactory_CreateEmptyPass(PassNames_BEFORE_RENAME_PROPERTIES)
	emptyBeforeVariableRenaming                   = PassFactory_CreateEmptyPass(PassNames_BEFORE_VARIABLE_RENAMING)
)

func NewDefaultPassConfig(options *CompilerOptions) *DefaultPassConfig {
	return &DefaultPassConfig{
		options:                        options,
		preprocessorSymbolTableFactory: NewPreprocessorSymbolTable_CachedInstanceFactory(),
		inferTypes:                     defaultPass_inferTypes,
	}
}

func (c *DefaultPassConfig) Options() *CompilerOptions {
	return c.options
}

func (c *DefaultPassConfig) GetBasePassConfig() PassConfig {
	return c
}

func (c *DefaultPassConfig) GetPreprocessorSymbolTable() *PreprocessorSymbolTable {
	return c.preprocessorSymbolTableFactory.GetInstanceOrNull()
}

func (c *DefaultPassConfig) GetTranspileOnlyPasses() *PassListBuilder {
	options := c.options
	passes := NewPassListBuilder(options)

	// Certain errors in block-scoped variable declarations prevent correct transpilation.
	passes.MaybeAdd(defaultPass_checkVariableReferences)
	passes.MaybeAdd(defaultPass_checkVars)

	passes.MaybeAdd(defaultPass_gatherModuleMetadataPass)
	passes.MaybeAdd(defaultPass_createModuleMapPass)

	if options.GetLanguageIn().ToFeatureSet().Has(parser.FeatureSet_Feature_MODULES) {
		passes.MaybeAdd(defaultPass_rewriteGoogJsImports)
		switch options.GetEs6ModuleTranspilation() {
		case CompilerOptions_Es6ModuleTranspilation_COMPILE:
			TranspilationPasses_AddEs6ModulePass(passes, c.preprocessorSymbolTableFactory)
		case CompilerOptions_Es6ModuleTranspilation_TO_COMMON_JS_LIKE_MODULES:
			TranspilationPasses_AddEs6ModuleToCjsPass(passes)
		case CompilerOptions_Es6ModuleTranspilation_RELATIVIZE_IMPORT_PATHS:
			TranspilationPasses_AddEs6RewriteImportPathPass(passes)
		case CompilerOptions_Es6ModuleTranspilation_NONE:
		}
	}

	passes.MaybeAdd(defaultPass_checkSuper)

	TranspilationPasses_AddTranspilationRuntimeLibraries(passes)

	passes.MaybeAdd(defaultPass_transpileOnlyClosureUnaware)

	if options.NeedsTranspilationFrom(parser.FeatureSet_ES2015) && options.GetRewritePolyfills() {
		if options.GetIsolatePolyfills() {
			panic(jutil.NewIllegalStateExceptionMsg("Polyfill isolation cannot be used in transpileOnly mode"))
		}
		TranspilationPasses_AddRewritePolyfillPass(passes)
	} else if options.GetInjectPolyfillsNewerThan() != nil {
		TranspilationPasses_AddRewritePolyfillPass(passes)
	}

	passes.MaybeAdd(defaultPass_injectRuntimeLibraries)

	// Passes below this point rely on normalization and must maintain it.
	passes.MaybeAdd(defaultPass_normalize)
	passes.MaybeAdd(defaultPass_gatherGettersAndSetters)

	TranspilationPasses_AddTranspilationPasses(passes, options)

	passes.MaybeAdd(defaultPass_markUnnormalized)
	// Restore original variable/parameter names (e.g. for frameworks doing DI by parameter name).
	passes.MaybeAdd(defaultPass_invertContextualRenaming)
	// Es6ConvertSuper may emit property renaming calls that must be stripped in transpile-only mode.
	passes.MaybeAdd(defaultPass_removePropertyRenamingCalls)

	passes.AssertAllOneTimePasses()
	assertValidOrderForChecks(passes, nil, nil, nil)
	return passes
}

func (c *DefaultPassConfig) GetWhitespaceOnlyPasses() *PassListBuilder {
	options := c.options
	passes := NewPassListBuilder(options)

	if options.GetProcessCommonJSModules() {
		passes.MaybeAdd(defaultPass_rewriteCommonJsModules)
	} else if options.GetLanguageIn().ToFeatureSet().Has(parser.FeatureSet_Feature_MODULES) {
		passes.MaybeAdd(defaultPass_rewriteScriptsToEs6Modules)
	}

	if options.ShouldWrapGoogModulesForWhitespaceOnly() {
		passes.MaybeAdd(defaultPass_whitespaceWrapGoogModules)
	}
	return passes
}

func (c *DefaultPassConfig) addModuleRewritingPasses(
	checks *PassListBuilder, options *CompilerOptions, closureRewriteModule **PassFactory,
) {
	if options.GetLanguageIn().ToFeatureSet().Has(parser.FeatureSet_Feature_MODULES) {
		checks.MaybeAdd(defaultPass_rewriteGoogJsImports)
		TranspilationPasses_AddEs6ModulePass(checks, c.preprocessorSymbolTableFactory)
	}
	if options.GetClosurePass() {
		if *closureRewriteModule == nil {
			*closureRewriteModule = c.makeClosureRewriteModule()
		}
		checks.MaybeAdd(*closureRewriteModule)
	}
}

func (c *DefaultPassConfig) GetChecks() *PassListBuilder {
	options := c.options
	checks := NewPassListBuilder(options)

	jutil.CheckStateMsg(
		!options.GetSkipNonTranspilationPasses(),
		"options.getSkipNonTranspilationPasses() cannot be mixed with PassConfig::getChecks. Call"+
			" PassConfig::getTranspileOnlyPasses instead.")

	if options.IsPropertyRenamingOnlyCompilationMode() {
		var closureRewriteModule *PassFactory
		checks.MaybeAdd(defaultPass_addSyntheticScript)
		checks.MaybeAdd(defaultPass_gatherGettersAndSetters)
		checks.MaybeAdd(defaultPass_gatherModuleMetadataPass)
		checks.MaybeAdd(defaultPass_createModuleMapPass)
		checks.MaybeAdd(defaultPass_declaredGlobalExternsOnWindow)
		checks.MaybeAdd(defaultPass_checkSideEffects)
		checks.MaybeAdd(defaultPass_angularPass)
		checks.MaybeAdd(c.makeClosureGoogScopeAliases())
		checks.MaybeAdd(c.makeClosurePrimitives())
		c.addModuleRewritingPasses(checks, options, &closureRewriteModule)
		checks.MaybeAdd(defaultPass_clearTypedScopeCreatorPass)
		checks.MaybeAdd(defaultPass_clearTopTypedScopePass)
		checks.MaybeAdd(defaultPass_generateExports)
		checks.MaybeAdd(emptyAfterStandardChecks)
		checks.MaybeAdd(defaultPass_mergeSyntheticScript)
		checks.MaybeAdd(defaultPass_gatherExternPropertiesCheck)
		checks.MaybeAdd(emptyBeforeSerialization)
		return checks
	}

	if options.ShouldGenerateTypedExterns() {
		checks.MaybeAdd(defaultPass_addSyntheticScript)
		checks.MaybeAdd(defaultPass_closureGoogScopeAliasesForIjs)
		checks.MaybeAdd(defaultPass_generateIjs)
		checks.MaybeAdd(defaultPass_removeExtraRequires)
		if options.ShouldWrapGoogModulesForWhitespaceOnly() {
			checks.MaybeAdd(defaultPass_whitespaceWrapGoogModules)
		}
		checks.MaybeAdd(defaultPass_removeSyntheticScript)
		return checks
	}

	// Run before any pass tries to inject runtime libraries.
	checks.MaybeAdd(defaultPass_addSyntheticScript)
	checks.MaybeAdd(defaultPass_gatherGettersAndSetters)

	inFeatures := options.GetLanguageIn().ToFeatureSet()
	if inFeatures.ContainsFeature(parser.FeatureSet_Feature_DYNAMIC_IMPORT) && !options.ShouldAllowDynamicImport() {
		checks.MaybeAdd(defaultPassConfig_forbidDynamicImportUsage)
	}

	checks.MaybeAdd(emptyBeforeStandardChecks)

	if !options.GetProcessCommonJSModules() && inFeatures.Has(parser.FeatureSet_Feature_MODULES) {
		checks.MaybeAdd(defaultPass_rewriteScriptsToEs6Modules)
	}

	// Run after promoting scripts to modules, but before rewriting any other modules.
	checks.MaybeAdd(defaultPass_gatherModuleMetadataPass)
	checks.MaybeAdd(defaultPass_createModuleMapPass)

	if options.GetProcessCommonJSModules() {
		checks.MaybeAdd(defaultPass_rewriteCommonJsModules)
	}

	// ChromePass can rewrite invalid @type annotations into valid ones, so run before JSDoc checks.
	if options.IsChromePassEnabled() {
		checks.MaybeAdd(defaultPass_chromePass)
	}

	checks.MaybeAdd(defaultPass_checkJsDocAndEs6Modules)
	checks.MaybeAdd(defaultPass_checkTypeImportCodeReferences)

	if options.Enables(DiagnosticGroups_LINT_CHECKS) {
		checks.MaybeAdd(defaultPass_lintChecks)
		if options.GetClosurePass() {
			checks.MaybeAdd(defaultPass_checkRequiresAndProvidesSorted)
		}
	}

	if options.Enables(DiagnosticGroups_EXTRA_REQUIRE) {
		checks.MaybeAdd(defaultPass_extraRequires)
	}
	if options.Enables(DiagnosticGroups_MISSING_REQUIRE) {
		checks.MaybeAdd(defaultPass_checkMissingRequires)
	}

	checks.MaybeAdd(defaultPass_declaredGlobalExternsOnWindow)

	if !options.GetProcessCommonJSModules() {
		checks.MaybeAdd(defaultPass_checkVariableReferences)
		checks.MaybeAdd(defaultPass_checkVars)
	}

	if options.GetClosurePass() {
		checks.MaybeAdd(defaultPass_checkClosureImports)
	}

	checks.MaybeAdd(defaultPass_checkStrictMode)

	if options.GetClosurePass() {
		checks.MaybeAdd(defaultPass_closureCheckModule)
	}

	checks.MaybeAdd(defaultPass_checkSuper)
	checks.MaybeAdd(defaultPass_checkSideEffects)

	if options.GetAngularPass() {
		checks.MaybeAdd(defaultPass_angularPass)
	}

	var closureGoogScopeAliases, closurePrimitives, closureRewriteModule *PassFactory
	if options.GetClosurePass() {
		closureGoogScopeAliases = c.makeClosureGoogScopeAliases()
		closurePrimitives = c.makeClosurePrimitives()
		checks.MaybeAdd(closureGoogScopeAliases)
		checks.MaybeAdd(closurePrimitives)
	}

	if options.ShouldRewriteModulesBeforeTypechecking() {
		c.addModuleRewritingPasses(checks, options, &closureRewriteModule)
	}

	// PolymerPass must run after ClosurePrimitives/ChromePass and before suspiciousCode checks.
	if options.IsPolymerPassEnabled() {
		checks.MaybeAdd(defaultPass_polymerPass)
	}

	if options.GetProcessCommonJSModules() {
		checks.MaybeAdd(defaultPass_checkVariableReferences)
		checks.MaybeAdd(defaultPass_checkVars)
	}

	if options.ShouldInferConsts() {
		checks.MaybeAdd(defaultPass_inferConsts)
	}
	if options.ShouldComputeFunctionSideEffects() {
		checks.MaybeAdd(defaultPass_checkRegExp)
	}

	checks.MaybeAdd(emptyBeforeTypeChecking)

	if options.GetCheckTypes() || options.GetInferTypes() {
		checks.MaybeAdd(defaultPass_inferTypes)
		if options.GetCheckTypes() {
			checks.MaybeAdd(defaultPass_checkTypes)
		} else {
			checks.MaybeAdd(defaultPass_inferJsDocInfo)
		}
	}

	if options.Enables(DiagnosticGroups_ANALYZER_CHECKS) && options.IsTypecheckingEnabled() {
		checks.MaybeAdd(defaultPass_analyzerChecks)
	}

	if !options.PreservesDetailedSourceInfo() {
		checks.MaybeAdd(defaultPass_clearTypedScopeCreatorPass)
	}

	if options.ShouldRewriteModulesAfterTypechecking() {
		c.addModuleRewritingPasses(checks, options, &closureRewriteModule)
	}

	if options.ShouldAllowDynamicImport() && inFeatures.Has(parser.FeatureSet_Feature_DYNAMIC_IMPORT) {
		checks.MaybeAdd(defaultPassConfig_rewriteDynamicImports)
	}

	if !options.PreservesDetailedSourceInfo() {
		checks.MaybeAdd(defaultPass_clearTopTypedScopePass)
	}

	if options.GetCheckSuspiciousCode() ||
		options.Enables(DiagnosticGroups_GLOBAL_THIS) ||
		options.Enables(DiagnosticGroups_DEBUGGER_STATEMENT_PRESENT) {
		checks.MaybeAdd(defaultPass_suspiciousCode)
	}

	if options.GetJ2clPass().ShouldAddJ2clPasses() {
		checks.MaybeAdd(defaultPass_j2clSourceFileChecker)
	}

	if !options.Disables(DiagnosticGroups_CHECK_USELESS_CODE) ||
		!options.Disables(DiagnosticGroups_MISSING_RETURN) {
		checks.MaybeAdd(defaultPass_checkControlFlow)
	}

	if options.IsTypecheckingEnabled() &&
		(!options.Disables(DiagnosticGroups_ACCESS_CONTROLS) ||
			options.Enables(DiagnosticGroups_CONSTANT_PROPERTY)) {
		checks.MaybeAdd(defaultPass_checkAccessControls)
	}

	checks.MaybeAdd(defaultPass_checkConsts)
	checks.MaybeAdd(defaultPass_rewriteCallerCodeLocation)

	if len(options.GetConformanceConfigs()) != 0 {
		checks.MaybeAdd(defaultPass_checkConformance)
	}
	if options.GetTweakProcessing().IsOn() {
		checks.MaybeAdd(defaultPass_processTweaks)
	}

	checks.MaybeAdd(defaultPass_processDefinesCheck)

	if options.GetJ2clPass().ShouldAddJ2clPasses() {
		checks.MaybeAdd(defaultPass_j2clChecksPass)
	}
	if options.ShouldGenerateExports() {
		checks.MaybeAdd(defaultPass_generateExports)
	}

	checks.MaybeAdd(emptyAfterStandardChecks)
	checks.MaybeAdd(defaultPass_mergeSyntheticScript)

	if options.GetExternExportsPath() != nil {
		checks.MaybeAdd(defaultPass_externExports)
	}
	if !options.IsChecksOnly() {
		checks.MaybeAdd(defaultPass_removeWeakSources)
	}

	checks.MaybeAdd(defaultPass_gatherExternPropertiesCheck)
	checks.MaybeAdd(emptyBeforeSerialization)

	if options.GetTypedAstOutputFile() != nil {
		checks.MaybeAdd(defaultPass_serializeTypedAst)
	}

	checks.AssertAllOneTimePasses()
	assertValidOrderForChecks(checks, closureRewriteModule, closurePrimitives, closureGoogScopeAliases)
	return checks
}

func (c *DefaultPassConfig) GetOptimizations() *PassListBuilder {
	options := c.options
	passes := NewPassListBuilder(options)

	if options.IsPropertyRenamingOnlyCompilationMode() {
		passes.MaybeAdd(defaultPass_removeUnnecessarySyntheticExterns)
		TranspilationPasses_AddTranspilationRuntimeLibraries(passes)
		passes.MaybeAdd(c.makeClosureProvidesRequires())
		passes.MaybeAdd(defaultPass_processDefinesOptimize)
		passes.MaybeAdd(defaultPass_normalize)
		passes.MaybeAdd(defaultPass_gatherGettersAndSetters)
		TranspilationPasses_AddTranspilationPasses(passes, options)
		passes.MaybeAdd(defaultPass_gatherExternPropertiesOptimize)
		passes.MaybeAdd(emptyBeforeStandardOptimizations)
		passes.MaybeAdd(defaultPass_inlineAndCollapseProperties)
		passes.MaybeAdd(defaultPass_closureOptimizePrimitives)
		passes.MaybeAdd(defaultPass_stripSideEffectProtection)
		return passes
	}

	if options.GetSkipNonTranspilationPasses() {
		return passes
	}

	passes.AddAll(c.getEarlyOptimizationPasses())
	passes.MaybeAdd(emptyOptimizationsHalfwayPoint)
	passes.AddAll(c.getLateOptimizationPasses())
	return passes
}

func validateRenamePrefixNamespace(ns string) {
	if !defaultPassConfig_GLOBAL_SYMBOL_NAMESPACE_PATTERN.Matches(ns) {
		panic(jutil.NewIllegalArgumentExceptionMsg("Illegal character in renamePrefixNamespace name: " + ns))
	}
}

func (c *DefaultPassConfig) GetFinalizations() *PassListBuilder {
	options := c.options
	passes := NewPassListBuilder(options)

	if options.IsPropertyRenamingOnlyCompilationMode() {
		if options.ShouldRewriteGlobalDeclarationsForTryCatchWrapping() || options.GetRenamePrefixNamespace() != nil {
			passes.MaybeAdd(defaultPass_rewriteGlobalDeclarationsForTryCatchWrapping)
		}
		passes.MaybeAdd(defaultPass_renameProperties)
		if ns := options.GetRenamePrefixNamespace(); ns != nil &&
			options.GetChunkOutputType() == CompilerOptions_ChunkOutputType_GLOBAL_NAMESPACE {
			validateRenamePrefixNamespace(*ns)
			passes.MaybeAdd(defaultPass_rescopeGlobalSymbols)
		}
		return passes
	}

	if options.DoLateLocalization() {
		if options.ShouldRunReplaceMessagesPass() {
			passes.MaybeAdd(defaultPass_replaceProtectedMessages)
		}
		passes.MaybeAdd(defaultPass_substituteLocaleData)
		passes.MaybeAdd(defaultPass_peepholeOptimizationsOnceNormalized)
	}

	if options.ShouldInlineVariables() || options.ShouldInlineLocalVariables() {
		passes.MaybeAdd(defaultPass_flowSensitiveInlineVariables)
		if !options.DoLateLocalization() && c.shouldRunRemoveUnusedCode() {
			passes.MaybeAdd(defaultPass_removeUnusedCodeOnce)
		}
	}

	if options.DoLateLocalization() {
		passes.AddAll(c.getPostL10nOptimizations())
	}

	passes.MaybeAdd(emptyBeforeChunkMotion)
	if options.ShouldRunCrossChunkCodeMotion() {
		passes.MaybeAdd(defaultPass_crossChunkCodeMotion)
	}
	if options.ShouldRunCrossChunkMethodMotion() {
		passes.MaybeAdd(defaultPass_crossChunkMethodMotion)
	}
	passes.MaybeAdd(emptyAfterChunkMotion)

	if options.GetOptimizeESClassConstructors() &&
		options.GetOutputFeatureSet().Contains(parser.FeatureSet_ES2015) {
		passes.MaybeAdd(defaultPass_optimizeConstructors)
	}

	if options.GetIsolatePolyfills() {
		passes.MaybeAdd(defaultPass_isolatePolyfills)
	}
	if options.ShouldCollapseAnonymousFunctions() {
		passes.MaybeAdd(defaultPass_collapseAnonymousFunctions)
	}
	if options.ShouldRewriteGlobalDeclarationsForTryCatchWrapping() || options.GetRenamePrefixNamespace() != nil {
		passes.MaybeAdd(defaultPass_rewriteGlobalDeclarationsForTryCatchWrapping)
	}

	passes.MaybeAdd(emptyBeforeExtractPrototypeMemberDeclarations)
	if options.GetExtractPrototypeMemberDeclarationsMode() != CompilerOptions_ExtractPrototypeMemberDeclarationsMode_OFF {
		passes.MaybeAdd(defaultPass_extractPrototypeMemberDeclarations)
	}

	if options.ShouldAmbiguateProperties() &&
		options.GetPropertyRenaming() == PropertyRenamingPolicy_ALL_UNQUOTED &&
		options.IsTypecheckingEnabled() {
		passes.MaybeAdd(defaultPass_ambiguateProperties)
	}

	passes.MaybeAdd(emptyBeforeRenameProperties)
	if options.GetPropertyRenaming() == PropertyRenamingPolicy_ALL_UNQUOTED {
		passes.MaybeAdd(defaultPass_renameProperties)
	} else {
		passes.MaybeAdd(defaultPass_removePropertyRenamingCalls)
	}

	if options.ShouldReserveRawExports() {
		passes.MaybeAdd(defaultPass_gatherRawExports)
	}
	if options.ShouldConvertToDottedProperties() {
		passes.MaybeAdd(defaultPass_convertToDottedProperties)
	}
	if options.ShouldRewriteFunctionExpressions() {
		passes.MaybeAdd(defaultPass_rewriteFunctionExpressions)
	}
	if options.GetAliasStringsMode() != CompilerOptions_AliasStringsMode_NONE {
		passes.MaybeAdd(defaultPass_aliasStrings)
	}

	if options.ShouldCoalesceVariableNames() {
		passes.MaybeAdd(defaultPass_coalesceVariableNames)
		if options.ShouldFoldConstants() {
			passes.MaybeAdd(defaultPass_peepholeOptimizationsOnceNonNormalized)
		}
	}

	passes.MaybeAdd(defaultPass_markUnnormalized)

	if options.ShouldCollapseVariableDeclarations() || options.ShouldOptimizeLetAndConst() {
		passes.MaybeAdd(defaultPass_postNormalizePeephole)
	}
	if options.ShouldCollapseVariableDeclarations() {
		passes.MaybeAdd(defaultPass_collapseVariableDeclarations)
	}

	passes.MaybeAdd(defaultPass_denormalize)
	passes.MaybeAdd(emptyBeforeVariableRenaming)

	if options.GetVariableRenaming() != VariableRenamingPolicy_ALL {
		passes.MaybeAdd(defaultPass_invertContextualRenaming)
	}
	if options.GetVariableRenaming() != VariableRenamingPolicy_OFF {
		passes.MaybeAdd(defaultPass_renameVars)
	}
	if options.ShouldRenameLabels() {
		passes.MaybeAdd(defaultPass_renameLabels)
	}
	if options.ShouldFoldConstants() {
		passes.MaybeAdd(defaultPass_latePeepholeOptimizations)
	}

	if options.ShouldProtectHiddenSideEffects() || options.GetMergedPrecompiledLibraries() {
		passes.MaybeAdd(defaultPass_stripSideEffectProtection)
	}

	if ns := options.GetRenamePrefixNamespace(); ns != nil &&
		options.GetChunkOutputType() == CompilerOptions_ChunkOutputType_GLOBAL_NAMESPACE {
		validateRenamePrefixNamespace(*ns)
		passes.MaybeAdd(defaultPass_rescopeGlobalSymbols)
	}

	if options.GetOutputFeatureSet().Contains(parser.FeatureSet_ES2015) {
		passes.MaybeAdd(defaultPass_optimizeToEs6)
	}
	if options.GetChunkOutputType() == CompilerOptions_ChunkOutputType_ES_MODULES {
		passes.MaybeAdd(defaultPass_convertChunksToESModules)
	}

	passes.MaybeAdd(defaultPass_checkAstValidity)
	passes.MaybeAdd(defaultPass_varCheckValidity)
	return passes
}

func (c *DefaultPassConfig) getEarlyOptimizationLoopPasses() *PassListBuilder {
	options := c.options
	earlyLoopPasses := NewPassListBuilder(options)

	if options.ShouldInlineVariables() || options.ShouldInlineLocalVariables() {
		earlyLoopPasses.MaybeAdd(defaultPass_inlineVariables)
	} else if options.ShouldInlineConstantVars() {
		earlyLoopPasses.MaybeAdd(defaultPass_inlineConstants)
	}
	if options.GetCollapseObjectLiterals() {
		earlyLoopPasses.MaybeAdd(defaultPass_collapseObjectLiterals)
	}
	if c.shouldRunRemoveUnusedCode() {
		earlyLoopPasses.MaybeAdd(defaultPass_removeUnusedCode)
	}
	if options.ShouldFoldConstants() {
		earlyLoopPasses.MaybeAdd(defaultPass_peepholeOptimizations)
	}

	earlyLoopPasses.AssertAllLoopablePasses()
	return earlyLoopPasses
}

func (c *DefaultPassConfig) getPostL10nOptimizations() *PassListBuilder {
	options := c.options
	loopPasses := NewPassListBuilder(options)

	if options.ShouldOptimizeCalls() {
		loopPasses.MaybeAdd(defaultPass_optimizeCalls)
	}
	if options.GetJ2clPass().ShouldAddJ2clPasses() {
		loopPasses.MaybeAdd(defaultPass_j2clConstantHoisterPass)
		loopPasses.MaybeAdd(defaultPass_j2clClinitPass)
	}
	if options.GetInlineFunctionsLevel() != CompilerOptions_Reach_NONE {
		loopPasses.MaybeAdd(defaultPass_inlineFunctions)
	}
	if options.ShouldInlineVariables() || options.ShouldInlineLocalVariables() {
		loopPasses.MaybeAdd(defaultPass_inlineVariables)
	} else if options.ShouldInlineConstantVars() {
		loopPasses.MaybeAdd(defaultPass_inlineConstants)
	}
	if c.shouldRunRemoveUnusedCode() {
		loopPasses.MaybeAdd(defaultPass_removeUnusedCode)
	}
	if options.ShouldFoldConstants() {
		loopPasses.MaybeAdd(defaultPass_peepholeOptimizations)
	}

	loopPasses.AssertAllLoopablePasses()
	return loopPasses
}

func (c *DefaultPassConfig) getEarlyOptimizationPasses() *PassListBuilder {
	options := c.options
	passes := NewPassListBuilder(options)

	if options.ShouldExportTestFunctions() {
		passes.MaybeAdd(defaultPass_exportTestFunctions)
	}

	if options.GetMergedPrecompiledLibraries() {
		passes.MaybeAdd(defaultPass_removeWeakSources)
		passes.MaybeAdd(defaultPass_checkRegExpForOptimizations)
		if options.GetJ2clPass().ShouldAddJ2clPasses() {
			passes.MaybeAdd(defaultPass_j2clSourceFileChecker)
		}
	} else {
		addNonTypedAstNormalizationPasses(passes)
	}

	passes.MaybeAdd(defaultPass_removeUnnecessarySyntheticExterns)

	if options.GetSyntheticBlockStartMarker() != nil {
		passes.MaybeAdd(defaultPass_createSyntheticBlocks)
	}
	if options.GetJ2clPass().ShouldAddJ2clPasses() {
		passes.MaybeAdd(defaultPass_j2clPass)
	}

	passes.MaybeAdd(defaultPass_transpileAndOptimizeClosureUnaware)

	TranspilationPasses_AddTranspilationRuntimeLibraries(passes)

	if options.GetRewritePolyfills() || options.GetIsolatePolyfills() || options.GetInjectPolyfillsNewerThan() != nil {
		TranspilationPasses_AddRewritePolyfillPass(passes)
	}

	passes.MaybeAdd(defaultPass_injectRuntimeLibraries)

	if options.GetClosurePass() {
		passes.MaybeAdd(c.makeClosureProvidesRequires())
	}

	if options.ShouldRunReplaceMessagesForChrome() {
		passes.MaybeAdd(defaultPass_replaceMessagesForChrome)
	} else if options.ShouldRunReplaceMessagesPass() {
		if options.DoLateLocalization() {
			passes.MaybeAdd(defaultPass_protectMessages)
		} else {
			passes.MaybeAdd(defaultPass_fullReplaceMessages)
		}
	}

	if options.DoLateLocalization() {
		passes.MaybeAdd(defaultPass_protectLocaleData)
	}

	if options.GetClosurePass() && !options.ShouldPreserveGoogLibraryPrimitives() {
		passes.MaybeAdd(defaultPass_closureReplaceGetCssName)
	}

	passes.MaybeAdd(defaultPass_processDefinesOptimize)
	passes.MaybeAdd(emptyBeforeEarlyOptimizationsTranspilation)
	passes.MaybeAdd(defaultPass_normalize)
	passes.MaybeAdd(defaultPass_gatherGettersAndSetters)

	TranspilationPasses_AddTranspilationPasses(passes, options)

	if options.GetJ2clPass().ShouldAddJ2clPasses() {
		passes.MaybeAdd(defaultPass_j2clUtilGetDefineRewriterPass)
	}
	if options.GetInstrumentForCoverageOption() != CompilerOptions_InstrumentOption_NONE {
		passes.MaybeAdd(defaultPass_instrumentForCodeCoverage)
	}

	passes.MaybeAdd(defaultPass_gatherExternPropertiesOptimize)
	passes.MaybeAdd(emptyBeforeStandardOptimizations)

	if options.GetClosurePass() && (options.ShouldRemoveAbstractMethods() || options.ShouldRemoveClosureAsserts()) {
		passes.MaybeAdd(defaultPass_closureCodeRemoval)
	}
	if options.ShouldRemoveJ2clAsserts() {
		passes.MaybeAdd(defaultPass_j2clAssertRemovalPass)
	}

	passes.MaybeAdd(defaultPass_replaceToggles)
	passes.MaybeAdd(defaultPass_inlineAndCollapseProperties)

	if options.GetTweakProcessing().ShouldStrip() ||
		!options.GetStripTypes().IsEmpty() ||
		!options.GetStripNameSuffixes().IsEmpty() ||
		!options.GetStripNamePrefixes().IsEmpty() {
		passes.MaybeAdd(defaultPass_stripCode)
	}

	if options.ShouldReplaceIdGenerators() {
		passes.MaybeAdd(defaultPass_replaceIdGenerators)
	}

	if options.GetJ2clPass().ShouldAddJ2clPasses() &&
		options.GetPropertyCollapseLevel() == CompilerOptions_PropertyCollapseLevel_ALL {
		passes.MaybeAdd(defaultPass_j2clPropertyInlinerPass)
	}

	if options.ShouldInferConsts() {
		passes.MaybeAdd(defaultPass_inferConsts)
	}

	passes.MaybeAdd(emptyObfuscationPassMarker)

	if options.GetSmartNameRemoval() {
		if options.ShouldFoldConstants() &&
			(options.ShouldInlineVariables() || options.ShouldInlineLocalVariables()) {
			passes.MaybeAdd(defaultPass_earlyInlineVariables)
			passes.MaybeAdd(defaultPass_earlyPeepholeOptimizations)
		}
		passes.MaybeAdd(defaultPass_removeUnusedCodeOnce)
	}

	if options.ShouldDisambiguateProperties() && options.IsTypecheckingEnabled() {
		passes.MaybeAdd(defaultPass_disambiguateProperties)
	}
	if options.ShouldComputeFunctionSideEffects() {
		passes.MaybeAdd(defaultPass_markPureFunctions)
	}

	passes.AssertAllOneTimePasses()
	return passes
}

func (c *DefaultPassConfig) getLateOptimizationPasses() *PassListBuilder {
	options := c.options
	passes := NewPassListBuilder(options)

	if options.GetSmartNameRemoval() {
		passes.MaybeAdd(emptyBeforeEarlyOptimizationLoop)
		passes.AddAll(c.getEarlyOptimizationLoopPasses())
		passes.MaybeAdd(emptyAfterEarlyOptimizationLoop)
	}

	if options.GetClosurePass() {
		passes.MaybeAdd(defaultPass_closureOptimizePrimitives)
	}
	if len(options.GetReplaceStringsFunctionDescriptions()) != 0 {
		passes.MaybeAdd(defaultPass_replaceStrings)
	}
	if options.ShouldRunCrossChunkCodeMotion() {
		passes.MaybeAdd(defaultPass_crossChunkCodeMotion)
	}
	if options.ShouldDevirtualizeMethods() {
		passes.MaybeAdd(defaultPass_devirtualizeMethods)
	}

	passes.MaybeAdd(defaultPass_customPassesBeforeOptimizationLoop)
	passes.MaybeAdd(emptyBeforeMainOptimizations)

	if options.ShouldInlineVariables() || options.ShouldInlineLocalVariables() {
		passes.MaybeAdd(defaultPass_flowSensitiveInlineVariables)
	}

	passes.AddAll(c.getMainOptimizationLoop())
	passes.MaybeAdd(emptyAfterMainOptimizations)
	passes.MaybeAdd(defaultPass_customPassesAfterOptimizationLoop)

	assertValidOrderForOptimizations(passes)
	return passes
}

func (c *DefaultPassConfig) getMainOptimizationLoop() *PassListBuilder {
	options := c.options
	passes := NewPassListBuilder(options)

	if options.ShouldInlineGetters() {
		passes.MaybeAdd(defaultPass_inlineSimpleMethods)
	}
	if options.ShouldInlineProperties() && options.IsTypecheckingEnabled() {
		passes.MaybeAdd(defaultPass_inlineProperties)
	}
	if options.ShouldRunDeadPropertyAssignmentElimination() {
		passes.MaybeAdd(defaultPass_deadPropertyAssignmentElimination)
	}
	if options.ShouldOptimizeCalls() {
		passes.MaybeAdd(defaultPass_optimizeCalls)
	}
	if options.GetJ2clPass().ShouldAddJ2clPasses() {
		passes.MaybeAdd(defaultPass_j2clConstantHoisterPass)
		passes.MaybeAdd(defaultPass_j2clClinitPass)
	}
	// inlineVariables and peepholeOptimizations must run after inlineFunctions to clean up patterns
	// it introduces (which affects the size-based loop-termination heuristic).
	if options.GetInlineFunctionsLevel() != CompilerOptions_Reach_NONE {
		passes.MaybeAdd(defaultPass_inlineFunctions)
	}
	if options.ShouldInlineVariables() || options.ShouldInlineLocalVariables() {
		passes.MaybeAdd(defaultPass_inlineVariables)
	} else if options.ShouldInlineConstantVars() {
		passes.MaybeAdd(defaultPass_inlineConstants)
	}
	if options.ShouldRunDeadAssignmentElimination() {
		passes.MaybeAdd(defaultPass_deadAssignmentsElimination)
	}
	if options.GetCollapseObjectLiterals() {
		passes.MaybeAdd(defaultPass_collapseObjectLiterals)
	}
	if c.shouldRunRemoveUnusedCode() {
		passes.MaybeAdd(defaultPass_removeUnusedCode)
	}
	if options.ShouldFoldConstants() {
		passes.MaybeAdd(defaultPass_peepholeOptimizations)
	}

	passes.AssertAllLoopablePasses()
	return passes
}

func addNonTypedAstNormalizationPasses(passes *PassListBuilder) {
	passes.MaybeAdd(defaultPass_removeCastNodes)
	passes.MaybeAdd(defaultPass_typesToColors)
}

func (c *DefaultPassConfig) shouldRunRemoveUnusedCode() bool {
	options := c.options
	return options.ShouldRemoveUnusedVariables() ||
		options.ShouldRemoveUnusedLocalVariables() ||
		options.ShouldRemoveUnusedPrototypeProperties() ||
		options.IsRemoveUnusedClassProperties() ||
		options.GetRewritePolyfills()
}

func assertValidOrderForChecks(
	checks *PassListBuilder,
	closureRewriteModule *PassFactory,
	closurePrimitives *PassFactory,
	closureGoogScopeAliases *PassFactory,
) {
	checks.AssertPassOrder(
		defaultPass_declaredGlobalExternsOnWindow,
		defaultPass_checkVars,
		"declaredGlobalExternsOnWindow must happen before VarCheck, which adds synthetic externs")
	checks.AssertPassOrder(
		defaultPass_chromePass,
		defaultPass_checkJsDocAndEs6Modules,
		"The ChromePass must run before after JsDoc and Es6 module checking.")
	if closureRewriteModule != nil {
		checks.AssertPassOrder(
			closureRewriteModule,
			defaultPass_processDefinesCheck,
			"Must rewrite goog.module before processing @define's, so that @defines in modules work.")
	}
	if closurePrimitives != nil {
		checks.AssertPassOrder(
			closurePrimitives, defaultPass_polymerPass, "The Polymer pass must run after goog.provide processing.")
	}
	checks.AssertPassOrder(
		defaultPass_chromePass, defaultPass_polymerPass, "The Polymer pass must run after ChromePass processing.")
	checks.AssertPassOrder(
		defaultPass_polymerPass, defaultPass_suspiciousCode, "The Polymer pass must run before suspiciousCode processing.")
	checks.AssertPassOrder(
		defaultPass_addSyntheticScript,
		defaultPass_gatherModuleMetadataPass,
		"Cannot add a synthetic script node after module metadata creation.")
	if closureRewriteModule != nil {
		checks.AssertPassOrder(
			closureRewriteModule,
			defaultPass_removeSyntheticScript,
			"Synthetic script node should be removed only after module rewriting.")
		checks.AssertPassOrder(
			closureRewriteModule,
			defaultPass_rewriteCallerCodeLocation,
			"ClosureRewriteModule must happen before RewriteCallerCodeLocation, so that exported"+
				" functions and call sites are rewritten correctly.")
		checks.AssertPassOrder(
			closureRewriteModule,
			TranspilationPasses_GetEs6RewriteDestructuring(
				Es6RewriteDestructuring_ObjectDestructuringRewriteMode_REWRITE_ALL_OBJECT_PATTERNS),
			"RewriteCallerCodeLocation must happen before Es6RewriteDestructuring, because we need"+
				" ReWriteCallerCodeLocation to run before default parameters get rewritten.")
		checks.AssertPassOrder(
			closureRewriteModule,
			TranspilationPasses_GetEs6RewriteDestructuring(
				Es6RewriteDestructuring_ObjectDestructuringRewriteMode_REWRITE_OBJECT_REST),
			"RewriteCallerCodeLocation must happen before Es6RewriteDestructuring, because we need"+
				" ReWriteCallerCodeLocation to run before default parameters get rewritten.")
	}

	if closureGoogScopeAliases != nil && checks.Contains(closureGoogScopeAliases) {
		jutil.CheckStateMsg(
			checks.Contains(defaultPass_checkVariableReferences),
			"goog.scope processing requires variable checking")
	}

	checks.AssertPassOrder(
		defaultPass_gatherModuleMetadataPass,
		defaultPass_closureCheckModule,
		"Need to gather module metadata before checking closure modules.")
	checks.AssertPassOrder(
		defaultPass_gatherModuleMetadataPass,
		defaultPass_createModuleMapPass,
		"Need to gather module metadata before scanning modules.")
	checks.AssertPassOrder(
		defaultPass_createModuleMapPass,
		defaultPass_rewriteCommonJsModules,
		"Need to gather module information before rewriting CommonJS modules.")
	checks.AssertPassOrder(
		defaultPass_rewriteScriptsToEs6Modules,
		defaultPass_gatherModuleMetadataPass,
		"Need to gather module information after rewriting scripts to modules.")
	checks.AssertPassOrder(
		defaultPass_gatherModuleMetadataPass,
		defaultPass_checkMissingRequires,
		"Need to gather module information before checking for missing requires.")
}

func assertValidOrderForOptimizations(optimizations *PassListBuilder) {
	optimizations.AssertPassOrder(
		defaultPass_j2clPass,
		TranspilationPasses_rewriteGenerators,
		"J2CL normalization should be done before generator re-writing.")
	optimizations.AssertPassOrder(
		TranspilationPasses_rewritePolyfills,
		defaultPass_processDefinesOptimize,
		"Polyfill injection must be done before processDefines as some polyfills reference goog.defines.")
	optimizations.AssertPassOrder(
		TranspilationPasses_injectTranspilationRuntimeLibraries,
		defaultPass_processDefinesOptimize,
		"Runtime library injection must be done before processDefines some runtime libraries reference goog.defines.")
	optimizations.AssertPassOrder(
		defaultPass_processDefinesOptimize,
		defaultPass_j2clUtilGetDefineRewriterPass,
		"J2CL define re-writing should be done after processDefines since it relies on Compiler#getDefineNames to have been populated by it.")
	optimizations.AssertPassOrder(
		defaultPass_removeUnusedCode,
		defaultPass_isolatePolyfills,
		"Polyfill isolation should be done after RemovedUnusedCode. Otherwise unused polyfill removal will not find any polyfill usages and will delete all polyfills.")
	optimizations.AssertPassOrder(
		TranspilationPasses_instrumentAsyncContext,
		TranspilationPasses_rewriteAsyncIteration,
		"AsyncContext should be instrumentated before await and/or yield is transpiled away")
	optimizations.AssertPassOrder(
		TranspilationPasses_instrumentAsyncContext,
		TranspilationPasses_rewriteAsyncFunctions,
		"AsyncContext should be instrumentated before await and/or yield is transpiled away")
	optimizations.AssertPassOrder(
		TranspilationPasses_instrumentAsyncContext,
		TranspilationPasses_rewriteGenerators,
		"AsyncContext should be instrumentated before await and/or yield is transpiled away")
}
