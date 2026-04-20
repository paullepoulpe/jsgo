// Port of com/google/javascript/rhino/jstype/UnionType.java @ closure-compiler v20260415
//
// Unions in Closure are commutative, associative, and collapse identical
// members. The Java implementation is much richer (supports typeclass
// refinements and intersection checks); this Tier-1/2 port covers
// primitive and function members only, which is what TypeCheckTest's
// `@type {A|B}` / `@param {A|B}` cases exercise.
//
// Three builder-time invariants, lifted from UnionType$Builder in Java:
// flatten nested unions, de-dup alternates, and sort alphabetically by
// TypeName so diagnostics spell `(A|B)` the same way every run. The
// canonical order was confirmed against oracle outputs (primitives
// alphabetical, nominal names alphabetical, anonymous records and
// functions last via their `{...}` / `function(...)` ASCII prefix).

package jstype

import (
	"sort"
	"strings"
)

type UnionType struct {
	members []JSType
}

// NewUnion flattens nested unions, drops duplicates, and sorts alternates
// deterministically. Returns a bare member when the union collapses to
// one element. Mirrors JSTypeRegistry.createUnionType + UnionType$Builder.
func NewUnion(members ...JSType) JSType {
	var flat []JSType
	for _, m := range members {
		if m == nil {
			continue
		}
		if u, ok := m.(*UnionType); ok {
			flat = append(flat, u.members...)
			continue
		}
		flat = append(flat, m)
	}

	unique := make([]JSType, 0, len(flat))
	for _, m := range flat {
		dup := false
		for _, k := range unique {
			if sameAlternate(k, m) {
				dup = true
				break
			}
		}
		if !dup {
			unique = append(unique, m)
		}
	}

	if len(unique) == 1 {
		return unique[0]
	}

	sort.SliceStable(unique, func(i, j int) bool {
		return unique[i].TypeName() < unique[j].TypeName()
	})
	return &UnionType{members: unique}
}

// sameAlternate is the dedup predicate. Nominal kinds (ObjectType,
// FunctionType, TemplateType) compare by pointer identity via the
// interface `a == b` check — matching Java's reference-equality on
// interned types. TemplatizedType needs a structural compare so two
// independently-built `Array<number>` values collapse. Everything else
// (primitives with value receivers) compares by TypeName.
func sameAlternate(a, b JSType) bool {
	if a == b {
		return true
	}
	if x, ok := a.(*TemplatizedType); ok {
		y, ok := b.(*TemplatizedType)
		if !ok || x.Base != y.Base || len(x.Args) != len(y.Args) {
			return false
		}
		for i := range x.Args {
			if !sameAlternate(x.Args[i], y.Args[i]) {
				return false
			}
		}
		return true
	}
	switch a.(type) {
	case *ObjectType, *FunctionType, *TemplateType:
		// Pointer identity already failed above; these are nominal.
		return false
	}
	return a.TypeName() == b.TypeName()
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
