package jutil

import "testing"

func TestHelpfulNullPointerException(t *testing.T) {
	const detail = `Cannot invoke "Object.toString()" because "value" is null`
	e := NewHelpfulNullPointerException(detail)
	if got, ok := ThrowableMessage(e); !ok || got != detail {
		t.Errorf("message (%q, %v), want (%q, true) (on by default, as on JDK 21)", got, ok, detail)
	}
	if got, want := e.Error(), "java.lang.NullPointerException: "+detail; got != want {
		t.Errorf("toString %q, want %q", got, want)
	}

	if !SetShowCodeDetailsInExceptionMessages(false) {
		t.Error("the option was off before")
	}
	defer SetShowCodeDetailsInExceptionMessages(true)
	e = NewHelpfulNullPointerException(detail)
	if got, ok := ThrowableMessage(e); ok {
		t.Errorf("with -XX:-ShowCodeDetailsInExceptionMessages: message %q, want null", got)
	}
	if got, want := e.Error(), "java.lang.NullPointerException"; got != want {
		t.Errorf("toString %q, want %q", got, want)
	}
	if _, ok := any(e).(INullPointerException); !ok {
		t.Error("not a NullPointerException")
	}
}
