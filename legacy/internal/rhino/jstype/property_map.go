// Port of com/google/javascript/rhino/jstype/PropertyMap.java @ closure-compiler v20260415
//
// PropertyMap holds the own-property table for an ObjectType and implements
// the prototype-chain traversal used by property lookup. This port covers:
//   - Own string-keyed and symbol-keyed property maps.
//   - Primary-parent traversal via ObjectType.getImplicitPrototype().
//   - Secondary-parent traversal for interfaces (@extends / @implements).
//   - A mutation-counter cache for the merged key set (getAllKeys).
//
// Differences from Java:
//   - ImmutableSortedSet<String> → []string (sorted slice, rebuilt on cache miss).
//   - ImmutableSet<KnownSymbolType> → []*KnownSymbolType (same approach).
//   - LinkedHashSet / LinkedHashMap → map[...] (Go map iteration order is
//     unspecified; insertion order is not preserved for knownSymbols either).
//   - checkState(counter >= 0) → panic on overflow kept as a TODO; counters
//     in practice never overflow for any realistic program.
//   - VisibleForTesting methods kept with the same names.

package jstype

import (
	"sort"
)

// PropertyMap mirrors PropertyMap.java — the property table for one object type.
type PropertyMap struct {
	// parentSource is the ObjectType whose implicit-prototype chain anchors
	// the parent lookup for this map. nil for the empty / root map.
	parentSource *ObjectType

	// properties holds own string-keyed properties.
	properties map[string]*Property

	// knownSymbols holds own well-known-symbol-keyed properties.
	// nil until the first symbol property is added.
	knownSymbols map[*KnownSymbolType]*Property

	// cachedKeySet / cachedKnownSymbolsKeySet cache the merged key sets
	// (own + all ancestors). Invalidated by incrementCachedKeySetCounter.
	cachedKeySet            []string           // sorted
	cachedKnownSymbolsKeySet []*KnownSymbolType // unordered

	// cachedKeySetCounter is bumped on every own mutation. A cached value is
	// valid only when this counter equals the maximum counter seen across all
	// ancestor maps.
	cachedKeySetCounter int
}

// emptyMap is the package-level singleton returned by ImmutableEmptyMap.
var emptyMap = &PropertyMap{
	properties:   map[string]*Property{},
	knownSymbols: map[*KnownSymbolType]*Property{},
}

// NewPropertyMap mirrors PropertyMap() — an empty mutable map backed by a
// regular (insertion-ordered-in-Java; unordered-in-Go) map.
func NewPropertyMap() *PropertyMap {
	return &PropertyMap{
		properties: map[string]*Property{},
	}
}

// newPropertyMapFromUnderlying mirrors the private PropertyMap(Map, Map)
// constructor.
func newPropertyMapFromUnderlying(props map[string]*Property, syms map[*KnownSymbolType]*Property) *PropertyMap {
	return &PropertyMap{
		properties:   props,
		knownSymbols: syms,
	}
}

// ImmutableEmptyMap mirrors PropertyMap.immutableEmptyMap() — returns a
// shared singleton empty map. Callers must not mutate the returned map.
func ImmutableEmptyMap() *PropertyMap { return emptyMap }

// SetParentSource mirrors PropertyMap.setParentSource(ObjectType).
func (pm *PropertyMap) SetParentSource(ownerType *ObjectType) {
	if pm == emptyMap {
		return
	}
	pm.parentSource = ownerType
	pm.incrementCachedKeySetCounter()
}

// SetParentForTesting mirrors PropertyMap.setParentForTesting(PropertyMap).
// VisibleForTesting.
func (pm *PropertyMap) SetParentForTesting(parent *PropertyMap) {
	pm.parentSource = parent.parentSource
	pm.incrementCachedKeySetCounter()
}

// GetPrimaryParent mirrors PropertyMap.getPrimaryParent() — the property map
// of the implicit prototype of parentSource. Returns nil at the root.
func (pm *PropertyMap) GetPrimaryParent() *PropertyMap {
	if pm.parentSource == nil {
		return nil
	}
	iProto := pm.parentSource.getImplicitPrototype()
	if iProto == nil {
		return nil
	}
	return iProto.GetPropertyMap()
}

// getSecondaryParentObjects mirrors PropertyMap.getSecondaryParentObjects().
// Returns the extended interfaces for parentSource (interface multiple
// inheritance).
//
// TODO(later-tier): implement getConstructor().isAbstract() +
// getOwnImplementedInterfaces() branch. For now only the extended-interfaces
// branch is wired.
func (pm *PropertyMap) getSecondaryParentObjects() []*ObjectType {
	if pm.parentSource == nil {
		return nil
	}
	return pm.parentSource.getCtorExtendedInterfaces()
}

// FindClosestByString mirrors PropertyMap.findClosest(String).
func (pm *PropertyMap) FindClosestByString(name string) *OwnedProperty {
	return pm.FindClosest(NewStringKey(name))
}

// FindClosest mirrors PropertyMap.findClosest(Property.Key).
func (pm *PropertyMap) FindClosest(name Key) *OwnedProperty {
	// Walk primary chain first (always takes precedence over secondary).
	for m := pm; m != nil; m = m.GetPrimaryParent() {
		prop := m.GetOwnProperty(name)
		if prop != nil {
			return NewOwnedProperty(m.parentSource, prop)
		}
	}

	// Recurse into secondary parents.
	for m := pm; m != nil; m = m.GetPrimaryParent() {
		for _, o := range m.getSecondaryParentObjects() {
			parent := o.GetPropertyMap()
			if parent != nil {
				e := parent.FindClosest(name)
				if e != nil {
					return e
				}
			}
		}
	}

	return nil
}

// GetOwnPropertyByString mirrors PropertyMap.getOwnProperty(String).
func (pm *PropertyMap) GetOwnPropertyByString(propertyName string) *Property {
	return pm.properties[propertyName]
}

// GetOwnProperty mirrors PropertyMap.getOwnProperty(Property.Key).
func (pm *PropertyMap) GetOwnProperty(propertyName Key) *Property {
	switch propertyName.Kind() {
	case KeyKindString:
		return pm.GetOwnPropertyByString(propertyName.StringVal())
	case KeyKindSymbol:
		if pm.knownSymbols == nil {
			return nil
		}
		return pm.knownSymbols[propertyName.SymbolVal()]
	default:
		return nil
	}
}

// GetPropertiesCount mirrors PropertyMap.getPropertiesCount(). When there is
// a parent, the merged (all-keys) string-key count is returned; otherwise
// just the own-property count.
func (pm *PropertyMap) GetPropertiesCount() int {
	if pm.GetPrimaryParent() == nil {
		return len(pm.properties)
	}
	return len(pm.GetAllKeys().StringKeys)
}

// GetOwnPropertyNames mirrors PropertyMap.getOwnPropertyNames().
// Returns the set of own string-keyed property names.
func (pm *PropertyMap) GetOwnPropertyNames() map[string]bool {
	out := make(map[string]bool, len(pm.properties))
	for k := range pm.properties {
		out[k] = true
	}
	return out
}

// GetOwnKnownSymbols mirrors PropertyMap.getOwnKnownSymbols().
func (pm *PropertyMap) GetOwnKnownSymbols() []*KnownSymbolType {
	if pm.knownSymbols == nil {
		return nil
	}
	out := make([]*KnownSymbolType, 0, len(pm.knownSymbols))
	for sym := range pm.knownSymbols {
		out = append(out, sym)
	}
	return out
}

// AllKeys mirrors PropertyMap.AllKeys — the merged key sets across all
// ancestors.
type AllKeys struct {
	StringKeys  []string           // sorted
	SymbolKeys  []*KnownSymbolType // unordered
}

// GetAllKeys mirrors PropertyMap.getAllKeys(). Rebuilds and caches the
// merged key set when any ancestor has been mutated since the last call.
func (pm *PropertyMap) GetAllKeys() AllKeys {
	// Collect all ancestors in BFS order.
	ancestors := pm.collectAllAncestors()

	// Find the maximum counter across all ancestors.
	maxCounter := 0
	for _, a := range ancestors {
		if a.cachedKeySetCounter > maxCounter {
			maxCounter = a.cachedKeySetCounter
		}
	}

	// If any counter is greater than ours, the cache is stale.
	if maxCounter != pm.cachedKeySetCounter || pm.cachedKeySet == nil {
		keySet := map[string]bool{}
		symSet := map[*KnownSymbolType]bool{}

		for _, a := range ancestors {
			// Propagate the counter to every ancestor so future calls
			// are fast as long as no one mutates. Mirrors Java's
			// convergent-counter scheme.
			a.cachedKeySetCounter = maxCounter
			a.cachedKeySet = nil
			a.cachedKnownSymbolsKeySet = nil

			for k := range a.properties {
				keySet[k] = true
			}
			for sym := range a.knownSymbols {
				symSet[sym] = true
			}
		}

		// Build sorted string slice.
		keys := make([]string, 0, len(keySet))
		for k := range keySet {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		pm.cachedKeySet = keys

		// Build symbol slice.
		syms := make([]*KnownSymbolType, 0, len(symSet))
		for sym := range symSet {
			syms = append(syms, sym)
		}
		pm.cachedKnownSymbolsKeySet = syms
	}

	return AllKeys{
		StringKeys: pm.cachedKeySet,
		SymbolKeys: pm.cachedKnownSymbolsKeySet,
	}
}

// collectAllAncestors mirrors PropertyMap.collectAllAncestors — returns all
// ancestor maps in BFS order (including self), deduplicated by pointer
// identity. Uses a slice-based visited set to preserve Java's LinkedHashSet
// iteration order (insertion order).
func (pm *PropertyMap) collectAllAncestors() []*PropertyMap {
	var ordered []*PropertyMap
	seen := map[*PropertyMap]bool{}

	var walk func(*PropertyMap)
	walk = func(m *PropertyMap) {
		if seen[m] {
			return
		}
		seen[m] = true
		ordered = append(ordered, m)

		if p := m.GetPrimaryParent(); p != nil {
			walk(p)
		}
		for _, parentType := range m.getSecondaryParentObjects() {
			if parentMap := parentType.GetPropertyMap(); parentMap != nil {
				walk(parentMap)
			}
		}
	}

	walk(pm)
	return ordered
}

// PutPropertyByString mirrors PropertyMap.putProperty(String, Property).
func (pm *PropertyMap) PutPropertyByString(name string, newProp *Property) {
	oldProp := pm.properties[name]
	if oldProp == nil {
		// New property name — invalidate cache.
		pm.incrementCachedKeySetCounter()
	} else {
		// Preserve previously inferred JSDoc (replaceScript scenario).
		newProp.SetJSDocInfo(oldProp.GetJSDocInfo())
	}
	pm.properties[name] = newProp
}

// PutPropertyBySymbol mirrors PropertyMap.putProperty(KnownSymbolType, Property).
func (pm *PropertyMap) PutPropertyBySymbol(symbol *KnownSymbolType, newProp *Property) {
	var oldProp *Property
	if pm.knownSymbols != nil {
		oldProp = pm.knownSymbols[symbol]
	}
	if oldProp == nil {
		pm.incrementCachedKeySetCounter()
	}
	if pm.knownSymbols == nil {
		pm.knownSymbols = map[*KnownSymbolType]*Property{}
	}
	pm.knownSymbols[symbol] = newProp
}

// PutProperty mirrors PropertyMap.putProperty(Property.Key, Property).
func (pm *PropertyMap) PutProperty(name Key, newProp *Property) {
	switch name.Kind() {
	case KeyKindString:
		pm.PutPropertyByString(name.StringVal(), newProp)
	case KeyKindSymbol:
		pm.PutPropertyBySymbol(name.SymbolVal(), newProp)
	}
}

// Values mirrors PropertyMap.values() — iterates over own string-keyed
// properties only (matches Java's properties.values()).
func (pm *PropertyMap) Values() []*Property {
	out := make([]*Property, 0, len(pm.properties))
	for _, p := range pm.properties {
		out = append(out, p)
	}
	return out
}

// HashCode mirrors PropertyMap.hashCode(). Based on property name sets only
// (not types) to avoid infinite recursion through ObjectType.hashCode.
func (pm *PropertyMap) HashCode() int {
	h := 17
	for k := range pm.properties {
		for _, c := range k {
			h = h*31 + int(c)
		}
	}
	for sym := range pm.knownSymbols {
		for _, c := range sym.GetDisplayName() {
			h = h*31 + int(c)
		}
	}
	return h
}

// incrementCachedKeySetCounter bumps the mutation counter and invalidates
// the cached key set. Mirrors PropertyMap.incrementCachedKeySetCounter().
func (pm *PropertyMap) incrementCachedKeySetCounter() {
	pm.cachedKeySetCounter++
	pm.cachedKeySet = nil
	// TODO: panic on overflow (checkState(counter >= 0)) if needed.
}

// SnapshotOwnProperties returns a point-in-time copy of this map's own
// string-keyed and symbol-keyed properties. The returned maps are
// independent of the original — mutations to the PropertyMap after
// the snapshot is taken do not affect the copy, and vice versa.
//
// This is used by interface-conflict diagnostics in SuperclassCheck and
// TypeValidator: those passes iterate over interface property slots while
// potentially modifying the live maps, so they operate on a snapshot to
// avoid concurrent-modification issues.
//
// There is no direct Java equivalent; in Java ImmutableMap views or a
// copy-on-enter pattern is used instead, but Go maps have no concurrent-
// modification guard so an explicit copy is necessary.
func (pm *PropertyMap) SnapshotOwnProperties() (strings map[string]*Property, symbols map[*KnownSymbolType]*Property) {
	strings = make(map[string]*Property, len(pm.properties))
	for k, v := range pm.properties {
		strings[k] = v
	}
	if pm.knownSymbols != nil {
		symbols = make(map[*KnownSymbolType]*Property, len(pm.knownSymbols))
		for k, v := range pm.knownSymbols {
			symbols[k] = v
		}
	}
	return strings, symbols
}
