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

func (u *UnionType) Members() []JSType { return u.members }

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
func (u *UnionType) MeetWith(that JSType) JSType {
	return UnionGetGreatestSubtype(u, that)
}
