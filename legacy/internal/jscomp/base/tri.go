// Port of com/google/javascript/jscomp/base/Tri.java
// @ closure-compiler v20260415.
//
// Tri is a three-valued (triboolean) type used throughout the Closure
// Compiler to represent values that may be TRUE, FALSE, or UNKNOWN
// (indeterminate at compile time). It is the canonical tri-state type for
// static analysis: Outcome.isNullish(), getBooleanValue, and similar.

package base

// Tri is a triboolean. Mirrors com.google.javascript.jscomp.base.Tri.
type Tri int

const (
	// TriFalse mirrors Tri.FALSE.
	TriFalse Tri = iota
	// TriTrue mirrors Tri.TRUE.
	TriTrue
	// TriUnknown mirrors Tri.UNKNOWN.
	TriUnknown
)

// Not mirrors Tri.not(): TRUE→FALSE, FALSE→TRUE, UNKNOWN→UNKNOWN.
func (t Tri) Not() Tri {
	switch t {
	case TriTrue:
		return TriFalse
	case TriFalse:
		return TriTrue
	}
	return TriUnknown
}

// And mirrors Tri.and(): returns FALSE if either operand is FALSE, TRUE if
// both are TRUE, UNKNOWN otherwise.
func (t Tri) And(other Tri) Tri {
	if t == TriFalse || other == TriFalse {
		return TriFalse
	}
	if t == TriTrue && other == TriTrue {
		return TriTrue
	}
	return TriUnknown
}

// Or mirrors Tri.or(): returns TRUE if either operand is TRUE, FALSE if
// both are FALSE, UNKNOWN otherwise.
func (t Tri) Or(other Tri) Tri {
	if t == TriTrue || other == TriTrue {
		return TriTrue
	}
	if t == TriFalse && other == TriFalse {
		return TriFalse
	}
	return TriUnknown
}
