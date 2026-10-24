// Port of com/google/javascript/jscomp/ManageClosureUnawareCode.java @ closure-compiler v20261015

/*
 * Copyright 2024 The Closure Compiler Authors.
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

// ManageClosureUnawareCode is the pass that "unwraps" closure-unaware code
// into a regular AST.
//
// ManageClosureUnawareCode.unwrap() should be run after all other passes have
// run, to unwrap the code and re-expose it to the code-printing stage of the
// compiler.
type ManageClosureUnawareCode struct {
	compiler AbstractCompiler
}

// Java: ManageClosureUnawareCode(AbstractCompiler)
func NewManageClosureUnawareCode(compiler AbstractCompiler) *ManageClosureUnawareCode {
	m := &ManageClosureUnawareCode{}
	m.compiler = compiler
	return m
}

// Java: ManageClosureUnawareCode.unwrap(AbstractCompiler)
func ManageClosureUnawareCode_Unwrap(compiler AbstractCompiler) *ManageClosureUnawareCode {
	return NewManageClosureUnawareCode(compiler)
}

// Java: ManageClosureUnawareCode.process(Node, Node)
func (m *ManageClosureUnawareCode) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(m.compiler, root, NewManageClosureUnawareCode_UnwrapConcealedClosureUnawareCode())
}

// TODO jameswr: Given that the CodePrinter supports printing out the shadow instead of the shadow
// host node, do we even need to revert the AST back to the original form at the end of
// compilation?

// ManageClosureUnawareCode_UnwrapConcealedClosureUnawareCode is the
// traversal that replaces each closure-unaware shadow host with the original
// code.
type ManageClosureUnawareCode_UnwrapConcealedClosureUnawareCode struct{}

// Java: UnwrapConcealedClosureUnawareCode()
func NewManageClosureUnawareCode_UnwrapConcealedClosureUnawareCode() *ManageClosureUnawareCode_UnwrapConcealedClosureUnawareCode {
	return &ManageClosureUnawareCode_UnwrapConcealedClosureUnawareCode{}
}

// Java: UnwrapConcealedClosureUnawareCode.shouldTraverse(NodeTraversal, Node, Node)
func (u *ManageClosureUnawareCode_UnwrapConcealedClosureUnawareCode) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if parent.IsNil() || !parent.IsScript() {
		return true // keep going
	}

	if parent.IsScript() && !parent.IsClosureUnawareCode() {
		return false
	}

	// Once inside a closureUnaware script, we want to traverse the entire thing to make sure we
	// find all the nodes marked as closure-unaware shadows.
	return true
}

// Java: UnwrapConcealedClosureUnawareCode.visit(NodeTraversal, Node, Node)
func (u *ManageClosureUnawareCode_UnwrapConcealedClosureUnawareCode) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	u.TryUnwrapClosureUnawareShadowedCode(n)
}

// Java: UnwrapConcealedClosureUnawareCode.tryUnwrapClosureUnawareShadowedCode(Node)
func (u *ManageClosureUnawareCode_UnwrapConcealedClosureUnawareCode) TryUnwrapClosureUnawareShadowedCode(n rhino.Node) {
	shadowAstRoot := n.GetClosureUnawareShadow()
	if shadowAstRoot.IsNil() {
		return
	}

	// ROOT -> SCRIPT -> EXPR_RESULT -> CALL -> FUNCTION
	shadowScript := shadowAstRoot.GetOnlyChild()
	jutil.CheckStateMsg1(shadowScript.IsScript(), "%s", shadowScript)
	jutil.CheckStateMsg1(shadowScript.HasOneChild(), "%s", shadowScript)
	exprResult := shadowScript.GetOnlyChild()
	jutil.CheckStateMsg1(exprResult.IsExprResult(), "%s", exprResult)
	callNode := exprResult.GetOnlyChild()
	jutil.CheckStateMsg1(callNode.IsCall(), "%s", callNode)
	originalCodeFunction := callNode.GetLastChild()
	jutil.CheckStateMsg1(originalCodeFunction.IsFunction(), "%s", originalCodeFunction)
	originalCodeFunction.Detach()
	n.ReplaceWith(originalCodeFunction)
}
