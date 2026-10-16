// Port of com/google/javascript/jscomp/PeepholeSubstituteAlternateSyntaxTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/PeepholeSubstituteAlternateSyntaxTest.java)

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

import "testing"

// PeepholeSubstituteAlternateSyntaxTest is the tests for
// PeepholeSubstituteAlternateSyntax in isolation. Tests for the interaction
// of multiple peephole passes are in PeepholeIntegrationTest.
type PeepholeSubstituteAlternateSyntaxTest struct {
	CompilerTestCaseBase

	late               bool
	retraverseOnChange bool
}

// Externs for built-in constructors
// Needed for testFoldLiteralObjectConstructors(),
// testFoldLiteralArrayConstructors() and testFoldRegExp...()
var peepholeSubstituteAlternateSyntaxTest_FOLD_CONSTANTS_TEST_EXTERNS = "var window = {};\n" +
	"var Object = function f(){};\n" +
	"var RegExp = function f(a){};\n" +
	"var Array = function f(a){};\n" +
	"window.foo = null;\n"

// Java: PeepholeSubstituteAlternateSyntaxTest()
func NewPeepholeSubstituteAlternateSyntaxTest() *PeepholeSubstituteAlternateSyntaxTest {
	c := &PeepholeSubstituteAlternateSyntaxTest{}
	c.InitCompilerTestCaseBaseWithExterns(c, peepholeSubstituteAlternateSyntaxTest_FOLD_CONSTANTS_TEST_EXTERNS)
	return c
}

func TestPeepholeSubstituteAlternateSyntaxTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewPeepholeSubstituteAlternateSyntaxTest() })
}

// Java: PeepholeSubstituteAlternateSyntaxTest.setUp()
func (c *PeepholeSubstituteAlternateSyntaxTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.late = true
	c.retraverseOnChange = false
	c.DisableNormalize()
}

// Java: PeepholeSubstituteAlternateSyntaxTest.getProcessor(Compiler)
func (c *PeepholeSubstituteAlternateSyntaxTest) GetProcessor(compiler *Compiler) CompilerPass {
	peepholePass := NewPeepholeOptimizationsPass(
		compiler, c.GetName(), NewPeepholeSubstituteAlternateSyntax(c.late))
	peepholePass.SetRetraverseOnChange(c.retraverseOnChange)
	return peepholePass
}

// Java: PeepholeSubstituteAlternateSyntaxTest.testFoldRegExpConstructor()
func (c *PeepholeSubstituteAlternateSyntaxTest) TestFoldRegExpConstructor() {
	c.EnableNormalize()
	// Cannot fold all the way to a literal because there are too few arguments.
	c.Test("x = new RegExp", "x = RegExp()")
	// Could fold all the way to /foobar/, but the complexity of that optimization wasn't worth the
	// code size improvements, mainly because there are a lot special cases to check.
	c.Test("x = new RegExp(\"foobar\")", "x = RegExp(\"foobar\")")
}

// Java: PeepholeSubstituteAlternateSyntaxTest.testFoldLiteralObjectConstructors()
func (c *PeepholeSubstituteAlternateSyntaxTest) TestFoldLiteralObjectConstructors() {
	c.EnableNormalize()
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	// Can fold when normalized
	c.Test("x = new Object", "x = ({})")
	c.Test("x = new Object()", "x = ({})")
	c.Test("x = Object()", "x = ({})")
	c.DisableNormalize()
	// Cannot fold above when not normalized
	c.TestSame("x = new Object")
	c.TestSame("x = new Object()")
	c.TestSame("x = Object()")
	c.EnableNormalize()
	// Cannot fold, the constructor being used is actually a local function
	c.TestSame(
		"x =\n" +
			"(function f(){function Object(){this.x=4};return new Object();})();\n")
}

// Java: PeepholeSubstituteAlternateSyntaxTest.testFoldLiteralObjectConstructors_onWindow()
func (c *PeepholeSubstituteAlternateSyntaxTest) TestFoldLiteralObjectConstructors_onWindow() {
	c.EnableNormalize()
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	// Can fold when normalized
	c.Test("x = new window.Object", "x = ({})")
	c.Test("x = new window.Object()", "x = ({})")
	// Mustn't fold optional chains
	c.Test("x = window.Object()", "x = ({})")
	c.Test("x = window.Object?.()", "x = Object?.()")
	c.DisableNormalize()
	// Cannot fold above when not normalized
	c.TestSame("x = new window.Object")
	c.TestSame("x = new window.Object()")
	c.TestSame("x = window.Object()")
	c.TestSame("x = window?.Object()")
	c.EnableNormalize()
	// Can fold, the window namespace ensures it's not a conflict with the local Object.
	c.Test(
		"x = (function f(){function Object(){this.x=4};return new window.Object;})();",
		"x = (function f(){function Object(){this.x=4};return {};})();")
}

// Java: PeepholeSubstituteAlternateSyntaxTest.testFoldLiteralArrayConstructors()
func (c *PeepholeSubstituteAlternateSyntaxTest) TestFoldLiteralArrayConstructors() {
	c.EnableNormalize()
	// No arguments - can fold when normalized
	c.Test("x = new Array", "x = []")
	c.Test("x = new Array()", "x = []")
	c.Test("x = Array()", "x = []")
	c.TestSame("x = Array?.()")
	// Mustn't fold optional chains
	// One argument - can be fold when normalized
	c.Test("x = new Array(0)", "x = []")
	c.Test("x = Array(0)", "x = []")
	c.Test("x = new Array(\"a\")", "x = [\"a\"]")
	c.Test("x = Array(\"a\")", "x = [\"a\"]")
	// One argument - cannot be fold when normalized
	c.Test("x = new Array(7)", "x = Array(7)")
	c.TestSame("x = Array(7)")
	c.Test("x = new Array(y)", "x = Array(y)")
	c.TestSame("x = Array(y)")
	c.Test("x = new Array(foo())", "x = Array(foo())")
	c.TestSame("x = Array(foo())")
	// More than one argument - can be fold when normalized
	c.Test("x = new Array(1, 2, 3, 4)", "x = [1, 2, 3, 4]")
	c.Test("x = Array(1, 2, 3, 4)", "x = [1, 2, 3, 4]")
	c.Test("x = new Array('a', 1, 2, 'bc', 3, {}, 'abc')", "x = ['a', 1, 2, 'bc', 3, {}, 'abc']")
	c.Test("x = Array('a', 1, 2, 'bc', 3, {}, 'abc')", "x = ['a', 1, 2, 'bc', 3, {}, 'abc']")
	c.Test("x = new Array(Array(1, '2', 3, '4'))", "x = [[1, '2', 3, '4']]")
	c.Test("x = Array(Array(1, '2', 3, '4'))", "x = [[1, '2', 3, '4']]")
	c.Test(
		"x = new Array(Object(), Array(\"abc\", Object(), Array(Array())))",
		"x = [{}, [\"abc\", {}, [[]]]]")
	c.Test(
		"x = new Array(Object(), Array(\"abc\", Object(), Array(Array())))",
		"x = [{}, [\"abc\", {}, [[]]]]")
	c.DisableNormalize()
	// Cannot fold above when not normalized
	c.TestSame("x = new Array")
	c.TestSame("x = new Array()")
	c.TestSame("x = Array()")
	c.TestSame("x = new Array(0)")
	c.TestSame("x = Array(0)")
	c.TestSame("x = new Array(\"a\")")
	c.TestSame("x = Array(\"a\")")
	c.TestSame("x = new Array(7)")
	c.TestSame("x = Array(7)")
	c.TestSame("x = new Array(foo())")
	c.TestSame("x = Array(foo())")
	c.TestSame("x = new Array(1, 2, 3, 4)")
	c.TestSame("x = Array(1, 2, 3, 4)")
	c.TestSame("x = new Array('a', 1, 2, 'bc', 3, {}, 'abc')")
	c.TestSame("x = Array('a', 1, 2, 'bc', 3, {}, 'abc')")
	c.TestSame("x = new Array(Array(1, '2', 3, '4'))")
	c.TestSame("x = Array(Array(1, '2', 3, '4'))")
	c.TestSame(
		"x = new Array(\n" +
			"Object(), Array(\"abc\", Object(), Array(Array())))\n")
	c.TestSame(
		"x = new Array(\n" +
			"Object(), Array(\"abc\", Object(), Array(Array())))\n")
}

// Java: PeepholeSubstituteAlternateSyntaxTest.testRemoveWindowRefs()
func (c *PeepholeSubstituteAlternateSyntaxTest) TestRemoveWindowRefs() {
	c.EnableNormalize()
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.Test("x = window.Object", "x = Object")
	c.Test("x = window.Object.keys", "x = Object.keys")
	c.Test("if (window.Object) {}", "if (Object) {}")
	c.Test("x = window.Object", "x = Object")
	c.Test("x = window.Array", "x = Array")
	c.Test("x = window.Error", "x = Error")
	c.Test("x = window.RegExp", "x = RegExp")
	c.Test("x = window.Math", "x = Math")
	// Not currently handled by the pass but should be folded in the future.
	c.TestSame("x = window.String")
	// Don't fold properties on the window.
	c.TestSame("x = window.foo")
	c.DisableNormalize()
	// Cannot fold when not normalized
	c.TestSame("x = window.Object")
	c.TestSame("x = window.Object.keys")
	c.EnableNormalize()
	c.TestSame(
		"var x =\n" +
			"(function f(){var window = {Object: function() {}};return new window.Object;})();\n")
}

// Tests that it's safe to fold access on window that is non-optional but is under an optional
// chain in the AST.
// Java: PeepholeSubstituteAlternateSyntaxTest.testRemoveWindowRef_childOfOptionalChain()
func (c *PeepholeSubstituteAlternateSyntaxTest) TestRemoveWindowRef_childOfOptionalChain() {
	c.EnableNormalize()
	// ref on window can be folded preserving the optional access.
	c.Test("x = window.Object?.keys", "x = Object?.keys")
	c.Test("x = window.Object?.(keys)", "x = Object?.(keys)")
	// Show that the above works only on `BUILTIN_EXTERNS` such as Object but not on regular prop
	// accesses
	c.TestSame("x = window.prop?.keys")
	c.TestSame("x = window.prop?.(keys)")
}

// There isn't an obvious reason to write `window?.Boolean()` or `window?.Error()`, but if one
// did, presumably there's some reason to believe that `window` won't be there. We don't optimize
// away these checks.
// Java: PeepholeSubstituteAlternateSyntaxTest.testDontRemoveWindowRefs_optChain()
func (c *PeepholeSubstituteAlternateSyntaxTest) TestDontRemoveWindowRefs_optChain() {
	c.EnableNormalize()
	// Don't fold the optional chain check on window
	c.TestSame("x = window?.Object")
	c.TestSame("if (window?.Object) {}")
	c.TestSame("x = window?.Object")
	c.TestSame("x = window?.Array")
	c.TestSame("x = window?.Error")
	c.TestSame("x = window?.RegExp")
	c.TestSame("x = window?.Math")
	c.TestSame("x = window?.String")
	// Don't fold properties on the window anyway (non-optional or optional).
	c.TestSame("x = window.foo")
	c.TestSame("x = window?.foo")
	c.DisableNormalize()
	// Cannot fold when not normalized
	c.TestSame("x = window?.Object")
	c.TestSame("x = window.Object?.keys")
}

// Java: PeepholeSubstituteAlternateSyntaxTest.testFoldStandardConstructors()
func (c *PeepholeSubstituteAlternateSyntaxTest) TestFoldStandardConstructors() {
	c.TestSame("new Foo('a')")
	c.TestSame("var x = new goog.Foo(1)")
	c.TestSame("var x = new String(1)")
	c.TestSame("var x = new Number(1)")
	c.TestSame("var x = new Boolean(1)")
	c.EnableNormalize()
	c.Test("var x = new Object('a')", "var x = Object('a')")
	c.Test("var x = new RegExp('')", "var x = RegExp('')")
	c.Test("var x = new Error('20')", "var x = Error(\"20\")")
	c.Test("var x = new Array(20)", "var x = Array(20)")
}

// Java: PeepholeSubstituteAlternateSyntaxTest.testFoldTrueFalse()
func (c *PeepholeSubstituteAlternateSyntaxTest) TestFoldTrueFalse() {
	c.Test("x = true", "x = !0")
	c.Test("x = false", "x = !1")
}

// Java: PeepholeSubstituteAlternateSyntaxTest.testFoldTrueFalseComparison()
func (c *PeepholeSubstituteAlternateSyntaxTest) TestFoldTrueFalseComparison() {
	c.Test("x == true", "x == 1")
	c.Test("x == false", "x == 0")
	c.Test("x != true", "x != 1")
	c.Test("x < true", "x < 1")
	c.Test("x <= true", "x <= 1")
	c.Test("x > true", "x > 1")
	c.Test("x >= true", "x >= 1")
}

// Java: PeepholeSubstituteAlternateSyntaxTest.testFoldSubtractionAssignment()
func (c *PeepholeSubstituteAlternateSyntaxTest) TestFoldSubtractionAssignment() {
	c.Test("x -= 1", "--x")
	c.Test("x -= -1", "++x")
}

// Java: PeepholeSubstituteAlternateSyntaxTest.testFoldReturnResult()
func (c *PeepholeSubstituteAlternateSyntaxTest) TestFoldReturnResult() {
	c.TestSame("function f(){return !1;}")
	c.TestSame("function f(){return null;}")
	c.Test("function f(){return void 0;}", "function f(){return}")
	c.TestSame("function f(){return void foo();}")
	c.Test("function f(){return undefined;}", "function f(){return}")
	c.Test("function f(){if(a()){return undefined;}}", "function f(){if(a()){return}}")
}

// Java: PeepholeSubstituteAlternateSyntaxTest.testUndefined()
func (c *PeepholeSubstituteAlternateSyntaxTest) TestUndefined() {
	c.TestSame("var x = undefined")
	c.TestSame("function f(f) {var undefined=2;var x = undefined;}")
	c.EnableNormalize()
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.Test("var x = undefined", "var x=void 0")
	c.TestSame(
		"var undefined = 1;\n" +
			"function f() {var undefined=2;var x = undefined;}\n")
	c.TestSame("function f(undefined) {}")
	c.TestSame("try {} catch(undefined) {}")
	c.TestSame("for (undefined in {}) {}")
	c.TestSame("undefined++;")
	c.DisableNormalize()
	c.TestSame("undefined += undefined;")
	c.EnableNormalize()
	c.Test("undefined += undefined;", "undefined = void 0 + void 0;")
}

// Java: PeepholeSubstituteAlternateSyntaxTest.testSplitCommaExpressions()
func (c *PeepholeSubstituteAlternateSyntaxTest) TestSplitCommaExpressions() {
	c.late = false
	// Don't try to split in expressions.
	c.TestSame("while (foo(), !0) boo()")
	c.TestSame("var a = (foo(), !0);")
	c.TestSame("a = (foo(), !0);")
	// Don't try to split COMMA under LABELs.
	c.TestSame("a:a(),b()")
	c.Test("1, 2, 3, 4", "1; 2; 3; 4")
	c.Test("x = 1, 2, 3", "x = 1; 2; 3")
	c.TestSame("x = (1, 2, 3)")
	c.Test("1, (2, 3), 4", "1; 2; 3; 4")
	c.Test("(x=2), foo()", "x=2; foo()")
	c.Test("foo(), boo();", "foo(); boo()")
	c.Test("(a(), b()), (c(), d());", "a(); b(); c(); d()")
	c.Test("a(); b(); (c(), d());", "a(); b(); c(); d();")
	c.Test("foo(), true", "foo();true")
	c.TestSame("foo();true")
	c.Test("function x(){foo(), !0}", "function x(){foo(); !0}")
	c.TestSame("function x(){foo(); !0}")
}

// Java: PeepholeSubstituteAlternateSyntaxTest.testComma1()
func (c *PeepholeSubstituteAlternateSyntaxTest) TestComma1() {
	c.late = false
	c.Test("1, 2", "1; 2")
	c.late = true
	c.TestSame("1, 2")
}

// Java: PeepholeSubstituteAlternateSyntaxTest.testComma2()
func (c *PeepholeSubstituteAlternateSyntaxTest) TestComma2() {
	c.late = false
	c.Test("1, a()", "1; a()")
	c.Test("1, a?.()", "1; a?.()")
	c.late = true
	c.TestSame("1, a()")
	c.TestSame("1, a?.()")
}

// Java: PeepholeSubstituteAlternateSyntaxTest.testComma3()
func (c *PeepholeSubstituteAlternateSyntaxTest) TestComma3() {
	c.late = false
	c.Test("1, a(), b()", "1; a(); b()")
	c.Test("1, a?.(), b?.()", "1; a?.(); b?.()")
	c.late = true
	c.TestSame("1, a(), b()")
	c.TestSame("1, a?.(), b?.()")
}

// Java: PeepholeSubstituteAlternateSyntaxTest.testComma4()
func (c *PeepholeSubstituteAlternateSyntaxTest) TestComma4() {
	c.late = false
	c.Test("a(), b()", "a();b()")
	c.Test("a?.(), b?.()", "a?.();b?.()")
	c.late = true
	c.TestSame("a(), b()")
	c.TestSame("a?.(), b?.()")
}

// Java: PeepholeSubstituteAlternateSyntaxTest.testComma5()
func (c *PeepholeSubstituteAlternateSyntaxTest) TestComma5() {
	c.late = false
	c.Test("a(), b(), 1", "a(); b(); 1")
	c.Test("a?.(), b?.(), 1", "a?.(); b?.(); 1")
	c.late = true
	c.TestSame("a(), b(), 1")
	c.TestSame("a?.(), b?.(), 1")
}

// Java: PeepholeSubstituteAlternateSyntaxTest.testStringArraySplitting()
func (c *PeepholeSubstituteAlternateSyntaxTest) TestStringArraySplitting() {
	c.TestSame("var x=['1','2','3','4']")
	c.TestSame("var x=['1','2','3','4','5']")
	c.Test("var x=['1','2','3','4','5','6']", "var x='123456'.split('')")
	c.Test("var x=['1','2','3','4','5','00']", "var x='1 2 3 4 5 00'.split(' ')")
	c.Test("var x=['1','2','3','4','5','6','7']", "var x='1234567'.split('')")
	c.Test("var x=['1','2','3','4','5','6','00']", "var x='1 2 3 4 5 6 00'.split(' ')")
	c.Test("var x=[' ,',',',',',',',',',',']", "var x=' ,;,;,;,;,;,'.split(';')")
	c.Test("var x=[',,',' ',',',',',',',',']", "var x=',,; ;,;,;,;,'.split(';')")
	c.Test("var x=['a,',' ',',',',',',',',']", "var x='a,; ;,;,;,;,'.split(';')")
	// all possible delimiters used, leave it alone
	c.TestSame("var x=[',', ' ', ';', '{', '}']")
}

// Java: PeepholeSubstituteAlternateSyntaxTest.testTemplateStringToString()
func (c *PeepholeSubstituteAlternateSyntaxTest) TestTemplateStringToString() {
	c.Test("`abcde`", "'abcde'")
	c.Test("`ab cd ef`", "'ab cd ef'")
	c.TestSame("`hello ${name}`")
	c.TestSame("tag `hello ${name}`")
	c.TestSame("tag `hello`")
	c.Test("`hello ${'foo'}`", "'hello foo'")
	c.Test("`${2} bananas`", "'2 bananas'")
	c.Test("`This is ${true}`", "'This is true'")
}

// Java: PeepholeSubstituteAlternateSyntaxTest.testBindToCall1()
func (c *PeepholeSubstituteAlternateSyntaxTest) TestBindToCall1() {
	c.Test("(goog.bind(f))()", "f()")
	c.Test("(goog.bind(f,a))()", "f.call(a)")
	c.Test("(goog.bind(f,a,b))()", "f.call(a,b)")
	c.Test("(goog.bind(f))(a)", "f(a)")
	c.Test("(goog.bind(f,a))(b)", "f.call(a,b)")
	c.Test("(goog.bind(f,a,b))(c)", "f.call(a,b,c)")
	c.Test("(goog.partial(f))()", "f()")
	c.Test("(goog.partial(f,a))()", "f(a)")
	c.Test("(goog.partial(f,a,b))()", "f(a,b)")
	c.Test("(goog.partial(f))(a)", "f(a)")
	c.Test("(goog.partial(f,a))(b)", "f(a,b)")
	c.Test("(goog.partial(f,a,b))(c)", "f(a,b,c)")
	c.Test("((function(){}).bind())()", "((function(){}))()")
	c.Test("((function(){}).bind(a))()", "((function(){})).call(a)")
	c.Test("((function(){}).bind(a,b))()", "((function(){})).call(a,b)")
	c.Test("((function(){}).bind())(a)", "((function(){}))(a)")
	c.Test("((function(){}).bind(a))(b)", "((function(){})).call(a,b)")
	c.Test("((function(){}).bind(a,b))(c)", "((function(){})).call(a,b,c)")
	// Without using type information we don't know "f" is a function.
	c.TestSame("(f.bind())()")
	c.TestSame("(f.bind(a))()")
	c.TestSame("(f.bind())(a)")
	c.TestSame("(f.bind(a))(b)")
	// Don't rewrite if the bind isn't the immediate call target
	c.TestSame("(goog.bind(f)).call(g)")
}

// Java: PeepholeSubstituteAlternateSyntaxTest.testBindToCall2()
func (c *PeepholeSubstituteAlternateSyntaxTest) TestBindToCall2() {
	c.Test("(goog$bind(f))()", "f()")
	c.Test("(goog$bind(f,a))()", "f.call(a)")
	c.Test("(goog$bind(f,a,b))()", "f.call(a,b)")
	c.Test("(goog$bind(f))(a)", "f(a)")
	c.Test("(goog$bind(f,a))(b)", "f.call(a,b)")
	c.Test("(goog$bind(f,a,b))(c)", "f.call(a,b,c)")
	c.Test("(goog$partial(f))()", "f()")
	c.Test("(goog$partial(f,a))()", "f(a)")
	c.Test("(goog$partial(f,a,b))()", "f(a,b)")
	c.Test("(goog$partial(f))(a)", "f(a)")
	c.Test("(goog$partial(f,a))(b)", "f(a,b)")
	c.Test("(goog$partial(f,a,b))(c)", "f(a,b,c)")
	// Don't rewrite if the bind isn't the immediate call target
	c.TestSame("(goog$bind(f)).call(g)")
}

// Java: PeepholeSubstituteAlternateSyntaxTest.testBindToCall3()
func (c *PeepholeSubstituteAlternateSyntaxTest) TestBindToCall3() {
	// TODO(johnlenz): The code generator wraps free calls with (0,...) to
	// prevent leaking "this", but the parser doesn't unfold it, making a
	// AST comparison fail.  For now do a string comparison to validate the
	// correct code is in fact generated.
	// The FREE call wrapping should be moved out of the code generator
	// and into a denormalizing pass.
	c.DisableCompareAsTree()
	c.retraverseOnChange = true
	c.late = false
	c.Test("(goog.bind(f.m))()", "(0,f.m)()")
	c.Test("(goog.bind(f.m,a))()", "f.m.call(a)")
	c.Test("(goog.bind(f.m))(a)", "(0,f.m)(a)")
	c.Test("(goog.bind(f.m,a))(b)", "f.m.call(a,b)")
	c.Test("(goog.partial(f.m))()", "(0,f.m)()")
	c.Test("(goog.partial(f.m,a))()", "(0,f.m)(a)")
	c.Test("(goog.partial(f.m))(a)", "(0,f.m)(a)")
	c.Test("(goog.partial(f.m,a))(b)", "(0,f.m)(a,b)")
	// Without using type information we don't know "f" is a function.
	c.TestSame("f.m.bind()()")
	c.TestSame("f.m.bind(a)()")
	c.TestSame("f.m.bind()(a)")
	c.TestSame("f.m.bind(a)(b)")
	// Don't rewrite if the bind isn't the immediate call target
	c.TestSame("goog.bind(f.m).call(g)")
}

// Java: PeepholeSubstituteAlternateSyntaxTest.testSimpleFunctionCall1()
func (c *PeepholeSubstituteAlternateSyntaxTest) TestSimpleFunctionCall1() {
	c.Test("var a = String(23)", "var a = '' + 23")
	// Don't fold the existence check to preserve behavior
	c.TestSame("var a = String?.(23)")
	c.Test("var a = String('hello')", "var a = '' + 'hello'")
	// Don't fold the existence check to preserve behavior
	c.TestSame("var a = String?.('hello')")
	c.TestSame("var a = String('hello', bar());")
	c.TestSame("var a = String({valueOf: function() { return 1; }});")
}

// Java: PeepholeSubstituteAlternateSyntaxTest.testSimpleFunctionCall2()
func (c *PeepholeSubstituteAlternateSyntaxTest) TestSimpleFunctionCall2() {
	c.Test("var a = Boolean(true)", "var a = !0")
	// Don't fold the existence check to preserve behavior
	c.Test("var a = Boolean?.(true)", "var a = Boolean?.(!0)")
	c.Test("var a = Boolean(false)", "var a = !1")
	// Don't fold the existence check to preserve behavior
	c.Test("var a = Boolean?.(false)", "var a = Boolean?.(!1)")
	c.Test("var a = Boolean(1)", "var a = !!1")
	// Don't fold the existence check to preserve behavior
	c.TestSame("var a = Boolean?.(1)")
	c.Test("var a = Boolean(x)", "var a = !!x")
	// Don't fold the existence check to preserve behavior
	c.TestSame("var a = Boolean?.(x)")
	c.Test("var a = Boolean({})", "var a = !!{}")
	// Don't fold the existence check to preserve behavior
	c.TestSame("var a = Boolean?.({})")
	c.TestSame("var a = Boolean()")
	c.TestSame("var a = Boolean(!0, !1);")
}

// Java: PeepholeSubstituteAlternateSyntaxTest.testRotateAssociativeOperators()
func (c *PeepholeSubstituteAlternateSyntaxTest) TestRotateAssociativeOperators() {
	// Multiplication is not associative because it can include floating point numbers e.g.
	// 1e-300 * 1e300 * 1e9 does not equal 1e-300 * (1e300 * 1e9).
	c.Test("a || (b || c); a * (b * c); a | (b | c)", "a || b || c; b * c * a; a | b | c")
	c.TestSame("a % (b % c); a / (b / c); a - (b - c);")
	c.TestSame("(a / b) & (c % d)")
	c.TestSame("(c = 5) & (c % d)")
	c.TestSame("(a + b) * c * (d % e)")
	c.Test("(a + b) * (c % d)", "c % d * (a + b)")
}

// Java: PeepholeSubstituteAlternateSyntaxTest.testRotateCommutativeeOperators()
func (c *PeepholeSubstituteAlternateSyntaxTest) TestRotateCommutativeeOperators() {
	c.Test("a * (b % c);", "b % c * a")
	c.TestSame("a * b * (c / d)")
	c.TestSame("!a * c * (d % e)")
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: PeepholeSubstituteAlternateSyntaxTest.nullishCoalesce()
func (c *PeepholeSubstituteAlternateSyntaxTest) TestNullishCoalesce() {
	c.Test("a ?? (b ?? c);", "(a ?? b) ?? c")
}

// Java: PeepholeSubstituteAlternateSyntaxTest.testNoRotateInfiniteLoop()
func (c *PeepholeSubstituteAlternateSyntaxTest) TestNoRotateInfiniteLoop() {
	c.Test("1/x || (y/1 ||(1/z))", "1/x || (y/1) || (1/z)")
	c.TestSame("1/x || (y/1) || (1/z)")
}
