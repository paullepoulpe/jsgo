// Port of com/google/javascript/rhino/jstype/NoObjectType.java
// @ closure-compiler v20260415
//
// NoObjectType is the bottom Object type — the subclass of all objects.
// It surfaces in static analysis when flow-narrowing produces an
// impossible object (e.g. `x instanceof Array` inside a branch where x
// was declared non-object). Every property access and method call on it
// is legal by definition, which is why DefineProperty always returns true.
// In the Java hierarchy NoObjectType extends FunctionType; in Go we embed
// FunctionType so that callers that want a *FunctionType receive nil via
// ToMaybeFunctionType (mirrors Java's override returning null).
//
// See: http://en.wikipedia.org/wiki/Bottom_type

package jstype

// NoObjectType is the bottom Object type, representing the subclass of
// all objects. Mirrors com.google.javascript.rhino.jstype.NoObjectType.
type NoObjectType struct{}

// NoObject is the canonical interned instance. Java's JSTypeRegistry
// returns the same NoObjectType for the lifetime of a compilation; we
// mirror that with a package-level variable.
var NoObject JSType = &NoObjectType{}

// TypeName returns "NoObject" (the internal toString representation).
// Java's appendTo emits "?" for annotation contexts and "NoObject"
// otherwise; we expose the non-annotation form for diagnostics.
func (n *NoObjectType) TypeName() string { return "NoObject" }

// IsNoObjectType mirrors NoObjectType.isNoObjectType() → true.
func (n *NoObjectType) IsNoObjectType() bool { return true }

// IsNoType mirrors the default JSType.isNoType() → false; overridden
// to false here since NoObjectType is distinct from the full bottom.
func (n *NoObjectType) IsNoType() bool { return false }

// MatchesNumberContext mirrors NoObjectType.matchesNumberContext() → true.
// Bottom objects are compatible with every context.
func (n *NoObjectType) MatchesNumberContext() bool { return true }

// MatchesObjectContext mirrors NoObjectType.matchesObjectContext() → true.
func (n *NoObjectType) MatchesObjectContext() bool { return true }

// MatchesStringContext mirrors NoObjectType.matchesStringContext() → true.
func (n *NoObjectType) MatchesStringContext() bool { return true }

// MatchesSymbolContext mirrors NoObjectType.matchesSymbolContext() → true.
func (n *NoObjectType) MatchesSymbolContext() bool { return true }

// IsSubtypeOf follows the bottom-type rule: ⊥ is a subtype of every
// type. This mirrors how Java's NoObjectType participates in the
// isSubtype lattice — it is below everything except itself.
func (n *NoObjectType) IsSubtypeOf(other JSType) bool { return true }
