package jutil

import (
	"strconv"
	"strings"
	"testing"
)

type pqTestElem struct{ key, id int }

// Vectors generated with tools/oracle/java/JutilPriorityQueueVectors.java
// (JDK 21 java.util.PriorityQueue with a key-only comparator, so that ties
// expose the heap's sift order).
func TestPriorityQueueMatchesJDK(t *testing.T) {
	rows := readTSV(t, "priority_queue.tsv")
	if len(rows) < 600 {
		t.Fatalf("only %d vectors", len(rows))
	}
	for _, r := range rows {
		capacity, _ := strconv.Atoi(r[0])
		q := NewPriorityQueueWithCapacity(capacity, func(a, b *pqTestElem) int {
			return a.key - b.key
		})
		var all []*pqTestElem
		var res []string
		for _, op := range strings.Split(r[1], ";") {
			arg, _ := strconv.Atoi(op[1:])
			switch op[0] {
			case 'A':
				e := &pqTestElem{key: arg, id: len(all)}
				all = append(all, e)
				if q.Add(e) {
					res = append(res, "t")
				} else {
					res = append(res, "f")
				}
			case 'P':
				if e := q.Poll(); e == nil {
					res = append(res, "null")
				} else {
					res = append(res, strconv.Itoa(e.id))
				}
			case 'R':
				func() {
					defer func() {
						if x := recover(); x != nil {
							if _, ok := x.(INoSuchElementException); !ok {
								panic(x)
							}
							res = append(res, "NSEE")
						}
					}()
					res = append(res, strconv.Itoa(q.Remove().id))
				}()
			case 'K':
				if e := q.Peek(); e == nil {
					res = append(res, "null")
				} else {
					res = append(res, strconv.Itoa(e.id))
				}
			case 'D':
				if len(all) > 0 && q.RemoveObject(all[arg]) {
					res = append(res, "t")
				} else {
					res = append(res, "f")
				}
			case 'C':
				if len(all) > 0 && q.Contains(all[arg]) {
					res = append(res, "t")
				} else {
					res = append(res, "f")
				}
			case 'S':
				res = append(res, strconv.Itoa(q.Size()))
			case 'X':
				q.Clear()
				res = append(res, "")
			default:
				t.Fatalf("bad op %q", op)
			}
		}
		var order []string
		for e := range q.All() {
			order = append(order, strconv.Itoa(e.id))
		}
		want := ""
		if len(r) > 3 {
			want = r[3]
		}
		if got := strings.Join(res, ";"); got != r[2] {
			t.Errorf("%s: results %s, JDK %s", r[1], got, r[2])
		}
		if got := strings.Join(order, ","); got != want {
			t.Errorf("%s: order %s, JDK %s", r[1], got, want)
		}
	}
}

func TestPriorityQueueRejectsNull(t *testing.T) {
	q := NewPriorityQueue(func(a, b *pqTestElem) int { return a.key - b.key })
	defer func() {
		if _, ok := recover().(*NullPointerException); !ok {
			t.Fatal("want NullPointerException")
		}
	}()
	q.Add(nil)
}
