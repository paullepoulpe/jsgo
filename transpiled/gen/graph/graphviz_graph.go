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
// Port of com/google/javascript/jscomp/graph/GraphvizGraph.java (generated mechanically by j2go).

package graph

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java interface: com.google.javascript.jscomp.graph.GraphvizGraph
type GraphvizGraph interface {
	Implements_graph_GraphvizGraph()
	IsJavaNil_() bool
	GetName() string
	IsDirected() bool
	GetGraphvizNodes() j2gort.Ext_java_util_List
	GetGraphvizEdges() j2gort.Ext_java_util_List
}
func ToIface_GraphvizGraph[T interface {
	GraphvizGraph
	comparable
}](v T) GraphvizGraph {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_GraphvizGraph[T GraphvizGraph](v T) GraphvizGraph {
	var vi GraphvizGraph = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.graph.GraphvizGraph", "com.google.javascript.jscomp.graph.GraphvizGraph", func(v_ any) bool {
		return j2gort.InstanceOf[GraphvizGraph](v_)
	})
}
// Java interface: com.google.javascript.jscomp.graph.GraphvizGraph.GraphvizNode
type GraphvizGraph_GraphvizNode interface {
	Implements_graph_GraphvizGraph_GraphvizNode()
	IsJavaNil_() bool
	GetId() string
	GetColor() string
	GetLabel() string
}
func ToIface_GraphvizGraph_GraphvizNode[T interface {
	GraphvizGraph_GraphvizNode
	comparable
}](v T) GraphvizGraph_GraphvizNode {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_GraphvizGraph_GraphvizNode[T GraphvizGraph_GraphvizNode](v T) GraphvizGraph_GraphvizNode {
	var vi GraphvizGraph_GraphvizNode = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.graph.GraphvizGraph.GraphvizNode", "com.google.javascript.jscomp.graph.GraphvizGraph$GraphvizNode", func(v_ any) bool {
		return j2gort.InstanceOf[GraphvizGraph_GraphvizNode](v_)
	})
}
// Java interface: com.google.javascript.jscomp.graph.GraphvizGraph.GraphvizEdge
type GraphvizGraph_GraphvizEdge interface {
	Implements_graph_GraphvizGraph_GraphvizEdge()
	IsJavaNil_() bool
	GetNode1Id() string
	GetNode2Id() string
	GetColor() string
	GetLabel() string
}
func ToIface_GraphvizGraph_GraphvizEdge[T interface {
	GraphvizGraph_GraphvizEdge
	comparable
}](v T) GraphvizGraph_GraphvizEdge {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_GraphvizGraph_GraphvizEdge[T GraphvizGraph_GraphvizEdge](v T) GraphvizGraph_GraphvizEdge {
	var vi GraphvizGraph_GraphvizEdge = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.graph.GraphvizGraph.GraphvizEdge", "com.google.javascript.jscomp.graph.GraphvizGraph$GraphvizEdge", func(v_ any) bool {
		return j2gort.InstanceOf[GraphvizGraph_GraphvizEdge](v_)
	})
}
