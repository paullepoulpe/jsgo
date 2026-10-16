package jscomp

import (
	"slices"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// assertThrows is org.junit.Assert.assertThrows: body must throw an
// exception assignable to E.
func assertThrows[E any](t *testing.T, body func()) E {
	t.Helper()
	var caught E
	ok := false
	jutil.Try(body, func(e E) {
		caught = e
		ok = true
	})
	if !ok {
		t.Fatalf("expected %T to be thrown", (*E)(nil))
	}
	return caught
}

// assertEqual is assertThat(got).isEqualTo(want) for comparable values.
func assertEqual[T comparable](t *testing.T, what string, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %v, want %v", what, got, want)
	}
}

// assertSliceEqual is containsExactly(...).inOrder() for slices.
func assertSliceEqual[T comparable](t *testing.T, what string, got, want []T) {
	t.Helper()
	if !slices.Equal(got, want) {
		t.Errorf("%s: got %v, want %v", what, got, want)
	}
}

// assertSameElements is containsExactly(...) (any order, with
// multiplicity).
func assertSameElements[T comparable](t *testing.T, what string, got, want []T) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("%s: got %v, want %v (any order)", what, got, want)
		return
	}
	counts := map[T]int{}
	for _, x := range want {
		counts[x]++
	}
	for _, x := range got {
		counts[x]--
	}
	for _, c := range counts {
		if c != 0 {
			t.Errorf("%s: got %v, want %v (any order)", what, got, want)
			return
		}
	}
}

// strPtrValue renders a nullable string for messages.
func strPtrValue(s *string) string {
	if s == nil {
		return "<null>"
	}
	return *s
}
