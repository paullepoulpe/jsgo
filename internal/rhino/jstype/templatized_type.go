// Port of com/google/javascript/rhino/jstype/TemplatizedType.java
// @ closure-compiler v20260415 (narrow; only the structural-template
// surface jsgo currently consumes).
//
// A TemplatizedType pairs a generic ObjectType (e.g. Array, Map)
// with a concrete substitution for its template parameters
// (`Array<number>`, `Map<string, Foo>`). At inference time we
// produce one whenever a JSDoc expression names a generic type with
// `<…>` slots; at call-site unification we recurse into its
// arguments so a `(Array<T>) → T` parameter binds T from the actual
// element type.
//
// Inheritance: when @extends carries template args
// (`class Foo<T> extends Bar<Array<T>>`), GetPropertyType walks the
// Base's templatized supertype with this instantiation's Args
// substituted — mirrors Java's
// ObjectType.getPropertyType's recursion through
// getImplicitPrototype() on a TemplatizedType receiver.

package jstype

import "strings"

// TemplatizedType is a generic instantiation: an underlying ObjectType
// (Array, Map, etc.) plus the per-position type arguments. Falls back
// to the underlying type for property lookup when the args don't
// affect a particular property.
type TemplatizedType struct {
	Base    *ObjectType
	Args    []JSType
}

func NewTemplatizedType(base *ObjectType, args []JSType) *TemplatizedType {
	return &TemplatizedType{Base: base, Args: args}
}

func (t *TemplatizedType) TypeName() string {
	if len(t.Args) == 0 {
		return t.Base.TypeName()
	}
	names := make([]string, len(t.Args))
	for i, a := range t.Args {
		if a == nil {
			names[i] = "?"
		} else {
			names[i] = a.TypeName()
		}
	}
	return t.Base.TypeName() + "<" + strings.Join(names, ",") + ">"
}

// IsSubtypeOf handles subclass-to-templatized-target checks by
// threading the substitution map up through the @extends chain. When
// `t` and `other` share a Base, args are compared positionally
// (invariant — Java's default for unspecified @template variance).
// When `t.Base` is a subclass of `other.Base`, we walk
// SubstitutedSupertype so the check reaches a same-Base comparison
// with args substituted against this type's map.
func (t *TemplatizedType) IsSubtypeOf(o JSType) bool {
	if other, ok := o.(*TemplatizedType); ok {
		if t.Base == other.Base {
			for i := 0; i < len(t.Args) && i < len(other.Args); i++ {
				if !IsSubtype(t.Args[i], other.Args[i]) {
					return false
				}
			}
			return true
		}
		// Different Bases: climb t's substituted supertype chain
		// looking for other.Base.
		if parent := t.SubstitutedSupertype(); parent != nil {
			if parent.IsSubtypeOf(other) {
				return true
			}
		}
		// Fall back to the base-level nominal walk — preserves the
		// pre-inheritance behaviour for bare generic targets.
		if !t.Base.IsSubtypeOf(other.Base) {
			return false
		}
		for i := 0; i < len(t.Args) && i < len(other.Args); i++ {
			if !IsSubtype(t.Args[i], other.Args[i]) {
				return false
			}
		}
		return true
	}
	return t.Base.IsSubtypeOf(o)
}

// substitutionMap pairs the Base's template parameters with this
// instantiation's Args. Returns nil when the Base has no template
// params wired or the Args don't line up — callers then treat the
// type as opaque. Mirrors the internal TemplateTypeMap built in
// TemplatizedType's Java constructor.
func (t *TemplatizedType) substitutionMap() *TemplateTypeMap {
	params := t.Base.TemplateParams
	if len(params) == 0 || len(t.Args) == 0 {
		return nil
	}
	m := NewTemplateTypeMap()
	n := len(params)
	if len(t.Args) < n {
		n = len(t.Args)
	}
	for i := 0; i < n; i++ {
		m.Bind(params[i], t.Args[i])
	}
	return m
}

// SubstitutedSupertype returns the Base's templatized @extends link
// with this type's substitution map applied — if Foo<T> extends
// Bar<Array<T>> and this is Foo<number>, returns Bar<Array<number>>.
// Returns nil when the Base has no templatized super (and so a
// nominal walk via Super() suffices). Mirrors Java's
// TemplatizedType.getImplicitPrototype, which returns the
// Base.getImplicitPrototype() re-substituted against the caller's
// own template map.
func (t *TemplatizedType) SubstitutedSupertype() *TemplatizedType {
	if t.Base == nil {
		return nil
	}
	parentTpl := t.Base.TemplatizedSuper()
	if parentTpl == nil {
		return nil
	}
	m := t.substitutionMap()
	if m == nil {
		return parentTpl
	}
	args := make([]JSType, len(parentTpl.Args))
	for i, a := range parentTpl.Args {
		args[i] = m.Substitute(a)
	}
	return NewTemplatizedType(parentTpl.Base, args)
}

// GetPropertyType looks up `name` on the Base's members, then walks
// the (substituted) templatized supertype when absent. Falls back to
// the Base's own super chain when no templatized link exists. The
// returned type already has this instantiation's template bindings
// substituted, so callers don't need to re-run substitution — this is
// what makes inherited `@return T` slots read out as the concrete arg
// at a generic call site. Mirrors Java's
// ObjectType.getPropertyType(String) when the receiver is a
// TemplatizedType.
func (t *TemplatizedType) GetPropertyType(name string) JSType {
	if t.Base == nil {
		return nil
	}
	m := t.substitutionMap()
	// Directly-declared members on Base come first — same order as
	// ObjectType.GetPropertyType.
	if pt := directMember(t.Base, name); pt != nil {
		if m != nil {
			return m.Substitute(pt)
		}
		return pt
	}
	// Templatized super wins over the raw nominal super when present:
	// we must propagate substitution as we climb.
	if parentTpl := t.SubstitutedSupertype(); parentTpl != nil {
		return parentTpl.GetPropertyType(name)
	}
	if t.Base.Super() != nil {
		pt := t.Base.Super().GetPropertyType(name)
		if pt != nil && m != nil {
			return m.Substitute(pt)
		}
		return pt
	}
	return nil
}

// directMember reads the property slot declared directly on o, with
// no super-chain walk. Kept out of object_type.go because its only
// caller is the templatized-type code path.
func directMember(o *ObjectType, name string) JSType {
	if o == nil {
		return nil
	}
	if t, ok := o.members[name]; ok {
		return t
	}
	return nil
}
