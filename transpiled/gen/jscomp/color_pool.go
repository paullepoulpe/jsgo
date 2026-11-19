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
// Port of com/google/javascript/jscomp/serialization/ColorPool.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_color_pool_1 j2gort.Ext_java_util_function_Function

func init() {
	lambdaSingleton_color_pool_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return j2gort.Ext_com_google_common_collect_ImmutableMap_CopyOf_1_Ext_java_util_Map(j2gort.AsIface[j2gort.Ext_java_util_Map](mrp0_))
	})
}

// Java class: com.google.javascript.jscomp.serialization.ColorPool
type ColorPool struct {
	idToColor     *j2gort.OrderedMap[pkg_colors.ColorId, *pkg_colors.Color] `gson:"idToColor"`
	colorRegistry *pkg_colors.ColorRegistry                                 `gson:"colorRegistry"`
	shardViews    j2gort.Ext_com_google_common_collect_ImmutableList        `gson:"shardViews"`
	idHash_       atomic.Int32
}

func (this_ *ColorPool) BindSelf_(s any) {
	_ = s
}
func (this_ *ColorPool) IsJavaNil_() bool { return this_ == nil }
func (*ColorPool) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.serialization.ColorPool", "ColorPool"
}
func (this_ *ColorPool) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ColorPool) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.ColorPool", "com.google.javascript.jscomp.serialization.ColorPool", func(v_ any) bool {
		return j2gort.InstanceOf[*ColorPool](v_)
	})
}

var static_ColorPool_PENDING_COLOR *pkg_colors.Color

func StaticGet_ColorPool_PENDING_COLOR() *pkg_colors.Color {
	clinit_ColorPool()
	return static_ColorPool_PENDING_COLOR
}
func StaticSet_ColorPool_PENDING_COLOR(v *pkg_colors.Color) {
	clinit_ColorPool()
	static_ColorPool_PENDING_COLOR = v
}

var clinitDone_ColorPool uint32
var clinitState_ColorPool j2gort.ClinitState

func clinit_ColorPool() {
	if atomic.LoadUint32(&clinitDone_ColorPool) != 0 {
		return
	}
	clinitSlow_ColorPool()
}

func clinitSlow_ColorPool() {
	if !j2gort.EnterClinit(&clinitState_ColorPool, &clinitDone_ColorPool) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ColorPool, &clinitDone_ColorPool)
	static_ColorPool_PENDING_COLOR = pkg_colors.Color_SingleBuilder().SetId(pkg_colors.ColorId_FromUnsigned_1_int32(int32(-559038737))).Build()
}

func (this_ *ColorPool) initFields_ColorPool() {
	_ = this_
}

// Java constructor: ColorPool(java.util.LinkedHashMap<com.google.javascript.jscomp.colors.ColorId,com.google.javascript.jscomp.colors.Color> idToColor, com.google.javascript.jscomp.colors.ColorRegistry colorRegistry, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.serialization.ColorPool.ShardView> shardViews)
func NewColorPool(idToColor *j2gort.OrderedMap[pkg_colors.ColorId, *pkg_colors.Color], colorRegistry *pkg_colors.ColorRegistry, shardViews j2gort.Ext_com_google_common_collect_ImmutableList) *ColorPool {
	clinit_ColorPool()
	this_ := &ColorPool{}
	InitColorPool(this_, idToColor, colorRegistry, shardViews)
	return this_
}

func InitColorPool(this_ *ColorPool, idToColor *j2gort.OrderedMap[pkg_colors.ColorId, *pkg_colors.Color], colorRegistry *pkg_colors.ColorRegistry, shardViews j2gort.Ext_com_google_common_collect_ImmutableList) {
	_ = idToColor
	_ = colorRegistry
	_ = shardViews
	this_.initFields_ColorPool()
	this_.idToColor = idToColor
	this_.colorRegistry = colorRegistry
	this_.shardViews = shardViews
}

// Java method: ColorPool.getColor(com.google.javascript.jscomp.colors.ColorId id)
func (this_ *ColorPool) GetColor(id pkg_colors.ColorId) *pkg_colors.Color {
	_ = id
	return this_.idToColor.Get(id)
}

// Java method: ColorPool.getRegistry()
func (this_ *ColorPool) GetRegistry() *pkg_colors.ColorRegistry {
	return this_.colorRegistry
}

// Java method: ColorPool.getOnlyShardForTesting()
func (this_ *ColorPool) GetOnlyShardForTesting() *ColorPool_ShardView {
	return j2gort.CastPtr[ColorPool_ShardView](j2gort.Ext_com_google_common_collect_Iterables_GetOnlyElement_1(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.shardViews)))
}

// Java static method: ColorPool.builder()
func ColorPool_S_Builder() *ColorPool_Builder {
	clinit_ColorPool()
	return NewColorPool_Builder()
}

// Java static method: ColorPool.fromOnlyShardForTesting(com.google.javascript.jscomp.serialization.TypePool typePool, com.google.javascript.jscomp.serialization.StringPool stringPool)
func ColorPool_FromOnlyShardForTesting(typePool *j2gort.Ext_com_google_javascript_jscomp_serialization_TypePool, stringPool *StringPool) *ColorPool {
	clinit_ColorPool()
	_ = typePool
	_ = stringPool
	return ColorPool_S_Builder().AddShardAnd(typePool, stringPool).priv_Builder_forTesting_7709().Build()
}

// Java static method: ColorPool.createTrimmedOffsetToId(com.google.javascript.jscomp.serialization.TypePool typePool)
func ColorPool_CreateTrimmedOffsetToId(typePool *j2gort.Ext_com_google_javascript_jscomp_serialization_TypePool) j2gort.Ext_com_google_common_collect_ImmutableList {
	clinit_ColorPool()
	_ = typePool
	var ids []pkg_colors.ColorId = make([]pkg_colors.ColorId, typePool.GetTypeCount())
	_ = ids
	{
		var i int32 = int32(0)
		_ = i
		for ; i < int32(len(ids)); i = i + 1 {
			var proto *j2gort.Ext_com_google_javascript_jscomp_serialization_TypeProto = typePool.GetType(i)
			_ = proto
			switch j2gort.CheckNotNullPtr(proto.GetKindCase()).Ordinal() {
			case int32(0) /* KindCase.OBJECT */ :
				{
					a_1 := ids
					i_2 := i
					a_1[i_2] = pkg_colors.ColorId_FromBytes_1_Ext_com_google_protobuf_ByteString(proto.GetObject().GetUuid())
				}
			case int32(1) /* KindCase.UNION */ :
			default:
				panic(NewMalformedTypedAstException_1(j2gort.PtrToAny(proto)))
			}
		}
	}
	{
		var i_2 int32 = int32(0)
		_ = i_2
		for ; i_2 < int32(len(ids)); i_2 = i_2 + 1 {
			var proto_2 *j2gort.Ext_com_google_javascript_jscomp_serialization_TypeProto = typePool.GetType(i_2)
			_ = proto_2
			switch j2gort.CheckNotNullPtr(proto_2.GetKindCase()).Ordinal() {
			case int32(0) /* KindCase.OBJECT */ :
			case int32(1) /* KindCase.UNION */ :
				{
					a_3 := ids
					i_4 := i_2
					a_3[i_4] = ColorPool_CreateUnionColorId(proto_2, ids)
				}
			default:
				panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError", j2gort.PtrToAny(proto_2))))
			}
		}
	}
	var seenIds j2gort.LocalOrderedSet[pkg_colors.ColorId] = j2gort.NewLocalOrderedSet[pkg_colors.ColorId]()
	_ = seenIds
	{
		var i_3 int32 = int32(0)
		_ = i_3
		for ; i_3 < int32(len(ids)); i_3 = i_3 + 1 {
			var proto_3 *j2gort.Ext_com_google_javascript_jscomp_serialization_TypeProto = typePool.GetType(i_3)
			_ = proto_3
			MalformedTypedAstException_CheckWellFormed_3(seenIds.Add(ids[i_3]), "Duplicate ID in single shard", j2gort.PtrToAny(proto_3))
		}
	}
	return j2gort.ImmListFromCompSlice(ids)
}

// Java static method: ColorPool.createUnionColorId(com.google.javascript.jscomp.serialization.TypeProto proto, com.google.javascript.jscomp.colors.ColorId[] allObjectIds)
func ColorPool_CreateUnionColorId(proto *j2gort.Ext_com_google_javascript_jscomp_serialization_TypeProto, allObjectIds []pkg_colors.ColorId) pkg_colors.ColorId {
	clinit_ColorPool()
	_ = proto
	_ = allObjectIds
	MalformedTypedAstException_CheckWellFormed_3((proto.GetUnion().GetUnionMemberCount() > int32(1)), "Union has too few members", j2gort.PtrToAny(proto))
	var members *j2gort.Ext_java_util_LinkedHashSet = j2gort.NewExt_java_util_LinkedHashSet()
	_ = members
	var unionMemberList j2gort.Ext_java_util_List = proto.GetUnion().GetUnionMemberList()
	_ = unionMemberList
	{
		var i int32 = int32(0)
		_ = i
		for ; i < unionMemberList.Size(); i = i + 1 {
			var memberPointer *j2gort.BoxInt = j2gort.CastPtr[j2gort.BoxInt](unionMemberList.Get(i))
			_ = memberPointer
			var memberId pkg_colors.ColorId
			_ = memberId
			if TypePointers_IsAxiomatic(j2gort.UnboxInt(memberPointer)) {
				memberId = j2gort.CastPtr[pkg_colors.Color](StaticGet_TypePointers_OFFSET_TO_AXIOMATIC_COLOR().Get(j2gort.UnboxInt(memberPointer))).GetId()
			} else {
				memberId = allObjectIds[TypePointers_TrimOffset(j2gort.UnboxInt(memberPointer))]
			}
			MalformedTypedAstException_CheckWellFormed_3((memberId != (pkg_colors.ColorId{})), "Union member not found", j2gort.PtrToAny(proto))
			members.Add_1(j2gort.BoxNullableVal(memberId))
		}
	}
	return pkg_colors.ColorId_Union(members)
}

// Java static method: ColorPool.validatePointer(int offset, com.google.javascript.jscomp.serialization.ColorPool.ShardView shard)
func ColorPool_ValidatePointer(offset int32, shard *ColorPool_ShardView) int32 {
	clinit_ColorPool()
	_ = offset
	_ = shard
	MalformedTypedAstException_CheckWellFormed_3(((int32(0) <= offset) && (offset < TypePointers_UntrimOffset(int32(len(shard.trimmedOffsetToId))))), "Pointer offset outside of shard", j2gort.ValueOfInt(offset))
	return offset
}

// Java class: com.google.javascript.jscomp.serialization.ColorPool.ShardView
type ColorPool_ShardView struct {
	trimmedOffsetToId []pkg_colors.ColorId                                            `gson:"trimmedOffsetToId"`
	typePool          *j2gort.Ext_com_google_javascript_jscomp_serialization_TypePool `gson:"typePool"`
	stringPool        *StringPool                                                     `gson:"stringPool"`
	colorPool         *ColorPool                                                      `gson:"colorPool"`
	idHash_           atomic.Int32
}

func (this_ *ColorPool_ShardView) BindSelf_(s any) {
	_ = s
}
func (this_ *ColorPool_ShardView) IsJavaNil_() bool { return this_ == nil }
func (*ColorPool_ShardView) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.serialization.ColorPool$ShardView", "ShardView"
}
func (this_ *ColorPool_ShardView) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ColorPool_ShardView) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.ColorPool.ShardView", "com.google.javascript.jscomp.serialization.ColorPool$ShardView", func(v_ any) bool {
		return j2gort.InstanceOf[*ColorPool_ShardView](v_)
	})
}

func (this_ *ColorPool_ShardView) initFields_ColorPool_ShardView() {
	_ = this_
}

// Java constructor: ShardView(com.google.javascript.jscomp.serialization.TypePool typePool, com.google.javascript.jscomp.serialization.StringPool stringPool, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.colors.ColorId> trimmedOffsetToId)
func NewColorPool_ShardView(typePool *j2gort.Ext_com_google_javascript_jscomp_serialization_TypePool, stringPool *StringPool, trimmedOffsetToId j2gort.Ext_com_google_common_collect_ImmutableList) *ColorPool_ShardView {
	this_ := &ColorPool_ShardView{}
	InitColorPool_ShardView(this_, typePool, stringPool, trimmedOffsetToId)
	return this_
}

func InitColorPool_ShardView(this_ *ColorPool_ShardView, typePool *j2gort.Ext_com_google_javascript_jscomp_serialization_TypePool, stringPool *StringPool, trimmedOffsetToId j2gort.Ext_com_google_common_collect_ImmutableList) {
	_ = typePool
	_ = stringPool
	_ = trimmedOffsetToId
	this_.initFields_ColorPool_ShardView()
	this_.typePool = typePool
	this_.stringPool = stringPool
	this_.trimmedOffsetToId = j2gort.ImmListToSliceNullable[

	// Java method: ShardView.getColor(int pointer)
	pkg_colors.ColorId](trimmedOffsetToId)
}

func (this_ *ColorPool_ShardView) GetColor(pointer int32) *pkg_colors.Color {
	_ = pointer
	if !(this_.colorPool != nil) {
		j2gort.CheckStateMsg(false, this_)
	}
	return this_.colorPool.GetColor(this_.priv_ShardView_getId_7706(pointer))
}

// Java method: ShardView.getId(int untrimmedOffset)
func (this_ *ColorPool_ShardView) priv_ShardView_getId_7706(untrimmedOffset int32) pkg_colors.ColorId {
	_ = untrimmedOffset
	if TypePointers_IsAxiomatic(untrimmedOffset) {
		return j2gort.CastPtr[pkg_colors.Color](StaticGet_TypePointers_OFFSET_TO_AXIOMATIC_COLOR().Get(untrimmedOffset)).GetId()
	} else {
		return this_.trimmedOffsetToId[TypePointers_TrimOffset(untrimmedOffset)]
	}
}

// Java method: ShardView.updateStateAfterColorPoolIsBuilt(com.google.javascript.jscomp.serialization.ColorPool colorPool)
func (this_ *ColorPool_ShardView) priv_ShardView_updateStateAfterColorPoolIsBuilt_7707(colorPool *ColorPool) {
	_ = colorPool
	this_.colorPool = colorPool
	this_.typePool = nil
	this_.stringPool = nil
}

// Java class: com.google.javascript.jscomp.serialization.ColorPool.Builder
type ColorPool_Builder struct {
	protoToShard            *pkg_base.LinkedIdentityHashMap[*j2gort.Ext_com_google_javascript_jscomp_serialization_TypePool, *ColorPool_ShardView] `gson:"protoToShard"`
	idToColor               j2gort.OrderedMap[pkg_colors.ColorId, *pkg_colors.Color]                                                               `gson:"idToColor"`
	registry                *pkg_colors.ColorRegistry_Builder                                                                                      `gson:"registry"`
	idToProto               *j2gort.Ext_com_google_common_collect_HashBasedTable                                                                   `gson:"idToProto"`
	reconcilationDebugStack *j2gort.Ext_java_util_ArrayDeque                                                                                       `gson:"reconcilationDebugStack"`
	scratchOwnProperties    *j2gort.Ext_java_util_LinkedHashSet                                                                                    `gson:"scratchOwnProperties"`
	idHash_                 atomic.Int32
	forTesting              bool `gson:"forTesting"`
}

func (this_ *ColorPool_Builder) BindSelf_(s any) {
	_ = s
}
func (this_ *ColorPool_Builder) IsJavaNil_() bool { return this_ == nil }
func (*ColorPool_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.serialization.ColorPool$Builder", "Builder"
}
func (this_ *ColorPool_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ColorPool_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.ColorPool.Builder", "com.google.javascript.jscomp.serialization.ColorPool$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*ColorPool_Builder](v_)
	})
}

func (this_ *ColorPool_Builder) initFields_ColorPool_Builder() {
	_ = this_
	this_.protoToShard = pkg_base.NewLinkedIdentityHashMap[*j2gort.Ext_com_google_javascript_jscomp_serialization_TypePool, *ColorPool_ShardView]()
	this_.idToColor = j2gort.NewOrderedMap[pkg_colors.ColorId, *pkg_colors.Color]()
	this_.registry = pkg_colors.ColorRegistry_S_Builder()
	this_.idToProto = j2gort.Ext_com_google_common_collect_HashBasedTable_Create_0()

	this_.reconcilationDebugStack = j2gort.NewExt_java_util_ArrayDeque()
	this_.scratchOwnProperties = j2gort.NewExt_java_util_LinkedHashSet()
}

// Java constructor: Builder()
func NewColorPool_Builder() *ColorPool_Builder {
	this_ := &ColorPool_Builder{}
	InitColorPool_Builder(this_)
	return this_
}

func InitColorPool_Builder(this_ *ColorPool_Builder) {
	this_.initFields_ColorPool_Builder()
	this_.idToColor.PutAllJavaMap(j2gort.AsIface[j2gort.Ext_java_util_Map](pkg_colors.StaticGet_StandardColors_AXIOMATIC_COLORS()))
}

// Java method: Builder.addShardAnd(com.google.javascript.jscomp.serialization.TypePool typePool, com.google.javascript.jscomp.serialization.StringPool stringPool)
func (this_ *ColorPool_Builder) AddShardAnd(typePool *j2gort.Ext_com_google_javascript_jscomp_serialization_TypePool, stringPool *StringPool) *ColorPool_Builder {
	_ = typePool
	_ = stringPool
	this_.AddShard(typePool, stringPool)
	return this_
}

// Java method: Builder.forTesting()
func (this_ *ColorPool_Builder) priv_Builder_forTesting_7709() *ColorPool_Builder {
	this_.forTesting = true
	return this_
}

// Java method: Builder.addShard(com.google.javascript.jscomp.serialization.TypePool typePool, com.google.javascript.jscomp.serialization.StringPool stringPool)
func (this_ *ColorPool_Builder) AddShard(typePool *j2gort.Ext_com_google_javascript_jscomp_serialization_TypePool, stringPool *StringPool) *ColorPool_ShardView {
	_ = typePool
	_ = stringPool
	if !(this_.idToProto.IsEmpty()) {
		j2gort.CheckStateStr(false, "build has already been called")
	}
	var existing *ColorPool_ShardView = this_.protoToShard.Get(typePool)
	_ = existing
	if existing != nil {
		if !(typePool == j2gort.Ext_com_google_javascript_jscomp_serialization_TypePool_GetDefaultInstance()) {
			j2gort.CheckStateMsg(false, j2gort.PtrToAny(typePool))
		}
		return existing
	}
	var trimmedOffsetToId j2gort.Ext_com_google_common_collect_ImmutableList = ColorPool_CreateTrimmedOffsetToId(typePool)
	_ = trimmedOffsetToId
	var val_5 ColorPool_ShardView
	InitColorPool_ShardView(&val_5, typePool, stringPool, trimmedOffsetToId)
	var shard *ColorPool_ShardView = &val_5
	_ = shard
	this_.protoToShard.Put(typePool, shard)
	if typePool.HasDebugInfo() {
		{
			cur_6 := j2gort.NewIterCursor(j2gort.CheckNotNull(typePool.GetDebugInfo().GetMismatchList()))
			for cur_6.HasNext() {
				rawElem_7 := cur_6.Next()
				m := j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_serialization_TypePool_DebugInfo_Mismatch](rawElem_7)
				_ = m
				{
					cur_8 := j2gort.NewIterCursor(j2gort.CheckNotNull(m.GetInvolvedColorList()))
					for cur_8.HasNext() {
						rawElem_9 := cur_8.Next()
						pointer := j2gort.CastPtr[j2gort.BoxInt](rawElem_9)
						_ = pointer
						this_.registry.AddMismatchLocation(shard.priv_ShardView_getId_7706(j2gort.UnboxInt(pointer)), m.GetSourceRef())
					}
				}
			}
		}
	}
	return shard
}

// Java method: Builder.build()
func (this_ *ColorPool_Builder) Build() *ColorPool {
	if !(this_.idToProto.IsEmpty()) {
		j2gort.CheckStateStr(false, "build has already been called")
	}
	this_.priv_Builder_forEachShard_7712(j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_10 any) {
		shard := j2gort.CastPtr[ColorPool_ShardView](lp_10)
		_ = shard
		{
			var i int32 = int32(0)
			_ = i
			for ; i < shard.typePool.GetTypeCount(); i = i + 1 {
				var id pkg_colors.ColorId = shard.trimmedOffsetToId[i]
				_ = id
				j2gort.TablePutColorId(this_.idToProto, id, j2gort.PtrToAny(shard), j2gort.PtrToAny(shard.typePool.GetType(i)))
			}
		}
	}))
	{
		cur_11 := j2gort.NewMapKeyCursor(pkg_colors.StaticGet_StandardColors_AXIOMATIC_COLORS())
		for cur_11.HasNext() {
			rawElem_12 := cur_11.Next()
			id := j2gort.Cast[pkg_colors.ColorId](rawElem_12)
			_ = id
			MalformedTypedAstException_CheckWellFormed_3(!(j2gort.TableContainsRowColorId(this_.idToProto, id)), "Found serialized definiton for axiomatic color", j2gort.BoxNullableVal(id))
		}
	}
	{
		cur_13 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.idToProto.RowKeySet()))
		for cur_13.HasNext() {
			rawElem_14 := cur_13.Next()
			id_2 := j2gort.Cast[pkg_colors.ColorId](rawElem_14)
			_ = id_2
			this_.priv_Builder_lookupOrReconcileColor_7713(id_2)
		}
	}
	{
		cur_15 := j2gort.NewIterCursor(j2gort.CheckNotNull(pkg_colors.StaticGet_ColorRegistry_REQUIRED_IDS()))
		for cur_15.HasNext() {
			rawElem_16 := cur_15.Next()
			colorId := j2gort.Cast[pkg_colors.ColorId](rawElem_16)
			_ = colorId
			this_.registry.SetNativeColor(j2gort.OrderedMapComputeIfAbsent(&this_.idToColor, colorId, nil, func(lp_21 pkg_colors.ColorId) *pkg_colors.Color {
				unused := lp_21
				_ = unused
				return pkg_colors.Color_SingleBuilder().SetId(colorId).Build()
			}),
			)
		}
	}
	this_.priv_Builder_forEachShard_7712(j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_22 any) {
		shard := j2gort.CastPtr[ColorPool_ShardView](lp_22)
		_ = shard
		{
			cur_23 := j2gort.NewIterCursor(j2gort.CheckNotNull(shard.typePool.GetDisambiguationEdgesList()))
			for cur_23.HasNext() {
				rawElem_24 := cur_23.Next()
				edge := j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_serialization_SubtypingEdge](rawElem_24)
				_ = edge
				this_.registry.AddDisambiguationEdge(this_.idToColor.Get(shard.priv_ShardView_getId_7706(ColorPool_ValidatePointer(edge.GetSubtype(), shard))), this_.idToColor.Get(shard.priv_ShardView_getId_7706(ColorPool_ValidatePointer(edge.GetSupertype(), shard))))
			}
		}
	}))
	var shardViewsForTesting *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = (&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{})
	_ = shardViewsForTesting
	if this_.forTesting {
		this_.priv_Builder_forEachShard_7712(func() j2gort.Ext_java_util_function_Consumer {
			bound_25 := shardViewsForTesting
			return j2gort.NewLambda_Ext_java_util_function_Consumer(func(mrp0_ any) { bound_25.Add_1_any(mrp0_) })
		}())
	}
	clinit_ColorPool()
	var val_26 ColorPool
	InitColorPool(&val_26, &this_.idToColor, this_.registry.Build(), shardViewsForTesting.BuildImmList_())
	var colorPool *ColorPool = &val_26
	_ = colorPool
	this_.priv_Builder_forEachShard_7712(j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_27 any) {
		shard := j2gort.CastPtr[ColorPool_ShardView](lp_27)
		_ = shard
		shard.priv_ShardView_updateStateAfterColorPoolIsBuilt_7707(colorPool)
	}))
	return colorPool
}

// Java method: Builder.forEachShard(java.util.function.Consumer<com.google.javascript.jscomp.serialization.ColorPool.ShardView> fn)
func (this_ *ColorPool_Builder) priv_Builder_forEachShard_7712(fn j2gort.Ext_java_util_function_Consumer) {
	_ = fn
	this_.protoToShard.ForEach(j2gort.NewLambda_Ext_java_util_function_BiConsumer(func(lp_28 any, lp_29 any) {
		unused := j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_serialization_TypePool](lp_28)
		_ = unused
		shard := j2gort.CastPtr[ColorPool_ShardView](lp_29)
		_ = shard
		fn.Accept(j2gort.PtrToAny(shard))
	}))
}

// Java method: Builder.lookupOrReconcileColor(com.google.javascript.jscomp.colors.ColorId id)
func (this_ *ColorPool_Builder) priv_Builder_lookupOrReconcileColor_7713(id pkg_colors.ColorId) *pkg_colors.Color {
	_ = id
	j2gort.DequeAddLastColorIdVoid(this_.reconcilationDebugStack, id)
	defer func() {
		j2gort.DequeRemoveLastVoid(this_.reconcilationDebugStack)
	}()
	{
		var existing *pkg_colors.Color = this_.idToColor.PutIfAbsent(id, StaticGet_ColorPool_PENDING_COLOR())
		_ = existing
		if existing != nil {
			if existing == StaticGet_ColorPool_PENDING_COLOR() {
				panic(NewMalformedTypedAstException_1("Cyclic Color structure detected: " + j2gort.StrVal(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableList](this_.reconcilationDebugStack.Stream().Map(func() j2gort.Ext_java_util_function_Function {
					bound_30 := this_.idToProto
					return j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return bound_30.Row(mrp0_) })
				}()).Map(lambdaSingleton_color_pool_1).Collect_1(j2gort.Ext_com_google_common_collect_ImmutableList_ToImmutableList())))))
			}
			return existing
		}
		var viewToProto j2gort.Ext_java_util_Map = j2gort.TableRowColorId(this_.idToProto, id)
		_ = viewToProto
		var sample *j2gort.Ext_com_google_javascript_jscomp_serialization_TypeProto = j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_serialization_TypeProto](j2gort.Ext_com_google_common_collect_Iterables_GetFirst(j2gort.AsIface[j2gort.Ext_java_util_Collection](viewToProto.Values()), nil))
		_ = sample
		j2gort.CheckNotNullPtr(sample)
		var result *pkg_colors.Color = nil
		_ = result
		switch j2gort.CheckNotNullPtr(sample.GetKindCase()).Ordinal() {
		case int32(0) /* KindCase.OBJECT */ :
			result = this_.priv_Builder_reconcileObjectProtos_7714(id, viewToProto)
		case int32(1) /* KindCase.UNION */ :
			result = this_.priv_Builder_reconcileUnionProtos_7715(id, viewToProto)
		default:
			panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError", j2gort.PtrToAny(sample))))
		}
		if !(result != nil) {
			j2gort.CheckStateMsg(false, j2gort.BoxNullableVal(id))
		}
		this_.idToColor.Put(id, result)
		return result
	}
	panic("unreachable")
}

// Java method: Builder.reconcileObjectProtos(com.google.javascript.jscomp.colors.ColorId id, java.util.Map<com.google.javascript.jscomp.serialization.ColorPool.ShardView,com.google.javascript.jscomp.serialization.TypeProto> viewToProto)
func (this_ *ColorPool_Builder) priv_Builder_reconcileObjectProtos_7714(id pkg_colors.ColorId, viewToProto j2gort.Ext_java_util_Map) *pkg_colors.Color {
	_ = id
	_ = viewToProto
	var totalInstanceCount int32 = int32(0)
	_ = totalInstanceCount
	var totalPrototypeCount int32 = int32(0)
	_ = totalPrototypeCount
	var totalOwnPropCount int32 = int32(0)
	_ = totalOwnPropCount
	{
		cur_31 := j2gort.NewMapValCursor(viewToProto)
		for cur_31.HasNext() {
			rawElem_32 := cur_31.Next()
			proto := j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_serialization_TypeProto](rawElem_32)
			_ = proto
			if !(proto.HasObject()) {
				j2gort.CheckState(false)
			}
			var objProto *j2gort.Ext_com_google_javascript_jscomp_serialization_ObjectTypeProto = proto.GetObject()
			_ = objProto
			totalInstanceCount = (totalInstanceCount + objProto.GetInstanceTypeCount())
			totalPrototypeCount = (totalPrototypeCount + objProto.GetPrototypeCount())
			totalOwnPropCount = (totalOwnPropCount + objProto.GetOwnPropertyCount())
		}
	}
	var instanceColorsSet *j2gort.Ext_java_util_LinkedHashSet
	_ = instanceColorsSet
	if totalInstanceCount > int32(1) {
		instanceColorsSet = j2gort.NewExt_java_util_LinkedHashSet_1_int32(totalInstanceCount)
	} else {
		instanceColorsSet = nil
	}
	var singleInstanceColor *pkg_colors.Color = nil
	_ = singleInstanceColor
	var prototypesSet *j2gort.Ext_java_util_LinkedHashSet
	_ = prototypesSet
	if totalPrototypeCount > int32(1) {
		prototypesSet = j2gort.NewExt_java_util_LinkedHashSet_1_int32(totalPrototypeCount)
	} else {
		prototypesSet = nil
	}
	var singlePrototypeColor *pkg_colors.Color = nil
	_ = singlePrototypeColor
	var isClosureAssert pkg_base.Tri = pkg_base.StaticGet_Tri_UNKNOWN()
	_ = isClosureAssert
	var isConstructor bool = false
	_ = isConstructor
	var isInvalidating bool = false
	_ = isInvalidating
	var propertiesKeepOriginalName bool = false
	_ = propertiesKeepOriginalName
	{
		cur_33 := j2gort.NewMapEntryCursor(viewToProto)
		for cur_33.NextEntry() {
			var shard *ColorPool_ShardView = j2gort.CastPtr[ColorPool_ShardView](cur_33.GetKey())
			_ = shard
			var proto_2 *j2gort.Ext_com_google_javascript_jscomp_serialization_TypeProto = j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_serialization_TypeProto](cur_33.GetValue())
			_ = proto_2
			var objProto_2 *j2gort.Ext_com_google_javascript_jscomp_serialization_ObjectTypeProto = proto_2.GetObject()
			_ = objProto_2
			var instanceTypeIds j2gort.Ext_java_util_List = objProto_2.GetInstanceTypeList()
			_ = instanceTypeIds
			{
				var i int32 = int32(0)
				_ = i
				var n int32 = instanceTypeIds.Size()
				_ = n
				for ; i < n; i = i + 1 {
					var instanceColor *pkg_colors.Color = this_.priv_Builder_lookupOrReconcileColor_7713(shard.priv_ShardView_getId_7706(j2gort.UnboxInt(j2gort.CastPtr[j2gort.BoxInt](instanceTypeIds.Get(i)))))
					_ = instanceColor
					if instanceColorsSet != nil {
						instanceColorsSet.Add_1(j2gort.PtrToAny(instanceColor))
					} else {
						singleInstanceColor = instanceColor
					}
				}
			}
			var isClosureAssertBool bool = objProto_2.GetClosureAssert()
			_ = isClosureAssertBool
			MalformedTypedAstException_CheckWellFormed_3((isClosureAssert.ToBoolean(isClosureAssertBool) == isClosureAssertBool), "Inconsistent values for closure_assert", j2gort.PtrToAny(objProto_2))
			isClosureAssert = pkg_base.Tri_ForBoolean(isClosureAssertBool)
			isConstructor = func() bool { ba_34 := isConstructor; bb_35 := objProto_2.GetMarkedConstructor(); return ba_34 || bb_35 }()
			isInvalidating = func() bool { ba_36 := isInvalidating; bb_37 := objProto_2.GetIsInvalidating(); return ba_36 || bb_37 }()
			propertiesKeepOriginalName = func() bool {
				ba_38 := propertiesKeepOriginalName
				bb_39 := objProto_2.GetPropertiesKeepOriginalName()
				return ba_38 || bb_39
			}()
			var prototypeIds j2gort.Ext_java_util_List = objProto_2.GetPrototypeList()
			_ = prototypeIds
			{
				var i_2 int32 = int32(0)
				_ = i_2
				var n_2 int32 = prototypeIds.Size()
				_ = n_2
				for ; i_2 < n_2; i_2 = i_2 + 1 {
					var prototypeColor *pkg_colors.Color = this_.priv_Builder_lookupOrReconcileColor_7713(shard.priv_ShardView_getId_7706(j2gort.UnboxInt(j2gort.CastPtr[j2gort.BoxInt](prototypeIds.Get(i_2)))))
					_ = prototypeColor
					if prototypesSet != nil {
						prototypesSet.Add_1(j2gort.PtrToAny(prototypeColor))
					} else {
						singlePrototypeColor = prototypeColor
					}
				}
			}
		}
	}
	var instanceColorsArr []*pkg_colors.Color = nil
	_ = instanceColorsArr
	if totalInstanceCount == int32(0) {
		instanceColorsArr = pkg_colors.StaticGet_Color_EMPTY_COLORS()
	} else if instanceColorsSet == nil {
		instanceColorsArr = []*pkg_colors.Color{singleInstanceColor}
	} else {
		instanceColorsArr = j2gort.ColToPtrArray[pkg_colors.Color](instanceColorsSet, pkg_colors.StaticGet_Color_EMPTY_COLORS())
	}
	var prototypesArr []*pkg_colors.Color = nil
	_ = prototypesArr
	if totalPrototypeCount == int32(0) {
		prototypesArr = pkg_colors.StaticGet_Color_EMPTY_COLORS()
	} else if prototypesSet == nil {
		prototypesArr = []*pkg_colors.Color{singlePrototypeColor}
	} else {
		prototypesArr = j2gort.ColToPtrArray[pkg_colors.Color](prototypesSet, pkg_colors.StaticGet_Color_EMPTY_COLORS())
	}
	var ownPropertiesArr []string = nil
	_ = ownPropertiesArr
	if totalOwnPropCount == int32(0) {
		ownPropertiesArr = pkg_colors.StaticGet_Color_EMPTY_STRINGS()
	} else if totalOwnPropCount == int32(1) {
		var onlyProp string = j2gort.NullString
		_ = onlyProp
		{
			cur_40 := j2gort.NewMapEntryCursor(viewToProto)
			for cur_40.NextEntry() {
				var objProto_3 *j2gort.Ext_com_google_javascript_jscomp_serialization_ObjectTypeProto = j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_serialization_TypeProto](cur_40.GetValue()).GetObject()
				_ = objProto_3
				if objProto_3.GetOwnPropertyCount() > int32(0) {
					onlyProp = j2gort.CastPtr[ColorPool_ShardView](cur_40.GetKey()).stringPool.Get(objProto_3.GetOwnProperty(int32(0)))
					break
				}
			}
		}
		ownPropertiesArr = []string{j2gort.CheckNotNullString(onlyProp)}
	} else {
		var ownPropsSet *j2gort.Ext_java_util_LinkedHashSet = this_.scratchOwnProperties
		_ = ownPropsSet
		ownPropsSet.Clear()
		{
			cur_41 := j2gort.NewMapEntryCursor(viewToProto)
			for cur_41.NextEntry() {
				var shard_2 *ColorPool_ShardView = j2gort.CastPtr[ColorPool_ShardView](cur_41.GetKey())
				_ = shard_2
				var objProto_4 *j2gort.Ext_com_google_javascript_jscomp_serialization_ObjectTypeProto = j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_serialization_TypeProto](cur_41.GetValue()).GetObject()
				_ = objProto_4
				{
					var i_3 int32 = int32(0)
					_ = i_3
					var n_3 int32 = objProto_4.GetOwnPropertyCount()
					_ = n_3
					for ; i_3 < n_3; i_3 = i_3 + 1 {
						j2gort.SetAddStr(ownPropsSet, shard_2.stringPool.Get(objProto_4.GetOwnProperty(i_3)))
					}
				}
			}
		}
		ownPropertiesArr = j2gort.ColToStringArray(ownPropsSet, pkg_colors.StaticGet_Color_EMPTY_STRINGS())
		ownPropsSet.Clear()
	}
	return pkg_colors.Color_SingleBuilder().SetId(id).SetInstanceColorsArray(instanceColorsArr).SetPrototypesArray(prototypesArr).SetOwnPropertiesArray(ownPropertiesArr).SetClosureAssert(isClosureAssert.ToBoolean(false)).SetConstructor(isConstructor).SetInvalidating(isInvalidating).SetPropertiesKeepOriginalName(propertiesKeepOriginalName).Build()
}

// Java method: Builder.reconcileUnionProtos(com.google.javascript.jscomp.colors.ColorId id, java.util.Map<com.google.javascript.jscomp.serialization.ColorPool.ShardView,com.google.javascript.jscomp.serialization.TypeProto> viewToProto)
func (this_ *ColorPool_Builder) priv_Builder_reconcileUnionProtos_7715(id pkg_colors.ColorId, viewToProto j2gort.Ext_java_util_Map) *pkg_colors.Color {
	_ = id
	_ = viewToProto
	if viewToProto.Size() == int32(1) {
		var onlyEntry j2gort.Ext_java_util_Map_Entry = j2gort.AsIface[j2gort.Ext_java_util_Map_Entry](j2gort.AsIface[j2gort.Ext_java_util_Iterator](j2gort.AsIface[j2gort.Ext_java_util_Set](viewToProto.EntrySet()).Iterator()).Next())
		_ = onlyEntry
		var shard *ColorPool_ShardView = j2gort.CastPtr[ColorPool_ShardView](onlyEntry.GetKey())
		_ = shard
		var proto *j2gort.Ext_com_google_javascript_jscomp_serialization_TypeProto = j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_serialization_TypeProto](onlyEntry.GetValue())
		_ = proto
		if !(proto.HasUnion()) {
			j2gort.CheckStateMsg(false, j2gort.PtrToAny(proto))
		}
		var unionMemberList j2gort.Ext_java_util_List = proto.GetUnion().GetUnionMemberList()
		_ = unionMemberList
		var memberCount int32 = unionMemberList.Size()
		_ = memberCount
		if memberCount == int32(2) {
			var m0 *pkg_colors.Color = this_.priv_Builder_lookupOrReconcileColor_7713(shard.priv_ShardView_getId_7706(j2gort.UnboxInt(j2gort.CastPtr[j2gort.BoxInt](unionMemberList.Get(int32(0))))))
			_ = m0
			MalformedTypedAstException_CheckWellFormed_3(!(m0.IsUnion()), "Reconciling union with non-union", j2gort.PtrToAny(proto))
			var m1 *pkg_colors.Color = this_.priv_Builder_lookupOrReconcileColor_7713(shard.priv_ShardView_getId_7706(j2gort.UnboxInt(j2gort.CastPtr[j2gort.BoxInt](unionMemberList.Get(int32(1))))))
			_ = m1
			MalformedTypedAstException_CheckWellFormed_3(!(m1.IsUnion()), "Reconciling union with non-union", j2gort.PtrToAny(proto))
			var result *pkg_colors.Color = pkg_colors.Color_CreateUnion2(m0, m1)
			_ = result
			{
				var iifeCond_1 bool
				{

					rv_42 := id
					av_43 := result.GetId()
					_ = j2gort.CheckNotNullVal(rv_42)
					iifeCond_1 = rv_42 == av_43
				}
				if !(iifeCond_1) {
					j2gort.CheckStateFmt(false, "%s == %s", j2gort.BoxNullableVal(id), j2gort.PtrToAny(result))
				}
			}
			return result
		}
		var members []*pkg_colors.Color = make([]*pkg_colors.Color, memberCount)
		_ = members
		{
			var i int32 = int32(0)
			_ = i
			for ; i < memberCount; i = i + 1 {
				var member *pkg_colors.Color = this_.priv_Builder_lookupOrReconcileColor_7713(shard.priv_ShardView_getId_7706(j2gort.UnboxInt(j2gort.CastPtr[j2gort.BoxInt](unionMemberList.Get(i)))))
				_ = member
				MalformedTypedAstException_CheckWellFormed_3(!(member.IsUnion()), "Reconciling union with non-union", j2gort.PtrToAny(proto))
				{
					a_44 := members
					i_45 := i
					a_44[i_45] = member
				}
			}
		}
		var result_2 *pkg_colors.Color = pkg_colors.Color_CreateUnion_1_Color(members)
		_ = result_2
		{
			var iifeCond_1 bool
			{

				rv_46 := id
				av_47 := result_2.GetId()
				_ = j2gort.CheckNotNullVal(rv_46)
				iifeCond_1 = rv_46 == av_47
			}
			if !(iifeCond_1) {
				j2gort.CheckStateFmt(false, "%s == %s", j2gort.BoxNullableVal(id), j2gort.PtrToAny(result_2))
			}
		}
		return result_2
	}
	var union *j2gort.Ext_java_util_LinkedHashSet = j2gort.NewExt_java_util_LinkedHashSet()
	_ = union
	{
		cur_48 := j2gort.NewMapEntryCursor(viewToProto)
		for cur_48.NextEntry() {
			var shard_2 *ColorPool_ShardView = j2gort.CastPtr[ColorPool_ShardView](cur_48.GetKey())
			_ = shard_2
			var proto_2 *j2gort.Ext_com_google_javascript_jscomp_serialization_TypeProto = j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_serialization_TypeProto](cur_48.GetValue())
			_ = proto_2
			if !(proto_2.HasUnion()) {
				j2gort.CheckStateMsg(false, j2gort.PtrToAny(proto_2))
			}
			var unionMemberList_2 j2gort.Ext_java_util_List = proto_2.GetUnion().GetUnionMemberList()
			_ = unionMemberList_2
			{
				var i_2 int32 = int32(0)
				_ = i_2
				var n int32 = unionMemberList_2.Size()
				_ = n
				for ; i_2 < n; i_2 = i_2 + 1 {
					var memberPointer *j2gort.BoxInt = j2gort.CastPtr[j2gort.BoxInt](unionMemberList_2.Get(i_2))
					_ = memberPointer
					var memberId pkg_colors.ColorId = shard_2.priv_ShardView_getId_7706(j2gort.UnboxInt(memberPointer))
					_ = memberId
					var member_2 *pkg_colors.Color = this_.priv_Builder_lookupOrReconcileColor_7713(memberId)
					_ = member_2
					MalformedTypedAstException_CheckWellFormed_3(!(member_2.IsUnion()), "Reconciling union with non-union", j2gort.PtrToAny(proto_2))
					union.Add_1(j2gort.PtrToAny(member_2))
				}
			}
		}
	}
	var result_3 *pkg_colors.Color = pkg_colors.Color_CreateUnion_1_Ext_java_util_Set(union)
	_ = result_3
	{
		var iifeCond_1 bool
		{

			rv_49 := id
			av_50 := result_3.GetId()
			_ = j2gort.CheckNotNullVal(rv_49)
			iifeCond_1 = rv_49 == av_50
		}
		if !(iifeCond_1) {
			j2gort.CheckStateFmt(false, "%s == %s", j2gort.BoxNullableVal(id), j2gort.PtrToAny(result_3))
		}
	}
	return result_3
}
