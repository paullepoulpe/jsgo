// Port of com/google/javascript/rhino/jstype/NoType.java
// @ closure-compiler v20260415
//
// NoType is the full bottom type — a subtype of every value or object.
// It surfaces in static analysis when flow-narrowing produces an
// impossible value (e.g. `if (x)` inside a branch where x was null).
// Since it has no possible runtime type, any call `f(x)` is legal no
// matter what type `f`'s parameter requires — this is the formal
// definition of ⊥ in a type lattice.
//
// In the Java hierarchy NoType extends NoObjectType. It overrides
// isNoObjectType() → false and isNoType() → true, and additionally
// reports isNullable() = true and isVoidable() = true (since the empty
// set of values trivially contains null and undefined). The
// getPossibleToBooleanOutcomes() is BOTH (any boolean outcome is
// possible for an unreachable value).
//
// See: http://en.wikipedia.org/wiki/Bottom_type

package jstype

// NoType is the bottom type, the subtype of every value or object.
// Mirrors com.google.javascript.rhino.jstype.NoType.
type NoType struct{}

// None is the canonical interned instance. Java's JSTypeRegistry
// returns the same NoType for the lifetime of a compilation; we mirror
// that with a package-level variable. Named "None" to follow the
// Java appendTo output ("None" in non-annotation contexts).
var None JSType = &NoType{}

// TypeName returns "None" — the internal toString representation.
// Java's appendTo emits "?" for annotation contexts and "None"
// otherwise; we expose the non-annotation form for diagnostics.
func (n *NoType) TypeName() string { return "None" }

// IsNoType mirrors NoType.isNoType() → true.
func (n *NoType) IsNoType() bool { return true }

// IsNoObjectType mirrors NoType.isNoObjectType() → false. NoType
// extends NoObjectType in Java but overrides this to false so that
// `instanceof NoObjectType` checks don't confuse the two.
func (n *NoType) IsNoObjectType() bool { return false }

// IsNullable mirrors NoType.isNullable() → true. The empty set of
// values trivially satisfies the nullable constraint.
func (n *NoType) IsNullable() bool { return true }

// IsVoidable mirrors NoType.isVoidable() → true. The empty set of
// values trivially satisfies the voidable constraint.
func (n *NoType) IsVoidable() bool { return true }

// IsSubtypeOf follows the bottom-type rule: ⊥ is a subtype of every
// type. This mirrors how Java's NoType participates in the isSubtype
// lattice — it is below everything, including NoObjectType.
func (n *NoType) IsSubtypeOf(other JSType) bool { return true }
