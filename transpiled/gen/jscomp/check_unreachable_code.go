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
// Port of com/google/javascript/jscomp/CheckUnreachableCode.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	pkg_graph "github.com/paullepoulpe/jsgo/transpiled/gen/graph"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.CheckUnreachableCode
type CheckUnreachableCode struct {
	NodeTraversal_AbstractCfgCallbackBase
	compiler  AbstractCompiler                                                                            `gson:"compiler"`
	workQueue []*pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch] `gson:"workQueue"`
}

func (this_ *CheckUnreachableCode) BindSelf_(s any) {
	_ = s
	this_.NodeTraversal_AbstractCfgCallbackBase.BindSelf_(s)
}
func (this_ *CheckUnreachableCode) IsJavaNil_() bool { return this_ == nil }
func (*CheckUnreachableCode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CheckUnreachableCode", "CheckUnreachableCode"
}
func (this_ *CheckUnreachableCode) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractCfgCallbackBase.IdentityHash_()
}
func (*CheckUnreachableCode) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CheckUnreachableCode", "com.google.javascript.jscomp.CheckUnreachableCode", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckUnreachableCode](v_)
	})
}
func (*CheckUnreachableCode) Implements_jscomp_NodeTraversal_AbstractCfgCallback() {}
func (*CheckUnreachableCode) Implements_jscomp_NodeTraversal_ScopedCallback()      {}
func (*CheckUnreachableCode) Implements_jscomp_NodeTraversal_Callback()            {}

var static_CheckUnreachableCode_UNREACHABLE_CODE *DiagnosticType

func StaticGet_CheckUnreachableCode_UNREACHABLE_CODE() *DiagnosticType {
	clinit_CheckUnreachableCode()
	return static_CheckUnreachableCode_UNREACHABLE_CODE
}
func StaticSet_CheckUnreachableCode_UNREACHABLE_CODE(v *DiagnosticType) {
	clinit_CheckUnreachableCode()
	static_CheckUnreachableCode_UNREACHABLE_CODE = v
}

var clinitDone_CheckUnreachableCode uint32
var clinitState_CheckUnreachableCode j2gort.ClinitState

func clinit_CheckUnreachableCode() {
	if atomic.LoadUint32(&clinitDone_CheckUnreachableCode) != 0 {
		return
	}
	clinitSlow_CheckUnreachableCode()
}

func clinitSlow_CheckUnreachableCode() {
	if !j2gort.EnterClinit(&clinitState_CheckUnreachableCode, &clinitDone_CheckUnreachableCode) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckUnreachableCode, &clinitDone_CheckUnreachableCode)
	static_CheckUnreachableCode_UNREACHABLE_CODE = DiagnosticType_Warning("JSC_UNREACHABLE_CODE", "unreachable code")
}

func (this_ *CheckUnreachableCode) initFields_CheckUnreachableCode() {
	_ = this_
	this_.workQueue = func() []*pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch] {
		arr_1 := make([]*pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[any, any], int32(32))
		if arr_1 == nil {
			return nil
		}
		out_2 := make([]*pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch], len(arr_1))
		for i_, el_3 := range arr_1 {
			out_2[i_] = j2gort.CastPtr[pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch]](el_3)
		}
		return out_2
	}()
}

// Java constructor: CheckUnreachableCode(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewCheckUnreachableCode(compiler AbstractCompiler) *CheckUnreachableCode {
	clinit_CheckUnreachableCode()
	this_ := &CheckUnreachableCode{}
	this_.BindSelf_(this_)
	InitCheckUnreachableCode(this_, compiler)
	return this_
}

func InitCheckUnreachableCode(this_ *CheckUnreachableCode, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractCfgCallback(&this_.NodeTraversal_AbstractCfgCallbackBase)
	this_.initFields_CheckUnreachableCode()
	this_.compiler = compiler
}

// Java method: CheckUnreachableCode.enterScopeWithCfg(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *CheckUnreachableCode) EnterScopeWithCfg(t *NodeTraversal) {
	_ = t
	if NodeUtil_IsValidCfgRoot(t.GetScopeRoot()) {
		this_.priv_CheckUnreachableCode_initScope_1128(this_.GetControlFlowGraph(this_.compiler))
	}
}

// Java method: CheckUnreachableCode.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckUnreachableCode) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	var cfg ControlFlowGraph[pkg_rhino.Node] = this_.GetControlFlowGraph(this_.compiler)
	_ = cfg
	var gNode *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch] = j2gort.CastPtr[pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch]](cfg.GetNode(n))
	_ = gNode
	if (gNode != nil) && !(j2gort.IfaceEquals[pkg_graph.Annotation](gNode.GetAnnotation(), pkg_graph.StaticGet_GraphReachability_REACHABLE())) {
		if ((n.GetLineno() != int32(-1)) && !(n.IsEmpty())) && !(n.IsBreak()) {
			this_.compiler.Report_1(JSError_Make_3(n, StaticGet_CheckUnreachableCode_UNREACHABLE_CODE(), []string{}))
			if !(!(j2gort.IfaceEquals[pkg_graph.Annotation](gNode.GetAnnotation(), pkg_graph.StaticGet_GraphReachability_REACHABLE()))) {
				j2gort.CheckState(false)
			}
			gNode.SetAnnotation(pkg_graph.StaticGet_GraphReachability_REACHABLE())
			this_.priv_CheckUnreachableCode_propagateReachability_1129(gNode, false)
			return false
		}
	}
	return true
}

// Java method: CheckUnreachableCode.initScope(com.google.javascript.jscomp.ControlFlowGraph<com.google.javascript.rhino.Node> controlFlowGraph)
func (this_ *CheckUnreachableCode) priv_CheckUnreachableCode_initScope_1128(controlFlowGraph ControlFlowGraph[pkg_rhino.Node]) {
	_ = controlFlowGraph
	controlFlowGraph.ClearNodeAnnotations()
	var entryNode *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch] = controlFlowGraph.GetEntryNode()
	_ = entryNode
	entryNode.SetAnnotation(pkg_graph.StaticGet_GraphReachability_REACHABLE())
	this_.priv_CheckUnreachableCode_propagateReachability_1129(entryNode, true)
}

// Java method: CheckUnreachableCode.propagateReachability(com.google.javascript.jscomp.graph.LinkedDirectedGraph.LinkedDiGraphNode<com.google.javascript.rhino.Node,com.google.javascript.jscomp.ControlFlowGraph.Branch> startNode, boolean checkConditional)
func (this_ *CheckUnreachableCode) priv_CheckUnreachableCode_propagateReachability_1129(startNode *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch], checkConditional bool) {
	_ = startNode
	_ = checkConditional
	var queue []*pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch] = this_.workQueue
	_ = queue
	{
		a_4 := queue
		i_5 := int32(0)
		a_4[i_5] = startNode
	}
	var head int32 = int32(0)
	_ = head
	var tail int32 = int32(1)
	_ = tail
	defer func() {
		j2gort.ArraysFillRangeZero(queue, head, tail)
	}()
	{
		if !(checkConditional) {
			for head < tail {
				var srcNode *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch] = queue[head]
				_ = srcNode
				{
					a_6 := queue
					var i_7 int32
					{

						old_8 := head
						nw_9 := old_8 + 1
						head = nw_9
						i_7 = old_8
					}
					a_6[i_7] = nil
				}
				{
					var i int32 = int32(0)
					_ = i
					var outCount int32 = srcNode.GetOutEdgeCount()
					_ = outCount
					for ; i < outCount; i = i + 1 {
						var destNode *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch] = j2gort.CastPtr[pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch]](j2gort.CastPtr[pkg_graph.LinkedDirectedGraph_LinkedDiGraphEdge[pkg_rhino.Node, ControlFlowGraph_Branch]](srcNode.GetOutEdgeAt(i)).GetDestination())
						_ = destNode
						if !(j2gort.IfaceEquals[pkg_graph.Annotation](destNode.GetAnnotation(), pkg_graph.StaticGet_GraphReachability_REACHABLE())) {
							destNode.SetAnnotation(pkg_graph.StaticGet_GraphReachability_REACHABLE())
							if tail == int32(len(queue)) {
								queue = this_.priv_CheckUnreachableCode_growQueue_1130(queue, head, tail)
								tail = (tail - head)
								head = int32(0)
							}
							{
								a_10 := queue
								var i_11 int32
								{

									old_12 := tail
									nw_13 := old_12 + 1
									tail = nw_13
									i_11 = old_12
								}
								a_10[i_11] = destNode
							}
						}
					}
				}
			}
		} else {
			for head < tail {
				var srcNode_2 *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch] = queue[head]
				_ = srcNode_2
				{
					a_14 := queue
					var i_15 int32
					{

						old_16 := head
						nw_17 := old_16 + 1
						head = nw_17
						i_15 = old_16
					}
					a_14[i_15] = nil
				}
				var srcValue pkg_rhino.Node = srcNode_2.GetValue()
				_ = srcValue
				{
					var i_2 int32 = int32(0)
					_ = i_2
					var outCount_2 int32 = srcNode_2.GetOutEdgeCount()
					_ = outCount_2
					for ; i_2 < outCount_2; i_2 = i_2 + 1 {
						var edge *pkg_graph.LinkedDirectedGraph_LinkedDiGraphEdge[pkg_rhino.Node, ControlFlowGraph_Branch] = j2gort.CastPtr[pkg_graph.LinkedDirectedGraph_LinkedDiGraphEdge[pkg_rhino.Node, ControlFlowGraph_Branch]](srcNode_2.GetOutEdgeAt(i_2))
						_ = edge
						var destNode_2 *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch] = j2gort.CastPtr[pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch]](edge.GetDestination())
						_ = destNode_2
						if !(j2gort.IfaceEquals[pkg_graph.Annotation](destNode_2.GetAnnotation(), pkg_graph.StaticGet_GraphReachability_REACHABLE())) && CheckUnreachableCode_IsReachable(srcValue, edge.GetValue()) {
							destNode_2.SetAnnotation(pkg_graph.StaticGet_GraphReachability_REACHABLE())
							if tail == int32(len(queue)) {
								queue = this_.priv_CheckUnreachableCode_growQueue_1130(queue, head, tail)
								tail = (tail - head)
								head = int32(0)
							}
							{
								a_18 := queue
								var i_19 int32
								{

									old_20 := tail
									nw_21 := old_20 + 1
									tail = nw_21
									i_19 = old_20
								}
								a_18[i_19] = destNode_2
							}
						}
					}
				}
			}
		}
	}
}

// Java method: CheckUnreachableCode.growQueue(com.google.javascript.jscomp.graph.LinkedDirectedGraph.LinkedDiGraphNode<com.google.javascript.rhino.Node,com.google.javascript.jscomp.ControlFlowGraph.Branch>[] queue, int head, int tail)
func (this_ *CheckUnreachableCode) priv_CheckUnreachableCode_growQueue_1130(queue []*pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch], head int32, tail int32) []*pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch] {
	_ = queue
	_ = head
	_ = tail
	var active int32 = (tail - head)
	_ = active
	if head > (javalang.Ushr32(int32(len(queue)), int(int32(1)))) {
		func() {
			acSrc_22 := queue
			acSP_23 := head
			acDst_24 := queue
			acDP_25 := int32(0)
			acLen_26 := active
			j2gort.CheckNotNullSlice(acSrc_22)
			j2gort.CheckNotNullSlice(acDst_24)
			if acSP_23 < 0 || acDP_25 < 0 || acLen_26 < 0 || int64(acSP_23)+int64(acLen_26) > int64(len(acSrc_22)) || int64(acDP_25)+int64(acLen_26) > int64(len(acDst_24)) {
				panic(j2gort.ConstructShimThrowable("java.lang.ArrayIndexOutOfBoundsException"))
			}
			copy(acDst_24[int(acDP_25):int(acDP_25)+int(acLen_26)], acSrc_22[int(acSP_23):int(acSP_23)+int(acLen_26)])
		}()
		j2gort.ArraysFillRangeZero(queue, active, tail)
		return queue
	}
	var next []*pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch]
iifeSwitch_1:
	switch {
	default:
		arr_27 := make([]*pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[any, any], (int32(len(queue)) * int32(2)))
		if arr_27 == nil {
			next = nil
			break iifeSwitch_1
		}
		out_28 := make([]*pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch], len(arr_27))
		for i_, el_29 := range arr_27 {
			out_28[i_] = j2gort.CastPtr[pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch]](el_29)
		}
		next = out_28
		break iifeSwitch_1
	}

	_ = next
	func() {
		acSrc_30 := queue
		acSP_31 := head
		acDst_32 := next
		acDP_33 := int32(0)
		acLen_34 := active
		j2gort.CheckNotNullSlice(acSrc_30)
		j2gort.CheckNotNullSlice(acDst_32)
		if acSP_31 < 0 || acDP_33 < 0 || acLen_34 < 0 || int64(acSP_31)+int64(acLen_34) > int64(len(acSrc_30)) || int64(acDP_33)+int64(acLen_34) > int64(len(acDst_32)) {
			panic(j2gort.ConstructShimThrowable("java.lang.ArrayIndexOutOfBoundsException"))
		}
		copy(acDst_32[int(acDP_33):int(acDP_33)+int(acLen_34)], acSrc_30[int(acSP_31):int(acSP_31)+int(acLen_34)])
	}()
	this_.workQueue = next
	return next
}

// Java static method: CheckUnreachableCode.isReachable(com.google.javascript.rhino.Node predecessor, com.google.javascript.jscomp.ControlFlowGraph.Branch branch)
func CheckUnreachableCode_IsReachable(predecessor pkg_rhino.Node, branch ControlFlowGraph_Branch) bool {
	clinit_CheckUnreachableCode()
	_ = predecessor
	_ = branch
	if !(branch.IsConditional()) {
		return true
	}
	var condition pkg_rhino.Node = NodeUtil_GetConditionExpression(predecessor)
	_ = condition
	if condition != 0 {
		var val pkg_base.Tri = NodeUtil_GetBooleanValue(condition)
		_ = val
		if val != pkg_base.StaticGet_Tri_UNKNOWN() {
			return (val.ToBoolean(true) == (branch == StaticGet_ControlFlowGraph_Branch_ON_TRUE()))
		}
	}
	return true
}
func (*CheckUnreachableCode) NoOpVisit_()                 {}
func (*CheckUnreachableCode) CanReleaseSyntacticScopes_() {}
