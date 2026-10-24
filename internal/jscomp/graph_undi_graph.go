// Port of com/google/javascript/jscomp/graph/UndiGraph.java @ closure-compiler v20261021

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
)

// UndiGraph is a generic undirected graph.
//
// N is the value type that the graph node stores; E is the value type that
// the graph edge stores.
type UndiGraph[N, E any] interface {
	Graph[N, E]

	// GetUndirectedGraphNodes gets an immutable collection of all the nodes
	// in this graph.
	//
	// Java: UndiGraph.getUndirectedGraphNodes()
	GetUndirectedGraphNodes() []UndiGraph_UndiGraphNode[N, E]

	// Java: UndiGraph.createUndirectedGraphNode(N)
	CreateUndirectedGraphNode(nodeValue N) UndiGraph_UndiGraphNode[N, E]

	// Java: UndiGraph.getUndirectedGraphNode(N)
	GetUndirectedGraphNode(nodeValue N) UndiGraph_UndiGraphNode[N, E]

	// GetUndirectedGraphEdges returns the edges between n1 and n2; nil
	// (Java null) if either node is not in the graph.
	//
	// Java: UndiGraph.getUndirectedGraphEdges(N, N)
	GetUndirectedGraphEdges(n1 N, n2 N) []UndiGraph_UndiGraphEdge[N, E]
}

// UndiGraphBase holds UndiGraph's base implementations (only Graph's).
type UndiGraphBase[N, E any] struct {
	GraphBase[N, E]
}

// InitUndiGraphBase is the UndiGraph() constructor; self is the graph
// being constructed.
func (g *UndiGraphBase[N, E]) InitUndiGraphBase(self UndiGraph[N, E]) {
	g.InitGraphBase(self)
}

// UndiGraph_UndiGraphNode is a generic undirected graph node.
//
// N is the value type that the graph node stores; E is the value type that
// the graph edge stores.
type UndiGraph_UndiGraphNode[N, E any] interface {
	GraphNode[N, E]

	// GetNeighborEdges returns the node's edge list.
	//
	// Go-only: language-adapter — the slice aliases the node's list (see
	// DiGraph_DiGraphNode).
	//
	// Java: UndiGraph.UndiGraphNode.getNeighborEdges()
	GetNeighborEdges() []UndiGraph_UndiGraphEdge[N, E]

	// GetNeighborEdgesIterator iterates the node's edges.
	//
	// Go-only: language-adapter — the Java Iterator is an iter.Seq (no
	// remove()).
	//
	// Java: UndiGraph.UndiGraphNode.getNeighborEdgesIterator()
	GetNeighborEdgesIterator() iter.Seq[UndiGraph_UndiGraphEdge[N, E]]
}

// UndiGraph_UndiGraphEdge is a generic undirected graph edge.
//
// N is the value type that the graph node stores; E is the value type that
// the graph edge stores.
type UndiGraph_UndiGraphEdge[N, E any] interface {
	Graph_GraphEdge[N, E]
}
