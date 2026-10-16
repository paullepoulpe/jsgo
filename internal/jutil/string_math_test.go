package jutil

import (
	"math"
	"strconv"
	"testing"
)

// Vectors from tools/oracle/java/JutilStringMathVectors.java (JDK 21):
// Math.round/max/min/signum, String.lastIndexOf(String, int) and
// Matcher.replaceAll/replaceFirst.

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

func doubleFromBits(t *testing.T, s string) float64 {
	t.Helper()
	b, err := strconv.ParseUint(s, 16, 64)
	if err != nil {
		t.Fatal(err)
	}
	return math.Float64frombits(b)
}

func TestJavaMathMatchesJava(t *testing.T) {
	for _, r := range stringMathRows(t, "round", 1000) {
		a := doubleFromBits(t, r[0])
		want, err := strconv.ParseInt(r[1], 10, 64)
		if err != nil {
			t.Fatal(err)
		}
		if got := MathRound(a); got != want {
			t.Errorf("Math.round(%v [%s]) = %d, want %d", a, r[0], got, want)
		}
	}
	for _, r := range stringMathRows(t, "signum", 1000) {
		a, want := doubleFromBits(t, r[0]), doubleFromBits(t, r[1])
		if got := MathSignum(a); math.Float64bits(got) != math.Float64bits(want) {
			t.Errorf("Math.signum(%s) = %016x, want %s", r[0], math.Float64bits(got), r[1])
		}
	}
	for _, op := range []struct {
		name string
		f    func(a, b float64) float64
	}{{"max", MathMax}, {"min", MathMin}} {
		for _, r := range stringMathRows(t, op.name, 1000) {
			a, b, want := doubleFromBits(t, r[0]), doubleFromBits(t, r[1]), doubleFromBits(t, r[2])
			got := op.f(a, b)
			if a != a && b != b {
				// Observed behaviour (black-box): with two NaNs the JDK
				// returns either one, depending on JIT compilation.
				if got == got {
					t.Errorf("Math.%s(%s, %s) = %v, want NaN", op.name, r[0], r[1], got)
				}
				continue
			}
			if math.Float64bits(got) != math.Float64bits(want) {
				t.Errorf("Math.%s(%s, %s) = %016x, want %s", op.name, r[0], r[1], math.Float64bits(got), r[2])
			}
		}
	}
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
