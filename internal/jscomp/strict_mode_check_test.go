// Port of com/google/javascript/jscomp/StrictModeCheckTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/StrictModeCheckTest.java)

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
	"testing"
)

type StrictModeCheckTest struct {
	CompilerTestCaseBase
}

// Java: StrictModeCheckTest()
func NewStrictModeCheckTest() *StrictModeCheckTest {
	c := &StrictModeCheckTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, strictModeCheckTest_EXTERNS)
	return c
}

func TestStrictModeCheckTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewStrictModeCheckTest() })
}

// Java: StrictModeCheckTest.EXTERNS
var strictModeCheckTest_EXTERNS = CompilerTestCase_DEFAULT_EXTERNS + "var arguments; function eval(str) {}"

// Java: StrictModeCheckTest.setUp()
func (c *StrictModeCheckTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.DisableTypeCheck()
}

// Java: StrictModeCheckTest.getProcessor(Compiler)
func (c *StrictModeCheckTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewStrictModeCheck(compiler, CheckLevel_WARNING)
}

// Java: StrictModeCheckTest.testSameEs6Strict(String)
func (c *StrictModeCheckTest) TestSameEs6Strict(js string) {
	c.TestSame(js)
}

// Java: StrictModeCheckTest.testUseOfWith1()
func (c *StrictModeCheckTest) TestUseOfWith1() {
	c.TestWarning("var a; with(a){}", StrictModeCheck_USE_OF_WITH)
}

// Java: StrictModeCheckTest.testUseOfWith2()
func (c *StrictModeCheckTest) TestUseOfWith2() {
	c.TestSame(
		`var a;
/** @suppress {with} */
with(a){}
`)
}

// Java: StrictModeCheckTest.testUseOfWith3()
func (c *StrictModeCheckTest) TestUseOfWith3() {
	c.TestSame(
		`function f(expr, context) {
  try {
    /** @suppress{with} */ with (context) {
      return eval('[' + expr + '][0]');
    }
  } catch (e) {
    return null;
  }
};
`)
}

// Java: StrictModeCheckTest.testEval2()
func (c *StrictModeCheckTest) TestEval2() {
	c.TestWarning("function foo(eval) {}", StrictModeCheck_EVAL_DECLARATION)
}

// Java: StrictModeCheckTest.testEval3()
func (c *StrictModeCheckTest) TestEval3() {
	c.TestSame("function foo() {} foo.eval = 3;")
}

// Java: StrictModeCheckTest.testEval4()
func (c *StrictModeCheckTest) TestEval4() {
	c.TestWarning("function foo() { var eval = 3; }", StrictModeCheck_EVAL_DECLARATION)
}

// Java: StrictModeCheckTest.testEval5()
func (c *StrictModeCheckTest) TestEval5() {
	c.TestWarning(
		"/** @suppress {duplicate} */ function eval() {}", StrictModeCheck_EVAL_DECLARATION)
}

// Java: StrictModeCheckTest.testEval6()
func (c *StrictModeCheckTest) TestEval6() {
	c.TestWarning("try {} catch (eval) {}", StrictModeCheck_EVAL_DECLARATION)
}

// Java: StrictModeCheckTest.testEval7()
func (c *StrictModeCheckTest) TestEval7() {
	c.TestSame("var o = {eval: 3};")
}

// Java: StrictModeCheckTest.testEval8()
func (c *StrictModeCheckTest) TestEval8() {
	c.TestSame("var a; eval: while (true) { a = 3; }")
}

// Java: StrictModeCheckTest.testUnknownVariable3()
func (c *StrictModeCheckTest) TestUnknownVariable3() {
	c.TestSame("try {} catch (ex) { ex = 3; }")
}

// Java: StrictModeCheckTest.testUnknownVariable4()
func (c *StrictModeCheckTest) TestUnknownVariable4() {
	c.DisableTypeCheck()
	c.TestSameEs6Strict("function foo(a) { let b; a = b; }")
	c.TestSameEs6Strict("function foo(a) { const b = 42; a = b; }")
}

// Java: StrictModeCheckTest.testArguments()
func (c *StrictModeCheckTest) TestArguments() {
	c.TestWarning("function foo(arguments) {}", StrictModeCheck_ARGUMENTS_DECLARATION)
}

// Java: StrictModeCheckTest.testArguments2()
func (c *StrictModeCheckTest) TestArguments2() {
	c.TestWarning("function foo() { var arguments = 3; }", StrictModeCheck_ARGUMENTS_DECLARATION)
}

// Java: StrictModeCheckTest.testArguments3()
func (c *StrictModeCheckTest) TestArguments3() {
	c.TestWarning(
		"/** @suppress {duplicate,checkTypes} */ function arguments() {}",
		StrictModeCheck_ARGUMENTS_DECLARATION)
}

// Java: StrictModeCheckTest.testArguments4()
func (c *StrictModeCheckTest) TestArguments4() {
	c.TestWarning("try {} catch (arguments) {}", StrictModeCheck_ARGUMENTS_DECLARATION)
}

// Java: StrictModeCheckTest.testArguments5()
func (c *StrictModeCheckTest) TestArguments5() {
	c.TestSame("var o = {arguments: 3};")
}

// Java: StrictModeCheckTest.testArguments6()
func (c *StrictModeCheckTest) TestArguments6() {
	c.DisableTypeCheck()
	c.TestSame("(() => arguments)();")
}

// Java: StrictModeCheckTest.testArgumentsCallee()
func (c *StrictModeCheckTest) TestArgumentsCallee() {
	c.TestWarning("function foo() {arguments.callee}", StrictModeCheck_ARGUMENTS_CALLEE_FORBIDDEN)
}

// Java: StrictModeCheckTest.testArgumentsCaller()
func (c *StrictModeCheckTest) TestArgumentsCaller() {
	c.TestWarning("function foo() {arguments.caller}", StrictModeCheck_ARGUMENTS_CALLER_FORBIDDEN)
}

// Java: StrictModeCheckTest.testFunctionCallerProp()
func (c *StrictModeCheckTest) TestFunctionCallerProp() {
	c.EnableTypeCheck()
	c.TestWarning("function foo() {foo.caller}", StrictModeCheck_FUNCTION_CALLER_FORBIDDEN)
}

// Java: StrictModeCheckTest.testFunctionArgumentsProp()
func (c *StrictModeCheckTest) TestFunctionArgumentsProp() {
	c.EnableTypeCheck()
	c.TestWarning(
		"function foo() {foo.arguments}", StrictModeCheck_FUNCTION_ARGUMENTS_PROP_FORBIDDEN)
}

// Java: StrictModeCheckTest.testEvalAssignment()
func (c *StrictModeCheckTest) TestEvalAssignment() {
	c.TestWarning(
		"/** @suppress {checkTypes} */ function foo() { eval = []; }",
		StrictModeCheck_EVAL_ASSIGNMENT)
}

// Java: StrictModeCheckTest.testAssignToArguments()
func (c *StrictModeCheckTest) TestAssignToArguments() {
	c.TestWarning("function foo() { arguments = []; }", StrictModeCheck_ARGUMENTS_ASSIGNMENT)
}

// Java: StrictModeCheckTest.testDeleteVar()
func (c *StrictModeCheckTest) TestDeleteVar() {
	c.TestWarning("var a; delete a", StrictModeCheck_DELETE_VARIABLE)
}

// Java: StrictModeCheckTest.testDeleteFunction()
func (c *StrictModeCheckTest) TestDeleteFunction() {
	c.TestWarning("function a() {} delete a", StrictModeCheck_DELETE_VARIABLE)
}

// Java: StrictModeCheckTest.testDeleteArgument()
func (c *StrictModeCheckTest) TestDeleteArgument() {
	c.TestWarning("function b(a) { delete a; }", StrictModeCheck_DELETE_VARIABLE)
}

// Java: StrictModeCheckTest.testValidDelete()
func (c *StrictModeCheckTest) TestValidDelete() {
	c.TestSame("var obj = { a: 0 }; delete obj.a;")
	c.TestSame("var obj = { a: function() {} }; delete obj.a;")
	c.TestSameEs6Strict("var obj = { a(){} }; delete obj.a;")
	c.TestSameEs6Strict("var obj = { a }; delete obj.a;")
}

// Java: StrictModeCheckTest.testDeleteProperty()
func (c *StrictModeCheckTest) TestDeleteProperty() {
	c.TestSame("/** @suppress {checkTypes} */ function f(obj) { delete obj.a; }")
}

// Java: StrictModeCheckTest.testAllowNumbersAsObjlitKeys()
func (c *StrictModeCheckTest) TestAllowNumbersAsObjlitKeys() {
	c.TestSame("var o = {1: 3, 2: 4};")
}

// Java: StrictModeCheckTest.testDuplicateObjectLiteralKey()
func (c *StrictModeCheckTest) TestDuplicateObjectLiteralKey() {
	c.TestSame("var o = {a: 1, b: 2, c: 3};")
	c.TestSame("var x = { get a() {}, set a(p) {} };")

	c.TestWarning("var o = {a: 1, b: 2, a: 3};", StrictModeCheck_DUPLICATE_MEMBER)
	c.TestWarning("var x = { get a() {}, get a() {} };", StrictModeCheck_DUPLICATE_MEMBER)
	c.TestWarning("var x = { get a() {}, a: 1 };", StrictModeCheck_DUPLICATE_MEMBER)
	c.TestWarning("var x = { set a(p) {}, a: 1 };", StrictModeCheck_DUPLICATE_MEMBER)

	c.TestSame(
		`'use strict';
/** @constructor */ function App() {}
App.prototype = {
  get appData() { return this.appData_; },
  set appData(data) { this.appData_ = data; }
};
`)

	c.TestWarning("var x = {a: 2, a(){}}", StrictModeCheck_DUPLICATE_MEMBER)
	c.TestWarning("var x = {a, a(){}}", StrictModeCheck_DUPLICATE_MEMBER)
	c.TestWarning("var x = {a(){}, a(){}}", StrictModeCheck_DUPLICATE_MEMBER)
}

// Java: StrictModeCheckTest.testFunctionDecl()
func (c *StrictModeCheckTest) TestFunctionDecl() {
	c.TestSame("function g() {}")
	c.TestSame("var g = function() {};")
	c.TestSame("(function() {})();")
	c.TestSame("(function() {});")
	c.TestSame(StrictModeCheckTest_InFn("function g() {}"))
	c.TestSame(StrictModeCheckTest_InFn("var g = function() {};"))
	c.TestSame(StrictModeCheckTest_InFn("(function() {})();"))
	c.TestSame(StrictModeCheckTest_InFn("(function() {});"))

	c.TestSame("{var g = function () {}}")
	c.TestSame("{(function g() {})()}")

	c.TestSame("var x;if (x) {var g = function () {}}")
	c.TestSame("var x;if (x) {(function g() {})()}")
}

// Java: StrictModeCheckTest.testClass()
func (c *StrictModeCheckTest) TestClass() {
	c.TestSame(
		`class A {
  method1() {}
  method2() {}
}
`)

	// Duplicate class methods test
	c.TestWarning(
		`class A {
  method1() {}
  method1() {}
}
`,
		StrictModeCheck_DUPLICATE_MEMBER)

	// Function declaration / call test.
	// The two following tests should have reported FUNCTION_CALLER_FORBIDDEN and
	// FUNCTION_ARGUMENTS_PROP_FORBIDDEN. Typecheck needed for them to work.
	// TODO(user): Add tests for these after typecheck supports class.
	c.TestSame(
		`class A {
  method() {this.method.caller}
}
`)
	c.TestSame(
		`class A {
  method() {this.method.arguments}
}
`)

	// Duplicate obj literal key in classes
	c.TestWarning(
		`class A {
  method() {
    var obj = {a : 1, a : 2}
  }
}
`,
		StrictModeCheck_DUPLICATE_MEMBER)

	// Delete test. Class methods are configurable, thus deletable.
	c.TestSame(
		`class A {
  methodA() {}
  methodB() {delete this.methodA}
}
`)

	// Use of with test
	c.TestWarning(
		`class A {
  constructor() {this.x = 1;}
  method() {
    with (this.x) {}
  }
}
`,
		StrictModeCheck_USE_OF_WITH)

	// Eval errors test
	c.TestWarning(
		`class A {
  method(eval) {}
}
`,
		StrictModeCheck_EVAL_DECLARATION)
	c.TestWarning(
		`class A {
  method() {var eval = 1;}
}
`,
		StrictModeCheck_EVAL_DECLARATION)
	c.TestWarning(
		`class A {
  method() {eval = 1}
}
`,
		StrictModeCheck_EVAL_ASSIGNMENT)

	// Use of 'arguments'
	c.TestWarning(
		`class A {
  method(arguments) {}
}
`,
		StrictModeCheck_ARGUMENTS_DECLARATION)
	c.TestWarning(
		`class A {
  method() {var arguments = 1;}
}
`,
		StrictModeCheck_ARGUMENTS_DECLARATION)
	c.TestWarning(
		`class A {
  method() {arguments = 1}
}
`,
		StrictModeCheck_ARGUMENTS_ASSIGNMENT)
	c.TestWarning(
		`class A {
  method() {arguments.callee}
}
`,
		StrictModeCheck_ARGUMENTS_CALLEE_FORBIDDEN)
	c.TestWarning(
		`class A {
  method() {arguments.caller}
}
`,
		StrictModeCheck_ARGUMENTS_CALLER_FORBIDDEN)
}

// Java: StrictModeCheckTest.testComputedPropInClass()
func (c *StrictModeCheckTest) TestComputedPropInClass() {
	c.TestSame(
		`class Example {
  [computed()]() {}
  [computed()]() {}
}
`)
}

// Java: StrictModeCheckTest.testStaticAndNonstaticMethodWithSameName()
func (c *StrictModeCheckTest) TestStaticAndNonstaticMethodWithSameName() {
	c.TestSame(
		`class Example {
  foo() {}
  static foo() {}
}
`)
}

// Java: StrictModeCheckTest.testStaticAndNonstaticGetterWithSameName()
func (c *StrictModeCheckTest) TestStaticAndNonstaticGetterWithSameName() {
	c.TestSame(
		`class Example {
  get foo() {}
  static get foo() {}
}
`)
}

// Java: StrictModeCheckTest.testStaticAndNonstaticSetterWithSameName()
func (c *StrictModeCheckTest) TestStaticAndNonstaticSetterWithSameName() {
	c.TestSame(
		`class Example {
  set foo(x) {}
  static set foo(x) {}
}
`)
}

// Java: StrictModeCheckTest.testClassWithEmptyMembers()
func (c *StrictModeCheckTest) TestClassWithEmptyMembers() {
	c.TestWarning("class Foo { dup() {}; dup() {}; }", StrictModeCheck_DUPLICATE_MEMBER)
}

// Java: StrictModeCheckTest.testClassDuplicateFieldError()
func (c *StrictModeCheckTest) TestClassDuplicateFieldError() {
	c.TestWarning(
		`class Example {
  a;
  a;
}
`,
		StrictModeCheck_DUPLICATE_MEMBER)
	c.TestWarning(
		`class Example {
  static a = 2;
  static a = 3;
}
`,
		StrictModeCheck_DUPLICATE_MEMBER)
}

// Java: StrictModeCheckTest.testClassDuplicateStaticFieldError()
func (c *StrictModeCheckTest) TestClassDuplicateStaticFieldError() {
	c.TestWarning(
		`class Example {
  static a = 2;
  static a = 3;
}
`,
		StrictModeCheck_DUPLICATE_MEMBER)
	c.TestWarning(
		`class Example {
  static a;
  static a;
}
`,
		StrictModeCheck_DUPLICATE_MEMBER)
}

// Java: StrictModeCheckTest.testClassFields_noError()
func (c *StrictModeCheckTest) TestClassFields_noError() {
	c.TestSame(
		`class Example {
  a;
  static a;
}
`)
}

// Java: StrictModeCheckTest.testClassStaticBlocksDuplicates()
func (c *StrictModeCheckTest) TestClassStaticBlocksDuplicates() {
	// testing that duplicates are ok in class static blocks
	c.TestSame("class Foo {static {this.a; this.a;}}")

	// even if they are in 2 different static blocks
	c.TestSame("class Foo {static {this.a} static {this.a}}")

	// testing that duplicates between class static blocks and public fields are ok
	c.TestSame("class Foo {static x; static {this.x}}")

	// duplicate functions are also okay
	c.TestSame("class Foo {static f() {}; static{this.f = function g() {}}}")

	// testing that duplicates in obj literals are still invalid
	c.TestWarning("class Foo {static {var obj = {a : 1, a : 2};}}", StrictModeCheck_DUPLICATE_MEMBER)
}

// Java: StrictModeCheckTest.testClassStaticBlocksWith()
func (c *StrictModeCheckTest) TestClassStaticBlocksWith() {
	c.TestWarning(
		`class A {
  static x;
  static {
    with (this.x) {}
  }
}
`,
		StrictModeCheck_USE_OF_WITH)
}

// Java: StrictModeCheckTest.testClassStaticBlocksEval()
func (c *StrictModeCheckTest) TestClassStaticBlocksEval() {
	c.TestWarning("class A { static {var eval;}}", StrictModeCheck_EVAL_DECLARATION)
	c.TestWarning("class A { static {function eval() {};}}", StrictModeCheck_EVAL_DECLARATION)
	c.TestWarning(
		"class A { static {this.e = function eval() {};}}", StrictModeCheck_EVAL_DECLARATION)
	c.TestWarning("class A { static {eval = 3;}}", StrictModeCheck_EVAL_ASSIGNMENT)

	// eval as a field is okay
	c.TestSame("class A { static {this.eval;}}")
}

// Java: StrictModeCheckTest.testClassStaticBlocksArguments()
func (c *StrictModeCheckTest) TestClassStaticBlocksArguments() {
	c.TestWarning(
		`class A {
  static {var arguments = 1;}
}
`,
		StrictModeCheck_ARGUMENTS_DECLARATION)
	c.TestWarning(
		`class A {
  static {arguments = 1}
}
`,
		StrictModeCheck_ARGUMENTS_ASSIGNMENT)
	c.TestWarning(
		`class A {
  static {arguments.callee}
}
`,
		StrictModeCheck_ARGUMENTS_CALLEE_FORBIDDEN)
	c.TestWarning(
		`class A {
  static {arguments.caller}
}
`,
		StrictModeCheck_ARGUMENTS_CALLER_FORBIDDEN)
}

// Java: StrictModeCheckTest.testClassStaticBlockDelete()
func (c *StrictModeCheckTest) TestClassStaticBlockDelete() {
	c.TestSame("class A {static{this.a; delete this.a;}}")
	c.TestSame("class A {static a; static{delete this.a;}}")
	c.TestWarning("class A {static a; static{var a; delete a;}}", StrictModeCheck_DELETE_VARIABLE)
	c.TestSame("class A {static f() {}; static{delete this.f;}}")
}

// Java: StrictModeCheckTest.inFn(String)
func StrictModeCheckTest_InFn(body string) string {
	return "function func() {" + body + "}"
}
