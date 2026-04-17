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

func isAll(t JSType) bool     { _, ok := t.(allType); return ok }
func isUnknown(t JSType) bool { _, ok := t.(unknownType); return ok }

// Canonical interned instances. Java JSTypeRegistry returns the same
// instance for each primitive; we mirror that with package-level vars.
var (
	All     JSType = allType{}
	Unknown JSType = unknownType{}
	Number  JSType = numberType{}
	String  JSType = stringType{}
)
