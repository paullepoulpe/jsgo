package jutil

import (
	"bufio"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
)

// TestInfraCollectionsMatchJava replays testdata/infra_collections.tsv
// (tools/oracle/java/JutilInfraVectors.java): PriorityQueue with ties,
// HashMultiset, ArrayListMultimap and hash-keyed SetMultimap scripts.
func TestInfraCollectionsMatchJava(t *testing.T) {
	f, err := os.Open("testdata/infra_collections.tsv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1<<20), 1<<20)
	rows := 0
	for sc.Scan() {
		cols := strings.Split(sc.Text(), "\t")
		kind, ops := cols[0], []string{}
		if cols[1] != "" {
			ops = strings.Split(cols[1], ";")
		}
		var got []string
		switch kind {
		case "pq":
			got = replayPriorityQueue(ops)
		case "ms":
			got = replayMultiset(ops)
		case "lm":
			got = replayListMultimap(ops)
		case "sm":
			got = replaySetMultimap(ops)
		default:
			t.Fatalf("unknown kind %q", kind)
		}
		if want := cols[2]; strings.Join(got, ";") != want {
			t.Errorf("%s %s:\n got %s\nwant %s", kind, cols[1], strings.Join(got, ";"), want)
		}
		rows++
	}
	if rows != 650 {
		t.Fatalf("replayed %d rows, want 650", rows)
	}
}

type pqElem struct{ key, id int }

func replayPriorityQueue(ops []string) []string {
	q := NewPriorityQueue(func(a, b *pqElem) int { return a.key - b.key })
	show := func(e *pqElem) string {
		if e == nil {
			return "null"
		}
		return strconv.Itoa(e.key) + "." + strconv.Itoa(e.id)
	}
	var res []string
	for _, op := range ops {
		switch op[0] {
		case 'a':
			k, id, _ := strings.Cut(op[1:], ".")
			ki, _ := strconv.Atoi(k)
			ii, _ := strconv.Atoi(id)
			res = append(res, strconv.FormatBool(q.Add(&pqElem{ki, ii})))
		case 'p':
			res = append(res, show(q.Poll()))
		case 'k':
			res = append(res, show(q.Peek()))
		}
	}
	res = append(res, "size="+strconv.Itoa(q.Size()))
	for !q.IsEmpty() {
		res = append(res, show(q.Poll()))
	}
	return res
}

func replayMultiset(ops []string) []string {
	m := HashMultiset_Create[string]()
	var res []string
	for _, op := range ops {
		arg := op[1:]
		switch op[0] {
		case 'a':
			res = append(res, strconv.FormatBool(m.Add(arg)))
		case 'r':
			res = append(res, strconv.FormatBool(m.Remove(arg)))
		case 'R':
			var ws []string
			if arg != "" {
				ws = strings.Split(arg, ",")
			}
			res = append(res, strconv.FormatBool(m.RemoveAll(ws)))
		case 'c':
			res = append(res, strconv.Itoa(m.Count(arg)))
		case 'h':
			res = append(res, strconv.FormatBool(m.Contains(arg)))
		}
	}
	return append(res, "size="+strconv.Itoa(m.Size()))
}

func replayListMultimap(ops []string) []string {
	m := ArrayListMultimap_Create[string, string]()
	var res []string
	for _, op := range ops {
		arg := op[1:]
		switch op[0] {
		case 'p':
			k, v, _ := strings.Cut(arg, "=")
			res = append(res, strconv.FormatBool(m.Put(k, v)))
		case 'R':
			res = append(res, CollectionToString(m.RemoveAll(arg)))
		case 'g':
			res = append(res, CollectionToString(m.Get(arg)))
		case 'c':
			res = append(res, strconv.FormatBool(m.ContainsKey(arg)))
		}
	}
	return append(res, "size="+strconv.Itoa(m.Size()))
}

func replaySetMultimap(ops []string) []string {
	m := NewSetMultimap[string, string]()
	sorted := func(s *ObjectLinkedHashSet[string]) string {
		l := slices.Collect(s.All())
		slices.Sort(l)
		return CollectionToString(l)
	}
	var res []string
	for _, op := range ops {
		arg := op[1:]
		switch op[0] {
		case 'p':
			k, v, _ := strings.Cut(arg, "=")
			res = append(res, strconv.FormatBool(m.Put(k, v)))
		case 'P':
			k, v, _ := strings.Cut(arg, "=")
			var vs []string
			if v != "" {
				vs = strings.Split(v, ",")
			}
			res = append(res, strconv.FormatBool(m.PutAll(k, slices.Values(vs))))
		case 'R':
			res = append(res, sorted(m.RemoveAll(arg)))
		case 'g':
			res = append(res, sorted(m.Get(arg)))
		case 'c':
			res = append(res, strconv.FormatBool(m.ContainsKey(arg)))
		}
	}
	return append(res, "size="+strconv.Itoa(m.Size()))
}
