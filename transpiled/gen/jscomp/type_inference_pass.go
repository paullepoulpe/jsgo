/*
 * Copyright 2009 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/TypeInferencePass.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_diagnostic "github.com/paullepoulpe/jsgo/transpiled/gen/diagnostic"
	pkg_graph "github.com/paullepoulpe/jsgo/transpiled/gen/graph"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_jscomptype "github.com/paullepoulpe/jsgo/transpiled/gen/jscomptype"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

var lambdaSingleton_type_inference_pass_1 j2gort.Ext_java_util_function_ToIntFunction
var lambdaSingleton_type_inference_pass_2 j2gort.Ext_java_util_function_Function

func init() {
	lambdaSingleton_type_inference_pass_1 = j2gort.NewLambda_Ext_java_util_function_ToIntFunction(func(mrp0_ any) int32 {
		return j2gort.AsIface[j2gort.Ext_com_google_common_collect_Multiset_Entry](mrp0_).GetCount()
	})
	lambdaSingleton_type_inference_pass_2 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_15 any) any {
		k := j2gort.CastPtr[j2gort.BoxInt](lp_15)
		_ = k
		return j2gort.PtrToAny(j2gort.Ext_com_google_common_collect_HashMultiset_Create_0())
	})
}

// Java class: com.google.javascript.jscomp.TypeInferencePass
type TypeInferencePass struct {
	compiler                AbstractCompiler                          `gson:"compiler"`
	registry                *pkg_rhino.JSTypeRegistry                 `gson:"registry"`
	reverseInterpreter      pkg_jscomptype.ReverseAbstractInterpreter `gson:"reverseInterpreter"`
	topScope                *TypedScope                               `gson:"topScope"`
	scopeCreator            *TypedScopeCreator                        `gson:"scopeCreator"`
	assertionFunctionLookup *CodingConvention_AssertionFunctionLookup `gson:"assertionFunctionLookup"`
	stepCountHistogram      *j2gort.Ext_java_util_LinkedHashMap       `gson:"stepCountHistogram"`
	reusableTypeInference   *TypeInference                            `gson:"reusableTypeInference"`
	idHash_                 atomic.Int32
	typeInferenceInUse      bool `gson:"typeInferenceInUse"`
}

func (this_ *TypeInferencePass) BindSelf_(s any) {
	_ = s
}
func (this_ *TypeInferencePass) IsJavaNil_() bool { return this_ == nil }
func (*TypeInferencePass) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TypeInferencePass", "TypeInferencePass"
}
func (this_ *TypeInferencePass) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TypeInferencePass) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TypeInferencePass", "com.google.javascript.jscomp.TypeInferencePass", func(v_ any) bool {
		return j2gort.InstanceOf[*TypeInferencePass](v_)
	})
}

func (this_ *TypeInferencePass) initFields_TypeInferencePass() {
	_ = this_
}

// Java constructor: TypeInferencePass(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.type.ReverseAbstractInterpreter reverseInterpreter, com.google.javascript.jscomp.TypedScopeCreator scopeCreator)
func NewTypeInferencePass(compiler AbstractCompiler, reverseInterpreter pkg_jscomptype.ReverseAbstractInterpreter, scopeCreator *TypedScopeCreator) *TypeInferencePass {
	this_ := &TypeInferencePass{}
	InitTypeInferencePass(this_, compiler, reverseInterpreter, scopeCreator)
	return this_
}

func InitTypeInferencePass(this_ *TypeInferencePass, compiler AbstractCompiler, reverseInterpreter pkg_jscomptype.ReverseAbstractInterpreter, scopeCreator *TypedScopeCreator) {
	_ = compiler
	_ = reverseInterpreter
	_ = scopeCreator
	this_.initFields_TypeInferencePass()
	this_.compiler = compiler
	this_.registry = compiler.GetTypeRegistry()
	this_.reverseInterpreter = reverseInterpreter
	this_.scopeCreator = scopeCreator
	this_.assertionFunctionLookup = CodingConvention_AssertionFunctionLookup_Of(j2gort.AsIface[j2gort.Ext_java_util_Collection](compiler.GetCodingConvention().GetAssertionFunctions()))
	if compiler.IsDebugLoggingEnabled() {
		this_.stepCountHistogram = j2gort.NewExt_java_util_LinkedHashMap()
	} else {
		this_.stepCountHistogram = nil
	}
}

// Java method: TypeInferencePass.inferAllScopes(com.google.javascript.rhino.Node inferenceRoot)
func (this_ *TypeInferencePass) InferAllScopes(inferenceRoot pkg_rhino.Node) *TypedScope {
	_ = inferenceRoot
	j2gort.PurgeAllRuntimePools()

	func() {
		var closer *pkg_rhino.JSTypeResolver_Closer = this_.registry.GetResolver().OpenForDefinition()
		_ = closer
		defer j2gort.CloseResource(closer)
		if !(inferenceRoot.IsRoot()) {
			j2gort.CheckState(false)
		}
		if !(inferenceRoot.GetParent() == 0) {
			j2gort.CheckState(false)
		}
		if !(this_.topScope == nil) {
			j2gort.CheckState(false)
		}
		this_.topScope = j2gort.CastPtr[TypedScope](this_.scopeCreator.CreateScope_2_Node_AbstractScope(inferenceRoot, nil))
		{
			clinit_NodeTraversal()
			var ntb_ NodeTraversal_Builder
			InitNodeTraversal_Builder(&ntb_)
			(&ntb_).SetCompiler(this_.compiler)
			(&ntb_).SetCallback_1_NodeTraversal_Callback(NewTypeInferencePass_FirstScopeBuildingCallback())
			(&ntb_).SetScopeCreator(ToIface_ScopeCreator(this_.scopeCreator))
			(&ntb_).TraverseWithScope(inferenceRoot, ToIface_AbstractScope(this_.topScope))
		}
		this_.scopeCreator.ResolveWeakImportsPreResolution()
	}()
	this_.scopeCreator.FinishAndFreeze()
	{
		clinit_NodeTraversal()
		var ntb_1_ NodeTraversal_Builder
		InitNodeTraversal_Builder(&ntb_1_)
		(&ntb_1_).SetCompiler(this_.compiler)
		(&ntb_1_).SetCallback_1_NodeTraversal_Callback(NewTypeInferencePass_SecondScopeBuildingCallback(this_))
		(&ntb_1_).SetScopeCreator(ToIface_ScopeCreator(this_.scopeCreator))
		(&ntb_1_).TraverseWithScope(inferenceRoot, ToIface_AbstractScope(this_.topScope))
	}
	var unknownType pkg_rhino.JSType = this_.registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
	_ = unknownType
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.scopeCreator.GetAllSymbols())))
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			var_ := j2gort.CastPtr[TypedVar](rawElem_2)
			_ = var_
			if j2gort.IfaceEquals[pkg_rhino.JSType](var_.GetType(), nil) {
				var_.SetType(unknownType)
			}
		}
	}
	if this_.stepCountHistogram != nil {
		func() {
			var histogram pkg_diagnostic.LogFile = this_.compiler.CreateOrReopenLog(j2gort.GetObjectClass(this_), "step_histogram.log", []string{})
			_ = histogram
			defer j2gort.CloseResource(histogram)
			histogram.Log_1_string("step_count token population")
			var totals []int32 = []int32{int32(0), int32(0)}
			_ = totals
			j2gort.AsIface[j2gort.Ext_java_util_Set](this_.stepCountHistogram.KeySet_0()).Stream().Sorted_1(j2gort.Ext_java_util_Comparator_NaturalOrder().Reversed()).ForEach(j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_3 any) {
				stepCount := j2gort.CastPtr[j2gort.BoxInt](lp_3)
				_ = stepCount
				j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.CastPtr[j2gort.Ext_com_google_common_collect_HashMultiset](this_.stepCountHistogram.Get(j2gort.PtrToAny(stepCount))).EntrySet()).Stream().Sorted_1(j2gort.Ext_java_util_Comparator_ComparingInt(lambdaSingleton_type_inference_pass_1)).ForEach(j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_4 any) {
					e := j2gort.AsIface[j2gort.Ext_com_google_common_collect_Multiset_Entry](lp_4)
					_ = e
					a_5 := totals
					i_6 := int32(0)
					a_5[i_6] = (a_5[i_6] + (j2gort.UnboxInt(stepCount) * e.GetCount()))
					a_7 := totals
					i_8 := int32(1)
					a_7[i_8] = (a_7[i_8] + e.GetCount())
					histogram.Log_2("%s %s %s", []any{j2gort.PtrToAny(stepCount), j2gort.BoxNullableVal(j2gort.Cast[pkg_rhino.Token](e.GetElement())), j2gort.ValueOfInt(e.GetCount())})
				}))
			}))
			if histogram.IsLogging() {
				histogram.Log_2("%s TOTAL %s", []any{j2gort.ValueOfInt(totals[int32(0)]), j2gort.ValueOfInt(totals[int32(1)])})
			}
		}()
	}
	this_.reusableTypeInference = nil
	return this_.topScope
}

// Java method: TypeInferencePass.inferScope(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.TypedScope scope)
func (this_ *TypeInferencePass) priv_TypeInferencePass_inferScope_5579(n pkg_rhino.Node, scope *TypedScope) {
	_ = n
	_ = scope
	var cfg ControlFlowGraph[pkg_rhino.Node] = this_.priv_TypeInferencePass_computeCfg_5580(n)
	_ = cfg
	var usedReusable bool = false
	_ = usedReusable
	var typeInference *TypeInference = nil
	_ = typeInference
	if !(this_.typeInferenceInUse) {
		this_.typeInferenceInUse = true
		usedReusable = true
		typeInference = this_.reusableTypeInference
		if typeInference == nil {
			typeInference = NewTypeInference(this_.compiler, cfg, this_.reverseInterpreter, scope, this_.scopeCreator, this_.assertionFunctionLookup)
			this_.reusableTypeInference = typeInference
		} else {
			typeInference.ResetForScope(cfg, scope)
		}
	} else {
		typeInference = NewTypeInference(this_.compiler, cfg, this_.reverseInterpreter, scope, this_.scopeCreator, this_.assertionFunctionLookup)
	}
	defer func() {
		typeInference.FinishScope()
		if usedReusable {
			this_.typeInferenceInUse = false
		}
		ControlFlowAnalysis_ReleaseCfg(cfg)
	}()
	{
		typeInference.Analyze()
		if this_.stepCountHistogram != nil {
			{
				cur_9 := j2gort.NewIterCursor(j2gort.CheckNotNull(cfg.GetNodes()))
				for cur_9.HasNext() {
					rawElem_10 := cur_9.Next()
					node := j2gort.AsIface[pkg_graph.DiGraph_DiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch]](rawElem_10)
					_ = node
					if j2gort.IfaceEquals[pkg_graph.DiGraph_DiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch]](node, cfg.GetImplicitReturn()) {
						continue
					}
					var state *DataFlowAnalysis_LinearFlowState = j2gort.CastPtr[DataFlowAnalysis_LinearFlowState](node.GetAnnotation())
					_ = state
					func() *j2gort.Ext_com_google_common_collect_HashMultiset {
						m_11 := this_.stepCountHistogram
						k_12 := j2gort.ValueOfInt(state.GetStepCount())
						if v_13, ok_14 := j2gort.MapTryGetComputeAny(m_11, k_12); ok_14 {
							return j2gort.CastPtr[j2gort.Ext_com_google_common_collect_HashMultiset](v_13)
						}
						return j2gort.CastPtr[j2gort.Ext_com_google_common_collect_HashMultiset](m_11.ComputeIfAbsent(k_12, lambdaSingleton_type_inference_pass_2))
					}().Add_1(j2gort.BoxNullableVal(node.GetValue().GetToken()))
				}
			}
		}
	}
}

// Java method: TypeInferencePass.computeCfg(com.google.javascript.rhino.Node n)
func (this_ *TypeInferencePass) priv_TypeInferencePass_computeCfg_5580(n pkg_rhino.Node) ControlFlowGraph[pkg_rhino.Node] {
	_ = n
	return ControlFlowAnalysis_ComputeCfgDirect(this_.compiler, n, false, true)
}

// Java class: com.google.javascript.jscomp.TypeInferencePass.FirstScopeBuildingCallback
type TypeInferencePass_FirstScopeBuildingCallback struct {
	NodeTraversal_AbstractScopedCallbackBase
}

func (this_ *TypeInferencePass_FirstScopeBuildingCallback) BindSelf_(s any) {
	_ = s
}
func (this_ *TypeInferencePass_FirstScopeBuildingCallback) IsJavaNil_() bool { return this_ == nil }
func (*TypeInferencePass_FirstScopeBuildingCallback) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TypeInferencePass$FirstScopeBuildingCallback", "FirstScopeBuildingCallback"
}
func (this_ *TypeInferencePass_FirstScopeBuildingCallback) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractScopedCallbackBase.IdentityHash_()
}
func (*TypeInferencePass_FirstScopeBuildingCallback) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TypeInferencePass.FirstScopeBuildingCallback", "com.google.javascript.jscomp.TypeInferencePass$FirstScopeBuildingCallback", func(v_ any) bool {
		return j2gort.InstanceOf[*TypeInferencePass_FirstScopeBuildingCallback](v_)
	})
}
func (*TypeInferencePass_FirstScopeBuildingCallback) Implements_jscomp_NodeTraversal_AbstractScopedCallback() {
}
func (*TypeInferencePass_FirstScopeBuildingCallback) Implements_jscomp_NodeTraversal_ScopedCallback() {
}
func (*TypeInferencePass_FirstScopeBuildingCallback) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *TypeInferencePass_FirstScopeBuildingCallback) initFields_TypeInferencePass_FirstScopeBuildingCallback() {
	_ = this_
}

// Java constructor: FirstScopeBuildingCallback()
func NewTypeInferencePass_FirstScopeBuildingCallback() *TypeInferencePass_FirstScopeBuildingCallback {
	this_ := &TypeInferencePass_FirstScopeBuildingCallback{}
	InitTypeInferencePass_FirstScopeBuildingCallback(this_)
	return this_
}

func InitTypeInferencePass_FirstScopeBuildingCallback(this_ *TypeInferencePass_FirstScopeBuildingCallback) {
	InitNodeTraversal_AbstractScopedCallback(&this_.NodeTraversal_AbstractScopedCallbackBase)
	this_.initFields_TypeInferencePass_FirstScopeBuildingCallback()
}

// Java method: FirstScopeBuildingCallback.enterScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *TypeInferencePass_FirstScopeBuildingCallback) EnterScope(t *NodeTraversal) {
	_ = t
	t.GetTypedScope()
}

// Java method: FirstScopeBuildingCallback.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *TypeInferencePass_FirstScopeBuildingCallback) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
}

// Java class: com.google.javascript.jscomp.TypeInferencePass.SecondScopeBuildingCallback
type TypeInferencePass_SecondScopeBuildingCallback struct {
	NodeTraversal_AbstractScopedCallbackBase
	outer_ *TypeInferencePass
}

func (this_ *TypeInferencePass_SecondScopeBuildingCallback) BindSelf_(s any) {
	_ = s
}
func (this_ *TypeInferencePass_SecondScopeBuildingCallback) IsJavaNil_() bool { return this_ == nil }
func (*TypeInferencePass_SecondScopeBuildingCallback) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TypeInferencePass$SecondScopeBuildingCallback", "SecondScopeBuildingCallback"
}
func (this_ *TypeInferencePass_SecondScopeBuildingCallback) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractScopedCallbackBase.IdentityHash_()
}
func (*TypeInferencePass_SecondScopeBuildingCallback) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TypeInferencePass.SecondScopeBuildingCallback", "com.google.javascript.jscomp.TypeInferencePass$SecondScopeBuildingCallback", func(v_ any) bool {
		return j2gort.InstanceOf[*TypeInferencePass_SecondScopeBuildingCallback](v_)
	})
}
func (this_ *TypeInferencePass_SecondScopeBuildingCallback) SetOuter_(o *TypeInferencePass) {
	this_.outer_ = o
}
func (this_ *TypeInferencePass_SecondScopeBuildingCallback) GetOuter_() *TypeInferencePass {
	return this_.outer_
}
func (*TypeInferencePass_SecondScopeBuildingCallback) Implements_jscomp_NodeTraversal_AbstractScopedCallback() {
}
func (*TypeInferencePass_SecondScopeBuildingCallback) Implements_jscomp_NodeTraversal_ScopedCallback() {
}
func (*TypeInferencePass_SecondScopeBuildingCallback) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *TypeInferencePass_SecondScopeBuildingCallback) initFields_TypeInferencePass_SecondScopeBuildingCallback() {
	_ = this_
}

// Java constructor: SecondScopeBuildingCallback()
func NewTypeInferencePass_SecondScopeBuildingCallback(outer_ *TypeInferencePass) *TypeInferencePass_SecondScopeBuildingCallback {
	this_ := &TypeInferencePass_SecondScopeBuildingCallback{}
	this_.outer_ = outer_
	InitTypeInferencePass_SecondScopeBuildingCallback(this_)
	return this_
}

func InitTypeInferencePass_SecondScopeBuildingCallback(this_ *TypeInferencePass_SecondScopeBuildingCallback) {
	InitNodeTraversal_AbstractScopedCallback(&this_.NodeTraversal_AbstractScopedCallbackBase)
	this_.initFields_TypeInferencePass_SecondScopeBuildingCallback()
}

// Java method: SecondScopeBuildingCallback.enterScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *TypeInferencePass_SecondScopeBuildingCallback) EnterScope(t *NodeTraversal) {
	_ = t
	var scope *TypedScope = t.GetTypedScope()
	_ = scope
	if scope.IsCfgRootScope() && !(scope.IsModuleScope()) {
		this_.outer_.priv_TypeInferencePass_inferScope_5579(t.GetCurrentNode(), scope)
	}
}

// Java method: SecondScopeBuildingCallback.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *TypeInferencePass_SecondScopeBuildingCallback) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
}
func (*TypeInferencePass_FirstScopeBuildingCallback) AlwaysShouldTraverse_()      {}
func (*TypeInferencePass_FirstScopeBuildingCallback) NoOpVisit_()                 {}
func (*TypeInferencePass_FirstScopeBuildingCallback) CanReleaseSyntacticScopes_() {}
func (*TypeInferencePass_SecondScopeBuildingCallback) AlwaysShouldTraverse_()     {}
func (*TypeInferencePass_SecondScopeBuildingCallback) NoOpVisit_()                {}
