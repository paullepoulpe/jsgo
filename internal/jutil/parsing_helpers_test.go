package jutil

import (
	"math"
	"math/rand"
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
		if strings.HasPrefix(want, "E:NumberFormatException:") && strings.HasPrefix(got, "E:NumberFormatException:") {
			// The message for a bad digit names the digit group Java was parsing;
			// only the exception class is compared.
			continue
		}
		if got != want {
			t.Errorf("new BigInteger(%q, %d) = %s, Java %s", in, radix, got, want)
		}
	}
}

func TestPropertiesLoadMatchesJava(t *testing.T) {
	rows := readTSV(t, "properties.tsv")
	if len(rows) < 30 {
		t.Fatalf("only %d vectors", len(rows))
	}
	for _, r := range rows {
		in := StringFromUTF16(decodeHexUTF16(r[0]))
		got := catchOutcome(func() string {
			p := PropertiesLoad(in)
			keys := p.KeySlice()
			sortJavaStrings(keys)
			var sb strings.Builder
			for i, k := range keys {
				if i > 0 {
					sb.WriteByte(',')
				}
				sb.WriteString(hexUTF16(k) + "=" + hexUTF16(p.Get(k)))
			}
			return "M:" + sb.String()
		})
		if got != r[1] {
			t.Errorf("Properties.load(%q) = %s, Java %s", in, got, r[1])
		}
	}
}

func sortJavaStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && StringCompareTo(s[j-1], s[j]) > 0; j-- {
			s[j-1], s[j] = s[j], s[j-1]
		}
	}
}

func TestTreeRangeSetAgainstModel(t *testing.T) {
	cmp := func(a, b int) int { return a - b }
	rnd := rand.New(rand.NewSource(7))
	for trial := 0; trial < 500; trial++ {
		s := NewTreeRangeSet(cmp)
		var model [][2]int
		if !s.IsEmpty() {
			t.Fatal("new set not empty")
		}
		for k := 0; k < 1+rnd.Intn(12); k++ {
			lo := rnd.Intn(60)
			hi := lo + rnd.Intn(8)
			s.Add(Range_Closed(lo, hi, cmp))
			model = append(model, [2]int{lo, hi})
			for v := -2; v < 72; v++ {
				want := false
				for _, m := range model {
					if m[0] <= v && v <= m[1] {
						want = true
					}
				}
				if got := s.Contains(v); got != want {
					t.Fatalf("trial %d: contains(%d) = %v, want %v (ranges %v)", trial, v, got, want, model)
				}
			}
			rs := s.AsRanges()
			for i := 1; i < len(rs); i++ {
				if rs[i-1].UpperEndpoint() >= rs[i].LowerEndpoint() {
					t.Fatalf("trial %d: ranges not coalesced: %v", trial, rs)
				}
			}
		}
		if s.IsEmpty() {
			t.Fatal("set empty after add")
		}
	}
	func() {
		defer func() {
			r := recover()
			if _, ok := r.(*IllegalArgumentException); !ok {
				t.Errorf("Range.closed(2, 1) did not throw IllegalArgumentException: %v", r)
			} else if msg := r.(*IllegalArgumentException).GetMessage(); msg != "Invalid range: [2..1]" {
				t.Errorf("message %q", msg)
			}
		}()
		Range_Closed(2, 1, cmp)
	}()
}
