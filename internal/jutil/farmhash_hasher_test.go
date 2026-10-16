package jutil

import (
	"encoding/hex"
	"strconv"
	"strings"
	"testing"
)

// Vectors generated with tools/oracle/java/JutilHasherVectors.java against
// the pinned closure-compiler.jar (shaded Guava).
func TestFarmHashFingerprint64HasherMatchesGuava(t *testing.T) {
	rows := readTSV(t, "farmhash_hasher.tsv")
	if len(rows) < 400 {
		t.Fatalf("only %d vectors", len(rows))
	}
	for _, r := range rows {
		want, err := strconv.ParseInt(r[1], 10, 64)
		if err != nil {
			t.Fatal(err)
		}
		h := Hashing_FarmHashFingerprint64NewHasher()
		for _, op := range strings.Split(r[0], ";") {
			arg := op[1:]
			switch op[0] {
			case 'I':
				v, _ := strconv.ParseInt(arg, 10, 32)
				h.PutInt(int32(v))
			case 'L':
				v, _ := strconv.ParseInt(arg, 10, 64)
				h.PutLong(v)
			case 'B':
				b, _ := hex.DecodeString(arg)
				h.PutBytes(b)
			case 'S':
				units := make([]uint16, len(arg)/4)
				for i := range units {
					v, _ := strconv.ParseUint(arg[4*i:4*i+4], 16, 16)
					units[i] = uint16(v)
				}
				h.PutStringUTF8(StringFromUTF16(units))
			default:
				t.Fatalf("bad op %q", op)
			}
		}
		if got := h.HashAsLong(); got != want {
			t.Errorf("hash(%s) = %d, Guava %d", r[0], got, want)
		}
	}
}
