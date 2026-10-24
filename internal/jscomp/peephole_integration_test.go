// Port of com/google/javascript/jscomp/PeepholeIntegrationTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/PeepholeIntegrationTest.java)

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

// PeepholeIntegrationTest is the tests for the interaction between multiple
// peephole passes.
type PeepholeIntegrationTest struct {
	CompilerTestCaseBase

	late           bool
	numRepetitions int
}

// Java: PeepholeIntegrationTest()
func NewPeepholeIntegrationTest() *PeepholeIntegrationTest {
	c := &PeepholeIntegrationTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestPeepholeIntegrationTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewPeepholeIntegrationTest() })
}

// Java: PeepholeIntegrationTest.setUp()
func (c *PeepholeIntegrationTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableNormalize()
	c.EnableComputeSideEffects()
	// TODO(bradfordcsmith): Stop normalizing the expected output or document why it is necessary.
	c.EnableNormalizeExpectedOutput()
	c.late = false
	c.numRepetitions = 2
}

// Java: PeepholeIntegrationTest.getProcessor(Compiler)
func (c *PeepholeIntegrationTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewPeepholeOptimizationsPass(
		compiler,
		c.GetName(),
		NewPeepholeMinimizeConditions(c.late),
		NewPeepholeSubstituteAlternateSyntax(c.late),
		NewPeepholeRemoveDeadCode(),
		NewPeepholeFoldConstants(c.late, false /* useTypes */),
		NewPeepholeReplaceKnownMethods(c.late /* useTypes= */, false),
		NewMinimizeExitPoints())
}

// Java: PeepholeIntegrationTest.getNumRepetitions()
func (c *PeepholeIntegrationTest) GetNumRepetitions() int {
	return c.numRepetitions
}

// Java: PeepholeIntegrationTest.testUselessLabels()
func (c *PeepholeIntegrationTest) TestUselessLabels() {
	c.late = false
	c.Test("a:b:{break a;}", "")
	c.Test("a:b:{break b;}", "")
	c.Test("a:{break a;}", "")
	c.late = true
	c.Test("a:b:{break a;}", "")
	c.Test("a:b:{break b;}", "")
	c.Test("a:{break a;}", "")
}

// Java: PeepholeIntegrationTest.testTrueFalse()
func (c *PeepholeIntegrationTest) TestTrueFalse() {
	c.late = false
	c.TestSame("x = true")
	c.TestSame("x = false")
	c.Test("x = !1", "x = false")
	c.Test("x = !0", "x = true")
	c.late = true
	c.Test("x = true", "x = !0")
	c.Test("x = false", "x = !1")
	c.TestSame("x = !1")
	c.TestSame("x = !0")
}

// Java: PeepholeIntegrationTest.testFoldFollowing()
func (c *PeepholeIntegrationTest) TestFoldFollowing() {
	c.Test("function f(){return; console.log(1);}", "function f(){}")
}

// Check that removing blocks with 1 child works
// Java: PeepholeIntegrationTest.testFoldOneChildBlocksIntegration()
func (c *PeepholeIntegrationTest) TestFoldOneChildBlocksIntegration() {
	c.Test("function f(){switch(foo()){default:{break}}}", "function f(){foo()}")
	c.Test("function f(){switch(x){default:{break}}}", "function f(){}")
	c.Test(
		"function f(){switch(x){default:x;case 1:return 2}}",
		"function f(){switch(x){default:case 1:return 2}}")
	// ensure that block folding does not break hook ifs
	c.Test("if(x){if(true){foo();foo()}else{bar();bar()}}", "if(x){foo();foo()}")
	c.Test("if(x){if(false){foo();foo()}else{bar();bar()}}", "if(x){bar();bar()}")
	// Cases where the then clause has no side effects.
	c.Test("if(x()){}", "x()")
	c.Test("if(x()){} else {x()}", "x()||x()")
	c.Test("if(x){}", "")
	// Even the condition has no side effect.
	c.Test("if(a()){A()} else if (b()) {} else {C()}", "a()?A():b()||C()")
	c.Test("if(a()){} else if (b()) {} else {C()}", "a() || (b() || C())")
	c.Test(
		"if(a()){A()} else if (b()) {} else if (c()) {} else{D()}",
		"a() ? A() : b() || (c() || D())")
	c.Test("if(a()){} else if (b()) {} else if (c()) {} else{D()}", "a() || (b() || (c() || D()))")
	c.Test("if(a()){A()} else if (b()) {} else if (c()) {} else{}", "a()?A():b()||c()")
	// Verify that non-global scope works.
	c.Test("function foo(){if(x()){}}", "function foo(){x()}")
	c.Test("function foo(){if(x?.()){}}", "function foo(){x?.()}")
}

// Java: PeepholeIntegrationTest.testFoldOneChildBlocksStringCompare()
func (c *PeepholeIntegrationTest) TestFoldOneChildBlocksStringCompare() {
	c.Test("if (x) {if (y) { var x; } } else{ var z; }", "if (x) { if (y) var x } else var z")
}

// Test a particularly hairy edge case.
// Java: PeepholeIntegrationTest.testNecessaryDanglingElse()
func (c *PeepholeIntegrationTest) TestNecessaryDanglingElse() {
	c.Test("if (x) if (y){ y(); z() } else; else x()", "if (x) { if(y) { y(); z() } } else x()")
}

// Try to minimize returns
// Java: PeepholeIntegrationTest.testFoldReturnsIntegration()
func (c *PeepholeIntegrationTest) TestFoldReturnsIntegration() {
	// if-then-else duplicate statement removal handles this case:
	c.Test("function f(){if(x)return;else return}", "function f(){}")
}

// Try to minimize returns
// Java: PeepholeIntegrationTest.testFoldReturnsIntegrationWithScoped()
func (c *PeepholeIntegrationTest) TestFoldReturnsIntegrationWithScoped() {
	c.late = true
	c.DisableComputeSideEffects()
	// if-then-else duplicate statement removal handles this case:
	c.TestSame("function test(a) {if (a) {const a = Math.random();if(a) {return a;}} return a; }")
}

// Java: PeepholeIntegrationTest.testBug1059649()
func (c *PeepholeIntegrationTest) TestBug1059649() {
	// ensure that folding blocks with a single var node doesn't explode
	c.Test("if(x){var y=3;}var z=5", "if(x)var y=3;var z=5")
	c.Test("for(var i=0;i<10;i++){var y=3;}var z=5", "for(var i=0;i<10;i++)var y=3;var z=5")
	c.Test("for(var i in x){var y=3;}var z=5", "for(var i in x)var y=3;var z=5")
	c.Test("do{var y=3;}while(x);var z=5", "do var y=3;while(x);var z=5")
}

// Java: PeepholeIntegrationTest.testHookIfIntegration()
func (c *PeepholeIntegrationTest) TestHookIfIntegration() {
	c.Test("if (false){ x = 1; } else if (cond) { x = 2; } else { x = 3; }", "x=cond?2:3")

	c.Test("x?void 0:y()", "x||y()")
	c.Test("!x?void 0:y()", "x&&y()")
	c.Test("x?y():void 0", "x&&y()")

	c.TestSame("window.f = x => x ? true : false;")
	c.TestSame("window.f = x => x ? false : true;")
}

// Java: PeepholeIntegrationTest.testRemoveDuplicateStatementsIntegration()
func (c *PeepholeIntegrationTest) TestRemoveDuplicateStatementsIntegration() {
	c.Test(
		"function z() {if (a) { return true }\n"+
			"else if (b) { return true }\n"+
			"else { return true }}\n",
		"function z() {return true;}")
	c.Test(
		"function z() {if (a()) { return true }\n"+
			"else if (b()) { return true }\n"+
			"else { return true }}\n",
		"function z() {a()||b();return true;}")
}

// Java: PeepholeIntegrationTest.testFoldLogicalOpIntegration()
func (c *PeepholeIntegrationTest) TestFoldLogicalOpIntegration() {
	c.Test("if(x && true) z()", "x&&z()")
	c.Test("if(x && false) z()", "")
	c.Test("if(x || 3) z()", "z()")
	c.Test("if(x || false) z()", "x&&z()")
	c.Test("if(x==y && false) z()", "")
	c.Test("if(y() || x || 3) z()", "y();z()")
}

// Java: PeepholeIntegrationTest.testFoldBitwiseOpStringCompareIntegration()
func (c *PeepholeIntegrationTest) TestFoldBitwiseOpStringCompareIntegration() {
	c.Test("for (;-1 | 0;) {}", "for (;;);")
}

// Java: PeepholeIntegrationTest.testVarLiftingIntegration()
func (c *PeepholeIntegrationTest) TestVarLiftingIntegration() {
	c.Test("if(true);else var a;", "var a")
	c.Test("if(false) foo();else var a;", "var a")
	c.Test("if(true)var a;else;", "var a")
	c.Test("if(false)var a;else;", "var a")
	c.Test("if(false)var a,b;", "var b; var a")
	c.Test("if(false){var a;var a;}", "var a")
	c.Test("if(false)var a=function(){var b};", "var a")
	c.Test("if(a)if(false)var a;else var b;", "var a;if(a)var b")
}

// Java: PeepholeIntegrationTest.testBug1438784()
func (c *PeepholeIntegrationTest) TestBug1438784() {
	c.Test("for(var i=0;i<10;i++)if(x)x.y;", "for(var i=0;i<10;i++);")
}

// Java: PeepholeIntegrationTest.testFoldUselessForIntegration()
func (c *PeepholeIntegrationTest) TestFoldUselessForIntegration() {
	c.Test("for(;!true;) { foo() }", "")
	c.Test("for(;void 0;) { foo() }", "")
	c.Test("for(;undefined;) { foo() }", "")
	c.Test("for(;1;) foo()", "for(;;) foo()")
	c.Test("for(;!void 0;) foo()", "for(;;) foo()")
	// Make sure proper empty nodes are inserted.
	c.Test("if(foo())for(;false;){foo()}else bar()", "foo()||bar()")
	c.Test("if(foo?.())for(;false;){foo?.()}else bar?.()", "foo?.()||bar?.()")
}

// Java: PeepholeIntegrationTest.testFoldUselessDoIntegration()
func (c *PeepholeIntegrationTest) TestFoldUselessDoIntegration() {
	c.Test("do { foo() } while(!true);", "foo()")
	c.Test("do { foo() } while(void 0);", "foo()")
	c.Test("do { foo() } while(undefined);", "foo()")
	c.Test("do { foo() } while(!void 0);", "do { foo() } while(1);")
	// Make sure proper empty nodes are inserted.
	c.Test("if(foo())do {foo()} while(false) else bar()", "foo()?foo():bar()")
	// Optional chaining version of these tests.
	c.Test("do { foo?.() } while(!true);", "foo?.()")
	c.Test("do { foo?.() } while(void 0);", "foo?.()")
	c.Test("do { foo?.() } while(undefined);", "foo?.()")
	c.Test("do { foo?.() } while(!void 0);", "do { foo?.() } while(1);")
	// Make sure proper empty nodes are inserted.
	c.Test("if(foo?.())do {foo?.()} while(false) else bar?.()", "foo?.() ? foo?.() : bar?.()")
}

// Java: PeepholeIntegrationTest.testMinimizeExpr()
func (c *PeepholeIntegrationTest) TestMinimizeExpr() {
	c.Test("!!true", "")
	c.Test("!!x()", "x()")
	c.Test("!(!x()&&!y())", "x()||y()")
	c.Test("x()||!!y()", "x()||y()")
	// This is similar to the !!true case
	c.Test("!!x()&&y()", "x()&&y()")
	c.Test("!!x?.()&&y?.()", "x?.()&&y?.()")
}

// Java: PeepholeIntegrationTest.testBug1509085()
func (c *PeepholeIntegrationTest) TestBug1509085() {
	c.numRepetitions = 1
	c.late = true
	// Such code can be replaced by using the simpler, equivalent optional chaining operator
	// `x?.()`.
	c.Test("x ? x() : void 0", "x&&x();")
	c.TestSame("y = x ? x() : void 0")
	c.TestSame("x?.()")
}

// Java: PeepholeIntegrationTest.testBugIssue3()
func (c *PeepholeIntegrationTest) TestBugIssue3() {
	c.TestSame(
		"function foo() {\n" +
			"  if(sections.length != 1) children[i] = 0;\n" +
			"  else var selectedid = children[i]\n" +
			"}\n")
}

// Java: PeepholeIntegrationTest.testBugIssue43()
func (c *PeepholeIntegrationTest) TestBugIssue43() {
	c.TestSame("function foo() {\n  if (a) { var b = 1; } else { a.b = 1; }\n}")
}

// Java: PeepholeIntegrationTest.testFoldNegativeBug()
func (c *PeepholeIntegrationTest) TestFoldNegativeBug() {
	c.Test("for (;-3;){};", "for (;;);")
}

// Java: PeepholeIntegrationTest.testNoNormalizeLabeledExpr()
func (c *PeepholeIntegrationTest) TestNoNormalizeLabeledExpr() {
	c.EnableNormalize()
	c.TestSame("var x; foo:{x = 3;}")
	c.TestSame("var x; foo:x = 3;")
	c.DisableNormalize()
}

// Java: PeepholeIntegrationTest.testShortCircuit1()
func (c *PeepholeIntegrationTest) TestShortCircuit1() {
	c.Test("1 && a()", "a()")
}

// Java: PeepholeIntegrationTest.testShortCircuit2()
func (c *PeepholeIntegrationTest) TestShortCircuit2() {
	c.Test("1 && a() && 2", "a()")
}

// Java: PeepholeIntegrationTest.testShortCircuit3()
func (c *PeepholeIntegrationTest) TestShortCircuit3() {
	c.Test("a() && 1 && 2", "a()")
}

// Java: PeepholeIntegrationTest.testShortCircuit4()
func (c *PeepholeIntegrationTest) TestShortCircuit4() {
	c.Test("a() && (1 && b())", "a() && b()")
	c.Test("a() && 1 && b()", "a() && b()")
	c.Test("(a() && 1) && b()", "a() && b()")
	c.Test("a?.() && (1 && b?.())", "a?.() && b?.()")
	c.Test("a?.() && 1 && b?.()", "a?.() && b?.()")
	c.Test("(a?.() && 1) && b?.()", "a?.() && b?.()")
}

// Java: PeepholeIntegrationTest.testMinimizeExprCondition()
func (c *PeepholeIntegrationTest) TestMinimizeExprCondition() {
	c.Test("(x || true) && y()", "y()")
	c.Test("(x || false) && y()", "x&&y()")
	c.Test("(x && true) && y()", "x && y()")
	c.Test("(x && false) && y()", "")
	c.Test("a = x || false ? b : c", "a=x?b:c")
	c.Test("do {x()} while((x && false) && y())", "x()")
	c.Test("do {x?.()} while((x && false) && y?.())", "x?.()")
}

// A few miscellaneous cases where one of the peephole passes increases the
// size, but it does it in such a way that a later pass can decrease it.
// Test to make sure the overall change is a decrease, not an increase.
// Java: PeepholeIntegrationTest.testMisc()
func (c *PeepholeIntegrationTest) TestMisc() {
	c.Test("x = [foo()] && x", "x = (foo(),x)")
	c.Test("x = foo() && false || bar()", "x = (foo(), bar())")
	c.Test("if(foo() && false) z()", "foo()")
}

// Java: PeepholeIntegrationTest.testTrueFalseFolding()
func (c *PeepholeIntegrationTest) TestTrueFalseFolding() {
	c.late = true
	c.Test("x = true", "x = !0")
	c.Test("x = false", "x = !1")
	c.Test("x = !3", "x = !1")
	c.Test("x = true && !0", "x = !0")
	c.Test("x = !!!!!!!!!!!!3", "x = !0")
	c.Test("if(!3){x()}", "")
	c.Test("if(!!3){x()}", "x()")
}

// Java: PeepholeIntegrationTest.testCommaSplitingConstantCondition()
func (c *PeepholeIntegrationTest) TestCommaSplitingConstantCondition() {
	c.late = false
	c.Test("(b=0,b=1);if(b)x=b;", "b=0;b=1;x=b;")
	c.Test("(b=0,b=1);if(b)x=b;", "b=0;b=1;x=b;")
}

// Java: PeepholeIntegrationTest.testAvoidCommaSplitting()
func (c *PeepholeIntegrationTest) TestAvoidCommaSplitting() {
	c.late = false
	c.Test("x(),y(),z()", "x();y();z()")
	c.Test("x?.(),y?.(),z?.()", "x?.();y?.();z?.()")
	c.late = true
	c.TestSame("x(),y(),z()")
	c.TestSame("x?.(),y?.(),z?.()")
}

// Java: PeepholeIntegrationTest.testObjectLiteral()
func (c *PeepholeIntegrationTest) TestObjectLiteral() {
	c.Test("({})", "")
	c.Test("({a:1})", "")
	c.Test("({a:foo()})", "foo()")
	c.Test("({'a':foo()})", "foo()")
	c.Test("({a:foo?.()})", "foo?.()")
	c.Test("({'a':foo?.()})", "foo?.()")
}

// Java: PeepholeIntegrationTest.testArrayLiteral()
func (c *PeepholeIntegrationTest) TestArrayLiteral() {
	c.Test("([])", "")
	c.Test("([1])", "")
	c.Test("([a])", "")
	c.Test("([foo()])", "foo()")
}

// Java: PeepholeIntegrationTest.testFoldIfs1()
func (c *PeepholeIntegrationTest) TestFoldIfs1() {
	c.Test(
		"function f() {if (x) return 1; else if (y) return 1;}",
		"function f() {if (x||y) return 1;}")
	c.Test(
		"function f() {if (x) return 1; else {if (y) return 1; else foo();}}",
		"function f() {if (x||y) return 1; foo();}")
}

// Java: PeepholeIntegrationTest.testFoldIfs2()
func (c *PeepholeIntegrationTest) TestFoldIfs2() {
	c.Test("function f() {if (x) { a(); } else if (y) { a() }}", "function f() {x?a():y&&a();}")
}

// Java: PeepholeIntegrationTest.testFoldHook2()
func (c *PeepholeIntegrationTest) TestFoldHook2() {
	c.Test("function f(a) {if (!a) return a; else return a;}", "function f(a) {return a}")
}

// Go-only: test-helper — renamed with the Test prefix that RunJUnit4 runs.
// Java: PeepholeIntegrationTest.disable_testFoldHook1()
func (c *PeepholeIntegrationTest) TestDisable_testFoldHook1() {
	c.Test("function f(a) {return (!a)?a:a;}", "function f(a) {return a}")
}

// Java: PeepholeIntegrationTest.testTemplateStringsKnownMethods()
func (c *PeepholeIntegrationTest) TestTemplateStringsKnownMethods() {
	c.Test("x = `abcdef`.indexOf('b')", "x = 1")
	c.Test("x = [`a`, `b`, `c`].join(``)", "x='abc'")
	c.Test("x = `abcdef`.substr(0,2)", "x = 'ab'")
	c.Test("x = `abcdef`.substring(0,2)", "x = 'ab'")
	c.Test("x = `abcdef`.slice(0,2)", "x = 'ab'")
	c.Test("x = `abcdef`.charAt(0)", "x = 'a'")
	c.Test("x = `abcdef`.charCodeAt(0)", "x = 97")
	c.Test("x = `abc`.toUpperCase()", "x = 'ABC'")
	c.Test("x = `ABC`.toLowerCase()", "x = 'abc'")
	c.Test("x = `\t\n\ufeff\t asd foo bar \r\n`.trim()", "x = 'asd foo bar'")
	c.Test("x = parseInt(`123`)", "x = 123")
	c.Test("x = parseFloat(`1.23`)", "x = 1.23")
}

// Java: PeepholeIntegrationTest.testDontFoldKnownMethodsWithOptionalChaining()
func (c *PeepholeIntegrationTest) TestDontFoldKnownMethodsWithOptionalChaining() {
	// Known methods guarded by an optional chain are not folded
	c.Test("x = `abcdef`.indexOf?.('b')", "x = \"abcdef\".indexOf?.(\"b\");")
	c.Test("x = [`a`, `b`, `c`].join?.(``)", "x = [\"a\", \"b\", \"c\"].join?.(\"\")")
	c.Test("x = `abcdef`.substr?.(0,2)", "x = \"abcdef\".substr?.(0,2)")
	c.Test("x = `abcdef`.substring?.(0,2)", "x = \"abcdef\".substring?.(0,2)")
	c.Test("x = `abcdef`.slice?.(0,2)", "x = \"abcdef\".slice?.(0,2)")
	c.Test("x = `abcdef`.charAt?.(0)", "x = \"abcdef\".charAt?.(0)")
	c.Test("x = `abcdef`.charCodeAt?.(0)", "x = \"abcdef\".charCodeAt?.(0)")
	c.Test("x = `abc`.toUpperCase?.()", "x = \"abc\".toUpperCase?.()")
	c.Test("x = `ABC`.toLowerCase?.()", "x = \"ABC\".toLowerCase?.()")
	c.Test("x = parseInt?.(`123`)", "x = parseInt?.(\"123\")")
	c.Test("x = parseFloat?.(`1.23`)", "x = parseFloat?.(\"1.23\")")
}

// Java: PeepholeIntegrationTest.testLabeledBlocks()
func (c *PeepholeIntegrationTest) TestLabeledBlocks() {
	c.Test(
		"function b(m) {\n"+
			" return m;\n"+
			" label: {\n"+
			"   START('debug');\n"+
			"   label2: {\n"+
			"     alert('Shouldnt be here' + m);\n"+
			"   }\n"+
			"   END('debug');\n"+
			"  }\n"+
			"}\n",
		"function b(m) {\n"+
			"  return m;\n"+
			"}\n")
}

// Java: PeepholeIntegrationTest.testDoNotRemoveDeclarationOfUsedVariable()
func (c *PeepholeIntegrationTest) TestDoNotRemoveDeclarationOfUsedVariable() {
	c.Test(
		"var f = function() {\n"+
			"  return 1;\n"+
			"  let b = 5;\n"+
			"  do {\n"+
			"    b--;\n"+
			"  } while (b);\n"+
			"  return 3;\n"+
			"};\n",
		"var f = function() {\n"+
			"  let b;\n"+
			"  return 1;\n"+
			"};\n")
}

// Java: PeepholeIntegrationTest.testDontRemoveExport()
func (c *PeepholeIntegrationTest) TestDontRemoveExport() {
	c.Test(
		"function foo() {\n"+
			"  return 1;\n"+
			"  alert(2);\n"+
			"}\n"+
			"export { foo as foo };\n",
		"function foo() {\n"+
			"  return 1;\n"+
			"}\n"+
			"export { foo as foo };\n")
}

// Java: PeepholeIntegrationTest.testRemoveUnreachableCode1()
func (c *PeepholeIntegrationTest) TestRemoveUnreachableCode1() {
	// switch statement with stuff after "return"
	c.Test(
		"function foo(){\n"+
			"  switch (foo) {\n"+
			"    case 1:\n"+
			"      x=1;\n"+
			"      return;\n"+
			"      break;\n"+
			"    case 2: {\n"+
			"      x=2;\n"+
			"      return;\n"+
			"      break;\n"+
			"    }\n"+
			"    default:\n"+
			"  }\n"+
			"}\n",
		"function foo() {\n"+
			"  switch (foo) {\n"+
			"    case 1:\n"+
			"      x=1;\n"+
			"      break;\n"+
			"    case 2:\n"+
			"      x=2;\n"+
			"  }\n"+
			"}\n")
}

// Java: PeepholeIntegrationTest.testRemoveUnreachableCode2()
func (c *PeepholeIntegrationTest) TestRemoveUnreachableCode2() {
	// if/else statements with returns
	c.Test(
		"function bar(){\n"+
			"  if (foo)\n"+
			"    x=1;\n"+
			"  else if(bar) {\n"+
			"    return;\n"+
			"    x=2;\n"+
			"  } else {\n"+
			"    x=3;\n"+
			"    return;\n"+
			"    x=4;\n"+
			"  }\n"+
			"  return 5;\n"+
			"  x=5;\n"+
			"}\n",
		"function bar() {\n"+
			"  if (foo) {\n"+
			"    x=1;\n"+
			"    return 5;\n"+
			"  }\n"+
			"  bar || (x = 3);\n"+
			"}\n")
	// if statements without blocks
	// NOTE: This pass should never see while-loops, because normalization replaces them all with
	// for-loops.
	c.Test(
		"function foo() {\n"+
			"  if (x == 3) return;\n"+
			"  x = 4;\n"+
			"  y++;\n"+
			"  for (; y == 4; ) {\n"+
			"    return;\n"+
			"    x = 3\n"+
			"  }\n"+
			"}\n",
		"function foo() {\n"+
			"  if (x != 3) {\n"+
			"    x = 4;\n"+
			"    y++;\n"+
			"    for (; y == 4; ) {\n"+
			"      break\n"+
			"    }\n"+
			"  }\n"+
			"}\n")
	// for/do/while loops
	c.Test(
		"function baz() {\n"+
			"// Normalize always moves the for-loop initializer out of the loop.\n"+
			"  i = 0;\n"+
			"  for (; i < n; i++) {\n"+
			"    x = 3;\n"+
			"    break;\n"+
			"    x = 4\n"+
			"  }\n"+
			"  do {\n"+
			"    x = 2;\n"+
			"    break;\n"+
			"    x = 4\n"+
			"  } while (x == 4);\n"+
			"  for (; i < 4; ) {\n"+
			"    x = 3;\n"+
			"    return;\n"+
			"    x = 6\n"+
			"  }\n"+
			"}\n",
		"function baz() {\n"+
			"  i = 0;\n"+
			"  for (; i < n; i++) {\n"+
			"    x = 3;\n"+
			"    break\n"+
			"  }\n"+
			"  do {\n"+
			"    x = 2;\n"+
			"    break\n"+
			"  } while (x == 4);\n"+
			"  for (; i < 4; ) {\n"+
			"    x = 3;\n"+
			"    break;\n"+
			"  }\n"+
			"}\n")
	// return statements on the same level as conditionals
	c.Test(
		"function foo() {\n"+
			"  if (x == 3) {\n"+
			"    return\n"+
			"  }\n"+
			"  return 5;\n"+
			"  while (y == 4) {\n"+
			"    x++;\n"+
			"    return;\n"+
			"    x = 4\n"+
			"  }\n"+
			"}\n",
		"function foo() {\n"+
			"  return x == 3 ? void 0 : 5;\n"+
			"}\n")
	// return statements on the same level as conditionals
	c.Test(
		"function foo() {\n"+
			"  return 3;\n"+
			"  for (; y == 4;) {\n"+
			"    x++;\n"+
			"    return;\n"+
			"    x = 4\n"+
			"  }\n"+
			"}\n",
		"function foo() {\n"+
			"  return 3\n"+
			"}\n")
	// try/catch statements
	c.Test(
		"function foo() {\n"+
			"  try {\n"+
			"    x = 3;\n"+
			"    return x + 1;\n"+
			"    x = 5\n"+
			"  } catch (e) {\n"+
			"    x = 4;\n"+
			"    return 5;\n"+
			"    x = 5\n"+
			"  }\n"+
			"}\n",
		"function foo() {\n"+
			"  try {\n"+
			"    x = 3;\n"+
			"    return x + 1\n"+
			"  } catch (e) {\n"+
			"    x = 4;\n"+
			"    return 5\n"+
			"  }\n"+
			"}\n")
	// try/finally statements
	c.Test(
		"function foo() {\n"+
			"  try {\n"+
			"    x = 3;\n"+
			"    return x + 1;\n"+
			"    x = 5\n"+
			"  } finally {\n"+
			"    x = 4;\n"+
			"    return 5;\n"+
			"    x = 5\n"+
			"  }\n"+
			"}\n",
		"function foo() {\n"+
			"  try {\n"+
			"    x = 3;\n"+
			"    return x + 1\n"+
			"  } finally {\n"+
			"    x = 4;\n"+
			"    return 5\n"+
			"  }\n"+
			"}\n")
	// try/catch/finally statements
	c.Test(
		"function foo() {\n"+
			"  try {\n"+
			"    x = 3;\n"+
			"    return x + 1;\n"+
			"    x = 5\n"+
			"  } catch (e) {\n"+
			"    x = 3;\n"+
			"    return;\n"+
			"    x = 2\n"+
			"  } finally {\n"+
			"    x = 4;\n"+
			"    return 5;\n"+
			"    x = 5\n"+
			"  }\n"+
			"}\n",
		"function foo() {\n"+
			"  try {\n"+
			"    x = 3;\n"+
			"    return x + 1\n"+
			"  } catch (e) {\n"+
			"    x = 3;\n"+
			"  } finally {\n"+
			"    x = 4;\n"+
			"    return 5\n"+
			"  }\n"+
			"}\n")
	// test a combination of blocks
	c.Test(
		"function foo() {\n"+
			"  x = 3;\n"+
			"  if (x == 4) {\n"+
			"    x = 5;\n"+
			"    return;\n"+
			"    x = 6\n"+
			"  } else {\n"+
			"    x = 7\n"+
			"  }\n"+
			"  return 5;\n"+
			"  x = 3\n"+
			"}\n",
		"function foo() {\n"+
			"  x = 3;\n"+
			"  if (x == 4) {\n"+
			"    x = 5;\n"+
			"  } else {\n"+
			"    x = 7\n"+
			"    return 5\n"+
			"  }\n"+
			"}\n")
	// test removing multiple statements
	c.Test(
		"function foo() {\n"+
			"  return 1;\n"+
			"  var x = 2;\n"+
			"  var y = 10;\n"+
			"  return 2;\n"+
			"}\n",
		"function foo() {\n"+
			"  var y;\n"+
			"  var x;\n"+
			"  return 1\n"+
			"}\n")
	c.Test(
		"function foo() {\n"+
			"  return 1;\n"+
			"  x = 2;\n"+
			"  y = 10;\n"+
			"  return 2;\n"+
			"}\n",
		"function foo() {\n"+
			"  return 1\n"+
			"}\n")
}

// Java: PeepholeIntegrationTest.testRemoveUselessNameStatements()
func (c *PeepholeIntegrationTest) TestRemoveUselessNameStatements() {
	c.Test("a;", "")
	c.Test("a.b;", "")
	c.Test("a.b.MyClass.prototype.memberName;", "")
}

// Java: PeepholeIntegrationTest.testRemoveUselessStrings()
func (c *PeepholeIntegrationTest) TestRemoveUselessStrings() {
	c.Test("'a';", "")
}

// Java: PeepholeIntegrationTest.testNoRemoveUseStrict()
func (c *PeepholeIntegrationTest) TestNoRemoveUseStrict() {
	c.Test("'use strict';", "'use strict'")
}

// Java: PeepholeIntegrationTest.testRemoveDo()
func (c *PeepholeIntegrationTest) TestRemoveDo() {
	c.TestSame("do { print(1); break } while(1)")
	// NOTE: This pass should never see while-loops, because normalization replaces them all with
	// for-loops.
	c.Test(
		"for (; 1;) { break; do { print(1); break } while(1) }",
		"for (;  ;) { break;                                 }")
}

// Java: PeepholeIntegrationTest.testRemoveUselessLiteralValueStatements()
func (c *PeepholeIntegrationTest) TestRemoveUselessLiteralValueStatements() {
	c.Test("true;", "")
	c.Test("'hi';", "")
	c.Test("if (x) 1;", "")
	// NOTE: This pass should never see while-loops, because normalization replaces them all with
	// for-loops.
	c.Test("for (; x;) 1;", "for (; x;);")
	c.Test("do 1; while (x);", "for (;x;);")
	c.Test("for (;;) 1;", "for (;;);")
	c.Test(
		"switch(x){case 1:true;case 2:'hi';default:true}", //
		"")
}

// Java: PeepholeIntegrationTest.testConditionalDeadCode()
func (c *PeepholeIntegrationTest) TestConditionalDeadCode() {
	c.Test("function f() { if (1) return 5; else return 5; x = 1}", "function f() { return 5;}")
}

// Java: PeepholeIntegrationTest.testSwitchCase()
func (c *PeepholeIntegrationTest) TestSwitchCase() {
	c.T().Skip("TODO(b/301641291): this was originally supported by UnreachableCodeElimination, which was removed in favor of peephole optimizations. Support this test case if found useful in the real code.") // @Ignore
	c.Test("function f() { switch(x) { default: return 5; foo()}}", "function f() { return 5; }")
	c.TestSame("function f() { switch(x) { default: return; case 1: foo(); bar()}}")
	c.Test(
		"function f() { switch(x) { default: return; case 1: return 5;bar()}}",
		"function f() { switch(x) { default: return; case 1: return 5;}}")
}

// Java: PeepholeIntegrationTest.testTryCatchFinally1()
func (c *PeepholeIntegrationTest) TestTryCatchFinally1() {
	c.TestSame("try {foo()} catch (e) {bar()}")
	c.TestSame("try { try {foo()} catch (e) {bar()}} catch (x) {bar()}")
}

// Java: PeepholeIntegrationTest.testTryCatchFinally2()
func (c *PeepholeIntegrationTest) TestTryCatchFinally2() {
	c.TestSame("try {var x = 1} catch (e) {e()}")
}

// Java: PeepholeIntegrationTest.testTryCatchFinally3()
func (c *PeepholeIntegrationTest) TestTryCatchFinally3() {
	c.TestSame("try {var x = 1} catch (e) {e()} finally {x()}")
}

// Java: PeepholeIntegrationTest.testTryCatchFinally4()
func (c *PeepholeIntegrationTest) TestTryCatchFinally4() {
	c.TestSame("try {var x = 1} catch (e) {e()} finally {}")
}

// Java: PeepholeIntegrationTest.testTryCatchFinally5()
func (c *PeepholeIntegrationTest) TestTryCatchFinally5() {
	c.TestSame("try {var x = 1} finally {x()}")
	c.TestSame("var x = 1")
}

// Java: PeepholeIntegrationTest.testTryCatchFinally6()
func (c *PeepholeIntegrationTest) TestTryCatchFinally6() {
	c.Test("function f() {return; try{var x = 1}catch(e){} }", "function f() {var x;}")
}

// Java: PeepholeIntegrationTest.testRemovalRequiresRedeclaration()
func (c *PeepholeIntegrationTest) TestRemovalRequiresRedeclaration() {
	// NOTE: This pass should never see while-loops, because normalization replaces them all with
	// for-loops.
	c.Test(
		"for (; 1;) {\n"+
			"  break;\n"+
			"  var x = 1\n"+
			"}\n",
		"var x;\n"+
			"for (;;) {\n"+
			"  break;\n"+
			"}\n")
	c.Test(
		"for (; 1;) {\n"+
			"  break;\n"+
			"  var x=1;\n"+
			"  var y=1;\n"+
			"}\n",
		"var y;\n"+
			"var x;\n"+
			"for (;;) {\n"+
			"  break;\n"+
			"}\n")
}

// Java: PeepholeIntegrationTest.testAssignPropertyOnCreatedObject()
func (c *PeepholeIntegrationTest) TestAssignPropertyOnCreatedObject() {
	c.TestSame("this.foo = 3;")
	c.TestSame("a.foo = 3;")
	c.TestSame("bar().foo = 3;")
	c.TestSame("({}).foo = bar();")
	c.TestSame("(new X()).foo = 3;")
	c.Test("({}).foo = 3;", "")
	c.Test("(function() {}).prototype.toString = function(){};", "")
	c.Test("(function() {}).prototype['toString'] = function(){};", "")
	c.Test("(function() {}).prototype[f] = function(){};", "")
}

// Java: PeepholeIntegrationTest.testUselessUnconditionalReturn1()
func (c *PeepholeIntegrationTest) TestUselessUnconditionalReturn1() {
	c.Test("function foo() { return }", "function foo() { }")
}

// Java: PeepholeIntegrationTest.testUselessUnconditionalReturn2()
func (c *PeepholeIntegrationTest) TestUselessUnconditionalReturn2() {
	c.Test("function foo() { return; return; x=1 }", "function foo() { }")
}

// Java: PeepholeIntegrationTest.testUselessUnconditionalReturn3()
func (c *PeepholeIntegrationTest) TestUselessUnconditionalReturn3() {
	c.Test("function foo() { return; return; var x=1}", "function foo() {var x}")
}

// Java: PeepholeIntegrationTest.testUselessUnconditionalReturn4()
func (c *PeepholeIntegrationTest) TestUselessUnconditionalReturn4() {
	c.Test(
		"function foo() { return; function bar() {} }",
		"function foo() {         function bar() {} }")
}

// Java: PeepholeIntegrationTest.testUselessUnconditionalReturn5()
func (c *PeepholeIntegrationTest) TestUselessUnconditionalReturn5() {
	c.TestSame("function foo() { return 5 }")
}

// Java: PeepholeIntegrationTest.testUselessUnconditionalReturn6()
func (c *PeepholeIntegrationTest) TestUselessUnconditionalReturn6() {
	c.Test("function f() {switch (a) { case 'a': return}}", "function f() {}")
}

// Java: PeepholeIntegrationTest.testUselessUnconditionalReturn7()
func (c *PeepholeIntegrationTest) TestUselessUnconditionalReturn7() {
	c.T().Skip("TODO(b/301641291): this was originally supported by UnreachableCodeElimination, which was removed in favor of peephole optimizations. Support this test case if found useful in the real code.") // @Ignore
	c.TestSame("function f() {switch (a) { default: return; case 'a': alert(1)}}")
}

// Java: PeepholeIntegrationTest.testUselessUnconditionalReturn8()
func (c *PeepholeIntegrationTest) TestUselessUnconditionalReturn8() {
	c.T().Skip("TODO(b/301641291): this was originally supported by UnreachableCodeElimination, which was removed in favor of peephole optimizations. Support this test case if found useful in the real code.") // @Ignore
	c.TestSame("function f() {switch (a) { case 'a': return; default: alert(1)}}")
}

// Java: PeepholeIntegrationTest.testUselessUnconditionalReturn9()
func (c *PeepholeIntegrationTest) TestUselessUnconditionalReturn9() {
	c.TestSame("function f() {switch (a) { case 'a': case foo(): }}")
}

// Java: PeepholeIntegrationTest.testUselessUnconditionalContinue()
func (c *PeepholeIntegrationTest) TestUselessUnconditionalContinue() {
	c.Test("for(;1;) {continue}", "for(;;) {}")
	c.Test("for(;0;) {continue}", "")
}

// Java: PeepholeIntegrationTest.testUselessUnconditionalContinue2()
func (c *PeepholeIntegrationTest) TestUselessUnconditionalContinue2() {
	c.Test(
		"X: for(;1;) { for(;1;) { if (x()) {continue X} x = 1}}",
		"X: for(; ;) { for(; ;) { if (x()) {continue X} x = 1}}")
}

// Java: PeepholeIntegrationTest.testUselessUnconditionalContinue3()
func (c *PeepholeIntegrationTest) TestUselessUnconditionalContinue3() {
	c.T().Skip("TODO(b/301641291): this was originally supported by UnreachableCodeElimination, which was removed in favor of peephole optimizations. Support this test case if found useful in the real code.") // @Ignore
	c.Test(
		"for(;1;) { X: for(;1;) { if (x()) {continue X} }}",
		"for(; ;) { X: for(; ;) { if (x()) {}}}")
}

// Java: PeepholeIntegrationTest.testUselessUnconditionalContinue4()
func (c *PeepholeIntegrationTest) TestUselessUnconditionalContinue4() {
	c.T().Skip("TODO(b/301641291): this was originally supported by UnreachableCodeElimination, which was removed in favor of peephole optimizations. Support this test case if found useful in the real code.") // @Ignore
	c.Test("do { continue } while(1);", "do {  } while(1);")
}

// Java: PeepholeIntegrationTest.testUselessUnconditionalBreak1()
func (c *PeepholeIntegrationTest) TestUselessUnconditionalBreak1() {
	c.T().Skip("@Ignore") // @Ignore
	// TODO - b/335145701: `case 'a'` can not be eliminated as that would force "foo()" to be
	// always executed.
	// `case foo()` can only be eliminated if "foo()" can be determined to be side-effect
	// free.
	c.Test(
		"switch (a) { case 'a': break; case foo(): }",
		"switch (a) { case 'a':        case foo(): }")
}

// Java: PeepholeIntegrationTest.testUselessUnconditionalBreak2()
func (c *PeepholeIntegrationTest) TestUselessUnconditionalBreak2() {
	c.Test("switch (a) { case 'a': break }", "")
	c.Test(
		"switch (a) { default: break; case 'a': }", //
		"")
}

// Java: PeepholeIntegrationTest.testUselessUnconditionalBreak3()
func (c *PeepholeIntegrationTest) TestUselessUnconditionalBreak3() {
	c.TestSame("switch (a) { case 'a': alert(a); break; default: alert(a); }")
	c.TestSame("switch (a) { default: alert(a); break; case 'a': alert(a); }")
}

// Java: PeepholeIntegrationTest.testUselessUnconditionalBreak4()
func (c *PeepholeIntegrationTest) TestUselessUnconditionalBreak4() {
	c.Test("X: {switch (a) { case 'a': break X}}", "")
	c.Test(
		"X: {switch (a) { case 'a': if (a()) {break X}  a = 1; }}", //
		"X: {switch (a) { case 'a': a() || (a = 1); }}")
}

// Java: PeepholeIntegrationTest.testUselessUnconditionalBreak5()
func (c *PeepholeIntegrationTest) TestUselessUnconditionalBreak5() {
	c.Test(
		"X: {switch (a) { case 'a': if (a()) {break X}}}", //
		"X: {switch (a) { case 'a':     a()           }}")
}

// Java: PeepholeIntegrationTest.testUselessUnconditionalBreak6()
func (c *PeepholeIntegrationTest) TestUselessUnconditionalBreak6() {
	c.Test(
		"X: {switch (a) { case 'a': if (a()) {break X}}}", //
		"X: {switch (a) { case 'a':     a();          }}")
}

// Java: PeepholeIntegrationTest.testUselessUnconditionalBreak7()
func (c *PeepholeIntegrationTest) TestUselessUnconditionalBreak7() {
	// There is no reason to keep these
	c.TestSame("do { break } while(1);")
	c.Test("for(;1;) { break }", "for(;;) { break; }")
}

// Java: PeepholeIntegrationTest.testIteratedRemoval1()
func (c *PeepholeIntegrationTest) TestIteratedRemoval1() {
	c.Test("switch (a) { case 'a': break; case 'b': break; case 'c': break }", "")
}

// Java: PeepholeIntegrationTest.testIteratedRemoval2()
func (c *PeepholeIntegrationTest) TestIteratedRemoval2() {
	c.Test(
		"function foo() { switch (a) { case 'a':return; case 'b':return; case 'c':return }}",
		"function foo() { }")
}

// Java: PeepholeIntegrationTest.testIteratedRemoval3()
func (c *PeepholeIntegrationTest) TestIteratedRemoval3() {
	c.Test(
		"for (;;) {\n"+
			"   switch (a) {\n"+
			"   case 'a': continue;\n"+
			"   case 'b': continue;\n"+
			"   case 'c': continue;\n"+
			"   }\n"+
			" }\n",
		" for (;;) { }")
}

// Java: PeepholeIntegrationTest.testIteratedRemoval4()
func (c *PeepholeIntegrationTest) TestIteratedRemoval4() {
	c.Test("function foo() { if (x) { return; } if (x) { return; }}", "function foo() {}")
}

// Java: PeepholeIntegrationTest.testIteratedRemoval5()
func (c *PeepholeIntegrationTest) TestIteratedRemoval5() {
	c.Test(
		"var x; \n"+
			" out: { \n"+
			"   try { break out; } catch (e) { break out; } \n"+
			"   x = undefined; \n"+
			" }\n",
		"var x;")
}

// Java: PeepholeIntegrationTest.testIssue311()
func (c *PeepholeIntegrationTest) TestIssue311() {
	c.Test(
		"function a(b) {\n"+
			"  switch (b.v) {\n"+
			"    case 'SWITCH':\n"+
			"      if (b.i >= 0) {\n"+
			"        return b.o;\n"+
			"      } else {\n"+
			"        return;\n"+
			"      }\n"+
			"      break;\n"+
			"  }\n"+
			"}\n",
		"function a(b) {\n"+
			"  switch (b.v) {\n"+
			"    case 'SWITCH':\n"+
			"      if (b.i >= 0) {\n"+
			"        return b.o;\n"+
			"      }\n"+
			"  }\n"+
			"}\n")
}

// Java: PeepholeIntegrationTest.testIssue4177428a()
func (c *PeepholeIntegrationTest) TestIssue4177428a() {
	c.TestSame(
		"f = function() {\n" +
			"  var action;\n" +
			"  a: {\n" +
			"    var proto = null;\n" +
			"    try {\n" +
			"      proto = new Proto\n" +
			"    } finally {\n" +
			"      action = proto;\n" +
			"      break a // Keep this...\n" +
			"    }\n" +
			"  }\n" +
			"  alert(action) // and this.\n" +
			"};\n")
}

// Java: PeepholeIntegrationTest.testIssue4177428b()
func (c *PeepholeIntegrationTest) TestIssue4177428b() {
	c.Test(
		"f = function() {\n"+
			"  var action;\n"+
			"  a: {\n"+
			"    var proto = null;\n"+
			"    try {\n"+
			"    try {\n"+
			"      proto = new Proto\n"+
			"    } finally {\n"+
			"      action = proto;\n"+
			"      break a // Keep this...\n"+
			"    }\n"+
			"    } finally {\n"+
			"    }\n"+
			"  }\n"+
			"  alert(action) // and this.\n"+
			"};\n",
		"f = function() {\n"+
			"  var action;\n"+
			"  a: {\n"+
			"    var proto = null;\n"+
			"    try {\n"+
			"      proto = new Proto\n"+
			"    } finally {\n"+
			"      action = proto;\n"+
			"      break a // Keep this...\n"+
			"    }\n"+
			"  }\n"+
			"  alert(action) // and this.\n"+
			"};\n")
}

// Java: PeepholeIntegrationTest.testIssue4177428c()
func (c *PeepholeIntegrationTest) TestIssue4177428c() {
	c.Test(
		"f = function() {\n"+
			"  var action;\n"+
			"  a: {\n"+
			"    var proto = null;\n"+
			"    try {\n"+
			"    } finally {\n"+
			"    try {\n"+
			"      proto = new Proto\n"+
			"    } finally {\n"+
			"      action = proto;\n"+
			"      break a // Keep this...\n"+
			"    }\n"+
			"    }\n"+
			"  }\n"+
			"  alert(action)\n"+
			"// and this.\n"+
			"};\n",
		"f = function() {\n"+
			"  var action;\n"+
			"  a: {\n"+
			"    var proto = null;\n"+
			"    try {\n"+
			"      proto = new Proto\n"+
			"    } finally {\n"+
			"      action = proto;\n"+
			"      break a // Keep this...\n"+
			"    }\n"+
			"  }\n"+
			"  alert(action) // and this.\n"+
			"};\n")
}

// Java: PeepholeIntegrationTest.testIssue4177428_continue()
func (c *PeepholeIntegrationTest) TestIssue4177428_continue() {
	c.Test(
		"f = function() {\n"+
			"  var action;\n"+
			"  a: do {\n"+
			"    var proto = null;\n"+
			"    try {\n"+
			"      proto = new Proto\n"+
			"    } finally {\n"+
			"      action = proto;\n"+
			"      continue a\n"+
			"// Keep this...\n"+
			"    }\n"+
			"  } while(false)\n"+
			"  alert(action)\n"+
			"// and this.\n"+
			"};\n",
		"f = function() {\n"+
			"  var action;\n"+
			"  a: do {\n"+
			"    var proto = null;\n"+
			"    try {\n"+
			"      proto = new Proto\n"+
			"    } finally {\n"+
			"      action = proto;\n"+
			"      continue a\n"+
			"// Keep this...\n"+
			"    }\n"+
			"  } while(0)\n"+
			"  alert(action)\n"+
			"// and this.\n"+
			"};\n")
}

// Java: PeepholeIntegrationTest.testIssue4177428_multifinally()
func (c *PeepholeIntegrationTest) TestIssue4177428_multifinally() {
	c.Test(
		"a: {\n"+
			" try {\n"+
			"   try {\n"+
			"   } finally {\n"+
			"     break a;\n"+
			"   }\n"+
			" } finally {\n"+
			"   x = 1;\n"+
			" }\n"+
			"}\n",
		"a: {\n"+
			"  x = 1;\n"+
			"}\n")
}

// Java: PeepholeIntegrationTest.testIssue5215541_deadVarDeclar()
func (c *PeepholeIntegrationTest) TestIssue5215541_deadVarDeclar() {
	c.Test(
		"       throw 1; var x;", //
		"var x; throw 1;       ")
	c.TestSame("throw 1; function x() {}")
	c.Test(
		"throw 1; var x; var y;                ", //
		"                var y; var x; throw 1;")
	c.Test(
		"       throw 1; var x = foo", //
		"var x; throw 1")
}

// Java: PeepholeIntegrationTest.testForInLoop()
func (c *PeepholeIntegrationTest) TestForInLoop() {
	c.TestSame("var x; for(x in y) {}")
}

// Java: PeepholeIntegrationTest.testDontRemoveBreakInTryFinally()
func (c *PeepholeIntegrationTest) TestDontRemoveBreakInTryFinally() {
	c.TestSame(
		"function f() {\n" +
			"  b: {\n" +
			"    try {\n" +
			"      throw 9;\n" +
			"    } finally {\n" +
			"      break b;\n" +
			"    }\n" +
			"  }\n" +
			"  return 1;\n" +
			"}\n")
}

// Java: PeepholeIntegrationTest.testDontRemoveBreakInTryFinallySwitch()
func (c *PeepholeIntegrationTest) TestDontRemoveBreakInTryFinallySwitch() {
	c.TestSame(
		"function f() {\n" +
			"  b: {\n" +
			"    try {\n" +
			"      throw 9;\n" +
			"    } finally {\n" +
			"      switch (x) {\n" +
			"        case 1:\n" +
			"          break b;\n" +
			"      }\n" +
			"    }\n" +
			"  }\n" +
			"  return 1;\n" +
			"}\n")
}

// This was originally supported by UnreachableCodeElimination, which was removed in favor
// of peephole optimizations. Support this test case if found useful in the real code.
// Java: PeepholeIntegrationTest.testIssue1001()
func (c *PeepholeIntegrationTest) TestIssue1001() {
	c.Test(
		"function f(x) { x.property = 3; } f({})", //
		"function f(x) { x.property = 3; }")
	c.Test(
		"function f(x) { x.property = 3; } new f({})", //
		"function f(x) { x.property = 3; }")
}

// Java: PeepholeIntegrationTest.testLetConstBlocks()
func (c *PeepholeIntegrationTest) TestLetConstBlocks() {
	c.Test(
		"function f() {return 1; let a; }", //
		"function f() {let a; return 1;}")
	c.Test(
		"function f() { return 1; const a = 1; }", //
		"function f() { let a;  return 1;}")
	c.Test(
		"function f() { x = 1; { let g; return x; } let y;}",
		"function f() { let y; x = 1;   let g; return x;         } ")
}

// Java: PeepholeIntegrationTest.testArrowFunctions()
func (c *PeepholeIntegrationTest) TestArrowFunctions() {
	c.Test("f(x => {return x; j = 1})", "f(x => {return x;})")
	c.TestSame("f( () => {return 1;})")
}

// Java: PeepholeIntegrationTest.testForOf()
func (c *PeepholeIntegrationTest) TestForOf() {
	c.Test("for(x of i){ 1; }", "for(x of i) {}")
	c.TestSame("for(x of i){}")
}

// Java: PeepholeIntegrationTest.testLetConstBlocks_inFunction_exportedFromEs6Module()
func (c *PeepholeIntegrationTest) TestLetConstBlocks_inFunction_exportedFromEs6Module() {
	c.Test(
		"function f() {\n"+
			"  return 1;\n"+
			"  let a;\n"+
			"}\n"+
			"export { f as f };\n",
		"function f() {\n"+
			"  let a;\n"+
			"  return 1;\n"+
			"}\n"+
			"export { f as f };\n")
	c.Test(
		"function f() {\n"+
			"  return 1;\n"+
			"  const a = 1;\n"+
			"}\n"+
			"export { f as f };\n",
		"function f() {\n"+
			"  let a;\n"+
			"  return 1;\n"+
			"}\n"+
			"export { f as f };\n")
	c.Test(
		"function f() {\n"+
			"  let z;\n"+
			"  x = 1;\n"+
			"  {\n"+
			"    let g;\n"+
			"    return x\n"+
			"  }\n"+
			"  let y\n"+
			"}\n"+
			"export { f as f };\n",
		"function f() {\n"+
			"  let y;\n"+
			"  let z;\n"+
			"  x = 1;\n"+
			"  let g;\n"+
			"  return x;\n"+
			"}\n"+
			"export { f as f };\n")
}

// Java: PeepholeIntegrationTest.testRemoveUnreachableCode_withES6Modules()
func (c *PeepholeIntegrationTest) TestRemoveUnreachableCode_withES6Modules() {
	// Switch statements
	c.Test(
		"function foo() {\n"+
			"  switch (foo) {\n"+
			"    case 1:\n"+
			"      x = 1;\n"+
			"      return;\n"+
			"      break;\n"+
			"    case 2: {\n"+
			"      x = 2;\n"+
			"      return;\n"+
			"      break;\n"+
			"    }\n"+
			"    default:\n"+
			"  }\n"+
			"}\n"+
			"export { foo as foo };\n",
		"function foo() {\n"+
			"  switch (foo) {\n"+
			"    case 1:\n"+
			"      x = 1;\n"+
			"      break;\n"+
			"    case 2:\n"+
			"      x = 2;\n"+
			"  }\n"+
			"}\n"+
			"export { foo as foo };\n")
	// if/else statements with returns
	c.Test(
		"function bar() {\n"+
			"  if (foo)\n"+
			"    x=1;\n"+
			"  else if(bar) {\n"+
			"    return;\n"+
			"    x=2;\n"+
			"  } else {\n"+
			"    x=3;\n"+
			"    return;\n"+
			"    x=4;\n"+
			"  }\n"+
			"  return 5;\n"+
			"  x=5;\n"+
			"}\n"+
			"export { bar as bar };\n",
		"function bar() {\n"+
			"  if (foo) {\n"+
			"    x=1;\n"+
			"    return 5;\n"+
			"  }\n"+
			"  bar || (x = 3);\n"+
			"}\n"+
			"export { bar as bar };\n")
}

// Java: PeepholeIntegrationTest.testComputedClassPropertyNotRemoved()
func (c *PeepholeIntegrationTest) TestComputedClassPropertyNotRemoved() {
	c.TestSame("class Foo { ['x']() {} }")
}

// Java: PeepholeIntegrationTest.testClassExtendsNotRemoved()
func (c *PeepholeIntegrationTest) TestClassExtendsNotRemoved() {
	c.TestSame(
		"function f() {}\n" +
			"class Foo extends f() {}\n")
}

// Java: PeepholeIntegrationTest.testRemoveUnreachableCodeInComputedPropertIife()
func (c *PeepholeIntegrationTest) TestRemoveUnreachableCodeInComputedPropertIife() {
	c.Test(
		"class Foo {\n"+
			"  [function() {\n"+
			"    1; return 'x';\n"+
			"  }()]() { return 1; }\n"+
			"}\n",
		"class Foo {\n"+
			"  [function() {\n"+
			"    return 'x';\n"+
			"  }()]() { return 1; }\n"+
			"}\n")
}

// Java: PeepholeIntegrationTest.testStaticBlockRemoved()
func (c *PeepholeIntegrationTest) TestStaticBlockRemoved() {
	c.Test("class Foo { static {} }", "class Foo {  }")
}

// Java: PeepholeIntegrationTest.testRemoveUnreachableCodeInStaticBlock1()
func (c *PeepholeIntegrationTest) TestRemoveUnreachableCodeInStaticBlock1() {
	// TODO(b/240443227): Unreachable/Useless code isn't removed in static blocks
	c.Test(
		"class Foo {\n"+
			"  static {\n"+
			"    switch (a) { case 'a': break }\n"+
			"    try {var x = 1} catch (e) {e()}\n"+
			"    true;\n"+
			"    if (x) 1;\n"+
			"  }\n"+
			"}\n",
		"class Foo {\n"+
			"  static {\n"+
			"    try {var x = 1} catch (e) {e()}\n"+
			"  }\n"+
			"}\n")
}

// Java: PeepholeIntegrationTest.testDontFoldDirectAndIndirectEval()
func (c *PeepholeIntegrationTest) TestDontFoldDirectAndIndirectEval() {
	c.Test(
		"function f(a, x) { if (a) { return eval(x); } else { return (0, eval)(x); } }",
		"function f(a, x) { return a ? eval(x) : (0, eval)(x); }")
}
