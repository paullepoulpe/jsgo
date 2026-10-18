package jutil

import (
	"bufio"
	"compress/gzip"
	unicode "github.com/paullepoulpe/jsgo/internal/unicode15"
	"math"
	"math/big"
	"os"
	"strconv"
	"strings"
	"testing"
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

// isOneUlpApart reports whether a and b are adjacent doubles.
func isOneUlpApart(a, b float64) bool {
	if math.Signbit(a) != math.Signbit(b) {
		return false
	}
	d := int64(math.Float64bits(a)) - int64(math.Float64bits(b))
	return d == 1 || d == -1
}

// powDivergenceAllowed reports whether MathPow may differ from HotSpot's
// Math.pow intrinsic on these arguments (see MathPow): the intrinsic is up to
// 1 ulp from the exact result when that is not representable, and it treats
// y = ±(2^52 - 0.5) as an odd integer for x = -0 and x = -Inf, unlike the
// Java specification (StrictMath.pow follows the specification there).
func powDivergenceAllowed(x, y, got, want float64) bool {
	if (x == 0 || math.IsInf(x, -1)) && math.Signbit(x) && math.Abs(y) == 1<<52-0.5 {
		return got == -want
	}
	if !isOneUlpApart(got, want) {
		return false
	}
	if x == math.Trunc(x) && y == math.Trunc(y) && math.Abs(x) <= 1<<20 && math.Abs(y) <= 64 {
		// Integer powers: the exact value must not be representable.
		exact := new(big.Rat).SetInt(new(big.Int).Exp(big.NewInt(int64(math.Abs(x))), big.NewInt(int64(math.Abs(y))), nil))
		if y < 0 {
			exact.Inv(exact)
		}
		if f, isExact := exact.Float64(); isExact && f != 0 && !math.IsInf(f, 0) {
			return false
		}
	}
	return true
}

func TestJavaMathMatchesJava(t *testing.T) {
	rows := readGzipTSV(t, "java_math.tsv.gz")
	if len(rows) < 40000 {
		t.Fatalf("only %d vectors", len(rows))
	}
	counts := map[string]int{}
	powInexact := 0
	for _, r := range rows {
		op := r[0]
		counts[op]++
		switch op {
		case "pow", "max", "min", "rem":
			x, y, want := hexDouble(t, r[1]), hexDouble(t, r[2]), hexDouble(t, r[3])
			var got float64
			switch op {
			case "pow":
				got = MathPow(x, y)
			case "max":
				got = MathMax(x, y)
			case "min":
				got = MathMin(x, y)
			case "rem":
				got = DRem(x, y)
			}
			if math.IsNaN(want) && math.IsNaN(got) {
				continue
			}
			if math.Float64bits(got) != math.Float64bits(want) {
				if op == "pow" && powDivergenceAllowed(x, y, got, want) {
					powInexact++
					continue
				}
				t.Errorf("%s(%v, %v) = %v (%x), Java %v (%x)", op, x, y, got, math.Float64bits(got), want, math.Float64bits(want))
			}
		case "round":
			x := hexDouble(t, r[1])
			want, _ := strconv.ParseInt(r[3], 10, 64)
			if got := MathRound(x); got != want {
				t.Errorf("round(%v) = %d, Java %d", x, got, want)
			}
		case "signum":
			x, want := hexDouble(t, r[1]), hexDouble(t, r[3])
			got := MathSignum(x)
			if !(math.IsNaN(got) && math.IsNaN(want)) && math.Float64bits(got) != math.Float64bits(want) {
				t.Errorf("signum(%v) = %v, Java %v", x, got, want)
			}
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
	// The intrinsic is 1 ulp from the correctly rounded value on ~0.05% of
	// random arguments, and rarely on integer powers.
	if powInexact*1000 > counts["pow"] {
		t.Errorf("%d of %d pow vectors are 1 ulp away from Java", powInexact, counts["pow"])
	}
	t.Logf("vectors: %v; pow vectors on which HotSpot's intrinsic differs: %d", counts, powInexact)
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
			want := r[4] == "true"
			if got := AsciiEqualsIgnoreCase(s, u); got != want {
				t.Errorf("Ascii.equalsIgnoreCase(%q, %q) = %v, Java %v", s, u, got, want)
			}
		default:
			t.Fatalf("unknown op %q", r[0])
		}
	}
}

func TestStringCaseRootMatchesJava(t *testing.T) {
	rows := readTSV(t, "string_case_root.tsv")
	simpleUpper := map[rune]rune{}
	simpleLower := map[rune]rune{}
	upper := map[rune]string{}
	parseCP := func(s string) rune {
		v, err := strconv.ParseUint(s, 16, 32)
		if err != nil {
			t.Fatalf("bad code point %q", s)
		}
		return rune(v)
	}
	strings_ := 0
	for _, r := range rows {
		switch r[0] {
		case "simple":
			cp := parseCP(r[1])
			simpleUpper[cp] = parseCP(r[2])
			simpleLower[cp] = parseCP(r[3])
		case "upper":
			upper[parseCP(r[1])] = StringFromUTF16(decodeUnits(t, r[2]))
		case "string":
			strings_++
			s := StringFromUTF16(decodeUnits(t, r[1]))
			if got, want := StringToUpperCaseRoot(s), StringFromUTF16(decodeUnits(t, r[2])); got != want {
				t.Errorf("%q.toUpperCase(ROOT) = %q, Java %q", s, got, want)
			}
		}
	}
	if len(simpleUpper) < 2000 || strings_ < 2000 {
		t.Fatalf("only %d simple mappings, %d strings", len(simpleUpper), strings_)
	}
	// Character.toUpperCase/toLowerCase(int) is Go's unicode.ToUpper/ToLower
	// on every code point.
	for cp := rune(0); cp <= unicode.MaxRune; cp++ {
		if cp >= 0xD800 && cp <= 0xDFFF {
			continue
		}
		wantU, wantL := cp, cp
		if v, ok := simpleUpper[cp]; ok {
			wantU, wantL = v, simpleLower[cp]
		}
		if unicode.ToUpper(cp) != wantU || unicode.ToLower(cp) != wantL {
			t.Errorf("U+%04X: Go upper/lower %X/%X, Java %X/%X", cp, unicode.ToUpper(cp), unicode.ToLower(cp), wantU, wantL)
		}
		s := string(cp)
		wantUS, ok := upper[cp]
		if !ok {
			wantUS = string(wantU)
		}
		if got := StringToUpperCaseRoot(s); got != wantUS {
			t.Errorf("U+%04X toUpperCase(ROOT) = %q, Java %q", cp, got, wantUS)
		}
	}
}
