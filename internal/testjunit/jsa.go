package testjunit

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// JsaEq is assertThat(got).isEqualTo(want) for comparable values.
func JsaEq[T comparable](t *testing.T, got, want T, ctx ...any) {
	t.Helper()
	if got != want {
		t.Errorf("%sgot %#v, want %#v", jsaCtx(ctx), got, want)
	}
}

// JsaDeepEq is assertThat(got).isEqualTo(want) for slices and maps.
func JsaDeepEq(t *testing.T, got, want any, ctx ...any) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%sgot %#v, want %#v", jsaCtx(ctx), got, want)
	}
}

// JsaJavaEq compares with the values' Java equals().
func JsaJavaEq(t *testing.T, got, want any, ctx ...any) {
	t.Helper()
	if !jutil.ObjectsEquals(got, want) {
		t.Errorf("%sgot %v, want %v", jsaCtx(ctx), got, want)
	}
}

func JsaTrue(t *testing.T, cond bool, ctx ...any) {
	t.Helper()
	if !cond {
		t.Errorf("%sexpected true", jsaCtx(ctx))
	}
}

func JsaFalse(t *testing.T, cond bool, ctx ...any) {
	t.Helper()
	if cond {
		t.Errorf("%sexpected false", jsaCtx(ctx))
	}
}

// JsaThrows is assertThrows(E.class, f): it returns the exception.
func JsaThrows[E any](t *testing.T, f func()) (caught E) {
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
