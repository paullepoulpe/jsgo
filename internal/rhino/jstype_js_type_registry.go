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
	"strings"
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
		slots := s.overflow.slots
		for _, slot := range slots {
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

// Java: JSTypeRegistry.DOT_SPLITTER (Splitter.on('.'))
func JSTypeRegistry_DOT_SPLITTER_Split(s string) []string {
	n := 1
	for i := 0; i < len(s); i++ {
		if s[i] == '.' {
			n++
		}
	}
	out := make([]string, 0, n)
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '.' {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	return append(out, s[start:])
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

// GetArrayElementKey returns the template variable for the element type of
// Arrays.
//
// Java: JSTypeRegistry.getArrayElementKey()
func (r *JSTypeRegistry) GetArrayElementKey() *TemplateType {
	return r.arrayElementTemplateKey
}

// GetReadonlyArrayElementKey returns the template variable for the element
// type of ReadonlyArrays.
//
// Java: JSTypeRegistry.getReadonlyArrayElementKey()
func (r *JSTypeRegistry) GetReadonlyArrayElementKey() *TemplateType {
	return r.readonlyArrayElementTemplateKey
}

// GetReadonlyMapKey returns the template variable for the key type of
// ReadonlyMaps.
//
// Java: JSTypeRegistry.getReadonlyMapKey()
func (r *JSTypeRegistry) GetReadonlyMapKey() *TemplateType {
	return r.readonlyMapKeyTemplateKey
}

// GetReadonlyMapValue returns the template variable for the value type of
// ReadonlyMaps.
//
// Java: JSTypeRegistry.getReadonlyMapValue()
func (r *JSTypeRegistry) GetReadonlyMapValue() *TemplateType {
	return r.readonlyMapValueTemplateKey
}

// GetObjectElementKey returns the template variable corresponding to the
// property value type for Javascript Objects and Arrays.
//
// Java: JSTypeRegistry.getObjectElementKey()
func (r *JSTypeRegistry) GetObjectElementKey() *TemplateType {
	return r.iObjectElementTemplateKey
}

// GetObjectIndexKey returns the template variable corresponding to the
// property key type of the built-in Javascript object.
//
// Java: JSTypeRegistry.getObjectIndexKey()
func (r *JSTypeRegistry) GetObjectIndexKey() *TemplateType {
	jutil.CheckNotNull(r.iObjectIndexTemplateKey)
	return r.iObjectIndexTemplateKey
}

// GetIterableValueTemplate returns the value template variable for the
// Iterable interface.
//
// Java: JSTypeRegistry.getIterableValueTemplate()
func (r *JSTypeRegistry) GetIterableValueTemplate() *TemplateType {
	return jutil.CheckNotNull(r.iterableValueTemplate)
}

// GetIterableReturnTemplate returns the return template variable for the
// Iterable interface.
//
// Java: JSTypeRegistry.getIterableReturnTemplate()
func (r *JSTypeRegistry) GetIterableReturnTemplate() *TemplateType {
	return jutil.CheckNotNull(r.iterableReturnTemplate)
}

// GetIteratorIterableValueTemplate returns the value template variable for
// the IteratorIterable interface.
//
// Java: JSTypeRegistry.getIteratorIterableValueTemplate()
func (r *JSTypeRegistry) GetIteratorIterableValueTemplate() *TemplateType {
	return jutil.CheckNotNull(r.iteratorIterableValueTemplate)
}

// GetIteratorIterableReturnTemplate returns the return template variable
// for the IteratorIterable interface.
//
// Java: JSTypeRegistry.getIteratorIterableReturnTemplate()
func (r *JSTypeRegistry) GetIteratorIterableReturnTemplate() *TemplateType {
	return jutil.CheckNotNull(r.iteratorIterableReturnTemplate)
}

// GetIIterableResultValueTemplate returns the value template variable for
// the IIterableResult interface.
//
// Java: JSTypeRegistry.getIIterableResultValueTemplate()
func (r *JSTypeRegistry) GetIIterableResultValueTemplate() *TemplateType {
	return jutil.CheckNotNull(r.iiterableResultValueTemplate)
}

// GetIteratorLikeValueTemplate returns the value template variable for the
// IteratorLike interface.
//
// Java: JSTypeRegistry.getIteratorLikeValueTemplate()
func (r *JSTypeRegistry) GetIteratorLikeValueTemplate() *TemplateType {
	return jutil.CheckNotNull(r.iteratorLikeValueTemplate)
}

// GetIteratorLikeReturnTemplate returns the return template variable for the
// IteratorLike interface.
//
// Java: JSTypeRegistry.getIteratorLikeReturnTemplate()
func (r *JSTypeRegistry) GetIteratorLikeReturnTemplate() *TemplateType {
	return jutil.CheckNotNull(r.iteratorLikeReturnTemplate)
}

// GetIteratorValueTemplate returns the value template variable for the
// Iterator interface.
//
// Java: JSTypeRegistry.getIteratorValueTemplate()
func (r *JSTypeRegistry) GetIteratorValueTemplate() *TemplateType {
	return jutil.CheckNotNull(r.iteratorValueTemplate)
}

// GetIteratorReturnTemplate returns the return template variable for the
// Iterator interface.
//
// Java: JSTypeRegistry.getIteratorReturnTemplate()
func (r *JSTypeRegistry) GetIteratorReturnTemplate() *TemplateType {
	return jutil.CheckNotNull(r.iteratorReturnTemplate)
}

// GetGeneratorValueTemplate returns the value template variable for the
// Generator interface.
//
// Java: JSTypeRegistry.getGeneratorValueTemplate()
func (r *JSTypeRegistry) GetGeneratorValueTemplate() *TemplateType {
	return jutil.CheckNotNull(r.generatorValueTemplate)
}

// GetGeneratorReturnTemplate returns the return template variable for the
// Generator interface.
//
// Java: JSTypeRegistry.getGeneratorReturnTemplate()
func (r *JSTypeRegistry) GetGeneratorReturnTemplate() *TemplateType {
	return jutil.CheckNotNull(r.generatorReturnTemplate)
}

// GetAsyncIterableValueTemplate returns the template variable for the
// AsyncIterable interface.
//
// Java: JSTypeRegistry.getAsyncIterableValueTemplate()
func (r *JSTypeRegistry) GetAsyncIterableValueTemplate() *TemplateType {
	return jutil.CheckNotNull(r.asyncIterableValueTemplate)
}

// GetAsyncIterableReturnTemplate returns the template variable for the
// AsyncIterableR interface.
//
// Java: JSTypeRegistry.getAsyncIterableReturnTemplate()
func (r *JSTypeRegistry) GetAsyncIterableReturnTemplate() *TemplateType {
	return jutil.CheckNotNull(r.asyncIterableReturnTemplate)
}

// GetAsyncIteratorValueTemplate returns the template variable for the
// AsyncIterator interface.
//
// Java: JSTypeRegistry.getAsyncIteratorValueTemplate()
func (r *JSTypeRegistry) GetAsyncIteratorValueTemplate() *TemplateType {
	return jutil.CheckNotNull(r.asyncIteratorValueTemplate)
}

// GetAsyncIteratorReturnTemplate returns the template variable for the
// AsyncIterator interface.
//
// Java: JSTypeRegistry.getAsyncIteratorReturnTemplate()
func (r *JSTypeRegistry) GetAsyncIteratorReturnTemplate() *TemplateType {
	return jutil.CheckNotNull(r.asyncIteratorReturnTemplate)
}

// GetIThenableTemplate returns the template variable for the IThenable
// interface.
//
// Java: JSTypeRegistry.getIThenableTemplate()
func (r *JSTypeRegistry) GetIThenableTemplate() *TemplateType {
	return jutil.CheckNotNull(r.iThenableTemplateKey)
}

// MaybeGetTemplateTypesOfBuiltin returns an immutable list of template types
// of the given builtin.
//
// Java: JSTypeRegistry.maybeGetTemplateTypesOfBuiltin(StaticScope, String)
func (r *JSTypeRegistry) MaybeGetTemplateTypesOfBuiltin(scope StaticScope, fnName string) []*TemplateType {
	type_ := r.GetType(scope, fnName)
	var objType ObjectType
	if type_ != nil {
		objType = type_.ToObjectType()
	}
	if objType != nil && objType.IsNativeObjectType() {
		return objType.GetTypeParameters()
	}
	return nil
}

// Java: JSTypeRegistry.getErrorReporter()
func (r *JSTypeRegistry) GetErrorReporter() ErrorReporter {
	return r.reporter
}

// Java: JSTypeRegistry.getRootElementOfName(String)
func JSTypeRegistry_GetRootElementOfName(name string) string {
	index := strings.IndexByte(name, '.')
	if index != -1 {
		return name[:index]
	}
	return name
}

// JSTypeRegistry_GetLookupScope returns which scope in the provided scope
// chain the provided name is declared in, or else null.
//
// This assumed that the Scope construction is complete. It can not be used
// during scope construction to determine if a name is already defined as a
// shadowed name from a parent scope would be returned.
//
// Java: JSTypeRegistry.getLookupScope(StaticScope, String)
func JSTypeRegistry_GetLookupScope(scope StaticScope, name string) StaticScope {
	if scope != nil && scope.GetParentScope() != nil {
		return scope.GetTopmostScopeOfEventualDeclaration(JSTypeRegistry_GetRootElementOfName(name))
	}
	return scope
}

// Java: JSTypeRegistry.getRootNodeForScope(StaticScope)
func (r *JSTypeRegistry) GetRootNodeForScope(scope StaticScope) Node {
	var root Node
	if scope != nil {
		root = scope.GetRootNode()
	}
	if root.IsNil() || root.IsRoot() || root.IsScript() {
		return r.nameTableGlobalRoot
	}
	return root
}

// Java: JSTypeRegistry.isDeclaredForScope(StaticScope, String)
func (r *JSTypeRegistry) IsDeclaredForScope(scope StaticScope, name string) bool {
	return r.GetTypeInternal(scope, name) != nil
}

// Java: JSTypeRegistry.checkTypeName(String)
func JSTypeRegistry_CheckTypeName(typeName string) {
	if strings.IndexByte(typeName, '<') >= 0 {
		panic(jutil.NewIllegalArgumentExceptionMsg("Type names cannot contain template annotations."))
	}
}

// Java: JSTypeRegistry.getTypeInternal(StaticScope, String)
func (r *JSTypeRegistry) GetTypeInternal(scope StaticScope, name string) JSType {
	JSTypeRegistry_CheckTypeName(name)
	if syntheticTemplateScope, ok := scope.(*JSTypeRegistry_SyntheticTemplateScope); ok {
		type_ := syntheticTemplateScope.GetTemplateType(name)
		if type_ != nil {
			return type_
		}
	}
	declarationScope := JSTypeRegistry_GetLookupScope(scope, name)
	resolvedViaTable := r.GetTypeForScopeInternal(declarationScope, name)
	if resolvedViaTable != nil {
		return resolvedViaTable
	}
	bestScope := scope
	if declarationScope != nil {
		bestScope = declarationScope
	}
	return r.ResolveViaComponents(bestScope, name)
}

// Java: JSTypeRegistry.resolveViaComponents(StaticScope, String)
func (r *JSTypeRegistry) ResolveViaComponents(scope StaticScope, qualifiedName string) JSType {
	resolutionScope, ok := scope.(StaticTypedScope)
	if qualifiedName == "" || !ok {
		return nil
	}
	// Skip closure namespace resolution of types whose root component is
	// defined in a local scope (not a global scope). Those will follow the
	// normal resolution scheme. (For legacy compatibility reasons we don't
	// check for global names that are the same as the module root).
	if !JSTypeRegistry_IsNameDefinedLocally(resolutionScope, JSTypeRegistry_GetRootElementOfName(qualifiedName)) {
		resolvedViaClosureNamespace := r.ResolveViaClosureNamespace(qualifiedName)
		if resolvedViaClosureNamespace != nil {
			return resolvedViaClosureNamespace
		}
	}
	return r.ResolveViaProperties(resolutionScope, qualifiedName)
}

// Java: JSTypeRegistry.isNameDefinedLocally(StaticTypedScope, String)
func JSTypeRegistry_IsNameDefinedLocally(resolutionScope StaticTypedScope, reference string) bool {
	slot := resolutionScope.GetSlot(reference)
	return slot != nil && slot.GetScope() != nil && slot.GetScope().GetParentScope() != nil
}

// ResolveViaClosureNamespace resolves a named type by checking for the
// longest prefix that matches some Closure namespace, if any, then
// attempting to resolve via properties based on the type of the `exports`
// object in that namespace.
//
// Java: JSTypeRegistry.resolveViaClosureNamespace(String)
func (r *JSTypeRegistry) ResolveViaClosureNamespace(reference string) JSType {
	if r.closureNamespaces.IsEmpty() {
		return nil
	}
	// Find the `exports` type of the longest prefix match of this namespace,
	// if any. Then resolve it via property.
	prefix := reference
	for {
		namespace := r.closureNamespaces.Get(prefix)
		if namespace != nil {
			if namespace.IsLegacy() {
				// Try to resolve this name via registry or properties.
				return nil
			}
			// Always stop resolution here whether successful or not, instead of
			// continuing with resolution via registry or via properties, to
			// match legacy behavior.
			var components []string
			if len(prefix) < len(reference) {
				components = JSTypeRegistry_DOT_SPLITTER_Split(reference[len(prefix)+1:])
			}
			return r.ResolveViaPropertyGivenSlot(
				namespace.Type(), namespace.DefinitionNode(), components)
		}

		lastDot := strings.LastIndexByte(prefix, '.')
		if lastDot < 0 {
			return nil
		}
		prefix = prefix[:lastDot]
	}
}

// Java: JSTypeRegistry.resolveViaProperties(StaticTypedScope, String)
func (r *JSTypeRegistry) ResolveViaProperties(declarationScope StaticTypedScope, qualifiedName string) JSType {
	jutil.CheckArgument(qualifiedName != "")
	dotIdx := strings.IndexByte(qualifiedName, '.')
	rootName := qualifiedName
	if dotIdx != -1 {
		rootName = qualifiedName[:dotIdx]
	}
	slot, _ := declarationScope.GetOwnSlot(rootName).(StaticTypedSlot)
	if slot == nil {
		return nil
	}
	var componentNames []string
	var componentBuf [16]string
	if dotIdx != -1 {
		rem := qualifiedName[dotIdx+1:]
		dots := 0
		for i := 0; i < len(rem); i++ {
			if rem[i] == '.' {
				dots++
			}
		}
		if dots < len(componentBuf) {
			componentNames = componentBuf[:0]
			start := 0
			for i := 0; i < len(rem); i++ {
				if rem[i] == '.' {
					componentNames = append(componentNames, rem[start:i])
					start = i + 1
				}
			}
			componentNames = append(componentNames, rem[start:])
		} else {
			componentNames = JSTypeRegistry_DOT_SPLITTER_Split(rem)
		}
	}
	slotType := slot.GetType()
	return r.ResolveViaPropertyGivenSlot(slotType, NilNode, componentNames)
}

// ResolveViaPropertyGivenSlot resolves a type using a given StaticTypedSlot
// and list of properties on that type.
//
// slotType: the JSType of the slot, possibly null. definitionNode: if known,
// the Node representing the type definition.
//
// Java: JSTypeRegistry.resolveViaPropertyGivenSlot(JSType, Node, List<String>)
func (r *JSTypeRegistry) ResolveViaPropertyGivenSlot(slotType JSType, definitionNode Node, componentNames []string) JSType {
	if len(componentNames) == 0 {
		typedefType := JSTypeRegistry_ResolveTypeFromNodeIfTypedef(definitionNode)
		if typedefType != nil {
			return typedefType
		}
	}

	// If the first component has a type of 'Unknown', then any type names
	// using it should be regarded as silently 'Unknown' rather than be noisy
	// about it.
	if slotType == nil || slotType.IsAllType() || slotType.IsNoType() {
		return nil
	}

	// resolving component by component
	for i := 0; i < len(componentNames); i++ {
		component := componentNames[i]
		parentObj := ObjectType_Cast(slotType)
		if parentObj == nil || component == "" || !parentObj.HasOwnProperty(component) {
			return nil
		}
		if i == len(componentNames)-1 {
			// Look for a typedefTypeProp on the definition node of the last
			// component.
			def := parentObj.GetPropertyDefSite(component)
			typedefType := JSTypeRegistry_ResolveTypeFromNodeIfTypedef(def)
			if typedefType != nil {
				return typedefType
			}
		}
		slotType = parentObj.GetPropertyType(component)
	}

	// Translate "constructor" types to "instance" types.
	if slotType == nil {
		return nil
	} else if slotType.IsFunctionType() && (slotType.IsConstructor() || slotType.IsInterface()) {
		return jsTypeOrNil(slotType.ToMaybeFunctionType().GetInstanceType())
	} else if slotType.IsNoObjectType() {
		return r.GetNativeObjectType(JSTypeNative_NO_OBJECT_TYPE)
	} else if enumType, ok := slotType.(*EnumType); ok {
		return enumType.GetElementsType()
	} else {
		return nil
	}
}

// JSTypeRegistry_ResolveTypeFromNodeIfTypedef checks the given Node for a
// typedef annotation, resolving to that type if existent.
//
// Java: JSTypeRegistry.resolveTypeFromNodeIfTypedef(Node)
func JSTypeRegistry_ResolveTypeFromNodeIfTypedef(node Node) JSType {
	if node.IsNil() {
		return nil
	}
	return node.GetTypedefTypeProp()
}

// Java: JSTypeRegistry.getTypeForScopeInternal(StaticScope, String)
func (r *JSTypeRegistry) GetTypeForScopeInternal(scope StaticScope, name string) JSType {
	rootNode := r.GetRootNodeForScope(scope)
	row := r.scopedNameTable.Get(rootNode)
	if row == nil {
		return nil
	}
	return row.Get(name)
}

// Java: JSTypeRegistry.registerGlobalType(JSType)
func (r *JSTypeRegistry) RegisterGlobalType(type_ JSType) {
	r.Register(nil, type_, type_.String())
}

// Java: JSTypeRegistry.registerGlobalType(JSType, String)
func (r *JSTypeRegistry) RegisterGlobalTypeWithName(type_ JSType, name string) {
	r.Register(nil, type_, name)
}

// Java: JSTypeRegistry.register(StaticScope, JSType, String)
func (r *JSTypeRegistry) Register(scope StaticScope, type_ JSType, name string) {
	JSTypeRegistry_CheckTypeName(name)
	r.RegisterForScope(JSTypeRegistry_GetLookupScope(scope, name), type_, name)
}

// Java: JSTypeRegistry.registerForScope(StaticScope, JSType, String)
func (r *JSTypeRegistry) RegisterForScope(scope StaticScope, type_ JSType, name string) {
	r.scopedNameTablePut(r.GetRootNodeForScope(scope), name, type_)
}

// scopedNameTablePut is Table.put(row, column, value), which rejects nulls.
func (r *JSTypeRegistry) scopedNameTablePut(rowKey Node, columnKey string, value JSType) {
	CheckNotNull(rowKey)
	jutil.CheckNotNull(value)
	row := r.scopedNameTable.Get(rowKey)
	if row == nil {
		row = jutil.NewLinkedHashMap[string, JSType]()
		r.scopedNameTable.Put(rowKey, row)
	}
	row.Put(columnKey, value)
}

// Java: JSTypeRegistry.registerNativeType(JSTypeNative, JSType)
func (r *JSTypeRegistry) RegisterNativeType(typeId JSTypeNative, type_ JSType) {
	r.nativeTypes[typeId.Ordinal()] = type_
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

// DeclareType records declared global type names. This makes resolution
// faster and more robust in the common case.
//
// name: the name of the type to be recorded. type: the actual type being
// associated with the name. Returns true if this name is not already
// defined, false otherwise.
//
// Java: JSTypeRegistry.declareType(StaticScope, String, JSType)
func (r *JSTypeRegistry) DeclareType(scope StaticScope, name string, type_ JSType) bool {
	jutil.CheckState(name != "")
	if r.GetTypeForScopeInternal(JSTypeRegistry_GetLookupScope(scope, name), name) != nil {
		return false
	}

	r.Register(scope, type_, name)
	return true
}

// DeclareTypeForExactScope records declared global type names. This makes
// resolution faster and more robust in the common case.
//
// Java: JSTypeRegistry.declareTypeForExactScope(StaticScope, String, JSType)
func (r *JSTypeRegistry) DeclareTypeForExactScope(scope StaticScope, name string, type_ JSType) bool {
	jutil.CheckState(name != "")
	if r.GetTypeForScopeInternal(scope, name) != nil {
		return false
	}

	r.RegisterForScope(scope, type_, name)
	return true
}

// OverwriteDeclaredType overrides a declared global type name. Throws an
// exception if this type name hasn't been declared yet.
//
// Java: JSTypeRegistry.overwriteDeclaredType(StaticScope, String, JSType)
func (r *JSTypeRegistry) OverwriteDeclaredType(scope StaticScope, name string, type_ JSType) {
	jutil.CheckStateMsg1(r.IsDeclaredForScope(scope, name), "missing name %s", name)
	r.Register(scope, type_, name)
}

// IsForwardDeclaredType reports whether this is a forward-declared type
// name.
//
// Java: JSTypeRegistry.isForwardDeclaredType(String)
func (r *JSTypeRegistry) IsForwardDeclaredType(name string) bool {
	return r.forwardDeclaredTypes.Contains(name)
}

// GetReadableTypeName first dereferences the JSType to remove
// null/undefined then returns a human-readable type name
//
// Java: JSTypeRegistry.getReadableTypeName(Node)
func (r *JSTypeRegistry) GetReadableTypeName(n Node) string {
	return r.GetReadableJSTypeName(n, true)
}

// Java: JSTypeRegistry.getReadableTypeNameNoDeref(Node)
func (r *JSTypeRegistry) GetReadableTypeNameNoDeref(n Node) string {
	return r.GetReadableJSTypeName(n, false)
}

// Java: JSTypeRegistry.getSimpleReadableJSTypeName(JSType)
func JSTypeRegistry_GetSimpleReadableJSTypeName(type_ JSType) *string {
	if _, ok := type_.(*AllType); ok {
		s := type_.String()
		return &s
	} else if _, ok := type_.(ValueType); ok {
		s := type_.String()
		return &s
	} else if type_.IsFunctionPrototypeType() {
		s := type_.String()
		return &s
	} else if _, ok := type_.(ObjectType); ok {
		if !JSTypeRegistry_IsReadableObjectType(type_.ToMaybeObjectType()) {
			return nil
		}

		source := type_.ToObjectType().GetConstructor().GetSource()
		if !source.IsNil() {
			jutil.CheckStateMsg1(source.IsFunction() || source.IsClass(), "%s", source)
			readable := source.GetFirstChild().GetOriginalName()
			if readable != nil {
				return readable
			}
		}
		s := type_.String()
		return &s

	} else if _, ok := type_.(*UnionType); ok {
		unionType := type_.ToMaybeUnionType()
		var union *string
		for _, alternate := range unionType.GetAlternates() {
			name := JSTypeRegistry_GetSimpleReadableJSTypeName(alternate)
			if name == nil {
				return nil
			}
			if union == nil {
				u := "(" + *name
				union = &u
			} else {
				u := *union + "|" + *name
				union = &u
			}
		}
		u := jutil.ToJavaString(union) + ")"
		return &u
	}
	return nil
}

// GetReadableJSTypeName: given a node, get a human-readable name for the
// type of that node so that will be easy for the programmer to find the
// original declaration.
//
// For example, if SubFoo's property "bar" might have the human-readable name
// "Foo.prototype.bar".
//
// n: the node. dereference: if true, the type of the node will be
// dereferenced to an Object type, if possible. Prefer to call
// #getReadableTypeName(String) or #getReadableTypeNameNoDeref(String)
// instead of passing this as an argument.
//
// Java: JSTypeRegistry.getReadableJSTypeName(Node, boolean)
func (r *JSTypeRegistry) GetReadableJSTypeName(n Node, dereference bool) string {
	type_ := r.GetJSTypeOrUnknown(n)
	if dereference {
		autoboxed := type_.Autobox()
		if !autoboxed.IsNoType() {
			type_ = autoboxed
		}
	}

	name := JSTypeRegistry_GetSimpleReadableJSTypeName(type_)
	if name != nil {
		return *name
	}

	// If we're analyzing a GETPROP, the property may be inherited by the
	// prototype chain. So climb the prototype chain and find out where the
	// property was originally defined.
	if n.IsGetProp() {
		objectType := r.GetJSTypeOrUnknown(n.GetFirstChild()).Dereference()
		if objectType != nil {
			propName := n.GetString()
			objectType = objectType.GetClosestDefiningType(propName)

			// Don't show complex function names or anonymous types. Instead, try to
			// get a human-readable type name.
			if JSTypeRegistry_IsReadableObjectType(objectType) {
				return objectType.String() + "." + propName
			}
		}
	}

	if n.IsQualifiedName() {
		return *n.GetQualifiedName()
	} else if type_.IsFunctionType() {
		// Don't show complex function names.
		return "function"
	} else {
		return type_.String()
	}
}

// Java: JSTypeRegistry.isReadableObjectType(ObjectType)
func JSTypeRegistry_IsReadableObjectType(type_ ObjectType) bool {
	if type_ == nil {
		return false
	} else if type_.IsInstanceType() {
		return true
	} else if type_.IsFunctionPrototypeType() {
		return true
	} else if type_.IsEnumElementType() {
		// TODO(b/147236174): This case should be deleted as impossible.
		return type_.ToMaybeEnumElementType().GetConstructor() != nil
	}
	return false
}

// Java: JSTypeRegistry.getJSTypeOrUnknown(Node)
func (r *JSTypeRegistry) GetJSTypeOrUnknown(n Node) JSType {
	jsType := n.GetJSType()
	if jsType == nil {
		return r.GetNativeType(JSTypeNative_UNKNOWN_TYPE)
	} else {
		return jsType
	}
}

// Java: JSTypeRegistry.getGlobalType(String)
func (r *JSTypeRegistry) GetGlobalType(jsTypeName string) JSType {
	return r.GetType(nil, jsTypeName)
}

// GetType looks up a native type by name. Returns the corresponding JSType
// object or null it cannot be found.
//
// Java: JSTypeRegistry.getType(StaticScope, String)
func (r *JSTypeRegistry) GetType(scope StaticScope, jsTypeName string) JSType {
	return r.GetTypeInternal(scope, jsTypeName)
}

// GetTypeWithLocation looks up a type by name. To allow for forward
// references to types, an unrecognized string has to be bound to a
// NamedType object that will be resolved later.
//
// scope: a scope for doing type name resolution. jsTypeName: the name
// string. sourceName: the name of the source file where this reference
// appears. lineno: the line number of the reference. Returns a NamedType if
// the string argument is not one of the known types, otherwise the
// corresponding JSType object.
func (r *JSTypeRegistry) getTypeResolvedOrNil(scope StaticTypedScope, jsTypeName string) JSType {
	switch jsTypeName {
	case "boolean":
		return r.GetNativeType(JSTypeNative_BOOLEAN_TYPE)
	case "number":
		return r.GetNativeType(JSTypeNative_NUMBER_TYPE)
	case "bigint":
		return r.GetNativeType(JSTypeNative_BIGINT_TYPE)
	case "string":
		return r.GetNativeType(JSTypeNative_STRING_TYPE)
	case "undefined", "void":
		return r.GetNativeType(JSTypeNative_VOID_TYPE)
	}
	// Resolve template type names
	var thisType JSType
	if scope != nil && scope.GetTypeOfThis() != nil {
		thisType = jsTypeOrNil(scope.GetTypeOfThis().ToObjectType())
	}
	if thisType != nil {
		// get the 'last' template type key with this name because if a
		// superclass template key shares the same name as a subclass key, the
		// subclass key is guaranteed to be last.
		key := thisType.GetTemplateTypeMap().GetLastTemplateTypeKeyByName(jsTypeName)
		if key != nil {
			type_ := JSType(key)
			jutil.CheckStateMsg1(type_.IsTemplateType(), "expected:%s", type_)
			return type_
		}
	}

	// TODO(sdh): The use of "getType" here is incorrect. This currently will
	// pick up a type in an outer scope if it will be shadowed by a local
	// type. But creating a unique NamedType object for every name referenced
	// (even if interned) in every scope would be expensive.
	//
	// Instead perhaps a standard (untyped) scope object might be used to pick
	// the right scope during type construction.
	var untypedScope StaticScope
	if scope != nil {
		untypedScope = scope
	}
	return r.GetType(untypedScope, jsTypeName)
}

func (r *JSTypeRegistry) getTypeWithLocationString(
	scope StaticTypedScope, jsTypeName string, sourceName string, lineno int, charno int) JSType {
	if type_ := r.getTypeResolvedOrNil(scope, jsTypeName); type_ != nil {
		return type_
	}
	return r.createNamedTypeByID(scope, jsTypeName, rhinoStringPoolIntern(sourceName), lineno, charno)
}

// Java: JSTypeRegistry.getType(StaticTypedScope, String, String, int, int)
func (r *JSTypeRegistry) GetTypeWithLocation(
	scope StaticTypedScope, jsTypeName string, sourceName *string, lineno int, charno int) JSType {
	if type_ := r.getTypeResolvedOrNil(scope, jsTypeName); type_ != nil {
		return type_
	}
	// TODO(user): Each instance should support named type creation using
	// interning.
	return r.CreateNamedType(scope, jsTypeName, sourceName, lineno, charno)
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

// Java: JSTypeRegistry.evaluateTypeExpressionInGlobalScope(JSTypeExpression)
func (r *JSTypeRegistry) EvaluateTypeExpressionInGlobalScope(expr *JSTypeExpression) JSType {
	return expr.Evaluate(nil, r)
}

// RegisterNonLegacyClosureNamespace registers a goog.module namespace (that
// does not have goog.module.declareLegacyNamespace)
//
// This allows JSTypeRegistry to resolve types that refer to goog.modules by
// namespace. These have unique handling because they exist only in the type
// space and do not have a corresponding value space value.
//
// Java: JSTypeRegistry.registerNonLegacyClosureNamespace(String, Node, JSType)
func (r *JSTypeRegistry) RegisterNonLegacyClosureNamespace(moduleName string, definitionNode Node, type_ JSType) {
	r.closureNamespaces.Put(
		moduleName, JSTypeRegistry_ClosureNamespace_Create( /* isLegacy= */ false, definitionNode, type_))
}

// RegisterLegacyClosureNamespace registers a goog.provide or legacy
// goog.module namespace with the type registry
//
// Java: JSTypeRegistry.registerLegacyClosureNamespace(String)
func (r *JSTypeRegistry) RegisterLegacyClosureNamespace(moduleName string) {
	r.closureNamespaces.Put(moduleName, JSTypeRegistry_ClosureNamespace_Create( /* isLegacy= */ true, NilNode, nil))
}

// JSTypeRegistry_ClosureNamespace stores information about a Closure
// namespace.
//
// Java: abstract class JSTypeRegistry.ClosureNamespace, with its generated
// AutoValue_JSTypeRegistry_ClosureNamespace implementation folded in.
type JSTypeRegistry_ClosureNamespace struct {
	isLegacy       bool
	definitionNode Node
	type_          JSType
}

// IsLegacy returns true if this is a goog.provide'd namespace or a
// goog.module namespace followed by `goog.module.declareLegacyNamespace()
//
// Java: JSTypeRegistry.ClosureNamespace.isLegacy()
func (c *JSTypeRegistry_ClosureNamespace) IsLegacy() bool { return c.isLegacy }

// Java: JSTypeRegistry.ClosureNamespace.definitionNode()
func (c *JSTypeRegistry_ClosureNamespace) DefinitionNode() Node { return c.definitionNode }

// Java: JSTypeRegistry.ClosureNamespace.type()
func (c *JSTypeRegistry_ClosureNamespace) Type() JSType { return c.type_ }

// Java: JSTypeRegistry.ClosureNamespace.create(boolean, Node, JSType)
func JSTypeRegistry_ClosureNamespace_Create(isLegacy bool, definitionNode Node, type_ JSType) *JSTypeRegistry_ClosureNamespace {
	return &JSTypeRegistry_ClosureNamespace{isLegacy: isLegacy, definitionNode: definitionNode, type_: type_}
}

// Java: AutoValue_JSTypeRegistry_ClosureNamespace.toString()
func (c *JSTypeRegistry_ClosureNamespace) String() string {
	return "ClosureNamespace{" +
		"isLegacy=" + jutil.ToJavaString(c.isLegacy) + ", " +
		"definitionNode=" + jutil.ToJavaString(c.definitionNode) + ", " +
		"type=" + jutil.ToJavaString(c.type_) +
		"}"
}

// Java: AutoValue_JSTypeRegistry_ClosureNamespace.equals(Object)
func (c *JSTypeRegistry_ClosureNamespace) Equals(o any) bool {
	if any(c) == o {
		return true
	}
	if that, ok := o.(*JSTypeRegistry_ClosureNamespace); ok && that != nil {
		return c.isLegacy == that.IsLegacy() &&
			jutil.ObjectsEquals(c.definitionNode, that.DefinitionNode()) &&
			jutil.ObjectsEquals(c.type_, that.Type())
	}
	return false
}

// Java: AutoValue_JSTypeRegistry_ClosureNamespace.hashCode()
func (c *JSTypeRegistry_ClosureNamespace) HashCode() int32 {
	var h int32 = 1
	h *= 1000003
	h ^= jutil.ObjectHashCode(c.isLegacy)
	h *= 1000003
	h ^= jutil.ObjectHashCode(c.definitionNode)
	h *= 1000003
	h ^= jutil.ObjectHashCode(c.type_)
	return h
}
