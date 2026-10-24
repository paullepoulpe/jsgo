// Port of com/google/javascript/jscomp/MinimizedConditionTest.java @ closure-compiler v20261021
// (test/com/google/javascript/jscomp/MinimizedConditionTest.java)

/*
 * Copyright 2013 The Closure Compiler Authors.
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

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// MinimizedConditionTest holds tests for MinimizedCondition in isolation.
// Tests for the containing PeepholeMinimizeConditions pass are in
// PeepholeMinimizeConditionsTest.
type MinimizedConditionTest struct{ t *testing.T }

func (c *MinimizedConditionTest) JUnitSetT(t *testing.T) { c.t = t }
func (c *MinimizedConditionTest) JUnitBefore()           {}
func (c *MinimizedConditionTest) JUnitAfter()            {}

func TestMinimizedConditionTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return &MinimizedConditionTest{} })
}

// Java: MinimizedConditionTest.parseExpr(String)
func minimizedConditionTest_parseExpr(code string) *rhino.Node {
	compiler := NewCompiler()
	input := []*SourceFile{SourceFile_FromCode("code", code)}
	externs := []*SourceFile{}
	options := NewCompilerOptions()
	compiler.Init(externs, input, options)
	root := compiler.ParseInputs()
	AssertWithMessage("Unexpected parse error(s): %s", compilerTestCase_LINE_JOINER(compiler.GetErrors())).
		That(root).
		IsNotNull()
	externsRoot := root.GetFirstChild()
	mainRoot := externsRoot.GetNext()
	script := mainRoot.GetFirstChild()
	exprResult := script.GetFirstChild()
	return exprResult.GetFirstChild()
}

// Java: MinimizedConditionTest.cloneAttachedTree(Node)
func minimizedConditionTest_cloneAttachedTree(n *rhino.Node) *rhino.Node {
	jutil.CheckState(n.GetParent().GetFirstChild() == n)
	return n.GetParent().CloneTree().GetFirstChild()
}

// minimizedConditionTest_minCond tests minimization of input condition.
//
// input is input code containing a condition; positive the representation
// expected to be produced when penalizing a leading NOT; negative the
// representation expected to be produced when not penalizing a leading NOT.
//
// Java: MinimizedConditionTest.minCond(String, String, String)
func minimizedConditionTest_minCond(input string, positive string, negative string) {
	inputNode := minimizedConditionTest_parseExpr(input)
	result1 := MinimizedCondition_FromConditionNode(minimizedConditionTest_cloneAttachedTree(inputNode))
	result2 := MinimizedCondition_FromConditionNode(minimizedConditionTest_cloneAttachedTree(inputNode))
	positiveNode := minimizedConditionTest_parseExpr(positive)
	negativeNode := minimizedConditionTest_parseExpr(negative)
	// With counting the leading NOT node:
	positiveResult := result1.GetMinimized(MinimizedCondition_MinimizationStyle_PREFER_UNNEGATED).BuildReplacement()
	// Without counting the leading NOT node:
	negativeResult := result2.GetMinimized(MinimizedCondition_MinimizationStyle_ALLOW_LEADING_NOT).BuildReplacement()
	if !positiveResult.IsEquivalentTo(positiveNode) {
		AssertWithMessage(
			"Not equal:"+"\nExpected: %s\nBut was : %s\nExpected tree:\n%s\nActual tree:\n%s",
			positive,
			NewCompiler().ToSourceOfNode(positiveResult),
			positiveNode.ToStringTree(),
			positiveResult.ToStringTree()).
			Fail()
	}
	if !negativeResult.IsEquivalentTo(negativeNode) {
		AssertWithMessage(
			"Not equal:"+"\nExpected: %s\nBut was : %s\nExpected tree:\n%s\nActual tree:\n%s",
			negative,
			NewCompiler().ToSourceOfNode(negativeResult),
			negativeNode.ToStringTree(),
			negativeResult.ToStringTree()).
			Fail()
	}
}

// Java: MinimizedConditionTest.testTryMinimizeCondSimple()
func (c *MinimizedConditionTest) TestTryMinimizeCondSimple() {
	minimizedConditionTest_minCond("x", "x", "x")
	minimizedConditionTest_minCond("!x", "!x", "!x")
	minimizedConditionTest_minCond("!!x", "x", "x")
	minimizedConditionTest_minCond("!(x && y)", "!x || !y", "!(x && y)")
}

// Java: MinimizedConditionTest.testMinimizeDemorganSimple()
func (c *MinimizedConditionTest) TestMinimizeDemorganSimple() {
	minimizedConditionTest_minCond("!(x&&y)", "!x||!y", "!(x&&y)")
	minimizedConditionTest_minCond("!(x||y)", "!x&&!y", "!(x||y)")
	minimizedConditionTest_minCond("!x||!y", "!x||!y", "!(x&&y)")
	minimizedConditionTest_minCond("!x&&!y", "!x&&!y", "!(x||y)")
	minimizedConditionTest_minCond("!(x && y && z)", "!(x && y && z)", "!(x && y && z)")
	minimizedConditionTest_minCond("(!a||!b)&&c", "(!a||!b)&&c", "!(a&&b||!c)")
	minimizedConditionTest_minCond("(!a||!b)&&(c||d)", "!(a&&b||!c&&!d)", "!(a&&b||!c&&!d)")
}

// Java: MinimizedConditionTest.testMinimizeBug8494751()
func (c *MinimizedConditionTest) TestMinimizeBug8494751() {
	minimizedConditionTest_minCond(
		"x && (y===2 || !f()) && (y===3 || !h())",
		// TODO(tbreisacher): The 'positive' option could be better:
		// "x && !((y!==2 && f()) || (y!==3 && h()))",
		"!(!x || (y!==2 && f()) || (y!==3 && h()))",
		"!(!x || (y!==2 && f()) || (y!==3 && h()))")

	minimizedConditionTest_minCond(
		"x && (y===2 || !f?.()) && (y===3 || !h?.())",
		"!(!x || (y!==2 && f?.()) || (y!==3 && h?.()))",
		"!(!x || (y!==2 && f?.()) || (y!==3 && h?.()))")
}

// Java: MinimizedConditionTest.testMinimizeComplementableOperator()
func (c *MinimizedConditionTest) TestMinimizeComplementableOperator() {
	minimizedConditionTest_minCond("0===c && (2===a || 1===a)", "0===c && (2===a || 1===a)", "!(0!==c || 2!==a && 1!==a)")
}

// Java: MinimizedConditionTest.testMinimizeHook()
func (c *MinimizedConditionTest) TestMinimizeHook() {
	minimizedConditionTest_minCond("!(x ? y : z)", "(x ? !y : !z)", "!(x ? y : z)")
}

// Java: MinimizedConditionTest.testMinimizeComma()
func (c *MinimizedConditionTest) TestMinimizeComma() {
	minimizedConditionTest_minCond("!(inc(), test())", "inc(), !test()", "!(inc(), test())")
	minimizedConditionTest_minCond("!(inc?.(), test?.())", "inc?.(), !test?.()", "!(inc?.(), test?.())")
	minimizedConditionTest_minCond("!((x,y)&&z)", "(x,!y)||!z", "!((x,y)&&z)")
}
