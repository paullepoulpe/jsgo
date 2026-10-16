package jutil

import (
	"os"
	"strings"
	"testing"
)

// TestThrowableMessageMatchesJava replays testdata/throwable_message.tsv
// (tools/oracle/java/JutilThrowableVectors.java): getMessage() (or null) and
// toString() of the exceptions GetMessageOrNull and InterruptedException
// model.
func TestThrowableMessageMatchesJava(t *testing.T) {
	data, err := os.ReadFile("testdata/throwable_message.tsv")
	if err != nil {
		t.Fatal(err)
	}
	cases := map[string]Throwable{
		"InterruptedException()":                   NewInterruptedException(),
		"RuntimeException(InterruptedException())": NewRuntimeExceptionWithCause(NewInterruptedException()),
		`RuntimeException("")`:                     NewRuntimeException(""),
		`RuntimeException("msg")`:                  NewRuntimeException("msg"),
		"NullPointerException()":                   NewNullPointerException(),
		"IllegalStateException()":                  NewIllegalStateException(),
		`IllegalStateException("x")`:               NewIllegalStateExceptionMsg("x"),
	}
	n := 0
	for _, line := range strings.Split(strings.TrimSuffix(string(data), "\n"), "\n") {
		f := strings.Split(line, "\t")
		e, ok := cases[f[0]]
		if !ok {
			t.Fatalf("no Go case for %q", f[0])
		}
		n++
		got := "<null>"
		if m := GetMessageOrNull(e); m != nil {
			got = "=" + *m
		}
		if got != f[1] {
			t.Errorf("%s: getMessage() = %q, Java %q", f[0], got, f[1])
		}
		if e.Error() != f[2] {
			t.Errorf("%s: toString() = %q, Java %q", f[0], e.Error(), f[2])
		}
	}
	if n != len(cases) {
		t.Errorf("compared %d of %d cases", n, len(cases))
	}
	if _, ok := any(NewInterruptedException()).(IException); !ok {
		t.Error("InterruptedException is not a checked exception")
	}
}
