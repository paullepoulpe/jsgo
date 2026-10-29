package jutil

import "testing"

func TestGoRuntimeFaultsMapToJavaExceptions(t *testing.T) {
	var nilPtr *int
	cases := []struct {
		name, class, msg string
		body             func()
	}{
		{"nil deref", "java.lang.NullPointerException", "", func() { _ = *nilPtr }},
		{"index", "java.lang.ArrayIndexOutOfBoundsException", "Index 5 out of bounds for length 3", func() {
			s := []int{1, 2, 3}
			i := 5
			_ = s[i]
		}},
		{"divide", "java.lang.ArithmeticException", "/ by zero", func() {
			z := 0
			_ = 1 / z
		}},
		{"cast", "java.lang.ClassCastException", "", func() {
			var x any = "s"
			_ = x.(int)
		}},
	}
	for _, c := range cases {
		var got Throwable
		Try(c.body, func(e IRuntimeException) { got = e })
		if got == nil {
			t.Fatalf("%s: not caught as RuntimeException", c.name)
		}
		if got.throwable().className != c.class {
			t.Errorf("%s: class %s, want %s", c.name, got.throwable().className, c.class)
		}
		if c.msg != "" && got.GetMessage() != c.msg {
			t.Errorf("%s: message %q, want %q", c.name, got.GetMessage(), c.msg)
		}
	}
	// A nil dereference is catchable as NullPointerException specifically.
	caught := false
	Try(func() { _ = *nilPtr }, func(e INullPointerException) { caught = true })
	if !caught {
		t.Fatal("nil dereference not caught as NullPointerException")
	}
}
