/*
 * Copyright 2020 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/disambiguate/ColorGraphBuilder.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	pkg_graph "github.com/paullepoulpe/jsgo/transpiled/gen/graph"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_color_graph_builder_1 j2gort.Ext_com_google_common_base_Predicate

func init() {
	lambdaSingleton_color_graph_builder_1 = j2gort.NewLambda_Ext_com_google_common_base_Predicate(func(lp_5 any) bool {
		t := lp_5
		_ = t
		return true
	})
}

// Java class: com.google.javascript.jscomp.disambiguate.ColorGraphBuilder
type ColorGraphBuilder struct {
	registry                *pkg_colors.ColorRegistry                                              `gson:"registry"`
	nodeFactory             *ColorGraphNodeFactory                                                 `gson:"nodeFactory"`
	lcaFinder               *pkg_graph.LowestCommonAncestorFinder[*ColorGraphNode, any]            `gson:"lcaFinder"`
	topNode                 *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[*ColorGraphNode, any] `gson:"topNode"`
	colorHoldsInstanceGraph pkg_graph.LinkedDirectedGraph[*ColorGraphNode, any]                    `gson:"colorHoldsInstanceGraph"`
	idHash_                 atomic.Int32
}

func (this_ *ColorGraphBuilder) BindSelf_(s any) {
	_ = s
}
func (this_ *ColorGraphBuilder) IsJavaNil_() bool { return this_ == nil }
func (*ColorGraphBuilder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.disambiguate.ColorGraphBuilder", "ColorGraphBuilder"
}
func (this_ *ColorGraphBuilder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ColorGraphBuilder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.disambiguate.ColorGraphBuilder", "com.google.javascript.jscomp.disambiguate.ColorGraphBuilder", func(v_ any) bool {
		return j2gort.InstanceOf[*ColorGraphBuilder](v_)
	})
}

func (this_ *ColorGraphBuilder) initFields_ColorGraphBuilder() {
	_ = this_
	this_.colorHoldsInstanceGraph = pkg_graph.LinkedDirectedGraph_CreateWithoutAnnotations[*ColorGraphNode, any]()
}

// Java constructor: ColorGraphBuilder(com.google.javascript.jscomp.disambiguate.ColorGraphNodeFactory nodeFactory, com.google.javascript.jscomp.graph.LowestCommonAncestorFinder.Factory<com.google.javascript.jscomp.disambiguate.ColorGraphNode,java.lang.Object> lcaFinderFactory, com.google.javascript.jscomp.colors.ColorRegistry registry)
func NewColorGraphBuilder(nodeFactory *ColorGraphNodeFactory, lcaFinderFactory pkg_graph.LowestCommonAncestorFinder_Factory[*ColorGraphNode, any], registry *pkg_colors.ColorRegistry) *ColorGraphBuilder {
	this_ := &ColorGraphBuilder{}
	InitColorGraphBuilder(this_, nodeFactory, lcaFinderFactory, registry)
	return this_
}

func InitColorGraphBuilder(this_ *ColorGraphBuilder, nodeFactory *ColorGraphNodeFactory, lcaFinderFactory pkg_graph.LowestCommonAncestorFinder_Factory[*ColorGraphNode, any], registry *pkg_colors.ColorRegistry) {
	_ = nodeFactory
	_ = lcaFinderFactory
	_ = registry
	this_.initFields_ColorGraphBuilder()
	this_.registry = registry
	this_.nodeFactory = nodeFactory
	this_.lcaFinder = lcaFinderFactory.Create(j2gort.AsIface[pkg_graph.DiGraph[*ColorGraphNode, any]](this_.colorHoldsInstanceGraph))
	this_.topNode = j2gort.CastPtr[pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[*ColorGraphNode, any]](this_.colorHoldsInstanceGraph.CreateNode(this_.nodeFactory.CreateNode(pkg_colors.StaticGet_StandardColors_UNKNOWN())))
}

// Java method: ColorGraphBuilder.add(com.google.javascript.jscomp.disambiguate.ColorGraphNode flat)
func (this_ *ColorGraphBuilder) Add(flat *ColorGraphNode) {
	_ = flat
	this_.priv_ColorGraphBuilder_addInternal_6207(flat)
}

// Java method: ColorGraphBuilder.addAll(java.util.Collection<com.google.javascript.jscomp.disambiguate.ColorGraphNode> flats)
func (this_ *ColorGraphBuilder) AddAll(flats j2gort.Ext_java_util_Collection) {
	_ = flats
	flats.ForEach(func() j2gort.Ext_java_util_function_Consumer {
		bound_1 := this_
		return j2gort.NewLambda_Ext_java_util_function_Consumer(func(mrp0_ any) { bound_1.Add(j2gort.CastPtr[ColorGraphNode](mrp0_)) })
	}())
}

// Java method: ColorGraphBuilder.build()
func (this_ *ColorGraphBuilder) Build() pkg_graph.LinkedDirectedGraph[*ColorGraphNode, any] {
	this_.colorHoldsInstanceGraph.GetNodes().ForEach(func() j2gort.Ext_java_util_function_Consumer {
		bound_2 := this_
		return j2gort.NewLambda_Ext_java_util_function_Consumer(func(mrp0_ any) {
			bound_2.priv_ColorGraphBuilder_connectUnionWithAncestors_6205(j2gort.CastPtr[pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[*ColorGraphNode, any]](mrp0_))
		})
	}())
	var temp pkg_graph.LinkedDirectedGraph[*ColorGraphNode, any] = this_.colorHoldsInstanceGraph
	_ = temp
	this_.colorHoldsInstanceGraph = nil
	return temp
}

// Java method: ColorGraphBuilder.connectUnionWithAncestors(com.google.javascript.jscomp.graph.LinkedDirectedGraph.LinkedDiGraphNode<com.google.javascript.jscomp.disambiguate.ColorGraphNode,java.lang.Object> unionNode)
func (this_ *ColorGraphBuilder) priv_ColorGraphBuilder_connectUnionWithAncestors_6205(unionNode *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[*ColorGraphNode, any]) {
	_ = unionNode
	var flatUnion *ColorGraphNode = unionNode.GetValue()
	_ = flatUnion
	if !(flatUnion.GetColor().IsUnion()) {
		return
	}
	if !(!(unionNode.GetOutEdges().IsEmpty())) {
		j2gort.CheckState(false)
	}
	var unionElems []*pkg_colors.Color = flatUnion.GetColor().GetUnionElementsArray()
	_ = unionElems
	var graphNodesBuilder j2gort.Ext_com_google_common_collect_ImmutableSet_Builder = j2gort.Ext_com_google_common_collect_ImmutableSet_BuilderWithExpectedSize(int32(len(unionElems)))
	_ = graphNodesBuilder
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = int32(len(unionElems))
		_ = n
		for ; i < n; i = i + 1 {
			graphNodesBuilder.Add_1_any(j2gort.PtrToAny(this_.nodeFactory.CreateNode(unionElems[i])))
		}
	}
	var graphNodes j2gort.Ext_com_google_common_collect_ImmutableSet = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](graphNodesBuilder.Build())
	_ = graphNodes
	{
		cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.lcaFinder.FindAll(j2gort.AsIface[j2gort.Ext_java_util_Set](graphNodes))))
		for cur_3.HasNext() {
			rawElem_4 := cur_3.Next()
			lca := j2gort.CastPtr[ColorGraphNode](rawElem_4)
			_ = lca
			this_.priv_ColorGraphBuilder_connectSourceToDest_6208(j2gort.CheckNotNullPtr(j2gort.CastPtr[pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[*ColorGraphNode, any]](this_.colorHoldsInstanceGraph.GetNode(lca))), StaticGet_ColorGraphBuilder_EdgeReason_ALGEBRAIC(), unionNode)
		}
	}
}

// Java method: ColorGraphBuilder.addInternal(com.google.javascript.jscomp.colors.Color color)
func (this_ *ColorGraphBuilder) priv_ColorGraphBuilder_addInternal_6206(color *pkg_colors.Color) *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[*ColorGraphNode, any] {
	_ = color
	return this_.priv_ColorGraphBuilder_addInternal_6207(this_.nodeFactory.CreateNode(color))
}

// Java method: ColorGraphBuilder.addInternal(com.google.javascript.jscomp.disambiguate.ColorGraphNode node)
func (this_ *ColorGraphBuilder) priv_ColorGraphBuilder_addInternal_6207(node *ColorGraphNode) *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[*ColorGraphNode, any] {
	_ = node
	var flatNode *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[*ColorGraphNode, any] = j2gort.CastPtr[pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[*ColorGraphNode, any]](this_.colorHoldsInstanceGraph.GetNode(node))
	_ = flatNode
	if flatNode != nil {
		return flatNode
	}
	flatNode = j2gort.CastPtr[pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[*ColorGraphNode, any]](this_.colorHoldsInstanceGraph.CreateNode(node))
	if node.GetColor().IsUnion() {
		var alts []*pkg_colors.Color = node.GetColor().GetUnionElementsArray()
		_ = alts
		{
			var i int32 = int32(0)
			_ = i
			var n int32 = int32(len(alts))
			_ = n
			for ; i < n; i = i + 1 {
				this_.priv_ColorGraphBuilder_connectSourceToDest_6208(flatNode, StaticGet_ColorGraphBuilder_EdgeReason_ALGEBRAIC(), this_.priv_ColorGraphBuilder_addInternal_6206(alts[i]))
			}
		}
		return flatNode
	}
	var color *pkg_colors.Color = node.GetColor()
	_ = color
	var supertypes []*pkg_colors.Color = this_.registry.GetDisambiguationSupertypesArray(color)
	_ = supertypes
	if int32(len(supertypes)) == int32(0) {
		this_.priv_ColorGraphBuilder_connectSourceToDest_6208(this_.topNode, StaticGet_ColorGraphBuilder_EdgeReason_ALGEBRAIC(), flatNode)
	} else {
		{
			var i_2 int32 = int32(0)
			_ = i_2
			var n_2 int32 = int32(len(supertypes))
			_ = n_2
			for ; i_2 < n_2; i_2 = i_2 + 1 {
				this_.priv_ColorGraphBuilder_connectSourceToDest_6208(this_.priv_ColorGraphBuilder_addInternal_6206(supertypes[i_2]), StaticGet_ColorGraphBuilder_EdgeReason_CAN_HOLD(), flatNode)
			}
		}
	}
	var prototypes []*pkg_colors.Color = color.GetPrototypesArray()
	_ = prototypes
	{
		var i_3 int32 = int32(0)
		_ = i_3
		var n_3 int32 = int32(len(prototypes))
		_ = n_3
		for ; i_3 < n_3; i_3 = i_3 + 1 {
			this_.priv_ColorGraphBuilder_addInternal_6206(prototypes[i_3])
		}
	}
	var instanceColors []*pkg_colors.Color = color.GetInstanceColorsArray()
	_ = instanceColors
	{
		var i_4 int32 = int32(0)
		_ = i_4
		var n_4 int32 = int32(len(instanceColors))
		_ = n_4
		for ; i_4 < n_4; i_4 = i_4 + 1 {
			this_.priv_ColorGraphBuilder_addInternal_6206(instanceColors[i_4])
		}
	}
	return flatNode
}

// Java method: ColorGraphBuilder.connectSourceToDest(com.google.javascript.jscomp.graph.LinkedDirectedGraph.LinkedDiGraphNode<com.google.javascript.jscomp.disambiguate.ColorGraphNode,java.lang.Object> source, com.google.javascript.jscomp.disambiguate.ColorGraphBuilder.EdgeReason reason, com.google.javascript.jscomp.graph.LinkedDirectedGraph.LinkedDiGraphNode<com.google.javascript.jscomp.disambiguate.ColorGraphNode,java.lang.Object> dest)
func (this_ *ColorGraphBuilder) priv_ColorGraphBuilder_connectSourceToDest_6208(source *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[*ColorGraphNode, any], reason ColorGraphBuilder_EdgeReason, dest *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[*ColorGraphNode, any]) {
	_ = source
	_ = reason
	_ = dest
	if (j2gort.CheckNotNullPtr(source) == dest) || this_.colorHoldsInstanceGraph.IsConnectedInDirection_3_LinkedDirectedGraph_LinkedDiGraphNode_any_any__Ext_com_google_common_base_Predicate_LinkedDirectedGraph_LinkedDiGraphNode_any_any_(source, lambdaSingleton_color_graph_builder_1, dest) {
		return
	}
	this_.colorHoldsInstanceGraph.Connect_3_DiGraph_DiGraphNode_any_any__any_DiGraph_DiGraphNode_any_any_(pkg_graph.ToIface_DiGraph_DiGraphNode[*ColorGraphNode, any](source), j2gort.BoxNullableVal(reason), pkg_graph.ToIface_DiGraph_DiGraphNode[*ColorGraphNode, any](dest))
}

// Java enum: com.google.javascript.jscomp.disambiguate.ColorGraphBuilder.EdgeReason
type ColorGraphBuilder_EdgeReason uint8

const (
	enum_ColorGraphBuilder_EdgeReason_ALGEBRAIC ColorGraphBuilder_EdgeReason = 1
	enum_ColorGraphBuilder_EdgeReason_CAN_HOLD  ColorGraphBuilder_EdgeReason = 2
)

var enumBaseTable_ColorGraphBuilder_EdgeReason = [...]j2gort.EnumBase{
	{Name_: "ALGEBRAIC", Ordinal_: 0},
	{Name_: "CAN_HOLD", Ordinal_: 1},
}
var enumValues_ColorGraphBuilder_EdgeReason = []ColorGraphBuilder_EdgeReason{1, 2}

func StaticGet_ColorGraphBuilder_EdgeReason_ALGEBRAIC() ColorGraphBuilder_EdgeReason {
	return enum_ColorGraphBuilder_EdgeReason_ALGEBRAIC
}
func StaticGet_ColorGraphBuilder_EdgeReason_CAN_HOLD() ColorGraphBuilder_EdgeReason {
	return enum_ColorGraphBuilder_EdgeReason_CAN_HOLD
}

func ColorGraphBuilder_EdgeReason_Values() []ColorGraphBuilder_EdgeReason {
	return j2gort.CloneSlice(enumValues_ColorGraphBuilder_EdgeReason)
}

func ColorGraphBuilder_EdgeReason_ValueOf(name string) ColorGraphBuilder_EdgeReason {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_ColorGraphBuilder_EdgeReason {
		if enumBaseTable_ColorGraphBuilder_EdgeReason[i_].Name_ == name {
			return ColorGraphBuilder_EdgeReason(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.disambiguate.ColorGraphBuilder.EdgeReason." + name))
}

func (this_ ColorGraphBuilder_EdgeReason) IsJavaNil_() bool                              { return this_ == 0 }
func (ColorGraphBuilder_EdgeReason) HasIdentityEquals_()                                 {}
func (ColorGraphBuilder_EdgeReason) Implements_j2gort_Ext_java_lang_Enum()               {}
func (ColorGraphBuilder_EdgeReason) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ ColorGraphBuilder_EdgeReason) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (ColorGraphBuilder_EdgeReason) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.disambiguate.ColorGraphBuilder$EdgeReason", "EdgeReason"
}
func (this_ ColorGraphBuilder_EdgeReason) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_ColorGraphBuilder_EdgeReason[this_-1]
}
func (this_ ColorGraphBuilder_EdgeReason) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ ColorGraphBuilder_EdgeReason) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ ColorGraphBuilder_EdgeReason) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ColorGraphBuilder_EdgeReason[this_-1].Name_
}
func (this_ ColorGraphBuilder_EdgeReason) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_ColorGraphBuilder_EdgeReason[this_-1].Name_
}
func (this_ ColorGraphBuilder_EdgeReason) Equals(other_ any) bool {
	o_, ok_ := other_.(ColorGraphBuilder_EdgeReason)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ ColorGraphBuilder_EdgeReason) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_ColorGraphBuilder_EdgeReason[this_-1].HashCode()
}
func (this_ ColorGraphBuilder_EdgeReason) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ColorGraphBuilder_EdgeReason[this_-1].CompareTo(other_)
}
func (this_ ColorGraphBuilder_EdgeReason) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_ColorGraphBuilder_EdgeReason[this_-1].Name_)
}
func init() {
	vals_ := []any{ColorGraphBuilder_EdgeReason(1), ColorGraphBuilder_EdgeReason(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.disambiguate.ColorGraphBuilder.EdgeReason", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.disambiguate.ColorGraphBuilder$EdgeReason", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.disambiguate.ColorGraphBuilder.EdgeReason", "com.google.javascript.jscomp.disambiguate.ColorGraphBuilder$EdgeReason", func(v_ any) bool {
		return j2gort.InstanceOf[ColorGraphBuilder_EdgeReason](v_)
	})
}
