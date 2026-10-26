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
		case "pq", "ms":
			rows++
			continue
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
