// Port of com/google/javascript/rhino/jstype/TemplateTypeTest.java @ closure-compiler v20261015

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

type templateTypeTest struct {
	*rhinotesting.BaseJSTypeTestCase
}

func newTemplateTypeTest(t *testing.T) *templateTypeTest {
	s := &templateTypeTest{BaseJSTypeTestCase: rhinotesting.NewBaseJSTypeTestCase(t)}
	return s
}

func TestTemplateType_TemplateTypes_boundedByUnknown_areNotEqual(t *testing.T) {
	s := newTemplateTypeTest(t)
	unknownType := s.Registry.GetNativeType(JSTypeNative_UNKNOWN_TYPE)
	t1 := NewTemplateTypeWithBound(s.Registry, "T", unknownType)
	t2 := NewTemplateTypeWithBound(s.Registry, "T", unknownType)

	t1.Resolve(s.Registry.GetErrorReporter())
	t2.Resolve(s.Registry.GetErrorReporter())

	s.AssertType(t1).IsNotEqualTo(t2)
	s.AssertType(t1).IsSubtypeOf(t2)
}

func TestTemplateType_TemplateTypes_boundedByNominalTypes_areNotEqual(t *testing.T) {
	s := newTemplateTypeTest(t)
	fooType :=
		s.Registry.CreateObjectType(ptr("Foo"), s.Registry.GetNativeObjectType(JSTypeNative_OBJECT_TYPE))
	t1 := NewTemplateTypeWithBound(s.Registry, "T", fooType)
	t2 := NewTemplateTypeWithBound(s.Registry, "T", fooType)

	t1.Resolve(s.Registry.GetErrorReporter())
	t2.Resolve(s.Registry.GetErrorReporter())

	s.AssertType(t1).IsNotEqualTo(t2)
	s.AssertType(t1).IsNotSubtypeOf(t2)
}

func TestTemplateType_ProxiesOf_templateTypes_boundedByNominalTypes_areNotEqual(t *testing.T) {
	s := newTemplateTypeTest(t)
	fooType :=
		s.Registry.CreateObjectType(ptr("Foo"), s.Registry.GetNativeObjectType(JSTypeNative_OBJECT_TYPE))

	t1 := NewTemplateTypeWithBound(s.Registry, "T", fooType)
	t2 := NewTemplateTypeWithBound(s.Registry, "T", fooType)

	t1Proxy :=
		NamedType_Builder_(s.Registry, "T").
			SetResolutionKind(NamedType_ResolutionKind_NONE).
			SetReferencedType(t1).
			Build()
	t2Proxy :=
		NamedType_Builder_(s.Registry, "T").
			SetResolutionKind(NamedType_ResolutionKind_NONE).
			SetReferencedType(t2).
			Build()

	s.AssertType(t1Proxy).IsNotEqualTo(t2Proxy)
	s.AssertType(t1Proxy).IsNotSubtypeOf(t2Proxy)
}

func TestTemplateType_TemplateTypes_areNotEqualToUnknown(t *testing.T) {
	s := newTemplateTypeTest(t)
	// Given
	t_ := NewTemplateType(s.Registry, "T")
	s.AssertThat(t_.IsUnknownType()).IsTrue()

	// Then
	s.AssertType(t_).IsNotEqualTo(s.UNKNOWN_TYPE)
	s.AssertType(t_).IsNotEqualTo(s.CHECKED_UNKNOWN_TYPE)
}

func TestTemplateType_TemplateTypes_withBound_areNotEqualToBound(t *testing.T) {
	s := newTemplateTypeTest(t)
	// Given
	t_ := NewTemplateTypeWithBound(s.Registry, "T", s.NUMBER_TYPE)

	// Then
	s.AssertType(t_).IsNotEqualTo(s.NUMBER_TYPE)
}
