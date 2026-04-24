// Port of com/google/javascript/rhino/Outcome.java +
// com/google/javascript/jscomp/base/Tri.java
// @ closure-compiler v20260415
//
// Tri is included here because Outcome.isNullish() returns it and both
// are consumed together in flow-narrowing logic. The jscomp package
// keeps its own local Tri copy for historical reasons; this one is the
// canonical rhino-level definition.

package rhino

// Tri is a triboolean used by Outcome.IsNullish. Mirrors
// com.google.javascript.jscomp.base.Tri.
type Tri int

const (
	// TriFalse mirrors Tri.FALSE.
	TriFalse Tri = iota
	// TriTrue mirrors Tri.TRUE.
	TriTrue
	// TriUnknown mirrors Tri.UNKNOWN.
	TriUnknown
)

// Outcome represents truthiness and nullishness of expression values.
// Mirrors com.google.javascript.rhino.Outcome.
//
// FALSE contains both FALSE_NOT_NULL and NULLISH (all falsy values).
// It specifically differentiates between falsy and explicitly nullish values:
// e.g. 0 is falsy but not nullish.
type Outcome int

const (
	// OutcomeTrue represents truthy values, e.g. {}, true, 1.
	OutcomeTrue Outcome = iota
	// OutcomeFalse represents all falsy values, e.g. '', 0, false, null.
	OutcomeFalse
	// OutcomeFalseNotNull represents falsy, non-nullish values, e.g. 0, ''.
	OutcomeFalseNotNull
	// OutcomeNullish represents null and undefined.
	OutcomeNullish
)

// IsTruthy reports whether the Outcome is truthy. Mirrors Outcome.isTruthy().
func (o Outcome) IsTruthy() bool {
	return o == OutcomeTrue
}

// IsNullish reports whether the Outcome is nullish. Returns TriUnknown for
// OutcomeFalse since it covers both nullish and non-nullish falsy values.
// Mirrors Outcome.isNullish().
func (o Outcome) IsNullish() Tri {
	switch o {
	case OutcomeTrue:
		return TriFalse
	case OutcomeFalse:
		return TriUnknown
	case OutcomeFalseNotNull:
		return TriFalse
	case OutcomeNullish:
		return TriTrue
	}
	return TriUnknown
}

// Not returns the logical negation of the outcome. Mirrors Outcome.not().
// Note that NOT of any falsy variant returns OutcomeTrue.
func (o Outcome) Not() Outcome {
	if o == OutcomeTrue {
		return OutcomeFalse
	}
	return OutcomeTrue
}

// OutcomeForBoolean returns OutcomeTrue for true, OutcomeFalse for false.
// Mirrors Outcome.forBoolean(boolean).
func OutcomeForBoolean(val bool) Outcome {
	if val {
		return OutcomeTrue
	}
	return OutcomeFalse
}
