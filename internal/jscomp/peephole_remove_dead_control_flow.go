// Port of com/google/javascript/jscomp/PeepholeRemoveDeadCode.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Java: PeepholeRemoveDeadCode.tryFoldLabel(Node)
func (p *PeepholeRemoveDeadCode) TryFoldLabel(n rhino.Node) rhino.Node {
	label, ok := ast.AsLabel(n)
	jutil.CheckState(ok)
	stmt := label.Body()
	if _, ok := stmt.AsEmpty(); ok {
		p.ReportChangeToEnclosingScope(n)
		n.Detach()
		return rhino.NilNode
	}

	if block, ok := stmt.AsBlock(); ok && !block.HasStatements() {
		p.ReportChangeToEnclosingScope(n)
		if n.GetParent().IsLabel() {
			// If the parent is itself a label, replace this label
			// with its contained block to keep the AST in a valid state.
			n.ReplaceWith(block.Node().Detach())
		} else {
			n.Detach()
		}
		return rhino.NilNode
	}

	if child, ok := PeepholeRemoveDeadCode_GetOnlyInterestingChild(stmt); ok {
		stmt = child
	}
	// Go-only: Java reads the label of the break without checking that it has
	// one, and throws on an unlabeled break such as `a: break;` in a switch
	// case. Label() makes the missing label explicit: such a label is kept.
	brk, ok := stmt.AsBreak()
	if !ok {
		return n
	}
	if target, ok := brk.Label(); ok && target.Node().GetStringID() == label.Name().Node().GetStringID() {
		if n.GetParent().IsLabel() {
			replacement := rhino.IR_Block().Srcref(n)
			n.ReplaceWith(replacement)
			p.ReportChangeToEnclosingScope(replacement)
			return replacement
		}
		parent := n.GetParent()
		n.Detach()
		p.ReportChangeToEnclosingScope(parent)
		return rhino.NilNode
	}
	return n
}

// PeepholeRemoveDeadCode_GetOnlyInterestingChild returns the only
// "interesting" child of block, if it has exactly one interesting child;
// ok is false otherwise (or if stmt is not a block). For purposes of this
// method, a node is considered "interesting" unless it is an empty
// synthetic block.
//
// Java: PeepholeRemoveDeadCode.getOnlyInterestingChild(Node)
func PeepholeRemoveDeadCode_GetOnlyInterestingChild(stmt ast.Statement) (ast.Statement, bool) {
	block, ok := stmt.AsBlock()
	if !ok {
		return ast.Statement{}, false
	}
	if block.Node().HasOneChild() {
		return block.FirstStatement()
	}

	var ret ast.Statement
	found := false
	for child := range block.Statements() {
		if child.Node().IsSyntheticBlock() && !child.Node().HasChildren() {
			// Uninteresting child.
		} else if found {
			// Found more than one interesting child.
			return ast.Statement{}, false
		} else {
			ret, found = child, true
		}
	}
	return ret, found
}

// TryFoldTry removes try blocks without catch blocks and with empty or not
// existent finally blocks. Or, only leave the finally blocks if try body
// blocks are empty.
//
// Returns the replacement node, if changed, or the original if not.
//
// Java: PeepholeRemoveDeadCode.tryFoldTry(Node)
func (p *PeepholeRemoveDeadCode) TryFoldTry(n rhino.Node) rhino.Node {
	try, ok := ast.AsTry(n)
	jutil.CheckStateMsg1(ok, "%s", n)
	body := try.Body()
	catch, hasCatch := try.Catch()
	finally, hasFinally := try.Finally()

	// Removes TRYs that had its CATCH removed and/or empty FINALLY.
	if !hasCatch && (!hasFinally || !finally.HasStatements()) {
		jutil.CheckState(!n.GetParent().IsLabel())
		body.Node().Detach()
		n.ReplaceWith(body.Node())
		p.ReportChangeToEnclosingScope(body.Node())
		return body.Node()
	}

	// Only leave FINALLYs if TRYs are not empty
	if !body.HasStatements() {
		if hasCatch {
			NodeUtil_RedeclareVarsInsideBranch(catch.Node())
		}
		p.ReportChangeToEnclosingScope(n)
		if hasFinally {
			finally.Node().Detach()
			jutil.CheckState(!n.GetParent().IsLabel())
			n.ReplaceWith(finally.Node())
			return finally.Node()
		}
		jutil.CheckState(!n.GetParent().IsLabel())
		n.Detach()
		return rhino.NilNode
	}

	return n
}

// peepholeRemoveDeadCode_isUnnamedBreak reports whether n is a break without
// a label.
//
// Go-only: Java writes n.isBreak() && !n.hasChildren() at each use.
func peepholeRemoveDeadCode_isUnnamedBreak(n rhino.Node) bool {
	brk, ok := ast.AsBreak(n)
	if !ok {
		return false
	}
	_, hasLabel := brk.Label()
	return !hasLabel
}

// peepholeRemoveDeadCode_isUnnamedContinue reports whether n is a continue
// without a label.
//
// Go-only: Java writes n.isContinue() && !n.hasChildren().
func peepholeRemoveDeadCode_isUnnamedContinue(n rhino.Node) bool {
	cont, ok := ast.AsContinue(n)
	if !ok {
		return false
	}
	_, hasLabel := cont.Label()
	return !hasLabel
}

func peepholeRemoveDeadCode_hasUnnamedBreakInScope(n rhino.Node) bool {
	if peepholeRemoveDeadCode_isUnnamedBreak(n) {
		return true
	}
	if n.IsFunction() {
		return false
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if peepholeRemoveDeadCode_hasUnnamedBreakInScope(c) {
			return true
		}
	}
	return false
}

// RemoveIfUnnamedBreak removes maybeBreak if it is a break without a label.
//
// Go-only: Java also accepts null; callers here only pass a statement.
//
// Java: PeepholeRemoveDeadCode.removeIfUnnamedBreak(Node)
func (p *PeepholeRemoveDeadCode) RemoveIfUnnamedBreak(maybeBreak ast.Statement) {
	if peepholeRemoveDeadCode_isUnnamedBreak(maybeBreak.Node()) {
		p.ReportChangeToEnclosingScope(maybeBreak.Node())
		maybeBreak.Node().Detach()
	}
}

// Java: PeepholeRemoveDeadCode.tryRemoveSwitchWithSingleCase(Node, boolean)
func (p *PeepholeRemoveDeadCode) TryRemoveSwitchWithSingleCase(sw ast.Switch, shouldHoistCondition bool) rhino.Node {
	n := sw.Node()
	cases := sw.Body()
	jutil.CheckState(cases.Node().HasOneChild())
	onlyCase, _ := cases.FirstCase() // the check above proves there is one
	caseBlock := onlyCase.Body()
	if last, ok := caseBlock.LastStatement(); ok {
		p.RemoveIfUnnamedBreak(last)
	}
	// Back off if the switch contains statements like "if (a) { break; }"
	if peepholeRemoveDeadCode_hasUnnamedBreakInScope(caseBlock.Node()) {
		return n
	}
	if shouldHoistCondition {
		cond := sw.Discriminant().Node()
		rhino.IR_ExprResult(cond.Detach()).Srcref(n).InsertBefore(n)
	}
	n.ReplaceWith(caseBlock.Node().Detach())
	p.ReportChangeToEnclosingScope(caseBlock.Node())
	return caseBlock.Node()
}

// Java: PeepholeRemoveDeadCode.tryRemoveSwitch(Node)
func (p *PeepholeRemoveDeadCode) TryRemoveSwitch(sw ast.Switch) rhino.Node {
	n := sw.Node()
	cases := sw.Body()
	if !cases.HasCases() {
		// Remove the switch if there are no remaining cases
		condition := sw.Discriminant().Node()
		replacement := rhino.IR_ExprResult(condition.Detach()).Srcref(n)
		n.ReplaceWith(replacement)
		p.ReportChangeToEnclosingScope(replacement)
		return replacement
	}
	if cases.Node().HasOneChild() {
		onlyCase, _ := cases.FirstCase() // the check above proves there is one
		if _, ok := onlyCase.AsDefaultCase(); ok {
			// Before removing switch, we must preserve the switch condition if it has side effects
			shouldHoistCondition := p.MayHaveSideEffects(sw.Discriminant().Node())
			return p.TryRemoveSwitchWithSingleCase(sw, shouldHoistCondition)
		}
	}
	return n
}

// peepholeRemoveDeadCode_caseTagOrBody returns the tag of a case, or the
// body of a default case.
//
// Go-only: Java reads the first child of each case, which for a default case
// (it has no tag) is its body: the switch optimizations below then compare
// that body with the condition and check it for side effects. This helper
// keeps that behavior, for the same output, and makes it visible.
func peepholeRemoveDeadCode_caseTagOrBody(c ast.SwitchCase) rhino.Node {
	if cs, ok := c.AsCase(); ok {
		return cs.Value().Node()
	}
	return c.Body().Node()
}

// TryOptimizeSwitch removes useless switches and cases.
//
// Java: PeepholeRemoveDeadCode.tryOptimizeSwitch(Node)
func (p *PeepholeRemoveDeadCode) TryOptimizeSwitch(n rhino.Node) rhino.Node {
	sw, ok := ast.AsSwitch(n)
	jutil.CheckStateMsg1(ok, "%s", n)

	cases := sw.Body()
	defaultCase := p.TryOptimizeDefaultCase(sw)

	// Generally, it is unsafe to remove other cases when the default case is not the last one.
	// (defaultCase, when there is one, is the only default case.)
	if (defaultCase.IsNil() || defaultCase.GetNext().IsNil()) && p.AreAllCaseTagsLiterals(cases) {
		cond := sw.Discriminant().Node()

		// First, remove empty cases where possible: always empty default cases; or when there is
		// no default case, other empty cases that are not the first matching case, may be
		// removable as well.
		var prev rhino.Node // the last case kept, if any
		foundMatchingCase := false
		for cur, ok := cases.FirstCase(); ok; {
			next, hasNext := cur.Next() // read it first: cur may be removed
			tag := peepholeRemoveDeadCode_caseTagOrBody(cur)
			foundMatchingCase = p.IsFirstSwitchMatch(foundMatchingCase, cond, tag)
			if !foundMatchingCase &&
				!p.MayHaveSideEffects(tag) &&
				p.IsUselessCase(cur, prev, defaultCase) {
				p.RemoveCase(n, cur.Node())
			} else {
				prev = cur.Node()
			}
			cur, ok = next, hasNext
		}

		// Next, optimize switches with constant condition
		if NodeUtil_IsLiteralValue(cond, false) {
			caseMatches := base.Tri_TRUE
			// Remove cases until you find one that may match
			cur, ok := cases.FirstCase()
			for ok {
				caseMatches = PeepholeFoldConstants_EvaluateComparison(
					&p.AbstractPeepholeOptimizationBase, rhino.Token_SHEQ, cond, peepholeRemoveDeadCode_caseTagOrBody(cur))
				if caseMatches == base.Tri_TRUE || caseMatches == base.Tri_UNKNOWN {
					break
				}
				next, hasNext := cur.Next()
				p.RemoveCase(n, cur.Node())
				cur, ok = next, hasNext
			}
			if ok && caseMatches == base.Tri_TRUE {
				// Skip cases until you find one whose last stm is a removable break
				matchingCase := cur
				matchingCaseBlock := matchingCase.Body()
				for ok {
					block := cur.Body()
					isLastStmRemovableBreak := false
					if lastStm, ok := block.LastStatement(); ok && PeepholeRemoveDeadCode_IsExit(lastStm.Node()) {
						p.RemoveIfUnnamedBreak(lastStm)
						isLastStmRemovableBreak = true
					}
					next, hasNext := cur.Next()
					// Remove the fallthrough case labels
					if cur != matchingCase {
						for block.HasStatements() {
							matchingCaseBlock.Node().AddChildToBack(block.Node().RemoveFirstChild())
						}
						p.ReportChangeToEnclosingScope(cur.Node())
						cur.Node().Detach()
					}
					cur, ok = next, hasNext
					if isLastStmRemovableBreak {
						break
					}
				}

				// Remove any remaining cases
				for ok {
					next, hasNext := cur.Next()
					p.RemoveCase(n, cur.Node())
					cur, ok = next, hasNext
				}
				// If there is one case left, we may be able to fold it
				// Go-only: Java tests cond.getNext(), which is the switch body (not
				// a case), so its test always passes; only the matching case is
				// left at this point anyway.
				if cases.Node().HasOneChild() {
					return p.TryRemoveSwitchWithSingleCase(sw, false)
				}
			}
		}
	}

	// Last, try to remove the entire switch if possible
	return p.TryRemoveSwitch(sw)
}

// TryOptimizeDefaultCase returns the default case node or null if there is
// no default case or if the default case is removed.
//
// Java: PeepholeRemoveDeadCode.tryOptimizeDefaultCase(Node)
func (p *PeepholeRemoveDeadCode) TryOptimizeDefaultCase(sw ast.Switch) rhino.Node {
	n := sw.Node()
	cases := sw.Body()
	var lastNonRemovable rhino.Node // the most recently iterated case known to not be removable.

	for c := range cases.Cases() {
		if cs, ok := c.AsCase(); ok {
			if cs.Body().HasStatements() || p.MayHaveSideEffects(cs.Value().Node()) {
				lastNonRemovable = c.Node()
			}
			continue
		}

		// c is the default case.
		// Remove any cases that fall-through to the default case
		caseToRemove := cases.Node().GetFirstChild()
		if !lastNonRemovable.IsNil() {
			caseToRemove = lastNonRemovable.GetNext()
		}
		for next := rhino.NilNode; !c.Is(caseToRemove); caseToRemove = next {
			next = caseToRemove.GetNext()
			p.RemoveCase(n, caseToRemove)
		}

		// Remove the default case if we can
		if p.IsUselessCase(c, lastNonRemovable, c.Node()) {
			p.RemoveCase(n, c.Node())
			return rhino.NilNode
		}
		return c.Node()
	}
	return rhino.NilNode
}

// RemoveCase removes the case from the switch redeclaring any variables
// declared in it.
//
// caseNode is the case to remove.
//
// Java: PeepholeRemoveDeadCode.removeCase(Node, Node)
func (p *PeepholeRemoveDeadCode) RemoveCase(switchNode rhino.Node, caseNode rhino.Node) {
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
// previousCase (the case before caseNode) and defaultCase are the nil node
// when there is none.
//
// Java: PeepholeRemoveDeadCode.isUselessCase(Node, Node, Node)
func (p *PeepholeRemoveDeadCode) IsUselessCase(
	caseNode ast.SwitchCase, previousCase rhino.Node, defaultCase rhino.Node) bool {
	jutil.CheckState(previousCase.IsNil() || previousCase.GetNext() == caseNode.Node())
	// A case isn't useless if a previous case falls through to it unless it happens to be the last
	// case in the switch.
	if prev, ok := ast.AsSwitchCase(previousCase); ok {
		if _, hasNext := caseNode.Next(); hasNext {
			previousLast, ok := prev.Body().LastStatement()
			if !ok || !PeepholeRemoveDeadCode_IsExit(previousLast.Node()) {
				return false
			}
		}
	}

	for executingCase, ok := caseNode, true; ok; executingCase, ok = executingCase.Next() {
		// We only expect a DEFAULT case if the case we are checking is the
		// DEFAULT case.  Otherwise, we assume the DEFAULT case has already
		// been removed.
		c, isCase := executingCase.AsCase()
		jutil.CheckState(caseNode == executingCase || isCase)
		if isCase && p.MayHaveSideEffects(c.Value().Node()) {
			// The case falls thru to a case whose condition has a potential side-effect,
			// removing the candidate case would skip that side-effect, so don't.
			return false
		}
		for blockChild := range executingCase.Body().Statements() {
			// If this is a block with a labelless break, it is useless.
			switch blockChild.Node().GetToken() {
			case rhino.Token_BREAK:
				// A case with a single labelless break is useless if it is the default case or if
				// there is no default case. A break to a different control structure isn't useless.
				return peepholeRemoveDeadCode_isUnnamedBreak(blockChild.Node()) &&
					(defaultCase.IsNil() || executingCase.Is(defaultCase))
			case rhino.Token_VAR:
				// Variable declarations without initializations are OK: one declarator,
				// and no child under it (a destructuring declarator has its pattern).
				if blockChild.Node().HasOneChild() && blockChild.Node().GetFirstFirstChild().IsNil() {
					continue
				}
				return false
			default:
				return false
			}
		}
		// Look at the fallthrough case
	}
	return true
}

// Java: PeepholeRemoveDeadCode.isFirstSwitchMatch(boolean, Node, Node)
func (p *PeepholeRemoveDeadCode) IsFirstSwitchMatch(foundMatchingCase bool, condition rhino.Node, tag rhino.Node) bool {
	if foundMatchingCase {
		return false
	}
	return PeepholeFoldConstants_EvaluateComparison(&p.AbstractPeepholeOptimizationBase, rhino.Token_SHEQ, condition, tag) == base.Tri_TRUE
}

// Java: PeepholeRemoveDeadCode.areAllCaseTagsLiterals(Node)
func (p *PeepholeRemoveDeadCode) AreAllCaseTagsLiterals(cases ast.SwitchBody) bool {
	for c := range cases.Cases() {
		// A default case has no tag.
		if cs, ok := c.AsCase(); ok && !NodeUtil_IsLiteralValue(cs.Value().Node(), false) {
			return false
		}
	}
	return true
}

// PeepholeRemoveDeadCode_IsExit returns whether the node is a control flow
// exit from the current block.
//
// Java: PeepholeRemoveDeadCode.isExit(Node)
func PeepholeRemoveDeadCode_IsExit(n rhino.Node) bool {
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
func PeepholeRemoveDeadCode_IsUnconditionalBlockExit(block ast.Block) bool {
	jutil.CheckStateMsg1(!block.Node().GetParent().IsLabel(), "%s", block.Node())

	// Last statement must lead out of the block.
	lastStm, ok := block.LastStatement()
	if !ok {
		return false
	}
	switch lastStm.Node().GetToken() {
	case rhino.Token_BREAK:
		if peepholeRemoveDeadCode_isUnnamedBreak(lastStm.Node()) {
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
	for child := range block.Statements() {
		if child == lastStm {
			break
		}
		switch child.Node().GetToken() {
		case rhino.Token_BREAK:
			if peepholeRemoveDeadCode_isUnnamedBreak(child.Node()) {
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
func PeepholeRemoveDeadCode_IsSwitchExit(n rhino.Node) bool {
	sw, ok := ast.AsSwitch(n)
	jutil.CheckStateMsg1(ok, "%s", n)

	hasDefaultCase := false

	for switchCase := range sw.Body().Cases() {
		if _, ok := switchCase.AsDefaultCase(); ok {
			hasDefaultCase = true
		}
		block := switchCase.Body()
		if !block.HasStatements() {
			if _, hasNext := switchCase.Next(); hasNext {
				continue // An empty case falls through to the next one.
			}
		}
		if !PeepholeRemoveDeadCode_IsUnconditionalBlockExit(block) {
			return false
		}
	}

	return hasDefaultCase
}

// PeepholeRemoveDeadCode_IsTryExit returns true if the try..catch always
// "exits" (return, throw, etc).
//
// Java: PeepholeRemoveDeadCode.isTryExit(Node)
func PeepholeRemoveDeadCode_IsTryExit(n rhino.Node) bool {
	try, ok := ast.AsTry(n)
	jutil.CheckStateMsg1(ok, "%s", n)

	// finally - regardless of the behavior of the other blocks,
	// an exit from the finally with guarantee that behavior.
	if finally, ok := try.Finally(); ok && PeepholeRemoveDeadCode_IsUnconditionalBlockExit(finally) {
		return true
	}
	// try
	if !PeepholeRemoveDeadCode_IsUnconditionalBlockExit(try.Body()) {
		return false
	}
	// catch
	catch, ok := try.Catch()
	return !ok || PeepholeRemoveDeadCode_IsUnconditionalBlockExit(catch.Body())
}

// TryFoldIf tries folding IF nodes by removing dead branches.
//
// Returns the replacement node, if changed, or the original if not.
//
// Java: PeepholeRemoveDeadCode.tryFoldIf(Node)
func (p *PeepholeRemoveDeadCode) TryFoldIf(n rhino.Node) rhino.Node {
	ifs, ok := ast.AsIf(n)
	jutil.CheckStateMsg1(ok, "%s", n)
	parent := n.GetParent()
	rhino.CheckNotNull(parent)
	cond := ifs.Cond().Node()
	thenBody := ifs.Then()
	elseBody, hasElse := ifs.Else()

	// if (x) { .. } else { } --> if (x) { ... }
	if hasElse && !p.MayHaveSideEffects(elseBody.Node()) {
		ifs.RemoveElse()
		p.ReportChangeToEnclosingScope(n)
		hasElse = false
	}

	// if (x) { } else { ... } --> if (!x) { ... }
	if !p.MayHaveSideEffects(thenBody.Node()) && hasElse {
		ifs.RemoveElse()
		ifs.SetThen(elseBody)
		notCond := rhino.NewNode(rhino.Token_NOT)
		cond.ReplaceWith(notCond)
		p.ReportChangeToEnclosingScope(n)
		notCond.AddChildToFront(cond)
		cond = notCond
		thenBody = elseBody
		hasElse = false
	}

	// `if (x()) { }` or `if (x?.()) { }`
	if !p.MayHaveSideEffects(thenBody.Node()) && !hasElse {
		if p.MayHaveSideEffects(cond) {
			// `x()` or `x?.()` has side effects, just leave the condition on its own.
			cond.Detach()
			replacement := NodeUtil_NewExpr(cond)
			n.ReplaceWith(replacement)
			p.ReportChangeToEnclosingScope(parent)
			return replacement
		}
		// `x()` or `x?.()` has no side effects, the whole tree is useless now.
		NodeUtil_RemoveChild(parent, n)
		p.ReportChangeToEnclosingScope(parent)
		return rhino.NilNode
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
		if !newConditionValue && !hasElse {
			elseBody = ast.NewBlock()
			elseBody.Node().Srcref(n)
			ifs.SetElse(elseBody)
			hasElse = true
		}
		cond.ReplaceWith(NodeUtil_BooleanNode(newConditionValue))
		branchToKeep := thenBody
		if !newConditionValue {
			branchToKeep = elseBody
		}
		branchToKeep.Node().AddChildToFront(rhino.IR_ExprResult(cond).Srcref(cond))
		p.ReportChangeToEnclosingScope(branchToKeep.Node())
	}

	condTrue := condValue.ToBoolean(true)
	if !hasElse {
		if condTrue {
			// Replace "if (true) { X }" with "X".
			thenStmt := thenBody.Node()
			thenStmt.Detach()
			n.ReplaceWith(thenStmt)
			p.ReportChangeToEnclosingScope(thenStmt)
			return thenStmt
		}
		// Remove "if (false) { X }" completely.
		NodeUtil_RedeclareVarsInsideBranch(n)
		NodeUtil_RemoveChild(parent, n)
		p.ReportChangeToEnclosingScope(parent)
		p.MarkFunctionsDeleted(n)
		return rhino.NilNode
	}

	// Replace "if (true) { X } else { Y }" with X, or
	// replace "if (false) { X } else { Y }" with Y.
	branchToKeep, branchToRemove := thenBody, elseBody
	if !condTrue {
		branchToKeep, branchToRemove = elseBody, thenBody
	}
	NodeUtil_RedeclareVarsInsideBranch(branchToRemove.Node())
	branchToKeep.Node().Detach()
	n.ReplaceWith(branchToKeep.Node())
	p.ReportChangeToEnclosingScope(branchToKeep.Node())
	p.MarkFunctionsDeleted(n)
	return branchToKeep.Node()
}

// TryFoldHook tries folding HOOK (?:) if the condition results of the
// condition is known.
//
// Returns the replacement node, if changed, or the original if not.
//
// Java: PeepholeRemoveDeadCode.tryFoldHook(Node)
func (p *PeepholeRemoveDeadCode) TryFoldHook(n rhino.Node) rhino.Node {
	hook, ok := ast.AsHook(n)
	jutil.CheckStateMsg1(ok, "%s", n)
	parent := n.GetParent()
	rhino.CheckNotNull(parent)
	cond := hook.Cond().Node()
	thenBody := hook.Then().Node()
	elseBody := hook.Else().Node()

	condValue := NodeUtil_GetBooleanValue(cond)
	if condValue == base.Tri_UNKNOWN {
		// If the result nodes are equivalent, then one of the nodes can be
		// removed and it doesn't matter which.
		if !p.AreNodesEqualForInlining(thenBody, elseBody) {
			return n // We can't remove branches otherwise!
		}
	}

	// Transform "(a = 2) ? x =2 : y" into "a=2,x=2"
	var branchToKeep rhino.Node
	var branchToRemove rhino.Node
	if condValue.ToBoolean(true) {
		branchToKeep = thenBody
		branchToRemove = elseBody
	} else {
		branchToKeep = elseBody
		branchToRemove = thenBody
	}

	var replacement rhino.Node
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
func (p *PeepholeRemoveDeadCode) TryFoldFor(n rhino.Node) rhino.Node {
	forNode, ok := ast.AsFor(n)
	jutil.CheckArgument(ok)

	// Read the parts first: removing the init or the update below shifts the
	// children that follow it.
	init, hasInit := forNode.Init()
	update, hasUpdate := forNode.Update()
	body := forNode.Body()

	if hasInit {
		if _, isDecl := init.AsNameDeclaration(); !isDecl && p.TrySimplifyUnusedResult(init.Node()).IsNil() {
			// The init was removed: put an EMPTY in its place.
			n.AddChildToFront(rhino.IR_Empty().Srcref(n))
		}
	}

	if hasUpdate && p.TrySimplifyUnusedResult(update.Node()).IsNil() {
		// The update was removed: put an EMPTY in its place, before the body.
		rhino.IR_Empty().Srcref(n).InsertBefore(body.Node())
	}

	// There is an initializer skip it
	if _, ok := forNode.Init(); ok {
		return n
	}

	condExpr, ok := forNode.Cond()
	if !ok || NodeUtil_GetBooleanValue(condExpr.Node()) != base.Tri_FALSE {
		return n
	}
	cond := condExpr.Node()

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

	return rhino.NilNode
}

// TryFoldDoAway removes DOs that always evaluate to false. This leaves the
// statements that were in the loop in a BLOCK node. The block will be removed
// in a later pass, if possible.
//
// Java: PeepholeRemoveDeadCode.tryFoldDoAway(Node)
func (p *PeepholeRemoveDeadCode) TryFoldDoAway(n rhino.Node) rhino.Node {
	do, ok := ast.AsDo(n)
	jutil.CheckArgument(ok)

	cond := do.Cond().Node()
	if NodeUtil_GetBooleanValue(cond) != base.Tri_FALSE {
		return n
	}

	block := do.Body().Node()
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
func (p *PeepholeRemoveDeadCode) TryFoldEmptyDo(n rhino.Node) rhino.Node {
	do, ok := ast.AsDo(n)
	jutil.CheckArgument(ok)

	// (Java also checks that the body is a block: it always is.)
	body := do.Body()
	if !body.HasStatements() {
		cond := do.Cond().Node()
		forNode :=
			rhino.IR_ForNode(rhino.IR_Empty().Srcref(n), cond.Detach(), rhino.IR_Empty().Srcref(n), body.Node().Detach())
		n.ReplaceWith(forNode)
		p.ReportChangeToEnclosingScope(forNode)
		return forNode
	}
	return n
}

// PeepholeRemoveDeadCode_HasUnnamedBreakOrContinue returns whether a node
// has any unhandled breaks or continue.
//
// Java: PeepholeRemoveDeadCode.hasUnnamedBreakOrContinue(Node)
func PeepholeRemoveDeadCode_HasUnnamedBreakOrContinue(n rhino.Node) bool {
	return peepholeRemoveDeadCode_hasUnnamedBreakOrContinue(n, true)
}

func peepholeRemoveDeadCode_hasUnnamedBreakOrContinue(n rhino.Node, allowBreak bool) bool {
	if (allowBreak && peepholeRemoveDeadCode_isUnnamedBreak(n)) ||
		peepholeRemoveDeadCode_isUnnamedContinue(n) {
		return true
	}
	if rhino.IR_MayBeExpression(n) || NodeUtil_IsLoopStructure(n) {
		return false
	}
	if n.IsSwitch() {
		allowBreak = false
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if peepholeRemoveDeadCode_hasUnnamedBreakOrContinue(c, allowBreak) {
			return true
		}
	}
	return false
}

// TryFoldForCondition removes always true loop conditions.
//
// Java: PeepholeRemoveDeadCode.tryFoldForCondition(Node)
func (p *PeepholeRemoveDeadCode) TryFoldForCondition(forCondition rhino.Node) {
	if p.GetSideEffectFreeBooleanValue(forCondition) == base.Tri_TRUE {
		p.ReportChangeToEnclosingScope(forCondition)
		forCondition.ReplaceWith(rhino.IR_Empty())
	}
}
