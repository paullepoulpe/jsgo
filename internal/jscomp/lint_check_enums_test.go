// Port of com/google/javascript/jscomp/lint/CheckEnumsTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/lint/CheckEnumsTest.java)

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
	"testing"
)

// Test case for {@link CheckEnums}.
type CheckEnumsTest struct {
	CompilerTestCaseBase
}

func NewCheckEnumsTest() *CheckEnumsTest {
	c := &CheckEnumsTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestCheckEnumsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCheckEnumsTest() })
}

// Java: CheckEnumsTest.getProcessor(Compiler)
func (c *CheckEnumsTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewCheckEnums(compiler)
}

// Java: CheckEnumsTest.getOptions()
func (c *CheckEnumsTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetWarningLevel(DiagnosticGroups_LINT_CHECKS, CheckLevel_WARNING)
	return options
}

// Java: CheckEnumsTest.testCheckEnums()
func (c *CheckEnumsTest) TestCheckEnums() {
	c.TestSame("/** @enum {number} */ ns.Enum = {A: 1, B: 2};")
	c.TestSame("/** @enum {string} */ ns.Enum = {A: 'foo', B: 'bar'};")
	c.TestWarning("/** @enum {number} */ ns.Enum = {A: 1, B: 1};", CheckEnums_DUPLICATE_ENUM_VALUE)
	c.TestWarning("/** @enum {string} */ ns.Enum = {A: 'foo', B: 'foo'};", CheckEnums_DUPLICATE_ENUM_VALUE)
	c.TestSame("/** @enum {number} */ let Enum = {A: 1, B: 2};")
	c.TestSame("/** @enum {string} */ let Enum = {A: 'foo', B: 'bar'};")
	c.TestWarning("/** @enum {number} */ let Enum = {A: 1, B: 1};", CheckEnums_DUPLICATE_ENUM_VALUE)
	c.TestNoWarning("/** @enum {number} */ let Enum = {A: 1+1, B: 1+1};")
	// uncaught dup values
	c.TestWarning("/** @enum {string} */ let Enum = {A: 'foo', B: 'foo'};", CheckEnums_DUPLICATE_ENUM_VALUE)
	c.TestNoWarning(
		"/** @enum {number} */ let Enum = {A: 'a'+'b', B: 'a'+'b'};")
	// uncaught dup values
	c.TestWarning("/** @enum {number} */ let Enum = {A};", CheckEnums_SHORTHAND_ASSIGNMENT_IN_ENUM)
	c.TestWarning(
		"/** @enum {string} */ let Enum = {['prop' + f()]: 'foo'};", CheckEnums_COMPUTED_PROP_NAME_IN_ENUM)
	c.TestWarning("/** @enum {number} */ let E = { a: 1 };", CheckEnums_ENUM_PROP_NOT_CONSTANT)
	c.TestWarning("/** @enum {number} */ let E = { ABC: 1, abc: 2 };", CheckEnums_ENUM_PROP_NOT_CONSTANT)
}

// Java: CheckEnumsTest.testCheckValidEnums_withES6Modules()
func (c *CheckEnumsTest) TestCheckValidEnums_withES6Modules() {
	c.TestSame("export /** @enum {number} */ let Enum = {A: 1, B: 2};")
}

// Java: CheckEnumsTest.testCheckInvalidEnums_withES6Modules01()
func (c *CheckEnumsTest) TestCheckInvalidEnums_withES6Modules01() {
	c.TestWarning("export /** @enum {number} */ let Enum = {A: 1, B: 1};", CheckEnums_DUPLICATE_ENUM_VALUE)
}

// Java: CheckEnumsTest.testCheckInvalidEnums_withES6Modules02()
func (c *CheckEnumsTest) TestCheckInvalidEnums_withES6Modules02() {
	c.TestWarning("export /** @enum {number} */ let Enum = {A};", CheckEnums_SHORTHAND_ASSIGNMENT_IN_ENUM)
}

// Java: CheckEnumsTest.testCheckInvalidEnums_withES6Modules03()
func (c *CheckEnumsTest) TestCheckInvalidEnums_withES6Modules03() {
	c.TestWarning(
		"export /** @enum {string} */ let Enum = {['prop' + f()]: 'foo'};",
		CheckEnums_COMPUTED_PROP_NAME_IN_ENUM)
}

// Java: CheckEnumsTest.testCheckInvalidEnums_withES6Modules04()
func (c *CheckEnumsTest) TestCheckInvalidEnums_withES6Modules04() {
	c.TestWarning("export /** @enum {number} */ let E = { a: 1 };", CheckEnums_ENUM_PROP_NOT_CONSTANT)
}

// Java: CheckEnumsTest.testEnumTypeIsDeclaredStringOrNumber()
func (c *CheckEnumsTest) TestEnumTypeIsDeclaredStringOrNumber() {
	c.TestNoWarning("/** @enum {number} */ let E = { A: 1 };")
	c.TestNoWarning("/** @enum {number} */ let E = { A: 1+1 };")
	c.TestNoWarning("/** @enum {string} */ let E = { A: 'static str' };")
	c.TestNoWarning("/** @enum {string} */ let E = { A: `template lit` };")
	c.TestWarning("/** @enum {?number} */ let E = { A: 1 };", CheckEnums_ENUM_TYPE_NOT_STRING_OR_NUMBER)
	c.TestWarning(
		"/** @enum {!string} */ let E = { A: 'static str' };", CheckEnums_ENUM_TYPE_NOT_STRING_OR_NUMBER)
	c.TestWarning("/** @enum {boolean} */ let E = {A: true };", CheckEnums_ENUM_TYPE_NOT_STRING_OR_NUMBER)
	c.TestWarning(
		"function foo() {} /** @enum {Function} */ let E = {A: foo };",
		CheckEnums_ENUM_TYPE_NOT_STRING_OR_NUMBER)
	c.TestWarning(
		"function foo() {} /** @enum {?} */ let E = {A: foo() };", CheckEnums_ENUM_TYPE_NOT_STRING_OR_NUMBER)
	c.TestWarning(
		`/** @typedef {number} */
var someName;
/** @enum {!someName} */
const EnumTwo = { FOO: EnumOne.FOO }
`,
		CheckEnums_ENUM_TYPE_NOT_STRING_OR_NUMBER)
	c.TestWarning(
		`/** @typedef {string} */
var someName;
/** @enum {!someName} */
const EnumTwo = { FOO: 'some' }
`,
		CheckEnums_ENUM_TYPE_NOT_STRING_OR_NUMBER)
	c.TestWarning(
		`/** @enum {number} */
let EnumOne = { FOO: 1, BAR: 2 };
/** @enum {!EnumOne}*/
let EnumTwo = { FOO: EnumOne.FOO }
`,
		CheckEnums_ENUM_TYPE_NOT_STRING_OR_NUMBER)
	c.TestWarning(
		`/** @enum {string} */
let EnumOne = { FOO: 'foo',  BAR: 'bar' };
/** @enum {!EnumOne}*/
let EnumTwo = { FOO: EnumOne.FOO }
`,
		CheckEnums_ENUM_TYPE_NOT_STRING_OR_NUMBER)
	c.TestWarning(" /** @enum {{x: number}} */ var E = {A: true };", CheckEnums_ENUM_TYPE_NOT_STRING_OR_NUMBER)
}

// Java: CheckEnumsTest.testNonStaticNumber_noWarning()
func (c *CheckEnumsTest) TestNonStaticNumber_noWarning() {
	// We do not warn on number enum values that are non-statically initialized.
	c.TestNoWarning("let fooVal = 1; /** @enum {number} */ let E = {A: fooVal };")
	c.TestNoWarning("let obj = {fooVal : 2}; /** @enum {number} */ let E = {A: obj.fooVal };")
	c.TestNoWarning("let obj = {fooVal : 2}; /** @enum {number} */ let E = {A: obj.fooVal + 1 };")
	c.TestNoWarning(
		"let obj = {fooVal : 2}; /** @enum {number} */ let E = {A: obj.fooVal - obj.fooVal};")
	c.TestNoWarning("function foo() {return 3;} /** @enum {number} */ let E = { A: foo() };")
}

// Java: CheckEnumsTest.testDefaultEnums_noWarning()
func (c *CheckEnumsTest) TestDefaultEnums_noWarning() {
	// Default enums are of number type; we must not warn when non-statically initialized
	c.TestNoWarning("let fooVal = 1; /** @enum */ let E = {A: fooVal };")
	c.TestNoWarning("let obj = {fooVal : 2}; /** @enum */ let E = {A: obj.fooVal };")
	c.TestNoWarning("function foo() {return 3;} /** @enum */ let E = { A: foo() };")
}

// Java: CheckEnumsTest.testStringEnums_withArithmeticOperations()
func (c *CheckEnumsTest) TestStringEnums_withArithmeticOperations() {
	// String enum values computed using arithmetic get reported
	c.TestWarning(
		"/** @enum {string} */ let E = { A: 'a'+ 2};",
		CheckEnums_NON_STATIC_INITIALIZER_STRING_VALUE_IN_ENUM)
	// 'a2'
	c.TestWarning(
		"/** @enum {string} */ let E = { A: 'a'+ 'b' };",
		CheckEnums_NON_STATIC_INITIALIZER_STRING_VALUE_IN_ENUM)
	// 'ab'
}

// Java: CheckEnumsTest.testNonStaticStringValue_reportsWarning()
func (c *CheckEnumsTest) TestNonStaticStringValue_reportsWarning() {
	c.TestWarning(
		"let fooVal = ''; /** @enum {string} */ let E = {A: fooVal };",
		CheckEnums_NON_STATIC_INITIALIZER_STRING_VALUE_IN_ENUM)
	// non-statically initialized with name
	c.TestWarning(
		`/** @return {string} */
function foo() {
  return ''
}
/** @enum {string} */
let E = { A: foo() };
`,
		CheckEnums_NON_STATIC_INITIALIZER_STRING_VALUE_IN_ENUM)
	// non-statically initialized with call
	// call to a function missing `@return` annotation
	c.TestWarning(
		`function foo() {
  return ''
}
/** @enum {string} */
let E = { A: foo() };
`,
		CheckEnums_NON_STATIC_INITIALIZER_STRING_VALUE_IN_ENUM)
	// call to a tagged template function
	c.TestWarning(
		"/** @return {string} */\nfunction tag(x) {\n  return x\n}\n/** @enum {string} */\nlet E = { A: tag`some` };\n",
		CheckEnums_NON_STATIC_INITIALIZER_STRING_VALUE_IN_ENUM)
	c.TestWarning(
		"let obj = {fooVal : ''}; /** @enum {string} */ let E = {A: obj.fooVal };",
		CheckEnums_NON_STATIC_INITIALIZER_STRING_VALUE_IN_ENUM)
	// non-statically initialized with GETPROP
	c.TestWarning(
		"let obj = {fooVal : ''}; /** @enum {string} */ let E = {A: obj.fooVal + '' };",
		CheckEnums_NON_STATIC_INITIALIZER_STRING_VALUE_IN_ENUM)
	// non-statically initialized with a
	// GETPROP operand in `+`
	c.TestWarning(
		"/** @type {number} */ let num = 5; /** @enum {string} */ let E = { A: 'a'+ num };",
		CheckEnums_NON_STATIC_INITIALIZER_STRING_VALUE_IN_ENUM)
	// non-statically initialized with var `num`.
	c.TestWarning(
		"/** @type {number} */ let num = 5; /** @enum {string} */ let E = { A: `a${num}`};",
		CheckEnums_NON_STATIC_INITIALIZER_STRING_VALUE_IN_ENUM)
	// non-statically initialized with
	// template lit substitution.
	c.TestWarning(
		`/** @enum {string} */
let EnumOne = {  FOO: 'foo',  BAR: 'bar'};
/** @enum {string}*/
let EnumTwo = { FOO: EnumOne.FOO}
`,
		CheckEnums_NON_STATIC_INITIALIZER_STRING_VALUE_IN_ENUM)
	// using another string enum as value gets
	// reported
}
