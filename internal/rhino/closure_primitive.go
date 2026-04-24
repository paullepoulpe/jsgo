// Port of com/google/javascript/rhino/ClosurePrimitive.java
// @ closure-compiler v20260415

package rhino

// ClosurePrimitive is an enum of primitive functions that the compiler
// recognises. These correspond to the @closurePrimitive tag in code.
//
// After typechecking all calls to one of these primitive types should have
// their FunctionType annotated with the corresponding enum member. This is
// intended to make identifying these calls more accurate than previous
// methods of finding primitives by qualified name.
type ClosurePrimitive int

const (
	// ClosurePrimitiveAssertsFail is a function that always throws an error.
	ClosurePrimitiveAssertsFail ClosurePrimitive = iota
	// ClosurePrimitiveAssertsMatchesReturn is a function that asserts its
	// first parameter matches the return type.
	ClosurePrimitiveAssertsMatchesReturn
	// ClosurePrimitiveAssertsTruthy is a function that asserts its first
	// parameter is truthy and returns the param.
	ClosurePrimitiveAssertsTruthy
)

// idToClosurePrimitive maps human-readable ids to enum members.
// The expected mapping of keys → values is that the enum member maps to a
// lowercase string with "_" replaced with ".".
var idToClosurePrimitive = map[string]ClosurePrimitive{
	"asserts.fail":          ClosurePrimitiveAssertsFail,
	"asserts.truthy":        ClosurePrimitiveAssertsTruthy,
	"asserts.matchesReturn": ClosurePrimitiveAssertsMatchesReturn,
}

// closurePrimitiveName maps each ClosurePrimitive to its Java enum name string.
var closurePrimitiveName = map[ClosurePrimitive]string{
	ClosurePrimitiveAssertsFail:          "ASSERTS_FAIL",
	ClosurePrimitiveAssertsMatchesReturn: "ASSERTS_MATCHES_RETURN",
	ClosurePrimitiveAssertsTruthy:        "ASSERTS_TRUTHY",
}

// Name mirrors the implicit Enum.name() method — returns the Java enum constant name.
func (p ClosurePrimitive) Name() string {
	if s, ok := closurePrimitiveName[p]; ok {
		return s
	}
	return "UNKNOWN"
}

// Go-only: language-adapter — ToString implements Java's Enum.toString() which ClosurePrimitive inherits;
// it is not declared explicitly in ClosurePrimitive.java but is implicitly available.
func (p ClosurePrimitive) ToString() string {
	return p.Name()
}

// FromStringId mirrors ClosurePrimitive.fromStringId(String) — returns the ClosurePrimitive
// corresponding to the given string id.
//
// Returns ok=false if id is empty; panics if id is non-empty but does not
// match any known enum member.
func FromStringId(id string) (ClosurePrimitive, bool) {
	if id == "" {
		return 0, false
	}
	p, ok := idToClosurePrimitive[id]
	if !ok {
		panic("String id " + id + " does not match any ClosurePrimitive")
	}
	return p, true
}

// Go-only: language-adapter — ClosurePrimitiveFromStringId is an alias for FromStringId using the
// package-prefixed naming convention; callers that predate the Java-parity rename use this.
func ClosurePrimitiveFromStringId(id string) (ClosurePrimitive, bool) {
	return FromStringId(id)
}
