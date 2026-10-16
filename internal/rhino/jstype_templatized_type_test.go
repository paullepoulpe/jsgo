// Port of com/google/javascript/rhino/jstype/TemplatizedTypeTest.java @ closure-compiler v20260415

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
 *   John Lenz
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

type templatizedTypeTest struct {
	*rhinotesting.BaseJSTypeTestCase
}

func newTemplatizedTypeTest(t *testing.T) *templatizedTypeTest {
	s := &templatizedTypeTest{BaseJSTypeTestCase: rhinotesting.NewBaseJSTypeTestCase(t)}
	return s
}

/** Tests the behavior of variants type. */
func TestTemplatizedType_TemplatizedType(t *testing.T) {
	s := newTemplatizedTypeTest(t)
	arrOfString := s.CreateTemplatizedType(
		s.ARRAY_TYPE, s.STRING_TYPE)
	s.AssertTypeCanAssignToItself(arrOfString)
	s.AssertThat(arrOfString.IsSubtype(s.ARRAY_TYPE)).IsTrue()
	s.AssertThat(s.ARRAY_TYPE.IsSubtypeOf(arrOfString)).IsTrue()

	arrOfNumber := s.CreateTemplatizedType(
		s.ARRAY_TYPE, s.NUMBER_TYPE)
	s.AssertTypeCanAssignToItself(arrOfNumber)
	s.AssertThat(arrOfNumber.IsSubtype(s.ARRAY_TYPE)).IsTrue()
	s.AssertThat(s.ARRAY_TYPE.IsSubtypeOf(arrOfNumber)).IsTrue()

	s.AssertThat(arrOfString.Equals(s.CreateTemplatizedType(s.ARRAY_TYPE, s.STRING_TYPE))).IsTrue()

	s.AssertThat(arrOfString.Equals(s.ARRAY_TYPE)).IsFalse()
	s.AssertThat(arrOfString.Equals(s.ARRAY_TYPE)).IsFalse()
	s.AssertThat(arrOfString.Equals(arrOfNumber)).IsFalse()
	s.AssertThat(arrOfNumber.Equals(arrOfString)).IsFalse()
}

func TestTemplatizedType_Print1(t *testing.T) {
	s := newTemplatizedTypeTest(t)
	arrOfString := s.CreateTemplatizedType(
		s.ARRAY_TYPE, s.STRING_TYPE)
	s.AssertThat(arrOfString.String()).IsEqualTo("Array<string>")
}

func TestTemplatizedType_Print2(t *testing.T) {
	s := newTemplatizedTypeTest(t)
	arrOfTemplateType := s.CreateTemplatizedType(
		s.ARRAY_TYPE, NewTemplateType(s.Registry, "T"))
	s.AssertThat(arrOfTemplateType.String()).IsEqualTo("Array<T>")
}

func TestTemplatizedType_Print3(t *testing.T) {
	s := newTemplatizedTypeTest(t)
	arrOfUnknown := s.CreateTemplatizedType(
		s.ARRAY_TYPE, s.UNKNOWN_TYPE)
	s.AssertThat(arrOfUnknown.String()).IsEqualTo("Array<?>")
}

func TestTemplatizedType_PrintingRawType(t *testing.T) {
	s := newTemplatizedTypeTest(t)
	rawType := s.CreateCustomTemplatizedType("Foo")

	s.AssertThat(rawType.String()).IsEqualTo("Foo")
}

func TestTemplatizedType_DifferentRawTypes(t *testing.T) {
	s := newTemplatizedTypeTest(t)
	arrOfNumber := s.CreateTemplatizedType(
		s.ARRAY_TYPE, s.NUMBER_TYPE)
	objType := s.CreateTemplatizedType(
		s.OBJECT_TYPE, s.UNKNOWN_TYPE)
	s.AssertThat(arrOfNumber.IsSubtype(objType)).IsTrue()
	s.AssertThat(objType.IsSubtype(arrOfNumber)).IsFalse()
}

func TestTemplatizedType_SubtypingAndEquivalenceAmongCustomTemplatizedTypes(t *testing.T) {
	s := newTemplatizedTypeTest(t)
	rawType := s.CreateCustomTemplatizedType("Baz")

	templatizedStringNumber :=
		s.Registry.CreateTemplatizedTypeOfList(rawType, []JSType{s.STRING_TYPE, s.NUMBER_TYPE})
	secondTemplatizedStringNumber :=
		s.Registry.CreateTemplatizedTypeOfList(rawType, []JSType{s.STRING_TYPE, s.NUMBER_TYPE})
	templatizedStringAll :=
		s.Registry.CreateTemplatizedTypeOfList(rawType, []JSType{s.STRING_TYPE, s.ALL_TYPE})
	templatizedStringUnknown :=
		s.Registry.CreateTemplatizedTypeOfList(rawType, []JSType{s.STRING_TYPE, s.UNKNOWN_TYPE})
	templatizedUnknownUnknown :=
		s.Registry.CreateTemplatizedTypeOfList(rawType, []JSType{s.UNKNOWN_TYPE, s.UNKNOWN_TYPE})

	s.AssertType(templatizedStringNumber).IsSubtypeOf(rawType)
	s.AssertType(templatizedStringAll).IsSubtypeOf(rawType)
	s.AssertType(templatizedStringUnknown).IsSubtypeOf(rawType)
	s.AssertType(templatizedUnknownUnknown).IsSubtypeOf(rawType)

	s.AssertType(templatizedStringNumber).IsNotEqualTo(rawType)
	s.AssertType(templatizedStringAll).IsNotEqualTo(rawType)
	s.AssertType(templatizedStringUnknown).IsNotEqualTo(rawType)

	s.AssertType(templatizedStringNumber).IsEqualTo(secondTemplatizedStringNumber)
	// TODO(b/110224889): This case should probably be `assertTypeNotEquals`.
	s.AssertType(templatizedUnknownUnknown).IsEqualTo(rawType)

	s.AssertType(rawType).IsSubtypeOf(templatizedStringNumber)
	s.AssertType(rawType).IsSubtypeOf(templatizedStringAll)
	s.AssertType(rawType).IsSubtypeOf(templatizedStringUnknown)
	s.AssertType(rawType).IsSubtypeOf(templatizedUnknownUnknown)

	s.AssertType(templatizedStringNumber).IsSubtypeOf(secondTemplatizedStringNumber)
	s.AssertType(templatizedStringNumber).IsNotSubtypeOf(templatizedStringAll)
	s.AssertType(templatizedStringAll).IsNotSubtypeOf(templatizedStringNumber)

	s.AssertType(templatizedStringAll).IsSubtypeOf(templatizedStringUnknown)
	s.AssertType(templatizedStringUnknown).IsSubtypeOf(templatizedStringAll)
}

func TestTemplatizedType_EqualityWithRawType_whenSpecializedOnUnknown_includingHashcode(t *testing.T) {
	s := newTemplatizedTypeTest(t)
	rawType := s.CreateCustomTemplatizedType("Baz")

	templatizedUnknownString :=
		s.Registry.CreateTemplatizedTypeOfList(rawType, []JSType{s.UNKNOWN_TYPE, s.STRING_TYPE})
	s.AssertType(rawType).IsNotEqualTo(templatizedUnknownString)

	templatizedStringUnknown :=
		s.Registry.CreateTemplatizedTypeOfList(rawType, []JSType{s.STRING_TYPE, s.UNKNOWN_TYPE})
	s.AssertType(rawType).IsNotEqualTo(templatizedStringUnknown)

	templatizedUnknownUnknown :=
		s.Registry.CreateTemplatizedTypeOfList(rawType, []JSType{s.UNKNOWN_TYPE, s.UNKNOWN_TYPE})
	s.AssertType(rawType).IsEqualTo(templatizedUnknownUnknown)

	templatizedCheckedUnknownUnknown :=
		s.Registry.CreateTemplatizedTypeOfList(
			rawType, []JSType{s.CHECKED_UNKNOWN_TYPE, s.UNKNOWN_TYPE})
	s.AssertType(rawType).IsNotEqualTo(templatizedCheckedUnknownUnknown)

	ramdomTemplate := s.Registry.CreateTemplateType("X")
	templatizedTemplateUnknown :=
		s.Registry.CreateTemplatizedTypeOfList(rawType, []JSType{ramdomTemplate, s.UNKNOWN_TYPE})
	s.AssertType(ramdomTemplate).IsUnknown()
	s.AssertType(rawType).IsNotEqualTo(templatizedTemplateUnknown)

	var resolvedNamedUnknown *NamedType
	func() {
		closer := s.Registry.GetResolver().OpenForDefinition()
		defer closer.Close()
		resolvedNamedUnknown = s.Registry.CreateNamedType(nil, "Y", ptr(""), -1, -1)
		resolvedNamedUnknown.SetReferencedType(s.UNKNOWN_TYPE)
	}()
	templatizedResolvedNamedUnknownUnknown :=
		s.Registry.CreateTemplatizedTypeOfList(
			rawType, []JSType{resolvedNamedUnknown, s.UNKNOWN_TYPE})
	s.AssertType(resolvedNamedUnknown).IsUnknown()
	s.AssertType(rawType).IsEqualTo(templatizedResolvedNamedUnknownUnknown)

	closer := s.Registry.GetResolver().OpenForDefinition()
	defer closer.Close()
	{
		unresolvedNamedUnknown := s.Registry.CreateNamedType(nil, "Y", ptr(""), -1, -1)
		unresolvedNamedUnknown.SetReferencedType(s.UNKNOWN_TYPE)

		templatizedUnresolvedNamedUnknownUnknown :=
			s.Registry.CreateTemplatizedTypeOfList(
				rawType, []JSType{unresolvedNamedUnknown, s.UNKNOWN_TYPE})
		s.AssertType(resolvedNamedUnknown).IsUnknown()
		// Not equal because we don't know yet whether the type args are all unknown.
		s.AssertType(rawType).IsNotEqualTo(templatizedUnresolvedNamedUnknownUnknown)
	}
}

func TestTemplatizedType_GetPropertyTypeOnTemplatizedType(t *testing.T) {
	s := newTemplatizedTypeTest(t)
	closer := s.Registry.GetResolver().OpenForDefinition()
	defer closer.Close()
	{
		templateT := s.Registry.CreateTemplateType("T")
		ctor := // function<T>(new:Foo<T>)
			s.Registry.CreateConstructorType(
				ptr("Foo"), nil, nil, nil, []*TemplateType{templateT}, false)
		rawType := ctor.GetInstanceType() // Foo<T> == Foo
		rawType.DefineDeclaredProperty("property", templateT, nil)

		templatizedNumber := s.Registry.CreateTemplatizedType(rawType, s.NUMBER_TYPE)
		s.AssertType(templatizedNumber.ToObjectType().GetPropertyType("property")).
			IsEqualTo(s.NUMBER_TYPE)
	}
}

func TestTemplatizedType_FindPropertyTypeOnTemplatizedType(t *testing.T) {
	s := newTemplatizedTypeTest(t)
	closer := s.Registry.GetResolver().OpenForDefinition()
	defer closer.Close()
	{
		templateT := s.Registry.CreateTemplateType("T")
		ctor := // function<T>(new:Foo<T>)
			s.Registry.CreateConstructorType(
				ptr("Foo"), nil, nil, nil, []*TemplateType{templateT}, false)
		rawType := ctor.GetInstanceType() // Foo<T> == Foo
		rawType.DefineDeclaredProperty("property", templateT, nil)

		templatizedNumber := s.Registry.CreateTemplatizedType(rawType, s.NUMBER_TYPE)
		s.AssertType(templatizedNumber.FindPropertyType("property")).IsEqualTo(s.NUMBER_TYPE)
	}
}

/** Returns an unspecialized type with the provided name and two type parameters. */
func (s *templatizedTypeTest) CreateCustomTemplatizedType(rawName string) ObjectType {
	closer := s.Registry.GetResolver().OpenForDefinition()
	defer closer.Close()
	{
		ctor := // function<T,U>(new:Foo<T,U>)
			s.Registry.CreateConstructorType(
				ptr(rawName),
				nil,
				nil,
				nil,
				[]*TemplateType{s.Registry.CreateTemplateType("T"), s.Registry.CreateTemplateType("U")},
				false)
		return ctor.GetInstanceType() // Foo<T, U> == Foo
	}
}

/** Assert that a type can assign to itself. */
func (s *templatizedTypeTest) AssertTypeCanAssignToItself(type_ JSType) {
	s.AssertThat(type_.IsSubtypeOf(type_)).IsTrue()
}
