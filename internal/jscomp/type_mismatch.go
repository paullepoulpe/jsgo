// Port of com/google/javascript/jscomp/TypeMismatch.java @ closure-compiler v20261021

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
	found    rhino.JSType
	required rhino.JSType
	location *rhino.Node
}

// NewTypeMismatch is the record's canonical (compact) constructor.
//
// Java: TypeMismatch(JSType, JSType, Node)
func NewTypeMismatch(found rhino.JSType, required rhino.JSType, location *rhino.Node) *TypeMismatch {
	jutil.CheckNotNullMsg(found, "found")
	jutil.CheckNotNullMsg(required, "required")
	jutil.CheckNotNullMsg(location, "location")
	return &TypeMismatch{found: found, required: required, location: location}
}

// Found is the record component accessor.
//
// Java: TypeMismatch.found()
func (m *TypeMismatch) Found() rhino.JSType { return m.found }

// Required is the record component accessor.
//
// Java: TypeMismatch.required()
func (m *TypeMismatch) Required() rhino.JSType { return m.required }

// Location is the record component accessor.
//
// Java: TypeMismatch.location()
func (m *TypeMismatch) Location() *rhino.Node { return m.location }

// Java: TypeMismatch.getFound()
func (m *TypeMismatch) GetFound() rhino.JSType {
	return m.Found()
}

// Java: TypeMismatch.getRequired()
func (m *TypeMismatch) GetRequired() rhino.JSType {
	return m.Required()
}

// Java: TypeMismatch.getLocation()
func (m *TypeMismatch) GetLocation() *rhino.Node {
	return m.Location()
}

// Java: TypeMismatch.create(JSType, JSType, Node)
func TypeMismatch_Create(found rhino.JSType, required rhino.JSType, location *rhino.Node) *TypeMismatch {
	return NewTypeMismatch(found, required, location)
}

// Java: TypeMismatch.createForTesting(JSType, JSType)
func TypeMismatch_CreateForTesting(found rhino.JSType, required rhino.JSType) *TypeMismatch {
	return TypeMismatch_Create(found, required, typeMismatch_TEST_LOCATION)
}

// Java: TypeMismatch.TEST_LOCATION
var typeMismatch_TEST_LOCATION = rhino.IR_Empty()

// Equals is the record's implicit equals: each component compared with
// Objects.equals (JSType.equals for the types; identity for the Node, which
// does not override equals).
//
// Java: TypeMismatch.equals(Object)
func (m *TypeMismatch) Equals(o any) bool {
	that, ok := o.(*TypeMismatch)
	if !ok || that == nil {
		return false
	}
	if m == that {
		return true
	}
	return m.found.Equals(that.found) &&
		m.required.Equals(that.required) &&
		m.location == that.location
}

// HashCode is the record's implicit hashCode: 31 * h + hash(component) over
// the components in order, starting from 0.
//
// Java: TypeMismatch.hashCode()
func (m *TypeMismatch) HashCode() int32 {
	var h int32
	h = 31*h + m.found.HashCode()
	h = 31*h + m.required.HashCode()
	h = 31*h + jutil.ObjectHashCode(m.location)
	return h
}

// String is the record's implicit toString.
//
// Java: TypeMismatch.toString()
func (m *TypeMismatch) String() string {
	return "TypeMismatch[found=" + m.found.String() +
		", required=" + m.required.String() +
		", location=" + m.location.String() + "]"
}

// TypeMismatch_Accumulator collects a set of related mismatches.
type TypeMismatch_Accumulator struct {
	mismatches []*TypeMismatch
}

// Java: TypeMismatch.Accumulator()
func NewTypeMismatch_Accumulator() *TypeMismatch_Accumulator {
	return &TypeMismatch_Accumulator{mismatches: []*TypeMismatch{}}
}

// Java: TypeMismatch.Accumulator.registerMismatch(Node, JSType, JSType)
func (a *TypeMismatch_Accumulator) RegisterMismatch(location *rhino.Node, found rhino.JSType, required rhino.JSType) {
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
		a.mismatches = append(a.mismatches, TypeMismatch_Create(found, required, location))
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
func (a *TypeMismatch_Accumulator) GetMismatches() []*TypeMismatch {
	return slices.Clone(a.mismatches)
}

// Java: TypeMismatch.Accumulator.registerIfMismatch(Node, JSType, JSType)
func (a *TypeMismatch_Accumulator) RegisterIfMismatch(location *rhino.Node, found rhino.JSType, required rhino.JSType) {
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
