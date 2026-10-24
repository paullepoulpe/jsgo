// Port of com/google/javascript/jscomp/CollapseAnonymousFunctions.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
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

// Java: CollapseAnonymousFunctions.visit(NodeTraversal, Node, Node)
func (c *CollapseAnonymousFunctions) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if !(n.IsVar() || n.IsLet() || n.IsConst()) {
		return
	}

	// It is only safe to collapse anonymous functions that appear
	// at top-level blocks.  In other cases the difference between
	// variable and function declarations can lead to problems or
	// expose subtle bugs in browser implementation as function
	// definitions are added to scopes before the start of execution.

	grandparent := parent.GetParent()
	if !(parent.IsScript() ||
		(!grandparent.IsNil() && grandparent.IsFunction() && parent.IsBlock())) {
		return
	}

	// Need to store the next name in case the current name is removed from
	// the linked list.
	name := n.GetOnlyChild()

	// Don't collapse if the lhs is a destructuring pattern.
	if !name.IsName() {
		return
	}

	value := name.GetFirstChild()
	if !value.IsNil() &&
		value.IsFunction() &&
		!value.IsArrowFunction() &&
		!c.IsRecursiveFunction(value) {
		fnName := value.GetFirstChild()
		if fnName.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME) &&
			!name.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME) {
			// When the LHS name `f` is already a const `const f = function() {};`, then we must
			// preserve its constness during this rewriting to `function f() {}`.
			// When the LHS name `f` is non-const `var f = function() {}` or `let f = function() {}`, we
			// must reset the constness as the function will no longer be an RHS function expression
			// during this rewriting to `function f() {}`
			fnName.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, false)
		}
		fnName.SetString(name.GetString())
		NodeUtil_CopyNameAnnotations(name, fnName)
		value.Detach()
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
	name := function.GetFirstChild()
	if name.GetString() == "" {
		return false
	}
	args := name.GetNext()
	body := args.GetNext()
	return c.ContainsName(body, name.GetString())
}

// Java: CollapseAnonymousFunctions.containsName(Node, String)
func (c *CollapseAnonymousFunctions) ContainsName(n rhino.Node, name string) bool {
	if n.IsName() && n.GetString() == name {
		return true
	}

	for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		if c.ContainsName(child, name) {
			return true
		}
	}
	return false
}
