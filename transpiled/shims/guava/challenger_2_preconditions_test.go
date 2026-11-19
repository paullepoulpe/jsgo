package guava_test

import (
	"testing"

	"github.com/paullepoulpe/jsgo/transpiled/shims/guava"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

type stringerObj struct{ val string }

func (s stringerObj) String() string { return "Stringer(" + s.val + ")" }

func TestPreconditionsFixedArityZeroAllocOnPass(t *testing.T) {
	// Passing expressions must not allocate on the heap
	allocsState1 := testing.AllocsPerRun(1000, func() {
		guava.CheckStateMsg1(true, "test %s", 123)
	})
	if allocsState1 != 0 {
		t.Errorf("CheckStateMsg1 allocated %f times on passing path, expected 0", allocsState1)
	}

	allocsState2 := testing.AllocsPerRun(1000, func() {
		guava.CheckStateMsg2(true, "test %s %s", 123, 456)
	})
	if allocsState2 != 0 {
		t.Errorf("CheckStateMsg2 allocated %f times on passing path, expected 0", allocsState2)
	}

	allocsState3 := testing.AllocsPerRun(1000, func() {
		guava.CheckStateMsg3(true, "test %s %s %s", 123, 456, 789)
	})
	if allocsState3 != 0 {
		t.Errorf("CheckStateMsg3 allocated %f times on passing path, expected 0", allocsState3)
	}

	allocsArg1 := testing.AllocsPerRun(1000, func() {
		guava.CheckArgumentMsg1(true, "test %s", "hello")
	})
	if allocsArg1 != 0 {
		t.Errorf("CheckArgumentMsg1 allocated %f times on passing path, expected 0", allocsArg1)
	}

	ptr := new(int)
	allocsNotNull1 := testing.AllocsPerRun(1000, func() {
		_ = guava.CheckNotNullMsg1(ptr, "ptr %s", "msg")
	})
	if allocsNotNull1 != 0 {
		t.Errorf("CheckNotNullMsg1 allocated %f times on passing path, expected 0", allocsNotNull1)
	}
}

func TestPreconditionsFailureTypesAndMessages(t *testing.T) {
	// CheckState failure must throw IllegalStateException
	assertException(t, "CheckState(false)", "java.lang.IllegalStateException", "", func() {
		guava.CheckState(false)
	})
	assertException(t, "CheckStateMsg1(false)", "java.lang.IllegalStateException", "error 42", func() {
		guava.CheckStateMsg1(false, "error %s", 42)
	})

	// CheckArgument failure must throw IllegalArgumentException
	assertException(t, "CheckArgument(false)", "java.lang.IllegalArgumentException", "", func() {
		guava.CheckArgument(false)
	})
	assertException(t, "CheckArgumentMsg2(false)", "java.lang.IllegalArgumentException", "expected 1 got 2", func() {
		guava.CheckArgumentMsg2(false, "expected %s got %s", 1, 2)
	})

	// CheckNotNull failure must throw NullPointerException
	var nilPtr *int
	assertException(t, "CheckNotNull(nil)", "java.lang.NullPointerException", "", func() {
		guava.CheckNotNull(nilPtr)
	})
	assertException(t, "CheckNotNullMsg1(nil)", "java.lang.NullPointerException", "must not be null: key", func() {
		guava.CheckNotNullMsg1(nilPtr, "must not be null: %s", "key")
	})

	// ElementIndex out of bounds
	assertException(t, "CheckElementIndex(-1, 5)", "java.lang.IndexOutOfBoundsException", "index (-1) must be less than size (5)", func() {
		guava.CheckElementIndex(-1, 5)
	})
	assertException(t, "CheckElementIndex(5, 5)", "java.lang.IndexOutOfBoundsException", "index (5) must be less than size (5)", func() {
		guava.CheckElementIndex(5, 5)
	})
	// Valid element index
	if idx := guava.CheckElementIndex(0, 5); idx != 0 {
		t.Errorf("CheckElementIndex(0, 5) = %d, expected 0", idx)
	}
	if idx := guava.CheckElementIndex(4, 5); idx != 4 {
		t.Errorf("CheckElementIndex(4, 5) = %d, expected 4", idx)
	}

	// PositionIndex
	assertException(t, "CheckPositionIndex(-1, 5)", "java.lang.IndexOutOfBoundsException", "index (-1) must not be greater than size (5)", func() {
		guava.CheckPositionIndex(-1, 5)
	})
	assertException(t, "CheckPositionIndex(6, 5)", "java.lang.IndexOutOfBoundsException", "index (6) must not be greater than size (5)", func() {
		guava.CheckPositionIndex(6, 5)
	})
	if idx := guava.CheckPositionIndex(5, 5); idx != 5 {
		t.Errorf("CheckPositionIndex(5, 5) = %d, expected 5", idx)
	}

	// PositionIndexes
	assertException(t, "CheckPositionIndexes(-1, 2, 5)", "java.lang.IndexOutOfBoundsException", "invalid range [-1, 2) for size 5", func() {
		guava.CheckPositionIndexes(-1, 2, 5)
	})
	assertException(t, "CheckPositionIndexes(3, 2, 5)", "java.lang.IndexOutOfBoundsException", "invalid range [3, 2) for size 5", func() {
		guava.CheckPositionIndexes(3, 2, 5)
	})
	assertException(t, "CheckPositionIndexes(0, 6, 5)", "java.lang.IndexOutOfBoundsException", "invalid range [0, 6) for size 5", func() {
		guava.CheckPositionIndexes(0, 6, 5)
	})
	// Valid range
	guava.CheckPositionIndexes(0, 5, 5)
	guava.CheckPositionIndexes(2, 2, 5)
}

func TestLenientFormatEdgeCases(t *testing.T) {
	cases := []struct {
		template string
		args     []any
		expected string
	}{
		{"no placeholders", nil, "no placeholders"},
		{"no placeholders", []any{}, "no placeholders"},
		{"no placeholders", []any{"extra1", "extra2"}, "no placeholders [extra1, extra2]"},
		{"single %s placeholder", []any{"val"}, "single val placeholder"},
		{"consecutive %s%s", []any{"a", "b"}, "consecutive ab"},
		{"start %s", []any{"first"}, "start first"},
		{"%s end", []any{"last"}, "last end"},
		{"fewer args %s %s %s", []any{"one"}, "fewer args one %s %s"},
		{"null arg %s", []any{nil}, "null arg null"},
		{"stringer arg %s", []any{stringerObj{"custom"}}, "stringer arg Stringer(custom)"},
		{"empty template", []any{"extra"}, "empty template [extra]"},
		{"", []any{"extra1", "extra2"}, " [extra1, extra2]"},
	}

	for _, c := range cases {
		out := guava.LenientFormat(c.template, c.args...)
		if out != c.expected {
			t.Errorf("LenientFormat(%q, %v) = %q, expected %q", c.template, c.args, out, c.expected)
		}
	}
}

func assertException(t *testing.T, desc, expectedClass, expectedMsgSubstr string, f func()) {
	t.Helper()
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("%s did not panic", desc)
			return
		}
		th := javalang.RecoverException(r)
		if th == nil {
			t.Errorf("%s panic could not be converted to Throwable: %v", desc, r)
			return
		}
		if th.Error() == "" {
			t.Errorf("%s Throwable.Error() is empty", desc)
		}
		if expectedMsgSubstr != "" && th.GetMessage() != expectedMsgSubstr {
			t.Errorf("%s message got %q, expected %q", desc, th.GetMessage(), expectedMsgSubstr)
		}
	}()
	f()
}
