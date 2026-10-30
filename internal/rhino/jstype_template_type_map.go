// Port of com/google/javascript/rhino/jstype/TemplateTypeMap.java @ closure-compiler v20261021

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
	"math/bits"
	"slices"

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
	registry               *JSTypeRegistry
	// Overflow bitset words used only when len(templateKeys) > 64.
	// Nil in the common <= 64 key case.
	subMapStarts []uint64
	// Inline 64-bit submap-start bitset for maps with <= 64 template keys.
	subMapStartsBits uint64
}

var emptyJSTypeSlice = []JSType{}
var emptyTemplateTypeSlice = []*TemplateType{}

func subMapBitWords(nBits int) []uint64 {
	if nBits <= 0 {
		return nil
	}
	return make([]uint64, (nBits+63)>>6)
}

func subMapBitGet(words []uint64, i int) bool {
	w := i >> 6
	return i >= 0 && w < len(words) && (words[w]&(uint64(1)<<uint(i&63))) != 0
}

func subMapBitSetTo(words []uint64, i int, v bool) {
	if v {
		words[i>>6] |= uint64(1) << uint(i&63)
	} else {
		words[i>>6] &^= uint64(1) << uint(i&63)
	}
}

func subMapBitLength(words []uint64) int {
	for w := len(words) - 1; w >= 0; w-- {
		if words[w] != 0 {
			return (w+1)*64 - bits.LeadingZeros64(words[w])
		}
	}
	return 0
}

// Java: TemplateTypeMap.createEmpty(JSTypeRegistry)
func TemplateTypeMap_CreateEmpty(registry *JSTypeRegistry) *TemplateTypeMap {
	// This method should only be called during registry initialization.
	jutil.CheckArgument(registry.GetEmptyTemplateTypeMap() == nil)

	return initTemplateTypeMap(&TemplateTypeMap{}, registry, emptyTemplateTypeSlice, emptyJSTypeSlice, 0, nil)
}

func bitSetToInline(bs []uint64, nKeys int) (uint64, []uint64) {
	jutil.CheckArgument(nKeys <= len(bs)*64)
	if subMapBitLength(bs) <= 64 && nKeys <= 64 {
		var mask uint64
		if len(bs) > 0 {
			mask = bs[0]
		}
		return mask, nil
	}
	return 0, bs
}

// Java: TemplateTypeMap(JSTypeRegistry, ImmutableList<TemplateType>, ImmutableList<JSType>, BitSet)
func NewTemplateTypeMap(
	registry *JSTypeRegistry,
	templateKeys []*TemplateType,
	templateValues []JSType,
	subMapStarts []uint64) *TemplateTypeMap {
	jutil.CheckNotNull(templateKeys)
	jutil.CheckNotNull(templateValues)
	bits64, overflow := bitSetToInline(subMapStarts, len(templateKeys))
	return initTemplateTypeMap(&TemplateTypeMap{}, registry, templateKeys, templateValues, bits64, overflow)
}

func initTemplateTypeMap(
	m *TemplateTypeMap,
	registry *JSTypeRegistry,
	templateKeys []*TemplateType,
	templateValues []JSType,
	subMapStartsBits uint64,
	subMapStarts []uint64) *TemplateTypeMap {
	jutil.CheckNotNull(templateKeys)
	jutil.CheckNotNull(templateValues)
	jutil.CheckArgument(len(templateValues) <= len(templateKeys))
	if subMapStarts != nil {
		jutil.CheckArgument(len(templateKeys) <= len(subMapStarts)*64)
	} else {
		jutil.CheckArgument(len(templateKeys) <= 64)
	}

	m.registry = registry
	m.templateKeys = templateKeys
	m.templateValues = templateValues
	m.subMapStartsBits = subMapStartsBits
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
	if nValues < nKeys {
		resolvedValues = make([]JSType, nKeys)
		for i := nValues; i < nKeys; i++ {
			resolvedValues[i] = m.templateKeys[i].GetBound()
		}
	}

	var simpleMask uint64
	for i := nValues - 1; i >= 0; i-- {
		templateValue := m.templateValues[i]
		var resolved JSType
		if r, ok := resolveSimpleTemplateValue(templateValue, registry); ok {
			resolved = r
			if i < 64 {
				simpleMask |= uint64(1) << uint(i)
			}
		} else if tt, ok := templateValue.(*TemplateType); ok && tt != m.templateKeys[i] {
			nextSubMap := m.NextSubMapStart(i + 1)
			if !m.HasTemplateKeyFrom(tt, nextSubMap) {
				resolved = tt
			} else {
				idx := m.GetTemplateTypeIndex(tt)
				if idx == -1 {
					resolved = registry.GetNativeType(JSTypeNative_UNKNOWN_TYPE)
					if i < 64 {
						simpleMask |= uint64(1) << uint(i)
					}
				} else if idx >= nextSubMap && idx < 64 && (simpleMask&(uint64(1)<<uint(idx))) != 0 {
					if resolvedValues != nil {
						resolved = resolvedValues[idx]
					} else {
						resolved = m.templateValues[idx]
					}
					if i < 64 {
						simpleMask |= uint64(1) << uint(i)
					}
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
					replacer.SetKeyType(m.templateKeys[i], nextSubMap)
					resolved = replacer.visitJSType(templateValue)
				}
			}
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
			replacer.SetKeyType(m.templateKeys[i], nextSubMap)
			resolved = replacer.visitJSType(templateValue)
		}
		if resolvedValues == nil {
			if resolved == templateValue {
				continue
			}
			resolvedValues = make([]JSType, nKeys)
			copy(resolvedValues[i+1:], m.templateValues[i+1:])
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
	case *FunctionTypeBase:
		if x.isNoObjectOrNo() {
			return t, true
		}
		return nil, false
	case *NamedType, *NumberType, *StringType, *BooleanType, *NullType, *VoidType, *AllType, *BigIntType, *EnumType, *NoResolvedType, *EnumElementType, *InstanceObjectType, *PrototypeObjectTypeBase, *RecordType, *SymbolTypeBase, *KnownSymbolType:
		return t, true
	case *UnionType:
		alts := x.GetAlternates()
		for i := 0; i < len(alts); i++ {
			alt := alts[i]
			r, ok := resolveSimpleTemplateValue(alt, registry)
			if !ok || r != alt {
				return nil, false
			}
		}
		return t, true
	default:
		return nil, false
	}
}

func (m *TemplateTypeMap) subMapStartAt(i int) bool {
	if m.subMapStarts != nil {
		return subMapBitGet(m.subMapStarts, i)
	}
	if uint(i) >= 64 {
		return false
	}
	return (m.subMapStartsBits & (uint64(1) << uint(i))) != 0
}

// NextSubMapStart returns the index of the next submap in subMapStarts,
// starting at index (inclusive); will return subMapStarts.length if there is
// no next submap & we are in the last one.
//
// Java: TemplateTypeMap.nextSubMapStart(int)
func (m *TemplateTypeMap) NextSubMapStart(index int) int {
	n := len(m.templateKeys)
	if index >= n {
		return index
	}
	if m.subMapStarts == nil && index >= 0 {
		if index >= 64 {
			return n
		}
		masked := m.subMapStartsBits >> uint(index)
		if masked == 0 {
			return n
		}
		return min(index+bits.TrailingZeros64(masked), n)
	}
	for index < n && !subMapBitGet(m.subMapStarts, index) {
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

	var t0, t1 *JSTypeBase
	var cacheIdx int
	canCache := false
	if len(values) <= 2 && m.registry != nil {
		canCache = true
		if len(values) >= 1 {
			if values[0] == nil {
				canCache = false
			} else {
				t0 = values[0].GetJSTypeBase()
			}
		}
		if len(values) == 2 && canCache {
			if values[1] == nil {
				canCache = false
			} else {
				t1 = values[1].GetJSTypeBase()
			}
		}
		if canCache && !m.registry.parallelReadOnly {
			cacheIdx = m.registry.templateTypeMapCacheIndex(m, t0, t1)
			if m.registry.templateTypeMapCache != nil {
				entry := &m.registry.templateTypeMapCache[cacheIdx]
				if entry.base == m && entry.t0 == t0 && entry.t1 == t1 {
					return entry.res
				}
			}
		}
	}

	nKeys := len(m.templateKeys)
	out := &TemplateTypeMap{}
	extendedValues := make([]JSType, nKeys)
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

	res := initTemplateTypeMap(out, m.registry, m.templateKeys, extendedValues, m.subMapStartsBits, m.subMapStarts)
	if canCache && !m.registry.parallelReadOnly {
		if m.registry.templateTypeMapCache == nil {
			m.registry.templateTypeMapCache = new([jsTypeRegistryTemplatizedCacheSize]templateTypeMapCacheEntry)
		}
		m.registry.templateTypeMapCache[cacheIdx] = templateTypeMapCacheEntry{
			base: m,
			t0:   t0,
			t1:   t1,
			res:  res,
		}
	}
	return res
}

// CopyWithExtension creates a new map in which the keys and values have been
// extended by extension.
//
// Before extension, any unfilled values in the initial map will be filled
// with `?`.
//
// Java: TemplateTypeMap.copyWithExtension(TemplateTypeMap)
func (m *TemplateTypeMap) CopyWithExtension(extension *TemplateTypeMap) *TemplateTypeMap {
	return m.copyWithExtensionInternal(
		extension.templateKeys, extension.templateValues, extension.subMapStartsBits, extension.subMapStarts)
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
		return m.copyWithExtensionInternal(keys, values, 0, nil)
	}
	if len(keys) <= 64 {
		return m.copyWithExtensionInternal(keys, values, 1, nil)
	}
	subMapStarts := subMapBitWords(len(keys))
	subMapStarts[0] = 1
	return m.copyWithExtensionInternal(keys, values, 0, subMapStarts)
}

// CopyWithExtensionOfKeysWithSubMapStarts creates a new map in which the
// keys and values have been extended by keys and values respectively.
//
// Before extension, any unfilled values in the initial map will be filled
// with `?`.
//
// Java: TemplateTypeMap.copyWithExtension(ImmutableList<TemplateType>, ImmutableList<JSType>, BitSet)
func (m *TemplateTypeMap) CopyWithExtensionOfKeysWithSubMapStarts(
	keys []*TemplateType, values []JSType, newsubMapStarts []uint64) *TemplateTypeMap {
	if subMapBitLength(newsubMapStarts) <= 64 && len(keys) <= 64 {
		var mask uint64
		if len(newsubMapStarts) > 0 {
			mask = newsubMapStarts[0]
		}
		return m.copyWithExtensionInternal(keys, values, mask, nil)
	}
	return m.copyWithExtensionInternal(keys, values, 0, newsubMapStarts)
}

func (m *TemplateTypeMap) copyWithExtensionInternal(
	keys []*TemplateType, values []JSType, newBits uint64, newOverflow []uint64) *TemplateTypeMap {
	extendedUnfilledCount := len(keys) - len(values)
	jutil.CheckArgumentMsg1(extendedUnfilledCount >= 0, "%s", extendedUnfilledCount)

	if m.NumUnfilledTemplateKeys() == 0 && len(keys) == 0 {
		return m // Nothing will change.
	}

	nOldKeys := len(m.templateKeys)
	nNewKeys := len(keys)
	var extendedKeys []*TemplateType
	if nNewKeys == 0 {
		extendedKeys = m.templateKeys
	} else {
		extendedKeys = make([]*TemplateType, nOldKeys+nNewKeys)
		copy(extendedKeys, m.templateKeys)
		for i, k := range keys {
			if k == nil {
				panic(jutil.NewNullPointerExceptionMsg("at index " + jutil.ToJavaString(i)))
			}
			extendedKeys[nOldKeys+i] = k
		}
	}

	out := &TemplateTypeMap{}
	nTotalValues := nOldKeys + len(values)
	extendedValues := make([]JSType, nTotalValues)
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

	if nNewKeys == 0 {
		return initTemplateTypeMap(out, m.registry, extendedKeys, extendedValues, m.subMapStartsBits, m.subMapStarts)
	}
	if nOldKeys+nNewKeys <= 64 && m.subMapStarts == nil && newOverflow == nil {
		var mask uint64
		if nOldKeys > 0 {
			mask = m.subMapStartsBits & ((uint64(1) << uint(nOldKeys)) - 1)
		}
		var extMask uint64
		if nNewKeys == 64 {
			extMask = newBits
		} else {
			extMask = (newBits & ((uint64(1) << uint(nNewKeys)) - 1)) << uint(nOldKeys)
		}
		return initTemplateTypeMap(out, m.registry, extendedKeys, extendedValues, mask|extMask, nil)
	}

	extendedsubMapStarts := subMapBitWords(nNewKeys + nOldKeys)
	for i := 0; i < nOldKeys; i++ {
		if m.subMapStartAt(i) {
			subMapBitSetTo(extendedsubMapStarts, i, true)
		}
	}
	for i := 0; i < nNewKeys; i++ {
		var bit bool
		if newOverflow != nil {
			bit = subMapBitGet(newOverflow, i)
		} else if i < 64 {
			bit = (newBits & (uint64(1) << uint(i))) != 0
		}
		subMapBitSetTo(extendedsubMapStarts, i+nOldKeys, bit)
	}

	return initTemplateTypeMap(out, m.registry, extendedKeys, extendedValues, 0, extendedsubMapStarts)
}

// CopyWithoutKeySlice creates a new map in which keys contained in removals
// are eliminated.
//
// The keys in removals will only be removed if they are unfilled.
func (m *TemplateTypeMap) CopyWithoutKeySlice(removals []*TemplateType) *TemplateTypeMap {
	return m.copyWithoutKeySlice(removals)
}

func (m *TemplateTypeMap) copyWithoutKeySlice(removals []*TemplateType) *TemplateTypeMap {
	if len(removals) == 0 {
		return m
	}
	return m.copyWithoutKeysWhere(func(t *TemplateType) bool {
		return slices.Contains(removals, t)
	})
}

func (m *TemplateTypeMap) copyWithoutKeysWhere(contains func(*TemplateType) bool) *TemplateTypeMap {
	hasRemovals := false
	for i := len(m.templateValues); i < len(m.templateKeys); i++ {
		if contains(m.templateKeys[i]) {
			hasRemovals = true
			break
		}
	}
	if !hasRemovals {
		return m // Nothing will change.
	}

	nFilled := len(m.templateValues)
	keys := make([]*TemplateType, nFilled, len(m.templateKeys)-1)
	copy(keys, m.templateKeys[:nFilled])

	if len(m.templateKeys) <= 64 && m.subMapStarts == nil {
		var newBits uint64
		if nFilled > 0 {
			newBits = m.subMapStartsBits & ((uint64(1) << uint(nFilled)) - 1)
		}
		newKeysSize := 0
		inNewSubmap := false
		for i := nFilled; i < len(m.templateKeys); i++ {
			key := m.templateKeys[i]
			inNewSubmap = inNewSubmap || ((m.subMapStartsBits & (uint64(1) << uint(i))) != 0)
			if !contains(key) {
				keys = append(keys, key)
				if inNewSubmap {
					newBits |= uint64(1) << uint(nFilled+newKeysSize)
				}
				newKeysSize++
				inNewSubmap = false
			}
		}
		return initTemplateTypeMap(&TemplateTypeMap{}, m.registry, keys, m.templateValues, newBits, nil)
	}

	newsubMapStarts := subMapBitWords(max(64, len(m.templateKeys)))
	for i := 0; i < nFilled; i++ {
		subMapBitSetTo(newsubMapStarts, i, m.subMapStartAt(i))
	}
	newKeysSize := 0
	inNewSubmap := false
	for i := nFilled; i < len(m.templateKeys); i++ {
		key := m.templateKeys[i]
		inNewSubmap = inNewSubmap || m.subMapStartAt(i)
		if !contains(key) {
			keys = append(keys, key)
			subMapBitSetTo(newsubMapStarts, nFilled+newKeysSize, inNewSubmap)
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
		if key.name == keyName {
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
func (m *TemplateTypeMap) HasAnyTemplateTypesInternal(st *JSTypeStack) bool {
	if m.resolvedTemplateValues != nil {
		for _, templateValue := range m.resolvedTemplateValues {
			if jsTypeHasAnyTemplateTypesIn(templateValue, st) {
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
	if m.subMapStarts == nil {
		for b := m.subMapStartsBits; b != 0; {
			tz := bits.TrailingZeros64(b)
			indices = append(indices, tz)
			b &= b - 1
		}
		return indices
	}
	for w, word := range m.subMapStarts {
		for word != 0 {
			tz := bits.TrailingZeros64(word)
			indices = append(indices, w*64+tz)
			word &= word - 1
		}
	}
	return indices
}
