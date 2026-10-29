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
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
)

type customPassGroup struct {
	time   CustomPassExecutionTime
	passes []CompilerPass
}

type customPassMap struct {
	groups []customPassGroup
}

func (m *customPassMap) get(time CustomPassExecutionTime) []CompilerPass {
	if m == nil {
		return []CompilerPass{}
	}
	for i := range m.groups {
		if m.groups[i].time == time {
			return slices.Clone(m.groups[i].passes)
		}
	}
	return []CompilerPass{}
}

func (m *customPassMap) String() string {
	if m == nil || len(m.groups) == 0 {
		return "{}"
	}
	var sb strings.Builder
	sb.WriteByte('{')
	for i, g := range m.groups {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(g.time.String())
		sb.WriteString("=[")
		for j, p := range g.passes {
			if j > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(jutil.ToJavaString(p))
		}
		sb.WriteByte(']')
	}
	sb.WriteByte('}')
	return sb.String()
}

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

// Java: CompilerOptions.getParentChunkCanSeeSymbolsDeclaredInChildren()
func (o *CompilerOptions) GetParentChunkCanSeeSymbolsDeclaredInChildren() bool {
	return o.parentChunkCanSeeSymbolsDeclaredInChildren
}

// Java: CompilerOptions.setCrossChunkMethodMotion(boolean)
func (o *CompilerOptions) SetCrossChunkMethodMotion(crossChunkMethodMotion bool) {
	o.crossChunkMethodMotion = crossChunkMethodMotion
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

// Java: CompilerOptions.shouldExportTestFunctions()
func (o *CompilerOptions) ShouldExportTestFunctions() bool {
	return o.exportTestFunctions
}

// Java: CompilerOptions.setSyntheticBlockStartMarker(String)
func (o *CompilerOptions) SetSyntheticBlockStartMarker(syntheticBlockStartMarker *string) {
	o.syntheticBlockStartMarker = syntheticBlockStartMarker
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
	o.stripTypes = cloneOrEmptyStringSet(stripTypes)
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
	o.stripNameSuffixes = cloneOrEmptyStringSet(stripNameSuffixes)
}

// Java: CompilerOptions.getStripNameSuffixes()
func (o *CompilerOptions) GetStripNameSuffixes() *jutil.LinkedHashSet[string] {
	return o.stripNameSuffixes
}

// @deprecated StripCode is deprecated. Code should be designed to be removed by other means.
//
// Java: CompilerOptions.setStripNamePrefixes(Set<String>)
func (o *CompilerOptions) SetStripNamePrefixes(stripNamePrefixes *jutil.LinkedHashSet[string]) {
	o.stripNamePrefixes = cloneOrEmptyStringSet(stripNamePrefixes)
}

// Java: CompilerOptions.getStripNamePrefixes()
func (o *CompilerOptions) GetStripNamePrefixes() *jutil.LinkedHashSet[string] {
	return o.stripNamePrefixes
}

// Java: CompilerOptions.getCustomPassesAt(CustomPassExecutionTime)
func (o *CompilerOptions) GetCustomPassesAt(executionTime CustomPassExecutionTime) []CompilerPass {
	return o.customPasses.get(executionTime)
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

// Java: CompilerOptions.getCssRenamingMap()
func (o *CompilerOptions) GetCssRenamingMap() CssRenamingMap {
	return o.cssRenamingMap
}

// Java: CompilerOptions.setCssRenamingSkiplist(Set<String>)
func (o *CompilerOptions) SetCssRenamingSkiplist(skiplist *jutil.LinkedHashSet[string]) {
	o.cssRenamingSkiplist = cloneOrEmptyStringSet(skiplist)
}

// Java: CompilerOptions.getCssRenamingSkiplist()
func (o *CompilerOptions) GetCssRenamingSkiplist() *jutil.LinkedHashSet[string] {
	return o.cssRenamingSkiplist
}

// Java: CompilerOptions.getReplaceStringsFunctionDescriptions()
func (o *CompilerOptions) GetReplaceStringsFunctionDescriptions() []string {
	return o.replaceStringsFunctionDescriptions
}

// Java: CompilerOptions.getReplaceStringsPlaceholderToken()
func (o *CompilerOptions) GetReplaceStringsPlaceholderToken() string {
	return o.replaceStringsPlaceholderToken
}
