package dtoa

import (
	"bufio"
	"math"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// testdata/dtoa.tsv is the output of DToA.numberToString from the pinned
// closure-compiler.jar, written by tools/oracle/java/DToAVectors.java.
func TestNumberToStringMatchesJava(t *testing.T) {
	f, err := os.Open("testdata/dtoa.tsv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	n, bad := 0, 0
	for sc.Scan() {
		parts := strings.Split(sc.Text(), "\t")
		bits, err := strconv.ParseUint(parts[0], 16, 64)
		if err != nil {
			t.Fatal(err)
		}
		d := math.Float64frombits(bits)
		n++
		got := ""
		func() {
			defer func() {
				if r := recover(); r != nil {
					got = "!" + jutil.AsThrowable(r).Error()
				}
			}()
			got = DToA_NumberToString(d)
		}()
		if got != parts[1] {
			bad++
			if bad <= 20 {
				t.Errorf("numberToString(%x = %v) = %q, Java %q", bits, d, got, parts[1])
			}
		}
	}
	if n < 50000 {
		t.Fatalf("only %d vectors", n)
	}
	if bad > 0 {
		t.Fatalf("%d/%d mismatches", bad, n)
	}
}

func TestNumberToStringExamples(t *testing.T) {
	cases := []struct {
		in   float64
		want string
	}{
		{0.0, "0"},
		{math.Copysign(0, -1), "0"},
		{1, "1"},
		{-1, "-1"},
		{0.5, "0.5"},
		{-0.5, "-0.5"},
		{12.34, "12.34"},
		{123456789, "123456789"},
		{1e+5, "100000"},
		{1.5e+5, "150000"},
		{1.5e+5 + 0.25, "150000.25"},
		{9007199254740991, "9007199254740991"},
		{-9007199254740991, "-9007199254740991"},
		{9007199254740992, "9007199254740992"},
		{-9007199254740992, "-9007199254740992"},
		{1e-6, "0.000001"},
		{-1e-6, "-0.000001"},
		{math.Nextafter(1e-6, 0), "9.999999999999997e-7"},
		{-math.Nextafter(1e-6, 0), "-9.999999999999997e-7"},
		{1e-7, "1e-7"},
		{-1e-7, "-1e-7"},
		{1.25e-7, "1.25e-7"},
		{1e20, "100000000000000000000"},
		{-1e20, "-100000000000000000000"},
		{999999999999999900000.0, "999999999999999900000"},
		{1e21, "1e+21"},
		{-1e21, "-1e+21"},
		{1.5e21, "1.5e+21"},
		{1.5e300, "1.5e+300"},
		{math.MaxFloat64, "1.7976931348623157e+308"},
		{-math.MaxFloat64, "-1.7976931348623157e+308"},
		{math.SmallestNonzeroFloat64, "5e-324"},
		{-math.SmallestNonzeroFloat64, "-5e-324"},
		{math.Float64frombits(0x000fffffffffffff), "2.225073858507201e-308"},  // max subnormal
		{math.Float64frombits(0x0010000000000000), "2.2250738585072014e-308"}, // min normal
		{math.NaN(), "NaN"},
		{math.Inf(1), "Infinity"},
		{math.Inf(-1), "-Infinity"},
	}
	for _, tc := range cases {
		got := DToA_NumberToString(tc.in)
		if got != tc.want {
			t.Errorf("DToA_NumberToString(%x = %g) = %q, want %q", math.Float64bits(tc.in), tc.in, got, tc.want)
		}
	}
}

var dtoaSink string

func TestNumberToStringAllocs(t *testing.T) {
	inputs := []float64{
		0.5, -12.34, 1e-6, 1e-7, math.Nextafter(1e-6, 0),
		1e20, 1e21, math.MaxFloat64, math.Float64frombits(0x000fffffffffffff),
	}
	for _, v := range inputs {
		allocs := testing.AllocsPerRun(100, func() {
			dtoaSink = DToA_NumberToString(v)
		})
		// At most 1 alloc for the returned string header/backing bytes (0 for short strings if optimized), never *big.Int allocs.
		if allocs > 1 {
			t.Fatalf("DToA_NumberToString(%g) allocs = %v, want <= 1", v, allocs)
		}
	}
}
