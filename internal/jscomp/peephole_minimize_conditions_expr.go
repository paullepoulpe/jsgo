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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

var peepholeMinimizeConditions_AND_PRECEDENCE = NodeUtil_Precedence(rhino.Token_AND)

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
		innerNot := ast.MustUnary(notChild)
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
		thenVal, thenOk := p.GetSideEffectFreeBooleanValue(thenBranch)
		elseVal, elseOk := p.GetSideEffectFreeBooleanValue(elseBranch)
		var replacement rhino.Node
		if thenOk && elseOk && thenVal && !elseVal {
			// x ? true : false -> x (if boolean) or !!x
			cond := originalCond.Detach()
			if NodeUtil_IsBooleanResult(cond) {
				replacement = cond
			} else {
				replacement = rhino.IR_Not(rhino.IR_Not(cond).Srcref(n)).Srcref(n)
			}
		} else if thenOk && elseOk && !thenVal && elseVal {
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

			if hasValue && isName && name1.Node().GetStringID() == maybeName2.Node().GetStringID() {
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

			if hasValue && isName && maybeName1.Node().GetStringID() == name2.Node().GetStringID() {
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
			if !isGetProp || strings.HasPrefix(getProp.Property(), "on") {
				return ast.ExprStmt{}, false
			}
		}
	}
	return maybeExpr, true
}

// PeepholeMinimizeConditions_IsLowerPrecedence reports whether the node type
// has lower precedence than "precedence".
//
// Java: PeepholeMinimizeConditions.isLowerPrecedence(Node, int)
func PeepholeMinimizeConditions_IsLowerPrecedence(n rhino.Node, precedence int) bool {
	return NodeUtil_Precedence(n.GetToken()) < precedence
}

// PeepholeMinimizeConditions_IsPropertyAssignmentInExpression reports
// whether the expression contains a property assignment.
//
// Java: PeepholeMinimizeConditions.isPropertyAssignmentInExpression(Node)
func PeepholeMinimizeConditions_IsPropertyAssignmentInExpression(n rhino.Node) bool {
	if n.IsGetProp() && n.GetParent().IsAssign() {
		return true
	}
	if n.IsFunction() {
		return false
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if PeepholeMinimizeConditions_IsPropertyAssignmentInExpression(c) {
			return true
		}
	}
	return false
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
func (p *PeepholeMinimizeConditions) ReplaceNode(original rhino.Node, measuredNodeReplacement MinimizedCondition_MeasuredNode) rhino.Node {
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
		binary := ast.MustBinary(n)
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
		if rval, ok := p.GetSideEffectFreeBooleanValue(right); ok {
			type_ := binary.Token()
			var replacement rhino.Node

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
		hook := ast.MustHook(n)
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
		trueNodeVal, trueOk := p.GetSideEffectFreeBooleanValue(trueNode)
		falseNodeVal, falseOk := p.GetSideEffectFreeBooleanValue(falseNode)
		if trueOk && trueNodeVal && falseOk && !falseNodeVal {
			// Remove useless conditionals, keep the condition
			condition.Detach()
			replacement = condition
		} else if trueOk && !trueNodeVal && falseOk && falseNodeVal {
			// Remove useless conditionals, keep the condition
			condition.Detach()
			replacement = rhino.IR_Not(condition)
		} else if trueOk && trueNodeVal {
			// Remove useless true case.
			n.DetachChildren()
			replacement = rhino.IR_Or(condition, falseNode)
		} else if falseOk && !falseNodeVal {
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
		if result, ok := p.GetSideEffectFreeBooleanValue(n); ok {
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
	if n.IsNumber() && n.GetDouble() == float64(num) && n.GetDeclaredTypeExpression().IsNil() {
		return n
	}
	newNode := rhino.IR_Number(float64(num))
	if !newNode.IsEquivalentTo(n) {
		n.ReplaceWith(newNode)
		p.ReportChangeToEnclosingScope(newNode)
		p.MarkFunctionsDeleted(n)

		return newNode
	}

	return n
}
