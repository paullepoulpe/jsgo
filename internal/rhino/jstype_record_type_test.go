// Port of com/google/javascript/rhino/jstype/RecordTypeTest.java @ closure-compiler v20261015

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

type recordTypeTest struct {
	*rhinotesting.BaseJSTypeTestCase
}

func newRecordTypeTest(t *testing.T) *recordTypeTest {
	s := &recordTypeTest{BaseJSTypeTestCase: rhinotesting.NewBaseJSTypeTestCase(t)}
	return s
}

func TestRecordType_RecursiveRecord(t *testing.T) {
	s := newRecordTypeTest(t)
	loop := NewProxyObjectType(s.Registry, s.NUMBER_TYPE)
	record := NewRecordTypeBuilder(s.Registry).
		AddProperty("loop", loop, NilNode).
		AddProperty("number", s.NUMBER_TYPE, NilNode).
		AddProperty("string", s.STRING_TYPE, NilNode).
		Build()
	s.AssertThat(record.String()).
		IsEqualTo("{\n  loop: number,\n  number: number,\n  string: string\n}")

	loop.SetReferencedType(record)
	s.AssertThat(record.String()).
		IsEqualTo("{\n  loop: {...},\n  number: number,\n  string: string\n}")
	s.AssertThat(record.ToAnnotationString(JSType_Nullability_EXPLICIT)).
		IsEqualTo("{loop: ?, number: number, string: string}")

	rhinotesting.Asserts_AssertEquivalenceOperations(s.T, record, loop)
}

func TestRecordType_ToString_largeRecord(t *testing.T) {
	s := newRecordTypeTest(t)
	record := NewRecordTypeBuilder(s.Registry).
		AddProperty("a01", s.NUMBER_TYPE, NilNode).
		AddProperty("a02", s.NUMBER_TYPE, NilNode).
		AddProperty("a03", s.NUMBER_TYPE, NilNode).
		AddProperty("a04", s.NUMBER_TYPE, NilNode).
		AddProperty("a05", s.NUMBER_TYPE, NilNode).
		AddProperty("a06", s.NUMBER_TYPE, NilNode).
		AddProperty("a07", s.NUMBER_TYPE, NilNode).
		AddProperty("a08", s.NUMBER_TYPE, NilNode).
		AddProperty("a09", s.NUMBER_TYPE, NilNode).
		AddProperty("a10", s.NUMBER_TYPE, NilNode).
		AddProperty("a11", s.NUMBER_TYPE, NilNode).
		Build()

	s.AssertThat(record.String()).
		IsEqualTo(
			rhinotesting.BaseJSTypeTestCase_Lines(
				"{",
				"  a01: number,",
				"  a02: number,",
				"  a03: number,",
				"  a04: number,",
				"  a05: number,",
				"  a06: number,",
				"  a07: number,",
				"  a08: number,",
				"  a09: number,",
				"  a10: number,",
				"  ...",
				"}"))

	s.AssertThat(record.ToAnnotationString(JSType_Nullability_EXPLICIT)).
		IsEqualTo(
			"{a01: number, a02: number, a03: number, a04: number, a05: number, a06: number," +
				" a07: number, a08: number, a09: number, a10: number, a11: number}")
}

func TestRecordType_ToString_nestedRecordIndentation(t *testing.T) {
	s := newRecordTypeTest(t)
	record :=
		NewRecordTypeBuilder(s.Registry).
			AddProperty(
				"a",
				NewRecordTypeBuilder(s.Registry).
					AddProperty(
						"aa",
						NewRecordTypeBuilder(s.Registry).
							AddProperty("aaa", s.NUMBER_TYPE, NilNode).
							AddProperty("aab", s.NUMBER_TYPE, NilNode).
							AddProperty("aac", s.NUMBER_TYPE, NilNode).
							Build(),
						NilNode).
					AddProperty("ab", s.NUMBER_TYPE, NilNode).
					Build(),
				NilNode).
			AddProperty("b", s.NUMBER_TYPE, NilNode).
			Build()

	s.AssertThat(record.String()).
		IsEqualTo(
			rhinotesting.BaseJSTypeTestCase_Lines(
				"{",
				"  a: {",
				"    aa: {",
				"      aaa: number,",
				"      aab: number,",
				"      aac: number",
				"    },",
				"    ab: number",
				"  },",
				"  b: number",
				"}"))

	s.AssertThat(record.ToAnnotationString(JSType_Nullability_EXPLICIT)).
		IsEqualTo("{a: {aa: {aaa: number, aab: number, aac: number}, ab: number}, b: number}")
}

func TestRecordType_ToString_nestedRecordIndentation_throughUnion(t *testing.T) {
	s := newRecordTypeTest(t)
	record :=
		NewRecordTypeBuilder(s.Registry).
			AddProperty("a", s.NUMBER_TYPE, NilNode).
			AddProperty(
				"b",
				s.Registry.CreateUnionType(
					s.NUMBER_TYPE,
					NewRecordTypeBuilder(s.Registry).
						AddProperty("ba", s.NUMBER_TYPE, NilNode).
						AddProperty("bb", s.NUMBER_TYPE, NilNode).
						Build()),
				NilNode).
			Build()

	s.AssertThat(record.String()).
		IsEqualTo(
			rhinotesting.BaseJSTypeTestCase_Lines(
				"{",
				"  a: number,",
				"  b: (number|{",
				"    ba: number,",
				"    bb: number",
				"  })",
				"}"))

	s.AssertThat(record.ToAnnotationString(JSType_Nullability_EXPLICIT)).
		IsEqualTo("{a: number, b: (number|{ba: number, bb: number})}")
}

func TestRecordType_ToString_nestedRecordIndentation_throughTemplateParam(t *testing.T) {
	s := newRecordTypeTest(t)
	record :=
		NewRecordTypeBuilder(s.Registry).
			AddProperty("a", s.NUMBER_TYPE, NilNode).
			AddProperty(
				"b",
				s.Registry.CreateTemplatizedTypeOfList(
					s.ARRAY_TYPE,
					[]JSType{
						NewRecordTypeBuilder(s.Registry).
							AddProperty("ba", s.NUMBER_TYPE, NilNode).
							AddProperty("bb", s.NUMBER_TYPE, NilNode).
							Build()}),
				NilNode).
			Build()

	s.AssertThat(record.String()).
		IsEqualTo(
			rhinotesting.BaseJSTypeTestCase_Lines(
				"{",
				"  a: number,",
				"  b: Array<{",
				"    ba: number,",
				"    bb: number",
				"  }>",
				"}"))

	s.AssertThat(record.ToAnnotationString(JSType_Nullability_EXPLICIT)).
		IsEqualTo("{a: number, b: !Array<{ba: number, bb: number}>}")
}

func TestRecordType_SupAndInf(t *testing.T) {
	s := newRecordTypeTest(t)
	recordA := NewRecordTypeBuilder(s.Registry).
		AddProperty("a", s.NUMBER_TYPE, NilNode).
		AddProperty("b", s.NUMBER_TYPE, NilNode).
		Build()
	recordC := NewRecordTypeBuilder(s.Registry).
		AddProperty("b", s.NUMBER_TYPE, NilNode).
		AddProperty("c", s.NUMBER_TYPE, NilNode).
		Build()
	proxyRecordA := NewProxyObjectType(s.Registry, recordA)
	proxyRecordC := NewProxyObjectType(s.Registry, recordC)

	aInfC := NewRecordTypeBuilder(s.Registry).
		AddProperty("a", s.NUMBER_TYPE, NilNode).
		AddProperty("b", s.NUMBER_TYPE, NilNode).
		AddProperty("c", s.NUMBER_TYPE, NilNode).
		Build()

	aSupC := s.Registry.CreateUnionType(recordA, recordC)

	s.AssertType(recordA.GetGreatestSubtype(recordC)).IsEqualTo(aInfC)
	s.AssertType(recordA.GetLeastSupertype(recordC)).IsEqualTo(aSupC)

	s.AssertType(proxyRecordA.GetGreatestSubtype(proxyRecordC)).IsEqualTo(aInfC)
	s.AssertType(proxyRecordA.GetLeastSupertype(proxyRecordC)).IsEqualTo(aSupC)
}

func TestRecordType_SubtypeWithUnknowns(t *testing.T) {
	s := newRecordTypeTest(t)
	recordA := NewRecordTypeBuilder(s.Registry).
		AddProperty("a", s.NUMBER_TYPE, NilNode).
		Build()
	recordB := NewRecordTypeBuilder(s.Registry).
		AddProperty("a", s.UNKNOWN_TYPE, NilNode).
		Build()
	s.AssertThat(recordA.IsSubtypeOf(recordB)).IsTrue()
	s.AssertThat(recordB.IsSubtypeOf(recordA)).IsTrue()
}

func TestRecordType_SubtypeWithUnknowns2(t *testing.T) {
	s := newRecordTypeTest(t)
	recordA :=
		NewRecordTypeBuilder(s.Registry).
			AddProperty(
				"a", FunctionType_Builder_(s.Registry).WithReturnType(s.NUMBER_TYPE).Build(), NilNode).
			Build()
	recordB :=
		NewRecordTypeBuilder(s.Registry).
			AddProperty(
				"a", FunctionType_Builder_(s.Registry).WithReturnType(s.UNKNOWN_TYPE).Build(), NilNode).
			Build()
	s.AssertThat(recordA.IsSubtypeOf(recordB)).IsTrue()
	s.AssertThat(recordB.IsSubtypeOf(recordA)).IsTrue()
}

func TestRecordType_SubtypeWithFunctionProps(t *testing.T) {
	s := newRecordTypeTest(t)
	recordA :=
		NewRecordTypeBuilder(s.Registry).
			AddProperty(
				"a", FunctionType_Builder_(s.Registry).WithReturnType(s.NUMBER_TYPE).Build(), NilNode).
			Build()
	recordB :=
		NewRecordTypeBuilder(s.Registry).
			AddProperty(
				"a", FunctionType_Builder_(s.Registry).WithReturnType(s.STRING_TYPE).Build(), NilNode).
			Build()
	s.AssertThat(recordA.IsSubtypeOf(recordB)).IsFalse()
	s.AssertThat(recordB.IsSubtypeOf(recordA)).IsFalse()
}

func TestRecordType_SubtypeWithManyProps(t *testing.T) {
	s := newRecordTypeTest(t)
	recordA := NewRecordTypeBuilder(s.Registry).
		AddProperty("a", s.NUMBER_TYPE, NilNode).
		AddProperty("b", s.NUMBER_TYPE, NilNode).
		Build()
	recordB := NewRecordTypeBuilder(s.Registry).
		AddProperty("a", s.NUMBER_TYPE, NilNode).
		AddProperty("b", s.STRING_TYPE, NilNode).
		Build()
	recordC := NewRecordTypeBuilder(s.Registry).
		AddProperty("a", s.NUMBER_TYPE, NilNode).
		AddProperty("b",
			s.Registry.CreateUnionType(s.NUMBER_TYPE, s.STRING_TYPE), NilNode).
		Build()
	s.AssertThat(recordA.IsSubtypeOf(recordB)).IsFalse()
	s.AssertThat(recordB.IsSubtypeOf(recordA)).IsFalse()
	s.AssertThat(recordC.IsSubtypeOf(recordB)).IsFalse()
	s.AssertThat(recordB.IsSubtypeOf(recordC)).IsTrue()
	s.AssertThat(recordA.IsSubtypeOf(recordC)).IsTrue()
}

func TestRecordType_Equality_onlyCompares_ownProperties(t *testing.T) {
	s := newRecordTypeTest(t)
	// Given
	recordA := //
		NewRecordTypeBuilder(s.Registry).AddProperty("a", s.NUMBER_TYPE, NilNode).Build().(*RecordType)
	recordB :=
		NewRecordTypeBuilder(s.Registry).
			AddProperty("a", s.NUMBER_TYPE, NilNode).
			AddProperty("toString", s.OBJECT_TYPE.GetPropertyType("toString"), NilNode).
			Build().(*RecordType)

	s.AssertType(recordA.GetPropertyType("toString")).IsEqualTo(recordB.GetPropertyType("toString"))

	// Then
	s.AssertType(recordA).IsNotEqualTo(recordB)
}
