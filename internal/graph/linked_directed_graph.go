// Port of com/google/javascript/jscomp/graph/LinkedDirectedGraph.java @ closure-compiler v20261021

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
	"slices"
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// LinkedDirectedGraph is a directed graph storing incoming and outgoing edges
// directly within slab-allocated nodes.
//
// N is the value type that the graph node stores; E is the value type that
// the graph edge stores.
type LinkedDirectedGraph[N comparable, E any] struct {
	Nodes jutil.LinkedHashMap[N, *LinkedDiGraphNode[N, E]]

	// Go-only: performance — see jutil.ReuseSlab.
	nodeSlab jutil.ReuseSlab[LinkedDiGraphNode[N, E]]
	edgeSlab jutil.ReuseSlab[LinkedDiGraphEdge[N, E]]

	nodeAnnotationStack *annotationStack
	edgeAnnotationStack *annotationStack

	nextNodeIndex      uint32
	nextEdgeIndex      uint32
	useNodeAnnotations bool
	useEdgeAnnotations bool
}

// Type aliases preserving callers that spell out the nested Java names.
type (
	LinkedDirectedGraph_LinkedDiGraphNode[N, E any] = LinkedDiGraphNode[N, E]
	LinkedDirectedGraph_LinkedDiGraphEdge[N, E any] = LinkedDiGraphEdge[N, E]
	DiGraph_DiGraphNode[N, E any]                   = *LinkedDiGraphNode[N, E]
	DiGraph_DiGraphEdge[N, E any]                   = *LinkedDiGraphEdge[N, E]
)

type linkedDirectedSubGraph[N comparable, E any] struct {
	graph *LinkedDirectedGraph[N, E]
	nodes []*LinkedDiGraphNode[N, E]
	buf   [8]*LinkedDiGraphNode[N, E]
}

// Java: LinkedDirectedGraph.newSubGraph()
func (g *LinkedDirectedGraph[N, E]) NewSubGraph() SubGraph[N, E] {
	s := &linkedDirectedSubGraph[N, E]{graph: g}
	s.nodes = s.buf[:0]
	return s
}

func (s *linkedDirectedSubGraph[N, E]) IsIndependentOf(value N) bool {
	node := s.graph.Nodes.Get(value)
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

// NewLinkedDirectedGraph constructs a LinkedDirectedGraph.
//
// Java: LinkedDirectedGraph(boolean, boolean)
func NewLinkedDirectedGraph[N comparable, E any](useNodeAnnotations bool, useEdgeAnnotations bool) *LinkedDirectedGraph[N, E] {
	g := &LinkedDirectedGraph[N, E]{}
	g.Init(useNodeAnnotations, useEdgeAnnotations)
	return g
}

// Init configures g's annotation flags and resets its node/edge index counters.
func (g *LinkedDirectedGraph[N, E]) Init(useNodeAnnotations bool, useEdgeAnnotations bool) {
	g.nextNodeIndex = 0
	g.nextEdgeIndex = 0
	g.useNodeAnnotations = useNodeAnnotations
	g.useEdgeAnnotations = useEdgeAnnotations
}

// NodeIndexCap returns an upper bound on Index + 1 for all nodes allocated by g.
func (g *LinkedDirectedGraph[N, E]) NodeIndexCap() int {
	return int(g.nextNodeIndex)
}

// EdgeIndexCap returns an upper bound on Index + 1 for all edges allocated by g.
func (g *LinkedDirectedGraph[N, E]) EdgeIndexCap() int {
	return int(g.nextEdgeIndex)
}

// AsLinkedDirectedGraph returns g, satisfying DiGraph for both
// *LinkedDirectedGraph and types embedding LinkedDirectedGraph.
func (g *LinkedDirectedGraph[N, E]) AsLinkedDirectedGraph() *LinkedDirectedGraph[N, E] {
	return g
}

// linkedDirectedGraphReuseMaxNodes is the most nodes a graph may have had for
// ClearForReuse to keep its node map (clearing a Go map costs time in
// proportion to its largest size).
const linkedDirectedGraphReuseMaxNodes = 1024

// ClearForReuse removes all nodes and edges from g, keeping its slab chunks
// and (when small) its node map for the next graph built in g, which must
// then be initialized again.
//
// Go-only: performance — see ReusableCfg.
func (g *LinkedDirectedGraph[N, E]) ClearForReuse() {
	g.nodeAnnotationStack = nil
	g.edgeAnnotationStack = nil
	g.nextNodeIndex = 0
	g.nextEdgeIndex = 0
	if g.Nodes.Size() > linkedDirectedGraphReuseMaxNodes {
		g.Nodes.Reset()
	} else {
		g.Nodes.Clear()
	}
	g.nodeSlab.Reset((*LinkedDiGraphNode[N, E]).recycle)
	g.edgeSlab.Reset(nil)
}

// getLinkedNodeOrFail is getNodeOrFail(N).
func (g *LinkedDirectedGraph[N, E]) getLinkedNodeOrFail(val N) *LinkedDiGraphNode[N, E] {
	if n := g.Nodes.Get(val); n != nil {
		return n
	}
	panic(jutil.NewIllegalArgumentExceptionMsg(formatGraphValue(val) + " does not exist in graph"))
}

// Java: Graph.getNodeOrFail(N)
func (g *LinkedDirectedGraph[N, E]) GetNodeOrFail(val N) *LinkedDiGraphNode[N, E] {
	return g.getLinkedNodeOrFail(val)
}

// AllocNode allocates a node from the graph's slab without inserting it into Nodes.
func (g *LinkedDirectedGraph[N, E]) AllocNode(nodeValue N) *LinkedDiGraphNode[N, E] {
	n := g.nodeSlab.Alloc()
	n.Value = nodeValue
	n.Index = g.nextNodeIndex
	g.nextNodeIndex++
	n.Priority = -1
	n.annotated = g.useNodeAnnotations
	return n
}

// newEdge allocates a directed edge from the graph's slab.
func (g *LinkedDirectedGraph[N, E]) newEdge(
	src *LinkedDiGraphNode[N, E], edgeValue E, dest *LinkedDiGraphNode[N, E]) *LinkedDiGraphEdge[N, E] {
	e := g.edgeSlab.Alloc()
	e.SourceNode = src
	e.DestNode = dest
	e.Value = edgeValue
	e.Index = g.nextEdgeIndex
	g.nextEdgeIndex++
	e.annotated = g.useEdgeAnnotations
	return e
}

// ConnectLinkedNodes connects two already-resolved linked nodes directly.
func (g *LinkedDirectedGraph[N, E]) ConnectLinkedNodes(
	src *LinkedDiGraphNode[N, E], edgeValue E, dest *LinkedDiGraphNode[N, E]) {
	edge := g.newEdge(src, edgeValue, dest)
	src.OutEdgeList.Add(edge)
	dest.InEdgeList.Add(edge)
}

// Java: LinkedDirectedGraph.connect(N, E, N)
func (g *LinkedDirectedGraph[N, E]) Connect(srcValue N, edgeValue E, destValue N) {
	src := g.getLinkedNodeOrFail(srcValue)
	dest := g.getLinkedNodeOrFail(destValue)
	g.ConnectLinkedNodes(src, edgeValue, dest)
}

// ConnectNodes connects two resolved graph nodes without map lookups.
//
// Java: LinkedDirectedGraph.connect(DiGraphNode<N, E>, E, DiGraphNode<N, E>)
func (g *LinkedDirectedGraph[N, E]) ConnectNodes(src *LinkedDiGraphNode[N, E], edgeValue E, dest *LinkedDiGraphNode[N, E]) {
	jutil.CheckNotNull(src)
	jutil.CheckNotNull(dest)
	g.ConnectLinkedNodes(src, edgeValue, dest)
}

// Java: LinkedDirectedGraph.connectIfNotConnectedInDirection(N, E, N)
func (g *LinkedDirectedGraph[N, E]) ConnectIfNotConnectedInDirection(srcValue N, edgeValue E, destValue N) {
	src := g.CreateLinkedDiGraphNode(srcValue)
	dest := g.CreateLinkedDiGraphNode(destValue)
	if !g.isConnectedInDirectionOfNodesWithEdge(src, edgeValue, dest) {
		g.ConnectLinkedNodes(src, edgeValue, dest)
	}
}

// Java: Graph.connectIfNotFound(N, E, N)
func (g *LinkedDirectedGraph[N, E]) ConnectIfNotFound(n1 N, edge E, n2 N) {
	dNode1 := g.getLinkedNodeOrFail(n1)
	dNode2 := g.getLinkedNodeOrFail(n2)
	if !g.isConnectedInDirectionOfNodesWithEdge(dNode1, edge, dNode2) &&
		!g.isConnectedInDirectionOfNodesWithEdge(dNode2, edge, dNode1) {
		g.ConnectLinkedNodes(dNode1, edge, dNode2)
	}
}

// Java: LinkedDirectedGraph.disconnect(N, N)
func (g *LinkedDirectedGraph[N, E]) Disconnect(n1 N, n2 N) {
	g.DisconnectInDirection(n1, n2)
	g.DisconnectInDirection(n2, n1)
}

// Java: LinkedDirectedGraph.disconnectInDirection(N, N)
func (g *LinkedDirectedGraph[N, E]) DisconnectInDirection(srcValue N, destValue N) {
	src := g.getLinkedNodeOrFail(srcValue)
	dest := g.getLinkedNodeOrFail(destValue)
	for _, e := range g.GetEdgesInDirection(srcValue, destValue) {
		src.OutEdgeList.Remove(e)
		dest.InEdgeList.Remove(e)
	}
}

// Java: LinkedDirectedGraph.getNodes()
func (g *LinkedDirectedGraph[N, E]) GetNodes() []*LinkedDiGraphNode[N, E] {
	out := make([]*LinkedDiGraphNode[N, E], 0, g.Nodes.Size())
	for n := range g.Nodes.Values() {
		out = append(out, n)
	}
	return out
}

// Java: LinkedDirectedGraph.getNodes() (as DiGraph.getNodes())
func (g *LinkedDirectedGraph[N, E]) GetDiGraphNodes() []*LinkedDiGraphNode[N, E] {
	return g.GetNodes()
}

// Java: LinkedDirectedGraph.getNodes() (Collection<LinkedDiGraphNode<N, E>>)
func (g *LinkedDirectedGraph[N, E]) GetLinkedDiGraphNodes() []*LinkedDiGraphNode[N, E] {
	return g.GetNodes()
}

// Java: Graph.hasNode(N)
func (g *LinkedDirectedGraph[N, E]) HasNode(n N) bool {
	return g.Nodes.Get(n) != nil
}

// Java: LinkedDirectedGraph.getNode(N)
func (g *LinkedDirectedGraph[N, E]) GetNode(nodeValue N) *LinkedDiGraphNode[N, E] {
	return g.Nodes.Get(nodeValue)
}

// Java: LinkedDirectedGraph.getNode(N) (as DiGraph.getNode(N))
func (g *LinkedDirectedGraph[N, E]) GetDiGraphNode(nodeValue N) *LinkedDiGraphNode[N, E] {
	return g.Nodes.Get(nodeValue)
}

// Java: LinkedDirectedGraph.getNode(N) (LinkedDiGraphNode<N, E>)
func (g *LinkedDirectedGraph[N, E]) GetLinkedDiGraphNode(nodeValue N) *LinkedDiGraphNode[N, E] {
	return g.Nodes.Get(nodeValue)
}

// Java: LinkedDirectedGraph.getInEdges(N)
func (g *LinkedDirectedGraph[N, E]) GetInEdges(nodeValue N) []*LinkedDiGraphEdge[N, E] {
	node := g.getLinkedNodeOrFail(nodeValue)
	return node.InEdgeList.Slice()
}

// Java: LinkedDirectedGraph.getOutEdges(N)
func (g *LinkedDirectedGraph[N, E]) GetOutEdges(nodeValue N) []*LinkedDiGraphEdge[N, E] {
	node := g.getLinkedNodeOrFail(nodeValue)
	return node.OutEdgeList.Slice()
}

// Java: LinkedDirectedGraph.createNode(N)
func (g *LinkedDirectedGraph[N, E]) CreateNode(nodeValue N) *LinkedDiGraphNode[N, E] {
	return g.CreateLinkedDiGraphNode(nodeValue)
}

// Java: LinkedDirectedGraph.createNode(N) (as DiGraph.createNode(N))
func (g *LinkedDirectedGraph[N, E]) CreateDiGraphNode(nodeValue N) *LinkedDiGraphNode[N, E] {
	return g.CreateLinkedDiGraphNode(nodeValue)
}

// Java: LinkedDirectedGraph.createNode(N) (LinkedDiGraphNode<N, E>)
func (g *LinkedDirectedGraph[N, E]) CreateLinkedDiGraphNode(nodeValue N) *LinkedDiGraphNode[N, E] {
	if n := g.Nodes.Get(nodeValue); n != nil {
		return n
	}
	n := g.AllocNode(nodeValue)
	g.Nodes.PutNew(nodeValue, n)
	return n
}

// Java: LinkedDirectedGraph.getEdges(N, N)
func (g *LinkedDirectedGraph[N, E]) GetEdgesBetween(n1 N, n2 N) []*LinkedDiGraphEdge[N, E] {
	return g.GetDiGraphEdgesBetween(n1, n2)
}

// Java: LinkedDirectedGraph.getEdges(N, N) (as DiGraph.getEdges(N, N))
func (g *LinkedDirectedGraph[N, E]) GetDiGraphEdgesBetween(n1 N, n2 N) []*LinkedDiGraphEdge[N, E] {
	// Since this is a method from a generic graph, edges from both
	// directions must be added to the returning list.
	forwardEdges := g.GetEdgesInDirection(n1, n2)
	backwardEdges := g.GetEdgesInDirection(n2, n1)
	edges := make([]*LinkedDiGraphEdge[N, E], 0, len(forwardEdges)+len(backwardEdges))
	edges = append(edges, forwardEdges...)
	edges = append(edges, backwardEdges...)
	return edges
}

// Java: LinkedDirectedGraph.getEdges()
func (g *LinkedDirectedGraph[N, E]) GetEdges() []*LinkedDiGraphEdge[N, E] {
	result := []*LinkedDiGraphEdge[N, E]{}
	for node := range g.Nodes.Values() {
		result = append(result, node.OutEdgeList.Slice()...)
	}
	return result
}

// Java: LinkedDirectedGraph.getEdges() (as DiGraph.getEdges())
func (g *LinkedDirectedGraph[N, E]) GetDiGraphEdges() []*LinkedDiGraphEdge[N, E] {
	return g.GetEdges()
}

// Java: LinkedDirectedGraph.getFirstEdge(N, N)
func (g *LinkedDirectedGraph[N, E]) GetFirstEdge(n1 N, n2 N) *LinkedDiGraphEdge[N, E] {
	dNode1 := g.getLinkedNodeOrFail(n1)
	dNode2 := g.getLinkedNodeOrFail(n2)
	for _, outEdge := range dNode1.OutEdgeList.Slice() {
		if outEdge.DestNode == dNode2 {
			return outEdge
		}
	}
	for _, outEdge := range dNode2.OutEdgeList.Slice() {
		if outEdge.DestNode == dNode1 {
			return outEdge
		}
	}
	return nil
}

// Java: LinkedDirectedGraph.getEdgesInDirection(N, N)
func (g *LinkedDirectedGraph[N, E]) GetEdgesInDirection(n1 N, n2 N) []*LinkedDiGraphEdge[N, E] {
	dNode1 := g.getLinkedNodeOrFail(n1)
	dNode2 := g.getLinkedNodeOrFail(n2)
	edges := []*LinkedDiGraphEdge[N, E]{}
	for _, outEdge := range dNode1.OutEdgeList.Slice() {
		if outEdge.DestNode == dNode2 {
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

// Java: LinkedDirectedGraph.isConnectedInDirection(N, N)
func (g *LinkedDirectedGraph[N, E]) IsConnectedInDirection(source N, dest N) bool {
	dNode1 := g.getLinkedNodeOrFail(source)
	dNode2 := g.getLinkedNodeOrFail(dest)
	return g.isConnectedInDirectionOfNodesAny(dNode1, dNode2)
}

// Java: LinkedDirectedGraph.isConnectedInDirection(N, E, N)
func (g *LinkedDirectedGraph[N, E]) IsConnectedInDirectionWithEdge(source N, edgeValue E, dest N) bool {
	dNode1 := g.getLinkedNodeOrFail(source)
	dNode2 := g.getLinkedNodeOrFail(dest)
	return g.isConnectedInDirectionOfNodesWithEdge(dNode1, edgeValue, dNode2)
}

// IsConnectedInDirectionOfNodesAny reports whether source has any directed edge to dest.
func (g *LinkedDirectedGraph[N, E]) IsConnectedInDirectionOfNodesAny(
	source *LinkedDiGraphNode[N, E], dest *LinkedDiGraphNode[N, E]) bool {
	return g.isConnectedInDirectionOfNodesAny(source, dest)
}

func (g *LinkedDirectedGraph[N, E]) isConnectedInDirectionOfNodesAny(
	source *LinkedDiGraphNode[N, E], dest *LinkedDiGraphNode[N, E]) bool {
	outEdges := source.OutEdgeList.Slice()
	inEdges := dest.InEdgeList.Slice()
	if len(outEdges) < len(inEdges) {
		for _, outEdge := range outEdges {
			if outEdge.DestNode == dest {
				return true
			}
		}
	} else {
		for _, inEdge := range inEdges {
			if inEdge.SourceNode == source {
				return true
			}
		}
	}
	return false
}

func (g *LinkedDirectedGraph[N, E]) isConnectedInDirectionOfNodesWithEdge(
	source *LinkedDiGraphNode[N, E], edgeValue E, dest *LinkedDiGraphNode[N, E]) bool {
	outEdges := source.OutEdgeList.Slice()
	inEdges := dest.InEdgeList.Slice()
	if len(outEdges) < len(inEdges) {
		for _, outEdge := range outEdges {
			if outEdge.DestNode == dest && jutil.ObjectsEqualsOf(edgeValue, outEdge.Value) {
				return true
			}
		}
	} else {
		for _, inEdge := range inEdges {
			if inEdge.SourceNode == source && jutil.ObjectsEqualsOf(edgeValue, inEdge.Value) {
				return true
			}
		}
	}
	return false
}

// IsConnectedInDirectionOfNodes checks whether source has a directed edge to
// dest matching edgeFilter, scanning whichever endpoint has fewer incident edges.
//
// Java: LinkedDirectedGraph.isConnectedInDirection(LinkedDiGraphNode<N, E>, Predicate<E>, LinkedDiGraphNode<N, E>)
func (g *LinkedDirectedGraph[N, E]) IsConnectedInDirectionOfNodes(
	source *LinkedDiGraphNode[N, E], edgeFilter func(E) bool, dest *LinkedDiGraphNode[N, E]) bool {
	outEdges := source.OutEdgeList.Slice()
	inEdges := dest.InEdgeList.Slice()
	if len(outEdges) < len(inEdges) {
		for _, outEdge := range outEdges {
			if outEdge.DestNode == dest && edgeFilter(outEdge.Value) {
				return true
			}
		}
	} else {
		for _, inEdge := range inEdges {
			if inEdge.SourceNode == source && edgeFilter(inEdge.Value) {
				return true
			}
		}
	}
	return false
}

// Java: Graph.clearNodeAnnotations()
func (g *LinkedDirectedGraph[N, E]) ClearNodeAnnotations() {
	for n := range g.Nodes.Values() {
		n.SetAnnotation(nil)
	}
}

// Java: Graph.clearEdgeAnnotations()
func (g *LinkedDirectedGraph[N, E]) ClearEdgeAnnotations() {
	for node := range g.Nodes.Values() {
		for _, e := range node.OutEdgeList.Slice() {
			e.SetAnnotation(nil)
		}
	}
}

// Java: Graph.pushNodeAnnotations()
func (g *LinkedDirectedGraph[N, E]) PushNodeAnnotations() {
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
func (g *LinkedDirectedGraph[N, E]) PopNodeAnnotations() {
	if g.nodeAnnotationStack == nil {
		panic(jutil.NewNullPointerExceptionMsg("Popping node annotations without pushing."))
	}
	popAnnotations(g.nodeAnnotationStack)
}

// Java: Graph.pushEdgeAnnotations()
func (g *LinkedDirectedGraph[N, E]) PushEdgeAnnotations() {
	if g.edgeAnnotationStack == nil {
		g.edgeAnnotationStack = &annotationStack{}
	}
	top := pushAnnotationFrame(g.edgeAnnotationStack, g.Nodes.Size())
	for node := range g.Nodes.Values() {
		for _, e := range node.OutEdgeList.Slice() {
			*top = append(*top, annotationState{first: e, second: e.GetAnnotation()})
			e.SetAnnotation(nil)
		}
	}
}

// Java: Graph.popEdgeAnnotations()
func (g *LinkedDirectedGraph[N, E]) PopEdgeAnnotations() {
	if g.edgeAnnotationStack == nil {
		panic(jutil.NewNullPointerExceptionMsg("Popping edge annotations without pushing."))
	}
	popAnnotations(g.edgeAnnotationStack)
}

// Java: LinkedDirectedGraph.getDirectedPredNodes(N)
func (g *LinkedDirectedGraph[N, E]) GetDirectedPredNodesOfValue(nodeValue N) []*LinkedDiGraphNode[N, E] {
	return g.GetDirectedPredNodes(g.Nodes.Get(nodeValue))
}

// Java: LinkedDirectedGraph.getDirectedPredNodes(DiGraphNode<N, E>)
func (g *LinkedDirectedGraph[N, E]) GetDirectedPredNodes(dNode *LinkedDiGraphNode[N, E]) []*LinkedDiGraphNode[N, E] {
	jutil.CheckNotNull(dNode)
	inEdges := dNode.InEdgeList.Slice()
	nodeList := make([]*LinkedDiGraphNode[N, E], len(inEdges))
	for i, edge := range inEdges {
		nodeList[i] = edge.SourceNode
	}
	return nodeList
}

// Java: LinkedDirectedGraph.getDirectedSuccNodes(N)
func (g *LinkedDirectedGraph[N, E]) GetDirectedSuccNodesOfValue(nodeValue N) []*LinkedDiGraphNode[N, E] {
	return g.GetDirectedSuccNodes(g.Nodes.Get(nodeValue))
}

// Java: LinkedDirectedGraph.getDirectedSuccNodes(DiGraphNode<N, E>)
func (g *LinkedDirectedGraph[N, E]) GetDirectedSuccNodes(dNode *LinkedDiGraphNode[N, E]) []*LinkedDiGraphNode[N, E] {
	jutil.CheckNotNull(dNode)
	outEdges := dNode.OutEdgeList.Slice()
	nodeList := make([]*LinkedDiGraphNode[N, E], len(outEdges))
	for i, edge := range outEdges {
		nodeList[i] = edge.DestNode
	}
	return nodeList
}

// Java: LinkedDirectedGraph.getGraphvizEdges()
func (g *LinkedDirectedGraph[N, E]) GetGraphvizEdges() []GraphvizGraph_GraphvizEdge {
	edgeList := []GraphvizGraph_GraphvizEdge{}
	for node := range g.Nodes.Values() {
		for _, edge := range node.OutEdgeList.Slice() {
			edgeList = append(edgeList, edge)
		}
	}
	return edgeList
}

// Java: LinkedDirectedGraph.getGraphvizNodes()
func (g *LinkedDirectedGraph[N, E]) GetGraphvizNodes() []GraphvizGraph_GraphvizNode {
	nodeList := make([]GraphvizGraph_GraphvizNode, 0, g.Nodes.Size())
	for node := range g.Nodes.Values() {
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
	return g.Nodes.Size()
}

// Java: LinkedDirectedGraph.getNeighborNodes(N)
func (g *LinkedDirectedGraph[N, E]) GetNeighborNodes(value N) []*LinkedDiGraphNode[N, E] {
	node := g.Nodes.Get(value)
	inEdges := node.InEdgeList.Slice()
	outEdges := node.OutEdgeList.Slice()
	result := make([]*LinkedDiGraphNode[N, E], 0, len(inEdges)+len(outEdges))
	for _, inEdge := range inEdges {
		result = append(result, inEdge.SourceNode)
	}
	for _, outEdge := range outEdges {
		result = append(result, outEdge.DestNode)
	}
	return result
}

// Java: LinkedDirectedGraph.getNodeDegree(N)
func (g *LinkedDirectedGraph[N, E]) GetNodeDegree(value N) int {
	node := g.getLinkedNodeOrFail(value)
	return node.InEdgeList.Size() + node.OutEdgeList.Size()
}

// Java: Graph.getWeight(N)
func (g *LinkedDirectedGraph[N, E]) GetWeight(value N) int {
	return g.GetNodeDegree(value)
}

// LinkedDiGraphNode is a directed graph node storing incoming and outgoing
// edges inline when at most two per direction.
type LinkedDiGraphNode[N, E any] struct {
	InEdgeList, OutEdgeList LinkedEdgeList[N, E]

	Annotation Annotation

	Value N

	Index      uint32
	Priority   int32
	annotated  bool
	inWorkSet  bool
	DfaInQueue bool
}

// recycle zeroes n for AllocNode to reuse, keeping small overflow arrays.
func (n *LinkedDiGraphNode[N, E]) recycle() {
	in, out := n.InEdgeList.spare(), n.OutEdgeList.spare()
	*n = LinkedDiGraphNode[N, E]{}
	n.InEdgeList.more, n.OutEdgeList.more = in, out
}

// LinkedEdgeList is the list of a graph node's in or out edges, in insertion order.
//
// Go-only: performance — most CFG nodes have at most two edges in each
// direction, so up to two are stored inline without a separate slice header.
type LinkedEdgeList[N, E any] struct {
	two  [2]*LinkedDiGraphEdge[N, E]
	more *[]*LinkedDiGraphEdge[N, E]
}

// Slice returns the edges, aliasing the list's storage.
func (l *LinkedEdgeList[N, E]) Slice() []*LinkedDiGraphEdge[N, E] {
	switch {
	case l.more != nil:
		return *l.more
	case l.two[1] != nil:
		return l.two[:2]
	case l.two[0] != nil:
		return l.two[:1]
	}
	return nil
}

// Size returns the number of edges.
func (l *LinkedEdgeList[N, E]) Size() int {
	return len(l.Slice())
}

// Add appends e.
func (l *LinkedEdgeList[N, E]) Add(e *LinkedDiGraphEdge[N, E]) {
	switch {
	case l.more != nil:
		*l.more = append(*l.more, e)
	case l.two[0] == nil:
		l.two[0] = e
	case l.two[1] == nil:
		l.two[1] = e
	default:
		all := make([]*LinkedDiGraphEdge[N, E], 3, 4)
		all[0], all[1], all[2] = l.two[0], l.two[1], e
		l.two = [2]*LinkedDiGraphEdge[N, E]{}
		l.more = &all
	}
}

// Remove removes the first occurrence of e, preserving order.
func (l *LinkedEdgeList[N, E]) Remove(e *LinkedDiGraphEdge[N, E]) {
	switch {
	case l.more != nil:
		if i := slices.Index(*l.more, e); i >= 0 {
			*l.more = slices.Delete(*l.more, i, i+1)
		}
	case l.two[0] == e:
		l.two[0], l.two[1] = l.two[1], nil
	case l.two[1] == e:
		l.two[1] = nil
	}
}

func (l *LinkedEdgeList[N, E]) spare() *[]*LinkedDiGraphEdge[N, E] {
	if l.more == nil || cap(*l.more) > 64 {
		return nil
	}
	clear(*l.more)
	*l.more = (*l.more)[:0]
	return l.more
}

// Java: LinkedDirectedGraph.LinkedDiGraphNode.getValue()
func (n *LinkedDiGraphNode[N, E]) GetValue() N {
	return n.Value
}

// Java: LinkedDirectedGraph.LinkedDiGraphNode.getAnnotation()
func (n *LinkedDiGraphNode[N, E]) GetAnnotation() Annotation {
	if n.annotated {
		return n.Annotation
	}
	panic(jutil.NewUnsupportedOperationExceptionMsg(
		"Graph initialized with node annotations turned off"))
}

// Java: LinkedDirectedGraph.LinkedDiGraphNode.setAnnotation(Annotation)
func (n *LinkedDiGraphNode[N, E]) SetAnnotation(data Annotation) {
	if n.annotated {
		n.Annotation = data
		return
	}
	panic(jutil.NewUnsupportedOperationExceptionMsg(
		"Graph initialized with node annotations turned off"))
}

// Java: LinkedDirectedGraph.LinkedDiGraphNode.getColor()
func (n *LinkedDiGraphNode[N, E]) GetColor() string {
	return "white"
}

// Java: LinkedDirectedGraph.LinkedDiGraphNode.getId()
func (n *LinkedDiGraphNode[N, E]) GetId() string {
	return "LDN" + strconv.FormatInt(int64(jutil.IdentityHashCode(n)), 10)
}

// Java: LinkedDirectedGraph.LinkedDiGraphNode.getLabel()
func (n *LinkedDiGraphNode[N, E]) GetLabel() string {
	if n.annotated {
		result := n.String()
		if n.Annotation != nil {
			result += "\n" + formatGraphValue(n.Annotation)
		}
		return result
	}
	return n.String()
}

// Java: LinkedDirectedGraph.LinkedDiGraphNode.toString()
func (n *LinkedDiGraphNode[N, E]) String() string {
	return formatGraphValue(n.Value)
}

// Java: LinkedDirectedGraph.LinkedDiGraphNode.getInEdges()
func (n *LinkedDiGraphNode[N, E]) GetInEdges() []*LinkedDiGraphEdge[N, E] {
	return n.InEdgeList.Slice()
}

// Java: LinkedDirectedGraph.LinkedDiGraphNode.getOutEdges()
func (n *LinkedDiGraphNode[N, E]) GetOutEdges() []*LinkedDiGraphEdge[N, E] {
	return n.OutEdgeList.Slice()
}

// Java: LinkedDirectedGraph.LinkedDiGraphNode.hasPriority()
func (n *LinkedDiGraphNode[N, E]) HasPriority() bool {
	return n.Priority >= 0
}

// Java: LinkedDirectedGraph.LinkedDiGraphNode.getPriority()
func (n *LinkedDiGraphNode[N, E]) GetPriority() int {
	if !(n.Priority >= 0) {
		panic(jutil.NewIllegalStateExceptionMsg("priority not set"))
	}
	return int(n.Priority)
}

// Java: LinkedDirectedGraph.LinkedDiGraphNode.setPriority(int)
func (n *LinkedDiGraphNode[N, E]) SetPriority(priority int) {
	if !(priority >= 0) {
		panic(jutil.NewIllegalArgumentExceptionMsg("priorities must be non-negative"))
	}
	n.Priority = int32(priority)
}

// LinkedDiGraphEdge is a directed graph edge storing concrete source and
// destination node pointers.
type LinkedDiGraphEdge[N, E any] struct {
	SourceNode *LinkedDiGraphNode[N, E]
	DestNode   *LinkedDiGraphNode[N, E]
	Value      E

	Annotation Annotation
	Index      uint32
	annotated  bool
}

// Java: LinkedDirectedGraph.LinkedDiGraphEdge.getSource()
func (e *LinkedDiGraphEdge[N, E]) GetSource() *LinkedDiGraphNode[N, E] {
	return e.SourceNode
}

// Java: LinkedDirectedGraph.LinkedDiGraphEdge.getDestination()
func (e *LinkedDiGraphEdge[N, E]) GetDestination() *LinkedDiGraphNode[N, E] {
	return e.DestNode
}

// GetLinkedSource returns SourceNode.
func (e *LinkedDiGraphEdge[N, E]) GetLinkedSource() *LinkedDiGraphNode[N, E] {
	return e.SourceNode
}

// GetLinkedDestination returns DestNode.
func (e *LinkedDiGraphEdge[N, E]) GetLinkedDestination() *LinkedDiGraphNode[N, E] {
	return e.DestNode
}

// Java: LinkedDirectedGraph.LinkedDiGraphEdge.getValue()
func (e *LinkedDiGraphEdge[N, E]) GetValue() E {
	return e.Value
}

// Java: LinkedDirectedGraph.LinkedDiGraphEdge.getAnnotation()
func (e *LinkedDiGraphEdge[N, E]) GetAnnotation() Annotation {
	if e.annotated {
		return e.Annotation
	}
	panic(jutil.NewUnsupportedOperationExceptionMsg(
		"Graph initialized with edge annotations turned off"))
}

// Java: LinkedDirectedGraph.LinkedDiGraphEdge.setAnnotation(Annotation)
func (e *LinkedDiGraphEdge[N, E]) SetAnnotation(data Annotation) {
	if e.annotated {
		e.Annotation = data
		return
	}
	panic(jutil.NewUnsupportedOperationExceptionMsg(
		"Graph initialized with edge annotations turned off"))
}

// Java: LinkedDirectedGraph.LinkedDiGraphEdge.getColor()
func (e *LinkedDiGraphEdge[N, E]) GetColor() string {
	return "black"
}

// Java: LinkedDirectedGraph.LinkedDiGraphEdge.getLabel()
func (e *LinkedDiGraphEdge[N, E]) GetLabel() string {
	return formatGraphValue(e.Value)
}

// Java: LinkedDirectedGraph.LinkedDiGraphEdge.getNode1Id()
func (e *LinkedDiGraphEdge[N, E]) GetNode1Id() string {
	return e.SourceNode.GetId()
}

// Java: LinkedDirectedGraph.LinkedDiGraphEdge.getNode2Id()
func (e *LinkedDiGraphEdge[N, E]) GetNode2Id() string {
	return e.DestNode.GetId()
}

// Java: LinkedDirectedGraph.LinkedDiGraphEdge.toString()
func (e *LinkedDiGraphEdge[N, E]) String() string {
	return e.SourceNode.String() + " -> " + e.DestNode.String()
}

// Java: LinkedDirectedGraph.LinkedDiGraphEdge.getNodeA()
func (e *LinkedDiGraphEdge[N, E]) GetNodeA() *LinkedDiGraphNode[N, E] {
	return e.SourceNode
}

// Java: LinkedDirectedGraph.LinkedDiGraphEdge.getNodeB()
func (e *LinkedDiGraphEdge[N, E]) GetNodeB() *LinkedDiGraphNode[N, E] {
	return e.DestNode
}

var (
	_ Graph[string, string]      = (*LinkedDirectedGraph[string, string])(nil)
	_ DiGraph[string, string]    = (*LinkedDirectedGraph[string, string])(nil)
	_ GraphvizGraph              = (*LinkedDirectedGraph[string, string])(nil)
	_ GraphNode[string, string]  = (*LinkedDiGraphNode[string, string])(nil)
	_ GraphvizGraph_GraphvizNode = (*LinkedDiGraphNode[string, string])(nil)
	_ GraphvizGraph_GraphvizEdge = (*LinkedDiGraphEdge[string, string])(nil)
)
