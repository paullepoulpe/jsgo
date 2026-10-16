package jutil

import "testing"

// new Error("Unknown op X").toString() is "java.lang.Error: Unknown op X" on
// the JDK; an Error is not a RuntimeException.
func TestNewJavaLangError(t *testing.T) {
	e := NewJavaLangError("Unknown op X")
	if got, want := e.Error(), "java.lang.Error: Unknown op X"; got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
	if got := e.GetMessage(); got != "Unknown op X" {
		t.Errorf("GetMessage() = %q", got)
	}
	caughtAsRuntime, caughtAsError := false, false
	Try(func() {
		Try(func() { panic(NewJavaLangError("x")) }, func(e IRuntimeException) { caughtAsRuntime = true })
	}, func(e IError) { caughtAsError = true })
	if caughtAsRuntime || !caughtAsError {
		t.Errorf("caughtAsRuntime=%v caughtAsError=%v", caughtAsRuntime, caughtAsError)
	}
}
