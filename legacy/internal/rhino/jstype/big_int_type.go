// Port of com/google/javascript/rhino/jstype/BigIntType.java
// @ closure-compiler v20260415
//
// BigIntType is the singleton primitive type for JavaScript's bigint
// value type.  It extends ValueType in Java, which itself extends
// JSType.  In Go, JSType is an interface; BigIntType implements it
// directly and embeds ValueType for the shared value-type behaviour.
//
// TODO: testForEquality returns Tri, which has not yet been ported;
// that method is omitted here and tracked separately.

package jstype

// bigIntType is the concrete runtime type for the bigint primitive.
// It is unexported so callers interact with the canonical BigInt
// singleton (package-level var) rather than constructing instances.
type bigIntType struct {
	ValueType
}

// TypeName returns "bigint", matching Java's BigIntType.getDisplayName()
// and the string produced by appendTo / toString in the Java hierarchy.
// Implements JSType.TypeName.
func (bigIntType) TypeName() string { return "bigint" }

// IsSubtypeOf reports whether bigint is a subtype of other.
// Bigint is a primitive singleton: it is a subtype of itself, AllType,
// and UnknownType only — mirroring the analogous rule in NumberType,
// StringType, and BooleanType.
// Implements JSType.IsSubtypeOf.
func (bigIntType) IsSubtypeOf(other JSType) bool {
	if _, ok := other.(bigIntType); ok {
		return true
	}
	return isAll(other) || isUnknown(other)
}

// IsBigIntValueType returns true, matching Java's
// BigIntType.isBigIntValueType().
func (bigIntType) IsBigIntValueType() bool { return true }

// MatchesNumberContext returns false.
// BigInt cannot be coerced to a Number, so it does not match a numeric
// context.  Matches Java BigIntType.matchesNumberContext().
func (bigIntType) MatchesNumberContext() bool { return false }

// MatchesStringContext returns true.
// BigInt can be coerced to a string.
// Matches Java BigIntType.matchesStringContext().
func (bigIntType) MatchesStringContext() bool { return true }

// MatchesObjectContext returns true.
// BigInt can be auto-boxed to BigInt Object.
// Matches Java BigIntType.matchesObjectContext().
func (bigIntType) MatchesObjectContext() bool { return true }

// GetPossibleToBooleanOutcomes returns BooleanLiteralSetBoth because a
// bigint value can be truthy (non-zero) or falsy (zero).
// Matches Java BigIntType.getPossibleToBooleanOutcomes().
func (bigIntType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return BooleanLiteralSetBoth
}

// BigInt is the canonical interned instance for the bigint primitive type.
// Java's JSTypeRegistry returns the same BigIntType pointer every time;
// we mirror that with a package-level var, consistent with Number, String,
// Boolean, Null, and Void.
var BigInt JSType = bigIntType{newValueType("bigint")}
