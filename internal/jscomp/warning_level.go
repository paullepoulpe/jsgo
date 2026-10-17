// Port of com/google/javascript/jscomp/WarningLevel.java @ closure-compiler v20260415

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
