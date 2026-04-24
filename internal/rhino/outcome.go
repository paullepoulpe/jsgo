// Port of com/google/javascript/rhino/Outcome.java
// @ closure-compiler v20260415.
//
// Outcome represents truthiness and nullishness of expression values.
// IsNullish() returns a Tri (triboolean) defined in
// internal/jscomp/base (com/google/javascript/jscomp/base/Tri.java).

package rhino

import "github.com/paullepoulpe/jsgo/internal/jscomp/base"

// Tri is a package-level alias so that existing code in this package
// can reference Tri without the full import path. The canonical definition
// lives in internal/jscomp/base.
// Go-only: language-adapter — Tri type alias; Java uses com.google.javascript.jscomp.base.Tri directly.
type Tri = base.Tri

const (
	// TriFalse mirrors Tri.FALSE.
	TriFalse = base.TriFalse
	// TriTrue mirrors Tri.TRUE.
	TriTrue = base.TriTrue
	// TriUnknown mirrors Tri.UNKNOWN.
	TriUnknown = base.TriUnknown
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
