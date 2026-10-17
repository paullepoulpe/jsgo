// Port of com/google/javascript/jscomp/RhinoErrorReporter.java @ closure-compiler v20260415

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
	"regexp"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// RhinoErrorReporterBase is an error reporter for serializing Rhino errors
// into our error format.
//
// Go-only: language-adapter — RhinoErrorReporter has a single, private
// subclass and no overridable methods, so the class is this embeddable
// struct (PORTING.md §4 without the interface).
type RhinoErrorReporterBase struct {
	internalReporter ErrorHandler
}

var (
	RhinoErrorReporter_PARSE_ERROR = DiagnosticType_Error("JSC_PARSE_ERROR", "Parse error. {0}")

	RhinoErrorReporter_TYPE_PARSE_ERROR = DiagnosticType_Warning("JSC_TYPE_PARSE_ERROR", "{0}")

	RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR = DiagnosticType_Warning("JSC_UNRECOGNIZED_TYPE_ERROR", "{0}")

	RhinoErrorReporter_UNRECOGNIZED_TYPEOF_ERROR = DiagnosticType_Warning("JSC_UNRECOGNIZED_TYPEOF_ERROR", "{0}")

	RhinoErrorReporter_CYCLIC_INHERITANCE_ERROR = DiagnosticType_Warning("JSC_CYCLIC_INHERITANCE_ERROR", "{0}")

	// This is separate from TYPE_PARSE_ERROR because there are many instances of this warning
	// and it is unfeasible to fix them all right away.
	RhinoErrorReporter_JSDOC_MISSING_BRACES_WARNING = DiagnosticType_Disabled("JSC_JSDOC_MISSING_BRACES_WARNING", "{0}")

	// This is separate from TYPE_PARSE_ERROR because there are many instances of this warning
	// and it is unfeasible to fix them all right away.
	RhinoErrorReporter_JSDOC_MISSING_TYPE_WARNING = DiagnosticType_Disabled("JSC_JSDOC_MISSING_TYPE_WARNING", "{0}")

	// Import is supported by VSCode and is pretty much a standard now.
	RhinoErrorReporter_JSDOC_IMPORT_TYPE_WARNING = DiagnosticType_Disabled("JSC_JSDOC_IMPORT_TYPE_WARNING", "{0}")

	RhinoErrorReporter_TOO_MANY_TEMPLATE_PARAMS = DiagnosticType_Warning("JSC_TOO_MANY_TEMPLATE_PARAMS", "{0}")

	// Special-cased errors, so that they can be configured via the
	// warnings API.
	RhinoErrorReporter_TRAILING_COMMA = DiagnosticType_Error(
		"JSC_TRAILING_COMMA",
		"Parse error. IE8 (and below) will parse trailing commas in "+
			"array and object literals incorrectly. "+
			"If you are targeting newer versions of JS, "+
			"set the appropriate language_in option.")

	RhinoErrorReporter_DUPLICATE_PARAM = DiagnosticType_Error("JSC_DUPLICATE_PARAM", "Parse error. {0}")

	RhinoErrorReporter_DUPLICATE_VISIBILITY = DiagnosticType_Warning("JSC_DUPLICATE_VISIBILITY", "{0}")

	RhinoErrorReporter_UNNECESSARY_ESCAPE = DiagnosticType_Disabled("JSC_UNNECESSARY_ESCAPE", "Parse error. {0}")

	RhinoErrorReporter_INVALID_PARAM = DiagnosticType_Warning("JSC_INVALID_PARAM", "Parse error. {0}")

	RhinoErrorReporter_BAD_JSDOC_ANNOTATION = DiagnosticType_Warning("JSC_BAD_JSDOC_ANNOTATION", "Parse error. {0}")

	RhinoErrorReporter_INVALID_ES3_PROP_NAME = DiagnosticType_Warning(
		"JSC_INVALID_ES3_PROP_NAME",
		"Keywords and reserved words are not allowed as unquoted property "+
			"names in older versions of JavaScript. "+
			"If you are targeting newer versions of JavaScript, "+
			"set the appropriate language_in option.")

	RhinoErrorReporter_PARSE_TREE_TOO_DEEP = DiagnosticType_Error("JSC_PARSE_TREE_TOO_DEEP", "Parse tree too deep.")

	RhinoErrorReporter_INVALID_OCTAL_LITERAL = DiagnosticType_Warning(
		"JSC_INVALID_OCTAL_LITERAL",
		"This style of octal literal is not supported in strict mode.")

	// This check should only be enabled in the lintChecks DiagnosticGroup as part of the linter
	RhinoErrorReporter_STRING_CONTINUATION = DiagnosticType_Disabled("JSC_STRING_CONTINUATION", "{0}")

	RhinoErrorReporter_LANGUAGE_FEATURE = DiagnosticType_Error("JSC_LANGUAGE_FEATURE", "{0}.")

	RhinoErrorReporter_UNSUPPORTED_LANGUAGE_FEATURE = DiagnosticType_Error("JSC_UNSUPPORTED_LANGUAGE_FEATURE", "{0}.")

	RhinoErrorReporter_UNSUPPORTED_BOUNDED_GENERIC_TYPES = DiagnosticType_Error(
		"JSC_UNSUPPORTED_BOUNDED_GENERIC_TYPES",
		"Bounded generic semantics are currently still in development")

	RhinoErrorReporter_BOUNDED_GENERIC_TYPE_ERROR = DiagnosticType_Error(
		"JSC_BOUNDED_GENERIC_TYPE_ERROR",
		"Bounded generic type error. "+
			"{0} assigned to template type {1} is not a subtype of bound {2}")

	RhinoErrorReporter_CLOSURE_UNAWARE_ANNOTATION_PRESENT = DiagnosticType_Disabled(
		"JSC_CLOSURE_UNAWARE_ANNOTATION_PRESENT",
		rhino.Msg_JSDOC_CLOSURE_UNAWARE_CODE_INVALID.Format())
)

// rhinoErrorReporter_typeMapEntry is an entry of the ImmutableMap typeMap.
type rhinoErrorReporter_typeMapEntry struct {
	pattern *jutil.Pattern
	type_   *DiagnosticType
}

// A map of Rhino messages to their DiagnosticType.
var rhinoErrorReporter_typeMap = []rhinoErrorReporter_typeMapEntry{
	// Trailing comma
	{
		jutil.Pattern_Compile("Trailing comma is not legal in an ECMA-262 object initializer"),
		RhinoErrorReporter_TRAILING_COMMA},
	// Duplicate parameter
	{rhinoErrorReporter_replacePlaceHolders("Duplicate parameter name \"{0}\""), RhinoErrorReporter_DUPLICATE_PARAM},
	// Duplicate visiblity
	{rhinoErrorReporter_replacePlaceHolders(rhino.Msg_JSDOC_EXTRA_VISIBILITY.Format()), RhinoErrorReporter_DUPLICATE_VISIBILITY},
	{jutil.Pattern_Compile("Unnecessary escape:.*"), RhinoErrorReporter_UNNECESSARY_ESCAPE},
	{jutil.Pattern_Compile("^invalid param name.*"), RhinoErrorReporter_INVALID_PARAM},
	// Unknown @annotations.
	{rhinoErrorReporter_replacePlaceHolders(rhino.Msg_BAD_JSDOC_TAG.Format()), RhinoErrorReporter_BAD_JSDOC_ANNOTATION},
	{
		jutil.Pattern_Compile(
			"^Keywords and reserved words are not allowed as unquoted property.*"),
		RhinoErrorReporter_INVALID_ES3_PROP_NAME},
	{jutil.Pattern_Compile("^Too many template parameters\n.*"), RhinoErrorReporter_TOO_MANY_TEMPLATE_PARAMS},
	// Type annotation warnings.
	{
		jutil.Pattern_Compile(".*Type annotations should have curly braces.*"),
		RhinoErrorReporter_JSDOC_MISSING_BRACES_WARNING},
	{jutil.Pattern_Compile("Missing type declaration\\."), RhinoErrorReporter_JSDOC_MISSING_TYPE_WARNING},
	// Unresolved types that aren't forward declared.
	{jutil.Pattern_Compile(".*Unknown type.*"), RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR},
	{jutil.Pattern_Compile(".*Unknown type.*\n.*"), RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR},
	// Unrecognized `typeof some.prop` errors
	{jutil.Pattern_Compile("^Missing type for `typeof` value.*"), RhinoErrorReporter_UNRECOGNIZED_TYPEOF_ERROR},
	// Cyclic inheritance errors
	{
		jutil.Pattern_Compile("^Cycle detected in inheritance chain of type .*"),
		RhinoErrorReporter_CYCLIC_INHERITANCE_ERROR},
	// Import annotation errors.
	{
		jutil.Pattern_Compile("^Bad type annotation. Import in typedef.*"),
		RhinoErrorReporter_JSDOC_IMPORT_TYPE_WARNING},
	// Type annotation errors.
	{jutil.Pattern_Compile("^Bad type annotation.*"), RhinoErrorReporter_TYPE_PARSE_ERROR},
	{jutil.Pattern_Compile("constructed type must be an object type"), RhinoErrorReporter_TYPE_PARSE_ERROR},
	// Parse tree too deep.
	{jutil.Pattern_Compile("Too deep recursion while parsing"), RhinoErrorReporter_PARSE_TREE_TOO_DEEP},
	// Old-style octal literals
	{jutil.Pattern_Compile("^Octal .*literal.*"), RhinoErrorReporter_INVALID_OCTAL_LITERAL},
	{jutil.Pattern_Compile("^String continuations.*"), RhinoErrorReporter_STRING_CONTINUATION},
	{jutil.Pattern_Compile("^This language feature is only supported for .*"), RhinoErrorReporter_LANGUAGE_FEATURE},
	{
		jutil.Pattern_Compile(
			"^This language feature is not currently supported by the compiler:" + " .*"),
		RhinoErrorReporter_UNSUPPORTED_LANGUAGE_FEATURE},
	{
		jutil.Pattern_Compile("Bounded generic semantics are currently still in development"),
		RhinoErrorReporter_UNSUPPORTED_BOUNDED_GENERIC_TYPES},
	{jutil.Pattern_Compile("^Bounded generic type error.*"), RhinoErrorReporter_BOUNDED_GENERIC_TYPE_ERROR},
	{
		rhinoErrorReporter_replacePlaceHolders(rhino.Msg_JSDOC_CLOSURE_UNAWARE_CODE_INVALID.Format()),
		RhinoErrorReporter_CLOSURE_UNAWARE_ANNOTATION_PRESENT},
}

var rhinoErrorReporter_placeHolder = regexp.MustCompile(`\{\d+\}`)

// rhinoErrorReporter_replacePlaceHolders replaces, for each message such as
// "Not a good use of {0}", the place holder {0} with a wild card that
// matches all possible strings. Also put the any non-place-holder in quotes
// for regex matching later.
//
// Java: RhinoErrorReporter.replacePlaceHolders(String)
func rhinoErrorReporter_replacePlaceHolders(s string) *jutil.Pattern {
	s = jutil.Pattern_Quote(s)
	// String.replaceAll("\\{\\d+\\}", "\\\\E.*\\\\Q"): the replacement is the literal \E.*\Q.
	return jutil.Pattern_Compile(rhinoErrorReporter_placeHolder.ReplaceAllLiteralString(s, `\E.*\Q`))
}

// Java: RhinoErrorReporter(ErrorHandler)
func (r *RhinoErrorReporterBase) InitRhinoErrorReporterBase(internalReporter ErrorHandler) {
	r.internalReporter = internalReporter
}

// Java: RhinoErrorReporter.forOldRhino(ErrorHandler)
func RhinoErrorReporter_ForOldRhino(internalReporter ErrorHandler) rhino.ErrorReporter {
	return newRhinoErrorReporter_OldRhinoErrorReporter(internalReporter)
}

// Java: RhinoErrorReporter.warningAtLine(String, String, int, int)
func (r *RhinoErrorReporterBase) WarningAtLine(message string, sourceName string, line int,
	lineOffset int) {
	r.internalReporter.Report(
		nil, rhinoErrorReporter_makeError(message, sourceName, line, lineOffset, CheckLevel_WARNING))
}

// Java: RhinoErrorReporter.errorAtLine(String, String, int, int)
func (r *RhinoErrorReporterBase) ErrorAtLine(message string, sourceName string, line int,
	lineOffset int) {
	r.internalReporter.Report(
		nil, rhinoErrorReporter_makeError(message, sourceName, line, lineOffset, CheckLevel_ERROR))
}

// RhinoErrorReporter_MapError may return nil.
//
// Java: RhinoErrorReporter.mapError(String)
func RhinoErrorReporter_MapError(message string) *DiagnosticType {
	for _, entry := range rhinoErrorReporter_typeMap {
		if entry.pattern.Matches(message) {
			return entry.type_
		}
	}
	return nil
}

// Java: RhinoErrorReporter.makeError(String, String, int, int, CheckLevel)
func rhinoErrorReporter_makeError(
	message string, sourceName string, line int, lineOffset int, defaultLevel CheckLevel) *JSError {
	// Try to see if the message is one of the rhino errors we want to
	// expose as DiagnosticType by matching it with the regex key.
	type_ := RhinoErrorReporter_MapError(message)
	t := type_
	if t == nil {
		t = RhinoErrorReporter_PARSE_ERROR
	}
	builder :=
		JSError_Builder_(t, message).
			SetSourceLocation(&sourceName, line, lineOffset)

	// If a specific DiagnosticType is present, its default level overrides the defaultLevel
	// passed to this method.
	if type_ == nil {
		builder.SetLevel(defaultLevel)
	}

	return builder.Build()
}

// rhinoErrorReporter_OldRhinoErrorReporter is RhinoErrorReporter's private
// subclass implementing rhino.ErrorReporter.
//
// Java: RhinoErrorReporter.OldRhinoErrorReporter
type rhinoErrorReporter_OldRhinoErrorReporter struct {
	RhinoErrorReporterBase
}

var _ rhino.ErrorReporter = (*rhinoErrorReporter_OldRhinoErrorReporter)(nil)

// Java: RhinoErrorReporter.OldRhinoErrorReporter(ErrorHandler)
func newRhinoErrorReporter_OldRhinoErrorReporter(internalReporter ErrorHandler) *rhinoErrorReporter_OldRhinoErrorReporter {
	r := &rhinoErrorReporter_OldRhinoErrorReporter{}
	r.InitRhinoErrorReporterBase(internalReporter)
	return r
}

// Java: RhinoErrorReporter.OldRhinoErrorReporter.error(String, String, int, int)
func (r *rhinoErrorReporter_OldRhinoErrorReporter) Error(message string, sourceName string, line int,
	lineOffset int) {
	r.RhinoErrorReporterBase.ErrorAtLine(message, sourceName, line, lineOffset)
}

// Java: RhinoErrorReporter.OldRhinoErrorReporter.warning(String, String, int, int)
func (r *rhinoErrorReporter_OldRhinoErrorReporter) Warning(message string, sourceName string, line int,
	lineOffset int) {
	r.RhinoErrorReporterBase.WarningAtLine(message, sourceName, line, lineOffset)
}
