package jutil

import (
	"math"
	"math/rand"
	"slices"
	"testing"
)

// orderedModel is a reference insertion-ordered map: a slice searched
// linearly, with Go == on keys like LinkedHashMap.
type orderedModel[K comparable] struct {
	keys []K
	vals []int
}

func (m *orderedModel[K]) index(k K) int {
	for i, x := range m.keys {
		if x == k {
			return i
		}
	}
	return -1
}

func (m *orderedModel[K]) put(k K, v int) (int, bool) {
	if i := m.index(k); i >= 0 {
		old := m.vals[i]
		m.vals[i] = v
		return old, true
	}
	m.keys = append(m.keys, k)
	m.vals = append(m.vals, v)
	return 0, false
}

func (m *orderedModel[K]) remove(k K) (int, bool) {
	i := m.index(k)
	if i < 0 {
		return 0, false
	}
	v := m.vals[i]
	m.keys = slices.Delete(m.keys, i, i+1)
	m.vals = slices.Delete(m.vals, i, i+1)
	return v, true
}

// checkLinkedHashMapAgainstModel runs random operations on a LinkedHashMap
// and a LinkedHashSet over keys drawn from keyOf(0..domain-1), so that the
// maps grow past lhmSmallSize, shrink back and are cleared, and compares
// every result and the iteration order with the model.
func checkLinkedHashMapAgainstModel[K comparable](t *testing.T, seed int64, domain int, keyOf func(int) K) {
	r := rand.New(rand.NewSource(seed))
	h := NewLinkedHashMap[K, int]()
	s := NewLinkedHashSet[K]()
	var m orderedModel[K]
	for step := 0; step < 20000; step++ {
		k := keyOf(r.Intn(domain))
		switch op := r.Intn(100); {
		case op < 45:
			v := r.Intn(1000) + 1
			gotOld := h.Put(k, v)
			wantOld, present := m.put(k, v)
			if gotOld != wantOld {
				t.Fatalf("step %d: Put(%v) returned %d, want %d", step, k, gotOld, wantOld)
			}
			if s.Add(k) == present {
				t.Fatalf("step %d: set Add(%v) = %v", step, k, present)
			}
		case op < 75:
			got := h.Remove(k)
			want, present := m.remove(k)
			if got != want {
				t.Fatalf("step %d: Remove(%v) returned %d, want %d", step, k, got, want)
			}
			if s.Remove(k) != present {
				t.Fatalf("step %d: set Remove(%v) != %v", step, k, present)
			}
		case op < 77:
			h.Clear()
			s.Clear()
			m = orderedModel[K]{}
		default:
			v, ok := h.Lookup(k)
			i := m.index(k)
			if ok != (i >= 0) || (ok && v != m.vals[i]) || h.ContainsKey(k) != ok || h.Get(k) != v {
				t.Fatalf("step %d: Lookup(%v) = %d, %v; model index %d", step, k, v, ok, i)
			}
			if s.Contains(k) != ok {
				t.Fatalf("step %d: set Contains(%v) != %v", step, k, ok)
			}
		}
		if h.Size() != len(m.keys) || s.Size() != len(m.keys) {
			t.Fatalf("step %d: sizes %d, %d, want %d", step, h.Size(), s.Size(), len(m.keys))
		}
		if step%50 == 0 {
			if got := h.KeySlice(); !slices.Equal(got, m.keys) && !(len(got) == 0 && len(m.keys) == 0) {
				t.Fatalf("step %d: keys %v, want %v", step, got, m.keys)
			}
			var vals []int
			for v := range h.Values() {
				vals = append(vals, v)
			}
			if !slices.Equal(vals, m.vals) && !(len(vals) == 0 && len(m.vals) == 0) {
				t.Fatalf("step %d: values %v, want %v", step, vals, m.vals)
			}
			if got := s.ToSlice(); !slices.Equal(got, m.keys) && !(len(got) == 0 && len(m.keys) == 0) {
				t.Fatalf("step %d: set elements %v, want %v", step, got, m.keys)
			}
		}
	}
}

func TestLinkedHashMapSmallAndIndexedModes(t *testing.T) {
	for seed := int64(1); seed <= 4; seed++ {
		for _, domain := range []int{3, 9, 12, 40} {
			checkLinkedHashMapAgainstModel(t, seed, domain, func(i int) string { return string(rune('a'+i%26)) + string(rune('0'+i/26)) })
			checkLinkedHashMapAgainstModel(t, seed, domain, func(i int) int32 { return int32(i * 7919) })
			type node struct{ x int }
			nodes := make([]*node, domain)
			for i := range nodes {
				nodes[i] = &node{i}
			}
			checkLinkedHashMapAgainstModel(t, seed, domain, func(i int) *node { return nodes[i] })
			checkLinkedHashMapAgainstModel(t, seed, domain, func(i int) any {
				if i%2 == 0 {
					return i
				}
				return string(rune('a' + i))
			})
		}
	}
}

// NaN is never equal to itself: like a Go map, a LinkedHashMap gives every
// NaN key its own entry and never finds one, before and after it builds
// its index.
func TestLinkedHashMapNaNKeys(t *testing.T) {
	h := NewLinkedHashMap[float64, int]()
	for i := 0; i < 2*lhmSmallSize; i++ {
		h.Put(math.NaN(), i)
		if h.Size() != i+1 || h.ContainsKey(math.NaN()) {
			t.Fatalf("after %d NaN puts: size %d, contains NaN %v", i+1, h.Size(), h.ContainsKey(math.NaN()))
		}
	}
	h.Put(0.0, 1)
	if v, ok := h.Lookup(math.Copysign(0, -1)); !ok || v != 1 {
		t.Fatalf("-0.0 does not find +0.0: %d, %v", v, ok)
	}
}

// Iteration skips entries removed during the loop and visits entries added
// during it, in both modes.
func TestLinkedHashMapMutationDuringIteration(t *testing.T) {
	for _, n := range []int{3, lhmSmallSize, 3 * lhmSmallSize} {
		h := NewLinkedHashMap[int, int]()
		for i := 0; i < n; i++ {
			h.Put(i, i)
		}
		var seen []int
		for k := range h.Keys() {
			seen = append(seen, k)
			if k == 0 {
				h.Remove(1)
				h.Put(100, 100)
			}
		}
		want := []int{0}
		for i := 2; i < n; i++ {
			want = append(want, i)
		}
		want = append(want, 100)
		if !slices.Equal(seen, want) {
			t.Fatalf("n=%d: iterated %v, want %v", n, seen, want)
		}
	}
}

type olhmKey struct{ id, hash int }

func (k olhmKey) Equals(o any) bool {
	x, ok := o.(olhmKey)
	return ok && x.id == k.id
}
func (k olhmKey) HashCode() int32 { return int32(k.hash) }

// ObjectLinkedHashMap keys that collide on hash code, across the switch to
// buckets and back.
func TestObjectLinkedHashMapSmallAndBucketModes(t *testing.T) {
	for seed := int64(1); seed <= 4; seed++ {
		r := rand.New(rand.NewSource(seed))
		h := NewObjectLinkedHashMap[olhmKey, int]()
		s := NewObjectLinkedHashSet[olhmKey]()
		var m orderedModel[olhmKey]
		for step := 0; step < 20000; step++ {
			id := r.Intn(30)
			k := olhmKey{id, id % 4}
			switch op := r.Intn(100); {
			case op < 45:
				v := r.Intn(1000) + 1
				gotOld := h.Put(k, v)
				wantOld, _ := m.put(k, v)
				if gotOld != wantOld {
					t.Fatalf("step %d: Put(%v) returned %d, want %d", step, k, gotOld, wantOld)
				}
				s.Add(k)
			case op < 75:
				got := h.Remove(k)
				want, present := m.remove(k)
				if got != want || s.Remove(k) != present {
					t.Fatalf("step %d: Remove(%v) returned %d, want %d", step, k, got, want)
				}
			case op < 77:
				h.Clear()
				s.Clear()
				m = orderedModel[olhmKey]{}
			default:
				v, ok := h.Lookup(k)
				i := m.index(k)
				if ok != (i >= 0) || (ok && v != m.vals[i]) || s.Contains(k) != ok {
					t.Fatalf("step %d: Lookup(%v) = %d, %v; model index %d", step, k, v, ok, i)
				}
			}
			if h.Size() != len(m.keys) || s.Size() != len(m.keys) {
				t.Fatalf("step %d: sizes %d, %d, want %d", step, h.Size(), s.Size(), len(m.keys))
			}
			if got := h.KeySlice(); !slices.Equal(got, m.keys) && !(len(got) == 0 && len(m.keys) == 0) {
				t.Fatalf("step %d: keys %v, want %v", step, got, m.keys)
			}
			if got := s.ToSlice(); !slices.Equal(got, m.keys) && !(len(got) == 0 && len(m.keys) == 0) {
				t.Fatalf("step %d: set elements %v, want %v", step, got, m.keys)
			}
		}
	}
}
