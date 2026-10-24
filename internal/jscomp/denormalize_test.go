// Port of com/google/javascript/jscomp/DenormalizeTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/DenormalizeTest.java)

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
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"testing"
)

type DenormalizeTest struct {
	CompilerTestCaseBase

	outputFeatureSet *parser.FeatureSet
}

// Java: DenormalizeTest()
func NewDenormalizeTest() *DenormalizeTest {
	c := &DenormalizeTest{}
	c.InitCompilerTestCaseBase(c)
	c.outputFeatureSet = parser.FeatureSet_ES_UNSUPPORTED
	return c
}

func TestDenormalizeTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewDenormalizeTest() })
}

// Java: DenormalizeTest.getProcessor(Compiler)
func (c *DenormalizeTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewDenormalizeTest_NormalizeAndDenormalizePass(compiler, c.outputFeatureSet)
}

// Java: DenormalizeTest.testInlineVarNullishCoalesce()
func (c *DenormalizeTest) TestInlineVarNullishCoalesce() {
	c.Test(
		"function f() {\n"+
			"  var x;\n"+
			"  function g() { x = 0 ?? \"hi\"; }\n"+
			"  if (y) { x = -1 ?? true; }\n"+
			"  alert(x);\n"+
			"}\n",
		"function f() {\n"+
			"  function g() { x = 0 ?? \"hi\"; }\n"+
			"  if (y) { var x = -1 ?? true; }\n"+
			"  alert(x);\n"+
			"}\n")
}

// Java: DenormalizeTest.testInlineVarKeyword1()
func (c *DenormalizeTest) TestInlineVarKeyword1() {
	c.Test(
		"function f() {\n"+
			"  var x;\n"+
			"  function g() { x = 2; }\n"+
			"  if (y) { x = -1; }\n"+
			"  alert(x);\n"+
			"}\n",
		"function f() {\n"+
			"  function g() { x = 2; }\n"+
			"  if (y) { var x = -1; }\n"+
			"  alert(x);\n"+
			"}\n")
}

// Java: DenormalizeTest.testInlineVarKeyword2()
func (c *DenormalizeTest) TestInlineVarKeyword2() {
	c.Test(
		"function f() {\n"+
			"  var x;\n"+
			"  function g() { x = 2; }\n"+
			"  if (y) { x = -1; } else { x = 3; }\n"+
			"  alert(x);\n"+
			"}\n",
		"function f() {\n"+
			"  function g() { x = 2; }\n"+
			"  if (y) { var x = -1; } else { x = 3; }\n"+
			"  alert(x);\n"+
			"}\n")
}

// Java: DenormalizeTest.testInlineVarKeywordArrowFunc1()
func (c *DenormalizeTest) TestInlineVarKeywordArrowFunc1() {
	c.Test(
		"var f = () => {\n"+
			"  var x;\n"+
			"  var g = () => { x = 2; }\n"+
			"  if (y) { x = -1; }\n"+
			"  alert(x);\n"+
			"}\n",
		"var f = () => {\n"+
			"  var g = () => { x = 2; }\n"+
			"  if (y) { var x = -1; }\n"+
			"  alert(x);\n"+
			"}\n")
}

// Java: DenormalizeTest.testInlineVarKeywordArrowFunc2()
func (c *DenormalizeTest) TestInlineVarKeywordArrowFunc2() {
	c.Test(
		"var f = () => {\n"+
			"  var x;\n"+
			"  var g = () => { x = 2; }\n"+
			"  if (y) { x = -1; } else { x = 3; }\n"+
			"  alert(x);\n"+
			"}\n",
		"var f = () => {\n"+
			"  var g = () => { x = 2; }\n"+
			"  if (y) { var x = -1; } else { x = 3; }\n"+
			"  alert(x);\n"+
			"}\n")
}

// Java: DenormalizeTest.testNotInlineConstLet()
func (c *DenormalizeTest) TestNotInlineConstLet() {
	c.TestSame(
		"let x;\n" +
			"if (y) { x = -1; }\n")

	c.TestSame(
		"const x = 1;\n" +
			"if (y) { x = -1; }\n")
}

// Java: DenormalizeTest.testFor()
func (c *DenormalizeTest) TestFor() {
	// Verify assignments are moved into the FOR init node.
	c.Test("a = 0; for(; a < 2 ; a++) foo()", "for(a = 0; a < 2 ; a++) foo();")
	// Verify vars are are moved into the FOR init node.
	c.Test("var a = 0; for(; c < b ; c++) foo()", "for(var a = 0; c < b ; c++) foo()")

	// We don't handle labels yet.
	c.TestSame("var a = 0; a:for(; c < b ; c++) foo()")
	c.TestSame("var a = 0; a:b:for(; c < b ; c++) foo()")

	// Do not inline let or const
	c.TestSame("let a = 0; for(; c < b ; c++) foo()")
	c.TestSame("const a = 0; for(; c < b ; c++) foo()")

	// Verify FOR inside IFs.
	c.Test("if(x){var a = 0; for(; c < b; c++) foo()}", "if(x){for(var a = 0; c < b; c++) foo()}")

	// Any other expression.
	c.Test("init(); for(; a < 2 ; a++) foo()", "for(init(); a < 2 ; a++) foo();")

	// Other statements are left as is.
	c.Test(
		"function f(){ var a; for(; a < 2 ; a++) foo() }",
		"function f(){ for(var a; a < 2 ; a++) foo() }")
	c.TestSame("function f(){ return; for(; a < 2 ; a++) foo() }")

	// Verify destructuring assignments are moved.
	c.Test(
		"[a, b] = [1, 2]; for (; a < 2; a = b++) foo();",
		"for ([a, b] = [1, 2]; a < 2; a = b++) foo();")

	c.Test(
		"var [a, b] = [1, 2]; for (; a < 2; a = b++) foo();",
		"var a; var b; for ([a, b] = [1, 2]; a < 2; a = b++) foo();")
}

// Java: DenormalizeTest.testForIn()
func (c *DenormalizeTest) TestForIn() {
	c.Test("var a; for(a in b) foo()", "for (var a in b) foo()")
	c.TestSame("a = 0; for(a in b) foo()")
	c.TestSame("var a = 0; for(a in b) foo()")

	// We don't handle labels yet.
	c.TestSame("var a; a:for(a in b) foo()")
	c.TestSame("var a; a:b:for(a in b) foo()")

	// Verify FOR inside IFs.
	c.Test("if(x){var a; for(a in b) foo()}", "if(x){for(var a in b) foo()}")

	// Any other expression.
	c.TestSame("init(); for(a in b) foo()")

	// Other statements are left as is.
	c.TestSame("function f(){ return; for(a in b) foo() }")

	// We don't handle destructuring patterns yet.
	c.TestSame("var a; var b; for ([a, b] in c) foo();")
}

// Java: DenormalizeTest.testForOf()
func (c *DenormalizeTest) TestForOf() {
	c.Test("var a; for (a of b) foo()", "for (var a of b) foo()")
	c.TestSame("a = 0; for (a of b) foo()")
	c.TestSame("var a = 0; for (a of b) foo()")

	// We don't handle labels yet.
	c.TestSame("var a; a: for (a of b) foo()")
	c.TestSame("var a; a: b: for (a of b) foo()")

	// Verify FOR inside IFs.
	c.Test("if (x) { var a; for (a of b) foo() }", "if (x) { for (var a of b) foo() }")

	// Any other expression.
	c.TestSame("init(); for (a of b) foo()")

	// Other statements are left as is.
	c.TestSame("function f() { return; for (a of b) foo() }")

	// We don't handle destructuring patterns yet.
	c.TestSame("var a; var b; for ([a, b] of c) foo();")
}

// Java: DenormalizeTest.testInOperatorNotInsideFor()
func (c *DenormalizeTest) TestInOperatorNotInsideFor() {
	// in operators shouldn't be moved into for loops.
	// Some JavaScript interpreters (such as the NetFront Access browser
	// embedded in the PlayStation 3) will not parse an in operator in
	// a for loop, even if it's protected by parentheses.

	// Make sure the in operator doesn't get moved into the for loop.
	c.TestSame(
		"function f(){ var a; var i=\"length\" in a;\n" +
			"for(; a < 2 ; a++) foo() }\n")
	// Same, but with parens around the operator.
	c.TestSame(
		"function f(){ var a; var i=(\"length\" in a);\n" +
			"for(; a < 2 ; a++) foo() }\n")
	// Make sure Normalize yanks the variable initializer out, and
	// Denormalize doesn't put it back.
	c.Test(
		"function f(){\n"+
			"var b,a=0; for (var i=(\"length\" in b);a<2; a++) foo()}\n",
		"function f(){var b; var a=0;var i=(\"length\" in b);\n"+
			"for (;a<2;a++) foo()}\n")
}

// Java: DenormalizeTest.testAssignShorthand()
func (c *DenormalizeTest) TestAssignShorthand() {
	c.Test("x = x | 1;", "x |= 1;")
	c.Test("x = x ^ 1;", "x ^= 1;")
	c.Test("x = x & 1;", "x &= 1;")
	c.Test("x = x << 1;", "x <<= 1;")
	c.Test("x = x >> 1;", "x >>= 1;")
	c.Test("x = x >>> 1;", "x >>>= 1;")
	c.Test("x = x + 1;", "x += 1;")
	c.Test("x = x - 1;", "x -= 1;")
	c.Test("x = x * 1;", "x *= 1;")
	c.Test("x = x / 1;", "x /= 1;")
	c.Test("x = x % 1;", "x %= 1;")

	c.Test("/** @suppress {const} */ x = x + 1;", "/** @suppress {const} */ x += 1;")
}

// Java: DenormalizeTest.testAssignShorthandDontDenormalizeWhenLHSNotName()
func (c *DenormalizeTest) TestAssignShorthandDontDenormalizeWhenLHSNotName() {
	c.TestSame("obj.x = obj.x + 1;")
}

// Java: DenormalizeTest.testLogicalAssignShorthand()
func (c *DenormalizeTest) TestLogicalAssignShorthand() {
	c.Test("x || (x = 1);", "x ||= 1;")
	c.Test("x && (x = 1);", "x &&= 1;")
	c.Test("x ?? (x = 1);", "x ??= 1;")
}

// Java: DenormalizeTest.testLogicalAssignShorthandDontDenormalizeWhenLHSNotName()
func (c *DenormalizeTest) TestLogicalAssignShorthandDontDenormalizeWhenLHSNotName() {
	c.TestSame("obj.x || (obj.x = 1);")
}

// Java: DenormalizeTest.testLogicalAssignShorthandRHSAssignmentOp()
func (c *DenormalizeTest) TestLogicalAssignShorthandRHSAssignmentOp() {
	// This test was intended to show that the following should not fold into `a ||= 1`
	// Both Normalize and Denormalize report to the compiler they have modified this code
	// (Normalize produces `a || (a = a + 1);` and Denormalize reverts it to `a || (a += 1);`).
	// The unit test wants to warn on unnecessary change reporting,
	// so validating AST change marking has been disabled.
	c.DisableValidateAstChangeMarking()
	c.TestSame("a || (a += 1);")
}

// Java: DenormalizeTest.testLogicalAssignShorthandBrowser2020()
func (c *DenormalizeTest) TestLogicalAssignShorthandBrowser2020() {
	c.outputFeatureSet = parser.FeatureSet_BROWSER_2020
	c.TestSame("x || (x = 1);")
	c.TestSame("x && (x = 1);")
	c.TestSame("x ?? (x = 1);")
}

// Java: DenormalizeTest.testNoCrashOnEs6Features()
func (c *DenormalizeTest) TestNoCrashOnEs6Features() {
	c.Test(
		"class C {\n"+
			"  constructor() {\n"+
			"    var x;\n"+
			"    if (y) { x = -1; }\n"+
			"  }\n"+
			"}\n",
		"class C {\n"+
			"  constructor() {\n"+
			"    if (y) { var x = -1; }\n"+
			"  }\n"+
			"}\n")

	c.Test(
		"var obj = {\n"+
			"  method() {\n"+
			"    var c; for (; c < b ; c++) foo()\n"+
			"  },\n"+
			"}\n",
		"var obj = {\n"+
			"  method() {\n"+
			"    for (var c; c < b ; c++) foo()\n"+
			"  },\n"+
			"}\n")

	c.TestSame(
		"var obj = {\n" +
			"  ['computed' + 'prop']: 42\n" +
			"}\n")

	// Denormalize does not revert shorthand object literals that were expanded in Normalize
	c.Test(
		"var obj = {\n"+
			"  key\n"+
			"}\n",
		"var obj = {\n"+
			"  key: key\n"+
			"}\n")

	c.Test(
		"function tag(strings) {\n"+
			"  var x;\n"+
			"  if (y) { x = x + 1; }\n"+
			"}\n"+
			"tag`template`\n",
		"function tag(strings) {\n"+
			"  var x;\n"+
			"  if (y) { x += 1; }\n"+
			"}\n"+
			"tag`template`\n")

	c.TestSame(
		"var x;\n" +
			"var y;\n" +
			"if (y) { [x, y] = [1, 2]; }\n")
}

/**
 * Create a class to combine the Normalize and Denormalize passes. This is needed because the
 * enableNormalize() call on CompilerTestCase causes normalization of the result *and* the
 * expected string, and we really don't want the compiler twisting the expected code around.
 */
// DenormalizeTest_NormalizeAndDenormalizePass: create a class to combine
// the Normalize and Denormalize passes. This is needed because the
// enableNormalize() call on CompilerTestCase causes normalization of the
// result *and* the expected string, and we really don't want the compiler
// twisting the expected code around.
type DenormalizeTest_NormalizeAndDenormalizePass struct {
	denormalizePass *Denormalize
	normalizePass   *Normalize_NormalizeStatements
	compiler        AbstractCompiler
}

// Java: DenormalizeTest.NormalizeAndDenormalizePass(AbstractCompiler, FeatureSet)
func NewDenormalizeTest_NormalizeAndDenormalizePass(compiler AbstractCompiler, outputFeatureSet *parser.FeatureSet) *DenormalizeTest_NormalizeAndDenormalizePass {
	p := &DenormalizeTest_NormalizeAndDenormalizePass{}
	p.compiler = compiler
	p.denormalizePass = NewDenormalize(compiler, outputFeatureSet)
	p.normalizePass = NewNormalize_NormalizeStatements(compiler, false, nil)
	return p
}

// Java: DenormalizeTest.NormalizeAndDenormalizePass.process(Node, Node)
func (p *DenormalizeTest_NormalizeAndDenormalizePass) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(p.compiler, root, p.normalizePass)
	p.denormalizePass.Process(externs, root)
}
