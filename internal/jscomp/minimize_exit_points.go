// Port of com/google/javascript/jscomp/MinimizeExitPoints.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// MinimizeExitPoints transforms the structure of the AST so that the number
// of explicit exits are minimized and instead flows to implicit exits
// conditions.
type MinimizeExitPoints struct {
	AbstractPeepholeOptimizationBase
}

// Java: MinimizeExitPoints()
func NewMinimizeExitPoints() *MinimizeExitPoints {
	m := &MinimizeExitPoints{}
	m.InitAbstractPeepholeOptimizationBase(m) // the implicit super()
	return m
}

// Java: MinimizeExitPoints.optimizeSubtree(Node)
func (m *MinimizeExitPoints) OptimizeSubtree(n *rhino.Node) *rhino.Node {
	jutil.CheckStateMsg(m.IsASTNormalized(), "MinimizeExitPoints requires the AST to be normalized")
	switch n.GetToken() {
	case rhino.Token_LABEL:
		labelName := n.GetFirstChild().GetString()
		m.TryMinimizeExits(n.GetLastChild(), rhino.Token_BREAK, &labelName)
	case rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF, rhino.Token_WHILE:
		m.TryMinimizeExits(NodeUtil_GetLoopCodeBlock(n), rhino.Token_CONTINUE, nil)
	case rhino.Token_DO:
		m.TryMinimizeExits(NodeUtil_GetLoopCodeBlock(n), rhino.Token_CONTINUE, nil)

		cond := NodeUtil_GetConditionExpression(n)
		if m.GetSideEffectFreeBooleanValue(cond) == base.Tri_FALSE {
			// Normally, we wouldn't be able to optimize BREAKs inside a loop
			// but as we know the condition will always be false, we can treat them
			// as we would a CONTINUE.
			m.TryMinimizeExits(n.GetFirstChild(), rhino.Token_BREAK, nil)
		}
	case rhino.Token_BLOCK:
		if n.HasParent() && n.GetParent().IsFunction() {
			m.TryMinimizeExits(n, rhino.Token_RETURN, nil)
		}
	case rhino.Token_SWITCH:
		m.TryMinimizeSwitchExits(n, rhino.Token_BREAK, nil)
	// TODO(johnlenz): Minimize any block that ends in a optimizable statements:
	//   break, continue, return
	default:
	}
	return n
}

// TryMinimizeExits attempts to minimize the number of explicit exit points
// in a control structure to take advantage of the implied exit at the end of
// the structure. This is accomplished by removing redundant statements, and
// moving statements following a qualifying IF node into that node. For
// example:
//
//	function () {
//	  if (x) return;
//	  else blah();
//	  foo();
//	}
//
// becomes:
//
//	function () {
//	 if (x) ;
//	 else {
//	   blah();
//	   foo();
//	 }
//
// n is the execution node of a parent to inspect; exitType is the type of
// exit to look for; labelName, if parent is a label, is the name of the label
// to look for, null otherwise. Non-null only for breaks within labels.
//
// Java: MinimizeExitPoints.tryMinimizeExits(Node, Token, String)
func (m *MinimizeExitPoints) TryMinimizeExits(n *rhino.Node, exitType rhino.Token, labelName *string) {

	// Just an 'exit'.
	if MinimizeExitPoints_MatchingExitNode(n, exitType, labelName) {
		m.ReportChangeToEnclosingScope(n)
		NodeUtil_RemoveChild(n.GetParent(), n)
		return
	}

	// Just an 'if'.
	if n.IsIf() {
		ifBlock := n.GetSecondChild()
		m.TryMinimizeExits(ifBlock, exitType, labelName)
		elseBlock := ifBlock.GetNext()
		if elseBlock != nil {
			m.TryMinimizeExits(elseBlock, exitType, labelName)
		}
		return
	}

	// Just a 'try/catch/finally'.
	if n.IsTry() {
		tryBlock := n.GetFirstChild()
		m.TryMinimizeExits(tryBlock, exitType, labelName)
		allCatchNodes := NodeUtil_GetCatchBlock(n)
		if NodeUtil_HasCatchHandler(allCatchNodes) {
			jutil.CheckState(allCatchNodes.HasOneChild())
			catchNode := allCatchNodes.GetFirstChild()
			catchCodeBlock := catchNode.GetLastChild()
			m.TryMinimizeExits(catchCodeBlock, exitType, labelName)
		}
		/* Don't try to minimize the exits of finally blocks, as this
		 * can cause problems if it changes the completion type of the finally
		 * block. See ECMA 262 Sections 8.9 & 12.14
		 */
	}

	// Just a 'label'.
	if n.IsLabel() {
		labelBlock := n.GetLastChild()
		m.TryMinimizeExits(labelBlock, exitType, labelName)
	}

	// We can only minimize switch cases if we are not trying to remove unlabeled breaks.
	if n.IsSwitch() && (exitType != rhino.Token_BREAK || labelName != nil) {
		m.TryMinimizeSwitchExits(n, exitType, labelName)
		return
	}

	// The rest assumes a block with at least one child, bail on anything else.
	if !n.IsBlock() || !n.HasChildren() {
		return
	}

	// Multiple if-exits can be converted in a single pass.
	// Convert "if (blah) break;  if (blah2) break; other_stmt;" to
	// become "if (blah); else { if (blah2); else { other_stmt; } }"
	// which will get converted to "if (!blah && !blah2) { other_stmt; }".
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		// An 'if' block to process below.
		if c.IsIf() {
			ifTree := c

			// First, the true condition block.
			trueBlock := ifTree.GetSecondChild()
			falseBlock := trueBlock.GetNext()
			m.TryMinimizeIfBlockExits(trueBlock, falseBlock,
				ifTree, exitType, labelName)

			// Now the else block.
			// The if blocks may have changed, get them again.
			trueBlock = ifTree.GetSecondChild()
			falseBlock = trueBlock.GetNext()
			if falseBlock != nil {
				m.TryMinimizeIfBlockExits(falseBlock, trueBlock,
					ifTree, exitType, labelName)
			}
		}

		if c == n.GetLastChild() {
			break
		}
	}

	// Now try to minimize the exits of the last child, if it is removed
	// look at what has become the last child.
	for c := n.GetLastChild(); c != nil; c = n.GetLastChild() {
		m.TryMinimizeExits(c, exitType, labelName)
		// If the node is still the last child, we are done.
		if c == n.GetLastChild() {
			break
		}
	}
}

// Java: MinimizeExitPoints.tryMinimizeSwitchExits(Node, Token, String)
func (m *MinimizeExitPoints) TryMinimizeSwitchExits(n *rhino.Node, exitType rhino.Token, labelName *string) {
	jutil.CheckState(n.IsSwitch())
	// Skipping the switch condition, visit all the children.
	switchBody := n.GetSecondChild()
	for c := switchBody.GetFirstChild(); c != nil; c = c.GetNext() {
		if c != switchBody.GetLastChild() {
			m.TryMinimizeSwitchCaseExits(c, exitType, labelName)
		} else {
			// Last case, the last case block can be optimized more aggressively.
			m.TryMinimizeExits(c.GetLastChild(), exitType, labelName)
		}
	}
}

// TryMinimizeSwitchCaseExits attempts to remove explicit exits from switch
// cases that also occur implicitly after the switch.
//
// Java: MinimizeExitPoints.tryMinimizeSwitchCaseExits(Node, Token, String)
func (m *MinimizeExitPoints) TryMinimizeSwitchCaseExits(n *rhino.Node, exitType rhino.Token, labelName *string) {
	jutil.CheckState(NodeUtil_IsSwitchCase(n))

	jutil.CheckState(n != n.GetParent().GetLastChild())
	block := n.GetLastChild()
	maybeBreak := block.GetLastChild()
	if maybeBreak == nil || !maybeBreak.IsBreak() || maybeBreak.HasChildren() {
		// Can not minimize exits from a case without an explicit break from the switch.
		return
	}

	// Now try to minimize the exits of the last child before the break, if it is removed
	// look at what has become the child before the break.
	childBeforeBreak := maybeBreak.GetPrevious()
	for childBeforeBreak != nil {
		c := childBeforeBreak
		m.TryMinimizeExits(c, exitType, labelName)
		// If the node is still the last child, we are done.
		childBeforeBreak = maybeBreak.GetPrevious()
		if c == childBeforeBreak {
			break
		}
	}
}

// TryMinimizeIfBlockExits looks for exits (returns, breaks, or continues,
// depending on the context) at the end of a block and removes them by moving
// the if node's siblings, if any, into the opposite condition block.
//
// srcBlock is the block to inspect; destBlock is the block to move sibling
// nodes into; ifNode is the if node to work with; exitType is the type of
// exit to look for; labelName is the name associated with the exit, if any.
// null for anything excepted for named-break associated with a label.
//
// Java: MinimizeExitPoints.tryMinimizeIfBlockExits(Node, Node, Node, Token, String)
func (m *MinimizeExitPoints) TryMinimizeIfBlockExits(srcBlock *rhino.Node, destBlock *rhino.Node,
	ifNode *rhino.Node, exitType rhino.Token, labelName *string) {
	var exitNodeParent *rhino.Node
	var exitNode *rhino.Node

	// Pick an exit node candidate.
	if srcBlock.IsBlock() {
		if !srcBlock.HasChildren() {
			return
		}
		exitNodeParent = srcBlock
		exitNode = exitNodeParent.GetLastChild()
	} else {
		// Just a single statement, if it isn't an exit bail.
		exitNodeParent = ifNode
		exitNode = srcBlock
	}
	_ = exitNodeParent

	// Verify the candidate.
	if !MinimizeExitPoints_MatchingExitNode(exitNode, exitType, labelName) {
		return
	}

	// Ensure no block-scoped declarations are moved into an inner block.
	if !MinimizeExitPoints_TryConvertAllBlockScopedFollowing(ifNode) {
		return
	}

	// Take case of the if nodes siblings, if any.
	if ifNode.GetNext() != nil {
		// Move siblings of the if block into the opposite
		// logic block of the exit.
		newDestBlock := rhino.IR_Block().Srcref(ifNode)
		if destBlock == nil {
			// Only possible if this is the false block.
			ifNode.AddChildToBack(newDestBlock)
		} else if destBlock.IsEmpty() {
			// Use the new block.
			destBlock.ReplaceWith(newDestBlock)
		} else if destBlock.IsBlock() {
			// Reuse the existing block.
			newDestBlock = destBlock
		} else {
			// Add the existing statement to the new block.
			destBlock.ReplaceWith(newDestBlock)
			newDestBlock.AddChildToBack(destBlock)
		}

		// Move all the if node's following siblings.
		MinimizeExitPoints_MoveAllFollowing(ifNode, newDestBlock)
		m.ReportChangeToEnclosingScope(ifNode)
	}
}

// MinimizeExitPoints_MatchingExitNode determines if n matches the type and
// name for the following types of "exits":
//   - return without values
//   - continues and breaks with or without names.
//
// n is the node to inspect; type_ is the Token type to look for; labelName is
// the name that must be associated with the exit type, non-null only for
// breaks associated with labels. Returns whether the node matches the
// specified block-exit type.
//
// Java: MinimizeExitPoints.matchingExitNode(Node, Token, String)
func MinimizeExitPoints_MatchingExitNode(n *rhino.Node, type_ rhino.Token, labelName *string) bool {
	if n.GetToken() == type_ {
		if type_ == rhino.Token_RETURN {
			// only returns without expressions.
			return !n.HasChildren()
		} else {
			if labelName == nil {
				return !n.HasChildren()
			} else {
				return n.HasChildren() &&
					*labelName == n.GetFirstChild().GetString()
			}
		}
	}
	return false
}

// MinimizeExitPoints_MoveAllFollowing moves all the child nodes following
// start in srcParent to the end of destParent's child list.
//
// start is the start point in the srcParent child list; destParent is the
// destination node.
//
// Java: MinimizeExitPoints.moveAllFollowing(Node, Node)
func MinimizeExitPoints_MoveAllFollowing(start *rhino.Node, destParent *rhino.Node) {
	for n := start.GetNext(); n != nil; n = start.GetNext() {
		isFunctionDeclaration := NodeUtil_IsFunctionDeclaration(n)
		n.Detach()
		if isFunctionDeclaration {
			destParent.AddChildToFront(n)
		} else {
			destParent.AddChildToBack(n)
		}
	}
}

// MinimizeExitPoints_TryConvertAllBlockScopedFollowing converts all
// let/const declarations following the start node to var declarations if
// possible.
//
// See the unit tests for examples of why this is necessary before moving code
// into an inner block, and why this is unsafe to do to declarations inside a
// loop.
//
// start is the start point. Returns whether all block-scoped declarations
// have been converted.
//
// Java: MinimizeExitPoints.tryConvertAllBlockScopedFollowing(Node)
func MinimizeExitPoints_TryConvertAllBlockScopedFollowing(start *rhino.Node) bool {
	if NodeUtil_IsWithinLoop(start) {
		// If in a loop, don't convert anything to a var. Return true only if there are no let/consts.
		return !MinimizeExitPoints_HasBlockScopedVarsFollowing(start)
	}
	for n := start.GetNext(); n != nil; n = n.GetNext() {
		if n.IsLet() || n.IsConst() {
			n.SetToken(rhino.Token_VAR)
		}
	}
	return true
}

// MinimizeExitPoints_HasBlockScopedVarsFollowing detects any block-scoped
// declarations that are younger siblings of the given starting point.
//
// Java: MinimizeExitPoints.hasBlockScopedVarsFollowing(Node)
func MinimizeExitPoints_HasBlockScopedVarsFollowing(start *rhino.Node) bool {
	for n := start.GetNext(); n != nil; n = n.GetNext() {
		if n.IsLet() || n.IsConst() {
			return true
		}
	}
	return false
}
