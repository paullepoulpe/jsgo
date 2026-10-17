// Port of com/google/javascript/jscomp/CollapseVariableDeclarationsTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/CollapseVariableDeclarationsTest.java)

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
	"testing"
)

type CollapseVariableDeclarationsTest struct {
	CompilerTestCaseBase
}

// Java: CollapseVariableDeclarationsTest()
func NewCollapseVariableDeclarationsTest() *CollapseVariableDeclarationsTest {
	c := &CollapseVariableDeclarationsTest{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func TestCollapseVariableDeclarationsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return NewCollapseVariableDeclarationsTest() })
}

// Java: CollapseVariableDeclarationsTest.testCollapsing()
func (c *CollapseVariableDeclarationsTest) TestCollapsing() {
	// Basic collapsing
	c.Test("var a;var b;", "var a,b;")

	// With initial values
	c.Test("var a = 1;var b = 1;", "var a=1,b=1;")

	// Already collapsed
	c.TestSame("var a, b;")

	// Already collapsed with values
	c.TestSame("var a = 1, b = 1;")

	// Some already collapsed
	c.Test("var a;var b, c;var d;", "var a,b,c,d;")

	// Some already collapsed with values
	c.Test("var a = 1;var b = 2, c = 3;var d = 4;", "var a=1,b=2,c=3,d=4;")

	c.Test(
		"var x = 2; foo(x); x = 3; x = 1; var y = 2; var z = 4; x = 5",
		"var x = 2; foo(x); x = 3; x = 1; var y = 2, z = 4; x = 5")
}

// Java: CollapseVariableDeclarationsTest.testIssue820()
func (c *CollapseVariableDeclarationsTest) TestIssue820() {
	// Don't redeclare function parameters, this is incompatible with
	// strict mode.
	c.TestSame("function f(a){ var b=1; a=2; var c; }")
}

// Java: CollapseVariableDeclarationsTest.testIfElseVarDeclarations()
func (c *CollapseVariableDeclarationsTest) TestIfElseVarDeclarations() {
	c.TestSame("if (x) var a = 1; else var b = 2;")
}

// Java: CollapseVariableDeclarationsTest.testAggressiveRedeclarationInFor()
func (c *CollapseVariableDeclarationsTest) TestAggressiveRedeclarationInFor() {
	c.TestSame("for(var x = 1; x = 2; x = 3) {x = 4}")
	c.TestSame("for(var x = 1; y = 2; z = 3) {var a = 4}")
	c.TestSame("var x; for(x = 1; x = 2; z = 3) {x = 4}")
}

// Java: CollapseVariableDeclarationsTest.testIssue397()
func (c *CollapseVariableDeclarationsTest) TestIssue397() {
	c.Test("var x; var y = 3; x = 5;", "var x, y = 3; x = 5;")

	c.TestSame("var x; x = 5; var z = 7;")

	c.Test("var x; var y = 3; x = 5; var z = 7;", "var x, y = 3; x = 5; var z = 7;")

	c.Test("var a = 1; var x; var y = 3; x = 5;", "var a = 1, x, y = 3; x = 5;")
}

// Java: CollapseVariableDeclarationsTest.testArgumentsAssignment()
func (c *CollapseVariableDeclarationsTest) TestArgumentsAssignment() {
	c.TestSame("function f() {arguments = 1;}")
}

// ES6 Tests
// Java: CollapseVariableDeclarationsTest.testCollapsingLetConst()
func (c *CollapseVariableDeclarationsTest) TestCollapsingLetConst() {
	// Basic collapsing
	c.Test("let a;let b;", "let a,b;")

	// With initial values
	c.Test("const a = 1;const b = 1;", "const a=1,b=1;")

	// Already collapsed
	c.TestSame("let a, b;")

	// Already collapsed with values
	c.TestSame("let a = 1, b = 1;")

	// Some already collapsed
	c.Test("let a;let b, c;let d;", "let a,b,c,d;")

	// Some already collapsed with values
	c.Test("let a = 1;let b = 2, c = 3;let d = 4;", "let a=1,b=2,c=3,d=4;")

	// Different variable types
	c.TestSame("let a = 1; const b = 2;")
}

// Java: CollapseVariableDeclarationsTest.testIfElseVarDeclarationsLet()
func (c *CollapseVariableDeclarationsTest) TestIfElseVarDeclarationsLet() {
	c.TestSame("if (x) { let a = 1; } else { let b = 2; }")
}

// Java: CollapseVariableDeclarationsTest.testAggressiveRedeclarationOfLetInFor()
func (c *CollapseVariableDeclarationsTest) TestAggressiveRedeclarationOfLetInFor() {
	c.TestSame("for(let x = 1; x = 2; x = 3) {x = 4}")
	c.TestSame("for(let x = 1; y = 2; z = 3) {let a = 4}")
	c.TestSame("let x; for(x = 1; x = 2; z = 3) {x = 4}")
}

// Java: CollapseVariableDeclarationsTest.testRedeclarationLetInFunction()
func (c *CollapseVariableDeclarationsTest) TestRedeclarationLetInFunction() {
	c.Test(
		"function f() { let x = 1; let y = 2; let z = 3; x + y + z; }",
		"function f() { let x = 1, y = 2, z = 3; x + y + z; } ")

	// recognize local scope version of x
	c.Test(
		"var x = 1; function f() { let x = 1; let y = 2; x + y; }",
		"var x = 1; function f() { let x = 1, y = 2; x + y } ")

	// do not redeclare function parameters
	// incompatible with strict mode
	c.TestSame("function f(x) { let y = 3; x = 4; x + y; }")
}

// Java: CollapseVariableDeclarationsTest.testArrowFunction()
func (c *CollapseVariableDeclarationsTest) TestArrowFunction() {
	c.Test("() => {let x = 1; let y = 2; x + y; }", "() => {let x = 1, y = 2; x + y; }")

	// do not redeclare function parameters
	// incompatible with strict mode
	c.TestSame("(x) => {x = 4; let y = 2; x + y; }")
}

// Java: CollapseVariableDeclarationsTest.testUncollapsableDeclarations()
func (c *CollapseVariableDeclarationsTest) TestUncollapsableDeclarations() {
	c.TestSame("let x = 1; var y = 2; const z = 3")

	c.TestSame("let x = 1; var y = 2; let z = 3;")
}

// Java: CollapseVariableDeclarationsTest.testMixedDeclarationTypes()
func (c *CollapseVariableDeclarationsTest) TestMixedDeclarationTypes() {
	// lets, vars, const declarations consecutive

	c.Test("let x = 1; let z = 3; var y = 2;", "let x = 1, z = 3; var y = 2;")

	c.Test("let x = 1; let y = 2; var z = 3; var a = 4;", "let x = 1, y = 2; var z = 3, a = 4")
}

// Java: CollapseVariableDeclarationsTest.getProcessor(Compiler)
func (c *CollapseVariableDeclarationsTest) GetProcessor(compiler *Compiler) CompilerPass {
	return NewCollapseVariableDeclarations(compiler)
}
