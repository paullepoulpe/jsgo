// Port of com/google/javascript/jscomp/type/Outcome.java @ closure-compiler v20260415

/*
 * Copyright 2020 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/base"
)

// Outcome is an enum for representing boolean outcomes. Currently they just
// represent typical booleans but this enum will be expanded to specifically
// differentiate between falsy and explicitly nullish values. e.g. 0 is falsy
// but not nullish.
type Outcome int32

const (
	// Outcome_TRUE represents truthy values. For example: {}, true, 1, etc.
	Outcome_TRUE Outcome = iota
	// Outcome_FALSE represents falsy values. For examples: '', 0, false,
	// null, etc.
	Outcome_FALSE
)

var outcome_names = [...]string{"TRUE", "FALSE"}

// Outcome_values is Outcome.values().
func Outcome_values() []Outcome { return []Outcome{Outcome_TRUE, Outcome_FALSE} }

// Name is Enum.name().
func (o Outcome) Name() string { return outcome_names[o] }

// String is Enum.toString().
func (o Outcome) String() string { return o.Name() }

// Ordinal is Enum.ordinal().
func (o Outcome) Ordinal() int { return int(o) }

// IsTruthy determines whether an Outcome enum value is truthy.
//
// Java: Outcome.isTruthy()
func (o Outcome) IsTruthy() bool {
	switch o {
	case Outcome_TRUE:
		return true
	default: // FALSE
		return false
	}
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
	default: // FALSE
		return base.Tri_UNKNOWN
	}
}

// Not gets the not of this.
//
// Java: Outcome.not()
func (o Outcome) Not() Outcome {
	switch o {
	case Outcome_TRUE:
		return Outcome_FALSE
	default: // FALSE
		return Outcome_TRUE
	}
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
