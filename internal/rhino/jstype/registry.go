// Port of com/google/javascript/rhino/jstype/JSTypeRegistry.java @ closure-compiler v20260415
//
// The real registry manages interning of all JSType values, named-type
// resolution, template instantiation, and much more. Tier-0 only needs a
// name-to-primitive lookup, which we provide here.

package jstype

import "strings"

type Registry struct {
	// resolver is the JSTypeResolver that shares this registry's lifecycle.
	// Nil until Create() wires it; exposed via GetResolver() so JSTypeResolver
	// can detect when a resolver already exists (mirrors the Java factory method
	// JSTypeResolver.create(JSTypeRegistry) which calls registry.getResolver()).
	resolver *JSTypeResolver
}

func NewRegistry() *Registry { return &Registry{} }

// GetResolver returns the JSTypeResolver associated with this registry, or nil
// if none has been attached yet. Mirrors JSTypeRegistry.getResolver() in Java.
func (r *Registry) GetResolver() *JSTypeResolver { return r.resolver }

// ResolveName mirrors JSTypeRegistry.getType for primitive names. Returns
// Unknown for names this scaffold doesn't recognise yet.
func (r *Registry) ResolveName(name string) JSType {
	switch name {
	case "number":
		return Number
	case "string":
		return String
	case "boolean":
		return Boolean
	case "null":
		return Null
	case "undefined", "void":
		return Void
	case "*":
		return All
	}
	return Unknown
}

// ResolveExpr parses the textual form captured from a JSDoc type
// annotation ({e.g. "number|string", "!Foo", "?number"}) and returns a
// JSType. Unknown members collapse the whole expression to Unknown so
// the checker stays conservative, which matches how Java resolves an
// unparseable type expression.
func (r *Registry) ResolveExpr(expr string) JSType {
	expr = strings.TrimSpace(expr)
	if expr == "" || expr == "?" {
		return Unknown
	}
	// `?Foo` = Foo | null (nullable). Strip the leading `?` and union
	// with Null.
	nullable := false
	if strings.HasPrefix(expr, "?") {
		nullable = true
		expr = strings.TrimSpace(expr[1:])
	}
	parts := strings.Split(expr, "|")
	var members []JSType
	for _, p := range parts {
		p = strings.TrimSpace(p)
		p = strings.TrimPrefix(p, "!") // non-null is a no-op on primitives
		if p == "" {
			return Unknown
		}
		m := r.ResolveName(p)
		if isUnknown(m) {
			// Any unresolved member makes the whole expression Unknown
			// — conservative, avoids partial false positives.
			return Unknown
		}
		members = append(members, m)
	}
	if nullable {
		members = append(members, Null)
	}
	if len(members) == 1 {
		return members[0]
	}
	return NewUnion(members...)
}
