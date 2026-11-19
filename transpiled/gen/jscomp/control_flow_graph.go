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
// Port of com/google/javascript/jscomp/ControlFlowGraph.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_graph "github.com/paullepoulpe/jsgo/transpiled/gen/graph"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java polymorphic class interface: com.google.javascript.jscomp.ControlFlowGraph
type ControlFlowGraph[N comparable] interface {
	Implements_jscomp_ControlFlowGraph()
	IsJavaNil_() bool
	GetControlFlowGraphBase() *ControlFlowGraphBase[N]
	pkg_graph.LinkedDirectedGraph[N, ControlFlowGraph_Branch]
	ResetCfgForPool() bool
	ReinitCfg(N, bool)
	GetImplicitReturn() pkg_graph.DiGraph_DiGraphNode[N, ControlFlowGraph_Branch]
	GetEntry() pkg_graph.DiGraph_DiGraphNode[N, ControlFlowGraph_Branch]
	GetEntryNode() *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch]
	IsImplicitReturn(pkg_graph.DiGraph_DiGraphNode[N, ControlFlowGraph_Branch]) bool
	GetOptionalNodeComparator(bool) j2gort.Ext_java_util_Comparator
	HasPriorityComparator() bool
	String() string
}

func ToIface_ControlFlowGraph[N comparable, T_iface_ interface {
	ControlFlowGraph[N]
	comparable
}](v T_iface_) ControlFlowGraph[N] {
	var z T_iface_
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.jscomp.ControlFlowGraph
type ControlFlowGraphBase[N comparable] struct {
	pkg_graph.LinkedDirectedGraphBase[N, ControlFlowGraph_Branch]
	implicitReturn       *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch] `gson:"implicitReturn"`
	entry                *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch] `gson:"entry"`
	slab_LinearFlowState j2gort.PooledSlab[DataFlowAnalysis_LinearFlowState]
}

func (this_ *ControlFlowGraphBase[N]) BindSelf_(s any) {
	_ = s
	this_.LinkedDirectedGraphBase.BindSelf_(s)
}
func (this_ *ControlFlowGraphBase[N]) IsJavaNil_() bool { return this_ == nil }
func (*ControlFlowGraphBase[N]) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ControlFlowGraph", "ControlFlowGraph"
}
func (this_ *ControlFlowGraphBase[N]) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.LinkedDirectedGraphBase.IdentityHash_()
}
func (*ControlFlowGraphBase[N]) IsGenericInst_jscomp_ControlFlowGraph_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ControlFlowGraph", "com.google.javascript.jscomp.ControlFlowGraph", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			IsGenericInst_jscomp_ControlFlowGraph_()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}
func (b *ControlFlowGraphBase[N]) GetControlFlowGraphBase() *ControlFlowGraphBase[N] { return b }
func (*ControlFlowGraphBase[N]) Implements_jscomp_ControlFlowGraph()                 {}
func (*ControlFlowGraphBase[N]) Implements_graph_LinkedDirectedGraph()               {}
func (*ControlFlowGraphBase[N]) Implements_graph_DiGraph()                           {}
func (*ControlFlowGraphBase[N]) Implements_graph_GraphvizGraph()                     {}
func (*ControlFlowGraphBase[N]) Implements_graph_Graph()                             {}
func (*ControlFlowGraphBase[N]) Implements_graph_AdjacencyGraph()                    {}

func (this_ *ControlFlowGraphBase[N]) initFields_ControlFlowGraph() {
	_ = this_
}

// Java constructor: ControlFlowGraph(N entry, boolean nodeAnnotations, boolean edgeAnnotations)
func NewControlFlowGraph[N comparable](entry N, nodeAnnotations bool, edgeAnnotations bool) ControlFlowGraph[N] {
	this_ := &ControlFlowGraphBase[N]{}
	this_.BindSelf_(this_)
	InitControlFlowGraph[N](this_, entry, nodeAnnotations, edgeAnnotations)
	return this_
}

func InitControlFlowGraph[N comparable](this_ *ControlFlowGraphBase[N], entry N, nodeAnnotations bool, edgeAnnotations bool) {
	_ = entry
	_ = nodeAnnotations
	_ = edgeAnnotations
	pkg_graph.InitLinkedDirectedGraph[N, ControlFlowGraph_Branch](&this_.LinkedDirectedGraphBase, nodeAnnotations, edgeAnnotations)
	this_.initFields_ControlFlowGraph()
	this_.implicitReturn = j2gort.CastPtr[pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch]](this_.CreateNode(func() N {
		if _, ok_ := any((*N)(nil)).(*string); ok_ {
			return any(j2gort.NullString).(N)
		}
		var z_ N
		return z_
	}()))
	this_.entry = j2gort.CastPtr[pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch]](this_.CreateNode(entry))
}

// Java method: ControlFlowGraph.resetCfgForPool()
func (this_ *ControlFlowGraphBase[N]) ResetCfgForPool() bool {
	this_.slab_LinearFlowState.Release()
	this_.implicitReturn = nil
	this_.entry = nil
	return this_.ResetForReuse(true, false)
}

// Java method: ControlFlowGraph.reinitCfg(N entryValue, boolean edgeAnnotations)
func (this_ *ControlFlowGraphBase[N]) ReinitCfg(entryValue N, edgeAnnotations bool) {
	_ = entryValue
	_ = edgeAnnotations
	this_.SetUseEdgeAnnotations(edgeAnnotations)
	this_.implicitReturn = j2gort.CastPtr[pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch]](this_.CreateNode(func() N {
		if _, ok_ := any((*N)(nil)).(*string); ok_ {
			return any(j2gort.NullString).(N)
		}
		var z_ N
		return z_
	}()))
	this_.entry = j2gort.CastPtr[pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch]](this_.CreateNode(entryValue))
}

// Java method: ControlFlowGraph.getImplicitReturn()
func (this_ *ControlFlowGraphBase[N]) GetImplicitReturn() pkg_graph.DiGraph_DiGraphNode[N, ControlFlowGraph_Branch] {
	return pkg_graph.ToIface_DiGraph_DiGraphNode[N, ControlFlowGraph_Branch](this_.implicitReturn)
}

// Java method: ControlFlowGraph.getEntry()
func (this_ *ControlFlowGraphBase[N]) GetEntry() pkg_graph.DiGraph_DiGraphNode[N, ControlFlowGraph_Branch] {
	return pkg_graph.ToIface_DiGraph_DiGraphNode[N, ControlFlowGraph_Branch](this_.entry)
}

// Java method: ControlFlowGraph.getEntryNode()
func (this_ *ControlFlowGraphBase[N]) GetEntryNode() *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch] {
	return this_.entry
}

// Java method: ControlFlowGraph.isImplicitReturn(com.google.javascript.jscomp.graph.DiGraph.DiGraphNode<N,com.google.javascript.jscomp.ControlFlowGraph.Branch> node)
func (this_ *ControlFlowGraphBase[N]) IsImplicitReturn(node pkg_graph.DiGraph_DiGraphNode[N, ControlFlowGraph_Branch]) bool {
	_ = node
	return node == this_.implicitReturn || ((node == nil || node.IsJavaNil_()) && (this_.implicitReturn == nil || this_.implicitReturn.IsJavaNil_(

	// Java method: ControlFlowGraph.getOptionalNodeComparator(boolean isForward)
	)))
}

func (this_ *ControlFlowGraphBase[N]) GetOptionalNodeComparator(isForward bool) j2gort.Ext_java_util_Comparator {
	_ = isForward
	return nil
}

// Java method: ControlFlowGraph.hasPriorityComparator()
func (this_ *ControlFlowGraphBase[N]) HasPriorityComparator() bool {
	return false
}

// Java static method: ControlFlowGraph.isEnteringNewCfgNode(com.google.javascript.rhino.Node n)
func ControlFlowGraph_IsEnteringNewCfgNode(n pkg_rhino.Node) bool {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	switch parent.TokenSwitchOrd_() {
	case int32(90) /* Token.BLOCK */, int32(89) /* Token.ROOT */, int32(94) /* Token.SCRIPT */, int32(48) /* Token.TRY */, int32(91) /* Token.SWITCH_BODY */ :
		return true
	case int32(73) /* Token.FUNCTION */ :
		return (n != parent.GetSecondChild())
	case int32(78) /* Token.WHILE */, int32(79) /* Token.DO */, int32(74) /* Token.IF */ :
		return (NodeUtil_GetConditionExpression(parent) != n)
	case int32(80) /* Token.FOR */ :
		return (NodeUtil_GetConditionExpression(parent) != n)
	case int32(81) /* Token.FOR_IN */, int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */ :
		return (n != parent.GetFirstChild())
	case int32(76) /* Token.CASE */, int32(86) /* Token.CATCH */, int32(85) /* Token.WITH */ :
		return (n != parent.GetFirstChild())
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: ControlFlowGraph.toString()
func (this_ *ControlFlowGraphBase[N]) String() string {
	var s string = "CFG:\n"
	_ = s
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.GetGraphvizEdges()))
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			e := j2gort.AsIface[pkg_graph.GraphvizGraph_GraphvizEdge](rawElem_2)
			_ = e
			s = (j2gort.StrValString(s) + javalang.ToJavaString(j2gort.CheckNotNull(e)) + "\n")
		}
	}
	return s
}

// Java enum: com.google.javascript.jscomp.ControlFlowGraph.Branch
type ControlFlowGraph_Branch uint8

const (
	enum_ControlFlowGraph_Branch_ON_TRUE   ControlFlowGraph_Branch = 1
	enum_ControlFlowGraph_Branch_ON_FALSE  ControlFlowGraph_Branch = 2
	enum_ControlFlowGraph_Branch_UNCOND    ControlFlowGraph_Branch = 3
	enum_ControlFlowGraph_Branch_ON_EX     ControlFlowGraph_Branch = 4
	enum_ControlFlowGraph_Branch_SYN_BLOCK ControlFlowGraph_Branch = 5
)

var enumBaseTable_ControlFlowGraph_Branch = [...]j2gort.EnumBase{
	{Name_: "ON_TRUE", Ordinal_: 0},
	{Name_: "ON_FALSE", Ordinal_: 1},
	{Name_: "UNCOND", Ordinal_: 2},
	{Name_: "ON_EX", Ordinal_: 3},
	{Name_: "SYN_BLOCK", Ordinal_: 4},
}
var enumValues_ControlFlowGraph_Branch = []ControlFlowGraph_Branch{1, 2, 3, 4, 5}

func StaticGet_ControlFlowGraph_Branch_ON_TRUE() ControlFlowGraph_Branch {
	return enum_ControlFlowGraph_Branch_ON_TRUE
}
func StaticGet_ControlFlowGraph_Branch_ON_FALSE() ControlFlowGraph_Branch {
	return enum_ControlFlowGraph_Branch_ON_FALSE
}
func StaticGet_ControlFlowGraph_Branch_UNCOND() ControlFlowGraph_Branch {
	return enum_ControlFlowGraph_Branch_UNCOND
}
func StaticGet_ControlFlowGraph_Branch_ON_EX() ControlFlowGraph_Branch {
	return enum_ControlFlowGraph_Branch_ON_EX
}
func StaticGet_ControlFlowGraph_Branch_SYN_BLOCK() ControlFlowGraph_Branch {
	return enum_ControlFlowGraph_Branch_SYN_BLOCK
}

func ControlFlowGraph_Branch_Values() []ControlFlowGraph_Branch {
	return j2gort.CloneSlice(enumValues_ControlFlowGraph_Branch)
}

func ControlFlowGraph_Branch_ValueOf(name string) ControlFlowGraph_Branch {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_ControlFlowGraph_Branch {
		if enumBaseTable_ControlFlowGraph_Branch[i_].Name_ == name {
			return ControlFlowGraph_Branch(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.ControlFlowGraph.Branch." + name))
}

// Java method: Branch.isConditional()
func (this_ ControlFlowGraph_Branch) IsConditional() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return ((this_ == StaticGet_ControlFlowGraph_Branch_ON_TRUE()) || (this_ == StaticGet_ControlFlowGraph_Branch_ON_FALSE()))
}

func (this_ ControlFlowGraph_Branch) IsJavaNil_() bool                              { return this_ == 0 }
func (ControlFlowGraph_Branch) HasIdentityEquals_()                                 {}
func (ControlFlowGraph_Branch) Implements_j2gort_Ext_java_lang_Enum()               {}
func (ControlFlowGraph_Branch) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ ControlFlowGraph_Branch) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (ControlFlowGraph_Branch) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ControlFlowGraph$Branch", "Branch"
}
func (this_ ControlFlowGraph_Branch) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_ControlFlowGraph_Branch[this_-1]
}
func (this_ ControlFlowGraph_Branch) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ ControlFlowGraph_Branch) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ ControlFlowGraph_Branch) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ControlFlowGraph_Branch[this_-1].Name_
}
func (this_ ControlFlowGraph_Branch) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_ControlFlowGraph_Branch[this_-1].Name_
}
func (this_ ControlFlowGraph_Branch) Equals(other_ any) bool {
	o_, ok_ := other_.(ControlFlowGraph_Branch)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ ControlFlowGraph_Branch) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_ControlFlowGraph_Branch[this_-1].HashCode()
}
func (this_ ControlFlowGraph_Branch) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ControlFlowGraph_Branch[this_-1].CompareTo(other_)
}
func (this_ ControlFlowGraph_Branch) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_ControlFlowGraph_Branch[this_-1].Name_)
}
func init() {
	vals_ := []any{ControlFlowGraph_Branch(1), ControlFlowGraph_Branch(2), ControlFlowGraph_Branch(3), ControlFlowGraph_Branch(4), ControlFlowGraph_Branch(5)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.ControlFlowGraph.Branch", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.ControlFlowGraph$Branch", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ControlFlowGraph.Branch", "com.google.javascript.jscomp.ControlFlowGraph$Branch", func(v_ any) bool {
		return j2gort.InstanceOf[ControlFlowGraph_Branch](v_)
	})
}

// Java polymorphic class interface: com.google.javascript.jscomp.ControlFlowGraph.AbstractCfgNodeTraversalCallback
type ControlFlowGraph_AbstractCfgNodeTraversalCallback interface {
	Implements_jscomp_ControlFlowGraph_AbstractCfgNodeTraversalCallback()
	IsJavaNil_() bool
	GetControlFlowGraph_AbstractCfgNodeTraversalCallbackBase() *ControlFlowGraph_AbstractCfgNodeTraversalCallbackBase
	NodeTraversal_Callback
}

func ToIface_ControlFlowGraph_AbstractCfgNodeTraversalCallback[T interface {
	ControlFlowGraph_AbstractCfgNodeTraversalCallback
	comparable
}](v T) ControlFlowGraph_AbstractCfgNodeTraversalCallback {
	var z T
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.jscomp.ControlFlowGraph.AbstractCfgNodeTraversalCallback
type ControlFlowGraph_AbstractCfgNodeTraversalCallbackBase struct {
	idHash_ atomic.Int32
}

func (this_ *ControlFlowGraph_AbstractCfgNodeTraversalCallbackBase) BindSelf_(s any) {
	_ = s
}
func (this_ *ControlFlowGraph_AbstractCfgNodeTraversalCallbackBase) IsJavaNil_() bool {
	return this_ == nil
}
func (*ControlFlowGraph_AbstractCfgNodeTraversalCallbackBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ControlFlowGraph$AbstractCfgNodeTraversalCallback", "AbstractCfgNodeTraversalCallback"
}
func (this_ *ControlFlowGraph_AbstractCfgNodeTraversalCallbackBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ControlFlowGraph.AbstractCfgNodeTraversalCallback", "com.google.javascript.jscomp.ControlFlowGraph$AbstractCfgNodeTraversalCallback", func(v_ any) bool {
		return j2gort.InstanceOf[ControlFlowGraph_AbstractCfgNodeTraversalCallback](v_)
	})
}
func (b *ControlFlowGraph_AbstractCfgNodeTraversalCallbackBase) GetControlFlowGraph_AbstractCfgNodeTraversalCallbackBase() *ControlFlowGraph_AbstractCfgNodeTraversalCallbackBase {
	return b
}
func (*ControlFlowGraph_AbstractCfgNodeTraversalCallbackBase) Implements_jscomp_ControlFlowGraph_AbstractCfgNodeTraversalCallback() {
}
func (*ControlFlowGraph_AbstractCfgNodeTraversalCallbackBase) Implements_jscomp_NodeTraversal_Callback() {
}

func (this_ *ControlFlowGraph_AbstractCfgNodeTraversalCallbackBase) initFields_ControlFlowGraph_AbstractCfgNodeTraversalCallback() {
	_ = this_
}

func InitControlFlowGraph_AbstractCfgNodeTraversalCallback(this_ *ControlFlowGraph_AbstractCfgNodeTraversalCallbackBase) {
	this_.initFields_ControlFlowGraph_AbstractCfgNodeTraversalCallback()
}

// Java method: AbstractCfgNodeTraversalCallback.shouldTraverse(com.google.javascript.jscomp.NodeTraversal nodeTraversal, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ControlFlowGraph_AbstractCfgNodeTraversalCallbackBase) ShouldTraverse_3(nodeTraversal *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = nodeTraversal
	_ = n
	_ = parent
	if parent == 0 {
		return true
	}
	return !(ControlFlowGraph_IsEnteringNewCfgNode(n))
}
func (this_ *ControlFlowGraphBase[N]) ReleaseSlabs_() {
	if this_ == nil {
		return
	}
	this_.slab_LinearFlowState.Release()
	this_.LinkedDirectedGraphBase.ReleaseSlabs_()
}
