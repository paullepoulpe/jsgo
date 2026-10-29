// Port of com/google/javascript/jscomp/serialization/JSTypeReconserializer.java @ closure-compiler v20261021

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

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

type jsTypeDebugMismatch struct {
	sourceRef        string
	involvedPointers []int32
}

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
	allowedProperties           *jutil.BitSet
	debugMismatches             []jsTypeDebugMismatch
	offsetToColor               []*colors.Color

	// Cache some commonly used types.
	unknownRecord     *JSTypeReconserializer_SeenTypeRecord
	topFunctionRecord *JSTypeReconserializer_SeenTypeRecord

	// The pointer of the record of each type recorded so far.
	//
	// Go-only: performance — an identity map from a type to its pointer,
	// equivalent to searching the records' jstypes. recordsByPointer holds the
	// records in exact insertion order alongside seenTypeRecords.
	typePointers           map[*rhino.JSTypeBase]int32
	recordsByPointer       []*JSTypeReconserializer_SeenTypeRecord
	seenTypeRecords        map[colors.ColorId]*JSTypeReconserializer_SeenTypeRecord
	disambiguateEdges      map[int32][]int32
	disambiguateSubtypes   []int32
	totalDisambiguateEdges int

	scratchInstances    []int32
	scratchPrototypes   []int32
	scratchOwnProps     []int32
	seenOwnProps        map[int32]uint32
	seenOwnPropsGen     uint32
	ownPropsUseMap      bool
	scratchAltRecords   []*JSTypeReconserializer_SeenTypeRecord
	seenAltRecords      map[*JSTypeReconserializer_SeenTypeRecord]uint32
	seenAltRecordsGen   uint32
	altRecordsUseMap    bool
	scratchAlternateIds []colors.ColorId
	scratchInt32s       []int32
	singleOwnPropFn     func(string)
	multiOwnPropFn      func(string)

	recordSlab   jutil.Slab[JSTypeReconserializer_SeenTypeRecord]
	int32Slab    []int32
	uuidByteSlab []byte

	state JSTypeReconserializer_State
}

type jsTypeNativeAxiomaticEntry struct {
	native rhino.JSTypeNative
	color  *colors.Color
}

// This is a one-way mapping because some JSTypes go to the same Color.
//
// Java: JSTypeReconserializer.JSTYPE_NATIVE_TO_AXIOMATIC_COLOR_MAP
var jsTypeReconserializer_JSTYPE_NATIVE_TO_AXIOMATIC_COLOR_ENTRIES = [...]jsTypeNativeAxiomaticEntry{
	// Merge all the various top/bottom-like/unknown types into a single unknown type.
	{rhino.JSTypeNative_ALL_TYPE, colors.StandardColors_UNKNOWN},
	{rhino.JSTypeNative_CHECKED_UNKNOWN_TYPE, colors.StandardColors_UNKNOWN},
	{rhino.JSTypeNative_NO_OBJECT_TYPE, colors.StandardColors_UNKNOWN},
	{rhino.JSTypeNative_NO_TYPE, colors.StandardColors_UNKNOWN},
	{rhino.JSTypeNative_UNKNOWN_TYPE, colors.StandardColors_UNKNOWN},
	// Map all the primitives in the obvious way.
	{rhino.JSTypeNative_BIGINT_TYPE, colors.StandardColors_BIGINT},
	{rhino.JSTypeNative_BOOLEAN_TYPE, colors.StandardColors_BOOLEAN},
	{rhino.JSTypeNative_GBIGINT_TYPE, colors.StandardColors_GBIGINT},
	{rhino.JSTypeNative_NULL_TYPE, colors.StandardColors_NULL_OR_VOID},
	{rhino.JSTypeNative_NUMBER_TYPE, colors.StandardColors_NUMBER},
	{rhino.JSTypeNative_STRING_TYPE, colors.StandardColors_STRING},
	{rhino.JSTypeNative_SYMBOL_TYPE, colors.StandardColors_SYMBOL},
	{rhino.JSTypeNative_VOID_TYPE, colors.StandardColors_NULL_OR_VOID},
	// Smoosh top-like function types into a single type.
	{rhino.JSTypeNative_FUNCTION_FUNCTION_TYPE, colors.StandardColors_TOP_FUNCTION},
	{rhino.JSTypeNative_FUNCTION_TYPE, colors.StandardColors_TOP_FUNCTION},
	// Smoosh top-like objects into a single type.
	{rhino.JSTypeNative_FUNCTION_PROTOTYPE, colors.StandardColors_TOP_OBJECT},
	{rhino.JSTypeNative_FUNCTION_INSTANCE_PROTOTYPE, colors.StandardColors_TOP_OBJECT},
	{rhino.JSTypeNative_OBJECT_FUNCTION_TYPE, colors.StandardColors_TOP_OBJECT},
	{rhino.JSTypeNative_OBJECT_PROTOTYPE, colors.StandardColors_TOP_OBJECT},
	{rhino.JSTypeNative_OBJECT_TYPE, colors.StandardColors_TOP_OBJECT},
}

// JSTypeReconserializer_State is JSTypeReconserializer.State.
type JSTypeReconserializer_State int32

const (
	JSTypeReconserializer_State_COLLECTING_TYPES JSTypeReconserializer_State = iota
	JSTypeReconserializer_State_GENERATING_POOL
	JSTypeReconserializer_State_FINISHED
)

var jsTypeReconserializer_State_names = [...]string{"COLLECTING_TYPES", "GENERATING_POOL", "FINISHED"}

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
	s := &JSTypeReconserializer{
		typePointers:                make(map[*rhino.JSTypeBase]int32, 4096),
		recordsByPointer:            make([]*JSTypeReconserializer_SeenTypeRecord, 0, 4096),
		seenTypeRecords:             make(map[colors.ColorId]*JSTypeReconserializer_SeenTypeRecord, 4096),
		disambiguateEdges:           make(map[int32][]int32, 2048),
		disambiguateSubtypes:        make([]int32, 0, 2048),
		state:                       JSTypeReconserializer_State_COLLECTING_TYPES,
		registry:                    registry,
		hasher:                      NewJSTypeColorIdHasher(registry),
		invalidatingTypes:           invalidatingTypes,
		stringPoolBuilder:           stringPoolBuilder,
		shouldPropagatePropertyName: shouldPropagatePropertyName,
		serializationMode:           serializationMode,
	}
	s.singleOwnPropFn = s.onSingleOwnProp
	s.multiOwnPropFn = s.onMultiOwnProp

	s.SeedCachesWithAxiomaticTypes()

	s.unknownRecord = s.seenTypeRecords[colors.StandardColors_UNKNOWN.GetId()]
	s.topFunctionRecord = s.seenTypeRecords[colors.StandardColors_TOP_FUNCTION.GetId()]
	return s
}

func (s *JSTypeReconserializer) onSingleOwnProp(ownProperty string) {
	if s.allowedProperties != nil {
		id := strpool.Lookup(ownProperty)
		if id != 0 && s.allowedProperties.Get(int(id)) {
			s.scratchInt32s = append(s.scratchInt32s, int32(s.stringPoolBuilder.PutID(uint32(id))))
		}
		return
	}
	if s.shouldPropagatePropertyName(ownProperty) {
		s.scratchInt32s = append(s.scratchInt32s, int32(s.stringPoolBuilder.Put(ownProperty)))
	}
}

func (s *JSTypeReconserializer) appendDedupOwnProp(ptr int32) {
	if !s.ownPropsUseMap {
		if slices.Contains(s.scratchOwnProps, ptr) {
			return
		}
		s.scratchOwnProps = append(s.scratchOwnProps, ptr)
		if len(s.scratchOwnProps) == 8 {
			s.ownPropsUseMap = true
			s.seenOwnPropsGen++
			if s.seenOwnProps == nil {
				s.seenOwnProps = make(map[int32]uint32, 64)
			} else if s.seenOwnPropsGen == 0 {
				clear(s.seenOwnProps)
				s.seenOwnPropsGen = 1
			}
			for _, existing := range s.scratchOwnProps {
				s.seenOwnProps[existing] = s.seenOwnPropsGen
			}
		}
		return
	}
	if s.seenOwnProps[ptr] == s.seenOwnPropsGen {
		return
	}
	s.seenOwnProps[ptr] = s.seenOwnPropsGen
	s.scratchOwnProps = append(s.scratchOwnProps, ptr)
}

func (s *JSTypeReconserializer) onMultiOwnProp(ownProperty string) {
	if s.allowedProperties != nil {
		id := strpool.Lookup(ownProperty)
		if id != 0 && s.allowedProperties.Get(int(id)) {
			s.appendDedupOwnProp(int32(s.stringPoolBuilder.PutID(uint32(id))))
		}
		return
	}
	if s.shouldPropagatePropertyName(ownProperty) {
		s.appendDedupOwnProp(int32(s.stringPoolBuilder.Put(ownProperty)))
	}
}

func (s *JSTypeReconserializer) appendDedupAltRecord(rec *JSTypeReconserializer_SeenTypeRecord) {
	if !s.altRecordsUseMap {
		if slices.Contains(s.scratchAltRecords, rec) {
			return
		}
		s.scratchAltRecords = append(s.scratchAltRecords, rec)
		if len(s.scratchAltRecords) == 8 {
			s.altRecordsUseMap = true
			s.seenAltRecordsGen++
			if s.seenAltRecords == nil {
				s.seenAltRecords = make(map[*JSTypeReconserializer_SeenTypeRecord]uint32, 32)
			} else if s.seenAltRecordsGen == 0 {
				clear(s.seenAltRecords)
				s.seenAltRecordsGen = 1
			}
			for _, existing := range s.scratchAltRecords {
				s.seenAltRecords[existing] = s.seenAltRecordsGen
			}
		}
		return
	}
	if s.seenAltRecords[rec] == s.seenAltRecordsGen {
		return
	}
	s.seenAltRecords[rec] = s.seenAltRecordsGen
	s.scratchAltRecords = append(s.scratchAltRecords, rec)
}

func (s *JSTypeReconserializer) allocSeenTypeRecord(colorId colors.ColorId, pointer int) *JSTypeReconserializer_SeenTypeRecord {
	r := s.recordSlab.New()
	r.colorId = colorId
	r.pointer = pointer
	r.jstypes = r.jstypeBuf[:0]
	return r
}

func (s *JSTypeReconserializer) allocInt32Slice(src []int32) []int32 {
	n := len(src)
	if n == 0 {
		return nil
	}
	if len(s.int32Slab) < n {
		s.int32Slab = make([]int32, max(n, 4096))
	}
	out := s.int32Slab[:n:n]
	s.int32Slab = s.int32Slab[n:]
	copy(out, src)
	return out
}

func (s *JSTypeReconserializer) allocSingleInt32(v int32) []int32 {
	if len(s.int32Slab) < 1 {
		s.int32Slab = make([]int32, 4096)
	}
	out := s.int32Slab[:1:1]
	s.int32Slab = s.int32Slab[1:]
	out[0] = v
	return out
}

func (s *JSTypeReconserializer) allocColorIdBytes(id colors.ColorId) []byte {
	if len(s.uuidByteSlab) < 8 {
		s.uuidByteSlab = make([]byte, 8192)
	}
	out := s.uuidByteSlab[:8:8]
	s.uuidByteSlab = s.uuidByteSlab[8:]
	id.PutBigEndian(out)
	return out
}

func (s *JSTypeReconserializer) shouldPropagate(prop string) bool {
	if s.allowedProperties != nil {
		id := strpool.Lookup(prop)
		return id != 0 && s.allowedProperties.Get(int(id))
	}
	return s.shouldPropagatePropertyName(prop)
}

func (s *JSTypeReconserializer) putDisambiguateEdge(subtype, supertype int) {
	sub, sup := int32(subtype), int32(supertype)
	edges, ok := s.disambiguateEdges[sub]
	if slices.Contains(edges, sup) {
		return
	}
	if !ok {
		s.disambiguateSubtypes = append(s.disambiguateSubtypes, sub)
	}
	s.disambiguateEdges[sub] = append(edges, sup)
	s.totalDisambiguateEdges++
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

// LookupTypePointer returns the cached pointer for type_ if it was serialized.
func (s *JSTypeReconserializer) LookupTypePointer(type_ rhino.JSType) (int, bool) {
	p, ok := s.typePointers[type_.GetJSTypeBase()]
	return int(p), ok
}

// hasTypePointer reports whether the type of tb has been recorded.
func (s *JSTypeReconserializer) hasTypePointer(tb *rhino.JSTypeBase) bool {
	_, ok := s.typePointers[tb]
	return ok
}

func (s *JSTypeReconserializer) cacheForwardedRecord(_ rhino.JSType, tb *rhino.JSTypeBase, rec *JSTypeReconserializer_SeenTypeRecord) *JSTypeReconserializer_SeenTypeRecord {
	s.typePointers[tb] = int32(rec.pointer)
	return rec
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
	tb := type_.GetJSTypeBase()
	if p, ok := s.typePointers[tb]; ok {
		return s.recordsByPointer[p]
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
	}
	if forwardedType != nil {
		rec := s.RecordType(forwardedType)
		return s.cacheForwardedRecord(type_, tb, rec)
	}

	if type_.IsUnknownType() || type_.IsNoResolvedType() || type_.IsTemplateType() {
		// template types are not serialized because optimizations don't seem to care about them.
		// serialize as the UNKNOWN_TYPE because bounded generics are unsupported
		return s.cacheForwardedRecord(type_, tb, s.unknownRecord)
	}
	if type_.IsFunctionType() &&
		!type_.ToMaybeFunctionType().HasInstanceType() &&
		type_.ToMaybeFunctionType().GetClosurePrimitive() == nil {
		// Distinguishing different function types does not matter for optimizations unless they
		// are a constructor/interface or have a @closurePrimitive tag associated. Optimization colors
		// don't track function parameter/return/template types, and function literals are all
		// invalidating types during property disambiguation.
		return s.cacheForwardedRecord(type_, tb, s.topFunctionRecord)
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

	alternates := type_.GetAlternates()
	var resolvedBuf [8]*JSTypeReconserializer_SeenTypeRecord
	var resolved []*JSTypeReconserializer_SeenTypeRecord
	if len(alternates) <= len(resolvedBuf) {
		resolved = resolvedBuf[:len(alternates)]
	} else {
		resolved = make([]*JSTypeReconserializer_SeenTypeRecord, len(alternates))
	}
	for i, altType := range alternates {
		resolved[i] = s.RecordType(altType)
	}

	s.scratchAltRecords = s.scratchAltRecords[:0]
	s.altRecordsUseMap = false
	for _, alt := range resolved {
		if len(alt.unionMembers) == 0 {
			s.appendDedupAltRecord(alt)
		} else {
			// Flatten out any nested unions. They are possible due to proxy-like types.
			for _, member := range alt.unionMembers {
				s.appendDedupAltRecord(member)
			}
		}
	}
	altRecords := s.scratchAltRecords

	// Some elements of the union may be equal as Colors
	if len(altRecords) == 1 {
		return s.cacheForwardedRecord(type_, type_.GetJSTypeBase(), altRecords[0])
	}

	alternateIds := s.scratchAlternateIds[:0]
	for _, altRecord := range altRecords {
		alternateIds = append(alternateIds, altRecord.colorId)
	}
	s.scratchAlternateIds = alternateIds
	unionId := unionColorIds(alternateIds)
	record := s.GetOrCreateRecord(unionId, type_)

	if len(record.unionMembers) == 0 {
		n := len(altRecords)
		if n <= len(record.unionMembersBuf) {
			record.unionMembers = append(record.unionMembersBuf[:0], altRecords...)
		} else {
			record.unionMembers = append(make([]*JSTypeReconserializer_SeenTypeRecord, 0, n), altRecords...)
		}
	} else if s.serializationMode.RunValidation() {
		same := len(altRecords) == len(record.unionMembers)
		if same {
			for _, m := range record.unionMembers {
				found := false
				if s.altRecordsUseMap {
					found = s.seenAltRecords[m] == s.seenAltRecordsGen
				} else {
					found = slices.Contains(altRecords, m)
				}
				if !found {
					same = false
					break
				}
			}
		}
		if !same {
			jutil.CheckStateMsg3(
				false,
				"Unions with same ID must have same members: %s => %s == %s",
				unionId,
				jsTypeReconserializer_formatColorIds(altRecords),
				jsTypeReconserializer_formatColorIds(record.unionMembers))
		}
	}

	return record
}

func jsTypeReconserializer_formatColorIds(records []*JSTypeReconserializer_SeenTypeRecord) string {
	var sb strings.Builder
	sb.WriteByte('[')
	seen := make([]colors.ColorId, 0, len(records))
	for _, r := range records {
		if !slices.Contains(seen, r.colorId) {
			if len(seen) > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(r.colorId.String())
			seen = append(seen, r.colorId)
		}
	}
	sb.WriteByte(']')
	return sb.String()
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

	record := s.seenTypeRecords[id]
	if record == nil {
		pointer := len(s.recordsByPointer)
		record = s.allocSeenTypeRecord(id, pointer)
		s.seenTypeRecords[id] = record
		s.recordsByPointer = append(s.recordsByPointer, record)
	}
	tb := jstype.GetJSTypeBase()
	if _, ok := s.typePointers[tb]; !ok {
		s.typePointers[tb] = int32(record.pointer)
		record.jstypes = append(record.jstypes, jstype)
	}

	return record
}

// Java: JSTypeReconserializer.reconcileUnionTypes(SeenTypeRecord)
func (s *JSTypeReconserializer) ReconcileUnionTypes(seen *JSTypeReconserializer_SeenTypeRecord) *TypeProto {
	var tp TypeProto
	var up UnionTypeProto
	s.reconcileUnionTypesInto(seen, &tp, &up)
	return &tp
}

func (s *JSTypeReconserializer) reconcileUnionTypesInto(seen *JSTypeReconserializer_SeenTypeRecord, tp *TypeProto, up *UnionTypeProto) {
	members := s.scratchInt32s[:0]
	for _, r := range seen.unionMembers {
		members = append(members, int32(r.pointer))
	}
	slices.Sort(members)
	s.scratchInt32s = members
	*up = UnionTypeProto{
		unionMember_: s.allocInt32Slice(members),
	}
	*tp = TypeProto{
		union_:    up,
		kindCase_: TypeProto_UNION_FIELD_NUMBER,
	}
}

// Java: JSTypeReconserializer.reconcileObjectTypes(SeenTypeRecord)
func (s *JSTypeReconserializer) ReconcileObjectTypes(seen *JSTypeReconserializer_SeenTypeRecord) *TypeProto {
	var tp TypeProto
	var op ObjectTypeProto
	s.reconcileObjectTypesInto(seen, &tp, &op)
	return &tp
}

func (s *JSTypeReconserializer) reconcileObjectTypesInto(seen *JSTypeReconserializer_SeenTypeRecord, tp *TypeProto, op *ObjectTypeProto) {
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
				instanceTypes = s.allocSingleInt32(int32(s.SerializeType(fnType.GetInstanceType())))
				prototypes = s.allocSingleInt32(int32(s.SerializeType(fnType.GetPrototype())))
				isConstructor = fnType.IsConstructor()
			}
			isClosureAssert = JSTypeReconserializer_IsClosureAssert(fnType.GetClosurePrimitive())
		}

		s.scratchInt32s = s.scratchInt32s[:0]
		objType.ForEachOwnPropertyName(s.singleOwnPropFn)
		ownProps := s.scratchInt32s

		isInvalidating := s.invalidatingTypes.IsInvalidating(objType)
		propertiesKeepOriginalName := objType.IsEnumType()

		*op = ObjectTypeProto{
			instanceType_:               instanceTypes,
			ownProperty_:                s.allocInt32Slice(ownProps),
			prototype_:                  prototypes,
			closureAssert_:              isClosureAssert,
			isInvalidating_:             isInvalidating,
			markedConstructor_:          isConstructor,
			propertiesKeepOriginalName_: propertiesKeepOriginalName,
			uuid_:                       s.allocColorIdBytes(seen.colorId),
		}
		*tp = TypeProto{
			object_:   op,
			kindCase_: TypeProto_OBJECT_FIELD_NUMBER,
		}
		return
	}

	s.scratchInstances = s.scratchInstances[:0]
	s.scratchPrototypes = s.scratchPrototypes[:0]
	s.scratchOwnProps = s.scratchOwnProps[:0]
	s.ownPropsUseMap = false
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
				instPtr := int32(s.SerializeType(fnType.GetInstanceType()))
				if !slices.Contains(s.scratchInstances, instPtr) {
					s.scratchInstances = append(s.scratchInstances, instPtr)
				}
				protoPtr := int32(s.SerializeType(fnType.GetPrototype()))
				if !slices.Contains(s.scratchPrototypes, protoPtr) {
					s.scratchPrototypes = append(s.scratchPrototypes, protoPtr)
				}
				isConstructor = isConstructor || fnType.IsConstructor()
			}

			isClosureAssert = isClosureAssert || JSTypeReconserializer_IsClosureAssert(fnType.GetClosurePrimitive())
		}

		objType.ForEachOwnPropertyName(s.multiOwnPropFn)

		isInvalidating = isInvalidating || s.invalidatingTypes.IsInvalidating(objType)

		/*
		 * To support legacy code, property disambiguation never renames properties of enums (e.g. 'A'
		 * in '/** @enum * / const E = {A: 0}`). In theory this would be safe to remove if we clean up
		 * code depending on the lack of renaming
		 */
		propertiesKeepOriginalName = propertiesKeepOriginalName || objType.IsEnumType()
	}

	*op = ObjectTypeProto{
		instanceType_:               s.allocInt32Slice(s.scratchInstances),
		ownProperty_:                s.allocInt32Slice(s.scratchOwnProps),
		prototype_:                  s.allocInt32Slice(s.scratchPrototypes),
		closureAssert_:              isClosureAssert,
		isInvalidating_:             isInvalidating,
		markedConstructor_:          isConstructor,
		propertiesKeepOriginalName_: propertiesKeepOriginalName,
		uuid_:                       s.allocColorIdBytes(seen.colorId),
	}
	*tp = TypeProto{
		object_:   op,
		kindCase_: TypeProto_OBJECT_FIELD_NUMBER,
	}
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
	jutil.CheckState(len(s.recordsByPointer) == 0)

	// Load all the axiomatic records in the right order without any types.
	for _, axiomatic := range TypePointers_OFFSET_TO_AXIOMATIC_COLOR {
		index := len(s.recordsByPointer)
		record := s.allocSeenTypeRecord(axiomatic.GetId(), index)
		s.seenTypeRecords[axiomatic.GetId()] = record
		s.recordsByPointer = append(s.recordsByPointer, record)
	}

	// Add JSTypes corresponding to axiomatic IDs.
	for _, entry := range jsTypeReconserializer_JSTYPE_NATIVE_TO_AXIOMATIC_COLOR_ENTRIES {
		s.GetOrCreateRecord(entry.color.GetId(), s.registry.GetNativeType(entry.native))
	}

	jutil.CheckState(len(s.recordsByPointer) == TypePointers_AXIOMATIC_COLOR_COUNT)
}

// CheckValidLinearTime checks that this instance is in a valid state.
//
// Java: JSTypeReconserializer.checkValidLinearTime()
func (s *JSTypeReconserializer) CheckValidLinearTime() {
	if !s.serializationMode.RunValidation() {
		return
	}

	totalTypeCount := len(s.recordsByPointer)
	for _, seen := range s.recordsByPointer {
		offset := seen.pointer
		jutil.CheckState(offset >= 0)
		jutil.CheckStateMsg2(
			offset <= totalTypeCount,
			"Found invalid pointer %s, out of a total of %s user-defined types",
			offset,
			totalTypeCount)
	}
}

func (s *JSTypeReconserializer) collectSortedInvolvedPointers(types *jutil.ObjectLinkedHashSet[rhino.JSType]) []int32 {
	involvedColors := make([]int32, 0, types.Size())
	for t := range types.All() {
		jutil.CheckStateMsg1(!t.IsUnionType(), "%s", t)
		// Ensure all types are recorded before reconciliation.
		involvedColors = append(involvedColors, int32(s.SerializeType(t)))
	}
	slices.Sort(involvedColors)
	return slices.Compact(involvedColors)
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
		var mismatches []*TypePool_DebugInfo_Mismatch
		// Key by source ref to deduplicate the strings, which are pretty long.
		for location, types := range s.invalidatingTypes.
			GetMismatchLocations().
			Inverse().
			AsMap().
			All() {
			mismatches = append(mismatches, &TypePool_DebugInfo_Mismatch{
				sourceRef_:     location.GetLocation(),
				involvedColor_: s.collectSortedInvolvedPointers(types),
			})
		}
		debugInfoProto = &TypePool_DebugInfo{
			mismatch_: mismatches,
		}
	}

	s.state = JSTypeReconserializer_State_GENERATING_POOL

	// Go-only: language-adapter — Java iterates seenTypeRecords.values(),
	// which would throw ConcurrentModificationException if reconciliation
	// recorded a new type; recordObjectType serializes instance types and
	// prototypes ahead of time so that it never does.
	numTypes := len(s.recordsByPointer) - TypePointers_AXIOMATIC_COLOR_COUNT
	var types []*TypeProto
	if numTypes > 0 {
		types = make([]*TypeProto, numTypes)
		typeProtoSlab := make([]TypeProto, numTypes)
		numUnions := 0
		for _, seen := range s.recordsByPointer[TypePointers_AXIOMATIC_COLOR_COUNT:] {
			if len(seen.unionMembers) != 0 {
				numUnions++
			}
		}
		numObjects := numTypes - numUnions
		var objProtoSlab []ObjectTypeProto
		if numObjects > 0 {
			objProtoSlab = make([]ObjectTypeProto, numObjects)
		}
		var unionProtoSlab []UnionTypeProto
		if numUnions > 0 {
			unionProtoSlab = make([]UnionTypeProto, numUnions)
		}
		typeIdx := 0
		objIdx := 0
		unionIdx := 0
		for _, seen := range s.recordsByPointer {
			if TypePointers_IsAxiomatic(seen.pointer) {
				jutil.CheckStateMsg1(isAxiomaticColorId(seen.colorId), "Missing .type for SeenTypeRecord %s", seen)
				continue
			}
			tp := &typeProtoSlab[typeIdx]
			types[typeIdx] = tp
			typeIdx++
			if len(seen.unionMembers) != 0 {
				s.reconcileUnionTypesInto(seen, tp, &unionProtoSlab[unionIdx])
				unionIdx++
			} else {
				s.reconcileObjectTypesInto(seen, tp, &objProtoSlab[objIdx])
				objIdx++
			}
		}
	} else {
		for _, seen := range s.recordsByPointer {
			if isAxiomaticColorId(seen.colorId) {
				jutil.CheckStateMsg1(TypePointers_IsAxiomatic(seen.pointer), "Missing .type for SeenTypeRecord %s", seen)
			}
		}
	}

	var disambiguationEdges []*SubtypingEdge
	if s.totalDisambiguateEdges > 0 {
		disambiguationEdges = make([]*SubtypingEdge, s.totalDisambiguateEdges)
		edgeSlab := make([]SubtypingEdge, s.totalDisambiguateEdges)
		edgeIdx := 0
		for _, subtype := range s.disambiguateSubtypes {
			for _, sup := range s.disambiguateEdges[subtype] {
				e := &edgeSlab[edgeIdx]
				*e = SubtypingEdge{
					subtype_:   subtype,
					supertype_: sup,
				}
				disambiguationEdges[edgeIdx] = e
				edgeIdx++
			}
		}
	}

	s.state = JSTypeReconserializer_State_FINISHED
	s.CheckValidLinearTime()
	return &TypePool{
		type_:                types,
		disambiguationEdges_: disambiguationEdges,
		debugInfo_:           debugInfoProto,
	}
}

func (s *JSTypeReconserializer) prepareForColorPool() {
	jutil.CheckState(s.state == JSTypeReconserializer_State_COLLECTING_TYPES)
	s.CheckValidLinearTime()

	if s.serializationMode.IncludeDebugInfo() {
		for location, types := range s.invalidatingTypes.
			GetMismatchLocations().
			Inverse().
			AsMap().
			All() {
			s.debugMismatches = append(s.debugMismatches, jsTypeDebugMismatch{
				sourceRef:        location.GetLocation(),
				involvedPointers: s.collectSortedInvolvedPointers(types),
			})
		}
	}

	s.state = JSTypeReconserializer_State_GENERATING_POOL
	if s.serializationMode.RunValidation() {
		for _, seen := range s.recordsByPointer {
			if TypePointers_IsAxiomatic(seen.pointer) {
				jutil.CheckStateMsg1(
					isAxiomaticColorId(seen.colorId),
					"Missing .type for SeenTypeRecord %s",
					seen)
			}
		}
	}
}

func (s *JSTypeReconserializer) finishColorPool() {
	s.state = JSTypeReconserializer_State_FINISHED
	s.CheckValidLinearTime()
}

// GetColorIdToJSTypeMapForDebugging returns a map from
// ObjectTypeProto#getUuid() to the originating JSTypes.
//
// Only intended to be used for debug logging.
//
// Java: JSTypeReconserializer.getColorIdToJSTypeMapForDebugging()
func (s *JSTypeReconserializer) GetColorIdToJSTypeMapForDebugging() *jutil.LinkedHashMap[string, []rhino.JSType] {
	// note: returns JSType values instead of String values, even though all that's needed for
	// debugging are the strings, to avoid the memory overhead of calculating all string
	// representations at once
	var keys []string
	byId := make(map[string][]rhino.JSType)
	for _, record := range s.recordsByPointer {
		if len(record.jstypes) == 0 {
			continue
		}
		idStr := record.colorId.String()
		if _, exists := byId[idStr]; !exists {
			keys = append(keys, idStr)
		}
		byId[idStr] = append(byId[idStr], record.jstypes...)
	}
	slices.Sort(keys)
	colorIdToTypes := jutil.NewLinkedHashMapWithCapacity[string, []rhino.JSType](len(keys))
	for _, k := range keys {
		colorIdToTypes.Put(k, byId[k])
	}
	return colorIdToTypes
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
	jstypeBuf [4]rhino.JSType

	unionMembers    []*JSTypeReconserializer_SeenTypeRecord
	unionMembersBuf [4]*JSTypeReconserializer_SeenTypeRecord
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
