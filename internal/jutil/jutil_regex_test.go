package jutil

import (
	"fmt"
	"strings"
	"testing"
)

// Expected values below were produced by the JDK 21 runtime that the pinned
// Closure Compiler oracle runs on.

func throwableString(f func()) (s string) {
	s = "no exception"
	Try(f, func(e Throwable) { s = e.Error() })
	return s
}

func TestBoundsExceptionMessagesMatchJava(t *testing.T) {
	cases := []struct {
		f    func()
		want string
	}{
		{func() { StringCharAt("abc", 3) }, "java.lang.StringIndexOutOfBoundsException: Index 3 out of bounds for length 3"},
		{func() { StringCharAt("abc", -1) }, "java.lang.StringIndexOutOfBoundsException: Index -1 out of bounds for length 3"},
		{func() { StringCharAt("abé", 3) }, "java.lang.StringIndexOutOfBoundsException: Index 3 out of bounds for length 3"},
		{func() { StringSubstring("abc", 2, 1) }, "java.lang.StringIndexOutOfBoundsException: Range [2, 1) out of bounds for length 3"},
		{func() { StringSubstringFrom("abc", 4) }, "java.lang.StringIndexOutOfBoundsException: Range [4, 3) out of bounds for length 3"},
		{func() { StringSubstring("abc", -1, 2) }, "java.lang.StringIndexOutOfBoundsException: Range [-1, 2) out of bounds for length 3"},
		{func() { StringSubstring("abc", 1, 5) }, "java.lang.StringIndexOutOfBoundsException: Range [1, 5) out of bounds for length 3"},
		{func() { NewStringBuilder().AppendCharSequenceRange(NewStringBuilderFrom("xyz"), 2, 5) }, "java.lang.IndexOutOfBoundsException: Range [2, 5) out of bounds for length 3"},
		{func() { CharacterToChars(0x110000) }, "java.lang.IllegalArgumentException: Not a valid Unicode code point: 0x110000"},
		{func() { CharacterToChars(-1) }, "java.lang.IllegalArgumentException: Not a valid Unicode code point: 0xFFFFFFFF"},
		{func() { panic(NewIndexOutOfBoundsExceptionNoMessage()) }, "java.lang.IndexOutOfBoundsException"},
	}
	for i, c := range cases {
		if got := throwableString(c.f); got != c.want {
			t.Errorf("case %d: got %q, want %q", i, got, c.want)
		}
	}
	// The new exceptions are caught as IndexOutOfBoundsException.
	caught := false
	Try(func() { StringCharAt("", 0) }, func(e IIndexOutOfBoundsException) { caught = true })
	if !caught {
		t.Error("StringIndexOutOfBoundsException is not an IndexOutOfBoundsException")
	}
}

func TestArraysBinarySearchIntsMatchesJava(t *testing.T) {
	arrs := [][]int{{}, {5}, {1, 3, 5, 7}, {1, 1, 1, 1, 1}, {0, 2, 2, 2, 4, 6, 6, 9}, {-5, -1, 0, 0x10000}}
	keys := []int{-6, -5, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 0x10000, 0x10001}
	want := []string{
		"-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,",
		"-1,-1,-1,-1,-1,-1,-1,-1,0,-2,-2,-2,-2,-2,-2,-2,",
		"-1,-1,-1,-1,0,-2,1,-3,2,-4,3,-5,-5,-5,-5,-5,",
		"-1,-1,-1,-1,2,-6,-6,-6,-6,-6,-6,-6,-6,-6,-6,-6,",
		"-1,-1,-1,0,-2,3,-5,4,-6,5,-8,-8,7,-9,-9,-9,",
		"-1,0,1,2,-4,-4,-4,-4,-4,-4,-4,-4,-4,-4,3,-5,",
	}
	for i, arr := range arrs {
		var sb strings.Builder
		for _, k := range keys {
			fmt.Fprintf(&sb, "%d,", ArraysBinarySearchInts(arr, k))
		}
		if sb.String() != want[i] {
			t.Errorf("binarySearch(%v) = %s, Java %s", arr, sb.String(), want[i])
		}
	}
}

func TestCharacterToChars(t *testing.T) {
	cases := map[int][]uint16{
		0:        {0},
		0x41:     {0x41},
		0xD800:   {0xD800},
		0xFFFF:   {0xFFFF},
		0x10000:  {0xD800, 0xDC00},
		0x1F600:  {0xD83D, 0xDE00},
		0x10FFFF: {0xDBFF, 0xDFFF},
	}
	for cp, want := range cases {
		got := CharacterToChars(cp)
		if !equalU16(got, want) {
			t.Errorf("toChars(%#x) = %x, want %x", cp, got, want)
		}
	}
}

func TestAppendCharSequenceRangeOfItself(t *testing.T) {
	// Java: sb.append(sb, start, end), as RegExpTree.Repetition does.
	sb := NewStringBuilderFrom("x(?:ab)")
	sb.AppendCharSequenceRange(sb, 1, sb.Length())
	if got := sb.String(); got != "x(?:ab)(?:ab)" {
		t.Errorf("got %q", got)
	}
}
