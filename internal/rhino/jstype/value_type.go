// Port of com/google/javascript/rhino/jstype/ValueType.java
// @ closure-compiler v20260415
//
// ValueType is the abstract Java superclass of the primitive value
// types: null, void, number, boolean, string, and bigint.  It adds
// shared behaviour on top of JSType:
//
//   - appendTo delegates to getDisplayName (so the printed form of a
//     value type is always its display name, e.g. "number").
//   - hasDisplayName always returns true.
//   - resolveInternal panics: value types are eagerly resolved and
//     should never go through the lazy-resolution path.
//   - recursionUnsafeHashCode uses identity, because each value type
//     is a unique singleton within a JSTypeRegistry.
//
// In Go, JSType is an interface rather than an abstract class, so
// ValueType cannot be embedded into the existing primitive types
// without a broader refactor (tracked separately).  This file
// therefore ports ValueType as a standalone named struct that exposes
// the shared semantics via a marker method and a display-name
// accessor.  Callers that need to check "is this a value type?" may
// use IsValueType(t).

package jstype

// ValueType holds the shared identity of a JavaScript value type.
// Java's ValueType is abstract; the Go struct is meant to be embedded
// in concrete value-type structs (e.g. BigIntType) or used as a
// standalone marker.  Existing primitive singletons (Number, String,
// Boolean, Null, Void) pre-date this port and do not yet embed it —
// that migration is a separate task.
type ValueType struct {
	// displayName is the human-readable name returned by getDisplayName
	// in Java (e.g. "bigint", "number").
	displayName string
}

// newValueType constructs a ValueType with the given display name.
// Mirrors the Java constructor chain: ValueType(registry) →
// super(registry) → eagerlyResolveToSelf().
func newValueType(displayName string) ValueType {
	return ValueType{displayName: displayName}
}

// DisplayName returns the type's display name.
// Mirrors Java ValueType.getDisplayName() (abstract, overridden by
// each concrete subclass to return a fixed string).
func (v ValueType) DisplayName() string {
	return v.displayName
}

// HasDisplayName always returns true for value types.
// Mirrors Java ValueType.hasDisplayName().
func (v ValueType) HasDisplayName() bool {
	return true
}

// IsValueType is a marker that lets callers determine at runtime
// whether a JSType was constructed as a ValueType.
func (v ValueType) IsValueType() bool { return true }

// IsValueType reports whether t is (or embeds) a ValueType.
// Analogous to an instanceof check in Java.
func IsValueType(t JSType) bool {
	type valueTyper interface {
		IsValueType() bool
	}
	vt, ok := t.(valueTyper)
	return ok && vt.IsValueType()
}
