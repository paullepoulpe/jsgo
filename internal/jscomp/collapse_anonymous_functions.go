// Port of com/google/javascript/jscomp/CollapseAnonymousFunctions.java @ closure-compiler v20261028

/*
 * Copyright 2008 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

var (
	cafEmptyStringID = uint32(strpool.Intern(""))
	cafVisitTokens   = func() ast.TokenSet {
		var s ast.TokenSet
		s.Add(rhino.Token_LET)
		s.Add(rhino.Token_CONST)
		s.Add(rhino.Token_VAR)
		return s
	}()
)

// CollapseAnonymousFunctions collapses anonymous function expressions into
// named function declarations, i.e. the following:
//
//	var f = function() {}
//
// becomes:
//
//	function f() {}
//
// This reduces the generated code size but changes the semantics because f
// will be defined before its definition is reached. Also, in ES6+, "var f"
// is visible in the entire function scope, whereas "function f" is block
// scoped, which may cause issues.
type CollapseAnonymousFunctions struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler
}

// Java: CollapseAnonymousFunctions(AbstractCompiler)
func NewCollapseAnonymousFunctions(compiler AbstractCompiler) *CollapseAnonymousFunctions {
	c := &CollapseAnonymousFunctions{}
	jutil.CheckArgument(compiler.GetLifeCycleStage().IsNormalized())
	c.compiler = compiler
	return c
}

// Java: CollapseAnonymousFunctions.process(Node, Node)
func (c *CollapseAnonymousFunctions) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}

var (
	_ nodeTraversal_noVisitLeaves = (*CollapseAnonymousFunctions)(nil)
	_ nodeTraversal_visitTokens   = (*CollapseAnonymousFunctions)(nil)
)

func (*CollapseAnonymousFunctions) noVisitLeaves() {}

func (*CollapseAnonymousFunctions) visitTokens() *ast.TokenSet {
	return &cafVisitTokens
}

// Java: CollapseAnonymousFunctions.visit(NodeTraversal, Node, Node)
func (c *CollapseAnonymousFunctions) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	decl, ok := ast.AsNameDeclaration(n)
	if !ok {
		return
	}

	// It is only safe to collapse anonymous functions that appear
	// at top-level blocks.  In other cases the difference between
	// variable and function declarations can lead to problems or
	// expose subtle bugs in browser implementation as function
	// definitions are added to scopes before the start of execution.

	_, inScript := ast.AsScript(parent)
	_, inBlock := ast.AsBlock(parent)
	_, inFunction := ast.AsFunction(parent.GetParent())
	if !inScript && !(inBlock && inFunction) {
		return
	}

	// Need to store the next name in case the current name is removed from
	// the linked list.
	declarator, ok := decl.OnlyDeclarator()
	jutil.CheckState(ok)

	// Don't collapse if the lhs is a destructuring pattern.
	name, ok := declarator.AsName()
	if !ok {
		return
	}

	initExpr, ok := name.Init()
	if !ok {
		return
	}
	fn, ok := initExpr.AsFunction()
	if ok && !fn.IsArrow() && !c.isRecursiveFunction(fn) {
		fnName := fn.NameNode()
		if fnName.IsConstantName() && !name.IsConstantName() {
			// When the LHS name `f` is already a const `const f = function() {};`, then we must
			// preserve its constness during this rewriting to `function f() {}`.
			// When the LHS name `f` is non-const `var f = function() {}` or `let f = function() {}`, we
			// must reset the constness as the function will no longer be an RHS function expression
			// during this rewriting to `function f() {}`
			fnName.SetConstantName(false)
		}
		fnName.SetIdent(name.Ident())
		NodeUtil_CopyNameAnnotations(name.Node(), fnName.Node())
		value := fn.Detach().Node()
		n.ReplaceWith(value)

		// Renormalize the code.
		if !t.InGlobalScope() && NodeUtil_IsHoistedFunctionDeclaration(value) {
			parent.AddChildToFront(value.Detach())
		}

		// report changes to both the change scopes
		c.compiler.ReportChangeToChangeScope(value)
		t.ReportCodeChange()
	}
}

// Java: CollapseAnonymousFunctions.isRecursiveFunction(Node)
func (c *CollapseAnonymousFunctions) IsRecursiveFunction(function rhino.Node) bool {
	return c.isRecursiveFunction(ast.MustFunction(function))
}

func (c *CollapseAnonymousFunctions) isRecursiveFunction(fn ast.Function) bool {
	name, ok := fn.Name()
	if !ok {
		return false
	}
	return cafContainsNameID(fn.Body().Node(), name.Node().GetStringID())
}

func cafContainsNameID(n rhino.Node, nameID uint32) bool {
	if n.IsName() && n.GetStringID() == nameID {
		return true
	}

	for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		if cafContainsNameID(child, nameID) {
			return true
		}
	}
	return false
}

// Java: CollapseAnonymousFunctions.containsName(Node, String)
func (c *CollapseAnonymousFunctions) ContainsName(n rhino.Node, name string) bool {
	if name == "" {
		return cafContainsNameID(n, cafEmptyStringID)
	}
	id := strpool.Lookup(name)
	if id == 0 {
		return false
	}
	return cafContainsNameID(n, uint32(id))
}
