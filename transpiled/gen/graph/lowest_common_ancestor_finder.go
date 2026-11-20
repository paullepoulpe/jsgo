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
// Port of com/google/javascript/jscomp/graph/LowestCommonAncestorFinder.java (generated mechanically by j2go).

package graph

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_lowest_common_ancestor_finder_1 j2gort.Ext_java_util_function_IntFunction

func init() {
	lambdaSingleton_lowest_common_ancestor_finder_1 = j2gort.NewLambda_Ext_java_util_function_IntFunction(func(mrp0_ int32) any { return NewLowestCommonAncestorFinder_Color(mrp0_) })
}

// Java class: com.google.javascript.jscomp.graph.LowestCommonAncestorFinder
type LowestCommonAncestorFinder[N comparable, E any] struct {
	graph          DiGraph[N, E]                           `gson:"graph"`
	searchColoring *j2gort.Ext_java_util_LinkedHashMap     `gson:"searchColoring"`
	searchQueue    j2gort.Deque[DiGraph_DiGraphNode[N, E]] `gson:"searchQueue"`
	idHash_        atomic.Int32
}

func (this_ *LowestCommonAncestorFinder[N, E]) BindSelf_(s any) {
	_ = s
}
func (this_ *LowestCommonAncestorFinder[N, E]) IsJavaNil_() bool { return this_ == nil }
func (*LowestCommonAncestorFinder[N, E]) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.graph.LowestCommonAncestorFinder", "LowestCommonAncestorFinder"
}
func (this_ *LowestCommonAncestorFinder[N, E]) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*LowestCommonAncestorFinder[N, E]) HasIdentityEquals_()                              {}
func (*LowestCommonAncestorFinder[N, E]) IsGenericInst_graph_LowestCommonAncestorFinder_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.graph.LowestCommonAncestorFinder", "com.google.javascript.jscomp.graph.LowestCommonAncestorFinder", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			IsGenericInst_graph_LowestCommonAncestorFinder_()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}

func (this_ *LowestCommonAncestorFinder[N, E]) initFields_LowestCommonAncestorFinder() {
	_ = this_
	this_.searchColoring = j2gort.NewExt_java_util_LinkedHashMap()
	this_.searchQueue = j2gort.Deque[DiGraph_DiGraphNode[N, E]]{}
}

// Java constructor: LowestCommonAncestorFinder(com.google.javascript.jscomp.graph.DiGraph<N,E> graph)
func NewLowestCommonAncestorFinder[N comparable, E any](graph DiGraph[N, E]) *LowestCommonAncestorFinder[N, E] {
	this_ := &LowestCommonAncestorFinder[N, E]{}
	InitLowestCommonAncestorFinder[N, E](this_, graph)
	return this_
}

func InitLowestCommonAncestorFinder[N comparable, E any](this_ *LowestCommonAncestorFinder[N, E], graph DiGraph[N, E]) {
	_ = graph
	this_.initFields_LowestCommonAncestorFinder()
	this_.graph = graph
}

// Java method: LowestCommonAncestorFinder.findAll(java.util.Set<N> roots)
func (this_ *LowestCommonAncestorFinder[N, E]) FindAll(roots j2gort.Ext_java_util_Set) j2gort.Ext_com_google_common_collect_ImmutableSet {
	_ = roots
	if !(roots.Size() <= (int32(32) - int32(1))) {
		j2gort.CheckArgumentStr(false, "Too many roots.")
	}
	if !(this_.searchColoring.IsEmpty()) {
		j2gort.CheckState(false)
	}
	var allColor *LowestCommonAncestorFinder_Color = LowestCommonAncestorFinder_Color_Create(((javalang.Shl32(int32(1), int(roots.Size()))) - int32(1)))
	_ = allColor
	var bitForRoot int32 = int32(1)
	_ = bitForRoot
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(roots))
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			var root N
		iifeSwitch_1:
			switch {
			default:

				tv_3 := rawElem_2
				if tv_3 == nil {
					{

						if _, ok_ := any((*N)(nil)).(*string); ok_ {
							root = any(j2gort.NullString).(N)
							break iifeSwitch_1
						}
						var z_ N
						root = z_
						break iifeSwitch_1
					}

				}
				root = j2gort.Cast[N](tv_3)
				break iifeSwitch_1
			}

			_ = root
			var rootNode DiGraph_DiGraphNode[N, E] = j2gort.AsIface[DiGraph_DiGraphNode[N, E]](this_.graph.GetNode(root))
			_ = rootNode
			j2gort.CheckNotNullIface(rootNode)
			var color *LowestCommonAncestorFinder_Color = LowestCommonAncestorFinder_Color_Create(bitForRoot)
			_ = color
			this_.searchColoring.Merge(rootNode, j2gort.PtrToAny(color), j2gort.NewLambda_Ext_java_util_function_BiFunction(func(mrp0_ any, mrp1_ any) any {
				return j2gort.PtrToAny(j2gort.CastPtr[LowestCommonAncestorFinder_Color](mrp0_).Mix(j2gort.CastPtr[LowestCommonAncestorFinder_Color](mrp1_)))
			}))
			this_.priv_LowestCommonAncestorFinder_paintAncestors_6424(rootNode, color)
			bitForRoot = javalang.Shl32(bitForRoot, int(int32(1)))
		}
	}
	this_.searchColoring.ForEach_1(j2gort.NewLambda_Ext_java_util_function_BiConsumer(func(lp_4 any, lp_5 any) {
		node := j2gort.AsIface[DiGraph_DiGraphNode[N, E]](lp_4)
		_ = node
		color_2 := j2gort.CastPtr[LowestCommonAncestorFinder_Color](lp_5)
		_ = color_2
		if color_2.Equals(j2gort.PtrToAny(allColor)) {
			this_.priv_LowestCommonAncestorFinder_paintAncestors_6424(node, StaticGet_LowestCommonAncestorFinder_Color_NOT_LOWEST())
		}
	}))
	var results j2gort.Ext_com_google_common_collect_ImmutableSet_Builder = j2gort.Ext_com_google_common_collect_ImmutableSet_S_Builder()
	_ = results
	this_.searchColoring.ForEach_1(j2gort.NewLambda_Ext_java_util_function_BiConsumer(func(lp_6 any, lp_7 any) {
		node := j2gort.AsIface[DiGraph_DiGraphNode[N, E]](lp_6)
		_ = node
		color_2 := j2gort.CastPtr[LowestCommonAncestorFinder_Color](lp_7)
		_ = color_2
		if color_2.Equals(j2gort.PtrToAny(allColor)) {
			results.Add_1_any(j2gort.ToAny(node.GetValue()))
		}
	}))
	this_.searchColoring.Clear()
	this_.searchQueue.Clear()
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](results.Build())
}

// Java method: LowestCommonAncestorFinder.paintAncestors(com.google.javascript.jscomp.graph.DiGraph.DiGraphNode<N,E> root, com.google.javascript.jscomp.graph.LowestCommonAncestorFinder.Color color)
func (this_ *LowestCommonAncestorFinder[N, E]) priv_LowestCommonAncestorFinder_paintAncestors_6424(root DiGraph_DiGraphNode[N, E], color *LowestCommonAncestorFinder_Color) {
	_ = root
	_ = color
	if !(this_.searchQueue.IsEmpty()) {
		j2gort.CheckState(false)
	}
	this_.searchQueue.AddLast(root)
	for !(this_.searchQueue.IsEmpty()) {
		var curr DiGraph_DiGraphNode[N, E] = this_.searchQueue.RemoveFirst()
		_ = curr
		var parentEdges j2gort.Ext_java_util_List = curr.GetInEdges()
		_ = parentEdges
		{
			cur_8 := j2gort.NewIterCursor(j2gort.CheckNotNull(parentEdges))
			for cur_8.HasNext() {
				rawElem_9 := cur_8.Next()
				parentEdge := j2gort.AsIface[DiGraph_DiGraphEdge[N, E]](rawElem_9)
				_ = parentEdge
				var parent DiGraph_DiGraphNode[N, E] = parentEdge.GetSource()
				_ = parent
				if javalang.ObjectsEquals(j2gort.CheckNotNull(parent), root) {
					continue
				}
				var oldColor *LowestCommonAncestorFinder_Color = j2gort.CastPtr[LowestCommonAncestorFinder_Color](this_.searchColoring.GetOrDefault(parent, j2gort.PtrToAny(StaticGet_LowestCommonAncestorFinder_Color_BLANK())))
				_ = oldColor
				if !(oldColor.Contains(color)) {
					this_.searchQueue.AddLast(parent)
					this_.searchColoring.Put_2(parent, j2gort.PtrToAny(oldColor.Mix(color)))
				}
			}
		}
	}
}

// Java interface: com.google.javascript.jscomp.graph.LowestCommonAncestorFinder.Factory
type LowestCommonAncestorFinder_Factory[N comparable, E any] interface {
	Implements_graph_LowestCommonAncestorFinder_Factory()
	IsJavaNil_() bool
	Create(DiGraph[N, E]) *LowestCommonAncestorFinder[N, E]
}

func ToIface_LowestCommonAncestorFinder_Factory[N comparable, E any, T_iface_ interface {
	LowestCommonAncestorFinder_Factory[N, E]
	comparable
}](v T_iface_) LowestCommonAncestorFinder_Factory[N, E] {
	var z T_iface_
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.graph.LowestCommonAncestorFinder.Factory", "com.google.javascript.jscomp.graph.LowestCommonAncestorFinder$Factory", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			Implements_graph_LowestCommonAncestorFinder_Factory()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}

type lambda_LowestCommonAncestorFinder_Factory[N comparable, E any] struct {
	fn func(DiGraph[N, E]) *LowestCommonAncestorFinder[N, E]
}

func (l *lambda_LowestCommonAncestorFinder_Factory[N, E]) IsJavaNil_() bool { return l == nil }
func NewLambda_LowestCommonAncestorFinder_Factory[N comparable, E any](fn func(DiGraph[N, E]) *LowestCommonAncestorFinder[N, E]) LowestCommonAncestorFinder_Factory[N, E] {
	if fn == nil {
		return nil
	}
	return &lambda_LowestCommonAncestorFinder_Factory[N, E]{fn: fn}
}

func (*lambda_LowestCommonAncestorFinder_Factory[N, E]) Implements_graph_LowestCommonAncestorFinder_Factory() {
}
func (l *lambda_LowestCommonAncestorFinder_Factory[N, E]) Create(p0_ DiGraph[N, E]) *LowestCommonAncestorFinder[N, E] {
	return l.fn(p0_)
}

// Java class: com.google.javascript.jscomp.graph.LowestCommonAncestorFinder.Color
type LowestCommonAncestorFinder_Color struct {
	bitset int32 `gson:"bitset"`
}

func (this_ *LowestCommonAncestorFinder_Color) BindSelf_(s any) {
	_ = s
}
func (this_ *LowestCommonAncestorFinder_Color) IsJavaNil_() bool { return this_ == nil }
func (*LowestCommonAncestorFinder_Color) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.graph.LowestCommonAncestorFinder$Color", "Color"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.graph.LowestCommonAncestorFinder.Color", "com.google.javascript.jscomp.graph.LowestCommonAncestorFinder$Color", func(v_ any) bool {
		return j2gort.InstanceOf[*LowestCommonAncestorFinder_Color](v_)
	})
}

var static_LowestCommonAncestorFinder_Color_COMMON_COLOR_CACHE []*LowestCommonAncestorFinder_Color

func StaticGet_LowestCommonAncestorFinder_Color_COMMON_COLOR_CACHE() []*LowestCommonAncestorFinder_Color {
	clinit_LowestCommonAncestorFinder_Color()
	return static_LowestCommonAncestorFinder_Color_COMMON_COLOR_CACHE
}
func StaticSet_LowestCommonAncestorFinder_Color_COMMON_COLOR_CACHE(v []*LowestCommonAncestorFinder_Color) {
	clinit_LowestCommonAncestorFinder_Color()
	static_LowestCommonAncestorFinder_Color_COMMON_COLOR_CACHE = v
}

var static_LowestCommonAncestorFinder_Color_BLANK *LowestCommonAncestorFinder_Color

func StaticGet_LowestCommonAncestorFinder_Color_BLANK() *LowestCommonAncestorFinder_Color {
	clinit_LowestCommonAncestorFinder_Color()
	return static_LowestCommonAncestorFinder_Color_BLANK
}
func StaticSet_LowestCommonAncestorFinder_Color_BLANK(v *LowestCommonAncestorFinder_Color) {
	clinit_LowestCommonAncestorFinder_Color()
	static_LowestCommonAncestorFinder_Color_BLANK = v
}

var static_LowestCommonAncestorFinder_Color_NOT_LOWEST *LowestCommonAncestorFinder_Color

func StaticGet_LowestCommonAncestorFinder_Color_NOT_LOWEST() *LowestCommonAncestorFinder_Color {
	clinit_LowestCommonAncestorFinder_Color()
	return static_LowestCommonAncestorFinder_Color_NOT_LOWEST
}
func StaticSet_LowestCommonAncestorFinder_Color_NOT_LOWEST(v *LowestCommonAncestorFinder_Color) {
	clinit_LowestCommonAncestorFinder_Color()
	static_LowestCommonAncestorFinder_Color_NOT_LOWEST = v
}

var clinitDone_LowestCommonAncestorFinder_Color uint32
var clinitState_LowestCommonAncestorFinder_Color j2gort.ClinitState

func clinit_LowestCommonAncestorFinder_Color() {
	if atomic.LoadUint32(&clinitDone_LowestCommonAncestorFinder_Color) != 0 {
		return
	}
	clinitSlow_LowestCommonAncestorFinder_Color()
}

func clinitSlow_LowestCommonAncestorFinder_Color() {
	if !j2gort.EnterClinit(&clinitState_LowestCommonAncestorFinder_Color, &clinitDone_LowestCommonAncestorFinder_Color) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_LowestCommonAncestorFinder_Color, &clinitDone_LowestCommonAncestorFinder_Color)
	static_LowestCommonAncestorFinder_Color_COMMON_COLOR_CACHE = make([]*LowestCommonAncestorFinder_Color, javalang.Shl32(int32(2), int(int32(5))))
	func() {
		arrMut_10 := StaticGet_LowestCommonAncestorFinder_Color_COMMON_COLOR_CACHE()
		var arrBox_11 []any
		if arrMut_10 != nil {
			arrBox_11 = make([]any, len(arrMut_10))
			for i_, el_12 := range arrMut_10 {
				arrBox_11[i_] = j2gort.PtrToAny(el_12)
			}
		}
		j2gort.Ext_java_util_Arrays_SetAll_2_slice_any_Ext_java_util_function_IntFunction(arrBox_11, lambdaSingleton_lowest_common_ancestor_finder_1)
		for i_, el_12 := range arrBox_11 {
			arrMut_10[i_] = j2gort.CastPtr[LowestCommonAncestorFinder_Color](el_12)
		}
	}()
	static_LowestCommonAncestorFinder_Color_BLANK = j2gort.CheckNotNullPtr(StaticGet_LowestCommonAncestorFinder_Color_COMMON_COLOR_CACHE()[int32(0)])
	static_LowestCommonAncestorFinder_Color_NOT_LOWEST = NewLowestCommonAncestorFinder_Color(int32(-1))
}

func (this_ *LowestCommonAncestorFinder_Color) initFields_LowestCommonAncestorFinder_Color() {
	_ = this_
}

// Java constructor: Color(int bitset)
func NewLowestCommonAncestorFinder_Color(bitset int32) *LowestCommonAncestorFinder_Color {
	clinit_LowestCommonAncestorFinder_Color()
	this_ := &LowestCommonAncestorFinder_Color{}
	InitLowestCommonAncestorFinder_Color(this_, bitset)
	return this_
}

func InitLowestCommonAncestorFinder_Color(this_ *LowestCommonAncestorFinder_Color, bitset int32) {
	_ = bitset
	this_.initFields_LowestCommonAncestorFinder_Color()
	this_.bitset = bitset
}

// Java static method: Color.create(int bitset)
func LowestCommonAncestorFinder_Color_Create(bitset int32) *LowestCommonAncestorFinder_Color {
	clinit_LowestCommonAncestorFinder_Color()
	_ = bitset
	if bitset < int32(0) {
		if !(bitset == int32(-1)) {
			j2gort.CheckArgument(false)
		}
		return StaticGet_LowestCommonAncestorFinder_Color_NOT_LOWEST()
	} else if bitset < int32(len(StaticGet_LowestCommonAncestorFinder_Color_COMMON_COLOR_CACHE())) {
		return StaticGet_LowestCommonAncestorFinder_Color_COMMON_COLOR_CACHE()[bitset]
	}
	return NewLowestCommonAncestorFinder_Color(bitset)
}

// Java method: Color.mix(com.google.javascript.jscomp.graph.LowestCommonAncestorFinder.Color other)
func (this_ *LowestCommonAncestorFinder_Color) Mix(other *LowestCommonAncestorFinder_Color) *LowestCommonAncestorFinder_Color {
	_ = other
	if this_.bitset == other.bitset {
		return this_
	}
	return LowestCommonAncestorFinder_Color_Create((this_.bitset | other.bitset))
}

// Java method: Color.contains(com.google.javascript.jscomp.graph.LowestCommonAncestorFinder.Color other)
func (this_ *LowestCommonAncestorFinder_Color) Contains(other *LowestCommonAncestorFinder_Color) bool {
	_ = other
	return ((this_.bitset & other.bitset) == other.bitset)
}

// Java method: Color.equals(java.lang.Object other)
func (this_ *LowestCommonAncestorFinder_Color) Equals(arg_13 any) bool {
	other := arg_13
	_ = other
	return (this_.bitset == (j2gort.CastPtr[LowestCommonAncestorFinder_Color](other)).bitset)
}

// Java method: Color.hashCode()
func (this_ *LowestCommonAncestorFinder_Color) HashCode() int32 {
	return this_.bitset
}
