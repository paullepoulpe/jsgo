// Port of com/google/javascript/rhino/testing/Asserts.java @ closure-compiler v20260415

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
