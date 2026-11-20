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
// Port of com/google/javascript/jscomp/CoalesceVariableNames.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_graph "github.com/paullepoulpe/jsgo/transpiled/gen/graph"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

var lambdaSingleton_coalesce_variable_names_1 j2gort.Ext_java_util_function_IntFunction
var lambdaSingleton_coalesce_variable_names_2 j2gort.Ext_com_google_common_base_Predicate
var lambdaSingleton_coalesce_variable_names_3 j2gort.Ext_com_google_common_base_Predicate
var lambdaSingleton_coalesce_variable_names_4 j2gort.Ext_com_google_common_base_Predicate

func init() {
	lambdaSingleton_coalesce_variable_names_1 = j2gort.NewLambda_Ext_java_util_function_IntFunction(func(i int32) any {
		_ = i
		return j2gort.NewExt_java_util_BitSet()
	})
	lambdaSingleton_coalesce_variable_names_2 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_ uint32) bool { return NodeUtil_IsNameDeclaration(pkg_rhino.Node(rawNode_)) })
	lambdaSingleton_coalesce_variable_names_3 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_ uint32) bool { return NodeUtil_IsNameDeclaration(pkg_rhino.Node(rawNode_)) })
	lambdaSingleton_coalesce_variable_names_4 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_ uint32) bool { return NodeUtil_IsNameDeclaration(pkg_rhino.Node(rawNode_)) })
}

// Java class: com.google.javascript.jscomp.CoalesceVariableNames
type CoalesceVariableNames struct {
	NodeTraversal_AbstractCfgCallbackBase
	compiler                 AbstractCompiler                                 `gson:"compiler"`
	scopeCreator             *MemoizedScopeCreator                            `gson:"scopeCreator"`
	colorings                j2gort.Deque[pkg_graph.GraphColoring[*Var, any]] `gson:"colorings"`
	liveAnalyses             j2gort.Deque[*LiveVariablesAnalysis]             `gson:"liveAnalyses"`
	freeLiveAnalyses         j2gort.Deque[*LiveVariablesAnalysis]             `gson:"freeLiveAnalyses"`
	astFactory               *AstFactory                                      `gson:"astFactory"`
	liveness                 *LiveVariablesAnalysis                           `gson:"liveness"`
	coloringTieBreaker       j2gort.Ext_java_util_Comparator                  `gson:"coloringTieBreaker"`
	shouldOptimizeScopeStack j2gort.Deque[bool]                               `gson:"shouldOptimizeScopeStack"`
	usePseudoNames           bool                                             `gson:"usePseudoNames"`
}

func (this_ *CoalesceVariableNames) BindSelf_(s any) {
	_ = s
	this_.NodeTraversal_AbstractCfgCallbackBase.BindSelf_(s)
}
func (this_ *CoalesceVariableNames) IsJavaNil_() bool { return this_ == nil }
func (*CoalesceVariableNames) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CoalesceVariableNames", "CoalesceVariableNames"
}
func (this_ *CoalesceVariableNames) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractCfgCallbackBase.IdentityHash_()
}
func (*CoalesceVariableNames) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CoalesceVariableNames", "com.google.javascript.jscomp.CoalesceVariableNames", func(v_ any) bool {
		return j2gort.InstanceOf[*CoalesceVariableNames](v_)
	})
}
func (*CoalesceVariableNames) Implements_jscomp_NodeTraversal_AbstractCfgCallback() {}
func (*CoalesceVariableNames) Implements_jscomp_CompilerPass()                      {}
func (*CoalesceVariableNames) Implements_jscomp_NodeTraversal_ScopedCallback()      {}
func (*CoalesceVariableNames) Implements_jscomp_NodeTraversal_Callback()            {}

func (this_ *CoalesceVariableNames) initFields_CoalesceVariableNames() {
	_ = this_
	this_.freeLiveAnalyses = j2gort.NewDequeCap[*LiveVariablesAnalysis](int32(8))
	this_.coloringTieBreaker = j2gort.Ext_java_util_Comparator_ComparingInt(j2gort.NewLambda_Ext_java_util_function_ToIntFunction(func(lp_1 any) int32 {
		arg := j2gort.CastPtr[Var](lp_1)
		_ = arg
		return this_.liveness.GetVarIndex(arg.GetName())
	}))
	this_.shouldOptimizeScopeStack = j2gort.Deque[bool]{}
}

// Java constructor: CoalesceVariableNames(com.google.javascript.jscomp.AbstractCompiler compiler, boolean usePseudoNames)
func NewCoalesceVariableNames(compiler AbstractCompiler, usePseudoNames bool) *CoalesceVariableNames {
	this_ := &CoalesceVariableNames{}
	this_.BindSelf_(this_)
	InitCoalesceVariableNames(this_, compiler, usePseudoNames)
	return this_
}

func InitCoalesceVariableNames(this_ *CoalesceVariableNames, compiler AbstractCompiler, usePseudoNames bool) {
	_ = compiler
	_ = usePseudoNames
	InitNodeTraversal_AbstractCfgCallback(&this_.NodeTraversal_AbstractCfgCallbackBase)
	this_.initFields_CoalesceVariableNames()
	if !(compiler.GetLifeCycleStage().IsNormalized()) {
		j2gort.CheckState(false)
	}
	this_.compiler = compiler
	this_.colorings = j2gort.Deque[pkg_graph.GraphColoring[*Var, any]]{}
	this_.liveAnalyses = j2gort.Deque[*LiveVariablesAnalysis]{}
	this_.usePseudoNames = usePseudoNames
	this_.astFactory = compiler.CreateAstFactory()
	this_.scopeCreator = NewMemoizedScopeCreator_2(NewSyntacticScopeCreator_1(compiler), max(int32(1024), (compiler.GetNumberOfInputs()*int32(2))))
}

// Java method: CoalesceVariableNames.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *CoalesceVariableNames) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	defer func() {
		for !this_.colorings.IsEmpty() {
			if c_ := this_.colorings.Pop(); c_ != nil {
				if r_, ok_ := any(c_.GetGraph()).(interface{ ReleaseSlabs_() }); ok_ {
					r_.ReleaseSlabs_()
				}
			}
		}

		for !this_.liveAnalyses.IsEmpty() {
			if lv_ := this_.liveAnalyses.Pop(); lv_ != nil {
				lv_.ReleaseSlabs_()
			}
		}
		this_.freeLiveAnalyses.Clear()
		this_.liveness = nil
		if this_.scopeCreator != nil {
			this_.scopeCreator.ReleaseSlabs_()
			this_.scopeCreator = nil
		}
	}()

	j2gort.CheckNotNullNode(externs)
	j2gort.CheckNotNullNode(root)
	{
		clinit_NodeTraversal()
		var ntb_ NodeTraversal_Builder
		InitNodeTraversal_Builder(&ntb_)
		(&ntb_).SetCompiler(this_.compiler)
		(&ntb_).SetCallback_1_NodeTraversal_Callback(this_)
		(&ntb_).SetScopeCreator(ToIface_ScopeCreator(this_.scopeCreator))
		(&ntb_).Traverse(root)
	}
	this_.compiler.SetLifeCycleStage(StaticGet_AbstractCompiler_LifeCycleStage_RAW())
}

// Java static method: CoalesceVariableNames.shouldOptimizeScope(com.google.javascript.jscomp.NodeTraversal t)
func CoalesceVariableNames_ShouldOptimizeScope(t *NodeTraversal) *NodeUtil_AllVarsDeclaredInFunction {
	_ = t
	if t.GetScopeRoot().IsFunction() {
		var allVarsDeclaredInFunction *NodeUtil_AllVarsDeclaredInFunction = NodeUtil_GetAllVarsDeclaredInFunction(t.GetCompiler(), t.GetScopeCreator(), t.GetScope())
		_ = allVarsDeclaredInFunction
		if int32(100) > allVarsDeclaredInFunction.GetAllVariablesInOrder().Size() {
			return allVarsDeclaredInFunction
		}
	}
	return nil
}

// Java method: CoalesceVariableNames.enterScopeWithCfg(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *CoalesceVariableNames) EnterScopeWithCfg(t *NodeTraversal) {
	_ = t
	var allVarsDeclaredInFunction *NodeUtil_AllVarsDeclaredInFunction = CoalesceVariableNames_ShouldOptimizeScope(t)
	_ = allVarsDeclaredInFunction
	if allVarsDeclaredInFunction == nil {
		this_.shouldOptimizeScopeStack.Push(false)
		return
	}
	this_.shouldOptimizeScopeStack.Push(true)
	var scope *Scope = t.GetScope()
	_ = scope
	if !(scope.IsFunctionScope()) {
		j2gort.CheckStateMsg(false, j2gort.PtrToAny(scope))
	}
	var cfg ControlFlowGraph[pkg_rhino.Node] = this_.GetControlFlowGraph(this_.compiler)
	_ = cfg
	var pooled *LiveVariablesAnalysis = this_.freeLiveAnalyses.PollLast()
	_ = pooled
	if pooled != nil {
		pooled.ResetForScope(cfg, scope, nil, this_.compiler, ToIface_ScopeCreator(this_.scopeCreator), allVarsDeclaredInFunction)
		this_.liveness = pooled
	} else {
		this_.liveness = NewLiveVariablesAnalysis(cfg, scope, nil, this_.compiler, ToIface_ScopeCreator(this_.scopeCreator), allVarsDeclaredInFunction)
	}
	if pkg_parser.StaticGet_FeatureSet_ES3().Contains_1_FeatureSet(this_.compiler.GetOptions().GetOutputFeatureSet()) {
		var enclosingFunction pkg_rhino.Node = scope.GetRootNode_0()
		_ = enclosingFunction
		if NodeUtil_GetFunctionParameters(enclosingFunction).HasTwoChildren() {
			this_.liveness.MarkAllParametersEscaped()
		}
	}
	this_.liveness.Analyze()
	this_.liveAnalyses.Push(this_.liveness)
	var interferenceGraph pkg_graph.UndiGraph[*Var, any] = this_.priv_CoalesceVariableNames_computeVariableNamesInterferenceGraph_1255(cfg, this_.liveness.GetEscapedLocals())
	_ = interferenceGraph
	var coloring pkg_graph.GraphColoring[*Var, any] = pkg_graph.ToIface_GraphColoring[*Var, any](pkg_graph.NewGraphColoring_GreedyGraphColoring_2[*Var, any](j2gort.AsIface[pkg_graph.AdjacencyGraph[*Var, any]](interferenceGraph), this_.coloringTieBreaker))
	_ = coloring
	coloring.Color()
	this_.colorings.Push(coloring)
}

// Java method: CoalesceVariableNames.exitScopeWithCfg(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *CoalesceVariableNames) ExitScopeWithCfg(t *NodeTraversal) {
	_ = t
	if !(this_.shouldOptimizeScopeStack.Pop()) {
		return
	}
	poppedColoring_ := this_.colorings.Pop()
	if poppedColoring_ != nil && !poppedColoring_.IsJavaNil_() {
		if r_, ok_ := any(poppedColoring_.GetGraph()).(interface{ ReleaseSlabs_() }); ok_ {
			r_.ReleaseSlabs_()
		}
	}

	var finished *LiveVariablesAnalysis = this_.liveAnalyses.Pop()
	_ = finished
	finished.ClearScopeState()
	if this_.freeLiveAnalyses.Size() < int32(8) {
		this_.freeLiveAnalyses.AddLast(finished)
	}
	this_.liveness = this_.liveAnalyses.Peek()
}

// Java method: CoalesceVariableNames.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CoalesceVariableNames) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if (this_.colorings.IsEmpty() || !(n.IsName())) || parent.IsFunction() {
		return
	}
	var var_ *Var = j2gort.CastPtr[Var](j2gort.MapGetStr(this_.liveness.GetAllVariables(), n.GetString()))
	_ = var_
	var vNode pkg_graph.GraphNode[*Var, any] = j2gort.AsIface[pkg_graph.GraphNode[*Var, any]](this_.colorings.Peek().GetGraph().GetNode(var_))
	_ = vNode
	if vNode == nil || vNode.IsJavaNil_() {
		return
	}
	var coalescedVar *Var = this_.colorings.Peek().GetPartitionSuperNode(var_)
	_ = coalescedVar
	if !(this_.usePseudoNames) {
		if vNode.GetValue().Equals(j2gort.PtrToAny(coalescedVar)) {
			return
		}
		n.SetString(coalescedVar.GetName())
		this_.compiler.ReportChangeToEnclosingScope(n)
		this_.priv_CoalesceVariableNames_updateDeclarationsPostCoalescing_1253(n, coalescedVar, var_, parent)
	} else {
		var allMergedNames j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_TreeSet()
		_ = allMergedNames
		{
			cur_2 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.liveness.GetAllVariablesInOrder()))
			for cur_2.HasNext() {
				rawElem_3 := cur_2.Next()
				iVar := j2gort.CastPtr[Var](rawElem_3)
				_ = iVar
				if !(j2gort.IfaceEquals[pkg_graph.GraphNode[*Var, any]](j2gort.AsIface[pkg_graph.GraphNode[*Var, any]](this_.colorings.Peek().GetGraph().GetNode(iVar)), nil)) && this_.colorings.Peek().HaveSameColor(iVar, coalescedVar) {
					j2gort.SetAddStr(allMergedNames, iVar.GetName())
				}
			}
		}
		if allMergedNames.Size() == int32(1) {
			return
		}
		var pseudoName string = j2gort.Ext_com_google_common_base_Joiner_On_1_string("_").Join_1_Ext_java_lang_Iterable(allMergedNames)
		_ = pseudoName
		for t.GetScope().HasSlot(pseudoName) {
			pseudoName = (j2gort.StrValString(pseudoName) + "$")
		}
		n.SetString(pseudoName)
		this_.compiler.ReportChangeToEnclosingScope(n)
		if vNode.GetValue().Equals(j2gort.PtrToAny(coalescedVar)) {
			return
		}
		this_.priv_CoalesceVariableNames_updateDeclarationsPostCoalescing_1253(n, coalescedVar, var_, parent)
	}
}

// Java method: CoalesceVariableNames.updateDeclarationsPostCoalescing(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.Var coalescedVar, com.google.javascript.jscomp.Var originalVar, com.google.javascript.rhino.Node parent)
func (this_ *CoalesceVariableNames) priv_CoalesceVariableNames_updateDeclarationsPostCoalescing_1253(n pkg_rhino.Node, coalescedVar *Var, originalVar *Var, parent pkg_rhino.Node) {
	_ = n
	_ = coalescedVar
	_ = originalVar
	_ = parent
	if !(n.IsName()) {
		j2gort.CheckStateStr(false, "trying to update the declaration of a non-name node")
	}
	if NodeUtil_IsNameDeclaration(parent) || (((NodeUtil_GetEnclosingType(n, pkg_rhino.StaticGet_Token_DESTRUCTURING_LHS()) != 0) && NodeUtil_IsLhsByDestructuring(n)) && !(this_.priv_CoalesceVariableNames_isNameInsideDestructuringAssignment_1254(n))) {
		this_.priv_CoalesceVariableNames_makeDeclarationVar_1258(coalescedVar)
		this_.priv_CoalesceVariableNames_removeVarDeclaration_1257(n, originalVar)
	}
}

// Java method: CoalesceVariableNames.isNameInsideDestructuringAssignment(com.google.javascript.rhino.Node n)
func (this_ *CoalesceVariableNames) priv_CoalesceVariableNames_isNameInsideDestructuringAssignment_1254(n pkg_rhino.Node) bool {
	_ = n
	var enclosingTarget pkg_rhino.Node = NodeUtil_GetRootTarget(n)
	_ = enclosingTarget
	if enclosingTarget.IsDestructuringPattern() && enclosingTarget.GetParent().IsAssign() {
		return true
	}
	return false
}

// Java method: CoalesceVariableNames.computeVariableNamesInterferenceGraph(com.google.javascript.jscomp.ControlFlowGraph<com.google.javascript.rhino.Node> cfg, java.util.Set<? extends com.google.javascript.jscomp.Var> escaped)
func (this_ *CoalesceVariableNames) priv_CoalesceVariableNames_computeVariableNamesInterferenceGraph_1255(cfg ControlFlowGraph[pkg_rhino.Node], escaped j2gort.Ext_java_util_Set) pkg_graph.UndiGraph[*Var, any] {
	_ = cfg
	_ = escaped
	var interferenceGraph pkg_graph.UndiGraph[*Var, any] = pkg_graph.ToIface_UndiGraph[*Var, any](pkg_graph.LinkedUndirectedGraph_Create[*Var, any]())
	_ = interferenceGraph
	var orderedVariables []*Var = j2gort.ColToNewPtrArray[Var](this_.liveness.GetAllVariablesInOrder())
	_ = orderedVariables
	var interferenceGraphNodes_bsVal j2gort.Ext_java_util_BitSet
	j2gort.InitExt_java_util_BitSet_0(&interferenceGraphNodes_bsVal)
	var interferenceGraphNodes *j2gort.Ext_java_util_BitSet = &interferenceGraphNodes_bsVal
	_ = interferenceGraphNodes
	var interferenceBitSet []*j2gort.Ext_java_util_BitSet = make([]*j2gort.Ext_java_util_BitSet, int32(len(orderedVariables)))
	_ = interferenceBitSet
	func() {
		arrMut_4 := interferenceBitSet
		var arrBox_5 []any
		if arrMut_4 != nil {
			arrBox_5 = make([]any, len(arrMut_4))
			for i_, el_6 := range arrMut_4 {
				arrBox_5[i_] = j2gort.PtrToAny(el_6)
			}
		}
		j2gort.Ext_java_util_Arrays_SetAll_2_slice_any_Ext_java_util_function_IntFunction(arrBox_5, lambdaSingleton_coalesce_variable_names_1)
		for i_, el_6 := range arrBox_5 {
			arrMut_4[i_] = j2gort.CastPtr[j2gort.Ext_java_util_BitSet](el_6)
		}
	}()
	var vIndex int32 = int32(-1)
	_ = vIndex
	for _, el_7 := range j2gort.CheckNotNullSlice(orderedVariables) {
		v := el_7
		_ = v
		vIndex = vIndex + 1
		if escaped.Contains(j2gort.PtrToAny(v)) {
			continue
		}
		if v.GetParentNode().IsFunction() {
			continue
		}
		if v.GetParentNode().IsClass() {
			continue
		}
		if this_.priv_CoalesceVariableNames_isInMultipleLvalueDecl_1256(v) {
			continue
		}
		interferenceGraph.CreateNode(v)
		interferenceGraphNodes.Set_1(vIndex)
	}
	{
		cur_8 := j2gort.NewIterCursor(j2gort.CheckNotNull(cfg.GetNodes()))
		for cur_8.HasNext() {
			rawElem_9 := cur_8.Next()
			cfgNode := j2gort.AsIface[pkg_graph.DiGraph_DiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch]](rawElem_9)
			_ = cfgNode
			if cfg.IsImplicitReturn(cfgNode) {
				continue
			}
			var state *DataFlowAnalysis_LinearFlowState = j2gort.CastPtr[DataFlowAnalysis_LinearFlowState](cfgNode.GetAnnotation())
			_ = state
			var livein *LiveVariablesAnalysis_LiveVariableLattice = j2gort.CastPtr[LiveVariablesAnalysis_LiveVariableLattice](state.GetIn())
			_ = livein
			{
				var i int32 = livein.NextSetBit(int32(0))
				_ = i
				for ; i >= int32(0); i = livein.NextSetBit((i + int32(1))) {
					{
						var j int32 = livein.NextSetBit(i)
						_ = j
						for ; j >= int32(0); j = livein.NextSetBit((j + int32(1))) {
							interferenceBitSet[i].Set_1(j)
						}
					}
				}
			}
			var liveout *LiveVariablesAnalysis_LiveVariableLattice = j2gort.CastPtr[LiveVariablesAnalysis_LiveVariableLattice](state.GetOut())
			_ = liveout
			{
				var i_2 int32 = liveout.NextSetBit(int32(0))
				_ = i_2
				for ; i_2 >= int32(0); i_2 = liveout.NextSetBit((i_2 + int32(1))) {
					{
						var j_2 int32 = liveout.NextSetBit(i_2)
						_ = j_2
						for ; j_2 >= int32(0); j_2 = liveout.NextSetBit((j_2 + int32(1))) {
							interferenceBitSet[i_2].Set_1(j_2)
						}
					}
				}
			}
			var liveRangeChecker *CoalesceVariableNames_LiveRangeChecker = NewCoalesceVariableNames_LiveRangeChecker(cfgNode.GetValue(), orderedVariables, state)
			_ = liveRangeChecker
			liveRangeChecker.Check(cfgNode.GetValue())
			liveRangeChecker.SetCrossingVariables(interferenceBitSet)
		}
	}
	var v1Index int32 = int32(-1)
	_ = v1Index
	for _, el_10 := range j2gort.CheckNotNullSlice(orderedVariables) {
		v1 := el_10
		_ = v1
		v1Index = v1Index + 1
		var v2Index int32 = int32(-1)
		_ = v2Index
		for _, el_11 := range j2gort.CheckNotNullSlice(orderedVariables) {
			v2 := el_11
			_ = v2
			v2Index = v2Index + 1
			if v1Index > v2Index {
				continue
			}
			if !(interferenceGraphNodes.Get_1(v1Index)) || !(interferenceGraphNodes.Get_1(v2Index)) {
				continue
			}
			if (v1.IsParam() && v2.IsParam()) || interferenceBitSet[v1Index].Get_1(v2Index) {
				interferenceGraph.ConnectIfNotFound(v1, nil, v2)
			}
		}
	}
	return interferenceGraph
}

// Java method: CoalesceVariableNames.isInMultipleLvalueDecl(com.google.javascript.jscomp.Var v)
func (this_ *CoalesceVariableNames) priv_CoalesceVariableNames_isInMultipleLvalueDecl_1256(v *Var) bool {
	_ = v
	var declarationType pkg_rhino.Token = v.DeclarationType()
	_ = declarationType
	switch (declarationType).Ordinal() {
	case int32(111) /* Token.LET */, int32(97) /* Token.CONST */, int32(84) /* Token.VAR */ :
		var nameDecl pkg_rhino.Node = NodeUtil_GetEnclosingNode(v.GetNode(), lambdaSingleton_coalesce_variable_names_2)
		_ = nameDecl
		var count []int32 = []int32{int32(0)}
		_ = count
		NodeUtil_VisitLhsNodesInNode(nameDecl, j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_12 any) {
			lhs := j2gort.CastNode[pkg_rhino.Node](lp_12)
			_ = lhs
			a_13 := count
			i_14 := int32(0)
			a_13[i_14] = a_13[i_14] + 1
		}))
		return (count[int32(0)] > int32(1))
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: CoalesceVariableNames.removeVarDeclaration(com.google.javascript.rhino.Node name, com.google.javascript.jscomp.Var originalVar)
func (this_ *CoalesceVariableNames) priv_CoalesceVariableNames_removeVarDeclaration_1257(name pkg_rhino.Node, originalVar *Var) {
	_ = name
	_ = originalVar
	var var_ pkg_rhino.Node = NodeUtil_GetEnclosingNode(name, lambdaSingleton_coalesce_variable_names_3)
	_ = var_
	var parent pkg_rhino.Node = var_.GetParent()
	_ = parent
	if var_.GetFirstChild().IsDestructuringLhs() {
		var destructuringLhs pkg_rhino.Node = var_.GetFirstChild()
		_ = destructuringLhs
		var pattern pkg_rhino.Node = destructuringLhs.RemoveFirstChild()
		_ = pattern
		if NodeUtil_IsEnhancedFor(parent) {
			var_.ReplaceWith(pattern)
		} else {
			var rvalue pkg_rhino.Node = var_.GetFirstFirstChild().Detach()
			_ = rvalue
			var_.ReplaceWith(NodeUtil_NewExpr(pkg_rhino.IR_Assign(pattern, rvalue).Srcref(var_)))
		}
	} else if NodeUtil_IsEnhancedFor(parent) {
		var_.ReplaceWith(name.Detach())
	} else {
		if !(var_.HasOneChild() && (var_.GetFirstChild() == name)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(var_))
		}
		if name.HasChildren() {
			var value pkg_rhino.Node = name.RemoveFirstChild()
			_ = value
			name.Detach()
			var assign pkg_rhino.Node = pkg_rhino.IR_Assign(name, value).Srcref(name)
			_ = assign
			if !(parent.IsVanillaFor()) {
				assign = NodeUtil_NewExpr(assign)
			}
			var_.ReplaceWith(assign)
		} else if this_.priv_CoalesceVariableNames_isUninitializedLetNameInLoopBodyAndLive_1259(name, originalVar) {
			name.Detach()
			var undefinedValue pkg_rhino.Node = this_.astFactory.CreateUndefinedValue().SrcrefTree(name)
			_ = undefinedValue
			var assign_2 pkg_rhino.Node = pkg_rhino.IR_Assign(name, undefinedValue).Srcref(name)
			_ = assign_2
			if !(parent.IsVanillaFor()) {
				assign_2 = NodeUtil_NewExpr(assign_2)
			}
			var_.ReplaceWith(assign_2)
		} else {
			NodeUtil_RemoveChild(parent, var_)
		}
	}
}

// Java method: CoalesceVariableNames.makeDeclarationVar(com.google.javascript.jscomp.Var coalescedName)
func (this_ *CoalesceVariableNames) priv_CoalesceVariableNames_makeDeclarationVar_1258(coalescedName *Var) {
	_ = coalescedName
	if coalescedName.IsConst() || coalescedName.IsLet() {
		var nameNode pkg_rhino.Node = j2gort.CheckNotNullNode(coalescedName.GetNameNode())
		_ = nameNode
		if CoalesceVariableNames_IsUninitializedLetNameInLoopBody(nameNode) {
			var undefinedValue pkg_rhino.Node = this_.astFactory.CreateUndefinedValue().SrcrefTree(nameNode)
			_ = undefinedValue
			nameNode.AddChildToFront(undefinedValue)
		}
		var declNode pkg_rhino.Node = NodeUtil_GetEnclosingNode(nameNode.GetParent(), lambdaSingleton_coalesce_variable_names_4)
		_ = declNode
		declNode.SetToken(pkg_rhino.StaticGet_Token_VAR())
	}
}

// Java static method: CoalesceVariableNames.isUninitializedLetNameInLoopBody(com.google.javascript.rhino.Node nameNode)
func CoalesceVariableNames_IsUninitializedLetNameInLoopBody(nameNode pkg_rhino.Node) bool {
	_ = nameNode
	if !(nameNode.IsName()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(nameNode))
	}
	var letNode pkg_rhino.Node = nameNode.GetParent()
	_ = letNode
	if !(letNode.IsLet()) {
		return false
	}
	if nameNode.HasOneChild() {
		return false
	}
	var letParent pkg_rhino.Node = letNode.GetParent()
	_ = letParent
	if NodeUtil_IsLoopStructure(letParent) {
		return false
	}
	return NodeUtil_IsWithinLoop(letParent)
}

// Java method: CoalesceVariableNames.isUninitializedLetNameInLoopBodyAndLive(com.google.javascript.rhino.Node nameNode, com.google.javascript.jscomp.Var originalVar)
func (this_ *CoalesceVariableNames) priv_CoalesceVariableNames_isUninitializedLetNameInLoopBodyAndLive_1259(nameNode pkg_rhino.Node, originalVar *Var) bool {
	_ = nameNode
	_ = originalVar
	if !(CoalesceVariableNames_IsUninitializedLetNameInLoopBody(nameNode)) {
		return false
	}
	if !(this_.liveness != nil) {
		j2gort.CheckState(false)
	}
	var letNode pkg_rhino.Node = nameNode.GetParent()
	_ = letNode
	var cfg ControlFlowGraph[pkg_rhino.Node] = this_.liveness.GetCfg()
	_ = cfg
	var cfgNode pkg_graph.DiGraph_DiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch] = pkg_graph.ToIface_DiGraph_DiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch](j2gort.CheckNotNullPtr(j2gort.CastPtr[pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch]](cfg.GetNode(letNode))))
	_ = cfgNode
	var state *DataFlowAnalysis_LinearFlowState = j2gort.CheckNotNullPtr(j2gort.CastPtr[DataFlowAnalysis_LinearFlowState](cfgNode.GetAnnotation()))
	_ = state
	var varIndex int32 = this_.liveness.GetVarIndex(originalVar.GetName())
	_ = varIndex
	return ((varIndex >= int32(0)) && j2gort.CastPtr[LiveVariablesAnalysis_LiveVariableLattice](state.GetOut()).IsLive(varIndex))
}

// Java class: com.google.javascript.jscomp.CoalesceVariableNames.LiveRangeChecker
type CoalesceVariableNames_LiveRangeChecker struct {
	state            *DataFlowAnalysis_LinearFlowState `gson:"state"`
	orderedVariables []*Var                            `gson:"orderedVariables"`
	isAssignToList   []int32                           `gson:"isAssignToList"`
	isReadFromList   []int32                           `gson:"isReadFromList"`
	root             pkg_rhino.Node                    `gson:"root"`
	idHash_          atomic.Int32
}

func (this_ *CoalesceVariableNames_LiveRangeChecker) BindSelf_(s any) {
	_ = s
}
func (this_ *CoalesceVariableNames_LiveRangeChecker) IsJavaNil_() bool { return this_ == nil }
func (*CoalesceVariableNames_LiveRangeChecker) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CoalesceVariableNames$LiveRangeChecker", "LiveRangeChecker"
}
func (this_ *CoalesceVariableNames_LiveRangeChecker) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CoalesceVariableNames_LiveRangeChecker) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CoalesceVariableNames.LiveRangeChecker", "com.google.javascript.jscomp.CoalesceVariableNames$LiveRangeChecker", func(v_ any) bool {
		return j2gort.InstanceOf[*CoalesceVariableNames_LiveRangeChecker](v_)
	})
}

func (this_ *CoalesceVariableNames_LiveRangeChecker) initFields_CoalesceVariableNames_LiveRangeChecker() {
	_ = this_
	this_.isAssignToList = []int32{}
	this_.isReadFromList = []int32{}
}

// Java constructor: LiveRangeChecker(com.google.javascript.rhino.Node root, com.google.javascript.jscomp.Var[] orderedVariables, com.google.javascript.jscomp.DataFlowAnalysis.LinearFlowState<com.google.javascript.jscomp.LiveVariablesAnalysis.LiveVariableLattice> state)
func NewCoalesceVariableNames_LiveRangeChecker(root pkg_rhino.Node, orderedVariables []*Var, state *DataFlowAnalysis_LinearFlowState) *CoalesceVariableNames_LiveRangeChecker {
	this_ := &CoalesceVariableNames_LiveRangeChecker{}
	InitCoalesceVariableNames_LiveRangeChecker(this_, root, orderedVariables, state)
	return this_
}

func InitCoalesceVariableNames_LiveRangeChecker(this_ *CoalesceVariableNames_LiveRangeChecker, root pkg_rhino.Node, orderedVariables []*Var, state *DataFlowAnalysis_LinearFlowState) {
	_ = root
	_ = orderedVariables
	_ = state
	this_.initFields_CoalesceVariableNames_LiveRangeChecker()
	this_.root = root
	this_.orderedVariables = orderedVariables
	this_.state = state
}

// Java method: LiveRangeChecker.check(com.google.javascript.rhino.Node n)
func (this_ *CoalesceVariableNames_LiveRangeChecker) Check(n pkg_rhino.Node) {
	_ = n
	if (n == this_.root) || !(ControlFlowGraph_IsEnteringNewCfgNode(n)) {
		if ((n.IsDestructuringLhs() && n.HasTwoChildren()) || (n.IsAssign() && n.GetFirstChild().IsDestructuringPattern())) || n.IsDefaultValue() {
			this_.Check(n.GetSecondChild())
			this_.Check(n.GetFirstChild())
		} else {
			{
				var c pkg_rhino.Node = n.GetFirstChild()
				_ = c
				for ; c != 0; c = c.GetNext() {
					this_.Check(c)
				}
			}
		}
		if CoalesceVariableNames_LiveRangeChecker_ShouldVisit(n) {
			this_.Visit(n, n.GetParent())
		}
	}
}

// Java method: LiveRangeChecker.visit(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CoalesceVariableNames_LiveRangeChecker) Visit(n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = n
	_ = parent
	{
		var iVar int32 = int32(0)
		_ = iVar
		for ; iVar < int32(len(this_.orderedVariables)); iVar = iVar + 1 {
			if CoalesceVariableNames_LiveRangeChecker_IsAssignTo(this_.orderedVariables[iVar], n, parent) {
				this_.isAssignToList = append(this_.isAssignToList, iVar)
			}
		}
	}
	if len(this_.isAssignToList) != 0 {
		{
			var iVar_2 int32 = int32(0)
			_ = iVar_2
			for ; iVar_2 < int32(len(this_.orderedVariables)); iVar_2 = iVar_2 + 1 {
				var varOutLive bool = j2gort.CastPtr[LiveVariablesAnalysis_LiveVariableLattice](this_.state.GetOut()).IsLive(iVar_2)
				_ = varOutLive
				if varOutLive || CoalesceVariableNames_LiveRangeChecker_IsReadFrom(this_.orderedVariables[iVar_2], n) {
					this_.isReadFromList = append(this_.isReadFromList, iVar_2)
				}
			}
		}
	}
}

// Java method: LiveRangeChecker.setCrossingVariables(java.util.BitSet[] interferenceBitSet)
func (this_ *CoalesceVariableNames_LiveRangeChecker) SetCrossingVariables(interferenceBitSet []*j2gort.Ext_java_util_BitSet) {
	_ = interferenceBitSet
	{
		var i int32 = int32(0)
		_ = i
		var wSize int32 = int32(len(this_.isAssignToList))
		_ = wSize
		for ; i < wSize; i = i + 1 {
			var iWrittenVar int32 = this_.isAssignToList[i]
			_ = iWrittenVar
			{
				var j int32 = int32(0)
				_ = j
				var rSize int32 = int32(len(this_.isReadFromList))
				_ = rSize
				for ; j < rSize; j = j + 1 {
					var iReadVar int32 = this_.isReadFromList[j]
					_ = iReadVar
					interferenceBitSet[iWrittenVar].Set_1(iReadVar)
					interferenceBitSet[iReadVar].Set_1(iWrittenVar)
				}
			}
		}
	}
}

// Java static method: LiveRangeChecker.shouldVisit(com.google.javascript.rhino.Node n)
func CoalesceVariableNames_LiveRangeChecker_ShouldVisit(n pkg_rhino.Node) bool {
	_ = n
	return (n.IsName() || (n.HasChildren() && n.GetFirstChild().IsName()))
}

// Java static method: LiveRangeChecker.isAssignTo(com.google.javascript.jscomp.Var var, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func CoalesceVariableNames_LiveRangeChecker_IsAssignTo(var_ *Var, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = var_
	_ = n
	_ = parent
	if n.IsName() {
		if parent.IsParamList() {
			{

				sr_15 := var_.GetName()
				se_16 := n.GetString()
				return j2gort.CheckNotNullString(sr_15) == se_16
			}

		} else if NodeUtil_IsNameDeclaration(parent) && n.HasChildren() {
			{

				sr_17 := var_.GetName()
				se_18 := n.GetString()
				return j2gort.CheckNotNullString(sr_17) == se_18
			}

		} else if NodeUtil_IsLhsByDestructuring(n) {
			{

				sr_19 := var_.GetName()
				se_20 := n.GetString()
				return j2gort.CheckNotNullString(sr_19) == se_20
			}

		}
	} else if NodeUtil_IsAssignmentOp(n) {
		var name pkg_rhino.Node = n.GetFirstChild()
		_ = name
		return (name.IsName() && func() bool {
			sr_21 := var_.GetName()
			se_22 := name.GetString()
			return j2gort.CheckNotNullString(sr_21) == se_22
		}())
	}
	return false
}

// Java static method: LiveRangeChecker.isReadFrom(com.google.javascript.jscomp.Var var, com.google.javascript.rhino.Node name)
func CoalesceVariableNames_LiveRangeChecker_IsReadFrom(var_ *Var, name pkg_rhino.Node) bool {
	_ = var_
	_ = name
	return ((name.IsName() && func() bool {
		sr_23 := var_.GetName()
		se_24 := name.GetString()
		return j2gort.CheckNotNullString(sr_23) == se_24
	}()) && !(NodeUtil_IsNameDeclOrSimpleAssignLhs(name, name.GetParent())))
}
func (*CoalesceVariableNames) AlwaysShouldTraverse_() {}
