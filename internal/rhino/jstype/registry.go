// Port of com/google/javascript/rhino/jstype/JSTypeRegistry.java @ closure-compiler v20260415
//
// The real registry manages interning of all JSType values, named-type
// resolution, template instantiation, and much more. Tier-0 only needs a
// name-to-primitive lookup, which we provide here.

package jstype

type Registry struct{}

func NewRegistry() *Registry { return &Registry{} }

// ResolveName mirrors JSTypeRegistry.getType(StaticScope, String). Returns
// Unknown for names this scaffold doesn't recognise yet (same fallback
// shape as Java's "unknown type" when a declaration is unresolved).
func (r *Registry) ResolveName(name string) JSType {
	switch name {
	case "number":
		return Number
	case "string":
		return String
	case "*":
		return All
	}
	return Unknown
}
