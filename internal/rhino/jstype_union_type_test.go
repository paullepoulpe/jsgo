// Port of com/google/javascript/rhino/jstype/UnionTypeTest.java @ closure-compiler v20260415

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

	"github.com/paullepoulpe/jsgo/internal/base"
	. "github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/rhinotesting"
)

type unionTypeTest struct {
	*rhinotesting.BaseJSTypeTestCase
	EMPTY_SCOPE      *rhinotesting.MapBasedScope
	base             ObjectType
	sub1             ObjectType
	sub2             ObjectType
	sub3             ObjectType
	baseNameConflict ObjectType
}

func newUnionTypeTest(t *testing.T) *unionTypeTest {
	s := &unionTypeTest{BaseJSTypeTestCase: rhinotesting.NewBaseJSTypeTestCase(t)}
	s.EMPTY_SCOPE = rhinotesting.MapBasedScope_EmptyScope()
	s.SetUp()
	return s
}

func (s *unionTypeTest) SetUp() {

	closer := s.Registry.GetResolver().OpenForDefinition()
	defer closer.Close()
	{
		s.base =
			FunctionType_Builder_(s.Registry).
				ForConstructor().
				WithName(ptr("Base")).
				Build().
				GetInstanceType()
		s.baseNameConflict =
			FunctionType_Builder_(s.Registry).
				ForConstructor().
				WithName(ptr("Base")).
				Build().
				GetInstanceType()
		s.sub1 =
			FunctionType_Builder_(s.Registry).
				ForConstructor().
				WithName(ptr("Sub1")).
				WithPrototypeBasedOn(s.base).
				Build().
				GetInstanceType()
		s.sub2 =
			FunctionType_Builder_(s.Registry).
				ForConstructor().
				WithName(ptr("Sub2")).
				WithPrototypeBasedOn(s.base).
				Build().
				GetInstanceType()
		s.sub3 =
			FunctionType_Builder_(s.Registry).
				ForConstructor().
				WithName(ptr("Sub3")).
				WithPrototypeBasedOn(s.base).
				Build().
				GetInstanceType()
	}
}

/**
* Assert that a type can assign to itself.
 */
func (s *unionTypeTest) AssertTypeCanAssignToItself(type_ JSType) {
	s.AssertThat(type_.IsSubtypeOf(type_)).IsTrue()
}

/** Tests the behavior of variants type. */
func TestUnionType_UnionType(t *testing.T) {
	s := newUnionTypeTest(t)
	nullOrString :=
		s.CreateUnionType(s.NULL_TYPE, s.STRING_OBJECT_TYPE).(*UnionType)
	stringOrNull :=
		s.CreateUnionType(s.STRING_OBJECT_TYPE, s.NULL_TYPE).(*UnionType)

	s.AssertType(stringOrNull).IsEqualTo(nullOrString)
	s.AssertType(nullOrString).IsEqualTo(stringOrNull)

	s.AssertTypeCanAssignToItself(s.CreateUnionType(s.VOID_TYPE, s.NUMBER_TYPE))
	s.AssertTypeCanAssignToItself(
		s.CreateUnionType(s.NUMBER_TYPE, s.STRING_TYPE, s.OBJECT_TYPE))
	s.AssertTypeCanAssignToItself(s.CreateUnionType(s.NUMBER_TYPE, s.BOOLEAN_TYPE))
	s.AssertTypeCanAssignToItself(s.CreateUnionType(s.VOID_TYPE))

	// findPropertyType
	s.AssertType(nullOrString.FindPropertyType("length")).IsEqualTo(s.NUMBER_TYPE)
	s.AssertThat(nullOrString.FindPropertyType("lengthx")).IsNull()

	rhinotesting.Asserts_AssertResolvesToSame(s.T, nullOrString)
}

func TestUnionType_UnionOfNullAndUnresolvedNamedType(t *testing.T) {
	s := newUnionTypeTest(t)
	s.ErrorReporter.ExpectAllWarnings("Bad type annotation. Unknown type not.resolved.named.type")
	unresolvedNamedType :=
		s.Registry.CreateNamedType(s.EMPTY_SCOPE, "not.resolved.named.type", ptr(""), -1, -1)
	nullOrUnknown := s.CreateUnionType(s.NULL_TYPE, unresolvedNamedType).(*UnionType)
	s.AssertThat(nullOrUnknown.IsUnknownType()).IsTrue()
	s.AssertType(s.NULL_TYPE.GetLeastSupertype(nullOrUnknown)).IsEqualTo(nullOrUnknown)
	s.AssertType(nullOrUnknown.GetLeastSupertype(s.NULL_TYPE)).IsEqualTo(nullOrUnknown)
	s.AssertType(s.NULL_TYPE.GetGreatestSubtype(nullOrUnknown)).IsEqualTo(s.NULL_TYPE)
	s.AssertType(nullOrUnknown.GetGreatestSubtype(s.NULL_TYPE)).IsEqualTo(s.NULL_TYPE)

	s.AssertThat(s.NULL_TYPE.DiffersFrom(nullOrUnknown)).IsTrue()
	s.AssertThat(nullOrUnknown.DiffersFrom(s.NULL_TYPE)).IsTrue()
	s.AssertThat(nullOrUnknown.DiffersFrom(unresolvedNamedType)).IsFalse()

	s.AssertThat(s.NULL_TYPE.IsSubtypeOf(nullOrUnknown)).IsTrue()
	s.AssertThat(unresolvedNamedType.IsSubtype(nullOrUnknown)).IsTrue()
	s.AssertThat(nullOrUnknown.IsSubtype(s.NULL_TYPE)).IsTrue()

	s.AssertType(nullOrUnknown.RestrictByNotNullOrUndefined()).
		IsEqualTo(unresolvedNamedType)
}

/** Tests {@link JSType#getGreatestSubtype(JSType)} on union types. */
func TestUnionType_GreatestSubtypeUnionTypes1(t *testing.T) {
	s := newUnionTypeTest(t)
	s.AssertType(s.CreateNullableType(s.STRING_TYPE).GetGreatestSubtype(s.CreateNullableType(s.NUMBER_TYPE))).
		IsEqualTo(s.NULL_TYPE)
}

/** Tests {@link JSType#getGreatestSubtype(JSType)} on union types. */
func TestUnionType_GreatestSubtypeUnionTypes2(t *testing.T) {
	s := newUnionTypeTest(t)
	subUnion := s.CreateUnionType(s.sub1, s.sub2).(*UnionType)
	s.AssertType(subUnion.GetGreatestSubtype(s.base)).IsEqualTo(subUnion)
}

/** Tests {@link JSType#getGreatestSubtype(JSType)} on union types. */
func TestUnionType_GreatestSubtypeUnionTypes3(t *testing.T) {
	s := newUnionTypeTest(t)
	// (number,undefined,null)
	nullableOptionalNumber :=
		s.CreateUnionType(s.NULL_TYPE, s.VOID_TYPE, s.NUMBER_TYPE).(*UnionType)
	// (null,undefined)
	nullUndefined :=
		s.CreateUnionType(s.VOID_TYPE, s.NULL_TYPE).(*UnionType)
	s.AssertType(nullUndefined.GetGreatestSubtype(nullableOptionalNumber)).
		IsEqualTo(nullUndefined)
	s.AssertType(nullableOptionalNumber.GetGreatestSubtype(nullUndefined)).
		IsEqualTo(nullUndefined)
}

/** Tests {@link JSType#getGreatestSubtype(JSType)} on union types. */
func TestUnionType_GreatestSubtypeUnionTypes4(t *testing.T) {
	s := newUnionTypeTest(t)
	union := s.CreateUnionType(s.NULL_TYPE, s.sub1, s.sub2).(*UnionType)
	s.AssertType(union.GetGreatestSubtype(s.base)).IsEqualTo(s.CreateUnionType(s.sub1, s.sub2))
}

/** Tests {@link JSType#getGreatestSubtype(JSType)} on union types. */
func TestUnionType_GreatestSubtypeUnionTypes5(t *testing.T) {
	s := newUnionTypeTest(t)
	subUnion := s.CreateUnionType(s.sub1, s.sub2)
	s.AssertType(subUnion.GetGreatestSubtype(s.STRING_OBJECT_TYPE)).
		IsEqualTo(s.NO_OBJECT_TYPE)
}

/** Tests subtyping of union types. */
func TestUnionType_SubtypingUnionTypes(t *testing.T) {
	s := newUnionTypeTest(t)
	// subtypes
	s.AssertThat(s.BOOLEAN_TYPE.IsSubtypeOf(s.CreateUnionType(s.BOOLEAN_TYPE, s.STRING_TYPE))).IsTrue()
	s.AssertThat(
		s.CreateUnionType(s.BOOLEAN_TYPE, s.STRING_TYPE).
			IsSubtypeOf(s.CreateUnionType(s.BOOLEAN_TYPE, s.STRING_TYPE))).
		IsTrue()
	s.AssertThat(
		s.CreateUnionType(s.BOOLEAN_TYPE, s.STRING_TYPE).
			IsSubtypeOf(s.CreateUnionType(s.BOOLEAN_TYPE, s.STRING_TYPE, s.NULL_TYPE))).
		IsTrue()
	s.AssertThat(
		s.CreateUnionType(s.BOOLEAN_TYPE, s.STRING_TYPE).
			IsSubtypeOf(s.CreateUnionType(s.BOOLEAN_TYPE, s.STRING_TYPE, s.NULL_TYPE))).
		IsTrue()
	s.AssertThat(
		s.CreateUnionType(s.BOOLEAN_TYPE).
			IsSubtypeOf(s.CreateUnionType(s.BOOLEAN_TYPE, s.STRING_TYPE, s.NULL_TYPE))).
		IsTrue()
	s.AssertThat(
		s.CreateUnionType(s.STRING_TYPE).
			IsSubtypeOf(s.CreateUnionType(s.BOOLEAN_TYPE, s.STRING_TYPE, s.NULL_TYPE))).
		IsTrue()
	s.AssertThat(s.CreateUnionType(s.STRING_TYPE, s.NULL_TYPE).IsSubtypeOf(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.CreateUnionType(s.DATE_TYPE, s.REGEXP_TYPE).IsSubtypeOf(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.CreateUnionType(s.sub1, s.sub2).IsSubtypeOf(s.base)).IsTrue()
	s.AssertThat(s.CreateUnionType(s.sub1, s.sub2).IsSubtypeOf(s.OBJECT_TYPE)).IsTrue()

	// not subtypes
	s.AssertThat(s.CreateUnionType(s.STRING_TYPE, s.NULL_TYPE).IsSubtypeOf(s.NO_TYPE)).IsFalse()
	s.AssertThat(s.CreateUnionType(s.STRING_TYPE, s.NULL_TYPE).IsSubtypeOf(s.NO_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.CreateUnionType(s.NO_OBJECT_TYPE, s.NULL_TYPE).IsSubtypeOf(s.OBJECT_TYPE)).IsFalse()

	// defined unions
	s.AssertThat(s.NUMBER_TYPE.IsSubtypeOf(s.OBJECT_NUMBER_STRING)).IsTrue()
	s.AssertThat(s.OBJECT_TYPE.IsSubtypeOf(s.OBJECT_NUMBER_STRING)).IsTrue()
	s.AssertThat(s.STRING_TYPE.IsSubtypeOf(s.OBJECT_NUMBER_STRING)).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.IsSubtypeOf(s.OBJECT_NUMBER_STRING)).IsTrue()

	s.AssertThat(s.NUMBER_TYPE.IsSubtypeOf(s.NUMBER_STRING_BOOLEAN)).IsTrue()
	s.AssertThat(s.BOOLEAN_TYPE.IsSubtypeOf(s.NUMBER_STRING_BOOLEAN)).IsTrue()
	s.AssertThat(s.STRING_TYPE.IsSubtypeOf(s.NUMBER_STRING_BOOLEAN)).IsTrue()

	s.AssertThat(s.NUMBER_TYPE.IsSubtypeOf(s.OBJECT_NUMBER_STRING_BOOLEAN)).IsTrue()
	s.AssertThat(s.OBJECT_TYPE.IsSubtypeOf(s.OBJECT_NUMBER_STRING_BOOLEAN)).IsTrue()
	s.AssertThat(s.STRING_TYPE.IsSubtypeOf(s.OBJECT_NUMBER_STRING_BOOLEAN)).IsTrue()
	s.AssertThat(s.BOOLEAN_TYPE.IsSubtypeOf(s.OBJECT_NUMBER_STRING_BOOLEAN)).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.IsSubtypeOf(s.OBJECT_NUMBER_STRING_BOOLEAN)).IsTrue()
}

/**
* Tests that special union types can assign to other types. Unions containing the unknown type
* should be able to assign to any other type.
 */
func TestUnionType_SpecialUnionCanAssignTo(t *testing.T) {
	s := newUnionTypeTest(t)
	// autoboxing quirks
	numbers := s.CreateUnionType(s.NUMBER_TYPE, s.NUMBER_OBJECT_TYPE).(*UnionType)
	s.AssertThat(numbers.IsSubtype(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(numbers.IsSubtype(s.NUMBER_OBJECT_TYPE)).IsFalse()
	s.AssertThat(numbers.IsSubtype(s.sub1)).IsFalse()

	strings_ := s.CreateUnionType(s.STRING_OBJECT_TYPE, s.STRING_TYPE).(*UnionType)
	s.AssertThat(strings_.IsSubtype(s.STRING_TYPE)).IsFalse()
	s.AssertThat(strings_.IsSubtype(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(strings_.IsSubtype(s.DATE_TYPE)).IsFalse()

	booleans := s.CreateUnionType(s.BOOLEAN_OBJECT_TYPE, s.BOOLEAN_TYPE).(*UnionType)
	s.AssertThat(booleans.IsSubtype(s.BOOLEAN_TYPE)).IsFalse()
	s.AssertThat(booleans.IsSubtype(s.BOOLEAN_OBJECT_TYPE)).IsFalse()
	s.AssertThat(booleans.IsSubtype(s.REGEXP_TYPE)).IsFalse()

	// unknown quirks
	unknown := s.CreateUnionType(s.UNKNOWN_TYPE, s.DATE_TYPE)
	s.AssertThat(unknown.IsSubtypeOf(s.STRING_TYPE)).IsTrue()

	// all members need to be assignable to
	stringDate := s.CreateUnionType(s.STRING_OBJECT_TYPE, s.DATE_TYPE).(*UnionType)
	s.AssertThat(stringDate.IsSubtype(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(stringDate.IsSubtype(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(stringDate.IsSubtype(s.DATE_TYPE)).IsFalse()
}

/** Tests the factory method {@link JSTypeRegistry#createUnionType(JSType...)}. */
func TestUnionType_CreateUnionType(t *testing.T) {
	s := newUnionTypeTest(t)
	// number
	optNumber :=
		s.Registry.CreateUnionType(s.NUMBER_TYPE, s.DATE_TYPE).(*UnionType)
	s.AssertThat(optNumber.Contains(s.NUMBER_TYPE)).IsTrue()
	s.AssertThat(optNumber.Contains(s.DATE_TYPE)).IsTrue()

	// union
	optUnion :=
		s.Registry.CreateUnionType(s.REGEXP_TYPE,
			s.Registry.CreateUnionType(s.STRING_OBJECT_TYPE, s.DATE_TYPE)).(*UnionType)
	s.AssertThat(optUnion.Contains(s.DATE_TYPE)).IsTrue()
	s.AssertThat(optUnion.Contains(s.STRING_OBJECT_TYPE)).IsTrue()
	s.AssertThat(optUnion.Contains(s.REGEXP_TYPE)).IsTrue()
}

func TestUnionType_UnionWithUnknown(t *testing.T) {
	s := newUnionTypeTest(t)
	s.AssertThat(s.CreateUnionType(s.UNKNOWN_TYPE, s.NULL_TYPE).IsUnknownType()).IsTrue()
}

func TestUnionType_GetRestrictedUnion1(t *testing.T) {
	s := newUnionTypeTest(t)
	numStr := s.CreateUnionType(s.NUMBER_TYPE, s.STRING_TYPE).(*UnionType)
	s.AssertType(numStr.GetRestrictedUnion(s.NUMBER_TYPE)).IsEqualTo(s.STRING_TYPE)
}

func TestUnionType_GetRestrictedUnion2(t *testing.T) {
	s := newUnionTypeTest(t)
	numStr := s.CreateUnionType(s.NULL_TYPE, s.sub1, s.sub2).(*UnionType)
	s.AssertType(numStr.GetRestrictedUnion(s.base)).IsEqualTo(s.NULL_TYPE)
}

func TestUnionType_Equals(t *testing.T) {
	s := newUnionTypeTest(t)
	type_ := s.CreateUnionType(s.NUMBER_TYPE, s.STRING_TYPE).(*UnionType)
	s.AssertThat(type_.Equals(nil)).IsFalse()
	s.AssertThat(type_.Equals(type_)).IsTrue()
}

func TestUnionType_ProxyUnionType(t *testing.T) {
	s := newUnionTypeTest(t)
	stringOrNumber :=
		s.CreateUnionType(s.NUMBER_TYPE, s.STRING_TYPE).(*UnionType)
	stringOrBoolean :=
		s.CreateUnionType(s.BOOLEAN_TYPE, s.STRING_TYPE).(*UnionType)

	s.AssertThat(stringOrNumber.GetLeastSupertype(stringOrBoolean).String()).
		IsEqualTo("(boolean|number|string)")
	s.AssertThat(stringOrNumber.GetGreatestSubtype(stringOrBoolean).String()).IsEqualTo("string")
	s.AssertThat(stringOrNumber.TestForEquality(stringOrBoolean)).IsEqualTo(base.Tri_UNKNOWN)
	s.AssertThat(stringOrNumber.GetTypesUnderEquality(stringOrBoolean).TypeA.String()).
		IsEqualTo("(number|string)")
	s.AssertThat(stringOrNumber.GetTypesUnderShallowEquality(stringOrBoolean).TypeA.String()).
		IsEqualTo("string")
	s.AssertThat(stringOrNumber.GetTypesUnderInequality(stringOrBoolean).TypeA.String()).
		IsEqualTo("(number|string)")
	s.AssertThat(stringOrNumber.GetTypesUnderShallowInequality(stringOrBoolean).TypeA.String()).
		IsEqualTo("(number|string)")

	stringOrNumberProxy :=
		NewProxyObjectType(s.Registry, stringOrNumber)
	stringOrBooleanProxy :=
		NewProxyObjectType(s.Registry, stringOrBoolean)
	s.AssertThat(stringOrNumberProxy.GetLeastSupertype(stringOrBooleanProxy).String()).
		IsEqualTo("(boolean|number|string)")
	s.AssertThat(stringOrNumberProxy.GetGreatestSubtype(stringOrBooleanProxy).String()).
		IsEqualTo("string")
	s.AssertThat(stringOrNumberProxy.TestForEquality(stringOrBooleanProxy)).IsEqualTo(base.Tri_UNKNOWN)
	s.AssertThat(stringOrNumberProxy.GetTypesUnderEquality(stringOrBooleanProxy).TypeA.String()).
		IsEqualTo("(number|string)")
	s.AssertThat(
		stringOrNumberProxy.GetTypesUnderShallowEquality(stringOrBooleanProxy).TypeA.String()).
		IsEqualTo("string")
	s.AssertThat(stringOrNumberProxy.GetTypesUnderInequality(stringOrBooleanProxy).TypeA.String()).
		IsEqualTo("(number|string)")
	s.AssertThat(
		stringOrNumberProxy.
			GetTypesUnderShallowInequality(stringOrBooleanProxy).
			TypeA.
			String()).
		IsEqualTo("(number|string)")
}

func TestUnionType_CollapseUnion1(t *testing.T) {
	s := newUnionTypeTest(t)
	s.AssertThat(s.Registry.CreateUnionType(s.NUMBER_TYPE, s.STRING_TYPE).CollapseUnion().String()).
		IsEqualTo("*")
}

func TestUnionType_CollapseUnion2(t *testing.T) {
	s := newUnionTypeTest(t)
	s.AssertThat(s.Registry.CreateUnionType(s.UNKNOWN_TYPE, s.NUMBER_TYPE).CollapseUnion().String()).
		IsEqualTo("?")
	s.AssertThat(s.Registry.CreateUnionType(s.NUMBER_TYPE, s.UNKNOWN_TYPE).CollapseUnion().String()).
		IsEqualTo("?")
}

func TestUnionType_CollapseUnion3(t *testing.T) {
	s := newUnionTypeTest(t)
	s.AssertThat(s.Registry.CreateUnionType(s.ARRAY_TYPE, s.DATE_TYPE).CollapseUnion().String()).
		IsEqualTo("Object")
	s.AssertThat(s.Registry.CreateUnionType(s.ARRAY_TYPE, s.OBJECT_TYPE).CollapseUnion().String()).
		IsEqualTo("Object")
	s.AssertThat(s.Registry.CreateUnionType(s.base, s.sub1).CollapseUnion().String()).IsEqualTo("Base")
	s.AssertThat(s.Registry.CreateUnionType(s.sub1, s.sub2).CollapseUnion().String()).IsEqualTo("Base")
	s.AssertThat(s.Registry.CreateUnionType(s.sub1, s.sub2, s.sub3).CollapseUnion().String()).
		IsEqualTo("Base")
}

func TestUnionType_CollapseUnion4(t *testing.T) {
	s := newUnionTypeTest(t)
	s.AssertThat(s.Registry.CreateUnionType(s.OBJECT_TYPE, s.STRING_TYPE).CollapseUnion().String()).
		IsEqualTo("*")
	s.AssertThat(s.Registry.CreateUnionType(s.STRING_TYPE, s.OBJECT_TYPE).CollapseUnion().String()).
		IsEqualTo("*")
}

func TestUnionType_CollapseProxyUnion(t *testing.T) {
	s := newUnionTypeTest(t)
	// Make sure we don't unbox the proxy.
	type_ := NewProxyObjectType(s.Registry, s.OBJECT_TYPE)
	s.AssertThat(type_.CollapseUnion()).IsSameInstanceAs(type_)
}

func TestUnionType_ShallowEquality(t *testing.T) {
	s := newUnionTypeTest(t)
	s.AssertThat(
		s.Registry.
			CreateUnionType(s.ARRAY_TYPE, s.STRING_TYPE).
			CanTestForShallowEqualityWith(s.OBJECT_TYPE)).
		IsTrue()
}

func TestUnionType_UnionOfSingleType_equalToNonUnion(t *testing.T) {
	s := newUnionTypeTest(t)
	// Create a union of a single type. UnionTypes collapse duplicate elements by default, so we
	// need to put in a proxy type then resolve it to number.

	var numberUnion *UnionType
	func() {
		closer := s.Registry.GetResolver().OpenForDefinition()
		defer closer.Close()
		numberProxy :=
			NamedType_Builder_(s.Registry, "NumberProxy").
				SetResolutionKind(NamedType_ResolutionKind_NONE).
				SetReferencedType(s.NUMBER_TYPE).
				Build()

		numberUnion = s.Registry.CreateUnionType(s.NUMBER_TYPE, numberProxy).(*UnionType)
	}()
	// verify our setup is correct
	s.AssertThat(numberUnion.GetAlternates()).HasSize(1)

	s.AssertType(numberUnion).IsEqualTo(s.NUMBER_TYPE)
	s.AssertType(s.NUMBER_TYPE).IsEqualTo(numberUnion)
}

func TestUnionType_UnionOfResolvedAlternates_isNotRebuilt_duringResolution(t *testing.T) {
	s := newUnionTypeTest(t)
	// Given
	s.AssertType(s.NUMBER_TYPE).IsResolved()
	s.AssertType(s.STRING_TYPE).IsResolved()

	// When
	var union *UnionType
	var originalAlternates []JSType
	func() {
		closer := s.Registry.GetResolver().OpenForDefinition()
		defer closer.Close()
		union = s.Registry.CreateUnionType(s.NUMBER_TYPE, s.STRING_TYPE).(*UnionType)
		s.AssertType(union).IsUnresolved()

		originalAlternates = union.GetAlternates()
	}()
	s.AssertType(union).IsResolved()

	// Then
	s.AssertThat(originalAlternates).IsSameInstanceAs(union.GetAlternates())
}

func TestUnionType_UnionOfResolvedAlternates_isRebuilt_ifAlternatedAddedWhileUnresolved(t *testing.T) {
	s := newUnionTypeTest(t)
	// Given
	s.AssertType(s.NUMBER_TYPE).IsResolved()
	builder := UnionType_Builder_(s.Registry).AddAlternate(s.NUMBER_TYPE)

	// When
	func() {
		closer := s.Registry.GetResolver().OpenForDefinition()
		defer closer.Close()
		numberProxy :=
			NamedType_Builder_(s.Registry, "NumberProxy").
				SetResolutionKind(NamedType_ResolutionKind_NONE).
				SetReferencedType(s.NUMBER_TYPE).
				Build()
		s.AssertType(numberProxy).IsUnresolved()

		builder.AddAlternate(numberProxy)
	}()
	result := builder.Build()

	// Then
	s.AssertThat(result.GetUnionMembers()).HasSize(1)
	s.AssertType(result).IsEqualTo(s.NUMBER_TYPE)
}

func TestUnionType_ToStringNameConflict(t *testing.T) {
	s := newUnionTypeTest(t)
	s.AssertThat(s.CreateUnionType(s.base, s.baseNameConflict).String()).IsEqualTo("(Base|Base)")
}
