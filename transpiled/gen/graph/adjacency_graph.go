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
// Port of com/google/javascript/jscomp/graph/AdjacencyGraph.java (generated mechanically by j2go).

package graph

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java interface: com.google.javascript.jscomp.graph.AdjacencyGraph
type AdjacencyGraph[N comparable, E any] interface {
	Implements_graph_AdjacencyGraph()
	IsJavaNil_() bool
	GetNodes() j2gort.Ext_java_util_Collection
	GetNodeCount() int32
	GetNode(N) any
	NewSubGraph() SubGraph[N, E]
	ClearNodeAnnotations()
	GetWeight(N) int32
}
func ToIface_AdjacencyGraph[N comparable, E any, T_iface_ interface {
	AdjacencyGraph[N, E]
	comparable
}](v T_iface_) AdjacencyGraph[N, E] {
	var z T_iface_
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.graph.AdjacencyGraph", "com.google.javascript.jscomp.graph.AdjacencyGraph", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			Implements_graph_AdjacencyGraph()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}
