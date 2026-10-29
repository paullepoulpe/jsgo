// Port of com/google/javascript/jscomp/serialization/ColorPool.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/serialization/TypePointers.java @ closure-compiler v20261021

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

/*
 * Copyright 2021 The Closure Compiler Authors.
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
	"encoding/binary"
	"fmt"
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Java: TypePointers.AXIOMATIC_COLOR_COUNT
const TypePointers_AXIOMATIC_COLOR_COUNT = int(PrimitiveType_UNRECOGNIZED)

// TypePointers_OFFSET_TO_AXIOMATIC_COLOR lists the axiomatic colors that the
// first N TypePointer pool offsets correspond to, in this order.
//
// These colors are never serialized because all their information is
// constant.
var TypePointers_OFFSET_TO_AXIOMATIC_COLOR = [TypePointers_AXIOMATIC_COLOR_COUNT]*colors.Color{
	PrimitiveType_UNKNOWN_TYPE:      colors.StandardColors_UNKNOWN,
	PrimitiveType_BOOLEAN_TYPE:      colors.StandardColors_BOOLEAN,
	PrimitiveType_STRING_TYPE:       colors.StandardColors_STRING,
	PrimitiveType_NUMBER_TYPE:       colors.StandardColors_NUMBER,
	PrimitiveType_NULL_OR_VOID_TYPE: colors.StandardColors_NULL_OR_VOID,
	PrimitiveType_SYMBOL_TYPE:       colors.StandardColors_SYMBOL,
	PrimitiveType_BIGINT_TYPE:       colors.StandardColors_BIGINT,
	PrimitiveType_TOP_OBJECT:        colors.StandardColors_TOP_OBJECT,
	PrimitiveType_TOP_FUNCTION:      colors.StandardColors_TOP_FUNCTION,
	PrimitiveType_GBIGINT_TYPE:      colors.StandardColors_GBIGINT,
}

func isAxiomaticColorId(id colors.ColorId) bool {
	for _, c := range TypePointers_OFFSET_TO_AXIOMATIC_COLOR {
		if c.GetId() == id {
			return true
		}
	}
	return false
}

// Java: TypePointers.trimOffset(int)
func TypePointers_TrimOffset(x int) int {
	jutil.CheckStateMsg1(x >= TypePointers_AXIOMATIC_COLOR_COUNT, "%s", x)
	return x - TypePointers_AXIOMATIC_COLOR_COUNT
}

// Java: TypePointers.untrimOffset(int)
func TypePointers_UntrimOffset(x int) int {
	jutil.CheckStateMsg1(x >= 0, "%s", x)
	return x + TypePointers_AXIOMATIC_COLOR_COUNT
}

// Java: TypePointers.isAxiomatic(int)
func TypePointers_IsAxiomatic(x int) bool {
	return x < TypePointers_AXIOMATIC_COLOR_COUNT
}

func colorIdRawInt64(id colors.ColorId) int64 {
	var b [8]byte
	id.PutBigEndian(b[:])
	return int64(binary.BigEndian.Uint64(b[:]))
}

// unionColorIds computes colors.ColorId_Union over a slice of already-unique ColorIds.
func unionColorIds(ids []colors.ColorId) colors.ColorId {
	n := len(ids)
	if n == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	if n == 1 {
		return ids[0]
	}
	if n == 2 {
		a, b := colorIdRawInt64(ids[0]), colorIdRawInt64(ids[1])
		if a > b {
			a, b = b, a
		}
		var buf [16]byte
		binary.LittleEndian.PutUint64(buf[0:8], uint64(a))
		binary.LittleEndian.PutUint64(buf[8:16], uint64(b))
		return colors.NewColorId(int64(jutil.FarmHashFingerprint64(buf[:])))
	}

	var stackBuf [32]int64
	var stackBytes [256]byte
	var sorted []int64
	var raw []byte
	if n <= len(stackBuf) {
		sorted = stackBuf[:n]
		raw = stackBytes[:n*8]
	} else {
		sorted = make([]int64, n)
		raw = make([]byte, n*8)
	}
	for i, cid := range ids {
		sorted[i] = colorIdRawInt64(cid)
	}
	slices.Sort(sorted)
	for i, v := range sorted {
		binary.LittleEndian.PutUint64(raw[i*8:], uint64(v))
	}
	return colors.NewColorId(int64(jutil.FarmHashFingerprint64(raw)))
}

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
		sb.WriteString(fmt.Sprint(list[0]))
		i := 1
		for ; i < 5 && i < len(list); i++ {
			sb.WriteString(", ")
			sb.WriteString(fmt.Sprint(list[i]))
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
	}
	return s.trimmedOffsetToId[TypePointers_TrimOffset(untrimmedOffset)]
}

// Java: ColorPool.ShardView.updateStateAfterColorPoolIsBuilt(ColorPool)
func (s *ColorPool_ShardView) UpdateStateAfterColorPoolIsBuilt(colorPool *ColorPool) {
	s.built = true
	s.typePool = nil
	s.stringPool = nil
	axiomaticCount := TypePointers_AXIOMATIC_COLOR_COUNT
	offsetToColor := make([]*colors.Color, axiomaticCount+len(s.trimmedOffsetToId))
	copy(offsetToColor, TypePointers_OFFSET_TO_AXIOMATIC_COLOR[:])
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
	id         colors.ColorId
	firstShard *ColorPool_ShardView
	firstProto *TypeProto
	extra      []colorPoolShardProto
	jsRecord   *JSTypeReconserializer_SeenTypeRecord
}

func (e colorPoolIdEntry) rowString() string {
	if e.firstProto == nil && len(e.extra) == 0 {
		if e.jsRecord != nil {
			return e.jsRecord.String()
		}
		return "{}"
	}
	var sb strings.Builder
	sb.WriteByte('{')
	first := true
	if e.firstProto != nil {
		sb.WriteString(e.firstShard.String())
		sb.WriteByte('=')
		sb.WriteString(e.firstProto.String())
		first = false
	}
	for _, sp := range e.extra {
		if !first {
			sb.WriteString(", ")
		}
		sb.WriteString(sp.shard.String())
		sb.WriteByte('=')
		sb.WriteString(sp.proto.String())
		first = false
	}
	sb.WriteByte('}')
	return sb.String()
}

// ColorPool_Builder collects TypePools and other data into ShardViews and
// then reconciles them into a single ColorPool.
type ColorPool_Builder struct {
	protoToShard map[*TypePool]*ColorPool_ShardView
	shards       []*ColorPool_ShardView
	jsTypeShard  *JSTypeReconserializer
	idToColor    map[colors.ColorId]*colors.Color
	registry     *colors.ColorRegistry_Builder
	idToProtoIdx map[colors.ColorId]int
	protoEntries []colorPoolIdEntry
	forTesting_  bool
	built        bool

	colorBuilder          colors.Color_Builder
	colorSlab             jutil.Slab[colors.Color]
	colorSliceArena       []*colors.Color
	stringSliceArena      []string
	allocColorSliceFn     func(n int) []*colors.Color
	allocStringSliceFn    func(n int) []string
	propScratch           []string
	seenProps             map[string]uint32
	seenPropsGen          uint32
	directPropMapActive   bool
	directSingleOwnPropFn func(string)
	directMultiOwnPropFn  func(string)

	reconcilationDebugStack []colors.ColorId
}

func (b *ColorPool_Builder) allocColorSlice(n int) []*colors.Color {
	if n == 0 {
		return colors.EmptyColorSlice()
	}
	const chunk = 4096 // 32 KiB
	if n > chunk {
		return make([]*colors.Color, n)
	}
	if len(b.colorSliceArena) < n {
		b.colorSliceArena = make([]*colors.Color, chunk)
	}
	out := b.colorSliceArena[:n:n]
	b.colorSliceArena = b.colorSliceArena[n:]
	return out
}

func (b *ColorPool_Builder) allocStringSlice(n int) []string {
	if n == 0 {
		return colors.EmptyStringSlice()
	}
	const chunk = 2048 // 32 KiB
	if n > chunk {
		return make([]string, n)
	}
	if len(b.stringSliceArena) < n {
		b.stringSliceArena = make([]string, chunk)
	}
	out := b.stringSliceArena[:n:n]
	b.stringSliceArena = b.stringSliceArena[n:]
	return out
}

func appendUniquePoolColor(dst []*colors.Color, c *colors.Color, byHash *map[int32][]*colors.Color) []*colors.Color {
	if *byHash != nil {
		h := c.HashCode()
		bucket := (*byHash)[h]
		for _, existing := range bucket {
			if existing == c || existing.Equals(c) {
				return dst
			}
		}
		(*byHash)[h] = append(bucket, c)
		return append(dst, c)
	}
	for _, existing := range dst {
		if existing == c || existing.Equals(c) {
			return dst
		}
	}
	dst = append(dst, c)
	if len(dst) > 16 {
		m := make(map[int32][]*colors.Color, len(dst)*2)
		for _, existing := range dst {
			eh := existing.HashCode()
			m[eh] = append(m[eh], existing)
		}
		*byHash = m
	}
	return dst
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

func (b *ColorPool_Builder) startDirectPropDedup() {
	b.propScratch = b.propScratch[:0]
	b.directPropMapActive = false
}

func (b *ColorPool_Builder) appendDirectDedupProp(prop string) {
	if !b.directPropMapActive {
		for _, existing := range b.propScratch {
			if existing == prop {
				return
			}
		}
		b.propScratch = append(b.propScratch, prop)
		if len(b.propScratch) == 8 {
			b.directPropMapActive = true
			b.seenPropsGen++
			if b.seenProps == nil {
				b.seenProps = make(map[string]uint32, 64)
			} else if b.seenPropsGen == 0 {
				clear(b.seenProps)
				b.seenPropsGen = 1
			}
			for _, existing := range b.propScratch {
				b.seenProps[existing] = b.seenPropsGen
			}
		}
		return
	}
	if b.seenProps[prop] == b.seenPropsGen {
		return
	}
	b.seenProps[prop] = b.seenPropsGen
	b.propScratch = append(b.propScratch, prop)
}

func (b *ColorPool_Builder) onDirectSingleOwnProp(prop string) {
	if b.jsTypeShard.shouldPropagate(prop) {
		b.propScratch = append(b.propScratch, prop)
	}
}

func (b *ColorPool_Builder) onDirectMultiOwnProp(prop string) {
	if b.jsTypeShard.shouldPropagate(prop) {
		b.appendDirectDedupProp(prop)
	}
}

func (b *ColorPool_Builder) addJSTypeShard(s *JSTypeReconserializer) {
	jutil.CheckStateMsg(!b.built, "build has already been called")
	jutil.CheckState(b.jsTypeShard == nil)
	b.jsTypeShard = s
	b.directSingleOwnPropFn = b.onDirectSingleOwnProp
	b.directMultiOwnPropFn = b.onDirectMultiOwnProp
	for _, m := range s.debugMismatches {
		for _, ptr := range m.involvedPointers {
			b.registry.AddMismatchLocation(s.recordsByPointer[ptr].colorId, m.sourceRef)
		}
	}
}

// Java: ColorPool.Builder()
func newColorPool_Builder() *ColorPool_Builder {
	b := &ColorPool_Builder{
		protoToShard: make(map[*TypePool]*ColorPool_ShardView, 4),
		registry:     colors.ColorRegistry_Builder_(),
		colorSlab:    jutil.NewFixedSlab[colors.Color](224), // 32.25 KiB <= 32 KiB
	}
	b.colorBuilder.SetSlab(&b.colorSlab)
	b.colorBuilder.SetTrustUniqueSlices(true)
	b.allocColorSliceFn = b.allocColorSlice
	b.allocStringSliceFn = b.allocStringSlice
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
		jutil.CheckStateMsg1(typePool == TypePool_GetDefaultInstance(), "%s", typePool)
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
	shardCount := len(b.shards)
	var singleShard *ColorPool_ShardView
	for _, shard := range b.shards {
		totalTypes += len(shard.typePool.GetTypeList())
		singleShard = shard
	}
	if b.jsTypeShard != nil {
		totalTypes += len(b.jsTypeShard.recordsByPointer)
	}

	b.idToColor = make(map[colors.ColorId]*colors.Color, totalTypes+16)

	if b.jsTypeShard != nil && shardCount == 0 {
		s := b.jsTypeShard
		records := s.recordsByPointer
		axiomaticCount := TypePointers_AXIOMATIC_COLOR_COUNT
		offsetToColor := make([]*colors.Color, len(records))
		copy(offsetToColor, TypePointers_OFFSET_TO_AXIOMATIC_COLOR[:])
		for _, c := range TypePointers_OFFSET_TO_AXIOMATIC_COLOR {
			b.idToColor[c.GetId()] = c
		}
		for ptr := axiomaticCount; ptr < len(records); ptr++ {
			b.reconcileDirectJSTypeAt(s, offsetToColor, ptr)
		}
		s.offsetToColor = offsetToColor
		s.finishColorPool()
	} else if b.jsTypeShard == nil && shardCount <= 1 {
		if singleShard != nil && len(singleShard.trimmedOffsetToId) > 0 {
			typeList := singleShard.typePool.GetTypeList()
			for i, id := range singleShard.trimmedOffsetToId {
				if b.idToColor[id] != nil {
					MalformedTypedAstException_CheckWellFormedWithParam(false, "Duplicate ID in single shard", typeList[i])
				}
				b.idToColor[id] = colorPool_UNRECONCILED_COLOR
			}
			for _, v := range TypePointers_OFFSET_TO_AXIOMATIC_COLOR {
				id := v.GetId()
				if b.idToColor[id] != nil {
					MalformedTypedAstException_CheckWellFormedWithParam(
						false, "Found serialized definiton for axiomatic color", id)
				}
				b.idToColor[id] = v
			}
			singleShard.trimmedOffsetToColor = make([]*colors.Color, len(typeList))
			for i := range typeList {
				b.reconcileSingleShardAt(singleShard, i)
			}
		} else {
			for _, v := range TypePointers_OFFSET_TO_AXIOMATIC_COLOR {
				b.idToColor[v.GetId()] = v
			}
		}
	} else {
		for _, v := range TypePointers_OFFSET_TO_AXIOMATIC_COLOR {
			b.idToColor[v.GetId()] = v
		}
		b.idToProtoIdx = make(map[colors.ColorId]int, totalTypes)
		b.protoEntries = make([]colorPoolIdEntry, 0, totalTypes)
		for _, shard := range b.shards {
			typeList := shard.typePool.GetTypeList()
			seenIds := make(map[colors.ColorId]struct{}, len(typeList))
			for i, proto := range typeList {
				id := shard.trimmedOffsetToId[i]
				if _, dup := seenIds[id]; dup {
					MalformedTypedAstException_CheckWellFormedWithParam(false, "Duplicate ID in single shard", proto)
				}
				seenIds[id] = struct{}{}
				if entryIdx, ok := b.idToProtoIdx[id]; !ok {
					b.idToProtoIdx[id] = len(b.protoEntries)
					b.protoEntries = append(b.protoEntries, colorPoolIdEntry{
						id:         id,
						firstShard: shard,
						firstProto: proto,
					})
				} else {
					entry := &b.protoEntries[entryIdx]
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
				}
			}
		}

		if s := b.jsTypeShard; s != nil {
			axiomaticCount := TypePointers_AXIOMATIC_COLOR_COUNT
			records := s.recordsByPointer
			for ptr := axiomaticCount; ptr < len(records); ptr++ {
				seen := records[ptr]
				id := seen.colorId
				if entryIdx, ok := b.idToProtoIdx[id]; !ok {
					b.idToProtoIdx[id] = len(b.protoEntries)
					b.protoEntries = append(b.protoEntries, colorPoolIdEntry{id: id, jsRecord: seen})
				} else {
					b.protoEntries[entryIdx].jsRecord = seen
				}
			}
		}

		for _, c := range TypePointers_OFFSET_TO_AXIOMATIC_COLOR {
			id := c.GetId()
			if _, exists := b.idToProtoIdx[id]; exists {
				MalformedTypedAstException_CheckWellFormedWithParam(
					false, "Found serialized definiton for axiomatic color", id)
			}
		}

		for i := range b.protoEntries {
			b.LookupOrReconcileColor(b.protoEntries[i].id)
		}

		if s := b.jsTypeShard; s != nil {
			axiomaticCount := TypePointers_AXIOMATIC_COLOR_COUNT
			records := s.recordsByPointer
			offsetToColor := make([]*colors.Color, len(records))
			copy(offsetToColor, TypePointers_OFFSET_TO_AXIOMATIC_COLOR[:])
			for ptr := axiomaticCount; ptr < len(records); ptr++ {
				offsetToColor[ptr] = b.idToColor[records[ptr].colorId]
			}
			s.offsetToColor = offsetToColor
			s.finishColorPool()
		}
	}

	emptyColors := colors.EmptyColorSlice()
	emptyStrings := colors.EmptyStringSlice()
	for colorId := range colors.ColorRegistry_REQUIRED_IDS.All() {
		c := b.idToColor[colorId]
		if c == nil {
			c = b.colorBuilder.
				ClearBoxId().
				SetId(colorId).
				SetInstanceColors(emptyColors).
				SetPrototypes(emptyColors).
				SetOwnProperties(emptyStrings).
				SetUnionElements(emptyColors).
				SetClosureAssert(false).
				SetConstructor(false).
				SetInvalidating(false).
				SetPropertiesKeepOriginalName(false).
				Build()
			b.idToColor[colorId] = c
		}
		b.registry.SetNativeColor(c)
	}

	if b.jsTypeShard == nil && shardCount == 1 && singleShard != nil {
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
		for _, shard := range b.shards {
			for _, edge := range shard.typePool.GetDisambiguationEdgesList() {
				b.registry.AddDisambiguationEdge(
					b.idToColor[shard.GetId(colorPool_validatePointer(int(edge.GetSubtype()), shard))],
					b.idToColor[shard.GetId(colorPool_validatePointer(int(edge.GetSupertype()), shard))])
			}
		}
	}
	if s := b.jsTypeShard; s != nil && s.totalDisambiguateEdges > 0 {
		offsetToColor := s.offsetToColor
		for _, subtype := range s.disambiguateSubtypes {
			subColor := offsetToColor[subtype]
			for _, sup := range s.disambiguateEdges[subtype] {
				b.registry.AddDisambiguationEdge(subColor, offsetToColor[sup])
			}
		}
	}

	shardViewsForTesting := []*ColorPool_ShardView{}
	if b.forTesting_ {
		shardViewsForTesting = append(shardViewsForTesting, b.shards...)
	}

	colorPool := newColorPool(b.idToColor, b.registry.Build(), shardViewsForTesting)

	for _, shard := range b.shards {
		shard.UpdateStateAfterColorPoolIsBuilt(colorPool)
	}
	b.idToProtoIdx = nil
	b.protoEntries = nil
	b.jsTypeShard = nil
	b.propScratch = nil
	b.seenProps = nil
	b.colorSliceArena = nil
	b.stringSliceArena = nil
	return colorPool
}
