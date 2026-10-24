// Port of com/google/javascript/jscomp/graph/GraphNode.java @ closure-compiler v20261021

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

// GraphNode is a generic node.
//
// N is the value type that the graph node stores; E is the value type that
// the graph edge stores.
type GraphNode[N, E any] interface {
	Annotatable

	// GetValue retrieves the node's value.
	//
	// Java: GraphNode.getValue()
	GetValue() N
}
