package jutil

import (
	"strconv"
	"strings"
	"testing"
)

// TestBitSetFlipMatchesJDK replays testdata/bitset_flip.tsv
// (tools/oracle/java/JutilBitSetFlipVectors.java).
func TestBitSetFlipMatchesJDK(t *testing.T) {
	for _, r := range readTSV(t, "bitset_flip.tsv") {
		b := NewBitSet()
		var got strings.Builder
		for _, op := range strings.Split(strings.TrimSuffix(r[0], ","), ",") {
			res := ""
			Try(func() {
				switch op[0] {
				case 'c':
					n, _ := strconv.Atoi(op[1:])
					b.Clear(n)
				case 's':
					n, _ := strconv.Atoi(op[1:])
					b.Set(n)
				case 'f':
					n, _ := strconv.Atoi(op[1:])
					b.Flip(n)
				case 'F':
					parts := strings.Split(op[1:], ":")
					from, _ := strconv.Atoi(parts[0])
					to, _ := strconv.Atoi(parts[1])
					b.FlipRange(from, to)
				}
				res = b.String() + "/" + strconv.Itoa(b.Length()) + "/" + strconv.Itoa(b.Cardinality())
			}, func(e IIndexOutOfBoundsException) { res = "!" + e.GetMessage() })
			got.WriteString(res + ";")
		}
		if got.String() != r[1] {
			t.Errorf("ops %s:\n got %s\nwant %s", r[0], got.String(), r[1])
		}
	}
}
