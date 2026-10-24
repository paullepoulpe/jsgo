// Port of com/google/javascript/jscomp/CollapseAnonymousFunctionsTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/CollapseAnonymousFunctionsTest.java)

/*
 * Copyright 2008 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"testing"
)

type CollapseAnonymousFunctionsTest struct {
	CompilerTestCaseBase
}

// Java: CollapseAnonymousFunctionsTest()
func NewCollapseAnonymousFunctionsTest() *CollapseAnonymousFunctionsTest {
	c := &CollapseAnonymousFunctionsTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestCollapseAnonymousFunctionsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCollapseAnonymousFunctionsTest() })
}

// JUnitBefore runs the @Before methods, superclass ones first:
// CompilerTestCase.setUp(), then customSetUp().
func (c *CollapseAnonymousFunctionsTest) JUnitBefore() {
	c.CompilerTestCaseBase.JUnitBefore()
	c.CustomSetUp()
}

// Java: CollapseAnonymousFunctionsTest.customSetUp()
func (c *CollapseAnonymousFunctionsTest) CustomSetUp() {
	c.EnableNormalize()
}

// Java: CollapseAnonymousFunctionsTest.getProcessor(Compiler)
func (c *CollapseAnonymousFunctionsTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewCollapseAnonymousFunctions(compiler)
}

// Java: CollapseAnonymousFunctionsTest.testGlobalScope()
func (c *CollapseAnonymousFunctionsTest) TestGlobalScope() {
	c.Test("var f = function(){}", "function f(){}")
}

// Java: CollapseAnonymousFunctionsTest.testSeparateDeclarationAndInitilization_unchanged()
func (c *CollapseAnonymousFunctionsTest) TestSeparateDeclarationAndInitilization_unchanged() {
	c.TestSame("var f; f = function(){}")
}

// Java: CollapseAnonymousFunctionsTest.testLocalScope1()
func (c *CollapseAnonymousFunctionsTest) TestLocalScope1() {
	c.Test("function f(){ var x = function(){}; x() }", "function f(){ function x(){} x() }")
}

// Java: CollapseAnonymousFunctionsTest.testLocalScope2()
func (c *CollapseAnonymousFunctionsTest) TestLocalScope2() {
	c.Test(
		"function f(){ var x = function(){}; return x }",
		"function f(){ function x(){} return x }")
}

// Java: CollapseAnonymousFunctionsTest.testLet()
func (c *CollapseAnonymousFunctionsTest) TestLet() {
	c.Test("let f = function() {};", "function f() {}")
	fnName := c.GetLastCompiler().GetJsRoot().GetFirstFirstChild().GetFirstChild()
	jutil.CheckState(fnName.IsName())
	jutil.CheckState(fnName.GetString() == "f")
	// The function name `f` is not marked as IS_CONSTANT_NAME property as LHS declaration is
	// `let`.
	jutil.CheckState(!fnName.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME))

	// The function name `f` and call name `f` are not marked as IS_CONSTANT_NAME property as LHS
	// declaration is `let`
	c.Test("let f = function() {}; f();", "function f(){} f();")
	fnName = c.GetLastCompiler().GetJsRoot().GetFirstFirstChild().GetFirstChild()
	jutil.CheckState(fnName.IsName())
	jutil.CheckState(fnName.GetString() == "f")
	jutil.CheckState(!fnName.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME))
	calledFnName := fnName.GetParent().GetNext().GetFirstFirstChild()
	jutil.CheckState(calledFnName.IsName())
	jutil.CheckState(calledFnName.GetString() == "f")
	jutil.CheckState(!calledFnName.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME))
}

// Java: CollapseAnonymousFunctionsTest.testConst()
func (c *CollapseAnonymousFunctionsTest) TestConst() {
	// The function name `f` gets marked as IS_CONSTANT_NAME property as LHS declaration is `const`.
	c.Test("const f = function() {};", "function f() {}")
	fnName := c.GetLastCompiler().GetJsRoot().GetFirstFirstChild().GetFirstChild()
	jutil.CheckState(fnName.IsName())
	jutil.CheckState(fnName.GetString() == "f")
	jutil.CheckState(fnName.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME))

	// The function name `f` and call name `f` stays marked as IS_CONSTANT_NAME property
	c.Test("const f = function() {}; f();", "function f(){} f();")
	fnName = c.GetLastCompiler().GetJsRoot().GetFirstFirstChild().GetFirstChild()
	jutil.CheckState(fnName.IsName())
	jutil.CheckState(fnName.GetString() == "f")
	jutil.CheckState(fnName.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME))
	calledFnName := fnName.GetParent().GetNext().GetFirstFirstChild()
	jutil.CheckState(calledFnName.IsName())
	jutil.CheckState(calledFnName.GetString() == "f")
	jutil.CheckState(calledFnName.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME))
}

// Java: CollapseAnonymousFunctionsTest.testArrow()
func (c *CollapseAnonymousFunctionsTest) TestArrow() {
	c.TestSame("function f() { var g = () => { return this; }; }")
	// It would be safe to collapse this one because it doesn't reference 'this' or 'arguments'
	// but (at least for now) we don't.
	c.TestSame("let f = () => { console.log(0); };")
}

// Java: CollapseAnonymousFunctionsTest.testDestructuring()
func (c *CollapseAnonymousFunctionsTest) TestDestructuring() {
	c.TestSame("var name; ({name} = function() {});")
	c.TestSame("let {name} = function() {};")
	c.TestSame("const {name} = function() {};")

	c.Test("var [x] = function() {};", "var x; [x] = function() {};")
	c.TestSame("let [x] = function() {};")
	c.TestSame("const [x] = function() {};")
}

// Java: CollapseAnonymousFunctionsTest.testVarNotImmediatelyBelowScriptOrBlock1()
func (c *CollapseAnonymousFunctionsTest) TestVarNotImmediatelyBelowScriptOrBlock1() {
	c.TestSame("if (x) var f = function(){}")
}

// Java: CollapseAnonymousFunctionsTest.testVarNotImmediatelyBelowScriptOrBlock2()
func (c *CollapseAnonymousFunctionsTest) TestVarNotImmediatelyBelowScriptOrBlock2() {
	c.TestSame(
		"var x = 1;\n" +
			"if (x == 1) {\n" +
			"  var f = function () { alert('b')}\n" +
			"} else {\n" +
			"  f = function() { alert('c')}\n" +
			"}\n" +
			"f();\n")
}

// Java: CollapseAnonymousFunctionsTest.testVarNotImmediatelyBelowScriptOrBlock3()
func (c *CollapseAnonymousFunctionsTest) TestVarNotImmediatelyBelowScriptOrBlock3() {
	c.TestSame("var x = 1; if (x) {var f = function(){return x}; f(); x--;}")
}

// TODO(tbreisacher): We could collapse in this case, but we probably need to tweak Normalize.
// Java: CollapseAnonymousFunctionsTest.testLetNotImmediatelyBelowScriptOrBlock1()
func (c *CollapseAnonymousFunctionsTest) TestLetNotImmediatelyBelowScriptOrBlock1() {
	c.TestSame("if (x) let f = function() {};")
}

// TODO(tbreisacher): We could collapse in this case, but we probably need to tweak Normalize.
// Java: CollapseAnonymousFunctionsTest.testLetNotImmediatelyBelowScriptOrBlock2()
func (c *CollapseAnonymousFunctionsTest) TestLetNotImmediatelyBelowScriptOrBlock2() {
	c.TestSame("let x = 1; if (x) { let f = function() {return x}; f(); x--;}")
}

// Java: CollapseAnonymousFunctionsTest.testMultipleVar()
func (c *CollapseAnonymousFunctionsTest) TestMultipleVar() {
	c.Test("var f = function(){}; var g = f", "function f(){} var g = f")
}

// Java: CollapseAnonymousFunctionsTest.testMultipleVar2()
func (c *CollapseAnonymousFunctionsTest) TestMultipleVar2() {
	c.Test(
		"var f = function(){}; var g = f; var h = function(){}",
		"function f(){}var g = f;function h(){}")
}

// Java: CollapseAnonymousFunctionsTest.testBothScopes()
func (c *CollapseAnonymousFunctionsTest) TestBothScopes() {
	c.Test("var x = function() { var y = function(){} }", "function x() { function y(){} }")
}

// Java: CollapseAnonymousFunctionsTest.testLocalScopeOnly1()
func (c *CollapseAnonymousFunctionsTest) TestLocalScopeOnly1() {
	c.Test(
		"if (x) var f = function(){ var g = function(){} }",
		"if (x) var f = function(){ function g(){} }")
}

// Java: CollapseAnonymousFunctionsTest.testLocalScopeOnly2()
func (c *CollapseAnonymousFunctionsTest) TestLocalScopeOnly2() {
	c.Test(
		"if (x) var f = function(){ var g = function(){} };",
		"if (x) var f = function(){ function g(){} }")
}

// Java: CollapseAnonymousFunctionsTest.testLocalScopeOnly3()
func (c *CollapseAnonymousFunctionsTest) TestLocalScopeOnly3() {
	c.Test(
		"if (x){ var f = function(){ var g = function(){} };}",
		"if (x){ var f = function(){ function g(){} }}")
}

// Java: CollapseAnonymousFunctionsTest.testReturn()
func (c *CollapseAnonymousFunctionsTest) TestReturn() {
	c.Test(
		"var f = function(x){return 2*x}; var g = f(2);", "function f(x){return 2*x} var g = f(2)")
}

// Java: CollapseAnonymousFunctionsTest.testAlert()
func (c *CollapseAnonymousFunctionsTest) TestAlert() {
	c.Test("var x = 1; var f = function(){alert(x)};", "var x = 1; function f(){alert(x)}")
}

// Java: CollapseAnonymousFunctionsTest.testRecursiveInternal1()
func (c *CollapseAnonymousFunctionsTest) TestRecursiveInternal1() {
	c.TestSame("var f = function foo() { foo() }")
}

// Java: CollapseAnonymousFunctionsTest.testRecursiveInternal2()
func (c *CollapseAnonymousFunctionsTest) TestRecursiveInternal2() {
	c.TestSame("var f = function foo() { function g(){foo()} g() }")
}

// Java: CollapseAnonymousFunctionsTest.testRecursiveExternal1()
func (c *CollapseAnonymousFunctionsTest) TestRecursiveExternal1() {
	c.Test("var f = function foo() { f() }", "function f() { f() }")
}

// Java: CollapseAnonymousFunctionsTest.testRecursiveExternal2()
func (c *CollapseAnonymousFunctionsTest) TestRecursiveExternal2() {
	c.Test(
		"var f = function foo() { function g(){f()} g() }",
		"function f() { function g(){f()} g() }")
}

// Java: CollapseAnonymousFunctionsTest.testConstantFunction1()
func (c *CollapseAnonymousFunctionsTest) TestConstantFunction1() {
	c.Test("var FOO = function(){};FOO()", "function FOO(){}FOO()")
}

// Java: CollapseAnonymousFunctionsTest.testInnerFunction1()
func (c *CollapseAnonymousFunctionsTest) TestInnerFunction1() {
	c.Test(
		"function f() {\n"+
			"  var x = 3;\n"+
			"  var y = function() { return 4; };\n"+
			"  return x + y();\n"+
			"}\n",
		"function f() {\n"+
			"  function y() { return 4; }\n"+
			"  var x = 3;\n"+
			"  return x + y();\n"+
			"}\n")
}
