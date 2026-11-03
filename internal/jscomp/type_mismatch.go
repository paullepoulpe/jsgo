// Port of com/google/javascript/jscomp/TypeMismatch.java @ closure-compiler v20261028

/*
 * Copyright 2017 The Closure Compiler Authors.
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
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// TypeMismatch signals that the first type and the second type have been
// used interchangeably.
//
// Type-based optimizations should take this into account so that they
// don't wreck code with type warnings.
//
// found is the RHS type; the type of the assignment target. required is the
// LHS type; the type being assigned. location is the location of the
// assignment.
//
// Java: record TypeMismatch(JSType found, JSType required, Node location)
type TypeMismatch struct {
	Found    rhino.JSType
	Required rhino.JSType
	Location rhino.Node
}

// NewTypeMismatch constructs a TypeMismatch after checking that none of its
// components are null.
//
// Java: TypeMismatch(JSType, JSType, Node) / TypeMismatch.create(JSType, JSType, Node)
func NewTypeMismatch(found rhino.JSType, required rhino.JSType, location rhino.Node) TypeMismatch {
	jutil.CheckNotNullMsg(found, "found")
	jutil.CheckNotNullMsg(required, "required")
	rhino.CheckNotNullMsg(location, "location")
	return TypeMismatch{Found: found, Required: required, Location: location}
}

// Java: TypeMismatch.createForTesting(JSType, JSType)
func NewTypeMismatchForTesting(found rhino.JSType, required rhino.JSType) TypeMismatch {
	return NewTypeMismatch(found, required, rhino.IR_Empty())
}

// String is the record's implicit toString.
//
// Java: TypeMismatch.toString()
func (m TypeMismatch) String() string {
	return "TypeMismatch[found=" + m.Found.String() +
		", required=" + m.Required.String() +
		", location=" + m.Location.String() + "]"
}

// TypeMismatch_Accumulator collects a set of related mismatches.
type TypeMismatch_Accumulator struct {
	mismatches []TypeMismatch
}

// Java: TypeMismatch.Accumulator()
func NewTypeMismatch_Accumulator() *TypeMismatch_Accumulator {
	return &TypeMismatch_Accumulator{}
}

// Java: TypeMismatch.Accumulator.registerMismatch(Node, JSType, JSType)
func (a *TypeMismatch_Accumulator) RegisterMismatch(location rhino.Node, found rhino.JSType, required rhino.JSType) {
	// Don't register a mismatch for differences in null or undefined or if the
	// code didn't downcast.
	found = TypeMismatch_Accumulator_RemoveNullUndefinedAndTemplates(found)
	required = TypeMismatch_Accumulator_RemoveNullUndefinedAndTemplates(required)
	if found.IsSubtypeOf(required) || required.IsSubtypeOf(found) {
		return
	}
	if required.IsUnionType() {
		for _, requiredAltType := range required.ToMaybeUnionType().GetAlternates() {
			if requiredAltType.IsSubtypeOf(found) {
				return
			}
		}
	}

	if TypeMismatch_Accumulator_BothAreNotTemplateTypes(found, required) {
		a.mismatches = append(a.mismatches, NewTypeMismatch(found, required, location))
	}

	if found.IsFunctionType() && required.IsFunctionType() {
		fnTypeA := found.ToMaybeFunctionType()
		fnTypeB := required.ToMaybeFunctionType()
		paramsA := fnTypeA.GetParameters()
		paramsB := fnTypeB.GetParameters()
		for i := 0; i < len(paramsA) && i < len(paramsB); i++ {
			a.RegisterIfMismatch(location, paramsA[i].GetJSType(), paramsB[i].GetJSType())
		}
		a.RegisterIfMismatch(location, fnTypeA.GetReturnType(), fnTypeB.GetReturnType())
	}
}

// Java: TypeMismatch.Accumulator.getMismatches()
func (a *TypeMismatch_Accumulator) GetMismatches() []TypeMismatch {
	if len(a.mismatches) == 0 {
		return []TypeMismatch{}
	}
	return slices.Clone(a.mismatches)
}

// Java: TypeMismatch.Accumulator.registerIfMismatch(Node, JSType, JSType)
func (a *TypeMismatch_Accumulator) RegisterIfMismatch(location rhino.Node, found rhino.JSType, required rhino.JSType) {
	if found != nil && required != nil && !found.IsSubtypeOf(required) {
		a.RegisterMismatch(location, found, required)
	}
}

// TypeMismatch_Accumulator_BothAreNotTemplateTypes: a type variable is
// considered unknown, so other types can be used as type variables, and vice
// versa, without warning. Otherwise, to get good disambiguation, we would
// have to add casts all over the code base.
//
// TODO(dimvar): this can be made safe in the distant future where we have
// bounded generics *and* we have switched all the unsafe uses of type
// variables in the code base to use bounded generics.
//
// Java: TypeMismatch.Accumulator.bothAreNotTemplateTypes(JSType, JSType)
func TypeMismatch_Accumulator_BothAreNotTemplateTypes(found rhino.JSType, required rhino.JSType) bool {
	return !found.IsTemplateType() && !required.IsTemplateType()
}

// Java: TypeMismatch.Accumulator.removeNullUndefinedAndTemplates(JSType)
func TypeMismatch_Accumulator_RemoveNullUndefinedAndTemplates(t rhino.JSType) rhino.JSType {
	result := t.RestrictByNotNullOrUndefined()
	obj := result.ToMaybeObjectType()
	if obj != nil && obj.IsTemplatizedType() {
		// We don't care about the specific specalization involved in the mismatch because all
		// specializations share the same JS code.
		return obj.ToMaybeTemplatizedType().GetRawType()
	}
	return result
}
