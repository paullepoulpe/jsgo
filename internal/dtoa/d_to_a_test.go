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
	for d, want := range map[float64]string{
		1: "1", 0.5: "0.5", 1e21: "1e+21", 1e-7: "1e-7", 123456789: "123456789", -0.0: "0",
		math.Inf(-1): "-Infinity", 1.5e300: "1.5e+300", 0.000001: "0.000001",
	} {
		if got := DToA_NumberToString(d); got != want {
			t.Errorf("numberToString(%v) = %q, want %q", d, got, want)
		}
	}
}
