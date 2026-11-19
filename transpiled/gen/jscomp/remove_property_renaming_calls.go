/*
 * Copyright 2023 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/RemovePropertyRenamingCalls.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.RemovePropertyRenamingCalls
type RemovePropertyRenamingCalls struct {
	compiler AbstractCompiler `gson:"compiler"`
	idHash_  atomic.Int32
}

func (this_ *RemovePropertyRenamingCalls) BindSelf_(s any) {
	_ = s
}
func (this_ *RemovePropertyRenamingCalls) IsJavaNil_() bool { return this_ == nil }
func (*RemovePropertyRenamingCalls) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RemovePropertyRenamingCalls", "RemovePropertyRenamingCalls"
}
func (this_ *RemovePropertyRenamingCalls) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RemovePropertyRenamingCalls) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RemovePropertyRenamingCalls", "com.google.javascript.jscomp.RemovePropertyRenamingCalls", func(v_ any) bool {
		return j2gort.InstanceOf[*RemovePropertyRenamingCalls](v_)
	})
}
func (*RemovePropertyRenamingCalls) Implements_jscomp_CompilerPass() {}

func (this_ *RemovePropertyRenamingCalls) initFields_RemovePropertyRenamingCalls() {
	_ = this_
}

// Java constructor: RemovePropertyRenamingCalls(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewRemovePropertyRenamingCalls(compiler AbstractCompiler) *RemovePropertyRenamingCalls {
	this_ := &RemovePropertyRenamingCalls{}
	InitRemovePropertyRenamingCalls(this_, compiler)
	return this_
}

func InitRemovePropertyRenamingCalls(this_ *RemovePropertyRenamingCalls, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_RemovePropertyRenamingCalls()
	this_.compiler = compiler
}

// Java method: RemovePropertyRenamingCalls.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *RemovePropertyRenamingCalls) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, NewRemovePropertyRenamingCalls_RemoveCompilerPropertyRenamingCallback(this_))
}

// Java class: com.google.javascript.jscomp.RemovePropertyRenamingCalls.RemoveCompilerPropertyRenamingCallback
type RemovePropertyRenamingCalls_RemoveCompilerPropertyRenamingCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_ *RemovePropertyRenamingCalls
}

func (this_ *RemovePropertyRenamingCalls_RemoveCompilerPropertyRenamingCallback) BindSelf_(s any) {
	_ = s
}
func (this_ *RemovePropertyRenamingCalls_RemoveCompilerPropertyRenamingCallback) IsJavaNil_() bool {
	return this_ == nil
}
func (*RemovePropertyRenamingCalls_RemoveCompilerPropertyRenamingCallback) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RemovePropertyRenamingCalls$RemoveCompilerPropertyRenamingCallback", "RemoveCompilerPropertyRenamingCallback"
}
func (this_ *RemovePropertyRenamingCalls_RemoveCompilerPropertyRenamingCallback) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*RemovePropertyRenamingCalls_RemoveCompilerPropertyRenamingCallback) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RemovePropertyRenamingCalls.RemoveCompilerPropertyRenamingCallback", "com.google.javascript.jscomp.RemovePropertyRenamingCalls$RemoveCompilerPropertyRenamingCallback", func(v_ any) bool {
		return j2gort.InstanceOf[*RemovePropertyRenamingCalls_RemoveCompilerPropertyRenamingCallback](v_)
	})
}
func (this_ *RemovePropertyRenamingCalls_RemoveCompilerPropertyRenamingCallback) SetOuter_(o *RemovePropertyRenamingCalls) {
	this_.outer_ = o
}
func (this_ *RemovePropertyRenamingCalls_RemoveCompilerPropertyRenamingCallback) GetOuter_() *RemovePropertyRenamingCalls {
	return this_.outer_
}
func (*RemovePropertyRenamingCalls_RemoveCompilerPropertyRenamingCallback) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*RemovePropertyRenamingCalls_RemoveCompilerPropertyRenamingCallback) Implements_jscomp_NodeTraversal_Callback() {
}

func (this_ *RemovePropertyRenamingCalls_RemoveCompilerPropertyRenamingCallback) initFields_RemovePropertyRenamingCalls_RemoveCompilerPropertyRenamingCallback() {
	_ = this_
}

// Java constructor: RemoveCompilerPropertyRenamingCallback()
func NewRemovePropertyRenamingCalls_RemoveCompilerPropertyRenamingCallback(outer_ *RemovePropertyRenamingCalls) *RemovePropertyRenamingCalls_RemoveCompilerPropertyRenamingCallback {
	this_ := &RemovePropertyRenamingCalls_RemoveCompilerPropertyRenamingCallback{}
	this_.outer_ = outer_
	InitRemovePropertyRenamingCalls_RemoveCompilerPropertyRenamingCallback(this_)
	return this_
}

func InitRemovePropertyRenamingCalls_RemoveCompilerPropertyRenamingCallback(this_ *RemovePropertyRenamingCalls_RemoveCompilerPropertyRenamingCallback) {
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_RemovePropertyRenamingCalls_RemoveCompilerPropertyRenamingCallback()
}

// Java method: RemoveCompilerPropertyRenamingCallback.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *RemovePropertyRenamingCalls_RemoveCompilerPropertyRenamingCallback) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsCall() {
		var fn pkg_rhino.Node = n.GetFirstChild()
		_ = fn
		if this_.outer_.compiler.GetCodingConvention().IsPropertyRenameFunction(fn) {
			this_.priv_RemoveCompilerPropertyRenamingCallback_replacePropertyRenamingCall_4677(n)
		}
	}
}

// Java method: RemoveCompilerPropertyRenamingCallback.replacePropertyRenamingCall(com.google.javascript.rhino.Node callNode)
func (this_ *RemovePropertyRenamingCalls_RemoveCompilerPropertyRenamingCallback) priv_RemoveCompilerPropertyRenamingCallback_replacePropertyRenamingCall_4677(callNode pkg_rhino.Node) {
	_ = callNode
	var propName pkg_rhino.Node = NodeUtil_GetArgumentForCallOrNew(callNode, int32(0))
	_ = propName
	if propName != 0 {
		callNode.ReplaceWith(propName.Detach())
		this_.outer_.compiler.ReportChangeToEnclosingScope(propName)
		return
	}
}
func (*RemovePropertyRenamingCalls_RemoveCompilerPropertyRenamingCallback) AlwaysShouldTraverse_() {}
func (*RemovePropertyRenamingCalls_RemoveCompilerPropertyRenamingCallback) CanReleaseSyntacticScopes_() {
}
