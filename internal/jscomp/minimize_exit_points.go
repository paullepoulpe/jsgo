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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: this pass reads and edits the tree through the typed view of
// package ast (ifs.Then(), try.Catch()...) instead of child positions.

// MinimizeExitPoints transforms the structure of the AST so that the number
// of explicit exits are minimized and instead flows to implicit exits
// conditions.
type MinimizeExitPoints struct {
	AbstractPeepholeOptimizationBase
}

// Java: MinimizeExitPoints()
func NewMinimizeExitPoints() *MinimizeExitPoints {
	return &MinimizeExitPoints{}
}

var minimizeExitPointsTokens = peepholeTokenSet(
	rhino.Token_LABEL,
	rhino.Token_FOR,
	rhino.Token_FOR_IN,
	rhino.Token_FOR_OF,
	rhino.Token_FOR_AWAIT_OF,
	rhino.Token_WHILE,
	rhino.Token_DO,
	rhino.Token_BLOCK,
	rhino.Token_SWITCH,
	rhino.Token_SCRIPT,
)

func (m *MinimizeExitPoints) peepholeConfig() peepholeOptConfig {
	return peepholeOptConfig{
		tokens:     minimizeExitPointsTokens,
		isStandard: true,
	}
}

func (m *MinimizeExitPoints) cloneForWorker() AbstractPeepholeOptimization {
	return NewMinimizeExitPoints()
}

// Java: MinimizeExitPoints.optimizeSubtree(Node)
func (m *MinimizeExitPoints) OptimizeSubtree(n rhino.Node) rhino.Node {
	jutil.CheckStateMsg(m.IsASTNormalized(), "MinimizeExitPoints requires the AST to be normalized")
	if label, ok := ast.AsLabel(n); ok {
		m.TryMinimizeExits(label.Body(), rhino.Token_BREAK, label.Name().Node().GetStringID(), true)
	} else if loop, ok := ast.AsLoop(n); ok {
		m.TryMinimizeExits(loop.Body().AsStatement(), rhino.Token_CONTINUE, 0, false)
		if do, ok := loop.AsDo(); ok {
			if condVal, condOk := m.GetSideEffectFreeBooleanValue(do.Cond().Node()); condOk && !condVal {
				// Normally, we wouldn't be able to optimize BREAKs inside a loop
				// but as we know the condition will always be false, we can treat them
				// as we would a CONTINUE.
				m.TryMinimizeExits(do.Body().AsStatement(), rhino.Token_BREAK, 0, false)
			}
		}
	} else if block, ok := ast.AsBlock(n); ok {
		if n.HasParent() && n.GetParent().IsFunction() {
			m.TryMinimizeExits(block.AsStatement(), rhino.Token_RETURN, 0, false)
		}
	} else if sw, ok := ast.AsSwitch(n); ok {
		m.TryMinimizeSwitchExits(sw, rhino.Token_BREAK, 0, false)
	}
	// TODO(johnlenz): Minimize any block that ends in a optimizable statements:
	//   break, continue, return
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
// exit to look for; labelID/hasLabel specify the label name ID to look for
// when parent is a label (hasLabel true only for breaks within labels).
//
// Java: MinimizeExitPoints.tryMinimizeExits(Node, Token, String)
func (m *MinimizeExitPoints) TryMinimizeExits(n ast.Statement, exitType rhino.Token, labelID uint32, hasLabel bool) {
	// Just an 'exit'.
	if MinimizeExitPoints_MatchingExitNode(n, exitType, labelID, hasLabel) {
		m.ReportChangeToEnclosingScope(n.Node())
		NodeUtil_RemoveChild(n.Node().GetParent(), n.Node())
		return
	}

	// Just an 'if'.
	if ifs, ok := n.AsIf(); ok {
		m.TryMinimizeExits(ifs.Then().AsStatement(), exitType, labelID, hasLabel)
		if els, ok := ifs.Else(); ok {
			m.TryMinimizeExits(els.AsStatement(), exitType, labelID, hasLabel)
		}
		return
	}

	// Just a 'try/catch/finally'.
	if try, ok := n.AsTry(); ok {
		m.TryMinimizeExits(try.Body().AsStatement(), exitType, labelID, hasLabel)
		if catch, ok := try.Catch(); ok {
			m.TryMinimizeExits(catch.Body().AsStatement(), exitType, labelID, hasLabel)
		}
		/* Don't try to minimize the exits of finally blocks, as this
		 * can cause problems if it changes the completion type of the finally
		 * block. See ECMA 262 Sections 8.9 & 12.14
		 */
		return
	}

	// Just a 'label'.
	if label, ok := n.AsLabel(); ok {
		m.TryMinimizeExits(label.Body(), exitType, labelID, hasLabel)
		return
	}

	// We can only minimize switch cases if we are not trying to remove unlabeled breaks.
	if sw, ok := n.AsSwitch(); ok && (exitType != rhino.Token_BREAK || hasLabel) {
		m.TryMinimizeSwitchExits(sw, exitType, labelID, hasLabel)
		return
	}

	// The rest assumes a block with at least one child, bail on anything else.
	block, ok := n.AsBlock()
	if !ok || !block.HasStatements() {
		return
	}

	// Multiple if-exits can be converted in a single pass.
	// Convert "if (blah) break;  if (blah2) break; other_stmt;" to
	// become "if (blah); else { if (blah2); else { other_stmt; } }"
	// which will get converted to "if (!blah && !blah2) { other_stmt; }".
	for c := range block.Statements() {
		// An 'if' block to process below.
		if ifs, ok := c.AsIf(); ok {
			// First, the true condition block.
			m.TryMinimizeIfBlockExits(ifs, true, exitType, labelID, hasLabel)
			// Now the else block (the if may have one now).
			if _, ok := ifs.Else(); ok {
				m.TryMinimizeIfBlockExits(ifs, false, exitType, labelID, hasLabel)
			}
		}

		if last, _ := block.LastStatement(); c == last {
			break
		}
	}

	// Now try to minimize the exits of the last child, if it is removed
	// look at what has become the last child.
	for c, ok := block.LastStatement(); ok; c, ok = block.LastStatement() {
		m.TryMinimizeExits(c, exitType, labelID, hasLabel)
		// If the node is still the last child, we are done.
		if last, _ := block.LastStatement(); c == last {
			break
		}
	}
}

// Java: MinimizeExitPoints.tryMinimizeSwitchExits(Node, Token, String)
func (m *MinimizeExitPoints) TryMinimizeSwitchExits(n ast.Switch, exitType rhino.Token, labelID uint32, hasLabel bool) {
	// Skipping the switch condition, visit all the cases.
	cases := n.Body()
	for c := range cases.Cases() {
		if last, _ := cases.LastCase(); c != last {
			m.TryMinimizeSwitchCaseExits(c, exitType, labelID, hasLabel)
		} else {
			// Last case, the last case block can be optimized more aggressively.
			m.TryMinimizeExits(c.Body().AsStatement(), exitType, labelID, hasLabel)
		}
	}
}

// TryMinimizeSwitchCaseExits attempts to remove explicit exits from switch
// cases that also occur implicitly after the switch.
//
// Java: MinimizeExitPoints.tryMinimizeSwitchCaseExits(Node, Token, String)
func (m *MinimizeExitPoints) TryMinimizeSwitchCaseExits(n ast.SwitchCase, exitType rhino.Token, labelID uint32, hasLabel bool) {
	_, hasNext := n.Next()
	jutil.CheckState(hasNext)
	maybeBreak, ok := n.Body().LastStatement()
	if !ok {
		return
	}
	if brk, ok := maybeBreak.AsBreak(); !ok || hasBreakLabel(brk) {
		// Can not minimize exits from a case without an explicit break from the switch.
		return
	}

	// Now try to minimize the exits of the last child before the break, if it is removed
	// look at what has become the child before the break.
	c, ok := maybeBreak.Previous()
	for ok {
		m.TryMinimizeExits(c, exitType, labelID, hasLabel)
		// If the node is still the last child before the break, we are done.
		before, stillOne := maybeBreak.Previous()
		if !stillOne || before == c {
			break
		}
		c = before
	}
}

func hasBreakLabel(brk ast.Break) bool {
	_, ok := brk.Label()
	return ok
}

// TryMinimizeIfBlockExits looks for exits (returns, breaks, or continues,
// depending on the context) at the end of a branch of ifNode and removes
// them by moving the if node's siblings, if any, into the other branch.
//
// thenBranch says which branch to inspect; exitType is the type of exit to
// look for; labelID/hasLabel specify the label name ID associated with the
// exit, if any (hasLabel true only for named-break associated with a label).
//
// Java: MinimizeExitPoints.tryMinimizeIfBlockExits(Node, Node, Node, Token, String)
// (Go: the branch is named by thenBranch; it is always a block, and so is
// the other branch, when there is one.)
func (m *MinimizeExitPoints) TryMinimizeIfBlockExits(ifNode ast.If, thenBranch bool,
	exitType rhino.Token, labelID uint32, hasLabel bool) {
	src := ifNode.Then()
	if !thenBranch {
		src, _ = ifNode.Else()
	}

	// Pick an exit node candidate, and verify it.
	exitNode, ok := src.LastStatement()
	if !ok || !MinimizeExitPoints_MatchingExitNode(exitNode, exitType, labelID, hasLabel) {
		return
	}

	// Ensure no block-scoped declarations are moved into an inner block.
	if !MinimizeExitPoints_TryConvertAllBlockScopedFollowing(ifNode.AsStatement()) {
		return
	}

	// Take case of the if nodes siblings, if any.
	if _, ok := ifNode.AsStatement().Next(); ok {
		// Move siblings of the if block into the opposite
		// logic block of the exit.
		var dest ast.Block
		if !thenBranch {
			dest = ifNode.Then()
		} else if els, ok := ifNode.Else(); ok {
			// Reuse the existing block.
			dest = els
		} else {
			dest = ast.NewBlock()
			dest.Node().Srcref(ifNode.Node())
			ifNode.SetElse(dest)
		}

		// Move all the if node's following siblings.
		MinimizeExitPoints_MoveAllFollowing(ifNode.AsStatement(), dest)
		m.ReportChangeToEnclosingScope(ifNode.Node())
	}
}

// MinimizeExitPoints_MatchingExitNode determines if n matches the type and
// name for the following types of "exits":
//   - return without values
//   - continues and breaks with or without names.
//
// n is the node to inspect; exitType is the Token type to look for;
// labelID/hasLabel specify the label name ID that must be associated with the
// exit type (hasLabel true only for breaks associated with labels). Returns
// whether the node matches the specified block-exit type.
//
// Java: MinimizeExitPoints.matchingExitNode(Node, Token, String)
func MinimizeExitPoints_MatchingExitNode(n ast.Statement, exitType rhino.Token, labelID uint32, hasLabel bool) bool {
	if n.Node().GetToken() != exitType {
		return false
	}
	if ret, ok := n.AsReturn(); ok {
		// only returns without expressions.
		_, hasValue := ret.Value()
		return !hasValue
	}
	jump, ok := n.AsJump()
	if !ok {
		return false
	}
	jumpLabel, hasJumpLabel := jump.Label()
	if !hasLabel {
		return !hasJumpLabel
	}
	return hasJumpLabel && labelID != 0 && labelID == jumpLabel.Node().GetStringID()
}

// MinimizeExitPoints_MoveAllFollowing moves all the statements following
// start to the end of dest (function declarations to its start).
//
// Java: MinimizeExitPoints.moveAllFollowing(Node, Node)
func MinimizeExitPoints_MoveAllFollowing(start ast.Statement, dest ast.Block) {
	for n, ok := start.Next(); ok; n, ok = start.Next() {
		isFunctionDeclaration := NodeUtil_IsFunctionDeclaration(n.Node())
		n.Node().Detach()
		if isFunctionDeclaration {
			dest.PrependStatement(n)
		} else {
			dest.AppendStatement(n)
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
func MinimizeExitPoints_TryConvertAllBlockScopedFollowing(start ast.Statement) bool {
	if NodeUtil_IsWithinLoop(start.Node()) {
		// If in a loop, don't convert anything to a var. Return true only if there are no let/consts.
		return !MinimizeExitPoints_HasBlockScopedVarsFollowing(start)
	}
	for n, ok := start.Next(); ok; n, ok = n.Next() {
		if decl, ok := n.AsNameDeclaration(); ok && decl.Token() != rhino.Token_VAR {
			decl.Node().SetToken(rhino.Token_VAR)
		}
	}
	return true
}

// MinimizeExitPoints_HasBlockScopedVarsFollowing detects any block-scoped
// declarations that are younger siblings of the given starting point.
//
// Java: MinimizeExitPoints.hasBlockScopedVarsFollowing(Node)
func MinimizeExitPoints_HasBlockScopedVarsFollowing(start ast.Statement) bool {
	for n, ok := start.Next(); ok; n, ok = n.Next() {
		if decl, ok := n.AsNameDeclaration(); ok && decl.Token() != rhino.Token_VAR {
			return true
		}
	}
	return false
}
