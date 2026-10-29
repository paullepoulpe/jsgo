package jutil

import (
	"bufio"
	"compress/gzip"
	"math"
	"math/big"
	"os"
	"strconv"
	"strings"
	"testing"
)

// Vectors generated with tools/oracle/java/JutilNumbersVectors.java (JDK 21).
// NUMBERS_VECTORS=<file.tsv.gz> replays another run (e.g. another seed).

// numbersOutcome runs f like JutilNumbersVectors.run: "=" + the result, or
// "E:" + the hex of the exception's toString().
func numbersOutcome(f func() string) (res string) {
	defer func() {
		if r := recover(); r != nil {
			if t, ok := r.(Throwable); ok {
				res = "E:" + hexUTF16(t.Error())
			} else {
				res = "go-panic"
			}
		}
	}()
	return "=" + f()
}

func TestNumbersMatchJava(t *testing.T) {
	path := os.Getenv("NUMBERS_VECTORS")
	if path == "" {
		path = "testdata/numbers.tsv.gz"
	}
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	sc := bufio.NewScanner(gz)
	sc.Buffer(make([]byte, 1<<20), 1<<20)
	counts := map[string]int{}
	failures := 0
	check := func(kind, desc, got, want string) {
		if got != want {
			failures++
			if failures <= 40 {
				t.Errorf("%s %s:\n   go: %s\n java: %s", kind, desc, got, want)
			}
		}
	}
	for sc.Scan() {
		r := strings.Split(sc.Text(), "\t")
		counts[r[0]]++
		switch r[0] {
		case "int", "long":
			s := StringFromUTF16(decodeHexUTF16(r[1]))
			radix, _ := strconv.Atoi(r[2])
			got := numbersOutcome(func() string {
				if r[0] == "int" {
					return strconv.FormatInt(int64(IntegerParseInt(s, radix)), 10)
				}
				return strconv.FormatInt(LongParseLong(s, radix), 10)
			})
			check(r[0], strconv.Quote(s)+" radix "+r[2], got, r[3])
		case "bd":
			s := StringFromUTF16(decodeHexUTF16(r[1]))
			got := numbersOutcome(func() string {
				d := NewBigDecimalFromString(s)
				tbi := "-"
				if scale := d.Scale(); (scale >= -2000 && scale <= 2000) || scale >= 536870920 || scale <= -536870920 {
					tbi = numbersOutcome(func() string { return d.ToBigInteger().String() })
				}
				return d.unscaled.String() + "|" + strconv.Itoa(int(d.Scale())) + "|" + strconv.Itoa(d.Precision()) +
					"|" + strconv.Itoa(d.Signum()) + "|" + strconv.Itoa(int(d.IntValue())) + "|" +
					strconv.FormatInt(d.LongValue(), 10) + "|" + tbi
			})
			check("BigDecimal", strconv.Quote(s), got, r[2])
		case "pd":
			s := StringFromUTF16(decodeHexUTF16(r[1]))
			got := numbersOutcome(func() string {
				return strconv.FormatUint(math.Float64bits(DoubleParseDouble(s)), 16)
			})
			check("Double.parseDouble", strconv.Quote(s), got, r[2])
		case "bir":
			s := StringFromUTF16(decodeHexUTF16(r[1]))
			radix, _ := strconv.Atoi(r[2])
			got := numbersOutcome(func() string { return NewBigIntegerFromStringRadix(s, radix).String() })
			check("new BigInteger", strconv.Quote(s)+" radix "+r[2], got, r[3])
		case "pow":
			base, _ := new(big.Int).SetString(r[1], 10)
			e, _ := strconv.Atoi(r[2])
			got := numbersOutcome(func() string { return BigIntegerPow(base, int32(e)).String() })
			check("BigInteger.pow", r[1]+"^"+r[2], got, r[3])
		case "fmt":
			// Formatter_FormatDecimalFloat moved to internal/jscomp.
		case "bs":
			var a []int
			if r[1] != "-" {
				for _, v := range strings.Split(r[1], ",") {
					n, _ := strconv.Atoi(v)
					a = append(a, n)
				}
			}
			var res []string
			for _, k := range strings.Split(r[2], ",") {
				key, _ := strconv.Atoi(k)
				res = append(res, strconv.Itoa(ArraysBinarySearchInts(a, key)))
			}
			check("Arrays.binarySearch", r[1], strings.Join(res, ","), r[3])
		case "rnd":
			seed, _ := strconv.ParseInt(r[1], 10, 64)
			x := NewRandomWithSeed(seed)
			bounds := strings.Split(r[2], ",")
			if strings.HasPrefix(r[3], "E:") {
				b, _ := strconv.Atoi(bounds[0])
				got := numbersOutcome(func() string { return strconv.Itoa(int(x.NextIntBound(int32(b)))) })
				check("Random.nextInt", r[2], got, r[3])
				continue
			}
			var res []string
			for _, bs := range bounds {
				b, _ := strconv.Atoi(bs)
				res = append(res, strconv.Itoa(int(x.NextIntBound(int32(b)))))
			}
			check("Random.nextInt", "seed "+r[1]+" bounds "+r[2], strings.Join(res, ","), r[3])
		}
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}
	for kind, min := range map[string]int{"int": 2000, "long": 2000, "bd": 2500, "pd": 10000, "bir": 2000,
		"pow": 400, "fmt": 2500, "bs": 1500, "rnd": 300} {
		if counts[kind] < min {
			t.Errorf("only %d %q vectors", counts[kind], kind)
		}
	}
	if failures > 0 {
		t.Errorf("%d vectors differ", failures)
	}
}
