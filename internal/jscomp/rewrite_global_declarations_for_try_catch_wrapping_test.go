// Port of com/google/javascript/jscomp/RewriteGlobalDeclarationsForTryCatchWrappingTest.java @ closure-compiler v20261015
// (test/com/google/javascript/jscomp/RewriteGlobalDeclarationsForTryCatchWrappingTest.java)

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
	"testing"
)

type RewriteGlobalDeclarationsForTryCatchWrappingTest struct {
	CompilerTestCaseBase
}

// Java: RewriteGlobalDeclarationsForTryCatchWrappingTest()
func NewRewriteGlobalDeclarationsForTryCatchWrappingTest() *RewriteGlobalDeclarationsForTryCatchWrappingTest {
	c := &RewriteGlobalDeclarationsForTryCatchWrappingTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestRewriteGlobalDeclarationsForTryCatchWrappingTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewRewriteGlobalDeclarationsForTryCatchWrappingTest() })
}

// Java: RewriteGlobalDeclarationsForTryCatchWrappingTest.getProcessor(Compiler)
func (c *RewriteGlobalDeclarationsForTryCatchWrappingTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewRewriteGlobalDeclarationsForTryCatchWrapping(compiler)
}

// Java: RewriteGlobalDeclarationsForTryCatchWrappingTest.testFunctionDeclarations()
func (c *RewriteGlobalDeclarationsForTryCatchWrappingTest) TestFunctionDeclarations() {
	c.Test("a; function f(){} function g(){}", "var f = function(){}; var g = function(){}; a;")
}

// Java: RewriteGlobalDeclarationsForTryCatchWrappingTest.testFunctionDeclarationsInBlock()
func (c *RewriteGlobalDeclarationsForTryCatchWrappingTest) TestFunctionDeclarationsInBlock() {
	c.TestSame("a; if (1) { function f(){} function g(){} }")
}

// Java: RewriteGlobalDeclarationsForTryCatchWrappingTest.testClassDeclarations()
func (c *RewriteGlobalDeclarationsForTryCatchWrappingTest) TestClassDeclarations() {
	// classes don't hoist
	c.Test("a; class B{} class C{}", "a; var B = class{}; var C = class{};")
}

// Java: RewriteGlobalDeclarationsForTryCatchWrappingTest.testClassDeclarationsInBlock()
func (c *RewriteGlobalDeclarationsForTryCatchWrappingTest) TestClassDeclarationsInBlock() {
	c.TestSame("a; if (1) { class B{} class C{} }")
}

// Java: RewriteGlobalDeclarationsForTryCatchWrappingTest.testFunctionDeclarationsInModule()
func (c *RewriteGlobalDeclarationsForTryCatchWrappingTest) TestFunctionDeclarationsInModule() {
	c.Test(
		c.Srcs(
			JSChunkGraphBuilder_ForUnordered().
				AddChunk("a; function f(){} function g(){}").
				Build()),
		c.Expected("var f = function(){}; var g = function(){}; a"))
}

// Java: RewriteGlobalDeclarationsForTryCatchWrappingTest.testGeneratorDeclarations()
func (c *RewriteGlobalDeclarationsForTryCatchWrappingTest) TestGeneratorDeclarations() {
	c.Test(
		"a; function *f(){} function *g(){}", "var f = function* (){}; var g = function* (){}; a;")
}

// Java: RewriteGlobalDeclarationsForTryCatchWrappingTest.testFunctionDeclarationsInEs6Module()
func (c *RewriteGlobalDeclarationsForTryCatchWrappingTest) TestFunctionDeclarationsInEs6Module() {
	// NOTE: Currently this pass always runs after module transpilation.
	// No current uses of this pass would benefit from ES6 module support:
	// a) RescopeGlobalSymbols, which relies on this pass, only cares about the global scope.
	// b) ES6 module output cannot be wrapped in a try/catch.
	c.TestSame("a; function f(){} function g(){} export default 5;")
}

// Java: RewriteGlobalDeclarationsForTryCatchWrappingTest.testFunctionsExpression()
func (c *RewriteGlobalDeclarationsForTryCatchWrappingTest) TestFunctionsExpression() {
	c.TestSame("a; f = function(){}")
}

// Java: RewriteGlobalDeclarationsForTryCatchWrappingTest.testNoMoveDeepFunctionDeclarations()
func (c *RewriteGlobalDeclarationsForTryCatchWrappingTest) TestNoMoveDeepFunctionDeclarations() {
	c.TestSame("a; if (a) function f(){};")
	c.TestSame("a; if (a) { function f(){} }")
}

// Java: RewriteGlobalDeclarationsForTryCatchWrappingTest.testMigrateLetConstToVar1()
func (c *RewriteGlobalDeclarationsForTryCatchWrappingTest) TestMigrateLetConstToVar1() {
	c.Test("let a = 1;", "var a = 1")
	c.Test("const a = 2;", "var a = 2")
}

// Java: RewriteGlobalDeclarationsForTryCatchWrappingTest.testMigrateLetConstToVar2()
func (c *RewriteGlobalDeclarationsForTryCatchWrappingTest) TestMigrateLetConstToVar2() {
	c.TestSame("{let a = 1;}")
	c.TestSame("{const a = 2;}")
}
