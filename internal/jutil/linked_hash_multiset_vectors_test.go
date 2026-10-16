package jutil

import (
	"strconv"
	"strings"
	"testing"
)

// Vectors generated with tools/oracle/java/JutilLinkedHashMultisetVectors.java
// against the pinned closure-compiler.jar (shaded Guava LinkedHashMultiset).

func TestLinkedHashMultisetMatchesGuava(t *testing.T) {
	rows := readTSV(t, "linked_hash_multiset.tsv")
	if len(rows) < 400 {
		t.Fatalf("only %d vectors", len(rows))
	}
	for _, r := range rows {
		if r[0] == "" {
			continue
		}
		ops := strings.Split(r[0], ";")
		want := strings.Split(r[1], "|")
		m := LinkedHashMultiset_Create[string]()
		for i, op := range ops {
			w := op[1:]
			var res bool
			if op[0] == 'a' {
				res = m.Add(w)
			} else {
				res = m.Remove(w)
			}
			var elems []string
			for e := range m.All() {
				elems = append(elems, e)
			}
			got := strconv.FormatBool(res) + "," + strconv.Itoa(m.Size()) + "," +
				strconv.Itoa(m.Count(w)) + "," + strconv.FormatBool(m.Contains(w)) + "," +
				m.String() + "," + strings.Join(elems, ":")
			if got != want[i] {
				t.Fatalf("%s: op %d (%s): got %q, want %q", r[0], i, op, got, want[i])
			}
		}
	}
}
