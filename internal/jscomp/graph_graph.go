// Port of com/google/javascript/jscomp/graph/Graph.java @ closure-compiler v20261015

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Graph is the base generic class for graph-like data structure and
// algorithms in the compiler.
//
// Nodes and edges in the graph can store a piece of data that this graph
// is used to represent. For example, a variable interference graph might
// store a variable in the node. This piece of data can be accessed with
// GraphNode#getValue and GraphEdge#getValue. However, in some cases
// wrapping the nodes and/or edges in a custom class may be preferrable, as
// this additional piece of data must be explicitly cast down from
// Annotation to its true type.
//
// Algorithms and analysis can annotate information on the nodes and edges
// using GraphNode#getValue and GraphEdge#getValue. For example, a graph
// coloring algorithm can store the color as an annotation. If multiple
// analyses are required, it is up to the user of the analysis to save the
// annotated solution between passes.
//
// We implemented our own graph data structure (as opposed to using
// com.google.common.graph) for three reasons. First, aside from the node's
// label value, we would like to annotate information on the nodes and
// edges. Using a map to annotate would introduce too much overhead during
// fix point analysis. Also, com.google.common.graph does not support
// labeling of edges. Second, avoiding using an external package limits our
// dependencies. Third, com.google.common.graph uses WeakRefs which J2CL
// does not support.
//
// N is the value type that the graph node stores; E is the value type that
// the graph edge stores.
//
// Go-only: language-adapter — Java's covariant overrides cannot be
// expressed by Go method sets. Each method keeps the signature of its
// topmost declaration (here, Graph's), and every narrowing override adds a
// method named after the narrowed type: DiGraph's getNodes() is
// GetDiGraphNodes(), LinkedDirectedGraph's is GetLinkedDiGraphNodes(), and
// so on (see DiGraph and LinkedDirectedGraph). Collections are returned as
// slices: snapshots of Java's unmodifiable views.
type Graph[N, E any] interface {
	AdjacencyGraph[N, E]

	// Connect connects two nodes in the graph with an edge.
	//
	// Java: Graph.connect(N, E, N)
	Connect(n1 N, edge E, n2 N)

	// Disconnect disconnects two nodes in the graph by removing all edges
	// between them.
	//
	// Java: Graph.disconnect(N, N)
	Disconnect(n1 N, n2 N)

	// Java: Graph.connectIfNotFound(N, E, N)
	ConnectIfNotFound(n1 N, edge E, n2 N)

	// CreateNode gets a node from the graph given a value. New nodes are
	// created if that value has not been assigned a graph node. Values
	// equality are compared using Object.equals.
	//
	// Java: Graph.createNode(N)
	CreateNode(value N) GraphNode[N, E]

	// GetEdges gets an immutable list of all edges.
	//
	// Java: Graph.getEdges()
	GetEdges() []Graph_GraphEdge[N, E]

	// GetEdgesBetween retrieves the edges between two values in the graph.
	//
	// Java: Graph.getEdges(N, N)
	GetEdgesBetween(n1 N, n2 N) []Graph_GraphEdge[N, E]

	// GetNodeDegree gets the degree of a node.
	//
	// Java: Graph.getNodeDegree(N)
	GetNodeDegree(value N) int

	// GetNeighborNodes gets the neighboring nodes.
	//
	// Java: Graph.getNeighborNodes(N)
	GetNeighborNodes(value N) []GraphNode[N, E]

	// GetFirstEdge retrieves any edge from the graph: the first edges
	// between those two values in the graph. null if there are none.
	//
	// Java: Graph.getFirstEdge(N, N)
	GetFirstEdge(n1 N, n2 N) Graph_GraphEdge[N, E]

	// Java: Graph.hasNode(N)
	HasNode(n N) bool

	// IsConnected checks whether two nodes in the graph are connected.
	//
	// Java: Graph.isConnected(N, N)
	IsConnected(n1 N, n2 N) bool

	// IsConnectedWithEdge checks whether two nodes in the graph are
	// connected by the given edge type.
	//
	// Java: Graph.isConnected(N, E, N)
	IsConnectedWithEdge(n1 N, e E, n2 N) bool

	// Java: Graph.getNodeOrFail(N)
	GetNodeOrFail(val N) GraphNode[N, E]

	// Java: Graph.clearEdgeAnnotations()
	ClearEdgeAnnotations()

	// Java: Graph.pushNodeAnnotations()
	PushNodeAnnotations()

	// Java: Graph.popNodeAnnotations()
	PopNodeAnnotations()

	// Java: Graph.pushEdgeAnnotations()
	PushEdgeAnnotations()

	// Java: Graph.popEdgeAnnotations()
	PopEdgeAnnotations()

	GetGraphBase() *GraphBase[N, E]
}

// graph_AnnotationState is a pseudo typedef for a pair of annotations.
// Record of an object's annotation at some state.
type graph_AnnotationState struct {
	first  Annotatable
	second Annotation
}

// graph_GraphAnnotationState is a pseudo typedef for
// ArrayList<AnnotationState>. Record of a collection of objects'
// annotations at some state.
type graph_GraphAnnotationState []graph_AnnotationState

// graph_annotationStack is a Deque<GraphAnnotationState> used as a stack.
type graph_annotationStack struct {
	states []graph_GraphAnnotationState
}

// GraphBase holds Graph's fields and final/base implementations.
type GraphBase[N, E any] struct {
	self Graph[N, E]

	// Used by pushNodeAnnotations() and popNodeAnnotations().
	nodeAnnotationStack *graph_annotationStack

	// Used by pushEdgeAnnotations() and popEdgeAnnotations().
	edgeAnnotationStack *graph_annotationStack
}

// InitGraphBase is the Graph() constructor; self is the graph being
// constructed.
func (g *GraphBase[N, E]) InitGraphBase(self Graph[N, E]) { g.self = self }

// GetGraphBase gives access to the base fields through the interface.
func (g *GraphBase[N, E]) GetGraphBase() *GraphBase[N, E] { return g }

// ConnectIfNotFound connects two nodes in the graph with an edge if such
// edge does not already exists between the nodes. Final in Java.
//
// Java: Graph.connectIfNotFound(N, E, N)
func (g *GraphBase[N, E]) ConnectIfNotFound(n1 N, edge E, n2 N) {
	if !g.self.IsConnectedWithEdge(n1, edge, n2) {
		g.self.Connect(n1, edge, n2)
	}
}

// Java: Graph.getWeight(N)
func (g *GraphBase[N, E]) GetWeight(value N) int {
	return g.self.GetNodeDegree(value)
}

// HasNode checks whether the node exists in the graph (createNode(Object)
// has been called with that value). Final in Java.
//
// Java: Graph.hasNode(N)
func (g *GraphBase[N, E]) HasNode(n N) bool {
	return g.self.GetNode(n) != nil
}

// GetNodeOrFail gets the node of the specified type, or throws an
// IllegalArgumentException.
//
// Go-only: language-adapter — Java's `<T extends GraphNode<N, E>> T` is an
// unchecked cast at the call site; callers type-assert the result.
//
// Java: Graph.getNodeOrFail(N)
func (g *GraphBase[N, E]) GetNodeOrFail(val N) GraphNode[N, E] {
	node := g.self.GetNode(val)
	if node == nil {
		panic(jutil.NewIllegalArgumentExceptionMsg(jutil.ToJavaString(val) + " does not exist in graph"))
	}
	return node
}

// ClearNodeAnnotations is final in Java.
//
// Java: Graph.clearNodeAnnotations()
func (g *GraphBase[N, E]) ClearNodeAnnotations() {
	for _, n := range g.self.GetNodes() {
		n.SetAnnotation(nil)
	}
}

// ClearEdgeAnnotations makes each edge's annotation null. Final in Java.
//
// Java: Graph.clearEdgeAnnotations()
func (g *GraphBase[N, E]) ClearEdgeAnnotations() {
	for _, e := range g.self.GetEdges() {
		e.SetAnnotation(nil)
	}
}

// PushNodeAnnotations pushes nodes' annotation values. Restored with
// popNodeAnnotations(). Nodes' annotation values are cleared. Final in
// Java.
//
// Java: Graph.pushNodeAnnotations()
func (g *GraphBase[N, E]) PushNodeAnnotations() {
	if g.nodeAnnotationStack == nil {
		g.nodeAnnotationStack = &graph_annotationStack{}
	}
	nodes := g.self.GetNodes()
	haveAnnotations := make([]Annotatable, len(nodes))
	for i, n := range nodes {
		haveAnnotations[i] = n
	}
	graph_pushAnnotations(g.nodeAnnotationStack, haveAnnotations)
}

// PopNodeAnnotations restores nodes' annotation values to state before
// last pushNodeAnnotations(). Final in Java.
//
// Java: Graph.popNodeAnnotations()
func (g *GraphBase[N, E]) PopNodeAnnotations() {
	if g.nodeAnnotationStack == nil {
		panic(jutil.NewNullPointerExceptionMsg("Popping node annotations without pushing."))
	}
	graph_popAnnotations(g.nodeAnnotationStack)
}

// PushEdgeAnnotations pushes edges' annotation values. Restored with
// popEdgeAnnotations(). Edges' annotation values are cleared. Final in
// Java.
//
// Java: Graph.pushEdgeAnnotations()
func (g *GraphBase[N, E]) PushEdgeAnnotations() {
	if g.edgeAnnotationStack == nil {
		g.edgeAnnotationStack = &graph_annotationStack{}
	}
	edges := g.self.GetEdges()
	haveAnnotations := make([]Annotatable, len(edges))
	for i, e := range edges {
		haveAnnotations[i] = e
	}
	graph_pushAnnotations(g.edgeAnnotationStack, haveAnnotations)
}

// PopEdgeAnnotations restores edges' annotation values to state before
// last pushEdgeAnnotations(). Final in Java.
//
// Java: Graph.popEdgeAnnotations()
func (g *GraphBase[N, E]) PopEdgeAnnotations() {
	if g.edgeAnnotationStack == nil {
		panic(jutil.NewNullPointerExceptionMsg("Popping edge annotations without pushing."))
	}
	graph_popAnnotations(g.edgeAnnotationStack)
}

// Graph_GraphEdge is a generic edge.
//
// N is the value type that the graph node stores; E is the value type that
// the graph edge stores.
type Graph_GraphEdge[N, E any] interface {
	Annotatable

	// GetValue retrieves the edge's value.
	//
	// Java: Graph.GraphEdge.getValue()
	GetValue() E

	// Java: Graph.GraphEdge.getNodeA()
	GetNodeA() GraphNode[N, E]

	// Java: Graph.GraphEdge.getNodeB()
	GetNodeB() GraphNode[N, E]
}

// Graph_SimpleSubGraph is a simple implementation of SubGraph that
// calculates adjacency by iterating over a node's neighbors.
type Graph_SimpleSubGraph[N, E any] struct {
	graph Graph[N, E]
	nodes []GraphNode[N, E]
}

// Java: Graph.SimpleSubGraph(Graph<N, E>)
func NewGraph_SimpleSubGraph[N, E any](graph Graph[N, E]) *Graph_SimpleSubGraph[N, E] {
	return &Graph_SimpleSubGraph[N, E]{graph: graph}
}

// Java: Graph.SimpleSubGraph.isIndependentOf(N)
func (s *Graph_SimpleSubGraph[N, E]) IsIndependentOf(value N) bool {
	node := s.graph.GetNode(value)
	for _, n := range s.nodes {
		for _, neighbor := range s.graph.GetNeighborNodes(n.GetValue()) {
			// List.contains(Object): GraphNodes do not override equals, so
			// this is identity (interface == on the node pointers).
			if node == neighbor {
				return false
			}
		}
	}
	return true
}

// Java: Graph.SimpleSubGraph.addNode(N)
func (s *Graph_SimpleSubGraph[N, E]) AddNode(value N) {
	s.nodes = append(s.nodes, s.graph.GetNodeOrFail(value))
}

// graph_pushAnnotationFrame pushes or reuses a frame on stack.
func graph_pushAnnotationFrame(stack *graph_annotationStack, capHint int) *graph_GraphAnnotationState {
	idx := len(stack.states)
	if idx < cap(stack.states) {
		stack.states = stack.states[:idx+1]
		top := &stack.states[idx]
		*top = (*top)[:0]
		return top
	}
	stack.states = append(stack.states, make(graph_GraphAnnotationState, 0, capHint))
	return &stack.states[idx]
}

// graph_pushAnnotations pushes a new list on stack and stores nodes
// annotations in the new list. Clears objects' annotations as well.
//
// Java: Graph.pushAnnotations(Deque<GraphAnnotationState>, Collection<? extends Annotatable>)
func graph_pushAnnotations(stack *graph_annotationStack, haveAnnotations []Annotatable) {
	top := graph_pushAnnotationFrame(stack, len(haveAnnotations))
	for _, h := range haveAnnotations {
		*top = append(*top, graph_AnnotationState{first: h, second: h.GetAnnotation()})
		h.SetAnnotation(nil)
	}
}

// graph_popAnnotations restores the node annotations on the top of stack
// and pops stack.
//
// Java: Graph.popAnnotations(Deque<GraphAnnotationState>)
func graph_popAnnotations(stack *graph_annotationStack) {
	if len(stack.states) == 0 {
		panic(jutil.NewNoSuchElementException()) // ArrayDeque.pop()
	}
	topIdx := len(stack.states) - 1
	top := stack.states[topIdx]
	for i := range top {
		top[i].first.SetAnnotation(top[i].second)
		top[i] = graph_AnnotationState{}
	}
	stack.states[topIdx] = top[:0]
	stack.states = stack.states[:topIdx]
}
