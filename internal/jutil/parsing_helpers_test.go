package jutil

import (
	"math"
	"strconv"
	"strings"
	"testing"
)

// Vectors generated with tools/oracle/java/JutilParsingVectors.java (JDK 21).

// javaOutcome renders a panic like JutilParsingVectors.outcome.
func javaOutcome(r any) string {
	t, ok := r.(Throwable)
	if !ok {
		return "go-panic"
	}
	name := t.Error()
	if i := strings.Index(name, ":"); i >= 0 {
		name = name[:i]
	}
	if i := strings.LastIndex(name, "."); i >= 0 {
		name = name[i+1:]
	}
	return "E:" + name + ":" + hexUTF16(t.GetMessage())
}

func hexUTF16(s string) string {
	var sb strings.Builder
	sb.WriteByte('x')
	for _, u := range StringToUTF16(s) {
		sb.WriteString(strings.Repeat("0", 4-len(strconv.FormatUint(uint64(u), 16))))
		sb.WriteString(strconv.FormatUint(uint64(u), 16))
	}
	return sb.String()
}

func catchOutcome(f func() string) (res string) {
	defer func() {
		if r := recover(); r != nil {
			res = javaOutcome(r)
		}
	}()
	return f()
}

func TestDoubleParseDoubleMatchesJava(t *testing.T) {
	rows := readTSV(t, "parse_double.tsv")
	if len(rows) < 6000 {
		t.Fatalf("only %d vectors", len(rows))
	}
	for _, r := range rows {
		in := StringFromUTF16(decodeHexUTF16(r[0]))
		got := catchOutcome(func() string {
			return strconv.FormatUint(math.Float64bits(DoubleParseDouble(in)), 16)
		})
		if got != r[1] {
			t.Errorf("Double.parseDouble(%q) = %s, Java %s", in, got, r[1])
		}
	}
}

func TestNewBigIntegerFromStringRadixMatchesJava(t *testing.T) {
	rows := readTSV(t, "bigint_radix.tsv")
	if len(rows) < 700 {
		t.Fatalf("only %d vectors", len(rows))
	}
	for _, r := range rows {
		in := StringFromUTF16(decodeHexUTF16(r[0]))
		radix, _ := strconv.Atoi(r[1])
		got := catchOutcome(func() string { return NewBigIntegerFromStringRadix(in, radix).String() })
		want := r[2]
		if got != want {
			t.Errorf("new BigInteger(%q, %d) = %s, Java %s", in, radix, got, want)
		}
	}
}

func TestCharMatcherWhitespace(t *testing.T) {
	want := map[uint16]bool{}
	for _, u := range StringToUTF16(charMatcherWhitespaceTable) {
		want[u] = true
	}
	if len(want) != 25 {
		t.Fatalf("table has %d distinct chars", len(want))
	}
	for c := 0; c < 0x10000; c++ {
		if got := CharMatcher_Whitespace_Matches(uint16(c)); got != want[uint16(c)] {
			t.Errorf("whitespace().matches(%#x) = %v", c, got)
		}
	}
}

func TestRuntimeExceptionWithMessageAndCause(t *testing.T) {
	cause := NewIllegalStateExceptionMsg("inner")
	e := NewRuntimeExceptionWithMessageAndCause("outer", cause)
	if e.Error() != "java.lang.RuntimeException: outer" || e.GetMessage() != "outer" || e.GetCause() != Throwable(cause) {
		t.Errorf("got %q %q %v", e.Error(), e.GetMessage(), e.GetCause())
	}
}
