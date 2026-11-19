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
// Port of com/google/javascript/jscomp/graph/DiGraph.java (generated mechanically by j2go).

package graph

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java polymorphic class interface: com.google.javascript.jscomp.graph.DiGraph
type DiGraph[N comparable, E any] interface {
	Implements_graph_DiGraph()
	IsJavaNil_() bool
	GetDiGraphBase() *DiGraphBase[N, E]
	Graph[N, E]
	GetEdgesInDirection(N, N) j2gort.Ext_java_util_List
	GetOutEdges(N) j2gort.Ext_java_util_List
	GetInEdges(N) j2gort.Ext_java_util_List
	GetDirectedPredNodes_1_DiGraph_DiGraphNode_any_any_(DiGraph_DiGraphNode[N, E]) j2gort.Ext_java_util_List
	GetDirectedPredNodes_1_any(N) j2gort.Ext_java_util_List
	GetDirectedSuccNodes_1_DiGraph_DiGraphNode_any_any_(DiGraph_DiGraphNode[N, E]) j2gort.Ext_java_util_List
	GetDirectedSuccNodes_1_any(N) j2gort.Ext_java_util_List
	DisconnectInDirection(N, N)
	IsConnectedInDirection_2(N, N) bool
	IsConnectedInDirection_3_any_any_any(N, E, N) bool
}
func ToIface_DiGraph[N comparable, E any, T_iface_ interface {
	DiGraph[N, E]
	comparable
}](v T_iface_) DiGraph[N, E] {
	var z T_iface_
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.jscomp.graph.DiGraph
type DiGraphBase[N comparable, E any] struct {
	GraphBase[N, E]
}

func (this_ *DiGraphBase[N, E]) BindSelf_(s any) {
	_ = s
	this_.GraphBase.BindSelf_(s)
}
func (this_ *DiGraphBase[N, E]) IsJavaNil_() bool { return this_ == nil }
func (*DiGraphBase[N, E]) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.graph.DiGraph", "DiGraph" }
func (this_ *DiGraphBase[N, E]) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.GraphBase.IdentityHash_()
}
func (*DiGraphBase[N, E]) IsGenericInst_graph_DiGraph_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.graph.DiGraph", "com.google.javascript.jscomp.graph.DiGraph", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			IsGenericInst_graph_DiGraph_()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}
func (b *DiGraphBase[N, E]) GetDiGraphBase() *DiGraphBase[N, E] { return b }
func (this_ *DiGraphBase[N, E]) self_DiGraph_() DiGraph[N, E] {
	if s := this_.GraphBase.self_; s != nil {
		return s.(DiGraph[N, E])
	}
	return nil
}
func (*DiGraphBase[N, E]) Implements_graph_DiGraph() {}
func (*DiGraphBase[N, E]) Implements_graph_Graph() {}
func (*DiGraphBase[N, E]) Implements_graph_AdjacencyGraph() {}

func (this_ *DiGraphBase[N, E]) initFields_DiGraph() {
	_ = this_
}

func InitDiGraph[N comparable, E any](this_ *DiGraphBase[N, E]) {
	InitGraph[N, E](&this_.GraphBase)
	this_.initFields_DiGraph()
}

// Java method: DiGraph.createNode(N nodeValue)
func (this_ *DiGraphBase[N, E]) CreateNode(nodeValue N) any {
	_ = nodeValue
	panic("abstract or native method invoked")
}

// Java method: DiGraph.getNodes()
func (this_ *DiGraphBase[N, E]) GetNodes() j2gort.Ext_java_util_Collection {
	panic("abstract or native method invoked")
}

// Java method: DiGraph.getNode(N nodeValue)
func (this_ *DiGraphBase[N, E]) GetNode(nodeValue N) any {
	_ = nodeValue
	panic("abstract or native method invoked")
}

// Java method: DiGraph.getEdges()
func (this_ *DiGraphBase[N, E]) GetEdges_0() j2gort.Ext_java_util_List {
	panic("abstract or native method invoked")
}

// Java method: DiGraph.getEdges(N n1, N n2)
func (this_ *DiGraphBase[N, E]) GetEdges_2(n1 N, n2 N) j2gort.Ext_java_util_List {
	_ = n1
	_ = n2
	panic("abstract or native method invoked")
}

// Java method: DiGraph.getEdgesInDirection(N n1, N n2)
func (this_ *DiGraphBase[N, E]) GetEdgesInDirection(n1 N, n2 N) j2gort.Ext_java_util_List {
	_ = n1
	_ = n2
	panic("abstract or native method invoked")
}

// Java method: DiGraph.getOutEdges(N nodeValue)
func (this_ *DiGraphBase[N, E]) GetOutEdges(nodeValue N) j2gort.Ext_java_util_List {
	_ = nodeValue
	panic("abstract or native method invoked")
}

// Java method: DiGraph.getInEdges(N nodeValue)
func (this_ *DiGraphBase[N, E]) GetInEdges(nodeValue N) j2gort.Ext_java_util_List {
	_ = nodeValue
	panic("abstract or native method invoked")
}

// Java method: DiGraph.getDirectedPredNodes(com.google.javascript.jscomp.graph.DiGraph.DiGraphNode<N,E> n)
func (this_ *DiGraphBase[N, E]) GetDirectedPredNodes_1_DiGraph_DiGraphNode_any_any_(n DiGraph_DiGraphNode[N, E]) j2gort.Ext_java_util_List {
	_ = n
	panic("abstract or native method invoked")
}

// Java method: DiGraph.getDirectedPredNodes(N nodeValue)
func (this_ *DiGraphBase[N, E]) GetDirectedPredNodes_1_any(nodeValue N) j2gort.Ext_java_util_List {
	_ = nodeValue
	panic("abstract or native method invoked")
}

// Java method: DiGraph.getDirectedSuccNodes(com.google.javascript.jscomp.graph.DiGraph.DiGraphNode<N,E> n)
func (this_ *DiGraphBase[N, E]) GetDirectedSuccNodes_1_DiGraph_DiGraphNode_any_any_(n DiGraph_DiGraphNode[N, E]) j2gort.Ext_java_util_List {
	_ = n
	panic("abstract or native method invoked")
}

// Java method: DiGraph.getDirectedSuccNodes(N nodeValue)
func (this_ *DiGraphBase[N, E]) GetDirectedSuccNodes_1_any(nodeValue N) j2gort.Ext_java_util_List {
	_ = nodeValue
	panic("abstract or native method invoked")
}

// Java method: DiGraph.disconnectInDirection(N n1, N n2)
func (this_ *DiGraphBase[N, E]) DisconnectInDirection(n1 N, n2 N) {
	_ = n1
	_ = n2
	panic("abstract or native method invoked")
}

// Java method: DiGraph.isConnectedInDirection(N n1, N n2)
func (this_ *DiGraphBase[N, E]) IsConnectedInDirection_2(n1 N, n2 N) bool {
	_ = n1
	_ = n2
	panic("abstract or native method invoked")
}

// Java method: DiGraph.isConnectedInDirection(N n1, E edgeValue, N n2)
func (this_ *DiGraphBase[N, E]) IsConnectedInDirection_3_any_any_any(n1 N, edgeValue E, n2 N) bool {
	_ = n1
	_ = edgeValue
	_ = n2
	panic("abstract or native method invoked")
}

// Java method: DiGraph.isConnected(N n1, N n2)
func (this_ *DiGraphBase[N, E]) IsConnected_2(n1 N, n2 N) bool {
	_ = n1
	_ = n2
	return (this_.self_DiGraph_().IsConnectedInDirection_2(n1, n2) || this_.self_DiGraph_().IsConnectedInDirection_2(n2, n1))
}

// Java method: DiGraph.isConnected(N n1, E e, N n2)
func (this_ *DiGraphBase[N, E]) IsConnected_3_any_any_any(n1 N, e E, n2 N) bool {
	_ = n1
	_ = e
	_ = n2
	return (this_.self_DiGraph_().IsConnectedInDirection_3_any_any_any(n1, e, n2) || this_.self_DiGraph_().IsConnectedInDirection_3_any_any_any(n2, e, n1))
}

// Java interface: com.google.javascript.jscomp.graph.DiGraph.DiGraphNode
type DiGraph_DiGraphNode[N comparable, E any] interface {
	Implements_graph_DiGraph_DiGraphNode()
	IsJavaNil_() bool
	GraphNode[N, E]
	GetOutEdges() j2gort.Ext_java_util_List
	GetInEdges() j2gort.Ext_java_util_List
	GetOutEdgeCount() int32
	GetOutEdgeAt(int32) DiGraph_DiGraphEdge[N, E]
	GetInEdgeCount() int32
	GetInEdgeAt(int32) DiGraph_DiGraphEdge[N, E]
	HasPriority() bool
	GetPriority() int32
	SetPriority(int32)
	GetAstPosition() int32
	SetAstPosition(int32)
	IsInDfaQueue() bool
	SetInDfaQueue(bool)
}
func ToIface_DiGraph_DiGraphNode[N comparable, E any, T_iface_ interface {
	DiGraph_DiGraphNode[N, E]
	comparable
}](v T_iface_) DiGraph_DiGraphNode[N, E] {
	var z T_iface_
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.graph.DiGraph.DiGraphNode", "com.google.javascript.jscomp.graph.DiGraph$DiGraphNode", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			Implements_graph_DiGraph_DiGraphNode()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}
// Java default interface method: DiGraphNode.getOutEdgeCount()
func IfaceDefault_DiGraph_DiGraphNode_GetOutEdgeCount[N comparable, E any](this_ DiGraph_DiGraphNode[N, E]) int32 {
	_ = this_
	return this_.GetOutEdges().Size()
}

// Java default interface method: DiGraphNode.getOutEdgeAt(int index)
func IfaceDefault_DiGraph_DiGraphNode_GetOutEdgeAt[N comparable, E any](this_ DiGraph_DiGraphNode[N, E], index int32) DiGraph_DiGraphEdge[N, E] {
	_ = this_
	_ = index
	return j2gort.AsIface[DiGraph_DiGraphEdge[N, E]](this_.GetOutEdges().Get(index))
}

// Java default interface method: DiGraphNode.getInEdgeCount()
func IfaceDefault_DiGraph_DiGraphNode_GetInEdgeCount[N comparable, E any](this_ DiGraph_DiGraphNode[N, E]) int32 {
	_ = this_
	return this_.GetInEdges().Size()
}

// Java default interface method: DiGraphNode.getInEdgeAt(int index)
func IfaceDefault_DiGraph_DiGraphNode_GetInEdgeAt[N comparable, E any](this_ DiGraph_DiGraphNode[N, E], index int32) DiGraph_DiGraphEdge[N, E] {
	_ = this_
	_ = index
	return j2gort.AsIface[DiGraph_DiGraphEdge[N, E]](this_.GetInEdges().Get(index))
}

// Java default interface method: DiGraphNode.getAstPosition()
func IfaceDefault_DiGraph_DiGraphNode_GetAstPosition[N comparable, E any](this_ DiGraph_DiGraphNode[N, E]) int32 {
	_ = this_
	return int32(-1)
}

// Java default interface method: DiGraphNode.setAstPosition(int astPosition)
func IfaceDefault_DiGraph_DiGraphNode_SetAstPosition[N comparable, E any](this_ DiGraph_DiGraphNode[N, E], astPosition int32) {
	_ = this_
	_ = astPosition
}

// Java default interface method: DiGraphNode.isInDfaQueue()
func IfaceDefault_DiGraph_DiGraphNode_IsInDfaQueue[N comparable, E any](this_ DiGraph_DiGraphNode[N, E]) bool {
	_ = this_
	return false
}

// Java default interface method: DiGraphNode.setInDfaQueue(boolean inDfaQueue)
func IfaceDefault_DiGraph_DiGraphNode_SetInDfaQueue[N comparable, E any](this_ DiGraph_DiGraphNode[N, E], inDfaQueue bool) {
	_ = this_
	_ = inDfaQueue
}

// Java interface: com.google.javascript.jscomp.graph.DiGraph.DiGraphEdge
type DiGraph_DiGraphEdge[N comparable, E any] interface {
	Implements_graph_DiGraph_DiGraphEdge()
	IsJavaNil_() bool
	Graph_GraphEdge[N, E]
	GetSource() DiGraph_DiGraphNode[N, E]
	GetDestination() DiGraph_DiGraphNode[N, E]
}
func ToIface_DiGraph_DiGraphEdge[N comparable, E any, T_iface_ interface {
	DiGraph_DiGraphEdge[N, E]
	comparable
}](v T_iface_) DiGraph_DiGraphEdge[N, E] {
	var z T_iface_
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.graph.DiGraph.DiGraphEdge", "com.google.javascript.jscomp.graph.DiGraph$DiGraphEdge", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			Implements_graph_DiGraph_DiGraphEdge()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}
