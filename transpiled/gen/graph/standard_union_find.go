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
// Port of com/google/javascript/jscomp/graph/StandardUnionFind.java (generated mechanically by j2go).

package graph

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.graph.StandardUnionFind
type StandardUnionFind[E comparable] struct {
	elmap                       j2gort.OrderedMap[E, *StandardUnionFind_Node[E]] `gson:"elmap"`
	idHash_                     atomic.Int32
	slab_StandardUnionFind_Node j2gort.Slab[StandardUnionFind_Node[E]]
}

func (this_ *StandardUnionFind[E]) BindSelf_(s any) {
	_ = s
}
func (this_ *StandardUnionFind[E]) IsJavaNil_() bool { return this_ == nil }
func (*StandardUnionFind[E]) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.graph.StandardUnionFind", "StandardUnionFind"
}
func (this_ *StandardUnionFind[E]) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*StandardUnionFind[E]) HasIdentityEquals_()                     {}
func (*StandardUnionFind[E]) IsGenericInst_graph_StandardUnionFind_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.graph.StandardUnionFind", "com.google.javascript.jscomp.graph.StandardUnionFind", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			IsGenericInst_graph_StandardUnionFind_()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}
func (*StandardUnionFind[E]) Implements_graph_UnionFind() {}

func (this_ *StandardUnionFind[E]) initFields_StandardUnionFind() {
	_ = this_
	this_.elmap = j2gort.NewOrderedMap[E, *StandardUnionFind_Node[E]](

	// Java constructor: StandardUnionFind()
	)
}

func NewStandardUnionFind[E comparable]() *StandardUnionFind[E] {
	this_ := &StandardUnionFind[E]{}
	InitStandardUnionFind[E](this_)
	return this_
}

func InitStandardUnionFind[E comparable](this_ *StandardUnionFind[E]) {
	this_.initFields_StandardUnionFind()
}

// Java constructor: StandardUnionFind(com.google.javascript.jscomp.graph.UnionFind<E> other)
func NewStandardUnionFind_1[E comparable](other UnionFind[E]) *StandardUnionFind[E] {
	this_ := &StandardUnionFind[E]{}
	InitStandardUnionFind_1[E](this_, other)
	return this_
}

func InitStandardUnionFind_1[E comparable](this_ *StandardUnionFind[E], other UnionFind[E]) {
	_ = other
	this_.initFields_StandardUnionFind()
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(other.Elements()))
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			var elem E
		iifeSwitch_1:
			switch {
			default:

				tv_3 := rawElem_2
				if tv_3 == nil {
					{

						if _, ok_ := any((*E)(nil)).(*string); ok_ {
							elem = any(j2gort.NullString).(E)
							break iifeSwitch_1
						}
						var z_ E
						elem = z_
						break iifeSwitch_1
					}

				}
				elem = j2gort.Cast[E](tv_3)
				break iifeSwitch_1
			}

			_ = elem
			this_.Union(other.Find(elem), elem)
		}
	}
}

// Java method: StandardUnionFind.add(@org.jspecify.annotations.Nullable E e)
func (this_ *StandardUnionFind[E]) Add(e E) {
	_ = e
	this_.Union(e, e)
}

// Java method: StandardUnionFind.addAll(java.lang.Iterable<E> es)
func (this_ *StandardUnionFind[E]) AddAll(es j2gort.Ext_java_lang_Iterable) {
	_ = es
	{
		cur_4 := j2gort.NewIterCursor(j2gort.CheckNotNull(es))
		for cur_4.HasNext() {
			rawElem_5 := cur_4.Next()
			var e E
		iifeSwitch_2:
			switch {
			default:

				tv_6 := rawElem_5
				if tv_6 == nil {
					{

						if _, ok_ := any((*E)(nil)).(*string); ok_ {
							e = any(j2gort.NullString).(E)
							break iifeSwitch_2
						}
						var z_ E
						e = z_
						break iifeSwitch_2
					}

				}
				e = j2gort.Cast[E](tv_6)
				break iifeSwitch_2
			}

			_ = e
			this_.Add(e)
		}
	}
}

// Java method: StandardUnionFind.union(@org.jspecify.annotations.Nullable E a, @org.jspecify.annotations.Nullable E b)
func (this_ *StandardUnionFind[E]) Union(a E, b E) E {
	_ = a
	_ = b
	var nodeA *StandardUnionFind_Node[E] = this_.priv_StandardUnionFind_findRootOrCreateNode_6426(a)
	_ = nodeA
	var nodeB *StandardUnionFind_Node[E] = this_.priv_StandardUnionFind_findRootOrCreateNode_6426(b)
	_ = nodeB
	if nodeA == nodeB {
		return nodeA.element
	}
	if nodeA.rank >= nodeB.rank {
		{
			r_7 := nodeB
			r_7.parent = nodeA
		}
		r_8 := nodeA
		r_8.size = (r_8.size + nodeB.size)
		if nodeA.rank == nodeB.rank {
			r_9 := nodeA
			r_9.rank = r_9.rank + 1
		}
		return nodeA.element
	}
	{
		r_10 := nodeA
		r_10.parent = nodeB
	}
	r_11 := nodeB
	r_11.size = (r_11.size + nodeA.size)
	var temp E = nodeB.element
	_ = temp
	{
		r_12 := nodeB
		r_12.element = nodeA.element
	}
	{
		r_13 := nodeA
		r_13.element = temp
	}
	return nodeB.element
}

// Java method: StandardUnionFind.find(@org.jspecify.annotations.Nullable E e)
func (this_ *StandardUnionFind[E]) Find(e E) E {
	_ = e
	var node *StandardUnionFind_Node[E] = this_.elmap.Get(e)
	_ = node
	if !(node != nil) {
		j2gort.CheckArgumentFmt(false, "Element does not exist: %s", j2gort.ToAny(e))
	}
	return this_.priv_StandardUnionFind_findRoot_6427(node).element
}

// Java method: StandardUnionFind.areEquivalent(@org.jspecify.annotations.Nullable E a, @org.jspecify.annotations.Nullable E b)
func (this_ *StandardUnionFind[E]) AreEquivalent(a E, b E) bool {
	_ = a
	_ = b
	var aRep E = this_.Find(a)
	_ = aRep
	var bRep E = this_.Find(b)
	_ = bRep
	return j2gort.RefEquals(aRep, bRep)
}

// Java method: StandardUnionFind.elements()
func (this_ *StandardUnionFind[E]) Elements() j2gort.Ext_java_util_Set {
	return j2gort.Ext_java_util_Collections_UnmodifiableSet(j2gort.AsIface[j2gort.Ext_java_util_Set](this_.elmap.LiveKeySet()))
}

// Java method: StandardUnionFind.allEquivalenceClasses()
func (this_ *StandardUnionFind[E]) AllEquivalenceClasses() j2gort.Ext_com_google_common_collect_ImmutableList {
	var groupsTmp j2gort.Ext_com_google_common_collect_SetMultimap = j2gort.AsIface[j2gort.Ext_com_google_common_collect_SetMultimap](j2gort.Ext_com_google_common_collect_MultimapBuilder_LinkedHashKeys_0().LinkedHashSetValues_0().Build_0())
	_ = groupsTmp
	{
		cur_14 := this_.elmap.IterValues()
		for cur_14.HasNext() {

			elem := cur_14.
				Next()
			_ = elem
			groupsTmp.Put(j2gort.PtrToAny(this_.priv_StandardUnionFind_findRoot_6427(elem)), j2gort.ToAny(elem.element))
		}
	}
	var result []j2gort.Ext_com_google_common_collect_ImmutableSet = []j2gort.Ext_com_google_common_collect_ImmutableSet{}
	_ = result
	{
		cur_16 := j2gort.NewMapValCursor(j2gort.Ext_com_google_common_collect_Multimaps_S_AsMap_1_Ext_com_google_common_collect_SetMultimap(groupsTmp))
		for cur_16.HasNext() {
			rawElem_17 := cur_16.Next()
			group := j2gort.AsIface[j2gort.Ext_java_util_Set](rawElem_17)
			_ = group
			result = append(result, j2gort.CheckNotNullElem(j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](group))))
		}
	}
	return j2gort.ImmListFromTypedSlice(result)
}

// Java method: StandardUnionFind.allRepresentatives()
func (this_ *StandardUnionFind[E]) AllRepresentatives() j2gort.Ext_com_google_common_collect_ImmutableSet {
	var builder j2gort.Ext_com_google_common_collect_ImmutableSet_Builder = j2gort.Ext_com_google_common_collect_ImmutableSet_S_Builder()
	_ = builder
	{
		cur_18 := this_.elmap.IterValues()
		for cur_18.HasNext() {

			n := cur_18.
				Next()
			_ = n
			if n == n.parent {
				builder.Add_1_any(j2gort.ToAny(n.element))
			}
		}
	}
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](builder.Build())
}

// Java method: StandardUnionFind.hasMultipleRepresentatives()
func (this_ *StandardUnionFind[E]) HasMultipleRepresentatives() bool {
	if this_.elmap.Size() <= int32(1) {
		return false
	}
	var count int32 = int32(0)
	_ = count
	{
		cur_20 := this_.elmap.IterValues()
		for cur_20.HasNext() {

			n := cur_20.
				Next()
			_ = n
			if (n == n.parent) && (func() int32 { old_22 := count; nw_23 := old_22 + 1; count = nw_23; return nw_23 }() > int32(1)) {
				return true
			}
		}
	}
	return false
}

// Java method: StandardUnionFind.getRepresentativeList()
func (this_ *StandardUnionFind[E]) GetRepresentativeList() j2gort.Ext_java_util_ArrayList {
	var reps []E = []E{}
	_ = reps
	{
		cur_24 := this_.elmap.IterValues()
		for cur_24.HasNext() {

			n := cur_24.
				Next()
			_ = n
			if n == n.parent {
				reps = append(reps, n.element)
			}
		}
	}
	return j2gort.WrapTypedArrayList(reps)
}

// Java method: StandardUnionFind.findRootOrCreateNode(E e)
func (this_ *StandardUnionFind[E]) priv_StandardUnionFind_findRootOrCreateNode_6426(e E) *StandardUnionFind_Node[E] {
	_ = e
	var node *StandardUnionFind_Node[E] = this_.elmap.Get(e)
	_ = node
	if node != nil {
		return this_.priv_StandardUnionFind_findRoot_6427(node)
	}
	node = alloc_StandardUnionFind_Node[E](&this_.slab_StandardUnionFind_Node, e)
	this_.elmap.Put(e, node)
	return node
}

// Java method: StandardUnionFind.findRoot(com.google.javascript.jscomp.graph.StandardUnionFind.Node<E> node)
func (this_ *StandardUnionFind[E]) priv_StandardUnionFind_findRoot_6427(node *StandardUnionFind_Node[E]) *StandardUnionFind_Node[E] {
	_ = node
	if node.parent != node {
		{
			r_26 := node
			r_26.parent = this_.priv_StandardUnionFind_findRoot_6427(node.parent)
		}
	}
	return node.parent
}

// Java method: StandardUnionFind.findAll(@org.jspecify.annotations.Nullable E value)
func (this_ *StandardUnionFind[E]) FindAll(value E) j2gort.Ext_java_util_Set {
	_ = value
	if !(this_.elmap.ContainsKey(value)) {
		j2gort.CheckArgumentFmt(false, "Element does not exist: %s", j2gort.ToAny(value))
	}
	var isSameRoot j2gort.Ext_com_google_common_base_Predicate = NewStandardUnionFind__Anon_119[E](this_, value)
	_ = isSameRoot
	return NewStandardUnionFind__Anon_120[E](this_, isSameRoot, value)
}

// Java class: com.google.javascript.jscomp.graph.StandardUnionFind.Node
type StandardUnionFind_Node[E comparable] struct {
	parent  *StandardUnionFind_Node[E] `gson:"parent"`
	element E                          `gson:"element"`
	rank    int32                      `gson:"rank"`
	size    int32                      `gson:"size"`
	idHash_ atomic.Int32
}

func (this_ *StandardUnionFind_Node[E]) BindSelf_(s any) {
	_ = s
}
func (this_ *StandardUnionFind_Node[E]) IsJavaNil_() bool { return this_ == nil }
func (*StandardUnionFind_Node[E]) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.graph.StandardUnionFind$Node", "Node"
}
func (this_ *StandardUnionFind_Node[E]) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*StandardUnionFind_Node[E]) HasIdentityEquals_()                          {}
func (*StandardUnionFind_Node[E]) IsGenericInst_graph_StandardUnionFind_Node_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.graph.StandardUnionFind.Node", "com.google.javascript.jscomp.graph.StandardUnionFind$Node", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			IsGenericInst_graph_StandardUnionFind_Node_()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}

func (this_ *StandardUnionFind_Node[E]) initFields_StandardUnionFind_Node() {
	_ = this_

	this_.size = int32(1)
}

// Java constructor: Node(E element)
func NewStandardUnionFind_Node[E comparable](element E) *StandardUnionFind_Node[E] {
	this_ := &StandardUnionFind_Node[E]{}
	InitStandardUnionFind_Node[E](this_, element)
	return this_
}
func alloc_StandardUnionFind_Node[E comparable](slab_ *j2gort.Slab[StandardUnionFind_Node[E]], element E) *StandardUnionFind_Node[E] {
	this_ := slab_.Alloc()
	InitStandardUnionFind_Node[E](this_, element)
	return this_
}

func InitStandardUnionFind_Node[E comparable](this_ *StandardUnionFind_Node[E], element E) {
	_ = element
	this_.initFields_StandardUnionFind_Node()
	this_.parent = this_
	this_.element = element
}

// Java class:
type StandardUnionFind__Anon_119[E comparable] struct {
	outer_       *StandardUnionFind[E]
	cap_0_value  E
	nodeForValue *StandardUnionFind_Node[E] `gson:"nodeForValue"`
	idHash_      atomic.Int32
}

func (this_ *StandardUnionFind__Anon_119[E]) BindSelf_(s any) {
	_ = s
}
func (this_ *StandardUnionFind__Anon_119[E]) IsJavaNil_() bool { return this_ == nil }
func (*StandardUnionFind__Anon_119[E]) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.graph.StandardUnionFind$1", ""
}
func (this_ *StandardUnionFind__Anon_119[E]) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*StandardUnionFind__Anon_119[E]) HasIdentityEquals_()                                     {}
func (*StandardUnionFind__Anon_119[E]) IsGenericInst_graph_StandardUnionFind__Anon_119_()       {}
func (this_ *StandardUnionFind__Anon_119[E]) SetOuter_(o *StandardUnionFind[E])                 { this_.outer_ = o }
func (this_ *StandardUnionFind__Anon_119[E]) GetOuter_() *StandardUnionFind[E]                  { return this_.outer_ }
func (*StandardUnionFind__Anon_119[E]) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
func (*StandardUnionFind__Anon_119[E]) Implements_j2gort_Ext_java_util_function_Predicate()     {}

func (this_ *StandardUnionFind__Anon_119[E]) initFields_StandardUnionFind__Anon_119() {
	_ = this_
	this_.nodeForValue = this_.outer_.elmap.Get(this_.cap_0_value)
}

func NewStandardUnionFind__Anon_119[E comparable](outer_ *StandardUnionFind[E], cap_p_value E) *StandardUnionFind__Anon_119[E] {
	this_ := &StandardUnionFind__Anon_119[E]{}
	this_.outer_ = outer_
	this_.cap_0_value = cap_p_value
	this_.initFields_StandardUnionFind__Anon_119()
	return this_
}

// Java method: .apply(java.lang.@org.jspecify.annotations.Nullable Object b)
func (this_ *StandardUnionFind__Anon_119[E]) Apply(arg_27 any) bool {
	b := arg_27
	_ = b
	if javalang.ObjectsEquals(j2gort.ToAny(this_.cap_0_value), b) {
		return true
	}
	var nodeForB *StandardUnionFind_Node[E] = j2gort.CastPtr[StandardUnionFind_Node[E]](j2gort.OrderedMapGetAny(&this_.outer_.elmap, b))
	_ = nodeForB
	if nodeForB == nil {
		return false
	}
	this_.nodeForValue = this_.outer_.priv_StandardUnionFind_findRoot_6427(this_.nodeForValue)
	return (this_.outer_.priv_StandardUnionFind_findRoot_6427(nodeForB) == this_.nodeForValue)
}

func (this_ *StandardUnionFind__Anon_119[E]) Equals(p0_ any) bool {
	return j2gort.RefEquals(this_, p0_)
}

func (this_ *StandardUnionFind__Anon_119[E]) Test(p0_ any) bool {
	return this_.Apply(p0_)
}

// Java class:
type StandardUnionFind__Anon_120[E comparable] struct {
	j2gort.Ext_java_util_AbstractSetBase
	outer_           *StandardUnionFind[E]
	cap_0_isSameRoot j2gort.Ext_com_google_common_base_Predicate
	cap_1_value      E
}

func (this_ *StandardUnionFind__Anon_120[E]) BindSelf_(s any) {
	_ = s
	this_.Ext_java_util_AbstractSetBase.BindExtSelf_(s)
}
func (this_ *StandardUnionFind__Anon_120[E]) IsJavaNil_() bool { return this_ == nil }
func (*StandardUnionFind__Anon_120[E]) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.graph.StandardUnionFind$2", ""
}
func (*StandardUnionFind__Anon_120[E]) IsGenericInst_graph_StandardUnionFind__Anon_120_()   {}
func (this_ *StandardUnionFind__Anon_120[E]) SetOuter_(o *StandardUnionFind[E])             { this_.outer_ = o }
func (this_ *StandardUnionFind__Anon_120[E]) GetOuter_() *StandardUnionFind[E]              { return this_.outer_ }
func (*StandardUnionFind__Anon_120[E]) Implements_j2gort_Ext_java_util_AbstractSet()        {}
func (*StandardUnionFind__Anon_120[E]) Implements_j2gort_Ext_java_util_AbstractCollection() {}
func (*StandardUnionFind__Anon_120[E]) Implements_j2gort_Ext_java_util_Set()                {}
func (*StandardUnionFind__Anon_120[E]) Implements_j2gort_Ext_java_util_Collection()         {}
func (*StandardUnionFind__Anon_120[E]) Implements_j2gort_Ext_java_lang_Iterable()           {}

func (this_ *StandardUnionFind__Anon_120[E]) initFields_StandardUnionFind__Anon_120() {
	_ = this_
}

func NewStandardUnionFind__Anon_120[E comparable](outer_ *StandardUnionFind[E], cap_p_isSameRoot j2gort.Ext_com_google_common_base_Predicate, cap_p_value E) *StandardUnionFind__Anon_120[E] {
	this_ := &StandardUnionFind__Anon_120[E]{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	this_.cap_0_isSameRoot = cap_p_isSameRoot
	this_.cap_1_value = cap_p_value
	j2gort.InitExt_java_util_AbstractSet(&this_.Ext_java_util_AbstractSetBase)
	this_.initFields_StandardUnionFind__Anon_120()
	return this_
}

// Java method: .contains(java.lang.Object o)
func (this_ *StandardUnionFind__Anon_120[E]) Contains(o any) bool {
	_ = o
	return this_.cap_0_isSameRoot.Apply(o)
}

// Java method: .iterator()
func (this_ *StandardUnionFind__Anon_120[E]) Iterator() any {
	return j2gort.Ext_com_google_common_collect_Iterators_Filter_2_Ext_java_util_Iterator_Ext_com_google_common_base_Predicate(j2gort.AsIface[j2gort.Ext_java_util_Iterator](j2gort.AsIface[j2gort.Ext_java_util_Set](this_.outer_.elmap.LiveKeySet()).Iterator()), this_.cap_0_isSameRoot)
}

// Java method: .size()
func (this_ *StandardUnionFind__Anon_120[E]) Size() int32 {
	return this_.outer_.priv_StandardUnionFind_findRoot_6427(this_.outer_.elmap.Get(this_.cap_1_value)).size
}
