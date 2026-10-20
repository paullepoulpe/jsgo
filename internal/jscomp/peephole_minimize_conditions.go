// Port of com/google/javascript/jscomp/PeepholeMinimizeConditions.java @ closure-compiler v20261015

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
	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

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
func (p *PeepholeMinimizeConditions) OptimizeSubtree(node *rhino.Node) *rhino.Node {
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
		p.TryMinimizeCondition(node.GetFirstChild())
		return p.TryMinimizeNot(node)
	case rhino.Token_IF:
		p.PerformConditionSubstitutions(node.GetFirstChild())
		return p.TryMinimizeIf(node)
	case rhino.Token_EXPR_RESULT:
		p.PerformConditionSubstitutions(node.GetFirstChild())
		return p.TryMinimizeExprResult(node)
	case rhino.Token_HOOK:
		p.PerformConditionSubstitutions(node.GetFirstChild())
		return p.TryMinimizeHook(node)
	case rhino.Token_WHILE, rhino.Token_DO:
		p.TryMinimizeCondition(NodeUtil_GetConditionExpression(node))
		return node
	case rhino.Token_FOR:
		p.TryJoinForCondition(node)
		p.TryMinimizeCondition(NodeUtil_GetConditionExpression(node))
		return node
	case rhino.Token_BLOCK:
		return p.TryReplaceIf(node)
	default:
		return node
		// Nothing changed
	}
}

// Java: PeepholeMinimizeConditions.tryJoinForCondition(Node)
func (p *PeepholeMinimizeConditions) TryJoinForCondition(n *rhino.Node) {
	if !p.late {
		return
	}

	block := n.GetLastChild()
	maybeIf := block.GetFirstChild()
	if maybeIf != nil && maybeIf.IsIf() {
		thenBlock := maybeIf.GetSecondChild()
		maybeBreak := thenBlock.GetFirstChild()
		if maybeBreak != nil && maybeBreak.IsBreak() &&
			!maybeBreak.HasChildren() {

			// Preserve the IF ELSE expression is there is one.
			if maybeIf.HasXChildren(3) {
				maybeIf.ReplaceWith(maybeIf.GetLastChild().Detach())
			} else {
				NodeUtil_RedeclareVarsInsideBranch(thenBlock)
				block.RemoveFirstChild()
			}

			ifCondition := maybeIf.RemoveFirstChild()
			fixedIfCondition := rhino.IR_Not(ifCondition).
				Srcref(ifCondition)

			// OK, join the IF expression with the FOR expression
			forCondition := NodeUtil_GetConditionExpression(n)
			if forCondition.IsEmpty() {
				forCondition.ReplaceWith(fixedIfCondition)
				p.ReportChangeToEnclosingScope(fixedIfCondition)
			} else {
				replacement := rhino.NewNode(rhino.Token_AND)
				forCondition.ReplaceWith(replacement)
				replacement.AddChildToBack(forCondition)
				replacement.AddChildToBack(fixedIfCondition)
				p.ReportChangeToEnclosingScope(replacement)
			}
		}
	}
}

// TryReplaceIf uses "return x?1:2;" in place of "if(x)return 1;return 2;".
//
// Java: PeepholeMinimizeConditions.tryReplaceIf(Node)
func (p *PeepholeMinimizeConditions) TryReplaceIf(n *rhino.Node) *rhino.Node {

	var next *rhino.Node
	for child := n.GetFirstChild(); child != nil; child = next {
		next = child.GetNext()
		if child.IsIf() {
			cond := child.GetFirstChild()
			thenBranch := cond.GetNext()
			elseBranch := thenBranch.GetNext()
			nextNode := child.GetNext()

			if nextNode != nil && elseBranch == nil &&
				PeepholeMinimizeConditions_IsReturnBlock(thenBranch) &&
				nextNode.IsIf() {
				nextCond := nextNode.GetFirstChild()
				nextThen := nextCond.GetNext()
				nextElse := nextThen.GetNext()
				if p.AreNodesEqualForInlining(thenBranch, nextThen) {
					// Transform
					//   if (x) return 1; if (y) return 1;
					// to
					//   if (x||y) return 1;
					child.Detach()
					child.DetachChildren()
					newCond := rhino.NewNodeWithChild(rhino.Token_OR, cond)
					nextCond.ReplaceWith(newCond)
					newCond.AddChildToBack(nextCond)
					p.ReportChangeToEnclosingScope(newCond)
				} else if nextElse != nil && p.AreNodesEqualForInlining(thenBranch, nextElse) {
					// Transform
					//   if (x) return 1; if (y) foo() else return 1;
					// to
					//   if (!x&&y) foo() else return 1;
					child.Detach()
					child.DetachChildren()
					newCond := rhino.NewNodeWithChild(rhino.Token_AND,
						rhino.IR_Not(cond).Srcref(cond))
					nextCond.ReplaceWith(newCond)
					newCond.AddChildToBack(nextCond)
					p.ReportChangeToEnclosingScope(newCond)
				}
			} else if nextNode != nil && elseBranch == nil &&
				PeepholeMinimizeConditions_IsReturnBlock(thenBranch) && PeepholeMinimizeConditions_IsReturnExpression(nextNode) {
				var thenExpr *rhino.Node
				// if(x)return; return 1 -> return x?void 0:1
				if PeepholeMinimizeConditions_IsReturnExpressBlock(thenBranch) {
					thenExpr = PeepholeMinimizeConditions_GetBlockReturnExpression(thenBranch)
					thenExpr.Detach()
				} else {
					thenExpr = NodeUtil_NewUndefinedNode(child)
				}

				elseExpr := nextNode.GetFirstChild()

				cond.Detach()
				elseExpr.Detach()

				returnNode := rhino.IR_ReturnNodeWithExpr(
					rhino.IR_Hook(cond, thenExpr, elseExpr).
						Srcref(child))
				child.ReplaceWith(returnNode)
				nextNode.Detach()
				p.ReportChangeToEnclosingScope(n)
				// everything else in the block is dead code.
				break
			} else if elseBranch != nil && PeepholeMinimizeConditions_StatementMustExitParent(thenBranch) {
				elseBranch.Detach()
				elseBranch.InsertAfter(child)
				p.ReportChangeToEnclosingScope(n)
			}
		}
	}
	return n
}

// Java: PeepholeMinimizeConditions.statementMustExitParent(Node)
func PeepholeMinimizeConditions_StatementMustExitParent(n *rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_THROW, rhino.Token_RETURN:
		return true
	case rhino.Token_BLOCK:
		if n.HasChildren() {
			child := n.GetLastChild()
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
func (p *PeepholeMinimizeConditions) TryReplaceExitWithBreak(n *rhino.Node) *rhino.Node {
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

	if follow == nil && (n.IsThrow() || result != nil) {
		// Can't complete remove a throw here or a return with a result.
		return n
	}

	// When follow is null, this mean the follow of a break target is the
	// end of a function. This means a break is same as return.
	if follow == nil || p.AreMatchingExits(n, follow) {
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
func (p *PeepholeMinimizeConditions) TryRemoveRedundantExit(n *rhino.Node) *rhino.Node {
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

	if follow == nil && (n.IsThrow() || exitExpr != nil) {
		// Can't complete remove a throw here or a return with a result.
		return n
	}

	// When follow is null, this mean the follow of a break target is the
	// end of a function. This means a break is same as return.
	if follow == nil || p.AreMatchingExits(n, follow) {
		p.ReportChangeToEnclosingScope(n)
		n.Detach()
		return nil
	}

	return n
}

// IsPure returns whether the expression does not produces and can not be
// affected by side-effects.
//
// Java: PeepholeMinimizeConditions.isPure(Node)
func (p *PeepholeMinimizeConditions) IsPure(n *rhino.Node) bool {
	return n == nil ||
		(!NodeUtil_CanBeSideEffected(n) &&
			!p.MayHaveSideEffects(n))
}

// PeepholeMinimizeConditions_SkipFinallyNodes returns n or the node
// following any following finally nodes.
//
// Java: PeepholeMinimizeConditions.skipFinallyNodes(Node)
func PeepholeMinimizeConditions_SkipFinallyNodes(n *rhino.Node) *rhino.Node {
	for n != nil && NodeUtil_IsTryFinallyNode(n.GetParent(), n) {
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
func (p *PeepholeMinimizeConditions) AreMatchingExits(nodeThis *rhino.Node, nodeThat *rhino.Node) bool {
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
func PeepholeMinimizeConditions_IsExceptionPossible(n *rhino.Node) bool {
	// TODO(johnlenz): maybe use ControlFlowAnalysis.mayThrowException?
	jutil.CheckStateMsg1(n.IsReturn() || n.IsThrow(), "%s", n)
	return n.IsThrow() ||
		(n.HasChildren() &&
			!NodeUtil_IsLiteralValue(n.GetLastChild(), true))
}

// Java: PeepholeMinimizeConditions.getExceptionHandler(Node)
func PeepholeMinimizeConditions_GetExceptionHandler(n *rhino.Node) *rhino.Node {
	return ControlFlowAnalysis_GetExceptionHandler(n)
}

// TryMinimizeNot tries to minimize NOT nodes such as !(x==y).
//
// # Returns the replacement for n or the original if no change was made
//
// Java: PeepholeMinimizeConditions.tryMinimizeNot(Node)
func (p *PeepholeMinimizeConditions) TryMinimizeNot(n *rhino.Node) *rhino.Node {
	jutil.CheckArgument(n.IsNot())
	parent := n.GetParent()

	notChild := n.GetFirstChild()
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
	// GT, GE, LT, LE are not handled in this because !(x<NaN) != x>=NaN.
	default:
		return n
	}
	newOperator := n.RemoveFirstChild()
	newOperator.SetToken(complementOperator)
	n.ReplaceWith(newOperator)
	p.ReportChangeToEnclosingScope(parent)
	return newOperator
}

func peepholeMinimizeConditions_canMinimizeConditionNode(n *rhino.Node) bool {
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
func (p *PeepholeMinimizeConditions) TryMinimizeExprResult(n *rhino.Node) *rhino.Node {
	originalExpr := n.GetFirstChild()
	if !peepholeMinimizeConditions_canMinimizeConditionNode(originalExpr) {
		return n
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
	return n
}

// TryMinimizeHook tries flipping HOOKs that have negated conditions.
//
// Returns the replacement for n or the original if no replacement was
// necessary.
//
// Java: PeepholeMinimizeConditions.tryMinimizeHook(Node)
func (p *PeepholeMinimizeConditions) TryMinimizeHook(n *rhino.Node) *rhino.Node {
	originalCond := n.GetFirstChild()
	if !peepholeMinimizeConditions_canMinimizeConditionNode(originalCond) {
		return n
	}
	minCond := MinimizedCondition_FromConditionNodeInto(originalCond, &p.scratchMC)
	mNode :=
		minCond.GetMinimized(MinimizedCondition_MinimizationStyle_ALLOW_LEADING_NOT)
	if mNode.IsNot() {
		// Swap the HOOK
		thenBranch := n.GetSecondChild()
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
func (p *PeepholeMinimizeConditions) TryMinimizeIf(n *rhino.Node) *rhino.Node {

	parent := n.GetParent()

	originalCond := n.GetFirstChild()

	/* If the condition is a literal, we'll let other
	 * optimizations try to remove useless code.
	 */
	if NodeUtil_IsLiteralValue(originalCond, true) {
		return n
	}

	thenBranch := originalCond.GetNext()
	elseBranch := thenBranch.GetNext()

	minCond := MinimizedCondition_FromConditionNodeInto(originalCond, &p.scratchMC)

	// Compute two minimized representations. The first representation counts
	// a leading NOT node, and the second ignores a leading NOT node.
	// If we can fold the if statement into a HOOK or boolean operation,
	// then the NOT node does not matter, and we prefer the second condition.
	// If we cannot fold the if statement, then we prefer the first condition.
	unnegatedCond := minCond.GetMinimized(MinimizedCondition_MinimizationStyle_PREFER_UNNEGATED)
	shortCond := minCond.GetMinimized(MinimizedCondition_MinimizationStyle_ALLOW_LEADING_NOT)

	if elseBranch == nil {
		if PeepholeMinimizeConditions_IsFoldableExpressBlock(thenBranch) {
			expr := PeepholeMinimizeConditions_GetBlockExpression(thenBranch)
			if !p.late && PeepholeMinimizeConditions_IsPropertyAssignmentInExpression(expr) {
				// Keep opportunities for CollapseProperties such as
				// a.longIdentifier || a.longIdentifier = ... -> var a = ...;
				// until CollapseProperties has been run.
				p.ReplaceNode(originalCond, unnegatedCond)
				return n
			}

			if shortCond.IsNot() {
				// if(!x)bar(); -> x||bar();
				replacementCond := p.ReplaceNode(originalCond, shortCond.WithoutNot()).Detach()
				or := rhino.IR_Or(
					replacementCond,
					expr.RemoveFirstChild()).Srcref(n)
				newExpr := NodeUtil_NewExpr(or)
				n.ReplaceWith(newExpr)
				p.ReportChangeToEnclosingScope(parent)

				return newExpr
			}
			// True, but removed for performance reasons.
			// Preconditions.checkState(shortCond.isEquivalentTo(unnegatedCond));

			// if(x)foo(); -> x&&foo();
			if shortCond.IsLowerPrecedenceThan(peepholeMinimizeConditions_AND_PRECEDENCE) &&
				PeepholeMinimizeConditions_IsLowerPrecedence(expr.GetFirstChild(), peepholeMinimizeConditions_AND_PRECEDENCE) {
				// One additional set of parentheses is worth the change even if
				// there is no immediate code size win. However, two extra pair of
				// {}, we would have to think twice. (unless we know for sure the
				// we can further optimize its parent.
				p.ReplaceNode(originalCond, shortCond)
				return n
			}

			replacementCond := p.ReplaceNode(originalCond, shortCond).Detach()
			and := rhino.IR_And(replacementCond, expr.RemoveFirstChild()).Srcref(n)
			newExpr := NodeUtil_NewExpr(and)
			n.ReplaceWith(newExpr)
			p.ReportChangeToEnclosingScope(parent)

			return newExpr
		} else {

			// Try to combine two IF-ELSE
			if NodeUtil_IsStatementBlock(thenBranch) && thenBranch.HasOneChild() {
				innerIf := thenBranch.GetFirstChild()

				if innerIf.IsIf() {
					innerCond := innerIf.GetFirstChild()
					innerThenBranch := innerCond.GetNext()
					innerElseBranch := innerThenBranch.GetNext()

					if innerElseBranch == nil &&
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
		}
		p.ReplaceNode(originalCond, unnegatedCond)
		return n
	}

	/* TODO(dcc) This modifies the siblings of n, which is undesirable for a
	 * peephole optimization. This should probably get moved to another pass.
	 */
	p.TryRemoveRepeatedStatements(n)

	// if(!x)foo();else bar(); -> if(x)bar();else foo();
	// An additional set of curly braces isn't worth it.
	if shortCond.IsNot() && !PeepholeMinimizeConditions_ConsumesDanglingElse(elseBranch) {
		p.ReplaceNode(originalCond, shortCond.WithoutNot())
		thenBranch.Detach()
		n.AddChildToBack(thenBranch)
		p.ReportChangeToEnclosingScope(n)
		return n
	}

	// if(x)return 1;else return 2; -> return x?1:2;
	if PeepholeMinimizeConditions_IsReturnExpressBlock(thenBranch) && PeepholeMinimizeConditions_IsReturnExpressBlock(elseBranch) {
		thenExpr := PeepholeMinimizeConditions_GetBlockReturnExpression(thenBranch)
		elseExpr := PeepholeMinimizeConditions_GetBlockReturnExpression(elseBranch)

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

	thenBranchIsExpressionBlock := PeepholeMinimizeConditions_IsFoldableExpressBlock(thenBranch)
	elseBranchIsExpressionBlock := PeepholeMinimizeConditions_IsFoldableExpressBlock(elseBranch)

	if thenBranchIsExpressionBlock && elseBranchIsExpressionBlock {
		thenOp := PeepholeMinimizeConditions_GetBlockExpression(thenBranch).GetFirstChild()
		elseOp := PeepholeMinimizeConditions_GetBlockExpression(elseBranch).GetFirstChild()
		if thenOp.GetToken() == elseOp.GetToken() {
			// if(x)a=1;else a=2; -> a=x?1:2;
			if NodeUtil_IsAssignmentOp(thenOp) {
				lhs := thenOp.GetFirstChild()
				if p.AreNodesEqualForInlining(lhs, elseOp.GetFirstChild()) &&
					// if LHS has side effects, don't proceed [since the optimization
					// evaluates LHS before cond]
					// NOTE - there are some circumstances where we can
					// proceed even if there are side effects...
					!p.MayEffectMutableState(lhs) &&
					(!p.MayHaveSideEffects(originalCond) ||
						(thenOp.IsAssign() && thenOp.GetFirstChild().IsName())) {

					replacementCond := p.ReplaceNode(originalCond, shortCond).Detach()
					assignName := thenOp.RemoveFirstChild()
					thenExpr := thenOp.RemoveFirstChild()
					elseExpr := elseOp.GetLastChild()
					elseExpr.Detach()

					hookNode := rhino.IR_Hook(replacementCond, thenExpr, elseExpr).
						Srcref(n)
					assign := rhino.NewNodeWith2Children(thenOp.GetToken(), assignName, hookNode).Srcref(thenOp)
					expr := NodeUtil_NewExpr(assign)
					n.ReplaceWith(expr)
					p.ReportChangeToEnclosingScope(parent)

					return expr
				}
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

	thenBranchIsVar := PeepholeMinimizeConditions_IsVarBlock(thenBranch)
	elseBranchIsVar := PeepholeMinimizeConditions_IsVarBlock(elseBranch)

	// if(x)var y=1;else y=2  ->  var y=x?1:2
	if thenBranchIsVar && elseBranchIsExpressionBlock &&
		PeepholeMinimizeConditions_GetBlockExpression(elseBranch).GetFirstChild().IsAssign() {

		var_ := PeepholeMinimizeConditions_GetBlockVar(thenBranch)
		elseAssign := PeepholeMinimizeConditions_GetBlockExpression(elseBranch).GetFirstChild()

		name1 := var_.GetFirstChild()
		maybeName2 := elseAssign.GetFirstChild()

		if name1.HasChildren() &&
			maybeName2.IsName() &&
			name1.GetString() == maybeName2.GetString() {
			jutil.CheckState(name1.HasOneChild())
			thenExpr := name1.RemoveFirstChild()
			elseExpr := elseAssign.GetLastChild().Detach()
			replacementCond := p.ReplaceNode(originalCond, shortCond).Detach()
			hookNode := rhino.IR_Hook(replacementCond, thenExpr, elseExpr).Srcref(n)
			var_.Detach()
			name1.AddChildToBack(hookNode)
			n.ReplaceWith(var_)
			p.ReportChangeToEnclosingScope(parent)
			return var_
		}

		// if(x)y=1;else var y=2  ->  var y=x?1:2
	} else if elseBranchIsVar && thenBranchIsExpressionBlock &&
		PeepholeMinimizeConditions_GetBlockExpression(thenBranch).GetFirstChild().IsAssign() {

		var_ := PeepholeMinimizeConditions_GetBlockVar(elseBranch)
		thenAssign := PeepholeMinimizeConditions_GetBlockExpression(thenBranch).GetFirstChild()

		maybeName1 := thenAssign.GetFirstChild()
		name2 := var_.GetFirstChild()

		if name2.HasChildren() &&
			maybeName1.IsName() &&
			maybeName1.GetString() == name2.GetString() {
			thenExpr := thenAssign.GetLastChild().Detach()
			jutil.CheckState(name2.HasOneChild())
			elseExpr := name2.RemoveFirstChild()
			replacementCond := p.ReplaceNode(originalCond, shortCond).Detach()
			hookNode := rhino.IR_Hook(replacementCond, thenExpr, elseExpr).Srcref(n)
			var_.Detach()
			name2.AddChildToBack(hookNode)
			n.ReplaceWith(var_)
			p.ReportChangeToEnclosingScope(parent)

			return var_
		}
	}

	p.ReplaceNode(originalCond, unnegatedCond)
	return n
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
func (p *PeepholeMinimizeConditions) TryRemoveRepeatedStatements(n *rhino.Node) {
	// Only run this if variable names are guaranteed to be unique. Otherwise bad things can happen:
	// see PeepholeMinimizeConditionsTest#testDontRemoveDuplicateStatementsWithoutNormalization
	if !p.IsASTNormalized() {
		return
	}
	jutil.CheckStateMsg1(n.IsIf(), "%s", n)

	parent := n.GetParent()
	if !NodeUtil_IsStatementBlock(parent) {
		// If the immediate parent is something like a label, we
		// can't move the statement, so bail.
		return
	}

	cond := n.GetFirstChild()
	trueBranch := cond.GetNext()
	falseBranch := trueBranch.GetNext()
	jutil.CheckNotNull(trueBranch)
	jutil.CheckNotNull(falseBranch)

	for {
		lastTrue := trueBranch.GetLastChild()
		lastFalse := falseBranch.GetLastChild()
		if lastTrue == nil || lastFalse == nil ||
			!p.AreNodesEqualForInlining(lastTrue, lastFalse) {
			break
		}
		lastTrue.Detach()
		lastFalse.Detach()
		lastTrue.InsertAfter(n)
		p.ReportChangeToEnclosingScope(parent)
	}
}

// PeepholeMinimizeConditions_IsFoldableExpressBlock returns whether the node
// is a block with a single statement that is an expression.
//
// Java: PeepholeMinimizeConditions.isFoldableExpressBlock(Node)
func PeepholeMinimizeConditions_IsFoldableExpressBlock(n *rhino.Node) bool {
	if n.IsBlock() {
		if n.HasOneChild() {
			maybeExpr := n.GetFirstChild()
			if maybeExpr.IsExprResult() {
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
				if maybeExpr.GetFirstChild().IsCall() || maybeExpr.GetFirstChild().IsOptChainCall() {
					calledFn := maybeExpr.GetFirstFirstChild()

					// We only have to worry about methods with an implicit 'this'
					// param, or this doesn't happen.
					if calledFn.IsGetElem() || calledFn.IsOptChainGetElem() {
						return false
					} else if (calledFn.IsGetProp() || calledFn.IsOptChainGetProp()) &&
						jutil.StringStartsWith(calledFn.GetString(), "on") {
						return false
					}
				}

				return true
			}
			return false
		}
	}

	return false
}

// PeepholeMinimizeConditions_GetBlockExpression returns the expression node.
//
// Java: PeepholeMinimizeConditions.getBlockExpression(Node)
func PeepholeMinimizeConditions_GetBlockExpression(n *rhino.Node) *rhino.Node {
	jutil.CheckState(PeepholeMinimizeConditions_IsFoldableExpressBlock(n))
	return n.GetFirstChild()
}

// PeepholeMinimizeConditions_IsReturnBlock returns whether the node is a
// block with a single statement that is an return with or without an
// expression.
//
// Java: PeepholeMinimizeConditions.isReturnBlock(Node)
func PeepholeMinimizeConditions_IsReturnBlock(n *rhino.Node) bool {
	if n.IsBlock() {
		if n.HasOneChild() {
			first := n.GetFirstChild()
			return first.IsReturn()
		}
	}

	return false
}

// PeepholeMinimizeConditions_IsReturnExpressBlock returns whether the node is
// a block with a single statement that is an return.
//
// Java: PeepholeMinimizeConditions.isReturnExpressBlock(Node)
func PeepholeMinimizeConditions_IsReturnExpressBlock(n *rhino.Node) bool {
	if n.IsBlock() {
		if n.HasOneChild() {
			first := n.GetFirstChild()
			if first.IsReturn() {
				return first.HasOneChild()
			}
		}
	}

	return false
}

// PeepholeMinimizeConditions_IsReturnExpression returns whether the node is a
// single return statement.
//
// Java: PeepholeMinimizeConditions.isReturnExpression(Node)
func PeepholeMinimizeConditions_IsReturnExpression(n *rhino.Node) bool {
	if n.IsReturn() {
		return n.HasOneChild()
	}
	return false
}

// PeepholeMinimizeConditions_GetBlockReturnExpression returns the expression
// that is part of the return.
//
// Java: PeepholeMinimizeConditions.getBlockReturnExpression(Node)
func PeepholeMinimizeConditions_GetBlockReturnExpression(n *rhino.Node) *rhino.Node {
	jutil.CheckState(PeepholeMinimizeConditions_IsReturnExpressBlock(n))
	return n.GetFirstFirstChild()
}

// PeepholeMinimizeConditions_IsVarBlock returns whether the node is a block
// with a single statement that is a VAR declaration of a single variable.
//
// Java: PeepholeMinimizeConditions.isVarBlock(Node)
func PeepholeMinimizeConditions_IsVarBlock(n *rhino.Node) bool {
	if n.IsBlock() {
		if n.HasOneChild() {
			first := n.GetFirstChild()
			if first.IsVar() {
				return first.HasOneChild()
			}
		}
	}

	return false
}

// PeepholeMinimizeConditions_GetBlockVar returns the var node.
//
// Java: PeepholeMinimizeConditions.getBlockVar(Node)
func PeepholeMinimizeConditions_GetBlockVar(n *rhino.Node) *rhino.Node {
	jutil.CheckState(PeepholeMinimizeConditions_IsVarBlock(n))
	return n.GetFirstChild()
}

// PeepholeMinimizeConditions_ConsumesDanglingElse reports whether a statement
// consumes a 'dangling else'. A statement consumes a 'dangling else' if an
// 'else' token following the statement would be considered by the parser to
// be part of the statement.
//
// Java: PeepholeMinimizeConditions.consumesDanglingElse(Node)
func PeepholeMinimizeConditions_ConsumesDanglingElse(n *rhino.Node) bool {
	for {
		switch n.GetToken() {
		case rhino.Token_IF:
			if n.GetChildCount() < 3 {
				return true
			}
			// This IF node has no else clause.
			n = n.GetLastChild()
			continue
		case rhino.Token_BLOCK:
			if !n.HasOneChild() {
				return false
			}
			// This BLOCK has no curly braces.
			n = n.GetLastChild()
			continue
		case rhino.Token_WITH, rhino.Token_WHILE, rhino.Token_FOR, rhino.Token_FOR_IN:
			n = n.GetLastChild()
			continue
		default:
			return false
		}
	}
}

// PeepholeMinimizeConditions_IsLowerPrecedence reports whether the node type
// has lower precedence than "precedence".
//
// Java: PeepholeMinimizeConditions.isLowerPrecedence(Node, int)
func PeepholeMinimizeConditions_IsLowerPrecedence(n *rhino.Node, precedence int) bool {
	return NodeUtil_Precedence(n.GetToken()) < precedence
}

func peepholeMinimizeConditions_isPropertyAssignmentInExpressionPredicate(input *rhino.Node) bool {
	return input.IsGetProp() && input.GetParent().IsAssign()
}

// PeepholeMinimizeConditions_IsPropertyAssignmentInExpression reports
// whether the expression contains a property assignment.
//
// Java: PeepholeMinimizeConditions.isPropertyAssignmentInExpression(Node)
func PeepholeMinimizeConditions_IsPropertyAssignmentInExpression(n *rhino.Node) bool {
	return NodeUtil_Has(n, peepholeMinimizeConditions_isPropertyAssignmentInExpressionPredicate, NodeUtil_MATCH_NOT_FUNCTION)
}

// TryMinimizeCondition tries to minimize condition expression, as there are
// additional assumptions that can be made when it is known that the final
// result is a boolean.
//
// Returns the replacement for n, or the original if no change was made.
//
// Java: PeepholeMinimizeConditions.tryMinimizeCondition(Node)
func (p *PeepholeMinimizeConditions) TryMinimizeCondition(n *rhino.Node) *rhino.Node {
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
func (p *PeepholeMinimizeConditions) ReplaceNode(original *rhino.Node, measuredNodeReplacement *MinimizedCondition_MeasuredNode) *rhino.Node {
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
func (p *PeepholeMinimizeConditions) PerformConditionSubstitutions(n *rhino.Node) *rhino.Node {
	parent := n.GetParent()

	switch n.GetToken() {
	case rhino.Token_OR, rhino.Token_AND:
		left := n.GetFirstChild()
		right := n.GetLastChild()

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
			type_ := n.GetToken()
			var replacement *rhino.Node
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

			if replacement != nil {
				n.DetachChildren()
				n.ReplaceWith(replacement)
				p.ReportChangeToEnclosingScope(parent)
				return replacement
			}
		}
		return n
	case rhino.Token_HOOK:
		condition := n.GetFirstChild()
		trueNode := n.GetSecondChild()
		falseNode := n.GetLastChild()

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
		var replacement *rhino.Node
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

		if replacement != nil {
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
func (p *PeepholeMinimizeConditions) MaybeReplaceChildWithNumber(n *rhino.Node, num int) *rhino.Node {
	newNode := rhino.IR_Number(float64(num))
	if !newNode.IsEquivalentTo(n) {
		n.ReplaceWith(newNode)
		p.ReportChangeToEnclosingScope(newNode)
		p.MarkFunctionsDeleted(n)

		return newNode
	}

	return n
}
