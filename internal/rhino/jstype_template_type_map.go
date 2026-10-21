// Port of com/google/javascript/rhino/jstype/TemplateTypeMap.java @ closure-compiler v20261015

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// TemplateTypeMap manages a mapping from TemplateType to its resolved
// JSType. Provides utility methods for cloning/extending the map.
type TemplateTypeMap struct {
	// The TemplateType keys of the map.
	templateKeys []*TemplateType
	// The JSType values, which are index-aligned with their corresponding
	// keys. These values are left as specified in the TemplateTypeMap
	// constructor; they may refer to TemplateTypes that are keys in this
	// TemplateTypeMap, requiring iterative type resolution to find their true,
	// resolved type.
	templateValues []JSType
	// The JSType values, which are index-aligned with their corresponding
	// keys. These values have been iteratively type-resolved using this
	// TemplateTypeMap instance. These fully-resolved values are necessary for
	// determining the equivalence of two TemplateTypeMap instances.
	resolvedTemplateValues []JSType
	// A positive boolean indicates that the corresponding key is the start of
	// a new submap, which corresponds to separate individual ObjectTypes with
	// their own keys.
	//
	// For example, consider:
	//
	//	class Foo<T> {}
	//	class Bar<U, V, X> extends Foo<U | V> {}
	//	class Baz<Y> extends Bar<Y, number, string> {}
	//
	// The TemplateTypeMap for Baz will have:
	//
	//	templateKeys:   [  T, U,      V,      X, Y]
	//	templateValues: [U|V, Y, number, string, ?]
	//	subMapStarts:   [  1, 1,      0,      0, 1]
	//
	// (resolvedTemplateValues would be all ? at the start) We would say there
	// are three "submaps" in the Baz TemplateMap:
	//
	//  1. [T -> U], for Foo
	//  2. [U -> Y, V -> number, X -> string] for Bar
	//  3. [Y -> ?] for Baz
	//
	// In theory, we could model this with having multiple arrays of keys and
	// values - a single submap would have pointers to the parent maps for its
	// extended/implemented types, and those could further extend other types.
	// But in practice, it has historically seemed more useful to view these as
	// a single template key/value list pair.
	//
	// A TemplateType key may only be referenced by a template value if the
	// key is either 1. to the "right" of the value, i.e. at an index greater
	// than or equal to the value 2. in the same submap as the value.
	//
	// For example, given class Parent<T> {} and class Child<U> extends
	// Parent<(whatever is here)> {} the TemplateTypeMap for Child can never
	// look like: templateKeys: [T, U] templateValues: [?, T] There's no way for
	// the templateValue of U to depend on T. It's not in scope. Of course,
	// it's fine for T to depend on U, if e.g. Child<U> extends Parent<U> {}
	//
	// Note: this map is only read during 1) construction and 2) this.copy*
	// methods.
	subMapStarts *jutil.BitSet
	registry     *JSTypeRegistry
}

var emptyJSTypeSlice = []JSType{}
var emptyTemplateTypeSlice = []*TemplateType{}

var sharedSubMapStartsEmpty = jutil.NewBitSetWithSize(64)

var sharedSubMapStartsZero = func() *jutil.BitSet {
	b := jutil.NewBitSetWithSize(64)
	b.Set(0)
	return b
}()

// newTemplateTypeMapBitSet is new BitSet().
func newTemplateTypeMapBitSet() *jutil.BitSet {
	// Java: new BitSet() reserves 64 bits.
	return jutil.NewBitSetWithSize(64)
}

// Java: TemplateTypeMap.createEmpty(JSTypeRegistry)
func TemplateTypeMap_CreateEmpty(registry *JSTypeRegistry) *TemplateTypeMap {
	// This method should only be called during registry initialization.
	jutil.CheckArgument(registry.GetEmptyTemplateTypeMap() == nil)

	return NewTemplateTypeMap(registry, emptyTemplateTypeSlice, emptyJSTypeSlice, sharedSubMapStartsEmpty)
}

// Java: TemplateTypeMap(JSTypeRegistry, ImmutableList<TemplateType>, ImmutableList<JSType>, BitSet)
func NewTemplateTypeMap(
	registry *JSTypeRegistry,
	templateKeys []*TemplateType,
	templateValues []JSType,
	subMapStarts *jutil.BitSet) *TemplateTypeMap {
	return initTemplateTypeMap(registry.allocTemplateTypeMap(), registry, templateKeys, templateValues, subMapStarts)
}

func initTemplateTypeMap(
	m *TemplateTypeMap,
	registry *JSTypeRegistry,
	templateKeys []*TemplateType,
	templateValues []JSType,
	subMapStarts *jutil.BitSet) *TemplateTypeMap {
	jutil.CheckNotNull(templateKeys)
	jutil.CheckNotNull(templateValues)
	jutil.CheckArgument(len(templateValues) <= len(templateKeys))
	jutil.CheckArgument(len(templateKeys) <= subMapStarts.Size())

	m.registry = registry
	m.templateKeys = templateKeys
	m.templateValues = templateValues
	m.subMapStarts = subMapStarts

	nValues := len(m.templateValues)
	nKeys := len(m.templateKeys)
	if nKeys == 0 {
		m.resolvedTemplateValues = emptyJSTypeSlice
		return m
	}

	// Iteratively resolve any JSType values that refer to the TemplateType
	// keys of this TemplateTypeMap.
	var replacer TemplateTypeReplacer
	replacerInit := false

	var resolvedValues []JSType
	for i := 0; i < nKeys; i++ {
		var resolved JSType
		if i < nValues {
			templateValue := m.templateValues[i]
			if r, ok := resolveSimpleTemplateValue(templateValue, registry); ok {
				resolved = r
			} else {
				if !replacerInit {
					replacer = TemplateTypeReplacer{
						registry:                   registry,
						bindings:                   m,
						useUnknownForMissingValues: true,
						ownSubMapBoundary:          -1,
					}
					replacerInit = true
				}
				nextSubMap := m.NextSubMapStart(i + 1)
				templateKey := m.templateKeys[i]
				replacer.SetKeyType(templateKey, nextSubMap)
				resolved = replacer.visitJSType(templateValue)
			}
			if resolvedValues == nil && resolved == templateValue {
				continue
			}
		} else {
			resolved = m.templateKeys[i].GetBound()
		}
		if resolvedValues == nil {
			resolvedValues = registry.allocJSTypeSlice(nKeys)
			copy(resolvedValues[:i], m.templateValues[:i])
		}
		resolvedValues[i] = resolved
	}
	if resolvedValues == nil {
		m.resolvedTemplateValues = m.templateValues
	} else {
		m.resolvedTemplateValues = resolvedValues
	}
	return m
}

func resolveSimpleTemplateValue(t JSType, registry *JSTypeRegistry) (JSType, bool) {
	switch x := t.(type) {
	case *UnknownType:
		if x.isChecked {
			return registry.GetNativeType(JSTypeNative_UNKNOWN_TYPE), true
		}
		return t, true
	case *NumberType, *StringType, *BooleanType, *NullType, *VoidType, *AllType, *BigIntType, *NoType, NoObjectType, *EnumElementType, *InstanceObjectType, PrototypeObjectType:
		return t, true
	default:
		return nil, false
	}
}

// NextSubMapStart returns the index of the next submap in subMapStarts,
// starting at index (inclusive); will return subMapStarts.length if there is
// no next submap & we are in the last one.
//
// Java: TemplateTypeMap.nextSubMapStart(int)
func (m *TemplateTypeMap) NextSubMapStart(index int) int {
	for index < len(m.templateKeys) && !m.subMapStarts.Get(index) {
		index++
	}
	return index
}

// CopyFilledWithValues creates a new map in which any unfilled values in this
// map have been filled with values.
//
// If there are fewer values than unfilled values, `?` will be used to fill
// the rest.
//
// Java: TemplateTypeMap.copyFilledWithValues(ImmutableList<JSType>)
func (m *TemplateTypeMap) CopyFilledWithValues(values []JSType) *TemplateTypeMap {
	requiredUnknownCount := m.NumUnfilledTemplateKeys() - len(values)
	jutil.CheckArgumentMsg1(requiredUnknownCount >= 0, "%s", requiredUnknownCount)

	if m.NumUnfilledTemplateKeys() == 0 {
		return m // Nothing will change.
	}

	nKeys := len(m.templateKeys)
	out := m.registry.allocTemplateTypeMap()
	extendedValues := m.registry.allocJSTypeSlice(nKeys)
	idx := 0
	for _, e := range m.templateValues {
		if e == nil {
			panic(jutil.NewNullPointerExceptionMsg("at index " + jutil.ToJavaString(idx)))
		}
		extendedValues[idx] = e
		idx++
	}
	for _, e := range values {
		if e == nil {
			panic(jutil.NewNullPointerExceptionMsg("at index " + jutil.ToJavaString(idx)))
		}
		extendedValues[idx] = e
		idx++
	}
	for idx < nKeys {
		e := m.DefaultValueType(m.templateKeys[idx])
		if e == nil {
			panic(jutil.NewNullPointerExceptionMsg("at index " + jutil.ToJavaString(idx)))
		}
		extendedValues[idx] = e
		idx++
	}

	return initTemplateTypeMap(out, m.registry, m.templateKeys, extendedValues, m.subMapStarts)
}

// templateTypeMapCopyOfTypes is ImmutableList.copyOf(list), which rejects
// null elements.
func templateTypeMapCopyOfTypes(list []JSType) []JSType {
	out := make([]JSType, len(list))
	for i, e := range list {
		if e == nil {
			panic(jutil.NewNullPointerExceptionMsg("at index " + jutil.ToJavaString(i)))
		}
		out[i] = e
	}
	return out
}

// CopyWithExtension creates a new map in which the keys and values have been
// extended by extension.
//
// Before extension, any unfilled values in the initial map will be filled
// with `?`.
//
// Java: TemplateTypeMap.copyWithExtension(TemplateTypeMap)
func (m *TemplateTypeMap) CopyWithExtension(extension *TemplateTypeMap) *TemplateTypeMap {
	return m.CopyWithExtensionOfKeysWithSubMapStarts(
		extension.templateKeys, extension.templateValues, extension.subMapStarts)
}

// CopyWithExtensionOfKeys creates a new map in which the keys and values
// have been extended by keys and values respectively.
//
// Before extension, any unfilled values in the initial map will be filled
// with `?`.
//
// Java: TemplateTypeMap.copyWithExtension(ImmutableList<TemplateType>, ImmutableList<JSType>)
func (m *TemplateTypeMap) CopyWithExtensionOfKeys(keys []*TemplateType, values []JSType) *TemplateTypeMap {
	if len(keys) == 0 {
		if len(values) == 0 && m.NumUnfilledTemplateKeys() == 0 {
			return m
		}
		return m.CopyWithExtensionOfKeysWithSubMapStarts(keys, values, sharedSubMapStartsEmpty)
	}
	if len(keys) <= 64 {
		return m.CopyWithExtensionOfKeysWithSubMapStarts(keys, values, sharedSubMapStartsZero)
	}
	var subMapStarts jutil.BitSet
	subMapStarts.InitWithSize(len(keys))
	subMapStarts.Set(0)
	return m.CopyWithExtensionOfKeysWithSubMapStarts(keys, values, &subMapStarts)
}

// CopyWithExtensionOfKeysWithSubMapStarts creates a new map in which the
// keys and values have been extended by keys and values respectively.
//
// Before extension, any unfilled values in the initial map will be filled
// with `?`.
//
// Java: TemplateTypeMap.copyWithExtension(ImmutableList<TemplateType>, ImmutableList<JSType>, BitSet)
func (m *TemplateTypeMap) CopyWithExtensionOfKeysWithSubMapStarts(
	keys []*TemplateType, values []JSType, newsubMapStarts *jutil.BitSet) *TemplateTypeMap {
	extendedUnfilledCount := len(keys) - len(values)
	jutil.CheckArgumentMsg1(extendedUnfilledCount >= 0, "%s", extendedUnfilledCount)

	if m.NumUnfilledTemplateKeys() == 0 && len(keys) == 0 {
		return m // Nothing will change.
	}

	nOldKeys := len(m.templateKeys)
	nNewKeys := len(keys)
	extendedKeys := make([]*TemplateType, nOldKeys+nNewKeys)
	copy(extendedKeys, m.templateKeys)
	for i, k := range keys {
		if k == nil {
			panic(jutil.NewNullPointerExceptionMsg("at index " + jutil.ToJavaString(i)))
		}
		extendedKeys[nOldKeys+i] = k
	}

	out := m.registry.allocTemplateTypeMap()
	nTotalValues := nOldKeys + len(values)
	extendedValues := m.registry.allocJSTypeSlice(nTotalValues)
	idx := 0
	for _, e := range m.templateValues {
		if e == nil {
			panic(jutil.NewNullPointerExceptionMsg("at index " + jutil.ToJavaString(idx)))
		}
		extendedValues[idx] = e
		idx++
	}
	for idx < nOldKeys {
		e := m.DefaultValueType(m.templateKeys[idx])
		if e == nil {
			panic(jutil.NewNullPointerExceptionMsg("at index " + jutil.ToJavaString(idx)))
		}
		extendedValues[idx] = e
		idx++
	}
	for _, e := range values {
		if e == nil {
			panic(jutil.NewNullPointerExceptionMsg("at index " + jutil.ToJavaString(idx)))
		}
		extendedValues[idx] = e
		idx++
	}

	var extendedsubMapStarts *jutil.BitSet
	if nOldKeys == 0 && newsubMapStarts == sharedSubMapStartsZero {
		extendedsubMapStarts = sharedSubMapStartsZero
	} else if nNewKeys == 0 {
		extendedsubMapStarts = m.subMapStarts
	} else {
		extendedsubMapStarts = jutil.NewBitSetWithSize(nNewKeys + nOldKeys)
		extendedsubMapStarts.Or(m.subMapStarts)
		for i := 0; i < nNewKeys; i++ {
			extendedsubMapStarts.SetTo(i+nOldKeys, newsubMapStarts.Get(i))
		}
	}

	return initTemplateTypeMap(
		out, m.registry, extendedKeys, extendedValues, extendedsubMapStarts)
}

// CopyWithoutKeys creates a new map in which keys contained in removals are
// eliminated.
//
// The keys in removals will only be removed if they are unfilled.
//
// Java: TemplateTypeMap.copyWithoutKeys(Set<TemplateType>)
func (m *TemplateTypeMap) CopyWithoutKeys(removals *jutil.LinkedHashSet[*TemplateType]) *TemplateTypeMap {
	hasRemovals := false
	if removals.Size() > 0 {
		for i := len(m.templateValues); i < len(m.templateKeys); i++ {
			if removals.Contains(m.templateKeys[i]) {
				hasRemovals = true
				break
			}
		}
	}
	if !hasRemovals {
		return m // Nothing will change.
	}

	keys := make([]*TemplateType, len(m.templateValues), len(m.templateKeys)-1)
	copy(keys, m.templateKeys[:len(m.templateValues)])
	newsubMapStarts := newTemplateTypeMapBitSet()
	for i := 0; i < len(m.templateValues); i++ {
		newsubMapStarts.SetTo(i, m.subMapStarts.Get(i))
	}
	newKeysSize := 0
	inNewSubmap := false
	for i := len(m.templateValues); i < len(m.templateKeys); i++ {
		key := m.templateKeys[i]
		inNewSubmap = inNewSubmap || m.subMapStarts.Get(i)
		if !removals.Contains(key) {
			keys = append(keys, key)
			newsubMapStarts.SetTo(len(m.templateValues)+newKeysSize, inNewSubmap)
			newKeysSize++
			inNewSubmap = false
		}
	}

	return NewTemplateTypeMap(m.registry, keys, m.templateValues, newsubMapStarts)
}

// Java: TemplateTypeMap.size()
func (m *TemplateTypeMap) Size() int {
	return len(m.templateKeys)
}

// IsEmpty returns true if the map is empty; false otherwise.
//
// Java: TemplateTypeMap.isEmpty()
func (m *TemplateTypeMap) IsEmpty() bool {
	return len(m.templateKeys) == 0
}

// GetTemplateKeys returns a list of all template keys.
//
// Java: TemplateTypeMap.getTemplateKeys()
func (m *TemplateTypeMap) GetTemplateKeys() []*TemplateType {
	return m.templateKeys
}

// Java: TemplateTypeMap.getTemplateValues()
func (m *TemplateTypeMap) GetTemplateValues() []JSType {
	return m.templateValues
}

// HasTemplateKey returns true if this map contains the specified template
// key, false otherwise.
//
// Java: TemplateTypeMap.hasTemplateKey(TemplateType)
func (m *TemplateTypeMap) HasTemplateKey(templateKey *TemplateType) bool {
	return m.HasTemplateKeyFrom(templateKey, 0)
}

// HasTemplateKeyFrom returns true if this map contains the specified
// template key, false otherwise.
//
// If max is non-negative, this search will not consider any keys in the map
// from [0, max]. The intended use is to exclude keys that are conceptually
// defined on a type earlier in the supertype chain.
//
// For use in TemplateTypeReplacer.
//
// The "nextSubMap" parameter tells this method to only consider template
// keys in the range [nextSubMap, templateKeys.size()). This is used in the
// TemplateTypeMap constructor to avoid "leaking" sibling generic types from
// the same sub map. It is /not/ used in TypeInference when further
// specializing templates. (See the Java source for the long worked example
// of why this matters.)
//
// Java: TemplateTypeMap.hasTemplateKey(TemplateType, int)
func (m *TemplateTypeMap) HasTemplateKeyFrom(templateKey *TemplateType, nextSubMap int) bool {
	// NOTE: avoid iterators, for-each for performance and GC reasons
	keyCount := len(m.templateKeys)
	start := max(0, nextSubMap)
	for i := start; i < keyCount; i++ {
		entry := m.templateKeys[i]
		// Note: match by identity, not equality
		if templateKey == entry {
			return true
		}
	}
	return false
}

// GetTemplateKeyCountThisShouldAlwaysBeOneOrZeroButIsnt: TODO(b/139230800):
// This method should be deleted. It checks what should be an impossible
// case.
//
// Java: TemplateTypeMap.getTemplateKeyCountThisShouldAlwaysBeOneOrZeroButIsnt(TemplateType)
func (m *TemplateTypeMap) GetTemplateKeyCountThisShouldAlwaysBeOneOrZeroButIsnt(templateKey *TemplateType) int {
	matches := 0
	keyCount := len(m.templateKeys)

	// NOTE: avoid iterators, for-each for performance and GC reasons
	for i := 0; i < keyCount; i++ {
		entry := m.templateKeys[i]
		if templateKey == entry {
			matches++
		}
	}
	return matches
}

// Java: TemplateTypeMap.numUnfilledTemplateKeys()
func (m *TemplateTypeMap) NumUnfilledTemplateKeys() int {
	return len(m.templateKeys) - len(m.templateValues)
}

// HasTemplateType returns true if there is a JSType value associated with
// the specified template key; false otherwise.
//
// Java: TemplateTypeMap.hasTemplateType(TemplateType)
func (m *TemplateTypeMap) HasTemplateType(key *TemplateType) bool {
	return m.GetTemplateTypeIndex(key) != -1
}

// Java: TemplateTypeMap.getUnresolvedOriginalTemplateType(TemplateType)
func (m *TemplateTypeMap) GetUnresolvedOriginalTemplateType(key *TemplateType) JSType {
	index := m.GetTemplateTypeIndex(key)
	if index == -1 {
		return m.registry.GetNativeType(JSTypeNative_UNKNOWN_TYPE)
	}
	return m.templateValues[index]
}

// GetLastTemplateTypeKeyByName returns the final template key matching this
// name in the ordered list of template keys
//
// Caution: there may be multiple template keys with the same name. Before
// using this method, consider whether you really want reference name string
// equality over TemplateType identity- based equality.
//
// Java: TemplateTypeMap.getLastTemplateTypeKeyByName(String)
func (m *TemplateTypeMap) GetLastTemplateTypeKeyByName(keyName string) *TemplateType {
	size := len(m.templateKeys)
	for i := size - 1; i >= 0; i-- {
		key := m.templateKeys[i]
		if *key.GetReferenceName() == keyName {
			return key
		}
	}
	return nil
}

// GetTemplateTypeIndex returns the index of the JSType value associated with
// the specified template key. If no JSType value is associated, returns -1.
//
// Java: TemplateTypeMap.getTemplateTypeIndex(TemplateType)
func (m *TemplateTypeMap) GetTemplateTypeIndex(key *TemplateType) int {
	maxIndex := min(len(m.templateKeys), len(m.templateValues))
	for i := maxIndex - 1; i >= 0; i-- {
		if m.templateKeys[i] == key {
			return i
		}
	}
	return -1
}

// GetResolvedTemplateType returns the JSType value associated with the
// specified template key. If no JSType value is associated, returns the
// upper bound for generic, UNKNOWN_TYPE if unspecified.
//
// Java: TemplateTypeMap.getResolvedTemplateType(TemplateType)
func (m *TemplateTypeMap) GetResolvedTemplateType(key *TemplateType) JSType {
	index := m.GetTemplateTypeIndex(key)
	if index == -1 {
		return m.DefaultValueType(key)
	}
	return m.resolvedTemplateValues[index]
}

// Java: TemplateTypeMap.hasAnyTemplateTypesInternal()
func (m *TemplateTypeMap) HasAnyTemplateTypesInternal() bool {
	if m.resolvedTemplateValues != nil {
		for _, templateValue := range m.resolvedTemplateValues {
			if templateValue.HasAnyTemplateTypes() {
				return true
			}
		}
	}
	return false
}

// Java: TemplateTypeMap.toString()
func (m *TemplateTypeMap) String() string {
	s := ""

	n := len(m.templateKeys)
	s += "{ "
	for i := 0; i < n; i++ {
		s += "("
		s += jutil.ToJavaString(m.templateKeys[i])
		s += ","
		if i < len(m.templateValues) {
			s += jutil.ToJavaString(m.templateValues[i])
		}
		s += ","
		if m.resolvedTemplateValues != nil && i < len(m.resolvedTemplateValues) {
			s += jutil.ToJavaString(m.resolvedTemplateValues[i])
		}
		s += ") "
	}
	s += "}"

	return s
}

// Java: TemplateTypeMap.padToSameLength(ImmutableList<TemplateType>, ArrayList<JSType>)
func (m *TemplateTypeMap) PadToSameLength(keys []*TemplateType, builder *jutil.ArrayList[JSType]) {
	jutil.CheckArgument(builder.Size() <= len(keys))

	for i := builder.Size(); i < len(keys); i++ {
		builder.Add(m.DefaultValueType(keys[i]))
	}
}

// DefaultValueType returns the default value type for the given key type.
// Since the bounded generics feature was removed, in practice this always
// returns `?`.
//
// Java: TemplateTypeMap.defaultValueType(TemplateType)
func (m *TemplateTypeMap) DefaultValueType(type_ *TemplateType) JSType {
	if type_.GetBound().IsUnknownType() {
		return m.registry.GetNativeType(JSTypeNative_UNKNOWN_TYPE)
	}
	return type_
}

// Java: TemplateTypeMap.getIndicesOfSubmapsForTesting()
func (m *TemplateTypeMap) GetIndicesOfSubmapsForTesting() []int {
	indices := []int{}
	for i := 0; i < m.subMapStarts.Length(); i++ {
		if m.subMapStarts.Get(i) {
			indices = append(indices, i)
		}
	}
	return indices
}
