/*
 * Copyright 2010 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/PeepholeOptimizationsPass.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.PeepholeOptimizationsPass
type PeepholeOptimizationsPass struct {
	compiler              AbstractCompiler               `gson:"compiler"`
	passName              string                         `gson:"passName"`
	peepholeOptimizations []AbstractPeepholeOptimization `gson:"peepholeOptimizations"`
	idHash_               atomic.Int32
	retraverseOnChange    bool `gson:"retraverseOnChange"`
}

func (this_ *PeepholeOptimizationsPass) BindSelf_(s any) {
	_ = s
}
func (this_ *PeepholeOptimizationsPass) IsJavaNil_() bool { return this_ == nil }
func (*PeepholeOptimizationsPass) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PeepholeOptimizationsPass", "PeepholeOptimizationsPass"
}
func (this_ *PeepholeOptimizationsPass) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PeepholeOptimizationsPass) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PeepholeOptimizationsPass", "com.google.javascript.jscomp.PeepholeOptimizationsPass", func(v_ any) bool {
		return j2gort.InstanceOf[*PeepholeOptimizationsPass](v_)
	})
}
func (*PeepholeOptimizationsPass) Implements_jscomp_CompilerPass() {}

func (this_ *PeepholeOptimizationsPass) initFields_PeepholeOptimizationsPass() {
	_ = this_
	this_.passName = j2gort.NullString
}

// Java constructor: PeepholeOptimizationsPass(com.google.javascript.jscomp.AbstractCompiler compiler, java.lang.String passName, com.google.javascript.jscomp.AbstractPeepholeOptimization[] optimizations)
func NewPeepholeOptimizationsPass_3_AbstractCompiler_string_AbstractPeepholeOptimization(compiler AbstractCompiler, passName string, optimizations []AbstractPeepholeOptimization) *PeepholeOptimizationsPass {
	this_ := &PeepholeOptimizationsPass{}
	InitPeepholeOptimizationsPass_3_AbstractCompiler_string_AbstractPeepholeOptimization(this_, compiler, passName, optimizations)
	return this_
}

func InitPeepholeOptimizationsPass_3_AbstractCompiler_string_AbstractPeepholeOptimization(this_ *PeepholeOptimizationsPass, compiler AbstractCompiler, passName string, optimizations []AbstractPeepholeOptimization) {
	_ = compiler
	_ = passName
	_ = optimizations
	var iifeArg_1 []any
iifeSwitch_1:
	switch {
	default:

		arr_1 := optimizations
		if arr_1 == nil {
			iifeArg_1 = nil
			break iifeSwitch_1
		}
		out_2 := make([]any, len(arr_1))
		for i_, el_3 := range arr_1 {
			out_2[i_] = el_3
		}
		iifeArg_1 = out_2
		break iifeSwitch_1
	}
	InitPeepholeOptimizationsPass_3_AbstractCompiler_string_Ext_java_util_List(this_, compiler, passName, j2gort.Ext_java_util_Arrays_AsList(iifeArg_1))
}

// Java constructor: PeepholeOptimizationsPass(com.google.javascript.jscomp.AbstractCompiler compiler, java.lang.String passName, java.util.List<com.google.javascript.jscomp.AbstractPeepholeOptimization> optimizations)
func NewPeepholeOptimizationsPass_3_AbstractCompiler_string_Ext_java_util_List(compiler AbstractCompiler, passName string, optimizations j2gort.Ext_java_util_List) *PeepholeOptimizationsPass {
	this_ := &PeepholeOptimizationsPass{}
	InitPeepholeOptimizationsPass_3_AbstractCompiler_string_Ext_java_util_List(this_, compiler, passName, optimizations)
	return this_
}

func InitPeepholeOptimizationsPass_3_AbstractCompiler_string_Ext_java_util_List(this_ *PeepholeOptimizationsPass, compiler AbstractCompiler, passName string, optimizations j2gort.Ext_java_util_List) {
	_ = compiler
	_ = passName
	_ = optimizations
	this_.initFields_PeepholeOptimizationsPass()
	this_.compiler = compiler
	this_.passName = passName
	this_.peepholeOptimizations = func() []AbstractPeepholeOptimization {
		arr_7 := optimizations.ToArray_1_slice_any(func() []any {
			arr_4 := make([]AbstractPeepholeOptimization, int32(0))
			if arr_4 == nil {
				return nil
			}
			out_5 := make([]any, len(arr_4))
			for i_, el_6 := range arr_4 {
				out_5[i_] = el_6
			}
			return out_5
		}())
		if arr_7 == nil {
			return nil
		}
		out_8 := make([]AbstractPeepholeOptimization, len(arr_7))
		for i_, el_9 := range arr_7 {
			out_8[i_] = j2gort.AsIface[AbstractPeepholeOptimization](el_9)
		}
		return out_8
	}()
	this_.retraverseOnChange = true
}

// Java method: PeepholeOptimizationsPass.setRetraverseOnChange(boolean retraverse)
func (this_ *PeepholeOptimizationsPass) SetRetraverseOnChange(retraverse bool) {
	_ = retraverse
	this_.retraverseOnChange = retraverse
}

// Java method: PeepholeOptimizationsPass.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *PeepholeOptimizationsPass) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	this_.priv_PeepholeOptimizationsPass_beginTraversal_4238()
	{
		var changedScopeNodes j2gort.Ext_java_util_List = this_.compiler.GetChangeTracker().GetChangedScopeNodesForPass(this_.passName)
		_ = changedScopeNodes
		for ; (changedScopeNodes == nil) || !(changedScopeNodes.IsEmpty()); changedScopeNodes = this_.compiler.GetChangeTracker().GetChangedScopeNodesForPass(this_.passName) {
			if changedScopeNodes == nil {
				NodeTraversal_Traverse(this_.compiler, root, NewPeepholeOptimizationsPass_PeepCallback(this_))
			} else {
				NodeTraversal_TraverseScopeRoots(this_.compiler, changedScopeNodes, NewPeepholeOptimizationsPass_PeepCallback(this_), false)
			}
			if !(this_.retraverseOnChange) {
				break
			}
		}
	}
	this_.priv_PeepholeOptimizationsPass_endTraversal_4239()
}

// Java method: PeepholeOptimizationsPass.beginTraversal()
func (this_ *PeepholeOptimizationsPass) priv_PeepholeOptimizationsPass_beginTraversal_4238() {
	for _, el_10 := range j2gort.CheckNotNullSlice(this_.peepholeOptimizations) {
		optimization := el_10
		_ = optimization
		optimization.BeginTraversal(this_.compiler)
	}
}

// Java method: PeepholeOptimizationsPass.endTraversal()
func (this_ *PeepholeOptimizationsPass) priv_PeepholeOptimizationsPass_endTraversal_4239() {
	for _, el_11 := range j2gort.CheckNotNullSlice(this_.peepholeOptimizations) {
		optimization := el_11
		_ = optimization
		optimization.EndTraversal()
	}
}

// Java class: com.google.javascript.jscomp.PeepholeOptimizationsPass.PeepCallback
type PeepholeOptimizationsPass_PeepCallback struct {
	NodeTraversal_AbstractScopedCallbackBase
	outer_ *PeepholeOptimizationsPass
}

func (this_ *PeepholeOptimizationsPass_PeepCallback) BindSelf_(s any) {
	_ = s
}
func (this_ *PeepholeOptimizationsPass_PeepCallback) IsJavaNil_() bool { return this_ == nil }
func (*PeepholeOptimizationsPass_PeepCallback) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PeepholeOptimizationsPass$PeepCallback", "PeepCallback"
}
func (this_ *PeepholeOptimizationsPass_PeepCallback) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractScopedCallbackBase.IdentityHash_()
}
func (*PeepholeOptimizationsPass_PeepCallback) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PeepholeOptimizationsPass.PeepCallback", "com.google.javascript.jscomp.PeepholeOptimizationsPass$PeepCallback", func(v_ any) bool {
		return j2gort.InstanceOf[*PeepholeOptimizationsPass_PeepCallback](v_)
	})
}
func (this_ *PeepholeOptimizationsPass_PeepCallback) SetOuter_(o *PeepholeOptimizationsPass) {
	this_.outer_ = o
}
func (this_ *PeepholeOptimizationsPass_PeepCallback) GetOuter_() *PeepholeOptimizationsPass {
	return this_.outer_
}
func (*PeepholeOptimizationsPass_PeepCallback) Implements_jscomp_NodeTraversal_AbstractScopedCallback() {
}
func (*PeepholeOptimizationsPass_PeepCallback) Implements_jscomp_NodeTraversal_ScopedCallback() {}
func (*PeepholeOptimizationsPass_PeepCallback) Implements_jscomp_NodeTraversal_Callback()       {}

func (this_ *PeepholeOptimizationsPass_PeepCallback) initFields_PeepholeOptimizationsPass_PeepCallback() {
	_ = this_
}

// Java constructor: PeepCallback()
func NewPeepholeOptimizationsPass_PeepCallback(outer_ *PeepholeOptimizationsPass) *PeepholeOptimizationsPass_PeepCallback {
	this_ := &PeepholeOptimizationsPass_PeepCallback{}
	this_.outer_ = outer_
	InitPeepholeOptimizationsPass_PeepCallback(this_)
	return this_
}

func InitPeepholeOptimizationsPass_PeepCallback(this_ *PeepholeOptimizationsPass_PeepCallback) {
	InitNodeTraversal_AbstractScopedCallback(&this_.NodeTraversal_AbstractScopedCallbackBase)
	this_.initFields_PeepholeOptimizationsPass_PeepCallback()
}

// Java method: PeepCallback.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *PeepholeOptimizationsPass_PeepCallback) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	var currentNode pkg_rhino.Node = n
	_ = currentNode
	for _, el_12 := range j2gort.CheckNotNullSlice(this_.outer_.peepholeOptimizations) {
		optim := el_12
		_ = optim
		currentNode = optim.OptimizeSubtree(currentNode)
		if currentNode == 0 {
			return
		}
	}
}

// Java method: PeepCallback.enterScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *PeepholeOptimizationsPass_PeepCallback) EnterScope(t *NodeTraversal) {
	_ = t
}

// Java method: PeepCallback.exitScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *PeepholeOptimizationsPass_PeepCallback) ExitScope(t *NodeTraversal) {
	_ = t
	this_.priv_PeepCallback_updateFeatures_4240(t)
}

// Java method: PeepCallback.updateFeatures(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *PeepholeOptimizationsPass_PeepCallback) priv_PeepCallback_updateFeatures_4240(t *NodeTraversal) {
	_ = t
	for _, el_13 := range j2gort.CheckNotNullSlice(this_.outer_.peepholeOptimizations) {
		optim := el_13
		_ = optim
		var newFeatures j2gort.Ext_com_google_common_collect_ImmutableSet = optim.GetAbstractPeepholeOptimizationBase().GetNewFeatures()
		_ = newFeatures
		if !(newFeatures.IsEmpty()) {
			NodeUtil_AddFeaturesToScript(t.GetCurrentScript(), pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM().With_1_Ext_java_util_Set(j2gort.AsIface[j2gort.Ext_java_util_Set](newFeatures)), this_.outer_.compiler)
		}
		optim.GetAbstractPeepholeOptimizationBase().ClearNewFeatures()
	}
}
func (*PeepholeOptimizationsPass_PeepCallback) AlwaysShouldTraverse_()      {}
func (*PeepholeOptimizationsPass_PeepCallback) CanReleaseSyntacticScopes_() {}
