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
// Port of com/google/javascript/jscomp/MaybeReachingVariableUse.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_graph "github.com/paullepoulpe/jsgo/transpiled/gen/graph"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.MaybeReachingVariableUse
type MaybeReachingVariableUse struct {
	DataFlowAnalysisBase[pkg_rhino.Node, *MaybeReachingVariableUse_ReachingUses]
	varNames                []string                               `gson:"varNames"`
	isEscaped               []bool                                 `gson:"isEscaped"`
	nameTable               []string                               `gson:"nameTable"`
	indexTable              []int32                                `gson:"indexTable"`
	emptyLattice            *MaybeReachingVariableUse_ReachingUses `gson:"emptyLattice"`
	flowOutput              *MaybeReachingVariableUse_ReachingUses `gson:"flowOutput"`
	numVars                 int32                                  `gson:"numVars"`
	slab_ReachingUses       j2gort.PooledSlab[MaybeReachingVariableUse_ReachingUses]
	slab_ReachingUsesJoinOp j2gort.PooledSlab[MaybeReachingVariableUse_ReachingUsesJoinOp]
}

func (this_ *MaybeReachingVariableUse) BindSelf_(s any) {
	_ = s
	this_.DataFlowAnalysisBase.BindSelf_(s)
}
func (this_ *MaybeReachingVariableUse) IsJavaNil_() bool { return this_ == nil }
func (*MaybeReachingVariableUse) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.MaybeReachingVariableUse", "MaybeReachingVariableUse"
}
func (this_ *MaybeReachingVariableUse) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.DataFlowAnalysisBase.IdentityHash_()
}
func (*MaybeReachingVariableUse) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.MaybeReachingVariableUse", "com.google.javascript.jscomp.MaybeReachingVariableUse", func(v_ any) bool {
		return j2gort.InstanceOf[*MaybeReachingVariableUse](v_)
	})
}
func (*MaybeReachingVariableUse) Implements_jscomp_DataFlowAnalysis() {}

func (this_ *MaybeReachingVariableUse) initFields_MaybeReachingVariableUse() {
	_ = this_
}

// Java constructor: MaybeReachingVariableUse(com.google.javascript.jscomp.ControlFlowGraph<com.google.javascript.rhino.Node> cfg, java.util.Set<com.google.javascript.jscomp.Var> escaped, java.util.Map<java.lang.String,com.google.javascript.jscomp.Var> allVarsInFn)
func NewMaybeReachingVariableUse(cfg ControlFlowGraph[pkg_rhino.Node], escaped j2gort.Ext_java_util_Set, allVarsInFn j2gort.Ext_java_util_Map) *MaybeReachingVariableUse {
	this_ := &MaybeReachingVariableUse{}
	this_.BindSelf_(this_)
	InitMaybeReachingVariableUse(this_, cfg, escaped, allVarsInFn)
	return this_
}

func InitMaybeReachingVariableUse(this_ *MaybeReachingVariableUse, cfg ControlFlowGraph[pkg_rhino.Node], escaped j2gort.Ext_java_util_Set, allVarsInFn j2gort.Ext_java_util_Map) {
	_ = cfg
	_ = escaped
	_ = allVarsInFn
	InitDataFlowAnalysis[pkg_rhino.Node, *MaybeReachingVariableUse_ReachingUses](&this_.DataFlowAnalysisBase, cfg)
	this_.initFields_MaybeReachingVariableUse()
	var n int32 = allVarsInFn.Size()
	_ = n
	this_.numVars = n
	this_.varNames = j2gort.NewStringArray(n)
	this_.isEscaped = make([]bool, n)
	var i int32 = int32(0)
	_ = i
	{
		cur_1 := j2gort.NewMapEntryCursor(allVarsInFn)
		for cur_1.NextEntry() {
			{
				a_2 := this_.varNames
				i_3 := i
				a_2[i_3] = cur_1.KeyStr()
			}
			{
				a_4 := this_.isEscaped
				i_5 := i
				a_4[i_5] = escaped.Contains(j2gort.PtrToAny(j2gort.CastPtr[Var](cur_1.GetValue())))
			}
			i = i + 1
		}
	}
	if n > int32(8) {
		var cap_ int32 = javalang.Shl32(j2gort.IntHighestOneBit(((n * int32(2)) - int32(1))), int(int32(1)))
		_ = cap_
		var mask int32 = (cap_ - int32(1))
		_ = mask
		var keys []string = j2gort.NewStringArray(cap_)
		_ = keys
		var indices []int32 = make([]int32, cap_)
		_ = indices
		{
			var j int32 = int32(0)
			_ = j
			for ; j < n; j = j + 1 {
				var key string = this_.varNames[j]
				_ = key
				var slot int32 = (javalang.StringHashCode(j2gort.CheckNotNullString(key)) & mask)
				_ = slot
				for keys[slot] != j2gort.NullString {
					slot = ((slot + int32(1)) & mask)
				}
				{
					a_6 := keys
					i_7 := slot
					a_6[i_7] = key
				}
				{
					a_8 := indices
					i_9 := slot
					a_8[i_9] = j
				}
			}
		}
		this_.nameTable = keys
		this_.indexTable = indices
	} else {
		this_.nameTable = nil
		this_.indexTable = nil
	}
	this_.emptyLattice = allocPooled_MaybeReachingVariableUse_ReachingUses_1_int32(&this_.

		// Java method: MaybeReachingVariableUse.lookupVarIndex(java.lang.String name)
		slab_ReachingUses, n)
}

func (this_ *MaybeReachingVariableUse) priv_MaybeReachingVariableUse_lookupVarIndex_3856(name string) int32 {
	_ = name
	var table []string = this_.nameTable
	_ = table
	if table != nil {
		var mask int32 = (int32(len(table)) - int32(1))
		_ = mask
		var slot int32 = (javalang.StringHashCode(j2gort.CheckNotNullString(name)) & mask)
		_ = slot
		for {
			var k string = table[slot]
			_ = k
			if k == j2gort.NullString {
				return int32(-1)
			}
			if j2gort.CheckNotNullString(k) == name {
				return this_.indexTable[slot]
			}
			slot = ((slot + int32(1)) & mask)
		}
	}
	var names []string = this_.varNames
	_ = names
	{
		var i int32 = int32(0)
		_ = i
		var len_ int32 = int32(len(names))
		_ = len_
		for ; i < len_; i = i + 1 {
			if j2gort.CheckNotNullString(names[i]) == name {
				return i
			}
		}
	}
	return int32(-1)
}

// Java method: MaybeReachingVariableUse.isForward()
func (this_ *MaybeReachingVariableUse) IsForward() bool {
	return false
}

// Java method: MaybeReachingVariableUse.createEntryLattice()
func (this_ *MaybeReachingVariableUse) CreateEntryLattice() *MaybeReachingVariableUse_ReachingUses {
	return this_.emptyLattice
}

// Java method: MaybeReachingVariableUse.createInitialEstimateLattice()
func (this_ *MaybeReachingVariableUse) CreateInitialEstimateLattice() *MaybeReachingVariableUse_ReachingUses {
	return this_.emptyLattice
}

// Java method: MaybeReachingVariableUse.createFlowJoiner()
func (this_ *MaybeReachingVariableUse) CreateFlowJoiner() DataFlowAnalysis_FlowJoiner[*MaybeReachingVariableUse_ReachingUses] {
	return ToIface_DataFlowAnalysis_FlowJoiner[*MaybeReachingVariableUse_ReachingUses](allocPooled_MaybeReachingVariableUse_ReachingUsesJoinOp(&this_.slab_ReachingUsesJoinOp,

		// Java method: MaybeReachingVariableUse.flowThrough(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.MaybeReachingVariableUse.ReachingUses input)
		&this_.slab_ReachingUses, this_.numVars))
}

func (this_ *MaybeReachingVariableUse) FlowThrough(n pkg_rhino.Node, input *MaybeReachingVariableUse_ReachingUses) *MaybeReachingVariableUse_ReachingUses {
	_ = n
	_ = input
	this_.flowOutput = nil
	var conditional bool = this_.priv_MaybeReachingVariableUse_hasExceptionHandler_3858(n)
	_ = conditional
	this_.priv_MaybeReachingVariableUse_computeMayUse_3859(n, n, input, conditional)
	var out *MaybeReachingVariableUse_ReachingUses = this_.flowOutput
	_ = out
	this_.flowOutput = nil
	if out != nil {
		return out
	} else {
		return input
	}
}

// Java method: MaybeReachingVariableUse.writableFlowOutput(com.google.javascript.jscomp.MaybeReachingVariableUse.ReachingUses input)
func (this_ *MaybeReachingVariableUse) priv_MaybeReachingVariableUse_writableFlowOutput_3857(input *MaybeReachingVariableUse_ReachingUses) *MaybeReachingVariableUse_ReachingUses {
	_ = input
	var out *MaybeReachingVariableUse_ReachingUses = this_.flowOutput
	_ = out
	if out == nil {
		out = allocPooled_MaybeReachingVariableUse_ReachingUses_1_MaybeReachingVariableUse_ReachingUses(&this_.slab_ReachingUses, input)
		this_.flowOutput = out
	}
	return out
}

// Java method: MaybeReachingVariableUse.hasExceptionHandler(com.google.javascript.rhino.Node cfgNode)
func (this_ *MaybeReachingVariableUse) priv_MaybeReachingVariableUse_hasExceptionHandler_3858(cfgNode pkg_rhino.Node) bool {
	_ = cfgNode
	var branchEdges j2gort.Ext_java_util_List = this_.GetCfg().GetOutEdges(cfgNode)
	_ = branchEdges
	{
		var i int32 = int32(0)
		_ = i
		var size int32 = branchEdges.Size()
		_ = size
		for ; i < size; i = i + 1 {
			if j2gort.AsIface[pkg_graph.DiGraph_DiGraphEdge[pkg_rhino.Node, ControlFlowGraph_Branch]](branchEdges.Get(i)).GetValue() == StaticGet_ControlFlowGraph_Branch_ON_EX() {
				return true
			}
		}
	}
	return false
}

// Java method: MaybeReachingVariableUse.computeMayUse(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node cfgNode, com.google.javascript.jscomp.MaybeReachingVariableUse.ReachingUses output, boolean conditional)
func (this_ *MaybeReachingVariableUse) priv_MaybeReachingVariableUse_computeMayUse_3859(n pkg_rhino.Node, cfgNode pkg_rhino.Node, output *MaybeReachingVariableUse_ReachingUses, conditional bool) {
	_ = n
	_ = cfgNode
	_ = output
	_ = conditional
	switch n.TokenSwitchOrd_() {
	case int32(90) /* Token.BLOCK */, int32(89) /* Token.ROOT */, int32(73) /* Token.FUNCTION */ :
		return
	case int32(32) /* Token.NAME */ :
		if NodeUtil_IsLhsByDestructuring(n) {
			if !(conditional) {
				this_.priv_MaybeReachingVariableUse_removeFromUseIfLocal_3861(n.GetString(), output)
			}
		} else {
			this_.priv_MaybeReachingVariableUse_addToUseIfLocal_3860(n.GetString(), cfgNode, output)
		}
		return
	case int32(78) /* Token.WHILE */, int32(79) /* Token.DO */, int32(74) /* Token.IF */, int32(80) /* Token.FOR */ :
		var condExpr pkg_rhino.Node = NodeUtil_GetConditionExpression(n)
		_ = condExpr
		this_.priv_MaybeReachingVariableUse_computeMayUse_3859(condExpr, cfgNode, output, conditional)
		return
	case int32(81) /* Token.FOR_IN */, int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */ :
		var lhs pkg_rhino.Node = n.GetFirstChild()
		_ = lhs
		var rhs pkg_rhino.Node = lhs.GetNext()
		_ = rhs
		if NodeUtil_IsNameDeclaration(lhs) {
			lhs = lhs.GetLastChild()
			if lhs.IsDestructuringLhs() {
				lhs = lhs.GetFirstChild()
			}
		}
		if lhs.IsDestructuringPattern() {
			this_.priv_MaybeReachingVariableUse_computeMayUse_3859(lhs, cfgNode, output, true)
		}
		this_.priv_MaybeReachingVariableUse_computeMayUse_3859(rhs, cfgNode, output, conditional)
		return
	case int32(69) /* Token.AND */, int32(68) /* Token.OR */, int32(70) /* Token.COALESCE */, int32(29) /* Token.OPTCHAIN_GETPROP */, int32(30) /* Token.OPTCHAIN_GETELEM */ :
		this_.priv_MaybeReachingVariableUse_computeMayUse_3859(n.GetLastChild(), cfgNode, output, true)
		this_.priv_MaybeReachingVariableUse_computeMayUse_3859(n.GetFirstChild(), cfgNode, output, conditional)
		return
	case int32(31) /* Token.OPTCHAIN_CALL */ :
		{
			var c pkg_rhino.Node = n.GetLastChild()
			_ = c
			for ; c != n.GetFirstChild(); c = c.GetPrevious() {
				this_.priv_MaybeReachingVariableUse_computeMayUse_3859(c, cfgNode, output, true)
			}
		}
		this_.priv_MaybeReachingVariableUse_computeMayUse_3859(n.GetFirstChild(), cfgNode, output, conditional)
		return
	case int32(67) /* Token.HOOK */ :
		this_.priv_MaybeReachingVariableUse_computeMayUse_3859(n.GetLastChild(), cfgNode, output, true)
		this_.priv_MaybeReachingVariableUse_computeMayUse_3859(n.GetSecondChild(), cfgNode, output, true)
		this_.priv_MaybeReachingVariableUse_computeMayUse_3859(n.GetFirstChild(), cfgNode, output, conditional)
		return
	case int32(84) /* Token.VAR */, int32(111) /* Token.LET */, int32(97) /* Token.CONST */ :
		var varName pkg_rhino.Node = n.GetFirstChild()
		_ = varName
		if !(n.HasChildren()) {
			j2gort.CheckStateFmt(false, "AST should be normalized (%s)", j2gort.BoxNode(n))
		}
		if varName.IsDestructuringLhs() {
			this_.priv_MaybeReachingVariableUse_computeMayUse_3859(varName.GetFirstChild(), cfgNode, output, conditional)
			this_.priv_MaybeReachingVariableUse_computeMayUse_3859(varName.GetSecondChild(), cfgNode, output, conditional)
		} else if varName.HasChildren() {
			this_.priv_MaybeReachingVariableUse_computeMayUse_3859(varName.GetFirstChild(), cfgNode, output, conditional)
			if !(conditional) {
				this_.priv_MaybeReachingVariableUse_removeFromUseIfLocal_3861(varName.GetString(), output)
			}
		}
		return
	case int32(134) /* Token.DEFAULT_VALUE */ :
		if n.GetFirstChild().IsDestructuringPattern() {
			this_.priv_MaybeReachingVariableUse_computeMayUse_3859(n.GetFirstChild(), cfgNode, output, conditional)
			this_.priv_MaybeReachingVariableUse_computeMayUse_3859(n.GetSecondChild(), cfgNode, output, true)
		} else if n.GetFirstChild().IsName() {
			if !(conditional) {
				this_.priv_MaybeReachingVariableUse_removeFromUseIfLocal_3861(n.GetFirstChild().GetString(), output)
			}
			this_.priv_MaybeReachingVariableUse_computeMayUse_3859(n.GetSecondChild(), cfgNode, output, true)
		} else {
			this_.priv_MaybeReachingVariableUse_computeMayUse_3859(n.GetSecondChild(), cfgNode, output, true)
			this_.priv_MaybeReachingVariableUse_computeMayUse_3859(n.GetFirstChild(), cfgNode, output, conditional)
		}
	default:
		if NodeUtil_IsAssignmentOp(n) && n.GetFirstChild().IsName() {
			if !(!(NodeUtil_IsLogicalAssignmentOp(n))) {
				j2gort.CheckState(false)
			}
			var name pkg_rhino.Node = n.GetFirstChild()
			_ = name
			if !(conditional) {
				this_.priv_MaybeReachingVariableUse_removeFromUseIfLocal_3861(name.GetString(), output)
			}
			if !(n.IsAssign()) {
				this_.priv_MaybeReachingVariableUse_addToUseIfLocal_3860(name.GetString(), cfgNode, output)
			}
			this_.priv_MaybeReachingVariableUse_computeMayUse_3859(name.GetNext(), cfgNode, output, conditional)
		} else if n.IsAssign() && n.GetFirstChild().IsDestructuringPattern() {
			this_.priv_MaybeReachingVariableUse_computeMayUse_3859(n.GetFirstChild(), cfgNode, output, conditional)
			this_.priv_MaybeReachingVariableUse_computeMayUse_3859(n.GetSecondChild(), cfgNode, output, conditional)
		} else {
			{
				var c_2 pkg_rhino.Node = n.GetLastChild()
				_ = c_2
				for ; c_2 != 0; c_2 = c_2.GetPrevious() {
					this_.priv_MaybeReachingVariableUse_computeMayUse_3859(c_2, cfgNode, output, conditional)
				}
			}
		}
	}
}

// Java method: MaybeReachingVariableUse.addToUseIfLocal(java.lang.String name, com.google.javascript.rhino.Node node, com.google.javascript.jscomp.MaybeReachingVariableUse.ReachingUses input)
func (this_ *MaybeReachingVariableUse) priv_MaybeReachingVariableUse_addToUseIfLocal_3860(name string, node pkg_rhino.Node, input *MaybeReachingVariableUse_ReachingUses) {
	_ = name
	_ = node
	_ = input
	var varIdx int32 = this_.priv_MaybeReachingVariableUse_lookupVarIndex_3856(name)
	_ = varIdx
	if (varIdx < int32(0)) || this_.isEscaped[varIdx] {
		return
	}
	var cur *MaybeReachingVariableUse_ReachingUses
	_ = cur
	if this_.flowOutput != nil {
		cur = this_.flowOutput
	} else {
		cur = input
	}
	if cur.ContainsUse(varIdx, node) {
		return
	}
	this_.priv_MaybeReachingVariableUse_writableFlowOutput_3857(input).Put(varIdx, node)
}

// Java method: MaybeReachingVariableUse.removeFromUseIfLocal(java.lang.String name, com.google.javascript.jscomp.MaybeReachingVariableUse.ReachingUses input)
func (this_ *MaybeReachingVariableUse) priv_MaybeReachingVariableUse_removeFromUseIfLocal_3861(name string, input *MaybeReachingVariableUse_ReachingUses) {
	_ = name
	_ = input
	var varIdx int32 = this_.priv_MaybeReachingVariableUse_lookupVarIndex_3856(name)
	_ = varIdx
	if (varIdx < int32(0)) || this_.isEscaped[varIdx] {
		return
	}
	var cur *MaybeReachingVariableUse_ReachingUses
	_ = cur
	if this_.flowOutput != nil {
		cur = this_.flowOutput
	} else {
		cur = input
	}
	if cur.IsEmptyAt(varIdx) {
		return
	}
	this_.priv_MaybeReachingVariableUse_writableFlowOutput_3857(input).RemoveAll(varIdx)
}

// Java method: MaybeReachingVariableUse.getUses(java.lang.String name, com.google.javascript.rhino.Node defNode)
func (this_ *MaybeReachingVariableUse) GetUses(name string, defNode pkg_rhino.Node) j2gort.Ext_java_lang_Iterable {
	_ = name
	_ = defNode
	var n *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch] = j2gort.CastPtr[pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch]](this_.GetCfg().GetNode(defNode))
	_ = n
	j2gort.CheckNotNullPtr(n)
	var varIdx int32 = this_.priv_MaybeReachingVariableUse_lookupVarIndex_3856(name)
	_ = varIdx
	if varIdx < int32(0) {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
	}
	return this_.GetNodeOutLattice(n).GetUsesAt(varIdx)
}

// Java method: MaybeReachingVariableUse.hasSingleUse(java.lang.String name, com.google.javascript.rhino.Node defNode)
func (this_ *MaybeReachingVariableUse) HasSingleUse(name string, defNode pkg_rhino.Node) bool {
	_ = name
	_ = defNode
	var n *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch] = j2gort.CastPtr[pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch]](this_.GetCfg().GetNode(defNode))
	_ = n
	j2gort.CheckNotNullPtr(n)
	var varIdx int32 = this_.priv_MaybeReachingVariableUse_lookupVarIndex_3856(name)
	_ = varIdx
	if varIdx < int32(0) {
		return false
	}
	return this_.GetNodeOutLattice(n).HasSingleUseAt(varIdx)
}

// Java method: MaybeReachingVariableUse.needsNodeAnnotations()
func (this_ *MaybeReachingVariableUse) NeedsNodeAnnotations() bool {
	return false
}

// Java method: MaybeReachingVariableUse.isLatticeEqual(com.google.javascript.jscomp.MaybeReachingVariableUse.ReachingUses a, com.google.javascript.jscomp.MaybeReachingVariableUse.ReachingUses b)
func (this_ *MaybeReachingVariableUse) IsLatticeEqual(a *MaybeReachingVariableUse_ReachingUses, b *MaybeReachingVariableUse_ReachingUses) bool {
	_ = a
	_ = b
	return a.EqualsUses(b)
}

// Java class: com.google.javascript.jscomp.MaybeReachingVariableUse.ReachingUses
type MaybeReachingVariableUse_ReachingUses struct {
	singleUse []pkg_rhino.Node   `gson:"singleUse"`
	multiUse  [][]pkg_rhino.Node `gson:"multiUse"`
	numVars   int32              `gson:"numVars"`
}

func (this_ *MaybeReachingVariableUse_ReachingUses) BindSelf_(s any) {
	_ = s
}
func (this_ *MaybeReachingVariableUse_ReachingUses) IsJavaNil_() bool { return this_ == nil }
func (*MaybeReachingVariableUse_ReachingUses) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.MaybeReachingVariableUse$ReachingUses", "ReachingUses"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.MaybeReachingVariableUse.ReachingUses", "com.google.javascript.jscomp.MaybeReachingVariableUse$ReachingUses", func(v_ any) bool {
		return j2gort.InstanceOf[*MaybeReachingVariableUse_ReachingUses](v_)
	})
}
func (*MaybeReachingVariableUse_ReachingUses) Implements_graph_LatticeElement() {}
func (*MaybeReachingVariableUse_ReachingUses) Implements_graph_Annotation()     {}

func (this_ *MaybeReachingVariableUse_ReachingUses) initFields_MaybeReachingVariableUse_ReachingUses() {
	_ = this_
}

// Java constructor: ReachingUses(int numVars)
func NewMaybeReachingVariableUse_ReachingUses_1_int32(numVars int32) *MaybeReachingVariableUse_ReachingUses {
	this_ := &MaybeReachingVariableUse_ReachingUses{}
	InitMaybeReachingVariableUse_ReachingUses_1_int32(this_, numVars)
	return this_
}
func allocPooled_MaybeReachingVariableUse_ReachingUses_1_int32(slab_ *j2gort.PooledSlab[MaybeReachingVariableUse_ReachingUses], numVars int32) *MaybeReachingVariableUse_ReachingUses {
	this_ := slab_.Alloc()
	InitMaybeReachingVariableUse_ReachingUses_1_int32(this_, numVars)
	return this_
}

func InitMaybeReachingVariableUse_ReachingUses_1_int32(this_ *MaybeReachingVariableUse_ReachingUses, numVars int32) {
	_ = numVars
	this_.initFields_MaybeReachingVariableUse_ReachingUses()
	this_.numVars = numVars
	this_.singleUse = nil
	this_.multiUse = nil
}

// Java constructor: ReachingUses(com.google.javascript.jscomp.MaybeReachingVariableUse.ReachingUses other)
func NewMaybeReachingVariableUse_ReachingUses_1_MaybeReachingVariableUse_ReachingUses(other *MaybeReachingVariableUse_ReachingUses) *MaybeReachingVariableUse_ReachingUses {
	this_ := &MaybeReachingVariableUse_ReachingUses{}
	InitMaybeReachingVariableUse_ReachingUses_1_MaybeReachingVariableUse_ReachingUses(this_, other)
	return this_
}
func allocPooled_MaybeReachingVariableUse_ReachingUses_1_MaybeReachingVariableUse_ReachingUses(slab_ *j2gort.PooledSlab[MaybeReachingVariableUse_ReachingUses], other *MaybeReachingVariableUse_ReachingUses) *MaybeReachingVariableUse_ReachingUses {
	this_ := slab_.Alloc()
	InitMaybeReachingVariableUse_ReachingUses_1_MaybeReachingVariableUse_ReachingUses(this_, other)
	return this_
}

func InitMaybeReachingVariableUse_ReachingUses_1_MaybeReachingVariableUse_ReachingUses(this_ *MaybeReachingVariableUse_ReachingUses, other *MaybeReachingVariableUse_ReachingUses) {
	_ = other
	this_.initFields_MaybeReachingVariableUse_ReachingUses()
	this_.numVars = other.numVars
	if other.singleUse != nil {
		this_.singleUse = j2gort.CloneSlice(other.singleUse)
	} else {
		this_.singleUse = nil
	}
	if other.multiUse != nil {
		this_.multiUse = j2gort.CloneSlice(other.multiUse)
	} else {
		this_.multiUse = nil
	}
}

// Java method: ReachingUses.isEmptyAt(int varIdx)
func (this_ *MaybeReachingVariableUse_ReachingUses) IsEmptyAt(varIdx int32) bool {
	_ = varIdx
	return (((this_.singleUse == nil) || (this_.singleUse[varIdx] == 0)) && ((this_.multiUse == nil) || (this_.multiUse[varIdx] == nil)))
}

// Java method: ReachingUses.hasSingleUseAt(int varIdx)
func (this_ *MaybeReachingVariableUse_ReachingUses) HasSingleUseAt(varIdx int32) bool {
	_ = varIdx
	return ((this_.singleUse != nil) && (this_.singleUse[varIdx] != 0))
}

// Java method: ReachingUses.containsUse(int varIdx, com.google.javascript.rhino.Node n)
func (this_ *MaybeReachingVariableUse_ReachingUses) ContainsUse(varIdx int32, n pkg_rhino.Node) bool {
	_ = varIdx
	_ = n
	if (this_.singleUse != nil) && (this_.singleUse[varIdx] == n) {
		return true
	}
	if this_.multiUse != nil {
		var ms []pkg_rhino.Node = this_.multiUse[varIdx]
		_ = ms
		if ms != nil {
			return MaybeReachingVariableUse_ReachingUses_ContainsSorted(ms, n)
		}
	}
	return false
}

// Java method: ReachingUses.getUsesAt(int varIdx)
func (this_ *MaybeReachingVariableUse_ReachingUses) GetUsesAt(varIdx int32) j2gort.Ext_java_lang_Iterable {
	_ = varIdx
	if this_.singleUse != nil {
		var single pkg_rhino.Node = this_.singleUse[varIdx]
		_ = single
		if single != 0 {
			return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_ImmutableList_Of_1(j2gort.BoxNode(single)))
		}
	}
	if this_.multiUse != nil {
		var ms []pkg_rhino.Node = this_.multiUse[varIdx]
		_ = ms
		if ms != nil {
			var iifeArg_1 []any
		iifeSwitch_1:
			switch {
			default:

				arr_10 := ms
				if arr_10 == nil {
					iifeArg_1 = nil
					break iifeSwitch_1
				}
				out_11 := make([]any, len(arr_10))
				for i_, el_12 := range arr_10 {
					out_11[i_] = j2gort.BoxNode(el_12)
				}
				iifeArg_1 = out_11
				break iifeSwitch_1
			}
			return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_java_util_Arrays_AsList(iifeArg_1))
		}
	}
	return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
}

// Java method: ReachingUses.removeAll(int varIdx)
func (this_ *MaybeReachingVariableUse_ReachingUses) RemoveAll(varIdx int32) {
	_ = varIdx
	if this_.singleUse != nil {
		{
			a_13 := this_.singleUse
			i_14 := varIdx
			a_13[i_14] = pkg_rhino.Node(0)
		}
	}
	if this_.multiUse != nil {
		{
			a_15 := this_.multiUse
			i_16 := varIdx
			a_15[i_16] = nil
		}
	}
}

// Java method: ReachingUses.put(int varIdx, com.google.javascript.rhino.Node n)
func (this_ *MaybeReachingVariableUse_ReachingUses) Put(varIdx int32, n pkg_rhino.Node) {
	_ = varIdx
	_ = n
	var mArr [][]pkg_rhino.Node = this_.multiUse
	_ = mArr
	if mArr != nil {
		var existingMulti []pkg_rhino.Node = mArr[varIdx]
		_ = existingMulti
		if existingMulti != nil {
			if !(MaybeReachingVariableUse_ReachingUses_ContainsSorted(existingMulti, n)) {
				{
					a_17 := mArr
					i_18 := varIdx
					a_17[i_18] = MaybeReachingVariableUse_ReachingUses_InsertSorted(existingMulti, n)
				}
			}
			return
		}
	}
	var sArr []pkg_rhino.Node = this_.singleUse
	_ = sArr
	if sArr == nil {
		sArr = make([]pkg_rhino.Node, this_.numVars)
		{
			a_19 := sArr
			i_20 := varIdx
			a_19[i_20] = n
		}
		this_.singleUse = sArr
		return
	}
	var existingSingle pkg_rhino.Node = sArr[varIdx]
	_ = existingSingle
	if existingSingle == 0 {
		{
			a_21 := sArr
			i_22 := varIdx
			a_21[i_22] = n
		}
	} else if existingSingle != n {
		{
			a_23 := sArr
			i_24 := varIdx
			a_23[i_24] = pkg_rhino.Node(0)
		}
		if mArr == nil {
			mArr = make([][]pkg_rhino.Node, this_.numVars)
			this_.multiUse = mArr
		}
		{
			a_25 := mArr
			i_26 := varIdx
			a_25[i_26] = MaybeReachingVariableUse_ReachingUses_PairSorted(existingSingle, n)
		}
	}
}

// Java method: ReachingUses.join(com.google.javascript.jscomp.MaybeReachingVariableUse.ReachingUses other)
func (this_ *MaybeReachingVariableUse_ReachingUses) Join(other *MaybeReachingVariableUse_ReachingUses) {
	_ = other
	var otherSingle []pkg_rhino.Node = other.singleUse
	_ = otherSingle
	var otherMulti [][]pkg_rhino.Node = other.multiUse
	_ = otherMulti
	if (otherSingle == nil) && (otherMulti == nil) {
		return
	}
	if (this_.singleUse == nil) && (this_.multiUse == nil) {
		if otherSingle != nil {
			this_.singleUse = j2gort.CloneSlice(otherSingle)
		} else {
			this_.singleUse = nil
		}
		if otherMulti != nil {
			this_.multiUse = j2gort.CloneSlice(otherMulti)
		} else {
			this_.multiUse = nil
		}
		return
	}
	{
		var i int32 = int32(0)
		_ = i
		var len_ int32 = this_.numVars
		_ = len_
		for ; i < len_; i = i + 1 {
			var os pkg_rhino.Node
			_ = os
			if otherSingle != nil {
				os = otherSingle[i]
			} else {
				os = pkg_rhino.Node(0)
			}
			var om []pkg_rhino.Node
			_ = om
			if otherMulti != nil {
				om = otherMulti[i]
			} else {
				om = nil
			}
			if (os == 0) && (om == nil) {
				continue
			}
			var ts pkg_rhino.Node
			_ = ts
			if this_.singleUse != nil {
				ts = this_.singleUse[i]
			} else {
				ts = pkg_rhino.Node(0)
			}
			var tm []pkg_rhino.Node
			_ = tm
			if this_.multiUse != nil {
				tm = this_.multiUse[i]
			} else {
				tm = nil
			}
			if (ts == 0) && (tm == nil) {
				if os != 0 {
					{
						a_27 := this_.priv_ReachingUses_ensureSingleArray_3871()
						i_28 := i
						a_27[i_28] = os
					}
				} else {
					{
						a_29 := this_.priv_ReachingUses_ensureMultiArray_3872()
						i_30 := i
						a_29[i_30] = om
					}
				}
			} else if ts != 0 {
				if os != 0 {
					if ts != os {
						{
							a_31 := this_.singleUse
							i_32 := i
							a_31[i_32] = pkg_rhino.Node(0)
						}
						{
							a_33 := this_.priv_ReachingUses_ensureMultiArray_3872()
							i_34 := i
							a_33[i_34] = MaybeReachingVariableUse_ReachingUses_PairSorted(ts, os)
						}
					}
				} else {
					{
						a_35 := this_.singleUse
						i_36 := i
						a_35[i_36] = pkg_rhino.Node(0)
					}
					{
						a_37 := this_.priv_ReachingUses_ensureMultiArray_3872()
						i_38 := i
						a_37[i_38] = func() []pkg_rhino.Node {
							if MaybeReachingVariableUse_ReachingUses_ContainsSorted(om, ts) {
								return om
							}
							return MaybeReachingVariableUse_ReachingUses_InsertSorted(om, ts)
						}()
					}
				}
			} else {
				if os != 0 {
					if !(MaybeReachingVariableUse_ReachingUses_ContainsSorted(tm, os)) {
						{
							a_39 := this_.multiUse
							i_40 := i
							a_39[i_40] = MaybeReachingVariableUse_ReachingUses_InsertSorted(tm, os)
						}
					}
				} else if !(j2gort.SliceRefEquals(tm, om)) {
					{
						a_41 := this_.multiUse
						i_42 := i
						a_41[i_42] = MaybeReachingVariableUse_ReachingUses_UnionSorted(tm, om)
					}
				}
			}
		}
	}
}

// Java method: ReachingUses.ensureSingleArray()
func (this_ *MaybeReachingVariableUse_ReachingUses) priv_ReachingUses_ensureSingleArray_3871() []pkg_rhino.Node {
	var s []pkg_rhino.Node = this_.singleUse
	_ = s
	if s == nil {
		s = make([]pkg_rhino.Node, this_.numVars)
		this_.singleUse = s
	}
	return s
}

// Java method: ReachingUses.ensureMultiArray()
func (this_ *MaybeReachingVariableUse_ReachingUses) priv_ReachingUses_ensureMultiArray_3872() [][]pkg_rhino.Node {
	var m [][]pkg_rhino.Node = this_.multiUse
	_ = m
	if m == nil {
		m = make([][]pkg_rhino.Node, this_.numVars)
		this_.multiUse = m
	}
	return m
}

// Java method: ReachingUses.equals(java.lang.Object other)
func (this_ *MaybeReachingVariableUse_ReachingUses) Equals(arg_43 any) bool {
	other := arg_43
	_ = other
	var otherUses *MaybeReachingVariableUse_ReachingUses
	_ = otherUses
	var ok_44 bool
	otherUses, ok_44 = j2gort.TryCast[*MaybeReachingVariableUse_ReachingUses](other)
	if !ok_44 {
		return false
	}
	return this_.EqualsUses(otherUses)
}

// Java method: ReachingUses.equalsUses(com.google.javascript.jscomp.MaybeReachingVariableUse.ReachingUses otherUses)
func (this_ *MaybeReachingVariableUse_ReachingUses) EqualsUses(otherUses *MaybeReachingVariableUse_ReachingUses) bool {
	_ = otherUses
	if this_ == otherUses {
		return true
	}
	if this_.numVars != otherUses.numVars {
		return false
	}
	var sa []pkg_rhino.Node = this_.singleUse
	_ = sa
	var sb []pkg_rhino.Node = otherUses.singleUse
	_ = sb
	if !(j2gort.SliceRefEquals(sa, sb)) {
		if sa == nil {
			{
				var i int32 = int32(0)
				_ = i
				var len_ int32 = int32(len(sb))
				_ = len_
				for ; i < len_; i = i + 1 {
					if sb[i] != 0 {
						return false
					}
				}
			}
		} else if sb == nil {
			{
				var i_2 int32 = int32(0)
				_ = i_2
				var len__2 int32 = int32(len(sa))
				_ = len__2
				for ; i_2 < len__2; i_2 = i_2 + 1 {
					if sa[i_2] != 0 {
						return false
					}
				}
			}
		} else {
			{
				var i_3 int32 = int32(0)
				_ = i_3
				var len__3 int32 = this_.numVars
				_ = len__3
				for ; i_3 < len__3; i_3 = i_3 + 1 {
					if sa[i_3] != sb[i_3] {
						return false
					}
				}
			}
		}
	}
	var ma [][]pkg_rhino.Node = this_.multiUse
	_ = ma
	var mb [][]pkg_rhino.Node = otherUses.multiUse
	_ = mb
	if !(j2gort.SliceRefEquals(ma, mb)) {
		if ma == nil {
			{
				var i_4 int32 = int32(0)
				_ = i_4
				var len__4 int32 = int32(len(mb))
				_ = len__4
				for ; i_4 < len__4; i_4 = i_4 + 1 {
					if mb[i_4] != nil {
						return false
					}
				}
			}
		} else if mb == nil {
			{
				var i_5 int32 = int32(0)
				_ = i_5
				var len__5 int32 = int32(len(ma))
				_ = len__5
				for ; i_5 < len__5; i_5 = i_5 + 1 {
					if ma[i_5] != nil {
						return false
					}
				}
			}
		} else {
			{
				var i_6 int32 = int32(0)
				_ = i_6
				var len__6 int32 = this_.numVars
				_ = len__6
				for ; i_6 < len__6; i_6 = i_6 + 1 {
					if !(MaybeReachingVariableUse_ReachingUses_EqualSortedSets(ma[i_6], mb[i_6])) {
						return false
					}
				}
			}
		}
	}
	return true
}

// Java static method: ReachingUses.pairSorted(com.google.javascript.rhino.Node n1, com.google.javascript.rhino.Node n2)
func MaybeReachingVariableUse_ReachingUses_PairSorted(n1 pkg_rhino.Node, n2 pkg_rhino.Node) []pkg_rhino.Node {
	_ = n1
	_ = n2
	if int32(j2gort.CheckNotNullNode(n1)) <= int32(j2gort.CheckNotNullNode(n2)) {
		return []pkg_rhino.Node{n1, n2}
	} else {
		return []pkg_rhino.Node{n2, n1}
	}
}

// Java static method: ReachingUses.containsSorted(com.google.javascript.rhino.Node[] arr, com.google.javascript.rhino.Node target)
func MaybeReachingVariableUse_ReachingUses_ContainsSorted(arr []pkg_rhino.Node, target pkg_rhino.Node) bool {
	_ = arr
	_ = target
	var h int32 = int32(j2gort.CheckNotNullNode(target))
	_ = h
	var lo int32 = int32(0)
	_ = lo
	var hi int32 = (int32(len(arr)) - int32(1))
	_ = hi
	for lo <= hi {
		var mid int32 = javalang.Ushr32((lo + hi), int(int32(1)))
		_ = mid
		var m pkg_rhino.Node = arr[mid]
		_ = m
		if m == target {
			return true
		}
		var mh int32 = int32(j2gort.CheckNotNullNode(m))
		_ = mh
		if mh < h {
			lo = (mid + int32(1))
		} else if mh > h {
			hi = (mid - int32(1))
		} else {
			{
				var k int32 = (mid - int32(1))
				_ = k
				for ; (k >= lo) && (int32(j2gort.CheckNotNullNode(arr[k])) == h); k = k - 1 {
					if arr[k] == target {
						return true
					}
				}
			}
			{
				var k_2 int32 = (mid + int32(1))
				_ = k_2
				for ; (k_2 <= hi) && (int32(j2gort.CheckNotNullNode(arr[k_2])) == h); k_2 = k_2 + 1 {
					if arr[k_2] == target {
						return true
					}
				}
			}
			return false
		}
	}
	return false
}

// Java static method: ReachingUses.insertSorted(com.google.javascript.rhino.Node[] arr, com.google.javascript.rhino.Node target)
func MaybeReachingVariableUse_ReachingUses_InsertSorted(arr []pkg_rhino.Node, target pkg_rhino.Node) []pkg_rhino.Node {
	_ = arr
	_ = target
	var h int32 = int32(j2gort.CheckNotNullNode(target))
	_ = h
	var len_ int32 = int32(len(arr))
	_ = len_
	var lo int32 = int32(0)
	_ = lo
	var hi int32 = len_
	_ = hi
	for lo < hi {
		var mid int32 = javalang.Ushr32((lo + hi), int(int32(1)))
		_ = mid
		if int32(j2gort.CheckNotNullNode(arr[mid])) <= h {
			lo = (mid + int32(1))
		} else {
			hi = mid
		}
	}
	var out []pkg_rhino.Node = make([]pkg_rhino.Node, (len_ + int32(1)))
	_ = out
	if lo > int32(0) {
		func() {
			acSrc_45 := arr
			acSP_46 := int32(0)
			acDst_47 := out
			acDP_48 := int32(0)
			acLen_49 := lo
			j2gort.CheckNotNullSlice(acSrc_45)
			j2gort.CheckNotNullSlice(acDst_47)
			if acSP_46 < 0 || acDP_48 < 0 || acLen_49 < 0 || int64(acSP_46)+int64(acLen_49) > int64(len(acSrc_45)) || int64(acDP_48)+int64(acLen_49) > int64(len(acDst_47)) {
				panic(j2gort.ConstructShimThrowable("java.lang.ArrayIndexOutOfBoundsException"))
			}
			copy(acDst_47[int(acDP_48):int(acDP_48)+int(acLen_49)], acSrc_45[int(acSP_46):int(acSP_46)+int(acLen_49)])
		}()
	}
	{
		a_50 := out
		i_51 := lo
		a_50[i_51] = target
	}
	if lo < len_ {
		func() {
			acSrc_52 := arr
			acSP_53 := lo
			acDst_54 := out
			acDP_55 := (lo + int32(1))
			acLen_56 := (len_ - lo)
			j2gort.CheckNotNullSlice(acSrc_52)
			j2gort.CheckNotNullSlice(acDst_54)
			if acSP_53 < 0 || acDP_55 < 0 || acLen_56 < 0 || int64(acSP_53)+int64(acLen_56) > int64(len(acSrc_52)) || int64(acDP_55)+int64(acLen_56) > int64(len(acDst_54)) {
				panic(j2gort.ConstructShimThrowable("java.lang.ArrayIndexOutOfBoundsException"))
			}
			copy(acDst_54[int(acDP_55):int(acDP_55)+int(acLen_56)], acSrc_52[int(acSP_53):int(acSP_53)+int(acLen_56)])
		}()
	}
	return out
}

// Java static method: ReachingUses.unionSorted(com.google.javascript.rhino.Node[] a, com.google.javascript.rhino.Node[] b)
func MaybeReachingVariableUse_ReachingUses_UnionSorted(a []pkg_rhino.Node, b []pkg_rhino.Node) []pkg_rhino.Node {
	_ = a
	_ = b
	var lenA int32 = int32(len(a))
	_ = lenA
	var lenB int32 = int32(len(b))
	_ = lenB
	var mergedLen int32 = MaybeReachingVariableUse_ReachingUses_CountMergedLength(a, lenA, b, lenB)
	_ = mergedLen
	if mergedLen == lenA {
		return a
	}
	if mergedLen == lenB {
		return b
	}
	var out []pkg_rhino.Node = make([]pkg_rhino.Node, mergedLen)
	_ = out
	var ia int32 = int32(0)
	_ = ia
	var ib int32 = int32(0)
	_ = ib
	var outIdx int32 = int32(0)
	_ = outIdx
	for (ia < lenA) && (ib < lenB) {
		var na pkg_rhino.Node = a[ia]
		_ = na
		var nb pkg_rhino.Node = b[ib]
		_ = nb
		var ha int32 = int32(j2gort.CheckNotNullNode(na))
		_ = ha
		var hb int32 = int32(j2gort.CheckNotNullNode(nb))
		_ = hb
		if ha < hb {
			{
				a_57 := out
				var i_58 int32
				{

					old_59 := outIdx
					nw_60 := old_59 + 1
					outIdx = nw_60
					i_58 = old_59
				}
				a_57[i_58] = na
			}
			ia = ia + 1
		} else if hb < ha {
			{
				a_61 := out
				var i_62 int32
				{

					old_63 := outIdx
					nw_64 := old_63 + 1
					outIdx = nw_64
					i_62 = old_63
				}
				a_61[i_62] = nb
			}
			ib = ib + 1
		} else if na == nb {
			{
				a_65 := out
				var i_66 int32
				{

					old_67 := outIdx
					nw_68 := old_67 + 1
					outIdx = nw_68
					i_66 = old_67
				}
				a_65[i_66] = na
			}
			ia = ia + 1
			ib = ib + 1
		} else {
			var endA int32 = (ia + int32(1))
			_ = endA
			for (endA < lenA) && (int32(j2gort.CheckNotNullNode(a[endA])) == ha) {
				endA = endA + 1
			}
			var endB int32 = (ib + int32(1))
			_ = endB
			for (endB < lenB) && (int32(j2gort.CheckNotNullNode(b[endB])) == hb) {
				endB = endB + 1
			}
			{
				var k int32 = ia
				_ = k
				for ; k < endA; k = k + 1 {
					{
						a_69 := out
						var i_70 int32
						{

							old_71 := outIdx
							nw_72 := old_71 + 1
							outIdx = nw_72
							i_70 = old_71
						}
						a_69[i_70] = a[k]
					}
				}
			}
			{
				var k_2 int32 = ib
				_ = k_2
				for ; k_2 < endB; k_2 = k_2 + 1 {
					var cand pkg_rhino.Node = b[k_2]
					_ = cand
					var dup bool = false
					_ = dup
					{
						var j int32 = ia
						_ = j
						for ; j < endA; j = j + 1 {
							if a[j] == cand {
								dup = true
								break
							}
						}
					}
					if !(dup) {
						{
							a_73 := out
							var i_74 int32
							{

								old_75 := outIdx
								nw_76 := old_75 + 1
								outIdx = nw_76
								i_74 = old_75
							}
							a_73[i_74] = cand
						}
					}
				}
			}
			ia = endA
			ib = endB
		}
	}
	if ia < lenA {
		func() {
			acSrc_77 := a
			acSP_78 := ia
			acDst_79 := out
			acDP_80 := outIdx
			acLen_81 := (lenA - ia)
			j2gort.CheckNotNullSlice(acSrc_77)
			j2gort.CheckNotNullSlice(acDst_79)
			if acSP_78 < 0 || acDP_80 < 0 || acLen_81 < 0 || int64(acSP_78)+int64(acLen_81) > int64(len(acSrc_77)) || int64(acDP_80)+int64(acLen_81) > int64(len(acDst_79)) {
				panic(j2gort.ConstructShimThrowable("java.lang.ArrayIndexOutOfBoundsException"))
			}
			copy(acDst_79[int(acDP_80):int(acDP_80)+int(acLen_81)], acSrc_77[int(acSP_78):int(acSP_78)+int(acLen_81)])
		}()
	} else if ib < lenB {
		func() {
			acSrc_82 := b
			acSP_83 := ib
			acDst_84 := out
			acDP_85 := outIdx
			acLen_86 := (lenB - ib)
			j2gort.CheckNotNullSlice(acSrc_82)
			j2gort.CheckNotNullSlice(acDst_84)
			if acSP_83 < 0 || acDP_85 < 0 || acLen_86 < 0 || int64(acSP_83)+int64(acLen_86) > int64(len(acSrc_82)) || int64(acDP_85)+int64(acLen_86) > int64(len(acDst_84)) {
				panic(j2gort.ConstructShimThrowable("java.lang.ArrayIndexOutOfBoundsException"))
			}
			copy(acDst_84[int(acDP_85):int(acDP_85)+int(acLen_86)], acSrc_82[int(acSP_83):int(acSP_83)+int(acLen_86)])
		}()
	}
	return out
}

// Java static method: ReachingUses.countMergedLength(com.google.javascript.rhino.Node[] a, int lenA, com.google.javascript.rhino.Node[] b, int lenB)
func MaybeReachingVariableUse_ReachingUses_CountMergedLength(a []pkg_rhino.Node, lenA int32, b []pkg_rhino.Node, lenB int32) int32 {
	_ = a
	_ = lenA
	_ = b
	_ = lenB
	var ia int32 = int32(0)
	_ = ia
	var ib int32 = int32(0)
	_ = ib
	var count int32 = int32(0)
	_ = count
	for (ia < lenA) && (ib < lenB) {
		var na pkg_rhino.Node = a[ia]
		_ = na
		var nb pkg_rhino.Node = b[ib]
		_ = nb
		var ha int32 = int32(j2gort.CheckNotNullNode(na))
		_ = ha
		var hb int32 = int32(j2gort.CheckNotNullNode(nb))
		_ = hb
		if ha < hb {
			count = count + 1
			ia = ia + 1
		} else if hb < ha {
			count = count + 1
			ib = ib + 1
		} else if na == nb {
			count = count + 1
			ia = ia + 1
			ib = ib + 1
		} else {
			var endA int32 = (ia + int32(1))
			_ = endA
			for (endA < lenA) && (int32(j2gort.CheckNotNullNode(a[endA])) == ha) {
				endA = endA + 1
			}
			var endB int32 = (ib + int32(1))
			_ = endB
			for (endB < lenB) && (int32(j2gort.CheckNotNullNode(b[endB])) == hb) {
				endB = endB + 1
			}
			count = (count + (endA - ia))
			{
				var k int32 = ib
				_ = k
				for ; k < endB; k = k + 1 {
					var cand pkg_rhino.Node = b[k]
					_ = cand
					var dup bool = false
					_ = dup
					{
						var j int32 = ia
						_ = j
						for ; j < endA; j = j + 1 {
							if a[j] == cand {
								dup = true
								break
							}
						}
					}
					if !(dup) {
						count = count + 1
					}
				}
			}
			ia = endA
			ib = endB
		}
	}
	return ((count + (lenA - ia)) + (lenB - ib))
}

// Java static method: ReachingUses.equalSortedSets(com.google.javascript.rhino.Node @org.jspecify.annotations.Nullable [] a, com.google.javascript.rhino.Node @org.jspecify.annotations.Nullable [] b)
func MaybeReachingVariableUse_ReachingUses_EqualSortedSets(a []pkg_rhino.Node, b []pkg_rhino.Node) bool {
	_ = a
	_ = b
	if j2gort.SliceRefEquals(a, b) {
		return true
	}
	if (a == nil) || (b == nil) {
		return false
	}
	var len_ int32 = int32(len(a))
	_ = len_
	if len_ != int32(len(b)) {
		return false
	}
	var i int32 = int32(0)
	_ = i
	for i < len_ {
		var na pkg_rhino.Node = a[i]
		_ = na
		var nb pkg_rhino.Node = b[i]
		_ = nb
		if na == nb {
			i = i + 1
			continue
		}
		var h int32 = int32(j2gort.CheckNotNullNode(na))
		_ = h
		if int32(j2gort.CheckNotNullNode(nb)) != h {
			return false
		}
		var end int32 = (i + int32(1))
		_ = end
		for (end < len_) && (int32(j2gort.CheckNotNullNode(a[end])) == h) {
			if int32(j2gort.CheckNotNullNode(b[end])) != h {
				return false
			}
			end = end + 1
		}
		if (end < len_) && (int32(j2gort.CheckNotNullNode(b[end])) == h) {
			return false
		}
		{
			var k int32 = i
			_ = k
			for ; k < end; k = k + 1 {
				var target pkg_rhino.Node = a[k]
				_ = target
				var found bool = false
				_ = found
				{
					var j int32 = i
					_ = j
					for ; j < end; j = j + 1 {
						if b[j] == target {
							found = true
							break
						}
					}
				}
				if !(found) {
					return false
				}
			}
		}
		i = end
	}
	return true
}

// Java method: ReachingUses.hashCode()
func (this_ *MaybeReachingVariableUse_ReachingUses) HashCode() int32 {
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException", "the hashcode of this object is not stable")))
}

// Java class: com.google.javascript.jscomp.MaybeReachingVariableUse.ReachingUsesJoinOp
type MaybeReachingVariableUse_ReachingUsesJoinOp struct {
	result  *MaybeReachingVariableUse_ReachingUses `gson:"result"`
	idHash_ atomic.Int32
}

func (this_ *MaybeReachingVariableUse_ReachingUsesJoinOp) BindSelf_(s any) {
	_ = s
}
func (this_ *MaybeReachingVariableUse_ReachingUsesJoinOp) IsJavaNil_() bool { return this_ == nil }
func (*MaybeReachingVariableUse_ReachingUsesJoinOp) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.MaybeReachingVariableUse$ReachingUsesJoinOp", "ReachingUsesJoinOp"
}
func (this_ *MaybeReachingVariableUse_ReachingUsesJoinOp) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*MaybeReachingVariableUse_ReachingUsesJoinOp) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.MaybeReachingVariableUse.ReachingUsesJoinOp", "com.google.javascript.jscomp.MaybeReachingVariableUse$ReachingUsesJoinOp", func(v_ any) bool {
		return j2gort.InstanceOf[*MaybeReachingVariableUse_ReachingUsesJoinOp](v_)
	})
}
func (*MaybeReachingVariableUse_ReachingUsesJoinOp) Implements_jscomp_DataFlowAnalysis_FlowJoiner() {}

func (this_ *MaybeReachingVariableUse_ReachingUsesJoinOp) initFields_MaybeReachingVariableUse_ReachingUsesJoinOp() {
	_ = this_
}

// Java constructor: ReachingUsesJoinOp(int numVars)
func NewMaybeReachingVariableUse_ReachingUsesJoinOp(numVars int32) *MaybeReachingVariableUse_ReachingUsesJoinOp {
	this_ := &MaybeReachingVariableUse_ReachingUsesJoinOp{}
	InitMaybeReachingVariableUse_ReachingUsesJoinOp(this_, numVars)
	return this_
}
func allocPooled_MaybeReachingVariableUse_ReachingUsesJoinOp(slab_ *j2gort.PooledSlab[MaybeReachingVariableUse_ReachingUsesJoinOp], slabUses_ *j2gort.PooledSlab[MaybeReachingVariableUse_ReachingUses], numVars int32) *MaybeReachingVariableUse_ReachingUsesJoinOp {
	this_ := slab_.Alloc()
	this_.initFields_MaybeReachingVariableUse_ReachingUsesJoinOp()
	this_.result = allocPooled_MaybeReachingVariableUse_ReachingUses_1_int32(slabUses_, numVars)
	return this_
}

func InitMaybeReachingVariableUse_ReachingUsesJoinOp(this_ *MaybeReachingVariableUse_ReachingUsesJoinOp, numVars int32) {
	_ = numVars
	this_.initFields_MaybeReachingVariableUse_ReachingUsesJoinOp()
	this_.result = NewMaybeReachingVariableUse_ReachingUses_1_int32(numVars)
}

// Java method: ReachingUsesJoinOp.joinFlow(com.google.javascript.jscomp.MaybeReachingVariableUse.ReachingUses uses)
func (this_ *MaybeReachingVariableUse_ReachingUsesJoinOp) JoinFlow(uses *MaybeReachingVariableUse_ReachingUses) {
	_ = uses
	this_.result.Join(uses)
}

// Java method: ReachingUsesJoinOp.finish()
func (this_ *MaybeReachingVariableUse_ReachingUsesJoinOp) Finish() *MaybeReachingVariableUse_ReachingUses {
	return this_.result
}
func (this_ *MaybeReachingVariableUse) ReleaseSlabs_() {
	if this_ == nil {
		return
	}
	this_.slab_ReachingUses.Release()
	this_.slab_ReachingUsesJoinOp.Release()
}
