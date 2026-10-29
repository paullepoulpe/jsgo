// Port of com/google/javascript/jscomp/CompilerOptions.java @ closure-compiler v20261021

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

// Java: CompilerOptions.UNLIMITED_FUN_SIZE_AFTER_INLINING
const CompilerOptions_UNLIMITED_FUN_SIZE_AFTER_INLINING = -1

// Skip all possible passes, to make the compiler as fast as possible.
//
// Java: CompilerOptions.skipAllCompilerPasses()
func (o *CompilerOptions) SkipAllCompilerPasses() {
	o.SetSkipNonTranspilationPasses(true)
}

// Sets the level of parallelism for compilation passes that can exploit multi-threading.
//
// Some compiler passes may take advantage of multi-threading, for example, parsing inputs.
// This sets the level of parallelism. The compiler will not start more than this number of
// threads.
//
// @param parallelism up to this number of parallel threads may be created.
//
// Java: CompilerOptions.setNumParallelThreads(int)
func (o *CompilerOptions) SetNumParallelThreads(parallelism int) {
	o.numParallelThreads = parallelism
}

// Java: CompilerOptions.getNumParallelThreads()
func (o *CompilerOptions) GetNumParallelThreads() int {
	return o.numParallelThreads
}

// SetFastMode enables or disables Go-to-Go deterministic fast mode.
func (o *CompilerOptions) SetFastMode(fastMode bool) {
	o.fastMode = fastMode
}

// isFastMode returns whether Go-to-Go deterministic fast mode is enabled.
func (o *CompilerOptions) isFastMode() bool {
	return o != nil && o.fastMode
}

// CompilerOptions_IsFastMode returns whether Go-to-Go deterministic fast mode is enabled.
func CompilerOptions_IsFastMode(o *CompilerOptions) bool {
	return o != nil && o.fastMode
}

// Java: CompilerOptions.setCodingConvention(CodingConvention)
func (o *CompilerOptions) SetCodingConvention(codingConvention CodingConvention) {
	o.codingConvention = codingConvention
}

// Java: CompilerOptions.getCodingConvention()
func (o *CompilerOptions) GetCodingConvention() CodingConvention {
	return o.codingConvention
}

// Should the compiler print its configuration options to stderr when they are initialized?
//
// Default false.
//
// Java: CompilerOptions.setPrintConfig(boolean)
func (o *CompilerOptions) SetPrintConfig(printConfig bool) {
	o.printConfig = printConfig
}

// Java: CompilerOptions.shouldPrintConfig()
func (o *CompilerOptions) ShouldPrintConfig() bool {
	return o.printConfig
}

// Java: CompilerOptions.getBrowserFeaturesetYearObject()
func (o *CompilerOptions) GetBrowserFeaturesetYearObject() *CompilerOptions_BrowserFeaturesetYear {
	return o.browserFeaturesetYear
}

// Java: CompilerOptions.setBrowserFeaturesetYear(int)
func (o *CompilerOptions) SetBrowserFeaturesetYear(year int) {
	o.SetBrowserFeaturesetYearOfBrowserFeaturesetYear(CompilerOptions_BrowserFeaturesetYear_From(year))
}

// Java: CompilerOptions.setBrowserFeaturesetYear(BrowserFeaturesetYear)
func (o *CompilerOptions) SetBrowserFeaturesetYearOfBrowserFeaturesetYear(year CompilerOptions_BrowserFeaturesetYear) {
	o.browserFeaturesetYear = year.Ptr()
	o.browserFeaturesetYear.SetDependentValuesFromYear(o)
}

// Sets ECMAScript version to use.
//
// Java: CompilerOptions.setLanguage(LanguageMode)
func (o *CompilerOptions) SetLanguage(language CompilerOptions_LanguageMode) {
	jutil.CheckState(language != CompilerOptions_LanguageMode_NO_TRANSPILE)
	o.SetLanguageIn(language)
	o.SetLanguageOut(language)
}

// Sets ECMAScript version to use for the input. If you are not transpiling from one version to
// another, use #setLanguage instead.
//
// Java: CompilerOptions.setLanguageIn(LanguageMode)
func (o *CompilerOptions) SetLanguageIn(languageIn CompilerOptions_LanguageMode) {
	jutil.CheckState(languageIn != CompilerOptions_LanguageMode_NO_TRANSPILE)
	if languageIn == CompilerOptions_LanguageMode_STABLE {
		o.languageIn = CompilerOptions_LanguageMode_STABLE_IN
	} else {
		o.languageIn = languageIn
	}
}

// Java: CompilerOptions.getLanguageIn()
func (o *CompilerOptions) GetLanguageIn() CompilerOptions_LanguageMode {
	return o.languageIn
}

// Sets ECMAScript version to use for the output.
//
// If you are not transpiling from one version to another, use #setLanguage instead.
//
// If you you need something more fine grained (e.g. "ES2017 without modules") use
// #setOutputFeatureSet.
//
// Java: CompilerOptions.setLanguageOut(LanguageMode)
func (o *CompilerOptions) SetLanguageOut(languageOut CompilerOptions_LanguageMode) {
	if languageOut == CompilerOptions_LanguageMode_NO_TRANSPILE {
		o.languageOutIsDefaultStrict = nil
		o.outputFeatureSet = nil
	} else {
		if languageOut == CompilerOptions_LanguageMode_STABLE {
			languageOut = CompilerOptions_LanguageMode_STABLE_OUT
		}
		isDefaultStrict := languageOut.IsDefaultStrict()
		o.languageOutIsDefaultStrict = &isDefaultStrict
		o.SetOutputFeatureSet(languageOut.ToFeatureSet())
	}
}

// Sets the features that allowed to appear in the output. Any feature in the input that is not in
// this output must be transpiled away.
//
// Note: this is a package private API since not every FeatureSet value can be properly output
// by the compiler without crashing. Both the `setBrowserFeaturesetYear` and `setLanguageOut` APIs
// are supported alternatives.
//
// Java: CompilerOptions.setOutputFeatureSet(FeatureSet)
func (o *CompilerOptions) SetOutputFeatureSet(featureSet *parser.FeatureSet) {
	o.outputFeatureSet = jutil.CheckNotNull(featureSet) // Optional.of
}

// Java: CompilerOptions.legacySetOutputFeatureSet(FeatureSet)
func (o *CompilerOptions) LegacySetOutputFeatureSet(featureSet *parser.FeatureSet) {
	o.SetOutputFeatureSet(featureSet)
}

// Gets the set of features that can appear in the output.
//
// Java: CompilerOptions.getOutputFeatureSet()
func (o *CompilerOptions) GetOutputFeatureSet() *parser.FeatureSet {
	if o.outputFeatureSet != nil {
		return o.outputFeatureSet
	}

	// Backwards compatibility for those that predate language out.
	return o.languageIn.ToFeatureSet()
}

// SetTranspilePublicClassFields disables emitting public class fields
// natively, forcing them to be transpiled into constructor assignments.
//
// Java: CompilerOptions.setTranspilePublicClassFields(boolean)
func (o *CompilerOptions) SetTranspilePublicClassFields(transpile bool) {
	if transpile {
		o.SetOutputFeatureSet(o.GetOutputFeatureSet().WithoutFeatures(parser.FeatureSet_Feature_PUBLIC_CLASS_FIELDS))
	} else {
		o.SetOutputFeatureSet(o.GetOutputFeatureSet().With(parser.FeatureSet_Feature_PUBLIC_CLASS_FIELDS))
	}
}

// Java: CompilerOptions.getTranspilePublicClassFields()
func (o *CompilerOptions) GetTranspilePublicClassFields() bool {
	return !o.GetOutputFeatureSet().Has(parser.FeatureSet_Feature_PUBLIC_CLASS_FIELDS)
}

// SetExperimentalOutputFeatureSet: a nil experimentalOutputFeatureSet is
// Java null (Optional.fromNullable).
//
// Java: CompilerOptions.setExperimentalOutputFeatureSet(ExperimentalOutputFeatureSet)
func (o *CompilerOptions) SetExperimentalOutputFeatureSet(experimentalOutputFeatureSet *CompilerOptions_ExperimentalOutputFeatureSet) {
	if experimentalOutputFeatureSet == nil {
		o.experimentalOutputFeatureSet = nil
	} else {
		v := *experimentalOutputFeatureSet
		o.experimentalOutputFeatureSet = &v
	}
}

// GetExperimentalOutputFeatureSet: nil is Optional.absent().
//
// Java: CompilerOptions.getExperimentalOutputFeatureSet()
func (o *CompilerOptions) GetExperimentalOutputFeatureSet() *CompilerOptions_ExperimentalOutputFeatureSet {
	return o.experimentalOutputFeatureSet
}

// Java: CompilerOptions.needsTranspilationFrom(FeatureSet)
func (o *CompilerOptions) NeedsTranspilationFrom(languageLevel *parser.FeatureSet) bool {
	// TODO(johnplaisted): This isn't really accurate. This should instead be the *parsed* language,
	// not the *input* language.
	return o.GetLanguageIn().ToFeatureSet().Contains(languageLevel) &&
		!o.GetOutputFeatureSet().Contains(languageLevel)
}

// Java: CompilerOptions.needsTranspilationOf(FeatureSet.Feature)
func (o *CompilerOptions) NeedsTranspilationOf(feature parser.FeatureSet_Feature) bool {
	return o.GetLanguageIn().ToFeatureSet().Has(feature) && !o.GetOutputFeatureSet().Has(feature)
}

// Set which set of builtin externs to use.
//
// Java: CompilerOptions.setEnvironment(Environment)
func (o *CompilerOptions) SetEnvironment(environment CompilerOptions_Environment) {
	o.environment = environment
}

// Java: CompilerOptions.getEnvironment()
func (o *CompilerOptions) GetEnvironment() CompilerOptions_Environment {
	return o.environment
}

// Sets the dependency management options.
//
// Java: CompilerOptions.setDependencyOptions(DependencyOptions)
func (o *CompilerOptions) SetDependencyOptions(dependencyOptions *DependencyOptions) {
	o.dependencyOptions = dependencyOptions
}

// Java: CompilerOptions.getDependencyOptions()
func (o *CompilerOptions) GetDependencyOptions() *DependencyOptions {
	return o.dependencyOptions
}

// Java: CompilerOptions.setMergedPrecompiledLibraries(boolean)
func (o *CompilerOptions) SetMergedPrecompiledLibraries(mergedPrecompiledLibraries bool) {
	o.mergedPrecompiledLibraries = mergedPrecompiledLibraries
}

// Java: CompilerOptions.getMergedPrecompiledLibraries()
func (o *CompilerOptions) GetMergedPrecompiledLibraries() bool {
	return o.mergedPrecompiledLibraries
}

// Java: CompilerOptions.setWrapGoogModulesForWhitespaceOnly(boolean)
func (o *CompilerOptions) SetWrapGoogModulesForWhitespaceOnly(enable bool) {
	o.wrapGoogModulesForWhitespaceOnly = enable
}

// Java: CompilerOptions.shouldWrapGoogModulesForWhitespaceOnly()
func (o *CompilerOptions) ShouldWrapGoogModulesForWhitespaceOnly() bool {
	return o.wrapGoogModulesForWhitespaceOnly
}

// Whether to enable the bad module rewriting before typechecking that we want to get rid of
//
// Java: CompilerOptions.setBadRewriteModulesBeforeTypecheckingThatWeWantToGetRidOf(boolean)
func (o *CompilerOptions) SetBadRewriteModulesBeforeTypecheckingThatWeWantToGetRidOf(b bool) {
	o.rewriteModulesBeforeTypechecking = b
}

// Java: CompilerOptions.shouldRewriteModulesBeforeTypechecking()
func (o *CompilerOptions) ShouldRewriteModulesBeforeTypechecking() bool {
	return o.enableModuleRewriting &&
		(o.rewriteModulesBeforeTypechecking || o.processCommonJSModules)
}

// Experimental option to disable all Closure and ES module and goog.provide rewriting
//
// Use at your own risk - disabling module rewriting is not fully tested yet.
//
// Java: CompilerOptions.setEnableModuleRewriting(boolean)
func (o *CompilerOptions) SetEnableModuleRewriting(enable bool) {
	o.enableModuleRewriting = enable
}

// Java: CompilerOptions.shouldRewriteModulesAfterTypechecking()
func (o *CompilerOptions) ShouldRewriteModulesAfterTypechecking() bool {
	return o.enableModuleRewriting && !o.rewriteModulesBeforeTypechecking
}

// Java: CompilerOptions.shouldRewriteModules()
func (o *CompilerOptions) ShouldRewriteModules() bool {
	return o.enableModuleRewriting
}

// Whether to enable support for dynamic import expressions
//
// Java: CompilerOptions.setAllowDynamicImport(boolean)
func (o *CompilerOptions) SetAllowDynamicImport(value bool) {
	o.allowDynamicImport = value
}

// Java: CompilerOptions.shouldAllowDynamicImport()
func (o *CompilerOptions) ShouldAllowDynamicImport() bool {
	return o.allowDynamicImport
}

// Set the alias name for dynamic import expressions
//
// Java: CompilerOptions.getDynamicImportAlias()
func (o *CompilerOptions) GetDynamicImportAlias() *string {
	return o.dynamicImportAlias
}

// Set the alias name for dynamic import expressions
//
// Java: CompilerOptions.setDynamicImportAlias(String)
func (o *CompilerOptions) SetDynamicImportAlias(value *string) {
	o.dynamicImportAlias = value
}

// Java: CompilerOptions.shouldAliasDynamicImport()
func (o *CompilerOptions) ShouldAliasDynamicImport() bool {
	return o.dynamicImportAlias != nil
}

// Java: CompilerOptions.shouldRunCrossChunkCodeMotion()
func (o *CompilerOptions) ShouldRunCrossChunkCodeMotion() bool {
	return o.crossChunkCodeMotion
}

// Java: CompilerOptions.shouldRunCrossChunkMethodMotion()
func (o *CompilerOptions) ShouldRunCrossChunkMethodMotion() bool {
	return o.crossChunkMethodMotion
}

// Java: CompilerOptions.getCrossChunkMethodMotion()
func (o *CompilerOptions) GetCrossChunkMethodMotion() bool {
	return o.crossChunkMethodMotion
}

// Java: CompilerOptions.getFlowSensitiveInlineVariables()
func (o *CompilerOptions) GetFlowSensitiveInlineVariables() bool {
	return o.flowSensitiveInlineVariables
}

// Java: CompilerOptions.shouldCollapseProperties()
func (o *CompilerOptions) ShouldCollapseProperties() bool {
	return o.collapsePropertiesLevel != CompilerOptions_PropertyCollapseLevel_NONE
}

// Java: CompilerOptions.getCollapsePropertiesLevel()
func (o *CompilerOptions) GetCollapsePropertiesLevel() CompilerOptions_PropertyCollapseLevel {
	return o.collapsePropertiesLevel
}

// Java: CompilerOptions.getPropertyCollapseLevel()
func (o *CompilerOptions) GetPropertyCollapseLevel() CompilerOptions_PropertyCollapseLevel {
	return o.collapsePropertiesLevel
}

// Java: CompilerOptions.setCollapseObjectLiterals(boolean)
func (o *CompilerOptions) SetCollapseObjectLiterals(enabled bool) {
	o.collapseObjectLiterals = enabled
}

// Java: CompilerOptions.getCollapseObjectLiterals()
func (o *CompilerOptions) GetCollapseObjectLiterals() bool {
	return o.collapseObjectLiterals
}

// @return Whether to attempt to remove unused class properties
//
// Java: CompilerOptions.isRemoveUnusedClassProperties()
func (o *CompilerOptions) IsRemoveUnusedClassProperties() bool {
	return o.removeUnusedClassProperties
}

// @param removeUnusedClassProperties Whether to attempt to remove unused class properties
//
// Java: CompilerOptions.setRemoveUnusedClassProperties(boolean)
func (o *CompilerOptions) SetRemoveUnusedClassProperties(removeUnusedClassProperties bool) {
	o.removeUnusedClassProperties = removeUnusedClassProperties
}

// Set the function inlining policy for the compiler.
//
// Java: CompilerOptions.setInlineFunctions(Reach)
func (o *CompilerOptions) SetInlineFunctions(reach CompilerOptions_Reach) {
	o.inlineFunctionsLevel = reach
}

// Get the function inlining policy for the compiler.
//
// Java: CompilerOptions.getInlineFunctionsLevel()
func (o *CompilerOptions) GetInlineFunctionsLevel() CompilerOptions_Reach {
	return o.inlineFunctionsLevel
}

// Java: CompilerOptions.setMaxFunctionSizeAfterInlining(int)
func (o *CompilerOptions) SetMaxFunctionSizeAfterInlining(funAstSize int) {
	jutil.CheckArgument(funAstSize > 0)
	o.maxFunctionSizeAfterInlining = funAstSize
}

// Java: CompilerOptions.getMaxFunctionSizeAfterInlining()
func (o *CompilerOptions) GetMaxFunctionSizeAfterInlining() int {
	return o.maxFunctionSizeAfterInlining
}

// Java: CompilerOptions.setInlineVariables(boolean)
func (o *CompilerOptions) SetInlineVariables(inlineVariables bool) {
	o.inlineVariables = inlineVariables
}

// Java: CompilerOptions.shouldInlineVariables()
func (o *CompilerOptions) ShouldInlineVariables() bool {
	return o.inlineVariables
}

// Set the variable inlining policy for the compiler.
//
// Java: CompilerOptions.setInlineVariables(Reach)
func (o *CompilerOptions) SetInlineVariablesOfReach(reach CompilerOptions_Reach) {
	switch reach {
	case CompilerOptions_Reach_ALL:
		o.inlineVariables = true
		o.inlineLocalVariables = true
	case CompilerOptions_Reach_LOCAL_ONLY:
		o.inlineVariables = false
		o.inlineLocalVariables = true
	case CompilerOptions_Reach_NONE:
		o.inlineVariables = false
		o.inlineLocalVariables = false
	}
}

// Set the function inlining policy for the compiler.
//
// Java: CompilerOptions.setInlineProperties(boolean)
func (o *CompilerOptions) SetInlineProperties(enable bool) {
	o.inlineProperties = enable
}

// Java: CompilerOptions.shouldInlineProperties()
func (o *CompilerOptions) ShouldInlineProperties() bool {
	return o.inlineProperties
}

// Set the variable removal policy for the compiler.
//
// Java: CompilerOptions.setRemoveUnusedVariables(Reach)
func (o *CompilerOptions) SetRemoveUnusedVariables(reach CompilerOptions_Reach) {
	switch reach {
	case CompilerOptions_Reach_ALL:
		o.removeUnusedVars = true
		o.removeUnusedLocalVars = true
	case CompilerOptions_Reach_LOCAL_ONLY:
		o.removeUnusedVars = false
		o.removeUnusedLocalVars = true
	case CompilerOptions_Reach_NONE:
		o.removeUnusedVars = false
		o.removeUnusedLocalVars = false
	}
}

// Java: CompilerOptions.shouldRemoveUnusedVariables()
func (o *CompilerOptions) ShouldRemoveUnusedVariables() bool {
	return o.removeUnusedVars
}

// Java: CompilerOptions.shouldRemoveUnusedLocalVariables()
func (o *CompilerOptions) ShouldRemoveUnusedLocalVariables() bool {
	return o.removeUnusedLocalVars
}

// Java: CompilerOptions.setRemoveAbstractMethods(boolean)
func (o *CompilerOptions) SetRemoveAbstractMethods(remove bool) {
	o.removeAbstractMethods = remove
}

// Java: CompilerOptions.shouldRemoveAbstractMethods()
func (o *CompilerOptions) ShouldRemoveAbstractMethods() bool {
	return o.removeAbstractMethods
}

// Java: CompilerOptions.setRemoveClosureAsserts(boolean)
func (o *CompilerOptions) SetRemoveClosureAsserts(remove bool) {
	o.removeClosureAsserts = remove
}

// Java: CompilerOptions.shouldRemoveClosureAsserts()
func (o *CompilerOptions) ShouldRemoveClosureAsserts() bool {
	return o.removeClosureAsserts
}

// Java: CompilerOptions.setRemoveJ2clAsserts(boolean)
func (o *CompilerOptions) SetRemoveJ2clAsserts(remove bool) {
	o.removeJ2clAsserts = remove
}

// Java: CompilerOptions.shouldRemoveJ2clAsserts()
func (o *CompilerOptions) ShouldRemoveJ2clAsserts() bool {
	return o.removeJ2clAsserts
}

// Sets how goog.tweak calls are processed.
//
// Java: CompilerOptions.setTweakProcessing(TweakProcessing)
func (o *CompilerOptions) SetTweakProcessing(tweakProcessing CompilerOptions_TweakProcessing) {
	o.tweakProcessing = tweakProcessing
}

// Java: CompilerOptions.getTweakProcessing()
func (o *CompilerOptions) GetTweakProcessing() CompilerOptions_TweakProcessing {
	return o.tweakProcessing
}
