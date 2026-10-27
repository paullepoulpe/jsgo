// Port of com/google/javascript/jscomp/graph/DiGraph.java @ closure-compiler v20261021

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

// DiGraph is a generic directed graph.
//
// N is the value type that the graph node stores; E is the value type that
// the graph edge stores.
//
// Go-only: language-adapter — DiGraph's covariant overrides of Graph's
// createNode, getNodes, getNode, getEdges() and getEdges(N, N) are the
// methods CreateDiGraphNode, GetDiGraphNodes, GetDiGraphNode,
// GetDiGraphEdges and GetDiGraphEdgesBetween (see Graph).
type DiGraph[N, E any] interface {
	Graph[N, E]

	// Java: DiGraph.createNode(N)
	CreateDiGraphNode(nodeValue N) DiGraph_DiGraphNode[N, E]

	// Java: DiGraph.getNodes()
	GetDiGraphNodes() []DiGraph_DiGraphNode[N, E]

	// Java: DiGraph.getNode(N)
	GetDiGraphNode(nodeValue N) DiGraph_DiGraphNode[N, E]

	// Java: DiGraph.getEdges()
	GetDiGraphEdges() []DiGraph_DiGraphEdge[N, E]

	// Java: DiGraph.getEdges(N, N)
	GetDiGraphEdgesBetween(n1 N, n2 N) []DiGraph_DiGraphEdge[N, E]

	// Java: DiGraph.getEdgesInDirection(N, N)
	GetEdgesInDirection(n1 N, n2 N) []DiGraph_DiGraphEdge[N, E]

	// GetOutEdges gets an immutable list of out edges of the given node.
	//
	// Java: DiGraph.getOutEdges(N)
	GetOutEdges(nodeValue N) []*LinkedDirectedGraph_LinkedDiGraphEdge[N, E]

	// GetInEdges gets an immutable list of in edges of the given node.
	//
	// Java: DiGraph.getInEdges(N)
	GetInEdges(nodeValue N) []*LinkedDirectedGraph_LinkedDiGraphEdge[N, E]

	// Java: DiGraph.getDirectedPredNodes(DiGraphNode<N, E>)
	GetDirectedPredNodes(n DiGraph_DiGraphNode[N, E]) []DiGraph_DiGraphNode[N, E]

	// Java: DiGraph.getDirectedPredNodes(N)
	GetDirectedPredNodesOfValue(nodeValue N) []DiGraph_DiGraphNode[N, E]

	// Java: DiGraph.getDirectedSuccNodes(DiGraphNode<N, E>)
	GetDirectedSuccNodes(n DiGraph_DiGraphNode[N, E]) []DiGraph_DiGraphNode[N, E]

	// Java: DiGraph.getDirectedSuccNodes(N)
	GetDirectedSuccNodesOfValue(nodeValue N) []DiGraph_DiGraphNode[N, E]

	// DisconnectInDirection disconnects all edges from n1 to n2.
	//
	// Java: DiGraph.disconnectInDirection(N, N)
	DisconnectInDirection(n1 N, n2 N)

	// IsConnectedInDirection checks whether two nodes in the graph are
	// connected via a directed edge: true if the graph contains edge from n1
	// to n2.
	//
	// Java: DiGraph.isConnectedInDirection(N, N)
	IsConnectedInDirection(n1 N, n2 N) bool

	// IsConnectedInDirectionWithEdge checks whether two nodes in the graph
	// are connected via a directed edge with the given value.
	//
	// Java: DiGraph.isConnectedInDirection(N, E, N)
	IsConnectedInDirectionWithEdge(n1 N, edgeValue E, n2 N) bool
}

// DiGraphBase holds DiGraph's base implementations.
type DiGraphBase[N, E any] struct {
	GraphBase[N, E]

	diSelf DiGraph[N, E]
}

// InitDiGraphBase is the DiGraph() constructor; self is the graph being
// constructed.
func (g *DiGraphBase[N, E]) InitDiGraphBase(self DiGraph[N, E]) {
	g.InitGraphBase(self)
	g.diSelf = self
}

// Java: DiGraph.isConnected(N, N)
func (g *DiGraphBase[N, E]) IsConnected(n1 N, n2 N) bool {
	return g.diSelf.IsConnectedInDirection(n1, n2) || g.diSelf.IsConnectedInDirection(n2, n1)
}

// Java: DiGraph.isConnected(N, E, N)
func (g *DiGraphBase[N, E]) IsConnectedWithEdge(n1 N, e E, n2 N) bool {
	return g.diSelf.IsConnectedInDirectionWithEdge(n1, e, n2) ||
		g.diSelf.IsConnectedInDirectionWithEdge(n2, e, n1)
}

// DiGraph_DiGraphNode is a generic directed graph node.
//
// N is the value type that the graph node stores; E is the value type that
// the graph edge stores.
//
// Go-only: language-adapter — the edge lists are returned as slices that
// alias the node's lists (no copy, as Java returns the lists themselves);
// callers must not modify them or keep them across graph mutations. They
// hold the edges of LinkedDirectedGraph, the only DiGraph: a slice of
// interfaces would take twice the memory in every graph node.
type DiGraph_DiGraphNode[N, E any] interface {
	GraphNode[N, E]

	// Java: DiGraph.DiGraphNode.getOutEdges()
	GetOutEdges() []*LinkedDirectedGraph_LinkedDiGraphEdge[N, E]

	// Java: DiGraph.DiGraphNode.getInEdges()
	GetInEdges() []*LinkedDirectedGraph_LinkedDiGraphEdge[N, E]

	// HasPriority returns whether a priority has been set.
	//
	// Java: DiGraph.DiGraphNode.hasPriority()
	HasPriority() bool

	// GetPriority returns a nonnegative integer priority which can be used
	// to order nodes.
	//
	// Throws if a priority has not been set.
	//
	// Java: DiGraph.DiGraphNode.getPriority()
	GetPriority() int

	// SetPriority sets a node priority, must be non-negative.
	//
	// Java: DiGraph.DiGraphNode.setPriority(int)
	SetPriority(priority int)
}

// DiGraph_DiGraphEdge is a generic directed graph edge.
//
// N is the value type that the graph node stores; E is the value type that
// the graph edge stores.
type DiGraph_DiGraphEdge[N, E any] interface {
	Graph_GraphEdge[N, E]

	// Java: DiGraph.DiGraphEdge.getSource()
	GetSource() DiGraph_DiGraphNode[N, E]

	// Java: DiGraph.DiGraphEdge.getDestination()
	GetDestination() DiGraph_DiGraphNode[N, E]
}
