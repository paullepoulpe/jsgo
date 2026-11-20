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
// Port of com/google/javascript/jscomp/ControlFlowAnalysis.java (generated mechanically by j2go).

package jscomp

import (
	"sync"
	"sync/atomic"

	pkg_graph "github.com/paullepoulpe/jsgo/transpiled/gen/graph"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_control_flow_analysis_1 j2gort.Ext_java_util_Comparator
var lambdaSingleton_control_flow_analysis_2 j2gort.Ext_java_util_Comparator

func init() {
	lambdaSingleton_control_flow_analysis_1 = j2gort.NewLambda_Ext_java_util_Comparator(func(lp_100 any, lp_101 any) int32 {
		a := j2gort.AsIface[pkg_graph.DiGraph_DiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch]](lp_100)
		_ = a
		b := j2gort.AsIface[pkg_graph.DiGraph_DiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch]](lp_101)
		_ = b
		var pa int32 = a.GetPriority()
		_ = pa
		var pb int32 = b.GetPriority()
		_ = pb
		if pa < pb {
			return int32(-1)
		} else {
			{

				if pa == pb {
					return int32(0)
				}
				return int32(1)
			}

		}
	})
	lambdaSingleton_control_flow_analysis_2 = j2gort.NewLambda_Ext_java_util_Comparator(func(lp_102 any, lp_103 any) int32 {
		a := j2gort.AsIface[pkg_graph.DiGraph_DiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch]](lp_102)
		_ = a
		b := j2gort.AsIface[pkg_graph.DiGraph_DiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch]](lp_103)
		_ = b
		var pa int32 = a.GetPriority()
		_ = pa
		var pb int32 = b.GetPriority()
		_ = pb
		if pb < pa {
			return int32(-1)
		} else {
			{

				if pb == pa {
					return int32(0)
				}
				return int32(1)
			}

		}
	})
}

// Java class: com.google.javascript.jscomp.ControlFlowAnalysis
type ControlFlowAnalysis struct {
	compiler                AbstractCompiler                                                         `gson:"compiler"`
	cfg                     ControlFlowGraph[pkg_rhino.Node]                                         `gson:"cfg"`
	astPosKeys              []pkg_rhino.Node                                                         `gson:"astPosKeys"`
	astPosValues            []int32                                                                  `gson:"astPosValues"`
	astPosEpochs            []int32                                                                  `gson:"astPosEpochs"`
	workHeap                []pkg_graph.DiGraph_DiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch] `gson:"workHeap"`
	exceptionHandler        j2gort.Deque[pkg_rhino.Node]                                             `gson:"exceptionHandler"`
	finallyMap              j2gort.Ext_com_google_common_collect_SetMultimap                         `gson:"finallyMap"`
	astPosEpoch             int32                                                                    `gson:"astPosEpoch"`
	astPosSize              int32                                                                    `gson:"astPosSize"`
	nullAstPos              int32                                                                    `gson:"nullAstPos"`
	workHeapMaxSize         int32                                                                    `gson:"workHeapMaxSize"`
	astPositionCounter      int32                                                                    `gson:"astPositionCounter"`
	priorityCounter         int32                                                                    `gson:"priorityCounter"`
	root                    pkg_rhino.Node                                                           `gson:"root"`
	idHash_                 atomic.Int32
	shouldTraverseFunctions bool `gson:"shouldTraverseFunctions"`
	edgeAnnotations         bool `gson:"edgeAnnotations"`
	workHeapDirty           bool `gson:"workHeapDirty"`
}

func (this_ *ControlFlowAnalysis) BindSelf_(s any) {
	_ = s
}
func (this_ *ControlFlowAnalysis) IsJavaNil_() bool { return this_ == nil }
func (*ControlFlowAnalysis) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ControlFlowAnalysis", "ControlFlowAnalysis"
}
func (this_ *ControlFlowAnalysis) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ControlFlowAnalysis) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ControlFlowAnalysis", "com.google.javascript.jscomp.ControlFlowAnalysis", func(v_ any) bool {
		return j2gort.InstanceOf[*ControlFlowAnalysis](v_)
	})
}
func (*ControlFlowAnalysis) Implements_jscomp_NodeTraversal_Callback() {}

var staticSyncMu_ControlFlowAnalysis sync.Mutex

func StaticGet_ControlFlowAnalysis_MAX_CFA_POOL_SIZE() int32             { return int32(16) }
func StaticGet_ControlFlowAnalysis_MAX_CFG_POOL_SIZE() int32             { return int32(16) }
func StaticGet_ControlFlowAnalysis_MAX_POOLED_AST_POS_CAPACITY() int32   { return int32(4096) }
func StaticGet_ControlFlowAnalysis_MAX_POOLED_WORK_HEAP_CAPACITY() int32 { return int32(2048) }

var static_ControlFlowAnalysis_CFA_POOL []*ControlFlowAnalysis

func StaticGet_ControlFlowAnalysis_CFA_POOL() []*ControlFlowAnalysis {
	clinit_ControlFlowAnalysis()
	return static_ControlFlowAnalysis_CFA_POOL
}
func StaticSet_ControlFlowAnalysis_CFA_POOL(v []*ControlFlowAnalysis) {
	clinit_ControlFlowAnalysis()
	static_ControlFlowAnalysis_CFA_POOL = v
}

var static_ControlFlowAnalysis_cfaPoolSize int32

func StaticGet_ControlFlowAnalysis_CfaPoolSize() int32 {
	clinit_ControlFlowAnalysis()
	return static_ControlFlowAnalysis_cfaPoolSize
}
func StaticSet_ControlFlowAnalysis_CfaPoolSize(v int32) {
	clinit_ControlFlowAnalysis()
	static_ControlFlowAnalysis_cfaPoolSize = v
}

var static_ControlFlowAnalysis_CFG_POOL []*ControlFlowAnalysis_AstControlFlowGraph

func StaticGet_ControlFlowAnalysis_CFG_POOL() []*ControlFlowAnalysis_AstControlFlowGraph {
	clinit_ControlFlowAnalysis()
	return static_ControlFlowAnalysis_CFG_POOL
}
func StaticSet_ControlFlowAnalysis_CFG_POOL(v []*ControlFlowAnalysis_AstControlFlowGraph) {
	clinit_ControlFlowAnalysis()
	static_ControlFlowAnalysis_CFG_POOL = v
}

var static_ControlFlowAnalysis_cfgPoolSize int32

func StaticGet_ControlFlowAnalysis_CfgPoolSize() int32 {
	clinit_ControlFlowAnalysis()
	return static_ControlFlowAnalysis_cfgPoolSize
}
func StaticSet_ControlFlowAnalysis_CfgPoolSize(v int32) {
	clinit_ControlFlowAnalysis()
	static_ControlFlowAnalysis_cfgPoolSize = v
}

var clinitDone_ControlFlowAnalysis uint32
var clinitState_ControlFlowAnalysis j2gort.ClinitState

func clinit_ControlFlowAnalysis() {
	if atomic.LoadUint32(&clinitDone_ControlFlowAnalysis) != 0 {
		return
	}
	clinitSlow_ControlFlowAnalysis()
}

func clinitSlow_ControlFlowAnalysis() {
	if !j2gort.EnterClinit(&clinitState_ControlFlowAnalysis, &clinitDone_ControlFlowAnalysis) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ControlFlowAnalysis, &clinitDone_ControlFlowAnalysis)
	static_ControlFlowAnalysis_CFA_POOL = make([]*ControlFlowAnalysis, int32(16))
	static_ControlFlowAnalysis_CFG_POOL = make([]*ControlFlowAnalysis_AstControlFlowGraph, int32(16))
}

func (this_ *ControlFlowAnalysis) initFields_ControlFlowAnalysis() {
	_ = this_
	this_.astPosEpoch = int32(1)
	this_.nullAstPos = int32(-1)
	this_.exceptionHandler = j2gort.Deque[pkg_rhino.Node]{}
}

// Java constructor: ControlFlowAnalysis(com.google.javascript.jscomp.AbstractCompiler compiler, boolean shouldTraverseFunctions, boolean edgeAnnotations)
func NewControlFlowAnalysis(compiler AbstractCompiler, shouldTraverseFunctions bool, edgeAnnotations bool) *ControlFlowAnalysis {
	clinit_ControlFlowAnalysis()
	this_ := &ControlFlowAnalysis{}
	InitControlFlowAnalysis(this_, compiler, shouldTraverseFunctions, edgeAnnotations)
	return this_
}

func InitControlFlowAnalysis(this_ *ControlFlowAnalysis, compiler AbstractCompiler, shouldTraverseFunctions bool, edgeAnnotations bool) {
	_ = compiler
	_ = shouldTraverseFunctions
	_ = edgeAnnotations
	this_.initFields_ControlFlowAnalysis()
	this_.compiler = compiler
	this_.shouldTraverseFunctions = shouldTraverseFunctions
	this_.edgeAnnotations = edgeAnnotations
}

// Java method: ControlFlowAnalysis.putFinallyEdge(com.google.javascript.rhino.Node from, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node to)
func (this_ *ControlFlowAnalysis) priv_ControlFlowAnalysis_putFinallyEdge_2180(from pkg_rhino.Node, to pkg_rhino.Node) {
	_ = from
	_ = to
	if this_.finallyMap == nil {
		this_.finallyMap = j2gort.AsIface[j2gort.Ext_com_google_common_collect_SetMultimap](j2gort.Ext_com_google_common_collect_LinkedHashMultimap_Create_0())
	}
	j2gort.MultimapPutNode(this_.finallyMap, from, j2gort.BoxNode(to))
}

// Java static method: ControlFlowAnalysis.computeCfgDirect(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node cfgRoot, boolean shouldTraverseFunctions, boolean edgeAnnotations)
func ControlFlowAnalysis_ComputeCfgDirect(compiler AbstractCompiler, cfgRoot pkg_rhino.Node, shouldTraverseFunctions bool, edgeAnnotations bool) ControlFlowGraph[pkg_rhino.Node] {
	clinit_ControlFlowAnalysis()
	_ = compiler
	_ = cfgRoot
	_ = shouldTraverseFunctions
	_ = edgeAnnotations
	var cfa *ControlFlowAnalysis = ControlFlowAnalysis_AcquireCfa(compiler, shouldTraverseFunctions, edgeAnnotations)
	_ = cfa
	defer func() {
		ControlFlowAnalysis_ReleaseCfa(cfa)
	}()
	{
		cfa.priv_ControlFlowAnalysis_computeCfg_2186(cfgRoot)
		return cfa.cfg
	}
	panic("unreachable")
}

// Java static method: ControlFlowAnalysis.popPooledAstCfg()
func ControlFlowAnalysis_PopPooledAstCfg() *ControlFlowAnalysis_AstControlFlowGraph {
	clinit_ControlFlowAnalysis()
	staticSyncMu_ControlFlowAnalysis.Lock()
	defer staticSyncMu_ControlFlowAnalysis.Unlock()
	if StaticGet_ControlFlowAnalysis_CfgPoolSize() > int32(0) {
		var idx int32
		{
			old_1 := StaticGet_ControlFlowAnalysis_CfgPoolSize()
			nw_2 := old_1 - 1
			StaticSet_ControlFlowAnalysis_CfgPoolSize(nw_2)
			idx = nw_2
		}

		_ = idx
		var pooled *ControlFlowAnalysis_AstControlFlowGraph = StaticGet_ControlFlowAnalysis_CFG_POOL()[idx]
		_ = pooled
		{
			a_3 := StaticGet_ControlFlowAnalysis_CFG_POOL()
			i_4 := idx
			a_3[i_4] = nil
		}
		return pooled
	}
	return nil
}

// Java static method: ControlFlowAnalysis.pushPooledAstCfg(com.google.javascript.jscomp.ControlFlowAnalysis.AstControlFlowGraph astCfg)
func ControlFlowAnalysis_PushPooledAstCfg(astCfg *ControlFlowAnalysis_AstControlFlowGraph) bool {
	clinit_ControlFlowAnalysis()
	staticSyncMu_ControlFlowAnalysis.Lock()
	defer staticSyncMu_ControlFlowAnalysis.Unlock()
	_ = astCfg
	if StaticGet_ControlFlowAnalysis_CfgPoolSize() < int32(16) {
		{
			a_5 := StaticGet_ControlFlowAnalysis_CFG_POOL()
			var i_6 int32
			{

				old_7 := StaticGet_ControlFlowAnalysis_CfgPoolSize()
				nw_8 := old_7 + 1
				StaticSet_ControlFlowAnalysis_CfgPoolSize(nw_8)
				i_6 = old_7
			}

			a_5[i_6] = astCfg
		}
		return true
	}
	return false
}

// Java static method: ControlFlowAnalysis.acquireAstCfg(com.google.javascript.rhino.Node entry, boolean edgeAnnotations)
func ControlFlowAnalysis_AcquireAstCfg(entry pkg_rhino.Node, edgeAnnotations bool) *ControlFlowAnalysis_AstControlFlowGraph {
	clinit_ControlFlowAnalysis()
	_ = entry
	_ = edgeAnnotations
	var pooled *ControlFlowAnalysis_AstControlFlowGraph = ControlFlowAnalysis_PopPooledAstCfg()
	_ = pooled
	if pooled != nil {
		pooled.ReinitCfg(entry, edgeAnnotations)
		return pooled
	}
	return NewControlFlowAnalysis_AstControlFlowGraph(entry, edgeAnnotations)
}

// Java static method: ControlFlowAnalysis.releaseCfg(java.lang.@org.jspecify.annotations.Nullable Object maybeCfg)
func ControlFlowAnalysis_ReleaseCfg(maybeCfg any) {
	clinit_ControlFlowAnalysis()
	_ = maybeCfg
	var astCfg *ControlFlowAnalysis_AstControlFlowGraph
	_ = astCfg
	var ok_9 bool
	astCfg, ok_9 = j2gort.TryCast[*ControlFlowAnalysis_AstControlFlowGraph](maybeCfg)
	if ok_9 {

		var pooled bool = (astCfg.ResetCfgForPool() && ControlFlowAnalysis_PushPooledAstCfg(astCfg))
		_ = pooled
		if !(pooled) {
			astCfg.ReleaseSlabs_(

			// Java static method: ControlFlowAnalysis.acquireCfa(com.google.javascript.jscomp.AbstractCompiler compiler, boolean shouldTraverseFunctions, boolean edgeAnnotations)
			)
		}
	}
}

func ControlFlowAnalysis_AcquireCfa(compiler AbstractCompiler, shouldTraverseFunctions bool, edgeAnnotations bool) *ControlFlowAnalysis {
	clinit_ControlFlowAnalysis()
	staticSyncMu_ControlFlowAnalysis.Lock()
	defer staticSyncMu_ControlFlowAnalysis.Unlock()
	_ = compiler
	_ = shouldTraverseFunctions
	_ = edgeAnnotations
	if StaticGet_ControlFlowAnalysis_CfaPoolSize() > int32(0) {
		var idx int32
		{
			old_10 := StaticGet_ControlFlowAnalysis_CfaPoolSize()
			nw_11 := old_10 - 1
			StaticSet_ControlFlowAnalysis_CfaPoolSize(nw_11)
			idx = nw_11
		}

		_ = idx
		var cfa *ControlFlowAnalysis = StaticGet_ControlFlowAnalysis_CFA_POOL()[idx]
		_ = cfa
		{
			a_12 := StaticGet_ControlFlowAnalysis_CFA_POOL()
			i_13 := idx
			a_12[i_13] = nil
		}
		{
			r_14 := cfa
			r_14.compiler = compiler
		}
		{
			r_15 := cfa
			r_15.shouldTraverseFunctions = shouldTraverseFunctions
		}
		{
			r_16 := cfa
			r_16.edgeAnnotations = edgeAnnotations
		}
		return cfa
	}
	return NewControlFlowAnalysis(compiler, shouldTraverseFunctions, edgeAnnotations)
}

// Java static method: ControlFlowAnalysis.releaseCfa(com.google.javascript.jscomp.ControlFlowAnalysis cfa)
func ControlFlowAnalysis_ReleaseCfa(cfa *ControlFlowAnalysis) {
	clinit_ControlFlowAnalysis()
	staticSyncMu_ControlFlowAnalysis.Lock()
	defer staticSyncMu_ControlFlowAnalysis.Unlock()
	_ = cfa
	{
		r_17 := cfa
		r_17.compiler = nil
	}
	{
		r_18 := cfa
		r_18.cfg = nil
	}
	{
		r_19 := cfa
		r_19.root = pkg_rhino.Node(0)
	}
	{
		r_20 := cfa
		r_20.finallyMap = nil
	}
	if !(cfa.exceptionHandler.IsEmpty()) {
		cfa.exceptionHandler.Clear()
	}
	var keys []pkg_rhino.Node = cfa.astPosKeys
	_ = keys
	if keys != nil {
		if int32(len(keys)) > int32(4096) {
			{
				r_21 := cfa
				r_21.astPosKeys = nil
			}
			{
				r_22 := cfa
				r_22.astPosValues = nil
			}
			{
				r_23 := cfa
				r_23.astPosEpochs = nil
			}
			{
				r_24 := cfa
				r_24.astPosEpoch = int32(1)
			}
		} else if cfa.astPosSize > int32(0) {
			r_25 := cfa
			r_25.astPosEpoch = r_25.astPosEpoch + 1
			if cfa.astPosEpoch <= int32(0) {
				{
					r_26 := cfa
					r_26.astPosEpoch = int32(1)
				}
				if cfa.astPosEpochs != nil {
					j2gort.ArraysFill(cfa.astPosEpochs, int32(0))
				}
			}
		}
	}
	{
		r_27 := cfa
		r_27.astPosSize = int32(0)
	}
	{
		r_28 := cfa
		r_28.nullAstPos = int32(-1)
	}
	{
		r_29 := cfa
		r_29.astPositionCounter = int32(0)
	}
	{
		r_30 := cfa
		r_30.priorityCounter = int32(0)
	}
	var heap []pkg_graph.DiGraph_DiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch] = cfa.workHeap
	_ = heap
	if heap != nil {
		if int32(len(heap)) > int32(2048) {
			{
				r_31 := cfa
				r_31.workHeap = nil
			}
		} else if cfa.workHeapDirty {
			var clearLimit int32 = min(cfa.workHeapMaxSize, int32(len(heap)))
			_ = clearLimit
			{
				var i int32 = int32(0)
				_ = i
				for ; i < clearLimit; i = i + 1 {
					{
						a_32 := heap
						i_33 := i
						a_32[i_33] = nil
					}
				}
			}
		}
	}
	{
		r_34 := cfa
		r_34.workHeapDirty = false
	}
	{
		r_35 := cfa
		r_35.workHeapMaxSize = int32(0)
	}
	if StaticGet_ControlFlowAnalysis_CfaPoolSize() < int32(16) {
		{
			a_36 := StaticGet_ControlFlowAnalysis_CFA_POOL()
			var i_37 int32
			{

				old_38 := StaticGet_ControlFlowAnalysis_CfaPoolSize()
				nw_39 := old_38 + 1
				StaticSet_ControlFlowAnalysis_CfaPoolSize(nw_39)
				i_37 = old_38
			}

			a_36[i_37] = cfa
		}
	}
}

// Java static method: ControlFlowAnalysis.builder()
func ControlFlowAnalysis_S_Builder() *ControlFlowAnalysis_Builder {
	clinit_ControlFlowAnalysis()
	return NewControlFlowAnalysis_Builder()
}

// Java method: ControlFlowAnalysis.getCfg()
func (this_ *ControlFlowAnalysis) GetCfg() ControlFlowGraph[pkg_rhino.Node] {
	return this_.cfg
}

// Java static method: ControlFlowAnalysis.mixNodeHash(int h)
func ControlFlowAnalysis_MixNodeHash(h int32) int32 {
	clinit_ControlFlowAnalysis()
	_ = h
	h = (h * int32(-1640531527))
	return (h ^ (javalang.Ushr32(h, int(int32(16)))))
}

// Java method: ControlFlowAnalysis.putAstPosition(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node n, int pos)
func (this_ *ControlFlowAnalysis) priv_ControlFlowAnalysis_putAstPosition_2182(n pkg_rhino.Node, pos int32) {
	_ = n
	_ = pos
	if n == 0 {
		this_.nullAstPos = pos
		return
	}
	var keys []pkg_rhino.Node = this_.astPosKeys
	_ = keys
	if keys == nil {
		keys = make([]pkg_rhino.Node, int32(64))
		this_.astPosKeys = keys
		this_.astPosValues = make([]int32, int32(64))
		this_.astPosEpochs = make([]int32, int32(64))
	} else if (this_.astPosSize * int32(4)) >= (int32(len(keys)) * int32(3)) {
		this_.priv_ControlFlowAnalysis_rehashAstPositions_2183((int32(len(keys)) * int32(2)))
		keys = this_.astPosKeys
	}
	var epochs []int32 = this_.astPosEpochs
	_ = epochs
	var ep int32 = this_.astPosEpoch
	_ = ep
	var mask int32 = (int32(len(keys)) - int32(1))
	_ = mask
	var slot int32 = (ControlFlowAnalysis_MixNodeHash(int32(j2gort.CheckNotNullNode(n))) & mask)
	_ = slot
	for {
		if epochs[slot] != ep {
			{
				a_40 := epochs
				i_41 := slot
				a_40[i_41] = ep
			}
			{
				a_42 := keys
				i_43 := slot
				a_42[i_43] = n
			}
			{
				a_44 := this_.astPosValues
				i_45 := slot
				a_44[i_45] = pos
			}
			this_.astPosSize = this_.astPosSize + 1
			return
		}
		var existing pkg_rhino.Node = keys[slot]
		_ = existing
		if existing == 0 {
			{
				a_46 := keys
				i_47 := slot
				a_46[i_47] = n
			}
			{
				a_48 := this_.astPosValues
				i_49 := slot
				a_48[i_49] = pos
			}
			this_.astPosSize = this_.astPosSize + 1
			return
		}
		if existing == n {
			{
				a_50 := this_.astPosValues
				i_51 := slot
				a_50[i_51] = pos
			}
			return
		}
		slot = ((slot + int32(1)) & mask)
	}
}

// Java method: ControlFlowAnalysis.rehashAstPositions(int newCap)
func (this_ *ControlFlowAnalysis) priv_ControlFlowAnalysis_rehashAstPositions_2183(newCap int32) {
	_ = newCap
	var oldKeys []pkg_rhino.Node = this_.astPosKeys
	_ = oldKeys
	var oldValues []int32 = this_.astPosValues
	_ = oldValues
	var oldEpochs []int32 = this_.astPosEpochs
	_ = oldEpochs
	var ep int32 = this_.astPosEpoch
	_ = ep
	var newKeys []pkg_rhino.Node = make([]pkg_rhino.Node, newCap)
	_ = newKeys
	var newValues []int32 = make([]int32, newCap)
	_ = newValues
	var newEpochs []int32 = make([]int32, newCap)
	_ = newEpochs
	var mask int32 = (newCap - int32(1))
	_ = mask
	{
		var i int32 = int32(0)
		_ = i
		var len_ int32 = int32(len(oldKeys))
		_ = len_
		for ; i < len_; i = i + 1 {
			if oldEpochs[i] == ep {
				var k pkg_rhino.Node = oldKeys[i]
				_ = k
				if k != 0 {
					var slot int32 = (ControlFlowAnalysis_MixNodeHash(int32(j2gort.CheckNotNullNode(k))) & mask)
					_ = slot
					for (newEpochs[slot] == ep) && (newKeys[slot] != 0) {
						slot = ((slot + int32(1)) & mask)
					}
					{
						a_52 := newEpochs
						i_53 := slot
						a_52[i_53] = ep
					}
					{
						a_54 := newKeys
						i_55 := slot
						a_54[i_55] = k
					}
					{
						a_56 := newValues
						i_57 := slot
						a_56[i_57] = oldValues[i]
					}
				}
			}
		}
	}
	this_.astPosKeys = newKeys
	this_.astPosValues = newValues
	this_.astPosEpochs = newEpochs
}

// Java method: ControlFlowAnalysis.getAstPosition(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node n)
func (this_ *ControlFlowAnalysis) priv_ControlFlowAnalysis_getAstPosition_2184(n pkg_rhino.Node) int32 {
	_ = n
	if n == 0 {
		if this_.nullAstPos >= int32(0) {
			return this_.nullAstPos
		}
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "null")))
	}
	var keys []pkg_rhino.Node = this_.astPosKeys
	_ = keys
	if keys != nil {
		var epochs []int32 = this_.astPosEpochs
		_ = epochs
		var ep int32 = this_.astPosEpoch
		_ = ep
		var mask int32 = (int32(len(keys)) - int32(1))
		_ = mask
		var slot int32 = (ControlFlowAnalysis_MixNodeHash(int32(j2gort.CheckNotNullNode(n))) & mask)
		_ = slot
		for epochs[slot] == ep {
			var existing pkg_rhino.Node = keys[slot]
			_ = existing
			if existing == 0 {
				break
			}
			if existing == n {
				return this_.astPosValues[slot]
			}
			slot = ((slot + int32(1)) & mask)
		}
	}
	panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", j2gort.StrVal(n))))
}

// Java method: ControlFlowAnalysis.resolveAstPosition(com.google.javascript.jscomp.graph.DiGraph.DiGraphNode<com.google.javascript.rhino.Node,com.google.javascript.jscomp.ControlFlowGraph.Branch> dNode)
func (this_ *ControlFlowAnalysis) priv_ControlFlowAnalysis_resolveAstPosition_2185(dNode pkg_graph.DiGraph_DiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch]) int32 {
	_ = dNode
	var pos int32 = dNode.GetAstPosition()
	_ = pos
	if pos < int32(0) {
		pos = this_.priv_ControlFlowAnalysis_getAstPosition_2184(dNode.GetValue())
		dNode.SetAstPosition(pos)
	}
	return pos
}

// Java method: ControlFlowAnalysis.computeCfg(com.google.javascript.rhino.Node root)
func (this_ *ControlFlowAnalysis) priv_ControlFlowAnalysis_computeCfg_2186(root pkg_rhino.Node) {
	_ = root
	if !(NodeUtil_IsValidCfgRoot(root)) {
		j2gort.CheckArgumentFmt(false, "Unexpected control flow graph root %s", j2gort.BoxNode(root))
	}
	this_.root = root
	this_.astPosSize = int32(0)
	this_.nullAstPos = int32(-1)
	this_.astPositionCounter = int32(0)
	this_.cfg = ToIface_ControlFlowGraph[pkg_rhino.Node](ControlFlowAnalysis_AcquireAstCfg(ControlFlowAnalysis_ComputeFallThrough(root), this_.edgeAnnotations))
	NodeTraversal_Traverse(this_.compiler, root, this_)
	var iifeArg_1 int32
	{

		old_58 := this_.astPositionCounter
		nw_59 := old_58 + 1
		this_.astPositionCounter = nw_59
		iifeArg_1 = nw_59
	}
	this_.priv_ControlFlowAnalysis_putAstPosition_2182(pkg_rhino.Node(0), iifeArg_1)
	this_.priorityCounter = int32(0)
	var entry pkg_graph.DiGraph_DiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch] = this_.cfg.GetEntry()
	_ = entry
	this_.priv_ControlFlowAnalysis_prioritizeFromEntryNode_2187(entry)
	if this_.shouldTraverseFunctions {
		{
			cur_60 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.cfg.GetNodes()))
			for cur_60.HasNext() {
				rawElem_61 := cur_60.Next()
				candidate := j2gort.AsIface[pkg_graph.DiGraph_DiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch]](rawElem_61)
				_ = candidate
				var value pkg_rhino.Node = candidate.GetValue()
				_ = value
				if (value != 0) && value.IsFunction() {
					this_.priv_ControlFlowAnalysis_prioritizeFromEntryNode_2187(candidate)
				}
			}
		}
	}
	{
		cur_62 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.cfg.GetNodes()))
		for cur_62.HasNext() {
			rawElem_63 := cur_62.Next()
			candidate_2 := j2gort.AsIface[pkg_graph.DiGraph_DiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch]](rawElem_63)
			_ = candidate_2
			if !(candidate_2.HasPriority()) {
				candidate_2.SetPriority(func() int32 {
					old_64 := this_.priorityCounter
					nw_65 := old_64 + 1
					this_.priorityCounter = nw_65
					return nw_65
				}())
			}
		}
	}
	this_.cfg.GetImplicitReturn().SetPriority(func() int32 {
		old_66 := this_.priorityCounter
		nw_67 := old_66 + 1
		this_.priorityCounter = nw_67
		return nw_67
	}())
}

// Java method: ControlFlowAnalysis.prioritizeFromEntryNode(com.google.javascript.jscomp.graph.DiGraph.DiGraphNode<com.google.javascript.rhino.Node,com.google.javascript.jscomp.ControlFlowGraph.Branch> entry)
func (this_ *ControlFlowAnalysis) priv_ControlFlowAnalysis_prioritizeFromEntryNode_2187(entry pkg_graph.DiGraph_DiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch]) {
	_ = entry
	var heap []pkg_graph.DiGraph_DiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch] = this_.workHeap
	_ = heap
	if heap == nil {
		heap = make([]pkg_graph.DiGraph_DiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch], int32(32))
		this_.workHeap = heap
	}
	this_.workHeapDirty = true
	{
		a_68 := heap
		i_69 := int32(0)
		a_68[i_69] = entry
	}
	var heapSize int32 = int32(1)
	_ = heapSize
	this_.workHeapMaxSize = int32(1)
	for heapSize > int32(0) {
		var current pkg_graph.DiGraph_DiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch] = heap[int32(0)]
		_ = current
		var n int32
		{
			old_70 := heapSize
			nw_71 := old_70 - 1
			heapSize = nw_71
			n = nw_71
		}
		_ = n
		var x pkg_graph.DiGraph_DiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch] = heap[n]
		_ = x
		{
			a_72 := heap
			i_73 := n
			a_72[i_73] = nil
		}
		if n > int32(0) {
			var half int32 = javalang.Ushr32(n, int(int32(1)))
			_ = half
			var k int32 = int32(0)
			_ = k
			var xPos int32 = this_.priv_ControlFlowAnalysis_resolveAstPosition_2185(x)
			_ = xPos
			for k < half {
				var child int32 = ((javalang.Shl32(k, int(int32(1)))) + int32(1))
				_ = child
				var c pkg_graph.DiGraph_DiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch] = heap[child]
				_ = c
				var right int32 = (child + int32(1))
				_ = right
				var cPos int32 = this_.priv_ControlFlowAnalysis_resolveAstPosition_2185(c)
				_ = cPos
				if right < n {
					var r pkg_graph.DiGraph_DiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch] = heap[right]
					_ = r
					var rPos int32 = this_.priv_ControlFlowAnalysis_resolveAstPosition_2185(r)
					_ = rPos
					if cPos > rPos {
						child = right
						c = r
						cPos = rPos
					}
				}
				if xPos <= cPos {
					break
				}
				{
					a_74 := heap
					i_75 := k
					a_74[i_75] = c
				}
				k = child
			}
			{
				a_76 := heap
				i_77 := k
				a_76[i_77] = x
			}
		}
		if current.HasPriority() {
			continue
		}
		current.SetPriority(func() int32 {
			old_78 := this_.priorityCounter
			nw_79 := old_78 + 1
			this_.priorityCounter = nw_79
			return nw_79
		}())
		{
			var i int32 = int32(0)
			_ = i
			var outCount int32 = current.GetOutEdgeCount()
			_ = outCount
			for ; i < outCount; i = i + 1 {
				var succ pkg_graph.DiGraph_DiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch] = current.GetOutEdgeAt(i).GetDestination()
				_ = succ
				if heapSize == int32(len(heap)) {
					{

						arrSrc_80 := heap
						arrLen_81 := (int32(len(heap)) * int32(2))
						j2gort.CheckNotNullSlice(arrSrc_80)
						if arrLen_81 < 0 {
							panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
						}
						arrOut_82 := make([]pkg_graph.DiGraph_DiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch], arrLen_81)
						copy(arrOut_82, arrSrc_80)
						heap = arrOut_82
					}

					this_.workHeap = heap
				}
				if heapSize == int32(0) {
					{
						a_83 := heap
						i_84 := int32(0)
						a_83[i_84] = succ
					}
					heapSize = int32(1)
				} else {
					var k_2 int32
					{
						old_85 := heapSize
						nw_86 := old_85 + 1
						heapSize = nw_86
						k_2 = old_85
					}
					_ = k_2
					var succPos int32 = this_.priv_ControlFlowAnalysis_resolveAstPosition_2185(succ)
					_ = succPos
					for k_2 > int32(0) {
						var parent int32 = javalang.Ushr32((k_2 - int32(1)), int(int32(1)))
						_ = parent
						var e pkg_graph.DiGraph_DiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch] = heap[parent]
						_ = e
						if succPos >= this_.priv_ControlFlowAnalysis_resolveAstPosition_2185(e) {
							break
						}
						{
							a_87 := heap
							i_88 := k_2
							a_87[i_88] = e
						}
						k_2 = parent
					}
					{
						a_89 := heap
						i_90 := k_2
						a_89[i_90] = succ
					}
				}
				if heapSize > this_.workHeapMaxSize {
					this_.workHeapMaxSize = heapSize
				}
			}
		}
	}
	this_.workHeapDirty = false
	this_.workHeapMaxSize = int32(0)
}

// Java method: ControlFlowAnalysis.shouldTraverse(com.google.javascript.jscomp.NodeTraversal nodeTraversal, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ControlFlowAnalysis) ShouldTraverse_3(nodeTraversal *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = nodeTraversal
	_ = n
	_ = parent
	if this_.priv_ControlFlowAnalysis_shouldTraverseIntoChildren_2188(n, parent) {
		var iifeArg_2 int32
		{

			old_91 := this_.astPositionCounter
			nw_92 := old_91 + 1
			this_.astPositionCounter = nw_92
			iifeArg_2 = old_91
		}
		this_.priv_ControlFlowAnalysis_putAstPosition_2182(n, iifeArg_2)
		return true
	}
	return false
}

// Java method: ControlFlowAnalysis.shouldTraverseIntoChildren(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ControlFlowAnalysis) priv_ControlFlowAnalysis_shouldTraverseIntoChildren_2188(n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(73) /* Token.FUNCTION */ :
		if this_.shouldTraverseFunctions || (n == this_.cfg.GetEntry().GetValue()) {
			this_.exceptionHandler.Push(j2gort.CheckNotNullNode(n))
			return true
		}
		return false
	case int32(48) /* Token.TRY */ :
		this_.exceptionHandler.Push(j2gort.CheckNotNullNode(n))
		return true
	default:
	}
	if parent != 0 {
		switch parent.TokenSwitchOrd_() {
		case int32(80) /* Token.FOR */, int32(81) /* Token.FOR_IN */, int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */ :
			var shouldTraverseForChild bool = (n == parent.GetLastChild())
			_ = shouldTraverseForChild
			if !(shouldTraverseForChild) {
				var iifeArg_3 int32
				{

					old_93 := this_.astPositionCounter
					nw_94 := old_93 + 1
					this_.astPositionCounter = nw_94
					iifeArg_3 = old_93
				}
				this_.priv_ControlFlowAnalysis_putAstPosition_2182(n, iifeArg_3)
			}
			return shouldTraverseForChild
		case int32(79) /* Token.DO */ :
			return (n != parent.GetSecondChild())
		case int32(74) /* Token.IF */, int32(78) /* Token.WHILE */, int32(85) /* Token.WITH */, int32(75) /* Token.SWITCH */, int32(76) /* Token.CASE */, int32(86) /* Token.CATCH */, int32(92) /* Token.LABEL */ :
			return (n != parent.GetFirstChild())
		case int32(73) /* Token.FUNCTION */ :
			return (n == parent.GetLastChild())
		case int32(105) /* Token.CLASS */ :
			return (this_.shouldTraverseFunctions && (n == parent.GetLastChild()))
		case int32(129) /* Token.COMPUTED_PROP */, int32(83) /* Token.CONTINUE */, int32(82) /* Token.BREAK */, int32(93) /* Token.EXPR_RESULT */, int32(84) /* Token.VAR */, int32(111) /* Token.LET */, int32(97) /* Token.CONST */, int32(120) /* Token.EXPORT */, int32(116) /* Token.IMPORT */, int32(0) /* Token.RETURN */, int32(43) /* Token.THROW */, int32(107) /* Token.MEMBER_FUNCTION_DEF */, int32(108) /* Token.MEMBER_FIELD_DEF */, int32(109) /* Token.COMPUTED_FIELD_DEF */ :
			return false
		case int32(48) /* Token.TRY */ :
			if (!(NodeUtil_HasFinally(parent)) && (n == NodeUtil_GetCatchBlock(parent))) || NodeUtil_IsTryFinallyNode(parent, n) {
				if !(this_.exceptionHandler.PeekFirst() == parent) {
					j2gort.CheckState(false)
				}
				this_.exceptionHandler.Pop()
			}
		default:
		}
		if (parent.HasParent() && parent.GetParent().IsArrowFunction()) && !(parent.IsBlock()) {
			return false
		}
	}
	return true
}

// Java method: ControlFlowAnalysis.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ControlFlowAnalysis) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(74) /* Token.IF */ :
		this_.priv_ControlFlowAnalysis_handleIf_2189(n)
	case int32(78) /* Token.WHILE */ :
		this_.priv_ControlFlowAnalysis_handleWhile_2190(n)
	case int32(79) /* Token.DO */ :
		this_.priv_ControlFlowAnalysis_handleDo_2191(n)
	case int32(80) /* Token.FOR */ :
		this_.priv_ControlFlowAnalysis_handleFor_2193(n)
	case int32(112) /* Token.FOR_OF */, int32(81) /* Token.FOR_IN */, int32(113) /* Token.FOR_AWAIT_OF */ :
		this_.priv_ControlFlowAnalysis_handleEnhancedFor_2192(n)
	case int32(75) /* Token.SWITCH */ :
		this_.priv_ControlFlowAnalysis_handleSwitch_2194(n)
	case int32(76) /* Token.CASE */ :
		this_.priv_ControlFlowAnalysis_handleCase_2195(n)
	case int32(77) /* Token.DEFAULT_CASE */ :
		this_.priv_ControlFlowAnalysis_handleDefault_2196(n)
	case int32(90) /* Token.BLOCK */, int32(89) /* Token.ROOT */, int32(94) /* Token.SCRIPT */, int32(123) /* Token.MODULE_BODY */ :
		this_.priv_ControlFlowAnalysis_handleStmtList_2198(n)
	case int32(73) /* Token.FUNCTION */ :
		this_.priv_ControlFlowAnalysis_handleFunction_2199(n)
	case int32(93) /* Token.EXPR_RESULT */ :
		this_.priv_ControlFlowAnalysis_handleExpr_2200(n)
	case int32(43) /* Token.THROW */ :
		this_.priv_ControlFlowAnalysis_handleThrow_2201(n)
	case int32(48) /* Token.TRY */ :
		this_.priv_ControlFlowAnalysis_handleTry_2202(n)
	case int32(86) /* Token.CATCH */ :
		this_.priv_ControlFlowAnalysis_handleCatch_2203(n)
	case int32(82) /* Token.BREAK */ :
		this_.priv_ControlFlowAnalysis_handleBreak_2204(n)
	case int32(83) /* Token.CONTINUE */ :
		this_.priv_ControlFlowAnalysis_handleContinue_2205(n)
	case int32(0) /* Token.RETURN */ :
		this_.priv_ControlFlowAnalysis_handleReturn_2206(n)
	case int32(85) /* Token.WITH */ :
		this_.priv_ControlFlowAnalysis_handleWith_2197(n)
	case int32(91) /* Token.SWITCH_BODY */, int32(92) /* Token.LABEL */, int32(106) /* Token.CLASS_MEMBERS */, int32(107) /* Token.MEMBER_FUNCTION_DEF */, int32(108) /* Token.MEMBER_FIELD_DEF */, int32(109) /* Token.COMPUTED_FIELD_DEF */ :
	default:
		this_.priv_ControlFlowAnalysis_handleStmt_2207(n)
	}
}

// Java method: ControlFlowAnalysis.handleIf(com.google.javascript.rhino.Node node)
func (this_ *ControlFlowAnalysis) priv_ControlFlowAnalysis_handleIf_2189(node pkg_rhino.Node) {
	_ = node
	var thenBlock pkg_rhino.Node = node.GetSecondChild()
	_ = thenBlock
	var elseBlock pkg_rhino.Node = thenBlock.GetNext()
	_ = elseBlock
	this_.priv_ControlFlowAnalysis_createEdge_2208(node, StaticGet_ControlFlowGraph_Branch_ON_TRUE(), ControlFlowAnalysis_ComputeFallThrough(thenBlock))
	if elseBlock == 0 {
		this_.priv_ControlFlowAnalysis_createEdge_2208(node, StaticGet_ControlFlowGraph_Branch_ON_FALSE(), ControlFlowAnalysis_ComputeFollowNode_2(node, this_))
	} else {
		this_.priv_ControlFlowAnalysis_createEdge_2208(node, StaticGet_ControlFlowGraph_Branch_ON_FALSE(), ControlFlowAnalysis_ComputeFallThrough(elseBlock))
	}
	this_.priv_ControlFlowAnalysis_connectToPossibleExceptionHandler_2209(node, NodeUtil_GetConditionExpression(node))
}

// Java method: ControlFlowAnalysis.handleWhile(com.google.javascript.rhino.Node node)
func (this_ *ControlFlowAnalysis) priv_ControlFlowAnalysis_handleWhile_2190(node pkg_rhino.Node) {
	_ = node
	var cond pkg_rhino.Node = node.GetFirstChild()
	_ = cond
	this_.priv_ControlFlowAnalysis_createEdge_2208(node, StaticGet_ControlFlowGraph_Branch_ON_TRUE(), ControlFlowAnalysis_ComputeFallThrough(cond.GetNext()))
	if !(cond.IsTrue()) {
		this_.priv_ControlFlowAnalysis_createEdge_2208(node, StaticGet_ControlFlowGraph_Branch_ON_FALSE(), ControlFlowAnalysis_ComputeFollowNode_2(node, this_))
	}
	this_.priv_ControlFlowAnalysis_connectToPossibleExceptionHandler_2209(node, NodeUtil_GetConditionExpression(node))
}

// Java method: ControlFlowAnalysis.handleDo(com.google.javascript.rhino.Node node)
func (this_ *ControlFlowAnalysis) priv_ControlFlowAnalysis_handleDo_2191(node pkg_rhino.Node) {
	_ = node
	var body pkg_rhino.Node = node.GetFirstChild()
	_ = body
	this_.priv_ControlFlowAnalysis_createEdge_2208(node, StaticGet_ControlFlowGraph_Branch_ON_TRUE(), ControlFlowAnalysis_ComputeFallThrough(body))
	var cond pkg_rhino.Node = body.GetNext()
	_ = cond
	if !(cond.IsTrue()) {
		this_.priv_ControlFlowAnalysis_createEdge_2208(node, StaticGet_ControlFlowGraph_Branch_ON_FALSE(), ControlFlowAnalysis_ComputeFollowNode_2(node, this_))
	}
	this_.priv_ControlFlowAnalysis_connectToPossibleExceptionHandler_2209(node, NodeUtil_GetConditionExpression(node))
}

// Java method: ControlFlowAnalysis.handleEnhancedFor(com.google.javascript.rhino.Node forNode)
func (this_ *ControlFlowAnalysis) priv_ControlFlowAnalysis_handleEnhancedFor_2192(forNode pkg_rhino.Node) {
	_ = forNode
	var item pkg_rhino.Node = forNode.GetFirstChild()
	_ = item
	var collection pkg_rhino.Node = item.GetNext()
	_ = collection
	var body pkg_rhino.Node = collection.GetNext()
	_ = body
	this_.priv_ControlFlowAnalysis_createEdge_2208(collection, StaticGet_ControlFlowGraph_Branch_UNCOND(), forNode)
	this_.priv_ControlFlowAnalysis_createEdge_2208(forNode, StaticGet_ControlFlowGraph_Branch_ON_TRUE(), ControlFlowAnalysis_ComputeFallThrough(body))
	this_.priv_ControlFlowAnalysis_createEdge_2208(forNode, StaticGet_ControlFlowGraph_Branch_ON_FALSE(), ControlFlowAnalysis_ComputeFollowNode_2(forNode, this_))
	if forNode.IsForOf() || forNode.IsForAwaitOf() {
		this_.priv_ControlFlowAnalysis_connectToPossibleExceptionHandlerUnconditionally_2210(forNode)
	} else {
		this_.priv_ControlFlowAnalysis_connectToPossibleExceptionHandler_2209(forNode, collection)
	}
}

// Java method: ControlFlowAnalysis.handleFor(com.google.javascript.rhino.Node forNode)
func (this_ *ControlFlowAnalysis) priv_ControlFlowAnalysis_handleFor_2193(forNode pkg_rhino.Node) {
	_ = forNode
	var init pkg_rhino.Node = forNode.GetFirstChild()
	_ = init
	var cond pkg_rhino.Node = init.GetNext()
	_ = cond
	var iter pkg_rhino.Node = cond.GetNext()
	_ = iter
	var body pkg_rhino.Node = iter.GetNext()
	_ = body
	this_.priv_ControlFlowAnalysis_createEdge_2208(init, StaticGet_ControlFlowGraph_Branch_UNCOND(), forNode)
	this_.priv_ControlFlowAnalysis_createEdge_2208(forNode, StaticGet_ControlFlowGraph_Branch_ON_TRUE(), ControlFlowAnalysis_ComputeFallThrough(body))
	if !(cond.IsEmpty()) && !(cond.IsTrue()) {
		this_.priv_ControlFlowAnalysis_createEdge_2208(forNode, StaticGet_ControlFlowGraph_Branch_ON_FALSE(), ControlFlowAnalysis_ComputeFollowNode_2(forNode, this_))
	}
	this_.priv_ControlFlowAnalysis_createEdge_2208(iter, StaticGet_ControlFlowGraph_Branch_UNCOND(), forNode)
	this_.priv_ControlFlowAnalysis_connectToPossibleExceptionHandler_2209(init, init)
	this_.priv_ControlFlowAnalysis_connectToPossibleExceptionHandler_2209(forNode, cond)
	this_.priv_ControlFlowAnalysis_connectToPossibleExceptionHandler_2209(iter, iter)
}

// Java method: ControlFlowAnalysis.handleSwitch(com.google.javascript.rhino.Node node)
func (this_ *ControlFlowAnalysis) priv_ControlFlowAnalysis_handleSwitch_2194(node pkg_rhino.Node) {
	_ = node
	var switchBody pkg_rhino.Node = node.GetLastChild()
	_ = switchBody
	this_.priv_ControlFlowAnalysis_createEdge_2208(node, StaticGet_ControlFlowGraph_Branch_UNCOND(), switchBody)
	var next pkg_rhino.Node = ControlFlowAnalysis_GetNextSiblingOfType(switchBody.GetFirstChild(), []pkg_rhino.Token{pkg_rhino.StaticGet_Token_CASE(), pkg_rhino.StaticGet_Token_EMPTY()})
	_ = next
	if next != 0 {
		this_.priv_ControlFlowAnalysis_createEdge_2208(switchBody, StaticGet_ControlFlowGraph_Branch_UNCOND(), next)
	} else {
		if switchBody.HasChildren() {
			this_.priv_ControlFlowAnalysis_createEdge_2208(switchBody, StaticGet_ControlFlowGraph_Branch_UNCOND(), switchBody.GetFirstChild())
		} else {
			this_.priv_ControlFlowAnalysis_createEdge_2208(switchBody, StaticGet_ControlFlowGraph_Branch_UNCOND(), ControlFlowAnalysis_ComputeFollowNode_2(node, this_))
		}
	}
	this_.priv_ControlFlowAnalysis_connectToPossibleExceptionHandler_2209(node, node.GetFirstChild())
}

// Java method: ControlFlowAnalysis.handleCase(com.google.javascript.rhino.Node caseNode)
func (this_ *ControlFlowAnalysis) priv_ControlFlowAnalysis_handleCase_2195(caseNode pkg_rhino.Node) {
	_ = caseNode
	this_.priv_ControlFlowAnalysis_createEdge_2208(caseNode, StaticGet_ControlFlowGraph_Branch_ON_TRUE(), caseNode.GetSecondChild())
	var next pkg_rhino.Node = ControlFlowAnalysis_GetNextSiblingOfType(caseNode.GetNext(), []pkg_rhino.Token{pkg_rhino.StaticGet_Token_CASE()})
	_ = next
	if next != 0 {
		if !(next.IsCase()) {
			j2gort.CheckState(false)
		}
		this_.priv_ControlFlowAnalysis_createEdge_2208(caseNode, StaticGet_ControlFlowGraph_Branch_ON_FALSE(), next)
	} else {
		var switchBody pkg_rhino.Node = caseNode.GetParent()
		_ = switchBody
		var deflt pkg_rhino.Node = ControlFlowAnalysis_GetNextSiblingOfType(switchBody.GetFirstChild(), []pkg_rhino.Token{pkg_rhino.StaticGet_Token_DEFAULT_CASE()})
		_ = deflt
		if deflt != 0 {
			this_.priv_ControlFlowAnalysis_createEdge_2208(caseNode, StaticGet_ControlFlowGraph_Branch_ON_FALSE(), deflt)
		} else {
			this_.priv_ControlFlowAnalysis_createEdge_2208(caseNode, StaticGet_ControlFlowGraph_Branch_ON_FALSE(), ControlFlowAnalysis_ComputeFollowNode_2(caseNode, this_))
		}
	}
	this_.priv_ControlFlowAnalysis_connectToPossibleExceptionHandler_2209(caseNode, caseNode.GetFirstChild())
}

// Java method: ControlFlowAnalysis.handleDefault(com.google.javascript.rhino.Node node)
func (this_ *ControlFlowAnalysis) priv_ControlFlowAnalysis_handleDefault_2196(node pkg_rhino.Node) {
	_ = node
	this_.priv_ControlFlowAnalysis_createEdge_2208(node, StaticGet_ControlFlowGraph_Branch_UNCOND(), node.GetFirstChild())
}

// Java method: ControlFlowAnalysis.handleWith(com.google.javascript.rhino.Node node)
func (this_ *ControlFlowAnalysis) priv_ControlFlowAnalysis_handleWith_2197(node pkg_rhino.Node) {
	_ = node
	this_.priv_ControlFlowAnalysis_createEdge_2208(node, StaticGet_ControlFlowGraph_Branch_UNCOND(), node.GetLastChild())
	this_.priv_ControlFlowAnalysis_connectToPossibleExceptionHandler_2209(node, node.GetFirstChild())
}

// Java method: ControlFlowAnalysis.handleStmtList(com.google.javascript.rhino.Node node)
func (this_ *ControlFlowAnalysis) priv_ControlFlowAnalysis_handleStmtList_2198(node pkg_rhino.Node) {
	_ = node
	var parent pkg_rhino.Node = node.GetParent()
	_ = parent
	if ((node.IsBlock() && parent.IsTry()) && (NodeUtil_GetCatchBlock(parent) == node)) && !(NodeUtil_HasCatchHandler(node)) {
		return
	}
	var child pkg_rhino.Node = node.GetFirstChild()
	_ = child
	for (child != 0) && child.IsFunction() {
		child = child.GetNext()
	}
	if child != 0 {
		this_.priv_ControlFlowAnalysis_createEdge_2208(node, StaticGet_ControlFlowGraph_Branch_UNCOND(), ControlFlowAnalysis_ComputeFallThrough(child))
	} else {
		this_.priv_ControlFlowAnalysis_createEdge_2208(node, StaticGet_ControlFlowGraph_Branch_UNCOND(), ControlFlowAnalysis_ComputeFollowNode_2(node, this_))
	}
	if parent != 0 {
		switch parent.TokenSwitchOrd_() {
		case int32(77) /* Token.DEFAULT_CASE */, int32(76) /* Token.CASE */, int32(48) /* Token.TRY */ :
		case int32(89) /* Token.ROOT */ :
			if node.IsRoot() && (node.GetNext() != 0) {
				this_.priv_ControlFlowAnalysis_createEdge_2208(node, StaticGet_ControlFlowGraph_Branch_UNCOND(), node.GetNext())
			}
		default:
			if node.IsBlock() && node.IsSyntheticBlock() {
				this_.priv_ControlFlowAnalysis_createEdge_2208(node, StaticGet_ControlFlowGraph_Branch_SYN_BLOCK(), ControlFlowAnalysis_ComputeFollowNode_2(node, this_))
			}
		}
	}
}

// Java method: ControlFlowAnalysis.handleFunction(com.google.javascript.rhino.Node node)
func (this_ *ControlFlowAnalysis) priv_ControlFlowAnalysis_handleFunction_2199(node pkg_rhino.Node) {
	_ = node
	if !(node.IsFunction()) {
		j2gort.CheckState(false)
	}
	if !(node.HasXChildren(int32(3))) {
		j2gort.CheckState(false)
	}
	this_.priv_ControlFlowAnalysis_createEdge_2208(node, StaticGet_ControlFlowGraph_Branch_UNCOND(), ControlFlowAnalysis_ComputeFallThrough(node.GetLastChild()))
	if !(this_.exceptionHandler.PeekFirst() == node) {
		j2gort.CheckState(false)
	}
	this_.exceptionHandler.Pop()
}

// Java method: ControlFlowAnalysis.handleExpr(com.google.javascript.rhino.Node node)
func (this_ *ControlFlowAnalysis) priv_ControlFlowAnalysis_handleExpr_2200(node pkg_rhino.Node) {
	_ = node
	this_.priv_ControlFlowAnalysis_createEdge_2208(node, StaticGet_ControlFlowGraph_Branch_UNCOND(), ControlFlowAnalysis_ComputeFollowNode_2(node, this_))
	this_.priv_ControlFlowAnalysis_connectToPossibleExceptionHandler_2209(node, node)
}

// Java method: ControlFlowAnalysis.handleThrow(com.google.javascript.rhino.Node node)
func (this_ *ControlFlowAnalysis) priv_ControlFlowAnalysis_handleThrow_2201(node pkg_rhino.Node) {
	_ = node
	this_.priv_ControlFlowAnalysis_connectToPossibleExceptionHandler_2209(node, node)
}

// Java method: ControlFlowAnalysis.handleTry(com.google.javascript.rhino.Node node)
func (this_ *ControlFlowAnalysis) priv_ControlFlowAnalysis_handleTry_2202(node pkg_rhino.Node) {
	_ = node
	this_.priv_ControlFlowAnalysis_createEdge_2208(node, StaticGet_ControlFlowGraph_Branch_UNCOND(), node.GetFirstChild())
}

// Java method: ControlFlowAnalysis.handleCatch(com.google.javascript.rhino.Node node)
func (this_ *ControlFlowAnalysis) priv_ControlFlowAnalysis_handleCatch_2203(node pkg_rhino.Node) {
	_ = node
	this_.priv_ControlFlowAnalysis_createEdge_2208(node, StaticGet_ControlFlowGraph_Branch_UNCOND(), node.GetLastChild())
}

// Java method: ControlFlowAnalysis.handleBreak(com.google.javascript.rhino.Node node)
func (this_ *ControlFlowAnalysis) priv_ControlFlowAnalysis_handleBreak_2204(node pkg_rhino.Node) {
	_ = node
	var label string = j2gort.NullString
	_ = label
	if node.HasChildren() {
		label = node.GetFirstChild().GetString()
	}
	var cur pkg_rhino.Node = pkg_rhino.Node(0)
	_ = cur
	var previous pkg_rhino.Node = pkg_rhino.Node(0)
	_ = previous
	var lastJump pkg_rhino.Node = pkg_rhino.Node(0)
	_ = lastJump
	var parent pkg_rhino.Node = node.GetParent()
	_ = parent
	{
		cur = node
		lastJump = node
		for ; !(ControlFlowAnalysis_IsBreakTarget(cur, label)); func() { cur = parent; parent = parent.GetParent() }() {
			if (cur.IsTry() && NodeUtil_HasFinally(cur)) && (cur.GetLastChild() != previous) {
				if lastJump == node {
					this_.priv_ControlFlowAnalysis_createEdge_2208(lastJump, StaticGet_ControlFlowGraph_Branch_UNCOND(), ControlFlowAnalysis_ComputeFallThrough(cur.GetLastChild()))
				} else {
					this_.priv_ControlFlowAnalysis_putFinallyEdge_2180(lastJump, ControlFlowAnalysis_ComputeFallThrough(cur.GetLastChild()))
				}
				lastJump = cur
			}
			if parent == 0 {
				if this_.compiler.GetOptions().CanContinueAfterErrors() {
					return
				} else {
					panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Cannot find break target.")))
				}
			}
			previous = cur
		}
	}
	if lastJump == node {
		this_.priv_ControlFlowAnalysis_createEdge_2208(lastJump, StaticGet_ControlFlowGraph_Branch_UNCOND(), ControlFlowAnalysis_ComputeFollowNode_2(cur, this_))
	} else {
		this_.priv_ControlFlowAnalysis_putFinallyEdge_2180(lastJump, ControlFlowAnalysis_ComputeFollowNode_2(cur, this_))
	}
}

// Java method: ControlFlowAnalysis.handleContinue(com.google.javascript.rhino.Node node)
func (this_ *ControlFlowAnalysis) priv_ControlFlowAnalysis_handleContinue_2205(node pkg_rhino.Node) {
	_ = node
	var label string = j2gort.NullString
	_ = label
	if node.HasChildren() {
		label = node.GetFirstChild().GetString()
	}
	var cur pkg_rhino.Node = pkg_rhino.Node(0)
	_ = cur
	var previous pkg_rhino.Node = pkg_rhino.Node(0)
	_ = previous
	var lastJump pkg_rhino.Node = pkg_rhino.Node(0)
	_ = lastJump
	{
		cur = node
		lastJump = node
		for ; !(ControlFlowAnalysis_IsContinueTarget(cur, label)); cur = cur.GetParent() {
			if (cur.IsTry() && NodeUtil_HasFinally(cur)) && (cur.GetLastChild() != previous) {
				if lastJump == node {
					this_.priv_ControlFlowAnalysis_createEdge_2208(lastJump, StaticGet_ControlFlowGraph_Branch_UNCOND(), cur.GetLastChild())
				} else {
					this_.priv_ControlFlowAnalysis_putFinallyEdge_2180(lastJump, ControlFlowAnalysis_ComputeFallThrough(cur.GetLastChild()))
				}
				lastJump = cur
			}
			if !(cur.HasParent()) {
				j2gort.CheckStateStr(false, "Cannot find continue target.")
			}
			previous = cur
		}
	}
	var iter pkg_rhino.Node = cur
	_ = iter
	if cur.IsVanillaFor() {
		iter = cur.GetChildAtIndex(int32(2))
	}
	if lastJump == node {
		this_.priv_ControlFlowAnalysis_createEdge_2208(node, StaticGet_ControlFlowGraph_Branch_UNCOND(), iter)
	} else {
		this_.priv_ControlFlowAnalysis_putFinallyEdge_2180(lastJump, iter)
	}
}

// Java method: ControlFlowAnalysis.handleReturn(com.google.javascript.rhino.Node node)
func (this_ *ControlFlowAnalysis) priv_ControlFlowAnalysis_handleReturn_2206(node pkg_rhino.Node) {
	_ = node
	var lastJump pkg_rhino.Node = pkg_rhino.Node(0)
	_ = lastJump
	{
		cur_95 := this_.exceptionHandler.Iter()
		for cur_95.HasNext() {
			curHandler := cur_95.Next()
			_ = curHandler
			if curHandler.IsFunction() {
				break
			}
			if NodeUtil_HasFinally(curHandler) {
				if lastJump == 0 {
					this_.priv_ControlFlowAnalysis_createEdge_2208(node, StaticGet_ControlFlowGraph_Branch_UNCOND(), curHandler.GetLastChild())
				} else {
					this_.priv_ControlFlowAnalysis_putFinallyEdge_2180(lastJump, ControlFlowAnalysis_ComputeFallThrough(curHandler.GetLastChild()))
				}
				lastJump = curHandler
			}
		}
	}
	if node.HasChildren() {
		this_.priv_ControlFlowAnalysis_connectToPossibleExceptionHandler_2209(node, node.GetFirstChild())
	}
	if lastJump == 0 {
		this_.priv_ControlFlowAnalysis_createEdge_2208(node, StaticGet_ControlFlowGraph_Branch_UNCOND(), pkg_rhino.Node(0))
	} else {
		this_.priv_ControlFlowAnalysis_putFinallyEdge_2180(lastJump, pkg_rhino.Node(0))
	}
}

// Java method: ControlFlowAnalysis.handleStmt(com.google.javascript.rhino.Node node)
func (this_ *ControlFlowAnalysis) priv_ControlFlowAnalysis_handleStmt_2207(node pkg_rhino.Node) {
	_ = node
	this_.priv_ControlFlowAnalysis_createEdge_2208(node, StaticGet_ControlFlowGraph_Branch_UNCOND(), ControlFlowAnalysis_ComputeFollowNode_2(node, this_))
	this_.priv_ControlFlowAnalysis_connectToPossibleExceptionHandler_2209(node, node)
}

// Java static method: ControlFlowAnalysis.computeFollowNode(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.ControlFlowAnalysis cfa)
func ControlFlowAnalysis_ComputeFollowNode_2(node pkg_rhino.Node, cfa *ControlFlowAnalysis) pkg_rhino.Node {
	clinit_ControlFlowAnalysis()
	_ = node
	_ = cfa
	return ControlFlowAnalysis_ComputeFollowNode_3(node, node, cfa)
}

// Java static method: ControlFlowAnalysis.computeFollowNode(com.google.javascript.rhino.Node node)
func ControlFlowAnalysis_ComputeFollowNode_1(node pkg_rhino.Node) pkg_rhino.Node {
	clinit_ControlFlowAnalysis()
	_ = node
	return ControlFlowAnalysis_ComputeFollowNode_3(node, node, nil)
}

// Java static method: ControlFlowAnalysis.computeFollowNode(com.google.javascript.rhino.Node fromNode, com.google.javascript.rhino.Node node, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable ControlFlowAnalysis cfa)
func ControlFlowAnalysis_ComputeFollowNode_3(fromNode pkg_rhino.Node, node pkg_rhino.Node, cfa *ControlFlowAnalysis) pkg_rhino.Node {
	clinit_ControlFlowAnalysis()
	_ = fromNode
	_ = node
	_ = cfa
	var parent pkg_rhino.Node = node.GetParent()
	_ = parent
	if ((parent == 0) || parent.IsFunction()) || ((cfa != nil) && (node == cfa.root)) {
		return pkg_rhino.Node(0)
	}
	switch parent.TokenSwitchOrd_() {
	case int32(74) /* Token.IF */ :
		return ControlFlowAnalysis_ComputeFollowNode_3(fromNode, parent, cfa)
	case int32(76) /* Token.CASE */, int32(77) /* Token.DEFAULT_CASE */ :
		if parent.GetNext() != 0 {
			if parent.GetNext().IsCase() {
				return parent.GetNext().GetSecondChild()
			} else if parent.GetNext().IsDefaultCase() {
				return parent.GetNext().GetFirstChild()
			} else {
				panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Not reachable")))
			}
		} else {
			return ControlFlowAnalysis_ComputeFollowNode_3(fromNode, parent, cfa)
		}
	case int32(80) /* Token.FOR */ :
		return parent.GetSecondChild().GetNext()
	case int32(81) /* Token.FOR_IN */, int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */, int32(78) /* Token.WHILE */, int32(79) /* Token.DO */ :
		return parent
	case int32(48) /* Token.TRY */ :
		if parent.GetFirstChild() == node {
			if NodeUtil_HasFinally(parent) {
				return ControlFlowAnalysis_ComputeFallThrough(parent.GetLastChild())
			} else {
				return ControlFlowAnalysis_ComputeFollowNode_3(fromNode, parent, cfa)
			}
		} else if NodeUtil_GetCatchBlock(parent) == node {
			if NodeUtil_HasFinally(parent) {
				return ControlFlowAnalysis_ComputeFallThrough(node.GetNext())
			} else {
				return ControlFlowAnalysis_ComputeFollowNode_3(fromNode, parent, cfa)
			}
		} else if parent.GetLastChild() == node {
			if (cfa != nil) && (cfa.finallyMap != nil) {
				{
					cur_96 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.MultimapGetNode(cfa.finallyMap, parent))))
					for cur_96.HasNext() {
						finallyNode := pkg_rhino.Node(cur_96.NextNode())
						_ = finallyNode
						cfa.priv_ControlFlowAnalysis_createEdge_2208(fromNode, StaticGet_ControlFlowGraph_Branch_ON_EX(), finallyNode)
					}
				}
			}
			return ControlFlowAnalysis_ComputeFollowNode_3(fromNode, parent, cfa)
		}
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected TRY child "+j2gort.StrVal(node))))
	default:
	}
	var nextSibling pkg_rhino.Node = node.GetNext()
	_ = nextSibling
	for (nextSibling != 0) && nextSibling.IsFunction() {
		nextSibling = nextSibling.GetNext()
	}
	if nextSibling != 0 {
		return ControlFlowAnalysis_ComputeFallThrough(nextSibling)
	} else {
		return ControlFlowAnalysis_ComputeFollowNode_3(fromNode, parent, cfa)
	}
}

// Java static method: ControlFlowAnalysis.computeFallThrough(com.google.javascript.rhino.Node n)
func ControlFlowAnalysis_ComputeFallThrough(n pkg_rhino.Node) pkg_rhino.Node {
	clinit_ControlFlowAnalysis()
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(79) /* Token.DO */, int32(80) /* Token.FOR */ :
		return ControlFlowAnalysis_ComputeFallThrough(n.GetFirstChild())
	case int32(81) /* Token.FOR_IN */, int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */ :
		return n.GetSecondChild()
	case int32(92) /* Token.LABEL */ :
		return ControlFlowAnalysis_ComputeFallThrough(n.GetLastChild())
	default:
		return n
	}
	panic("unreachable switch expression")
}

// Java method: ControlFlowAnalysis.createEdge(com.google.javascript.rhino.Node fromNode, com.google.javascript.jscomp.ControlFlowGraph.Branch branch, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node toNode)
func (this_ *ControlFlowAnalysis) priv_ControlFlowAnalysis_createEdge_2208(fromNode pkg_rhino.Node, branch ControlFlowGraph_Branch, toNode pkg_rhino.Node) {
	_ = fromNode
	_ = branch
	_ = toNode
	this_.cfg.CreateNode(fromNode)
	this_.cfg.CreateNode(toNode)
	this_.cfg.ConnectIfNotFound(fromNode, branch, toNode)
}

// Java method: ControlFlowAnalysis.connectToPossibleExceptionHandler(com.google.javascript.rhino.Node cfgNode, com.google.javascript.rhino.Node target)
func (this_ *ControlFlowAnalysis) priv_ControlFlowAnalysis_connectToPossibleExceptionHandler_2209(cfgNode pkg_rhino.Node, target pkg_rhino.Node) {
	_ = cfgNode
	_ = target
	if ControlFlowAnalysis_MayThrowException(target) {
		this_.priv_ControlFlowAnalysis_connectToPossibleExceptionHandlerUnconditionally_2210(cfgNode)
	}
}

// Java method: ControlFlowAnalysis.connectToPossibleExceptionHandlerUnconditionally(com.google.javascript.rhino.Node cfgNode)
func (this_ *ControlFlowAnalysis) priv_ControlFlowAnalysis_connectToPossibleExceptionHandlerUnconditionally_2210(cfgNode pkg_rhino.Node) {
	_ = cfgNode
	if this_.exceptionHandler.IsEmpty() {
		return
	}
	var lastJump pkg_rhino.Node = cfgNode
	_ = lastJump
	{
		cur_97 := this_.exceptionHandler.Iter()
		for cur_97.HasNext() {
			handler := cur_97.Next()
			_ = handler
			if handler.IsFunction() {
				return
			}
			if !(handler.IsTry()) {
				j2gort.CheckState(false)
			}
			var catchBlock pkg_rhino.Node = NodeUtil_GetCatchBlock(handler)
			_ = catchBlock
			var lastJumpInCatchBlock bool = false
			_ = lastJumpInCatchBlock
			{
				curAnc_ := lastJump
				if curAnc_ == 0 {
					j2gort.ThrowNPE()
				}
				curAnc_ = curAnc_.GetParent()
				for curAnc_ != 0 {
					ancestor := curAnc_
					curAnc_ = curAnc_.GetParent()

					_ = ancestor
					if ancestor == handler {
						break
					} else if ancestor == catchBlock {
						lastJumpInCatchBlock = true
						break
					}
				}
			}
			if !(NodeUtil_HasCatchHandler(catchBlock)) || lastJumpInCatchBlock {
				if lastJump == cfgNode {
					this_.priv_ControlFlowAnalysis_createEdge_2208(cfgNode, StaticGet_ControlFlowGraph_Branch_ON_EX(), handler.GetLastChild())
				} else {
					this_.priv_ControlFlowAnalysis_putFinallyEdge_2180(lastJump, handler.GetLastChild())
				}
			} else {
				if lastJump == cfgNode {
					this_.priv_ControlFlowAnalysis_createEdge_2208(cfgNode, StaticGet_ControlFlowGraph_Branch_ON_EX(), catchBlock)
					return
				} else {
					this_.priv_ControlFlowAnalysis_putFinallyEdge_2180(lastJump, catchBlock)
				}
			}
			lastJump = handler
		}
	}
}

// Java static method: ControlFlowAnalysis.getNextSiblingOfType(com.google.javascript.rhino.Node first, com.google.javascript.rhino.Token[] types)
func ControlFlowAnalysis_GetNextSiblingOfType(first pkg_rhino.Node, types []pkg_rhino.Token) pkg_rhino.Node {
	clinit_ControlFlowAnalysis()
	_ = first
	_ = types
	{
		var c pkg_rhino.Node = first
		_ = c
		for ; c != 0; c = c.GetNext() {
			for _, el_99 := range j2gort.CheckNotNullSlice(types) {
				type_ := el_99
				_ = type_
				if c.GetToken() == type_ {
					return c
				}
			}
		}
	}
	return pkg_rhino.Node(0)
}

// Java static method: ControlFlowAnalysis.isBreakTarget(com.google.javascript.rhino.Node target, java.lang.String label)
func ControlFlowAnalysis_IsBreakTarget(target pkg_rhino.Node, label string) bool {
	clinit_ControlFlowAnalysis()
	_ = target
	_ = label
	return (ControlFlowAnalysis_IsBreakStructure(target, (label != j2gort.NullString)) && ControlFlowAnalysis_MatchLabel(target.GetParent(), label))
}

// Java static method: ControlFlowAnalysis.isContinueTarget(com.google.javascript.rhino.Node target, java.lang.String label)
func ControlFlowAnalysis_IsContinueTarget(target pkg_rhino.Node, label string) bool {
	clinit_ControlFlowAnalysis()
	_ = target
	_ = label
	return (NodeUtil_IsLoopStructure(target) && ControlFlowAnalysis_MatchLabel(target.GetParent(), label))
}

// Java static method: ControlFlowAnalysis.matchLabel(com.google.javascript.rhino.Node target, java.lang.String label)
func ControlFlowAnalysis_MatchLabel(target pkg_rhino.Node, label string) bool {
	clinit_ControlFlowAnalysis()
	_ = target
	_ = label
	if label == j2gort.NullString {
		return true
	}
	for target.IsLabel() {
		if j2gort.CheckNotNullString(target.GetFirstChild().GetString()) == label {
			return true
		}
		target = target.GetParent()
	}
	return false
}

// Java static method: ControlFlowAnalysis.mayThrowException(com.google.javascript.rhino.Node n)
func ControlFlowAnalysis_MayThrowException(n pkg_rhino.Node) bool {
	clinit_ControlFlowAnalysis()
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(28) /* Token.CALL */, int32(31) /* Token.OPTCHAIN_CALL */, int32(130) /* Token.TAGGED_TEMPLATELIT */, int32(26) /* Token.GETPROP */, int32(29) /* Token.OPTCHAIN_GETPROP */, int32(27) /* Token.GETELEM */, int32(30) /* Token.OPTCHAIN_GETELEM */, int32(43) /* Token.THROW */, int32(23) /* Token.NEW */, int32(51) /* Token.ASSIGN */, int32(71) /* Token.INC */, int32(72) /* Token.DEC */, int32(45) /* Token.INSTANCEOF */, int32(44) /* Token.IN */, int32(114) /* Token.YIELD */, int32(115) /* Token.AWAIT */ :
		return true
	case int32(73) /* Token.FUNCTION */ :
		return false
	default:
	}
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			if !(ControlFlowGraph_IsEnteringNewCfgNode(c)) && ControlFlowAnalysis_MayThrowException(c) {
				return true
			}
		}
	}
	return false
}

// Java static method: ControlFlowAnalysis.isBreakStructure(com.google.javascript.rhino.Node n, boolean labeled)
func ControlFlowAnalysis_IsBreakStructure(n pkg_rhino.Node, labeled bool) bool {
	clinit_ControlFlowAnalysis()
	_ = n
	_ = labeled
	switch n.TokenSwitchOrd_() {
	case int32(80) /* Token.FOR */, int32(81) /* Token.FOR_IN */, int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */, int32(79) /* Token.DO */, int32(78) /* Token.WHILE */, int32(75) /* Token.SWITCH */ :
		return true
	case int32(90) /* Token.BLOCK */, int32(74) /* Token.IF */, int32(48) /* Token.TRY */, int32(82) /* Token.BREAK */ :
		return labeled
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: ControlFlowAnalysis.getExceptionHandler(com.google.javascript.rhino.Node n)
func ControlFlowAnalysis_GetExceptionHandler(n pkg_rhino.Node) pkg_rhino.Node {
	clinit_ControlFlowAnalysis()
	_ = n
	{
		var cur pkg_rhino.Node = n
		_ = cur
		for ; !(cur.IsScript()) && !(cur.IsFunction()); cur = cur.GetParent() {
			var catchNode pkg_rhino.Node = ControlFlowAnalysis_GetCatchHandlerForBlock(cur)
			_ = catchNode
			if catchNode != 0 {
				return catchNode
			}
		}
	}
	return pkg_rhino.Node(0)
}

// Java static method: ControlFlowAnalysis.getCatchHandlerForBlock(com.google.javascript.rhino.Node block)
func ControlFlowAnalysis_GetCatchHandlerForBlock(block pkg_rhino.Node) pkg_rhino.Node {
	clinit_ControlFlowAnalysis()
	_ = block
	if (block.IsBlock() && block.GetParent().IsTry()) && (block.GetParent().GetFirstChild() == block) {
		{
			var s pkg_rhino.Node = block.GetNext()
			_ = s
			for ; s != 0; s = s.GetNext() {
				if NodeUtil_HasCatchHandler(s) {
					return s.GetFirstChild()
				}
			}
		}
	}
	return pkg_rhino.Node(0)
}

// Java class: com.google.javascript.jscomp.ControlFlowAnalysis.Builder
type ControlFlowAnalysis_Builder struct {
	compiler                AbstractCompiler `gson:"compiler"`
	cfgRoot                 pkg_rhino.Node   `gson:"cfgRoot"`
	idHash_                 atomic.Int32
	shouldTraverseFunctions bool `gson:"shouldTraverseFunctions"`
	edgeAnnotations         bool `gson:"edgeAnnotations"`
}

func (this_ *ControlFlowAnalysis_Builder) BindSelf_(s any) {
	_ = s
}
func (this_ *ControlFlowAnalysis_Builder) IsJavaNil_() bool { return this_ == nil }
func (*ControlFlowAnalysis_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ControlFlowAnalysis$Builder", "Builder"
}
func (this_ *ControlFlowAnalysis_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ControlFlowAnalysis_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ControlFlowAnalysis.Builder", "com.google.javascript.jscomp.ControlFlowAnalysis$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*ControlFlowAnalysis_Builder](v_)
	})
}

func (this_ *ControlFlowAnalysis_Builder) initFields_ControlFlowAnalysis_Builder() {
	_ = this_

}

// Java constructor: Builder()
func NewControlFlowAnalysis_Builder() *ControlFlowAnalysis_Builder {
	this_ := &ControlFlowAnalysis_Builder{}
	InitControlFlowAnalysis_Builder(this_)
	return this_
}

func InitControlFlowAnalysis_Builder(this_ *ControlFlowAnalysis_Builder) {
	this_.initFields_ControlFlowAnalysis_Builder()
}

// Java method: Builder.setCompiler(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *ControlFlowAnalysis_Builder) SetCompiler(compiler AbstractCompiler) *ControlFlowAnalysis_Builder {
	_ = compiler
	this_.compiler = compiler
	return this_
}

// Java method: Builder.setCfgRoot(com.google.javascript.rhino.Node cfgRoot)
func (this_ *ControlFlowAnalysis_Builder) SetCfgRoot(cfgRoot pkg_rhino.Node) *ControlFlowAnalysis_Builder {
	_ = cfgRoot
	this_.cfgRoot = cfgRoot
	return this_
}

// Java method: Builder.setTraverseFunctions(boolean shouldTraverseFunctions)
func (this_ *ControlFlowAnalysis_Builder) SetTraverseFunctions(shouldTraverseFunctions bool) *ControlFlowAnalysis_Builder {
	_ = shouldTraverseFunctions
	this_.shouldTraverseFunctions = shouldTraverseFunctions
	return this_
}

// Java method: Builder.setIncludeEdgeAnnotations(boolean includeEdgeAnnotations)
func (this_ *ControlFlowAnalysis_Builder) SetIncludeEdgeAnnotations(includeEdgeAnnotations bool) *ControlFlowAnalysis_Builder {
	_ = includeEdgeAnnotations
	this_.edgeAnnotations = includeEdgeAnnotations
	return this_
}

// Java method: Builder.computeCfg()
func (this_ *ControlFlowAnalysis_Builder) ComputeCfg() ControlFlowGraph[pkg_rhino.Node] {
	j2gort.CheckNotNullIface(this_.compiler)
	j2gort.CheckNotNullNode(this_.cfgRoot)
	return ControlFlowAnalysis_ComputeCfgDirect(this_.compiler, this_.cfgRoot, this_.shouldTraverseFunctions, this_.edgeAnnotations)
}

// Java class: com.google.javascript.jscomp.ControlFlowAnalysis.AstControlFlowGraph
type ControlFlowAnalysis_AstControlFlowGraph struct {
	ControlFlowGraphBase[pkg_rhino.Node]
}

func (this_ *ControlFlowAnalysis_AstControlFlowGraph) BindSelf_(s any) {
	_ = s
	this_.ControlFlowGraphBase.BindSelf_(s)
}
func (this_ *ControlFlowAnalysis_AstControlFlowGraph) IsJavaNil_() bool { return this_ == nil }
func (*ControlFlowAnalysis_AstControlFlowGraph) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ControlFlowAnalysis$AstControlFlowGraph", "AstControlFlowGraph"
}
func (this_ *ControlFlowAnalysis_AstControlFlowGraph) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ControlFlowGraphBase.IdentityHash_()
}
func (*ControlFlowAnalysis_AstControlFlowGraph) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ControlFlowAnalysis.AstControlFlowGraph", "com.google.javascript.jscomp.ControlFlowAnalysis$AstControlFlowGraph", func(v_ any) bool {
		return j2gort.InstanceOf[*ControlFlowAnalysis_AstControlFlowGraph](v_)
	})
}
func (*ControlFlowAnalysis_AstControlFlowGraph) Implements_jscomp_ControlFlowGraph()   {}
func (*ControlFlowAnalysis_AstControlFlowGraph) Implements_graph_LinkedDirectedGraph() {}
func (*ControlFlowAnalysis_AstControlFlowGraph) Implements_graph_DiGraph()             {}
func (*ControlFlowAnalysis_AstControlFlowGraph) Implements_graph_GraphvizGraph()       {}
func (*ControlFlowAnalysis_AstControlFlowGraph) Implements_graph_Graph()               {}
func (*ControlFlowAnalysis_AstControlFlowGraph) Implements_graph_AdjacencyGraph()      {}

var static_ControlFlowAnalysis_AstControlFlowGraph_FORWARD_COMPARATOR j2gort.Ext_java_util_Comparator

func StaticGet_ControlFlowAnalysis_AstControlFlowGraph_FORWARD_COMPARATOR() j2gort.Ext_java_util_Comparator {
	clinit_ControlFlowAnalysis_AstControlFlowGraph()
	return static_ControlFlowAnalysis_AstControlFlowGraph_FORWARD_COMPARATOR
}
func StaticSet_ControlFlowAnalysis_AstControlFlowGraph_FORWARD_COMPARATOR(v j2gort.Ext_java_util_Comparator) {
	clinit_ControlFlowAnalysis_AstControlFlowGraph()
	static_ControlFlowAnalysis_AstControlFlowGraph_FORWARD_COMPARATOR = v
}

var static_ControlFlowAnalysis_AstControlFlowGraph_BACKWARD_COMPARATOR j2gort.Ext_java_util_Comparator

func StaticGet_ControlFlowAnalysis_AstControlFlowGraph_BACKWARD_COMPARATOR() j2gort.Ext_java_util_Comparator {
	clinit_ControlFlowAnalysis_AstControlFlowGraph()
	return static_ControlFlowAnalysis_AstControlFlowGraph_BACKWARD_COMPARATOR
}
func StaticSet_ControlFlowAnalysis_AstControlFlowGraph_BACKWARD_COMPARATOR(v j2gort.Ext_java_util_Comparator) {
	clinit_ControlFlowAnalysis_AstControlFlowGraph()
	static_ControlFlowAnalysis_AstControlFlowGraph_BACKWARD_COMPARATOR = v
}

var clinitDone_ControlFlowAnalysis_AstControlFlowGraph uint32
var clinitState_ControlFlowAnalysis_AstControlFlowGraph j2gort.ClinitState

func clinit_ControlFlowAnalysis_AstControlFlowGraph() {
	if atomic.LoadUint32(&clinitDone_ControlFlowAnalysis_AstControlFlowGraph) != 0 {
		return
	}
	clinitSlow_ControlFlowAnalysis_AstControlFlowGraph()
}

func clinitSlow_ControlFlowAnalysis_AstControlFlowGraph() {
	if !j2gort.EnterClinit(&clinitState_ControlFlowAnalysis_AstControlFlowGraph, &clinitDone_ControlFlowAnalysis_AstControlFlowGraph) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ControlFlowAnalysis_AstControlFlowGraph, &clinitDone_ControlFlowAnalysis_AstControlFlowGraph)
	static_ControlFlowAnalysis_AstControlFlowGraph_FORWARD_COMPARATOR = lambdaSingleton_control_flow_analysis_1
	static_ControlFlowAnalysis_AstControlFlowGraph_BACKWARD_COMPARATOR = lambdaSingleton_control_flow_analysis_2
}

func (this_ *ControlFlowAnalysis_AstControlFlowGraph) initFields_ControlFlowAnalysis_AstControlFlowGraph() {
	_ = this_
}

// Java constructor: AstControlFlowGraph(com.google.javascript.rhino.Node entry, boolean edgeAnnotations)
func NewControlFlowAnalysis_AstControlFlowGraph(entry pkg_rhino.Node, edgeAnnotations bool) *ControlFlowAnalysis_AstControlFlowGraph {
	this_ := &ControlFlowAnalysis_AstControlFlowGraph{}
	this_.BindSelf_(this_)
	InitControlFlowAnalysis_AstControlFlowGraph(this_, entry, edgeAnnotations)
	return this_
}

func InitControlFlowAnalysis_AstControlFlowGraph(this_ *ControlFlowAnalysis_AstControlFlowGraph, entry pkg_rhino.Node, edgeAnnotations bool) {
	_ = entry
	_ = edgeAnnotations
	InitControlFlowGraph[pkg_rhino.Node](&this_.ControlFlowGraphBase, entry, true, edgeAnnotations)
	this_.initFields_ControlFlowAnalysis_AstControlFlowGraph()
}

// Java method: AstControlFlowGraph.hasPriorityComparator()
func (this_ *ControlFlowAnalysis_AstControlFlowGraph) HasPriorityComparator() bool {
	return true
}

// Java method: AstControlFlowGraph.getOptionalNodeComparator(boolean isForward)
func (this_ *ControlFlowAnalysis_AstControlFlowGraph) GetOptionalNodeComparator(isForward bool) j2gort.Ext_java_util_Comparator {
	_ = isForward
	if isForward {
		return StaticGet_ControlFlowAnalysis_AstControlFlowGraph_FORWARD_COMPARATOR()
	} else {
		return StaticGet_ControlFlowAnalysis_AstControlFlowGraph_BACKWARD_COMPARATOR()
	}
}
func (*ControlFlowAnalysis) CanReleaseSyntacticScopes_() {}
