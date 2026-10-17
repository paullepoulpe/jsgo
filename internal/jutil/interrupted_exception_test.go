package jutil

import "testing"

// Go-only: test of java.lang.InterruptedException's toString()/getMessage()
// and of IllegalStateException(String, Throwable), as printed by JDK 21.
func TestInterruptedExceptionAndIllegalStateWithCause(t *testing.T) {
	e := NewInterruptedException()
	if got := e.Error(); got != "java.lang.InterruptedException" {
		t.Errorf("InterruptedException().toString() = %q", got)
	}
	wrapped := NewRuntimeExceptionWithCause(e)
	if got := wrapped.Error(); got != "java.lang.RuntimeException: java.lang.InterruptedException" {
		t.Errorf("new RuntimeException(new InterruptedException()).toString() = %q", got)
	}
	var caught bool
	Try(func() { panic(e) }, func(IException) { caught = true })
	if !caught {
		t.Errorf("InterruptedException is not caught as a checked exception")
	}
	ise := NewIllegalStateExceptionWithMessageAndCause("msg", NewRuntimeException("cause"))
	if got := ise.Error(); got != "java.lang.IllegalStateException: msg" {
		t.Errorf("IllegalStateException(msg, cause).toString() = %q", got)
	}
	if ise.GetCause() == nil || ise.GetCause().GetMessage() != "cause" {
		t.Errorf("IllegalStateException(msg, cause).getCause() wrong")
	}
	var isIse bool
	Try(func() { panic(ise) }, func(IIllegalStateException) { isIse = true })
	if !isIse {
		t.Errorf("not caught as IllegalStateException")
	}
}
