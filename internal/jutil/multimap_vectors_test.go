package jutil

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

// Vectors generated with tools/oracle/java/JutilMultimapVectors.java against
// the pinned closure-compiler.jar (shaded Guava).

func multimapVectorsUnhex(t *testing.T, s string) string {
	units := make([]uint16, len(s)/4)
	for i := range units {
		v, err := strconv.ParseUint(s[4*i:4*i+4], 16, 16)
		if err != nil {
			t.Fatal(err)
		}
		units[i] = uint16(v)
	}
	return StringFromUTF16(units)
}

func multimapVectorsEsc(s string) string {
	var sb strings.Builder
	for _, c := range StringToUTF16(s) {
		if c < 0x20 || c >= 0x7f || c == '\\' {
			fmt.Fprintf(&sb, "\\u%04x", c)
		} else {
			sb.WriteByte(byte(c))
		}
	}
	return sb.String()
}

func multimapVectorsPuts(t *testing.T, s string, put func(k, v string)) {
	if s == "" {
		return
	}
	for _, p := range strings.Split(s, ";") {
		kv := strings.SplitN(p, "=", 2)
		put(multimapVectorsUnhex(t, kv[0]), multimapVectorsUnhex(t, kv[1]))
	}
}

func TestMultimapsMatchGuava(t *testing.T) {
	rows := readTSV(t, "multimap.tsv")
	if len(rows) < 600 {
		t.Fatalf("only %d vectors", len(rows))
	}
	for _, r := range rows {
		switch r[0] {
		case "set":
			mm := NewSetMultimap[string, string]()
			multimapVectorsPuts(t, r[1], func(k, v string) { mm.Put(k, v) })
			m := SetMultimapCopyOf(mm)
			if got := multimapVectorsEsc(m.String()); got != r[2] {
				t.Errorf("set %s = %s, Guava %s", r[1], got, r[2])
			}
			if got := multimapVectorsEsc(m.Inverse().String()); got != r[3] {
				t.Errorf("set %s inverse = %s, Guava %s", r[1], got, r[3])
			}
		case "list":
			continue
		default:
			t.Fatalf("bad row %q", r[0])
		}
	}
}
