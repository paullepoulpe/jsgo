package jutil

import (
	"strconv"
	"strings"
	"testing"
)

// TestArrayDequeMatchesJDK replays testdata/array_deque.tsv
// (tools/oracle/java/JutilArrayDequeVectors.java).
func TestArrayDequeMatchesJDK(t *testing.T) {
	str := func(p *int) string {
		if p == nil {
			return "null"
		}
		return strconv.Itoa(*p)
	}
	for _, r := range readTSV(t, "array_deque.tsv") {
		d := NewArrayDeque[*int]()
		var got strings.Builder
		for _, op := range strings.Split(strings.TrimSuffix(r[0], ","), ",") {
			name := strings.TrimRight(op, "0123456789:")
			arg := op[len(name):]
			v, _ := strconv.Atoi(arg)
			res := ""
			Try(func() {
				switch name {
				case "push":
					d.Push(&v)
					res = "-"
				case "pop":
					res = str(d.Pop())
				case "peek":
					res = str(d.Peek())
				case "add":
					res = strconv.FormatBool(d.Add(&v))
				case "addFirst":
					d.AddFirst(&v)
					res = "-"
				case "addLast":
					d.AddLast(&v)
					res = "-"
				case "poll":
					res = str(d.Poll())
				case "peekLast":
					res = str(d.PeekLast())
				case "pollLast":
					res = str(d.PollLast())
				case "removeFirst":
					res = str(d.RemoveFirst())
				case "removeLast":
					res = str(d.RemoveLast())
				case "getFirst":
					res = str(d.GetFirst())
				case "getLast":
					res = str(d.GetLast())
				case "copy":
					d = NewArrayDequeFrom(d.All())
					res = "-"
				case "addAll":
					var l []*int
					if arg != "" {
						for _, x := range strings.Split(arg, ":") {
							n, _ := strconv.Atoi(x)
							l = append(l, &n)
						}
					}
					res = strconv.FormatBool(d.AddAll(func(yield func(*int) bool) {
						for _, e := range l {
							if !yield(e) {
								return
							}
						}
					}))
				case "clear":
					d.Clear()
					res = "-"
				case "contains":
					found := false
					for e := range d.All() {
						if *e == v {
							found = true
						}
					}
					res = strconv.FormatBool(found)
				case "desc":
					var parts []string
					for e := range d.Descending() {
						parts = append(parts, strconv.Itoa(*e))
					}
					res = strings.Join(parts, " ")
				case "pollFirst":
					res = str(d.PollFirst())
				default:
					t.Fatalf("unknown op %q", op)
				}
			}, func(e INoSuchElementException) { res = "!NoSuchElementException" })
			var parts []string
			for e := range d.All() {
				parts = append(parts, strconv.Itoa(*e))
			}
			got.WriteString(res + "/" + strconv.Itoa(d.Size()) + "/[" + strings.Join(parts, ", ") + "];")
		}
		if got.String() != r[1] {
			t.Errorf("ops %s:\n got %s\nwant %s", r[0], got.String(), r[1])
		}
	}
}
