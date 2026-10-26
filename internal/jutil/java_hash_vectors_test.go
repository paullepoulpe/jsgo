package jutil

import (
	"iter"
	"strconv"
	"strings"
	"testing"
)

// Vectors generated with tools/oracle/java/JutilJavaHashVectors.java against
// the JDK and the shaded Guava of the pinned closure-compiler.jar. The
// generator takes a seed; this test only replays what the files contain, and
// checks coverage (tree bins built, split and dissolved) from the replay.

// jhCP is the generator's CP: a Comparable value class with a colliding
// hashCode.
type jhCP struct{ a, b int32 }

func (c *jhCP) Equals(o any) bool {
	x, ok := o.(*jhCP)
	return ok && x != nil && x.a == c.a && x.b == c.b
}
func (c *jhCP) HashCode() int32 { return 31*c.a + c.b }
func (c *jhCP) CompareTo(o *jhCP) int {
	if c.a != o.a {
		return javaCompareInts(c.a, o.a)
	}
	return javaCompareInts(c.b, o.b)
}
func (c *jhCP) String() string {
	return "CP(" + strconv.Itoa(int(c.a)) + "," + strconv.Itoa(int(c.b)) + ")"
}

// jhNP is the generator's NP: a non-Comparable value class.
type jhNP struct{ h int32 }

func (n *jhNP) Equals(o any) bool {
	x, ok := o.(*jhNP)
	return ok && x != nil && x.h == n.h
}
func (n *jhNP) HashCode() int32 { return n.h }
func (n *jhNP) String() string  { return "NP(" + strconv.Itoa(int(n.h)) + ")" }

// jhPool is a decoded key pool: fresh(i) returns a new key equal to entry i
// (a new object for value classes), index maps a key back to its entry.
type jhPool[K any] struct {
	n     int
	fresh func(i int) K
	index func(k K) int
}

func jhAtoi(t *testing.T, s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		t.Fatalf("bad int %q", s)
	}
	return n
}

func jhParseInt32(t *testing.T, s string) int32 {
	n, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		t.Fatalf("bad int %q", s)
	}
	return int32(n)
}

func jhParseInt64(t *testing.T, s string) int64 {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		t.Fatalf("bad long %q", s)
	}
	return n
}

func jhSplit(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, ",")
}

func newJHPool[K any](t *testing.T, encs []string, decode func(string) K) jhPool[K] {
	idx := make(map[string]int, len(encs))
	for i, e := range encs {
		idx[e] = i
	}
	return jhPool[K]{
		n:     len(encs),
		fresh: func(i int) K { return decode(encs[i]) },
		index: func(k K) int {
			i, ok := idx[jhEncode(k)]
			if !ok {
				t.Fatalf("key %v not in pool", k)
			}
			return i
		},
	}
}

// jhEncode is the generator's encode() for the Go key types.
func jhEncode(k any) string {
	switch x := k.(type) {
	case string:
		var sb strings.Builder
		for _, c := range StringToUTF16(x) {
			sb.WriteString(strconv.FormatUint(uint64(c)|0x10000, 16)[1:])
		}
		return sb.String()
	case int32:
		return strconv.Itoa(int(x))
	case int64:
		return strconv.FormatInt(x, 10)
	case *jhCP:
		return strconv.Itoa(int(x.a)) + ":" + strconv.Itoa(int(x.b))
	case *jhNP:
		return strconv.Itoa(int(x.h))
	}
	panic("bad key")
}

func jhEncodeMixed(k any) string {
	switch k.(type) {
	case string:
		return "s" + jhEncode(k)
	case int32:
		return "i" + jhEncode(k)
	}
	return "l" + jhEncode(k)
}

// jhWithPool decodes a pool of the given type and calls f with it.
func jhWithPool(t *testing.T, typ, pool string, f func(p any)) {
	encs := jhSplit(pool)
	switch typ {
	case "S":
		f(newJHPool(t, encs, func(s string) string { return multimapVectorsUnhex(t, s) }))
	case "I":
		f(newJHPool(t, encs, func(s string) int32 { return jhParseInt32(t, s) }))
	case "L":
		f(newJHPool(t, encs, func(s string) int64 { return jhParseInt64(t, s) }))
	case "P":
		f(newJHPool(t, encs, func(s string) *jhCP {
			a, b, _ := strings.Cut(s, ":")
			return &jhCP{jhParseInt32(t, a), jhParseInt32(t, b)}
		}))
	case "N":
		f(newJHPool(t, encs, func(s string) *jhNP { return &jhNP{jhParseInt32(t, s)} }))
	case "O":
		idx := make(map[string]int, len(encs))
		for i, e := range encs {
			idx[e] = i
		}
		f(jhPool[any]{
			n: len(encs),
			fresh: func(i int) any {
				e := encs[i]
				switch e[0] {
				case 's':
					return multimapVectorsUnhex(t, e[1:])
				case 'i':
					return jhParseInt32(t, e[1:])
				}
				return jhParseInt64(t, e[1:])
			},
			index: func(k any) int { return idx[jhEncodeMixed(k)] },
		})
	default:
		t.Fatalf("bad key type %q", typ)
	}
}

// jhScenarios splits rows into scenarios, each starting with a B row and
// ending with an E row.
func jhScenarios(t *testing.T, name string, min int) [][][]string {
	rows := readTSV(t, name)
	var out [][][]string
	var cur [][]string
	for _, r := range rows {
		if r[0] == "B" {
			cur = nil
		}
		cur = append(cur, r)
		if r[0] == "E" {
			out = append(out, cur)
		}
	}
	if len(out) < min {
		t.Fatalf("%s: only %d scenarios", name, len(out))
	}
	return out
}

// jhV encodes a map value (an Integer or null) like the generator.
func jhV(v any) string {
	if IsNil(v) {
		return "n"
	}
	return ToJavaString(v)
}

func jhParseV(t *testing.T, s string) any {
	if s == "n" {
		return nil
	}
	return jhParseInt32(t, s)
}

func jhB(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

func jhMergeFn(a, b any) any {
	s := a.(int32) + b.(int32)
	if s%5 == 0 {
		return nil
	}
	return s
}

// jhSrcMap is a MapSource over k=v pairs in order (the generator's
// LinkedHashMap sources).
type jhSrcMap[K any] struct {
	keys []K
	vals []any
}

func (m *jhSrcMap[K]) Size() int { return len(m.keys) }
func (m *jhSrcMap[K]) All() iter.Seq2[K, any] {
	return func(yield func(K, any) bool) {
		for i, k := range m.keys {
			if !yield(k, m.vals[i]) {
				return
			}
		}
	}
}

func jhParseSrc[K any](t *testing.T, p jhPool[K], s string) *jhSrcMap[K] {
	src := &jhSrcMap[K]{}
	for _, kv := range jhSplit(s) {
		k, v, _ := strings.Cut(kv, "=")
		src.keys = append(src.keys, p.fresh(jhAtoi(t, k)))
		src.vals = append(src.vals, jhParseV(t, v))
	}
	return src
}

func jhMapDump[K any](p jhPool[K], m *JavaHashMap[K, any]) string {
	var sb strings.Builder
	for k, v := range m.All() {
		if sb.Len() > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(strconv.Itoa(p.index(k)))
		sb.WriteByte('=')
		sb.WriteString(jhV(v))
	}
	return sb.String()
}

// ---------------------------------------------------------------------------
// Structure checks and coverage

type jhCoverage struct {
	treeStates     int // operations after which some bin is a tree
	treeResizes    int // resizes of a table holding a tree bin
	treeDissolved  int // operations after which there were fewer tree bins
	treeIterRemove int // iterator removals from a tree bin
	treeSplitTwo   int // resizes that split a tree bin into two trees
}

func jhTreeBins[K any, V any](m *JavaHashMap[K, V]) int {
	n := 0
	for _, e := range m.table {
		if e != nil && e.rb != nil {
			n++
		}
	}
	return n
}

// jhCheck verifies the invariants of JavaHashMap's representation: every
// node in its bin, the size, every bin either a list or a tree bin, and for
// tree bins the chain's prev links, the red-black properties, parent links,
// hash order, and that the tree holds exactly the chain's nodes.
func jhCheck[K any, V any](t *testing.T, ctx string, m *JavaHashMap[K, V]) {
	t.Helper()
	count := 0
	n := int32(len(m.table))
	for i, e := range m.table {
		if e == nil {
			continue
		}
		isTree := e.rb != nil
		var prev *hashMapNode[K, V]
		chain := map[*hashMapNode[K, V]]bool{}
		for x := e; x != nil; x = x.next {
			count++
			chain[x] = true
			if x.hash&(n-1) != int32(i) {
				t.Fatalf("%s: node in wrong bin", ctx)
			}
			if (x.rb != nil) != isTree {
				t.Fatalf("%s: mixed tree/list bin", ctx)
			}
			if isTree && x.rb.prev != prev {
				t.Fatalf("%s: bad prev link", ctx)
			}
			prev = x
		}
		if !isTree {
			continue
		}
		root := e.root()
		if root.rb.red {
			t.Fatalf("%s: red root", ctx)
		}
		inTree := 0
		var walk func(x *hashMapNode[K, V]) int
		walk = func(x *hashMapNode[K, V]) int {
			if x == nil {
				return 1
			}
			inTree++
			if !chain[x] {
				t.Fatalf("%s: tree node not in the chain", ctx)
			}
			l, r := x.rb.left, x.rb.right
			if l != nil && (l.rb.parent != x || l.hash > x.hash) {
				t.Fatalf("%s: bad left child", ctx)
			}
			if r != nil && (r.rb.parent != x || r.hash < x.hash) {
				t.Fatalf("%s: bad right child", ctx)
			}
			if x.rb.red && (javaHashMapIsRed(l) || javaHashMapIsRed(r)) {
				t.Fatalf("%s: red node with red child", ctx)
			}
			bl, br := walk(l), walk(r)
			if bl != br {
				t.Fatalf("%s: unequal black heights", ctx)
			}
			if !x.rb.red {
				bl++
			}
			return bl
		}
		walk(root)
		if len(chain) != inTree {
			t.Fatalf("%s: tree has %d nodes, chain %d", ctx, inTree, len(chain))
		}
	}
	if count != m.size {
		t.Fatalf("%s: size %d, %d nodes", ctx, m.size, count)
	}
}

// ---------------------------------------------------------------------------
// HashMap

func TestJavaHashMapMatchesJDK(t *testing.T) {
	var cov jhCoverage
	scenarios := jhScenarios(t, "java_hash_map.tsv", 180)
	// The compute-driven family (JutilJavaHashVectors.computeScenario).
	scenarios = append(scenarios, jhScenarios(t, "java_hash_map_compute.tsv", 90)...)
	for si, sc := range scenarios {
		jhWithPool(t, sc[0][1], sc[0][2], func(p any) {
			switch p := p.(type) {
			case jhPool[string]:
				jhReplayMap(t, si, sc, p, &cov)
			case jhPool[int32]:
				jhReplayMap(t, si, sc, p, &cov)
			case jhPool[int64]:
				jhReplayMap(t, si, sc, p, &cov)
			case jhPool[*jhCP]:
				jhReplayMap(t, si, sc, p, &cov)
			case jhPool[*jhNP]:
				jhReplayMap(t, si, sc, p, &cov)
			case jhPool[any]:
				jhReplayMap(t, si, sc, p, &cov)
			}
		})
	}
	t.Logf("coverage: %+v", cov)
	if cov.treeStates < 200 || cov.treeResizes < 10 || cov.treeSplitTwo < 1 || cov.treeDissolved < 10 || cov.treeIterRemove < 5 {
		t.Errorf("tree bins not covered enough: %+v", cov)
	}
}

func jhReplayMap[K any](t *testing.T, si int, sc [][]string, p jhPool[K], cov *jhCoverage) {
	var m *JavaHashMap[K, any]
	key := func(s string) K { return p.fresh(jhAtoi(t, s)) }
	for li, r := range sc[1:] {
		ctx := "map scenario " + strconv.Itoa(si) + " line " + strconv.Itoa(li+2) + " " + strings.Join(r, " ")
		if len(ctx) > 300 {
			ctx = ctx[:300]
		}
		want := r[len(r)-1]
		got := ""
		treesBefore, capBefore := 0, 0
		if m != nil {
			treesBefore, capBefore = jhTreeBins(m), len(m.table)
		}
		switch r[0] {
		case "C":
			switch r[1] {
			case "d":
				m = NewJavaHashMap[K, any]()
			case "c":
				m = NewJavaHashMapWithCapacity[K, any](jhAtoi(t, r[2]))
			case "cl":
				lf, _ := strconv.ParseFloat(r[3], 32)
				m = NewJavaHashMapWithCapacityAndLoadFactor[K, any](jhAtoi(t, r[2]), float32(lf))
			case "m":
				m = NewJavaHashMapFrom[K, any](jhParseSrc(t, p, r[2]))
			case "gd":
				m = Maps_NewHashMap[K, any]()
			case "gm":
				m = Maps_NewHashMapFrom[K, any](jhParseSrc(t, p, r[2]))
			case "ge":
				m = Maps_NewHashMapWithExpectedSize[K, any](jhAtoi(t, r[2]))
			case "nh":
				m = JavaHashMap_NewHashMap[K, any](jhAtoi(t, r[2]))
			default:
				t.Fatalf("%s: bad constructor", ctx)
			}
			got = want
		case "p":
			got = jhV(m.Put(key(r[1]), jhParseV(t, r[2])))
		case "pi":
			got = jhV(m.PutIfAbsent(key(r[1]), jhParseV(t, r[2])))
		case "r":
			got = jhV(m.Remove(key(r[1])))
		case "re":
			got = jhB(m.RemoveEntry(key(r[1]), jhParseV(t, r[2])))
		case "g":
			got = jhV(m.Get(key(r[1])))
		case "ck":
			got = jhB(m.ContainsKey(key(r[1])))
		case "kr":
			got = jhB(m.KeySet().Remove(key(r[1])))
		case "x":
			m.Clear()
			got = want
		case "ca":
			y := jhParseV(t, r[2])
			got = jhV(m.ComputeIfAbsent(key(r[1]), func(K) any { return y }))
		case "cp":
			y := jhParseV(t, r[2])
			got = jhV(m.ComputeIfPresent(key(r[1]), func(K, any) any { return y }))
		case "co":
			y := jhParseV(t, r[2])
			got = jhV(m.Compute(key(r[1]), func(K, any) any { return y }))
		case "mg":
			got = jhV(m.Merge(key(r[1]), jhParseV(t, r[2]), jhMergeFn))
		case "pa":
			m.PutAll(jhParseSrc(t, p, r[1]))
			got = want
		case "ir":
			mod, rem := jhAtoi(t, r[1]), jhAtoi(t, r[2])
			var vis []string
			it := m.Iterator()
			for i := 0; it.HasNext(); i++ {
				k, _ := it.Next()
				vis = append(vis, strconv.Itoa(p.index(k)))
				if i%mod == rem {
					if it.lastReturned.rb != nil {
						cov.treeIterRemove++
					}
					it.Remove()
				}
			}
			got = strings.Join(vis, ",")
		case "cme":
			k, pos := key(r[1]), jhAtoi(t, r[2])
			got = "ok"
			i := 0
			Try(func() {
				for range m.Keys() {
					if i == pos {
						m.Put(k, int32(-1))
					}
					i++
				}
			}, func(e IRuntimeException) {
				if _, ok := e.(IConcurrentModificationException); !ok {
					panic(e)
				}
				got = "cme" + strconv.Itoa(i)
			})
		case "cpy":
			m = NewJavaHashMapFrom[K, any](m)
			got = want
		case "cln":
			m = m.Clone()
			got = want
		case "D":
			got = jhMapDump(p, m)
		case "E":
			if jhAtoi(t, r[1]) != m.Size() || jhParseInt32(t, r[2]) != m.HashCode() {
				t.Errorf("%s: size %d hashCode %d", ctx, m.Size(), m.HashCode())
			}
			if want != "-" {
				got = multimapVectorsEsc(m.String())
			} else {
				got = want
			}
		default:
			t.Fatalf("%s: bad op", ctx)
		}
		if got != want {
			t.Fatalf("%s:\n got %s\nwant %s", ctx, got, want)
		}
		jhCheck(t, ctx, m)
		trees := jhTreeBins(m)
		if trees > 0 {
			cov.treeStates++
		}
		if treesBefore > 0 && len(m.table) > capBefore && capBefore > 0 {
			cov.treeResizes++
			if trees > treesBefore {
				cov.treeSplitTwo++
			}
		}
		if trees < treesBefore {
			cov.treeDissolved++
		}
	}
}
