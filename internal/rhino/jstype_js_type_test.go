// Port of com/google/javascript/rhino/jstype/JSTypeTest.java @ closure-compiler v20260415

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
	"fmt"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	. "github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/rhinotesting"
)

type jsTypeTest struct {
	*rhinotesting.BaseJSTypeTestCase
	dateMethod                    FunctionType
	functionType                  FunctionType
	unresolvedNamedType           *NamedType
	googBar                       FunctionType
	googSubBar                    FunctionType
	googSubSubBar                 FunctionType
	googBarInst                   ObjectType
	googSubBarInst                ObjectType
	googSubSubBarInst             ObjectType
	namedGoogBar                  *NamedType
	subclassOfUnresolvedNamedType ObjectType
	subclassCtor                  FunctionType
	interfaceType                 FunctionType
	interfaceInstType             ObjectType
	subInterfaceType              FunctionType
	subInterfaceInstType          ObjectType
	recordType                    JSType
	enumType                      *EnumType
	elementsType                  *EnumElementType
	forwardDeclaredNamedType      *NamedType
	EMPTY_SCOPE                   StaticTypedScope
	types                         []JSType
}

func newJsTypeTest(t *testing.T) *jsTypeTest {
	s := &jsTypeTest{BaseJSTypeTestCase: rhinotesting.NewBaseJSTypeTestCase(t)}
	s.EMPTY_SCOPE = rhinotesting.MapBasedScope_EmptyScope()
	s.SetUp()
	return s
}

/**
* A non exhaustive list of representative types used to test simple properties that should hold
* for all types (such as the reflexivity of subtyping).
 */

func (s *jsTypeTest) SetUp() {
	s.ResetRegistryWithForwardDeclaredName("forwardDeclared")
	func() {
		closer := s.Registry.GetResolver().OpenForDefinition()
		defer closer.Close()
		googObject := s.Registry.CreateAnonymousObjectType(nil)
		scope := rhinotesting.NewMapBasedScope(mapOf("goog", googObject))

		builder := NewRecordTypeBuilder(s.Registry)
		builder.AddProperty("a", s.NUMBER_TYPE, nil)
		builder.AddProperty("b", s.STRING_TYPE, nil)
		s.recordType = builder.Build()

		s.enumType = EnumType_Builder_(s.Registry).SetName(ptr("Enum")).SetElementType(s.NUMBER_TYPE).Build()
		s.elementsType = s.enumType.GetElementsType()
		s.functionType = FunctionType_Builder_(s.Registry).WithReturnType(s.NUMBER_TYPE).Build()
		s.dateMethod =
			FunctionType_Builder_(s.Registry).
				WithParametersOfList(s.Registry.CreateParameters()).
				WithReturnType(s.NUMBER_TYPE).
				WithTypeOfThis(s.DATE_TYPE).
				Build()

		s.ErrorReporter.ExpectAllWarnings("Bad type annotation. Unknown type not.resolved.named.type")
		s.unresolvedNamedType = s.Registry.CreateNamedType(scope, "not.resolved.named.type", ptr(""), -1, -1)
		s.namedGoogBar = s.Registry.CreateNamedType(scope, "goog.Bar", nil, -1, -1)

		s.subclassCtor = FunctionType_Builder_(s.Registry).ForConstructor().Build()
		s.subclassCtor.SetPrototypeBasedOn(s.unresolvedNamedType)
		s.subclassOfUnresolvedNamedType = s.subclassCtor.GetInstanceType()

		s.interfaceType = FunctionType_Builder_(s.Registry).ForInterface().WithName(ptr("Interface")).Build()
		s.interfaceInstType = s.interfaceType.GetInstanceType()

		s.subInterfaceType =
			FunctionType_Builder_(s.Registry).ForInterface().WithName(ptr("SubInterfacce")).Build()
		s.subInterfaceType.SetExtendedInterfaces([]ObjectType{s.interfaceInstType})
		s.subInterfaceInstType = s.subInterfaceType.GetInstanceType()

		s.googBar = s.Registry.CreateConstructorType(ptr("goog.Bar"), nil, nil, nil, nil, false)
		s.googBar.GetPrototype().DefineDeclaredProperty("date", s.DATE_TYPE, nil)
		s.googBar.SetImplementedInterfaces([]ObjectType{s.interfaceInstType})
		s.googBarInst = s.googBar.GetInstanceType()

		s.googSubBar = s.Registry.CreateConstructorType(ptr("googSubBar"), nil, nil, nil, nil, false)
		s.googSubBar.SetPrototypeBasedOn(s.googBar.GetInstanceType())
		s.googSubBarInst = s.googSubBar.GetInstanceType()

		s.googSubSubBar =
			s.Registry.CreateConstructorType(ptr("googSubSubBar"), nil, nil, nil, nil, false)
		s.googSubSubBar.SetPrototypeBasedOn(s.googSubBar.GetInstanceType())
		s.googSubSubBarInst = s.googSubSubBar.GetInstanceType()

		googObject.DefineDeclaredProperty("Bar", s.googBar, nil)

		s.forwardDeclaredNamedType = s.Registry.CreateNamedType(scope, "forwardDeclared", ptr("source"), 1, 0)
	}()

	s.AssertThat(s.namedGoogBar.GetImplicitPrototype()).IsNotNull()

	s.types =
		[]JSType{
			s.NO_OBJECT_TYPE,
			NewNoResolvedType(s.Registry, "Foo", nil),
			s.NO_TYPE,
			s.BOOLEAN_OBJECT_TYPE,
			s.BOOLEAN_TYPE,
			s.STRING_OBJECT_TYPE,
			s.STRING_TYPE,
			s.SYMBOL_OBJECT_TYPE,
			s.SYMBOL_TYPE,
			s.VOID_TYPE,
			s.UNKNOWN_TYPE,
			s.NULL_TYPE,
			s.NUMBER_OBJECT_TYPE,
			s.NUMBER_TYPE,
			s.DATE_TYPE,
			s.dateMethod,
			s.functionType,
			s.unresolvedNamedType,
			s.googBar,
			s.googSubBar,
			s.googSubSubBar,
			s.namedGoogBar,
			s.googBar.GetInstanceType(),
			s.subclassOfUnresolvedNamedType,
			s.subclassCtor,
			s.recordType,
			s.enumType,
			s.elementsType,
			s.googBar,
			s.googSubBar,
			s.forwardDeclaredNamedType}
}

/** Tests the behavior of the top constructor type. */
func TestJSType_UniversalConstructorType(t *testing.T) {
	s := newJsTypeTest(t)
	// isXxx
	s.AssertThat(s.FUNCTION_TYPE.IsNoObjectType()).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsNoType()).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsArrayType()).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsBooleanValueType()).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsDateType()).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsEnumElementType()).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsNullType()).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsNamedType()).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsNullType()).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsNumber()).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsNumberObjectType()).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsNumberValueType()).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsObject()).IsTrue()
	s.AssertThat(s.FUNCTION_TYPE.IsFunctionPrototypeType()).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsRegexpType()).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsString()).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsStringObjectType()).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsStringValueType()).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsSymbol()).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsSymbolObjectType()).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsSymbolValueType()).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsEnumType()).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsUnionType()).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsStruct()).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsDict()).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsAllType()).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsVoidType()).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsConstructor()).IsTrue()
	s.AssertThat(s.FUNCTION_TYPE.IsInstanceType()).IsTrue()

	// isSubtype
	s.AssertThat(s.FUNCTION_TYPE.IsSubtype(s.NO_TYPE)).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsSubtype(s.NO_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsSubtype(s.ARRAY_TYPE)).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsSubtype(s.BOOLEAN_TYPE)).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsSubtype(s.BOOLEAN_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsSubtype(s.DATE_TYPE)).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsSubtype(s.functionType)).IsTrue()
	s.AssertThat(s.FUNCTION_TYPE.IsSubtype(s.recordType)).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsSubtype(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsSubtype(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsSubtype(s.NUMBER_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsSubtype(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.FUNCTION_TYPE.IsSubtype(s.REGEXP_TYPE)).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsSubtype(s.STRING_TYPE)).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsSubtype(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsSubtype(s.SYMBOL_TYPE)).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsSubtype(s.SYMBOL_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsSubtype(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.FUNCTION_TYPE.IsSubtype(s.VOID_TYPE)).IsFalse()

	// canTestForEqualityWith
	s.AssertThat(s.FUNCTION_TYPE.CanTestForEqualityWith(s.NO_TYPE)).IsTrue()
	s.AssertThat(s.FUNCTION_TYPE.CanTestForEqualityWith(s.NO_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.FUNCTION_TYPE.CanTestForEqualityWith(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.FUNCTION_TYPE.CanTestForEqualityWith(s.ARRAY_TYPE)).IsTrue()
	s.AssertThat(s.FUNCTION_TYPE.CanTestForEqualityWith(s.BOOLEAN_TYPE)).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.CanTestForEqualityWith(s.BOOLEAN_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.FUNCTION_TYPE.CanTestForEqualityWith(s.DATE_TYPE)).IsTrue()
	s.AssertThat(s.FUNCTION_TYPE.CanTestForEqualityWith(s.functionType)).IsTrue()
	s.AssertThat(s.FUNCTION_TYPE.CanTestForEqualityWith(s.recordType)).IsTrue()
	s.AssertThat(s.FUNCTION_TYPE.CanTestForEqualityWith(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.CanTestForEqualityWith(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.CanTestForEqualityWith(s.NUMBER_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.FUNCTION_TYPE.CanTestForEqualityWith(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.FUNCTION_TYPE.CanTestForEqualityWith(s.REGEXP_TYPE)).IsTrue()
	s.AssertThat(s.FUNCTION_TYPE.CanTestForEqualityWith(s.STRING_TYPE)).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.CanTestForEqualityWith(s.STRING_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.FUNCTION_TYPE.CanTestForEqualityWith(s.SYMBOL_TYPE)).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.CanTestForEqualityWith(s.SYMBOL_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.CanTestForEqualityWith(s.VOID_TYPE)).IsFalse()

	// canTestForShallowEqualityWith
	s.AssertThat(s.FUNCTION_TYPE.CanTestForShallowEqualityWith(s.NO_TYPE)).IsTrue()
	s.AssertThat(s.FUNCTION_TYPE.CanTestForShallowEqualityWith(s.NO_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.FUNCTION_TYPE.CanTestForShallowEqualityWith(s.ARRAY_TYPE)).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.CanTestForShallowEqualityWith(s.BOOLEAN_TYPE)).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.CanTestForShallowEqualityWith(s.BOOLEAN_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.CanTestForShallowEqualityWith(s.DATE_TYPE)).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.CanTestForShallowEqualityWith(s.functionType)).IsTrue()
	s.AssertThat(s.FUNCTION_TYPE.CanTestForShallowEqualityWith(s.recordType)).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.CanTestForShallowEqualityWith(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.CanTestForShallowEqualityWith(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.CanTestForShallowEqualityWith(s.NUMBER_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.CanTestForShallowEqualityWith(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.FUNCTION_TYPE.CanTestForShallowEqualityWith(s.REGEXP_TYPE)).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.CanTestForShallowEqualityWith(s.STRING_TYPE)).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.CanTestForShallowEqualityWith(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.CanTestForShallowEqualityWith(s.SYMBOL_TYPE)).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.CanTestForShallowEqualityWith(s.SYMBOL_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.CanTestForShallowEqualityWith(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.FUNCTION_TYPE.CanTestForShallowEqualityWith(s.VOID_TYPE)).IsFalse()

	// isNullable
	s.AssertThat(s.FUNCTION_TYPE.IsNullable()).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.IsVoidable()).IsFalse()

	// isObject
	s.AssertThat(s.FUNCTION_TYPE.IsObject()).IsTrue()

	// matchesXxx
	s.AssertThat(s.FUNCTION_TYPE.MatchesNumberContext()).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.MatchesObjectContext()).IsTrue()
	s.AssertThat(s.FUNCTION_TYPE.MatchesStringContext()).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.MatchesSymbolContext()).IsFalse()

	// toString
	s.AssertThat(s.FUNCTION_TYPE.String()).IsEqualTo("Function")
	s.AssertThat(s.FUNCTION_TYPE.HasDisplayName()).IsTrue()
	s.AssertThat(s.FUNCTION_TYPE.GetDisplayName()).IsEqualTo("Function")

	// getPropertyType
	s.AssertTypeEquals(s.UNKNOWN_TYPE, s.FUNCTION_TYPE.GetPropertyType("anyProperty"))

	s.AssertThat(s.FUNCTION_TYPE.IsNativeObjectType()).IsTrue()

	rhinotesting.Asserts_AssertResolvesToSame(s.T, s.FUNCTION_TYPE)

	s.AssertThat(s.FUNCTION_TYPE.IsNominalConstructorOrInterface()).IsTrue()
}

/** Tests the behavior of the Bottom Object type. */
func TestJSType_NoObjectType(t *testing.T) {
	s := newJsTypeTest(t)
	// isXxx
	s.AssertThat(s.NO_OBJECT_TYPE.IsNoObjectType()).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.IsNoType()).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.IsArrayType()).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.IsBooleanValueType()).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.IsDateType()).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.IsEnumElementType()).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.IsNullType()).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.IsNamedType()).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.IsNullType()).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.IsNumber()).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.IsNumberObjectType()).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.IsNumberValueType()).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.IsObject()).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.IsFunctionPrototypeType()).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.IsRegexpType()).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.IsString()).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.IsStringObjectType()).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.IsStringValueType()).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.IsSymbol()).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.IsSymbolObjectType()).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.IsSymbolValueType()).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.IsEnumType()).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.IsUnionType()).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.IsStruct()).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.IsDict()).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.IsAllType()).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.IsVoidType()).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.IsConstructor()).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.IsInstanceType()).IsFalse()

	// isSubtype
	s.AssertThat(s.NO_OBJECT_TYPE.IsSubtypeOf(s.NO_TYPE)).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.IsSubtypeOf(s.NO_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.IsSubtypeOf(s.ARRAY_TYPE)).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.IsSubtypeOf(s.BOOLEAN_TYPE)).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.IsSubtypeOf(s.BOOLEAN_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.IsSubtypeOf(s.DATE_TYPE)).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.IsSubtypeOf(s.functionType)).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.IsSubtypeOf(s.recordType)).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.IsSubtypeOf(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.IsSubtypeOf(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.IsSubtypeOf(s.NUMBER_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.IsSubtypeOf(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.IsSubtypeOf(s.REGEXP_TYPE)).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.IsSubtypeOf(s.STRING_TYPE)).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.IsSubtypeOf(s.STRING_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.IsSubtypeOf(s.SYMBOL_TYPE)).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.IsSubtypeOf(s.SYMBOL_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.IsSubtypeOf(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.IsSubtypeOf(s.VOID_TYPE)).IsFalse()

	// canTestForEqualityWith
	s.AssertCannotTestForEqualityWith(s.NO_OBJECT_TYPE, s.NO_TYPE)
	s.AssertCannotTestForEqualityWith(s.NO_OBJECT_TYPE, s.NO_OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.NO_OBJECT_TYPE, s.ALL_TYPE)
	s.AssertCanTestForEqualityWith(s.NO_OBJECT_TYPE, s.ARRAY_TYPE)
	s.AssertCanTestForEqualityWith(s.NO_OBJECT_TYPE, s.BOOLEAN_TYPE)
	s.AssertCanTestForEqualityWith(s.NO_OBJECT_TYPE, s.BOOLEAN_OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.NO_OBJECT_TYPE, s.DATE_TYPE)
	s.AssertCanTestForEqualityWith(s.NO_OBJECT_TYPE, s.functionType)
	s.AssertCanTestForEqualityWith(s.NO_OBJECT_TYPE, s.recordType)
	s.AssertCanTestForEqualityWith(s.NO_OBJECT_TYPE, s.NULL_TYPE)
	s.AssertCanTestForEqualityWith(s.NO_OBJECT_TYPE, s.NUMBER_TYPE)
	s.AssertCanTestForEqualityWith(s.NO_OBJECT_TYPE, s.NUMBER_OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.NO_OBJECT_TYPE, s.OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.NO_OBJECT_TYPE, s.REGEXP_TYPE)
	s.AssertCanTestForEqualityWith(s.NO_OBJECT_TYPE, s.STRING_TYPE)
	s.AssertCanTestForEqualityWith(s.NO_OBJECT_TYPE, s.STRING_OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.NO_OBJECT_TYPE, s.SYMBOL_TYPE)
	s.AssertCanTestForEqualityWith(s.NO_OBJECT_TYPE, s.SYMBOL_OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.NO_OBJECT_TYPE, s.VOID_TYPE)

	// canTestForShallowEqualityWith
	s.AssertThat(s.NO_OBJECT_TYPE.CanTestForShallowEqualityWith(s.NO_TYPE)).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.CanTestForShallowEqualityWith(s.NO_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.CanTestForShallowEqualityWith(s.ARRAY_TYPE)).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.CanTestForShallowEqualityWith(s.BOOLEAN_TYPE)).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.CanTestForShallowEqualityWith(s.BOOLEAN_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.CanTestForShallowEqualityWith(s.DATE_TYPE)).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.CanTestForShallowEqualityWith(s.functionType)).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.CanTestForShallowEqualityWith(s.recordType)).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.CanTestForShallowEqualityWith(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.CanTestForShallowEqualityWith(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.CanTestForShallowEqualityWith(s.NUMBER_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.CanTestForShallowEqualityWith(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.CanTestForShallowEqualityWith(s.REGEXP_TYPE)).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.CanTestForShallowEqualityWith(s.STRING_TYPE)).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.CanTestForShallowEqualityWith(s.STRING_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.CanTestForShallowEqualityWith(s.SYMBOL_TYPE)).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.CanTestForShallowEqualityWith(s.SYMBOL_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.CanTestForShallowEqualityWith(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.CanTestForShallowEqualityWith(s.VOID_TYPE)).IsFalse()

	// isNullable
	s.AssertThat(s.NO_OBJECT_TYPE.IsNullable()).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.IsVoidable()).IsFalse()

	// isObject
	s.AssertThat(s.NO_OBJECT_TYPE.IsObject()).IsTrue()

	// matchesXxx
	s.AssertThat(s.NO_OBJECT_TYPE.MatchesNumberContext()).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.MatchesObjectContext()).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.MatchesStringContext()).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.MatchesSymbolContext()).IsTrue()

	// toString
	s.AssertThat(s.NO_OBJECT_TYPE.String()).IsEqualTo("NoObject")
	s.AssertThat(s.NO_OBJECT_TYPE.HasDisplayName()).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.GetDisplayName()).IsNull()

	// getPropertyType
	s.AssertTypeEquals(s.NO_TYPE,
		s.NO_OBJECT_TYPE.GetPropertyType("anyProperty"))

	rhinotesting.Asserts_AssertResolvesToSame(s.T, s.NO_OBJECT_TYPE)

	s.AssertThat(s.NO_OBJECT_TYPE.IsNominalConstructorOrInterface()).IsFalse()
}

/** Tests the behavior of the Bottom type. */
func TestJSType_NoType(t *testing.T) {
	s := newJsTypeTest(t)
	// isXxx
	s.AssertThat(s.NO_TYPE.IsNoObjectType()).IsFalse()
	s.AssertThat(s.NO_TYPE.IsNoType()).IsTrue()
	s.AssertThat(s.NO_TYPE.IsArrayType()).IsFalse()
	s.AssertThat(s.NO_TYPE.IsBooleanValueType()).IsFalse()
	s.AssertThat(s.NO_TYPE.IsDateType()).IsFalse()
	s.AssertThat(s.NO_TYPE.IsEnumElementType()).IsFalse()
	s.AssertThat(s.NO_TYPE.IsNullType()).IsFalse()
	s.AssertThat(s.NO_TYPE.IsNamedType()).IsFalse()
	s.AssertThat(s.NO_TYPE.IsNullType()).IsFalse()
	s.AssertThat(s.NO_TYPE.IsNumber()).IsTrue()
	s.AssertThat(s.NO_TYPE.IsNumberObjectType()).IsFalse()
	s.AssertThat(s.NO_TYPE.IsNumberValueType()).IsFalse()
	s.AssertThat(s.NO_TYPE.IsObject()).IsTrue()
	s.AssertThat(s.NO_TYPE.IsFunctionPrototypeType()).IsFalse()
	s.AssertThat(s.NO_TYPE.IsRegexpType()).IsFalse()
	s.AssertThat(s.NO_TYPE.IsString()).IsTrue()
	s.AssertThat(s.NO_TYPE.IsStringObjectType()).IsFalse()
	s.AssertThat(s.NO_TYPE.IsStringValueType()).IsFalse()
	s.AssertThat(s.NO_TYPE.IsSymbol()).IsTrue()
	s.AssertThat(s.NO_TYPE.IsSymbolObjectType()).IsFalse()
	s.AssertThat(s.NO_TYPE.IsSymbolValueType()).IsFalse()
	s.AssertThat(s.NO_TYPE.IsEnumType()).IsFalse()
	s.AssertThat(s.NO_TYPE.IsUnionType()).IsFalse()
	s.AssertThat(s.NO_TYPE.IsStruct()).IsFalse()
	s.AssertThat(s.NO_TYPE.IsDict()).IsFalse()
	s.AssertThat(s.NO_TYPE.IsAllType()).IsFalse()
	s.AssertThat(s.NO_TYPE.IsVoidType()).IsFalse()
	s.AssertThat(s.NO_TYPE.IsConstructor()).IsFalse()
	s.AssertThat(s.NO_TYPE.IsInstanceType()).IsFalse()

	// isSubtype
	s.AssertThat(s.NO_TYPE.IsSubtypeOf(s.NO_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.IsSubtypeOf(s.NO_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.IsSubtypeOf(s.ARRAY_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.IsSubtypeOf(s.BOOLEAN_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.IsSubtypeOf(s.BOOLEAN_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.IsSubtypeOf(s.DATE_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.IsSubtypeOf(s.functionType)).IsTrue()
	s.AssertThat(s.NO_TYPE.IsSubtypeOf(s.NULL_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.IsSubtypeOf(s.NUMBER_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.IsSubtypeOf(s.NUMBER_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.IsSubtypeOf(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.IsSubtypeOf(s.REGEXP_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.IsSubtypeOf(s.STRING_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.IsSubtypeOf(s.STRING_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.IsSubtypeOf(s.SYMBOL_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.IsSubtypeOf(s.SYMBOL_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.IsSubtypeOf(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.IsSubtypeOf(s.VOID_TYPE)).IsTrue()

	// canTestForEqualityWith
	s.AssertCannotTestForEqualityWith(s.NO_TYPE, s.NO_TYPE)
	s.AssertCannotTestForEqualityWith(s.NO_TYPE, s.NO_OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.NO_TYPE, s.ARRAY_TYPE)
	s.AssertCanTestForEqualityWith(s.NO_TYPE, s.BOOLEAN_TYPE)
	s.AssertCanTestForEqualityWith(s.NO_TYPE, s.BOOLEAN_OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.NO_TYPE, s.DATE_TYPE)
	s.AssertCanTestForEqualityWith(s.NO_TYPE, s.functionType)
	s.AssertCanTestForEqualityWith(s.NO_TYPE, s.NULL_TYPE)
	s.AssertCanTestForEqualityWith(s.NO_TYPE, s.NUMBER_TYPE)
	s.AssertCanTestForEqualityWith(s.NO_TYPE, s.NUMBER_OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.NO_TYPE, s.OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.NO_TYPE, s.REGEXP_TYPE)
	s.AssertCanTestForEqualityWith(s.NO_TYPE, s.STRING_TYPE)
	s.AssertCanTestForEqualityWith(s.NO_TYPE, s.STRING_OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.NO_TYPE, s.SYMBOL_TYPE)
	s.AssertCanTestForEqualityWith(s.NO_TYPE, s.SYMBOL_OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.NO_TYPE, s.ALL_TYPE)
	s.AssertCanTestForEqualityWith(s.NO_TYPE, s.VOID_TYPE)

	// canTestForShallowEqualityWith
	s.AssertThat(s.NO_TYPE.CanTestForShallowEqualityWith(s.NO_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.CanTestForShallowEqualityWith(s.NO_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.CanTestForShallowEqualityWith(s.ARRAY_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.CanTestForShallowEqualityWith(s.BOOLEAN_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.CanTestForShallowEqualityWith(s.BOOLEAN_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.CanTestForShallowEqualityWith(s.DATE_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.CanTestForShallowEqualityWith(s.functionType)).IsTrue()
	s.AssertThat(s.NO_TYPE.CanTestForShallowEqualityWith(s.NULL_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.CanTestForShallowEqualityWith(s.NUMBER_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.CanTestForShallowEqualityWith(s.NUMBER_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.CanTestForShallowEqualityWith(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.CanTestForShallowEqualityWith(s.REGEXP_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.CanTestForShallowEqualityWith(s.STRING_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.CanTestForShallowEqualityWith(s.STRING_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.CanTestForShallowEqualityWith(s.SYMBOL_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.CanTestForShallowEqualityWith(s.SYMBOL_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.CanTestForShallowEqualityWith(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.CanTestForShallowEqualityWith(s.VOID_TYPE)).IsTrue()

	// isNullable
	s.AssertThat(s.NO_TYPE.IsNullable()).IsTrue()
	s.AssertThat(s.NO_TYPE.IsVoidable()).IsTrue()

	// isObject
	s.AssertThat(s.NO_TYPE.IsObject()).IsTrue()

	// matchesXxx
	s.AssertThat(s.NO_TYPE.MatchesNumberContext()).IsTrue()
	s.AssertThat(s.NO_TYPE.MatchesObjectContext()).IsTrue()
	s.AssertThat(s.NO_TYPE.MatchesStringContext()).IsTrue()
	s.AssertThat(s.NO_TYPE.MatchesSymbolContext()).IsTrue()

	// toString
	s.AssertThat(s.NO_TYPE.String()).IsEqualTo("None")
	s.AssertThat(s.NO_TYPE.GetDisplayName()).IsNull()
	s.AssertThat(s.NO_TYPE.HasDisplayName()).IsFalse()

	// getPropertyType
	s.AssertTypeEquals(s.NO_TYPE,
		s.NO_TYPE.GetPropertyType("anyProperty"))

	rhinotesting.Asserts_AssertResolvesToSame(s.T, s.NO_TYPE)

	s.AssertThat(s.NO_TYPE.IsNominalConstructorOrInterface()).IsFalse()
}

/** Tests the behavior of the unresolved Bottom type. */
func TestJSType_NoResolvedType(t *testing.T) {
	s := newJsTypeTest(t)
	// isXxx
	noResolvedFooType := NewNoResolvedType(s.Registry, "Foo", nil)
	s.AssertThat(noResolvedFooType.IsNoObjectType()).IsFalse()
	s.AssertThat(noResolvedFooType.IsNoType()).IsFalse()
	s.AssertThat(noResolvedFooType.IsNoResolvedType()).IsTrue()
	s.AssertThat(noResolvedFooType.IsArrayType()).IsFalse()
	s.AssertThat(noResolvedFooType.IsBooleanValueType()).IsFalse()
	s.AssertThat(noResolvedFooType.IsDateType()).IsFalse()
	s.AssertThat(noResolvedFooType.IsEnumElementType()).IsFalse()
	s.AssertThat(noResolvedFooType.IsNullType()).IsFalse()
	s.AssertThat(noResolvedFooType.IsNamedType()).IsFalse()
	s.AssertThat(noResolvedFooType.IsNumber()).IsFalse()
	s.AssertThat(noResolvedFooType.IsNumberObjectType()).IsFalse()
	s.AssertThat(noResolvedFooType.IsNumberValueType()).IsFalse()
	s.AssertThat(noResolvedFooType.IsObject()).IsTrue()
	s.AssertThat(noResolvedFooType.IsFunctionPrototypeType()).IsFalse()
	s.AssertThat(noResolvedFooType.IsRegexpType()).IsFalse()
	s.AssertThat(noResolvedFooType.IsString()).IsFalse()
	s.AssertThat(noResolvedFooType.IsStringObjectType()).IsFalse()
	s.AssertThat(noResolvedFooType.IsStringValueType()).IsFalse()
	s.AssertThat(noResolvedFooType.IsSymbol()).IsFalse()
	s.AssertThat(noResolvedFooType.IsSymbolObjectType()).IsFalse()
	s.AssertThat(noResolvedFooType.IsSymbolValueType()).IsFalse()
	s.AssertThat(noResolvedFooType.IsEnumType()).IsFalse()
	s.AssertThat(noResolvedFooType.IsUnionType()).IsFalse()
	s.AssertThat(noResolvedFooType.IsStruct()).IsFalse()
	s.AssertThat(noResolvedFooType.IsDict()).IsFalse()
	s.AssertThat(noResolvedFooType.IsAllType()).IsFalse()
	s.AssertThat(noResolvedFooType.IsVoidType()).IsFalse()
	s.AssertThat(noResolvedFooType.IsConstructor()).IsFalse()
	s.AssertThat(noResolvedFooType.IsInstanceType()).IsFalse()

	// isSubtype
	s.AssertType(noResolvedFooType).IsSubtypeOf(noResolvedFooType)
	s.AssertType(noResolvedFooType).IsSubtypeOf(s.ALL_TYPE)
	s.AssertType(noResolvedFooType).IsNotSubtypeOf(s.OBJECT_TYPE)
	s.AssertType(noResolvedFooType).IsNotSubtypeOf(s.NO_OBJECT_TYPE)
	s.AssertType(noResolvedFooType).IsNotSubtypeOf(s.ARRAY_TYPE)
	s.AssertType(noResolvedFooType).IsNotSubtypeOf(s.BOOLEAN_TYPE)
	s.AssertType(noResolvedFooType).IsNotSubtypeOf(s.BOOLEAN_OBJECT_TYPE)
	s.AssertType(noResolvedFooType).IsNotSubtypeOf(s.DATE_TYPE)
	s.AssertType(noResolvedFooType).IsNotSubtypeOf(s.functionType)
	s.AssertType(noResolvedFooType).IsNotSubtypeOf(s.NULL_TYPE)
	s.AssertType(noResolvedFooType).IsNotSubtypeOf(s.NUMBER_TYPE)
	s.AssertType(noResolvedFooType).IsNotSubtypeOf(s.NUMBER_OBJECT_TYPE)
	s.AssertType(noResolvedFooType).IsNotSubtypeOf(s.REGEXP_TYPE)
	s.AssertType(noResolvedFooType).IsNotSubtypeOf(s.STRING_TYPE)
	s.AssertType(noResolvedFooType).IsNotSubtypeOf(s.STRING_OBJECT_TYPE)
	s.AssertType(noResolvedFooType).IsNotSubtypeOf(s.SYMBOL_TYPE)
	s.AssertType(noResolvedFooType).IsNotSubtypeOf(s.SYMBOL_OBJECT_TYPE)
	s.AssertType(noResolvedFooType).IsNotSubtypeOf(s.VOID_TYPE)

	// canTestForEqualityWith
	s.AssertCanTestForEqualityWith(noResolvedFooType, noResolvedFooType)
	s.AssertCanTestForEqualityWith(noResolvedFooType, s.NO_TYPE)
	s.AssertCanTestForEqualityWith(noResolvedFooType, s.NO_OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(noResolvedFooType, s.ARRAY_TYPE)
	s.AssertCanTestForEqualityWith(noResolvedFooType, s.BOOLEAN_TYPE)
	s.AssertCanTestForEqualityWith(noResolvedFooType, s.BOOLEAN_OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(noResolvedFooType, s.DATE_TYPE)
	s.AssertCanTestForEqualityWith(noResolvedFooType, s.functionType)
	s.AssertCanTestForEqualityWith(noResolvedFooType, s.NULL_TYPE)
	s.AssertCanTestForEqualityWith(noResolvedFooType, s.NUMBER_TYPE)
	s.AssertCanTestForEqualityWith(noResolvedFooType, s.NUMBER_OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(noResolvedFooType, s.OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(noResolvedFooType, s.REGEXP_TYPE)
	s.AssertCanTestForEqualityWith(noResolvedFooType, s.STRING_TYPE)
	s.AssertCanTestForEqualityWith(noResolvedFooType, s.STRING_OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(noResolvedFooType, s.SYMBOL_TYPE)
	s.AssertCanTestForEqualityWith(noResolvedFooType, s.SYMBOL_OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(noResolvedFooType, s.ALL_TYPE)
	s.AssertCanTestForEqualityWith(noResolvedFooType, s.VOID_TYPE)

	// canTestForShallowEqualityWith
	s.AssertType(noResolvedFooType).CanTestForShallowEqualityWith(noResolvedFooType)
	s.AssertType(noResolvedFooType).CanTestForShallowEqualityWith(s.NO_OBJECT_TYPE)
	s.AssertType(noResolvedFooType).CanTestForShallowEqualityWith(s.ARRAY_TYPE)
	s.AssertType(noResolvedFooType).CanTestForShallowEqualityWith(s.BOOLEAN_TYPE)
	s.AssertType(noResolvedFooType).CanTestForShallowEqualityWith(s.BOOLEAN_OBJECT_TYPE)
	s.AssertType(noResolvedFooType).CanTestForShallowEqualityWith(s.DATE_TYPE)
	s.AssertType(noResolvedFooType).CanTestForShallowEqualityWith(s.functionType)
	s.AssertType(noResolvedFooType).CanTestForShallowEqualityWith(s.NULL_TYPE)
	s.AssertType(noResolvedFooType).CanTestForShallowEqualityWith(s.NUMBER_TYPE)
	s.AssertType(noResolvedFooType).CanTestForShallowEqualityWith(s.NUMBER_OBJECT_TYPE)
	s.AssertType(noResolvedFooType).CanTestForShallowEqualityWith(s.OBJECT_TYPE)
	s.AssertType(noResolvedFooType).CanTestForShallowEqualityWith(s.REGEXP_TYPE)
	s.AssertType(noResolvedFooType).CanTestForShallowEqualityWith(s.STRING_TYPE)
	s.AssertType(noResolvedFooType).CanTestForShallowEqualityWith(s.STRING_OBJECT_TYPE)
	s.AssertType(noResolvedFooType).CanTestForShallowEqualityWith(s.SYMBOL_TYPE)
	s.AssertType(noResolvedFooType).CanTestForShallowEqualityWith(s.SYMBOL_OBJECT_TYPE)
	s.AssertType(noResolvedFooType).CanTestForShallowEqualityWith(s.ALL_TYPE)
	s.AssertType(noResolvedFooType).CanTestForShallowEqualityWith(s.VOID_TYPE)

	// isNullable
	s.AssertThat(noResolvedFooType.IsNullable()).IsFalse()
	s.AssertThat(noResolvedFooType.IsVoidable()).IsFalse()

	// isObject
	s.AssertThat(noResolvedFooType.IsObject()).IsTrue()

	// matchesXxx
	s.AssertThat(noResolvedFooType.MatchesNumberContext()).IsFalse()
	s.AssertThat(noResolvedFooType.MatchesObjectContext()).IsFalse()
	s.AssertThat(noResolvedFooType.MatchesStringContext()).IsFalse()
	s.AssertThat(noResolvedFooType.MatchesSymbolContext()).IsFalse()

	// toString
	s.AssertThat(noResolvedFooType.String()).IsEqualTo("NoResolvedType<Foo>")
	s.AssertThat(noResolvedFooType.GetDisplayName()).IsEqualTo("Foo")
	s.AssertThat(noResolvedFooType.HasDisplayName()).IsTrue()

	rhinotesting.Asserts_AssertResolvesToSame(s.T, noResolvedFooType)

	s.AssertThat(s.forwardDeclaredNamedType.IsEmptyType()).IsTrue()
	s.AssertThat(s.forwardDeclaredNamedType.IsNoResolvedType()).IsTrue()

	nullable := s.Registry.CreateNullableType(noResolvedFooType).(*UnionType)
	s.AssertType(nullable.GetGreatestSubtype(s.NULL_TYPE)).IsEqualTo(s.NULL_TYPE)
	s.AssertType(nullable.GetRestrictedUnion(s.NULL_TYPE)).IsEqualTo(noResolvedFooType)

	testIsVoidable := s.Registry.CreateUnionType(noResolvedFooType, s.VOID_TYPE).(*UnionType)
	s.AssertThat(testIsVoidable.GetPossibleToBooleanOutcomes()).IsEqualTo(BooleanLiteralSet_BOTH)
}

/** Tests the behavior of the Array type. */
func TestJSType_ArrayType(t *testing.T) {
	s := newJsTypeTest(t)
	// isXxx
	s.AssertThat(s.ARRAY_TYPE.IsArrayType()).IsTrue()
	s.AssertThat(s.ARRAY_TYPE.IsBooleanValueType()).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsDateType()).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsEnumElementType()).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsNamedType()).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsNullType()).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsNumber()).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsNumberObjectType()).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsNumberValueType()).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsObject()).IsTrue()
	s.AssertThat(s.ARRAY_TYPE.IsFunctionPrototypeType()).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.GetImplicitPrototype().IsFunctionPrototypeType()).IsTrue()
	s.AssertThat(s.ARRAY_TYPE.IsRegexpType()).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsString()).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsStringObjectType()).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsStringValueType()).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsSymbol()).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsSymbolObjectType()).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsSymbolValueType()).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsEnumType()).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsUnionType()).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsStruct()).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsDict()).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsAllType()).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsVoidType()).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsConstructor()).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsInstanceType()).IsTrue()

	// isSubtype
	s.AssertThat(s.ARRAY_TYPE.IsSubtypeOf(s.NO_TYPE)).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsSubtypeOf(s.NO_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsSubtypeOf(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.ARRAY_TYPE.IsSubtypeOf(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsSubtypeOf(s.SYMBOL_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsSubtypeOf(s.SYMBOL_TYPE)).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsSubtypeOf(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsSubtypeOf(s.functionType)).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsSubtypeOf(s.recordType)).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsSubtypeOf(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsSubtypeOf(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.ARRAY_TYPE.IsSubtypeOf(s.DATE_TYPE)).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsSubtypeOf(s.unresolvedNamedType)).IsTrue()
	s.AssertThat(s.ARRAY_TYPE.IsSubtypeOf(s.namedGoogBar)).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsSubtypeOf(s.REGEXP_TYPE)).IsFalse()

	// canBeCalled
	s.AssertThat(s.ARRAY_TYPE.CanBeCalled()).IsFalse()

	// canTestForEqualityWith
	s.AssertCanTestForEqualityWith(s.ARRAY_TYPE, s.NO_TYPE)
	s.AssertCanTestForEqualityWith(s.ARRAY_TYPE, s.NO_OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.ARRAY_TYPE, s.ALL_TYPE)
	s.AssertCanTestForEqualityWith(s.ARRAY_TYPE, s.STRING_OBJECT_TYPE)
	s.AssertCannotTestForEqualityWith(s.ARRAY_TYPE, s.SYMBOL_TYPE)
	s.AssertCannotTestForEqualityWith(s.ARRAY_TYPE, s.SYMBOL_OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.ARRAY_TYPE, s.NUMBER_TYPE)
	s.AssertCanTestForEqualityWith(s.ARRAY_TYPE, s.functionType)
	s.AssertCanTestForEqualityWith(s.ARRAY_TYPE, s.recordType)
	s.AssertCannotTestForEqualityWith(s.ARRAY_TYPE, s.VOID_TYPE)
	s.AssertCanTestForEqualityWith(s.ARRAY_TYPE, s.OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.ARRAY_TYPE, s.DATE_TYPE)
	s.AssertCanTestForEqualityWith(s.ARRAY_TYPE, s.REGEXP_TYPE)

	// canTestForShallowEqualityWith
	s.AssertThat(s.ARRAY_TYPE.CanTestForShallowEqualityWith(s.NO_TYPE)).IsTrue()
	s.AssertThat(s.ARRAY_TYPE.CanTestForShallowEqualityWith(s.NO_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.ARRAY_TYPE.CanTestForShallowEqualityWith(s.ARRAY_TYPE)).IsTrue()
	s.AssertThat(s.ARRAY_TYPE.CanTestForShallowEqualityWith(s.BOOLEAN_TYPE)).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.CanTestForShallowEqualityWith(s.BOOLEAN_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.CanTestForShallowEqualityWith(s.DATE_TYPE)).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.CanTestForShallowEqualityWith(s.functionType)).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.CanTestForShallowEqualityWith(s.recordType)).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.CanTestForShallowEqualityWith(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.CanTestForShallowEqualityWith(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.CanTestForShallowEqualityWith(s.NUMBER_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.CanTestForShallowEqualityWith(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.ARRAY_TYPE.CanTestForShallowEqualityWith(s.REGEXP_TYPE)).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.CanTestForShallowEqualityWith(s.STRING_TYPE)).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.CanTestForShallowEqualityWith(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.CanTestForShallowEqualityWith(s.SYMBOL_TYPE)).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.CanTestForShallowEqualityWith(s.SYMBOL_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.CanTestForShallowEqualityWith(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.ARRAY_TYPE.CanTestForShallowEqualityWith(s.VOID_TYPE)).IsFalse()

	// isNullable
	s.AssertThat(s.ARRAY_TYPE.IsNullable()).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsVoidable()).IsFalse()
	s.AssertThat(s.CreateUnionType(s.ARRAY_TYPE, s.NULL_TYPE).IsNullable()).IsTrue()
	s.AssertThat(s.CreateUnionType(s.ARRAY_TYPE, s.VOID_TYPE).IsVoidable()).IsTrue()

	// isObject
	s.AssertThat(s.ARRAY_TYPE.IsObject()).IsTrue()

	// getLeastSupertype
	s.AssertTypeEquals(s.ALL_TYPE,
		s.ARRAY_TYPE.GetLeastSupertype(s.ALL_TYPE))
	s.AssertTypeEquals(s.CreateUnionType(s.STRING_OBJECT_TYPE, s.ARRAY_TYPE),
		s.ARRAY_TYPE.GetLeastSupertype(s.STRING_OBJECT_TYPE))
	s.AssertTypeEquals(
		s.CreateUnionType(s.SYMBOL_OBJECT_TYPE, s.ARRAY_TYPE),
		s.ARRAY_TYPE.GetLeastSupertype(s.SYMBOL_OBJECT_TYPE))
	s.AssertTypeEquals(s.CreateUnionType(s.NUMBER_TYPE, s.ARRAY_TYPE),
		s.ARRAY_TYPE.GetLeastSupertype(s.NUMBER_TYPE))
	s.AssertTypeEquals(s.CreateUnionType(s.ARRAY_TYPE, s.functionType),
		s.ARRAY_TYPE.GetLeastSupertype(s.functionType))
	s.AssertTypeEquals(s.OBJECT_TYPE, s.ARRAY_TYPE.GetLeastSupertype(s.OBJECT_TYPE))
	s.AssertTypeEquals(s.CreateUnionType(s.DATE_TYPE, s.ARRAY_TYPE),
		s.ARRAY_TYPE.GetLeastSupertype(s.DATE_TYPE))
	s.AssertTypeEquals(s.CreateUnionType(s.REGEXP_TYPE, s.ARRAY_TYPE),
		s.ARRAY_TYPE.GetLeastSupertype(s.REGEXP_TYPE))

	// getPropertyType
	s.AssertThat(s.ARRAY_TYPE.GetImplicitPrototype().GetPropertiesCount()).IsEqualTo(17)
	s.AssertThat(s.ARRAY_TYPE.GetPropertiesCount()).IsEqualTo(18)
	s.AssertReturnTypeEquals(s.ARRAY_TYPE, s.ARRAY_TYPE.GetPropertyType("constructor"))
	s.AssertReturnTypeEquals(s.STRING_TYPE,
		s.ARRAY_TYPE.GetPropertyType("toString"))
	s.AssertReturnTypeEquals(s.STRING_TYPE,
		s.ARRAY_TYPE.GetPropertyType("toLocaleString"))
	s.AssertReturnTypeEquals(s.ARRAY_TYPE, s.ARRAY_TYPE.GetPropertyType("concat"))
	s.AssertReturnTypeEquals(s.STRING_TYPE,
		s.ARRAY_TYPE.GetPropertyType("join"))
	s.AssertReturnTypeEquals(s.UNKNOWN_TYPE, s.ARRAY_TYPE.GetPropertyType("pop"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE, s.ARRAY_TYPE.GetPropertyType("push"))
	s.AssertReturnTypeEquals(s.ARRAY_TYPE, s.ARRAY_TYPE.GetPropertyType("reverse"))
	s.AssertReturnTypeEquals(s.UNKNOWN_TYPE, s.ARRAY_TYPE.GetPropertyType("shift"))
	s.AssertReturnTypeEquals(s.ARRAY_TYPE, s.ARRAY_TYPE.GetPropertyType("slice"))
	s.AssertReturnTypeEquals(s.ARRAY_TYPE, s.ARRAY_TYPE.GetPropertyType("sort"))
	s.AssertReturnTypeEquals(s.ARRAY_TYPE, s.ARRAY_TYPE.GetPropertyType("splice"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE, s.ARRAY_TYPE.GetPropertyType("unshift"))
	s.AssertTypeEquals(s.NUMBER_TYPE, s.ARRAY_TYPE.GetPropertyType("length"))

	// isPropertyType*
	s.AssertPropertyTypeDeclared(s.ARRAY_TYPE, "pop")

	// matchesXxx
	s.AssertThat(s.ARRAY_TYPE.MatchesNumberContext()).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.MatchesObjectContext()).IsTrue()
	s.AssertThat(s.ARRAY_TYPE.MatchesStringContext()).IsTrue()
	s.AssertThat(s.ARRAY_TYPE.MatchesSymbolContext()).IsFalse()

	// toString
	s.AssertThat(s.ARRAY_TYPE.String()).IsEqualTo("Array")
	s.AssertThat(s.ARRAY_TYPE.HasDisplayName()).IsTrue()
	s.AssertThat(s.ARRAY_TYPE.GetDisplayName()).IsEqualTo("Array")

	s.AssertThat(s.ARRAY_TYPE.IsNativeObjectType()).IsTrue()

	rhinotesting.Asserts_AssertResolvesToSame(s.T, s.ARRAY_TYPE)

	s.AssertThat(s.ARRAY_TYPE.IsNominalConstructorOrInterface()).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.GetConstructor().IsNominalConstructorOrInterface()).IsTrue()
}

/** Tests the behavior of the unknown type. */
func TestJSType_UnknownType(t *testing.T) {
	s := newJsTypeTest(t)
	// isXxx
	s.AssertThat(s.UNKNOWN_TYPE.IsArrayType()).IsFalse()
	s.AssertThat(s.UNKNOWN_TYPE.IsBooleanObjectType()).IsFalse()
	s.AssertThat(s.UNKNOWN_TYPE.IsBooleanValueType()).IsFalse()
	s.AssertThat(s.UNKNOWN_TYPE.IsDateType()).IsFalse()
	s.AssertThat(s.UNKNOWN_TYPE.IsEnumElementType()).IsFalse()
	s.AssertThat(s.UNKNOWN_TYPE.IsNamedType()).IsFalse()
	s.AssertThat(s.UNKNOWN_TYPE.IsNullType()).IsFalse()
	s.AssertThat(s.UNKNOWN_TYPE.IsNumberObjectType()).IsFalse()
	s.AssertThat(s.UNKNOWN_TYPE.IsNumberValueType()).IsFalse()
	s.AssertThat(s.UNKNOWN_TYPE.IsObject()).IsTrue()
	s.AssertThat(s.UNKNOWN_TYPE.IsFunctionPrototypeType()).IsFalse()
	s.AssertThat(s.UNKNOWN_TYPE.IsRegexpType()).IsFalse()
	s.AssertThat(s.UNKNOWN_TYPE.IsStringObjectType()).IsFalse()
	s.AssertThat(s.UNKNOWN_TYPE.IsStringValueType()).IsFalse()
	s.AssertThat(s.UNKNOWN_TYPE.IsSymbolObjectType()).IsFalse()
	s.AssertThat(s.UNKNOWN_TYPE.IsSymbolValueType()).IsFalse()
	s.AssertThat(s.UNKNOWN_TYPE.IsEnumType()).IsFalse()
	s.AssertThat(s.UNKNOWN_TYPE.IsUnionType()).IsFalse()
	s.AssertThat(s.UNKNOWN_TYPE.IsStruct()).IsFalse()
	s.AssertThat(s.UNKNOWN_TYPE.IsDict()).IsFalse()
	s.AssertThat(s.UNKNOWN_TYPE.IsUnknownType()).IsTrue()
	s.AssertThat(s.UNKNOWN_TYPE.IsVoidType()).IsFalse()
	s.AssertThat(s.UNKNOWN_TYPE.IsConstructor()).IsFalse()
	s.AssertThat(s.UNKNOWN_TYPE.IsInstanceType()).IsFalse()

	// autoboxesTo
	s.AssertThat(s.UNKNOWN_TYPE.AutoboxesTo()).IsNull()

	// isSubtype
	s.AssertThat(s.UNKNOWN_TYPE.IsSubtypeOf(s.UNKNOWN_TYPE)).IsTrue()
	s.AssertThat(s.UNKNOWN_TYPE.IsSubtypeOf(s.STRING_TYPE)).IsTrue()
	s.AssertThat(s.UNKNOWN_TYPE.IsSubtypeOf(s.SYMBOL_TYPE)).IsTrue()
	s.AssertThat(s.UNKNOWN_TYPE.IsSubtypeOf(s.NUMBER_TYPE)).IsTrue()
	s.AssertThat(s.UNKNOWN_TYPE.IsSubtypeOf(s.functionType)).IsTrue()
	s.AssertThat(s.UNKNOWN_TYPE.IsSubtypeOf(s.recordType)).IsTrue()
	s.AssertThat(s.UNKNOWN_TYPE.IsSubtypeOf(s.NULL_TYPE)).IsTrue()
	s.AssertThat(s.UNKNOWN_TYPE.IsSubtypeOf(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.UNKNOWN_TYPE.IsSubtypeOf(s.DATE_TYPE)).IsTrue()
	s.AssertThat(s.UNKNOWN_TYPE.IsSubtypeOf(s.namedGoogBar)).IsTrue()
	s.AssertThat(s.UNKNOWN_TYPE.IsSubtypeOf(s.unresolvedNamedType)).IsTrue()
	s.AssertThat(s.UNKNOWN_TYPE.IsSubtypeOf(s.REGEXP_TYPE)).IsTrue()
	s.AssertThat(s.UNKNOWN_TYPE.IsSubtypeOf(s.VOID_TYPE)).IsTrue()

	// canBeCalled
	s.AssertThat(s.UNKNOWN_TYPE.CanBeCalled()).IsTrue()

	// canTestForEqualityWith
	s.AssertCanTestForEqualityWith(s.UNKNOWN_TYPE, s.UNKNOWN_TYPE)
	s.AssertCanTestForEqualityWith(s.UNKNOWN_TYPE, s.STRING_TYPE)
	s.AssertCanTestForEqualityWith(s.UNKNOWN_TYPE, s.SYMBOL_TYPE)
	s.AssertCanTestForEqualityWith(s.UNKNOWN_TYPE, s.NUMBER_TYPE)
	s.AssertCanTestForEqualityWith(s.UNKNOWN_TYPE, s.functionType)
	s.AssertCanTestForEqualityWith(s.UNKNOWN_TYPE, s.recordType)
	s.AssertCanTestForEqualityWith(s.UNKNOWN_TYPE, s.VOID_TYPE)
	s.AssertCanTestForEqualityWith(s.UNKNOWN_TYPE, s.OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.UNKNOWN_TYPE, s.DATE_TYPE)
	s.AssertCanTestForEqualityWith(s.UNKNOWN_TYPE, s.REGEXP_TYPE)
	s.AssertCanTestForEqualityWith(s.UNKNOWN_TYPE, s.BOOLEAN_TYPE)

	// canTestForShallowEqualityWith
	s.AssertThat(s.UNKNOWN_TYPE.CanTestForShallowEqualityWith(s.UNKNOWN_TYPE)).IsTrue()
	s.AssertThat(s.UNKNOWN_TYPE.CanTestForShallowEqualityWith(s.STRING_TYPE)).IsTrue()
	s.AssertThat(s.UNKNOWN_TYPE.CanTestForShallowEqualityWith(s.SYMBOL_TYPE)).IsTrue()
	s.AssertThat(s.UNKNOWN_TYPE.CanTestForShallowEqualityWith(s.NUMBER_TYPE)).IsTrue()
	s.AssertThat(s.UNKNOWN_TYPE.CanTestForShallowEqualityWith(s.functionType)).IsTrue()
	s.AssertThat(s.UNKNOWN_TYPE.CanTestForShallowEqualityWith(s.recordType)).IsTrue()
	s.AssertThat(s.UNKNOWN_TYPE.CanTestForShallowEqualityWith(s.VOID_TYPE)).IsTrue()
	s.AssertThat(s.UNKNOWN_TYPE.CanTestForShallowEqualityWith(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.UNKNOWN_TYPE.CanTestForShallowEqualityWith(s.DATE_TYPE)).IsTrue()
	s.AssertThat(s.UNKNOWN_TYPE.CanTestForShallowEqualityWith(s.REGEXP_TYPE)).IsTrue()

	// canHaveNullValue
	s.AssertThat(s.UNKNOWN_TYPE.IsNullable()).IsTrue()
	s.AssertThat(s.UNKNOWN_TYPE.IsVoidable()).IsTrue()

	// getGreatestCommonType
	s.AssertTypeEquals(s.UNKNOWN_TYPE,
		s.UNKNOWN_TYPE.GetLeastSupertype(s.UNKNOWN_TYPE))
	s.AssertTypeEquals(s.UNKNOWN_TYPE,
		s.UNKNOWN_TYPE.GetLeastSupertype(s.STRING_TYPE))
	s.AssertTypeEquals(s.UNKNOWN_TYPE, s.UNKNOWN_TYPE.GetLeastSupertype(s.SYMBOL_TYPE))
	s.AssertTypeEquals(s.UNKNOWN_TYPE, s.UNKNOWN_TYPE.GetLeastSupertype(s.NUMBER_TYPE))
	s.AssertTypeEquals(s.UNKNOWN_TYPE,
		s.UNKNOWN_TYPE.GetLeastSupertype(s.functionType))
	s.AssertTypeEquals(s.UNKNOWN_TYPE,
		s.UNKNOWN_TYPE.GetLeastSupertype(s.OBJECT_TYPE))
	s.AssertTypeEquals(s.UNKNOWN_TYPE,
		s.UNKNOWN_TYPE.GetLeastSupertype(s.DATE_TYPE))
	s.AssertTypeEquals(s.UNKNOWN_TYPE,
		s.UNKNOWN_TYPE.GetLeastSupertype(s.REGEXP_TYPE))

	// matchesXxx
	s.AssertThat(s.UNKNOWN_TYPE.MatchesNumberContext()).IsTrue()
	s.AssertThat(s.UNKNOWN_TYPE.MatchesObjectContext()).IsTrue()
	s.AssertThat(s.UNKNOWN_TYPE.MatchesStringContext()).IsTrue()
	s.AssertThat(s.UNKNOWN_TYPE.MatchesSymbolContext()).IsTrue()

	// isPropertyType*
	s.AssertPropertyTypeUnknown(s.UNKNOWN_TYPE, "XXX")

	// toString
	s.AssertThat(s.UNKNOWN_TYPE.String()).IsEqualTo("?")
	s.AssertThat(s.UNKNOWN_TYPE.HasDisplayName()).IsTrue()
	s.AssertThat(s.UNKNOWN_TYPE.GetDisplayName()).IsEqualTo("Unknown")

	rhinotesting.Asserts_AssertResolvesToSame(s.T, s.UNKNOWN_TYPE)
	s.AssertThat(s.UNKNOWN_TYPE.IsNominalConstructorOrInterface()).IsFalse()

	s.AssertThat(s.UNKNOWN_TYPE.GetPropertyType("abc")).IsEqualTo(s.UNKNOWN_TYPE)
}

/** Tests the behavior of the checked unknown type. */
func TestJSType_CheckedUnknownType(t *testing.T) {
	s := newJsTypeTest(t)
	// isPropertyType*
	s.AssertPropertyTypeUnknown(s.CHECKED_UNKNOWN_TYPE, "XXX")

	// toString
	s.AssertThat(s.CHECKED_UNKNOWN_TYPE.String()).IsEqualTo("??")
	s.AssertThat(s.CHECKED_UNKNOWN_TYPE.HasDisplayName()).IsTrue()
	s.AssertThat(s.CHECKED_UNKNOWN_TYPE.GetDisplayName()).IsEqualTo("Unknown")

	rhinotesting.Asserts_AssertResolvesToSame(s.T, s.CHECKED_UNKNOWN_TYPE)
	s.AssertThat(s.CHECKED_UNKNOWN_TYPE.IsNominalConstructorOrInterface()).IsFalse()

	s.AssertThat(s.CHECKED_UNKNOWN_TYPE.GetPropertyType("abc")).IsEqualTo(s.CHECKED_UNKNOWN_TYPE)
}

/** Tests the behavior of the unknown type. */
func TestJSType_AllType(t *testing.T) {
	s := newJsTypeTest(t)
	// isXxx
	s.AssertThat(s.ALL_TYPE.IsArrayType()).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsBooleanValueType()).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsDateType()).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsEnumElementType()).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsNamedType()).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsNullType()).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsNumber()).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsNumberObjectType()).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsNumberValueType()).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsObject()).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsFunctionPrototypeType()).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsRegexpType()).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsString()).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsStringObjectType()).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsStringValueType()).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsSymbol()).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsSymbolObjectType()).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsSymbolValueType()).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsEnumType()).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsUnionType()).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsStruct()).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsDict()).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsAllType()).IsTrue()
	s.AssertThat(s.ALL_TYPE.IsVoidType()).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsConstructor()).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsInstanceType()).IsFalse()

	// isSubtype
	s.AssertThat(s.ALL_TYPE.IsSubtypeOf(s.NO_TYPE)).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsSubtypeOf(s.NO_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsSubtypeOf(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.ALL_TYPE.IsSubtypeOf(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsSubtypeOf(s.SYMBOL_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsSubtypeOf(s.SYMBOL_TYPE)).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsSubtypeOf(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsSubtypeOf(s.functionType)).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsSubtypeOf(s.recordType)).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsSubtypeOf(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsSubtypeOf(s.OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsSubtypeOf(s.DATE_TYPE)).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsSubtypeOf(s.unresolvedNamedType)).IsTrue()
	s.AssertThat(s.ALL_TYPE.IsSubtypeOf(s.namedGoogBar)).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsSubtypeOf(s.REGEXP_TYPE)).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsSubtypeOf(s.VOID_TYPE)).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsSubtypeOf(s.UNKNOWN_TYPE)).IsTrue()

	// canBeCalled
	s.AssertThat(s.ALL_TYPE.CanBeCalled()).IsFalse()

	// canTestForEqualityWith
	s.AssertCanTestForEqualityWith(s.ALL_TYPE, s.ALL_TYPE)
	s.AssertCanTestForEqualityWith(s.ALL_TYPE, s.STRING_OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.ALL_TYPE, s.SYMBOL_OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.ALL_TYPE, s.SYMBOL_TYPE)
	s.AssertCanTestForEqualityWith(s.ALL_TYPE, s.NUMBER_TYPE)
	s.AssertCanTestForEqualityWith(s.ALL_TYPE, s.functionType)
	s.AssertCanTestForEqualityWith(s.ALL_TYPE, s.recordType)
	s.AssertCanTestForEqualityWith(s.ALL_TYPE, s.VOID_TYPE)
	s.AssertCanTestForEqualityWith(s.ALL_TYPE, s.OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.ALL_TYPE, s.DATE_TYPE)
	s.AssertCanTestForEqualityWith(s.ALL_TYPE, s.REGEXP_TYPE)

	// canTestForShallowEqualityWith
	s.AssertThat(s.ALL_TYPE.CanTestForShallowEqualityWith(s.NO_TYPE)).IsTrue()
	s.AssertThat(s.ALL_TYPE.CanTestForShallowEqualityWith(s.NO_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.ALL_TYPE.CanTestForShallowEqualityWith(s.ARRAY_TYPE)).IsTrue()
	s.AssertThat(s.ALL_TYPE.CanTestForShallowEqualityWith(s.BOOLEAN_TYPE)).IsTrue()
	s.AssertThat(s.ALL_TYPE.CanTestForShallowEqualityWith(s.BOOLEAN_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.ALL_TYPE.CanTestForShallowEqualityWith(s.DATE_TYPE)).IsTrue()
	s.AssertThat(s.ALL_TYPE.CanTestForShallowEqualityWith(s.functionType)).IsTrue()
	s.AssertThat(s.ALL_TYPE.CanTestForShallowEqualityWith(s.recordType)).IsTrue()
	s.AssertThat(s.ALL_TYPE.CanTestForShallowEqualityWith(s.NULL_TYPE)).IsTrue()
	s.AssertThat(s.ALL_TYPE.CanTestForShallowEqualityWith(s.NUMBER_TYPE)).IsTrue()
	s.AssertThat(s.ALL_TYPE.CanTestForShallowEqualityWith(s.NUMBER_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.ALL_TYPE.CanTestForShallowEqualityWith(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.ALL_TYPE.CanTestForShallowEqualityWith(s.REGEXP_TYPE)).IsTrue()
	s.AssertThat(s.ALL_TYPE.CanTestForShallowEqualityWith(s.STRING_TYPE)).IsTrue()
	s.AssertThat(s.ALL_TYPE.CanTestForShallowEqualityWith(s.STRING_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.ALL_TYPE.CanTestForShallowEqualityWith(s.SYMBOL_TYPE)).IsTrue()
	s.AssertThat(s.ALL_TYPE.CanTestForShallowEqualityWith(s.SYMBOL_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.ALL_TYPE.CanTestForShallowEqualityWith(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.ALL_TYPE.CanTestForShallowEqualityWith(s.VOID_TYPE)).IsTrue()

	// isNullable
	s.AssertThat(s.ALL_TYPE.IsNullable()).IsTrue()
	s.AssertThat(s.ALL_TYPE.IsVoidable()).IsTrue()

	// getLeastSupertype
	s.AssertTypeEquals(s.ALL_TYPE,
		s.ALL_TYPE.GetLeastSupertype(s.ALL_TYPE))
	s.AssertTypeEquals(s.ALL_TYPE,
		s.ALL_TYPE.GetLeastSupertype(s.UNKNOWN_TYPE))
	s.AssertTypeEquals(s.ALL_TYPE,
		s.ALL_TYPE.GetLeastSupertype(s.STRING_OBJECT_TYPE))
	s.AssertTypeEquals(s.ALL_TYPE, s.ALL_TYPE.GetLeastSupertype(s.SYMBOL_TYPE))
	s.AssertTypeEquals(s.ALL_TYPE, s.ALL_TYPE.GetLeastSupertype(s.SYMBOL_OBJECT_TYPE))
	s.AssertTypeEquals(s.ALL_TYPE, s.ALL_TYPE.GetLeastSupertype(s.NUMBER_TYPE))
	s.AssertTypeEquals(s.ALL_TYPE,
		s.ALL_TYPE.GetLeastSupertype(s.functionType))
	s.AssertTypeEquals(s.ALL_TYPE,
		s.ALL_TYPE.GetLeastSupertype(s.OBJECT_TYPE))
	s.AssertTypeEquals(s.ALL_TYPE,
		s.ALL_TYPE.GetLeastSupertype(s.DATE_TYPE))
	s.AssertTypeEquals(s.ALL_TYPE,
		s.ALL_TYPE.GetLeastSupertype(s.REGEXP_TYPE))

	// matchesXxx
	s.AssertThat(s.ALL_TYPE.MatchesNumberContext()).IsFalse() // ?
	s.AssertThat(s.ALL_TYPE.MatchesObjectContext()).IsTrue()
	s.AssertThat(s.ALL_TYPE.MatchesStringContext()).IsTrue()
	s.AssertThat(s.ALL_TYPE.MatchesSymbolContext()).IsFalse() // ?

	// toString
	s.AssertThat(s.ALL_TYPE.String()).IsEqualTo("*")

	s.AssertThat(s.ALL_TYPE.HasDisplayName()).IsTrue()
	s.AssertThat(s.ALL_TYPE.GetDisplayName()).IsEqualTo("<Any Type>")

	rhinotesting.Asserts_AssertResolvesToSame(s.T, s.ALL_TYPE)
	s.AssertThat(s.ALL_TYPE.IsNominalConstructorOrInterface()).IsFalse()
}

/** Tests the behavior of the Object type (the object at the top of the JavaScript hierarchy). */
func TestJSType_TheObjectType(t *testing.T) {
	s := newJsTypeTest(t)
	// implicit prototype
	s.AssertTypeEquals(s.OBJECT_PROTOTYPE, s.OBJECT_TYPE.GetImplicitPrototype())

	// isXxx
	s.AssertThat(s.OBJECT_TYPE.IsNoObjectType()).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.IsNoType()).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.IsArrayType()).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.IsBooleanValueType()).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.IsDateType()).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.IsEnumElementType()).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.IsNullType()).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.IsNamedType()).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.IsNullType()).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.IsNumber()).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.IsNumberObjectType()).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.IsNumberValueType()).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.IsObject()).IsTrue()
	s.AssertThat(s.OBJECT_TYPE.IsFunctionPrototypeType()).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.GetImplicitPrototype().IsFunctionPrototypeType()).IsTrue()
	s.AssertThat(s.OBJECT_TYPE.IsRegexpType()).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.IsString()).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.IsStringObjectType()).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.IsStringValueType()).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.IsSymbol()).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.IsSymbolObjectType()).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.IsSymbolValueType()).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.IsEnumType()).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.IsUnionType()).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.IsStruct()).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.IsDict()).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.IsAllType()).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.IsVoidType()).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.IsConstructor()).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.IsInstanceType()).IsTrue()

	// isSubtype
	s.AssertThat(s.OBJECT_TYPE.IsSubtypeOf(s.NO_TYPE)).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.IsSubtypeOf(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.OBJECT_TYPE.IsSubtypeOf(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.IsSubtypeOf(s.SYMBOL_TYPE)).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.IsSubtypeOf(s.SYMBOL_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.IsSubtypeOf(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.IsSubtypeOf(s.functionType)).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.IsSubtypeOf(s.recordType)).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.IsSubtypeOf(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.IsSubtypeOf(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.OBJECT_TYPE.IsSubtypeOf(s.DATE_TYPE)).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.IsSubtypeOf(s.namedGoogBar)).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.IsSubtypeOf(s.unresolvedNamedType)).IsTrue()
	s.AssertThat(s.OBJECT_TYPE.IsSubtypeOf(s.REGEXP_TYPE)).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.IsSubtypeOf(s.ARRAY_TYPE)).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.IsSubtypeOf(s.UNKNOWN_TYPE)).IsTrue()

	// canBeCalled
	s.AssertThat(s.OBJECT_TYPE.CanBeCalled()).IsFalse()

	// canTestForEqualityWith
	s.AssertCanTestForEqualityWith(s.OBJECT_TYPE, s.ALL_TYPE)
	s.AssertCanTestForEqualityWith(s.OBJECT_TYPE, s.STRING_OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.OBJECT_TYPE, s.SYMBOL_OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.OBJECT_TYPE, s.SYMBOL_TYPE)
	s.AssertCanTestForEqualityWith(s.OBJECT_TYPE, s.NUMBER_TYPE)
	s.AssertCanTestForEqualityWith(s.OBJECT_TYPE, s.STRING_TYPE)
	s.AssertCanTestForEqualityWith(s.OBJECT_TYPE, s.BOOLEAN_TYPE)
	s.AssertCanTestForEqualityWith(s.OBJECT_TYPE, s.functionType)
	s.AssertCanTestForEqualityWith(s.OBJECT_TYPE, s.recordType)
	s.AssertCannotTestForEqualityWith(s.OBJECT_TYPE, s.VOID_TYPE)
	s.AssertCanTestForEqualityWith(s.OBJECT_TYPE, s.OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.OBJECT_TYPE, s.DATE_TYPE)
	s.AssertCanTestForEqualityWith(s.OBJECT_TYPE, s.REGEXP_TYPE)
	s.AssertCanTestForEqualityWith(s.OBJECT_TYPE, s.ARRAY_TYPE)
	s.AssertCanTestForEqualityWith(s.OBJECT_TYPE, s.UNKNOWN_TYPE)

	// canTestForShallowEqualityWith
	s.AssertThat(s.OBJECT_TYPE.CanTestForShallowEqualityWith(s.NO_TYPE)).IsTrue()
	s.AssertThat(s.OBJECT_TYPE.CanTestForShallowEqualityWith(s.NO_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.OBJECT_TYPE.CanTestForShallowEqualityWith(s.ARRAY_TYPE)).IsTrue()
	s.AssertThat(s.OBJECT_TYPE.CanTestForShallowEqualityWith(s.BOOLEAN_TYPE)).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.CanTestForShallowEqualityWith(s.BOOLEAN_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.OBJECT_TYPE.CanTestForShallowEqualityWith(s.DATE_TYPE)).IsTrue()
	s.AssertThat(s.OBJECT_TYPE.CanTestForShallowEqualityWith(s.functionType)).IsTrue()
	s.AssertThat(s.OBJECT_TYPE.CanTestForShallowEqualityWith(s.recordType)).IsTrue()
	s.AssertThat(s.OBJECT_TYPE.CanTestForShallowEqualityWith(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.CanTestForShallowEqualityWith(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.CanTestForShallowEqualityWith(s.NUMBER_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.OBJECT_TYPE.CanTestForShallowEqualityWith(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.OBJECT_TYPE.CanTestForShallowEqualityWith(s.REGEXP_TYPE)).IsTrue()
	s.AssertThat(s.OBJECT_TYPE.CanTestForShallowEqualityWith(s.STRING_TYPE)).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.CanTestForShallowEqualityWith(s.STRING_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.OBJECT_TYPE.CanTestForShallowEqualityWith(s.SYMBOL_TYPE)).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.CanTestForShallowEqualityWith(s.SYMBOL_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.OBJECT_TYPE.CanTestForShallowEqualityWith(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.OBJECT_TYPE.CanTestForShallowEqualityWith(s.VOID_TYPE)).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.CanTestForShallowEqualityWith(s.UNKNOWN_TYPE)).IsTrue()

	// isNullable
	s.AssertThat(s.OBJECT_TYPE.IsNullable()).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.IsVoidable()).IsFalse()

	// getLeastSupertype
	s.AssertTypeEquals(s.ALL_TYPE,
		s.OBJECT_TYPE.GetLeastSupertype(s.ALL_TYPE))
	s.AssertTypeEquals(s.OBJECT_TYPE,
		s.OBJECT_TYPE.GetLeastSupertype(s.STRING_OBJECT_TYPE))
	s.AssertTypeEquals(s.OBJECT_TYPE, s.OBJECT_TYPE.GetLeastSupertype(s.SYMBOL_OBJECT_TYPE))
	s.AssertTypeEquals(
		s.CreateUnionType(s.OBJECT_TYPE, s.SYMBOL_TYPE), s.OBJECT_TYPE.GetLeastSupertype(s.SYMBOL_TYPE))
	s.AssertTypeEquals(s.CreateUnionType(s.OBJECT_TYPE, s.NUMBER_TYPE),
		s.OBJECT_TYPE.GetLeastSupertype(s.NUMBER_TYPE))
	s.AssertTypeEquals(s.OBJECT_TYPE,
		s.OBJECT_TYPE.GetLeastSupertype(s.functionType))
	s.AssertTypeEquals(s.OBJECT_TYPE,
		s.OBJECT_TYPE.GetLeastSupertype(s.OBJECT_TYPE))
	s.AssertTypeEquals(s.OBJECT_TYPE,
		s.OBJECT_TYPE.GetLeastSupertype(s.DATE_TYPE))
	s.AssertTypeEquals(s.OBJECT_TYPE,
		s.OBJECT_TYPE.GetLeastSupertype(s.REGEXP_TYPE))

	// getPropertyType
	s.AssertThat(s.OBJECT_TYPE.GetPropertiesCount()).IsEqualTo(7)
	s.AssertReturnTypeEquals(s.OBJECT_TYPE, s.OBJECT_TYPE.GetPropertyType("constructor"))
	s.AssertReturnTypeEquals(s.STRING_TYPE,
		s.OBJECT_TYPE.GetPropertyType("toString"))
	s.AssertReturnTypeEquals(s.STRING_TYPE,
		s.OBJECT_TYPE.GetPropertyType("toLocaleString"))
	s.AssertReturnTypeEquals(s.UNKNOWN_TYPE,
		s.OBJECT_TYPE.GetPropertyType("valueOf"))
	s.AssertReturnTypeEquals(s.BOOLEAN_TYPE,
		s.OBJECT_TYPE.GetPropertyType("hasOwnProperty"))
	s.AssertReturnTypeEquals(s.BOOLEAN_TYPE,
		s.OBJECT_TYPE.GetPropertyType("isPrototypeOf"))
	s.AssertReturnTypeEquals(s.BOOLEAN_TYPE,
		s.OBJECT_TYPE.GetPropertyType("propertyIsEnumerable"))

	// matchesXxx
	s.AssertThat(s.OBJECT_TYPE.MatchesNumberContext()).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.MatchesObjectContext()).IsTrue()
	s.AssertThat(s.OBJECT_TYPE.MatchesStringContext()).IsTrue()
	s.AssertThat(s.OBJECT_TYPE.MatchesSymbolContext()).IsFalse()

	// implicit prototype
	s.AssertTypeEquals(s.OBJECT_PROTOTYPE, s.OBJECT_TYPE.GetImplicitPrototype())

	// toString
	s.AssertThat(s.OBJECT_TYPE.String()).IsEqualTo("Object")

	s.AssertThat(s.OBJECT_TYPE.IsNativeObjectType()).IsTrue()
	s.AssertThat(s.OBJECT_TYPE.GetImplicitPrototype().IsNativeObjectType()).IsTrue()

	rhinotesting.Asserts_AssertResolvesToSame(s.T, s.OBJECT_TYPE)
	s.AssertThat(s.OBJECT_TYPE.IsNominalConstructorOrInterface()).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.GetConstructor().IsNominalConstructorOrInterface()).IsTrue()
}

/** Tests the behavior of the number value type. */
func TestJSType_NumberObjectType(t *testing.T) {
	s := newJsTypeTest(t)
	// isXxx
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsArrayType()).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsBooleanObjectType()).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsBooleanValueType()).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsDateType()).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsEnumElementType()).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsNamedType()).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsNullType()).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsNumber()).IsTrue()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsNumberObjectType()).IsTrue()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsNumberValueType()).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsObject()).IsTrue()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsFunctionPrototypeType()).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.GetImplicitPrototype().IsFunctionPrototypeType()).IsTrue()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsRegexpType()).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsString()).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsStringObjectType()).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsStringValueType()).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsSymbol()).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsSymbolObjectType()).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsSymbolValueType()).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsEnumType()).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsUnionType()).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsStruct()).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsDict()).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsAllType()).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsVoidType()).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsConstructor()).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsInstanceType()).IsTrue()

	// autoboxesTo
	s.AssertTypeEquals(s.NUMBER_OBJECT_TYPE, s.NUMBER_TYPE.AutoboxesTo())

	// isSubtype
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsSubtypeOf(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsSubtypeOf(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsSubtypeOf(s.SYMBOL_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsSubtypeOf(s.SYMBOL_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsSubtypeOf(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsSubtypeOf(s.functionType)).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsSubtypeOf(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsSubtypeOf(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsSubtypeOf(s.DATE_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsSubtypeOf(s.unresolvedNamedType)).IsTrue()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsSubtypeOf(s.namedGoogBar)).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsSubtypeOf(s.CreateUnionType(s.NUMBER_OBJECT_TYPE, s.NULL_TYPE))).
		IsTrue()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsSubtypeOf(s.CreateUnionType(s.NUMBER_TYPE, s.NULL_TYPE))).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsSubtypeOf(s.UNKNOWN_TYPE)).IsTrue()

	// canBeCalled
	s.AssertThat(s.NUMBER_OBJECT_TYPE.CanBeCalled()).IsFalse()

	// canTestForEqualityWith
	s.AssertCanTestForEqualityWith(s.NUMBER_OBJECT_TYPE, s.NO_TYPE)
	s.AssertCanTestForEqualityWith(s.NUMBER_OBJECT_TYPE, s.NO_OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.NUMBER_OBJECT_TYPE, s.ALL_TYPE)
	s.AssertCanTestForEqualityWith(s.NUMBER_OBJECT_TYPE, s.NUMBER_TYPE)
	s.AssertCanTestForEqualityWith(s.NUMBER_OBJECT_TYPE, s.STRING_OBJECT_TYPE)
	s.AssertCannotTestForEqualityWith(s.NUMBER_OBJECT_TYPE, s.SYMBOL_TYPE)
	s.AssertCannotTestForEqualityWith(s.NUMBER_OBJECT_TYPE, s.SYMBOL_OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.NUMBER_OBJECT_TYPE, s.functionType)
	s.AssertCanTestForEqualityWith(s.NUMBER_OBJECT_TYPE, s.elementsType)
	s.AssertCannotTestForEqualityWith(s.NUMBER_OBJECT_TYPE, s.VOID_TYPE)
	s.AssertCanTestForEqualityWith(s.NUMBER_OBJECT_TYPE, s.OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.NUMBER_OBJECT_TYPE, s.DATE_TYPE)
	s.AssertCanTestForEqualityWith(s.NUMBER_OBJECT_TYPE, s.REGEXP_TYPE)
	s.AssertCanTestForEqualityWith(s.NUMBER_OBJECT_TYPE, s.ARRAY_TYPE)

	// canTestForShallowEqualityWith
	s.AssertThat(s.NUMBER_OBJECT_TYPE.CanTestForShallowEqualityWith(s.NO_TYPE)).IsTrue()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.CanTestForShallowEqualityWith(s.NO_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.CanTestForShallowEqualityWith(s.ARRAY_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.CanTestForShallowEqualityWith(s.BOOLEAN_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.CanTestForShallowEqualityWith(s.BOOLEAN_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.CanTestForShallowEqualityWith(s.DATE_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.CanTestForShallowEqualityWith(s.functionType)).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.CanTestForShallowEqualityWith(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.CanTestForShallowEqualityWith(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.CanTestForShallowEqualityWith(s.NUMBER_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.CanTestForShallowEqualityWith(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.CanTestForShallowEqualityWith(s.REGEXP_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.CanTestForShallowEqualityWith(s.STRING_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.CanTestForShallowEqualityWith(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.CanTestForShallowEqualityWith(s.SYMBOL_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.CanTestForShallowEqualityWith(s.SYMBOL_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.CanTestForShallowEqualityWith(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.CanTestForShallowEqualityWith(s.VOID_TYPE)).IsFalse()

	// isNullable
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsNullable()).IsFalse()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsVoidable()).IsFalse()

	// getLeastSupertype
	s.AssertTypeEquals(s.ALL_TYPE,
		s.NUMBER_OBJECT_TYPE.GetLeastSupertype(s.ALL_TYPE))
	s.AssertTypeEquals(s.CreateUnionType(s.NUMBER_OBJECT_TYPE, s.STRING_OBJECT_TYPE),
		s.NUMBER_OBJECT_TYPE.GetLeastSupertype(s.STRING_OBJECT_TYPE))
	s.AssertTypeEquals(
		s.CreateUnionType(s.NUMBER_OBJECT_TYPE, s.SYMBOL_OBJECT_TYPE),
		s.NUMBER_OBJECT_TYPE.GetLeastSupertype(s.SYMBOL_OBJECT_TYPE))
	s.AssertTypeEquals(
		s.CreateUnionType(s.NUMBER_OBJECT_TYPE, s.SYMBOL_TYPE),
		s.NUMBER_OBJECT_TYPE.GetLeastSupertype(s.SYMBOL_TYPE))
	s.AssertTypeEquals(s.CreateUnionType(s.NUMBER_OBJECT_TYPE, s.NUMBER_TYPE),
		s.NUMBER_OBJECT_TYPE.GetLeastSupertype(s.NUMBER_TYPE))
	s.AssertTypeEquals(s.CreateUnionType(s.NUMBER_OBJECT_TYPE, s.functionType),
		s.NUMBER_OBJECT_TYPE.GetLeastSupertype(s.functionType))
	s.AssertTypeEquals(s.OBJECT_TYPE,
		s.NUMBER_OBJECT_TYPE.GetLeastSupertype(s.OBJECT_TYPE))
	s.AssertTypeEquals(s.CreateUnionType(s.NUMBER_OBJECT_TYPE, s.DATE_TYPE),
		s.NUMBER_OBJECT_TYPE.GetLeastSupertype(s.DATE_TYPE))
	s.AssertTypeEquals(s.CreateUnionType(s.NUMBER_OBJECT_TYPE, s.REGEXP_TYPE),
		s.NUMBER_OBJECT_TYPE.GetLeastSupertype(s.REGEXP_TYPE))

	// matchesXxx
	s.AssertThat(s.NUMBER_OBJECT_TYPE.MatchesNumberContext()).IsTrue()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.MatchesObjectContext()).IsTrue()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.MatchesStringContext()).IsTrue()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.MatchesSymbolContext()).IsFalse()

	// toString
	s.AssertThat(s.NUMBER_OBJECT_TYPE.String()).IsEqualTo("Number")
	s.AssertThat(s.NUMBER_OBJECT_TYPE.HasDisplayName()).IsTrue()
	s.AssertThat(s.NUMBER_OBJECT_TYPE.GetDisplayName()).IsEqualTo("Number")

	s.AssertThat(s.NUMBER_OBJECT_TYPE.IsNativeObjectType()).IsTrue()

	rhinotesting.Asserts_AssertResolvesToSame(s.T, s.NUMBER_OBJECT_TYPE)
}

/** Tests the behavior of the number value type. */
func TestJSType_NumberValueType(t *testing.T) {
	s := newJsTypeTest(t)
	// isXxx
	s.AssertThat(s.NUMBER_TYPE.IsArrayType()).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.IsBooleanObjectType()).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.IsBooleanValueType()).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.IsDateType()).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.IsEnumElementType()).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.IsNamedType()).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.IsNullType()).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.IsNumber()).IsTrue()
	s.AssertThat(s.NUMBER_TYPE.IsNumberObjectType()).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.IsNumberValueType()).IsTrue()
	s.AssertThat(s.NUMBER_TYPE.IsFunctionPrototypeType()).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.IsRegexpType()).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.IsString()).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.IsStringObjectType()).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.IsStringValueType()).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.IsSymbol()).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.IsSymbolObjectType()).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.IsSymbolValueType()).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.IsEnumType()).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.IsUnionType()).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.IsStruct()).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.IsDict()).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.IsAllType()).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.IsVoidType()).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.IsConstructor()).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.IsInstanceType()).IsFalse()

	// autoboxesTo
	s.AssertTypeEquals(s.NUMBER_OBJECT_TYPE, s.NUMBER_TYPE.AutoboxesTo())

	// isSubtype
	s.AssertThat(s.NUMBER_TYPE.IsSubtypeOf(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.NUMBER_TYPE.IsSubtypeOf(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.IsSubtypeOf(s.SYMBOL_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.IsSubtypeOf(s.SYMBOL_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.IsSubtypeOf(s.NUMBER_TYPE)).IsTrue()
	s.AssertThat(s.NUMBER_TYPE.IsSubtypeOf(s.functionType)).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.IsSubtypeOf(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.IsSubtypeOf(s.OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.IsSubtypeOf(s.DATE_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.IsSubtypeOf(s.unresolvedNamedType)).IsTrue()
	s.AssertThat(s.NUMBER_TYPE.IsSubtypeOf(s.namedGoogBar)).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.IsSubtypeOf(s.CreateUnionType(s.NUMBER_TYPE, s.NULL_TYPE))).IsTrue()
	s.AssertThat(s.NUMBER_TYPE.IsSubtypeOf(s.UNKNOWN_TYPE)).IsTrue()

	// canBeCalled
	s.AssertThat(s.NUMBER_TYPE.CanBeCalled()).IsFalse()

	// canTestForEqualityWith
	s.AssertCanTestForEqualityWith(s.NUMBER_TYPE, s.NO_TYPE)
	s.AssertCanTestForEqualityWith(s.NUMBER_TYPE, s.NO_OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.NUMBER_TYPE, s.ALL_TYPE)
	s.AssertCanTestForEqualityWith(s.NUMBER_TYPE, s.NUMBER_TYPE)
	s.AssertCanTestForEqualityWith(s.NUMBER_TYPE, s.BIGINT_TYPE)
	s.AssertCanTestForEqualityWith(s.NUMBER_TYPE, s.STRING_OBJECT_TYPE)
	s.AssertCannotTestForEqualityWith(s.NUMBER_TYPE, s.SYMBOL_TYPE)
	s.AssertCannotTestForEqualityWith(s.NUMBER_TYPE, s.SYMBOL_OBJECT_TYPE)
	s.AssertCannotTestForEqualityWith(s.NUMBER_TYPE, s.functionType)
	s.AssertCannotTestForEqualityWith(s.NUMBER_TYPE, s.VOID_TYPE)
	s.AssertCanTestForEqualityWith(s.NUMBER_TYPE, s.OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.NUMBER_TYPE, s.DATE_TYPE)
	s.AssertCanTestForEqualityWith(s.NUMBER_TYPE, s.REGEXP_TYPE)
	s.AssertCanTestForEqualityWith(s.NUMBER_TYPE, s.ARRAY_TYPE)
	s.AssertCanTestForEqualityWith(s.NUMBER_TYPE, s.UNKNOWN_TYPE)

	// canTestForShallowEqualityWith
	s.AssertThat(s.NUMBER_TYPE.CanTestForShallowEqualityWith(s.NO_TYPE)).IsTrue()
	s.AssertThat(s.NUMBER_TYPE.CanTestForShallowEqualityWith(s.NO_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.CanTestForShallowEqualityWith(s.ARRAY_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.CanTestForShallowEqualityWith(s.BOOLEAN_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.CanTestForShallowEqualityWith(s.BOOLEAN_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.CanTestForShallowEqualityWith(s.DATE_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.CanTestForShallowEqualityWith(s.functionType)).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.CanTestForShallowEqualityWith(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.CanTestForShallowEqualityWith(s.NUMBER_TYPE)).IsTrue()
	s.AssertThat(s.NUMBER_TYPE.CanTestForShallowEqualityWith(s.NUMBER_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.CanTestForShallowEqualityWith(s.OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.CanTestForShallowEqualityWith(s.REGEXP_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.CanTestForShallowEqualityWith(s.SYMBOL_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.CanTestForShallowEqualityWith(s.SYMBOL_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.CanTestForShallowEqualityWith(s.STRING_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.CanTestForShallowEqualityWith(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.CanTestForShallowEqualityWith(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.NUMBER_TYPE.CanTestForShallowEqualityWith(s.VOID_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.CanTestForShallowEqualityWith(s.UNKNOWN_TYPE)).IsTrue()

	// isNullable
	s.AssertThat(s.NUMBER_TYPE.IsNullable()).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.IsVoidable()).IsFalse()

	// getLeastSupertype
	s.AssertTypeEquals(s.ALL_TYPE,
		s.NUMBER_TYPE.GetLeastSupertype(s.ALL_TYPE))
	s.AssertTypeEquals(s.CreateUnionType(s.NUMBER_TYPE, s.STRING_OBJECT_TYPE),
		s.NUMBER_TYPE.GetLeastSupertype(s.STRING_OBJECT_TYPE))
	s.AssertTypeEquals(
		s.CreateUnionType(s.NUMBER_TYPE, s.SYMBOL_OBJECT_TYPE),
		s.NUMBER_TYPE.GetLeastSupertype(s.SYMBOL_OBJECT_TYPE))
	s.AssertTypeEquals(
		s.CreateUnionType(s.NUMBER_TYPE, s.SYMBOL_TYPE), s.NUMBER_TYPE.GetLeastSupertype(s.SYMBOL_TYPE))
	s.AssertTypeEquals(s.NUMBER_TYPE,
		s.NUMBER_TYPE.GetLeastSupertype(s.NUMBER_TYPE))
	s.AssertTypeEquals(s.CreateUnionType(s.NUMBER_TYPE, s.functionType),
		s.NUMBER_TYPE.GetLeastSupertype(s.functionType))
	s.AssertTypeEquals(s.CreateUnionType(s.NUMBER_TYPE, s.OBJECT_TYPE),
		s.NUMBER_TYPE.GetLeastSupertype(s.OBJECT_TYPE))
	s.AssertTypeEquals(s.CreateUnionType(s.NUMBER_TYPE, s.DATE_TYPE),
		s.NUMBER_TYPE.GetLeastSupertype(s.DATE_TYPE))
	s.AssertTypeEquals(s.CreateUnionType(s.NUMBER_TYPE, s.REGEXP_TYPE),
		s.NUMBER_TYPE.GetLeastSupertype(s.REGEXP_TYPE))

	// matchesXxx
	s.AssertThat(s.NUMBER_TYPE.MatchesNumberContext()).IsTrue()
	s.AssertThat(s.NUMBER_TYPE.MatchesObjectContext()).IsTrue()
	s.AssertThat(s.NUMBER_TYPE.MatchesStringContext()).IsTrue()
	s.AssertThat(s.NUMBER_TYPE.MatchesSymbolContext()).IsFalse()

	// toString
	s.AssertThat(s.NUMBER_TYPE.String()).IsEqualTo("number")
	s.AssertThat(s.NUMBER_TYPE.HasDisplayName()).IsTrue()
	s.AssertThat(s.NUMBER_TYPE.GetDisplayName()).IsEqualTo("number")

	rhinotesting.Asserts_AssertResolvesToSame(s.T, s.NUMBER_TYPE)
	s.AssertThat(s.NUMBER_TYPE.IsNominalConstructorOrInterface()).IsFalse()
}

/** Tests the behavior of the null type. */
func TestJSType_NullType(t *testing.T) {
	s := newJsTypeTest(t)

	// isXxx
	s.AssertThat(s.NULL_TYPE.IsArrayType()).IsFalse()
	s.AssertThat(s.NULL_TYPE.IsBooleanValueType()).IsFalse()
	s.AssertThat(s.NULL_TYPE.IsDateType()).IsFalse()
	s.AssertThat(s.NULL_TYPE.IsEnumElementType()).IsFalse()
	s.AssertThat(s.NULL_TYPE.IsNamedType()).IsFalse()
	s.AssertThat(s.NULL_TYPE.IsNullType()).IsTrue()
	s.AssertThat(s.NULL_TYPE.IsNumber()).IsFalse()
	s.AssertThat(s.NULL_TYPE.IsNumberObjectType()).IsFalse()
	s.AssertThat(s.NULL_TYPE.IsNumberValueType()).IsFalse()
	s.AssertThat(s.NULL_TYPE.IsFunctionPrototypeType()).IsFalse()
	s.AssertThat(s.NULL_TYPE.IsRegexpType()).IsFalse()
	s.AssertThat(s.NULL_TYPE.IsString()).IsFalse()
	s.AssertThat(s.NULL_TYPE.IsStringObjectType()).IsFalse()
	s.AssertThat(s.NULL_TYPE.IsStringValueType()).IsFalse()
	s.AssertThat(s.NULL_TYPE.IsSymbol()).IsFalse()
	s.AssertThat(s.NULL_TYPE.IsSymbolObjectType()).IsFalse()
	s.AssertThat(s.NULL_TYPE.IsSymbolValueType()).IsFalse()
	s.AssertThat(s.NULL_TYPE.IsEnumType()).IsFalse()
	s.AssertThat(s.NULL_TYPE.IsUnionType()).IsFalse()
	s.AssertThat(s.NULL_TYPE.IsStruct()).IsFalse()
	s.AssertThat(s.NULL_TYPE.IsDict()).IsFalse()
	s.AssertThat(s.NULL_TYPE.IsAllType()).IsFalse()
	s.AssertThat(s.NULL_TYPE.IsVoidType()).IsFalse()
	s.AssertThat(s.NULL_TYPE.IsConstructor()).IsFalse()
	s.AssertThat(s.NULL_TYPE.IsInstanceType()).IsFalse()

	// autoboxesTo
	s.AssertThat(s.NULL_TYPE.AutoboxesTo()).IsNull()

	// isSubtype
	s.AssertThat(s.NULL_TYPE.IsSubtypeOf(s.NO_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.NULL_TYPE.IsSubtypeOf(s.NO_TYPE)).IsFalse()
	s.AssertThat(s.NULL_TYPE.IsSubtypeOf(s.NULL_TYPE)).IsTrue()
	s.AssertThat(s.NULL_TYPE.IsSubtypeOf(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.NULL_TYPE.IsSubtypeOf(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.NULL_TYPE.IsSubtypeOf(s.SYMBOL_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.NULL_TYPE.IsSubtypeOf(s.SYMBOL_TYPE)).IsFalse()
	s.AssertThat(s.NULL_TYPE.IsSubtypeOf(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.NULL_TYPE.IsSubtypeOf(s.functionType)).IsFalse()
	s.AssertThat(s.NULL_TYPE.IsSubtypeOf(s.OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.NULL_TYPE.IsSubtypeOf(s.DATE_TYPE)).IsFalse()
	s.AssertThat(s.NULL_TYPE.IsSubtypeOf(s.REGEXP_TYPE)).IsFalse()
	s.AssertThat(s.NULL_TYPE.IsSubtypeOf(s.ARRAY_TYPE)).IsFalse()
	s.AssertThat(s.NULL_TYPE.IsSubtypeOf(s.UNKNOWN_TYPE)).IsTrue()
	s.AssertType(s.NULL_TYPE).IsNotSubtypeOf(s.forwardDeclaredNamedType)

	s.AssertThat(s.NULL_TYPE.IsSubtypeOf(s.CreateNullableType(s.NO_OBJECT_TYPE))).IsTrue()
	s.AssertThat(s.NULL_TYPE.IsSubtypeOf(s.CreateNullableType(s.NO_TYPE))).IsTrue()
	s.AssertThat(s.NULL_TYPE.IsSubtypeOf(s.CreateNullableType(s.NULL_TYPE))).IsTrue()
	s.AssertThat(s.NULL_TYPE.IsSubtypeOf(s.CreateNullableType(s.ALL_TYPE))).IsTrue()
	s.AssertThat(s.NULL_TYPE.IsSubtypeOf(s.CreateNullableType(s.STRING_OBJECT_TYPE))).IsTrue()
	s.AssertThat(s.NULL_TYPE.IsSubtypeOf(s.CreateNullableType(s.SYMBOL_OBJECT_TYPE))).IsTrue()
	s.AssertThat(s.NULL_TYPE.IsSubtypeOf(s.CreateNullableType(s.SYMBOL_TYPE))).IsTrue()
	s.AssertThat(s.NULL_TYPE.IsSubtypeOf(s.CreateNullableType(s.NUMBER_TYPE))).IsTrue()
	s.AssertThat(s.NULL_TYPE.IsSubtypeOf(s.CreateNullableType(s.functionType))).IsTrue()
	s.AssertThat(s.NULL_TYPE.IsSubtypeOf(s.CreateNullableType(s.OBJECT_TYPE))).IsTrue()
	s.AssertThat(s.NULL_TYPE.IsSubtypeOf(s.CreateNullableType(s.DATE_TYPE))).IsTrue()
	s.AssertThat(s.NULL_TYPE.IsSubtypeOf(s.CreateNullableType(s.REGEXP_TYPE))).IsTrue()
	s.AssertThat(s.NULL_TYPE.IsSubtypeOf(s.CreateNullableType(s.ARRAY_TYPE))).IsTrue()

	// canBeCalled
	s.AssertThat(s.NULL_TYPE.CanBeCalled()).IsFalse()

	// canTestForEqualityWith
	s.AssertCanTestForEqualityWith(s.NULL_TYPE, s.NO_TYPE)
	s.AssertCanTestForEqualityWith(s.NULL_TYPE, s.NO_OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.NULL_TYPE, s.ALL_TYPE)
	s.AssertCannotTestForEqualityWith(s.NULL_TYPE, s.ARRAY_TYPE)
	s.AssertCannotTestForEqualityWith(s.NULL_TYPE, s.BOOLEAN_TYPE)
	s.AssertCannotTestForEqualityWith(s.NULL_TYPE, s.BOOLEAN_OBJECT_TYPE)
	s.AssertCannotTestForEqualityWith(s.NULL_TYPE, s.DATE_TYPE)
	s.AssertCannotTestForEqualityWith(s.NULL_TYPE, s.functionType)
	s.AssertCannotTestForEqualityWith(s.NULL_TYPE, s.NULL_TYPE)
	s.AssertCannotTestForEqualityWith(s.NULL_TYPE, s.NUMBER_TYPE)
	s.AssertCannotTestForEqualityWith(s.NULL_TYPE, s.NUMBER_OBJECT_TYPE)
	s.AssertCannotTestForEqualityWith(s.NULL_TYPE, s.OBJECT_TYPE)
	s.AssertCannotTestForEqualityWith(s.NULL_TYPE, s.REGEXP_TYPE)
	s.AssertCannotTestForEqualityWith(s.NULL_TYPE, s.STRING_TYPE)
	s.AssertCannotTestForEqualityWith(s.NULL_TYPE, s.STRING_OBJECT_TYPE)
	s.AssertCannotTestForEqualityWith(s.NULL_TYPE, s.SYMBOL_TYPE)
	s.AssertCannotTestForEqualityWith(s.NULL_TYPE, s.SYMBOL_OBJECT_TYPE)
	s.AssertCannotTestForEqualityWith(s.NULL_TYPE, s.VOID_TYPE)

	// canTestForShallowEqualityWith
	s.AssertThat(s.NULL_TYPE.CanTestForShallowEqualityWith(s.NO_TYPE)).IsTrue()
	s.AssertThat(s.NULL_TYPE.CanTestForShallowEqualityWith(s.NO_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.NULL_TYPE.CanTestForShallowEqualityWith(s.ARRAY_TYPE)).IsFalse()
	s.AssertThat(s.NULL_TYPE.CanTestForShallowEqualityWith(s.BOOLEAN_TYPE)).IsFalse()
	s.AssertThat(s.NULL_TYPE.CanTestForShallowEqualityWith(s.BOOLEAN_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.NULL_TYPE.CanTestForShallowEqualityWith(s.DATE_TYPE)).IsFalse()
	s.AssertThat(s.NULL_TYPE.CanTestForShallowEqualityWith(s.functionType)).IsFalse()
	s.AssertThat(s.NULL_TYPE.CanTestForShallowEqualityWith(s.NULL_TYPE)).IsTrue()
	s.AssertThat(s.NULL_TYPE.CanTestForShallowEqualityWith(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.NULL_TYPE.CanTestForShallowEqualityWith(s.NUMBER_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.NULL_TYPE.CanTestForShallowEqualityWith(s.OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.NULL_TYPE.CanTestForShallowEqualityWith(s.REGEXP_TYPE)).IsFalse()
	s.AssertThat(s.NULL_TYPE.CanTestForShallowEqualityWith(s.STRING_TYPE)).IsFalse()
	s.AssertThat(s.NULL_TYPE.CanTestForShallowEqualityWith(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.NULL_TYPE.CanTestForShallowEqualityWith(s.SYMBOL_TYPE)).IsFalse()
	s.AssertThat(s.NULL_TYPE.CanTestForShallowEqualityWith(s.SYMBOL_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.NULL_TYPE.CanTestForShallowEqualityWith(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.NULL_TYPE.CanTestForShallowEqualityWith(s.VOID_TYPE)).IsFalse()
	s.AssertThat(s.NULL_TYPE.CanTestForShallowEqualityWith(s.CreateNullableType(s.STRING_OBJECT_TYPE))).
		IsTrue()

	// getLeastSupertype
	s.AssertTypeEquals(s.NULL_TYPE, s.NULL_TYPE.GetLeastSupertype(s.NULL_TYPE))
	s.AssertTypeEquals(s.ALL_TYPE, s.NULL_TYPE.GetLeastSupertype(s.ALL_TYPE))
	s.AssertTypeEquals(s.CreateNullableType(s.STRING_OBJECT_TYPE),
		s.NULL_TYPE.GetLeastSupertype(s.STRING_OBJECT_TYPE))
	s.AssertTypeEquals(
		s.CreateNullableType(s.SYMBOL_OBJECT_TYPE), s.NULL_TYPE.GetLeastSupertype(s.SYMBOL_OBJECT_TYPE))
	s.AssertTypeEquals(s.CreateNullableType(s.SYMBOL_TYPE), s.NULL_TYPE.GetLeastSupertype(s.SYMBOL_TYPE))
	s.AssertTypeEquals(s.CreateNullableType(s.NUMBER_TYPE),
		s.NULL_TYPE.GetLeastSupertype(s.NUMBER_TYPE))
	s.AssertTypeEquals(s.CreateNullableType(s.functionType),
		s.NULL_TYPE.GetLeastSupertype(s.functionType))
	s.AssertTypeEquals(s.CreateNullableType(s.OBJECT_TYPE),
		s.NULL_TYPE.GetLeastSupertype(s.OBJECT_TYPE))
	s.AssertTypeEquals(s.CreateNullableType(s.DATE_TYPE),
		s.NULL_TYPE.GetLeastSupertype(s.DATE_TYPE))
	s.AssertTypeEquals(s.CreateNullableType(s.REGEXP_TYPE),
		s.NULL_TYPE.GetLeastSupertype(s.REGEXP_TYPE))

	// matchesXxx
	s.AssertThat(s.NULL_TYPE.MatchesNumberContext()).IsTrue()
	s.AssertThat(s.NULL_TYPE.MatchesObjectContext()).IsFalse()
	s.AssertThat(s.NULL_TYPE.MatchesStringContext()).IsTrue()
	s.AssertThat(s.NULL_TYPE.MatchesSymbolContext()).IsFalse()

	// toString
	s.AssertThat(s.NULL_TYPE.String()).IsEqualTo("null")
	s.AssertThat(s.NULL_TYPE.HasDisplayName()).IsTrue()
	s.AssertThat(s.NULL_TYPE.GetDisplayName()).IsEqualTo("null")

	rhinotesting.Asserts_AssertResolvesToSame(s.T, s.NULL_TYPE)

	// getGreatestSubtype
	s.AssertThat(s.NULL_TYPE.IsSubtypeOf(s.CreateUnionType(s.forwardDeclaredNamedType, s.NULL_TYPE))).
		IsTrue()
	fwdDeclaredNullUnion := s.CreateUnionType(s.forwardDeclaredNamedType, s.NULL_TYPE)
	s.AssertType(s.NULL_TYPE.GetGreatestSubtype(fwdDeclaredNullUnion)).IsEqualTo(s.NULL_TYPE)
	s.AssertThat(s.NULL_TYPE.IsNominalConstructorOrInterface()).IsFalse()

	s.AssertThat(s.NULL_TYPE.DiffersFrom(s.UNKNOWN_TYPE)).IsTrue()
}

func TestJSType_DateTypeX(t *testing.T) {
	s := newJsTypeTest(t)
	s.AssertCannotTestForEqualityWith(s.SYMBOL_OBJECT_TYPE, s.functionType)
}

/** Tests the behavior of the Date type. */
func TestJSType_DateType(t *testing.T) {
	s := newJsTypeTest(t)
	// isXxx
	s.AssertThat(s.DATE_TYPE.IsArrayType()).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsBooleanValueType()).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsDateType()).IsTrue()
	s.AssertThat(s.DATE_TYPE.IsEnumElementType()).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsNamedType()).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsNullType()).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsNumberValueType()).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsFunctionPrototypeType()).IsFalse()
	s.AssertThat(s.DATE_TYPE.GetImplicitPrototype().IsFunctionPrototypeType()).IsTrue()
	s.AssertThat(s.DATE_TYPE.IsRegexpType()).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsStringValueType()).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsSymbolValueType()).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsEnumType()).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsUnionType()).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsStruct()).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsDict()).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsAllType()).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsVoidType()).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsConstructor()).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsInstanceType()).IsTrue()

	// autoboxesTo
	s.AssertThat(s.DATE_TYPE.AutoboxesTo()).IsNull()

	// isSubtype
	s.AssertThat(s.DATE_TYPE.IsSubtypeOf(s.NO_TYPE)).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsSubtypeOf(s.NO_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsSubtypeOf(s.ARRAY_TYPE)).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsSubtypeOf(s.BOOLEAN_TYPE)).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsSubtypeOf(s.BOOLEAN_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsSubtypeOf(s.DATE_TYPE)).IsTrue()
	s.AssertThat(s.DATE_TYPE.IsSubtypeOf(s.functionType)).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsSubtypeOf(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsSubtypeOf(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsSubtypeOf(s.NUMBER_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsSubtypeOf(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.DATE_TYPE.IsSubtypeOf(s.REGEXP_TYPE)).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsSubtypeOf(s.STRING_TYPE)).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsSubtypeOf(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsSubtypeOf(s.SYMBOL_TYPE)).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsSubtypeOf(s.SYMBOL_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsSubtypeOf(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.DATE_TYPE.IsSubtypeOf(s.VOID_TYPE)).IsFalse()

	// canBeCalled
	s.AssertThat(s.DATE_TYPE.CanBeCalled()).IsFalse()

	// canTestForEqualityWith
	s.AssertCanTestForEqualityWith(s.DATE_TYPE, s.ALL_TYPE)
	s.AssertCanTestForEqualityWith(s.DATE_TYPE, s.STRING_OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.DATE_TYPE, s.NUMBER_TYPE)
	s.AssertCannotTestForEqualityWith(s.DATE_TYPE, s.SYMBOL_OBJECT_TYPE) // fix this
	s.AssertCannotTestForEqualityWith(s.DATE_TYPE, s.SYMBOL_TYPE)        // fix this
	s.AssertCanTestForEqualityWith(s.DATE_TYPE, s.functionType)
	s.AssertCannotTestForEqualityWith(s.DATE_TYPE, s.VOID_TYPE)
	s.AssertCanTestForEqualityWith(s.DATE_TYPE, s.OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.DATE_TYPE, s.DATE_TYPE)
	s.AssertCanTestForEqualityWith(s.DATE_TYPE, s.REGEXP_TYPE)
	s.AssertCanTestForEqualityWith(s.DATE_TYPE, s.ARRAY_TYPE)

	// canTestForShallowEqualityWith
	s.AssertThat(s.DATE_TYPE.CanTestForShallowEqualityWith(s.NO_TYPE)).IsTrue()
	s.AssertThat(s.DATE_TYPE.CanTestForShallowEqualityWith(s.NO_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.DATE_TYPE.CanTestForShallowEqualityWith(s.ARRAY_TYPE)).IsFalse()
	s.AssertThat(s.DATE_TYPE.CanTestForShallowEqualityWith(s.BOOLEAN_TYPE)).IsFalse()
	s.AssertThat(s.DATE_TYPE.CanTestForShallowEqualityWith(s.BOOLEAN_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.DATE_TYPE.CanTestForShallowEqualityWith(s.DATE_TYPE)).IsTrue()
	s.AssertThat(s.DATE_TYPE.CanTestForShallowEqualityWith(s.functionType)).IsFalse()
	s.AssertThat(s.DATE_TYPE.CanTestForShallowEqualityWith(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.DATE_TYPE.CanTestForShallowEqualityWith(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.DATE_TYPE.CanTestForShallowEqualityWith(s.NUMBER_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.DATE_TYPE.CanTestForShallowEqualityWith(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.DATE_TYPE.CanTestForShallowEqualityWith(s.REGEXP_TYPE)).IsFalse()
	s.AssertThat(s.DATE_TYPE.CanTestForShallowEqualityWith(s.STRING_TYPE)).IsFalse()
	s.AssertThat(s.DATE_TYPE.CanTestForShallowEqualityWith(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.DATE_TYPE.CanTestForShallowEqualityWith(s.SYMBOL_TYPE)).IsFalse()
	s.AssertThat(s.DATE_TYPE.CanTestForShallowEqualityWith(s.SYMBOL_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.DATE_TYPE.CanTestForShallowEqualityWith(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.DATE_TYPE.CanTestForShallowEqualityWith(s.VOID_TYPE)).IsFalse()

	// isNullable
	s.AssertThat(s.DATE_TYPE.IsNullable()).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsVoidable()).IsFalse()
	s.AssertThat(s.CreateNullableType(s.DATE_TYPE).IsNullable()).IsTrue()
	s.AssertThat(s.CreateUnionType(s.DATE_TYPE, s.VOID_TYPE).IsVoidable()).IsTrue()

	// getLeastSupertype
	s.AssertTypeEquals(s.ALL_TYPE,
		s.DATE_TYPE.GetLeastSupertype(s.ALL_TYPE))
	s.AssertTypeEquals(s.CreateUnionType(s.DATE_TYPE, s.STRING_OBJECT_TYPE),
		s.DATE_TYPE.GetLeastSupertype(s.STRING_OBJECT_TYPE))
	s.AssertTypeEquals(s.CreateUnionType(s.DATE_TYPE, s.NUMBER_TYPE),
		s.DATE_TYPE.GetLeastSupertype(s.NUMBER_TYPE))
	s.AssertTypeEquals(s.CreateUnionType(s.DATE_TYPE, s.functionType),
		s.DATE_TYPE.GetLeastSupertype(s.functionType))
	s.AssertTypeEquals(s.OBJECT_TYPE, s.DATE_TYPE.GetLeastSupertype(s.OBJECT_TYPE))
	s.AssertTypeEquals(s.DATE_TYPE, s.DATE_TYPE.GetLeastSupertype(s.DATE_TYPE))
	s.AssertTypeEquals(s.CreateUnionType(s.DATE_TYPE, s.REGEXP_TYPE),
		s.DATE_TYPE.GetLeastSupertype(s.REGEXP_TYPE))

	// getPropertyType
	s.AssertThat(s.DATE_TYPE.GetImplicitPrototype().GetPropertiesCount()).IsEqualTo(46)
	s.AssertThat(s.DATE_TYPE.GetPropertiesCount()).IsEqualTo(46)
	s.AssertReturnTypeEquals(s.DATE_TYPE, s.DATE_TYPE.GetPropertyType("constructor"))
	s.AssertReturnTypeEquals(s.STRING_TYPE,
		s.DATE_TYPE.GetPropertyType("toString"))
	s.AssertReturnTypeEquals(s.STRING_TYPE,
		s.DATE_TYPE.GetPropertyType("toDateString"))
	s.AssertReturnTypeEquals(s.STRING_TYPE,
		s.DATE_TYPE.GetPropertyType("toTimeString"))
	s.AssertReturnTypeEquals(s.STRING_TYPE,
		s.DATE_TYPE.GetPropertyType("toLocaleString"))
	s.AssertReturnTypeEquals(s.STRING_TYPE,
		s.DATE_TYPE.GetPropertyType("toLocaleDateString"))
	s.AssertReturnTypeEquals(s.STRING_TYPE,
		s.DATE_TYPE.GetPropertyType("toLocaleTimeString"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE, s.DATE_TYPE.GetPropertyType("valueOf"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE, s.DATE_TYPE.GetPropertyType("getTime"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE,
		s.DATE_TYPE.GetPropertyType("getFullYear"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE,
		s.DATE_TYPE.GetPropertyType("getUTCFullYear"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE, s.DATE_TYPE.GetPropertyType("getMonth"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE,
		s.DATE_TYPE.GetPropertyType("getUTCMonth"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE, s.DATE_TYPE.GetPropertyType("getDate"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE,
		s.DATE_TYPE.GetPropertyType("getUTCDate"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE, s.DATE_TYPE.GetPropertyType("getDay"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE, s.DATE_TYPE.GetPropertyType("getUTCDay"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE, s.DATE_TYPE.GetPropertyType("getHours"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE,
		s.DATE_TYPE.GetPropertyType("getUTCHours"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE,
		s.DATE_TYPE.GetPropertyType("getMinutes"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE,
		s.DATE_TYPE.GetPropertyType("getUTCMinutes"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE,
		s.DATE_TYPE.GetPropertyType("getSeconds"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE,
		s.DATE_TYPE.GetPropertyType("getUTCSeconds"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE,
		s.DATE_TYPE.GetPropertyType("getMilliseconds"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE,
		s.DATE_TYPE.GetPropertyType("getUTCMilliseconds"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE,
		s.DATE_TYPE.GetPropertyType("getTimezoneOffset"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE, s.DATE_TYPE.GetPropertyType("setTime"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE,
		s.DATE_TYPE.GetPropertyType("setMilliseconds"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE,
		s.DATE_TYPE.GetPropertyType("setUTCMilliseconds"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE,
		s.DATE_TYPE.GetPropertyType("setSeconds"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE,
		s.DATE_TYPE.GetPropertyType("setUTCSeconds"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE,
		s.DATE_TYPE.GetPropertyType("setUTCSeconds"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE,
		s.DATE_TYPE.GetPropertyType("setMinutes"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE,
		s.DATE_TYPE.GetPropertyType("setUTCMinutes"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE, s.DATE_TYPE.GetPropertyType("setHours"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE,
		s.DATE_TYPE.GetPropertyType("setUTCHours"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE, s.DATE_TYPE.GetPropertyType("setDate"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE,
		s.DATE_TYPE.GetPropertyType("setUTCDate"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE, s.DATE_TYPE.GetPropertyType("setMonth"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE,
		s.DATE_TYPE.GetPropertyType("setUTCMonth"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE,
		s.DATE_TYPE.GetPropertyType("setFullYear"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE,
		s.DATE_TYPE.GetPropertyType("setUTCFullYear"))
	s.AssertReturnTypeEquals(s.STRING_TYPE,
		s.DATE_TYPE.GetPropertyType("toUTCString"))
	s.AssertReturnTypeEquals(s.STRING_TYPE,
		s.DATE_TYPE.GetPropertyType("toGMTString"))

	// matchesXxx
	s.AssertThat(s.DATE_TYPE.MatchesNumberContext()).IsTrue()
	s.AssertThat(s.DATE_TYPE.MatchesObjectContext()).IsTrue()
	s.AssertThat(s.DATE_TYPE.MatchesStringContext()).IsTrue()
	s.AssertThat(s.DATE_TYPE.MatchesSymbolContext()).IsFalse()

	// toString
	s.AssertThat(s.DATE_TYPE.String()).IsEqualTo("Date")
	s.AssertThat(s.DATE_TYPE.HasDisplayName()).IsTrue()
	s.AssertThat(s.DATE_TYPE.GetDisplayName()).IsEqualTo("Date")

	s.AssertThat(s.DATE_TYPE.IsNativeObjectType()).IsTrue()

	rhinotesting.Asserts_AssertResolvesToSame(s.T, s.DATE_TYPE)
	s.AssertThat(s.DATE_TYPE.IsNominalConstructorOrInterface()).IsFalse()
	s.AssertThat(s.DATE_TYPE.GetConstructor().IsNominalConstructorOrInterface()).IsTrue()
}

/** Tests the behavior of the RegExp type. */
func TestJSType_RegExpType(t *testing.T) {
	s := newJsTypeTest(t)
	// isXxx
	s.AssertThat(s.REGEXP_TYPE.IsNoType()).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.IsNoObjectType()).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.IsArrayType()).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.IsBooleanValueType()).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.IsDateType()).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.IsEnumElementType()).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.IsNamedType()).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.IsNullType()).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.IsNumberValueType()).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.IsFunctionPrototypeType()).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.GetImplicitPrototype().IsFunctionPrototypeType()).IsTrue()
	s.AssertThat(s.REGEXP_TYPE.IsRegexpType()).IsTrue()
	s.AssertThat(s.REGEXP_TYPE.IsStringValueType()).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.IsSymbolValueType()).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.IsEnumType()).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.IsUnionType()).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.IsStruct()).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.IsDict()).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.IsAllType()).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.IsVoidType()).IsFalse()

	// autoboxesTo
	s.AssertThat(s.REGEXP_TYPE.AutoboxesTo()).IsNull()

	// isSubtype
	s.AssertThat(s.REGEXP_TYPE.IsSubtypeOf(s.NO_TYPE)).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.IsSubtypeOf(s.NO_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.IsSubtypeOf(s.ARRAY_TYPE)).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.IsSubtypeOf(s.BOOLEAN_TYPE)).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.IsSubtypeOf(s.BOOLEAN_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.IsSubtypeOf(s.DATE_TYPE)).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.IsSubtypeOf(s.functionType)).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.IsSubtypeOf(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.IsSubtypeOf(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.IsSubtypeOf(s.NUMBER_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.IsSubtypeOf(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.REGEXP_TYPE.IsSubtypeOf(s.REGEXP_TYPE)).IsTrue()
	s.AssertThat(s.REGEXP_TYPE.IsSubtypeOf(s.STRING_TYPE)).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.IsSubtypeOf(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.IsSubtypeOf(s.SYMBOL_TYPE)).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.IsSubtypeOf(s.SYMBOL_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.IsSubtypeOf(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.REGEXP_TYPE.IsSubtypeOf(s.VOID_TYPE)).IsFalse()

	// canBeCalled
	s.AssertThat(s.REGEXP_TYPE.CanBeCalled()).IsFalse()

	// canTestForEqualityWith
	s.AssertCanTestForEqualityWith(s.REGEXP_TYPE, s.ALL_TYPE)
	s.AssertCanTestForEqualityWith(s.REGEXP_TYPE, s.STRING_OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.REGEXP_TYPE, s.NUMBER_TYPE)
	s.AssertCannotTestForEqualityWith(s.REGEXP_TYPE, s.SYMBOL_OBJECT_TYPE)
	s.AssertCannotTestForEqualityWith(s.REGEXP_TYPE, s.SYMBOL_TYPE)
	s.AssertCanTestForEqualityWith(s.REGEXP_TYPE, s.functionType)
	s.AssertCannotTestForEqualityWith(s.REGEXP_TYPE, s.VOID_TYPE)
	s.AssertCanTestForEqualityWith(s.REGEXP_TYPE, s.OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.REGEXP_TYPE, s.DATE_TYPE)
	s.AssertCanTestForEqualityWith(s.REGEXP_TYPE, s.REGEXP_TYPE)
	s.AssertCanTestForEqualityWith(s.REGEXP_TYPE, s.ARRAY_TYPE)

	// canTestForShallowEqualityWith
	s.AssertThat(s.REGEXP_TYPE.CanTestForShallowEqualityWith(s.NO_TYPE)).IsTrue()
	s.AssertThat(s.REGEXP_TYPE.CanTestForShallowEqualityWith(s.NO_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.REGEXP_TYPE.CanTestForShallowEqualityWith(s.ARRAY_TYPE)).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.CanTestForShallowEqualityWith(s.BOOLEAN_TYPE)).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.CanTestForShallowEqualityWith(s.BOOLEAN_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.CanTestForShallowEqualityWith(s.DATE_TYPE)).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.CanTestForShallowEqualityWith(s.functionType)).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.CanTestForShallowEqualityWith(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.CanTestForShallowEqualityWith(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.CanTestForShallowEqualityWith(s.NUMBER_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.CanTestForShallowEqualityWith(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.REGEXP_TYPE.CanTestForShallowEqualityWith(s.REGEXP_TYPE)).IsTrue()
	s.AssertThat(s.REGEXP_TYPE.CanTestForShallowEqualityWith(s.STRING_TYPE)).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.CanTestForShallowEqualityWith(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.CanTestForShallowEqualityWith(s.SYMBOL_TYPE)).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.CanTestForShallowEqualityWith(s.SYMBOL_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.CanTestForShallowEqualityWith(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.REGEXP_TYPE.CanTestForShallowEqualityWith(s.VOID_TYPE)).IsFalse()

	// isNullable
	s.AssertThat(s.REGEXP_TYPE.IsNullable()).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.IsVoidable()).IsFalse()
	s.AssertThat(s.CreateNullableType(s.REGEXP_TYPE).IsNullable()).IsTrue()
	s.AssertThat(s.CreateUnionType(s.REGEXP_TYPE, s.VOID_TYPE).IsVoidable()).IsTrue()

	// getLeastSupertype
	s.AssertTypeEquals(s.ALL_TYPE,
		s.REGEXP_TYPE.GetLeastSupertype(s.ALL_TYPE))
	s.AssertTypeEquals(s.CreateUnionType(s.REGEXP_TYPE, s.STRING_OBJECT_TYPE),
		s.REGEXP_TYPE.GetLeastSupertype(s.STRING_OBJECT_TYPE))
	s.AssertTypeEquals(s.CreateUnionType(s.REGEXP_TYPE, s.NUMBER_TYPE),
		s.REGEXP_TYPE.GetLeastSupertype(s.NUMBER_TYPE))
	s.AssertTypeEquals(s.CreateUnionType(s.REGEXP_TYPE, s.functionType),
		s.REGEXP_TYPE.GetLeastSupertype(s.functionType))
	s.AssertTypeEquals(s.OBJECT_TYPE, s.REGEXP_TYPE.GetLeastSupertype(s.OBJECT_TYPE))
	s.AssertTypeEquals(s.CreateUnionType(s.DATE_TYPE, s.REGEXP_TYPE),
		s.REGEXP_TYPE.GetLeastSupertype(s.DATE_TYPE))
	s.AssertTypeEquals(s.REGEXP_TYPE,
		s.REGEXP_TYPE.GetLeastSupertype(s.REGEXP_TYPE))

	// getPropertyType
	s.AssertThat(s.REGEXP_TYPE.GetImplicitPrototype().GetPropertiesCount()).IsEqualTo(9)
	s.AssertThat(s.REGEXP_TYPE.GetPropertiesCount()).IsEqualTo(14)
	s.AssertReturnTypeEquals(s.REGEXP_TYPE, s.REGEXP_TYPE.GetPropertyType("constructor"))
	s.AssertReturnTypeEquals(s.CreateNullableType(s.ARRAY_TYPE),
		s.REGEXP_TYPE.GetPropertyType("exec"))
	s.AssertReturnTypeEquals(s.BOOLEAN_TYPE,
		s.REGEXP_TYPE.GetPropertyType("test"))
	s.AssertReturnTypeEquals(s.STRING_TYPE,
		s.REGEXP_TYPE.GetPropertyType("toString"))
	s.AssertTypeEquals(s.STRING_TYPE, s.REGEXP_TYPE.GetPropertyType("source"))
	s.AssertTypeEquals(s.BOOLEAN_TYPE, s.REGEXP_TYPE.GetPropertyType("global"))
	s.AssertTypeEquals(s.BOOLEAN_TYPE, s.REGEXP_TYPE.GetPropertyType("ignoreCase"))
	s.AssertTypeEquals(s.BOOLEAN_TYPE, s.REGEXP_TYPE.GetPropertyType("multiline"))
	s.AssertTypeEquals(s.NUMBER_TYPE, s.REGEXP_TYPE.GetPropertyType("lastIndex"))

	// matchesXxx
	s.AssertThat(s.REGEXP_TYPE.MatchesNumberContext()).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.MatchesObjectContext()).IsTrue()
	s.AssertThat(s.REGEXP_TYPE.MatchesStringContext()).IsTrue()
	s.AssertThat(s.REGEXP_TYPE.MatchesSymbolContext()).IsFalse()

	// toString
	s.AssertThat(s.REGEXP_TYPE.String()).IsEqualTo("RegExp")
	s.AssertThat(s.REGEXP_TYPE.HasDisplayName()).IsTrue()
	s.AssertThat(s.REGEXP_TYPE.GetDisplayName()).IsEqualTo("RegExp")

	s.AssertThat(s.REGEXP_TYPE.IsNativeObjectType()).IsTrue()

	rhinotesting.Asserts_AssertResolvesToSame(s.T, s.REGEXP_TYPE)
	s.AssertThat(s.REGEXP_TYPE.IsNominalConstructorOrInterface()).IsFalse()
	s.AssertThat(s.REGEXP_TYPE.GetConstructor().IsNominalConstructorOrInterface()).IsTrue()
}

/** Tests the behavior of the string object type. */
func TestJSType_StringObjectType(t *testing.T) {
	s := newJsTypeTest(t)
	// isXxx
	s.AssertThat(s.STRING_OBJECT_TYPE.IsArrayType()).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.IsBooleanObjectType()).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.IsBooleanValueType()).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.IsDateType()).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.IsEnumElementType()).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.IsNamedType()).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.IsNullType()).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.IsNumber()).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.IsNumberObjectType()).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.IsNumberValueType()).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.IsFunctionPrototypeType()).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.GetImplicitPrototype().IsFunctionPrototypeType()).IsTrue()
	s.AssertThat(s.STRING_OBJECT_TYPE.IsRegexpType()).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.IsString()).IsTrue()
	s.AssertThat(s.STRING_OBJECT_TYPE.IsStringObjectType()).IsTrue()
	s.AssertThat(s.STRING_OBJECT_TYPE.IsStringValueType()).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.IsSymbol()).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.IsSymbolObjectType()).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.IsSymbolValueType()).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.IsEnumType()).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.IsUnionType()).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.IsStruct()).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.IsDict()).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.IsAllType()).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.IsVoidType()).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.IsConstructor()).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.IsInstanceType()).IsTrue()

	// autoboxesTo
	s.AssertTypeEquals(s.STRING_OBJECT_TYPE, s.STRING_TYPE.AutoboxesTo())

	// isSubtype
	s.AssertThat(s.STRING_OBJECT_TYPE.IsSubtypeOf(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.STRING_OBJECT_TYPE.IsSubtypeOf(s.STRING_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.STRING_OBJECT_TYPE.IsSubtypeOf(s.STRING_TYPE)).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.IsSubtypeOf(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.STRING_OBJECT_TYPE.IsSubtypeOf(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.IsSubtypeOf(s.SYMBOL_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.IsSubtypeOf(s.SYMBOL_TYPE)).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.IsSubtypeOf(s.DATE_TYPE)).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.IsSubtypeOf(s.REGEXP_TYPE)).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.IsSubtypeOf(s.ARRAY_TYPE)).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.IsSubtypeOf(s.STRING_TYPE)).IsFalse()

	// canBeCalled
	s.AssertThat(s.STRING_OBJECT_TYPE.CanBeCalled()).IsFalse()

	// canTestForEqualityWith
	s.AssertCanTestForEqualityWith(s.STRING_OBJECT_TYPE, s.ALL_TYPE)
	s.AssertCanTestForEqualityWith(s.STRING_OBJECT_TYPE, s.STRING_OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.STRING_OBJECT_TYPE, s.STRING_TYPE)
	s.AssertCannotTestForEqualityWith(s.STRING_OBJECT_TYPE, s.SYMBOL_OBJECT_TYPE)
	s.AssertCannotTestForEqualityWith(s.STRING_OBJECT_TYPE, s.SYMBOL_TYPE)
	s.AssertCanTestForEqualityWith(s.STRING_OBJECT_TYPE, s.functionType)
	s.AssertCanTestForEqualityWith(s.STRING_OBJECT_TYPE, s.OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.STRING_OBJECT_TYPE, s.NUMBER_TYPE)
	s.AssertCanTestForEqualityWith(s.STRING_OBJECT_TYPE, s.BOOLEAN_TYPE)
	s.AssertCanTestForEqualityWith(s.STRING_OBJECT_TYPE, s.BOOLEAN_OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.STRING_OBJECT_TYPE, s.DATE_TYPE)
	s.AssertCanTestForEqualityWith(s.STRING_OBJECT_TYPE, s.REGEXP_TYPE)
	s.AssertCanTestForEqualityWith(s.STRING_OBJECT_TYPE, s.ARRAY_TYPE)
	s.AssertCanTestForEqualityWith(s.STRING_OBJECT_TYPE, s.UNKNOWN_TYPE)

	// canTestForShallowEqualityWith
	s.AssertThat(s.STRING_OBJECT_TYPE.CanTestForShallowEqualityWith(s.NO_TYPE)).IsTrue()
	s.AssertThat(s.STRING_OBJECT_TYPE.CanTestForShallowEqualityWith(s.NO_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.STRING_OBJECT_TYPE.CanTestForShallowEqualityWith(s.ARRAY_TYPE)).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.CanTestForShallowEqualityWith(s.BOOLEAN_TYPE)).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.CanTestForShallowEqualityWith(s.BOOLEAN_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.CanTestForShallowEqualityWith(s.DATE_TYPE)).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.CanTestForShallowEqualityWith(s.functionType)).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.CanTestForShallowEqualityWith(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.CanTestForShallowEqualityWith(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.CanTestForShallowEqualityWith(s.NUMBER_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.CanTestForShallowEqualityWith(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.STRING_OBJECT_TYPE.CanTestForShallowEqualityWith(s.REGEXP_TYPE)).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.CanTestForShallowEqualityWith(s.STRING_TYPE)).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.CanTestForShallowEqualityWith(s.STRING_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.STRING_OBJECT_TYPE.CanTestForShallowEqualityWith(s.SYMBOL_TYPE)).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.CanTestForShallowEqualityWith(s.SYMBOL_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.CanTestForShallowEqualityWith(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.STRING_OBJECT_TYPE.CanTestForShallowEqualityWith(s.VOID_TYPE)).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.CanTestForShallowEqualityWith(s.UNKNOWN_TYPE)).IsTrue()

	// properties (ECMA-262 page 98 - 106)
	s.AssertThat(s.STRING_OBJECT_TYPE.GetImplicitPrototype().GetPropertiesCount()).IsEqualTo(24)
	s.AssertThat(s.STRING_OBJECT_TYPE.GetPropertiesCount()).IsEqualTo(25)

	s.AssertReturnTypeEquals(s.STRING_TYPE,
		s.STRING_OBJECT_TYPE.GetPropertyType("toString"))
	s.AssertReturnTypeEquals(s.STRING_TYPE,
		s.STRING_OBJECT_TYPE.GetPropertyType("valueOf"))
	s.AssertReturnTypeEquals(s.STRING_TYPE,
		s.STRING_OBJECT_TYPE.GetPropertyType("charAt"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE,
		s.STRING_OBJECT_TYPE.GetPropertyType("charCodeAt"))
	s.AssertReturnTypeEquals(s.STRING_TYPE,
		s.STRING_OBJECT_TYPE.GetPropertyType("concat"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE,
		s.STRING_OBJECT_TYPE.GetPropertyType("indexOf"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE,
		s.STRING_OBJECT_TYPE.GetPropertyType("lastIndexOf"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE,
		s.STRING_OBJECT_TYPE.GetPropertyType("localeCompare"))
	s.AssertReturnTypeEquals(s.CreateNullableType(s.ARRAY_TYPE),
		s.STRING_OBJECT_TYPE.GetPropertyType("match"))
	s.AssertReturnTypeEquals(s.STRING_TYPE,
		s.STRING_OBJECT_TYPE.GetPropertyType("replace"))
	s.AssertReturnTypeEquals(s.NUMBER_TYPE,
		s.STRING_OBJECT_TYPE.GetPropertyType("search"))
	s.AssertReturnTypeEquals(s.STRING_TYPE,
		s.STRING_OBJECT_TYPE.GetPropertyType("slice"))
	s.AssertReturnTypeEquals(s.ARRAY_TYPE,
		s.STRING_OBJECT_TYPE.GetPropertyType("split"))
	s.AssertReturnTypeEquals(s.STRING_TYPE,
		s.STRING_OBJECT_TYPE.GetPropertyType("substr"))
	s.AssertReturnTypeEquals(s.STRING_TYPE,
		s.STRING_OBJECT_TYPE.GetPropertyType("substring"))
	s.AssertReturnTypeEquals(s.STRING_TYPE,
		s.STRING_OBJECT_TYPE.GetPropertyType("toLowerCase"))
	s.AssertReturnTypeEquals(s.STRING_TYPE,
		s.STRING_OBJECT_TYPE.GetPropertyType("toLocaleLowerCase"))
	s.AssertReturnTypeEquals(s.STRING_TYPE,
		s.STRING_OBJECT_TYPE.GetPropertyType("toUpperCase"))
	s.AssertReturnTypeEquals(s.STRING_TYPE,
		s.STRING_OBJECT_TYPE.GetPropertyType("toLocaleUpperCase"))
	s.AssertTypeEquals(s.NUMBER_TYPE, s.STRING_OBJECT_TYPE.GetPropertyType("length"))

	// matchesXxx
	s.AssertThat(s.STRING_OBJECT_TYPE.MatchesNumberContext()).IsTrue()
	s.AssertThat(s.STRING_OBJECT_TYPE.MatchesObjectContext()).IsTrue()
	s.AssertThat(s.STRING_OBJECT_TYPE.MatchesStringContext()).IsTrue()
	s.AssertThat(s.STRING_OBJECT_TYPE.MatchesSymbolContext()).IsFalse()

	// isNullable
	s.AssertThat(s.STRING_OBJECT_TYPE.IsNullable()).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.IsVoidable()).IsFalse()
	s.AssertThat(s.CreateNullableType(s.STRING_OBJECT_TYPE).IsNullable()).IsTrue()
	s.AssertThat(s.CreateUnionType(s.STRING_OBJECT_TYPE, s.VOID_TYPE).IsVoidable()).IsTrue()

	s.AssertThat(s.STRING_OBJECT_TYPE.IsNativeObjectType()).IsTrue()

	rhinotesting.Asserts_AssertResolvesToSame(s.T, s.STRING_OBJECT_TYPE)

	s.AssertThat(s.STRING_OBJECT_TYPE.HasDisplayName()).IsTrue()
	s.AssertThat(s.STRING_OBJECT_TYPE.GetDisplayName()).IsEqualTo("String")
	s.AssertThat(s.STRING_OBJECT_TYPE.IsNominalConstructorOrInterface()).IsFalse()
	s.AssertThat(s.STRING_OBJECT_TYPE.GetConstructor().IsNominalConstructorOrInterface()).IsTrue()
}

/** Tests the behavior of the string value type. */
func TestJSType_StringValueType(t *testing.T) {
	s := newJsTypeTest(t)
	// isXxx
	s.AssertThat(s.STRING_TYPE.IsArrayType()).IsFalse()
	s.AssertThat(s.STRING_TYPE.IsBooleanObjectType()).IsFalse()
	s.AssertThat(s.STRING_TYPE.IsBooleanValueType()).IsFalse()
	s.AssertThat(s.STRING_TYPE.IsDateType()).IsFalse()
	s.AssertThat(s.STRING_TYPE.IsEnumElementType()).IsFalse()
	s.AssertThat(s.STRING_TYPE.IsNamedType()).IsFalse()
	s.AssertThat(s.STRING_TYPE.IsNullType()).IsFalse()
	s.AssertThat(s.STRING_TYPE.IsNumber()).IsFalse()
	s.AssertThat(s.STRING_TYPE.IsNumberObjectType()).IsFalse()
	s.AssertThat(s.STRING_TYPE.IsNumberValueType()).IsFalse()
	s.AssertThat(s.STRING_TYPE.IsFunctionPrototypeType()).IsFalse()
	s.AssertThat(s.STRING_TYPE.IsRegexpType()).IsFalse()
	s.AssertThat(s.STRING_TYPE.IsString()).IsTrue()
	s.AssertThat(s.STRING_TYPE.IsStringObjectType()).IsFalse()
	s.AssertThat(s.STRING_TYPE.IsStringValueType()).IsTrue()
	s.AssertThat(s.STRING_TYPE.IsSymbol()).IsFalse()
	s.AssertThat(s.STRING_TYPE.IsSymbolObjectType()).IsFalse()
	s.AssertThat(s.STRING_TYPE.IsSymbolValueType()).IsFalse()
	s.AssertThat(s.STRING_TYPE.IsEnumType()).IsFalse()
	s.AssertThat(s.STRING_TYPE.IsUnionType()).IsFalse()
	s.AssertThat(s.STRING_TYPE.IsStruct()).IsFalse()
	s.AssertThat(s.STRING_TYPE.IsDict()).IsFalse()
	s.AssertThat(s.STRING_TYPE.IsAllType()).IsFalse()
	s.AssertThat(s.STRING_TYPE.IsVoidType()).IsFalse()
	s.AssertThat(s.STRING_TYPE.IsConstructor()).IsFalse()
	s.AssertThat(s.STRING_TYPE.IsInstanceType()).IsFalse()

	// autoboxesTo
	s.AssertTypeEquals(s.STRING_OBJECT_TYPE, s.STRING_TYPE.AutoboxesTo())

	// isSubtype
	s.AssertThat(s.STRING_TYPE.IsSubtypeOf(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.STRING_TYPE.IsSubtypeOf(s.STRING_TYPE)).IsTrue()
	s.AssertThat(s.STRING_TYPE.IsSubtypeOf(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.STRING_TYPE.IsSubtypeOf(s.SYMBOL_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.STRING_TYPE.IsSubtypeOf(s.SYMBOL_TYPE)).IsFalse()
	s.AssertThat(s.STRING_TYPE.IsSubtypeOf(s.OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.STRING_TYPE.IsSubtypeOf(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.STRING_TYPE.IsSubtypeOf(s.DATE_TYPE)).IsFalse()
	s.AssertThat(s.STRING_TYPE.IsSubtypeOf(s.REGEXP_TYPE)).IsFalse()
	s.AssertThat(s.STRING_TYPE.IsSubtypeOf(s.ARRAY_TYPE)).IsFalse()
	s.AssertThat(s.STRING_TYPE.IsSubtypeOf(s.UNKNOWN_TYPE)).IsTrue()

	// canBeCalled
	s.AssertThat(s.STRING_TYPE.CanBeCalled()).IsFalse()

	// canTestForEqualityWith
	s.AssertCanTestForEqualityWith(s.STRING_TYPE, s.ALL_TYPE)
	s.AssertCanTestForEqualityWith(s.STRING_TYPE, s.STRING_OBJECT_TYPE)
	s.AssertCannotTestForEqualityWith(s.STRING_TYPE, s.functionType)
	s.AssertCanTestForEqualityWith(s.STRING_TYPE, s.OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.STRING_TYPE, s.NUMBER_TYPE)
	s.AssertCanTestForEqualityWith(s.STRING_TYPE, s.BIGINT_TYPE)
	s.AssertCanTestForEqualityWith(s.STRING_TYPE, s.BOOLEAN_TYPE)
	s.AssertCanTestForEqualityWith(s.STRING_TYPE, s.BOOLEAN_OBJECT_TYPE)
	s.AssertCannotTestForEqualityWith(s.STRING_TYPE, s.SYMBOL_TYPE)
	s.AssertCannotTestForEqualityWith(s.STRING_TYPE, s.SYMBOL_OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.STRING_TYPE, s.DATE_TYPE)
	s.AssertCanTestForEqualityWith(s.STRING_TYPE, s.REGEXP_TYPE)
	s.AssertCanTestForEqualityWith(s.STRING_TYPE, s.ARRAY_TYPE)
	s.AssertCanTestForEqualityWith(s.STRING_TYPE, s.UNKNOWN_TYPE)

	// canTestForShallowEqualityWith
	s.AssertThat(s.STRING_TYPE.CanTestForShallowEqualityWith(s.NO_TYPE)).IsTrue()
	s.AssertThat(s.STRING_TYPE.CanTestForShallowEqualityWith(s.NO_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.STRING_TYPE.CanTestForShallowEqualityWith(s.ARRAY_TYPE)).IsFalse()
	s.AssertThat(s.STRING_TYPE.CanTestForShallowEqualityWith(s.BOOLEAN_TYPE)).IsFalse()
	s.AssertThat(s.STRING_TYPE.CanTestForShallowEqualityWith(s.BOOLEAN_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.STRING_TYPE.CanTestForShallowEqualityWith(s.DATE_TYPE)).IsFalse()
	s.AssertThat(s.STRING_TYPE.CanTestForShallowEqualityWith(s.functionType)).IsFalse()
	s.AssertThat(s.STRING_TYPE.CanTestForShallowEqualityWith(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.STRING_TYPE.CanTestForShallowEqualityWith(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.STRING_TYPE.CanTestForShallowEqualityWith(s.NUMBER_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.STRING_TYPE.CanTestForShallowEqualityWith(s.OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.STRING_TYPE.CanTestForShallowEqualityWith(s.REGEXP_TYPE)).IsFalse()
	s.AssertThat(s.STRING_TYPE.CanTestForShallowEqualityWith(s.STRING_TYPE)).IsTrue()
	s.AssertThat(s.STRING_TYPE.CanTestForShallowEqualityWith(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.STRING_TYPE.CanTestForShallowEqualityWith(s.SYMBOL_TYPE)).IsFalse()
	s.AssertThat(s.STRING_TYPE.CanTestForShallowEqualityWith(s.SYMBOL_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.STRING_TYPE.CanTestForShallowEqualityWith(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.STRING_TYPE.CanTestForShallowEqualityWith(s.VOID_TYPE)).IsFalse()
	s.AssertThat(s.STRING_TYPE.CanTestForShallowEqualityWith(s.UNKNOWN_TYPE)).IsTrue()

	// matchesXxx
	s.AssertThat(s.STRING_TYPE.MatchesNumberContext()).IsTrue()
	s.AssertThat(s.STRING_TYPE.MatchesObjectContext()).IsTrue()
	s.AssertThat(s.STRING_TYPE.MatchesStringContext()).IsTrue()
	s.AssertThat(s.STRING_TYPE.MatchesSymbolContext()).IsFalse()

	// isNullable
	s.AssertThat(s.STRING_TYPE.IsNullable()).IsFalse()
	s.AssertThat(s.STRING_TYPE.IsVoidable()).IsFalse()
	s.AssertThat(s.CreateNullableType(s.STRING_TYPE).IsNullable()).IsTrue()
	s.AssertThat(s.CreateUnionType(s.STRING_TYPE, s.VOID_TYPE).IsVoidable()).IsTrue()

	// toString
	s.AssertThat(s.STRING_TYPE.String()).IsEqualTo("string")
	s.AssertThat(s.STRING_TYPE.HasDisplayName()).IsTrue()
	s.AssertThat(s.STRING_TYPE.GetDisplayName()).IsEqualTo("string")

	// findPropertyType
	s.AssertTypeEquals(s.NUMBER_TYPE, s.STRING_TYPE.FindPropertyType("length"))
	s.AssertThat(s.STRING_TYPE.FindPropertyType("unknownProperty")).IsNull()

	rhinotesting.Asserts_AssertResolvesToSame(s.T, s.STRING_TYPE)
	s.AssertThat(s.STRING_TYPE.IsNominalConstructorOrInterface()).IsFalse()
}

/** Tests the behavior of the symbol type. */
func TestJSType_SymbolValueType(t *testing.T) {
	s := newJsTypeTest(t)
	// isXxx
	s.AssertThat(s.SYMBOL_TYPE.IsArrayType()).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.IsBooleanObjectType()).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.IsBooleanValueType()).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.IsDateType()).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.IsEnumElementType()).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.IsNamedType()).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.IsNullType()).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.IsNumberObjectType()).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.IsNumberValueType()).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.IsFunctionPrototypeType()).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.IsRegexpType()).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.IsStringObjectType()).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.IsStringValueType()).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.IsEnumType()).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.IsUnionType()).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.IsStruct()).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.IsDict()).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.IsAllType()).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.IsVoidType()).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.IsConstructor()).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.IsInstanceType()).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.IsSymbol()).IsTrue()
	s.AssertThat(s.SYMBOL_TYPE.IsSymbolObjectType()).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.IsSymbolValueType()).IsTrue()

	// autoboxesTo
	s.AssertTypeEquals(s.SYMBOL_OBJECT_TYPE, s.SYMBOL_TYPE.AutoboxesTo())

	// isSubtype
	s.AssertThat(s.SYMBOL_TYPE.IsSubtypeOf(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.SYMBOL_TYPE.IsSubtypeOf(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.IsSubtypeOf(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.IsSubtypeOf(s.functionType)).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.IsSubtypeOf(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.IsSubtypeOf(s.OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.IsSubtypeOf(s.DATE_TYPE)).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.IsSubtypeOf(s.unresolvedNamedType)).IsTrue()
	s.AssertThat(s.SYMBOL_TYPE.IsSubtypeOf(s.namedGoogBar)).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.IsSubtypeOf(s.REGEXP_TYPE)).IsFalse()

	// canBeCalled
	s.AssertThat(s.SYMBOL_TYPE.CanBeCalled()).IsFalse()

	// canTestForEqualityWith
	s.AssertCanTestForEqualityWith(s.SYMBOL_TYPE, s.SYMBOL_TYPE)
	s.AssertCanTestForEqualityWith(s.SYMBOL_TYPE, s.SYMBOL_OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.SYMBOL_TYPE, s.ALL_TYPE)
	s.AssertCannotTestForEqualityWith(s.SYMBOL_TYPE, s.STRING_OBJECT_TYPE)
	s.AssertCannotTestForEqualityWith(s.SYMBOL_TYPE, s.NUMBER_TYPE)
	s.AssertCannotTestForEqualityWith(s.SYMBOL_TYPE, s.functionType)
	s.AssertCannotTestForEqualityWith(s.SYMBOL_TYPE, s.VOID_TYPE)
	s.AssertCanTestForEqualityWith(s.SYMBOL_TYPE, s.OBJECT_TYPE)
	s.AssertCannotTestForEqualityWith(s.SYMBOL_TYPE, s.DATE_TYPE)
	s.AssertCannotTestForEqualityWith(s.SYMBOL_TYPE, s.REGEXP_TYPE)
	s.AssertCanTestForEqualityWith(s.SYMBOL_TYPE, s.UNKNOWN_TYPE)

	// canTestForShallowEqualityWith
	s.AssertThat(s.SYMBOL_TYPE.CanTestForShallowEqualityWith(s.SYMBOL_TYPE)).IsTrue()
	s.AssertThat(s.SYMBOL_TYPE.CanTestForShallowEqualityWith(s.SYMBOL_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.CanTestForShallowEqualityWith(s.NO_TYPE)).IsTrue()
	s.AssertThat(s.SYMBOL_TYPE.CanTestForShallowEqualityWith(s.NO_OBJECT_TYPE)).IsFalse() // ?
	s.AssertThat(s.SYMBOL_TYPE.CanTestForShallowEqualityWith(s.ARRAY_TYPE)).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.CanTestForShallowEqualityWith(s.SYMBOL_TYPE)).IsTrue()
	s.AssertThat(s.SYMBOL_TYPE.CanTestForShallowEqualityWith(s.BOOLEAN_TYPE)).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.CanTestForShallowEqualityWith(s.DATE_TYPE)).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.CanTestForShallowEqualityWith(s.functionType)).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.CanTestForShallowEqualityWith(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.CanTestForShallowEqualityWith(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.CanTestForShallowEqualityWith(s.NUMBER_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.CanTestForShallowEqualityWith(s.OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.CanTestForShallowEqualityWith(s.REGEXP_TYPE)).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.CanTestForShallowEqualityWith(s.STRING_TYPE)).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.CanTestForShallowEqualityWith(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.CanTestForShallowEqualityWith(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.SYMBOL_TYPE.CanTestForShallowEqualityWith(s.VOID_TYPE)).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.CanTestForShallowEqualityWith(s.UNKNOWN_TYPE)).IsTrue()

	// isNullable
	s.AssertThat(s.SYMBOL_TYPE.IsNullable()).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.IsVoidable()).IsFalse()

	// matchesXxx
	s.AssertThat(s.SYMBOL_TYPE.MatchesNumberContext()).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.MatchesObjectContext()).IsTrue()
	s.AssertThat(s.SYMBOL_TYPE.MatchesStringContext()).IsFalse()
	s.AssertThat(s.SYMBOL_TYPE.MatchesSymbolContext()).IsTrue()

	// toString
	s.AssertThat(s.SYMBOL_TYPE.String()).IsEqualTo("symbol")
	s.AssertThat(s.SYMBOL_TYPE.HasDisplayName()).IsTrue()
	s.AssertThat(s.SYMBOL_TYPE.GetDisplayName()).IsEqualTo("symbol")

	rhinotesting.Asserts_AssertResolvesToSame(s.T, s.SYMBOL_TYPE)
}

/** Tests the behavior of the Symbol type. */
func TestJSType_SymbolObjectType(t *testing.T) {
	s := newJsTypeTest(t)
	// isXxx
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.IsArrayType()).IsFalse()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.IsBooleanObjectType()).IsFalse()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.IsBooleanValueType()).IsFalse()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.IsDateType()).IsFalse()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.IsEnumElementType()).IsFalse()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.IsNamedType()).IsFalse()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.IsNullType()).IsFalse()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.IsNumberObjectType()).IsFalse()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.IsNumberValueType()).IsFalse()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.IsFunctionPrototypeType()).IsFalse()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.GetImplicitPrototype().IsFunctionPrototypeType()).IsTrue()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.IsRegexpType()).IsFalse()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.IsStringObjectType()).IsFalse()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.IsStringValueType()).IsFalse()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.IsEnumType()).IsFalse()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.IsUnionType()).IsFalse()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.IsStruct()).IsFalse()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.IsDict()).IsFalse()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.IsAllType()).IsFalse()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.IsVoidType()).IsFalse()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.IsConstructor()).IsFalse()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.IsInstanceType()).IsTrue()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.IsSymbol()).IsTrue()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.IsSymbolObjectType()).IsTrue()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.IsSymbolValueType()).IsFalse()

	// isSubtype
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.IsSubtypeOf(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.IsSubtypeOf(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.IsSubtypeOf(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.IsSubtypeOf(s.functionType)).IsFalse()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.IsSubtypeOf(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.IsSubtypeOf(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.IsSubtypeOf(s.DATE_TYPE)).IsFalse()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.IsSubtypeOf(s.unresolvedNamedType)).IsTrue()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.IsSubtypeOf(s.namedGoogBar)).IsFalse()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.IsSubtypeOf(s.REGEXP_TYPE)).IsFalse()
	// canBeCalled
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.CanBeCalled()).IsFalse()

	// canTestForEqualityWith
	s.AssertCanTestForEqualityWith(s.SYMBOL_OBJECT_TYPE, s.ALL_TYPE)
	s.AssertCannotTestForEqualityWith(s.SYMBOL_OBJECT_TYPE, s.STRING_OBJECT_TYPE)
	s.AssertCannotTestForEqualityWith(s.SYMBOL_OBJECT_TYPE, s.NUMBER_TYPE)
	s.AssertCannotTestForEqualityWith(s.SYMBOL_OBJECT_TYPE, s.functionType)
	s.AssertCannotTestForEqualityWith(s.SYMBOL_OBJECT_TYPE, s.VOID_TYPE)
	s.AssertCanTestForEqualityWith(s.SYMBOL_OBJECT_TYPE, s.OBJECT_TYPE)
	s.AssertCannotTestForEqualityWith(s.SYMBOL_OBJECT_TYPE, s.DATE_TYPE)
	s.AssertCannotTestForEqualityWith(s.SYMBOL_OBJECT_TYPE, s.REGEXP_TYPE)

	// canTestForShallowEqualityWith
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.CanTestForShallowEqualityWith(s.NO_TYPE)).IsTrue()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.CanTestForShallowEqualityWith(s.NO_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.CanTestForShallowEqualityWith(s.ARRAY_TYPE)).IsFalse()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.CanTestForShallowEqualityWith(s.SYMBOL_TYPE)).IsFalse()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.CanTestForShallowEqualityWith(s.SYMBOL_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.CanTestForShallowEqualityWith(s.DATE_TYPE)).IsFalse()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.CanTestForShallowEqualityWith(s.functionType)).IsFalse()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.CanTestForShallowEqualityWith(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.CanTestForShallowEqualityWith(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.CanTestForShallowEqualityWith(s.NUMBER_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.CanTestForShallowEqualityWith(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.CanTestForShallowEqualityWith(s.REGEXP_TYPE)).IsFalse()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.CanTestForShallowEqualityWith(s.STRING_TYPE)).IsFalse()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.CanTestForShallowEqualityWith(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.CanTestForShallowEqualityWith(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.CanTestForShallowEqualityWith(s.VOID_TYPE)).IsFalse()

	// isNullable
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.IsNullable()).IsFalse()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.IsVoidable()).IsFalse()

	// matchesXxx
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.MatchesNumberContext()).IsFalse()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.MatchesObjectContext()).IsTrue()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.MatchesStringContext()).IsFalse()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.MatchesSymbolContext()).IsTrue()

	// toString
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.String()).IsEqualTo("Symbol")
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.HasDisplayName()).IsTrue()
	s.AssertThat(s.SYMBOL_OBJECT_TYPE.GetDisplayName()).IsEqualTo("Symbol")

	s.AssertThat(s.SYMBOL_OBJECT_TYPE.IsNativeObjectType()).IsTrue()

	rhinotesting.Asserts_AssertResolvesToSame(s.T, s.SYMBOL_OBJECT_TYPE)
}

func (s *jsTypeTest) AssertPropertyTypeDeclared(ownerType ObjectType, prop string) {
	s.AssertThat(ownerType.IsPropertyTypeDeclared(prop)).IsTrue()
	s.AssertThat(ownerType.IsPropertyTypeInferred(prop)).IsFalse()
}

func (s *jsTypeTest) AssertPropertyTypeInferred(ownerType ObjectType, prop string) {
	s.AssertThat(ownerType.IsPropertyTypeDeclared(prop)).IsFalse()
	s.AssertThat(ownerType.IsPropertyTypeInferred(prop)).IsTrue()
}

func (s *jsTypeTest) AssertPropertyTypeUnknown(ownerType ObjectType, prop string) {
	s.AssertThat(ownerType.IsPropertyTypeDeclared(prop)).IsFalse()
	s.AssertThat(ownerType.IsPropertyTypeInferred(prop)).IsFalse()
	s.AssertThat(ownerType.GetPropertyType(prop).IsUnknownType()).IsTrue()
}

func (s *jsTypeTest) AssertReturnTypeEquals(expectedReturnType JSType, function JSType) {
	rhinotesting.IsInstanceOf[FunctionType](s.AssertThat(function))
	s.AssertTypeEquals(expectedReturnType,
		(function.(FunctionType)).GetReturnType())
}

/** Tests the behavior of record types. */
func TestJSType_RecordType(t *testing.T) {
	s := newJsTypeTest(t)
	// isXxx
	s.AssertThat(s.recordType.IsObject()).IsTrue()
	s.AssertThat(s.recordType.IsFunctionPrototypeType()).IsFalse()

	// isSubtype
	s.AssertThat(s.recordType.IsSubtypeOf(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.recordType.IsSubtypeOf(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.recordType.IsSubtypeOf(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.recordType.IsSubtypeOf(s.SYMBOL_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.recordType.IsSubtypeOf(s.SYMBOL_TYPE)).IsFalse()
	s.AssertThat(s.recordType.IsSubtypeOf(s.DATE_TYPE)).IsFalse()
	s.AssertThat(s.recordType.IsSubtypeOf(s.REGEXP_TYPE)).IsFalse()
	s.AssertThat(s.recordType.IsSubtypeOf(s.UNKNOWN_TYPE)).IsTrue()
	s.AssertThat(s.recordType.IsSubtypeOf(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.recordType.IsSubtypeOf(s.FUNCTION_TYPE)).IsFalse()

	// autoboxesTo
	s.AssertThat(s.recordType.AutoboxesTo()).IsNull()

	// canBeCalled
	s.AssertThat(s.recordType.CanBeCalled()).IsFalse()

	// canTestForEqualityWith
	s.AssertCanTestForEqualityWith(s.recordType, s.ALL_TYPE)
	s.AssertCanTestForEqualityWith(s.recordType, s.STRING_OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.recordType, s.recordType)
	s.AssertCanTestForEqualityWith(s.recordType, s.functionType)
	s.AssertCanTestForEqualityWith(s.recordType, s.OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.recordType, s.NUMBER_TYPE)
	s.AssertCanTestForEqualityWith(s.recordType, s.DATE_TYPE)
	s.AssertCanTestForEqualityWith(s.recordType, s.REGEXP_TYPE)
	s.AssertCanTestForEqualityWith(s.recordType, s.SYMBOL_OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.recordType, s.SYMBOL_TYPE)

	// canTestForShallowEqualityWith
	s.AssertThat(s.recordType.CanTestForShallowEqualityWith(s.NO_TYPE)).IsTrue()
	s.AssertThat(s.recordType.CanTestForShallowEqualityWith(s.NO_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.recordType.CanTestForShallowEqualityWith(s.ARRAY_TYPE)).IsFalse()
	s.AssertThat(s.recordType.CanTestForShallowEqualityWith(s.BOOLEAN_TYPE)).IsFalse()
	s.AssertThat(s.recordType.CanTestForShallowEqualityWith(s.BOOLEAN_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.recordType.CanTestForShallowEqualityWith(s.DATE_TYPE)).IsFalse()
	s.AssertThat(s.recordType.CanTestForShallowEqualityWith(s.recordType)).IsTrue()
	s.AssertThat(s.recordType.CanTestForShallowEqualityWith(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.recordType.CanTestForShallowEqualityWith(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.recordType.CanTestForShallowEqualityWith(s.NUMBER_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.recordType.CanTestForShallowEqualityWith(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.recordType.CanTestForShallowEqualityWith(s.REGEXP_TYPE)).IsFalse()
	s.AssertThat(s.recordType.CanTestForShallowEqualityWith(s.STRING_TYPE)).IsFalse()
	s.AssertThat(s.recordType.CanTestForShallowEqualityWith(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.recordType.CanTestForShallowEqualityWith(s.SYMBOL_TYPE)).IsFalse()
	s.AssertThat(s.recordType.CanTestForShallowEqualityWith(s.SYMBOL_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.recordType.CanTestForShallowEqualityWith(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.recordType.CanTestForShallowEqualityWith(s.VOID_TYPE)).IsFalse()
	s.AssertThat(s.recordType.CanTestForShallowEqualityWith(s.UNKNOWN_TYPE)).IsTrue()

	// matchesXxx
	s.AssertThat(s.recordType.MatchesNumberContext()).IsFalse()
	s.AssertThat(s.recordType.MatchesObjectContext()).IsTrue()
	s.AssertThat(s.recordType.MatchesStringContext()).IsFalse()
	s.AssertThat(s.recordType.MatchesSymbolContext()).IsFalse()

	rhinotesting.Asserts_AssertResolvesToSame(s.T, s.recordType)
}

/** Tests the behavior of the instance of Function. */
func TestJSType_FunctionInstanceType(t *testing.T) {
	s := newJsTypeTest(t)
	functionInst := s.FUNCTION_TYPE

	// isXxx
	s.AssertThat(functionInst.IsObject()).IsTrue()
	s.AssertThat(functionInst.IsFunctionPrototypeType()).IsFalse()
	s.AssertThat(functionInst.GetImplicitPrototype().IsFunctionPrototypeType()).IsTrue()

	// isSubtype
	s.AssertThat(functionInst.IsSubtype(s.ALL_TYPE)).IsTrue()
	s.AssertThat(functionInst.IsSubtype(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(functionInst.IsSubtype(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(functionInst.IsSubtype(s.SYMBOL_OBJECT_TYPE)).IsFalse()
	s.AssertThat(functionInst.IsSubtype(s.SYMBOL_TYPE)).IsFalse()
	s.AssertThat(functionInst.IsSubtype(s.DATE_TYPE)).IsFalse()
	s.AssertThat(functionInst.IsSubtype(s.REGEXP_TYPE)).IsFalse()
	s.AssertThat(functionInst.IsSubtype(s.UNKNOWN_TYPE)).IsTrue()
	s.AssertThat(functionInst.IsSubtype(s.FUNCTION_TYPE)).IsTrue()

	// autoboxesTo
	s.AssertThat(functionInst.AutoboxesTo()).IsNull()

	// canBeCalled
	s.AssertThat(functionInst.CanBeCalled()).IsTrue()

	// canTestForEqualityWith
	s.AssertCanTestForEqualityWith(functionInst, s.ALL_TYPE)
	s.AssertCanTestForEqualityWith(functionInst, s.STRING_OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(functionInst, functionInst)
	s.AssertCanTestForEqualityWith(functionInst, s.OBJECT_TYPE)
	s.AssertCannotTestForEqualityWith(functionInst, s.NUMBER_TYPE)
	s.AssertCannotTestForEqualityWith(functionInst, s.SYMBOL_OBJECT_TYPE)
	s.AssertCannotTestForEqualityWith(functionInst, s.SYMBOL_TYPE)
	s.AssertCanTestForEqualityWith(functionInst, s.DATE_TYPE)
	s.AssertCanTestForEqualityWith(functionInst, s.REGEXP_TYPE)

	// canTestForShallowEqualityWith
	s.AssertThat(functionInst.CanTestForShallowEqualityWith(s.NO_TYPE)).IsTrue()
	s.AssertThat(functionInst.CanTestForShallowEqualityWith(s.NO_OBJECT_TYPE)).IsTrue()
	s.AssertThat(functionInst.CanTestForShallowEqualityWith(s.ARRAY_TYPE)).IsFalse()
	s.AssertThat(functionInst.CanTestForShallowEqualityWith(s.BOOLEAN_TYPE)).IsFalse()
	s.AssertThat(functionInst.CanTestForShallowEqualityWith(s.BOOLEAN_OBJECT_TYPE)).IsFalse()
	s.AssertThat(functionInst.CanTestForShallowEqualityWith(s.DATE_TYPE)).IsFalse()
	s.AssertThat(functionInst.CanTestForShallowEqualityWith(functionInst)).IsTrue()
	s.AssertThat(functionInst.CanTestForShallowEqualityWith(s.NULL_TYPE)).IsFalse()
	s.AssertThat(functionInst.CanTestForShallowEqualityWith(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(functionInst.CanTestForShallowEqualityWith(s.NUMBER_OBJECT_TYPE)).IsFalse()
	s.AssertThat(functionInst.CanTestForShallowEqualityWith(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(functionInst.CanTestForShallowEqualityWith(s.REGEXP_TYPE)).IsFalse()
	s.AssertThat(functionInst.CanTestForShallowEqualityWith(s.STRING_TYPE)).IsFalse()
	s.AssertThat(functionInst.CanTestForShallowEqualityWith(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(functionInst.CanTestForShallowEqualityWith(s.SYMBOL_TYPE)).IsFalse()
	s.AssertThat(functionInst.CanTestForShallowEqualityWith(s.SYMBOL_OBJECT_TYPE)).IsFalse()
	s.AssertThat(functionInst.CanTestForShallowEqualityWith(s.ALL_TYPE)).IsTrue()
	s.AssertThat(functionInst.CanTestForShallowEqualityWith(s.VOID_TYPE)).IsFalse()
	s.AssertThat(functionInst.CanTestForShallowEqualityWith(s.UNKNOWN_TYPE)).IsTrue()

	// matchesXxx
	s.AssertThat(functionInst.MatchesNumberContext()).IsFalse()
	s.AssertThat(functionInst.MatchesObjectContext()).IsTrue()
	s.AssertThat(functionInst.MatchesStringContext()).IsFalse()
	s.AssertThat(functionInst.MatchesSymbolContext()).IsFalse()

	// hasProperty
	s.AssertThat(functionInst.HasProperty("prototype")).IsTrue()
	s.AssertPropertyTypeInferred(functionInst, "prototype")

	// misc
	s.AssertTypeEquals(s.FUNCTION_FUNCTION_TYPE, functionInst.GetConstructor())
	s.AssertTypeEquals(s.FUNCTION_PROTOTYPE, functionInst.GetImplicitPrototype())
	s.AssertTypeEquals(functionInst, s.FUNCTION_FUNCTION_TYPE.GetInstanceType())

	rhinotesting.Asserts_AssertResolvesToSame(s.T, functionInst)
}

/** Tests the behavior of functional types. */
func TestJSType_FunctionType(t *testing.T) {
	s := newJsTypeTest(t)
	// isXxx
	s.AssertThat(s.functionType.IsObject()).IsTrue()
	s.AssertThat(s.functionType.IsFunctionPrototypeType()).IsFalse()
	s.AssertThat(s.functionType.GetImplicitPrototype().IsFunctionPrototypeType()).IsTrue()

	// isSubtype
	s.AssertThat(s.functionType.IsSubtype(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.functionType.IsSubtype(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.functionType.IsSubtype(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.functionType.IsSubtype(s.SYMBOL_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.functionType.IsSubtype(s.SYMBOL_TYPE)).IsFalse()
	s.AssertThat(s.functionType.IsSubtype(s.DATE_TYPE)).IsFalse()
	s.AssertThat(s.functionType.IsSubtype(s.REGEXP_TYPE)).IsFalse()
	s.AssertThat(s.functionType.IsSubtype(s.UNKNOWN_TYPE)).IsTrue()
	s.AssertThat(s.functionType.IsSubtype(s.FUNCTION_TYPE)).IsTrue()

	// autoboxesTo
	s.AssertThat(s.functionType.AutoboxesTo()).IsNull()

	// canBeCalled
	s.AssertThat(s.functionType.CanBeCalled()).IsTrue()

	// canTestForEqualityWith
	s.AssertCanTestForEqualityWith(s.functionType, s.ALL_TYPE)
	s.AssertCanTestForEqualityWith(s.functionType, s.STRING_OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.functionType, s.functionType)
	s.AssertCanTestForEqualityWith(s.functionType, s.OBJECT_TYPE)
	s.AssertCannotTestForEqualityWith(s.functionType, s.NUMBER_TYPE)
	s.AssertCanTestForEqualityWith(s.functionType, s.DATE_TYPE)
	s.AssertCanTestForEqualityWith(s.functionType, s.REGEXP_TYPE)
	s.AssertCannotTestForEqualityWith(s.functionType, s.SYMBOL_OBJECT_TYPE)
	s.AssertCannotTestForEqualityWith(s.functionType, s.SYMBOL_TYPE)

	// canTestForShallowEqualityWith
	s.AssertThat(s.functionType.CanTestForShallowEqualityWith(s.NO_TYPE)).IsTrue()
	s.AssertThat(s.functionType.CanTestForShallowEqualityWith(s.NO_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.functionType.CanTestForShallowEqualityWith(s.ARRAY_TYPE)).IsFalse()
	s.AssertThat(s.functionType.CanTestForShallowEqualityWith(s.BOOLEAN_TYPE)).IsFalse()
	s.AssertThat(s.functionType.CanTestForShallowEqualityWith(s.BOOLEAN_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.functionType.CanTestForShallowEqualityWith(s.DATE_TYPE)).IsFalse()
	s.AssertThat(s.functionType.CanTestForShallowEqualityWith(s.functionType)).IsTrue()
	s.AssertThat(s.functionType.CanTestForShallowEqualityWith(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.functionType.CanTestForShallowEqualityWith(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.functionType.CanTestForShallowEqualityWith(s.NUMBER_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.functionType.CanTestForShallowEqualityWith(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.functionType.CanTestForShallowEqualityWith(s.REGEXP_TYPE)).IsFalse()
	s.AssertThat(s.functionType.CanTestForShallowEqualityWith(s.STRING_TYPE)).IsFalse()
	s.AssertThat(s.functionType.CanTestForShallowEqualityWith(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.functionType.CanTestForShallowEqualityWith(s.SYMBOL_TYPE)).IsFalse()
	s.AssertThat(s.functionType.CanTestForShallowEqualityWith(s.SYMBOL_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.functionType.CanTestForShallowEqualityWith(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.functionType.CanTestForShallowEqualityWith(s.VOID_TYPE)).IsFalse()
	s.AssertThat(s.functionType.CanTestForShallowEqualityWith(s.UNKNOWN_TYPE)).IsTrue()

	// matchesXxx
	s.AssertThat(s.functionType.MatchesNumberContext()).IsFalse()
	s.AssertThat(s.functionType.MatchesObjectContext()).IsTrue()
	s.AssertThat(s.functionType.MatchesStringContext()).IsFalse()
	s.AssertThat(s.functionType.MatchesSymbolContext()).IsFalse()

	// hasProperty
	s.AssertThat(s.functionType.HasProperty("prototype")).IsTrue()
	s.AssertPropertyTypeInferred(s.functionType, "prototype")

	rhinotesting.Asserts_AssertResolvesToSame(s.T, s.functionType)

	s.AssertThat(FunctionType_Builder_(s.Registry).WithName(ptr("aFunctionName")).Build().GetDisplayName()).
		IsEqualTo("aFunctionName")
}

/** Tests the subtyping relation of record types. */
func TestJSType_RecordTypeSubtyping(t *testing.T) {
	s := newJsTypeTest(t)
	builder := NewRecordTypeBuilder(s.Registry)
	builder.AddProperty("a", s.NUMBER_TYPE, nil)
	builder.AddProperty("b", s.STRING_TYPE, nil)
	builder.AddProperty("c", s.STRING_TYPE, nil)
	subRecordType := builder.Build()

	s.AssertThat(subRecordType.IsSubtypeOf(s.recordType)).IsTrue()
	s.AssertThat(s.recordType.IsSubtypeOf(subRecordType)).IsFalse()

	builder = NewRecordTypeBuilder(s.Registry)
	builder.AddProperty("a", s.OBJECT_TYPE, nil)
	builder.AddProperty("b", s.STRING_TYPE, nil)
	differentRecordType := builder.Build()

	s.AssertThat(differentRecordType.IsSubtypeOf(s.recordType)).IsFalse()
	s.AssertThat(s.recordType.IsSubtypeOf(differentRecordType)).IsFalse()
}

/** Tests the subtyping relation of record types when an object has an inferred property.. */
func TestJSType_RecordTypeSubtypingWithInferredProperties(t *testing.T) {
	s := newJsTypeTest(t)
	builder := NewRecordTypeBuilder(s.Registry)
	builder.AddProperty("a", s.googSubBarInst, nil)
	record := builder.Build()

	subtypeProp := s.Registry.CreateAnonymousObjectType(nil)
	subtypeProp.DefineInferredProperty("a", s.googSubSubBarInst, nil)
	s.AssertThat(subtypeProp.IsSubtypeOf(record)).IsTrue()
	s.AssertThat(record.IsSubtypeOf(subtypeProp)).IsFalse()

	supertypeProp := s.Registry.CreateAnonymousObjectType(nil)
	supertypeProp.DefineInferredProperty("a", s.googBarInst, nil)
	s.AssertThat(supertypeProp.IsSubtypeOf(record)).IsFalse()
	s.AssertThat(record.IsSubtypeOf(supertypeProp)).IsFalse()

	declaredSubtypeProp := s.Registry.CreateAnonymousObjectType(nil)
	declaredSubtypeProp.DefineDeclaredProperty("a", s.googSubSubBarInst, nil)
	s.AssertThat(declaredSubtypeProp.IsSubtypeOf(record)).IsTrue()
	s.AssertThat(record.IsSubtypeOf(declaredSubtypeProp)).IsFalse()
}

/** Tests the getLeastSupertype method for record types. */
func TestJSType_RecordTypeLeastSuperType1(t *testing.T) {
	s := newJsTypeTest(t)
	builder := NewRecordTypeBuilder(s.Registry)
	builder.AddProperty("a", s.NUMBER_TYPE, nil)
	builder.AddProperty("b", s.STRING_TYPE, nil)
	builder.AddProperty("c", s.STRING_TYPE, nil)
	subRecordType := builder.Build()

	leastSupertype := s.recordType.GetLeastSupertype(subRecordType)
	s.AssertTypeEquals(leastSupertype, s.Registry.CreateUnionType(s.recordType, subRecordType))
}

func TestJSType_RecordTypeLeastSuperType2(t *testing.T) {
	s := newJsTypeTest(t)
	builder := NewRecordTypeBuilder(s.Registry)
	builder.AddProperty("e", s.NUMBER_TYPE, nil)
	builder.AddProperty("b", s.STRING_TYPE, nil)
	builder.AddProperty("c", s.STRING_TYPE, nil)
	otherRecordType := builder.Build()

	s.AssertTypeEquals(
		s.Registry.CreateUnionType(s.recordType, otherRecordType),
		s.recordType.GetLeastSupertype(otherRecordType))
}

func TestJSType_RecordTypeLeastSuperType3(t *testing.T) {
	s := newJsTypeTest(t)
	builder := NewRecordTypeBuilder(s.Registry)
	builder.AddProperty("d", s.NUMBER_TYPE, nil)
	builder.AddProperty("e", s.STRING_TYPE, nil)
	builder.AddProperty("f", s.STRING_TYPE, nil)
	otherRecordType := builder.Build()

	s.AssertTypeEquals(
		s.Registry.CreateUnionType(s.recordType, otherRecordType),
		s.recordType.GetLeastSupertype(otherRecordType))
}

func TestJSType_RecordTypeLeastSuperType4(t *testing.T) {
	s := newJsTypeTest(t)
	leastSupertype := s.recordType.GetLeastSupertype(s.OBJECT_TYPE)
	s.AssertTypeEquals(leastSupertype, s.OBJECT_TYPE)
}

/** Tests the getGreatestSubtype method for record types. */
func TestJSType_RecordTypeGreatestSubType1(t *testing.T) {
	s := newJsTypeTest(t)
	builder := NewRecordTypeBuilder(s.Registry)
	builder.AddProperty("d", s.NUMBER_TYPE, nil)
	builder.AddProperty("e", s.STRING_TYPE, nil)
	builder.AddProperty("f", s.STRING_TYPE, nil)

	subRecordType := builder.Build()

	subtype := s.recordType.GetGreatestSubtype(subRecordType)

	builder = NewRecordTypeBuilder(s.Registry)
	builder.AddProperty("d", s.NUMBER_TYPE, nil)
	builder.AddProperty("e", s.STRING_TYPE, nil)
	builder.AddProperty("f", s.STRING_TYPE, nil)
	builder.AddProperty("a", s.NUMBER_TYPE, nil)
	builder.AddProperty("b", s.STRING_TYPE, nil)

	s.AssertTypeEquals(subtype, builder.Build())
}

func TestJSType_RecordTypeGreatestSubType2(t *testing.T) {
	s := newJsTypeTest(t)
	builder := NewRecordTypeBuilder(s.Registry)

	subRecordType := builder.Build()

	subtype := s.recordType.GetGreatestSubtype(subRecordType)

	builder = NewRecordTypeBuilder(s.Registry)
	builder.AddProperty("a", s.NUMBER_TYPE, nil)
	builder.AddProperty("b", s.STRING_TYPE, nil)

	s.AssertTypeEquals(subtype, builder.Build())
}

func TestJSType_RecordTypeGreatestSubType3(t *testing.T) {
	s := newJsTypeTest(t)
	builder := NewRecordTypeBuilder(s.Registry)
	builder.AddProperty("a", s.NUMBER_TYPE, nil)
	builder.AddProperty("b", s.STRING_TYPE, nil)
	builder.AddProperty("c", s.STRING_TYPE, nil)

	subRecordType := builder.Build()

	subtype := s.recordType.GetGreatestSubtype(subRecordType)

	builder = NewRecordTypeBuilder(s.Registry)
	builder.AddProperty("a", s.NUMBER_TYPE, nil)
	builder.AddProperty("b", s.STRING_TYPE, nil)
	builder.AddProperty("c", s.STRING_TYPE, nil)

	s.AssertTypeEquals(subtype, builder.Build())
}

func TestJSType_RecordTypeGreatestSubType4(t *testing.T) {
	s := newJsTypeTest(t)
	builder := NewRecordTypeBuilder(s.Registry)
	builder.AddProperty("a", s.STRING_TYPE, nil)
	builder.AddProperty("b", s.STRING_TYPE, nil)
	builder.AddProperty("c", s.STRING_TYPE, nil)

	subRecordType := builder.Build()

	subtype := s.recordType.GetGreatestSubtype(subRecordType)
	s.AssertTypeEquals(subtype, s.NO_TYPE)
}

func TestJSType_RecordTypeGreatestSubType5(t *testing.T) {
	s := newJsTypeTest(t)
	builder := NewRecordTypeBuilder(s.Registry)
	builder.AddProperty("a", s.STRING_TYPE, nil)

	recordType := builder.Build()

	s.AssertTypeEquals(s.NO_OBJECT_TYPE, recordType.GetGreatestSubtype(s.FUNCTION_TYPE))

	// if Function is given a property "a" of type "string", then it's
	// a subtype of the record type {a: string}.
	s.FUNCTION_TYPE.DefineDeclaredProperty("a", s.STRING_TYPE, nil)
	s.AssertTypeEquals(s.FUNCTION_TYPE, recordType.GetGreatestSubtype(s.FUNCTION_TYPE))
	s.AssertTypeEquals(s.FUNCTION_TYPE, s.FUNCTION_TYPE.GetGreatestSubtype(recordType))
}

func TestJSType_RecordTypeGreatestSubType6(t *testing.T) {
	s := newJsTypeTest(t)
	builder := NewRecordTypeBuilder(s.Registry)
	builder.AddProperty("x", s.UNKNOWN_TYPE, nil)

	recordType := builder.Build()

	s.AssertTypeEquals(s.NO_OBJECT_TYPE, recordType.GetGreatestSubtype(s.FUNCTION_TYPE))

	// if Function is given a property "x" of type "string", then it's
	// also a subtype of the record type {x: ?}.
	s.FUNCTION_TYPE.DefineDeclaredProperty("x", s.STRING_TYPE, nil)
	s.AssertTypeEquals(s.FUNCTION_TYPE, recordType.GetGreatestSubtype(s.FUNCTION_TYPE))
	s.AssertTypeEquals(s.FUNCTION_TYPE, s.FUNCTION_TYPE.GetGreatestSubtype(recordType))
}

func TestJSType_RecordTypeGreatestSubType7(t *testing.T) {
	s := newJsTypeTest(t)
	builder := NewRecordTypeBuilder(s.Registry)
	builder.AddProperty("x", s.NUMBER_TYPE, nil)

	recordType := builder.Build()

	// if Function is given a property "x" of type "string", then it's
	// not a subtype of the record type {x: number}.
	s.FUNCTION_TYPE.DefineDeclaredProperty("x", s.STRING_TYPE, nil)
	s.AssertTypeEquals(s.NO_OBJECT_TYPE, recordType.GetGreatestSubtype(s.FUNCTION_TYPE))
}

func TestJSType_RecordTypeGreatestSubType8(t *testing.T) {
	s := newJsTypeTest(t)
	builder := NewRecordTypeBuilder(s.Registry)
	builder.AddProperty("xyz", s.UNKNOWN_TYPE, nil)

	recordType := builder.Build()

	s.AssertTypeEquals(s.NO_OBJECT_TYPE, recordType.GetGreatestSubtype(s.FUNCTION_TYPE))

	// if goog.Bar is given a property "xyz" of type "string", then it's
	// also a subtype of the record type {x: ?}.
	s.googBar.DefineDeclaredProperty("xyz", s.STRING_TYPE, nil)

	s.AssertTypeEquals(s.googBar, recordType.GetGreatestSubtype(s.FUNCTION_TYPE))
	s.AssertTypeEquals(s.googBar, s.FUNCTION_TYPE.GetGreatestSubtype(recordType))

	googBarInst := s.googBar.GetInstanceType()
	s.AssertTypeEquals(s.NO_OBJECT_TYPE,
		recordType.GetGreatestSubtype(googBarInst))
	s.AssertTypeEquals(s.NO_OBJECT_TYPE,
		googBarInst.GetGreatestSubtype(recordType))
}

func TestJSType_RecordTypeGreatestSubType9(t *testing.T) {
	s := newJsTypeTest(t)
	builder := NewRecordTypeBuilder(s.Registry)
	builder.AddProperty("d", s.googBar.GetPrototype(), nil)
	builder.AddProperty("e", s.STRING_TYPE, nil)
	builder.AddProperty("f", s.STRING_TYPE, nil)

	recordType1 := builder.Build()

	builder = NewRecordTypeBuilder(s.Registry)
	builder.AddProperty("d", s.googBar.GetInstanceType(), nil)
	builder.AddProperty("e", s.STRING_TYPE, nil)
	builder.AddProperty("f", s.STRING_TYPE, nil)

	recordType2 := builder.Build()

	subtype := recordType1.GetGreatestSubtype(recordType2)

	builder = NewRecordTypeBuilder(s.Registry)
	builder.AddProperty("d", s.googBar.GetInstanceType(), nil)
	builder.AddProperty("e", s.STRING_TYPE, nil)
	builder.AddProperty("f", s.STRING_TYPE, nil)

	s.AssertTypeEquals(subtype, builder.Build())
}

func TestJSType_RecordTypeGreatestSubType10(t *testing.T) {
	s := newJsTypeTest(t)
	builder := NewRecordTypeBuilder(s.Registry)
	builder.AddProperty("d", s.googBar.GetPrototype(), nil)
	builder.AddProperty("e", s.STRING_TYPE, nil)

	recordType1 := builder.Build()

	builder = NewRecordTypeBuilder(s.Registry)
	builder.AddProperty("d", s.googBar.GetInstanceType(), nil)
	builder.AddProperty("e", s.STRING_TYPE, nil)
	builder.AddProperty("f", s.STRING_TYPE, nil)

	recordType2 := builder.Build()

	subtype := recordType2.GetGreatestSubtype(recordType1)

	builder = NewRecordTypeBuilder(s.Registry)
	builder.AddProperty("d", s.googBar.GetInstanceType(), nil)
	builder.AddProperty("e", s.STRING_TYPE, nil)
	builder.AddProperty("f", s.STRING_TYPE, nil)

	s.AssertTypeEquals(subtype, builder.Build())
}

func TestJSType_RecordTypeGreatestSubType11(t *testing.T) {
	s := newJsTypeTest(t)
	builder := NewRecordTypeBuilder(s.Registry)
	builder.AddProperty("d", s.CreateUnionType(s.NUMBER_TYPE, s.STRING_TYPE), nil)
	builder.AddProperty("e", s.STRING_TYPE, nil)

	recordType1 := builder.Build()

	builder = NewRecordTypeBuilder(s.Registry)
	builder.AddProperty("d", s.CreateUnionType(s.NUMBER_TYPE, s.BOOLEAN_TYPE), nil)
	builder.AddProperty("e", s.STRING_TYPE, nil)
	builder.AddProperty("f", s.STRING_TYPE, nil)

	recordType2 := builder.Build()

	subtype := recordType2.GetGreatestSubtype(recordType1)

	builder = NewRecordTypeBuilder(s.Registry)
	builder.AddProperty("d", s.NUMBER_TYPE, nil)
	builder.AddProperty("e", s.STRING_TYPE, nil)
	builder.AddProperty("f", s.STRING_TYPE, nil)

	s.AssertTypeEquals(subtype, builder.Build())
}

func TestJSType_RecordTypeGreatestSubType12(t *testing.T) {
	s := newJsTypeTest(t)
	builder := NewRecordTypeBuilder(s.Registry)
	builder.AddProperty("d", s.googBar.GetPrototype(), nil)
	builder.AddProperty("e", s.STRING_TYPE, nil)

	recordType1 := builder.Build()

	googBarArgConstructor :=
		withOpenRegistry(s, func() FunctionType {
			return s.Registry.CreateConstructorType(
				ptr("barArg"), nil, s.Registry.CreateParameters(s.googBar), nil, nil, false)
		})

	builder = NewRecordTypeBuilder(s.Registry)
	builder.AddProperty("d", googBarArgConstructor, nil)
	builder.AddProperty("e", s.STRING_TYPE, nil)
	builder.AddProperty("f", s.STRING_TYPE, nil)

	recordType2 := builder.Build()

	subtype := recordType2.GetGreatestSubtype(recordType1)

	builder = NewRecordTypeBuilder(s.Registry)
	builder.AddProperty("d", s.Registry.GetNativeObjectType(
		JSTypeNative_NO_OBJECT_TYPE), nil)
	builder.AddProperty("e", s.STRING_TYPE, nil)
	builder.AddProperty("f", s.STRING_TYPE, nil)

	s.AssertTypeEquals(subtype, builder.Build())
}

/** Tests the "apply" method on the function type. */
func TestJSType_ApplyOfDateMethod(t *testing.T) {
	s := newJsTypeTest(t)
	applyType := s.dateMethod.GetPropertyType("apply")
	s.AssertWithMessage("apply should be a function").
		That(instanceOf[FunctionType](applyType)).
		IsTrue()

	applyFn := applyType.(FunctionType)
	s.AssertTypeEqualsMsg("apply should have the same return type as its function",
		s.NUMBER_TYPE, applyFn.GetReturnType())

	params := applyFn.GetParameters()
	s.AssertWithMessage("apply takes two args").That(params).HasSize(2)
	s.AssertTypeEqualsMsg(
		"apply's first arg is the @this type",
		s.Registry.CreateOptionalNullableType(s.DATE_TYPE),
		params[0].GetJSType())
	s.AssertTypeEqualsMsg(
		"apply's second arg is an Array",
		s.Registry.CreateOptionalNullableType(s.OBJECT_TYPE),
		params[1].GetJSType())
	s.AssertWithMessage("apply's args must be optional").That(params[0].IsOptional()).IsTrue()
	s.AssertWithMessage("apply's args must be optional").That(params[1].IsOptional()).IsTrue()
}

/** Tests the "call" method on the function type. */
func TestJSType_CallOfDateMethod(t *testing.T) {
	s := newJsTypeTest(t)
	callType := s.dateMethod.GetPropertyType("call")
	s.AssertWithMessage("call should be a function").That(instanceOf[FunctionType](callType)).IsTrue()

	callFn := callType.(FunctionType)
	s.AssertTypeEqualsMsg("call should have the same return type as its function",
		s.NUMBER_TYPE, callFn.GetReturnType())

	params := callFn.GetParameters()
	s.AssertWithMessage("call takes one argument in this case").That(params).HasSize(1)
	s.AssertTypeEqualsMsg(
		"call's first arg is the @this type",
		s.Registry.CreateOptionalNullableType(s.DATE_TYPE),
		params[0].GetJSType())
	s.AssertWithMessage("call's args must be optional").That(params[0].IsOptional()).IsTrue()
}

/** Tests the representation of function types. */
func TestJSType_FunctionTypeRepresentation(t *testing.T) {
	s := newJsTypeTest(t)
	s.AssertThat(s.Registry.CreateFunctionType(s.BOOLEAN_TYPE, s.NUMBER_TYPE, s.STRING_TYPE).String()).
		IsEqualTo("function(number, string): boolean")

	s.AssertThat(s.ARRAY_FUNCTION_TYPE.String()).IsEqualTo("function(new:Array, ...*): Array")

	s.AssertThat(s.BOOLEAN_OBJECT_FUNCTION_TYPE.String()).
		IsEqualTo("function(new:Boolean, *=): boolean")

	s.AssertThat(s.NUMBER_OBJECT_FUNCTION_TYPE.String()).
		IsEqualTo("function(new:Number, *=): number")

	s.AssertThat(s.STRING_OBJECT_FUNCTION_TYPE.String()).
		IsEqualTo("function(new:String, *=): string")

	s.AssertThat(s.Registry.CreateFunctionTypeWithVarArgs(s.BOOLEAN_TYPE, s.NUMBER_TYPE).String()).
		IsEqualTo("function(...number): boolean")

	s.AssertThat(
		s.Registry.
			CreateFunctionTypeWithVarArgs(s.BOOLEAN_TYPE, s.NUMBER_TYPE, s.STRING_TYPE).
			String()).
		IsEqualTo("function(number, ...string): boolean")

	s.AssertThat(
		FunctionType_Builder_(s.Registry).
			WithParametersOfList(s.Registry.CreateParameters(s.NUMBER_TYPE)).
			WithReturnType(s.NUMBER_STRING_BOOLEAN).
			WithTypeOfThis(s.DATE_TYPE).
			Build().
			String()).
		IsEqualTo("function(this:Date, number): (boolean|number|string)")
}

/** Tests relationships between structural function types. */
func TestJSType_FunctionTypeRelationships(t *testing.T) {
	s := newJsTypeTest(t)
	dateMethodEmpty :=
		FunctionType_Builder_(s.Registry).
			WithParametersOfList(s.Registry.CreateParameters()).
			WithTypeOfThis(s.DATE_TYPE).
			Build()
	dateMethodWithParam :=
		FunctionType_Builder_(s.Registry).
			WithParametersOfList(s.Registry.CreateOptionalParameters(s.NUMBER_TYPE)).
			WithTypeOfThis(s.DATE_TYPE).
			Build()
	dateMethodWithReturn :=
		FunctionType_Builder_(s.Registry).
			WithReturnType(s.NUMBER_TYPE).
			WithTypeOfThis(s.DATE_TYPE).
			Build()
	stringMethodEmpty :=
		FunctionType_Builder_(s.Registry).
			WithParametersOfList(s.Registry.CreateParameters()).
			WithTypeOfThis(s.STRING_OBJECT_TYPE).
			Build()
	stringMethodWithParam :=
		FunctionType_Builder_(s.Registry).
			WithParametersOfList(s.Registry.CreateOptionalParameters(s.NUMBER_TYPE)).
			WithTypeOfThis(s.STRING_OBJECT_TYPE).
			Build()
	stringMethodWithReturn :=
		FunctionType_Builder_(s.Registry).
			WithReturnType(s.NUMBER_TYPE).
			WithTypeOfThis(s.STRING_OBJECT_TYPE).
			Build()

	// One-off tests.
	s.AssertThat(stringMethodEmpty.IsSubtype(dateMethodEmpty)).IsFalse()

	// Systemic tests.
	allFunctions :=
		[]FunctionType{
			dateMethodEmpty,
			dateMethodWithParam,
			dateMethodWithReturn,
			stringMethodEmpty,
			stringMethodWithParam,
			stringMethodWithReturn}
	for i := 0; i < len(allFunctions); i++ {
		for j := 0; j < len(allFunctions); j++ {
			typeA := allFunctions[i]
			typeB := allFunctions[j]
			s.AssertWithMessage("equals(%s, %s)", typeA, typeB).
				That(typeA.Equals(typeB)).
				IsEqualTo(i == j)

			// For this particular set of functions, the functions are subtypes
			// of each other iff they have the same "this" type.
			s.AssertWithMessage("isSubtype(%s, %s)", typeA, typeB).
				That(typeA.IsSubtype(typeB)).
				IsEqualTo(typeA.GetTypeOfThis().Equals(typeB.GetTypeOfThis()))

			if i == j {
				s.AssertTypeEquals(typeA, typeA.GetLeastSupertype(typeB))
				s.AssertTypeEquals(typeA, typeA.GetGreatestSubtype(typeB))
			} else {
				s.AssertTypeEqualsMsg(
					fmt.Sprintf("sup(%s, %s)", typeA, typeB),
					s.FUNCTION_TYPE,
					typeA.GetLeastSupertype(typeB))
				s.AssertTypeEqualsMsg(fmt.Sprintf("inf(%s, %s)", typeA, typeB),
					s.LEAST_FUNCTION_TYPE, typeA.GetGreatestSubtype(typeB))
			}
		}
	}
}

func TestJSType_ProxiedFunctionTypeRelationships(t *testing.T) {
	s := newJsTypeTest(t)
	dateMethodEmpty :=
		FunctionType_Builder_(s.Registry).
			WithParametersOfList(s.Registry.CreateParameters()).
			WithTypeOfThis(s.DATE_TYPE).
			Build().
			ToMaybeFunctionType()
	dateMethodWithParam :=
		FunctionType_Builder_(s.Registry).
			WithParametersOfList(s.Registry.CreateParameters(s.NUMBER_TYPE)).
			WithTypeOfThis(s.DATE_TYPE).
			Build().
			ToMaybeFunctionType()
	proxyDateMethodEmpty :=
		NewProxyObjectType(s.Registry, dateMethodEmpty)
	proxyDateMethodWithParam :=
		NewProxyObjectType(s.Registry, dateMethodWithParam)

	s.AssertTypeEquals(
		s.FUNCTION_TYPE, proxyDateMethodEmpty.GetLeastSupertype(proxyDateMethodWithParam))
	s.AssertTypeEquals(s.LEAST_FUNCTION_TYPE,
		proxyDateMethodEmpty.GetGreatestSubtype(proxyDateMethodWithParam))
}

/** Tests relationships between structural function types. */
func TestJSType_FunctionSubTypeRelationships(t *testing.T) {
	s := newJsTypeTest(t)
	googBarMethod := FunctionType_Builder_(s.Registry).WithTypeOfThis(s.googBar).Build()
	googBarParamFn :=
		FunctionType_Builder_(s.Registry).WithParametersOfList(s.Registry.CreateParameters(s.googBar)).Build()
	googBarReturnFn :=
		FunctionType_Builder_(s.Registry).
			WithParametersOfList(s.Registry.CreateParameters()).
			WithReturnType(s.googBar).
			Build()
	googSubBarMethod :=
		FunctionType_Builder_(s.Registry).WithTypeOfThis(s.googSubBar).Build()
	googSubBarParamFn :=
		FunctionType_Builder_(s.Registry).
			WithParametersOfList(s.Registry.CreateParameters(s.googSubBar)).
			Build()
	googSubBarReturnFn :=
		FunctionType_Builder_(s.Registry).WithReturnType(s.googSubBar).Build()

	s.AssertThat(googBarMethod.IsSubtype(googSubBarMethod)).IsTrue()
	s.AssertThat(googBarReturnFn.IsSubtype(googSubBarReturnFn)).IsTrue()

	allFunctions :=
		[]JSType{
			googBarMethod,
			googBarParamFn,
			googBarReturnFn,
			googSubBarMethod,
			googSubBarParamFn,
			googSubBarReturnFn}
	for i := 0; i < len(allFunctions); i++ {
		for j := 0; j < len(allFunctions); j++ {
			typeA := allFunctions[i]
			typeB := allFunctions[j]
			s.AssertWithMessage("equals(%s, %s)", typeA, typeB).
				That(typeA.Equals(typeB)).
				IsEqualTo(i == j)

			// TODO(nicksantos): This formulation of least subtype and greatest
			// supertype is a bit loose. We might want to tighten it up later.
			if i == j {
				s.AssertTypeEquals(typeA, typeA.GetLeastSupertype(typeB))
				s.AssertTypeEquals(typeA, typeA.GetGreatestSubtype(typeB))
			} else {
				s.AssertTypeEqualsMsg(
					fmt.Sprintf("sup(%s, %s)", typeA, typeB),
					s.FUNCTION_TYPE,
					typeA.GetLeastSupertype(typeB))
				s.AssertTypeEqualsMsg(fmt.Sprintf("inf(%s, %s)", typeA, typeB),
					s.LEAST_FUNCTION_TYPE, typeA.GetGreatestSubtype(typeB))
			}
		}
	}
}

/**
* Tests that defining a property of a function's {@code prototype} adds the property to it
* instance type.
 */
func TestJSType_FunctionPrototypeAndImplicitPrototype1(t *testing.T) {
	s := newJsTypeTest(t)

	constructor :=
		withOpenRegistry(s, func() FunctionType {
			return s.Registry.CreateConstructorType(
				ptr("Foo") /* source= */, nil, nil, nil, nil /* isAbstract= */, false)
		})
	instance := constructor.GetInstanceType()

	// adding one property on the prototype
	prototype :=
		constructor.GetPropertyType("prototype").(ObjectType)
	prototype.DefineDeclaredProperty("foo", s.DATE_TYPE, nil)

	s.AssertThat(instance.GetPropertiesCount()).IsEqualTo(s.NATIVE_PROPERTIES_COUNT + 1)
}

/**
* Tests that replacing a function's {@code prototype} changes the visible properties of its
* instance type.
 */
func TestJSType_FunctionPrototypeAndImplicitPrototype2(t *testing.T) {
	s := newJsTypeTest(t)

	constructor :=
		withOpenRegistry(s, func() FunctionType {
			return s.Registry.CreateConstructorType(
				ptr("Bar"),
				/* source= */ nil,
				s.Registry.CreateParameters(s.ALL_TYPE, s.ALL_TYPE, s.ALL_TYPE),
				nil,
				nil,
				/* isAbstract= */ false)
		})
	instance := constructor.GetInstanceType()

	// replacing the prototype
	prototype := s.Registry.CreateAnonymousObjectType(nil)
	prototype.DefineDeclaredProperty("foo", s.DATE_TYPE, nil)
	constructor.DefineDeclaredProperty("prototype", prototype, nil)

	s.AssertThat(instance.GetPropertiesCount()).IsEqualTo(s.NATIVE_PROPERTIES_COUNT + 1)
}

/** Tests assigning JsDoc on a prototype property. */
func TestJSType_JSDocOnPrototypeProperty(t *testing.T) {
	s := newJsTypeTest(t)
	s.subclassCtor.SetPropertyJSDocInfo("prototype", JSDocInfo_Builder_().Build())
	s.AssertThat(s.subclassCtor.GetOwnPropertyJSDocInfo("prototype")).IsNull()
}

/** Tests operation of {@code isVoidable}. */
func TestJSType_IsVoidable(t *testing.T) {
	s := newJsTypeTest(t)
	s.AssertThat(s.VOID_TYPE.IsVoidable()).IsTrue()
	s.AssertThat(s.NULL_TYPE.IsVoidable()).IsFalse()
	s.AssertThat(s.CreateUnionType(s.NUMBER_TYPE, s.VOID_TYPE).IsVoidable()).IsTrue()
}

/** Tests operation of {@code isBigIntOrNumber}. */
func TestJSType_IsBigIntOrNumber(t *testing.T) {
	s := newJsTypeTest(t)
	s.AssertThat(s.BIGINT_TYPE.IsBigIntOrNumber()).IsTrue()
	s.AssertThat(s.BIGINT_OBJECT_TYPE.IsBigIntOrNumber()).IsTrue()
	// testing {bigint,number}
	s.AssertThat(s.BIGINT_NUMBER.IsBigIntOrNumber()).IsTrue()
	// testing {BigInt,Number}
	s.AssertThat(s.BIGINT_NUMBER_OBJECT.IsBigIntOrNumber()).IsTrue()
	s.AssertThat(s.BIGINT_NUMBER_STRING.IsBigIntOrNumber()).IsFalse()
	// This check is designed to return false when mixing object and primitive types, e.g.
	// (bigint|BigInt). Even though it should technically be true, these situation are usually a
	// result of user error, so we want to report them when encountered.
	s.AssertThat(s.CreateUnionType(s.BIGINT_TYPE, s.BIGINT_OBJECT_TYPE).IsBigIntOrNumber()).IsFalse()
}

/** Tests the behavior of the void type. */
func TestJSType_VoidType(t *testing.T) {
	s := newJsTypeTest(t)
	// isSubtype
	s.AssertThat(s.VOID_TYPE.IsSubtypeOf(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.VOID_TYPE.IsSubtypeOf(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.VOID_TYPE.IsSubtypeOf(s.REGEXP_TYPE)).IsFalse()

	// autoboxesTo
	s.AssertThat(s.VOID_TYPE.AutoboxesTo()).IsNull()

	// canTestForEqualityWith
	s.AssertCanTestForEqualityWith(s.VOID_TYPE, s.ALL_TYPE)
	s.AssertCannotTestForEqualityWith(s.VOID_TYPE, s.REGEXP_TYPE)

	// canTestForShallowEqualityWith
	s.AssertThat(s.VOID_TYPE.CanTestForShallowEqualityWith(s.NO_TYPE)).IsTrue()
	s.AssertThat(s.VOID_TYPE.CanTestForShallowEqualityWith(s.NO_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.VOID_TYPE.CanTestForShallowEqualityWith(s.ARRAY_TYPE)).IsFalse()
	s.AssertThat(s.VOID_TYPE.CanTestForShallowEqualityWith(s.BOOLEAN_TYPE)).IsFalse()
	s.AssertThat(s.VOID_TYPE.CanTestForShallowEqualityWith(s.BOOLEAN_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.VOID_TYPE.CanTestForShallowEqualityWith(s.DATE_TYPE)).IsFalse()
	s.AssertThat(s.VOID_TYPE.CanTestForShallowEqualityWith(s.functionType)).IsFalse()
	s.AssertThat(s.VOID_TYPE.CanTestForShallowEqualityWith(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.VOID_TYPE.CanTestForShallowEqualityWith(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.VOID_TYPE.CanTestForShallowEqualityWith(s.NUMBER_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.VOID_TYPE.CanTestForShallowEqualityWith(s.OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.VOID_TYPE.CanTestForShallowEqualityWith(s.REGEXP_TYPE)).IsFalse()
	s.AssertThat(s.VOID_TYPE.CanTestForShallowEqualityWith(s.STRING_TYPE)).IsFalse()
	s.AssertThat(s.VOID_TYPE.CanTestForShallowEqualityWith(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.VOID_TYPE.CanTestForShallowEqualityWith(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.VOID_TYPE.CanTestForShallowEqualityWith(s.VOID_TYPE)).IsTrue()
	s.AssertThat(s.VOID_TYPE.CanTestForShallowEqualityWith(s.CreateUnionType(s.NUMBER_TYPE, s.VOID_TYPE))).
		IsTrue()

	// matchesXxx
	s.AssertThat(s.VOID_TYPE.MatchesNumberContext()).IsFalse()
	s.AssertThat(s.VOID_TYPE.MatchesNumberContext()).IsFalse()
	s.AssertThat(s.VOID_TYPE.MatchesObjectContext()).IsFalse()
	s.AssertThat(s.VOID_TYPE.MatchesStringContext()).IsTrue()
	s.AssertThat(s.VOID_TYPE.MatchesNumberContext()).IsFalse()

	rhinotesting.Asserts_AssertResolvesToSame(s.T, s.VOID_TYPE)
}

/** Tests the behavior of the boolean type. */
func TestJSType_BooleanValueType(t *testing.T) {
	s := newJsTypeTest(t)
	// isXxx
	s.AssertThat(s.BOOLEAN_TYPE.IsArrayType()).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.IsBooleanObjectType()).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.IsBooleanValueType()).IsTrue()
	s.AssertThat(s.BOOLEAN_TYPE.IsDateType()).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.IsEnumElementType()).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.IsNamedType()).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.IsNullType()).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.IsNumberObjectType()).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.IsNumberValueType()).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.IsFunctionPrototypeType()).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.IsRegexpType()).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.IsStringObjectType()).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.IsStringValueType()).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.IsEnumType()).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.IsUnionType()).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.IsStruct()).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.IsDict()).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.IsAllType()).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.IsVoidType()).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.IsConstructor()).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.IsInstanceType()).IsFalse()

	// autoboxesTo
	s.AssertTypeEquals(s.BOOLEAN_OBJECT_TYPE, s.BOOLEAN_TYPE.AutoboxesTo())

	// isSubtype
	s.AssertThat(s.BOOLEAN_TYPE.IsSubtypeOf(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.BOOLEAN_TYPE.IsSubtypeOf(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.IsSubtypeOf(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.IsSubtypeOf(s.functionType)).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.IsSubtypeOf(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.IsSubtypeOf(s.OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.IsSubtypeOf(s.DATE_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.IsSubtypeOf(s.unresolvedNamedType)).IsTrue()
	s.AssertThat(s.BOOLEAN_TYPE.IsSubtypeOf(s.namedGoogBar)).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.IsSubtypeOf(s.REGEXP_TYPE)).IsFalse()

	// canBeCalled
	s.AssertThat(s.BOOLEAN_TYPE.CanBeCalled()).IsFalse()

	// canTestForEqualityWith
	s.AssertCanTestForEqualityWith(s.BOOLEAN_TYPE, s.ALL_TYPE)
	s.AssertCanTestForEqualityWith(s.BOOLEAN_TYPE, s.STRING_OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.BOOLEAN_TYPE, s.NUMBER_TYPE)
	s.AssertCannotTestForEqualityWith(s.BOOLEAN_TYPE, s.functionType)
	s.AssertCannotTestForEqualityWith(s.BOOLEAN_TYPE, s.VOID_TYPE)
	s.AssertCanTestForEqualityWith(s.BOOLEAN_TYPE, s.OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.BOOLEAN_TYPE, s.DATE_TYPE)
	s.AssertCanTestForEqualityWith(s.BOOLEAN_TYPE, s.REGEXP_TYPE)
	s.AssertCanTestForEqualityWith(s.BOOLEAN_TYPE, s.UNKNOWN_TYPE)

	// canTestForShallowEqualityWith
	s.AssertThat(s.BOOLEAN_TYPE.CanTestForShallowEqualityWith(s.NO_TYPE)).IsTrue()
	s.AssertThat(s.BOOLEAN_TYPE.CanTestForShallowEqualityWith(s.NO_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.CanTestForShallowEqualityWith(s.ARRAY_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.CanTestForShallowEqualityWith(s.BOOLEAN_TYPE)).IsTrue()
	s.AssertThat(s.BOOLEAN_TYPE.CanTestForShallowEqualityWith(s.BOOLEAN_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.CanTestForShallowEqualityWith(s.DATE_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.CanTestForShallowEqualityWith(s.functionType)).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.CanTestForShallowEqualityWith(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.CanTestForShallowEqualityWith(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.CanTestForShallowEqualityWith(s.NUMBER_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.CanTestForShallowEqualityWith(s.OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.CanTestForShallowEqualityWith(s.REGEXP_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.CanTestForShallowEqualityWith(s.STRING_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.CanTestForShallowEqualityWith(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.CanTestForShallowEqualityWith(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.BOOLEAN_TYPE.CanTestForShallowEqualityWith(s.VOID_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.CanTestForShallowEqualityWith(s.UNKNOWN_TYPE)).IsTrue()

	// isNullable
	s.AssertThat(s.BOOLEAN_TYPE.IsNullable()).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.IsVoidable()).IsFalse()

	// matchesXxx
	s.AssertThat(s.BOOLEAN_TYPE.MatchesNumberContext()).IsTrue()
	s.AssertThat(s.BOOLEAN_TYPE.MatchesNumberContext()).IsTrue()
	s.AssertThat(s.BOOLEAN_TYPE.MatchesObjectContext()).IsTrue()
	s.AssertThat(s.BOOLEAN_TYPE.MatchesStringContext()).IsTrue()
	s.AssertThat(s.BOOLEAN_TYPE.MatchesNumberContext()).IsTrue()

	// toString
	s.AssertThat(s.BOOLEAN_TYPE.String()).IsEqualTo("boolean")
	s.AssertThat(s.BOOLEAN_TYPE.HasDisplayName()).IsTrue()
	s.AssertThat(s.BOOLEAN_TYPE.GetDisplayName()).IsEqualTo("boolean")

	rhinotesting.Asserts_AssertResolvesToSame(s.T, s.BOOLEAN_TYPE)
}

/** Tests the behavior of the Boolean type. */
func TestJSType_BooleanObjectType(t *testing.T) {
	s := newJsTypeTest(t)
	// isXxx
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsArrayType()).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsBooleanObjectType()).IsTrue()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsBooleanValueType()).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsDateType()).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsEnumElementType()).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsNamedType()).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsNullType()).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsNumberObjectType()).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsNumberValueType()).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsFunctionPrototypeType()).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.GetImplicitPrototype().IsFunctionPrototypeType()).IsTrue()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsRegexpType()).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsStringObjectType()).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsStringValueType()).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsEnumType()).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsUnionType()).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsStruct()).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsDict()).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsAllType()).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsVoidType()).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsConstructor()).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsInstanceType()).IsTrue()

	// isSubtype
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsSubtypeOf(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsSubtypeOf(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsSubtypeOf(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsSubtypeOf(s.functionType)).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsSubtypeOf(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsSubtypeOf(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsSubtypeOf(s.DATE_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsSubtypeOf(s.unresolvedNamedType)).IsTrue()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsSubtypeOf(s.namedGoogBar)).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsSubtypeOf(s.REGEXP_TYPE)).IsFalse()
	// canBeCalled
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.CanBeCalled()).IsFalse()

	// canTestForEqualityWith
	s.AssertCanTestForEqualityWith(s.BOOLEAN_OBJECT_TYPE, s.ALL_TYPE)
	s.AssertCanTestForEqualityWith(s.BOOLEAN_OBJECT_TYPE, s.STRING_OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.BOOLEAN_OBJECT_TYPE, s.NUMBER_TYPE)
	s.AssertCanTestForEqualityWith(s.BOOLEAN_OBJECT_TYPE, s.functionType)
	s.AssertCannotTestForEqualityWith(s.BOOLEAN_OBJECT_TYPE, s.VOID_TYPE)
	s.AssertCanTestForEqualityWith(s.BOOLEAN_OBJECT_TYPE, s.OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.BOOLEAN_OBJECT_TYPE, s.DATE_TYPE)
	s.AssertCanTestForEqualityWith(s.BOOLEAN_OBJECT_TYPE, s.REGEXP_TYPE)

	// canTestForShallowEqualityWith
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.CanTestForShallowEqualityWith(s.NO_TYPE)).IsTrue()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.CanTestForShallowEqualityWith(s.NO_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.CanTestForShallowEqualityWith(s.ARRAY_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.CanTestForShallowEqualityWith(s.BOOLEAN_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.CanTestForShallowEqualityWith(s.BOOLEAN_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.CanTestForShallowEqualityWith(s.DATE_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.CanTestForShallowEqualityWith(s.functionType)).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.CanTestForShallowEqualityWith(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.CanTestForShallowEqualityWith(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.CanTestForShallowEqualityWith(s.NUMBER_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.CanTestForShallowEqualityWith(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.CanTestForShallowEqualityWith(s.REGEXP_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.CanTestForShallowEqualityWith(s.STRING_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.CanTestForShallowEqualityWith(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.CanTestForShallowEqualityWith(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.CanTestForShallowEqualityWith(s.VOID_TYPE)).IsFalse()

	// isNullable
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsNullable()).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsVoidable()).IsFalse()

	// matchesXxx
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.MatchesNumberContext()).IsTrue()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.MatchesNumberContext()).IsTrue()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.MatchesObjectContext()).IsTrue()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.MatchesStringContext()).IsTrue()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.MatchesNumberContext()).IsTrue()

	// toString
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.String()).IsEqualTo("Boolean")
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.HasDisplayName()).IsTrue()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.GetDisplayName()).IsEqualTo("Boolean")

	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsNativeObjectType()).IsTrue()

	rhinotesting.Asserts_AssertResolvesToSame(s.T, s.BOOLEAN_OBJECT_TYPE)
}

/** Tests the behavior of the enum type. */
func TestJSType_EnumType(t *testing.T) {
	s := newJsTypeTest(t)
	enumType :=
		EnumType_Builder_(s.Registry).SetName(ptr("Enum")).SetElementType(s.NUMBER_TYPE).Build()

	// isXxx
	s.AssertThat(enumType.IsArrayType()).IsFalse()
	s.AssertThat(enumType.IsBooleanObjectType()).IsFalse()
	s.AssertThat(enumType.IsBooleanValueType()).IsFalse()
	s.AssertThat(enumType.IsDateType()).IsFalse()
	s.AssertThat(enumType.IsEnumElementType()).IsFalse()
	s.AssertThat(enumType.IsNamedType()).IsFalse()
	s.AssertThat(enumType.IsNullType()).IsFalse()
	s.AssertThat(enumType.IsNumberObjectType()).IsFalse()
	s.AssertThat(enumType.IsNumberValueType()).IsFalse()
	s.AssertThat(enumType.IsFunctionPrototypeType()).IsFalse()
	s.AssertThat(enumType.IsRegexpType()).IsFalse()
	s.AssertThat(enumType.IsStringObjectType()).IsFalse()
	s.AssertThat(enumType.IsStringValueType()).IsFalse()
	s.AssertThat(enumType.IsEnumType()).IsTrue()
	s.AssertThat(enumType.IsUnionType()).IsFalse()
	s.AssertThat(enumType.IsStruct()).IsFalse()
	s.AssertThat(enumType.IsDict()).IsFalse()
	s.AssertThat(enumType.IsAllType()).IsFalse()
	s.AssertThat(enumType.IsVoidType()).IsFalse()
	s.AssertThat(enumType.IsConstructor()).IsFalse()
	s.AssertThat(enumType.IsInstanceType()).IsFalse()

	// isSubtype
	s.AssertThat(enumType.IsSubtype(s.ALL_TYPE)).IsTrue()
	s.AssertThat(enumType.IsSubtype(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(enumType.IsSubtype(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(enumType.IsSubtype(s.functionType)).IsFalse()
	s.AssertThat(enumType.IsSubtype(s.NULL_TYPE)).IsFalse()
	s.AssertThat(enumType.IsSubtype(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(enumType.IsSubtype(s.DATE_TYPE)).IsFalse()
	s.AssertThat(enumType.IsSubtype(s.unresolvedNamedType)).IsTrue()
	s.AssertThat(enumType.IsSubtype(s.namedGoogBar)).IsFalse()
	s.AssertThat(enumType.IsSubtype(s.REGEXP_TYPE)).IsFalse()

	// canBeCalled
	s.AssertThat(enumType.CanBeCalled()).IsFalse()

	// canTestForEqualityWith
	s.AssertCanTestForEqualityWith(enumType, s.ALL_TYPE)
	s.AssertCanTestForEqualityWith(enumType, s.STRING_OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(enumType, s.NUMBER_TYPE)
	s.AssertCanTestForEqualityWith(enumType, s.functionType)
	s.AssertCannotTestForEqualityWith(enumType, s.VOID_TYPE)
	s.AssertCanTestForEqualityWith(enumType, s.OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(enumType, s.DATE_TYPE)
	s.AssertCanTestForEqualityWith(enumType, s.REGEXP_TYPE)

	// canTestForShallowEqualityWith
	s.AssertThat(enumType.CanTestForShallowEqualityWith(s.NO_TYPE)).IsTrue()
	s.AssertThat(enumType.CanTestForShallowEqualityWith(s.NO_OBJECT_TYPE)).IsTrue()
	s.AssertThat(enumType.CanTestForShallowEqualityWith(s.ARRAY_TYPE)).IsFalse()
	s.AssertThat(enumType.CanTestForShallowEqualityWith(s.BOOLEAN_TYPE)).IsFalse()
	s.AssertThat(enumType.CanTestForShallowEqualityWith(enumType)).IsTrue()
	s.AssertThat(enumType.CanTestForShallowEqualityWith(s.DATE_TYPE)).IsFalse()
	s.AssertThat(enumType.CanTestForShallowEqualityWith(s.functionType)).IsFalse()
	s.AssertThat(enumType.CanTestForShallowEqualityWith(s.NULL_TYPE)).IsFalse()
	s.AssertThat(enumType.CanTestForShallowEqualityWith(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(enumType.CanTestForShallowEqualityWith(s.NUMBER_OBJECT_TYPE)).IsFalse()
	s.AssertThat(enumType.CanTestForShallowEqualityWith(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(enumType.CanTestForShallowEqualityWith(s.REGEXP_TYPE)).IsFalse()
	s.AssertThat(enumType.CanTestForShallowEqualityWith(s.STRING_TYPE)).IsFalse()
	s.AssertThat(enumType.CanTestForShallowEqualityWith(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(enumType.CanTestForShallowEqualityWith(s.ALL_TYPE)).IsTrue()
	s.AssertThat(enumType.CanTestForShallowEqualityWith(s.VOID_TYPE)).IsFalse()

	// isNullable
	s.AssertThat(enumType.IsNullable()).IsFalse()
	s.AssertThat(enumType.IsVoidable()).IsFalse()

	// matchesXxx
	s.AssertThat(enumType.MatchesNumberContext()).IsFalse()
	s.AssertThat(enumType.MatchesNumberContext()).IsFalse()
	s.AssertThat(enumType.MatchesObjectContext()).IsTrue()
	s.AssertThat(enumType.MatchesStringContext()).IsTrue()
	s.AssertThat(enumType.MatchesNumberContext()).IsFalse()

	// toString
	s.AssertThat(enumType.String()).IsEqualTo("enum{Enum}")
	s.AssertThat(enumType.HasDisplayName()).IsTrue()
	s.AssertThat(enumType.GetDisplayName()).IsEqualTo("Enum")

	s.AssertThat(
		EnumType_Builder_(s.Registry).
			SetName(ptr("AnotherEnum")).
			SetElementType(s.NUMBER_TYPE).
			Build().
			GetDisplayName()).
		IsEqualTo("AnotherEnum")
	s.AssertThat(EnumType_Builder_(s.Registry).SetElementType(s.NUMBER_TYPE).Build().HasDisplayName()).
		IsFalse()

	rhinotesting.Asserts_AssertResolvesToSame(s.T, enumType)
}

/** Tests the behavior of the enum element type. */
func TestJSType_EnumElementType(t *testing.T) {
	s := newJsTypeTest(t)
	// isXxx
	s.AssertThat(s.elementsType.IsArrayType()).IsFalse()
	s.AssertThat(s.elementsType.IsBooleanObjectType()).IsFalse()
	s.AssertThat(s.elementsType.IsBooleanValueType()).IsFalse()
	s.AssertThat(s.elementsType.IsDateType()).IsFalse()
	s.AssertThat(s.elementsType.IsEnumElementType()).IsTrue()
	s.AssertThat(s.elementsType.IsNamedType()).IsFalse()
	s.AssertThat(s.elementsType.IsNullType()).IsFalse()
	s.AssertThat(s.elementsType.IsNumberObjectType()).IsFalse()
	s.AssertThat(s.elementsType.IsNumberValueType()).IsFalse()
	s.AssertThat(s.elementsType.IsFunctionPrototypeType()).IsFalse()
	s.AssertThat(s.elementsType.IsRegexpType()).IsFalse()
	s.AssertThat(s.elementsType.IsStringObjectType()).IsFalse()
	s.AssertThat(s.elementsType.IsStringValueType()).IsFalse()
	s.AssertThat(s.elementsType.IsEnumType()).IsFalse()
	s.AssertThat(s.elementsType.IsUnionType()).IsFalse()
	s.AssertThat(s.elementsType.IsStruct()).IsFalse()
	s.AssertThat(s.elementsType.IsDict()).IsFalse()
	s.AssertThat(s.elementsType.IsAllType()).IsFalse()
	s.AssertThat(s.elementsType.IsVoidType()).IsFalse()
	s.AssertThat(s.elementsType.IsConstructor()).IsFalse()
	s.AssertThat(s.elementsType.IsInstanceType()).IsFalse()

	// isSubtype
	s.AssertThat(s.elementsType.IsSubtype(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.elementsType.IsSubtype(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.elementsType.IsSubtype(s.NUMBER_TYPE)).IsTrue()
	s.AssertThat(s.elementsType.IsSubtype(s.functionType)).IsFalse()
	s.AssertThat(s.elementsType.IsSubtype(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.elementsType.IsSubtype(s.OBJECT_TYPE)).IsFalse() // no more autoboxing
	s.AssertThat(s.elementsType.IsSubtype(s.DATE_TYPE)).IsFalse()
	s.AssertThat(s.elementsType.IsSubtype(s.unresolvedNamedType)).IsTrue()
	s.AssertThat(s.elementsType.IsSubtype(s.namedGoogBar)).IsFalse()
	s.AssertThat(s.elementsType.IsSubtype(s.REGEXP_TYPE)).IsFalse()

	// canBeCalled
	s.AssertThat(s.elementsType.CanBeCalled()).IsFalse()

	// canTestForEqualityWith
	s.AssertCanTestForEqualityWith(s.elementsType, s.ALL_TYPE)
	s.AssertCanTestForEqualityWith(s.elementsType, s.STRING_OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.elementsType, s.NUMBER_TYPE)
	s.AssertCanTestForEqualityWith(s.elementsType, s.NUMBER_OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.elementsType, s.elementsType)
	s.AssertCannotTestForEqualityWith(s.elementsType, s.functionType)
	s.AssertCannotTestForEqualityWith(s.elementsType, s.VOID_TYPE)
	s.AssertCanTestForEqualityWith(s.elementsType, s.OBJECT_TYPE)
	s.AssertCanTestForEqualityWith(s.elementsType, s.DATE_TYPE)
	s.AssertCanTestForEqualityWith(s.elementsType, s.REGEXP_TYPE)

	// canTestForShallowEqualityWith
	s.AssertThat(s.elementsType.CanTestForShallowEqualityWith(s.NO_TYPE)).IsTrue()
	s.AssertThat(s.elementsType.CanTestForShallowEqualityWith(s.NO_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.elementsType.CanTestForShallowEqualityWith(s.ARRAY_TYPE)).IsFalse()
	s.AssertThat(s.elementsType.CanTestForShallowEqualityWith(s.BOOLEAN_TYPE)).IsFalse()
	s.AssertThat(s.elementsType.CanTestForShallowEqualityWith(s.elementsType)).IsTrue()
	s.AssertThat(s.elementsType.CanTestForShallowEqualityWith(s.DATE_TYPE)).IsFalse()
	s.AssertThat(s.elementsType.CanTestForShallowEqualityWith(s.functionType)).IsFalse()
	s.AssertThat(s.elementsType.CanTestForShallowEqualityWith(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.elementsType.CanTestForShallowEqualityWith(s.NUMBER_TYPE)).IsTrue()
	s.AssertThat(s.elementsType.CanTestForShallowEqualityWith(s.NUMBER_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.elementsType.CanTestForShallowEqualityWith(s.OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.elementsType.CanTestForShallowEqualityWith(s.REGEXP_TYPE)).IsFalse()
	s.AssertThat(s.elementsType.CanTestForShallowEqualityWith(s.STRING_TYPE)).IsFalse()
	s.AssertThat(s.elementsType.CanTestForShallowEqualityWith(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.elementsType.CanTestForShallowEqualityWith(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.elementsType.CanTestForShallowEqualityWith(s.VOID_TYPE)).IsFalse()

	// isNullable
	s.AssertThat(s.elementsType.IsNullable()).IsFalse()
	s.AssertThat(s.elementsType.IsVoidable()).IsFalse()

	// matchesXxx
	s.AssertThat(s.elementsType.MatchesNumberContext()).IsTrue()
	s.AssertThat(s.elementsType.MatchesNumberContext()).IsTrue()
	s.AssertThat(s.elementsType.MatchesObjectContext()).IsTrue()
	s.AssertThat(s.elementsType.MatchesStringContext()).IsTrue()
	s.AssertThat(s.elementsType.MatchesNumberContext()).IsTrue()

	// toString
	s.AssertThat(s.elementsType.String()).IsEqualTo("Enum<number>")
	s.AssertThat(s.elementsType.HasDisplayName()).IsTrue()
	s.AssertThat(s.elementsType.GetDisplayName()).IsEqualTo("Enum")

	rhinotesting.Asserts_AssertResolvesToSame(s.T, s.elementsType)
}

func TestJSType_StringEnumType(t *testing.T) {
	s := newJsTypeTest(t)
	stringEnum :=
		EnumType_Builder_(s.Registry).
			SetName(ptr("Enum")).
			SetElementType(s.STRING_TYPE).
			Build().
			GetElementsType()

	s.AssertTypeEquals(s.UNKNOWN_TYPE, stringEnum.GetPropertyType("length"))
	s.AssertTypeEquals(s.NUMBER_TYPE, stringEnum.FindPropertyType("length"))
	s.AssertThat(stringEnum.HasProperty("length")).IsFalse()
	s.AssertTypeEquals(s.STRING_OBJECT_TYPE, stringEnum.AutoboxesTo())
	s.AssertThat(stringEnum.GetConstructor()).IsNull()

	rhinotesting.Asserts_AssertResolvesToSame(s.T, stringEnum)
}

func TestJSType_StringObjectEnumType(t *testing.T) {
	s := newJsTypeTest(t)
	stringEnum :=
		EnumType_Builder_(s.Registry).
			SetName(ptr("Enum")).
			SetElementType(s.STRING_OBJECT_TYPE).
			Build().
			GetElementsType()

	s.AssertTypeEquals(s.NUMBER_TYPE, stringEnum.GetPropertyType("length"))
	s.AssertTypeEquals(s.NUMBER_TYPE, stringEnum.FindPropertyType("length"))
	s.AssertThat(stringEnum.HasProperty("length")).IsTrue()
	s.AssertTypeEquals(s.STRING_OBJECT_FUNCTION_TYPE, stringEnum.GetConstructor())
}

/** Tests object types. */
func TestJSType_ObjectType(t *testing.T) {
	s := newJsTypeTest(t)
	objectType := PrototypeObjectType_Builder_(s.Registry).Build()

	// isXxx
	s.AssertThat(objectType.IsAllType()).IsFalse()
	s.AssertThat(objectType.IsArrayType()).IsFalse()
	s.AssertThat(objectType.IsDateType()).IsFalse()
	s.AssertThat(objectType.IsFunctionPrototypeType()).IsFalse()
	s.AssertThat(s.OBJECT_TYPE).IsSameInstanceAs(objectType.GetImplicitPrototype())

	// isSubtype
	s.AssertThat(objectType.IsSubtypeOf(s.ALL_TYPE)).IsTrue()
	s.AssertThat(objectType.IsSubtypeOf(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(objectType.IsSubtypeOf(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(objectType.IsSubtypeOf(s.functionType)).IsFalse()
	s.AssertThat(objectType.IsSubtypeOf(s.NULL_TYPE)).IsFalse()
	s.AssertThat(objectType.IsSubtypeOf(s.DATE_TYPE)).IsFalse()
	s.AssertThat(objectType.IsSubtypeOf(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(objectType.IsSubtypeOf(s.unresolvedNamedType)).IsTrue()
	s.AssertThat(objectType.IsSubtypeOf(s.namedGoogBar)).IsFalse()
	s.AssertThat(objectType.IsSubtypeOf(s.REGEXP_TYPE)).IsFalse()

	// autoboxesTo
	s.AssertThat(objectType.AutoboxesTo()).IsNull()

	// canTestForEqualityWith
	s.AssertCanTestForEqualityWith(objectType, s.NUMBER_TYPE)
	s.AssertCanTestForEqualityWith(objectType, s.BIGINT_TYPE)

	// matchesXxxContext
	s.AssertThat(objectType.MatchesNumberContext()).IsFalse()
	s.AssertThat(objectType.MatchesNumberContext()).IsFalse()
	s.AssertThat(objectType.MatchesObjectContext()).IsTrue()
	s.AssertThat(objectType.MatchesStringContext()).IsFalse()
	s.AssertThat(objectType.MatchesNumberContext()).IsFalse()

	// isNullable
	s.AssertThat(objectType.IsNullable()).IsFalse()
	s.AssertThat(objectType.IsVoidable()).IsFalse()
	s.AssertThat(s.CreateNullableType(objectType).IsNullable()).IsTrue()
	s.AssertThat(s.CreateUnionType(objectType, s.VOID_TYPE).IsVoidable()).IsTrue()

	// toString
	s.AssertThat(objectType.String()).IsEqualTo("{...}")
	s.AssertThat(objectType.GetDisplayName()).IsNull()
	s.AssertType(objectType).GetReferenceNameIsNull()
	s.AssertThat(PrototypeObjectType_Builder_(s.Registry).SetName(ptr("anObject")).Build().GetDisplayName()).
		IsEqualTo("anObject")

	rhinotesting.Asserts_AssertResolvesToSame(s.T, objectType)
}

/** Tests the goog.Bar type. */
func TestJSType_GoogBar(t *testing.T) {
	s := newJsTypeTest(t)
	s.AssertThat(s.namedGoogBar.IsInstanceType()).IsTrue()
	s.AssertThat(s.googBar.IsInstanceType()).IsFalse()
	s.AssertThat(s.namedGoogBar.IsConstructor()).IsFalse()
	s.AssertThat(s.googBar.IsConstructor()).IsTrue()
	s.AssertThat(s.googBar.GetInstanceType().IsInstanceType()).IsTrue()
	s.AssertThat(s.namedGoogBar.GetConstructor().IsConstructor()).IsTrue()
	s.AssertThat(s.namedGoogBar.GetImplicitPrototype().IsFunctionPrototypeType()).IsTrue()

	// isSubtype
	s.AssertTypeCanAssignToItself(s.googBar)
	s.AssertTypeCanAssignToItself(s.namedGoogBar)
	s.googBar.IsSubtype(s.namedGoogBar)
	s.namedGoogBar.IsSubtype(s.googBar)
	s.AssertTypeEquals(s.googBar, s.googBar)
	s.AssertTypeNotEquals(s.googBar, s.googSubBar)

	rhinotesting.Asserts_AssertResolvesToSame(s.T, s.googBar)
	rhinotesting.Asserts_AssertResolvesToSame(s.T, s.googSubBar)
}

/** Tests how properties are counted for object types. */
func TestJSType_ObjectTypePropertiesCount(t *testing.T) {
	s := newJsTypeTest(t)
	sup := s.Registry.CreateAnonymousObjectType(nil)
	nativeProperties := sup.GetPropertiesCount()

	sup.DefineDeclaredProperty("a", s.DATE_TYPE, nil)
	s.AssertThat(sup.GetPropertiesCount()).IsEqualTo(nativeProperties + 1)

	sup.DefineDeclaredProperty("b", s.DATE_TYPE, nil)
	s.AssertThat(sup.GetPropertiesCount()).IsEqualTo(nativeProperties + 2)

	sub := s.Registry.CreateObjectType(nil, sup)
	s.AssertThat(sub.GetPropertiesCount()).IsEqualTo(nativeProperties + 2)
}

/** Tests how properties are defined. */
func TestJSType_DefineProperties(t *testing.T) {
	s := newJsTypeTest(t)
	prototype := s.googBar.GetPrototype()
	instance := s.googBar.GetInstanceType()

	s.AssertTypeEquals(instance.GetImplicitPrototype(), prototype)

	// Test declarations.
	s.AssertThat(prototype.DefineDeclaredProperty("declared", s.NUMBER_TYPE, nil)).IsTrue()
	s.AssertThat(prototype.DefineDeclaredProperty("declared", s.NUMBER_TYPE, nil)).IsFalse()
	s.AssertThat(instance.DefineDeclaredProperty("declared", s.NUMBER_TYPE, nil)).IsFalse()
	s.AssertTypeEquals(s.NUMBER_TYPE, instance.GetPropertyType("declared"))

	// Test inferring different types.
	s.AssertThat(prototype.DefineInferredProperty("inferred1", s.STRING_TYPE, nil)).IsTrue()
	s.AssertThat(prototype.DefineInferredProperty("inferred1", s.NUMBER_TYPE, nil)).IsTrue()
	s.AssertTypeEquals(
		s.CreateUnionType(s.NUMBER_TYPE, s.STRING_TYPE),
		instance.GetPropertyType("inferred1"))

	// Test inferring different types on different objects.
	s.AssertThat(prototype.DefineInferredProperty("inferred2", s.STRING_TYPE, nil)).IsTrue()
	s.AssertThat(instance.DefineInferredProperty("inferred2", s.NUMBER_TYPE, nil)).IsTrue()
	s.AssertTypeEquals(
		s.CreateUnionType(s.NUMBER_TYPE, s.STRING_TYPE),
		instance.GetPropertyType("inferred2"))

	// Test inferring on the supertype and declaring on the subtype.
	s.AssertThat(prototype.DefineInferredProperty("prop", s.STRING_TYPE, nil)).IsTrue()
	s.AssertThat(instance.DefineDeclaredProperty("prop", s.NUMBER_TYPE, nil)).IsTrue()
	s.AssertTypeEquals(s.NUMBER_TYPE, instance.GetPropertyType("prop"))
	s.AssertTypeEquals(s.STRING_TYPE, prototype.GetPropertyType("prop"))
}

/** Tests that properties are correctly counted even when shadowing occurs. */
func TestJSType_ObjectTypePropertiesCountWithShadowing(t *testing.T) {
	s := newJsTypeTest(t)
	sup := s.Registry.CreateAnonymousObjectType(nil)
	nativeProperties := sup.GetPropertiesCount()

	sup.DefineDeclaredProperty("a", s.OBJECT_TYPE, nil)
	s.AssertThat(sup.GetPropertiesCount()).IsEqualTo(nativeProperties + 1)

	sub := s.Registry.CreateObjectType(nil, sup)
	sub.DefineDeclaredProperty("a", s.OBJECT_TYPE, nil)
	s.AssertThat(sub.GetPropertiesCount()).IsEqualTo(nativeProperties + 1)
}

/** Tests the named type goog.Bar. */
func TestJSType_NamedGoogBar(t *testing.T) {
	s := newJsTypeTest(t)
	// isXxx
	s.AssertThat(s.namedGoogBar.IsFunctionPrototypeType()).IsFalse()
	s.AssertThat(s.namedGoogBar.GetImplicitPrototype().IsFunctionPrototypeType()).IsTrue()

	// isSubtype
	s.AssertThat(s.namedGoogBar.IsSubtype(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.namedGoogBar.IsSubtype(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.namedGoogBar.IsSubtype(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.namedGoogBar.IsSubtype(s.functionType)).IsFalse()
	s.AssertThat(s.namedGoogBar.IsSubtype(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.namedGoogBar.IsSubtype(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.namedGoogBar.IsSubtype(s.DATE_TYPE)).IsFalse()
	s.AssertThat(s.namedGoogBar.IsSubtype(s.namedGoogBar)).IsTrue()
	s.AssertThat(s.namedGoogBar.IsSubtype(s.unresolvedNamedType)).IsTrue()
	s.AssertThat(s.namedGoogBar.IsSubtype(s.REGEXP_TYPE)).IsFalse()
	s.AssertThat(s.namedGoogBar.IsSubtype(s.ARRAY_TYPE)).IsFalse()

	// autoboxesTo
	s.AssertThat(s.namedGoogBar.AutoboxesTo()).IsNull()

	// properties
	s.AssertTypeEquals(s.DATE_TYPE, s.namedGoogBar.GetPropertyType("date"))

	s.AssertThat(s.namedGoogBar.IsNativeObjectType()).IsFalse()
	s.AssertThat(s.namedGoogBar.GetImplicitPrototype().IsNativeObjectType()).IsFalse()

	resolvedNamedGoogBar := rhinotesting.Asserts_AssertValidResolve(s.T, s.namedGoogBar)
	s.AssertThat(s.namedGoogBar).IsNotSameInstanceAs(resolvedNamedGoogBar)
	s.AssertThat(s.googBar.GetInstanceType()).IsSameInstanceAs(resolvedNamedGoogBar)
}

/** Tests the prototype chaining of native objects. */
func TestJSType_PrototypeChaining(t *testing.T) {
	s := newJsTypeTest(t)
	// equals
	s.AssertTypeEquals(
		s.ARRAY_TYPE.GetImplicitPrototype().GetImplicitPrototype(),
		s.OBJECT_TYPE)
	s.AssertTypeEquals(
		s.BOOLEAN_OBJECT_TYPE.GetImplicitPrototype().
			GetImplicitPrototype(), s.OBJECT_TYPE)
	s.AssertTypeEquals(
		s.DATE_TYPE.GetImplicitPrototype().GetImplicitPrototype(),
		s.OBJECT_TYPE)
	s.AssertTypeEquals(
		s.NUMBER_OBJECT_TYPE.GetImplicitPrototype().
			GetImplicitPrototype(), s.OBJECT_TYPE)
	s.AssertTypeEquals(
		s.STRING_OBJECT_TYPE.GetImplicitPrototype().
			GetImplicitPrototype(), s.OBJECT_TYPE)
	s.AssertTypeEquals(
		s.REGEXP_TYPE.GetImplicitPrototype().GetImplicitPrototype(),
		s.OBJECT_TYPE)
}

/**
* Tests that function instances have their constructor pointer back at the function that created
* them.
 */
func TestJSType_InstanceFunctionChaining(t *testing.T) {
	s := newJsTypeTest(t)
	// Array
	s.AssertTypeEquals(
		s.ARRAY_FUNCTION_TYPE, s.ARRAY_TYPE.GetConstructor())

	// Boolean
	s.AssertTypeEquals(
		s.BOOLEAN_OBJECT_FUNCTION_TYPE,
		s.BOOLEAN_OBJECT_TYPE.GetConstructor())

	// Date
	s.AssertTypeEquals(
		s.DATE_FUNCTION_TYPE, s.DATE_TYPE.GetConstructor())

	// Number
	s.AssertTypeEquals(
		s.NUMBER_OBJECT_FUNCTION_TYPE,
		s.NUMBER_OBJECT_TYPE.GetConstructor())

	// Object
	s.AssertTypeEquals(
		s.OBJECT_FUNCTION_TYPE, s.OBJECT_TYPE.GetConstructor())

	// RegExp
	s.AssertTypeEquals(s.REGEXP_FUNCTION_TYPE, s.REGEXP_TYPE.GetConstructor())

	// String
	s.AssertTypeEquals(
		s.STRING_OBJECT_FUNCTION_TYPE,
		s.STRING_OBJECT_TYPE.GetConstructor())
}

/**
* Tests that the method {@link JSType#canTestForEqualityWith(JSType)} handles special corner
* cases.
 */
func TestJSType_CanTestForEqualityWithCornerCases(t *testing.T) {
	s := newJsTypeTest(t)
	// null == undefined is always true
	s.AssertCannotTestForEqualityWith(s.NULL_TYPE, s.VOID_TYPE)

	// (Object,null) == undefined could be true or false
	nullableObject :=
		s.CreateUnionType(s.OBJECT_TYPE, s.NULL_TYPE).(*UnionType)
	s.AssertCanTestForEqualityWith(nullableObject, s.VOID_TYPE)
	s.AssertCanTestForEqualityWith(s.VOID_TYPE, nullableObject)
}

/** Tests the {@link JSType#testForEquality(JSType)} method. */
func TestJSType_TestForEquality(t *testing.T) {
	s := newJsTypeTest(t)
	s.Compare(base.Tri_TRUE, s.NO_OBJECT_TYPE, s.NO_OBJECT_TYPE)
	s.Compare(base.Tri_UNKNOWN, s.ALL_TYPE, s.ALL_TYPE)
	s.Compare(base.Tri_TRUE, s.NO_TYPE, s.NO_TYPE)
	fooNoResolvedType := NewNoResolvedType(s.Registry, "Foo", nil)
	s.Compare(base.Tri_UNKNOWN, fooNoResolvedType, fooNoResolvedType)
	s.Compare(base.Tri_UNKNOWN, s.NO_OBJECT_TYPE, s.NUMBER_TYPE)
	s.Compare(base.Tri_UNKNOWN, s.ALL_TYPE, s.NUMBER_TYPE)
	s.Compare(base.Tri_UNKNOWN, s.NO_TYPE, s.NUMBER_TYPE)

	s.Compare(base.Tri_FALSE, s.NULL_TYPE, s.BOOLEAN_TYPE)
	s.Compare(base.Tri_TRUE, s.NULL_TYPE, s.NULL_TYPE)
	s.Compare(base.Tri_FALSE, s.NULL_TYPE, s.NUMBER_TYPE)
	s.Compare(base.Tri_FALSE, s.NULL_TYPE, s.OBJECT_TYPE)
	s.Compare(base.Tri_FALSE, s.NULL_TYPE, s.STRING_TYPE)
	s.Compare(base.Tri_TRUE, s.NULL_TYPE, s.VOID_TYPE)
	s.Compare(base.Tri_UNKNOWN, s.NULL_TYPE, s.CreateUnionType(s.UNKNOWN_TYPE, s.VOID_TYPE))
	s.Compare(base.Tri_UNKNOWN, s.NULL_TYPE, s.CreateUnionType(s.OBJECT_TYPE, s.VOID_TYPE))
	s.Compare(base.Tri_UNKNOWN, s.NULL_TYPE, s.unresolvedNamedType)
	s.Compare(base.Tri_UNKNOWN,
		s.NULL_TYPE, s.CreateUnionType(s.unresolvedNamedType, s.DATE_TYPE))

	s.Compare(base.Tri_FALSE, s.VOID_TYPE, s.REGEXP_TYPE)
	s.Compare(base.Tri_TRUE, s.VOID_TYPE, s.VOID_TYPE)
	s.Compare(base.Tri_UNKNOWN, s.VOID_TYPE, s.CreateUnionType(s.REGEXP_TYPE, s.VOID_TYPE))

	s.Compare(base.Tri_UNKNOWN, s.NUMBER_TYPE, s.BOOLEAN_TYPE)
	s.Compare(base.Tri_UNKNOWN, s.NUMBER_TYPE, s.NUMBER_TYPE)
	s.Compare(base.Tri_UNKNOWN, s.NUMBER_TYPE, s.OBJECT_TYPE)

	s.Compare(base.Tri_UNKNOWN, s.ARRAY_TYPE, s.BOOLEAN_TYPE)
	s.Compare(base.Tri_UNKNOWN, s.OBJECT_TYPE, s.BOOLEAN_TYPE)
	s.Compare(base.Tri_UNKNOWN, s.OBJECT_TYPE, s.STRING_TYPE)

	s.Compare(base.Tri_UNKNOWN, s.STRING_TYPE, s.STRING_TYPE)

	s.Compare(base.Tri_UNKNOWN, s.STRING_TYPE, s.BOOLEAN_TYPE)
	s.Compare(base.Tri_UNKNOWN, s.STRING_TYPE, s.NUMBER_TYPE)
	s.Compare(base.Tri_FALSE, s.STRING_TYPE, s.VOID_TYPE)
	s.Compare(base.Tri_FALSE, s.STRING_TYPE, s.NULL_TYPE)
	s.Compare(base.Tri_FALSE, s.STRING_TYPE, s.CreateUnionType(s.NULL_TYPE, s.VOID_TYPE))

	s.Compare(base.Tri_UNKNOWN, s.UNKNOWN_TYPE, s.BOOLEAN_TYPE)
	s.Compare(base.Tri_UNKNOWN, s.UNKNOWN_TYPE, s.NULL_TYPE)
	s.Compare(base.Tri_UNKNOWN, s.UNKNOWN_TYPE, s.VOID_TYPE)

	s.Compare(base.Tri_FALSE, s.FUNCTION_TYPE, s.BOOLEAN_TYPE)
	s.Compare(base.Tri_FALSE, s.FUNCTION_TYPE, s.NUMBER_TYPE)
	s.Compare(base.Tri_FALSE, s.FUNCTION_TYPE, s.STRING_TYPE)
	s.Compare(base.Tri_FALSE, s.FUNCTION_TYPE, s.VOID_TYPE)
	s.Compare(base.Tri_FALSE, s.FUNCTION_TYPE, s.NULL_TYPE)
	s.Compare(base.Tri_UNKNOWN, s.FUNCTION_TYPE, s.OBJECT_TYPE)
	s.Compare(base.Tri_UNKNOWN, s.FUNCTION_TYPE, s.ALL_TYPE)

	s.Compare(base.Tri_UNKNOWN, s.NULL_TYPE, s.subclassOfUnresolvedNamedType)

	functionAndNull := s.CreateUnionType(s.NULL_TYPE, s.dateMethod)
	s.Compare(base.Tri_UNKNOWN, functionAndNull, s.dateMethod)

	s.Compare(base.Tri_UNKNOWN, s.NULL_TYPE, s.NO_TYPE)
	s.Compare(base.Tri_UNKNOWN, s.VOID_TYPE, s.NO_TYPE)
	s.Compare(base.Tri_UNKNOWN, s.NULL_TYPE, s.unresolvedNamedType)
	s.Compare(base.Tri_UNKNOWN, s.VOID_TYPE, s.unresolvedNamedType)
	s.Compare(base.Tri_TRUE, s.NO_TYPE, s.NO_TYPE)
}

func (s *jsTypeTest) Compare(r base.Tri, t1 JSType, t2 JSType) {
	s.AssertThat(t1.TestForEquality(t2)).IsEqualTo(r)
	s.AssertThat(t2.TestForEquality(t1)).IsEqualTo(r)
}

/** Tests the subtyping relationships among simple types. */
func TestJSType_SubtypingSimpleTypes(t *testing.T) {
	s := newJsTypeTest(t)
	// Any
	s.AssertThat(s.NO_TYPE.IsSubtypeOf(s.NO_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.IsSubtypeOf(s.NO_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.IsSubtypeOf(s.ARRAY_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.IsSubtypeOf(s.BOOLEAN_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.IsSubtypeOf(s.BOOLEAN_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.IsSubtypeOf(s.DATE_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.IsSubtypeOf(s.functionType)).IsTrue()
	s.AssertThat(s.NO_TYPE.IsSubtypeOf(s.NULL_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.IsSubtypeOf(s.NUMBER_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.IsSubtypeOf(s.NUMBER_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.IsSubtypeOf(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.IsSubtypeOf(s.REGEXP_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.IsSubtypeOf(s.STRING_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.IsSubtypeOf(s.STRING_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.IsSubtypeOf(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.NO_TYPE.IsSubtypeOf(s.VOID_TYPE)).IsTrue()

	// AnyObject
	s.AssertThat(s.NO_OBJECT_TYPE.IsSubtypeOf(s.NO_TYPE)).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.IsSubtypeOf(s.NO_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.IsSubtypeOf(s.ARRAY_TYPE)).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.IsSubtypeOf(s.BOOLEAN_TYPE)).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.IsSubtypeOf(s.BOOLEAN_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.IsSubtypeOf(s.DATE_TYPE)).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.IsSubtypeOf(s.functionType)).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.IsSubtypeOf(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.IsSubtypeOf(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.IsSubtypeOf(s.NUMBER_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.IsSubtypeOf(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.IsSubtypeOf(s.REGEXP_TYPE)).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.IsSubtypeOf(s.STRING_TYPE)).IsFalse()
	s.AssertThat(s.NO_OBJECT_TYPE.IsSubtypeOf(s.STRING_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.IsSubtypeOf(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.NO_OBJECT_TYPE.IsSubtypeOf(s.VOID_TYPE)).IsFalse()

	// Array
	s.AssertThat(s.ARRAY_TYPE.IsSubtypeOf(s.NO_TYPE)).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsSubtypeOf(s.NO_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsSubtypeOf(s.ARRAY_TYPE)).IsTrue()
	s.AssertThat(s.ARRAY_TYPE.IsSubtypeOf(s.BOOLEAN_TYPE)).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsSubtypeOf(s.BOOLEAN_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsSubtypeOf(s.DATE_TYPE)).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsSubtypeOf(s.functionType)).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsSubtypeOf(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsSubtypeOf(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsSubtypeOf(s.NUMBER_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsSubtypeOf(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.ARRAY_TYPE.IsSubtypeOf(s.REGEXP_TYPE)).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsSubtypeOf(s.STRING_TYPE)).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsSubtypeOf(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.ARRAY_TYPE.IsSubtypeOf(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.ARRAY_TYPE.IsSubtypeOf(s.VOID_TYPE)).IsFalse()

	// boolean
	s.AssertThat(s.BOOLEAN_TYPE.IsSubtypeOf(s.NO_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.IsSubtypeOf(s.NO_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.IsSubtypeOf(s.ARRAY_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.IsSubtypeOf(s.BOOLEAN_TYPE)).IsTrue()
	s.AssertThat(s.BOOLEAN_TYPE.IsSubtypeOf(s.BOOLEAN_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.IsSubtypeOf(s.DATE_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.IsSubtypeOf(s.functionType)).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.IsSubtypeOf(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.IsSubtypeOf(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.IsSubtypeOf(s.NUMBER_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.IsSubtypeOf(s.OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.IsSubtypeOf(s.REGEXP_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.IsSubtypeOf(s.STRING_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.IsSubtypeOf(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.IsSubtypeOf(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.BOOLEAN_TYPE.IsSubtypeOf(s.VOID_TYPE)).IsFalse()

	// Boolean
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsSubtypeOf(s.NO_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsSubtypeOf(s.NO_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsSubtypeOf(s.ARRAY_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsSubtypeOf(s.BOOLEAN_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsSubtypeOf(s.BOOLEAN_OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsSubtypeOf(s.DATE_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsSubtypeOf(s.functionType)).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsSubtypeOf(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsSubtypeOf(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsSubtypeOf(s.NUMBER_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsSubtypeOf(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsSubtypeOf(s.REGEXP_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsSubtypeOf(s.STRING_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsSubtypeOf(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsSubtypeOf(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.BOOLEAN_OBJECT_TYPE.IsSubtypeOf(s.VOID_TYPE)).IsFalse()

	// Date
	s.AssertThat(s.DATE_TYPE.IsSubtypeOf(s.NO_TYPE)).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsSubtypeOf(s.NO_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsSubtypeOf(s.ARRAY_TYPE)).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsSubtypeOf(s.BOOLEAN_TYPE)).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsSubtypeOf(s.BOOLEAN_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsSubtypeOf(s.DATE_TYPE)).IsTrue()
	s.AssertThat(s.DATE_TYPE.IsSubtypeOf(s.functionType)).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsSubtypeOf(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsSubtypeOf(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsSubtypeOf(s.NUMBER_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsSubtypeOf(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.DATE_TYPE.IsSubtypeOf(s.REGEXP_TYPE)).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsSubtypeOf(s.STRING_TYPE)).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsSubtypeOf(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.DATE_TYPE.IsSubtypeOf(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.DATE_TYPE.IsSubtypeOf(s.VOID_TYPE)).IsFalse()

	// Unknown
	s.AssertThat(s.ALL_TYPE.IsSubtypeOf(s.NO_TYPE)).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsSubtypeOf(s.NO_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsSubtypeOf(s.ARRAY_TYPE)).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsSubtypeOf(s.BOOLEAN_TYPE)).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsSubtypeOf(s.BOOLEAN_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsSubtypeOf(s.functionType)).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsSubtypeOf(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsSubtypeOf(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsSubtypeOf(s.NUMBER_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsSubtypeOf(s.OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsSubtypeOf(s.REGEXP_TYPE)).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsSubtypeOf(s.STRING_TYPE)).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsSubtypeOf(s.STRING_OBJECT_TYPE)).IsFalse()
	s.AssertThat(s.ALL_TYPE.IsSubtypeOf(s.ALL_TYPE)).IsTrue()
	s.AssertThat(s.ALL_TYPE.IsSubtypeOf(s.VOID_TYPE)).IsFalse()
}

/** Tests that the Object type is the greatest element (top) of the object hierarchy. */
func TestJSType_SubtypingObjectTopOfObjects(t *testing.T) {
	s := newJsTypeTest(t)
	s.AssertThat(s.OBJECT_TYPE.IsSubtypeOf(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.CreateUnionType(s.DATE_TYPE, s.REGEXP_TYPE).IsSubtypeOf(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.CreateUnionType(s.OBJECT_TYPE, s.NO_OBJECT_TYPE).IsSubtypeOf(s.OBJECT_TYPE)).IsTrue()
	s.AssertThat(s.functionType.IsSubtype(s.OBJECT_TYPE)).IsTrue()
}

func TestJSType_SubtypingFunctionPrototypeType(t *testing.T) {
	s := newJsTypeTest(t)

	sub1 :=
		withOpenRegistry(s, func() FunctionType {
			return s.Registry.CreateConstructorType(
				ptr("Foo1"),
				/* source= */ nil,
				s.Registry.CreateParameters(s.ALL_TYPE),
				nil,
				nil,
				/* isAbstract= */ false)
		})
	sub1.SetPrototypeBasedOn(s.googBar)

	sub2 :=
		withOpenRegistry(s, func() FunctionType {
			return s.Registry.CreateConstructorType(
				ptr("Foo2"),
				/* source= */ nil,
				s.Registry.CreateParameters(s.ALL_TYPE),
				nil,
				nil,
				/* isAbstract= */ false)
		})
	sub2.SetPrototypeBasedOn(s.googBar)

	o1 := sub1.GetInstanceType()
	o2 := sub2.GetInstanceType()

	s.AssertThat(o1.IsSubtypeOf(o2)).IsFalse()
	s.AssertThat(o1.GetImplicitPrototype().IsSubtypeOf(o2.GetImplicitPrototype())).IsFalse()
	s.AssertThat(o1.GetImplicitPrototype().IsSubtypeOf(s.googBar)).IsTrue()
	s.AssertThat(o2.GetImplicitPrototype().IsSubtypeOf(s.googBar)).IsTrue()
}

func TestJSType_SubtypingFunctionFixedArgs(t *testing.T) {
	s := newJsTypeTest(t)
	f1 := s.Registry.CreateFunctionType(s.OBJECT_TYPE, s.BOOLEAN_TYPE)
	f2 := s.Registry.CreateFunctionType(s.STRING_OBJECT_TYPE, s.BOOLEAN_TYPE)

	s.AssertThat(f1.IsSubtype(f1)).IsTrue()
	s.AssertThat(f1.IsSubtype(f2)).IsFalse()
	s.AssertThat(f2.IsSubtype(f1)).IsTrue()
	s.AssertThat(f2.IsSubtype(f2)).IsTrue()

	s.AssertThat(f1.IsSubtype(s.FUNCTION_TYPE)).IsTrue()
	s.AssertThat(f2.IsSubtype(s.FUNCTION_TYPE)).IsTrue()
	s.AssertThat(s.FUNCTION_TYPE.IsSubtype(f1)).IsTrue()
	s.AssertThat(s.FUNCTION_TYPE.IsSubtype(f2)).IsTrue()
}

func TestJSType_SubtypingFunctionMultipleFixedArgs(t *testing.T) {
	s := newJsTypeTest(t)
	f1 := s.Registry.CreateFunctionType(s.OBJECT_TYPE, s.NUMBER_TYPE, s.STRING_TYPE)
	f2 := s.Registry.CreateFunctionType(s.STRING_OBJECT_TYPE, s.NUMBER_STRING, s.ALL_TYPE)

	s.AssertThat(f1.IsSubtype(f1)).IsTrue()
	s.AssertThat(f1.IsSubtype(f2)).IsFalse()
	s.AssertThat(f2.IsSubtype(f1)).IsTrue()
	s.AssertThat(f2.IsSubtype(f2)).IsTrue()

	s.AssertThat(f1.IsSubtype(s.FUNCTION_TYPE)).IsTrue()
	s.AssertThat(f2.IsSubtype(s.FUNCTION_TYPE)).IsTrue()
	s.AssertThat(s.FUNCTION_TYPE.IsSubtype(f1)).IsTrue()
	s.AssertThat(s.FUNCTION_TYPE.IsSubtype(f2)).IsTrue()
}

func TestJSType_SubtypingFunctionFixedArgsNotMatching(t *testing.T) {
	s := newJsTypeTest(t)
	f1 := s.Registry.CreateFunctionType(s.OBJECT_TYPE, s.STRING_TYPE, s.UNKNOWN_TYPE)
	f2 := s.Registry.CreateFunctionType(s.STRING_OBJECT_TYPE, s.NUMBER_STRING, s.ALL_TYPE)

	s.AssertThat(f1.IsSubtype(f1)).IsTrue()
	s.AssertThat(f1.IsSubtype(f2)).IsFalse()
	s.AssertThat(f2.IsSubtype(f1)).IsTrue()
	s.AssertThat(f2.IsSubtype(f2)).IsTrue()

	s.AssertThat(f1.IsSubtype(s.FUNCTION_TYPE)).IsTrue()
	s.AssertThat(f2.IsSubtype(s.FUNCTION_TYPE)).IsTrue()
	s.AssertThat(s.FUNCTION_TYPE.IsSubtype(f1)).IsTrue()
	s.AssertThat(s.FUNCTION_TYPE.IsSubtype(f2)).IsTrue()
}

func TestJSType_SubtypingFunctionVariableArgsOneOnly(t *testing.T) {
	s := newJsTypeTest(t)
	// f1 = (string...) -> Object
	f1 := s.Registry.CreateFunctionTypeWithVarArgs(s.OBJECT_TYPE, s.STRING_TYPE)
	// f2 = (string|number, Object) -> String
	f2 := s.Registry.CreateFunctionType(s.STRING_OBJECT_TYPE, s.NUMBER_STRING, s.OBJECT_TYPE)

	s.AssertThat(f1.IsSubtype(f1)).IsTrue()
	s.AssertThat(f1.IsSubtype(f2)).IsFalse()
	s.AssertThat(f2.IsSubtype(f1)).IsFalse()
	s.AssertThat(f2.IsSubtype(f2)).IsTrue()

	s.AssertThat(f1.IsSubtype(s.FUNCTION_TYPE)).IsTrue()
	s.AssertThat(f2.IsSubtype(s.FUNCTION_TYPE)).IsTrue()
	s.AssertThat(s.FUNCTION_TYPE.IsSubtype(f1)).IsTrue()
	s.AssertThat(s.FUNCTION_TYPE.IsSubtype(f2)).IsTrue()
}

func TestJSType_SubtypingFunctionVariableArgsBoth(t *testing.T) {
	s := newJsTypeTest(t)
	// f1 = (number, String, string...) -> Object
	f1 := s.Registry.CreateFunctionTypeWithVarArgs(
		s.OBJECT_TYPE, s.NUMBER_TYPE, s.STRING_OBJECT_TYPE, s.STRING_TYPE)
	// f2 = (string|number, Object, string...) -> String
	f2 := s.Registry.CreateFunctionTypeWithVarArgs(
		s.STRING_OBJECT_TYPE, s.NUMBER_STRING, s.OBJECT_TYPE, s.STRING_TYPE)

	s.AssertThat(f1.IsSubtype(f1)).IsTrue()
	s.AssertThat(f1.IsSubtype(f2)).IsFalse()
	s.AssertThat(f2.IsSubtype(f1)).IsTrue()
	s.AssertThat(f2.IsSubtype(f2)).IsTrue()

	s.AssertThat(f1.IsSubtype(s.FUNCTION_TYPE)).IsTrue()
	s.AssertThat(f2.IsSubtype(s.FUNCTION_TYPE)).IsTrue()
	s.AssertThat(s.FUNCTION_TYPE.IsSubtype(f1)).IsTrue()
	s.AssertThat(s.FUNCTION_TYPE.IsSubtype(f2)).IsTrue()
}

func TestJSType_SubtypingMostGeneralFunction(t *testing.T) {
	s := newJsTypeTest(t)
	// (string, string) -> Object
	f1 := s.Registry.CreateFunctionType(s.OBJECT_TYPE, s.STRING_TYPE, s.STRING_TYPE)
	// (string, void) -> number
	f2 := s.Registry.CreateFunctionType(s.NUMBER_TYPE, s.STRING_TYPE, s.VOID_TYPE)
	// (Date, string, number) -> AnyObject
	f3 := s.Registry.CreateFunctionType(
		s.NO_OBJECT_TYPE, s.DATE_TYPE, s.STRING_TYPE, s.NUMBER_TYPE)
	// (Number) -> Any
	f4 := s.Registry.CreateFunctionType(s.NO_TYPE, s.NUMBER_OBJECT_TYPE)
	// f1 = (string...) -> Object
	f5 := s.Registry.CreateFunctionTypeWithVarArgs(s.OBJECT_TYPE, s.STRING_TYPE)
	// f2 = (string|number, Object) -> String
	f6 := s.Registry.CreateFunctionType(s.STRING_OBJECT_TYPE, s.NUMBER_STRING, s.OBJECT_TYPE)
	// f1 = (number, string...) -> Object
	f7 := s.Registry.CreateFunctionTypeWithVarArgs(
		s.OBJECT_TYPE, s.NUMBER_TYPE, s.STRING_TYPE)
	// f2 = (string|number, Object, string...) -> String
	f8 := s.Registry.CreateFunctionTypeWithVarArgs(
		s.STRING_OBJECT_TYPE, s.NUMBER_STRING, s.OBJECT_TYPE, s.STRING_TYPE)

	s.AssertThat(s.LEAST_FUNCTION_TYPE.IsSubtypeOf(s.GREATEST_FUNCTION_TYPE)).IsTrue()
	s.AssertThat(s.LEAST_FUNCTION_TYPE.IsSubtypeOf(s.FUNCTION_TYPE)).IsTrue()
	s.AssertThat(s.FUNCTION_TYPE.IsSubtype(s.LEAST_FUNCTION_TYPE)).IsTrue()

	s.AssertThat(s.GREATEST_FUNCTION_TYPE.IsSubtypeOf(s.LEAST_FUNCTION_TYPE)).IsFalse()
	s.AssertThat(s.GREATEST_FUNCTION_TYPE.IsSubtypeOf(s.FUNCTION_TYPE)).IsTrue()
	s.AssertThat(s.FUNCTION_TYPE.IsSubtype(s.GREATEST_FUNCTION_TYPE)).IsTrue()

	s.AssertThat(f1.IsSubtype(s.GREATEST_FUNCTION_TYPE)).IsTrue()
	s.AssertThat(f2.IsSubtype(s.GREATEST_FUNCTION_TYPE)).IsTrue()
	s.AssertThat(f3.IsSubtype(s.GREATEST_FUNCTION_TYPE)).IsTrue()
	s.AssertThat(f4.IsSubtype(s.GREATEST_FUNCTION_TYPE)).IsTrue()
	s.AssertThat(f5.IsSubtype(s.GREATEST_FUNCTION_TYPE)).IsTrue()
	s.AssertThat(f6.IsSubtype(s.GREATEST_FUNCTION_TYPE)).IsTrue()
	s.AssertThat(f7.IsSubtype(s.GREATEST_FUNCTION_TYPE)).IsTrue()
	s.AssertThat(f8.IsSubtype(s.GREATEST_FUNCTION_TYPE)).IsTrue()

	s.AssertThat(f1.IsSubtype(s.LEAST_FUNCTION_TYPE)).IsFalse()
	s.AssertThat(f2.IsSubtype(s.LEAST_FUNCTION_TYPE)).IsFalse()
	s.AssertThat(f3.IsSubtype(s.LEAST_FUNCTION_TYPE)).IsFalse()
	s.AssertThat(f4.IsSubtype(s.LEAST_FUNCTION_TYPE)).IsFalse()
	s.AssertThat(f5.IsSubtype(s.LEAST_FUNCTION_TYPE)).IsFalse()
	s.AssertThat(f6.IsSubtype(s.LEAST_FUNCTION_TYPE)).IsFalse()
	s.AssertThat(f7.IsSubtype(s.LEAST_FUNCTION_TYPE)).IsFalse()
	s.AssertThat(f8.IsSubtype(s.LEAST_FUNCTION_TYPE)).IsFalse()

	s.AssertThat(s.LEAST_FUNCTION_TYPE.IsSubtypeOf(f1)).IsTrue()
	s.AssertThat(s.LEAST_FUNCTION_TYPE.IsSubtypeOf(f2)).IsTrue()
	s.AssertThat(s.LEAST_FUNCTION_TYPE.IsSubtypeOf(f3)).IsTrue()
	s.AssertThat(s.LEAST_FUNCTION_TYPE.IsSubtypeOf(f4)).IsTrue()
	s.AssertThat(s.LEAST_FUNCTION_TYPE.IsSubtypeOf(f5)).IsTrue()
	s.AssertThat(s.LEAST_FUNCTION_TYPE.IsSubtypeOf(f6)).IsTrue()
	s.AssertThat(s.LEAST_FUNCTION_TYPE.IsSubtypeOf(f7)).IsTrue()
	s.AssertThat(s.LEAST_FUNCTION_TYPE.IsSubtypeOf(f8)).IsTrue()

	s.AssertThat(s.GREATEST_FUNCTION_TYPE.IsSubtypeOf(f1)).IsFalse()
	s.AssertThat(s.GREATEST_FUNCTION_TYPE.IsSubtypeOf(f2)).IsFalse()
	s.AssertThat(s.GREATEST_FUNCTION_TYPE.IsSubtypeOf(f3)).IsFalse()
	s.AssertThat(s.GREATEST_FUNCTION_TYPE.IsSubtypeOf(f4)).IsFalse()
	s.AssertThat(s.GREATEST_FUNCTION_TYPE.IsSubtypeOf(f5)).IsFalse()
	s.AssertThat(s.GREATEST_FUNCTION_TYPE.IsSubtypeOf(f6)).IsFalse()
	s.AssertThat(s.GREATEST_FUNCTION_TYPE.IsSubtypeOf(f7)).IsFalse()
	s.AssertThat(s.GREATEST_FUNCTION_TYPE.IsSubtypeOf(f8)).IsFalse()
}

/** Types to test for symmetrical relationships. */
func (s *jsTypeTest) GetTypesToTestForSymmetry() []JSType {
	return []JSType{
		s.UNKNOWN_TYPE,
		s.CHECKED_UNKNOWN_TYPE,
		s.NULL_TYPE,
		s.VOID_TYPE,
		s.NUMBER_TYPE,
		s.STRING_TYPE,
		s.BOOLEAN_TYPE,
		s.OBJECT_TYPE,
		s.FUNCTION_TYPE,
		s.LEAST_FUNCTION_TYPE,
		s.GREATEST_FUNCTION_TYPE,
		s.ALL_TYPE,
		s.NO_TYPE,
		s.NO_OBJECT_TYPE,
		NewNoResolvedType(s.Registry, "Foo", nil),
		NewNoResolvedType(s.Registry, "Foo", []JSType{s.STRING_TYPE}),
		NewNoResolvedType(s.Registry, "Bar", nil),
		s.CreateUnionType(NewNoResolvedType(s.Registry, "Foo", nil), s.NULL_TYPE),
		s.CreateUnionType(s.BOOLEAN_TYPE, s.STRING_TYPE),
		s.CreateUnionType(s.NUMBER_TYPE, s.STRING_TYPE),
		s.CreateUnionType(s.NULL_TYPE, s.dateMethod),
		s.CreateUnionType(s.UNKNOWN_TYPE, s.dateMethod),
		s.CreateUnionType(s.namedGoogBar, s.dateMethod),
		s.CreateUnionType(s.NULL_TYPE, s.unresolvedNamedType),
		s.enumType,
		s.elementsType,
		s.dateMethod,
		s.functionType,
		s.unresolvedNamedType,
		s.googBar,
		s.namedGoogBar,
		s.googBar.GetInstanceType(),
		s.namedGoogBar,
		s.subclassOfUnresolvedNamedType,
		s.subclassCtor,
		s.recordType,
		s.forwardDeclaredNamedType,
		s.CreateUnionType(s.forwardDeclaredNamedType, s.NULL_TYPE),
		s.CreateTemplatizedType(s.OBJECT_TYPE, s.STRING_TYPE),
		s.CreateTemplatizedType(s.OBJECT_TYPE, s.NUMBER_TYPE),
		s.CreateTemplatizedType(s.ARRAY_TYPE, s.STRING_TYPE),
		s.CreateTemplatizedType(s.ARRAY_TYPE, s.NUMBER_TYPE),
		s.CreateUnionType(s.CreateTemplatizedType(s.ARRAY_TYPE, s.BOOLEAN_TYPE), s.NULL_TYPE),
		s.CreateUnionType(s.CreateTemplatizedType(s.OBJECT_TYPE, s.BOOLEAN_TYPE), s.NULL_TYPE)}
}

func TestJSType_SymmetryOfTestForEquality(t *testing.T) {
	s := newJsTypeTest(t)
	listA := s.GetTypesToTestForSymmetry()
	listB := s.GetTypesToTestForSymmetry()
	for _, typeA := range listA {
		for _, typeB := range listB {
			aOnB := typeA.TestForEquality(typeB)
			bOnA := typeB.TestForEquality(typeA)
			s.AssertWithMessage(
				rhinotesting.BaseJSTypeTestCase_Lines(
					"testForEquality not symmetrical:",
					"typeA: %s\ntypeB: %s",
					"a.testForEquality(b): %s",
					"b.testForEquality(a): %s"),
				typeA,
				typeB,
				aOnB,
				bOnA).
				That(aOnB == bOnA).
				IsTrue()
		}
	}
}

/** Tests that getLeastSupertype is a symmetric relation. */
func TestJSType_SymmetryOfLeastSupertype(t *testing.T) {
	s := newJsTypeTest(t)
	listA := s.GetTypesToTestForSymmetry()
	listB := s.GetTypesToTestForSymmetry()
	for _, typeA := range listA {
		for _, typeB := range listB {
			aOnB := typeA.GetLeastSupertype(typeB)
			bOnA := typeB.GetLeastSupertype(typeA)

			// TODO(b/110226422): This should use `assertTypeEquals` but at least one of the cases
			// doesn't have matching `hashCode()` values.
			s.AssertWithMessage(
				rhinotesting.BaseJSTypeTestCase_Lines(
					"getLeastSupertype not symmetrical:",
					"typeA: %s",
					"typeB: %s",
					"a.getLeastSupertype(b): %s",
					"b.getLeastSupertype(a): %s"),
				typeA,
				typeB,
				aOnB,
				bOnA).
				That(aOnB.Equals(bOnA)).
				IsTrue()
		}
	}
}

func TestJSType_WeirdBug(t *testing.T) {
	s := newJsTypeTest(t)
	s.AssertTypeNotEquals(s.googBar, s.googBar.GetInstanceType())
	s.AssertThat(s.googBar.IsSubtype(s.googBar.GetInstanceType())).IsFalse()
	s.AssertThat(s.googBar.GetInstanceType().IsSubtypeOf(s.googBar)).IsFalse()
}

/** Tests that getGreatestSubtype is a symmetric relation. */
func TestJSType_SymmetryOfGreatestSubtype(t *testing.T) {
	s := newJsTypeTest(t)
	listA := s.GetTypesToTestForSymmetry()
	listB := s.GetTypesToTestForSymmetry()
	for _, typeA := range listA {
		for _, typeB := range listB {
			aOnB := typeA.GetGreatestSubtype(typeB)
			bOnA := typeB.GetGreatestSubtype(typeA)

			// TODO(b/110226422): This should use `assertTypeEquals` but at least one of the cases
			// doesn't have matching `hashCode()` values.
			s.AssertWithMessage(
				rhinotesting.BaseJSTypeTestCase_Lines(
					"getGreatestSubtype not symmetrical:",
					"typeA: %s",
					"typeB: %s",
					"a.getGreatestSubtype(b): %s",
					"b.getGreatestSubtype(a): %s"),
				typeA,
				typeB,
				aOnB,
				bOnA).
				That(aOnB.Equals(bOnA)).
				IsTrue()
		}
	}
}

/** Tests that getLeastSupertype is a reflexive relation. */
func TestJSType_ReflexivityOfLeastSupertype(t *testing.T) {
	s := newJsTypeTest(t)
	list := s.GetTypesToTestForSymmetry()
	for _, type_ := range list {
		s.AssertTypeEqualsMsg("getLeastSupertype not reflexive",
			type_, type_.GetLeastSupertype(type_))
	}
}

/** Tests that getGreatestSubtype is a reflexive relation. */
func TestJSType_ReflexivityOfGreatestSubtype(t *testing.T) {
	s := newJsTypeTest(t)
	list := s.GetTypesToTestForSymmetry()
	for _, type_ := range list {
		s.AssertTypeEqualsMsg("getGreatestSubtype not reflexive",
			type_, type_.GetGreatestSubtype(type_))
	}
}

/** Tests {@link JSType#getLeastSupertype(JSType)} for unresolved named types. */
func TestJSType_LeastSupertypeUnresolvedNamedType(t *testing.T) {
	s := newJsTypeTest(t)
	// (undefined,function(?):?) and ? unresolved named type
	expected := s.Registry.CreateUnionType(s.unresolvedNamedType, s.FUNCTION_TYPE)
	s.AssertTypeEquals(expected, s.unresolvedNamedType.GetLeastSupertype(s.FUNCTION_TYPE))
	s.AssertTypeEquals(expected, s.FUNCTION_TYPE.GetLeastSupertype(s.unresolvedNamedType))
	s.AssertThat(expected.String()).IsEqualTo("(?|Function)")
}

func TestJSType_LeastSupertypeUnresolvedNamedType2(t *testing.T) {
	s := newJsTypeTest(t)
	expected := s.Registry.CreateUnionType(
		s.unresolvedNamedType, s.UNKNOWN_TYPE)
	s.AssertTypeEquals(expected,
		s.unresolvedNamedType.GetLeastSupertype(s.UNKNOWN_TYPE))
	s.AssertTypeEquals(expected,
		s.UNKNOWN_TYPE.GetLeastSupertype(s.unresolvedNamedType))
	s.AssertTypeEquals(s.UNKNOWN_TYPE, expected)
}

func TestJSType_LeastSupertypeUnresolvedNamedType3(t *testing.T) {
	s := newJsTypeTest(t)
	expected := s.Registry.CreateUnionType(
		s.unresolvedNamedType, s.CHECKED_UNKNOWN_TYPE)
	s.AssertTypeEquals(expected,
		s.unresolvedNamedType.GetLeastSupertype(s.CHECKED_UNKNOWN_TYPE))
	s.AssertTypeEquals(expected,
		s.CHECKED_UNKNOWN_TYPE.GetLeastSupertype(s.unresolvedNamedType))
	s.AssertTypeEquals(s.CHECKED_UNKNOWN_TYPE, expected)
}

/** Tests the subclass of an unresolved named type */
func TestJSType_SubclassOfUnresolvedNamedType(t *testing.T) {
	s := newJsTypeTest(t)
	s.AssertThat(s.subclassOfUnresolvedNamedType.IsUnknownType()).IsTrue()
}

/**
* Tests that Proxied FunctionTypes behave the same over getLeastSupertype and getGreatestSubtype
* as non proxied FunctionTypes
 */
func TestJSType_SupertypeOfProxiedFunctionTypes(t *testing.T) {
	s := newJsTypeTest(t)
	fn1 :=
		FunctionType_Builder_(s.Registry).
			WithParametersOfList(s.Registry.CreateParameters()).
			WithReturnType(s.NUMBER_TYPE).
			Build()
	fn2 :=
		FunctionType_Builder_(s.Registry).
			WithParametersOfList(s.Registry.CreateParameters()).
			WithReturnType(s.STRING_TYPE).
			Build()
	p1 := NewProxyObjectType(s.Registry, fn1)
	p2 := NewProxyObjectType(s.Registry, fn2)
	supremum :=
		FunctionType_Builder_(s.Registry).
			WithParametersOfList(s.Registry.CreateParameters()).
			WithReturnType(s.Registry.CreateUnionType(s.STRING_TYPE, s.NUMBER_TYPE)).
			Build()

	s.AssertTypeEquals(fn1.GetLeastSupertype(fn2), p1.GetLeastSupertype(p2))
	s.AssertTypeEquals(supremum, fn1.GetLeastSupertype(fn2))
	s.AssertTypeEquals(supremum, fn1.GetLeastSupertype(p2))
	s.AssertTypeEquals(supremum, p1.GetLeastSupertype(fn2))
	s.AssertTypeEquals(supremum, p1.GetLeastSupertype(p2))
}

func TestJSType_TypeOfThisIsProxied(t *testing.T) {
	s := newJsTypeTest(t)
	fnType :=
		FunctionType_Builder_(s.Registry).
			WithReturnType(s.NUMBER_TYPE).
			WithTypeOfThis(s.OBJECT_TYPE).
			Build()
	proxyType := NewProxyObjectType(s.Registry, fnType)
	s.AssertTypeEquals(fnType.GetTypeOfThis(), proxyType.GetTypeOfThis())
}

/** Tests the {@link NamedType#equals} function, which had a bug in it. */
func TestJSType_NamedTypeEquals(t *testing.T) {
	s := newJsTypeTest(t)
	closer := s.Registry.GetResolver().OpenForDefinition()
	defer closer.Close()
	{
		// test == if references are equal
		s.ErrorReporter.ExpectAllWarnings("Bad type annotation. Unknown type type1")
		a := s.Registry.CreateNamedType(s.EMPTY_SCOPE, "type1", ptr("source"), 1, 0)
		b := s.Registry.CreateNamedType(s.EMPTY_SCOPE, "type1", ptr("source"), 1, 0)
		s.AssertThat(a.Equals(b)).IsTrue()

		// test == instance of referenced type
		s.AssertThat(s.namedGoogBar.Equals(s.googBar.GetInstanceType())).IsTrue()
		s.AssertThat(s.googBar.GetInstanceType().Equals(s.namedGoogBar)).IsTrue()
	}
}

/** Tests the {@link NamedType#equals} function against other types. */
func TestJSType_NamedTypeEquals2(t *testing.T) {
	s := newJsTypeTest(t)
	// test == if references are equal
	closer := s.Registry.GetResolver().OpenForDefinition()
	a := s.Registry.CreateNamedType(s.EMPTY_SCOPE, "typeA", ptr("source"), 1, 0)
	b := s.Registry.CreateNamedType(s.EMPTY_SCOPE, "typeB", ptr("source"), 1, 0)

	realA :=
		s.Registry.CreateConstructorType(ptr("typeA"), nil, nil, nil, nil, false).GetInstanceType()
	realB :=
		EnumType_Builder_(s.Registry).
			SetName(ptr("typeB")).
			SetElementType(s.NUMBER_TYPE).
			Build().
			GetElementsType()
	s.Registry.DeclareType(nil, "typeA", realA)
	s.Registry.DeclareType(nil, "typeB", realB)
	closer.Close()

	s.AssertThat(a.IsResolved()).IsTrue()
	s.AssertThat(b.IsResolved()).IsTrue()
	s.AssertTypeEquals(a, realA)
	s.AssertTypeEquals(b, realB)

	resolvedA := rhinotesting.Asserts_AssertValidResolve(s.T, a)
	s.AssertThat(a).IsNotSameInstanceAs(resolvedA)
	s.AssertThat(realA).IsSameInstanceAs(resolvedA)

	resolvedB := rhinotesting.Asserts_AssertValidResolve(s.T, b)
	s.AssertThat(b).IsNotSameInstanceAs(resolvedB)
	s.AssertThat(realB).IsSameInstanceAs(resolvedB)
}

func TestJSType_MeaningOfUnresolved(t *testing.T) {
	s := newJsTypeTest(t)
	s.Registry.GetResolver().OpenForDefinition()
	// Given
	underTest := newUnitTestingJSType(s.Registry)

	// When
	// No resolution.

	// Then
	s.AssertThat(underTest.IsResolved()).IsFalse()
	s.AssertThat(underTest.IsSuccessfullyResolved()).IsFalse()
	s.AssertThat(underTest.IsUnsuccessfullyResolved()).IsFalse()
}

func TestJSType_MeaningOfSuccessfullyResolved(t *testing.T) {
	s := newJsTypeTest(t)
	// Given
	underTest :=
		&unitTestingJSType{noResolvedType: ptr(false), typeClass: ptr(JSTypeClass_NO)}
	underTest.InitJSTypeBase(underTest, s.Registry)

	// When
	underTest.EagerlyResolveToSelf()

	// Then
	s.AssertThat(underTest.IsResolved()).IsTrue()
	s.AssertThat(underTest.IsSuccessfullyResolved()).IsTrue()
	s.AssertThat(underTest.IsUnsuccessfullyResolved()).IsFalse()
}

func TestJSType_MeaningOfUnsuccessfullyResolved(t *testing.T) {
	s := newJsTypeTest(t)
	// Given
	underTest :=
		&unitTestingJSType{noResolvedType: ptr(true), typeClass: ptr(JSTypeClass_NO)}
	underTest.InitJSTypeBase(underTest, s.Registry)
	underTest.EagerlyResolveToSelf()

	// When

	// Then
	s.AssertThat(underTest.IsResolved()).IsTrue()
	s.AssertThat(underTest.IsSuccessfullyResolved()).IsFalse()
	s.AssertThat(underTest.IsUnsuccessfullyResolved()).IsTrue()
}

/** Tests {@link JSType#getGreatestSubtype(JSType)} on simple types. */
func TestJSType_GreatestSubtypeSimpleTypes(t *testing.T) {
	s := newJsTypeTest(t)
	s.AssertTypeEquals(s.ARRAY_TYPE,
		s.ARRAY_TYPE.GetGreatestSubtype(s.ALL_TYPE))
	s.AssertTypeEquals(s.ARRAY_TYPE,
		s.ALL_TYPE.GetGreatestSubtype(s.ARRAY_TYPE))
	s.AssertTypeEquals(s.NO_OBJECT_TYPE,
		s.REGEXP_TYPE.GetGreatestSubtype(s.NO_OBJECT_TYPE))
	s.AssertTypeEquals(s.NO_OBJECT_TYPE,
		s.NO_OBJECT_TYPE.GetGreatestSubtype(s.REGEXP_TYPE))
	s.AssertTypeEquals(s.NO_OBJECT_TYPE,
		s.ARRAY_TYPE.GetGreatestSubtype(s.STRING_OBJECT_TYPE))
	s.AssertTypeEquals(s.NO_TYPE, s.ARRAY_TYPE.GetGreatestSubtype(s.NUMBER_TYPE))
	s.AssertTypeEquals(s.NO_OBJECT_TYPE,
		s.ARRAY_TYPE.GetGreatestSubtype(s.functionType))
	s.AssertTypeEquals(s.STRING_OBJECT_TYPE,
		s.STRING_OBJECT_TYPE.GetGreatestSubtype(s.OBJECT_TYPE))
	s.AssertTypeEquals(s.STRING_OBJECT_TYPE,
		s.OBJECT_TYPE.GetGreatestSubtype(s.STRING_OBJECT_TYPE))
	s.AssertTypeEquals(s.NO_OBJECT_TYPE,
		s.ARRAY_TYPE.GetGreatestSubtype(s.DATE_TYPE))
	s.AssertTypeEquals(s.NO_OBJECT_TYPE,
		s.ARRAY_TYPE.GetGreatestSubtype(s.REGEXP_TYPE))
	s.AssertTypeEquals(s.NO_TYPE,
		s.NULL_TYPE.GetGreatestSubtype(s.ARRAY_TYPE))
	s.AssertTypeEquals(s.NUMBER_TYPE,
		s.NUMBER_TYPE.GetGreatestSubtype(s.UNKNOWN_TYPE))
	s.AssertTypeEquals(s.NUMBER_TYPE,
		s.NUMBER_TYPE.GetGreatestSubtype(s.CHECKED_UNKNOWN_TYPE))

	s.AssertTypeEquals(s.NO_OBJECT_TYPE, s.NO_OBJECT_TYPE.GetGreatestSubtype(s.forwardDeclaredNamedType))
	s.AssertTypeEquals(s.NO_OBJECT_TYPE, s.forwardDeclaredNamedType.GetGreatestSubtype(s.NO_OBJECT_TYPE))

	s.AssertTypeEquals(s.CHECKED_UNKNOWN_TYPE,
		s.CHECKED_UNKNOWN_TYPE.GetGreatestSubtype(s.CHECKED_UNKNOWN_TYPE))
	s.AssertTypeEquals(s.CHECKED_UNKNOWN_TYPE,
		s.CHECKED_UNKNOWN_TYPE.GetGreatestSubtype(s.UNKNOWN_TYPE))

}

/** Tests that a derived class extending a type via a named type is a subtype of it. */
func TestJSType_SubtypingDerivedExtendsNamedBaseType(t *testing.T) {
	s := newJsTypeTest(t)
	derived :=
		s.Registry.CreateObjectType(nil, s.Registry.CreateObjectType(nil, s.namedGoogBar))

	s.AssertThat(derived.IsSubtypeOf(s.googBar.GetInstanceType())).IsTrue()
}

func TestJSType_NamedSubtypeChain(t *testing.T) {
	s := newJsTypeTest(t)
	typeChain :=
		[]JSType{
			s.Registry.GetNativeType(JSTypeNative_ALL_TYPE),
			s.Registry.GetNativeType(JSTypeNative_OBJECT_PROTOTYPE),
			s.Registry.GetNativeType(JSTypeNative_OBJECT_TYPE),
			s.googBar.GetPrototype(),
			s.googBar.GetInstanceType(),
			s.googSubBar.GetPrototype(),
			s.googSubBar.GetInstanceType(),
			s.googSubSubBar.GetPrototype(),
			s.googSubSubBar.GetInstanceType(),
			s.Registry.GetNativeType(JSTypeNative_NO_OBJECT_TYPE),
			s.Registry.GetNativeType(JSTypeNative_NO_TYPE)}
	s.VerifySubtypeChain(typeChain)
}

func TestJSType_RecordSubtypeChain(t *testing.T) {
	s := newJsTypeTest(t)
	builder := NewRecordTypeBuilder(s.Registry)
	builder.AddProperty("a", s.STRING_TYPE, nil)
	aType := builder.Build()

	builder = NewRecordTypeBuilder(s.Registry)
	builder.AddProperty("a", s.STRING_TYPE, nil)
	builder.AddProperty("b", s.STRING_TYPE, nil)
	abType := builder.Build()

	builder = NewRecordTypeBuilder(s.Registry)
	builder.AddProperty("a", s.STRING_TYPE, nil)
	builder.AddProperty("c", s.STRING_TYPE, nil)
	acType := builder.Build()
	abOrAcType := s.Registry.CreateUnionType(abType, acType)

	builder = NewRecordTypeBuilder(s.Registry)
	builder.AddProperty("a", s.STRING_TYPE, nil)
	builder.AddProperty("b", s.STRING_TYPE, nil)
	builder.AddProperty("c", s.NUMBER_TYPE, nil)
	abcType := builder.Build()

	typeChain :=
		[]JSType{
			s.Registry.GetNativeType(JSTypeNative_ALL_TYPE),
			s.Registry.GetNativeType(JSTypeNative_OBJECT_PROTOTYPE),
			s.Registry.GetNativeType(JSTypeNative_OBJECT_TYPE),
			aType,
			abOrAcType,
			abType,
			abcType,
			s.Registry.GetNativeType(JSTypeNative_NO_OBJECT_TYPE),
			s.Registry.GetNativeType(JSTypeNative_NO_TYPE)}
	s.VerifySubtypeChain(typeChain)
}

func TestJSType_RecordAndObjectChain2(t *testing.T) {
	s := newJsTypeTest(t)
	builder := NewRecordTypeBuilder(s.Registry)
	builder.AddProperty("date", s.DATE_TYPE, nil)
	hasDateProperty := builder.Build()

	typeChain :=
		[]JSType{
			s.Registry.GetNativeType(JSTypeNative_OBJECT_TYPE),
			hasDateProperty,
			s.googBar.GetInstanceType(),
			s.Registry.GetNativeType(JSTypeNative_NO_OBJECT_TYPE),
			s.Registry.GetNativeType(JSTypeNative_NO_TYPE)}
	s.VerifySubtypeChain(typeChain)
}

func TestJSType_RecordAndObjectChain3(t *testing.T) {
	s := newJsTypeTest(t)
	builder := NewRecordTypeBuilder(s.Registry)
	builder.AddProperty("date", s.UNKNOWN_TYPE, nil)
	hasUnknownDateProperty := builder.Build()

	typeChain :=
		[]JSType{
			s.Registry.GetNativeType(JSTypeNative_OBJECT_TYPE),
			hasUnknownDateProperty,
			s.googBar.GetInstanceType(),
			s.Registry.GetNativeType(JSTypeNative_NO_OBJECT_TYPE),
			s.Registry.GetNativeType(JSTypeNative_NO_TYPE)}
	s.VerifySubtypeChain(typeChain)
}

func TestJSType_NullableNamedTypeChain(t *testing.T) {
	s := newJsTypeTest(t)
	typeChain :=
		[]JSType{
			s.Registry.CreateOptionalNullableType(s.Registry.GetNativeType(JSTypeNative_ALL_TYPE)),
			s.Registry.CreateOptionalNullableType(
				s.Registry.GetNativeType(JSTypeNative_OBJECT_PROTOTYPE)),
			s.Registry.CreateOptionalNullableType(s.Registry.GetNativeType(JSTypeNative_OBJECT_TYPE)),
			s.Registry.CreateOptionalNullableType(s.googBar.GetPrototype()),
			s.Registry.CreateOptionalNullableType(s.googBar.GetInstanceType()),
			s.Registry.CreateNullableType(s.googSubBar.GetPrototype()),
			s.Registry.CreateNullableType(s.googSubBar.GetInstanceType()),
			s.googSubSubBar.GetPrototype(),
			s.googSubSubBar.GetInstanceType(),
			s.Registry.GetNativeType(JSTypeNative_NO_OBJECT_TYPE),
			s.Registry.GetNativeType(JSTypeNative_NO_TYPE)}
	s.VerifySubtypeChain(typeChain)
}

func TestJSType_EnumTypeChain(t *testing.T) {
	s := newJsTypeTest(t)
	typeChain :=
		[]JSType{
			s.Registry.GetNativeType(JSTypeNative_ALL_TYPE),
			s.Registry.GetNativeType(JSTypeNative_OBJECT_PROTOTYPE),
			s.Registry.GetNativeType(JSTypeNative_OBJECT_TYPE),
			s.enumType,
			s.Registry.GetNativeType(JSTypeNative_NO_OBJECT_TYPE),
			s.Registry.GetNativeType(JSTypeNative_NO_TYPE)}
	s.VerifySubtypeChain(typeChain)
}

func TestJSType_FunctionSubtypeChain(t *testing.T) {
	s := newJsTypeTest(t)
	typeChain :=
		[]JSType{
			s.Registry.GetNativeType(JSTypeNative_ALL_TYPE),
			s.Registry.GetNativeType(JSTypeNative_OBJECT_PROTOTYPE),
			s.Registry.GetNativeType(JSTypeNative_OBJECT_TYPE),
			s.Registry.GetNativeType(JSTypeNative_FUNCTION_PROTOTYPE),
			s.Registry.GetNativeType(JSTypeNative_GREATEST_FUNCTION_TYPE),
			s.dateMethod,
			s.Registry.GetNativeType(JSTypeNative_LEAST_FUNCTION_TYPE),
			s.Registry.GetNativeType(JSTypeNative_NO_OBJECT_TYPE),
			s.Registry.GetNativeType(JSTypeNative_NO_TYPE)}
	s.VerifySubtypeChain(typeChain)
}

func TestJSType_FunctionUnionSubtypeChain(t *testing.T) {
	s := newJsTypeTest(t)
	typeChain :=
		[]JSType{
			s.CreateUnionType(s.OBJECT_TYPE, s.STRING_TYPE),
			s.CreateUnionType(s.GREATEST_FUNCTION_TYPE, s.googBarInst, s.STRING_TYPE),
			s.CreateUnionType(
				s.STRING_TYPE,
				s.Registry.CreateFunctionType(s.CreateUnionType(s.STRING_TYPE, s.NUMBER_TYPE)),
				s.googBarInst),
			s.CreateUnionType(s.Registry.CreateFunctionType(s.NUMBER_TYPE), s.googSubBarInst),
			s.LEAST_FUNCTION_TYPE,
			s.NO_OBJECT_TYPE,
			s.NO_TYPE}
	s.VerifySubtypeChain(typeChain)
}

func TestJSType_ConstructorSubtypeChain(t *testing.T) {
	s := newJsTypeTest(t)
	typeChain :=
		[]JSType{
			s.Registry.GetNativeType(JSTypeNative_ALL_TYPE),
			s.Registry.GetNativeType(JSTypeNative_OBJECT_PROTOTYPE),
			s.Registry.GetNativeType(JSTypeNative_OBJECT_TYPE),
			s.Registry.GetNativeType(JSTypeNative_FUNCTION_PROTOTYPE),
			s.Registry.GetNativeType(JSTypeNative_FUNCTION_TYPE),
			s.Registry.GetNativeType(JSTypeNative_NO_OBJECT_TYPE),
			s.Registry.GetNativeType(JSTypeNative_NO_TYPE)}
	s.VerifySubtypeChain(typeChain)
}

func TestJSType_GoogBarSubtypeChain(t *testing.T) {
	s := newJsTypeTest(t)
	typeChain :=
		[]JSType{
			s.Registry.GetNativeType(JSTypeNative_FUNCTION_TYPE),
			s.googBar,
			s.googSubBar,
			s.googSubSubBar,
			s.Registry.GetNativeType(JSTypeNative_NO_OBJECT_TYPE)}
	s.VerifySubtypeChainWithCheck(typeChain, false)
}

func TestJSType_ConstructorWithArgSubtypeChain(t *testing.T) {
	s := newJsTypeTest(t)
	googBarArgConstructor :=
		withOpenRegistry(s, func() FunctionType {
			return s.Registry.CreateConstructorType(
				ptr("barArg"),
				/* source= */ nil,
				s.Registry.CreateParameters(s.googBar),
				nil,
				nil,
				/* isAbstract= */ false)
		})

	googSubBarArgConstructor :=
		withOpenRegistry(s, func() FunctionType {
			return s.Registry.CreateConstructorType(
				ptr("subBarArg"),
				/* source= */ nil,
				s.Registry.CreateParameters(s.googSubBar),
				nil,
				nil,
				/* isAbstract= */ false)
		})

	typeChain :=
		[]JSType{
			s.Registry.GetNativeType(JSTypeNative_FUNCTION_TYPE),
			googBarArgConstructor,
			googSubBarArgConstructor,
			s.Registry.GetNativeType(JSTypeNative_NO_OBJECT_TYPE)}
	s.VerifySubtypeChainWithCheck(typeChain, false)
}

func TestJSType_InterfaceInstanceSubtypeChain(t *testing.T) {
	s := newJsTypeTest(t)
	typeChain :=
		[]JSType{
			s.ALL_TYPE,
			s.OBJECT_TYPE,
			s.interfaceInstType,
			s.googBar.GetPrototype(),
			s.googBarInst,
			s.googSubBar.GetPrototype(),
			s.googSubBarInst,
			s.Registry.GetNativeType(JSTypeNative_NO_OBJECT_TYPE),
			s.Registry.GetNativeType(JSTypeNative_NO_TYPE)}
	s.VerifySubtypeChain(typeChain)
}

func TestJSType_InterfaceInheritanceSubtypeChain(t *testing.T) {
	s := newJsTypeTest(t)
	tempType :=
		withOpenRegistry(s, func() FunctionType {
			return s.Registry.CreateConstructorType(
				ptr("goog.TempType"),
				/* source= */ nil,
				nil,
				nil,
				nil,
				/* isAbstract= */ false)
		})
	tempType.SetImplementedInterfaces(
		[]ObjectType{s.subInterfaceInstType})
	typeChain :=
		[]JSType{
			s.ALL_TYPE,
			s.OBJECT_TYPE,
			s.interfaceInstType,
			s.subInterfaceInstType,
			tempType.GetPrototype(),
			tempType.GetInstanceType(),
			s.Registry.GetNativeType(JSTypeNative_NO_OBJECT_TYPE),
			s.Registry.GetNativeType(JSTypeNative_NO_TYPE)}
	s.VerifySubtypeChain(typeChain)
}

func TestJSType_AnonymousObjectChain(t *testing.T) {
	s := newJsTypeTest(t)
	typeChain :=
		[]JSType{
			s.ALL_TYPE,
			s.CreateNullableType(s.OBJECT_TYPE),
			s.OBJECT_TYPE,
			s.Registry.CreateAnonymousObjectType(nil),
			s.Registry.GetNativeType(JSTypeNative_NO_OBJECT_TYPE),
			s.Registry.GetNativeType(JSTypeNative_NO_TYPE)}
	s.VerifySubtypeChain(typeChain)
}

func TestJSType_AnonymousEnumElementChain(t *testing.T) {
	s := newJsTypeTest(t)
	enumElemType :=
		EnumType_Builder_(s.Registry).
			SetName(ptr("typeB")).
			SetElementType(s.Registry.CreateAnonymousObjectType(nil)).
			Build().
			GetElementsType()
	typeChain :=
		[]JSType{
			s.ALL_TYPE,
			s.CreateNullableType(s.OBJECT_TYPE),
			s.OBJECT_TYPE,
			enumElemType,
			s.Registry.GetNativeType(JSTypeNative_NO_OBJECT_TYPE),
			s.Registry.GetNativeType(JSTypeNative_NO_TYPE)}
	s.VerifySubtypeChain(typeChain)
}

func TestJSType_TemplatizedArrayChain(t *testing.T) {
	s := newJsTypeTest(t)
	arrayOfNoType := s.CreateTemplatizedType(
		s.ARRAY_TYPE, s.NO_TYPE)
	arrayOfString := s.CreateTemplatizedType(
		s.ARRAY_TYPE, s.STRING_TYPE)
	arrayOfStringOrNumber := s.CreateTemplatizedType(
		s.ARRAY_TYPE, s.CreateUnionType(s.STRING_TYPE, s.NUMBER_TYPE))
	arrayOfAllType := s.CreateTemplatizedType(
		s.ARRAY_TYPE, s.ALL_TYPE)

	typeChain :=
		[]JSType{
			s.Registry.GetNativeType(JSTypeNative_ALL_TYPE),
			s.Registry.GetNativeType(JSTypeNative_OBJECT_TYPE),
			arrayOfAllType,
			arrayOfStringOrNumber,
			arrayOfString,
			arrayOfNoType,
			s.Registry.GetNativeType(JSTypeNative_NO_OBJECT_TYPE),
			s.Registry.GetNativeType(JSTypeNative_NO_TYPE)}
	s.VerifySubtypeChainWithCheck(typeChain, false)
}

func TestJSType_TemplatizedArrayChain2(t *testing.T) {
	s := newJsTypeTest(t)
	arrayOfNoType := s.CreateTemplatizedType(
		s.ARRAY_TYPE, s.NO_TYPE)
	arrayOfNoObjectType := s.CreateTemplatizedType(
		s.ARRAY_TYPE, s.NO_OBJECT_TYPE)
	arrayOfArray := s.CreateTemplatizedType(
		s.ARRAY_TYPE, s.ARRAY_TYPE)
	arrayOfObject := s.CreateTemplatizedType(
		s.ARRAY_TYPE, s.OBJECT_TYPE)
	arrayOfAllType := s.CreateTemplatizedType(
		s.ARRAY_TYPE, s.ALL_TYPE)

	typeChain :=
		[]JSType{
			s.Registry.GetNativeType(JSTypeNative_ALL_TYPE),
			s.Registry.GetNativeType(JSTypeNative_OBJECT_TYPE),
			arrayOfAllType,
			arrayOfObject,
			arrayOfArray,
			arrayOfNoObjectType,
			arrayOfNoType,
			s.Registry.GetNativeType(JSTypeNative_NO_OBJECT_TYPE),
			s.Registry.GetNativeType(JSTypeNative_NO_TYPE)}
	s.VerifySubtypeChainWithCheck(typeChain, false)
}

func TestJSType_TemplatizedObjectChain(t *testing.T) {
	s := newJsTypeTest(t)
	objectOfNoType := s.CreateTemplatizedType(
		s.OBJECT_TYPE, s.NO_TYPE)
	objectOfString := s.CreateTemplatizedType(
		s.OBJECT_TYPE, s.STRING_TYPE)
	objectOfStringOrNumber := s.CreateTemplatizedType(
		s.OBJECT_TYPE, s.CreateUnionType(s.STRING_TYPE, s.NUMBER_TYPE))
	objectOfAllType := s.CreateTemplatizedType(
		s.OBJECT_TYPE, s.ALL_TYPE)

	typeChain :=
		[]JSType{
			s.Registry.GetNativeType(JSTypeNative_ALL_TYPE),
			s.Registry.GetNativeType(JSTypeNative_OBJECT_TYPE),
			objectOfAllType,
			objectOfStringOrNumber,
			objectOfString,
			objectOfNoType,
			s.Registry.GetNativeType(JSTypeNative_NO_OBJECT_TYPE),
			s.Registry.GetNativeType(JSTypeNative_NO_TYPE)}
	s.VerifySubtypeChainWithCheck(typeChain, false)
}

func TestJSType_MixedTemplatizedTypeChain(t *testing.T) {
	s := newJsTypeTest(t)
	arrayOfNoType := s.CreateTemplatizedType(
		s.ARRAY_TYPE, s.NO_TYPE)
	arrayOfString := s.CreateTemplatizedType(
		s.ARRAY_TYPE, s.STRING_TYPE)
	objectOfString := s.CreateTemplatizedType(
		s.OBJECT_TYPE, s.STRING_TYPE)
	objectOfStringOrNumber := s.CreateTemplatizedType(
		s.OBJECT_TYPE, s.CreateUnionType(s.STRING_TYPE, s.NUMBER_TYPE))
	objectOfAllType := s.CreateTemplatizedType(
		s.OBJECT_TYPE, s.ALL_TYPE)

	typeChain :=
		[]JSType{
			s.Registry.GetNativeType(JSTypeNative_ALL_TYPE),
			s.Registry.GetNativeType(JSTypeNative_OBJECT_TYPE),
			objectOfAllType,
			objectOfStringOrNumber,
			objectOfString,
			arrayOfString,
			arrayOfNoType,
			s.Registry.GetNativeType(JSTypeNative_NO_OBJECT_TYPE),
			s.Registry.GetNativeType(JSTypeNative_NO_TYPE)}
	s.VerifySubtypeChainWithCheck(typeChain, false)
}

func TestJSType_TemplatizedTypeSubtypes(t *testing.T) {
	s := newJsTypeTest(t)
	objectOfString := s.CreateTemplatizedType(
		s.OBJECT_TYPE, s.STRING_TYPE)
	arrayOfString := s.CreateTemplatizedType(
		s.ARRAY_TYPE, s.STRING_TYPE)
	arrayOfNumber := s.CreateTemplatizedType(
		s.ARRAY_TYPE, s.NUMBER_TYPE)
	arrayOfUnknown := s.CreateTemplatizedType(
		s.ARRAY_TYPE, s.UNKNOWN_TYPE)

	s.AssertThat(objectOfString.IsSubtypeOf(s.ARRAY_TYPE)).IsFalse()
	// TODO(johnlenz): should this be false?
	s.AssertThat(s.ARRAY_TYPE.IsSubtypeOf(objectOfString)).IsTrue()
	s.AssertThat(objectOfString.IsSubtypeOf(s.ARRAY_TYPE)).IsFalse()
	// TODO(johnlenz): should this be false?
	s.AssertThat(s.ARRAY_TYPE.IsSubtypeOf(objectOfString)).IsTrue()

	s.AssertThat(arrayOfString.IsSubtypeOf(s.ARRAY_TYPE)).IsTrue()
	s.AssertThat(s.ARRAY_TYPE.IsSubtypeOf(arrayOfString)).IsTrue()
	s.AssertThat(arrayOfString.IsSubtypeOf(arrayOfUnknown)).IsTrue()
	s.AssertThat(arrayOfUnknown.IsSubtypeOf(arrayOfString)).IsTrue()

	s.AssertThat(arrayOfString.IsSubtypeOf(arrayOfNumber)).IsFalse()
	s.AssertThat(arrayOfNumber.IsSubtypeOf(arrayOfString)).IsFalse()

	s.AssertThat(arrayOfNumber.IsSubtypeOf(s.CreateUnionType(arrayOfNumber, s.NULL_VOID))).IsTrue()
	s.AssertThat(s.CreateUnionType(arrayOfNumber, s.NULL_VOID).IsSubtypeOf(arrayOfNumber)).IsFalse()
	s.AssertThat(arrayOfString.IsSubtypeOf(s.CreateUnionType(arrayOfNumber, s.NULL_VOID))).IsFalse()
}

func TestJSType_TemplatizedTypeRelations(t *testing.T) {
	s := newJsTypeTest(t)
	objectOfString := s.CreateTemplatizedType(
		s.OBJECT_TYPE, s.STRING_TYPE)
	arrayOfString := s.CreateTemplatizedType(
		s.ARRAY_TYPE, s.STRING_TYPE)
	arrayOfNumber := s.CreateTemplatizedType(
		s.ARRAY_TYPE, s.NUMBER_TYPE)

	// Union and least super type cases:
	//
	// 1) alternate:Array<string> and current:Object ==> Object
	// 2) alternate:Array<string> and current:Array ==> Array
	// 3) alternate:Object<string> and current:Array ==> Array|Object<string>
	// 4) alternate:Object and current:Array<string> ==> Object
	// 5) alternate:Array and current:Array<string> ==> Array
	// 6) alternate:Array and current:Object<string> ==> Array|Object<string>
	// 7) alternate:Array<string> and current:Array<number> ==> Array<?>
	// 8) alternate:Array<string> and current:Array<string> ==> Array<string>
	// 9) alternate:Array<string> and
	//    current:Object<string> ==> Object<string>|Array<string>

	s.AssertTypeEquals(
		s.OBJECT_TYPE,
		JSType.GetLeastSupertype(arrayOfString, s.OBJECT_TYPE))
	s.AssertTypeEquals(
		s.OBJECT_TYPE,
		JSType.GetLeastSupertype(s.OBJECT_TYPE, arrayOfString))

	s.AssertTypeEquals(
		s.ARRAY_TYPE,
		JSType.GetLeastSupertype(arrayOfString, s.ARRAY_TYPE))
	s.AssertTypeEquals(
		s.ARRAY_TYPE,
		JSType.GetLeastSupertype(s.ARRAY_TYPE, arrayOfString))

	s.AssertThat(JSType.GetLeastSupertype(objectOfString, s.ARRAY_TYPE).String()).
		IsEqualTo("(Array|Object<string,?>)")
	s.AssertThat(JSType.GetLeastSupertype(s.ARRAY_TYPE, objectOfString).String()).
		IsEqualTo("(Array|Object<string,?>)")

	s.AssertType(JSType.GetLeastSupertype(arrayOfString, arrayOfNumber)).
		ToStringIsEqualTo("Array<?>")
	s.AssertType(JSType.GetLeastSupertype(arrayOfNumber, arrayOfString)).
		ToStringIsEqualTo("Array<?>")

	s.AssertTypeEquals(
		arrayOfString,
		JSType.GetLeastSupertype(arrayOfString, arrayOfString))

	s.AssertThat(JSType.GetLeastSupertype(objectOfString, arrayOfString).String()).
		IsEqualTo("(Array<string>|Object<string,?>)")
	s.AssertThat(JSType.GetLeastSupertype(arrayOfString, objectOfString).String()).
		IsEqualTo("(Array<string>|Object<string,?>)")

	s.AssertTypeEquals(
		objectOfString,
		JSType.GetGreatestSubtype(s.OBJECT_TYPE, objectOfString))

	s.AssertTypeEquals(
		objectOfString,
		JSType.GetGreatestSubtype(objectOfString, s.OBJECT_TYPE))

	s.AssertTypeEquals(
		s.ARRAY_TYPE,
		JSType.GetGreatestSubtype(objectOfString, s.ARRAY_TYPE))

	s.AssertTypeEquals(
		JSType.GetGreatestSubtype(objectOfString, arrayOfString),
		s.NO_OBJECT_TYPE)

	s.AssertTypeEquals(
		JSType.GetGreatestSubtype(s.OBJECT_TYPE, arrayOfString),
		arrayOfString)
}

func TestJSType_TemplatizedTypesWithBoundedGenerics(t *testing.T) {
	s := newJsTypeTest(t)
	templatizedCtor :=
		withOpenRegistry(s, func() FunctionType {
			return s.Registry.CreateConstructorType(
				ptr("TestingType"),
				/* source= */ nil,
				nil,
				s.UNKNOWN_TYPE,
				[]*TemplateType{
					s.Registry.CreateTemplateTypeWithBound("A", s.NUMBER_TYPE),
					s.Registry.CreateTemplateTypeWithBound("B", s.STRING_TYPE)},
				/* isAbstract= */ false)
		})
	templatizedInstance := s.Registry.CreateTemplatizedTypeOfList(
		templatizedCtor.GetInstanceType(),
		[]JSType{s.NUMBER_TYPE})

	ctrTypeMap := templatizedCtor.GetTemplateTypeMap()
	keyA := ctrTypeMap.GetLastTemplateTypeKeyByName("A")
	s.AssertThat(keyA).IsNotNull()
	keyB := ctrTypeMap.GetLastTemplateTypeKeyByName("B")
	s.AssertThat(keyB).IsNotNull()
	keyC := ctrTypeMap.GetLastTemplateTypeKeyByName("C")
	s.AssertThat(keyC).IsNull()
	unknownKey := s.Registry.CreateTemplateType("C")

	templateTypeMap := templatizedInstance.GetTemplateTypeMap()
	s.AssertThat(templateTypeMap.HasTemplateKey(keyA)).IsTrue()
	s.AssertThat(templateTypeMap.HasTemplateKey(keyB)).IsTrue()
	s.AssertThat(templateTypeMap.HasTemplateKey(unknownKey)).IsFalse()

	s.AssertThat(templateTypeMap.GetResolvedTemplateType(keyA)).IsEqualTo(s.NUMBER_TYPE)
	s.AssertThat(templateTypeMap.GetResolvedTemplateType(keyB)).IsEqualTo(keyB)
	s.AssertThat(templateTypeMap.GetResolvedTemplateType(unknownKey)).IsEqualTo(s.UNKNOWN_TYPE)
}

/**
* Tests that the given chain of types has a total ordering defined
* by the subtype relationship, with types at the top of the lattice
* listed first.
*
* Also verifies that the infimum of any two types on the chain
* is the lower type, and the supremum of any two types on the chain
* is the higher type.
 */
func (s *jsTypeTest) VerifySubtypeChain(typeChain []JSType) {
	s.VerifySubtypeChainWithCheck(typeChain, true)
}

func (s *jsTypeTest) VerifySubtypeChainWithCheck(typeChain []JSType, checkSubtyping bool) {
	// Ugh. This wouldn't require so much copy-and-paste if we had a functional
	// programming language.
	for i := 0; i < len(typeChain); i++ {
		for j := 0; j < len(typeChain); j++ {
			typeI := typeChain[i]
			typeJ := typeChain[j]

			namedTypeI := s.GetNamedWrapper("TypeI", typeI)
			namedTypeJ := s.GetNamedWrapper("TypeJ", typeJ)
			proxyTypeI := NewProxyObjectType(s.Registry, typeI)
			proxyTypeJ := NewProxyObjectType(s.Registry, typeJ)

			if i == j {
				s.AssertWithMessage("%s should equal itself", typeI).That(typeI.Equals(typeI)).IsTrue()
				s.AssertWithMessage("Named %s should equal itself", typeI).
					That(namedTypeI.Equals(namedTypeI)).
					IsTrue()
				s.AssertWithMessage("Proxy %s should equal itself", typeI).
					That(proxyTypeI.Equals(proxyTypeI)).
					IsTrue()
			} else {
				shouldCheck := true
				// due to structural interface matching and its updated equivalence
				// checking, a subtype interface could be considered as equivalent
				// to its super type interface (if they are structurally the same)
				// when this happens, the following checks are skipped.
				objectI := typeI.ToObjectType()
				objectJ := typeJ.ToObjectType()
				if objectI != nil && objectJ != nil {
					constructorI := objectI.GetConstructor()
					constructorJ := objectJ.GetConstructor()
					if constructorI != nil && constructorJ != nil &&
						constructorI.IsStructuralInterface() &&
						constructorJ.IsStructuralInterface() {
						if constructorI.Equals(constructorJ) {
							shouldCheck = false
						}
					}
				}
				if shouldCheck {
					s.AssertWithMessage("%s should not equal %s", typeI, typeJ).
						That(typeI.Equals(typeJ)).
						IsFalse()
					s.AssertWithMessage("Named %s should not equal %s", typeI, typeJ).
						That(namedTypeI.Equals(namedTypeJ)).
						IsFalse()
					s.AssertWithMessage("Proxy %s should not equal %s", typeI, typeJ).
						That(proxyTypeI.Equals(proxyTypeJ)).
						IsFalse()
				}
			}

			s.AssertWithMessage("%s should be castable to %s", typeJ, typeI).
				That(typeJ.CanCastTo(typeI)).
				IsTrue()
			s.AssertWithMessage("%s should be castable to Named %s", typeJ, namedTypeI).
				That(typeJ.CanCastTo(namedTypeI)).
				IsTrue()
			s.AssertWithMessage("%s should be castable to Proxy %s", typeJ, proxyTypeI).
				That(typeJ.CanCastTo(proxyTypeI)).
				IsTrue()

			s.AssertWithMessage("Named %s should be castable to %s", typeJ, typeI).
				That(namedTypeJ.CanCastTo(typeI)).
				IsTrue()
			s.AssertWithMessage("Named %s should be castable to Named %s", typeJ, typeI).
				That(namedTypeJ.CanCastTo(namedTypeI)).
				IsTrue()
			s.AssertWithMessage("Named %s should be castable to Proxy %s", typeJ, typeI).
				That(namedTypeJ.CanCastTo(proxyTypeI)).
				IsTrue()

			s.AssertWithMessage("Proxy %s should be castable to %s", typeJ, typeI).
				That(proxyTypeJ.CanCastTo(typeI)).
				IsTrue()
			s.AssertWithMessage("Proxy %s should be castable to Named %s", typeJ, typeI).
				That(proxyTypeJ.CanCastTo(namedTypeI)).
				IsTrue()
			s.AssertWithMessage("Proxy %s should be castable to Proxy %s", typeJ, typeI).
				That(proxyTypeJ.CanCastTo(proxyTypeI)).
				IsTrue()

			// due to structural interface matching, a subtype could be considered
			// as the super type of its super type (if they are structurally the same)
			// when this happens, the following checks are skipped.
			if typeI.IsSubtypeOf(typeJ) && typeJ.IsSubtypeOf(typeI) {
				continue
			}

			if checkSubtyping {
				if i <= j {
					s.AssertWithMessage("%s should be a subtype of %s", typeJ, typeI).
						That(typeJ.IsSubtypeOf(typeI)).
						IsTrue()
					s.AssertWithMessage("Named %s should be a subtype of Named %s", typeJ, typeI).
						That(namedTypeJ.IsSubtypeOf(namedTypeI)).
						IsTrue()
					s.AssertWithMessage("Proxy %s should be a subtype of Proxy %s", typeJ, typeI).
						That(proxyTypeJ.IsSubtypeOf(proxyTypeI)).
						IsTrue()
				} else {
					s.AssertWithMessage("%s should not be a subtype of %s", typeJ, typeI).
						That(typeJ.IsSubtypeOf(typeI)).
						IsFalse()
					s.AssertWithMessage("Named %s should not be a subtype of Named %s", typeJ, typeI).
						That(namedTypeJ.IsSubtypeOf(namedTypeI)).
						IsFalse()
					s.AssertWithMessage("Named %s should not be a subtype of Named %s", typeJ, typeI).
						That(proxyTypeJ.IsSubtypeOf(proxyTypeI)).
						IsFalse()
				}

				expectedSupremum := typeJ
				if i < j {
					expectedSupremum = typeI
				}
				expectedInfimum := typeJ
				if i > j {
					expectedInfimum = typeI
				}
				typeIleastSupertype := typeI.GetLeastSupertype(typeJ)
				typeJleastSupertype := typeJ.GetLeastSupertype(typeI)

				// The `isSubtypeWithoutStructuralTyping` function now considers object literal typedefs
				// as structural types. To ensure accuracy, this check must be extended to accommodate
				// scenarios where the least supertype is a union type.
				//
				// For example, the least supertype of `{a: number}` and `{a: string, b: number}` is
				// now the union type `{a: number} | {a: string, b: number}`.

				if !typeIleastSupertype.IsSubtypeOf(typeJleastSupertype) {
					s.AssertTypeEqualsMsg(
						jutil.ToJavaString(expectedSupremum)+" should be the least supertype of "+jutil.ToJavaString(typeI)+" and "+jutil.ToJavaString(typeJ),
						expectedSupremum,
						typeI.GetLeastSupertype(typeJ))
				}

				// TODO(nicksantos): Should these tests pass?
				//assertTypeEquals(
				//    expectedSupremum + " should be the least supertype of Named " +
				//    typeI + " and Named " + typeJ,
				//    expectedSupremum, namedTypeI.getLeastSupertype(namedTypeJ));
				//assertTypeEquals(
				//    expectedSupremum + " should be the least supertype of Proxy " +
				//    typeI + " and Proxy " + typeJ,
				//    expectedSupremum, proxyTypeI.getLeastSupertype(proxyTypeJ));

				s.AssertTypeEqualsMsg(
					jutil.ToJavaString(expectedInfimum)+" should be the greatest subtype of "+jutil.ToJavaString(typeI)+
						" and "+jutil.ToJavaString(typeJ),
					expectedInfimum, typeI.GetGreatestSubtype(typeJ))

				// TODO(nicksantos): Should these tests pass?
				//assertTypeEquals(
				//    expectedInfimum + " should be the greatest subtype of Named " +
				//    typeI + " and Named " + typeJ,
				//    expectedInfimum, namedTypeI.getGreatestSubtype(namedTypeJ));
				//assertTypeEquals(
				//    expectedInfimum + " should be the greatest subtype of Proxy " +
				//    typeI + " and Proxy " + typeJ,
				//    expectedInfimum, proxyTypeI.getGreatestSubtype(proxyTypeJ));
			}
		}
	}
}

func (s *jsTypeTest) GetNamedWrapper(name string, jstype JSType) JSType {
	// Normally, there is no way to create a Named NoType alias so avoid confusing things by doing
	// it here explicitly.
	if !jstype.IsNoType() {
		return NamedType_Builder_(s.Registry, name).
			SetResolutionKind(NamedType_ResolutionKind_NONE).
			SetReferencedType(jstype).
			Build()
	} else {
		return jstype
	}
}

/** Tests the behavior of {@link JSType#getRestrictedTypeGivenOutcome(Outcome)} ()}. */
func TestJSType_RestrictedTypeGivenOutcome(t *testing.T) {
	s := newJsTypeTest(t)
	// simple cases
	s.AssertTypeEquals(s.BOOLEAN_TYPE,
		s.BOOLEAN_TYPE.GetRestrictedTypeGivenOutcome(Outcome_TRUE))
	s.AssertTypeEquals(s.BOOLEAN_TYPE,
		s.BOOLEAN_TYPE.GetRestrictedTypeGivenOutcome(Outcome_FALSE))

	s.AssertTypeEquals(s.NO_TYPE,
		s.NULL_TYPE.GetRestrictedTypeGivenOutcome(Outcome_TRUE))
	s.AssertTypeEquals(s.NULL_TYPE,
		s.NULL_TYPE.GetRestrictedTypeGivenOutcome(Outcome_FALSE))

	s.AssertTypeEquals(s.NUMBER_TYPE,
		s.NUMBER_TYPE.GetRestrictedTypeGivenOutcome(Outcome_TRUE))
	s.AssertTypeEquals(s.NUMBER_TYPE,
		s.NUMBER_TYPE.GetRestrictedTypeGivenOutcome(Outcome_FALSE))

	s.AssertTypeEquals(s.STRING_TYPE,
		s.STRING_TYPE.GetRestrictedTypeGivenOutcome(Outcome_TRUE))
	s.AssertTypeEquals(s.STRING_TYPE,
		s.STRING_TYPE.GetRestrictedTypeGivenOutcome(Outcome_FALSE))

	s.AssertTypeEquals(s.STRING_OBJECT_TYPE,
		s.STRING_OBJECT_TYPE.GetRestrictedTypeGivenOutcome(Outcome_TRUE))
	s.AssertTypeEquals(s.NO_TYPE,
		s.STRING_OBJECT_TYPE.GetRestrictedTypeGivenOutcome(Outcome_FALSE))

	s.AssertTypeEquals(s.NO_TYPE,
		s.VOID_TYPE.GetRestrictedTypeGivenOutcome(Outcome_TRUE))
	s.AssertTypeEquals(s.VOID_TYPE,
		s.VOID_TYPE.GetRestrictedTypeGivenOutcome(Outcome_FALSE))

	s.AssertTypeEquals(s.NO_OBJECT_TYPE,
		s.NO_OBJECT_TYPE.GetRestrictedTypeGivenOutcome(Outcome_TRUE))
	s.AssertTypeEquals(s.NO_TYPE,
		s.NO_OBJECT_TYPE.GetRestrictedTypeGivenOutcome(Outcome_FALSE))

	s.AssertTypeEquals(s.NO_TYPE,
		s.NO_TYPE.GetRestrictedTypeGivenOutcome(Outcome_TRUE))
	s.AssertTypeEquals(s.NO_TYPE,
		s.NO_TYPE.GetRestrictedTypeGivenOutcome(Outcome_FALSE))

	s.AssertTypeEquals(s.ALL_TYPE,
		s.ALL_TYPE.GetRestrictedTypeGivenOutcome(Outcome_TRUE))
	s.AssertTypeEquals(s.ALL_TYPE,
		s.ALL_TYPE.GetRestrictedTypeGivenOutcome(Outcome_FALSE))

	s.AssertTypeEquals(s.CHECKED_UNKNOWN_TYPE,
		s.UNKNOWN_TYPE.GetRestrictedTypeGivenOutcome(Outcome_TRUE))
	s.AssertTypeEquals(s.UNKNOWN_TYPE,
		s.UNKNOWN_TYPE.GetRestrictedTypeGivenOutcome(Outcome_FALSE))

	// unions
	nullableStringValue :=
		s.CreateNullableType(s.STRING_TYPE).(*UnionType)
	s.AssertTypeEquals(s.STRING_TYPE,
		nullableStringValue.GetRestrictedTypeGivenOutcome(Outcome_TRUE))
	s.AssertTypeEquals(nullableStringValue,
		nullableStringValue.GetRestrictedTypeGivenOutcome(Outcome_FALSE))

	nullableStringObject :=
		s.CreateNullableType(s.STRING_OBJECT_TYPE).(*UnionType)
	s.AssertTypeEquals(s.STRING_OBJECT_TYPE,
		nullableStringObject.GetRestrictedTypeGivenOutcome(Outcome_TRUE))
	s.AssertTypeEquals(s.NULL_TYPE,
		nullableStringObject.GetRestrictedTypeGivenOutcome(Outcome_FALSE))
}

func TestJSType_NullishOutcomeGetRestrictedTypeGivenOutcome(t *testing.T) {
	s := newJsTypeTest(t)
	// simple cases
	s.AssertTypeEquals(s.NO_TYPE,
		s.BOOLEAN_TYPE.GetRestrictedTypeGivenOutcome(Outcome_NULLISH))

	s.AssertTypeEquals(s.NULL_TYPE,
		s.NULL_TYPE.GetRestrictedTypeGivenOutcome(Outcome_NULLISH))

	s.AssertTypeEquals(s.NO_TYPE,
		s.NUMBER_TYPE.GetRestrictedTypeGivenOutcome(Outcome_NULLISH))

	s.AssertTypeEquals(s.NO_TYPE,
		s.STRING_TYPE.GetRestrictedTypeGivenOutcome(Outcome_NULLISH))

	s.AssertTypeEquals(s.NO_TYPE,
		s.STRING_OBJECT_TYPE.GetRestrictedTypeGivenOutcome(Outcome_NULLISH))

	s.AssertTypeEquals(s.VOID_TYPE,
		s.VOID_TYPE.GetRestrictedTypeGivenOutcome(Outcome_NULLISH))

	s.AssertTypeEquals(s.NO_TYPE,
		s.NO_OBJECT_TYPE.GetRestrictedTypeGivenOutcome(Outcome_NULLISH))

	s.AssertTypeEquals(s.NO_TYPE,
		s.NO_TYPE.GetRestrictedTypeGivenOutcome(Outcome_NULLISH))

	s.AssertTypeEquals(s.NO_TYPE,
		s.ALL_TYPE.GetRestrictedTypeGivenOutcome(Outcome_NULLISH))

	s.AssertTypeEquals(s.NO_TYPE,
		s.UNKNOWN_TYPE.GetRestrictedTypeGivenOutcome(Outcome_NULLISH))

	// unions
	nullableStringValue :=
		s.CreateNullableType(s.STRING_TYPE).(*UnionType)
	s.AssertTypeEquals(s.NULL_TYPE,
		nullableStringValue.GetRestrictedTypeGivenOutcome(Outcome_NULLISH))

	nullableStringObject :=
		s.CreateNullableType(s.STRING_OBJECT_TYPE).(*UnionType)
	s.AssertTypeEquals(s.NULL_TYPE,
		nullableStringObject.GetRestrictedTypeGivenOutcome(Outcome_NULLISH))
}

func TestJSType_FalseNotNullOutcomeGetRestrictedTypeGivenOutcome(t *testing.T) {
	s := newJsTypeTest(t)
	// simple cases
	s.AssertTypeEquals(s.BOOLEAN_TYPE,
		s.BOOLEAN_TYPE.GetRestrictedTypeGivenOutcome(Outcome_FALSE_NOT_NULL))

	s.AssertTypeEquals(s.NULL_TYPE,
		s.NULL_TYPE.GetRestrictedTypeGivenOutcome(Outcome_FALSE_NOT_NULL))

	s.AssertTypeEquals(s.NUMBER_TYPE,
		s.NUMBER_TYPE.GetRestrictedTypeGivenOutcome(Outcome_FALSE_NOT_NULL))

	s.AssertTypeEquals(s.STRING_TYPE,
		s.STRING_TYPE.GetRestrictedTypeGivenOutcome(Outcome_FALSE_NOT_NULL))

	s.AssertTypeEquals(s.NO_TYPE,
		s.STRING_OBJECT_TYPE.GetRestrictedTypeGivenOutcome(Outcome_FALSE_NOT_NULL))

	s.AssertTypeEquals(s.VOID_TYPE,
		s.VOID_TYPE.GetRestrictedTypeGivenOutcome(Outcome_FALSE_NOT_NULL))

	s.AssertTypeEquals(s.NO_TYPE,
		s.NO_OBJECT_TYPE.GetRestrictedTypeGivenOutcome(Outcome_FALSE_NOT_NULL))

	s.AssertTypeEquals(s.NO_TYPE,
		s.NO_TYPE.GetRestrictedTypeGivenOutcome(Outcome_FALSE_NOT_NULL))

	s.AssertTypeEquals(s.ALL_TYPE,
		s.ALL_TYPE.GetRestrictedTypeGivenOutcome(Outcome_FALSE_NOT_NULL))

	s.AssertTypeEquals(s.UNKNOWN_TYPE,
		s.UNKNOWN_TYPE.GetRestrictedTypeGivenOutcome(Outcome_FALSE_NOT_NULL))

	// unions
	nullableStringValue :=
		s.CreateNullableType(s.STRING_TYPE).(*UnionType)
	s.AssertTypeEquals(nullableStringValue,
		nullableStringValue.GetRestrictedTypeGivenOutcome(Outcome_FALSE_NOT_NULL))

	nullableStringObject :=
		s.CreateNullableType(s.STRING_OBJECT_TYPE).(*UnionType)
	s.AssertTypeEquals(s.NULL_TYPE,
		nullableStringObject.GetRestrictedTypeGivenOutcome(Outcome_FALSE_NOT_NULL))
}

func TestJSType_GoodSetPrototypeBasedOn(t *testing.T) {
	s := newJsTypeTest(t)
	fun :=
		withOpenRegistry(s, func() FunctionType {
			return s.Registry.CreateConstructorType(
				ptr("fun") /* source= */, nil, nil, nil, nil /* isAbstract= */, false)
		})
	fun.SetPrototypeBasedOn(s.unresolvedNamedType)
	s.AssertThat(fun.GetInstanceType().IsUnknownType()).IsTrue()
}

func TestJSType_LateSetPrototypeBasedOn(t *testing.T) {
	s := newJsTypeTest(t)
	fun :=
		withOpenRegistry(s, func() FunctionType {
			return s.Registry.CreateConstructorType(
				ptr("fun") /* source= */, nil, nil, nil, nil /* isAbstract= */, false)
		})
	s.AssertThat(fun.GetInstanceType().IsUnknownType()).IsFalse()

	fun.SetPrototypeBasedOn(s.unresolvedNamedType)
	s.AssertThat(fun.GetInstanceType().IsUnknownType()).IsTrue()
}

func TestJSType_GetTypeUnderEquality1(t *testing.T) {
	s := newJsTypeTest(t)
	for _, type_ := range s.types {
		s.TestGetTypeUnderEquality(type_, type_, type_, type_)
	}
}

func TestJSType_GetTypesUnderEquality2(t *testing.T) {
	s := newJsTypeTest(t)
	// objects can be equal to numbers
	s.TestGetTypeUnderEquality(
		s.NUMBER_TYPE, s.OBJECT_TYPE,
		s.NUMBER_TYPE, s.OBJECT_TYPE)
}

func TestJSType_GetTypesUnderEquality3(t *testing.T) {
	s := newJsTypeTest(t)
	// null == undefined
	s.TestGetTypeUnderEquality(
		s.NULL_TYPE, s.VOID_TYPE,
		s.NULL_TYPE, s.VOID_TYPE)
}

func TestJSType_GetTypesUnderEquality4(t *testing.T) {
	s := newJsTypeTest(t)
	// (number,string) and number/string
	stringNumber :=
		s.CreateUnionType(s.NUMBER_TYPE, s.STRING_TYPE).(*UnionType)
	s.TestGetTypeUnderEquality(
		stringNumber, s.STRING_TYPE,
		stringNumber, s.STRING_TYPE)
	s.TestGetTypeUnderEquality(
		stringNumber, s.NUMBER_TYPE,
		stringNumber, s.NUMBER_TYPE)
}

func TestJSType_GetTypesUnderEquality5(t *testing.T) {
	s := newJsTypeTest(t)
	// (number,null) and undefined
	nullUndefined := s.CreateUnionType(s.VOID_TYPE, s.NULL_TYPE)
	s.TestGetTypeUnderEquality(
		nullUndefined, s.NULL_TYPE,
		nullUndefined, s.NULL_TYPE)
	s.TestGetTypeUnderEquality(
		nullUndefined, s.VOID_TYPE,
		nullUndefined, s.VOID_TYPE)
}

func TestJSType_GetTypesUnderEquality6(t *testing.T) {
	s := newJsTypeTest(t)
	// (number,undefined,null) == null
	optNullNumber := s.CreateUnionType(s.VOID_TYPE, s.NULL_TYPE, s.NUMBER_TYPE)
	s.TestGetTypeUnderEquality(
		optNullNumber, s.NULL_TYPE,
		s.CreateUnionType(s.NULL_TYPE, s.VOID_TYPE), s.NULL_TYPE)
}

func (s *jsTypeTest) TestGetTypeUnderEquality(
	t1 JSType, t2 JSType, t1Eq JSType, t2Eq JSType) {
	// creating the pairs
	p12 := t1.GetTypesUnderEquality(t2)
	p21 := t2.GetTypesUnderEquality(t1)

	// t1Eq
	s.AssertTypeEquals(t1Eq, p12.TypeA)
	s.AssertTypeEquals(t1Eq, p21.TypeB)

	// t2Eq
	s.AssertTypeEquals(t2Eq, p12.TypeB)
	s.AssertTypeEquals(t2Eq, p21.TypeA)
}

func TestJSType_GetTypesUnderInequality1(t *testing.T) {
	s := newJsTypeTest(t)
	// objects can be not equal to numbers
	numberObject :=
		s.CreateUnionType(s.NUMBER_TYPE, s.OBJECT_TYPE).(*UnionType)
	s.TestGetTypesUnderInequality(
		numberObject, s.NUMBER_TYPE,
		numberObject, s.NUMBER_TYPE)
	s.TestGetTypesUnderInequality(
		numberObject, s.OBJECT_TYPE,
		numberObject, s.OBJECT_TYPE)
}

func TestJSType_GetTypesUnderInequality2(t *testing.T) {
	s := newJsTypeTest(t)
	// null == undefined
	nullUndefined :=
		s.CreateUnionType(s.VOID_TYPE, s.NULL_TYPE).(*UnionType)
	s.TestGetTypesUnderInequality(
		nullUndefined, s.NULL_TYPE,
		s.NO_TYPE, s.NO_TYPE)
	s.TestGetTypesUnderInequality(
		nullUndefined, s.VOID_TYPE,
		s.NO_TYPE, s.NO_TYPE)
}

func TestJSType_GetTypesUnderInequality3(t *testing.T) {
	s := newJsTypeTest(t)
	// (number,string)
	stringNumber :=
		s.CreateUnionType(s.NUMBER_TYPE, s.STRING_TYPE).(*UnionType)
	s.TestGetTypesUnderInequality(
		stringNumber, s.NUMBER_TYPE,
		stringNumber, s.NUMBER_TYPE)
	s.TestGetTypesUnderInequality(
		stringNumber, s.STRING_TYPE,
		stringNumber, s.STRING_TYPE)
}

func TestJSType_GetTypesUnderInequality4(t *testing.T) {
	s := newJsTypeTest(t)
	// (number,undefined,null) and null
	nullableOptionalNumber :=
		s.CreateUnionType(s.NULL_TYPE, s.VOID_TYPE, s.NUMBER_TYPE).(*UnionType)
	s.TestGetTypesUnderInequality(
		nullableOptionalNumber, s.NULL_TYPE,
		s.NUMBER_TYPE, s.NULL_TYPE)
}

func (s *jsTypeTest) TestGetTypesUnderInequality(
	t1 JSType, t2 JSType, t1Eq JSType, t2Eq JSType) {
	// creating the pairs
	p12 := t1.GetTypesUnderInequality(t2)
	p21 := t2.GetTypesUnderInequality(t1)

	// t1Eq
	s.AssertTypeEquals(t1Eq, p12.TypeA)
	s.AssertTypeEquals(t1Eq, p21.TypeB)

	// t2Eq
	s.AssertTypeEquals(t2Eq, p12.TypeB)
	s.AssertTypeEquals(t2Eq, p21.TypeA)
}

/** Tests the factory method {@link JSTypeRegistry#createOptionalType(JSType)}. */
func TestJSType_CreateOptionalType(t *testing.T) {
	s := newJsTypeTest(t)
	// number
	optNumber := s.Registry.CreateOptionalType(s.NUMBER_TYPE).(*UnionType)
	s.AssertUnionContains(optNumber, s.NUMBER_TYPE)
	s.AssertUnionContains(optNumber, s.VOID_TYPE)

	// union
	optUnion :=
		s.Registry.CreateOptionalType(
			s.CreateUnionType(s.STRING_OBJECT_TYPE, s.DATE_TYPE)).(*UnionType)
	s.AssertUnionContains(optUnion, s.DATE_TYPE)
	s.AssertUnionContains(optUnion, s.STRING_OBJECT_TYPE)
	s.AssertUnionContains(optUnion, s.VOID_TYPE)
}

func (s *jsTypeTest) AssertUnionContains(union *UnionType, type_ JSType) {
	s.AssertWithMessage("%s should contain %s", union, type_).That(union.Contains(type_)).IsTrue()
}

/** Tests the factory method {@link JSTypeRegistry#createAnonymousObjectType}}. */
func TestJSType_CreateAnonymousObjectType(t *testing.T) {
	s := newJsTypeTest(t)
	// anonymous
	anonymous := s.Registry.CreateAnonymousObjectType(nil)
	s.AssertTypeEquals(s.OBJECT_TYPE, anonymous.GetImplicitPrototype())
	s.AssertType(anonymous).GetReferenceNameIsNull()
	s.AssertThat(anonymous.String()).IsEqualTo("{}")
}

/**
* Tests the factory method {@link JSTypeRegistry#createAnonymousObjectType}} and adds some
* properties to it.
 */
func TestJSType_CreateAnonymousObjectType2(t *testing.T) {
	s := newJsTypeTest(t)
	// anonymous
	anonymous := s.Registry.CreateAnonymousObjectType(nil)
	anonymous.DefineDeclaredProperty(
		"a", s.NUMBER_TYPE, nil)
	anonymous.DefineDeclaredProperty(
		"b", s.NUMBER_TYPE, nil)
	anonymous.DefineDeclaredProperty(
		"c", s.NUMBER_TYPE, nil)
	anonymous.DefineDeclaredProperty(
		"d", s.NUMBER_TYPE, nil)
	anonymous.DefineDeclaredProperty(
		"e", s.NUMBER_TYPE, nil)
	anonymous.DefineDeclaredProperty(
		"f", s.NUMBER_TYPE, nil)
	s.AssertThat(anonymous.String()).
		IsEqualTo(
			rhinotesting.BaseJSTypeTestCase_Lines(
				"{",
				"  a: number,",
				"  b: number,",
				"  c: number,",
				"  d: number,",
				"  e: number,",
				"  f: number",
				"}"))
}

/** Tests the factory method {@link JSTypeRegistry#createObjectType(String, ObjectType)}}. */
func TestJSType_CreateObjectType(t *testing.T) {
	s := newJsTypeTest(t)
	// simple
	subDate :=
		s.Registry.CreateObjectType(nil, s.DATE_TYPE.GetImplicitPrototype())
	s.AssertTypeEquals(s.DATE_TYPE.GetImplicitPrototype(),
		subDate.GetImplicitPrototype())
	s.AssertType(subDate).GetReferenceNameIsNull()
	s.AssertThat(subDate.String()).IsEqualTo("{...}")

	// name, node, prototype
	subArray := s.Registry.CreateObjectType(ptr("Foo"),
		s.ARRAY_TYPE.GetImplicitPrototype())
	s.AssertTypeEquals(s.ARRAY_TYPE.GetImplicitPrototype(),
		subArray.GetImplicitPrototype())
	s.AssertType(subArray).GetReferenceNameIsEqualTo("Foo")
}

/** Tests {@code (U2U_CONSTRUCTOR,undefined) <: (U2U_CONSTRUCTOR,undefined)}. */
func TestJSType_Bug903110(t *testing.T) {
	s := newJsTypeTest(t)
	union := s.CreateUnionType(s.FUNCTION_TYPE, s.VOID_TYPE).(*UnionType)
	s.AssertThat(s.VOID_TYPE.IsSubtypeOf(union)).IsTrue()
	s.AssertThat(s.FUNCTION_TYPE.IsSubtype(union)).IsTrue()
	s.AssertThat(union.IsSubtype(union)).IsTrue()
}

/**
* Assert that a type can assign to itself.
 */
func (s *jsTypeTest) AssertTypeCanAssignToItself(type_ JSType) {
	s.AssertThat(type_.IsSubtypeOf(type_)).IsTrue()
}

/**
* Tests that hasOwnProperty returns true when a property is defined directly on a class and false
* if the property is defined on the supertype or not at all.
 */
func TestJSType_HasOwnProperty(t *testing.T) {
	s := newJsTypeTest(t)
	sup :=
		s.Registry.CreateObjectType(nil, s.Registry.CreateAnonymousObjectType(nil))
	sub := s.Registry.CreateObjectType(nil, sup)

	sup.DefineProperty("base", s.GetBottomType(), false, nil)
	sub.DefineProperty("sub", s.GetBottomType(), false, nil)

	s.AssertThat(sup.HasProperty("base")).IsTrue()
	s.AssertThat(sup.HasProperty("sub")).IsFalse()
	s.AssertThat(sup.HasOwnProperty("base")).IsTrue()
	s.AssertThat(sup.HasOwnProperty("sub")).IsFalse()
	s.AssertThat(sup.HasOwnProperty("none")).IsFalse()

	s.AssertThat(sub.HasProperty("base")).IsTrue()
	s.AssertThat(sub.HasProperty("sub")).IsTrue()
	s.AssertThat(sub.HasOwnProperty("base")).IsFalse()
	s.AssertThat(sub.HasOwnProperty("sub")).IsTrue()
	s.AssertThat(sub.HasOwnProperty("none")).IsFalse()
}

/**
* Tests that hasOwnProperty returns true when a property is defined directly on a class and false
* if the property is defined on the supertype or not at all.
 */
func TestJSType_HasOwnProperty_symbol(t *testing.T) {
	s := newJsTypeTest(t)
	sup := s.Registry.CreateObjectType(nil, s.Registry.CreateAnonymousObjectType(nil))
	sub := s.Registry.CreateObjectType(nil, sup)
	baseSymbol := NewProperty_SymbolKey(NewKnownSymbolType(s.Registry, "base"))
	subSymbol := NewProperty_SymbolKey(NewKnownSymbolType(s.Registry, "sub"))
	noneSymbol := NewProperty_SymbolKey(NewKnownSymbolType(s.Registry, "none"))

	sup.DefinePropertyByKey(baseSymbol, s.GetBottomType(), false, nil)
	sub.DefinePropertyByKey(subSymbol, s.GetBottomType(), false, nil)

	s.AssertThat(sup.HasPropertyByKey(baseSymbol)).IsTrue()
	s.AssertThat(sup.HasPropertyByKey(subSymbol)).IsFalse()
	s.AssertThat(sup.HasOwnPropertyByKey(baseSymbol)).IsTrue()
	s.AssertThat(sup.HasOwnPropertyByKey(subSymbol)).IsFalse()
	s.AssertThat(sup.HasOwnPropertyByKey(noneSymbol)).IsFalse()
	s.AssertThat(sup.HasOwnProperty("none")).IsFalse()

	s.AssertThat(sub.HasPropertyByKey(baseSymbol)).IsTrue()
	s.AssertThat(sub.HasPropertyByKey(subSymbol)).IsTrue()
	s.AssertThat(sub.HasOwnPropertyByKey(baseSymbol)).IsFalse()
	s.AssertThat(sub.HasOwnPropertyByKey(subSymbol)).IsTrue()
	s.AssertThat(sub.HasOwnPropertyByKey(noneSymbol)).IsFalse()
	s.AssertThat(sub.HasOwnProperty("none")).IsFalse()
}

func TestJSType_NamedTypeHasOwnProperty(t *testing.T) {
	s := newJsTypeTest(t)
	s.namedGoogBar.GetImplicitPrototype().DefineProperty("base", s.GetBottomType(), false, nil)
	s.namedGoogBar.DefineProperty("sub", s.GetBottomType(), false, nil)

	s.AssertThat(s.namedGoogBar.HasOwnProperty("base")).IsFalse()
	s.AssertThat(s.namedGoogBar.HasProperty("base")).IsTrue()
	s.AssertThat(s.namedGoogBar.HasOwnProperty("sub")).IsTrue()
	s.AssertThat(s.namedGoogBar.HasProperty("sub")).IsTrue()
}

func TestJSType_InterfaceHasOwnProperty(t *testing.T) {
	s := newJsTypeTest(t)
	s.interfaceInstType.DefineProperty("base", s.GetBottomType(), false, nil)
	s.subInterfaceInstType.DefineProperty("sub", s.GetBottomType(), false, nil)

	s.AssertThat(s.interfaceInstType.HasProperty("base")).IsTrue()
	s.AssertThat(s.interfaceInstType.HasProperty("sub")).IsFalse()
	s.AssertThat(s.interfaceInstType.HasOwnProperty("base")).IsTrue()
	s.AssertThat(s.interfaceInstType.HasOwnProperty("sub")).IsFalse()
	s.AssertThat(s.interfaceInstType.HasOwnProperty("none")).IsFalse()

	s.AssertThat(s.subInterfaceInstType.HasProperty("base")).IsTrue()
	s.AssertThat(s.subInterfaceInstType.HasProperty("sub")).IsTrue()
	s.AssertThat(s.subInterfaceInstType.HasOwnProperty("base")).IsFalse()
	s.AssertThat(s.subInterfaceInstType.HasOwnProperty("sub")).IsTrue()
	s.AssertThat(s.subInterfaceInstType.HasOwnProperty("none")).IsFalse()
}

func TestJSType_InterfaceHasOwnProperty_symbol(t *testing.T) {
	s := newJsTypeTest(t)
	baseSymbol := NewProperty_SymbolKey(NewKnownSymbolType(s.Registry, "base"))
	subSymbol := NewProperty_SymbolKey(NewKnownSymbolType(s.Registry, "sub"))
	noneSymbol := NewProperty_SymbolKey(NewKnownSymbolType(s.Registry, "none"))
	s.interfaceInstType.DefinePropertyByKey(baseSymbol, s.GetBottomType(), false, nil)
	s.subInterfaceInstType.DefinePropertyByKey(subSymbol, s.GetBottomType(), false, nil)

	s.AssertThat(s.interfaceInstType.HasPropertyByKey(baseSymbol)).IsTrue()
	s.AssertThat(s.interfaceInstType.HasPropertyByKey(subSymbol)).IsFalse()
	s.AssertThat(s.interfaceInstType.HasOwnPropertyByKey(baseSymbol)).IsTrue()
	s.AssertThat(s.interfaceInstType.HasOwnPropertyByKey(subSymbol)).IsFalse()
	s.AssertThat(s.interfaceInstType.HasOwnPropertyByKey(noneSymbol)).IsFalse()
	s.AssertThat(s.interfaceInstType.HasOwnProperty("none")).IsFalse()

	s.AssertThat(s.subInterfaceInstType.HasPropertyByKey(baseSymbol)).IsTrue()
	s.AssertThat(s.subInterfaceInstType.HasPropertyByKey(subSymbol)).IsTrue()
	s.AssertThat(s.subInterfaceInstType.HasOwnPropertyByKey(baseSymbol)).IsFalse()
	s.AssertThat(s.subInterfaceInstType.HasOwnPropertyByKey(subSymbol)).IsTrue()
	s.AssertThat(s.subInterfaceInstType.HasOwnPropertyByKey(noneSymbol)).IsFalse()
	s.AssertThat(s.subInterfaceInstType.HasOwnProperty("none")).IsFalse()
}

func TestJSType_GetPropertyNames(t *testing.T) {
	s := newJsTypeTest(t)
	sup :=
		s.Registry.CreateObjectType(nil, s.Registry.CreateAnonymousObjectType(nil))
	sub := s.Registry.CreateObjectType(nil, sup)
	bottomType := s.GetBottomType()

	sup.DefineProperty("base", bottomType, false, nil)
	sub.DefineProperty("sub", bottomType, false, nil)

	s.AssertThat(sub.GetPropertyNames()).
		ContainsExactlyElementsIn(
			[]string{
				"isPrototypeOf",
				"toLocaleString",
				"propertyIsEnumerable",
				"toString",
				"valueOf",
				"hasOwnProperty",
				"constructor",
				"base",
				"sub"})
	s.AssertThat(sup.GetPropertyNames()).
		ContainsExactlyElementsIn(
			[]string{
				"isPrototypeOf",
				"toLocaleString",
				"propertyIsEnumerable",
				"toString",
				"valueOf",
				"hasOwnProperty",
				"constructor",
				"base"})

	s.AssertThat(s.NO_OBJECT_TYPE.GetPropertyNames()).IsEmpty()
}

func TestJSType_GetPropertyNames_symbol(t *testing.T) {
	s := newJsTypeTest(t)
	sup := s.Registry.CreateObjectType(nil, s.Registry.CreateAnonymousObjectType(nil))
	sub := s.Registry.CreateObjectType(nil, sup)
	bottomType := s.GetBottomType()
	baseSymbol := NewProperty_SymbolKey(NewKnownSymbolType(s.Registry, "base"))
	subSymbol := NewProperty_SymbolKey(NewKnownSymbolType(s.Registry, "sub"))

	sup.DefinePropertyByKey(baseSymbol, bottomType, false, nil)
	sub.DefinePropertyByKey(subSymbol, bottomType, false, nil)

	s.AssertThat(sub.GetPropertyNames()).
		ContainsExactlyElementsIn(
			[]string{
				"isPrototypeOf",
				"toLocaleString",
				"propertyIsEnumerable",
				"toString",
				"valueOf",
				"hasOwnProperty",
				"constructor"})
	s.AssertThat(sup.GetPropertyNames()).
		ContainsExactlyElementsIn(
			[]string{
				"isPrototypeOf",
				"toLocaleString",
				"propertyIsEnumerable",
				"toString",
				"valueOf",
				"hasOwnProperty",
				"constructor"})
	s.AssertThat(sup.GetAllKeys().KnownSymbolKeys()).ContainsExactly(baseSymbol.Symbol())
	s.AssertThat(sub.GetAllKeys().KnownSymbolKeys()).
		ContainsExactly(baseSymbol.Symbol(), subSymbol.Symbol())
	s.AssertThat(sub.GetOwnPropertyKnownSymbols()).ContainsExactly(subSymbol.Symbol())

	s.AssertThat(s.NO_OBJECT_TYPE.GetPropertyNames()).IsEmpty()
}

func (s *jsTypeTest) GetBottomType() ObjectType {
	return s.Registry.GetNativeObjectType(JSTypeNative_NO_TYPE)
}

func TestJSType_GetAndSetJSDocInfoWithNamedType(t *testing.T) {
	s := newJsTypeTest(t)
	builder := JSDocInfo_Builder_()
	builder.RecordDeprecated()
	info := builder.Build()

	s.AssertThat(s.namedGoogBar.GetOwnPropertyJSDocInfo("X")).IsNull()
	s.namedGoogBar.SetPropertyJSDocInfo("X", info)
	s.AssertThat(s.namedGoogBar.GetOwnPropertyJSDocInfo("X").IsDeprecated()).IsTrue()
	s.AssertPropertyTypeInferred(s.namedGoogBar, "X")
	s.AssertTypeEquals(s.UNKNOWN_TYPE, s.namedGoogBar.GetPropertyType("X"))
}

func TestJSType_GetAndSetJSDocInfoWithObjectTypes(t *testing.T) {
	s := newJsTypeTest(t)
	sup :=
		s.Registry.CreateObjectType(nil, s.Registry.CreateAnonymousObjectType(nil))
	sub := s.Registry.CreateObjectType(nil, sup)

	builder := JSDocInfo_Builder_()
	builder.RecordDeprecated()
	deprecated := builder.Build()

	builder = JSDocInfo_Builder_()
	builder.RecordVisibility(JSDocInfo_Visibility_PRIVATE)
	privateInfo := builder.Build()

	sup.DefineProperty("X", s.NUMBER_TYPE, true, nil)
	sup.SetPropertyJSDocInfo("X", privateInfo)

	sub.DefineProperty("X", s.NUMBER_TYPE, true, nil)
	sub.SetPropertyJSDocInfo("X", deprecated)

	s.AssertThat(sup.GetOwnPropertyJSDocInfo("X").IsDeprecated()).IsFalse()
	s.AssertThat(sup.GetOwnPropertyJSDocInfo("X").GetVisibility()).IsEqualTo(JSDocInfo_Visibility_PRIVATE)
	s.AssertTypeEquals(s.NUMBER_TYPE, sup.GetPropertyType("X"))
	s.AssertThat(sub.GetOwnPropertyJSDocInfo("X").IsDeprecated()).IsTrue()
	s.AssertThat(sub.GetOwnPropertyJSDocInfo("X").GetVisibility()).IsEqualTo(JSDocInfo_Visibility_INHERITED)
	s.AssertTypeEquals(s.NUMBER_TYPE, sub.GetPropertyType("X"))
}

func TestJSType_GetAndSetJSDocInfoWithNoType(t *testing.T) {
	s := newJsTypeTest(t)
	builder := JSDocInfo_Builder_()
	builder.RecordDeprecated()
	deprecated := builder.Build()

	s.NO_TYPE.SetPropertyJSDocInfo("X", deprecated)
	s.AssertThat(s.NO_TYPE.GetOwnPropertyJSDocInfo("X")).IsNull()
}

func TestJSType_IsTemplatedType(t *testing.T) {
	s := newJsTypeTest(t)
	s.AssertThat(NewTemplateType(s.Registry, "T").HasAnyTemplateTypes()).IsTrue()
	s.AssertThat(s.ARRAY_TYPE.HasAnyTemplateTypes()).IsFalse()

	s.AssertThat(
		s.CreateTemplatizedType(s.ARRAY_TYPE, NewTemplateType(s.Registry, "T")).
			HasAnyTemplateTypes()).
		IsTrue()
	s.AssertThat(s.CreateTemplatizedType(s.ARRAY_TYPE, s.STRING_TYPE).HasAnyTemplateTypes()).IsFalse()

	s.AssertThat(
		FunctionType_Builder_(s.Registry).
			WithReturnType(NewTemplateType(s.Registry, "T")).
			Build().
			HasAnyTemplateTypes()).
		IsTrue()
	s.AssertThat(
		FunctionType_Builder_(s.Registry).
			WithTypeOfThis(NewTemplateType(s.Registry, "T")).
			Build().
			HasAnyTemplateTypes()).
		IsTrue()
	s.AssertThat(
		FunctionType_Builder_(s.Registry).
			WithReturnType(s.STRING_TYPE).
			Build().
			HasAnyTemplateTypes()).
		IsFalse()

	s.AssertThat(
		s.Registry.
			CreateUnionType(s.NULL_TYPE, NewTemplateType(s.Registry, "T"), s.STRING_TYPE).
			HasAnyTemplateTypes()).
		IsTrue()
	s.AssertThat(s.Registry.CreateUnionType(s.NULL_TYPE, s.ARRAY_TYPE, s.STRING_TYPE).HasAnyTemplateTypes()).
		IsFalse()
}

func TestJSType_TemplatizedType(t *testing.T) {
	s := newJsTypeTest(t)

	templatizedCtor :=
		withOpenRegistry(s, func() FunctionType {
			return s.Registry.CreateConstructorType(
				ptr("TestingType"),
				/* source= */ nil,
				nil,
				s.UNKNOWN_TYPE,
				[]*TemplateType{
					s.Registry.CreateTemplateType("A"), s.Registry.CreateTemplateType("B")},
				/* isAbstract= */ false)
		})
	templatizedInstance := s.Registry.CreateTemplatizedTypeOfList(
		templatizedCtor.GetInstanceType(),
		[]JSType{s.NUMBER_TYPE, s.STRING_TYPE})

	ctrTypeMap := templatizedCtor.GetTemplateTypeMap()
	keyA := ctrTypeMap.GetLastTemplateTypeKeyByName("A")
	s.AssertThat(keyA).IsNotNull()
	keyB := ctrTypeMap.GetLastTemplateTypeKeyByName("B")
	s.AssertThat(keyB).IsNotNull()
	keyC := ctrTypeMap.GetLastTemplateTypeKeyByName("C")
	s.AssertThat(keyC).IsNull()
	unknownKey := s.Registry.CreateTemplateType("C")

	templateTypeMap := templatizedInstance.GetTemplateTypeMap()
	s.AssertThat(templateTypeMap.HasTemplateKey(keyA)).IsTrue()
	s.AssertThat(templateTypeMap.HasTemplateKey(keyB)).IsTrue()
	s.AssertThat(templateTypeMap.HasTemplateKey(unknownKey)).IsFalse()

	s.AssertThat(templateTypeMap.GetResolvedTemplateType(keyA)).IsEqualTo(s.NUMBER_TYPE)
	s.AssertThat(templateTypeMap.GetResolvedTemplateType(keyB)).IsEqualTo(s.STRING_TYPE)
	s.AssertThat(templateTypeMap.GetResolvedTemplateType(unknownKey)).IsEqualTo(s.UNKNOWN_TYPE)

	s.AssertThat(templatizedInstance.String()).IsEqualTo("TestingType<number,string>")
}

func TestJSType_PartiallyTemplatizedType(t *testing.T) {
	s := newJsTypeTest(t)

	templatizedCtor :=
		withOpenRegistry(s, func() FunctionType {
			return s.Registry.CreateConstructorType(
				ptr("TestingType"),
				/* source= */ nil,
				nil,
				s.UNKNOWN_TYPE,
				[]*TemplateType{
					s.Registry.CreateTemplateType("A"), s.Registry.CreateTemplateType("B")},
				/* isAbstract= */ false)
		})
	templatizedInstance := s.Registry.CreateTemplatizedTypeOfList(
		templatizedCtor.GetInstanceType(),
		[]JSType{s.NUMBER_TYPE})

	ctrTypeMap := templatizedCtor.GetTemplateTypeMap()
	keyA := ctrTypeMap.GetLastTemplateTypeKeyByName("A")
	s.AssertThat(keyA).IsNotNull()
	keyB := ctrTypeMap.GetLastTemplateTypeKeyByName("B")
	s.AssertThat(keyB).IsNotNull()
	keyC := ctrTypeMap.GetLastTemplateTypeKeyByName("C")
	s.AssertThat(keyC).IsNull()
	unknownKey := s.Registry.CreateTemplateType("C")

	templateTypeMap := templatizedInstance.GetTemplateTypeMap()
	s.AssertThat(templateTypeMap.HasTemplateKey(keyA)).IsTrue()
	s.AssertThat(templateTypeMap.HasTemplateKey(keyB)).IsTrue()
	s.AssertThat(templateTypeMap.HasTemplateKey(unknownKey)).IsFalse()

	s.AssertThat(templateTypeMap.GetResolvedTemplateType(keyA)).IsEqualTo(s.NUMBER_TYPE)
	s.AssertThat(templateTypeMap.GetResolvedTemplateType(keyB)).IsEqualTo(s.UNKNOWN_TYPE)
	s.AssertThat(templateTypeMap.GetResolvedTemplateType(unknownKey)).IsEqualTo(s.UNKNOWN_TYPE)

	s.AssertThat(templatizedInstance.String()).IsEqualTo("TestingType<number,?>")
}

func TestJSType_TemplatizedTypeWithSubclass(t *testing.T) {
	s := newJsTypeTest(t)
	parentKeyA := s.Registry.CreateTemplateType("A")
	childKeyA := s.Registry.CreateTemplateType("A")
	childKeyB := s.Registry.CreateTemplateType("B")

	templatizedCtor :=
		withOpenRegistry(s, func() FunctionType {
			parent :=
				FunctionType_Builder_(s.Registry).
					WithName(ptr("ParentType")).
					ForConstructor().
					WithTemplateKeysOfList([]*TemplateType{parentKeyA}).
					Build()
			baseType :=
				s.Registry.CreateTemplatizedTypeOfList(
					parent.GetInstanceType(), []JSType{s.NULL_TYPE})
			// create a subclass that `@extends {ParentType<null>}`
			// and has a template key with the same name 'A' as the parent.
			child :=
				FunctionType_Builder_(s.Registry).
					WithPrototypeBasedOn(baseType).
					WithName(ptr("TestingType")).
					ForConstructor().
					WithTemplateKeysOfList([]*TemplateType{childKeyA, childKeyB}).
					Build()
			child.GetInstanceType().MergeSupertypeTemplateTypes(baseType)
			return child
		})

	templatizedInstance :=
		s.Registry.CreateTemplatizedTypeOfList(
			templatizedCtor.GetInstanceType(), []JSType{s.NUMBER_TYPE, s.STRING_TYPE})

	templateTypeMap := templatizedInstance.GetTemplateTypeMap()

	s.AssertThat(templateTypeMap.GetLastTemplateTypeKeyByName("A")).IsEqualTo(childKeyA)
	s.AssertThat(templateTypeMap.GetLastTemplateTypeKeyByName("B")).IsEqualTo(childKeyB)

	s.AssertThat(templateTypeMap.HasTemplateKey(parentKeyA)).IsTrue()
	s.AssertThat(templateTypeMap.HasTemplateKey(childKeyA)).IsTrue()
	s.AssertThat(templateTypeMap.HasTemplateKey(childKeyB)).IsTrue()

	s.AssertThat(templateTypeMap.GetResolvedTemplateType(parentKeyA)).IsEqualTo(s.NULL_TYPE)
	s.AssertThat(templateTypeMap.GetResolvedTemplateType(childKeyA)).IsEqualTo(s.NUMBER_TYPE)
	s.AssertThat(templateTypeMap.GetResolvedTemplateType(childKeyB)).IsEqualTo(s.STRING_TYPE)

	s.AssertThat(templatizedInstance.String()).IsEqualTo("TestingType<number,string>")
}

func TestJSType_TemplateTypeValidator(t *testing.T) {
	s := newJsTypeTest(t)
	// The template type setValidator() will see the TemplateType, not the referenced unknown type
	// like other ProxyObjectTypes do.
	t_ := NewTemplateType(s.Registry, "T")

	s.AssertThat(t_.SetValidator(func(type_ JSType) bool { return type_.IsTemplateType() })).IsTrue()
	s.AssertThat(t_.SetValidator(func(type_ JSType) bool { return !type_.IsTemplateType() })).IsFalse()
}

func TestJSType_TemplateTypeHasReferenceName(t *testing.T) {
	s := newJsTypeTest(t)
	t_ := NewTemplateType(s.Registry, "T")

	s.AssertType(t_).GetReferenceNameIsEqualTo("T")
}

func TestJSType_CanCastTo(t *testing.T) {
	s := newJsTypeTest(t)
	s.AssertThat(s.ALL_TYPE.CanCastTo(s.NULL_TYPE)).IsTrue()
	s.AssertThat(s.ALL_TYPE.CanCastTo(s.VOID_TYPE)).IsTrue()
	s.AssertThat(s.ALL_TYPE.CanCastTo(s.STRING_TYPE)).IsTrue()
	s.AssertThat(s.ALL_TYPE.CanCastTo(s.NUMBER_TYPE)).IsTrue()
	s.AssertThat(s.ALL_TYPE.CanCastTo(s.BOOLEAN_TYPE)).IsTrue()
	s.AssertThat(s.ALL_TYPE.CanCastTo(s.OBJECT_TYPE)).IsTrue()

	s.AssertThat(s.NUMBER_TYPE.CanCastTo(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.CanCastTo(s.VOID_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.CanCastTo(s.STRING_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.CanCastTo(s.NUMBER_TYPE)).IsTrue()
	s.AssertThat(s.NUMBER_TYPE.CanCastTo(s.BOOLEAN_TYPE)).IsFalse()
	s.AssertThat(s.NUMBER_TYPE.CanCastTo(s.OBJECT_TYPE)).IsFalse()

	s.AssertThat(s.STRING_TYPE.CanCastTo(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.STRING_TYPE.CanCastTo(s.VOID_TYPE)).IsFalse()
	s.AssertThat(s.STRING_TYPE.CanCastTo(s.STRING_TYPE)).IsTrue()
	s.AssertThat(s.STRING_TYPE.CanCastTo(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.STRING_TYPE.CanCastTo(s.BOOLEAN_TYPE)).IsFalse()
	s.AssertThat(s.STRING_TYPE.CanCastTo(s.OBJECT_TYPE)).IsFalse()

	s.AssertThat(s.BOOLEAN_TYPE.CanCastTo(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.CanCastTo(s.VOID_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.CanCastTo(s.STRING_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.CanCastTo(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.BOOLEAN_TYPE.CanCastTo(s.BOOLEAN_TYPE)).IsTrue()
	s.AssertThat(s.BOOLEAN_TYPE.CanCastTo(s.OBJECT_TYPE)).IsFalse()

	s.AssertThat(s.OBJECT_TYPE.CanCastTo(s.NULL_TYPE)).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.CanCastTo(s.VOID_TYPE)).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.CanCastTo(s.STRING_TYPE)).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.CanCastTo(s.NUMBER_TYPE)).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.CanCastTo(s.BOOLEAN_TYPE)).IsFalse()
	s.AssertThat(s.OBJECT_TYPE.CanCastTo(s.OBJECT_TYPE)).IsTrue()

	s.AssertThat(s.BOOLEAN_TYPE.CanCastTo(s.OBJECT_NUMBER_STRING)).IsFalse()
	s.AssertThat(s.OBJECT_NUMBER_STRING.CanCastTo(s.BOOLEAN_TYPE)).IsFalse()

	s.AssertThat(s.ARRAY_TYPE.CanCastTo(s.FUNCTION_TYPE)).IsFalse()
	s.AssertThat(s.FUNCTION_TYPE.CanCastTo(s.ARRAY_TYPE)).IsFalse()

	s.AssertThat(s.NULL_VOID.CanCastTo(s.ARRAY_TYPE)).IsFalse()
	s.AssertThat(s.NULL_VOID.CanCastTo(s.CreateUnionType(s.ARRAY_TYPE, s.NULL_TYPE))).IsTrue()

	// We currently allow any function to be cast to any other function type
	s.AssertThat(s.ARRAY_FUNCTION_TYPE.CanCastTo(s.BOOLEAN_OBJECT_FUNCTION_TYPE)).IsTrue()
}

func TestJSType_EqualityOfClassTypes_withSameReferenceName_preResolution(t *testing.T) {
	s := newJsTypeTest(t)
	closer := s.Registry.GetResolver().OpenForDefinition()
	defer closer.Close()
	{
		classACtor :=
			FunctionType_Builder_(s.Registry).ForConstructor().WithName(ptr("Foo")).Build()

		classBCtor :=
			FunctionType_Builder_(s.Registry).ForConstructor().WithName(ptr("Foo")).Build()

		// Currently, to handle NamedTypes, we treat unresolved type equality as purely based on
		// reference name.
		s.AssertType(classACtor.GetInstanceType()).IsEqualTo(classBCtor.GetInstanceType())
	}
}

func TestJSType_EqualityOfClassTypes_withSameReferenceName_differentGoogModuleId_preResolution(t *testing.T) {
	s := newJsTypeTest(t)
	closer := s.Registry.GetResolver().OpenForDefinition()
	defer closer.Close()
	{
		classACtor :=
			FunctionType_Builder_(s.Registry).
				ForConstructor().
				WithName(ptr("Foo")).
				SetGoogModuleId(ptr("module1")).
				Build()

		classBCtor :=
			FunctionType_Builder_(s.Registry).ForConstructor().WithName(ptr("Foo")).Build()

		// the different goog module id indicate the types are not equal.
		s.AssertType(classACtor.GetInstanceType()).IsNotEqualTo(classBCtor.GetInstanceType())
	}
}

func TestJSType_EqualityOfClassTypes_withSameReferenceName_sameGoogModuleId_preResolution(t *testing.T) {
	s := newJsTypeTest(t)
	closer := s.Registry.GetResolver().OpenForDefinition()
	defer closer.Close()
	{
		classACtor :=
			FunctionType_Builder_(s.Registry).
				ForConstructor().
				WithName(ptr("Foo")).
				SetGoogModuleId(ptr("module1")).
				Build()

		classBCtor :=
			FunctionType_Builder_(s.Registry).
				ForConstructor().
				WithName(ptr("Foo")).
				SetGoogModuleId(ptr("module1")).
				Build()

		// Currently, to handle NamedTypes, we treat unresolved type equality as purely based on
		// reference name & goog module id.
		s.AssertType(classACtor.GetInstanceType()).IsEqualTo(classBCtor.GetInstanceType())
	}
}

func TestJSType_EqualityOfClassTypes_withSameReferenceName_postResolution(t *testing.T) {
	s := newJsTypeTest(t)
	classACtor :=
		withOpenRegistry(s, func() FunctionType {
			return FunctionType_Builder_(s.Registry).ForConstructor().WithName(ptr("Foo")).Build()
		})
	classBCtor :=
		withOpenRegistry(s, func() FunctionType {
			return FunctionType_Builder_(s.Registry).ForConstructor().WithName(ptr("Foo")).Build()
		})

	s.AssertType(classACtor).IsNotEqualTo(classBCtor)
	s.AssertType(classACtor.GetInstanceType()).IsNotEqualTo(classBCtor.GetInstanceType())
	s.AssertType(classACtor.GetPrototype()).IsNotEqualTo(classBCtor.GetPrototype())
}

func TestJSType_EqualityOfInterfaceTypes_withSameReferenceName_postResolution(t *testing.T) {
	s := newJsTypeTest(t)
	interfaceACtor :=
		FunctionType_Builder_(s.Registry).ForInterface().WithName(ptr("Foo")).Build()

	interfaceBCtor :=
		FunctionType_Builder_(s.Registry).ForInterface().WithName(ptr("Foo")).Build()

	s.AssertType(interfaceACtor).IsNotEqualTo(interfaceBCtor)
	s.AssertType(interfaceACtor.GetInstanceType()).IsNotEqualTo(interfaceBCtor.GetInstanceType())
	s.AssertType(interfaceACtor.GetPrototype()).IsNotEqualTo(interfaceBCtor.GetPrototype())
}

func TestJSType_RecordTypeEquality(t *testing.T) {
	s := newJsTypeTest(t)
	// {x: number}
	firstType := s.Registry.CreateRecordType(mapOf("x", s.NUMBER_TYPE))
	secondType := s.Registry.CreateRecordType(mapOf("x", s.NUMBER_TYPE))

	s.AssertType(firstType).IsNotSameInstanceAs(secondType)
	s.AssertType(firstType).IsEqualTo(secondType)
	s.AssertType(firstType).IsEqualTo(secondType)
}

func TestJSType_RecordAndInterfaceObjectTypeNotEqualWithSameProperties(t *testing.T) {
	s := newJsTypeTest(t)
	// {x: number}
	firstType := s.Registry.CreateRecordType(mapOf("x", s.NUMBER_TYPE))
	// /** @interface */
	// function Foo() {}
	// /** @type {number} */
	// Foo.prototype.x;
	secondTypeConstructor :=
		FunctionType_Builder_(s.Registry).ForInterface().WithName(ptr("Foo")).Build()

	secondTypeConstructor.GetPrototype().DefineProperty("x", s.NUMBER_TYPE, false, nil)
	secondTypeConstructor.SetImplicitMatch(true)
	secondType := secondTypeConstructor.GetInstanceType()

	// These are not equal but are structurally equivalent
	s.AssertType(firstType).IsNotEqualTo(secondType)
}

func TestJSType_RecordAndObjectLiteralWithSameProperties(t *testing.T) {
	s := newJsTypeTest(t)
	// {x: number}
	firstType := s.Registry.CreateRecordType(mapOf("x", s.NUMBER_TYPE))
	s.AssertType(firstType).ToStringIsEqualTo("{x: number}")
	// the type inferred for `const obj = {x: 1};`
	secondType := s.Registry.CreateAnonymousObjectType(nil)
	secondType.DefineDeclaredProperty("x", s.NUMBER_TYPE, nil)
	s.AssertType(secondType).ToStringIsEqualTo("{x: number}")

	// These are neither equal nor structurally equivalent because the second type is not a
	// structural type.
	s.AssertType(firstType).IsNotEqualTo(secondType)

	// The second type is a subtype of the first type but not vice versa because only the first type
	// is structural.
	s.AssertType(firstType).IsNotSubtypeOf(secondType)
	s.AssertType(secondType).IsSubtypeOf(firstType)
}

func withOpenRegistry[T any](s *jsTypeTest, cb func() T) T {
	closer := s.Registry.GetResolver().OpenForDefinition()
	defer closer.Close()
	{
		return cb()
	}
}
