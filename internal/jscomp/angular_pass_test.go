// Port of com/google/javascript/jscomp/AngularPassTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/AngularPassTest.java)

/*
 * Copyright 2012 The Closure Compiler Authors.
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

// AngularPassTest tests AngularPass.
type AngularPassTest struct {
	CompilerTestCaseBase
}

// Java: AngularPassTest()
func NewAngularPassTest() *AngularPassTest {
	c := &AngularPassTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestAngularPassTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewAngularPassTest() })
}

// Java: AngularPassTest.getProcessor(Compiler)
func (c *AngularPassTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewAngularPass(compiler)
}

// Java: AngularPassTest.getOptions()
func (c *AngularPassTest) GetOptions() *CompilerOptions {
	options := c.CompilerTestCaseBase.GetOptions()
	options.SetWarningLevel(DiagnosticGroups_MODULE_LOAD, CheckLevel_OFF)
	// enables angularPass.
	options.SetAngularPass(true)
	return options
}

// Java: AngularPassTest.testNgInjectAddsInjectToFunctions()
func (c *AngularPassTest) TestNgInjectAddsInjectToFunctions() {
	c.Test(
		"/** @ngInject */ function fn(a, b) {}",
		"/** @ngInject */ function fn(a, b) {} /** @public */ fn['$inject']=['a', 'b']")

	c.TestSame("function fn(a, b) {}")
}

// Java: AngularPassTest.testNgInjectSetVisibility()
func (c *AngularPassTest) TestNgInjectSetVisibility() {
	c.Test(
		"/** @ngInject */ function fn(a, b) {}",
		"/** @ngInject */ function fn(a, b) {} /** @public */ fn['$inject']=['a', 'b']")
}

// Java: AngularPassTest.testNgInjectAddsInjectAfterGoogInherits()
func (c *AngularPassTest) TestNgInjectAddsInjectAfterGoogInherits() {
	c.Test(
		"/** @ngInject @constructor */\n"+
			"function fn(a, b) {}\n"+
			"goog.inherits(fn, parent);\n",
		"/** @ngInject @constructor */\n"+
			"function fn(a, b) {}\n"+
			"goog.inherits(fn, parent);\n"+
			"/** @public */\n"+
			"fn['$inject']=['a', 'b']\n")

	c.Test(
		"/** @ngInject @constructor */\n"+
			"function fn(a, b) {}\n"+
			"goog.inherits(fn, parent);\n"+
			"var foo = 42;\n",
		"/** @ngInject @constructor */\n"+
			"function fn(a, b) {}\n"+
			"goog.inherits(fn, parent);\n"+
			"/** @public */\n"+
			"fn['$inject']=['a', 'b'];\n"+
			"var foo = 42;\n")
}

// Java: AngularPassTest.testNgInjectAddsInjectToProps()
func (c *AngularPassTest) TestNgInjectAddsInjectToProps() {
	c.Test(
		"var ns = {}; /** @ngInject */ ns.fn = function (a, b) {}",
		"var ns = {}; /** @ngInject */ ns.fn = function (a, b) {};\n"+
			"/** @public */ ns.fn['$inject']=['a', 'b']\n")

	c.TestSame("var ns = {}; ns.fn = function (a, b) {}")
}

// Java: AngularPassTest.testNgInjectAddsInjectToNestedProps()
func (c *AngularPassTest) TestNgInjectAddsInjectToNestedProps() {
	c.Test(
		"var ns = {}; ns.subns = {};\n"+
			"/** @ngInject */ ns.subns.fn = function (a, b) {}\n",
		"var ns = {}; ns.subns = {};\n"+
			"/** @ngInject */\n"+
			"ns.subns.fn = function (a, b) {};\n"+
			"/** @public */\n"+
			"ns.subns.fn['$inject']=['a', 'b']\n")

	c.TestSame("var ns = {}; ns.fn = function (a, b) {}")
}

// Java: AngularPassTest.testNgInjectAddsInjectToVars()
func (c *AngularPassTest) TestNgInjectAddsInjectToVars() {
	c.Test(
		"/** @ngInject */ var fn = function (a, b) {}",
		"/** @ngInject */ var fn = function (a, b) {}; /** @public */ fn['$inject']=['a', 'b']")

	c.TestSame("var fn = function (a, b) {}")
}

// Java: AngularPassTest.testNgInjectAddsInjectToLet()
func (c *AngularPassTest) TestNgInjectAddsInjectToLet() {
	c.Test(
		"/** @ngInject */ let fn = function (a, b) {}",
		"/** @ngInject */ let fn = function (a, b) {}; /** @public */ fn['$inject']=['a', 'b']")

	c.TestSame("let fn = function (a, b) {}")
}

// Java: AngularPassTest.testNgInjectAddsInjectToConst()
func (c *AngularPassTest) TestNgInjectAddsInjectToConst() {
	c.Test(
		"/** @ngInject */ const fn = function (a, b) {}",
		"/** @ngInject */ const fn = function (a, b) {}; /** @public */ fn['$inject']=['a', 'b']")

	c.TestSame("const fn = function (a, b) {}")
}

// Java: AngularPassTest.testNgInjectAddsInjectToVarsWithChainedAssignment()
func (c *AngularPassTest) TestNgInjectAddsInjectToVarsWithChainedAssignment() {
	c.Test(
		"var ns = {}; /** @ngInject */ var fn = ns.func = function (a, b) {}",
		"var ns = {};\n"+
			"/** @ngInject */\n"+
			"var fn = ns.func = function (a, b) {};\n"+
			"/** @public */\n"+
			"fn['$inject']=['a', 'b']\n")

	c.TestSame("var ns = {}; var fn = ns.func = function (a, b) {}")
}

// Java: AngularPassTest.testNgInjectInBlock()
func (c *AngularPassTest) TestNgInjectInBlock() {
	c.Test(
		"(function() {\n"+
			"  var ns = {};\n"+
			"  /** @ngInject */ var fn = ns.func = function (a, b) {}\n"+
			"})()\n",
		"(function() {\n"+
			"  var ns = {};\n"+
			"  /** @ngInject */\n"+
			"  var fn = ns.func = function (a, b) {};\n"+
			"  /** @public */\n"+
			"  fn['$inject']=['a', 'b']\n"+
			"})()\n")

	c.TestSame(
		"(function() {\n" +
			"  var ns = {}; var fn = ns.func = function (a, b) {}\n" +
			"})()\n")
}

// Java: AngularPassTest.testNgInjectAddsToTheRightBlock()
func (c *AngularPassTest) TestNgInjectAddsToTheRightBlock() {
	c.Test(
		"var fn = 10;\n"+
			"(function() {\n"+
			"  var ns = {};\n"+
			"  /** @ngInject */ var fn = ns.func = function (a, b) {};\n"+
			"})()\n",
		"var fn = 10;\n"+
			"(function() {\n"+
			"  var ns = {};\n"+
			"  /** @ngInject */\n"+
			"  var fn = ns.func = function (a, b) {};\n"+
			"  /** @public */\n"+
			"  fn['$inject']=['a', 'b'];\n"+
			"})()\n")
}

// Java: AngularPassTest.testNgInjectInNonBlock()
func (c *AngularPassTest) TestNgInjectInNonBlock() {
	c.TestError(
		"function fake(){};\n"+
			"var ns = {};\n"+
			"fake( /** @ngInject */ ns.func = function (a, b) {} )\n",
		AngularPass_INJECT_IN_NON_GLOBAL_OR_BLOCK_ERROR)

	c.TestError(
		"/** @ngInject */( function (a, b) {} )", AngularPass_INJECT_IN_NON_GLOBAL_OR_BLOCK_ERROR)
}

// Java: AngularPassTest.testNgInjectNonFunction()
func (c *AngularPassTest) TestNgInjectNonFunction() {
	c.TestError(
		"var ns = {}; ns.subns = {};\n"+
			"ns.subns.fake = function(x, y){};\n"+
			"/** @ngInject */ ns.subns.fake(1);\n",
		AngularPass_INJECT_NON_FUNCTION_ERROR)

	c.TestError("/** @ngInject */ var a = 10", AngularPass_INJECT_NON_FUNCTION_ERROR)

	c.TestError("/** @ngInject */ var x", AngularPass_INJECT_NON_FUNCTION_ERROR)

	c.TestError(
		"class FnClass {constructor(a, b) {/** @ngInject */ this.x = 42}}",
		AngularPass_INJECT_NON_FUNCTION_ERROR)

	c.TestError(
		"class FnClass {constructor(a, b) {/** @ngInject */ this.x}}",
		AngularPass_INJECT_NON_FUNCTION_ERROR)
}

// Java: AngularPassTest.testNgInjectOnGetElem()
func (c *AngularPassTest) TestNgInjectOnGetElem() {
	c.TestError(
		"/** @ngInject */ foo.bar['baz'] = function(a) {};",
		AngularPass_INJECTED_FUNCTION_ON_NON_QNAME)
}

// Java: AngularPassTest.testNgInjectAddsInjectToClass()
func (c *AngularPassTest) TestNgInjectAddsInjectToClass() {
	c.TestError(
		"/** @ngInject */ class FnClass {constructor(a, b) {}}",
		AngularPass_INJECT_NON_FUNCTION_ERROR)
}

// Java: AngularPassTest.testNgInjectAddsInjectToClassConstructor()
func (c *AngularPassTest) TestNgInjectAddsInjectToClassConstructor() {
	c.Test(
		"class FnClass {/** @ngInject */ constructor(a, b) {}}",
		"class FnClass{ /** @ngInject */ constructor(a, b){}}\n"+
			"/** @public */ FnClass['$inject'] = ['a', 'b'];\n")
}

// Java: AngularPassTest.testNgInjectAddsInjectToClassMethod1()
func (c *AngularPassTest) TestNgInjectAddsInjectToClassMethod1() {
	c.Test(
		"class FnClass {\n"+
			"  constructor(a, b) {}\n"+
			"  /** @ngInject */\n"+
			"  methodA(c, d){}\n"+
			"}\n",
		"class FnClass {\n"+
			"  constructor(a, b){}\n"+
			"  /** @ngInject */\n"+
			"  methodA(c, d){}\n"+
			"}\n"+
			"/** @public */\n"+
			"FnClass.prototype.methodA['$inject'] = ['c','d']\n")
}

// Java: AngularPassTest.testNgInjectAddsInjectToClassMethod2()
func (c *AngularPassTest) TestNgInjectAddsInjectToClassMethod2() {
	c.Test(
		"FnClass.foo = class {\n"+
			"  /** @ngInject */\n"+
			"  constructor(a, b) {}\n"+
			"};\n",
		"FnClass.foo = class {\n"+
			"  /** @ngInject */\n"+
			"  constructor(a, b){}\n"+
			"};\n"+
			"/** @public */\n"+
			"FnClass.foo['$inject'] = ['a','b'];\n")
}

// Java: AngularPassTest.testNgInjectAddsInjectToClassMethod3()
func (c *AngularPassTest) TestNgInjectAddsInjectToClassMethod3() {
	c.Test(
		"var foo = class {\n"+
			"  /** @ngInject */\n"+
			"  constructor(a, b) {}\n"+
			"};\n",
		"var foo = class {\n"+
			"  /** @ngInject */\n"+
			"  constructor(a, b){}\n"+
			"};\n"+
			"/** @public */\n"+
			"foo['$inject'] = ['a','b'];\n")

	c.Test(
		"let foo = class {\n"+
			"  /** @ngInject */\n"+
			"  constructor(a, b) {}\n"+
			"};\n",
		"let foo = class {\n"+
			"  /** @ngInject */\n"+
			"  constructor(a, b){}\n"+
			"};\n"+
			"/** @public */\n"+
			"foo['$inject'] = ['a','b'];\n")

	c.Test(
		"const foo = class {\n"+
			"  /** @ngInject */\n"+
			"  constructor(a, b) {}\n"+
			"};\n",
		"const foo = class {\n"+
			"  /** @ngInject */\n"+
			"  constructor(a, b){}\n"+
			"};\n"+
			"/** @public */\n"+
			"foo['$inject'] = ['a','b'];\n")
}

// Java: AngularPassTest.testNgInjectAddsInjectToStaticMethod()
func (c *AngularPassTest) TestNgInjectAddsInjectToStaticMethod() {
	c.Test(
		"class FnClass {\n"+
			"  constructor(a, b) {}\n"+
			"  /** @ngInject */\n"+
			"  static methodA(c, d) {}\n"+
			"}\n",
		"class FnClass {\n"+
			"  constructor(a, b) {}\n"+
			"  /** @ngInject */\n"+
			"  static methodA(c, d) {}\n"+
			"}\n"+
			"/** @public */\n"+
			"FnClass.methodA['$inject'] = ['c','d']\n")
}

// Java: AngularPassTest.testNgInjectAddsInjectToClassGenerator()
func (c *AngularPassTest) TestNgInjectAddsInjectToClassGenerator() {
	c.Test(
		"class FnClass {\n"+
			"  constructor(a, b) {}\n"+
			"  /** @ngInject */\n"+
			"  * methodA(c, d){}\n"+
			"}\n",
		"class FnClass {\n"+
			"  constructor(a, b){}\n"+
			"  /** @ngInject */\n"+
			"  *methodA(c, d){}\n"+
			"}\n"+
			"/** @public */\n"+
			"FnClass.prototype.methodA['$inject'] = ['c','d']\n")
}

// Java: AngularPassTest.testNgInjectAddsInjectToClassMixOldStyle()
func (c *AngularPassTest) TestNgInjectAddsInjectToClassMixOldStyle() {
	c.Test(
		"class FnClass {\n"+
			"  constructor() {\n"+
			"    /** @ngInject */\n"+
			"    this.someMethod = function(a, b){}\n"+
			"  }\n"+
			"}\n",
		"class FnClass {\n"+
			"  constructor() {\n"+
			"    /** @ngInject */\n"+
			"    this.someMethod = function(a, b){}\n"+
			"    /** @public */\n"+
			"    this.someMethod['$inject'] = ['a','b']\n"+
			"  }\n"+
			"}\n")
}

// Java: AngularPassTest.testNgInjectAddsInjectToClassWithExtraName()
func (c *AngularPassTest) TestNgInjectAddsInjectToClassWithExtraName() {
	c.Test(
		"var foo = class bar{\n"+
			"  /** @ngInject */\n"+
			"  constructor(a, b) {}\n"+
			"};\n",
		"var foo = class bar{\n"+
			"  /** @ngInject */\n"+
			"  constructor(a, b){}\n"+
			"};\n"+
			"/** @public */\n"+
			"foo['$inject'] = ['a','b'];\n")

	c.Test(
		"let foo = class bar{\n"+
			"  /** @ngInject */\n"+
			"  constructor(a, b) {}\n"+
			"};\n",
		"let foo = class bar{\n"+
			"  /** @ngInject */\n"+
			"  constructor(a, b){}\n"+
			"};\n"+
			"/** @public */\n"+
			"foo['$inject'] = ['a','b'];\n")

	c.Test(
		"const foo = class bar{\n"+
			"  /** @ngInject */\n"+
			"  constructor(a, b) {}\n"+
			"};\n",
		"const foo = class bar{\n"+
			"  /** @ngInject */\n"+
			"  constructor(a, b){}\n"+
			"};\n"+
			"/** @public */\n"+
			"foo['$inject'] = ['a','b'];\n")

	c.Test(
		"x.y = class bar{\n"+
			"  /** @ngInject */\n"+
			"  constructor(a, b) {}\n"+
			"};\n",
		"x.y = class bar{\n"+
			"  /** @ngInject */\n"+
			"  constructor(a, b){}\n"+
			"};\n"+
			"/** @public */\n"+
			"x.y['$inject'] = ['a','b'];\n")
}

// Java: AngularPassTest.testNgInjectAddsInjectToClassArrowFunc()
func (c *AngularPassTest) TestNgInjectAddsInjectToClassArrowFunc() {
	c.Test(
		"class FnClass {\n"+
			"  constructor() {\n"+
			"    /** @ngInject */\n"+
			"    this.someMethod = (a, b) => 42\n"+
			"  }\n"+
			"}\n",
		"class FnClass {\n"+
			"  constructor() {\n"+
			"    /** @ngInject */\n"+
			"    this.someMethod = (a, b) => 42\n"+
			"    /** @public */\n"+
			"    this.someMethod['$inject'] = ['a','b']\n"+
			"  }\n"+
			"}\n")
}

// Java: AngularPassTest.testNgInjectAddsInjectToClassCompMethodName()
func (c *AngularPassTest) TestNgInjectAddsInjectToClassCompMethodName() {
	c.TestError(
		"class FnClass {\n"+
			"  constructor() {}\n"+
			"    /** @ngInject */\n"+
			"  ['comp' + 'MethodName'](a, b){}\n"+
			"}\n",
		AngularPass_INJECT_NON_FUNCTION_ERROR)
}

// Java: AngularPassTest.testNgInjectToArrowFunctions()
func (c *AngularPassTest) TestNgInjectToArrowFunctions() {
	c.Test(
		"/** @ngInject */ var fn = (a, b, c)=>{};",
		"/** @ngInject */ var fn = (a, b, c)=>{}; /** @public */ fn['$inject']=['a', 'b', 'c'];")
	c.TestSame("/** @ngInject */ var fn = ()=>{}")
}

// Java: AngularPassTest.testNgInjectToFunctionsWithDestructuredParam()
func (c *AngularPassTest) TestNgInjectToFunctionsWithDestructuredParam() {
	c.TestError(
		"/** @ngInject */ function fn(a, {b, c}){}",
		AngularPass_INJECTED_FUNCTION_HAS_DESTRUCTURED_PARAM)
	c.TestError(
		"/** @ngInject */ function fn(a, [b, c]){}",
		AngularPass_INJECTED_FUNCTION_HAS_DESTRUCTURED_PARAM)
	c.TestError(
		"/** @ngInject */ function fn(a, {b, c}, d){}",
		AngularPass_INJECTED_FUNCTION_HAS_DESTRUCTURED_PARAM)
}

// Java: AngularPassTest.testNgInjectToFunctionsWithDefaultValue()
func (c *AngularPassTest) TestNgInjectToFunctionsWithDefaultValue() {
	c.TestError(
		"/** @ngInject */ function fn(a, b = 1){}",
		AngularPass_INJECTED_FUNCTION_HAS_DEFAULT_VALUE)
	c.TestError(
		"/** @ngInject */ function fn(a, {b, c} = {b: 1, c: 2}){}",
		AngularPass_INJECTED_FUNCTION_HAS_DEFAULT_VALUE)
	c.TestError(
		"/** @ngInject */ function fn(a, [b, c] = [1, 2]){}",
		AngularPass_INJECTED_FUNCTION_HAS_DEFAULT_VALUE)
}

// Java: AngularPassTest.testInGoogModule()
func (c *AngularPassTest) TestInGoogModule() {
	c.EnableRewriteClosureCode()
	c.Test(
		"goog.module('my.module');\n"+
			"/** @ngInject */\n"+
			"function fn(a, b) {}\n",
		"goog.module('my.module');\n"+
			"/** @ngInject */\n"+
			"function fn(a, b) {}\n"+
			"/** @public */ fn['$inject'] = ['a', 'b'];\n")
}

// Java: AngularPassTest.testInEsModule()
func (c *AngularPassTest) TestInEsModule() {
	js :=
		"import {Foo} from './foo';\n" +
			"\n" +
			"class Bar extends Foo { /** @ngInject */ constructor(x, y) {} }\n"
	c.Test(
		js,
		js+
			"/** @public */\n"+
			"Bar['$inject'] = ['x', 'y'];\n")
}

// Java: AngularPassTest.testInExportInEsModule()
func (c *AngularPassTest) TestInExportInEsModule() {
	js :=
		"import {Foo} from './foo';\n" +
			"\n" +
			"export class Bar extends Foo { /** @ngInject */ constructor(x, y) {} }\n"
	c.Test(
		js,
		js+
			"/** @public */\n"+
			"Bar['$inject'] = ['x', 'y'];\n")
}

// Java: AngularPassTest.testInGoogScope()
func (c *AngularPassTest) TestInGoogScope() {
	c.Test(
		"goog.scope(function() {\n"+
			"/** @ngInject */\n"+
			"function fn(a, b) {}\n"+
			"});\n",
		"goog.scope(function() {\n"+
			"/** @ngInject */\n"+
			"function fn(a, b) {}\n"+
			"/** @public */ fn['$inject'] = ['a', 'b'];\n"+
			"});\n")
}

// Java: AngularPassTest.testNameDeclarationAndAssign()
func (c *AngularPassTest) TestNameDeclarationAndAssign() {
	c.Test(
		"let A = A_1 = class A {\n"+
			"  /**\n"+
			"   * @ngInject\n"+
			"   * @param {?} foo\n"+
			"   */\n"+
			"  constructor(foo) {}\n"+
			"};\n",
		"let A = A_1 = class A {\n"+
			"  /**\n"+
			"   * @ngInject\n"+
			"   * @param {?} foo\n"+
			"   */\n"+
			"  constructor(foo) {}\n"+
			"};\n"+
			"/** @public */\n"+
			"A_1[\"$inject\"] = [\"foo\"];\n")
}
