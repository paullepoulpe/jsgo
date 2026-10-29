package jutil

import (
	"bufio"
	"compress/gzip"
	"math"
	"os"
	"strconv"
	"strings"
	"testing"
	"unicode"
)

// Vectors generated with tools/oracle/java/JutilPeepholeVectors.java (JDK 21
// on x86_64, shaded Guava from the pinned closure-compiler.jar).

func readGzipTSV(t *testing.T, name string) [][]string {
	t.Helper()
	f, err := os.Open("testdata/" + name)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	var rows [][]string
	sc := bufio.NewScanner(gz)
	sc.Buffer(make([]byte, 1<<20), 1<<20)
	for sc.Scan() {
		rows = append(rows, strings.Split(sc.Text(), "\t"))
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}
	return rows
}

func hexDouble(t *testing.T, s string) float64 {
	t.Helper()
	v, err := strconv.ParseUint(s, 16, 64)
	if err != nil {
		t.Fatalf("bad double %q", s)
	}
	return math.Float64frombits(v)
}

func catchArithmetic(f func() string) (res string) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(IArithmeticException); ok {
				res = "AE:" + e.GetMessage()
				return
			}
			panic(r)
		}
	}()
	return f()
}

func TestJavaMathMatchesJava(t *testing.T) {
	rows := readGzipTSV(t, "java_math.tsv.gz")
	if len(rows) < 40000 {
		t.Fatalf("only %d vectors", len(rows))
	}
	counts := map[string]int{}
	for _, r := range rows {
		op := r[0]
		counts[op]++
		switch op {
		case "rem":
			x, y, want := hexDouble(t, r[1]), hexDouble(t, r[2]), hexDouble(t, r[3])
			got := DRem(x, y)
			if math.IsNaN(want) && math.IsNaN(got) {
				continue
			}
			if math.Float64bits(got) != math.Float64bits(want) {
				t.Errorf("rem(%v, %v) = %v (%x), Java %v (%x)", x, y, got, math.Float64bits(got), want, math.Float64bits(want))
			}
		case "pow", "max", "min", "round", "signum":
			// Math* moved to internal/base.
		case "intValueExact":
			a := NewBigIntegerFromString(r[1])
			got := catchArithmetic(func() string { return strconv.Itoa(int(BigIntegerIntValueExact(a))) })
			if got != r[3] {
				t.Errorf("intValueExact(%s) = %s, Java %s", r[1], got, r[3])
			}
		case "mod":
			a, m := NewBigIntegerFromString(r[1]), NewBigIntegerFromString(r[2])
			got := catchArithmetic(func() string { return BigIntegerMod(a, m).String() })
			if got != r[3] {
				t.Errorf("%s.mod(%s) = %s, Java %s", r[1], r[2], got, r[3])
			}
		default:
			t.Fatalf("unknown op %q", op)
		}
	}
}

func TestStringOpsMatchJava(t *testing.T) {
	rows := readTSV(t, "string_ops.tsv")
	if len(rows) < 3000 {
		t.Fatalf("only %d vectors", len(rows))
	}
	for _, r := range rows {
		s := StringFromUTF16(decodeUnits(t, r[1]))
		u := StringFromUTF16(decodeUnits(t, r[2]))
		switch r[0] {
		case "lastIndexOf":
			from, _ := strconv.Atoi(r[3])
			want, _ := strconv.Atoi(r[4])
			if got := StringLastIndexOfFrom(s, u, from); got != want {
				t.Errorf("%q.lastIndexOf(%q, %d) = %d, Java %d", s, u, from, got, want)
			}
		case "equalsIgnoreCase":
			// Ascii.equalsIgnoreCase was removed with ascii.go.
		default:
			t.Fatalf("unknown op %q", r[0])
		}
	}
}

func TestStringCaseRootMatchesJava(t *testing.T) {
	rows := readTSV(t, "string_case_root.tsv")
	simpleUpper := map[rune]rune{}
	simpleLower := map[rune]rune{}
	parseCP := func(s string) rune {
		v, err := strconv.ParseUint(s, 16, 32)
		if err != nil {
			t.Fatalf("bad code point %q", s)
		}
		return rune(v)
	}
	for _, r := range rows {
		if r[0] == "simple" {
			cp := parseCP(r[1])
			simpleUpper[cp] = parseCP(r[2])
			simpleLower[cp] = parseCP(r[3])
		}
	}
	if len(simpleUpper) < 2000 {
		t.Fatalf("only %d simple mappings", len(simpleUpper))
	}
	for cp := rune(0); cp <= unicode.MaxRune; cp++ {
		if cp >= 0xD800 && cp <= 0xDFFF {
			continue
		}
		wantU, wantL := cp, cp
		if v, ok := simpleUpper[cp]; ok {
			wantU, wantL = v, simpleLower[cp]
		}
		if unicode15ToUpper(cp) != wantU || unicode15ToLower(cp) != wantL {
			t.Errorf("U+%04X: Go upper/lower %X/%X, Java %X/%X", cp, unicode15ToUpper(cp), unicode15ToLower(cp), wantU, wantL)
		}
	}
}
