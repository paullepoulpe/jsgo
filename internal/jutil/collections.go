package jutil

import (
	"hash/maphash"
	"iter"
	"slices"
	"strings"
	"sync/atomic"
	"unsafe"
)

// Go-only: language-adapter — java.util / Guava collections with Java
// iteration-order semantics. Ported code uses these (never a ranged Go map)
// wherever Java iteration order is observable.

// ---------------------------------------------------------------------------
// ArrayList

// ArrayList is java.util.ArrayList: a growable list with reference
// semantics (mutations through any alias are visible to all holders, unlike
// a bare Go slice). Use plain []T for Java ImmutableList / List.of values
// that are never mutated after construction.
type ArrayList[T any] struct {
	elems []T
}

// NewArrayList returns an empty list.
func NewArrayList[T any]() *ArrayList[T] { return &ArrayList[T]{} }

// NewArrayListWithCapacity is new ArrayList<>(capacity).
func NewArrayListWithCapacity[T any](capacity int) *ArrayList[T] {
	return &ArrayList[T]{elems: make([]T, 0, capacity)}
}

// InitWithBuffer makes l an empty list (new ArrayList<>(len(buf))) that
// stores its first len(buf) elements in buf; it grows into a new array, as
// any list does, past that.
//
// Go-only: performance — for a list held by value in another object with a
// small array next to it: the object, the list and its initial storage are
// then one allocation instead of three.
func (l *ArrayList[T]) InitWithBuffer(buf []T) {
	l.elems = buf[:0:len(buf)]
}

// NewArrayListFrom is new ArrayList<>(collection) / Lists.newArrayList(...).
func NewArrayListFrom[T any](elems []T) *ArrayList[T] {
	return &ArrayList[T]{elems: slices.Clone(elems)}
}

func (l *ArrayList[T]) Size() int     { return len(l.elems) }
func (l *ArrayList[T]) IsEmpty() bool { return len(l.elems) == 0 }

// Add appends e; always returns true like Collection.add.
func (l *ArrayList[T]) Add(e T) bool {
	l.elems = append(l.elems, e)
	return true
}

// AddAt is List.add(int, E).
func (l *ArrayList[T]) AddAt(index int, e T) {
	if index < 0 || index > len(l.elems) {
		panic(NewIndexOutOfBoundsException(LenientFormat("Index: %s, Size: %s", index, len(l.elems))))
	}
	l.elems = slices.Insert(l.elems, index, e)
}

// AddAll appends all elements.
func (l *ArrayList[T]) AddAll(es ...T) bool {
	l.elems = append(l.elems, es...)
	return len(es) > 0
}

func (l *ArrayList[T]) checkIndex(index int) {
	if uint(index) >= uint(len(l.elems)) {
		l.indexOutOfBounds(index)
	}
}

// indexOutOfBounds is kept out of line so that the checks above inline.
//
//go:noinline
func (l *ArrayList[T]) indexOutOfBounds(index int) {
	panic(NewIndexOutOfBoundsException(LenientFormat("Index %s out of bounds for length %s", index, len(l.elems))))
}

// Get is List.get(int).
func (l *ArrayList[T]) Get(index int) T {
	if uint(index) >= uint(len(l.elems)) {
		l.indexOutOfBounds(index)
	}
	return l.elems[index]
}

// Set is List.set(int, E); returns the previous element.
func (l *ArrayList[T]) Set(index int, e T) T {
	l.checkIndex(index)
	old := l.elems[index]
	l.elems[index] = e
	return old
}

// RemoveAt is List.remove(int); returns the removed element.
func (l *ArrayList[T]) RemoveAt(index int) T {
	l.checkIndex(index)
	old := l.elems[index]
	l.elems = slices.Delete(l.elems, index, index+1)
	return old
}

// Clear is List.clear().
func (l *ArrayList[T]) Clear() { clear(l.elems); l.elems = l.elems[:0] }

// IndexOf is List.indexOf(Object) using Go == on the dynamic values
// (identity for pointers, value equality for strings and numbers).
func (l *ArrayList[T]) IndexOf(e T) int {
	for i, x := range l.elems {
		if any(x) == any(e) {
			return i
		}
	}
	return -1
}

// Contains is List.contains(Object) with the same equality as IndexOf.
func (l *ArrayList[T]) Contains(e T) bool { return l.IndexOf(e) >= 0 }

// Remove is List.remove(Object): removes the first equal element.
func (l *ArrayList[T]) Remove(e T) bool {
	if i := l.IndexOf(e); i >= 0 {
		l.elems = slices.Delete(l.elems, i, i+1)
		return true
	}
	return false
}

// GetFirst / GetLast are List.getFirst()/getLast() (Java 21).
func (l *ArrayList[T]) GetFirst() T { l.checkIndex(0); return l.elems[0] }
func (l *ArrayList[T]) GetLast() T  { l.checkIndex(len(l.elems) - 1); return l.elems[len(l.elems)-1] }

// RemoveLast is List.removeLast() (Java 21).
func (l *ArrayList[T]) RemoveLast() T { return l.RemoveAt(len(l.elems) - 1) }

// Slice returns the backing slice. Callers must not retain it across
// mutations of the list.
func (l *ArrayList[T]) Slice() []T {
	if l == nil {
		return nil
	}
	return l.elems
}

// ToSlice returns a copy (Java: ImmutableList.copyOf(list)).
func (l *ArrayList[T]) ToSlice() []T { return slices.Clone(l.elems) }

// All iterates in index order (Java: for (T e : list)).
func (l *ArrayList[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, e := range l.elems {
			if !yield(e) {
				return
			}
		}
	}
}

// String renders like AbstractCollection.toString(): "[a, b]".
func (l *ArrayList[T]) String() string { return CollectionToString(l.elems) }

// CollectionToString renders elements like Java's AbstractCollection.toString.
func CollectionToString[T any](elems []T) string {
	var sb strings.Builder
	sb.WriteByte('[')
	for i, e := range elems {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(ToJavaString(e))
	}
	sb.WriteByte(']')
	return sb.String()
}

// ---------------------------------------------------------------------------
// LinkedHashMap

// LhmSlot is an entry slot of a LinkedHashMap.
//
// Go-only: performance — storing entries in a flat contiguous slice instead
// of a doubly-linked list of heap-allocated nodes eliminates per-entry heap
// allocations and the 16-byte prev/next pointer overhead per entry. The
// value comes first because Go pads a struct that ends in a zero-size field
// (the struct{} values of a LinkedHashSet) by a word.
type LhmSlot[K comparable, V any] struct {
	value V
	key   K
}

// Key returns the slot's key.
func (s *LhmSlot[K, V]) Key() K { return s.key }

// Value returns the slot's value.
func (s *LhmSlot[K, V]) Value() V { return s.value }

// LinkedHashMap is java.util.LinkedHashMap (insertion order). It also models
// Guava ImmutableMap and any Java map whose iteration order is insertion
// order. Keys use Go equality, which matches Java equals() for strings,
// boxed numbers, enums and identity-compared objects (pointers).
//
// Go-only: performance — most maps of a compilation stay small (in the
// typed-app benchmark, half of them never get an entry and 98% never hold
// more than 8), so the hash index of the entries (lhmIndex) is only built
// once the slot slice holds more than lhmSmallSize entries; until then
// lookups scan the contiguous slots slice, comparing keys with == as a Go
// map does.
type LinkedHashMap[K comparable, V any] struct {
	m     *lhmIndex // nil until len(slots) exceeds lhmSmallSize or a non-front slot is removed
	slots []LhmSlot[K, V]
	size  int32
	first int32
}

// lhmIndex is the hash index of a LinkedHashMap's slots.
//
// Go-only: performance — an open-addressing table of 32-bit slot numbers
// costs 5 to 16 bytes per entry and holds no pointers for the GC to scan;
// a Go map[K]int costs 25 to 60 bytes per entry and, for string or pointer
// keys, is scanned by every GC cycle.
type lhmIndex struct {
	// table is a linear-probing hash table whose length is a power of two.
	// A bucket holds 0 (empty), lhmTombstone (its entry was removed) or the
	// slot index + 1 of an entry. At most 3/4 of the buckets are non-empty,
	// so every probe sequence ends at an empty bucket.
	table   []uint32
	used    int    // non-empty buckets, live or tombstones
	removed []bool // non-nil only when a slot > first is removed
}

const lhmTombstone = ^uint32(0)

// lhmSmallSize is the most entries a LinkedHashMap or ObjectLinkedHashMap
// holds without an index.
const lhmSmallSize = 8

// lhmSeed seeds the hash of LinkedHashMap keys. Iteration order never
// depends on it.
var lhmSeed = maphash.MakeSeed()

// lhmHash hashes k as a Go map would.
func lhmHash[K comparable](k K) uint {
	// Go-only: language-adapter — Go 1.24 compiles maphash.Comparable[K] to
	// a panic when K is a generic interface type instantiated in generic
	// code (it does not apply its escapeForHash intrinsic there), so
	// interface keys are hashed as any.
	var zero K
	if any(zero) == nil {
		return uint(maphash.Comparable(lhmSeed, any(k)))
	}
	return uint(maphash.Comparable(lhmSeed, k))
}

// lhmTableLen is the table length for an index of n entries: at most half
// full, so that the table takes at least n/2 more insertions to fill to 3/4.
func lhmTableLen(n int) int {
	l := 2 * lhmSmallSize
	for l < 2*n {
		l *= 2
	}
	return l
}

// insert records slot (whose key hashes to hash and is absent from the
// index) in the first empty or tombstone bucket of its probe sequence. The
// caller makes sure that the table has room (needsGrow).
func (x *lhmIndex) insert(hash uint, slot int) {
	t := x.table
	mask := uint(len(t) - 1)
	b := hash & mask
	for t[b] != 0 && t[b] != lhmTombstone {
		b = (b + 1) & mask
	}
	if t[b] == 0 {
		x.used++
	}
	t[b] = uint32(slot + 1)
}

// erase turns the bucket of slot (whose key hashes to hash) into a tombstone.
func (x *lhmIndex) erase(hash uint, slot int) {
	t := x.table
	mask := uint(len(t) - 1)
	b := hash & mask
	for t[b] != uint32(slot+1) {
		b = (b + 1) & mask
	}
	t[b] = lhmTombstone
}

// reset empties the table and gives it n buckets.
func (x *lhmIndex) reset(n int) {
	if n == len(x.table) {
		clear(x.table)
	} else {
		x.table = make([]uint32, n)
	}
	x.used = 0
}

// needsGrow reports whether one more insertion would fill more than 3/4 of
// the table.
func (x *lhmIndex) needsGrow() bool {
	return 4*(x.used+1) > 3*len(x.table)
}

func (x *lhmIndex) isRemoved(i int) bool {
	return i < len(x.removed) && x.removed[i]
}

// NewLinkedHashMap returns an empty map.
func NewLinkedHashMap[K comparable, V any]() *LinkedHashMap[K, V] {
	return &LinkedHashMap[K, V]{}
}

// NewLinkedHashMapWithCapacity returns an empty map pre-sized for cap entries
// (or the few more that the allocation's size class holds).
func NewLinkedHashMapWithCapacity[K comparable, V any](cap int) *LinkedHashMap[K, V] {
	if cap > 0 {
		return &LinkedHashMap[K, V]{
			slots: slices.Grow([]LhmSlot[K, V](nil), cap),
		}
	}
	return &LinkedHashMap[K, V]{}
}

// InitWithBuffer makes h an empty map that stores its first len(buf) entries
// in buf; it grows into a new slice past that.
func (h *LinkedHashMap[K, V]) InitWithBuffer(buf []LhmSlot[K, V]) {
	h.m = nil
	h.slots = buf[:0:len(buf)]
	h.size = 0
	h.first = 0
}

// CompactLen returns len(h.slots) when every slot is live (no entry was
// removed), or -1 otherwise.
func (h *LinkedHashMap[K, V]) CompactLen() int {
	if h == nil {
		return 0
	}
	if int(h.size) == len(h.slots) {
		return len(h.slots)
	}
	return -1
}

// ValueAtSlot returns h.slots[i].value for 0 <= i < h.CompactLen().
func (h *LinkedHashMap[K, V]) ValueAtSlot(i int) V {
	return h.slots[i].value
}

func (h *LinkedHashMap[K, V]) loadIndex() *lhmIndex {
	return (*lhmIndex)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&h.m))))
}

// find returns k's slot index in h.slots, or -1 if absent.
func (h *LinkedHashMap[K, V]) find(k K) int {
	if m := h.loadIndex(); m != nil {
		return h.lookupIn(m, k)
	}
	if int(h.size) > lhmSmallSize || len(h.slots) > 2*lhmSmallSize {
		return h.lookupIn(h.buildIndexShared(), k)
	}
	for i := int(h.first); i < len(h.slots); i++ {
		if h.slots[i].key == k {
			return i
		}
	}
	return -1
}

// lookup is find through the index.
func (h *LinkedHashMap[K, V]) lookup(k K) int {
	return h.lookupIn(h.m, k)
}

func (h *LinkedHashMap[K, V]) lookupIn(m *lhmIndex, k K) int {
	t := m.table
	mask := uint(len(t) - 1)
	for b := lhmHash(k) & mask; ; b = (b + 1) & mask {
		switch e := t[b]; e {
		case 0:
			return -1
		case lhmTombstone:
		default:
			if h.slots[e-1].key == k {
				return int(e - 1)
			}
		}
	}
}

func (h *LinkedHashMap[K, V]) Size() int {
	if h == nil {
		return 0
	}
	return int(h.size)
}
func (h *LinkedHashMap[K, V]) IsEmpty() bool { return h.Size() == 0 }

// Get is Map.get: the zero value (Java null) when absent.
func (h *LinkedHashMap[K, V]) Get(k K) V {
	if h != nil {
		if i := h.find(k); i >= 0 {
			return h.slots[i].value
		}
	}
	var zero V
	return zero
}

// Lookup returns the value and whether the key is present.
func (h *LinkedHashMap[K, V]) Lookup(k K) (V, bool) {
	if h != nil {
		if i := h.find(k); i >= 0 {
			return h.slots[i].value, true
		}
	}
	var zero V
	return zero, false
}

// GetOrDefault is Map.getOrDefault.
func (h *LinkedHashMap[K, V]) GetOrDefault(k K, def V) V {
	if v, ok := h.Lookup(k); ok {
		return v
	}
	return def
}

// ContainsKey is Map.containsKey.
func (h *LinkedHashMap[K, V]) ContainsKey(k K) bool {
	return h != nil && h.find(k) >= 0
}

// Put is Map.put; re-putting an existing key keeps its position. Returns the
// previous value (zero if none).
func (h *LinkedHashMap[K, V]) Put(k K, v V) V {
	if i := h.find(k); i >= 0 {
		old := h.slots[i].value
		h.slots[i].value = v
		return old
	}
	h.add(k, v)
	var zero V
	return zero
}

// PutNew inserts (k, v) when the caller has already verified that k is absent.
func (h *LinkedHashMap[K, V]) PutNew(k K, v V) {
	h.add(k, v)
}

// HasBacking reports whether h already has backing slot capacity.
func (h *LinkedHashMap[K, V]) HasBacking() bool {
	return cap(h.slots) > 0
}

// Cap returns the current slot slice capacity of h.
func (h *LinkedHashMap[K, V]) Cap() int {
	if h == nil {
		return 0
	}
	return cap(h.slots)
}

// add appends an entry for k, which must be absent.
func (h *LinkedHashMap[K, V]) add(k K, v V) {
	if len(h.slots) < cap(h.slots) && h.m == nil {
		h.slots = append(h.slots, LhmSlot[K, V]{key: k, value: v})
		h.size++
		return
	}
	h.addSlow(k, v)
}

//go:noinline
func (h *LinkedHashMap[K, V]) addSlow(k K, v V) {
	if cap(h.slots) == 0 {
		h.slots = make([]LhmSlot[K, V], 0, 2)
		h.first = 0
	} else if len(h.slots) >= 64 && int(h.size)*4 < len(h.slots) {
		h.compact()
	}
	idx := len(h.slots)
	h.slots = append(h.slots, LhmSlot[K, V]{key: k, value: v})
	h.size++
	if h.m != nil {
		if h.m.needsGrow() {
			h.reindex()
		} else {
			h.m.insert(lhmHash(k), idx)
		}
	}
}

// buildIndex gives h an index of its live slots.
//
//go:noinline
func (h *LinkedHashMap[K, V]) buildIndex() {
	x := &lhmIndex{}
	x.reset(lhmTableLen(int(h.size)))
	for i := int(h.first); i < len(h.slots); i++ {
		x.insert(lhmHash(h.slots[i].key), i)
	}
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&h.m)), unsafe.Pointer(x))
}

// buildIndexShared builds an index of h's live slots and publishes it
// atomically so concurrent read-only lookups on an unindexed map are race-free.
//
//go:noinline
func (h *LinkedHashMap[K, V]) buildIndexShared() *lhmIndex {
	x := &lhmIndex{}
	x.reset(lhmTableLen(int(h.size)))
	for i := int(h.first); i < len(h.slots); i++ {
		x.insert(lhmHash(h.slots[i].key), i)
	}
	p := (*unsafe.Pointer)(unsafe.Pointer(&h.m))
	if atomic.CompareAndSwapPointer(p, nil, unsafe.Pointer(x)) {
		return x
	}
	return (*lhmIndex)(atomic.LoadPointer(p))
}

// reindex rebuilds the index table from the live slots, sized for them.
func (h *LinkedHashMap[K, V]) reindex() {
	x := h.m
	x.reset(lhmTableLen(int(h.size)))
	for i := int(h.first); i < len(h.slots); i++ {
		if !x.isRemoved(i) {
			x.insert(lhmHash(h.slots[i].key), i)
		}
	}
}

// compact drops the removed slots, keeping the order of the others.
//
//go:noinline
func (h *LinkedHashMap[K, V]) compact() {
	w := 0
	for i := int(h.first); i < len(h.slots); i++ {
		if h.m == nil || !h.m.isRemoved(i) {
			h.slots[w] = h.slots[i]
			w++
		}
	}
	clear(h.slots[w:])
	h.slots = h.slots[:w]
	h.first = 0
	if h.m != nil {
		clear(h.m.removed)
		h.m.removed = h.m.removed[:0]
		h.reindex()
	}
}

// PutIfAbsent is Map.putIfAbsent: if the key is absent or mapped to null
// (jutil.IsNil) it is (re)mapped to v and null is returned; otherwise the
// current value is returned.
func (h *LinkedHashMap[K, V]) PutIfAbsent(k K, v V) V {
	if i := h.find(k); i >= 0 {
		if !IsNil(h.slots[i].value) {
			return h.slots[i].value
		}
		h.slots[i].value = v
		var zero V
		return zero
	}
	h.add(k, v)
	var zero V
	return zero
}

// LoadOrStore returns the existing value and true, or inserts v and returns
// it with false. Unlike Java putIfAbsent it treats a present null as present.
// Go-only: performance — one lookup for analysis bookkeeping, including misses.
func (h *LinkedHashMap[K, V]) LoadOrStore(k K, v V) (V, bool) {
	if i := h.find(k); i >= 0 {
		return h.slots[i].value, true
	}
	h.add(k, v)
	return v, false
}

// ComputeIfAbsent is Map.computeIfAbsent (the mapping function must not
// return Java null; Go has no way to express that distinction generically).
func (h *LinkedHashMap[K, V]) ComputeIfAbsent(k K, f func(K) V) V {
	if i := h.find(k); i >= 0 {
		return h.slots[i].value
	}
	v := f(k)
	h.Put(k, v)
	return v
}

// Remove is Map.remove; returns the removed value (zero if absent).
func (h *LinkedHashMap[K, V]) Remove(k K) V {
	i := h.find(k)
	if i < 0 {
		var zero V
		return zero
	}
	old := h.slots[i].value
	h.removeSlot(i)
	return old
}

// removeSlot removes the entry at slot index i.
func (h *LinkedHashMap[K, V]) removeSlot(i int) {
	if i != int(h.first) {
		h.removeSlotMiddle(i)
		return
	}
	if h.m != nil {
		h.m.erase(lhmHash(h.slots[i].key), i)
	}
	h.size--
	h.slots[i] = LhmSlot[K, V]{}
	first := i + 1
	if h.m != nil {
		for first < len(h.slots) && h.m.isRemoved(first) {
			first++
		}
	}
	h.first = int32(first)
}

//go:noinline
func (h *LinkedHashMap[K, V]) removeSlotMiddle(i int) {
	if h.m == nil {
		h.buildIndex()
	}
	h.m.erase(lhmHash(h.slots[i].key), i)
	h.size--
	h.slots[i] = LhmSlot[K, V]{}
	if i >= len(h.m.removed) {
		newLen := max(i+1, len(h.m.removed)*2, len(h.slots))
		newRemoved := make([]bool, newLen)
		copy(newRemoved, h.m.removed)
		h.m.removed = newRemoved
	}
	h.m.removed[i] = true
}

// Clear is Map.clear.
func (h *LinkedHashMap[K, V]) Clear() {
	clear(h.slots)
	h.slots = h.slots[:0]
	if h.m != nil {
		h.m.reset(len(h.m.table))
		clear(h.m.removed)
		h.m.removed = h.m.removed[:0]
	}
	h.size = 0
	h.first = 0
}

// ForEachValue calls fn for each value in insertion order without allocating an iterator closure.
func (h *LinkedHashMap[K, V]) ForEachValue(fn func(V)) {
	if h == nil || h.size == 0 {
		return
	}
	for i := int(h.first); i < len(h.slots); i++ {
		if m := h.loadIndex(); i < int(h.first) || (m != nil && i < len(m.removed) && m.removed[i]) {
			continue
		}
		fn(h.slots[i].value)
	}
}

// CompactSlots returns h's contiguous slot slice and true when h has no removed
// entries (including when h is nil or empty), enabling zero-allocation loops.
func (h *LinkedHashMap[K, V]) CompactSlots() ([]LhmSlot[K, V], bool) {
	if h == nil {
		return nil, true
	}
	if int(h.size) == len(h.slots) {
		return h.slots, true
	}
	return nil, false
}

// All iterates entries in insertion order. Entries removed during iteration
// are skipped; entries added during iteration are visited.
func (h *LinkedHashMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if h == nil {
			return
		}
		for i := int(h.first); i < len(h.slots); i++ {
			if m := h.loadIndex(); i < int(h.first) || m != nil && m.isRemoved(i) {
				continue
			}
			if !yield(h.slots[i].key, h.slots[i].value) {
				return
			}
		}
	}
}

// KeyAt returns the index-th key in insertion order without allocating a slice.
func (h *LinkedHashMap[K, V]) KeyAt(index int) (K, bool) {
	var zero K
	if h == nil || index < 0 || index >= int(h.size) {
		return zero, false
	}
	m := h.loadIndex()
	if h.first == 0 && (m == nil || len(m.removed) == 0) {
		return h.slots[index].key, true
	}
	cur := 0
	for i := int(h.first); i < len(h.slots); i++ {
		if m != nil && m.isRemoved(i) {
			continue
		}
		if cur == index {
			return h.slots[i].key, true
		}
		cur++
	}
	return zero, false
}

// ValueAt returns the index-th value in insertion order without allocating a slice.
func (h *LinkedHashMap[K, V]) ValueAt(index int) (V, bool) {
	var zero V
	if h == nil || index < 0 || index >= int(h.size) {
		return zero, false
	}
	m := h.loadIndex()
	if h.first == 0 && (m == nil || len(m.removed) == 0) {
		return h.slots[index].value, true
	}
	cur := 0
	for i := int(h.first); i < len(h.slots); i++ {
		if m != nil && m.isRemoved(i) {
			continue
		}
		if cur == index {
			return h.slots[i].value, true
		}
		cur++
	}
	return zero, false
}

// NextSlotValue advances from raw slot index `slot` (pass 0 initially) and
// returns the next live value, the next raw slot index to pass on the
// following call, and true; or the zero value, len(h.slots), and false.
func (h *LinkedHashMap[K, V]) NextSlotValue(slot int) (V, int, bool) {
	var zero V
	if h == nil {
		return zero, 0, false
	}
	if slot < int(h.first) {
		slot = int(h.first)
	}
	m := h.loadIndex()
	for i := slot; i < len(h.slots); i++ {
		if m != nil && m.isRemoved(i) {
			continue
		}
		return h.slots[i].value, i + 1, true
	}
	return zero, len(h.slots), false
}

// Keys iterates keys in insertion order (Java: map.keySet()).
//
// Go-only: performance — Keys, Values and the sets' All walk the slot slice
// themselves, like All, instead of ranging over All: nested iterator
// closures escape to the heap (seven allocations per loop over a
// LinkedHashSet), flat ones are inlined into the caller's loop.
func (h *LinkedHashMap[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		if h == nil {
			return
		}
		for i := int(h.first); i < len(h.slots); i++ {
			if m := h.loadIndex(); i < int(h.first) || m != nil && m.isRemoved(i) {
				continue
			}
			if !yield(h.slots[i].key) {
				return
			}
		}
	}
}

// Values iterates values in insertion order (Java: map.values()).
func (h *LinkedHashMap[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		if h == nil {
			return
		}
		for i := int(h.first); i < len(h.slots); i++ {
			if m := h.loadIndex(); i < int(h.first) || m != nil && m.isRemoved(i) {
				continue
			}
			if !yield(h.slots[i].value) {
				return
			}
		}
	}
}

// KeySlice returns the keys in insertion order.
func (h *LinkedHashMap[K, V]) KeySlice() []K {
	n := h.Size()
	if n == 0 {
		return nil
	}
	out := make([]K, 0, n)
	slots := h.slots
	if n == len(slots) {
		for i := range slots {
			out = append(out, slots[i].key)
		}
		return out
	}
	m := h.loadIndex()
	for i := int(h.first); i < len(slots); i++ {
		if m == nil || !m.isRemoved(i) {
			out = append(out, slots[i].key)
		}
	}
	return out
}

// String renders like AbstractMap.toString(): "{k1=v1, k2=v2}".
func (h *LinkedHashMap[K, V]) String() string {
	var sb strings.Builder
	sb.WriteByte('{')
	first := true
	for k, v := range h.All() {
		if !first {
			sb.WriteString(", ")
		}
		first = false
		sb.WriteString(ToJavaString(k))
		sb.WriteByte('=')
		sb.WriteString(ToJavaString(v))
	}
	sb.WriteByte('}')
	return sb.String()
}

// ---------------------------------------------------------------------------
// LinkedHashSet

// LinkedHashSet is java.util.LinkedHashSet (insertion order); it also models
// Guava ImmutableSet.
type LinkedHashSet[T comparable] struct {
	// Go-only: performance — held by value: one allocation per set.
	m LinkedHashMap[T, struct{}]
}

// NewLinkedHashSet returns a set containing elems in order.
func NewLinkedHashSet[T comparable](elems ...T) *LinkedHashSet[T] {
	s := &LinkedHashSet[T]{}
	for _, e := range elems {
		s.Add(e)
	}
	return s
}

// InitWithBuffer makes s an empty set that stores its first len(buf) elements
// in buf; it grows into a new slice past that.
func (s *LinkedHashSet[T]) InitWithBuffer(buf []LhmSlot[T, struct{}]) {
	s.m.InitWithBuffer(buf)
}

func (s *LinkedHashSet[T]) Size() int {
	if s == nil {
		return 0
	}
	return int(s.m.size)
}
func (s *LinkedHashSet[T]) IsEmpty() bool { return s.Size() == 0 }

// Add is Set.add; reports whether the set changed.
func (s *LinkedHashSet[T]) Add(e T) bool {
	if s.m.find(e) >= 0 {
		return false
	}
	s.m.add(e, struct{}{})
	return true
}

// AddNew inserts e when the caller has already verified that e is absent.
func (s *LinkedHashSet[T]) AddNew(e T) {
	s.m.add(e, struct{}{})
}

// AddAll adds every element; reports whether the set changed.
func (s *LinkedHashSet[T]) AddAll(es ...T) bool {
	changed := false
	for _, e := range es {
		if s.Add(e) {
			changed = true
		}
	}
	return changed
}

// Contains is Set.contains.
func (s *LinkedHashSet[T]) Contains(e T) bool { return s != nil && s.m.find(e) >= 0 }

// Remove is Set.remove; reports whether the element was present.
func (s *LinkedHashSet[T]) Remove(e T) bool {
	i := s.m.find(e)
	if i < 0 {
		return false
	}
	s.m.removeSlot(i)
	return true
}

// Clear is Set.clear.
func (s *LinkedHashSet[T]) Clear() { s.m.Clear() }

// All iterates in insertion order (see LinkedHashMap.Keys).
func (s *LinkedHashSet[T]) All() iter.Seq[T] {
	var h *LinkedHashMap[T, struct{}]
	if s != nil {
		h = &s.m
	}
	return func(yield func(T) bool) {
		if h == nil {
			return
		}
		for i := int(h.first); i < len(h.slots); i++ {
			if m := h.loadIndex(); i < int(h.first) || m != nil && m.isRemoved(i) {
				continue
			}
			if !yield(h.slots[i].key) {
				return
			}
		}
	}
}

// ToSlice returns the elements in insertion order.
func (s *LinkedHashSet[T]) ToSlice() []T {
	if s == nil {
		return nil
	}
	return s.m.KeySlice()
}

// First returns the first element in insertion order, if any, without closures.
func (s *LinkedHashSet[T]) First() (T, bool) {
	if s != nil && s.m.size > 0 {
		h := &s.m
		m := h.loadIndex()
		if m == nil || len(m.removed) == 0 {
			return h.slots[h.first].key, true
		}
		for i := int(h.first); i < len(h.slots); i++ {
			if i >= len(m.removed) || !m.removed[i] {
				return h.slots[i].key, true
			}
		}
	}
	var zero T
	return zero, false
}

// FirstTwo returns the first two elements in insertion order when Size() >= 2.
func (s *LinkedHashSet[T]) FirstTwo() (T, T, bool) {
	if s != nil && s.m.size >= 2 {
		h := &s.m
		m := h.loadIndex()
		if m == nil || len(m.removed) == 0 {
			f := int(h.first)
			return h.slots[f].key, h.slots[f+1].key, true
		}
		var a T
		gotFirst := false
		for i := int(h.first); i < len(h.slots); i++ {
			if i < len(m.removed) && m.removed[i] {
				continue
			}
			if !gotFirst {
				a = h.slots[i].key
				gotFirst = true
			} else {
				return a, h.slots[i].key, true
			}
		}
	}
	var zero T
	return zero, zero, false
}

// CopyTo copies up to len(dst) elements in insertion order into dst without
// allocating an iterator closure, returning the number of elements copied.
func (s *LinkedHashSet[T]) CopyTo(dst []T) int {
	if s == nil || s.m.size == 0 || len(dst) == 0 {
		return 0
	}
	h := &s.m
	m := h.loadIndex()
	n := 0
	if m == nil || len(m.removed) == 0 {
		for i := int(h.first); i < len(h.slots) && n < len(dst); i++ {
			dst[n] = h.slots[i].key
			n++
		}
		return n
	}
	for i := int(h.first); i < len(h.slots) && n < len(dst); i++ {
		if i < len(m.removed) && m.removed[i] {
			continue
		}
		dst[n] = h.slots[i].key
		n++
	}
	return n
}

// String renders like AbstractCollection.toString(): "[a, b]".
func (s *LinkedHashSet[T]) String() string { return CollectionToString(s.ToSlice()) }

// ---------------------------------------------------------------------------
// Collection / Map markers and value equality

// JavaCollection is implemented by the jutil types that model a
// java.util.Collection, so that `x instanceof Collection` can be expressed
// as a Go type assertion.
type JavaCollection interface {
	Size() int
	IsEmpty() bool
	javaCollection()
}

// JavaMap is implemented by the jutil types that model a java.util.Map.
type JavaMap interface {
	Size() int
	IsEmpty() bool
	javaMap()
}

func (*ArrayList[T]) javaCollection()           {}
func (*LinkedHashSet[T]) javaCollection()       {}
func (*LinkedHashMap[K, V]) javaMap()           {}
func (*ObjectLinkedHashSet[T]) javaCollection() {}
func (*ObjectLinkedHashMap[K, V]) javaMap()     {}
func (*TreeMap[K, V]) javaMap()                 {}

// javaList is implemented by list types so that List.equals can compare
// across implementations.
type javaList interface {
	javaListElems() []any
}

func (l *ArrayList[T]) javaListElems() []any {
	out := make([]any, len(l.elems))
	for i, e := range l.elems {
		out[i] = e
	}
	return out
}

// Equals is List.equals(Object): same size and pairwise Objects.equals.
func (l *ArrayList[T]) Equals(o any) bool {
	if any(l) == o {
		return true
	}
	other, ok := o.(javaList)
	if !ok || IsNil(o) {
		return false
	}
	oe := other.javaListElems()
	if len(oe) != len(l.elems) {
		return false
	}
	for i, e := range l.elems {
		if !ObjectsEquals(e, oe[i]) {
			return false
		}
	}
	return true
}

// HashCode is List.hashCode().
func (l *ArrayList[T]) HashCode() int32 {
	var h int32 = 1
	for _, e := range l.elems {
		h = 31*h + ObjectHashCodeOf(e)
	}
	return h
}

// javaSet is implemented by set types so that Set.equals can compare across
// implementations.
type javaSet interface {
	Size() int
	containsAny(e any) bool
}

func (s *LinkedHashSet[T]) containsAny(e any) bool {
	x, ok := e.(T)
	return ok && s.Contains(x)
}

// Equals is Set.equals(Object): same size and every element contained.
func (s *LinkedHashSet[T]) Equals(o any) bool {
	if any(s) == o {
		return true
	}
	other, ok := o.(javaSet)
	if !ok || IsNil(o) || other.Size() != s.Size() {
		return false
	}
	for e := range s.All() {
		if !other.containsAny(e) {
			return false
		}
	}
	return true
}

// HashCode is Set.hashCode(): the sum of the element hash codes.
func (s *LinkedHashSet[T]) HashCode() int32 {
	if s == nil || s.m.size == 0 {
		return 0
	}
	var h int32
	var removed []bool
	if m := s.m.loadIndex(); m != nil {
		removed = m.removed
	}
	for i := int(s.m.first); i < len(s.m.slots); i++ {
		if i < len(removed) && removed[i] {
			continue
		}
		h += ObjectHashCodeOf(s.m.slots[i].key)
	}
	return h
}

// javaMapEntries is implemented by map types so that Map.equals can compare
// across implementations.
type javaMapEntries interface {
	Size() int
	lookupAny(k any) (any, bool)
}

func (h *LinkedHashMap[K, V]) lookupAny(k any) (any, bool) {
	key, ok := k.(K)
	if !ok {
		return nil, false
	}
	return h.Lookup(key)
}

// Equals is Map.equals(Object): same size and equal values for every key.
func (h *LinkedHashMap[K, V]) Equals(o any) bool {
	if any(h) == o {
		return true
	}
	other, ok := o.(javaMapEntries)
	if !ok || IsNil(o) || other.Size() != h.Size() {
		return false
	}
	for k, v := range h.All() {
		ov, present := other.lookupAny(k)
		if !present || !ObjectsEquals(v, ov) {
			return false
		}
	}
	return true
}

// HashCode is Map.hashCode(): the sum of key^value hash codes.
func (h *LinkedHashMap[K, V]) HashCode() int32 {
	if h == nil || h.size == 0 {
		return 0
	}
	var sum int32
	var removed []bool
	if m := h.loadIndex(); m != nil {
		removed = m.removed
	}
	for i := int(h.first); i < len(h.slots); i++ {
		if i < len(removed) && removed[i] {
			continue
		}
		sum += ObjectHashCodeOf(h.slots[i].key) ^ ObjectHashCodeOf(h.slots[i].value)
	}
	return sum
}

// NewLinkedHashMapFrom is new LinkedHashMap<>(map) / ImmutableMap.copyOf.
func NewLinkedHashMapFrom[K comparable, V any](src *LinkedHashMap[K, V]) *LinkedHashMap[K, V] {
	n := src.Size()
	if n == 0 {
		return &LinkedHashMap[K, V]{}
	}
	if int(src.size) == len(src.slots) && src.first == 0 {
		slots := make([]LhmSlot[K, V], n)
		copy(slots, src.slots)
		return &LinkedHashMap[K, V]{
			slots: slots,
			size:  int32(n),
		}
	}
	out := NewLinkedHashMapWithCapacity[K, V](n)
	for k, v := range src.All() {
		out.PutNew(k, v)
	}
	return out
}

// ---------------------------------------------------------------------------
// ObjectLinkedHashMap / ObjectLinkedHashSet

type olhmSlot[K any, V any] struct {
	key     K
	value   V
	hash    int32 // ObjectHashCode(key)
	removed bool
}

// ObjectLinkedHashMap is java.util.LinkedHashMap (and Guava ImmutableMap)
// for keys whose classes override equals()/hashCode(): keys are compared
// with ObjectsEquals and hashed with ObjectHashCode instead of Go ==.
// Iteration is in insertion order.
//
// Go-only: performance — as for LinkedHashMap, the index is only built once
// the slot slice holds more than lhmSmallSize entries (in the typed-app
// benchmark, 65% of these maps stay empty); until then lookups walk the
// contiguous slots slice and call ObjectsEquals on the entries whose hash
// matches. The index hashes the hash codes stored in the slots.
type ObjectLinkedHashMap[K any, V any] struct {
	m     *lhmIndex // nil until len(slots) exceeds lhmSmallSize; m.removed is unused
	slots []olhmSlot[K, V]
	size  int32
	first int32
}

// OlhmSlot is an entry slot of ObjectLinkedHashMap / ObjectLinkedHashSet for
// callers that provide an inline or slab-allocated initial buffer via InitWithBuffer.
type OlhmSlot[K any, V any] = olhmSlot[K, V]

// NewObjectLinkedHashMap returns an empty map.
func NewObjectLinkedHashMap[K any, V any]() *ObjectLinkedHashMap[K, V] {
	return &ObjectLinkedHashMap[K, V]{}
}

// InitWithBuffer makes h an empty map that stores its first len(buf) entries
// in buf; it grows into a new slice past that.
func (h *ObjectLinkedHashMap[K, V]) InitWithBuffer(buf []OlhmSlot[K, V]) {
	h.m = nil
	h.slots = buf[:0:len(buf)]
	h.size = 0
	h.first = 0
}

// find returns k's slot index (or -1) and its hash code.
func (h *ObjectLinkedHashMap[K, V]) find(k K) (int, int32) {
	// Go-only: performance — the ...Of forms do not box String keys.
	hash := ObjectHashCodeOf(k)
	if h == nil {
		return -1, hash
	}
	slots := h.slots
	if h.m != nil {
		t := h.m.table
		mask := uint(len(t) - 1)
		for b := lhmHash(hash) & mask; ; b = (b + 1) & mask {
			switch e := t[b]; e {
			case 0:
				return -1, hash
			case lhmTombstone:
			default:
				if slots[e-1].hash == hash && ObjectsEqualsOf(slots[e-1].key, k) {
					return int(e - 1), hash
				}
			}
		}
	}
	for i := int(h.first); i < len(slots); i++ {
		if slots[i].hash == hash && !slots[i].removed && ObjectsEqualsOf(slots[i].key, k) {
			return i, hash
		}
	}
	return -1, hash
}

func (h *ObjectLinkedHashMap[K, V]) Size() int {
	if h == nil {
		return 0
	}
	return int(h.size)
}
func (h *ObjectLinkedHashMap[K, V]) IsEmpty() bool { return h.Size() == 0 }

// Get is Map.get: the zero value (Java null) when absent.
func (h *ObjectLinkedHashMap[K, V]) Get(k K) V {
	v, _ := h.Lookup(k)
	return v
}

// Lookup returns the value and whether the key is present.
func (h *ObjectLinkedHashMap[K, V]) Lookup(k K) (V, bool) {
	if i, _ := h.find(k); i >= 0 {
		return h.slots[i].value, true
	}
	var zero V
	return zero, false
}

// ContainsKey is Map.containsKey.
func (h *ObjectLinkedHashMap[K, V]) ContainsKey(k K) bool {
	i, _ := h.find(k)
	return i >= 0
}

// Put is Map.put; re-putting an existing key keeps its position (and its
// original key object). Returns the previous value (zero if none).
func (h *ObjectLinkedHashMap[K, V]) Put(k K, v V) V {
	i, hash := h.find(k)
	if i >= 0 {
		old := h.slots[i].value
		h.slots[i].value = v
		return old
	}
	h.addWithHash(k, v, hash)
	var zero V
	return zero
}

// addWithHash appends an entry for k, which must be absent and have hash
// code hash.
func (h *ObjectLinkedHashMap[K, V]) addWithHash(k K, v V, hash int32) {
	if cap(h.slots) == 0 {
		h.slots = make([]olhmSlot[K, V], 0, 2)
		h.first = 0
	} else if len(h.slots) >= 64 && int(h.size)*4 < len(h.slots) {
		h.compact()
	}
	idx := len(h.slots)
	h.slots = append(h.slots, olhmSlot[K, V]{key: k, value: v, hash: hash})
	h.size++
	switch {
	case h.m != nil && !h.m.needsGrow():
		h.m.insert(lhmHash(hash), idx)
	case h.m != nil || int(h.size) > lhmSmallSize || len(h.slots) > 2*lhmSmallSize:
		h.reindex()
	}
}

// reindex (re)builds the index from the live slots, sized for them.
func (h *ObjectLinkedHashMap[K, V]) reindex() {
	if h.m == nil {
		h.m = &lhmIndex{}
	}
	h.m.reset(lhmTableLen(int(h.size)))
	for i := int(h.first); i < len(h.slots); i++ {
		if !h.slots[i].removed {
			h.m.insert(lhmHash(h.slots[i].hash), i)
		}
	}
}

// compact drops the removed slots, keeping the order of the others.
//
//go:noinline
func (h *ObjectLinkedHashMap[K, V]) compact() {
	w := 0
	for i := int(h.first); i < len(h.slots); i++ {
		if !h.slots[i].removed {
			h.slots[w] = h.slots[i]
			w++
		}
	}
	clear(h.slots[w:])
	h.slots = h.slots[:w]
	h.first = 0
	if h.m != nil {
		h.reindex()
	}
}

// Remove is Map.remove; returns the removed value (zero if absent).
func (h *ObjectLinkedHashMap[K, V]) Remove(k K) V {
	i, hash := h.find(k)
	if i < 0 {
		var zero V
		return zero
	}
	old := h.slots[i].value
	h.removeSlot(i, hash)
	return old
}

// removeSlot removes the entry at slot index i, whose hash code is hash.
func (h *ObjectLinkedHashMap[K, V]) removeSlot(i int, hash int32) {
	if h.m != nil {
		h.m.erase(lhmHash(hash), i)
	}
	h.size--
	h.slots[i] = olhmSlot[K, V]{removed: true}
	if i == int(h.first) {
		first := i + 1
		for first < len(h.slots) && h.slots[first].removed {
			first++
		}
		h.first = int32(first)
	}
}

// All iterates entries in insertion order.
func (h *ObjectLinkedHashMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if h == nil || h.size == 0 {
			return
		}
		for i := int(h.first); i < len(h.slots); i++ {
			if h.slots[i].removed {
				continue
			}
			if !yield(h.slots[i].key, h.slots[i].value) {
				return
			}
		}
	}
}

// Keys iterates keys in insertion order (see LinkedHashMap.Keys).
func (h *ObjectLinkedHashMap[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		if h == nil || h.size == 0 {
			return
		}
		for i := int(h.first); i < len(h.slots); i++ {
			if h.slots[i].removed {
				continue
			}
			if !yield(h.slots[i].key) {
				return
			}
		}
	}
}

// FirstKey returns the first non-removed key in insertion order.
func (h *ObjectLinkedHashMap[K, V]) FirstKey() (K, bool) {
	if h != nil && h.size > 0 {
		for i := int(h.first); i < len(h.slots); i++ {
			if !h.slots[i].removed {
				return h.slots[i].key, true
			}
		}
	}
	var zero K
	return zero, false
}

// Values iterates values in insertion order.
func (h *ObjectLinkedHashMap[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		if h == nil || h.size == 0 {
			return
		}
		for i := int(h.first); i < len(h.slots); i++ {
			if h.slots[i].removed {
				continue
			}
			if !yield(h.slots[i].value) {
				return
			}
		}
	}
}

// KeySlice returns the keys in insertion order.
func (h *ObjectLinkedHashMap[K, V]) KeySlice() []K {
	n := h.Size()
	if n == 0 {
		return nil
	}
	out := make([]K, 0, n)
	slots := h.slots
	if n == len(slots) {
		for i := range slots {
			out = append(out, slots[i].key)
		}
		return out
	}
	for i := int(h.first); i < len(slots); i++ {
		if !slots[i].removed {
			out = append(out, slots[i].key)
		}
	}
	return out
}

// CompactLen returns len(h.slots) when h has no tombstones, or -1 otherwise.
func (h *ObjectLinkedHashMap[K, V]) CompactLen() int {
	if h == nil {
		return 0
	}
	if int(h.size) == len(h.slots) {
		return len(h.slots)
	}
	return -1
}

// ValueAtSlot returns h.slots[i].value for 0 <= i < h.CompactLen().
func (h *ObjectLinkedHashMap[K, V]) ValueAtSlot(i int) V {
	return h.slots[i].value
}

func (h *ObjectLinkedHashMap[K, V]) lookupAny(k any) (any, bool) {
	key, ok := k.(K)
	if !ok {
		return nil, false
	}
	return h.Lookup(key)
}

// Equals is Map.equals(Object).
func (h *ObjectLinkedHashMap[K, V]) Equals(o any) bool {
	if any(h) == o {
		return true
	}
	other, ok := o.(javaMapEntries)
	if !ok || IsNil(o) || other.Size() != h.Size() {
		return false
	}
	for k, v := range h.All() {
		ov, present := other.lookupAny(k)
		if !present || !ObjectsEquals(v, ov) {
			return false
		}
	}
	return true
}

// HashCode is Map.hashCode().
func (h *ObjectLinkedHashMap[K, V]) HashCode() int32 {
	if h == nil || h.size == 0 {
		return 0
	}
	var sum int32
	for i := int(h.first); i < len(h.slots); i++ {
		if h.slots[i].removed {
			continue
		}
		sum += ObjectHashCodeOf(h.slots[i].key) ^ ObjectHashCodeOf(h.slots[i].value)
	}
	return sum
}

// String renders like AbstractMap.toString().
func (h *ObjectLinkedHashMap[K, V]) String() string {
	var sb strings.Builder
	sb.WriteByte('{')
	first := true
	for k, v := range h.All() {
		if !first {
			sb.WriteString(", ")
		}
		first = false
		sb.WriteString(ToJavaString(k))
		sb.WriteByte('=')
		sb.WriteString(ToJavaString(v))
	}
	sb.WriteByte('}')
	return sb.String()
}

// ObjectLinkedHashSet is java.util.LinkedHashSet (and Guava ImmutableSet)
// for elements whose classes override equals()/hashCode().
type ObjectLinkedHashSet[T any] struct {
	// Go-only: performance — held by value: one allocation per set.
	m ObjectLinkedHashMap[T, struct{}]
}

// NewObjectLinkedHashSet returns a set containing elems in order.
func NewObjectLinkedHashSet[T any](elems ...T) *ObjectLinkedHashSet[T] {
	s := &ObjectLinkedHashSet[T]{}
	for _, e := range elems {
		s.Add(e)
	}
	return s
}

// InitWithBuffer makes s an empty set that stores its first len(buf) elements
// in buf; it grows into a new slice past that.
func (s *ObjectLinkedHashSet[T]) InitWithBuffer(buf []OlhmSlot[T, struct{}]) {
	s.m.InitWithBuffer(buf)
}

func (s *ObjectLinkedHashSet[T]) Size() int {
	if s == nil {
		return 0
	}
	return s.m.Size()
}
func (s *ObjectLinkedHashSet[T]) IsEmpty() bool { return s.Size() == 0 }

// Add is Set.add; reports whether the set changed.
func (s *ObjectLinkedHashSet[T]) Add(e T) bool {
	i, hash := s.m.find(e)
	if i >= 0 {
		return false
	}
	s.m.addWithHash(e, struct{}{}, hash)
	return true
}

// AddAll adds every element of the sequence; reports whether the set changed.
func (s *ObjectLinkedHashSet[T]) AddAll(es iter.Seq[T]) bool {
	changed := false
	for e := range es {
		if s.Add(e) {
			changed = true
		}
	}
	return changed
}

// Contains is Set.contains.
func (s *ObjectLinkedHashSet[T]) Contains(e T) bool { return s != nil && s.m.ContainsKey(e) }

// Remove is Set.remove.
func (s *ObjectLinkedHashSet[T]) Remove(e T) bool {
	i, hash := s.m.find(e)
	if i < 0 {
		return false
	}
	s.m.removeSlot(i, hash)
	return true
}

// All iterates in insertion order (see LinkedHashMap.Keys).
func (s *ObjectLinkedHashSet[T]) All() iter.Seq[T] {
	var h *ObjectLinkedHashMap[T, struct{}]
	if s != nil {
		h = &s.m
	}
	return func(yield func(T) bool) {
		if h == nil || h.size == 0 {
			return
		}
		for i := int(h.first); i < len(h.slots); i++ {
			if h.slots[i].removed {
				continue
			}
			if !yield(h.slots[i].key) {
				return
			}
		}
	}
}

// ToSlice returns the elements in insertion order.
func (s *ObjectLinkedHashSet[T]) ToSlice() []T {
	if s == nil {
		return nil
	}
	return s.m.KeySlice()
}

func (s *ObjectLinkedHashSet[T]) containsAny(e any) bool {
	x, ok := e.(T)
	return ok && s.Contains(x)
}

// Equals is Set.equals(Object).
func (s *ObjectLinkedHashSet[T]) Equals(o any) bool {
	if any(s) == o {
		return true
	}
	other, ok := o.(javaSet)
	if !ok || IsNil(o) || other.Size() != s.Size() {
		return false
	}
	for e := range s.All() {
		if !other.containsAny(e) {
			return false
		}
	}
	return true
}

// HashCode is Set.hashCode().
func (s *ObjectLinkedHashSet[T]) HashCode() int32 {
	if s == nil || s.m.size == 0 {
		return 0
	}
	var h int32
	for i := int(s.m.first); i < len(s.m.slots); i++ {
		if s.m.slots[i].removed {
			continue
		}
		h += ObjectHashCodeOf(s.m.slots[i].key)
	}
	return h
}

// String renders like AbstractCollection.toString().
func (s *ObjectLinkedHashSet[T]) String() string { return CollectionToString(s.ToSlice()) }
