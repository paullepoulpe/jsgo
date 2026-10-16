// Port of com/google/javascript/jscomp/serialization/ColorPool.java @ closure-compiler v20260415

package jscomp

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// ColorPool is a set of Colors reconstructed from possibly many TypePool
// protos.
//
// Protos representing the same Color are reconciled while the pool is
// built. Reconciliation ensures that, within a pool, each ColorId
// corresponds to at most one Color.
//
// It's possible for two Color objects with the same ID to have different
// definitions between different ColorPools, but mixing multiple ColorPools
// is likely a design error.
type ColorPool struct {
	// This could be an ImmutableMap, but because it's very large, creating an ImmutableMap copy has
	// caused OOMs in large projects. So we just use a LinkedHashMap and hide the mutability behind
	// accessor methods.
	idToColor     *jutil.LinkedHashMap[colors.ColorId, *colors.Color]
	colorRegistry *colors.ColorRegistry
	// Non-empty for testing only. Normally empty to save on memory.
	shardViews []*ColorPool_ShardView
}

// Java: ColorPool(LinkedHashMap<ColorId, Color>, ColorRegistry, ImmutableList<ShardView>)
func newColorPool(
	idToColor *jutil.LinkedHashMap[colors.ColorId, *colors.Color],
	colorRegistry *colors.ColorRegistry,
	shardViews []*ColorPool_ShardView) *ColorPool {
	return &ColorPool{idToColor: idToColor, colorRegistry: colorRegistry, shardViews: shardViews}
}

// Java: ColorPool.getColor(ColorId)
func (p *ColorPool) GetColor(id colors.ColorId) *colors.Color {
	return p.idToColor.Get(id)
}

// Java: ColorPool.getRegistry()
func (p *ColorPool) GetRegistry() *colors.ColorRegistry {
	return p.colorRegistry
}

// Java: ColorPool.getOnlyShardForTesting()
func (p *ColorPool) GetOnlyShardForTesting() *ColorPool_ShardView {
	return colorPool_getOnlyElement(p.shardViews)
}

// colorPool_getOnlyElement is Guava's Iterables.getOnlyElement(Iterable).
//
// Go-only: language-adapter — Guava helper.
func colorPool_getOnlyElement[T any](list []T) T {
	if len(list) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	if len(list) > 1 {
		var sb strings.Builder
		sb.WriteString("expected one element but was: <")
		sb.WriteString(jutil.ToJavaString(list[0]))
		i := 1
		for ; i < 5 && i < len(list); i++ {
			sb.WriteString(", ")
			sb.WriteString(jutil.ToJavaString(list[i]))
		}
		if i < len(list) {
			sb.WriteString(", ...")
		}
		sb.WriteByte('>')
		panic(jutil.NewIllegalArgumentExceptionMsg(sb.String()))
	}
	return list[0]
}

// ColorPool_ShardView is a view of the pool based on one of the input
// shards.
type ColorPool_ShardView struct {
	trimmedOffsetToId []colors.ColorId

	// Fields only present before/while the ColorPool is being built. Null afterwards.
	typePool   *TypePool
	stringPool *StringPool

	// Set once the complete pool is built.
	colorPool *ColorPool
}

// Java: ColorPool.ShardView(TypePool, StringPool, ImmutableList<ColorId>)
func newColorPool_ShardView(
	typePool *TypePool, stringPool *StringPool, trimmedOffsetToId []colors.ColorId) *ColorPool_ShardView {
	return &ColorPool_ShardView{
		typePool:          typePool,
		stringPool:        stringPool,
		trimmedOffsetToId: trimmedOffsetToId,
	}
}

// Java: ColorPool.ShardView.getColor(int)
func (s *ColorPool_ShardView) GetColor(pointer int) *colors.Color {
	jutil.CheckStateMsg1(s.colorPool != nil, "%s", s)
	return s.colorPool.GetColor(s.GetId(pointer))
}

// Java: ColorPool.ShardView.getId(int)
func (s *ColorPool_ShardView) GetId(untrimmedOffset int) colors.ColorId {
	if TypePointers_IsAxiomatic(untrimmedOffset) {
		return TypePointers_OFFSET_TO_AXIOMATIC_COLOR[untrimmedOffset].GetId()
	} else {
		return s.trimmedOffsetToId[TypePointers_TrimOffset(untrimmedOffset)]
	}
}

// Java: ColorPool.ShardView.updateStateAfterColorPoolIsBuilt(ColorPool)
func (s *ColorPool_ShardView) UpdateStateAfterColorPoolIsBuilt(colorPool *ColorPool) {
	s.colorPool = colorPool
	s.typePool = nil
	s.stringPool = nil
}

// String is Object.toString().
func (s *ColorPool_ShardView) String() string {
	return jutil.ObjectToString("com.google.javascript.jscomp.serialization.ColorPool$ShardView", s)
}

// Java: ColorPool.builder()
func ColorPool_Builder_() *ColorPool_Builder {
	return newColorPool_Builder()
}

// Java: ColorPool.fromOnlyShardForTesting(TypePool, StringPool)
func ColorPool_FromOnlyShardForTesting(typePool *TypePool, stringPool *StringPool) *ColorPool {
	return ColorPool_Builder_().AddShardAnd(typePool, stringPool).ForTesting().Build()
}

// ColorPool_Builder collects TypePools and other data into ShardViews and
// then reconciles them into a single ColorPool.
type ColorPool_Builder struct {
	protoToShard *base.LinkedIdentityHashMap[*TypePool, *ColorPool_ShardView]
	idToColor    *jutil.LinkedHashMap[colors.ColorId, *colors.Color]
	registry     *colors.ColorRegistry_Builder
	idToProto    *jutil.HashBasedTable[colors.ColorId, *ColorPool_ShardView, *TypeProto]
	forTesting_  bool

	reconcilationDebugStack []colors.ColorId
}

// Java: ColorPool.Builder()
func newColorPool_Builder() *ColorPool_Builder {
	b := &ColorPool_Builder{
		protoToShard: base.NewLinkedIdentityHashMap[*TypePool, *ColorPool_ShardView](),
		idToColor:    jutil.NewLinkedHashMap[colors.ColorId, *colors.Color](),
		registry:     colors.ColorRegistry_Builder_(),
		idToProto:    jutil.HashBasedTable_Create[colors.ColorId, *ColorPool_ShardView, *TypeProto](),
	}
	for k, v := range colors.StandardColors_AXIOMATIC_COLORS.All() {
		b.idToColor.Put(k, v)
	}
	return b
}

// Java: ColorPool.Builder.addShardAnd(TypePool, StringPool)
func (b *ColorPool_Builder) AddShardAnd(typePool *TypePool, stringPool *StringPool) *ColorPool_Builder {
	b.AddShard(typePool, stringPool)
	return b
}

// Java: ColorPool.Builder.forTesting()
func (b *ColorPool_Builder) ForTesting() *ColorPool_Builder {
	b.forTesting_ = true
	return b
}

// Java: ColorPool.Builder.addShard(TypePool, StringPool)
func (b *ColorPool_Builder) AddShard(typePool *TypePool, stringPool *StringPool) *ColorPool_ShardView {
	jutil.CheckStateMsg(b.idToProto.IsEmpty(), "build has already been called")

	existing := b.protoToShard.Get(typePool)
	if existing != nil {
		jutil.CheckStateMsg1(base.JSCompObjects_Identical(typePool, TypePool_GetDefaultInstance()), "%s", typePool)
		return existing
	}

	trimmedOffsetToId := colorPool_createTrimmedOffsetToId(typePool)
	shard := newColorPool_ShardView(typePool, stringPool, trimmedOffsetToId)
	b.protoToShard.Put(typePool, shard)

	if typePool.HasDebugInfo() {
		for _, m := range typePool.GetDebugInfo().GetMismatchList() {
			for _, pointer := range m.GetInvolvedColorList() {
				b.registry.AddMismatchLocation(shard.GetId(int(pointer)), m.GetSourceRef())
			}
		}
	}

	return shard
}

// Java: ColorPool.Builder.build()
func (b *ColorPool_Builder) Build() *ColorPool {
	jutil.CheckStateMsg(b.idToProto.IsEmpty(), "build has already been called")

	b.ForEachShard(
		func(shard *ColorPool_ShardView) {
			for i := 0; i < shard.typePool.GetTypeCount(); i++ {
				id := shard.trimmedOffsetToId[i]
				b.idToProto.Put(id, shard, shard.typePool.GetType(i))
			}
		})

	for id := range colors.StandardColors_AXIOMATIC_COLORS.Keys() {
		MalformedTypedAstException_CheckWellFormedWithParam(
			!b.idToProto.ContainsRow(id), "Found serialized definiton for axiomatic color", id)
	}

	for id := range b.idToProto.RowKeySet() {
		b.LookupOrReconcileColor(id)
	}

	for colorId := range colors.ColorRegistry_REQUIRED_IDS.All() {
		b.registry.SetNativeColor(
			b.idToColor.ComputeIfAbsent(
				colorId, func(unused colors.ColorId) *colors.Color { return colors.Color_SingleBuilder().SetId(colorId).Build() }))
	}

	b.ForEachShard(
		func(shard *ColorPool_ShardView) {
			for _, edge := range shard.typePool.GetDisambiguationEdgesList() {
				b.registry.AddDisambiguationEdge(
					b.idToColor.Get(shard.GetId(colorPool_validatePointer(int(edge.GetSubtype()), shard))), //
					b.idToColor.Get(shard.GetId(colorPool_validatePointer(int(edge.GetSupertype()), shard))))
			}
		})

	shardViewsForTesting := []*ColorPool_ShardView{}
	if b.forTesting_ {
		b.ForEachShard(func(shard *ColorPool_ShardView) { shardViewsForTesting = append(shardViewsForTesting, shard) })
	}

	colorPool := newColorPool(b.idToColor, b.registry.Build(), shardViewsForTesting)

	b.ForEachShard(func(shard *ColorPool_ShardView) { shard.UpdateStateAfterColorPoolIsBuilt(colorPool) })
	return colorPool
}

// Java: ColorPool.Builder.forEachShard(Consumer<ShardView>)
func (b *ColorPool_Builder) ForEachShard(fn func(*ColorPool_ShardView)) {
	b.protoToShard.ForEach(func(unused *TypePool, shard *ColorPool_ShardView) { fn(shard) })
}

// Java: ColorPool.Builder.lookupOrReconcileColor(ColorId)
func (b *ColorPool_Builder) LookupOrReconcileColor(id colors.ColorId) *colors.Color {
	b.reconcilationDebugStack = append(b.reconcilationDebugStack, id)
	defer func() {
		b.reconcilationDebugStack = b.reconcilationDebugStack[:len(b.reconcilationDebugStack)-1]
	}()
	existing := b.idToColor.PutIfAbsent(id, colorPool_PENDING_COLOR)
	if existing != nil {
		if base.JSCompObjects_Identical(existing, colorPool_PENDING_COLOR) {
			rows := make([]string, 0, len(b.reconcilationDebugStack))
			for _, stackId := range b.reconcilationDebugStack {
				rows = append(rows, b.idToProto.Row(stackId).String())
			}
			panic(NewMalformedTypedAstException(
				"Cyclic Color structure detected: [" + strings.Join(rows, ", ") + "]"))
		}
		return existing
	}

	viewToProto := b.idToProto.Row(id)
	var sample *TypeProto
	for v := range viewToProto.Values() {
		sample = v
		break
	}
	jutil.CheckNotNullMsg1(sample, "%s", id)

	var result *colors.Color
	switch sample.GetKindCase() {
	case TypeProto_KindCase_OBJECT:
		result = b.ReconcileObjectProtos(id, viewToProto)
	case TypeProto_KindCase_UNION:
		result = b.ReconcileUnionProtos(id, viewToProto)
	default:
		panic(jutil.NewAssertionError(sample.String()))
	}

	jutil.CheckStateMsg1(result != nil, "%s", id)
	b.idToColor.Put(id, result)
	return result
}

// Java: ColorPool.Builder.reconcileObjectProtos(ColorId, Map<ShardView, TypeProto>)
func (b *ColorPool_Builder) ReconcileObjectProtos(id colors.ColorId, viewToProto *jutil.LinkedHashMap[*ColorPool_ShardView, *TypeProto]) *colors.Color {
	instanceColors := jutil.NewObjectLinkedHashSet[*colors.Color]()
	prototypes := jutil.NewObjectLinkedHashSet[*colors.Color]()
	ownProperties := jutil.NewLinkedHashSet[string]()
	isClosureAssert := base.Tri_UNKNOWN
	isConstructor := false
	isInvalidating := false
	propertiesKeepOriginalName := false

	for shard, proto := range viewToProto.All() {
		jutil.CheckState(proto.HasObject())
		objProto := proto.GetObject()

		instanceTypeIds := objProto.GetInstanceTypeList()
		for i := 0; i < len(instanceTypeIds); i++ {
			instanceTypePointer := instanceTypeIds[i]
			instanceColors.Add(b.LookupOrReconcileColor(shard.GetId(int(instanceTypePointer))))
		}

		isClosureAssertBool := objProto.GetClosureAssert()
		MalformedTypedAstException_CheckWellFormedWithParam(
			isClosureAssert.ToBoolean(isClosureAssertBool) == isClosureAssertBool,
			"Inconsistent values for closure_assert",
			objProto)
		isClosureAssert = base.Tri_ForBoolean(isClosureAssertBool)

		isConstructor = isConstructor || objProto.GetMarkedConstructor()
		isInvalidating = isInvalidating || objProto.GetIsInvalidating()
		propertiesKeepOriginalName = propertiesKeepOriginalName || objProto.GetPropertiesKeepOriginalName()

		prototypeIds := objProto.GetPrototypeList()
		for i := 0; i < len(prototypeIds); i++ {
			prototypePointer := prototypeIds[i]
			prototypes.Add(b.LookupOrReconcileColor(shard.GetId(int(prototypePointer))))
		}
		for i := 0; i < objProto.GetOwnPropertyCount(); i++ {
			ownProperties.Add(shard.stringPool.Get(int(objProto.GetOwnProperty(i))))
		}
	}

	return colors.Color_SingleBuilder().
		SetId(id).
		SetInstanceColors(instanceColors).
		SetPrototypes(prototypes).
		SetOwnProperties(ownProperties).
		SetClosureAssert(isClosureAssert.ToBoolean(false)).
		SetConstructor(isConstructor).
		SetInvalidating(isInvalidating).
		SetPropertiesKeepOriginalName(propertiesKeepOriginalName).
		Build()
}

// Java: ColorPool.Builder.reconcileUnionProtos(ColorId, Map<ShardView, TypeProto>)
func (b *ColorPool_Builder) ReconcileUnionProtos(id colors.ColorId, viewToProto *jutil.LinkedHashMap[*ColorPool_ShardView, *TypeProto]) *colors.Color {
	union := jutil.NewObjectLinkedHashSet[*colors.Color]()
	for shard, proto := range viewToProto.All() {
		jutil.CheckStateMsg1(proto.HasUnion(), "%s", proto)
		unionMemberList := proto.GetUnion().GetUnionMemberList()
		for i := 0; i < len(unionMemberList); i++ {
			memberPointer := unionMemberList[i]
			memberId := shard.GetId(int(memberPointer))
			member := b.LookupOrReconcileColor(memberId)
			MalformedTypedAstException_CheckWellFormedWithParam(!member.IsUnion(), "Reconciling union with non-union", proto)
			union.Add(member)
		}
	}

	result := colors.Color_CreateUnion(union)
	jutil.CheckStateMsg2(id.Equals(result.GetId()), "%s == %s", id, result)
	return result
}

// Java: ColorPool.createTrimmedOffsetToId(TypePool)
func colorPool_createTrimmedOffsetToId(typePool *TypePool) []colors.ColorId {
	ids := make([]*colors.ColorId, typePool.GetTypeCount())

	for i := 0; i < len(ids); i++ {
		proto := typePool.GetType(i)
		switch proto.GetKindCase() {
		case TypeProto_KindCase_OBJECT:
			id := colors.ColorId_FromBytes(proto.GetObject().GetUuid())
			ids[i] = &id
		case TypeProto_KindCase_UNION:
			// Defer generating union IDs until we have the IDs of all the element types.
		default:
			panic(NewMalformedTypedAstException(proto))
		}
	}

	for i := 0; i < len(ids); i++ {
		proto := typePool.GetType(i)
		switch proto.GetKindCase() {
		case TypeProto_KindCase_OBJECT:
		case TypeProto_KindCase_UNION:
			id := colorPool_createUnionColorId(proto, ids)
			ids[i] = &id
		default:
			panic(jutil.NewAssertionError(proto.String()))
		}
	}

	seenIds := jutil.NewLinkedHashSet[colors.ColorId]()
	for i := 0; i < len(ids); i++ {
		proto := typePool.GetType(i)
		MalformedTypedAstException_CheckWellFormedWithParam(seenIds.Add(*ids[i]), "Duplicate ID in single shard", proto)
	}

	out := make([]colors.ColorId, len(ids))
	for i, id := range ids {
		out[i] = *id
	}
	return out
}

// Java: ColorPool.createUnionColorId(TypeProto, ColorId[])
func colorPool_createUnionColorId(proto *TypeProto, allObjectIds []*colors.ColorId) colors.ColorId {
	MalformedTypedAstException_CheckWellFormedWithParam(proto.GetUnion().GetUnionMemberCount() > 1, "Union has too few members", proto)
	members := jutil.NewLinkedHashSet[colors.ColorId]()
	unionMemberList := proto.GetUnion().GetUnionMemberList()
	for i := 0; i < len(unionMemberList); i++ {
		memberPointer := int(unionMemberList[i])
		var memberId *colors.ColorId
		if TypePointers_IsAxiomatic(memberPointer) {
			id := TypePointers_OFFSET_TO_AXIOMATIC_COLOR[memberPointer].GetId()
			memberId = &id
		} else {
			memberId = allObjectIds[TypePointers_TrimOffset(memberPointer)]
		}
		MalformedTypedAstException_CheckWellFormedWithParam(memberId != nil, "Union member not found", proto)
		members.Add(*memberId)
	}
	return colors.ColorId_Union(members)
}

// Java: ColorPool.validatePointer(int, ShardView)
func colorPool_validatePointer(offset int, shard *ColorPool_ShardView) int {
	MalformedTypedAstException_CheckWellFormedWithParam(
		0 <= offset && offset < TypePointers_UntrimOffset(len(shard.trimmedOffsetToId)),
		"Pointer offset outside of shard",
		offset)
	return offset
}

// 0xDEADBEEF as a Java int.
var colorPool_PENDING_COLOR = colors.Color_SingleBuilder().SetId(colors.ColorId_FromUnsignedInt(-0x21524111)).Build()
