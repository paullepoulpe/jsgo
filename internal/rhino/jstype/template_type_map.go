// Port of com/google/javascript/rhino/jstype/TemplateTypeMap.java
// @ closure-compiler v20260415
//
// Manages a mapping from TemplateType to its resolved JSType.  Provides
// utility methods for cloning/extending the map.  The full Java class
// tracks submap boundaries (one submap per generic class in an
// inheritance chain) via a BitSet; here we use a plain []bool.
//
// Subset ported: all public methods plus the internal helpers needed to
// implement them.  The JSTypeRegistry parameter is accepted for
// compatibility but only the getNativeType(UNKNOWN_TYPE) path is
// exercised — bounded generics were removed upstream so getBound()
// always returns Unknown.
//
// TODO(template-bounds): port getBound / bounded-generics support when
// @template T extends Base syntax is exercised (Tier-5+).
// TODO(registry-integration): wire to the real Registry.getNativeType
// when JSTypeNative enum is ported (Tier-4+).

package jstype

// TemplateTypeMap manages a mapping from TemplateType to its resolved
// JSType.  Mirrors com.google.javascript.rhino.jstype.TemplateTypeMap.
//
// The map stores two parallel slices:
//   - templateKeys   — the ordered TemplateType keys (one per @template var)
//   - templateValues — the substituted JSType per key (length ≤ len(keys))
//
// Keys without a corresponding value are "unfilled"; they default to the
// key's bound (Unknown in the current port, since bounded generics are
// removed upstream).
//
// subMapStarts is a boolean slice of the same length as templateKeys.
// A true entry at index i means that index starts a new "submap" — i.e.
// the template keys declared directly on one class in the @extends chain.
// This is used to prevent premature substitution of sibling type vars
// during TemplateTypeMap construction.
type TemplateTypeMap struct {
	// templateKeys are the TemplateType keys of the map.
	templateKeys []*TemplateType
	// templateValues are the JSType values, index-aligned with templateKeys.
	// These are the unresolved values as provided at construction time.
	templateValues []JSType
	// resolvedTemplateValues are iteratively type-resolved values; used for
	// equivalence checks between two TemplateTypeMap instances.
	resolvedTemplateValues []JSType
	// subMapStarts[i] == true means index i starts a new submap.
	subMapStarts []bool
	// registry is kept for getNativeType compatibility; currently unused in
	// the narrow port (Unknown / bound resolution is inlined).
	registry *Registry
}

// createEmpty creates the empty TemplateTypeMap used as the registry
// sentinel.  Mirrors TemplateTypeMap.createEmpty in Java.
func createEmpty(registry *Registry) *TemplateTypeMap {
	return newTemplateTypeMap(registry, nil, nil, nil)
}

// newTemplateTypeMap is the private constructor.  Mirrors
// TemplateTypeMap(JSTypeRegistry, ImmutableList, ImmutableList, BitSet)
// in Java.
func newTemplateTypeMap(
	registry *Registry,
	templateKeys []*TemplateType,
	templateValues []JSType,
	subMapStarts []bool,
) *TemplateTypeMap {
	if len(templateValues) > len(templateKeys) {
		panic("templateValues.size() must be <= templateKeys.size()")
	}
	if len(subMapStarts) < len(templateKeys) {
		// Extend subMapStarts to at least len(templateKeys), padding with false.
		ext := make([]bool, len(templateKeys))
		copy(ext, subMapStarts)
		subMapStarts = ext
	}

	m := &TemplateTypeMap{
		registry:       registry,
		templateKeys:   templateKeys,
		templateValues: templateValues,
		subMapStarts:   subMapStarts,
	}

	// Iteratively resolve any JSType values that refer to the TemplateType
	// keys of this TemplateTypeMap — mirrors the constructor body in Java.
	replacer := forTotalReplacement(registry, m)

	nValues := len(m.templateValues)
	nKeys := len(m.templateKeys)
	resolved := make([]JSType, nKeys)
	for i := 0; i < nKeys; i++ {
		if i < nValues {
			nextSub := m.nextSubMapStart(i + 1)
			templateKey := m.templateKeys[i]
			replacer.setKeyType(templateKey, nextSub)
			templateValue := m.templateValues[i]
			resolved[i] = visitJSType(replacer, templateValue)
		} else {
			resolved[i] = m.templateKeys[i].getBound()
		}
	}
	m.resolvedTemplateValues = resolved
	return m
}

// nextSubMapStart returns the index of the next submap start at or after
// index, or len(templateKeys) if there is no later submap start.
// Mirrors TemplateTypeMap.nextSubMapStart in Java.
func (m *TemplateTypeMap) nextSubMapStart(index int) int {
	for index < len(m.templateKeys) && !m.subMapStarts[index] {
		index++
	}
	return index
}

// CopyFilledWithValues creates a new map in which any unfilled values in
// this map have been filled with values.  If there are fewer values than
// unfilled keys, Unknown is used for the rest.
// Mirrors TemplateTypeMap.copyFilledWithValues in Java.
func (m *TemplateTypeMap) CopyFilledWithValues(values []JSType) *TemplateTypeMap {
	requiredUnknownCount := m.numUnfilledTemplateKeys() - len(values)
	if requiredUnknownCount < 0 {
		panic("too many values")
	}
	if m.numUnfilledTemplateKeys() == 0 {
		return m
	}

	extendedValues := make([]JSType, 0, len(m.templateValues)+len(values)+requiredUnknownCount)
	extendedValues = append(extendedValues, m.templateValues...)
	extendedValues = append(extendedValues, values...)
	padToSameLengthWith(m.templateKeys, &extendedValues)

	return newTemplateTypeMap(m.registry, m.templateKeys, extendedValues, m.subMapStarts)
}

// CopyWithExtensionMap creates a new map in which the keys and values
// have been extended by extension.  Before extension, any unfilled values
// in the initial map will be filled with Unknown.
// Mirrors TemplateTypeMap.copyWithExtension(TemplateTypeMap) in Java.
func (m *TemplateTypeMap) CopyWithExtensionMap(extension *TemplateTypeMap) *TemplateTypeMap {
	return m.copyWithExtensionRaw(
		extension.templateKeys,
		extension.templateValues,
		extension.subMapStarts,
	)
}

// CopyWithExtension creates a new map in which the keys and values have
// been extended by keys and values respectively.  Before extension, any
// unfilled values in the initial map will be filled with Unknown.
// Mirrors TemplateTypeMap.copyWithExtension(ImmutableList, ImmutableList) in Java.
func (m *TemplateTypeMap) CopyWithExtension(keys []*TemplateType, values []JSType) *TemplateTypeMap {
	if len(keys) == 0 {
		return m.copyWithExtensionRaw(keys, values, nil)
	}
	sub := make([]bool, len(keys))
	sub[0] = true
	return m.copyWithExtensionRaw(keys, values, sub)
}

// copyWithExtensionRaw is the private three-argument form, mirroring
// TemplateTypeMap.copyWithExtension(ImmutableList, ImmutableList, BitSet) in Java.
func (m *TemplateTypeMap) copyWithExtensionRaw(
	keys []*TemplateType,
	values []JSType,
	newSubMapStarts []bool,
) *TemplateTypeMap {
	extendedUnfilledCount := len(keys) - len(values)
	if extendedUnfilledCount < 0 {
		panic("values exceed keys in extension")
	}
	if m.numUnfilledTemplateKeys() == 0 && len(keys) == 0 {
		return m
	}

	extendedKeys := make([]*TemplateType, 0, len(m.templateKeys)+len(keys))
	extendedKeys = append(extendedKeys, m.templateKeys...)
	extendedKeys = append(extendedKeys, keys...)

	extendedValues := make([]JSType, 0, len(m.templateValues)+len(values))
	extendedValues = append(extendedValues, m.templateValues...)
	padToSameLengthWith(m.templateKeys, &extendedValues)
	extendedValues = append(extendedValues, values...)

	totalLen := len(m.templateKeys) + len(keys)
	extSub := make([]bool, totalLen)
	copy(extSub, m.subMapStarts)
	for i := 0; i < len(keys); i++ {
		if newSubMapStarts != nil && i < len(newSubMapStarts) {
			extSub[i+len(m.templateKeys)] = newSubMapStarts[i]
		}
	}

	return newTemplateTypeMap(m.registry, extendedKeys, extendedValues, extSub)
}

// CopyWithoutKeys creates a new map in which keys contained in removals
// are eliminated.  Only unfilled keys may be removed.
// Mirrors TemplateTypeMap.copyWithoutKeys in Java.
func (m *TemplateTypeMap) CopyWithoutKeys(removals map[*TemplateType]bool) *TemplateTypeMap {
	newKeys := make([]*TemplateType, 0, len(m.templateKeys))
	newSub := make([]bool, 0, len(m.templateKeys))

	// Keep all filled keys as-is.
	for i := 0; i < len(m.templateValues); i++ {
		newKeys = append(newKeys, m.templateKeys[i])
		newSub = append(newSub, m.subMapStarts[i])
	}

	newKeysSize := 0
	inNewSubmap := false
	for i := len(m.templateValues); i < len(m.templateKeys); i++ {
		key := m.templateKeys[i]
		inNewSubmap = inNewSubmap || m.subMapStarts[i]
		if !removals[key] {
			newKeys = append(newKeys, key)
			newSub = append(newSub, inNewSubmap)
			newKeysSize++
			inNewSubmap = false
		}
	}

	if len(newKeys) == len(m.templateKeys) {
		return m
	}
	return newTemplateTypeMap(m.registry, newKeys, m.templateValues, newSub)
}

// Size returns the number of template keys in the map.
func (m *TemplateTypeMap) Size() int { return len(m.templateKeys) }

// IsEmpty returns true if the map has no template keys.
func (m *TemplateTypeMap) IsEmpty() bool { return len(m.templateKeys) == 0 }

// GetTemplateKeys returns all template keys in the map.
func (m *TemplateTypeMap) GetTemplateKeys() []*TemplateType { return m.templateKeys }

// GetTemplateValues returns all template values in the map.
func (m *TemplateTypeMap) GetTemplateValues() []JSType { return m.templateValues }

// HasTemplateKey reports whether this map contains the specified template
// key (using pointer identity, matching Java's identical() check).
// Mirrors TemplateTypeMap.hasTemplateKey(TemplateType) in Java.
func (m *TemplateTypeMap) HasTemplateKey(templateKey *TemplateType) bool {
	return m.hasTemplateKeyFrom(templateKey, 0)
}

// hasTemplateKeyFrom reports whether this map contains templateKey, only
// considering keys at indices >= nextSubMap.  Mirrors
// TemplateTypeMap.hasTemplateKey(TemplateType, int) in Java.
func (m *TemplateTypeMap) hasTemplateKeyFrom(templateKey *TemplateType, nextSubMap int) bool {
	start := nextSubMap
	if start < 0 {
		start = 0
	}
	for i := start; i < len(m.templateKeys); i++ {
		if m.templateKeys[i] == templateKey { // pointer identity
			return true
		}
	}
	return false
}

// GetTemplateKeyCountThisShouldAlwaysBeOneOrZeroButIsnt counts the
// occurrences of templateKey in the map.  Java TODO b/139230800.
func (m *TemplateTypeMap) GetTemplateKeyCountThisShouldAlwaysBeOneOrZeroButIsnt(templateKey *TemplateType) int {
	count := 0
	for _, k := range m.templateKeys {
		if k == templateKey {
			count++
		}
	}
	return count
}

// numUnfilledTemplateKeys returns the number of keys without a value.
func (m *TemplateTypeMap) numUnfilledTemplateKeys() int {
	return len(m.templateKeys) - len(m.templateValues)
}

// HasTemplateType reports whether there is a JSType value associated
// with the specified template key.
// Mirrors TemplateTypeMap.hasTemplateType in Java.
func (m *TemplateTypeMap) HasTemplateType(key *TemplateType) bool {
	return m.getTemplateTypeIndex(key) != -1
}

// GetUnresolvedOriginalTemplateType returns the unresolved JSType value
// for key, or Unknown if no value is associated.
// Mirrors TemplateTypeMap.getUnresolvedOriginalTemplateType in Java.
func (m *TemplateTypeMap) GetUnresolvedOriginalTemplateType(key *TemplateType) JSType {
	index := m.getTemplateTypeIndex(key)
	if index == -1 {
		return Unknown
	}
	return m.templateValues[index]
}

// GetLastTemplateTypeKeyByName returns the last TemplateType key whose
// name equals keyName, or nil if none.
// Mirrors TemplateTypeMap.getLastTemplateTypeKeyByName in Java.
func (m *TemplateTypeMap) GetLastTemplateTypeKeyByName(keyName string) *TemplateType {
	for i := len(m.templateKeys) - 1; i >= 0; i-- {
		if m.templateKeys[i].Name == keyName {
			return m.templateKeys[i]
		}
	}
	return nil
}

// getTemplateTypeIndex returns the last index in templateKeys (up to
// len(templateValues)) where key appears, or -1 if not found.
// Mirrors TemplateTypeMap.getTemplateTypeIndex in Java.
func (m *TemplateTypeMap) getTemplateTypeIndex(key *TemplateType) int {
	maxIndex := len(m.templateValues)
	if len(m.templateKeys) < maxIndex {
		maxIndex = len(m.templateKeys)
	}
	for i := maxIndex - 1; i >= 0; i-- {
		if m.templateKeys[i] == key {
			return i
		}
	}
	return -1
}

// GetResolvedTemplateType returns the resolved JSType for key.  If no
// value is bound, returns the default (Unknown in the current port).
// Mirrors TemplateTypeMap.getResolvedTemplateType in Java.
func (m *TemplateTypeMap) GetResolvedTemplateType(key *TemplateType) JSType {
	index := m.getTemplateTypeIndex(key)
	if index == -1 {
		return defaultValueType(key)
	}
	return m.resolvedTemplateValues[index]
}

// HasAnyTemplateTypesInternal reports whether any resolved value still
// contains an unresolved TemplateType.
// Mirrors TemplateTypeMap.hasAnyTemplateTypesInternal in Java.
func (m *TemplateTypeMap) HasAnyTemplateTypesInternal() bool {
	for _, tv := range m.resolvedTemplateValues {
		if hasAnyTemplateTypes(tv) {
			return true
		}
	}
	return false
}

// String renders the map for diagnostics.  Mirrors
// TemplateTypeMap.toString in Java.
func (m *TemplateTypeMap) String() string {
	s := "{ "
	for i, k := range m.templateKeys {
		s += "("
		s += k.Name
		s += ","
		if i < len(m.templateValues) {
			s += typeName(m.templateValues[i])
		}
		s += ","
		if i < len(m.resolvedTemplateValues) {
			s += typeName(m.resolvedTemplateValues[i])
		}
		s += ") "
	}
	s += "}"
	return s
}

// padToSameLengthWith pads *values with defaultValueType entries until
// len(*values) == len(keys).
func padToSameLengthWith(keys []*TemplateType, values *[]JSType) {
	for i := len(*values); i < len(keys); i++ {
		*values = append(*values, defaultValueType(keys[i]))
	}
}

// defaultValueType returns the default fill value for a template key.
// Since bounded generics were removed upstream, this always returns
// Unknown.  Mirrors TemplateTypeMap.defaultValueType in Java.
func defaultValueType(t *TemplateType) JSType {
	bound := t.getBound()
	if isUnknown(bound) {
		return Unknown
	}
	return t
}

// GetIndicesOfSubmapsForTesting returns the indices where submaps begin.
// Mirrors TemplateTypeMap.getIndicesOfSubmapsForTesting in Java.
func (m *TemplateTypeMap) GetIndicesOfSubmapsForTesting() []int {
	var out []int
	for i, b := range m.subMapStarts {
		if b {
			out = append(out, i)
		}
	}
	return out
}

// hasAnyTemplateTypes reports whether t or any of its component types
// is an unresolved TemplateType.  Mirrors JSType.hasAnyTemplateTypes
// in Java (structural walk).
func hasAnyTemplateTypes(t JSType) bool {
	if t == nil {
		return false
	}
	switch v := t.(type) {
	case *TemplateType:
		return true
	case *UnionType:
		for _, m := range v.Members() {
			if hasAnyTemplateTypes(m) {
				return true
			}
		}
	case *TemplatizedType:
		for _, a := range v.Args {
			if hasAnyTemplateTypes(a) {
				return true
			}
		}
	case *FunctionType:
		for _, p := range v.Params {
			if hasAnyTemplateTypes(p) {
				return true
			}
		}
		if hasAnyTemplateTypes(v.Return) {
			return true
		}
	}
	return false
}

// typeName returns the TypeName() of t, or "" if t is nil.
func typeName(t JSType) string {
	if t == nil {
		return ""
	}
	return t.TypeName()
}

// getBound returns the upper bound of the TemplateType.  Since bounded
// generics were removed in the upstream Java source, this always returns
// Unknown.
// TODO(template-bounds): wire to real bound when @template T extends Base
// is ported.
func (t *TemplateType) getBound() JSType { return Unknown }

