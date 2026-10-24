// Port of com/google/javascript/jscomp/graph/SubGraph.java @ closure-compiler v20261021

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

// SubGraph is an interface representing a subgraph that provides adjacency
// calculation to a node.
//
// N is the value type that the graph node stores; E is the value type that
// the graph edge stores.
type SubGraph[N, E any] interface {
	// IsIndependentOf returns true if the node is a neighbor of any node in
	// this SubGraph.
	//
	// Java: SubGraph.isIndependentOf(N)
	IsIndependentOf(node N) bool

	// AddNode adds the node into this subgraph.
	//
	// Java: SubGraph.addNode(N)
	AddNode(value N)
}
