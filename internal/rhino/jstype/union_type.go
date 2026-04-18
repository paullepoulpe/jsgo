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

// NewUnion flattens nested unions and drops duplicates. Returns a bare
// member when the union collapses to one element (matches Java's
// `JSTypeRegistry.createUnionType` behaviour).
func NewUnion(members ...JSType) JSType {
	seen := map[string]JSType{}
	for _, m := range members {
		if u, ok := m.(*UnionType); ok {
			for _, n := range u.members {
				seen[n.TypeName()] = n
			}
			continue
		}
		seen[m.TypeName()] = m
	}
	if len(seen) == 1 {
		for _, v := range seen {
			return v
		}
	}
	var out []JSType
	for _, v := range seen {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].TypeName() < out[j].TypeName()
	})
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
