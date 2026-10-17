// Port of com/google/javascript/refactoring/ErrorToFixMapperTest.java @ closure-compiler v20261015
// (test/com/google/javascript/refactoring/ErrorToFixMapperTest.java)

/*
 * Copyright 2014 The Closure Compiler Authors.
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
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
)

// TODO(tjgq): Re-enable the disabled tests once
// CheckRequiresAndProvidesSorted is updated to use RequiresFixer. Currently,
// CheckRequiresAndProvidesSorted and ErrorToFixMapper sometimes disagree on
// whether a fix is required, which causes no fix to be suggested when it
// should.

// ErrorToFixMapperTest is the test case for ErrorToFixMapper.
//
// Java: public class ErrorToFixMapperTest
type ErrorToFixMapperTest struct {
	t *testing.T

	errorManager    *FixingErrorManager
	options         *CompilerOptions
	compiler        *Compiler
	preexistingCode string
}

func (c *ErrorToFixMapperTest) JUnitSetT(t *testing.T) { c.t = t }
func (c *ErrorToFixMapperTest) JUnitBefore()           { c.SetUp() }
func (c *ErrorToFixMapperTest) JUnitAfter()            {}

// Java: ErrorToFixMapperTest()
func NewErrorToFixMapperTest() *ErrorToFixMapperTest {
	return &ErrorToFixMapperTest{}
}

func TestErrorToFixMapperTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewErrorToFixMapperTest() })
}

// Java: ErrorToFixMapperTest.setUp() (@Before)
func (c *ErrorToFixMapperTest) SetUp() {
	c.errorManager = NewFixingErrorManager()
	c.compiler = NewCompilerWithErrorManager(c.errorManager)
	c.preexistingCode = ""
	c.compiler.DisableThreads()
	c.errorManager.SetCompiler(c.compiler)

	c.options = RefactoringDriver_GetCompilerOptions()
	c.options.SetWarningLevel(DiagnosticGroups_ANALYZER_CHECKS, CheckLevel_WARNING)
	c.options.SetWarningLevel(DiagnosticGroups_CHECK_VARIABLES, CheckLevel_ERROR)
	c.options.SetWarningLevel(DiagnosticGroups_DEBUGGER_STATEMENT_PRESENT, CheckLevel_ERROR)
	c.options.SetWarningLevel(DiagnosticGroups_LINT_CHECKS, CheckLevel_WARNING)
	c.options.SetWarningLevel(DiagnosticGroups_MISSING_REQUIRE, CheckLevel_ERROR)
	c.options.SetWarningLevel(DiagnosticGroups_EXTRA_REQUIRE, CheckLevel_ERROR)
	c.options.SetWarningLevel(DiagnosticGroups_STRICT_MODULE_CHECKS, CheckLevel_WARNING)
	c.options.SetCodingConvention(NewGoogleCodingConvention())
}

// ErrorToFixMapperTest_ExpectedFix is the expected outcome of one fix.
//
// Java: @AutoValue abstract static class ErrorToFixMapperTest.ExpectedFix
type ErrorToFixMapperTest_ExpectedFix struct {
	description *string
	fixedCode   string
}

// Description is the optional string describing the fix.
//
// Java: ErrorToFixMapperTest.ExpectedFix.description()
func (f *ErrorToFixMapperTest_ExpectedFix) Description() *string { return f.description }

// FixedCode is what the code should look like after applying the fix.
//
// Java: ErrorToFixMapperTest.ExpectedFix.fixedCode()
func (f *ErrorToFixMapperTest_ExpectedFix) FixedCode() string { return f.fixedCode }

// ErrorToFixMapperTest_ExpectedFix_Builder_ is ExpectedFix.builder() (the
// trailing underscore keeps it apart from the type
// ErrorToFixMapperTest_ExpectedFix_Builder).
//
// Java: ErrorToFixMapperTest.ExpectedFix.builder()
func ErrorToFixMapperTest_ExpectedFix_Builder_() *ErrorToFixMapperTest_ExpectedFix_Builder {
	return &ErrorToFixMapperTest_ExpectedFix_Builder{}
}

// ErrorToFixMapperTest_ExpectedFix_Builder is the AutoValue builder.
//
// Java: @AutoValue.Builder abstract static class ErrorToFixMapperTest.ExpectedFix.Builder
type ErrorToFixMapperTest_ExpectedFix_Builder struct {
	description  *string
	fixedCode    string
	fixedCodeSet bool
}

// Java: ErrorToFixMapperTest.ExpectedFix.Builder.description(String)
func (b *ErrorToFixMapperTest_ExpectedFix_Builder) Description(description string) *ErrorToFixMapperTest_ExpectedFix_Builder {
	b.description = &description
	return b
}

// Java: ErrorToFixMapperTest.ExpectedFix.Builder.fixedCode(String)
func (b *ErrorToFixMapperTest_ExpectedFix_Builder) FixedCode(fixedCode string) *ErrorToFixMapperTest_ExpectedFix_Builder {
	b.fixedCode = fixedCode
	b.fixedCodeSet = true
	return b
}

// Java: ErrorToFixMapperTest.ExpectedFix.Builder.build()
func (b *ErrorToFixMapperTest_ExpectedFix_Builder) Build() *ErrorToFixMapperTest_ExpectedFix {
	if !b.fixedCodeSet {
		panic(jutil.NewIllegalStateExceptionMsg("Missing required properties: fixedCode"))
	}
	return &ErrorToFixMapperTest_ExpectedFix{description: b.description, fixedCode: b.fixedCode}
}

// Java: ErrorToFixMapperTest.testMissingSuperCall()
func (c *ErrorToFixMapperTest) TestMissingSuperCall() {
	c.AssertExpectedFixes(
		// Must have a super call inside the constructor.
		`class C {
}
class D extends C {
  constructor() {
  }
}
`,
		ErrorToFixMapperTest_ExpectedFix_Builder_().
			FixedCode(
				`class C {
}
class D extends C {
  constructor() {
super();
  }
}
`).
			Build())
}

// Java: ErrorToFixMapperTest.testInvalidSuperCall()
func (c *ErrorToFixMapperTest) TestInvalidSuperCall() {
	// super() constructor call invalid from within method()
	c.AssertExpectedFixes(
		`class C {
  method() {}
}
class D extends C {
  method() {
    return super();
  }
}
`,
		ErrorToFixMapperTest_ExpectedFix_Builder_().
			Description("Call 'super.method' instead").
			FixedCode(
				`class C {
  method() {}
}
class D extends C {
  method() {
    return super.method();
  }
}
`).
			Build())
}

// Java: ErrorToFixMapperTest.testDebugger()
func (c *ErrorToFixMapperTest) TestDebugger() {
	code :=
		`function f() {
  debugger;
}
`
	expectedCode :=
		`function f() {
  
}
`
	c.AssertExpectedFixes(
		code,
		ErrorToFixMapperTest_ExpectedFix_Builder_().
			Description("Remove debugger statement").
			FixedCode(expectedCode).
			Build())
}

// Java: ErrorToFixMapperTest.testEmptyStatement1()
func (c *ErrorToFixMapperTest) TestEmptyStatement1() {
	c.AssertExpectedFixes(
		"var x;;",
		ErrorToFixMapperTest_ExpectedFix_Builder_().Description("Remove empty statement").FixedCode("var x;").Build())
}

// Java: ErrorToFixMapperTest.testEmptyStatement2()
func (c *ErrorToFixMapperTest) TestEmptyStatement2() {
	c.AssertExpectedFixes(
		"var x;;\nvar y;",
		ErrorToFixMapperTest_ExpectedFix_Builder_().
			Description("Remove empty statement").
			FixedCode("var x;\nvar y;").
			Build())
}

// Java: ErrorToFixMapperTest.testEmptyStatement3()
func (c *ErrorToFixMapperTest) TestEmptyStatement3() {
	c.AssertExpectedFixes(
		"function f() {};\nf();",
		ErrorToFixMapperTest_ExpectedFix_Builder_().
			Description("Remove empty statement").
			FixedCode("function f() {}\nf();").
			Build())
}

// Java: ErrorToFixMapperTest.testImplicitNullability1()
func (c *ErrorToFixMapperTest) TestImplicitNullability1() {
	c.options.SetWarningLevel(DiagnosticGroups_LINT_CHECKS, CheckLevel_OFF)
	c.AssertChangesWithExpectedFixes(
		"/** @type {Object} */ var o;",
		"/** @type {?Object} */ var o;",
		"/** @type {!Object} */ var o;")
}

// Java: ErrorToFixMapperTest.testImplicitNullability2()
func (c *ErrorToFixMapperTest) TestImplicitNullability2() {
	c.options.SetWarningLevel(DiagnosticGroups_LINT_CHECKS, CheckLevel_OFF)
	c.AssertChangesWithExpectedFixes(
		"/** @param {Object} o */ function f(o) {}",
		"/** @param {?Object} o */ function f(o) {}",
		"/** @param {!Object} o */ function f(o) {}")
}

// Java: ErrorToFixMapperTest.testImplicitNullability3()
func (c *ErrorToFixMapperTest) TestImplicitNullability3() {
	c.options.SetWarningLevel(DiagnosticGroups_LINT_CHECKS, CheckLevel_OFF)
	originalCode :=
		`/**
 * Some non-ASCII characters: αβγδε
 * @param {Object} o
 */
function f(o) {}
`
	expected1 :=
		`/**
 * Some non-ASCII characters: αβγδε
 * @param {?Object} o
 */
function f(o) {}
`
	expected2 :=
		`/**
 * Some non-ASCII characters: αβγδε
 * @param {!Object} o
 */
function f(o) {}
`
	c.AssertChangesWithExpectedFixes(originalCode, expected1, expected2)
}

// Java: ErrorToFixMapperTest.testMissingNullabilityModifier1()
func (c *ErrorToFixMapperTest) TestMissingNullabilityModifier1() {
	c.options.SetWarningLevel(DiagnosticGroups_ANALYZER_CHECKS, CheckLevel_OFF)
	c.AssertChangesWithExpectedFixes(
		"/** @type {Object} */ var o;",
		"/** @type {!Object} */ var o;",
		"/** @type {?Object} */ var o;")
}

// Java: ErrorToFixMapperTest.testMissingNullabilityModifier2()
func (c *ErrorToFixMapperTest) TestMissingNullabilityModifier2() {
	c.options.SetWarningLevel(DiagnosticGroups_ANALYZER_CHECKS, CheckLevel_OFF)
	c.AssertChangesWithExpectedFixes(
		"/** @param {Object} o */ function f(o) {}",
		"/** @param {!Object} o */ function f(o) {}",
		"/** @param {?Object} o */ function f(o) {}")
}

// Java: ErrorToFixMapperTest.testMissingBangOnEnum()
func (c *ErrorToFixMapperTest) TestMissingBangOnEnum() {
	c.options.SetWarningLevel(DiagnosticGroups_LINT_CHECKS, CheckLevel_OFF)
	prelude := "/** @enum {number} */ var Enum;\n"
	c.AssertChanges(prelude+"/** @type {Enum} */ var o;", prelude+"/** @type {!Enum} */ var o;")
}

// Java: ErrorToFixMapperTest.testMissingBangOnTypedef()
func (c *ErrorToFixMapperTest) TestMissingBangOnTypedef() {
	c.options.SetWarningLevel(DiagnosticGroups_LINT_CHECKS, CheckLevel_OFF)
	prelude := "/** @typedef {number} */ var Num;\n"
	c.AssertChanges(prelude+"/** @type {Num} */ var o;", prelude+"/** @type {!Num} */ var o;")
}

// Java: ErrorToFixMapperTest.testMissingNullabilityModifier3()
func (c *ErrorToFixMapperTest) TestMissingNullabilityModifier3() {
	c.options.SetWarningLevel(DiagnosticGroups_ANALYZER_CHECKS, CheckLevel_OFF)
	originalCode :=
		`/**
 * Some non-ASCII characters: αβγδε
 * @param {Object} o
 */
function f(o) {}
`
	expected1 :=
		`/**
 * Some non-ASCII characters: αβγδε
 * @param {!Object} o
 */
function f(o) {}
`
	expected2 :=
		`/**
 * Some non-ASCII characters: αβγδε
 * @param {?Object} o
 */
function f(o) {}
`
	c.AssertChangesWithExpectedFixes(originalCode, expected1, expected2)
}

// Java: ErrorToFixMapperTest.testNullMissingNullabilityModifier1()
func (c *ErrorToFixMapperTest) TestNullMissingNullabilityModifier1() {
	c.options.SetWarningLevel(DiagnosticGroups_ANALYZER_CHECKS, CheckLevel_OFF)
	c.AssertChanges("/** @type {Object} */ var x = null;", "/** @type {?Object} */ var x = null;")
}

// Java: ErrorToFixMapperTest.testNullMissingNullabilityModifier2()
func (c *ErrorToFixMapperTest) TestNullMissingNullabilityModifier2() {
	c.options.SetWarningLevel(DiagnosticGroups_ANALYZER_CHECKS, CheckLevel_OFF)
	c.AssertNoChanges("/** @type {?Object} */ var x = null;")
}

// Java: ErrorToFixMapperTest.testMissingNullabilityModifier_nonNullValue()
func (c *ErrorToFixMapperTest) TestMissingNullabilityModifier_nonNullValue() {
	c.options.SetWarningLevel(DiagnosticGroups_ANALYZER_CHECKS, CheckLevel_OFF)
	c.AssertChangesWithExpectedFixes(
		"/** @type {Object} */ var o = {};",
		"/** @type {!Object} */ var o = {};",
		"/** @type {?Object} */ var o = {};")
}

// Java: ErrorToFixMapperTest.testRedundantNullabilityModifier1()
func (c *ErrorToFixMapperTest) TestRedundantNullabilityModifier1() {
	c.AssertChanges("/** @type {!string} */ var x;", "/** @type {string} */ var x;")
}

// Java: ErrorToFixMapperTest.testRedundantNullabilityModifier2()
func (c *ErrorToFixMapperTest) TestRedundantNullabilityModifier2() {
	c.AssertChanges(
		"/** @type {!{foo: string, bar: !string}} */ var x;",
		"/** @type {{foo: string, bar: string}} */ var x;")
}

// Java: ErrorToFixMapperTest.testRedeclaration()
func (c *ErrorToFixMapperTest) TestRedeclaration() {
	code := "function f() { var x; var x; }"
	expectedCode := "function f() { var x; }"
	c.AssertExpectedFixes(
		code,
		ErrorToFixMapperTest_ExpectedFix_Builder_().
			Description("Remove redundant declaration").
			FixedCode(expectedCode).
			Build())
}

// Java: ErrorToFixMapperTest.testRedeclaration_multipleVars1()
func (c *ErrorToFixMapperTest) TestRedeclaration_multipleVars1() {
	code := "function f() { var x; var x, y; }"
	expectedCode := "function f() { var x; var y; }"
	c.AssertChanges(code, expectedCode)
}

// Java: ErrorToFixMapperTest.testRedeclaration_multipleVars2()
func (c *ErrorToFixMapperTest) TestRedeclaration_multipleVars2() {
	code := "function f() { var x; var y, x; }"
	expectedCode := "function f() { var x; var y; }"
	c.AssertChanges(code, expectedCode)
}

// Java: ErrorToFixMapperTest.testRedeclaration_withValue()
func (c *ErrorToFixMapperTest) TestRedeclaration_withValue() {
	code :=
		`function f() {
  var x;
  var x = 0;
}
`
	expectedCode :=
		`function f() {
  var x;
  x = 0;
}
`
	c.AssertChanges(code, expectedCode)
}

// Java: ErrorToFixMapperTest.testRedeclaration_multipleVars_withValue1()
func (c *ErrorToFixMapperTest) TestRedeclaration_multipleVars_withValue1() {
	code :=
		`function f() {
  var x;
  var x = 0, y;
}
`
	expectedCode :=
		`function f() {
  var x;
  x = 0;
var y;
}
`
	c.AssertChanges(code, expectedCode)
}

// Java: ErrorToFixMapperTest.testRedeclaration_multipleVars_withValue2()
func (c *ErrorToFixMapperTest) TestRedeclaration_multipleVars_withValue2() {
	code :=
		`function f() {
  var x;
  var y, x = 0;
}
`
	expectedCode :=
		`function f() {
  var x;
  var y;
x = 0;
}
`
	c.AssertChanges(code, expectedCode)
}

// Make sure the vars stay in the same order, so that in case the get*
// functions have side effects, we don't change the order they're called in.
// Java: ErrorToFixMapperTest.testRedeclaration_multipleVars_withValue3()
func (c *ErrorToFixMapperTest) TestRedeclaration_multipleVars_withValue3() {
	code :=
		`function f() {
  var y;
  var x = getX(), y = getY(), z = getZ();
}
`
	expectedCode :=
		`function f() {
  var y;
  var x = getX();
y = getY();
var z = getZ();
}
`
	c.AssertChanges(code, expectedCode)
}

// Java: ErrorToFixMapperTest.testRedeclaration_multipleVars_withValue4()
func (c *ErrorToFixMapperTest) TestRedeclaration_multipleVars_withValue4() {
	code :=
		`function f() {
  var x;
  var x = getX(), y = getY(), z = getZ();
}
`
	expectedCode :=
		`function f() {
  var x;
  x = getX();
var y = getY(), z = getZ();
}
`
	c.AssertChanges(code, expectedCode)
}

// Java: ErrorToFixMapperTest.testRedeclaration_multipleVars_withValue5()
func (c *ErrorToFixMapperTest) TestRedeclaration_multipleVars_withValue5() {
	code :=
		`function f() {
  var z;
  var x = getX(), y = getY(), z = getZ();
}
`
	expectedCode :=
		`function f() {
  var z;
  var x = getX(), y = getY();
z = getZ();
}
`
	c.AssertChanges(code, expectedCode)
}

// Java: ErrorToFixMapperTest.testRedeclarationOfParam()
func (c *ErrorToFixMapperTest) TestRedeclarationOfParam() {
	c.AssertExpectedFixes(
		"function f(x) { var x = 3; }",
		ErrorToFixMapperTest_ExpectedFix_Builder_().
			Description("Convert redundant declaration to assignment").
			FixedCode("function f(x) { x = 3; }").
			Build())
}

// Java: ErrorToFixMapperTest.testRedeclaration_params()
func (c *ErrorToFixMapperTest) TestRedeclaration_params() {
	c.AssertNoChanges("function f(x, x) {}")
}

// Java: ErrorToFixMapperTest.testEarlyReference()
func (c *ErrorToFixMapperTest) TestEarlyReference() {
	code := "if (x < 0) alert(1);\nvar x;"
	expectedCode := "var x;\n" + code
	c.AssertExpectedFixes(
		code,
		ErrorToFixMapperTest_ExpectedFix_Builder_().
			Description("Insert var declaration statement").
			FixedCode(expectedCode).
			Build())
}

// Java: ErrorToFixMapperTest.testEarlyReferenceInFunction()
func (c *ErrorToFixMapperTest) TestEarlyReferenceInFunction() {
	code := "function f() {\n  if (x < 0) alert(1);\nvar x;\n}"
	expectedCode := "function f() {\n  var x;\nif (x < 0) alert(1);\nvar x;\n}"
	c.AssertExpectedFixes(
		code,
		ErrorToFixMapperTest_ExpectedFix_Builder_().
			Description("Insert var declaration statement").
			FixedCode(expectedCode).
			Build())
}

// Java: ErrorToFixMapperTest.testInsertSemicolon1()
func (c *ErrorToFixMapperTest) TestInsertSemicolon1() {
	code := "var x = 3"
	expectedCode := "var x = 3;"
	c.AssertChanges(code, expectedCode)
}

// Java: ErrorToFixMapperTest.testInsertSemicolon2()
func (c *ErrorToFixMapperTest) TestInsertSemicolon2() {
	code := "function f() { return 'it' }"
	expectedCode := "function f() { return 'it'; }"
	c.AssertChanges(code, expectedCode)
}

// Java: ErrorToFixMapperTest.testFixProvides_sort()
func (c *ErrorToFixMapperTest) TestFixProvides_sort() {
	c.AssertChanges(
		`/** @fileoverview foo */

goog.provide('b');
goog.provide('a');
goog.provide('c');

alert(1);
`,
		`/** @fileoverview foo */

goog.provide('a');
goog.provide('b');
goog.provide('c');

alert(1);
`)
}

// Java: ErrorToFixMapperTest.testFixProvides_deduplicate()
func (c *ErrorToFixMapperTest) TestFixProvides_deduplicate() {
	c.AssertChanges(
		`/** @fileoverview foo */

goog.provide('a');
goog.provide('b');
goog.provide('a');

alert(1);
`,
		`/** @fileoverview foo */

goog.provide('a');
goog.provide('b');

alert(1);
`)
}

// Java: ErrorToFixMapperTest.testFixProvides_alreadySorted()
func (c *ErrorToFixMapperTest) TestFixProvides_alreadySorted() {
	c.AssertNoChanges(
		`/** @fileoverview foo */

goog.provide('a');
goog.provide('b');
goog.provide('c');

alert(1);
`)
}

// Java: ErrorToFixMapperTest.testFixNonAliasedRequire()
func (c *ErrorToFixMapperTest) TestFixNonAliasedRequire() {
	c.AssertChanges(
		`/** @fileoverview foo */

goog.module('m');

goog.require('x');
`+
			c.UseInCode("x"),
		`/** @fileoverview foo */

goog.module('m');

const x = goog.require('x');
`+
			c.UseInCode("x"))
}

// Java: ErrorToFixMapperTest.testFixNonAliasedRequireType()
func (c *ErrorToFixMapperTest) TestFixNonAliasedRequireType() {
	c.AssertChanges(
		`/** @fileoverview foo */

goog.module('m');

goog.requireType('x');
`+
			c.UseInType("x"),
		`/** @fileoverview foo */

goog.module('m');

const x = goog.requireType('x');
`+
			c.UseInType("x"))
}

// Java: ErrorToFixMapperTest.testFixRequires_sortStandaloneOnly()
func (c *ErrorToFixMapperTest) TestFixRequires_sortStandaloneOnly() {
	c.AssertChanges(
		c.FileWithImports(
			`goog.require('b');
goog.requireType('d');
goog.requireType('c');
goog.forwardDeclare('f');
goog.require('a');
goog.forwardDeclare('e');
`+
				c.UseInCode("a", "b")+
				c.UseInType("c", "d", "e", "f")),
		c.FileWithImports(
			`goog.forwardDeclare('e');
goog.forwardDeclare('f');
goog.require('a');
goog.require('b');
goog.requireType('c');
goog.requireType('d');
`+
				c.UseInCode("a", "b")+
				c.UseInType("c", "d", "e", "f")))
}

// Java: ErrorToFixMapperTest.testFixRequires_sortAllTypes()
func (c *ErrorToFixMapperTest) TestFixRequires_sortAllTypes() {
	c.AssertChanges(
		c.FileWithImports(
			`goog.requireType('a');
goog.require('b');
const f = goog.require('f');
const {d} = goog.require('d');
const e = goog.requireType('e');
const {c} = goog.requireType('c');
goog.forwardDeclare('g');
`+
				c.UseInCode("a", "d", "f")+
				c.UseInType("b", "c", "e", "g")),
		c.FileWithImports(
			`const e = goog.requireType('e');
const f = goog.require('f');
const {c} = goog.requireType('c');
const {d} = goog.require('d');
goog.forwardDeclare('g');
goog.require('b');
goog.requireType('a');
`+
				c.UseInCode("a", "d", "f")+
				c.UseInType("b", "c", "e", "g")))
}

// Java: ErrorToFixMapperTest.testLineCommentDoesNotGetDeleted()
func (c *ErrorToFixMapperTest) TestLineCommentDoesNotGetDeleted() {
	c.options.SetParseJsDocDocumentation(parsing.Config_JsDocParsing_INCLUDE_ALL_COMMENTS)
	c.AssertChanges(
		c.FileWithImports(
			`goog.require('d');
// dummy
goog.require('b');
`+
				c.UseInCode("d", "b")), //
		c.FileWithImports(
			`// dummy
goog.require('b');
goog.require('d');
`+
				c.UseInCode("d", "b")))
}

// Java: ErrorToFixMapperTest.testBlockCommentDoesNotGetDeleted()
func (c *ErrorToFixMapperTest) TestBlockCommentDoesNotGetDeleted() {
	c.options.SetParseJsDocDocumentation(parsing.Config_JsDocParsing_INCLUDE_ALL_COMMENTS)
	c.AssertChanges(
		c.FileWithImports(
			`goog.require('d');
/* dummy */
goog.require('b');
`+
				c.UseInCode("d", "b")),
		c.FileWithImports(
			`/* dummy */
goog.require('b');
goog.require('d');
`+
				c.UseInCode("d", "b")))
}

// Java: ErrorToFixMapperTest.testFirstBlockCommentDoesNotGetDeleted()
func (c *ErrorToFixMapperTest) TestFirstBlockCommentDoesNotGetDeleted() {
	c.options.SetParseJsDocDocumentation(parsing.Config_JsDocParsing_INCLUDE_ALL_COMMENTS)
	c.AssertChanges(
		c.FileWithImports(
			`/* dummy */
goog.require('d');
goog.require('b');
`+
				c.UseInCode("d", "b")),
		c.FileWithImports(
			`goog.require('b');
/* dummy */
goog.require('d');
`+
				c.UseInCode("d", "b")))
}

// Java: ErrorToFixMapperTest.testIndividualCommentsMoveCorrectly()
func (c *ErrorToFixMapperTest) TestIndividualCommentsMoveCorrectly() {
	c.options.SetParseJsDocDocumentation(parsing.Config_JsDocParsing_INCLUDE_ALL_COMMENTS)
	c.AssertChanges(
		c.FileWithImports(
			`// dummy 1
goog.require('d');
// dummy 2
goog.require('b');
`+
				c.UseInCode("d", "b")),
		c.FileWithImports(
			`// dummy 2
goog.require('b');
// dummy 1
goog.require('d');
`+
				c.UseInCode("d", "b")))
}

// Java: ErrorToFixMapperTest.testIndividualCommentsMoveCorrectly2()
func (c *ErrorToFixMapperTest) TestIndividualCommentsMoveCorrectly2() {
	c.options.SetParseJsDocDocumentation(parsing.Config_JsDocParsing_INCLUDE_ALL_COMMENTS)
	c.AssertChanges(
		c.FileWithImports(
			`// dummy 1
goog.require('d');
// dummy 2
goog.require('b');
// dummy 3
goog.require('a');
`+
				c.UseInCode("d", "b", "a")),
		c.FileWithImports(
			`// dummy 3
goog.require('a');
// dummy 2
goog.require('b');
// dummy 1
goog.require('d');
`+
				c.UseInCode("d", "b", "a")))
}

// Java: ErrorToFixMapperTest.testIndividualCommentsMoveCorrectly3()
func (c *ErrorToFixMapperTest) TestIndividualCommentsMoveCorrectly3() {
	c.options.SetParseJsDocDocumentation(parsing.Config_JsDocParsing_INCLUDE_ALL_COMMENTS)
	c.AssertChanges(
		c.FileWithImports(
			`// dummy 1
goog.require('d');
// dummy 2
goog.require('b');
goog.require('a');
`+
				c.UseInCode("d", "b", "a")),
		c.FileWithImports(
			`goog.require('a');
// dummy 2
goog.require('b');
// dummy 1
goog.require('d');
`+
				c.UseInCode("d", "b", "a")))
}

// Java: ErrorToFixMapperTest.testIndividualCommentsMoveCorrectly4()
func (c *ErrorToFixMapperTest) TestIndividualCommentsMoveCorrectly4() {
	c.options.SetParseJsDocDocumentation(parsing.Config_JsDocParsing_INCLUDE_ALL_COMMENTS)
	c.AssertChanges(
		c.FileWithImports(
			`// dummy 1
goog.require('d');
goog.require('b');
goog.require('a');
`+
				c.UseInCode("d", "b", "a")),
		c.FileWithImports(
			`goog.require('a');
goog.require('b');
// dummy 1
goog.require('d');
`+
				c.UseInCode("d", "b", "a")))
}

// Java: ErrorToFixMapperTest.testMultipleCommentsMoveCorrectly()
func (c *ErrorToFixMapperTest) TestMultipleCommentsMoveCorrectly() {
	c.options.SetParseJsDocDocumentation(parsing.Config_JsDocParsing_INCLUDE_ALL_COMMENTS)
	c.AssertChanges(
		c.FileWithImports(
			`// dummy 1
// dummy 2
goog.require('d');
goog.require('b');
goog.require('a');
`+
				c.UseInCode("d", "b", "a")),
		c.FileWithImports(
			`goog.require('a');
goog.require('b');
// dummy 1
// dummy 2
goog.require('d');
`+
				c.UseInCode("d", "b", "a")))
}

// Java: ErrorToFixMapperTest.testMultipleCommentsMoveCorrectly2()
func (c *ErrorToFixMapperTest) TestMultipleCommentsMoveCorrectly2() {
	c.options.SetParseJsDocDocumentation(parsing.Config_JsDocParsing_INCLUDE_ALL_COMMENTS)
	c.AssertChanges(
		c.FileWithImports(
			`// dummy 1


// dummy 2
goog.require('d');
goog.require('b');
goog.require('a');
`+
				c.UseInCode("d", "b", "a")),
		c.FileWithImports(
			`goog.require('a');
goog.require('b');
// dummy 1


// dummy 2
goog.require('d');
`+
				c.UseInCode("d", "b", "a")))
}

// Java: ErrorToFixMapperTest.testFirstLineCommentDoesNotGetDeleted()
func (c *ErrorToFixMapperTest) TestFirstLineCommentDoesNotGetDeleted() {
	c.options.SetParseJsDocDocumentation(parsing.Config_JsDocParsing_INCLUDE_ALL_COMMENTS)
	c.AssertChanges(
		c.FileWithImports(
			`// dummy
goog.require('d');
goog.require('b');
`+
				c.UseInCode("d", "b")),
		c.FileWithImports(
			`goog.require('b');
// dummy
goog.require('d');
`+
				c.UseInCode("d", "b")))
}

// Java: ErrorToFixMapperTest.testBothJSDocAndNonJSDocCommentsTogether()
func (c *ErrorToFixMapperTest) TestBothJSDocAndNonJSDocCommentsTogether() {
	c.options.SetParseJsDocDocumentation(parsing.Config_JsDocParsing_INCLUDE_ALL_COMMENTS)
	c.AssertChanges(
		c.FileWithImports(
			`/** JSDoc gets preserved*/
// dummy
goog.require('d');
goog.require('b');
`+
				c.UseInCode("d", "b")),
		c.FileWithImports(
			`goog.require('b');
// dummy
/** JSDoc gets preserved*/
goog.require('d');
`+
				c.UseInCode("d", "b")))
}

// Java: ErrorToFixMapperTest.testFileOverviewCommentsPreservedAfterSortingProvides()
func (c *ErrorToFixMapperTest) TestFileOverviewCommentsPreservedAfterSortingProvides() {
	c.options.SetParseJsDocDocumentation(parsing.Config_JsDocParsing_INCLUDE_ALL_COMMENTS)
	c.AssertChanges(
		`// Copyright 2011 The Closure Library Authors. All Rights Reserved.
/**
 * @fileoverview Date/time formatting symbols for all locales.
 * @suppress {const} */
// clang-format off
goog.provide('goog.i18n.DateTimeSymbols_af');
goog.provide('goog.i18n.DateTimeSymbols');
goog.provide('goog.i18n.DateTimeSymbolsType');
`,
		`// Copyright 2011 The Closure Library Authors. All Rights Reserved.
/**
 * @fileoverview Date/time formatting symbols for all locales.
 * @suppress {const} */
// clang-format off
goog.provide('goog.i18n.DateTimeSymbols');
goog.provide('goog.i18n.DateTimeSymbolsType');
goog.provide('goog.i18n.DateTimeSymbols_af');
`)
}

// Java: ErrorToFixMapperTest.testFileOverviewCommentsPreservedAfterSortingProvides2()
func (c *ErrorToFixMapperTest) TestFileOverviewCommentsPreservedAfterSortingProvides2() {
	c.options.SetParseJsDocDocumentation(parsing.Config_JsDocParsing_INCLUDE_ALL_COMMENTS)
	c.AssertChanges(
		`// clang-format off
/**
 * @fileoverview Date/time formatting symbols for all locales.
 * @suppress {const} */
goog.provide('goog.i18n.DateTimeSymbols_af');
goog.provide('goog.i18n.DateTimeSymbols');
goog.provide('goog.i18n.DateTimeSymbolsType');
`,
		`// clang-format off
/**
 * @fileoverview Date/time formatting symbols for all locales.
 * @suppress {const} */
goog.provide('goog.i18n.DateTimeSymbols');
goog.provide('goog.i18n.DateTimeSymbolsType');
goog.provide('goog.i18n.DateTimeSymbols_af');
`)
}

// Java: ErrorToFixMapperTest.testBothJSDocAndNonJSDocCommentsTogether2()
func (c *ErrorToFixMapperTest) TestBothJSDocAndNonJSDocCommentsTogether2() {
	c.options.SetParseJsDocDocumentation(parsing.Config_JsDocParsing_INCLUDE_ALL_COMMENTS)
	c.AssertChanges(
		c.FileWithImports(
			`// dummy
/** JSDoc gets preserved*/
goog.require('d');
goog.require('b');
`+
				c.UseInCode("d", "b")),
		c.FileWithImports(
			`goog.require('b');
// dummy
/** JSDoc gets preserved*/
goog.require('d');
`+
				c.UseInCode("d", "b")))
}

// Java: ErrorToFixMapperTest.testFixRequires_deduplicate_standalone()
func (c *ErrorToFixMapperTest) TestFixRequires_deduplicate_standalone() {
	c.AssertChanges(
		c.FileWithImports(
			`goog.require('a');
goog.require('a');
goog.require('b');
goog.requireType('b');
goog.requireType('c');
goog.requireType('c');
goog.forwardDeclare('b');
goog.forwardDeclare('c');
goog.forwardDeclare('d');
`+
				c.UseInCode("a", "b")+
				c.UseInType("c", "d")),
		c.FileWithImports(
			`goog.forwardDeclare('d');
goog.require('a');
goog.require('b');
goog.requireType('c');
`+
				c.UseInCode("a", "b")+
				c.UseInType("c", "d")))
}

// Java: ErrorToFixMapperTest.testFixRequires_preserveMultipleAliases()
func (c *ErrorToFixMapperTest) TestFixRequires_preserveMultipleAliases() {
	// a3 changes from goog.requireType to goog.require because we enforce all imports for the same
	// namespace to have the same strength when rewriting.
	c.AssertChanges(
		c.FileWithImports(
			`const a3 = goog.requireType('a');
const a2 = goog.require('a');
const a1 = goog.require('a');
`+
				c.UseInCode("a1", "a2")+
				c.UseInType("a3")),
		c.FileWithImports(
			`const a1 = goog.require('a');
const a2 = goog.require('a');
const a3 = goog.require('a');
`+
				c.UseInCode("a1", "a2")+
				c.UseInType("a3")))
}

// Java: ErrorToFixMapperTest.testFixRequires_mergeDestructuring()
func (c *ErrorToFixMapperTest) TestFixRequires_mergeDestructuring() {
	c.AssertChanges(
		c.FileWithImports(
			`const {c, a: a2} = goog.require('a');
const {b, a: a1} = goog.require('a');
const {a} = goog.require('a');
const {} = goog.require('a');
`+
				c.UseInCode("a", "a1", "a2", "b", "c")),
		c.FileWithImports(
			`const {a, a: a1, a: a2, b, c} = goog.require('a');
`+
				c.UseInCode("a", "a1", "a2", "b", "c")))
}

// Java: ErrorToFixMapperTest.testFixRequires_mergeDestructuring_multiplePrimitives()
func (c *ErrorToFixMapperTest) TestFixRequires_mergeDestructuring_multiplePrimitives() {
	c.AssertChanges(
		c.FileWithImports(
			`const {c, a} = goog.require('a');
const {b, a} = goog.requireType('a');
`+
				c.UseInCode("a", "c")+
				c.UseInType("a", "b")),
		c.FileWithImports(
			`const {a, b, c} = goog.require('a');
`+
				c.UseInCode("a", "c")+
				c.UseInType("a", "b")))
}

// Java: ErrorToFixMapperTest.testFixRequires_nonAliasedRequire()
func (c *ErrorToFixMapperTest) TestFixRequires_nonAliasedRequire() {
	c.AssertChanges(
		c.FileWithImports(
			`goog.require('a');
`+
				c.UseInCode("a")),
		c.FileWithImports(
			`const a = goog.require('a');
`+
				c.UseInCode("a")))
}

// Java: ErrorToFixMapperTest.testFixRequires_mutlipleFixesSpecifySameRequire()
func (c *ErrorToFixMapperTest) TestFixRequires_mutlipleFixesSpecifySameRequire() {
	c.AssertChanges(
		c.FileWithImports(
			`goog.require('a.b');
`+
				c.UseInCode("a.b", "a.b", "a.b")),
		c.FileWithImports(
			`const b = goog.require('a.b');
`+
				c.UseInCode("b", "b", "b")))
}

// Java: ErrorToFixMapperTest.testFixRequires_emptyDestructuring_alone()
func (c *ErrorToFixMapperTest) TestFixRequires_emptyDestructuring_alone() {
	c.t.Skip("@Ignore")
	// It would be nice to dedeuplicate the destructuring require with the aliased require,
	// but it's not clear that this pattern is common enough to warrant special casing
	c.AssertChanges(
		c.FileWithImports(
			`const {} = goog.require('a');
`+
				c.UseInCode("a")),
		c.FileWithImports(
			`goog.require('a');
`+
				c.UseInCode("a")))
}

// Java: ErrorToFixMapperTest.testFixRequires_emptyDestructuringStandaloneBySamePrimitive()
func (c *ErrorToFixMapperTest) TestFixRequires_emptyDestructuringStandaloneBySamePrimitive() {
	c.AssertChanges(
		c.FileWithImports(
			`const {} = goog.require('a');
goog.require('a');
`+
				c.UseInCode("a")),
		c.FileWithImports(
			`goog.require('a');
`+
				c.UseInCode("a")))
}

// Java: ErrorToFixMapperTest.testFixRequires_emptyDestructuringStandaloneByStrongerPrimitive()
func (c *ErrorToFixMapperTest) TestFixRequires_emptyDestructuringStandaloneByStrongerPrimitive() {
	c.AssertChanges(
		c.FileWithImports(
			`const {} = goog.requireType('a');
goog.require('a');
`+
				c.UseInCode("a")),
		c.FileWithImports(
			`goog.require('a');
`+
				c.UseInCode("a")))
}

// Java: ErrorToFixMapperTest.testFixRequires_emptyDestructuringStandaloneByWeakerPrimitive()
func (c *ErrorToFixMapperTest) TestFixRequires_emptyDestructuringStandaloneByWeakerPrimitive() {
	c.t.Skip("@Ignore")
	// It would be nice to dedeuplicate the destructuring require with the aliased require,
	// but it's not clear that this pattern is common enough to warrant special casing
	c.AssertChanges(
		c.FileWithImports(
			`const {} = goog.require('a');
goog.requireType('a');
`+
				c.UseInCode("a")),
		c.FileWithImports(
			`goog.require('a');
`+
				c.UseInCode("a")))
}

// Java: ErrorToFixMapperTest.testFixRequires_emptyDestructuringAliasedBySamePrimitive()
func (c *ErrorToFixMapperTest) TestFixRequires_emptyDestructuringAliasedBySamePrimitive() {
	c.AssertChanges(
		c.FileWithImports(
			`const {} = goog.require('a');
const a = goog.require('a');
`+
				c.UseInCode("a")),
		c.FileWithImports(
			`const a = goog.require('a');
`+
				c.UseInCode("a")))
}

// Java: ErrorToFixMapperTest.testFixRequires_emptyDestructuringAliasedByStrongerPrimitive()
func (c *ErrorToFixMapperTest) TestFixRequires_emptyDestructuringAliasedByStrongerPrimitive() {
	c.AssertChanges(
		c.FileWithImports(
			`const {} = goog.requireType('a');
const a = goog.require('a');
`+
				c.UseInCode("a")),
		c.FileWithImports(
			`const a = goog.require('a');
`+
				c.UseInCode("a")))
}

// Java: ErrorToFixMapperTest.testFixRequires_emptyDestructuringAliasedByWeakerPrimitive()
func (c *ErrorToFixMapperTest) TestFixRequires_emptyDestructuringAliasedByWeakerPrimitive() {
	c.AssertChanges(
		c.FileWithImports(
			`const {} = goog.require('a');
const a = goog.requireType('a');
`+
				c.UseInCode("a")),
		c.FileWithImports(
			`const a = goog.require('a');
`+
				c.UseInCode("a")))
}

// Java: ErrorToFixMapperTest.testFixRequires_aliasPreservedWhenDestructuring()
func (c *ErrorToFixMapperTest) TestFixRequires_aliasPreservedWhenDestructuring() {
	c.AssertNoChanges(
		c.FileWithImports(
			`const a = goog.require('a');
const b = goog.require('b');
const {a1, a2} = goog.require('a');
const {b1, b2} = goog.require('b');
` +
				c.UseInCode("a1", "a2", "b", "b1", "b2") +
				c.UseInType("a")))
}

// The Go method name starts with Test so that RunJUnit4 runs it
// (Go-only: language-adapter).
//
// Java: ErrorToFixMapperTest.attachesInlineJsDocToParam()
func (c *ErrorToFixMapperTest) TestAttachesInlineJsDocToParam() {
	c.AssertNoChanges(
		` class C {
  inlineJSDocWithDefault(/** boolean= */ isSomething) {}
}
`)
}

// The Go method name starts with Test so that RunJUnit4 runs it
// (Go-only: language-adapter).
//
// Java: ErrorToFixMapperTest.attachesInlineJsDocToDefaultParam()
func (c *ErrorToFixMapperTest) TestAttachesInlineJsDocToDefaultParam() {
	c.AssertNoChanges(
		` class C {
  inlineJSDocWithDefault(/** boolean= */ isSomething = true) {}
}
`)
}

// Java: ErrorToFixMapperTest.testFixRequires_standaloneAliasedBySamePrimitive()
func (c *ErrorToFixMapperTest) TestFixRequires_standaloneAliasedBySamePrimitive() {
	c.AssertChanges(
		c.FileWithImports(
			`const a = goog.require('a');
goog.require('a');
`+
				c.UseInCode("a")),
		c.FileWithImports(
			`const a = goog.require('a');
`+
				c.UseInCode("a")))
}

// Java: ErrorToFixMapperTest.testFixRequires_standaloneAliasedByStrongerPrimitive()
func (c *ErrorToFixMapperTest) TestFixRequires_standaloneAliasedByStrongerPrimitive() {
	c.AssertChanges(
		c.FileWithImports(
			`const a = goog.require('a');
goog.requireType('a');
`+
				c.UseInCode("a")),
		c.FileWithImports(
			`const a = goog.require('a');
`+
				c.UseInCode("a")))
}

// Java: ErrorToFixMapperTest.testFixRequires_standaloneAliasedByWeakerPrimitive()
func (c *ErrorToFixMapperTest) TestFixRequires_standaloneAliasedByWeakerPrimitive() {
	c.AssertChanges(
		c.FileWithImports(
			`const a = goog.requireType('a');
goog.require('a');
`+
				c.UseInCode("a")),
		c.FileWithImports(
			`const a = goog.require('a');
`+
				c.UseInCode("a")))
}

// Java: ErrorToFixMapperTest.testFixRequires_standaloneDestructuredBySamePrimitive()
func (c *ErrorToFixMapperTest) TestFixRequires_standaloneDestructuredBySamePrimitive() {
	c.AssertChanges(
		c.FileWithImports(
			`const {a} = goog.require('a');
goog.require('a');
`+
				c.UseInCode("a")),
		c.FileWithImports(
			`const {a} = goog.require('a');
`+
				c.UseInCode("a")))
}

// Java: ErrorToFixMapperTest.testFixRequires_standaloneDestructuredByStrongerPrimitive()
func (c *ErrorToFixMapperTest) TestFixRequires_standaloneDestructuredByStrongerPrimitive() {
	c.AssertChanges(
		c.FileWithImports(
			`const {a} = goog.require('a');
goog.requireType('a');
`+
				c.UseInCode("a")),
		c.FileWithImports(
			`const {a} = goog.require('a');
`+
				c.UseInCode("a")))
}

// Java: ErrorToFixMapperTest.testFixRequires_standaloneDestructuredByWeakerPrimitive()
func (c *ErrorToFixMapperTest) TestFixRequires_standaloneDestructuredByWeakerPrimitive() {
	c.AssertChanges(
		c.FileWithImports(
			`const {a} = goog.requireType('a');
goog.require('a');
`+
				c.UseInCode("a")),
		c.FileWithImports(
			`const {a} = goog.require('a');
`+
				c.UseInCode("a")))
}

// Java: ErrorToFixMapperTest.testFixRequires_varAndLetBecomeConstIfUnsorted()
func (c *ErrorToFixMapperTest) TestFixRequires_varAndLetBecomeConstIfUnsorted() {
	c.AssertChanges(
		c.FileWithImports(
			`var b = goog.require('b');
let a = goog.require('a');
`+
				c.UseInCode("a", "b")),
		c.FileWithImports(
			`const a = goog.require('a');
const b = goog.require('b');
`+
				c.UseInCode("a", "b")))
}

// Java: ErrorToFixMapperTest.testFixRequires_varAndLetDoNotBecomeConstIfAlreadySorted()
func (c *ErrorToFixMapperTest) TestFixRequires_varAndLetDoNotBecomeConstIfAlreadySorted() {
	// TODO(tjgq): Consider rewriting to const even when already sorted.
	c.AssertNoChanges(
		c.FileWithImports(
			`let a = goog.require('a');
var b = goog.require('b');
` +
				c.UseInCode("a", "b")))
}

// Java: ErrorToFixMapperTest.testFixRequires_preserveJsDoc_whenSorting()
func (c *ErrorToFixMapperTest) TestFixRequires_preserveJsDoc_whenSorting() {
	c.AssertChanges(
		c.FileWithImports(
			`const c = goog.require('c');
/**
 * @suppress {extraRequire} Because I said so.
 */
const b = goog.require('b');
const a = goog.require('a');
`+
				c.UseInCode("a", "b", "c")),
		c.FileWithImports(
			`const a = goog.require('a');
/**
 * @suppress {extraRequire} Because I said so.
 */
const b = goog.require('b');
const c = goog.require('c');
`+
				c.UseInCode("a", "b", "c")))
}

// Java: ErrorToFixMapperTest.testFixRequires_preserveJsDoc_whenMergingStandalone()
func (c *ErrorToFixMapperTest) TestFixRequires_preserveJsDoc_whenMergingStandalone() {
	c.AssertChanges(
		c.FileWithImports(
			`/**
 * @suppress {extraRequire} Because I said so.
 */
goog.require('a');
goog.requireType('a');
`+
				c.UseInCode("a")),
		c.FileWithImports(
			`/**
 * @suppress {extraRequire} Because I said so.
 */
goog.require('a');
`+
				c.UseInCode("a")))
}

// Java: ErrorToFixMapperTest.testFixRequires_singleDestructure_shorthandProperty()
func (c *ErrorToFixMapperTest) TestFixRequires_singleDestructure_shorthandProperty() {
	c.AssertChanges(
		c.FileWithImports(
			`const {b: b} = goog.require('a');
`+
				c.UseInCode("b")),
		c.FileWithImports(
			`const {b} = goog.require('a');
`+
				c.UseInCode("b")))
}

// Java: ErrorToFixMapperTest.testFixRequires_preserveJsDoc_whenMergingDestructures_single()
func (c *ErrorToFixMapperTest) TestFixRequires_preserveJsDoc_whenMergingDestructures_single() {
	c.AssertChanges(
		c.FileWithImports(
			`/**
 * @suppress {extraRequire} Because I said so.
 */
const {b} = goog.require('a');
const {c} = goog.require('a');
`+
				c.UseInCode("b", "c")),
		c.FileWithImports(
			`/**
 * @suppress {extraRequire} Because I said so.
 */
const {b, c} = goog.require('a');
`+
				c.UseInCode("b", "c")))
}

// Java: ErrorToFixMapperTest.testFixRequires_preserveJsDoc_whenMergingDestructures_multiple()
func (c *ErrorToFixMapperTest) TestFixRequires_preserveJsDoc_whenMergingDestructures_multiple() {
	// TODO(tjgq): Consider merging multiple @suppress annotations into a single comment.
	c.AssertChanges(
		c.FileWithImports(
			`/**
 * @suppress {extraRequire} Because I said so.
 */
const {b} = goog.require('a');
/**
 * @suppress {extraRequire} Because I rule.
 */
const {c} = goog.require('a');
`+
				c.UseInCode("b", "c")),
		c.FileWithImports(
			`/**
 * @suppress {extraRequire} Because I said so.
 */
/**
 * @suppress {extraRequire} Because I rule.
 */
const {b, c} = goog.require('a');
`+
				c.UseInCode("b", "c")))
}

// Java: ErrorToFixMapperTest.testFixRequires_veryLongNames()
func (c *ErrorToFixMapperTest) TestFixRequires_veryLongNames() {
	c.AssertNoChanges(
		c.FileWithImports(
			`const veryLongIdentifierSoLongThatItGoesPastThe80CharactersLimitAndYetWeShouldNotLineBreak = goog.require('b');
const {anotherVeryLongIdentifierSoLongThatItGoesPastThe80CharactersLimitAndYetWeShouldNotLineBreak} = goog.require('a');
` +
				c.UseInCode(
					"veryLongIdentifierSoLongThatItGoesPastThe80CharactersLimitAndYetWeShouldNotLineBreak",
					"anotherVeryLongIdentifierSoLongThatItGoesPastThe80CharactersLimitAndYetWeShouldNotLineBreak")))
}

// Java: ErrorToFixMapperTest.testFixRequires_veryLongNames_whenMergingDestructures()
func (c *ErrorToFixMapperTest) TestFixRequires_veryLongNames_whenMergingDestructures() {
	c.AssertChanges(
		c.FileWithImports(
			`const {veryLongSymbolThatMapsTo: veryLongLocalNameForIt} = goog.require('a');
const {anotherVeryLongSymbolThatMapsTo: veryLongLocalNameForItAlso} = goog.require('a');
`+
				c.UseInCode("veryLongLocalNameForIt", "veryLongLocalNameForItAlso")),
		c.FileWithImports(
			`const {anotherVeryLongSymbolThatMapsTo: veryLongLocalNameForItAlso, veryLongSymbolThatMapsTo: veryLongLocalNameForIt} = goog.require('a');
`+
				c.UseInCode("veryLongLocalNameForIt", "veryLongLocalNameForItAlso")))
}

// Java: ErrorToFixMapperTest.testFixRequires_noRequires()
func (c *ErrorToFixMapperTest) TestFixRequires_noRequires() {
	c.AssertNoChanges(c.FileWithImports(""))
}

// Java: ErrorToFixMapperTest.testMissingRequire_inJSDoc_withWhitespace()
func (c *ErrorToFixMapperTest) TestMissingRequire_inJSDoc_withWhitespace() {
	c.preexistingCode = "goog.provide('some.really.very.long.namespace.SuperInt');"
	c.AssertExpectedFixes(
		`goog.module('m');

/** @interface @implements {some.really.very.long.
                            namespace.SuperInt} */
class Bar {}
`,
		ErrorToFixMapperTest_ExpectedFix_Builder_().
			FixedCode(
				`goog.module('m');
const SuperInt = goog.require('some.really.very.long.namespace.SuperInt');

/** @interface @implements {SuperInt} */
class Bar {}
`).
			Build())
}

// Java: ErrorToFixMapperTest.testMissingRequire_unsorted1()
func (c *ErrorToFixMapperTest) TestMissingRequire_unsorted1() {
	// Both the fix for requires being unsorted, and the fix for the missing require, are applied.
	// However, the end result is still out of order.
	c.preexistingCode = "goog.provide('goog.dom.DomHelper');"
	c.AssertChanges(
		`goog.module('module');

const Xray = goog.require('goog.Xray');
const Anteater = goog.require('goog.Anteater');

alert(new Anteater());
alert(new Xray());
alert(new goog.dom.DomHelper());
`,
		`goog.module('module');

const DomHelper = goog.require('goog.dom.DomHelper');
const Anteater = goog.require('goog.Anteater');
const Xray = goog.require('goog.Xray');

alert(new Anteater());
alert(new Xray());
alert(new DomHelper());
`)
}

// Java: ErrorToFixMapperTest.testMissingRequire_unsorted2()
func (c *ErrorToFixMapperTest) TestMissingRequire_unsorted2() {
	// Both the fix for requires being unsorted, and the fix for the missing require, are applied.
	// The end result is ordered.
	c.preexistingCode = "goog.provide('goog.rays.Xray');"
	c.AssertChanges(
		`goog.module('module');

const DomHelper = goog.require('goog.dom.DomHelper');
const Anteater = goog.require('goog.Anteater');

alert(new Anteater());
alert(new goog.rays.Xray());
alert(new DomHelper());
`,
		`goog.module('module');

const Anteater = goog.require('goog.Anteater');
const DomHelper = goog.require('goog.dom.DomHelper');
const Xray = goog.require('goog.rays.Xray');

alert(new Anteater());
alert(new Xray());
alert(new DomHelper());
`)
}

// Java: ErrorToFixMapperTest.testMissingRequireInGoogModule()
func (c *ErrorToFixMapperTest) TestMissingRequireInGoogModule() {
	c.preexistingCode = "goog.provide('a.b.C');"
	c.AssertChanges(
		`goog.module('m');

alert(new a.b.C());
`,
		`goog.module('m');
const C = goog.require('a.b.C');

alert(new C());
`)
}

// Java: ErrorToFixMapperTest.testMissingRequireInGoogModuleTwice()
func (c *ErrorToFixMapperTest) TestMissingRequireInGoogModuleTwice() {
	c.preexistingCode = "goog.provide('a.b.C');"
	c.AssertChanges(
		`goog.module('m');

alert(new a.b.C());
alert(new a.b.C());
`,
		`goog.module('m');
const C = goog.require('a.b.C');

alert(new C());
alert(new C());
`)
}

// Java: ErrorToFixMapperTest.testMissingRequireInGoogModule_call()
func (c *ErrorToFixMapperTest) TestMissingRequireInGoogModule_call() {
	c.preexistingCode = "goog.provide('a.b');"
	c.AssertChanges(
		`goog.module('m');

alert(a.b.c());
`,
		`goog.module('m');
const b = goog.require('a.b');

alert(b.c());
`)
}

// Java: ErrorToFixMapperTest.testMissingRequireInGoogModule_extends()
func (c *ErrorToFixMapperTest) TestMissingRequireInGoogModule_extends() {
	c.preexistingCode = "goog.provide('world.util.Animal');"
	c.AssertChanges(
		`goog.module('m');

class Cat extends world.util.Animal {}
`,
		`goog.module('m');
const Animal = goog.require('world.util.Animal');

class Cat extends Animal {}
`)
}

// Java: ErrorToFixMapperTest.testMissingRequireInGoogModule_atExtends()
func (c *ErrorToFixMapperTest) TestMissingRequireInGoogModule_atExtends() {
	c.preexistingCode = "goog.provide('world.util.Animal');"
	c.AssertChanges(
		`goog.module('m');

/** @constructor @extends {world.util.Animal} */
function Cat() {}
`,
		`goog.module('m');
const Animal = goog.require('world.util.Animal');

/** @constructor @extends {Animal} */
function Cat() {}
`)
}

// Java: ErrorToFixMapperTest.testMissingRequire_inProvidesFile()
func (c *ErrorToFixMapperTest) TestMissingRequire_inProvidesFile() {
	c.preexistingCode = "goog.provide('world.util.Animal');"
	c.AssertChanges(
		`var f = world.util.Animal.create();
`,
		`goog.require('world.util.Animal');
var f = world.util.Animal.create();
`)
}

// Java: ErrorToFixMapperTest.testStandaloneVarDoesntCrashMissingRequire()
func (c *ErrorToFixMapperTest) TestStandaloneVarDoesntCrashMissingRequire() {
	c.preexistingCode = "goog.provide('goog.Animal');"
	c.AssertChanges(
		`goog.module('m');

var x;

class Cat extends goog.Animal {}
`,
		`goog.module('m');
const Animal = goog.require('goog.Animal');

var x;

class Cat extends Animal {}
`)
}

// Java: ErrorToFixMapperTest.testAddLhsToGoogRequire()
func (c *ErrorToFixMapperTest) TestAddLhsToGoogRequire() {
	c.preexistingCode = "goog.provide('world.util.Animal');"
	c.AssertChanges(
		`goog.module('m');

goog.require('world.util.Animal');

class Cat extends world.util.Animal {}
`,
		`goog.module('m');

const Animal = goog.require('world.util.Animal');

class Cat extends Animal {}
`)
}

// Java: ErrorToFixMapperTest.testAddLhsToGoogRequire_conflictingName()
func (c *ErrorToFixMapperTest) TestAddLhsToGoogRequire_conflictingName() {
	c.preexistingCode = "goog.provide('world.util.Animal');"
	c.AssertChanges(
		`goog.module('m');

goog.require('world.util.Animal');

const Animal = null;

class Cat extends world.util.Animal {}
`,
		`goog.module('m');

const UtilAnimal = goog.require('world.util.Animal');

const Animal = null;

class Cat extends UtilAnimal {}
`)
}

// Java: ErrorToFixMapperTest.testAddLhsToGoogRequire_conflictingName_fromOtherSuggestion()
func (c *ErrorToFixMapperTest) TestAddLhsToGoogRequire_conflictingName_fromOtherSuggestion() {
	c.preexistingCode = "goog.provide('world.util.Animal'); goog.provide('rara.exotic.Animal');"
	c.AssertChanges(
		`goog.module('m');

goog.require('rare.exotic.Animal');
goog.require('world.util.Animal');

/** @implements {rare.exotic.Animal} */
class Cat extends world.util.Animal {}
`,
		`goog.module('m');

const ExoticAnimal = goog.require('rare.exotic.Animal');
const Animal = goog.require('world.util.Animal');

/** @implements {ExoticAnimal} */
class Cat extends Animal {}
`)
}

// Java: ErrorToFixMapperTest.testAddLhsToGoogRequire_new()
func (c *ErrorToFixMapperTest) TestAddLhsToGoogRequire_new() {
	c.preexistingCode = "goog.provide('world.util.Animal');"
	c.AssertChanges(
		`goog.module('m');

goog.require('world.util.Animal');

let cat = new world.util.Animal();
`,
		`goog.module('m');

const Animal = goog.require('world.util.Animal');

let cat = new Animal();
`)
}

// Java: ErrorToFixMapperTest.testAddLhsToGoogRequire_getprop()
func (c *ErrorToFixMapperTest) TestAddLhsToGoogRequire_getprop() {
	c.preexistingCode = "goog.provide('magical.factories'); goog.provide('world.util.Animal');"
	c.AssertChanges(
		`goog.module('m');

goog.require('magical.factories');
goog.require('world.util.AnimalType');

let cat = magical.factories.createAnimal(world.util.AnimalType.CAT);
`,
		`goog.module('m');

const factories = goog.require('magical.factories');
const AnimalType = goog.require('world.util.AnimalType');

let cat = factories.createAnimal(AnimalType.CAT);
`)
}

// Java: ErrorToFixMapperTest.testAddLhsToGoogRequire_jsdoc()
func (c *ErrorToFixMapperTest) TestAddLhsToGoogRequire_jsdoc() {
	c.preexistingCode = "goog.provide('world.util.Animal');"
	c.AssertChanges(
		`goog.module('m');

goog.require('world.util.Animal');

/** @type {!world.util.Animal} */
var cat;
`,
		`goog.module('m');

const Animal = goog.require('world.util.Animal');

/** @type {!Animal} */
var cat;
`)
}

// Java: ErrorToFixMapperTest.testSwitchToShorthand_JSDoc1()
func (c *ErrorToFixMapperTest) TestSwitchToShorthand_JSDoc1() {
	c.preexistingCode = "goog.provide('world.util.Animal');"
	c.AssertChanges(
		`goog.module('m');
var Animal = goog.require('world.util.Animal');

/** @constructor @implements {world.util.Animal} */
function Cat() {}
`,
		`goog.module('m');
var Animal = goog.require('world.util.Animal');

/** @constructor @implements {Animal} */
function Cat() {}
`)
}

// Java: ErrorToFixMapperTest.testSwitchToShorthand_JSDoc2()
func (c *ErrorToFixMapperTest) TestSwitchToShorthand_JSDoc2() {
	c.preexistingCode = "goog.provide('world.util.Animal');"
	c.AssertChanges(
		`goog.module('m');
var Animal = goog.require('world.util.Animal');

/** @constructor @extends {world.util.Animal} */
function Cat() {}
`,
		`goog.module('m');
var Animal = goog.require('world.util.Animal');

/** @constructor @extends {Animal} */
function Cat() {}
`)
}

// Java: ErrorToFixMapperTest.testSwitchToShorthand_JSDoc3()
func (c *ErrorToFixMapperTest) TestSwitchToShorthand_JSDoc3() {
	c.preexistingCode = "goog.provide('world.util.Animal');"
	c.AssertChanges(
		`goog.module('m');
var Animal = goog.require('world.util.Animal');

/** @type {!world.util.Animal} */
var animal;
`,
		`goog.module('m');
var Animal = goog.require('world.util.Animal');

/** @type {!Animal} */
var animal;
`)
}

// Java: ErrorToFixMapperTest.testSwitchToShorthand_JSDoc4()
func (c *ErrorToFixMapperTest) TestSwitchToShorthand_JSDoc4() {
	c.preexistingCode = "goog.provide('world.util.Animal');"
	c.AssertChanges(
		`goog.module('m');
var Animal = goog.require('world.util.Animal');

/** @type {?world.util.Animal} */
var animal;
`,
		`goog.module('m');
var Animal = goog.require('world.util.Animal');

/** @type {?Animal} */
var animal;
`)
}

// Java: ErrorToFixMapperTest.testSwitchToShorthand_JSDoc5()
func (c *ErrorToFixMapperTest) TestSwitchToShorthand_JSDoc5() {
	c.preexistingCode = "goog.provide('world.util.Animal');"
	c.AssertChanges(
		`goog.module('m');
var Animal = goog.require('world.util.Animal');

/** @type {?Array<!world.util.Animal>} */
var animals;
`,
		`goog.module('m');
var Animal = goog.require('world.util.Animal');

/** @type {?Array<!Animal>} */
var animals;
`)
}

// Java: ErrorToFixMapperTest.testSwitchToShorthand_JSDoc6()
func (c *ErrorToFixMapperTest) TestSwitchToShorthand_JSDoc6() {
	c.preexistingCode = "goog.provide('world.util.Animal');"
	c.AssertChanges(
		`goog.module('m');
var Animal = goog.require('world.util.Animal');

/** @type {?Array<!world.util.Animal.Turtle>} */
var turtles;
`,
		`goog.module('m');
var Animal = goog.require('world.util.Animal');

/** @type {?Array<!Animal.Turtle>} */
var turtles;
`)
}

// Java: ErrorToFixMapperTest.testSwitchToShorthand_JSDoc7()
func (c *ErrorToFixMapperTest) TestSwitchToShorthand_JSDoc7() {
	c.preexistingCode = "goog.provide('world.util.Animal');"
	c.AssertChanges(
		`goog.module('m');
var AnimalAltName = goog.require('world.util.Animal');

/** @type {?Array<!world.util.Animal.Turtle>} */
var turtles;
`,
		`goog.module('m');
var AnimalAltName = goog.require('world.util.Animal');

/** @type {?Array<!AnimalAltName.Turtle>} */
var turtles;
`)
}

// Java: ErrorToFixMapperTest.testMissingRequireInGoogModule_atExtends_qname()
func (c *ErrorToFixMapperTest) TestMissingRequireInGoogModule_atExtends_qname() {
	c.preexistingCode = "goog.provide('world.util.Animal');"
	c.AssertChanges(
		`goog.module('m');

/** @constructor @extends {world.util.Animal} */
world.util.Cat = function() {};
`,
		`goog.module('m');
const Animal = goog.require('world.util.Animal');

/** @constructor @extends {Animal} */
world.util.Cat = function() {};
`)
}

// Java: ErrorToFixMapperTest.testMissingRequireTypeInGoogModule_atType_qname()
func (c *ErrorToFixMapperTest) TestMissingRequireTypeInGoogModule_atType_qname() {
	c.preexistingCode = "goog.provide('world.util.Animal');"
	c.AssertChanges(
		`goog.module('m');

/** @type {?world.util.Animal} */
let x = null;
`,
		`goog.module('m');
const Animal = goog.requireType('world.util.Animal');

/** @type {?Animal} */
let x = null;
`)
}

// Java: ErrorToFixMapperTest.testMissingRequireInGoogModule_googString()
func (c *ErrorToFixMapperTest) TestMissingRequireInGoogModule_googString() {
	c.preexistingCode = "goog.provide('goog.string');"
	c.AssertChanges(
		`goog.module('m');

alert(goog.string.trim('   str    '));
`,
		`goog.module('m');
const googString = goog.require('goog.string');

alert(googString.trim('   str    '));
`)
}

// Java: ErrorToFixMapperTest.testMissingRequireInGoogModule_googStructsMap()
func (c *ErrorToFixMapperTest) TestMissingRequireInGoogModule_googStructsMap() {
	c.preexistingCode = "goog.provide('goog.structs.Map');"
	c.AssertChanges(
		`goog.module('m');

alert(new goog.structs.Map());
`,
		`goog.module('m');
const StructsMap = goog.require('goog.structs.Map');

alert(new StructsMap());
`)
}

// Java: ErrorToFixMapperTest.testMissingRequireInGoogModule_insertedInCorrectOrder()
func (c *ErrorToFixMapperTest) TestMissingRequireInGoogModule_insertedInCorrectOrder() {
	c.preexistingCode = "goog.provide('x.B');"
	c.AssertChangesWithExpectedFixes(
		`goog.module('m');

const A = goog.require('a.A');
const C = goog.require('c.C');

alert(new A(new x.B(new C())));
`,
		// Requires are sorted by the short name, not the full namespace.
		`goog.module('m');

const A = goog.require('a.A');
const B = goog.require('x.B');
const C = goog.require('c.C');

alert(new A(new B(new C())));
`)
}

// Java: ErrorToFixMapperTest.testMissingRequireInGoogModule_alwaysInsertsConst()
func (c *ErrorToFixMapperTest) TestMissingRequireInGoogModule_alwaysInsertsConst() {
	c.preexistingCode = "goog.provide('x.B');"
	c.AssertChanges(
		`goog.module('m');

var A = goog.require('a.A');
var C = goog.require('c.C');

alert(new A(new x.B(new C())));
`,
		`goog.module('m');

var A = goog.require('a.A');
const B = goog.require('x.B');
var C = goog.require('c.C');

alert(new A(new B(new C())));
`)
}

// Java: ErrorToFixMapperTest.testGoogForwardDeclare_fixedToRequire()
func (c *ErrorToFixMapperTest) TestGoogForwardDeclare_fixedToRequire() {
	c.preexistingCode = "goog.provide('goog.dom.DomHelper');"
	c.AssertChanges(
		`goog.module('module');

goog.forwardDeclare('goog.dom.DomHelper');

new goog.dom.DomHelper();
`,
		`goog.module('module');

const DomHelper = goog.require('goog.dom.DomHelper');

new DomHelper();
`)
}

// Java: ErrorToFixMapperTest.testGoogForwardDeclare_fixedToRequireType()
func (c *ErrorToFixMapperTest) TestGoogForwardDeclare_fixedToRequireType() {
	c.preexistingCode = "goog.provide('goog.dom.DomHelper');"
	c.AssertChanges(
		`goog.module('module');

goog.forwardDeclare('goog.dom.DomHelper');

function f(/** !goog.dom.DomHelper */ x) {}
`,
		`goog.module('module');

const DomHelper = goog.requireType('goog.dom.DomHelper');

function f(/** !DomHelper */ x) {}
`)
}

// Java: ErrorToFixMapperTest.testShortRequireInGoogModule1()
func (c *ErrorToFixMapperTest) TestShortRequireInGoogModule1() {
	c.AssertChanges(
		`goog.module('m');

var c = goog.require('a.b.c');

alert(a.b.c);
`,
		`goog.module('m');

var c = goog.require('a.b.c');

alert(c);
`)
}

// Java: ErrorToFixMapperTest.testShortRequireInGoogModule2()
func (c *ErrorToFixMapperTest) TestShortRequireInGoogModule2() {
	c.AssertChanges(
		`goog.module('m');

var Classname = goog.require('a.b.Classname');

alert(a.b.Classname.instance_.foo());
`,
		`goog.module('m');

var Classname = goog.require('a.b.Classname');

alert(Classname.instance_.foo());
`)
}

// Java: ErrorToFixMapperTest.testShortRequireInGoogModule3()
func (c *ErrorToFixMapperTest) TestShortRequireInGoogModule3() {
	c.AssertChanges(
		`goog.module('m');

var Classname = goog.require('a.b.Classname');

alert(a.b.Classname.INSTANCE_.foo());
`,
		`goog.module('m');

var Classname = goog.require('a.b.Classname');

alert(Classname.INSTANCE_.foo());
`)
}

// Java: ErrorToFixMapperTest.testShortRequireInGoogModule4()
func (c *ErrorToFixMapperTest) TestShortRequireInGoogModule4() {
	c.AssertChanges(
		`goog.module('m');

var object = goog.require('goog.object');

alert(goog.object.values({x:1}));
`,
		`goog.module('m');

var object = goog.require('goog.object');

alert(object.values({x:1}));
`)
}

// Java: ErrorToFixMapperTest.testShortRequireInGoogModule5()
func (c *ErrorToFixMapperTest) TestShortRequireInGoogModule5() {
	c.AssertChanges(
		`goog.module('m');

var Widget = goog.require('goog.Widget');

alert(new goog.Widget());
`,
		`goog.module('m');

var Widget = goog.require('goog.Widget');

alert(new Widget());
`)
}

// Java: ErrorToFixMapperTest.testShortRequireInGoogModule6()
func (c *ErrorToFixMapperTest) TestShortRequireInGoogModule6() {
	c.AssertChanges(
		`goog.module('m');

var GoogWidget = goog.require('goog.Widget');

alert(new goog.Widget());
`,
		`goog.module('m');

var GoogWidget = goog.require('goog.Widget');

alert(new GoogWidget());
`)
}

/**
* Here, if the short name weren't provided the suggested fix would use 'Table' for both, but
* since there is a short name provided for each one, it uses those names.
 */
// Java: ErrorToFixMapperTest.testShortRequireInGoogModule7()
func (c *ErrorToFixMapperTest) TestShortRequireInGoogModule7() {
	c.AssertChanges(
		`goog.module('m');

var CoffeeTable = goog.require('coffee.Table');
var KitchenTable = goog.require('kitchen.Table');

alert(new coffee.Table(), new kitchen.Table());
`,
		`goog.module('m');

var CoffeeTable = goog.require('coffee.Table');
var KitchenTable = goog.require('kitchen.Table');

alert(new CoffeeTable(), new KitchenTable());
`)
}

// Java: ErrorToFixMapperTest.testBug65602711a()
func (c *ErrorToFixMapperTest) TestBug65602711a() {
	c.AssertChanges(
		`goog.module('x');

const {X} = goog.require('ns.abc.xyz');

use(ns.abc.xyz.X);
`,
		`goog.module('x');

const {X} = goog.require('ns.abc.xyz');

use(X);
`)
}

// Java: ErrorToFixMapperTest.testBug65602711b()
func (c *ErrorToFixMapperTest) TestBug65602711b() {
	c.AssertChanges(
		`goog.module('x');

const {X: X2} = goog.require('ns.abc.xyz');

use(ns.abc.xyz.X);
`,
		`goog.module('x');

const {X: X2} = goog.require('ns.abc.xyz');

use(X2);
`)
}

// Java: ErrorToFixMapperTest.testExtraRequire()
func (c *ErrorToFixMapperTest) TestExtraRequire() {
	c.AssertExpectedFixes(
		`goog.require('goog.object');
goog.require('goog.string');

alert(goog.string.parseInt('7'));
`,
		ErrorToFixMapperTest_ExpectedFix_Builder_().
			Description("Delete extra require").
			FixedCode(
				`goog.require('goog.string');

alert(goog.string.parseInt('7'));
`).
			Build())
}

// Java: ErrorToFixMapperTest.testExtraRequireType()
func (c *ErrorToFixMapperTest) TestExtraRequireType() {
	c.AssertChanges(
		`goog.requireType('goog.events.Listenable');
goog.require('goog.string');

alert(goog.string.parseInt('7'));
`,
		`goog.require('goog.string');

alert(goog.string.parseInt('7'));
`)
}

// Java: ErrorToFixMapperTest.testExtraRequire_module()
func (c *ErrorToFixMapperTest) TestExtraRequire_module() {
	c.AssertChanges(
		`goog.module('x');

const googString = goog.require('goog.string');
const object = goog.require('goog.object');
alert(googString.parseInt('7'));
`,
		`goog.module('x');

const googString = goog.require('goog.string');
alert(googString.parseInt('7'));
`)
}

// Java: ErrorToFixMapperTest.testExtraRequireType_module()
func (c *ErrorToFixMapperTest) TestExtraRequireType_module() {
	c.AssertChanges(
		`goog.module('x');

const googString = goog.require('goog.string');
const Listenable = goog.requireType('goog.events.Listenable');
alert(googString.parseInt('7'));
`,
		`goog.module('x');

const googString = goog.require('goog.string');
alert(googString.parseInt('7'));
`)
}

// Java: ErrorToFixMapperTest.testExtraRequire_destructuring_unusedInitialMember()
func (c *ErrorToFixMapperTest) TestExtraRequire_destructuring_unusedInitialMember() {
	c.AssertChanges(
		`goog.module('x');

const {foo, bar} = goog.require('goog.util');

alert(bar(7));
`,
		`goog.module('x');

const {bar} = goog.require('goog.util');

alert(bar(7));
`)
}

// Java: ErrorToFixMapperTest.testExtraRequire_destructuring_unusedFinalMember()
func (c *ErrorToFixMapperTest) TestExtraRequire_destructuring_unusedFinalMember() {
	c.AssertChanges(
		`goog.module('x');

const {foo, bar} = goog.require('goog.util');

alert(foo(7));
`,
		`goog.module('x');

const {foo} = goog.require('goog.util');

alert(foo(7));
`)
}

// Java: ErrorToFixMapperTest.testExtraRequire_destructuring_unusedMiddleMember()
func (c *ErrorToFixMapperTest) TestExtraRequire_destructuring_unusedMiddleMember() {
	c.AssertChanges(
		`goog.module('x');

const {foo, bar, qux} = goog.require('goog.util');

alert(foo(qux(7)));
`,
		`goog.module('x');

const {foo, qux} = goog.require('goog.util');

alert(foo(qux(7)));
`)
}

/** Because of overlapping replacements, it takes two runs to fully fix this case. */
// Java: ErrorToFixMapperTest.testExtraRequire_destructuring_multipleUnusedMembers()
func (c *ErrorToFixMapperTest) TestExtraRequire_destructuring_multipleUnusedMembers() {
	c.AssertChanges(
		`goog.module('x');

const {foo, bar, qux} = goog.require('goog.util');
`,
		`goog.module('x');

const {qux} = goog.require('goog.util');
`)
}

// Java: ErrorToFixMapperTest.testExtraRequire_destructuring_allUnusedMembers()
func (c *ErrorToFixMapperTest) TestExtraRequire_destructuring_allUnusedMembers() {
	c.AssertChanges(
		`goog.module('x');

const {qux} = goog.require('goog.util');
`,
		`goog.module('x');

const {} = goog.require('goog.util');
`)
}

// Java: ErrorToFixMapperTest.testExtraRequire_destructuring_unusedShortnameMember()
func (c *ErrorToFixMapperTest) TestExtraRequire_destructuring_unusedShortnameMember() {
	c.AssertExpectedFixes(
		`goog.module('x');

const {bar, foo: googUtilFoo} = goog.require('goog.util');

alert(bar(7));
`,
		ErrorToFixMapperTest_ExpectedFix_Builder_().
			Description("Delete unused symbol").
			FixedCode(
				`goog.module('x');

const {bar} = goog.require('goog.util');

alert(bar(7));
`).
			Build())
}

// Java: ErrorToFixMapperTest.testExtraRequire_destructuring_keepShortnameMember()
func (c *ErrorToFixMapperTest) TestExtraRequire_destructuring_keepShortnameMember() {
	c.AssertChanges(
		`goog.module('x');

const {foo: googUtilFoo, bar} = goog.require('goog.util');

alert(googUtilFoo(7));
`,
		`goog.module('x');

const {foo: googUtilFoo} = goog.require('goog.util');

alert(googUtilFoo(7));
`)
}

// Java: ErrorToFixMapperTest.testExtraRequire_destructuring_onlyUnusedShortnameMember()
func (c *ErrorToFixMapperTest) TestExtraRequire_destructuring_onlyUnusedShortnameMember() {
	c.AssertChanges(
		`goog.module('x');

const {foo: googUtilFoo} = goog.require('goog.util');
`,
		`goog.module('x');

const {} = goog.require('goog.util');
`)
}

// Java: ErrorToFixMapperTest.testExtraRequire_destructuring_noMembers()
func (c *ErrorToFixMapperTest) TestExtraRequire_destructuring_noMembers() {
	c.AssertChanges(
		`goog.module('x');

const {} = goog.require('goog.util');
`,
		`goog.module('x');


`)
}

// Java: ErrorToFixMapperTest.testExtraRequireType_destructuring()
func (c *ErrorToFixMapperTest) TestExtraRequireType_destructuring() {
	c.AssertChanges(
		`goog.module('x');

const {foo: googUtilFoo, bar, baz: googUtilBaz, qux} = goog.requireType('goog.util');

/** @type {!googUtilFoo} */ let x;
/** @type {!qux} */ let x;
`,
		`goog.module('x');

const {foo: googUtilFoo, qux} = goog.requireType('goog.util');

/** @type {!googUtilFoo} */ let x;
/** @type {!qux} */ let x;
`)
}

// Java: ErrorToFixMapperTest.testExtraRequire_unsorted()
func (c *ErrorToFixMapperTest) TestExtraRequire_unsorted() {
	// There is also a warning because requires are not sorted. That one is not fixed because
	// the fix would conflict with the extra-require fix.
	c.AssertChanges(
		`goog.require('goog.string');
goog.require('goog.object');
goog.require('goog.dom');

alert(goog.string.parseInt('7'));
alert(goog.dom.createElement('div'));
`,
		`goog.require('goog.string');
goog.require('goog.dom');

alert(goog.string.parseInt('7'));
alert(goog.dom.createElement('div'));
`)
}

// TODO(tjgq): Make this not crash on ClosureRewriteModule#updateGoogRequire.
// Java: ErrorToFixMapperTest.testNoCrashOnInvalidMultiRequireStatement()
func (c *ErrorToFixMapperTest) TestNoCrashOnInvalidMultiRequireStatement() {
	c.t.Skip("@Ignore")
	c.AssertNoChanges(
		c.FileWithImports(
			`const a = goog.require('a'), b = goog.require('b');
` +
				c.UseInCode("a", "b")))
}

// Java: ErrorToFixMapperTest.testConstantCaseName_let()
func (c *ErrorToFixMapperTest) TestConstantCaseName_let() {
	c.AssertExpectedFixes(
		"goog.module('m'); let CONSTANT_CASE = 'value';",
		ErrorToFixMapperTest_ExpectedFix_Builder_().
			Description("Make explicitly constant").
			FixedCode("goog.module('m'); const CONSTANT_CASE = 'value';").
			Build())
}

// Java: ErrorToFixMapperTest.testConstantCaseName_var()
func (c *ErrorToFixMapperTest) TestConstantCaseName_var() {
	c.AssertExpectedFixes(
		`goog.module('m');
var CONSTANT_CASE = 'value';
`,
		ErrorToFixMapperTest_ExpectedFix_Builder_().
			Description("Make explicitly constant").
			FixedCode(
				`goog.module('m');
/** @const */
var CONSTANT_CASE = 'value';
`).
			Build())
}

// Java: ErrorToFixMapperTest.testConstantCaseName_varWithExistingJSDoc()
func (c *ErrorToFixMapperTest) TestConstantCaseName_varWithExistingJSDoc() {
	c.AssertNoChanges(
		"goog.module('m'); /** @type {string} Some description */ var CONSTANT_CASE = 'value';")
}

// Java: ErrorToFixMapperTest.fileWithImports(String)
func (c *ErrorToFixMapperTest) FileWithImports(imports string) string {
	return "/*\n" +
		" * @fileoverview\n" +
		" */\n" +
		"goog.module('x');\n" +
		imports +
		"module.exports = function() {};\n"
}

// Java: ErrorToFixMapperTest.useInCode(String...)
func (c *ErrorToFixMapperTest) UseInCode(names ...string) string {
	var sb strings.Builder
	sb.WriteString("\n")
	for _, name := range names {
		sb.WriteString("use(")
		sb.WriteString(name)
		sb.WriteString(");\n")
	}
	return sb.String()
}

// Java: ErrorToFixMapperTest.useInType(String...)
func (c *ErrorToFixMapperTest) UseInType(names ...string) string {
	var sb strings.Builder
	sb.WriteString("\n")
	for _, name := range names {
		sb.WriteString("/** @type {!")
		sb.WriteString(name)
		sb.WriteString("} */ var _")
		sb.WriteString(name)
		sb.WriteString(";\n")
	}
	return sb.String()
}

// Java: ErrorToFixMapperTest.compileExpectingAtLeastOneWarningOrError(String)
func (c *ErrorToFixMapperTest) CompileExpectingAtLeastOneWarningOrError(originalCode string) {
	c.CompileExpectingAtLeastOneWarningOrErrorOfMap(c.GetInputMap(originalCode))
}

// Java: ErrorToFixMapperTest.compileExpectingAtLeastOneWarningOrError(ImmutableMap<String, String>)
func (c *ErrorToFixMapperTest) CompileExpectingAtLeastOneWarningOrErrorOfMap(inputMap *jutil.LinkedHashMap[string, string]) {
	inputs := []*SourceFile{}
	for k, v := range inputMap.All() {
		inputs = append(inputs, SourceFile_FromCode(k, v))
	}
	c.compiler.Compile(
		[]*SourceFile{}, // Externs
		inputs,
		c.options)
	warningsAndErrors := []*JSError{}
	warningsAndErrors = append(warningsAndErrors, c.compiler.GetWarnings()...)
	warningsAndErrors = append(warningsAndErrors, c.compiler.GetErrors()...)
	AssertWithMessage("warnings/errors").That(warningsAndErrors).IsNotEmpty()
}

// AssertChanges applies all the fixes together.
//
// Deprecated: Use AssertExpectedFixes() for new cases or when amending
// existing ones.
//
// Java: ErrorToFixMapperTest.assertChanges(String, String)
func (c *ErrorToFixMapperTest) AssertChanges(originalCode string, expectedCode string) {
	c.CompileExpectingAtLeastOneWarningOrError(originalCode)
	fixes := c.errorManager.GetAllFixes()
	AssertWithMessage("fixes").That(fixes).IsNotEmpty()
	newCode :=
		errorToFixMapperTest_get(ApplySuggestedFixes_ApplySuggestedFixesToCode(slices.Values(fixes), c.GetInputMap(originalCode)), "test")
	AssertThat(newCode).IsEqualTo(expectedCode)
}

// AssertChangesWithExpectedFixes applies each fix on its own.
//
// Deprecated: Use AssertExpectedFixes() for new cases or when amending
// existing ones.
//
// Java: ErrorToFixMapperTest.assertChanges(String, String...)
func (c *ErrorToFixMapperTest) AssertChangesWithExpectedFixes(originalCode string, expectedFixes ...string) {
	c.CompileExpectingAtLeastOneWarningOrError(originalCode)
	fixes := c.errorManager.GetAllFixes()
	AssertWithMessage("fixes").That(fixes).HasSize(len(expectedFixes))
	for i := 0; i < len(fixes); i++ {
		newCode :=
			errorToFixMapperTest_get(ApplySuggestedFixes_ApplySuggestedFixesToCode(
				slices.Values([]*SuggestedFix{fixes[i]}), c.GetInputMap(originalCode)),
				"test")
		AssertThat(newCode).IsEqualTo(expectedFixes[i])
	}
}

// Java: ErrorToFixMapperTest.assertExpectedFixes(String, ExpectedFix...)
func (c *ErrorToFixMapperTest) AssertExpectedFixes(originalCode string, expectedFixes ...*ErrorToFixMapperTest_ExpectedFix) {
	c.CompileExpectingAtLeastOneWarningOrError(originalCode)
	fixes := c.errorManager.GetAllFixes()
	AssertWithMessage("Unexpected number of fixes").That(fixes).HasSize(len(expectedFixes))
	for i := 0; i < len(fixes); i++ {
		expectedFix := expectedFixes[i]
		actualFix := fixes[i]
		AssertWithMessage("Actual fix[%s]: %s", i, actualFix).
			That(errorToFixMapperTest_nullable(actualFix.GetDescription())).
			IsEqualTo(errorToFixMapperTest_nullable(expectedFix.Description()))
		AssertThat(
			ApplySuggestedFixes_ApplySuggestedFixesToCode(
				slices.Values([]*SuggestedFix{fixes[i]}), c.GetInputMap(originalCode))).
			ContainsEntry("test", expectedFix.FixedCode())
	}
}

// Java: ErrorToFixMapperTest.assertNoChanges(String)
func (c *ErrorToFixMapperTest) AssertNoChanges(originalCode string) {
	c.compiler.Compile(
		[]*SourceFile{}, // Externs
		[]*SourceFile{SourceFile_FromCode("test", originalCode)},
		c.options)
	fixes := c.errorManager.GetAllFixes()
	AssertThat(fixes).IsEmpty()
}

// Java: ErrorToFixMapperTest.getInputMap(String)
func (c *ErrorToFixMapperTest) GetInputMap(originalCode string) *jutil.LinkedHashMap[string, string] {
	return jutil.NewImmutableMapBuilder[string, string]().
		Put("preexistingCode", c.preexistingCode).
		Put("test", originalCode).
		BuildOrThrow()
}

// errorToFixMapperTest_get is ImmutableMap.get(key): nil when absent.
//
// Go-only: test-helper.
func errorToFixMapperTest_get(m *jutil.LinkedHashMap[string, string], key string) any {
	if v, ok := m.Lookup(key); ok {
		return v
	}
	return nil
}

// errorToFixMapperTest_nullable is a nullable String as Truth compares it.
//
// Go-only: test-helper.
func errorToFixMapperTest_nullable(s *string) any {
	if s == nil {
		return nil
	}
	return *s
}
