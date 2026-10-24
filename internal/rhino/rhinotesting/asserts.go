// Port of com/google/javascript/rhino/testing/Asserts.java @ closure-compiler v20261021

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
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhinotesting

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Helper methods for making assertions about the validity of types.

// Java: Asserts.assertResolvesToSame(JSType)
func Asserts_AssertResolvesToSame(t testing.TB, type_ rhino.JSType) rhino.JSType {
	t.Helper()
	Truth{T: t}.AssertThat(Asserts_AssertValidResolve(t, type_)).IsSameInstanceAs(type_)
	return type_
}

// Returns the resolved type.
//
// Java: Asserts.assertValidResolve(JSType)
func Asserts_AssertValidResolve(t testing.TB, type_ rhino.JSType) rhino.JSType {
	t.Helper()
	reporter := NewTestErrorReporter()
	resolvedType := type_.Resolve(reporter)
	Truth{T: t}.AssertWithMessage("JSType#resolve should not affect object equality").
		AboutTypes().
		That(resolvedType).
		IsEqualTo(type_)
	reporter.VerifyHasEncounteredAllWarningsAndErrors(t)
	return resolvedType
}

// Java: Asserts.assertTypeCollectionEquals(Iterable<T>, Iterable<S>)
func Asserts_AssertTypeCollectionEquals[T rhino.JSType, S rhino.JSType](t testing.TB, a []T, b []S) {
	t.Helper()
	tr := Truth{T: t}
	tr.AssertThat(b).HasSize(len(a))
	for i := range a {
		tr.AssertType(b[i]).IsEqualTo(a[i])
	}
}

// For the given equivalent types, run all type operations that should have
// trivially predictable results. Proves that the two types are equivalent.
//
// Java: Asserts.assertEquivalenceOperations(JSType, JSType)
func Asserts_AssertEquivalenceOperations(t testing.TB, a rhino.JSType, b rhino.JSType) {
	t.Helper()
	tr := Truth{T: t}
	tr.AssertType(a).IsEqualTo(a)
	tr.AssertType(b).IsEqualTo(a)
	tr.AssertType(b).IsEqualTo(b)

	tr.AssertThat(a.IsSubtypeOf(b)).IsTrue()
	tr.AssertThat(a.IsSubtypeOf(a)).IsTrue()
	tr.AssertThat(b.IsSubtypeOf(b)).IsTrue()
	tr.AssertThat(b.IsSubtypeOf(a)).IsTrue()

	tr.AssertType(a.GetGreatestSubtype(b)).IsEqualTo(a)
	tr.AssertType(a.GetGreatestSubtype(a)).IsEqualTo(a)
	tr.AssertType(b.GetGreatestSubtype(b)).IsEqualTo(a)
	tr.AssertType(b.GetGreatestSubtype(a)).IsEqualTo(a)

	tr.AssertType(a.GetLeastSupertype(b)).IsEqualTo(a)
	tr.AssertType(a.GetLeastSupertype(a)).IsEqualTo(a)
	tr.AssertType(b.GetLeastSupertype(b)).IsEqualTo(a)
	tr.AssertType(b.GetLeastSupertype(a)).IsEqualTo(a)

	tr.AssertThat(a.CanCastTo(b)).IsTrue()
	tr.AssertThat(a.CanCastTo(a)).IsTrue()
	tr.AssertThat(b.CanCastTo(b)).IsTrue()
	tr.AssertThat(b.CanCastTo(a)).IsTrue()
}
