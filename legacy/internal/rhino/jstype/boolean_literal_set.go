// Port of com/google/javascript/rhino/jstype/BooleanLiteralSet.java
// @ closure-compiler v20260415
//
// BooleanLiteralSet represents a subset of {true, false}.  There are
// exactly four possible values:
//
//   BooleanLiteralSetEmpty — {} (the empty set)
//   BooleanLiteralSetTrue  — {true}
//   BooleanLiteralSetFalse — {false}
//   BooleanLiteralSetBoth  — {true, false}
//
// The Java original is an enum whose ordinal encodes the set as a
// two-bit field (bit 0 = true present, bit 1 = false present).
// Intersection and union are therefore bitwise AND/OR on the ordinal,
// which we replicate exactly here.

package jstype

// BooleanLiteralSet is the Go equivalent of Java's BooleanLiteralSet enum.
// The underlying int uses the same two-bit encoding as the Java ordinal so
// that Intersection and Union can be implemented with bitwise operators.
type BooleanLiteralSet int

const (
	// BooleanLiteralSetEmpty is the empty set {}.
	BooleanLiteralSetEmpty BooleanLiteralSet = 0
	// BooleanLiteralSetTrue is {true}.
	BooleanLiteralSetTrue BooleanLiteralSet = 1
	// BooleanLiteralSetFalse is {false}.
	BooleanLiteralSetFalse BooleanLiteralSet = 2
	// BooleanLiteralSetBoth is {true, false}.
	BooleanLiteralSetBoth BooleanLiteralSet = 3
)

// Intersection returns the set of boolean literals present in both b and other.
// Mirrors Java BooleanLiteralSet.intersection.
func (b BooleanLiteralSet) Intersection(other BooleanLiteralSet) BooleanLiteralSet {
	return b & other
}

// Union returns the set of boolean literals present in b or other.
// Mirrors Java BooleanLiteralSet.union.
func (b BooleanLiteralSet) Union(other BooleanLiteralSet) BooleanLiteralSet {
	return b | other
}

// Contains reports whether literalValue is an element of b.
// Mirrors Java BooleanLiteralSet.contains.
func (b BooleanLiteralSet) Contains(literalValue bool) bool {
	switch b {
	case BooleanLiteralSetEmpty:
		return false
	case BooleanLiteralSetTrue:
		return literalValue
	case BooleanLiteralSetFalse:
		return !literalValue
	case BooleanLiteralSetBoth:
		return true
	default:
		panic("BooleanLiteralSet: invalid ordinal")
	}
}

// BooleanLiteralSetGet returns the singleton set {literalValue}.
// Mirrors Java BooleanLiteralSet.get.
func BooleanLiteralSetGet(literalValue bool) BooleanLiteralSet {
	if literalValue {
		return BooleanLiteralSetTrue
	}
	return BooleanLiteralSetFalse
}
