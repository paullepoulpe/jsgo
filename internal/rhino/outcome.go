// Port of com/google/javascript/rhino/Outcome.java @ closure-compiler v20261015

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Annie Wang
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

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
