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
// Port of com/google/javascript/jscomp/PeepholeTranspilationsPass.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.PeepholeTranspilationsPass
type PeepholeTranspilationsPass struct {
	compiler                         AbstractCompiler                `gson:"compiler"`
	peepholeTranspilations           []AbstractPeepholeTranspilation `gson:"peepholeTranspilations"`
	featureSetToMarkAsTranspiledAway *pkg_parser.FeatureSet          `gson:"featureSetToMarkAsTranspiledAway"`
	featureSetToRunOn                *pkg_parser.FeatureSet          `gson:"featureSetToRunOn"`
	idHash_                          atomic.Int32
}

func (this_ *PeepholeTranspilationsPass) BindSelf_(s any) {
	_ = s
}
func (this_ *PeepholeTranspilationsPass) IsJavaNil_() bool { return this_ == nil }
func (*PeepholeTranspilationsPass) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PeepholeTranspilationsPass", "PeepholeTranspilationsPass"
}
func (this_ *PeepholeTranspilationsPass) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PeepholeTranspilationsPass) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PeepholeTranspilationsPass", "com.google.javascript.jscomp.PeepholeTranspilationsPass", func(v_ any) bool {
		return j2gort.InstanceOf[*PeepholeTranspilationsPass](v_)
	})
}
func (*PeepholeTranspilationsPass) Implements_jscomp_CompilerPass() {}

func (this_ *PeepholeTranspilationsPass) initFields_PeepholeTranspilationsPass() {
	_ = this_
}

// Java constructor: PeepholeTranspilationsPass(com.google.javascript.jscomp.AbstractCompiler compiler, java.util.List<com.google.javascript.jscomp.AbstractPeepholeTranspilation> transpilations, com.google.javascript.jscomp.parsing.parser.FeatureSet featureSetToRunOn, com.google.javascript.jscomp.parsing.parser.FeatureSet featureSetToMarkAsTranspiledAway)
func NewPeepholeTranspilationsPass(compiler AbstractCompiler, transpilations j2gort.Ext_java_util_List, featureSetToRunOn *pkg_parser.FeatureSet, featureSetToMarkAsTranspiledAway *pkg_parser.FeatureSet) *PeepholeTranspilationsPass {
	this_ := &PeepholeTranspilationsPass{}
	InitPeepholeTranspilationsPass(this_, compiler, transpilations, featureSetToRunOn, featureSetToMarkAsTranspiledAway)
	return this_
}

func InitPeepholeTranspilationsPass(this_ *PeepholeTranspilationsPass, compiler AbstractCompiler, transpilations j2gort.Ext_java_util_List, featureSetToRunOn *pkg_parser.FeatureSet, featureSetToMarkAsTranspiledAway *pkg_parser.FeatureSet) {
	_ = compiler
	_ = transpilations
	_ = featureSetToRunOn
	_ = featureSetToMarkAsTranspiledAway
	this_.initFields_PeepholeTranspilationsPass()
	this_.compiler = compiler
	this_.peepholeTranspilations = func() []AbstractPeepholeTranspilation {
		arr_4 := transpilations.ToArray_1_slice_any(func() []any {
			arr_1 := make([]AbstractPeepholeTranspilation, int32(0))
			if arr_1 == nil {
				return nil
			}
			out_2 := make([]any, len(arr_1))
			for i_, el_3 := range arr_1 {
				out_2[i_] = el_3
			}
			return out_2
		}())
		if arr_4 == nil {
			return nil
		}
		out_5 := make([]AbstractPeepholeTranspilation, len(arr_4))
		for i_, el_6 := range arr_4 {
			out_5[i_] = j2gort.AsIface[AbstractPeepholeTranspilation](el_6)
		}
		return out_5
	}()
	this_.featureSetToRunOn = featureSetToRunOn
	this_.featureSetToMarkAsTranspiledAway = featureSetToMarkAsTranspiledAway
}

// Java static method: PeepholeTranspilationsPass.create(com.google.javascript.jscomp.AbstractCompiler compiler, java.util.List<com.google.javascript.jscomp.AbstractPeepholeTranspilation> transpilations)
func PeepholeTranspilationsPass_Create(compiler AbstractCompiler, transpilations j2gort.Ext_java_util_List) *PeepholeTranspilationsPass {
	_ = compiler
	_ = transpilations
	var featureSetToRunOn *pkg_parser.FeatureSet = pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM()
	_ = featureSetToRunOn
	var featureSetToMarkAsTranspiledAway *pkg_parser.FeatureSet = pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM()
	_ = featureSetToMarkAsTranspiledAway
	{
		cur_7 := j2gort.NewIterCursor(j2gort.CheckNotNull(transpilations))
		for cur_7.HasNext() {
			rawElem_8 := cur_7.Next()
			transpilation := j2gort.AsIface[AbstractPeepholeTranspilation](rawElem_8)
			_ = transpilation
			j2gort.CheckStateFmt(!(transpilation.GetAdditionalFeaturesToRunOn().ContainsAtLeastOneOf(transpilation.GetTranspiledAwayFeatures())), "Transpilation pass %s has getAdditionalFeaturesToRunOn() that contains features that it transpiles away.", j2gort.BoxString(j2gort.CheckNotNullPtr(j2gort.GetObjectClass(transpilation)).GetSimpleName()))
			featureSetToRunOn = featureSetToRunOn.With_1_FeatureSet(transpilation.GetTranspiledAwayFeatures())
			featureSetToRunOn = featureSetToRunOn.With_1_FeatureSet(transpilation.GetAdditionalFeaturesToRunOn())
			featureSetToMarkAsTranspiledAway = featureSetToMarkAsTranspiledAway.With_1_FeatureSet(transpilation.GetTranspiledAwayFeatures())
		}
	}
	return NewPeepholeTranspilationsPass(compiler, transpilations, featureSetToRunOn, featureSetToMarkAsTranspiledAway)
}

// Java method: PeepholeTranspilationsPass.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *PeepholeTranspilationsPass) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, NewPeepholeTranspilationsPass_PeepCallback(this_))
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(this_.compiler, root, this_.featureSetToMarkAsTranspiledAway)
}

// Java class: com.google.javascript.jscomp.PeepholeTranspilationsPass.PeepCallback
type PeepholeTranspilationsPass_PeepCallback struct {
	outer_  *PeepholeTranspilationsPass
	idHash_ atomic.Int32
}

func (this_ *PeepholeTranspilationsPass_PeepCallback) BindSelf_(s any) {
	_ = s
}
func (this_ *PeepholeTranspilationsPass_PeepCallback) IsJavaNil_() bool { return this_ == nil }
func (*PeepholeTranspilationsPass_PeepCallback) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PeepholeTranspilationsPass$PeepCallback", "PeepCallback"
}
func (this_ *PeepholeTranspilationsPass_PeepCallback) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PeepholeTranspilationsPass_PeepCallback) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PeepholeTranspilationsPass.PeepCallback", "com.google.javascript.jscomp.PeepholeTranspilationsPass$PeepCallback", func(v_ any) bool {
		return j2gort.InstanceOf[*PeepholeTranspilationsPass_PeepCallback](v_)
	})
}
func (this_ *PeepholeTranspilationsPass_PeepCallback) SetOuter_(o *PeepholeTranspilationsPass) {
	this_.outer_ = o
}
func (this_ *PeepholeTranspilationsPass_PeepCallback) GetOuter_() *PeepholeTranspilationsPass {
	return this_.outer_
}
func (*PeepholeTranspilationsPass_PeepCallback) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *PeepholeTranspilationsPass_PeepCallback) initFields_PeepholeTranspilationsPass_PeepCallback() {
	_ = this_
}

// Java constructor: PeepCallback()
func NewPeepholeTranspilationsPass_PeepCallback(outer_ *PeepholeTranspilationsPass) *PeepholeTranspilationsPass_PeepCallback {
	this_ := &PeepholeTranspilationsPass_PeepCallback{}
	this_.outer_ = outer_
	InitPeepholeTranspilationsPass_PeepCallback(this_)
	return this_
}

func InitPeepholeTranspilationsPass_PeepCallback(this_ *PeepholeTranspilationsPass_PeepCallback) {
	this_.initFields_PeepholeTranspilationsPass_PeepCallback()
}

// Java method: PeepCallback.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *PeepholeTranspilationsPass_PeepCallback) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	var currentNode pkg_rhino.Node = n
	_ = currentNode
	for _, el_9 := range j2gort.CheckNotNullSlice(this_.outer_.peepholeTranspilations) {
		transpilation := el_9
		_ = transpilation
		currentNode = transpilation.TranspileSubtree(currentNode)
		if currentNode == 0 {
			return
		}
	}
}

// Java method: PeepCallback.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *PeepholeTranspilationsPass_PeepCallback) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(94) /* Token.SCRIPT */ :
		var scriptFeatures *pkg_parser.FeatureSet = NodeUtil_GetFeatureSetOfScript(n)
		_ = scriptFeatures
		return scriptFeatures.ContainsAtLeastOneOf(this_.outer_.featureSetToRunOn)
	default:
		return true
	}
	panic("unreachable switch expression")
}
func (*PeepholeTranspilationsPass_PeepCallback) CanReleaseSyntacticScopes_() {}
