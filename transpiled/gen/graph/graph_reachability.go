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
// Port of com/google/javascript/jscomp/graph/GraphReachability.java (generated mechanically by j2go).

package graph

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.graph.GraphReachability
type GraphReachability[N comparable, E any] struct {
	graph         DiGraph[N, E]                               `gson:"graph"`
	edgePredicate j2gort.Ext_com_google_common_base_Predicate `gson:"edgePredicate"`
	edgeFilter    GraphReachability_EdgeFilter                `gson:"edgeFilter"`
	idHash_       atomic.Int32
}

func (this_ *GraphReachability[N, E]) BindSelf_(s any) {
	_ = s
}
func (this_ *GraphReachability[N, E]) IsJavaNil_() bool { return this_ == nil }
func (*GraphReachability[N, E]) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.graph.GraphReachability", "GraphReachability"
}
func (this_ *GraphReachability[N, E]) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*GraphReachability[N, E]) HasIdentityEquals_()                     {}
func (*GraphReachability[N, E]) IsGenericInst_graph_GraphReachability_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.graph.GraphReachability", "com.google.javascript.jscomp.graph.GraphReachability", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			IsGenericInst_graph_GraphReachability_()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}
func (*GraphReachability[N, E]) Implements_graph_FixedPointGraphTraversal_EdgeCallback() {}

var static_GraphReachability_REACHABLE Annotation

func StaticGet_GraphReachability_REACHABLE() Annotation {
	clinit_GraphReachability()
	return static_GraphReachability_REACHABLE
}
func StaticSet_GraphReachability_REACHABLE(v Annotation) {
	clinit_GraphReachability()
	static_GraphReachability_REACHABLE = v
}

var clinitDone_GraphReachability uint32
var clinitState_GraphReachability j2gort.ClinitState

func clinit_GraphReachability() {
	if atomic.LoadUint32(&clinitDone_GraphReachability) != 0 {
		return
	}
	clinitSlow_GraphReachability()
}

func clinitSlow_GraphReachability() {
	if !j2gort.EnterClinit(&clinitState_GraphReachability, &clinitDone_GraphReachability) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_GraphReachability, &clinitDone_GraphReachability)
	static_GraphReachability_REACHABLE = NewGraphReachability__Anon_118()
}

func (this_ *GraphReachability[N, E]) initFields_GraphReachability() {
	_ = this_
}

// Java constructor: GraphReachability(com.google.javascript.jscomp.graph.DiGraph<N,E> graph)
func NewGraphReachability_1[N comparable, E any](graph DiGraph[N, E]) *GraphReachability[N, E] {
	clinit_GraphReachability()
	this_ := &GraphReachability[N, E]{}
	InitGraphReachability_1[N, E](this_, graph)
	return this_
}

func InitGraphReachability_1[N comparable, E any](this_ *GraphReachability[N, E], graph DiGraph[N, E]) {
	_ = graph
	this_.initFields_GraphReachability()
	this_.graph = graph
	this_.edgePredicate = nil
	this_.edgeFilter = nil
}

// Java constructor: GraphReachability(com.google.javascript.jscomp.graph.DiGraph<N,E> graph, com.google.common.base.@org.jspecify.annotations.Nullable Predicate<com.google.javascript.jscomp.graph.GraphReachability.EdgeTuple<N,E>> edgePredicate)
func NewGraphReachability_2_DiGraph_any_any__Ext_com_google_common_base_Predicate[N comparable, E any](graph DiGraph[N, E], edgePredicate j2gort.Ext_com_google_common_base_Predicate) *GraphReachability[N, E] {
	clinit_GraphReachability()
	this_ := &GraphReachability[N, E]{}
	InitGraphReachability_2_DiGraph_any_any__Ext_com_google_common_base_Predicate[N, E](this_, graph, edgePredicate)
	return this_
}

func InitGraphReachability_2_DiGraph_any_any__Ext_com_google_common_base_Predicate[N comparable, E any](this_ *GraphReachability[N, E], graph DiGraph[N, E], edgePredicate j2gort.Ext_com_google_common_base_Predicate) {
	_ = graph
	_ = edgePredicate
	this_.initFields_GraphReachability()
	this_.graph = graph
	this_.edgePredicate = edgePredicate
	this_.edgeFilter = nil
}

// Java constructor: GraphReachability(com.google.javascript.jscomp.graph.DiGraph<N,E> graph, com.google.javascript.jscomp.graph.GraphReachability.@org.jspecify.annotations.Nullable EdgeFilter<N,E> edgeFilter)
func NewGraphReachability_2_DiGraph_any_any__GraphReachability_EdgeFilter[N comparable, E any](graph DiGraph[N, E], edgeFilter GraphReachability_EdgeFilter) *GraphReachability[N, E] {
	clinit_GraphReachability()
	this_ := &GraphReachability[N, E]{}
	InitGraphReachability_2_DiGraph_any_any__GraphReachability_EdgeFilter[N, E](this_, graph, edgeFilter)
	return this_
}

func InitGraphReachability_2_DiGraph_any_any__GraphReachability_EdgeFilter[N comparable, E any](this_ *GraphReachability[N, E], graph DiGraph[N, E], edgeFilter GraphReachability_EdgeFilter) {
	_ = graph
	_ = edgeFilter
	this_.initFields_GraphReachability()
	this_.graph = graph
	this_.edgePredicate = nil
	this_.edgeFilter = edgeFilter
}

// Java method: GraphReachability.compute(N entry)
func (this_ *GraphReachability[N, E]) Compute(entry N) {
	_ = entry
	this_.graph.ClearNodeAnnotations()
	var entryNode DiGraph_DiGraphNode[N, E] = j2gort.AsIface[DiGraph_DiGraphNode[N, E]](this_.graph.GetNode(entry))
	_ = entryNode
	entryNode.SetAnnotation(StaticGet_GraphReachability_REACHABLE())
	this_.priv_GraphReachability_propagate_6315(entryNode, j2gort.NewExt_java_util_ArrayDeque())
}

// Java method: GraphReachability.recompute(N reachableNode)
func (this_ *GraphReachability[N, E]) Recompute(reachableNode N) {
	_ = reachableNode
	var newReachable DiGraph_DiGraphNode[N, E] = j2gort.AsIface[DiGraph_DiGraphNode[N, E]](this_.graph.GetNode(reachableNode))
	_ = newReachable
	if !(!(j2gort.IfaceEquals[Annotation](newReachable.GetAnnotation(), StaticGet_GraphReachability_REACHABLE()))) {
		j2gort.CheckState(false)
	}
	newReachable.SetAnnotation(StaticGet_GraphReachability_REACHABLE())
	this_.priv_GraphReachability_propagate_6315(newReachable, j2gort.NewExt_java_util_ArrayDeque())
}

// Java static method: GraphReachability.computeReachability(com.google.javascript.jscomp.graph.DiGraph<N,E> graph, N entry, com.google.javascript.jscomp.graph.GraphReachability.@org.jspecify.annotations.Nullable EdgeFilter<N,E> edgeFilter, java.util.ArrayDeque<com.google.javascript.jscomp.graph.DiGraph.DiGraphNode<N,E>> workQueue)
func GraphReachability_ComputeReachability[N comparable, E any](graph DiGraph[N, E], entry N, edgeFilter GraphReachability_EdgeFilter, workQueue *j2gort.Ext_java_util_ArrayDeque) {
	clinit_GraphReachability()
	_ = graph
	_ = entry
	_ = edgeFilter
	_ = workQueue
	graph.ClearNodeAnnotations()
	var entryNode DiGraph_DiGraphNode[N, E] = j2gort.AsIface[DiGraph_DiGraphNode[N, E]](graph.GetNode(entry))
	_ = entryNode
	entryNode.SetAnnotation(StaticGet_GraphReachability_REACHABLE())
	GraphReachability_PropagateWithFilter[N, E](entryNode, edgeFilter, workQueue)
}

// Java static method: GraphReachability.recomputeReachability(com.google.javascript.jscomp.graph.DiGraph<N,E> graph, N reachableNode, java.util.ArrayDeque<com.google.javascript.jscomp.graph.DiGraph.DiGraphNode<N,E>> workQueue)
func GraphReachability_RecomputeReachability[N comparable, E any](graph DiGraph[N, E], reachableNode N, workQueue *j2gort.Ext_java_util_ArrayDeque) {
	clinit_GraphReachability()
	_ = graph
	_ = reachableNode
	_ = workQueue
	var newReachable DiGraph_DiGraphNode[N, E] = j2gort.AsIface[DiGraph_DiGraphNode[N, E]](graph.GetNode(reachableNode))
	_ = newReachable
	if !(!(j2gort.IfaceEquals[Annotation](newReachable.GetAnnotation(), StaticGet_GraphReachability_REACHABLE()))) {
		j2gort.CheckState(false)
	}
	newReachable.SetAnnotation(StaticGet_GraphReachability_REACHABLE())
	GraphReachability_PropagateWithFilter[N, E](newReachable, nil, workQueue)
}

// Java method: GraphReachability.propagate(com.google.javascript.jscomp.graph.DiGraph.DiGraphNode<N,E> startNode, java.util.ArrayDeque<com.google.javascript.jscomp.graph.DiGraph.DiGraphNode<N,E>> workQueue)
func (this_ *GraphReachability[N, E]) priv_GraphReachability_propagate_6315(startNode DiGraph_DiGraphNode[N, E], workQueue *j2gort.Ext_java_util_ArrayDeque) {
	_ = startNode
	_ = workQueue
	if this_.edgePredicate == nil {
		GraphReachability_PropagateWithFilter[N, E](startNode, this_.edgeFilter, workQueue)
		return
	}
	workQueue.Clear()
	workQueue.AddLast(startNode)
	for !(workQueue.IsEmpty()) {
		var srcNode DiGraph_DiGraphNode[N, E] = j2gort.AsIface[DiGraph_DiGraphNode[N, E]](workQueue.RemoveFirst())
		_ = srcNode
		var srcValue N = srcNode.GetValue()
		_ = srcValue
		{
			var i int32 = int32(0)
			_ = i
			var outCount int32 = srcNode.GetOutEdgeCount()
			_ = outCount
			for ; i < outCount; i = i + 1 {
				var edge DiGraph_DiGraphEdge[N, E] = srcNode.GetOutEdgeAt(i)
				_ = edge
				var destNode DiGraph_DiGraphNode[N, E] = edge.GetDestination()
				_ = destNode
				if !(j2gort.IfaceEquals[Annotation](destNode.GetAnnotation(), StaticGet_GraphReachability_REACHABLE())) && this_.edgePredicate.Apply(NewGraphReachability_EdgeTuple[N, E](srcValue, edge.GetValue(), destNode.GetValue())) {
					destNode.SetAnnotation(StaticGet_GraphReachability_REACHABLE())
					workQueue.AddLast(destNode)
				}
			}
		}
	}
}

// Java static method: GraphReachability.propagateWithFilter(com.google.javascript.jscomp.graph.DiGraph.DiGraphNode<N,E> startNode, com.google.javascript.jscomp.graph.GraphReachability.@org.jspecify.annotations.Nullable EdgeFilter<N,E> edgeFilter, java.util.ArrayDeque<com.google.javascript.jscomp.graph.DiGraph.DiGraphNode<N,E>> workQueue)
func GraphReachability_PropagateWithFilter[N comparable, E any](startNode DiGraph_DiGraphNode[N, E], edgeFilter GraphReachability_EdgeFilter, workQueue *j2gort.Ext_java_util_ArrayDeque) {
	clinit_GraphReachability()
	_ = startNode
	_ = edgeFilter
	_ = workQueue
	workQueue.Clear()
	workQueue.AddLast(startNode)
	if edgeFilter == nil || edgeFilter.IsJavaNil_() {
		for !(workQueue.IsEmpty()) {
			var srcNode DiGraph_DiGraphNode[N, E] = j2gort.AsIface[DiGraph_DiGraphNode[N, E]](workQueue.RemoveFirst())
			_ = srcNode
			{
				var i int32 = int32(0)
				_ = i
				var outCount int32 = srcNode.GetOutEdgeCount()
				_ = outCount
				for ; i < outCount; i = i + 1 {
					var destNode DiGraph_DiGraphNode[N, E] = srcNode.GetOutEdgeAt(i).GetDestination()
					_ = destNode
					if !(j2gort.IfaceEquals[Annotation](destNode.GetAnnotation(), StaticGet_GraphReachability_REACHABLE())) {
						destNode.SetAnnotation(StaticGet_GraphReachability_REACHABLE())
						workQueue.AddLast(destNode)
					}
				}
			}
		}
	} else {
		for !(workQueue.IsEmpty()) {
			var srcNode_2 DiGraph_DiGraphNode[N, E] = j2gort.AsIface[DiGraph_DiGraphNode[N, E]](workQueue.RemoveFirst())
			_ = srcNode_2
			var srcValue N = srcNode_2.GetValue()
			_ = srcValue
			{
				var i_2 int32 = int32(0)
				_ = i_2
				var outCount_2 int32 = srcNode_2.GetOutEdgeCount()
				_ = outCount_2
				for ; i_2 < outCount_2; i_2 = i_2 + 1 {
					var edge DiGraph_DiGraphEdge[N, E] = srcNode_2.GetOutEdgeAt(i_2)
					_ = edge
					var destNode_2 DiGraph_DiGraphNode[N, E] = edge.GetDestination()
					_ = destNode_2
					if !(j2gort.IfaceEquals[Annotation](destNode_2.GetAnnotation(), StaticGet_GraphReachability_REACHABLE())) && edgeFilter.Apply(j2gort.ToAny(srcValue), j2gort.ToAny(edge.GetValue()), j2gort.ToAny(destNode_2.GetValue())) {
						destNode_2.SetAnnotation(StaticGet_GraphReachability_REACHABLE())
						workQueue.AddLast(destNode_2)
					}
				}
			}
		}
	}
}

// Java method: GraphReachability.traverseEdge(N source, E e, N destination)
func (this_ *GraphReachability[N, E]) TraverseEdge(source N, e E, destination N) bool {
	_ = source
	_ = e
	_ = destination
	if (j2gort.IfaceEquals[Annotation](j2gort.AsIface[DiGraph_DiGraphNode[N, E]](this_.graph.GetNode(source)).GetAnnotation(), StaticGet_GraphReachability_REACHABLE()) && ((this_.edgeFilter == nil || this_.edgeFilter.IsJavaNil_()) || this_.edgeFilter.Apply(j2gort.ToAny(source), j2gort.ToAny(e), j2gort.ToAny(destination)))) && ((this_.edgePredicate == nil) || this_.edgePredicate.Apply(NewGraphReachability_EdgeTuple[N, E](source, e, destination))) {
		var destNode GraphNode[N, E] = j2gort.AsIface[DiGraph_DiGraphNode[N, E]](this_.graph.GetNode(destination))
		_ = destNode
		if !(j2gort.IfaceEquals[Annotation](destNode.GetAnnotation(), StaticGet_GraphReachability_REACHABLE())) {
			destNode.SetAnnotation(StaticGet_GraphReachability_REACHABLE())
			return true
		}
	}
	return false
}

// Java interface: com.google.javascript.jscomp.graph.GraphReachability.EdgeFilter
type GraphReachability_EdgeFilter interface {
	Implements_graph_GraphReachability_EdgeFilter()
	IsJavaNil_() bool
	Apply(any, any, any) bool
}

func ToIface_GraphReachability_EdgeFilter[T interface {
	GraphReachability_EdgeFilter
	comparable
}](v T) GraphReachability_EdgeFilter {
	var z T
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.graph.GraphReachability.EdgeFilter", "com.google.javascript.jscomp.graph.GraphReachability$EdgeFilter", func(v_ any) bool {
		return j2gort.InstanceOf[GraphReachability_EdgeFilter](v_)
	})
}

// Java class: com.google.javascript.jscomp.graph.GraphReachability.EdgeTuple
type GraphReachability_EdgeTuple[N comparable, E any] struct {
	sourceNode N `gson:"sourceNode"`
	edge       E `gson:"edge"`
	idHash_    atomic.Int32
}

func (this_ *GraphReachability_EdgeTuple[N, E]) BindSelf_(s any) {
	_ = s
}
func (this_ *GraphReachability_EdgeTuple[N, E]) IsJavaNil_() bool { return this_ == nil }
func (*GraphReachability_EdgeTuple[N, E]) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.graph.GraphReachability$EdgeTuple", "EdgeTuple"
}
func (this_ *GraphReachability_EdgeTuple[N, E]) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*GraphReachability_EdgeTuple[N, E]) HasIdentityEquals_()                               {}
func (*GraphReachability_EdgeTuple[N, E]) IsGenericInst_graph_GraphReachability_EdgeTuple_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.graph.GraphReachability.EdgeTuple", "com.google.javascript.jscomp.graph.GraphReachability$EdgeTuple", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			IsGenericInst_graph_GraphReachability_EdgeTuple_()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}

func (this_ *GraphReachability_EdgeTuple[N, E]) initFields_GraphReachability_EdgeTuple() {
	_ = this_
}

// Java constructor: EdgeTuple(N sourceNode, E edge, N destNode)
func NewGraphReachability_EdgeTuple[N comparable, E any](sourceNode N, edge E, destNode N) *GraphReachability_EdgeTuple[N, E] {
	this_ := &GraphReachability_EdgeTuple[N, E]{}
	InitGraphReachability_EdgeTuple[N, E](this_, sourceNode, edge, destNode)
	return this_
}

func InitGraphReachability_EdgeTuple[N comparable, E any](this_ *GraphReachability_EdgeTuple[N, E], sourceNode N, edge E, destNode N) {
	_ = sourceNode
	_ = edge
	_ = destNode
	this_.initFields_GraphReachability_EdgeTuple()
	this_.sourceNode = sourceNode
	this_.edge = edge
}

// Java class:
type GraphReachability__Anon_118 struct {
	idHash_ atomic.Int32
}

func (this_ *GraphReachability__Anon_118) BindSelf_(s any) {
	_ = s
}
func (this_ *GraphReachability__Anon_118) IsJavaNil_() bool { return this_ == nil }
func (*GraphReachability__Anon_118) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.graph.GraphReachability$1", ""
}
func (this_ *GraphReachability__Anon_118) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*GraphReachability__Anon_118) HasIdentityEquals_()          {}
func (*GraphReachability__Anon_118) Implements_graph_Annotation() {}

func (this_ *GraphReachability__Anon_118) initFields_GraphReachability__Anon_118() {
	_ = this_
}

func NewGraphReachability__Anon_118() *GraphReachability__Anon_118 {
	this_ := &GraphReachability__Anon_118{}
	this_.initFields_GraphReachability__Anon_118()
	return this_
}
