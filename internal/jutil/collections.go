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

// MapSource is the part of java.util.Map<K, V> that copy constructors and
// putAll read: the size and the entries in iteration order. *LinkedHashMap
// implements it.
type MapSource[K any, V any] interface {
	Size() int
	All() iter.Seq2[K, V]
}

var _ MapSource[string, int] = (*LinkedHashMap[string, int])(nil)

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

// ---------------------------------------------------------------------------
// Shared insertion-ordered slot table & open-addressing index

// lhmSlot is an entry slot of an orderedTable.
//
// Go-only: performance — storing entries in a flat contiguous slice instead
// of a doubly-linked list of heap-allocated nodes eliminates per-entry heap
// allocations and the 16-byte prev/next pointer overhead per entry. value
// precedes key because Go pads a struct that ends in a zero-size field
// (struct{} value of LinkedHashSet) by a word.
type lhmSlot[K comparable, V any] struct {
	value V
	key   K
}

// orderedTable is the insertion-ordered slot storage and lazy open-addressed
// index shared by LinkedHashMap and LinkedHashSet.
type orderedTable[K comparable, V any] struct {
	m     atomic.Pointer[lhmIndex] // nil until len(slots) exceeds lhmSmallSize or a non-front slot is removed
	slots []lhmSlot[K, V]
	size  int32
	first int32
}

// lhmIndex is the hash index of an orderedTable's slots.
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

// lhmSmallSize is the most entries a LinkedHashMap holds without an index.
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

// find probes the table starting at hash and returns the first live slot for
// which match returns true, or -1 if absent.
func (x *lhmIndex) find(hash uint, match func(slot int) bool) int {
	t := x.table
	mask := uint(len(t) - 1)
	for b := hash & mask; ; b = (b + 1) & mask {
		switch e := t[b]; e {
		case 0:
			return -1
		case lhmTombstone:
		default:
			if match(int(e - 1)) {
				return int(e - 1)
			}
		}
	}
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

func (t *orderedTable[K, V]) reset() {
	t.m.Store(nil)
	t.slots = nil
	t.size = 0
	t.first = 0
}

func (t *orderedTable[K, V]) clear() {
	clear(t.slots)
	t.slots = t.slots[:0]
	if m := t.m.Load(); m != nil {
		m.reset(len(m.table))
		clear(m.removed)
		m.removed = m.removed[:0]
	}
	t.size = 0
	t.first = 0
}

func (t *orderedTable[K, V]) loadIndex() *lhmIndex {
	return t.m.Load()
}

func (t *orderedTable[K, V]) removedSlice() []bool {
	if m := t.m.Load(); m != nil {
		return m.removed
	}
	return nil
}

func (t *orderedTable[K, V]) sizeOrZero() int {
	if t == nil {
		return 0
	}
	return int(t.size)
}

func (t *orderedTable[K, V]) capOrZero() int {
	if t == nil {
		return 0
	}
	return cap(t.slots)
}

func (t *orderedTable[K, V]) add(k K, v V) {
	if len(t.slots) < cap(t.slots) && t.m.Load() == nil {
		t.slots = append(t.slots, lhmSlot[K, V]{value: v, key: k})
		t.size++
		return
	}
	t.addSlow(k, v)
}

//go:noinline
func (t *orderedTable[K, V]) addSlow(k K, v V) {
	if cap(t.slots) == 0 {
		t.slots = make([]lhmSlot[K, V], 0, 2)
		t.first = 0
	} else if len(t.slots) >= 64 && int(t.size)*4 < len(t.slots) {
		t.compact()
	}
	idx := len(t.slots)
	t.slots = append(t.slots, lhmSlot[K, V]{value: v, key: k})
	t.size++
	if m := t.m.Load(); m != nil {
		if m.needsGrow() {
			t.reindex()
		} else {
			m.insert(lhmHash(t.slots[idx].key), idx)
		}
	}
}

//go:noinline
func (t *orderedTable[K, V]) buildIndex() {
	x := &lhmIndex{}
	x.reset(lhmTableLen(int(t.size)))
	for i := int(t.first); i < len(t.slots); i++ {
		x.insert(lhmHash(t.slots[i].key), i)
	}
	t.m.Store(x)
}

// buildIndexShared builds an index of t's live slots and publishes it
// atomically so concurrent read-only lookups on an unindexed map are race-free.
//
//go:noinline
func (t *orderedTable[K, V]) buildIndexShared() *lhmIndex {
	x := &lhmIndex{}
	x.reset(lhmTableLen(int(t.size)))
	for i := int(t.first); i < len(t.slots); i++ {
		x.insert(lhmHash(t.slots[i].key), i)
	}
	if t.m.CompareAndSwap(nil, x) {
		return x
	}
	return t.m.Load()
}

func (t *orderedTable[K, V]) reindex() {
	x := t.m.Load()
	x.reset(lhmTableLen(int(t.size)))
	for i := int(t.first); i < len(t.slots); i++ {
		if !x.isRemoved(i) {
			x.insert(lhmHash(t.slots[i].key), i)
		}
	}
}

//go:noinline
func (t *orderedTable[K, V]) compact() {
	m := t.m.Load()
	w := 0
	for i := int(t.first); i < len(t.slots); i++ {
		if m == nil || !m.isRemoved(i) {
			t.slots[w] = t.slots[i]
			w++
		}
	}
	clear(t.slots[w:])
	t.slots = t.slots[:w]
	t.first = 0
	if m != nil {
		clear(m.removed)
		m.removed = m.removed[:0]
		t.reindex()
	}
}

func (t *orderedTable[K, V]) removeSlot(i int) {
	if i != int(t.first) {
		t.removeSlotMiddle(i)
		return
	}
	m := t.m.Load()
	if m != nil {
		m.erase(lhmHash(t.slots[i].key), i)
	}
	t.size--
	t.slots[i] = lhmSlot[K, V]{}
	first := i + 1
	if m != nil {
		for first < len(t.slots) && m.isRemoved(first) {
			first++
		}
	}
	t.first = int32(first)
}

//go:noinline
func (t *orderedTable[K, V]) removeSlotMiddle(i int) {
	m := t.m.Load()
	if m == nil {
		t.buildIndex()
		m = t.m.Load()
	}
	m.erase(lhmHash(t.slots[i].key), i)
	t.size--
	t.slots[i] = lhmSlot[K, V]{}
	if i >= len(m.removed) {
		newLen := max(i+1, len(m.removed)*2, len(t.slots))
		newRemoved := make([]bool, newLen)
		copy(newRemoved, m.removed)
		m.removed = newRemoved
	}
	m.removed[i] = true
}

func (t *orderedTable[K, V]) all() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if t == nil {
			return
		}
		for i := int(t.first); i < len(t.slots); i++ {
			if m := t.m.Load(); i < int(t.first) || m != nil && m.isRemoved(i) {
				continue
			}
			if !yield(t.slots[i].key, t.slots[i].value) {
				return
			}
		}
	}
}

// Go-only: performance — keys and values walk the slot slice directly, like
// all, instead of ranging over all: nested iterator closures escape to the
// heap, whereas flat ones are inlined into the caller's loop.
func (t *orderedTable[K, V]) keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		if t == nil {
			return
		}
		for i := int(t.first); i < len(t.slots); i++ {
			if m := t.m.Load(); i < int(t.first) || m != nil && m.isRemoved(i) {
				continue
			}
			if !yield(t.slots[i].key) {
				return
			}
		}
	}
}

func (t *orderedTable[K, V]) values() iter.Seq[V] {
	return func(yield func(V) bool) {
		if t == nil {
			return
		}
		for i := int(t.first); i < len(t.slots); i++ {
			if m := t.m.Load(); i < int(t.first) || m != nil && m.isRemoved(i) {
				continue
			}
			if !yield(t.slots[i].value) {
				return
			}
		}
	}
}

func (t *orderedTable[K, V]) valueAt(index int) (V, bool) {
	var zero V
	if t == nil || index < 0 || index >= int(t.size) {
		return zero, false
	}
	m := t.m.Load()
	if t.first == 0 && (m == nil || len(m.removed) == 0) {
		return t.slots[index].value, true
	}
	cur := 0
	for i := int(t.first); i < len(t.slots); i++ {
		if m != nil && m.isRemoved(i) {
			continue
		}
		if cur == index {
			return t.slots[i].value, true
		}
		cur++
	}
	return zero, false
}

func (t *orderedTable[K, V]) firstKey() (K, bool) {
	if t != nil && t.size > 0 {
		m := t.m.Load()
		if m == nil || len(m.removed) == 0 {
			return t.slots[t.first].key, true
		}
		for i := int(t.first); i < len(t.slots); i++ {
			if !m.isRemoved(i) {
				return t.slots[i].key, true
			}
		}
	}
	var zero K
	return zero, false
}

func (t *orderedTable[K, V]) keySlice() []K {
	if t == nil || t.size == 0 {
		return nil
	}
	n := int(t.size)
	out := make([]K, 0, n)
	slots := t.slots
	if n == len(slots) {
		for i := range slots {
			out = append(out, slots[i].key)
		}
		return out
	}
	m := t.m.Load()
	for i := int(t.first); i < len(slots); i++ {
		if m == nil || !m.isRemoved(i) {
			out = append(out, slots[i].key)
		}
	}
	return out
}

// ---------------------------------------------------------------------------
// LinkedHashMap

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
	orderedTable[K, V]
}

func (h *LinkedHashMap[K, V]) table() *orderedTable[K, V] {
	if h == nil {
		return nil
	}
	return &h.orderedTable
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
			orderedTable: orderedTable[K, V]{
				slots: slices.Grow([]lhmSlot[K, V](nil), cap),
			},
		}
	}
	return &LinkedHashMap[K, V]{}
}

// Reset drops all entries and releases the backing slice and index.
func (h *LinkedHashMap[K, V]) Reset() { h.reset() }

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

func (h *LinkedHashMap[K, V]) lookupIn(m *lhmIndex, k K) int {
	slots := h.slots
	return m.find(lhmHash(k), func(slot int) bool {
		return slots[slot].key == k
	})
}

func (h *LinkedHashMap[K, V]) Size() int     { return h.table().sizeOrZero() }
func (h *LinkedHashMap[K, V]) IsEmpty() bool { return h.Size() == 0 }
func (h *LinkedHashMap[K, V]) Cap() int      { return h.table().capOrZero() }

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

// Clear is Map.clear.
func (h *LinkedHashMap[K, V]) Clear() { h.clear() }

// All iterates entries in insertion order. Entries removed during iteration
// are skipped; entries added during iteration are visited.
func (h *LinkedHashMap[K, V]) All() iter.Seq2[K, V] { return h.table().all() }

// ValueAt returns the index-th value in insertion order without allocating a slice.
func (h *LinkedHashMap[K, V]) ValueAt(index int) (V, bool) { return h.table().valueAt(index) }

// Keys iterates keys in insertion order (Java: map.keySet()).
func (h *LinkedHashMap[K, V]) Keys() iter.Seq[K] { return h.table().keys() }

// Values iterates values in insertion order (Java: map.values()).
func (h *LinkedHashMap[K, V]) Values() iter.Seq[V] { return h.table().values() }

// KeySlice returns the keys in insertion order.
func (h *LinkedHashMap[K, V]) KeySlice() []K { return h.table().keySlice() }

func formatMapString[K any, V any](seq iter.Seq2[K, V]) string {
	var sb strings.Builder
	sb.WriteByte('{')
	first := true
	for k, v := range seq {
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

// String renders like AbstractMap.toString(): "{k1=v1, k2=v2}".
func (h *LinkedHashMap[K, V]) String() string { return formatMapString(h.All()) }

// javaMapEntries is implemented by map types so that Map.equals can compare
// across implementations.
type javaMapEntries interface {
	Size() int
	lookupAny(k any) (any, bool)
}

func mapEntriesEquals[K any, V any](size int, seq iter.Seq2[K, V], o any) bool {
	other, ok := o.(javaMapEntries)
	if !ok || IsNil(o) || other.Size() != size {
		return false
	}
	for k, v := range seq {
		ov, present := other.lookupAny(k)
		if !present || !ObjectsEquals(v, ov) {
			return false
		}
	}
	return true
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
	return mapEntriesEquals(h.Size(), h.All(), o)
}

// HashCode is Map.hashCode(): the sum of key^value hash codes.
func (h *LinkedHashMap[K, V]) HashCode() int32 {
	if h == nil || h.size == 0 {
		return 0
	}
	var sum int32
	removed := h.removedSlice()
	for i := int(h.first); i < len(h.slots); i++ {
		if i < len(removed) && removed[i] {
			continue
		}
		sum += ObjectHashCode(any(h.slots[i].key)) ^ ObjectHashCode(any(h.slots[i].value))
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
		slots := make([]lhmSlot[K, V], n)
		copy(slots, src.slots)
		return &LinkedHashMap[K, V]{
			orderedTable: orderedTable[K, V]{
				slots: slots,
				size:  int32(n),
			},
		}
	}
	out := NewLinkedHashMapWithCapacity[K, V](n)
	for k, v := range src.All() {
		out.PutNew(k, v)
	}
	return out
}

// ---------------------------------------------------------------------------
// ImmutableMapBuilder

// ImmutableMapBuilder is Guava's ImmutableMap.Builder: entries keep their
// insertion order, and duplicate keys are rejected when the map is built.
type ImmutableMapBuilder[K comparable, V any] struct {
	slots []lhmSlot[K, V]
	buf   [4]lhmSlot[K, V]
	out   LinkedHashMap[K, V]
	built bool
}

// NewImmutableMapBuilder is ImmutableMap.builder().
func NewImmutableMapBuilder[K comparable, V any]() *ImmutableMapBuilder[K, V] {
	b := &ImmutableMapBuilder[K, V]{}
	b.slots = b.buf[:0]
	return b
}

// Put is ImmutableMap.Builder.put(K, V).
func (b *ImmutableMapBuilder[K, V]) Put(k K, v V) *ImmutableMapBuilder[K, V] {
	// Java: CollectPreconditions.checkEntryNotNull(key, value)
	if IsNil(k) {
		panic(NewNullPointerExceptionMsg("null key in entry: null=" + ToJavaString(v)))
	} else if IsNil(v) {
		panic(NewNullPointerExceptionMsg("null value in entry: " + ToJavaString(k) + "=null"))
	}
	if b.slots == nil {
		b.slots = b.buf[:0]
	}
	b.slots = append(b.slots, lhmSlot[K, V]{key: k, value: v})
	return b
}

// BuildOrThrow is ImmutableMap.Builder.buildOrThrow() (and build()).
// Guava's RegularImmutableMap.fromEntryArray inserts the entries from the
// last to the first, so the duplicate reported is the last entry (in
// insertion order) whose key is put again later, together with the nearest
// later entry with that key: IllegalArgumentException "Multiple entries
// with same key: k=later and k=earlier".
func (b *ImmutableMapBuilder[K, V]) BuildOrThrow() *LinkedHashMap[K, V] {
	n := len(b.slots)
	if n == 0 {
		if !b.built {
			b.built = true
			b.out.Reset()
			return &b.out
		}
		return &LinkedHashMap[K, V]{}
	}
	src := b.slots
	if n <= lhmSmallSize {
		for i := n - 1; i >= 0; i-- {
			k := src[i].key
			for j := i + 1; j < n; j++ {
				if src[j].key == k {
					panic(NewIllegalArgumentExceptionMsg("Multiple entries with same key: " +
						ToJavaString(k) + "=" + ToJavaString(src[j].value) + " and " +
						ToJavaString(k) + "=" + ToJavaString(src[i].value)))
				}
			}
		}
		if !b.built {
			b.built = true
			slots := src[:n:n]
			b.slots = slots
			b.out.m.Store(nil)
			b.out.slots = slots
			b.out.size = int32(n)
			b.out.first = 0
			return &b.out
		}
		slots := make([]lhmSlot[K, V], n)
		copy(slots, src)
		return &LinkedHashMap[K, V]{
			orderedTable: orderedTable[K, V]{
				slots: slots,
				size:  int32(n),
			},
		}
	}
	var out *LinkedHashMap[K, V]
	if !b.built {
		out = &b.out
	} else {
		out = &LinkedHashMap[K, V]{}
	}
	idx := &lhmIndex{table: make([]uint32, lhmTableLen(n))}
	out.m.Store(idx)
	out.slots = src
	out.size = int32(n)
	out.first = 0
	for i := n - 1; i >= 0; i-- {
		k := src[i].key
		if later := out.lookupIn(idx, k); later >= 0 {
			out.Reset()
			panic(NewIllegalArgumentExceptionMsg("Multiple entries with same key: " +
				ToJavaString(k) + "=" + ToJavaString(src[later].value) + " and " +
				ToJavaString(k) + "=" + ToJavaString(src[i].value)))
		}
		idx.insert(lhmHash(k), i)
	}
	if !b.built {
		b.built = true
		out.slots = src[:n:n]
		b.slots = out.slots
		return out
	}
	out.slots = make([]lhmSlot[K, V], n)
	copy(out.slots, src)
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

// ---------------------------------------------------------------------------
// OrderedSetMultimap

// OrderedSetView is an insertion-ordered, deduplicated read-only slice view of
// comparable elements.
type OrderedSetView[V comparable] []V

// Contains reports whether v is in s.
func (s OrderedSetView[V]) Contains(v V) bool { return slices.Contains(s, v) }

// ToSlice returns a defensive copy of s (non-nil empty slice when empty).
func (s OrderedSetView[V]) ToSlice() []V {
	if len(s) == 0 {
		return []V{}
	}
	return slices.Clone([]V(s))
}

// All iterates elements in insertion order.
func (s OrderedSetView[V]) All() iter.Seq[V] { return slices.Values([]V(s)) }

// Size returns the number of elements in s.
func (s OrderedSetView[V]) Size() int { return len(s) }

// IsEmpty reports whether s has no elements.
func (s OrderedSetView[V]) IsEmpty() bool { return len(s) == 0 }

// OrderedSetMultimap maps keys to insertion-ordered, deduplicated value sets,
// preserving first-seen key insertion order (Guava LinkedHashMultimap /
// ImmutableSetMultimap semantics). Its zero value is ready to use.
type OrderedSetMultimap[K, V comparable] struct {
	vals  map[K][]V
	order []K
	size  int
}

// NewOrderedSetMultimap returns an empty OrderedSetMultimap.
func NewOrderedSetMultimap[K, V comparable]() *OrderedSetMultimap[K, V] {
	return &OrderedSetMultimap[K, V]{}
}

// Put adds (k, v) if not already present and reports whether the multimap changed.
func (m *OrderedSetMultimap[K, V]) Put(k K, v V) bool {
	if m.vals == nil {
		m.vals = make(map[K][]V)
	}
	existing, ok := m.vals[k]
	if !ok {
		m.order = append(m.order, k)
	} else if slices.Contains(existing, v) {
		return false
	}
	m.vals[k] = append(existing, v)
	m.size++
	return true
}

// Get returns the ordered value set for k (empty view when absent).
func (m *OrderedSetMultimap[K, V]) Get(k K) OrderedSetView[V] {
	if m == nil {
		return nil
	}
	return OrderedSetView[V](m.vals[k])
}

// Size returns the total number of (key, value) pairs.
func (m *OrderedSetMultimap[K, V]) Size() int {
	if m == nil {
		return 0
	}
	return m.size
}

// IsEmpty reports whether m contains no entries.
func (m *OrderedSetMultimap[K, V]) IsEmpty() bool {
	return m == nil || m.size == 0
}

// Entries iterates (key, value) pairs in key and per-key insertion order.
func (m *OrderedSetMultimap[K, V]) Entries() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if m == nil {
			return
		}
		for _, k := range m.order {
			for _, v := range m.vals[k] {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

// AsMap iterates (key, valueSet) pairs in key insertion order.
func (m *OrderedSetMultimap[K, V]) AsMap() iter.Seq2[K, OrderedSetView[V]] {
	return func(yield func(K, OrderedSetView[V]) bool) {
		if m == nil {
			return
		}
		for _, k := range m.order {
			if !yield(k, OrderedSetView[V](m.vals[k])) {
				return
			}
		}
	}
}

// Clone returns an independent deep copy of m's key and value slices.
func (m *OrderedSetMultimap[K, V]) Clone() OrderedSetMultimap[K, V] {
	if m == nil || len(m.order) == 0 {
		return OrderedSetMultimap[K, V]{}
	}
	vals := make(map[K][]V, len(m.vals))
	for _, k := range m.order {
		vals[k] = slices.Clone(m.vals[k])
	}
	return OrderedSetMultimap[K, V]{
		vals:  vals,
		order: slices.Clone(m.order),
		size:  m.size,
	}
}

// Inverse returns a new multimap mapping each value to its associated keys in
// insertion order.
func (m *OrderedSetMultimap[K, V]) Inverse() *OrderedSetMultimap[V, K] {
	out := &OrderedSetMultimap[V, K]{}
	if m == nil || len(m.order) == 0 {
		return out
	}
	out.vals = make(map[V][]K, len(m.order))
	out.order = make([]V, 0, len(m.order))
	for _, k := range m.order {
		for _, v := range m.vals[k] {
			ks, ok := out.vals[v]
			if !ok {
				out.order = append(out.order, v)
			} else if slices.Contains(ks, k) {
				continue
			}
			out.vals[v] = append(ks, k)
			out.size++
		}
	}
	return out
}

// String formats m like Guava's AbstractMultimap.toString(): "{k1=[v1, v2]}".
func (m *OrderedSetMultimap[K, V]) String() string {
	if m == nil || len(m.order) == 0 {
		return "{}"
	}
	var sb strings.Builder
	sb.WriteByte('{')
	for i, k := range m.order {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(ToJavaString(k))
		sb.WriteString("=[")
		for j, v := range m.vals[k] {
			if j > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(ToJavaString(v))
		}
		sb.WriteByte(']')
	}
	sb.WriteByte('}')
	return sb.String()
}
