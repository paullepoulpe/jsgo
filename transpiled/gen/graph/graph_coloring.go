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
// Port of com/google/javascript/jscomp/graph/GraphColoring.java (generated mechanically by j2go).

package graph

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java polymorphic class interface: com.google.javascript.jscomp.graph.GraphColoring
type GraphColoring[N comparable, E any] interface {
	Implements_graph_GraphColoring()
	IsJavaNil_() bool
	GetGraphColoringBase() *GraphColoringBase[N, E]
	Color() int32
	GetPartitionSuperNode(N) N
	HaveSameColor(N, N) bool
	GetGraph() AdjacencyGraph[N, E]
}

func ToIface_GraphColoring[N comparable, E any, T_iface_ interface {
	GraphColoring[N, E]
	comparable
}](v T_iface_) GraphColoring[N, E] {
	var z T_iface_
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.jscomp.graph.GraphColoring
type GraphColoringBase[N comparable, E any] struct {
	colorToNodeMap []N                  `gson:"colorToNodeMap"`
	graph          AdjacencyGraph[N, E] `gson:"graph"`
	idHash_        atomic.Int32
}

func (this_ *GraphColoringBase[N, E]) BindSelf_(s any) {
	_ = s
}
func (this_ *GraphColoringBase[N, E]) IsJavaNil_() bool { return this_ == nil }
func (*GraphColoringBase[N, E]) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.graph.GraphColoring", "GraphColoring"
}
func (this_ *GraphColoringBase[N, E]) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*GraphColoringBase[N, E]) IsGenericInst_graph_GraphColoring_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.graph.GraphColoring", "com.google.javascript.jscomp.graph.GraphColoring", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			IsGenericInst_graph_GraphColoring_()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}
func (b *GraphColoringBase[N, E]) GetGraphColoringBase() *GraphColoringBase[N, E] { return b }
func (*GraphColoringBase[N, E]) Implements_graph_GraphColoring()                  {}

func (this_ *GraphColoringBase[N, E]) initFields_GraphColoring() {
	_ = this_
}

func InitGraphColoring[N comparable, E any](this_ *GraphColoringBase[N, E], graph AdjacencyGraph[N, E]) {
	_ = graph
	this_.initFields_GraphColoring()
	this_.graph = graph
}

// Java method: GraphColoring.color()
func (this_ *GraphColoringBase[N, E]) Color() int32 {
	panic("abstract or native method invoked")
}

// Java method: GraphColoring.getPartitionSuperNode(N node)
func (this_ *GraphColoringBase[N, E]) GetPartitionSuperNode(node N) N {
	_ = node
	j2gort.CheckNotNullSlice(this_.colorToNodeMap)
	var color *GraphColoring_Color = j2gort.CastPtr[GraphColoring_Color](j2gort.AsIface[GraphNode[N, E]](this_.graph.GetNode(node)).GetAnnotation())
	_ = color
	var headNode N = this_.colorToNodeMap[color.value]
	_ = headNode
	if headNode == func() N {
		if _, ok_ := any((*N)(nil)).(*string); ok_ {
			return any(j2gort.NullString).(N)
		}
		var z_ N
		return z_
	}() {
		{
			a_1 := this_.colorToNodeMap
			i_2 := color.value
			a_1[i_2] = node
		}
		return node
	} else {
		return headNode
	}
}

// Java method: GraphColoring.haveSameColor(N first, N second)
func (this_ *GraphColoringBase[N, E]) HaveSameColor(first N, second N) bool {
	_ = first
	_ = second
	j2gort.CheckNotNullSlice(this_.colorToNodeMap)
	var colorFirst *GraphColoring_Color = j2gort.CastPtr[GraphColoring_Color](j2gort.AsIface[GraphNode[N, E]](this_.graph.GetNode(first)).GetAnnotation())
	_ = colorFirst
	var colorSecond *GraphColoring_Color = j2gort.CastPtr[GraphColoring_Color](j2gort.AsIface[GraphNode[N, E]](this_.graph.GetNode(second)).GetAnnotation())
	_ = colorSecond
	j2gort.CheckNotNullPtr(colorFirst)
	j2gort.CheckNotNullPtr(colorSecond)
	return colorFirst.Equals(j2gort.PtrToAny(colorSecond))
}

// Java method: GraphColoring.getGraph()
func (this_ *GraphColoringBase[N, E]) GetGraph() AdjacencyGraph[N, E] {
	return this_.graph
}

// Java class: com.google.javascript.jscomp.graph.GraphColoring.Color
type GraphColoring_Color struct {
	value int32 `gson:"value"`
}

func (this_ *GraphColoring_Color) BindSelf_(s any) {
	_ = s
}
func (this_ *GraphColoring_Color) IsJavaNil_() bool { return this_ == nil }
func (*GraphColoring_Color) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.graph.GraphColoring$Color", "Color"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.graph.GraphColoring.Color", "com.google.javascript.jscomp.graph.GraphColoring$Color", func(v_ any) bool {
		return j2gort.InstanceOf[*GraphColoring_Color](v_)
	})
}
func (*GraphColoring_Color) Implements_graph_Annotation() {}

func (this_ *GraphColoring_Color) initFields_GraphColoring_Color() {
	_ = this_

}

// Java constructor: Color(int value)
func NewGraphColoring_Color(value int32) *GraphColoring_Color {
	this_ := &GraphColoring_Color{}
	InitGraphColoring_Color(this_, value)
	return this_
}

func InitGraphColoring_Color(this_ *GraphColoring_Color, value int32) {
	_ = value
	this_.initFields_GraphColoring_Color()
	this_.value = value
}

// Java method: Color.equals(java.lang.Object other)
func (this_ *GraphColoring_Color) Equals(arg_3 any) bool {
	other := arg_3
	_ = other
	var color *GraphColoring_Color
	_ = color
	var ok_4 bool
	color, ok_4 = j2gort.TryCast[*GraphColoring_Color](other)
	if !ok_4 {
		return false
	} else {
		return (this_.value == color.value)
	}
}

// Java method: Color.hashCode()
func (this_ *GraphColoring_Color) HashCode() int32 {
	return this_.value
}

// Java class: com.google.javascript.jscomp.graph.GraphColoring.GreedyGraphColoring
type GraphColoring_GreedyGraphColoring[N comparable, E any] struct {
	GraphColoringBase[N, E]
	tieBreaker j2gort.Ext_java_util_Comparator `gson:"tieBreaker"`
}

func (this_ *GraphColoring_GreedyGraphColoring[N, E]) BindSelf_(s any) {
	_ = s
}
func (this_ *GraphColoring_GreedyGraphColoring[N, E]) IsJavaNil_() bool { return this_ == nil }
func (*GraphColoring_GreedyGraphColoring[N, E]) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.graph.GraphColoring$GreedyGraphColoring", "GreedyGraphColoring"
}
func (this_ *GraphColoring_GreedyGraphColoring[N, E]) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.GraphColoringBase.IdentityHash_()
}
func (*GraphColoring_GreedyGraphColoring[N, E]) HasIdentityEquals_() {}
func (*GraphColoring_GreedyGraphColoring[N, E]) IsGenericInst_graph_GraphColoring_GreedyGraphColoring_() {
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.graph.GraphColoring.GreedyGraphColoring", "com.google.javascript.jscomp.graph.GraphColoring$GreedyGraphColoring", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			IsGenericInst_graph_GraphColoring_GreedyGraphColoring_()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}
func (*GraphColoring_GreedyGraphColoring[N, E]) Implements_graph_GraphColoring() {}

func (this_ *GraphColoring_GreedyGraphColoring[N, E]) initFields_GraphColoring_GreedyGraphColoring() {
	_ = this_
}

// Java constructor: GreedyGraphColoring(com.google.javascript.jscomp.graph.AdjacencyGraph<N,E> graph)
func NewGraphColoring_GreedyGraphColoring_1[N comparable, E any](graph AdjacencyGraph[N, E]) *GraphColoring_GreedyGraphColoring[N, E] {
	this_ := &GraphColoring_GreedyGraphColoring[N, E]{}
	InitGraphColoring_GreedyGraphColoring_1[N, E](this_, graph)
	return this_
}

func InitGraphColoring_GreedyGraphColoring_1[N comparable, E any](this_ *GraphColoring_GreedyGraphColoring[N, E], graph AdjacencyGraph[N, E]) {
	_ = graph
	InitGraphColoring_GreedyGraphColoring_2[N, E](this_, graph, nil)
}

// Java constructor: GreedyGraphColoring(com.google.javascript.jscomp.graph.AdjacencyGraph<N,E> graph, java.util.@org.jspecify.annotations.Nullable Comparator<N> tieBreaker)
func NewGraphColoring_GreedyGraphColoring_2[N comparable, E any](graph AdjacencyGraph[N, E], tieBreaker j2gort.Ext_java_util_Comparator) *GraphColoring_GreedyGraphColoring[N, E] {
	this_ := &GraphColoring_GreedyGraphColoring[N, E]{}
	InitGraphColoring_GreedyGraphColoring_2[N, E](this_, graph, tieBreaker)
	return this_
}

func InitGraphColoring_GreedyGraphColoring_2[N comparable, E any](this_ *GraphColoring_GreedyGraphColoring[N, E], graph AdjacencyGraph[N, E], tieBreaker j2gort.Ext_java_util_Comparator) {
	_ = graph
	_ = tieBreaker
	InitGraphColoring[N, E](&this_.GraphColoringBase, graph)
	this_.initFields_GraphColoring_GreedyGraphColoring()
	this_.tieBreaker = tieBreaker
}

// Java method: GreedyGraphColoring.color()
func (this_ *GraphColoring_GreedyGraphColoring[N, E]) Color() int32 {
	var worklist j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList_1_Ext_java_util_Collection(this_.GetGraphColoringBase().graph.GetNodes())
	_ = worklist
	j2gort.Ext_java_util_Collections_Sort_2(worklist, j2gort.NewLambda_Ext_java_util_Comparator(func(lp_5 any, lp_6 any) int32 {
		leftNode := j2gort.AsIface[GraphNode[N, E]](lp_5)
		_ = leftNode
		rightNode := j2gort.AsIface[GraphNode[N, E]](lp_6)
		_ = rightNode
		var leftProperty N = leftNode.GetValue()
		_ = leftProperty
		var rightProperty N = rightNode.GetValue()
		_ = rightProperty
		var result int32 = (this_.GetGraphColoringBase().graph.GetWeight(rightProperty) - this_.GetGraphColoringBase().graph.GetWeight(leftProperty))
		_ = result
		if (result == int32(0)) && (this_.tieBreaker != nil) {
			return this_.tieBreaker.Compare(j2gort.ToAny(leftProperty), j2gort.ToAny(rightProperty))
		} else {
			return result
		}
	}))
	var count int32 = int32(0)
	_ = count
	for doFirst_7 := true; doFirst_7 || (!(worklist.IsEmpty())); doFirst_7 = false {
		var color *GraphColoring_Color = NewGraphColoring_Color(count)
		_ = color
		var subgraph SubGraph[N, E] = this_.GetGraphColoringBase().graph.NewSubGraph()
		_ = subgraph
		{
			var i j2gort.Ext_java_util_Iterator = j2gort.AsIface[j2gort.Ext_java_util_Iterator](worklist.Iterator())
			_ = i
			for i.HasNext() {
				var node GraphNode[N, E] = j2gort.AsIface[GraphNode[N, E]](i.Next())
				_ = node
				if subgraph.IsIndependentOf(node.GetValue()) {
					subgraph.AddNode(node.GetValue())
					node.SetAnnotation(color)
					i.Remove()
				}
			}
		}
		count = count + 1
	}
	var map_ []N
	{
		arr_8 := make([]N, count)
		if _, ok_ := any((*N)(nil)).(*string); ok_ {
			ns_ := any(j2gort.NullString).(N)
			for i_ := range arr_8 {
				arr_8[i_] = ns_
			}
		}
		map_ = arr_8
	}
	_ = map_
	this_.GetGraphColoringBase().colorToNodeMap = map_
	return count
}
