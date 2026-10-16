package jscomp

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Test helpers shared by the tests ported in wave jscomp-a (prefix jsa to
// stay clear of the helpers other waves add to this package).

// jsaEq is assertThat(got).isEqualTo(want) for comparable values.
func jsaEq[T comparable](t *testing.T, got, want T, ctx ...any) {
	t.Helper()
	if got != want {
		t.Errorf("%sgot %#v, want %#v", jsaCtx(ctx), got, want)
	}
}

// jsaDeepEq is assertThat(got).isEqualTo(want) for slices and maps.
func jsaDeepEq(t *testing.T, got, want any, ctx ...any) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%sgot %#v, want %#v", jsaCtx(ctx), got, want)
	}
}

// jsaJavaEq compares with the values' Java equals().
func jsaJavaEq(t *testing.T, got, want any, ctx ...any) {
	t.Helper()
	if !jutil.ObjectsEquals(got, want) {
		t.Errorf("%sgot %v, want %v", jsaCtx(ctx), got, want)
	}
}

func jsaTrue(t *testing.T, cond bool, ctx ...any) {
	t.Helper()
	if !cond {
		t.Errorf("%sexpected true", jsaCtx(ctx))
	}
}

func jsaFalse(t *testing.T, cond bool, ctx ...any) {
	t.Helper()
	if cond {
		t.Errorf("%sexpected false", jsaCtx(ctx))
	}
}

// jsaThrows is assertThrows(E.class, f): it returns the exception.
func jsaThrows[E any](t *testing.T, f func()) (caught E) {
	t.Helper()
	thrown := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				thrown = true
				if e, ok := r.(E); ok {
					caught = e
					return
				}
				if e, ok := any(jutil.AsThrowable(r)).(E); ok {
					caught = e
					return
				}
				t.Errorf("unexpected panic %T: %v", r, r)
			}
		}()
		f()
	}()
	if !thrown {
		t.Errorf("expected %T to be thrown", caught)
	}
	return caught
}

func jsaCtx(ctx []any) string {
	if len(ctx) == 0 {
		return ""
	}
	return fmt.Sprint(ctx...) + ": "
}
