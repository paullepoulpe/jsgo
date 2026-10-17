// Port of com/google/javascript/jscomp/PeepholeRemoveDeadCode.java @ closure-compiler v20260415

/*
 * Copyright 2004 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// PeepholeRemoveDeadCode is a peephole optimization to remove useless code
// such as IF's with false guard conditions, comma operator's left hand sides
// with no side effects, etc.
type PeepholeRemoveDeadCode struct {
	AbstractPeepholeOptimizationBase
}

// Java: PeepholeRemoveDeadCode()
func NewPeepholeRemoveDeadCode() *PeepholeRemoveDeadCode {
	p := &PeepholeRemoveDeadCode{}
	p.InitAbstractPeepholeOptimizationBase(p) // the implicit super()
	return p
}

// Java: PeepholeRemoveDeadCode.optimizeSubtree(Node)
func (p *PeepholeRemoveDeadCode) OptimizeSubtree(subtree *rhino.Node) *rhino.Node {
	switch subtree.GetToken() {
	case rhino.Token_ASSIGN:
		return p.TryFoldAssignment(subtree)
	case rhino.Token_COMMA:
		return p.TryFoldComma(subtree)
	case rhino.Token_SCRIPT, rhino.Token_BLOCK:
		return p.TryOptimizeBlock(subtree)
	case rhino.Token_EXPR_RESULT:
		return p.TryFoldExpr(subtree)
	case rhino.Token_HOOK:
		return p.TryFoldHook(subtree)
	case rhino.Token_SWITCH:
		return p.TryOptimizeSwitch(subtree)
	case rhino.Token_IF:
		return p.TryFoldIf(subtree)
	case rhino.Token_WHILE:
		// This pass gets run both before and after denormalize. Hence, the AST could potentially
		// contain WHILE (denormalized).
		// TODO: Ideally, we should optimize this case instead of returning
		return subtree
	case rhino.Token_FOR:
		condition := NodeUtil_GetConditionExpression(subtree)
		if condition != nil {
			p.TryFoldForCondition(condition)
		}
		return p.TryFoldFor(subtree)
	case rhino.Token_DO:
		foldedDo := p.TryFoldDoAway(subtree)
		if foldedDo.IsDo() {
			return p.TryFoldEmptyDo(foldedDo)
		}
		return foldedDo
	case rhino.Token_TRY:
		return p.TryFoldTry(subtree)
	case rhino.Token_LABEL:
		return p.TryFoldLabel(subtree)
	case rhino.Token_ARRAY_PATTERN:
		return p.TryOptimizeArrayPattern(subtree)
	case rhino.Token_OBJECT_PATTERN:
		return p.TryOptimizeObjectPattern(subtree)
	case rhino.Token_VAR, rhino.Token_CONST, rhino.Token_LET:
		return p.TryOptimizeNameDeclaration(subtree)
	case rhino.Token_DEFAULT_VALUE:
		return p.TryRemoveDefaultValue(subtree)
	case rhino.Token_OPTCHAIN_GETPROP, rhino.Token_OPTCHAIN_CALL, rhino.Token_OPTCHAIN_GETELEM:
		return p.TryRemoveOptionalChaining(subtree)
	default:
		return subtree
	}
}

// Java: PeepholeRemoveDeadCode.tryRemoveDefaultValue(Node)
func (p *PeepholeRemoveDeadCode) TryRemoveDefaultValue(defaultValue *rhino.Node) *rhino.Node {
	jutil.CheckArgumentMsg1(defaultValue.IsDefaultValue(), "%s", defaultValue)

	lValue := defaultValue.GetFirstChild()
	val := defaultValue.GetSecondChild()
	removeVal := false

	// If the default is `undefined` always remove the value
	if val.IsName() && val.GetString() == "undefined" {
		removeVal = true
	}

	// If the `void` application is pure, remove the value
	if val.IsVoid() {
		voidArg := val.GetFirstChild()
		removeVal = !p.MayHaveSideEffects(voidArg)
	}

	if removeVal {
		defaultValue.ReplaceWith(lValue.Detach())
		p.ReportChangeToEnclosingScope(lValue)
		return lValue
	}

	return defaultValue
}

// Java: PeepholeRemoveDeadCode.tryFoldLabel(Node)
func (p *PeepholeRemoveDeadCode) TryFoldLabel(n *rhino.Node) *rhino.Node {
	labelName := n.GetFirstChild().GetString()
	stmt := n.GetLastChild()
	if stmt.IsEmpty() {
		p.ReportChangeToEnclosingScope(n)
		n.Detach()
		return nil
	}

	if stmt.IsBlock() && !stmt.HasChildren() {
		p.ReportChangeToEnclosingScope(n)
		if n.GetParent().IsLabel() {
			// If the parent is itself a label, replace this label
			// with its contained block to keep the AST in a valid state.
			n.ReplaceWith(stmt.Detach())
		} else {
			n.Detach()
		}
		return nil
	}

	child := PeepholeRemoveDeadCode_GetOnlyInterestingChild(stmt)
	if child != nil {
		stmt = child
	}
	if stmt.IsBreak() && stmt.GetFirstChild().GetString() == labelName {
		if n.GetParent().IsLabel() {
			replacement := rhino.IR_Block().Srcref(n)
			n.ReplaceWith(replacement)
			p.ReportChangeToEnclosingScope(replacement)
			return replacement
		} else {
			parent := n.GetParent()
			n.Detach()
			p.ReportChangeToEnclosingScope(parent)
			return nil
		}
	}
	return n
}

// PeepholeRemoveDeadCode_GetOnlyInterestingChild returns the only
// "interesting" child of block, if it has exactly one interesting child,
// otherwise return null. For purposes of this method, a node is considered
// "interesting" unless it is an empty synthetic block.
//
// Java: PeepholeRemoveDeadCode.getOnlyInterestingChild(Node)
func PeepholeRemoveDeadCode_GetOnlyInterestingChild(block *rhino.Node) *rhino.Node {
	if !block.IsBlock() {
		return nil
	}
	if block.HasOneChild() {
		return block.GetOnlyChild()
	}

	var ret *rhino.Node
	for child := block.GetFirstChild(); child != nil; child = child.GetNext() {
		if child.IsSyntheticBlock() && !child.HasChildren() {
			// Uninteresting child.
		} else if ret != nil {
			// Found more than one interesting child.
			return nil
		} else {
			ret = child
		}
	}
	return ret
}

// TryFoldTry removes try blocks without catch blocks and with empty or not
// existent finally blocks. Or, only leave the finally blocks if try body
// blocks are empty.
//
// Returns the replacement node, if changed, or the original if not.
//
// Java: PeepholeRemoveDeadCode.tryFoldTry(Node)
func (p *PeepholeRemoveDeadCode) TryFoldTry(n *rhino.Node) *rhino.Node {
	jutil.CheckStateMsg1(n.IsTry(), "%s", n)
	body := n.GetFirstChild()
	catchBlock := body.GetNext()
	finallyBlock := catchBlock.GetNext()

	// Removes TRYs that had its CATCH removed and/or empty FINALLY.
	if !catchBlock.HasChildren() && (finallyBlock == nil || !finallyBlock.HasChildren()) {
		jutil.CheckState(!n.GetParent().IsLabel())
		body.Detach()
		n.ReplaceWith(body)
		p.ReportChangeToEnclosingScope(body)
		return body
	}

	// Only leave FINALLYs if TRYs are not empty
	if !body.HasChildren() {
		NodeUtil_RedeclareVarsInsideBranch(catchBlock)
		p.ReportChangeToEnclosingScope(n)
		if finallyBlock != nil {
			finallyBlock.Detach()
			jutil.CheckState(!n.GetParent().IsLabel())
			n.ReplaceWith(finallyBlock)
			return finallyBlock
		} else {
			jutil.CheckState(!n.GetParent().IsLabel())
			n.Detach()
			return nil
		}
	}

	return n
}

// TryFoldAssignment tries removing identity assignments and empty
// destructuring pattern assignments.
//
// Returns the replacement node, if changed, or the original if not.
//
// Java: PeepholeRemoveDeadCode.tryFoldAssignment(Node)
func (p *PeepholeRemoveDeadCode) TryFoldAssignment(subtree *rhino.Node) *rhino.Node {
	jutil.CheckState(subtree.IsAssign())
	left := subtree.GetFirstChild()
	right := subtree.GetLastChild()
	if left.IsName() && right.IsName() && left.GetString() == right.GetString() {
		// Only names
		subtree.ReplaceWith(right.Detach())
		p.ReportChangeToEnclosingScope(right)
		return right
	} else if left.IsDestructuringPattern() && !left.HasChildren() {
		// `[] = <expr>` becomes `<expr>`
		// Note that this does potentially change behavior. If `<expr>` is not iterable and this
		// code originally threw, it will no longer throw.
		subtree.ReplaceWith(right.Detach())
		p.ReportChangeToEnclosingScope(right)
		return right
	}
	return subtree
}

// TryOptimizeNameDeclaration tries removing identity assignments and empty
// destructuring pattern assignments.
//
// Returns the replacement node, if changed, or the original if not.
//
// Java: PeepholeRemoveDeadCode.tryOptimizeNameDeclaration(Node)
func (p *PeepholeRemoveDeadCode) TryOptimizeNameDeclaration(subtree *rhino.Node) *rhino.Node {
	jutil.CheckState(NodeUtil_IsNameDeclaration(subtree))
	left := subtree.GetFirstChild()
	if left.IsDestructuringLhs() && left.HasTwoChildren() {
		pattern := left.GetFirstChild()
		if !pattern.HasChildren() {
			// `var [] = foo();` becomes `foo();`
			value := left.GetSecondChild()
			subtree.ReplaceWith(rhino.IR_ExprResult(value.Detach()).Srcref(value))
			p.ReportChangeToEnclosingScope(value)
		}
	}
	return subtree
}

// TryFoldExpr tries folding EXPR_RESULT nodes by removing useless Ops and
// expressions.
//
// Returns the replacement node, if changed, or the original if not.
//
// Java: PeepholeRemoveDeadCode.tryFoldExpr(Node)
func (p *PeepholeRemoveDeadCode) TryFoldExpr(subtree *rhino.Node) *rhino.Node {
	result := p.TrySimplifyUnusedResult(subtree.GetFirstChild())
	if result == nil {
		parent := subtree.GetParent()
		// If the EXPR_RESULT no longer has any children, remove it as well.
		if parent.IsLabel() {
			replacement := rhino.IR_Block().Srcref(subtree)
			subtree.ReplaceWith(replacement)
			subtree = replacement
		} else {
			subtree.Detach()
			subtree = nil
		}
	}
	return subtree
}

// TrySimplifyUnusedResult replaces expression with an expression that
// contains only side-effects of the original.
//
// This replacement is made under the assumption that the result of
// expression is unused and therefore it is correct to eliminate
// non-side-effectful nodes.
//
// Returns the replacement expression, or null if there were no side-effects
// to preserve.
//
// Java: PeepholeRemoveDeadCode.trySimplifyUnusedResult(Node)
func (p *PeepholeRemoveDeadCode) TrySimplifyUnusedResult(expression *rhino.Node) *rhino.Node {
	sideEffectRoots := &peepholeRemoveDeadCode_nodeDeque{}
	atFixedPoint := p.TrySimplifyUnusedResultInternal(expression, sideEffectRoots)

	if atFixedPoint {
		// `expression` is in a form that cannot be further optimized.
		return expression
	} else if sideEffectRoots.isEmpty() {
		p.DeleteNode(expression)
		return nil
	} else if sideEffectRoots.peekFirst() == expression {
		// Expression was a conditional that was transformed. There can't be any other side-effects,
		// but we also can't detach the transformed root.
		jutil.CheckStateMsg1(sideEffectRoots.size() == 1, "%s", sideEffectRoots)
		p.ReportChangeToEnclosingScope(expression)
		return expression
	} else {
		sideEffects := PeepholeRemoveDeadCode_AsDetachedExpression(sideEffectRoots.pollFirst())

		// Assemble a tree of comma expressions for all the side-effects. The tree must execute the
		// side-effects in FIFO order with respect to the queue. It must also be left leaning to match
		// the parser's preferred structure.
		for !sideEffectRoots.isEmpty() {
			next := PeepholeRemoveDeadCode_AsDetachedExpression(sideEffectRoots.pollFirst())
			sideEffects = rhino.IR_Comma(sideEffects, next).Srcref(next)
		}

		sideEffects.InsertBefore(expression)
		p.DeleteNode(expression)
		return sideEffects
	}
}

// peepholeRemoveDeadCode_nodeDeque is the ArrayDeque<Node> used by
// trySimplifyUnusedResult (only addLast, peekFirst, pollFirst, size and
// isEmpty are used).
//
// Go-only: language-adapter — java.util.ArrayDeque.
type peepholeRemoveDeadCode_nodeDeque struct {
	elems []*rhino.Node
	head  int
}

func (d *peepholeRemoveDeadCode_nodeDeque) addLast(n *rhino.Node) {
	jutil.CheckNotNull(n)
	d.elems = append(d.elems, n)
}

func (d *peepholeRemoveDeadCode_nodeDeque) peekFirst() *rhino.Node {
	if d.head == len(d.elems) {
		return nil
	}
	return d.elems[d.head]
}

func (d *peepholeRemoveDeadCode_nodeDeque) pollFirst() *rhino.Node {
	if d.head == len(d.elems) {
		return nil
	}
	n := d.elems[d.head]
	d.head++
	return n
}

func (d *peepholeRemoveDeadCode_nodeDeque) size() int { return len(d.elems) - d.head }

func (d *peepholeRemoveDeadCode_nodeDeque) isEmpty() bool { return d.head == len(d.elems) }

// String is AbstractCollection.toString().
func (d *peepholeRemoveDeadCode_nodeDeque) String() string {
	return jutil.CollectionToString(d.elems[d.head:])
}

// TrySimplifyUnusedResultInternal collects any potentially side-effectful
// subtrees within tree into sideEffectRoots.
//
// When a node is determined to have side-effects its descendants are not
// explored. This method assumes the entire subtree of such a node must be
// preserved. As a corollary, the contents of sideEffectRoots are a forest.
//
// This operation generally does not mutate tree; however, exceptions are
// made for expressions that alter control-flow. Such expression will be
// pruned of their side-effectless branches. Even in this case, tree is never
// detached.
//
// sideEffectRoots are the roots of subtrees determined to have side-effects,
// in execution order. Returns true iff there is no code to be removed from
// within tree; it is already at a fixed point for code removal.
//
// Java: PeepholeRemoveDeadCode.trySimplifyUnusedResultInternal(Node, ArrayDeque<Node>)
func (p *PeepholeRemoveDeadCode) TrySimplifyUnusedResultInternal(tree *rhino.Node, sideEffectRoots *peepholeRemoveDeadCode_nodeDeque) bool {
	// Special cases for conditional expressions that may be using results.
	switch tree.GetToken() {
	case rhino.Token_HOOK:
		// Try to remove one or more of the conditional children and transform the HOOK to an
		// equivalent operation. Remember that if either value branch still exists, the result of
		// the predicate expression is being used, and so cannot be removed.
		//    x() ? foo() : 1 --> x() && foo()
		//    x() ? 1 : foo() --> x() || foo()
		//    x() ? 1 : 1 --> x()
		//    x ? 1 : 1 --> null
		trueNode := p.TrySimplifyUnusedResult(tree.GetSecondChild())
		falseNode := p.TrySimplifyUnusedResult(tree.GetLastChild())
		if trueNode == nil && falseNode != nil {
			jutil.CheckStateMsg1(tree.HasTwoChildren(), "%s", tree)

			tree.SetToken(rhino.Token_OR)
			sideEffectRoots.addLast(tree)
			return false // The node type was changed.
		} else if trueNode != nil && falseNode == nil {
			jutil.CheckStateMsg1(tree.HasTwoChildren(), "%s", tree)

			tree.SetToken(rhino.Token_AND)
			sideEffectRoots.addLast(tree)
			return false // The node type was changed.
		} else if trueNode == nil && falseNode == nil {
			// Don't bother adding true and false branch children to make the AST valid; this HOOK is
			// going to be deleted. We just need to collect any side-effects from the predicate
			// expression.
			p.TrySimplifyUnusedResultInternal(tree.GetOnlyChild(), sideEffectRoots)
			return false // This HOOK must be cleaned up.
		} else {
			sideEffectRoots.addLast(tree)
			return PeepholeRemoveDeadCode_HasFixedPointParent(tree)
		}
	case rhino.Token_AND, rhino.Token_OR, rhino.Token_COALESCE:
		// Try to remove the second operand from a AND, OR, and COALESCE operations. Remember that
		// if the second
		// child still exists, the result of the first expression is being used, and so cannot be
		// removed.
		//    x() ?? f --> x()
		//    x() || f --> x()
		//    x() && f --> x()
		conditionalResultNode := p.TrySimplifyUnusedResult(tree.GetLastChild())
		if conditionalResultNode == nil {
			// Don't bother adding a second child to make the AST valid; this op is going to be
			// deleted. We just need to collect any side-effects from the predicate first child.
			p.TrySimplifyUnusedResultInternal(tree.GetOnlyChild(), sideEffectRoots)
			return false // This op must be cleaned up.
		} else {
			sideEffectRoots.addLast(tree)
			return PeepholeRemoveDeadCode_HasFixedPointParent(tree)
		}
	case rhino.Token_FUNCTION:
		// Functions that aren't being invoked are dead. If they were invoked we'd see the CALL
		// before arriving here. We don't want to look at any children since they'll never execute.
		return false
	default:
		// This is the meat of this function. It covers the general case of nodes which are unused
		if p.NodeTypeMayHaveSideEffects(tree) {
			sideEffectRoots.addLast(tree)
			return PeepholeRemoveDeadCode_HasFixedPointParent(tree)
		} else if !tree.HasChildren() {
			return false // A node must have children or side-effects to be at fixed-point.
		}

		atFixedPoint := PeepholeRemoveDeadCode_HasFixedPointParent(tree)
		for child := tree.GetFirstChild(); child != nil; child = child.GetNext() {
			atFixedPoint = p.TrySimplifyUnusedResultInternal(child, sideEffectRoots) && atFixedPoint
		}
		return atFixedPoint
	}
}

// PeepholeRemoveDeadCode_AsDetachedExpression returns an expression
// executing expr which is legal in any expression context.
//
// expr is an attached expression. Returns a detached expression.
//
// Java: PeepholeRemoveDeadCode.asDetachedExpression(Node)
func PeepholeRemoveDeadCode_AsDetachedExpression(expr *rhino.Node) *rhino.Node {
	switch expr.GetToken() {
	case rhino.Token_ITER_SPREAD, rhino.Token_OBJECT_SPREAD:
		switch expr.GetParent().GetToken() {
		case rhino.Token_ARRAYLIT, rhino.Token_NEW,
			rhino.Token_CALL,          // `Math.sin(...c)`
			rhino.Token_OPTCHAIN_CALL: // `Math?.sin(...c)`
			expr = rhino.IR_Arraylit(expr.Detach()).Srcref(expr)
		case rhino.Token_OBJECTLIT:
			expr = rhino.IR_Objectlit(expr.Detach()).Srcref(expr)
		default:
			panic(jutil.NewIllegalStateExceptionMsg(expr.ToStringTree()))
		}
	default:
	}

	if expr.HasParent() {
		expr.Detach()
	}

	jutil.CheckStateMsg1(rhino.IR_MayBeExpression(expr), "%s", expr)
	return expr
}

// PeepholeRemoveDeadCode_HasFixedPointParent returns true iff expr is
// parented such that it is valid in a fixed-point representation of an
// unused expression tree.
//
// A fixed-point representation is one in which no further nodes should be
// changed or removed when removing unused code. This method assumes that the
// expression tree in question is unused, so only side-effects are relevant.
//
// Java: PeepholeRemoveDeadCode.hasFixedPointParent(Node)
func PeepholeRemoveDeadCode_HasFixedPointParent(expr *rhino.Node) bool {
	// Most kinds of nodes shouldn't be branches in the fixed-point tree of an unused
	// expression. Those listed below are the only valid kinds.
	switch expr.GetParent().GetToken() {
	case rhino.Token_AND, rhino.Token_COMMA, rhino.Token_HOOK, rhino.Token_OR, rhino.Token_COALESCE:
		return true
	case rhino.Token_ARRAYLIT, rhino.Token_OBJECTLIT:
		// Make a special allowance for SPREADs so they remain in a legal context. Parent types
		// other than ARRAYLIT and OBJECTLIT are not fixed-point because they are the tersest
		// legal
		// parents and are known to be side-effect free.
		return expr.IsSpread()
	default:
		// Statments are always fixed-point parents. All other expressions are not.
		return NodeUtil_IsStatement(expr.GetParent())
	}
}

// peepholeRemoveDeadCode_MatchUnnamedBreak is a predicate for matching
// anything except function nodes.
type peepholeRemoveDeadCode_MatchUnnamedBreak struct {
	// Go-only: language-adapter — distinct identity
	_ byte
}

// Java: PeepholeRemoveDeadCode.MatchUnnamedBreak.apply(Node)
func (m *peepholeRemoveDeadCode_MatchUnnamedBreak) Apply(n *rhino.Node) bool {
	return n.IsBreak() && !n.HasChildren()
}

// Java: PeepholeRemoveDeadCode.MATCH_UNNAMED_BREAK
var PeepholeRemoveDeadCode_MATCH_UNNAMED_BREAK func(*rhino.Node) bool = (&peepholeRemoveDeadCode_MatchUnnamedBreak{}).Apply

// Java: PeepholeRemoveDeadCode.removeIfUnnamedBreak(Node)
func (p *PeepholeRemoveDeadCode) RemoveIfUnnamedBreak(maybeBreak *rhino.Node) {
	if maybeBreak != nil && maybeBreak.IsBreak() && !maybeBreak.HasChildren() {
		p.ReportChangeToEnclosingScope(maybeBreak)
		maybeBreak.Detach()
	}
}

// Java: PeepholeRemoveDeadCode.tryRemoveSwitchWithSingleCase(Node, boolean)
func (p *PeepholeRemoveDeadCode) TryRemoveSwitchWithSingleCase(switchNode *rhino.Node, shouldHoistCondition bool) *rhino.Node {
	switchBody := switchNode.GetSecondChild()
	caseBlock := switchBody.GetOnlyChild().GetLastChild()
	p.RemoveIfUnnamedBreak(caseBlock.GetLastChild())
	// Back off if the switch contains statements like "if (a) { break; }"
	if NodeUtil_Has(caseBlock, PeepholeRemoveDeadCode_MATCH_UNNAMED_BREAK, NodeUtil_MATCH_NOT_FUNCTION) {
		return switchNode
	}
	if shouldHoistCondition {
		rhino.IR_ExprResult(switchNode.RemoveFirstChild()).Srcref(switchNode).InsertBefore(switchNode)
	}
	switchNode.ReplaceWith(caseBlock.Detach())
	p.ReportChangeToEnclosingScope(caseBlock)
	return caseBlock
}

// Java: PeepholeRemoveDeadCode.tryRemoveSwitch(Node)
func (p *PeepholeRemoveDeadCode) TryRemoveSwitch(n *rhino.Node) *rhino.Node {
	switchBody := n.GetSecondChild()
	if !switchBody.HasChildren() {
		// Remove the switch if there are no remaining cases
		condition := n.RemoveFirstChild()
		replacement := rhino.IR_ExprResult(condition).Srcref(n)
		n.ReplaceWith(replacement)
		p.ReportChangeToEnclosingScope(replacement)
		return replacement
	} else if switchBody.HasOneChild() && switchBody.GetOnlyChild().IsDefaultCase() {
		if n.GetFirstChild().IsCall() || n.GetFirstChild().IsOptChainCall() {
			// Before removing switch, we must preserve the switch condition if it is a call
			return p.TryRemoveSwitchWithSingleCase(n, true)
		} else {
			return p.TryRemoveSwitchWithSingleCase(n, false)
		}
	} else {
		return n
	}
}

// TryOptimizeSwitch removes useless switches and cases.
//
// Java: PeepholeRemoveDeadCode.tryOptimizeSwitch(Node)
func (p *PeepholeRemoveDeadCode) TryOptimizeSwitch(n *rhino.Node) *rhino.Node {
	jutil.CheckStateMsg1(n.IsSwitch(), "%s", n)

	switchBody := n.GetSecondChild()
	defaultCase := p.TryOptimizeDefaultCase(n)

	// Generally, it is unsafe to remove other cases when the default case is not the last one.
	if (defaultCase == nil || switchBody.GetLastChild().IsDefaultCase()) &&
		p.AreAllCaseTagsLiterals(switchBody) {
		cond := n.GetFirstChild()
		var prev *rhino.Node
		var next *rhino.Node
		var cur *rhino.Node

		// First, remove empty cases where possible: always empty default cases; or when there is
		// no default case, other empty cases that are not the first matching case, may be
		// removable as well.
		foundMatchingCase := false
		for cur = switchBody.GetFirstChild(); cur != nil; cur = next {
			next = cur.GetNext()
			firstChild := cur.GetFirstChild()
			foundMatchingCase = p.IsFirstSwitchMatch(foundMatchingCase, cond, firstChild)
			if !foundMatchingCase &&
				!p.MayHaveSideEffects(firstChild) &&
				p.IsUselessCase(cur, prev, defaultCase) {
				p.RemoveCase(n, cur)
			} else {
				prev = cur
			}
		}

		// Next, optimize switches with constant condition
		if NodeUtil_IsLiteralValue(cond, false) {
			var caseLabel *rhino.Node
			caseMatches := base.Tri_TRUE
			// Remove cases until you find one that may match
			for cur = switchBody.GetFirstChild(); cur != nil; cur = next {
				next = cur.GetNext()
				caseLabel = cur.GetFirstChild()
				caseMatches = PeepholeFoldConstants_EvaluateComparison(p, rhino.Token_SHEQ, cond, caseLabel)
				if caseMatches == base.Tri_TRUE {
					break
				} else if caseMatches == base.Tri_UNKNOWN {
					break
				} else {
					p.RemoveCase(n, cur)
				}
			}
			if cur != nil && caseMatches == base.Tri_TRUE {
				// Skip cases until you find one whose last stm is a removable break
				matchingCase := cur
				matchingCaseBlock := matchingCase.GetLastChild()
				for cur != nil {
					block := cur.GetLastChild()
					lastStm := block.GetLastChild()
					isLastStmRemovableBreak := false
					if lastStm != nil && PeepholeRemoveDeadCode_IsExit(lastStm) {
						p.RemoveIfUnnamedBreak(lastStm)
						isLastStmRemovableBreak = true
					}
					next = cur.GetNext()
					// Remove the fallthrough case labels
					if cur != matchingCase {
						for block.HasChildren() {
							matchingCaseBlock.AddChildToBack(block.RemoveFirstChild())
						}
						p.ReportChangeToEnclosingScope(cur)
						cur.Detach()
					}
					cur = next
					if isLastStmRemovableBreak {
						break
					}
				}

				// Remove any remaining cases
				for ; cur != nil; cur = next {
					next = cur.GetNext()
					p.RemoveCase(n, cur)
				}
				// If there is one case left, we may be able to fold it
				cur = cond.GetNext()
				if cur != nil && cur.GetNext() == nil {
					return p.TryRemoveSwitchWithSingleCase(n, false)
				}
			}
		}
	}

	// Last, try to remove the entire switch if possible
	return p.TryRemoveSwitch(n)
}

// TryOptimizeDefaultCase returns the default case node or null if there is
// no default case or if the default case is removed.
//
// Java: PeepholeRemoveDeadCode.tryOptimizeDefaultCase(Node)
func (p *PeepholeRemoveDeadCode) TryOptimizeDefaultCase(n *rhino.Node) *rhino.Node {
	jutil.CheckStateMsg1(n.IsSwitch(), "%s", n)

	switchBody := n.GetSecondChild()
	var lastNonRemovable *rhino.Node // the most recently iterated case known to not be removable.

	for c := switchBody.GetFirstChild(); c != nil; c = c.GetNext() {
		if c.IsDefaultCase() {
			// Remove any cases that fall-through to the default case
			var caseToRemove *rhino.Node
			if lastNonRemovable != nil {
				caseToRemove = lastNonRemovable.GetNext()
			} else {
				caseToRemove = switchBody.GetFirstChild()
			}
			for next := (*rhino.Node)(nil); caseToRemove != c; caseToRemove = next {
				next = caseToRemove.GetNext()
				p.RemoveCase(n, caseToRemove)
			}

			// Remove the default case if we can
			if p.IsUselessCase(c, lastNonRemovable, c) {
				p.RemoveCase(n, c)
				return nil
			}
			return c
		} else {
			jutil.CheckState(c.IsCase())
			if c.GetLastChild().HasChildren() || p.MayHaveSideEffects(c.GetFirstChild()) {
				lastNonRemovable = c
			}
		}
	}
	return nil
}

// RemoveCase removes the case from the switch redeclaring any variables
// declared in it.
//
// caseNode is the case to remove.
//
// Java: PeepholeRemoveDeadCode.removeCase(Node, Node)
func (p *PeepholeRemoveDeadCode) RemoveCase(switchNode *rhino.Node, caseNode *rhino.Node) {
	NodeUtil_RedeclareVarsInsideBranch(caseNode)
	caseNode.Detach()
	p.ReportChangeToEnclosingScope(switchNode)
}

// IsUselessCase reports whether the CASE or DEFAULT_CASE block does anything
// useful.
//
// The function assumes that when checking a CASE node there is no
// DEFAULT_CASE node in the SWITCH, or the DEFAULT_CASE is the last case in
// the SWITCH.
//
// Java: PeepholeRemoveDeadCode.isUselessCase(Node, Node, Node)
func (p *PeepholeRemoveDeadCode) IsUselessCase(
	caseNode *rhino.Node, previousCase *rhino.Node, defaultCase *rhino.Node) bool {
	jutil.CheckState(previousCase == nil || previousCase.GetNext() == caseNode)
	// A case isn't useless if a previous case falls through to it unless it happens to be the last
	// case in the switch.
	switchBody := caseNode.GetParent()
	if switchBody.GetLastChild() != caseNode && previousCase != nil {
		previousBlock := previousCase.GetLastChild()
		if !previousBlock.HasChildren() || !PeepholeRemoveDeadCode_IsExit(previousBlock.GetLastChild()) {
			return false
		}
	}

	executingCase := caseNode
	for executingCase != nil {
		jutil.CheckState(executingCase.IsDefaultCase() || executingCase.IsCase())
		// We only expect a DEFAULT case if the case we are checking is the
		// DEFAULT case.  Otherwise, we assume the DEFAULT case has already
		// been removed.
		jutil.CheckState(caseNode == executingCase || !executingCase.IsDefaultCase())
		if !executingCase.IsDefaultCase() && p.MayHaveSideEffects(executingCase.GetFirstChild()) {
			// The case falls thru to a case whose condition has a potential side-effect,
			// removing the candidate case would skip that side-effect, so don't.
			return false
		}
		block := executingCase.GetLastChild()
		jutil.CheckState(block.IsBlock())
		if block.HasChildren() {
			for blockChild := block.GetFirstChild(); blockChild != nil; blockChild = blockChild.GetNext() {
				// If this is a block with a labelless break, it is useless.
				switch blockChild.GetToken() {
				case rhino.Token_BREAK:
					// A case with a single labelless break is useless if it is the default case or if
					// there is no default case. A break to a different control structure isn't useless.
					return !blockChild.HasChildren() &&
						(defaultCase == nil || defaultCase == executingCase)
				case rhino.Token_VAR:
					if blockChild.HasOneChild() && blockChild.GetFirstFirstChild() == nil {
						// Variable declarations without initializations are OK.
						continue
					}
					return false
				default:
					return false
				}
			}
		}
		// Look at the fallthrough case
		executingCase = executingCase.GetNext()
	}
	return true
}

// Java: PeepholeRemoveDeadCode.isFirstSwitchMatch(boolean, Node, Node)
func (p *PeepholeRemoveDeadCode) IsFirstSwitchMatch(foundMatchingCase bool, condition *rhino.Node, tag *rhino.Node) bool {
	if foundMatchingCase {
		return false
	}
	return PeepholeFoldConstants_EvaluateComparison(p, rhino.Token_SHEQ, condition, tag) == base.Tri_TRUE
}

// Java: PeepholeRemoveDeadCode.areAllCaseTagsLiterals(Node)
func (p *PeepholeRemoveDeadCode) AreAllCaseTagsLiterals(switchBody *rhino.Node) bool {
	for caseNode := switchBody.GetFirstChild(); caseNode != nil; caseNode = caseNode.GetNext() {
		if caseNode.IsDefaultCase() {
			continue
		}
		caseTag := caseNode.GetFirstChild()
		if !NodeUtil_IsLiteralValue(caseTag, false) {
			return false
		}
	}
	return true
}

// PeepholeRemoveDeadCode_IsExit returns whether the node is a control flow
// exit from the current block.
//
// Java: PeepholeRemoveDeadCode.isExit(Node)
func PeepholeRemoveDeadCode_IsExit(n *rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_BREAK, rhino.Token_CONTINUE, rhino.Token_RETURN, rhino.Token_THROW:
		return true
	case rhino.Token_SWITCH:
		return PeepholeRemoveDeadCode_IsSwitchExit(n)
	case rhino.Token_TRY:
		return PeepholeRemoveDeadCode_IsTryExit(n)
	default:
		return false
	}
}

// PeepholeRemoveDeadCode_IsUnconditionalBlockExit returns whether the block
// is a control flow exit from the block containing current switch or
// try..catch statement.
//
// Java: PeepholeRemoveDeadCode.isUnconditionalBlockExit(Node)
func PeepholeRemoveDeadCode_IsUnconditionalBlockExit(n *rhino.Node) bool {
	jutil.CheckStateMsg1(n.IsBlock(), "%s", n)
	jutil.CheckStateMsg1(!n.GetParent().IsLabel(), "%s", n)

	// Last statement must lead out of the block.
	lastStm := n.GetLastChild()
	if lastStm == nil {
		return false
	}
	switch lastStm.GetToken() {
	case rhino.Token_BREAK:
		if !lastStm.HasChildren() {
			return false
		}
		// Last statement is OK - continue with checking others.
	case rhino.Token_RETURN, rhino.Token_THROW:
		// Last statement is OK - continue with checking others.
	default:
		return false
	}

	// Other statements can be anything except for unlabeled "break". But for simplicity, don't go
	// into inner blocks and complex constructs - instead, allow only the simplest statements.
	for child := n.GetFirstChild(); child != lastStm; child = child.GetNext() {
		switch child.GetToken() {
		case rhino.Token_BREAK:
			if !child.HasChildren() {
				return false
			}
			// This break is OK - continue with checking others.
		case rhino.Token_RETURN, rhino.Token_THROW, rhino.Token_FUNCTION, rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST, rhino.Token_EXPR_RESULT:
			// This statement is OK - continue with checking others.
		default:
			return false
		}
	}

	return true
}

// PeepholeRemoveDeadCode_IsSwitchExit returns true if the switch always
// "exits" (return, throw, etc).
//
// Java: PeepholeRemoveDeadCode.isSwitchExit(Node)
func PeepholeRemoveDeadCode_IsSwitchExit(n *rhino.Node) bool {
	jutil.CheckStateMsg1(n.IsSwitch(), "%s", n)

	hasDefaultCase := false

	switchBody := n.GetSecondChild()
	for switchCase := switchBody.GetFirstChild(); switchCase != nil; switchCase = switchCase.GetNext() {
		if switchCase.IsDefaultCase() {
			hasDefaultCase = true
		}
		block := switchCase.GetLastChild()
		if (block.HasChildren() || switchCase.GetNext() == nil) &&
			!PeepholeRemoveDeadCode_IsUnconditionalBlockExit(block) {
			return false
		}
	}

	return hasDefaultCase
}

// PeepholeRemoveDeadCode_IsTryExit returns true if the try..catch always
// "exits" (return, throw, etc).
//
// Java: PeepholeRemoveDeadCode.isTryExit(Node)
func PeepholeRemoveDeadCode_IsTryExit(n *rhino.Node) bool {
	jutil.CheckStateMsg1(n.IsTry(), "%s", n)

	// finally - regardless of the behavior of the other blocks,
	// an exit from the finally with guarantee that behavior.
	if n.HasXChildren(3) {
		if PeepholeRemoveDeadCode_IsUnconditionalBlockExit(n.GetLastChild()) {
			return true
		}
	}
	// try
	if !PeepholeRemoveDeadCode_IsUnconditionalBlockExit(n.GetFirstChild()) {
		return false
	}
	// catch
	catches := n.GetSecondChild()
	return !catches.HasChildren() ||
		PeepholeRemoveDeadCode_IsUnconditionalBlockExit(catches.GetOnlyChild().GetLastChild())
}

// Java: PeepholeRemoveDeadCode.tryFoldComma(Node)
func (p *PeepholeRemoveDeadCode) TryFoldComma(n *rhino.Node) *rhino.Node {
	// If the left side does nothing replace the comma with the result.
	parent := n.GetParent()
	left := n.GetFirstChild()
	right := left.GetNext()

	left = p.TrySimplifyUnusedResult(left)
	if left == nil || !p.MayHaveSideEffects(left) {
		// Fold it!
		right.Detach()
		n.ReplaceWith(right)
		p.ReportChangeToEnclosingScope(parent)
		return right
	}
	return n
}

// TryOptimizeBlock tries removing unneeded block nodes and their useless
// children.
//
// Java: PeepholeRemoveDeadCode.tryOptimizeBlock(Node)
func (p *PeepholeRemoveDeadCode) TryOptimizeBlock(n *rhino.Node) *rhino.Node {
	// Remove any useless children
	for c := n.GetFirstChild(); c != nil; {
		next := c.GetNext() // save c.next, since 'c' may be removed
		if !PeepholeRemoveDeadCode_IsUnremovableNode(c) && !p.MayHaveSideEffects(c) {
			// TODO(johnlenz): determine what this is actually removing. Candidates
			//    include: EMPTY nodes, control structures without children
			//    (removing infinite loops), empty try blocks.  What else?
			c.Detach()
			p.ReportChangeToEnclosingScope(n)
			p.MarkFunctionsDeleted(c)
		} else if PeepholeRemoveDeadCode_IsExit(c) {
			p.RemoveFollowingNodes(c)
			break
		} else {
			p.TryOptimizeConditionalAfterAssign(c)
		}
		c = next
	}

	if n.IsSyntheticBlock() || n.IsScript() || n.GetParent() == nil {
		return n
	}

	// Try to merge the block with its parent, or remove it if it is an empty class static block.
	parent := n.GetParent()
	if NodeUtil_TryMergeBlock(n, p.IsASTNormalized()) {
		p.ReportChangeToEnclosingScope(parent)
		return nil
	} else if parent.IsClassMembers() && !n.HasChildren() {
		n.Detach()
		p.ReportChangeToEnclosingScope(parent)
		return nil
	}

	return n
}

// Java: PeepholeRemoveDeadCode.removeFollowingNodes(Node)
func (p *PeepholeRemoveDeadCode) RemoveFollowingNodes(n *rhino.Node) {
	parent := n.GetParent()
	if parent.GetLastChild() != n {
		changed := false
		for n.GetNext() != nil {
			dead := n.GetNext()
			if NodeUtil_IsFunctionDeclaration(dead) {
				// Don't remove function declarations as they are hoisted
				n = dead
				continue
			}
			changed = true
			NodeUtil_RedeclareVarsInsideBranch(dead)
			p.RedeclareIfBlockScopedVar(dead, parent)
			dead.Detach()
			p.MarkFunctionsDeleted(dead)
		}
		if changed {
			p.ReportChangeToEnclosingScope(parent)
		}
	}
}

// RedeclareIfBlockScopedVar redeclares the given node's names in the parent
// scope if they are block-scoped declarations; otherwise does nothing.
//
// Java: PeepholeRemoveDeadCode.redeclareIfBlockScopedVar(Node, Node)
func (p *PeepholeRemoveDeadCode) RedeclareIfBlockScopedVar(decl *rhino.Node, parent *rhino.Node) {
	// This isn't called on function declarations ever: this method is only used in
	// removeFollowingNodes, which preserves function declarations as they're hoisted. This
	// preconditions check is to make sure nothing else tries reusing this code, since in that case
	// maybe they do want to handle function declarations.
	jutil.CheckStateMsg1(!NodeUtil_IsFunctionDeclaration(decl), "Unexpected function declaration %s", decl)
	if !NodeUtil_IsBlockScopedDeclaration(decl) {
		return
	}
	var blockScopedVars []*rhino.Node
	if decl.IsClass() {
		blockScopedVars = append(blockScopedVars, decl.GetFirstChild())
	} else {
		NodeUtil_VisitLhsNodesInNode(
			decl,
			func(n *rhino.Node) {
				if n.IsName() {
					blockScopedVars = append(blockScopedVars, n)
				}
			})
	}

	addBefore := NodeUtil_GetInsertionPointAfterAllInnerFunctionDeclarations(parent)
	for _, nameNode := range blockScopedVars {
		var_ := rhino.IR_Let(rhino.IR_Name(nameNode.GetString()).Srcref(nameNode)).Srcref(nameNode)
		NodeUtil_CopyNameAnnotations(nameNode, var_.GetFirstChild())
		if addBefore != nil {
			var_.InsertBefore(addBefore)
		} else {
			parent.AddChildToBack(var_)
		}
	}
	p.AddFeatureToEnclosingScript(parser.FeatureSet_Feature_LET_DECLARATIONS)
}

// PeepholeRemoveDeadCode_IsUnremovableNode reports whether n is one of the
// nodes that are unremovable but don't have side effects so they aren't
// caught by mayHaveSideEffects.
//
// Java: PeepholeRemoveDeadCode.isUnremovableNode(Node)
func PeepholeRemoveDeadCode_IsUnremovableNode(n *rhino.Node) bool {
	return (n.IsBlock() && n.IsSyntheticBlock()) || n.IsScript()
}

// TryOptimizeConditionalAfterAssign attempts to replace the condition of if
// or hook immediately that is a reference to a name that is assigned
// immediately before.
//
// TODO(johnlenz): Consider moving this to a separate peephole pass.
//
// Java: PeepholeRemoveDeadCode.tryOptimizeConditionalAfterAssign(Node)
func (p *PeepholeRemoveDeadCode) TryOptimizeConditionalAfterAssign(n *rhino.Node) {
	next := n.GetNext()

	// Look for patterns like the following and replace the if-condition with
	// a constant value so it can later be folded:
	//   var a = /a/;
	//   if (a) {foo(a)}
	// or
	//   a = 0;
	//   a ? foo(a) : c;
	// or
	//   a = 0;
	//   a || foo(a);
	// or
	//   a = 0;
	//   a && foo(a)
	// or
	//   a = 0;
	//   a ?? foo(a)
	// TODO(johnlenz): This would be better handled by control-flow sensitive
	// constant propagation. As the other case that I want to handle is:
	//   i=0; for(;i<0;i++){}
	// as right now nothing facilitates removing a loop like that.
	// This is here simply to remove the cruft left behind goog.userAgent and
	// similar cases.

	if !PeepholeRemoveDeadCode_IsSimpleAssignment(n) {
		return
	}
	conditionalRoot := p.GetConditionalRoot(next)
	if conditionalRoot == nil {
		return
	}
	lhsAssign := p.GetSimpleAssignmentName(n)

	condition := p.GetConditionalStatementCondition(next)
	if !lhsAssign.MatchesNameNode(condition) {
		return
	}
	rhsAssign := p.GetSimpleAssignmentValue(n)
	switch conditionalRoot.GetToken() {
	case rhino.Token_AND, rhino.Token_OR, rhino.Token_IF, rhino.Token_HOOK:
		// conditionals that coerce their condition to a boolean
		value := NodeUtil_GetBooleanValue(rhsAssign)
		if value != base.Tri_UNKNOWN {
			replacementConditionNode := NodeUtil_BooleanNode(value.ToBoolean(true))
			condition.ReplaceWith(replacementConditionNode)
			p.ReportChangeToEnclosingScope(replacementConditionNode)
		}
		return
	case rhino.Token_COALESCE:
		// conditional that checks whether its operand is nullish
		valueType := NodeUtil_GetKnownValueType(rhsAssign)
		switch valueType {
		case NodeUtil_ValueType_NULL, NodeUtil_ValueType_VOID:
			condition.ReplaceWith(NodeUtil_NewUndefinedNode(condition))
			p.ReportChangeToEnclosingScope(conditionalRoot)
		case NodeUtil_ValueType_NUMBER, NodeUtil_ValueType_BIGINT, NodeUtil_ValueType_STRING, NodeUtil_ValueType_BOOLEAN, NodeUtil_ValueType_OBJECT:
			// semi-arbitrarily use '0' as a short non-nullish conditional
			condition.ReplaceWith(rhino.IR_Number(0).Srcref(condition))
			p.ReportChangeToEnclosingScope(conditionalRoot)
		case NodeUtil_ValueType_UNDETERMINED:
		}
		return
	default:
		panic(jutil.NewAssertionError("Unhandled condition " + conditionalRoot.String()))
	}
}

// PeepholeRemoveDeadCode_IsSimpleAssignment returns whether the node is a
// assignment to a simple name, or simple var declaration with
// initialization.
//
// Java: PeepholeRemoveDeadCode.isSimpleAssignment(Node)
func PeepholeRemoveDeadCode_IsSimpleAssignment(n *rhino.Node) bool {
	// For our purposes we define a simple assignment to be a assignment
	// to a NAME node, or a VAR declaration with one child and a initializer.
	if NodeUtil_IsExprAssign(n) && n.GetFirstFirstChild().IsName() {
		return true
	} else if NodeUtil_IsNameDeclaration(n) && n.HasOneChild() && n.GetFirstFirstChild() != nil {
		return true
	}

	return false
}

// GetSimpleAssignmentName returns the name being assigned to.
//
// Java: PeepholeRemoveDeadCode.getSimpleAssignmentName(Node)
func (p *PeepholeRemoveDeadCode) GetSimpleAssignmentName(n *rhino.Node) *rhino.Node {
	jutil.CheckState(PeepholeRemoveDeadCode_IsSimpleAssignment(n))
	if NodeUtil_IsExprAssign(n) {
		return n.GetFirstFirstChild()
	} else {
		// A var declaration.
		return n.GetFirstChild()
	}
}

// GetSimpleAssignmentValue returns the value assigned in the simple
// assignment.
//
// Java: PeepholeRemoveDeadCode.getSimpleAssignmentValue(Node)
func (p *PeepholeRemoveDeadCode) GetSimpleAssignmentValue(n *rhino.Node) *rhino.Node {
	jutil.CheckState(PeepholeRemoveDeadCode_IsSimpleAssignment(n))
	return n.GetFirstChild().GetLastChild()
}

// GetConditionalRoot returns the root node if a conditional statement or
// else null.
//
// Java: PeepholeRemoveDeadCode.getConditionalRoot(Node)
func (p *PeepholeRemoveDeadCode) GetConditionalRoot(n *rhino.Node) *rhino.Node {
	// We defined a conditional statement to be a IF or EXPR_RESULT rooted with
	// a HOOK, AND, or OR node.
	if n == nil {
		return nil
	}
	if n.IsIf() {
		return n
	} else if PeepholeRemoveDeadCode_IsExprConditional(n) {
		return n.GetFirstChild()
	}
	return nil
}

// PeepholeRemoveDeadCode_IsExprConditional returns whether the node is a
// rooted with a HOOK, AND, OR, or COALESCE node.
//
// Java: PeepholeRemoveDeadCode.isExprConditional(Node)
func PeepholeRemoveDeadCode_IsExprConditional(n *rhino.Node) bool {
	if n.IsExprResult() {
		switch n.GetFirstChild().GetToken() {
		case rhino.Token_HOOK, rhino.Token_AND, rhino.Token_OR, rhino.Token_COALESCE:
			return true
		default:
		}
	}
	return false
}

// GetConditionalStatementCondition returns the condition of a conditional
// statement.
//
// Java: PeepholeRemoveDeadCode.getConditionalStatementCondition(Node)
func (p *PeepholeRemoveDeadCode) GetConditionalStatementCondition(n *rhino.Node) *rhino.Node {
	if n.IsIf() {
		return NodeUtil_GetConditionExpression(n)
	} else {
		jutil.CheckState(PeepholeRemoveDeadCode_IsExprConditional(n))
		return n.GetFirstFirstChild()
	}
}

// TryFoldIf tries folding IF nodes by removing dead branches.
//
// Returns the replacement node, if changed, or the original if not.
//
// Java: PeepholeRemoveDeadCode.tryFoldIf(Node)
func (p *PeepholeRemoveDeadCode) TryFoldIf(n *rhino.Node) *rhino.Node {
	jutil.CheckStateMsg1(n.IsIf(), "%s", n)
	parent := n.GetParent()
	jutil.CheckNotNull(parent)
	type_ := n.GetToken()
	cond := n.GetFirstChild()
	thenBody := cond.GetNext()
	elseBody := thenBody.GetNext()

	// if (x) { .. } else { } --> if (x) { ... }
	if elseBody != nil && !p.MayHaveSideEffects(elseBody) {
		elseBody.Detach()
		p.ReportChangeToEnclosingScope(n)
		elseBody = nil
	}

	// if (x) { } else { ... } --> if (!x) { ... }
	if !p.MayHaveSideEffects(thenBody) && elseBody != nil {
		elseBody.Detach()
		thenBody.ReplaceWith(elseBody)
		notCond := rhino.NewNode(rhino.Token_NOT)
		cond.ReplaceWith(notCond)
		p.ReportChangeToEnclosingScope(n)
		notCond.AddChildToFront(cond)
		cond = notCond
		thenBody = cond.GetNext()
		elseBody = nil
	}

	// `if (x()) { }` or `if (x?.()) { }`
	if !p.MayHaveSideEffects(thenBody) && elseBody == nil {
		if p.MayHaveSideEffects(cond) {
			// `x()` or `x?.()` has side effects, just leave the condition on its own.
			cond.Detach()
			replacement := NodeUtil_NewExpr(cond)
			n.ReplaceWith(replacement)
			p.ReportChangeToEnclosingScope(parent)
			return replacement
		} else {
			// `x()` or `x?.()` has no side effects, the whole tree is useless now.
			NodeUtil_RemoveChild(parent, n)
			p.ReportChangeToEnclosingScope(parent)
			return nil
		}
	}

	// Try transforms that apply to both IF and HOOK.
	condValue := NodeUtil_GetBooleanValue(cond)
	if condValue == base.Tri_UNKNOWN {
		return n // We can't remove branches otherwise!
	}

	if p.MayHaveSideEffects(cond) {
		// Transform "if (a = 2) {x =2}" into "if (true) {a=2;x=2}"
		newConditionValue := condValue == base.Tri_TRUE
		// Add an elseBody if it is needed.
		if !newConditionValue && elseBody == nil {
			elseBody = rhino.IR_Block().Srcref(n)
			n.AddChildToBack(elseBody)
		}
		newCond := NodeUtil_BooleanNode(newConditionValue)
		cond.ReplaceWith(newCond)
		var branchToKeep *rhino.Node
		if newConditionValue {
			branchToKeep = thenBody
		} else {
			branchToKeep = elseBody
		}
		branchToKeep.AddChildToFront(rhino.IR_ExprResult(cond).Srcref(cond))
		p.ReportChangeToEnclosingScope(branchToKeep)
		cond = newCond
	}

	condTrue := condValue.ToBoolean(true)
	if n.HasTwoChildren() {
		jutil.CheckState(type_ == rhino.Token_IF)

		if condTrue {
			// Replace "if (true) { X }" with "X".
			thenStmt := n.GetSecondChild()
			thenStmt.Detach()
			n.ReplaceWith(thenStmt)
			p.ReportChangeToEnclosingScope(thenStmt)
			return thenStmt
		} else {
			// Remove "if (false) { X }" completely.
			NodeUtil_RedeclareVarsInsideBranch(n)
			NodeUtil_RemoveChild(parent, n)
			p.ReportChangeToEnclosingScope(parent)
			p.MarkFunctionsDeleted(n)
			return nil
		}
	} else {
		// Replace "if (true) { X } else { Y }" with X, or
		// replace "if (false) { X } else { Y }" with Y.
		trueBranch := n.GetSecondChild()
		falseBranch := trueBranch.GetNext()
		var branchToKeep, branchToRemove *rhino.Node
		if condTrue {
			branchToKeep = trueBranch
			branchToRemove = falseBranch
		} else {
			branchToKeep = falseBranch
			branchToRemove = trueBranch
		}
		NodeUtil_RedeclareVarsInsideBranch(branchToRemove)
		branchToKeep.Detach()
		n.ReplaceWith(branchToKeep)
		p.ReportChangeToEnclosingScope(branchToKeep)
		p.MarkFunctionsDeleted(n)
		return branchToKeep
	}
}

// TryFoldHook tries folding HOOK (?:) if the condition results of the
// condition is known.
//
// Returns the replacement node, if changed, or the original if not.
//
// Java: PeepholeRemoveDeadCode.tryFoldHook(Node)
func (p *PeepholeRemoveDeadCode) TryFoldHook(n *rhino.Node) *rhino.Node {
	jutil.CheckStateMsg1(n.IsHook(), "%s", n)
	parent := n.GetParent()
	jutil.CheckNotNull(parent)
	cond := n.GetFirstChild()
	thenBody := cond.GetNext()
	elseBody := thenBody.GetNext()

	condValue := NodeUtil_GetBooleanValue(cond)
	if condValue == base.Tri_UNKNOWN {
		// If the result nodes are equivalent, then one of the nodes can be
		// removed and it doesn't matter which.
		if !p.AreNodesEqualForInlining(thenBody, elseBody) {
			return n // We can't remove branches otherwise!
		}
	}

	// Transform "(a = 2) ? x =2 : y" into "a=2,x=2"
	var branchToKeep *rhino.Node
	var branchToRemove *rhino.Node
	if condValue.ToBoolean(true) {
		branchToKeep = thenBody
		branchToRemove = elseBody
	} else {
		branchToKeep = elseBody
		branchToRemove = thenBody
	}

	var replacement *rhino.Node
	condHasSideEffects := p.MayHaveSideEffects(cond)
	// Must detach after checking for side effects, to ensure that the parents
	// of nodes are set correctly.
	n.DetachChildren()
	if condHasSideEffects {
		replacement = rhino.IR_Comma(cond, branchToKeep).Srcref(n)
	} else {
		replacement = branchToKeep
		p.MarkFunctionsDeleted(cond)
	}

	n.ReplaceWith(replacement)
	p.ReportChangeToEnclosingScope(replacement)
	p.MarkFunctionsDeleted(branchToRemove)
	return replacement
}

// TryFoldFor removes FORs that always evaluate to false.
//
// Java: PeepholeRemoveDeadCode.tryFoldFor(Node)
func (p *PeepholeRemoveDeadCode) TryFoldFor(n *rhino.Node) *rhino.Node {
	jutil.CheckArgument(n.IsVanillaFor())

	init := n.GetFirstChild()
	cond := init.GetNext()
	increment := cond.GetNext()

	if !init.IsEmpty() && !NodeUtil_IsNameDeclaration(init) {
		init = p.TrySimplifyUnusedResult(init)
		if init == nil {
			init = rhino.IR_Empty().Srcref(n)
			n.AddChildToFront(init)
		}
	}

	if !increment.IsEmpty() {
		increment = p.TrySimplifyUnusedResult(increment)
		if increment == nil {
			increment = rhino.IR_Empty().Srcref(n)
			increment.InsertAfter(cond)
		}
	}

	// There is an initializer skip it
	if !n.GetFirstChild().IsEmpty() {
		return n
	}

	if NodeUtil_GetBooleanValue(cond) != base.Tri_FALSE {
		return n
	}

	parent := n.GetParent()
	NodeUtil_RedeclareVarsInsideBranch(n)

	if !p.MayHaveSideEffects(cond) {
		// Remove the entire loop and any associated labels.
		for parent.IsLabel() {
			n = parent
			parent = parent.GetParent()
		}
		n.Detach()
	} else {
		statement := rhino.IR_ExprResult(cond.Detach()).SrcrefIfMissing(cond)
		if parent.IsLabel() {
			block := rhino.IR_Block()
			block.SrcrefIfMissing(statement)
			block.AddChildToFront(statement)
			statement = block
		}
		n.ReplaceWith(statement)
	}
	p.ReportChangeToEnclosingScope(parent)

	return nil
}

// TryFoldDoAway removes DOs that always evaluate to false. This leaves the
// statements that were in the loop in a BLOCK node. The block will be removed
// in a later pass, if possible.
//
// Java: PeepholeRemoveDeadCode.tryFoldDoAway(Node)
func (p *PeepholeRemoveDeadCode) TryFoldDoAway(n *rhino.Node) *rhino.Node {
	jutil.CheckArgument(n.IsDo())

	cond := NodeUtil_GetConditionExpression(n)
	if NodeUtil_GetBooleanValue(cond) != base.Tri_FALSE {
		return n
	}

	block := NodeUtil_GetLoopCodeBlock(n)
	if n.GetParent().IsLabel() || PeepholeRemoveDeadCode_HasUnnamedBreakOrContinue(block) {
		return n
	}

	parent := n.GetParent()
	n.ReplaceWith(block.Detach())
	if p.MayHaveSideEffects(cond) {
		condStatement := rhino.IR_ExprResult(cond.Detach()).Srcref(cond)
		condStatement.InsertAfter(block)
	}
	p.ReportChangeToEnclosingScope(parent)

	return block
}

// TryFoldEmptyDo removes DOs that have empty bodies into FORs, which are
// much easier for the CFA to analyze.
//
// Java: PeepholeRemoveDeadCode.tryFoldEmptyDo(Node)
func (p *PeepholeRemoveDeadCode) TryFoldEmptyDo(n *rhino.Node) *rhino.Node {
	jutil.CheckArgument(n.IsDo())

	body := NodeUtil_GetLoopCodeBlock(n)
	if body.IsBlock() && !body.HasChildren() {
		cond := NodeUtil_GetConditionExpression(n)
		forNode :=
			rhino.IR_ForNode(rhino.IR_Empty().Srcref(n), cond.Detach(), rhino.IR_Empty().Srcref(n), body.Detach())
		n.ReplaceWith(forNode)
		p.ReportChangeToEnclosingScope(forNode)
		return forNode
	}
	return n
}

// TryOptimizeObjectPattern removes string keys with an empty pattern as
// their child.
//
// Java: PeepholeRemoveDeadCode.tryOptimizeObjectPattern(Node)
func (p *PeepholeRemoveDeadCode) TryOptimizeObjectPattern(pattern *rhino.Node) *rhino.Node {
	jutil.CheckArgumentMsg1(pattern.IsObjectPattern(), "%s", pattern)

	if pattern.HasChildren() && pattern.GetLastChild().IsRest() {
		// don't remove any elements in `const {f: [], ...rest} = obj` because that affects what's
		// assigned to `rest`. only the last element can be object rest.
		return pattern
	}

	// remove trailing EMPTY nodes and empty destructuring patterns
	for child := pattern.GetFirstChild(); child != nil; {
		key := child
		child = key.GetNext() // don't put this in the for loop since we might remove `child`

		if !key.IsStringKey() {
			// don't try to remove rest or computed properties, since they might have side effects
			continue
		}
		if p.IsRemovableDestructuringTarget(key.GetOnlyChild()) {
			// e.g. `const {f: {}} = obj;`
			key.Detach()
			p.ReportChangeToEnclosingScope(pattern)
		}
	}
	return pattern
}

// TryOptimizeArrayPattern removes trailing EMPTY nodes and empty array
// patterns.
//
// Java: PeepholeRemoveDeadCode.tryOptimizeArrayPattern(Node)
func (p *PeepholeRemoveDeadCode) TryOptimizeArrayPattern(pattern *rhino.Node) *rhino.Node {
	jutil.CheckArgumentMsg1(pattern.IsArrayPattern(), "%s", pattern)

	for lastChild := pattern.GetLastChild(); lastChild != nil; {
		if lastChild.IsEmpty() || p.IsRemovableDestructuringTarget(lastChild) {
			prev := lastChild.GetPrevious()
			lastChild.Detach()
			lastChild = prev
			p.ReportChangeToEnclosingScope(pattern)
		} else {
			// don't remove any non-trailing empty nodes because that will change the ordering of the
			// other assignments
			// note that this case also covers array pattern rest, which must be the final element
			break
		}
	}
	return pattern
}

// Java: PeepholeRemoveDeadCode.isRemovableDestructuringTarget(Node)
func (p *PeepholeRemoveDeadCode) IsRemovableDestructuringTarget(destructruringElement *rhino.Node) bool {
	target := destructruringElement
	var defaultValue *rhino.Node
	if destructruringElement.IsDefaultValue() {
		target = destructruringElement.GetFirstChild()
		defaultValue = destructruringElement.GetSecondChild()
	}
	if !target.IsDestructuringPattern() || target.HasChildren() {
		return false
	}
	// only remove default values without side effects
	return defaultValue == nil || !p.MayHaveSideEffects(defaultValue)
}

// PeepholeRemoveDeadCode_HasUnnamedBreakOrContinue returns whether a node
// has any unhandled breaks or continue.
//
// Java: PeepholeRemoveDeadCode.hasUnnamedBreakOrContinue(Node)
func PeepholeRemoveDeadCode_HasUnnamedBreakOrContinue(n *rhino.Node) bool {
	return NodeUtil_Has(
		n,
		// Check for unlabeled breaks
		func(node *rhino.Node) bool { return node.IsBreak() && !node.HasChildren() },
		// ...inside contexts that can contain breaks.
		func(node *rhino.Node) bool {
			return !rhino.IR_MayBeExpression(node) && // Functions are not visited
				!NodeUtil_IsLoopStructure(node) &&
				!node.IsSwitch()
		}) ||
		NodeUtil_Has(
			n,
			// Check for unlabeled continues
			func(node *rhino.Node) bool { return node.IsContinue() && !node.HasChildren() },
			// ...inside contexts that can contain continues.
			func(node *rhino.Node) bool {
				return !rhino.IR_MayBeExpression(node) && // Functions are not visited
					!NodeUtil_IsLoopStructure(node)
			})
}

// TryFoldForCondition removes always true loop conditions.
//
// Java: PeepholeRemoveDeadCode.tryFoldForCondition(Node)
func (p *PeepholeRemoveDeadCode) TryFoldForCondition(forCondition *rhino.Node) {
	if p.GetSideEffectFreeBooleanValue(forCondition) == base.Tri_TRUE {
		p.ReportChangeToEnclosingScope(forCondition)
		forCondition.ReplaceWith(rhino.IR_Empty())
	}
}

// Java: PeepholeRemoveDeadCode.tryRemoveOptionalChaining(Node)
func (p *PeepholeRemoveDeadCode) TryRemoveOptionalChaining(optionalChain *rhino.Node) *rhino.Node {
	callee := optionalChain.GetFirstChild()
	if !NodeUtil_IsNullOrUndefined(callee) {
		return optionalChain
	}
	var result *rhino.Node
	if p.MayHaveSideEffects(callee) {
		// Simplify `(void sideEffectFunction())?.()` to `(void sideEffectFunction())`
		// The optional chain call won't execute but sideEffectFunction() is still evaluated.
		optionalChain.ReplaceWith(callee.Detach())
		result = callee
	} else {
		// Remove `(void 0)?.()` and (null)?.() and simplify `(void 0)?.x and null?.x` to void 0
		result = NodeUtil_NewUndefinedNode(callee)
		optionalChain.ReplaceWith(result)
	}
	p.MarkFunctionsDeleted(optionalChain)
	p.ReportChangeToEnclosingScope(result)
	return result
}
