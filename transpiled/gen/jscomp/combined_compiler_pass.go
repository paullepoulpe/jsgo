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
// Port of com/google/javascript/jscomp/CombinedCompilerPass.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.CombinedCompilerPass
type CombinedCompilerPass struct {
	callbacks []*CombinedCompilerPass_CallbackWrapper `gson:"callbacks"`
	compiler  AbstractCompiler                        `gson:"compiler"`
	idHash_   atomic.Int32
}

func (this_ *CombinedCompilerPass) BindSelf_(s any) {
	_ = s
}
func (this_ *CombinedCompilerPass) IsJavaNil_() bool { return this_ == nil }
func (*CombinedCompilerPass) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CombinedCompilerPass", "CombinedCompilerPass"
}
func (this_ *CombinedCompilerPass) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CombinedCompilerPass) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CombinedCompilerPass", "com.google.javascript.jscomp.CombinedCompilerPass", func(v_ any) bool {
		return j2gort.InstanceOf[*CombinedCompilerPass](v_)
	})
}
func (*CombinedCompilerPass) Implements_jscomp_CompilerPass()                 {}
func (*CombinedCompilerPass) Implements_jscomp_NodeTraversal_ScopedCallback() {}
func (*CombinedCompilerPass) Implements_jscomp_NodeTraversal_Callback()       {}

func (this_ *CombinedCompilerPass) initFields_CombinedCompilerPass() {
	_ = this_
}

// Java constructor: CombinedCompilerPass(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.NodeTraversal.Callback[] callbacks)
func NewCombinedCompilerPass_2_AbstractCompiler_NodeTraversal_Callback(compiler AbstractCompiler, callbacks []NodeTraversal_Callback) *CombinedCompilerPass {
	this_ := &CombinedCompilerPass{}
	InitCombinedCompilerPass_2_AbstractCompiler_NodeTraversal_Callback(this_, compiler, callbacks)
	return this_
}

func InitCombinedCompilerPass_2_AbstractCompiler_NodeTraversal_Callback(this_ *CombinedCompilerPass, compiler AbstractCompiler, callbacks []NodeTraversal_Callback) {
	_ = compiler
	_ = callbacks
	InitCombinedCompilerPass_2_AbstractCompiler_Ext_java_util_List(this_, compiler, j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_slice_any(func() []any {
		arr_1 := callbacks
		if arr_1 == nil {
			return nil
		}
		out_2 := make([]any, len(arr_1))
		for i_, el_3 := range arr_1 {
			out_2[i_] = el_3
		}
		return out_2
	}())))
}

// Java constructor: CombinedCompilerPass(com.google.javascript.jscomp.AbstractCompiler compiler, java.util.List<com.google.javascript.jscomp.NodeTraversal.Callback> callbacks)
func NewCombinedCompilerPass_2_AbstractCompiler_Ext_java_util_List(compiler AbstractCompiler, callbacks j2gort.Ext_java_util_List) *CombinedCompilerPass {
	this_ := &CombinedCompilerPass{}
	InitCombinedCompilerPass_2_AbstractCompiler_Ext_java_util_List(this_, compiler, callbacks)
	return this_
}

func InitCombinedCompilerPass_2_AbstractCompiler_Ext_java_util_List(this_ *CombinedCompilerPass, compiler AbstractCompiler, callbacks j2gort.Ext_java_util_List) {
	_ = compiler
	_ = callbacks
	this_.initFields_CombinedCompilerPass()
	this_.compiler = compiler
	this_.callbacks = make([]*CombinedCompilerPass_CallbackWrapper, callbacks.Size())
	{
		var i int32 = int32(0)
		_ = i
		for ; i < callbacks.Size(); i = i + 1 {
			{
				a_4 := this_.callbacks
				i_5 := i
				a_4[i_5] = NewCombinedCompilerPass_CallbackWrapper(j2gort.AsIface[NodeTraversal_Callback](callbacks.Get(i)))
			}
		}
	}
}

// Java static method: CombinedCompilerPass.traverse(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node root, java.util.List<com.google.javascript.jscomp.NodeTraversal.Callback> callbacks)
func CombinedCompilerPass_Traverse(compiler AbstractCompiler, root pkg_rhino.Node, callbacks j2gort.Ext_java_util_List) {
	_ = compiler
	_ = root
	_ = callbacks
	if callbacks.Size() == int32(1) {
		NodeTraversal_Traverse(compiler, root, j2gort.AsIface[NodeTraversal_Callback](callbacks.Get(int32(0))))
	} else {
		{
			var val_6 CombinedCompilerPass
			InitCombinedCompilerPass_2_AbstractCompiler_Ext_java_util_List(&val_6, compiler, callbacks)
			(&val_6).Process_2(pkg_rhino.Node(0), root)
		}
	}
}

// Java method: CombinedCompilerPass.process(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node externs, com.google.javascript.rhino.Node root)
func (this_ *CombinedCompilerPass) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
}

// Java method: CombinedCompilerPass.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CombinedCompilerPass) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	if this_.compiler.HasHaltingErrors() {
		return false
	}
	for _, el_7 := range j2gort.CheckNotNullSlice(this_.callbacks) {
		callback := el_7
		_ = callback
		callback.ShouldTraverseIfActive(t, n, parent)
	}
	return true
}

// Java method: CombinedCompilerPass.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CombinedCompilerPass) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if this_.compiler.HasHaltingErrors() {
		return
	}
	for _, el_8 := range j2gort.CheckNotNullSlice(this_.callbacks) {
		callback := el_8
		_ = callback
		callback.VisitOrMaybeActivate(t, n, parent)
	}
}

// Java method: CombinedCompilerPass.enterScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *CombinedCompilerPass) EnterScope(t *NodeTraversal) {
	_ = t
	for _, el_9 := range j2gort.CheckNotNullSlice(this_.callbacks) {
		callback := el_9
		_ = callback
		callback.EnterScopeIfActive(t)
	}
}

// Java method: CombinedCompilerPass.exitScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *CombinedCompilerPass) ExitScope(t *NodeTraversal) {
	_ = t
	for _, el_10 := range j2gort.CheckNotNullSlice(this_.callbacks) {
		callback := el_10
		_ = callback
		callback.ExitScopeIfActive(t)
	}
}

// Java class: com.google.javascript.jscomp.CombinedCompilerPass.CallbackWrapper
type CombinedCompilerPass_CallbackWrapper struct {
	callback       NodeTraversal_Callback       `gson:"callback"`
	scopedCallback NodeTraversal_ScopedCallback `gson:"scopedCallback"`
	waiting        pkg_rhino.Node               `gson:"waiting"`
	idHash_        atomic.Int32
}

func (this_ *CombinedCompilerPass_CallbackWrapper) BindSelf_(s any) {
	_ = s
}
func (this_ *CombinedCompilerPass_CallbackWrapper) IsJavaNil_() bool { return this_ == nil }
func (*CombinedCompilerPass_CallbackWrapper) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CombinedCompilerPass$CallbackWrapper", "CallbackWrapper"
}
func (this_ *CombinedCompilerPass_CallbackWrapper) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CombinedCompilerPass_CallbackWrapper) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CombinedCompilerPass.CallbackWrapper", "com.google.javascript.jscomp.CombinedCompilerPass$CallbackWrapper", func(v_ any) bool {
		return j2gort.InstanceOf[*CombinedCompilerPass_CallbackWrapper](v_)
	})
}

func (this_ *CombinedCompilerPass_CallbackWrapper) initFields_CombinedCompilerPass_CallbackWrapper() {
	_ = this_

}

// Java constructor: CallbackWrapper(com.google.javascript.jscomp.NodeTraversal.Callback callback)
func NewCombinedCompilerPass_CallbackWrapper(callback NodeTraversal_Callback) *CombinedCompilerPass_CallbackWrapper {
	this_ := &CombinedCompilerPass_CallbackWrapper{}
	InitCombinedCompilerPass_CallbackWrapper(this_, callback)
	return this_
}

func InitCombinedCompilerPass_CallbackWrapper(this_ *CombinedCompilerPass_CallbackWrapper, callback NodeTraversal_Callback) {
	_ = callback
	this_.initFields_CombinedCompilerPass_CallbackWrapper()
	this_.callback = callback
	var scopedCallback NodeTraversal_ScopedCallback
	_ = scopedCallback
	var ok_11 bool
	scopedCallback, ok_11 = j2gort.TryCast[NodeTraversal_ScopedCallback](callback)
	if ok_11 {
		this_.scopedCallback = scopedCallback
	} else {
		this_.scopedCallback = nil
	}
}

// Java method: CallbackWrapper.visitOrMaybeActivate(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CombinedCompilerPass_CallbackWrapper) VisitOrMaybeActivate(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if this_.IsActive() {
		this_.callback.Visit_3(t, n, parent)
	} else if this_.waiting == n {
		this_.waiting = pkg_rhino.Node(0)
	}
}

// Java method: CallbackWrapper.shouldTraverseIfActive(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CombinedCompilerPass_CallbackWrapper) ShouldTraverseIfActive(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if this_.IsActive() && !(this_.callback.ShouldTraverse_3(t, n, parent)) {
		this_.waiting = n
	}
}

// Java method: CallbackWrapper.enterScopeIfActive(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *CombinedCompilerPass_CallbackWrapper) EnterScopeIfActive(t *NodeTraversal) {
	_ = t
	if this_.IsActive() && (this_.scopedCallback != nil && !this_.scopedCallback.IsJavaNil_()) {
		this_.scopedCallback.EnterScope(t)
	}
}

// Java method: CallbackWrapper.exitScopeIfActive(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *CombinedCompilerPass_CallbackWrapper) ExitScopeIfActive(t *NodeTraversal) {
	_ = t
	if this_.IsActive() && (this_.scopedCallback != nil && !this_.scopedCallback.IsJavaNil_()) {
		this_.scopedCallback.ExitScope(t)
	}
}

// Java method: CallbackWrapper.isActive()
func (this_ *CombinedCompilerPass_CallbackWrapper) IsActive() bool {
	return (this_.waiting == 0)
}
func (this_ *CombinedCompilerPass) CanReleaseSyntacticScopesDynamic_() bool {
	if this_ == nil || len(this_.callbacks) == 0 {
		return false
	}
	for _, w_ := range this_.callbacks {
		if w_ == nil || !canReleaseSyntacticScopesForCallback_(w_.callback) {
			return false
		}
	}
	return true
}
