// Port of com/google/javascript/jscomp/base/Tri.java @ closure-compiler v20261021

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

// tri_value is the `value` field of each constant.
var tri_value = [...]int32{-1, 0, 1}

// Tri_values is Tri.values().
func Tri_values() []Tri { return []Tri{Tri_FALSE, Tri_UNKNOWN, Tri_TRUE} }

// Name is Enum.name().
func (t Tri) Name() string { return tri_names[t] }

// Ordinal is Enum.ordinal().
func (t Tri) Ordinal() int { return int(t) }

// Or is ternary or.
//
// Java: Tri.or(Tri)
func (t Tri) Or(x Tri) Tri {
	if tri_value[t] > tri_value[x] {
		return t
	}
	return x
}

// And is ternary and.
//
// Java: Tri.and(Tri)
func (t Tri) And(x Tri) Tri {
	if tri_value[t] < tri_value[x] {
		return t
	}
	return x
}

// Xor is ternary xor.
//
// Java: Tri.xor(Tri)
func (t Tri) Xor(x Tri) Tri {
	return Tri_ForInt(-tri_value[t] * tri_value[x])
}

// Not is ternary negation.
//
// Java: Tri.not()
func (t Tri) Not() Tri {
	return Tri_ForInt(-tri_value[t])
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

// Java: Tri.forInt(int)
func Tri_ForInt(x int32) Tri {
	return Tri_VALUES[x+1]
}

// Java: Tri.VALUES
var Tri_VALUES = Tri_values()
