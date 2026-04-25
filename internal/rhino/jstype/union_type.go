// Port of com/google/javascript/rhino/jstype/UnionType.java @ closure-compiler v20260415
//
// Unions in Closure are commutative, associative, and collapse identical
// members. The Java implementation is much richer (supports typeclass
// refinements and intersection checks); this Tier-1/2 port covers
// primitive and function members only, which is what TypeCheckTest's
// `@type {A|B}` / `@param {A|B}` cases exercise.

package jstype

import (
	"sort"
	"strings"
)

type UnionType struct {
	members []JSType
}

// NewUnion flattens nested unions and drops duplicates. Returns a
// bare member when the union collapses to one element. Dedup keys
// combine TypeName with pointer identity so two anonymous structural
// records (both TypeName "{...}") stay distinct. Mirrors
// JSTypeRegistry.createUnionType interning semantics.
// Go-only: language-adapter
func NewUnion(members ...JSType) JSType {
	type key struct {
		name string
		ptr  any
	}
	seen := map[key]JSType{}
	var order []key
	add := func(m JSType) {
		var k key
		switch v := m.(type) {
		case *ObjectType:
			k = key{name: v.TypeName(), ptr: v}
		case *TemplatizedType:
			k = key{name: v.TypeName(), ptr: v}
		case *FunctionType:
			k = key{name: v.TypeName(), ptr: v}
		default:
			k = key{name: m.TypeName()}
		}
		if _, ok := seen[k]; ok {
			return
		}
		seen[k] = m
		order = append(order, k)
	}
	for _, m := range members {
		if u, ok := m.(*UnionType); ok {
			for _, n := range u.members {
				add(n)
			}
			continue
		}
		add(m)
	}
	if len(seen) == 1 {
		for _, v := range seen {
			return v
		}
	}
	sort.Slice(order, func(i, j int) bool {
		return order[i].name < order[j].name
	})
	out := make([]JSType, 0, len(order))
	for _, k := range order {
		out = append(out, seen[k])
	}
	return &UnionType{members: out}
}

// Go-only: language-adapter
func (u *UnionType) Members() []JSType { return u.members }

// GetAlternates returns the slice of member types in this union.
// Mirrors UnionType.getAlternates() in Java.
func (u *UnionType) GetAlternates() []JSType { return u.members }

// Go-only: language-adapter
func (u *UnionType) TypeName() string {
	names := make([]string, len(u.members))
	for i, m := range u.members {
		names[i] = m.TypeName()
	}
	return "(" + strings.Join(names, "|") + ")"
}

// IsSubtypeOf: a union X is a subtype of Y iff every member of X is a
// subtype of Y. (The asymmetric case, X a single type and Y a union, is
// handled in the package-level IsSubtype.)
func (u *UnionType) IsSubtypeOf(o JSType) bool {
	for _, m := range u.members {
		if !IsSubtype(m, o) {
			return false
		}
	}
	return true
}

// RestrictByNotNullOrUndefined strips null and undefined from the union's
// alternates and returns the collapsed result. Mirrors
// UnionType.restrictByNotNullOrUndefined() in Java, which delegates to
// mapTypes(JSType::restrictByNotNullOrUndefined, this): each alternate's
// own restrictByNotNullOrUndefined is called and the results are re-unioned.
// For primitive null/void the Java override returns NoType (bottom);
// for all other types it returns self. The net effect for a union is that
// null and undefined members are removed and the rest are kept.
func (u *UnionType) RestrictByNotNullOrUndefined() JSType {
	b := NewUnionTypeBuilder()
	for _, m := range u.members {
		if isNullType(m) || isVoidType(m) {
			continue // those two restrict to NoType, contributing nothing
		}
		b.AddAlternate(m)
	}
	return b.Build()
}

// UnionGetGreatestSubtype returns the greatest common subtype (lattice meet)
// of a UnionType and an arbitrary JSType. Mirrors the static
// UnionType.getGreatestSubtype(UnionType, JSType) in Java.
//
// Algorithm:
//  1. Collect all alternates of union that are subtypes of that.
//  2. If that is itself a union, also collect its alternates that are
//     subtypes of union.
//  3. Otherwise, if that is a subtype of union, add that.
//
// Returns NoType when the intersection is empty unless both sides are
// object types, in which case NoObjectType is returned instead.
// Go-only: language-adapter
func UnionGetGreatestSubtype(union *UnionType, that JSType) JSType {
	b := NewUnionTypeBuilder()

	for _, alt := range union.members {
		if IsSubtype(alt, that) {
			b.AddAlternate(alt)
		}
	}

	if tu, ok := that.(*UnionType); ok {
		for _, otherAlt := range tu.members {
			if IsSubtype(otherAlt, union) {
				b.AddAlternate(otherAlt)
			}
		}
	} else if IsSubtype(that, union) {
		b.AddAlternate(that)
	}

	result := b.Build()
	if !isNoType(result) {
		return result
	}
	// Both sides are objects → return the no-object bottom.
	unionIsObj := AllTypesMatchUnion(func(m JSType) bool {
		_, isObj := m.(*ObjectType)
		return isObj
	}, union)
	if unionIsObj {
		if _, ok := that.(*ObjectType); ok {
			return NoObject
		}
	}
	return None
}

// MeetWith returns the lattice meet (greatest common subtype) of this
// union and another type. Mirrors JSType.meetWith / the RelationshipVisitor
// dispatch that ultimately calls UnionType.getGreatestSubtype in Java.
// Go-only: language-adapter
func (u *UnionType) MeetWith(that JSType) JSType {
	return UnionGetGreatestSubtype(u, that)
}

// CollapseUnion returns the least non-union supertype of u.
// If any alternate is unknown → Unknown.
// If all are objects → common super (approximated as Unknown when no single
// super is known). If mixed values and objects → All.
// Mirrors UnionType.collapseUnion() @ closure-compiler v20260415.
func (u *UnionType) CollapseUnion() JSType {
	var currentValue JSType
	var currentObj *ObjectType
	for _, a := range u.members {
		if isUnknown(a) {
			return Unknown
		}
		obj, isObj := a.(*ObjectType)
		if !isObj {
			if currentValue == nil && currentObj == nil {
				currentValue = a
			} else {
				return All // mixed or multiple values
			}
		} else if currentValue != nil {
			return All // value + object
		} else if currentObj == nil {
			currentObj = obj
		} else {
			// Two different objects without registry — conservatively return Unknown.
			return Unknown
		}
	}
	if currentObj != nil {
		return currentObj
	}
	return currentValue
}

// MatchesObjectContext reports whether any alternate of u matches an object
// context. Mirrors UnionType.matchesObjectContext() @ closure-compiler v20260415.
func (u *UnionType) MatchesObjectContext() bool {
	return AnyTypeMatchesUnion(MatchesObjectContext, u)
}

// GetPossibleToBooleanOutcomes returns the union of boolean outcomes across
// all alternates. Mirrors UnionType.getPossibleToBooleanOutcomes().
func (u *UnionType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	result := BooleanLiteralSetEmpty
	for _, m := range u.members {
		result = result.Union(GetPossibleToBooleanOutcomes(m))
		if result == BooleanLiteralSetBoth {
			break
		}
	}
	return result
}

// ToMaybeUnionType mirrors JSType.toMaybeUnionType() — returns this union
// itself. Mirrors the Java UnionType override.
func (u *UnionType) ToMaybeUnionType() *UnionType { return u }

// Contains mirrors UnionType.contains(JSType) — returns true when the
// given type is a member of this union (by TypeName equality for value
// types, by pointer identity for object types). Mirrors Java's
// anyTypeMatches(type::equals, this).
func (u *UnionType) Contains(t JSType) bool {
	tName := t.TypeName()
	for _, m := range u.members {
		if m == t || m.TypeName() == tName {
			return true
		}
	}
	return false
}

// GetRestrictedUnion mirrors UnionType.getRestrictedUnion(JSType) in Java.
// Returns a new union with all subtypes of `restrictType` removed, keeping
// unknown/unresolved members. Mirrors Java's logic exactly.
func (u *UnionType) GetRestrictedUnion(restrictType JSType) JSType {
	b := NewUnionTypeBuilder()
	for _, m := range u.members {
		// Keep unknown types and any type that is NOT a subtype of restrictType.
		if isUnknown(m) || !IsSubtype(m, restrictType) {
			b.AddAlternate(m)
		}
	}
	return b.Build()
}

// RestrictByNotUndefined strips undefined (void) from the union's alternates.
// Mirrors UnionType.restrictByNotUndefined() in Java.
func (u *UnionType) RestrictByNotUndefined() JSType {
	b := NewUnionTypeBuilder()
	for _, m := range u.members {
		if isVoidType(m) {
			continue
		}
		b.AddAlternate(m)
	}
	return b.Build()
}

// RestrictByNotNull strips null from the union's alternates.
// Mirrors UnionType.restrictByNotNull() in Java.
func (u *UnionType) RestrictByNotNull() JSType {
	b := NewUnionTypeBuilder()
	for _, m := range u.members {
		if isNullType(m) {
			continue
		}
		b.AddAlternate(m)
	}
	return b.Build()
}

// IsNullable mirrors UnionType.isNullable() — true when any alternate is
// nullable (i.e. has null in its type). Mirrors Java's
// anyTypeMatches(JSType::isNullable, this).
func (u *UnionType) IsNullable() bool {
	for _, m := range u.members {
		if isNullType(m) {
			return true
		}
	}
	return false
}

// IsVoidable mirrors UnionType.isVoidable() — true when any alternate
// includes undefined/void. Mirrors Java's anyTypeMatches(JSType::isVoidable, this).
func (u *UnionType) IsVoidable() bool {
	for _, m := range u.members {
		if isVoidType(m) {
			return true
		}
		if inner, ok := m.(*UnionType); ok {
			if inner.IsVoidable() {
				return true
			}
		}
	}
	return false
}

// HasPropertyKind describes whether a property is absent, present on all
// non-null/undefined union members, or present on only some of them.
// Mirrors Java's JSType.HasPropertyKind enum.
type HasPropertyKind int

const (
	HasPropertyAbsent       HasPropertyKind = iota // property not found on any member
	HasPropertyKnownPresent                        // property found on all non-null/void members
	HasPropertyMaybePresent                        // property found on some but not all members
)

// propertyHaver is satisfied by *ObjectType (and *RecordType via embedding).
type propertyHaver interface {
	HasProperty(name string) bool
}

// GetPropertyKind mirrors Java's UnionType.getPropertyKind(String) — it
// iterates the union's members (skipping null/void) and returns:
//   - HasPropertyKnownPresent when every non-null/void member has the property,
//   - HasPropertyMaybePresent when some but not all do,
//   - HasPropertyAbsent when none do.
//
// Mirrors UnionType.getPropertyKind(Property.Key, boolean) @ closure-compiler v20260415.
func (u *UnionType) GetPropertyKind(name string) HasPropertyKind {
	found := false
	always := true
	for _, m := range u.members {
		if isNullType(m) || isVoidType(m) {
			continue
		}
		if ph, ok := m.(propertyHaver); ok {
			if ph.HasProperty(name) {
				found = true
			} else {
				always = false
			}
		} else {
			// Member is not an object type (e.g. primitive); treat as absent.
			always = false
		}
	}
	if !found {
		return HasPropertyAbsent
	}
	if always {
		return HasPropertyKnownPresent
	}
	return HasPropertyMaybePresent
}

// IsObject mirrors UnionType.isObject() — true when all alternates are
// object types. Mirrors Java's allTypesMatch(JSType::isObject, this).
func (u *UnionType) IsObject() bool {
	for _, m := range u.members {
		switch m.(type) {
		case *ObjectType, *FunctionType, *TemplatizedType:
			// ok
		default:
			return false
		}
	}
	return len(u.members) > 0
}

// findPropertyTypeOwner is satisfied by types that can resolve a named
// property to its JSType (returning nil when absent).
type findPropertyTypeOwner interface {
	GetPropertyType(name string) JSType
}

// FindPropertyType ports UnionType.findPropertyTypeWithoutConsideringTemplateTypes
// (rhino/jstype/UnionType.java:213). Iterates alternates, skipping null/void
// and members lacking the property, and returns the least-supertype of the
// per-member property types. Returns nil when no member has the property.
func (u *UnionType) FindPropertyType(propertyName string) JSType {
	var propertyType JSType
	for i := 0; i < len(u.members); i++ {
		alternate := u.members[i]
		// Filter out the null/undefined type — matches Java's `isNullType
		// || isVoidType` skip. Java does this so a `?T` (T|null) read of a
		// property defined on T doesn't lose the property.
		if isNullType(alternate) || isVoidType(alternate) {
			continue
		}
		var altPropertyType JSType
		if owner, ok := alternate.(findPropertyTypeOwner); ok {
			altPropertyType = owner.GetPropertyType(propertyName)
		}
		if altPropertyType == nil {
			continue
		}
		if propertyType == nil {
			propertyType = altPropertyType
		} else {
			// Java does propertyType.getLeastSupertype(altPropertyType);
			// our rough equivalent is unioning the two and letting
			// NewUnion's dedup collapse identical members.
			propertyType = NewUnion(propertyType, altPropertyType)
		}
	}
	return propertyType
}

// CanBeCalled ports UnionType.canBeCalled(): true iff every alternate can
// be called. Mirrors `allTypesMatch(JSType::canBeCalled, this)`.
func (u *UnionType) CanBeCalled() bool {
	for _, m := range u.members {
		// Function types and Unknown are callable; anything else isn't.
		switch m.(type) {
		case *FunctionType:
			continue
		default:
			if isUnknown(m) {
				continue
			}
			return false
		}
	}
	return len(u.members) > 0
}

// IsStruct ports UnionType.isStruct() — true when any alternate is @struct.
// Mirrors `anyTypeMatches(JSType::isStruct, this)`.
func (u *UnionType) IsStruct() bool {
	for _, m := range u.members {
		if obj, ok := m.(*ObjectType); ok && obj.IsStruct {
			return true
		}
	}
	return false
}

// IsDict ports UnionType.isDict() — true when any alternate is @dict.
// Mirrors `anyTypeMatches(JSType::isDict, this)`.
func (u *UnionType) IsDict() bool {
	for _, m := range u.members {
		if obj, ok := m.(*ObjectType); ok && obj.IsDict {
			return true
		}
	}
	return false
}

// IsExplicitlyVoidable ports UnionType.isExplicitlyVoidable() — true when
// any alternate is the explicit Void type (not the Unknown wildcard).
// Mirrors `anyTypeMatches(JSType::isExplicitlyVoidable, this)`.
func (u *UnionType) IsExplicitlyVoidable() bool {
	for _, m := range u.members {
		if isVoidType(m) {
			return true
		}
	}
	return false
}

// IsUnknownType ports UnionType.isUnknownType() — true when any alternate
// is Unknown. Mirrors `anyTypeMatches(JSType::isUnknownType, this)`.
func (u *UnionType) IsUnknownType() bool {
	for _, m := range u.members {
		if isUnknown(m) {
			return true
		}
	}
	return false
}

// MatchesNumberContext ports UnionType.matchesNumberContext() — true when
// any alternate may appear in a numeric context. Mirrors
// `anyTypeMatches(JSType::matchesNumberContext, this)`.
func (u *UnionType) MatchesNumberContext() bool {
	for _, m := range u.members {
		switch m.(type) {
		case numberType, bigIntType, booleanType, nullType, voidType:
			return true
		}
		if isUnknown(m) || isAll(m) {
			return true
		}
	}
	return false
}

// MatchesStringContext ports UnionType.matchesStringContext() — true when
// any alternate may convert to String. Mirrors
// `anyTypeMatches(JSType::matchesStringContext, this)`.
func (u *UnionType) MatchesStringContext() bool {
	// Java says: TODO reverse this logic — currently most types convert to
	// String, so this returns true for every non-empty union.
	return len(u.members) > 0
}

// MatchesSymbolContext ports UnionType.matchesSymbolContext() — true when
// any alternate is a Symbol. Mirrors `anyTypeMatches(JSType::matchesSymbolContext, this)`.
func (u *UnionType) MatchesSymbolContext() bool {
	for _, m := range u.members {
		if _, ok := m.(SymbolType); ok {
			return true
		}
		if isUnknown(m) || isAll(m) {
			return true
		}
	}
	return false
}

// HasAnyTemplateTypesInternal ports UnionType.hasAnyTemplateTypesInternal()
// — true when any alternate has @template parameters. Mirrors
// `anyTypeMatches(JSType::hasAnyTemplateTypes, this)`.
func (u *UnionType) HasAnyTemplateTypesInternal() bool {
	for _, m := range u.members {
		if _, ok := m.(*TemplateType); ok {
			return true
		}
		if t, ok := m.(*TemplatizedType); ok && len(t.Args) > 0 {
			for _, a := range t.Args {
				if _, ok := a.(*TemplateType); ok {
					return true
				}
			}
		}
	}
	return false
}

// TestForEquality ports UnionType.testForEquality(JSType) — returns the
// (Yes/No/Unknown) outcome of `this == that`. Mirrors the Java loop
// that aggregates per-alternate `testForEquality` calls and returns
// Unknown the moment two alternates disagree.
func (u *UnionType) TestForEquality(that JSType) eqResult {
	first := true
	var result eqResult
	for i := 0; i < len(u.members); i++ {
		alt := u.members[i]
		test := testForEquality(alt, that)
		if first {
			result = test
			first = false
		} else if result != test {
			return eqUnknown
		}
	}
	if first {
		return eqUnknown
	}
	return result
}

// GetLeastSupertype ports UnionType.getLeastSupertype(JSType) — the union
// override checks if `that` is a subtype of any alternate (excluding
// unknowns) and short-circuits to `this`; otherwise delegates to the
// generic `JSType.getLeastSupertype(this, that)`.
func (u *UnionType) GetLeastSupertype(that JSType) JSType {
	if !isUnknown(that) {
		if _, ok := that.(*UnionType); !ok {
			for i := 0; i < len(u.members); i++ {
				alt := u.members[i]
				if !isUnknown(alt) && IsSubtype(that, alt) {
					return u
				}
			}
		}
	}
	// Generic fallback: union of (this, that).
	return NewUnion(u, that)
}
