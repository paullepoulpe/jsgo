// Port of com/google/javascript/jscomp/CollapseVariableDeclarations.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: this pass uses the typed AST view (internal/ast): it gathers and
// collapses ast.NameDeclarations, and the "next declaration of the same
// kind" test that both phases need is one helper (nextOfSameKind) instead of
// token comparisons. The test that runs on every node stays a token
// comparison (see GatherCollapses.Visit), and moving the declarators stays
// one rhino splice (AddChildrenToBack of RemoveChildren).

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

	// Go-only: the declarations that the following ones collapse into (the
	// startNode of each Collapse), in a reused buffer.
	collapseNodes []ast.NameDeclaration
	collapseBuf   [32]ast.NameDeclaration

	// Nodes we've already looked at for collapsing, so that we don't look at
	// them again (we look ahead when examining what nodes can be collapsed,
	// and the node traversal may give them to us again)
	collapsedSet map[rhino.Node]struct{}
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
func (c *CollapseVariableDeclarations) Process(externs rhino.Node, root rhino.Node) {
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
func (g *CollapseVariableDeclarations_GatherCollapses) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	// Go-only: this runs on every node, so the kind is tested on the token
	// (three comparisons are a little cheaper than the token-set lookup of
	// ast.AsNameDeclaration, and a typed visitor costs more still).
	if tok := n.GetToken(); tok != rhino.Token_VAR && tok != rhino.Token_LET && tok != rhino.Token_CONST {
		return
	}
	decl := ast.MustNameDeclaration(n)

	// Adjacent VAR children of an IF node are the if and else parts and can't
	// be collapsed
	// Go-only: the branches of an if are blocks (grammar.txt), so a
	// declaration is never a child of an IF; the check is kept from Java.
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

	next, ok := collapseVariableDeclarations_nextOfSameKind(decl)
	if !ok {
		return
	}

	if outer.collapsedSet == nil {
		outer.collapsedSet = make(map[rhino.Node]struct{}, 32)
	}
	for ok {
		outer.collapsedSet[next.Node()] = struct{}{}
		next, ok = collapseVariableDeclarations_nextOfSameKind(next)
	}
	outer.collapseNodes = append(outer.collapseNodes, decl)
}

// collapseVariableDeclarations_nextOfSameKind returns the statement after
// decl if it is a declaration of the same kind (var, let or const); ok is
// false otherwise.
//
// Go-only: Java compares the tokens where it needs this.
func collapseVariableDeclarations_nextOfSameKind(decl ast.NameDeclaration) (ast.NameDeclaration, bool) {
	next, ok := ast.AsNameDeclaration(decl.Node().GetNext())
	return next, ok && next.Token() == decl.Token()
}

// Java: CollapseVariableDeclarations.applyCollapses()
func (c *CollapseVariableDeclarations) ApplyCollapses() {
	for _, decl := range c.collapseNodes {
		c.compiler.ReportChangeToEnclosingScope(decl.Node())

		for next, ok := collapseVariableDeclarations_nextOfSameKind(decl); ok; next, ok = collapseVariableDeclarations_nextOfSameKind(decl) {
			next.Detach()

			// Move all children of the next var node into the first one.
			decl.Node().AddChildrenToBack(next.Node().RemoveChildren())
		}
	}
}
