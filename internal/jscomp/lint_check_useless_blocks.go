// Port of com/google/javascript/jscomp/lint/CheckUselessBlocks.java @ closure-compiler v20261015

/*
 * Copyright 2016 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CheckUselessBlocks checks for useless blocks. A block is considered useful
// if it is part of a control structure like if / else / while / switch / etc.
// or if it contains any block-scoped variables (let, const, class, or
// function declarations). Otherwise there is no reason to use it so it is
// likely a mistake. This would catch the classic error:
//
// return {foo: 'bar'};
//
// or more contrived cases like:
//
// if (denied) { showAccessDenied(); } { grantAccess(); }
//
// Inspired by ESLint (https://github.com/eslint/eslint/blob/master/lib/rules/no-lone-blocks.js)
type CheckUselessBlocks struct {
	compiler AbstractCompiler
	// Go-only: language-adapter — the ArrayDeque used as a stack (push, pop
	// and peek at its head) is a slice whose last element is the head.
	loneBlocks []rhino.Node
}

var CheckUselessBlocks_USELESS_BLOCK = DiagnosticType_Disabled(
	"JSC_USELESS_BLOCK", "Useless block.")

// NewCheckUselessBlocks is new CheckUselessBlocks(compiler).
//
// Java: CheckUselessBlocks(AbstractCompiler)
func NewCheckUselessBlocks(compiler AbstractCompiler) *CheckUselessBlocks {
	c := &CheckUselessBlocks{}
	c.compiler = compiler
	c.loneBlocks = []rhino.Node{}
	return c
}

// Java: CheckUselessBlocks.process(Node, Node)
func (c *CheckUselessBlocks) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}

// IsLoneBlock: a lone block is a non-synthetic, not-added BLOCK that is a
// direct child of another non-synthetic, not-added BLOCK or a SCRIPT node.
//
// Java: CheckUselessBlocks.isLoneBlock(Node)
func (c *CheckUselessBlocks) IsLoneBlock(n rhino.Node) bool {
	parent := n.GetParent()
	if !parent.IsNil() && (parent.IsScript() ||
		(parent.IsBlock() && !parent.IsSyntheticBlock() && !parent.IsAddedBlock())) {
		return !n.IsSyntheticBlock() && !n.IsAddedBlock()
	}
	return false
}

// AllowLoneBlock removes the enclosing block of a block-scoped declaration
// from the loneBlocks stack.
//
// Java: CheckUselessBlocks.allowLoneBlock(Node)
func (c *CheckUselessBlocks) AllowLoneBlock(parent rhino.Node) {
	if len(c.loneBlocks) == 0 {
		return
	}

	if c.loneBlocks[len(c.loneBlocks)-1] == parent {
		c.popLoneBlock()
	}
}

// popLoneBlock is loneBlocks.pop() on a non-empty deque.
//
// Go-only: language-adapter — see loneBlocks.
func (c *CheckUselessBlocks) popLoneBlock() {
	last := len(c.loneBlocks) - 1
	c.loneBlocks[last] = rhino.NilNode
	c.loneBlocks = c.loneBlocks[:last]
}

// Java: CheckUselessBlocks.shouldTraverse(NodeTraversal, Node, Node)
func (c *CheckUselessBlocks) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_BLOCK:
		if c.IsLoneBlock(n) {
			c.loneBlocks = append(c.loneBlocks, n)
		}
	case rhino.Token_LET, rhino.Token_CONST:
		c.AllowLoneBlock(parent)
	case rhino.Token_CLASS:
		if NodeUtil_IsClassDeclaration(n) {
			c.AllowLoneBlock(parent)
		}
	case rhino.Token_FUNCTION:
		if NodeUtil_IsFunctionDeclaration(n) {
			c.AllowLoneBlock(parent)
		}
	default:
	}
	return true
}

// Java: CheckUselessBlocks.visit(NodeTraversal, Node, Node)
func (c *CheckUselessBlocks) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsBlock() && len(c.loneBlocks) != 0 && c.loneBlocks[len(c.loneBlocks)-1] == n {
		c.popLoneBlock()
		t.Report(n, CheckUselessBlocks_USELESS_BLOCK)
	}
}
