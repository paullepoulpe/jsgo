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
// Port of com/google/javascript/jscomp/graph/LinkedUndirectedGraph.java (generated mechanically by j2go).

package graph

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.graph.LinkedUndirectedGraph
type LinkedUndirectedGraph[N comparable, E any] struct {
	UndiGraphBase[N, E]
	nodes                                                j2gort.OrderedMap[N, *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]] `gson:"nodes"`
	useNodeAnnotations                                   bool                                                                         `gson:"useNodeAnnotations"`
	useEdgeAnnotations                                   bool                                                                         `gson:"useEdgeAnnotations"`
	slab_LinkedUndirectedGraph_LinkedUndirectedGraphEdge j2gort.PooledSlab[LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E]]
	slab_LinkedUndirectedGraph_LinkedUndirectedGraphNode j2gort.PooledSlab[LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]]
}

func (this_ *LinkedUndirectedGraph[N, E]) BindSelf_(s any) {
	_ = s
	this_.UndiGraphBase.BindSelf_(s)
}
func (this_ *LinkedUndirectedGraph[N, E]) IsJavaNil_() bool { return this_ == nil }
func (*LinkedUndirectedGraph[N, E]) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.graph.LinkedUndirectedGraph", "LinkedUndirectedGraph"
}
func (this_ *LinkedUndirectedGraph[N, E]) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.UndiGraphBase.IdentityHash_()
}
func (*LinkedUndirectedGraph[N, E]) HasIdentityEquals_()                         {}
func (*LinkedUndirectedGraph[N, E]) IsGenericInst_graph_LinkedUndirectedGraph_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.graph.LinkedUndirectedGraph", "com.google.javascript.jscomp.graph.LinkedUndirectedGraph", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			IsGenericInst_graph_LinkedUndirectedGraph_()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}
func (*LinkedUndirectedGraph[N, E]) Implements_graph_UndiGraph()      {}
func (*LinkedUndirectedGraph[N, E]) Implements_graph_GraphvizGraph()  {}
func (*LinkedUndirectedGraph[N, E]) Implements_graph_Graph()          {}
func (*LinkedUndirectedGraph[N, E]) Implements_graph_AdjacencyGraph() {}

func (this_ *LinkedUndirectedGraph[N, E]) initFields_LinkedUndirectedGraph() {
	_ = this_
	this_.nodes = j2gort.AcquireOrderedMap[N, *LinkedUndirectedGraph_LinkedUndirectedGraphNode[

		// Java constructor: LinkedUndirectedGraph(boolean useNodeAnnotations, boolean useEdgeAnnotations)
		N, E]]()
}

func NewLinkedUndirectedGraph[N comparable, E any](useNodeAnnotations bool, useEdgeAnnotations bool) *LinkedUndirectedGraph[N, E] {
	this_ := &LinkedUndirectedGraph[N, E]{}
	this_.BindSelf_(this_)
	InitLinkedUndirectedGraph[N, E](this_, useNodeAnnotations, useEdgeAnnotations)
	return this_
}

func InitLinkedUndirectedGraph[N comparable, E any](this_ *LinkedUndirectedGraph[N, E], useNodeAnnotations bool, useEdgeAnnotations bool) {
	_ = useNodeAnnotations
	_ = useEdgeAnnotations
	InitUndiGraph[N, E](&this_.UndiGraphBase)
	this_.initFields_LinkedUndirectedGraph()
	this_.useNodeAnnotations = useNodeAnnotations
	this_.useEdgeAnnotations = useEdgeAnnotations
}

// Java method: LinkedUndirectedGraph.newSubGraph()
func (this_ *LinkedUndirectedGraph[N, E]) NewSubGraph() SubGraph[N, E] {
	return ToIface_SubGraph[N, E](NewGraph_SimpleSubGraph[N, E](ToIface_Graph[N, E](this_)))
}

// Java static method: LinkedUndirectedGraph.create()
func LinkedUndirectedGraph_Create[N comparable, E any]() *LinkedUndirectedGraph[N, E] {
	return NewLinkedUndirectedGraph[N, E](true, true)
}

// Java method: LinkedUndirectedGraph.connect(N srcValue, E edgeValue, N destValue)
func (this_ *LinkedUndirectedGraph[N, E]) Connect_3_any_any_any(srcValue N, edgeValue E, destValue N) {
	_ = srcValue
	_ = edgeValue
	_ = destValue
	var src *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E] = j2gort.CastPtr[LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]](this_.GetNodeOrFail(srcValue))
	_ = src
	var dest *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E] = j2gort.CastPtr[LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]](this_.GetNodeOrFail(destValue))
	_ = dest
	var edge *LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E] = allocPooled_LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E](&this_.slab_LinkedUndirectedGraph_LinkedUndirectedGraphEdge, ToIface_UndiGraph_UndiGraphNode[N, E](src), edgeValue, ToIface_UndiGraph_UndiGraphNode[N, E](dest), this_.useEdgeAnnotations)
	_ = edge
	src.GetNeighborEdges().Add_1(j2gort.PtrToAny(edge))
	dest.GetNeighborEdges().Add_1(j2gort.PtrToAny(edge))
}

// Java method: LinkedUndirectedGraph.disconnect(N srcValue, N destValue)
func (this_ *LinkedUndirectedGraph[N, E]) Disconnect(srcValue N, destValue N) {
	_ = srcValue
	_ = destValue
	var src *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E] = j2gort.CastPtr[LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]](this_.GetNodeOrFail(srcValue))
	_ = src
	var dest *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E] = j2gort.CastPtr[LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]](this_.GetNodeOrFail(destValue))
	_ = dest
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.GetUndirectedGraphEdges(srcValue, destValue)))
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			edge := j2gort.AsIface[UndiGraph_UndiGraphEdge[N, E]](rawElem_2)
			_ = edge
			src.GetNeighborEdges().Remove_1_any(edge)
			dest.GetNeighborEdges().Remove_1_any(edge)
		}
	}
}

// Java method: LinkedUndirectedGraph.createUndirectedGraphNode(N nodeValue)
func (this_ *LinkedUndirectedGraph[N, E]) CreateUndirectedGraphNode(nodeValue N) UndiGraph_UndiGraphNode[N, E] {
	_ = nodeValue
	var node *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E] = this_.nodes.Get(nodeValue)
	_ = node
	if node == nil {
		node = allocPooled_LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E](&this_.slab_LinkedUndirectedGraph_LinkedUndirectedGraphNode, nodeValue, this_.useNodeAnnotations)
		this_.nodes.Put(nodeValue, node)
	}
	return ToIface_UndiGraph_UndiGraphNode[N, E](node)
}

// Java method: LinkedUndirectedGraph.getNeighborNodes(N value)
func (this_ *LinkedUndirectedGraph[N, E]) GetNeighborNodes(value N) j2gort.Ext_java_util_List {
	_ = value
	var uNode UndiGraph_UndiGraphNode[N, E] = this_.GetUndirectedGraphNode(value)
	_ = uNode
	return (j2gort.CastPtr[LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]](uNode)).priv_LinkedUndirectedGraphNode_neighborList_6370()
}

// Java method: LinkedUndirectedGraph.getUndirectedGraphEdges(N n1, N n2)
func (this_ *LinkedUndirectedGraph[N, E]) GetUndirectedGraphEdges(n1 N, n2 N) j2gort.Ext_java_util_List {
	_ = n1
	_ = n2
	var dNode1 *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E] = this_.nodes.Get(n1)
	_ = dNode1
	if dNode1 == nil {
		return nil
	}
	var dNode2 *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E] = this_.nodes.Get(n2)
	_ = dNode2
	if dNode2 == nil {
		return nil
	}
	var edges []UndiGraph_UndiGraphEdge[N, E] = []UndiGraph_UndiGraphEdge[N, E]{}
	_ = edges
	{
		cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(dNode1.GetNeighborEdges()))
		for cur_3.HasNext() {
			rawElem_4 := cur_3.Next()
			outEdge := j2gort.AsIface[UndiGraph_UndiGraphEdge[N, E]](rawElem_4)
			_ = outEdge
			if j2gort.IfaceEquals[GraphNode[N, E]](outEdge.GetNodeA(), dNode2) || j2gort.IfaceEquals[GraphNode[N, E]](outEdge.GetNodeB(), dNode2) {
				edges = append(edges, outEdge)
			}
		}
	}
	return j2gort.WrapTypedArrayList(edges)
}

// Java method: LinkedUndirectedGraph.getUndirectedGraphNode(N nodeValue)
func (this_ *LinkedUndirectedGraph[N, E]) GetUndirectedGraphNode(nodeValue N) UndiGraph_UndiGraphNode[N, E] {
	_ = nodeValue
	var node *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E] = this_.nodes.Get(nodeValue)
	_ = node
	return ToIface_UndiGraph_UndiGraphNode[N, E](node)
}

// Java method: LinkedUndirectedGraph.getUndirectedGraphNodes()
func (this_ *LinkedUndirectedGraph[N, E]) GetUndirectedGraphNodes() j2gort.Ext_java_util_Collection {
	return j2gort.Ext_java_util_Collections_UnmodifiableCollection(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.nodes.LiveValues()))
}

// Java method: LinkedUndirectedGraph.createNode(N value)
func (this_ *LinkedUndirectedGraph[N, E]) CreateNode(value N) any {
	_ = value
	return this_.CreateUndirectedGraphNode(value)
}

// Java method: LinkedUndirectedGraph.getEdges(N n1, N n2)
func (this_ *LinkedUndirectedGraph[N, E]) GetEdges_2(n1 N, n2 N) j2gort.Ext_java_util_List {
	_ = n1
	_ = n2
	return j2gort.Ext_java_util_Collections_UnmodifiableList(this_.GetUndirectedGraphEdges(n1, n2))
}

// Java method: LinkedUndirectedGraph.getEdges()
func (this_ *LinkedUndirectedGraph[N, E]) GetEdges_0() j2gort.Ext_java_util_List {
	var result []UndiGraph_UndiGraphEdge[N, E] = []UndiGraph_UndiGraphEdge[N, E]{}
	_ = result
	{
		cur_5 := this_.nodes.IterValues()
		for cur_5.HasNext() {

			node := cur_5.
				Next()
			_ = node
			{
				cur_7 := j2gort.NewIterCursor(j2gort.CheckNotNull(node.GetNeighborEdges()))
				for cur_7.HasNext() {
					rawElem_8 := cur_7.Next()
					edge := j2gort.AsIface[UndiGraph_UndiGraphEdge[N, E]](rawElem_8)
					_ = edge
					if j2gort.IfaceEquals[GraphNode[N, E]](edge.GetNodeA(), node) {
						result = append(result, edge)
					}
				}
			}
		}
	}
	return j2gort.WrapTypedArrayList(result)
}

// Java method: LinkedUndirectedGraph.getFirstEdge(N n1, N n2)
func (this_ *LinkedUndirectedGraph[N, E]) GetFirstEdge(n1 N, n2 N) Graph_GraphEdge[N, E] {
	_ = n1
	_ = n2
	var dNode1 UndiGraph_UndiGraphNode[N, E] = j2gort.AsIface[UndiGraph_UndiGraphNode[N, E]](this_.GetNodeOrFail(n1))
	_ = dNode1
	var dNode2 UndiGraph_UndiGraphNode[N, E] = j2gort.AsIface[UndiGraph_UndiGraphNode[N, E]](this_.GetNodeOrFail(n2))
	_ = dNode2
	{
		cur_9 := j2gort.NewIterCursor(j2gort.CheckNotNull(dNode1.GetNeighborEdges()))
		for cur_9.HasNext() {
			rawElem_10 := cur_9.Next()
			outEdge := j2gort.AsIface[UndiGraph_UndiGraphEdge[N, E]](rawElem_10)
			_ = outEdge
			if j2gort.IfaceEquals[GraphNode[N, E]](outEdge.GetNodeA(), dNode2) || j2gort.IfaceEquals[GraphNode[N, E]](outEdge.GetNodeB(), dNode2) {
				return j2gort.AsIface[Graph_GraphEdge[N, E]](outEdge)
			}
		}
	}
	return nil
}

// Java method: LinkedUndirectedGraph.getNode(N value)
func (this_ *LinkedUndirectedGraph[N, E]) GetNode(value N) any {
	_ = value
	return this_.GetUndirectedGraphNode(value)
}

// Java method: LinkedUndirectedGraph.isConnected(N n1, N n2)
func (this_ *LinkedUndirectedGraph[N, E]) IsConnected_2(n1 N, n2 N) bool {
	_ = n1
	_ = n2
	return this_.priv_LinkedUndirectedGraph_isConnected_6367(n1, j2gort.Ext_com_google_common_base_Predicates_AlwaysTrue(), n2)
}

// Java method: LinkedUndirectedGraph.isConnected(N n1, E e, N n2)
func (this_ *LinkedUndirectedGraph[N, E]) IsConnected_3_any_any_any(n1 N, e E, n2 N) bool {
	_ = n1
	_ = e
	_ = n2
	return this_.priv_LinkedUndirectedGraph_isConnected_6367(n1, j2gort.Ext_com_google_common_base_Predicates_EqualTo(j2gort.ToAny(e)), n2)
}

// Java method: LinkedUndirectedGraph.isConnected(N n1, com.google.common.base.Predicate<E> edgePredicate, N n2)
func (this_ *LinkedUndirectedGraph[N, E]) priv_LinkedUndirectedGraph_isConnected_6367(n1 N, edgePredicate j2gort.Ext_com_google_common_base_Predicate, n2 N) bool {
	_ = n1
	_ = edgePredicate
	_ = n2
	var dNode1 *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E] = this_.nodes.Get(n1)
	_ = dNode1
	if dNode1 == nil {
		return false
	}
	var dNode2 *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E] = this_.nodes.Get(n2)
	_ = dNode2
	if dNode2 == nil {
		return false
	}
	{
		cur_11 := j2gort.NewIterCursor(j2gort.CheckNotNull(dNode1.GetNeighborEdges()))
		for cur_11.HasNext() {
			rawElem_12 := cur_11.Next()
			outEdge := j2gort.AsIface[UndiGraph_UndiGraphEdge[N, E]](rawElem_12)
			_ = outEdge
			if (j2gort.IfaceEquals[GraphNode[N, E]](outEdge.GetNodeA(), dNode1) && j2gort.IfaceEquals[GraphNode[N, E]](outEdge.GetNodeB(), dNode2)) || (j2gort.IfaceEquals[GraphNode[N, E]](outEdge.GetNodeA(), dNode2) && j2gort.IfaceEquals[GraphNode[N, E]](outEdge.GetNodeB(), dNode1)) {
				if edgePredicate.Apply(j2gort.ToAny(outEdge.GetValue())) {
					return true
				}
			}
		}
	}
	return false
}

// Java method: LinkedUndirectedGraph.getGraphvizEdges()
func (this_ *LinkedUndirectedGraph[N, E]) GetGraphvizEdges() j2gort.Ext_java_util_List {
	var edgeList []GraphvizGraph_GraphvizEdge = []GraphvizGraph_GraphvizEdge{}
	_ = edgeList
	{
		cur_13 := this_.nodes.IterValues()
		for cur_13.HasNext() {

			node := cur_13.
				Next()
			_ = node
			{
				cur_15 := j2gort.NewIterCursor(j2gort.CheckNotNull(node.GetNeighborEdges()))
				for cur_15.HasNext() {
					rawElem_16 := cur_15.Next()
					edge := j2gort.AsIface[UndiGraph_UndiGraphEdge[N, E]](rawElem_16)
					_ = edge
					if j2gort.IfaceEquals[GraphNode[N, E]](edge.GetNodeA(), node) {
						edgeList = append(edgeList, j2gort.AsIface[GraphvizGraph_GraphvizEdge](edge))
					}
				}
			}
		}
	}
	return j2gort.WrapTypedArrayList(edgeList)
}

// Java method: LinkedUndirectedGraph.getName()
func (this_ *LinkedUndirectedGraph[N, E]) GetName() string {
	return "LinkedUndirectedGraph"
}

// Java method: LinkedUndirectedGraph.getGraphvizNodes()
func (this_ *LinkedUndirectedGraph[N, E]) GetGraphvizNodes() j2gort.Ext_java_util_List {
	var nodeList []*LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E] = j2gort.MakeNonNilSlice[*LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]](this_.nodes.Size())
	_ = nodeList
	{
		cur_17 := this_.nodes.IterValues()
		for cur_17.HasNext() {

			node := cur_17.
				Next()
			_ = node
			nodeList = append(nodeList, node)
		}
	}
	return j2gort.WrapTypedArrayList(nodeList)
}

// Java method: LinkedUndirectedGraph.isDirected()
func (this_ *LinkedUndirectedGraph[N, E]) IsDirected() bool {
	return false
}

// Java method: LinkedUndirectedGraph.getNodes()
func (this_ *LinkedUndirectedGraph[N, E]) GetNodes() j2gort.Ext_java_util_Collection {
	return j2gort.Ext_java_util_Collections_UnmodifiableCollection(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.nodes.LiveValues()))
}

// Java method: LinkedUndirectedGraph.getNodeCount()
func (this_ *LinkedUndirectedGraph[N, E]) GetNodeCount() int32 {
	return this_.nodes.Size()
}

// Java method: LinkedUndirectedGraph.getNodeDegree(N value)
func (this_ *LinkedUndirectedGraph[N, E]) GetNodeDegree(value N) int32 {
	_ = value
	var uNode UndiGraph_UndiGraphNode[N, E] = this_.GetUndirectedGraphNode(value)
	_ = uNode
	if uNode == nil || uNode.IsJavaNil_() {
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", j2gort.StrVal(j2gort.ToAny(value))+" not found in graph")))
	}
	return uNode.GetNeighborEdges().Size()
}

// Java class: com.google.javascript.jscomp.graph.LinkedUndirectedGraph.LinkedUndirectedGraphNode
type LinkedUndirectedGraph_LinkedUndirectedGraphNode[N comparable, E any] struct {
	neighborEdges j2gort.Ext_java_util_ArrayList `gson:"neighborEdges"`
	value         N                              `gson:"value"`
	annotation    Annotation                     `gson:"annotation"`
	idHash_       atomic.Int32
	useAnnotation bool `gson:"useAnnotation"`
}

func (this_ *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]) BindSelf_(s any) {
	_ = s
}
func (this_ *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]) IsJavaNil_() bool {
	return this_ == nil
}
func (*LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.graph.LinkedUndirectedGraph$LinkedUndirectedGraphNode", "LinkedUndirectedGraphNode"
}
func (this_ *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]) HasIdentityEquals_() {}
func (*LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]) IsGenericInst_graph_LinkedUndirectedGraph_LinkedUndirectedGraphNode_() {
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.graph.LinkedUndirectedGraph.LinkedUndirectedGraphNode", "com.google.javascript.jscomp.graph.LinkedUndirectedGraph$LinkedUndirectedGraphNode", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			IsGenericInst_graph_LinkedUndirectedGraph_LinkedUndirectedGraphNode_()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}
func (*LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]) Implements_graph_UndiGraph_UndiGraphNode() {
}
func (*LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]) Implements_graph_GraphvizGraph_GraphvizNode() {
}
func (*LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]) Implements_graph_GraphNode()   {}
func (*LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]) Implements_graph_Annotatable() {}

func (this_ *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]) initFields_LinkedUndirectedGraph_LinkedUndirectedGraphNode() {
	_ = this_
	this_.neighborEdges = j2gort.NewExt_java_util_ArrayList()
}

// Java constructor: LinkedUndirectedGraphNode(N nodeValue, boolean useAnnotation)
func NewLinkedUndirectedGraph_LinkedUndirectedGraphNode[N comparable, E any](nodeValue N, useAnnotation bool) *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E] {
	this_ := &LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]{}
	InitLinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E](this_, nodeValue, useAnnotation)
	return this_
}
func allocPooled_LinkedUndirectedGraph_LinkedUndirectedGraphNode[N comparable, E any](slab_ *j2gort.PooledSlab[LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]], nodeValue N, useAnnotation bool) *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E] {
	this_ := slab_.Alloc()
	InitLinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E](this_, nodeValue, useAnnotation)
	return this_
}

func InitLinkedUndirectedGraph_LinkedUndirectedGraphNode[N comparable, E any](this_ *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E], nodeValue N, useAnnotation bool) {
	_ = nodeValue
	_ = useAnnotation
	this_.initFields_LinkedUndirectedGraph_LinkedUndirectedGraphNode()
	this_.value = nodeValue
	this_.useAnnotation = useAnnotation
}

// Java method: LinkedUndirectedGraphNode.getNeighborEdges()
func (this_ *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]) GetNeighborEdges() j2gort.Ext_java_util_List {
	return j2gort.AsIface[j2gort.Ext_java_util_List](this_.neighborEdges)
}

// Java method: LinkedUndirectedGraphNode.getNeighborEdgesIterator()
func (this_ *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]) GetNeighborEdgesIterator() j2gort.Ext_java_util_Iterator {
	return j2gort.AsIface[j2gort.Ext_java_util_Iterator](this_.neighborEdges.Iterator())
}

// Java method: LinkedUndirectedGraphNode.getAnnotation()
func (this_ *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]) GetAnnotation() Annotation {
	if !(this_.useAnnotation) {
		panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException", "Graph initialized with node annotations turned off")))
	}
	return this_.annotation
}

// Java method: LinkedUndirectedGraphNode.setAnnotation(com.google.javascript.jscomp.graph.Annotation data)
func (this_ *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]) SetAnnotation(data Annotation) {
	_ = data
	if !(this_.useAnnotation) {
		panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException", "Graph initialized with node annotations turned off")))
	}
	this_.annotation = data
}

// Java method: LinkedUndirectedGraphNode.getValue()
func (this_ *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]) GetValue() N {
	return this_.value
}

// Java method: LinkedUndirectedGraphNode.getColor()
func (this_ *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]) GetColor() string {
	return "white"
}

// Java method: LinkedUndirectedGraphNode.getId()
func (this_ *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]) GetId() string {
	return ("LDN" + j2gort.StrValInt32(j2gort.IdentityHash(this_)))
}

// Java method: LinkedUndirectedGraphNode.getLabel()
func (this_ *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]) GetLabel() string {
	return j2gort.StrVal(this_.value)
}

// Java method: LinkedUndirectedGraphNode.neighborList()
func (this_ *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]) priv_LinkedUndirectedGraphNode_neighborList_6370() j2gort.Ext_java_util_List {
	var result []GraphNode[N, E] = j2gort.MakeNonNilSlice[GraphNode[N, E]](this_.neighborEdges.Size())
	_ = result
	{
		var i int32 = int32(0)
		_ = i
		var size int32 = this_.neighborEdges.Size()
		_ = size
		for ; i < size; i = i + 1 {
			var edge UndiGraph_UndiGraphEdge[N, E] = j2gort.AsIface[UndiGraph_UndiGraphEdge[N, E]](this_.neighborEdges.Get(i))
			_ = edge
			if j2gort.IfaceEquals[GraphNode[N, E]](edge.GetNodeA(), this_) {
				result = append(result, edge.GetNodeB())
			} else {
				result = append(result, edge.GetNodeA())
			}
		}
	}
	return j2gort.WrapTypedArrayList(result)
}

// Java class: com.google.javascript.jscomp.graph.LinkedUndirectedGraph.LinkedUndirectedGraphEdge
type LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N comparable, E any] struct {
	nodeA         UndiGraph_UndiGraphNode[N, E] `gson:"nodeA"`
	nodeB         UndiGraph_UndiGraphNode[N, E] `gson:"nodeB"`
	value         E                             `gson:"value"`
	annotation    Annotation                    `gson:"annotation"`
	idHash_       atomic.Int32
	useAnnotation bool `gson:"useAnnotation"`
}

func (this_ *LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E]) BindSelf_(s any) {
	_ = s
}
func (this_ *LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E]) IsJavaNil_() bool {
	return this_ == nil
}
func (*LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E]) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.graph.LinkedUndirectedGraph$LinkedUndirectedGraphEdge", "LinkedUndirectedGraphEdge"
}
func (this_ *LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E]) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E]) HasIdentityEquals_() {}
func (*LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E]) IsGenericInst_graph_LinkedUndirectedGraph_LinkedUndirectedGraphEdge_() {
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.graph.LinkedUndirectedGraph.LinkedUndirectedGraphEdge", "com.google.javascript.jscomp.graph.LinkedUndirectedGraph$LinkedUndirectedGraphEdge", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			IsGenericInst_graph_LinkedUndirectedGraph_LinkedUndirectedGraphEdge_()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}
func (*LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E]) Implements_graph_UndiGraph_UndiGraphEdge() {
}
func (*LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E]) Implements_graph_GraphvizGraph_GraphvizEdge() {
}
func (*LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E]) Implements_graph_Graph_GraphEdge() {}
func (*LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E]) Implements_graph_Annotatable()     {}

func (this_ *LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E]) initFields_LinkedUndirectedGraph_LinkedUndirectedGraphEdge() {
	_ = this_
}

// Java constructor: LinkedUndirectedGraphEdge(com.google.javascript.jscomp.graph.UndiGraph.UndiGraphNode<N,E> nodeA, E edgeValue, com.google.javascript.jscomp.graph.UndiGraph.UndiGraphNode<N,E> nodeB, boolean useAnnotation)
func NewLinkedUndirectedGraph_LinkedUndirectedGraphEdge[N comparable, E any](nodeA UndiGraph_UndiGraphNode[N, E], edgeValue E, nodeB UndiGraph_UndiGraphNode[N, E], useAnnotation bool) *LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E] {
	this_ := &LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E]{}
	InitLinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E](this_, nodeA, edgeValue, nodeB, useAnnotation)
	return this_
}
func allocPooled_LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N comparable, E any](slab_ *j2gort.PooledSlab[LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E]], nodeA UndiGraph_UndiGraphNode[N, E], edgeValue E, nodeB UndiGraph_UndiGraphNode[N, E], useAnnotation bool) *LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E] {
	this_ := slab_.Alloc()
	InitLinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E](this_, nodeA, edgeValue, nodeB, useAnnotation)
	return this_
}

func InitLinkedUndirectedGraph_LinkedUndirectedGraphEdge[N comparable, E any](this_ *LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E], nodeA UndiGraph_UndiGraphNode[N, E], edgeValue E, nodeB UndiGraph_UndiGraphNode[N, E], useAnnotation bool) {
	_ = nodeA
	_ = edgeValue
	_ = nodeB
	_ = useAnnotation
	this_.initFields_LinkedUndirectedGraph_LinkedUndirectedGraphEdge()
	this_.value = edgeValue
	this_.nodeA = nodeA
	this_.nodeB = nodeB
	this_.useAnnotation = useAnnotation
}

// Java method: LinkedUndirectedGraphEdge.getValue()
func (this_ *LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E]) GetValue() E {
	return this_.value
}

// Java method: LinkedUndirectedGraphEdge.getNodeA()
func (this_ *LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E]) GetNodeA() GraphNode[N, E] {
	return j2gort.AsIface[GraphNode[N, E]](this_.nodeA)
}

// Java method: LinkedUndirectedGraphEdge.getNodeB()
func (this_ *LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E]) GetNodeB() GraphNode[N, E] {
	return j2gort.AsIface[GraphNode[N, E]](this_.nodeB)
}

// Java method: LinkedUndirectedGraphEdge.getAnnotation()
func (this_ *LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E]) GetAnnotation() Annotation {
	if !(this_.useAnnotation) {
		panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException", "Graph initialized with edge annotations turned off")))
	}
	return this_.annotation
}

// Java method: LinkedUndirectedGraphEdge.setAnnotation(com.google.javascript.jscomp.graph.Annotation data)
func (this_ *LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E]) SetAnnotation(data Annotation) {
	_ = data
	if !(this_.useAnnotation) {
		panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException", "Graph initialized with edge annotations turned off")))
	}
	this_.annotation = data
}

// Java method: LinkedUndirectedGraphEdge.getColor()
func (this_ *LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E]) GetColor() string {
	return "black"
}

// Java method: LinkedUndirectedGraphEdge.getLabel()
func (this_ *LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E]) GetLabel() string {
	return j2gort.StrVal(this_.value)
}

// Java method: LinkedUndirectedGraphEdge.getNode1Id()
func (this_ *LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E]) GetNode1Id() string {
	return (j2gort.CastPtr[LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]](this_.nodeA)).GetId()
}

// Java method: LinkedUndirectedGraphEdge.getNode2Id()
func (this_ *LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E]) GetNode2Id() string {
	return (j2gort.CastPtr[LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]](this_.nodeB)).GetId()
}

// Java method: LinkedUndirectedGraphEdge.toString()
func (this_ *LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E]) String() string {
	return (j2gort.StrVal(this_.nodeA) + " -- " + j2gort.StrVal(this_.nodeB))
}
func (this_ *LinkedUndirectedGraph[N, E]) ReleaseSlabs_() {
	if this_ == nil {
		return
	}
	this_.slab_LinkedUndirectedGraph_LinkedUndirectedGraphEdge.Release()
	this_.slab_LinkedUndirectedGraph_LinkedUndirectedGraphNode.Release()
	this_.nodes.ReleaseToPool()
}
