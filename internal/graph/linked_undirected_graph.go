// Port of com/google/javascript/jscomp/graph/LinkedUndirectedGraph.java @ closure-compiler v20261028

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

// LinkedUndirectedGraph is an undirected graph storing incident edges within
// slab-allocated nodes.
//
// N is the value type that the graph node stores; E is the value type that
// the graph edge stores.
type LinkedUndirectedGraph[N, E comparable] struct {
	Nodes GraphNodeMap[N, *LinkedUndirectedGraphNode[N, E]]

	// Go-only: performance — see jutil.ReuseSlab.
	nodeSlab jutil.ReuseSlab[LinkedUndirectedGraphNode[N, E]]
	edgeSlab jutil.ReuseSlab[LinkedUndirectedGraphEdge[N, E]]

	useNodeAnnotations bool
	useEdgeAnnotations bool
}

// Java: LinkedUndirectedGraph.create()
func LinkedUndirectedGraph_Create[N, E comparable]() *LinkedUndirectedGraph[N, E] {
	return NewLinkedUndirectedGraph[N, E](true, true)
}

// Java: LinkedUndirectedGraph(boolean, boolean)
func NewLinkedUndirectedGraph[N, E comparable](useNodeAnnotations bool, useEdgeAnnotations bool) *LinkedUndirectedGraph[N, E] {
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

func (g *LinkedUndirectedGraph[N, E]) allocNode(nodeValue N) *LinkedUndirectedGraphNode[N, E] {
	n := g.nodeSlab.Alloc()
	n.Value = nodeValue
	n.annotationSlot = makeNodeAnnotation(g.useNodeAnnotations)
	return n
}

func (g *LinkedUndirectedGraph[N, E]) newEdge(
	src *LinkedUndirectedGraphNode[N, E], edgeValue E, dest *LinkedUndirectedGraphNode[N, E]) *LinkedUndirectedGraphEdge[N, E] {
	e := g.edgeSlab.Alloc()
	e.nodeA = src
	e.nodeB = dest
	e.value = edgeValue
	e.annotationSlot = makeEdgeAnnotation(g.useEdgeAnnotations)
	return e
}

// ConnectLinkedNodes connects two already-resolved linked nodes directly.
func (g *LinkedUndirectedGraph[N, E]) ConnectLinkedNodes(
	src *LinkedUndirectedGraphNode[N, E], edgeValue E, dest *LinkedUndirectedGraphNode[N, E]) {
	edge := g.newEdge(src, edgeValue, dest)
	src.NeighborEdges.Add(edge)
	dest.NeighborEdges.Add(edge)
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
		src.NeighborEdges.Remove(edge)
		dest.NeighborEdges.Remove(edge)
	}
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
	for _, outEdge := range dNode1.NeighborEdges.Slice() {
		if outEdge.nodeA == dNode2 || outEdge.nodeB == dNode2 {
			edges = append(edges, outEdge)
		}
	}
	return edges
}

// Java: LinkedUndirectedGraph.createNode(N)
func (g *LinkedUndirectedGraph[N, E]) CreateNode(value N) *LinkedUndirectedGraphNode[N, E] {
	if n := g.Nodes.Get(value); n != nil {
		return n
	}
	n := g.allocNode(value)
	g.Nodes.PutNew(value, n)
	return n
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
	for _, node := range g.Nodes.Slice() {
		for _, edge := range node.NeighborEdges.Slice() {
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
	for _, outEdge := range dNode1.NeighborEdges.Slice() {
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

// IsConnectedNodesAny reports whether ln1 and ln2 share any edge.
func (g *LinkedUndirectedGraph[N, E]) IsConnectedNodesAny(
	ln1 *LinkedUndirectedGraphNode[N, E], ln2 *LinkedUndirectedGraphNode[N, E]) bool {
	edges := ln1.NeighborEdges.Slice()
	if edges2 := ln2.NeighborEdges.Slice(); len(edges2) < len(edges) {
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
	edges := ln1.NeighborEdges.Slice()
	if edges2 := ln2.NeighborEdges.Slice(); len(edges2) < len(edges) {
		edges = edges2
	}
	for _, outEdge := range edges {
		if (outEdge.nodeA == ln1 && outEdge.nodeB == ln2) ||
			(outEdge.nodeA == ln2 && outEdge.nodeB == ln1) {
			if e == outEdge.value {
				return true
			}
		}
	}
	return false
}

// Java: Graph.clearNodeAnnotations()
func (g *LinkedUndirectedGraph[N, E]) ClearNodeAnnotations() {
	for _, n := range g.Nodes.Slice() {
		n.SetAnnotation(nil)
	}
}

// Java: Graph.clearEdgeAnnotations()
func (g *LinkedUndirectedGraph[N, E]) ClearEdgeAnnotations() {
	for _, node := range g.Nodes.Slice() {
		for _, edge := range node.NeighborEdges.Slice() {
			if edge.nodeA == node {
				edge.SetAnnotation(nil)
			}
		}
	}
}

// Java: LinkedUndirectedGraph.getGraphvizEdges()
func (g *LinkedUndirectedGraph[N, E]) GetGraphvizEdges() []GraphvizGraph_GraphvizEdge {
	edgeList := []GraphvizGraph_GraphvizEdge{}
	for _, node := range g.Nodes.Slice() {
		for _, edge := range node.NeighborEdges.Slice() {
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
	nodes := g.Nodes.Slice()
	nodeList := make([]GraphvizGraph_GraphvizNode, len(nodes))
	for i, node := range nodes {
		nodeList[i] = node
	}
	return nodeList
}

// Java: LinkedUndirectedGraph.isDirected()
func (g *LinkedUndirectedGraph[N, E]) IsDirected() bool {
	return false
}

// Java: LinkedUndirectedGraph.getNodes()
func (g *LinkedUndirectedGraph[N, E]) GetNodes() []*LinkedUndirectedGraphNode[N, E] {
	order := g.Nodes.Slice()
	if len(order) == 0 {
		return []*LinkedUndirectedGraphNode[N, E]{}
	}
	return slices.Clone(order)
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
	return uNode.NeighborEdges.Size()
}

// LinkedUndirectedGraphNode is an undirected graph node storing incident edges
// inline when at most two.
type LinkedUndirectedGraphNode[N, E comparable] struct {
	NeighborEdges LinkedEdgeList[LinkedUndirectedGraphEdge[N, E]]
	annotationSlot
	Value    N
	colorIdx int32
}

// recycle zeroes n for allocNode to reuse, keeping small overflow arrays.
func (n *LinkedUndirectedGraphNode[N, E]) recycle() {
	more := n.NeighborEdges.spare()
	*n = LinkedUndirectedGraphNode[N, E]{}
	n.NeighborEdges.more = more
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
	return "LDN" + strconv.FormatInt(int64(jutil.IdentityHashCode(n)), 10)
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphNode.getLabel()
func (n *LinkedUndirectedGraphNode[N, E]) GetLabel() string {
	return jutil.ToJavaString(n.Value)
}

func (n *LinkedUndirectedGraphNode[N, E]) String() string {
	className := "com.google.javascript.jscomp.graph.LinkedUndirectedGraph$LinkedUndirectedGraphNode"
	if n.isAnnotated() {
		className = "com.google.javascript.jscomp.graph.LinkedUndirectedGraph$AnnotatedLinkedUndirectedGraphNode"
	}
	return jutil.ObjectToString(className, n)
}

// Java: LinkedUndirectedGraph.LinkedUndirectedGraphNode.neighborList()
func (n *LinkedUndirectedGraphNode[N, E]) neighborList() []*LinkedUndirectedGraphNode[N, E] {
	edges := n.NeighborEdges.Slice()
	result := make([]*LinkedUndirectedGraphNode[N, E], len(edges))
	for i, edge := range edges {
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
type LinkedUndirectedGraphEdge[N, E comparable] struct {
	nodeA *LinkedUndirectedGraphNode[N, E]
	nodeB *LinkedUndirectedGraphNode[N, E]
	value E

	annotationSlot
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
