// Port of com/google/javascript/jscomp/StatementFusion.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: this pass reads block, control-flow, and expression-statement shapes
// through package ast while keeping its rhino.Node peephole interface intact.

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
	return &StatementFusion{}
}

// Java: StatementFusion.optimizeSubtree(Node)
func (s *StatementFusion) OptimizeSubtree(n rhino.Node) rhino.Node {
	if _, inFn := ast.AsFunction(n.GetParent()); inFn || !s.CanFuseIntoOneStatement(n) {
		return n
	}

	block := ast.MustBlock(n)
	start := block.MustFirstStatement().Node()
	end, _ := block.LastStatement()
	result := StatementFusion_FuseIntoOneStatement(start, end.Node())
	last, _ := block.LastStatement()
	StatementFusion_FuseExpressionIntoControlFlowStatement(result, last.Node())

	s.ReportChangeToEnclosingScope(n)
	return n
}

// Java: StatementFusion.canFuseIntoOneStatement(Node)
func (s *StatementFusion) CanFuseIntoOneStatement(blockNode rhino.Node) bool {
	block, ok := ast.AsBlock(blockNode)
	if !ok {
		return false
	}

	// Nothing to do here.
	last, hasStmts := block.LastStatement()
	if _, onlyOne := block.OnlyStatement(); !hasStmts || onlyOne {
		return false
	}

	for stmt := range block.Statements() {
		if _, isExpr := stmt.AsExprStmt(); !isExpr && !stmt.Is(last.Node()) {
			return false
		}
	}

	return s.IsFusableControlStatement(last.Node())
}

// Java: StatementFusion.isFusableControlStatement(Node)
func (s *StatementFusion) IsFusableControlStatement(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_IF, rhino.Token_THROW, rhino.Token_SWITCH, rhino.Token_EXPR_RESULT:
		return true
	case rhino.Token_RETURN:
		// We don't want to add a new return value.
		_, hasValue := ast.MustReturn(n).Value()
		return hasValue
	case rhino.Token_FOR:
		// Avoid cases where we have for(var x;_;_) { ....
		init, hasInit := ast.MustFor(n).Init()
		if !hasInit {
			return true
		}
		_, isDecl := init.AsNameDeclaration()
		return !isDecl
	case rhino.Token_FOR_IN:
		// Avoid cases where we have for(var x = foo() in a) { ....
		return !s.MayHaveSideEffects(ast.MustForIn(n).Target().Node())
	case rhino.Token_LABEL:
		return s.IsFusableControlStatement(ast.MustLabel(n).Body().Node())
	case rhino.Token_BLOCK:
		blk := ast.MustBlock(n)
		first, hasFirst := blk.FirstStatement()
		return (s.IsASTNormalized() || NodeUtil_CanMergeBlock(n)) &&
			!blk.IsSynthetic() &&
			hasFirst &&
			s.IsFusableControlStatement(first.Node())
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
func StatementFusion_FuseIntoOneStatement(first rhino.Node, last rhino.Node) rhino.Node {
	// Nothing to fuse if there is only one statement.
	if first.GetNext() == last {
		return first
	}

	// Step one: Create a comma tree that contains all the statements.
	commaTree := ast.MustExprStmt(first).Expr().Detach().Node()

	var next rhino.Node
	for cur := first.GetNext(); cur != last; cur = next {
		curExpr := ast.MustExprStmt(cur).Expr().Detach().Node()
		commaTree = AstManipulations_FuseExpressions(commaTree, curExpr)
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
	before rhino.Node, control rhino.Node) {
	beforeStmt, ok := ast.AsExprStmt(before)
	jutil.CheckArgumentMsg(ok, "before must be expression result")

	// Now we are just left with two statements. The comma tree of the first
	// n - 1 statements (which can be used in an expression) and the last
	// statement. We perform specific fusion based on the last statement's type.
	switch control.GetToken() {
	case rhino.Token_IF, rhino.Token_RETURN, rhino.Token_THROW, rhino.Token_SWITCH, rhino.Token_EXPR_RESULT, rhino.Token_FOR:
		beforeStmt.Detach()
		StatementFusion_FuseExpressionIntoFirstChild(beforeStmt.Expr().Detach().Node(), control)
	case rhino.Token_FOR_IN:
		beforeStmt.Detach()
		StatementFusion_FuseExpressionIntoSecondChild(beforeStmt.Expr().Detach().Node(), control)
	case rhino.Token_LABEL:
		StatementFusion_FuseExpressionIntoControlFlowStatement(before, ast.MustLabel(control).Body().Node())
	case rhino.Token_BLOCK:
		StatementFusion_FuseExpressionIntoControlFlowStatement(before, ast.MustBlock(control).MustFirstStatement().Node())
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Statement fusion missing."))
	}
}

// Java: StatementFusion.fuseExpressionIntoFirstChild(Node, Node)
func StatementFusion_FuseExpressionIntoFirstChild(exp rhino.Node, stmt rhino.Node) {
	val := stmt.RemoveFirstChild()
	comma := AstManipulations_FuseExpressions(exp, val)
	stmt.AddChildToFront(comma)
}

// Java: StatementFusion.fuseExpressionIntoSecondChild(Node, Node)
func StatementFusion_FuseExpressionIntoSecondChild(exp rhino.Node, stmt rhino.Node) {
	forIn := ast.MustForIn(stmt)
	target := forIn.Target().Node()
	val := forIn.Object().Detach().Node()
	comma := AstManipulations_FuseExpressions(exp, val)
	comma.InsertAfter(target)
}
