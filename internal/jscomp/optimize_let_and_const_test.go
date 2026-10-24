// Port of com/google/javascript/jscomp/OptimizeLetAndConstTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/OptimizeLetAndConstTest.java)

/*
 * Copyright 2026 The Closure Compiler Authors.
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

// OptimizeLetAndConstTest is the tests for OptimizeLetAndConst.
type OptimizeLetAndConstTest struct {
	CompilerTestCaseBase

	assumeOutputIsWrapped bool
}

// Java: OptimizeLetAndConstTest()
func NewOptimizeLetAndConstTest() *OptimizeLetAndConstTest {
	c := &OptimizeLetAndConstTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestOptimizeLetAndConstTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewOptimizeLetAndConstTest() })
}

// Java: OptimizeLetAndConstTest.setUp()
func (c *OptimizeLetAndConstTest) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.assumeOutputIsWrapped = false
}

// Java: OptimizeLetAndConstTest.getProcessor(Compiler)
func (c *OptimizeLetAndConstTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewPeepholeOptimizationsPass(
		compiler,
		"testOptimizeLetAndConst",
		NewOptimizeLetAndConstPeephole(c.assumeOutputIsWrapped))
}

// Java: OptimizeLetAndConstTest.testConstToLetInBlock()
func (c *OptimizeLetAndConstTest) TestConstToLetInBlock() {
	c.Test("if (true) { const x = 1; }", "if (true) { let x = 1; }")
}

// Java: OptimizeLetAndConstTest.testConstToVarInFunction()
func (c *OptimizeLetAndConstTest) TestConstToVarInFunction() {
	c.Test("function f() { const x = 1; }", "function f() { var x = 1; }")
}

// Java: OptimizeLetAndConstTest.testLetToVarInFunction()
func (c *OptimizeLetAndConstTest) TestLetToVarInFunction() {
	c.Test("function f() { let x = 1; }", "function f() { var x = 1; }")
}

// Java: OptimizeLetAndConstTest.testLetToVarInBlock()
func (c *OptimizeLetAndConstTest) TestLetToVarInBlock() {
	c.TestSame("function f() { if (true) { let x = 1; } }")
}

// Java: OptimizeLetAndConstTest.testConstToLetGlobalUnwrapped()
func (c *OptimizeLetAndConstTest) TestConstToLetGlobalUnwrapped() {
	c.assumeOutputIsWrapped = false
	c.TestSame("const x = 1;")
}

// Java: OptimizeLetAndConstTest.testConstToVarGlobalWrapped()
func (c *OptimizeLetAndConstTest) TestConstToVarGlobalWrapped() {
	c.assumeOutputIsWrapped = true
	c.Test("const x = 1;", "var x = 1;")
}

// Java: OptimizeLetAndConstTest.testLetToVarGlobalWrapped()
func (c *OptimizeLetAndConstTest) TestLetToVarGlobalWrapped() {
	c.assumeOutputIsWrapped = true
	c.Test("let x = 1;", "var x = 1;")
}

// Java: OptimizeLetAndConstTest.testLetInLoopNotConvertedToVar()
func (c *OptimizeLetAndConstTest) TestLetInLoopNotConvertedToVar() {
	// let in loop header and loop body stay let
	c.TestSame("function f() { for (let i = 0; i < 10; i++) { let x = 1; } }")
}

// Java: OptimizeLetAndConstTest.testConstToLetInLoopConvertedToLetButNotVar()
func (c *OptimizeLetAndConstTest) TestConstToLetInLoopConvertedToLetButNotVar() {
	c.Test(
		"function f() { for (let i = 0; i < 10; i++) { const x = 1; } }",
		"function f() { for (let i = 0; i < 10; i++) { let x = 1; } }")
}

// Java: OptimizeLetAndConstTest.testDestructuringConstToLet()
func (c *OptimizeLetAndConstTest) TestDestructuringConstToLet() {
	c.Test("if (true) { const [a, b] = [1, 2]; }", "if (true) { let [a, b] = [1, 2]; }")
}

// Java: OptimizeLetAndConstTest.testDestructuringLetToVarInFunction()
func (c *OptimizeLetAndConstTest) TestDestructuringLetToVarInFunction() {
	c.Test("function f() { let [a, b] = [1, 2]; }", "function f() { var [a, b] = [1, 2]; }")
}

// Java: OptimizeLetAndConstTest.testClassStaticBlockLetAndConst()
func (c *OptimizeLetAndConstTest) TestClassStaticBlockLetAndConst() {
	c.Test(
		"class C { static { const x = 1; let y = 2; } }",
		"class C { static { var x = 1; var y = 2; } }")
}

// Java: OptimizeLetAndConstTest.testCatchBlockConstAndLet()
func (c *OptimizeLetAndConstTest) TestCatchBlockConstAndLet() {
	c.Test(
		"function f() { try {} catch (e) { const x = 1; let y = 2; } }",
		"function f() { try {} catch (e) { let x = 1; let y = 2; } }")
}

// Java: OptimizeLetAndConstTest.testSwitchConstAndLet()
func (c *OptimizeLetAndConstTest) TestSwitchConstAndLet() {
	c.Test(
		"function f(x) { switch (x) { case 1: const y = 1; let z = 2; } }",
		"function f(x) { switch (x) { case 1: let y = 1; let z = 2; } }")
}
