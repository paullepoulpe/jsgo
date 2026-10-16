package rhino

import (
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
