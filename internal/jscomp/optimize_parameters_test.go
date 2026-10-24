// Port of com/google/javascript/jscomp/OptimizeParametersTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/OptimizeParametersTest.java)

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
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"testing"
)

type OptimizeParametersTest struct {
	CompilerTestCaseBase
}

// Java: OptimizeParametersTest()
func NewOptimizeParametersTest() *OptimizeParametersTest {
	c := &OptimizeParametersTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, CompilerTestCase_DEFAULT_EXTERNS+"var alert;")
	return c
}

func TestOptimizeParametersTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewOptimizeParametersTest() })
}

// Java: OptimizeParametersTest.getProcessor(Compiler)
func (c *OptimizeParametersTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewOptimizeParameters(compiler)
}

// Java: OptimizeParametersTest.setUp()
func (c *OptimizeParametersTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableNormalize()
	c.DisableMultistageCompilation()
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.EnableGatherExternProperties()
}

// Java: OptimizeParametersTest.testAliasingAssignment()
func (c *OptimizeParametersTest) TestAliasingAssignment() {
	c.TestSame(
		"/** @constructor */\n" +
			"function MyClass() {\n" +
			"  this.myField = null;\n" +
			"}\n" +
			"\n" +
			"// This assignment creates an alias, so we can't know all of the callers and cannot\n" +
			"// safely optimize away `myArgument`.\n" +
			"MyClass.prototype[\"myMethod\"] =\n" +
			"    MyClass.prototype.myMethod = function (myArgument) {\n" +
			"  if (undefined === myArgument) {\n" +
			"      myArgument = this.myField;\n" +
			"  }\n" +
			"  return \"myMethod with argument: \" + myArgument;\n" +
			"};\n" +
			"\n" +
			"function globalMyMethod(oMyClass) {\n" +
			"// One call to `myMethod` exists, and it doesn't use the optional argument.\n" +
			"  return oMyClass.myMethod();\n" +
			"}\n")
}

// Java: OptimizeParametersTest.nullishCoalesce()
func (c *OptimizeParametersTest) TestNullishCoalesce() {
	c.Test(
		"var f = (function(...p1){            }) ?? (function(...p2){           }); f()",
		"var f = (function(     ){ var p1 = []}) ?? (function(     ){var p2 = []}); f()")
}

// Java: OptimizeParametersTest.testNoRemoval()
func (c *OptimizeParametersTest) TestNoRemoval() {
	c.TestSame("function foo(p1) { } foo(1); foo(2)") // required "remove unused vars"
	c.TestSame("function foo(p1) { } foo(this);")
	c.TestSame("function foo(p1) { } function g() {foo(arguments)}; g();")
	// Can't move a reference to a local.
	c.TestSame("function foo(p1) { use(p1); } function g() {var x = 1; foo(x);}; g();")

	// optional versions
	c.TestSame("function foo(p1) { } foo?.(1); foo?.(2)") // required "remove unused vars"
	c.TestSame("function foo(p1) { } foo?.(this);")
	c.TestSame("function foo(p1) { } function g() {foo?.(arguments)}; g?.();")
	// Can't move a reference to a local.
	c.TestSame("function foo(p1) { use?.(p1); } function g() {var x = 1; foo?.(x);}; g?.();")
}

// Java: OptimizeParametersTest.testNoRemoval_arrowFunctions()
func (c *OptimizeParametersTest) TestNoRemoval_arrowFunctions() {
	c.TestSame("var foo = (p1)=>{ }; foo(1); foo(2)") // required "remove unused vars"
	c.TestSame("var foo = (p1)=>{ }; foo(this);")
	c.TestSame("var foo = (p1)=>{ }; function g() {foo(arguments)}; g();")
	// Can't move a reference to a local.
	c.TestSame("var foo = (p1)=>{ use(p1); }; function g() {var x = 1; foo(x);}; g();")

	// optional versions
	c.TestSame("var foo = (p1)=>{ }; foo?.(1); foo?.(2)") // required "remove unused vars"
	c.TestSame("var foo = (p1)=>{ }; foo?.(this);")
	c.TestSame("var foo = (p1)=>{ }; function g() {foo?.(arguments)}; g?.();")
	// Can't move a reference to a local.
	c.TestSame("var foo = (p1)=>{ use(p1); }; function g() {var x = 1; foo?.(x);}; g?.();")
}

// Java: OptimizeParametersTest.testNoRemovalSpread()
func (c *OptimizeParametersTest) TestNoRemovalSpread() {
	// TODO(johnlenz): make spread removable
	c.TestSame("function f(p1) {} f(...x);")
	c.TestSame("function f(...p1) {} f(...x);")
	c.Test("function f(p1, ...p2) {} f(1, ...x);", "function f(...p2) {var p1 = 1;} f(...x);")
	c.Test("function f(p1, ...p2) {} f?.(1, ...x);", "function f(...p2) {var p1 = 1;} f?.(...x);")

	c.Test(
		"function f(p1, p2) {} f(1, ...x); f(1, 2);",
		"function f(p2) {var p1 = 1;} f(...x); f(2);")
	c.Test(
		"function f(p1, p2) {} f?.(1, ...x); f?.(1, 2);",
		"function f(p2) {var p1 = 1;} f?.(...x); f?.(2);")

	c.TestSame("function f(p1, p2) {} f(1, ...x); f(2, ...y);")
	c.TestSame("function f(p1, p2) {} f(p1, ...args);")
	c.TestSame("function f(p1, p2) {} f.call(thisArg, p1, ...args);")
	c.TestSame("function f(p1, p2) {} f?.call(thisArg, p1, ...args);")
	// Test spread argument with side effects
	c.TestSame(
		"function foo(x) {sideEffect(); return x;}\n" +
			"function f(p1, p2) {}\n" +
			"f(foo(0), ...[foo(1)]);\n" +
			"f(foo(0), 1);\n")
	// Test should not remove arguments following spread
	c.Test(
		"function f(p1, p2, p3) {} f(1, ...[2], 3); f(1, 2, 3);",
		"function f(p2, p3) {var p1 = 1;} f(...[2], 3); f(2, 3);")
	c.Test(
		"function f(p1, p2, p3) {} f?.(1, ...[2], 3); f?.(1, 2, 3);",
		"function f(p2, p3) {var p1 = 1;} f?.(...[2], 3); f?.(2, 3);")
}

// Java: OptimizeParametersTest.testRemovalRest_singleParam()
func (c *OptimizeParametersTest) TestRemovalRest_singleParam() {
	// rest as the first parameter
	c.Test(
		"function f(...p1){          } f();", //
		"function f(     ){var p1=[];} f()")
	c.Test(
		"function f(...p1){           } f(1);", //
		"function f(     ){var p1=[1];} f( )")
	c.Test(
		"function f(...p1){            use(p1)} f(1);",
		"function f(     ){var p1=[1]; use(p1)} f( )")

	c.Test(
		"function f(...p1){                 } f(alert());",
		"function f(     ){var p1=[alert()];} f(       );")
	c.Test(
		"function f(...p1){                 } f?.(alert());",
		"function f(     ){var p1=[alert()];} f?.(       );")
}

// Java: OptimizeParametersTest.testRemovalRest_secondParam()
func (c *OptimizeParametersTest) TestRemovalRest_secondParam() {
	// rest as the second parameter
	c.Test(
		"function f(p1, ...p2){          } f(x); f(y);",
		"function f(p1,      ){var p2=[];} f(x); f(y);")
	c.Test(
		"function f(p1, ...p2){           } f(x, 1); f(y, 1);",
		"function f(p1,      ){var p2=[1];} f(x   ); f(y   );")
	c.Test(
		"function f(p1, ...p2){            use(p2)} f(x, 1); f(y, 1);",
		"function f(p1,      ){var p2=[1]; use(p2)} f(x   ); f(y   );")
	c.Test(
		"function f(p1, ...p2){            use(p2)} f?.(x, 1); f?.(y, 1);",
		"function f(p1,      ){var p2=[1]; use(p2)} f?.(x   ); f?.(y   );")
	c.Test(
		"function f(p1, ...p2){                 } f(x, alert()); f(y, alert());",
		"function f(p1,      ){var p2=[alert()];} f(x         ); f(y         );")
}

// Java: OptimizeParametersTest.testRemovalRest3()
func (c *OptimizeParametersTest) TestRemovalRest3() {
	// Can't move a reference to a local.
	c.TestSame("function f(...p1){} function _g(x) { f(x); f(x); }")
}

// Java: OptimizeParametersTest.testRemovalRest_rollbackPreservesTrailingArgs()
func (c *OptimizeParametersTest) TestRemovalRest_rollbackPreservesTrailingArgs() {
	c.TestSame("function f(...p) { use(p); } function g(x) { f(1, x, 2); f(1, x, 2); }")
	c.Test(
		"function f(a, ...p) { use(a, p); } function g(x) { f(1, 2, x, 3); f(1, 2, x, 3); }",
		"function f(...p) { var a = 1; use(a, p); } function g(x) { f(2, x, 3); f(2, x, 3); }")
}

// Java: OptimizeParametersTest.testRemovalRestWithDestructuring1()
func (c *OptimizeParametersTest) TestRemovalRestWithDestructuring1() {
	c.Test(
		"function f(...[a]){          } f();", //
		"function f(      ){var a; [a]=[];} f()")
	c.Test(
		"function f(...[a]){          } f(1);", //
		"function f(      ){var a; [a]=[1];} f()")
	c.Test(
		"function f(...{a}){          } f();", //
		"function f(      ){var {a}=[];} f()")
	c.Test(
		"function f(...{a}){          } f?.();", //
		"function f(      ){var {a}=[];} f?.()")

	c.Test(
		"function f(...{a}){            } f(1);", //
		"function f(      ){var a; ({a}=[1]);} f( )")
}

// Java: OptimizeParametersTest.testRemoveParamWithDefault1()
func (c *OptimizeParametersTest) TestRemoveParamWithDefault1() {
	c.Test(
		"function f(a = 1){          } f();", //
		"function f(     ){var a = 1;} f()")
	c.Test(
		"function f(a = []){           } f();", //
		"function f(      ){var a = [];} f()")
	c.Test(
		"function f(a = alert()){                } f();",
		"function f(           ){var a = alert();} f()")
	c.Test(
		"function f([a] = []){             } f();", //
		"function f(        ){var a; [a] = [];} f()")
	c.Test(
		"function f([a] = []){             } f?.();", //
		"function f(        ){var a; [a] = [];} f?.()")
	c.Test(
		"function f([a = 1] = []){                 } f();",
		"function f(            ){var a; [a = 1] = [];} f()")
	c.Test(
		"function f([a = 1]){                 } f([]);",
		"function f(       ){var a; [a = 1] = [];} f(  );")
	c.Test(
		"function f([a = 1]){                 } f?.([]);",
		"function f(       ){var a; [a = 1] = [];} f?.(  );")
}

// Java: OptimizeParametersTest.testRemoveParam_defaultValueOverriden()
func (c *OptimizeParametersTest) TestRemoveParam_defaultValueOverriden() {
	c.Test(
		"function f(a = 0){          } f(1);", //
		"function f(     ){var a = 1;} f( )")
	c.Test(
		"function f({a} = 0){           } f([]);", //
		"function f(       ){var a; ({a}=[]);} f(  )")
	c.Test(
		"function f({a = 1} = 0){               } f([]);",
		"function f(           ){var a; ({a = 1}=[]);} f(  )")
	c.Test(
		"function f({a = 1} = 0){               } f?.([]);",
		"function f(           ){var a; ({a = 1}=[]);} f?.(  )")
	c.Test(
		"function f({a = 1}){               } f({});",
		"function f(       ){var a; ({a = 1}={});} f(  )")
	c.Test(
		"function f({a: a = 1}){               } f({});",
		"function f(          ){var a; ({a = 1}={});} f(  )")
}

// Java: OptimizeParametersTest.testRemoveMultipleParams_withDestructuring_withDefaultValue()
func (c *OptimizeParametersTest) TestRemoveMultipleParams_withDestructuring_withDefaultValue() {
	// test removing multiple parameters
	c.Test(
		"var x = 0; function f(a, {b} = {b: 1}) {} f(1, {b:4});", //
		"var x = 0; function f() { var a = 1; var b; ({b} = {b: 4});} f();")
}

// Java: OptimizeParametersTest.testInlineParamWithDefault_referencingOutsideVariable()
func (c *OptimizeParametersTest) TestInlineParamWithDefault_referencingOutsideVariable() {
	c.Test(
		"var x = 0; function f(a = x) {} f(1);", //
		"var x = 0; function f() { var a = 1; } f();")
	c.Test(
		"var x = 0; function f(a = x + x) {} f(1);", //
		"var x = 0; function f() { var a = 1; } f();")
	c.Test(
		"function g(x) { return x; } function f(a = g(0)) {} f(1);", //
		"function g() { var x = 0; return x; } function f() { var a = 1; } f();")
}

// Java: OptimizeParametersTest.testNoRemoveParamWithDefault_whenArgPossiblyUndefined()
func (c *OptimizeParametersTest) TestNoRemoveParamWithDefault_whenArgPossiblyUndefined() {
	// different scopes
	c.TestSame("function f(p = 1){} function _g(x) { f(x); f(x); }")

	// TODO(johnlenz): add logic for adding an undefined check to the body of the function
	// so that this can be inlined into the function body
	c.TestSame("function f(a = 2){} f(alert);")

	c.Test("function f(a = 2){} f(void 0);", "function f(a = 2){} f();")

	// Make sure `sideEffects()` is always evaluated before `x`;
	c.TestSame("var x = 0; function f(a = sideEffects(), b = x) {}; f(void 0, something);")
	c.TestSame("var x = 0; function f(a = sideEffects(), b = x) {}; f?.(void 0, something);")
}

// Java: OptimizeParametersTest.testNoRemoveParam_beingUsedInSubsequentDefault()
func (c *OptimizeParametersTest) TestNoRemoveParam_beingUsedInSubsequentDefault() {
	c.TestSame("function f(a, b = a) { return b; }; f(x);")

	c.TestSame("function f(a, b = foo(a)) { return b; }; f(x);")

	c.TestSame("function f(a, b = (c) => a) { return b; }; f(x);")

	c.TestSame("function f({a}, {b: [{c = a}]}) { return c; }; f(x);")
}

// Java: OptimizeParametersTest.testNoRemoveParam_beingUsedInSubsequentDefault_fromSingleFormal()
func (c *OptimizeParametersTest) TestNoRemoveParam_beingUsedInSubsequentDefault_fromSingleFormal() {
	c.TestSame("function f([a, b = a]) { return b; }; f(x);")

	c.TestSame("function f({a, [a]: b}) { return b; }; f(x);")
}

// Java: OptimizeParametersTest.testNoInlineParam_beingUsedInSubsequentDefault()
func (c *OptimizeParametersTest) TestNoInlineParam_beingUsedInSubsequentDefault() {
	c.TestSame("function f(a, b = a) { return a + b; }; f(1, x);")

	c.TestSame("function f(a, b = (a = 9)) { return a + b; }; f(1, x);")

	// These would actually be fine to inline, but it's hard to detect that in general.
	c.TestSame("function f(a, b = a) { return a + b; }; f(1, 1);")
	c.TestSame("function f(a, b = a) { return a + b; }; f(1);")
}

// Java: OptimizeParametersTest.testNoInlineYieldExpression()
func (c *OptimizeParametersTest) TestNoInlineYieldExpression() {
	c.TestSame(
		"function f(a) { return a; }\n" +
			"function *g() { f(yield 1); }\n" +
			"use(g().next());\n")
}

// Java: OptimizeParametersTest.testNoRemoveYieldExpression()
func (c *OptimizeParametersTest) TestNoRemoveYieldExpression() {
	c.TestSame(
		"function f(a) { }\n" +
			"function *g() { f(yield 1); }\n" +
			"use(g().next());\n")
}

// Java: OptimizeParametersTest.testNoInlineAwaitExpression()
func (c *OptimizeParametersTest) TestNoInlineAwaitExpression() {
	c.TestSame(
		"function f(a) { return a; }\n" +
			"async function g() { f(await 1); }\n" +
			"use(g().next());\n")
}

// Java: OptimizeParametersTest.testNoRemoveAwaitExpression()
func (c *OptimizeParametersTest) TestNoRemoveAwaitExpression() {
	c.TestSame(
		"function f(a) { }\n" +
			"async function g() { f(await 1); }\n" +
			"use(g().next());\n")
}

// Java: OptimizeParametersTest.testSimpleRemoval0()
func (c *OptimizeParametersTest) TestSimpleRemoval0() {
	c.Test(
		"function foo(p1) {       } foo(); foo()", //
		"function foo(  ) {var p1;} foo(); foo()")
	c.Test(
		"function foo(p1) {           } foo(1); foo(1)",
		"function foo(  ) {var p1 = 1;} foo( ); foo( )")
	c.Test(
		"function foo(p1) {           } foo(1,2); foo(1,4)",
		"function foo(  ) {var p1 = 1;} foo(   ); foo(   )")
	c.Test(
		"function foo(p1) { } foo(1,2); foo(3,4);", //
		"function foo(p1) { } foo(1  ); foo(3  );")
	c.Test(
		"function foo(p1) {           } foo(1,x()); foo(1,y())",
		"function foo(  ) {var p1 = 1;} foo(  x()); foo(  y())")
}

// Java: OptimizeParametersTest.testSimpleRemoval1()
func (c *OptimizeParametersTest) TestSimpleRemoval1() {
	// parameter never supplied
	c.Test(
		"var foo = (p1)=>{       }; foo(); foo()", //
		"var foo = (  )=>{var p1;}; foo(); foo()")
	c.Test(
		"let foo = (p1)=>{       }; foo(); foo()", //
		"let foo = (  )=>{var p1;}; foo(); foo()")
	c.Test(
		"const foo = (p1)=>{       }; foo(); foo()", //
		"const foo = (  )=>{var p1;}; foo(); foo()")
	c.TestSame("class foo { /** @usedViaDotConstructor */ constructor(p1) {}} new foo(); new foo()")

	// constant parameter
	c.Test(
		"var foo = (p1)=>{           }; foo(1); foo(1)",
		"var foo = (  )=>{var p1 = 1;}; foo( ); foo( )")
	c.Test(
		"let foo = (p1)=>{           }; foo(1); foo(1)",
		"let foo = (  )=>{var p1 = 1;}; foo( ); foo( )")
	c.Test(
		"const foo = (p1)=>{           }; foo(1); foo(1)",
		"const foo = (  )=>{var p1 = 1;}; foo( ); foo( )")
	c.TestSame(
		"class foo { /** @usedViaDotConstructor */ constructor(p1) {}} new foo(1); new foo(1)")
}

// Java: OptimizeParametersTest.testSimpleRemoval2()
func (c *OptimizeParametersTest) TestSimpleRemoval2() {
	c.Test(
		"function f(p1) {       } new f(); new f()", //
		"function f(  ) {var p1;} new f(); new f()")
	c.Test(
		"function f(p1) {           } new f(1); new f(1)",
		"function f(  ) {var p1 = 1;} new f( ); new f( )")
	c.Test(
		"function f(p1) {           } new f(1,2); new f(1,4)",
		"function f(  ) {var p1 = 1;} new f(   ); new f(   )")
	c.Test(
		"function f(p1) { } new f(1,2); new f(3,4);", //
		"function f(p1) { } new f(1  ); new f(3  );")
	c.Test(
		"function f(p1) {           } new f(1,x()); new f(1,y())",
		"function f(  ) {var p1 = 1;} new f(  x()); new f(  y())")
	c.TestSame("/** @usedViaDotConstructor */ function f(p1) {} new f(); new f()")
}

// Java: OptimizeParametersTest.testSimpleRemovalInstanceof()
func (c *OptimizeParametersTest) TestSimpleRemovalInstanceof() {
	c.Test(
		"function f(p1) {       } x instanceof f; new f(); new f()",
		"function f(  ) {var p1;} x instanceof f; new f(); new f()")
	c.Test(
		"function f(p1) {           } x instanceof f; new f(1); new f(1)",
		"function f(  ) {var p1 = 1;} x instanceof f; new f( ); new f( )")
}

// Java: OptimizeParametersTest.testSimpleRemovalTypeof()
func (c *OptimizeParametersTest) TestSimpleRemovalTypeof() {
	c.Test(
		"function f(p1) {       } typeof f != 'undefined' && f();",
		"function f(  ) {var p1;} typeof f != 'undefined' && f();")
	c.Test(
		"function f(p1) {           } typeof f != 'undefined'; f(1);",
		"function f(  ) {var p1 = 1;} typeof f != 'undefined'; f();")
}

// Java: OptimizeParametersTest.testSimpleRemoval4()
func (c *OptimizeParametersTest) TestSimpleRemoval4() {
	c.Test(
		"function f(p1) {       } f.prop = 1; new f(); new f()",
		"function f(  ) {var p1;} f.prop = 1; new f(); new f()")
	c.TestSame("/** @usedViaDotConstructor */ function f(p1) {} f.prop = 1; new f(); new f()")

	c.Test(
		"function f(p1) {           } f.prop = 1; new f(1); new f(1)",
		"function f(  ) {var p1 = 1;} f.prop = 1; new f( ); new f( )")
	c.TestSame("/** @usedViaDotConstructor */ function f(p1) {} f.prop = 1; new f(1); new f(1)")

	c.Test(
		"function f(p1) {       } f['prop'] = 1; new f(); new f()",
		"function f(  ) {var p1;} f['prop'] = 1; new f(); new f()")
	c.TestSame("/** @usedViaDotConstructor */ function f(p1) {} f['prop'] = 1; new f(); new f()")

	c.Test(
		"function f(p1) {           } f['prop'] = 1; new f(1); new f(1)",
		"function f(  ) {var p1 = 1;} f['prop'] = 1; new f( ); new f( )")
	c.TestSame("/** @usedViaDotConstructor */ function f(p1) {} f['prop'] = 1; new f(1); new f(1)")
}

// Java: OptimizeParametersTest.testSimpleRemovalAsync()
func (c *OptimizeParametersTest) TestSimpleRemovalAsync() {
	// parameter never supplied
	c.Test(
		"var f = async function (p1) {       }; f(); f()",
		"var f = async function (  ) {var p1;}; f(); f()")
	c.Test(
		"let f = async function (p1) {       }; f(); f()",
		"let f = async function (  ) {var p1;}; f(); f()")
	c.Test(
		"const f = async function (p1) {       }; f(); f()",
		"const f = async function (  ) {var p1;}; f(); f()")

	// constant parameter
	c.Test(
		"var f = async function (p1) {          }; f(1); f(1)",
		"var f = async function (  ) {var p1 = 1;}; f( ); f( )")
	c.Test(
		"let f = async function (p1) {           }; f(1); f(1)",
		"let f = async function (  ) {var p1 = 1;}; f( ); f( )")
	c.Test(
		"const f = async function (p1) {          }; f(1); f(1)",
		"const f = async function (  ) {var p1 = 1;}; f( ); f( )")
}

// Java: OptimizeParametersTest.testSimpleRemovalGenerator()
func (c *OptimizeParametersTest) TestSimpleRemovalGenerator() {
	// parameter never supplied
	c.Test(
		"var f = function * (p1) {       }; f(); f()",
		"var f = function * (  ) {var p1;}; f(); f()")
	c.Test(
		"let f = function * (p1) {       }; f(); f()",
		"let f = function * (  ) {var p1;}; f(); f()")
	c.Test(
		"const f = function * (p1) {       }; f(); f()",
		"const f = function * (  ) {var p1;}; f(); f()")

	// constant parameter
	c.Test(
		"var f = function * (p1) {          }; f(1); f(1)",
		"var f = function * (  ) {var p1 = 1;}; f( ); f( )")
	c.Test(
		"let f = function * (p1) {           }; f(1); f(1)",
		"let f = function * (  ) {var p1 = 1;}; f( ); f( )")
	c.Test(
		"const f = function * (p1) {          }; f(1); f(1)",
		"const f = function * (  ) {var p1 = 1;}; f( ); f( )")
}

// Java: OptimizeParametersTest.testNotAFunction()
func (c *OptimizeParametersTest) TestNotAFunction() {
	c.TestSame("var x = 1; x; x = 2")
}

// Java: OptimizeParametersTest.testRemoveOneOptionalNamedFunction()
func (c *OptimizeParametersTest) TestRemoveOneOptionalNamedFunction() {
	c.Test("function foo(p1) { } foo()", "function foo() {var p1} foo()")
}

// Java: OptimizeParametersTest.testDifferentScopes()
func (c *OptimizeParametersTest) TestDifferentScopes() {
	c.Test(
		"function f(a, b) {} f(1, 2); f(1, 3);\n"+
			"function h() {function g(a) {} g(4); g(5);} f(1, 2);\n",
		"function f(b) {var a = 1} f(2); f(3);\n"+
			"function h() {function g(a) {} g(4); g(5);} f(2);\n")
}

// Java: OptimizeParametersTest.testOptimizeOnlyImmutableValues1()
func (c *OptimizeParametersTest) TestOptimizeOnlyImmutableValues1() {
	c.Test(
		"function foo(a) {}; foo(undefined);", //
		"function foo() {var a = undefined}; foo()")
}

// Java: OptimizeParametersTest.testOptimizeOnlyImmutableValues2()
func (c *OptimizeParametersTest) TestOptimizeOnlyImmutableValues2() {
	c.Test(
		"function foo(a) {}; foo(null);", //
		"function foo() {var a = null}; foo()")
	c.Test(
		"function foo(a) {}; foo(1);", //
		"function foo() {var a = 1}; foo()")
	c.Test(
		"function foo(a) {}; foo('abc');", //
		"function foo() {var a = 'abc'}; foo()")

	c.Test(
		"var foo = function(a) {}; foo(undefined);",
		"var foo = function() {var a = undefined}; foo()")
	c.Test(
		"var foo = function(a) {}; foo(null);", //
		"var foo = function() {var a = null}; foo()")
	c.Test(
		"var foo = function(a) {}; foo(1);", //
		"var foo = function() {var a = 1}; foo()")
	c.Test(
		"var foo = function(a) {}; foo('abc');", //
		"var foo = function() {var a = 'abc'}; foo()")
}

// Java: OptimizeParametersTest.testOptimizeOnlyImmutableValues3()
func (c *OptimizeParametersTest) TestOptimizeOnlyImmutableValues3() {
	// "var a = null;" gets inserted after the declaration of 'goo' so the tree stays normalized.
	c.Test(
		"function foo(a) {\n"+
			"  function goo() {}\n"+
			"  goo(a);\n"+
			"};\n"+
			"foo(null);\n",
		"function foo() {\n"+
			"  function goo() {}\n"+
			"  var a = null;\n"+
			"  goo(a);\n"+
			"};\n"+
			"foo();\n")

	c.Test(
		"function foo(a) {\n"+
			"  function goo() {}\n"+
			"  function boo() {}\n"+
			"  goo(a);\n"+
			"};\n"+
			"foo(null);\n",
		"function foo() {\n"+
			"  function goo() {}\n"+
			"  function boo() {}\n"+
			"  var a = null;\n"+
			"  goo(a);\n"+
			"};\n"+
			"foo();\n")
}

// Java: OptimizeParametersTest.testOptimizeOnlyImmutableValues4()
func (c *OptimizeParametersTest) TestOptimizeOnlyImmutableValues4() {
	c.Test(
		"function foo(a) {\n"+
			"  function goo() { return a; }\n"+
			"};\n"+
			"foo(null);\n",
		"function foo() {\n"+
			"  function goo() { return a; }\n"+
			"  var a = null;\n"+
			"};\n"+
			"foo();\n")
}

// Java: OptimizeParametersTest.testRemoveOneOptionalVarAssignment()
func (c *OptimizeParametersTest) TestRemoveOneOptionalVarAssignment() {
	c.Test(
		"var foo = function (p1) { }; foo()", //
		"var foo = function () {var p1}; foo()")
}

// Java: OptimizeParametersTest.testDoOptimizeCall()
func (c *OptimizeParametersTest) TestDoOptimizeCall() {
	c.TestSame("var foo = function () {}; foo(); foo.call();")
	c.TestSame("var foo = function () {}; foo(); foo?.call();")

	// TODO(johnlenz): support removing unused "this" from .call
	c.TestSame("var foo = function () {}; foo(); foo.call(this);")
	c.TestSame("var foo = function () {}; foo(); foo?.call(this);")

	c.Test(
		"var foo = function (a, b) {}; foo(1); foo.call(this, 1);",
		"var foo = function () {var a = 1;var b;}; foo(); foo.call(this);")
	c.Test(
		"var foo = function (a, b) {}; foo?.(1); foo?.call(this, 1);",
		"var foo = function () {var a = 1;var b;}; foo?.(); foo?.call(this);")

	c.TestSame("var foo = function () {}; foo(); foo.call(null);")
	c.TestSame("var foo = function () {}; foo(); foo?.call(null);")

	c.Test(
		"var foo = function (a, b) {}; foo(1); foo.call(null, 1);",
		"var foo = function () {var a = 1;var b;}; foo(); foo.call(null);")

	c.TestSame("var foo = function () {}; foo.call();")

	c.TestSame("var foo = function () {}; foo.call(this);")
	c.Test(
		"var foo = function (a) {}; foo.call(this, 1);",
		"var foo = function () {var a = 1;}; foo.call(this);")
	c.TestSame("var foo = function () {}; foo.call(null);")
	c.Test(
		"var foo = function (a) {}; foo.call(null, 1);",
		"var foo = function () {var a = 1;}; foo.call(null);")
	c.Test(
		"var foo = function (a) {}; foo?.call(null, 1);",
		"var foo = function () {var a = 1;}; foo?.call(null);")
	c.TestSame("var foo = function (a, b) {}; foo.call(...arr);")
	c.TestSame("var foo = function (a, b, c) {}; foo.call(...arr, 1, 2);")
	c.TestSame("var foo = function (a, b) {}; foo?.call(...arr);")
	c.TestSame("var foo = function (p1, p2) {}; foo(p1, ...args);")
	c.TestSame("var foo = function (p1, p2) {}; foo.call(this, p1, ...args);")
	c.TestSame("var foo = function (p1, p2) {}; foo?.call(this, p1, ...args);")
}

// Java: OptimizeParametersTest.testDoOptimizeApply()
func (c *OptimizeParametersTest) TestDoOptimizeApply() {
	c.TestSame("var foo = function () {}; foo(); foo.apply();")
	c.TestSame("var foo = function () {}; foo(); foo.apply(this);")
	c.TestSame("var foo = function (a, b) {}; foo(1); foo.apply(this, 1);")
	c.TestSame("var foo = function () {}; foo(); foo.apply(null);")
	c.TestSame("var foo = function (a, b) {}; foo(1); foo.apply(null, []);")

	c.TestSame("var foo = function () {}; foo.apply();")
	c.TestSame("var foo = function () {}; foo.apply(this);")
	c.TestSame("var foo = function (a, b) {}; foo.apply(this, 1);")
	c.TestSame("var foo = function () {}; foo.apply(null);")
	c.TestSame("var foo = function (a, b) {}; foo.apply(null, []);")

	// optional versions of the above tests
	c.TestSame("var foo = function () {}; foo(); foo?.apply();")
	c.TestSame("var foo = function () {}; foo(); foo?.apply(this);")
	c.TestSame("var foo = function (a, b) {}; foo(1); foo?.apply(this, 1);")
	c.TestSame("var foo = function () {}; foo(); foo?.apply(null);")
	c.TestSame("var foo = function (a, b) {}; foo(1); foo?.apply(null, []);")

	c.TestSame("var foo = function () {}; foo?.apply();")
	c.TestSame("var foo = function () {}; foo?.apply(this);")
	c.TestSame("var foo = function (a, b) {}; foo?.apply(this, 1);")
	c.TestSame("var foo = function () {}; foo?.apply(null);")
	c.TestSame("var foo = function (a, b) {}; foo?.apply(null, []);")
}

// Java: OptimizeParametersTest.testRemoveOneOptionalExpressionAssign()
func (c *OptimizeParametersTest) TestRemoveOneOptionalExpressionAssign() {
	c.Test(
		"var foo; foo = function (p1) { }; foo()", //
		"var foo; foo = function () { var p1; }; foo()")
}

// Java: OptimizeParametersTest.testRemoveOneOptionalOneRequired()
func (c *OptimizeParametersTest) TestRemoveOneOptionalOneRequired() {
	c.Test(
		"function foo(p1, p2) { } foo(1); foo(2)", //
		"function foo(p1) {var p2} foo(1); foo(2)")
}

// Java: OptimizeParametersTest.testRemoveOneOptionalMultipleCalls()
func (c *OptimizeParametersTest) TestRemoveOneOptionalMultipleCalls() {
	c.Test(
		"function foo(p1, p2) { } foo(1); foo(2); foo()",
		"function foo(p1) {var p2} foo(1); foo(2); foo()")
	c.Test(
		"function foo(p1, p2) { } foo?.(1); foo(2); foo?.()",
		"function foo(p1) {var p2} foo?.(1); foo(2); foo?.()")
}

// Java: OptimizeParametersTest.testRemoveOneOptionalMultiplePossibleDefinition()
func (c *OptimizeParametersTest) TestRemoveOneOptionalMultiplePossibleDefinition() {
	src :=
		"var goog = {};\n" +
			"goog.foo = function (p1, p2) { };\n" +
			"goog.foo = function (q1, q2) { };\n" +
			"goog.foo = function (r1, r2) { };\n" +
			"goog.foo(1); goog.foo(2); goog.foo()\n"
	result :=
		"var goog = {};\n" +
			"goog.foo = function (p1) { var p2; };\n" +
			"goog.foo = function (q1) { var q2; };\n" +
			"goog.foo = function (r1) { var r2; };\n" +
			"goog.foo(1); goog.foo(2); goog.foo()\n"
	c.Test(src, result)
}

// Java: OptimizeParametersTest.testRemoveTwoOptionalMultiplePossibleDefinition()
func (c *OptimizeParametersTest) TestRemoveTwoOptionalMultiplePossibleDefinition() {
	src :=
		"var goog = {};\n" +
			"goog.foo = function (p1, p2, p3, p4) { };\n" +
			"goog.foo = function (q1, q2, q3, q4) { };\n" +
			"goog.foo = function (r1, r2, r3, r4) { };\n" +
			"goog.foo(1,0); goog.foo(2,1); goog.foo()\n"
	result :=
		"var goog = {};\n" +
			"goog.foo = function (p1, p2) { var p3; var p4 };\n" +
			"goog.foo = function (q1, q2) { var q3; var q4 };\n" +
			"goog.foo = function (r1, r2) { var r3; var r4 };\n" +
			"goog.foo(1,0); goog.foo(2,1); goog.foo()\n"
	c.Test(src, result)
}

// Java: OptimizeParametersTest.testMultipleCalls()
func (c *OptimizeParametersTest) TestMultipleCalls() {
	src :=
		"function f(p1, p2, p3, p4) { };\n" +
			"f(1,0); f(2,1); f()\n"
	result :=
		"function f(p1, p2) { var p3; var p4 };\n" +
			"f(1,0); f(2,1); f()\n"
	c.Test(src, result)
}

// Java: OptimizeParametersTest.testConstructorOptArgsNotRemoved()
func (c *OptimizeParametersTest) TestConstructorOptArgsNotRemoved() {
	src :=
		"/** @constructor */\n" +
			"var goog = function(){};\n" +
			"goog.prototype.foo = function(a,b) {};\n" +
			"goog.prototype.bar = function(a) {};\n" +
			"goog.bar.inherits(goog.foo);\n" +
			"new goog.foo(2,3);\n" +
			"new goog.foo(1,2);\n"
	c.TestSame(src)
}

// Java: OptimizeParametersTest.testMultipleUnknown()
func (c *OptimizeParametersTest) TestMultipleUnknown() {
	src :=
		"var goog1 = {};\n" +
			"goog1.foo = function () { };\n" +
			"var goog2 = {};\n" +
			"goog2.foo = function (p1) { };\n" +
			"var x = getGoog();\n" +
			"x.foo()\n"
	result :=
		"var goog1 = {};\n" +
			"goog1.foo = function () { };\n" +
			"var goog2 = {};\n" +
			"goog2.foo = function () { var p1; };\n" +
			"var x = getGoog();\n" +
			"x.foo()\n"
	c.Test(src, result)
}

// Java: OptimizeParametersTest.testSingleUnknown()
func (c *OptimizeParametersTest) TestSingleUnknown() {
	src :=
		"var goog2 = {};\n" +
			"goog2.foo = function (p1) { };\n" +
			"var x = getGoog();x.foo()\n"

	expected :=
		"var goog2 = {};\n" +
			"goog2.foo = function () { var p1 };\n" +
			"var x = getGoog();x.foo()\n"
	c.Test(src, expected)
}

// Java: OptimizeParametersTest.testRemoveVarArg()
func (c *OptimizeParametersTest) TestRemoveVarArg() {
	c.Test(
		"function foo(p1, var_args) { } foo(1); foo(2)",
		"function foo(p1) { var var_args } foo(1); foo(2)")
}

// Java: OptimizeParametersTest.testAliasMethodsDontGetOptimize()
func (c *OptimizeParametersTest) TestAliasMethodsDontGetOptimize() {
	src :=
		"var foo = function(a, b) {};\n" +
			"var goog = {};\n" +
			"goog.foo = foo;\n" +
			"goog.prototype.bar = goog.foo;\n" +
			"new goog().bar(1,2);\n" +
			"foo(2);\n"
	c.TestSame(src)
}

// Java: OptimizeParametersTest.testAliasMethodsDontGetOptimize2()
func (c *OptimizeParametersTest) TestAliasMethodsDontGetOptimize2() {
	src :=
		"var foo = function(a, b) {};\n" +
			"var bar = foo;foo(1);bar(2,3);\n"
	c.TestSame(src)
}

// Java: OptimizeParametersTest.testAliasMethodsDontGetOptimize3()
func (c *OptimizeParametersTest) TestAliasMethodsDontGetOptimize3() {
	src :=
		"var array = {};\n" +
			"array[0] = function(a, b) {};\n" +
			"var foo = array[0]; // foo should be marked as aliased.\n" +
			"foo(1);\n"
	c.TestSame(src)

	srcOptChain :=
		"var array = {};\n" +
			"array[0] = function(a, b) {};\n" +
			"var foo = array[0]; // foo should be marked as aliased.\n" +
			"foo?.(1);\n"
	c.TestSame(srcOptChain)
}

// Java: OptimizeParametersTest.testAliasMethodsDontGetOptimize4()
func (c *OptimizeParametersTest) TestAliasMethodsDontGetOptimize4() {
	// Don't change the call to baz as it has been aliased.

	c.Test(
		"function foo(bar) {};\n"+
			"var baz = function(a) {};\n"+
			"baz(1);\n"+
			"foo(baz);\n",
		"function foo() {var bar = baz};\n"+
			"var baz = function(a) {};\n"+
			"baz(1);\n"+
			"foo();\n")
}

// Java: OptimizeParametersTest.testMethodsDefinedInArraysDontGetOptimized()
func (c *OptimizeParametersTest) TestMethodsDefinedInArraysDontGetOptimized() {
	src :=
		"var array = [true, function (a) {}];\n" +
			"array[1](1)\n"
	c.TestSame(src)

	srcOptChain :=
		"var array = [true, function (a) {}];\n" +
			"array[1]?.(1)\n"
	c.TestSame(srcOptChain)
}

// Java: OptimizeParametersTest.testMethodsDefinedInObjectDontGetOptimized()
func (c *OptimizeParametersTest) TestMethodsDefinedInObjectDontGetOptimized() {
	src :=
		"var object = { foo: function bar() {} };\n" +
			"object.foo(1)\n"
	c.TestSame(src)
	src =
		"var object = { foo: function bar() {} };\n" +
			"object['foo'](1)\n"
	c.TestSame(src)
	src =
		"var object = { foo: function bar() {} };\n" +
			"object['foo']?.(1)\n"
	c.TestSame(src)
}

// Java: OptimizeParametersTest.testRemoveConstantArgument()
func (c *OptimizeParametersTest) TestRemoveConstantArgument() {
	// Remove only one parameter
	c.Test(
		"function foo(p1, p2) {}; foo(1,2); foo(2,2);",
		"function foo(p1) {var p2 = 2}; foo(1); foo(2)")

	// @noinline prevents constant inlining
	c.TestSame("/** @noinline */ function foo(p1, p2) {}; foo(1,2); foo(2,2);")

	// Remove nothing
	c.TestSame("function foo(p1, p2) {}; foo(1); foo(2,3);")

	// Remove middle parameter
	c.Test(
		"function foo(a,b,c){}; foo(1, 2, 3); foo(1, 2, 4); foo(2, 2, 3)",
		"function foo(a,c){var b=2}; foo(1, 3); foo(1, 4); foo(2, 3)")

	// Number are equals
	c.Test("function foo(a) {}; foo(1); foo(1.0);", "function foo() {var a = 1;}; foo(); foo();")

	// A more OO test
	c.Test(
		"/** @constructor */\n"+
			"function Person() {}; Person.prototype.run = function(a, b) {};\n"+
			"Person.run(1, 'a'); Person.run(2, 'a');\n",
		"/** @constructor */\n"+
			"function Person() {}; Person.prototype.run = function(a) {var b = 'a'};\n"+
			"Person.run(1); Person.run(2);\n")
}

// Java: OptimizeParametersTest.testCanDeleteArgumentsAtAnyPosition()
func (c *OptimizeParametersTest) TestCanDeleteArgumentsAtAnyPosition() {
	// Argument removed in middle and end
	src :=
		"function foo(a,b,c,d,e) {};\n" +
			"foo(1,2,3,4,5);\n" +
			"foo(2,2,4,4,5);\n"
	expected :=
		"function foo(a,c) {var b=2; var d=4; var e=5;};\n" +
			"foo(1,3);\n" +
			"foo(2,4);\n"
	c.Test(src, expected)
}

// Java: OptimizeParametersTest.testNoOptimizationForExternsFunctions()
func (c *OptimizeParametersTest) TestNoOptimizationForExternsFunctions() {
	c.TestSame("function _foo(x, y, z){}; _foo(1);")
}

// Java: OptimizeParametersTest.testNoOptimizationForGoogExportSymbol()
func (c *OptimizeParametersTest) TestNoOptimizationForGoogExportSymbol() {
	c.TestSame(
		"goog.exportSymbol('foo', foo);\n" +
			"function foo(x, y, z){}; foo(1);\n")
}

// Java: OptimizeParametersTest.testNoArgumentRemovalNonEqualNodes()
func (c *OptimizeParametersTest) TestNoArgumentRemovalNonEqualNodes() {
	c.TestSame("function foo(a){}; foo('bar'); foo('baz');")
	c.TestSame("function foo(a){}; foo(1.0); foo(2.0);")
	c.TestSame("function foo(a){}; foo(true); foo(false);")
	c.TestSame("var a = 1, b = 2; function foo(a){}; foo(a); foo(b);")
	c.TestSame("function foo(a){}; foo(/&/g); foo(/</g);")
}

// Java: OptimizeParametersTest.testFunctionPassedAsParam()
func (c *OptimizeParametersTest) TestFunctionPassedAsParam() {
	c.Test(
		"/** @constructor */ function person() {};\n"+
			"person.prototype.run = function(a, b) {};\n"+
			"person.prototype.walk = function() {};\n"+
			"person.prototype.foo = function() { this.run(this.walk, 0.1); };\n"+
			"person.foo();\n",
		"/** @constructor */ function person() {};\n"+
			"person.prototype.run = function(a) { var b = 0.1; };\n"+
			"person.prototype.walk = function() {};\n"+
			"person.prototype.foo = function() { this.run(this.walk); };\n"+
			"person.foo();\n")
}

// Java: OptimizeParametersTest.testCallIsIgnore()
func (c *OptimizeParametersTest) TestCallIsIgnore() {
	c.Test(
		"var goog;\n"+
			"goog.foo = function(a, opt) {};\n"+
			"var bar = function(){goog.foo.call(this, 1)};\n"+
			"goog.foo(1);\n",
		"var goog;\n"+
			"goog.foo = function() {var a = 1;var opt;};\n"+
			"var bar = function(){goog.foo.call(this)};\n"+
			"goog.foo();\n")
}

// Java: OptimizeParametersTest.testApplyIsIgnore()
func (c *OptimizeParametersTest) TestApplyIsIgnore() {
	c.TestSame(
		"var goog;\n" +
			"goog.foo = function(a, opt) {};var bar = function(){goog.foo.apply(this, 1)};goog.foo(1);\n")
}

// Java: OptimizeParametersTest.testFunctionWithReferenceToArgumentsShouldNotBeOptimized()
func (c *OptimizeParametersTest) TestFunctionWithReferenceToArgumentsShouldNotBeOptimized() {
	c.TestSame("function foo(a,b,c) { return arguments.size; }; foo(1);")
	c.TestSame("var foo = function(a,b,c) { return arguments.size }; foo(1);")
	c.TestSame("var foo = function bar(a,b,c) { return arguments.size }; foo(2);")
}

// Java: OptimizeParametersTest.testClassMemberWithReferenceToArgumentsShouldNotBeOptimize()
func (c *OptimizeParametersTest) TestClassMemberWithReferenceToArgumentsShouldNotBeOptimize() {
	c.TestSame(
		"class C {\n" +
			"  constructor() {\n" +
			"  }\n" +
			"  setValue(value) {\n" +
			"    if (!arguments.length) {\n" +
			"      return 0;\n" +
			"    }\n" +
			"    return value;\n" +
			"  }\n" +
			"}\n" +
			"var c = new C();\n" +
			"alert(c.setValue(42));\n")
}

// Java: OptimizeParametersTest.testFunctionWithTwoNames()
func (c *OptimizeParametersTest) TestFunctionWithTwoNames() {
	c.TestSame("var foo = function bar(a,b) {};")
	c.TestSame("var foo = function bar(a,b) {}; foo(1)")
	c.TestSame("var foo = function bar(a,b) {}; foo(1); foo(2)")
}

// Java: OptimizeParametersTest.testRecursion()
func (c *OptimizeParametersTest) TestRecursion() {
	c.Test(
		"var foo = function (a,b) {foo(1, b)}; foo(1, 2)",
		"var foo = function (b) {var a=1; foo(b)}; foo(2)")
}

// Java: OptimizeParametersTest.testConstantArgumentsToConstructorCanBeOptimized()
func (c *OptimizeParametersTest) TestConstantArgumentsToConstructorCanBeOptimized() {
	src :=
		"function foo(a) {};\n" +
			"var bar = new foo(1);\n"
	expected :=
		"function foo() {var a=1;};\n" +
			"var bar = new foo();\n"
	c.Test(src, expected)
}

// Java: OptimizeParametersTest.testOptionalArgumentsToConstructorCanBeOptimized()
func (c *OptimizeParametersTest) TestOptionalArgumentsToConstructorCanBeOptimized() {
	src :=
		"function foo(a) {};\n" +
			"var bar = new foo();\n"
	expected :=
		"function foo() {var a;};\n" +
			"var bar = new foo();\n"
	c.Test(src, expected)
}

// Java: OptimizeParametersTest.testRegexesCanBeInlined()
func (c *OptimizeParametersTest) TestRegexesCanBeInlined() {
	c.Test("function foo(a) {}; foo(/abc/);", "function foo() {var a = /abc/}; foo();")
}

// Java: OptimizeParametersTest.testConstructorUsedAsFunctionCanBeOptimized()
func (c *OptimizeParametersTest) TestConstructorUsedAsFunctionCanBeOptimized() {
	src :=
		"function foo(a) {};\n" +
			"var bar = new foo(1);\n" +
			"foo(1);\n"
	expected :=
		"function foo() {var a=1;};\n" +
			"var bar = new foo();\n" +
			"foo();\n"
	c.Test(src, expected)
}

// Java: OptimizeParametersTest.testDoNotOptimizeConstructorWhenArgumentsAreNotEqual()
func (c *OptimizeParametersTest) TestDoNotOptimizeConstructorWhenArgumentsAreNotEqual() {
	c.TestSame(
		"function Foo(a) {};\n" +
			"var bar = new Foo(1);\n" +
			"var baz = new Foo(2);\n")
}

// Java: OptimizeParametersTest.testDoNotOptimizeArrayElements()
func (c *OptimizeParametersTest) TestDoNotOptimizeArrayElements() {
	c.TestSame("var array = [function (a, b) {}];")
	c.TestSame("var array = [function f(a, b) {}]")

	c.TestSame(
		"var array = [function (a, b) {}];\n" +
			"array[0](1, 2);\n" +
			"array[0](1);\n")

	c.TestSame(
		"var array = [];\n" +
			"function foo(a, b) {};\n" +
			"array[0] = foo;\n")
}

// Java: OptimizeParametersTest.testOptimizeThis1()
func (c *OptimizeParametersTest) TestOptimizeThis1() {
	src :=
		"var bar = function (a, b) {};\n" +
			"function foo() {\n" +
			"  this.bar = function (a, b) {};\n" +
			"  this.bar(3);\n" +
			"  bar(2);\n" +
			"}\n"
	expected :=
		"var bar = function () {var a = 2;var b;};\n" +
			"function foo() {\n" +
			"  this.bar = function () {var a = 3;var b;};\n" +
			"  this.bar();\n" +
			"  bar();\n" +
			"}\n"
	c.Test(src, expected)
}

// Java: OptimizeParametersTest.testOptimizeThis2()
func (c *OptimizeParametersTest) TestOptimizeThis2() {
	src :=
		"function foo() {\n" +
			"  var bar = function (a, b) {};\n" +
			"  this.bar = function (a, b) {};\n" +
			"  this.bar(3);\n" +
			"  bar(2);\n" +
			"}\n"
	expected :=
		"function foo() {\n" +
			"  var bar = function (a, b) {};\n" +
			"  this.bar = function () {var a = 3;var b;};\n" +
			"  this.bar();\n" +
			"  bar(2);\n" +
			"}\n"
	c.Test(src, expected)
}

// Java: OptimizeParametersTest.testDoNotOptimizeWhenArgumentsPassedAsParameter()
func (c *OptimizeParametersTest) TestDoNotOptimizeWhenArgumentsPassedAsParameter() {
	c.TestSame("function foo(a) {}; foo(arguments)")
	c.TestSame("function foo(a) {}; foo(arguments[0])")

	c.Test("function foo(a, b) {}; foo(arguments, 1)", "function foo(a) {var b = 1}; foo(arguments)")

	c.Test("function foo(a, b) {}; foo(arguments)", "function foo(a) {var b}; foo(arguments)")
}

// Java: OptimizeParametersTest.testDoNotOptimizeGoogExportFunctions()
func (c *OptimizeParametersTest) TestDoNotOptimizeGoogExportFunctions() {
	c.TestSame("function foo(a, b) {}; foo(); goog.export_function(foo);")
}

// Java: OptimizeParametersTest.testDoNotOptimizeJSCompiler_renameProperty()
func (c *OptimizeParametersTest) TestDoNotOptimizeJSCompiler_renameProperty() {
	c.TestSame(
		"function JSCompiler_renameProperty(a) {return a};\n" +
			"JSCompiler_renameProperty('a');\n")
}

// Java: OptimizeParametersTest.testMutableValues1()
func (c *OptimizeParametersTest) TestMutableValues1() {
	c.Test("function foo(p1) {} foo()", "function foo() {var p1} foo()")
	c.Test("function foo(p1) {} foo(1)", "function foo() {var p1=1} foo()")
	c.Test("function foo(p1) {} foo([])", "function foo() {var p1=[]} foo()")
	c.Test("function foo(p1) {} foo({})", "function foo() {var p1={}} foo()")
	c.Test("var x;function foo(p1) {} foo(x)", "var x;function foo() {var p1=x} foo()")
	c.Test("var x;function foo(p1) {} foo(x())", "var x;function foo() {var p1=x()} foo()")
	c.Test("var x;function foo(p1) {} foo(new x())", "var x;function foo() {var p1=new x()} foo()")
	c.Test("var x;function foo(p1) {} foo('' + x)", "var x;function foo() {var p1='' + x} foo()")

	c.TestSame("function foo(p1) {} foo(this)")
	c.TestSame("function foo(p1) {} foo(arguments)")
	c.TestSame("function foo(p1) {} foo(function(){})")
	c.TestSame("function foo(p1) {} (function () {var x;foo(x)})()")
}

// Java: OptimizeParametersTest.testMutableValues2()
func (c *OptimizeParametersTest) TestMutableValues2() {
	c.Test("function foo(p1, p2) {} foo(1, 2)", "function foo() {var p1=1; var p2 = 2} foo()")
	c.Test(
		"var x; var y; function foo(p1, p2) {} foo(x(), y())",
		"var x; var y; function foo() {var p1=x(); var p2 = y()} foo()")
}

// Java: OptimizeParametersTest.testMutableValues3()
func (c *OptimizeParametersTest) TestMutableValues3() {
	c.Test(
		"var x; var y; var z;\n"+
			"function foo(p1, p2) {}\n"+
			"foo(x(), y()); foo(x(),y())\n",
		"var x; var y; var z;\n"+
			"function foo() {var p1=x(); var p2=y()}\n"+
			"foo(); foo()\n")
}

// Java: OptimizeParametersTest.testMutableValues4()
func (c *OptimizeParametersTest) TestMutableValues4() {
	// Preserve the ordering of side-effects.
	// If z(), can't be moved into the function then z() may change the value
	// of x and y.
	c.TestSame(
		"var x; var y; var z;\n" +
			"function foo(p1, p2, p3) {}\n" +
			"foo(x(), y(), z()); foo(x(),y(),3)\n")

	// If z(), can't be moved into the function then z() may change the value
	// of x and y.
	c.TestSame(
		"var x; var y; var z;\n" +
			"function foo(p1, p2, p3) {}\n" +
			"foo(x, y(), z()); foo(x,y(),3)\n")

	// Mutable object that can not be effect by side-effects are movable,
	// however.
	c.Test(
		"var x; var y; var z;\n"+
			"function foo(p1, p2, p3) {}\n"+
			"foo([], y(), z()); foo([],y(),3)\n",
		"var x; var y; var z;\n"+
			"function foo(p2, p3) {var p1=[]}\n"+
			"foo(y(), z()); foo(y(),3)\n")
	// Pure literal at the first call site and side-effectful call at the second call site.
	c.TestSame(
		"var x; var y; var z;\n" +
			"function foo(p1, p2, p3) { alert(p1); }\n" +
			"foo(x, y(), 3);\n" +
			"foo(x, y(), z());\n")

	// Mutable variable passed along with an argument that mutates it at the second call site.
	c.TestSame(
		"var out = [];\n" +
			"var policy = 'STRICT';\n" +
			"function mutateAndReturn() { policy = 'LENIENT'; return {}; }\n" +
			"function render(p, opts) { out.push(p); }\n" +
			"render(policy, {});\n" +
			"render(policy, mutateAndReturn());\n")
}

// Java: OptimizeParametersTest.testMutableValues5()
func (c *OptimizeParametersTest) TestMutableValues5() {
	c.Test(
		"var x; var y; var z;\n"+
			"function foo(p1, p2) {}\n"+
			"new foo(new x(), y()); new foo(new x(),y())\n",
		"var x; var y; var z;\n"+
			"function foo() {var p1=new x(); var p2=y()}\n"+
			"new foo(); new foo()\n")

	c.Test(
		"var x; var y; var z;\n"+
			"function foo(p1, p2) {}\n"+
			"new foo(x(), y()); new foo(x(),y())\n",
		"var x; var y; var z;\n"+
			"function foo() {var p1=x(); var p2=y()}\n"+
			"new foo(); new foo()\n")

	c.TestSame(
		"var x; var y; var z;\n" +
			"function foo(p1, p2, p3) {}\n" +
			"new foo(x(), y(), z()); new foo(x(),y(),3)\n")

	c.TestSame(
		"var x; var y; var z;\n" +
			"function foo(p1, p2, p3) {}\n" +
			"new foo(x, y(), z()); new foo(x,y(),3)\n")

	c.Test(
		"var x; var y; var z;\n"+
			"function foo(p1, p2, p3) {}\n"+
			"new foo([], y(), z()); new foo([],y(),3)\n",
		"var x; var y; var z;\n"+
			"function foo(p2, p3) {var p1=[]}\n"+
			"new foo(y(), z()); new foo(y(),3)\n")
}

// Java: OptimizeParametersTest.testMutableValuesDoNotMoveSuper()
func (c *OptimizeParametersTest) TestMutableValuesDoNotMoveSuper() {
	c.TestSame(
		"var A;\n" +
			"function fn(p1) {}\n" +
			"class B extends A { constructor() { fn(super.x); } }\n")
}

// Java: OptimizeParametersTest.testShadows()
func (c *OptimizeParametersTest) TestShadows() {
	c.TestSame(
		"function foo(a) {}\n" +
			"var x;\n" +
			"function f() {\n" +
			"  var x;\n" +
			"  function g() {\n" +
			"    foo(x());\n" +
			"  }\n" +
			"};\n" +
			"foo(x())\n")
}

// Java: OptimizeParametersTest.testNoCrash()
func (c *OptimizeParametersTest) TestNoCrash() {
	c.Test(
		"function foo(a) {}\n"+
			"foo({o:1});\n"+
			"foo({o:1})\n",
		"function foo() {var a = {o:1}}\n"+
			"foo();\n"+
			"foo()\n")
}

// Java: OptimizeParametersTest.testGlobalCatch()
func (c *OptimizeParametersTest) TestGlobalCatch() {
	c.TestSame("function foo(a) {} try {} catch (e) {foo(e)}")
}

// Java: OptimizeParametersTest.testNamelessParameter1()
func (c *OptimizeParametersTest) TestNamelessParameter1() {
	c.Test(c.Externs("var g;"), c.Srcs("f(g()); function f(){}"), c.Expected("f(); function f(){g()}"))
}

// Java: OptimizeParametersTest.testNamelessParameter2()
func (c *OptimizeParametersTest) TestNamelessParameter2() {
	c.Test(
		c.Externs("var g, h;"),
		c.Srcs("f(g(),h()); function f(){}"),
		c.Expected("f(); function f(){g();h()}"))
}

// Java: OptimizeParametersTest.testRewriteUsedClassConstructor1()
func (c *OptimizeParametersTest) TestRewriteUsedClassConstructor1() {
	c.Test(
		"class C {\n"+
			"  constructor(a) {\n"+
			"    use(a);\n"+
			"  }\n"+
			"}\n"+
			"var c = new C();\n",
		"class C {\n"+
			"  constructor( ) {\n"+
			"    var a; // moved from parameter list\n"+
			"    use(a);\n"+
			"  }\n"+
			"}\n"+
			"var c = new C();\n")
}

// Java: OptimizeParametersTest.testConstructorEscapesThroughThis()
func (c *OptimizeParametersTest) TestConstructorEscapesThroughThis() {
	// Demonstrate b/174875103
	c.Test(
		"class C {\n"+
			"  constructor(a) {\n"+
			"    use(a);\n"+
			"  }\n"+
			"  static create() { return new this(2); };\n"+
			"}\n"+
			"class D extends C {\n"+
			"  constructor(b) { super(1); }\n"+
			"\n"+
			"}\n"+
			"var c = new C(1);\n"+
			"var d = new D(1)\n"+
			"var e = D.create();\n",
		"class C {\n"+
			"  constructor() {\n"+
			"    var a = 1; // <-- bad optimization\n"+
			"    use(a);\n"+
			"  }\n"+
			"  static create() { return new this(2); }; // <-- also a call here\n"+
			"}\n"+
			"class D extends C {\n"+
			"  constructor(b) { super(); }\n"+
			"\n"+
			"}\n"+
			"var c = new C();\n"+
			"var d = new D(1)\n"+
			"var e = D.create();\n")
}

// Java: OptimizeParametersTest.testRewriteUsedES5Constructor1()
func (c *OptimizeParametersTest) TestRewriteUsedES5Constructor1() {
	c.Test(
		"/** @constructor */\n"+
			"function C(a) {\n"+
			"  use(a);\n"+
			"}\n"+
			"var c = new C();\n",
		"/** @constructor */\n"+
			"function C( ) {\n"+
			"  var a; // moved from parameter list\n"+
			"  use(a);\n"+
			"}\n"+
			"var c = new C();\n")
}

// Java: OptimizeParametersTest.testRewriteUsedClassConstructor2()
func (c *OptimizeParametersTest) TestRewriteUsedClassConstructor2() {
	// `constructor` aliases the class constructor
	c.Test(
		"class C {\n"+
			"  constructor(a) {\n"+
			"    use(a);\n"+
			"  }\n"+
			"}\n"+
			"var c = new C();\n"+
			"new c.constructor(1);\n",
		"class C {\n"+
			"  constructor( ) {\n"+
			"    var a; // moved from parameter list\n"+
			"    use(a);\n"+
			"  }\n"+
			"}\n"+
			"var c = new C();\n"+
			"// TODO(bradfordcsmith): This call is now broken, since the parameter is ignored.\n"+
			"//     For now we consider the code size savings worth the risk of breaking this\n"+
			"//     coding pattern that we consider bad practice anyway.\n"+
			"new c.constructor(1);\n")
}

// Java: OptimizeParametersTest.testRewriteUsedES5Constructor2()
func (c *OptimizeParametersTest) TestRewriteUsedES5Constructor2() {
	// `constructor` aliases the class constructor
	c.Test(
		"/** @constructor */\n"+
			" function C(a) {\n"+
			"  use(a);\n"+
			"}\n"+
			"var c = new C();\n"+
			"new c.constructor(1);\n",
		"/** @constructor */\n"+
			" function C( ) {\n"+
			"  var a; // moved from parameter list\n"+
			"  use(a);\n"+
			"}\n"+
			"var c = new C();\n"+
			"// TODO(bradfordcsmith): This call is now broken, since the parameter is ignored.\n"+
			"//     For now we consider the code size savings worth the risk of breaking this\n"+
			"//     coding pattern that we consider bad practice anyway.\n"+
			"new c.constructor(1);\n")
}

// Java: OptimizeParametersTest.testNoRewriteUsedClassConstructor3()
func (c *OptimizeParametersTest) TestNoRewriteUsedClassConstructor3() {
	// `super` aliases the super type constructor
	c.TestSame(
		"class C { constructor(a) { use(a); } }\n" +
			"class D extends C { constructor() { super(1); } }\n" +
			"var d = new D(); new C();\n")
}

// Java: OptimizeParametersTest.testRewriteUsedClassConstructor4()
func (c *OptimizeParametersTest) TestRewriteUsedClassConstructor4() {
	// `new.target` aliases self and subtype constructors
	c.Test(
		"class C {\n"+
			"  constructor() {\n"+
			"    var x = new new.target(1);\n"+
			"  }\n"+
			"}\n"+
			"class D extends C {\n"+
			"  constructor(a) {\n"+
			"    super();\n"+
			"  }\n"+
			"}\n"+
			"var d = new D(); new C();\n",
		"class C {\n"+
			"  constructor() {\n"+
			"// TODO(bradfordcsmith): This call is now broken, since the parameter is ignored.\n"+
			"//     For now we consider the code size savings worth the risk of breaking this\n"+
			"//     coding pattern that we consider bad practice anyway.\n"+
			"    var x = new new.target(1);\n"+
			"  }\n"+
			"}\n"+
			"class D extends C {\n"+
			"  constructor( ) {\n"+
			"    var a; // moved from parameter list\n"+
			"    super();\n"+
			"  }\n"+
			"}\n"+
			"var d = new D(); new C();\n")
}

// Java: OptimizeParametersTest.testNoRewriteUsedClassConstructor5()
func (c *OptimizeParametersTest) TestNoRewriteUsedClassConstructor5() {
	// Static class methods "this" values can alias constructors.
	c.TestSame(
		"class C { constructor(a) { use(a); }; static create(a) { new this(1); } }\n" +
			"var c = new C();\n" +
			"C.create();\n")
}

// Java: OptimizeParametersTest.testRewriteUsedClassConstructorWithClassStaticField()
func (c *OptimizeParametersTest) TestRewriteUsedClassConstructorWithClassStaticField() {
	c.Test(
		"class C {\n"+
			"  static field2 = alert();\n"+
			"  constructor(a) {\n"+
			"    use(a);\n"+
			"  }\n"+
			"}\n"+
			"var c = new C(1);\n",
		"class C {\n"+
			"  static field2 = alert();\n"+
			"  constructor( ) {\n"+
			"    var a = 1; // moved from parameter list\n"+
			"    use(a);\n"+
			"  }\n"+
			"}\n"+
			"var c = new C();\n")

	c.Test(
		"class C {\n"+
			"  static field2 = alert();\n"+
			"  constructor(a) {\n"+
			"    use(a);\n"+
			"  }\n"+
			"}\n"+
			"var c = new C(alert());\n",
		"class C {\n"+
			"  static field2 = alert();\n"+
			"  constructor() {\n"+
			"    var a = alert();\n"+
			"    use(a);\n"+
			"  }\n"+
			"}\n"+
			"var c = new C();\n")
}

// Java: OptimizeParametersTest.testRewriteClassStaticBlock_removeOptional()
func (c *OptimizeParametersTest) TestRewriteClassStaticBlock_removeOptional() {
	c.Test(
		"function foo(a,b=1){\n"+
			"  return a * b;\n"+
			"}\n"+
			"class C {\n"+
			"  static {\n"+
			"    use(foo(1));\n"+
			"    use(foo(2));\n"+
			"  }\n"+
			"}\n",
		"function foo(a){\n"+
			"  var b = 1;\n"+
			"  return a * b;\n"+
			"}\n"+
			"class C {\n"+
			"  static {\n"+
			"    use(foo(1));\n"+
			"    use(foo(2));\n"+
			"  }\n"+
			"}\n")
	// TODO(b/240443227): Function parameters inside class static blocks not optimized
	c.TestSame(
		"class C {\n" +
			"  static {\n" +
			"    function foo(a,b=1){\n" +
			"      return(a * b);\n" +
			"    }\n" +
			"    use(foo(1));\n" +
			"    use(foo(2));\n" +
			"  }\n" +
			"}\n")
}

// Java: OptimizeParametersTest.testRewriteClassStaticBlock_trailingUndefinedLiterals()
func (c *OptimizeParametersTest) TestRewriteClassStaticBlock_trailingUndefinedLiterals() {
	c.Test(
		"function foo(a,b){\n"+
			"  return a;\n"+
			"}\n"+
			"class C {\n"+
			"  static {\n"+
			"    use(foo(1, undefined, 2));\n"+
			"    use(foo(2));\n"+
			"  }\n"+
			"}\n",
		"function foo(a,b){\n"+
			"  return a;\n"+
			"}\n"+
			"class C {\n"+
			"  static {\n"+
			"    use(foo(1));\n"+
			"    use(foo(2));\n"+
			"  }\n"+
			"}\n")
	// TODO(b/240443227): Function parameters inside class static blocks not optimized
	c.TestSame(
		"class C {\n" +
			"  static {\n" +
			"    function foo(a,b){\n" +
			"      return a;\n" +
			"    }\n" +
			"    use(foo(1, undefined, 2));\n" +
			"    use(foo(2));\n" +
			"  }\n" +
			"}\n")
}

// Java: OptimizeParametersTest.testRewriteClassStaticBlock_inlineParameter()
func (c *OptimizeParametersTest) TestRewriteClassStaticBlock_inlineParameter() {
	c.Test(
		"function foo(a){\n"+
			"  return a;\n"+
			"}\n"+
			"class C {\n"+
			"  static {\n"+
			"    use(foo(1));\n"+
			"    use(foo(1));\n"+
			"    use(foo(1));\n"+
			"  }\n"+
			"}\n",
		"function foo(){\n"+
			"  var a = 1;\n"+
			"  return a;\n"+
			"}\n"+
			"class C {\n"+
			"  static {\n"+
			"    use(foo());\n"+
			"    use(foo());\n"+
			"    use(foo());\n"+
			"  }\n"+
			"}\n")
	// TODO(b/240443227): Function parameters inside class static blocks not optimized
	c.TestSame(
		"class C {\n" +
			"  static {\n" +
			"    function foo(a){\n" +
			"      return(a);\n" +
			"    }\n" +
			"    use(foo(1));\n" +
			"    use(foo(1));\n" +
			"    use(foo(1));\n" +
			"  }\n" +
			"}\n")
}

// Java: OptimizeParametersTest.testNoRewriteUsedClassConstructorWithClassNonstaticField_sideEffect()
func (c *OptimizeParametersTest) TestNoRewriteUsedClassConstructorWithClassNonstaticField_sideEffect() {
	c.TestSame(
		"class C {\n" +
			"  field2 = alert();\n" +
			"  constructor(a) {\n" +
			"    use(a);\n" +
			"  }\n" +
			"}\n" +
			"var c = new C(1);\n")

	c.TestSame(
		"class C {\n" +
			"  field2 = alert();\n" +
			"  constructor(a) {\n" +
			"    use(a);\n" +
			"  }\n" +
			"}\n" +
			"var c = new C(alert());\n")
}

// Java: OptimizeParametersTest.testNoRewriteUsedClassConstructorWithClassNonstaticComputedField_sideEffect()
func (c *OptimizeParametersTest) TestNoRewriteUsedClassConstructorWithClassNonstaticComputedField_sideEffect() {
	c.TestSame(
		"class C {\n" +
			"  ['field2'] = alert(2);\n" +
			"  constructor(a) {\n" +
			"    use(a);\n" +
			"  }\n" +
			"}\n" +
			"var c = new C(alert(1)); // alerts in order 1, 2\n")
}

// Java: OptimizeParametersTest.testNoRewriteUsedClassMethodParam1()
func (c *OptimizeParametersTest) TestNoRewriteUsedClassMethodParam1() {
	c.TestSame(
		"class C { method(a) { use(a); } }\n" +
			"var c = new C(); c.method(1); c.method(2)\n")
}

// Java: OptimizeParametersTest.testNoRewriteUnusedClassComputedMethodParam1()
func (c *OptimizeParametersTest) TestNoRewriteUnusedClassComputedMethodParam1() {
	c.TestSame(
		"class C { [method](a) { } }\n" +
			"var c = new C(); c[method](1); c[method](2)\n")
}

// Java: OptimizeParametersTest.testRewriteUsedClassMethodParam1()
func (c *OptimizeParametersTest) TestRewriteUsedClassMethodParam1() {
	c.Test(
		"class C { method(a) {          }} new C().method(1)",
		"class C { method( ) {var a = 1;}} new C().method( )")
}

// Java: OptimizeParametersTest.testNoRewriteUnsedObjectMethodParam()
func (c *OptimizeParametersTest) TestNoRewriteUnsedObjectMethodParam() {
	c.TestSame("var o = { method(a) {          }}; o.method(1)")
}

// Java: OptimizeParametersTest.testNoRewriteDestructured1()
func (c *OptimizeParametersTest) TestNoRewriteDestructured1() {
	c.TestSame(
		"class C { m(a) {}};\n" +
			"var c = new C();\n" +
			"({m} = c);\n" +
			"c.m(1)\n")
	c.TestSame(
		"class C { m(a) {}};\n" +
			"var c = new C();\n" +
			"({m:x} = c);\n" +
			"c.m(1)\n")
	c.TestSame(
		"class C { m(a) {}};\n" +
			"var c = new C();\n" +
			"({xx:C.prototype.m} = {xx:function(a) {}});\n" +
			"c.m(1)\n")
}

// Java: OptimizeParametersTest.testNoRewriteDestructured2()
func (c *OptimizeParametersTest) TestNoRewriteDestructured2() {
	c.TestSame(
		"var x = function(a) {};\n" +
			"({x} = {})\n" +
			"x(1)\n")
	c.TestSame(
		"var x = function(a) {};\n" +
			"({x:x} = {})\n" +
			"x(1)\n")
	c.TestSame(
		"var x = function(a) {};\n" +
			"({x:x = function(a) {}} = {})\n" +
			"x(1)\n")
}

// Java: OptimizeParametersTest.testNoRewriteDestructured3()
func (c *OptimizeParametersTest) TestNoRewriteDestructured3() {
	c.TestSame(
		"var x = function(a) {};\n" +
			"[x = function() {}] = []\n" +
			"x(1)\n")
	c.TestSame(
		"class C { method() { return 1 }}\n" +
			"var c = new C();\n" +
			"var y = C.prototype;\n" +
			"[y.method] = []\n" +
			"c.method()\n")
	c.TestSame(
		"class C { method() { return 1 }}\n" +
			"[x.method] = []\n" +
			"y.method()\n")
}

// Java: OptimizeParametersTest.testNoRewriteTagged1()
func (c *OptimizeParametersTest) TestNoRewriteTagged1() {
	// Optimizing methods called though tagged template literal requires
	// specific knowledge of how tagged templated is supplied to the method.
	c.TestSame(
		"var f = function(a, b, c) {};\n" +
			"f`tagged`\n")

	c.TestSame(
		"var f = function(a, b, c) {};\n" +
			"f`tagged`\n" +
			"f()\n")
}

// Java: OptimizeParametersTest.testArrow()
func (c *OptimizeParametersTest) TestArrow() {
	// Optimizing methods called though tagged template literal requires
	// specific knowledge of how tagged templated is supplied to the method.
	c.TestSame(
		"var f = (a)=>{};\n" +
			"f`tagged`\n")

	c.TestSame(
		"var f = function(a, b, c) {};\n" +
			"f`tagged`\n" +
			"f()\n")
}

// Java: OptimizeParametersTest.testNoRewriteTagged_methodCall()
func (c *OptimizeParametersTest) TestNoRewriteTagged_methodCall() {
	c.TestSame(
		"var obj = {\n" +
			"  f: function(a, b, c) {}\n" +
			"};\n" +
			"obj.f`tagged`;\n")
}

// Java: OptimizeParametersTest.testSuperInvocation_preventsParamInlining_whenImplicit()
func (c *OptimizeParametersTest) TestSuperInvocation_preventsParamInlining_whenImplicit() {
	c.TestSame(
		"class Foo {\n" +
			"  constructor(x) {\n" +
			"    this.x = x;\n" +
			"  }\n" +
			"}\n" +
			"\n" +
			"// lack of explicit constructor prevents optimizing calls to both classes.\n" +
			"class Bar extends Foo { }\n" +
			"\n" +
			"new Foo(4);\n" +
			"new Bar(4);\n")
}

// Java: OptimizeParametersTest.testSuperInvocationCanBeInlined()
func (c *OptimizeParametersTest) TestSuperInvocationCanBeInlined() {
	c.Test(
		"class Foo {\n"+
			"  constructor(x) {\n"+
			"    this.x = x;\n"+
			"  }\n"+
			"}\n"+
			"\n"+
			"class Bar extends Foo {\n"+
			"  constructor() {\n"+
			"    super(4);\n"+
			"  }\n"+
			"}\n"+
			"\n"+
			"new Foo(4);\n"+
			"new Bar();\n",
		"class Foo {\n"+
			"  constructor( ) {\n"+
			"    var x = 4; // moved from the parameter list\n"+
			"    this.x = x;\n"+
			"  }\n"+
			"}\n"+
			"\n"+
			"class Bar extends Foo {\n"+
			"  constructor() {\n"+
			"    super( );\n"+
			"  }\n"+
			"}\n"+
			"\n"+
			"new Foo( );\n"+
			"new Bar();\n")
}

// Java: OptimizeParametersTest.testSuperInvocationCannotBeInlinedWhenExtendingExpression()
func (c *OptimizeParametersTest) TestSuperInvocationCannotBeInlinedWhenExtendingExpression() {
	c.TestSame(
		"class Foo {\n" +
			"  constructor(x) {\n" +
			"    this.x = x;\n" +
			"  }\n" +
			"}\n" +
			"\n" +
			"// This bit of indirection prevents us from recognizing what\n" +
			"// is being extended, so `new Foo()` cannot be optimized.\n" +
			"class Bar extends (() => Foo)() {\n" +
			"  constructor() {\n" +
			"    super(4);\n" +
			"  }\n" +
			"}\n" +
			"\n" +
			"new Foo(4);\n" +
			"new Bar();\n")
}

// Java: OptimizeParametersTest.testRemoveOptionalDestructuringParam()
func (c *OptimizeParametersTest) TestRemoveOptionalDestructuringParam() {
	c.Test("function f({x}) {} f();", "function f() { var x; ({x} = void 0); } f();")
	c.Test("function f({x} = {}) {} f();", "function f() { var x; ({x} = {}); } f();")
}

// Java: OptimizeParametersTest.testClassField()
func (c *OptimizeParametersTest) TestClassField() {
	c.TestSame(
		"class C {\n" +
			"  x = function(a,b) { return a + b };\n" +
			"}\n" +
			"new C().x(1,2)\n")
}

// Java: OptimizeParametersTest.testTrailingUndefinedLiterals()
func (c *OptimizeParametersTest) TestTrailingUndefinedLiterals() {
	c.Test(
		"function foo(a) { use(a);}; foo(undefined); foo(2);",
		"function foo(a) { use(a);}; foo(         ); foo(2);")
}

// Java: OptimizeParametersTest.testTrailingUndefinedLiterals_multiple()
func (c *OptimizeParametersTest) TestTrailingUndefinedLiterals_multiple() {
	c.Test(
		"function foo(a, b, c) { use(a, b, c); }\n"+
			"foo(undefined);\n"+
			"foo(undefined, void 0);\n"+
			"foo(undefined, void 0, undefined);\n"+
			"foo(2);\n",
		"function foo(a, b, c) { use(a, b, c); }\n"+
			"foo();\n"+
			"foo();\n"+
			"foo();\n"+
			"foo(2);\n")
}

// Java: OptimizeParametersTest.testTrailingUndefinedLiterals_functionRefsArguments()
func (c *OptimizeParametersTest) TestTrailingUndefinedLiterals_functionRefsArguments() {
	c.TestSame("function foo(a) { use(arguments);}; foo(undefined); foo(2);")
}

// Java: OptimizeParametersTest.testTrailingUndefinedLiterals_afterASpread()
func (c *OptimizeParametersTest) TestTrailingUndefinedLiterals_afterASpread() {
	c.TestSame("function foo(a,b) { use(a)}; foo(...[1], undefined, undefined);")
	c.TestSame("function foo(a,b) { use(a)}; foo(undefined, ...[1], undefined); foo(2);")
}

// Java: OptimizeParametersTest.testTrailingUndefinedLiterals_afterAllFormalParameters()
func (c *OptimizeParametersTest) TestTrailingUndefinedLiterals_afterAllFormalParameters() {
	c.Test(
		"function foo(a, b) { use(a)}; foo('used', undefined, undefined, 2, 'a'); foo(2);",
		"function foo(a, b) { use(a)}; foo('used');                               foo(2);")
}

// Java: OptimizeParametersTest.testTrailingUndefinedLiterals_afterAllFormalParameters_sideEffects()
func (c *OptimizeParametersTest) TestTrailingUndefinedLiterals_afterAllFormalParameters_sideEffects() {
	c.TestSame("function foo(a, b) { use(a)}; foo('used', undefined, sideEffects()); foo(2);")
}

// Java: OptimizeParametersTest.testInliningSideEffectfulArg_updatesInvocationSideEffects()
func (c *OptimizeParametersTest) TestInliningSideEffectfulArg_updatesInvocationSideEffects() {
	c.EnableComputeSideEffects()
	c.Test(
		c.Externs("function sideEffects() {}"),
		c.Srcs("function foo() {} foo(sideEffects()); foo(sideEffects());"),
		c.Expected("function foo() { sideEffects(); } foo(); foo();"))

	// Inspect the AST to verify both calls to `foo()` are marked as mutating global state
	jsRoot := c.GetLastCompiler().GetJsRoot()
	calls :=
		CodeSubTree_FindNodesNonEmpty(
			jsRoot, func(n rhino.Node) bool { return n.IsCall() && n.GetFirstChild().MatchesName("foo") })
	for _, call := range calls {
		AssertThat(call.GetSideEffectFlags()).
			IsEqualTo(
				rhino.NewNode_SideEffectFlagsWithValue(rhino.Node_SideEffectFlags_NO_SIDE_EFFECTS).
					SetMutatesGlobalState().
					ValueOf())
	}
}

// Java: OptimizeParametersTest.testNoHoistNewTarget()
func (c *OptimizeParametersTest) TestNoHoistNewTarget() {
	c.TestSame(
		"function check(x) { if (x) throw 1; }\n" +
			"class C { constructor() { check(!new.target); } }\n" +
			"new C();\n")
}

// Java: OptimizeParametersTest.testNoHoistImportMeta()
func (c *OptimizeParametersTest) TestNoHoistImportMeta() {
	c.TestSame(
		"function check(x) { if (x) throw 1; }\n" +
			"function f() { check(import.meta.url); }\n" +
			"f();\n")
	c.TestSame(
		"function check(x) { if (x) throw 1; }\n" +
			"function f() { check(import.meta); }\n" +
			"f();\n")
}

// Java: OptimizeParametersTest.testPublicClassFields_propertyReadInInitializerCanBeSideEffected()
func (c *OptimizeParametersTest) TestPublicClassFields_propertyReadInInitializerCanBeSideEffected() {
	// TODO(b/572053067): side-effectful constructor argument moved into constructor body after
	// non-static field initializer that reads modified state
	c.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2022)
	c.Test(
		"function use(val) {}\n"+
			"const obj = {p: 1};\n"+
			"class C {\n"+
			"  x = obj.p;\n"+
			"  constructor(a) {\n"+
			"    use(this.x + a);\n"+
			"  }\n"+
			"}\n"+
			"new C(obj.p = 2);\n",
		"function use(val) {}\n"+
			"const obj = {p: 1};\n"+
			"class C {\n"+
			"  x = obj.p;\n"+
			"  constructor() {\n"+
			"    var a = obj.p = 2; // <----- this is wrong, obj.p = 2 *should* execute before x = obj.p\n"+
			"    use(this.x + a);\n"+
			"  }\n"+
			"}\n"+
			"new C();\n")
}

// Java: OptimizeParametersTest.testAsyncGenerator_doNotMoveSideEffectfulCallArgIntoBody()
func (c *OptimizeParametersTest) TestAsyncGenerator_doNotMoveSideEffectfulCallArgIntoBody() {
	// TODO(b/573183650): side-effectful call argument moved into async generator body
	c.Test(
		"async function* g(p) {\n"+
			"  yield await p;\n"+
			"}\n"+
			"g(alert(1));\n",
		"async function* g() {\n"+
			"  var p = alert(1);\n"+
			"  yield await p;\n"+
			"}\n"+
			"g();\n")
}

// Java: OptimizeParametersTest.testAsyncGeneratorMethod_doNotMoveSideEffectfulCallArgIntoBody()
func (c *OptimizeParametersTest) TestAsyncGeneratorMethod_doNotMoveSideEffectfulCallArgIntoBody() {
	// TODO(b/573183650): side-effectful call argument moved into async generator method body
	c.Test(
		"class C {\n"+
			"  async *g(p) {\n"+
			"    yield await p;\n"+
			"  }\n"+
			"}\n"+
			"new C().g(alert(1));\n",
		"class C {\n"+
			"  async *g() {\n"+
			"    var p = alert(1);\n"+
			"    yield await p;\n"+
			"  }\n"+
			"}\n"+
			"(new C()).g();\n")
}

// Java: OptimizeParametersTest.testAsyncGenerator_doNotMoveMutableCallArgAcrossAwait()
func (c *OptimizeParametersTest) TestAsyncGenerator_doNotMoveMutableCallArgAcrossAwait() {
	// TODO(b/573183650): mutable call argument moved into async generator body across suspension
	// point
	c.Test(
		"var x = 1;\n"+
			"async function* g(a) {\n"+
			"  yield a;\n"+
			"}\n"+
			"async function run() {\n"+
			"  const it = g(x);\n"+
			"  await alert();\n"+
			"  x = 2;\n"+
			"  return (await it.next()).value;\n"+
			"}\n"+
			"run();\n",
		"var x = 1;\n"+
			"async function* g() {\n"+
			"  var a = x;\n"+
			"  yield a;\n"+
			"}\n"+
			"async function run() {\n"+
			"  const it = g();\n"+
			"  await alert();\n"+
			"  x = 2;\n"+
			"  return (await it.next()).value;\n"+
			"}\n"+
			"run();\n")
}

// Java: OptimizeParametersTest.testAsyncGenerator_doNotMoveDefaultOrDestructuringParamIntoBody()
func (c *OptimizeParametersTest) TestAsyncGenerator_doNotMoveDefaultOrDestructuringParamIntoBody() {
	// TODO(b/573183650): default or destructuring parameter moved into async generator body
	c.Test(
		"async function* g(a = alert(1)) {\n"+
			"  yield a;\n"+
			"}\n"+
			"g();\n",
		"async function* g() {\n"+
			"  var a = alert(1);\n"+
			"  yield a;\n"+
			"}\n"+
			"g();\n")
	c.Test(
		"async function* g({a} = {a: alert(1)}) {\n"+
			"  yield a;\n"+
			"}\n"+
			"g();\n",
		"async function* g() {\n"+
			"  ({a} = {a: alert(1)});\n"+
			"  var a;\n"+
			"  yield a;\n"+
			"}\n"+
			"g();\n")
	c.Test(
		"async function* g({a}) {\n"+
			"  yield a;\n"+
			"}\n"+
			"g();\n",
		"async function* g() {\n"+
			"  ({a} = void 0);\n"+
			"  var a;\n"+
			"  yield a;\n"+
			"}\n"+
			"g();\n")
}

// Java: OptimizeParametersTest.testGenerator_doNotMoveSideEffectfulCallArgIntoBody()
func (c *OptimizeParametersTest) TestGenerator_doNotMoveSideEffectfulCallArgIntoBody() {
	// TODO(b/573183650): side-effectful call argument moved into generator body
	c.Test(
		"function* g(p) {\n"+
			"  yield p;\n"+
			"}\n"+
			"g(alert(1));\n",
		"function* g() {\n"+
			"  var p = alert(1);\n"+
			"  yield p;\n"+
			"}\n"+
			"g();\n")
}

// Java: OptimizeParametersTest.testGeneratorMethod_doNotMoveSideEffectfulCallArgIntoBody()
func (c *OptimizeParametersTest) TestGeneratorMethod_doNotMoveSideEffectfulCallArgIntoBody() {
	// TODO(b/573183650): side-effectful call argument moved into generator method body
	c.Test(
		"class C {\n"+
			"  *g(p) {\n"+
			"    yield p;\n"+
			"  }\n"+
			"}\n"+
			"new C().g(alert(1));\n",
		"class C {\n"+
			"  *g() {\n"+
			"    var p = alert(1);\n"+
			"    yield p;\n"+
			"  }\n"+
			"}\n"+
			"(new C()).g();\n")
}

// Java: OptimizeParametersTest.testGenerator_doNotMoveMutableCallArgAcrossNext()
func (c *OptimizeParametersTest) TestGenerator_doNotMoveMutableCallArgAcrossNext() {
	// TODO(b/573183650): mutable call argument moved into generator body across suspension point
	c.Test(
		"var x = 1;\n"+
			"function* g(a) {\n"+
			"  yield a;\n"+
			"}\n"+
			"function run() {\n"+
			"  const it = g(x);\n"+
			"  x = 2;\n"+
			"  return it.next().value;\n"+
			"}\n"+
			"run();\n",
		"var x = 1;\n"+
			"function* g() {\n"+
			"  var a = x;\n"+
			"  yield a;\n"+
			"}\n"+
			"function run() {\n"+
			"  const it = g();\n"+
			"  x = 2;\n"+
			"  return it.next().value;\n"+
			"}\n"+
			"run();\n")
}

// Java: OptimizeParametersTest.testGenerator_doNotMoveDefaultOrDestructuringParamIntoBody()
func (c *OptimizeParametersTest) TestGenerator_doNotMoveDefaultOrDestructuringParamIntoBody() {
	// TODO(b/573183650): default or destructuring parameter moved into generator body
	c.Test(
		"function* g(a = alert(1)) {\n"+
			"  yield a;\n"+
			"}\n"+
			"g();\n",
		"function* g() {\n"+
			"  var a = alert(1);\n"+
			"  yield a;\n"+
			"}\n"+
			"g();\n")
	c.Test(
		"function* g({a} = {a: alert(1)}) {\n"+
			"  yield a;\n"+
			"}\n"+
			"g();\n",
		"function* g() {\n"+
			"  ({a} = {a: alert(1)});\n"+
			"  var a;\n"+
			"  yield a;\n"+
			"}\n"+
			"g();\n")
	c.Test(
		"function* g({a}) {\n"+
			"  yield a;\n"+
			"}\n"+
			"g();\n",
		"function* g() {\n"+
			"  ({a} = void 0);\n"+
			"  var a;\n"+
			"  yield a;\n"+
			"}\n"+
			"g();\n")
}
