// Clean-room specification based on java.util.Random, Collections.shuffle, TreeMap, TreeSet, and PriorityQueue Javadoc.

package javautil

import (
	"cmp"
	"iter"
	"strings"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// ---------------------------------------------------------------------------
// java.util.Random (exact 48-bit LCG per java.util.Random Javadoc)
// ---------------------------------------------------------------------------

const (
	randomMultiplier int64 = 0x5DEECE66D
	randomAddend     int64 = 0xB
	randomMask       int64 = (1 << 48) - 1
)

// Random implements java.util.Random's 48-bit linear congruential generator.
type Random struct {
	seed atomic.Int64
}

// NewRandomSeed creates a Random initialized with seed per java.util.Random(long) Javadoc.
func NewRandomSeed(seed int64) *Random {
	r := &Random{}
	r.SetSeed(seed)
	return r
}

// NewRandom creates a Random seeded from System_NanoTime.
func NewRandom() *Random {
	return NewRandomSeed(javalang.System_NanoTime())
}

// SetSeed resets the seed to (seed ^ 0x5DEECE66DL) & ((1L << 48) - 1).
func (r *Random) SetSeed(seed int64) {
	r.seed.Store((seed ^ randomMultiplier) & randomMask)
}

func (r *Random) next(bits int) int32 {
	for {
		oldSeed := r.seed.Load()
		nextSeed := (oldSeed*randomMultiplier + randomAddend) & randomMask
		if r.seed.CompareAndSwap(oldSeed, nextSeed) {
			return int32(uint64(nextSeed) >> (48 - bits))
		}
	}
}

func (r *Random) NextInt() int32 {
	return r.next(32)
}

// NextIntBound implements Random.nextInt(int bound) per Java Javadoc.
func (r *Random) NextIntBound(bound int32) int32 {
	if bound <= 0 {
		panic(javalang.NewIllegalArgumentException("bound must be positive"))
	}
	if (bound & -bound) == bound { // power of 2
		return int32((int64(bound) * int64(r.next(31))) >> 31)
	}
	for {
		bits := r.next(31)
		val := bits % bound
		if bits-val+(bound-1) >= 0 {
			return val
		}
	}
}

func (r *Random) NextLong() int64 {
	return (int64(r.next(32)) << 32) + int64(r.next(32))
}

func (r *Random) NextBoolean() bool {
	return r.next(1) != 0
}

// Collections_Shuffle shuffles list in-place using rnd per java.util.Collections.shuffle Javadoc.
func Collections_Shuffle[T any](list *ArrayList[T], rnd *Random) {
	if list == nil || len(list.elems) <= 1 {
		return
	}
	if rnd == nil {
		rnd = NewRandom()
	}
	for i := len(list.elems); i > 1; i-- {
		j := int(rnd.NextIntBound(int32(i)))
		list.elems[i-1], list.elems[j] = list.elems[j], list.elems[i-1]
	}
}

// IdentityHashMap preserves deterministic insertion order while comparing pointer keys by identity.
type IdentityHashMap[K comparable, V any] = LinkedHashMap[K, V]

func NewIdentityHashMap[K comparable, V any]() *IdentityHashMap[K, V] {
	return NewLinkedHashMap[K, V]()
}

// ---------------------------------------------------------------------------
// java.util.TreeMap & java.util.TreeSet
// ---------------------------------------------------------------------------

type treeEntry[K any, V any] struct {
	key K
	val V
}

// TreeMap is an ordered map sorted according to a Comparator[K].
type TreeMap[K any, V any] struct {
	cmp     Comparator[K]
	entries []treeEntry[K, V]
}

func NewTreeMap[K cmp.Ordered, V any]() *TreeMap[K, V] {
	return &TreeMap[K, V]{cmp: NaturalOrder[K]()}
}

func NewTreeMapComparator[K any, V any](comparator Comparator[K]) *TreeMap[K, V] {
	if comparator == nil {
		panic(javalang.NewNullPointerExceptionMsg("comparator is null"))
	}
	return &TreeMap[K, V]{cmp: comparator}
}

func (m *TreeMap[K, V]) findIndex(k K) (int, bool) {
	lo, hi := 0, len(m.entries)
	for lo < hi {
		mid := int(uint(lo+hi) >> 1)
		c := m.cmp(m.entries[mid].key, k)
		if c < 0 {
			lo = mid + 1
		} else if c > 0 {
			hi = mid
		} else {
			return mid, true
		}
	}
	return lo, false
}

func (m *TreeMap[K, V]) Size() int {
	if m == nil {
		return 0
	}
	return len(m.entries)
}

func (m *TreeMap[K, V]) IsEmpty() bool { return m.Size() == 0 }

func (m *TreeMap[K, V]) Clear() {
	if m != nil {
		clear(m.entries)
		m.entries = m.entries[:0]
	}
}

func (m *TreeMap[K, V]) ContainsKey(k K) bool {
	if m == nil {
		return false
	}
	_, ok := m.findIndex(k)
	return ok
}

func (m *TreeMap[K, V]) Get(k K) (V, bool) {
	var zero V
	if m == nil {
		return zero, false
	}
	idx, ok := m.findIndex(k)
	if !ok {
		return zero, false
	}
	return m.entries[idx].val, true
}

func (m *TreeMap[K, V]) GetVal(k K) V {
	v, _ := m.Get(k)
	return v
}

func (m *TreeMap[K, V]) Put(k K, v V) (V, bool) {
	idx, ok := m.findIndex(k)
	if ok {
		prev := m.entries[idx].val
		m.entries[idx].val = v
		return prev, true
	}
	var zero V
	m.entries = append(m.entries, treeEntry[K, V]{})
	copy(m.entries[idx+1:], m.entries[idx:])
	m.entries[idx] = treeEntry[K, V]{key: k, val: v}
	return zero, false
}

func (m *TreeMap[K, V]) Remove(k K) (V, bool) {
	var zero V
	idx, ok := m.findIndex(k)
	if !ok {
		return zero, false
	}
	prev := m.entries[idx].val
	copy(m.entries[idx:], m.entries[idx+1:])
	m.entries[len(m.entries)-1] = treeEntry[K, V]{}
	m.entries = m.entries[:len(m.entries)-1]
	return prev, true
}

func (m *TreeMap[K, V]) FirstKey() K {
	if len(m.entries) == 0 {
		panic(javalang.NewNoSuchElementException(""))
	}
	return m.entries[0].key
}

func (m *TreeMap[K, V]) LastKey() K {
	if len(m.entries) == 0 {
		panic(javalang.NewNoSuchElementException(""))
	}
	return m.entries[len(m.entries)-1].key
}

func (m *TreeMap[K, V]) FloorKey(k K) (K, bool) {
	var zero K
	idx, ok := m.findIndex(k)
	if ok {
		return m.entries[idx].key, true
	}
	if idx == 0 {
		return zero, false
	}
	return m.entries[idx-1].key, true
}

func (m *TreeMap[K, V]) CeilingKey(k K) (K, bool) {
	var zero K
	idx, ok := m.findIndex(k)
	if ok {
		return m.entries[idx].key, true
	}
	if idx >= len(m.entries) {
		return zero, false
	}
	return m.entries[idx].key, true
}

func (m *TreeMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if m == nil {
			return
		}
		snap := make([]treeEntry[K, V], len(m.entries))
		copy(snap, m.entries)
		for _, e := range snap {
			if !yield(e.key, e.val) {
				return
			}
		}
	}
}

func (m *TreeMap[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for k := range m.All() {
			if !yield(k) {
				return
			}
		}
	}
}

func (m *TreeMap[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range m.All() {
			if !yield(v) {
				return
			}
		}
	}
}

func (m *TreeMap[K, V]) String() string {
	if m == nil {
		return "null"
	}
	var sb strings.Builder
	sb.WriteByte('{')
	for i, e := range m.entries {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(javalang.ToJavaString(e.key))
		sb.WriteByte('=')
		sb.WriteString(javalang.ToJavaString(e.val))
	}
	sb.WriteByte('}')
	return sb.String()
}

// TreeSet is an ordered set sorted according to a Comparator[T].
type TreeSet[T any] struct {
	tm *TreeMap[T, struct{}]
}

func NewTreeSet[T cmp.Ordered]() *TreeSet[T] {
	return &TreeSet[T]{tm: NewTreeMap[T, struct{}]()}
}

func NewTreeSetComparator[T any](comparator Comparator[T]) *TreeSet[T] {
	return &TreeSet[T]{tm: NewTreeMapComparator[T, struct{}](comparator)}
}

func (s *TreeSet[T]) Size() int          { return s.tm.Size() }
func (s *TreeSet[T]) IsEmpty() bool      { return s.tm.IsEmpty() }
func (s *TreeSet[T]) Clear()             { s.tm.Clear() }
func (s *TreeSet[T]) Contains(v T) bool  { return s.tm.ContainsKey(v) }
func (s *TreeSet[T]) First() T           { return s.tm.FirstKey() }
func (s *TreeSet[T]) Last() T            { return s.tm.LastKey() }
func (s *TreeSet[T]) All() iter.Seq[T]   { return s.tm.Keys() }
func (s *TreeSet[T]) Add(v T) bool       { _, existed := s.tm.Put(v, struct{}{}); return !existed }
func (s *TreeSet[T]) Remove(v T) bool    { _, existed := s.tm.Remove(v); return existed }
func (s *TreeSet[T]) Floor(v T) (T, bool)   { return s.tm.FloorKey(v) }
func (s *TreeSet[T]) Ceiling(v T) (T, bool) { return s.tm.CeilingKey(v) }

func (s *TreeSet[T]) PollFirst() (T, bool) {
	var zero T
	if s.IsEmpty() {
		return zero, false
	}
	first := s.First()
	s.Remove(first)
	return first, true
}

func (s *TreeSet[T]) ToSlice() []T {
	out := make([]T, len(s.tm.entries))
	for i, e := range s.tm.entries {
		out[i] = e.key
	}
	return out
}

// ---------------------------------------------------------------------------
// java.util.PriorityQueue (binary min-heap)
// ---------------------------------------------------------------------------

type PriorityQueue[T any] struct {
	cmp  Comparator[T]
	heap []T
}

func NewPriorityQueue[T cmp.Ordered]() *PriorityQueue[T] {
	return &PriorityQueue[T]{cmp: NaturalOrder[T]()}
}

func NewPriorityQueueComparator[T any](comparator Comparator[T]) *PriorityQueue[T] {
	if comparator == nil {
		panic(javalang.NewNullPointerExceptionMsg("comparator is null"))
	}
	return &PriorityQueue[T]{cmp: comparator}
}

func (pq *PriorityQueue[T]) Size() int {
	if pq == nil {
		return 0
	}
	return len(pq.heap)
}

func (pq *PriorityQueue[T]) IsEmpty() bool { return pq.Size() == 0 }

func (pq *PriorityQueue[T]) Clear() {
	if pq != nil {
		clear(pq.heap)
		pq.heap = pq.heap[:0]
	}
}

func (pq *PriorityQueue[T]) Offer(v T) bool {
	if javalang.IsNil(v) {
		panic(javalang.NewNullPointerException())
	}
	pq.heap = append(pq.heap, v)
	pq.siftUp(len(pq.heap) - 1)
	return true
}

func (pq *PriorityQueue[T]) Add(v T) bool { return pq.Offer(v) }

func (pq *PriorityQueue[T]) Peek() (T, bool) {
	var zero T
	if len(pq.heap) == 0 {
		return zero, false
	}
	return pq.heap[0], true
}

func (pq *PriorityQueue[T]) Poll() (T, bool) {
	var zero T
	n := len(pq.heap)
	if n == 0 {
		return zero, false
	}
	top := pq.heap[0]
	last := pq.heap[n-1]
	pq.heap[n-1] = zero
	pq.heap = pq.heap[:n-1]
	if n > 1 {
		pq.heap[0] = last
		pq.siftDown(0)
	}
	return top, true
}

func (pq *PriorityQueue[T]) Remove() T {
	v, ok := pq.Poll()
	if !ok {
		panic(javalang.NewNoSuchElementException(""))
	}
	return v
}

func (pq *PriorityQueue[T]) siftUp(k int) {
	x := pq.heap[k]
	for k > 0 {
		parent := int(uint(k-1) >> 1)
		e := pq.heap[parent]
		if pq.cmp(x, e) >= 0 {
			break
		}
		pq.heap[k] = e
		k = parent
	}
	pq.heap[k] = x
}

func (pq *PriorityQueue[T]) siftDown(k int) {
	n := len(pq.heap)
	half := int(uint(n) >> 1)
	x := pq.heap[k]
	for k < half {
		child := (k << 1) + 1
		c := pq.heap[child]
		right := child + 1
		if right < n && pq.cmp(c, pq.heap[right]) > 0 {
			child = right
			c = pq.heap[child]
		}
		if pq.cmp(x, c) <= 0 {
			break
		}
		pq.heap[k] = c
		k = child
	}
	pq.heap[k] = x
}

func (pq *PriorityQueue[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		if pq == nil {
			return
		}
		snap := make([]T, len(pq.heap))
		copy(snap, pq.heap)
		for _, v := range snap {
			if !yield(v) {
				return
			}
		}
	}
}
