// Port of com/google/javascript/jscomp/RemoveUnusedCodeTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/RemoveUnusedCodeTest.java)

/*
 * Copyright 2004 The Closure Compiler Authors.
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
	"strings"
	"testing"
)

type RemoveUnusedCodeTest struct {
	CompilerTestCaseBase
	removeGlobal                    bool
	preserveFunctionExpressionNames bool
}

// Java: RemoveUnusedCodeTest()
func NewRemoveUnusedCodeTest() *RemoveUnusedCodeTest {
	c := &RemoveUnusedCodeTest{}
	// Set up externs to be used in the test cases.
	c.InitCompilerTestCaseBaseWithExterns(c, RemoveUnusedCodeTest_EXTERNS)
	return c
}

func TestRemoveUnusedCodeTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewRemoveUnusedCodeTest() })
}

// Java: RemoveUnusedCodeTest() (the externs passed to super)
const RemoveUnusedCodeTest_EXTERNS = "var undefined;\n" +
	"var goog = {};\n" +
	"/** @const {!Global} */ goog.global;\n" +
	"goog.reflect = {};\n" +
	"goog.reflect.object = function(obj, propertiesObj) {};\n" +
	"goog.reflect.objectProperty = function(prop, obj) {};\n" +
	"goog.weakUsage = function(nameArg) {};\n" +
	"function goog$inherits(subClass, superClass) {}\n" +
	"function valueType$mixin(dstPrototype, srcPrototype, flags, ...args) {}\n" +
	"function alert() {}\n" +
	"function use() {}\n" +
	"function externFunction() {}\n" +
	"var externVar;\n" +
	"function validate() {}\n" +
	"var userInput;\n" +
	"function gen() {}\n" +
	"var window;\n" +
	"var console = {};\n" +
	"console.log = function(var_args) {};\n" +
	"/** @constructor @return {!Array} */ function Array(/** ...* */ var_args) {}\n" +
	"/** @constructor @return {string} */ function String(/** *= */ opt_arg) {}\n" +
	"/** @constructor */ function Set() {}\n" +
	"/** @constructor */ function WeakMap() {}\n"

const RemoveUnusedCodeTest_JSCOMP_POLYFILL = "var $jscomp = {};\n" +
	"$jscomp.polyfill = function(\n" +
	"    /** string */ name, /** Function */ func, /** string */ from, /** string */ to) {};\n" +
	"$jscomp.polyfillTypedArrayMethod = function(\n" +
	"    /** string */ name, /** Function */ func, /** string */ from, /** string */ to) {};\n"
const RemoveUnusedCodeTest_JSCOMP_PATCH = "var $jscomp = {};\n" +
	"$jscomp.patch = function(/** string */ name, /** Function */ func) {};\n"

// Java: RemoveUnusedCodeTest.setUp()
func (c *RemoveUnusedCodeTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	// Allow testing of features that aren't supported for output yet.
	c.EnableNormalize()
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.EnableGatherExternProperties()
	c.removeGlobal = true
	c.preserveFunctionExpressionNames = false
}

// Java: RemoveUnusedCodeTest.getOptions()
func (c *RemoveUnusedCodeTest) GetOptions() *CompilerOptions {
	return c.CompilerTestCaseBase.GetOptions()
}

// Java: RemoveUnusedCodeTest.getProcessor(Compiler)
func (c *RemoveUnusedCodeTest) GetProcessor(compiler *Compiler) CompilerPass {
	return CompilerPassFunc(func(externs rhino.Node, root rhino.Node) {
		NewRemoveUnusedCode_Builder(compiler).
			RemoveLocalVars(true).
			RemoveGlobals(c.removeGlobal).
			RemoveUnusedPolyfills(true).
			PreserveFunctionExpressionNames(c.preserveFunctionExpressionNames).
			Build().
			Process(externs, root)
	})
}

// Java: RemoveUnusedCodeTest.testDoNotRemoveUnusedVarDeclaredInObjectPatternUsingRest()
func (c *RemoveUnusedCodeTest) TestDoNotRemoveUnusedVarDeclaredInObjectPatternUsingRest() {
	c.TestSame(
		"const data = {\n" +
			"  hello: 'abc',\n" +
			"  world: 'def',\n" +
			"}\n" +
			"\n" +
			"const {\n" +
			"// If `hello` were removed, then `rest` could end up getting a `hello` property.\n" +
			"  hello,\n" +
			"  ...rest\n" +
			"} = data\n" +
			"\n" +
			"console.log(rest);\n")
}

// Java: RemoveUnusedCodeTest.testUnusedPrototypeFieldReference()
func (c *RemoveUnusedCodeTest) TestUnusedPrototypeFieldReference() {
	// Simply mentioning a prototype property without using it doesn't count as a reference.
	c.Test("function C() {} C.prototype.x;", "")
}

// Java: RemoveUnusedCodeTest.testLeaveZeroBehind()
func (c *RemoveUnusedCodeTest) TestLeaveZeroBehind() {
	// We don't need the assignment or the assigned value, but we need to keep the AST valid.
	c.Test(
		"var x; (x = 15, externFunction())", // preserve formatting
		"                externFunction()")
}

// Java: RemoveUnusedCodeTest.testRemoveInBlock()
func (c *RemoveUnusedCodeTest) TestRemoveInBlock() {
	c.Test(
		"if (true) {\n"+
			"  if (true) {\n"+
			"    var foo = function() {};\n"+
			"  }\n"+
			"}\n",
		"if (true) {\n"+
			"  if (true) {\n"+
			"  }\n"+
			"}\n")

	c.Test("if (true) { let foo = function() {} }", "if (true);")
}

// Java: RemoveUnusedCodeTest.testDeclarationInSwitch()
func (c *RemoveUnusedCodeTest) TestDeclarationInSwitch() {
	c.Test(
		"const x = 1;\n"+
			"const y = 2;\n"+
			"switch (x) {\n"+
			"  case 1:\n"+
			"    let y = 3;\n"+
			"    break;\n"+
			"  default:\n"+
			"    let z = 5;\n"+
			"    alert(z);\n"+
			"    break;\n"+
			"}\n"+
			"alert(y);\n",
		"const x = 1;\n"+
			"const y = 2;\n"+
			"switch (x) {\n"+
			"  case 1:\n"+
			"    break;\n"+
			"  default:\n"+
			"    let z = 5;\n"+
			"    alert(z);\n"+
			"    break;\n"+
			"}\n"+
			"alert(y);\n")
}

// Java: RemoveUnusedCodeTest.testPrototypeIsAliased()
func (c *RemoveUnusedCodeTest) TestPrototypeIsAliased() {
	// without alias C is removed
	c.Test("function C() {} C.prototype = {};", "")
	// with alias C must stay
	c.TestSame("function C() {} C.prototype = {}; var x = C.prototype; x;")
	c.TestSame("var x; function C() {} C.prototype = {}; x = C.prototype; x;")
	c.TestSame("var x; function C() {} x = C.prototype = {}; x;")
}

// Java: RemoveUnusedCodeTest.testWholePrototypeAssignment()
func (c *RemoveUnusedCodeTest) TestWholePrototypeAssignment() {
	c.Test("function C() {} C.prototype = { constructor: C };", "")
}

// Java: RemoveUnusedCodeTest.testUsageBeforeDefinition()
func (c *RemoveUnusedCodeTest) TestUsageBeforeDefinition() {
	c.Test("function f(a) { x[a] = 1; } var x; x = {}; f();", "function f() {} f();")
}

// Java: RemoveUnusedCodeTest.testReferencedPropertiesOnUnreferencedVar()
func (c *RemoveUnusedCodeTest) TestReferencedPropertiesOnUnreferencedVar() {
	c.Test(
		"var x = {}; x.a = 1; var y = {a: 2}; y.a;", // preserve format
		"                     var y = {a: 2}; y.a;")
}

// Java: RemoveUnusedCodeTest.testPropertyValuesAddedAfterReferenceAreRemoved()
func (c *RemoveUnusedCodeTest) TestPropertyValuesAddedAfterReferenceAreRemoved() {
	// Make sure property assignments added after the first reference to the var are kept and their
	// values traversed.
	c.TestSame("var x = 1; var y = {}; y; y.foo = x;")
}

// Java: RemoveUnusedCodeTest.testReferenceInObjectLiteral()
func (c *RemoveUnusedCodeTest) TestReferenceInObjectLiteral() {
	c.TestSame(
		"function f(a) {\n" +
			"  return {a: a};\n" +
			"}\n" +
			"f(1);\n")
}

// Java: RemoveUnusedCodeTest.testSelfOverwrite()
func (c *RemoveUnusedCodeTest) TestSelfOverwrite() {
	// Test for possible ConcurrentModificationException
	// Reference to `a` triggers traversal of its function value, which recursively adds another
	// definition of `a` to consider.
	c.TestSame("var a = function() { a = function() {}; }; a();")
}

// Java: RemoveUnusedCodeTest.testPropertyReferenceAddsPropertyReference()
func (c *RemoveUnusedCodeTest) TestPropertyReferenceAddsPropertyReference() {
	// Test for possible ConcurrentModificationException
	// Reference to `a.foo()` triggers traversal of its function value, which recursively adds
	// the `foo` property to another variable.
	c.TestSame("var a = {}; a.foo = function() { b.foo = 1; }; var b = {}; a.foo(); b.foo;")
}

// Java: RemoveUnusedCodeTest.testExternVarDestructuredAssign()
func (c *RemoveUnusedCodeTest) TestExternVarDestructuredAssign() {
	c.TestSame("({a:externVar} = {a:1});")
	c.TestSame("({a:externVar = 1} = {});")
	c.TestSame("({['a']:externVar} = {a:1});")
	c.TestSame("({['a']:externVar = 1} = {});")
	c.TestSame("[externVar] = [1];")
	c.TestSame("[externVar = 1] = [];")
	c.TestSame("[, ...externVar] = [1];")
	c.TestSame("[, ...[externVar = 1]] = [1];")
}

// Java: RemoveUnusedCodeTest.testRemoveVarDeclaration1()
func (c *RemoveUnusedCodeTest) TestRemoveVarDeclaration1() {
	c.Test("var a = 0, b = a = 1", "")
}

// Java: RemoveUnusedCodeTest.testRemoveVarDeclaration2()
func (c *RemoveUnusedCodeTest) TestRemoveVarDeclaration2() {
	c.Test("var a;var b = 0, c = a = b = 1", "")
}

// Java: RemoveUnusedCodeTest.testRemoveUnusedSymbol()
func (c *RemoveUnusedCodeTest) TestRemoveUnusedSymbol() {
	externs := c.Externs("var Symbol = function(opt_desc) {}; var foo = function() {};")
	c.Test(externs, c.Srcs("var x = Symbol();"), c.Expected(""))
	c.Test(externs, c.Srcs("const x = Symbol('desc');"), c.Expected(""))
	c.Test(externs, c.Srcs("let x = Symbol();"), c.Expected(""))
	c.Test(externs, c.Srcs("var x = Symbol(foo());"), c.Expected("Symbol(foo());"))
}

// Java: RemoveUnusedCodeTest.testRemoveUnusedVarsFn0()
func (c *RemoveUnusedCodeTest) TestRemoveUnusedVarsFn0() {
	// Test with function expressions in another function call
	c.Test(
		"function A(){} if(0){function B(){}} window.setTimeout(function(){A()})",
		"function A(){} if(0);                window.setTimeout(function(){A()})")
}

// Java: RemoveUnusedCodeTest.testRemoveUnusedVarsFn1s()
func (c *RemoveUnusedCodeTest) TestRemoveUnusedVarsFn1s() {
	// Test with function expressions in another function call
	c.Test(
		"function A(){}\n"+
			"if(0){var B = function(){}}\n"+
			"A();\n",
		"function A(){}\n"+
			"if(0){}\n"+
			"A();\n")
}

// Java: RemoveUnusedCodeTest.testRemoveUnusedVars1()
func (c *RemoveUnusedCodeTest) TestRemoveUnusedVars1() {
	c.Test(
		"var a;\n"+
			"var b=3;\n"+
			"var c=function(){};\n"+
			"var x=A();\n"+
			"var y;\n"+
			"var z;\n"+
			"function A(){B()}\n"+
			"function B(){C(b)}\n"+
			"function C(){}\n"+
			"function X(){Y()}\n"+
			"function Y(z){Z(x)}\n"+
			"function Z(){y}\n"+
			"externVar = function(){A()};\n"+
			"try{0}catch(e){a}\n",
		"var a;\n"+
			"var b=3;\n"+
			"A();\n"+
			"function A(){B()}\n"+
			"function B(){C(b)}\n"+
			"function C(){}\n"+
			"externVar = function(){A()};\n"+
			"try{0}catch(e){a}\n")

	// Test removal from if {} blocks
	c.Test(
		"var i=0; var j=0; if(i>0){var k=1;}", // preserve alignment
		"var i=0;          if(i>0);")

	// Test with for loop
	c.Test(
		"var x = '';\n"+
			"for (var i in externVar) {\n"+
			"  if (i > 0) x += ', ';\n"+
			"  var arg = 'foo';\n"+
			"  if (arg.length > 40) {\n"+
			"    var unused = 'bar';\n"+
			"    arg = arg.substr(0, 40) + '...';\n"+
			"  }\n"+
			"  x += arg;\n"+
			"}\n"+
			"alert(x);\n",
		"var x = '';\n"+
			"for(var i in externVar) {\n"+
			"  if (i > 0) x += ', ';\n"+
			"  var arg = 'foo';\n"+
			"  if (arg.length > 40) arg = arg.substr(0,40) + '...';\n"+
			"  x += arg\n"+
			"}\n"+
			"alert(x);\n")

	// Test with recursive functions
	c.Test(
		"function A(){A()}function B(){B()}B()", // preserve alignment
		"                 function B(){B()}B()")

	// Test with multiple var declarations.
	c.Test(
		"var x,y=2,z=3; alert(x); alert(z); var a,b,c=4;", // preserve alignment
		"var x,z=3;     alert(x); alert(z);")

	// Test with for loop declarations
	c.Test(
		"for(var i=0,j=0;i<10;){}\n"+
			"for(var x=0,y=0;;y++){}\n"+
			"for(var a,b;;){a}\n"+
			"for(var c,d;;);\n"+
			"for(var item in externVar){}\n",
		"for(var i=0;i<10;);\n"+
			"for(var y=0;;y++);\n"+
			"for(var a;;)a;\n"+
			"for(;;);\n"+
			"for(var item in externVar);\n")

	// Test multiple passes required
	c.Test(
		"var a,b,c,d; var e=[b,c]; var x=e[3]; var f=[d]; alert(f[0])",
		"var       d;                          var f=[d]; alert(f[0])")

	// Test proper scoping (static vs dynamic)
	c.Test(
		"var x; function A() {var x; B()} function B(){ alert(x); } A();",
		"var x; function A() {       B()} function B(){ alert(x); } A();")

	// Test closures in a return statement
	c.TestSame("function A(){var x;return function(){alert(x)}}A()")

	// Test other closures, multiple passes
	c.Test(
		"function A(){}\n"+
			"function B() {\n"+
			"  var c,d,e,f,g,h;\n"+
			"  function C(){ alert(c) }\n"+
			"  var handler = function(){ alert(d) };\n"+
			"  var handler2 = function(){ handler() };\n"+
			"  e = function(){ alert(e) };\n"+
			"  if (1) { function G(){ alert(g) } }\n"+
			"  externVar = [function(){ alert(h) }];\n"+
			"  return function(){ alert(f) };\n"+
			"}\n"+
			"B()\n",
		"function B() {\n"+
			"  var f,h;\n"+
			"  if(1);\n"+
			"  externVar = [function(){ alert(h) }];\n"+
			"  return function(){ alert(f) }\n"+
			"}\n"+
			"B()\n")

	// Test exported names
	c.Test(
		"var a,b=1; function _A1() {this.foo(a)}", // preserve alignment
		"var a;     function _A1() {this.foo(a)}")

	// Test undefined (i.e. externally defined) names
	c.TestSame("externVar = 1")

	// Test unused vars with side effects
	c.Test(
		"var a, b = alert(), c = externVar++,     d; var e = alert(); var f; alert(d);",
		"           alert();     externVar++; var d;         alert();        alert(d);")

	c.Test("var a, b = alert()", "alert()")
	c.Test("var b = alert(), a", "alert()")
	c.Test("var a, b = alert(a)", "var a; alert(a);")
}

// Java: RemoveUnusedCodeTest.testFunctionArgRemoval()
func (c *RemoveUnusedCodeTest) TestFunctionArgRemoval() {
	// remove all function arguments
	c.Test(
		"var b=function(c,d){return};b(1,2)", // preserve alignment
		"var b=function(   ){return};b(1,2)")

	// remove no function arguments
	c.TestSame("var b=function(c,d){return c+d};b(1,2)")
	c.TestSame("var b=function(e,f,c,d){return c+d};b(1,2)")

	// remove some function arguments
	c.Test(
		"var b=function(c,d,e,f){return c+d};b(1,2)", // preserve alignment
		"var b=function(c,d    ){return c+d};b(1,2)")
	c.Test(
		"var b=function(e,c,f,d,g){return c+d};b(1,2)",
		"var b=function(e,c,f,d){return c+d};b(1,2)")
}

// Java: RemoveUnusedCodeTest.testDollarSuperParameterNotRemoved()
func (c *RemoveUnusedCodeTest) TestDollarSuperParameterNotRemoved() {
	// supports special parameter expected by the prototype open-source library
	c.TestSame("function f($super) {} f();")
}

// Java: RemoveUnusedCodeTest.testFunctionArgRemovalWithLeadingUnderscore()
func (c *RemoveUnusedCodeTest) TestFunctionArgRemovalWithLeadingUnderscore() {
	// Coding convention usually prevents removal of variables beginning with a leading underscore,
	// but that makes no sense for parameter names.
	c.Test(
		"function f(__$jscomp$1) {__$jscomp$1 = 1;} f();", // preserve alignment
		"function f(           ) {                } f();")
}

// Java: RemoveUnusedCodeTest.testComputedPropertyDestructuring()
func (c *RemoveUnusedCodeTest) TestComputedPropertyDestructuring() {
	// Don't remove variables accessed by computed properties
	c.TestSame("var {['a']:a, ['b']:b} = {a:1, b:2}; a; b;")

	c.Test(
		"var {['a']:a, ['b']:b} = {a:1, b:2};", // preserve alignment
		"var {                } = {a:1, b:2};")

	c.TestSame("var {[alert()]:a, [alert()]:b} = {a:1, b:2};")

	c.TestSame("var a = {['foo' + 1]:1, ['foo' + 2]:2}; alert(a.foo1);")
}

// Java: RemoveUnusedCodeTest.testFunctionArgRemoval_defaultValue1()
func (c *RemoveUnusedCodeTest) TestFunctionArgRemoval_defaultValue1() {
	c.Test(
		"function f(unusedParam = undefined) {}; f();", // preserve alignment
		"function f(                       ) {}; f();")
}

// Java: RemoveUnusedCodeTest.testFunctionArgRemoval_defaultValue2()
func (c *RemoveUnusedCodeTest) TestFunctionArgRemoval_defaultValue2() {
	c.Test(
		"function f(unusedParam = 0) {}; f();", // preserve alignment
		"function f(               ) {}; f();")
}

// Java: RemoveUnusedCodeTest.testFunctionArgRemoval_defaultValue3()
func (c *RemoveUnusedCodeTest) TestFunctionArgRemoval_defaultValue3() {
	// Parameters already encountered can be used by later parameters
	c.Test(
		"function f(x, y = x + 0) {}; f();", // preserve alignment
		"function f(            ) {}; f();")
}

// Java: RemoveUnusedCodeTest.testFunctionArgRemoval_defaultValue4()
func (c *RemoveUnusedCodeTest) TestFunctionArgRemoval_defaultValue4() {
	// Parameters already encountered can be used by later parameters
	c.TestSame("function f(x, y = x + 0) { y; }; f();")
}

// Java: RemoveUnusedCodeTest.testFunctionArgRemoval_defaultValue5()
func (c *RemoveUnusedCodeTest) TestFunctionArgRemoval_defaultValue5() {
	// Default value inside arrow function param list
	c.Test(
		"var f = (unusedParam = 0) => {}; f();", // preserve alignment
		"var f = (               ) => {}; f();")

	c.TestSame("var f = (usedParam = 0) => { usedParam; }; f();")

	c.Test("var f = (usedParam = 0) => {usedParam;};", "")
}

// Java: RemoveUnusedCodeTest.testFunctionArgRemoval_defaultValue6()
func (c *RemoveUnusedCodeTest) TestFunctionArgRemoval_defaultValue6() {
	// Parameters already encountered can be used by later parameters
	c.TestSame("var x = 2; function f(y = x) { use(y); }; f();")
}

// Java: RemoveUnusedCodeTest.testFunctionArgRemoval_defaultValue7()
func (c *RemoveUnusedCodeTest) TestFunctionArgRemoval_defaultValue7() {
	// Parameters already encountered can be used by later parameters
	c.Test(
		"var x = 2; function f(y = x) {}; f();", // preserve alignment
		"           function f(     ) {}; f();")
}

// Java: RemoveUnusedCodeTest.testDestructuringParams()
func (c *RemoveUnusedCodeTest) TestDestructuringParams() {
	// Default value not used
	c.Test(
		"function f({a:{b:b}} = {a:{}}) { /* b is unused */ }; f();",
		"function f({a:{   }} = {a:{}}) { /* b is unused */ }; f();")

	// Default value with nested value used in default value assignment
	c.Test(
		"function f({a:{b:b}} = {a:{b:1}}) { /* b is unused */ }; f();",
		"function f({a:{   }} = {a:{b:1}}) { /* b is unused */ }; f();")

	// Default value with nested value used in function body
	c.TestSame("function f({a:{b:b}} = {a:{b:1}}) { b; }; f();")

	c.Test(
		"function f({a:{b:b}} = {a:{b:1}}) {}; f();", // preserve alignment
		"function f({a:{   }} = {a:{b:1}}) {}; f();")

	c.Test(
		"function f({a:{b:b}} = {a:{}}) {}; f();", // preserve alignment
		"function f({a:{   }} = {a:{}}) {}; f();")

	// Destructuring pattern not default and parameter not used
	c.Test(
		"function f({a:{b:b}}) { /* b is unused */ }; f({c:{d:1}});",
		"function f({a:{   }}) { /* b is unused */ }; f({c:{d:1}});")

	// Destructuring pattern not default and parameter used
	c.TestSame("function f({a:{b:b}}) { b; }; f({c:{d:1}});")
}

// Java: RemoveUnusedCodeTest.testMixedParamTypes()
func (c *RemoveUnusedCodeTest) TestMixedParamTypes() {
	// Traditional and destructuring pattern
	c.Test(
		"function f({a:{b:b}}, c, d) { c; }; f();", // preserve alignment
		"function f({a:{   }}, c   ) { c; }; f();")

	c.Test(
		"function f({a:{b:b = 5}}, c, d) { c; }; f();", // preserve alignment
		"function f({a:{       }}, c   ) { c; }; f()")

	c.TestSame("function f({}, c, {d:{e:e}}) { c; e; }; f();")

	// Parent is the parameter list
	c.Test(
		"function f({a}, b, {c}) {use(a)}; f({}, {});", // preserve alignment
		"function f({a}        ) {use(a)}; f({}, {});")

	// Default and traditional
	c.Test(
		"function f(unusedParam = undefined, z) { z; }; f();",
		"function f(unusedParam, z            ) { z; }; f();")
}

// Java: RemoveUnusedCodeTest.testDefaultParams()
func (c *RemoveUnusedCodeTest) TestDefaultParams() {
	c.Test(
		"function f(x = undefined) {}; f();", // preserve alignment
		"function f(             ) {}; f();")
	c.Test(
		"function f(x = 0) {}; f();", // preserve alignment
		"function f(     ) {}; f();")
	c.TestSame("function f(x = 0, y = x) { alert(y); }; f();")
}

// Java: RemoveUnusedCodeTest.testDefaultParamsInClass()
func (c *RemoveUnusedCodeTest) TestDefaultParamsInClass() {
	c.Test(
		"class Foo { constructor(value = undefined) {} }; new Foo;",
		"class Foo { constructor(                 ) {} }; new Foo;")

	c.TestSame("class Foo { constructor(value = undefined) { value; } }; new Foo;")
	c.TestSame(
		"class Bar {}\n" +
			"class Foo extends Bar {\n" +
			"  constructor(value = undefined) { super(); value; }\n" +
			"}\n" +
			"new Foo;\n")
}

// Java: RemoveUnusedCodeTest.testDefaultParamsInClassThatReferencesArguments()
func (c *RemoveUnusedCodeTest) TestDefaultParamsInClassThatReferencesArguments() {
	c.TestSame(
		"class Bar {}\n" +
			"class Foo extends Bar {\n" +
			"  constructor(value = undefined) {\n" +
			"    super();\n" +
			"    if (arguments.length)\n" +
			"      value;\n" +
			"  }\n" +
			"};\n" +
			"new Foo;\n")
}

// Java: RemoveUnusedCodeTest.testDefaultParamsWithoutSideEffects0()
func (c *RemoveUnusedCodeTest) TestDefaultParamsWithoutSideEffects0() {
	c.Test(
		"function f({} = {}){}; f()", // preserve alignment
		"function f(       ){}; f()")
}

// Java: RemoveUnusedCodeTest.testDefaultParamsWithoutSideEffects1()
func (c *RemoveUnusedCodeTest) TestDefaultParamsWithoutSideEffects1() {
	c.Test(
		"function f(a = 1) {}; f();", // preserve alignment
		"function f(     ) {}; f();")

	c.Test(
		"function f({a:b = 1} = 1){}; f();", // preserve alignment
		"function f(             ){}; f();")

	c.Test(
		"function f({a:b} = {}){}; f();", // preserve alignment
		"function f(          ){}; f();")
}

// Java: RemoveUnusedCodeTest.testDefaultParamsWithSideEffects1()
func (c *RemoveUnusedCodeTest) TestDefaultParamsWithSideEffects1() {
	c.TestSame("function f(a = alert('foo')) {}; f();")

	c.TestSame("function f({} = alert('foo')){}; f()")

	c.Test(
		"function f(){var x; var {a:b} = x()}; f();", // preserve alignment
		"function f(){var x; var {   } = x()}; f();")

	c.Test(
		"function f(){var {a:b} = alert('foo')}; f();",
		"function f(){var {} = alert('foo')}; f();")
	c.Test(
		"function f({a:b} = alert('foo')){}; f();", // preserve alignment
		"function f({   } = alert('foo')){}; f();")

	c.TestSame("function f({a:b = alert('bar')} = alert('foo')){}; f();")
}

// Java: RemoveUnusedCodeTest.testDefaultParamsWithSideEffects2()
func (c *RemoveUnusedCodeTest) TestDefaultParamsWithSideEffects2() {
	c.Test(
		"function f({a:b} = alert('foo')){}; f();", // preserve alignment
		"function f({   } = alert('foo')){}; f();")
}

// Java: RemoveUnusedCodeTest.testUnusedSetterParam_isRetained()
func (c *RemoveUnusedCodeTest) TestUnusedSetterParam_isRetained() {
	// These params are a syntactic requirement.
	c.TestSame("class Foo { set foo(x)     { } }; use(new Foo);")
	c.TestSame("class Foo { set ['foo'](x) { } }; use(new Foo);")
	c.TestSame("class Foo { set 'foo'(x)   { } }; use(new Foo);")
}

// Java: RemoveUnusedCodeTest.testArrayDestructuringParams()
func (c *RemoveUnusedCodeTest) TestArrayDestructuringParams() {
	// Default values in array pattern unused
	c.Test(
		"function f([x,y] = [1,2]) {}; f();", // preserve alignment
		"function f(             ) {}; f();")

	c.Test(
		"function f([x,y] = [1,2], z) { z; }; f();", // preserve alignment
		"function f([   ] = [1,2], z) { z; }; f();")

	// Default values in array pattern used
	c.Test(
		"function f([x,y,z] = [1,2,3]) { y; }; f();", // preserve alignment
		"function f([ ,y  ] = [1,2,3]) { y; }; f();")

	// Side effects
	c.Test(
		"function f([x] = [alert()]) {}; f();", // preserve alignment
		"function f([ ] = [alert()]) {}; f();")
}

// Java: RemoveUnusedCodeTest.testRestPattern()
func (c *RemoveUnusedCodeTest) TestRestPattern() {
	c.TestSame("var x; [...x] = externVar;")
	c.TestSame("var [...x] = externVar;")
	c.TestSame("var x; [...x] = externVar; use(x);")
	c.TestSame("var [...x] = externVar; use(x);")
	c.TestSame("var x; [...x.y] = externVar;")
	c.TestSame("var x; [...x['y']] = externVar;")
	c.TestSame("var x; [...x().y] = externVar;")
	c.TestSame("var x; [...x()['y']] = externVar;")
}

// Java: RemoveUnusedCodeTest.testRestParams()
func (c *RemoveUnusedCodeTest) TestRestParams() {
	c.Test(
		"function foo(...args) {/* rest param unused*/}; foo();",
		"function foo(       ) {/* rest param unused*/}; foo();")

	c.TestSame("function foo(a, ...args) { args[0]; }; foo();")

	// Rest param in pattern
	c.TestSame(
		"function countArgs(x, ...{length}) {\n" +
			"  return length;\n" +
			"}\n" +
			"alert(countArgs(1, 1, 1, 1, 1));\n")

	c.TestSame("function foo([...rest]) {/* rest unused*/}; foo();")

	c.TestSame("function foo([x, ...rest]) { x; }; foo();")
}

// Java: RemoveUnusedCodeTest.testFunctionsDeadButEscaped()
func (c *RemoveUnusedCodeTest) TestFunctionsDeadButEscaped() {
	c.TestSame("function b(a) { a = 1; alert(arguments[0]) }; b(6)")
	c.TestSame("function b(a) { var c = 2; a = c; alert(arguments[0]) }; b(6)")
}

// Java: RemoveUnusedCodeTest.testVarInControlStructure()
func (c *RemoveUnusedCodeTest) TestVarInControlStructure() {
	c.Test("if (true) var b = 3;", "if(true);")
	c.Test("if (true) var b = 3; else var c = 5;", "if(true);else;")
	c.Test("while (true) var b = 3;", "while(true);")
	c.Test("for (;;) var b = 3;", "for(;;);")
	c.Test("do var b = 3; while(true)", "do;while(true)")
	c.Test("with (true) var b = 3;", "with(true);")
	c.Test("f: var b = 3;", "f:{}")
}

// Java: RemoveUnusedCodeTest.testRValueHoisting()
func (c *RemoveUnusedCodeTest) TestRValueHoisting() {
	c.Test("var x = alert();", "alert()")
	c.Test("var x = {a: alert()};", "({a:alert()})")

	c.Test("var x=function y(){}", "")
}

// Java: RemoveUnusedCodeTest.testModule()
func (c *RemoveUnusedCodeTest) TestModule() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk(
					"var unreferenced=1; function x() { foo(); }\n"+
						"function uncalled() { var x; return 2; }\n").
				AddChunk("var a,b; function foo() { this.foo(a); } x()").
				Build()),
		c.Expected("function x(){foo()}", "var a;function foo(){this.foo(a)}x()"))
}

// Java: RemoveUnusedCodeTest.testRecursiveFunction1()
func (c *RemoveUnusedCodeTest) TestRecursiveFunction1() {
	c.TestSame("(function x(){return x()})()")
}

// Java: RemoveUnusedCodeTest.testRecursiveFunction2()
func (c *RemoveUnusedCodeTest) TestRecursiveFunction2() {
	c.Test(
		"var x = 3; (function          x() { return          x(); })();",
		"           (function x$jscomp$1() { return x$jscomp$1(); })()")
}

// Java: RemoveUnusedCodeTest.testFunctionWithName1()
func (c *RemoveUnusedCodeTest) TestFunctionWithName1() {
	c.Test(
		"var x=function f(){};x()", // preserve alignment
		"var x=function  (){};x()")

	c.preserveFunctionExpressionNames = true
	c.TestSame("var x=function f(){};x()")
}

// Java: RemoveUnusedCodeTest.testFunctionWithName2()
func (c *RemoveUnusedCodeTest) TestFunctionWithName2() {
	c.Test("alert(function bar(){})", "alert(function(){})")

	c.preserveFunctionExpressionNames = true
	c.TestSame("alert(function bar(){})")
}

// Java: RemoveUnusedCodeTest.testRemoveGlobal1()
func (c *RemoveUnusedCodeTest) TestRemoveGlobal1() {
	c.removeGlobal = false
	c.TestSame("var x=1")
	c.Test("var y=function(x){var z;}", "var y=function(x){}")
}

// Java: RemoveUnusedCodeTest.testRemoveGlobal2()
func (c *RemoveUnusedCodeTest) TestRemoveGlobal2() {
	c.removeGlobal = false
	c.TestSame("var x=1")
	c.Test("function y(x){var z;}", "function y(x){}")
}

// Java: RemoveUnusedCodeTest.testRemoveGlobal3()
func (c *RemoveUnusedCodeTest) TestRemoveGlobal3() {
	c.removeGlobal = false
	c.TestSame("var x=1")
	c.Test(
		"function x(){function y(x){var z;}y()}", // preserve alignment
		"function x(){function y(x){      }y()}")
}

// Java: RemoveUnusedCodeTest.testRemoveGlobal4()
func (c *RemoveUnusedCodeTest) TestRemoveGlobal4() {
	c.removeGlobal = false
	c.TestSame("var x=1")
	c.Test(
		"function x(){function y(x){var z;}}", // preserve alignment
		"function x(){                     }")
}

// Java: RemoveUnusedCodeTest.testIssue168a()
func (c *RemoveUnusedCodeTest) TestIssue168a() {
	c.Test(
		"function _a(){\n"+
			"  (function(x){ _b(); })(1);\n"+
			"}\n"+
			"function _b(){\n"+
			"  _a();\n"+
			"}\n",
		"function _a(){(function(){_b()})(1)}\n"+
			"function _b(){_a()}\n")
}

// Java: RemoveUnusedCodeTest.testIssue168b()
func (c *RemoveUnusedCodeTest) TestIssue168b() {
	c.removeGlobal = false
	c.Test(
		"function a(){\n"+
			"  (function(x){ b(); })(1);\n"+
			"}\n"+
			"function b(){\n"+
			"  a();\n"+
			"}\n",
		"function a(){(function(x){b()})(1)}\n"+
			"function b(){a()}\n")
}

// Java: RemoveUnusedCodeTest.testUnusedAssign1()
func (c *RemoveUnusedCodeTest) TestUnusedAssign1() {
	c.Test("var x = 3; x = 5;", "")
}

// Java: RemoveUnusedCodeTest.testUnusedAssign2()
func (c *RemoveUnusedCodeTest) TestUnusedAssign2() {
	c.Test(
		"function f(a) { a = 3; } this.x = f;", //
		"function f(){} this.x=f")
}

// Java: RemoveUnusedCodeTest.testUnusedAssign3()
func (c *RemoveUnusedCodeTest) TestUnusedAssign3() {
	// e can't be removed, so we don't try to remove the dead assign.
	// We might be able to improve on this case.
	c.Test(
		"try { throw ''; } catch (e) { e = 3; }", //
		"try{throw\"\";}catch(e){e=3}")
}

// Java: RemoveUnusedCodeTest.testUnusedAssign4()
func (c *RemoveUnusedCodeTest) TestUnusedAssign4() {
	c.Test(
		"function f(a, b) { this.foo(b); a = 3; } this.x = f;",
		"function f(a,b){this.foo(b);}this.x=f")
}

// Java: RemoveUnusedCodeTest.testUnusedAssign5()
func (c *RemoveUnusedCodeTest) TestUnusedAssign5() {
	c.Test(
		"var z = function f() { f = 3; }; z();", //
		"var z=function(){};z()")
}

// Java: RemoveUnusedCodeTest.testUnusedAssign5b()
func (c *RemoveUnusedCodeTest) TestUnusedAssign5b() {
	c.Test(
		"var z = function f() { f = alert(); }; z();", //
		"var z=function(){alert()};z()")
}

// Java: RemoveUnusedCodeTest.testUnusedAssign6()
func (c *RemoveUnusedCodeTest) TestUnusedAssign6() {
	c.Test("var z; z = 3;", "")
}

// Java: RemoveUnusedCodeTest.testUnusedAssign6b()
func (c *RemoveUnusedCodeTest) TestUnusedAssign6b() {
	c.Test("var z; z = alert();", "alert()")
}

// Java: RemoveUnusedCodeTest.testUnusedAssign_toString()
func (c *RemoveUnusedCodeTest) TestUnusedAssign_toString() {
	c.Test("var x = (10).toString(16);", "")
	c.Test("var x = (10).toString();", "")
	c.Test("var x = (10).toString(16, 2);", "(10).toString(16, 2);")
}

// Java: RemoveUnusedCodeTest.testUnusedAssign7()
func (c *RemoveUnusedCodeTest) TestUnusedAssign7() {
	// This loop is normalized to "var i;for(i in..."
	c.Test(
		"var a = 3; for (var i in {}) { i = a; }",
		// TODO(johnlenz): "i = a" should be removed here.
		"var a = 3; var i; for (i in {}) {i = a;}")
}

// Java: RemoveUnusedCodeTest.testUnusedAssign8()
func (c *RemoveUnusedCodeTest) TestUnusedAssign8() {
	// This loop is normalized to "var i;for(i in..."
	c.Test(
		"var a = 3; for (var i in {}) { i = a; } alert(a);",
		// TODO(johnlenz): "i = a" should be removed here.
		"var a = 3; var i; for (i in {}) {i = a} alert(a);")
}

// Java: RemoveUnusedCodeTest.testUnusedAssign9a()
func (c *RemoveUnusedCodeTest) TestUnusedAssign9a() {
	c.TestSame("function b(a) { a = 1; arguments; }; b(6)")
}

// Java: RemoveUnusedCodeTest.testUnusedAssign9()
func (c *RemoveUnusedCodeTest) TestUnusedAssign9() {
	c.TestSame("function b(a) { a = 1; arguments=1; }; b(6)")
}

// Java: RemoveUnusedCodeTest.testES6ModuleExports()
func (c *RemoveUnusedCodeTest) TestES6ModuleExports() {
	c.Test("const X = 1; function f() {}", "")
	c.Test("const X = 1; export function f() {}", "function f() {} export {f as f}")
	c.Test("const X = 1; export class C {}", "class C {} export { C as C }")
	c.Test("const X = 1; export default function f() {};", "export default function f() {}")
	c.Test("const X = 1; export default class C {}", "export default class C {}")
}

// Java: RemoveUnusedCodeTest.testUnusedPropAssign1()
func (c *RemoveUnusedCodeTest) TestUnusedPropAssign1() {
	c.Test("var x = {}; x.foo = 3;", "")
}

// Java: RemoveUnusedCodeTest.testUnusedPropAssign1b()
func (c *RemoveUnusedCodeTest) TestUnusedPropAssign1b() {
	c.Test("var x = {}; x.foo = alert();", "alert()")
}

// Java: RemoveUnusedCodeTest.testUnusedPropAssign2()
func (c *RemoveUnusedCodeTest) TestUnusedPropAssign2() {
	c.Test("var x = {}; x['foo'] = 3;", "")
}

// Java: RemoveUnusedCodeTest.testUnusedPropAssign2b()
func (c *RemoveUnusedCodeTest) TestUnusedPropAssign2b() {
	c.Test("var x = {}; x[alert()] = alert();", "alert(),alert()")
}

// Java: RemoveUnusedCodeTest.testUnusedPropAssign3()
func (c *RemoveUnusedCodeTest) TestUnusedPropAssign3() {
	c.Test("var x = {}; x['foo'] = {}; x['bar'] = 3", "")
}

// Java: RemoveUnusedCodeTest.testUnusedPropAssign3b()
func (c *RemoveUnusedCodeTest) TestUnusedPropAssign3b() {
	c.Test(
		"var x = {}; x[alert()] = alert(); x[alert() + alert()] = alert()",
		"alert(),alert();(alert() + alert()),alert()")
}

// Java: RemoveUnusedCodeTest.testUnusedPropAssign4()
func (c *RemoveUnusedCodeTest) TestUnusedPropAssign4() {
	c.Test("var x = {foo: 3}; x['foo'] = 5;", "")
}

// Java: RemoveUnusedCodeTest.testUnusedPropAssign5()
func (c *RemoveUnusedCodeTest) TestUnusedPropAssign5() {
	c.Test(
		"var x = {foo: alert()}; x['foo'] = 5;", // preserve newline
		"       ({foo: alert()})             ;")
}

// Java: RemoveUnusedCodeTest.testUnusedPropAssign6()
func (c *RemoveUnusedCodeTest) TestUnusedPropAssign6() {
	c.Test("var x = function() {}; x.prototype.bar = function() {};", "")
}

// Java: RemoveUnusedCodeTest.testUnusedPropAssign6b()
func (c *RemoveUnusedCodeTest) TestUnusedPropAssign6b() {
	c.Test("function x() {} x.prototype.bar = function() {};", "")
}

// Java: RemoveUnusedCodeTest.testUnusedPropAssign6c()
func (c *RemoveUnusedCodeTest) TestUnusedPropAssign6c() {
	c.Test("function x() {} x.prototype['bar'] = function() {};", "")
}

// Java: RemoveUnusedCodeTest.testUnusedPropAssign7()
func (c *RemoveUnusedCodeTest) TestUnusedPropAssign7() {
	c.Test("var x = {}; x[x.foo] = x.bar;", "")
}

// Java: RemoveUnusedCodeTest.testUnusedPropAssign7b()
func (c *RemoveUnusedCodeTest) TestUnusedPropAssign7b() {
	c.TestSame("var x = {}; x[x.foo] = alert(x.bar);")
}

// Java: RemoveUnusedCodeTest.testUnusedPropAssign7c()
func (c *RemoveUnusedCodeTest) TestUnusedPropAssign7c() {
	c.Test(
		"var x = {}; x[alert(x.foo)] = x.bar;", //
		"var x={};x[alert(x.foo)]=x.bar")
}

// Java: RemoveUnusedCodeTest.testUsedPropAssign1()
func (c *RemoveUnusedCodeTest) TestUsedPropAssign1() {
	c.TestSame("function f(x) { x.bar = 3; } f({});")
}

// Java: RemoveUnusedCodeTest.testUsedPropAssign2()
func (c *RemoveUnusedCodeTest) TestUsedPropAssign2() {
	c.TestSame("try { throw {}; } catch (e) { e.bar = 3; }")
}

// Java: RemoveUnusedCodeTest.testUsedPropInDestructuringCatch_objectPattern()
func (c *RemoveUnusedCodeTest) TestUsedPropInDestructuringCatch_objectPattern() {
	c.TestSame("this.p = 1; try {} catch ({p}) { alert(p); }")
}

// Java: RemoveUnusedCodeTest.testUsedPropInDestructuringCatch_arrayPattern()
func (c *RemoveUnusedCodeTest) TestUsedPropInDestructuringCatch_arrayPattern() {
	c.TestSame("this.p = 1; try {} catch ([{p}]) { alert(p); }")
}

// Java: RemoveUnusedCodeTest.testUsedPropAssign3()
func (c *RemoveUnusedCodeTest) TestUsedPropAssign3() {
	// This pass does not do flow analysis.
	c.TestSame("var x = {}; x.foo = 3; x = alert();")
}

// Java: RemoveUnusedCodeTest.testUsedPropAssign4()
func (c *RemoveUnusedCodeTest) TestUsedPropAssign4() {
	c.TestSame("var y = alert(); var x = {}; x.foo = 3; y[x.foo] = 5;")
}

// Java: RemoveUnusedCodeTest.testUsedPropAssign5()
func (c *RemoveUnusedCodeTest) TestUsedPropAssign5() {
	c.TestSame("var y = alert(); var x = 3; y[x] = 5;")
}

// Java: RemoveUnusedCodeTest.testUsedPropAssign6()
func (c *RemoveUnusedCodeTest) TestUsedPropAssign6() {
	c.TestSame("var x = alert(externVar); x.innerHTML = 'new text';")
}

// Java: RemoveUnusedCodeTest.testUsedPropAssign7()
func (c *RemoveUnusedCodeTest) TestUsedPropAssign7() {
	c.TestSame("var x = {}; for (x in alert()) { x.foo = 3; }")
}

// Java: RemoveUnusedCodeTest.testUsedPropAssign8()
func (c *RemoveUnusedCodeTest) TestUsedPropAssign8() {
	c.TestSame("for (var x in alert()) { x.foo = 3; }")
}

// Java: RemoveUnusedCodeTest.testUsedPropAssign9()
func (c *RemoveUnusedCodeTest) TestUsedPropAssign9() {
	c.TestSame("var x = {}; x.foo = alert(externVar); x.foo.innerHTML = 'new test';")
}

// Java: RemoveUnusedCodeTest.testUsedPropNotAssign()
func (c *RemoveUnusedCodeTest) TestUsedPropNotAssign() {
	c.TestSame("function f(x) { x.a; }; f(externVar);")
}

// Java: RemoveUnusedCodeTest.testDependencies1()
func (c *RemoveUnusedCodeTest) TestDependencies1() {
	c.Test("var a = 3; var b = function() { alert(a); };", "")
}

// Java: RemoveUnusedCodeTest.testDependencies1b()
func (c *RemoveUnusedCodeTest) TestDependencies1b() {
	c.Test(
		"var a = 3; var b = alert(function() { alert(a); });",
		"var a=3;alert(function(){alert(a)})")
}

// Java: RemoveUnusedCodeTest.testDependencies1c()
func (c *RemoveUnusedCodeTest) TestDependencies1c() {
	c.Test(
		"var a = 3; var _b = function() { alert(a); };", //
		"var a=3;var _b=function(){alert(a)}")
}

// Java: RemoveUnusedCodeTest.testDependencies2()
func (c *RemoveUnusedCodeTest) TestDependencies2() {
	c.Test("var a = 3; var b = 3; b = function() { alert(a); };", "")
}

// Java: RemoveUnusedCodeTest.testDependencies2b()
func (c *RemoveUnusedCodeTest) TestDependencies2b() {
	c.Test(
		"var a = 3; var b = 3; b = alert(function() { alert(a); });",
		"var a=3;alert(function(){alert(a)})")
}

// Java: RemoveUnusedCodeTest.testDependencies2c()
func (c *RemoveUnusedCodeTest) TestDependencies2c() {
	c.TestSame("var a=3;var _b=3;_b=function(){alert(a)}")
}

// Java: RemoveUnusedCodeTest.testGlobalVarReferencesLocalVar()
func (c *RemoveUnusedCodeTest) TestGlobalVarReferencesLocalVar() {
	c.TestSame("var a=3;function f(){var b=4;a=b}alert(a + f())")
}

// Java: RemoveUnusedCodeTest.testLocalVarReferencesGlobalVar1()
func (c *RemoveUnusedCodeTest) TestLocalVarReferencesGlobalVar1() {
	c.TestSame("var a=3;function f(b, c){b=a; alert(b + c);} f();")
}

// Java: RemoveUnusedCodeTest.testLocalVarReferencesGlobalVar2()
func (c *RemoveUnusedCodeTest) TestLocalVarReferencesGlobalVar2() {
	c.Test(
		"var a=3;function f(b, c){b=a; alert(c);} f();", //
		"function f(b, c) { alert(c); } f();")
}

// Java: RemoveUnusedCodeTest.testNestedAssign1()
func (c *RemoveUnusedCodeTest) TestNestedAssign1() {
	c.Test(
		"var b = null; var a = (b = 3); alert(a);", //
		"var a = 3; alert(a);")
}

// Java: RemoveUnusedCodeTest.testNestedAssign2()
func (c *RemoveUnusedCodeTest) TestNestedAssign2() {
	c.Test(
		"var a = 1; var b = 2; var c = (b = a); alert(c);", //
		"var a = 1; var c = a; alert(c);")
}

// Java: RemoveUnusedCodeTest.testNestedAssign3()
func (c *RemoveUnusedCodeTest) TestNestedAssign3() {
	c.Test(
		"var b = 0; var z; z = z = b = 1; alert(b);", //
		"var b = 0; b = 1; alert(b);")
}

// Java: RemoveUnusedCodeTest.testWeakUsageRemoved()
func (c *RemoveUnusedCodeTest) TestWeakUsageRemoved() {
	// There are no other references to `a`, so it can be removed.
	c.Test(
		"var a=function(){return}; use(goog.weakUsage(a));", //
		"use(void 0);")
}

// Java: RemoveUnusedCodeTest.testWeakUsageNotRemoved()
func (c *RemoveUnusedCodeTest) TestWeakUsageNotRemoved() {
	// There is another reference to `a`, so it cannot be removed.
	//
	// Normally the PeepholeReplaceKnownMethods pass would remove the call to `goog.weakUsage`
	// entirely, but in this unit test we are only testing RemoveUnusedCode.
	c.TestSame("var a=function(){return}; use(goog.weakUsage(a)); use(a);")
}

// Java: RemoveUnusedCodeTest.testCallSiteInteraction()
func (c *RemoveUnusedCodeTest) TestCallSiteInteraction() {
	c.TestSame("var b=function(){return};b()")
	c.TestSame("var b=function(c){return c};b(1)")
	c.Test(
		"var b=function(c){return};b(1)", //
		"var b=function(){return};b(1)")
	c.Test(
		"var b=function(c){};b.call(null, externVar)",
		"var b=function( ){};b.call(null, externVar)")
	c.Test(
		"var b=function(c){};b.apply(null, externVar)",
		"var b=function( ){};b.apply(null, externVar)")

	// Recursive calls
	c.TestSame("var b=function(c,d){b(1, 2);return d};b(3,4);b(5,6)")

	c.TestSame("var b=function(c){return arguments};b(1,2);b(3,4)")
}

// Java: RemoveUnusedCodeTest.testCallSiteInteraction_constructors()
func (c *RemoveUnusedCodeTest) TestCallSiteInteraction_constructors() {
	// The third level tests that the functions which have already been looked
	// at get re-visited if they are changed by a call site removal.
	c.Test(
		"var Ctor1=function(a, b){return a};\n"+
			"var Ctor2=function(x, y){Ctor1.call(this, x, y)};\n"+
			"goog$inherits(Ctor2, Ctor1);\n"+
			"new Ctor2(1, 2)\n",
		"var Ctor1=function(a){return a};\n"+
			"var Ctor2=function(x, y){Ctor1.call(this, x, y)};\n"+
			"goog$inherits(Ctor2, Ctor1);\n"+
			"new Ctor2(1, 2)\n")
}

// Java: RemoveUnusedCodeTest.testRemoveUnusedVarsPossibleNpeCase()
func (c *RemoveUnusedCodeTest) TestRemoveUnusedVarsPossibleNpeCase() {
	c.Test(
		"var a = [];\n"+
			"var register = function(callback) {a[0] = callback};\n"+
			"register(function(transformer) {});\n"+
			"register(function(transformer) {});\n",
		"var register = function() {};\n"+
			"register(function() {});\n"+
			"register(function() {});\n")
}

// Java: RemoveUnusedCodeTest.testDoNotOptimizeJSCompiler_renameProperty()
func (c *RemoveUnusedCodeTest) TestDoNotOptimizeJSCompiler_renameProperty() {
	// Only the function definition can be modified, none of the call sites.
	c.Test(
		"function JSCompiler_renameProperty(a) {};\n"+
			"JSCompiler_renameProperty('a');\n",
		"function JSCompiler_renameProperty() {};\n"+
			"JSCompiler_renameProperty('a');\n")
}

// Java: RemoveUnusedCodeTest.testDoNotOptimizeSetters()
func (c *RemoveUnusedCodeTest) TestDoNotOptimizeSetters() {
	c.TestSame("({set s(a) {}})")
}

// Java: RemoveUnusedCodeTest.testRemoveSingletonClass1()
func (c *RemoveUnusedCodeTest) TestRemoveSingletonClass1() {
	c.Test(
		"function goog$addSingletonGetter(a){}\n"+
			"/**@constructor*/function a(){}\n"+
			"goog$addSingletonGetter(a);\n",
		"")
}

// Java: RemoveUnusedCodeTest.testRemoveInheritedClass1()
func (c *RemoveUnusedCodeTest) TestRemoveInheritedClass1() {
	c.Test(
		"/**@constructor*/ function a() {}\n"+
			"/**@constructor*/ function b() {}\n"+
			"goog$inherits(b, a);\n"+
			"new a;\n",
		"/**@constructor*/ function a() {}\n"+
			"new a;\n")
}

// Java: RemoveUnusedCodeTest.testRemoveInheritedClass2()
func (c *RemoveUnusedCodeTest) TestRemoveInheritedClass2() {
	c.Test(
		"/**@constructor*/function a(){}\n"+
			"/**@constructor*/function b(){}\n"+
			"/**@constructor*/function c(){}\n"+
			"goog$inherits(b,a);\n",
		"")
}

// Java: RemoveUnusedCodeTest.testRemoveInheritedClass3()
func (c *RemoveUnusedCodeTest) TestRemoveInheritedClass3() {
	c.TestSame(
		"/**@constructor*/function a(){}\n" +
			"/**@constructor*/function b(){}\n" +
			"goog$inherits(b,a); new b\n")
}

// Java: RemoveUnusedCodeTest.testRemoveInheritedClass4()
func (c *RemoveUnusedCodeTest) TestRemoveInheritedClass4() {
	c.TestSame(
		"function goog$inherits(){}\n" +
			"/**@constructor*/function a(){}\n" +
			"/**@constructor*/function b(){}\n" +
			"goog$inherits(b,a);\n" +
			"/**@constructor*/function c(){}\n" +
			"goog$inherits(c,b); new c\n")
}

// Java: RemoveUnusedCodeTest.testRemoveInheritedClass5()
func (c *RemoveUnusedCodeTest) TestRemoveInheritedClass5() {
	c.Test(
		"function goog$inherits() {}\n"+
			"/**@constructor*/ function a() {}\n"+
			"/**@constructor*/ function b() {}\n"+
			"goog$inherits(b,a);\n"+
			"/**@constructor*/ function c() {}\n"+
			"goog$inherits(c,b); new b\n",
		"function goog$inherits(){}\n"+
			"/**@constructor*/ function a(){}\n"+
			"/**@constructor*/ function b(){}\n"+
			"goog$inherits(b,a); new b\n")
}

// Java: RemoveUnusedCodeTest.testRemoveInheritedClass8()
func (c *RemoveUnusedCodeTest) TestRemoveInheritedClass8() {
	c.TestSame(
		"/**@constructor*/function a(){}\n" +
			"/**@constructor*/function b(){}\n" +
			"/**@constructor*/function c(){}\n" +
			"b.inherits(a);c.mixin(b.prototype);new c\n")
}

// Java: RemoveUnusedCodeTest.testRemoveInheritedClass9()
func (c *RemoveUnusedCodeTest) TestRemoveInheritedClass9() {
	c.Test(
		"function goog$inherits() {}\n"+
			"/**@constructor*/ function a() {}\n"+
			"/**@constructor*/ function b() {}\n"+
			"goog$inherits(b,a); new a;\n"+
			"var c = a; var d = a.g; new b\n",
		"function goog$inherits(){}\n"+
			"/**@constructor*/ function a(){}\n"+
			"/**@constructor*/ function b(){}\n"+
			"goog$inherits(b,a);\n"+
			"new a; new b\n")
}

// Java: RemoveUnusedCodeTest.testRemoveInheritedClass11()
func (c *RemoveUnusedCodeTest) TestRemoveInheritedClass11() {
	c.TestSame(
		"/**@constructor*/function a(){}\n" +
			"var b = {};\n" +
			"// goog$inherits not treated specially when derived class is a property\n" +
			"goog$inherits(b.foo, a)\n")
}

// Java: RemoveUnusedCodeTest.testRemoveInheritedClass12()
func (c *RemoveUnusedCodeTest) TestRemoveInheritedClass12() {
	// An inherits call must be a statement unto itself or the left side of a comma to be considered
	// specially. These calls have no return values, so this restriction avoids false positives.
	// It also simplifies removal logic.
	c.TestSame(
		"function a(){}\n" +
			"function b(){}\n" +
			"goog$inherits(b, a) + 1;\n")

	// Although a human is unlikely to write code like this, some optimizations may end up
	// converting an inherits call statement into the left side of a comma. We should still
	// remove this case.
	c.Test(
		"function a(){}\n"+
			"function b(){}\n"+
			"(goog$inherits(b, a), 1);\n",
		"1")
}

// Java: RemoveUnusedCodeTest.testRemoveInheritedClass13()
func (c *RemoveUnusedCodeTest) TestRemoveInheritedClass13() {
	c.Test(
		"class D {}\n"+
			"class C {}\n"+
			"valueType$mixin(C, D, 1, goog.reflect.objectProperty('a', C))\n",
		"")
}

// Java: RemoveUnusedCodeTest.testReflectedMethods()
func (c *RemoveUnusedCodeTest) TestReflectedMethods() {
	c.TestSame(
		"/** @constructor */\n" +
			"function Foo() {}\n" +
			"Foo.prototype.handle = function(x, y) { alert(y); };\n" +
			"var x = goog.reflect.object(Foo, {handle: 1});\n" +
			"for (var i in x) { x[i].call(x); }\n" +
			"window['Foo'] = Foo;\n")
}

// Java: RemoveUnusedCodeTest.testIssue618_1()
func (c *RemoveUnusedCodeTest) TestIssue618_1() {
	c.removeGlobal = false
	c.TestSame(
		"function f() {\n" +
			"\n" +
			"  var a = [], b;\n" +
			"\n" +
			"  a.push(b = []);\n" +
			"\n" +
			"  b[0] = 1;\n" +
			"\n" +
			"  return a;\n" +
			"\n" +
			"}\n")
}

// Java: RemoveUnusedCodeTest.testIssue618_2()
func (c *RemoveUnusedCodeTest) TestIssue618_2() {
	c.removeGlobal = false
	c.TestSame("var b; externVar.push(b = []); b[0] = 1;")
}

// Java: RemoveUnusedCodeTest.testBug38457201()
func (c *RemoveUnusedCodeTest) TestBug38457201() {
	c.removeGlobal = true
	c.Test(
		"var DOCUMENT_MODE = 1;\n"+
			"var temp;\n"+
			"false || (temp = (Number(DOCUMENT_MODE) >= 9));\n",
		"false || 0")
}

// Java: RemoveUnusedCodeTest.testBlockScoping()
func (c *RemoveUnusedCodeTest) TestBlockScoping() {
	c.Test("{let x;}", "{}")

	c.Test("{const x = 1;}", "{}")

	c.TestSame("{const x =1; alert(x)}")

	c.TestSame("{let x; alert(x)}")

	c.Test("let x = 1;", "")

	c.Test("let x; x = 3;", "")

	c.Test(
		"let x;\n"+
			"{ let x; } // unused in block\n"+
			"alert(x);\n",
		"let x; // outer x kept\n"+
			"{} // inner x removed\n"+
			"alert(x);\n")

	c.TestSame(
		"var x;\n" +
			"{\n" +
			"  var y = 1;\n" +
			"  alert(y); // keeps y alive\n" +
			"}\n" +
			"alert(x);\n")

	c.TestSame(
		"let x;\n" +
			"{\n" +
			"  let x = 1;\n" +
			"  alert(x); // keeps inner x alive\n" +
			"}\n" +
			"alert(x);\n")

	c.Test(
		"let x;\n"+
			"{ let x; } // inner x unused\n"+
			"let y; // outer y unused\n"+
			"alert(x);\n",
		"let x; // only outer x was used\n"+
			"{}\n"+
			"alert(x);\n")

	c.Test(
		"let x;\n"+
			"{\n"+
			"  let g;\n"+
			"  {\n"+
			"    const y = 1;\n"+
			"    alert(y); // keeps y alive\n"+
			"  }\n"+
			"}\n"+
			"let z;\n"+
			"alert(x);\n",
		"let x; // z removed\n"+
			"{\n"+
			"  { // g removed\n"+
			"    const y = 1;\n"+
			"    alert(y);\n"+
			"  }\n"+
			"}\n"+
			"alert(x);\n")

	c.Test(
		"let x;\n"+
			"{\n"+
			"  let x = 1;\n"+
			"  {\n"+
			"    let y;\n"+
			"  }\n"+
			"  alert(x); // keeps inner x alive\n"+
			"}\n"+
			"alert(x);\n",
		"let x; // outer x kept\n"+
			"{\n"+
			"  let x = 1;  // inner x kept\n"+
			"  {} // y removed\n"+
			"  alert(x);\n"+
			"}\n"+
			"alert(x);\n")

	c.Test(
		"let x;\n"+
			"{\n"+
			"  let x;\n"+
			"  alert(x); // keeps inner x alive\n"+
			"}\n",
		"{\n"+
			"  let x$jscomp$1; // inner x was renamed\n"+
			"  alert(x$jscomp$1);\n"+
			"}\n")
}

// Java: RemoveUnusedCodeTest.testArrowFunctions()
func (c *RemoveUnusedCodeTest) TestArrowFunctions() {
	c.Test(
		"class C {\n"+
			"  g() {\n"+
			"    var x;\n"+
			"  }\n"+
			"}\n"+
			"new C\n",
		"class C {\n"+
			"  g() {\n"+
			"  }\n"+
			"}\n"+
			"new C\n")

	c.Test("() => {var x}", "() => {};")

	c.Test("(x) => {}", "() => {};")

	c.TestSame("(x) => {alert(x)}")

	c.TestSame("(x) => {x + 1}")
}

// Java: RemoveUnusedCodeTest.testClasses()
func (c *RemoveUnusedCodeTest) TestClasses() {
	c.Test("var C = class {};", "")

	c.Test("class C {}", "")

	c.Test("class C {constructor(){} g() {}}", "")

	c.TestSame("class C {}; new C;")

	c.Test("var C = class X {}; new C;", "var C = class {}; new C;")

	c.TestSame(
		"class C{\n" +
			"  g() {\n" +
			"    var x;\n" +
			"    alert(x); // use in alert prevents removal of x\n" +
			"  }\n" +
			"}\n" +
			"new C\n")

	c.TestSame(
		"class C{\n" +
			"  g() {\n" +
			"    let x; // same as above but with 'let'\n" +
			"    alert(x);\n" +
			"  }\n" +
			"}\n" +
			"new C\n")

	c.Test(
		"class C {\n"+
			"  g() {\n"+
			"    let x; // unused x\n"+
			"  }\n"+
			"}\n"+
			"new C\n",
		"class C {\n"+
			"  g() { // x is removed\n"+
			"  }\n"+
			"}\n"+
			"new C\n")
}

// Java: RemoveUnusedCodeTest.testClassFields()
func (c *RemoveUnusedCodeTest) TestClassFields() {
	// Removal of individual fields is not enabled in this test class.
	// These tests confirm that the existence of fields does not prevent removal of a class that
	// is never instantiated.
	c.TestSame(
		"class C {\n" +
			"  x;\n" +
			"  y = 2;\n" +
			"  z = 'hi';\n" +
			"  static x;\n" +
			"  static y = 2;\n" +
			"  static z = 'hi';\n" +
			"}\n" +
			"new C\n")
	c.Test(
		"class C {\n"+
			"  x;\n"+
			"  y = 2;\n"+
			"  z = 'hi';\n"+
			"  static x;\n"+
			"  static y = 2;\n"+
			"  static z = 'hi';\n"+
			"}\n",
		"")
	c.TestSame(
		"class C {\n" +
			"  ['x'];\n" +
			"  1 = 2\n" +
			"  'a' = 'foo';\n" +
			"  static ['x'];\n" +
			"  static 1 = 2\n" +
			"  static 'a' = 'foo';\n" +
			"}\n" +
			"new C\n")
	c.Test(
		"class C {\n"+
			"  ['x'];\n"+
			"  1 = 2\n"+
			"  'a' = 'foo';\n"+
			"  static ['x'];\n"+
			"  static 1 = 2\n"+
			"  static 'a' = 'foo';\n"+
			"}\n",
		"")
}

// Java: RemoveUnusedCodeTest.testClassStaticBlocksDoesRemove()
func (c *RemoveUnusedCodeTest) TestClassStaticBlocksDoesRemove() {
	c.Test(
		"class C {\n"+
			"  static {\n"+
			"  }\n"+
			"}\n",
		"")

	// TODO(bradfordcsmith): Would be nice to remove the whole class at this point
	c.Test(
		"class C {\n"+
			"  static {\n"+
			"    var x = 1;\n"+
			"    let y = 2;\n"+
			"    const z = 3;\n"+
			"  }\n"+
			"}\n",
		"class C {\n"+
			"  static {\n"+
			"  }\n"+
			"}\n")

	c.Test(
		"class C {\n"+
			"  static {\n"+
			"    if (true) {\n"+
			"      if (true) {\n"+
			"        var foo = function() {};\n"+
			"      }\n"+
			"    }\n"+
			"  }\n"+
			"}\n",
		"class C {\n"+
			"  static {\n"+
			"    if (true) {\n"+
			"      if (true) {\n"+
			"      }\n"+
			"    }\n"+
			"  }\n"+
			"}\n")

	c.Test(
		"class C {\n"+
			"  static {\n"+
			"    function f() {}\n"+
			"  }\n"+
			"}\n",
		"class C {\n"+
			"  static {\n"+
			"  }\n"+
			"}\n")

	c.Test(
		"const x = 1;\n"+
			"class C {\n"+
			"  static {\n"+
			"    function f() {\n"+
			"      x;\n"+
			"    }\n"+
			"  }\n"+
			"}\n",
		"class C {\n"+
			"  static {\n"+
			"  }\n"+
			"}\n")
}

// Java: RemoveUnusedCodeTest.testClassStaticBlockDoesntRemove()
func (c *RemoveUnusedCodeTest) TestClassStaticBlockDoesntRemove() {
	c.TestSame(
		"class C {\n" +
			"  static {\n" +
			"    let x;\n" +
			"    alert(x);\n" +
			"    console.log(x);\n" +
			"  }\n" +
			"}\n")

	c.TestSame(
		"class C {\n" +
			"  static {\n" +
			"    this.x=1;\n" +
			"  }\n" +
			"}\n")

	c.TestSame(
		"const x = 1;\n" +
			"class C {\n" +
			"  static {\n" +
			"    x; // reference prevents `const x = 1;` from being removed.\n" +
			"  }\n" +
			"}\n" +
			"new C;\n")

	c.TestSame(
		"const x = 1;\n" +
			"class C {\n" +
			"  static {\n" +
			"// side-effect of `alert` prevents `x` from being removed\n" +
			"    alert(x); // reference prevents `x` from being removed.\n" +
			"  }\n" +
			"}\n")

	c.TestSame(
		"const x = 1;\n" +
			"class C {\n" +
			"  static {\n" +
			"    function f() {\n" +
			"      x; // reference prevents `const x = 1;` from being removed.\n" +
			"    }\n" +
			"    f(); // call prevents f() from being removed.\n" +
			"  }\n" +
			"}\n")
}

// Java: RemoveUnusedCodeTest.testComputedPropSideEffects()
func (c *RemoveUnusedCodeTest) TestComputedPropSideEffects() {
	c.TestSame(
		"class C {\n" +
			"  [alert(1)](){}\n" +
			"}\n")
}

// Java: RemoveUnusedCodeTest.testNonstaticClassFieldWithSideEffectsDoesRemove()
func (c *RemoveUnusedCodeTest) TestNonstaticClassFieldWithSideEffectsDoesRemove() {
	c.Test(
		"class C {\n"+
			"  y = alert(1);\n"+
			"}\n",
		"")

	c.Test(
		"class C {\n"+
			"  ['y'] = alert(1);\n"+
			"}\n",
		"")
}

// Java: RemoveUnusedCodeTest.testStaticClassFieldDetectsSideEffects()
func (c *RemoveUnusedCodeTest) TestStaticClassFieldDetectsSideEffects() {
	c.TestSame(
		"class C {\n" +
			"  static y = alert(1);\n" +
			"}\n")
}

// Java: RemoveUnusedCodeTest.testComputedClassFieldDetectsSideEffects()
func (c *RemoveUnusedCodeTest) TestComputedClassFieldDetectsSideEffects() {
	c.TestSame(
		"class C {\n" +
			"  [alert(1)];\n" +
			"}\n")

	c.TestSame(
		"class C {\n" +
			"  [alert(1)] = 'str';\n" +
			"}\n")
}

// Java: RemoveUnusedCodeTest.testStaticComputedClassFieldDetectsSideEffects()
func (c *RemoveUnusedCodeTest) TestStaticComputedClassFieldDetectsSideEffects() {
	c.TestSame(
		"class C {\n" +
			"  static [alert(1)];\n" +
			"}\n")

	c.TestSame(
		"class C {\n" +
			"  static [alert(1)] = 'str';\n" +
			"}\n")
}

// Java: RemoveUnusedCodeTest.testReferencesInClasses()
func (c *RemoveUnusedCodeTest) TestReferencesInClasses() {
	c.TestSame(
		"const A = 15;\n" +
			"const C = class {\n" +
			"  constructor() {\n" +
			"    this.a = A;\n" +
			"  }\n" +
			"}\n" +
			"new C;\n")
}

// Java: RemoveUnusedCodeTest.testRemoveGlobalClasses()
func (c *RemoveUnusedCodeTest) TestRemoveGlobalClasses() {
	c.removeGlobal = false
	c.TestSame("class C{}")
}

// Java: RemoveUnusedCodeTest.testSubClasses()
func (c *RemoveUnusedCodeTest) TestSubClasses() {
	c.Test("class D {} class C extends D {}", "")

	c.Test("class D {} class C extends D {} new D;", "class D {} new D;")

	c.TestSame("class D {} class C extends D {} new C;")
}

// Java: RemoveUnusedCodeTest.testGenerators()
func (c *RemoveUnusedCodeTest) TestGenerators() {
	c.Test(
		"function* f() {\n"+
			"  var x;\n"+
			"  yield x;\n"+
			"}\n",
		"")
	c.Test(
		"function* f() {\n"+
			"  var x;\n"+
			"  yield x;\n"+
			"}\n",
		"")
	c.Test(
		"function* f() {\n"+
			"  var x;\n"+
			"  var y;\n"+
			"  yield x;\n"+
			"}\n"+
			"f();\n",
		"function* f() {\n"+
			"  var x;\n"+
			"  yield x;\n"+
			"}\n"+
			"f()\n")
}

// Java: RemoveUnusedCodeTest.testForLet()
func (c *RemoveUnusedCodeTest) TestForLet() {
	// Could be optimized so that unused lets in the for loop are removed
	c.Test("for(let x; ;){}", "for(;;){}")

	c.Test("for(let x, y; ;) {x}", "for(let x; ;) {x}")
	c.Test("for(let x, y; ;) {y}", "for(let y; ;) {y}")

	c.Test("for(let x=0,y=0;;y++){}", "for(let y=0;;y++){}")
}

// Java: RemoveUnusedCodeTest.testForInLet()
func (c *RemoveUnusedCodeTest) TestForInLet() {
	c.TestSame("let item; for(item in externVar){}")
	c.TestSame("for(let item in externVar){}")
}

// Java: RemoveUnusedCodeTest.testForOf()
func (c *RemoveUnusedCodeTest) TestForOf() {
	c.TestSame("let item; for(item of externVar){}")

	c.TestSame("for(var item of externVar){}")

	c.TestSame("for(let item of externVar){}")

	c.TestSame(
		"var x;\n" +
			"for (var n of externVar) {\n" +
			"  if (x > n) {}; // keeps x alive\n" +
			"}\n")

	c.Test(
		"var x; // no references to this\n"+
			"for (var n of externVar) {\n"+
			"  if (n) {};\n"+
			"}\n",
		"for (var n of externVar) { // x removed\n"+
			"  if (n) {};\n"+
			"}\n")
}

// Java: RemoveUnusedCodeTest.testForAwaitOf()
func (c *RemoveUnusedCodeTest) TestForAwaitOf() {
	c.TestSame("async () => { let item; for await (item of externVar){} }")

	c.TestSame("async () => { for(var item of externVar){} }")

	c.TestSame("async () => { for(let item of externVar){} }")

	c.TestSame(
		"async () => {\n" +
			"  var x;\n" +
			"  for (var n of externVar) {\n" +
			"    if (x > n) {}; // keeps x alive\n" +
			"  }\n" +
			"};\n")

	c.Test(
		"async () => {\n"+
			"  var x; // no references to this\n"+
			"  for (var n of externVar) {\n"+
			"    if (n) {};\n"+
			"  }\n"+
			"};\n",
		"async () => {\n"+
			"  for (var n of externVar) { // x removed\n"+
			"    if (n) {};\n"+
			"  }\n"+
			"};\n")
}

// Java: RemoveUnusedCodeTest.testEnhancedForWithExternVar()
func (c *RemoveUnusedCodeTest) TestEnhancedForWithExternVar() {
	c.TestSame("for(externVar in externVar){}")
	c.TestSame("for(externVar of externVar){}")
}

// Java: RemoveUnusedCodeTest.testExternVarNotRemovable()
func (c *RemoveUnusedCodeTest) TestExternVarNotRemovable() {
	c.TestSame("externVar = 5;")
}

// Java: RemoveUnusedCodeTest.testTemplateStrings()
func (c *RemoveUnusedCodeTest) TestTemplateStrings() {
	c.TestSame(
		"var name = 'foo';\n" +
			"`Hello ${name}`\n")
}

// Java: RemoveUnusedCodeTest.testDestructuringArrayPattern0()
func (c *RemoveUnusedCodeTest) TestDestructuringArrayPattern0() {
	c.Test("function f(a) {} f();", "function f() {} f();")
	c.Test("function f([a]) {} f();", "function f() {} f();")
	c.Test("function f(...a) {} f();", "function f() {} f();")
	c.Test("function f(...[a]) {} f();", "function f() {} f();")
	c.TestSame("function f(...[...a]) {} f();")
	c.Test("function f(...{length:a}) {} f();", "function f() {} f();")

	c.Test("function f(a) {} f();", "function f() {} f();")
	c.Test("function f([a] = 1) {} f();", "function f() {} f();")
	c.Test("function f([a] = alert()) {} f();", "function f([] = alert()) {} f();")

	c.Test("function f(a = 1) {} f();", "function f() {} f();")
	c.Test("function f([a = 1]) {} f();", "function f() {} f();")
	c.Test("function f(...[a = 1]) {} f();", "function f() {} f();")
	c.Test("function f(...{length:a = 1}) {} f();", "function f() {} f();")

	c.TestSame("function f(a = alert()) {} f();")
	c.TestSame("function f([a = alert()]) {} f();")
	c.TestSame("function f(...[a = alert()]) {} f();")
	c.TestSame("function f(...{length:a = alert()}) {} f();")

	c.Test("function f([a] = []) {} f();", "function f() {} f();")
	c.Test("function f([a]) {} f();", "function f() {} f();")
	c.Test("function f([a], b) {} f();", "function f() {} f();")
	c.Test("function f([a], ...b) {} f();", "function f() {} f();")
	c.TestSame("function f([...a]) {} f();")
	c.Test("function f([[]]) {} f();", "function f([[]]) {} f();")
	c.TestSame("function f([[],...a]) {} f();")

	c.Test("var [a] = [];", "var [] = [];")
	c.TestSame("var [...a] = [];")
	c.TestSame("var [...[...a]] = [];")

	c.Test("var [a, b] = [];", "var [] = [];")
	c.Test("var [a, ...b] = [];", "var [ , ...b] = [];")
	c.Test("var [a, ...[...b]] = [];", "var [ ,...[...b]] = [];")

	c.Test("var [a, b, c] = []; use(a, b);", "var [a,b  ] = []; use(a, b);")
	c.Test("var [a, b, c] = []; use(a, c);", "var [a, ,c] = []; use(a, c);")
	c.Test("var [a, b, c] = []; use(b, c);", "var [ ,b,c] = []; use(b, c);")

	c.Test("var [a, b, c] = []; use(a);", "var [a] = []; use(a);")
	c.Test("var [a, b, c] = []; use(b);", "var [ ,b, ] = []; use(b);")
	c.Test("var [a, b, c] = []; use(c);", "var [ , ,c] = []; use(c);")

	c.Test("var a, b, c; [a, b, c] = []; use(a);", "var a; [a] = []; use(a);")
	c.Test("var a, b, c; [a, b, c] = []; use(b);", "var b; [ ,b, ] = []; use(b);")
	c.Test("var a, b, c; [a, b, c] = []; use(c);", "var c; [ , ,c] = []; use(c);")

	c.Test("var a, b, c; [[a, b, c]] = []; use(a);", "var a; [[a]] = []; use(a);")
	c.Test("var a, b, c; [{a, b, c}] = []; use(a);", "var a; [{a}] = []; use(a);")
	c.Test("var a, b, c; ({x:[a, b, c]} = []); use(a);", "var a; ({x:[a]} = []); use(a);")

	c.Test("var a, b, c; [[a, b, c]] = []; use(b);", "var b; [[ ,b]] = []; use(b);")
	c.Test("var a, b, c; [[a, b, c]] = []; use(c);", "var c; [[ , ,c]] = []; use(c);")

	c.Test("var a, b, c; [a, b, ...c] = []; use(a);", "var a, c; [a,  , ...c] = []; use(a);")
	c.Test("var a, b, c; [a, b, ...c] = []; use(b);", "var b, c; [ , b, ...c] = []; use(b);")
	c.Test("var a, b, c; [a, b, ...c] = []; use(c);", "var c;    [ ,  , ...c] = []; use(c);")

	c.Test("var a, b, c; [a=1, b=2, c=3] = []; use(a);", "var a; [a=1] = []; use(a);")
	c.Test("var a, b, c; [a=1, b=2, c=3] = []; use(b);", "var b; [   ,b=2, ] = []; use(b);")
	c.Test("var a, b, c; [a=1, b=2, c=3] = []; use(c);", "var c; [   ,   ,c=3] = []; use(c);")

	c.TestSame("var a, b, c; [a.x, b.y, c.z] = []; use(a);") // unnecessary retention of b,c
	c.TestSame("var a, b, c; [a.x, b.y, c.z] = []; use(b);") // unnecessary retention of a,c
	c.TestSame("var a, b, c; [a.x, b.y, c.z] = []; use(c);") // unnecessary retention of a,b

	c.TestSame("var a, b, c; [a().x, b().y, c().z] = []; use(a);")
	c.TestSame("var a, b, c; [a().x, b().y, c().z] = []; use(b);")
	c.TestSame("var a, b, c; [a().x, b().y, c().z] = []; use(c);")
}

// Java: RemoveUnusedCodeTest.testDestructuringArrayPattern1()
func (c *RemoveUnusedCodeTest) TestDestructuringArrayPattern1() {
	c.Test(
		"var a; var b\n"+
			"[a, b] = [1, 2]\n",
		"[] = [1, 2]")

	c.Test(
		"var b; var a\n"+
			"[a, b] = [1, 2]\n",
		"[] = [1, 2]")

	c.Test(
		"var a; var b;\n"+
			"[a] = [1]\n",
		"[] = [1]")

	c.TestSame("var [a, b] = [1, 2]; alert(a); alert(b);")
}

// TODO(b/538150824): RemoveUnusedCode: preserve trailing empty slots in destructuring array
// patterns
// Java: RemoveUnusedCodeTest.testDestructuringArrayPatternImpureRhs()
func (c *RemoveUnusedCodeTest) TestDestructuringArrayPatternImpureRhs() {
	c.Test("const [ok] = validate(userInput);", "const [] = validate(userInput);")
	c.Test("var [a, b, c] = gen(); use(a);", "var [a] = gen(); use(a);")
	c.Test("var a, b, c; [a, b, c] = gen(); use(a);", "var a; [a] = gen(); use(a);")
	c.Test("var [a, b, c] = [1, 2, 3]; use(a);", "var [a] = [1, 2, 3]; use(a);")
}

// Java: RemoveUnusedCodeTest.testDestructuringObjectPattern0()
func (c *RemoveUnusedCodeTest) TestDestructuringObjectPattern0() {
	c.Test("function f({a}) {} f();", "function f() {} f();")
	c.Test("function f({a:b}) {} f();", "function f() {} f();")
	c.Test("function f({[a]:b}) {} f();", "function f() {} f();")
	c.Test("function f({['a']:b}) {} f();", "function f() {} f();")
	c.TestSame("function f({[alert()]:b}) {} f();")
	c.Test("function f({a} = {}) {} f();", "function f() {} f();")
	c.Test("function f({a:b} = {}) {} f();", "function f() {} f();")
	c.TestSame("function f({[alert()]:b} = {}) {} f();")
	c.Test("function f({a} = alert()) {} f();", "function f({} = alert()) {} f();")
	c.Test("function f({a:b} = alert()) {} f();", "function f({} = alert()) {} f();")
	c.TestSame("function f({[alert()]:b} = alert()) {} f();")
	c.Test("function f({a = 1}) {} f();", "function f() {} f();")
	c.Test("function f({a:b = 1}) {} f();", "function f() {} f();")
	c.Test("function f({[a]:b = 1}) {} f();", "function f() {} f();")
	c.Test("function f({['a']:b = 1}) {} f();", "function f() {} f();")
	c.TestSame("function f({[alert()]:b = 1}) {} f();") // fix me, remove "= 1"
	c.Test("function f({a = 1}) {} f();", "function f() {} f();")
	c.TestSame("function f({a:b = alert()}) {} f();")
	c.TestSame("function f({[externVar]:b = alert()}) {} f();")
	c.TestSame("function f({['a']:b = alert()}) {} f();")
	c.TestSame("function f({[alert()]:b = alert()}) {} f();")

	c.Test("function f({a}, c) {use(c)} f();", "function f({}, c) {use(c)} f();")
	c.Test("function f({a:b}, c) {use(c)} f();", "function f({}, c) {use(c)} f();")
	c.Test("function f({[externVar]:b}, c) {use(c)} f();", "function f({}, c) {use(c)} f();")
	c.Test("function f({['a']:b}, c) {use(c)} f();", "function f({}, c) {use(c)} f();")
	c.TestSame("function f({[use()]:b}, c) {use(c)} f();")
	c.Test("function f({a} = {}, c) {use(c)} f();", "function f({} = {}, c) {use(c)} f();")
	c.Test("function f({a:b} = {}, c) {use(c)} f();", "function f({} = {}, c) {use(c)} f();")
	c.TestSame("function f({[use()]:b} = {}, c) {use(c)} f();")
	c.Test(
		"function f({a} = use(), c) {use(c)} f();", // preserve alignment
		"function f({ } = use(), c) {use(c)} f();")
	c.Test(
		"function f({a:b} = use(), c) {use(c)} f();", // preserve alignment
		"function f({   } = use(), c) {use(c)} f();")
	c.TestSame("function f({[use()]:b} = use(), c) {use(c)} f();")
	c.Test("function f({a = 1}, c) {use(c)} f();", "function f({}, c) {use(c)} f();")
	c.Test("function f({a:b = 1}, c) {use(c)} f();", "function f({}, c) {use(c)} f();")
	c.Test("function f({[a]:b = 1}, c) {use(c)} f();", "function f({}, c) {use(c)} f();")
	c.Test("function f({['a']:b = 1}, c) {use(c)} f();", "function f({}, c) {use(c)} f();")
	c.TestSame("function f({[use()]:b = 1}, c) {use(c)} f();") // fix me, remove "= 1"

	c.Test("var {a} = {a:1};", "var {} = {a:1};")
	c.Test("var {a:a} = {a:1};", "var {} = {a:1};")
	c.Test("var {['a']:a} = {a:1};", "var {} = {a:1};")
	c.Test("var {['a']:a = 1} = {a:1};", "var {} = {a:1};")
	c.TestSame("var {[alert()]:a = 1} = {a:1};")
	c.Test("var {a:a = 1} = {a:1};", "var {} = {a:1};")
	c.TestSame("var {a:a = alert()} = {a:1};")

	c.Test("var a; ({a} = {a:1});", "({} = {a:1});")
	c.Test("var a; ({a:a} = {a:1});", "({} = {a:1});")
	c.Test("var a; ({['a']:a} = {a:1});", "({} = {a:1});")
	c.Test("var a; ({['a']:a = 1} = {a:1});", "({} = {a:1});")
	c.TestSame("var a; ({[alert()]:a = 1} = {a:1});")
	c.Test("var a; ({a:a = 1} = {a:1});", "({} = {a:1});")
	c.TestSame("var a; ({a:a = alert()} = {a:1});")

	c.TestSame("var a = {}; ({a:a.foo} = {a:1});")
	c.TestSame("var a = {}; ({['a']:a.foo} = {a:1});")
	c.TestSame("var a = {}; ({['a']:a.foo = 1} = {a:1});")
	c.TestSame("var a = {}; ({[alert()]:a.foo = 1} = {a:1});")
	c.TestSame("var a = {}; ({a:a.foo = 1} = {a:1});")
	c.TestSame("var a = {}; ({a:a.foo = alert()} = {a:1});")

	c.TestSame("function a() {} ({a:a().foo} = {a:1});")
	c.TestSame("function a() {} ({['a']:a().foo} = {a:1});")
	c.TestSame("function a() {} ({['a']:a().foo = 1} = {a:1});")
	c.TestSame("function a() {} ({[alert()]:a().foo = 1} = {a:1});")
	c.TestSame("function a() {} ({a:a().foo = 1} = {a:1});")
	c.TestSame("function a() {} ({a:a().foo = alert()} = {a:1});")
}

// Java: RemoveUnusedCodeTest.testDestructuringObjectPattern1()
func (c *RemoveUnusedCodeTest) TestDestructuringObjectPattern1() {
	c.Test("var a; var b; ({a, b} = {a:1, b:2})", "({} = {a:1, b:2});")
	c.Test("var a; var b; ({a:a, b:b} = {a:1, b:2})", "({} = {a:1, b:2});")
	c.TestSame("var a; var b; ({[alert()]:a, [alert()]:b} = {x:1, y:2})")

	c.Test("var a; var b; ({a} = {a:1})", "({} = {a:1})")
	c.Test("var a; var b; ({a:a.foo} = {a:1})", "var a; ({a:a.foo} = {a:1})")

	c.TestSame("var {a, b} = {a:1, b:2}; alert(a); alert(b);")

	c.TestSame("var {a} = {p:{}, q:{}}; a.q = 4;")

	// Nested Destructuring
	c.Test(
		"const someObject = {a:{ b:1 }}; var {a: {b}} = someObject; someObject.a.b;",
		"const someObject = {a:{ b:1 }}; var {a: {}} = someObject; someObject.a.b;")

	c.TestSame("const someObject = {a:{ b:1 }}; var {a: {b}} = someObject; b;")
}

// Java: RemoveUnusedCodeTest.testRemoveUnusedVarsDeclaredInDestructuring0()
func (c *RemoveUnusedCodeTest) TestRemoveUnusedVarsDeclaredInDestructuring0() {
	// Array destructuring
	c.Test("var [a, b] = [1, 2]; alert(a);", "var [a] = [1, 2]; alert(a);")

	c.Test("var [a, b] = [1, 2];", "var [] = [1, 2];")

	// Object pattern destructuring
	c.Test(
		"var {a, b} = {a:1, b:2}; alert(a);", // comment to keep before/after alignment
		"var {a,  } = {a:1, b:2}; alert(a);")

	c.Test(
		"var {a, b} = {a:1, b:2};", // comment to keep before/after alignment
		"var {    } = {a:1, b:2};")
}

// Java: RemoveUnusedCodeTest.testRemoveUnusedVarsDeclaredInDestructuring1()
func (c *RemoveUnusedCodeTest) TestRemoveUnusedVarsDeclaredInDestructuring1() {
	// Nested pattern
	c.Test(
		"var {a, b:{c}} = {a:1, b:{c:5}}; alert(a);", "var {a, b:{ }} = {a:1, b:{c:5}}; alert(a);")

	c.TestSame("var {a, b:{c}} = {a:1, b:{c:5}}; alert(a, c);")

	c.Test("var {a, b:{c}} = externVar;", "var {b:{}} = externVar;")

	// Value may have side effects
	c.Test(
		"var {a, b} = {a:alert(), b:alert()};", // preserve before / after alignment
		"var {    } = {a:alert(), b:alert()};")

	c.Test(
		"var {a, b} = alert();", // preserve before / after alignment
		"var {    } = alert();")

	// Same as above case without the destructuring declaration
	c.Test("var a, b = alert();", "alert();")
}

// RemoveUnusedCodeTest_PolyfillRemovalTester handles the boilerplate for
// testing whether a polyfill is removed or not.
type RemoveUnusedCodeTest_PolyfillRemovalTester struct {
	outer     *RemoveUnusedCodeTest
	externs   []string
	polyfills *jutil.LinkedHashSet[string]

	// Input source code, which will be prefixed with the polyfills.
	// A value of `null` just means this hasn't been set yet.
	inputSource *string

	// Expected output source code, which will be prefixed with the output version of each polyfill.
	// A value of `null` means this value hasn't been set yet or needs to be reset because
	// inputSource has changed.
	expectedSource *string

	// Set of polyfills that are expected to be removed by RemoveUnusedCode.
	// A value of `null` indicates that no expectation has been set since the last time inputSource
	// was modified.
	polyfillsExpectedToBeRemoved map[string]bool
}

// Java: RemoveUnusedCodeTest.PolyfillRemovalTester()
func NewRemoveUnusedCodeTest_PolyfillRemovalTester(outer *RemoveUnusedCodeTest) *RemoveUnusedCodeTest_PolyfillRemovalTester {
	return &RemoveUnusedCodeTest_PolyfillRemovalTester{
		outer:                        outer,
		externs:                      []string{},
		polyfills:                    jutil.NewLinkedHashSet[string](),
		polyfillsExpectedToBeRemoved: map[string]bool{},
	}
}

// Java: RemoveUnusedCodeTest.PolyfillRemovalTester.addExterns(String)
func (p *RemoveUnusedCodeTest_PolyfillRemovalTester) AddExterns(moreExterns string) *RemoveUnusedCodeTest_PolyfillRemovalTester {
	p.externs = append(p.externs, moreExterns)
	return p
}

// Java: RemoveUnusedCodeTest.PolyfillRemovalTester.addPolyfill(String)
func (p *RemoveUnusedCodeTest_PolyfillRemovalTester) AddPolyfill(polyfill string) *RemoveUnusedCodeTest_PolyfillRemovalTester {
	jutil.CheckArgumentMsg1(!p.polyfills.Contains(polyfill), "duplicate polyfill added: >%s<", polyfill)
	p.polyfills.Add(polyfill)
	// force update of expectations whenever the polyfills list is updated
	p.polyfillsExpectedToBeRemoved = nil
	return p
}

// Java: RemoveUnusedCodeTest.PolyfillRemovalTester.inputSource(String)
func (p *RemoveUnusedCodeTest_PolyfillRemovalTester) InputSource(inputSource string) *RemoveUnusedCodeTest_PolyfillRemovalTester {
	p.inputSource = &inputSource
	// Force updates for the expected output source and polyfills
	p.expectedSource = nil
	p.polyfillsExpectedToBeRemoved = nil
	return p
}

// Java: RemoveUnusedCodeTest.PolyfillRemovalTester.expectSourceUnchanged()
func (p *RemoveUnusedCodeTest_PolyfillRemovalTester) ExpectSourceUnchanged() *RemoveUnusedCodeTest_PolyfillRemovalTester {
	jutil.CheckNotNull(p.inputSource)
	p.expectedSource = p.inputSource
	return p
}

// Java: RemoveUnusedCodeTest.PolyfillRemovalTester.expectSource(String)
func (p *RemoveUnusedCodeTest_PolyfillRemovalTester) ExpectSource(expectedSource string) *RemoveUnusedCodeTest_PolyfillRemovalTester {
	p.expectedSource = &expectedSource
	return p
}

// Java: RemoveUnusedCodeTest.PolyfillRemovalTester.expectNoPolyfillsRemoved()
func (p *RemoveUnusedCodeTest_PolyfillRemovalTester) ExpectNoPolyfillsRemoved() *RemoveUnusedCodeTest_PolyfillRemovalTester {
	p.polyfillsExpectedToBeRemoved = map[string]bool{}
	return p
}

// Java: RemoveUnusedCodeTest.PolyfillRemovalTester.expectPolyfillsRemoved(String...)
func (p *RemoveUnusedCodeTest_PolyfillRemovalTester) ExpectPolyfillsRemoved(removedPolyfills ...string) *RemoveUnusedCodeTest_PolyfillRemovalTester {
	for _, polyfillToRemove := range removedPolyfills {
		p.ExpectPolyfillRemoved(polyfillToRemove)
	}
	return p
}

// Java: RemoveUnusedCodeTest.PolyfillRemovalTester.expectPolyfillRemoved(String)
func (p *RemoveUnusedCodeTest_PolyfillRemovalTester) ExpectPolyfillRemoved(polyfill string) *RemoveUnusedCodeTest_PolyfillRemovalTester {
	jutil.CheckArgumentMsg1(
		p.polyfills.Contains(polyfill), "non-existent polyfill cannot be removed: >%s<", polyfill)
	if p.polyfillsExpectedToBeRemoved == nil {
		p.polyfillsExpectedToBeRemoved = map[string]bool{}
	} else {
		jutil.CheckArgumentMsg1(
			!p.polyfillsExpectedToBeRemoved[polyfill],
			"polyfill cannot be removed twice: >%s<",
			polyfill)
	}
	p.polyfillsExpectedToBeRemoved[polyfill] = true
	return p
}

// Java: RemoveUnusedCodeTest.PolyfillRemovalTester.test()
func (p *RemoveUnusedCodeTest_PolyfillRemovalTester) Test() {
	externsText := strings.Join(p.externs, "\n")

	// Java: Stream.concat(polyfills.stream(), Stream.of(inputSource)).collect(joining("\n"));
	// a null inputSource joins as "null".
	inputSourceText := strings.Join(append(p.polyfills.ToSlice(), removeUnusedCodeTest_nullableString(p.inputSource)), "\n")

	jutil.CheckNotNullMsg(p.polyfillsExpectedToBeRemoved, "unset/undefined polyfill removal behavior")
	expectedPolyfills := []string{}
	for polyfill := range p.polyfills.All() {
		if !p.polyfillsExpectedToBeRemoved[polyfill] {
			expectedPolyfills = append(expectedPolyfills, polyfill)
		}
	}
	expectedSourceText := strings.Join(append(expectedPolyfills, removeUnusedCodeTest_nullableString(p.expectedSource)), "\n")

	p.outer.Test(
		p.outer.Externs(externsText), p.outer.Srcs(inputSourceText), p.outer.Expected(expectedSourceText))
}

// removeUnusedCodeTest_nullableString is String.valueOf for a @Nullable String
// (Go-only: language-adapter).
func removeUnusedCodeTest_nullableString(s *string) string {
	if s == nil {
		return "null"
	}
	return *s
}

// Java: RemoveUnusedCodeTest.PolyfillRemovalTester.expectNoRemovalTest(String)
func (p *RemoveUnusedCodeTest_PolyfillRemovalTester) ExpectNoRemovalTest(src string) {
	/* no polyfills removed */
	p.InputSource(src). //
				ExpectSourceUnchanged().
				ExpectNoPolyfillsRemoved().
				Test()
}

// Java: RemoveUnusedCodeTest.PolyfillRemovalTester.expectPolyfillsRemovedTest(String, String...)
func (p *RemoveUnusedCodeTest_PolyfillRemovalTester) ExpectPolyfillsRemovedTest(src string, removedPolyfills ...string) {
	p.InputSource(src). //
				ExpectSourceUnchanged().
				ExpectPolyfillsRemoved(removedPolyfills...).
				Test()
}

// Java: RemoveUnusedCodeTest.testRemoveUnusedPolyfills_global_untyped()
func (c *RemoveUnusedCodeTest) TestRemoveUnusedPolyfills_global_untyped() {
	mapPolyfill := "$jscomp.polyfill('Map', function() {}, 'es6', 'es3');"
	tester :=
		NewRemoveUnusedCodeTest_PolyfillRemovalTester(c).
			AddExterns(
				NewTestExternsBuilder().AddConsole().AddExtra(RemoveUnusedCodeTest_JSCOMP_POLYFILL).AddMap().Build()).
			AddPolyfill(mapPolyfill)

	// unused polyfill is removed
	tester.ExpectPolyfillsRemovedTest("console.log()", mapPolyfill)
	// used polyfill is not removed
	tester.ExpectNoRemovalTest("console.log(new Map())")

	// Local names shadowing global polyfills are not themselves polyfill references.
	tester.ExpectPolyfillsRemovedTest(
		"console.log(function(Map$jscomp$1) {\n"+
			"  console.log(new Map$jscomp$1());\n"+
			"});\n",
		mapPolyfill)
}

// Java: RemoveUnusedCodeTest.testRemoveUnusedPolyfills_global_typed()
func (c *RemoveUnusedCodeTest) TestRemoveUnusedPolyfills_global_typed() {
	c.EnableTypeCheck()
	// Type information is no longer used to make decisions for polyfill removal.
	c.TestRemoveUnusedPolyfills_global_untyped()
}

// Java: RemoveUnusedCodeTest.testRemoveUnusedPolyfills_propertyOfGlobalObject_untyped()
func (c *RemoveUnusedCodeTest) TestRemoveUnusedPolyfills_propertyOfGlobalObject_untyped() {
	mapPolyfill := "$jscomp.polyfill('Map', function() {}, 'es6', 'es3');"
	tester :=
		NewRemoveUnusedCodeTest_PolyfillRemovalTester(c).
			AddExterns(
				NewTestExternsBuilder().
					AddConsole().
					AddExtra(
						RemoveUnusedCodeTest_JSCOMP_POLYFILL,
						"/** @const */ var goog = {};",
						"/** @const {!Global} */ goog.global;",
						"/** @const */ goog.structs = {};",
						"/** @constructor */ goog.structs.Map = function() {};",
						"/** @type {!Global} */ var someGlobal;",
						"/** @const */ var notGlobal = {};",
						"/** @constructor */ notGlobal.Map = function() {};",
						"").
					AddMap().
					Build()).
			AddPolyfill(mapPolyfill)

	// Global polyfills may be accessed as properties on the global object.
	tester.ExpectNoRemovalTest("console.log(new goog.global.Map());")

	// NOTE: Without type information we don't know whether a property access is on the global
	// object or not.
	tester.ExpectNoRemovalTest("console.log(new goog.structs.Map());")
	tester.ExpectNoRemovalTest("console.log(new notGlobal.Map());")
	tester.ExpectNoRemovalTest(
		"var x = {Map: /** @constructor */ function() {}};\n" +
			"console.log(new x.Map());\n")
}

// Java: RemoveUnusedCodeTest.testRemoveUnusedPolyfills_propertyOfGlobalObject_withOptionalChain()
func (c *RemoveUnusedCodeTest) TestRemoveUnusedPolyfills_propertyOfGlobalObject_withOptionalChain() {
	promisePolyfill := "$jscomp.polyfill('Promise', function() {}, 'es6', 'es3');"
	tester :=
		NewRemoveUnusedCodeTest_PolyfillRemovalTester(c).
			AddExterns(
				NewTestExternsBuilder().
					AddConsole().
					AddPromise().
					AddExtra(RemoveUnusedCodeTest_JSCOMP_POLYFILL, "/** @const */ var goog = {};").
					Build()).
			AddPolyfill(promisePolyfill)

	// The optional chain here isn't guarding against Promise being undefined,
	// just the object on which it is referenced.
	tester.ExpectNoRemovalTest("console.log(goog.global?.Promise.resolve());")

	// The optional chain here _is_ guarding against Promise being undefined.
	tester.ExpectPolyfillsRemovedTest(
		"console.log(goog.global.Promise?.resolve());", promisePolyfill)

	// The optional chain here _is_ guarding against Promise being undefined,
	// in addition to the object on which it is referenced.
	tester.ExpectPolyfillsRemovedTest(
		"console.log(goog.global?.Promise?.resolve());", promisePolyfill)
}

// Java: RemoveUnusedCodeTest.testRemoveUnusedPolyfills_propertyOfGlobalObject_typed()
func (c *RemoveUnusedCodeTest) TestRemoveUnusedPolyfills_propertyOfGlobalObject_typed() {
	c.EnableTypeCheck()
	// We no longer use type information to make polyfill removal decisions
	c.TestRemoveUnusedPolyfills_propertyOfGlobalObject_untyped()
}

// Java: RemoveUnusedCodeTest.testRemoveUnusedPolyfills_staticProperty_untyped()
func (c *RemoveUnusedCodeTest) TestRemoveUnusedPolyfills_staticProperty_untyped() {
	arrayDotFromPolyfill :=
		"$jscomp.polyfill('Array.from', function() {}, 'es6', 'es3');"
	tester :=
		NewRemoveUnusedCodeTest_PolyfillRemovalTester(c).
			AddExterns(
				NewTestExternsBuilder().
					AddConsole().
					AddArray().
					AddExtra(
						RemoveUnusedCodeTest_JSCOMP_POLYFILL,
						"/** @constructor */ function Set() {}",
						// NOTE: this is not a real API but it allows testing that we can tell it
						// apart.
						"Set.from = function() {};",
						"/** @const */ var goog = {};",
						"/** @const {!Global} */ goog.global;",
						"").
					Build()).
			AddPolyfill(arrayDotFromPolyfill)

	// Used polyfill is retained.
	tester.ExpectNoRemovalTest("console.log(Array.from([]));")

	// Used polyfill is retained, even if accessed as a property of global.
	tester.ExpectNoRemovalTest("console.log(goog.global.Array.from([]));")

	// Unused polyfill is not removed if there is another static property with the same name
	tester.ExpectNoRemovalTest(
		"class NotArray { static from() {} }\n" +
			"console.log(NotArray.from());\n")

	// Without type information, we can't correctly remove this polyfill.
	tester.ExpectNoRemovalTest(
		"var x = {Array: {from: function() {}}};\n" +
			"console.log(x.Array.from());\n")

	// Used polyfill via aliased owner: retains definition.
	tester.ExpectNoRemovalTest(
		"/** @const */ var MyArray = Array;\n" +
			"// polyfill is kept even though called via an alias\n" +
			"console.log(MyArray.from([]));\n")

	// Used polyfill via subclass: retains definition.
	tester.ExpectNoRemovalTest(
		"class SubArray extends Array {}\n" +
			"// polyfill is kept even though called via a subclass.\n" +
			"console.log(SubArray.from([]));\n")

	// Cannot distinguish between Set.from and Array.from,
	// so the Set.from polyfill will also be kept.
	setDotFromPolyfill := "$jscomp.polyfill('Set.from', function() {}, 'es6', 'es3');"
	tester.AddPolyfill(setDotFromPolyfill)
	tester.ExpectNoRemovalTest("console.log(Array.from([]));")
}

// Java: RemoveUnusedCodeTest.testRemoveUnusedPolyfills_staticProperty_withOptionalChain()
func (c *RemoveUnusedCodeTest) TestRemoveUnusedPolyfills_staticProperty_withOptionalChain() {
	promisePolyfill := "$jscomp.polyfill('Promise', function() {}, 'es6', 'es3');"
	allSettledPolyfill :=
		"$jscomp.polyfill('Promise.allSettled', function() {}, 'es6', 'es3');"
	tester :=
		NewRemoveUnusedCodeTest_PolyfillRemovalTester(c).
			AddExterns(
				NewTestExternsBuilder().
					AddConsole().
					AddPromise().
					AddExtra(RemoveUnusedCodeTest_JSCOMP_POLYFILL, "/** @const */ var goog = {};").
					Build()).
			AddPolyfill(promisePolyfill).
			AddPolyfill(allSettledPolyfill)

	tester.ExpectNoRemovalTest("console.log(Promise.allSettled());") // neither guarded

	// `?.` guards the name on its left, allowing it to be removed
	tester.ExpectPolyfillsRemovedTest(
		"console.log(Promise.allSettled?.());", // allSettled guarded
		allSettledPolyfill)
	tester.ExpectPolyfillsRemovedTest(
		"console.log(Promise.allSettled?.(Promise.allSettled));", // allSettled guarded
		allSettledPolyfill)

	// TODO(b/163394833): Note that this behavior could result in the Promise polyfill being
	// removed, then the Promise.allSettled polyfill code having nowhere to hang its value.
	// At runtime the `$jscomp.polyfill('Promise.allSettled',...)` call will silently fail
	// to create the polyfill if `Promise` isn't defined.
	// This is consistent with the way guarding with `&&` would behave, as demonstrated by the
	// following test case.
	tester.ExpectPolyfillsRemovedTest(
		"console.log(Promise?.allSettled());", // Promise guarded
		promisePolyfill)

	tester.ExpectPolyfillsRemovedTest(
		"console.log(Promise && Promise.allSettled());", // Promise guarded
		promisePolyfill)

	tester.ExpectPolyfillsRemovedTest(
		"console.log(Promise?.allSettled?.([Promise.allSettled]));", // both guarded
		promisePolyfill,
		allSettledPolyfill)

	// Access via a (potentially) global variable should behave the same way, even if the
	// global is referenced via optional chaining
	tester.ExpectNoRemovalTest(
		"console.log(goog.global?.Promise.allSettled());") // neither guarded

	tester.ExpectPolyfillsRemovedTest(
		"console.log(goog.global?.Promise.allSettled?.());", // allSettled guarded
		allSettledPolyfill)

	// TODO(b/163394833): Note that this behavior could result in the Promise polyfill being
	// removed, then the Promise.allSettled polyfill code having nowhere to hang its value.
	// At runtime the `$jscomp.polyfill('Promise.allSettled',...)` call will silently fail
	// to create the polyfill if `Promise` isn't defined.
	// This is consistent with the way guarding with `&&` would behave, as demonstrated by the
	// following test case.
	tester.ExpectPolyfillsRemovedTest(
		"console.log(goog.global?.Promise?.allSettled());", // Promise guarded
		promisePolyfill)

	tester.ExpectPolyfillsRemovedTest(
		"console.log(goog.global && goog.global.Promise && goog.global.Promise.allSettled());",
		promisePolyfill)

	tester.ExpectPolyfillsRemovedTest(
		"console.log(goog.global?.Promise?.allSettled?.());", // both guarded
		promisePolyfill,
		allSettledPolyfill)
}

// Java: RemoveUnusedCodeTest.testRemoveUnusedPolyfills_staticProperty_typed()
func (c *RemoveUnusedCodeTest) TestRemoveUnusedPolyfills_staticProperty_typed() {
	c.EnableTypeCheck()
	// NOTE: We no longer use type information to make polyfill removal decisions.
	c.TestRemoveUnusedPolyfills_staticProperty_untyped()
}

// Java: RemoveUnusedCodeTest.testRemoveUnusedPolyfills_prototypeProperty_untyped()
func (c *RemoveUnusedCodeTest) TestRemoveUnusedPolyfills_prototypeProperty_untyped() {
	stringRepeatPolyfill :=
		"$jscomp.polyfill('String.prototype.repeat', function() {}, 'es6', 'es3');"
	tester :=
		NewRemoveUnusedCodeTest_PolyfillRemovalTester(c).
			AddExterns(
				NewTestExternsBuilder().
					AddArray().
					AddConsole().
					AddString().
					AddFunction().
					AddAlert().
					AddExtra(
						RemoveUnusedCodeTest_JSCOMP_POLYFILL, //
						"function externFunction() {}",
						"").
					Build()).
			AddPolyfill(stringRepeatPolyfill)

	// Unused polyfill is removed.
	tester.ExpectPolyfillsRemovedTest("console.log();", stringRepeatPolyfill)

	// Used polyfill is retained.
	tester.ExpectNoRemovalTest("''.repeat(1);")
	tester.ExpectNoRemovalTest("''?.repeat(1);")

	// Guarded polyfill is removed.
	tester.ExpectPolyfillsRemovedTest("''.repeat?.(1);", stringRepeatPolyfill)

	// Used polyfill (directly via String.prototype) is retained.
	tester.ExpectNoRemovalTest("String.prototype.repeat(1);")
	tester.ExpectNoRemovalTest("String.prototype?.repeat(1);")
	tester.ExpectPolyfillsRemovedTest("String.prototype.repeat?.(1);", stringRepeatPolyfill)

	// Used polyfill (directly String.prototype and Function.prototype.call) is retained.
	tester.ExpectNoRemovalTest("String.prototype.repeat.call('', 1);")
	tester.ExpectNoRemovalTest("String.prototype.repeat.call?.('', 1);")
	tester.ExpectPolyfillsRemovedTest(
		"String.prototype.repeat?.call('', 1);", stringRepeatPolyfill)

	// Unused polyfill is not removed if there is another property with the same name on an unknown
	// type
	tester.ExpectNoRemovalTest(
		"var x = externFunction();\n" +
			"x.repeat();\n")
	tester.ExpectPolyfillsRemovedTest(
		"var x = externFunction();\n"+
			"x.repeat?.();\n",
		stringRepeatPolyfill)

	// Without type information, cannot remove the polyfill.
	tester.ExpectNoRemovalTest(
		"class Repeatable {\n" +
			"  static repeat() {}\n" +
			"};\n" +
			"Repeatable.repeat();\n")

	// Without type information, cannot remove the polyfill.
	tester.ExpectNoRemovalTest(
		"class Repeatable {\n" +
			"  repeat() {}\n" +
			"};\n" +
			"var x = new Repeatable();\n" +
			"x.repeat();\n")

	// Multiple same-name methods
	stringIncludesPolyfill :=
		"$jscomp.polyfill('String.prototype.includes', function() {}, 'es6', 'es3');"
	arrayIncludesPolyfill :=
		"$jscomp.polyfill('Array.prototype.includes', function() {}, 'es6', 'es3');"
	tester.
		AddPolyfill(stringIncludesPolyfill). //
		AddPolyfill(arrayIncludesPolyfill)

	// The unused `String.prototype.repeat` polyfill is removed, but both of the
	// `(Array|String).prototype.includes` polyfills are kept, since we aren't using type
	// information to recognize which one is being called.
	tester.ExpectPolyfillsRemovedTest("[].includes(5);", stringRepeatPolyfill)
}

// Java: RemoveUnusedCodeTest.testRemoveUnusedPolyfills_prototypeProperty_typed()
func (c *RemoveUnusedCodeTest) TestRemoveUnusedPolyfills_prototypeProperty_typed() {
	c.EnableTypeCheck()
	// We no longer use type information to make polyfill removal decisions.
	c.TestRemoveUnusedPolyfills_prototypeProperty_untyped()
}

// Java: RemoveUnusedCodeTest.testRemoveUnusedPolyfills_globalWithPrototypePolyfill_untyped()
func (c *RemoveUnusedCodeTest) TestRemoveUnusedPolyfills_globalWithPrototypePolyfill_untyped() {
	promisePolyfill := "$jscomp.polyfill('Promise', function() {}, 'es6', 'es3');"
	finallyPolyfill :=
		"$jscomp.polyfill('Promise.prototype.finally', function() {}, 'es8', 'es3');"

	tester :=
		NewRemoveUnusedCodeTest_PolyfillRemovalTester(c).
			AddExterns(
				NewTestExternsBuilder(). //
								AddPromise().
								AddConsole().
								AddExtra(RemoveUnusedCodeTest_JSCOMP_POLYFILL).
								Build()).
			AddPolyfill(promisePolyfill).
			AddPolyfill(finallyPolyfill)

	// Both the base polyfill and the extra method are removed when unused.
	tester.ExpectPolyfillsRemovedTest("console.log();", promisePolyfill, finallyPolyfill)

	// The extra method polyfill is removed if not used, even when the base is retained.
	tester.ExpectPolyfillsRemovedTest("console.log(Promise.resolve());", finallyPolyfill)

	// Promise is guarded by an optional chain.
	tester.ExpectPolyfillsRemovedTest(
		"console.log(Promise?.resolve());", promisePolyfill, finallyPolyfill)

	// Can't remove finally without type information
	tester.ExpectPolyfillsRemovedTest(
		"const p = {finally() {}};\n"+
			"p.finally();\n",
		// NOTE: In reality the Promise.prototype.finally polyfill references Promise, so
		// it would actually prevent the Promise polyfill from being removed.
		promisePolyfill)

	// `finally` is guarded by an optional chain
	tester.ExpectPolyfillsRemovedTest(
		"const p = {finally() {}};\n"+
			"p.finally?.();\n",
		promisePolyfill,
		finallyPolyfill)

	// Retain both the base and the extra method when both are used.
	tester.ExpectNoRemovalTest("console.log(Promise.resolve().finally(() => {}));")

	tester.ExpectPolyfillsRemovedTest(
		"console.log(Promise?.resolve()?.finally(() => {}));",
		promisePolyfill) // finally is not guarded

	tester.ExpectPolyfillsRemovedTest(
		"console.log(Promise?.resolve().finally?.(() => {}));",
		promisePolyfill,
		finallyPolyfill) // both are guarded

	// TODO(b/163394833): Note that this behavior could result in the Promise polyfill being
	// removed, then the Promise.prototype.finally polyfill code having nowhere to hang its value.
	// This is consistent with the way guarding with `&&` would behave, as demonstrated by the
	// following test case.
	// NOTE: In reality the Promise.prototype.finally polyfill references Promise, so
	// it would actually prevent the Promise polyfill from being removed.
	tester.ExpectPolyfillsRemovedTest(
		"console.log(Promise?.resolve().finally(() => {}));",
		promisePolyfill) // only Promise guarded by optional chain

	tester.ExpectPolyfillsRemovedTest(
		"console.log(Promise && Promise.resolve().finally(() => {}));",
		promisePolyfill) // only Promise guarded by optional chain

	tester.ExpectPolyfillsRemovedTest(
		"console.log(Promise.resolve().finally?.(() => {}));",
		finallyPolyfill) // only `finally` guarded by optional chain

	// The base polyfill is removed and the extra method retained if the constructor never shows up
	// anywhere in the source code.  NOTE: this is probably the wrong thing to do.  This situation
	// should only be possible if async function transpilation happens *after* RemoveUnusedCode
	// (since we have an async function).  The fact that we inserted the Promise polyfill in the
	// first case indicates the output language is < ES6.  When that later transpilation occurs, we
	// will end up adding uses of the Promise constructor.  We need to keep this in mind when moving
	// transpilation after optimizations.
	tester.ExpectPolyfillsRemovedTest(
		"async function f() {}\n"+
			"f().finally(() => {});\n",
		promisePolyfill)
}

// Java: RemoveUnusedCodeTest.testRemoveUnusedPolyfills_globalWithPrototypePolyfill_typed()
func (c *RemoveUnusedCodeTest) TestRemoveUnusedPolyfills_globalWithPrototypePolyfill_typed() {
	c.EnableTypeCheck()
	// We no longer use type information to make polyfill removal decisions.
	c.TestRemoveUnusedPolyfills_globalWithPrototypePolyfill_untyped()
}

// Java: RemoveUnusedCodeTest.testRemoveUnusedPolyfills_typedArrayMethod_untyped()
func (c *RemoveUnusedCodeTest) TestRemoveUnusedPolyfills_typedArrayMethod_untyped() {
	typedArrayAtPolyfill :=
		"$jscomp.polyfillTypedArrayMethod('at', function() {}, 'es_2022', 'es5');"
	tester :=
		NewRemoveUnusedCodeTest_PolyfillRemovalTester(c).
			AddExterns(
				NewTestExternsBuilder().
					AddConsole().
					AddExtra(
						RemoveUnusedCodeTest_JSCOMP_POLYFILL,
						"/** @constructor */ function Int8Array() {}",
						"/** @type {function(number):number} */ Int8Array.prototype.at;",
						"/** @constructor */ function Float32Array() {}",
						"/** @type {function(number):number} */ Float32Array.prototype.at;",
						"/** @constructor */ function BigInt64Array() {}",
						"/** @type {function(number):number} */ BigInt64Array.prototype.at;").
					Build()).
			AddPolyfill(typedArrayAtPolyfill)

	// unused polyfill is removed
	tester.ExpectPolyfillsRemovedTest("console.log('.at() is not called')", typedArrayAtPolyfill)
	// used polyfill is not removed - usage via one TypedArray type is sufficient to keep it
	tester.ExpectNoRemovalTest("console.log(new Int8Array().at(0))")
	tester.ExpectNoRemovalTest("console.log(new Float32Array().at(0))")
	tester.ExpectNoRemovalTest("console.log(new BigInt64Array().at(0))")
}

// Java: RemoveUnusedCodeTest.testRemoveUnusedPolyfills_chained()
func (c *RemoveUnusedCodeTest) TestRemoveUnusedPolyfills_chained() {
	weakMapPolyfill :=
		"$jscomp.polyfill('WeakMap', function() { console.log(); }, 'es6', 'es3');"
	// Polyfill of Map depends on WeakMap
	mapPolyfill :=
		"$jscomp.polyfill('Map', function() { new WeakMap(); }, 'es6', 'es3');"
	// Polyfill of Set depends on Map
	setPolyfill := "$jscomp.polyfill('Set', function() { new Map(); }, 'es6', 'es3');"

	tester :=
		NewRemoveUnusedCodeTest_PolyfillRemovalTester(c).
			AddExterns(
				NewTestExternsBuilder().
					AddConsole().
					AddMap().
					AddExtra(
						RemoveUnusedCodeTest_JSCOMP_POLYFILL,
						"/** @constructor */ function Set() {}",
						"/** @constructor */ function WeakMap() {}").
					Build()).
			AddPolyfill(weakMapPolyfill).
			AddPolyfill(mapPolyfill).
			AddPolyfill(setPolyfill)

	// Removes polyfills that are only referenced in other (removed) polyfills' definitions.
	tester.
		InputSource(
			"function unused() { new Set(); }\n"+
				"console.log();\n").
		ExpectSource("console.log();").
		// Unused method gets removed, allowing all 3 polyfills to be removed.
		ExpectPolyfillsRemoved(weakMapPolyfill, mapPolyfill, setPolyfill).
		Test()

	// Chains can be partially removed if just an outer-most symbol is unreferenced.
	tester.ExpectPolyfillsRemovedTest("console.log(new Map());", setPolyfill)

	// Only requires a single reference to the outermost symbol to retain the whole chain.
	tester.ExpectNoRemovalTest("console.log(new Set())")
}

// Java: RemoveUnusedCodeTest.testRemoveUnusedPolyfills_continued()
func (c *RemoveUnusedCodeTest) TestRemoveUnusedPolyfills_continued() {
	externs :=
		c.Externs(NewTestExternsBuilder().AddConsole().AddMap().AddExtra(RemoveUnusedCodeTest_JSCOMP_POLYFILL).Build())

	// Ensure that continuations occur so that retained polyfill definitions are still optimized.
	c.Test(
		externs,
		c.Srcs(
			"$jscomp.polyfill('Map', function() { var x; }, 'es6', 'es3');\n"+
				"console.log(new Map());\n"),
		c.Expected(
			"$jscomp.polyfill('Map', function() {        }, 'es6', 'es3');\n"+
				"console.log(new Map());\n"))
}

// Java: RemoveUnusedCodeTest.testRemoveUnusedPolyfills_collapsedPolyfillFunction()
func (c *RemoveUnusedCodeTest) TestRemoveUnusedPolyfills_collapsedPolyfillFunction() {
	externs :=
		c.Externs(
			NewTestExternsBuilder().
				AddConsole().
				AddMap().
				AddExtra("function $jscomp$polyfill() {}").
				Build())

	// The pass should also work after CollapseProperties.
	c.Test(
		externs,
		c.Srcs(
			"$jscomp$polyfill('Map', function() {}, 'es6', 'es3');\n"+
				"console.log();\n"),
		c.Expected("console.log();"))

	c.TestSame(
		externs,
		c.Srcs(
			"$jscomp$polyfill('Map', function() {}, 'es6', 'es3');\n"+
				"console.log(new Map());\n"))
}

// Java: RemoveUnusedCodeTest.testRemoveUnusedPolyfills_guardedGlobals()
func (c *RemoveUnusedCodeTest) TestRemoveUnusedPolyfills_guardedGlobals() {
	mapPolyfill := "$jscomp.polyfill('Map', function() {}, 'es6', 'es3');"

	tester :=
		NewRemoveUnusedCodeTest_PolyfillRemovalTester(c).
			AddExterns(
				NewTestExternsBuilder().
					AddConsole().
					AddMap().
					AddExtra(RemoveUnusedCodeTest_JSCOMP_POLYFILL, "/** @constructor */ function Promise() {}").
					Build()).
			AddPolyfill(mapPolyfill)

	tester.ExpectPolyfillsRemovedTest(
		"if (typeof Map !== 'undefined') {\n"+
			"  console.log(Map);\n"+
			"}\n",
		mapPolyfill)

	tester.ExpectPolyfillsRemovedTest(
		"if (Map) {\n"+
			"  console.log(Map);\n"+
			"}\n",
		mapPolyfill)

	promisePolyfill := "$jscomp.polyfill('Promise', function() {}, 'es6', 'es3');"
	tester.AddPolyfill(promisePolyfill)

	tester.ExpectPolyfillsRemovedTest(
		"if (typeof Map !== 'undefined') {\n"+
			"  console.log(Map);\n"+
			"  console.log(Promise);\n"+
			"}\n",
		mapPolyfill)
}

// Java: RemoveUnusedCodeTest.testRemoveUnusedPolyfills_guardedStatics()
func (c *RemoveUnusedCodeTest) TestRemoveUnusedPolyfills_guardedStatics() {
	arrayFromPolyfill := "$jscomp.polyfill('Array.from', function() {}, 'es6', 'es3');"
	promisePolyfill := "$jscomp.polyfill('Promise', function() {}, 'es6', 'es3');"
	allSettledPolyfill :=
		"$jscomp.polyfill('Promise.allSettled', function() {}, 'es8', 'es3');"
	symbolPolyfill := "$jscomp.polyfill('Symbol', function() {}, 'es6', 'es3');"
	symbolIteratorPolyfill :=
		"$jscomp.polyfill('Symbol.iterator', function() {}, 'es6', 'es3');"
	tester :=
		NewRemoveUnusedCodeTest_PolyfillRemovalTester(c).
			AddExterns(
				NewTestExternsBuilder().
					AddConsole().
					AddArray().
					AddPromise().
					AddExtra(RemoveUnusedCodeTest_JSCOMP_POLYFILL).
					Build()).
			AddPolyfill(arrayFromPolyfill).
			AddPolyfill(promisePolyfill).
			AddPolyfill(allSettledPolyfill).
			AddPolyfill(symbolPolyfill).
			AddPolyfill(symbolIteratorPolyfill)

	tester.ExpectPolyfillsRemovedTest(
		"if (typeof Array.from !== 'undefined') {\n"+
			"  console.log(Array.from);\n"+
			"}\n",
		arrayFromPolyfill, // guarded & all others unused
		promisePolyfill,
		allSettledPolyfill,
		symbolPolyfill,
		symbolIteratorPolyfill)

	tester.ExpectPolyfillsRemovedTest(
		"var a;\n"+
			"if (Promise && Promise.allSettled) {\n"+
			"  Promise.allSettled(a);\n"+
			"}\n",
		promisePolyfill,    // guarded
		allSettledPolyfill, // guarded
		arrayFromPolyfill,  // this and following unused
		symbolPolyfill,
		symbolIteratorPolyfill)

	tester.ExpectPolyfillsRemovedTest(
		"if (Symbol && Symbol.iterator) {\n"+
			"  console.log(Symbol.iterator);\n"+
			"}\n",
		symbolPolyfill,         // guarded
		symbolIteratorPolyfill, // guarded
		arrayFromPolyfill,      // this and following unused
		promisePolyfill,
		allSettledPolyfill)
}

// Java: RemoveUnusedCodeTest.testRemoveUnusedPolyfills_guardedMethods()
func (c *RemoveUnusedCodeTest) TestRemoveUnusedPolyfills_guardedMethods() {
	arrayFindPolyfill :=
		"$jscomp.polyfill('Array.prototype.find', function() {}, 'es6', 'es3');"
	tester :=
		NewRemoveUnusedCodeTest_PolyfillRemovalTester(c).
			AddExterns(
				NewTestExternsBuilder().AddConsole().AddArray().AddExtra(RemoveUnusedCodeTest_JSCOMP_POLYFILL).Build()).
			AddPolyfill(arrayFindPolyfill)
	tester.ExpectPolyfillsRemovedTest(
		"const arr = [];\n"+
			"if (typeof arr.find !== 'undefined') {\n"+
			"  console.log(arr.find(0));\n"+
			"}\n",
		arrayFindPolyfill)
}

// Java: RemoveUnusedCodeTest.testRemoveUnusedPolyfills_unguardedAndGuarded()
func (c *RemoveUnusedCodeTest) TestRemoveUnusedPolyfills_unguardedAndGuarded() {
	mapPolyfill := "$jscomp.polyfill('Map', function() {}, 'es6', 'es3');"
	tester :=
		NewRemoveUnusedCodeTest_PolyfillRemovalTester(c).
			AddExterns(
				NewTestExternsBuilder().AddConsole().AddMap().AddExtra(RemoveUnusedCodeTest_JSCOMP_POLYFILL).Build()).
			AddPolyfill(mapPolyfill)

	// Map is not removed because it has an unguarded usage.
	tester.ExpectNoRemovalTest(
		"if (typeof Map == 'undefined') {\n" +
			"  console.log(Map);\n" +
			"}\n" +
			"console.log(Map);\n")
}

// Java: RemoveUnusedCodeTest.testRemoveUnusedPatches()
func (c *RemoveUnusedCodeTest) TestRemoveUnusedPatches() {
	mapPatch := "$jscomp.patch('Map', function() {});"
	tester :=
		NewRemoveUnusedCodeTest_PolyfillRemovalTester(c).
			AddExterns(
				NewTestExternsBuilder().AddConsole().AddExtra(RemoveUnusedCodeTest_JSCOMP_PATCH).AddMap().Build()).
			AddPolyfill(mapPatch)

	// unused polyfill is removed
	tester.ExpectPolyfillsRemovedTest("console.log()", mapPatch)
	// used polyfill is not removed
	tester.ExpectNoRemovalTest("console.log(new Map())")

	// Local names shadowing global polyfills are not themselves polyfill references.
	tester.ExpectPolyfillsRemovedTest(
		"console.log(function(Map$jscomp$1) {\n"+
			"  console.log(new Map$jscomp$1());\n"+
			"});\n",
		mapPatch)
}

// Java: RemoveUnusedCodeTest.testRemoveUnusedPatches_guardedUsage()
func (c *RemoveUnusedCodeTest) TestRemoveUnusedPatches_guardedUsage() {
	mapPatch := "$jscomp.patch('Map', function() {});"
	tester :=
		NewRemoveUnusedCodeTest_PolyfillRemovalTester(c).
			AddExterns(
				NewTestExternsBuilder().AddConsole().AddExtra(RemoveUnusedCodeTest_JSCOMP_PATCH).AddMap().Build()).
			AddPolyfill(mapPatch)

	// Map is not removed because it is a patch.
	tester.ExpectNoRemovalTest(
		"if (typeof Map == 'undefined') {\n" +
			"  console.log(Map);\n" +
			"}\n")
}

// Java: RemoveUnusedCodeTest.testNoCatchBinding()
func (c *RemoveUnusedCodeTest) TestNoCatchBinding() {
	c.TestSame("function doNothing() {} try { doNothing(); } catch { doNothing(); }")
	c.TestSame("function doNothing() {} try { throw 0; } catch { doNothing(); }")
	c.TestSame("function doNothing() {} try { doNothing(); } catch { console.log('stuff'); }")
}

// Java: RemoveUnusedCodeTest.testDoNotRemoveSetterAssignmentObject()
func (c *RemoveUnusedCodeTest) TestDoNotRemoveSetterAssignmentObject() {
	c.TestSame(
		"var a = {\n" +
			"  set property(x) {}\n" +
			"};\n" +
			"a.property = 1;\n")
}

// Java: RemoveUnusedCodeTest.testDoNotRemoveSetterAssignmentClass()
func (c *RemoveUnusedCodeTest) TestDoNotRemoveSetterAssignmentClass() {
	c.TestSame(
		"class Class {\n" +
			"  set property(x) {}\n" +
			"}\n" +
			"const a = new Class();\n" +
			"a.property = 1;\n")

	c.TestSame(
		"class Class {\n" +
			"  set property(x) {}\n" +
			"}\n" +
			"new Class().property = 1;\n")

	c.Test(
		"class Class {\n"+
			"  set property(x) {}\n"+
			"}\n"+
			"function foo() {}\n"+
			"foo().property = 1;\n",
		"function foo() {}\n"+
			"foo().property = 1;\n")

	c.Test(
		"class Class {\n"+
			"  set property(x) {}\n"+
			"}\n"+
			"var obj;\n"+
			"obj.property = 1;\n",
		"var obj;\n"+
			"obj.property = 1;\n")
}

// Java: RemoveUnusedCodeTest.testDoNotRemoveStaticSetterAssignment()
func (c *RemoveUnusedCodeTest) TestDoNotRemoveStaticSetterAssignment() {
	c.TestSame(
		"class Class {\n" +
			"  static set property(x) {}\n" +
			"}\n" +
			"Class.property = 1;\n")

	c.Test(
		"class Class {\n"+
			"  static set property(x) {}\n"+
			"}\n"+
			"function foo() {}\n"+
			"foo().property = 1;\n",
		"function foo() {}\n"+
			"foo().property = 1;\n")

	c.Test(
		"class Class {\n"+
			"  static set property(x) {}\n"+
			"}\n"+
			"var obj;\n"+
			"obj.property = 1;\n",
		"var obj;\n"+
			"obj.property = 1;\n")
}

// Java: RemoveUnusedCodeTest.testMethodCallingSetterHasSideEffects()
func (c *RemoveUnusedCodeTest) TestMethodCallingSetterHasSideEffects() {
	c.TestSame(
		"class Class {\n" +
			"  set property(x) {}\n" +
			"}\n" +
			"function foo() { new Class().property = 1; }\n" +
			"foo();\n")

	c.TestSame(
		"class Class {\n" +
			"  static set property(x) {}\n" +
			"}\n" +
			"function foo() { Class.property = 1; }\n" +
			"foo();\n")

	c.TestSame(
		"class Class {\n" +
			"  setProperty(v) { this.property = v; }\n" +
			"  set property(x) {}\n" +
			"}\n" +
			"new Class().setProperty(0);\n")

	c.TestSame(
		"class Class {\n" +
			"  static setProperty(v) { this.property = v; }\n" +
			"  static set property(x) {}\n" +
			"}\n" +
			"Class.setProperty(0);\n")

	c.TestSame(
		"class Class {\n" +
			"  setProperty(v) { Class.property = v; }\n" +
			"  static set property(x) {}\n" +
			"}\n" +
			"Class.setProperty(0);\n")
}

// Java: RemoveUnusedCodeTest.testDoNotRemoveSetter_fromObjectLiteral_inCompoundAssignment_onName()
func (c *RemoveUnusedCodeTest) TestDoNotRemoveSetter_fromObjectLiteral_inCompoundAssignment_onName() {
	c.TestSame(
		"var a = {\n" +
			"  set property(x) {}\n" +
			"};\n" +
			"a.property += 1;\n")
}

// Java: RemoveUnusedCodeTest.testDoNotRemoveSetter_fromObjectLiteral_inCompoundAssignment_onThis()
func (c *RemoveUnusedCodeTest) TestDoNotRemoveSetter_fromObjectLiteral_inCompoundAssignment_onThis() {
	c.TestSame(
		"var a = {\n" +
			"  set property(x) {},\n" +
			"\n" +
			"  method() {\n" +
			"    this.property += 1;\n" +
			" },\n" +
			"};\n" +
			"a.method();\n")
}

// Java: RemoveUnusedCodeTest.testDoNotRemoveSetterCompoundAssignmentClass()
func (c *RemoveUnusedCodeTest) TestDoNotRemoveSetterCompoundAssignmentClass() {
	c.TestSame(
		"class Class {\n" +
			"  set property(x) {}\n" +
			"}\n" +
			"const a = new Class();\n" +
			"a.property += 1;\n")
}

// Java: RemoveUnusedCodeTest.testDoNotRemoveSetter_fromObjectLiteral_inUnaryOp_onName()
func (c *RemoveUnusedCodeTest) TestDoNotRemoveSetter_fromObjectLiteral_inUnaryOp_onName() {
	c.TestSame(
		"var a = {\n" +
			"  set property(x) {}\n" +
			"};\n" +
			"a.property++;\n")
}

// Java: RemoveUnusedCodeTest.testDoNotRemoveSetter_fromObjectLiteral_inUnaryOp_onThis()
func (c *RemoveUnusedCodeTest) TestDoNotRemoveSetter_fromObjectLiteral_inUnaryOp_onThis() {
	c.TestSame(
		"var a = {\n" +
			"  set property(x) {},\n" +
			"\n" +
			"  method() {\n" +
			"    this.property++;\n" +
			"  },\n" +
			"};\n" +
			"a.method();\n")
}

// Java: RemoveUnusedCodeTest.testDoNotRemoveSetterUnaryOperatorClass()
func (c *RemoveUnusedCodeTest) TestDoNotRemoveSetterUnaryOperatorClass() {
	c.TestSame(
		"class Class {\n" +
			"  set property(x) {}\n" +
			"}\n" +
			"const a = new Class();\n" +
			"a.property++;\n")
}

// Java: RemoveUnusedCodeTest.testDoNotRemoveAssignmentIfOtherPropertyIsSetterObject()
func (c *RemoveUnusedCodeTest) TestDoNotRemoveAssignmentIfOtherPropertyIsSetterObject() {
	c.Test(
		"var a = {\n"+
			"  set property(x) {}\n"+
			"};\n"+
			"var b = {\n"+
			"  property: 0\n"+
			"};\n"+
			"b.property = 1;\n",
		"var b = {\n"+
			"  property: 0\n"+
			"};\n"+
			"b.property = 1;\n")

	// Test that this gets cleaned up on a second pass...
	c.Test(
		"var b = {\n"+
			"  property: 0\n"+
			"};\n"+
			"b.property = 1;\n",
		"")
}

// Java: RemoveUnusedCodeTest.testDoNotRemoveAssignmentIfOtherPropertyIsSetterClass()
func (c *RemoveUnusedCodeTest) TestDoNotRemoveAssignmentIfOtherPropertyIsSetterClass() {
	c.Test(
		"class Class {\n"+
			"  set property(x) {}\n"+
			"}\n"+
			"var b = {\n"+
			"  property: 0\n"+
			"};\n"+
			"b.property = 1;\n",
		"var b = {\n"+
			"  property: 0\n"+
			"};\n"+
			"b.property = 1;\n")

	// Test that this gets cleaned up on a second pass...
	c.Test(
		"var b = {\n"+
			"  property: 0\n"+
			"};\n"+
			"b.property = 1;\n",
		"")
}

// Java: RemoveUnusedCodeTest.testRemovePropertyUnrelatedFromSetterObject()
func (c *RemoveUnusedCodeTest) TestRemovePropertyUnrelatedFromSetterObject() {
	c.Test(
		"var a = {\n"+
			"  set property(x) {},\n"+
			"  aUnrelated: 0\n"+
			"};\n"+
			"var b = {\n"+
			"  unrelated: 0\n"+
			"};\n"+
			"a.property = 1;\n"+
			"a.aUnrelated = 1;\n"+
			"b.unrelated = 1;\n",
		"var a = {\n"+
			"  set property(x) {},\n"+
			"  aUnrelated: 0\n"+
			"};\n"+
			"a.property = 1;\n"+
			"a.aUnrelated = 1;\n")
}

// Java: RemoveUnusedCodeTest.testFunctionCallReferencesGetterIsNotRemoved()
func (c *RemoveUnusedCodeTest) TestFunctionCallReferencesGetterIsNotRemoved() {
	c.TestSame(
		"var a = {\n" +
			"  get property() {}\n" +
			"};\n" +
			"function foo() { a.property; }\n" +
			"foo();\n")
}

// Java: RemoveUnusedCodeTest.testFunctionCallReferencesSetterIsNotRemoved()
func (c *RemoveUnusedCodeTest) TestFunctionCallReferencesSetterIsNotRemoved() {
	c.TestSame(
		"var a = {\n" +
			"  set property(v) {}\n" +
			"};\n" +
			"function foo() { a.property = 0; }\n" +
			"foo();\n")
}

// Java: RemoveUnusedCodeTest.testRemoveUnusedGettersAndSetters()
func (c *RemoveUnusedCodeTest) TestRemoveUnusedGettersAndSetters() {
	c.TestSame(
		"class C {\n" +
			"  get usedProperty() {}\n" +
			"  set usedProperty(v) {}\n" +
			"  get unUsedProperty() {}\n" +
			"  set unUsedProperty(v) {}\n" +
			"};\n" +
			"function foo() {\n" +
			"  const c = new C();\n" +
			"  c.usedProperty = 0;\n" +
			"  return c.usedProperty;\n" +
			"}\n" +
			"foo();\n")
}

// Java: RemoveUnusedCodeTest.testRemovalFromRHSOfComma()
func (c *RemoveUnusedCodeTest) TestRemovalFromRHSOfComma() {
	// This is the repro for github issue 3612
	c.Test(
		"function a() {\n"+
			"    var a = {}, b = null;\n"+
			"    a.a = 1,\n"+
			"    b.a = 2, // Note the comma here.\n"+
			"    Object.defineProperties(a, b);\n"+
			"};\n"+
			"alert(a);\n",
		"function a() {\n"+
			"};\n"+
			"alert(a);\n")
}

// Java: RemoveUnusedCodeTest.testRemovalFromRHSOfAND()
func (c *RemoveUnusedCodeTest) TestRemovalFromRHSOfAND() {
	c.Test(
		"function a() {\n"+
			"    var a = {};\n"+
			"    var CONDITION = true;\n"+
			"    CONDITION && Object.defineProperties(a, b);\n"+
			"};\n"+
			"alert(a);\n",
		"function a() {\n"+
			"  var CONDITION = true;\n"+
			"  CONDITION;\n"+
			"};\n"+
			"alert(a);\n")
}

// Java: RemoveUnusedCodeTest.testRemovalFromRHSOfOR()
func (c *RemoveUnusedCodeTest) TestRemovalFromRHSOfOR() {
	c.Test(
		"function a() {\n"+
			"    var a = {};\n"+
			"    var CONDITION = true;\n"+
			"    CONDITION || Object.defineProperties(a, b);\n"+
			"};\n"+
			"alert(a);\n",
		"function a() {\n"+
			"  var CONDITION = true;\n"+
			"  CONDITION;\n"+
			"};\n"+
			"alert(a);\n")
}

// Java: RemoveUnusedCodeTest.testVariableAssignedToItselfInORNode()
func (c *RemoveUnusedCodeTest) TestVariableAssignedToItselfInORNode() {
	// Fix for b/359932022, where variable is assigned to itself in an OR node.
	// 1) the node to which a value is being assigned
	// 2) the value being assigned and the
	// Handle the case in an OR node,where both (1) and (2) are the same variable.
	c.Test("var a = a || {}; a[\"hi\"] = true;", "")
	c.Test("var a = {}; a[\"hi\"] = true;", "") // without the OR node

	// test RHS of OR is not removable
	c.TestSame("var not_removable = 5; var a = a || not_removable ; a[\"hi\"] = true;")
}

// Java: RemoveUnusedCodeTest.testRemovalFromExpression()
func (c *RemoveUnusedCodeTest) TestRemovalFromExpression() {
	c.Test(
		"function a() {\n"+
			"    var a = {};\n"+
			"    var CONDITION = true;\n"+
			"    CONDITION ? Object.defineProperties(a, b) : 'something';\n"+
			"};\n"+
			"alert(a);\n",
		"function a() {\n"+
			"  var CONDITION = true;\n"+
			"  CONDITION ? 0 : 'something';\n"+
			"};\n"+
			"alert(a);\n")
}

// Java: RemoveUnusedCodeTest.testPreserveDestructuringWithObjectRest()
func (c *RemoveUnusedCodeTest) TestPreserveDestructuringWithObjectRest() {
	c.TestSame(
		"const obj = {'one': 1, 'two': 2};\n" +
			"const {['one']: unused, ...remaining} = obj;\n" +
			"console.log(remaining);\n")
}

// Java: RemoveUnusedCodeTest.testArrayReferenceInForLoopNameDeclarations()
func (c *RemoveUnusedCodeTest) TestArrayReferenceInForLoopNameDeclarations() {
	c.TestSame(
		"const arr=[];\n" +
			"for (let i = 0, ref = arr; i<5; i++)\n" +
			"  ref[i] = \"test\";\n" +
			"console.log(arr)\n")
}
