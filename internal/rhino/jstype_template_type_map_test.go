// Port of com/google/javascript/rhino/jstype/TemplateTypeMapTest.java @ closure-compiler v20261021

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

	"github.com/paullepoulpe/jsgo/internal/jutil"
	. "github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/rhinotesting"
)

type templateTypeMapTest struct {
	*rhinotesting.BaseJSTypeTestCase
	emptyMap *TemplateTypeMap
}

func newTemplateTypeMapTest(t *testing.T) *templateTypeMapTest {
	s := &templateTypeMapTest{BaseJSTypeTestCase: rhinotesting.NewBaseJSTypeTestCase(t)}
	s.SetUp()
	return s
}

func (s *templateTypeMapTest) SetUp() {
	s.emptyMap = s.Registry.GetEmptyTemplateTypeMap()
}

// assertThatMap is the statically imported TemplateTypeMapSubject.assertThat.
func (s *templateTypeMapTest) assertThatMap(m *TemplateTypeMap) *rhinotesting.TemplateTypeMapSubject {
	return rhinotesting.TemplateTypeMapSubject_AssertThat(s.T, m)
}

func TestTemplateTypeMap_CreateEmptyMap_throwsWhenDuplicateRequested(t *testing.T) {
	s := newTemplateTypeMapTest(t)
	rhinotesting.AssertThrows[jutil.Throwable](s.T, func() { TemplateTypeMap_CreateEmpty(s.Registry) })
}

func TestTemplateTypeMap_EmptyMap_hasNoKeysValuesOrSubmaps(t *testing.T) {
	s := newTemplateTypeMapTest(t)
	s.assertThatMap(s.emptyMap).HasKeysAndValues([]*TemplateType{}, []JSType{})
	s.assertThatMap(s.emptyMap).HasSubmapsAt()
}

func TestTemplateTypeMap_CopyExtend_createsNewMap(t *testing.T) {
	s := newTemplateTypeMapTest(t)
	// When
	result :=
		s.emptyMap.CopyWithExtensionOfKeys([]*TemplateType{s.key("A")}, []JSType{s.NUMBER_TYPE})

	// Then
	s.AssertThat(result).IsNotSameInstanceAs(s.emptyMap)
	s.assertThatMap(result).HasSubmapsAt(0)
}

func TestTemplateTypeMap_CopyExtend_validatesCountOfKeysVsValues(t *testing.T) {
	s := newTemplateTypeMapTest(t)
	rhinotesting.AssertThrows[jutil.Throwable](
		s.T,
		func() { s.emptyMap.CopyWithExtensionOfKeys([]*TemplateType{}, []JSType{s.NUMBER_TYPE}) })
}

func TestTemplateTypeMap_CopyExtend_emptyMap(t *testing.T) {
	s := newTemplateTypeMapTest(t)
	// Given
	keyT := s.key("T")
	keyU := s.key("U")

	// When
	result :=
		s.emptyMap.CopyWithExtensionOfKeys([]*TemplateType{keyT, keyU}, []JSType{s.NUMBER_TYPE})

	// Then
	s.assertThatMap(result).
		HasKeysAndValues([]*TemplateType{keyT, keyU}, []JSType{s.NUMBER_TYPE})
	s.assertThatMap(result).HasSubmapsAt(0)
}

func TestTemplateTypeMap_CopyExtend_partialMap_fillsExitingUnfilledKeysWithUnknown(t *testing.T) {
	s := newTemplateTypeMapTest(t)
	// Given
	keyA := s.key("A")
	keyB := s.key("B")
	keyC := s.key("C")
	existing :=
		s.createMap([]*TemplateType{keyA, keyB, keyC}, []JSType{s.NUMBER_TYPE})

	keyX := s.key("X")
	keyY := s.key("Y")
	keyZ := s.key("Z")

	// When
	result :=
		existing.CopyWithExtensionOfKeys(
			[]*TemplateType{keyX, keyY, keyZ}, []JSType{s.STRING_TYPE})

	// Then
	s.assertThatMap(result).
		HasKeysAndValues(
			[]*TemplateType{keyA, keyB, keyC, keyX, keyY, keyZ},
			[]JSType{s.NUMBER_TYPE, s.UNKNOWN_TYPE, s.UNKNOWN_TYPE, s.STRING_TYPE})
	s.assertThatMap(result).HasSubmapsAt(0, 3)
}

func TestTemplateTypeMap_CopyExtend_partialMap_emptyExtension_fillsExitingUnfilledKeysWithUnknown(t *testing.T) {
	s := newTemplateTypeMapTest(t)
	// Given
	keyA := s.key("A")
	keyB := s.key("B")
	keyC := s.key("C")
	existing :=
		s.createMap([]*TemplateType{keyA, keyB, keyC}, []JSType{s.NUMBER_TYPE})
	// When
	result := existing.CopyWithExtensionOfKeys([]*TemplateType{}, []JSType{})

	// Then
	s.assertThatMap(result).
		HasKeysAndValues(
			[]*TemplateType{keyA, keyB, keyC},
			[]JSType{s.NUMBER_TYPE, s.UNKNOWN_TYPE, s.UNKNOWN_TYPE})
	s.assertThatMap(result).HasSubmapsAt(0)
}

func TestTemplateTypeMap_CopyExtend_ifMapFilled_ifExtensionEmpty_returnsSelf(t *testing.T) {
	s := newTemplateTypeMapTest(t)
	// Given
	existing := s.createMap([]*TemplateType{s.key("A")}, []JSType{s.NUMBER_TYPE})

	// When
	result := s.emptyMap.CopyWithExtensionOfKeys([]*TemplateType{}, []JSType{})

	// Then
	s.AssertThat(result).IsNotSameInstanceAs(existing)
}

func TestTemplateTypeMap_CopyFill_validatesValueCount_againstUnfilledKeys(t *testing.T) {
	s := newTemplateTypeMapTest(t)
	rhinotesting.AssertThrows[jutil.Throwable](
		s.T, func() { s.emptyMap.CopyFilledWithValues([]JSType{s.NUMBER_TYPE}) })
}

func TestTemplateTypeMap_CopyFill_ifNoSlotsEmpty_returnsSelf(t *testing.T) {
	s := newTemplateTypeMapTest(t)
	// Given
	keyA := s.key("A")
	keyB := s.key("B")
	existing :=
		s.createMap([]*TemplateType{keyA, keyB}, []JSType{s.NUMBER_TYPE, s.STRING_TYPE})

	// When
	result := existing.CopyFilledWithValues([]JSType{})

	// Then
	s.AssertThat(result).IsSameInstanceAs(existing)
}

func TestTemplateTypeMap_CopyFill_partialMap_fillsExtraKeysWithUnknown(t *testing.T) {
	s := newTemplateTypeMapTest(t)
	// Given
	keyA := s.key("A")
	keyB := s.key("B")
	keyC := s.key("C")
	existing :=
		s.createMap([]*TemplateType{keyA, keyB, keyC}, []JSType{s.NUMBER_TYPE})

	// When
	result := existing.CopyFilledWithValues([]JSType{s.STRING_TYPE})

	// Then
	s.assertThatMap(result).
		HasKeysAndValues(
			[]*TemplateType{keyA, keyB, keyC},
			[]JSType{s.NUMBER_TYPE, s.STRING_TYPE, s.UNKNOWN_TYPE})
	s.assertThatMap(result).HasSubmapsAt(0)
}

func TestTemplateTypeMap_CopyWithout_canRemoveUnfilledKeys(t *testing.T) {
	s := newTemplateTypeMapTest(t)
	// Given
	keyA := s.key("A")
	keyB := s.key("B")
	keyC := s.key("C")
	existing :=
		s.createMap([]*TemplateType{keyA, keyB, keyC}, []JSType{s.NUMBER_TYPE})

	// When
	result := existing.CopyWithoutKeySlice([]*TemplateType{keyB})

	// Then
	s.assertThatMap(result).
		HasKeysAndValues([]*TemplateType{keyA, keyC}, []JSType{s.NUMBER_TYPE})
	s.assertThatMap(result).HasSubmapsAt(0)
}

func TestTemplateTypeMap_CopyWithout_preservesSubmapBoundary(t *testing.T) {
	s := newTemplateTypeMapTest(t)
	// Given
	keyA := s.key("A")
	keyB := s.key("B")
	keyC := s.key("C")
	keyD := s.key("D")
	keyE := s.key("E")
	existing :=
		s.createMap([]*TemplateType{keyA, keyB}, []JSType{s.NUMBER_TYPE, s.STRING_TYPE})
	extension :=
		existing.CopyWithExtensionOfKeys([]*TemplateType{keyC, keyD, keyE}, []JSType{})

	// When
	result := extension.CopyWithoutKeySlice([]*TemplateType{keyC, keyE})

	// Then
	s.assertThatMap(result).
		HasKeysAndValues(
			[]*TemplateType{keyA, keyB, keyD}, []JSType{s.NUMBER_TYPE, s.STRING_TYPE})
	s.assertThatMap(result).HasSubmapsAt(0, 2)
}

func TestTemplateTypeMap_CopyWithout_retainFilledKeys(t *testing.T) {
	s := newTemplateTypeMapTest(t)
	// Given
	keyA := s.key("A")
	keyB := s.key("B")
	keyC := s.key("C")
	existing :=
		s.createMap([]*TemplateType{keyA, keyB, keyC}, []JSType{s.NUMBER_TYPE})

	// When
	result := existing.CopyWithoutKeySlice([]*TemplateType{keyA})

	// Then
	s.assertThatMap(result).
		HasKeysAndValues([]*TemplateType{keyA, keyB, keyC}, []JSType{s.NUMBER_TYPE})
	s.assertThatMap(result).HasSubmapsAt(0)
}

func TestTemplateTypeMap_CopyWithout_ifNothingRemoved_returnsSelf(t *testing.T) {
	s := newTemplateTypeMapTest(t)
	// Given
	keyA := s.key("A")
	keyB := s.key("B")
	keyC := s.key("C")
	existing :=
		s.createMap([]*TemplateType{keyA, keyB, keyC}, []JSType{s.NUMBER_TYPE})

	// When
	result := existing.CopyWithoutKeySlice([]*TemplateType{keyA})

	// Then
	s.AssertThat(result).IsSameInstanceAs(existing)
}

func TestTemplateTypeMap_GetLastTemplateTypeKeyByName_returnsLastKeyIfDuplicates(t *testing.T) {
	s := newTemplateTypeMapTest(t)
	key1 := s.key("A")
	key2 := s.key("A")

	map_ := s.createMap([]*TemplateType{key1, key2}, []JSType{s.NUMBER_TYPE})

	result := map_.GetLastTemplateTypeKeyByName("A")

	s.AssertType(result).IsSameInstanceAs(key2)
}

func (s *templateTypeMapTest) createMap(keys []*TemplateType, values []JSType) *TemplateTypeMap {
	return s.emptyMap.CopyWithExtensionOfKeys(keys, values)
}

func (s *templateTypeMapTest) key(name string) *TemplateType {
	return s.Registry.CreateTemplateType(name)
}
