package jutil

import (
	"strconv"
	"testing"
)

// Vectors from tools/oracle/java/JutilStringMathVectors.java (JDK 21):
// String.lastIndexOf(String, int) and Matcher.replaceAll/replaceFirst.

func stringMathRows(t *testing.T, kind string, minRows int) [][]string {
	t.Helper()
	var rows [][]string
	for _, r := range readTSV(t, "string_math.tsv") {
		if r[0] == kind {
			rows = append(rows, r[1:])
		}
	}
	if len(rows) < minRows {
		t.Fatalf("only %d %s vectors", len(rows), kind)
	}
	return rows
}

func TestStringLastIndexOfFromMatchesJava(t *testing.T) {
	for _, r := range stringMathRows(t, "lastidx", 1000) {
		s, sub := StringFromUTF16(decodeUnits(t, r[0])), StringFromUTF16(decodeUnits(t, r[1]))
		from, err := strconv.Atoi(r[2])
		if err != nil {
			t.Fatal(err)
		}
		want, err := strconv.Atoi(r[3])
		if err != nil {
			t.Fatal(err)
		}
		if got := StringLastIndexOfFrom(s, sub, from); got != want {
			t.Errorf("%s.lastIndexOf(%s, %d) = %d, want %d", r[0], r[1], from, got, want)
		}
	}
}

func TestMatcherReplaceMatchesJava(t *testing.T) {
	outcome := func(f func() string) (res string) {
		Try(func() { res = encodeUnits(StringToUTF16(f())) }, func(e Throwable) { res = "!" + e.Error() })
		return res
	}
	for _, r := range stringMathRows(t, "replace", 1000) {
		regex := StringFromUTF16(decodeUnits(t, r[0]))
		input := StringFromUTF16(decodeUnits(t, r[1]))
		repl := StringFromUTF16(decodeUnits(t, r[2]))
		p := Pattern_Compile(regex)
		if got := outcome(func() string { return p.Matcher(input).ReplaceAll(repl) }); got != r[3] {
			t.Errorf("Pattern.compile(%q).matcher(%q).replaceAll(%q) = %s, want %s", regex, input, repl, got, r[3])
		}
		if got := outcome(func() string { return p.Matcher(input).ReplaceFirst(repl) }); got != r[4] {
			t.Errorf("Pattern.compile(%q).matcher(%q).replaceFirst(%q) = %s, want %s", regex, input, repl, got, r[4])
		}
	}
}
