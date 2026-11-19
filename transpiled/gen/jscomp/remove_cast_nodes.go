/*
 * Copyright 2021 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/RemoveCastNodes.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.RemoveCastNodes
type RemoveCastNodes struct {
	compiler AbstractCompiler `gson:"compiler"`
	idHash_  atomic.Int32
}

func (this_ *RemoveCastNodes) BindSelf_(s any) {
	_ = s
}
func (this_ *RemoveCastNodes) IsJavaNil_() bool { return this_ == nil }
func (*RemoveCastNodes) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RemoveCastNodes", "RemoveCastNodes"
}
func (this_ *RemoveCastNodes) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RemoveCastNodes) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RemoveCastNodes", "com.google.javascript.jscomp.RemoveCastNodes", func(v_ any) bool {
		return j2gort.InstanceOf[*RemoveCastNodes](v_)
	})
}
func (*RemoveCastNodes) Implements_jscomp_CompilerPass() {}

func (this_ *RemoveCastNodes) initFields_RemoveCastNodes() {
	_ = this_
}

// Java constructor: RemoveCastNodes(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewRemoveCastNodes(compiler AbstractCompiler) *RemoveCastNodes {
	this_ := &RemoveCastNodes{}
	InitRemoveCastNodes(this_, compiler)
	return this_
}

func InitRemoveCastNodes(this_ *RemoveCastNodes, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_RemoveCastNodes()
	this_.compiler = compiler
}

// Java method: RemoveCastNodes.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *RemoveCastNodes) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_TraverseRoots(this_.compiler, NewRemoveCastNodes_RemoveCastsCallback(this_), externs, root)
}

// Java class: com.google.javascript.jscomp.RemoveCastNodes.RemoveCastsCallback
type RemoveCastNodes_RemoveCastsCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_ *RemoveCastNodes
}

func (this_ *RemoveCastNodes_RemoveCastsCallback) BindSelf_(s any) {
	_ = s
}
func (this_ *RemoveCastNodes_RemoveCastsCallback) IsJavaNil_() bool { return this_ == nil }
func (*RemoveCastNodes_RemoveCastsCallback) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RemoveCastNodes$RemoveCastsCallback", "RemoveCastsCallback"
}
func (this_ *RemoveCastNodes_RemoveCastsCallback) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*RemoveCastNodes_RemoveCastsCallback) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RemoveCastNodes.RemoveCastsCallback", "com.google.javascript.jscomp.RemoveCastNodes$RemoveCastsCallback", func(v_ any) bool {
		return j2gort.InstanceOf[*RemoveCastNodes_RemoveCastsCallback](v_)
	})
}
func (this_ *RemoveCastNodes_RemoveCastsCallback) SetOuter_(o *RemoveCastNodes) { this_.outer_ = o }
func (this_ *RemoveCastNodes_RemoveCastsCallback) GetOuter_() *RemoveCastNodes  { return this_.outer_ }
func (*RemoveCastNodes_RemoveCastsCallback) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*RemoveCastNodes_RemoveCastsCallback) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *RemoveCastNodes_RemoveCastsCallback) initFields_RemoveCastNodes_RemoveCastsCallback() {
	_ = this_
}

// Java constructor: RemoveCastsCallback()
func NewRemoveCastNodes_RemoveCastsCallback(outer_ *RemoveCastNodes) *RemoveCastNodes_RemoveCastsCallback {
	this_ := &RemoveCastNodes_RemoveCastsCallback{}
	this_.outer_ = outer_
	InitRemoveCastNodes_RemoveCastsCallback(this_)
	return this_
}

func InitRemoveCastNodes_RemoveCastsCallback(this_ *RemoveCastNodes_RemoveCastsCallback) {
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_RemoveCastNodes_RemoveCastsCallback()
}

// Java method: RemoveCastsCallback.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *RemoveCastNodes_RemoveCastsCallback) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsCast() {
		this_.outer_.compiler.ReportChangeToEnclosingScope(n)
		n.ReplaceWith(n.RemoveFirstChild())
	}
}
func (*RemoveCastNodes_RemoveCastsCallback) AlwaysShouldTraverse_()      {}
func (*RemoveCastNodes_RemoveCastsCallback) CanReleaseSyntacticScopes_() {}
