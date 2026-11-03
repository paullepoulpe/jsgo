// Port of com/google/javascript/jscomp/base/Tri.java @ closure-compiler v20261028

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

package base

import "github.com/paullepoulpe/jsgo/internal/jutil"

// Tri is an enum for ternary logic.
//
// The TRUE and FALSE values are equivalent to typical booleans, and the
// UNKNOWN value plays the role of a placeholder, which can be either TRUE or
// FALSE.
//
// A ternary value expression evaluates to TRUE or FALSE only if all
// replacements of UNKNOWN in this expression yield the same result.
// Therefore, the ternary logic coincides with typical Boolean logic if the
// UNKNOWN value is not present in an expression.
//
// See http://en.wikipedia.org/wiki/Ternary_logic
type Tri int32

const (
	Tri_FALSE Tri = iota
	Tri_UNKNOWN
	Tri_TRUE
)

var tri_names = [...]string{"FALSE", "UNKNOWN", "TRUE"}

// Tri_values is Tri.values().
func Tri_values() []Tri { return []Tri{Tri_FALSE, Tri_UNKNOWN, Tri_TRUE} }

// Name is Enum.name().
func (t Tri) Name() string { return tri_names[t] }

// Or is ternary or.
//
// Java: Tri.or(Tri)
func (t Tri) Or(x Tri) Tri {
	if t == Tri_TRUE || x == Tri_TRUE {
		return Tri_TRUE
	}
	if t == Tri_UNKNOWN || x == Tri_UNKNOWN {
		return Tri_UNKNOWN
	}
	return Tri_FALSE
}

// And is ternary and.
//
// Java: Tri.and(Tri)
func (t Tri) And(x Tri) Tri {
	if t == Tri_FALSE || x == Tri_FALSE {
		return Tri_FALSE
	}
	if t == Tri_UNKNOWN || x == Tri_UNKNOWN {
		return Tri_UNKNOWN
	}
	return Tri_TRUE
}

// Xor is ternary xor.
//
// Java: Tri.xor(Tri)
func (t Tri) Xor(x Tri) Tri {
	if t == Tri_UNKNOWN || x == Tri_UNKNOWN {
		return Tri_UNKNOWN
	}
	if t == x {
		return Tri_FALSE
	}
	return Tri_TRUE
}

// Not is ternary negation.
//
// Java: Tri.not()
func (t Tri) Not() Tri {
	switch t {
	case Tri_FALSE:
		return Tri_TRUE
	case Tri_TRUE:
		return Tri_FALSE
	default:
		return Tri_UNKNOWN
	}
}

// ToBoolean returns the boolean value, or x for UNKNOWN.
//
// Java: Tri.toBoolean(boolean)
func (t Tri) ToBoolean(x bool) bool {
	switch t {
	case Tri_FALSE:
		return false
	case Tri_UNKNOWN:
		return x
	case Tri_TRUE:
		return true
	}
	panic(jutil.NewIllegalStateException())
}

// Tri_ForBoolean returns TRUE or FALSE.
//
// Java: Tri.forBoolean(boolean)
func Tri_ForBoolean(x bool) Tri {
	if x {
		return Tri_TRUE
	}
	return Tri_FALSE
}

// String is Tri.toString(): the lower-case name.
//
// Java: Tri.toString()
func (t Tri) String() string {
	return jutil.StringToLowerCaseASCII(t.Name())
}
