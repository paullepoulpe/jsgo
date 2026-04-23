// Port of com/google/javascript/rhino/jstype/JSType.java @ closure-compiler v20260415
//
// The Java hierarchy is an abstract class `JSType` with many subclasses.
// In Go we expose an interface and concrete types; subtype logic is
// centralised in isSubtypeOf. Primitive type definitions live in per-type
// files matching the Java layout (all_type.go, unknown_type.go, …).

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

func isAll(t JSType) bool     { _, ok := t.(allType); return ok }
func isUnknown(t JSType) bool { _, ok := t.(unknownType); return ok }

// Canonical interned instances. Java JSTypeRegistry returns the same
// instance for each primitive; we mirror that with package-level vars.
var (
	All      JSType = allType{}
	Unknown  JSType = unknownType{}
	Number   JSType = numberType{}
	String   JSType = stringType{}
	Boolean  JSType = booleanType{}
	Null     JSType = nullType{}
	Void     JSType = voidType{}
)
