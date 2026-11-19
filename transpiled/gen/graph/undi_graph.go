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
// Port of com/google/javascript/jscomp/graph/UndiGraph.java (generated mechanically by j2go).

package graph

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java polymorphic class interface: com.google.javascript.jscomp.graph.UndiGraph
type UndiGraph[N comparable, E any] interface {
	Implements_graph_UndiGraph()
	IsJavaNil_() bool
	GetUndiGraphBase() *UndiGraphBase[N, E]
	Graph[N, E]
	GetUndirectedGraphNodes() j2gort.Ext_java_util_Collection
	CreateUndirectedGraphNode(N) UndiGraph_UndiGraphNode[N, E]
	GetUndirectedGraphNode(N) UndiGraph_UndiGraphNode[N, E]
	GetUndirectedGraphEdges(N, N) j2gort.Ext_java_util_List
}
func ToIface_UndiGraph[N comparable, E any, T_iface_ interface {
	UndiGraph[N, E]
	comparable
}](v T_iface_) UndiGraph[N, E] {
	var z T_iface_
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.jscomp.graph.UndiGraph
type UndiGraphBase[N comparable, E any] struct {
	GraphBase[N, E]
}

func (this_ *UndiGraphBase[N, E]) BindSelf_(s any) {
	_ = s
	this_.GraphBase.BindSelf_(s)
}
func (this_ *UndiGraphBase[N, E]) IsJavaNil_() bool { return this_ == nil }
func (*UndiGraphBase[N, E]) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.graph.UndiGraph", "UndiGraph" }
func (this_ *UndiGraphBase[N, E]) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.GraphBase.IdentityHash_()
}
func (*UndiGraphBase[N, E]) IsGenericInst_graph_UndiGraph_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.graph.UndiGraph", "com.google.javascript.jscomp.graph.UndiGraph", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			IsGenericInst_graph_UndiGraph_()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}
func (b *UndiGraphBase[N, E]) GetUndiGraphBase() *UndiGraphBase[N, E] { return b }
func (*UndiGraphBase[N, E]) Implements_graph_UndiGraph() {}
func (*UndiGraphBase[N, E]) Implements_graph_Graph() {}
func (*UndiGraphBase[N, E]) Implements_graph_AdjacencyGraph() {}

func (this_ *UndiGraphBase[N, E]) initFields_UndiGraph() {
	_ = this_
}

func InitUndiGraph[N comparable, E any](this_ *UndiGraphBase[N, E]) {
	InitGraph[N, E](&this_.GraphBase)
	this_.initFields_UndiGraph()
}

// Java method: UndiGraph.getUndirectedGraphNodes()
func (this_ *UndiGraphBase[N, E]) GetUndirectedGraphNodes() j2gort.Ext_java_util_Collection {
	panic("abstract or native method invoked")
}

// Java method: UndiGraph.createUndirectedGraphNode(N nodeValue)
func (this_ *UndiGraphBase[N, E]) CreateUndirectedGraphNode(nodeValue N) UndiGraph_UndiGraphNode[N, E] {
	_ = nodeValue
	panic("abstract or native method invoked")
}

// Java method: UndiGraph.getUndirectedGraphNode(N nodeValue)
func (this_ *UndiGraphBase[N, E]) GetUndirectedGraphNode(nodeValue N) UndiGraph_UndiGraphNode[N, E] {
	_ = nodeValue
	panic("abstract or native method invoked")
}

// Java method: UndiGraph.getUndirectedGraphEdges(N n1, N n2)
func (this_ *UndiGraphBase[N, E]) GetUndirectedGraphEdges(n1 N, n2 N) j2gort.Ext_java_util_List {
	_ = n1
	_ = n2
	panic("abstract or native method invoked")
}

// Java interface: com.google.javascript.jscomp.graph.UndiGraph.UndiGraphNode
type UndiGraph_UndiGraphNode[N comparable, E any] interface {
	Implements_graph_UndiGraph_UndiGraphNode()
	IsJavaNil_() bool
	GraphNode[N, E]
	GetNeighborEdges() j2gort.Ext_java_util_List
	GetNeighborEdgesIterator() j2gort.Ext_java_util_Iterator
}
func ToIface_UndiGraph_UndiGraphNode[N comparable, E any, T_iface_ interface {
	UndiGraph_UndiGraphNode[N, E]
	comparable
}](v T_iface_) UndiGraph_UndiGraphNode[N, E] {
	var z T_iface_
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.graph.UndiGraph.UndiGraphNode", "com.google.javascript.jscomp.graph.UndiGraph$UndiGraphNode", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			Implements_graph_UndiGraph_UndiGraphNode()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}
// Java interface: com.google.javascript.jscomp.graph.UndiGraph.UndiGraphEdge
type UndiGraph_UndiGraphEdge[N comparable, E any] interface {
	Implements_graph_UndiGraph_UndiGraphEdge()
	IsJavaNil_() bool
	Graph_GraphEdge[N, E]
}
func ToIface_UndiGraph_UndiGraphEdge[N comparable, E any, T_iface_ interface {
	UndiGraph_UndiGraphEdge[N, E]
	comparable
}](v T_iface_) UndiGraph_UndiGraphEdge[N, E] {
	var z T_iface_
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.graph.UndiGraph.UndiGraphEdge", "com.google.javascript.jscomp.graph.UndiGraph$UndiGraphEdge", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			Implements_graph_UndiGraph_UndiGraphEdge()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}
