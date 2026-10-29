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
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func (b *ColorPool_Builder) reconcileDirectJSTypeAt(
	s *JSTypeReconserializer,
	offsetToColor []*colors.Color,
	ptr int,
) *colors.Color {
	if existing := offsetToColor[ptr]; existing != nil {
		if existing == colorPool_PENDING_COLOR {
			rows := make([]string, 0, len(b.reconcilationDebugStack)+1)
			for _, stackId := range b.reconcilationDebugStack {
				rows = append(rows, stackId.String())
			}
			rows = append(rows, s.recordsByPointer[ptr].colorId.String())
			panic(NewMalformedTypedAstException(
				"Cyclic Color structure detected: [" + strings.Join(rows, ", ") + "]"))
		}
		return existing
	}

	seen := s.recordsByPointer[ptr]
	id := seen.colorId
	b.reconcilationDebugStack = append(b.reconcilationDebugStack, id)
	defer func() {
		b.reconcilationDebugStack = b.reconcilationDebugStack[:len(b.reconcilationDebugStack)-1]
	}()
	offsetToColor[ptr] = colorPool_PENDING_COLOR

	var result *colors.Color
	if len(seen.unionMembers) != 0 {
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

		var unionBuf [16]*colors.Color
		union := unionBuf[:0]
		for _, memberPtr := range ptrs {
			member := b.reconcileDirectJSTypeAt(s, offsetToColor, int(memberPtr))
			MalformedTypedAstException_CheckWellFormedWithParam(!member.IsUnion(), "Reconciling union with non-union", seen)
			union = append(union, member)
		}
		result = colors.Color_CreateUnionSliceWithBuilder(union, &b.colorBuilder, b.allocColorSliceFn, b.allocStringSliceFn)
		if id != result.GetId() {
			jutil.CheckStateMsg2(false, "%s == %s", id, result)
		}
	} else if len(seen.jstypes) == 1 {
		type_ := seen.jstypes[0]
		objType := jutil.CheckNotNullMsg1(type_.ToMaybeObjectType(), "%s", type_)

		instanceColors := colors.EmptyColorSlice()
		prototypes := colors.EmptyColorSlice()
		isClosureAssert := false
		isConstructor := false

		if objType.IsFunctionType() {
			fnType := objType.ToMaybeFunctionType()
			if fnType.HasInstanceType() && fnType.GetInstanceType() != nil {
				instColor := b.reconcileDirectJSTypeAt(s, offsetToColor, s.SerializeType(fnType.GetInstanceType()))
				instanceColors = b.allocColorSlice(1)
				instanceColors[0] = instColor
				protoColor := b.reconcileDirectJSTypeAt(s, offsetToColor, s.SerializeType(fnType.GetPrototype()))
				prototypes = b.allocColorSlice(1)
				prototypes[0] = protoColor
				isConstructor = fnType.IsConstructor()
			}
			isClosureAssert = JSTypeReconserializer_IsClosureAssert(fnType.GetClosurePrimitive())
		}

		b.propScratch = b.propScratch[:0]
		objType.ForEachOwnPropertyName(b.directSingleOwnPropFn)
		ownProperties := colors.EmptyStringSlice()
		if len(b.propScratch) > 0 {
			ownProperties = b.allocStringSlice(len(b.propScratch))
			copy(ownProperties, b.propScratch)
		}

		result = b.colorBuilder.
			ClearBoxId().
			SetId(id).
			SetInstanceColors(instanceColors).
			SetPrototypes(prototypes).
			SetOwnProperties(ownProperties).
			SetUnionElements(colors.EmptyColorSlice()).
			SetClosureAssert(isClosureAssert).
			SetConstructor(isConstructor).
			SetInvalidating(s.invalidatingTypes.IsInvalidating(objType)).
			SetPropertiesKeepOriginalName(objType.IsEnumType()).
			Build()
	} else {
		var instBuf [8]*colors.Color
		inst := instBuf[:0]
		var instByHash map[int32][]*colors.Color

		var protoBuf [8]*colors.Color
		protos := protoBuf[:0]
		var protoByHash map[int32][]*colors.Color

		isClosureAssert := false
		isConstructor := false
		isInvalidating := false
		propertiesKeepOriginalName := false

		for _, type_ := range seen.jstypes {
			objType := jutil.CheckNotNullMsg1(type_.ToMaybeObjectType(), "%s", type_)
			if objType.IsFunctionType() {
				fnType := objType.ToMaybeFunctionType()
				if fnType.HasInstanceType() && fnType.GetInstanceType() != nil {
					inst = appendUniquePoolColor(
						inst,
						b.reconcileDirectJSTypeAt(s, offsetToColor, s.SerializeType(fnType.GetInstanceType())),
						&instByHash)
					protos = appendUniquePoolColor(
						protos,
						b.reconcileDirectJSTypeAt(s, offsetToColor, s.SerializeType(fnType.GetPrototype())),
						&protoByHash)
					isConstructor = isConstructor || fnType.IsConstructor()
				}
				isClosureAssert = isClosureAssert || JSTypeReconserializer_IsClosureAssert(fnType.GetClosurePrimitive())
			}
			isInvalidating = isInvalidating || s.invalidatingTypes.IsInvalidating(objType)
			propertiesKeepOriginalName = propertiesKeepOriginalName || objType.IsEnumType()
		}

		instanceColors := colors.EmptyColorSlice()
		if len(inst) > 0 {
			instanceColors = b.allocColorSlice(len(inst))
			copy(instanceColors, inst)
		}
		prototypes := colors.EmptyColorSlice()
		if len(protos) > 0 {
			prototypes = b.allocColorSlice(len(protos))
			copy(prototypes, protos)
		}

		b.startDirectPropDedup()
		for _, type_ := range seen.jstypes {
			type_.ToMaybeObjectType().ForEachOwnPropertyName(b.directMultiOwnPropFn)
		}
		ownProperties := colors.EmptyStringSlice()
		if len(b.propScratch) > 0 {
			ownProperties = b.allocStringSlice(len(b.propScratch))
			copy(ownProperties, b.propScratch)
		}

		result = b.colorBuilder.
			ClearBoxId().
			SetId(id).
			SetInstanceColors(instanceColors).
			SetPrototypes(prototypes).
			SetOwnProperties(ownProperties).
			SetUnionElements(colors.EmptyColorSlice()).
			SetClosureAssert(isClosureAssert).
			SetConstructor(isConstructor).
			SetInvalidating(isInvalidating).
			SetPropertiesKeepOriginalName(propertiesKeepOriginalName).
			Build()
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
	b.reconcilationDebugStack = append(b.reconcilationDebugStack, id)
	defer func() {
		b.reconcilationDebugStack = b.reconcilationDebugStack[:len(b.reconcilationDebugStack)-1]
	}()
	existing := shard.trimmedOffsetToColor[trimmedIdx]
	if existing != nil {
		if existing == colorPool_PENDING_COLOR {
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
		instanceColors := colors.EmptyColorSlice()
		if instanceTypeIds := objProto.GetInstanceTypeList(); len(instanceTypeIds) == 1 {
			resolved := b.resolveSingleShardPointer(shard, int(instanceTypeIds[0]))
			instanceColors = b.allocColorSlice(1)
			instanceColors[0] = resolved
		} else if len(instanceTypeIds) > 1 {
			var instBuf [8]*colors.Color
			inst := instBuf[:0]
			var instByHash map[int32][]*colors.Color
			for i := 0; i < len(instanceTypeIds); i++ {
				inst = appendUniquePoolColor(inst, b.resolveSingleShardPointer(shard, int(instanceTypeIds[i])), &instByHash)
			}
			instanceColors = b.allocColorSlice(len(inst))
			copy(instanceColors, inst)
		}

		prototypes := colors.EmptyColorSlice()
		if prototypeIds := objProto.GetPrototypeList(); len(prototypeIds) == 1 {
			resolved := b.resolveSingleShardPointer(shard, int(prototypeIds[0]))
			prototypes = b.allocColorSlice(1)
			prototypes[0] = resolved
		} else if len(prototypeIds) > 1 {
			var protoBuf [8]*colors.Color
			protos := protoBuf[:0]
			var protoByHash map[int32][]*colors.Color
			for i := 0; i < len(prototypeIds); i++ {
				protos = appendUniquePoolColor(protos, b.resolveSingleShardPointer(shard, int(prototypeIds[i])), &protoByHash)
			}
			prototypes = b.allocColorSlice(len(protos))
			copy(prototypes, protos)
		}

		ownProperties := colors.EmptyStringSlice()
		if ownPropList := objProto.GetOwnPropertyList(); len(ownPropList) == 1 {
			prop := shard.stringPool.Get(int(ownPropList[0]))
			ownProperties = b.allocStringSlice(1)
			ownProperties[0] = prop
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
			ownProperties = b.allocStringSlice(len(b.propScratch))
			copy(ownProperties, b.propScratch)
		}
		result = b.colorBuilder.
			ClearBoxId().
			SetId(id).
			SetInstanceColors(instanceColors).
			SetPrototypes(prototypes).
			SetOwnProperties(ownProperties).
			SetUnionElements(colors.EmptyColorSlice()).
			SetClosureAssert(objProto.GetClosureAssert()).
			SetConstructor(objProto.GetMarkedConstructor()).
			SetInvalidating(objProto.GetIsInvalidating()).
			SetPropertiesKeepOriginalName(objProto.GetPropertiesKeepOriginalName()).
			Build()
	} else if proto.HasUnion() {
		unionMemberList := proto.GetUnion().GetUnionMemberList()
		var unionBuf [8]*colors.Color
		union := unionBuf[:0]
		for i := 0; i < len(unionMemberList); i++ {
			member := b.resolveSingleShardPointer(shard, int(unionMemberList[i]))
			MalformedTypedAstException_CheckWellFormedWithParam(!member.IsUnion(), "Reconciling union with non-union", proto)
			union = append(union, member)
		}
		result = colors.Color_CreateUnionSliceWithBuilder(union, &b.colorBuilder, b.allocColorSliceFn, b.allocStringSliceFn)
		if id != result.GetId() {
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
		if existing == colorPool_PENDING_COLOR {
			rows := make([]string, 0, len(b.reconcilationDebugStack))
			for _, stackId := range b.reconcilationDebugStack {
				rows = append(rows, b.protoEntries[b.idToProtoIdx[stackId]].rowString())
			}
			panic(NewMalformedTypedAstException(
				"Cyclic Color structure detected: [" + strings.Join(rows, ", ") + "]"))
		}
		return existing
	}
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
	var instBuf [8]*colors.Color
	inst := instBuf[:0]
	var instByHash map[int32][]*colors.Color

	var protoBuf [8]*colors.Color
	protos := protoBuf[:0]
	var protoByHash map[int32][]*colors.Color

	isClosureAssert := base.Tri_UNKNOWN
	isConstructor := false
	isInvalidating := false
	propertiesKeepOriginalName := false
	totalOwnProps := 0

	if entry.firstProto != nil {
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
			for i := 0; i < len(instanceTypeIds); i++ {
				instanceTypePointer := instanceTypeIds[i]
				inst = appendUniquePoolColor(inst, b.LookupOrReconcileColor(shard.GetId(int(instanceTypePointer))), &instByHash)
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
				protos = appendUniquePoolColor(protos, b.LookupOrReconcileColor(shard.GetId(int(prototypePointer))), &protoByHash)
			}
			totalOwnProps += objProto.GetOwnPropertyCount()
		}
	}

	if seen := entry.jsRecord; seen != nil {
		jutil.CheckState(len(seen.unionMembers) == 0)
		s := b.jsTypeShard
		jsClosureAssert := false
		for _, type_ := range seen.jstypes {
			objType := jutil.CheckNotNullMsg1(type_.ToMaybeObjectType(), "%s", type_)
			if objType.IsFunctionType() {
				fnType := objType.ToMaybeFunctionType()
				if fnType.HasInstanceType() && fnType.GetInstanceType() != nil {
					instPtr := s.SerializeType(fnType.GetInstanceType())
					inst = appendUniquePoolColor(inst, b.LookupOrReconcileColor(s.recordsByPointer[instPtr].colorId), &instByHash)
					protoPtr := s.SerializeType(fnType.GetPrototype())
					protos = appendUniquePoolColor(protos, b.LookupOrReconcileColor(s.recordsByPointer[protoPtr].colorId), &protoByHash)
					isConstructor = isConstructor || fnType.IsConstructor()
				}
				jsClosureAssert = jsClosureAssert || JSTypeReconserializer_IsClosureAssert(fnType.GetClosurePrimitive())
			}
			isInvalidating = isInvalidating || s.invalidatingTypes.IsInvalidating(objType)
			propertiesKeepOriginalName = propertiesKeepOriginalName || objType.IsEnumType()
		}
		MalformedTypedAstException_CheckWellFormedWithParam(
			isClosureAssert.ToBoolean(jsClosureAssert) == jsClosureAssert,
			"Inconsistent values for closure_assert",
			seen)
		isClosureAssert = base.Tri_ForBoolean(jsClosureAssert)
	}

	instanceColors := colors.EmptyColorSlice()
	if len(inst) > 0 {
		instanceColors = b.allocColorSlice(len(inst))
		copy(instanceColors, inst)
	}

	prototypes := colors.EmptyColorSlice()
	if len(protos) > 0 {
		prototypes = b.allocColorSlice(len(protos))
		copy(prototypes, protos)
	}

	ownProperties := colors.EmptyStringSlice()
	if entry.jsRecord == nil {
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
					prop := shard.stringPool.Get(int(objProto.GetOwnProperty(0)))
					ownProperties = b.allocStringSlice(1)
					ownProperties[0] = prop
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
			ownProperties = b.allocStringSlice(len(b.propScratch))
			copy(ownProperties, b.propScratch)
		}
	} else if entry.firstProto == nil && len(entry.jsRecord.jstypes) == 1 {
		b.propScratch = b.propScratch[:0]
		entry.jsRecord.jstypes[0].ToMaybeObjectType().ForEachOwnPropertyName(b.directSingleOwnPropFn)
		if len(b.propScratch) > 0 {
			ownProperties = b.allocStringSlice(len(b.propScratch))
			copy(ownProperties, b.propScratch)
		}
	} else {
		b.startDirectPropDedup()
		if entry.firstProto != nil {
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
					b.appendDirectDedupProp(shard.stringPool.Get(int(objProto.GetOwnProperty(i))))
				}
			}
		}
		for _, type_ := range entry.jsRecord.jstypes {
			type_.ToMaybeObjectType().ForEachOwnPropertyName(b.directMultiOwnPropFn)
		}
		if len(b.propScratch) > 0 {
			ownProperties = b.allocStringSlice(len(b.propScratch))
			copy(ownProperties, b.propScratch)
		}
	}

	return b.colorBuilder.
		ClearBoxId().
		SetId(id).
		SetInstanceColors(instanceColors).
		SetPrototypes(prototypes).
		SetOwnProperties(ownProperties).
		SetUnionElements(colors.EmptyColorSlice()).
		SetClosureAssert(isClosureAssert.ToBoolean(false)).
		SetConstructor(isConstructor).
		SetInvalidating(isInvalidating).
		SetPropertiesKeepOriginalName(propertiesKeepOriginalName).
		Build()
}

func (b *ColorPool_Builder) reconcileUnionProtos(id colors.ColorId, entry colorPoolIdEntry) *colors.Color {
	var unionBuf [16]*colors.Color
	union := unionBuf[:0]
	if entry.firstProto != nil {
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
				union = append(union, member)
			}
		}
	}
	if seen := entry.jsRecord; seen != nil {
		jutil.CheckStateMsg1(len(seen.unionMembers) > 1, "%s", seen)
		records := b.jsTypeShard.recordsByPointer
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
			member := b.LookupOrReconcileColor(records[memberPtr].colorId)
			MalformedTypedAstException_CheckWellFormedWithParam(!member.IsUnion(), "Reconciling union with non-union", seen)
			union = append(union, member)
		}
	}

	result := colors.Color_CreateUnionSliceWithBuilder(union, &b.colorBuilder, b.allocColorSliceFn, b.allocStringSliceFn)
	if id != result.GetId() {
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
		var scratchBuf [16]colors.ColorId
		scratch := scratchBuf[:0]
		var seenMap map[colors.ColorId]uint32
		var gen uint32
		for i, proto := range typeList {
			if proto.HasUnion() {
				out[i] = colorPool_createUnionColorId(proto, i, typeList, out, &scratch, &seenMap, &gen)
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
	scratch *[]colors.ColorId,
	seenMap *map[colors.ColorId]uint32,
	gen *uint32,
) colors.ColorId {
	unionMemberList := proto.GetUnion().GetUnionMemberList()
	MalformedTypedAstException_CheckWellFormedWithParam(len(unionMemberList) > 1, "Union has too few members", proto)
	members := (*scratch)[:0]
	useMap := len(unionMemberList) > 16
	var curGen uint32
	if useMap {
		*gen++
		if *seenMap == nil {
			*seenMap = make(map[colors.ColorId]uint32, 32)
		} else if *gen == 0 {
			clear(*seenMap)
			*gen = 1
		}
		curGen = *gen
	}
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
		if useMap {
			if (*seenMap)[cid] != curGen {
				(*seenMap)[cid] = curGen
				members = append(members, cid)
			}
		} else if !slices.Contains(members, cid) {
			members = append(members, cid)
		}
	}
	*scratch = members
	return unionColorIds(members)
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
