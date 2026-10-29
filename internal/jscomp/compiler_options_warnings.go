// Port of com/google/javascript/jscomp/CompilerOptions.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/WarningLevel.java @ closure-compiler v20261021

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

// Java: CompilerOptions.shouldInferConsts()
func (o *CompilerOptions) ShouldInferConsts() bool {
	return o.inferConsts
}

// Gets the inferTypes flag. Note that if checkTypes is enabled, this flag is ignored when
// configuring the compiler.
//
// Java: CompilerOptions.getInferTypes()
func (o *CompilerOptions) GetInferTypes() bool {
	return o.inferTypes
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

// WarningLevel converts the warnings level to an Options object.
type WarningLevel int32

const (
	WarningLevel_QUIET WarningLevel = iota

	WarningLevel_DEFAULT

	WarningLevel_VERBOSE
)

var warningLevel_names = [...]string{"QUIET", "DEFAULT", "VERBOSE"}

// Java: WarningLevel.values()
func WarningLevel_values() []WarningLevel {
	return []WarningLevel{WarningLevel_QUIET, WarningLevel_DEFAULT, WarningLevel_VERBOSE}
}

// String is Enum.name().
func (l WarningLevel) String() string { return warningLevel_names[l] }

// Ordinal is Enum.ordinal().
func (l WarningLevel) Ordinal() int { return int(l) }

// Java: WarningLevel.setOptionsForWarningLevel(CompilerOptions)
func (l WarningLevel) SetOptionsForWarningLevel(options *CompilerOptions) {
	switch l {
	case WarningLevel_QUIET:
		WarningLevel_SilenceAllWarnings(options)
	case WarningLevel_DEFAULT:
		WarningLevel_AddDefaultWarnings(options)
	case WarningLevel_VERBOSE:
		WarningLevel_AddVerboseWarnings(options)
	}
}

// WarningLevel_SilenceAllWarnings silences all non-essential warnings.
//
// Java: WarningLevel.silenceAllWarnings(CompilerOptions)
func WarningLevel_SilenceAllWarnings(options *CompilerOptions) {
	// Just use a ShowByPath warnings guard, so that we don't have
	// to maintain a separate class of warnings guards for silencing warnings.
	options.AddWarningsGuard(
		NewShowByPathWarningsGuard(
			"the_longest_path_that_cannot_be_expressed_as_a_string"))

	// Allow passes that aren't going to report anything to be skipped.

	options.SetWarningLevel(DiagnosticGroups_MISSING_PROVIDE, CheckLevel_OFF)
	options.SetCheckTypes(false)
	options.SetWarningLevel(DiagnosticGroups_CHECK_TYPES, CheckLevel_OFF)
	options.SetWarningLevel(DiagnosticGroups_CHECK_USELESS_CODE, CheckLevel_OFF)
	options.SetWarningLevel(DiagnosticGroups_MISSING_RETURN, CheckLevel_OFF)
	options.SetWarningLevel(DiagnosticGroups_ACCESS_CONTROLS, CheckLevel_OFF)
	options.SetWarningLevel(DiagnosticGroups_CONST, CheckLevel_OFF)
	options.SetWarningLevel(DiagnosticGroups_CONSTANT_PROPERTY, CheckLevel_OFF)
	options.SetCheckSuspiciousCode(false)
	options.SetWarningLevel(DiagnosticGroups_GLOBAL_THIS, CheckLevel_OFF)
	options.SetWarningLevel(DiagnosticGroups_GLOBAL_THIS, CheckLevel_OFF)
	options.SetWarningLevel(DiagnosticGroups_ES5_STRICT, CheckLevel_OFF)

	// Allows annotations that are not standard.
	options.SetWarningLevel(DiagnosticGroups_NON_STANDARD_JSDOC,
		CheckLevel_OFF)
}

// WarningLevel_AddDefaultWarnings adds the default checking pass to the
// compilation options.
//
// Java: WarningLevel.addDefaultWarnings(CompilerOptions)
func WarningLevel_AddDefaultWarnings(options *CompilerOptions) {
	options.SetCheckSuspiciousCode(true)

	// Allows annotations that are not standard.
	options.SetWarningLevel(DiagnosticGroups_NON_STANDARD_JSDOC,
		CheckLevel_OFF)
}

// WarningLevel_AddVerboseWarnings adds all the check pass that are possibly
// relevant to a non-googler.
//
// Java: WarningLevel.addVerboseWarnings(CompilerOptions)
func WarningLevel_AddVerboseWarnings(options *CompilerOptions) {
	WarningLevel_AddDefaultWarnings(options)

	// checkSuspiciousCode needs to be enabled for CheckGlobalThis to get run.
	options.SetCheckSuspiciousCode(true)
	options.SetWarningLevel(DiagnosticGroups_GLOBAL_THIS, CheckLevel_WARNING)
	options.SetCheckSymbols(true)

	// checkTypes has the side-effect of asserting that the
	// correct number of arguments are passed to a function.
	// Because the CodingConvention used with the web service does not provide a
	// way for optional arguments to be specified, these warnings may result in
	// false positives.
	options.SetCheckTypes(true)
	options.SetWarningLevel(DiagnosticGroups_MISSING_PROPERTIES, CheckLevel_WARNING)
	options.SetWarningLevel(
		DiagnosticGroups_DEPRECATED, CheckLevel_WARNING)
	options.SetWarningLevel(
		DiagnosticGroups_ES5_STRICT, CheckLevel_WARNING)
	options.SetWarningLevel(
		DiagnosticGroups_VISIBILITY, CheckLevel_WARNING)
	options.SetWarningLevel(
		DiagnosticGroups_CONST, CheckLevel_WARNING)
	options.SetWarningLevel(
		DiagnosticGroups_CHECK_REGEXP, CheckLevel_WARNING)
	options.SetWarningLevel(
		DiagnosticGroups_STRICT_MODULE_DEP_CHECK, CheckLevel_WARNING)
	options.SetWarningLevel(
		DiagnosticGroups_MISSING_RETURN, CheckLevel_WARNING)

	// Kindly tell the user that they have JsDocs that we don't understand.
	options.SetWarningLevel(DiagnosticGroups_NON_STANDARD_JSDOC,
		CheckLevel_WARNING)
}
