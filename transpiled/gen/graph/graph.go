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
// Port of com/google/javascript/jscomp/graph/Graph.java (generated mechanically by j2go).

package graph

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java polymorphic class interface: com.google.javascript.jscomp.graph.Graph
type Graph[N comparable, E any] interface {
	Implements_graph_Graph()
	IsJavaNil_() bool
	GetGraphBase() *GraphBase[N, E]
	AdjacencyGraph[N, E]
	Connect_3_any_any_any(N, E, N)
	Disconnect(N, N)
	ConnectIfNotFound(N, E, N)
	CreateNode(N) any
	GetEdges_0() j2gort.Ext_java_util_List
	GetEdges_2(N, N) j2gort.Ext_java_util_List
	GetNodeDegree(N) int32
	GetNeighborNodes(N) j2gort.Ext_java_util_List
	GetFirstEdge(N, N) Graph_GraphEdge[N, E]
	HasNode(N) bool
	IsConnected_2(N, N) bool
	IsConnected_3_any_any_any(N, E, N) bool
	GetNodeOrFail(N) GraphNode[N, E]
	ClearEdgeAnnotations()
	PushNodeAnnotations()
	PopNodeAnnotations()
	PushEdgeAnnotations()
	PopEdgeAnnotations()
}

func ToIface_Graph[N comparable, E any, T_iface_ interface {
	Graph[N, E]
	comparable
}](v T_iface_) Graph[N, E] {
	var z T_iface_
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.jscomp.graph.Graph
type GraphBase[N comparable, E any] struct {
	self_               Graph[N, E]
	nodeAnnotationStack j2gort.Ext_java_util_Deque `gson:"nodeAnnotationStack"`
	edgeAnnotationStack j2gort.Ext_java_util_Deque `gson:"edgeAnnotationStack"`
	idHash_             atomic.Int32
}

func (this_ *GraphBase[N, E]) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(Graph[N, E])
}
func (this_ *GraphBase[N, E]) IsJavaNil_() bool { return this_ == nil }
func (*GraphBase[N, E]) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.graph.Graph", "Graph"
}
func (this_ *GraphBase[N, E]) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*GraphBase[N, E]) IsGenericInst_graph_Graph_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.graph.Graph", "com.google.javascript.jscomp.graph.Graph", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			IsGenericInst_graph_Graph_()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}
func (b *GraphBase[N, E]) GetGraphBase() *GraphBase[N, E] { return b }
func (*GraphBase[N, E]) Implements_graph_Graph()          {}
func (*GraphBase[N, E]) Implements_graph_AdjacencyGraph() {}

func (this_ *GraphBase[N, E]) initFields_Graph() {
	_ = this_
}

func InitGraph[N comparable, E any](this_ *GraphBase[N, E]) {
	this_.initFields_Graph()
}

// Java method: Graph.connect(N n1, E edge, N n2)
func (this_ *GraphBase[N, E]) Connect_3_any_any_any(n1 N, edge E, n2 N) {
	_ = n1
	_ = edge
	_ = n2
	panic("abstract or native method invoked")
}

// Java method: Graph.disconnect(N n1, N n2)
func (this_ *GraphBase[N, E]) Disconnect(n1 N, n2 N) {
	_ = n1
	_ = n2
	panic("abstract or native method invoked")
}

// Java method: Graph.connectIfNotFound(N n1, E edge, N n2)
func (this_ *GraphBase[N, E]) ConnectIfNotFound(n1 N, edge E, n2 N) {
	_ = n1
	_ = edge
	_ = n2
	if !(this_.self_.IsConnected_3_any_any_any(n1, edge, n2)) {
		this_.self_.Connect_3_any_any_any(n1, edge, n2)
	}
}

// Java method: Graph.createNode(N value)
func (this_ *GraphBase[N, E]) CreateNode(value N) any {
	_ = value
	panic("abstract or native method invoked")
}

// Java method: Graph.getNodes()
func (this_ *GraphBase[N, E]) GetNodes() j2gort.Ext_java_util_Collection {
	panic("abstract or native method invoked")
}

// Java method: Graph.getNodeCount()
func (this_ *GraphBase[N, E]) GetNodeCount() int32 {
	panic("abstract or native method invoked")
}

// Java method: Graph.getEdges()
func (this_ *GraphBase[N, E]) GetEdges_0() j2gort.Ext_java_util_List {
	panic("abstract or native method invoked")
}

// Java method: Graph.getEdges(N n1, N n2)
func (this_ *GraphBase[N, E]) GetEdges_2(n1 N, n2 N) j2gort.Ext_java_util_List {
	_ = n1
	_ = n2
	panic("abstract or native method invoked")
}

// Java method: Graph.getNodeDegree(N value)
func (this_ *GraphBase[N, E]) GetNodeDegree(value N) int32 {
	_ = value
	panic("abstract or native method invoked")
}

// Java method: Graph.getWeight(N value)
func (this_ *GraphBase[N, E]) GetWeight(value N) int32 {
	_ = value
	return this_.self_.GetNodeDegree(value)
}

// Java method: Graph.getNeighborNodes(N value)
func (this_ *GraphBase[N, E]) GetNeighborNodes(value N) j2gort.Ext_java_util_List {
	_ = value
	panic("abstract or native method invoked")
}

// Java method: Graph.getFirstEdge(N n1, N n2)
func (this_ *GraphBase[N, E]) GetFirstEdge(n1 N, n2 N) Graph_GraphEdge[N, E] {
	_ = n1
	_ = n2
	panic("abstract or native method invoked")
}

// Java method: Graph.hasNode(N n)
func (this_ *GraphBase[N, E]) HasNode(n N) bool {
	_ = n
	return !(j2gort.IfaceEquals[GraphNode[N, E]](j2gort.AsIface[GraphNode[N, E]](this_.self_.GetNode(n)), nil))
}

// Java method: Graph.isConnected(N n1, N n2)
func (this_ *GraphBase[N, E]) IsConnected_2(n1 N, n2 N) bool {
	_ = n1
	_ = n2
	panic("abstract or native method invoked")
}

// Java method: Graph.isConnected(N n1, E e, N n2)
func (this_ *GraphBase[N, E]) IsConnected_3_any_any_any(n1 N, e E, n2 N) bool {
	_ = n1
	_ = e
	_ = n2
	panic("abstract or native method invoked")
}

// Java method: Graph.getNodeOrFail(N val)
func (this_ *GraphBase[N, E]) GetNodeOrFail(val N) GraphNode[N, E] {
	_ = val
	var node GraphNode[N, E] = j2gort.AsIface[GraphNode[N, E]](this_.self_.GetNode(val))
	_ = node
	if node == nil || node.IsJavaNil_() {
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", j2gort.StrVal(j2gort.ToAny(val))+" does not exist in graph")))
	}
	return node
}

// Java method: Graph.clearNodeAnnotations()
func (this_ *GraphBase[N, E]) ClearNodeAnnotations() {
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.self_.GetNodes()))
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			n := j2gort.AsIface[GraphNode[N, E]](rawElem_2)
			_ = n
			n.SetAnnotation(nil)
		}
	}
}

// Java method: Graph.clearEdgeAnnotations()
func (this_ *GraphBase[N, E]) ClearEdgeAnnotations() {
	{
		cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.self_.GetEdges_0()))
		for cur_3.HasNext() {
			rawElem_4 := cur_3.Next()
			e := j2gort.AsIface[Graph_GraphEdge[N, E]](rawElem_4)
			_ = e
			e.SetAnnotation(nil)
		}
	}
}

// Java method: Graph.pushNodeAnnotations()
func (this_ *GraphBase[N, E]) PushNodeAnnotations() {
	if this_.nodeAnnotationStack == nil {
		this_.nodeAnnotationStack = j2gort.NewExt_java_util_ArrayDeque()
	}
	Graph_PushAnnotations(this_.nodeAnnotationStack, this_.self_.GetNodes())
}

// Java method: Graph.popNodeAnnotations()
func (this_ *GraphBase[N, E]) PopNodeAnnotations() {
	j2gort.CheckNotNull(this_.nodeAnnotationStack)
	Graph_PopAnnotations(this_.nodeAnnotationStack)
}

// Java method: Graph.pushEdgeAnnotations()
func (this_ *GraphBase[N, E]) PushEdgeAnnotations() {
	if this_.edgeAnnotationStack == nil {
		this_.edgeAnnotationStack = j2gort.NewExt_java_util_ArrayDeque()
	}
	Graph_PushAnnotations(this_.edgeAnnotationStack, j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.self_.GetEdges_0()))
}

// Java method: Graph.popEdgeAnnotations()
func (this_ *GraphBase[N, E]) PopEdgeAnnotations() {
	j2gort.CheckNotNull(this_.edgeAnnotationStack)
	Graph_PopAnnotations(this_.edgeAnnotationStack)
}

// Java static method: Graph.pushAnnotations(java.util.Deque<com.google.javascript.jscomp.graph.Graph.GraphAnnotationState> stack, java.util.Collection<? extends com.google.javascript.jscomp.graph.Annotatable> haveAnnotations)
func Graph_PushAnnotations(stack j2gort.Ext_java_util_Deque, haveAnnotations j2gort.Ext_java_util_Collection) {
	_ = stack
	_ = haveAnnotations
	stack.Push(NewGraph_GraphAnnotationState(haveAnnotations.Size()))
	{
		cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(haveAnnotations))
		for cur_5.HasNext() {
			rawElem_6 := cur_5.Next()
			h := j2gort.AsIface[Annotatable](rawElem_6)
			_ = h
			j2gort.CastPtr[Graph_GraphAnnotationState](stack.Peek()).Add_1(NewGraph_AnnotationState(h, h.GetAnnotation()))
			h.SetAnnotation(nil)
		}
	}
}

// Java static method: Graph.popAnnotations(java.util.Deque<com.google.javascript.jscomp.graph.Graph.GraphAnnotationState> stack)
func Graph_PopAnnotations(stack j2gort.Ext_java_util_Deque) {
	_ = stack
	{
		cur_7 := j2gort.NewIterCursor(j2gort.CheckNotNullPtr(j2gort.CastPtr[Graph_GraphAnnotationState](stack.Pop())))
		for cur_7.HasNext() {
			rawElem_8 := cur_7.Next()
			as := j2gort.CastPtr[Graph_AnnotationState](rawElem_8)
			_ = as
			as.first.SetAnnotation(as.second)
		}
	}
}

// Java class: com.google.javascript.jscomp.graph.Graph.AnnotationState
type Graph_AnnotationState struct {
	first   Annotatable `gson:"first"`
	second  Annotation  `gson:"second"`
	idHash_ atomic.Int32
}

func (this_ *Graph_AnnotationState) BindSelf_(s any) {
	_ = s
}
func (this_ *Graph_AnnotationState) IsJavaNil_() bool { return this_ == nil }
func (*Graph_AnnotationState) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.graph.Graph$AnnotationState", "AnnotationState"
}
func (this_ *Graph_AnnotationState) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Graph_AnnotationState) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.graph.Graph.AnnotationState", "com.google.javascript.jscomp.graph.Graph$AnnotationState", func(v_ any) bool {
		return j2gort.InstanceOf[*Graph_AnnotationState](v_)
	})
}

func (this_ *Graph_AnnotationState) initFields_Graph_AnnotationState() {
	_ = this_
}

// Java constructor: AnnotationState(com.google.javascript.jscomp.graph.Annotatable annotatable, com.google.javascript.jscomp.graph.Annotation annotation)
func NewGraph_AnnotationState(annotatable Annotatable, annotation Annotation) *Graph_AnnotationState {
	this_ := &Graph_AnnotationState{}
	InitGraph_AnnotationState(this_, annotatable, annotation)
	return this_
}

func InitGraph_AnnotationState(this_ *Graph_AnnotationState, annotatable Annotatable, annotation Annotation) {
	_ = annotatable
	_ = annotation
	this_.initFields_Graph_AnnotationState()
	this_.first = annotatable
	this_.second = annotation
}

// Java class: com.google.javascript.jscomp.graph.Graph.GraphAnnotationState
type Graph_GraphAnnotationState struct {
	j2gort.Ext_java_util_ArrayListBase
}

func (this_ *Graph_GraphAnnotationState) BindSelf_(s any) {
	_ = s
	this_.Ext_java_util_ArrayListBase.BindExtSelf_(s)
}
func (this_ *Graph_GraphAnnotationState) IsJavaNil_() bool { return this_ == nil }
func (*Graph_GraphAnnotationState) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.graph.Graph$GraphAnnotationState", "GraphAnnotationState"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.graph.Graph.GraphAnnotationState", "com.google.javascript.jscomp.graph.Graph$GraphAnnotationState", func(v_ any) bool {
		return j2gort.InstanceOf[*Graph_GraphAnnotationState](v_)
	})
}
func (*Graph_GraphAnnotationState) Implements_j2gort_Ext_java_util_ArrayList()           {}
func (*Graph_GraphAnnotationState) Implements_j2gort_Ext_java_util_AbstractList()        {}
func (*Graph_GraphAnnotationState) Implements_j2gort_Ext_java_util_List()                {}
func (*Graph_GraphAnnotationState) Implements_j2gort_Ext_java_util_RandomAccess()        {}
func (*Graph_GraphAnnotationState) Implements_j2gort_Ext_java_util_AbstractCollection()  {}
func (*Graph_GraphAnnotationState) Implements_j2gort_Ext_java_util_SequencedCollection() {}
func (*Graph_GraphAnnotationState) Implements_j2gort_Ext_java_util_Collection()          {}
func (*Graph_GraphAnnotationState) Implements_j2gort_Ext_java_lang_Iterable()            {}

func StaticGet_Graph_GraphAnnotationState_SerialVersionUID() int64 { return int64(1) }
func (this_ *Graph_GraphAnnotationState) initFields_Graph_GraphAnnotationState() {
	_ = this_
}

// Java constructor: GraphAnnotationState(int size)
func NewGraph_GraphAnnotationState(size int32) *Graph_GraphAnnotationState {
	this_ := &Graph_GraphAnnotationState{}
	this_.BindSelf_(this_)
	InitGraph_GraphAnnotationState(this_, size)
	return this_
}

func InitGraph_GraphAnnotationState(this_ *Graph_GraphAnnotationState, size int32) {
	_ = size
	j2gort.InitExt_java_util_ArrayList_1_int32(&this_.Ext_java_util_ArrayListBase, size)
	this_.initFields_Graph_GraphAnnotationState()
}

// Java interface: com.google.javascript.jscomp.graph.Graph.GraphEdge
type Graph_GraphEdge[N comparable, E any] interface {
	Implements_graph_Graph_GraphEdge()
	IsJavaNil_() bool
	Annotatable
	GetValue() E
	GetNodeA() GraphNode[N, E]
	GetNodeB() GraphNode[N, E]
}

func ToIface_Graph_GraphEdge[N comparable, E any, T_iface_ interface {
	Graph_GraphEdge[N, E]
	comparable
}](v T_iface_) Graph_GraphEdge[N, E] {
	var z T_iface_
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.graph.Graph.GraphEdge", "com.google.javascript.jscomp.graph.Graph$GraphEdge", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			Implements_graph_Graph_GraphEdge()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}

// Java class: com.google.javascript.jscomp.graph.Graph.SimpleSubGraph
type Graph_SimpleSubGraph[N comparable, E any] struct {
	graph   Graph[N, E]       `gson:"graph"`
	nodes   []GraphNode[N, E] `gson:"nodes"`
	idHash_ atomic.Int32
}

func (this_ *Graph_SimpleSubGraph[N, E]) BindSelf_(s any) {
	_ = s
}
func (this_ *Graph_SimpleSubGraph[N, E]) IsJavaNil_() bool { return this_ == nil }
func (*Graph_SimpleSubGraph[N, E]) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.graph.Graph$SimpleSubGraph", "SimpleSubGraph"
}
func (this_ *Graph_SimpleSubGraph[N, E]) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Graph_SimpleSubGraph[N, E]) HasIdentityEquals_()                        {}
func (*Graph_SimpleSubGraph[N, E]) IsGenericInst_graph_Graph_SimpleSubGraph_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.graph.Graph.SimpleSubGraph", "com.google.javascript.jscomp.graph.Graph$SimpleSubGraph", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			IsGenericInst_graph_Graph_SimpleSubGraph_()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}
func (*Graph_SimpleSubGraph[N, E]) Implements_graph_SubGraph() {}

func (this_ *Graph_SimpleSubGraph[N, E]) initFields_Graph_SimpleSubGraph() {
	_ = this_
	this_.nodes = []GraphNode[N, E]{}
}

// Java constructor: SimpleSubGraph(com.google.javascript.jscomp.graph.Graph<N,E> graph)
func NewGraph_SimpleSubGraph[N comparable, E any](graph Graph[N, E]) *Graph_SimpleSubGraph[N, E] {
	this_ := &Graph_SimpleSubGraph[N, E]{}
	InitGraph_SimpleSubGraph[N, E](this_, graph)
	return this_
}

func InitGraph_SimpleSubGraph[N comparable, E any](this_ *Graph_SimpleSubGraph[N, E], graph Graph[N, E]) {
	_ = graph
	this_.initFields_Graph_SimpleSubGraph()
	this_.graph = graph
}

// Java method: SimpleSubGraph.isIndependentOf(N value)
func (this_ *Graph_SimpleSubGraph[N, E]) IsIndependentOf(value N) bool {
	_ = value
	var node GraphNode[N, E] = j2gort.AsIface[GraphNode[N, E]](this_.graph.GetNode(value))
	_ = node
	for _, n := range this_.nodes {
		_ = n
		if this_.graph.GetNeighborNodes(n.GetValue()).Contains(node) {
			return false
		}
	}

	return true
}

// Java method: SimpleSubGraph.addNode(N value)
func (this_ *Graph_SimpleSubGraph[N, E]) AddNode(value N) {
	_ = value
	this_.nodes = append(this_.nodes, this_.graph.GetNodeOrFail(value))
}
