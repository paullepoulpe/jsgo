package javalang_test

import (
	"strings"
	"testing"
	"unicode/utf16"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

func TestBMPCharacters(t *testing.T) {
	s := "Hello, 世界! \u03A9"
	expectedLen := 12

	if l := javalang.StringLength(s); l != expectedLen {
		t.Errorf("StringLength(%q) = %d, expected %d", s, l, expectedLen)
	}

	expectedUnits := []uint16{
		'H', 'e', 'l', 'l', 'o', ',', ' ',
		0x4E16, 0x754C, '!', ' ', 0x03A9,
	}

	u16 := javalang.StringToUTF16(s)
	if len(u16) != len(expectedUnits) {
		t.Fatalf("StringToUTF16 length mismatch: got %d, expected %d", len(u16), len(expectedUnits))
	}
	for i, u := range u16 {
		if u != expectedUnits[i] {
			t.Errorf("StringToUTF16[%d] = 0x%04x, expected 0x%04x", i, u, expectedUnits[i])
		}
		ch := javalang.StringCharAt(s, i)
		if ch != expectedUnits[i] {
			t.Errorf("StringCharAt(%d) = 0x%04x, expected 0x%04x", i, ch, expectedUnits[i])
		}
	}

	sub := javalang.StringSubstring(s, 7, 9)
	if sub != "世界" {
		t.Errorf("StringSubstring(7, 9) = %q, expected %q", sub, "世界")
	}

	if idx := javalang.StringIndexOf(s, "世界"); idx != 7 {
		t.Errorf("StringIndexOf('世界') = %d, expected 7", idx)
	}
	if idx := javalang.StringIndexOf(s, "\u03A9"); idx != 11 {
		t.Errorf("StringIndexOf('Ω') = %d, expected 11", idx)
	}
}

func TestSupplementaryCodePointsSurrogatePairs(t *testing.T) {
	s := "A😀B𐐷C"
	expectedUnits := []uint16{
		'A', 0xD83D, 0xDE00, 'B', 0xD801, 0xDC37, 'C',
	}

	if l := javalang.StringLength(s); l != 7 {
		t.Errorf("StringLength(%q) = %d, expected 7", s, l)
	}

	for i, exp := range expectedUnits {
		ch := javalang.StringCharAt(s, i)
		if ch != exp {
			t.Errorf("StringCharAt(%d) = 0x%04x, expected 0x%04x", i, ch, exp)
		}
	}

	emojiSub := javalang.StringSubstring(s, 1, 3)
	if emojiSub != "😀" {
		t.Errorf("StringSubstring(1, 3) = %q, expected '😀'", emojiSub)
	}

	roundTrip := javalang.StringFromUTF16(expectedUnits)
	if roundTrip != s {
		t.Errorf("StringFromUTF16 = %q, expected %q", roundTrip, s)
	}
}

func TestInvalidUTF8Sequences(t *testing.T) {
	s := "bad:\xff\xfe"
	l := javalang.StringLength(s)
	if l != 6 {
		t.Errorf("StringLength of invalid utf-8 = %d, expected 6", l)
	}

	if ch := javalang.StringCharAt(s, 4); ch != 0xFFFD {
		t.Errorf("StringCharAt(4) for \\xff = 0x%04x, expected 0xFFFD", ch)
	}
	if ch := javalang.StringCharAt(s, 5); ch != 0xFFFD {
		t.Errorf("StringCharAt(5) for \\xfe = 0x%04x, expected 0xFFFD", ch)
	}
}

func TestEmptyStringSemantics(t *testing.T) {
	s := ""
	if l := javalang.StringLength(s); l != 0 {
		t.Errorf("StringLength('') = %d, expected 0", l)
	}

	assertPanicsStringIndex(t, "CharAt(0) on empty", func() { javalang.StringCharAt(s, 0) })
	assertPanicsStringIndex(t, "CharAt(-1) on empty", func() { javalang.StringCharAt(s, -1) })

	if sub := javalang.StringSubstring(s, 0, 0); sub != "" {
		t.Errorf("StringSubstring(0, 0) = %q, expected ''", sub)
	}

	assertPanicsStringIndex(t, "Substring(0, 1) on empty", func() { javalang.StringSubstring(s, 0, 1) })

	if idx := javalang.StringIndexOf(s, ""); idx != 0 {
		t.Errorf("StringIndexOf('', '') = %d, expected 0", idx)
	}
	if idx := javalang.StringIndexOf(s, "abc"); idx != -1 {
		t.Errorf("StringIndexOf('', 'abc') = %d, expected -1", idx)
	}
}

func TestASCIIFastPaths(t *testing.T) {
	s := "The quick brown fox jumps over the lazy dog"
	if !javalang.IsASCII(s) {
		t.Fatalf("Expected IsASCII(%q) to be true", s)
	}

	if javalang.StringLength(s) != len(s) {
		t.Errorf("StringLength = %d, expected %d", javalang.StringLength(s), len(s))
	}

	for i := 0; i < len(s); i++ {
		if ch := javalang.StringCharAt(s, i); ch != uint16(s[i]) {
			t.Fatalf("StringCharAt(%d) = %c, expected %c", i, ch, s[i])
		}
	}

	if sub := javalang.StringSubstring(s, 4, 9); sub != "quick" {
		t.Errorf("StringSubstring(4, 9) = %q, expected 'quick'", sub)
	}

	if idx := javalang.StringIndexOf(s, "brown"); idx != 10 {
		t.Errorf("StringIndexOf('brown') = %d, expected 10", idx)
	}

	if idx := javalang.StringLastIndexOf(s, "o"); idx != 41 {
		t.Errorf("StringLastIndexOf('o') = %d, expected 41", idx)
	}
}

// ---------------------------------------------------------------------------
// Empirical Challenger Defect Reproductions
// ---------------------------------------------------------------------------

// DEFECT 1: StringSubstring surrogate splitting corrupts code units to 0xFFFD
func TestDefect1_SurrogateSplitCorruptsCodeUnits(t *testing.T) {
	s := "a\U0001F600bc" // [ 'a', 0xD83D, 0xDE00, 'b', 'c' ]

	// Substring extracting only high surrogate [1, 2)
	hiSub := javalang.StringSubstring(s, 1, 2)
	u16Hi := javalang.StringToUTF16(hiSub)
	if len(u16Hi) != 1 || u16Hi[0] != 0xD83D {
		t.Errorf("[DEFECT 1 CONFIRMED] High surrogate substring corrupted: got u16=%x, expected [d83d]", u16Hi)
	}

	// Substring starting from low surrogate and spanning into 'b': [2, 4)
	loSub := javalang.StringSubstring(s, 2, 4)
	u16Lo := javalang.StringToUTF16(loSub)
	if len(u16Lo) != 2 || u16Lo[0] != 0xDE00 || u16Lo[1] != 'b' {
		t.Errorf("[DEFECT 1 CONFIRMED] Low surrogate substring corrupted: got u16=%x, expected [de00 62]", u16Lo)
	}
}

// DEFECT 2: StringSubstring panics on begin == end == length for non-ASCII strings
func TestDefect2_SubstringAtLengthPanicsOnNonASCII(t *testing.T) {
	s := "a\u0100b" // length 3
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("[DEFECT 2 CONFIRMED] StringSubstring(s, 3, 3) panicked: %v", r)
		}
	}()
	res := javalang.StringSubstring(s, 3, 3)
	if res != "" {
		t.Errorf("Expected empty string, got %q", res)
	}
}

// DEFECT 3: StringLastIndexOfFrom returns 0 instead of -1 when sub == "" and fromIndex < 0
func TestDefect3_LastIndexOfEmptyNegativeFromIndex(t *testing.T) {
	idx := javalang.StringLastIndexOfFrom("hello", "", -1)
	if idx != -1 {
		t.Errorf("[DEFECT 3 CONFIRMED] StringLastIndexOfFrom('hello', '', -1) = %d, expected -1", idx)
	}
}

// DEFECT 4: StringSplitLiteralLimit returns nil when limit == 0
func TestDefect4_StringSplitLimitZero(t *testing.T) {
	parts := javalang.StringSplitLiteralLimit("a:b", ":", 0)
	if parts == nil || len(parts) != 2 {
		t.Errorf("[DEFECT 4 CONFIRMED] StringSplitLiteralLimit('a:b', ':', 0) returned %v, expected [a b]", parts)
	}
}

func assertPanicsStringIndex(t *testing.T, desc string, f func()) {
	t.Helper()
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("%s did not panic", desc)
		}
	}()
	f()
}

var _ = utf16.DecodeRune
var _ = strings.Index
