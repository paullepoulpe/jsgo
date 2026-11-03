package base

import "testing"

func TestOptString(t *testing.T) {
	var zero OptString
	if zero != NoString() {
		t.Fatalf("zero value = %#v, want NoString()", zero)
	}
	if zero.IsPresent() {
		t.Fatal("NoString().IsPresent() = true, want false")
	}
	if got, ok := zero.Get(); ok || got != "" {
		t.Fatalf("NoString().Get() = (%q, %v), want (\"\", false)", got, ok)
	}
	if got := zero.ValueOrEmpty(); got != "" {
		t.Fatalf("NoString().ValueOrEmpty() = %q, want \"\"", got)
	}
	if got := zero.ValueOr("fallback"); got != "fallback" {
		t.Fatalf("NoString().ValueOr(\"fallback\") = %q, want \"fallback\"", got)
	}
	if got := zero.Ptr(); got != nil {
		t.Fatalf("NoString().Ptr() = %v, want nil", got)
	}
	if got := zero.String(); got != "null" {
		t.Fatalf("NoString().String() = %q, want \"null\"", got)
	}
	if got := OptStringFromPtr(nil); got != NoString() {
		t.Fatalf("OptStringFromPtr(nil) = %#v, want NoString()", got)
	}

	empty := SomeString("")
	if empty == NoString() {
		t.Fatal("SomeString(\"\") == NoString(), want distinct values")
	}
	if !empty.IsPresent() {
		t.Fatal("SomeString(\"\").IsPresent() = false, want true")
	}
	if got, ok := empty.Get(); !ok || got != "" {
		t.Fatalf("SomeString(\"\").Get() = (%q, %v), want (\"\", true)", got, ok)
	}
	if got := empty.ValueOrEmpty(); got != "" {
		t.Fatalf("SomeString(\"\").ValueOrEmpty() = %q, want \"\"", got)
	}
	if got := empty.ValueOr("fallback"); got != "" {
		t.Fatalf("SomeString(\"\").ValueOr(\"fallback\") = %q, want \"\"", got)
	}
	if p := empty.Ptr(); p == nil || *p != "" {
		t.Fatalf("SomeString(\"\").Ptr() = %v, want non-nil \"\"", p)
	}
	if got := empty.String(); got != "" {
		t.Fatalf("SomeString(\"\").String() = %q, want \"\"", got)
	}

	raw := "hello"
	some := OptStringFromPtr(&raw)
	if some != SomeString("hello") {
		t.Fatalf("OptStringFromPtr(&%q) = %#v, want SomeString(\"hello\")", raw, some)
	}
	raw = "mutated"
	if got, ok := some.Get(); !ok || got != "hello" {
		t.Fatalf("after mutating input string, some.Get() = (%q, %v), want (\"hello\", true)", got, ok)
	}
	if got := some.ValueOrEmpty(); got != "hello" {
		t.Fatalf("some.ValueOrEmpty() = %q, want \"hello\"", got)
	}
	if got := some.ValueOr("fallback"); got != "hello" {
		t.Fatalf("some.ValueOr(\"fallback\") = %q, want \"hello\"", got)
	}
	if p := some.Ptr(); p == nil || *p != "hello" {
		t.Fatalf("some.Ptr() = %v, want non-nil \"hello\"", p)
	}
	if got := some.String(); got != "hello" {
		t.Fatalf("some.String() = %q, want \"hello\"", got)
	}
}
