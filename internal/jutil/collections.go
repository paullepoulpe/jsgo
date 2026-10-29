package jutil

import (
	"hash/maphash"
	"iter"
	"slices"
	"strings"
	"sync/atomic"
)

// Go-only: language-adapter — java.util / Guava collections with Java
// iteration-order semantics. Ported code uses these (never a ranged Go map)
// wherever Java iteration order is observable.

// Collection is the part of java.util.Collection<T> that bulk operations
// read: the size, membership, and the elements in iteration order.
// *ArrayList, *LinkedHashSet, and *ObjectLinkedHashSet implement it.
type Collection[T any] interface {
	Size() int
	Contains(e T) bool
	All() iter.Seq[T]
}

// MapSource is the part of java.util.Map<K, V> that copy constructors and
// putAll read: the size and the entries in iteration order. *LinkedHashMap
// and *ObjectLinkedHashMap implement it.
type MapSource[K any, V any] interface {
	Size() int
	All() iter.Seq2[K, V]
}

var (
	_ Collection[string]     = (*ArrayList[string])(nil)
	_ Collection[string]     = (*LinkedHashSet[string])(nil)
	_ Collection[string]     = (*ObjectLinkedHashSet[string])(nil)
	_ MapSource[string, int] = (*LinkedHashMap[string, int])(nil)
	_ MapSource[string, int] = (*ObjectLinkedHashMap[string, int])(nil)
)

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

// ListHashCode is List.hashCode() for a slice standing in for a Java list.
func ListHashCode[T any](elems []T) int32 {
	var h int32 = 1
	for _, e := range elems {
		h = 31*h + ObjectHashCode(e)
	}
	return h
}

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
	m     atomic.Pointer[lhmIndex] // nil until len(slots) exceeds lhmSmallSize or a non-front slot is removed
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
	h.m.Store(nil)
	h.slots = buf[:0:len(buf)]
	h.size = 0
	h.first = 0
}

// Reset drops all entries and releases the backing slice and index.
func (h *LinkedHashMap[K, V]) Reset() {
	h.m.Store(nil)
	h.slots = nil
	h.size = 0
	h.first = 0
}

// EnsureCapacity grows h's backing slot slice so it can hold at least totalCap
// entries without reallocating.
func (h *LinkedHashMap[K, V]) EnsureCapacity(totalCap int) {
	if rem := totalCap - len(h.slots); rem > 0 {
		h.slots = slices.Grow(h.slots, rem)
	}
}

func (h *LinkedHashMap[K, V]) loadIndex() *lhmIndex {
	return h.m.Load()
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
	return h.lookupIn(h.m.Load(), k)
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

// Cap returns the current slot slice capacity of h.
func (h *LinkedHashMap[K, V]) Cap() int {
	if h == nil {
		return 0
	}
	return cap(h.slots)
}

// add appends an entry for k, which must be absent.
func (h *LinkedHashMap[K, V]) add(k K, v V) {
	if len(h.slots) < cap(h.slots) && h.m.Load() == nil {
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
	if m := h.m.Load(); m != nil {
		if m.needsGrow() {
			h.reindex()
		} else {
			m.insert(lhmHash(k), idx)
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
	h.m.Store(x)
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
	if h.m.CompareAndSwap(nil, x) {
		return x
	}
	return h.m.Load()
}

// reindex rebuilds the index table from the live slots, sized for them.
func (h *LinkedHashMap[K, V]) reindex() {
	x := h.m.Load()
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
	m := h.m.Load()
	w := 0
	for i := int(h.first); i < len(h.slots); i++ {
		if m == nil || !m.isRemoved(i) {
			h.slots[w] = h.slots[i]
			w++
		}
	}
	clear(h.slots[w:])
	h.slots = h.slots[:w]
	h.first = 0
	if m != nil {
		clear(m.removed)
		m.removed = m.removed[:0]
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
	m := h.m.Load()
	if m != nil {
		m.erase(lhmHash(h.slots[i].key), i)
	}
	h.size--
	h.slots[i] = LhmSlot[K, V]{}
	first := i + 1
	if m != nil {
		for first < len(h.slots) && m.isRemoved(first) {
			first++
		}
	}
	h.first = int32(first)
}

//go:noinline
func (h *LinkedHashMap[K, V]) removeSlotMiddle(i int) {
	m := h.m.Load()
	if m == nil {
		h.buildIndex()
		m = h.m.Load()
	}
	m.erase(lhmHash(h.slots[i].key), i)
	h.size--
	h.slots[i] = LhmSlot[K, V]{}
	if i >= len(m.removed) {
		newLen := max(i+1, len(m.removed)*2, len(h.slots))
		newRemoved := make([]bool, newLen)
		copy(newRemoved, m.removed)
		m.removed = newRemoved
	}
	m.removed[i] = true
}

// Clear is Map.clear.
func (h *LinkedHashMap[K, V]) Clear() {
	clear(h.slots)
	h.slots = h.slots[:0]
	if m := h.m.Load(); m != nil {
		m.reset(len(m.table))
		clear(m.removed)
		m.removed = m.removed[:0]
	}
	h.size = 0
	h.first = 0
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

// Go-only: language-adapter — Arrays.binarySearch(int[], int).

// ArraysBinarySearchInts is Arrays.binarySearch(int[] a, int key) on a
// sorted slice: the index of key, or -(insertion point) - 1 when key is
// absent, the insertion point being the index of the first element greater
// than key (len(a) if there is none), as the Javadoc specifies.
//
// The Javadoc leaves open which index is returned when key occurs several
// times. Observed behaviour (black-box): it is the one the textbook
// bisection finds, probing the middle index (low+high)/2 (rounded down) of
// the remaining range [low, high] and returning on the first equal element;
// checked against 200,000 random sorted arrays with many duplicates.
func ArraysBinarySearchInts(a []int, key int) int {
	low, high := 0, len(a)-1
	for low <= high {
		mid := int(uint(low+high) >> 1)
		switch v := a[mid]; {
		case v < key:
			low = mid + 1
		case v > key:
			high = mid - 1
		default:
			return mid
		}
	}
	return -(low + 1)
}
