// Port of com/google/javascript/jscomp/graph/LinkedUndirectedGraph.java @ closure-compiler v20261015

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

package jscomp

import (
	"iter"
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// LinkedUndirectedGraph is an undirected graph using linked list within
// nodes to store edge information.
//
// N is the value type that the graph node stores; E is the value type that
// the graph edge stores.
//
// Go-only: language-adapter — node values are keyed with Go == (Java
// equals() for strings, numbers and identity-compared objects).
type LinkedUndirectedGraph[N comparable, E any] struct {
	UndiGraphBase[N, E]

	nodes *jutil.LinkedHashMap[N, *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]]

	useNodeAnnotations bool
	useEdgeAnnotations bool
}

// Java: LinkedUndirectedGraph.newSubGraph()
func (g *LinkedUndirectedGraph[N, E]) NewSubGraph() SubGraph[N, E] {
	return NewGraph_SimpleSubGraph[N, E](g)
}

// Java: LinkedUndirectedGraph.create()
func LinkedUndirectedGraph_Create[N comparable, E any]() *LinkedUndirectedGraph[N, E] {
	return NewLinkedUndirectedGraph[N, E](true, true)
}

// Java: LinkedUndirectedGraph(boolean, boolean)
func NewLinkedUndirectedGraph[N comparable, E any](useNodeAnnotations bool, useEdgeAnnotations bool) *LinkedUndirectedGraph[N, E] {
	g := &LinkedUndirectedGraph[N, E]{
		nodes:              jutil.NewLinkedHashMap[N, *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]](),
		useNodeAnnotations: useNodeAnnotations,
		useEdgeAnnotations: useEdgeAnnotations,
	}
	g.InitUndiGraphBase(g)
	return g
}

// getLinkedNodeOrFail is getNodeOrFail(N) with the
// LinkedUndirectedGraphNode cast.
func (g *LinkedUndirectedGraph[N, E]) getLinkedNodeOrFail(val N) *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E] {
	return g.GetNodeOrFail(val).(*LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E])
}

// Java: LinkedUndirectedGraph.connect(N, E, N)
func (g *LinkedUndirectedGraph[N, E]) Connect(srcValue N, edgeValue E, destValue N) {
	src := g.getLinkedNodeOrFail(srcValue)
	dest := g.getLinkedNodeOrFail(destValue)
	var edge *LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E]
	if g.useEdgeAnnotations {
		edge = newLinkedUndirectedGraph_AnnotatedLinkedUndirectedGraphEdge[N, E](src, edgeValue, dest)
	} else {
		edge = newLinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E](src, edgeValue, dest)
	}
	src.neighborEdges.Add(edge)
	dest.neighborEdges.Add(edge)
}

// Java: LinkedUndirectedGraph.disconnect(N, N)
func (g *LinkedUndirectedGraph[N, E]) Disconnect(srcValue N, destValue N) {
	src := g.getLinkedNodeOrFail(srcValue)
	dest := g.getLinkedNodeOrFail(destValue)
	for _, edge := range g.GetUndirectedGraphEdges(srcValue, destValue) {
		src.neighborEdges.Remove(edge)
		dest.neighborEdges.Remove(edge)
	}
}

// Java: LinkedUndirectedGraph.createUndirectedGraphNode(N)
func (g *LinkedUndirectedGraph[N, E]) CreateUndirectedGraphNode(nodeValue N) UndiGraph_UndiGraphNode[N, E] {
	return g.createLinkedNode(nodeValue)
}

func (g *LinkedUndirectedGraph[N, E]) createLinkedNode(nodeValue N) *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E] {
	return g.nodes.ComputeIfAbsent(
		nodeValue,
		func(k N) *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E] {
			if g.useNodeAnnotations {
				return newLinkedUndirectedGraph_AnnotatedLinkedUndirectedGraphNode[N, E](k)
			}
			return newLinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E](k)
		})
}

// Java: LinkedUndirectedGraph.getNeighborNodes(N)
func (g *LinkedUndirectedGraph[N, E]) GetNeighborNodes(value N) []GraphNode[N, E] {
	uNode := g.nodes.Get(value)
	if uNode == nil {
		panic(jutil.NewNullPointerException()) // ((LinkedUndirectedGraphNode) null).neighborList()
	}
	return uNode.neighborList()
}

// GetUndirectedGraphEdges returns nil (Java null) when either node is not
// in the graph.
//
// Java: LinkedUndirectedGraph.getUndirectedGraphEdges(N, N)
func (g *LinkedUndirectedGraph[N, E]) GetUndirectedGraphEdges(n1 N, n2 N) []UndiGraph_UndiGraphEdge[N, E] {
	dNode1 := g.nodes.Get(n1)
	if dNode1 == nil {
		return nil
	}
	dNode2 := g.nodes.Get(n2)
	if dNode2 == nil {
		return nil
	}
	edges := []UndiGraph_UndiGraphEdge[N, E]{}
	for _, outEdge := range dNode1.neighborEdges.Slice() {
		if outEdge.GetNodeA() == GraphNode[N, E](dNode2) || outEdge.GetNodeB() == GraphNode[N, E](dNode2) {
			edges = append(edges, outEdge)
		}
	}
	return edges
}

// Java: LinkedUndirectedGraph.getUndirectedGraphNode(N)
func (g *LinkedUndirectedGraph[N, E]) GetUndirectedGraphNode(nodeValue N) UndiGraph_UndiGraphNode[N, E] {
	if n := g.nodes.Get(nodeValue); n != nil {
		return n
	}
	return nil
}

// Java: LinkedUndirectedGraph.getUndirectedGraphNodes()
func (g *LinkedUndirectedGraph[N, E]) GetUndirectedGraphNodes() []UndiGraph_UndiGraphNode[N, E] {
	out := make([]UndiGraph_UndiGraphNode[N, E], 0, g.nodes.Size())
	for n := range g.nodes.Values() {
		out = append(out, n)
	}
	return out
}

// Java: LinkedUndirectedGraph.createNode(N)
func (g *LinkedUndirectedGraph[N, E]) CreateNode(value N) GraphNode[N, E] {
	return g.createLinkedNode(value)
}

// GetEdgesBetween panics with a NullPointerException when either node is
// not in the graph (Collections.unmodifiableList(null)).
//
// Java: LinkedUndirectedGraph.getEdges(N, N)
func (g *LinkedUndirectedGraph[N, E]) GetEdgesBetween(n1 N, n2 N) []Graph_GraphEdge[N, E] {
	edges := g.GetUndirectedGraphEdges(n1, n2)
	if edges == nil {
		panic(jutil.NewNullPointerException())
	}
	out := make([]Graph_GraphEdge[N, E], len(edges))
	for i, e := range edges {
		out[i] = e
	}
	return out
}

// Java: LinkedUndirectedGraph.getEdges()
func (g *LinkedUndirectedGraph[N, E]) GetEdges() []Graph_GraphEdge[N, E] {
	result := []Graph_GraphEdge[N, E]{}
	for node := range g.nodes.Values() {
		for _, edge := range node.neighborEdges.Slice() {
			if edge.GetNodeA() == GraphNode[N, E](node) {
				result = append(result, edge)
			}
		}
	}
	return result
}

// Java: LinkedUndirectedGraph.getFirstEdge(N, N)
func (g *LinkedUndirectedGraph[N, E]) GetFirstEdge(n1 N, n2 N) Graph_GraphEdge[N, E] {
	dNode1 := g.getLinkedNodeOrFail(n1)
	dNode2 := GraphNode[N, E](g.getLinkedNodeOrFail(n2))
	for _, outEdge := range dNode1.neighborEdges.Slice() {
		if outEdge.GetNodeA() == dNode2 || outEdge.GetNodeB() == dNode2 {
			return outEdge
		}
	}
	return nil
}

// Java: LinkedUndirectedGraph.getNode(N)
func (g *LinkedUndirectedGraph[N, E]) GetNode(value N) GraphNode[N, E] {
	if n := g.nodes.Get(value); n != nil {
		return n
	}
	return nil
}

// Java: LinkedUndirectedGraph.isConnected(N, N)
func (g *LinkedUndirectedGraph[N, E]) IsConnected(n1 N, n2 N) bool {
	return g.isConnectedWithPredicate(n1, linkedDirectedGraph_alwaysTrue[E], n2)
}

// Java: LinkedUndirectedGraph.isConnected(N, E, N)
func (g *LinkedUndirectedGraph[N, E]) IsConnectedWithEdge(n1 N, e E, n2 N) bool {
	return g.isConnectedWithPredicate(n1, linkedDirectedGraph_equalTo(e), n2)
}

// Java: LinkedUndirectedGraph.isConnected(N, Predicate<E>, N)
func (g *LinkedUndirectedGraph[N, E]) isConnectedWithPredicate(n1 N, edgePredicate func(E) bool, n2 N) bool {
	ln1 := g.nodes.Get(n1)
	if ln1 == nil {
		return false
	}
	ln2 := g.nodes.Get(n2)
	if ln2 == nil {
		return false
	}
	dNode1, dNode2 := GraphNode[N, E](ln1), GraphNode[N, E](ln2)
	for _, outEdge := range ln1.neighborEdges.Slice() {
		if (outEdge.GetNodeA() == dNode1 && outEdge.GetNodeB() == dNode2) ||
			(outEdge.GetNodeA() == dNode2 && outEdge.GetNodeB() == dNode1) {
			if edgePredicate(outEdge.GetValue()) {
				return true
			}
		}
	}
	return false
}

// Java: LinkedUndirectedGraph.getGraphvizEdges()
func (g *LinkedUndirectedGraph[N, E]) GetGraphvizEdges() []GraphvizGraph_GraphvizEdge {
	edgeList := []GraphvizGraph_GraphvizEdge{}
	for node := range g.nodes.Values() {
		for _, edge := range node.neighborEdges.Slice() {
			if edge.GetNodeA() == GraphNode[N, E](node) {
				edgeList = append(edgeList, edge.(GraphvizGraph_GraphvizEdge))
			}
		}
	}
	return edgeList
}

// Java: LinkedUndirectedGraph.getName()
func (g *LinkedUndirectedGraph[N, E]) GetName() string {
	return "LinkedUndirectedGraph"
}

// Java: LinkedUndirectedGraph.getGraphvizNodes()
func (g *LinkedUndirectedGraph[N, E]) GetGraphvizNodes() []GraphvizGraph_GraphvizNode {
	nodeList := make([]GraphvizGraph_GraphvizNode, 0, g.nodes.Size())
	for node := range g.nodes.Values() {
		nodeList = append(nodeList, node)
	}
	return nodeList
}

// Java: LinkedUndirectedGraph.isDirected()
func (g *LinkedUndirectedGraph[N, E]) IsDirected() bool {
	return false
}

// Java: LinkedUndirectedGraph.getNodes()
func (g *LinkedUndirectedGraph[N, E]) GetNodes() []GraphNode[N, E] {
	out := make([]GraphNode[N, E], 0, g.nodes.Size())
	for n := range g.nodes.Values() {
		out = append(out, n)
	}
	return out
}

// Java: LinkedUndirectedGraph.getNodeCount()
func (g *LinkedUndirectedGraph[N, E]) GetNodeCount() int {
	return g.nodes.Size()
}

// Java: LinkedUndirectedGraph.getNodeDegree(N)
func (g *LinkedUndirectedGraph[N, E]) GetNodeDegree(value N) int {
	uNode := g.nodes.Get(value)
	if uNode == nil {
		panic(jutil.NewIllegalArgumentExceptionMsg(jutil.ToJavaString(value) + " not found in graph"))
	}
	return uNode.neighborEdges.Size()
}

// LinkedUndirectedGraph_LinkedUndirectedGraphNode is an undirected graph
// node that stores outgoing edges and incoming edges as an list within the
// node itself.
//
// Go-only: language-adapter — the subclass AnnotatedLinkedUndirectedGraphNode
// only overrides the annotation accessors, so it is folded in (PORTING.md
// §4): annotated selects its behaviour.
type LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E any] struct {
	neighborEdges *jutil.ArrayList[UndiGraph_UndiGraphEdge[N, E]]
	value         N

	// AnnotatedLinkedUndirectedGraphNode
	annotated  bool
	annotation Annotation
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphNode(N)
func newLinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E any](nodeValue N) *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E] {
	return &LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]{
		neighborEdges: jutil.NewArrayList[UndiGraph_UndiGraphEdge[N, E]](),
		value:         nodeValue,
	}
}

// Java: LinkedUndirectedGraph.AnnotatedLinkedUndirectedGraphNode(N)
func newLinkedUndirectedGraph_AnnotatedLinkedUndirectedGraphNode[N, E any](nodeValue N) *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E] {
	n := newLinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E](nodeValue)
	n.annotated = true
	return n
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphNode.getNeighborEdges()
func (n *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]) GetNeighborEdges() []UndiGraph_UndiGraphEdge[N, E] {
	return n.neighborEdges.Slice()
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphNode.getNeighborEdgesIterator()
func (n *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]) GetNeighborEdgesIterator() iter.Seq[UndiGraph_UndiGraphEdge[N, E]] {
	return slices.Values(n.neighborEdges.Slice())
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphNode.getAnnotation(), AnnotatedLinkedUndirectedGraphNode.getAnnotation()
func (n *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]) GetAnnotation() Annotation {
	if n.annotated {
		return n.annotation
	}
	panic(jutil.NewUnsupportedOperationExceptionMsg(
		"Graph initialized with node annotations turned off"))
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphNode.setAnnotation(Annotation), AnnotatedLinkedUndirectedGraphNode.setAnnotation(Annotation)
func (n *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]) SetAnnotation(data Annotation) {
	if n.annotated {
		n.annotation = data
		return
	}
	panic(jutil.NewUnsupportedOperationExceptionMsg(
		"Graph initialized with node annotations turned off"))
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphNode.getValue()
func (n *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]) GetValue() N {
	return n.value
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphNode.getColor()
func (n *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]) GetColor() string {
	return "white"
}

// GetId is "LDN" + hashCode() (the identity hash code).
//
// Java: LinkedUndirectedGraph.LinkedUndirectedGraphNode.getId()
func (n *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]) GetId() string {
	return "LDN" + jutil.ToJavaString(jutil.IdentityHashCode(n))
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphNode.getLabel()
func (n *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]) GetLabel() string {
	return jutil.ToJavaString(n.value)
}

// String is the inherited Object.toString() (used by the edge's toString).
func (n *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]) String() string {
	className := "com.google.javascript.jscomp.graph.LinkedUndirectedGraph$LinkedUndirectedGraphNode"
	if n.annotated {
		className = "com.google.javascript.jscomp.graph.LinkedUndirectedGraph$AnnotatedLinkedUndirectedGraphNode"
	}
	return jutil.ObjectToString(className, n)
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphNode.neighborList()
func (n *LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]) neighborList() []GraphNode[N, E] {
	self := GraphNode[N, E](n)
	result := make([]GraphNode[N, E], 0, n.neighborEdges.Size())
	for _, edge := range n.neighborEdges.Slice() {
		if edge.GetNodeA() == self {
			result = append(result, edge.GetNodeB())
		} else {
			result = append(result, edge.GetNodeA())
		}
	}
	return result
}

// LinkedUndirectedGraph_LinkedUndirectedGraphEdge is an undirected graph
// edge that stores two nodes at each edge.
//
// Go-only: language-adapter — the subclass AnnotatedLinkedUndirectedGraphEdge
// only overrides the annotation accessors, so it is folded in (PORTING.md
// §4): annotated selects its behaviour.
type LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E any] struct {
	nodeA UndiGraph_UndiGraphNode[N, E]
	nodeB UndiGraph_UndiGraphNode[N, E]
	value E

	// AnnotatedLinkedUndirectedGraphEdge
	annotated  bool
	annotation Annotation
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphEdge(UndiGraphNode<N, E>, E, UndiGraphNode<N, E>)
func newLinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E any](nodeA UndiGraph_UndiGraphNode[N, E], edgeValue E, nodeB UndiGraph_UndiGraphNode[N, E]) *LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E] {
	return &LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E]{
		value: edgeValue,
		nodeA: nodeA,
		nodeB: nodeB,
	}
}

// Java: LinkedUndirectedGraph.AnnotatedLinkedUndirectedGraphEdge(UndiGraphNode<N, E>, E, UndiGraphNode<N, E>)
func newLinkedUndirectedGraph_AnnotatedLinkedUndirectedGraphEdge[N, E any](nodeA UndiGraph_UndiGraphNode[N, E], edgeValue E, nodeB UndiGraph_UndiGraphNode[N, E]) *LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E] {
	e := newLinkedUndirectedGraph_LinkedUndirectedGraphEdge(nodeA, edgeValue, nodeB)
	e.annotated = true
	return e
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphEdge.getValue()
func (e *LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E]) GetValue() E {
	return e.value
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphEdge.getNodeA()
func (e *LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E]) GetNodeA() GraphNode[N, E] {
	return e.nodeA
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphEdge.getNodeB()
func (e *LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E]) GetNodeB() GraphNode[N, E] {
	return e.nodeB
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphEdge.getAnnotation(), AnnotatedLinkedUndirectedGraphEdge.getAnnotation()
func (e *LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E]) GetAnnotation() Annotation {
	if e.annotated {
		return e.annotation
	}
	panic(jutil.NewUnsupportedOperationExceptionMsg(
		"Graph initialized with edge annotations turned off"))
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphEdge.setAnnotation(Annotation), AnnotatedLinkedUndirectedGraphEdge.setAnnotation(Annotation)
func (e *LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E]) SetAnnotation(data Annotation) {
	if e.annotated {
		e.annotation = data
		return
	}
	panic(jutil.NewUnsupportedOperationExceptionMsg(
		"Graph initialized with edge annotations turned off"))
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphEdge.getColor()
func (e *LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E]) GetColor() string {
	return "black"
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphEdge.getLabel()
func (e *LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E]) GetLabel() string {
	return jutil.ToJavaString(e.value)
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphEdge.getNode1Id()
func (e *LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E]) GetNode1Id() string {
	return e.nodeA.(*LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]).GetId()
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphEdge.getNode2Id()
func (e *LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E]) GetNode2Id() string {
	return e.nodeB.(*LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E]).GetId()
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphEdge.toString()
func (e *LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E]) String() string {
	return jutil.ToJavaString(e.nodeA) + " -- " + jutil.ToJavaString(e.nodeB)
}

var (
	_ UndiGraph[string, string]               = (*LinkedUndirectedGraph[string, string])(nil)
	_ GraphvizGraph                           = (*LinkedUndirectedGraph[string, string])(nil)
	_ UndiGraph_UndiGraphNode[string, string] = (*LinkedUndirectedGraph_LinkedUndirectedGraphNode[string, string])(nil)
	_ GraphvizGraph_GraphvizNode              = (*LinkedUndirectedGraph_LinkedUndirectedGraphNode[string, string])(nil)
	_ UndiGraph_UndiGraphEdge[string, string] = (*LinkedUndirectedGraph_LinkedUndirectedGraphEdge[string, string])(nil)
	_ GraphvizGraph_GraphvizEdge              = (*LinkedUndirectedGraph_LinkedUndirectedGraphEdge[string, string])(nil)
)
