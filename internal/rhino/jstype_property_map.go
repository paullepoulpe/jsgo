// Port of com/google/javascript/rhino/jstype/PropertyMap.java @ closure-compiler v20260415

package rhino

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// PropertyMap is the representation for a collection of properties on an
// object.
type PropertyMap struct {
	// A place to get the inheritance structure. Because the extended
	// interfaces are resolved dynamically, this gets messy :(. If
	// type-resolution was more well-defined, we could just reference primary
	// parents and secondary parents directly.
	parentSource ObjectType

	// The map of our own properties (a TreeMap, ordered by String.compareTo;
	// the empty map is an ImmutableMap).
	properties *jutil.TreeMap[string, *Property]

	knownSymbols *jutil.ObjectLinkedHashMap[*KnownSymbolType, *Property] // lazily initialized

	// The set of keys for this map and its ancestors.
	//
	// Collecting the set of properties turns out to be expensive for some
	// high volume callers (e.g. structural type equality). Since the results
	// don't change often, this cache eliminates most of the cost.
	//
	// Go: nil is Java null; the sorted keys otherwise.
	cachedKeySet []string

	// The set of known symbols keys for this map and its ancestors.
	cachedKnownSymbolsKeySet []*KnownSymbolType

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
	immutable bool
}

// Java: PropertyMap.EMPTY_MAP
var propertyMap_EMPTY_MAP = func() *PropertyMap {
	m := newPropertyMapWithMaps(
		jutil.NewTreeMap[string, *Property](jutil.StringCompareTo),
		jutil.NewObjectLinkedHashMap[*KnownSymbolType, *Property]())
	m.immutable = true
	return m
}()

// Java: PropertyMap()
func NewPropertyMap() *PropertyMap {
	return newPropertyMapWithMaps(jutil.NewTreeMap[string, *Property](jutil.StringCompareTo), nil)
}

// Java: PropertyMap(Map<String, Property>, Map<KnownSymbolType, Property>)
func newPropertyMapWithMaps(
	underlyingMap *jutil.TreeMap[string, *Property],
	underlyingknownSymbols *jutil.ObjectLinkedHashMap[*KnownSymbolType, *Property]) *PropertyMap {
	return &PropertyMap{properties: underlyingMap, knownSymbols: underlyingknownSymbols}
}

// Java: PropertyMap.immutableEmptyMap()
func PropertyMap_ImmutableEmptyMap() *PropertyMap {
	return propertyMap_EMPTY_MAP
}

// Java: PropertyMap.setParentSource(ObjectType)
func (m *PropertyMap) SetParentSource(ownerType ObjectType) {
	if m == propertyMap_EMPTY_MAP {
		return
	}

	m.parentSource = ownerType
	m.incrementCachedKeySetCounter()
}

// Java: PropertyMap.setParentForTesting(PropertyMap)
func (m *PropertyMap) SetParentForTesting(parent *PropertyMap) {
	m.parentSource = parent.parentSource
	m.incrementCachedKeySetCounter()
}

// GetPrimaryParent returns the direct parent of this property map.
//
// Java: PropertyMap.getPrimaryParent()
func (m *PropertyMap) GetPrimaryParent() *PropertyMap {
	if m.parentSource == nil {
		return nil
	}
	iProto := m.parentSource.GetImplicitPrototype()
	if iProto == nil {
		return nil
	}
	return iProto.GetPropertyMap()
}

// getSecondaryParentObjects returns the secondary parents of this property
// map, for interfaces that need multiple inheritance or for interfaces of
// abstract classes.
//
// Java: PropertyMap.getSecondaryParentObjects()
func (m *PropertyMap) getSecondaryParentObjects() []ObjectType {
	if m.parentSource == nil {
		return nil
	}
	if m.parentSource.GetConstructor() != nil && m.parentSource.GetConstructor().IsAbstract() {
		return m.parentSource.GetConstructor().GetOwnImplementedInterfaces()
	}
	return m.parentSource.GetCtorExtendedInterfaces()
}

// Java: PropertyMap.findClosest(String)
func (m *PropertyMap) FindClosest(name string) *Property_OwnedProperty {
	return m.FindClosestByKey(NewProperty_StringKey(name))
}

// Java: PropertyMap.findClosest(Property.Key)
func (m *PropertyMap) FindClosestByKey(name Property_Key) *Property_OwnedProperty {
	// Check primary parents which always has precendence over secondary.
	for pm := m; pm != nil; pm = pm.GetPrimaryParent() {
		prop := pm.GetOwnPropertyByKey(name)
		if prop != nil {
			return NewProperty_OwnedProperty(pm.parentSource, prop)
		}
	}

	// Recurse into secondary parents.
	for pm := m; pm != nil; pm = pm.GetPrimaryParent() {
		for _, o := range pm.getSecondaryParentObjects() {
			parent := o.GetPropertyMap()
			if parent != nil {
				e := parent.FindClosestByKey(name)
				if e != nil {
					return e
				}
			}
		}
	}

	return nil
}

// Java: PropertyMap.getOwnProperty(String)
func (m *PropertyMap) GetOwnProperty(propertyName string) *Property {
	return m.properties.Get(propertyName)
}

// Java: PropertyMap.getOwnProperty(Property.Key)
func (m *PropertyMap) GetOwnPropertyByKey(propertyName Property_Key) *Property {
	switch propertyName.Kind() {
	case Property_KeyKind_STRING:
		return m.GetOwnProperty(propertyName.String_())
	}
	if m.knownSymbols != nil {
		return m.knownSymbols.Get(propertyName.Symbol())
	}
	return nil
}

// Java: PropertyMap.getPropertiesCount()
func (m *PropertyMap) GetPropertiesCount() int {
	primaryParent := m.GetPrimaryParent()
	if primaryParent == nil {
		return m.properties.Size()
	}

	return len(m.GetAllKeys().StringKeys())
}

// GetOwnPropertyNames returns properties.keySet(), in key order.
//
// Go-only: language-adapter — a snapshot of the live keySet() view.
//
// Java: PropertyMap.getOwnPropertyNames()
func (m *PropertyMap) GetOwnPropertyNames() []string {
	return m.properties.KeySlice()
}

// Java: PropertyMap.getOwnKnownSymbols()
func (m *PropertyMap) GetOwnKnownSymbols() []*KnownSymbolType {
	if m.knownSymbols != nil {
		return m.knownSymbols.KeySlice()
	}
	return nil
}

// Java: PropertyMap.getAllKeys()
func (m *PropertyMap) GetAllKeys() *PropertyMap_AllKeys {
	ancestors := jutil.NewLinkedHashSet[*PropertyMap]()
	m.collectAllAncestors(ancestors)

	var maxAncestorCounter int32 = 0
	for ancestor := range ancestors.All() {
		if ancestor.cachedKeySetCounter > maxAncestorCounter {
			maxAncestorCounter = ancestor.cachedKeySetCounter
		}
	}

	// If any counter is greater than this counter, there has been a mutation
	// and the cache must be rebuilt.
	if maxAncestorCounter != m.cachedKeySetCounter || m.cachedKeySet == nil {
		keys := jutil.NewTreeMap[string, bool](jutil.StringCompareTo)
		knownSymbolsKeys := jutil.NewObjectLinkedHashSet[*KnownSymbolType]()
		for ancestor := range ancestors.All() {
			// Update the counters in all ancestors.
			//
			// This update scheme is convergent. As long as there are no
			// mutations, calls #keySet will eventually set stable caches on all
			// maps.
			ancestor.cachedKeySetCounter = maxAncestorCounter
			ancestor.cachedKeySet = nil
			ancestor.cachedKnownSymbolsKeySet = nil

			for _, k := range ancestor.GetOwnPropertyNames() {
				keys.Put(k, true)
			}
			for _, s := range ancestor.GetOwnKnownSymbols() {
				knownSymbolsKeys.Add(s)
			}
		}
		m.cachedKeySet = keys.KeySlice()
		m.cachedKnownSymbolsKeySet = knownSymbolsKeys.ToSlice()
	}

	return NewPropertyMap_AllKeys(m.cachedKeySet, m.cachedKnownSymbolsKeySet)
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

// Java: PropertyMap.collectAllAncestors(LinkedHashSet<PropertyMap>)
func (m *PropertyMap) collectAllAncestors(ancestors *jutil.LinkedHashSet[*PropertyMap]) {
	if !ancestors.Add(m) {
		return
	}

	primaryParent := m.GetPrimaryParent()
	if primaryParent != nil {
		primaryParent.collectAllAncestors(ancestors)
	}

	for _, parentType := range m.getSecondaryParentObjects() {
		parentMap := parentType.GetPropertyMap()
		if parentMap != nil {
			parentMap.collectAllAncestors(ancestors)
		}
	}
}

// Java: PropertyMap.putProperty(String, Property)
func (m *PropertyMap) PutProperty(name string, newProp *Property) {
	oldProp := m.properties.Get(name)

	if oldProp == nil {
		// The cache is only invalidated if this is a new property name.
		m.incrementCachedKeySetCounter()
	} else {
		// This is to keep previously inferred JsDoc info, e.g., in a
		// replaceScript scenario.
		newProp.SetJSDocInfo(oldProp.GetJSDocInfo())
	}

	if m.immutable {
		panic(jutil.NewUnsupportedOperationException())
	}
	m.properties.Put(name, newProp)
}

// Java: PropertyMap.putProperty(KnownSymbolType, Property)
func (m *PropertyMap) PutPropertyBySymbol(symbol *KnownSymbolType, newProp *Property) {
	var oldProp *Property
	if m.knownSymbols != nil {
		oldProp = m.knownSymbols.Get(symbol)
	}

	if oldProp == nil {
		// The cache is only invalidated if this is a new property name.
		m.incrementCachedKeySetCounter()
	}

	if m.knownSymbols == nil {
		m.knownSymbols = jutil.NewObjectLinkedHashMap[*KnownSymbolType, *Property]()
	}
	if m.immutable {
		panic(jutil.NewUnsupportedOperationException())
	}
	m.knownSymbols.Put(symbol, newProp)
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
	return m.properties.Values()
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
	if m.knownSymbols == nil {
		return jutil.SetHashCode(m.properties.KeySlice())
	}
	return 31*(31*1+jutil.SetHashCode(m.properties.KeySlice())) +
		jutil.SetHashCode(m.knownSymbols.KeySlice())
}

// Equals is Object.equals (PropertyMap does not override it): identity.
func (m *PropertyMap) Equals(o any) bool {
	return any(m) == o
}

// Java: PropertyMap.incrementCachedKeySetCounter()
func (m *PropertyMap) incrementCachedKeySetCounter() {
	m.cachedKeySetCounter++
	m.cachedKeySet = nil

	jutil.CheckState(m.cachedKeySetCounter >= 0)
}
