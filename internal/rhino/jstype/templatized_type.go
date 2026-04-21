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
// Recursive substitution through nested generics (class Foo<T>
// extends Bar<List<T>> — substituting T → string should descend into
// Bar<List<string>>) is handled by TemplateTypeMap.Substitute, which
// recurses into each Arg.

package jstype

import "strings"

// TemplatizedType is a generic instantiation: an underlying ObjectType
// (Array, Map, etc.) plus the per-position type arguments. Falls back
// to the underlying type for property lookup when the args don't
// affect a particular property.
type TemplatizedType struct {
	Base *ObjectType
	Args []JSType
}

// NewTemplatizedType pairs a base ObjectType with concrete type arguments.
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

// IsSubtypeOf delegates to the underlying ObjectType — assignment
// compatibility with concrete classes ignores template arguments
// (Java's behaviour for invariant generics with structural wrapping).
func (t *TemplatizedType) IsSubtypeOf(o JSType) bool {
	if other, ok := o.(*TemplatizedType); ok {
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
