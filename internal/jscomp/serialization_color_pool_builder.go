// Port of com/google/javascript/jscomp/serialization/ColorPool.java @ closure-compiler v20261028

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
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

type colorDedupBuf struct {
	buf    [8]*colors.Color
	extra  []*colors.Color
	n      int
	byHash map[int32][]*colors.Color
}

func (d *colorDedupBuf) add(c *colors.Color) {
	if d.extra != nil {
		d.extra = appendUniquePoolColor(d.extra, c, &d.byHash)
		return
	}
	for i := range d.n {
		existing := d.buf[i]
		if existing == c || existing.Equals(c) {
			return
		}
	}
	if d.n < len(d.buf) {
		d.buf[d.n] = c
		d.n++
		return
	}
	d.extra = make([]*colors.Color, len(d.buf)+1, 16)
	copy(d.extra, d.buf[:])
	d.extra[len(d.buf)] = c
}

func (d *colorDedupBuf) toSlice(b *ColorPool_Builder) []*colors.Color {
	if d.extra != nil {
		return b.copyColorSlice(d.extra)
	}
	if d.n == 0 {
		return colors.EmptyColorSlice()
	}
	return b.copyColorSlice(d.buf[:d.n])
}

type colorObjectAcc struct {
	inst                       colorDedupBuf
	protos                     colorDedupBuf
	closureAssert              base.Tri
	isConstructor              bool
	isInvalidating             bool
	propertiesKeepOriginalName bool
}

func (acc *colorObjectAcc) setClosureAssert(val bool, errParam any) {
	if acc.closureAssert.ToBoolean(val) != val {
		MalformedTypedAstException_CheckWellFormedWithParam(false, "Inconsistent values for closure_assert", errParam)
	}
	acc.closureAssert = base.Tri_ForBoolean(val)
}

func (acc *colorObjectAcc) addObjectProto(
	b *ColorPool_Builder,
	shard *ColorPool_ShardView,
	objProto *ObjectTypeProto,
	singleShard bool,
) int {
	for _, ptr := range objProto.GetInstanceTypeList() {
		acc.inst.add(b.resolveShardPtr(shard, int(ptr), singleShard))
	}
	acc.setClosureAssert(objProto.GetClosureAssert(), objProto)
	acc.isConstructor = acc.isConstructor || objProto.GetMarkedConstructor()
	acc.isInvalidating = acc.isInvalidating || objProto.GetIsInvalidating()
	acc.propertiesKeepOriginalName = acc.propertiesKeepOriginalName || objProto.GetPropertiesKeepOriginalName()
	for _, ptr := range objProto.GetPrototypeList() {
		acc.protos.add(b.resolveShardPtr(shard, int(ptr), singleShard))
	}
	return len(objProto.GetOwnPropertyList())
}

func (acc *colorObjectAcc) addJSTypes(
	b *ColorPool_Builder,
	s *JSTypeReconserializer,
	offsetToColor []*colors.Color,
	seen *JSTypeReconserializer_SeenTypeRecord,
) {
	jsClosureAssert := false
	for _, type_ := range seen.jstypes {
		objType := jutil.CheckNotNullMsg1(type_.ToMaybeObjectType(), "%s", type_)
		if objType.IsFunctionType() {
			fnType := objType.ToMaybeFunctionType()
			if fnType.HasInstanceType() && fnType.GetInstanceType() != nil {
				acc.inst.add(b.resolveJSTypePtr(s, offsetToColor, s.SerializeType(fnType.GetInstanceType())))
				acc.protos.add(b.resolveJSTypePtr(s, offsetToColor, s.SerializeType(fnType.GetPrototype())))
				acc.isConstructor = acc.isConstructor || fnType.IsConstructor()
			}
			jsClosureAssert = jsClosureAssert || JSTypeReconserializer_IsClosureAssert(fnType.GetClosurePrimitive())
		}
		acc.isInvalidating = acc.isInvalidating || s.invalidatingTypes.IsInvalidating(objType)
		acc.propertiesKeepOriginalName = acc.propertiesKeepOriginalName || objType.IsEnumType()
	}
	acc.setClosureAssert(jsClosureAssert, seen)
}

func (acc *colorObjectAcc) build(b *ColorPool_Builder, id colors.ColorId, ownProperties []string) *colors.Color {
	return b.colorBuilder.
		ClearBoxId().
		SetId(id).
		SetInstanceColors(acc.inst.toSlice(b)).
		SetPrototypes(acc.protos.toSlice(b)).
		SetOwnProperties(ownProperties).
		SetUnionElements(colors.EmptyColorSlice()).
		SetClosureAssert(acc.closureAssert.ToBoolean(false)).
		SetConstructor(acc.isConstructor).
		SetInvalidating(acc.isInvalidating).
		SetPropertiesKeepOriginalName(acc.propertiesKeepOriginalName).
		Build()
}

func (e colorPoolIdEntry) protoCount() int {
	if e.firstProto == nil {
		return 0
	}
	return 1 + len(e.extra)
}

func (e colorPoolIdEntry) shardProtoAt(i int) (*ColorPool_ShardView, *TypeProto) {
	if i == 0 {
		return e.firstShard, e.firstProto
	}
	sp := e.extra[i-1]
	return sp.shard, sp.proto
}

func (b *ColorPool_Builder) pushDebugStack(id colors.ColorId) {
	b.reconcilationDebugStack = append(b.reconcilationDebugStack, id)
}

func (b *ColorPool_Builder) popDebugStack() {
	b.reconcilationDebugStack = b.reconcilationDebugStack[:len(b.reconcilationDebugStack)-1]
}

func (b *ColorPool_Builder) panicCyclicColor(tail colors.ColorId, formatRow func(colors.ColorId) string) {
	rows := make([]string, 0, len(b.reconcilationDebugStack)+1)
	for _, stackId := range b.reconcilationDebugStack {
		rows = append(rows, formatRow(stackId))
	}
	rows = append(rows, formatRow(tail))
	panic(NewMalformedTypedAstException(
		"Cyclic Color structure detected: [" + strings.Join(rows, ", ") + "]"))
}

func (b *ColorPool_Builder) copyColorSlice(src []*colors.Color) []*colors.Color {
	if len(src) == 0 {
		return colors.EmptyColorSlice()
	}
	out := b.allocColorSlice(len(src))
	copy(out, src)
	return out
}

func (b *ColorPool_Builder) singleStringSlice(s string) []string {
	out := b.allocStringSlice(1)
	out[0] = s
	return out
}

func (b *ColorPool_Builder) flushPropScratch() []string {
	if len(b.propScratch) == 0 {
		return colors.EmptyStringSlice()
	}
	out := b.allocStringSlice(len(b.propScratch))
	copy(out, b.propScratch)
	return out
}

func (b *ColorPool_Builder) resetPropScratch(totalProps int) bool {
	b.propScratch = b.propScratch[:0]
	useMap := totalProps > 8
	if useMap {
		b.seenPropsGen++
		if b.seenProps == nil {
			b.seenProps = make(map[string]uint32, 64)
		} else if b.seenPropsGen == 0 {
			clear(b.seenProps)
			b.seenPropsGen = 1
		}
	}
	return useMap
}

func (b *ColorPool_Builder) collectProtoOwnProperties(shard *ColorPool_ShardView, ownPropList []int32) []string {
	if len(ownPropList) == 0 {
		return colors.EmptyStringSlice()
	}
	if len(ownPropList) == 1 {
		return b.singleStringSlice(shard.stringPool.Get(int(ownPropList[0])))
	}
	useMap := b.resetPropScratch(len(ownPropList))
	for _, p := range ownPropList {
		b.appendDedupProp(shard.stringPool.Get(int(p)), useMap)
	}
	return b.flushPropScratch()
}

func (b *ColorPool_Builder) collectJSTypeOwnProperties(seen *JSTypeReconserializer_SeenTypeRecord) []string {
	if len(seen.jstypes) == 1 {
		b.propScratch = b.propScratch[:0]
		seen.jstypes[0].ToMaybeObjectType().ForEachOwnPropertyName(b.directSingleOwnPropFn)
		return b.flushPropScratch()
	}
	b.startDirectPropDedup()
	for _, type_ := range seen.jstypes {
		type_.ToMaybeObjectType().ForEachOwnPropertyName(b.directMultiOwnPropFn)
	}
	return b.flushPropScratch()
}

func (b *ColorPool_Builder) collectEntryOwnProperties(entry colorPoolIdEntry, totalOwnProps int) []string {
	protoCount := entry.protoCount()
	if entry.jsRecord == nil {
		if totalOwnProps == 0 {
			return colors.EmptyStringSlice()
		}
		if totalOwnProps == 1 {
			for i := range protoCount {
				shard, proto := entry.shardProtoAt(i)
				if ownProps := proto.GetObject().GetOwnPropertyList(); len(ownProps) == 1 {
					return b.singleStringSlice(shard.stringPool.Get(int(ownProps[0])))
				}
			}
			return colors.EmptyStringSlice()
		}
		useMap := b.resetPropScratch(totalOwnProps)
		for i := range protoCount {
			shard, proto := entry.shardProtoAt(i)
			for _, p := range proto.GetObject().GetOwnPropertyList() {
				b.appendDedupProp(shard.stringPool.Get(int(p)), useMap)
			}
		}
		return b.flushPropScratch()
	}
	if protoCount == 0 {
		return b.collectJSTypeOwnProperties(entry.jsRecord)
	}
	b.startDirectPropDedup()
	for i := range protoCount {
		shard, proto := entry.shardProtoAt(i)
		for _, p := range proto.GetObject().GetOwnPropertyList() {
			b.appendDirectDedupProp(shard.stringPool.Get(int(p)))
		}
	}
	for _, type_ := range entry.jsRecord.jstypes {
		type_.ToMaybeObjectType().ForEachOwnPropertyName(b.directMultiOwnPropFn)
	}
	return b.flushPropScratch()
}

func (b *ColorPool_Builder) resolveJSTypePtr(s *JSTypeReconserializer, offsetToColor []*colors.Color, ptr int) *colors.Color {
	if offsetToColor != nil {
		return b.reconcileDirectJSTypeAt(s, offsetToColor, ptr)
	}
	return b.LookupOrReconcileColor(s.recordsByPointer[ptr].colorId)
}

func (b *ColorPool_Builder) resolveShardPtr(shard *ColorPool_ShardView, untrimmedOffset int, singleShard bool) *colors.Color {
	if singleShard {
		return b.resolveSingleShardPointer(shard, untrimmedOffset)
	}
	return b.LookupOrReconcileColor(shard.GetId(untrimmedOffset))
}

type colorUnionBuf struct {
	buf   [16]*colors.Color
	extra []*colors.Color
	n     int
}

func (u *colorUnionBuf) add(member *colors.Color, errParam any) {
	if member.IsUnion() {
		MalformedTypedAstException_CheckWellFormedWithParam(false, "Reconciling union with non-union", errParam)
	}
	if u.extra != nil {
		u.extra = append(u.extra, member)
		return
	}
	if u.n < len(u.buf) {
		u.buf[u.n] = member
		u.n++
		return
	}
	u.extra = make([]*colors.Color, len(u.buf)+1, 32)
	copy(u.extra, u.buf[:])
	u.extra[len(u.buf)] = member
}

func (u *colorUnionBuf) addJSTypeMembers(
	b *ColorPool_Builder,
	s *JSTypeReconserializer,
	offsetToColor []*colors.Color,
	seen *JSTypeReconserializer_SeenTypeRecord,
) {
	var ptrBuf [16]int32
	var ptrs []int32
	if len(seen.unionMembers) <= len(ptrBuf) {
		ptrs = ptrBuf[:len(seen.unionMembers)]
	} else {
		ptrs = make([]int32, len(seen.unionMembers))
	}
	for i, r := range seen.unionMembers {
		ptrs[i] = int32(r.pointer)
	}
	slices.Sort(ptrs)
	for _, memberPtr := range ptrs {
		u.add(b.resolveJSTypePtr(s, offsetToColor, int(memberPtr)), seen)
	}
}

func (u *colorUnionBuf) build(b *ColorPool_Builder, id colors.ColorId) *colors.Color {
	slice := u.buf[:u.n]
	if u.extra != nil {
		slice = u.extra
	}
	result := colors.Color_CreateUnionSliceWithBuilder(slice, &b.colorBuilder, b.allocColorSliceFn, b.allocStringSliceFn)
	if id != result.GetId() {
		jutil.CheckStateMsg2(false, "%s == %s", id, result)
	}
	return result
}

func (b *ColorPool_Builder) reconcileDirectJSTypeAt(
	s *JSTypeReconserializer,
	offsetToColor []*colors.Color,
	ptr int,
) *colors.Color {
	if existing := offsetToColor[ptr]; existing != nil {
		if existing == colorPool_PENDING_COLOR {
			b.panicCyclicColor(s.recordsByPointer[ptr].colorId, colors.ColorId.String)
		}
		return existing
	}

	seen := s.recordsByPointer[ptr]
	id := seen.colorId
	b.pushDebugStack(id)
	defer b.popDebugStack()
	offsetToColor[ptr] = colorPool_PENDING_COLOR

	var result *colors.Color
	if len(seen.unionMembers) != 0 {
		var u colorUnionBuf
		u.addJSTypeMembers(b, s, offsetToColor, seen)
		result = u.build(b, id)
	} else {
		acc := colorObjectAcc{closureAssert: base.Tri_UNKNOWN}
		acc.addJSTypes(b, s, offsetToColor, seen)
		result = acc.build(b, id, b.collectJSTypeOwnProperties(seen))
	}

	if result == nil {
		jutil.CheckStateMsg1(false, "%s", id)
	}
	offsetToColor[ptr] = result
	b.idToColor[id] = result
	return result
}

func (b *ColorPool_Builder) resolveSingleShardPointer(shard *ColorPool_ShardView, untrimmedOffset int) *colors.Color {
	if TypePointers_IsAxiomatic(untrimmedOffset) {
		return TypePointers_OFFSET_TO_AXIOMATIC_COLOR[untrimmedOffset]
	}
	return b.reconcileSingleShardAt(shard, TypePointers_TrimOffset(untrimmedOffset))
}

func (b *ColorPool_Builder) reconcileSingleShardAt(shard *ColorPool_ShardView, trimmedIdx int) *colors.Color {
	id := shard.trimmedOffsetToId[trimmedIdx]
	if existing := shard.trimmedOffsetToColor[trimmedIdx]; existing != nil {
		if existing == colorPool_PENDING_COLOR {
			typeList := shard.typePool.GetTypeList()
			b.panicCyclicColor(id, func(stackId colors.ColorId) string {
				for i, sid := range shard.trimmedOffsetToId {
					if sid == stackId {
						return (colorPoolIdEntry{firstShard: shard, firstProto: typeList[i]}).rowString()
					}
				}
				return ""
			})
		}
		return existing
	}
	b.pushDebugStack(id)
	defer b.popDebugStack()
	shard.trimmedOffsetToColor[trimmedIdx] = colorPool_PENDING_COLOR

	proto := shard.typePool.GetTypeList()[trimmedIdx]
	var result *colors.Color
	if proto.HasObject() {
		objProto := proto.GetObject()
		acc := colorObjectAcc{closureAssert: base.Tri_UNKNOWN}
		acc.addObjectProto(b, shard, objProto, true)
		result = acc.build(b, id, b.collectProtoOwnProperties(shard, objProto.GetOwnPropertyList()))
	} else if proto.HasUnion() {
		var u colorUnionBuf
		for _, memberPtr := range proto.GetUnion().GetUnionMemberList() {
			u.add(b.resolveSingleShardPointer(shard, int(memberPtr)), proto)
		}
		result = u.build(b, id)
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
	if existing := b.idToColor[id]; existing != nil {
		if existing == colorPool_PENDING_COLOR {
			b.panicCyclicColor(id, func(stackId colors.ColorId) string {
				return b.protoEntries[b.idToProtoIdx[stackId]].rowString()
			})
		}
		return existing
	}
	b.pushDebugStack(id)
	defer b.popDebugStack()
	b.idToColor[id] = colorPool_PENDING_COLOR

	entry := b.protoEntries[b.idToProtoIdx[id]]
	var result *colors.Color
	if sample := entry.firstProto; sample != nil {
		if sample.HasObject() {
			result = b.reconcileObjectProtos(id, entry)
		} else if sample.HasUnion() {
			result = b.reconcileUnionProtos(id, entry)
		} else {
			panic(jutil.NewAssertionError(sample.String()))
		}
	} else if seen := entry.jsRecord; seen != nil {
		if len(seen.unionMembers) == 0 {
			result = b.reconcileObjectProtos(id, entry)
		} else {
			result = b.reconcileUnionProtos(id, entry)
		}
	} else {
		jutil.CheckNotNullMsg1(entry.firstProto, "%s", id)
	}

	if result == nil {
		jutil.CheckStateMsg1(false, "%s", id)
	}
	b.idToColor[id] = result
	return result
}

func (b *ColorPool_Builder) reconcileObjectProtos(id colors.ColorId, entry colorPoolIdEntry) *colors.Color {
	acc := colorObjectAcc{closureAssert: base.Tri_UNKNOWN}
	totalOwnProps := 0
	protoCount := entry.protoCount()
	for i := range protoCount {
		shard, proto := entry.shardProtoAt(i)
		jutil.CheckState(proto.HasObject())
		totalOwnProps += acc.addObjectProto(b, shard, proto.GetObject(), false)
	}
	if seen := entry.jsRecord; seen != nil {
		jutil.CheckState(len(seen.unionMembers) == 0)
		acc.addJSTypes(b, b.jsTypeShard, nil, seen)
	}
	return acc.build(b, id, b.collectEntryOwnProperties(entry, totalOwnProps))
}

func (b *ColorPool_Builder) reconcileUnionProtos(id colors.ColorId, entry colorPoolIdEntry) *colors.Color {
	var u colorUnionBuf
	protoCount := entry.protoCount()
	for i := range protoCount {
		shard, proto := entry.shardProtoAt(i)
		jutil.CheckStateMsg1(proto.HasUnion(), "%s", proto)
		for _, memberPtr := range proto.GetUnion().GetUnionMemberList() {
			u.add(b.LookupOrReconcileColor(shard.GetId(int(memberPtr))), proto)
		}
	}
	if seen := entry.jsRecord; seen != nil {
		jutil.CheckStateMsg1(len(seen.unionMembers) > 1, "%s", seen)
		u.addJSTypeMembers(b, b.jsTypeShard, nil, seen)
	}
	return u.build(b, id)
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
		var members smallSet[colors.ColorId]
		for i, proto := range typeList {
			if proto.HasUnion() {
				out[i] = colorPool_createUnionColorId(proto, i, typeList, out, &members)
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
	members *smallSet[colors.ColorId],
) colors.ColorId {
	unionMemberList := proto.GetUnion().GetUnionMemberList()
	MalformedTypedAstException_CheckWellFormedWithParam(len(unionMemberList) > 1, "Union has too few members", proto)
	members.clear()
	for _, rawPtr := range unionMemberList {
		memberPointer := int(rawPtr)
		var cid colors.ColorId
		if TypePointers_IsAxiomatic(memberPointer) {
			MalformedTypedAstException_CheckWellFormedWithParam(
				memberPointer >= 0,
				"Union member not found",
				proto)
			cid = TypePointers_OFFSET_TO_AXIOMATIC_COLOR[memberPointer].GetId()
		} else {
			trimmed := TypePointers_TrimOffset(memberPointer)
			MalformedTypedAstException_CheckWellFormedWithParam(
				trimmed >= 0 && trimmed < len(typeList) && (trimmed < unionIndex || typeList[trimmed].HasObject()),
				"Union member not found",
				proto)
			cid = allIds[trimmed]
		}
		members.add(cid)
	}
	return unionColorIds(members.items)
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
