// Port of test/com/google/javascript/jscomp/RhinoErrorReporterTest.java @ closure-compiler v20261015

/*
 * Copyright 2010 The Closure Compiler Authors.
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
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
)

// Tests for error message filtering.
//
// Java: public final class RhinoErrorReporterTest
type RhinoErrorReporterTest struct {
	// Go-only: test-helper — the Go test running the current test method.
	t *testing.T

	reportLintWarnings bool
	languageIn         CompilerOptions_LanguageMode
}

func TestRhinoErrorReporterTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return &RhinoErrorReporterTest{} })
}

// JUnitSetT gives the instance the Go test running its test method.
func (c *RhinoErrorReporterTest) JUnitSetT(t *testing.T) { c.t = t }

// JUnitBefore runs the @Before method.
func (c *RhinoErrorReporterTest) JUnitBefore() { c.SetUp() }

// JUnitAfter: the class has no @After method.
func (c *RhinoErrorReporterTest) JUnitAfter() {}

// Java: RhinoErrorReporterTest.setUp()
func (c *RhinoErrorReporterTest) SetUp() {
	c.reportLintWarnings = true
	c.languageIn = CompilerOptions_LanguageMode_UNSUPPORTED
}

// Java: RhinoErrorReporterTest.testMissingTypeWarnings()
func (c *RhinoErrorReporterTest) TestMissingTypeWarnings() {
	c.reportLintWarnings = false

	c.assertNoWarningOrError("/** @return */ function f() {}")

	c.reportLintWarnings = true

	message := "Missing type declaration."
	error_ :=
		c.assertWarning(
			"/** @return */ function f() {}",
			RhinoErrorReporter_JSDOC_MISSING_TYPE_WARNING,
			message)

	AssertThat(error_.GetLineNumber()).IsEqualTo(1)
	AssertThat(error_.Charno()).IsEqualTo(4)
}

// Java: RhinoErrorReporterTest.testMissingCurlyBraceWarning()
func (c *RhinoErrorReporterTest) TestMissingCurlyBraceWarning() {
	c.reportLintWarnings = false
	c.assertNoWarningOrError("/** @type string */ var x;")

	c.reportLintWarnings = true
	c.assertWarning(
		"/** @type string */ var x;",
		RhinoErrorReporter_JSDOC_MISSING_BRACES_WARNING,
		"Bad type annotation. Type annotations should have curly braces."+parsing.JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

// Java: RhinoErrorReporterTest.testLanguageFeatureInHigherLanguageInError()
func (c *RhinoErrorReporterTest) TestLanguageFeatureInHigherLanguageInError() {
	c.languageIn = CompilerOptions_LanguageMode_ECMASCRIPT_2015
	c.assertError(
		"2 ** 3",
		RhinoErrorReporter_LANGUAGE_FEATURE,
		"This language feature is only supported for ECMASCRIPT_2016 mode or better: "+
			"exponent operator (**).")
}

// Verifies that the compiler emits an error for the given code.
//
// Java: RhinoErrorReporterTest.assertNoWarningOrError(String)
func (c *RhinoErrorReporterTest) assertNoWarningOrError(code string) {
	compiler := c.parseCode(code)
	AssertWithMessage("Expected error").That(compiler.GetErrorCount()).IsEqualTo(0)
	AssertWithMessage("Expected warning").That(compiler.GetErrorCount()).IsEqualTo(0)
}

// Verifies that the compiler emits an error for the given code.
//
// Java: RhinoErrorReporterTest.assertError(String, DiagnosticType, String)
func (c *RhinoErrorReporterTest) assertError(code string, type_ *DiagnosticType, description string) *JSError {
	compiler := c.parseCode(code)
	AssertWithMessage("Expected error").That(compiler.GetErrorCount()).IsEqualTo(1)

	error_ := rhinoErrorReporterTest_getOnlyElement(compiler.GetErrors())
	AssertThat(error_.Type()).IsEqualTo(type_)
	AssertThat(error_.Description()).IsEqualTo(description)
	return error_
}

// Verifies that the compiler emits a warning for the given code.
//
// Java: RhinoErrorReporterTest.assertWarning(String, DiagnosticType, String)
func (c *RhinoErrorReporterTest) assertWarning(code string, type_ *DiagnosticType, description string) *JSError {
	compiler := c.parseCode(code)
	AssertWithMessage("Expected warning").That(compiler.GetWarningCount()).IsEqualTo(1)

	error_ := rhinoErrorReporterTest_getOnlyElement(compiler.GetWarnings())
	AssertThat(error_.Type()).IsEqualTo(type_)
	AssertThat(error_.Description()).IsEqualTo(description)
	return error_
}

// rhinoErrorReporterTest_getOnlyElement is Iterables.getOnlyElement.
func rhinoErrorReporterTest_getOnlyElement(errors []*JSError) *JSError {
	if len(errors) != 1 {
		panic(jutil.NewIllegalArgumentExceptionMsg(
			"expected one element but was: " + jutil.ToJavaString(errors)))
	}
	return errors[0]
}

// Java: RhinoErrorReporterTest.parseCode(String)
func (c *RhinoErrorReporterTest) parseCode(code string) *Compiler {
	compiler := NewCompiler()
	options := NewCompilerOptions()

	if !c.reportLintWarnings {
		options.SetWarningLevel(DiagnosticGroups_LINT_CHECKS, CheckLevel_OFF)
	} else {
		options.SetWarningLevel(DiagnosticGroups_LINT_CHECKS, CheckLevel_WARNING)
		options.SetWarningLevel(DiagnosticGroups_JSDOC_MISSING_TYPE, CheckLevel_WARNING)
	}

	options.SetLanguageIn(c.languageIn)

	externs := []*SourceFile{}
	inputs := []*SourceFile{SourceFile_FromCode("input", code)}
	compiler.Init(externs, inputs, options)
	compiler.ParseInputs()
	return compiler
}

// jsaRecordingErrorHandler records the reports it receives.
type jsaRecordingErrorHandler struct {
	levels []*CheckLevel
	errors []*JSError
}

func (h *jsaRecordingErrorHandler) Report(level *CheckLevel, error *JSError) {
	h.levels = append(h.levels, level)
	h.errors = append(h.errors, error)
}

// Go-only: extra coverage — the reporter maps Rhino messages to diagnostic
// types and keeps the default level only for unmapped messages.
func TestRhinoErrorReporter_ForOldRhino(t *testing.T) {
	h := &jsaRecordingErrorHandler{}
	r := RhinoErrorReporter_ForOldRhino(h)
	r.Warning("Missing type declaration.", "a.js", 1, 4)
	r.Error("something odd", "b.js", 2, 3)
	r.Warning("Duplicate parameter name \"x\"", "c.js", 3, 0)
	jsaEq(t, len(h.errors), 3)
	jsaTrue(t, h.levels[0] == nil && h.levels[1] == nil)
	jsaTrue(t, h.errors[0].Type() == RhinoErrorReporter_JSDOC_MISSING_TYPE_WARNING)
	jsaEq(t, h.errors[0].GetDefaultLevel(), CheckLevel_OFF)
	jsaEq(t, *h.errors[0].SourceName(), "a.js")
	jsaEq(t, h.errors[0].Charno(), 4)
	jsaTrue(t, h.errors[1].Type() == RhinoErrorReporter_PARSE_ERROR)
	jsaEq(t, h.errors[1].Description(), "Parse error. something odd")
	jsaEq(t, h.errors[1].GetDefaultLevel(), CheckLevel_ERROR)
	jsaTrue(t, h.errors[2].Type() == RhinoErrorReporter_DUPLICATE_PARAM)
	jsaEq(t, h.errors[2].GetDefaultLevel(), CheckLevel_ERROR)
	jsaTrue(t, RhinoErrorReporter_MapError("no match") == nil)
	jsaTrue(t, RhinoErrorReporter_MapError("Too many template parameters\nx") == RhinoErrorReporter_TOO_MANY_TEMPLATE_PARAMS)
	jsaTrue(t, RhinoErrorReporter_MapError("a Unknown type b\nc") == RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR)
}
