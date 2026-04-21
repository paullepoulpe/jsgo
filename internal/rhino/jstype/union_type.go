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

// NewUnion flattens nested unions (A | (B | C) → A | B | C, including
// deeper nestings), drops duplicates, and returns a bare member when
// the union collapses to one element. Members are sorted by TypeName
// so diagnostics are deterministic regardless of construction order.
// Dedup keys combine TypeName with pointer identity so two anonymous
// structural records (both TypeName "{...}") stay distinct. Mirrors
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
	// Recursive flatten: unions nested arbitrarily deep all collapse
	// to a single level. Matches UnionType.Builder.addAlternate, which
	// re-enters the builder for each member.
	var visit func(JSType)
	visit = func(m JSType) {
		if u, ok := m.(*UnionType); ok {
			for _, n := range u.members {
				visit(n)
			}
			return
		}
		add(m)
	}
	for _, m := range members {
		visit(m)
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
