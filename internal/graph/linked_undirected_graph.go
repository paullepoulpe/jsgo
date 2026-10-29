// Port of com/google/javascript/jscomp/graph/LinkedUndirectedGraph.java @ closure-compiler v20261021

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

package graph

import (
	"iter"
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// LinkedUndirectedGraph is an undirected graph storing incident edges within
// slab-allocated nodes.
//
// N is the value type that the graph node stores; E is the value type that
// the graph edge stores.
type LinkedUndirectedGraph[N comparable, E any] struct {
	Nodes jutil.LinkedHashMap[N, *LinkedUndirectedGraphNode[N, E]]

	// Go-only: performance — see jutil.ReuseSlab.
	nodeSlab jutil.ReuseSlab[LinkedUndirectedGraphNode[N, E]]
	edgeSlab jutil.ReuseSlab[LinkedUndirectedGraphEdge[N, E]]

	nodeAnnotationStack *annotationStack
	edgeAnnotationStack *annotationStack

	useNodeAnnotations bool
	useEdgeAnnotations bool
}

// Type aliases preserving callers that spell out the nested Java names.
type (
	LinkedUndirectedGraph_LinkedUndirectedGraphNode[N, E any] = LinkedUndirectedGraphNode[N, E]
	LinkedUndirectedGraph_LinkedUndirectedGraphEdge[N, E any] = LinkedUndirectedGraphEdge[N, E]
	UndiGraph_UndiGraphNode[N, E any]                         = *LinkedUndirectedGraphNode[N, E]
	UndiGraph_UndiGraphEdge[N, E any]                         = *LinkedUndirectedGraphEdge[N, E]
)

type linkedUndirectedSubGraph[N comparable, E any] struct {
	graph *LinkedUndirectedGraph[N, E]
	nodes []*LinkedUndirectedGraphNode[N, E]
	buf   [8]*LinkedUndirectedGraphNode[N, E]
}

// Java: LinkedUndirectedGraph.newSubGraph()
func (g *LinkedUndirectedGraph[N, E]) NewSubGraph() SubGraph[N, E] {
	s := &linkedUndirectedSubGraph[N, E]{graph: g}
	s.nodes = s.buf[:0]
	return s
}

func (s *linkedUndirectedSubGraph[N, E]) IsIndependentOf(value N) bool {
	node := s.graph.Nodes.Get(value)
	if node == nil || len(node.NeighborEdges) == 0 {
		return true
	}
	for _, n := range s.nodes {
		if s.graph.IsConnectedNodesAny(n, node) {
			return false
		}
	}
	return true
}

func (s *linkedUndirectedSubGraph[N, E]) AddNode(value N) {
	s.nodes = append(s.nodes, s.graph.getLinkedNodeOrFail(value))
}

// Java: LinkedUndirectedGraph.create()
func LinkedUndirectedGraph_Create[N comparable, E any]() *LinkedUndirectedGraph[N, E] {
	return NewLinkedUndirectedGraph[N, E](true, true)
}

// Java: LinkedUndirectedGraph(boolean, boolean)
func NewLinkedUndirectedGraph[N comparable, E any](useNodeAnnotations bool, useEdgeAnnotations bool) *LinkedUndirectedGraph[N, E] {
	g := &LinkedUndirectedGraph[N, E]{}
	g.Init(useNodeAnnotations, useEdgeAnnotations)
	return g
}

// Init configures g's annotation flags.
func (g *LinkedUndirectedGraph[N, E]) Init(useNodeAnnotations bool, useEdgeAnnotations bool) {
	g.useNodeAnnotations = useNodeAnnotations
	g.useEdgeAnnotations = useEdgeAnnotations
}

// ClearForReuse removes all nodes and edges from g, keeping its slab chunks
// and (when small) its node map for the next graph built in g.
func (g *LinkedUndirectedGraph[N, E]) ClearForReuse() {
	g.nodeAnnotationStack = nil
	g.edgeAnnotationStack = nil
	if g.Nodes.Size() > linkedDirectedGraphReuseMaxNodes {
		g.Nodes.Reset()
	} else {
		g.Nodes.Clear()
	}
	g.nodeSlab.Reset((*LinkedUndirectedGraphNode[N, E]).recycle)
	g.edgeSlab.Reset(nil)
}

func (g *LinkedUndirectedGraph[N, E]) getLinkedNodeOrFail(val N) *LinkedUndirectedGraphNode[N, E] {
	if n := g.Nodes.Get(val); n != nil {
		return n
	}
	panic(jutil.NewIllegalArgumentExceptionMsg(jutil.ToJavaString(val) + " does not exist in graph"))
}

// Java: Graph.getNodeOrFail(N)
func (g *LinkedUndirectedGraph[N, E]) GetNodeOrFail(val N) *LinkedUndirectedGraphNode[N, E] {
	return g.getLinkedNodeOrFail(val)
}

func (g *LinkedUndirectedGraph[N, E]) allocNode(nodeValue N) *LinkedUndirectedGraphNode[N, E] {
	n := g.nodeSlab.Alloc()
	n.Value = nodeValue
	n.annotated = g.useNodeAnnotations
	if cap(n.NeighborEdges) == 0 {
		n.NeighborEdges = n.neighborEdgeBuffer[:0:len(n.neighborEdgeBuffer)]
	}
	return n
}

func (g *LinkedUndirectedGraph[N, E]) newEdge(
	src *LinkedUndirectedGraphNode[N, E], edgeValue E, dest *LinkedUndirectedGraphNode[N, E]) *LinkedUndirectedGraphEdge[N, E] {
	e := g.edgeSlab.Alloc()
	e.nodeA = src
	e.nodeB = dest
	e.value = edgeValue
	e.annotated = g.useEdgeAnnotations
	return e
}

// ConnectLinkedNodes connects two already-resolved linked nodes directly.
func (g *LinkedUndirectedGraph[N, E]) ConnectLinkedNodes(
	src *LinkedUndirectedGraphNode[N, E], edgeValue E, dest *LinkedUndirectedGraphNode[N, E]) {
	edge := g.newEdge(src, edgeValue, dest)
	src.NeighborEdges = append(src.NeighborEdges, edge)
	dest.NeighborEdges = append(dest.NeighborEdges, edge)
}

// Java: LinkedUndirectedGraph.connect(N, E, N)
func (g *LinkedUndirectedGraph[N, E]) Connect(srcValue N, edgeValue E, destValue N) {
	src := g.getLinkedNodeOrFail(srcValue)
	dest := g.getLinkedNodeOrFail(destValue)
	g.ConnectLinkedNodes(src, edgeValue, dest)
}

// Java: Graph.connectIfNotFound(N, E, N)
func (g *LinkedUndirectedGraph[N, E]) ConnectIfNotFound(n1 N, edge E, n2 N) {
	src := g.getLinkedNodeOrFail(n1)
	dest := g.getLinkedNodeOrFail(n2)
	if !g.isConnectedNodesWithEdge(src, edge, dest) {
		g.ConnectLinkedNodes(src, edge, dest)
	}
}

// Java: LinkedUndirectedGraph.disconnect(N, N)
func (g *LinkedUndirectedGraph[N, E]) Disconnect(srcValue N, destValue N) {
	src := g.getLinkedNodeOrFail(srcValue)
	dest := g.getLinkedNodeOrFail(destValue)
	for _, edge := range g.GetUndirectedGraphEdges(srcValue, destValue) {
		if i := slices.Index(src.NeighborEdges, edge); i >= 0 {
			src.NeighborEdges = slices.Delete(src.NeighborEdges, i, i+1)
		}
		if i := slices.Index(dest.NeighborEdges, edge); i >= 0 {
			dest.NeighborEdges = slices.Delete(dest.NeighborEdges, i, i+1)
		}
	}
}

// Java: LinkedUndirectedGraph.createUndirectedGraphNode(N)
func (g *LinkedUndirectedGraph[N, E]) CreateUndirectedGraphNode(nodeValue N) *LinkedUndirectedGraphNode[N, E] {
	return g.CreateLinkedNode(nodeValue)
}

// CreateLinkedNode returns or creates the concrete linked undirected graph node for nodeValue.
func (g *LinkedUndirectedGraph[N, E]) CreateLinkedNode(nodeValue N) *LinkedUndirectedGraphNode[N, E] {
	if n := g.Nodes.Get(nodeValue); n != nil {
		return n
	}
	n := g.allocNode(nodeValue)
	g.Nodes.PutNew(nodeValue, n)
	return n
}

// Java: LinkedUndirectedGraph.getNeighborNodes(N)
func (g *LinkedUndirectedGraph[N, E]) GetNeighborNodes(value N) []*LinkedUndirectedGraphNode[N, E] {
	uNode := g.Nodes.Get(value)
	if uNode == nil {
		panic(jutil.NewNullPointerException()) // ((LinkedUndirectedGraphNode) null).neighborList()
	}
	return uNode.neighborList()
}

// GetUndirectedGraphEdges returns nil (Java null) when either node is not
// in the graph.
//
// Java: LinkedUndirectedGraph.getUndirectedGraphEdges(N, N)
func (g *LinkedUndirectedGraph[N, E]) GetUndirectedGraphEdges(n1 N, n2 N) []*LinkedUndirectedGraphEdge[N, E] {
	dNode1 := g.Nodes.Get(n1)
	if dNode1 == nil {
		return nil
	}
	dNode2 := g.Nodes.Get(n2)
	if dNode2 == nil {
		return nil
	}
	edges := []*LinkedUndirectedGraphEdge[N, E]{}
	for _, outEdge := range dNode1.NeighborEdges {
		if outEdge.nodeA == dNode2 || outEdge.nodeB == dNode2 {
			edges = append(edges, outEdge)
		}
	}
	return edges
}

// Java: LinkedUndirectedGraph.getUndirectedGraphNode(N)
func (g *LinkedUndirectedGraph[N, E]) GetUndirectedGraphNode(nodeValue N) *LinkedUndirectedGraphNode[N, E] {
	return g.Nodes.Get(nodeValue)
}

// Java: LinkedUndirectedGraph.getUndirectedGraphNodes()
func (g *LinkedUndirectedGraph[N, E]) GetUndirectedGraphNodes() []*LinkedUndirectedGraphNode[N, E] {
	return g.GetNodes()
}

// Java: LinkedUndirectedGraph.createNode(N)
func (g *LinkedUndirectedGraph[N, E]) CreateNode(value N) *LinkedUndirectedGraphNode[N, E] {
	return g.CreateLinkedNode(value)
}

// GetEdgesBetween panics with a NullPointerException when either node is
// not in the graph (Collections.unmodifiableList(null)).
//
// Java: LinkedUndirectedGraph.getEdges(N, N)
func (g *LinkedUndirectedGraph[N, E]) GetEdgesBetween(n1 N, n2 N) []*LinkedUndirectedGraphEdge[N, E] {
	edges := g.GetUndirectedGraphEdges(n1, n2)
	if edges == nil {
		panic(jutil.NewNullPointerException())
	}
	return edges
}

// Java: LinkedUndirectedGraph.getEdges()
func (g *LinkedUndirectedGraph[N, E]) GetEdges() []*LinkedUndirectedGraphEdge[N, E] {
	result := []*LinkedUndirectedGraphEdge[N, E]{}
	for node := range g.Nodes.Values() {
		for _, edge := range node.NeighborEdges {
			if edge.nodeA == node {
				result = append(result, edge)
			}
		}
	}
	return result
}

// Java: LinkedUndirectedGraph.getFirstEdge(N, N)
func (g *LinkedUndirectedGraph[N, E]) GetFirstEdge(n1 N, n2 N) *LinkedUndirectedGraphEdge[N, E] {
	dNode1 := g.getLinkedNodeOrFail(n1)
	dNode2 := g.getLinkedNodeOrFail(n2)
	for _, outEdge := range dNode1.NeighborEdges {
		if outEdge.nodeA == dNode2 || outEdge.nodeB == dNode2 {
			return outEdge
		}
	}
	return nil
}

// Java: Graph.hasNode(N)
func (g *LinkedUndirectedGraph[N, E]) HasNode(n N) bool {
	return g.Nodes.Get(n) != nil
}

// Java: LinkedUndirectedGraph.getNode(N)
func (g *LinkedUndirectedGraph[N, E]) GetNode(value N) *LinkedUndirectedGraphNode[N, E] {
	return g.Nodes.Get(value)
}

// Java: LinkedUndirectedGraph.isConnected(N, N)
func (g *LinkedUndirectedGraph[N, E]) IsConnected(n1 N, n2 N) bool {
	ln1 := g.Nodes.Get(n1)
	if ln1 == nil {
		return false
	}
	ln2 := g.Nodes.Get(n2)
	if ln2 == nil {
		return false
	}
	return g.IsConnectedNodesAny(ln1, ln2)
}

// Java: LinkedUndirectedGraph.isConnected(N, E, N)
func (g *LinkedUndirectedGraph[N, E]) IsConnectedWithEdge(n1 N, e E, n2 N) bool {
	ln1 := g.Nodes.Get(n1)
	if ln1 == nil {
		return false
	}
	ln2 := g.Nodes.Get(n2)
	if ln2 == nil {
		return false
	}
	return g.isConnectedNodesWithEdge(ln1, e, ln2)
}

// IsConnectedNodesAny reports whether ln1 and ln2 share any edge.
func (g *LinkedUndirectedGraph[N, E]) IsConnectedNodesAny(
	ln1 *LinkedUndirectedGraphNode[N, E], ln2 *LinkedUndirectedGraphNode[N, E]) bool {
	edges := ln1.NeighborEdges
	if edges2 := ln2.NeighborEdges; len(edges2) < len(edges) {
		edges = edges2
	}
	for _, outEdge := range edges {
		if (outEdge.nodeA == ln1 && outEdge.nodeB == ln2) ||
			(outEdge.nodeA == ln2 && outEdge.nodeB == ln1) {
			return true
		}
	}
	return false
}

func (g *LinkedUndirectedGraph[N, E]) isConnectedNodesWithEdge(
	ln1 *LinkedUndirectedGraphNode[N, E], e E, ln2 *LinkedUndirectedGraphNode[N, E]) bool {
	edges := ln1.NeighborEdges
	if edges2 := ln2.NeighborEdges; len(edges2) < len(edges) {
		edges = edges2
	}
	for _, outEdge := range edges {
		if (outEdge.nodeA == ln1 && outEdge.nodeB == ln2) ||
			(outEdge.nodeA == ln2 && outEdge.nodeB == ln1) {
			if jutil.ObjectsEqualsOf(e, outEdge.value) {
				return true
			}
		}
	}
	return false
}

// Java: Graph.clearNodeAnnotations()
func (g *LinkedUndirectedGraph[N, E]) ClearNodeAnnotations() {
	for n := range g.Nodes.Values() {
		n.SetAnnotation(nil)
	}
}

// Java: Graph.clearEdgeAnnotations()
func (g *LinkedUndirectedGraph[N, E]) ClearEdgeAnnotations() {
	for node := range g.Nodes.Values() {
		for _, edge := range node.NeighborEdges {
			if edge.nodeA == node {
				edge.SetAnnotation(nil)
			}
		}
	}
}

// Java: Graph.pushNodeAnnotations()
func (g *LinkedUndirectedGraph[N, E]) PushNodeAnnotations() {
	if g.nodeAnnotationStack == nil {
		g.nodeAnnotationStack = &annotationStack{}
	}
	top := pushAnnotationFrame(g.nodeAnnotationStack, g.Nodes.Size())
	for n := range g.Nodes.Values() {
		*top = append(*top, annotationState{first: n, second: n.GetAnnotation()})
		n.SetAnnotation(nil)
	}
}

// Java: Graph.popNodeAnnotations()
func (g *LinkedUndirectedGraph[N, E]) PopNodeAnnotations() {
	if g.nodeAnnotationStack == nil {
		panic(jutil.NewNullPointerExceptionMsg("Popping node annotations without pushing."))
	}
	popAnnotations(g.nodeAnnotationStack)
}

// Java: Graph.pushEdgeAnnotations()
func (g *LinkedUndirectedGraph[N, E]) PushEdgeAnnotations() {
	if g.edgeAnnotationStack == nil {
		g.edgeAnnotationStack = &annotationStack{}
	}
	top := pushAnnotationFrame(g.edgeAnnotationStack, g.Nodes.Size())
	for node := range g.Nodes.Values() {
		for _, edge := range node.NeighborEdges {
			if edge.nodeA == node {
				*top = append(*top, annotationState{first: edge, second: edge.GetAnnotation()})
				edge.SetAnnotation(nil)
			}
		}
	}
}

// Java: Graph.popEdgeAnnotations()
func (g *LinkedUndirectedGraph[N, E]) PopEdgeAnnotations() {
	if g.edgeAnnotationStack == nil {
		panic(jutil.NewNullPointerExceptionMsg("Popping edge annotations without pushing."))
	}
	popAnnotations(g.edgeAnnotationStack)
}

// Java: LinkedUndirectedGraph.getGraphvizEdges()
func (g *LinkedUndirectedGraph[N, E]) GetGraphvizEdges() []GraphvizGraph_GraphvizEdge {
	edgeList := []GraphvizGraph_GraphvizEdge{}
	for node := range g.Nodes.Values() {
		for _, edge := range node.NeighborEdges {
			if edge.nodeA == node {
				edgeList = append(edgeList, edge)
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
	nodeList := make([]GraphvizGraph_GraphvizNode, 0, g.Nodes.Size())
	for node := range g.Nodes.Values() {
		nodeList = append(nodeList, node)
	}
	return nodeList
}

// Java: LinkedUndirectedGraph.isDirected()
func (g *LinkedUndirectedGraph[N, E]) IsDirected() bool {
	return false
}

// Java: LinkedUndirectedGraph.getNodes()
func (g *LinkedUndirectedGraph[N, E]) GetNodes() []*LinkedUndirectedGraphNode[N, E] {
	out := make([]*LinkedUndirectedGraphNode[N, E], 0, g.Nodes.Size())
	for n := range g.Nodes.Values() {
		out = append(out, n)
	}
	return out
}

// Java: LinkedUndirectedGraph.getNodeCount()
func (g *LinkedUndirectedGraph[N, E]) GetNodeCount() int {
	return g.Nodes.Size()
}

// Java: LinkedUndirectedGraph.getNodeDegree(N)
func (g *LinkedUndirectedGraph[N, E]) GetNodeDegree(value N) int {
	uNode := g.Nodes.Get(value)
	if uNode == nil {
		panic(jutil.NewIllegalArgumentExceptionMsg(jutil.ToJavaString(value) + " not found in graph"))
	}
	return len(uNode.NeighborEdges)
}

// Java: Graph.getWeight(N)
func (g *LinkedUndirectedGraph[N, E]) GetWeight(value N) int {
	return g.GetNodeDegree(value)
}

// LinkedUndirectedGraphNode is an undirected graph node storing incident edges
// inline when at most two.
type LinkedUndirectedGraphNode[N, E any] struct {
	NeighborEdges      []*LinkedUndirectedGraphEdge[N, E]
	neighborEdgeBuffer [2]*LinkedUndirectedGraphEdge[N, E]
	Value              N

	annotated  bool
	Annotation Annotation
}

// recycle zeroes n for allocNode to reuse, keeping small backing arrays.
func (n *LinkedUndirectedGraphNode[N, E]) recycle() {
	edges := n.NeighborEdges
	*n = LinkedUndirectedGraphNode[N, E]{}
	if c := cap(edges); c > len(n.neighborEdgeBuffer) && c <= 64 {
		clear(edges)
		n.NeighborEdges = edges[:0]
	}
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphNode.getNeighborEdges()
func (n *LinkedUndirectedGraphNode[N, E]) GetNeighborEdges() []*LinkedUndirectedGraphEdge[N, E] {
	return n.NeighborEdges
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphNode.getNeighborEdgesIterator()
func (n *LinkedUndirectedGraphNode[N, E]) GetNeighborEdgesIterator() iter.Seq[*LinkedUndirectedGraphEdge[N, E]] {
	return slices.Values(n.NeighborEdges)
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphNode.getAnnotation()
func (n *LinkedUndirectedGraphNode[N, E]) GetAnnotation() Annotation {
	if n.annotated {
		return n.Annotation
	}
	panic(jutil.NewUnsupportedOperationExceptionMsg(
		"Graph initialized with node annotations turned off"))
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphNode.setAnnotation(Annotation)
func (n *LinkedUndirectedGraphNode[N, E]) SetAnnotation(data Annotation) {
	if n.annotated {
		n.Annotation = data
		return
	}
	panic(jutil.NewUnsupportedOperationExceptionMsg(
		"Graph initialized with node annotations turned off"))
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphNode.getValue()
func (n *LinkedUndirectedGraphNode[N, E]) GetValue() N {
	return n.Value
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphNode.getColor()
func (n *LinkedUndirectedGraphNode[N, E]) GetColor() string {
	return "white"
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphNode.getId()
func (n *LinkedUndirectedGraphNode[N, E]) GetId() string {
	return "LDN" + jutil.ToJavaString(jutil.IdentityHashCode(n))
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphNode.getLabel()
func (n *LinkedUndirectedGraphNode[N, E]) GetLabel() string {
	return jutil.ToJavaString(n.Value)
}

func (n *LinkedUndirectedGraphNode[N, E]) String() string {
	className := "com.google.javascript.jscomp.graph.LinkedUndirectedGraph$LinkedUndirectedGraphNode"
	if n.annotated {
		className = "com.google.javascript.jscomp.graph.LinkedUndirectedGraph$AnnotatedLinkedUndirectedGraphNode"
	}
	return jutil.ObjectToString(className, n)
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphNode.neighborList()
func (n *LinkedUndirectedGraphNode[N, E]) neighborList() []*LinkedUndirectedGraphNode[N, E] {
	result := make([]*LinkedUndirectedGraphNode[N, E], len(n.NeighborEdges))
	for i, edge := range n.NeighborEdges {
		if edge.nodeA == n {
			result[i] = edge.nodeB
		} else {
			result[i] = edge.nodeA
		}
	}
	return result
}

// LinkedUndirectedGraphEdge is an undirected graph edge storing concrete
// endpoint node pointers.
type LinkedUndirectedGraphEdge[N, E any] struct {
	nodeA *LinkedUndirectedGraphNode[N, E]
	nodeB *LinkedUndirectedGraphNode[N, E]
	value E

	annotated  bool
	annotation Annotation
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphEdge.getValue()
func (e *LinkedUndirectedGraphEdge[N, E]) GetValue() E {
	return e.value
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphEdge.getNodeA()
func (e *LinkedUndirectedGraphEdge[N, E]) GetNodeA() *LinkedUndirectedGraphNode[N, E] {
	return e.nodeA
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphEdge.getNodeB()
func (e *LinkedUndirectedGraphEdge[N, E]) GetNodeB() *LinkedUndirectedGraphNode[N, E] {
	return e.nodeB
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphEdge.getAnnotation()
func (e *LinkedUndirectedGraphEdge[N, E]) GetAnnotation() Annotation {
	if e.annotated {
		return e.annotation
	}
	panic(jutil.NewUnsupportedOperationExceptionMsg(
		"Graph initialized with edge annotations turned off"))
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphEdge.setAnnotation(Annotation)
func (e *LinkedUndirectedGraphEdge[N, E]) SetAnnotation(data Annotation) {
	if e.annotated {
		e.annotation = data
		return
	}
	panic(jutil.NewUnsupportedOperationExceptionMsg(
		"Graph initialized with edge annotations turned off"))
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphEdge.getColor()
func (e *LinkedUndirectedGraphEdge[N, E]) GetColor() string {
	return "black"
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphEdge.getLabel()
func (e *LinkedUndirectedGraphEdge[N, E]) GetLabel() string {
	return jutil.ToJavaString(e.value)
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphEdge.getNode1Id()
func (e *LinkedUndirectedGraphEdge[N, E]) GetNode1Id() string {
	return e.nodeA.GetId()
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphEdge.getNode2Id()
func (e *LinkedUndirectedGraphEdge[N, E]) GetNode2Id() string {
	return e.nodeB.GetId()
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphEdge.toString()
func (e *LinkedUndirectedGraphEdge[N, E]) String() string {
	return jutil.ToJavaString(e.nodeA) + " -- " + jutil.ToJavaString(e.nodeB)
}

var (
	_ Graph[string, string]      = (*LinkedUndirectedGraph[string, string])(nil)
	_ GraphvizGraph              = (*LinkedUndirectedGraph[string, string])(nil)
	_ GraphNode[string, string]  = (*LinkedUndirectedGraphNode[string, string])(nil)
	_ GraphvizGraph_GraphvizNode = (*LinkedUndirectedGraphNode[string, string])(nil)
	_ GraphvizGraph_GraphvizEdge = (*LinkedUndirectedGraphEdge[string, string])(nil)
)
