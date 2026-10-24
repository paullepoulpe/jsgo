// Port of com/google/javascript/jscomp/graph/LinkedDirectedGraph.java @ closure-compiler v20261015

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
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// LinkedDirectedGraph is a directed graph using ArrayLists within nodes to
// store edge information.
//
// This implementation favors directed graph operations inherited from
// DirectedGraph. Operations from Graph would tend to be slower.
//
// N is the value type that the graph node stores; E is the value type that
// the graph edge stores.
//
// Go-only: language-adapter — node values are keyed with Go == (Java
// equals() for strings, numbers and identity-compared objects). The
// covariant overrides narrowing to LinkedDiGraphNode/LinkedDiGraphEdge are
// the methods named after them (CreateLinkedDiGraphNode,
// GetLinkedDiGraphNodes, GetLinkedDiGraphNode); see Graph.
type LinkedDirectedGraph[N comparable, E any] struct {
	DiGraphBase[N, E]

	nodes jutil.LinkedHashMap[N, *LinkedDirectedGraph_LinkedDiGraphNode[N, E]]

	// Go-only: performance — see reuseSlab.
	nodeSlab reuseSlab[LinkedDirectedGraph_LinkedDiGraphNode[N, E]]
	edgeSlab reuseSlab[LinkedDirectedGraph_LinkedDiGraphEdge[N, E]]

	useNodeAnnotations bool
	useEdgeAnnotations bool
}

type linkedDirectedSubGraph[N comparable, E any] struct {
	graph *LinkedDirectedGraph[N, E]
	nodes []*LinkedDirectedGraph_LinkedDiGraphNode[N, E]
	buf   [8]*LinkedDirectedGraph_LinkedDiGraphNode[N, E]
}

// Java: LinkedDirectedGraph.newSubGraph()
func (g *LinkedDirectedGraph[N, E]) NewSubGraph() SubGraph[N, E] {
	s := &linkedDirectedSubGraph[N, E]{graph: g}
	s.nodes = s.buf[:0]
	return s
}

func (s *linkedDirectedSubGraph[N, E]) IsIndependentOf(value N) bool {
	node := s.graph.nodes.Get(value)
	if node == nil {
		return true
	}
	for _, n := range s.nodes {
		if s.graph.isConnectedInDirectionOfNodesAny(n, node) ||
			s.graph.isConnectedInDirectionOfNodesAny(node, n) {
			return false
		}
	}
	return true
}

func (s *linkedDirectedSubGraph[N, E]) AddNode(value N) {
	s.nodes = append(s.nodes, s.graph.getLinkedNodeOrFail(value))
}

// Java: LinkedDirectedGraph.createWithoutAnnotations()
func LinkedDirectedGraph_CreateWithoutAnnotations[N comparable, E any]() *LinkedDirectedGraph[N, E] {
	return NewLinkedDirectedGraph[N, E](false, false)
}

// Java: LinkedDirectedGraph.create()
func LinkedDirectedGraph_Create[N comparable, E any]() *LinkedDirectedGraph[N, E] {
	return NewLinkedDirectedGraph[N, E](true, true)
}

// NewLinkedDirectedGraph is the protected constructor.
//
// Java: LinkedDirectedGraph(boolean, boolean)
func NewLinkedDirectedGraph[N comparable, E any](useNodeAnnotations bool, useEdgeAnnotations bool) *LinkedDirectedGraph[N, E] {
	g := &LinkedDirectedGraph[N, E]{}
	g.InitLinkedDirectedGraph(g, useNodeAnnotations, useEdgeAnnotations)
	return g
}

// InitLinkedDirectedGraph runs the LinkedDirectedGraph(boolean, boolean)
// constructor on g (for subclasses, which embed LinkedDirectedGraph); self
// is the object being constructed.
func (g *LinkedDirectedGraph[N, E]) InitLinkedDirectedGraph(self DiGraph[N, E], useNodeAnnotations bool, useEdgeAnnotations bool) {
	g.InitDiGraphBase(self)
	g.useNodeAnnotations = useNodeAnnotations
	g.useEdgeAnnotations = useEdgeAnnotations
}

// linkedDirectedGraphReuseMaxNodes is the most nodes a graph may have had for
// clearForReuse to keep its node map (clearing a Go map costs time in
// proportion to its largest size).
const linkedDirectedGraphReuseMaxNodes = 1024

// clearForReuse removes all nodes and edges from g, keeping its slab chunks
// and (when small) its node map for the next graph built in g, which must
// then be initialized again.
//
// Go-only: performance — see ReusableCfg.
func (g *LinkedDirectedGraph[N, E]) clearForReuse() {
	g.nodeAnnotationStack = nil
	g.edgeAnnotationStack = nil
	if g.nodes.Size() > linkedDirectedGraphReuseMaxNodes {
		g.nodes.Reset()
	} else {
		g.nodes.Clear()
	}
	g.nodeSlab.reset((*LinkedDirectedGraph_LinkedDiGraphNode[N, E]).recycle)
	g.edgeSlab.reset(nil)
}

// getLinkedNodeOrFail is getNodeOrFail(N) with the LinkedDiGraphNode cast.
func (g *LinkedDirectedGraph[N, E]) getLinkedNodeOrFail(val N) *LinkedDirectedGraph_LinkedDiGraphNode[N, E] {
	if n := g.nodes.Get(val); n != nil {
		return n
	}
	panic(jutil.NewIllegalArgumentExceptionMsg(jutil.ToJavaString(val) + " does not exist in graph"))
}

func (g *LinkedDirectedGraph[N, E]) allocNode(nodeValue N) *LinkedDirectedGraph_LinkedDiGraphNode[N, E] {
	n := g.nodeSlab.alloc()
	n.value = nodeValue
	n.priority = -1
	n.annotated = g.useNodeAnnotations
	if cap(n.inEdgeList.Slice()) == 0 {
		n.inEdgeList.InitWithBuffer(n.inEdgeBuffer[:]) // initialCapacity= 1
	}
	if cap(n.outEdgeList.Slice()) == 0 {
		n.outEdgeList.InitWithBuffer(n.outEdgeBuffer[:]) // initialCapacity= 1
	}
	return n
}

// newEdge is the conditional edge construction of connect().
func (g *LinkedDirectedGraph[N, E]) newEdge(
	src *LinkedDirectedGraph_LinkedDiGraphNode[N, E], edgeValue E, dest *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) *LinkedDirectedGraph_LinkedDiGraphEdge[N, E] {
	e := g.edgeSlab.alloc()
	e.sourceNode = src
	e.destNode = dest
	e.value = edgeValue
	e.annotated = g.useEdgeAnnotations
	return e
}

func (g *LinkedDirectedGraph[N, E]) connectLinkedNodes(
	src *LinkedDirectedGraph_LinkedDiGraphNode[N, E], edgeValue E, dest *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) {
	edge := g.newEdge(src, edgeValue, dest)
	src.outEdgeList.Add(edge)
	dest.inEdgeList.Add(edge)
}

// Java: LinkedDirectedGraph.connect(N, E, N)
func (g *LinkedDirectedGraph[N, E]) Connect(srcValue N, edgeValue E, destValue N) {
	src := g.getLinkedNodeOrFail(srcValue)
	dest := g.getLinkedNodeOrFail(destValue)
	g.connectLinkedNodes(src, edgeValue, dest)
}

// ConnectNodes: DiGraphNode look ups can be expensive for a large graph
// operation, prefer this method if you have the DiGraphNode available.
//
// TODO(johnlenz): make this part of the general graph interface.
//
// Java: LinkedDirectedGraph.connect(DiGraphNode<N, E>, E, DiGraphNode<N, E>)
func (g *LinkedDirectedGraph[N, E]) ConnectNodes(src DiGraph_DiGraphNode[N, E], edgeValue E, dest DiGraph_DiGraphNode[N, E]) {
	lSrc := src.(*LinkedDirectedGraph_LinkedDiGraphNode[N, E])
	lDest := dest.(*LinkedDirectedGraph_LinkedDiGraphNode[N, E])
	g.connectLinkedNodes(lSrc, edgeValue, lDest)
}

// ConnectIfNotConnectedInDirection: DiGraphNode look ups can be expensive
// for a large graph operation, prefer this method if you have the
// DiGraphNode available.
//
// TODO(johnlenz): make this part of the general graph interface.
//
// Java: LinkedDirectedGraph.connectIfNotConnectedInDirection(N, E, N)
func (g *LinkedDirectedGraph[N, E]) ConnectIfNotConnectedInDirection(srcValue N, edgeValue E, destValue N) {
	src := g.CreateLinkedDiGraphNode(srcValue)
	dest := g.CreateLinkedDiGraphNode(destValue)
	if !g.isConnectedInDirectionOfNodesWithEdge(src, edgeValue, dest) {
		g.connectLinkedNodes(src, edgeValue, dest)
	}
}

// Java: Graph.connectIfNotFound(N, E, N)
func (g *LinkedDirectedGraph[N, E]) ConnectIfNotFound(n1 N, edge E, n2 N) {
	dNode1 := g.getLinkedNodeOrFail(n1)
	dNode2 := g.getLinkedNodeOrFail(n2)
	if !g.isConnectedInDirectionOfNodesWithEdge(dNode1, edge, dNode2) &&
		!g.isConnectedInDirectionOfNodesWithEdge(dNode2, edge, dNode1) {
		g.connectLinkedNodes(dNode1, edge, dNode2)
	}
}

// linkedDirectedGraph_equalTo is Predicates.equalTo(target).
func linkedDirectedGraph_equalTo[E any](target E) func(E) bool {
	return func(e E) bool { return jutil.ObjectsEquals(target, e) }
}

// linkedDirectedGraph_alwaysTrue is Predicates.alwaysTrue().
func linkedDirectedGraph_alwaysTrue[E any](E) bool { return true }

// Java: LinkedDirectedGraph.disconnect(N, N)
func (g *LinkedDirectedGraph[N, E]) Disconnect(n1 N, n2 N) {
	g.diSelf.DisconnectInDirection(n1, n2)
	g.diSelf.DisconnectInDirection(n2, n1)
}

// Java: LinkedDirectedGraph.disconnectInDirection(N, N)
func (g *LinkedDirectedGraph[N, E]) DisconnectInDirection(srcValue N, destValue N) {
	src := g.getLinkedNodeOrFail(srcValue)
	dest := g.getLinkedNodeOrFail(destValue)
	for _, edge := range g.diSelf.GetEdgesInDirection(srcValue, destValue) {
		src.outEdgeList.Remove(edge)
		dest.inEdgeList.Remove(edge)
	}
}

// Java: LinkedDirectedGraph.getNodes()
func (g *LinkedDirectedGraph[N, E]) GetNodes() []GraphNode[N, E] {
	out := make([]GraphNode[N, E], 0, g.nodes.Size())
	for n := range g.nodes.Values() {
		out = append(out, n)
	}
	return out
}

// Java: LinkedDirectedGraph.getNodes() (as DiGraph.getNodes())
func (g *LinkedDirectedGraph[N, E]) GetDiGraphNodes() []DiGraph_DiGraphNode[N, E] {
	out := make([]DiGraph_DiGraphNode[N, E], 0, g.nodes.Size())
	for n := range g.nodes.Values() {
		out = append(out, n)
	}
	return out
}

// Java: LinkedDirectedGraph.getNodes() (Collection<LinkedDiGraphNode<N, E>>)
func (g *LinkedDirectedGraph[N, E]) GetLinkedDiGraphNodes() []*LinkedDirectedGraph_LinkedDiGraphNode[N, E] {
	return slices.Collect(g.nodes.Values())
}

// Java: Graph.hasNode(N)
func (g *LinkedDirectedGraph[N, E]) HasNode(n N) bool {
	return g.nodes.Get(n) != nil
}

// Java: LinkedDirectedGraph.getNode(N)
func (g *LinkedDirectedGraph[N, E]) GetNode(nodeValue N) GraphNode[N, E] {
	if n := g.nodes.Get(nodeValue); n != nil {
		return n
	}
	return nil
}

// Java: LinkedDirectedGraph.getNode(N) (as DiGraph.getNode(N))
func (g *LinkedDirectedGraph[N, E]) GetDiGraphNode(nodeValue N) DiGraph_DiGraphNode[N, E] {
	if n := g.nodes.Get(nodeValue); n != nil {
		return n
	}
	return nil
}

// Java: LinkedDirectedGraph.getNode(N) (LinkedDiGraphNode<N, E>)
func (g *LinkedDirectedGraph[N, E]) GetLinkedDiGraphNode(nodeValue N) *LinkedDirectedGraph_LinkedDiGraphNode[N, E] {
	return g.nodes.Get(nodeValue)
}

// Java: LinkedDirectedGraph.getInEdges(N)
func (g *LinkedDirectedGraph[N, E]) GetInEdges(nodeValue N) []DiGraph_DiGraphEdge[N, E] {
	node := g.getLinkedNodeOrFail(nodeValue)
	return node.GetInEdges()
}

// Java: LinkedDirectedGraph.getOutEdges(N)
func (g *LinkedDirectedGraph[N, E]) GetOutEdges(nodeValue N) []DiGraph_DiGraphEdge[N, E] {
	node := g.getLinkedNodeOrFail(nodeValue)
	return node.GetOutEdges()
}

// Java: LinkedDirectedGraph.createNode(N)
func (g *LinkedDirectedGraph[N, E]) CreateNode(nodeValue N) GraphNode[N, E] {
	return g.CreateLinkedDiGraphNode(nodeValue)
}

// Java: LinkedDirectedGraph.createNode(N) (as DiGraph.createNode(N))
func (g *LinkedDirectedGraph[N, E]) CreateDiGraphNode(nodeValue N) DiGraph_DiGraphNode[N, E] {
	return g.CreateLinkedDiGraphNode(nodeValue)
}

// Java: LinkedDirectedGraph.createNode(N) (LinkedDiGraphNode<N, E>)
func (g *LinkedDirectedGraph[N, E]) CreateLinkedDiGraphNode(nodeValue N) *LinkedDirectedGraph_LinkedDiGraphNode[N, E] {
	if n := g.nodes.Get(nodeValue); n != nil {
		return n
	}
	n := g.allocNode(nodeValue)
	g.nodes.PutNew(nodeValue, n)
	return n
}

// Java: LinkedDirectedGraph.getEdges(N, N)
func (g *LinkedDirectedGraph[N, E]) GetEdgesBetween(n1 N, n2 N) []Graph_GraphEdge[N, E] {
	edges := g.GetDiGraphEdgesBetween(n1, n2)
	out := make([]Graph_GraphEdge[N, E], len(edges))
	for i, e := range edges {
		out[i] = e
	}
	return out
}

// Java: LinkedDirectedGraph.getEdges(N, N) (as DiGraph.getEdges(N, N))
func (g *LinkedDirectedGraph[N, E]) GetDiGraphEdgesBetween(n1 N, n2 N) []DiGraph_DiGraphEdge[N, E] {
	// Since this is a method from a generic graph, edges from both
	// directions must be added to the returning list.
	forwardEdges := g.diSelf.GetEdgesInDirection(n1, n2)
	backwardEdges := g.diSelf.GetEdgesInDirection(n2, n1)
	totalSize := len(forwardEdges) + len(backwardEdges)
	edges := make([]DiGraph_DiGraphEdge[N, E], 0, totalSize)
	edges = append(edges, forwardEdges...)
	edges = append(edges, backwardEdges...)
	return edges
}

// Java: LinkedDirectedGraph.getEdges()
func (g *LinkedDirectedGraph[N, E]) GetEdges() []Graph_GraphEdge[N, E] {
	result := []Graph_GraphEdge[N, E]{}
	for node := range g.nodes.Values() {
		for _, e := range node.outEdgeList.Slice() {
			result = append(result, e)
		}
	}
	return result
}

// Java: LinkedDirectedGraph.getEdges() (as DiGraph.getEdges())
func (g *LinkedDirectedGraph[N, E]) GetDiGraphEdges() []DiGraph_DiGraphEdge[N, E] {
	result := []DiGraph_DiGraphEdge[N, E]{}
	for node := range g.nodes.Values() {
		result = append(result, node.outEdgeList.Slice()...)
	}
	return result
}

// Java: LinkedDirectedGraph.getFirstEdge(N, N)
func (g *LinkedDirectedGraph[N, E]) GetFirstEdge(n1 N, n2 N) Graph_GraphEdge[N, E] {
	dNode1 := g.getLinkedNodeOrFail(n1)
	dNode2 := g.getLinkedNodeOrFail(n2)
	for _, outEdge := range dNode1.outEdgeList.Slice() {
		if outEdge.GetDestination() == DiGraph_DiGraphNode[N, E](dNode2) {
			return outEdge
		}
	}
	for _, outEdge := range dNode2.outEdgeList.Slice() {
		if outEdge.GetDestination() == DiGraph_DiGraphNode[N, E](dNode1) {
			return outEdge
		}
	}
	return nil
}

// Java: LinkedDirectedGraph.getEdgesInDirection(N, N)
func (g *LinkedDirectedGraph[N, E]) GetEdgesInDirection(n1 N, n2 N) []DiGraph_DiGraphEdge[N, E] {
	dNode1 := g.getLinkedNodeOrFail(n1)
	dNode2 := g.getLinkedNodeOrFail(n2)
	edges := []DiGraph_DiGraphEdge[N, E]{}
	for _, outEdge := range dNode1.outEdgeList.Slice() {
		if outEdge.GetDestination() == DiGraph_DiGraphNode[N, E](dNode2) {
			edges = append(edges, outEdge)
		}
	}
	return edges
}

// Java: DiGraph.isConnected(N, N)
func (g *LinkedDirectedGraph[N, E]) IsConnected(n1 N, n2 N) bool {
	dNode1 := g.getLinkedNodeOrFail(n1)
	dNode2 := g.getLinkedNodeOrFail(n2)
	return g.isConnectedInDirectionOfNodesAny(dNode1, dNode2) ||
		g.isConnectedInDirectionOfNodesAny(dNode2, dNode1)
}

// Java: DiGraph.isConnected(N, E, N)
func (g *LinkedDirectedGraph[N, E]) IsConnectedWithEdge(n1 N, e E, n2 N) bool {
	dNode1 := g.getLinkedNodeOrFail(n1)
	dNode2 := g.getLinkedNodeOrFail(n2)
	return g.isConnectedInDirectionOfNodesWithEdge(dNode1, e, dNode2) ||
		g.isConnectedInDirectionOfNodesWithEdge(dNode2, e, dNode1)
}

// IsConnectedInDirection: DiGraphNode look ups can be expensive for a large
// graph operation, prefer the version below that takes DiGraphNodes, if you
// have them available.
//
// source is the source node from which we traverse outwards.
//
// Java: LinkedDirectedGraph.isConnectedInDirection(N, N)
func (g *LinkedDirectedGraph[N, E]) IsConnectedInDirection(source N, dest N) bool {
	dNode1 := g.getLinkedNodeOrFail(source)
	dNode2 := g.getLinkedNodeOrFail(dest)
	return g.isConnectedInDirectionOfNodesAny(dNode1, dNode2)
}

// IsConnectedInDirectionWithEdge: DiGraphNode look ups can be expensive for
// a large graph operation, prefer the version below that takes
// DiGraphNodes, if you have them available.
//
// source is the source node from which we traverse outwards; only edges
// equal to edgeValue will be traversed.
//
// Java: LinkedDirectedGraph.isConnectedInDirection(N, E, N)
func (g *LinkedDirectedGraph[N, E]) IsConnectedInDirectionWithEdge(source N, edgeValue E, dest N) bool {
	dNode1 := g.getLinkedNodeOrFail(source)
	dNode2 := g.getLinkedNodeOrFail(dest)
	return g.isConnectedInDirectionOfNodesWithEdge(dNode1, edgeValue, dNode2)
}

func (g *LinkedDirectedGraph[N, E]) isConnectedInDirectionOfNodesAny(
	source *LinkedDirectedGraph_LinkedDiGraphNode[N, E], dest *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) bool {
	outEdges := source.outEdgeList.Slice()
	inEdges := dest.inEdgeList.Slice()
	if len(outEdges) < len(inEdges) {
		for _, outEdge := range outEdges {
			if outEdge.GetDestination() == DiGraph_DiGraphNode[N, E](dest) {
				return true
			}
		}
	} else {
		for _, inEdge := range inEdges {
			if inEdge.GetSource() == DiGraph_DiGraphNode[N, E](source) {
				return true
			}
		}
	}
	return false
}

func (g *LinkedDirectedGraph[N, E]) isConnectedInDirectionOfNodesWithEdge(
	source *LinkedDirectedGraph_LinkedDiGraphNode[N, E], edgeValue E, dest *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) bool {
	outEdges := source.outEdgeList.Slice()
	inEdges := dest.inEdgeList.Slice()
	if len(outEdges) < len(inEdges) {
		for _, outEdge := range outEdges {
			if outEdge.GetDestination() == DiGraph_DiGraphNode[N, E](dest) && jutil.ObjectsEqualsOf(edgeValue, outEdge.GetValue()) {
				return true
			}
		}
	} else {
		for _, inEdge := range inEdges {
			if inEdge.GetSource() == DiGraph_DiGraphNode[N, E](source) && jutil.ObjectsEqualsOf(edgeValue, inEdge.GetValue()) {
				return true
			}
		}
	}
	return false
}

// IsConnectedInDirectionOfNodes: DiGraphNode look ups can be expensive for
// a large graph operation, prefer this method if you have the DiGraphNodes
// available.
//
// source is the source node from which we traverse outwards; only edges
// matching edgeFilter will be traversed; dest is the destination node.
//
// Java: LinkedDirectedGraph.isConnectedInDirection(LinkedDiGraphNode<N, E>, Predicate<E>, LinkedDiGraphNode<N, E>)
func (g *LinkedDirectedGraph[N, E]) IsConnectedInDirectionOfNodes(
	source *LinkedDirectedGraph_LinkedDiGraphNode[N, E], edgeFilter func(E) bool, dest *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) bool {
	outEdges := source.outEdgeList.Slice()
	outEdgesLen := len(outEdges)
	inEdges := dest.inEdgeList.Slice()
	inEdgesLen := len(inEdges)
	// It is possible that there is a large asymmetry between the nodes, so
	// pick the direction to search based on the shorter list since the edge
	// lists should be symmetric.
	if outEdgesLen < inEdgesLen {
		for _, outEdge := range outEdges {
			if outEdge.GetDestination() == DiGraph_DiGraphNode[N, E](dest) && edgeFilter(outEdge.GetValue()) {
				return true
			}
		}
	} else {
		for _, inEdge := range inEdges {
			if inEdge.GetSource() == DiGraph_DiGraphNode[N, E](source) && edgeFilter(inEdge.GetValue()) {
				return true
			}
		}
	}

	return false
}

// Java: LinkedDirectedGraph.isConnectedInDirection(N, Predicate<E>, N)
func (g *LinkedDirectedGraph[N, E]) isConnectedInDirectionWithPredicate(n1 N, edgeMatcher func(E) bool, n2 N) bool {
	// Verify the nodes.
	dNode1 := g.getLinkedNodeOrFail(n1)
	dNode2 := g.getLinkedNodeOrFail(n2)
	return g.IsConnectedInDirectionOfNodes(dNode1, edgeMatcher, dNode2)
}

// Java: Graph.clearNodeAnnotations()
func (g *LinkedDirectedGraph[N, E]) ClearNodeAnnotations() {
	for n := range g.nodes.Values() {
		n.SetAnnotation(nil)
	}
}

// Java: Graph.clearEdgeAnnotations()
func (g *LinkedDirectedGraph[N, E]) ClearEdgeAnnotations() {
	for node := range g.nodes.Values() {
		for _, e := range node.outEdgeList.Slice() {
			e.SetAnnotation(nil)
		}
	}
}

// Java: Graph.pushNodeAnnotations()
func (g *LinkedDirectedGraph[N, E]) PushNodeAnnotations() {
	if g.nodeAnnotationStack == nil {
		g.nodeAnnotationStack = &graph_annotationStack{}
	}
	top := graph_pushAnnotationFrame(g.nodeAnnotationStack, g.nodes.Size())
	for n := range g.nodes.Values() {
		*top = append(*top, graph_AnnotationState{first: n, second: n.GetAnnotation()})
		n.SetAnnotation(nil)
	}
}

// Java: Graph.pushEdgeAnnotations()
func (g *LinkedDirectedGraph[N, E]) PushEdgeAnnotations() {
	if g.edgeAnnotationStack == nil {
		g.edgeAnnotationStack = &graph_annotationStack{}
	}
	top := graph_pushAnnotationFrame(g.edgeAnnotationStack, g.nodes.Size())
	for node := range g.nodes.Values() {
		for _, e := range node.outEdgeList.Slice() {
			*top = append(*top, graph_AnnotationState{first: e, second: e.GetAnnotation()})
			e.SetAnnotation(nil)
		}
	}
}

// Java: LinkedDirectedGraph.getDirectedPredNodes(N)
func (g *LinkedDirectedGraph[N, E]) GetDirectedPredNodesOfValue(nodeValue N) []DiGraph_DiGraphNode[N, E] {
	return g.diSelf.GetDirectedPredNodes(g.GetDiGraphNode(nodeValue))
}

// Java: LinkedDirectedGraph.getDirectedPredNodes(DiGraphNode<N, E>)
func (g *LinkedDirectedGraph[N, E]) GetDirectedPredNodes(dNode DiGraph_DiGraphNode[N, E]) []DiGraph_DiGraphNode[N, E] {
	jutil.CheckNotNull(dNode)
	lNode := dNode.(*LinkedDirectedGraph_LinkedDiGraphNode[N, E])

	nodeList := make([]DiGraph_DiGraphNode[N, E], 0, lNode.inEdgeList.Size())
	for _, edge := range lNode.inEdgeList.Slice() {
		nodeList = append(nodeList, edge.GetSource())
	}
	return nodeList
}

// Java: LinkedDirectedGraph.getDirectedSuccNodes(N)
func (g *LinkedDirectedGraph[N, E]) GetDirectedSuccNodesOfValue(nodeValue N) []DiGraph_DiGraphNode[N, E] {
	return g.diSelf.GetDirectedSuccNodes(g.GetDiGraphNode(nodeValue))
}

// Java: LinkedDirectedGraph.getDirectedSuccNodes(DiGraphNode<N, E>)
func (g *LinkedDirectedGraph[N, E]) GetDirectedSuccNodes(dNode DiGraph_DiGraphNode[N, E]) []DiGraph_DiGraphNode[N, E] {
	jutil.CheckNotNull(dNode)
	lNode := dNode.(*LinkedDirectedGraph_LinkedDiGraphNode[N, E])

	nodeList := make([]DiGraph_DiGraphNode[N, E], 0, lNode.outEdgeList.Size())
	for _, edge := range lNode.outEdgeList.Slice() {
		nodeList = append(nodeList, edge.GetDestination())
	}
	return nodeList
}

// Java: LinkedDirectedGraph.getGraphvizEdges()
func (g *LinkedDirectedGraph[N, E]) GetGraphvizEdges() []GraphvizGraph_GraphvizEdge {
	edgeList := []GraphvizGraph_GraphvizEdge{}
	for node := range g.nodes.Values() {
		for _, edge := range node.outEdgeList.Slice() {
			edgeList = append(edgeList, edge.(*LinkedDirectedGraph_LinkedDiGraphEdge[N, E]))
		}
	}
	return edgeList
}

// Java: LinkedDirectedGraph.getGraphvizNodes()
func (g *LinkedDirectedGraph[N, E]) GetGraphvizNodes() []GraphvizGraph_GraphvizNode {
	nodeList := make([]GraphvizGraph_GraphvizNode, 0, g.nodes.Size())
	for node := range g.nodes.Values() {
		nodeList = append(nodeList, node)
	}
	return nodeList
}

// Java: LinkedDirectedGraph.getName()
func (g *LinkedDirectedGraph[N, E]) GetName() string {
	return "LinkedGraph"
}

// Java: LinkedDirectedGraph.isDirected()
func (g *LinkedDirectedGraph[N, E]) IsDirected() bool {
	return true
}

// Java: LinkedDirectedGraph.getNodeCount()
func (g *LinkedDirectedGraph[N, E]) GetNodeCount() int {
	return g.nodes.Size()
}

// Java: LinkedDirectedGraph.getNeighborNodes(N)
func (g *LinkedDirectedGraph[N, E]) GetNeighborNodes(value N) []GraphNode[N, E] {
	node := g.GetLinkedDiGraphNode(value)
	result := make([]GraphNode[N, E], 0, node.inEdgeList.Size()+node.outEdgeList.Size())
	for _, inEdge := range node.inEdgeList.Slice() {
		result = append(result, inEdge.GetSource())
	}
	for _, outEdge := range node.outEdgeList.Slice() {
		result = append(result, outEdge.GetDestination())
	}
	return result
}

// Java: LinkedDirectedGraph.getNodeDegree(N)
func (g *LinkedDirectedGraph[N, E]) GetNodeDegree(value N) int {
	node := g.getLinkedNodeOrFail(value)
	return node.inEdgeList.Size() + node.outEdgeList.Size()
}

// LinkedDirectedGraph_LinkedDiGraphNode is a directed graph node that
// stores outgoing edges and incoming edges as an list within the node
// itself.
//
// Go-only: language-adapter — the private subclass
// AnnotatedLinkedDiGraphNode only overrides accessors, so it is folded in
// (PORTING.md §4): annotated selects its behaviour. The edge lists hold
// DiGraphEdges (the Java lists hold LinkedDiGraphEdges) so that
// GetInEdges/GetOutEdges can return them without a copy.
type LinkedDirectedGraph_LinkedDiGraphNode[N, E any] struct {
	// The overwhelming majority of nodes have in/out degree == 1.
	// Initialize our lists to account for that.
	//
	// Go-only: performance — the lists and their capacity-1 arrays are held
	// in the node (ArrayList.InitWithBuffer): one allocation per graph node
	// instead of five, and scalar fields sit at the tail to minimize GC ptrdata.
	inEdgeBuffer, outEdgeBuffer [2]DiGraph_DiGraphEdge[N, E]

	// AnnotatedLinkedDiGraphNode
	annotation Annotation

	value N

	inEdgeList jutil.ArrayList[DiGraph_DiGraphEdge[N, E]]

	outEdgeList jutil.ArrayList[DiGraph_DiGraphEdge[N, E]]

	priority   int32
	annotated  bool
	inWorkSet  bool
	dfaInQueue bool
}

// recycle zeroes n for allocNode to reuse, but keeps its edge lists' backing
// arrays (emptied) unless they are large.
func (n *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) recycle() {
	in, out := n.inEdgeList, n.outEdgeList
	*n = LinkedDirectedGraph_LinkedDiGraphNode[N, E]{}
	if c := cap(in.Slice()); c > 1 && c <= 64 {
		in.Clear()
		n.inEdgeList = in
	}
	if c := cap(out.Slice()); c > 1 && c <= 64 {
		out.Clear()
		n.outEdgeList = out
	}
}

// newLinkedDirectedGraph_LinkedDiGraphNode is the constructor; nodeValue is
// the node's value.
//
// Java: LinkedDirectedGraph.LinkedDiGraphNode(N)
func newLinkedDirectedGraph_LinkedDiGraphNode[N, E any](nodeValue N) *LinkedDirectedGraph_LinkedDiGraphNode[N, E] {
	n := &LinkedDirectedGraph_LinkedDiGraphNode[N, E]{
		value:    nodeValue,
		priority: -1,
	}
	n.inEdgeList.InitWithBuffer(n.inEdgeBuffer[:])   // initialCapacity= 1
	n.outEdgeList.InitWithBuffer(n.outEdgeBuffer[:]) // initialCapacity= 1
	return n
}

// Java: LinkedDirectedGraph.AnnotatedLinkedDiGraphNode(N)
func newLinkedDirectedGraph_AnnotatedLinkedDiGraphNode[N, E any](nodeValue N) *LinkedDirectedGraph_LinkedDiGraphNode[N, E] {
	n := newLinkedDirectedGraph_LinkedDiGraphNode[N, E](nodeValue)
	n.annotated = true
	return n
}

// Java: LinkedDirectedGraph.LinkedDiGraphNode.getValue()
func (n *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) GetValue() N {
	return n.value
}

// Java: LinkedDirectedGraph.LinkedDiGraphNode.getAnnotation(), AnnotatedLinkedDiGraphNode.getAnnotation()
func (n *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) GetAnnotation() Annotation {
	if n.annotated {
		return n.annotation
	}
	panic(jutil.NewUnsupportedOperationExceptionMsg(
		"Graph initialized with node annotations turned off"))
}

// Java: LinkedDirectedGraph.LinkedDiGraphNode.setAnnotation(Annotation), AnnotatedLinkedDiGraphNode.setAnnotation(Annotation)
func (n *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) SetAnnotation(data Annotation) {
	if n.annotated {
		n.annotation = data
		return
	}
	panic(jutil.NewUnsupportedOperationExceptionMsg(
		"Graph initialized with node annotations turned off"))
}

// Java: LinkedDirectedGraph.LinkedDiGraphNode.getColor()
func (n *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) GetColor() string {
	return "white"
}

// GetId is "LDN" + hashCode() (the identity hash code).
//
// Java: LinkedDirectedGraph.LinkedDiGraphNode.getId()
func (n *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) GetId() string {
	return "LDN" + jutil.ToJavaString(jutil.IdentityHashCode(n))
}

// Java: LinkedDirectedGraph.LinkedDiGraphNode.getLabel(), AnnotatedLinkedDiGraphNode.getLabel()
func (n *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) GetLabel() string {
	if n.annotated {
		result := n.String()
		if n.annotation != nil {
			result += "\n" + jutil.ToJavaString(n.annotation)
		}
		return result
	}
	return n.String()
}

// Java: LinkedDirectedGraph.LinkedDiGraphNode.toString()
func (n *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) String() string {
	return jutil.ToJavaString(n.value)
}

// Java: LinkedDirectedGraph.LinkedDiGraphNode.getInEdges()
func (n *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) GetInEdges() []DiGraph_DiGraphEdge[N, E] {
	return n.inEdgeList.Slice()
}

// Java: LinkedDirectedGraph.LinkedDiGraphNode.getOutEdges()
func (n *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) GetOutEdges() []DiGraph_DiGraphEdge[N, E] {
	return n.outEdgeList.Slice()
}

// Java: LinkedDirectedGraph.LinkedDiGraphNode.hasPriority()
func (n *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) HasPriority() bool {
	return n.priority >= 0
}

// Java: LinkedDirectedGraph.LinkedDiGraphNode.getPriority()
func (n *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) GetPriority() int {
	if !(n.priority >= 0) {
		// checkState(boolean, Object errorMessage)
		panic(jutil.NewIllegalStateExceptionMsg("priority not set"))
	}
	return int(n.priority)
}

// Java: LinkedDirectedGraph.LinkedDiGraphNode.setPriority(int)
func (n *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) SetPriority(priority int) {
	if !(priority >= 0) {
		// checkArgument(boolean, Object errorMessage)
		panic(jutil.NewIllegalArgumentExceptionMsg("priorities must be non-negative"))
	}
	n.priority = int32(priority)
}

// LinkedDirectedGraph_LinkedDiGraphEdge is a directed graph edge that
// stores the source and destination nodes at each edge.
//
// Go-only: language-adapter — the private subclass
// AnnotatedLinkedDiGraphEdge only overrides accessors, so it is folded in
// (annotated selects its behaviour).
type LinkedDirectedGraph_LinkedDiGraphEdge[N, E any] struct {
	sourceNode *LinkedDirectedGraph_LinkedDiGraphNode[N, E]

	destNode *LinkedDirectedGraph_LinkedDiGraphNode[N, E]

	value E

	// AnnotatedLinkedDiGraphEdge
	annotated  bool
	annotation Annotation
}

// newLinkedDirectedGraph_LinkedDiGraphEdge is the constructor; edgeValue is
// the edge value.
//
// Java: LinkedDirectedGraph.LinkedDiGraphEdge(LinkedDiGraphNode<N, E>, E, LinkedDiGraphNode<N, E>)
func newLinkedDirectedGraph_LinkedDiGraphEdge[N, E any](
	sourceNode *LinkedDirectedGraph_LinkedDiGraphNode[N, E], edgeValue E, destNode *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) *LinkedDirectedGraph_LinkedDiGraphEdge[N, E] {
	return &LinkedDirectedGraph_LinkedDiGraphEdge[N, E]{value: edgeValue, sourceNode: sourceNode, destNode: destNode}
}

// Java: LinkedDirectedGraph.AnnotatedLinkedDiGraphEdge(LinkedDiGraphNode<N, E>, E, LinkedDiGraphNode<N, E>)
func newLinkedDirectedGraph_AnnotatedLinkedDiGraphEdge[N, E any](
	sourceNode *LinkedDirectedGraph_LinkedDiGraphNode[N, E], edgeValue E, destNode *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) *LinkedDirectedGraph_LinkedDiGraphEdge[N, E] {
	e := newLinkedDirectedGraph_LinkedDiGraphEdge(sourceNode, edgeValue, destNode)
	e.annotated = true
	return e
}

// Java: LinkedDirectedGraph.LinkedDiGraphEdge.getSource()
func (e *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) GetSource() DiGraph_DiGraphNode[N, E] {
	return e.sourceNode
}

// Java: LinkedDirectedGraph.LinkedDiGraphEdge.getDestination()
func (e *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) GetDestination() DiGraph_DiGraphNode[N, E] {
	return e.destNode
}

// GetLinkedSource is getSource() with its LinkedDiGraphNode return type.
func (e *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) GetLinkedSource() *LinkedDirectedGraph_LinkedDiGraphNode[N, E] {
	return e.sourceNode
}

// GetLinkedDestination is getDestination() with its LinkedDiGraphNode
// return type.
func (e *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) GetLinkedDestination() *LinkedDirectedGraph_LinkedDiGraphNode[N, E] {
	return e.destNode
}

// Java: LinkedDirectedGraph.LinkedDiGraphEdge.getValue()
func (e *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) GetValue() E {
	return e.value
}

// Java: LinkedDirectedGraph.LinkedDiGraphEdge.getAnnotation(), AnnotatedLinkedDiGraphEdge.getAnnotation()
func (e *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) GetAnnotation() Annotation {
	if e.annotated {
		return e.annotation
	}
	panic(jutil.NewUnsupportedOperationExceptionMsg(
		"Graph initialized with edge annotations turned off"))
}

// Java: LinkedDirectedGraph.LinkedDiGraphEdge.setAnnotation(Annotation), AnnotatedLinkedDiGraphEdge.setAnnotation(Annotation)
func (e *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) SetAnnotation(data Annotation) {
	if e.annotated {
		e.annotation = data
		return
	}
	panic(jutil.NewUnsupportedOperationExceptionMsg(
		"Graph initialized with edge annotations turned off"))
}

// Java: LinkedDirectedGraph.LinkedDiGraphEdge.getColor()
func (e *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) GetColor() string {
	return "black"
}

// Java: LinkedDirectedGraph.LinkedDiGraphEdge.getLabel()
func (e *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) GetLabel() string {
	return jutil.ToJavaString(e.value)
}

// Java: LinkedDirectedGraph.LinkedDiGraphEdge.getNode1Id()
func (e *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) GetNode1Id() string {
	return e.sourceNode.GetId()
}

// Java: LinkedDirectedGraph.LinkedDiGraphEdge.getNode2Id()
func (e *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) GetNode2Id() string {
	return e.destNode.GetId()
}

// Java: LinkedDirectedGraph.LinkedDiGraphEdge.toString()
func (e *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) String() string {
	return e.sourceNode.String() + " -> " + e.destNode.String()
}

// Java: LinkedDirectedGraph.LinkedDiGraphEdge.getNodeA()
func (e *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) GetNodeA() GraphNode[N, E] {
	return e.sourceNode
}

// Java: LinkedDirectedGraph.LinkedDiGraphEdge.getNodeB()
func (e *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) GetNodeB() GraphNode[N, E] {
	return e.destNode
}

var (
	_ DiGraph[string, string]             = (*LinkedDirectedGraph[string, string])(nil)
	_ GraphvizGraph                       = (*LinkedDirectedGraph[string, string])(nil)
	_ GraphvizGraph_GraphvizNode          = (*LinkedDirectedGraph_LinkedDiGraphNode[string, string])(nil)
	_ GraphvizGraph_GraphvizEdge          = (*LinkedDirectedGraph_LinkedDiGraphEdge[string, string])(nil)
	_ DiGraph_DiGraphEdge[string, string] = (*LinkedDirectedGraph_LinkedDiGraphEdge[string, string])(nil)
)
