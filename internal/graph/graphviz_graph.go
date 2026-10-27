// Port of com/google/javascript/jscomp/graph/GraphvizGraph.java @ closure-compiler v20261021

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

// GraphvizGraph is a graph that can be dumped to a Graphviz DOT file.
//
// An object which can be visualized as a graph should implement this
// interface. The com.google.javascript.jscomp.DotFormatter#toDot function
// can be used to get a visualization of the object for debugging purpose.
type GraphvizGraph interface {
	// GetName is the name of the graph.
	//
	// Java: GraphvizGraph.getName()
	GetName() string

	// IsDirected returns true if the graph is a directed graph.
	//
	// Java: GraphvizGraph.isDirected()
	IsDirected() bool

	// GetGraphvizNodes retrieves a list of nodes in the graph.
	//
	// Java: GraphvizGraph.getGraphvizNodes()
	GetGraphvizNodes() []GraphvizGraph_GraphvizNode

	// GetGraphvizEdges retrieves a list of edges in the graph.
	//
	// Java: GraphvizGraph.getGraphvizEdges()
	GetGraphvizEdges() []GraphvizGraph_GraphvizEdge
}

// GraphvizGraph_GraphvizNode is a Graphviz node.
type GraphvizGraph_GraphvizNode interface {
	// GetId retrieves the unique ID.
	//
	// Java: GraphvizGraph.GraphvizNode.getId()
	GetId() string

	// GetColor retrieves color of the node.
	//
	// Java: GraphvizGraph.GraphvizNode.getColor()
	GetColor() string

	// GetLabel retrieves the label of the node.
	//
	// Java: GraphvizGraph.GraphvizNode.getLabel()
	GetLabel() string
}

// GraphvizGraph_GraphvizEdge is a Graphviz edge.
type GraphvizGraph_GraphvizEdge interface {
	// GetNode1Id gets the first node in the edge. In a directed node, this
	// will be the source node.
	//
	// Java: GraphvizGraph.GraphvizEdge.getNode1Id()
	GetNode1Id() string

	// GetNode2Id gets the second node in the edge. In a directed node, this
	// will be the destination node.
	//
	// Java: GraphvizGraph.GraphvizEdge.getNode2Id()
	GetNode2Id() string

	// GetColor retrieves color of the edge.
	//
	// Java: GraphvizGraph.GraphvizEdge.getColor()
	GetColor() string

	// GetLabel retrieves the label of the edge.
	//
	// Java: GraphvizGraph.GraphvizEdge.getLabel()
	GetLabel() string
}
