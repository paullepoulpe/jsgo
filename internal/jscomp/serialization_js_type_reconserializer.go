// Port of com/google/javascript/jscomp/serialization/JSTypeReconserializer.java @ closure-compiler v20261015

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

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// JSTypeReconserializer takes JSTypes produced by JSCompiler's typechecker
// and deduplicates and serializes them into the TypedAST colors format.
//
// The deduplication phase is called "reconciliation". It's necessary
// because the the TypedAST colors format is simpler than the JSType format,
// and so there is a many-to-one relationship between JSTypes and TypedAST
// colors.
type JSTypeReconserializer struct {
	registry                    *rhino.JSTypeRegistry
	serializationMode           *SerializationOptions
	invalidatingTypes           *InvalidatingTypes
	stringPoolBuilder           *StringPool_Builder
	hasher                      *JSTypeColorIdHasher
	shouldPropagatePropertyName func(string) bool

	// Cache some commonly used types.
	unknownRecord     *JSTypeReconserializer_SeenTypeRecord
	topFunctionRecord *JSTypeReconserializer_SeenTypeRecord

	typeToRecordCache *base.LinkedIdentityHashMap[rhino.JSType, *JSTypeReconserializer_SeenTypeRecord]
	seenTypeRecords   *jutil.LinkedHashMap[colors.ColorId, *JSTypeReconserializer_SeenTypeRecord]
	// Go-only: language-adapter — LinkedHashMultimap<Integer, Integer>.
	disambiguateEdges *jutil.LinkedHashMap[int, *jutil.LinkedHashSet[int]]

	scratchInstances  *jutil.LinkedHashSet[int]
	scratchPrototypes *jutil.LinkedHashSet[int]
	scratchOwnProps   *jutil.LinkedHashSet[int]

	state JSTypeReconserializer_State
}

// This is a one-way mapping because some JSTypes go to the same Color.
//
// Java: JSTypeReconserializer.JSTYPE_NATIVE_TO_AXIOMATIC_COLOR_MAP
var jsTypeReconserializer_JSTYPE_NATIVE_TO_AXIOMATIC_COLOR_MAP = jutil.NewImmutableMapBuilder[rhino.JSTypeNative, *colors.Color]().
	// Merge all the various top/bottom-like/unknown types into a single unknown type.
	Put(rhino.JSTypeNative_ALL_TYPE, colors.StandardColors_UNKNOWN).
	Put(rhino.JSTypeNative_CHECKED_UNKNOWN_TYPE, colors.StandardColors_UNKNOWN).
	Put(rhino.JSTypeNative_NO_OBJECT_TYPE, colors.StandardColors_UNKNOWN).
	Put(rhino.JSTypeNative_NO_TYPE, colors.StandardColors_UNKNOWN).
	Put(rhino.JSTypeNative_UNKNOWN_TYPE, colors.StandardColors_UNKNOWN).
	// Map all the primitives in the obvious way.
	Put(rhino.JSTypeNative_BIGINT_TYPE, colors.StandardColors_BIGINT).
	Put(rhino.JSTypeNative_BOOLEAN_TYPE, colors.StandardColors_BOOLEAN).
	Put(rhino.JSTypeNative_GBIGINT_TYPE, colors.StandardColors_GBIGINT).
	Put(rhino.JSTypeNative_NULL_TYPE, colors.StandardColors_NULL_OR_VOID).
	Put(rhino.JSTypeNative_NUMBER_TYPE, colors.StandardColors_NUMBER).
	Put(rhino.JSTypeNative_STRING_TYPE, colors.StandardColors_STRING).
	Put(rhino.JSTypeNative_SYMBOL_TYPE, colors.StandardColors_SYMBOL).
	Put(rhino.JSTypeNative_VOID_TYPE, colors.StandardColors_NULL_OR_VOID).
	// Smoosh top-like function types into a single type.
	Put(rhino.JSTypeNative_FUNCTION_FUNCTION_TYPE, colors.StandardColors_TOP_FUNCTION).
	Put(rhino.JSTypeNative_FUNCTION_TYPE, colors.StandardColors_TOP_FUNCTION).
	// Smoosh top-like objects into a single type.
	Put(rhino.JSTypeNative_FUNCTION_PROTOTYPE, colors.StandardColors_TOP_OBJECT).
	Put(rhino.JSTypeNative_FUNCTION_INSTANCE_PROTOTYPE, colors.StandardColors_TOP_OBJECT).
	Put(rhino.JSTypeNative_OBJECT_FUNCTION_TYPE, colors.StandardColors_TOP_OBJECT).
	Put(rhino.JSTypeNative_OBJECT_PROTOTYPE, colors.StandardColors_TOP_OBJECT).
	Put(rhino.JSTypeNative_OBJECT_TYPE, colors.StandardColors_TOP_OBJECT).
	BuildOrThrow()

// JSTypeReconserializer_State is JSTypeReconserializer.State.
type JSTypeReconserializer_State int32

const (
	JSTypeReconserializer_State_COLLECTING_TYPES JSTypeReconserializer_State = iota
	JSTypeReconserializer_State_GENERATING_POOL
	JSTypeReconserializer_State_FINISHED
)

var jsTypeReconserializer_State_names = [...]string{"COLLECTING_TYPES", "GENERATING_POOL", "FINISHED"}

// JSTypeReconserializer_State_values is State.values().
func JSTypeReconserializer_State_values() []JSTypeReconserializer_State {
	return []JSTypeReconserializer_State{
		JSTypeReconserializer_State_COLLECTING_TYPES,
		JSTypeReconserializer_State_GENERATING_POOL,
		JSTypeReconserializer_State_FINISHED,
	}
}

// Name is Enum.name().
func (s JSTypeReconserializer_State) Name() string { return jsTypeReconserializer_State_names[s] }

// String is Enum.toString().
func (s JSTypeReconserializer_State) String() string { return s.Name() }

// Ordinal is Enum.ordinal().
func (s JSTypeReconserializer_State) Ordinal() int { return int(s) }

// Java: JSTypeReconserializer(JSTypeRegistry, InvalidatingTypes, StringPool.Builder, Predicate<String>, SerializationOptions)
func newJSTypeReconserializer(
	registry *rhino.JSTypeRegistry,
	invalidatingTypes *InvalidatingTypes,
	stringPoolBuilder *StringPool_Builder,
	shouldPropagatePropertyName func(string) bool,
	serializationMode *SerializationOptions) *JSTypeReconserializer {
	s := &JSTypeReconserializer{}
	// Field initializers.
	s.typeToRecordCache = base.NewLinkedIdentityHashMap[rhino.JSType, *JSTypeReconserializer_SeenTypeRecord]()
	s.seenTypeRecords = jutil.NewLinkedHashMap[colors.ColorId, *JSTypeReconserializer_SeenTypeRecord]()
	s.disambiguateEdges = jutil.NewLinkedHashMap[int, *jutil.LinkedHashSet[int]]()
	s.state = JSTypeReconserializer_State_COLLECTING_TYPES
	// Constructor body.
	s.registry = registry
	s.hasher = NewJSTypeColorIdHasher(registry)
	s.invalidatingTypes = invalidatingTypes
	s.stringPoolBuilder = stringPoolBuilder
	s.shouldPropagatePropertyName = shouldPropagatePropertyName
	s.serializationMode = serializationMode

	s.SeedCachesWithAxiomaticTypes()

	s.unknownRecord = s.seenTypeRecords.Get(colors.StandardColors_UNKNOWN.GetId())
	s.topFunctionRecord = s.seenTypeRecords.Get(colors.StandardColors_TOP_FUNCTION.GetId())
	return s
}

func (s *JSTypeReconserializer) putDisambiguateEdge(subtype, supertype int) {
	set := s.disambiguateEdges.Get(subtype)
	if set == nil {
		set = jutil.NewLinkedHashSet[int]()
		s.disambiguateEdges.PutNew(subtype, set)
	}
	set.Add(supertype)
}

// JSTypeReconserializer_Create initializes a JSTypeReconserializer
//
// shouldPropagatePropertyName decides whether a property present on some
// ObjectType should actually be serialized. Used to avoid serializing
// properties that won't impact optimizations (because they aren't present
// in the AST)
//
// Java: JSTypeReconserializer.create(JSTypeRegistry, InvalidatingTypes, StringPool.Builder, Predicate<String>, SerializationOptions)
func JSTypeReconserializer_Create(
	registry *rhino.JSTypeRegistry,
	invalidatingTypes *InvalidatingTypes,
	stringPoolBuilder *StringPool_Builder,
	shouldPropagatePropertyName func(string) bool,
	serializationMode *SerializationOptions) *JSTypeReconserializer {
	serializer := newJSTypeReconserializer(
		registry,
		invalidatingTypes,
		stringPoolBuilder,
		shouldPropagatePropertyName,
		serializationMode)
	serializer.CheckValidLinearTime()
	return serializer
}

// SerializeType returns a pointer to the given type. If it is not already
// serialized, serializes it too
//
// Java: JSTypeReconserializer.serializeType(JSType)
func (s *JSTypeReconserializer) SerializeType(type_ rhino.JSType) int {
	record := s.RecordType(type_)
	return record.pointer
}

// Java: JSTypeReconserializer.recordType(JSType)
func (s *JSTypeReconserializer) RecordType(type_ rhino.JSType) *JSTypeReconserializer_SeenTypeRecord {
	if jstypeRecord := s.typeToRecordCache.Get(type_); jstypeRecord != nil {
		return jstypeRecord
	}

	var forwardedType rhino.JSType
	if type_.IsNamedType() {
		forwardedType = type_.ToMaybeNamedType().GetReferencedType()
	} else if type_.IsEnumElementType() {
		forwardedType = type_.ToMaybeEnumElementType().GetPrimitiveType()
	} else if type_.IsKnownSymbolValueType() {
		forwardedType = s.registry.GetNativeType(rhino.JSTypeNative_SYMBOL_TYPE)
	} else if type_.IsTemplatizedType() {
		forwardedType = type_.ToMaybeTemplatizedType().GetReferencedType()
	} else if type_.IsFunctionType() &&
		type_.ToMaybeFunctionType().GetCanonicalRepresentation() != nil {
		forwardedType = type_.ToMaybeFunctionType().GetCanonicalRepresentation()
	} else {
		forwardedType = nil
	}
	if forwardedType != nil {
		return s.RecordType(forwardedType)
	}

	if type_.IsUnknownType() || type_.IsNoResolvedType() || type_.IsTemplateType() {
		// template types are not serialized because optimizations don't seem to care about them.
		// serialize as the UNKNOWN_TYPE because bounded generics are unsupported
		return s.unknownRecord
	}
	if type_.IsFunctionType() &&
		!type_.ToMaybeFunctionType().HasInstanceType() &&
		type_.ToMaybeFunctionType().GetClosurePrimitive() == nil {
		// Distinguishing different function types does not matter for optimizations unless they
		// are a constructor/interface or have a @closurePrimitive tag associated. Optimization colors
		// don't track function parameter/return/template types, and function literals are all
		// invalidating types during property disambiguation.
		return s.topFunctionRecord
	}

	if type_.IsUnionType() {
		return s.RecordUnionType(type_.ToMaybeUnionType())
	} else if type_.IsObjectType() {
		return s.RecordObjectType(type_.ToMaybeObjectType())
	}

	panic(jutil.NewAssertionError(type_.String()))
}

// Java: JSTypeReconserializer.recordUnionType(UnionType)
func (s *JSTypeReconserializer) RecordUnionType(type_ *rhino.UnionType) *JSTypeReconserializer_SeenTypeRecord {
	jutil.CheckNotNull(type_)

	altRecords := jutil.NewLinkedHashSet[*JSTypeReconserializer_SeenTypeRecord]()
	for _, altType := range type_.GetAlternates() {
		alt := s.RecordType(altType)
		if alt.unionMembers == nil {
			altRecords.Add(alt)
		} else {
			// Flatten out any nested unions. They are possible due to proxy-like types.
			for member := range alt.unionMembers.All() {
				altRecords.Add(member)
			}
		}
	}

	// Some elements of the union may be equal as Colors
	if altRecords.Size() == 1 {
		return jsTypeReconserializer_getOnlyElement(altRecords)
	}

	alternateIds := jutil.NewLinkedHashSet[colors.ColorId]()
	for altRecord := range altRecords.All() {
		alternateIds.Add(altRecord.colorId)
	}
	unionId := colors.ColorId_Union(alternateIds)
	record := s.GetOrCreateRecord(unionId, type_)

	if record.unionMembers == nil {
		record.unionMembers = altRecords
	} else if s.serializationMode.RunValidation() {
		if !altRecords.Equals(record.unionMembers) {
			jutil.CheckStateMsg3(
				false,
				"Unions with same ID must have same members: %s => %s == %s",
				unionId,
				jsTypeReconserializer_colorIdsOf(altRecords),
				jsTypeReconserializer_colorIdsOf(record.unionMembers))
		}
	}

	return record
}

// jsTypeReconserializer_getOnlyElement is Iterables.getOnlyElement on a set
// of size 1.
func jsTypeReconserializer_getOnlyElement(
	set *jutil.LinkedHashSet[*JSTypeReconserializer_SeenTypeRecord]) *JSTypeReconserializer_SeenTypeRecord {
	for e := range set.All() {
		return e
	}
	panic(jutil.NewNoSuchElementException())
}

// jsTypeReconserializer_colorIdsOf is records.stream().map((r) ->
// r.colorId).collect(toImmutableSet()).
func jsTypeReconserializer_colorIdsOf(
	records *jutil.LinkedHashSet[*JSTypeReconserializer_SeenTypeRecord]) *jutil.LinkedHashSet[colors.ColorId] {
	ids := jutil.NewLinkedHashSet[colors.ColorId]()
	for r := range records.All() {
		ids.Add(r.colorId)
	}
	return ids
}

// Java: JSTypeReconserializer.recordObjectType(ObjectType)
func (s *JSTypeReconserializer) RecordObjectType(type_ rhino.ObjectType) *JSTypeReconserializer_SeenTypeRecord {
	jutil.CheckNotNull(type_)

	id := s.hasher.HashObjectType(type_)
	record := s.GetOrCreateRecord(id, type_)
	s.AddSupertypeEdges(type_, record.pointer)

	if type_.IsFunctionType() {
		fnType := type_.ToMaybeFunctionType()
		if fnType.HasInstanceType() && fnType.GetInstanceType() != nil {
			// We have to serialize these here ahead of time to avoid a ConcurrentModificationException
			// during reconciliation.
			s.SerializeType(fnType.GetInstanceType())
			s.SerializeType(fnType.GetPrototype())
		}
	}

	return record
}

// Java: JSTypeReconserializer.addSupertypeEdges(ObjectType, Integer)
func (s *JSTypeReconserializer) AddSupertypeEdges(subtype rhino.ObjectType, serializedSubtype int) {
	// disambiguateEdges.putAll(serializedSubtype, ownAncestorInterfacesOf(subtype))
	if ctorType := subtype.GetConstructor(); ctorType != nil {
		ext := ctorType.GetExtendedInterfaces()
		imp := ctorType.GetOwnImplementedInterfaces()
		if len(ext)+len(imp) > 0 {
			var buf [8]int
			ancestors := buf[:0]
			for _, ancestor := range ext {
				ancestors = append(ancestors, s.SerializeType(ancestor))
			}
			for _, ancestor := range imp {
				ancestors = append(ancestors, s.SerializeType(ancestor))
			}
			for _, ancestor := range ancestors {
				s.putDisambiguateEdge(serializedSubtype, ancestor)
			}
		}
	}
	if implicitProto := subtype.GetImplicitPrototype(); implicitProto != nil {
		supertype := s.SerializeType(implicitProto)
		s.putDisambiguateEdge(serializedSubtype, supertype)
	}
}

// Java: JSTypeReconserializer.getOrCreateRecord(ColorId, JSType)
func (s *JSTypeReconserializer) GetOrCreateRecord(id colors.ColorId, jstype rhino.JSType) *JSTypeReconserializer_SeenTypeRecord {
	jutil.CheckNotNull(jstype)
	jutil.CheckState(JSTypeReconserializer_State_COLLECTING_TYPES == s.state || JSTypeReconserializer_State_GENERATING_POOL == s.state)

	record := s.seenTypeRecords.Get(id)
	if record == nil {
		pointer := s.seenTypeRecords.Size()
		record = newJSTypeReconserializer_SeenTypeRecord(id, pointer)
		s.seenTypeRecords.PutNew(id, record)
	}
	s.typeToRecordCache.Put(jstype, record)
	record.jstypes = append(record.jstypes, jstype)

	return record
}

// Java: JSTypeReconserializer.reconcileUnionTypes(SeenTypeRecord)
func (s *JSTypeReconserializer) ReconcileUnionTypes(seen *JSTypeReconserializer_SeenTypeRecord) *TypeProto {
	members := make([]int32, 0, seen.unionMembers.Size())
	for r := range seen.unionMembers.All() {
		members = append(members, int32(r.pointer))
	}
	slices.Sort(members)
	unionProto := &UnionTypeProto{
		unionMember_: members,
		frozen:       true,
		memoizedSize: -1,
	}
	return &TypeProto{
		union_:       unionProto,
		kindCase_:    TypeProto_UNION_FIELD_NUMBER,
		frozen:       true,
		memoizedSize: -1,
	}
}

// Java: JSTypeReconserializer.reconcileObjectTypes(SeenTypeRecord)
func (s *JSTypeReconserializer) ReconcileObjectTypes(seen *JSTypeReconserializer_SeenTypeRecord) *TypeProto {
	if len(seen.jstypes) == 1 {
		type_ := seen.jstypes[0]
		objType := jutil.CheckNotNullMsg1(type_.ToMaybeObjectType(), "%s", type_)

		var instanceTypes []int32
		var prototypes []int32
		isClosureAssert := false
		isConstructor := false

		if objType.IsFunctionType() {
			fnType := objType.ToMaybeFunctionType()
			if fnType.HasInstanceType() && fnType.GetInstanceType() != nil {
				instanceTypes = []int32{int32(s.SerializeType(fnType.GetInstanceType()))}
				prototypes = []int32{int32(s.SerializeType(fnType.GetPrototype()))}
				isConstructor = fnType.IsConstructor()
			}
			isClosureAssert = JSTypeReconserializer_IsClosureAssert(fnType.GetClosurePrimitive())
		}

		var ownProps []int32
		for _, ownProperty := range objType.GetOwnPropertyNames() {
			if s.shouldPropagatePropertyName(ownProperty) {
				ownProps = append(ownProps, int32(s.stringPoolBuilder.Put(ownProperty)))
			}
		}

		isInvalidating := s.invalidatingTypes.IsInvalidating(objType)
		propertiesKeepOriginalName := objType.IsEnumType()

		objectProto := &ObjectTypeProto{
			instanceType_:               instanceTypes,
			ownProperty_:                ownProps,
			prototype_:                  prototypes,
			closureAssert_:              isClosureAssert,
			isInvalidating_:             isInvalidating,
			markedConstructor_:          isConstructor,
			propertiesKeepOriginalName_: propertiesKeepOriginalName,
			uuid_:                       seen.colorId.AsByteString(),
			frozen:                      true,
			memoizedSize:                -1,
		}
		return &TypeProto{
			object_:      objectProto,
			kindCase_:    TypeProto_OBJECT_FIELD_NUMBER,
			frozen:       true,
			memoizedSize: -1,
		}
	}

	if s.scratchInstances == nil {
		s.scratchInstances = jutil.NewLinkedHashSet[int]()
		s.scratchPrototypes = jutil.NewLinkedHashSet[int]()
		s.scratchOwnProps = jutil.NewLinkedHashSet[int]()
	} else {
		s.scratchInstances.Clear()
		s.scratchPrototypes.Clear()
		s.scratchOwnProps.Clear()
	}
	instancePointers := s.scratchInstances
	prototypePointers := s.scratchPrototypes
	ownProperties := s.scratchOwnProps
	isClosureAssert := false
	isConstructor := false
	isInvalidating := false
	propertiesKeepOriginalName := false

	for _, type_ := range seen.jstypes {
		objType := jutil.CheckNotNullMsg1(type_.ToMaybeObjectType(), "%s", type_)

		if objType.IsFunctionType() {
			fnType := objType.ToMaybeFunctionType()

			// Serialize prototypes and instance types for instantiable types. Even if these types never
			// appear on the AST, optimizations need to know that at runtime these types may be present.
			if fnType.HasInstanceType() && fnType.GetInstanceType() != nil {
				instancePointers.Add(s.SerializeType(fnType.GetInstanceType()))
				prototypePointers.Add(s.SerializeType(fnType.GetPrototype()))
				isConstructor = isConstructor || fnType.IsConstructor()
			}

			isClosureAssert = isClosureAssert || JSTypeReconserializer_IsClosureAssert(fnType.GetClosurePrimitive())
		}

		for _, ownProperty := range objType.GetOwnPropertyNames() {
			// TODO(b/169899789): consider omitting common, well-known properties like "prototype" to
			// save space.
			if s.shouldPropagatePropertyName(ownProperty) {
				ownProperties.Add(s.stringPoolBuilder.Put(ownProperty))
			}
		}

		isInvalidating = isInvalidating || s.invalidatingTypes.IsInvalidating(objType)

		/*
		 * To support legacy code, property disambiguation never renames properties of enums (e.g. 'A'
		 * in '/** @enum * / const E = {A: 0}`). In theory this would be safe to remove if we clean up
		 * code depending on the lack of renaming
		 */
		propertiesKeepOriginalName = propertiesKeepOriginalName || objType.IsEnumType()
	}

	objectProto := &ObjectTypeProto{
		instanceType_:               jsTypeReconserializer_int32s(instancePointers),
		ownProperty_:                jsTypeReconserializer_int32s(ownProperties),
		prototype_:                  jsTypeReconserializer_int32s(prototypePointers),
		closureAssert_:              isClosureAssert,
		isInvalidating_:             isInvalidating,
		markedConstructor_:          isConstructor,
		propertiesKeepOriginalName_: propertiesKeepOriginalName,
		uuid_:                       seen.colorId.AsByteString(),
		frozen:                      true,
		memoizedSize:                -1,
	}
	return &TypeProto{
		object_:      objectProto,
		kindCase_:    TypeProto_OBJECT_FIELD_NUMBER,
		frozen:       true,
		memoizedSize: -1,
	}
}

// jsTypeReconserializer_int32s is the Iterable<Integer> handed to a proto's
// addAll: the set's elements in order, as int32.
func jsTypeReconserializer_int32s(set *jutil.LinkedHashSet[int]) []int32 {
	if set.IsEmpty() {
		return nil
	}
	out := make([]int32, 0, set.Size())
	for v := range set.All() {
		out = append(out, int32(v))
	}
	return out
}

// OwnAncestorInterfacesOf returns the interfaces directly implemented and
// extended by type.
//
// Some of these relationships represent type errors; however, the graph
// needs to contain those edges for safe disambiguation. In particular, code
// generated from other languages (e.g TS) might have more flexible subtyping
// rules.
//
// Java: JSTypeReconserializer.ownAncestorInterfacesOf(ObjectType)
func (s *JSTypeReconserializer) OwnAncestorInterfacesOf(type_ rhino.ObjectType) []int {
	ctorType := type_.GetConstructor()
	if ctorType == nil {
		return nil
	}
	ext := ctorType.GetExtendedInterfaces()
	imp := ctorType.GetOwnImplementedInterfaces()
	if len(ext)+len(imp) == 0 {
		return nil
	}

	ancestors := make([]int, 0, len(ext)+len(imp))
	for _, ancestor := range ext {
		ancestors = append(ancestors, s.SerializeType(ancestor))
	}
	for _, ancestor := range imp {
		ancestors = append(ancestors, s.SerializeType(ancestor))
	}
	return ancestors
}

// SeedCachesWithAxiomaticTypes inserts dummy pointers corresponding to all
// PrimitiveTypes in the type pool.
//
// These types will never correspond to an actual TypeProto. Instead, all
// normal Integer offsets into the pool are offset by a number equivalent to
// the number of PrimitiveType enum elements.
//
// Java: JSTypeReconserializer.seedCachesWithAxiomaticTypes()
func (s *JSTypeReconserializer) SeedCachesWithAxiomaticTypes() {
	jutil.CheckState(s.seenTypeRecords.IsEmpty())

	// Load all the axiomatic records in the right order without any types.
	for _, axiomatic := range TypePointers_OFFSET_TO_AXIOMATIC_COLOR {
		index := s.seenTypeRecords.Size()
		record := newJSTypeReconserializer_SeenTypeRecord(axiomatic.GetId(), index)
		s.seenTypeRecords.Put(axiomatic.GetId(), record)
	}

	// Add JSTypes corresponding to axiomatic IDs.
	for jstypeNative, axiomatic := range jsTypeReconserializer_JSTYPE_NATIVE_TO_AXIOMATIC_COLOR_MAP.All() {
		s.GetOrCreateRecord(axiomatic.GetId(), s.registry.GetNativeType(jstypeNative))
	}

	jutil.CheckState(s.seenTypeRecords.Size() == len(TypePointers_OFFSET_TO_AXIOMATIC_COLOR))
}

// CheckValidLinearTime checks that this instance is in a valid state.
//
// Java: JSTypeReconserializer.checkValidLinearTime()
func (s *JSTypeReconserializer) CheckValidLinearTime() {
	if !s.serializationMode.RunValidation() {
		return
	}

	totalTypeCount := s.seenTypeRecords.Size()
	for seen := range s.seenTypeRecords.Values() {
		offset := seen.pointer
		jutil.CheckState(offset >= 0)
		jutil.CheckStateMsg2(
			offset <= totalTypeCount,
			"Found invalid pointer %s, out of a total of %s user-defined types",
			offset,
			totalTypeCount)
	}
}

// GenerateTypePool generates a "type-pool" representing all the types that
// this class has encountered through calls to SerializeType.
//
// After generation, no new types can be added, so subsequent calls to
// SerializeType can only be used to retrieve pointers to existing types in
// the type pool.
//
// Java: JSTypeReconserializer.generateTypePool()
func (s *JSTypeReconserializer) GenerateTypePool() *TypePool {
	jutil.CheckState(s.state == JSTypeReconserializer_State_COLLECTING_TYPES)
	s.CheckValidLinearTime()

	var debugInfoProto *TypePool_DebugInfo

	if s.serializationMode.IncludeDebugInfo() {
		debugInfo := TypePool_DebugInfo_NewBuilder()
		// Key by source ref to deduplicate the strings, which are pretty long.
		for location, types := range s.invalidatingTypes.
			GetMismatchLocations().
			Inverse().
			AsMap().
			All() {
			mismatch := debugInfo.
				AddMismatchBuilder().
				SetSourceRef(location.GetLocation())
			// types.stream().peek(...).map(this::serializeType).distinct().sorted()
			involved := jutil.NewLinkedHashSet[int]()
			for t := range types.All() {
				jutil.CheckStateMsg1(!t.IsUnionType(), "%s", t)
				// Ensure all types are recorded before reconciliation.
				involved.Add(s.SerializeType(t))
			}
			involvedColors := jsTypeReconserializer_int32s(involved)
			slices.Sort(involvedColors)
			mismatch.AddAllInvolvedColor(involvedColors)
		}
		debugInfoProto = debugInfo.Build()
	}

	s.state = JSTypeReconserializer_State_GENERATING_POOL

	// Go-only: language-adapter — Java iterates seenTypeRecords.values(),
	// which would throw ConcurrentModificationException if reconciliation
	// recorded a new type; recordObjectType serializes instance types and
	// prototypes ahead of time so that it never does.
	numTypes := s.seenTypeRecords.Size() - len(TypePointers_OFFSET_TO_AXIOMATIC_COLOR)
	var types []*TypeProto
	if numTypes > 0 {
		types = make([]*TypeProto, 0, numTypes)
	}
	for seen := range s.seenTypeRecords.Values() {
		if colors.StandardColors_AXIOMATIC_COLORS.ContainsKey(seen.colorId) {
			jutil.CheckStateMsg1(TypePointers_IsAxiomatic(seen.pointer), "Missing .type for SeenTypeRecord %s", seen)
			continue
		}
		if seen.unionMembers != nil {
			types = append(types, s.ReconcileUnionTypes(seen))
		} else {
			types = append(types, s.ReconcileObjectTypes(seen))
		}
	}

	var disambiguationEdges []*SubtypingEdge
	for subtype, supertypes := range s.disambiguateEdges.All() {
		subtypeI32 := int32(subtype)
		for supertype := range supertypes.All() {
			disambiguationEdges = append(disambiguationEdges, &SubtypingEdge{
				subtype_:     subtypeI32,
				supertype_:   int32(supertype),
				frozen:       true,
				memoizedSize: -1,
			})
		}
	}

	s.state = JSTypeReconserializer_State_FINISHED
	s.CheckValidLinearTime()
	return &TypePool{
		type_:                types,
		disambiguationEdges_: disambiguationEdges,
		debugInfo_:           debugInfoProto,
		frozen:               true,
		memoizedSize:         -1,
	}
}

// GetColorIdToJSTypeMapForDebugging returns a map from
// ObjectTypeProto#getUuid() to the originating JSTypes.
//
// Only intended to be used for debug logging.
//
// Java: JSTypeReconserializer.getColorIdToJSTypeMapForDebugging()
func (s *JSTypeReconserializer) GetColorIdToJSTypeMapForDebugging() *jutil.ListMultimap[string, rhino.JSType] {
	// note: returns JSType values instead of String values, even though all that's needed for
	// debugging are the strings, to avoid the memory overhead of calculating all string
	// representations at once
	colorIdToTypes := jutil.NewImmutableListMultimapBuilder[string, rhino.JSType]().OrderKeysBy(jutil.StringCompareTo)
	s.typeToRecordCache.ForEach(func(jstype rhino.JSType, record *JSTypeReconserializer_SeenTypeRecord) {
		colorIdToTypes.Put(record.colorId.String(), jstype)
	})
	return colorIdToTypes.Build()
}

// JSTypeReconserializer_SeenTypeRecord is JSTypeReconserializer.SeenTypeRecord.
type JSTypeReconserializer_SeenTypeRecord struct {
	colorId colors.ColorId
	pointer int

	// 2021-05-25: It's faster to build a list and reconcile duplicates than
	// to deduplicate using a set. The likely cause is that the caches head of
	// these lists have low hit-rates for unions, and that union
	// reconciliation doesn't actually look at most entries in this list.
	jstypes   []rhino.JSType
	jstypeBuf [2]rhino.JSType

	// Java: @Nullable ImmutableSet<SeenTypeRecord> unionMembers (nil is null).
	unionMembers *jutil.LinkedHashSet[*JSTypeReconserializer_SeenTypeRecord]
}

// Java: JSTypeReconserializer.SeenTypeRecord(ColorId, int)
func newJSTypeReconserializer_SeenTypeRecord(colorId colors.ColorId, pointer int) *JSTypeReconserializer_SeenTypeRecord {
	r := &JSTypeReconserializer_SeenTypeRecord{
		colorId: colorId,
		pointer: pointer,
	}
	r.jstypes = r.jstypeBuf[:0]
	return r
}

// String is the inherited Object.toString().
//
// Java: Object.toString()
func (r *JSTypeReconserializer_SeenTypeRecord) String() string {
	return jutil.ObjectToString(
		"com.google.javascript.jscomp.serialization.JSTypeReconserializer$SeenTypeRecord", r)
}

// JSTypeReconserializer_IsClosureAssert returns whether this is some
// assertion call that should be removed by optimizations when
// --remove_closure_asserts is enabled.
//
// Java: JSTypeReconserializer.isClosureAssert(ClosurePrimitive)
func JSTypeReconserializer_IsClosureAssert(primitive *rhino.ClosurePrimitive) bool {
	if primitive == nil {
		return false
	}

	switch *primitive {
	case rhino.ClosurePrimitive_ASSERTS_TRUTHY, rhino.ClosurePrimitive_ASSERTS_MATCHES_RETURN:
		return true
	default: // ASSERTS_FAIL: the Java switch expression is exhaustive over the enum.
		// technically an assertion function, but not removed by ClosureCodeRemoval
		return false
	}
}
