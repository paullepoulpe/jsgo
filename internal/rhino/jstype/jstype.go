// Port of com/google/javascript/rhino/jstype/JSType.java (plus
// NumberType.java, StringType.java, AllType.java, UnknownType.java,
// NoType.java, NoResolvedType.java, NoObjectType.java) @
// closure-compiler v20260415
//
// The Java hierarchy is an abstract class `JSType` with many subclasses.
// In Go we expose an interface and concrete types; subtype logic is
// centralised in isSubtypeOf. Only the types Tier-0 exercises are here.

package jstype

type JSType interface {
	// TypeName matches Java JSType.toString() for primitive types (e.g.
	// "number"). Not used for conformance, only for diagnostics.
	TypeName() string

	// IsSubtypeOf ports JSType.isSubtypeOf(JSType). Semantics match Java:
	// every type is a subtype of AllType and Unknown; primitives are
	// subtypes only of themselves.
	IsSubtypeOf(other JSType) bool
}

// IsSubtype is the package-level entry point; callers should prefer this
// over t.IsSubtypeOf so union-type handling stays in one place. Mirrors
// how JSType.isSubtype in Java dispatches through the registry.
func IsSubtype(x, y JSType) bool {
	if x == nil || y == nil {
		return true
	}
	// Bottom types (NoType, NoObjectType, NoResolvedType) are subtypes of
	// everything — mirrors Java's NoType.isSubtype returning true
	// unconditionally. Placed ahead of TemplateType so "no" dominates.
	if isBottom(x) {
		return true
	}
	// Unresolved template vars accept any value — mirrors Java's
	// "open type" rule during unification. Once the substitution
	// resolves at a CALL site, the concrete type is used instead,
	// and this branch stops applying.
	if _, ok := y.(*TemplateType); ok {
		return true
	}
	if _, ok := x.(*TemplateType); ok {
		return true
	}
	// TemplatizedType ↔ underlying nominal: a `Foo<X>` is
	// assignment-compatible with a bare `Foo` (Java's invariant
	// generic rule degrades the args when one side has none).
	if xt, ok := x.(*TemplatizedType); ok {
		if y == xt.Base {
			return true
		}
	}
	if yt, ok := y.(*TemplatizedType); ok {
		if x == yt.Base {
			return true
		}
	}
	// When the target is a union, X matches if it's a subtype of any
	// member. Left-side unions stay in UnionType.IsSubtypeOf.
	if u, ok := y.(*UnionType); ok {
		if _, isUnion := x.(*UnionType); isUnion {
			return x.IsSubtypeOf(u)
		}
		for _, m := range u.members {
			if IsSubtype(x, m) {
				return true
			}
		}
		return false
	}
	return x.IsSubtypeOf(y)
}

type allType struct{}

func (allType) TypeName() string          { return "*" }
func (allType) IsSubtypeOf(o JSType) bool { _, ok := o.(allType); return ok || isUnknown(o) }

type unknownType struct{}

func (unknownType) TypeName() string          { return "?" }
func (unknownType) IsSubtypeOf(o JSType) bool { return true }

type numberType struct{}

func (numberType) TypeName() string { return "number" }
func (numberType) IsSubtypeOf(o JSType) bool {
	if _, ok := o.(numberType); ok {
		return true
	}
	return isAll(o) || isUnknown(o)
}

type stringType struct{}

func (stringType) TypeName() string { return "string" }
func (stringType) IsSubtypeOf(o JSType) bool {
	if _, ok := o.(stringType); ok {
		return true
	}
	return isAll(o) || isUnknown(o)
}

type booleanType struct{}

func (booleanType) TypeName() string { return "boolean" }
func (booleanType) IsSubtypeOf(o JSType) bool {
	if _, ok := o.(booleanType); ok {
		return true
	}
	return isAll(o) || isUnknown(o)
}

type nullType struct{}

func (nullType) TypeName() string { return "null" }
func (nullType) IsSubtypeOf(o JSType) bool {
	if _, ok := o.(nullType); ok {
		return true
	}
	return isAll(o) || isUnknown(o)
}

type voidType struct{}

func (voidType) TypeName() string { return "undefined" }
func (voidType) IsSubtypeOf(o JSType) bool {
	if _, ok := o.(voidType); ok {
		return true
	}
	return isAll(o) || isUnknown(o)
}

// noType is Closure's NoType — the bottom of the lattice. A subtype of
// every type; nothing (except itself) is a subtype of it. Produced by
// impossible-type refinements (e.g. `x` known to be both string and
// number). Java's JSType.NoType.
type noType struct{}

func (noType) TypeName() string          { return "None" }
func (noType) IsSubtypeOf(o JSType) bool { return true }

// noObjectType is Closure's NoObjectType — the bottom of the object
// lattice. Subtype of every ObjectType; used for things like the
// inferred type of `{}` before any property adds. Java's
// JSType.NoObjectType.
type noObjectType struct{}

func (noObjectType) TypeName() string          { return "NoObject" }
func (noObjectType) IsSubtypeOf(o JSType) bool { return true }

// noResolvedType is Closure's NoResolvedType — the type assigned to a
// NamedType reference whose target never resolved (typo'd `@type
// {Fooo}`). Behaves like NoType for subtype purposes but carries a
// distinct diagnostic name. Java's JSType.NoResolvedType.
type noResolvedType struct{}

func (noResolvedType) TypeName() string          { return "NoResolvedType" }
func (noResolvedType) IsSubtypeOf(o JSType) bool { return true }

func isAll(t JSType) bool     { _, ok := t.(allType); return ok }
func isUnknown(t JSType) bool { _, ok := t.(unknownType); return ok }

// isBottom returns true for any of the three Java "no" types. Used by
// IsSubtype so a bottom on the left passes before the TemplateType
// shortcut kicks in.
func isBottom(t JSType) bool {
	switch t.(type) {
	case noType, noObjectType, noResolvedType:
		return true
	}
	return false
}

// Canonical interned instances. Java JSTypeRegistry returns the same
// instance for each primitive; we mirror that with package-level vars.
var (
	All            JSType = allType{}
	Unknown        JSType = unknownType{}
	Number         JSType = numberType{}
	String         JSType = stringType{}
	Boolean        JSType = booleanType{}
	Null           JSType = nullType{}
	Void           JSType = voidType{}
	NoType         JSType = noType{}
	NoObjectType   JSType = noObjectType{}
	NoResolvedType JSType = noResolvedType{}
)
