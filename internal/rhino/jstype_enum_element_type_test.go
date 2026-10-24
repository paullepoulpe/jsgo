// Port of com/google/javascript/rhino/jstype/EnumElementTypeTest.java @ closure-compiler v20261021

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
 *   Nick Santos
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino_test

import (
	"testing"

	. "github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/rhinotesting"
)

type enumElementTypeTest struct {
	*rhinotesting.BaseJSTypeTestCase
}

func newEnumElementTypeTest(t *testing.T) *enumElementTypeTest {
	s := &enumElementTypeTest{BaseJSTypeTestCase: rhinotesting.NewBaseJSTypeTestCase(t)}
	return s
}

func TestEnumElementType_SubtypeRelation(t *testing.T) {
	s := newEnumElementTypeTest(t)
	typeA := s.CreateEnumType("typeA", s.NUMBER_TYPE).GetElementsType()
	typeB := s.CreateEnumType("typeB", s.NUMBER_TYPE).GetElementsType()

	s.AssertThat(typeA.IsSubtype(typeB)).IsFalse()
	s.AssertThat(typeB.IsSubtype(typeA)).IsFalse()

	s.AssertThat(s.NUMBER_TYPE.IsSubtypeOf(typeB)).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.IsSubtypeOf(typeA)).IsFalse()

	s.AssertThat(typeA.IsSubtype(s.NUMBER_TYPE)).IsTrue()
	s.AssertThat(typeB.IsSubtype(s.NUMBER_TYPE)).IsTrue()
}

func TestEnumElementType_GetGreatestSubtype(t *testing.T) {
	s := newEnumElementTypeTest(t)
	typeA :=
		s.CreateEnumType("typeA", s.CreateUnionType(s.NUMBER_TYPE, s.STRING_TYPE)).GetElementsType()

	stringsOfA := typeA.GetGreatestSubtype(s.STRING_TYPE)
	s.AssertThat(stringsOfA.IsEmptyType()).IsFalse()
	s.AssertThat(stringsOfA.String()).IsEqualTo("typeA<string>")
	s.AssertThat(stringsOfA.IsSubtypeOf(typeA)).IsTrue()

	numbersOfA := s.NUMBER_TYPE.GetGreatestSubtype(typeA)
	s.AssertThat(numbersOfA.IsEmptyType()).IsFalse()
	s.AssertThat(numbersOfA.String()).IsEqualTo("typeA<number>")
	s.AssertThat(numbersOfA.IsSubtypeOf(typeA)).IsTrue()
}

func TestEnumElementType_GetGreatestSubtype_twoEnumElementTypes(t *testing.T) {
	s := newEnumElementTypeTest(t)
	typeA := s.CreateEnumType("typeA", s.NUMBER_TYPE).GetElementsType()
	typeB := s.CreateEnumType("typeB", s.NUMBER_TYPE).GetElementsType()

	greatestSubtype := EnumElementType_GetGreatestSubtype(typeA, typeB)

	s.AssertType(greatestSubtype).IsSubtypeOf(typeA)
	s.AssertType(greatestSubtype).IsSubtypeOf(typeB)
	s.AssertType(greatestSubtype).IsNotEqualTo(typeA)
	s.AssertType(greatestSubtype).IsNotEqualTo(typeB)

	s.AssertType(typeA).IsNotSubtypeOf(greatestSubtype)
	s.AssertType(typeB).IsNotSubtypeOf(greatestSubtype)
}

func TestEnumElementType_GetGreatestSubtype_twoEnumElementTypes_postResolution(t *testing.T) {
	s := newEnumElementTypeTest(t)
	typeA := s.CreateEnumType("typeA", s.NUMBER_TYPE).GetElementsType()
	typeB := s.CreateEnumType("typeB", s.NUMBER_TYPE).GetElementsType()

	greatestSubtype := EnumElementType_GetGreatestSubtype(typeA, typeB)

	s.AssertType(greatestSubtype).IsSubtypeOf(typeA)
	s.AssertType(greatestSubtype).IsSubtypeOf(typeB)
	s.AssertType(greatestSubtype).IsNotEqualTo(typeA)
	s.AssertType(greatestSubtype).IsNotEqualTo(typeB)

	s.AssertType(typeA).IsNotSubtypeOf(greatestSubtype)
	s.AssertType(typeB).IsNotSubtypeOf(greatestSubtype)
}

func TestEnumElementType_EqualityOfEnumTypes_withSameReferenceName(t *testing.T) {
	s := newEnumElementTypeTest(t)
	firstFoo := s.CreateEnumType("Foo", s.NUMBER_TYPE)
	secondFoo := s.CreateEnumType("Foo", s.NUMBER_TYPE)

	s.AssertType(firstFoo).IsNotEqualTo(secondFoo)
	s.AssertType(firstFoo.GetElementsType()).IsNotEqualTo(secondFoo.GetElementsType())
}

func (s *enumElementTypeTest) CreateEnumType(name string, elementType JSType) *EnumType {
	return EnumType_Builder_(s.Registry).SetName(ptr(name)).SetElementType(elementType).Build()
}
