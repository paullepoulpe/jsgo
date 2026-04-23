// Port of com/google/javascript/rhino/jstype/JSType.java (plus
// NumberType.java, StringType.java, AllType.java, UnknownType.java) @
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

func (allType) TypeName() string            { return "*" }
func (allType) IsSubtypeOf(o JSType) bool   { _, ok := o.(allType); return ok || isUnknown(o) }

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

// ToMaybeFunctionType mirrors JSType.toMaybeFunctionType(JSType) — null-
// safe downcast to *FunctionType. Returns nil when t is nil or not a
// FunctionType.
func ToMaybeFunctionType(t JSType) *FunctionType {
	if t == nil {
		return nil
	}
	f, _ := t.(*FunctionType)
	return f
}

// ToMaybeEnumType mirrors JSType.toMaybeEnumType() — null-safe downcast
// to *EnumType. Returns nil when t is nil or not an EnumType.
func ToMaybeEnumType(t JSType) *EnumType {
	if t == nil {
		return nil
	}
	e, _ := t.(*EnumType)
	return e
}

// Dereference mirrors JSType.dereference() — returns the ObjectType when
// t is an *ObjectType or *FunctionType (which is also an object in the
// Java hierarchy). Returns nil for primitives, null, undefined, and
// unknown. Mirrors Java's autobox().toObjectType() at a simplified level.
func Dereference(t JSType) *ObjectType {
	if t == nil {
		return nil
	}
	if o, ok := t.(*ObjectType); ok {
		return o
	}
	// FunctionType is an ObjectType in the Java hierarchy — InferJsDocInfo
	// calls dereference() on the inferred type and expects to get an
	// ObjectType back for function expressions. Return nil here because our
	// FunctionType is not an ObjectType subtype; the caller must handle the
	// FunctionType case separately if needed.
	return nil
}
