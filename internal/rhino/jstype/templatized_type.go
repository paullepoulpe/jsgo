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

package jstype

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// TemplatizedType is a generic instantiation: an underlying ObjectType
// (Array, Map, etc.) plus the per-position type arguments. Falls back
// to the underlying type for property lookup when the args don't
// affect a particular property.
type TemplatizedType struct {
	Base    *ObjectType
	Args    []JSType
}

func NewTemplatizedType(base *ObjectType, args []JSType) *TemplatizedType {
	t := &TemplatizedType{Base: base, Args: args}
	if base != nil {
		base.templatizedSelf = t
	}
	return t
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

// IsInstantiatedGenericType reports whether this TemplatizedType has at
// least one template argument.  Mirrors TemplatizedType.isInstantiatedGenericType
// in Java.
func (t *TemplatizedType) IsInstantiatedGenericType() bool {
	return len(t.Args) > 0
}

// GetTemplateArg returns the resolved type argument for the given template
// key, looked up by pointer identity through the base object's type
// parameter list.  Returns nil if key is not a parameter of the base type.
// Mirrors TemplatizedType.getTemplateTypes() (keyed version) in Java.
func (t *TemplatizedType) GetTemplateArg(key *TemplateType) JSType {
	if key == nil || t.Base == nil {
		return nil
	}
	// Walk the base's super chain looking for a TemplatizedType supertype
	// so we can resolve inherited template parameters.
	params := t.Base.GetTypeParameters()
	for i, p := range params {
		if p == key && i < len(t.Args) {
			return t.Args[i]
		}
	}
	// Fall through to inherited templatized supertypes: if the supertype is
	// itself a TemplatizedType, delegate to it.
	if t.Base.super != nil {
		superTemplatized := findTemplatizedSuper(t.Base.super)
		if superTemplatized != nil {
			return superTemplatized.GetTemplateArg(key)
		}
	}
	return nil
}

// findTemplatizedSuper returns the first TemplatizedType encountered
// while walking up the super chain of obj, or nil.
func findTemplatizedSuper(obj *ObjectType) *TemplatizedType {
	for cur := obj; cur != nil; cur = cur.super {
		if cur.templatizedSelf != nil {
			return cur.templatizedSelf
		}
	}
	return nil
}

// GetTemplateTypes mirrors TemplatizedType.getTemplateTypes() in Java.
// Returns the ordered list of template type arguments for this specialization.
func (t *TemplatizedType) GetTemplateTypes() []JSType {
	return t.Args
}

// ToMaybeTemplatizedType mirrors JSType.toMaybeTemplatizedType() — returns
// this TemplatizedType itself. Mirrors the Java override.
func (t *TemplatizedType) ToMaybeTemplatizedType() *TemplatizedType {
	return t
}

// GetReferencedType mirrors TemplatizedType.getReferencedType() in Java —
// returns the underlying ObjectType (the unparameterized base type).
func (t *TemplatizedType) GetReferencedType() *ObjectType {
	return t.Base
}

// MatchesObjectContext delegates to the base ObjectType. Mirrors
// TemplatizedType (ProxyObjectType) matchesObjectContext().
func (t *TemplatizedType) MatchesObjectContext() bool {
	return true
}

// GetPossibleToBooleanOutcomes delegates to the base ObjectType.
func (t *TemplatizedType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return BooleanLiteralSetTrue
}

// GetJSDocInfo mirrors ObjectType.getJSDocInfo() for TemplatizedType by
// delegating to the base ObjectType.
func (t *TemplatizedType) GetJSDocInfo() *rhino.JSDocInfo {
	if t.Base != nil {
		return t.Base.GetJSDocInfo()
	}
	return nil
}

// GetPropertyType mirrors TemplatizedType.getPropertyType(name) — looks up
// property on the base ObjectType. Mirrors Java's ProxyObjectType delegation.
func (t *TemplatizedType) GetPropertyType(name string) JSType {
	if t.Base != nil {
		return t.Base.GetPropertyType(name)
	}
	return nil
}

// HasProperty mirrors TemplatizedType.hasProperty(name).
func (t *TemplatizedType) HasProperty(name string) bool {
	return t.GetPropertyType(name) != nil
}

// IsSubtypeOf delegates to the underlying ObjectType — assignment
// compatibility with concrete classes ignores template arguments
// (Java's behaviour for invariant generics with structural wrapping).
//
// Template-arg comparison strips null/undefined before checking
// subtype, mirroring Java's TemplateTypeMap.checkEquivalenceHelper
// which always passes SubtypingMode.IGNORE_NULL_UNDEFINED when
// recursing into invariant template arguments — this is the
// "modulo null/undef" rule that holds independent of the outer
// subtyping mode (see TypeCheckTest.testModuloNullUndefThatWorked
// WithoutSpecialSubtypingRules*).
func (t *TemplatizedType) IsSubtypeOf(o JSType) bool {
	if other, ok := o.(*TemplatizedType); ok {
		if !t.Base.IsSubtypeOf(other.Base) {
			return false
		}
		for i := 0; i < len(t.Args) && i < len(other.Args); i++ {
			a := stripNullUndef(t.Args[i])
			b := stripNullUndef(other.Args[i])
			if !IsSubtype(a, b) {
				return false
			}
		}
		return true
	}
	return t.Base.IsSubtypeOf(o)
}

// stripNullUndef removes Null and Void members from a UnionType, returning
// the remaining type. Mirrors the null/undef stripping Java applies via
// SubtypingMode.IGNORE_NULL_UNDEFINED when comparing template arguments.
func stripNullUndef(t JSType) JSType {
	if t == nil {
		return t
	}
	if t == Null || t == Void {
		return t
	}
	u, ok := t.(*UnionType)
	if !ok {
		return t
	}
	var members []JSType
	for _, m := range u.members {
		if m != Null && m != Void {
			members = append(members, m)
		}
	}
	if len(members) == 0 {
		return t
	}
	if len(members) == 1 {
		return members[0]
	}
	return &UnionType{members: members}
}
