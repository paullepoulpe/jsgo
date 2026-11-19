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
// Port of com/google/javascript/jscomp/graph/CheckPathsBetweenNodes.java (generated mechanically by j2go).

package graph

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.graph.CheckPathsBetweenNodes
type CheckPathsBetweenNodes[N comparable, E any] struct {
	nodePredicate j2gort.Ext_com_google_common_base_Predicate `gson:"nodePredicate"`
	edgePredicate j2gort.Ext_com_google_common_base_Predicate `gson:"edgePredicate"`
	graph         DiGraph[N, E]                               `gson:"graph"`
	start         DiGraph_DiGraphNode[N, E]                   `gson:"start"`
	end           DiGraph_DiGraphNode[N, E]                   `gson:"end"`
	idHash_       atomic.Int32
	inclusive     bool `gson:"inclusive"`
}

func (this_ *CheckPathsBetweenNodes[N, E]) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckPathsBetweenNodes[N, E]) IsJavaNil_() bool { return this_ == nil }
func (*CheckPathsBetweenNodes[N, E]) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.graph.CheckPathsBetweenNodes", "CheckPathsBetweenNodes"
}
func (this_ *CheckPathsBetweenNodes[N, E]) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CheckPathsBetweenNodes[N, E]) HasIdentityEquals_()                          {}
func (*CheckPathsBetweenNodes[N, E]) IsGenericInst_graph_CheckPathsBetweenNodes_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.graph.CheckPathsBetweenNodes", "com.google.javascript.jscomp.graph.CheckPathsBetweenNodes", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			IsGenericInst_graph_CheckPathsBetweenNodes_()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}

var static_CheckPathsBetweenNodes_BACK_EDGE Annotation

func StaticGet_CheckPathsBetweenNodes_BACK_EDGE() Annotation {
	clinit_CheckPathsBetweenNodes()
	return static_CheckPathsBetweenNodes_BACK_EDGE
}
func StaticSet_CheckPathsBetweenNodes_BACK_EDGE(v Annotation) {
	clinit_CheckPathsBetweenNodes()
	static_CheckPathsBetweenNodes_BACK_EDGE = v
}

var static_CheckPathsBetweenNodes_VISITED_EDGE Annotation

func StaticGet_CheckPathsBetweenNodes_VISITED_EDGE() Annotation {
	clinit_CheckPathsBetweenNodes()
	return static_CheckPathsBetweenNodes_VISITED_EDGE
}
func StaticSet_CheckPathsBetweenNodes_VISITED_EDGE(v Annotation) {
	clinit_CheckPathsBetweenNodes()
	static_CheckPathsBetweenNodes_VISITED_EDGE = v
}

var static_CheckPathsBetweenNodes_WHITE Annotation

func StaticGet_CheckPathsBetweenNodes_WHITE() Annotation {
	clinit_CheckPathsBetweenNodes()
	return static_CheckPathsBetweenNodes_WHITE
}
func StaticSet_CheckPathsBetweenNodes_WHITE(v Annotation) {
	clinit_CheckPathsBetweenNodes()
	static_CheckPathsBetweenNodes_WHITE = v
}

var static_CheckPathsBetweenNodes_GRAY Annotation

func StaticGet_CheckPathsBetweenNodes_GRAY() Annotation {
	clinit_CheckPathsBetweenNodes()
	return static_CheckPathsBetweenNodes_GRAY
}
func StaticSet_CheckPathsBetweenNodes_GRAY(v Annotation) {
	clinit_CheckPathsBetweenNodes()
	static_CheckPathsBetweenNodes_GRAY = v
}

var static_CheckPathsBetweenNodes_BLACK Annotation

func StaticGet_CheckPathsBetweenNodes_BLACK() Annotation {
	clinit_CheckPathsBetweenNodes()
	return static_CheckPathsBetweenNodes_BLACK
}
func StaticSet_CheckPathsBetweenNodes_BLACK(v Annotation) {
	clinit_CheckPathsBetweenNodes()
	static_CheckPathsBetweenNodes_BLACK = v
}

var clinitDone_CheckPathsBetweenNodes uint32
var clinitState_CheckPathsBetweenNodes j2gort.ClinitState

func clinit_CheckPathsBetweenNodes() {
	if atomic.LoadUint32(&clinitDone_CheckPathsBetweenNodes) != 0 {
		return
	}
	clinitSlow_CheckPathsBetweenNodes()
}

func clinitSlow_CheckPathsBetweenNodes() {
	if !j2gort.EnterClinit(&clinitState_CheckPathsBetweenNodes, &clinitDone_CheckPathsBetweenNodes) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckPathsBetweenNodes, &clinitDone_CheckPathsBetweenNodes)
	static_CheckPathsBetweenNodes_BACK_EDGE = NewCheckPathsBetweenNodes__Anon_114()
	static_CheckPathsBetweenNodes_VISITED_EDGE = NewCheckPathsBetweenNodes__Anon_115()
	static_CheckPathsBetweenNodes_WHITE = nil
	static_CheckPathsBetweenNodes_GRAY = NewCheckPathsBetweenNodes__Anon_116()
	static_CheckPathsBetweenNodes_BLACK = NewCheckPathsBetweenNodes__Anon_117()
}

func (this_ *CheckPathsBetweenNodes[N, E]) initFields_CheckPathsBetweenNodes() {
	_ = this_
}

// Java constructor: CheckPathsBetweenNodes(com.google.javascript.jscomp.graph.DiGraph<N,E> graph, com.google.javascript.jscomp.graph.DiGraph.DiGraphNode<N,E> a, com.google.javascript.jscomp.graph.DiGraph.DiGraphNode<N,E> b, com.google.common.base.Predicate<N> nodePredicate, com.google.common.base.Predicate<com.google.javascript.jscomp.graph.DiGraph.DiGraphEdge<N,E>> edgePredicate, boolean inclusive)
func NewCheckPathsBetweenNodes_6[N comparable, E any](graph DiGraph[N, E], a DiGraph_DiGraphNode[N, E], b DiGraph_DiGraphNode[N, E], nodePredicate j2gort.Ext_com_google_common_base_Predicate, edgePredicate j2gort.Ext_com_google_common_base_Predicate, inclusive bool) *CheckPathsBetweenNodes[N, E] {
	clinit_CheckPathsBetweenNodes()
	this_ := &CheckPathsBetweenNodes[N, E]{}
	InitCheckPathsBetweenNodes_6[N, E](this_, graph, a, b, nodePredicate, edgePredicate, inclusive)
	return this_
}

func InitCheckPathsBetweenNodes_6[N comparable, E any](this_ *CheckPathsBetweenNodes[N, E], graph DiGraph[N, E], a DiGraph_DiGraphNode[N, E], b DiGraph_DiGraphNode[N, E], nodePredicate j2gort.Ext_com_google_common_base_Predicate, edgePredicate j2gort.Ext_com_google_common_base_Predicate, inclusive bool) {
	_ = graph
	_ = a
	_ = b
	_ = nodePredicate
	_ = edgePredicate
	_ = inclusive
	this_.initFields_CheckPathsBetweenNodes()
	this_.graph = graph
	this_.start = a
	this_.end = b
	this_.nodePredicate = nodePredicate
	this_.edgePredicate = edgePredicate
	this_.inclusive = inclusive
}

// Java constructor: CheckPathsBetweenNodes(com.google.javascript.jscomp.graph.DiGraph<N,E> graph, com.google.javascript.jscomp.graph.DiGraph.DiGraphNode<N,E> a, com.google.javascript.jscomp.graph.DiGraph.DiGraphNode<N,E> b, com.google.common.base.Predicate<N> nodePredicate, com.google.common.base.Predicate<com.google.javascript.jscomp.graph.DiGraph.DiGraphEdge<N,E>> edgePredicate)
func NewCheckPathsBetweenNodes_5[N comparable, E any](graph DiGraph[N, E], a DiGraph_DiGraphNode[N, E], b DiGraph_DiGraphNode[N, E], nodePredicate j2gort.Ext_com_google_common_base_Predicate, edgePredicate j2gort.Ext_com_google_common_base_Predicate) *CheckPathsBetweenNodes[N, E] {
	clinit_CheckPathsBetweenNodes()
	this_ := &CheckPathsBetweenNodes[N, E]{}
	InitCheckPathsBetweenNodes_5[N, E](this_, graph, a, b, nodePredicate, edgePredicate)
	return this_
}

func InitCheckPathsBetweenNodes_5[N comparable, E any](this_ *CheckPathsBetweenNodes[N, E], graph DiGraph[N, E], a DiGraph_DiGraphNode[N, E], b DiGraph_DiGraphNode[N, E], nodePredicate j2gort.Ext_com_google_common_base_Predicate, edgePredicate j2gort.Ext_com_google_common_base_Predicate) {
	_ = graph
	_ = a
	_ = b
	_ = nodePredicate
	_ = edgePredicate
	InitCheckPathsBetweenNodes_6[N, E](this_, graph, a, b, nodePredicate, edgePredicate, true)
}

// Java method: CheckPathsBetweenNodes.allPathsSatisfyPredicate()
func (this_ *CheckPathsBetweenNodes[N, E]) AllPathsSatisfyPredicate() bool {
	this_.priv_CheckPathsBetweenNodes_setUp_6259()
	var result bool = this_.priv_CheckPathsBetweenNodes_checkAllPathsWithoutBackEdges_6263(this_.start, this_.end)
	_ = result
	this_.priv_CheckPathsBetweenNodes_tearDown_6260()
	return result
}

// Java method: CheckPathsBetweenNodes.somePathsSatisfyPredicate()
func (this_ *CheckPathsBetweenNodes[N, E]) SomePathsSatisfyPredicate() bool {
	this_.priv_CheckPathsBetweenNodes_setUp_6259()
	var result bool = this_.priv_CheckPathsBetweenNodes_checkSomePathsWithoutBackEdges_6264(this_.start, this_.end)
	_ = result
	this_.priv_CheckPathsBetweenNodes_tearDown_6260()
	return result
}

// Java method: CheckPathsBetweenNodes.setUp()
func (this_ *CheckPathsBetweenNodes[N, E]) priv_CheckPathsBetweenNodes_setUp_6259() {
	this_.graph.PushNodeAnnotations()
	this_.graph.PushEdgeAnnotations()
	this_.priv_CheckPathsBetweenNodes_discoverBackEdges_6261(this_.start)
}

// Java method: CheckPathsBetweenNodes.tearDown()
func (this_ *CheckPathsBetweenNodes[N, E]) priv_CheckPathsBetweenNodes_tearDown_6260() {
	this_.graph.PopNodeAnnotations()
	this_.graph.PopEdgeAnnotations()
}

// Java method: CheckPathsBetweenNodes.discoverBackEdges(com.google.javascript.jscomp.graph.DiGraph.DiGraphNode<N,E> u)
func (this_ *CheckPathsBetweenNodes[N, E]) priv_CheckPathsBetweenNodes_discoverBackEdges_6261(u DiGraph_DiGraphNode[N, E]) {
	_ = u
	u.SetAnnotation(StaticGet_CheckPathsBetweenNodes_GRAY())
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(u.GetOutEdges()))
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			e := j2gort.AsIface[DiGraph_DiGraphEdge[N, E]](rawElem_2)
			_ = e
			if this_.priv_CheckPathsBetweenNodes_ignoreEdge_6262(e) {
				continue
			}
			var v DiGraph_DiGraphNode[N, E] = e.GetDestination()
			_ = v
			if j2gort.IfaceEquals[Annotation](v.GetAnnotation(), StaticGet_CheckPathsBetweenNodes_WHITE()) {
				this_.priv_CheckPathsBetweenNodes_discoverBackEdges_6261(v)
			} else if j2gort.IfaceEquals[Annotation](v.GetAnnotation(), StaticGet_CheckPathsBetweenNodes_GRAY()) {
				e.SetAnnotation(StaticGet_CheckPathsBetweenNodes_BACK_EDGE())
			}
		}
	}
	u.SetAnnotation(StaticGet_CheckPathsBetweenNodes_BLACK())
}

// Java method: CheckPathsBetweenNodes.ignoreEdge(com.google.javascript.jscomp.graph.DiGraph.DiGraphEdge<N,E> e)
func (this_ *CheckPathsBetweenNodes[N, E]) priv_CheckPathsBetweenNodes_ignoreEdge_6262(e DiGraph_DiGraphEdge[N, E]) bool {
	_ = e
	return !(this_.edgePredicate.Apply(e))
}

// Java method: CheckPathsBetweenNodes.checkAllPathsWithoutBackEdges(com.google.javascript.jscomp.graph.DiGraph.DiGraphNode<N,E> a, com.google.javascript.jscomp.graph.DiGraph.DiGraphNode<N,E> b)
func (this_ *CheckPathsBetweenNodes[N, E]) priv_CheckPathsBetweenNodes_checkAllPathsWithoutBackEdges_6263(a DiGraph_DiGraphNode[N, E], b DiGraph_DiGraphNode[N, E]) bool {
	_ = a
	_ = b
	if this_.nodePredicate.Apply(j2gort.ToAny(a.GetValue())) && (this_.inclusive || (!(a == this_.start || ((a == nil || a.IsJavaNil_()) && (this_.start == nil || this_.start.IsJavaNil_()))) && !(a == this_.end || ((a == nil || a.IsJavaNil_()) && (this_.end == nil || this_.end.IsJavaNil_()))))) {
		return true
	}
	if a == b || ((a == nil || a.IsJavaNil_()) && (b == nil || b.IsJavaNil_())) {
		return false
	}
	{
		cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(a.GetOutEdges()))
		for cur_3.HasNext() {
			rawElem_4 := cur_3.Next()
			e := j2gort.AsIface[DiGraph_DiGraphEdge[N, E]](rawElem_4)
			_ = e
			if j2gort.IfaceEquals[Annotation](e.GetAnnotation(), StaticGet_CheckPathsBetweenNodes_VISITED_EDGE()) {
				continue
			}
			if this_.priv_CheckPathsBetweenNodes_ignoreEdge_6262(e) {
				continue
			}
			if j2gort.IfaceEquals[Annotation](e.GetAnnotation(), StaticGet_CheckPathsBetweenNodes_BACK_EDGE()) {
				continue
			}
			e.SetAnnotation(StaticGet_CheckPathsBetweenNodes_VISITED_EDGE())
			var next DiGraph_DiGraphNode[N, E] = e.GetDestination()
			_ = next
			if !(this_.priv_CheckPathsBetweenNodes_checkAllPathsWithoutBackEdges_6263(next, b)) {
				return false
			}
		}
	}
	return true
}

// Java method: CheckPathsBetweenNodes.checkSomePathsWithoutBackEdges(com.google.javascript.jscomp.graph.DiGraph.DiGraphNode<N,E> a, com.google.javascript.jscomp.graph.DiGraph.DiGraphNode<N,E> b)
func (this_ *CheckPathsBetweenNodes[N, E]) priv_CheckPathsBetweenNodes_checkSomePathsWithoutBackEdges_6264(a DiGraph_DiGraphNode[N, E], b DiGraph_DiGraphNode[N, E]) bool {
	_ = a
	_ = b
	if this_.nodePredicate.Apply(j2gort.ToAny(a.GetValue())) && (this_.inclusive || (!(a == this_.start || ((a == nil || a.IsJavaNil_()) && (this_.start == nil || this_.start.IsJavaNil_()))) && !(a == this_.end || ((a == nil || a.IsJavaNil_()) && (this_.end == nil || this_.end.IsJavaNil_()))))) {
		return true
	}
	if a == b || ((a == nil || a.IsJavaNil_()) && (b == nil || b.IsJavaNil_())) {
		return false
	}
	{
		cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(a.GetOutEdges()))
		for cur_5.HasNext() {
			rawElem_6 := cur_5.Next()
			e := j2gort.AsIface[DiGraph_DiGraphEdge[N, E]](rawElem_6)
			_ = e
			if j2gort.IfaceEquals[Annotation](e.GetAnnotation(), StaticGet_CheckPathsBetweenNodes_VISITED_EDGE()) {
				continue
			}
			if this_.priv_CheckPathsBetweenNodes_ignoreEdge_6262(e) {
				continue
			}
			if j2gort.IfaceEquals[Annotation](e.GetAnnotation(), StaticGet_CheckPathsBetweenNodes_BACK_EDGE()) {
				continue
			}
			e.SetAnnotation(StaticGet_CheckPathsBetweenNodes_VISITED_EDGE())
			var next DiGraph_DiGraphNode[N, E] = e.GetDestination()
			_ = next
			if this_.priv_CheckPathsBetweenNodes_checkSomePathsWithoutBackEdges_6264(next, b) {
				return true
			}
		}
	}
	return false
}

// Java class:
type CheckPathsBetweenNodes__Anon_114 struct {
	idHash_ atomic.Int32
}

func (this_ *CheckPathsBetweenNodes__Anon_114) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckPathsBetweenNodes__Anon_114) IsJavaNil_() bool { return this_ == nil }
func (*CheckPathsBetweenNodes__Anon_114) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.graph.CheckPathsBetweenNodes$1", ""
}
func (this_ *CheckPathsBetweenNodes__Anon_114) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CheckPathsBetweenNodes__Anon_114) HasIdentityEquals_()          {}
func (*CheckPathsBetweenNodes__Anon_114) Implements_graph_Annotation() {}

func (this_ *CheckPathsBetweenNodes__Anon_114) initFields_CheckPathsBetweenNodes__Anon_114() {
	_ = this_
}

func NewCheckPathsBetweenNodes__Anon_114() *CheckPathsBetweenNodes__Anon_114 {
	this_ := &CheckPathsBetweenNodes__Anon_114{}
	this_.initFields_CheckPathsBetweenNodes__Anon_114()
	return this_
}

// Java class:
type CheckPathsBetweenNodes__Anon_115 struct {
	idHash_ atomic.Int32
}

func (this_ *CheckPathsBetweenNodes__Anon_115) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckPathsBetweenNodes__Anon_115) IsJavaNil_() bool { return this_ == nil }
func (*CheckPathsBetweenNodes__Anon_115) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.graph.CheckPathsBetweenNodes$2", ""
}
func (this_ *CheckPathsBetweenNodes__Anon_115) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CheckPathsBetweenNodes__Anon_115) HasIdentityEquals_()          {}
func (*CheckPathsBetweenNodes__Anon_115) Implements_graph_Annotation() {}

func (this_ *CheckPathsBetweenNodes__Anon_115) initFields_CheckPathsBetweenNodes__Anon_115() {
	_ = this_
}

func NewCheckPathsBetweenNodes__Anon_115() *CheckPathsBetweenNodes__Anon_115 {
	this_ := &CheckPathsBetweenNodes__Anon_115{}
	this_.initFields_CheckPathsBetweenNodes__Anon_115()
	return this_
}

// Java class:
type CheckPathsBetweenNodes__Anon_116 struct {
	idHash_ atomic.Int32
}

func (this_ *CheckPathsBetweenNodes__Anon_116) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckPathsBetweenNodes__Anon_116) IsJavaNil_() bool { return this_ == nil }
func (*CheckPathsBetweenNodes__Anon_116) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.graph.CheckPathsBetweenNodes$3", ""
}
func (this_ *CheckPathsBetweenNodes__Anon_116) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CheckPathsBetweenNodes__Anon_116) HasIdentityEquals_()          {}
func (*CheckPathsBetweenNodes__Anon_116) Implements_graph_Annotation() {}

func (this_ *CheckPathsBetweenNodes__Anon_116) initFields_CheckPathsBetweenNodes__Anon_116() {
	_ = this_
}

func NewCheckPathsBetweenNodes__Anon_116() *CheckPathsBetweenNodes__Anon_116 {
	this_ := &CheckPathsBetweenNodes__Anon_116{}
	this_.initFields_CheckPathsBetweenNodes__Anon_116()
	return this_
}

// Java class:
type CheckPathsBetweenNodes__Anon_117 struct {
	idHash_ atomic.Int32
}

func (this_ *CheckPathsBetweenNodes__Anon_117) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckPathsBetweenNodes__Anon_117) IsJavaNil_() bool { return this_ == nil }
func (*CheckPathsBetweenNodes__Anon_117) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.graph.CheckPathsBetweenNodes$4", ""
}
func (this_ *CheckPathsBetweenNodes__Anon_117) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CheckPathsBetweenNodes__Anon_117) HasIdentityEquals_()          {}
func (*CheckPathsBetweenNodes__Anon_117) Implements_graph_Annotation() {}

func (this_ *CheckPathsBetweenNodes__Anon_117) initFields_CheckPathsBetweenNodes__Anon_117() {
	_ = this_
}

func NewCheckPathsBetweenNodes__Anon_117() *CheckPathsBetweenNodes__Anon_117 {
	this_ := &CheckPathsBetweenNodes__Anon_117{}
	this_.initFields_CheckPathsBetweenNodes__Anon_117()
	return this_
}
