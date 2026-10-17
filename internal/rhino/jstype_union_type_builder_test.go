// Port of com/google/javascript/rhino/jstype/UnionTypeBuilderTest.java @ closure-compiler v20261015

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

type unionTypeBuilderTest struct {
	*rhinotesting.BaseJSTypeTestCase
	EMPTY_SCOPE *rhinotesting.MapBasedScope
	base        ObjectType
	sub         ObjectType
}

func newUnionTypeBuilderTest(t *testing.T) *unionTypeBuilderTest {
	s := &unionTypeBuilderTest{BaseJSTypeTestCase: rhinotesting.NewBaseJSTypeTestCase(t)}
	s.EMPTY_SCOPE = rhinotesting.MapBasedScope_EmptyScope()
	s.SetUp()
	return s
}

func (s *unionTypeBuilderTest) SetUp() {
	closer := s.Registry.GetResolver().OpenForDefinition()
	defer closer.Close()
	{
		baseCtor :=
			FunctionType_Builder_(s.Registry).ForConstructor().WithName(ptr("Base")).Build()
		s.base = baseCtor.GetInstanceType()
		subCtor :=
			FunctionType_Builder_(s.Registry).
				ForConstructor().
				WithName(ptr("Sub")).
				WithPrototypeBasedOn(s.base).
				Build()
		s.sub = subCtor.GetInstanceType()
	}
}

func TestUnionTypeBuilder_WildcardType_allType(t *testing.T) {
	s := newUnionTypeBuilderTest(t)
	s.AssertUnion("*", s.ALL_TYPE)
	s.AssertUnion("*", s.ALL_TYPE, s.ALL_TYPE)
	s.AssertUnion("*", s.NUMBER_TYPE, s.ALL_TYPE)
	s.AssertUnion("*", s.ALL_TYPE, s.NUMBER_TYPE)
	s.AssertUnion("*", s.ALL_TYPE, s.NUMBER_TYPE, s.NO_TYPE)

	s.AssertUnion("(*|undefined)", s.ALL_TYPE, s.VOID_TYPE)

	s.AssertUnion("*", s.ALL_TYPE, s.UNKNOWN_TYPE)         // TODO(b/137892871): Should be `?`.
	s.AssertUnion("*", s.ALL_TYPE, s.CHECKED_UNKNOWN_TYPE) // TODO(b/137892871): Should be `??`.
}

func TestUnionTypeBuilder_WildcardType_unknownType(t *testing.T) {
	s := newUnionTypeBuilderTest(t)
	s.AssertUnion("?", s.UNKNOWN_TYPE)
	s.AssertUnion("?", s.UNKNOWN_TYPE, s.UNKNOWN_TYPE)
	s.AssertUnion("?", s.NUMBER_TYPE, s.UNKNOWN_TYPE)
	s.AssertUnion("?", s.UNKNOWN_TYPE, s.NUMBER_TYPE)
	s.AssertUnion("?", s.UNKNOWN_TYPE, s.NUMBER_TYPE, s.NO_TYPE)

	s.AssertUnion("(?|undefined)", s.UNKNOWN_TYPE, s.VOID_TYPE)

	s.AssertUnion("*", s.ALL_TYPE, s.UNKNOWN_TYPE) // TODO(b/137892871): Should be `?`.
	s.AssertUnion("?", s.CHECKED_UNKNOWN_TYPE, s.UNKNOWN_TYPE)
}

func TestUnionTypeBuilder_WildcardType_checkedUnknownType(t *testing.T) {
	s := newUnionTypeBuilderTest(t)
	s.AssertUnion("??", s.CHECKED_UNKNOWN_TYPE)
	s.AssertUnion("??", s.CHECKED_UNKNOWN_TYPE, s.CHECKED_UNKNOWN_TYPE)
	s.AssertUnion("??", s.NUMBER_TYPE, s.CHECKED_UNKNOWN_TYPE)
	s.AssertUnion("??", s.CHECKED_UNKNOWN_TYPE, s.NUMBER_TYPE)
	s.AssertUnion("??", s.CHECKED_UNKNOWN_TYPE, s.NUMBER_TYPE, s.NO_TYPE)

	s.AssertUnion("(??|undefined)", s.CHECKED_UNKNOWN_TYPE, s.VOID_TYPE)

	s.AssertUnion("*", s.ALL_TYPE, s.CHECKED_UNKNOWN_TYPE) // TODO(b/137892871): Should be `??`.
	s.AssertUnion("?", s.UNKNOWN_TYPE, s.CHECKED_UNKNOWN_TYPE)
}

func TestUnionTypeBuilder_EmptyUnion(t *testing.T) {
	s := newUnionTypeBuilderTest(t)
	s.AssertUnion("None")
	s.AssertUnion("None", s.NO_TYPE, s.NO_TYPE)
}

func TestUnionTypeBuilder_UnionTypes(t *testing.T) {
	s := newUnionTypeBuilderTest(t)
	union := s.Registry.CreateUnionType(s.STRING_TYPE, s.OBJECT_TYPE)

	s.AssertUnion("*", s.ALL_TYPE, union)
	s.AssertUnion("(Object|string)", s.OBJECT_TYPE, union)
	s.AssertUnion("(Object|string)", union, s.OBJECT_TYPE)
	s.AssertUnion("(Object|number|string)", s.NUMBER_TYPE, union)
	s.AssertUnion("(Object|number|string)", union, s.NUMBER_TYPE)
	s.AssertUnion("(Object|boolean|number|string)", union,
		s.Registry.CreateUnionType(s.NUMBER_TYPE, s.BOOLEAN_TYPE))
	s.AssertUnion("(Object|boolean|number|string)",
		s.Registry.CreateUnionType(s.NUMBER_TYPE, s.BOOLEAN_TYPE), union)
	s.AssertUnion("(Object|string)", union, s.STRING_OBJECT_TYPE)
}

func TestUnionTypeBuilder_UnknownTypes(t *testing.T) {
	s := newUnionTypeBuilderTest(t)
	s.AssertUnion("?", s.UNKNOWN_TYPE)
	s.AssertUnion("?", s.UNKNOWN_TYPE, s.UNKNOWN_TYPE)
	s.AssertUnion("(?|undefined)", s.UNKNOWN_TYPE, s.VOID_TYPE)
	s.AssertUnion("(?|undefined)", s.VOID_TYPE, s.UNKNOWN_TYPE)
	s.AssertUnion("(?|undefined)", s.VOID_TYPE, s.NUMBER_TYPE, s.UNKNOWN_TYPE)

	s.AssertUnion("(*|undefined)", s.ALL_TYPE, s.VOID_TYPE, s.NULL_TYPE)
}

func TestUnionTypeBuilder_UnresolvedNamedTypes(t *testing.T) {
	s := newUnionTypeBuilderTest(t)
	s.ErrorReporter.ExpectAllWarnings(
		"Bad type annotation. Unknown type not.resolved.A",
		"Bad type annotation. Unknown type not.resolved.B")

	closer := s.Registry.GetResolver().OpenForDefinition()
	defer closer.Close()
	{
		unresolvedNameA1 := s.Registry.CreateNamedType(s.EMPTY_SCOPE, "not.resolved.A", ptr(""), -1, -1)
		unresolvedNameA2 := s.Registry.CreateNamedType(s.EMPTY_SCOPE, "not.resolved.A", ptr(""), -1, -1)
		unresolvedNameB := s.Registry.CreateNamedType(s.EMPTY_SCOPE, "not.resolved.B", ptr(""), -1, -1)

		// NOTE: "(?)" means there are multiple unknown types in the union.
		s.AssertUnion("?", s.UNKNOWN_TYPE, unresolvedNameA1)
		s.AssertUnion("(not.resolved.A|not.resolved.A)", unresolvedNameA1, unresolvedNameA2)
		s.AssertUnion("(not.resolved.A|not.resolved.B)", unresolvedNameA1, unresolvedNameB)
		s.AssertUnion("(Object|not.resolved.A)", unresolvedNameA1, s.OBJECT_TYPE)
	}
}

func TestUnionTypeBuilder_RemovalOfDupes(t *testing.T) {
	s := newUnionTypeBuilderTest(t)
	stringAndObject := s.Registry.CreateUnionType(s.STRING_TYPE, s.OBJECT_TYPE)
	s.AssertUnion("(Object|string)", stringAndObject, s.STRING_OBJECT_TYPE)
	s.AssertUnion("(Object|string)", s.STRING_OBJECT_TYPE, stringAndObject)
}

func TestUnionTypeBuilder_RemovalOfDupes2(t *testing.T) {
	s := newUnionTypeBuilderTest(t)
	union :=
		s.Registry.CreateUnionType(
			s.sub,
			s.CreateFunctionWithReturn(s.base),
			s.base,
			s.CreateFunctionWithReturn(s.sub))
	s.AssertThat(union.String()).IsEqualTo("(Base|function(): Base)")
}

func TestUnionTypeBuilder_RemovalOfDupes3(t *testing.T) {
	s := newUnionTypeBuilderTest(t)
	union :=
		s.Registry.CreateUnionType(
			s.base,
			s.CreateFunctionWithReturn(s.sub),
			s.sub,
			s.CreateFunctionWithReturn(s.base))
	s.AssertThat(union.String()).IsEqualTo("(Base|function(): Base)")
}

func TestUnionTypeBuilder_RemovalOfIdenticalRecordTypes_beforeResolution(t *testing.T) {
	s := newUnionTypeBuilderTest(t)
	var recordUnion *UnionType
	closer := s.Registry.GetResolver().OpenForDefinition()
	defer closer.Close()
	{
		// Create two identical anonymous record types
		type_ := s.Registry.CreateRecordType(mapOf("prop", s.STRING_TYPE))

		recordUnion = s.Registry.CreateUnionType(type_, type_, s.NULL_TYPE).(*UnionType)

		// drop the identical references to 'type'
		s.AssertThat(recordUnion.GetAlternates()).HasSize(2)
	}
}

func TestUnionTypeBuilder_RemovalOfDuplicateRecordTypes_deferredUntilResolution(t *testing.T) {
	s := newUnionTypeBuilderTest(t)
	var recordUnion *UnionType
	func() {
		closer := s.Registry.GetResolver().OpenForDefinition()
		defer closer.Close()
		// Create two identical anonymous record types
		typeA := s.Registry.CreateRecordType(mapOf("prop", s.STRING_TYPE))
		typeB := s.Registry.CreateRecordType(mapOf("prop", s.STRING_TYPE))

		recordUnion = s.Registry.CreateUnionType(typeA, typeB).(*UnionType)

		s.AssertThat(recordUnion.GetAlternates()).HasSize(2)
	}()
	s.AssertThat(recordUnion.GetAlternates()).HasSize(1)
}

func TestUnionTypeBuilder_RemovalOfDuplicateRecordTypes1(t *testing.T) {
	s := newUnionTypeBuilderTest(t)
	builder := UnionType_Builder_(s.Registry)

	s.AddRecordType(builder, false)
	s.AddRecordType(builder, false)

	s.AssertThat(builder.Build().ToMaybeUnionType()).IsNull()
}

func TestUnionTypeBuilder_RemovalOfDuplicateRecordTypes2(t *testing.T) {
	s := newUnionTypeBuilderTest(t)
	builder := UnionType_Builder_(s.Registry)

	s.AddRecordType(builder, true)
	s.AddRecordType(builder, true)

	s.AssertThat(builder.Build().ToMaybeUnionType()).IsNull()
}

func TestUnionTypeBuilder_DifferentTemplateSpecializations_whenUnioned_doNotLeakRawType(t *testing.T) {
	s := newUnionTypeBuilderTest(t)
	// Given
	arrayOfString := s.Registry.CreateTemplatizedType(s.ARRAY_TYPE, s.STRING_TYPE)
	arrayOfNumber := s.Registry.CreateTemplatizedType(s.ARRAY_TYPE, s.NUMBER_TYPE)
	arrayOfUnknown := s.Registry.CreateTemplatizedType(s.ARRAY_TYPE, s.UNKNOWN_TYPE)

	builder :=
		UnionType_Builder_(s.Registry).AddAlternate(arrayOfString).AddAlternate(arrayOfNumber)

	// When
	result := builder.Build()

	// Then
	s.AssertThat(result.IsRawTypeOfTemplatizedType()).IsFalse()
	s.AssertType(result).IsEqualTo(arrayOfUnknown)
}

func TestUnionTypeBuilder_AfterBuild_cannotRebuild(t *testing.T) {
	s := newUnionTypeBuilderTest(t)
	builder := UnionType_Builder_(s.Registry)

	builder.Build()

	rhinotesting.AssertThrows[jutil.Throwable](s.T, func() { builder.Build() })
}

func TestUnionTypeBuilder_AfterBuild_cannotAdd(t *testing.T) {
	s := newUnionTypeBuilderTest(t)
	builder := UnionType_Builder_(s.Registry)

	builder.Build()

	rhinotesting.AssertThrows[jutil.Throwable](s.T, func() { builder.AddAlternate(s.NUMBER_TYPE) })
}

func (s *unionTypeBuilderTest) AddRecordType(builder *UnionType_Builder, inferred bool) {
	recBuilder := NewRecordTypeBuilder(s.Registry)
	recBuilder.SetSynthesized(inferred)
	recBuilder.AddProperty("prop", s.NUMBER_TYPE, nil)
	builder.AddAlternate(recBuilder.Build())
}

func (s *unionTypeBuilderTest) AssertUnion(expected string, types ...JSType) {
	builder := UnionType_Builder_(s.Registry)
	for _, type_ := range types {
		builder.AddAlternate(type_)
	}
	s.AssertThat(builder.Build().String()).IsEqualTo(expected)
}

func (s *unionTypeBuilderTest) CreateFunctionWithReturn(type_ JSType) FunctionType {
	return FunctionType_Builder_(s.Registry).
		WithParametersOfList(s.Registry.CreateParameters()).
		WithReturnType(type_).
		Build()
}
