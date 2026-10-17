// Port of com/google/javascript/jscomp/FunctionRewriterTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/FunctionRewriterTest.java)

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"testing"
)

type FunctionRewriterTest struct {
	CompilerTestCaseBase
}

// Java: FunctionRewriterTest()
func NewFunctionRewriterTest() *FunctionRewriterTest {
	c := &FunctionRewriterTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestFunctionRewriterTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewFunctionRewriterTest() })
}

const FunctionRewriterTest_RETURNARG_HELPER = "function JSCompiler_returnArg(JSCompiler_returnArg_value){\n" +
	"  return function() { return JSCompiler_returnArg_value }\n" +
	"}\n"
const FunctionRewriterTest_GET_HELPER = "function JSCompiler_get(JSCompiler_get_name){\n" +
	"  return function() { return this[JSCompiler_get_name] }\n" +
	"}\n"
const FunctionRewriterTest_SET_HELPER = "function JSCompiler_set(JSCompiler_set_name) {\n" +
	"  return function(JSCompiler_set_value){\n" +
	"    this[JSCompiler_set_name]=JSCompiler_set_value\n" +
	"  }\n" +
	"}\n"
const FunctionRewriterTest_EMPTY_HELPER = "function JSCompiler_emptyFn() {\n" +
	"  return function(){}\n" +
	"}\n"
const FunctionRewriterTest_IDENTITY_HELPER = "function JSCompiler_identityFn() {\n" +
	"  return function(JSCompiler_identityFn_value) {\n" +
	"      return JSCompiler_identityFn_value\n" +
	"  }\n" +
	"}\n"

// Java: FunctionRewriterTest.setUp()
func (c *FunctionRewriterTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableNormalize()
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
}

// Java: FunctionRewriterTest.getProcessor(Compiler)
func (c *FunctionRewriterTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewFunctionRewriter(compiler)
}

// Java: FunctionRewriterTest.testEs6Class()
func (c *FunctionRewriterTest) TestEs6Class() {
	// There is never any benefit to replacing ES6 class methods
	c.CheckCompilesToSame(
		"class C {\n"+
			"  constructor(x = 1) { // looks like a setter\n"+
			"    this.x_ = x;\n"+
			"  }\n"+
			"  get x() {\n"+
			"    return this.x_\n"+
			"  }\n"+
			"  getX() {\n"+
			"    return this.x_\n"+
			"  }\n"+
			"  set x(x) {\n"+
			"    this.x_ = x;\n"+
			"  }\n"+
			"  setX(x = 3) {\n"+
			"    this.x_ = x;\n"+
			"  }\n"+
			"  getConst() {\n"+
			"    return 1;\n"+
			"  }\n"+
			"  empty() {}\n"+
			"  identity(x) { return x; }\n"+
			"}\n",
		10)
}

// Java: FunctionRewriterTest.testReplaceReturnConst1()
func (c *FunctionRewriterTest) TestReplaceReturnConst1() {
	source := "a.prototype.foo = function() {return \"foobar\"}"
	c.CheckCompilesToSame(source, 3)
	c.CheckCompilesTo(
		source, FunctionRewriterTest_RETURNARG_HELPER, "a.prototype.foo = JSCompiler_returnArg(\"foobar\")", 4)
}

// Java: FunctionRewriterTest.testReplaceReturnConst2()
func (c *FunctionRewriterTest) TestReplaceReturnConst2() {
	c.CheckCompilesToSame("a.prototype.foo = function() {return foobar}", 10)
}

// Java: FunctionRewriterTest.testReplaceReturnConst3()
func (c *FunctionRewriterTest) TestReplaceReturnConst3() {
	source := "a.prototype.foo = function() {return void 0;}"
	c.CheckCompilesToSame(source, 3)
	c.CheckCompilesTo(source, FunctionRewriterTest_RETURNARG_HELPER, "a.prototype.foo = JSCompiler_returnArg(void 0)", 4)
}

// Java: FunctionRewriterTest.testReplaceGetter1()
func (c *FunctionRewriterTest) TestReplaceGetter1() {
	source := "a.prototype.foo = function() {return this.foo_}"
	c.CheckCompilesToSame(source, 3)
	c.CheckCompilesTo(source, FunctionRewriterTest_GET_HELPER, "a.prototype.foo = JSCompiler_get(\"foo_\")", 4)
}

// Java: FunctionRewriterTest.testReplaceGetter2()
func (c *FunctionRewriterTest) TestReplaceGetter2() {
	c.CheckCompilesToSame("a.prototype.foo = function() {return}", 10)
}

// Java: FunctionRewriterTest.testDontReplaceGetterInArrow()
func (c *FunctionRewriterTest) TestDontReplaceGetterInArrow() {
	c.CheckCompilesToSame("class C { foo() { return (x) => { return this.x; }; } }", 10)
	c.CheckCompilesToSame("class C { foo() { return (x) => this.x; } }", 10)
}

// Java: FunctionRewriterTest.testReplaceSetter1()
func (c *FunctionRewriterTest) TestReplaceSetter1() {
	source := "a.prototype.foo = function(v) {this.foo_ = v}"
	c.CheckCompilesToSame(source, 4)
	c.CheckCompilesTo(source, FunctionRewriterTest_SET_HELPER, "a.prototype.foo = JSCompiler_set(\"foo_\")", 5)
}

// Java: FunctionRewriterTest.testReplaceSetterWithDefault()
func (c *FunctionRewriterTest) TestReplaceSetterWithDefault() {
	source := "a.prototype.foo = function(v = 1) {this.foo_ = v}"
	c.CheckCompilesToSame(source, 10)
}

// Java: FunctionRewriterTest.testReplaceSetterWithDestructuring()
func (c *FunctionRewriterTest) TestReplaceSetterWithDestructuring() {
	source := "a.prototype.foo = function({v}) {this.foo_ = v}"
	c.CheckCompilesToSame(source, 10)
}

// Java: FunctionRewriterTest.testReplaceSetter2()
func (c *FunctionRewriterTest) TestReplaceSetter2() {
	source := "a.prototype.foo = function(v, v2) {this.foo_ = v}"
	c.CheckCompilesToSame(source, 3)
	c.CheckCompilesTo(source, FunctionRewriterTest_SET_HELPER, "a.prototype.foo = JSCompiler_set(\"foo_\")", 4)
}

// Java: FunctionRewriterTest.testReplaceSetter3()
func (c *FunctionRewriterTest) TestReplaceSetter3() {
	c.CheckCompilesToSame("a.prototype.foo = function() {this.foo_ = v}", 10)
}

// Java: FunctionRewriterTest.testReplaceSetter4()
func (c *FunctionRewriterTest) TestReplaceSetter4() {
	c.CheckCompilesToSame("a.prototype.foo = function(v, v2) {this.foo_ = v2}", 10)
}

// Java: FunctionRewriterTest.testDontReplaceSetterInArrow()
func (c *FunctionRewriterTest) TestDontReplaceSetterInArrow() {
	c.CheckCompilesToSame("class C { foo() { return (x) => { this.x = x; }; } }", 10)
	c.CheckCompilesToSame("class C { foo() { return (x) => this.x = x; } }", 10)
}

// Java: FunctionRewriterTest.testReplaceEmptyFunction1()
func (c *FunctionRewriterTest) TestReplaceEmptyFunction1() {
	source := "a.prototype.foo = function() {}"
	c.CheckCompilesToSame(source, 4)
	c.CheckCompilesTo(source, FunctionRewriterTest_EMPTY_HELPER, "a.prototype.foo = JSCompiler_emptyFn()", 5)
}

// Java: FunctionRewriterTest.testReplaceEmptyFunction2()
func (c *FunctionRewriterTest) TestReplaceEmptyFunction2() {
	c.CheckCompilesToSame("function foo() {}", 10)
}

// Java: FunctionRewriterTest.testReplaceEmptyFunction3()
func (c *FunctionRewriterTest) TestReplaceEmptyFunction3() {
	source := "var foo = function() {}"
	c.CheckCompilesToSame(source, 4)
	c.CheckCompilesTo(source, FunctionRewriterTest_EMPTY_HELPER, "var foo = JSCompiler_emptyFn()", 5)
}

// Java: FunctionRewriterTest.testReplaceEmptyFunction4()
func (c *FunctionRewriterTest) TestReplaceEmptyFunction4() {
	source := "var foo = async function() {}"
	c.CheckCompilesToSame(source, 10)
}

// Java: FunctionRewriterTest.testReplaceEmptyFunction5()
func (c *FunctionRewriterTest) TestReplaceEmptyFunction5() {
	source := "var foo = function *() {}"
	c.CheckCompilesToSame(source, 10)
}

// Java: FunctionRewriterTest.testReplaceEmptyArrowFunction()
func (c *FunctionRewriterTest) TestReplaceEmptyArrowFunction() {
	source := "var foo = () => {}"
	c.CheckCompilesToSame(source, 10)
	c.CheckCompilesTo(source, FunctionRewriterTest_EMPTY_HELPER, "var foo = JSCompiler_emptyFn()", 30)
}

// Java: FunctionRewriterTest.testReplaceIdentityFunction1()
func (c *FunctionRewriterTest) TestReplaceIdentityFunction1() {
	source := "a.prototype.foo = function(a) {return a}"
	c.CheckCompilesToSame(source, 2)
	c.CheckCompilesTo(source, FunctionRewriterTest_IDENTITY_HELPER, "a.prototype.foo = JSCompiler_identityFn()", 3)
}

// Java: FunctionRewriterTest.testReplaceIdentityFunctionWithDefault()
func (c *FunctionRewriterTest) TestReplaceIdentityFunctionWithDefault() {
	c.CheckCompilesToSame("a.prototype.foo = function(a = 1) {return a}", 10)
}

// Java: FunctionRewriterTest.testReplaceIdentityFunctionWithDestructuring()
func (c *FunctionRewriterTest) TestReplaceIdentityFunctionWithDestructuring() {
	c.CheckCompilesToSame("a.prototype.foo = function({a}) {return a}", 10)
}

// Java: FunctionRewriterTest.testReplaceIdentityFunction2()
func (c *FunctionRewriterTest) TestReplaceIdentityFunction2() {
	c.CheckCompilesToSame("a.prototype.foo = function(a) {return a + 1}", 10)
}

// Java: FunctionRewriterTest.testReplaceIdentityFunctionAsArrow()
func (c *FunctionRewriterTest) TestReplaceIdentityFunctionAsArrow() {
	c.CheckCompilesTo(
		"a.prototype.foo = (a) => { return a; }",
		FunctionRewriterTest_IDENTITY_HELPER,
		"a.prototype.foo = JSCompiler_identityFn()",
		20)
}

// Java: FunctionRewriterTest.testReplaceIdentityFunctionAsBlocklessArrow()
func (c *FunctionRewriterTest) TestReplaceIdentityFunctionAsBlocklessArrow() {
	source := "a.prototype.foo = (a) => a"
	c.CheckCompilesToSame(source, 2)
	c.CheckCompilesTo(source, FunctionRewriterTest_IDENTITY_HELPER, "a.prototype.foo = JSCompiler_identityFn()", 50)
}

// Java: FunctionRewriterTest.testIssue538()
func (c *FunctionRewriterTest) TestIssue538() {
	c.CheckCompilesToSame(
		"/** @constructor */\n"+
			"WebInspector.Setting = function() {}\n"+
			"WebInspector.Setting.prototype = {\n"+
			"    get name0(){return this._name;},\n"+
			"    get name1(){return this._name;},\n"+
			"    get name2(){return this._name;},\n"+
			"    get name3(){return this._name;},\n"+
			"    get name4(){return this._name;},\n"+
			"    get name5(){return this._name;},\n"+
			"    get name6(){return this._name;},\n"+
			"    get name7(){return this._name;},\n"+
			"    get name8(){return this._name;},\n"+
			"    get name9(){return this._name;},\n"+
			"}\n",
		1)
}

// Java: FunctionRewriterTest.testNoRewriteWithNonSimpleParameters()
func (c *FunctionRewriterTest) TestNoRewriteWithNonSimpleParameters() {
	// Empty function with default parameter side effect
	c.CheckCompilesToSame("a.prototype.foo = function(a = audit()) {}", 10)
	// Empty function with rest parameter
	c.CheckCompilesToSame("a.prototype.foo = function(...a) {}", 10)
	// Empty function with destructuring parameter
	c.CheckCompilesToSame("a.prototype.foo = function({a}) {}", 10)

	// Return constant with default parameter side effect
	c.CheckCompilesToSame("a.prototype.foo = function(a = audit()) { return 1; }", 10)
	// Return constant with rest parameter
	c.CheckCompilesToSame("a.prototype.foo = function(...a) { return 1; }", 10)

	// Getter with default parameter side effect
	c.CheckCompilesToSame("a.prototype.foo = function(a = audit()) { return this.x; }", 10)
	// Getter with destructuring parameter
	c.CheckCompilesToSame("a.prototype.foo = function({a}) { return this.x; }", 10)

	// Setter with default parameter side effect in a second parameter
	c.CheckCompilesToSame("a.prototype.foo = function(v, a = audit()) { this.x = v; }", 10)
	// Setter with destructuring in a second parameter
	c.CheckCompilesToSame("a.prototype.foo = function(v, {a}) { this.x = v; }", 10)

	// Identity with default parameter side effect in a second parameter
	c.CheckCompilesToSame("a.prototype.foo = function(v, a = audit()) { return v; }", 10)
	// Identity with destructuring in a second parameter
	c.CheckCompilesToSame("a.prototype.foo = function(v, {a}) { return v; }", 10)
}

// Java: FunctionRewriterTest.checkCompilesTo(String, String, String, int)
func (c *FunctionRewriterTest) CheckCompilesTo(src string, expectedHdr string, expectedBody string, repetitions int) {
	srcBuffer := jutil.NewStringBuilder()
	expectedBuffer := jutil.NewStringBuilder()

	expectedBuffer.Append(expectedHdr)

	for idx := 0; idx < repetitions; idx++ {
		if idx != 0 {
			srcBuffer.Append(";")
			expectedBuffer.Append(";")
		}
		srcBuffer.Append(src)
		expectedBuffer.Append(expectedBody)
	}
	c.Test(srcBuffer.String(), expectedBuffer.String())
}

// Java: FunctionRewriterTest.checkCompilesToSame(String, int)
func (c *FunctionRewriterTest) CheckCompilesToSame(src string, repetitions int) {
	c.CheckCompilesTo(src, "", src, repetitions)
}
