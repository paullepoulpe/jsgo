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
	"fmt"
	"io"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/parsing"
)

// Ensures that functions labelled with `@requireInlining` JSDoc are, in fact, actually inlined by
// the optimizer.
//
// This may only be set when using advanced optimizations, as otherwise the compiler is unable
// to do most inlining. The compiler will throw an exception during initialization if this is set
// to true, but advanced optimizations are not being used.
//
// Java: CompilerOptions.setValidateRequiredInlinings(boolean)
func (o *CompilerOptions) SetValidateRequiredInlinings(validateRequireInliningAnnotation bool) {
	o.validateRequireInliningAnnotation = base.Tri_ForBoolean(validateRequireInliningAnnotation)
}

// Java: CompilerOptions.getShouldValidateRequiredInlinings()
func (o *CompilerOptions) GetShouldValidateRequiredInlinings() base.Tri {
	return o.validateRequireInliningAnnotation
}

// Java: CompilerOptions.setPreserveDetailedSourceInfo(boolean)
func (o *CompilerOptions) SetPreserveDetailedSourceInfo(preserveDetailedSourceInfo bool) {
	o.preserveDetailedSourceInfo = preserveDetailedSourceInfo
}

// Java: CompilerOptions.preservesDetailedSourceInfo()
func (o *CompilerOptions) PreservesDetailedSourceInfo() bool {
	return o.preserveDetailedSourceInfo
}

// Java: CompilerOptions.setPreserveNonJSDocComments(boolean)
func (o *CompilerOptions) SetPreserveNonJSDocComments(preserveNonJSDocComments bool) {
	o.preserveNonJSDocComments = preserveNonJSDocComments
}

// Java: CompilerOptions.getPreserveNonJSDocComments()
func (o *CompilerOptions) GetPreserveNonJSDocComments() bool {
	return o.preserveNonJSDocComments
}

// Java: CompilerOptions.setContinueAfterErrors(boolean)
func (o *CompilerOptions) SetContinueAfterErrors(continueAfterErrors bool) {
	o.continueAfterErrors = continueAfterErrors
}

// Java: CompilerOptions.canContinueAfterErrors()
func (o *CompilerOptions) CanContinueAfterErrors() bool {
	return o.continueAfterErrors
}

// Enables or disables the parsing of JSDoc documentation, and optionally also the preservation of
// all whitespace and formatting within a JSDoc comment. By default, whitespace is collapsed for
// all comments except {@literal @license} and {@literal @preserve} blocks,
//
// Java: CompilerOptions.setParseJsDocDocumentation(Config.JsDocParsing)
func (o *CompilerOptions) SetParseJsDocDocumentation(parseJsDocDocumentation parsing.Config_JsDocParsing) {
	o.parseJsDocDocumentation = parseJsDocDocumentation
}

// Checks JSDoc documentation will be parsed.
//
// @return True when JSDoc documentation will be parsed, false if not.
//
// Java: CompilerOptions.isParseJsDocDocumentation()
func (o *CompilerOptions) IsParseJsDocDocumentation() parsing.Config_JsDocParsing {
	return o.parseJsDocDocumentation
}

// Skip all passes (other than transpilation, if requested). Don't inject any runtime libraries
// (unless explicitly requested) or do any checks/optimizations (this is useful for per-file
// transpilation).
//
// Java: CompilerOptions.setSkipNonTranspilationPasses(boolean)
func (o *CompilerOptions) SetSkipNonTranspilationPasses(skipNonTranspilationPasses bool) {
	o.skipNonTranspilationPasses = skipNonTranspilationPasses
}

// Java: CompilerOptions.getSkipNonTranspilationPasses()
func (o *CompilerOptions) GetSkipNonTranspilationPasses() bool {
	return o.skipNonTranspilationPasses
}

// Java: CompilerOptions.setDevMode(DevMode)
func (o *CompilerOptions) SetDevMode(devMode CompilerOptions_DevMode) {
	o.devMode = devMode
}

// Java: CompilerOptions.getDevMode()
func (o *CompilerOptions) GetDevMode() CompilerOptions_DevMode {
	return o.devMode
}

// Java: CompilerOptions.setCheckDeterminism(boolean)
func (o *CompilerOptions) SetCheckDeterminism(checkDeterminism bool) {
	o.checkDeterminism = checkDeterminism
}

// Java: CompilerOptions.getCheckDeterminism()
func (o *CompilerOptions) GetCheckDeterminism() bool {
	return o.checkDeterminism
}

// Java: CompilerOptions.setMessageBundle(MessageBundle)
func (o *CompilerOptions) SetMessageBundle(messageBundle MessageBundle) {
	o.messageBundle = messageBundle
}

// Java: CompilerOptions.getMessageBundle()
func (o *CompilerOptions) GetMessageBundle() MessageBundle {
	return o.messageBundle
}

// Java: CompilerOptions.setCheckSymbols(boolean)
func (o *CompilerOptions) SetCheckSymbols(checkSymbols bool) {
	o.checkSymbols = checkSymbols
}

// Java: CompilerOptions.getCheckSymbols()
func (o *CompilerOptions) GetCheckSymbols() bool {
	return o.checkSymbols
}

// Java: CompilerOptions.setCheckSuspiciousCode(boolean)
func (o *CompilerOptions) SetCheckSuspiciousCode(checkSuspiciousCode bool) {
	o.checkSuspiciousCode = checkSuspiciousCode
}

// Java: CompilerOptions.getCheckSuspiciousCode()
func (o *CompilerOptions) GetCheckSuspiciousCode() bool {
	return o.checkSuspiciousCode
}

// Java: CompilerOptions.setCheckTypes(boolean)
func (o *CompilerOptions) SetCheckTypes(checkTypes bool) {
	o.checkTypes = checkTypes
}

// Java: CompilerOptions.getCheckTypes()
func (o *CompilerOptions) GetCheckTypes() bool {
	return o.checkTypes
}

// Java: CompilerOptions.setFoldConstants(boolean)
func (o *CompilerOptions) SetFoldConstants(foldConstants bool) {
	o.foldConstants = foldConstants
}

// Java: CompilerOptions.shouldFoldConstants()
func (o *CompilerOptions) ShouldFoldConstants() bool {
	return o.foldConstants
}

// Java: CompilerOptions.setDeadAssignmentElimination(boolean)
func (o *CompilerOptions) SetDeadAssignmentElimination(deadAssignmentElimination bool) {
	o.deadAssignmentElimination = deadAssignmentElimination
}

// Java: CompilerOptions.shouldRunDeadAssignmentElimination()
func (o *CompilerOptions) ShouldRunDeadAssignmentElimination() bool {
	return o.deadAssignmentElimination
}

// Java: CompilerOptions.setDeadPropertyAssignmentElimination(boolean)
func (o *CompilerOptions) SetDeadPropertyAssignmentElimination(deadPropertyAssignmentElimination bool) {
	o.deadPropertyAssignmentElimination = base.Tri_ForBoolean(deadPropertyAssignmentElimination)
}

// Java: CompilerOptions.shouldRunDeadPropertyAssignmentElimination()
func (o *CompilerOptions) ShouldRunDeadPropertyAssignmentElimination() bool {
	// By default, run dead property assignment elimination when both
	//   1. regular dead assignment elimiantion is enabled.
	//   2. the polymer pass is disabled.
	// This is because Polymer relies heavily on getters and setters. Those are very incompatible
	// eliminating seemingly "dead" property assignments. See b/28912819.
	return o.deadPropertyAssignmentElimination.ToBoolean(
		o.deadAssignmentElimination && !o.polymerPass)
}

// Java: CompilerOptions.setInlineConstantVars(boolean)
func (o *CompilerOptions) SetInlineConstantVars(inlineConstantVars bool) {
	o.inlineConstantVars = inlineConstantVars
}

// Java: CompilerOptions.shouldInlineConstantVars()
func (o *CompilerOptions) ShouldInlineConstantVars() bool {
	return o.inlineConstantVars
}

// Java: CompilerOptions.setCrossChunkCodeMotion(boolean)
func (o *CompilerOptions) SetCrossChunkCodeMotion(crossChunkCodeMotion bool) {
	o.crossChunkCodeMotion = crossChunkCodeMotion
}

// Java: CompilerOptions.setCrossChunkCodeMotionNoStubMethods(boolean)
func (o *CompilerOptions) SetCrossChunkCodeMotionNoStubMethods(crossChunkCodeMotionNoStubMethods bool) {
	o.crossChunkCodeMotionNoStubMethods = crossChunkCodeMotionNoStubMethods
}

// Java: CompilerOptions.getCrossChunkCodeMotionNoStubMethods()
func (o *CompilerOptions) GetCrossChunkCodeMotionNoStubMethods() bool {
	return o.crossChunkCodeMotionNoStubMethods
}

// Java: CompilerOptions.setParentChunkCanSeeSymbolsDeclaredInChildren(boolean)
func (o *CompilerOptions) SetParentChunkCanSeeSymbolsDeclaredInChildren(parentChunkCanSeeSymbolsDeclaredInChildren bool) {
	o.parentChunkCanSeeSymbolsDeclaredInChildren = parentChunkCanSeeSymbolsDeclaredInChildren
}

// Java: CompilerOptions.getParentChunkCanSeeSymbolsDeclaredInChildren()
func (o *CompilerOptions) GetParentChunkCanSeeSymbolsDeclaredInChildren() bool {
	return o.parentChunkCanSeeSymbolsDeclaredInChildren
}

// Java: CompilerOptions.setCrossChunkMethodMotion(boolean)
func (o *CompilerOptions) SetCrossChunkMethodMotion(crossChunkMethodMotion bool) {
	o.crossChunkMethodMotion = crossChunkMethodMotion
}

// Java: CompilerOptions.getCrossChunkMethodMotion()
func (o *CompilerOptions) GetCrossChunkMethodMotion() bool {
	return o.crossChunkMethodMotion
}

// Java: CompilerOptions.setCoalesceVariableNames(boolean)
func (o *CompilerOptions) SetCoalesceVariableNames(coalesceVariableNames bool) {
	o.coalesceVariableNames = coalesceVariableNames
}

// Java: CompilerOptions.shouldCoalesceVariableNames()
func (o *CompilerOptions) ShouldCoalesceVariableNames() bool {
	return o.coalesceVariableNames
}

// Java: CompilerOptions.setOptimizeLetAndConst(boolean)
func (o *CompilerOptions) SetOptimizeLetAndConst(optimizeLetAndConst bool) {
	o.optimizeLetAndConst = optimizeLetAndConst
}

// Java: CompilerOptions.shouldOptimizeLetAndConst()
func (o *CompilerOptions) ShouldOptimizeLetAndConst() bool {
	return o.optimizeLetAndConst
}

// Java: CompilerOptions.setAssumeGlobalScopeIsIsolated(boolean)
func (o *CompilerOptions) SetAssumeGlobalScopeIsIsolated(assumeGlobalScopeIsIsolated bool) {
	o.assumeGlobalScopeIsIsolated = assumeGlobalScopeIsIsolated
}

// Java: CompilerOptions.shouldTreatGlobalScopeAsIsolated()
func (o *CompilerOptions) ShouldTreatGlobalScopeAsIsolated() bool {
	return o.assumeGlobalScopeIsIsolated || o.renamePrefixNamespace != nil
}

// Java: CompilerOptions.setInlineLocalVariables(boolean)
func (o *CompilerOptions) SetInlineLocalVariables(inlineLocalVariables bool) {
	o.inlineLocalVariables = inlineLocalVariables
}

// Java: CompilerOptions.shouldInlineLocalVariables()
func (o *CompilerOptions) ShouldInlineLocalVariables() bool {
	return o.inlineLocalVariables
}

// Java: CompilerOptions.setFlowSensitiveInlineVariables(boolean)
func (o *CompilerOptions) SetFlowSensitiveInlineVariables(enabled bool) {
	o.flowSensitiveInlineVariables = enabled
}

// Java: CompilerOptions.getFlowSensitiveInlineVariables()
func (o *CompilerOptions) GetFlowSensitiveInlineVariables() bool {
	return o.flowSensitiveInlineVariables
}

// Java: CompilerOptions.setSmartNameRemoval(boolean)
func (o *CompilerOptions) SetSmartNameRemoval(smartNameRemoval bool) {
	// TODO(bradfordcsmith): Remove the smart name removal option.
	o.smartNameRemoval = smartNameRemoval
	if smartNameRemoval {
		// To get the effect this option used to have we need to enable these options.
		// Don't disable them here if they were set explicitly, though.
		o.removeUnusedVars = true
		o.removeUnusedPrototypeProperties = true
	}
}

// Java: CompilerOptions.getSmartNameRemoval()
func (o *CompilerOptions) GetSmartNameRemoval() bool {
	return o.smartNameRemoval
}

// Java: CompilerOptions.setExtractPrototypeMemberDeclarations(boolean)
func (o *CompilerOptions) SetExtractPrototypeMemberDeclarations(enabled bool) {
	if enabled {
		o.extractPrototypeMemberDeclarations = CompilerOptions_ExtractPrototypeMemberDeclarationsMode_USE_GLOBAL_TEMP
	} else {
		o.extractPrototypeMemberDeclarations = CompilerOptions_ExtractPrototypeMemberDeclarationsMode_OFF
	}
}

// USE_IIFE is currently unused. Consider removing support for it and
// deleting this setter.
//
// Java: CompilerOptions.setExtractPrototypeMemberDeclarations(ExtractPrototypeMemberDeclarationsMode)
func (o *CompilerOptions) SetExtractPrototypeMemberDeclarationsOfMode(mode CompilerOptions_ExtractPrototypeMemberDeclarationsMode) {
	o.extractPrototypeMemberDeclarations = mode
}

// Java: CompilerOptions.getExtractPrototypeMemberDeclarationsMode()
func (o *CompilerOptions) GetExtractPrototypeMemberDeclarationsMode() CompilerOptions_ExtractPrototypeMemberDeclarationsMode {
	return o.extractPrototypeMemberDeclarations
}

// Java: CompilerOptions.setRemoveUnusedPrototypeProperties(boolean)
func (o *CompilerOptions) SetRemoveUnusedPrototypeProperties(enabled bool) {
	o.removeUnusedPrototypeProperties = enabled
	// InlineSimpleMethods makes similar assumptions to
	// RemoveUnusedCode, so they are enabled together.
	o.inlineGetters = enabled
}

// Java: CompilerOptions.shouldRemoveUnusedPrototypeProperties()
func (o *CompilerOptions) ShouldRemoveUnusedPrototypeProperties() bool {
	return o.removeUnusedPrototypeProperties
}

// Java: CompilerOptions.shouldInlineGetters()
func (o *CompilerOptions) ShouldInlineGetters() bool {
	return o.inlineGetters
}

// Java: CompilerOptions.setCollapseVariableDeclarations(boolean)
func (o *CompilerOptions) SetCollapseVariableDeclarations(enabled bool) {
	o.collapseVariableDeclarations = enabled
}

// Java: CompilerOptions.shouldCollapseVariableDeclarations()
func (o *CompilerOptions) ShouldCollapseVariableDeclarations() bool {
	return o.collapseVariableDeclarations
}

// Java: CompilerOptions.setCollapseAnonymousFunctions(boolean)
func (o *CompilerOptions) SetCollapseAnonymousFunctions(enabled bool) {
	o.collapseAnonymousFunctions = enabled
}

// Java: CompilerOptions.shouldCollapseAnonymousFunctions()
func (o *CompilerOptions) ShouldCollapseAnonymousFunctions() bool {
	return o.collapseAnonymousFunctions
}

// Java: CompilerOptions.setAliasStringsMode(AliasStringsMode)
func (o *CompilerOptions) SetAliasStringsMode(aliasStringsMode CompilerOptions_AliasStringsMode) {
	o.aliasStringsMode = aliasStringsMode
}

// Java: CompilerOptions.getAliasStringsMode()
func (o *CompilerOptions) GetAliasStringsMode() CompilerOptions_AliasStringsMode {
	return o.aliasStringsMode
}

// Java: CompilerOptions.setOutputJsStringUsage(boolean)
func (o *CompilerOptions) SetOutputJsStringUsage(outputJsStringUsage bool) {
	o.outputJsStringUsage = outputJsStringUsage
}

// Java: CompilerOptions.shouldOutputJsStringUsage()
func (o *CompilerOptions) ShouldOutputJsStringUsage() bool {
	return o.outputJsStringUsage
}

// Java: CompilerOptions.setConvertToDottedProperties(boolean)
func (o *CompilerOptions) SetConvertToDottedProperties(convertToDottedProperties bool) {
	o.convertToDottedProperties = convertToDottedProperties
}

// Java: CompilerOptions.shouldConvertToDottedProperties()
func (o *CompilerOptions) ShouldConvertToDottedProperties() bool {
	return o.convertToDottedProperties
}

// Java: CompilerOptions.setUseTypesForLocalOptimization(boolean)
func (o *CompilerOptions) SetUseTypesForLocalOptimization(useTypesForLocalOptimization bool) {
	o.useTypesForLocalOptimization = useTypesForLocalOptimization
}

// Java: CompilerOptions.shouldUseTypesForLocalOptimization()
func (o *CompilerOptions) ShouldUseTypesForLocalOptimization() bool {
	return o.useTypesForLocalOptimization
}

// Java: CompilerOptions.setUseTypesForOptimization(boolean)
func (o *CompilerOptions) SetUseTypesForOptimization(useTypesForOptimization bool) {
	if useTypesForOptimization {
		o.disambiguateProperties = true
		o.ambiguateProperties = true
		o.inlineProperties = true
		o.useTypesForLocalOptimization = true
	}
}

// Whether any type-based optimizations are enabled
//
// Java: CompilerOptions.requiresTypesForOptimization()
func (o *CompilerOptions) RequiresTypesForOptimization() bool {
	// Add new type-based optimization options to this list, or types may be erased from the AST
	// before the optimziations run.
	return o.disambiguateProperties ||
		o.ambiguateProperties ||
		o.inlineProperties ||
		o.useTypesForLocalOptimization ||
		(o.closurePass && (o.removeAbstractMethods || o.removeClosureAsserts))
}

// Java: CompilerOptions.setRewriteFunctionExpressions(boolean)
func (o *CompilerOptions) SetRewriteFunctionExpressions(rewriteFunctionExpressions bool) {
	o.rewriteFunctionExpressions = rewriteFunctionExpressions
}

// Java: CompilerOptions.shouldRewriteFunctionExpressions()
func (o *CompilerOptions) ShouldRewriteFunctionExpressions() bool {
	return o.rewriteFunctionExpressions
}

// Java: CompilerOptions.setOptimizeCalls(boolean)
func (o *CompilerOptions) SetOptimizeCalls(optimizeCalls bool) {
	o.optimizeCalls = optimizeCalls
}

// Java: CompilerOptions.shouldOptimizeCalls()
func (o *CompilerOptions) ShouldOptimizeCalls() bool {
	return o.optimizeCalls
}

// Java: CompilerOptions.getOptimizeESClassConstructors()
func (o *CompilerOptions) GetOptimizeESClassConstructors() bool {
	return o.optimizeESClassConstructors
}

// Java: CompilerOptions.setOptimizeESClassConstructors(boolean)
func (o *CompilerOptions) SetOptimizeESClassConstructors(optimizeESClassConstructors bool) {
	o.optimizeESClassConstructors = optimizeESClassConstructors
}

// Java: CompilerOptions.setVariableRenaming(VariableRenamingPolicy)
func (o *CompilerOptions) SetVariableRenaming(variableRenaming VariableRenamingPolicy) {
	o.variableRenaming = variableRenaming
}

// Java: CompilerOptions.getVariableRenaming()
func (o *CompilerOptions) GetVariableRenaming() VariableRenamingPolicy {
	return o.variableRenaming
}

// Java: CompilerOptions.setPropertyRenaming(PropertyRenamingPolicy)
func (o *CompilerOptions) SetPropertyRenaming(propertyRenaming PropertyRenamingPolicy) {
	o.propertyRenaming = propertyRenaming
}

// Java: CompilerOptions.getPropertyRenaming()
func (o *CompilerOptions) GetPropertyRenaming() PropertyRenamingPolicy {
	return o.propertyRenaming
}

// Java: CompilerOptions.setLabelRenaming(boolean)
func (o *CompilerOptions) SetLabelRenaming(labelRenaming bool) {
	o.labelRenaming = labelRenaming
}

// Java: CompilerOptions.shouldRenameLabels()
func (o *CompilerOptions) ShouldRenameLabels() bool {
	return o.labelRenaming
}

// Java: CompilerOptions.setReserveRawExports(boolean)
func (o *CompilerOptions) SetReserveRawExports(reserveRawExports bool) {
	o.reserveRawExports = reserveRawExports
}

// Java: CompilerOptions.shouldReserveRawExports()
func (o *CompilerOptions) ShouldReserveRawExports() bool {
	return o.reserveRawExports
}

// Java: CompilerOptions.setPreferStableNames(boolean)
func (o *CompilerOptions) SetPreferStableNames(preferStableNames bool) {
	o.preferStableNames = preferStableNames
}

// Java: CompilerOptions.shouldPreferStableNames()
func (o *CompilerOptions) ShouldPreferStableNames() bool {
	return o.preferStableNames
}

// Java: CompilerOptions.setGeneratePseudoNames(boolean)
func (o *CompilerOptions) SetGeneratePseudoNames(generatePseudoNames bool) {
	o.generatePseudoNames = generatePseudoNames
}

// Java: CompilerOptions.shouldGeneratePseudoNames()
func (o *CompilerOptions) ShouldGeneratePseudoNames() bool {
	return o.generatePseudoNames
}

// Java: CompilerOptions.setPropertyRenamingOnlyCompilationMode(boolean)
func (o *CompilerOptions) SetPropertyRenamingOnlyCompilationMode(propertyRenamingOnlyCompilationMode bool) {
	o.propertyRenamingOnlyCompilationMode = propertyRenamingOnlyCompilationMode
}

// Java: CompilerOptions.isPropertyRenamingOnlyCompilationMode()
func (o *CompilerOptions) IsPropertyRenamingOnlyCompilationMode() bool {
	return o.propertyRenamingOnlyCompilationMode
}

// Java: CompilerOptions.setRenamePrefix(String)
func (o *CompilerOptions) SetRenamePrefix(renamePrefix *string) {
	o.renamePrefix = renamePrefix
}

// Java: CompilerOptions.getRenamePrefix()
func (o *CompilerOptions) GetRenamePrefix() *string {
	return o.renamePrefix
}

// Java: CompilerOptions.getRenamePrefixNamespace()
func (o *CompilerOptions) GetRenamePrefixNamespace() *string {
	return o.renamePrefixNamespace
}

// Java: CompilerOptions.setRenamePrefixNamespace(String)
func (o *CompilerOptions) SetRenamePrefixNamespace(renamePrefixNamespace *string) {
	o.renamePrefixNamespace = renamePrefixNamespace
}

// Java: CompilerOptions.setCollapsePropertiesLevel(PropertyCollapseLevel)
func (o *CompilerOptions) SetCollapsePropertiesLevel(level CompilerOptions_PropertyCollapseLevel) {
	o.collapsePropertiesLevel = level
}

// Java: CompilerOptions.setCollapseProperties(boolean)
func (o *CompilerOptions) SetCollapseProperties(fullyCollapse bool) {
	if fullyCollapse {
		o.collapsePropertiesLevel = CompilerOptions_PropertyCollapseLevel_ALL
	} else {
		o.collapsePropertiesLevel = CompilerOptions_PropertyCollapseLevel_NONE
	}
}

// Java: CompilerOptions.getCollapsePropertiesLevel()
func (o *CompilerOptions) GetCollapsePropertiesLevel() CompilerOptions_PropertyCollapseLevel {
	return o.collapsePropertiesLevel
}

// Java: CompilerOptions.setDevirtualizeMethods(boolean)
func (o *CompilerOptions) SetDevirtualizeMethods(devirtualizeMethods bool) {
	o.devirtualizeMethods = devirtualizeMethods
}

// Java: CompilerOptions.shouldDevirtualizeMethods()
func (o *CompilerOptions) ShouldDevirtualizeMethods() bool {
	return o.devirtualizeMethods
}

// Java: CompilerOptions.setComputeFunctionSideEffects(boolean)
func (o *CompilerOptions) SetComputeFunctionSideEffects(computeFunctionSideEffects bool) {
	o.computeFunctionSideEffects = computeFunctionSideEffects
}

// Java: CompilerOptions.shouldComputeFunctionSideEffects()
func (o *CompilerOptions) ShouldComputeFunctionSideEffects() bool {
	return o.computeFunctionSideEffects
}

// Java: CompilerOptions.setDisambiguateProperties(boolean)
func (o *CompilerOptions) SetDisambiguateProperties(disambiguateProperties bool) {
	o.disambiguateProperties = disambiguateProperties
}

// Java: CompilerOptions.shouldDisambiguateProperties()
func (o *CompilerOptions) ShouldDisambiguateProperties() bool {
	return o.disambiguateProperties
}

// Java: CompilerOptions.setAmbiguateProperties(boolean)
func (o *CompilerOptions) SetAmbiguateProperties(ambiguateProperties bool) {
	o.ambiguateProperties = ambiguateProperties
}

// Java: CompilerOptions.shouldAmbiguateProperties()
func (o *CompilerOptions) ShouldAmbiguateProperties() bool {
	return o.ambiguateProperties
}

// Java: CompilerOptions.setInputVariableMap(VariableMap)
func (o *CompilerOptions) SetInputVariableMap(inputVariableMap *VariableMap) {
	o.inputVariableMap = inputVariableMap
}

// Java: CompilerOptions.getInputVariableMap()
func (o *CompilerOptions) GetInputVariableMap() *VariableMap {
	return o.inputVariableMap
}

// Java: CompilerOptions.setInputPropertyMap(VariableMap)
func (o *CompilerOptions) SetInputPropertyMap(inputPropertyMap *VariableMap) {
	o.inputPropertyMap = inputPropertyMap
}

// Java: CompilerOptions.getInputPropertyMap()
func (o *CompilerOptions) GetInputPropertyMap() *VariableMap {
	return o.inputPropertyMap
}

// Java: CompilerOptions.setExportTestFunctions(boolean)
func (o *CompilerOptions) SetExportTestFunctions(exportTestFunctions bool) {
	o.exportTestFunctions = exportTestFunctions
}

// Java: CompilerOptions.shouldExportTestFunctions()
func (o *CompilerOptions) ShouldExportTestFunctions() bool {
	return o.exportTestFunctions
}

// Java: CompilerOptions.setSyntheticBlockStartMarker(String)
func (o *CompilerOptions) SetSyntheticBlockStartMarker(syntheticBlockStartMarker *string) {
	o.syntheticBlockStartMarker = syntheticBlockStartMarker
}

// Java: CompilerOptions.setSyntheticBlockEndMarker(String)
func (o *CompilerOptions) SetSyntheticBlockEndMarker(syntheticBlockEndMarker *string) {
	o.syntheticBlockEndMarker = syntheticBlockEndMarker
}

// Java: CompilerOptions.getSyntheticBlockStartMarker()
func (o *CompilerOptions) GetSyntheticBlockStartMarker() *string {
	return o.syntheticBlockStartMarker
}

// Java: CompilerOptions.getSyntheticBlockEndMarker()
func (o *CompilerOptions) GetSyntheticBlockEndMarker() *string {
	return o.syntheticBlockEndMarker
}

// Java: CompilerOptions.setLocale(String)
func (o *CompilerOptions) SetLocale(locale *string) {
	o.locale = locale
}

// Java: CompilerOptions.getLocale()
func (o *CompilerOptions) GetLocale() *string {
	return o.locale
}

// Java: CompilerOptions.setDoLateLocalization(boolean)
func (o *CompilerOptions) SetDoLateLocalization(doLateLocalization bool) {
	o.doLateLocalization = doLateLocalization
}

// Java: CompilerOptions.doLateLocalization()
func (o *CompilerOptions) DoLateLocalization() bool {
	return o.doLateLocalization
}

// Should we run any form of the `ReplaceMessages` pass?
//
// Java: CompilerOptions.shouldRunReplaceMessagesPass()
func (o *CompilerOptions) ShouldRunReplaceMessagesPass() bool {
	return !o.ShouldRunReplaceMessagesForChrome() && o.messageBundle != nil
}

// Java: CompilerOptions.setMarkAsCompiled(boolean)
func (o *CompilerOptions) SetMarkAsCompiled(markAsCompiled bool) {
	o.markAsCompiled = markAsCompiled
}

// Java: CompilerOptions.shouldMarkAsCompiled()
func (o *CompilerOptions) ShouldMarkAsCompiled() bool {
	return o.markAsCompiled
}

// Java: CompilerOptions.setClosurePass(boolean)
func (o *CompilerOptions) SetClosurePass(closurePass bool) {
	o.closurePass = closurePass
}

// Java: CompilerOptions.getClosurePass()
func (o *CompilerOptions) GetClosurePass() bool {
	return o.closurePass
}

// Preserve closure primitives.
//
// For now, this only preserves goog.provide(), goog.require() and goog.module() calls.
//
// Java: CompilerOptions.setPreserveClosurePrimitives(boolean)
func (o *CompilerOptions) SetPreserveClosurePrimitives(preserveClosurePrimitives bool) {
	o.preserveClosurePrimitives = preserveClosurePrimitives
}

// Java: CompilerOptions.shouldPreservesGoogProvidesAndRequires()
func (o *CompilerOptions) ShouldPreservesGoogProvidesAndRequires() bool {
	return o.preserveClosurePrimitives
}

// Java: CompilerOptions.shouldPreserveGoogModule()
func (o *CompilerOptions) ShouldPreserveGoogModule() bool {
	return o.preserveClosurePrimitives
}

// Do not process goog. intrinsics, such as goog.getCssName().
//
// Java: CompilerOptions.shouldPreserveGoogLibraryPrimitives()
func (o *CompilerOptions) ShouldPreserveGoogLibraryPrimitives() bool {
	return o.preserveClosurePrimitives
}

// Java: CompilerOptions.setPreserveTypeAnnotations(boolean)
func (o *CompilerOptions) SetPreserveTypeAnnotations(preserveTypeAnnotations bool) {
	o.preserveTypeAnnotations = preserveTypeAnnotations
}

// Java: CompilerOptions.shouldPreserveTypeAnnotations()
func (o *CompilerOptions) ShouldPreserveTypeAnnotations() bool {
	return o.preserveTypeAnnotations
}

// Java: CompilerOptions.setGentsMode(boolean)
func (o *CompilerOptions) SetGentsMode(gentsMode bool) {
	o.gentsMode = gentsMode
}

// Java: CompilerOptions.getGentsMode()
func (o *CompilerOptions) GetGentsMode() bool {
	return o.gentsMode
}

// Java: CompilerOptions.setGatherCssNames(boolean)
func (o *CompilerOptions) SetGatherCssNames(gatherCssNames bool) {
	o.gatherCssNames = gatherCssNames
}

// Java: CompilerOptions.shouldGatherCssNames()
func (o *CompilerOptions) ShouldGatherCssNames() bool {
	return o.gatherCssNames
}

// @deprecated StripCode is deprecated. Code should be designed to be removed by other means.
//
// Java: CompilerOptions.setStripTypes(Set<String>)
func (o *CompilerOptions) SetStripTypes(stripTypes *jutil.LinkedHashSet[string]) {
	o.stripTypes = jutil.NewLinkedHashSet(stripTypes.ToSlice()...)
}

// @deprecated StripCode is deprecated. Code should be designed to be removed by other means.
//
// Java: CompilerOptions.getStripTypes()
func (o *CompilerOptions) GetStripTypes() *jutil.LinkedHashSet[string] {
	return o.stripTypes
}

// @deprecated StripCode is deprecated. Code should be designed to be removed by other means.
//
// Java: CompilerOptions.setStripNameSuffixes(Set<String>)
func (o *CompilerOptions) SetStripNameSuffixes(stripNameSuffixes *jutil.LinkedHashSet[string]) {
	o.stripNameSuffixes = jutil.NewLinkedHashSet(stripNameSuffixes.ToSlice()...)
}

// Java: CompilerOptions.getStripNameSuffixes()
func (o *CompilerOptions) GetStripNameSuffixes() *jutil.LinkedHashSet[string] {
	return o.stripNameSuffixes
}

// @deprecated StripCode is deprecated. Code should be designed to be removed by other means.
//
// Java: CompilerOptions.setStripNamePrefixes(Set<String>)
func (o *CompilerOptions) SetStripNamePrefixes(stripNamePrefixes *jutil.LinkedHashSet[string]) {
	o.stripNamePrefixes = jutil.NewLinkedHashSet(stripNamePrefixes.ToSlice()...)
}

// Java: CompilerOptions.getStripNamePrefixes()
func (o *CompilerOptions) GetStripNamePrefixes() *jutil.LinkedHashSet[string] {
	return o.stripNamePrefixes
}

// Java: CompilerOptions.addCustomPass(CustomPassExecutionTime, CompilerPass)
func (o *CompilerOptions) AddCustomPass(time CustomPassExecutionTime, customPass CompilerPass) {
	if o.customPasses == nil {
		o.customPasses = jutil.NewSetMultimap[CustomPassExecutionTime, CompilerPass]()
	}
	o.customPasses.Put(time, customPass)
}

// Java: CompilerOptions.getCustomPassesAt(CustomPassExecutionTime)
func (o *CompilerOptions) GetCustomPassesAt(executionTime CustomPassExecutionTime) []CompilerPass {
	if o.customPasses == nil {
		return []CompilerPass{}
	}
	return o.customPasses.Get(executionTime).ToSlice()
}

// Java: CompilerOptions.getEnableZonesDefineName()
func (o *CompilerOptions) GetEnableZonesDefineName() *string {
	return o.enableZonesDefineName
}

// Java: CompilerOptions.setEnableZonesDefineName(String)
func (o *CompilerOptions) SetEnableZonesDefineName(enableZonesDefineName *string) {
	o.enableZonesDefineName = enableZonesDefineName
}

// Java: CompilerOptions.getZoneInputPattern()
func (o *CompilerOptions) GetZoneInputPattern() *jutil.Pattern {
	return o.zoneInputPattern
}

// Java: CompilerOptions.setZoneInputPattern(Pattern)
func (o *CompilerOptions) SetZoneInputPattern(zoneInputPattern *jutil.Pattern) {
	o.zoneInputPattern = zoneInputPattern
}

// Java: CompilerOptions.setRewriteGlobalDeclarationsForTryCatchWrapping(boolean)
func (o *CompilerOptions) SetRewriteGlobalDeclarationsForTryCatchWrapping(rewrite bool) {
	o.rewriteGlobalDeclarationsForTryCatchWrapping = rewrite
}

// Java: CompilerOptions.shouldRewriteGlobalDeclarationsForTryCatchWrapping()
func (o *CompilerOptions) ShouldRewriteGlobalDeclarationsForTryCatchWrapping() bool {
	return o.rewriteGlobalDeclarationsForTryCatchWrapping
}

// Java: CompilerOptions.setCssRenamingMap(CssRenamingMap)
func (o *CompilerOptions) SetCssRenamingMap(cssRenamingMap CssRenamingMap) {
	o.cssRenamingMap = cssRenamingMap
}

// Java: CompilerOptions.getCssRenamingMap()
func (o *CompilerOptions) GetCssRenamingMap() CssRenamingMap {
	return o.cssRenamingMap
}

// Java: CompilerOptions.setCssRenamingSkiplist(Set<String>)
func (o *CompilerOptions) SetCssRenamingSkiplist(skiplist *jutil.LinkedHashSet[string]) {
	o.cssRenamingSkiplist = jutil.NewLinkedHashSet(skiplist.ToSlice()...)
}

// Java: CompilerOptions.getCssRenamingSkiplist()
func (o *CompilerOptions) GetCssRenamingSkiplist() *jutil.LinkedHashSet[string] {
	return o.cssRenamingSkiplist
}

// Java: CompilerOptions.setReplaceStringsFunctionDescriptions(List<String>)
func (o *CompilerOptions) SetReplaceStringsFunctionDescriptions(replaceStringsFunctionDescriptions []string) {
	o.replaceStringsFunctionDescriptions = append([]string{}, replaceStringsFunctionDescriptions...)
}

// Java: CompilerOptions.getReplaceStringsFunctionDescriptions()
func (o *CompilerOptions) GetReplaceStringsFunctionDescriptions() []string {
	return o.replaceStringsFunctionDescriptions
}

// Java: CompilerOptions.setReplaceStringsPlaceholderToken(String)
func (o *CompilerOptions) SetReplaceStringsPlaceholderToken(replaceStringsPlaceholderToken string) {
	o.replaceStringsPlaceholderToken = replaceStringsPlaceholderToken
}

// Java: CompilerOptions.getReplaceStringsPlaceholderToken()
func (o *CompilerOptions) GetReplaceStringsPlaceholderToken() string {
	return o.replaceStringsPlaceholderToken
}

// Java: CompilerOptions.setPrettyPrint(boolean)
func (o *CompilerOptions) SetPrettyPrint(prettyPrint bool) {
	o.prettyPrint = prettyPrint
}

// Java: CompilerOptions.isPrettyPrint()
func (o *CompilerOptions) IsPrettyPrint() bool {
	return o.prettyPrint
}

// Java: CompilerOptions.setLineBreak(boolean)
func (o *CompilerOptions) SetLineBreak(lineBreak bool) {
	o.lineBreak = lineBreak
}

// Java: CompilerOptions.shouldAddLineBreak()
func (o *CompilerOptions) ShouldAddLineBreak() bool {
	return o.lineBreak
}

// Java: CompilerOptions.setPrintInputDelimiter(boolean)
func (o *CompilerOptions) SetPrintInputDelimiter(printInputDelimiter bool) {
	o.printInputDelimiter = printInputDelimiter
}

// Java: CompilerOptions.shouldPrintInputDelimiter()
func (o *CompilerOptions) ShouldPrintInputDelimiter() bool {
	return o.printInputDelimiter
}

// Java: CompilerOptions.setInputDelimiter(String)
func (o *CompilerOptions) SetInputDelimiter(inputDelimiter string) {
	o.inputDelimiter = inputDelimiter
}

// Java: CompilerOptions.getInputDelimiter()
func (o *CompilerOptions) GetInputDelimiter() string {
	return o.inputDelimiter
}

// Java: CompilerOptions.setDebugLogDirectory(Path)
func (o *CompilerOptions) SetDebugLogDirectory(dir *jutil.Path) {
	o.debugLogDirectory = dir
}

// Java: CompilerOptions.getDebugLogDirectory()
func (o *CompilerOptions) GetDebugLogDirectory() *jutil.Path {
	return o.debugLogDirectory
}

// Java: CompilerOptions.setDebugLogFilter(String)
func (o *CompilerOptions) SetDebugLogFilter(filter *string) {
	o.debugLogFilter = filter
}

// Java: CompilerOptions.getDebugLogFilter()
func (o *CompilerOptions) GetDebugLogFilter() *string {
	return o.debugLogFilter
}

// Java: CompilerOptions.shouldSerializeExtraDebugInfo()
func (o *CompilerOptions) ShouldSerializeExtraDebugInfo() bool {
	// NOTE: this is tied to debug logging (via checking getDebugLogDirectory()) for convenience,
	// and it's possible to split up the two options if this begins to trigger memory issues.
	return o.serializeExtraDebugInfo || o.GetDebugLogDirectory() != nil
}

// Java: CompilerOptions.setSerializeExtraDebugInfo(boolean)
func (o *CompilerOptions) SetSerializeExtraDebugInfo(serializeExtraDebugInfo bool) {
	o.serializeExtraDebugInfo = serializeExtraDebugInfo
}

// Java: CompilerOptions.setStateCompressionWrapper(Function<OutputStream, OutputStream>)
func (o *CompilerOptions) SetStateCompressionWrapper(wrapper func(io.Writer) io.Writer) {
	o.stateCompressionWrapper = wrapper
}

// Java: CompilerOptions.getStateCompressionWrapper()
func (o *CompilerOptions) GetStateCompressionWrapper() func(io.Writer) io.Writer {
	return o.stateCompressionWrapper
}

// Java: CompilerOptions.setStateDecompressionWrapper(Function<InputStream, InputStream>)
func (o *CompilerOptions) SetStateDecompressionWrapper(wrapper func(io.Reader) io.Reader) {
	o.stateDecompressionWrapper = wrapper
}

// Java: CompilerOptions.getStateDecompressionWrapper()
func (o *CompilerOptions) GetStateDecompressionWrapper() func(io.Reader) io.Reader {
	return o.stateDecompressionWrapper
}

// Java: CompilerOptions.setQuoteKeywordProperties(boolean)
func (o *CompilerOptions) SetQuoteKeywordProperties(quoteKeywordProperties bool) {
	o.quoteKeywordProperties = quoteKeywordProperties
}

// Java: CompilerOptions.shouldQuoteKeywordProperties()
func (o *CompilerOptions) ShouldQuoteKeywordProperties() bool {
	// Never quote properties in .i.js files
	if o.incrementalCheckMode == CompilerOptions_IncrementalCheckMode_GENERATE_IJS {
		return false
	}
	return o.quoteKeywordProperties || parser.FeatureSet_ES3.Contains(o.GetOutputFeatureSet())
}

// Java: CompilerOptions.setErrorFormat(ErrorFormat)
func (o *CompilerOptions) SetErrorFormat(errorFormat ErrorFormat) {
	o.errorFormat = errorFormat
}

// Java: CompilerOptions.getErrorFormat()
func (o *CompilerOptions) GetErrorFormat() ErrorFormat {
	return o.errorFormat
}

// Java: CompilerOptions.setWarningsGuard(ComposeWarningsGuard)
func (o *CompilerOptions) SetWarningsGuard(warningsGuard *ComposeWarningsGuard) {
	o.warningsGuard = warningsGuard
}

// Java: CompilerOptions.setLineLengthThreshold(int)
func (o *CompilerOptions) SetLineLengthThreshold(lineLengthThreshold int) {
	o.lineLengthThreshold = lineLengthThreshold
}

// Java: CompilerOptions.getLineLengthThreshold()
func (o *CompilerOptions) GetLineLengthThreshold() int {
	return o.lineLengthThreshold
}

// Java: CompilerOptions.setUseOriginalNamesInOutput(boolean)
func (o *CompilerOptions) SetUseOriginalNamesInOutput(useOriginalNamesInOutput bool) {
	o.useOriginalNamesInOutput = useOriginalNamesInOutput
}

// Java: CompilerOptions.getUseOriginalNamesInOutput()
func (o *CompilerOptions) GetUseOriginalNamesInOutput() bool {
	return o.useOriginalNamesInOutput
}

// Java: CompilerOptions.setExternExportsPath(String)
func (o *CompilerOptions) SetExternExportsPath(externExportsPath *string) {
	o.externExportsPath = externExportsPath
}

// Java: CompilerOptions.getExternExportsPath()
func (o *CompilerOptions) GetExternExportsPath() *string {
	return o.externExportsPath
}

// Java: CompilerOptions.setSourceMapOutputPath(String)
func (o *CompilerOptions) SetSourceMapOutputPath(sourceMapOutputPath *string) {
	o.sourceMapOutputPath = sourceMapOutputPath
}

// Java: CompilerOptions.setApplyInputSourceMaps(boolean)
func (o *CompilerOptions) SetApplyInputSourceMaps(applyInputSourceMaps bool) {
	o.applyInputSourceMaps = applyInputSourceMaps
}

// Java: CompilerOptions.getApplyInputSourceMaps()
func (o *CompilerOptions) GetApplyInputSourceMaps() bool {
	return o.applyInputSourceMaps
}

// Java: CompilerOptions.setResolveSourceMapAnnotations(boolean)
func (o *CompilerOptions) SetResolveSourceMapAnnotations(resolveSourceMapAnnotations bool) {
	o.resolveSourceMapAnnotations = resolveSourceMapAnnotations
}

// Java: CompilerOptions.getResolveSourceMapAnnotations()
func (o *CompilerOptions) GetResolveSourceMapAnnotations() bool {
	return o.resolveSourceMapAnnotations
}

// Java: CompilerOptions.setSourceMapIncludeSourcesContent(boolean)
func (o *CompilerOptions) SetSourceMapIncludeSourcesContent(sourceMapIncludeSourcesContent bool) {
	o.sourceMapIncludeSourcesContent = sourceMapIncludeSourcesContent
}

// Java: CompilerOptions.getSourceMapIncludeSourcesContent()
func (o *CompilerOptions) GetSourceMapIncludeSourcesContent() bool {
	return o.sourceMapIncludeSourcesContent
}

// Java: CompilerOptions.setParseInlineSourceMaps(boolean)
func (o *CompilerOptions) SetParseInlineSourceMaps(parseInlineSourceMaps bool) {
	o.parseInlineSourceMaps = parseInlineSourceMaps
}

// Java: CompilerOptions.getParseInlineSourceMaps()
func (o *CompilerOptions) GetParseInlineSourceMaps() bool {
	return o.parseInlineSourceMaps
}

// Java: CompilerOptions.setSourceMapDetailLevel(SourceMap.DetailLevel)
func (o *CompilerOptions) SetSourceMapDetailLevel(sourceMapDetailLevel SourceMap_DetailLevel) {
	o.sourceMapDetailLevel = sourceMapDetailLevel
}

// Java: CompilerOptions.getSourceMapDetailLevel()
func (o *CompilerOptions) GetSourceMapDetailLevel() SourceMap_DetailLevel {
	return o.sourceMapDetailLevel
}

// Java: CompilerOptions.setSourceMapFormat(SourceMap.Format)
func (o *CompilerOptions) SetSourceMapFormat(sourceMapFormat SourceMap_Format) {
	o.sourceMapFormat = sourceMapFormat
}

// Java: CompilerOptions.getSourceMapFormat()
func (o *CompilerOptions) GetSourceMapFormat() SourceMap_Format {
	return o.sourceMapFormat
}

// Java: CompilerOptions.setSourceMapLocationMappings(List<? extends SourceMap.LocationMapping>)
func (o *CompilerOptions) SetSourceMapLocationMappings(sourceMapLocationMappings []SourceMap_LocationMapping) {
	o.sourceMapLocationMappings = append([]SourceMap_LocationMapping{}, sourceMapLocationMappings...)
}

// Java: CompilerOptions.getSourceMapLocationMappings()
func (o *CompilerOptions) GetSourceMapLocationMappings() []SourceMap_LocationMapping {
	return o.sourceMapLocationMappings
}

// Rewrites CommonJS modules so that modules can be concatenated together, by renaming all globals
// to avoid conflicting with other modules.
//
// Java: CompilerOptions.setProcessCommonJSModules(boolean)
func (o *CompilerOptions) SetProcessCommonJSModules(processCommonJSModules bool) {
	o.processCommonJSModules = processCommonJSModules
}

// Java: CompilerOptions.getProcessCommonJSModules()
func (o *CompilerOptions) GetProcessCommonJSModules() bool {
	return o.processCommonJSModules
}

// Java: CompilerOptions.setEs6ModuleTranspilation(Es6ModuleTranspilation)
func (o *CompilerOptions) SetEs6ModuleTranspilation(value CompilerOptions_Es6ModuleTranspilation) {
	o.es6ModuleTranspilation = value
}

// Java: CompilerOptions.getEs6ModuleTranspilation()
func (o *CompilerOptions) GetEs6ModuleTranspilation() CompilerOptions_Es6ModuleTranspilation {
	return o.es6ModuleTranspilation
}

// Sets a path prefix for CommonJS modules (maps to #setModuleRoots(List)).
//
// Java: CompilerOptions.setCommonJSModulePathPrefix(String)
func (o *CompilerOptions) SetCommonJSModulePathPrefix(commonJSModulePathPrefix string) {
	o.SetModuleRoots([]string{commonJSModulePathPrefix})
}

// Sets the module roots.
//
// Java: CompilerOptions.setModuleRoots(List<String>)
func (o *CompilerOptions) SetModuleRoots(moduleRoots []string) {
	o.moduleRoots = append([]string{}, moduleRoots...)
}

// Java: CompilerOptions.getModuleRoots()
func (o *CompilerOptions) GetModuleRoots() []string {
	return o.moduleRoots
}

// Sets whether to rewrite polyfills.
//
// Java: CompilerOptions.setRewritePolyfills(boolean)
func (o *CompilerOptions) SetRewritePolyfills(rewritePolyfills bool) {
	o.rewritePolyfills = rewritePolyfills
}

// Java: CompilerOptions.getRewritePolyfills()
func (o *CompilerOptions) GetRewritePolyfills() bool {
	return o.rewritePolyfills
}

// Sets whether to isolate polyfills from the global scope.
//
// Java: CompilerOptions.setIsolatePolyfills(boolean)
func (o *CompilerOptions) SetIsolatePolyfills(isolatePolyfills bool) {
	o.isolatePolyfills = isolatePolyfills
	if o.isolatePolyfills {
		o.SetDefineToBooleanLiteral("$jscomp.ISOLATE_POLYFILLS", isolatePolyfills)
	}
}

// Java: CompilerOptions.getIsolatePolyfills()
func (o *CompilerOptions) GetIsolatePolyfills() bool {
	return o.isolatePolyfills
}

// Java: CompilerOptions.setInjectPolyfillsNewerThan(LanguageMode)
func (o *CompilerOptions) SetInjectPolyfillsNewerThan(injectPolyfillsNewerThan *CompilerOptions_LanguageMode) {
	o.injectPolyfillsNewerThan = injectPolyfillsNewerThan
}

// Java: CompilerOptions.getInjectPolyfillsNewerThan()
func (o *CompilerOptions) GetInjectPolyfillsNewerThan() *CompilerOptions_LanguageMode {
	return o.injectPolyfillsNewerThan
}

// Sets whether to transpile async functions and generators for AsyncContext.
//
// Java: CompilerOptions.setInstrumentAsyncContext(boolean)
func (o *CompilerOptions) SetInstrumentAsyncContext(instrumentAsyncContext bool) {
	o.instrumentAsyncContext = instrumentAsyncContext
	o.SetDefineToBooleanLiteral("$jscomp.INSTRUMENT_ASYNC_CONTEXT", instrumentAsyncContext)
}

// Java: CompilerOptions.getInstrumentAsyncContext()
func (o *CompilerOptions) GetInstrumentAsyncContext() bool {
	return o.instrumentAsyncContext
}

// Java: CompilerOptions.setEs6SubclassTranspilation(Es6SubclassTranspilation)
func (o *CompilerOptions) SetEs6SubclassTranspilation(es6SubclassTranspilation CompilerOptions_Es6SubclassTranspilation) {
	o.es6SubclassTranspilation = es6SubclassTranspilation
}

// Java: CompilerOptions.getEs6SubclassTranspilation()
func (o *CompilerOptions) GetEs6SubclassTranspilation() CompilerOptions_Es6SubclassTranspilation {
	return o.es6SubclassTranspilation
}

// Sets list of libraries to always inject, even if not needed.
//
// Java: CompilerOptions.setForceLibraryInjection(Iterable<String>)
func (o *CompilerOptions) SetForceLibraryInjection(libraries []string) {
	o.forceLibraryInjection = append([]string{}, libraries...)
}

// Java: CompilerOptions.getForceLibraryInjectionList()
func (o *CompilerOptions) GetForceLibraryInjectionList() []string {
	return o.forceLibraryInjection
}

// Sets the set of libraries to never inject, even if required.
//
// Java: CompilerOptions.setPreventLibraryInjection(boolean)
func (o *CompilerOptions) SetPreventLibraryInjection(preventLibraryInjection bool) {
	if preventLibraryInjection {
		o.runtimeLibraryMode = RuntimeJsLibManager_RuntimeLibraryMode_NO_OP
	} else {
		o.runtimeLibraryMode = RuntimeJsLibManager_RuntimeLibraryMode_INJECT
	}
}

// Java: CompilerOptions.setRuntimeLibraryMode(RuntimeLibraryMode)
func (o *CompilerOptions) SetRuntimeLibraryMode(runtimeLibraryMode RuntimeJsLibManager_RuntimeLibraryMode) {
	o.runtimeLibraryMode = runtimeLibraryMode
}

// Java: CompilerOptions.getRuntimeLibraryMode()
func (o *CompilerOptions) GetRuntimeLibraryMode() RuntimeJsLibManager_RuntimeLibraryMode {
	return o.runtimeLibraryMode
}

// Java: CompilerOptions.setUnusedImportsToRemove(ImmutableSet<String>)
func (o *CompilerOptions) SetUnusedImportsToRemove(unusedImportsToRemove *jutil.LinkedHashSet[string]) {
	o.unusedImportsToRemove = unusedImportsToRemove
}

// Java: CompilerOptions.getUnusedImportsToRemove()
func (o *CompilerOptions) GetUnusedImportsToRemove() *jutil.LinkedHashSet[string] {
	return o.unusedImportsToRemove
}

// Java: CompilerOptions.setInstrumentForCoverageOption(InstrumentOption)
func (o *CompilerOptions) SetInstrumentForCoverageOption(instrumentForCoverageOption CompilerOptions_InstrumentOption) {
	o.instrumentForCoverageOption = instrumentForCoverageOption
}

// Java: CompilerOptions.getInstrumentForCoverageOption()
func (o *CompilerOptions) GetInstrumentForCoverageOption() CompilerOptions_InstrumentOption {
	return o.instrumentForCoverageOption
}

// Sets the name for the global array which is used by PRODUCTION instrumentation. The array is
// declared during the instrumentation pass with the name provided through this setter.
//
// Java: CompilerOptions.setProductionInstrumentationArrayName(String)
func (o *CompilerOptions) SetProductionInstrumentationArrayName(productionInstrumentationArrayName string) {
	o.productionInstrumentationArrayName = productionInstrumentationArrayName
}

// Java: CompilerOptions.getProductionInstrumentationArrayName()
func (o *CompilerOptions) GetProductionInstrumentationArrayName() string {
	return o.productionInstrumentationArrayName
}

// Java: CompilerOptions.getConformanceConfigs()
func (o *CompilerOptions) GetConformanceConfigs() []*ConformanceConfig {
	return o.conformanceConfigs
}

// Both enable and configure conformance checks, if non-null.
//
// @deprecated See go/binary-level-conformance-deprecated.
//
// Java: CompilerOptions.setConformanceConfig(ConformanceConfig)
func (o *CompilerOptions) SetConformanceConfig(conformanceConfig *ConformanceConfig) {
	o.SetConformanceConfigs([]*ConformanceConfig{jutil.CheckNotNull(conformanceConfig)})
}

// Both enable and configure conformance checks, if non-null.
//
// @deprecated See go/binary-level-conformance-deprecated.
//
// Java: CompilerOptions.setConformanceConfigs(List<ConformanceConfig>)
func (o *CompilerOptions) SetConformanceConfigs(configs []*ConformanceConfig) {
	for _, c := range configs {
		jutil.CheckNotNull(c) // ImmutableList.copyOf rejects null elements.
	}
	o.conformanceConfigs = append([]*ConformanceConfig{}, configs...)
}

// Java: CompilerOptions.clearConformanceConfigs()
func (o *CompilerOptions) ClearConformanceConfigs() {
	o.conformanceConfigs = []*ConformanceConfig{}
}

// Whether the output should contain a 'use strict' directive.
//
// Java: CompilerOptions.shouldEmitUseStrict()
func (o *CompilerOptions) ShouldEmitUseStrict() bool {
	// Fall back to the language in's strictness if there is no output language explicitly set
	// for backwards compatibility.
	if o.emitUseStrict != nil {
		return *o.emitUseStrict
	}
	if o.languageOutIsDefaultStrict != nil {
		return *o.languageOutIsDefaultStrict
	}
	return o.languageIn.IsDefaultStrict()
}

// Java: CompilerOptions.setEmitUseStrict(boolean)
func (o *CompilerOptions) SetEmitUseStrict(emitUseStrict bool) *CompilerOptions {
	o.emitUseStrict = &emitUseStrict
	return o
}

// Java: CompilerOptions.getModuleResolutionMode()
func (o *CompilerOptions) GetModuleResolutionMode() ModuleLoader_ResolutionMode {
	return o.moduleResolutionMode
}

// Java: CompilerOptions.setModuleResolutionMode(ResolutionMode)
func (o *CompilerOptions) SetModuleResolutionMode(moduleResolutionMode ModuleLoader_ResolutionMode) {
	o.moduleResolutionMode = moduleResolutionMode
}

// Java: CompilerOptions.getBrowserResolverPrefixReplacements()
func (o *CompilerOptions) GetBrowserResolverPrefixReplacements() *jutil.LinkedHashMap[string, string] {
	return o.browserResolverPrefixReplacements
}

// Java: CompilerOptions.setBrowserResolverPrefixReplacements(ImmutableMap<String, String>)
func (o *CompilerOptions) SetBrowserResolverPrefixReplacements(browserResolverPrefixReplacements *jutil.LinkedHashMap[string, string]) {
	o.browserResolverPrefixReplacements = browserResolverPrefixReplacements
}

// Java: CompilerOptions.setPathEscaper(ModuleLoader.PathEscaper)
func (o *CompilerOptions) SetPathEscaper(pathEscaper ModuleLoader_PathEscaper) {
	o.pathEscaper = pathEscaper
}

// Java: CompilerOptions.getPathEscaper()
func (o *CompilerOptions) GetPathEscaper() ModuleLoader_PathEscaper {
	return o.pathEscaper
}

// Java: CompilerOptions.getPackageJsonEntryNames()
func (o *CompilerOptions) GetPackageJsonEntryNames() []string {
	return o.packageJsonEntryNames
}

// Java: CompilerOptions.setPackageJsonEntryNames(List<String>)
func (o *CompilerOptions) SetPackageJsonEntryNames(names []string) {
	o.packageJsonEntryNames = names
}

// Java: CompilerOptions.setUseSizeHeuristicToStopOptimizationLoop(boolean)
func (o *CompilerOptions) SetUseSizeHeuristicToStopOptimizationLoop(mayStopEarly bool) {
	o.useSizeHeuristicToStopOptimizationLoop = mayStopEarly
}

// Java: CompilerOptions.shouldUseSizeHeuristicToStopOptimizationLoop()
func (o *CompilerOptions) ShouldUseSizeHeuristicToStopOptimizationLoop() bool {
	return o.useSizeHeuristicToStopOptimizationLoop
}

// Java: CompilerOptions.setMaxOptimizationLoopIterations(int)
func (o *CompilerOptions) SetMaxOptimizationLoopIterations(maxIterations int) {
	o.optimizationLoopMaxIterations = maxIterations
}

// Java: CompilerOptions.getMaxOptimizationLoopIterations()
func (o *CompilerOptions) GetMaxOptimizationLoopIterations() int {
	return o.optimizationLoopMaxIterations
}

// Java: CompilerOptions.getChunkOutputType()
func (o *CompilerOptions) GetChunkOutputType() CompilerOptions_ChunkOutputType {
	return o.chunkOutputType
}

// Java: CompilerOptions.setChunkOutputType(ChunkOutputType)
func (o *CompilerOptions) SetChunkOutputType(chunkOutputType CompilerOptions_ChunkOutputType) {
	o.chunkOutputType = chunkOutputType
}

// Java: CompilerOptions.setStrictMessageReplacement(boolean)
func (o *CompilerOptions) SetStrictMessageReplacement(strictMessageReplacement bool) {
	o.strictMessageReplacement = strictMessageReplacement
}

// Java: CompilerOptions.getStrictMessageReplacement()
func (o *CompilerOptions) GetStrictMessageReplacement() bool {
	return o.strictMessageReplacement
}

// Java: CompilerOptions.toString()
func (o *CompilerOptions) String() string {
	h := jutil.MoreObjects_ToStringHelper("CompilerOptions").OmitNullValues()
	h.Add("aliasStringsMode", compilerOptions_valueOf(o.GetAliasStringsMode()))
	h.AddBool("ambiguateProperties", o.ambiguateProperties)
	h.AddBool("angularPass", o.angularPass)
	h.AddBool("assumeGettersArePure", o.assumeGettersArePure)
	h.AddBool("assumePropertiesAreStaticallyAnalyzable", o.assumePropertiesAreStaticallyAnalyzable)
	h.AddBool("assumeStrictThis", o.AssumeStrictThis())
	h.Add("browserResolverPrefixReplacements", compilerOptions_valueOf(o.browserResolverPrefixReplacements))
	h.AddBool("checkDeterminism", o.GetCheckDeterminism())
	h.AddBool("checkSuspiciousCode", o.checkSuspiciousCode)
	h.AddBool("checkSymbols", o.checkSymbols)
	h.AddBool("checkTypes", o.checkTypes)
	h.AddBool("checksOnly", o.checksOnly)
	h.Add("chunksToPrintAfterEachPassRegexList", compilerOptions_valueOf(o.chunksToPrintAfterEachPassRegexList))
	h.AddBool("closurePass", o.closurePass)
	h.AddBool("coalesceVariableNames", o.coalesceVariableNames)
	h.Add("codingConvention", compilerOptions_valueOf(o.GetCodingConvention()))
	h.AddBool("collapseAnonymousFunctions", o.collapseAnonymousFunctions)
	h.AddBool("collapseObjectLiterals", o.collapseObjectLiterals)
	h.Add("collapseProperties", compilerOptions_valueOf(o.collapsePropertiesLevel))
	h.AddBool("collapseVariableDeclarations", o.collapseVariableDeclarations)
	h.AddBool("colorizeErrorOutput", o.ShouldColorizeErrorOutput())
	h.AddBool("computeFunctionSideEffects", o.computeFunctionSideEffects)
	h.Add("conformanceConfigs", compilerOptions_valueOf(o.GetConformanceConfigs()))
	h.Add("conformanceRemoveRegexFromPath", compilerOptions_valueOf(compilerOptions_optionalOf(o.conformanceRemoveRegexFromPath)))
	h.Add("conformanceReportingMode", compilerOptions_valueOf(o.conformanceReportingMode))
	h.AddBool("continueAfterErrors", o.CanContinueAfterErrors())
	h.AddBool("convertToDottedProperties", o.convertToDottedProperties)
	h.AddBool("crossChunkCodeMotion", o.crossChunkCodeMotion)
	h.AddBool("crossChunkCodeMotionNoStubMethods", o.crossChunkCodeMotionNoStubMethods)
	h.AddBool("crossChunkMethodMotion", o.crossChunkMethodMotion)
	h.Add("cssRenamingMap", compilerOptions_valueOf(o.cssRenamingMap))
	h.Add("cssRenamingSkiplist", compilerOptions_valueOf(o.cssRenamingSkiplist))
	h.Add("customPasses", compilerOptions_valueOf(o.customPasses))
	h.AddBool("deadAssignmentElimination", o.deadAssignmentElimination)
	h.Add("debugLogDirectory", compilerOptions_valueOf(o.debugLogDirectory))
	h.Add("defineReplacements", compilerOptions_valueOf(o.GetDefineReplacements()))
	h.Add("dependencyOptions", compilerOptions_valueOf(o.GetDependencyOptions()))
	h.Add("devMode", compilerOptions_valueOf(o.devMode))
	h.AddBool("devirtualizeMethods", o.devirtualizeMethods)
	h.AddBool("disambiguateProperties", o.disambiguateProperties)
	h.Add("emitUseStrict", compilerOptions_valueOf(compilerOptions_optionalOf(o.emitUseStrict)))
	h.AddBool("enableModuleRewriting", o.enableModuleRewriting)
	h.Add("environment", compilerOptions_valueOf(o.GetEnvironment()))
	h.Add("errorFormat", compilerOptions_valueOf(o.errorFormat))
	h.Add("errorHandler", compilerOptions_valueOf(o.errorHandler))
	h.Add("es6ModuleTranspilation", compilerOptions_valueOf(o.es6ModuleTranspilation))
	h.Add("es6SubclassTranspilation", compilerOptions_valueOf(o.es6SubclassTranspilation))
	h.AddBool("exportLocalPropertyDefinitions", o.exportLocalPropertyDefinitions)
	h.AddBool("exportTestFunctions", o.exportTestFunctions)
	h.Add("externExportsPath", compilerOptions_valueOf(o.externExportsPath))
	h.Add("extraAnnotationNames", compilerOptions_valueOf(o.extraAnnotationNames))
	h.Add("extractPrototypeMemberDeclarations", compilerOptions_valueOf(o.extractPrototypeMemberDeclarations))
	h.Add("filesToPrintAfterEachPassRegexList", compilerOptions_valueOf(o.filesToPrintAfterEachPassRegexList))
	h.AddBool("flowSensitiveInlineVariables", o.flowSensitiveInlineVariables)
	h.AddBool("foldConstants", o.foldConstants)
	h.Add("forceLibraryInjection", compilerOptions_valueOf(o.forceLibraryInjection))
	h.AddBool("gatherCssNames", o.gatherCssNames)
	h.AddBool("generateExports", o.generateExports)
	h.AddBool("generatePseudoNames", o.generatePseudoNames)
	h.AddBool("generateTypedExterns", o.ShouldGenerateTypedExterns())
	h.Add("idGenerators", compilerOptions_valueOf(o.idGenerators))
	h.Add("idGeneratorsMapSerialized", compilerOptions_valueOf(o.idGeneratorsMapSerialized))
	h.Add("incrementalCheckMode", compilerOptions_valueOf(o.incrementalCheckMode))
	h.AddBool("inferConsts", o.inferConsts)
	h.AddBool("inferTypes", o.inferTypes)
	h.AddBool("inlineConstantVars", o.inlineConstantVars)
	h.Add("inlineFunctionsLevel", compilerOptions_valueOf(o.inlineFunctionsLevel))
	h.AddBool("inlineGetters", o.inlineGetters)
	h.AddBool("inlineLocalVariables", o.inlineLocalVariables)
	h.AddBool("inlineProperties", o.inlineProperties)
	h.AddBool("inlineVariables", o.inlineVariables)
	h.Add("inputDelimiter", compilerOptions_valueOf(o.inputDelimiter))
	h.Add("inputPropertyMap", compilerOptions_valueOf(o.inputPropertyMap))
	h.Add("inputSourceMaps", compilerOptions_valueOf(o.inputSourceMaps))
	h.Add("inputVariableMap", compilerOptions_valueOf(o.inputVariableMap))
	h.AddBool("instrumentForCoverageOnly", o.instrumentForCoverageOnly)
	h.Add("instrumentForCoverageOption", compilerOptions_valueOf(o.instrumentForCoverageOption.String()))
	h.AddBool("isolatePolyfills", o.isolatePolyfills)
	h.AddBool("j2clMinifierEnabled", o.j2clMinifierEnabled)
	h.Add("j2clMinifierPruningManifest", compilerOptions_valueOf(o.j2clMinifierPruningManifest))
	h.Add("j2clPassMode", compilerOptions_valueOf(o.j2clPassMode))
	h.AddBool("labelRenaming", o.labelRenaming)
	h.Add("languageIn", compilerOptions_valueOf(o.GetLanguageIn()))
	h.Add("languageOutIsDefaultStrict", compilerOptions_valueOf(compilerOptions_optionalOf(o.languageOutIsDefaultStrict)))
	h.AddBool("lineBreak", o.lineBreak)
	h.Add("lineLengthThreshold", o.lineLengthThreshold) // add(String, int): never omitted
	h.Add("locale", compilerOptions_valueOf(o.locale))
	h.AddBool("markAsCompiled", o.markAsCompiled)
	h.Add("maxFunctionSizeAfterInlining", o.maxFunctionSizeAfterInlining) // add(String, int): never omitted
	h.Add("messageBundle", compilerOptions_valueOf(o.messageBundle))
	h.Add("moduleRoots", compilerOptions_valueOf(o.moduleRoots))
	h.Add("nameGenerator", compilerOptions_valueOf(o.nameGenerator))
	h.Add("numParallelThreads", o.numParallelThreads) // add(String, int): never omitted
	h.AddBool("optimizeCalls", o.optimizeCalls)
	h.AddBool("optimizeESClassConstructors", o.optimizeESClassConstructors)
	h.Add(
		"optimizeLocalAccessForGlobalSymbolNamespace",
		compilerOptions_valueOf(o.optimizeLocalAccessForGlobalSymbolNamespace))
	h.Add("outputCharset", compilerOptions_valueOf(o.outputCharset))
	h.Add("outputFeatureSet", compilerOptions_valueOf(compilerOptions_optionalOf(o.outputFeatureSet)))
	h.Add("outputJs", compilerOptions_valueOf(o.outputJs))
	h.AddBool("outputJsStringUsage", o.outputJsStringUsage)
	h.AddBool("parentChunkCanSeeSymbolsDeclaredInChildren", o.parentChunkCanSeeSymbolsDeclaredInChildren)
	h.Add("parseJsDocDocumentation", compilerOptions_valueOf(o.IsParseJsDocDocumentation()))
	h.Add("pathEscaper", compilerOptions_valueOf(o.pathEscaper))
	h.AddBool("polymerPass", o.polymerPass)
	h.AddBool("preferSingleQuotes", o.preferSingleQuotes)
	h.AddBool("preferStableNames", o.preferStableNames)
	h.AddBool("preserveDetailedSourceInfo", o.PreservesDetailedSourceInfo())
	h.AddBool("preserveGoogProvidesAndRequires", o.preserveClosurePrimitives)
	h.AddBool("preserveNonJSDocComments", o.GetPreserveNonJSDocComments())
	h.AddBool("preserveTypeAnnotations", o.preserveTypeAnnotations)
	h.AddBool("prettyPrint", o.prettyPrint)
	h.AddBool("printConfig", o.printConfig)
	h.AddBool("printInputDelimiter", o.printInputDelimiter)
	h.AddBool("printSourceAfterEachPass", o.printSourceAfterEachPass)
	h.AddBool("processCommonJSModules", o.processCommonJSModules)
	h.Add("productionInstrumentationArrayName", compilerOptions_valueOf(o.productionInstrumentationArrayName))
	h.Add("propertiesThatMustDisambiguate", compilerOptions_valueOf(o.propertiesThatMustDisambiguate))
	h.Add("propertyRenaming", compilerOptions_valueOf(o.propertyRenaming))
	h.AddBool("propertyRenamingOnlyCompilationMode", o.propertyRenamingOnlyCompilationMode)
	h.AddBool("protectHiddenSideEffects", o.protectHiddenSideEffects)
	h.Add("qnameUsesToPrintAfterEachPassRegexList", compilerOptions_valueOf(o.qnameUsesToPrintAfterEachPassList))
	h.AddBool("quoteKeywordProperties", o.quoteKeywordProperties)
	h.AddBool("removeAbstractMethods", o.removeAbstractMethods)
	h.AddBool("removeClosureAsserts", o.removeClosureAsserts)
	h.AddBool("removeJ2clAsserts", o.removeJ2clAsserts)
	h.AddBool("removeUnusedClassProperties", o.removeUnusedClassProperties)
	h.AddBool("removeUnusedLocalVars", o.removeUnusedLocalVars)
	h.AddBool("removeUnusedPrototypeProperties", o.removeUnusedPrototypeProperties)
	h.AddBool("removeUnusedVars", o.removeUnusedVars)
	h.Add("renamePrefix", compilerOptions_valueOf(o.renamePrefix))
	h.Add("renamePrefixNamespace", compilerOptions_valueOf(o.renamePrefixNamespace))
	h.AddBool("renamePrefixNamespaceAssumeCrossChunkNames", o.renamePrefixNamespaceAssumeCrossChunkNames)
	h.AddBool("replaceIdGenerators", o.replaceIdGenerators)
	h.AddBool("replaceMessagesWithChromeI18n", o.replaceMessagesWithChromeI18n)
	h.Add("replaceStringsFunctionDescriptions", compilerOptions_valueOf(o.replaceStringsFunctionDescriptions))
	h.Add("replaceStringsPlaceholderToken", compilerOptions_valueOf(o.replaceStringsPlaceholderToken))
	h.AddBool("reserveRawExports", o.reserveRawExports)
	h.AddBool("rewriteFunctionExpressions", o.rewriteFunctionExpressions)
	h.AddBool("rewriteGlobalDeclarationsForTryCatchWrapping", o.rewriteGlobalDeclarationsForTryCatchWrapping)
	h.AddBool("rewriteModulesBeforeTypechecking", o.rewriteModulesBeforeTypechecking)
	h.AddBool("rewritePolyfills", o.rewritePolyfills)
	h.Add("runtimeLibraryMode", compilerOptions_valueOf(o.runtimeLibraryMode))
	h.AddBool("skipNonTranspilationPasses", o.skipNonTranspilationPasses)
	h.AddBool("smartNameRemoval", o.smartNameRemoval)
	h.Add("sourceMapDetailLevel", compilerOptions_valueOf(o.sourceMapDetailLevel))
	h.Add("sourceMapFormat", compilerOptions_valueOf(o.sourceMapFormat))
	h.Add("sourceMapLocationMappings", compilerOptions_valueOf(o.sourceMapLocationMappings))
	h.Add("sourceMapOutputPath", compilerOptions_valueOf(o.sourceMapOutputPath))
	h.AddBool("strictMessageReplacement", o.strictMessageReplacement)
	h.Add("stripNamePrefixes", compilerOptions_valueOf(o.stripNamePrefixes))
	h.Add("stripNameSuffixes", compilerOptions_valueOf(o.stripNameSuffixes))
	h.Add("stripTypes", compilerOptions_valueOf(o.stripTypes))
	h.Add("summaryDetailLevel", o.summaryDetailLevel) // add(String, int): never omitted
	h.Add("syntheticBlockEndMarker", compilerOptions_valueOf(o.syntheticBlockEndMarker))
	h.Add("syntheticBlockStartMarker", compilerOptions_valueOf(o.syntheticBlockStartMarker))
	h.Add("tcProjectId", compilerOptions_valueOf(o.tcProjectId))
	h.Add("tracer", compilerOptions_valueOf(o.tracer))
	h.AddBool("trustedStrings", o.trustedStrings)
	h.Add("tweakProcessing", compilerOptions_valueOf(o.GetTweakProcessing()))
	h.Add("unusedImportsToRemove", compilerOptions_valueOf(o.unusedImportsToRemove))
	h.AddBool("useTypesForLocalOptimization", o.useTypesForLocalOptimization)
	h.Add("variableRenaming", compilerOptions_valueOf(o.variableRenaming))
	h.Add("warningsGuard", compilerOptions_valueOf(o.GetWarningsGuard()))
	h.AddBool("wrapGoogModulesForWhitespaceOnly", o.wrapGoogModulesForWhitespaceOnly)
	return h.String()
}

// Java: CompilerOptions.expectStrictModeInput()
func (o *CompilerOptions) ExpectStrictModeInput() bool {
	if o.isStrictModeInput != nil {
		return *o.isStrictModeInput
	}
	return o.GetLanguageIn().IsDefaultStrict()
}

// Java: CompilerOptions.setStrictModeInput(boolean)
func (o *CompilerOptions) SetStrictModeInput(isStrictModeInput bool) *CompilerOptions {
	o.isStrictModeInput = &isStrictModeInput
	return o
}

// Java: CompilerOptions.getPropertyReservedNamingFirstChars()
func (o *CompilerOptions) GetPropertyReservedNamingFirstChars() *jutil.LinkedHashSet[uint16] {
	if o.polymerPass {
		return compilerOptions_POLYMER_PROPERTY_RESERVED_FIRST_CHARS
	} else if o.angularPass {
		return compilerOptions_ANGULAR_PROPERTY_RESERVED_FIRST_CHARS
	}
	return jutil.NewLinkedHashSet[uint16]()
}

// Java: CompilerOptions.getPropertyReservedNamingNonFirstChars()
func (o *CompilerOptions) GetPropertyReservedNamingNonFirstChars() *jutil.LinkedHashSet[uint16] {
	if o.polymerPass {
		return compilerOptions_POLYMER_PROPERTY_RESERVED_NON_FIRST_CHARS
	}
	return jutil.NewLinkedHashSet[uint16]()
}

// Java: CompilerOptions.shouldOptimize()
func (o *CompilerOptions) ShouldOptimize() bool {
	return !o.skipNonTranspilationPasses &&
		!o.checksOnly &&
		!o.ShouldGenerateTypedExterns() &&
		!o.instrumentForCoverageOnly
}

// compilerOptions_valueOf is the value ToStringHelper.add(String, Object)
// receives: nil for a Java null, otherwise its String.valueOf.
func compilerOptions_valueOf(value any) any {
	if jutil.IsNil(value) {
		return nil
	}
	return compilerOptions_toJavaString(value)
}

// compilerOptions_optional is a present or absent Guava Optional.
type compilerOptions_optional struct {
	value any // nil when absent
}

// compilerOptions_optionalOf wraps a nil-able pointer as an Optional.
func compilerOptions_optionalOf[T any](p *T) compilerOptions_optional {
	if p == nil {
		return compilerOptions_optional{}
	}
	return compilerOptions_optional{value: p}
}

func compilerOptions_formatSlice[T any](s []T) string {
	parts := make([]string, len(s))
	for i, elem := range s {
		parts[i] = compilerOptions_toJavaString(elem)
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

// compilerOptions_toJavaString is String.valueOf(Object) for the option
// values.
//
// Go-only: language-adapter — slices are printed like Java lists, and a
// value whose Java class does not override toString() prints as
// Object.toString() (class name, '@', identity hash).
func compilerOptions_toJavaString(value any) string {
	switch v := value.(type) {
	case compilerOptions_optional:
		if v.value == nil {
			return "Optional.absent()"
		}
		return "Optional.of(" + compilerOptions_toJavaString(v.value) + ")"
	case *bool:
		return jutil.ToJavaString(*v)
	case *string:
		return *v
	case string, bool, int, int32, int64, float64:
		return jutil.ToJavaString(v)
	case []string:
		return "[" + strings.Join(v, ", ") + "]"
	case []*ConformanceConfig:
		return compilerOptions_formatSlice(v)
	case []SourceMap_LocationMapping:
		return compilerOptions_formatSlice(v)
	case fmt.Stringer:
		return v.String()
	}
	typeName := fmt.Sprintf("%T", value)
	typeName = strings.TrimPrefix(typeName, "*")
	if dot := strings.LastIndexByte(typeName, '.'); dot >= 0 {
		typeName = typeName[dot+1:]
	}
	return jutil.ObjectToString(
		"com.google.javascript.jscomp."+strings.ReplaceAll(typeName, "_", "$"), value)
}
