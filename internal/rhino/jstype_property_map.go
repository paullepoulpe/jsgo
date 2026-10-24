// Port of com/google/javascript/rhino/jstype/PropertyMap.java @ closure-compiler v20261015

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
 *   Nick Santos
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino

import (
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

type propMapExt struct {
	byID                     map[rhinoStringID]*Property
	knownSymbols             *jutil.ObjectLinkedHashMap[*KnownSymbolType, *Property]
	cachedKeySet             []string
	cachedKnownSymbolsKeySet []*KnownSymbolType
}

// PropertyMap is the representation for a collection of properties on an
// object.
type PropertyMap struct {
	// Own string properties, the slice entries(). Small maps (<=24 entries)
	// scan entries directly by 32-bit rhinoStringID; larger maps maintain
	// ext.byID for O(1) lookup and sort entries by StringCompareTo lazily when
	// key-ordered iteration is requested.
	entries_     []*Property
	parentSource *ObjectTypeBase
	ext          *propMapExt

	// A "timestamp" for map mutations to validate cachedKeySet and
	// cachedKnownSymbolsKeySet.
	//
	// If this value is less than the counter in any ancestor map, the cache
	// is invalid. The update algorithm is similar to using a global counter,
	// but uses a distributed approach that prevents the need for a single
	// source of truth.
	cachedKeySetCounter int32

	// Go-only: language-adapter — EMPTY_MAP is built from ImmutableMaps,
	// whose put() throws UnsupportedOperationException.
	immutable   bool
	sorted      bool
	hasNonASCII bool
}

// entries is the slice of own string properties.
func (m *PropertyMap) entries() []*Property {
	return m.entries_
}

func (m *PropertyMap) setEntries(entries []*Property) {
	m.entries_ = entries
}

func (m *PropertyMap) getParentSource() ObjectType {
	if m.parentSource == nil {
		return nil
	}
	return m.parentSource.selfObj()
}

func (m *PropertyMap) ensureExt() *propMapExt {
	if m.ext == nil {
		m.ext = &propMapExt{}
	}
	return m.ext
}

func propMapIsASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= 0x80 {
			return false
		}
	}
	return true
}

// Java: PropertyMap.EMPTY_MAP
var PropertyMap_EMPTY_MAP = &PropertyMap{
	ext: &propMapExt{
		knownSymbols: jutil.NewObjectLinkedHashMap[*KnownSymbolType, *Property](),
	},
	immutable: true,
	sorted:    true,
}

// Java: PropertyMap()
func NewPropertyMap() *PropertyMap {
	return &PropertyMap{sorted: true}
}

func (m *PropertyMap) ensureSorted() {
	if m.sorted {
		return
	}
	if len(m.entries()) <= 1 {
		m.sorted = true
		return
	}
	slices.SortFunc(m.entries(), func(a, b *Property) int {
		return jutil.StringCompareTo(rhinoStringFromID(a.nameID), rhinoStringFromID(b.nameID))
	})
	m.sorted = true
}

// Java: PropertyMap.immutableEmptyMap()
func PropertyMap_ImmutableEmptyMap() *PropertyMap {
	return PropertyMap_EMPTY_MAP
}

// Java: PropertyMap.setParentSource(ObjectType)
func (m *PropertyMap) SetParentSource(ownerType ObjectType) {
	if m == PropertyMap_EMPTY_MAP {
		return
	}
	m.parentSource = objTypeBase(ownerType)
	m.IncrementCachedKeySetCounter()
}

// Java: PropertyMap.setParentForTesting(PropertyMap)
func (m *PropertyMap) SetParentForTesting(parent *PropertyMap) {
	m.parentSource = parent.parentSource
	m.IncrementCachedKeySetCounter()
}

// GetPrimaryParent returns the direct parent of this property map.
//
// Java: PropertyMap.getPrimaryParent()
func (m *PropertyMap) GetPrimaryParent() *PropertyMap {
	parentSource := m.getParentSource()
	if parentSource == nil {
		return nil
	}
	iProto := parentSource.GetImplicitPrototype()
	if iProto == nil {
		return nil
	}
	return iProto.GetPropertyMap()
}

// getSecondaryParentObjects returns the secondary parents of this property
// map, for interfaces that need multiple inheritance or for interfaces of
// abstract classes.
//
// Java: PropertyMap.GetSecondaryParentObjects()
func (m *PropertyMap) GetSecondaryParentObjects() []ObjectType {
	parentSource := m.getParentSource()
	if parentSource == nil {
		return nil
	}
	if ctor := parentSource.GetConstructor(); ctor != nil && ctor.IsAbstract() {
		return ctor.GetOwnImplementedInterfaces()
	}
	return parentSource.GetCtorExtendedInterfaces()
}

// Java: PropertyMap.findClosest(String)
func (m *PropertyMap) FindClosest(name string) *Property_OwnedProperty {
	return m.FindClosestByKey(NewProperty_StringKey(name))
}

// findClosestPairByKey resolves the closest (owner, Property) pair without
// heap-allocating a Property_OwnedProperty wrapper.
func (m *PropertyMap) findClosestPairByKey(name Property_Key) (ObjectType, *Property) {
	if name.Kind() == Property_KeyKind_STRING {
		str := name.String_()
		if len(m.entries()) == 1 && !m.hasNonASCII {
			if p := m.entries()[0]; rhinoStringFromID(p.nameID) == str {
				return m.getParentSource(), p
			}
		}
		return m.findClosestPairByStringID(str, rhinoStringPoolLookup(str))
	}
	// Check primary parents which always has precedence over secondary.
	for pm := m; pm != nil; pm = pm.GetPrimaryParent() {
		if pm.ext != nil && pm.ext.knownSymbols != nil {
			if prop := pm.ext.knownSymbols.Get(name.Symbol()); prop != nil {
				return pm.getParentSource(), prop
			}
		}
	}

	// Recurse into secondary parents.
	for pm := m; pm != nil; pm = pm.GetPrimaryParent() {
		for _, o := range pm.GetSecondaryParentObjects() {
			parent := o.GetPropertyMap()
			if parent != nil {
				owner, prop := parent.findClosestPairByKey(name)
				if prop != nil {
					return owner, prop
				}
			}
		}
	}

	return nil, nil
}

func (m *PropertyMap) findClosestPairByStringID(str string, id rhinoStringID) (ObjectType, *Property) {
	for pm := m; pm != nil; pm = pm.GetPrimaryParent() {
		if prop := pm.getOwnPropertyByID(str, id); prop != nil {
			return pm.getParentSource(), prop
		}
	}
	for pm := m; pm != nil; pm = pm.GetPrimaryParent() {
		for _, o := range pm.GetSecondaryParentObjects() {
			if parent := o.GetPropertyMap(); parent != nil {
				if owner, prop := parent.findClosestPairByStringID(str, id); prop != nil {
					return owner, prop
				}
			}
		}
	}
	return nil, nil
}

// Java: PropertyMap.findClosest(Property.Key)
func (m *PropertyMap) FindClosestByKey(name Property_Key) *Property_OwnedProperty {
	owner, prop := m.findClosestPairByKey(name)
	if prop == nil {
		return nil
	}
	return NewProperty_OwnedProperty(owner, prop)
}

func (m *PropertyMap) getOwnPropertyByID(propertyName string, id rhinoStringID) *Property {
	if len(m.entries()) == 0 {
		return nil
	}
	if !m.hasNonASCII {
		if id == 0 {
			return nil
		}
		if m.ext != nil && m.ext.byID != nil {
			return m.ext.byID[id]
		}
		for _, p := range m.entries() {
			if p.nameID == id {
				return p
			}
		}
		return nil
	}
	for _, p := range m.entries() {
		if (id != 0 && p.nameID == id) || jutil.StringCompareTo(rhinoStringFromID(p.nameID), propertyName) == 0 {
			return p
		}
	}
	return nil
}

// Java: PropertyMap.getOwnProperty(String)
func (m *PropertyMap) GetOwnProperty(propertyName string) *Property {
	if len(m.entries()) == 0 {
		return nil
	}
	if len(m.entries()) == 1 && !m.hasNonASCII {
		if p := m.entries()[0]; rhinoStringFromID(p.nameID) == propertyName {
			return p
		}
		return nil
	}
	return m.getOwnPropertyByID(propertyName, rhinoStringPoolLookup(propertyName))
}

// Java: PropertyMap.getOwnProperty(Property.Key)
func (m *PropertyMap) GetOwnPropertyByKey(propertyName Property_Key) *Property {
	switch propertyName.Kind() {
	case Property_KeyKind_STRING:
		return m.GetOwnProperty(propertyName.String_())
	}
	if m.ext != nil && m.ext.knownSymbols != nil {
		return m.ext.knownSymbols.Get(propertyName.Symbol())
	}
	return nil
}

// Java: PropertyMap.getPropertiesCount()
func (m *PropertyMap) GetPropertiesCount() int {
	primaryParent := m.GetPrimaryParent()
	if primaryParent == nil {
		return len(m.entries())
	}

	return len(m.GetAllKeys().StringKeys())
}

// GetOwnPropertyNames returns properties.keySet(), in key order.
//
// Go-only: language-adapter — a snapshot of the live keySet() view.
//
// Java: PropertyMap.getOwnPropertyNames()
func (m *PropertyMap) GetOwnPropertyNames() []string {
	m.ensureSorted()
	entries := m.entries()
	if len(entries) == 0 {
		return nil
	}
	out := make([]string, len(entries))
	for i, p := range entries {
		out[i] = rhinoStringFromID(p.nameID)
	}
	return out
}

// ForEachOwnPropertyName calls fn for each own property name in sorted order without allocating a slice.
func (m *PropertyMap) ForEachOwnPropertyName(fn func(string)) {
	if len(m.entries()) == 0 {
		return
	}
	m.ensureSorted()
	for _, p := range m.entries() {
		fn(rhinoStringFromID(p.nameID))
	}
}

// HasAnyOwnProperties reports whether m has at least one own string property.
func (m *PropertyMap) HasAnyOwnProperties() bool {
	return len(m.entries()) != 0
}

// Java: PropertyMap.getOwnKnownSymbols()
func (m *PropertyMap) GetOwnKnownSymbols() []*KnownSymbolType {
	if m.ext != nil && m.ext.knownSymbols != nil {
		return m.ext.knownSymbols.KeySlice()
	}
	return nil
}

// Java: PropertyMap.getAllKeys()
func (m *PropertyMap) GetAllKeys() *PropertyMap_AllKeys {
	ancestors := jutil.NewLinkedHashSet[*PropertyMap]()
	m.CollectAllAncestors(ancestors)

	var maxAncestorCounter int32 = 0
	for ancestor := range ancestors.All() {
		if ancestor.cachedKeySetCounter > maxAncestorCounter {
			maxAncestorCounter = ancestor.cachedKeySetCounter
		}
	}

	// If any counter is greater than this counter, there has been a mutation
	// and the cache must be rebuilt.
	ext := m.ensureExt()
	if maxAncestorCounter != m.cachedKeySetCounter || ext.cachedKeySet == nil {
		keys := jutil.NewTreeSet[string](jutil.StringCompareTo)
		knownSymbolsKeys := jutil.NewObjectLinkedHashSet[*KnownSymbolType]()
		for ancestor := range ancestors.All() {
			// Update the counters in all ancestors.
			//
			// This update scheme is convergent. As long as there are no
			// mutations, calls #keySet will eventually set stable caches on all
			// maps.
			ancestor.cachedKeySetCounter = maxAncestorCounter
			if ancestor.ext != nil {
				ancestor.ext.cachedKeySet = nil
				ancestor.ext.cachedKnownSymbolsKeySet = nil
			}

			ancestor.ensureSorted()
			for _, p := range ancestor.entries() {
				keys.Add(rhinoStringFromID(p.nameID))
			}
			if ancestor.ext != nil && ancestor.ext.knownSymbols != nil {
				for s := range ancestor.ext.knownSymbols.Keys() {
					knownSymbolsKeys.Add(s)
				}
			}
		}
		if keys.IsEmpty() {
			ext.cachedKeySet = []string{}
		} else {
			ext.cachedKeySet = keys.ToSlice()
		}
		ext.cachedKnownSymbolsKeySet = knownSymbolsKeys.ToSlice()
	}

	return NewPropertyMap_AllKeys(ext.cachedKeySet, ext.cachedKnownSymbolsKeySet)
}

// PropertyMap_AllKeys is the record of all the keys of a map and its
// ancestors.
type PropertyMap_AllKeys struct {
	stringKeys      []string
	knownSymbolKeys []*KnownSymbolType
}

// Java: PropertyMap.AllKeys(ImmutableSortedSet<String>, Set<KnownSymbolType>)
func NewPropertyMap_AllKeys(stringKeys []string, knownSymbolKeys []*KnownSymbolType) *PropertyMap_AllKeys {
	return &PropertyMap_AllKeys{stringKeys: stringKeys, knownSymbolKeys: knownSymbolKeys}
}

// StringKeys is the sorted set of string keys.
//
// Java: PropertyMap.AllKeys.stringKeys()
func (k *PropertyMap_AllKeys) StringKeys() []string { return k.stringKeys }

// KnownSymbolKeys is the set of known symbol keys.
//
// Java: PropertyMap.AllKeys.knownSymbolKeys()
func (k *PropertyMap_AllKeys) KnownSymbolKeys() []*KnownSymbolType { return k.knownSymbolKeys }

// Equals is the record's component-wise equality; both components are
// sets.
//
// Java: PropertyMap.AllKeys.equals(Object)
func (k *PropertyMap_AllKeys) Equals(o any) bool {
	other, ok := o.(*PropertyMap_AllKeys)
	if !ok || other == nil {
		return false
	}
	return propertyMap_setsEqual(k.stringKeys, other.stringKeys) &&
		propertyMap_setsEqual(k.knownSymbolKeys, other.knownSymbolKeys)
}

// HashCode is the record's hash of its components.
//
// Java: PropertyMap.AllKeys.hashCode()
func (k *PropertyMap_AllKeys) HashCode() int32 {
	return 31*jutil.SetHashCode(k.stringKeys) + jutil.SetHashCode(k.knownSymbolKeys)
}

// propertyMap_setsEqual is Set.equals for two duplicate-free slices.
func propertyMap_setsEqual[T any](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for _, x := range a {
		found := false
		for _, y := range b {
			if jutil.ObjectsEquals(x, y) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// Java: PropertyMap.CollectAllAncestors(LinkedHashSet<PropertyMap>)
func (m *PropertyMap) CollectAllAncestors(ancestors *jutil.LinkedHashSet[*PropertyMap]) {
	if !ancestors.Add(m) {
		return
	}

	primaryParent := m.GetPrimaryParent()
	if primaryParent != nil {
		primaryParent.CollectAllAncestors(ancestors)
	}

	for _, parentType := range m.GetSecondaryParentObjects() {
		parentMap := parentType.GetPropertyMap()
		if parentMap != nil {
			parentMap.CollectAllAncestors(ancestors)
		}
	}
}

// Java: PropertyMap.putProperty(String, Property)
func (m *PropertyMap) PutProperty(name string, newProp *Property) {
	nameID := newProp.nameID
	if newProp.isSymbol || rhinoStringFromID(nameID) != name {
		nameID = rhinoStringPoolIntern(name)
	}
	isASCII := !m.hasNonASCII && propMapIsASCII(name)
	var oldProp *Property
	oldIdx := -1
	if m.ext != nil && m.ext.byID != nil && isASCII {
		oldProp = m.ext.byID[nameID]
	} else if !m.hasNonASCII && isASCII {
		for i, p := range m.entries() {
			if p.nameID == nameID {
				oldProp = p
				oldIdx = i
				break
			}
		}
	} else {
		for i, p := range m.entries() {
			if p.nameID == nameID || jutil.StringCompareTo(rhinoStringFromID(p.nameID), name) == 0 {
				oldProp = p
				oldIdx = i
				break
			}
		}
	}

	if oldProp == nil {
		// The cache is only invalidated if this is a new property name.
		m.IncrementCachedKeySetCounter()
	} else {
		// This is to keep previously inferred JsDoc info, e.g., in a
		// replaceScript scenario.
		newProp.SetJSDocInfo(oldProp.GetJSDocInfo())
	}

	if m.immutable {
		panic(jutil.NewUnsupportedOperationException())
	}
	if oldProp != nil {
		if oldIdx >= 0 {
			m.entries()[oldIdx] = newProp
		} else {
			m.ext.byID[nameID] = newProp
			for i, p := range m.entries() {
				if p.nameID == nameID {
					m.entries()[i] = newProp
					break
				}
			}
		}
		return
	}
	if !isASCII && !m.hasNonASCII {
		m.hasNonASCII = true
		if m.ext != nil {
			m.ext.byID = nil
		}
	}
	if len(m.entries()) == 0 {
		m.sorted = true
	} else if m.sorted && jutil.StringCompareTo(rhinoStringFromID(m.entries()[len(m.entries())-1].nameID), name) >= 0 {
		m.sorted = false
	}
	m.setEntries(append(m.entries(), newProp))
	if !m.hasNonASCII {
		if m.ext != nil && m.ext.byID != nil {
			m.ext.byID[nameID] = newProp
		} else if len(m.entries()) > 24 {
			ext := m.ensureExt()
			ext.byID = make(map[rhinoStringID]*Property, 32)
			for _, p := range m.entries() {
				ext.byID[p.nameID] = p
			}
		}
	}
}

// Java: PropertyMap.putProperty(KnownSymbolType, Property)
func (m *PropertyMap) PutPropertyBySymbol(symbol *KnownSymbolType, newProp *Property) {
	var oldProp *Property
	if m.ext != nil && m.ext.knownSymbols != nil {
		oldProp = m.ext.knownSymbols.Get(symbol)
	}

	if oldProp == nil {
		// The cache is only invalidated if this is a new property name.
		m.IncrementCachedKeySetCounter()
	}

	ext := m.ensureExt()
	if ext.knownSymbols == nil {
		ext.knownSymbols = jutil.NewObjectLinkedHashMap[*KnownSymbolType, *Property]()
	}
	if m.immutable {
		panic(jutil.NewUnsupportedOperationException())
	}
	ext.knownSymbols.Put(symbol, newProp)
}

// Java: PropertyMap.putProperty(Property.Key, Property)
func (m *PropertyMap) PutPropertyByKey(name Property_Key, newProp *Property) {
	switch name.Kind() {
	case Property_KeyKind_STRING:
		m.PutProperty(name.String_(), newProp)
	case Property_KeyKind_SYMBOL:
		m.PutPropertyBySymbol(name.Symbol(), newProp)
	}
}

// Java: PropertyMap.values()
func (m *PropertyMap) Values() func(yield func(*Property) bool) {
	m.ensureSorted()
	return func(yield func(*Property) bool) {
		for i := 0; i < len(m.entries()); i++ {
			if !yield(m.entries()[i]) {
				return
			}
		}
	}
}

// Java: PropertyMap.hashCode()
func (m *PropertyMap) HashCode() int32 {
	// We need to override hashCode so that JSType.hashCode is consistent.
	// JSType uses the hashcode of its PropertyMap. Otherwise PropertyMap uses
	// identity equality.
	//
	// Calculate the hash just based on the property names, not their types.
	// Otherwise we can get into an infinite loop because the ObjectType
	// hashCode method calls this one.
	var propHash int32
	for _, p := range m.entries() {
		propHash += jutil.StringHashCode(rhinoStringFromID(p.nameID))
	}
	if m.ext == nil || m.ext.knownSymbols == nil {
		return propHash
	}
	var symHash int32
	for s := range m.ext.knownSymbols.Keys() {
		symHash += jutil.ObjectHashCode(s)
	}
	return 31*(31*1+propHash) + symHash
}

// Equals is Object.equals (PropertyMap does not override it): identity.
func (m *PropertyMap) Equals(o any) bool {
	return any(m) == o
}

// Java: PropertyMap.IncrementCachedKeySetCounter()
func (m *PropertyMap) IncrementCachedKeySetCounter() {
	m.cachedKeySetCounter++
	if m.ext != nil {
		m.ext.cachedKeySet = nil
	}

	jutil.CheckState(m.cachedKeySetCounter >= 0)
}
