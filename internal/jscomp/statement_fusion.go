// Port of com/google/javascript/jscomp/StatementFusion.java @ closure-compiler v20261015

/*
 * Copyright 2011 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// StatementFusion tries to fuse all the statements in a block into a one
// statement by using COMMAs or statements.
//
// Because COMMAs has the lowest precedence, we never need to insert extra ()
// around. Once we have only one statement in a block, we can then eliminate a
// pair of {}'s. Further more, we can also fold a single statement IF into &&
// or create further opportunities for all the other goodies in
// PeepholeMinimizeConditions.
//
// NOTE(user): The current compiler assumes that there are more ;'s than ,'s
// in a real program, and so it makes sense to prefer fusing statements with
// semicolons rather than commas. This assumption has never been validated on
// a real program.
type StatementFusion struct {
	AbstractPeepholeOptimizationBase
}

// Java: StatementFusion()
func NewStatementFusion() *StatementFusion {
	s := &StatementFusion{}
	s.InitAbstractPeepholeOptimizationBase(s) // the implicit super()
	return s
}

// Java: StatementFusion.optimizeSubtree(Node)
func (s *StatementFusion) OptimizeSubtree(n *rhino.Node) *rhino.Node {
	if n.GetParent().IsFunction() || !s.CanFuseIntoOneStatement(n) {
		return n
	}

	start := n.GetFirstChild()
	end := n.GetLastChild()
	result := StatementFusion_FuseIntoOneStatement(start, end)
	StatementFusion_FuseExpressionIntoControlFlowStatement(result, n.GetLastChild())

	s.ReportChangeToEnclosingScope(n)
	return n
}

// Java: StatementFusion.canFuseIntoOneStatement(Node)
func (s *StatementFusion) CanFuseIntoOneStatement(block *rhino.Node) bool {
	if !block.IsBlock() {
		return false
	}

	// Nothing to do here.
	if !block.HasChildren() || block.HasOneChild() {
		return false
	}

	last := block.GetLastChild()

	for c := block.GetFirstChild(); c != nil; c = c.GetNext() {
		if !c.IsExprResult() && c != last {
			return false
		}
	}

	return s.IsFusableControlStatement(last)
}

// Java: StatementFusion.isFusableControlStatement(Node)
func (s *StatementFusion) IsFusableControlStatement(n *rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_IF, rhino.Token_THROW, rhino.Token_SWITCH, rhino.Token_EXPR_RESULT:
		return true
	case rhino.Token_RETURN:
		// We don't want to add a new return value.
		return n.HasChildren()
	case rhino.Token_FOR:
		// Avoid cases where we have for(var x;_;_) { ....
		return !NodeUtil_IsNameDeclaration(n.GetFirstChild())
	case rhino.Token_FOR_IN:
		// Avoid cases where we have for(var x = foo() in a) { ....
		return !s.MayHaveSideEffects(n.GetFirstChild())
	case rhino.Token_LABEL:
		return s.IsFusableControlStatement(n.GetLastChild())
	case rhino.Token_BLOCK:
		return (s.IsASTNormalized() || NodeUtil_CanMergeBlock(n)) &&
			!n.IsSyntheticBlock() &&
			s.IsFusableControlStatement(n.GetFirstChild())
	default:
	}
	return false
}

// StatementFusion_FuseIntoOneStatement fuses, in a block, a list of
// statements with comma's.
//
// first is the first statement to fuse (inclusive); last is the last
// statement to fuse (exclusive). Returns a single statement that contains all
// the fused statement as one.
//
// Java: StatementFusion.fuseIntoOneStatement(Node, Node)
func StatementFusion_FuseIntoOneStatement(first *rhino.Node, last *rhino.Node) *rhino.Node {
	// Nothing to fuse if there is only one statement.
	if first.GetNext() == last {
		return first
	}

	// Step one: Create a comma tree that contains all the statements.
	commaTree := first.RemoveFirstChild()

	var next *rhino.Node
	for cur := first.GetNext(); cur != last; cur = next {
		commaTree = AstManipulations_FuseExpressions(commaTree, cur.RemoveFirstChild())
		next = cur.GetNext()
		cur.Detach()
	}

	// Step two: The last EXPR_RESULT will now hold the comma tree with all
	// the fused statements.
	first.AddChildToBack(commaTree)
	return first
}

// Java: StatementFusion.fuseExpressionIntoControlFlowStatement(Node, Node)
func StatementFusion_FuseExpressionIntoControlFlowStatement(
	before *rhino.Node, control *rhino.Node) {
	jutil.CheckArgumentMsg(before.IsExprResult(), "before must be expression result")

	// Now we are just left with two statements. The comma tree of the first
	// n - 1 statements (which can be used in an expression) and the last
	// statement. We perform specific fusion based on the last statement's type.
	switch control.GetToken() {
	case rhino.Token_IF, rhino.Token_RETURN, rhino.Token_THROW, rhino.Token_SWITCH, rhino.Token_EXPR_RESULT, rhino.Token_FOR:
		before.Detach()
		StatementFusion_FuseExpressionIntoFirstChild(before.RemoveFirstChild(), control)
	case rhino.Token_FOR_IN:
		before.Detach()
		StatementFusion_FuseExpressionIntoSecondChild(before.RemoveFirstChild(), control)
	case rhino.Token_LABEL:
		StatementFusion_FuseExpressionIntoControlFlowStatement(before, control.GetLastChild())
	case rhino.Token_BLOCK:
		StatementFusion_FuseExpressionIntoControlFlowStatement(before, control.GetFirstChild())
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Statement fusion missing."))
	}
}

// Java: StatementFusion.fuseExpressionIntoFirstChild(Node, Node)
func StatementFusion_FuseExpressionIntoFirstChild(exp *rhino.Node, stmt *rhino.Node) {
	val := stmt.RemoveFirstChild()
	comma := AstManipulations_FuseExpressions(exp, val)
	stmt.AddChildToFront(comma)
}

// Java: StatementFusion.fuseExpressionIntoSecondChild(Node, Node)
func StatementFusion_FuseExpressionIntoSecondChild(exp *rhino.Node, stmt *rhino.Node) {
	val := stmt.GetSecondChild().Detach()
	comma := AstManipulations_FuseExpressions(exp, val)
	comma.InsertAfter(stmt.GetFirstChild())
}
