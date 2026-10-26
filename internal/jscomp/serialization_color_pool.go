// Port of com/google/javascript/jscomp/serialization/ColorPool.java @ closure-compiler v20261021

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
	// caused OOMs in large projects. So we just use a map and hide the mutability behind
	// accessor methods.
	idToColor     map[colors.ColorId]*colors.Color
	colorRegistry *colors.ColorRegistry
	// Non-empty for testing only. Normally empty to save on memory.
	shardViews []*ColorPool_ShardView
}

var (
	colorPool_emptyColorSet  = jutil.NewObjectLinkedHashSet[*colors.Color]()
	colorPool_emptyStringSet = jutil.NewLinkedHashSet[string]()
)

// Java: ColorPool(LinkedHashMap<ColorId, Color>, ColorRegistry, ImmutableList<ShardView>)
func newColorPool(
	idToColor map[colors.ColorId]*colors.Color,
	colorRegistry *colors.ColorRegistry,
	shardViews []*ColorPool_ShardView) *ColorPool {
	return &ColorPool{idToColor: idToColor, colorRegistry: colorRegistry, shardViews: shardViews}
}

// Java: ColorPool.getColor(ColorId)
func (p *ColorPool) GetColor(id colors.ColorId) *colors.Color {
	return p.idToColor[id]
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
	trimmedOffsetToId    []colors.ColorId
	trimmedOffsetToColor []*colors.Color
	offsetToColor        []*colors.Color

	// Fields only present before/while the ColorPool is being built. Null afterwards.
	typePool   *TypePool
	stringPool *StringPool

	// Set once the complete pool is built.
	built bool
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
	if uint(pointer) < uint(len(s.offsetToColor)) {
		return s.offsetToColor[pointer]
	}
	jutil.CheckStateMsg1(s.built, "%s", s)
	return s.offsetToColor[pointer]
}

// Java: ColorPool.ShardView.getId(int)
func (s *ColorPool_ShardView) GetId(untrimmedOffset int) colors.ColorId {
	if s.offsetToColor != nil {
		return s.offsetToColor[untrimmedOffset].GetId()
	}
	if TypePointers_IsAxiomatic(untrimmedOffset) {
		return TypePointers_OFFSET_TO_AXIOMATIC_COLOR[untrimmedOffset].GetId()
	} else {
		return s.trimmedOffsetToId[TypePointers_TrimOffset(untrimmedOffset)]
	}
}

// Java: ColorPool.ShardView.updateStateAfterColorPoolIsBuilt(ColorPool)
func (s *ColorPool_ShardView) UpdateStateAfterColorPoolIsBuilt(colorPool *ColorPool) {
	s.built = true
	s.typePool = nil
	s.stringPool = nil
	axiomaticCount := len(TypePointers_OFFSET_TO_AXIOMATIC_COLOR)
	offsetToColor := make([]*colors.Color, axiomaticCount+len(s.trimmedOffsetToId))
	copy(offsetToColor, TypePointers_OFFSET_TO_AXIOMATIC_COLOR)
	if s.trimmedOffsetToColor != nil {
		copy(offsetToColor[axiomaticCount:], s.trimmedOffsetToColor)
		s.trimmedOffsetToColor = nil
	} else {
		idToColor := colorPool.idToColor
		for i, id := range s.trimmedOffsetToId {
			offsetToColor[axiomaticCount+i] = idToColor[id]
		}
	}
	s.trimmedOffsetToId = nil
	s.offsetToColor = offsetToColor
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

type colorPoolShardProto struct {
	shard *ColorPool_ShardView
	proto *TypeProto
}

type colorPoolIdEntry struct {
	firstShard *ColorPool_ShardView
	firstProto *TypeProto
	extra      []colorPoolShardProto
}

func (e colorPoolIdEntry) rowString() string {
	m := jutil.NewLinkedHashMap[*ColorPool_ShardView, *TypeProto]()
	m.Put(e.firstShard, e.firstProto)
	for _, sp := range e.extra {
		m.Put(sp.shard, sp.proto)
	}
	return m.String()
}

// ColorPool_Builder collects TypePools and other data into ShardViews and
// then reconciles them into a single ColorPool.
type ColorPool_Builder struct {
	protoToShard map[*TypePool]*ColorPool_ShardView
	shards       []*ColorPool_ShardView
	idToColor    map[colors.ColorId]*colors.Color
	registry     *colors.ColorRegistry_Builder
	idToProto    *jutil.LinkedHashMap[colors.ColorId, colorPoolIdEntry]
	forTesting_  bool
	built        bool

	colorBuilder     colors.Color_Builder
	colorSlab        jutil.Slab[colors.Color]
	stringSetSlab    jutil.Slab[jutil.LinkedHashSet[string]]
	stringSlotSlab   []jutil.LhmSlot[string, struct{}]
	colorSetSlab     jutil.Slab[jutil.ObjectLinkedHashSet[*colors.Color]]
	colorSlotSlab    []jutil.OlhmSlot[*colors.Color, struct{}]
	allocColorSetFn  func(capHint int) *jutil.ObjectLinkedHashSet[*colors.Color]
	singlePropSets   map[int32]*jutil.LinkedHashSet[string]
	singlePropByName map[string]*jutil.LinkedHashSet[string]
	propScratch      []string
	seenProps        map[string]uint32
	seenPropsGen     uint32

	reconcilationDebugStack []colors.ColorId
}

func (b *ColorPool_Builder) allocStringSet() *jutil.LinkedHashSet[string] {
	return b.stringSetSlab.New()
}

func (b *ColorPool_Builder) allocStringSlots(n int) []jutil.LhmSlot[string, struct{}] {
	const slotSlabSize = 1024 // 24 KiB <= 32 KiB
	if n > slotSlabSize {
		return make([]jutil.LhmSlot[string, struct{}], n)
	}
	if len(b.stringSlotSlab) < n {
		b.stringSlotSlab = make([]jutil.LhmSlot[string, struct{}], slotSlabSize)
	}
	buf := b.stringSlotSlab[:n:n]
	b.stringSlotSlab = b.stringSlotSlab[n:]
	return buf
}

func (b *ColorPool_Builder) internSinglePropSet(prop string) *jutil.LinkedHashSet[string] {
	if b.singlePropByName == nil {
		b.singlePropByName = make(map[string]*jutil.LinkedHashSet[string], 256)
	} else if cached := b.singlePropByName[prop]; cached != nil {
		return cached
	}
	s := b.allocStringSet()
	s.InitWithBuffer(b.allocStringSlots(1))
	s.AddNew(prop)
	b.singlePropByName[prop] = s
	return s
}

func (b *ColorPool_Builder) buildUniqueStringSet(props []string) *jutil.LinkedHashSet[string] {
	if len(props) == 0 {
		return colorPool_emptyStringSet
	}
	if len(props) == 1 {
		return b.internSinglePropSet(props[0])
	}
	s := b.allocStringSet()
	s.InitWithBuffer(b.allocStringSlots(len(props)))
	for _, p := range props {
		s.AddNew(p)
	}
	return s
}

func (b *ColorPool_Builder) appendDedupProp(prop string, useMap bool) {
	if useMap {
		if b.seenProps[prop] == b.seenPropsGen {
			return
		}
		b.seenProps[prop] = b.seenPropsGen
		b.propScratch = append(b.propScratch, prop)
		return
	}
	for _, existing := range b.propScratch {
		if existing == prop {
			return
		}
	}
	b.propScratch = append(b.propScratch, prop)
}

func (b *ColorPool_Builder) allocColorSlots(n int) []jutil.OlhmSlot[*colors.Color, struct{}] {
	const slotSlabSize = 1024 // 16 KiB <= 32 KiB
	if n > slotSlabSize {
		return make([]jutil.OlhmSlot[*colors.Color, struct{}], n)
	}
	if len(b.colorSlotSlab) < n {
		b.colorSlotSlab = make([]jutil.OlhmSlot[*colors.Color, struct{}], slotSlabSize)
	}
	buf := b.colorSlotSlab[:n:n]
	b.colorSlotSlab = b.colorSlotSlab[n:]
	return buf
}

func (b *ColorPool_Builder) allocColorSetWithCap(capHint int) *jutil.ObjectLinkedHashSet[*colors.Color] {
	s := b.colorSetSlab.New()
	if capHint < 1 {
		capHint = 2
	}
	s.InitWithBuffer(b.allocColorSlots(capHint))
	return s
}

func (b *ColorPool_Builder) allocColorSet() *jutil.ObjectLinkedHashSet[*colors.Color] {
	return b.allocColorSetWithCap(2)
}

// Java: ColorPool.Builder()
func newColorPool_Builder() *ColorPool_Builder {
	b := &ColorPool_Builder{
		protoToShard: make(map[*TypePool]*ColorPool_ShardView, 4),
		registry:     colors.ColorRegistry_Builder_(),
		colorSlab:    jutil.NewFixedSlab[colors.Color](256), // 18.4 KiB <= 32 KiB
	}
	b.colorBuilder.SetSlab(&b.colorSlab)
	b.allocColorSetFn = b.allocColorSetWithCap
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
	jutil.CheckStateMsg(!b.built, "build has already been called")

	existing := b.protoToShard[typePool]
	if existing != nil {
		jutil.CheckStateMsg1(base.JSCompObjects_Identical(typePool, TypePool_GetDefaultInstance()), "%s", typePool)
		return existing
	}

	trimmedOffsetToId := colorPool_createTrimmedOffsetToId(typePool)
	shard := newColorPool_ShardView(typePool, stringPool, trimmedOffsetToId)
	b.protoToShard[typePool] = shard
	b.shards = append(b.shards, shard)

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
	jutil.CheckStateMsg(!b.built, "build has already been called")
	b.built = true

	totalTypes := 0
	shardCount := 0
	var singleShard *ColorPool_ShardView
	b.ForEachShard(func(shard *ColorPool_ShardView) {
		totalTypes += shard.typePool.GetTypeCount()
		shardCount++
		singleShard = shard
	})

	b.idToColor = make(map[colors.ColorId]*colors.Color, totalTypes+16)

	if shardCount <= 1 {
		if singleShard != nil && len(singleShard.trimmedOffsetToId) > 0 {
			typeList := singleShard.typePool.GetTypeList()
			for i, id := range singleShard.trimmedOffsetToId {
				if b.idToColor[id] != nil {
					MalformedTypedAstException_CheckWellFormedWithParam(false, "Duplicate ID in single shard", typeList[i])
				}
				b.idToColor[id] = colorPool_UNRECONCILED_COLOR
			}
			for id, v := range colors.StandardColors_AXIOMATIC_COLORS.All() {
				if b.idToColor[id] != nil {
					MalformedTypedAstException_CheckWellFormedWithParam(
						false, "Found serialized definiton for axiomatic color", id)
				}
				b.idToColor[id] = v
			}
			singleShard.trimmedOffsetToColor = make([]*colors.Color, len(typeList))
			b.singlePropSets = make(map[int32]*jutil.LinkedHashSet[string], 256)
			for i := range typeList {
				b.reconcileSingleShardAt(singleShard, i)
			}
			b.singlePropSets = nil
		} else {
			for k, v := range colors.StandardColors_AXIOMATIC_COLORS.All() {
				b.idToColor[k] = v
			}
		}
	} else {
		for k, v := range colors.StandardColors_AXIOMATIC_COLORS.All() {
			b.idToColor[k] = v
		}
		b.idToProto = jutil.NewLinkedHashMapWithCapacity[colors.ColorId, colorPoolIdEntry](totalTypes)
		b.ForEachShard(
			func(shard *ColorPool_ShardView) {
				typeList := shard.typePool.GetTypeList()
				seenIds := make(map[colors.ColorId]struct{}, len(typeList))
				for i, proto := range typeList {
					id := shard.trimmedOffsetToId[i]
					if _, dup := seenIds[id]; dup {
						MalformedTypedAstException_CheckWellFormedWithParam(false, "Duplicate ID in single shard", proto)
					}
					seenIds[id] = struct{}{}
					if entry, ok := b.idToProto.Lookup(id); !ok {
						b.idToProto.PutNew(id, colorPoolIdEntry{
							firstShard: shard,
							firstProto: proto,
						})
					} else {
						if entry.firstShard == shard {
							entry.firstProto = proto
						} else {
							replaced := false
							for j := range entry.extra {
								if entry.extra[j].shard == shard {
									entry.extra[j].proto = proto
									replaced = true
									break
								}
							}
							if !replaced {
								entry.extra = append(entry.extra, colorPoolShardProto{shard: shard, proto: proto})
							}
						}
						b.idToProto.Put(id, entry)
					}
				}
			})

		for id := range colors.StandardColors_AXIOMATIC_COLORS.Keys() {
			if b.idToProto.ContainsKey(id) {
				MalformedTypedAstException_CheckWellFormedWithParam(
					false, "Found serialized definiton for axiomatic color", id)
			}
		}

		for id := range b.idToProto.Keys() {
			b.LookupOrReconcileColor(id)
		}
	}

	for colorId := range colors.ColorRegistry_REQUIRED_IDS.All() {
		c := b.idToColor[colorId]
		if c == nil {
			c = b.colorBuilder.
				SetId(colorId).
				SetInstanceColors(colorPool_emptyColorSet).
				SetPrototypes(colorPool_emptyColorSet).
				SetOwnProperties(colorPool_emptyStringSet).
				SetUnionElements(colorPool_emptyColorSet).
				SetClosureAssert(false).
				SetConstructor(false).
				SetInvalidating(false).
				SetPropertiesKeepOriginalName(false).
				Build()
			b.idToColor[colorId] = c
		}
		b.registry.SetNativeColor(c)
	}

	if shardCount == 1 && singleShard != nil {
		for _, edge := range singleShard.typePool.GetDisambiguationEdgesList() {
			subOffset := colorPool_validatePointer(int(edge.GetSubtype()), singleShard)
			supOffset := colorPool_validatePointer(int(edge.GetSupertype()), singleShard)
			var subColor, supColor *colors.Color
			if TypePointers_IsAxiomatic(subOffset) {
				subColor = TypePointers_OFFSET_TO_AXIOMATIC_COLOR[subOffset]
			} else {
				subColor = singleShard.trimmedOffsetToColor[TypePointers_TrimOffset(subOffset)]
			}
			if TypePointers_IsAxiomatic(supOffset) {
				supColor = TypePointers_OFFSET_TO_AXIOMATIC_COLOR[supOffset]
			} else {
				supColor = singleShard.trimmedOffsetToColor[TypePointers_TrimOffset(supOffset)]
			}
			b.registry.AddDisambiguationEdge(subColor, supColor)
		}
	} else {
		b.ForEachShard(
			func(shard *ColorPool_ShardView) {
				for _, edge := range shard.typePool.GetDisambiguationEdgesList() {
					b.registry.AddDisambiguationEdge(
						b.idToColor[shard.GetId(colorPool_validatePointer(int(edge.GetSubtype()), shard))], //
						b.idToColor[shard.GetId(colorPool_validatePointer(int(edge.GetSupertype()), shard))])
				}
			})
	}

	shardViewsForTesting := []*ColorPool_ShardView{}
	if b.forTesting_ {
		b.ForEachShard(func(shard *ColorPool_ShardView) { shardViewsForTesting = append(shardViewsForTesting, shard) })
	}

	colorPool := newColorPool(b.idToColor, b.registry.Build(), shardViewsForTesting)

	b.ForEachShard(func(shard *ColorPool_ShardView) { shard.UpdateStateAfterColorPoolIsBuilt(colorPool) })
	b.idToProto = nil
	b.singlePropSets = nil
	b.singlePropByName = nil
	b.propScratch = nil
	b.seenProps = nil
	return colorPool
}

func (b *ColorPool_Builder) resolveSingleShardPointer(shard *ColorPool_ShardView, untrimmedOffset int) *colors.Color {
	if TypePointers_IsAxiomatic(untrimmedOffset) {
		return TypePointers_OFFSET_TO_AXIOMATIC_COLOR[untrimmedOffset]
	}
	return b.reconcileSingleShardAt(shard, TypePointers_TrimOffset(untrimmedOffset))
}

func (b *ColorPool_Builder) reconcileSingleShardAt(shard *ColorPool_ShardView, trimmedIdx int) *colors.Color {
	id := shard.trimmedOffsetToId[trimmedIdx]
	b.reconcilationDebugStack = append(b.reconcilationDebugStack, id)
	defer func() {
		b.reconcilationDebugStack = b.reconcilationDebugStack[:len(b.reconcilationDebugStack)-1]
	}()
	existing := shard.trimmedOffsetToColor[trimmedIdx]
	if existing != nil {
		if base.JSCompObjects_Identical(existing, colorPool_PENDING_COLOR) {
			typeList := shard.typePool.GetTypeList()
			rows := make([]string, 0, len(b.reconcilationDebugStack))
			for _, stackId := range b.reconcilationDebugStack {
				for i, sid := range shard.trimmedOffsetToId {
					if sid == stackId {
						rows = append(rows, (colorPoolIdEntry{firstShard: shard, firstProto: typeList[i]}).rowString())
						break
					}
				}
			}
			panic(NewMalformedTypedAstException(
				"Cyclic Color structure detected: [" + strings.Join(rows, ", ") + "]"))
		}
		return existing
	}
	shard.trimmedOffsetToColor[trimmedIdx] = colorPool_PENDING_COLOR

	proto := shard.typePool.GetTypeList()[trimmedIdx]
	var result *colors.Color
	if proto.HasObject() {
		objProto := proto.GetObject()
		instanceColors := colorPool_emptyColorSet
		if instanceTypeIds := objProto.GetInstanceTypeList(); len(instanceTypeIds) > 0 {
			instanceColors = b.allocColorSetWithCap(len(instanceTypeIds))
			for i := 0; i < len(instanceTypeIds); i++ {
				instanceColors.Add(b.resolveSingleShardPointer(shard, int(instanceTypeIds[i])))
			}
		}
		prototypes := colorPool_emptyColorSet
		if prototypeIds := objProto.GetPrototypeList(); len(prototypeIds) > 0 {
			prototypes = b.allocColorSetWithCap(len(prototypeIds))
			for i := 0; i < len(prototypeIds); i++ {
				prototypes.Add(b.resolveSingleShardPointer(shard, int(prototypeIds[i])))
			}
		}
		ownProperties := colorPool_emptyStringSet
		if ownPropList := objProto.GetOwnPropertyList(); len(ownPropList) == 1 {
			propIdx := ownPropList[0]
			if cached := b.singlePropSets[propIdx]; cached != nil {
				ownProperties = cached
			} else {
				ownProperties = b.allocStringSet()
				ownProperties.InitWithBuffer(b.allocStringSlots(1))
				ownProperties.AddNew(shard.stringPool.Get(int(propIdx)))
				b.singlePropSets[propIdx] = ownProperties
			}
		} else if len(ownPropList) > 1 {
			b.propScratch = b.propScratch[:0]
			useMap := len(ownPropList) > 8
			if useMap {
				b.seenPropsGen++
				if b.seenProps == nil {
					b.seenProps = make(map[string]uint32, 64)
				} else if b.seenPropsGen == 0 {
					clear(b.seenProps)
					b.seenPropsGen = 1
				}
			}
			for i := 0; i < len(ownPropList); i++ {
				b.appendDedupProp(shard.stringPool.Get(int(ownPropList[i])), useMap)
			}
			ownProperties = b.buildUniqueStringSet(b.propScratch)
		}
		result = b.colorBuilder.
			SetId(id).
			SetInstanceColors(instanceColors).
			SetPrototypes(prototypes).
			SetOwnProperties(ownProperties).
			SetUnionElements(colorPool_emptyColorSet).
			SetClosureAssert(objProto.GetClosureAssert()).
			SetConstructor(objProto.GetMarkedConstructor()).
			SetInvalidating(objProto.GetIsInvalidating()).
			SetPropertiesKeepOriginalName(objProto.GetPropertiesKeepOriginalName()).
			Build()
	} else if proto.HasUnion() {
		unionMemberList := proto.GetUnion().GetUnionMemberList()
		var unionBuf [8]jutil.OlhmSlot[*colors.Color, struct{}]
		var union jutil.ObjectLinkedHashSet[*colors.Color]
		union.InitWithBuffer(unionBuf[:])
		for i := 0; i < len(unionMemberList); i++ {
			member := b.resolveSingleShardPointer(shard, int(unionMemberList[i]))
			MalformedTypedAstException_CheckWellFormedWithParam(!member.IsUnion(), "Reconciling union with non-union", proto)
			union.Add(member)
		}
		result = colors.Color_CreateUnionWithBuilder(&union, &b.colorBuilder, b.allocColorSetFn)
		if !id.Equals(result.GetId()) {
			jutil.CheckStateMsg2(false, "%s == %s", id, result)
		}
	} else {
		panic(jutil.NewAssertionError(proto.String()))
	}

	if result == nil {
		jutil.CheckStateMsg1(false, "%s", id)
	}
	shard.trimmedOffsetToColor[trimmedIdx] = result
	b.idToColor[id] = result
	return result
}

// Java: ColorPool.Builder.forEachShard(Consumer<ShardView>)
func (b *ColorPool_Builder) ForEachShard(fn func(*ColorPool_ShardView)) {
	for _, shard := range b.shards {
		fn(shard)
	}
}

// Java: ColorPool.Builder.lookupOrReconcileColor(ColorId)
func (b *ColorPool_Builder) LookupOrReconcileColor(id colors.ColorId) *colors.Color {
	b.reconcilationDebugStack = append(b.reconcilationDebugStack, id)
	defer func() {
		b.reconcilationDebugStack = b.reconcilationDebugStack[:len(b.reconcilationDebugStack)-1]
	}()
	existing := b.idToColor[id]
	if existing != nil {
		if base.JSCompObjects_Identical(existing, colorPool_PENDING_COLOR) {
			rows := make([]string, 0, len(b.reconcilationDebugStack))
			for _, stackId := range b.reconcilationDebugStack {
				rows = append(rows, b.idToProto.Get(stackId).rowString())
			}
			panic(NewMalformedTypedAstException(
				"Cyclic Color structure detected: [" + strings.Join(rows, ", ") + "]"))
		}
		return existing
	}
	b.idToColor[id] = colorPool_PENDING_COLOR

	entry := b.idToProto.Get(id)
	sample := entry.firstProto
	if sample == nil {
		jutil.CheckNotNullMsg1(sample, "%s", id)
	}

	var result *colors.Color
	if sample.HasObject() {
		result = b.reconcileObjectProtos(id, entry)
	} else if sample.HasUnion() {
		result = b.reconcileUnionProtos(id, entry)
	} else {
		panic(jutil.NewAssertionError(sample.String()))
	}

	if result == nil {
		jutil.CheckStateMsg1(false, "%s", id)
	}
	b.idToColor[id] = result
	return result
}

func (b *ColorPool_Builder) reconcileObjectProtos(id colors.ColorId, entry colorPoolIdEntry) *colors.Color {
	var instanceColors *jutil.ObjectLinkedHashSet[*colors.Color]
	var prototypes *jutil.ObjectLinkedHashSet[*colors.Color]
	var ownProperties *jutil.LinkedHashSet[string]
	isClosureAssert := base.Tri_UNKNOWN
	isConstructor := false
	isInvalidating := false
	propertiesKeepOriginalName := false
	totalOwnProps := 0

	for idx := -1; idx < len(entry.extra); idx++ {
		var shard *ColorPool_ShardView
		var proto *TypeProto
		if idx < 0 {
			shard = entry.firstShard
			proto = entry.firstProto
		} else {
			shard = entry.extra[idx].shard
			proto = entry.extra[idx].proto
		}
		jutil.CheckState(proto.HasObject())
		objProto := proto.GetObject()

		instanceTypeIds := objProto.GetInstanceTypeList()
		if len(instanceTypeIds) > 0 {
			if instanceColors == nil {
				instanceColors = b.allocColorSet()
			}
			for i := 0; i < len(instanceTypeIds); i++ {
				instanceTypePointer := instanceTypeIds[i]
				instanceColors.Add(b.LookupOrReconcileColor(shard.GetId(int(instanceTypePointer))))
			}
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
		if len(prototypeIds) > 0 {
			if prototypes == nil {
				prototypes = b.allocColorSet()
			}
			for i := 0; i < len(prototypeIds); i++ {
				prototypePointer := prototypeIds[i]
				prototypes.Add(b.LookupOrReconcileColor(shard.GetId(int(prototypePointer))))
			}
		}
		totalOwnProps += objProto.GetOwnPropertyCount()
	}

	if totalOwnProps == 1 {
		for idx := -1; idx < len(entry.extra); idx++ {
			shard := entry.firstShard
			proto := entry.firstProto
			if idx >= 0 {
				shard = entry.extra[idx].shard
				proto = entry.extra[idx].proto
			}
			objProto := proto.GetObject()
			if objProto.GetOwnPropertyCount() == 1 {
				ownProperties = b.internSinglePropSet(shard.stringPool.Get(int(objProto.GetOwnProperty(0))))
				break
			}
		}
	} else if totalOwnProps > 1 {
		b.propScratch = b.propScratch[:0]
		useMap := totalOwnProps > 8
		if useMap {
			b.seenPropsGen++
			if b.seenProps == nil {
				b.seenProps = make(map[string]uint32, 64)
			} else if b.seenPropsGen == 0 {
				clear(b.seenProps)
				b.seenPropsGen = 1
			}
		}
		for idx := -1; idx < len(entry.extra); idx++ {
			shard := entry.firstShard
			proto := entry.firstProto
			if idx >= 0 {
				shard = entry.extra[idx].shard
				proto = entry.extra[idx].proto
			}
			objProto := proto.GetObject()
			ownPropCount := objProto.GetOwnPropertyCount()
			for i := 0; i < ownPropCount; i++ {
				b.appendDedupProp(shard.stringPool.Get(int(objProto.GetOwnProperty(i))), useMap)
			}
		}
		ownProperties = b.buildUniqueStringSet(b.propScratch)
	}

	if instanceColors == nil {
		instanceColors = colorPool_emptyColorSet
	}
	if prototypes == nil {
		prototypes = colorPool_emptyColorSet
	}
	if ownProperties == nil {
		ownProperties = colorPool_emptyStringSet
	}

	return b.colorBuilder.
		SetId(id).
		SetInstanceColors(instanceColors).
		SetPrototypes(prototypes).
		SetOwnProperties(ownProperties).
		SetUnionElements(colorPool_emptyColorSet).
		SetClosureAssert(isClosureAssert.ToBoolean(false)).
		SetConstructor(isConstructor).
		SetInvalidating(isInvalidating).
		SetPropertiesKeepOriginalName(propertiesKeepOriginalName).
		Build()
}

func (b *ColorPool_Builder) reconcileUnionProtos(id colors.ColorId, entry colorPoolIdEntry) *colors.Color {
	var unionBuf [8]jutil.OlhmSlot[*colors.Color, struct{}]
	var union jutil.ObjectLinkedHashSet[*colors.Color]
	union.InitWithBuffer(unionBuf[:])
	for idx := -1; idx < len(entry.extra); idx++ {
		var shard *ColorPool_ShardView
		var proto *TypeProto
		if idx < 0 {
			shard = entry.firstShard
			proto = entry.firstProto
		} else {
			shard = entry.extra[idx].shard
			proto = entry.extra[idx].proto
		}
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

	result := colors.Color_CreateUnionWithBuilder(&union, &b.colorBuilder, b.allocColorSetFn)
	if !id.Equals(result.GetId()) {
		jutil.CheckStateMsg2(false, "%s == %s", id, result)
	}
	return result
}

// Java: ColorPool.createTrimmedOffsetToId(TypePool)
func colorPool_createTrimmedOffsetToId(typePool *TypePool) []colors.ColorId {
	typeList := typePool.GetTypeList()
	n := len(typeList)
	if n == 0 {
		return nil
	}
	out := make([]colors.ColorId, n)
	hasUnion := false

	for i, proto := range typeList {
		if proto.HasObject() {
			out[i] = colors.ColorId_FromBytes(proto.GetObject().GetUuid())
		} else if proto.HasUnion() {
			hasUnion = true
		} else {
			panic(NewMalformedTypedAstException(proto))
		}
	}

	if hasUnion {
		scratchMembers := jutil.NewLinkedHashSet[colors.ColorId]()
		for i, proto := range typeList {
			if proto.HasUnion() {
				out[i] = colorPool_createUnionColorId(proto, i, typeList, out, scratchMembers)
			}
		}
	}

	return out
}

func colorPool_createUnionColorId(
	proto *TypeProto,
	unionIndex int,
	typeList []*TypeProto,
	allIds []colors.ColorId,
	scratchMembers *jutil.LinkedHashSet[colors.ColorId],
) colors.ColorId {
	MalformedTypedAstException_CheckWellFormedWithParam(proto.GetUnion().GetUnionMemberCount() > 1, "Union has too few members", proto)
	scratchMembers.Clear()
	unionMemberList := proto.GetUnion().GetUnionMemberList()
	for i := 0; i < len(unionMemberList); i++ {
		memberPointer := int(unionMemberList[i])
		if TypePointers_IsAxiomatic(memberPointer) {
			scratchMembers.Add(TypePointers_OFFSET_TO_AXIOMATIC_COLOR[memberPointer].GetId())
		} else {
			trimmed := TypePointers_TrimOffset(memberPointer)
			MalformedTypedAstException_CheckWellFormedWithParam(
				trimmed >= 0 && trimmed < len(typeList) && (trimmed < unionIndex || typeList[trimmed].HasObject()),
				"Union member not found",
				proto)
			scratchMembers.Add(allIds[trimmed])
		}
	}
	return colors.ColorId_Union(scratchMembers)
}

// Java: ColorPool.validatePointer(int, ShardView)
func colorPool_validatePointer(offset int, shard *ColorPool_ShardView) int {
	if offset < 0 || offset >= TypePointers_UntrimOffset(len(shard.trimmedOffsetToId)) {
		MalformedTypedAstException_CheckWellFormedWithParam(
			false,
			"Pointer offset outside of shard",
			offset)
	}
	return offset
}

// 0xDEADBEEF as a Java int.
var colorPool_PENDING_COLOR = colors.Color_SingleBuilder().SetId(colors.ColorId_FromUnsignedInt(-0x21524111)).Build()
var colorPool_UNRECONCILED_COLOR = colors.Color_SingleBuilder().SetId(colors.ColorId_FromUnsignedInt(-0x21524112)).Build()
