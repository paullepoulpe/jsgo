// Port of com/google/javascript/rhino/jstype/JSTypeRegistry.java @ closure-compiler v20261021

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Bob Jervis
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino

import (
	"hash/maphash"
	"sync"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

const (
	jsTypeRegistryLeastSupertypeCacheSize = 4096
	jsTypeRegistryTemplatizedCacheSize    = 4096
)

type leastSupertypeCacheEntry struct {
	a, b *JSTypeBase
	res  JSType
	gen  uint32
}

type templateTypeMapCacheEntry struct {
	base   *TemplateTypeMap
	t0, t1 *JSTypeBase
	res    *TemplateTypeMap
}

type propTypeSlot struct {
	t    *JSTypeBase
	hash int32
}

type propTypeOverflow struct {
	slots []propTypeSlot
	table []uint32
}

// propTypeSet is an insertion-ordered set of *JSTypeBase stored by value in
// JSTypeRegistry's property-indexed maps: up to two types live inline without
// heap allocation or interior slice pointers, spilling to overflow on >=3.
type propTypeSet struct {
	t0, t1   *JSTypeBase
	overflow *propTypeOverflow
	h0, h1   int32
}

func propTypeBasesEqual(a, b *JSTypeBase) bool {
	return a == b || equalityChecker_AreIdentityEqualWithBases(a.self_, a, b.self_, b)
}

func propTypeHashMix(h int32) uint {
	u := uint32(h)
	return uint(u ^ (u >> 16))
}

func (s *propTypeSet) add(tb *JSTypeBase) bool {
	return s.addWithHash(tb, tb.HashCode())
}

func (s *propTypeSet) addWithHash(tb *JSTypeBase, hash int32) bool {
	if s.t0 == nil {
		s.t0 = tb
		s.h0 = hash
		return true
	}
	if s.h0 == hash && propTypeBasesEqual(s.t0, tb) {
		return false
	}
	if s.t1 == nil {
		s.t1 = tb
		s.h1 = hash
		return true
	}
	if s.h1 == hash && propTypeBasesEqual(s.t1, tb) {
		return false
	}
	if s.overflow == nil {
		o := &propTypeOverflow{slots: make([]propTypeSlot, 1, 4)}
		o.slots[0] = propTypeSlot{t: tb, hash: hash}
		s.overflow = o
		return true
	}
	s.overflow.add(tb, hash)
	return false
}

func (o *propTypeOverflow) add(tb *JSTypeBase, hash int32) {
	if o.table == nil {
		for i := range o.slots {
			if o.slots[i].hash == hash && propTypeBasesEqual(o.slots[i].t, tb) {
				return
			}
		}
		o.slots = append(o.slots, propTypeSlot{t: tb, hash: hash})
		if len(o.slots) > 8 {
			o.rebuildTable(32)
		}
		return
	}
	mask := uint(len(o.table) - 1)
	b := propTypeHashMix(hash) & mask
	for {
		e := o.table[b]
		if e == 0 {
			o.slots = append(o.slots, propTypeSlot{t: tb, hash: hash})
			if len(o.slots)*2 > len(o.table) {
				o.rebuildTable(len(o.table) * 2)
			} else {
				o.table[b] = uint32(len(o.slots))
			}
			return
		}
		slot := o.slots[e-1]
		if slot.hash == hash && propTypeBasesEqual(slot.t, tb) {
			return
		}
		b = (b + 1) & mask
	}
}

func (o *propTypeOverflow) rebuildTable(size int) {
	table := make([]uint32, size)
	mask := uint(size - 1)
	for i := range o.slots {
		b := propTypeHashMix(o.slots[i].hash) & mask
		for table[b] != 0 {
			b = (b + 1) & mask
		}
		table[b] = uint32(i + 1)
	}
	o.table = table
}

func (s propTypeSet) toObjectTypeSlice() []ObjectType {
	if s.t0 == nil {
		return []ObjectType{}
	}
	if s.t1 == nil {
		return []ObjectType{s.t0.self_.(ObjectType)}
	}
	if s.overflow == nil {
		return []ObjectType{s.t0.self_.(ObjectType), s.t1.self_.(ObjectType)}
	}
	slots := s.overflow.slots
	out := make([]ObjectType, 2+len(slots))
	out[0] = s.t0.self_.(ObjectType)
	out[1] = s.t1.self_.(ObjectType)
	for i, slot := range slots {
		out[2+i] = slot.t.self_.(ObjectType)
	}
	return out
}

func (s propTypeSet) forEachObjectType(fn func(ObjectType)) {
	if s.t0 == nil {
		return
	}
	fn(s.t0.self_.(ObjectType))
	if s.t1 == nil {
		return
	}
	fn(s.t1.self_.(ObjectType))
	if s.overflow != nil {
		for _, slot := range s.overflow.slots {
			fn(slot.t.self_.(ObjectType))
		}
	}
}

// JSTypeRegistry is the type registry used to resolve named types.
//
// This class is not thread-safe.
type JSTypeRegistry struct {
	// The template variable in IObject<IOBJECT_KEY, IOBJECT_VALUE> (plus the
	// builtin Javascript Object).
	iObjectIndexTemplateKey *TemplateType

	// The template variable in IObject<IOBJECT_KEY, IOBJECT_VALUE> (plus the
	// builtin Javascript Object).
	iObjectElementTemplateKey *TemplateType

	// The template variable corresponding to the T type in Iterable<T,
	// TReturn, TNext>
	iterableValueTemplate *TemplateType

	// The template variable corresponding to the T type in Iterable<T,
	// TReturn, TNext>
	iterableReturnTemplate *TemplateType

	// The template variable corresponding to the T type in IteratorIterable<T,
	// TReturn, TNext>
	iteratorIterableValueTemplate *TemplateType

	// The template variable corresponding to the TReturn type in
	// IteratorIterable<T, TReturn, TNext>
	iteratorIterableReturnTemplate *TemplateType

	// The template variable corresponding to the T type in IteratorLike<T,
	// TReturn, TNext>.
	iteratorLikeValueTemplate *TemplateType

	// The template variable corresponding to the TReturn type in
	// IteratorLike<T, TReturn, TNext>.
	iteratorLikeReturnTemplate *TemplateType

	// The template variable corresponding to the T type in Iterator<T,
	// TReturn, TNext>.
	iteratorValueTemplate *TemplateType

	// The template variable corresponding to the TReturn type in Iterator<T,
	// TReturn, TNext>.
	iteratorReturnTemplate *TemplateType

	// The template variable corresponding to the TYield type in
	// IIterableResult<TYield>
	iiterableResultValueTemplate *TemplateType

	// The template variable corresponding to the T type in AsyncIterable<T,
	// TReturn, TNext>
	asyncIterableValueTemplate *TemplateType

	// The template variable corresponding to the TReturn type in
	// AsyncIterable<T, TReturn, TNext>
	asyncIterableReturnTemplate *TemplateType

	// The template variable corresponding to the T type in AsyncIterator<T,
	// TReturn, TNext>
	asyncIteratorValueTemplate *TemplateType

	// The template variable corresponding to the TReturn type in
	// AsyncIterator<T, TReturn, TNext>
	asyncIteratorReturnTemplate *TemplateType

	// The template variable corresponding to the T type in Generator<T,
	// TReturn, TNext>
	generatorValueTemplate *TemplateType

	// The template variable corresponding to the TReturn type in Generator<T,
	// TReturn, TNext>
	generatorReturnTemplate *TemplateType

	// The template variable corresponding to the TYPE type in IThenable<TYPE>
	iThenableTemplateKey *TemplateType

	// The template variable corresponding to the TYPE in Promise<TYPE>
	promiseTemplateKey *TemplateType

	// The template variable in Array<T>.
	arrayElementTemplateKey *TemplateType

	// The template variable in ReadonlyArray<T>.
	readonlyArrayElementTemplateKey *TemplateType

	// The key variable in ReadonlyMap<K, V>.
	readonlyMapKeyTemplateKey *TemplateType

	// The value variable in ReadonlyMap<K, V>.
	readonlyMapValueTemplateKey *TemplateType

	// The key variable in Map<K, V>.
	mapKeyTemplateKey *TemplateType

	// The value variable in Map<K, V>.
	mapValueTemplateKey *TemplateType

	// TODO(user): An instance of this class should be used during
	// compilation. We also want to make all types' constructors package
	// private and force usage of this registry instead. This will allow us to
	// evolve the types without being tied by an open API.

	reporter ErrorReporter

	// nodeStores are the node stores whose type tables hold this registry's
	// types (node_tables.go); cleared by ReleaseNodeTypeHandles.
	nodeStoresMu sync.Mutex
	nodeStores   []*NodeStore

	// We use an Array instead of an immutable list because this lookup needs
	// to be very fast. When it was an immutable list, we were spending 5% of
	// CPU time on bounds checking inside get().
	nativeTypes []JSType

	// Table<Node, String, JSType> (a HashBasedTable).
	scopedNameTable *jutil.LinkedHashMap[Node, *jutil.LinkedHashMap[string, JSType]]

	// Only needed for type resolution at the moment
	closureNamespaces *jutil.LinkedHashMap[string, *JSTypeRegistry_ClosureNamespace]

	// NOTE: This would normally be "static final" but that causes unit test
	// failures when serializing and deserializing compiler state for
	// multistage builds.
	nameTableGlobalRoot Node

	// NOTE(nicksantos): This is a terrible terrible hack. When type
	// expressions are evaluated, we need to be able to decide whether that
	// type name resolves to a nullable type or a non-nullable type. Object
	// types are nullable, but enum types and typedefs are not.
	//
	// Notice that it's not good enough to just declare enum types and
	// typedefs sooner. For example, if we have
	//   /** @enum {MyObject} */ var MyEnum = ...;
	// we won't be to declare "MyEnum" without evaluating the expression
	// {MyObject}, and following those dependencies starts to lead us into
	// undecidable territory. Instead, we "pre-declare" enum types and
	// typedefs, so that the expression resolver can decide whether a given
	// name is nullable or not. Also, note that this solution is buggy and the
	// type resolution still gets default nullability wrong in some cases. See
	// b/116853368. Probably NamedType should be responsible for applying
	// nullability to forward references instead of the type expression
	// evaluation.
	//
	// Multimap<Node, String> with hash keys and hash set values.
	nonNullableTypeNames *jutil.LinkedHashMap[Node, *jutil.LinkedHashSet[string]]

	// Types that have been "forward-declared." If these types are not
	// declared anywhere in the binary, we shouldn't try to type-check them at
	// all.
	forwardDeclaredTypes jutil.Set[string]

	// A map of properties to the types on which those properties have been
	// declared. "Reference" types are excluded because those already exist in
	// eachRefTypeIndexedByProperty to avoid blowing up the size of this map.
	nonRefTypesIndexedByProperty map[rhinoStringID]propTypeSet

	sentinelObjectLiteral     JSType
	sentinelObjectLiteralHash int32

	// To avoid blowing up the size of nonRefTypesIndexedByProperty, we use the
	// sentinel object literal instead of registering arbitrarily many types.
	// But because of the way unions are constructed, some properties of
	// record types in unions are getting dropped and cause spurious
	// "non-existent property" warnings. The next two fields avoid the
	// warnings. The first field contains property names of records that
	// participate in unions, and have caused properties to be dropped. The
	// second field contains the names of the dropped properties. When
	// checking canPropertyBeDefined, if the type has a property in
	// propertiesOfSupertypesInUnions, we consider it to possibly have any
	// property in droppedPropertiesOfUnions. This is a loose check, but we
	// restrict it to records that may be present in unions, and it allows us
	// to keep nonRefTypesIndexedByProperty small.
	propertiesOfSupertypesInUnions map[rhinoStringID]struct{}
	droppedPropertiesOfUnions      map[rhinoStringID]struct{}

	// A map of properties to each reference type on which those properties
	// have been declared.
	eachRefTypeIndexedByProperty map[rhinoStringID]propTypeSet

	// A single empty TemplateTypeMap, which can be safely reused in cases
	// where there are no template types.
	emptyTemplateTypeMap *TemplateTypeMap

	resolver *JSTypeResolver

	// The symbols used as property keys (see propertySymbolID).
	propertySymbols []*KnownSymbolType

	unionBuilderStack    [16]UnionType_Builder
	unionBuilderDepth    int
	fnBuilderStack       [8]FunctionType_Builder
	fnBuilderDepth       int
	leastSupertypeGen    uint32
	leastSupertypeCache  *[jsTypeRegistryLeastSupertypeCacheSize]leastSupertypeCacheEntry
	templateTypeMapCache *[jsTypeRegistryTemplatizedCacheSize]templateTypeMapCacheEntry
	objectDocInfo        map[*ObjectTypeBase]*JSDocInfo
	parallelReadOnly     bool
	parallelMu           sync.Mutex
}

// SetParallelReadOnly toggles concurrent read-only query mode on JSTypeRegistry,
// bypassing single-threaded builder/cache mutations and guarding union
// dropped-property sets with a mutex. Returns the previous mode.
func (r *JSTypeRegistry) SetParallelReadOnly(enabled bool) bool {
	if r == nil {
		return false
	}
	if enabled {
		_ = r.GetSentinelObjectLiteral()
	}
	prev := r.parallelReadOnly
	r.parallelReadOnly = enabled
	return prev
}

func (r *JSTypeRegistry) getObjectDocInfo(o *ObjectTypeBase) *JSDocInfo {
	if r == nil || r.objectDocInfo == nil {
		return nil
	}
	return r.objectDocInfo[o]
}

func (r *JSTypeRegistry) setObjectDocInfo(o *ObjectTypeBase, info *JSDocInfo) {
	if r.objectDocInfo == nil {
		r.objectDocInfo = make(map[*ObjectTypeBase]*JSDocInfo, 64)
	}
	r.objectDocInfo[o] = info
}

func (r *JSTypeRegistry) deleteObjectDocInfo(o *ObjectTypeBase) {
	if r != nil && r.objectDocInfo != nil {
		delete(r.objectDocInfo, o)
	}
}

func (r *JSTypeRegistry) acquireFunctionBuilder() *FunctionType_Builder {
	if r != nil && !r.parallelReadOnly && r.fnBuilderDepth < len(r.fnBuilderStack) {
		b := &r.fnBuilderStack[r.fnBuilderDepth]
		r.fnBuilderDepth++
		*b = FunctionType_Builder{}
		b.pooled = true
		return b
	}
	return &FunctionType_Builder{}
}

func (r *JSTypeRegistry) releaseFunctionBuilder(b *FunctionType_Builder) {
	if b.pooled && r != nil && !r.parallelReadOnly && r.fnBuilderDepth > 0 && b == &r.fnBuilderStack[r.fnBuilderDepth-1] {
		b.pooled = false
		r.fnBuilderDepth--
	}
}

func (r *JSTypeRegistry) acquireUnionBuilder() *UnionType_Builder {
	if r != nil && !r.parallelReadOnly && r.unionBuilderDepth < len(r.unionBuilderStack) {
		b := &r.unionBuilderStack[r.unionBuilderDepth]
		r.unionBuilderDepth++
		initUnionType_Builder(b, r)
		return b
	}
	return NewUnionType_Builder(r)
}

func (r *JSTypeRegistry) releaseUnionBuilder(b *UnionType_Builder) {
	if r != nil && !r.parallelReadOnly && r.unionBuilderDepth > 0 && b == &r.unionBuilderStack[r.unionBuilderDepth-1] {
		switch b.inlineLen {
		case 0:
		case 1:
			b.inlineAlternates[0] = nil
		case 2:
			b.inlineAlternates[0] = nil
			b.inlineAlternates[1] = nil
		default:
			clear(b.inlineAlternates[:b.inlineLen])
		}
		b.rebuildTarget = nil
		b.alternates = nil
		b.finalAlternates = nil
		r.unionBuilderDepth--
	}
}

var jsTypeRegistryHashSeed = maphash.MakeSeed()

func jsTypeBasePairCacheIndex(a, b *JSTypeBase) int {
	h := maphash.Comparable(jsTypeRegistryHashSeed, struct{ a, b *JSTypeBase }{a, b})
	return int(h & (jsTypeRegistryLeastSupertypeCacheSize - 1))
}

func (r *JSTypeRegistry) lookupLeastSupertype(a, b *JSTypeBase) JSType {
	if r.parallelReadOnly || r.leastSupertypeCache == nil {
		return nil
	}
	entry := &r.leastSupertypeCache[jsTypeBasePairCacheIndex(a, b)]
	if entry.gen == r.leastSupertypeGen && entry.a == a && entry.b == b {
		return entry.res
	}
	return nil
}

func (r *JSTypeRegistry) storeLeastSupertype(a, b *JSTypeBase, res JSType) {
	if r.parallelReadOnly {
		return
	}
	if r.leastSupertypeCache == nil {
		r.leastSupertypeCache = new([jsTypeRegistryLeastSupertypeCacheSize]leastSupertypeCacheEntry)
	}
	r.leastSupertypeCache[jsTypeBasePairCacheIndex(a, b)] = leastSupertypeCacheEntry{a: a, b: b, res: res, gen: r.leastSupertypeGen}
}

func (r *JSTypeRegistry) templateTypeMapCacheIndex(base *TemplateTypeMap, t0, t1 *JSTypeBase) int {
	type cacheKey struct {
		base   *TemplateTypeMap
		t0, t1 *JSTypeBase
	}
	h := maphash.Comparable(jsTypeRegistryHashSeed, cacheKey{base, t0, t1})
	return int(h & (jsTypeRegistryTemplatizedCacheSize - 1))
}

// Java: JSTypeRegistry.I_OBJECT_ELEMENT_TEMPLATE
const JSTypeRegistry_I_OBJECT_ELEMENT_TEMPLATE = "IOBJECT_VALUE"

// Deprecated: Java: JSTypeRegistry.OBJECT_ELEMENT_TEMPLATE
const JSTypeRegistry_OBJECT_ELEMENT_TEMPLATE = JSTypeRegistry_I_OBJECT_ELEMENT_TEMPLATE

// Java: JSTypeRegistry(ErrorReporter)
func NewJSTypeRegistry(reporter ErrorReporter) *JSTypeRegistry {
	return NewJSTypeRegistryWithForwardDeclaredTypes(reporter, jutil.NewLinkedHashSet[string]())
}

// NewJSTypeRegistryWithForwardDeclaredTypes constructs a new type registry
// populated with the built-in types.
//
// Java: JSTypeRegistry(ErrorReporter, Set<String>)
func NewJSTypeRegistryWithForwardDeclaredTypes(reporter ErrorReporter, forwardDeclaredTypes jutil.Set[string]) *JSTypeRegistry {
	r := &JSTypeRegistry{
		scopedNameTable:              jutil.NewLinkedHashMap[Node, *jutil.LinkedHashMap[string, JSType]](),
		closureNamespaces:            jutil.NewLinkedHashMap[string, *JSTypeRegistry_ClosureNamespace](),
		nameTableGlobalRoot:          NewNode(Token_ROOT),
		nonNullableTypeNames:         jutil.NewLinkedHashMap[Node, *jutil.LinkedHashSet[string]](),
		nonRefTypesIndexedByProperty: make(map[rhinoStringID]propTypeSet),
		eachRefTypeIndexedByProperty: make(map[rhinoStringID]propTypeSet),
		leastSupertypeGen:            1,
	}
	r.reporter = reporter
	r.forwardDeclaredTypes = forwardDeclaredTypes
	r.emptyTemplateTypeMap = TemplateTypeMap_CreateEmpty(r)
	r.resolver = JSTypeResolver_Create(r)
	r.nativeTypes = make([]JSType, len(JSTypeNative_values()))

	func() {
		closer := r.resolver.OpenForDefinition()
		defer closer.Close()
		r.InitializeBuiltInTypes()
		r.InitializeRegistry()
	}()
	return r
}

// Java: JSTypeRegistry.getSentinelObjectLiteral()
func (r *JSTypeRegistry) GetSentinelObjectLiteral() JSType {
	if r.sentinelObjectLiteral == nil {
		r.sentinelObjectLiteral = r.CreateAnonymousObjectType(nil)
		r.sentinelObjectLiteralHash = r.sentinelObjectLiteral.HashCode()
	}
	return r.sentinelObjectLiteral
}

// Java: JSTypeRegistry.getErrorReporter()
func (r *JSTypeRegistry) GetErrorReporter() ErrorReporter {
	return r.reporter
}

// Java: JSTypeRegistry.registerNativeType(JSTypeNative, JSType)
func (r *JSTypeRegistry) RegisterNativeType(typeId JSTypeNative, type_ JSType) {
	r.nativeTypes[typeId.Ordinal()] = type_
}

// Java: JSTypeRegistry.getNativeType(JSTypeNative)
func (r *JSTypeRegistry) GetNativeType(typeId JSTypeNative) JSType {
	return r.nativeTypes[typeId.Ordinal()]
}

// Java: JSTypeRegistry.getNativeObjectType(JSTypeNative)
func (r *JSTypeRegistry) GetNativeObjectType(typeId JSTypeNative) ObjectType {
	t := r.GetNativeType(typeId)
	if t == nil {
		return nil
	}
	return t.(ObjectType)
}

// Java: JSTypeRegistry.getNativeFunctionType(JSTypeNative)
func (r *JSTypeRegistry) GetNativeFunctionType(typeId JSTypeNative) FunctionType {
	t := r.GetNativeType(typeId)
	if t == nil {
		return nil
	}
	return t.(FunctionType)
}

// Java: JSTypeRegistry.getResolver()
func (r *JSTypeRegistry) GetResolver() *JSTypeResolver {
	return r.resolver
}

// JSTypeRegistry_IsObjectLiteralThatCanBeSkipped: when t is an object that
// is not the prototype of some class, and its nominal type is Object, and it
// has some properties, we don't need to store these properties in the
// propertyIndex separately.
//
// Java: JSTypeRegistry.isObjectLiteralThatCanBeSkipped(JSType)
func JSTypeRegistry_IsObjectLiteralThatCanBeSkipped(t JSType) bool {
	t = t.RestrictByNotNullOrUndefined()
	return t.IsRecordType() || t.IsLiteralObject()
}

// Java: JSTypeRegistry.registerDroppedPropertiesInUnion(RecordType, RecordType)
func (r *JSTypeRegistry) RegisterDroppedPropertiesInUnion(subtype *RecordType, supertype *RecordType) {
	if r.parallelReadOnly {
		r.parallelMu.Lock()
		defer r.parallelMu.Unlock()
	}
	foundDroppedProperty := false
	subMap := subtype.GetPropertyMap()
	subMap.ensureSorted()
	for i := 0; i < len(subMap.entries()); i++ {
		propID := subMap.entries()[i].nameID
		if !HasPropertyByID(supertype, uint32(propID)) {
			foundDroppedProperty = true
			if r.droppedPropertiesOfUnions == nil {
				r.droppedPropertiesOfUnions = make(map[rhinoStringID]struct{})
			}
			r.droppedPropertiesOfUnions[propID] = struct{}{}
		}
	}
	if foundDroppedProperty {
		superMap := supertype.GetPropertyMap()
		superMap.ensureSorted()
		if len(superMap.entries()) > 0 && r.propertiesOfSupertypesInUnions == nil {
			r.propertiesOfSupertypesInUnions = make(map[rhinoStringID]struct{}, len(superMap.entries()))
		}
		for i := 0; i < len(superMap.entries()); i++ {
			r.propertiesOfSupertypesInUnions[superMap.entries()[i].nameID] = struct{}{}
		}
	}
}

// RegisterPropertyOnType tells the type system that owner may have a
// property named propertyName. This allows the registry to keep track of
// what types a property is defined upon.
//
// This is NOT the same as saying that owner must have a property named type.
// ObjectType#hasProperty attempts to minimize false positives ("if we're not
// sure, then don't type check this property"). The type registry, on the
// other hand, should attempt to minimize false negatives ("if this property
// is assigned anywhere in the program, it must show up in the type
// registry").
//
// Only string property keys are registered, not symbol property keys. This
// is largely for historical reasons - JSCompiler has allowed accessing
// properties on a struct via any symbol type, instead of reporting missing
// properties; and even if we do implement missing properties for symbol prop
// accesses, this method is not needed for "STRICT_MISSING_PROPERTIES" errors
// which we might default all symbol props to.
//
// Java: JSTypeRegistry.registerPropertyOnType(String, JSType)
func (r *JSTypeRegistry) RegisterPropertyOnType(propertyName string, type_ JSType) {
	r.registerPropertyOnTypeByID(rhinoStringPoolIntern(propertyName), type_)
}

func (r *JSTypeRegistry) registerPropertyOnTypeByID(id rhinoStringID, type_ JSType) {
	if type_.IsUnionType() {
		for _, alternate := range type_.ToMaybeUnionType().GetAlternates() {
			r.registerPropertyOnTypeByID(id, alternate)
		}
		return
	}

	if JSTypeRegistry_IsObjectLiteralThatCanBeSkipped(type_) {
		tb := r.GetSentinelObjectLiteral().GetJSTypeBase()
		s := r.nonRefTypesIndexedByProperty[id]
		if s.addWithHash(tb, r.sentinelObjectLiteralHash) {
			r.nonRefTypesIndexedByProperty[id] = s
		}
		return
	}

	tb := type_.GetJSTypeBase()
	if objType, ok := type_.(ObjectType); ok && objType.HasReferenceName() {
		s := r.eachRefTypeIndexedByProperty[id]
		if s.add(tb) {
			r.eachRefTypeIndexedByProperty[id] = s
		}
	} else {
		s := r.nonRefTypesIndexedByProperty[id]
		if s.add(tb) {
			r.nonRefTypesIndexedByProperty[id] = s
		}
	}
}

// JSTypeRegistry_PropDefinitionKind is a tristate value returned from
// canPropertyBeDefined.
type JSTypeRegistry_PropDefinitionKind int32

const (
	JSTypeRegistry_PropDefinitionKind_UNKNOWN     JSTypeRegistry_PropDefinitionKind = iota // The property is not known to be part of this type
	JSTypeRegistry_PropDefinitionKind_KNOWN                                                // The properties is known to be defined on a type or its super types
	JSTypeRegistry_PropDefinitionKind_LOOSE                                                // The property is loosely associated with a type, typically one of its subtypes
	JSTypeRegistry_PropDefinitionKind_LOOSE_UNION                                          // The property is loosely associated with a union type
)

var jsTypeRegistry_PropDefinitionKind_names = [...]string{"UNKNOWN", "KNOWN", "LOOSE", "LOOSE_UNION"}

// String is Enum.name().
func (k JSTypeRegistry_PropDefinitionKind) String() string {
	return jsTypeRegistry_PropDefinitionKind_names[k]
}

// CanPropertyBeDefined returns whether the given property can possibly be
// set on the given type.
//
// Java: JSTypeRegistry.canPropertyBeDefined(JSType, String)
func (r *JSTypeRegistry) CanPropertyBeDefined(type_ JSType, propertyName string) JSTypeRegistry_PropDefinitionKind {
	return r.canPropertyBeDefinedWithKey(
		type_,
		Property_Key{kind: Property_KeyKind_STRING, stringID: rhinoStringPoolLookup(propertyName), string_: propertyName})
}

// CanPropertyBeDefinedByID returns whether the property identified by propID
// can possibly be set on the given type without re-looking up propID in the
// string pool.
func (r *JSTypeRegistry) CanPropertyBeDefinedByID(type_ JSType, propID uint32) JSTypeRegistry_PropDefinitionKind {
	return r.canPropertyBeDefinedWithKey(type_, NewProperty_StringKeyByID(propID))
}

func isNonSyntheticSubtypeMatch(alternative JSType, type_ JSType) bool {
	greatestSubtype := alternative.GetGreatestSubtype(type_)
	if greatestSubtype.IsEmptyType() {
		return false
	}
	// We've found a type with this property. Now we just have to make
	// sure it's not a type used for internal bookkeeping.
	if maybeRecordType := greatestSubtype.ToMaybeRecordType(); maybeRecordType != nil && maybeRecordType.IsSynthetic() {
		return false
	}
	return true
}

func hasLooseSubtypeInPropSet(s propTypeSet, type_ JSType) bool {
	if s.t0 == nil {
		return false
	}
	if isNonSyntheticSubtypeMatch(s.t0.self_, type_) {
		return true
	}
	if s.t1 == nil {
		return false
	}
	if isNonSyntheticSubtypeMatch(s.t1.self_, type_) {
		return true
	}
	if s.overflow != nil {
		for _, slot := range s.overflow.slots {
			if isNonSyntheticSubtypeMatch(slot.t.self_, type_) {
				return true
			}
		}
	}
	return false
}

func (r *JSTypeRegistry) canPropertyBeDefinedWithKey(type_ JSType, key Property_Key) JSTypeRegistry_PropDefinitionKind {
	if type_.IsStruct() {
		// We are stricter about "struct" types and only allow access to
		// properties that to the best of our knowledge are available at
		// creation time and specifically not properties only defined on
		// subtypes.
		switch type_.GetPropertyKindByKey(key) {
		case JSType_HasPropertyKind_KNOWN_PRESENT:
			return JSTypeRegistry_PropDefinitionKind_KNOWN
		case JSType_HasPropertyKind_MAYBE_PRESENT:
			return JSTypeRegistry_PropDefinitionKind_LOOSE_UNION
		}
		return JSTypeRegistry_PropDefinitionKind_UNKNOWN // ABSENT
	}
	if !type_.IsEmptyType() && !type_.IsUnknownType() {
		switch type_.GetPropertyKindByKey(key) {
		case JSType_HasPropertyKind_KNOWN_PRESENT:
			return JSTypeRegistry_PropDefinitionKind_KNOWN
		case JSType_HasPropertyKind_MAYBE_PRESENT:
			return JSTypeRegistry_PropDefinitionKind_LOOSE_UNION
		case JSType_HasPropertyKind_ABSENT:
			// check for loose properties below.
		}
	}

	id := key.stringID
	if id == 0 {
		return JSTypeRegistry_PropDefinitionKind_UNKNOWN
	}

	if s, ok := r.nonRefTypesIndexedByProperty[id]; ok && hasLooseSubtypeInPropSet(s, type_) {
		return JSTypeRegistry_PropDefinitionKind_LOOSE
	}
	if s, ok := r.eachRefTypeIndexedByProperty[id]; ok && hasLooseSubtypeInPropSet(s, type_) {
		return JSTypeRegistry_PropDefinitionKind_LOOSE
	}

	if rec := type_.ToMaybeRecordType(); rec != nil {
		if r.parallelReadOnly {
			r.parallelMu.Lock()
		}
		loose := false
		if _, ok := r.droppedPropertiesOfUnions[id]; ok {
			propMap := rec.GetPropertyMap()
			for i := 0; i < len(propMap.entries()); i++ {
				if _, ok := r.propertiesOfSupertypesInUnions[propMap.entries()[i].nameID]; ok {
					loose = true
					break
				}
			}
		}
		if r.parallelReadOnly {
			r.parallelMu.Unlock()
		}
		if loose {
			return JSTypeRegistry_PropDefinitionKind_LOOSE
		}
	}
	return JSTypeRegistry_PropDefinitionKind_UNKNOWN
}

// GetEachReferenceTypeWithProperty returns each reference type that has a
// property propertyName defined on it.
//
// Unlike most types in our type system, the collection of types returned
// will not be collapsed. This means that if a type is defined on Object and
// on Array, this method must return [Object, Array]. It would not be correct
// to collapse them to [Object].
//
// Java: JSTypeRegistry.getEachReferenceTypeWithProperty(String)
func (r *JSTypeRegistry) GetEachReferenceTypeWithProperty(propertyName string) []ObjectType {
	if id := rhinoStringPoolLookup(propertyName); id != 0 {
		if s, ok := r.eachRefTypeIndexedByProperty[id]; ok {
			return s.toObjectTypeSlice()
		}
	}
	return []ObjectType{}
}

func (r *JSTypeRegistry) forEachRefTypeWithPropertyByID(id rhinoStringID, fn func(ObjectType)) {
	if id != 0 {
		if s, ok := r.eachRefTypeIndexedByProperty[id]; ok {
			s.forEachObjectType(fn)
		}
	}
}

// FindCommonSuperObject finds the common supertype of the two given object
// types.
//
// Java: JSTypeRegistry.findCommonSuperObject(ObjectType, ObjectType)
func (r *JSTypeRegistry) FindCommonSuperObject(a ObjectType, b ObjectType) ObjectType {
	var bufA, bufB [16]ObjectType
	stackA := bufA[:0]
	for current := a; current != nil; current = current.GetImplicitPrototype() {
		stackA = append(stackA, current)
	}
	stackB := bufB[:0]
	for current := b; current != nil; current = current.GetImplicitPrototype() {
		stackB = append(stackB, current)
	}

	result := r.GetNativeObjectType(JSTypeNative_OBJECT_TYPE)
	for len(stackA) != 0 && len(stackB) != 0 {
		currentA := stackA[len(stackA)-1]
		stackA = stackA[:len(stackA)-1]
		currentB := stackB[len(stackB)-1]
		stackB = stackB[:len(stackB)-1]
		if currentA.Equals(currentB) {
			result = currentA
		} else {
			return result
		}
	}
	return result
}

// Java: JSTypeRegistry.getSuperStack(ObjectType)
func JSTypeRegistry_GetSuperStack(a ObjectType) []ObjectType {
	stack := make([]ObjectType, 0, 5)
	for current := a; current != nil; current = current.GetImplicitPrototype() {
		stack = append(stack, current)
	}
	return stack
}
