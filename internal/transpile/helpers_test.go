package transpile

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/testjunit"
)

func jsaEq[T comparable](t *testing.T, got, want T, ctx ...any) {
	t.Helper()
	testjunit.JsaEq(t, got, want, ctx...)
}

func jsaDeepEq(t *testing.T, got, want any, ctx ...any) {
	t.Helper()
	testjunit.JsaDeepEq(t, got, want, ctx...)
}

func jsaJavaEq(t *testing.T, got, want any, ctx ...any) {
	t.Helper()
	testjunit.JsaJavaEq(t, got, want, ctx...)
}

func jsaTrue(t *testing.T, cond bool, ctx ...any) {
	t.Helper()
	testjunit.JsaTrue(t, cond, ctx...)
}

func jsaFalse(t *testing.T, cond bool, ctx ...any) {
	t.Helper()
	testjunit.JsaFalse(t, cond, ctx...)
}

func jsaThrows[E any](t *testing.T, f func()) E {
	t.Helper()
	return testjunit.JsaThrows[E](t, f)
}
