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
// Port of com/google/javascript/jscomp/DataFlowAnalysis.java (generated mechanically by j2go).

package jscomp

import (
	"sync"
	"sync/atomic"

	pkg_graph "github.com/paullepoulpe/jsgo/transpiled/gen/graph"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java polymorphic class interface: com.google.javascript.jscomp.DataFlowAnalysis
type DataFlowAnalysis[N comparable, L any] interface {
	Implements_jscomp_DataFlowAnalysis()
	IsJavaNil_() bool
	GetDataFlowAnalysisBase() *DataFlowAnalysisBase[N, L]
	InitCfg(ControlFlowGraph[N])
	ReleaseCfgReference()
	GetCfg() ControlFlowGraph[N]
	IsForward() bool
	IsBranched() bool
	NeedsNodeAnnotations() bool
	IsLatticeEqual(L, L) bool
	GetNodeInLattice(*pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch]) L
	GetNodeOutLattice(*pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch]) L
	Join(L, L) L
	CreateFlowJoiner() DataFlowAnalysis_FlowJoiner[L]
	CreateFlowBrancher(N, L) DataFlowAnalysis_FlowBrancher[L]
	FlowThrough(N, L) L
	Analyze()
	CreateInitialEstimateLattice() L
	CreateEntryLattice() L
}

func ToIface_DataFlowAnalysis[N comparable, L any, T_iface_ interface {
	DataFlowAnalysis[N, L]
	comparable
}](v T_iface_) DataFlowAnalysis[N, L] {
	var z T_iface_
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.jscomp.DataFlowAnalysis
type DataFlowAnalysisBase[N comparable, L any] struct {
	self_            DataFlowAnalysis[N, L]
	cfg              ControlFlowGraph[N]                                                            `gson:"cfg"`
	customComparator j2gort.Ext_java_util_Comparator                                                `gson:"customComparator"`
	workElements     []*pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch] `gson:"workElements"`
	cfgNodes         []*pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch] `gson:"cfgNodes"`
	stateIn          []L                                                                            `gson:"stateIn"`
	stateOut         []L                                                                            `gson:"stateOut"`
	edgeState        []L                                                                            `gson:"edgeState"`
	stepCounts       []int32                                                                        `gson:"stepCounts"`
	cfgNodeCount     int32                                                                          `gson:"cfgNodeCount"`
	cfgEdgeCount     int32                                                                          `gson:"cfgEdgeCount"`
	workHead         int32                                                                          `gson:"workHead"`
	workSize         int32                                                                          `gson:"workSize"`
	idHash_          atomic.Int32
	forward          bool `gson:"forward"`
	branched         bool `gson:"branched"`
	ordered          bool `gson:"ordered"`
}

func (this_ *DataFlowAnalysisBase[N, L]) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(DataFlowAnalysis[N, L])
}
func (this_ *DataFlowAnalysisBase[N, L]) IsJavaNil_() bool { return this_ == nil }
func (*DataFlowAnalysisBase[N, L]) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DataFlowAnalysis", "DataFlowAnalysis"
}
func (this_ *DataFlowAnalysisBase[N, L]) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DataFlowAnalysisBase[N, L]) IsGenericInst_jscomp_DataFlowAnalysis_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.DataFlowAnalysis", "com.google.javascript.jscomp.DataFlowAnalysis", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			IsGenericInst_jscomp_DataFlowAnalysis_()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}
func (b *DataFlowAnalysisBase[N, L]) GetDataFlowAnalysisBase() *DataFlowAnalysisBase[N, L] { return b }
func (*DataFlowAnalysisBase[N, L]) Implements_jscomp_DataFlowAnalysis()                    {}

var staticSyncMu_DataFlowAnalysis sync.Mutex

func StaticGet_DataFlowAnalysis_MAX_FLOW_STATE_POOL_SIZE() int32          { return int32(4096) }
func StaticGet_DataFlowAnalysis_MAX_POOLED_WORK_ELEMENTS_CAPACITY() int32 { return int32(4096) }

var static_DataFlowAnalysis_FLOW_STATE_POOL []*DataFlowAnalysis_LinearFlowState

func StaticGet_DataFlowAnalysis_FLOW_STATE_POOL() []*DataFlowAnalysis_LinearFlowState {
	clinit_DataFlowAnalysis()
	return static_DataFlowAnalysis_FLOW_STATE_POOL
}
func StaticSet_DataFlowAnalysis_FLOW_STATE_POOL(v []*DataFlowAnalysis_LinearFlowState) {
	clinit_DataFlowAnalysis()
	static_DataFlowAnalysis_FLOW_STATE_POOL = v
}

var static_DataFlowAnalysis_flowStatePoolSize int32

func StaticGet_DataFlowAnalysis_FlowStatePoolSize() int32 {
	clinit_DataFlowAnalysis()
	return static_DataFlowAnalysis_flowStatePoolSize
}
func StaticSet_DataFlowAnalysis_FlowStatePoolSize(v int32) {
	clinit_DataFlowAnalysis()
	static_DataFlowAnalysis_flowStatePoolSize = v
}
func StaticGet_DataFlowAnalysis_MAX_STEPS_PER_NODE() int32 { return int32(20000) }

var clinitDone_DataFlowAnalysis uint32
var clinitState_DataFlowAnalysis j2gort.ClinitState

func clinit_DataFlowAnalysis() {
	if atomic.LoadUint32(&clinitDone_DataFlowAnalysis) != 0 {
		return
	}
	clinitSlow_DataFlowAnalysis()
}

func clinitSlow_DataFlowAnalysis() {
	if !j2gort.EnterClinit(&clinitState_DataFlowAnalysis, &clinitDone_DataFlowAnalysis) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_DataFlowAnalysis, &clinitDone_DataFlowAnalysis)
	static_DataFlowAnalysis_FLOW_STATE_POOL = make([]*DataFlowAnalysis_LinearFlowState, int32(4096))
}

func (this_ *DataFlowAnalysisBase[N, L]) initFields_DataFlowAnalysis() {
	_ = this_
}

func InitDataFlowAnalysis[N comparable, L any](this_ *DataFlowAnalysisBase[N, L], cfg ControlFlowGraph[N]) {
	_ = cfg
	this_.initFields_DataFlowAnalysis()
	var forward bool = this_.self_.IsForward()
	_ = forward
	var branched bool = this_.self_.IsBranched()
	_ = branched
	this_.forward = forward
	this_.branched = branched
	if branched {
		if !(forward) {
			j2gort.CheckState(false)
		}
	}
	this_.InitCfg(cfg)
}

// Java method: DataFlowAnalysis.initCfg(com.google.javascript.jscomp.ControlFlowGraph<N> cfg)
func (this_ *DataFlowAnalysisBase[N, L]) InitCfg(cfg ControlFlowGraph[N]) {
	_ = cfg
	this_.cfg = cfg
	if cfg.HasPriorityComparator() {
		this_.ordered = true
		this_.customComparator = nil
	} else {
		var cmp j2gort.Ext_java_util_Comparator = cfg.GetOptionalNodeComparator(this_.forward)
		_ = cmp
		this_.ordered = (cmp != nil)
		this_.customComparator = cmp
	}
}

// Java method: DataFlowAnalysis.releaseCfgReference()
func (this_ *DataFlowAnalysisBase[N, L]) ReleaseCfgReference() {
	this_.priv_DataFlowAnalysis_clearWorkQueue_2328()
	if (this_.workElements != nil) && (int32(len(this_.workElements)) > int32(4096)) {
		this_.workElements = nil
	}
	var nCount int32 = this_.cfgNodeCount
	_ = nCount
	if this_.cfgNodes != nil {
		if int32(len(this_.cfgNodes)) > int32(4096) {
			this_.cfgNodes = nil
		} else if nCount > int32(0) {
			j2gort.ArraysFillRangeZero(this_.cfgNodes, int32(0), nCount)
		}
	}
	if this_.stateIn != nil {
		if int32(len(this_.stateIn)) > int32(4096) {
			this_.stateIn = nil
			this_.stateOut = nil
			this_.stepCounts = nil
		} else if nCount > int32(0) {
			var iifeArg_1 L
		iifeSwitch_1:
			switch {
			default:

				if _, ok_ := any((*L)(nil)).(*string); ok_ {
					iifeArg_1 = any(j2gort.NullString).(L)
					break iifeSwitch_1
				}
				var z_ L
				iifeArg_1 = z_
				break iifeSwitch_1
			}
			j2gort.ArraysFillRange(this_.stateIn, int32(0), nCount, iifeArg_1)
			var iifeArg_2 L
		iifeSwitch_2:
			switch {
			default:

				if _, ok_ := any((*L)(nil)).(*string); ok_ {
					iifeArg_2 = any(j2gort.NullString).(L)
					break iifeSwitch_2
				}
				var z_ L
				iifeArg_2 = z_
				break iifeSwitch_2
			}
			j2gort.ArraysFillRange(this_.stateOut, int32(0), nCount, iifeArg_2)
		}
	}
	var eCount int32 = this_.cfgEdgeCount
	_ = eCount
	if this_.edgeState != nil {
		if int32(len(this_.edgeState)) > (int32(4096) * int32(2)) {
			this_.edgeState = nil
		} else if eCount > int32(0) {
			var iifeArg_3 L
		iifeSwitch_3:
			switch {
			default:

				if _, ok_ := any((*L)(nil)).(*string); ok_ {
					iifeArg_3 = any(j2gort.NullString).(L)
					break iifeSwitch_3
				}
				var z_ L
				iifeArg_3 = z_
				break iifeSwitch_3
			}
			j2gort.ArraysFillRange(this_.edgeState, int32(0), eCount, iifeArg_3)
		}
	}
	this_.cfgNodeCount = int32(0)
	this_.cfgEdgeCount = int32(0)
	this_.cfg = nil
	this_.customComparator = nil
}

// Java static method: DataFlowAnalysis.popPooledFlowState()
func DataFlowAnalysis_PopPooledFlowState() *DataFlowAnalysis_LinearFlowState {
	clinit_DataFlowAnalysis()
	staticSyncMu_DataFlowAnalysis.Lock()
	defer staticSyncMu_DataFlowAnalysis.Unlock()
	if StaticGet_DataFlowAnalysis_FlowStatePoolSize() > int32(0) {
		var idx int32
		{
			old_1 := StaticGet_DataFlowAnalysis_FlowStatePoolSize()
			nw_2 := old_1 - 1
			StaticSet_DataFlowAnalysis_FlowStatePoolSize(nw_2)
			idx = nw_2
		}

		_ = idx
		var pooled *DataFlowAnalysis_LinearFlowState = StaticGet_DataFlowAnalysis_FLOW_STATE_POOL()[idx]
		_ = pooled
		{
			a_3 := StaticGet_DataFlowAnalysis_FLOW_STATE_POOL()
			i_4 := idx
			a_3[i_4] = nil
		}
		return pooled
	}
	return nil
}

// Java static method: DataFlowAnalysis.recycleGraphFlowStates(com.google.javascript.jscomp.ControlFlowGraph<N> cfg)
func DataFlowAnalysis_RecycleGraphFlowStates[N comparable](cfg ControlFlowGraph[N]) {
	_ = cfg

}

// Java method: DataFlowAnalysis.getCfg()
func (this_ *DataFlowAnalysisBase[N, L]) GetCfg() ControlFlowGraph[N] {
	return this_.cfg
}

// Java method: DataFlowAnalysis.isForward()
func (this_ *DataFlowAnalysisBase[N, L]) IsForward() bool {
	panic("abstract or native method invoked")
}

// Java method: DataFlowAnalysis.isBranched()
func (this_ *DataFlowAnalysisBase[N, L]) IsBranched() bool {
	return false
}

// Java method: DataFlowAnalysis.needsNodeAnnotations()
func (this_ *DataFlowAnalysisBase[N, L]) NeedsNodeAnnotations() bool {
	return true
}

// Java method: DataFlowAnalysis.isLatticeEqual(L a, L b)
func (this_ *DataFlowAnalysisBase[N, L]) IsLatticeEqual(a L, b L) bool {
	_ = a
	_ = b
	panic("abstract or native method invoked")
}

// Java method: DataFlowAnalysis.getNodeInLattice(com.google.javascript.jscomp.graph.LinkedDirectedGraph.LinkedDiGraphNode<N,com.google.javascript.jscomp.ControlFlowGraph.Branch> node)
func (this_ *DataFlowAnalysisBase[N, L]) GetNodeInLattice(node *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch]) L {
	_ = node
	return this_.stateIn[node.GetDfaIndex()]
}

// Java method: DataFlowAnalysis.getNodeOutLattice(com.google.javascript.jscomp.graph.LinkedDirectedGraph.LinkedDiGraphNode<N,com.google.javascript.jscomp.ControlFlowGraph.Branch> node)
func (this_ *DataFlowAnalysisBase[N, L]) GetNodeOutLattice(node *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch]) L {
	_ = node
	return this_.stateOut[node.GetDfaIndex()]
}

// Java method: DataFlowAnalysis.join(L latticeA, L latticeB)
func (this_ *DataFlowAnalysisBase[N, L]) Join(latticeA L, latticeB L) L {
	_ = latticeA
	_ = latticeB
	var joiner DataFlowAnalysis_FlowJoiner[L] = this_.self_.CreateFlowJoiner()
	_ = joiner
	joiner.JoinFlow(latticeA)
	joiner.JoinFlow(latticeB)
	return joiner.Finish()
}

// Java method: DataFlowAnalysis.createFlowJoiner()
func (this_ *DataFlowAnalysisBase[N, L]) CreateFlowJoiner() DataFlowAnalysis_FlowJoiner[L] {
	panic("abstract or native method invoked")
}

// Java method: DataFlowAnalysis.createFlowBrancher(N node, L output)
func (this_ *DataFlowAnalysisBase[N, L]) CreateFlowBrancher(node N, output L) DataFlowAnalysis_FlowBrancher[L] {
	_ = node
	_ = output
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException")))
}

// Java method: DataFlowAnalysis.flowThrough(N node, L input)
func (this_ *DataFlowAnalysisBase[N, L]) FlowThrough(node N, input L) L {
	_ = node
	_ = input
	panic("abstract or native method invoked")
}

// Java method: DataFlowAnalysis.analyze()
func (this_ *DataFlowAnalysisBase[N, L]) Analyze() {
	this_.priv_DataFlowAnalysis_initialize_2321()
	var forward bool = this_.forward
	_ = forward
	var implicitReturn *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch] = this_.cfg.GetImplicitReturnNode()
	_ = implicitReturn
	var steps []int32 = this_.stepCounts
	_ = steps
	for this_.workSize > int32(0) {
		var curNode *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch] = this_.priv_DataFlowAnalysis_dequeueWorkNode_2327()
		_ = curNode
		var curIdx int32 = curNode.GetDfaIndex()
		_ = curIdx
		{
			var iifeCond_1 int32
			{

				a_12 := steps
				i_13 := curIdx
				old_14 := a_12[i_13]
				nw_15 := old_14 + 1
				a_12[i_13] = nw_15
				iifeCond_1 = old_14
			}
			if iifeCond_1 > int32(20000) {
				this_.priv_DataFlowAnalysis_clearWorkQueue_2328()
				panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Dataflow analysis appears to diverge around: "+j2gort.StrVal(curNode))))
			}
		}
		this_.priv_DataFlowAnalysis_joinInputs_2323(curNode, curIdx)
		if this_.priv_DataFlowAnalysis_flow_2322(curNode, curIdx) {
			if forward {
				{
					var i int32 = int32(0)
					_ = i
					var count int32 = curNode.GetOutEdgeCount()
					_ = count
					for ; i < count; i = i + 1 {
						var nextNode *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch] = j2gort.CastPtr[pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch]](j2gort.CastPtr[pkg_graph.LinkedDirectedGraph_LinkedDiGraphEdge[N, ControlFlowGraph_Branch]](curNode.GetOutEdgeAt(i)).GetDestination())
						_ = nextNode
						if nextNode != implicitReturn {
							this_.priv_DataFlowAnalysis_enqueueWorkNode_2326(nextNode)
						}
					}
				}
			} else {
				{
					var i_2 int32 = int32(0)
					_ = i_2
					var count_2 int32 = curNode.GetInEdgeCount()
					_ = count_2
					for ; i_2 < count_2; i_2 = i_2 + 1 {
						var nextNode_2 *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch] = j2gort.CastPtr[pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch]](j2gort.CastPtr[pkg_graph.LinkedDirectedGraph_LinkedDiGraphEdge[N, ControlFlowGraph_Branch]](curNode.GetInEdgeAt(i_2)).GetSource())
						_ = nextNode_2
						if nextNode_2 != implicitReturn {
							this_.priv_DataFlowAnalysis_enqueueWorkNode_2326(nextNode_2)
						}
					}
				}
			}
		}
	}
	if forward {
		this_.priv_DataFlowAnalysis_joinInputs_2323(implicitReturn, implicitReturn.GetDfaIndex())
	}
	if this_.self_.NeedsNodeAnnotations() {
		this_.priv_DataFlowAnalysis_syncNodeAnnotations_2318()
	}
}

// Java method: DataFlowAnalysis.syncNodeAnnotations()
func (this_ *DataFlowAnalysisBase[N, L]) priv_DataFlowAnalysis_syncNodeAnnotations_2318() {
	var count int32 = this_.cfgNodeCount
	_ = count
	var nodes []*pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch] = this_.cfgNodes
	_ = nodes
	var inArr []L = this_.stateIn
	_ = inArr
	var outArr []L = this_.stateOut
	_ = outArr
	var stepsArr []int32 = this_.stepCounts
	_ = stepsArr
	{
		var i int32 = int32(0)
		_ = i
		for ; i < count; i = i + 1 {
			var state *DataFlowAnalysis_LinearFlowState = j2gort.CastPtr[DataFlowAnalysis_LinearFlowState](nodes[i].GetAnnotation())
			_ = state
			state.SyncResult(stepsArr[i], j2gort.ToAny(inArr[i]), j2gort.ToAny(outArr[i]))
		}
	}
}

// Java method: DataFlowAnalysis.createInitialEstimateLattice()
func (this_ *DataFlowAnalysisBase[N, L]) CreateInitialEstimateLattice() L {
	panic("abstract or native method invoked")
}

// Java method: DataFlowAnalysis.createEntryLattice()
func (this_ *DataFlowAnalysisBase[N, L]) CreateEntryLattice() L {
	panic("abstract or native method invoked")
}

// Java method: DataFlowAnalysis.initialize()
func (this_ *DataFlowAnalysisBase[N, L]) priv_DataFlowAnalysis_initialize_2321() {
	this_.priv_DataFlowAnalysis_clearWorkQueue_2328()
	var prevNodeCount int32 = this_.cfgNodeCount
	_ = prevNodeCount
	if (prevNodeCount > int32(0)) && (this_.cfgNodes != nil) {
		j2gort.ArraysFillRangeZero(this_.cfgNodes, int32(0), prevNodeCount)
	}
	var nodeCount int32 = this_.cfg.GetNodeCount()
	_ = nodeCount
	this_.cfgNodeCount = nodeCount
	if this_.ordered {
		var cap_ int32
		_ = cap_
		if nodeCount > int32(4) {
			cap_ = nodeCount
		} else {
			cap_ = int32(4)
		}
		if (this_.workElements == nil) || (int32(len(this_.workElements)) < cap_) {
			this_.workElements = make([]*pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch], cap_)
		}
	} else {
		var cap__2 int32 = int32(8)
		_ = cap__2
		for cap__2 < nodeCount {
			cap__2 = javalang.Shl32(cap__2, int(int32(1)))
		}
		if (this_.workElements == nil) || (int32(len(this_.workElements)) < cap__2) {
			this_.workElements = make([]*pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch], cap__2)
		}
	}
	var stateCap int32
	_ = stateCap
	if nodeCount > int32(8) {
		stateCap = nodeCount
	} else {
		stateCap = int32(8)
	}
	if (this_.cfgNodes == nil) || (int32(len(this_.cfgNodes)) < stateCap) {
		this_.cfgNodes = make([]*pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch], stateCap)
	}
	if (this_.stateIn == nil) || (int32(len(this_.stateIn)) < stateCap) {
		this_.stateIn = func() []L {
			arr_16 := make([]L, stateCap)
			if _, ok_ := any((*L)(nil)).(*string); ok_ {
				ns_ := any(j2gort.NullString).(L)
				for i_ := range arr_16 {
					arr_16[i_] = ns_
				}
			}
			return arr_16
		}()
		this_.stateOut = func() []L {
			arr_17 := make([]L, stateCap)
			if _, ok_ := any((*L)(nil)).(*string); ok_ {
				ns_ := any(j2gort.NullString).(L)
				for i_ := range arr_17 {
					arr_17[i_] = ns_
				}
			}
			return arr_17
		}()
		this_.stepCounts = make([]int32, stateCap)
	} else if prevNodeCount > nodeCount {
		var iifeArg_4 L
	iifeSwitch_4:
		switch {
		default:

			if _, ok_ := any((*L)(nil)).(*string); ok_ {
				iifeArg_4 = any(j2gort.NullString).(L)
				break iifeSwitch_4
			}
			var z_ L
			iifeArg_4 = z_
			break iifeSwitch_4
		}
		j2gort.ArraysFillRange(this_.stateIn, nodeCount, prevNodeCount, iifeArg_4)
		var iifeArg_5 L
	iifeSwitch_5:
		switch {
		default:

			if _, ok_ := any((*L)(nil)).(*string); ok_ {
				iifeArg_5 = any(j2gort.NullString).(L)
				break iifeSwitch_5
			}
			var z_ L
			iifeArg_5 = z_
			break iifeSwitch_5
		}
		j2gort.ArraysFillRange(this_.stateOut, nodeCount, prevNodeCount, iifeArg_5)
	}
	var collected int32 = this_.cfg.CollectNodes(this_.cfgNodes)
	_ = collected
	if !(collected == nodeCount) {
		j2gort.CheckState(false)
	}
	var recordAnnotations bool = this_.self_.NeedsNodeAnnotations()
	_ = recordAnnotations
	var implicitReturn *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch] = this_.cfg.GetImplicitReturnNode()
	_ = implicitReturn
	var nodes []*pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch] = this_.cfgNodes
	_ = nodes
	var inArr []L = this_.stateIn
	_ = inArr
	var outArr []L = this_.stateOut
	_ = outArr
	var stepsArr []int32 = this_.stepCounts
	_ = stepsArr
	{
		var i int32 = int32(0)
		_ = i
		for ; i < nodeCount; i = i + 1 {
			var node *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch] = nodes[i]
			_ = node
			node.SetDfaIndex(i)
			node.SetInDfaQueue(false)
			var initIn L = this_.self_.CreateInitialEstimateLattice()
			_ = initIn
			var initOut L = this_.self_.CreateInitialEstimateLattice()
			_ = initOut
			{
				a_18 := inArr
				i_19 := i
				a_18[i_19] = initIn
			}
			{
				a_20 := outArr
				i_21 := i
				a_20[i_21] = initOut
			}
			{
				a_22 := stepsArr
				i_23 := i
				a_22[i_23] = int32(0)
			}
			if recordAnnotations {
				var existing pkg_graph.Annotation = node.GetAnnotation()
				_ = existing
				var raw *DataFlowAnalysis_LinearFlowState
				_ = raw
				var ok_24 bool
				raw, ok_24 = j2gort.TryCast[*DataFlowAnalysis_LinearFlowState](existing)
				if ok_24 {
					var reused *DataFlowAnalysis_LinearFlowState = raw
					_ = reused
					reused.Reinit(j2gort.ToAny(initIn), j2gort.ToAny(initOut))
				} else {
					if cfgBase_ := this_.cfg.GetControlFlowGraphBase(); cfgBase_ != nil {
						node.SetAnnotation(allocPooled_DataFlowAnalysis_LinearFlowState(&cfgBase_.slab_LinearFlowState, j2gort.ToAny(initIn), j2gort.ToAny(initOut)))
					} else {
						node.SetAnnotation(NewDataFlowAnalysis_LinearFlowState(j2gort.ToAny(initIn), j2gort.ToAny(initOut)))
					}

				}
			}
			if node != implicitReturn {
				this_.priv_DataFlowAnalysis_enqueueWorkNode_2326(node)
			}
		}
	}
	var prevEdgeCount int32 = this_.cfgEdgeCount
	_ = prevEdgeCount
	if this_.branched {
		var totalEdges int32 = int32(0)
		_ = totalEdges
		{
			var i_2 int32 = int32(0)
			_ = i_2
			for ; i_2 < nodeCount; i_2 = i_2 + 1 {
				totalEdges = (totalEdges + nodes[i_2].GetOutEdgeCount())
			}
		}
		this_.cfgEdgeCount = totalEdges
		var edgeCap int32
		_ = edgeCap
		if totalEdges > int32(8) {
			edgeCap = totalEdges
		} else {
			edgeCap = int32(8)
		}
		if (this_.edgeState == nil) || (int32(len(this_.edgeState)) < edgeCap) {
			this_.edgeState = func() []L {
				arr_25 := make([]L, edgeCap)
				if _, ok_ := any((*L)(nil)).(*string); ok_ {
					ns_ := any(j2gort.NullString).(L)
					for i_ := range arr_25 {
						arr_25[i_] = ns_
					}
				}
				return arr_25
			}()
		} else if prevEdgeCount > totalEdges {
			var iifeArg_6 L
		iifeSwitch_6:
			switch {
			default:

				if _, ok_ := any((*L)(nil)).(*string); ok_ {
					iifeArg_6 = any(j2gort.NullString).(L)
					break iifeSwitch_6
				}
				var z_ L
				iifeArg_6 = z_
				break iifeSwitch_6
			}
			j2gort.ArraysFillRange(this_.edgeState, totalEdges, prevEdgeCount, iifeArg_6)
		}
		var edges []L = this_.edgeState
		_ = edges
		var edgeIdx int32 = int32(0)
		_ = edgeIdx
		{
			var i_3 int32 = int32(0)
			_ = i_3
			for ; i_3 < nodeCount; i_3 = i_3 + 1 {
				var node_2 *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch] = nodes[i_3]
				_ = node_2
				{
					var j int32 = int32(0)
					_ = j
					var outCount int32 = node_2.GetOutEdgeCount()
					_ = outCount
					for ; j < outCount; j = j + 1 {
						var outEdge *pkg_graph.LinkedDirectedGraph_LinkedDiGraphEdge[N, ControlFlowGraph_Branch] = j2gort.CastPtr[pkg_graph.LinkedDirectedGraph_LinkedDiGraphEdge[N, ControlFlowGraph_Branch]](node_2.GetOutEdgeAt(j))
						_ = outEdge
						outEdge.SetDfaEdgeIndex(edgeIdx)
						{
							a_26 := edges
							var i_27 int32
							{

								old_28 := edgeIdx
								nw_29 := old_28 + 1
								edgeIdx = nw_29
								i_27 = old_28
							}
							a_26[i_27] = this_.self_.CreateInitialEstimateLattice()
						}
					}
				}
			}
		}
	} else {
		if (prevEdgeCount > int32(0)) && (this_.edgeState != nil) {
			var iifeArg_7 L
		iifeSwitch_7:
			switch {
			default:

				if _, ok_ := any((*L)(nil)).(*string); ok_ {
					iifeArg_7 = any(j2gort.NullString).(L)
					break iifeSwitch_7
				}
				var z_ L
				iifeArg_7 = z_
				break iifeSwitch_7
			}
			j2gort.ArraysFillRange(this_.edgeState, int32(0), prevEdgeCount, iifeArg_7)
		}
		this_.cfgEdgeCount = int32(0)
	}
}

// Java method: DataFlowAnalysis.flow(com.google.javascript.jscomp.graph.LinkedDirectedGraph.LinkedDiGraphNode<N,com.google.javascript.jscomp.ControlFlowGraph.Branch> node, int nodeIdx)
func (this_ *DataFlowAnalysisBase[N, L]) priv_DataFlowAnalysis_flow_2322(node *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch], nodeIdx int32) bool {
	_ = node
	_ = nodeIdx
	if this_.forward {
		var outBefore L = this_.stateOut[nodeIdx]
		_ = outBefore
		var outAfter L
		{
			tv_30 := this_.self_.FlowThrough(node.GetValue(), this_.stateIn[nodeIdx])
			_ = j2gort.CheckNotNull(j2gort.ToAny(tv_30))
			outAfter = tv_30
		}

		_ = outAfter
		{
			a_31 := this_.stateOut
			i_32 := nodeIdx
			a_31[i_32] = outAfter
		}
		var changed bool = !(this_.self_.IsLatticeEqual(outBefore, outAfter))
		_ = changed
		if this_.branched {
			var brancher DataFlowAnalysis_FlowBrancher[L] = this_.self_.CreateFlowBrancher(node.GetValue(), outAfter)
			_ = brancher
			var edges []L = this_.edgeState
			_ = edges
			{
				var i int32 = int32(0)
				_ = i
				var outCount int32 = node.GetOutEdgeCount()
				_ = outCount
				for ; i < outCount; i = i + 1 {
					var outEdge *pkg_graph.LinkedDirectedGraph_LinkedDiGraphEdge[N, ControlFlowGraph_Branch] = j2gort.CastPtr[pkg_graph.LinkedDirectedGraph_LinkedDiGraphEdge[N, ControlFlowGraph_Branch]](node.GetOutEdgeAt(i))
					_ = outEdge
					var edgeIdx int32 = outEdge.GetDfaEdgeIndex()
					_ = edgeIdx
					var outBranchBefore L = edges[edgeIdx]
					_ = outBranchBefore
					var outBranchAfter L
					{
						tv_33 := brancher.BranchFlow(outEdge.GetValue())
						_ = j2gort.CheckNotNull(j2gort.ToAny(tv_33))
						outBranchAfter = tv_33
					}

					_ = outBranchAfter
					{
						a_34 := edges
						i_35 := edgeIdx
						a_34[i_35] = outBranchAfter
					}
					if !(changed) {
						changed = !(this_.self_.IsLatticeEqual(outBranchBefore, outBranchAfter))
					}
				}
			}
		}
		return changed
	} else {
		var inBefore L = this_.stateIn[nodeIdx]
		_ = inBefore
		var inAfter L
		{
			tv_36 := this_.self_.FlowThrough(node.GetValue(), this_.stateOut[nodeIdx])
			_ = j2gort.CheckNotNull(j2gort.ToAny(tv_36))
			inAfter = tv_36
		}

		_ = inAfter
		{
			a_37 := this_.stateIn
			i_38 := nodeIdx
			a_37[i_38] = inAfter
		}
		return !(this_.self_.IsLatticeEqual(inBefore, inAfter))
	}
}

// Java method: DataFlowAnalysis.joinInputs(com.google.javascript.jscomp.graph.LinkedDirectedGraph.LinkedDiGraphNode<N,com.google.javascript.jscomp.ControlFlowGraph.Branch> node, int nodeIdx)
func (this_ *DataFlowAnalysisBase[N, L]) priv_DataFlowAnalysis_joinInputs_2323(node *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch], nodeIdx int32) {
	_ = node
	_ = nodeIdx
	var forward bool = this_.forward
	_ = forward
	if forward && (this_.cfg.GetEntryNode() == node) {
		{
			a_40 := this_.stateIn
			i_41 := nodeIdx
			a_40[i_41] = func() L {
				tv_39 := this_.self_.CreateEntryLattice()
				_ = j2gort.CheckNotNull(j2gort.ToAny(tv_39))
				return tv_39
			}()
		}
		return
	}
	var edgeCount int32
	_ = edgeCount
	if forward {
		edgeCount = node.GetInEdgeCount()
	} else {
		edgeCount = node.GetOutEdgeCount()
	}
	var result L
iifeSwitch_9:
	switch {
	default:
		if _, ok_ := any((*L)(nil)).(*string); ok_ {
			result = any(j2gort.NullString).(L)
			break iifeSwitch_9
		}
		var z_ L
		result = z_
		break iifeSwitch_9
	}

	_ = result
	switch edgeCount {
	case int32(0):
		return
	case int32(1):
		var iifeArg_8 *pkg_graph.
			LinkedDirectedGraph_LinkedDiGraphEdge[N,
			ControlFlowGraph_Branch,
		]
	iifeSwitch_8:
		switch {
		default:

			if forward {
				iifeArg_8 = j2gort.CastPtr[pkg_graph.LinkedDirectedGraph_LinkedDiGraphEdge[N, ControlFlowGraph_Branch]](node.GetInEdgeAt(int32(0)))
				break iifeSwitch_8
			}
			iifeArg_8 = j2gort.CastPtr[pkg_graph.LinkedDirectedGraph_LinkedDiGraphEdge[N, ControlFlowGraph_Branch]](node.GetOutEdgeAt(int32(0)))
			break iifeSwitch_8
		}
		result = this_.priv_DataFlowAnalysis_getInputFromEdge_2324(iifeArg_8)
	default:
		var joiner DataFlowAnalysis_FlowJoiner[L] = this_.self_.CreateFlowJoiner()
		_ = joiner
		if forward {
			{
				var i int32 = int32(0)
				_ = i
				for ; i < edgeCount; i = i + 1 {
					joiner.JoinFlow(this_.priv_DataFlowAnalysis_getInputFromEdge_2324(j2gort.CastPtr[pkg_graph.LinkedDirectedGraph_LinkedDiGraphEdge[N, ControlFlowGraph_Branch]](node.GetInEdgeAt(i))))
				}
			}
		} else {
			{
				var i_2 int32 = int32(0)
				_ = i_2
				for ; i_2 < edgeCount; i_2 = i_2 + 1 {
					joiner.JoinFlow(this_.priv_DataFlowAnalysis_getInputFromEdge_2324(j2gort.CastPtr[pkg_graph.LinkedDirectedGraph_LinkedDiGraphEdge[N, ControlFlowGraph_Branch]](node.GetOutEdgeAt(i_2))))
				}
			}
		}
		result = joiner.Finish()
	}
	if forward {
		{
			a_43 := this_.stateIn
			i_44 := nodeIdx
			a_43[i_44] = func() L { tv_42 := result; _ = j2gort.CheckNotNull(j2gort.ToAny(tv_42)); return tv_42 }()
		}
	} else {
		{
			a_46 := this_.stateOut
			i_47 := nodeIdx
			a_46[i_47] = func() L { tv_45 := result; _ = j2gort.CheckNotNull(j2gort.ToAny(tv_45)); return tv_45 }()
		}
	}
}

// Java method: DataFlowAnalysis.getInputFromEdge(com.google.javascript.jscomp.graph.LinkedDirectedGraph.LinkedDiGraphEdge<N,com.google.javascript.jscomp.ControlFlowGraph.Branch> edge)
func (this_ *DataFlowAnalysisBase[N, L]) priv_DataFlowAnalysis_getInputFromEdge_2324(edge *pkg_graph.LinkedDirectedGraph_LinkedDiGraphEdge[N, ControlFlowGraph_Branch]) L {
	_ = edge
	if this_.branched {
		return this_.edgeState[edge.GetDfaEdgeIndex()]
	} else if this_.forward {
		return this_.stateOut[j2gort.CastPtr[pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch]](edge.GetSource()).GetDfaIndex()]
	} else {
		var node *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch] = j2gort.CastPtr[pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch]](edge.GetDestination())
		_ = node
		if node == this_.cfg.GetImplicitReturnNode() {
			return this_.self_.CreateEntryLattice()
		}
		return this_.stateIn[node.GetDfaIndex()]
	}
}

// Java static method: DataFlowAnalysis.computeEscaped(com.google.javascript.jscomp.Scope jsScope, java.util.Set<com.google.javascript.jscomp.Var> escaped, com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.ScopeCreator scopeCreator, java.util.Map<java.lang.String,com.google.javascript.jscomp.Var> allVarsInFn)
func DataFlowAnalysis_ComputeEscaped(jsScope *Scope, escaped j2gort.Ext_java_util_Set, compiler AbstractCompiler, scopeCreator ScopeCreator, allVarsInFn j2gort.Ext_java_util_Map) {
	clinit_DataFlowAnalysis()
	_ = jsScope
	_ = escaped
	_ = compiler
	_ = scopeCreator
	_ = allVarsInFn
	if !(jsScope.IsFunctionScope()) {
		j2gort.CheckArgument(false)
	}
	var finder NodeTraversal_AbstractPostOrderCallback = NewDataFlowAnalysis__Anon_21(jsScope, escaped)
	_ = finder
	{
		clinit_NodeTraversal()
		var ntb_ NodeTraversal_Builder
		InitNodeTraversal_Builder(&ntb_)
		(&ntb_).SetCompiler(compiler)
		(&ntb_).SetCallback_1_NodeTraversal_Callback(Upcast_NodeTraversal_Callback(finder))
		(&ntb_).SetScopeCreator(scopeCreator)
		(&ntb_).TraverseAtScope(ToIface_AbstractScope(jsScope))
	}
	{
		cur_48 := j2gort.NewMapValCursor(allVarsInFn)
		for cur_48.HasNext() {
			rawElem_49 := cur_48.Next()
			var_ := j2gort.CastPtr[Var](rawElem_49)
			_ = var_
			if var_.GetParentNode().IsCatch() || compiler.GetCodingConvention().IsExported_1(var_.GetName()) {
				escaped.Add_1(j2gort.PtrToAny(var_))
			}
		}
	}
}

// Java static method: DataFlowAnalysis.isReferencedInNonStaticClassField(com.google.javascript.jscomp.Scope refScope, com.google.javascript.jscomp.Scope targetScope)
func DataFlowAnalysis_IsReferencedInNonStaticClassField(refScope *Scope, targetScope *Scope) bool {
	clinit_DataFlowAnalysis()
	_ = refScope
	_ = targetScope
	{
		var s *Scope = refScope
		_ = s
		for ; (s != nil) && (s != targetScope); s = j2gort.CastPtr[Scope](s.GetParent()) {
			var root pkg_rhino.Node = s.GetRootNode_0()
			_ = root
			if (root.IsMemberFieldDef() || root.IsComputedFieldDef()) && !(root.IsStaticMember()) {
				return true
			}
		}
	}
	return false
}

// Java method: DataFlowAnalysis.compareWorkNodes(com.google.javascript.jscomp.graph.LinkedDirectedGraph.LinkedDiGraphNode<N,com.google.javascript.jscomp.ControlFlowGraph.Branch> a, com.google.javascript.jscomp.graph.LinkedDirectedGraph.LinkedDiGraphNode<N,com.google.javascript.jscomp.ControlFlowGraph.Branch> b)
func (this_ *DataFlowAnalysisBase[N, L]) priv_DataFlowAnalysis_compareWorkNodes_2325(a *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch], b *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch]) int32 {
	_ = a
	_ = b
	var cmp j2gort.Ext_java_util_Comparator = this_.customComparator
	_ = cmp
	if cmp != nil {
		return cmp.Compare(j2gort.PtrToAny(a), j2gort.PtrToAny(b))
	}
	var pa int32 = a.GetPriority()
	_ = pa
	var pb int32 = b.GetPriority()
	_ = pb
	if this_.forward {
		{

			if pa < pb {
				return int32(-1)
			}
			{

				if pa == pb {
					return int32(0)
				}
				return int32(1)
			}
		}

	} else {
		{

			if pb < pa {
				return int32(-1)
			}
			{

				if pb == pa {
					return int32(0)
				}
				return int32(1)
			}
		}

	}
}

// Java method: DataFlowAnalysis.enqueueWorkNode(com.google.javascript.jscomp.graph.LinkedDirectedGraph.LinkedDiGraphNode<N,com.google.javascript.jscomp.ControlFlowGraph.Branch> node)
func (this_ *DataFlowAnalysisBase[N, L]) priv_DataFlowAnalysis_enqueueWorkNode_2326(node *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch]) {
	_ = node
	if node.IsInDfaQueue() {
		return
	}
	node.SetInDfaQueue(true)
	var elems []*pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch] = this_.workElements
	_ = elems
	if this_.ordered {
		var k int32 = this_.workSize
		_ = k
		if k == int32(len(elems)) {
			{

				arrSrc_50 := elems
				arrLen_51 := (int32(len(elems)) * int32(2))
				j2gort.CheckNotNullSlice(arrSrc_50)
				if arrLen_51 < 0 {
					panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
				}
				arrOut_52 := make([]*pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch], arrLen_51)
				copy(arrOut_52, arrSrc_50)
				elems = arrOut_52
			}

			this_.workElements = elems
		}
		this_.workSize = (k + int32(1))
		for k > int32(0) {
			var parent int32 = javalang.Ushr32((k - int32(1)), int(int32(1)))
			_ = parent
			var e *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch] = elems[parent]
			_ = e
			if this_.priv_DataFlowAnalysis_compareWorkNodes_2325(node, e) >= int32(0) {
				break
			}
			{
				a_53 := elems
				i_54 := k
				a_53[i_54] = e
			}
			k = parent
		}
		{
			a_55 := elems
			i_56 := k
			a_55[i_56] = node
		}
	} else {
		if this_.workSize == int32(len(elems)) {
			var oldLen int32 = int32(len(elems))
			_ = oldLen
			var grown []*pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch] = make([]*pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch], (oldLen * int32(2)))
			_ = grown
			var mask int32 = (oldLen - int32(1))
			_ = mask
			{
				var i int32 = int32(0)
				_ = i
				for ; i < oldLen; i = i + 1 {
					{
						a_57 := grown
						i_58 := i
						a_57[i_58] = elems[((this_.workHead + i) & mask)]
					}
				}
			}
			elems = grown
			this_.workElements = elems
			this_.workHead = int32(0)
		}
		var tail int32 = ((this_.workHead + this_.workSize) & (int32(len(elems)) - int32(1)))
		_ = tail
		{
			a_59 := elems
			i_60 := tail
			a_59[i_60] = node
		}
		this_.workSize = this_.workSize + 1
	}
}

// Java method: DataFlowAnalysis.dequeueWorkNode()
func (this_ *DataFlowAnalysisBase[N, L]) priv_DataFlowAnalysis_dequeueWorkNode_2327() *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch] {
	var elems []*pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch] = this_.workElements
	_ = elems
	if this_.ordered {
		var result *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch] = elems[int32(0)]
		_ = result
		result.SetInDfaQueue(false)
		var n int32
		{
			old_61 := this_.workSize
			nw_62 := old_61 - 1
			this_.workSize = nw_62
			n = nw_62
		}
		_ = n
		var x *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch] = elems[n]
		_ = x
		{
			a_63 := elems
			i_64 := n
			a_63[i_64] = nil
		}
		if n > int32(0) {
			var half int32 = javalang.Ushr32(n, int(int32(1)))
			_ = half
			var k int32 = int32(0)
			_ = k
			for k < half {
				var child int32 = ((javalang.Shl32(k, int(int32(1)))) + int32(1))
				_ = child
				var c *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch] = elems[child]
				_ = c
				var right int32 = (child + int32(1))
				_ = right
				if (right < n) && (this_.priv_DataFlowAnalysis_compareWorkNodes_2325(c, elems[right]) > int32(0)) {
					child = right
					c = elems[right]
				}
				if this_.priv_DataFlowAnalysis_compareWorkNodes_2325(x, c) <= int32(0) {
					break
				}
				{
					a_65 := elems
					i_66 := k
					a_65[i_66] = c
				}
				k = child
			}
			{
				a_67 := elems
				i_68 := k
				a_67[i_68] = x
			}
		}
		return result
	} else {
		var head int32 = this_.workHead
		_ = head
		var result_2 *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch] = elems[head]
		_ = result_2
		{
			a_69 := elems
			i_70 := head
			a_69[i_70] = nil
		}
		result_2.SetInDfaQueue(false)
		this_.workHead = ((head + int32(1)) & (int32(len(elems)) - int32(1)))
		this_.workSize = this_.workSize - 1
		return result_2
	}
}

// Java method: DataFlowAnalysis.clearWorkQueue()
func (this_ *DataFlowAnalysisBase[N, L]) priv_DataFlowAnalysis_clearWorkQueue_2328() {
	var elems []*pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch] = this_.workElements
	_ = elems
	if (elems != nil) && (this_.workSize > int32(0)) {
		if this_.ordered {
			{
				var i int32 = int32(0)
				_ = i
				var sz int32 = this_.workSize
				_ = sz
				for ; i < sz; i = i + 1 {
					var n *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch] = elems[i]
					_ = n
					if n != nil {
						n.SetInDfaQueue(false)
						{
							a_71 := elems
							i_72 := i
							a_71[i_72] = nil
						}
					}
				}
			}
		} else {
			var mask int32 = (int32(len(elems)) - int32(1))
			_ = mask
			{
				var i_2 int32 = int32(0)
				_ = i_2
				var sz_2 int32 = this_.workSize
				_ = sz_2
				for ; i_2 < sz_2; i_2 = i_2 + 1 {
					var idx int32 = ((this_.workHead + i_2) & mask)
					_ = idx
					var n_2 *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch] = elems[idx]
					_ = n_2
					if n_2 != nil {
						n_2.SetInDfaQueue(false)
						{
							a_73 := elems
							i_74 := idx
							a_73[i_74] = nil
						}
					}
				}
			}
		}
	}
	this_.workHead = int32(0)
	this_.workSize = int32(0)
}

// Java interface: com.google.javascript.jscomp.DataFlowAnalysis.FlowJoiner
type DataFlowAnalysis_FlowJoiner[L any] interface {
	Implements_jscomp_DataFlowAnalysis_FlowJoiner()
	IsJavaNil_() bool
	JoinFlow(L)
	Finish() L
}

func ToIface_DataFlowAnalysis_FlowJoiner[L any, T_iface_ interface {
	DataFlowAnalysis_FlowJoiner[L]
	comparable
}](v T_iface_) DataFlowAnalysis_FlowJoiner[L] {
	var z T_iface_
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.DataFlowAnalysis.FlowJoiner", "com.google.javascript.jscomp.DataFlowAnalysis$FlowJoiner", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			Implements_jscomp_DataFlowAnalysis_FlowJoiner()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}

// Java interface: com.google.javascript.jscomp.DataFlowAnalysis.FlowBrancher
type DataFlowAnalysis_FlowBrancher[L any] interface {
	Implements_jscomp_DataFlowAnalysis_FlowBrancher()
	IsJavaNil_() bool
	BranchFlow(ControlFlowGraph_Branch) L
}

func ToIface_DataFlowAnalysis_FlowBrancher[L any, T_iface_ interface {
	DataFlowAnalysis_FlowBrancher[L]
	comparable
}](v T_iface_) DataFlowAnalysis_FlowBrancher[L] {
	var z T_iface_
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.DataFlowAnalysis.FlowBrancher", "com.google.javascript.jscomp.DataFlowAnalysis$FlowBrancher", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			Implements_jscomp_DataFlowAnalysis_FlowBrancher()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}

// Java class: com.google.javascript.jscomp.DataFlowAnalysis.LinearFlowState
type DataFlowAnalysis_LinearFlowState struct {
	in        any   `gson:"in"`
	out       any   `gson:"out"`
	stepCount int32 `gson:"stepCount"`
	idHash_   atomic.Int32
}

func (this_ *DataFlowAnalysis_LinearFlowState) BindSelf_(s any) {
	_ = s
}
func (this_ *DataFlowAnalysis_LinearFlowState) IsJavaNil_() bool { return this_ == nil }
func (*DataFlowAnalysis_LinearFlowState) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DataFlowAnalysis$LinearFlowState", "LinearFlowState"
}
func (this_ *DataFlowAnalysis_LinearFlowState) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DataFlowAnalysis_LinearFlowState) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.DataFlowAnalysis.LinearFlowState", "com.google.javascript.jscomp.DataFlowAnalysis$LinearFlowState", func(v_ any) bool {
		return j2gort.InstanceOf[*DataFlowAnalysis_LinearFlowState](v_)
	})
}
func (*DataFlowAnalysis_LinearFlowState) Implements_graph_Annotation() {}

func (this_ *DataFlowAnalysis_LinearFlowState) initFields_DataFlowAnalysis_LinearFlowState() {
	_ = this_

}

// Java constructor: LinearFlowState(L in, L out)
func NewDataFlowAnalysis_LinearFlowState(in any, out any) *DataFlowAnalysis_LinearFlowState {
	this_ := &DataFlowAnalysis_LinearFlowState{}
	InitDataFlowAnalysis_LinearFlowState(this_, in, out)
	return this_
}
func allocPooled_DataFlowAnalysis_LinearFlowState(slab_ *j2gort.PooledSlab[DataFlowAnalysis_LinearFlowState], in any, out any) *DataFlowAnalysis_LinearFlowState {
	this_ := slab_.Alloc()
	InitDataFlowAnalysis_LinearFlowState(this_, in, out)
	return this_
}

func InitDataFlowAnalysis_LinearFlowState(this_ *DataFlowAnalysis_LinearFlowState, in any, out any) {
	_ = in
	_ = out
	this_.initFields_DataFlowAnalysis_LinearFlowState()
	j2gort.CheckNotNull(in)
	j2gort.CheckNotNull(out)
	this_.in = in
	this_.out = out
}

// Java method: LinearFlowState.reinit(L in, L out)
func (this_ *DataFlowAnalysis_LinearFlowState) Reinit(in any, out any) {
	_ = in
	_ = out
	j2gort.CheckNotNull(in)
	j2gort.CheckNotNull(out)
	this_.stepCount = int32(0)
	this_.in = in
	this_.out = out
}

// Java method: LinearFlowState.syncResult(int stepCount, L in, L out)
func (this_ *DataFlowAnalysis_LinearFlowState) SyncResult(stepCount int32, in any, out any) {
	_ = stepCount
	_ = in
	_ = out
	this_.stepCount = stepCount
	this_.in = in
	this_.out = out
}

// Java method: LinearFlowState.clearForPool()
func (this_ *DataFlowAnalysis_LinearFlowState) ClearForPool() {
	this_.stepCount = int32(0)
	this_.in = nil
	this_.out = nil
}

// Java method: LinearFlowState.getStepCount()
func (this_ *DataFlowAnalysis_LinearFlowState) GetStepCount() int32 {
	return this_.stepCount
}

// Java method: LinearFlowState.getIn()
func (this_ *DataFlowAnalysis_LinearFlowState) GetIn() any {
	return this_.in
}

// Java method: LinearFlowState.setIn(L in)
func (this_ *DataFlowAnalysis_LinearFlowState) priv_LinearFlowState_setIn_2337(in any) {
	_ = in
	j2gort.CheckNotNull(in)
	this_.in = in
}

// Java method: LinearFlowState.getOut()
func (this_ *DataFlowAnalysis_LinearFlowState) GetOut() any {
	return this_.out
}

// Java method: LinearFlowState.setOut(L out)
func (this_ *DataFlowAnalysis_LinearFlowState) priv_LinearFlowState_setOut_2339(out any) {
	_ = out
	j2gort.CheckNotNull(out)
	this_.out = out
}

// Java method: LinearFlowState.toString()
func (this_ *DataFlowAnalysis_LinearFlowState) String() string {
	return j2gort.Ext_java_lang_String_Format_2("IN: %s OUT: %s", []any{this_.in, this_.out})
}

// Java class:
type DataFlowAnalysis__Anon_21 struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	cap_0_jsScope *Scope
	cap_1_escaped j2gort.Ext_java_util_Set
}

func (this_ *DataFlowAnalysis__Anon_21) BindSelf_(s any) {
	_ = s
}
func (this_ *DataFlowAnalysis__Anon_21) IsJavaNil_() bool { return this_ == nil }
func (*DataFlowAnalysis__Anon_21) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DataFlowAnalysis$1", ""
}
func (this_ *DataFlowAnalysis__Anon_21) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*DataFlowAnalysis__Anon_21) HasIdentityEquals_()                                        {}
func (*DataFlowAnalysis__Anon_21) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*DataFlowAnalysis__Anon_21) Implements_jscomp_NodeTraversal_Callback()                  {}

func (this_ *DataFlowAnalysis__Anon_21) initFields_DataFlowAnalysis__Anon_21() {
	_ = this_
}

func NewDataFlowAnalysis__Anon_21(cap_p_jsScope *Scope, cap_p_escaped j2gort.Ext_java_util_Set) *DataFlowAnalysis__Anon_21 {
	this_ := &DataFlowAnalysis__Anon_21{}
	this_.cap_0_jsScope = cap_p_jsScope
	this_.cap_1_escaped = cap_p_escaped
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_DataFlowAnalysis__Anon_21()
	return this_
}

// Java method: .visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *DataFlowAnalysis__Anon_21) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if !(n.IsName()) || parent.IsFunction() {
		return
	}
	var name string = n.GetString()
	_ = name
	var var_ *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(name))
	_ = var_
	if var_ == nil {
		return
	}
	var variableCfgScope *Scope = j2gort.CastPtr[Scope](j2gort.CastPtr[Scope](var_.GetScope()).GetClosestCfgRootScope())
	_ = variableCfgScope
	if variableCfgScope != this_.cap_0_jsScope {
		return
	}
	var referenceCfgScope *Scope = j2gort.CastPtr[Scope](t.GetScope().GetClosestCfgRootScope())
	_ = referenceCfgScope
	if (referenceCfgScope != variableCfgScope) || DataFlowAnalysis_IsReferencedInNonStaticClassField(t.GetScope(), variableCfgScope) {
		this_.cap_1_escaped.Add_1(j2gort.PtrToAny(var_))
	}
}
func (*DataFlowAnalysis__Anon_21) AlwaysShouldTraverse_() {}
