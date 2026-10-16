package jutil

import (
	"bufio"
	"math"
	"os"
	"strconv"
	"strings"
	"testing"
)

// testdata/character.tsv was generated with JDK 21: run-length-encoded lines
// "<startHex> <endHex> <getType> <isLetter> <isDigit> <digit16> <digit10>
// <digit36>" covering every UTF-16 code unit.
func TestCharacterMatchesJava(t *testing.T) {
	f, err := os.Open("testdata/character.tsv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	covered := 0
	bad := 0
	for sc.Scan() {
		p := strings.Fields(sc.Text())
		lo, _ := strconv.ParseUint(p[0], 16, 32)
		hi, _ := strconv.ParseUint(p[1], 16, 32)
		typ, _ := strconv.Atoi(p[2])
		letter := p[3] == "1"
		digit := p[4] == "1"
		d16, _ := strconv.Atoi(p[5])
		d10, _ := strconv.Atoi(p[6])
		d36, _ := strconv.Atoi(p[7])
		for c := lo; c <= hi; c++ {
			ch := uint16(c)
			covered++
			if got := CharacterGetType(ch); got != typ {
				bad++
				if bad < 20 {
					t.Errorf("CharacterGetType(%04X) = %d, Java %d", c, got, typ)
				}
			}
			if got := CharacterIsLetter(ch); got != letter {
				bad++
				if bad < 20 {
					t.Errorf("CharacterIsLetter(%04X) = %v, Java %v", c, got, letter)
				}
			}
			if got := CharacterIsDigit(ch); got != digit {
				bad++
				if bad < 20 {
					t.Errorf("CharacterIsDigit(%04X) = %v, Java %v", c, got, digit)
				}
			}
			if got := CharacterDigit(ch, 16); got != d16 {
				bad++
				if bad < 20 {
					t.Errorf("CharacterDigit(%04X, 16) = %d, Java %d", c, got, d16)
				}
			}
			if got := CharacterDigit(ch, 10); got != d10 {
				bad++
				if bad < 20 {
					t.Errorf("CharacterDigit(%04X, 10) = %d, Java %d", c, got, d10)
				}
			}
			if got := CharacterDigit(ch, 36); got != d36 {
				bad++
				if bad < 20 {
					t.Errorf("CharacterDigit(%04X, 36) = %d, Java %d", c, got, d36)
				}
			}
		}
	}
	if covered != 0x10000 {
		t.Errorf("table covers %d code units, want 65536", covered)
	}
}

// Vectors generated with JDK 21: Integer.parseInt(s, 16).
func TestIntegerParseIntMatchesJava(t *testing.T) {
	cases := []struct {
		in   []uint16
		want string // decimal result, or "NFE"
	}{
		{[]uint16{}, "NFE"},
		{[]uint16{'+'}, "NFE"},
		{[]uint16{'-'}, "NFE"},
		{[]uint16{'0'}, "0"},
		{[]uint16{'f', 'F'}, "255"},
		{[]uint16{'+', '7', 'f'}, "127"},
		{[]uint16{'-', '7', 'f'}, "-127"},
		{StringToUTF16("7fffffff"), "2147483647"},
		{StringToUTF16("80000000"), "NFE"},
		{StringToUTF16("-80000000"), "-2147483648"},
		{StringToUTF16("-80000001"), "NFE"},
		{StringToUTF16("ffffffff"), "NFE"},
		{StringToUTF16("100000000"), "NFE"},
		{[]uint16{0x0663, 0x0664}, "52"},
		{[]uint16{0xFF21, 0xFF42}, "171"},
		{[]uint16{0xFF10}, "0"},
		{StringToUTF16("12z"), "NFE"},
		{StringToUTF16(" 1"), "NFE"},
		{StringToUTF16("0x10"), "NFE"},
		{StringToUTF16("00000000041"), "65"},
		{StringToUTF16("110000"), "1114112"},
		{[]uint16{0x0966, 0x0967}, "1"},
		{[]uint16{0xD835, 0xDFCE}, "NFE"},
		{StringToUTF16("--1"), "NFE"},
		{StringToUTF16("+-1"), "NFE"},
		{StringToUTF16("D800"), "55296"},
		{StringToUTF16("{41}"), "NFE"},
	}
	for _, c := range cases {
		s := StringFromUTF16(c.in)
		got := ""
		Try(func() {
			got = strconv.FormatInt(int64(IntegerParseInt(s, 16)), 10)
		}, func(e INumberFormatException) {
			got = "NFE"
		})
		if got != c.want {
			t.Errorf("IntegerParseInt(%q, 16) = %s, Java %s", s, got, c.want)
		}
	}
	var msg string
	Try(func() { IntegerParseInt("", 16) }, func(e INumberFormatException) { msg = e.GetMessage() })
	if msg != `For input string: "" under radix 16` {
		t.Errorf("message = %q", msg)
	}
}

// Expected outputs generated with JDK 21 String.format.
func TestStringFormatMatchesJava(t *testing.T) {
	cases := []struct {
		format string
		args   []any
		want   string
	}{
		{"Character '%c' (U+%04X) is not a valid identifier start char", []any{uint16('x'), int('x')},
			"Character 'x' (U+0078) is not a valid identifier start char"},
		{"Character '%c' (U+%04X) is not a valid identifier start char", []any{uint16(0xe9), 0xe9},
			"Character 'é' (U+00E9) is not a valid identifier start char"},
		{"Character '%c' (U+%04X)", []any{uint16(0xd800), 0xd800}, "Character '" + StringFromChar(0xd800) + "' (U+D800)"},
		{"%04X", []any{0x12345}, "12345"},
		{"%04X", []any{-1}, "FFFFFFFF"},
		{"%x", []any{int64(-1)}, "ffffffffffffffff"},
		{"Setter must have exactly 1 parameter, found %d", []any{3}, "Setter must have exactly 1 parameter, found 3"},
		{"'%s' expected", []any{nil}, "'null' expected"},
		{"%s(%d, %d)", []any{"a.js", 1, 2}, "a.js(1, 2)"},
		{"%5s|%-5s|%.2s|%%|%n", []any{"ab", "cd", "efg"}, "   ab|cd   |ef|%|\n"},
		{"%2$s %1$s", []any{"a", "b"}, "b a"},
		{"%05d|%-4d|%d", []any{-42, 7, int64(math.MinInt64)}, "-0042|7   |-9223372036854775808"},
		{"%b %b %b", []any{true, nil, "x"}, "true false true"},
		{"%o", []any{8}, "10"},
	}
	for _, c := range cases {
		if got := StringFormat(c.format, c.args...); got != c.want {
			t.Errorf("StringFormat(%q) = %q, Java %q", c.format, got, c.want)
		}
	}
	for _, c := range []struct {
		format string
		args   []any
		want   string
	}{
		{"%s %s", []any{"a"}, "java.util.MissingFormatArgumentException: Format specifier '%s'"},
		{"%q", []any{1}, "java.util.UnknownFormatConversionException: Conversion = 'q'"},
	} {
		got := ""
		Try(func() { StringFormat(c.format, c.args...) }, func(e IIllegalArgumentException) { got = e.Error() })
		if got != c.want {
			t.Errorf("StringFormat(%q) threw %q, Java %q", c.format, got, c.want)
		}
	}
}

type testEnum int32

func (e testEnum) String() string { return "e" + strconv.Itoa(int(e)) }

func TestEnumSet(t *testing.T) {
	s := NewEnumSet[testEnum](70, 3, 1, 3)
	if s.Size() != 3 || !s.Contains(70) || s.Contains(2) || s.Contains(200) {
		t.Fatalf("bad set %v", s)
	}
	if got := s.String(); got != "[e1, e3, e70]" {
		t.Errorf("String() = %q", got)
	}
	o := EnumSetCopyOf(s)
	o.Remove(70)
	if !s.ContainsAll(o) || o.ContainsAll(s) || s.Equals(o) {
		t.Errorf("containsAll/equals wrong")
	}
	o.Add(70)
	if !s.Equals(o) || s.HashCode() != o.HashCode() {
		t.Errorf("equals after re-add wrong")
	}
	o.RemoveAll(NewEnumSet[testEnum](1, 70))
	if got := o.String(); got != "[e3]" {
		t.Errorf("after RemoveAll: %q", got)
	}
	if !NewEnumSet[testEnum]().Equals(&EnumSet[testEnum]{words: []uint64{0, 0}}) {
		t.Errorf("empty sets with different word counts must be equal")
	}
}

func TestStringTrim(t *testing.T) {
	for in, want := range map[string]string{
		"":                  "",
		"  \t\n x y \r\x00": "x y",
		" x ":               " x ", // Java trim() only strips <= U+0020
		"\x01\x1f":          "",
	} {
		if got := StringTrim(in); got != want {
			t.Errorf("StringTrim(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestAssertionErrorIsNotRuntimeException(t *testing.T) {
	caughtAsRuntime := false
	caughtAsError := false
	Try(func() {
		Try(func() { panic(NewAssertionError("x")) }, func(e IRuntimeException) { caughtAsRuntime = true })
	}, func(e IError) { caughtAsError = e.Error() == "java.lang.AssertionError: x" })
	if caughtAsRuntime || !caughtAsError {
		t.Errorf("caughtAsRuntime=%v caughtAsError=%v", caughtAsRuntime, caughtAsError)
	}
}
