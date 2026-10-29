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
	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Deprecated. Please use setWarningLevel(DiagnosticGroups.GLOBAL_THIS, level) instead.
//
// Java: CompilerOptions.setCheckGlobalThisLevel(CheckLevel)
func (o *CompilerOptions) SetCheckGlobalThisLevel(level CheckLevel) {
}

// Whether the warnings guard in this Options object enables the given group of warnings.
//
// Java: CompilerOptions.enables(DiagnosticGroup)
func (o *CompilerOptions) Enables(group *DiagnosticGroup) bool {
	return o.warningsGuard.MustRunChecks(group) == base.Tri_TRUE
}

// Whether the warnings guard in this Options object disables the given group of warnings.
//
// Java: CompilerOptions.disables(DiagnosticGroup)
func (o *CompilerOptions) Disables(group *DiagnosticGroup) bool {
	return o.warningsGuard.MustRunChecks(group) == base.Tri_FALSE
}

// Configure the given type of warning to the given level.
//
// Java: CompilerOptions.setWarningLevel(DiagnosticGroup, CheckLevel)
func (o *CompilerOptions) SetWarningLevel(type_ *DiagnosticGroup, level CheckLevel) {
	o.AddWarningsGuard(NewDiagnosticGroupWarningsGuard(type_, level))
}

// Java: CompilerOptions.getWarningsGuard()
func (o *CompilerOptions) GetWarningsGuard() WarningsGuard {
	return o.warningsGuard
}

// Reset the warnings guard.
//
// Java: CompilerOptions.resetWarningsGuard()
func (o *CompilerOptions) ResetWarningsGuard() {
	o.warningsGuard = NewComposeWarningsGuard()
}

// Add a guard to the set of warnings guards.
//
// Java: CompilerOptions.addWarningsGuard(WarningsGuard)
func (o *CompilerOptions) AddWarningsGuard(guard WarningsGuard) {
	o.warningsGuard.AddGuard(guard)
}

// Set a custom handler for warnings and errors.
//
// This is mostly used for piping the warnings and errors to a file behind the scenes.
//
// If you want to filter warnings and errors, you should use a WarningsGuard.
//
// If you want to change how warnings and errors are reported to the user, you should set a
// ErrorManager on the Compiler. An ErrorManager is intended to summarize the errors for a single
// compile job.
//
// Java: CompilerOptions.setErrorHandler(ErrorHandler)
func (o *CompilerOptions) SetErrorHandler(handler ErrorHandler) {
	o.errorHandler = handler
}

// Java: CompilerOptions.getErrorHandler()
func (o *CompilerOptions) GetErrorHandler() ErrorHandler {
	return o.errorHandler
}

// Java: CompilerOptions.setColorizeErrorOutput(boolean)
func (o *CompilerOptions) SetColorizeErrorOutput(colorizeErrorOutput bool) {
	o.colorizeErrorOutput = colorizeErrorOutput
}

// Java: CompilerOptions.shouldColorizeErrorOutput()
func (o *CompilerOptions) ShouldColorizeErrorOutput() bool {
	return o.colorizeErrorOutput
}

// Java: CompilerOptions.getExtraReportGenerators()
func (o *CompilerOptions) GetExtraReportGenerators() []SortingErrorManager_ErrorReportGenerator {
	return o.extraReportGenerators
}

// Java: CompilerOptions.addReportGenerator(SortingErrorManager.ErrorReportGenerator)
func (o *CompilerOptions) AddReportGenerator(generator SortingErrorManager_ErrorReportGenerator) {
	o.extraReportGenerators = append(o.extraReportGenerators, generator)
}

// Controls how detailed the compilation summary is. Values: 0 (never print summary), 1 (print
// summary only if there are errors or warnings), 2 (print summary if type checking is on), 3
// (always print summary). The default level is 1
//
// Java: CompilerOptions.setSummaryDetailLevel(int)
func (o *CompilerOptions) SetSummaryDetailLevel(summaryDetailLevel int) {
	o.summaryDetailLevel = summaryDetailLevel
}

// Java: CompilerOptions.getSummaryDetailLevel()
func (o *CompilerOptions) GetSummaryDetailLevel() int {
	return o.summaryDetailLevel
}

// Java: CompilerOptions.setExtraAnnotationNames(Iterable<String>)
func (o *CompilerOptions) SetExtraAnnotationNames(extraAnnotationNames []string) {
	o.extraAnnotationNames = jutil.NewLinkedHashSet(extraAnnotationNames...)
}

// Java: CompilerOptions.getExtraAnnotationNames()
func (o *CompilerOptions) GetExtraAnnotationNames() *jutil.LinkedHashSet[string] {
	return o.extraAnnotationNames
}

// Java: CompilerOptions.setChecksOnly(boolean)
func (o *CompilerOptions) SetChecksOnly(checksOnly bool) {
	o.checksOnly = checksOnly
}

// Java: CompilerOptions.isChecksOnly()
func (o *CompilerOptions) IsChecksOnly() bool {
	return o.checksOnly
}

// Java: CompilerOptions.setIncrementalChecks(IncrementalCheckMode)
func (o *CompilerOptions) SetIncrementalChecks(value CompilerOptions_IncrementalCheckMode) {
	o.incrementalCheckMode = value
	switch value {
	case CompilerOptions_IncrementalCheckMode_OFF:
	case CompilerOptions_IncrementalCheckMode_GENERATE_IJS:
		o.SetPreserveTypeAnnotations(true)
		o.SetOutputJs(CompilerOptions_OutputJs_NORMAL)
	}
}

// Java: CompilerOptions.shouldGenerateTypedExterns()
func (o *CompilerOptions) ShouldGenerateTypedExterns() bool {
	return o.incrementalCheckMode == CompilerOptions_IncrementalCheckMode_GENERATE_IJS
}

// TODO(tbreisacher): Remove this method after ctemplate issues are solved.
//
// Java: CompilerOptions.setInferConst(boolean)
func (o *CompilerOptions) SetInferConst(value bool) {
	o.inferConsts = value
}

// Java: CompilerOptions.shouldInferConsts()
func (o *CompilerOptions) ShouldInferConsts() bool {
	return o.inferConsts
}

// If true, enables type inference. If checkTypes is enabled, this flag has no effect.
//
// Java: CompilerOptions.setInferTypes(boolean)
func (o *CompilerOptions) SetInferTypes(enable bool) {
	o.inferTypes = enable
}

// Gets the inferTypes flag. Note that if checkTypes is enabled, this flag is ignored when
// configuring the compiler.
//
// Java: CompilerOptions.getInferTypes()
func (o *CompilerOptions) GetInferTypes() bool {
	return o.inferTypes
}

// @deprecated This is a no-op.
//
// Java: CompilerOptions.setNewTypeInference(boolean)
func (o *CompilerOptions) SetNewTypeInference(enable bool) {
}

// @return true if either typechecker is ON.
//
// Java: CompilerOptions.isTypecheckingEnabled()
func (o *CompilerOptions) IsTypecheckingEnabled() bool {
	return o.checkTypes
}

// Java: CompilerOptions.getConformanceReportingMode()
func (o *CompilerOptions) GetConformanceReportingMode() CompilerOptions_ConformanceReportingMode {
	return o.conformanceReportingMode
}

// Java: CompilerOptions.setConformanceReportingMode(ConformanceReportingMode)
func (o *CompilerOptions) SetConformanceReportingMode(mode CompilerOptions_ConformanceReportingMode) {
	o.conformanceReportingMode = mode
}

// Java: CompilerOptions.setConformanceRemoveRegexFromPath(Optional<Pattern>)
func (o *CompilerOptions) SetConformanceRemoveRegexFromPath(pattern *jutil.Pattern) {
	o.conformanceRemoveRegexFromPath = pattern
}

// Java: CompilerOptions.getConformanceRemoveRegexFromPath()
func (o *CompilerOptions) GetConformanceRemoveRegexFromPath() *jutil.Pattern {
	return o.conformanceRemoveRegexFromPath
}

// @return Whether assumeStrictThis is set.
//
// Java: CompilerOptions.assumeStrictThis()
func (o *CompilerOptions) AssumeStrictThis() bool {
	return o.assumeStrictThis
}

// If true, enables enables additional optimizations.
//
// Java: CompilerOptions.setAssumeStrictThis(boolean)
func (o *CompilerOptions) SetAssumeStrictThis(enable bool) {
	o.assumeStrictThis = enable
}

// If true, considers all missing types to be forward declared (useful for partial
// compilation).
//
// Java: CompilerOptions.setAssumeForwardDeclaredForMissingTypes(boolean)
func (o *CompilerOptions) SetAssumeForwardDeclaredForMissingTypes(assumeForwardDeclaredForMissingTypes bool) {
	o.assumeForwardDeclaredForMissingTypes = assumeForwardDeclaredForMissingTypes
}

// Java: CompilerOptions.assumeForwardDeclaredForMissingTypes()
func (o *CompilerOptions) AssumeForwardDeclaredForMissingTypes() bool {
	return o.assumeForwardDeclaredForMissingTypes
}

// When enabled, assume that apparently side-effect free code is meaningful.
//
// Java: CompilerOptions.setProtectHiddenSideEffects(boolean)
func (o *CompilerOptions) SetProtectHiddenSideEffects(enable bool) {
	o.protectHiddenSideEffects = enable
}

// Whether or not the compiler should wrap apparently side-effect free code to prevent it from
// being removed
//
// Java: CompilerOptions.shouldProtectHiddenSideEffects()
func (o *CompilerOptions) ShouldProtectHiddenSideEffects() bool {
	// TODO(lharker): would it make sense to make 'checksOnly' a tristate instead? either we're
	// running only checks for refactoring; running checks to generate a .typedast; or running
	// checks as part of a fully optimizing build?
	return o.protectHiddenSideEffects && (!o.checksOnly || o.GetTypedAstOutputFile() != nil)
}

// Java: CompilerOptions.setAssumeGettersArePure(boolean)
func (o *CompilerOptions) SetAssumeGettersArePure(x bool) {
	o.assumeGettersArePure = x
}

// Java: CompilerOptions.getAssumeGettersArePure()
func (o *CompilerOptions) GetAssumeGettersArePure() bool {
	return o.assumeGettersArePure
}

// Java: CompilerOptions.setAssumePropertiesAreStaticallyAnalyzable(boolean)
func (o *CompilerOptions) SetAssumePropertiesAreStaticallyAnalyzable(x bool) {
	o.assumePropertiesAreStaticallyAnalyzable = x
}

// Java: CompilerOptions.getAssumePropertiesAreStaticallyAnalyzable()
func (o *CompilerOptions) GetAssumePropertiesAreStaticallyAnalyzable() bool {
	return o.assumePropertiesAreStaticallyAnalyzable
}

// Java: CompilerOptions.setAssumeStaticInheritanceIsNotUsed(boolean)
func (o *CompilerOptions) SetAssumeStaticInheritanceIsNotUsed(x bool) {
	o.assumeStaticInheritanceIsNotUsed = x
}

// Java: CompilerOptions.getAssumeStaticInheritanceIsNotUsed()
func (o *CompilerOptions) GetAssumeStaticInheritanceIsNotUsed() bool {
	return o.assumeStaticInheritanceIsNotUsed
}

// Sets the list of properties that we report property invalidation errors for.
//
// Java: CompilerOptions.setPropertiesThatMustDisambiguate(Set<String>)
func (o *CompilerOptions) SetPropertiesThatMustDisambiguate(names *jutil.LinkedHashSet[string]) {
	o.propertiesThatMustDisambiguate = jutil.NewLinkedHashSet(names.ToSlice()...)
}

// Java: CompilerOptions.getPropertiesThatMustDisambiguate()
func (o *CompilerOptions) GetPropertiesThatMustDisambiguate() *jutil.LinkedHashSet[string] {
	return o.propertiesThatMustDisambiguate
}
