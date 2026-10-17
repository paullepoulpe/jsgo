// Port of com/google/javascript/jscomp/RemoveUnusedCodeNameAnalyzerTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/RemoveUnusedCodeNameAnalyzerTest.java)

/*
 * Copyright 2006 The Closure Compiler Authors.
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

type RemoveUnusedCodeNameAnalyzerTest struct {
	CompilerTestCaseBase
}

// Java: RemoveUnusedCodeNameAnalyzerTest()
func NewRemoveUnusedCodeNameAnalyzerTest() *RemoveUnusedCodeNameAnalyzerTest {
	c := &RemoveUnusedCodeNameAnalyzerTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, RemoveUnusedCodeNameAnalyzerTest_EXTERNS)
	return c
}

func TestRemoveUnusedCodeNameAnalyzerTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewRemoveUnusedCodeNameAnalyzerTest() })
}

const RemoveUnusedCodeNameAnalyzerTest_EXTERNS = "/**\n" +
	" * @constructor\n" +
	" * @param {*=} opt_value\n" +
	" * @return {!Object}\n" +
	" */\n" +
	"function Object(opt_value) {}\n" +
	"/** @type {Function} */\n" +
	"Object.defineProperties = function() {};\n" +
	"/**\n" +
	" * @constructor\n" +
	" * @param {string} message\n" +
	" */\n" +
	"function Error(message) {}\n" +
	"var window, top, console;\n" +
	"var document;\n" +
	"var Function;\n" +
	"var Array;\n" +
	"var goog = {};\n" +
	"goog.inherits = function(childClass, parentClass) {};\n" +
	"goog.mixin = function(target, base) {};\n" +
	"function goog$addSingletonGetter(a){}\n" +
	"var externfoo = {};\n" +
	"externfoo.externProp1 = 1;\n" +
	"externfoo.externProp2 = 2;\n" +
	"function doThing1() {}\n" +
	"function doThing2() {}\n" +
	"function use(something) {}\n" +
	"function alert(something) {}\n" +
	"function sideEffect() {}\n"

// Java: RemoveUnusedCodeNameAnalyzerTest.getProcessor(Compiler)
func (c *RemoveUnusedCodeNameAnalyzerTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewRemoveUnusedCodeNameAnalyzerTest_MarkNoSideEffectCallsAndRemoveUnusedCodeRunner(compiler)
}

type RemoveUnusedCodeNameAnalyzerTest_MarkNoSideEffectCallsAndRemoveUnusedCodeRunner struct {
	pureFunctionIdentifier *PureFunctionIdentifier_Driver
	removeUnusedCode       *RemoveUnusedCode
}

// Java: RemoveUnusedCodeNameAnalyzerTest.MarkNoSideEffectCallsAndRemoveUnusedCodeRunner(Compiler)
func NewRemoveUnusedCodeNameAnalyzerTest_MarkNoSideEffectCallsAndRemoveUnusedCodeRunner(compiler *Compiler) *RemoveUnusedCodeNameAnalyzerTest_MarkNoSideEffectCallsAndRemoveUnusedCodeRunner {
	r := &RemoveUnusedCodeNameAnalyzerTest_MarkNoSideEffectCallsAndRemoveUnusedCodeRunner{}
	r.pureFunctionIdentifier = NewPureFunctionIdentifier_Driver(compiler)
	r.removeUnusedCode =
		NewRemoveUnusedCode_Builder(compiler).
			RemoveGlobals(true).
			RemoveLocalVars(true).
			RemoveUnusedPrototypeProperties(true).
			RemoveUnusedThisProperties(true).
			RemoveUnusedObjectDefinePropertiesDefinitions(true).
			// Removal of function expression names isn't what these tests are about.
			// It just adds noise to the tests when we can't use testSame() because of it.
			PreserveFunctionExpressionNames(true).
			Build()
	return r
}

// Java: RemoveUnusedCodeNameAnalyzerTest.MarkNoSideEffectCallsAndRemoveUnusedCodeRunner.process(Node, Node)
func (r *RemoveUnusedCodeNameAnalyzerTest_MarkNoSideEffectCallsAndRemoveUnusedCodeRunner) Process(externs *rhino.Node, root *rhino.Node) {
	r.pureFunctionIdentifier.Process(externs, root)
	r.removeUnusedCode.Process(externs, root)
}

// Java: RemoveUnusedCodeNameAnalyzerTest.setUp()
func (c *RemoveUnusedCodeNameAnalyzerTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	// Allow testing of features that aren't fully supported for output yet.
	c.EnableNormalize()
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.EnableGatherExternProperties()
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testDefaultingAssignmentWithAssignedProperty()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestDefaultingAssignmentWithAssignedProperty() {
	c.Test("var x = function() {}; x.externProp1 = function() {}; var x = x || {};", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveVarDeclaration1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveVarDeclaration1() {
	c.Test("var foo = 3;", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveVarDeclaration2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveVarDeclaration2() {
	c.Test("var foo = 3, bar = 4; externfoo = foo;", "var foo = 3; externfoo = foo;")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveVarDeclaration3()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveVarDeclaration3() {
	// TODO(b/66971163): Don't block removal when the reference to a name is not actually used.
	// e.g. `b; c` below.
	c.Test(
		"var a = doThing1(),     b = 1, c = 2; b; c", // preserve newline
		"        doThing1(); var b = 1, c = 2; b; c")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveVarDeclaration4()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveVarDeclaration4() {
	c.Test(
		"var a = 0, b = doThing1(),     c = 2; a; c", // preserve newline
		"var a = 0;     doThing1(); var c = 2; a; c")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveVarDeclaration5()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveVarDeclaration5() {
	c.Test(
		"var a = 0, b = 1, c = doThing1(); a; b", // preserve newline
		"var a = 0, b = 1;     doThing1(); a; b")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveVarDeclaration6()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveVarDeclaration6() {
	c.Test("var a = 0, b = a = 1; a", "var a = 0; a = 1; a")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveVarDeclaration7()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveVarDeclaration7() {
	c.Test("var a = 0, b = a = 1", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveVarDeclaration8()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveVarDeclaration8() {
	c.Test("var a;var b = 0, c = a = b = 1", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveLetDeclaration1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveLetDeclaration1() {
	c.Test("let foo = 3;", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveLetDeclaration2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveLetDeclaration2() {
	c.Test("let foo = 3, bar = 4; externfoo = foo;", "let foo = 3; externfoo = foo;")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveLetDeclaration3()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveLetDeclaration3() {
	c.Test(
		"let a = doThing1(),     b = 1, c = 2; b; c", // preserve newline
		"        doThing1(); let b = 1, c = 2; b; c")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveLetDeclaration4()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveLetDeclaration4() {
	c.Test(
		"let a = 0, b = doThing1(),     c = 2; a; c", // preserve newline
		"let a = 0;     doThing1(); let c = 2; a; c")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveLetDeclaration5()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveLetDeclaration5() {
	c.Test(
		"let a = 0, b = 1, c = doThing1(); a; b", // preserve newline
		"let a = 0, b = 1;     doThing1(); a; b")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveLetDeclaration6()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveLetDeclaration6() {
	c.Test("let a = 0, b = a = 1; a", "let a = 0; a = 1; a")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveLetDeclaration7()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveLetDeclaration7() {
	c.Test("let a = 0, b = a = 1", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveLetDeclaration8()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveLetDeclaration8() {
	c.Test("let a;let b = 0, c = a = b = 1", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveLetDeclaration9()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveLetDeclaration9() {
	// The variable inside the block doesn't get removed (but does get renamed by Normalize).
	c.Test("let x = 1; if (true) { let x = 2; x; }", "if (true) { let x$jscomp$1 = 2; x$jscomp$1; }")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveLetInBlock1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveLetInBlock1() {
	c.TestSame(
		"if (true) { // preserve newline\n" +
			"  let x = 1; alert(x);\n" +
			"}\n")

	c.Test(
		"if (true) { let x = 1; }", // preserve newline
		"if (true)              ;")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveLetInBlock2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveLetInBlock2() {
	c.TestSame(
		"if (true) { // preserve newline\n" +
			"  let x = 1; alert(x);\n" +
			"} else {\n" +
			"  let x = 1; alert(x);\n" +
			"}\n")

	c.Test(
		c.Srcs(
			"if (true) { // preserve newline\n"+
				"  let x = 1;\n"+
				"} else {\n"+
				"  let x = 1;\n"+
				"}\n"),
		c.Expected("if (true); else;"))
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveConstDeclaration1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveConstDeclaration1() {
	c.Test("const a = 4;", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveConstDeclaration2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveConstDeclaration2() {
	c.TestSame("const a = 4; window.x = a;")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveConstDeclaration3()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveConstDeclaration3() {
	// The variable inside the block doesn't get removed (but does get renamed by Normalize).
	c.Test(
		"const x = 1; if (true) { const x = 2; x; }",
		"if (true) { const x$jscomp$1 = 2; x$jscomp$1; }")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveDeclaration1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveDeclaration1() {
	c.Test("var a;var b = 0, c = a = b = 1", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveDeclaration2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveDeclaration2() {
	c.Test("var a,b,c; c = a = b = 1", "")
}

// TODO(b/66971163): enable or remove
// Java: RemoveUnusedCodeNameAnalyzerTest.disabledTestRemoveDeclaration3()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestDisabledTestRemoveDeclaration3() {
	c.T().Skip("@Ignore")
	c.Test("var a,b,c; c = a = b = {}; a.x = 1;", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveDeclaration4()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveDeclaration4() {
	c.Test(
		"var a,b,c; c = a = b = {}; a.x = 1;alert(c.x);",
		"var a,  c; c = a     = {}; a.x = 1;alert(c.x);")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveDeclaration5()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveDeclaration5() {
	c.Test("var a,b,c; c = a = b = null; use(b)", "var b;b=null;use(b)")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveDeclaration6()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveDeclaration6() {
	c.Test("var a,b,c; c = a = b = 'str';use(b)", "var b;b='str';use(b)")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveDeclaration7()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveDeclaration7() {
	c.Test("var a,b,c; c = a = b = true;use(b)", "var b;b=true;use(b)")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveFunction1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveFunction1() {
	c.Test("var foo = function(){};", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveFunction2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveFunction2() {
	c.Test("var foo; foo = function(){};", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveFunction3()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveFunction3() {
	c.Test("var foo = {}; foo.bar = function() {};", "")
}

// TODO(b/66971163): enable or remove this test
// Java: RemoveUnusedCodeNameAnalyzerTest.disabledTestRemoveFunction4()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestDisabledTestRemoveFunction4() {
	c.T().Skip("@Ignore")
	c.Test("var a = {}; a.b = {}; a.b.c = function() {};", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testDontRemoveFunctionOnNamespaceThatEscapes()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestDontRemoveFunctionOnNamespaceThatEscapes() {
	c.TestSame("var a = doThing1(); a.b = {}; a.b.c = function() {};")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testReferredToByWindow()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestReferredToByWindow() {
	c.TestSame("var foo = {}; foo.bar = function() {}; window['fooz'] = foo.bar")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testExtern()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestExtern() {
	c.TestSame("externfoo = 5")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveNamedFunction()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveNamedFunction() {
	c.Test("function foo(){}", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveRecursiveFunction1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveRecursiveFunction1() {
	c.Test("function f(){f()}", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveRecursiveFunction2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveRecursiveFunction2() {
	c.Test("var f = function (){f()}", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveRecursiveFunction2a()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveRecursiveFunction2a() {
	c.Test("var f = function g(){g()}", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveRecursiveFunction3()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveRecursiveFunction3() {
	c.Test("var f;f = function (){f()}", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveRecursiveFunction4()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveRecursiveFunction4() {
	// don't remove redefinition of an external variable
	c.TestSame("doThing1 = function (){doThing1()}")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveRecursiveFunction5()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveRecursiveFunction5() {
	c.Test("function g(){f()}function f(){g()}", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveRecursiveFunction6()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveRecursiveFunction6() {
	c.Test("var f=function(){g()};function g(){f()}", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveRecursiveFunction7()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveRecursiveFunction7() {
	c.Test("var g = function(){f()};var f = function(){g()}", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveRecursiveFunction8()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveRecursiveFunction8() {
	c.Test("var o = {};o.f = function(){o.f()}", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveRecursiveFunction9()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveRecursiveFunction9() {
	c.TestSame("var o = {};o.f = function(){o.f()};o.f()")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSideEffectClassification1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSideEffectClassification1() {
	c.TestSame("doThing1();")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSideEffectClassification2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSideEffectClassification2() {
	c.Test("var a = doThing1();", "doThing1();")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSideEffectClassification3()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSideEffectClassification3() {
	c.TestSame("var a = doThing1();window['b']=a;")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSideEffectClassification4()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSideEffectClassification4() {
	c.TestSame("function sef(){} sef();")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSideEffectClassification5()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSideEffectClassification5() {
	c.TestSame("function nsef(){} var a = nsef();window['b']=a;")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSideEffectClassification6()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSideEffectClassification6() {
	c.TestSame("function sef(){} sef();")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSideEffectClassification7()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSideEffectClassification7() {
	c.TestSame("function sef(){} var a = sef();window['b']=a;")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNoSideEffectAnnotation1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNoSideEffectAnnotation1() {
	c.Test(c.Externs("function f(){}"), c.Srcs("var a = f();"), c.Expected("f()"))

	c.Test(c.Externs("function f(){}"), c.Srcs("let a = f();"), c.Expected("f()"))

	c.Test(c.Externs("function f(){}"), c.Srcs("const a = f();"), c.Expected("f()"))
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNoSideEffectAnnotation2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNoSideEffectAnnotation2() {
	c.Test(c.Externs("/** @nosideeffects */ function f(){}"), c.Srcs("var a = f();"), c.Expected(""))
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNoSideEffectAnnotation3()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNoSideEffectAnnotation3() {
	c.Test(c.Externs("var f = function(){};"), c.Srcs("var a = f();"), c.Expected("f()"))
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNoSideEffectAnnotation4()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNoSideEffectAnnotation4() {
	c.Test(
		c.Externs("var f = /** @nosideeffects */ function(){};"), c.Srcs("var a = f();"), c.Expected(""))

	c.Test(
		c.Externs("/** @nosideeffects */ var f = function(){};"), c.Srcs("var a = f();"), c.Expected(""))
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNoSideEffectAnnotation5()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNoSideEffectAnnotation5() {
	c.Test(
		"var f; f = function(){alert('a')}; var a = f();",
		"var f; f = function(){alert('a')}; f();")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNoSideEffectAnnotation6()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNoSideEffectAnnotation6() {
	c.Test(
		c.Externs("var f = /** @nosideeffects */ function(){};"), //
		c.Srcs("var a = f();"),
		c.Expected(""))
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNoSideEffectAnnotation7()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNoSideEffectAnnotation7() {
	c.Test(
		c.Externs("var f = /** @nosideeffects */ function(){};"),
		c.Srcs("f = function(){};var a = f();"),
		c.Expected("f = function(){};"))

	c.Test(
		c.Externs("function sideEffect() {}; var f = /** @nosideeffects */ function(){};"),
		c.Srcs("f = function(){sideEffect()};var a = f();"),
		c.Expected("f = function(){sideEffect()};f()"))
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNoSideEffectAnnotation8()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNoSideEffectAnnotation8() {
	c.Test(
		c.Externs("var f = function(){}; f = /** @nosideeffects */ function(){};"),
		c.Srcs("var a = f();"),
		c.Expected("f();"))
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNoSideEffectAnnotation_whenUsedOnDuplicateDefinitions_eliminatesSideEffects()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNoSideEffectAnnotation_whenUsedOnDuplicateDefinitions_eliminatesSideEffects() {
	c.Test(
		c.Externs(
			"var f = /** @nosideeffects */ function(){};",
			"var f = /** @nosideeffects */ function(){};"),
		c.Srcs("var a = f();"),
		c.Expected(""))
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNoSideEffectAnnotation10()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNoSideEffectAnnotation10() {
	c.Test(
		c.Externs("var o = {}; o.f = function(){};"), //
		c.Srcs("var a = o.f();"),
		c.Expected("o.f();"))
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNoSideEffectAnnotation11()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNoSideEffectAnnotation11() {
	c.Test(
		c.Externs("var o = {}; o.f = /** @nosideeffects */ function(){};"),
		c.Srcs("var a = o.f();"),
		c.Expected(""))
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNoSideEffectAnnotation12()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNoSideEffectAnnotation12() {
	c.Test(c.Externs("function c(){}"), c.Srcs("var a = new c()"), c.Expected("new c()"))
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNoSideEffectAnnotation13()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNoSideEffectAnnotation13() {
	c.Test(c.Externs("/** @nosideeffects */ function c(){}"), c.Srcs("var a = new c"), c.Expected(""))
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNoSideEffectAnnotation14()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNoSideEffectAnnotation14() {
	externs :=
		"function c(){};\n" +
			"c.prototype.f = /**@nosideeffects*/function(){};\n"
	c.Test(c.Externs(externs), c.Srcs("var o = new c; var a = o.f()"), c.Expected("new c"))
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNoSideEffectAnnotation15()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNoSideEffectAnnotation15() {
	c.Test(
		c.Externs("/** @nosideeffects */ function c(){}", "c.prototype.f = function(){};"),
		c.Srcs("var a = (new c).f()"),
		c.Expected("(new c).f()"))
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNoSideEffectAnnotation16()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNoSideEffectAnnotation16() {
	c.Test(
		c.Externs(
			"/** @nosideeffects */ function c(){}",
			"c.prototype.f = /** @nosideeffects */ function(){};"),
		c.Srcs("var a = (new c).f()"),
		c.Expected(""))
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testFunctionPrototype()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestFunctionPrototype() {
	c.Test("var a = 5; Function.prototype.foo = function() {return a;}", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testTopLevelClass1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestTopLevelClass1() {
	c.Test("var Point = function() {}; Point.prototype.foo = function() {}", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testTopLevelClass2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestTopLevelClass2() {
	c.Test(
		"var Point = {}; Point.prototype.foo = function() {};externfoo = new Point()",
		"var Point = {};                                     externfoo = new Point()")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testTopLevelClass3()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestTopLevelClass3() {
	c.Test("function Point() {this.me_ = Point}", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testTopLevelClass4()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestTopLevelClass4() {
	c.Test(
		"function f(){} function A(){} A.prototype = {x: function() {}}; f();",
		"function f(){} f();")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testTopLevelClass5()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestTopLevelClass5() {
	c.Test(
		"function f(){} function A(){} A.prototype = {x: function() { f(); }}; new A();",
		"               function A(){} A.prototype = {                      }; new A();")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testTopLevelClass6()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestTopLevelClass6() {
	c.TestSame(
		"function f(){} function A(){}\n" +
			"A.prototype = {x: function() { f(); }}; new A().x();\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testTopLevelClass7()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestTopLevelClass7() {
	c.Test("A.prototype.foo = function(){}; function A() {}", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNamespacedClass1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNamespacedClass1() {
	c.Test("var foo = {};foo.bar = {};foo.bar.prototype.baz = {}", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNamespacedClass2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNamespacedClass2() {
	c.Test(
		"var foo = {}; foo.bar = {}; foo.bar.prototype.baz = {}; window.z = new foo.bar()",
		"var foo = {}; foo.bar = {};                             window.z = new foo.bar()")
}

// TODO(b/66971163): enable or remove this
// Java: RemoveUnusedCodeNameAnalyzerTest.disabledTestNamespacedClass3()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestDisabledTestNamespacedClass3() {
	c.T().Skip("@Ignore")
	c.Test("var a = {}; a.b = function() {}; a.b.prototype = {x: function() {}};", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNamespacedClass4()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNamespacedClass4() {
	c.Test(
		"function f(){}\n"+
			"var a = {};\n"+
			"a.b = function() {};\n"+
			"a.b.prototype = {x: function() { f(); }};\n"+
			"new a.b();\n",
		"var a = {};\n"+
			"a.b = function() {};\n"+
			"a.b.prototype = {                      };\n"+
			"new a.b();\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNamespacedClass5()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNamespacedClass5() {
	c.TestSame(
		"function f(){} var a = {}; a.b = function() {};\n" +
			"a.b.prototype = {x: function() { f(); }}; new a.b().x();\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testEs6Class()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestEs6Class() {
	c.Test("class C {}", "")

	c.Test("class C {constructor() {} }", "")

	c.TestSame("class C {} var c = new C(); use(c);")

	c.TestSame(
		"class C {\n" +
			"  constructor() {\n" +
			"    this.x = 1;\n" +
			"  }\n" +
			"  add() {\n" +
			"    this.x++\n" +
			"  }\n" +
			"}\n" +
			"var c = new C;\n" +
			"c.add();\n" +
			"use(c.x);\n")

	c.Test("class C{} class D{} var d = new D;", "")

	c.Test("{class C{} }", "{}")
	c.TestSame("{class C{} var c = new C; use(c)}")

	c.Test(
		"function f() { return class C {} } doThing1();", // preserve newline
		"                                   doThing1();")
	c.TestSame("export class C {}")

	// class expressions

	c.Test("var c = class{}", "")
	c.TestSame("var c = class{}; use(c);")

	c.Test("var c = class C {}", "")
	c.Test(
		"var c = class C {}; use(c);", // preserve newline
		"var c = class   {}; use(c);")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testInnerClassNameInstanceofCheck()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestInnerClassNameInstanceofCheck() {
	c.TestSame(
		"window.Class = class MyClass {\n" +
			"  constructor() {\n" +
			"    if (this instanceof MyClass) {\n" +
			"      console.log(\"test\");\n" +
			"    }\n" +
			"  }\n" +
			"};\n" +
			"new window.Class();\n")
	c.TestSame(
		"/** @constructor */\n" +
			"window.Class = function MyClass() {\n" +
			"  if (this instanceof MyClass) {\n" +
			"    console.log(\"test\");\n" +
			"  }\n" +
			"};\n" +
			"new window.Class();\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testEs6ClassExtends()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestEs6ClassExtends() {
	c.TestSame("class D {} class C extends D {} var c = new C; c.g();")
	c.Test("class D {} class C extends D {}", "")
}

/**
 * @bug 67430253
 */
// Java: RemoveUnusedCodeNameAnalyzerTest.testEs6ClassExtendsQualifiedName1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestEs6ClassExtendsQualifiedName1() {
	c.TestSame("var ns = {}; ns.Class1 = class {}; class Class2 extends ns.Class1 {}; use(Class2);")
}

/**
 * @bug 67430253
 */
// Java: RemoveUnusedCodeNameAnalyzerTest.testEs6ClassExtendsQualifiedName2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestEs6ClassExtendsQualifiedName2() {
	c.Test(
		"var ns = {}; ns.Class1 = class {}; use(ns.Class1); class Class2 extends ns.Class1 {}",
		"var ns = {}; ns.Class1 = class {}; use(ns.Class1);")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testAssignmentToThisPrototype()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestAssignmentToThisPrototype() {
	c.TestSame(
		"Function.prototype.inherits = function(parentCtor) {\n" +
			"  function tempCtor() {};\n" +
			"  tempCtor.prototype = parentCtor.prototype;\n" +
			"  this.superClass_ = parentCtor.prototype;\n" +
			"  this.prototype = new tempCtor();\n" +
			"  this.prototype.constructor = this;\n" +
			"};\n" +
			"/** @constructor */ function A() {}\n" +
			"/** @constructor */ function B() {}\n" +
			"B.inherits(A);\n" +
			"use(B.superClass_);\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testAssignmentToCallResultPrototype()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestAssignmentToCallResultPrototype() {
	c.TestSame("function f() { return function(){}; } f().prototype = {};")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testAssignmentToExternPrototype()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestAssignmentToExternPrototype() {
	c.TestSame("externfoo.prototype = {};")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testAssignmentToUnknownPrototype()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestAssignmentToUnknownPrototype() {
	c.DisableCompareJsDoc() // multistage compilation simplifies suppressions
	c.TestSame(
		"/** @suppress {duplicate} */ var window;\n" +
			"window['a'].prototype = {};\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testBug2099540()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestBug2099540() {
	c.DisableCompareJsDoc() // multistage compilation simplifies suppressions
	c.TestSame(
		"/** @suppress {duplicate} */ var document;\n" +
			"/** @suppress {duplicate} */ var window;\n" +
			"var klass;\n" +
			"window[klass].prototype = document.createElement('p')['__proto__'];\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testOtherGlobal()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestOtherGlobal() {
	c.TestSame("goog.global.foo = bar(); function bar(){}")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testExternName1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestExternName1() {
	c.TestSame("top.z = bar(); function bar(){}")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testExternName2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestExternName2() {
	c.TestSame("top['z'] = bar(); function bar(){}")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testInherits1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestInherits1() {
	c.Test("var a = {}; var b = {}; goog.inherits(b, a)", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testInherits2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestInherits2() {
	c.TestSame("var a = {}; externfoo.inherits(a);")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testInherits3()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestInherits3() {
	c.TestSame("var a = {}; goog.inherits(externfoo, a);")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testInherits4()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestInherits4() {
	c.Test("var b = {}; goog.inherits(b, externfoo);", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testInherits5()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestInherits5() {
	c.TestSame("var a = {}; goog.inherits(externfoo, a);")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testConstants1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestConstants1() {
	c.TestSame("var bar = function(){}; var EXP_FOO = true; if (EXP_FOO) bar();")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testConstants2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestConstants2() {
	c.Test(
		"var bar = function(){}; var EXP_FOO = true; var EXP_BAR = true;\n"+
			"if (EXP_FOO) bar();\n",
		"var bar = function(){}; var EXP_FOO = true; if (EXP_FOO) bar();")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testExpressions1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestExpressions1() {
	c.Test("var foo={}; foo.A='A'; foo.AB=foo.A+'B'; foo.ABC=foo.AB+'C'", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testExpressions2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestExpressions2() {
	c.Test("var foo={}; foo.A='A'; foo.AB=foo.A+'B'; this.ABC=foo.AB+'C'", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testExpressions3()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestExpressions3() {
	c.TestSame("var foo = 2; window.bar(foo + 3)")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSetCreatingReference()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSetCreatingReference() {
	c.Test(
		"var foo; var bar = function(){foo=6;}; bar();",
		"         var bar = function(){      }; bar();")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testAnonymous1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestAnonymous1() {
	c.TestSame("function foo() {}; function bar() {}; foo(function() {bar()})")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testAnonymous2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestAnonymous2() {
	c.Test("var foo;(function(){foo=6;})()", "(function(){})()")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testAnonymous3()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestAnonymous3() {
	c.TestSame("var foo; (function(){ if(!foo)foo=6; })()")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testAnonymous4()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestAnonymous4() {
	c.TestSame("var foo; (function(){ foo=6; })(); externfoo=foo;")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testAnonymous5()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestAnonymous5() {
	c.TestSame(
		"var foo;\n" +
			"(function(){ foo=function(){ bar() }; function bar(){} })();\n" +
			"foo();\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testAnonymous6()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestAnonymous6() {
	c.TestSame(
		"function foo(){}\n" +
			"function bar(){}\n" +
			"foo(function(){externfoo = bar});\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testAnonymous7()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestAnonymous7() {
	c.TestSame(
		"var foo;\n" +
			"(function (){ function bar(){ externfoo = foo; } bar(); })();\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testAnonymous8()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestAnonymous8() {
	c.TestSame(
		"var foo;\n" +
			"(function (){ var g=function(){ externfoo = foo; }; g(); })();\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testAnonymous9()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestAnonymous9() {
	c.TestSame(
		"function foo(){}\n" +
			"function bar(){}\n" +
			"foo(function(){ function baz(){ externfoo = bar; } baz(); });\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testFunctions1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestFunctions1() {
	c.Test(
		"var foo = null; function baz() {sideEffect()} function bar() {foo=baz();} bar();",
		"                function baz() {sideEffect()} function bar() {    baz();} bar();")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testFunctions2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestFunctions2() {
	c.Test(
		"var foo; foo = function() {var a = bar()}; var bar = function(){sideEffect()}; foo();",
		"var foo; foo = function() {        bar()}; var bar = function(){sideEffect()}; foo();")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testGetElem1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestGetElem1() {
	c.TestSame(
		"var foo = {}; foo.bar = {}; foo.bar.baz = {a: 5, b: 10};\n" +
			"var fn = function() {window[foo.bar.baz.a] = 5;}; fn()\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testGetElem2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestGetElem2() {
	c.TestSame(
		"var foo = {}; foo.bar = {}; foo.bar.baz = {a: 5, b: 10};\n" +
			"var fn = function() {this[foo.bar.baz.a] = 5;}; fn()\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testGetElem3()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestGetElem3() {
	c.TestSame("var foo = {'i': 0, 'j': 1}; foo['k'] = 2; top.foo = foo;")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testIf1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestIf1() {
	c.Test("var foo = {};if(externfoo)foo.bar=function(){};", "if(externfoo);")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testIf2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestIf2() {
	c.Test("var e = false;var foo = {};if(e)foo.bar=function(){};", "var e = false;if(e);")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testIf3()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestIf3() {
	c.Test("var e = false;var foo = {};if(e + 1)foo.bar=function(){};", "var e = false;if(e + 1);")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testIf4()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestIf4() {
	c.Test("var e = false, f;var foo = {};if(f=e)foo.bar=function(){};", "var e = false;if(e);")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testIf4a()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestIf4a() {
	// TODO(johnlenz): fix this.
	c.TestSame("var e = [], f;if(f=e);f[0] = 1;")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testIf4b()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestIf4b() {
	// TODO(johnlenz): fix this.
	c.Test("var e = [], f;if(e=f);f[0] = 1;", "var f;if(f);f[0] = 1;")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testIf4c()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestIf4c() {
	c.Test("var e = [], f;if(f=e);e[0] = 1;", "var e = [];if(e);e[0] = 1;")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testIf5()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestIf5() {
	c.Test(
		"var e = false, f;var foo = {};if(f = e + 1)foo.bar=function(){};",
		"var e = false;if(e + 1);")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testIfElse()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestIfElse() {
	c.Test(
		"var foo = {}; if(externfoo) foo.bar = function(){}; else foo.bar = function(){};",
		"              if(externfoo)                       ; else                       ;")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testWhile()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestWhile() {
	c.Test(
		"var foo = {}; while(doThing1()) foo.bar=function(){};",
		"              while(doThing1())                     ;")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testForIn()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestForIn() {
	c.Test(
		"var foo = {};for(var e in externfoo)foo.bar=function(){};",
		"             for(var e in externfoo)                    ;")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testForOf()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestForOf() {
	c.Test(
		"var foo = {};for(var e of externfoo)foo.bar=function(){};",
		"             for(var e of externfoo)                    ;")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testDo()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestDo() {
	c.Test("var cond = false;do {var a = 1} while (cond)", "var cond = false;do {} while (cond)")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSetterInForStruct1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSetterInForStruct1() {
	c.Test("var j = 0; for (var i = 1; i = 0; j++);", "var j = 0; for (; 0; j++);")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSetterInForStruct2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSetterInForStruct2() {
	c.Test(
		"var Class = function() {};\n"+
			"for (var i = 1; Class.prototype.property_ = 0; i++);\n",
		"for (var i = 1; 0; i++);")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSetterInForStruct3()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSetterInForStruct3() {
	c.Test(
		c.Externs("function f(){} function g() {} function h() {}"),
		c.Srcs("var j = 0;                      for (var i = 1 + f() + g() + h(); i = 0; j++);"),
		c.Expected("var j = 0; 1 + f() + g() + h(); for (                           ;     0; j++);"))
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSetterInForStruct4()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSetterInForStruct4() {
	c.Test(
		c.Externs("function f(){} function g() {} function h() {}"),
		c.Srcs(
			"var i = 0; var j = 0;                      for (i = 1 + f() + g() + h(); i = 0;\n"+
				" j++);\n"),
		c.Expected(
			"          var j = 0; 1 + f() + g() + h(); for (                       ;     0;\n"+
				"j++);\n"))
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSetterInForStruct5()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSetterInForStruct5() {
	c.Test(
		c.Externs("function f(){} function g() {} function h() {}"),
		c.Srcs("var i = 0, j = 0; for (i = f(), j = g(); 0;);"),
		c.Expected("                  for (    f(),     g(); 0;);"))
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSetterInForStruct6()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSetterInForStruct6() {
	c.Test(
		c.Externs("function f(){} function g() {} function h() {}"),
		c.Srcs("var i = 0, j = 0, k = 0; for (i = f(), j = g(), k = h(); i = 0;);"),
		c.Expected("                         for (    f(),     g(),     h();     0;);"))
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSetterInForStruct7()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSetterInForStruct7() {
	c.Test(
		c.Externs("function f(){} function g() {} function h() {}"),
		c.Srcs("var i = 0, j = 0, k = 0; for (i = 1, j = 2, k = 3; i = 0;);"),
		c.Expected("                         for (                   ;     0;);"))
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSetterInForStruct8()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSetterInForStruct8() {
	c.Test(
		c.Externs("function f(){} function g() {} function h() {}"),
		c.Srcs("var i = 0, j = 0, k = 0; for (i = 1, j = i, k = 2; i = 0;);"),
		c.Expected("                         for (                   ;     0;);"))
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSetterInForStruct9()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSetterInForStruct9() {
	c.Test(
		"var Class = function() {};\n"+
			"for (var i = 1; Class.property_ = 0; i++);\n",
		"for (var i = 1; 0; i++);")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSetterInForStruct10()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSetterInForStruct10() {
	c.Test(
		"var Class = function() {}; for (var i = 1; Class.property_ = 0; i = 2);",
		"                           for (         ; 0                  ; 0    );")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSetterInForStruct11()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSetterInForStruct11() {
	c.Test(
		"var Class = function() {};\n"+
			"for (;Class.property_ = 0;);\n",
		"for (;0;);")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSetterInForStruct12()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSetterInForStruct12() {
	c.Test(
		"var a = 1; var Class = function() {};\n"+
			"for (;Class.property_ = a;);\n",
		"var a = 1; for (; a;);")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSetterInForStruct13()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSetterInForStruct13() {
	c.Test(
		"var a = 1; var Class = function() {};\n"+
			"for (Class.property_ = a; 0 ;);\n",
		"for (; 0;);")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSetterInForStruct14()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSetterInForStruct14() {
	c.Test(
		"var a = 1; var Class = function() {}; for (; 0; Class.property_ = a);",
		"                                      for (; 0; 0                  );")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSetterInForStruct15()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSetterInForStruct15() {
	c.Test(
		"var Class = function() {};\n"+
			"for (var i = 1; 0; Class.prototype.property_ = 0);\n",
		"for (; 0; 0);")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSetterInForStruct16()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSetterInForStruct16() {
	c.Test(
		"var Class = function() {};\n"+
			"for (var i = 1; i = 0; Class.prototype.property_ = 0);\n",
		"for (; 0; 0);")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSetterInForIn1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSetterInForIn1() {
	c.Test(
		"var foo = {}; var bar; for(var e in bar = foo.a);",
		"var foo = {};          for(var e in       foo.a);")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSetterInForIn2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSetterInForIn2() {
	c.TestSame("var foo = {}; var bar; for(var e in bar = foo.a); bar")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSetterInForIn3()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSetterInForIn3() {
	c.TestSame("var foo = {}; var bar; for(var e in bar = foo.a); bar.b = 3")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSetterInForIn4()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSetterInForIn4() {
	c.TestSame("var foo = {}; var bar; for (var e in bar = foo.a); bar.b = 3; foo.a")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSetterInForIn5()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSetterInForIn5() {
	c.TestSame("var foo = {}; var bar; for (var e in foo.a) { bar = e } bar.b = 3; foo.a")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSetterInForIn6()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSetterInForIn6() {
	c.TestSame("var foo = {}; for(var e in foo);")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSetterInForOf1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSetterInForOf1() {
	c.Test(
		"var foo = {}; var bar; for(var e of bar = foo.a);",
		"var foo = {};          for(var e of       foo.a);")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSetterInForOf2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSetterInForOf2() {
	c.TestSame("var foo = {}; var bar; for(var e of bar = foo.a); bar")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSetterInForOf3()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSetterInForOf3() {
	c.TestSame("var foo = {}; var bar; for(var e of bar = foo.a); bar.b = 3")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSetterInForOf4()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSetterInForOf4() {
	c.TestSame("var foo = {}; var bar; for (var e of bar = foo.a); bar.b = 3; foo.a")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSetterInForOf5()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSetterInForOf5() {
	c.TestSame("var foo = {}; var bar; for (var e of foo.a) { bar = e } bar.b = 3; foo.a")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSetterInForOf6()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSetterInForOf6() {
	c.TestSame("var foo = {}; for(var e of foo);")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSetterInIfPredicate()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSetterInIfPredicate() {
	c.Test(
		"var a = 1; var Class = function() {}; if (Class.property_ = a);",
		"var a = 1;                            if (                  a);")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSetterInWhilePredicate()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSetterInWhilePredicate() {
	c.Test(
		"var a = 1;\n"+
			"var Class = function() {};\n"+
			"while (Class.property_ = a);\n",
		"var a = 1; for (;a;) {}")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSetterInDoWhilePredicate()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSetterInDoWhilePredicate() {
	c.Test(
		"var a = 1; var Class = function() {}; do {} while (Class.property_ = a);",
		"var a = 1;                            do ;  while (                  a);")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSetterInSwitchInput()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSetterInSwitchInput() {
	c.Test(
		"var a = 1;var Class = function() {}; switch (Class.property_ = a) {default:}",
		"var a = 1;                           switch (                  a) {default:}")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testComplexAssigns()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestComplexAssigns() {
	c.Test("var x = 0; x += 3; x *= 5;", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNestedAssigns1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNestedAssigns1() {
	c.Test("var x = 0; var y = x = 3; window.alert(y);", "var y = 3; window.alert(y);")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNestedAssigns2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNestedAssigns2() {
	c.TestSame("var x = 0; var y = x = {}; x.b = 3; window.alert(y);")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testComplexNestedAssigns1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestComplexNestedAssigns1() {
	// TODO(bradfordcsmith): Make RemoveUnusedCode smarter, so that we can eliminate y.
	c.TestSame("var x = 0; var y = 2; y += x = 3; window.alert(x);")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testComplexNestedAssigns2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestComplexNestedAssigns2() {
	c.Test(
		"var x = 0; var y = 2; y += x = 3; window.alert(y);",
		"var y = 2; y += 3; window.alert(y);")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testComplexNestedAssigns3()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestComplexNestedAssigns3() {
	c.Test("var x = 0; var y = x += 3; window.alert(x);", "var x = 0; x += 3; window.alert(x);")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testComplexNestedAssigns4()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestComplexNestedAssigns4() {
	c.TestSame("var x = 0; var y = x += 3; window.alert(y);")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testPrototypePropertySetInLocalScope1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestPrototypePropertySetInLocalScope1() {
	c.Test(
		"(function() { var x = function(){}; x.prototype.bar = 3; })();",
		"(function() {                                            })();")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testPrototypePropertySetInLocalScope2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestPrototypePropertySetInLocalScope2() {
	c.Test(
		"var x = function(){}; (function() { x.prototype.bar = 3; })();",
		"                      (function() {                      })();")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testPrototypePropertySetInLocalScope3()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestPrototypePropertySetInLocalScope3() {
	c.Test("var x = function(){ x.prototype.bar = 3; };", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testPrototypePropertySetInLocalScope4()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestPrototypePropertySetInLocalScope4() {
	c.Test("var x = {}; x.foo = function(){ x.foo.prototype.bar = 3; };", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testPrototypePropertySetInLocalScope5()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestPrototypePropertySetInLocalScope5() {
	c.Test("var x = {}; x.prototype.foo = 3;", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testPrototypePropertySetInLocalScope6()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestPrototypePropertySetInLocalScope6() {
	c.TestSame("var x = {}; x.prototype.foo = 3; use(x.prototype.foo)")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testPrototypePropertySetInLocalScope7()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestPrototypePropertySetInLocalScope7() {
	c.TestSame("var x = {}; x.foo = 3; use(x.foo)")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRValueReference1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRValueReference1() {
	c.TestSame("var a = 1; a")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRValueReference2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRValueReference2() {
	c.TestSame("var a = 1; 1+a")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRValueReference3()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRValueReference3() {
	c.TestSame("var x = {}; x.prototype.foo = 3; var a = x.prototype.foo; 1+a")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRValueReference4()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRValueReference4() {
	c.TestSame("var x = {}; x.prototype.foo = 3; use(x.prototype.foo);")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRValueReference5()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRValueReference5() {
	c.TestSame("var x = {}; x.prototype.foo = 3; 1+x.prototype.foo")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRValueReference6()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRValueReference6() {
	c.TestSame("var x = {}; var idx = 2; x[idx]")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testUnhandledTopNode()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestUnhandledTopNode() {
	c.TestSame(
		"function Foo() {}; Foo.prototype.isBar = function() {};\n" +
			"function Bar() {}; Bar.prototype.isFoo = function() {};\n" +
			"var foo = new Foo(); var bar = new Bar();\n" +
			"// The boolean AND here is currently unhandled by this pass, but it\n" +
			"// it should not cause it to blow up.\n" +
			"var cond = foo.isBar() && bar.isFoo();\n" +
			"if (cond) {window.alert('hello');}")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testPropertyDefinedInGlobalScope()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestPropertyDefinedInGlobalScope() {
	c.TestSame(
		"function Foo() {}; var x = new Foo(); x.cssClass = 'bar';\n" +
			"window.alert(x);\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testConditionallyDefinedFunction1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestConditionallyDefinedFunction1() {
	c.TestSame("var g; externfoo.x || (externfoo.x = function() { g; })")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testConditionallyDefinedFunction2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestConditionallyDefinedFunction2() {
	c.TestSame("var g; 1 || (externfoo.x = function() { g; })")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testConditionallyDefinedFunction3()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestConditionallyDefinedFunction3() {
	c.Test(
		"var a = {}; doThing1() || (a.f = function() { externfoo = 1; } || doThing2());",
		"            doThing1() || (      function() { externfoo = 1; } || doThing2());")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testGetElemOnThis()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestGetElemOnThis() {
	c.TestSame("var a = 3; this['foo'] = a;")
	c.TestSame("this['foo'] = 3;")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveInstanceOfOnly()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveInstanceOfOnly() {
	c.Test(
		"function Foo() {}\n"+
			"Foo.prototype.isBar = function() {};\n"+
			"var x;\n"+
			"if (x instanceof Foo) { window.alert(x); }\n",
		"var x;\n"+
			"if (false           ) { window.alert(x); }\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveLocalScopedInstanceOfOnly()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveLocalScopedInstanceOfOnly() {
	c.Test(
		"function Foo() {}\n"+
			"function Bar(x) { use(x instanceof Foo); }\n"+
			"externfoo.x = new Bar({});\n",
		"function Bar(x) { use(false           ); }\n"+
			"externfoo.x = new Bar({});\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveInstanceOfWithReferencedMethod()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveInstanceOfWithReferencedMethod() {
	c.Test(
		"function Foo() {}\n"+
			"Foo.prototype.isBar = function() {};\n"+
			"var x;\n"+
			"if (x instanceof Foo) { window.alert(x.isBar()); }\n",
		"var x;\n"+
			"if (false           ) { window.alert(x.isBar()); }\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testDoNotChangeReferencedInstanceOf()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestDoNotChangeReferencedInstanceOf() {
	c.TestSame(
		"function Foo() {}\n" +
			"var x = new Foo();\n" +
			"if (x instanceof Foo) { window.alert(x); }\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testDoNotChangeReferencedLocalScopedInstanceOf()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestDoNotChangeReferencedLocalScopedInstanceOf() {
	c.TestSame(
		"function Foo() {};\n" +
			"externfoo.x = new Foo();\n" +
			"function Bar() {\n" +
			"  var x;\n" +
			"  if (x instanceof Foo) { window.alert(x); }\n" +
			"}\n" +
			"externfoo.y = new Bar();\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testDoNotChangeLocalScopeReferencedInstanceOf()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestDoNotChangeLocalScopeReferencedInstanceOf() {
	c.TestSame(
		"function Foo() {}\n" +
			"function Bar() { use(new Foo()); }\n" +
			"externfoo.x = new Bar();\n" +
			"var x;\n" +
			"if (x instanceof Foo) { window.alert(x); }\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testDoNotChangeLocalScopeReferencedLocalScopedInstanceOf()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestDoNotChangeLocalScopeReferencedLocalScopedInstanceOf() {
	c.TestSame(
		"function Foo() {}\n" +
			"function Bar() { new Foo(); }\n" +
			"Bar.prototype.func = function(x) {\n" +
			"  if (x instanceof Foo) { window.alert(x); }\n" +
			"};\n" +
			"new Bar().func();\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testDoNotChangeLocalScopeReferencedLocalScopedInstanceOf2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestDoNotChangeLocalScopeReferencedLocalScopedInstanceOf2() {
	c.Test(
		"function Foo() {}\n"+
			"var createAxis = function(f) { return window.passThru(f); };\n"+
			"var axis = createAxis(function(test) {\n"+
			"  return test instanceof Foo;\n"+
			"});\n",
		"var createAxis = function(f) { return window.passThru(f); };\n"+
			"createAxis(function(test) {\n"+
			"  return false;\n"+
			"});\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testDoNotChangeInstanceOfGetElem()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestDoNotChangeInstanceOfGetElem() {
	c.TestSame(
		"var goog = {};\n" +
			"function f(obj, name) {\n" +
			"  if (obj instanceof goog[name]) {\n" +
			"    return name;\n" +
			"  }\n" +
			"}\n" +
			"window['f'] = f;\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testIssue2822()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestIssue2822() {
	c.TestSame(
		"var C = function C() {\n" +
			"  if (!(this instanceof C)) {\n" +
			"    throw new Error('not an instance');\n" +
			"  }\n" +
			"}\n" +
			"use(new C());\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testWeirdnessOnLeftSideOfPrototype()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestWeirdnessOnLeftSideOfPrototype() {
	// This checks a bug where 'x' was removed, but the function referencing
	// it was not, causing problems.
	c.TestSame("var x = 3; (function() {}).z = function() { return x; };")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testDoNotChangeInstanceOfGetprop()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestDoNotChangeInstanceOfGetprop() {
	c.TestSame(
		"function f(obj) {\n" +
			"  if (obj instanceof window.MouseEvent) obj.preventDefault();\n" +
			"}\n" +
			"window['f'] = f;\n")
}

// TODO(b/66971163): Enable or remove this test.
// Java: RemoveUnusedCodeNameAnalyzerTest.disabledTestShortCircuit1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestDisabledTestShortCircuit1() {
	c.T().Skip("@Ignore")
	c.Test("var a = doThing1() || 1", "doThing1()")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testShortCircuit2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestShortCircuit2() {
	c.Test("var a = 1 || doThing1()", "1 || doThing1()")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testShortCircuit3()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestShortCircuit3() {
	c.Test("var a = doThing1() || doThing2()", "doThing1() || doThing2()")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testShortCircuit4()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestShortCircuit4() {
	c.Test("var a = doThing1() || 3 || doThing2()", "doThing1() || 3 || doThing2()")
}

// TODO(b/66971163): Enable or remove this test.
// Java: RemoveUnusedCodeNameAnalyzerTest.disabledTestShortCircuit5()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestDisabledTestShortCircuit5() {
	c.T().Skip("@Ignore")
	c.Test("var a = doThing1() && 1", "doThing1()")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testShortCircuit6()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestShortCircuit6() {
	c.Test("var a = 1 && doThing1()", "1 && doThing1()")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testShortCircuit7()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestShortCircuit7() {
	c.Test("var a = doThing1() && doThing2()", "doThing1() && doThing2()")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testShortCircuit8()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestShortCircuit8() {
	c.Test("var a = doThing1() && 3 && doThing2()", "doThing1() && 3 && doThing2()")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRhsReference1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRhsReference1() {
	c.TestSame("var a = 1; a")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRhsReference2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRhsReference2() {
	c.TestSame("var a = 1; a || doThing1()")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRhsReference3()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRhsReference3() {
	c.TestSame("var a = 1; 1 || a")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRhsReference4()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRhsReference4() {
	c.Test("var a = 1; var b = a || doThing1()", "var a = 1; a || doThing1()")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRhsReference5()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRhsReference5() {
	c.TestSame("var a = 1, b = 5; a; use(b)")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRhsAssign1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRhsAssign1() {
	c.Test(
		"var foo, bar; foo || (bar = 1)", // preserve newline
		"var foo;      foo || 0        ")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRhsAssign2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRhsAssign2() {
	c.Test(
		"var foo, bar, baz; foo || (baz = bar = 1)", // preserve newline
		"var foo;           foo || 0              ")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRhsAssign3()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRhsAssign3() {
	c.TestSame("var foo = null; foo || (foo = 1)")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRhsAssign4()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRhsAssign4() {
	c.Test("var foo = null; foo = (foo || 1)", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRhsAssign5()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRhsAssign5() {
	c.Test("var a = 3, foo, bar; foo || (bar = a)", "var        foo     ; foo || 0        ")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRhsAssign6()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRhsAssign6() {
	c.Test(
		"function Foo(){} var foo = null;\n"+
			"var f = function () {foo || (foo = new Foo()); return foo}\n",
		"")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRhsAssign7()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRhsAssign7() {
	c.TestSame(
		"function Foo(){} var foo = null;\n" +
			"var f = function () {foo || (foo = new Foo())}; f()\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRhsAssign8()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRhsAssign8() {
	c.Test(
		"function Foo(){}\n"+
			"var foo = null;\n"+
			"var f = function () {\n"+
			"  (foo = new Foo()) || doThing1();\n"+
			"};\n"+
			"f()\n",
		"function Foo(){}\n"+
			"\n"+
			"var f = function () {\n"+
			"         new Foo()  || doThing1();\n"+
			"};\n"+
			"f()\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRhsAssign9()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRhsAssign9() {
	c.Test(
		"function Foo(){} var foo = null;\n"+
			"var f = function () {1 + (foo = new Foo()); return foo}\n",
		"")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testAssignWithOr1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestAssignWithOr1() {
	c.TestSame(
		"var foo = null;\n" +
			"var f = window.a || function () {return foo}; f()\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testAssignWithOr2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestAssignWithOr2() {
	c.Test("var foo = null; var f = window.a || function () {return foo};", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testAssignWithAnd1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestAssignWithAnd1() {
	c.TestSame(
		"var foo = null;\n" +
			"var f = window.a && function () {return foo}; f()\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testAssignWithAnd2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestAssignWithAnd2() {
	c.Test("var foo = null; var f = window.a && function () {return foo};", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testAssignWithHook1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestAssignWithHook1() {
	c.TestSame(
		"function Foo(){} var foo = null;\n" +
			"var f = window.a ?\n" +
			"    function () {return new Foo()} : function () {return foo}; f()\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testAssignWithHook2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestAssignWithHook2() {
	c.Test(
		"function Foo(){} var foo = null;\n"+
			"var f = window.a ?\n"+
			"    function () {return new Foo()} : function () {return foo};\n",
		"")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testAssignWithHook2a()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestAssignWithHook2a() {
	c.Test(
		"function Foo(){} var foo = null;\n"+
			"var f; f = window.a ?\n"+
			"    function () {return new Foo()} : function () {return foo};\n",
		"")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testAssignWithHook3()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestAssignWithHook3() {
	c.TestSame(
		"function Foo(){} var foo = null; var f = {};\n" +
			"f.b = window.a ?\n" +
			"    function () {return new Foo()} : function () {return foo}; f.b()\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testAssignWithHook4()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestAssignWithHook4() {
	c.Test(
		"function Foo(){} var foo = null; var f = {};\n"+
			"f.b = window.a ?\n"+
			"    function () {return new Foo()} : function () {return foo};\n",
		"")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testAssignWithHook5()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestAssignWithHook5() {
	c.TestSame(
		"function Foo(){} var foo = null; var f = {};\n" +
			"f.b = window.a ? function () {return new Foo()} :\n" +
			"    window.b ? function () {return foo} :\n" +
			"    function() { return Foo }; f.b()\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testAssignWithHook6()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestAssignWithHook6() {
	c.Test(
		"function Foo(){} var foo = null; var f = {};\n"+
			"f.b = window.a ? function () {return new Foo()} :\n"+
			"    window.b ? function () {return foo} :\n"+
			"    function() { return Foo };\n",
		"")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testAssignWithHook7()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestAssignWithHook7() {
	c.TestSame(
		"function Foo(){} var foo = null;\n" +
			"var f = window.a ? new Foo() : foo;\n" +
			"f()\n")
}

// TODO(b/66971163): enable or remove this test
// Java: RemoveUnusedCodeNameAnalyzerTest.disabledTestAssignWithHook8()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestDisabledTestAssignWithHook8() {
	c.T().Skip("@Ignore")
	c.Test(
		"function Foo(){} var foo = null; var f = window.a  ? new Foo() : foo;",
		"function Foo(){}                         window.a && new Foo()      ;")
}

// TODO(b/66971163): enable or remove this test
// Java: RemoveUnusedCodeNameAnalyzerTest.disabledTestAssignWithHook9()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestDisabledTestAssignWithHook9() {
	c.T().Skip("@Ignore")
	c.Test(
		"function Foo(){} var foo = null; var f = {};f.b = window.a  ? new Foo() : foo;",
		"function Foo(){}                                  window.a && new Foo()      ;")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testAssign1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestAssign1() {
	c.Test(
		"function Foo(){} var foo = null; var f = {};\n"+
			"f.b = window.a;\n",
		"")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testAssign2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestAssign2() {
	c.Test(
		"function Foo(){} var foo = null; var f = {};\n"+
			"f.b = window;\n",
		"")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testAssign3()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestAssign3() {
	c.Test(
		"var f = {};\n"+
			"f.b = window;\n",
		"")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testAssign4()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestAssign4() {
	c.Test(
		"function Foo(){sideEffect()} var foo = null; var f = {};\n"+
			"f.b = new Foo();\n",
		"function Foo(){sideEffect()} new Foo()")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testAssign5()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestAssign5() {
	c.Test(
		"function Foo(){} var foo = null; var f = {};\n"+
			"f.b = foo;\n",
		"")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testAssignWithCall()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestAssignWithCall() {
	c.Test("var fun, x; (fun = function(){ x; })();", "var x; (function(){ x; })();")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testAssignWithCall2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestAssignWithCall2() {
	c.Test(
		"var fun, x; (123, fun = function(){ x; })();",
		"var      x; (123,       function(){ x; })();")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNestedAssign1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNestedAssign1() {
	c.Test("var a, b = a = 1, c = 2", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNestedAssign2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNestedAssign2() {
	c.Test(
		"var a, b = a = 1; use(b)", // preserve newline
		"var b = 1;        use(b)")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNestedAssign3()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNestedAssign3() {
	c.Test(
		"var a, b = a = 1; a = b = 2; use(b)", // preserve newline
		"var    b =     1;     b = 2; use(b)")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNestedAssign4()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNestedAssign4() {
	c.Test(
		"var a, b = a = 1; b = a = 2; use(b)", // preserve newline
		"var    b     = 1; b =     2; use(b)")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNestedAssign5()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNestedAssign5() {
	c.Test("var a, b = a = 1; b = a = 2", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNestedAssign15()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNestedAssign15() {
	c.Test("var a, b, c; c = b = a = 2", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNestedAssign6()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNestedAssign6() {
	c.TestSame("var a, b, c; a = b = c = 1; use(a, b, c)")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNestedAssign7()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNestedAssign7() {
	c.TestSame("var a = 0, j = 0, i = []; a = i[j] = 1; use(a, i[j])")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNestedAssign8()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNestedAssign8() {
	c.TestSame(
		"function f(){\n" +
			"  this.externProp1 = this.externProp2 = use(this.hiddenInput_, externfoo);\n" +
			"}\n" +
			"f()\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRefChain1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRefChain1() {
	c.Test("var a = 1; var b = a; var c = b; var d = c", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRefChain2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRefChain2() {
	c.Test(
		"var a = 1; var b = a; var c = b; var d = c || doThing1()",
		"var a = 1; var b = a; var c = b;         c || doThing1()")
}

// TODO(b/66971163): Enable or remove
// Java: RemoveUnusedCodeNameAnalyzerTest.disabledTestRefChain3()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestDisabledTestRefChain3() {
	c.T().Skip("@Ignore")
	c.Test(
		"var a = 1; var b = a; var c = b; var d = c + doThing1()",
		"                                             doThing1()")
}

// TODO(b/66971163): Enable or remove
// Java: RemoveUnusedCodeNameAnalyzerTest.disabledTestRefChain4()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestDisabledTestRefChain4() {
	c.T().Skip("@Ignore")
	c.Test(
		"var a = 1; var b = a; var c = b; var d = doThing1() || c",
		"                                         doThing1()     ")
}

// TODO(b/66971163): Enable or remove
// Java: RemoveUnusedCodeNameAnalyzerTest.disabledTestRefChain5()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestDisabledTestRefChain5() {
	c.T().Skip("@Ignore")
	c.Test(
		"var a = 1; var b = a; var c = b; var d = doThing1() ?  doThing2() : c",
		"                                         doThing1() && doThing2()    ")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRefChain6()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRefChain6() {
	c.Test(
		"var a = 1; var b = a; var c = b; var d = c ? doThing1() : doThing2()",
		"var a = 1; var b = a; var c = b;         c ? doThing1() : doThing2()")
}

// TODO(b/66971163): Enable or remove
// Java: RemoveUnusedCodeNameAnalyzerTest.disabledTestRefChain7()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestDisabledTestRefChain7() {
	c.T().Skip("@Ignore")
	c.Test(
		"var a = 1; var b = a; var c = b; var d = (b + doThing1())  ? doThing2() : c",
		"var a = 1; var b = a;                    (b + doThing1()) && doThing2()    ")
}

// TODO(b/66971163): Enable or remove
// Java: RemoveUnusedCodeNameAnalyzerTest.disabledTestRefChain8()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestDisabledTestRefChain8() {
	c.T().Skip("@Ignore")
	c.Test(
		"var a = 1; var b = a; var c = b; var d = doThing1()[b]  ? doThing2() : 0",
		"var a = 1; var b = a;                    doThing1()[b] && doThing2()    ")
}

// TODO(b/66971163): Enable or remove
// Java: RemoveUnusedCodeNameAnalyzerTest.disabledTestRefChain9()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestDisabledTestRefChain9() {
	c.T().Skip("@Ignore")
	c.Test(
		"var a = 1; var b = a; var c = 5; var d = doThing1()[b+c]  ? doThing2() : 0",
		"var a = 1; var b = a; var c = 5;         doThing1()[b+c] && doThing2()    ")
}

// TODO(b/66971163): Enable or remove
// Java: RemoveUnusedCodeNameAnalyzerTest.disabledTestRefChain10()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestDisabledTestRefChain10() {
	c.T().Skip("@Ignore")
	c.Test(
		"var a = 1; var b = a; var c = b; var d = doThing1()[b]  ? doThing2() : 0",
		"var a = 1; var b = a;                    doThing1()[b] && doThing2()    ")
}

// TODO(b/66971163): Enable or remove
// Java: RemoveUnusedCodeNameAnalyzerTest.disabledTestRefChain11()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestDisabledTestRefChain11() {
	c.T().Skip("@Ignore")
	c.Test(
		"var a = 1; var b = a; var d = doThing1()[b]  ? doThing2() : 0",
		"var a = 1; var b = a;         doThing1()[b] && doThing2()    ")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRefChain12()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRefChain12() {
	c.TestSame("var a = 1; var b = a; doThing1()[b] ? doThing2() : 0")
}

// TODO(b/66971163): Enable or remove
// Java: RemoveUnusedCodeNameAnalyzerTest.disabledTestRefChain13()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestDisabledTestRefChain13() {
	c.T().Skip("@Ignore")
	c.Test(
		"function f(){} var a = 1; var b = a; var d = f()[b]  ? doThing1() : 0",
		"function f(){} var a = 1; var b = a;         f()[b] && doThing1()    ")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRefChain14()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRefChain14() {
	c.TestSame("function f(){}var a = 1; var b = a; f()[b] ? doThing1() : 0")
}

// TODO(b/66971163): Enable or remove
// Java: RemoveUnusedCodeNameAnalyzerTest.disabledTestRefChain15()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestDisabledTestRefChain15() {
	c.T().Skip("@Ignore")
	c.Test(
		"function f(){} var a = 1, b = a; var c = f(); var d = c[b]  ? doThing1() : 0",
		"function f(){} var a = 1, b = a; var c = f();         c[b] && doThing1()    ")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRefChain16()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRefChain16() {
	c.TestSame("function f(){}var a = 1; var b = a; var c = f(); c[b] ? doThing1() : 0")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRefChain17()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRefChain17() {
	c.Test(
		"function f(){sideEffect()} var a = 1; var b = a; var c = f(); var d = c[b]",
		"function f(){sideEffect()}                               f()              ")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRefChain18()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRefChain18() {
	c.TestSame("var a = 1; doThing1()[a] && doThing2()")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRefChain19()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRefChain19() {
	c.Test(
		"var a = 1; var b = [a]; var c = b; b[doThing1()] ? doThing2() : 0",
		"var a = 1; var b = [a];            b[doThing1()] ? doThing2() : 0")
}

// TODO(b/66971163): enable or remove this case
// Java: RemoveUnusedCodeNameAnalyzerTest.disabledTestRefChain20()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestDisabledTestRefChain20() {
	c.T().Skip("@Ignore")
	c.Test(
		"var a = 1; var b = [a]; var c = b; var d = b[doThing1()]  ? doThing2() : 0",
		"var a = 1; var b = [a];                    b[doThing1()] && doThing2()    ")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRefChain21()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRefChain21() {
	c.TestSame("var a = 1; var b = 2; var c = a + b; use(c)")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRefChain22()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRefChain22() {
	c.Test(
		"var a = 2; var b = a = 4; use(a)", // preserve newline
		"var a = 2;         a = 4; use(a)")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRefChain23()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRefChain23() {
	c.Test(
		"var a = {}; var b = a[1] || doThing1()", // preserve newline
		"var a = {};         a[1] || doThing1()")
}

/**
 * Expressions that cannot be attributed to any enclosing dependency scope should be treated as
 * global references.
 *
 * @bug 1739062
 */
// Java: RemoveUnusedCodeNameAnalyzerTest.testAssignmentWithComplexLhs()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestAssignmentWithComplexLhs() {
	c.TestSame(
		"function f() { return this; }\n" +
			"var o = {'key': 'val'};\n" +
			"f().x_ = o['key'];\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testAssignmentWithComplexLhs2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestAssignmentWithComplexLhs2() {
	c.TestSame(
		"function f() { return this; }\n" +
			"var o = {'key': 'val'};\n" +
			"f().foo = function() {\n" +
			"  o\n" +
			"};\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testAssignmentWithComplexLhs3()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestAssignmentWithComplexLhs3() {
	source :=
		"var o = {'key': 'val'}; // preserve newline\n" +
			"function init_() {\n" +
			"  use(o['key'])\n" +
			"}\n"

	c.Test(source, "")
	c.TestSame(source + ";init_()")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testAssignmentWithComplexLhs4()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestAssignmentWithComplexLhs4() {
	c.TestSame(
		"function f() { return this; }\n" +
			"var o = {'key': 'val'};\n" +
			"f().foo = function() {\n" +
			"  use(o['key']);\n" +
			"};\n")
}

/**
 * Do not "prototype" property of variables that are not being tracked (because they are local).
 *
 * @bug 1809442
 */
// Java: RemoveUnusedCodeNameAnalyzerTest.testNoRemovePrototypeDefinitionsOutsideGlobalScope1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNoRemovePrototypeDefinitionsOutsideGlobalScope1() {
	c.TestSame(
		"function f(arg){ use(arg); }\n" +
			"(function(){\n" +
			"  function O() {}\n" +
			"  O.prototype = { constructor: O };\n" +
			"  f(O);\n" +
			"})()\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNoRemovePrototypeDefinitionsOutsideGlobalScope2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNoRemovePrototypeDefinitionsOutsideGlobalScope2() {
	c.TestSame(
		"function f(arg){ use(arg); }\n" +
			"(function h(){\n" +
			"  function L() {}\n" +
			"  L.prototype = { constructor: L };\n" +
			"  f(L);\n" +
			"})()\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNoRemovePrototypeDefinitionsOutsideGlobalScope4()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNoRemovePrototypeDefinitionsOutsideGlobalScope4() {
	c.TestSame(
		"function f(arg){ use(arg); }\n" +
			"function g(){\n" +
			"  function N() {}\n" +
			"  N.prototype = { constructor: N };\n" +
			"  f(N);\n" +
			"}\n" +
			"g()\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNoRemovePrototypeDefinitionsOutsideGlobalScope5()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNoRemovePrototypeDefinitionsOutsideGlobalScope5() {
	c.Test(
		"function g(){ function R() {} R.prototype = { constructor: R }; } g()",
		"function g(){                                               } g()")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemovePrototypeDefinitionsInGlobalScope1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemovePrototypeDefinitionsInGlobalScope1() {
	c.TestSame("function M() {} M.prototype = { constructor: M }; use(M);")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemovePrototypeDefinitionsInGlobalScope2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemovePrototypeDefinitionsInGlobalScope2() {
	c.Test("function Q() {} Q.prototype = { constructor: Q };", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveLabeledStatement()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveLabeledStatement() {
	c.Test("LBL: var x = 1;", "LBL: {}")
}

// TODO(b/66971163): enable or remove this
// Java: RemoveUnusedCodeNameAnalyzerTest.disabledTestRemoveLabeledStatement2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestDisabledTestRemoveLabeledStatement2() {
	c.T().Skip("@Ignore")
	c.Test(
		"var x; LBL: x = doThing1() + doThing2()  ", // preserve newline
		"       LBL:   { doThing1() ; doThing2() }")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveLabeledStatement3()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveLabeledStatement3() {
	c.Test("var x; LBL: x = 1;", "LBL: {}")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveLabeledStatement4()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveLabeledStatement4() {
	c.Test(
		"var a; LBL: a = doThing1()", // preserve newline
		"       LBL:     doThing1()")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testPreservePropertyMutationsToAlias1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestPreservePropertyMutationsToAlias1() {
	// Test for issue b/2316773 - property get case
	// Since a is referenced, property mutations via a's alias b must
	// be preserved.
	c.TestSame("var a = {}; var b = a; b.x = 1; a")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testPreservePropertyMutationsToAlias2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestPreservePropertyMutationsToAlias2() {
	// Test for issue b/2316773 - property get case, don't keep 'c'
	c.Test(
		"var a = {}; var b = a; var c = a; b.x = 1; a", // preserve newline
		"var a = {}; var b = a;            b.x = 1; a")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testPreservePropertyMutationsToAlias3()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestPreservePropertyMutationsToAlias3() {
	// Test for issue b/2316773 - property get case, chain
	c.TestSame("var a = {}; var b = a; var c = b; c.x = 1; a")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testPreservePropertyMutationsToAlias4()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestPreservePropertyMutationsToAlias4() {
	// Test for issue b/2316773 - element get case
	c.TestSame("var a = {}; var b = a; b['x'] = 1; a")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testPreservePropertyMutationsToAlias5()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestPreservePropertyMutationsToAlias5() {
	// From issue b/2316773 description
	c.TestSame(
		"function testCall(o){ use(o); }\n" +
			"var DATA = {'prop': 'foo','attr': {}};\n" +
			"var SUBDATA = DATA['attr'];\n" +
			"SUBDATA['subprop'] = 'bar';\n" +
			"testCall(DATA);\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testPreservePropertyMutationsToAlias6()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestPreservePropertyMutationsToAlias6() {
	// Longer GETELEM chain
	c.TestSame(
		"function testCall(o){ use(o); }\n" +
			"var DATA = {'prop': 'foo','attr': {}};\n" +
			"var SUBDATA = DATA['attr'];\n" +
			"var SUBSUBDATA = SUBDATA['subprop'];\n" +
			"SUBSUBDATA['subsubprop'] = 'bar';\n" +
			"testCall(DATA);\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testPreservePropertyMutationsToAlias7()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestPreservePropertyMutationsToAlias7() {
	// Make sure that the base class does not depend on the derived class.
	c.Test(
		"var a = {}; var b = {}; b.x = 0; goog.inherits(b, a); use(a);",
		"var a = {};                                           use(a);")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testPreservePropertyMutationsToAlias8()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestPreservePropertyMutationsToAlias8() {
	// Make sure that the derived classes don't end up depending on each other.
	c.Test(
		"var a = {}; // preserve newline\n"+
			"var b = {}; b.x = 0;\n"+
			"var c = {}; c.y = 0;\n"+
			"goog.inherits(b, a);\n"+
			"goog.inherits(c, a);\n"+
			"c\n",
		"var a = {}; // preserve newline\n"+
			"\n"+
			"var c = {}; c.y = 0;\n"+
			"\n"+
			"goog.inherits(c, a);\n"+
			"c\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testPreservePropertyMutationsToAlias9()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestPreservePropertyMutationsToAlias9() {
	c.TestSame(
		"var a = {b: {}}; // preserve newline\n" +
			"var c = a.b; c.d = 3;\n" +
			"a.d = 3; a.d;\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemoveAlias()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemoveAlias() {
	c.Test(
		"var a = {b: {}}; var c = a.b; a.d = 3; a.d;",
		"var a = {b: {}};              a.d = 3; a.d;")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSingletonGetter1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSingletonGetter1() {
	c.Test("function Foo() {} goog.addSingletonGetter(Foo);", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSingletonGetter2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSingletonGetter2() {
	c.Test("function Foo() {} goog$addSingletonGetter(Foo);", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testSingletonGetter3()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSingletonGetter3() {
	// addSingletonGetter adds a getInstance method to a class.
	c.TestSame("function Foo() {} goog$addSingletonGetter(Foo); Foo.getInstance();")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testObjectDefineProperty()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestObjectDefineProperty() {
	// TODO(bradfordcsmith): Remove Object.defineProperty() like we do Object.defineProperties().
	c.TestSame("var a = {}; Object.defineProperty(a, 'prop', {value: 5});")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testObjectDefinePropertiesOnNamespaceThatEscapes()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestObjectDefinePropertiesOnNamespaceThatEscapes() {
	c.TestSame("var a = doThing1(); Object.defineProperties(a, {'prop': {value: 5}});")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testObjectDefinePropertiesOnConstructorThatEscapes()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestObjectDefinePropertiesOnConstructorThatEscapes() {
	c.TestSame("var Foo = doThing1(); Object.defineProperties(Foo.prototype, {'prop': {value: 5}});")
	c.TestSame("Object.defineProperties(doThing1(), {'prop': {value: 5}});")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRegularAssignPropOnNamespaceThatEscapes()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRegularAssignPropOnNamespaceThatEscapes() {
	c.TestSame("var a = doThing1(); a.prop = 5;")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRegularAssignPropOnPropFromAVar()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRegularAssignPropOnPropFromAVar() {
	c.TestSame("var b = 5; var a = {}; a.prop = b; use(a.prop);")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testUnanalyzableObjectDefineProperties()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestUnanalyzableObjectDefineProperties() {
	c.Test("var a = {}; Object.defineProperties(a, externfoo);", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testObjectDefinePropertiesOnNamespace1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestObjectDefinePropertiesOnNamespace1() {
	c.TestSame("var a = {}; Object.defineProperties(a, {prop: {value: 5}}); use(a.prop);")
	c.Test("var a = {}; Object.defineProperties(a, {prop: {value: 5}});", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testObjectDefinePropertiesOnNamespace2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestObjectDefinePropertiesOnNamespace2() {
	c.Test(
		"var a = {};\n"+
			"Object.defineProperties(a, {p1: {value: 5}, p2: {value: 3} });\n"+
			"use(a.p1);\n",
		"var a = {};\n"+
			"Object.defineProperties(a, {p1: {value: 5}                 });\n"+
			"use(a.p1);\n")

	c.Test(
		"var a = {}; // preserve newline\n"+
			"Object.defineProperties(a, {p1: {value: 5}, p2: {value: 3} });\n",
		"")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNonAnalyzableObjectDefinePropertiesCall()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNonAnalyzableObjectDefinePropertiesCall() {
	c.TestSame("var a = {}; var z = Object.defineProperties(a, {'prop': {value: 5}}); use(z);")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testObjectDefinePropertiesOnNamespace3()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestObjectDefinePropertiesOnNamespace3() {
	c.TestSame(
		"var b = 5;\n" +
			"var a = {};\n" +
			"Object.defineProperties(a, {prop: {value: b}});\n" +
			"use(a.prop);\n")

	c.Test(
		"var b = 5;\n"+
			"var a = {};\n"+
			"Object.defineProperties(a, {prop: {value: b}});\n"+
			"use(b);\n",
		"var b = 5; use(b);")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testObjectDefinePropertiesOnNamespace4()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestObjectDefinePropertiesOnNamespace4() {
	c.Test(
		"function b() { alert('hello'); };\n"+
			"var a = {};\n"+
			"Object.defineProperties(a, {prop: {value: b()}});\n",
		"function b() { alert('hello'); }; ({prop: {value: b()}});")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testObjectDefinePropertiesOnNamespace5()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestObjectDefinePropertiesOnNamespace5() {
	c.Test(
		"function b() { alert('hello'); }; // preserve newline\n"+
			"function c() { alert('world'); };\n"+
			"var a = {};\n"+
			"Object.defineProperties(a, {p1: {value: b()}, p2: {value: c()}});\n",
		"function b() { alert('hello'); }; // preserve newline\n"+
			"function c() { alert('world'); };\n"+
			"\n"+
			"                          ({p1: {value: b()}, p2: {value: c()}});\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testObjectDefinePropertiesOnConstructor()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestObjectDefinePropertiesOnConstructor() {
	c.TestSame("function Foo() {} Object.defineProperties(Foo, {prop: {value: 5}}); use(Foo.prop);")
	c.Test("function Foo() {} Object.defineProperties(Foo, {prop: {value: 5}});", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testObjectDefinePropertiesOnPrototype1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestObjectDefinePropertiesOnPrototype1() {
	c.TestSame(
		"function Foo() {}\n" +
			"Object.defineProperties(Foo.prototype, {prop: {value: 5}});\n" +
			"use((new Foo).prop);\n")

	c.Test("function Foo() {} Object.defineProperties(Foo.prototype, {prop: {value: 5}});", "")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testObjectDefinePropertiesOnPrototype2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestObjectDefinePropertiesOnPrototype2() {
	c.Test(
		"var b = 5;\n"+
			"function Foo() {}\n"+
			"Object.defineProperties(Foo.prototype, {prop: {value: b}});\n"+
			"use(b)\n",
		"var b = 5; use(b);")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testObjectDefinePropertiesOnPrototype3()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestObjectDefinePropertiesOnPrototype3() {
	c.Test(
		"var b = function() {sideEffect()};\n"+
			"function Foo() {}\n"+
			"Object.defineProperties(Foo.prototype, {prop: {value: b()}});\n",
		"var b = function() {sideEffect()}; ({prop: {value: b()}});")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testObjectDefineGetters()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestObjectDefineGetters() {
	c.Test("function Foo() {} Object.defineProperties(Foo, {prop: {get: function() {}}});", "")

	c.Test(
		"function Foo() {} Object.defineProperties(Foo.prototype, {prop: {get: function() {}}});",
		"")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testObjectDefineSetters()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestObjectDefineSetters() {
	c.Test("function Foo() {} Object.defineProperties(Foo, {prop: {set: function() {}}});", "")

	c.Test(
		"function Foo() {} Object.defineProperties(Foo.prototype, {prop: {set: function() {}}});",
		"")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testObjectDefineSetters_global()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestObjectDefineSetters_global() {
	c.Test(
		"function Foo() {}\n"+
			"$jscomp.global.Object.defineProperties(Foo, {prop: {set: function() {}}});\n",
		"")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNoRemoveWindowPropertyAlias1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNoRemoveWindowPropertyAlias1() {
	c.TestSame("var self_ = window.gbar; self_.qs = function() {};")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNoRemoveWindowPropertyAlias2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNoRemoveWindowPropertyAlias2() {
	c.TestSame("var self_ = window; self_.qs = function() {};")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNoRemoveWindowPropertyAlias3()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNoRemoveWindowPropertyAlias3() {
	c.TestSame("var self_ = window; self_['qs'] = function() {};")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNoRemoveWindowPropertyAlias4()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNoRemoveWindowPropertyAlias4() {
	c.TestSame("var self_ = window['gbar'] || {}; self_.qs = function() {};")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNoRemoveWindowPropertyAlias4a()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNoRemoveWindowPropertyAlias4a() {
	c.TestSame("var self_; self_ = window.gbar || {}; self_.qs = function() {};")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNoRemoveWindowPropertyAlias5()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNoRemoveWindowPropertyAlias5() {
	c.TestSame("var self_ = window || {}; self_['qs'] = function() {};")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNoRemoveWindowPropertyAlias5a()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNoRemoveWindowPropertyAlias5a() {
	c.TestSame("var self_; self_ = window || {}; self_['qs'] = function() {};")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNoRemoveWindowPropertyAlias6()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNoRemoveWindowPropertyAlias6() {
	c.TestSame("var self_ = (window.gbar = window.gbar || {}); self_.qs = function() {};")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNoRemoveWindowPropertyAlias6a()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNoRemoveWindowPropertyAlias6a() {
	c.TestSame("var self_; self_ = (window.gbar = window.gbar || {}); self_.qs = function() {};")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNoRemoveWindowPropertyAlias7()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNoRemoveWindowPropertyAlias7() {
	c.TestSame("var self_ = (window = window || {}); self_['qs'] = function() {};")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNoRemoveWindowPropertyAlias7a()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNoRemoveWindowPropertyAlias7a() {
	c.TestSame("var self_; self_ = (window = window || {}); self_['qs'] = function() {};")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNoRemoveAlias0()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNoRemoveAlias0() {
	c.TestSame(
		"var x = {}; function f() { return x; };\n" +
			"f().style.display = 'block';\n" +
			"alert(x.style)\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNoRemoveAlias1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNoRemoveAlias1() {
	c.TestSame(
		"var x = {}; function f() { return x; };\n" +
			"var map = f();\n" +
			"map.style.display = 'block';\n" +
			"alert(x.style)\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNoRemoveAlias2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNoRemoveAlias2() {
	c.TestSame(
		"var x = {};\n" +
			"var map = (function () { return x; })();\n" +
			"map.style = 'block';\n" +
			"alert(x.style)\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNoRemoveAlias3()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNoRemoveAlias3() {
	c.TestSame(
		"var x = {}; function f() { return x; };\n" +
			"var map = {};\n" +
			"map[1] = f();\n" +
			"map[1].style.display = 'block';\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNoRemoveAliasOfExternal0()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNoRemoveAliasOfExternal0() {
	c.TestSame("document.getElementById('foo').style.display = 'block';")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNoRemoveAliasOfExternal1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNoRemoveAliasOfExternal1() {
	c.TestSame("var map = document.getElementById('foo'); map.style.display = 'block';")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNoRemoveAliasOfExternal2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNoRemoveAliasOfExternal2() {
	c.TestSame(
		"var map = {} // preserve newline\n" +
			"map[1] = document.getElementById('foo');\n" +
			"map[1].style.display = 'block';\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNoRemoveThrowReference1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNoRemoveThrowReference1() {
	c.TestSame("var e = {}; throw e;")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testNoRemoveThrowReference2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestNoRemoveThrowReference2() {
	c.TestSame("function e() {} throw new e();")
}

// TODO(b/66971163): enable or remove this
// Java: RemoveUnusedCodeNameAnalyzerTest.disabledTestClassDefinedInObjectLit1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestDisabledTestClassDefinedInObjectLit1() {
	c.T().Skip("@Ignore")
	c.Test(
		"var data = {Foo: function() {}}; // preserve newline\n"+
			"data.Foo.prototype.toString = function() {};\n",
		"")
}

// TODO(b/66971163): enable or remove this
// Java: RemoveUnusedCodeNameAnalyzerTest.disabledTestClassDefinedInObjectLit2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestDisabledTestClassDefinedInObjectLit2() {
	c.T().Skip("@Ignore")
	c.Test(
		"var data = {}; data.bar = {Foo: function() {}};\n"+
			"data.bar.Foo.prototype.toString = function() {};\n",
		"")
}

// TODO(b/66971163): enable or remove this
// Java: RemoveUnusedCodeNameAnalyzerTest.disabledTestClassDefinedInObjectLit3()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestDisabledTestClassDefinedInObjectLit3() {
	c.T().Skip("@Ignore")
	c.Test(
		"var data = {bar: {Foo: function() {}}}; // preserve newline\n"+
			"data.bar.Foo.prototype.toString = function() {};\n",
		"")
}

// TODO(b/66971163): enable or remove this
// Java: RemoveUnusedCodeNameAnalyzerTest.disabledTestClassDefinedInObjectLit4()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestDisabledTestClassDefinedInObjectLit4() {
	c.T().Skip("@Ignore")
	c.Test(
		"var data = {};\n"+
			"data.baz = {bar: {Foo: function() {}}};\n"+
			"data.baz.bar.Foo.prototype.toString = function() {};\n",
		"")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testVarReferencedInClassDefinedInObjectLit1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestVarReferencedInClassDefinedInObjectLit1() {
	c.TestSame(
		"var ref = 3; // preserve newline\n" +
			"var data = {Foo: function() { use(ref); }};\n" +
			"window.Foo = data.Foo;\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testVarReferencedInClassDefinedInObjectLit2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestVarReferencedInClassDefinedInObjectLit2() {
	c.TestSame(
		"var ref = 3;\n" +
			"var data = {\n" +
			"  Foo: function() { use(ref); },\n" +
			"  Bar: function() {}\n" +
			"};\n" +
			"window.Bar = data.Bar;\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testArrayExt()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestArrayExt() {
	c.TestSame(
		"Array.prototype.foo = function() { return 1 };\n" +
			"var y = [];\n" +
			"switch (y.foo()) {\n" +
			"}\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testArrayAliasExt()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestArrayAliasExt() {
	c.TestSame(
		"Array$X = Array;\n" +
			"Array$X.prototype.foo = function() { return 1 };\n" +
			"function Array$X() {}\n" +
			"var y = [];\n" +
			"switch (y.foo()) {\n" +
			"}\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testExternalAliasInstanceof1()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestExternalAliasInstanceof1() {
	c.Test(
		"Array$X = Array; // preserve newline\n"+
			"function Array$X() {}\n"+
			"var y = [];\n"+
			"if (y instanceof Array) {}\n",
		"var y = []; // preserve newline\n"+
			"if (y instanceof Array) {}\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testExternalAliasInstanceof2()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestExternalAliasInstanceof2() {
	c.TestSame(
		"Array$X = Array;\n" +
			"function Array$X() {}\n" +
			"var y = [];\n" +
			"if (y instanceof Array$X) {}\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testExternalAliasInstanceof3()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestExternalAliasInstanceof3() {
	c.TestSame("var b = Array; var y = []; if (y instanceof b) {}")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testAliasInstanceof4()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestAliasInstanceof4() {
	c.TestSame("function Foo() {}; var b = Foo; var y = new Foo(); if (y instanceof b) {}")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testAliasInstanceof5()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestAliasInstanceof5() {
	c.TestSame(
		"var x;\n" +
			"function Foo() {}\n" +
			"function Bar() {}\n" +
			"var b = x ? Foo : Bar;\n" +
			"var y = new Foo();\n" +
			"if (y instanceof b) {}\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testRemovePrototypeAliases()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestRemovePrototypeAliases() {
	c.Test(
		"function g() {} function F() {} F.prototype.bar = g; window.g = g;",
		"function g() {}                                      window.g = g;")
}

// TODO(b/66971163): Enable or remove this test.
// Java: RemoveUnusedCodeNameAnalyzerTest.disabledTestIssue284()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestDisabledTestIssue284() {
	c.T().Skip("@Ignore")
	c.Test(
		"var ns = {};\n"+
			"/** @constructor */\n"+
			"ns.PageSelectionModel = function() {};\n"+
			"/** @constructor */\n"+
			"ns.PageSelectionModel.FooEvent = function() {};\n"+
			"/** @constructor */\n"+
			"ns.PageSelectionModel.SelectEvent = function() {};\n"+
			"goog.inherits(ns.PageSelectionModel.ChangeEvent, ns.PageSelectionModel.FooEvent);\n",
		"")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testIssue838a()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestIssue838a() {
	c.TestSame(
		"var z = window['z'] || (window['z'] = {}); // preserve newline\n" +
			"z['hello'] = 'Hello';\n" +
			"z['world'] = 'World';\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testIssue838b()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestIssue838b() {
	c.TestSame(
		"var z; // preserve newline\n" +
			"window['z'] = z || (z = {});\n" +
			"z['hello'] = 'Hello';\n" +
			"z['world'] = 'World';\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testIssue874a()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestIssue874a() {
	c.TestSame(
		"var a = a || {};\n" +
			"var b = a;\n" +
			"b.View = b.View || {}\n" +
			"var c = b.View;\n" +
			"c.Editor = function f(d, e) {\n" +
			"  return d + e\n" +
			"};\n" +
			"window.ImageEditor.View.Editor = a.View.Editor;\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testIssue874b()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestIssue874b() {
	c.TestSame(
		"var b;\n" +
			"var c = b = {};\n" +
			"c.Editor = function f(d, e) {\n" +
			"  return d + e\n" +
			"};\n" +
			"window['Editor'] = b.Editor;\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testIssue874c()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestIssue874c() {
	c.TestSame(
		"var b, c;\n" +
			"c = b = {};\n" +
			"c.Editor = function f(d, e) {\n" +
			"  return d + e\n" +
			"};\n" +
			"window['Editor'] = b.Editor;\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testIssue874d()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestIssue874d() {
	c.TestSame(
		"var b = {}, c;\n" +
			"c = b;\n" +
			"c.Editor = function f(d, e) {\n" +
			"  return d + e\n" +
			"};\n" +
			"window['Editor'] = b.Editor;\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testIssue874e()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestIssue874e() {
	c.TestSame(
		"var a;\n" +
			"var b = a || (a = {});\n" +
			"var c = b.View || (b.View = {});\n" +
			"c.Editor = function f(d, e) {\n" +
			"  return d + e\n" +
			"};\n" +
			"window.ImageEditor.View.Editor = a.View.Editor;\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testBug6575051()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestBug6575051() {
	c.TestSame(
		"var hackhack = window['__o_o_o__'] = window['__o_o_o__'] || {};\n" +
			"window['__o_o_o__']['va'] = 1;\n" +
			"hackhack['Vb'] = 1;\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testBug37975351a()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestBug37975351a() {
	// The original repro case from the bug.
	c.TestSame(
		"function noop() {}\n" +
			"var x = window['magic'];\n" +
			"var FormData = window['FormData'] || noop;\n" +
			"function f() { return x instanceof FormData; }\n" +
			"console.log(f());\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testBug37975351b()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestBug37975351b() {
	// The simplified repro that still repro'd the problem.
	c.TestSame(
		"var FormData = window['FormData'] || function() {};\n" +
			"function f() { return window['magic'] instanceof FormData; }\n" +
			"console.log(f());\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testBug37975351c()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestBug37975351c() {
	// This simpliification did not reproduce the problematic behavior.
	c.TestSame(
		"var FormData = window['FormData'];\n" +
			"function f() { return window['magic'] instanceof FormData; }\n" +
			"console.log(f());\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testBug30868041()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestBug30868041() {
	c.TestSame(
		"function Base() {};\n" +
			"/** @nosideeffects */\n" +
			"Base.prototype.foo =  function() {\n" +
			"}\n" +
			"var x = new Base();\n" +
			"x.foo()\n")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testGenerators()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestGenerators() {
	c.Test("function* g() {yield 1}", "")

	c.TestSame("function* g() {yield 1} var g = g(); g.next().value()")
}

/**
 * Just check that we don't crash in this case.
 *
 * @bug 65489464
 */
// Java: RemoveUnusedCodeNameAnalyzerTest.testSpread()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestSpread() {
	c.Test(
		"const ns = {};\n"+
			"const X = [];\n"+
			"ns.Y = [{}, ...X];\n",
		"const X = [];\n"+
			"       [{}, ...X];\n")
}

// TODO(b/66971163): enable or remove this test
// Java: RemoveUnusedCodeNameAnalyzerTest.disabledTestObjectDestructuring()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestDisabledTestObjectDestructuring() {
	c.T().Skip("@Ignore")
	c.Test(
		"var {a: a, x: y} = {a:1, x:2} ", // preserve newline
		"                   ({a:1,x:2})")

	c.Test(
		"var {a: a, x: y} = {a:1, x:2}; f(a);", // preserve newline
		"var {a: a      } = {a:1, x:2}; f(a);")

	c.Test(
		"var {a: a, x: y} = {a:1, x:2}; f(y);", // preserve newline
		"var {      x: y} = {a:1, x:2}; f(y);")

	c.Test(
		"var {a: a, x: y = 3} = {a:1, x:2}; f(y);", // preserve newline
		"var {      x: y = 3} = {a:1, x:2}; f(y);")

	c.Test(
		"var {a: a, x: y = 3} = {a:1}; f(y);", // preserve newline
		"var {      x: y = 3} = {a:1}; f(y);")

	c.Test(
		"function f() {} var {a: a, x: y} = f();", // preserve newline
		"function f() {}                    f();")

	c.Test(
		"function f() {} var {a: a, x: y} = f(); g(a)", // preserve newline
		"function f() {} var {a: a      } = f(); g(a)")

	c.Test(
		"function f() {} var {a: a, x: y} = f(); g(y)", // preserve newline
		"function f() {} var {      x: y} = f(); g(y)")

	// complicated destructuring cases
	// TODO(blickly): Look into adding a pass to completely remove empty destructuring patterns
	c.Test(
		"var {a: a, b: [{c: d}]} = o;", // preserve newline
		"var {      b: [{    }]} = o; ")

	c.Test(
		"var {a: a, b: {c: d}} = o; f(d)", // preserve newline
		"var {      b: {c: d}} = o; f(d)")

	c.Test(
		"var {a: a, b: [key]} = o;", // preserve newline
		"var {      b: [key]} = o;")

	c.Test(
		"var {a: a, [key]: foo} = o;", // preserve newline
		"var {      [key]: foo} = o;")

	c.Test(
		"var { a: a, b: { c: { d: y}}} = o", // preserve newline
		"var {       b: { c: {     }}} = o")

	c.Test(
		"var {[foo()] : { p : x } } = o;", // preserve newline
		"var {[foo()] : {       } } = o;")

	c.TestSame("var {x = foo()} = o;")

	c.TestSame("var {p : x = foo()} = o;")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testArrayDestructuring()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestArrayDestructuring() {
	c.Test(
		"var [a, b = 3, ...c] = [1, 2, 3]", // preserve newline
		"var [ ,      , ...c] = [1, 2, 3]")

	c.Test(
		"var [a, b = 3, ...c] = [1, 2, 3]; use(b);", // preserve newline
		"var [ , b = 3  ...c] = [1, 2, 3]; use(b);")

	c.Test(
		"var [a, b = 3, ...c] = [1, 2, 3]; use(c);", // preserve newline
		"var [ ,      , ...c] = [1, 2, 3]; use(c);")

	c.Test(
		"var a, b, c; [a, b, ...c] = [1, 2, 3]", // preserve newline
		"var c      ; [ ,  , ...c] = [1, 2, 3]")

	c.Test(
		"var [a, [b, [c, d]]] = [1, [2, [[[3, 4], 5], 6]]];", // preserve newline
		"var [ , [ , [    ]]] = [1, [2, [[[3, 4], 5], 6]]];")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testBlock()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestBlock() {
	// Currently after normalization this becomes {var f = function f() {}}
	// Will no longer be able to be removed after that normalize change
	c.Test("{function g() {}}", "{}")

	c.TestSame("{function g() {} g()}")

	c.TestSame("function g() {} {let a = g(); use(a)}")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testTemplateLit()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestTemplateLit() {
	c.Test("let a = `hello`", "")
	c.Test("var name = 'foo'; let a = `hello ${name}`", "")
	c.Test(
		"function Base() {} // preserve newline\n"+
			"Base.prototype.foo =  `hello`;\n",
		"")

	c.Test(
		"var bar = 'foo';\n"+
			"function Base() {} // preserve newline\n"+
			"Base.prototype.foo =  `foo ${bar}`;\n",
		"")
}

// Java: RemoveUnusedCodeNameAnalyzerTest.testPureOrBreakMyCode()
func (c *RemoveUnusedCodeNameAnalyzerTest) TestPureOrBreakMyCode() {
	c.Test("const a = /** @pureOrBreakMyCode */(alert());", "")
	c.Test("let a = /** @pureOrBreakMyCode */(alert());", "")
	c.Test("var a = /** @pureOrBreakMyCode */(alert());", "")
}
