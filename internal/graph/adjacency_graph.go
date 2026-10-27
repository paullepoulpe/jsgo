// Port of com/google/javascript/jscomp/graph/AdjacencyGraph.java @ closure-compiler v20261021

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

// AdjacencyGraph is a minimal graph interface. Provided is add nodes to the
// graph, adjacency calculation between a SubGraph and a GraphNode, and
// adding node annotations.
//
// For a more extensive interface, see Graph.
//
// N is the value type that the graph node stores; E is the value type that
// the graph edge stores.
type AdjacencyGraph[N, E any] interface {
	// GetNodes gets an immutable list of all nodes.
	//
	// Java: AdjacencyGraph.getNodes()
	GetNodes() []GraphNode[N, E]

	// Java: AdjacencyGraph.getNodeCount()
	GetNodeCount() int

	// GetNode gets a node from the graph given a value. Values equality are
	// compared using Object.equals.
	//
	// Returns the corresponding node in the graph, null if there value has
	// no corresponding node.
	//
	// Java: AdjacencyGraph.getNode(N)
	GetNode(value N) GraphNode[N, E]

	// NewSubGraph returns an empty SubGraph for this Graph.
	//
	// Java: AdjacencyGraph.newSubGraph()
	NewSubGraph() SubGraph[N, E]

	// ClearNodeAnnotations makes each node's annotation null.
	//
	// Java: AdjacencyGraph.clearNodeAnnotations()
	ClearNodeAnnotations()

	// GetWeight returns a weight for the given value to be used in ordering
	// nodes, e.g. in GraphColoring.
	//
	// Java: AdjacencyGraph.getWeight(N)
	GetWeight(value N) int
}
