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
// Port of com/google/javascript/jscomp/graph/LinkedDirectedGraph.java (generated mechanically by j2go).

package graph

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java polymorphic class interface: com.google.javascript.jscomp.graph.LinkedDirectedGraph
type LinkedDirectedGraph[N comparable, E any] interface {
	Implements_graph_LinkedDirectedGraph()
	IsJavaNil_() bool
	GetLinkedDirectedGraphBase() *LinkedDirectedGraphBase[N, E]
	DiGraph[N, E]
	GraphvizGraph
	SetUseEdgeAnnotations(bool)
	ResetForReuse(bool, bool) bool
	Connect_3_DiGraph_DiGraphNode_any_any__any_DiGraph_DiGraphNode_any_any_(DiGraph_DiGraphNode[N, E], E, DiGraph_DiGraphNode[N, E])
	ConnectIfNotConnectedInDirection(N, E, N)
	CollectNodes([]*LinkedDirectedGraph_LinkedDiGraphNode[N, E]) int32
	IsConnectedInDirection_3_LinkedDirectedGraph_LinkedDiGraphNode_any_any__Ext_com_google_common_base_Predicate_LinkedDirectedGraph_LinkedDiGraphNode_any_any_(*LinkedDirectedGraph_LinkedDiGraphNode[N, E], j2gort.Ext_com_google_common_base_Predicate, *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) bool
}

func ToIface_LinkedDirectedGraph[N comparable, E any, T_iface_ interface {
	LinkedDirectedGraph[N, E]
	comparable
}](v T_iface_) LinkedDirectedGraph[N, E] {
	var z T_iface_
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.jscomp.graph.LinkedDirectedGraph
type LinkedDirectedGraphBase[N comparable, E any] struct {
	DiGraphBase[N, E]
	nodes                                      j2gort.OrderedMap[N, *LinkedDirectedGraph_LinkedDiGraphNode[N, E]] `gson:"nodes"`
	nodesView                                  j2gort.Ext_java_util_Collection                                    `gson:"nodesView"`
	freeNodes                                  j2gort.Ext_java_util_ArrayList                                     `gson:"freeNodes"`
	freeEdges                                  j2gort.Ext_java_util_ArrayList                                     `gson:"freeEdges"`
	useNodeAnnotations                         bool                                                               `gson:"useNodeAnnotations"`
	useEdgeAnnotations                         bool                                                               `gson:"useEdgeAnnotations"`
	slab_LinkedDirectedGraph_LinkedDiGraphEdge j2gort.PooledSlab[LinkedDirectedGraph_LinkedDiGraphEdge[N, E]]
	slab_LinkedDirectedGraph_LinkedDiGraphNode j2gort.PooledSlab[LinkedDirectedGraph_LinkedDiGraphNode[N, E]]
}

func (this_ *LinkedDirectedGraphBase[N, E]) BindSelf_(s any) {
	_ = s
	this_.DiGraphBase.BindSelf_(s)
}
func (this_ *LinkedDirectedGraphBase[N, E]) IsJavaNil_() bool { return this_ == nil }
func (*LinkedDirectedGraphBase[N, E]) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.graph.LinkedDirectedGraph", "LinkedDirectedGraph"
}
func (this_ *LinkedDirectedGraphBase[N, E]) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.DiGraphBase.IdentityHash_()
}
func (*LinkedDirectedGraphBase[N, E]) IsGenericInst_graph_LinkedDirectedGraph_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.graph.LinkedDirectedGraph", "com.google.javascript.jscomp.graph.LinkedDirectedGraph", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			IsGenericInst_graph_LinkedDirectedGraph_()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}
func (b *LinkedDirectedGraphBase[N, E]) GetLinkedDirectedGraphBase() *LinkedDirectedGraphBase[N, E] {
	return b
}
func (*LinkedDirectedGraphBase[N, E]) Implements_graph_LinkedDirectedGraph() {}
func (*LinkedDirectedGraphBase[N, E]) Implements_graph_DiGraph()             {}
func (*LinkedDirectedGraphBase[N, E]) Implements_graph_GraphvizGraph()       {}
func (*LinkedDirectedGraphBase[N, E]) Implements_graph_Graph()               {}
func (*LinkedDirectedGraphBase[N, E]) Implements_graph_AdjacencyGraph()      {}

func StaticGet_LinkedDirectedGraph_MAX_POOLED_GRAPH_NODES() int32 { return int32(4096) }
func StaticGet_LinkedDirectedGraph_MAX_POOLED_GRAPH_EDGES() int32 { return int32(8192) }
func (this_ *LinkedDirectedGraphBase[N, E]) initFields_LinkedDirectedGraph() {
	_ = this_
	this_.nodes = j2gort.AcquireOrderedMap[N, *LinkedDirectedGraph_LinkedDiGraphNode[

		// Java constructor: LinkedDirectedGraph(boolean useNodeAnnotations, boolean useEdgeAnnotations)
		N, E]]()
}

func NewLinkedDirectedGraph[N comparable, E any](useNodeAnnotations bool, useEdgeAnnotations bool) LinkedDirectedGraph[N, E] {
	this_ := &LinkedDirectedGraphBase[N, E]{}
	this_.BindSelf_(this_)
	InitLinkedDirectedGraph[N, E](this_, useNodeAnnotations, useEdgeAnnotations)
	return this_
}

func InitLinkedDirectedGraph[N comparable, E any](this_ *LinkedDirectedGraphBase[N, E], useNodeAnnotations bool, useEdgeAnnotations bool) {
	_ = useNodeAnnotations
	_ = useEdgeAnnotations
	InitDiGraph[N, E](&this_.DiGraphBase)
	this_.initFields_LinkedDirectedGraph()
	this_.useNodeAnnotations = useNodeAnnotations
	this_.useEdgeAnnotations = useEdgeAnnotations
}

// Java method: LinkedDirectedGraph.newSubGraph()
func (this_ *LinkedDirectedGraphBase[N, E]) NewSubGraph() SubGraph[N, E] {
	return ToIface_SubGraph[N, E](NewGraph_SimpleSubGraph[N, E](this_.GraphBase.self_))
}

// Java static method: LinkedDirectedGraph.createWithoutAnnotations()
func LinkedDirectedGraph_CreateWithoutAnnotations[N comparable, E any]() LinkedDirectedGraph[N, E] {
	return NewLinkedDirectedGraph[N, E](false, false)
}

// Java static method: LinkedDirectedGraph.create()
func LinkedDirectedGraph_Create[N comparable, E any]() LinkedDirectedGraph[N, E] {
	return NewLinkedDirectedGraph[N, E](true, true)
}

// Java method: LinkedDirectedGraph.setUseEdgeAnnotations(boolean useEdgeAnnotations)
func (this_ *LinkedDirectedGraphBase[N, E]) SetUseEdgeAnnotations(useEdgeAnnotations bool) {
	_ = useEdgeAnnotations
	this_.useEdgeAnnotations = useEdgeAnnotations
}

// Java method: LinkedDirectedGraph.resetForReuse(boolean useNodeAnnotations, boolean useEdgeAnnotations)
func (this_ *LinkedDirectedGraphBase[N, E]) ResetForReuse(useNodeAnnotations bool, useEdgeAnnotations bool) bool {
	_ = useNodeAnnotations
	_ = useEdgeAnnotations
	var size int32 = this_.nodes.Size()
	_ = size
	if size > int32(4096) {
		return false
	}
	if size > int32(0) {
		this_.freeNodes = nil
		this_.freeEdges = nil
		this_.slab_LinkedDirectedGraph_LinkedDiGraphEdge.Reset()
		this_.slab_LinkedDirectedGraph_LinkedDiGraphNode.Reset()
		this_.nodes.Clear()

	}
	this_.useNodeAnnotations = useNodeAnnotations
	this_.useEdgeAnnotations = useEdgeAnnotations
	return true
}

// Java method: LinkedDirectedGraph.acquireEdge(com.google.javascript.jscomp.graph.LinkedDirectedGraph.LinkedDiGraphNode<N,E> src, E edgeValue, com.google.javascript.jscomp.graph.LinkedDirectedGraph.LinkedDiGraphNode<N,E> dest)
func (this_ *LinkedDirectedGraphBase[N, E]) priv_LinkedDirectedGraph_acquireEdge_6386(src *LinkedDirectedGraph_LinkedDiGraphNode[N, E], edgeValue E, dest *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) *LinkedDirectedGraph_LinkedDiGraphEdge[N, E] {
	_ = src
	_ = edgeValue
	_ = dest
	if (this_.freeEdges != nil) && !(this_.freeEdges.IsEmpty()) {
		var last int32 = (this_.freeEdges.Size() - int32(1))
		_ = last
		var edge *LinkedDirectedGraph_LinkedDiGraphEdge[N, E] = j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphEdge[N, E]](this_.freeEdges.Get(last))
		_ = edge
		this_.freeEdges.Remove_1_int32(last)
		edge.Reinit(src, edgeValue, dest, this_.useEdgeAnnotations)
		return edge
	}
	return allocPooled_LinkedDirectedGraph_LinkedDiGraphEdge[N, E](&this_.slab_LinkedDirectedGraph_LinkedDiGraphEdge, src, edgeValue, dest, this_.useEdgeAnnotations)
}

// Java method: LinkedDirectedGraph.connect(N srcValue, E edgeValue, N destValue)
func (this_ *LinkedDirectedGraphBase[N, E]) Connect_3_any_any_any(srcValue N, edgeValue E, destValue N) {
	_ = srcValue
	_ = edgeValue
	_ = destValue
	var src *LinkedDirectedGraph_LinkedDiGraphNode[N, E] = j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphNode[N, E]](this_.GetNodeOrFail(srcValue))
	_ = src
	var dest *LinkedDirectedGraph_LinkedDiGraphNode[N, E] = j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphNode[N, E]](this_.GetNodeOrFail(destValue))
	_ = dest
	var edge *LinkedDirectedGraph_LinkedDiGraphEdge[N, E] = this_.priv_LinkedDirectedGraph_acquireEdge_6386(src, edgeValue, dest)
	_ = edge
	src.AddOutEdge(edge)
	dest.AddInEdge(edge)
}

// Java method: LinkedDirectedGraph.connect(com.google.javascript.jscomp.graph.DiGraph.DiGraphNode<N,E> src, E edgeValue, com.google.javascript.jscomp.graph.DiGraph.DiGraphNode<N,E> dest)
func (this_ *LinkedDirectedGraphBase[N, E]) Connect_3_DiGraph_DiGraphNode_any_any__any_DiGraph_DiGraphNode_any_any_(src DiGraph_DiGraphNode[N, E], edgeValue E, dest DiGraph_DiGraphNode[N, E]) {
	_ = src
	_ = edgeValue
	_ = dest
	var lSrc *LinkedDirectedGraph_LinkedDiGraphNode[N, E] = j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphNode[N, E]](src)
	_ = lSrc
	var lDest *LinkedDirectedGraph_LinkedDiGraphNode[N, E] = j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphNode[N, E]](dest)
	_ = lDest
	var edge *LinkedDirectedGraph_LinkedDiGraphEdge[N, E] = this_.priv_LinkedDirectedGraph_acquireEdge_6386(lSrc, edgeValue, lDest)
	_ = edge
	lSrc.AddOutEdge(edge)
	lDest.AddInEdge(edge)
}

// Java method: LinkedDirectedGraph.connectIfNotConnectedInDirection(N srcValue, E edgeValue, N destValue)
func (this_ *LinkedDirectedGraphBase[N, E]) ConnectIfNotConnectedInDirection(srcValue N, edgeValue E, destValue N) {
	_ = srcValue
	_ = edgeValue
	_ = destValue
	var src *LinkedDirectedGraph_LinkedDiGraphNode[N, E] = j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphNode[N, E]](this_.CreateNode(srcValue))
	_ = src
	var dest *LinkedDirectedGraph_LinkedDiGraphNode[N, E] = j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphNode[N, E]](this_.CreateNode(destValue))
	_ = dest
	if !(this_.IsConnectedInDirection_3_LinkedDirectedGraph_LinkedDiGraphNode_any_any__Ext_com_google_common_base_Predicate_LinkedDirectedGraph_LinkedDiGraphNode_any_any_(src, j2gort.Ext_com_google_common_base_Predicates_EqualTo(j2gort.ToAny(edgeValue)), dest)) {
		this_.Connect_3_DiGraph_DiGraphNode_any_any__any_DiGraph_DiGraphNode_any_any_(ToIface_DiGraph_DiGraphNode[N, E](src), edgeValue, ToIface_DiGraph_DiGraphNode[N, E](dest))
	}
}

// Java method: LinkedDirectedGraph.disconnect(N n1, N n2)
func (this_ *LinkedDirectedGraphBase[N, E]) Disconnect(n1 N, n2 N) {
	_ = n1
	_ = n2
	this_.DisconnectInDirection(n1, n2)
	this_.DisconnectInDirection(n2, n1)
}

// Java method: LinkedDirectedGraph.disconnectInDirection(N srcValue, N destValue)
func (this_ *LinkedDirectedGraphBase[N, E]) DisconnectInDirection(srcValue N, destValue N) {
	_ = srcValue
	_ = destValue
	var src *LinkedDirectedGraph_LinkedDiGraphNode[N, E] = j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphNode[N, E]](this_.GetNodeOrFail(srcValue))
	_ = src
	var dest *LinkedDirectedGraph_LinkedDiGraphNode[N, E] = j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphNode[N, E]](this_.GetNodeOrFail(destValue))
	_ = dest
	{
		cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.GetEdgesInDirection(srcValue, destValue)))
		for cur_3.HasNext() {
			rawElem_4 := cur_3.Next()
			edge := j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphEdge[N, E]](rawElem_4)
			_ = edge
			src.RemoveOutEdge(edge)
			dest.RemoveInEdge(edge)
		}
	}
}

// Java method: LinkedDirectedGraph.getNodes()
func (this_ *LinkedDirectedGraphBase[N, E]) GetNodes() j2gort.Ext_java_util_Collection {
	var view j2gort.Ext_java_util_Collection = this_.nodesView
	_ = view
	if view == nil {
		view = j2gort.Ext_java_util_Collections_UnmodifiableCollection(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.nodes.LiveValues()))
		this_.nodesView = view
	}
	return view
}

// Java method: LinkedDirectedGraph.collectNodes(com.google.javascript.jscomp.graph.LinkedDirectedGraph.LinkedDiGraphNode<N,E>[] out)
func (this_ *LinkedDirectedGraphBase[N, E]) CollectNodes(out []*LinkedDirectedGraph_LinkedDiGraphNode[N, E]) int32 {
	_ = out
	var i int32 = int32(0)
	_ = i
	{
		cur_5 := this_.nodes.IterValues()
		for cur_5.HasNext() {

			node := cur_5.
				Next()
			_ = node
			{
				a_7 := out
				var i_8 int32
				{

					old_9 := i
					nw_10 := old_9 + 1
					i = nw_10
					i_8 = old_9
				}
				a_7[i_8] = node
			}
		}
	}
	return i
}

// Java method: LinkedDirectedGraph.getNode(N nodeValue)
func (this_ *LinkedDirectedGraphBase[N, E]) GetNode(nodeValue N) any {
	_ = nodeValue
	return j2gort.PtrToAny(this_.nodes.Get(nodeValue))
}

// Java method: LinkedDirectedGraph.getInEdges(N nodeValue)
func (this_ *LinkedDirectedGraphBase[N, E]) GetInEdges(nodeValue N) j2gort.Ext_java_util_List {
	_ = nodeValue
	var node *LinkedDirectedGraph_LinkedDiGraphNode[N, E] = j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphNode[N, E]](this_.GetNodeOrFail(nodeValue))
	_ = node
	if node.inEdges != nil {
		return j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.ImmListFromSlice(node.inEdges))
	}
	if node.firstInEdge == nil {
		return j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
	} else {
		return node.firstInEdge.AsSingletonList()
	}
}

// Java method: LinkedDirectedGraph.getOutEdges(N nodeValue)
func (this_ *LinkedDirectedGraphBase[N, E]) GetOutEdges(nodeValue N) j2gort.Ext_java_util_List {
	_ = nodeValue
	var node *LinkedDirectedGraph_LinkedDiGraphNode[N, E] = j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphNode[N, E]](this_.GetNodeOrFail(nodeValue))
	_ = node
	if node.outEdges != nil {
		return j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.ImmListFromSlice(node.outEdges))
	}
	if node.firstOutEdge == nil {
		return j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
	} else {
		return node.firstOutEdge.AsSingletonList()
	}
}

// Java method: LinkedDirectedGraph.createNode(N nodeValue)
func (this_ *LinkedDirectedGraphBase[N, E]) CreateNode(nodeValue N) any {
	_ = nodeValue
	var node *LinkedDirectedGraph_LinkedDiGraphNode[N, E] = this_.nodes.Get(nodeValue)
	_ = node
	if node == nil {
		if (this_.freeNodes != nil) && !(this_.freeNodes.IsEmpty()) {
			var last int32 = (this_.freeNodes.Size() - int32(1))
			_ = last
			node = j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphNode[N, E]](this_.freeNodes.Get(last))
			this_.freeNodes.Remove_1_int32(last)
			node.Reinit(nodeValue, this_.useNodeAnnotations)
		} else {
			node = allocPooled_LinkedDirectedGraph_LinkedDiGraphNode[N, E](&this_.slab_LinkedDirectedGraph_LinkedDiGraphNode, nodeValue, this_.useNodeAnnotations)
		}
		this_.nodes.Put(nodeValue, node)
	}
	return j2gort.PtrToAny(node)
}

// Java method: LinkedDirectedGraph.getEdges(N n1, N n2)
func (this_ *LinkedDirectedGraphBase[N, E]) GetEdges_2(n1 N, n2 N) j2gort.Ext_java_util_List {
	_ = n1
	_ = n2
	var forwardEdges j2gort.Ext_java_util_List = this_.GetEdgesInDirection(n1, n2)
	_ = forwardEdges
	var backwardEdges j2gort.Ext_java_util_List = this_.GetEdgesInDirection(n2, n1)
	_ = backwardEdges
	var totalSize int32 = (forwardEdges.Size() + backwardEdges.Size())
	_ = totalSize
	var edges []*LinkedDirectedGraph_LinkedDiGraphEdge[N, E] = j2gort.MakeNonNilSlice[*LinkedDirectedGraph_LinkedDiGraphEdge[N, E]](totalSize)
	_ = edges
	{
		var i int32 = int32(0)
		_ = i
		var size int32 = forwardEdges.Size()
		_ = size
		for ; i < size; i = i + 1 {
			edges = append(edges, j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphEdge[N, E]](forwardEdges.Get(i)))
		}
	}
	{
		var i_2 int32 = int32(0)
		_ = i_2
		var size_2 int32 = backwardEdges.Size()
		_ = size_2
		for ; i_2 < size_2; i_2 = i_2 + 1 {
			edges = append(edges, j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphEdge[N, E]](backwardEdges.Get(i_2)))
		}
	}
	return j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.ImmListFromSlice(edges))
}

// Java method: LinkedDirectedGraph.getEdges()
func (this_ *LinkedDirectedGraphBase[N, E]) GetEdges_0() j2gort.Ext_java_util_List {
	var result []*LinkedDirectedGraph_LinkedDiGraphEdge[N, E] = []*LinkedDirectedGraph_LinkedDiGraphEdge[N, E]{}
	_ = result
	{
		cur_11 := this_.nodes.IterValues()
		for cur_11.HasNext() {

			node := cur_11.
				Next()
			_ = node
			if node.outEdges != nil {
				{
					var i int32 = int32(0)
					_ = i
					var size int32 = int32(len(node.outEdges))
					_ = size
					for ; i < size; i = i + 1 {
						result = append(result, node.outEdges[i])
					}
				}
			} else if node.firstOutEdge != nil {
				result = append(result, node.firstOutEdge)
			}
		}
	}
	return j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.ImmListFromSlice(result))
}

// Java method: LinkedDirectedGraph.getFirstEdge(N n1, N n2)
func (this_ *LinkedDirectedGraphBase[N, E]) GetFirstEdge(n1 N, n2 N) Graph_GraphEdge[N, E] {
	_ = n1
	_ = n2
	var dNode1 *LinkedDirectedGraph_LinkedDiGraphNode[N, E] = j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphNode[N, E]](this_.GetNodeOrFail(n1))
	_ = dNode1
	var dNode2 *LinkedDirectedGraph_LinkedDiGraphNode[N, E] = j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphNode[N, E]](this_.GetNodeOrFail(n2))
	_ = dNode2
	if dNode1.outEdges != nil {
		{
			var i int32 = int32(0)
			_ = i
			var size int32 = int32(len(dNode1.outEdges))
			_ = size
			for ; i < size; i = i + 1 {
				var outEdge *LinkedDirectedGraph_LinkedDiGraphEdge[N, E] = dNode1.outEdges[i]
				_ = outEdge
				if j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphNode[N, E]](outEdge.GetDestination()) == dNode2 {
					return ToIface_Graph_GraphEdge[N, E](outEdge)
				}
			}
		}
	} else if (dNode1.firstOutEdge != nil) && (j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphNode[N, E]](dNode1.firstOutEdge.GetDestination()) == dNode2) {
		return ToIface_Graph_GraphEdge[N, E](dNode1.firstOutEdge)
	}
	if dNode2.outEdges != nil {
		{
			var i_2 int32 = int32(0)
			_ = i_2
			var size_2 int32 = int32(len(dNode2.outEdges))
			_ = size_2
			for ; i_2 < size_2; i_2 = i_2 + 1 {
				var outEdge_2 *LinkedDirectedGraph_LinkedDiGraphEdge[N, E] = dNode2.outEdges[i_2]
				_ = outEdge_2
				if j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphNode[N, E]](outEdge_2.GetDestination()) == dNode1 {
					return ToIface_Graph_GraphEdge[N, E](outEdge_2)
				}
			}
		}
	} else if (dNode2.firstOutEdge != nil) && (j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphNode[N, E]](dNode2.firstOutEdge.GetDestination()) == dNode1) {
		return ToIface_Graph_GraphEdge[N, E](dNode2.firstOutEdge)
	}
	return nil
}

// Java method: LinkedDirectedGraph.getEdgesInDirection(N n1, N n2)
func (this_ *LinkedDirectedGraphBase[N, E]) GetEdgesInDirection(n1 N, n2 N) j2gort.Ext_java_util_List {
	_ = n1
	_ = n2
	var dNode1 *LinkedDirectedGraph_LinkedDiGraphNode[N, E] = j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphNode[N, E]](this_.GetNodeOrFail(n1))
	_ = dNode1
	var dNode2 *LinkedDirectedGraph_LinkedDiGraphNode[N, E] = j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphNode[N, E]](this_.GetNodeOrFail(n2))
	_ = dNode2
	if dNode1.outEdges == nil {
		if (dNode1.firstOutEdge != nil) && (j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphNode[N, E]](dNode1.firstOutEdge.GetDestination()) == dNode2) {
			return dNode1.firstOutEdge.AsSingletonList()
		}
		return j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
	}
	var edges []*LinkedDirectedGraph_LinkedDiGraphEdge[N, E] = []*LinkedDirectedGraph_LinkedDiGraphEdge[N, E]{}
	_ = edges
	{
		var i int32 = int32(0)
		_ = i
		var size int32 = int32(len(dNode1.outEdges))
		_ = size
		for ; i < size; i = i + 1 {
			var outEdge *LinkedDirectedGraph_LinkedDiGraphEdge[N, E] = dNode1.outEdges[i]
			_ = outEdge
			if j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphNode[N, E]](outEdge.GetDestination()) == dNode2 {
				edges = append(edges, outEdge)
			}
		}
	}
	return j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.ImmListFromSlice(edges))
}

// Java method: LinkedDirectedGraph.isConnectedInDirection(N source, N dest)
func (this_ *LinkedDirectedGraphBase[N, E]) IsConnectedInDirection_2(source N, dest N) bool {
	_ = source
	_ = dest
	return this_.priv_LinkedDirectedGraph_isConnectedInDirection_6391(source, j2gort.Ext_com_google_common_base_Predicates_AlwaysTrue(), dest)
}

// Java method: LinkedDirectedGraph.isConnectedInDirection(N source, E edgeValue, N dest)
func (this_ *LinkedDirectedGraphBase[N, E]) IsConnectedInDirection_3_any_any_any(source N, edgeValue E, dest N) bool {
	_ = source
	_ = edgeValue
	_ = dest
	return this_.priv_LinkedDirectedGraph_isConnectedInDirection_6391(source, j2gort.Ext_com_google_common_base_Predicates_EqualTo(j2gort.ToAny(edgeValue)), dest)
}

// Java method: LinkedDirectedGraph.isConnectedInDirection(com.google.javascript.jscomp.graph.LinkedDirectedGraph.LinkedDiGraphNode<N,E> source, com.google.common.base.Predicate<E> edgeFilter, com.google.javascript.jscomp.graph.LinkedDirectedGraph.LinkedDiGraphNode<N,E> dest)
func (this_ *LinkedDirectedGraphBase[N, E]) IsConnectedInDirection_3_LinkedDirectedGraph_LinkedDiGraphNode_any_any__Ext_com_google_common_base_Predicate_LinkedDirectedGraph_LinkedDiGraphNode_any_any_(source *LinkedDirectedGraph_LinkedDiGraphNode[N, E], edgeFilter j2gort.Ext_com_google_common_base_Predicate, dest *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) bool {
	_ = source
	_ = edgeFilter
	_ = dest
	var outEdgesLen int32 = source.OutEdgeCount()
	_ = outEdgesLen
	var inEdgesLen int32 = dest.InEdgeCount()
	_ = inEdgesLen
	if outEdgesLen < inEdgesLen {
		if source.outEdges != nil {
			{
				var i int32 = int32(0)
				_ = i
				for ; i < outEdgesLen; i = i + 1 {
					var outEdge *LinkedDirectedGraph_LinkedDiGraphEdge[N, E] = source.outEdges[i]
					_ = outEdge
					if (j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphNode[N, E]](outEdge.GetDestination()) == dest) && edgeFilter.Apply(j2gort.ToAny(outEdge.GetValue())) {
						return true
					}
				}
			}
		} else if source.firstOutEdge != nil {
			return ((j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphNode[N, E]](source.firstOutEdge.GetDestination()) == dest) && edgeFilter.Apply(j2gort.ToAny(source.firstOutEdge.GetValue())))
		}
	} else {
		if dest.inEdges != nil {
			{
				var i_2 int32 = int32(0)
				_ = i_2
				for ; i_2 < inEdgesLen; i_2 = i_2 + 1 {
					var inEdge *LinkedDirectedGraph_LinkedDiGraphEdge[N, E] = dest.inEdges[i_2]
					_ = inEdge
					if (j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphNode[N, E]](inEdge.GetSource()) == source) && edgeFilter.Apply(j2gort.ToAny(inEdge.GetValue())) {
						return true
					}
				}
			}
		} else if dest.firstInEdge != nil {
			return ((j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphNode[N, E]](dest.firstInEdge.GetSource()) == source) && edgeFilter.Apply(j2gort.ToAny(dest.firstInEdge.GetValue())))
		}
	}
	return false
}

// Java method: LinkedDirectedGraph.isConnectedInDirection(N n1, com.google.common.base.Predicate<E> edgeMatcher, N n2)
func (this_ *LinkedDirectedGraphBase[N, E]) priv_LinkedDirectedGraph_isConnectedInDirection_6391(n1 N, edgeMatcher j2gort.Ext_com_google_common_base_Predicate, n2 N) bool {
	_ = n1
	_ = edgeMatcher
	_ = n2
	var dNode1 *LinkedDirectedGraph_LinkedDiGraphNode[N, E] = j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphNode[N, E]](this_.GetNodeOrFail(n1))
	_ = dNode1
	var dNode2 *LinkedDirectedGraph_LinkedDiGraphNode[N, E] = j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphNode[N, E]](this_.GetNodeOrFail(n2))
	_ = dNode2
	return this_.IsConnectedInDirection_3_LinkedDirectedGraph_LinkedDiGraphNode_any_any__Ext_com_google_common_base_Predicate_LinkedDirectedGraph_LinkedDiGraphNode_any_any_(dNode1, edgeMatcher, dNode2)
}

// Java method: LinkedDirectedGraph.getDirectedPredNodes(N nodeValue)
func (this_ *LinkedDirectedGraphBase[N, E]) GetDirectedPredNodes_1_any(nodeValue N) j2gort.Ext_java_util_List {
	_ = nodeValue
	return this_.GetDirectedPredNodes_1_DiGraph_DiGraphNode_any_any_(ToIface_DiGraph_DiGraphNode[N, E](this_.nodes.Get(nodeValue)))
}

// Java method: LinkedDirectedGraph.getDirectedPredNodes(com.google.javascript.jscomp.graph.DiGraph.DiGraphNode<N,E> dNode)
func (this_ *LinkedDirectedGraphBase[N, E]) GetDirectedPredNodes_1_DiGraph_DiGraphNode_any_any_(dNode DiGraph_DiGraphNode[N, E]) j2gort.Ext_java_util_List {
	_ = dNode
	j2gort.CheckNotNullIface(dNode)
	var lNode *LinkedDirectedGraph_LinkedDiGraphNode[N, E] = j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphNode[N, E]](dNode)
	_ = lNode
	if lNode.inEdges == nil {
		if lNode.firstInEdge != nil {
			return j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphNode[N, E]](lNode.firstInEdge.GetSource()).AsSingletonList()
		}
		return j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
	}
	var size int32 = int32(len(lNode.inEdges))
	_ = size
	var nodeList []*LinkedDirectedGraph_LinkedDiGraphNode[N, E] = j2gort.MakeNonNilSlice[*LinkedDirectedGraph_LinkedDiGraphNode[N, E]](size)
	_ = nodeList
	{
		var i int32 = int32(0)
		_ = i
		for ; i < size; i = i + 1 {
			nodeList = append(nodeList, j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphNode[N, E]](lNode.inEdges[i].GetSource()))
		}
	}
	return j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.ImmListFromSlice(nodeList))
}

// Java method: LinkedDirectedGraph.getDirectedSuccNodes(N nodeValue)
func (this_ *LinkedDirectedGraphBase[N, E]) GetDirectedSuccNodes_1_any(nodeValue N) j2gort.Ext_java_util_List {
	_ = nodeValue
	return this_.GetDirectedSuccNodes_1_DiGraph_DiGraphNode_any_any_(ToIface_DiGraph_DiGraphNode[N, E](this_.nodes.Get(nodeValue)))
}

// Java method: LinkedDirectedGraph.getDirectedSuccNodes(com.google.javascript.jscomp.graph.DiGraph.DiGraphNode<N,E> dNode)
func (this_ *LinkedDirectedGraphBase[N, E]) GetDirectedSuccNodes_1_DiGraph_DiGraphNode_any_any_(dNode DiGraph_DiGraphNode[N, E]) j2gort.Ext_java_util_List {
	_ = dNode
	j2gort.CheckNotNullIface(dNode)
	var lNode *LinkedDirectedGraph_LinkedDiGraphNode[N, E] = j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphNode[N, E]](dNode)
	_ = lNode
	if lNode.outEdges == nil {
		if lNode.firstOutEdge != nil {
			return j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphNode[N, E]](lNode.firstOutEdge.GetDestination()).AsSingletonList()
		}
		return j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
	}
	var size int32 = int32(len(lNode.outEdges))
	_ = size
	var nodeList []*LinkedDirectedGraph_LinkedDiGraphNode[N, E] = j2gort.MakeNonNilSlice[*LinkedDirectedGraph_LinkedDiGraphNode[N, E]](size)
	_ = nodeList
	{
		var i int32 = int32(0)
		_ = i
		for ; i < size; i = i + 1 {
			nodeList = append(nodeList, j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphNode[N, E]](lNode.outEdges[i].GetDestination()))
		}
	}
	return j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.ImmListFromSlice(nodeList))
}

// Java method: LinkedDirectedGraph.getGraphvizEdges()
func (this_ *LinkedDirectedGraphBase[N, E]) GetGraphvizEdges() j2gort.Ext_java_util_List {
	var edgeList []*LinkedDirectedGraph_LinkedDiGraphEdge[N, E] = []*LinkedDirectedGraph_LinkedDiGraphEdge[N, E]{}
	_ = edgeList
	{
		cur_13 := this_.nodes.IterValues()
		for cur_13.HasNext() {

			node := cur_13.
				Next()
			_ = node
			if node.outEdges != nil {
				{
					var i int32 = int32(0)
					_ = i
					var size int32 = int32(len(node.outEdges))
					_ = size
					for ; i < size; i = i + 1 {
						edgeList = append(edgeList, node.outEdges[i])
					}
				}
			} else if node.firstOutEdge != nil {
				edgeList = append(edgeList, node.firstOutEdge)
			}
		}
	}
	return j2gort.WrapTypedArrayList(edgeList)
}

// Java method: LinkedDirectedGraph.getGraphvizNodes()
func (this_ *LinkedDirectedGraphBase[N, E]) GetGraphvizNodes() j2gort.Ext_java_util_List {
	var nodeList []*LinkedDirectedGraph_LinkedDiGraphNode[N, E] = j2gort.MakeNonNilSlice[*LinkedDirectedGraph_LinkedDiGraphNode[N, E]](this_.nodes.Size())
	_ = nodeList
	{
		cur_15 := this_.nodes.IterValues()
		for cur_15.HasNext() {

			node := cur_15.
				Next()
			_ = node
			nodeList = append(nodeList, node)
		}
	}
	return j2gort.WrapTypedArrayList(nodeList)
}

// Java method: LinkedDirectedGraph.getName()
func (this_ *LinkedDirectedGraphBase[N, E]) GetName() string {
	return "LinkedGraph"
}

// Java method: LinkedDirectedGraph.isDirected()
func (this_ *LinkedDirectedGraphBase[N, E]) IsDirected() bool {
	return true
}

// Java method: LinkedDirectedGraph.getNodeCount()
func (this_ *LinkedDirectedGraphBase[N, E]) GetNodeCount() int32 {
	return this_.nodes.Size()
}

// Java method: LinkedDirectedGraph.getNeighborNodes(N value)
func (this_ *LinkedDirectedGraphBase[N, E]) GetNeighborNodes(value N) j2gort.Ext_java_util_List {
	_ = value
	var node *LinkedDirectedGraph_LinkedDiGraphNode[N, E] = j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphNode[N, E]](this_.GetNode(value))
	_ = node
	var result []*LinkedDirectedGraph_LinkedDiGraphNode[N, E] = j2gort.MakeNonNilSlice[*LinkedDirectedGraph_LinkedDiGraphNode[N, E]]((node.InEdgeCount() + node.OutEdgeCount()))
	_ = result
	if node.inEdges != nil {
		{
			var i int32 = int32(0)
			_ = i
			var size int32 = int32(len(node.inEdges))
			_ = size
			for ; i < size; i = i + 1 {
				result = append(result, j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphNode[N, E]](node.inEdges[i].GetSource()))
			}
		}
	} else if node.firstInEdge != nil {
		result = append(result, j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphNode[N, E]](node.firstInEdge.GetSource()))
	}
	if node.outEdges != nil {
		{
			var i_2 int32 = int32(0)
			_ = i_2
			var size_2 int32 = int32(len(node.outEdges))
			_ = size_2
			for ; i_2 < size_2; i_2 = i_2 + 1 {
				result = append(result, j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphNode[N, E]](node.outEdges[i_2].GetDestination()))
			}
		}
	} else if node.firstOutEdge != nil {
		result = append(result, j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphNode[N, E]](node.firstOutEdge.GetDestination()))
	}
	return j2gort.WrapTypedArrayList(result)
}

// Java method: LinkedDirectedGraph.getNodeDegree(N value)
func (this_ *LinkedDirectedGraphBase[N, E]) GetNodeDegree(value N) int32 {
	_ = value
	var node *LinkedDirectedGraph_LinkedDiGraphNode[N, E] = j2gort.CastPtr[LinkedDirectedGraph_LinkedDiGraphNode[N, E]](this_.GetNodeOrFail(value))
	_ = node
	return (node.InEdgeCount() + node.OutEdgeCount())
}

// Java class: com.google.javascript.jscomp.graph.LinkedDirectedGraph.LinkedDiGraphNode
type LinkedDirectedGraph_LinkedDiGraphNode[N comparable, E any] struct {
	firstInEdge       *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]   `gson:"firstInEdge"`
	inEdges           []*LinkedDirectedGraph_LinkedDiGraphEdge[N, E] `gson:"inEdges"`
	firstOutEdge      *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]   `gson:"firstOutEdge"`
	outEdges          []*LinkedDirectedGraph_LinkedDiGraphEdge[N, E] `gson:"outEdges"`
	singletonSelfList j2gort.Ext_java_util_List                      `gson:"singletonSelfList"`
	cachedId          string                                         `gson:"cachedId"`
	value             N                                              `gson:"value"`
	annotation        Annotation                                     `gson:"annotation"`
	priority          int32                                          `gson:"priority"`
	astPosition       int32                                          `gson:"astPosition"`
	dfaIndex          int32                                          `gson:"dfaIndex"`
	idHash_           atomic.Int32
	useAnnotation     bool `gson:"useAnnotation"`
	inDfaQueue        bool `gson:"inDfaQueue"`
}

func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) BindSelf_(s any) {
	_ = s
}
func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) IsJavaNil_() bool { return this_ == nil }
func (*LinkedDirectedGraph_LinkedDiGraphNode[N, E]) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.graph.LinkedDirectedGraph$LinkedDiGraphNode", "LinkedDiGraphNode"
}
func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*LinkedDirectedGraph_LinkedDiGraphNode[N, E]) HasIdentityEquals_() {}
func (*LinkedDirectedGraph_LinkedDiGraphNode[N, E]) IsGenericInst_graph_LinkedDirectedGraph_LinkedDiGraphNode_() {
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.graph.LinkedDirectedGraph.LinkedDiGraphNode", "com.google.javascript.jscomp.graph.LinkedDirectedGraph$LinkedDiGraphNode", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			IsGenericInst_graph_LinkedDirectedGraph_LinkedDiGraphNode_()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}
func (*LinkedDirectedGraph_LinkedDiGraphNode[N, E]) Implements_graph_DiGraph_DiGraphNode()        {}
func (*LinkedDirectedGraph_LinkedDiGraphNode[N, E]) Implements_graph_GraphvizGraph_GraphvizNode() {}
func (*LinkedDirectedGraph_LinkedDiGraphNode[N, E]) Implements_graph_GraphNode()                  {}
func (*LinkedDirectedGraph_LinkedDiGraphNode[N, E]) Implements_graph_Annotatable()                {}

func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) initFields_LinkedDirectedGraph_LinkedDiGraphNode() {
	_ = this_
	this_.cachedId = j2gort.NullString
	this_.priority = int32(-1)
	this_.astPosition = int32(-1)
	this_.dfaIndex = int32(-1)
}

// Java constructor: LinkedDiGraphNode(N nodeValue, boolean useAnnotation)
func NewLinkedDirectedGraph_LinkedDiGraphNode[N comparable, E any](nodeValue N, useAnnotation bool) *LinkedDirectedGraph_LinkedDiGraphNode[N, E] {
	this_ := &LinkedDirectedGraph_LinkedDiGraphNode[N, E]{}
	InitLinkedDirectedGraph_LinkedDiGraphNode[N, E](this_, nodeValue, useAnnotation)
	return this_
}
func allocPooled_LinkedDirectedGraph_LinkedDiGraphNode[N comparable, E any](slab_ *j2gort.PooledSlab[LinkedDirectedGraph_LinkedDiGraphNode[N, E]], nodeValue N, useAnnotation bool) *LinkedDirectedGraph_LinkedDiGraphNode[N, E] {
	this_ := slab_.Alloc()
	this_.cachedId = j2gort.NullString
	this_.priority = int32(-1)
	this_.astPosition = int32(-1)
	this_.dfaIndex = int32(-1)
	this_.value = nodeValue
	if useAnnotation {
		this_.useAnnotation = true
	}
	return this_
}

func InitLinkedDirectedGraph_LinkedDiGraphNode[N comparable, E any](this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E], nodeValue N, useAnnotation bool) {
	_ = nodeValue
	_ = useAnnotation
	this_.initFields_LinkedDirectedGraph_LinkedDiGraphNode()
	this_.value = nodeValue
	this_.useAnnotation = useAnnotation
}

// Java method: LinkedDiGraphNode.reinit(N nodeValue, boolean useAnnotation)
func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) Reinit(nodeValue N, useAnnotation bool) {
	_ = nodeValue
	_ = useAnnotation
	this_.value = nodeValue
	this_.useAnnotation = useAnnotation
}

// Java method: LinkedDiGraphNode.clearForPool()
func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) ClearForPool() {
	this_.firstInEdge = nil
	if this_.inEdges != nil {
		j2gort.ReleasePooledSlice(this_.inEdges)
		this_.inEdges = nil
	}

	this_.firstOutEdge = nil
	if this_.outEdges != nil {
		j2gort.ReleasePooledSlice(this_.outEdges)
		this_.outEdges = nil
	}

	this_.value = func() N {
		if _, ok_ := any((*N)(nil)).(*string); ok_ {
			return any(j2gort.NullString).(N)
		}
		var z_ N
		return z_
	}()
	this_.annotation = nil
	this_.priority = int32(-1)
	this_.astPosition = int32(-1)
	this_.dfaIndex = int32(-1)
	this_.inDfaQueue = false
}

// Java method: LinkedDiGraphNode.inEdgeCount()
func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) InEdgeCount() int32 {
	if this_.inEdges != nil {
		return int32(len(this_.inEdges))
	} else {
		{

			if this_.firstInEdge != nil {
				return int32(1)
			}
			return int32(0)
		}

	}
}

// Java method: LinkedDiGraphNode.outEdgeCount()
func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) OutEdgeCount() int32 {
	if this_.outEdges != nil {
		return int32(len(this_.outEdges))
	} else {
		{

			if this_.firstOutEdge != nil {
				return int32(1)
			}
			return int32(0)
		}

	}
}

// Java method: LinkedDiGraphNode.getInEdgeCount()
func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) GetInEdgeCount() int32 {
	if this_.inEdges != nil {
		return int32(len(this_.inEdges))
	} else {
		{

			if this_.firstInEdge != nil {
				return int32(1)
			}
			return int32(0)
		}

	}
}

// Java method: LinkedDiGraphNode.getInEdgeAt(int index)
func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) GetInEdgeAt(index int32) DiGraph_DiGraphEdge[N, E] {
	_ = index
	return ToIface_DiGraph_DiGraphEdge[N, E](this_.GetLinkedInEdgeAt(index))
}

// Java method: LinkedDiGraphNode.getLinkedInEdgeAt(int index)
func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) GetLinkedInEdgeAt(index int32) *LinkedDirectedGraph_LinkedDiGraphEdge[N, E] {
	_ = index
	if this_.inEdges != nil {
		return this_.inEdges[index]
	}
	if (index == int32(0)) && (this_.firstInEdge != nil) {
		return this_.firstInEdge
	}
	panic(j2gort.AsIface[javalang.IIndexOutOfBoundsException](j2gort.ConstructShimThrowable("java.lang.IndexOutOfBoundsException", j2gort.ValueOfInt(index))))
}

// Java method: LinkedDiGraphNode.getOutEdgeCount()
func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) GetOutEdgeCount() int32 {
	if this_.outEdges != nil {
		return int32(len(this_.outEdges))
	} else {
		{

			if this_.firstOutEdge != nil {
				return int32(1)
			}
			return int32(0)
		}

	}
}

// Java method: LinkedDiGraphNode.getOutEdgeAt(int index)
func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) GetOutEdgeAt(index int32) DiGraph_DiGraphEdge[N, E] {
	_ = index
	return ToIface_DiGraph_DiGraphEdge[N, E](this_.GetLinkedOutEdgeAt(index))
}

// Java method: LinkedDiGraphNode.getLinkedOutEdgeAt(int index)
func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) GetLinkedOutEdgeAt(index int32) *LinkedDirectedGraph_LinkedDiGraphEdge[N, E] {
	_ = index
	if this_.outEdges != nil {
		return this_.outEdges[index]
	}
	if (index == int32(0)) && (this_.firstOutEdge != nil) {
		return this_.firstOutEdge
	}
	panic(j2gort.AsIface[javalang.IIndexOutOfBoundsException](j2gort.ConstructShimThrowable("java.lang.IndexOutOfBoundsException", j2gort.ValueOfInt(index))))
}

// Java method: LinkedDiGraphNode.addInEdge(com.google.javascript.jscomp.graph.LinkedDirectedGraph.LinkedDiGraphEdge<N,E> edge)
func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) AddInEdge(edge *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) {
	if this_.inEdges != nil {
		j2gort.AppendPooledSlice(&this_.inEdges, edge)
	} else if this_.firstInEdge == nil {
		this_.firstInEdge = edge
	} else {
		s_ := j2gort.AcquirePooledSlice[*LinkedDirectedGraph_LinkedDiGraphEdge[N, E]](4)
		s_ = append(s_, this_.firstInEdge, edge)
		this_.inEdges = s_
		this_.firstInEdge = nil
	}

}

// Java method: LinkedDiGraphNode.addOutEdge(com.google.javascript.jscomp.graph.LinkedDirectedGraph.LinkedDiGraphEdge<N,E> edge)
func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) AddOutEdge(edge *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) {
	if this_.outEdges != nil {
		j2gort.AppendPooledSlice(&this_.outEdges, edge)
	} else if this_.firstOutEdge == nil {
		this_.firstOutEdge = edge
	} else {
		s_ := j2gort.AcquirePooledSlice[*LinkedDirectedGraph_LinkedDiGraphEdge[N, E]](4)
		s_ = append(s_, this_.firstOutEdge, edge)
		this_.outEdges = s_
		this_.firstOutEdge = nil
	}

}

// Java method: LinkedDiGraphNode.removeInEdge(com.google.javascript.jscomp.graph.LinkedDirectedGraph.LinkedDiGraphEdge<N,E> edge)
func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) RemoveInEdge(edge *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) {
	if this_.inEdges != nil {
		for i, size := int32(0), int32(len(this_.inEdges)); i < size; i++ {
			if j2gort.CheckNotNullPtr(this_.inEdges[i]) == edge {
				this_.inEdges = j2gort.SliceRemoveAt(this_.inEdges, i)
				remaining := int32(len(this_.inEdges))
				if remaining == 1 {
					this_.firstInEdge = this_.inEdges[0]
					j2gort.ReleasePooledSlice(this_.inEdges)
					this_.inEdges = nil
				} else if remaining == 0 {
					j2gort.ReleasePooledSlice(this_.inEdges)
					this_.inEdges = nil
				}
				return
			}
		}
	} else if this_.firstInEdge != nil && j2gort.CheckNotNullPtr(this_.firstInEdge) == edge {
		this_.firstInEdge = nil
	}

}

// Java method: LinkedDiGraphNode.removeOutEdge(com.google.javascript.jscomp.graph.LinkedDirectedGraph.LinkedDiGraphEdge<N,E> edge)
func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) RemoveOutEdge(edge *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) {
	if this_.outEdges != nil {
		for i, size := int32(0), int32(len(this_.outEdges)); i < size; i++ {
			if j2gort.CheckNotNullPtr(this_.outEdges[i]) == edge {
				this_.outEdges = j2gort.SliceRemoveAt(this_.outEdges, i)
				remaining := int32(len(this_.outEdges))
				if remaining == 1 {
					this_.firstOutEdge = this_.outEdges[0]
					j2gort.ReleasePooledSlice(this_.outEdges)
					this_.outEdges = nil
				} else if remaining == 0 {
					j2gort.ReleasePooledSlice(this_.outEdges)
					this_.outEdges = nil
				}
				return
			}
		}
	} else if this_.firstOutEdge != nil && j2gort.CheckNotNullPtr(this_.firstOutEdge) == edge {
		this_.firstOutEdge = nil
	}

}

// Java method: LinkedDiGraphNode.asSingletonList()
func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) AsSingletonList() j2gort.Ext_java_util_List {
	var list j2gort.Ext_java_util_List = this_.singletonSelfList
	_ = list
	if list == nil {
		list = j2gort.SingletonListPtr(this_)
		this_.singletonSelfList = list
	}
	return list
}

// Java method: LinkedDiGraphNode.getValue()
func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) GetValue() N {
	return this_.value
}

// Java method: LinkedDiGraphNode.getAnnotation()
func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) GetAnnotation() Annotation {
	if !(this_.useAnnotation) {
		panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException", "Graph initialized with node annotations turned off")))
	}
	return this_.annotation
}

// Java method: LinkedDiGraphNode.setAnnotation(com.google.javascript.jscomp.graph.Annotation data)
func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) SetAnnotation(data Annotation) {
	_ = data
	if !(this_.useAnnotation) {
		panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException", "Graph initialized with node annotations turned off")))
	}
	this_.annotation = data
}

// Java method: LinkedDiGraphNode.getColor()
func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) GetColor() string {
	return "white"
}

// Java method: LinkedDiGraphNode.getId()
func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) GetId() string {
	var id string = this_.cachedId
	_ = id
	if id == j2gort.NullString {
		id = ("LDN" + j2gort.StrValInt32(j2gort.IdentityHash(this_)))
		this_.cachedId = id
	}
	return id
}

// Java method: LinkedDiGraphNode.getLabel()
func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) GetLabel() string {
	var result string = this_.String()
	_ = result
	if this_.annotation != nil && !this_.annotation.IsJavaNil_() {
		result = (j2gort.StrValString(result) + "\n" + j2gort.StrVal(this_.annotation))
	}
	return result
}

// Java method: LinkedDiGraphNode.toString()
func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) String() string {
	return j2gort.StrVal(this_.value)
}

// Java method: LinkedDiGraphNode.getInEdges()
func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) GetInEdges() j2gort.Ext_java_util_List {
	if this_.inEdges != nil {
		return j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.ImmListFromSlice(this_.inEdges))
	}
	if this_.firstInEdge == nil {
		return j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
	} else {
		return this_.firstInEdge.AsSingletonList()
	}
}

// Java method: LinkedDiGraphNode.getOutEdges()
func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) GetOutEdges() j2gort.Ext_java_util_List {
	if this_.outEdges != nil {
		return j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.ImmListFromSlice(this_.outEdges))
	}
	if this_.firstOutEdge == nil {
		return j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
	} else {
		return this_.firstOutEdge.AsSingletonList()
	}
}

// Java method: LinkedDiGraphNode.hasPriority()
func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) HasPriority() bool {
	return (this_.priority >= int32(0))
}

// Java method: LinkedDiGraphNode.getPriority()
func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) GetPriority() int32 {
	if !(this_.priority >= int32(0)) {
		j2gort.CheckStateStr(false, "priority not set")
	}
	return this_.priority
}

// Java method: LinkedDiGraphNode.setPriority(int priority)
func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) SetPriority(priority int32) {
	_ = priority
	if !(priority >= int32(0)) {
		j2gort.CheckArgumentStr(false, "priorities must be non-negative")
	}
	this_.priority = priority
}

// Java method: LinkedDiGraphNode.getAstPosition()
func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) GetAstPosition() int32 {
	return this_.astPosition
}

// Java method: LinkedDiGraphNode.setAstPosition(int astPosition)
func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) SetAstPosition(astPosition int32) {
	_ = astPosition
	this_.astPosition = astPosition
}

// Java method: LinkedDiGraphNode.isInDfaQueue()
func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) IsInDfaQueue() bool {
	return this_.inDfaQueue
}

// Java method: LinkedDiGraphNode.setInDfaQueue(boolean inDfaQueue)
func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) SetInDfaQueue(inDfaQueue bool) {
	_ = inDfaQueue
	this_.inDfaQueue = inDfaQueue
}

// Java method: LinkedDiGraphNode.getDfaIndex()
func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) GetDfaIndex() int32 {
	return this_.dfaIndex
}

// Java method: LinkedDiGraphNode.setDfaIndex(int dfaIndex)
func (this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) SetDfaIndex(dfaIndex int32) {
	_ = dfaIndex
	this_.dfaIndex = dfaIndex
}

// Java class: com.google.javascript.jscomp.graph.LinkedDirectedGraph.LinkedDiGraphEdge
type LinkedDirectedGraph_LinkedDiGraphEdge[N comparable, E any] struct {
	sourceNode        *LinkedDirectedGraph_LinkedDiGraphNode[N, E] `gson:"sourceNode"`
	destNode          *LinkedDirectedGraph_LinkedDiGraphNode[N, E] `gson:"destNode"`
	value             E                                            `gson:"value"`
	singletonSelfList j2gort.Ext_java_util_List                    `gson:"singletonSelfList"`
	annotation        Annotation                                   `gson:"annotation"`
	dfaEdgeIndex      int32                                        `gson:"dfaEdgeIndex"`
	idHash_           atomic.Int32
	useAnnotation     bool `gson:"useAnnotation"`
}

func (this_ *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) BindSelf_(s any) {
	_ = s
}
func (this_ *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) IsJavaNil_() bool { return this_ == nil }
func (*LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.graph.LinkedDirectedGraph$LinkedDiGraphEdge", "LinkedDiGraphEdge"
}
func (this_ *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) HasIdentityEquals_() {}
func (*LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) IsGenericInst_graph_LinkedDirectedGraph_LinkedDiGraphEdge_() {
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.graph.LinkedDirectedGraph.LinkedDiGraphEdge", "com.google.javascript.jscomp.graph.LinkedDirectedGraph$LinkedDiGraphEdge", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			IsGenericInst_graph_LinkedDirectedGraph_LinkedDiGraphEdge_()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}
func (*LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) Implements_graph_DiGraph_DiGraphEdge()        {}
func (*LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) Implements_graph_GraphvizGraph_GraphvizEdge() {}
func (*LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) Implements_graph_Graph_GraphEdge()            {}
func (*LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) Implements_graph_Annotatable()                {}

func (this_ *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) initFields_LinkedDirectedGraph_LinkedDiGraphEdge() {
	_ = this_
	this_.dfaEdgeIndex = int32(-1)
}

// Java constructor: LinkedDiGraphEdge(com.google.javascript.jscomp.graph.LinkedDirectedGraph.LinkedDiGraphNode<N,E> sourceNode, E edgeValue, com.google.javascript.jscomp.graph.LinkedDirectedGraph.LinkedDiGraphNode<N,E> destNode, boolean useAnnotation)
func NewLinkedDirectedGraph_LinkedDiGraphEdge[N comparable, E any](sourceNode *LinkedDirectedGraph_LinkedDiGraphNode[N, E], edgeValue E, destNode *LinkedDirectedGraph_LinkedDiGraphNode[N, E], useAnnotation bool) *LinkedDirectedGraph_LinkedDiGraphEdge[N, E] {
	this_ := &LinkedDirectedGraph_LinkedDiGraphEdge[N, E]{}
	InitLinkedDirectedGraph_LinkedDiGraphEdge[N, E](this_, sourceNode, edgeValue, destNode, useAnnotation)
	return this_
}
func allocPooled_LinkedDirectedGraph_LinkedDiGraphEdge[N comparable, E any](slab_ *j2gort.PooledSlab[LinkedDirectedGraph_LinkedDiGraphEdge[N, E]], sourceNode *LinkedDirectedGraph_LinkedDiGraphNode[N, E], edgeValue E, destNode *LinkedDirectedGraph_LinkedDiGraphNode[N, E], useAnnotation bool) *LinkedDirectedGraph_LinkedDiGraphEdge[N, E] {
	this_ := slab_.Alloc()
	this_.dfaEdgeIndex = int32(-1)
	this_.value = edgeValue
	this_.sourceNode = sourceNode
	this_.destNode = destNode
	if useAnnotation {
		this_.useAnnotation = true
	}
	return this_
}

func InitLinkedDirectedGraph_LinkedDiGraphEdge[N comparable, E any](this_ *LinkedDirectedGraph_LinkedDiGraphEdge[N, E], sourceNode *LinkedDirectedGraph_LinkedDiGraphNode[N, E], edgeValue E, destNode *LinkedDirectedGraph_LinkedDiGraphNode[N, E], useAnnotation bool) {
	_ = sourceNode
	_ = edgeValue
	_ = destNode
	_ = useAnnotation
	this_.initFields_LinkedDirectedGraph_LinkedDiGraphEdge()
	this_.value = edgeValue
	this_.sourceNode = sourceNode
	this_.destNode = destNode
	this_.useAnnotation = useAnnotation
}

// Java method: LinkedDiGraphEdge.reinit(com.google.javascript.jscomp.graph.LinkedDirectedGraph.LinkedDiGraphNode<N,E> sourceNode, E edgeValue, com.google.javascript.jscomp.graph.LinkedDirectedGraph.LinkedDiGraphNode<N,E> destNode, boolean useAnnotation)
func (this_ *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) Reinit(sourceNode *LinkedDirectedGraph_LinkedDiGraphNode[N, E], edgeValue E, destNode *LinkedDirectedGraph_LinkedDiGraphNode[N, E], useAnnotation bool) {
	_ = sourceNode
	_ = edgeValue
	_ = destNode
	_ = useAnnotation
	this_.value = edgeValue
	this_.sourceNode = sourceNode
	this_.destNode = destNode
	this_.useAnnotation = useAnnotation
}

// Java method: LinkedDiGraphEdge.clearForPool()
func (this_ *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) ClearForPool() {
	this_.sourceNode = nil
	this_.destNode = nil
	this_.value = func() E {
		if _, ok_ := any((*E)(nil)).(*string); ok_ {
			return any(j2gort.NullString).(E)
		}
		var z_ E
		return z_
	}()
	this_.annotation = nil
	this_.dfaEdgeIndex = int32(-1)
}

// Java method: LinkedDiGraphEdge.asSingletonList()
func (this_ *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) AsSingletonList() j2gort.Ext_java_util_List {
	var list j2gort.Ext_java_util_List = this_.singletonSelfList
	_ = list
	if list == nil {
		list = j2gort.SingletonListPtr(this_)
		this_.singletonSelfList = list
	}
	return list
}

// Java method: LinkedDiGraphEdge.getSource()
func (this_ *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) GetSource() DiGraph_DiGraphNode[N, E] {
	return ToIface_DiGraph_DiGraphNode[N, E](this_.sourceNode)
}

// Java method: LinkedDiGraphEdge.getLinkedSource()
func (this_ *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) GetLinkedSource() *LinkedDirectedGraph_LinkedDiGraphNode[N, E] {
	return this_.sourceNode
}

// Java method: LinkedDiGraphEdge.getDestination()
func (this_ *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) GetDestination() DiGraph_DiGraphNode[N, E] {
	return ToIface_DiGraph_DiGraphNode[N, E](this_.destNode)
}

// Java method: LinkedDiGraphEdge.getLinkedDestination()
func (this_ *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) GetLinkedDestination() *LinkedDirectedGraph_LinkedDiGraphNode[N, E] {
	return this_.destNode
}

// Java method: LinkedDiGraphEdge.getValue()
func (this_ *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) GetValue() E {
	return this_.value
}

// Java method: LinkedDiGraphEdge.getAnnotation()
func (this_ *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) GetAnnotation() Annotation {
	if !(this_.useAnnotation) {
		panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException", "Graph initialized with edge annotations turned off")))
	}
	return this_.annotation
}

// Java method: LinkedDiGraphEdge.setAnnotation(com.google.javascript.jscomp.graph.Annotation data)
func (this_ *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) SetAnnotation(data Annotation) {
	_ = data
	if !(this_.useAnnotation) {
		panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException", "Graph initialized with edge annotations turned off")))
	}
	this_.annotation = data
}

// Java method: LinkedDiGraphEdge.getColor()
func (this_ *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) GetColor() string {
	return "black"
}

// Java method: LinkedDiGraphEdge.getLabel()
func (this_ *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) GetLabel() string {
	return j2gort.StrVal(this_.value)
}

// Java method: LinkedDiGraphEdge.getNode1Id()
func (this_ *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) GetNode1Id() string {
	return this_.sourceNode.GetId()
}

// Java method: LinkedDiGraphEdge.getNode2Id()
func (this_ *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) GetNode2Id() string {
	return this_.destNode.GetId()
}

// Java method: LinkedDiGraphEdge.toString()
func (this_ *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) String() string {
	return (j2gort.StrVal(this_.sourceNode) + " -> " + j2gort.StrVal(this_.destNode))
}

// Java method: LinkedDiGraphEdge.getNodeA()
func (this_ *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) GetNodeA() GraphNode[N, E] {
	return ToIface_GraphNode[N, E](this_.sourceNode)
}

// Java method: LinkedDiGraphEdge.getNodeB()
func (this_ *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) GetNodeB() GraphNode[N, E] {
	return ToIface_GraphNode[N, E](this_.destNode)
}

// Java method: LinkedDiGraphEdge.getDfaEdgeIndex()
func (this_ *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) GetDfaEdgeIndex() int32 {
	return this_.dfaEdgeIndex
}

// Java method: LinkedDiGraphEdge.setDfaEdgeIndex(int dfaEdgeIndex)
func (this_ *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) SetDfaEdgeIndex(dfaEdgeIndex int32) {
	_ = dfaEdgeIndex
	this_.dfaEdgeIndex = dfaEdgeIndex
}
func (this_ *LinkedDirectedGraphBase[N, E]) ReleaseSlabs_() {
	if this_ == nil {
		return
	}
	this_.slab_LinkedDirectedGraph_LinkedDiGraphEdge.Release()
	this_.slab_LinkedDirectedGraph_LinkedDiGraphNode.ReleaseWith(func(n_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) {
		if n_.inEdges != nil {
			j2gort.ReleasePooledSlice(n_.inEdges)
			n_.inEdges = nil
		}
		if n_.outEdges != nil {
			j2gort.ReleasePooledSlice(n_.outEdges)
			n_.outEdges = nil
		}
	})

	this_.nodes.ReleaseToPool()
}
