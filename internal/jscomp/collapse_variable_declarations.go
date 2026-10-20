// Port of com/google/javascript/jscomp/CollapseVariableDeclarations.java @ closure-compiler v20261015

/*
 * Copyright 2006 The Closure Compiler Authors.
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

// CollapseVariableDeclarations collapses multiple variable declarations into
// a single one. i.e the following:
//
//	var a;
//	var b = 1;
//	var c = 2;
//
// becomes:
//
//	var a, b = 1, c = 2;
//
// This reduces the generated code size. More optimizations are possible:
//
//   - Group all variable declarations inside a function into one such
//     variable. declaration block.
//
//   - Re-use variables instead of declaring a new one if they are used for
//     only part of a function. Similarly, also collapses assigns like:
//
//     a = true;
//     b = true;
//     var c = true;
//
//     becomes:
//
//     var c = b = a = true;
type CollapseVariableDeclarations struct {
	// Reference to JS Compiler
	compiler AbstractCompiler

	// Collapses to do in this pass.
	collapses     *jutil.ArrayList[*CollapseVariableDeclarations_Collapse]
	collapseNodes []*rhino.Node
	collapseBuf   [32]*rhino.Node

	// Nodes we've already looked at for collapsing, so that we don't look at
	// them again (we look ahead when examining what nodes can be collapsed,
	// and the node traversal may give them to us again)
	nodesToCollapse *jutil.LinkedHashSet[*rhino.Node]
	collapsedSet    map[*rhino.Node]struct{}
}

// CollapseVariableDeclarations_Collapse is the encapsulation of information
// about a variable declaration collapse.
type CollapseVariableDeclarations_Collapse struct {
	// Variable declaration that any following var nodes should be collapsed
	// into
	startNode *rhino.Node
}

// Java: CollapseVariableDeclarations.Collapse(Node)
func NewCollapseVariableDeclarations_Collapse(startNode *rhino.Node) *CollapseVariableDeclarations_Collapse {
	return &CollapseVariableDeclarations_Collapse{startNode: startNode}
}

// Java: CollapseVariableDeclarations(AbstractCompiler)
func NewCollapseVariableDeclarations(compiler AbstractCompiler) *CollapseVariableDeclarations {
	jutil.CheckState(!compiler.GetLifeCycleStage().IsNormalized())
	c := &CollapseVariableDeclarations{
		compiler: compiler,
	}
	c.collapseNodes = c.collapseBuf[:0]
	return c
}

// Java: CollapseVariableDeclarations.process(Node, Node)
func (c *CollapseVariableDeclarations) Process(externs *rhino.Node, root *rhino.Node) {
	c.collapseNodes = c.collapseNodes[:0]
	clear(c.collapsedSet)

	NodeTraversal_Traverse(c.compiler, root, &CollapseVariableDeclarations_GatherCollapses{outer: c})

	if len(c.collapseNodes) > 0 {
		c.ApplyCollapses()
	}
}

// CollapseVariableDeclarations_GatherCollapses gathers all of the variable
// declarations that should be collapsed into one.
//
// We do not do the collapsing as we go since node traversal would be
// affected by the changes we are making to the parse tree.
type CollapseVariableDeclarations_GatherCollapses struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer *CollapseVariableDeclarations
}

// Java: CollapseVariableDeclarations.GatherCollapses.visit(NodeTraversal, Node, Node)
func (g *CollapseVariableDeclarations_GatherCollapses) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	nType := n.GetToken()
	if nType != rhino.Token_VAR && nType != rhino.Token_LET && nType != rhino.Token_CONST {
		return
	}

	// Adjacent VAR children of an IF node are the if and else parts and can't
	// be collapsed
	if parent.IsIf() {
		return
	}

	outer := g.outer
	// If we've already looked at this node, skip it
	if len(outer.collapsedSet) > 0 {
		if _, ok := outer.collapsedSet[n]; ok {
			return
		}
	}

	next := n.GetNext()
	if next == nil || next.GetToken() != nType {
		return
	}

	if outer.collapsedSet == nil {
		outer.collapsedSet = make(map[*rhino.Node]struct{}, 32)
	}
	for next != nil && next.GetToken() == nType {
		outer.collapsedSet[next] = struct{}{}
		next = next.GetNext()
	}
	outer.collapseNodes = append(outer.collapseNodes, n)
}

// Java: CollapseVariableDeclarations.applyCollapses()
func (c *CollapseVariableDeclarations) ApplyCollapses() {
	for _, var_ := range c.collapseNodes {
		c.compiler.ReportChangeToEnclosingScope(var_)

		varTok := var_.GetToken()
		for var_.GetNext() != nil && var_.GetNext().GetToken() == varTok {
			next := var_.GetNext().Detach()

			// Move all children of the next var node into the first one.
			var_.AddChildrenToBack(next.RemoveChildren())
		}
	}
}
