// Port of com/google/javascript/jscomp/NodeIteratorsTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/NodeIteratorsTest.java)

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
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// NodeIteratorsTest holds tests for NodeIterators.
//
// In each test, we find the declaration of "X" in the local scope,
// construct a list of all nodes where X is guaranteed to retain its
// original value, and compare those nodes against an expected list of
// tokens.
type NodeIteratorsTest struct{ t *testing.T }

func (c *NodeIteratorsTest) JUnitSetT(t *testing.T) { c.t = t }
func (c *NodeIteratorsTest) JUnitBefore()           {}
func (c *NodeIteratorsTest) JUnitAfter()            {}

func TestNodeIteratorsTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return &NodeIteratorsTest{} })
}

// Java: NodeIteratorsTest.testBasic()
func (c *NodeIteratorsTest) TestBasic() {
	c.testVarMotionWithCode("var X = 3;", rhino.Token_VAR, rhino.Token_SCRIPT)
}

// Java: NodeIteratorsTest.testNamedFunction()
func (c *NodeIteratorsTest) TestNamedFunction() {
	c.testVarMotionWithCode("var X = 3; function f() {}", rhino.Token_VAR, rhino.Token_SCRIPT)
}

// Java: NodeIteratorsTest.testNamedFunction2()
func (c *NodeIteratorsTest) TestNamedFunction2() {
	c.testVarMotionWithCode(
		"var X = 3; function f() {} var Y;", rhino.Token_VAR, rhino.Token_NAME, rhino.Token_VAR, rhino.Token_SCRIPT)
}

// Java: NodeIteratorsTest.testFunctionExpression()
func (c *NodeIteratorsTest) TestFunctionExpression() {
	c.testVarMotionWithCode(
		"var X = 3, Y = function() {}; 3;",
		rhino.Token_NAME,
		rhino.Token_VAR,
		rhino.Token_NUMBER,
		rhino.Token_EXPR_RESULT,
		rhino.Token_SCRIPT)
}

// Java: NodeIteratorsTest.testFunctionExpression2()
func (c *NodeIteratorsTest) TestFunctionExpression2() {
	c.testVarMotionWithCode(
		"var X = 3; var Y = function() {}; 3;",
		rhino.Token_VAR,
		rhino.Token_NAME,
		rhino.Token_VAR,
		rhino.Token_NUMBER,
		rhino.Token_EXPR_RESULT,
		rhino.Token_SCRIPT)
}

// Java: NodeIteratorsTest.testHaltAtVarRef()
func (c *NodeIteratorsTest) TestHaltAtVarRef() {
	c.testVarMotionWithCode(
		"var X, Y = 3; var Z = X;", rhino.Token_NUMBER, rhino.Token_NAME, rhino.Token_VAR, rhino.Token_NAME)
}

// Java: NodeIteratorsTest.testHaltAtVarRef2()
func (c *NodeIteratorsTest) TestHaltAtVarRef2() {
	c.testVarMotionWithCode(
		"var X, Y = 3; (function() {})(3, X);",
		rhino.Token_NUMBER,
		rhino.Token_NAME,
		rhino.Token_VAR,
		rhino.Token_NUMBER,
		rhino.Token_NAME)
}

// Java: NodeIteratorsTest.testHaltAtVarRef3()
func (c *NodeIteratorsTest) TestHaltAtVarRef3() {
	c.testVarMotionWithCode("var X, Y = 3; X;", rhino.Token_NUMBER, rhino.Token_NAME, rhino.Token_VAR, rhino.Token_NAME)
}

// Java: NodeIteratorsTest.testHaltAtSideEffects()
func (c *NodeIteratorsTest) TestHaltAtSideEffects() {
	c.testVarMotionWithCode(
		"var X, Y = 3; var Z = B(3);",
		rhino.Token_NUMBER,
		rhino.Token_NAME,
		rhino.Token_VAR,
		rhino.Token_NAME,
		rhino.Token_NUMBER)
}

// Java: NodeIteratorsTest.testHaltAtSideEffects2()
func (c *NodeIteratorsTest) TestHaltAtSideEffects2() {
	c.testVarMotionWithCode(
		"var A = 1, X = A, Y = 3; delete A;", rhino.Token_NUMBER, rhino.Token_NAME, rhino.Token_VAR, rhino.Token_NAME)
}

// Java: NodeIteratorsTest.testHaltAtSideEffects3()
func (c *NodeIteratorsTest) TestHaltAtSideEffects3() {
	c.testVarMotionWithCode(
		"var A = 1, X = A, Y = 3; A++;", rhino.Token_NUMBER, rhino.Token_NAME, rhino.Token_VAR, rhino.Token_NAME)
}

// Java: NodeIteratorsTest.testHaltAtSideEffects4()
func (c *NodeIteratorsTest) TestHaltAtSideEffects4() {
	c.testVarMotionWithCode(
		"var A = 1, X = A, Y = 3; A--;", rhino.Token_NUMBER, rhino.Token_NAME, rhino.Token_VAR, rhino.Token_NAME)
}

// Java: NodeIteratorsTest.testHaltAtSideEffects5()
func (c *NodeIteratorsTest) TestHaltAtSideEffects5() {
	c.testVarMotionWithCode(
		"var A = 1, X = A, Y = 3; A = 'a';",
		rhino.Token_NUMBER,
		rhino.Token_NAME,
		rhino.Token_VAR,
		rhino.Token_NAME,
		rhino.Token_STRINGLIT)
}

// Java: NodeIteratorsTest.testNoHaltReadWhenValueIsImmutable()
func (c *NodeIteratorsTest) TestNoHaltReadWhenValueIsImmutable() {
	c.testVarMotionWithCode(
		"var X = 1, Y = 3; alert();", rhino.Token_NUMBER, rhino.Token_NAME, rhino.Token_VAR, rhino.Token_NAME)
}

// Java: NodeIteratorsTest.testHaltReadWhenValueHasSideEffects()
func (c *NodeIteratorsTest) TestHaltReadWhenValueHasSideEffects() {
	c.testVarMotionWithCode("var X = f(), Y = 3; alert();", rhino.Token_NUMBER, rhino.Token_NAME, rhino.Token_VAR)
}

// Java: NodeIteratorsTest.testCatchBlock()
func (c *NodeIteratorsTest) TestCatchBlock() {
	c.testVarMotionWithCode(
		"var X = 1; try { 4; } catch (X) {}",
		rhino.Token_VAR,
		rhino.Token_NUMBER,
		rhino.Token_EXPR_RESULT,
		rhino.Token_BLOCK)
}

// Java: NodeIteratorsTest.testIfBranch()
func (c *NodeIteratorsTest) TestIfBranch() {
	c.testVarMotionWithCode("var X = foo(); if (X) {}", rhino.Token_VAR, rhino.Token_NAME)
}

// Java: NodeIteratorsTest.testLet()
func (c *NodeIteratorsTest) TestLet() {
	c.testVarMotionWithCode("let X = foo(); if (X) {}", rhino.Token_LET, rhino.Token_NAME)
}

// Java: NodeIteratorsTest.testConst()
func (c *NodeIteratorsTest) TestConst() {
	c.testVarMotionWithCode("const X = foo(); if (X) {}", rhino.Token_CONST, rhino.Token_NAME)
}

// testVarMotionWithCode parses the given code, finds the variable X in the
// global scope, and runs the VarMotion iterator. Asserts that the iteration
// order matches the tokens given.
//
// Java: NodeIteratorsTest.testVarMotionWithCode(String, Token...)
func (c *NodeIteratorsTest) testVarMotionWithCode(code string, expectedTokens ...rhino.Token) {
	expectedList := expectedTokens
	c.testVarMotionWithCodeOfList(code, expectedList)
}

// Java: NodeIteratorsTest.testVarMotionWithCode(String, List<Token>)
func (c *NodeIteratorsTest) testVarMotionWithCodeOfList(code string, expectedTokens []rhino.Token) {
	var ancestors []*rhino.Node

	// Add an empty node to the beginning of the code and start there.
	compiler := NewCompiler()
	root := compiler.ParseTestCode(";" + code)
	for n := root; n != nil; n = n.GetFirstChild() {
		ancestors = append([]*rhino.Node{n}, ancestors...)
	}

	searchIt := NewNodeIterators_FunctionlessLocalScope(ancestors...)

	found := false
	for searchIt.HasNext() {
		n := searchIt.Next()
		if n.IsName() &&
			NodeUtil_IsNameDeclaration(searchIt.CurrentParent()) &&
			n.GetString() == "X" {
			found = true
			break
		}
	}

	AssertWithMessage("Variable X not found! %s", root.ToStringTree()).That(found).IsTrue()

	currentAncestors := searchIt.CurrentAncestors()
	AssertThat(currentAncestors.Size()).IsAtLeast(3)
	moveIt := NodeIterators_LocalVarMotion_ForVar(
		compiler, currentAncestors.Get(0), currentAncestors.Get(1), currentAncestors.Get(2))
	var actualTokens []rhino.Token
	for moveIt.HasNext() {
		actualTokens = append(actualTokens, moveIt.Next().GetToken())
	}

	AssertThat(actualTokens).IsEqualTo(expectedTokens)
}
