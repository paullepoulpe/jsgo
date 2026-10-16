// Port of com/google/javascript/rhino/ClosurePrimitive.java @ closure-compiler v20260415

package rhino

import "github.com/paullepoulpe/jsgo/internal/jutil"

// ClosurePrimitive enumerates the primitive functions that the compiler
// recognizes.
//
// These correspond to the @closurePrimitive tag in code; in order to parse
// new primitives, add any entry to the list in
// parsing/ParserConfig.properties, then map it to an enum member in the
// idToEnum map.
//
// After typechecking is done, all calls to one of these primitive types
// should have their FunctionType annotated with the corresponding enum
// member. This is intended to make identifying these calls more accurate
// than previous methods of finding primitives by qualified name.
type ClosurePrimitive int32

const (
	ClosurePrimitive_ASSERTS_FAIL           ClosurePrimitive = iota // A function that always throws an error
	ClosurePrimitive_ASSERTS_MATCHES_RETURN                         // A function that asserts its first parameter matches the return type
	ClosurePrimitive_ASSERTS_TRUTHY                                 // A function that asserts its first parameter is truthy and returns the param
)

var closurePrimitive_names = [...]string{"ASSERTS_FAIL", "ASSERTS_MATCHES_RETURN", "ASSERTS_TRUTHY"}

// ClosurePrimitive_values is ClosurePrimitive.values().
func ClosurePrimitive_values() []ClosurePrimitive {
	return []ClosurePrimitive{ClosurePrimitive_ASSERTS_FAIL, ClosurePrimitive_ASSERTS_MATCHES_RETURN, ClosurePrimitive_ASSERTS_TRUTHY}
}

// String is Enum.name().
func (p ClosurePrimitive) String() string { return closurePrimitive_names[p] }

// ClosurePrimitive_ID_TO_ENUM maps human-readable ids to enum members.
//
// The expected mapping (although not enforced) of keys -> values is that
// the enum member maps to a lowercase string with "_" replaced with "."
//
// Java: ClosurePrimitive.ID_TO_ENUM
var ClosurePrimitive_ID_TO_ENUM = func() *jutil.LinkedHashMap[string, ClosurePrimitive] {
	m := jutil.NewLinkedHashMap[string, ClosurePrimitive]()
	m.Put("asserts.fail", ClosurePrimitive_ASSERTS_FAIL)
	m.Put("asserts.truthy", ClosurePrimitive_ASSERTS_TRUTHY)
	m.Put("asserts.matchesReturn", ClosurePrimitive_ASSERTS_MATCHES_RETURN)
	return m
}()

// ClosurePrimitive_FromStringId returns the ClosurePrimitive corresponding to
// the given string id.
//
// This is to make reading @closurePrimitive easier in code. Using
// Enum.valueOf to parse closure primitive identifiers from JSDoc directly
// would require code like `closurePrimitive {ASSERTS_FAIL}`; instead we
// separate the string ids from the enum names.
//
// Returns null if the argument is null, otherwise the corresponding enum
// member. Throws IllegalArgumentException if the id is non-null but does not
// match an enum member.
//
// Java: ClosurePrimitive.fromStringId(String)
func ClosurePrimitive_FromStringId(id *string) *ClosurePrimitive {
	if id == nil {
		return nil
	}

	primitive, ok := ClosurePrimitive_ID_TO_ENUM.Lookup(*id)
	if !ok {
		panic(jutil.NewIllegalArgumentExceptionMsg(
			"String id " + *id + " does not match any ClosurePrimitive"))
	}

	return &primitive
}
