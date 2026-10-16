// Port of com/google/javascript/rhino/Outcome.java @ closure-compiler v20260415

package rhino

import "github.com/paullepoulpe/jsgo/internal/base"

// Outcome is an enum for representing truthiness and nullishness of
// outcomes. The values are TRUE, FALSE, FALSE_NOT_NULL, and NULLISH. FALSE
// represents all falsy values so it contains both FALSE_NOT_NULL and
// NULLISH. It specifically differentiates between falsy and explicitly
// nullish values. e.g. 0 is falsy but not nullish.
type Outcome int32

const (
	// Represents truthy values. For example: {}, true, 1, etc.
	Outcome_TRUE Outcome = iota
	// Represents falsy values. For examples: '', 0, false, null, etc.
	Outcome_FALSE
	Outcome_FALSE_NOT_NULL
	Outcome_NULLISH
)

var outcome_names = [...]string{"TRUE", "FALSE", "FALSE_NOT_NULL", "NULLISH"}

// Outcome_values is Outcome.values().
func Outcome_values() []Outcome {
	return []Outcome{Outcome_TRUE, Outcome_FALSE, Outcome_FALSE_NOT_NULL, Outcome_NULLISH}
}

// String is Enum.name().
func (o Outcome) String() string { return outcome_names[o] }

// IsTruthy determines whether an Outcome enum value is truthy.
//
// Java: Outcome.isTruthy()
func (o Outcome) IsTruthy() bool {
	switch o {
	case Outcome_TRUE:
		return true
	}
	return false
}

// IsNullish determines whether an Outcome enum value is nullish. Using Tri
// instead of a boolean because 0 is Outcome.FALSE but not nullish so
// sometimes it is unclear.
//
// Java: Outcome.isNullish()
func (o Outcome) IsNullish() base.Tri {
	switch o {
	case Outcome_TRUE:
		return base.Tri_FALSE
	case Outcome_FALSE:
		return base.Tri_UNKNOWN
	case Outcome_FALSE_NOT_NULL:
		return base.Tri_FALSE
	}
	return base.Tri_TRUE // NULLISH
}

// Not gets the `not` of `this`.
//
// Java: Outcome.not()
func (o Outcome) Not() Outcome {
	switch o {
	case Outcome_TRUE:
		return Outcome_FALSE
	}
	return Outcome_TRUE // FALSE, FALSE_NOT_NULL, NULLISH
}

// Outcome_ForBoolean gets the Outcome for the given boolean.
//
// Java: Outcome.forBoolean(boolean)
func Outcome_ForBoolean(val bool) Outcome {
	if val {
		return Outcome_TRUE
	}
	return Outcome_FALSE
}
