// Port of com/google/javascript/jscomp/PeepholeMinimizeConditions.java @ closure-compiler v20261021

/*
 * Copyright 2010 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: this pass reads the shapes it rewrites (if, hook, &&, ||, the
// blocks of an if) through the typed AST view of package ast (ifs.Cond(),
// ifs.Else(), hook.Then(), block.OnlyStatement()...) instead of child
// positions. The Java isXBlock/getXBlock pairs are single comma-ok
// functions that return the typed statement or expression. The edits still
// use the rhino API, in the Java's order.

// PeepholeMinimizeConditions is a peephole optimization that minimizes
// conditional expressions according to De Morgan's laws. Also rewrites
// conditional statements as expressions by replacing them with HOOKs and
// short-circuit binary operators.
//
// Based on PeepholeSubstituteAlternateSyntax
type PeepholeMinimizeConditions struct {
	AbstractPeepholeOptimizationBase
	late      bool
	scratchMC MinimizedCondition
}

var peepholeMinimizeConditions_AND_PRECEDENCE = NodeUtil_Precedence(rhino.Token_AND)

// NewPeepholeMinimizeConditions creates the optimization.
//
// late: When late is false, this mean we are currently running before most
// of the other optimizations. In this case we would avoid optimizations that
// would make the code harder to analyze (such as using string splitting,
// merging statements with commas, etc). When late is true, we would do
// anything to minimize for size.
//
// Java: PeepholeMinimizeConditions(boolean)
func NewPeepholeMinimizeConditions(late bool) *PeepholeMinimizeConditions {
	p := &PeepholeMinimizeConditions{}
	p.InitAbstractPeepholeOptimizationBase(p) // the implicit super()
	p.late = late
	return p
}

// OptimizeSubtree tries to apply our various peephole minimizations on the
// passed in node.
//
// Java: PeepholeMinimizeConditions.optimizeSubtree(Node)
func (p *PeepholeMinimizeConditions) OptimizeSubtree(node rhino.Node) rhino.Node {
	switch node.GetToken() {
	case rhino.Token_THROW, rhino.Token_RETURN:
		result := p.TryRemoveRedundantExit(node)
		if result != node {
			return result
		}
		return p.TryReplaceExitWithBreak(node)
		// TODO(johnlenz): Maybe remove redundant BREAK and CONTINUE. Overlaps
		// with MinimizeExitPoints.
	case rhino.Token_NOT:
		not, _ := ast.AsUnary(node) // the case proves the kind
		p.TryMinimizeCondition(not.Operand().Node())
		return p.TryMinimizeNot(node)
	case rhino.Token_IF:
		ifs, _ := ast.AsIf(node) // the case proves the kind
		p.PerformConditionSubstitutions(ifs.Cond().Node())
		return p.TryMinimizeIf(ifs)
	case rhino.Token_EXPR_RESULT:
		expr, _ := ast.AsExprStmt(node) // the case proves the kind
		p.PerformConditionSubstitutions(expr.Expr().Node())
		return p.TryMinimizeExprResult(expr)
	case rhino.Token_HOOK:
		hook, _ := ast.AsHook(node) // the case proves the kind
		p.PerformConditionSubstitutions(hook.Cond().Node())
		return p.TryMinimizeHook(hook)
	case rhino.Token_WHILE:
		while, _ := ast.AsWhile(node) // the case proves the kind
		p.TryMinimizeCondition(while.Cond().Node())
		return node
	case rhino.Token_DO:
		do, _ := ast.AsDo(node) // the case proves the kind
		p.TryMinimizeCondition(do.Cond().Node())
		return node
	case rhino.Token_FOR:
		forNode, _ := ast.AsFor(node) // the case proves the kind
		p.TryJoinForCondition(forNode)
		// (Go: for (;;) has no condition to minimize; the Java tries its EMPTY
		// placeholder, which has no boolean value.)
		if cond, ok := forNode.Cond(); ok {
			p.TryMinimizeCondition(cond.Node())
		}
		return node
	case rhino.Token_BLOCK:
		block, _ := ast.AsBlock(node) // the case proves the kind
		return p.TryReplaceIf(block)
	default:
		return node
		// Nothing changed
	}
}

// Java: PeepholeMinimizeConditions.tryJoinForCondition(Node)
func (p *PeepholeMinimizeConditions) TryJoinForCondition(n ast.For) {
	if !p.late {
		return
	}

	block := n.Body()
	first, ok := block.FirstStatement()
	if !ok {
		return
	}
	maybeIf, ok := first.AsIf()
	if !ok {
		return
	}
	thenBlock := maybeIf.Then()
	maybeBreak, ok := thenBlock.FirstStatement()
	if !ok {
		return
	}
	brk, ok := maybeBreak.AsBreak()
	if !ok {
		return
	}
	if _, hasLabel := brk.Label(); hasLabel {
		return
	}

	ifCondition := maybeIf.Cond()
	// Preserve the IF ELSE expression is there is one.
	if els, ok := maybeIf.Else(); ok {
		maybeIf.Node().ReplaceWith(els.Node().Detach())
	} else {
		NodeUtil_RedeclareVarsInsideBranch(thenBlock.Node())
		maybeIf.Node().Detach()
	}

	ifCondition.Node().Detach()
	fixedIfCondition := ast.NewUnary(rhino.Token_NOT, ifCondition)
	fixedIfCondition.Node().Srcref(ifCondition.Node())

	// OK, join the IF expression with the FOR expression
	if forCondition, ok := n.Cond(); ok {
		replacement := rhino.NewNode(rhino.Token_AND)
		forCondition.Node().ReplaceWith(replacement)
		replacement.AddChildToBack(forCondition.Node())
		replacement.AddChildToBack(fixedIfCondition.Node())
		p.ReportChangeToEnclosingScope(replacement)
	} else {
		// for (init; ; update): the condition replaces the EMPTY placeholder.
		n.SetCond(fixedIfCondition.AsExpr())
		p.ReportChangeToEnclosingScope(fixedIfCondition.Node())
	}
}

// TryReplaceIf uses "return x?1:2;" in place of "if(x)return 1;return 2;".
//
// Java: PeepholeMinimizeConditions.tryReplaceIf(Node)
func (p *PeepholeMinimizeConditions) TryReplaceIf(n ast.Block) rhino.Node {
	// The next statement is read before the body runs: the body may detach
	// the current one.
	next, hasNext := n.FirstStatement()
	for hasNext {
		child := next
		next, hasNext = child.Next()
		ifs, ok := child.AsIf()
		if !ok {
			continue
		}
		cond := ifs.Cond().Node()
		thenBranch := ifs.Then()
		elseBranch, hasElse := ifs.Else()

		if hasNext && !hasElse && PeepholeMinimizeConditions_IsReturnBlock(thenBranch) {
			if nextIf, ok := next.AsIf(); ok {
				nextCond := nextIf.Cond().Node()
				nextThen := nextIf.Then()
				nextElse, nextHasElse := nextIf.Else()
				if p.AreNodesEqualForInlining(thenBranch.Node(), nextThen.Node()) {
					// Transform
					//   if (x) return 1; if (y) return 1;
					// to
					//   if (x||y) return 1;
					ifs.Node().Detach()
					ifs.Node().DetachChildren()
					newCond := rhino.NewNodeWithChild(rhino.Token_OR, cond)
					nextCond.ReplaceWith(newCond)
					newCond.AddChildToBack(nextCond)
					p.ReportChangeToEnclosingScope(newCond)
				} else if nextHasElse && p.AreNodesEqualForInlining(thenBranch.Node(), nextElse.Node()) {
					// Transform
					//   if (x) return 1; if (y) foo() else return 1;
					// to
					//   if (!x&&y) foo() else return 1;
					ifs.Node().Detach()
					ifs.Node().DetachChildren()
					newCond := rhino.NewNodeWithChild(rhino.Token_AND,
						rhino.IR_Not(cond).Srcref(cond))
					nextCond.ReplaceWith(newCond)
					newCond.AddChildToBack(nextCond)
					p.ReportChangeToEnclosingScope(newCond)
				}
			} else if elseExpr, ok := PeepholeMinimizeConditions_GetReturnExpression(next); ok {
				var thenExpr rhino.Node
				// if(x)return; return 1 -> return x?void 0:1
				if value, ok := PeepholeMinimizeConditions_GetBlockReturnExpression(thenBranch); ok {
					thenExpr = value.Node()
					thenExpr.Detach()
				} else {
					thenExpr = NodeUtil_NewUndefinedNode(ifs.Node())
				}

				cond.Detach()
				elseExpr.Node().Detach()

				returnNode := rhino.IR_ReturnNodeWithExpr(
					rhino.IR_Hook(cond, thenExpr, elseExpr.Node()).
						Srcref(ifs.Node()))
				ifs.Node().ReplaceWith(returnNode)
				next.Node().Detach()
				p.ReportChangeToEnclosingScope(n.Node())
				// everything else in the block is dead code.
				break
			}
		} else if hasElse && PeepholeMinimizeConditions_StatementMustExitParent(thenBranch.AsStatement()) {
			ifs.RemoveElse()
			elseBranch.Node().InsertAfter(ifs.Node())
			p.ReportChangeToEnclosingScope(n.Node())
		}
	}
	return n.Node()
}

// Java: PeepholeMinimizeConditions.statementMustExitParent(Node)
func PeepholeMinimizeConditions_StatementMustExitParent(n ast.Statement) bool {
	switch n.Node().GetToken() {
	case rhino.Token_THROW, rhino.Token_RETURN:
		return true
	case rhino.Token_BLOCK:
		block, _ := n.AsBlock() // the case proves the kind
		if child, ok := block.LastStatement(); ok {
			return PeepholeMinimizeConditions_StatementMustExitParent(child)
		}
		return false
	// TODO(johnlenz): handle TRY/FINALLY
	case rhino.Token_FUNCTION:
		return false
	default:
		return false
	}
}

// TryReplaceExitWithBreak replaces duplicate exits in control structures.
// If the node following the exit node expression has the same effect as exit
// node, the node can be replaced or removed. For example:
//
//	"while (a) {return f()} return f();" ==> "while (a) {break} return f();"
//	"while (a) {throw 'ow'} throw 'ow';" ==> "while (a) {break} throw 'ow';"
//
// n is an follow control exit expression (a THROW or RETURN node). Returns
// the replacement for n, or the original if no change was made.
//
// Java: PeepholeMinimizeConditions.tryReplaceExitWithBreak(Node)
func (p *PeepholeMinimizeConditions) TryReplaceExitWithBreak(n rhino.Node) rhino.Node {
	result := n.GetFirstChild()

	// Find the enclosing control structure, if any, that a "break" would exit
	// from.
	breakTarget := n
	for ; !ControlFlowAnalysis_IsBreakTarget(breakTarget, nil /* no label */); breakTarget = breakTarget.GetParent() {
		if breakTarget.IsFunction() || breakTarget.IsScript() {
			// No break target.
			return n
		}
	}

	follow := ControlFlowAnalysis_ComputeFollowNode(breakTarget)

	// Skip pass all the finally blocks because both the break and return will
	// also trigger all the finally blocks. However, the order of execution is
	// slightly changed. Consider:
	//
	// return a() -> finally { b() } -> return a()
	//
	// which would call a() first. However, changing the first return to a
	// break will result in calling b().

	prefinallyFollows := follow
	follow = PeepholeMinimizeConditions_SkipFinallyNodes(follow)

	innerFollow := ControlFlowAnalysis_ComputeFollowNode(n)
	if prefinallyFollows != follow || innerFollow != PeepholeMinimizeConditions_SkipFinallyNodes(innerFollow) {
		// There were finally clauses
		if !p.IsPure(result) {
			// Can't defer the exit
			return n
		}
	}

	if follow.IsNil() && (n.IsThrow() || !result.IsNil()) {
		// Can't complete remove a throw here or a return with a result.
		return n
	}

	// When follow is null, this mean the follow of a break target is the
	// end of a function. This means a break is same as return.
	if follow.IsNil() || p.AreMatchingExits(n, follow) {
		replacement := rhino.IR_BreakNode()
		n.ReplaceWith(replacement)
		p.ReportChangeToEnclosingScope(replacement)
		return replacement
	}

	return n
}

// TryRemoveRedundantExit removes duplicate exits. If the node following the
// exit node expression has the same effect as exit node, the node can be
// removed. For example: "if (a) {return f()} return f();" ==> "if (a) {}
// return f();" "if (a) {throw 'ow'} throw 'ow';" ==> "if (a) {} throw 'ow';"
//
// n is an follow control exit expression (a THROW or RETURN node). Returns
// the replacement for n, or the original if no change was made.
//
// Java: PeepholeMinimizeConditions.tryRemoveRedundantExit(Node)
func (p *PeepholeMinimizeConditions) TryRemoveRedundantExit(n rhino.Node) rhino.Node {
	exitExpr := n.GetFirstChild()

	follow := ControlFlowAnalysis_ComputeFollowNode(n)

	// Skip pass all the finally blocks because both the fall through and return
	// will also trigger all the finally blocks.
	prefinallyFollows := follow
	follow = PeepholeMinimizeConditions_SkipFinallyNodes(follow)
	if prefinallyFollows != follow {
		// There were finally clauses
		if !p.IsPure(exitExpr) {
			// Can't replace the return
			return n
		}
	}

	if follow.IsNil() && (n.IsThrow() || !exitExpr.IsNil()) {
		// Can't complete remove a throw here or a return with a result.
		return n
	}

	// When follow is null, this mean the follow of a break target is the
	// end of a function. This means a break is same as return.
	if follow.IsNil() || p.AreMatchingExits(n, follow) {
		p.ReportChangeToEnclosingScope(n)
		n.Detach()
		return rhino.NilNode
	}

	return n
}

// IsPure returns whether the expression does not produces and can not be
// affected by side-effects.
//
// Java: PeepholeMinimizeConditions.isPure(Node)
func (p *PeepholeMinimizeConditions) IsPure(n rhino.Node) bool {
	return n.IsNil() ||
		(!NodeUtil_CanBeSideEffected(n) &&
			!p.MayHaveSideEffects(n))
}

// PeepholeMinimizeConditions_SkipFinallyNodes returns n or the node
// following any following finally nodes.
//
// Java: PeepholeMinimizeConditions.skipFinallyNodes(Node)
func PeepholeMinimizeConditions_SkipFinallyNodes(n rhino.Node) rhino.Node {
	for !n.IsNil() && NodeUtil_IsTryFinallyNode(n.GetParent(), n) {
		n = ControlFlowAnalysis_ComputeFollowNode(n)
	}
	return n
}

// AreMatchingExits checks whether one exit can be replaced with another.
// Verify:
//  1. They are identical expressions
//  2. If an exception is possible that the statements, the original and the
//     potential replacement are in the same exception handler.
//
// Java: PeepholeMinimizeConditions.areMatchingExits(Node, Node)
func (p *PeepholeMinimizeConditions) AreMatchingExits(nodeThis rhino.Node, nodeThat rhino.Node) bool {
	if !p.IsASTNormalized() &&
		(nodeThis.IsThrow() || nodeThis.IsReturn()) &&
		nodeThis.HasChildren() {
		// if the ast isn't normalized "return a" or "throw a" may not mean the same thing in
		// different blocks.
		return false
	}
	return nodeThis.IsEquivalentTo(nodeThat) &&
		(!PeepholeMinimizeConditions_IsExceptionPossible(nodeThis) ||
			PeepholeMinimizeConditions_GetExceptionHandler(nodeThis) == PeepholeMinimizeConditions_GetExceptionHandler(nodeThat))
}

// Java: PeepholeMinimizeConditions.isExceptionPossible(Node)
func PeepholeMinimizeConditions_IsExceptionPossible(n rhino.Node) bool {
	// TODO(johnlenz): maybe use ControlFlowAnalysis.mayThrowException?
	if _, ok := ast.AsThrow(n); ok {
		return true
	}
	ret, ok := ast.AsReturn(n)
	jutil.CheckStateMsg1(ok, "%s", n)
	value, ok := ret.Value()
	return ok && !NodeUtil_IsLiteralValue(value.Node(), true)
}

// Java: PeepholeMinimizeConditions.getExceptionHandler(Node)
func PeepholeMinimizeConditions_GetExceptionHandler(n rhino.Node) rhino.Node {
	return ControlFlowAnalysis_GetExceptionHandler(n)
}

// TryMinimizeNot tries to minimize NOT nodes such as !(x==y).
//
// # Returns the replacement for n or the original if no change was made
//
// Java: PeepholeMinimizeConditions.tryMinimizeNot(Node)
func (p *PeepholeMinimizeConditions) TryMinimizeNot(n rhino.Node) rhino.Node {
	not, ok := ast.AsUnary(n)
	jutil.CheckArgument(ok && not.Token() == rhino.Token_NOT)
	parent := n.GetParent()

	notChild := not.Operand().Node()
	// negative operator of the current one : == -> != for instance.
	var complementOperator rhino.Token
	switch notChild.GetToken() {
	case rhino.Token_EQ:
		complementOperator = rhino.Token_NE
	case rhino.Token_NE:
		complementOperator = rhino.Token_EQ
	case rhino.Token_SHEQ:
		complementOperator = rhino.Token_SHNE
	case rhino.Token_SHNE:
		complementOperator = rhino.Token_SHEQ
	case rhino.Token_NOT:
		// !!(a < b) -> a < b (when the inner expression already produces a boolean)
		innerNot, _ := ast.AsUnary(notChild) // the case proves the kind
		if expr := innerNot.Operand().Node(); NodeUtil_IsBooleanResult(expr) {
			expr.Detach()
			n.ReplaceWith(expr)
			p.ReportChangeToEnclosingScope(parent)
			return expr
		}
		return n
	// GT, GE, LT, LE are not handled in this because !(x<NaN) != x>=NaN.
	default:
		return n
	}
	newOperator := notChild
	newOperator.Detach()
	newOperator.SetToken(complementOperator)
	n.ReplaceWith(newOperator)
	p.ReportChangeToEnclosingScope(parent)
	return newOperator
}

func peepholeMinimizeConditions_canMinimizeConditionNode(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_NOT, rhino.Token_AND, rhino.Token_OR, rhino.Token_HOOK, rhino.Token_COMMA:
		return true
	default:
		return false
	}
}

// TryMinimizeExprResult tries to remove leading NOTs from EXPR_RESULTS.
//
// Returns the replacement for n or the original if no replacement was
// necessary.
//
// Java: PeepholeMinimizeConditions.tryMinimizeExprResult(Node)
func (p *PeepholeMinimizeConditions) TryMinimizeExprResult(n ast.ExprStmt) rhino.Node {
	originalExpr := n.Expr().Node()
	if !peepholeMinimizeConditions_canMinimizeConditionNode(originalExpr) {
		return n.Node()
	}
	minCond := MinimizedCondition_FromConditionNodeInto(originalExpr, &p.scratchMC)
	mNode :=
		minCond.GetMinimized(MinimizedCondition_MinimizationStyle_ALLOW_LEADING_NOT)
	if mNode.IsNot() {
		// Remove the leading NOT in the EXPR_RESULT.
		p.ReplaceNode(originalExpr, mNode.WithoutNot())
	} else {
		p.ReplaceNode(originalExpr, mNode)
	}
	return n.Node()
}

// TryMinimizeHook tries flipping HOOKs that have negated conditions.
//
// Returns the replacement for n or the original if no replacement was
// necessary.
//
// Java: PeepholeMinimizeConditions.tryMinimizeHook(Node)
func (p *PeepholeMinimizeConditions) TryMinimizeHook(hook ast.Hook) rhino.Node {
	n := hook.Node()
	originalCond := hook.Cond().Node()
	thenBranch := hook.Then().Node()
	elseBranch := hook.Else().Node()

	// Fold boolean-literal HOOKs:
	//   x ? true : false -> x   (if x is already a boolean result)
	//   x ? true : false -> !!x (otherwise)
	//   x ? false : true -> !x
	// Note: identical branches (e.g. x ? true : true -> true or (x, true)) are handled
	// generally for any equal subtrees (x ? a : a) in PeepholeRemoveDeadCode.tryFoldHook.
	if NodeUtil_IsBooleanResult(thenBranch) && NodeUtil_IsBooleanResult(elseBranch) {
		thenVal := p.GetSideEffectFreeBooleanValue(thenBranch)
		elseVal := p.GetSideEffectFreeBooleanValue(elseBranch)
		var replacement rhino.Node
		if thenVal == base.Tri_TRUE && elseVal == base.Tri_FALSE {
			// x ? true : false -> x (if boolean) or !!x
			cond := originalCond.Detach()
			if NodeUtil_IsBooleanResult(cond) {
				replacement = cond
			} else {
				replacement = rhino.IR_Not(rhino.IR_Not(cond).Srcref(n)).Srcref(n)
			}
		} else if thenVal == base.Tri_FALSE && elseVal == base.Tri_TRUE {
			// x ? false : true -> !x
			cond := originalCond.Detach()
			replacement = rhino.IR_Not(cond).Srcref(n)
		}
		if !replacement.IsNil() {
			// Replace `x ? true : false` or `x ? false : true` with `replacement` (`x`, `!!x`, or `!x`)
			n.ReplaceWith(replacement)
			p.ReportChangeToEnclosingScope(replacement)
			p.MarkFunctionsDeleted(n)
			if replacement.IsNot() {
				return p.TryMinimizeNot(replacement)
			}
			return replacement
		}
	}

	if !peepholeMinimizeConditions_canMinimizeConditionNode(originalCond) {
		return n
	}
	minCond := MinimizedCondition_FromConditionNodeInto(originalCond, &p.scratchMC)
	mNode :=
		minCond.GetMinimized(MinimizedCondition_MinimizationStyle_ALLOW_LEADING_NOT)
	if mNode.IsNot() {
		// Swap the HOOK
		p.ReplaceNode(originalCond, mNode.WithoutNot())
		thenBranch.Detach()
		n.AddChildToBack(thenBranch)
		p.ReportChangeToEnclosingScope(n)
	} else {
		p.ReplaceNode(originalCond, mNode)
	}
	return n
}

// TryMinimizeIf tries turning IF nodes into smaller HOOKs.
//
// Returns the replacement for n or the original if no replacement was
// necessary.
//
// Java: PeepholeMinimizeConditions.tryMinimizeIf(Node)
func (p *PeepholeMinimizeConditions) TryMinimizeIf(ifs ast.If) rhino.Node {
	n := ifs.Node()
	parent := n.GetParent()

	originalCond := ifs.Cond().Node()

	/* If the condition is a literal, we'll let other
	 * optimizations try to remove useless code.
	 */
	if NodeUtil_IsLiteralValue(originalCond, true) {
		return n
	}

	// Read both branches now: the edits below detach the condition, after
	// which ifs.Then() would be the wrong child.
	thenBranch := ifs.Then()
	elseBranch, hasElse := ifs.Else()

	minCond := MinimizedCondition_FromConditionNodeInto(originalCond, &p.scratchMC)

	// Compute two minimized representations. The first representation counts
	// a leading NOT node, and the second ignores a leading NOT node.
	// If we can fold the if statement into a HOOK or boolean operation,
	// then the NOT node does not matter, and we prefer the second condition.
	// If we cannot fold the if statement, then we prefer the first condition.
	unnegatedCond := minCond.GetMinimized(MinimizedCondition_MinimizationStyle_PREFER_UNNEGATED)
	shortCond := minCond.GetMinimized(MinimizedCondition_MinimizationStyle_ALLOW_LEADING_NOT)

	if !hasElse {
		if expr, ok := PeepholeMinimizeConditions_GetBlockExpression(thenBranch); ok {
			if !p.late && PeepholeMinimizeConditions_IsPropertyAssignmentInExpression(expr.Node()) {
				// Keep opportunities for CollapseProperties such as
				// a.longIdentifier || a.longIdentifier = ... -> var a = ...;
				// until CollapseProperties has been run.
				p.ReplaceNode(originalCond, unnegatedCond)
				return n
			}

			if shortCond.IsNot() {
				// if(!x)bar(); -> x||bar();
				replacementCond := p.ReplaceNode(originalCond, shortCond.WithoutNot()).Detach()
				return p.replaceIfWithExpr(ifs, thenBranch, expr, rhino.Token_OR, replacementCond)
			}
			// True, but removed for performance reasons.
			// Preconditions.checkState(shortCond.isEquivalentTo(unnegatedCond));

			// if(x)foo(); -> x&&foo();
			if shortCond.IsLowerPrecedenceThan(peepholeMinimizeConditions_AND_PRECEDENCE) &&
				PeepholeMinimizeConditions_IsLowerPrecedence(expr.Expr().Node(), peepholeMinimizeConditions_AND_PRECEDENCE) {
				// One additional set of parentheses is worth the change even if
				// there is no immediate code size win. However, two extra pair of
				// {}, we would have to think twice. (unless we know for sure the
				// we can further optimize its parent.
				p.ReplaceNode(originalCond, shortCond)
				return n
			}

			replacementCond := p.ReplaceNode(originalCond, shortCond).Detach()
			return p.replaceIfWithExpr(ifs, thenBranch, expr, rhino.Token_AND, replacementCond)
		}

		// Try to combine two IF-ELSE
		if inner, ok := thenBranch.OnlyStatement(); ok {
			if innerIf, ok := inner.AsIf(); ok {
				innerCond := innerIf.Cond().Node()
				innerThenBranch := innerIf.Then().Node()
				_, innerHasElse := innerIf.Else()

				if !innerHasElse &&
					!(unnegatedCond.IsLowerPrecedenceThan(peepholeMinimizeConditions_AND_PRECEDENCE) &&
						PeepholeMinimizeConditions_IsLowerPrecedence(innerCond, peepholeMinimizeConditions_AND_PRECEDENCE)) {
					replacementCond := p.ReplaceNode(originalCond, unnegatedCond).Detach()
					n.DetachChildren()
					n.AddChildToBack(
						rhino.IR_And(
							replacementCond,
							innerCond.Detach()).
							Srcref(originalCond))
					n.AddChildToBack(innerThenBranch.Detach())
					p.ReportChangeToEnclosingScope(n)
					// Not worth trying to fold the current IF-ELSE into && because
					// the inner IF-ELSE wasn't able to be folded into && anyways.
					return n
				}
			}
		}
		p.ReplaceNode(originalCond, unnegatedCond)
		return n
	}

	/* TODO(dcc) This modifies the siblings of n, which is undesirable for a
	 * peephole optimization. This should probably get moved to another pass.
	 */
	p.TryRemoveRepeatedStatements(ifs)

	// if(!x)foo();else bar(); -> if(x)bar();else foo();
	// An additional set of curly braces isn't worth it.
	if shortCond.IsNot() && !PeepholeMinimizeConditions_ConsumesDanglingElse(elseBranch.AsStatement()) {
		p.ReplaceNode(originalCond, shortCond.WithoutNot())
		thenBranch.Node().Detach()
		n.AddChildToBack(thenBranch.Node())
		p.ReportChangeToEnclosingScope(n)
		return n
	}

	// if(x)return 1;else return 2; -> return x?1:2;
	thenReturned, thenIsReturn := PeepholeMinimizeConditions_GetBlockReturnExpression(thenBranch)
	elseReturned, elseIsReturn := PeepholeMinimizeConditions_GetBlockReturnExpression(elseBranch)
	if thenIsReturn && elseIsReturn {
		thenExpr := thenReturned.Node()
		elseExpr := elseReturned.Node()

		replacementCond := p.ReplaceNode(originalCond, shortCond).Detach()
		thenExpr.Detach()
		elseExpr.Detach()

		// note - we ignore any cases with "return;", technically this
		// can be converted to "return undefined;" or some variant, but
		// that does not help code size.
		returnNode := rhino.IR_ReturnNodeWithExpr(
			rhino.IR_Hook(replacementCond, thenExpr, elseExpr).
				Srcref(n))
		n.ReplaceWith(returnNode)
		p.ReportChangeToEnclosingScope(returnNode)
		return returnNode
	}

	thenStmt, thenBranchIsExpressionBlock := PeepholeMinimizeConditions_GetBlockExpression(thenBranch)
	elseStmt, elseBranchIsExpressionBlock := PeepholeMinimizeConditions_GetBlockExpression(elseBranch)

	if thenBranchIsExpressionBlock && elseBranchIsExpressionBlock {
		thenOp := thenStmt.Expr().Node()
		elseOp := elseStmt.Expr().Node()
		// if(x)a=1;else a=2; -> a=x?1:2;
		thenAssign, thenIsAssign := thenStmt.Expr().AsAssignment()
		elseAssign, elseIsAssign := elseStmt.Expr().AsAssignment()
		if thenIsAssign && elseIsAssign && thenAssign.Token() == elseAssign.Token() {
			lhs := thenAssign.Target()
			_, isPlainAssign := thenAssign.AsAssign()
			_, lhsIsName := lhs.AsName()
			if p.AreNodesEqualForInlining(lhs.Node(), elseAssign.Target().Node()) &&
				// if LHS has side effects, don't proceed [since the optimization
				// evaluates LHS before cond]
				// NOTE - there are some circumstances where we can
				// proceed even if there are side effects...
				!p.MayEffectMutableState(lhs.Node()) &&
				(!p.MayHaveSideEffects(originalCond) ||
					(isPlainAssign && lhsIsName)) {

				assignName := lhs.Node()
				thenExpr := thenAssign.Value().Node()
				elseExpr := elseAssign.Value().Node()

				replacementCond := p.ReplaceNode(originalCond, shortCond).Detach()
				assignName.Detach()
				thenExpr.Detach()
				elseExpr.Detach()

				hookNode := rhino.IR_Hook(replacementCond, thenExpr, elseExpr).
					Srcref(n)
				assign := rhino.NewNodeWith2Children(thenAssign.Token(), assignName, hookNode).Srcref(thenOp)
				expr := NodeUtil_NewExpr(assign)
				n.ReplaceWith(expr)
				p.ReportChangeToEnclosingScope(parent)

				return expr
			}
		}
		// if(x)foo();else bar(); -> x?foo():bar()
		replacementCond := p.ReplaceNode(originalCond, shortCond).Detach()
		thenOp.Detach()
		elseOp.Detach()
		expr := rhino.IR_ExprResult(
			rhino.IR_Hook(replacementCond, thenOp, elseOp).Srcref(n))
		n.ReplaceWith(expr)
		p.ReportChangeToEnclosingScope(parent)
		return expr
	}

	thenVar, name1, thenBranchIsVar := PeepholeMinimizeConditions_GetBlockVar(thenBranch)
	elseVar, name2, elseBranchIsVar := PeepholeMinimizeConditions_GetBlockVar(elseBranch)

	// if(x)var y=1;else y=2  ->  var y=x?1:2
	if thenBranchIsVar && elseBranchIsExpressionBlock {
		if elseAssign, ok := elseStmt.Expr().AsAssign(); ok {
			thenValue, hasValue := name1.Init()
			maybeName2, isName := elseAssign.Target().AsName()

			if hasValue && isName && name1.Ident() == maybeName2.Ident() {
				thenExpr := thenValue.Node()
				thenExpr.Detach()
				elseExpr := elseAssign.Value().Node()
				elseExpr.Detach()
				replacementCond := p.ReplaceNode(originalCond, shortCond).Detach()
				hookNode := rhino.IR_Hook(replacementCond, thenExpr, elseExpr).Srcref(n)
				thenVar.Node().Detach()
				name1.Node().AddChildToBack(hookNode)
				n.ReplaceWith(thenVar.Node())
				p.ReportChangeToEnclosingScope(parent)
				return thenVar.Node()
			}
		}

		// if(x)y=1;else var y=2  ->  var y=x?1:2
	} else if elseBranchIsVar && thenBranchIsExpressionBlock {
		if thenAssign, ok := thenStmt.Expr().AsAssign(); ok {
			maybeName1, isName := thenAssign.Target().AsName()
			elseValue, hasValue := name2.Init()

			if hasValue && isName && maybeName1.Ident() == name2.Ident() {
				thenExpr := thenAssign.Value().Node()
				thenExpr.Detach()
				elseExpr := elseValue.Node()
				elseExpr.Detach()
				replacementCond := p.ReplaceNode(originalCond, shortCond).Detach()
				hookNode := rhino.IR_Hook(replacementCond, thenExpr, elseExpr).Srcref(n)
				elseVar.Node().Detach()
				name2.Node().AddChildToBack(hookNode)
				n.ReplaceWith(elseVar.Node())
				p.ReportChangeToEnclosingScope(parent)

				return elseVar.Node()
			}
		}
	}

	p.ReplaceNode(originalCond, unnegatedCond)
	return n
}

// replaceIfWithExpr replaces ifs, an if without else whose then branch is
// the block {expr;}, with the statement cond op expr; (op is OR or AND,
// cond the new condition, already detached), and returns that statement.
//
// When the if and its block carry no cold metadata, they are reused in
// place: the block becomes the OR/AND and the if its EXPR_RESULT.
//
// Go-only: the Java builds NodeUtil.newExpr(IR.or(cond, expr)) (or IR.and)
// in tryMinimizeIf; thenBranch is passed in because ifs.Then() is no longer
// the then branch once the condition is detached.
func (p *PeepholeMinimizeConditions) replaceIfWithExpr(ifs ast.If, thenBranch ast.Block, expr ast.ExprStmt, op rhino.Token, cond rhino.Node) rhino.Node {
	n := ifs.Node()
	parent := n.GetParent()
	rhs := expr.Expr().Node()
	rhs.Detach()
	if block := thenBranch.Node(); block.CanResetInPlace() && n.CanResetInPlace() {
		block.Detach()
		binary := block.ResetDetachedToBinaryOpFromSrcref(op, cond, rhs, n)
		newExpr := n.MutateAttachedToExprResultInPlace(binary)
		p.ReportChangeToEnclosingScope(parent)
		return newExpr
	}
	binary := rhino.IR_BinaryOp(op, cond, rhs).Srcref(n)
	newExpr := NodeUtil_NewExpr(binary)
	n.ReplaceWith(newExpr)
	p.ReportChangeToEnclosingScope(parent)

	return newExpr
}

// TryRemoveRepeatedStatements tries to remove duplicate statements from IF
// blocks. For example:
//
//	if (a) {
//	  x = 1;
//	  return true;
//	} else {
//	  x = 2;
//	  return true;
//	}
//
// becomes:
//
//	if (a) {
//	  x = 1;
//	} else {
//	  x = 2;
//	}
//	return true;
//
// n is the IF node to examine.
//
// Java: PeepholeMinimizeConditions.tryRemoveRepeatedStatements(Node)
func (p *PeepholeMinimizeConditions) TryRemoveRepeatedStatements(n ast.If) {
	// Only run this if variable names are guaranteed to be unique. Otherwise bad things can happen:
	// see PeepholeMinimizeConditionsTest#testDontRemoveDuplicateStatementsWithoutNormalization
	if !p.IsASTNormalized() {
		return
	}

	parent := n.Node().GetParent()
	if !NodeUtil_IsStatementBlock(parent) {
		// If the immediate parent is something like a label, we
		// can't move the statement, so bail.
		return
	}

	trueBranch := n.Then()
	falseBranch, ok := n.Else()
	jutil.CheckState(ok)

	for {
		lastTrue, ok1 := trueBranch.LastStatement()
		lastFalse, ok2 := falseBranch.LastStatement()
		if !ok1 || !ok2 ||
			!p.AreNodesEqualForInlining(lastTrue.Node(), lastFalse.Node()) {
			break
		}
		lastTrue.Node().Detach()
		lastFalse.Node().Detach()
		lastTrue.Node().InsertAfter(n.Node())
		p.ReportChangeToEnclosingScope(parent)
	}
}

// PeepholeMinimizeConditions_GetBlockExpression returns the expression
// statement of n if n is a block with a single statement that is an
// expression, and that expression can be folded into another; ok is false
// otherwise.
//
// Java: PeepholeMinimizeConditions.isFoldableExpressBlock(Node) and
// getBlockExpression(Node) (Go: one comma-ok function).
func PeepholeMinimizeConditions_GetBlockExpression(n ast.Block) (ast.ExprStmt, bool) {
	s, ok := n.OnlyStatement()
	if !ok {
		return ast.ExprStmt{}, false
	}
	maybeExpr, ok := s.AsExprStmt()
	if !ok {
		return ast.ExprStmt{}, false
	}
	// IE has a bug where event handlers behave differently when
	// their return value is used vs. when their return value is in
	// an EXPR_RESULT. It's pretty freaking weird. See:
	// http://blickly.github.io/closure-compiler-issues/#291
	// We try to detect this case, and not fold EXPR_RESULTs
	// into other expressions.
	// e.g.:
	// if (e.onchange) {
	//    e.onchange({
	//        _extendedByPrototype: Prototype.emptyFunction,
	//        target: e
	//    });
	//  }
	if call, ok := maybeExpr.Expr().AsNormalOrOptChainCall(); ok {
		// We only have to worry about methods with an implicit 'this'
		// param, or this doesn't happen: calls of obj[key] or obj.onX.
		if calledFn, ok := call.Callee().AsNormalOrOptChainGet(); ok {
			getProp, isGetProp := calledFn.AsNormalOrOptChainGetProp()
			if !isGetProp || jutil.StringStartsWith(getProp.Property(), "on") {
				return ast.ExprStmt{}, false
			}
		}
	}
	return maybeExpr, true
}

// PeepholeMinimizeConditions_IsReturnBlock returns whether the node is a
// block with a single statement that is an return with or without an
// expression.
//
// Java: PeepholeMinimizeConditions.isReturnBlock(Node)
func PeepholeMinimizeConditions_IsReturnBlock(n ast.Block) bool {
	first, ok := n.OnlyStatement()
	if !ok {
		return false
	}
	_, ok = first.AsReturn()
	return ok
}

// PeepholeMinimizeConditions_GetBlockReturnExpression returns the
// expression of n if n is a block with a single statement that is a return
// with an expression; ok is false otherwise.
//
// Java: PeepholeMinimizeConditions.isReturnExpressBlock(Node) and
// getBlockReturnExpression(Node) (Go: one comma-ok function).
func PeepholeMinimizeConditions_GetBlockReturnExpression(n ast.Block) (ast.Expr, bool) {
	first, ok := n.OnlyStatement()
	if !ok {
		return ast.Expr{}, false
	}
	return PeepholeMinimizeConditions_GetReturnExpression(first)
}

// PeepholeMinimizeConditions_GetReturnExpression returns the expression of
// n if n is a return statement with an expression; ok is false otherwise.
//
// Java: PeepholeMinimizeConditions.isReturnExpression(Node) (Go: comma-ok,
// with the expression the callers read next).
func PeepholeMinimizeConditions_GetReturnExpression(n ast.Statement) (ast.Expr, bool) {
	ret, ok := n.AsReturn()
	if !ok {
		return ast.Expr{}, false
	}
	return ret.Value()
}

// PeepholeMinimizeConditions_GetBlockVar returns the var statement of n and
// the name it declares if n is a block with a single statement that is a
// VAR declaration of a single variable; ok is false otherwise.
//
// Java: PeepholeMinimizeConditions.isVarBlock(Node) and getBlockVar(Node)
// (Go: one comma-ok function. A var that declares a destructuring pattern
// is not one: the callers compare the declared name, where the Java would
// throw.)
func PeepholeMinimizeConditions_GetBlockVar(n ast.Block) (ast.NameDeclaration, ast.Name, bool) {
	first, ok := n.OnlyStatement()
	if !ok {
		return ast.NameDeclaration{}, ast.Name{}, false
	}
	decl, ok := first.AsNameDeclaration()
	if !ok || decl.Token() != rhino.Token_VAR {
		return ast.NameDeclaration{}, ast.Name{}, false
	}
	declarator, ok := decl.OnlyDeclarator()
	if !ok {
		return ast.NameDeclaration{}, ast.Name{}, false
	}
	name, ok := declarator.AsName()
	return decl, name, ok
}

// PeepholeMinimizeConditions_ConsumesDanglingElse reports whether a statement
// consumes a 'dangling else'. A statement consumes a 'dangling else' if an
// 'else' token following the statement would be considered by the parser to
// be part of the statement.
//
// Java: PeepholeMinimizeConditions.consumesDanglingElse(Node)
func PeepholeMinimizeConditions_ConsumesDanglingElse(n ast.Statement) bool {
	for {
		switch n.Node().GetToken() {
		case rhino.Token_IF:
			ifs, _ := n.AsIf() // the case proves the kind
			els, ok := ifs.Else()
			if !ok {
				return true
			}
			// This IF node has an else clause: look into it.
			n = els.AsStatement()
		case rhino.Token_BLOCK:
			block, _ := n.AsBlock() // the case proves the kind
			only, ok := block.OnlyStatement()
			if !ok {
				return false
			}
			// This BLOCK has no curly braces.
			n = only
		case rhino.Token_WITH:
			with, _ := n.AsWith() // the case proves the kind
			n = with.Body().AsStatement()
		case rhino.Token_WHILE, rhino.Token_FOR, rhino.Token_FOR_IN:
			loop, _ := n.AsLoop() // the case proves the kind
			n = loop.Body().AsStatement()
		default:
			return false
		}
	}
}

// PeepholeMinimizeConditions_IsLowerPrecedence reports whether the node type
// has lower precedence than "precedence".
//
// Java: PeepholeMinimizeConditions.isLowerPrecedence(Node, int)
func PeepholeMinimizeConditions_IsLowerPrecedence(n rhino.Node, precedence int) bool {
	return NodeUtil_Precedence(n.GetToken()) < precedence
}

func peepholeMinimizeConditions_isPropertyAssignmentInExpressionPredicate(input rhino.Node) bool {
	return input.IsGetProp() && input.GetParent().IsAssign()
}

// PeepholeMinimizeConditions_IsPropertyAssignmentInExpression reports
// whether the expression contains a property assignment.
//
// Java: PeepholeMinimizeConditions.isPropertyAssignmentInExpression(Node)
func PeepholeMinimizeConditions_IsPropertyAssignmentInExpression(n rhino.Node) bool {
	return NodeUtil_Has(n, peepholeMinimizeConditions_isPropertyAssignmentInExpressionPredicate, NodeUtil_MATCH_NOT_FUNCTION)
}

// TryMinimizeCondition tries to minimize condition expression, as there are
// additional assumptions that can be made when it is known that the final
// result is a boolean.
//
// Returns the replacement for n, or the original if no change was made.
//
// Java: PeepholeMinimizeConditions.tryMinimizeCondition(Node)
func (p *PeepholeMinimizeConditions) TryMinimizeCondition(n rhino.Node) rhino.Node {
	n = p.PerformConditionSubstitutions(n)
	if !peepholeMinimizeConditions_canMinimizeConditionNode(n) {
		return n
	}
	minCond := MinimizedCondition_FromConditionNodeInto(n, &p.scratchMC)
	return p.ReplaceNode(
		n,
		minCond.GetMinimized(MinimizedCondition_MinimizationStyle_PREFER_UNNEGATED))
}

// Java: PeepholeMinimizeConditions.replaceNode(Node, MeasuredNode)
func (p *PeepholeMinimizeConditions) ReplaceNode(original rhino.Node, measuredNodeReplacement *MinimizedCondition_MeasuredNode) rhino.Node {
	if measuredNodeReplacement.WillChange(original) {
		replacement := measuredNodeReplacement.ApplyTo(original)
		p.ReportChangeToEnclosingScope(replacement)
		return replacement
	}
	return original
}

// PerformConditionSubstitutions tries to minimize the given condition by
// applying local substitutions.
//
// The following types of transformations are performed:
//
//	x || true        --> true
//	x && true        --> x
//	x ? false : true --> !x
//	x ? true : y     --> x || y
//	x ? x : y        --> x || y
//
// # Returns the replacement for n, or the original if no change was made
//
// Java: PeepholeMinimizeConditions.performConditionSubstitutions(Node)
func (p *PeepholeMinimizeConditions) PerformConditionSubstitutions(n rhino.Node) rhino.Node {
	parent := n.GetParent()

	switch n.GetToken() {
	case rhino.Token_OR, rhino.Token_AND:
		binary, _ := ast.AsBinary(n) // the case proves the kind
		left := binary.Left().Node()
		right := binary.Right().Node()

		// Because the expression is in a boolean context minimize
		// the children, this can't be done in the general case.
		left = p.PerformConditionSubstitutions(left)
		right = p.PerformConditionSubstitutions(right)

		// Remove useless conditionals
		// Handle the following cases:
		//   x || false --> x
		//   x && true --> x
		// This works regardless of whether x has side effects.
		//
		// If x does not have side effects:
		//   x || true  --> true
		//   x && false  --> false
		//
		// If x may have side effects:
		//   x || true  --> x,true
		//   x && false  --> x,false
		//
		// In the last two cases, code size may increase slightly (adding
		// some parens because the comma operator has a low precedence) but
		// the new AST is easier for other passes to handle.
		rightVal := p.GetSideEffectFreeBooleanValue(right)
		if rightVal != base.Tri_UNKNOWN {
			type_ := binary.Token()
			var replacement rhino.Node
			rval := rightVal.ToBoolean(true)

			// (x || FALSE) => x
			// (x && TRUE) => x
			if (type_ == rhino.Token_OR && !rval) || (type_ == rhino.Token_AND && rval) {
				replacement = left
			} else if !p.MayHaveSideEffects(left) {
				replacement = right
			} else {
				// expr_with_sideeffects || true  =>  expr_with_sideeffects, true
				// expr_with_sideeffects && false  =>  expr_with_sideeffects, false
				n.DetachChildren()
				replacement = rhino.IR_Comma(left, right)
			}

			if !replacement.IsNil() {
				n.DetachChildren()
				n.ReplaceWith(replacement)
				p.ReportChangeToEnclosingScope(parent)
				return replacement
			}
		}
		return n
	case rhino.Token_HOOK:
		hook, _ := ast.AsHook(n) // the case proves the kind
		condition := hook.Cond().Node()
		trueNode := hook.Then().Node()
		falseNode := hook.Else().Node()

		// Because the expression is in a boolean context minimize
		// the result children, this can't be done in the general case.
		// The condition is handled in the general case in #optimizeSubtree
		trueNode = p.PerformConditionSubstitutions(trueNode)
		falseNode = p.PerformConditionSubstitutions(falseNode)

		// Handle five cases:
		//   x ? true : false --> x
		//   x ? false : true --> !x
		//   x ? true : y     --> x || y
		//   x ? y : false    --> x && y

		//   Only when x is NAME, hence x does not have side effects
		//   x ? x : y        --> x || y
		var replacement rhino.Node
		trueNodeVal := p.GetSideEffectFreeBooleanValue(trueNode)
		falseNodeVal := p.GetSideEffectFreeBooleanValue(falseNode)
		if trueNodeVal == base.Tri_TRUE && falseNodeVal == base.Tri_FALSE {
			// Remove useless conditionals, keep the condition
			condition.Detach()
			replacement = condition
		} else if trueNodeVal == base.Tri_FALSE && falseNodeVal == base.Tri_TRUE {
			// Remove useless conditionals, keep the condition
			condition.Detach()
			replacement = rhino.IR_Not(condition)
		} else if trueNodeVal == base.Tri_TRUE {
			// Remove useless true case.
			n.DetachChildren()
			replacement = rhino.IR_Or(condition, falseNode)
		} else if falseNodeVal == base.Tri_FALSE {
			// Remove useless false case
			n.DetachChildren()
			replacement = rhino.IR_And(condition, trueNode)
		} else if !p.MayHaveSideEffects(condition) &&
			!p.MayHaveSideEffects(trueNode) &&
			condition.IsEquivalentTo(trueNode) {
			// Remove redundant condition
			n.DetachChildren()
			replacement = rhino.IR_Or(trueNode, falseNode)
		}

		if !replacement.IsNil() {
			n.ReplaceWith(replacement)
			p.ReportChangeToEnclosingScope(replacement)
			n = replacement
		}

		return n
	default:
		// while(true) --> while(1)
		nVal := p.GetSideEffectFreeBooleanValue(n)
		if nVal != base.Tri_UNKNOWN {
			result := nVal.ToBoolean(true)
			equivalentResult := 0
			if result {
				equivalentResult = 1
			}
			return p.MaybeReplaceChildWithNumber(n, equivalentResult)
		}
		// We can't do anything else currently.
		return n
	}
}

// MaybeReplaceChildWithNumber replaces a node with a number node if the new
// number node is not equivalent to the current node.
//
// Returns the replacement for n if it was replaced, otherwise returns n.
//
// Java: PeepholeMinimizeConditions.maybeReplaceChildWithNumber(Node, int)
func (p *PeepholeMinimizeConditions) MaybeReplaceChildWithNumber(n rhino.Node, num int) rhino.Node {
	newNode := rhino.IR_Number(float64(num))
	if !newNode.IsEquivalentTo(n) {
		n.ReplaceWith(newNode)
		p.ReportChangeToEnclosingScope(newNode)
		p.MarkFunctionsDeleted(n)

		return newNode
	}

	return n
}
