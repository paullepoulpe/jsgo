package jutil

import (
	"math/rand/v2"
	"strings"
	"testing"
)

func refStringCompareToUTF16(a, b string) int {
	ua, ub := StringToUTF16(a), StringToUTF16(b)
	m := min(len(ua), len(ub))
	for k := 0; k < m; k++ {
		if ua[k] != ub[k] {
			return int(ua[k]) - int(ub[k])
		}
	}
	return len(ua) - len(ub)
}

func TestStringCompareToZeroAllocAndEquivalence(t *testing.T) {
	hi := StringFromChar(0xD83D)
	lo := StringFromChar(0xDE00)
	nonCanonSmile := hi + lo // 6-byte WTF-8 pair equal in UTF-16 to 4-byte "😀"

	alphabet := []string{
		"",
		"a",
		"b",
		"abc",
		"é",
		"éb",
		"éa",
		"日本語",
		"日本",
		"\uFFFD",
		"\uFFFF",
		"😀",
		"😁",
		"😀x",
		nonCanonSmile,
		nonCanonSmile + "x",
		hi,
		lo,
		hi + "a",
		"\x80",
		"\x80\x80\x80\x80a",
		"\x80\x80\x80\x80b",
		"\xe0\xa0X",
		"\xe0\xa0\x80",
		"\xf0\x9f\x98",
		"\xf0\x9f\x98\x80",
	}

	rng := rand.New(rand.NewPCG(42, 1337))
	samples := make([]string, 0, 400)
	samples = append(samples, alphabet...)
	for range 350 {
		var sb strings.Builder
		parts := rng.IntN(6)
		for range parts {
			sb.WriteString(alphabet[rng.IntN(len(alphabet))])
		}
		samples = append(samples, sb.String())
	}

	for i, a := range samples {
		for j := i; j < len(samples); j++ {
			b := samples[j]
			gotAB := StringCompareTo(a, b)
			wantAB := refStringCompareToUTF16(a, b)
			if gotAB != wantAB {
				t.Fatalf("StringCompareTo(%x, %x) = %d, want %d", a, b, gotAB, wantAB)
			}
			gotBA := StringCompareTo(b, a)
			if gotBA != -wantAB {
				t.Fatalf("StringCompareTo(%x, %x) = %d, want %d", b, a, gotBA, -wantAB)
			}
		}
	}

	a := "prefix_😀_日本語_café_alpha"
	b := "prefix_😀_日本語_café_beta"
	c := "prefix_😀_日本語_\uFFFD"
	allocs := testing.AllocsPerRun(1000, func() {
		if StringCompareTo(a, b) >= 0 {
			t.Fatal("unexpected order")
		}
		if StringCompareTo(a, c) >= 0 {
			t.Fatal("expected astral surrogate hi (0xD83D) < 0xFFFD")
		}
		if StringCompareTo("😀", nonCanonSmile) != 0 {
			t.Fatal("expected UTF-16 equality for non-canonical surrogate pair")
		}
	})
	if allocs != 0 {
		t.Fatalf("StringCompareTo allocated %.2f objects/op, want 0", allocs)
	}
}

func BenchmarkStringCompareToUnicode(b *testing.B) {
	keys := []string{
		" emoji_😀_key_1",
		" emoji_😁_key_2",
		" cjk_日本語_key_3",
		" latin_café_key_4",
		" bmp_\uFFFD_key_5",
	}
	b.ReportAllocs()
	for b.Loop() {
		for i := range keys {
			for j := range keys {
				_ = StringCompareTo(keys[i], keys[j])
			}
		}
	}
}
