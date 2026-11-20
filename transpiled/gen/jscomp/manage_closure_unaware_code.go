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
// Port of com/google/javascript/jscomp/ManageClosureUnawareCode.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.ManageClosureUnawareCode
type ManageClosureUnawareCode struct {
	compiler AbstractCompiler `gson:"compiler"`
	idHash_  atomic.Int32
}

func (this_ *ManageClosureUnawareCode) BindSelf_(s any) {
	_ = s
}
func (this_ *ManageClosureUnawareCode) IsJavaNil_() bool { return this_ == nil }
func (*ManageClosureUnawareCode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ManageClosureUnawareCode", "ManageClosureUnawareCode"
}
func (this_ *ManageClosureUnawareCode) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ManageClosureUnawareCode) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ManageClosureUnawareCode", "com.google.javascript.jscomp.ManageClosureUnawareCode", func(v_ any) bool {
		return j2gort.InstanceOf[*ManageClosureUnawareCode](v_)
	})
}
func (*ManageClosureUnawareCode) Implements_jscomp_CompilerPass() {}

func (this_ *ManageClosureUnawareCode) initFields_ManageClosureUnawareCode() {
	_ = this_
}

// Java constructor: ManageClosureUnawareCode(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewManageClosureUnawareCode(compiler AbstractCompiler) *ManageClosureUnawareCode {
	this_ := &ManageClosureUnawareCode{}
	InitManageClosureUnawareCode(this_, compiler)
	return this_
}

func InitManageClosureUnawareCode(this_ *ManageClosureUnawareCode, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_ManageClosureUnawareCode()
	this_.compiler = compiler
}

// Java static method: ManageClosureUnawareCode.unwrap(com.google.javascript.jscomp.AbstractCompiler compiler)
func ManageClosureUnawareCode_Unwrap(compiler AbstractCompiler) *ManageClosureUnawareCode {
	_ = compiler
	return NewManageClosureUnawareCode(compiler)
}

// Java method: ManageClosureUnawareCode.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ManageClosureUnawareCode) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, NewManageClosureUnawareCode_UnwrapConcealedClosureUnawareCode())
}

// Java class: com.google.javascript.jscomp.ManageClosureUnawareCode.UnwrapConcealedClosureUnawareCode
type ManageClosureUnawareCode_UnwrapConcealedClosureUnawareCode struct {
	idHash_ atomic.Int32
}

func (this_ *ManageClosureUnawareCode_UnwrapConcealedClosureUnawareCode) BindSelf_(s any) {
	_ = s
}
func (this_ *ManageClosureUnawareCode_UnwrapConcealedClosureUnawareCode) IsJavaNil_() bool {
	return this_ == nil
}
func (*ManageClosureUnawareCode_UnwrapConcealedClosureUnawareCode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ManageClosureUnawareCode$UnwrapConcealedClosureUnawareCode", "UnwrapConcealedClosureUnawareCode"
}
func (this_ *ManageClosureUnawareCode_UnwrapConcealedClosureUnawareCode) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ManageClosureUnawareCode_UnwrapConcealedClosureUnawareCode) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ManageClosureUnawareCode.UnwrapConcealedClosureUnawareCode", "com.google.javascript.jscomp.ManageClosureUnawareCode$UnwrapConcealedClosureUnawareCode", func(v_ any) bool {
		return j2gort.InstanceOf[*ManageClosureUnawareCode_UnwrapConcealedClosureUnawareCode](v_)
	})
}
func (*ManageClosureUnawareCode_UnwrapConcealedClosureUnawareCode) Implements_jscomp_NodeTraversal_Callback() {
}

func (this_ *ManageClosureUnawareCode_UnwrapConcealedClosureUnawareCode) initFields_ManageClosureUnawareCode_UnwrapConcealedClosureUnawareCode() {
	_ = this_
}

// Java constructor: UnwrapConcealedClosureUnawareCode()
func NewManageClosureUnawareCode_UnwrapConcealedClosureUnawareCode() *ManageClosureUnawareCode_UnwrapConcealedClosureUnawareCode {
	this_ := &ManageClosureUnawareCode_UnwrapConcealedClosureUnawareCode{}
	InitManageClosureUnawareCode_UnwrapConcealedClosureUnawareCode(this_)
	return this_
}

func InitManageClosureUnawareCode_UnwrapConcealedClosureUnawareCode(this_ *ManageClosureUnawareCode_UnwrapConcealedClosureUnawareCode) {
	this_.initFields_ManageClosureUnawareCode_UnwrapConcealedClosureUnawareCode()
}

// Java method: UnwrapConcealedClosureUnawareCode.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node parent)
func (this_ *ManageClosureUnawareCode_UnwrapConcealedClosureUnawareCode) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	if (parent == 0) || !(parent.IsScript()) {
		return true
	}
	if parent.IsScript() && !(parent.IsClosureUnawareCode()) {
		return false
	}
	return true
}

// Java method: UnwrapConcealedClosureUnawareCode.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node parent)
func (this_ *ManageClosureUnawareCode_UnwrapConcealedClosureUnawareCode) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	this_.priv_UnwrapConcealedClosureUnawareCode_tryUnwrapClosureUnawareShadowedCode_3862(n)
}

// Java method: UnwrapConcealedClosureUnawareCode.tryUnwrapClosureUnawareShadowedCode(com.google.javascript.rhino.Node n)
func (this_ *ManageClosureUnawareCode_UnwrapConcealedClosureUnawareCode) priv_UnwrapConcealedClosureUnawareCode_tryUnwrapClosureUnawareShadowedCode_3862(n pkg_rhino.Node) {
	_ = n
	var shadowAstRoot pkg_rhino.Node = n.GetClosureUnawareShadow()
	_ = shadowAstRoot
	if shadowAstRoot == 0 {
		return
	}
	var shadowScript pkg_rhino.Node = shadowAstRoot.GetOnlyChild()
	_ = shadowScript
	if !(shadowScript.IsScript()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(shadowScript))
	}
	if !(shadowScript.HasOneChild()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(shadowScript))
	}
	var exprResult pkg_rhino.Node = shadowScript.GetOnlyChild()
	_ = exprResult
	if !(exprResult.IsExprResult()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(exprResult))
	}
	var callNode pkg_rhino.Node = exprResult.GetOnlyChild()
	_ = callNode
	if !(callNode.IsCall()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(callNode))
	}
	var originalCodeFunction pkg_rhino.Node = callNode.GetLastChild()
	_ = originalCodeFunction
	if !(originalCodeFunction.IsFunction()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(originalCodeFunction))
	}
	originalCodeFunction.Detach()
	n.ReplaceWith(originalCodeFunction)
}
func (*ManageClosureUnawareCode_UnwrapConcealedClosureUnawareCode) CanReleaseSyntacticScopes_() {}
