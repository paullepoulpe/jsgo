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
// Port of com/google/javascript/jscomp/graph/FixedPointGraphTraversal.java (generated mechanically by j2go).

package graph

import (
	"math"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.graph.FixedPointGraphTraversal
type FixedPointGraphTraversal[N comparable, E any] struct {
	callback           FixedPointGraphTraversal_EdgeCallback[N, E] `gson:"callback"`
	idHash_            atomic.Int32
	traversalDirection FixedPointGraphTraversal_TraversalDirection `gson:"traversalDirection"`
}

func (this_ *FixedPointGraphTraversal[N, E]) BindSelf_(s any) {
	_ = s
}
func (this_ *FixedPointGraphTraversal[N, E]) IsJavaNil_() bool { return this_ == nil }
func (*FixedPointGraphTraversal[N, E]) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.graph.FixedPointGraphTraversal", "FixedPointGraphTraversal"
}
func (this_ *FixedPointGraphTraversal[N, E]) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*FixedPointGraphTraversal[N, E]) HasIdentityEquals_()                            {}
func (*FixedPointGraphTraversal[N, E]) IsGenericInst_graph_FixedPointGraphTraversal_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.graph.FixedPointGraphTraversal", "com.google.javascript.jscomp.graph.FixedPointGraphTraversal", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			IsGenericInst_graph_FixedPointGraphTraversal_()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}

func StaticGet_FixedPointGraphTraversal_NON_HALTING_ERROR_MSG() string {
	return "Fixed point computation not halting"
}

var static_FixedPointGraphTraversal_MAX_NODE_COUNT_FOR_ITERATION_LIMIT int64

func StaticGet_FixedPointGraphTraversal_MAX_NODE_COUNT_FOR_ITERATION_LIMIT() int64 {
	clinit_FixedPointGraphTraversal()
	return static_FixedPointGraphTraversal_MAX_NODE_COUNT_FOR_ITERATION_LIMIT
}
func StaticSet_FixedPointGraphTraversal_MAX_NODE_COUNT_FOR_ITERATION_LIMIT(v int64) {
	clinit_FixedPointGraphTraversal()
	static_FixedPointGraphTraversal_MAX_NODE_COUNT_FOR_ITERATION_LIMIT = v
}

var clinitDone_FixedPointGraphTraversal uint32
var clinitState_FixedPointGraphTraversal j2gort.ClinitState

func clinit_FixedPointGraphTraversal() {
	if atomic.LoadUint32(&clinitDone_FixedPointGraphTraversal) != 0 {
		return
	}
	clinitSlow_FixedPointGraphTraversal()
}

func clinitSlow_FixedPointGraphTraversal() {
	if !j2gort.EnterClinit(&clinitState_FixedPointGraphTraversal, &clinitDone_FixedPointGraphTraversal) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_FixedPointGraphTraversal, &clinitDone_FixedPointGraphTraversal)
	static_FixedPointGraphTraversal_MAX_NODE_COUNT_FOR_ITERATION_LIMIT = javalang.D2L(math.Floor(j2gort.Ext_java_lang_Math_Cbrt(float64(6.0e10))))
}

func (this_ *FixedPointGraphTraversal[N, E]) initFields_FixedPointGraphTraversal() {
	_ = this_
}

// Java constructor: FixedPointGraphTraversal(com.google.javascript.jscomp.graph.FixedPointGraphTraversal.EdgeCallback<N,E> callback, com.google.javascript.jscomp.graph.FixedPointGraphTraversal.TraversalDirection traversalDirection)
func NewFixedPointGraphTraversal[N comparable, E any](callback FixedPointGraphTraversal_EdgeCallback[N, E], traversalDirection FixedPointGraphTraversal_TraversalDirection) *FixedPointGraphTraversal[N, E] {
	clinit_FixedPointGraphTraversal()
	this_ := &FixedPointGraphTraversal[N, E]{}
	InitFixedPointGraphTraversal[N, E](this_, callback, traversalDirection)
	return this_
}

func InitFixedPointGraphTraversal[N comparable, E any](this_ *FixedPointGraphTraversal[N, E], callback FixedPointGraphTraversal_EdgeCallback[N, E], traversalDirection FixedPointGraphTraversal_TraversalDirection) {
	_ = callback
	_ = traversalDirection
	this_.initFields_FixedPointGraphTraversal()
	this_.callback = callback
	this_.traversalDirection = traversalDirection
}

// Java static method: FixedPointGraphTraversal.newTraversal(com.google.javascript.jscomp.graph.FixedPointGraphTraversal.EdgeCallback<NodeT,EdgeT> callback)
func FixedPointGraphTraversal_NewTraversal[NodeT comparable, EdgeT any](callback FixedPointGraphTraversal_EdgeCallback[NodeT, EdgeT]) *FixedPointGraphTraversal[NodeT, EdgeT] {
	clinit_FixedPointGraphTraversal()
	_ = callback
	return NewFixedPointGraphTraversal[NodeT, EdgeT](callback, StaticGet_FixedPointGraphTraversal_TraversalDirection_OUTWARDS())
}

// Java static method: FixedPointGraphTraversal.newReverseTraversal(com.google.javascript.jscomp.graph.FixedPointGraphTraversal.EdgeCallback<NodeT,EdgeT> callback)
func FixedPointGraphTraversal_NewReverseTraversal[NodeT comparable, EdgeT any](callback FixedPointGraphTraversal_EdgeCallback[NodeT, EdgeT]) *FixedPointGraphTraversal[NodeT, EdgeT] {
	clinit_FixedPointGraphTraversal()
	_ = callback
	return NewFixedPointGraphTraversal[NodeT, EdgeT](callback, StaticGet_FixedPointGraphTraversal_TraversalDirection_INWARDS())
}

// Java method: FixedPointGraphTraversal.computeFixedPoint(com.google.javascript.jscomp.graph.DiGraph<N,E> graph)
func (this_ *FixedPointGraphTraversal[N, E]) ComputeFixedPoint_1(graph DiGraph[N, E]) {
	_ = graph
	{
		var iifeCond_1 bool
		{

			v_, ok_ := graph.(interface {
				Implements_graph_LinkedDirectedGraph()
				IsJavaNil_() bool
			})
			iifeCond_1 = ok_ && !v_.IsJavaNil_()
		}
		if iifeCond_1 {
			var workQueue j2gort.Deque[DiGraph_DiGraphNode[N, E]] = j2gort.NewDequeCap[DiGraph_DiGraphNode[N, E]](max(graph.GetNodeCount(), int32(8)))
			_ = workQueue
			{
				cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(graph.GetNodes()))
				for cur_1.HasNext() {
					rawElem_2 := cur_1.Next()
					node := j2gort.AsIface[DiGraph_DiGraphNode[N, E]](rawElem_2)
					_ = node
					if !(node.IsInDfaQueue()) {
						node.SetInDfaQueue(true)
						workQueue.AddLast(node)
					}
				}
			}
			this_.priv_FixedPointGraphTraversal_runWorkQueueFixedPoint_6287(graph, &workQueue)
			return
		}
	}
	var workSet *j2gort.Ext_java_util_LinkedHashSet = j2gort.NewExt_java_util_LinkedHashSet_1_int32(graph.GetNodeCount())
	_ = workSet
	workSet.AddAll_1(graph.GetNodes())
	this_.priv_FixedPointGraphTraversal_runWorkSetFixedPoint_6288(graph, workSet)
}

// Java method: FixedPointGraphTraversal.computeFixedPoint(com.google.javascript.jscomp.graph.DiGraph<N,E> graph, N entry)
func (this_ *FixedPointGraphTraversal[N, E]) ComputeFixedPoint_2_DiGraph_any_any__any(graph DiGraph[N, E], entry N) {
	_ = graph
	_ = entry
	var entryNode DiGraph_DiGraphNode[N, E] = j2gort.AsIface[DiGraph_DiGraphNode[N, E]](graph.GetNode(entry))
	_ = entryNode
	{
		var iifeCond_1 bool
		{

			v_, ok_ := graph.(interface {
				Implements_graph_LinkedDirectedGraph()
				IsJavaNil_() bool
			})
			iifeCond_1 = ok_ && !v_.IsJavaNil_()
		}
		if (iifeCond_1 && (entryNode != nil && !entryNode.IsJavaNil_())) && !(entryNode.IsInDfaQueue()) {
			var workQueue j2gort.Deque[DiGraph_DiGraphNode[N, E]]
			_ = workQueue
			entryNode.SetInDfaQueue(true)
			workQueue.AddLast(entryNode)
			this_.priv_FixedPointGraphTraversal_runWorkQueueFixedPoint_6287(graph, &workQueue)
			return
		}
	}
	var workSet *j2gort.Ext_java_util_LinkedHashSet = j2gort.NewExt_java_util_LinkedHashSet_1_int32(int32(4))
	_ = workSet
	workSet.Add_1(entryNode)
	this_.priv_FixedPointGraphTraversal_runWorkSetFixedPoint_6288(graph, workSet)
}

// Java method: FixedPointGraphTraversal.computeFixedPoint(com.google.javascript.jscomp.graph.DiGraph<N,E> graph, java.util.Set<N> entrySet)
func (this_ *FixedPointGraphTraversal[N, E]) ComputeFixedPoint_2_DiGraph_any_any__Ext_java_util_Set(graph DiGraph[N, E], entrySet j2gort.Ext_java_util_Set) {
	_ = graph
	_ = entrySet
	{
		var iifeCond_1 bool
		{

			v_, ok_ := graph.(interface {
				Implements_graph_LinkedDirectedGraph()
				IsJavaNil_() bool
			})
			iifeCond_1 = ok_ && !v_.IsJavaNil_()
		}
		if iifeCond_1 {
			var workQueue j2gort.Deque[DiGraph_DiGraphNode[N, E]] = j2gort.NewDequeCap[DiGraph_DiGraphNode[N, E]](max(entrySet.Size(), int32(8)))
			_ = workQueue
			var ok bool = true
			_ = ok
			{
				cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(entrySet))
				for cur_3.HasNext() {
					rawElem_4 := cur_3.Next()
					n := func() N {
						tv_5 := rawElem_4
						if tv_5 == nil {
							return func() N {
								if _, ok_ := any((*N)(nil)).(*string); ok_ {
									return any(j2gort.NullString).(N)
								}
								var z_ N
								return z_
							}()
						}
						return j2gort.Cast[N](tv_5)
					}()
					_ = n
					var node DiGraph_DiGraphNode[N, E] = j2gort.AsIface[DiGraph_DiGraphNode[N, E]](graph.GetNode(n))
					_ = node
					if node == nil || node.IsJavaNil_() {
						ok = false
						break
					}
					if !(node.IsInDfaQueue()) {
						node.SetInDfaQueue(true)
						workQueue.AddLast(node)
					}
				}
			}
			if ok {
				this_.priv_FixedPointGraphTraversal_runWorkQueueFixedPoint_6287(graph, &workQueue)
				return
			}
			for !(workQueue.IsEmpty()) {
				workQueue.RemoveFirst().SetInDfaQueue(false)
			}
		}
	}
	var workSet *j2gort.Ext_java_util_LinkedHashSet = j2gort.NewExt_java_util_LinkedHashSet_1_int32(entrySet.Size())
	_ = workSet
	{
		cur_6 := j2gort.NewIterCursor(j2gort.CheckNotNull(entrySet))
		for cur_6.HasNext() {
			rawElem_7 := cur_6.Next()
			n_2 := func() N {
				tv_8 := rawElem_7
				if tv_8 == nil {
					return func() N {
						if _, ok_ := any((*N)(nil)).(*string); ok_ {
							return any(j2gort.NullString).(N)
						}
						var z_ N
						return z_
					}()
				}
				return j2gort.Cast[N](tv_8)
			}()
			_ = n_2
			workSet.Add_1(j2gort.AsIface[DiGraph_DiGraphNode[N, E]](graph.GetNode(n_2)))
		}
	}
	this_.priv_FixedPointGraphTraversal_runWorkSetFixedPoint_6288(graph, workSet)
}

// Java method: FixedPointGraphTraversal.runWorkQueueFixedPoint(com.google.javascript.jscomp.graph.DiGraph<N,E> graph, java.util.ArrayDeque<com.google.javascript.jscomp.graph.DiGraph.DiGraphNode<N,E>> workQueue)
func (this_ *FixedPointGraphTraversal[N, E]) priv_FixedPointGraphTraversal_runWorkQueueFixedPoint_6287(graph DiGraph[N, E], workQueue *j2gort.Deque[DiGraph_DiGraphNode[N, E]]) {
	_ = graph
	_ = workQueue
	var cycleCount int64 = int64(int32(0))
	_ = cycleCount
	var nodeCount int64 = min(int64(graph.GetNodeCount()), StaticGet_FixedPointGraphTraversal_MAX_NODE_COUNT_FOR_ITERATION_LIMIT())
	_ = nodeCount
	var maxIterations int64 = max(((nodeCount * nodeCount) * nodeCount), int64(100))
	_ = maxIterations
	defer func() {
		for !(workQueue.IsEmpty()) {
			workQueue.RemoveFirst().SetInDfaQueue(false)
		}
	}()
	{
		for ; !(workQueue.IsEmpty()) && (cycleCount < maxIterations); cycleCount = cycleCount + 1 {
			var node DiGraph_DiGraphNode[N, E] = workQueue.RemoveFirst()
			_ = node
			node.SetInDfaQueue(false)
			this_.priv_FixedPointGraphTraversal_visitNodeQueue_6289(node, workQueue)
		}
		if !(cycleCount != maxIterations) {
			j2gort.CheckStateStr(false, "Fixed point computation not halting")
		}
	}
}

// Java method: FixedPointGraphTraversal.runWorkSetFixedPoint(com.google.javascript.jscomp.graph.DiGraph<N,E> graph, java.util.LinkedHashSet<com.google.javascript.jscomp.graph.DiGraph.DiGraphNode<N,E>> workSet)
func (this_ *FixedPointGraphTraversal[N, E]) priv_FixedPointGraphTraversal_runWorkSetFixedPoint_6288(graph DiGraph[N, E], workSet *j2gort.Ext_java_util_LinkedHashSet) {
	_ = graph
	_ = workSet
	var cycleCount int64 = int64(int32(0))
	_ = cycleCount
	var nodeCount int64 = min(int64(graph.GetNodeCount()), StaticGet_FixedPointGraphTraversal_MAX_NODE_COUNT_FOR_ITERATION_LIMIT())
	_ = nodeCount
	var maxIterations int64 = max(((nodeCount * nodeCount) * nodeCount), int64(100))
	_ = maxIterations
	for ; !(workSet.IsEmpty()) && (cycleCount < maxIterations); cycleCount = cycleCount + 1 {
		this_.priv_FixedPointGraphTraversal_visitNode_6290(j2gort.AsIface[DiGraph_DiGraphNode[N, E]](workSet.GetFirst()), workSet)
	}
	if !(cycleCount != maxIterations) {
		j2gort.CheckStateStr(false, "Fixed point computation not halting")
	}
}

// Java method: FixedPointGraphTraversal.visitNodeQueue(com.google.javascript.jscomp.graph.DiGraph.DiGraphNode<N,E> node, java.util.ArrayDeque<com.google.javascript.jscomp.graph.DiGraph.DiGraphNode<N,E>> workQueue)
func (this_ *FixedPointGraphTraversal[N, E]) priv_FixedPointGraphTraversal_visitNodeQueue_6289(node DiGraph_DiGraphNode[N, E], workQueue *j2gort.Deque[DiGraph_DiGraphNode[N, E]]) {
	_ = node
	_ = workQueue
	switch (this_.traversalDirection).Ordinal() {
	case int32(1) /* TraversalDirection.OUTWARDS */ :
		var sourceValue N = node.GetValue()
		_ = sourceValue
		{
			var i int32 = int32(0)
			_ = i
			var outCount int32 = node.GetOutEdgeCount()
			_ = outCount
			for ; i < outCount; i = i + 1 {
				var edge DiGraph_DiGraphEdge[N, E] = node.GetOutEdgeAt(i)
				_ = edge
				var dest DiGraph_DiGraphNode[N, E] = edge.GetDestination()
				_ = dest
				if this_.callback.TraverseEdge(sourceValue, edge.GetValue(), dest.GetValue()) {
					if !(dest.IsInDfaQueue()) {
						dest.SetInDfaQueue(true)
						workQueue.AddLast(dest)
					}
				}
			}
		}
		return
	case int32(0) /* TraversalDirection.INWARDS */ :
		var revSourceValue N = node.GetValue()
		_ = revSourceValue
		{
			var i_2 int32 = int32(0)
			_ = i_2
			var inCount int32 = node.GetInEdgeCount()
			_ = inCount
			for ; i_2 < inCount; i_2 = i_2 + 1 {
				var edge_2 DiGraph_DiGraphEdge[N, E] = node.GetInEdgeAt(i_2)
				_ = edge_2
				var src DiGraph_DiGraphNode[N, E] = edge_2.GetSource()
				_ = src
				if this_.callback.TraverseEdge(revSourceValue, edge_2.GetValue(), src.GetValue()) {
					if !(src.IsInDfaQueue()) {
						src.SetInDfaQueue(true)
						workQueue.AddLast(src)
					}
				}
			}
		}
		return
	}
	panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError", "Unrecognized direction "+j2gort.StrVal(this_.traversalDirection))))
}

// Java method: FixedPointGraphTraversal.visitNode(com.google.javascript.jscomp.graph.DiGraph.DiGraphNode<N,E> node, java.util.LinkedHashSet<com.google.javascript.jscomp.graph.DiGraph.DiGraphNode<N,E>> workSet)
func (this_ *FixedPointGraphTraversal[N, E]) priv_FixedPointGraphTraversal_visitNode_6290(node DiGraph_DiGraphNode[N, E], workSet *j2gort.Ext_java_util_LinkedHashSet) {
	_ = node
	_ = workSet
	workSet.Remove_1_any(node)
	switch (this_.traversalDirection).Ordinal() {
	case int32(1) /* TraversalDirection.OUTWARDS */ :
		var sourceValue N = node.GetValue()
		_ = sourceValue
		{
			var i int32 = int32(0)
			_ = i
			var outCount int32 = node.GetOutEdgeCount()
			_ = outCount
			for ; i < outCount; i = i + 1 {
				var edge DiGraph_DiGraphEdge[N, E] = node.GetOutEdgeAt(i)
				_ = edge
				var dest DiGraph_DiGraphNode[N, E] = edge.GetDestination()
				_ = dest
				if this_.callback.TraverseEdge(sourceValue, edge.GetValue(), dest.GetValue()) {
					workSet.Add_1(dest)
				}
			}
		}
		return
	case int32(0) /* TraversalDirection.INWARDS */ :
		var revSourceValue N = node.GetValue()
		_ = revSourceValue
		{
			var i_2 int32 = int32(0)
			_ = i_2
			var inCount int32 = node.GetInEdgeCount()
			_ = inCount
			for ; i_2 < inCount; i_2 = i_2 + 1 {
				var edge_2 DiGraph_DiGraphEdge[N, E] = node.GetInEdgeAt(i_2)
				_ = edge_2
				var src DiGraph_DiGraphNode[N, E] = edge_2.GetSource()
				_ = src
				if this_.callback.TraverseEdge(revSourceValue, edge_2.GetValue(), src.GetValue()) {
					workSet.Add_1(src)
				}
			}
		}
		return
	}
	panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError", "Unrecognized direction "+j2gort.StrVal(this_.traversalDirection))))
}

// Java enum: com.google.javascript.jscomp.graph.FixedPointGraphTraversal.TraversalDirection
type FixedPointGraphTraversal_TraversalDirection uint8

const (
	enum_FixedPointGraphTraversal_TraversalDirection_INWARDS  FixedPointGraphTraversal_TraversalDirection = 1
	enum_FixedPointGraphTraversal_TraversalDirection_OUTWARDS FixedPointGraphTraversal_TraversalDirection = 2
)

var enumBaseTable_FixedPointGraphTraversal_TraversalDirection = [...]j2gort.EnumBase{
	{Name_: "INWARDS", Ordinal_: 0},
	{Name_: "OUTWARDS", Ordinal_: 1},
}
var enumValues_FixedPointGraphTraversal_TraversalDirection = []FixedPointGraphTraversal_TraversalDirection{1, 2}

func StaticGet_FixedPointGraphTraversal_TraversalDirection_INWARDS() FixedPointGraphTraversal_TraversalDirection {
	return enum_FixedPointGraphTraversal_TraversalDirection_INWARDS
}
func StaticGet_FixedPointGraphTraversal_TraversalDirection_OUTWARDS() FixedPointGraphTraversal_TraversalDirection {
	return enum_FixedPointGraphTraversal_TraversalDirection_OUTWARDS
}

func FixedPointGraphTraversal_TraversalDirection_Values() []FixedPointGraphTraversal_TraversalDirection {
	return j2gort.CloneSlice(enumValues_FixedPointGraphTraversal_TraversalDirection)
}

func FixedPointGraphTraversal_TraversalDirection_ValueOf(name string) FixedPointGraphTraversal_TraversalDirection {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_FixedPointGraphTraversal_TraversalDirection {
		if enumBaseTable_FixedPointGraphTraversal_TraversalDirection[i_].Name_ == name {
			return FixedPointGraphTraversal_TraversalDirection(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.graph.FixedPointGraphTraversal.TraversalDirection." + name))
}

func (this_ FixedPointGraphTraversal_TraversalDirection) IsJavaNil_() bool                { return this_ == 0 }
func (FixedPointGraphTraversal_TraversalDirection) HasIdentityEquals_()                   {}
func (FixedPointGraphTraversal_TraversalDirection) Implements_j2gort_Ext_java_lang_Enum() {}
func (FixedPointGraphTraversal_TraversalDirection) Implements_j2gort_Ext_java_lang_constant_Constable() {
}
func (this_ FixedPointGraphTraversal_TraversalDirection) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (FixedPointGraphTraversal_TraversalDirection) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.graph.FixedPointGraphTraversal$TraversalDirection", "TraversalDirection"
}
func (this_ FixedPointGraphTraversal_TraversalDirection) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_FixedPointGraphTraversal_TraversalDirection[this_-1]
}
func (this_ FixedPointGraphTraversal_TraversalDirection) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ FixedPointGraphTraversal_TraversalDirection) EnumOrdinal_() int32 {
	return int32(this_) - 1
}
func (this_ FixedPointGraphTraversal_TraversalDirection) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_FixedPointGraphTraversal_TraversalDirection[this_-1].Name_
}
func (this_ FixedPointGraphTraversal_TraversalDirection) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_FixedPointGraphTraversal_TraversalDirection[this_-1].Name_
}
func (this_ FixedPointGraphTraversal_TraversalDirection) Equals(other_ any) bool {
	o_, ok_ := other_.(FixedPointGraphTraversal_TraversalDirection)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ FixedPointGraphTraversal_TraversalDirection) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_FixedPointGraphTraversal_TraversalDirection[this_-1].HashCode()
}
func (this_ FixedPointGraphTraversal_TraversalDirection) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_FixedPointGraphTraversal_TraversalDirection[this_-1].CompareTo(other_)
}
func (this_ FixedPointGraphTraversal_TraversalDirection) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_FixedPointGraphTraversal_TraversalDirection[this_-1].Name_)
}
func init() {
	vals_ := []any{FixedPointGraphTraversal_TraversalDirection(1), FixedPointGraphTraversal_TraversalDirection(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.graph.FixedPointGraphTraversal.TraversalDirection", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.graph.FixedPointGraphTraversal$TraversalDirection", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.graph.FixedPointGraphTraversal.TraversalDirection", "com.google.javascript.jscomp.graph.FixedPointGraphTraversal$TraversalDirection", func(v_ any) bool {
		return j2gort.InstanceOf[FixedPointGraphTraversal_TraversalDirection](v_)
	})
}

// Java interface: com.google.javascript.jscomp.graph.FixedPointGraphTraversal.EdgeCallback
type FixedPointGraphTraversal_EdgeCallback[NodeT comparable, EdgeT any] interface {
	Implements_graph_FixedPointGraphTraversal_EdgeCallback()
	IsJavaNil_() bool
	TraverseEdge(NodeT, EdgeT, NodeT) bool
}

func ToIface_FixedPointGraphTraversal_EdgeCallback[NodeT comparable, EdgeT any, T_iface_ interface {
	FixedPointGraphTraversal_EdgeCallback[NodeT, EdgeT]
	comparable
}](v T_iface_) FixedPointGraphTraversal_EdgeCallback[NodeT, EdgeT] {
	var z T_iface_
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.graph.FixedPointGraphTraversal.EdgeCallback", "com.google.javascript.jscomp.graph.FixedPointGraphTraversal$EdgeCallback", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			Implements_graph_FixedPointGraphTraversal_EdgeCallback()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}

type lambda_FixedPointGraphTraversal_EdgeCallback[NodeT comparable, EdgeT any] struct {
	fn func(NodeT, EdgeT, NodeT) bool
}

func (l *lambda_FixedPointGraphTraversal_EdgeCallback[NodeT, EdgeT]) IsJavaNil_() bool {
	return l == nil
}
func NewLambda_FixedPointGraphTraversal_EdgeCallback[NodeT comparable, EdgeT any](fn func(NodeT, EdgeT, NodeT) bool) FixedPointGraphTraversal_EdgeCallback[NodeT, EdgeT] {
	if fn == nil {
		return nil
	}
	return &lambda_FixedPointGraphTraversal_EdgeCallback[NodeT, EdgeT]{fn: fn}
}

func (*lambda_FixedPointGraphTraversal_EdgeCallback[NodeT, EdgeT]) Implements_graph_FixedPointGraphTraversal_EdgeCallback() {
}
func (l *lambda_FixedPointGraphTraversal_EdgeCallback[NodeT, EdgeT]) TraverseEdge(p0_ NodeT, p1_ EdgeT, p2_ NodeT) bool {
	return l.fn(p0_, p1_, p2_)
}
