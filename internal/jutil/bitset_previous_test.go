package jutil

import (
	"strconv"
	"strings"
	"testing"
)

// TestBitSetPreviousSetBitMatchesJDK replays testdata/bitset_previous.tsv
// (tools/oracle/java/JutilBitSetPreviousVectors.java).
func TestBitSetPreviousSetBitMatchesJDK(t *testing.T) {
	for _, r := range readTSV(t, "bitset_previous.tsv") {
		b := NewBitSet()
		for _, op := range strings.Split(strings.TrimSuffix(r[0], ","), ",") {
			if op == "" {
				continue
			}
			n, _ := strconv.Atoi(op[1:])
			if op[0] == 's' {
				b.Set(n)
			} else {
				b.Clear(n)
			}
		}
		var got strings.Builder
		for _, q := range strings.Split(strings.TrimSuffix(r[1], ";"), ";") {
			from, _ := strconv.Atoi(q[:strings.IndexByte(q, '=')])
			res := ""
			Try(func() { res = strconv.Itoa(b.PreviousSetBit(from)) },
				func(e IIndexOutOfBoundsException) { res = "!" + e.GetMessage() })
			got.WriteString(strconv.Itoa(from) + "=" + res + ";")
		}
		if got.String() != r[1] {
			t.Errorf("ops %s:\n got %s\nwant %s", r[0], got.String(), r[1])
		}
	}
}
