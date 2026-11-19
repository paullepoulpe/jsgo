/*
 * Copyright 2007 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/InvocationsCallback.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java polymorphic class interface: com.google.javascript.jscomp.InvocationsCallback
type InvocationsCallback interface {
	Implements_jscomp_InvocationsCallback()
	IsJavaNil_() bool
	GetInvocationsCallbackBase() *InvocationsCallbackBase
	NodeTraversal_AbstractPostOrderCallback
	Visit_4(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node, string)
}
func ToIface_InvocationsCallback[T interface {
	InvocationsCallback
	comparable
}](v T) InvocationsCallback {
	var z T
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.jscomp.InvocationsCallback
type InvocationsCallbackBase struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	self_ InvocationsCallback
}

func (this_ *InvocationsCallbackBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(InvocationsCallback)
}
func (this_ *InvocationsCallbackBase) IsJavaNil_() bool { return this_ == nil }
func (*InvocationsCallbackBase) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.InvocationsCallback", "InvocationsCallback" }
func (this_ *InvocationsCallbackBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InvocationsCallback", "com.google.javascript.jscomp.InvocationsCallback", func(v_ any) bool {
		return j2gort.InstanceOf[InvocationsCallback](v_)
	})
}
func (b *InvocationsCallbackBase) GetInvocationsCallbackBase() *InvocationsCallbackBase { return b }
func (*InvocationsCallbackBase) Implements_jscomp_InvocationsCallback() {}
func (*InvocationsCallbackBase) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*InvocationsCallbackBase) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *InvocationsCallbackBase) initFields_InvocationsCallback() {
	_ = this_
}

func InitInvocationsCallback(this_ *InvocationsCallbackBase) {
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_InvocationsCallback()
}

// Java method: InvocationsCallback.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *InvocationsCallbackBase) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if (!(n.IsCall())) {
		return
	}
	var callee pkg_rhino.Node = n.GetFirstChild()
	_ = callee
	if (!(callee.IsGetProp())) {
		return
	}
	this_.self_.Visit_4(t, n, parent, callee.GetString())
}

// Java method: InvocationsCallback.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node callNode, com.google.javascript.rhino.Node parent, java.lang.String callName)
func (this_ *InvocationsCallbackBase) Visit_4(t *NodeTraversal, callNode pkg_rhino.Node, parent pkg_rhino.Node, callName string) {
	_ = t
	_ = callNode
	_ = parent
	_ = callName
	panic("abstract or native method invoked")
}

