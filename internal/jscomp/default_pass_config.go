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

// DefaultPassConfig holds the pass factories and meta-data for native
// JSCompiler passes.
//
// NOTE(johnlenz): this needs some non-trivial refactoring. The pass config
// should use as little state as possible. The recommended way for a pass to
// leave behind some state for a subsequent pass is through the compiler
// object. Any other state remaining here should only be used when the pass
// config is creating the list of checks and optimizations, not after passes
// have started executing.
//
// The general goal is for this class to be as minimal as possible. Option
// validation should occur before this class configures the compiler,
// business logic should live here (passes should not be inlined, etc).
//
// The PassFactory fields are initialized by NewDefaultPassConfig, in Java
// declaration order (initFields).
type DefaultPassConfig struct {
	PassConfigBase

	// A symbol table for registering references that get removed during
	// preprocessing.
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

// For the --mark-as-compiled pass
//
// Java: DefaultPassConfig.COMPILED_CONSTANT_NAME
const defaultPassConfig_COMPILED_CONSTANT_NAME = "COMPILED"

// Constant name for Closure's locale
//
// Java: DefaultPassConfig.CLOSURE_LOCALE_CONSTANT_NAME
const defaultPassConfig_CLOSURE_LOCALE_CONSTANT_NAME = "goog.LOCALE"

// Miscellaneous errors.
//
// Java: DefaultPassConfig.GLOBAL_SYMBOL_NAMESPACE_PATTERN
var defaultPassConfig_GLOBAL_SYMBOL_NAMESPACE_PATTERN = jutil.Pattern_Compile("^[a-zA-Z0-9$_]+$")

// Java: DefaultPassConfig(CompilerOptions)
func NewDefaultPassConfig(options *CompilerOptions) *DefaultPassConfig {
	c := &DefaultPassConfig{}
	c.self = c
	c.InitPassConfigBase(options)
	c.initFields()
	return c
}

// Java: DefaultPassConfig.getPreprocessorSymbolTable()
func (c *DefaultPassConfig) GetPreprocessorSymbolTable() *PreprocessorSymbolTable {
	return c.preprocessorSymbolTableFactory.GetInstanceOrNull()
}

// Java: DefaultPassConfig.getTranspileOnlyPasses()
func (c *DefaultPassConfig) GetTranspileOnlyPasses() *PassListBuilder {
	options := c.options
	passes := NewPassListBuilder(options)

	// Certain errors in block-scoped variable declarations will prevent correct transpilation
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
			// nothing
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

	// Passes below this point may rely on normalization and must maintain normalization.
	passes.MaybeAdd(c.normalize)

	passes.MaybeAdd(c.gatherGettersAndSetters)

	TranspilationPasses_AddTranspilationPasses(passes, options)

	passes.MaybeAdd(c.markUnnormalized)
	// The transpilation passes may rely on normalize making all variables unique,
	// but we're doing only transpilation, so we want to put back the original variable names
	// wherever we can to meet user expectations.
	//
	// Also, if we don't do this, we could end up breaking runtime behavior that depends on specific
	// variable names.
	//
	// The primary concern is function parameter names, because some frameworks, like Angular,
	// do runtime injection of function call arguments based on the function parameter names.
	passes.MaybeAdd(c.invertContextualRenaming)

	// Es6ConvertSuper may add this so it needs to be removed for transpile-only mode.
	passes.MaybeAdd(c.removePropertyRenamingCalls)

	passes.AssertAllOneTimePasses()
	c.assertValidOrderForChecks(passes)
	return passes
}

// Java: DefaultPassConfig.getWhitespaceOnlyPasses()
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

// Java: DefaultPassConfig.addModuleRewritingPasses(PassListBuilder, CompilerOptions)
func (c *DefaultPassConfig) addModuleRewritingPasses(checks *PassListBuilder, options *CompilerOptions) {
	if options.GetLanguageIn().ToFeatureSet().Has(parser.FeatureSet_Feature_MODULES) {
		checks.MaybeAdd(c.rewriteGoogJsImports)
		TranspilationPasses_AddEs6ModulePass(checks, c.preprocessorSymbolTableFactory)
	}

	if options.GetClosurePass() {
		checks.MaybeAdd(c.closureRewriteModule)
	}
}

// Java: DefaultPassConfig.getChecks()
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
		checks.MaybeAdd(PassFactory_CreateEmptyPass(PassNames_AFTER_STANDARD_CHECKS))
		checks.MaybeAdd(c.mergeSyntheticScript)
		checks.MaybeAdd(c.gatherExternPropertiesCheck)
		checks.MaybeAdd(PassFactory_CreateEmptyPass(PassNames_BEFORE_SERIALIZATION))
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

	// Run this pass before any pass tries to inject new runtime libraries
	checks.MaybeAdd(c.addSyntheticScript)

	checks.MaybeAdd(c.gatherGettersAndSetters)

	if options.GetLanguageIn().ToFeatureSet().ContainsFeature(parser.FeatureSet_Feature_DYNAMIC_IMPORT) &&
		!options.ShouldAllowDynamicImport() {
		checks.MaybeAdd(defaultPassConfig_forbidDynamicImportUsage)
	}

	checks.MaybeAdd(PassFactory_CreateEmptyPass("beforeStandardChecks"))

	if !options.GetProcessCommonJSModules() &&
		options.GetLanguageIn().ToFeatureSet().Has(parser.FeatureSet_Feature_MODULES) {
		checks.MaybeAdd(c.rewriteScriptsToEs6Modules)
	}

	// Run these passes after promoting scripts to modules, but before rewriting any other modules.
	checks.MaybeAdd(c.gatherModuleMetadataPass)
	checks.MaybeAdd(c.createModuleMapPass)

	if options.GetProcessCommonJSModules() {
		checks.MaybeAdd(c.rewriteCommonJsModules)
	}

	// Note: ChromePass can rewrite invalid @type annotations into valid ones, so should run before
	// JsDoc checks.
	if options.IsChromePassEnabled() {
		checks.MaybeAdd(c.chromePass)
	}

	// Verify JsDoc annotations and check ES modules
	checks.MaybeAdd(c.checkJsDocAndEs6Modules)

	checks.MaybeAdd(c.checkTypeImportCodeReferences)

	if options.Enables(DiagnosticGroups_LINT_CHECKS) {
		checks.MaybeAdd(c.lintChecks)
	}

	if options.GetClosurePass() && options.Enables(DiagnosticGroups_LINT_CHECKS) {
		checks.MaybeAdd(c.checkRequiresAndProvidesSorted)
	}

	if options.Enables(DiagnosticGroups_EXTRA_REQUIRE) {
		checks.MaybeAdd(c.extraRequires)
	}

	if options.Enables(DiagnosticGroups_MISSING_REQUIRE) {
		checks.MaybeAdd(c.checkMissingRequires)
	}

	checks.MaybeAdd(c.declaredGlobalExternsOnWindow)

	if !options.GetProcessCommonJSModules() {
		// TODO(ChadKillingsworth): move CommonJS module rewriting after VarCheck
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
	}

	if options.GetClosurePass() {
		checks.MaybeAdd(c.closurePrimitives)
	}

	// TODO(b/141389184): Move this after the Polymer pass
	if options.ShouldRewriteModulesBeforeTypechecking() {
		c.addModuleRewritingPasses(checks, options)
	}

	// It's important that the PolymerPass run *after* the ClosurePrimitives and ChromePass rewrites
	// and *before* the suspicious code checks. This is enforced in the assertValidOrder method.
	if options.IsPolymerPassEnabled() {
		checks.MaybeAdd(c.polymerPass)
	}

	if options.GetProcessCommonJSModules() {
		// TODO(ChadKillingsworth): remove this branch.
		checks.MaybeAdd(c.checkVariableReferences)
		checks.MaybeAdd(c.checkVars)
	}

	if options.ShouldInferConsts() {
		checks.MaybeAdd(c.inferConsts)
	}

	if options.ShouldComputeFunctionSideEffects() {
		checks.MaybeAdd(c.checkRegExp)
	}

	checks.MaybeAdd(PassFactory_CreateEmptyPass(PassNames_BEFORE_TYPE_CHECKING))

	if options.GetCheckTypes() || options.GetInferTypes() {
		checks.MaybeAdd(c.inferTypes)
		if options.GetCheckTypes() {
			checks.MaybeAdd(c.checkTypes)
		} else {
			checks.MaybeAdd(c.inferJsDocInfo)
		}
	}

	// Analyzer checks must be run after typechecking but before module rewriting.
	if options.Enables(DiagnosticGroups_ANALYZER_CHECKS) && options.IsTypecheckingEnabled() {
		checks.MaybeAdd(c.analyzerChecks)
	}

	// We assume that only clients who are going to re-compile, or do in-depth static analysis,
	// will need the typed scope creator after the compile job.
	if !options.PreservesDetailedSourceInfo() {
		checks.MaybeAdd(c.clearTypedScopeCreatorPass)
	}

	if options.ShouldRewriteModulesAfterTypechecking() {
		c.addModuleRewritingPasses(checks, options)
	}

	// Dynamic import rewriting must always occur after type checking
	if options.ShouldAllowDynamicImport() &&
		options.GetLanguageIn().ToFeatureSet().Has(parser.FeatureSet_Feature_DYNAMIC_IMPORT) {
		checks.MaybeAdd(defaultPassConfig_rewriteDynamicImports)
	}

	// We assume that only clients who are going to re-compile, or do in-depth static analysis,
	// will need the typed scope creator after the compile job.
	if !options.PreservesDetailedSourceInfo() {
		checks.MaybeAdd(c.clearTopTypedScopePass)
	}

	// CheckSuspiciousCode requires type information, so must run after the type checker.
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

	// CheckAccessControls only works if check types is on.
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

	checks.MaybeAdd(PassFactory_CreateEmptyPass(PassNames_AFTER_STANDARD_CHECKS))

	checks.MaybeAdd(c.mergeSyntheticScript)

	// Create extern exports after the normalize because externExports depends on unique names.
	if options.GetExternExportsPath() != nil {
		checks.MaybeAdd(c.externExports)
	}

	if !options.IsChecksOnly() {
		checks.MaybeAdd(c.removeWeakSources)
	}

	// Gather property names in externs so they can be queried by the optimizing passes.
	// See b/180424427 for why this runs in stage 1 and not stage 2.
	checks.MaybeAdd(c.gatherExternPropertiesCheck)

	checks.MaybeAdd(PassFactory_CreateEmptyPass(PassNames_BEFORE_SERIALIZATION))

	if options.GetTypedAstOutputFile() != nil {
		checks.MaybeAdd(c.serializeTypedAst)
	}

	// Validation should be after all pass additions
	checks.AssertAllOneTimePasses()
	c.assertValidOrderForChecks(checks)

	return checks
}

// Java: DefaultPassConfig.getOptimizations()
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
		passes.MaybeAdd(PassFactory_CreateEmptyPass(PassNames_BEFORE_STANDARD_OPTIMIZATIONS))
		passes.MaybeAdd(c.inlineAndCollapseProperties)
		passes.MaybeAdd(c.closureOptimizePrimitives)
		// If side-effects were protected, remove the protection now.
		passes.MaybeAdd(c.stripSideEffectProtection)
		return passes
	}

	if options.GetSkipNonTranspilationPasses() {
		// Reaching this if-condition means the 'getChecks()' phase has been skipped in favor of
		// 'getTranspileOnlyPasses'.
		return passes
	}

	passes.AddAll(c.getEarlyOptimizationPasses())
	passes.MaybeAdd(PassFactory_CreateEmptyPass(PassNames_OPTIMIZATIONS_HALFWAY_POINT))
	passes.AddAll(c.getLateOptimizationPasses())

	return passes
}

// Java: DefaultPassConfig.getFinalizations()
func (c *DefaultPassConfig) GetFinalizations() *PassListBuilder {
	options := c.options
	passes := NewPassListBuilder(options)

	if options.IsPropertyRenamingOnlyCompilationMode() {
		if options.ShouldRewriteGlobalDeclarationsForTryCatchWrapping() ||
			options.GetRenamePrefixNamespace() != nil {
			passes.MaybeAdd(c.rewriteGlobalDeclarationsForTryCatchWrapping)
		}

		passes.MaybeAdd(c.renameProperties)

		if options.GetRenamePrefixNamespace() != nil &&
			options.GetChunkOutputType() == CompilerOptions_ChunkOutputType_GLOBAL_NAMESPACE {
			if !defaultPassConfig_GLOBAL_SYMBOL_NAMESPACE_PATTERN.Matches(*options.GetRenamePrefixNamespace()) {
				panic(jutil.NewIllegalArgumentExceptionMsg(
					"Illegal character in renamePrefixNamespace name: " +
						*options.GetRenamePrefixNamespace()))
			}
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

		// After inlining variable uses, some variables may be unused.
		// If we're doing late localization, the simple code removal pass runs added below will clean
		// those up. Otherwise, clean them up now.
		if !options.DoLateLocalization() && c.shouldRunRemoveUnusedCode() {
			passes.MaybeAdd(c.removeUnusedCodeOnce)
		}
	}

	if options.DoLateLocalization() {
		passes.AddAll(c.getPostL10nOptimizations())
	}

	passes.MaybeAdd(PassFactory_CreateEmptyPass("beforeChunkMotion"))

	if options.ShouldRunCrossChunkCodeMotion() {
		passes.MaybeAdd(c.crossChunkCodeMotion)
	}

	if options.ShouldRunCrossChunkMethodMotion() {
		passes.MaybeAdd(c.crossChunkMethodMotion)
	}

	passes.MaybeAdd(PassFactory_CreateEmptyPass("afterChunkMotion"))

	if options.GetOptimizeESClassConstructors() &&
		options.GetOutputFeatureSet().Contains(parser.FeatureSet_ES2015) {
		passes.MaybeAdd(c.optimizeConstructors)
	}

	// Isolate injected polyfills from the global scope. Runs late in the optimization loop
	// to take advantage of property renaming & RemoveUnusedCode, as this pass will increase code
	// size by wrapping all potential polyfill usages.
	if options.GetIsolatePolyfills() {
		passes.MaybeAdd(c.isolatePolyfills)
	}

	if options.ShouldCollapseAnonymousFunctions() {
		// TODO: b/197349249 - Maybe we should just move this pass after denormalization. It seems
		// weird to convert from function expression to function declaration while we're still
		// supposed to be in a normalized state. But it requires testing as perhaps some optimizations
		// in that range will get affected if we skip this rewriting.
		passes.MaybeAdd(c.collapseAnonymousFunctions)
	}

	// Move functions before extracting prototype member declarations.
	if options.ShouldRewriteGlobalDeclarationsForTryCatchWrapping() ||
		// renamePrefixNamescape relies on rewriteGlobalDeclarationsForTryCatchWrapping
		// to preserve semantics.
		options.GetRenamePrefixNamespace() != nil {
		passes.MaybeAdd(c.rewriteGlobalDeclarationsForTryCatchWrapping)
	}

	passes.MaybeAdd(PassFactory_CreateEmptyPass(PassNames_BEFORE_EXTRACT_PROTOTYPE_MEMBER_DECLARATIONS))
	// The mapped name anonymous function pass makes use of information that
	// the extract prototype member declarations pass removes so the former
	// happens before the latter.
	if options.GetExtractPrototypeMemberDeclarationsMode() !=
		CompilerOptions_ExtractPrototypeMemberDeclarationsMode_OFF {
		passes.MaybeAdd(c.extractPrototypeMemberDeclarations)
	}

	if options.ShouldAmbiguateProperties() &&
		options.GetPropertyRenaming() == PropertyRenamingPolicy_ALL_UNQUOTED &&
		options.IsTypecheckingEnabled() {
		passes.MaybeAdd(c.ambiguateProperties)
	}

	passes.MaybeAdd(PassFactory_CreateEmptyPass(PassNames_BEFORE_RENAME_PROPERTIES))
	if options.GetPropertyRenaming() == PropertyRenamingPolicy_ALL_UNQUOTED {
		passes.MaybeAdd(c.renameProperties)
	} else {
		passes.MaybeAdd(c.removePropertyRenamingCalls)
	}

	// Reserve global names added to the "windows" object.
	if options.ShouldReserveRawExports() {
		passes.MaybeAdd(c.gatherRawExports)
	}

	// This comes after property renaming because quoted property names must
	// not be renamed.
	if options.ShouldConvertToDottedProperties() {
		passes.MaybeAdd(c.convertToDottedProperties)
	}

	// Property renaming must happen before this pass runs since this
	// pass may convert dotted properties into quoted properties.  It
	// is beneficial to run before alias strings, alias keywords and
	// variable renaming.
	if options.ShouldRewriteFunctionExpressions() {
		passes.MaybeAdd(c.rewriteFunctionExpressions)
	}

	// This comes after converting quoted property accesses to dotted property
	// accesses in order to avoid aliasing property names.
	if options.GetAliasStringsMode() != CompilerOptions_AliasStringsMode_NONE {
		passes.MaybeAdd(c.aliasStrings)
	}

	if options.ShouldCoalesceVariableNames() {
		// Passes after this point can no longer depend on normalized AST
		// assumptions because the code is marked as un-normalized
		passes.MaybeAdd(c.coalesceVariableNames)

		// coalesceVariables creates identity assignments and more redundant code
		// that can be removed, rerun the peephole optimizations to clean them
		// up.
		if options.ShouldFoldConstants() {
			passes.MaybeAdd(c.peepholeOptimizationsOnceNonNormalized)
		}
	}

	// Passes after this point can no longer depend on normalized AST assumptions.
	passes.MaybeAdd(c.markUnnormalized)

	if options.ShouldCollapseVariableDeclarations() || options.ShouldOptimizeLetAndConst() {
		passes.MaybeAdd(c.postNormalizePeephole)
	}
	if options.ShouldCollapseVariableDeclarations() {
		passes.MaybeAdd(c.collapseVariableDeclarations)
	}

	// This pass works best after collapseVariableDeclarations.
	passes.MaybeAdd(c.denormalize)

	passes.MaybeAdd(PassFactory_CreateEmptyPass(PassNames_BEFORE_VARIABLE_RENAMING))

	if options.GetVariableRenaming() != VariableRenamingPolicy_ALL {
		// If we're leaving some (or all) variables with their old names,
		// then we need to undo any of the markers we added for distinguishing
		// local variables ("x" -> "x$jscomp$1").
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

	// If side-effects were protected, remove the protection now.
	// Note that when using precompiled libraries we always run this pass regardless of the
	// 'shouldProtectHiddenSideEffects' option: the library compilation may have run
	// with side effect protection enabled even if the binary disables it, so we assume we may
	// always need to strip side effect protection.
	if options.ShouldProtectHiddenSideEffects() || options.GetMergedPrecompiledLibraries() {
		passes.MaybeAdd(c.stripSideEffectProtection)
	}

	if options.GetRenamePrefixNamespace() != nil &&
		options.GetChunkOutputType() == CompilerOptions_ChunkOutputType_GLOBAL_NAMESPACE {
		if !defaultPassConfig_GLOBAL_SYMBOL_NAMESPACE_PATTERN.Matches(*options.GetRenamePrefixNamespace()) {
			panic(jutil.NewIllegalArgumentExceptionMsg(
				"Illegal character in renamePrefixNamespace name: " +
					*options.GetRenamePrefixNamespace()))
		}
		passes.MaybeAdd(c.rescopeGlobalSymbols)
	}

	// Raise to ES2015, if allowed
	if options.GetOutputFeatureSet().Contains(parser.FeatureSet_ES2015) {
		passes.MaybeAdd(c.optimizeToEs6)
	}

	// Must run after all non-safety-check passes as the optimizations do not support modules.
	if options.GetChunkOutputType() == CompilerOptions_ChunkOutputType_ES_MODULES {
		passes.MaybeAdd(c.convertChunksToESModules)
	}

	// Safety checks.  These should always be the last passes.
	passes.MaybeAdd(c.checkAstValidity)
	passes.MaybeAdd(c.varCheckValidity)
	return passes
}

// Java: DefaultPassConfig.getEarlyOptimizationLoopPasses()
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

// getPostL10nOptimizations adds optimization passes that need to run after
// late localization has been done.
//
// Once we've replaced message references with string constants and
// `goog.LOCALE` with a constant value, we need to re-run some optimizations,
// so they can throw away code for locales that aren't relevant and perform
// constant folding on the message strings we've now inserted.
//
// Java: DefaultPassConfig.getPostL10nOptimizations()
func (c *DefaultPassConfig) getPostL10nOptimizations() *PassListBuilder {
	options := c.options
	// Localization replaced lots of function calls with constants to get statements like these.
	//
	// `goog.LOCALE = 'es-419';`
	// `const MSG_GREETING = 'Hola';`
	// `const MSG_GREETING_WITH_NAME = 'Hola, ' + person.getName();`
	//
	// Before calling this method we should also have run `flowSensitiveInlineVariables`
	// To make optimization opportunities for peepholeOptimizations.
	//
	// For example, `flowSensitiveInlineVariables` will change this
	// ```
	// var x = 'localized version of message';
	// x = x + '&nbsp;';
	// ```
	// to this
	// ```
	// var x = 'localized version of message' + '&nbsp;';
	// ```
	// Which constant folding can then turn into this
	// ```
	// var x = 'localized version of message&nbsp;';
	// ```
	// Now we should have unblocked a lot of potential optimizations,
	// so do an optimization loop to perform those.
	// This loop is similar to the one created by getMainOptimizationLoop().
	// These should be in the same order as those, but only optimizations we expect to need
	// doing to clean up after localization are included.
	loopPasses := NewPassListBuilder(options)

	if options.ShouldOptimizeCalls() {
		loopPasses.MaybeAdd(c.optimizeCalls)
	}

	if options.GetJ2clPass().ShouldAddJ2clPasses() {
		loopPasses.MaybeAdd(c.j2clConstantHoisterPass)
		loopPasses.MaybeAdd(c.j2clClinitPass)
	}

	// It is important that inlineVariables and peepholeOptimizations run after inlineFunctions,
	// because inlineFunctions relies on them to clean up patterns it introduces. This affects our
	// size-based loop-termination heuristic.
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

// getEarlyOptimizationPasses: these are the passes run in the first half of
// optimizations, which consists of transpilation and some early
// optimization passes.
//
// Java: DefaultPassConfig.getEarlyOptimizationPasses()
func (c *DefaultPassConfig) getEarlyOptimizationPasses() *PassListBuilder {
	options := c.options
	passes := NewPassListBuilder(options)
	// At this point all checks have been done.
	if options.ShouldExportTestFunctions() {
		passes.MaybeAdd(c.exportTestFunctions)
	}

	if options.GetMergedPrecompiledLibraries() {
		// Weak sources aren't removed at the library level
		passes.MaybeAdd(c.removeWeakSources)

		// it would be safe to always recompute side effects even if not using precompiled libraries
		// (the else case) but it's unnecessary so skip it to improve build times.
		passes.MaybeAdd(c.checkRegExpForOptimizations)

		// This runs during getChecks(), so only needs to be run here if using precompiled .typedasts
		if options.GetJ2clPass().ShouldAddJ2clPasses() {
			passes.MaybeAdd(c.j2clSourceFileChecker)
		}
	} else {
		c.addNonTypedAstNormalizationPasses(passes)
	}

	// Remove synthetic extern declarations of names that are now defined in source
	// This is expected to do nothing when in a monolithic build
	passes.MaybeAdd(c.removeUnnecessarySyntheticExterns)

	if options.GetSyntheticBlockStartMarker() != nil {
		// This pass must run before the first fold constants pass.
		passes.MaybeAdd(c.createSyntheticBlocks)
	}

	if options.GetJ2clPass().ShouldAddJ2clPasses() {
		passes.MaybeAdd(c.j2clPass)
	}

	passes.MaybeAdd(c.transpileAndOptimizeClosureUnaware)

	TranspilationPasses_AddTranspilationRuntimeLibraries(passes)

	if options.GetRewritePolyfills() ||
		options.GetIsolatePolyfills() ||
		options.GetInjectPolyfillsNewerThan() != nil {
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
			// With late localization we protect the messages from mangling by optimizations now,
			// then actually replace them after optimizations.
			// The purpose of doing this is to separate localization from optimization so we can
			// optimize just once for all locales.
			passes.MaybeAdd(c.getProtectMessagesPass())
		} else {
			// TODO(bradfordcsmith): At the moment we expect the optimized output may be slightly
			// smaller if you replace messages before optimizing, but if we can change that, it would
			// be good to drop this early replacement entirely.
			passes.MaybeAdd(c.getFullReplaceMessagesPass())
		}
	}

	if options.DoLateLocalization() {
		passes.MaybeAdd(c.protectLocaleData)
	}

	// Replace 'goog.getCssName' before processing defines
	if options.GetClosurePass() && !options.ShouldPreserveGoogLibraryPrimitives() {
		passes.MaybeAdd(c.closureReplaceGetCssName)
	}

	// Defines in code always need to be processed.
	passes.MaybeAdd(c.processDefinesOptimize)
	passes.MaybeAdd(PassFactory_CreateEmptyPass(PassNames_BEFORE_EARLY_OPTIMIZATIONS_TRANSPILATION))

	passes.MaybeAdd(c.normalize)

	passes.MaybeAdd(c.gatherGettersAndSetters)

	// TODO(b/329447979): Add an early removeUnusedCode pass here
	TranspilationPasses_AddTranspilationPasses(passes, options)

	if options.GetJ2clPass().ShouldAddJ2clPasses() {
		passes.MaybeAdd(c.j2clUtilGetDefineRewriterPass)
	}

	if options.GetInstrumentForCoverageOption() != CompilerOptions_InstrumentOption_NONE {
		passes.MaybeAdd(c.instrumentForCodeCoverage)
	}

	passes.MaybeAdd(c.gatherExternPropertiesOptimize)

	passes.MaybeAdd(PassFactory_CreateEmptyPass(PassNames_BEFORE_STANDARD_OPTIMIZATIONS))

	// Abstract method removal works best on minimally modified code, and also
	// only needs to run once.
	if options.GetClosurePass() &&
		(options.ShouldRemoveAbstractMethods() || options.ShouldRemoveClosureAsserts()) {
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

	// Ideally this pass would run before transpilation which would allow it to be simplified.
	// It needs to run after `inlineAndCollapseProperties` in order to identify idGenerator calls.
	if options.ShouldReplaceIdGenerators() {
		passes.MaybeAdd(c.replaceIdGenerators)
	}

	// Inline getters/setters in J2CL classes so that Object.defineProperties() calls (resulting
	// from desugaring) don't block class stripping.
	if options.GetJ2clPass().ShouldAddJ2clPasses() &&
		options.GetPropertyCollapseLevel() == CompilerOptions_PropertyCollapseLevel_ALL {
		// Relies on collapseProperties-triggered aggressive alias inlining.
		passes.MaybeAdd(c.j2clPropertyInlinerPass)
	}

	if options.ShouldInferConsts() {
		passes.MaybeAdd(c.inferConsts)
	}

	// A marker pass to allow {@code ExtraPassConfig} passes to order themselves before
	// RemoveUnusedCode.
	passes.MaybeAdd(PassFactory_CreateEmptyPass(PassNames_OBFUSCATION_PASS_MARKER))

	// Running RemoveUnusedCode before disambiguate properties allows disambiguate properties to be
	// more effective if code that would prevent disambiguation can be removed.
	// TODO(b/66971163): Rename options since we're not actually using smartNameRemoval here now.
	if options.GetSmartNameRemoval() {

		// These passes remove code that is dead because of define flags.
		// If the dead code is weakly typed, running these passes before property
		// disambiguation results in more code removal.
		// The passes are one-time on purpose. (The later runs are loopable.)
		if options.ShouldFoldConstants() &&
			(options.ShouldInlineVariables() || options.ShouldInlineLocalVariables()) {
			passes.MaybeAdd(c.earlyInlineVariables)
			passes.MaybeAdd(c.earlyPeepholeOptimizations)
		}

		passes.MaybeAdd(c.removeUnusedCodeOnce)
	}

	// Property disambiguation should only run once and needs to be done
	// soon after type checking, both so that it can make use of type
	// information and so that other passes can take advantage of the renamed
	// properties.
	if options.ShouldDisambiguateProperties() && options.IsTypecheckingEnabled() {
		passes.MaybeAdd(c.disambiguateProperties)
	}

	if options.ShouldComputeFunctionSideEffects() {
		passes.MaybeAdd(c.markPureFunctions)
	}

	passes.AssertAllOneTimePasses()
	return passes
}

// getLateOptimizationPasses: these are the passes run in the second half of
// optimizations, which consists of the early optimization loop and the main
// optimization loop.
//
// Java: DefaultPassConfig.getLateOptimizationPasses()
func (c *DefaultPassConfig) getLateOptimizationPasses() *PassListBuilder {
	options := c.options
	passes := NewPassListBuilder(options)
	if options.GetSmartNameRemoval() {
		// Place one-time marker passes around this loop to prevent the addition of a looping pass
		// above or below from accidentally becoming part of the loop.
		passes.MaybeAdd(PassFactory_CreateEmptyPass(PassNames_BEFORE_EARLY_OPTIMIZATION_LOOP))
		passes.AddAll(c.getEarlyOptimizationLoopPasses())
		// TODO(): Remove this early loop or rename the option that enables it
		// to something more appropriate.
		passes.MaybeAdd(PassFactory_CreateEmptyPass(PassNames_AFTER_EARLY_OPTIMIZATION_LOOP))
	}

	// This needs to come after the inline constants pass, which is run within
	// the code removing passes.
	if options.GetClosurePass() {
		passes.MaybeAdd(c.closureOptimizePrimitives)
	}

	// ReplaceStrings runs after CollapseProperties in order to simplify
	// pulling in values of constants defined in enums structures. It also runs
	// after disambiguate properties and smart name removal so that it can
	// correctly identify logging types and can replace references to string
	// expressions.
	if len(options.GetReplaceStringsFunctionDescriptions()) != 0 {
		passes.MaybeAdd(c.replaceStrings)
	}

	// TODO(user): This forces a first crack at crossChunkCodeMotion
	// before devirtualization. Once certain functions are devirtualized,
	// it confuses crossChunkCodeMotion ability to recognized that
	// it is recursive.

	// TODO(user): This is meant for a temporary quick win.
	// In the future, we might want to improve our analysis in
	// CrossChunkCodeMotion so we don't need to do this.
	if options.ShouldRunCrossChunkCodeMotion() {
		passes.MaybeAdd(c.crossChunkCodeMotion)
	}

	// Method devirtualization benefits from property disambiguation so
	// it should run after that pass but before passes that do
	// optimizations based on global names (like cross-chunk code motion
	// and inline functions).  RemoveUnusedCode does better if run before
	// this pass.
	if options.ShouldDevirtualizeMethods() {
		passes.MaybeAdd(c.devirtualizeMethods)
	}

	passes.MaybeAdd(c.getCustomPasses(CustomPassExecutionTime_BEFORE_OPTIMIZATION_LOOP))

	passes.MaybeAdd(PassFactory_CreateEmptyPass(PassNames_BEFORE_MAIN_OPTIMIZATIONS))

	// Because FlowSensitiveInlineVariables does not operate on the global scope due to compilation
	// time, we need to run it once before InlineFunctions so that we don't miss inlining
	// opportunities when a function will be inlined into the global scope.
	if options.ShouldInlineVariables() || options.ShouldInlineLocalVariables() {
		passes.MaybeAdd(c.flowSensitiveInlineVariables)
	}

	passes.AddAll(c.getMainOptimizationLoop())
	passes.MaybeAdd(PassFactory_CreateEmptyPass(PassNames_AFTER_MAIN_OPTIMIZATIONS))

	// Some optimizations belong outside the loop because running them more
	// than once would either have no benefit or be incorrect.
	passes.MaybeAdd(c.getCustomPasses(CustomPassExecutionTime_AFTER_OPTIMIZATION_LOOP))

	c.assertValidOrderForOptimizations(passes)
	return passes
}

// getMainOptimizationLoop creates the passes for the main optimization
// loop.
//
// Java: DefaultPassConfig.getMainOptimizationLoop()
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

	// It is important that inlineVariables and peepholeOptimizations run after inlineFunctions,
	// because inlineFunctions relies on them to clean up patterns it introduces. This affects our
	// size-based loop-termination heuristic.
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

// addNonTypedAstNormalizationPasses is for use in builds that run
// getOptimizations() on the result of an AST directly from getChecks(),
// that has not gone through TypedAST serialization/deserialization.
//
// Java: DefaultPassConfig.addNonTypedAstNormalizationPasses(PassListBuilder)
func (c *DefaultPassConfig) addNonTypedAstNormalizationPasses(passes *PassListBuilder) {
	passes.MaybeAdd(c.removeCastNodes)
	passes.MaybeAdd(c.typesToColors)
}

// Java: DefaultPassConfig.shouldRunRemoveUnusedCode()
func (c *DefaultPassConfig) shouldRunRemoveUnusedCode() bool {
	options := c.options
	return options.ShouldRemoveUnusedVariables() ||
		options.ShouldRemoveUnusedLocalVariables() ||
		options.ShouldRemoveUnusedPrototypeProperties() ||
		options.IsRemoveUnusedClassProperties() ||
		options.GetRewritePolyfills()
}

// assertValidOrderForChecks: certain checks and rewriting passes need to
// run in a particular order. For example, the PolymerPass will not work
// correctly unless it runs after the goog.provide() processing. This
// enforces those constraints.
//
// checks is the list of check passes.
//
// Java: DefaultPassConfig.assertValidOrderForChecks(PassListBuilder)
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
	// TODO(lharker): add this back once fixing the ProcessCommonJSModules VarCheck branch
	// checks.assertPassOrder(
	//     checkVariableReferences,
	//     closureGoogScopeAliases,
	//     "Variable checking must happen before goog.scope processing.");

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

// assertValidOrderForOptimizations: certain optimizations need to run in a
// particular order. For example, OptimizeCalls must run before
// RemoveSuperMethodsPass, because the former can invalidate assumptions in
// the latter. This enforces those constraints.
//
// optimizations is the list of optimization passes.
//
// Java: DefaultPassConfig.assertValidOrderForOptimizations(PassListBuilder)
func (c *DefaultPassConfig) assertValidOrderForOptimizations(optimizations *PassListBuilder) {
	optimizations.AssertPassOrder(
		c.j2clPass,
		TranspilationPasses_rewriteGenerators,
		"J2CL normalization should be done before generator re-writing.")

	optimizations.AssertPassOrder(
		TranspilationPasses_rewritePolyfills,
		c.processDefinesOptimize,
		"Polyfill injection must be done before processDefines as some polyfills reference "+
			"goog.defines.")
	optimizations.AssertPassOrder(
		TranspilationPasses_injectTranspilationRuntimeLibraries,
		c.processDefinesOptimize,
		"Runtime library injection must be done before processDefines some runtime libraries "+
			"reference goog.defines.")

	optimizations.AssertPassOrder(
		c.processDefinesOptimize,
		c.j2clUtilGetDefineRewriterPass,
		"J2CL define re-writing should be done after processDefines since it relies on "+
			"Compiler#getDefineNames to have been populated by it.")

	optimizations.AssertPassOrder(
		c.removeUnusedCode,
		c.isolatePolyfills,
		"Polyfill isolation should be done after RemovedUnusedCode. Otherwise unused polyfill"+
			" removal will not find any polyfill usages and will delete all polyfills.")

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

// Java: DefaultPassConfig.GENERATE_EXPORTS_ERROR
var defaultPassConfig_GENERATE_EXPORTS_ERROR = DiagnosticType_Error(
	"JSC_GENERATE_EXPORTS_ERROR",
	"Exports can only be generated if export symbol/property functions are set.")
