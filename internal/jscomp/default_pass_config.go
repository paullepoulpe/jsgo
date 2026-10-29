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

	checkSideEffects                             *PassFactory
	stripSideEffectProtection                    *PassFactory
	suspiciousCode                               *PassFactory
	extraRequires                                *PassFactory
	checkMissingRequires                         *PassFactory
	checkJsDocAndEs6Modules                      *PassFactory
	generateExports                              *PassFactory
	generateIjs                                  *PassFactory
	removeExtraRequires                          *PassFactory
	exportTestFunctions                          *PassFactory
	gatherRawExports                             *PassFactory
	closurePrimitives                            *PassFactory
	closureProvidesRequires                      *PassFactory
	angularPass                                  *PassFactory
	replaceMessagesForChrome                     *PassFactory
	closureGoogScopeAliasesForIjs                *PassFactory
	closureGoogScopeAliases                      *PassFactory
	injectRuntimeLibraries                       *PassFactory
	removeWeakSources                            *PassFactory
	declaredGlobalExternsOnWindow                *PassFactory
	checkTypeImportCodeReferences                *PassFactory
	closureCheckModule                           *PassFactory
	closureRewriteModule                         *PassFactory
	checkClosureImports                          *PassFactory
	rewriteGoogJsImports                         *PassFactory
	closureReplaceGetCssName                     *PassFactory
	createSyntheticBlocks                        *PassFactory
	earlyPeepholeOptimizations                   *PassFactory
	earlyInlineVariables                         *PassFactory
	peepholeOptimizations                        *PassFactory
	peepholeOptimizationsOnceNormalized          *PassFactory
	peepholeOptimizationsOnceNonNormalized       *PassFactory
	latePeepholeOptimizations                    *PassFactory
	checkVars                                    *PassFactory
	inferConsts                                  *PassFactory
	checkRegExp                                  *PassFactory
	checkRegExpForOptimizations                  *PassFactory
	checkVariableReferences                      *PassFactory
	checkSuper                                   *PassFactory
	clearTypedScopeCreatorPass                   *PassFactory
	clearTopTypedScopePass                       *PassFactory
	inferTypes                                   *PassFactory
	inferJsDocInfo                               *PassFactory
	checkTypes                                   *PassFactory
	checkControlFlow                             *PassFactory
	checkAccessControls                          *PassFactory
	lintChecks                                   *PassFactory
	analyzerChecks                               *PassFactory
	checkRequiresAndProvidesSorted               *PassFactory
	checkStrictMode                              *PassFactory
	processTweaks                                *PassFactory
	processDefinesCheck                          *PassFactory
	processDefinesOptimize                       *PassFactory
	stripCode                                    *PassFactory
	checkConsts                                  *PassFactory
	rewriteCallerCodeLocation                    *PassFactory
	replaceToggles                               *PassFactory
	replaceIdGenerators                          *PassFactory
	replaceStrings                               *PassFactory
	closureCodeRemoval                           *PassFactory
	closureOptimizePrimitives                    *PassFactory
	rescopeGlobalSymbols                         *PassFactory
	convertChunksToESModules                     *PassFactory
	inlineAndCollapseProperties                  *PassFactory
	collapseObjectLiterals                       *PassFactory
	disambiguateProperties                       *PassFactory
	devirtualizeMethods                          *PassFactory
	optimizeCalls                                *PassFactory
	optimizeConstructors                         *PassFactory
	markPureFunctions                            *PassFactory
	inlineVariables                              *PassFactory
	inlineConstants                              *PassFactory
	inlineSimpleMethods                          *PassFactory
	deadAssignmentsElimination                   *PassFactory
	deadPropertyAssignmentElimination            *PassFactory
	inlineFunctions                              *PassFactory
	inlineProperties                             *PassFactory
	isolatePolyfills                             *PassFactory
	removeUnusedCode                             *PassFactory
	removeUnusedCodeOnce                         *PassFactory
	crossChunkCodeMotion                         *PassFactory
	crossChunkMethodMotion                       *PassFactory
	flowSensitiveInlineVariables                 *PassFactory
	coalesceVariableNames                        *PassFactory
	postNormalizePeephole                        *PassFactory
	collapseVariableDeclarations                 *PassFactory
	extractPrototypeMemberDeclarations           *PassFactory
	rewriteFunctionExpressions                   *PassFactory
	collapseAnonymousFunctions                   *PassFactory
	rewriteGlobalDeclarationsForTryCatchWrapping *PassFactory
	aliasStrings                                 *PassFactory
	ambiguateProperties                          *PassFactory
	markUnnormalized                             *PassFactory
	normalize                                    *PassFactory
	externExports                                *PassFactory
	denormalize                                  *PassFactory
	invertContextualRenaming                     *PassFactory
	renameProperties                             *PassFactory
	removePropertyRenamingCalls                  *PassFactory
	renameVars                                   *PassFactory
	renameLabels                                 *PassFactory
	convertToDottedProperties                    *PassFactory
	checkAstValidity                             *PassFactory
	varCheckValidity                             *PassFactory
	instrumentForCodeCoverage                    *PassFactory
	gatherExternPropertiesCheck                  *PassFactory
	gatherExternPropertiesOptimize               *PassFactory
	polymerPass                                  *PassFactory
	chromePass                                   *PassFactory
	j2clConstantHoisterPass                      *PassFactory
	j2clClinitPass                               *PassFactory
	j2clPropertyInlinerPass                      *PassFactory
	j2clPass                                     *PassFactory
	j2clUtilGetDefineRewriterPass                *PassFactory
	j2clAssertRemovalPass                        *PassFactory
	j2clSourceFileChecker                        *PassFactory
	j2clChecksPass                               *PassFactory
	checkConformance                             *PassFactory
	removeCastNodes                              *PassFactory
	typesToColors                                *PassFactory
	serializeTypedAst                            *PassFactory
	removeUnnecessarySyntheticExterns            *PassFactory
	optimizeToEs6                                *PassFactory
	whitespaceWrapGoogModules                    *PassFactory
	rewriteCommonJsModules                       *PassFactory
	rewriteScriptsToEs6Modules                   *PassFactory
	gatherModuleMetadataPass                     *PassFactory
	createModuleMapPass                          *PassFactory
	gatherGettersAndSetters                      *PassFactory
	addSyntheticScript                           *PassFactory
	removeSyntheticScript                        *PassFactory
	mergeSyntheticScript                         *PassFactory
	protectLocaleData                            *PassFactory
	substituteLocaleData                         *PassFactory
	transpileOnlyClosureUnaware                  *PassFactory
	transpileAndOptimizeClosureUnaware           *PassFactory
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
	c := &DefaultPassConfig{options: options}
	c.initFields()
	return c
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
	passes.MaybeAdd(c.checkVariableReferences)
	passes.MaybeAdd(c.checkVars)

	passes.MaybeAdd(c.gatherModuleMetadataPass)
	passes.MaybeAdd(c.createModuleMapPass)

	if options.GetLanguageIn().ToFeatureSet().Has(parser.FeatureSet_Feature_MODULES) {
		passes.MaybeAdd(c.rewriteGoogJsImports)
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

	passes.MaybeAdd(c.checkSuper)

	TranspilationPasses_AddTranspilationRuntimeLibraries(passes)

	passes.MaybeAdd(c.transpileOnlyClosureUnaware)

	if options.NeedsTranspilationFrom(parser.FeatureSet_ES2015) && options.GetRewritePolyfills() {
		if options.GetIsolatePolyfills() {
			panic(jutil.NewIllegalStateExceptionMsg("Polyfill isolation cannot be used in transpileOnly mode"))
		}
		TranspilationPasses_AddRewritePolyfillPass(passes)
	} else if options.GetInjectPolyfillsNewerThan() != nil {
		TranspilationPasses_AddRewritePolyfillPass(passes)
	}

	passes.MaybeAdd(c.injectRuntimeLibraries)

	// Passes below this point rely on normalization and must maintain it.
	passes.MaybeAdd(c.normalize)
	passes.MaybeAdd(c.gatherGettersAndSetters)

	TranspilationPasses_AddTranspilationPasses(passes, options)

	passes.MaybeAdd(c.markUnnormalized)
	// Restore original variable/parameter names (e.g. for frameworks doing DI by parameter name).
	passes.MaybeAdd(c.invertContextualRenaming)
	// Es6ConvertSuper may emit property renaming calls that must be stripped in transpile-only mode.
	passes.MaybeAdd(c.removePropertyRenamingCalls)

	passes.AssertAllOneTimePasses()
	c.assertValidOrderForChecks(passes)
	return passes
}

func (c *DefaultPassConfig) GetWhitespaceOnlyPasses() *PassListBuilder {
	options := c.options
	passes := NewPassListBuilder(options)

	if options.GetProcessCommonJSModules() {
		passes.MaybeAdd(c.rewriteCommonJsModules)
	} else if options.GetLanguageIn().ToFeatureSet().Has(parser.FeatureSet_Feature_MODULES) {
		passes.MaybeAdd(c.rewriteScriptsToEs6Modules)
	}

	if options.ShouldWrapGoogModulesForWhitespaceOnly() {
		passes.MaybeAdd(c.whitespaceWrapGoogModules)
	}
	return passes
}

func (c *DefaultPassConfig) addModuleRewritingPasses(checks *PassListBuilder, options *CompilerOptions) {
	if options.GetLanguageIn().ToFeatureSet().Has(parser.FeatureSet_Feature_MODULES) {
		checks.MaybeAdd(c.rewriteGoogJsImports)
		TranspilationPasses_AddEs6ModulePass(checks, c.preprocessorSymbolTableFactory)
	}
	if options.GetClosurePass() {
		checks.MaybeAdd(c.closureRewriteModule)
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
		checks.MaybeAdd(c.addSyntheticScript)
		checks.MaybeAdd(c.gatherGettersAndSetters)
		checks.MaybeAdd(c.gatherModuleMetadataPass)
		checks.MaybeAdd(c.createModuleMapPass)
		checks.MaybeAdd(c.declaredGlobalExternsOnWindow)
		checks.MaybeAdd(c.checkSideEffects)
		checks.MaybeAdd(c.angularPass)
		checks.MaybeAdd(c.closureGoogScopeAliases)
		checks.MaybeAdd(c.closurePrimitives)
		c.addModuleRewritingPasses(checks, options)
		checks.MaybeAdd(c.clearTypedScopeCreatorPass)
		checks.MaybeAdd(c.clearTopTypedScopePass)
		checks.MaybeAdd(c.generateExports)
		checks.MaybeAdd(emptyAfterStandardChecks)
		checks.MaybeAdd(c.mergeSyntheticScript)
		checks.MaybeAdd(c.gatherExternPropertiesCheck)
		checks.MaybeAdd(emptyBeforeSerialization)
		return checks
	}

	if options.ShouldGenerateTypedExterns() {
		checks.MaybeAdd(c.addSyntheticScript)
		checks.MaybeAdd(c.closureGoogScopeAliasesForIjs)
		checks.MaybeAdd(c.generateIjs)
		checks.MaybeAdd(c.removeExtraRequires)
		if options.ShouldWrapGoogModulesForWhitespaceOnly() {
			checks.MaybeAdd(c.whitespaceWrapGoogModules)
		}
		checks.MaybeAdd(c.removeSyntheticScript)
		return checks
	}

	// Run before any pass tries to inject runtime libraries.
	checks.MaybeAdd(c.addSyntheticScript)
	checks.MaybeAdd(c.gatherGettersAndSetters)

	inFeatures := options.GetLanguageIn().ToFeatureSet()
	if inFeatures.ContainsFeature(parser.FeatureSet_Feature_DYNAMIC_IMPORT) && !options.ShouldAllowDynamicImport() {
		checks.MaybeAdd(defaultPassConfig_forbidDynamicImportUsage)
	}

	checks.MaybeAdd(emptyBeforeStandardChecks)

	if !options.GetProcessCommonJSModules() && inFeatures.Has(parser.FeatureSet_Feature_MODULES) {
		checks.MaybeAdd(c.rewriteScriptsToEs6Modules)
	}

	// Run after promoting scripts to modules, but before rewriting any other modules.
	checks.MaybeAdd(c.gatherModuleMetadataPass)
	checks.MaybeAdd(c.createModuleMapPass)

	if options.GetProcessCommonJSModules() {
		checks.MaybeAdd(c.rewriteCommonJsModules)
	}

	// ChromePass can rewrite invalid @type annotations into valid ones, so run before JSDoc checks.
	if options.IsChromePassEnabled() {
		checks.MaybeAdd(c.chromePass)
	}

	checks.MaybeAdd(c.checkJsDocAndEs6Modules)
	checks.MaybeAdd(c.checkTypeImportCodeReferences)

	if options.Enables(DiagnosticGroups_LINT_CHECKS) {
		checks.MaybeAdd(c.lintChecks)
		if options.GetClosurePass() {
			checks.MaybeAdd(c.checkRequiresAndProvidesSorted)
		}
	}

	if options.Enables(DiagnosticGroups_EXTRA_REQUIRE) {
		checks.MaybeAdd(c.extraRequires)
	}
	if options.Enables(DiagnosticGroups_MISSING_REQUIRE) {
		checks.MaybeAdd(c.checkMissingRequires)
	}

	checks.MaybeAdd(c.declaredGlobalExternsOnWindow)

	if !options.GetProcessCommonJSModules() {
		checks.MaybeAdd(c.checkVariableReferences)
		checks.MaybeAdd(c.checkVars)
	}

	if options.GetClosurePass() {
		checks.MaybeAdd(c.checkClosureImports)
	}

	checks.MaybeAdd(c.checkStrictMode)

	if options.GetClosurePass() {
		checks.MaybeAdd(c.closureCheckModule)
	}

	checks.MaybeAdd(c.checkSuper)
	checks.MaybeAdd(c.checkSideEffects)

	if options.GetAngularPass() {
		checks.MaybeAdd(c.angularPass)
	}

	if options.GetClosurePass() {
		checks.MaybeAdd(c.closureGoogScopeAliases)
		checks.MaybeAdd(c.closurePrimitives)
	}

	if options.ShouldRewriteModulesBeforeTypechecking() {
		c.addModuleRewritingPasses(checks, options)
	}

	// PolymerPass must run after ClosurePrimitives/ChromePass and before suspiciousCode checks.
	if options.IsPolymerPassEnabled() {
		checks.MaybeAdd(c.polymerPass)
	}

	if options.GetProcessCommonJSModules() {
		checks.MaybeAdd(c.checkVariableReferences)
		checks.MaybeAdd(c.checkVars)
	}

	if options.ShouldInferConsts() {
		checks.MaybeAdd(c.inferConsts)
	}
	if options.ShouldComputeFunctionSideEffects() {
		checks.MaybeAdd(c.checkRegExp)
	}

	checks.MaybeAdd(emptyBeforeTypeChecking)

	if options.GetCheckTypes() || options.GetInferTypes() {
		checks.MaybeAdd(c.inferTypes)
		if options.GetCheckTypes() {
			checks.MaybeAdd(c.checkTypes)
		} else {
			checks.MaybeAdd(c.inferJsDocInfo)
		}
	}

	if options.Enables(DiagnosticGroups_ANALYZER_CHECKS) && options.IsTypecheckingEnabled() {
		checks.MaybeAdd(c.analyzerChecks)
	}

	if !options.PreservesDetailedSourceInfo() {
		checks.MaybeAdd(c.clearTypedScopeCreatorPass)
	}

	if options.ShouldRewriteModulesAfterTypechecking() {
		c.addModuleRewritingPasses(checks, options)
	}

	if options.ShouldAllowDynamicImport() && inFeatures.Has(parser.FeatureSet_Feature_DYNAMIC_IMPORT) {
		checks.MaybeAdd(defaultPassConfig_rewriteDynamicImports)
	}

	if !options.PreservesDetailedSourceInfo() {
		checks.MaybeAdd(c.clearTopTypedScopePass)
	}

	if options.GetCheckSuspiciousCode() ||
		options.Enables(DiagnosticGroups_GLOBAL_THIS) ||
		options.Enables(DiagnosticGroups_DEBUGGER_STATEMENT_PRESENT) {
		checks.MaybeAdd(c.suspiciousCode)
	}

	if options.GetJ2clPass().ShouldAddJ2clPasses() {
		checks.MaybeAdd(c.j2clSourceFileChecker)
	}

	if !options.Disables(DiagnosticGroups_CHECK_USELESS_CODE) ||
		!options.Disables(DiagnosticGroups_MISSING_RETURN) {
		checks.MaybeAdd(c.checkControlFlow)
	}

	if options.IsTypecheckingEnabled() &&
		(!options.Disables(DiagnosticGroups_ACCESS_CONTROLS) ||
			options.Enables(DiagnosticGroups_CONSTANT_PROPERTY)) {
		checks.MaybeAdd(c.checkAccessControls)
	}

	checks.MaybeAdd(c.checkConsts)
	checks.MaybeAdd(c.rewriteCallerCodeLocation)

	if len(options.GetConformanceConfigs()) != 0 {
		checks.MaybeAdd(c.checkConformance)
	}
	if options.GetTweakProcessing().IsOn() {
		checks.MaybeAdd(c.processTweaks)
	}

	checks.MaybeAdd(c.processDefinesCheck)

	if options.GetJ2clPass().ShouldAddJ2clPasses() {
		checks.MaybeAdd(c.j2clChecksPass)
	}
	if options.ShouldGenerateExports() {
		checks.MaybeAdd(c.generateExports)
	}

	checks.MaybeAdd(emptyAfterStandardChecks)
	checks.MaybeAdd(c.mergeSyntheticScript)

	if options.GetExternExportsPath() != nil {
		checks.MaybeAdd(c.externExports)
	}
	if !options.IsChecksOnly() {
		checks.MaybeAdd(c.removeWeakSources)
	}

	checks.MaybeAdd(c.gatherExternPropertiesCheck)
	checks.MaybeAdd(emptyBeforeSerialization)

	if options.GetTypedAstOutputFile() != nil {
		checks.MaybeAdd(c.serializeTypedAst)
	}

	checks.AssertAllOneTimePasses()
	c.assertValidOrderForChecks(checks)
	return checks
}

func (c *DefaultPassConfig) GetOptimizations() *PassListBuilder {
	options := c.options
	passes := NewPassListBuilder(options)

	if options.IsPropertyRenamingOnlyCompilationMode() {
		passes.MaybeAdd(c.removeUnnecessarySyntheticExterns)
		TranspilationPasses_AddTranspilationRuntimeLibraries(passes)
		passes.MaybeAdd(c.closureProvidesRequires)
		passes.MaybeAdd(c.processDefinesOptimize)
		passes.MaybeAdd(c.normalize)
		passes.MaybeAdd(c.gatherGettersAndSetters)
		TranspilationPasses_AddTranspilationPasses(passes, options)
		passes.MaybeAdd(c.gatherExternPropertiesOptimize)
		passes.MaybeAdd(emptyBeforeStandardOptimizations)
		passes.MaybeAdd(c.inlineAndCollapseProperties)
		passes.MaybeAdd(c.closureOptimizePrimitives)
		passes.MaybeAdd(c.stripSideEffectProtection)
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
			passes.MaybeAdd(c.rewriteGlobalDeclarationsForTryCatchWrapping)
		}
		passes.MaybeAdd(c.renameProperties)
		if ns := options.GetRenamePrefixNamespace(); ns != nil &&
			options.GetChunkOutputType() == CompilerOptions_ChunkOutputType_GLOBAL_NAMESPACE {
			validateRenamePrefixNamespace(*ns)
			passes.MaybeAdd(c.rescopeGlobalSymbols)
		}
		return passes
	}

	if options.DoLateLocalization() {
		if options.ShouldRunReplaceMessagesPass() {
			passes.MaybeAdd(c.getReplaceProtectedMessagesPass())
		}
		passes.MaybeAdd(c.substituteLocaleData)
		passes.MaybeAdd(c.peepholeOptimizationsOnceNormalized)
	}

	if options.ShouldInlineVariables() || options.ShouldInlineLocalVariables() {
		passes.MaybeAdd(c.flowSensitiveInlineVariables)
		if !options.DoLateLocalization() && c.shouldRunRemoveUnusedCode() {
			passes.MaybeAdd(c.removeUnusedCodeOnce)
		}
	}

	if options.DoLateLocalization() {
		passes.AddAll(c.getPostL10nOptimizations())
	}

	passes.MaybeAdd(emptyBeforeChunkMotion)
	if options.ShouldRunCrossChunkCodeMotion() {
		passes.MaybeAdd(c.crossChunkCodeMotion)
	}
	if options.ShouldRunCrossChunkMethodMotion() {
		passes.MaybeAdd(c.crossChunkMethodMotion)
	}
	passes.MaybeAdd(emptyAfterChunkMotion)

	if options.GetOptimizeESClassConstructors() &&
		options.GetOutputFeatureSet().Contains(parser.FeatureSet_ES2015) {
		passes.MaybeAdd(c.optimizeConstructors)
	}

	if options.GetIsolatePolyfills() {
		passes.MaybeAdd(c.isolatePolyfills)
	}
	if options.ShouldCollapseAnonymousFunctions() {
		passes.MaybeAdd(c.collapseAnonymousFunctions)
	}
	if options.ShouldRewriteGlobalDeclarationsForTryCatchWrapping() || options.GetRenamePrefixNamespace() != nil {
		passes.MaybeAdd(c.rewriteGlobalDeclarationsForTryCatchWrapping)
	}

	passes.MaybeAdd(emptyBeforeExtractPrototypeMemberDeclarations)
	if options.GetExtractPrototypeMemberDeclarationsMode() != CompilerOptions_ExtractPrototypeMemberDeclarationsMode_OFF {
		passes.MaybeAdd(c.extractPrototypeMemberDeclarations)
	}

	if options.ShouldAmbiguateProperties() &&
		options.GetPropertyRenaming() == PropertyRenamingPolicy_ALL_UNQUOTED &&
		options.IsTypecheckingEnabled() {
		passes.MaybeAdd(c.ambiguateProperties)
	}

	passes.MaybeAdd(emptyBeforeRenameProperties)
	if options.GetPropertyRenaming() == PropertyRenamingPolicy_ALL_UNQUOTED {
		passes.MaybeAdd(c.renameProperties)
	} else {
		passes.MaybeAdd(c.removePropertyRenamingCalls)
	}

	if options.ShouldReserveRawExports() {
		passes.MaybeAdd(c.gatherRawExports)
	}
	if options.ShouldConvertToDottedProperties() {
		passes.MaybeAdd(c.convertToDottedProperties)
	}
	if options.ShouldRewriteFunctionExpressions() {
		passes.MaybeAdd(c.rewriteFunctionExpressions)
	}
	if options.GetAliasStringsMode() != CompilerOptions_AliasStringsMode_NONE {
		passes.MaybeAdd(c.aliasStrings)
	}

	if options.ShouldCoalesceVariableNames() {
		passes.MaybeAdd(c.coalesceVariableNames)
		if options.ShouldFoldConstants() {
			passes.MaybeAdd(c.peepholeOptimizationsOnceNonNormalized)
		}
	}

	passes.MaybeAdd(c.markUnnormalized)

	if options.ShouldCollapseVariableDeclarations() || options.ShouldOptimizeLetAndConst() {
		passes.MaybeAdd(c.postNormalizePeephole)
	}
	if options.ShouldCollapseVariableDeclarations() {
		passes.MaybeAdd(c.collapseVariableDeclarations)
	}

	passes.MaybeAdd(c.denormalize)
	passes.MaybeAdd(emptyBeforeVariableRenaming)

	if options.GetVariableRenaming() != VariableRenamingPolicy_ALL {
		passes.MaybeAdd(c.invertContextualRenaming)
	}
	if options.GetVariableRenaming() != VariableRenamingPolicy_OFF {
		passes.MaybeAdd(c.renameVars)
	}
	if options.ShouldRenameLabels() {
		passes.MaybeAdd(c.renameLabels)
	}
	if options.ShouldFoldConstants() {
		passes.MaybeAdd(c.latePeepholeOptimizations)
	}

	if options.ShouldProtectHiddenSideEffects() || options.GetMergedPrecompiledLibraries() {
		passes.MaybeAdd(c.stripSideEffectProtection)
	}

	if ns := options.GetRenamePrefixNamespace(); ns != nil &&
		options.GetChunkOutputType() == CompilerOptions_ChunkOutputType_GLOBAL_NAMESPACE {
		validateRenamePrefixNamespace(*ns)
		passes.MaybeAdd(c.rescopeGlobalSymbols)
	}

	if options.GetOutputFeatureSet().Contains(parser.FeatureSet_ES2015) {
		passes.MaybeAdd(c.optimizeToEs6)
	}
	if options.GetChunkOutputType() == CompilerOptions_ChunkOutputType_ES_MODULES {
		passes.MaybeAdd(c.convertChunksToESModules)
	}

	passes.MaybeAdd(c.checkAstValidity)
	passes.MaybeAdd(c.varCheckValidity)
	return passes
}

func (c *DefaultPassConfig) getEarlyOptimizationLoopPasses() *PassListBuilder {
	options := c.options
	earlyLoopPasses := NewPassListBuilder(options)

	if options.ShouldInlineVariables() || options.ShouldInlineLocalVariables() {
		earlyLoopPasses.MaybeAdd(c.inlineVariables)
	} else if options.ShouldInlineConstantVars() {
		earlyLoopPasses.MaybeAdd(c.inlineConstants)
	}
	if options.GetCollapseObjectLiterals() {
		earlyLoopPasses.MaybeAdd(c.collapseObjectLiterals)
	}
	if c.shouldRunRemoveUnusedCode() {
		earlyLoopPasses.MaybeAdd(c.removeUnusedCode)
	}
	if options.ShouldFoldConstants() {
		earlyLoopPasses.MaybeAdd(c.peepholeOptimizations)
	}

	earlyLoopPasses.AssertAllLoopablePasses()
	return earlyLoopPasses
}

func (c *DefaultPassConfig) getPostL10nOptimizations() *PassListBuilder {
	options := c.options
	loopPasses := NewPassListBuilder(options)

	if options.ShouldOptimizeCalls() {
		loopPasses.MaybeAdd(c.optimizeCalls)
	}
	if options.GetJ2clPass().ShouldAddJ2clPasses() {
		loopPasses.MaybeAdd(c.j2clConstantHoisterPass)
		loopPasses.MaybeAdd(c.j2clClinitPass)
	}
	if options.GetInlineFunctionsLevel() != CompilerOptions_Reach_NONE {
		loopPasses.MaybeAdd(c.inlineFunctions)
	}
	if options.ShouldInlineVariables() || options.ShouldInlineLocalVariables() {
		loopPasses.MaybeAdd(c.inlineVariables)
	} else if options.ShouldInlineConstantVars() {
		loopPasses.MaybeAdd(c.inlineConstants)
	}
	if c.shouldRunRemoveUnusedCode() {
		loopPasses.MaybeAdd(c.removeUnusedCode)
	}
	if options.ShouldFoldConstants() {
		loopPasses.MaybeAdd(c.peepholeOptimizations)
	}

	loopPasses.AssertAllLoopablePasses()
	return loopPasses
}

func (c *DefaultPassConfig) getEarlyOptimizationPasses() *PassListBuilder {
	options := c.options
	passes := NewPassListBuilder(options)

	if options.ShouldExportTestFunctions() {
		passes.MaybeAdd(c.exportTestFunctions)
	}

	if options.GetMergedPrecompiledLibraries() {
		passes.MaybeAdd(c.removeWeakSources)
		passes.MaybeAdd(c.checkRegExpForOptimizations)
		if options.GetJ2clPass().ShouldAddJ2clPasses() {
			passes.MaybeAdd(c.j2clSourceFileChecker)
		}
	} else {
		c.addNonTypedAstNormalizationPasses(passes)
	}

	passes.MaybeAdd(c.removeUnnecessarySyntheticExterns)

	if options.GetSyntheticBlockStartMarker() != nil {
		passes.MaybeAdd(c.createSyntheticBlocks)
	}
	if options.GetJ2clPass().ShouldAddJ2clPasses() {
		passes.MaybeAdd(c.j2clPass)
	}

	passes.MaybeAdd(c.transpileAndOptimizeClosureUnaware)

	TranspilationPasses_AddTranspilationRuntimeLibraries(passes)

	if options.GetRewritePolyfills() || options.GetIsolatePolyfills() || options.GetInjectPolyfillsNewerThan() != nil {
		TranspilationPasses_AddRewritePolyfillPass(passes)
	}

	passes.MaybeAdd(c.injectRuntimeLibraries)

	if options.GetClosurePass() {
		passes.MaybeAdd(c.closureProvidesRequires)
	}

	if options.ShouldRunReplaceMessagesForChrome() {
		passes.MaybeAdd(c.replaceMessagesForChrome)
	} else if options.ShouldRunReplaceMessagesPass() {
		if options.DoLateLocalization() {
			passes.MaybeAdd(c.getProtectMessagesPass())
		} else {
			passes.MaybeAdd(c.getFullReplaceMessagesPass())
		}
	}

	if options.DoLateLocalization() {
		passes.MaybeAdd(c.protectLocaleData)
	}

	if options.GetClosurePass() && !options.ShouldPreserveGoogLibraryPrimitives() {
		passes.MaybeAdd(c.closureReplaceGetCssName)
	}

	passes.MaybeAdd(c.processDefinesOptimize)
	passes.MaybeAdd(emptyBeforeEarlyOptimizationsTranspilation)
	passes.MaybeAdd(c.normalize)
	passes.MaybeAdd(c.gatherGettersAndSetters)

	TranspilationPasses_AddTranspilationPasses(passes, options)

	if options.GetJ2clPass().ShouldAddJ2clPasses() {
		passes.MaybeAdd(c.j2clUtilGetDefineRewriterPass)
	}
	if options.GetInstrumentForCoverageOption() != CompilerOptions_InstrumentOption_NONE {
		passes.MaybeAdd(c.instrumentForCodeCoverage)
	}

	passes.MaybeAdd(c.gatherExternPropertiesOptimize)
	passes.MaybeAdd(emptyBeforeStandardOptimizations)

	if options.GetClosurePass() && (options.ShouldRemoveAbstractMethods() || options.ShouldRemoveClosureAsserts()) {
		passes.MaybeAdd(c.closureCodeRemoval)
	}
	if options.ShouldRemoveJ2clAsserts() {
		passes.MaybeAdd(c.j2clAssertRemovalPass)
	}

	passes.MaybeAdd(c.replaceToggles)
	passes.MaybeAdd(c.inlineAndCollapseProperties)

	if options.GetTweakProcessing().ShouldStrip() ||
		!options.GetStripTypes().IsEmpty() ||
		!options.GetStripNameSuffixes().IsEmpty() ||
		!options.GetStripNamePrefixes().IsEmpty() {
		passes.MaybeAdd(c.stripCode)
	}

	if options.ShouldReplaceIdGenerators() {
		passes.MaybeAdd(c.replaceIdGenerators)
	}

	if options.GetJ2clPass().ShouldAddJ2clPasses() &&
		options.GetPropertyCollapseLevel() == CompilerOptions_PropertyCollapseLevel_ALL {
		passes.MaybeAdd(c.j2clPropertyInlinerPass)
	}

	if options.ShouldInferConsts() {
		passes.MaybeAdd(c.inferConsts)
	}

	passes.MaybeAdd(emptyObfuscationPassMarker)

	if options.GetSmartNameRemoval() {
		if options.ShouldFoldConstants() &&
			(options.ShouldInlineVariables() || options.ShouldInlineLocalVariables()) {
			passes.MaybeAdd(c.earlyInlineVariables)
			passes.MaybeAdd(c.earlyPeepholeOptimizations)
		}
		passes.MaybeAdd(c.removeUnusedCodeOnce)
	}

	if options.ShouldDisambiguateProperties() && options.IsTypecheckingEnabled() {
		passes.MaybeAdd(c.disambiguateProperties)
	}
	if options.ShouldComputeFunctionSideEffects() {
		passes.MaybeAdd(c.markPureFunctions)
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
		passes.MaybeAdd(c.closureOptimizePrimitives)
	}
	if len(options.GetReplaceStringsFunctionDescriptions()) != 0 {
		passes.MaybeAdd(c.replaceStrings)
	}
	if options.ShouldRunCrossChunkCodeMotion() {
		passes.MaybeAdd(c.crossChunkCodeMotion)
	}
	if options.ShouldDevirtualizeMethods() {
		passes.MaybeAdd(c.devirtualizeMethods)
	}

	passes.MaybeAdd(c.getCustomPasses(CustomPassExecutionTime_BEFORE_OPTIMIZATION_LOOP))
	passes.MaybeAdd(emptyBeforeMainOptimizations)

	if options.ShouldInlineVariables() || options.ShouldInlineLocalVariables() {
		passes.MaybeAdd(c.flowSensitiveInlineVariables)
	}

	passes.AddAll(c.getMainOptimizationLoop())
	passes.MaybeAdd(emptyAfterMainOptimizations)
	passes.MaybeAdd(c.getCustomPasses(CustomPassExecutionTime_AFTER_OPTIMIZATION_LOOP))

	c.assertValidOrderForOptimizations(passes)
	return passes
}

func (c *DefaultPassConfig) getMainOptimizationLoop() *PassListBuilder {
	options := c.options
	passes := NewPassListBuilder(options)

	if options.ShouldInlineGetters() {
		passes.MaybeAdd(c.inlineSimpleMethods)
	}
	if options.ShouldInlineProperties() && options.IsTypecheckingEnabled() {
		passes.MaybeAdd(c.inlineProperties)
	}
	if options.ShouldRunDeadPropertyAssignmentElimination() {
		passes.MaybeAdd(c.deadPropertyAssignmentElimination)
	}
	if options.ShouldOptimizeCalls() {
		passes.MaybeAdd(c.optimizeCalls)
	}
	if options.GetJ2clPass().ShouldAddJ2clPasses() {
		passes.MaybeAdd(c.j2clConstantHoisterPass)
		passes.MaybeAdd(c.j2clClinitPass)
	}
	// inlineVariables and peepholeOptimizations must run after inlineFunctions to clean up patterns
	// it introduces (which affects the size-based loop-termination heuristic).
	if options.GetInlineFunctionsLevel() != CompilerOptions_Reach_NONE {
		passes.MaybeAdd(c.inlineFunctions)
	}
	if options.ShouldInlineVariables() || options.ShouldInlineLocalVariables() {
		passes.MaybeAdd(c.inlineVariables)
	} else if options.ShouldInlineConstantVars() {
		passes.MaybeAdd(c.inlineConstants)
	}
	if options.ShouldRunDeadAssignmentElimination() {
		passes.MaybeAdd(c.deadAssignmentsElimination)
	}
	if options.GetCollapseObjectLiterals() {
		passes.MaybeAdd(c.collapseObjectLiterals)
	}
	if c.shouldRunRemoveUnusedCode() {
		passes.MaybeAdd(c.removeUnusedCode)
	}
	if options.ShouldFoldConstants() {
		passes.MaybeAdd(c.peepholeOptimizations)
	}

	passes.AssertAllLoopablePasses()
	return passes
}

func (c *DefaultPassConfig) addNonTypedAstNormalizationPasses(passes *PassListBuilder) {
	passes.MaybeAdd(c.removeCastNodes)
	passes.MaybeAdd(c.typesToColors)
}

func (c *DefaultPassConfig) shouldRunRemoveUnusedCode() bool {
	options := c.options
	return options.ShouldRemoveUnusedVariables() ||
		options.ShouldRemoveUnusedLocalVariables() ||
		options.ShouldRemoveUnusedPrototypeProperties() ||
		options.IsRemoveUnusedClassProperties() ||
		options.GetRewritePolyfills()
}

func (c *DefaultPassConfig) assertValidOrderForChecks(checks *PassListBuilder) {
	checks.AssertPassOrder(
		c.declaredGlobalExternsOnWindow,
		c.checkVars,
		"declaredGlobalExternsOnWindow must happen before VarCheck, which adds synthetic externs")
	checks.AssertPassOrder(
		c.chromePass,
		c.checkJsDocAndEs6Modules,
		"The ChromePass must run before after JsDoc and Es6 module checking.")
	checks.AssertPassOrder(
		c.closureRewriteModule,
		c.processDefinesCheck,
		"Must rewrite goog.module before processing @define's, so that @defines in modules work.")
	checks.AssertPassOrder(
		c.closurePrimitives, c.polymerPass, "The Polymer pass must run after goog.provide processing.")
	checks.AssertPassOrder(
		c.chromePass, c.polymerPass, "The Polymer pass must run after ChromePass processing.")
	checks.AssertPassOrder(
		c.polymerPass, c.suspiciousCode, "The Polymer pass must run before suspiciousCode processing.")
	checks.AssertPassOrder(
		c.addSyntheticScript,
		c.gatherModuleMetadataPass,
		"Cannot add a synthetic script node after module metadata creation.")
	checks.AssertPassOrder(
		c.closureRewriteModule,
		c.removeSyntheticScript,
		"Synthetic script node should be removed only after module rewriting.")
	checks.AssertPassOrder(
		c.closureRewriteModule,
		c.rewriteCallerCodeLocation,
		"ClosureRewriteModule must happen before RewriteCallerCodeLocation, so that exported"+
			" functions and call sites are rewritten correctly.")
	checks.AssertPassOrder(
		c.closureRewriteModule,
		TranspilationPasses_GetEs6RewriteDestructuring(
			Es6RewriteDestructuring_ObjectDestructuringRewriteMode_REWRITE_ALL_OBJECT_PATTERNS),
		"RewriteCallerCodeLocation must happen before Es6RewriteDestructuring, because we need"+
			" ReWriteCallerCodeLocation to run before default parameters get rewritten.")
	checks.AssertPassOrder(
		c.closureRewriteModule,
		TranspilationPasses_GetEs6RewriteDestructuring(
			Es6RewriteDestructuring_ObjectDestructuringRewriteMode_REWRITE_OBJECT_REST),
		"RewriteCallerCodeLocation must happen before Es6RewriteDestructuring, because we need"+
			" ReWriteCallerCodeLocation to run before default parameters get rewritten.")

	if checks.Contains(c.closureGoogScopeAliases) {
		jutil.CheckStateMsg(
			checks.Contains(c.checkVariableReferences),
			"goog.scope processing requires variable checking")
	}

	checks.AssertPassOrder(
		c.gatherModuleMetadataPass,
		c.closureCheckModule,
		"Need to gather module metadata before checking closure modules.")
	checks.AssertPassOrder(
		c.gatherModuleMetadataPass,
		c.createModuleMapPass,
		"Need to gather module metadata before scanning modules.")
	checks.AssertPassOrder(
		c.createModuleMapPass,
		c.rewriteCommonJsModules,
		"Need to gather module information before rewriting CommonJS modules.")
	checks.AssertPassOrder(
		c.rewriteScriptsToEs6Modules,
		c.gatherModuleMetadataPass,
		"Need to gather module information after rewriting scripts to modules.")
	checks.AssertPassOrder(
		c.gatherModuleMetadataPass,
		c.checkMissingRequires,
		"Need to gather module information before checking for missing requires.")
}

func (c *DefaultPassConfig) assertValidOrderForOptimizations(optimizations *PassListBuilder) {
	optimizations.AssertPassOrder(
		c.j2clPass,
		TranspilationPasses_rewriteGenerators,
		"J2CL normalization should be done before generator re-writing.")
	optimizations.AssertPassOrder(
		TranspilationPasses_rewritePolyfills,
		c.processDefinesOptimize,
		"Polyfill injection must be done before processDefines as some polyfills reference goog.defines.")
	optimizations.AssertPassOrder(
		TranspilationPasses_injectTranspilationRuntimeLibraries,
		c.processDefinesOptimize,
		"Runtime library injection must be done before processDefines some runtime libraries reference goog.defines.")
	optimizations.AssertPassOrder(
		c.processDefinesOptimize,
		c.j2clUtilGetDefineRewriterPass,
		"J2CL define re-writing should be done after processDefines since it relies on Compiler#getDefineNames to have been populated by it.")
	optimizations.AssertPassOrder(
		c.removeUnusedCode,
		c.isolatePolyfills,
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
