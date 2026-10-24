// Port of com/google/javascript/rhino/jstype/PrototypeObjectTypeTest.java @ closure-compiler v20261021

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
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/rhinotesting"
)

type prototypeObjectTypeTest struct {
	*rhinotesting.BaseJSTypeTestCase
}

func newPrototypeObjectTypeTest(t *testing.T) *prototypeObjectTypeTest {
	s := &prototypeObjectTypeTest{BaseJSTypeTestCase: rhinotesting.NewBaseJSTypeTestCase(t)}
	return s
}

func TestPrototypeObjectType_ToString(t *testing.T) {
	s := newPrototypeObjectTypeTest(t)
	type_ := s.Registry.CreateAnonymousObjectType(nil)
	s.AssertThat(type_.String()).IsEqualTo("{}")

	type_.DefineDeclaredProperty("foo", s.NUMBER_TYPE, rhino.NilNode)
	s.AssertThat(type_.String()).IsEqualTo("{foo: number}")

	type_.DefineDeclaredProperty("bar", type_, rhino.NilNode)
	s.AssertThat(type_.String()).
		IsEqualTo(
			strings.Join([]string{
				"{", //
				"  bar: {...},",
				"  foo: number",
				"}"}, "\n"))
}
