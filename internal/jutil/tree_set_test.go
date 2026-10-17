package jutil

import (
	"strconv"
	"strings"
	"testing"
)

// TestTreeSetMatchesJDK replays testdata/tree_set.tsv
// (tools/oracle/java/JutilTreeSetVectors.java): TreeSet scripts with a
// comparator that changes over time ("v" rows, which expose the shape of the
// red-black tree) and with Integer order ("n" rows).
func TestTreeSetMatchesJDK(t *testing.T) {
	rows := readTSV(t, "tree_set.tsv")
	if len(rows) < 400 {
		t.Fatalf("only %d vectors", len(rows))
	}
	for ri, r := range rows {
		v := r[0] == "v"
		var prio []int
		if v {
			for _, p := range strings.Split(r[2], ",") {
				n, _ := strconv.Atoi(p)
				prio = append(prio, n)
			}
		}
		order := map[int]int{}
		compare := func(a, b int) int {
			if a < b {
				return -1
			} else if a > b {
				return 1
			}
			return 0
		}
		if v {
			compare = func(a, b int) int {
				if d := prio[a] - prio[b]; d != 0 {
					return d
				}
				return order[b] - order[a]
			}
		}
		set := NewTreeSet(compare)
		tf := func(b bool) string {
			if b {
				return "t"
			}
			return "f"
		}
		orNSEE := func(f func() int) (s string) {
			defer func() {
				if x := recover(); x != nil {
					if _, ok := x.(INoSuchElementException); !ok {
						panic(x)
					}
					s = "NSEE"
				}
			}()
			return strconv.Itoa(f())
		}
		adds := 0
		var res []string
		for _, op := range strings.Fields(r[3]) {
			g, _ := strconv.Atoi(op[1:])
			switch op[0] {
			case 'a':
				if v {
					adds++
					order[g] = adds
					removed := set.Remove(g)
					res = append(res, tf(removed)+tf(set.Add(g)))
				} else {
					res = append(res, tf(set.Add(g)))
				}
			case 'r':
				res = append(res, tf(set.Remove(g)))
			case 'c':
				res = append(res, tf(set.Contains(g)))
			case 'f':
				res = append(res, orNSEE(set.First))
			case 'l':
				res = append(res, orNSEE(set.Last))
			case 'x':
				set.Clear()
				res = append(res, "-")
			default:
				t.Fatalf("bad op %q", op)
			}
		}
		join := func(xs []int) string {
			var sb strings.Builder
			for i, x := range xs {
				if i > 0 {
					sb.WriteByte(',')
				}
				sb.WriteString(strconv.Itoa(x))
			}
			return sb.String()
		}
		var desc []int
		for x := range set.Descending() {
			desc = append(desc, x)
		}
		got := []string{strings.Join(res, " "), join(set.ToSlice()), join(desc), strconv.Itoa(set.Size())}
		if want := r[4:8]; strings.Join(got, "\t") != strings.Join(want, "\t") {
			t.Errorf("row %d (%s %s):\n got %q\nwant %q", ri+1, r[0], r[3], got, want)
		}
	}
}

func TestTreeSetString(t *testing.T) {
	s := NewTreeSet(StringCompareTo)
	if s.String() != "[]" {
		t.Errorf("empty %q", s.String())
	}
	for _, e := range []string{"b", "c", "a"} {
		s.Add(e)
	}
	if s.String() != "[a, b, c]" || s.First() != "a" || s.Last() != "c" {
		t.Errorf("got %q first %q last %q", s.String(), s.First(), s.Last())
	}
	calls := 0
	o := NewTreeSet(func(a, b int) int { calls++; return a - b })
	o.Contains(1)
	o.Remove(1)
	if calls != 0 {
		t.Errorf("lookups in an empty set compared %d times", calls)
	}
	o.Add(5)
	if calls != 1 {
		t.Errorf("first add compared %d times, want 1 (compare(e, e))", calls)
	}
}
