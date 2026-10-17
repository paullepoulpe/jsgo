package jutil

import (
	"iter"
	"slices"
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

func jhParseList[K any](t *testing.T, p jhPool[K], s string) []K {
	var out []K
	for _, i := range jhSplit(s) {
		out = append(out, p.fresh(jhAtoi(t, i)))
	}
	return out
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

func jhSeqDump[K any](p jhPool[K], s iter.Seq[K]) string {
	var sb strings.Builder
	for k := range s {
		if sb.Len() > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(strconv.Itoa(p.index(k)))
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

// ---------------------------------------------------------------------------
// HashSet

func TestJavaHashSetMatchesJDK(t *testing.T) {
	for si, sc := range jhScenarios(t, "java_hash_set.tsv", 80) {
		jhWithPool(t, sc[0][1], sc[0][2], func(p any) {
			switch p := p.(type) {
			case jhPool[string]:
				jhReplaySet(t, si, sc, p)
			case jhPool[int32]:
				jhReplaySet(t, si, sc, p)
			default:
				t.Fatalf("set scenario %d: unexpected key type", si)
			}
		})
	}
}

func jhReplaySet[K comparable](t *testing.T, si int, sc [][]string, p jhPool[K]) {
	var s *JavaHashSet[K]
	key := func(x string) K { return p.fresh(jhAtoi(t, x)) }
	for li, r := range sc[1:] {
		ctx := "set scenario " + strconv.Itoa(si) + " line " + strconv.Itoa(li+2) + " " + strings.Join(r, " ")
		if len(ctx) > 300 {
			ctx = ctx[:300]
		}
		want := r[len(r)-1]
		got := want
		switch r[0] {
		case "C":
			switch r[1] {
			case "d":
				s = NewJavaHashSet[K]()
			case "c":
				s = NewJavaHashSetWithCapacity[K](jhAtoi(t, r[2]))
			case "cl":
				lf, _ := strconv.ParseFloat(r[3], 32)
				s = NewJavaHashSetWithCapacityAndLoadFactor[K](jhAtoi(t, r[2]), float32(lf))
			case "col":
				s = NewJavaHashSetFrom[K](NewArrayListFrom(jhParseList(t, p, r[2])))
			case "g0":
				s = Sets_NewHashSet[K]()
			case "gv":
				s = Sets_NewHashSetOf(jhParseList(t, p, r[2])...)
			case "gi":
				s = Sets_NewHashSetFromIterator(slices.Values(jhParseList(t, p, r[2])))
			case "gc":
				s = Sets_NewHashSetFromCollection[K](NewArrayListFrom(jhParseList(t, p, r[2])))
			case "ge":
				s = Sets_NewHashSetWithExpectedSize[K](jhAtoi(t, r[2]))
			case "nh":
				s = JavaHashSet_NewHashSet[K](jhAtoi(t, r[2]))
			default:
				t.Fatalf("%s: bad constructor", ctx)
			}
		case "a":
			got = jhB(s.Add(key(r[1])))
		case "r":
			got = jhB(s.Remove(key(r[1])))
		case "h":
			got = jhB(s.Contains(key(r[1])))
		case "x":
			s.Clear()
		case "aa":
			got = jhB(s.AddAll(slices.Values(jhParseList(t, p, r[1]))))
		case "ra":
			got = jhB(s.RemoveAll(NewArrayListFrom(jhParseList(t, p, r[1]))))
		case "ta":
			got = jhB(s.RetainAll(NewArrayListFrom(jhParseList(t, p, r[1]))))
		case "ir":
			mod, rem := jhAtoi(t, r[1]), jhAtoi(t, r[2])
			var vis []string
			it := s.Iterator()
			for i := 0; it.HasNext(); i++ {
				vis = append(vis, strconv.Itoa(p.index(it.Next())))
				if i%mod == rem {
					it.Remove()
				}
			}
			got = strings.Join(vis, ",")
		case "cpy":
			s = NewJavaHashSetFrom[K](s)
		case "gcp":
			s = Sets_NewHashSetFromCollection[K](s)
		case "D":
			got = jhSeqDump(p, s.All())
		case "E":
			if jhAtoi(t, r[1]) != s.Size() || jhParseInt32(t, r[2]) != s.HashCode() {
				t.Errorf("%s: size %d hashCode %d", ctx, s.Size(), s.HashCode())
			}
			if want != "-" {
				got = multimapVectorsEsc(s.String())
			}
		default:
			t.Fatalf("%s: bad op", ctx)
		}
		if got != want {
			t.Fatalf("%s:\n got %s\nwant %s", ctx, got, want)
		}
		jhCheck(t, ctx, s.m)
	}
}

// ---------------------------------------------------------------------------
// HashMultimap

// jhSrcMultimap is a MultimapSource over entries in order (the generator's
// LinkedHashMultimap sources).
type jhSrcMultimap[K any, V any] struct {
	keys []K
	vals []V
	nKey int
}

func (m *jhSrcMultimap[K, V]) KeySetSize() int { return m.nKey }
func (m *jhSrcMultimap[K, V]) Entries() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for i, k := range m.keys {
			if !yield(k, m.vals[i]) {
				return
			}
		}
	}
}

func TestHashMultimapMatchesGuava(t *testing.T) {
	for si, sc := range jhScenarios(t, "java_hash_multimap.tsv", 80) {
		b := sc[0]
		jhWithPool(t, b[1], b[2], func(kp any) {
			jhWithPool(t, b[3], b[4], func(vp any) {
				switch kp := kp.(type) {
				case jhPool[string]:
					switch vp := vp.(type) {
					case jhPool[string]:
						jhReplayMultimap(t, si, sc, kp, vp)
					case jhPool[int32]:
						jhReplayMultimap(t, si, sc, kp, vp)
					}
				case jhPool[int32]:
					switch vp := vp.(type) {
					case jhPool[string]:
						jhReplayMultimap(t, si, sc, kp, vp)
					case jhPool[int32]:
						jhReplayMultimap(t, si, sc, kp, vp)
					}
				}
			})
		})
	}
}

func jhReplayMultimap[K comparable, V any](t *testing.T, si int, sc [][]string, kp jhPool[K], vp jhPool[V]) {
	var mm *HashMultimap[K, V]
	key := func(x string) K { return kp.fresh(jhAtoi(t, x)) }
	val := func(x string) V { return vp.fresh(jhAtoi(t, x)) }
	for li, r := range sc[1:] {
		ctx := "multimap scenario " + strconv.Itoa(si) + " line " + strconv.Itoa(li+2) + " " + strings.Join(r, " ")
		if len(ctx) > 300 {
			ctx = ctx[:300]
		}
		want := r[len(r)-1]
		got := want
		switch r[0] {
		case "C":
			switch r[1] {
			case "hc":
				mm = HashMultimap_Create[K, V]()
			case "hk":
				mm = HashMultimap_CreateWithExpectedSizes[K, V](jhAtoi(t, r[2]), jhAtoi(t, r[3]))
			case "mb":
				mm = MultimapBuilder_HashKeys_HashSetValues[K, V]()
			case "mbn":
				mm = MultimapBuilder_HashKeysN_HashSetValuesN[K, V](jhAtoi(t, r[2]), jhAtoi(t, r[3]))
			case "hm":
				src := &jhSrcMultimap[K, V]{}
				seen := map[K]bool{}
				for _, e := range jhSplit(r[2]) {
					k, v, _ := strings.Cut(e, ":")
					src.keys = append(src.keys, key(k))
					src.vals = append(src.vals, val(v))
					if !seen[key(k)] {
						seen[key(k)] = true
						src.nKey++
					}
				}
				mm = HashMultimap_CreateFrom[K, V](src)
			default:
				t.Fatalf("%s: bad constructor", ctx)
			}
		case "p":
			got = jhB(mm.Put(key(r[1]), val(r[2])))
		case "r":
			got = jhB(mm.Remove(key(r[1]), val(r[2])))
		case "ra":
			got = jhSeqDump(vp, mm.RemoveAll(key(r[1])).All())
		case "pa":
			got = jhB(mm.PutAll(key(r[1]), slices.Values(jhParseList(t, vp, r[2]))))
		case "ce":
			got = jhB(mm.ContainsEntry(key(r[1]), val(r[2])))
		case "ck":
			got = jhB(mm.ContainsKey(key(r[1])))
		case "g":
			got = jhSeqDump(vp, mm.Get(key(r[1])).All())
		case "x":
			mm.Clear()
		case "D":
			var sb strings.Builder
			var keys, distinct, vals []string
			for k, v := range mm.Entries() {
				if sb.Len() > 0 {
					sb.WriteByte(',')
				}
				ks := strconv.Itoa(kp.index(k))
				sb.WriteString(ks + ":" + strconv.Itoa(vp.index(v)))
				keys = append(keys, ks)
				if len(distinct) == 0 || distinct[len(distinct)-1] != ks {
					distinct = append(distinct, ks)
				}
				vals = append(vals, strconv.Itoa(vp.index(v)))
			}
			got = sb.String()
			if g := jhSeqDump(kp, mm.KeySet()); g != strings.Join(distinct, ",") {
				t.Fatalf("%s: keySet %s", ctx, g)
			}
			if g := jhSeqDump(kp, mm.AsMap().Keys()); g != strings.Join(distinct, ",") {
				t.Fatalf("%s: asMap %s", ctx, g)
			}
			if g := jhSeqDump(kp, mm.Keys()); g != strings.Join(keys, ",") {
				t.Fatalf("%s: keys %s", ctx, g)
			}
			if g := jhSeqDump(vp, mm.Values()); g != strings.Join(vals, ",") {
				t.Fatalf("%s: values %s", ctx, g)
			}
			if mm.Size() != len(keys) || mm.KeySetSize() != len(distinct) {
				t.Fatalf("%s: size %d keySetSize %d", ctx, mm.Size(), mm.KeySetSize())
			}
		case "E":
			if jhAtoi(t, r[1]) != mm.Size() || jhParseInt32(t, r[2]) != mm.HashCode() {
				t.Errorf("%s: size %d hashCode %d", ctx, mm.Size(), mm.HashCode())
			}
			if want != "-" {
				got = multimapVectorsEsc(mm.String())
			}
		default:
			t.Fatalf("%s: bad op", ctx)
		}
		if got != want {
			t.Fatalf("%s:\n got %s\nwant %s", ctx, got, want)
		}
		jhCheck(t, ctx, mm.m)
	}
}
