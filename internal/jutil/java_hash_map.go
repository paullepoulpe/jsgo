package jutil

import (
	"iter"
	"math"
	"reflect"
	"strconv"
	"strings"
)

// Go-only: language-adapter — java.util.HashMap as JDK 21 behaves, written
// from its public specification and black-box experiments only (PORTING.md:
// no JDK code). Callers depend on its iteration order bit for bit, so every
// rule that decides where an entry sits is reproduced.
//
// Specified by the Javadoc of java.util.HashMap and java.util.Map:
//   - the table has a power-of-two number of buckets (the "capacity"); the
//     default capacity is 16 and the default load factor 0.75;
//   - when the number of entries exceeds capacity * load factor the table is
//     rehashed into about twice as many buckets;
//   - HashMap(Map) and HashMap.newHashMap(int) size the table for the given
//     number of mappings at load factor 0.75;
//   - keys whose classes are Comparable may be compared to order a bucket that
//     holds many colliding keys;
//   - iterators are fail-fast (ConcurrentModificationException), and the
//     compute/merge methods throw it when their function modifies the map;
//   - equals/hashCode/toString follow Map and AbstractMap.
//
// Observed behaviour (black-box, JDK 21, public API only; the experiments
// print the iteration order of crafted key sets):
//   - a key with hashCode h goes to bucket (h ^ (h >>> 16)) & (capacity - 1);
//     iteration visits the buckets in index order;
//   - put and putIfAbsent append a new key at the tail of its bucket and then,
//     if the size exceeds the threshold, double the table; computeIfAbsent,
//     compute and merge instead insert a new key at the head of its bucket and
//     never grow the table afterwards: they grow it first, before looking the
//     key up, when the size already exceeds the threshold (or there is no
//     table yet);
//   - doubling keeps the relative order of the keys of a bucket in both halves
//     it splits into; the threshold is recomputed as capacity * load factor
//     (truncated) while the old capacity is below 16, and doubled from there
//     on; the first table of a map without an initial capacity has 16 buckets
//     and threshold 12, whatever the load factor;
//   - the table is allocated by the first insertion; removals and clear()
//     never shrink it;
//   - putAll/HashMap(Map) on a map without a table size it for
//     ceil(n / loadFactor) mappings (but never below its initial capacity);
//     on a map with a table they double it until n fits the threshold, then
//     insert one by one; clone() re-inserts the entries, in iteration order,
//     into a fresh map with the same load factor sized the same way;
//   - a bucket becomes a "sorted bucket" when a put makes it hold 9 keys
//     (8 for computeIfAbsent/compute/merge), if the table has at least 64
//     buckets; with fewer buckets the table is doubled instead;
//   - a sorted bucket is a red-black tree (the textbook algorithms: insertion
//     fix-up, deletion through the in-order successor) ordered by spread
//     hash (signed), then, for two keys of one class C implementing
//     Comparable<C>, by compareTo, and otherwise by class name and finally
//     by identity hash code. Its iteration order is a separate chain: when
//     the bucket is converted the chain keeps the bucket's order, a new key is
//     linked right after its parent in the tree, a removed key is unlinked,
//     and after every insertion or removal (except removal through an
//     iterator) the root of the tree is moved to the head of the chain;
//   - a removal (not through an iterator) from a tree whose root lacks a
//     child, or whose left child lacks a left child, turns the bucket back
//     into a plain list (the chain order is kept), before the key is removed;
//   - when the table doubles, a sorted bucket splits into its two halves
//     (chain order kept): a half of at most 6 keys becomes a list; a half
//     of more keys keeps the tree if the other half is empty, and otherwise
//     is rebuilt as a new tree from its chain order (root moved to the head);
//   - iterators: next() on an exhausted iterator of a modified map throws
//     ConcurrentModificationException, not NoSuchElementException; remove()
//     with nothing to remove throws IllegalStateException even when the map
//     was modified.

const (
	javaHashMapDefaultCapacity = 16
	javaHashMapMaxCapacity     = 1 << 30
	javaHashMapDefaultLoad     = float32(0.75)
	// A list bucket that a put makes reach this many keys is converted to a tree
	// (or the table grows); computeIfAbsent/compute/merge convert at one less.
	javaHashMapTreeifyAt = 9
	// The smallest table whose buckets are converted to trees.
	javaHashMapMinTreeTable = 64
	// A half of a split sorted bucket with at most this many keys becomes a list.
	javaHashMapListBelow = 6
)

// hashMapNode is one mapping. The nodes of a bucket form a singly linked
// chain (next) in iteration order. While their bucket is sorted, the nodes
// also carry tree links.
type hashMapNode[K any, V any] struct {
	hash  int32 // the spread hash code of key
	key   K
	value V
	next  *hashMapNode[K, V]
	rb    *hashTreeLinks[K, V] // non-nil exactly while the bucket is sorted
}

// hashTreeLinks hold what a node needs only while its bucket is sorted:
// its place in the red-black tree, its colour, and a back pointer into the
// iteration chain, so that unlinking a node or promoting the tree's root to
// the head of the chain does not have to walk the chain.
type hashTreeLinks[K any, V any] struct {
	parent, left, right *hashMapNode[K, V]
	prev                *hashMapNode[K, V]
	red                 bool
}

// JavaHashMap is java.util.HashMap: keys hash with ObjectHashCode and
// compare with ObjectsEquals; iteration order is exactly the JDK's.
type JavaHashMap[K any, V any] struct {
	table      []*hashMapNode[K, V] // nil until the first insertion
	size       int
	loadFactor float32
	// firstBins is the number of buckets the first table will get: set from
	// an initial capacity or by a bulk insertion into a map with no table
	// yet; 0 asks for the default.
	firstBins int
	// threshold is the largest size the current table may hold; one more
	// entry makes it grow. Unused while table is nil.
	threshold int
	modCount  int
}

// javaHashMapHash is the spread hash HashMap places a key by: the key's
// hashCode with its high half folded into the low half.
func javaHashMapHash(key any) int32 {
	h := ObjectHashCode(key)
	return h ^ int32(uint32(h)>>16)
}

// hashMapHashOf is javaHashMapHash without boxing the common key types.
func hashMapHashOf[K any](key K) int32 {
	var h int32
	switch k := any(key).(type) {
	case string:
		h = StringHashCode(k)
	case int32:
		h = k
	case int:
		h = int32(k)
	case int64:
		h = int32(k ^ int64(uint64(k)>>32))
	default:
		return javaHashMapHash(any(key))
	}
	return h ^ int32(uint32(h)>>16)
}

// hashMapKeysEqual is key.equals(other) for a non-null probe key, and
// identity for null (ObjectsEquals), without boxing the common key types.
func hashMapKeysEqual[K any](key, other K) bool {
	switch k := any(key).(type) {
	case string:
		o, ok := any(other).(string)
		return ok && k == o
	case int32:
		o, ok := any(other).(int32)
		return ok && k == o
	case int:
		o, ok := any(other).(int)
		return ok && k == o
	case int64:
		o, ok := any(other).(int64)
		return ok && k == o
	}
	return ObjectsEquals(any(key), any(other))
}

// pow2AtLeast turns a requested capacity into a bucket count: c rounded up
// to a power of two, and kept within [1, 2^30]. (The Javadoc only says the
// capacity is a power of two; the rounding and the bounds were observed.)
func pow2AtLeast(c int32) int32 {
	if c >= javaHashMapMaxCapacity {
		return javaHashMapMaxCapacity
	}
	p := int32(1)
	for p < c {
		p *= 2
	}
	return p
}

// capacityForMappings is the initial capacity that HashMap.newHashMap (and
// HashSet's copy constructor) request for n mappings at the default load
// factor 3/4: the least integer >= 4n/3, capped at Integer.MAX_VALUE.
// Observed: 12 mappings give a 16-bucket table, 13 give 32.
func capacityForMappings(n int) int {
	return int(min((4*int64(n)+2)/3, math.MaxInt32))
}

// NewJavaHashMap is new HashMap<>().
func NewJavaHashMap[K any, V any]() *JavaHashMap[K, V] {
	return &JavaHashMap[K, V]{loadFactor: javaHashMapDefaultLoad}
}

// NewJavaHashMapWithCapacity is new HashMap<>(initialCapacity).
func NewJavaHashMapWithCapacity[K any, V any](initialCapacity int) *JavaHashMap[K, V] {
	return NewJavaHashMapWithCapacityAndLoadFactor[K, V](initialCapacity, javaHashMapDefaultLoad)
}

// NewJavaHashMapWithCapacityAndLoadFactor is new HashMap<>(initialCapacity,
// loadFactor). The messages of the IllegalArgumentExceptions are the JDK's
// (observed).
func NewJavaHashMapWithCapacityAndLoadFactor[K any, V any](initialCapacity int, loadFactor float32) *JavaHashMap[K, V] {
	if initialCapacity < 0 {
		panic(NewIllegalArgumentExceptionMsg("Illegal initial capacity: " + strconv.Itoa(initialCapacity)))
	}
	if initialCapacity > javaHashMapMaxCapacity {
		initialCapacity = javaHashMapMaxCapacity
	}
	if !(loadFactor > 0) { // <= 0 or NaN
		panic(NewIllegalArgumentExceptionMsg("Illegal load factor: " + javaFloatString(loadFactor)))
	}
	return &JavaHashMap[K, V]{
		loadFactor: loadFactor,
		firstBins:  int(pow2AtLeast(int32(initialCapacity))),
	}
}

// javaFloatString renders a float like Float.toString for the values a load
// factor takes (shortest digits of the float, printed as a double).
func javaFloatString(f float32) string {
	d, _ := strconv.ParseFloat(strconv.FormatFloat(float64(f), 'g', -1, 32), 64)
	if f != f {
		d = math.NaN()
	}
	return DoubleToString(d)
}

// NewJavaHashMapFrom is new HashMap<>(m): load factor 0.75, a table sized
// for m's mappings, then m's entries in its iteration order.
func NewJavaHashMapFrom[K any, V any](m MapSource[K, V]) *JavaHashMap[K, V] {
	h := NewJavaHashMap[K, V]()
	h.bulkInsert(m)
	return h
}

// JavaHashMap_NewHashMap is HashMap.newHashMap(numMappings).
func JavaHashMap_NewHashMap[K any, V any](numMappings int) *JavaHashMap[K, V] {
	if numMappings < 0 {
		panic(NewIllegalArgumentExceptionMsg("Negative number of mappings: " + strconv.Itoa(numMappings)))
	}
	return NewJavaHashMapWithCapacity[K, V](capacityForMappings(numMappings))
}

// bulkInsert adds the mappings of m one by one, in m's iteration order: it
// is putAll, the HashMap(Map) constructor and clone(). Before inserting it
// makes room for all of them at once (header bullet on putAll): a map that
// has no table yet will allocate one with bucketsFor(n) buckets, unless it
// was already asked for more; a map that has a table doubles it until n
// entries fit under the threshold.
func (h *JavaHashMap[K, V]) bulkInsert(m MapSource[K, V]) {
	n := m.Size()
	if n == 0 {
		return
	}
	if h.table == nil {
		h.firstBins = max(h.firstBins, h.bucketsFor(n))
	} else {
		for n > h.threshold && len(h.table) < javaHashMapMaxCapacity {
			h.grow()
		}
	}
	for k, v := range m.All() {
		h.insertKey(hashMapHashOf(k), k, v, false)
	}
}

// bucketsFor is the bucket count a first table needs for n mappings at this
// map's load factor: ceil(n / loadFactor), at most 2^30, rounded up to a
// power of two. Observed: putAll of 6 mappings into new HashMap<>() gives 8
// buckets and 12 mappings give 16 (so it is the ceiling, not +1).
// The quotient is taken in float64; no n below 2^24 gives a different
// bucket count than float32 arithmetic would (checked numerically).
func (h *JavaHashMap[K, V]) bucketsFor(n int) int {
	need := math.Ceil(float64(n) / float64(h.loadFactor))
	if need >= javaHashMapMaxCapacity {
		return javaHashMapMaxCapacity
	}
	return int(pow2AtLeast(int32(need)))
}

// grow gives the map its first table, or a table with twice as many
// buckets, and returns it; at 2^30 buckets it only lifts the threshold to
// Integer.MAX_VALUE. On doubling, the keys of bucket j go to bucket j or to
// bucket j + (old bucket count), whichever their hash selects, and each of
// the two keeps the keys in the order they had (header bullet).
func (h *JavaHashMap[K, V]) grow() []*hashMapNode[K, V] {
	prev := h.table
	prevBins := len(prev)
	if prevBins >= javaHashMapMaxCapacity {
		h.threshold = math.MaxInt32
		return prev
	}
	bins := javaHashMapDefaultCapacity
	if prevBins > 0 {
		bins = 2 * prevBins
	} else if h.firstBins > 0 {
		bins = h.firstBins
	}
	h.threshold = h.nextThreshold(prevBins, bins)
	tab := make([]*hashMapNode[K, V], bins)
	h.table = tab
	for j, head := range prev {
		switch {
		case head == nil:
		case head.rb != nil:
			h.splitSortedBucket(tab, head, j, prevBins)
		default:
			splitListBucket(tab, head, j, prevBins)
		}
	}
	return tab
}

// nextThreshold is the threshold of the new table (of bins buckets) that grow
// installs in place of one of prevBins buckets (0: there was none).
// Observed (header bullets): the default first table (16 buckets) always
// gets 12, whatever the load factor; a table grown from 16 or more buckets
// to fewer than 2^30 gets twice the old threshold, in 32-bit arithmetic (a
// saturated threshold wraps to a negative one, after which every insertion
// grows the table again); every other table gets thresholdFor(bins).
func (h *JavaHashMap[K, V]) nextThreshold(prevBins, bins int) int {
	switch {
	case prevBins == 0 && h.firstBins == 0:
		return javaHashMapDefaultCapacity / 4 * 3
	case prevBins >= javaHashMapDefaultCapacity && bins < javaHashMapMaxCapacity:
		return int(int32(h.threshold) * 2)
	}
	return h.thresholdFor(bins)
}

// thresholdFor is bins * loadFactor rounded toward zero, or
// Integer.MAX_VALUE once bins or that product reaches 2^30. bins is a power
// of two, so the float64 product is exact.
func (h *JavaHashMap[K, V]) thresholdFor(bins int) int {
	limit := float64(bins) * float64(h.loadFactor)
	if bins >= javaHashMapMaxCapacity || limit >= javaHashMapMaxCapacity {
		return math.MaxInt32
	}
	return int(limit)
}

// splitListBucket moves the keys of list bucket j of a table of prevBins
// buckets into buckets j and j+prevBins of the doubled table tab, by the one
// hash bit the doubled mask adds. Both sides keep the keys' relative order.
func splitListBucket[K any, V any](tab []*hashMapNode[K, V], head *hashMapNode[K, V], j, prevBins int) {
	var heads, tails [2]*hashMapNode[K, V]
	for x := head; x != nil; {
		following := x.next
		x.next = nil
		side := 0
		if int(x.hash)&prevBins != 0 {
			side = 1
		}
		if tails[side] == nil {
			heads[side] = x
		} else {
			tails[side].next = x
		}
		tails[side] = x
		x = following
	}
	tab[j], tab[j+prevBins] = heads[0], heads[1]
}

// Size is Map.size().
func (h *JavaHashMap[K, V]) Size() int {
	if h == nil {
		return 0
	}
	return h.size
}

// IsEmpty is Map.isEmpty().
func (h *JavaHashMap[K, V]) IsEmpty() bool { return h.Size() == 0 }

// lookupNode finds the node for key, or nil.
func (h *JavaHashMap[K, V]) lookupNode(hash int32, key K) *hashMapNode[K, V] {
	if h == nil || len(h.table) == 0 {
		return nil
	}
	first := h.table[int(hash)&(len(h.table)-1)]
	if first == nil {
		return nil
	}
	if first.rb != nil {
		return first.root().searchTree(hash, key, javaComparableClassOf(any(key)))
	}
	for e := first; e != nil; e = e.next {
		if e.hash == hash && hashMapKeysEqual(key, e.key) {
			return e
		}
	}
	return nil
}

// Get is Map.get: the zero value (Java null) when absent.
func (h *JavaHashMap[K, V]) Get(key K) V {
	if e := h.lookupNode(hashMapHashOf(key), key); e != nil {
		return e.value
	}
	var zero V
	return zero
}

// Lookup returns the value for key and whether the key is present.
func (h *JavaHashMap[K, V]) Lookup(key K) (V, bool) {
	if e := h.lookupNode(hashMapHashOf(key), key); e != nil {
		return e.value, true
	}
	var zero V
	return zero, false
}

// GetOrDefault is Map.getOrDefault: the value of a present key (even null),
// else defaultValue.
func (h *JavaHashMap[K, V]) GetOrDefault(key K, defaultValue V) V {
	if e := h.lookupNode(hashMapHashOf(key), key); e != nil {
		return e.value
	}
	return defaultValue
}

// ContainsKey is Map.containsKey.
func (h *JavaHashMap[K, V]) ContainsKey(key K) bool {
	return h.lookupNode(hashMapHashOf(key), key) != nil
}

// ContainsValue is Map.containsValue: value.equals(v) for some value v.
func (h *JavaHashMap[K, V]) ContainsValue(value V) bool {
	if h == nil {
		return false
	}
	for _, e := range h.table {
		for ; e != nil; e = e.next {
			if ObjectsEquals(any(value), any(e.value)) {
				return true
			}
		}
	}
	return false
}

// isNullValue reports whether v is Java null.
func isNullValue[V any](v V) bool { return IsNil(any(v)) }

// Put is Map.put; returns the previous value (zero if none).
func (h *JavaHashMap[K, V]) Put(key K, value V) V {
	old, _ := h.insertKey(hashMapHashOf(key), key, value, false)
	return old
}

// PutIfAbsent is Map.putIfAbsent: maps key to value unless it is mapped to
// a non-null value; returns the previous value.
func (h *JavaHashMap[K, V]) PutIfAbsent(key K, value V) V {
	old, _ := h.insertKey(hashMapHashOf(key), key, value, true)
	return old
}

// insertKey is put (keepNonNull false) and putIfAbsent (keepNonNull true).
// A key already present stays where it is, and gets the new value unless
// keepNonNull is set and its current value is non-null; a new key goes at
// the end of its bucket's list (or into its sorted bucket). It returns the
// previous value and whether the key was present.
func (h *JavaHashMap[K, V]) insertKey(hash int32, key K, value V, keepNonNull bool) (V, bool) {
	buckets := h.table
	if len(buckets) == 0 {
		buckets = h.grow()
	}
	i := int(hash) & (len(buckets) - 1)
	var hit *hashMapNode[K, V]
	switch head := buckets[i]; {
	case head == nil:
		buckets[i] = &hashMapNode[K, V]{hash: hash, key: key, value: value}
	case head.rb != nil:
		hit = h.treePut(i, hash, key, value)
	default:
		walked := 0
		tail := head
		for x := head; x != nil; x = x.next {
			if x.hash == hash && hashMapKeysEqual(key, x.key) {
				hit = x
				break
			}
			walked++
			tail = x
		}
		if hit == nil {
			tail.next = &hashMapNode[K, V]{hash: hash, key: key, value: value}
			if walked+1 >= javaHashMapTreeifyAt {
				h.growOrConvertBucket(i)
			}
		}
	}
	if hit != nil {
		prior := hit.value
		if !keepNonNull || isNullValue(prior) {
			hit.value = value
		}
		return prior, true
	}
	h.modCount++
	h.size++
	if h.size > h.threshold {
		h.grow()
	}
	var zero V
	return zero, false
}

// PutAll is Map.putAll.
func (h *JavaHashMap[K, V]) PutAll(m MapSource[K, V]) { h.bulkInsert(m) }

// Remove is Map.remove(Object): returns the removed value (zero if none).
func (h *JavaHashMap[K, V]) Remove(key K) V {
	var zero V
	if e := h.deleteKey(hashMapHashOf(key), key, zero, false, false); e != nil {
		return e.value
	}
	return zero
}

// RemoveEntry is Map.remove(Object, Object): removes key only while it is
// mapped to value.
func (h *JavaHashMap[K, V]) RemoveEntry(key K, value V) bool {
	return h.deleteKey(hashMapHashOf(key), key, value, true, false) != nil
}

// deleteKey takes key's node out of the map and returns it (nil if key is
// absent). With valueMustMatch it does so only while the key is mapped to
// value. byIterator marks a removal through an iterator: it leaves the
// order of the other keys of a sorted bucket alone (the root is not
// promoted and the bucket is not turned back into a list).
func (h *JavaHashMap[K, V]) deleteKey(hash int32, key K, value V, valueMustMatch, byIterator bool) *hashMapNode[K, V] {
	if h == nil || len(h.table) == 0 {
		return nil
	}
	i := int(hash) & (len(h.table) - 1)
	head := h.table[i]
	if head == nil {
		return nil
	}
	var node, before *hashMapNode[K, V]
	if head.rb != nil {
		node = head.root().searchTree(hash, key, javaComparableClassOf(any(key)))
	} else {
		for x := head; x != nil; before, x = x, x.next {
			if x.hash == hash && hashMapKeysEqual(key, x.key) {
				node = x
				break
			}
		}
	}
	if node == nil || valueMustMatch && !ObjectsEquals(any(value), any(node.value)) {
		return nil
	}
	if node.rb != nil {
		h.treeRemove(i, node, byIterator)
	} else if before == nil {
		h.table[i] = node.next
	} else {
		before.next = node.next
	}
	h.modCount++
	h.size--
	return node
}

// Clear is Map.clear(); the table keeps its size (observed).
func (h *JavaHashMap[K, V]) Clear() {
	h.modCount++
	if h.size > 0 {
		h.size = 0
		clear(h.table)
	}
}

// Replace is Map.replace(K, V): replaces the value of a present key;
// returns the previous value.
func (h *JavaHashMap[K, V]) Replace(key K, value V) V {
	if e := h.lookupNode(hashMapHashOf(key), key); e != nil {
		old := e.value
		e.value = value
		return old
	}
	var zero V
	return zero
}

// prepareCompute is the start of computeIfAbsent/compute/merge (observed):
// grow the table if the size exceeds the threshold (or there is no table),
// then find key's node, also returning its bucket index and the number of
// nodes of a list bucket.
func (h *JavaHashMap[K, V]) prepareCompute(hash int32, key K) (i int, node *hashMapNode[K, V], count int) {
	if h.size > h.threshold || len(h.table) == 0 {
		h.grow()
	}
	i = int(hash) & (len(h.table) - 1)
	first := h.table[i]
	if first == nil {
		return i, nil, 0
	}
	if first.rb != nil {
		return i, first.root().searchTree(hash, key, javaComparableClassOf(any(key))), 0
	}
	for e := first; e != nil; e = e.next {
		if e.hash == hash && hashMapKeysEqual(key, e.key) {
			return i, e, 0
		}
		count++
	}
	return i, nil, count
}

// insertComputed adds a new key for computeIfAbsent/compute/merge: into the
// tree of a sorted bucket, else at the head of the list (observed),
// converting a list that reaches 8 keys.
func (h *JavaHashMap[K, V]) insertComputed(i int, hash int32, key K, value V, count int) {
	first := h.table[i]
	if first != nil && first.rb != nil {
		h.treePut(i, hash, key, value)
	} else {
		h.table[i] = &hashMapNode[K, V]{hash: hash, key: key, value: value, next: first}
		if count+1 >= javaHashMapTreeifyAt-1 {
			h.growOrConvertBucket(i)
		}
	}
	h.modCount++
	h.size++
}

// ComputeIfAbsent is Map.computeIfAbsent: if key is absent or mapped to
// null, maps it to f(key) unless that is null; returns the current value.
func (h *JavaHashMap[K, V]) ComputeIfAbsent(key K, f func(K) V) V {
	if f == nil {
		panic(NewNullPointerException())
	}
	hash := hashMapHashOf(key)
	i, node, count := h.prepareCompute(hash, key)
	if node != nil && !isNullValue(node.value) {
		return node.value
	}
	mc := h.modCount
	v := f(key)
	if h.modCount != mc {
		panic(NewConcurrentModificationException())
	}
	if isNullValue(v) {
		return v
	}
	if node != nil {
		node.value = v
		return v
	}
	h.insertComputed(i, hash, key, v, count)
	return v
}

// ComputeIfPresent is Map.computeIfPresent: for a key mapped to a non-null
// value, maps it to f(key, value), or removes it if that is null.
func (h *JavaHashMap[K, V]) ComputeIfPresent(key K, f func(K, V) V) V {
	if f == nil {
		panic(NewNullPointerException())
	}
	hash := hashMapHashOf(key)
	e := h.lookupNode(hash, key)
	if e == nil || isNullValue(e.value) {
		var zero V
		return zero
	}
	mc := h.modCount
	v := f(key, e.value)
	if h.modCount != mc {
		panic(NewConcurrentModificationException())
	}
	if isNullValue(v) {
		var zero V
		h.deleteKey(hash, key, zero, false, false)
	} else {
		e.value = v
	}
	return v
}

// Compute is Map.compute: maps key to f(key, current value or null), or
// removes it if that is null.
func (h *JavaHashMap[K, V]) Compute(key K, f func(K, V) V) V {
	if f == nil {
		panic(NewNullPointerException())
	}
	hash := hashMapHashOf(key)
	i, node, count := h.prepareCompute(hash, key)
	var old V
	if node != nil {
		old = node.value
	}
	mc := h.modCount
	v := f(key, old)
	if h.modCount != mc {
		panic(NewConcurrentModificationException())
	}
	switch {
	case node != nil && !isNullValue(v):
		node.value = v
	case node != nil:
		var zero V
		h.deleteKey(hash, key, zero, false, false)
	case !isNullValue(v):
		h.insertComputed(i, hash, key, v, count)
	}
	return v
}

// Merge is Map.merge: maps an absent (or null-mapped) key to value, else to
// f(old, value), removing it if that is null.
func (h *JavaHashMap[K, V]) Merge(key K, value V, f func(V, V) V) V {
	if isNullValue(value) || f == nil {
		panic(NewNullPointerException())
	}
	hash := hashMapHashOf(key)
	i, node, count := h.prepareCompute(hash, key)
	if node == nil {
		h.insertComputed(i, hash, key, value, count)
		return value
	}
	v := value
	if !isNullValue(node.value) {
		mc := h.modCount
		v = f(node.value, value)
		if h.modCount != mc {
			panic(NewConcurrentModificationException())
		}
	}
	if isNullValue(v) {
		var zero V
		h.deleteKey(hash, key, zero, false, false)
	} else {
		node.value = v
	}
	return v
}

// ForEach is Map.forEach: visits the whole table, then throws
// ConcurrentModificationException if the action changed the map's structure
// (observed: the entries added meanwhile in later buckets are visited too).
func (h *JavaHashMap[K, V]) ForEach(action func(K, V)) {
	if action == nil {
		panic(NewNullPointerException())
	}
	if h.size == 0 {
		return
	}
	mc := h.modCount
	for _, e := range h.table {
		for ; e != nil; e = e.next {
			action(e.key, e.value)
		}
	}
	if h.modCount != mc {
		panic(NewConcurrentModificationException())
	}
}

// ReplaceAll is Map.replaceAll: every value v becomes f(key, v).
func (h *JavaHashMap[K, V]) ReplaceAll(f func(K, V) V) {
	if f == nil {
		panic(NewNullPointerException())
	}
	if h.size == 0 {
		return
	}
	mc := h.modCount
	for _, e := range h.table {
		for ; e != nil; e = e.next {
			e.value = f(e.key, e.value)
		}
	}
	if h.modCount != mc {
		panic(NewConcurrentModificationException())
	}
}

// Clone is HashMap.clone(): a shallow copy with the same load factor, its
// table sized for the current size (not copied from this map; observed).
func (h *JavaHashMap[K, V]) Clone() *JavaHashMap[K, V] {
	c := &JavaHashMap[K, V]{loadFactor: h.loadFactor}
	c.bulkInsert(h)
	return c
}

// JavaHashMapIterator is an iterator over the entries (entrySet(),
// keySet() and values() iterators share it).
type JavaHashMapIterator[K any, V any] struct {
	m *JavaHashMap[K, V]
	// pending is the entry the next call to Next returns; nil once the
	// iteration is over. It is looked up as soon as the previous entry is
	// returned.
	pending *hashMapNode[K, V]
	// lastReturned is the entry Remove and SetValue act on: the one Next
	// returned last, nil before the first Next and after a Remove.
	lastReturned *hashMapNode[K, V]
	nextBucket   int // where the scan for the next non-empty bucket resumes
	expectedMods int // the map's modCount as of this iterator's last change
}

// Iterator is entrySet().iterator().
func (h *JavaHashMap[K, V]) Iterator() *JavaHashMapIterator[K, V] {
	it := &JavaHashMapIterator[K, V]{}
	it.start(h)
	return it
}

func (it *JavaHashMapIterator[K, V]) start(h *JavaHashMap[K, V]) {
	it.m = h
	if h == nil {
		return
	}
	it.expectedMods = h.modCount
	if h.size > 0 {
		it.advance()
	}
}

// advance sets pending to the head of the first non-empty bucket at or
// after nextBucket (leaving it nil if there is none).
func (it *JavaHashMapIterator[K, V]) advance() {
	buckets := it.m.table
	for it.nextBucket < len(buckets) {
		head := buckets[it.nextBucket]
		it.nextBucket++
		if head != nil {
			it.pending = head
			return
		}
	}
}

// HasNext is Iterator.hasNext(). It only reports whether an entry was found
// when the previous one was returned, so a change to the map after the last
// entry was returned does not make it fail (observed).
func (it *JavaHashMapIterator[K, V]) HasNext() bool { return it.pending != nil }

// step returns the pending entry and looks up the one after it. Iterators
// are fail-fast (Javadoc of java.util.HashMap). Observed (black-box): when
// the map was structurally modified and the iteration is also exhausted,
// ConcurrentModificationException wins over NoSuchElementException.
func (it *JavaHashMapIterator[K, V]) step() *hashMapNode[K, V] {
	if it.m.modCount != it.expectedMods {
		panic(NewConcurrentModificationException())
	}
	out := it.pending
	if out == nil {
		panic(NewNoSuchElementException())
	}
	it.lastReturned = out
	if it.pending = out.next; it.pending == nil {
		it.advance()
	}
	return out
}

// Next is Iterator.next(): the next entry's key and value.
func (it *JavaHashMapIterator[K, V]) Next() (K, V) {
	out := it.step()
	return out.key, out.value
}

// SetValue is Map.Entry.setValue on the entry Next returned last.
func (it *JavaHashMapIterator[K, V]) SetValue(value V) V {
	if it.lastReturned == nil {
		panic(NewIllegalStateException())
	}
	prior := it.lastReturned.value
	it.lastReturned.value = value
	return prior
}

// Remove is Iterator.remove(): removes the entry Next returned last.
// Observed (black-box): with no such entry (no Next yet, or Remove twice)
// it throws IllegalStateException even if the map was modified meanwhile;
// only then is a modification reported as ConcurrentModificationException.
func (it *JavaHashMapIterator[K, V]) Remove() {
	victim := it.lastReturned
	if victim == nil {
		panic(NewIllegalStateException())
	}
	if it.m.modCount != it.expectedMods {
		panic(NewConcurrentModificationException())
	}
	it.lastReturned = nil
	var zero V
	it.m.deleteKey(victim.hash, victim.key, zero, false, true)
	it.expectedMods = it.m.modCount
}

// All iterates the entries in Java's order (for (Map.Entry e : entrySet())).
func (h *JavaHashMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		var it JavaHashMapIterator[K, V]
		it.start(h)
		for it.HasNext() {
			e := it.step()
			if !yield(e.key, e.value) {
				return
			}
		}
	}
}

// Keys iterates the keys (keySet()).
func (h *JavaHashMap[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		var it JavaHashMapIterator[K, V]
		it.start(h)
		for it.HasNext() {
			if !yield(it.step().key) {
				return
			}
		}
	}
}

// Values iterates the values (values()).
func (h *JavaHashMap[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		var it JavaHashMapIterator[K, V]
		it.start(h)
		for it.HasNext() {
			if !yield(it.step().value) {
				return
			}
		}
	}
}

// KeySlice returns the keys in iteration order.
func (h *JavaHashMap[K, V]) KeySlice() []K {
	n := h.Size()
	if n == 0 {
		return nil
	}
	out := make([]K, 0, n)
	for k := range h.Keys() {
		out = append(out, k)
	}
	return out
}

// RemoveIf is entrySet().removeIf(e -> pred(key, value)): removes through
// the iterator; reports whether anything was removed.
func (h *JavaHashMap[K, V]) RemoveIf(pred func(K, V) bool) bool {
	if pred == nil {
		panic(NewNullPointerException())
	}
	removed := false
	it := h.Iterator()
	for it.HasNext() {
		if k, v := it.Next(); pred(k, v) {
			it.Remove()
			removed = true
		}
	}
	return removed
}

// JavaHashMap_KeySet is the keySet() view.
type JavaHashMap_KeySet[K any, V any] struct{ m *JavaHashMap[K, V] }

// KeySet is HashMap.keySet().
func (h *JavaHashMap[K, V]) KeySet() JavaHashMap_KeySet[K, V] { return JavaHashMap_KeySet[K, V]{h} }

func (s JavaHashMap_KeySet[K, V]) Size() int           { return s.m.Size() }
func (s JavaHashMap_KeySet[K, V]) IsEmpty() bool       { return s.m.IsEmpty() }
func (s JavaHashMap_KeySet[K, V]) Contains(key K) bool { return s.m.ContainsKey(key) }

// Remove is keySet().remove(Object): Map.remove, reporting presence.
func (s JavaHashMap_KeySet[K, V]) Remove(key K) bool {
	var zero V
	return s.m.deleteKey(hashMapHashOf(key), key, zero, false, false) != nil
}

func (s JavaHashMap_KeySet[K, V]) Clear()           { s.m.Clear() }
func (s JavaHashMap_KeySet[K, V]) All() iter.Seq[K] { return s.m.Keys() }
func (s JavaHashMap_KeySet[K, V]) ToSlice() []K     { return s.m.KeySlice() }

// String is AbstractCollection.toString(): "[k1, k2]".
func (s JavaHashMap_KeySet[K, V]) String() string {
	return javaCollectionString(any(s), s.All())
}

func (s JavaHashMap_KeySet[K, V]) containsAny(e any) bool {
	k, ok := e.(K)
	return ok && s.m.ContainsKey(k)
}

func (JavaHashMap_KeySet[K, V]) javaCollection() {}

// Equals is Set.equals(Object).
func (s JavaHashMap_KeySet[K, V]) Equals(o any) bool {
	if other, ok := o.(JavaHashMap_KeySet[K, V]); ok && other.m == s.m {
		return true
	}
	return javaSetEquals(s.Size(), s.All(), o)
}

// HashCode is Set.hashCode(): the sum of the keys' hash codes.
func (s JavaHashMap_KeySet[K, V]) HashCode() int32 {
	var sum int32
	for k := range s.All() {
		sum += ObjectHashCode(any(k))
	}
	return sum
}

func (*JavaHashMap[K, V]) javaMap() {}

func (h *JavaHashMap[K, V]) lookupAny(k any) (any, bool) {
	key, ok := k.(K)
	if !ok {
		// Only a null k can match a key that is not a K: the zero K, when K
		// is nullable.
		var zero K
		if !IsNil(k) || !IsNil(any(zero)) {
			return nil, false
		}
		key = zero
	}
	v, present := h.Lookup(key)
	return v, present
}

// Equals is Map.equals(Object): same size and equal values for every key.
func (h *JavaHashMap[K, V]) Equals(o any) bool {
	if any(h) == o {
		return true
	}
	other, ok := o.(javaMapEntries)
	if !ok || IsNil(o) || other.Size() != h.Size() {
		return false
	}
	for k, v := range h.All() {
		ov, present := other.lookupAny(any(k))
		if !present || !ObjectsEquals(any(v), ov) {
			return false
		}
	}
	return true
}

// HashCode is Map.hashCode(): the sum of hashCode(key) ^ hashCode(value).
func (h *JavaHashMap[K, V]) HashCode() int32 {
	var sum int32
	for k, v := range h.All() {
		sum += ObjectHashCode(any(k)) ^ ObjectHashCode(any(v))
	}
	return sum
}

// String is AbstractMap.toString(): "{k1=v1, k2=v2}", with "(this Map)" for
// the map itself.
func (h *JavaHashMap[K, V]) String() string {
	var sb strings.Builder
	sb.WriteByte('{')
	first := true
	for k, v := range h.All() {
		if !first {
			sb.WriteString(", ")
		}
		first = false
		if any(k) == any(h) {
			sb.WriteString("(this Map)")
		} else {
			sb.WriteString(ToJavaString(any(k)))
		}
		sb.WriteByte('=')
		if any(v) == any(h) {
			sb.WriteString("(this Map)")
		} else {
			sb.WriteString(ToJavaString(any(v)))
		}
	}
	sb.WriteByte('}')
	return sb.String()
}

// javaCollectionString is AbstractCollection.toString(): "[a, b]", with
// "(this Collection)" for the collection itself.
func javaCollectionString[T any](self any, elems iter.Seq[T]) string {
	var sb strings.Builder
	sb.WriteByte('[')
	first := true
	for e := range elems {
		if !first {
			sb.WriteString(", ")
		}
		first = false
		if any(e) == self {
			sb.WriteString("(this Collection)")
		} else {
			sb.WriteString(ToJavaString(any(e)))
		}
	}
	sb.WriteByte(']')
	return sb.String()
}

// javaSetEquals is Set.equals(o) for a set of the given size and elements.
func javaSetEquals[T any](size int, elems iter.Seq[T], o any) bool {
	other, ok := o.(javaSet)
	if !ok || IsNil(o) || other.Size() != size {
		return false
	}
	for e := range elems {
		if !other.containsAny(any(e)) {
			return false
		}
	}
	return true
}

var (
	_ JavaMap                          = (*JavaHashMap[string, int])(nil)
	_ MapSource[string, int]           = (*JavaHashMap[string, int])(nil)
	_ Collection[string]               = JavaHashMap_KeySet[string, int]{}
	_ Set[string]                      = JavaHashMap_KeySet[string, int]{}
	_ javaMapEntries                   = (*JavaHashMap[string, int])(nil)
	_ JavaCollection                   = JavaHashMap_KeySet[string, int]{}
	_ IConcurrentModificationException = (*ConcurrentModificationException)(nil)
)

// javaCompareInts is Integer.compare / Long.compare and friends.
func javaCompareInts[T int | int32 | int64 | int16 | int8 | uint16](a, b T) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}

// ---------------------------------------------------------------------------
// Sorted buckets

// root returns the root of the tree that sorted-bucket node x belongs to.
func (x *hashMapNode[K, V]) root() *hashMapNode[K, V] {
	for x.rb.parent != nil {
		x = x.rb.parent
	}
	return x
}

// searchTree looks for key k (spread hash h, Comparable class kc as given
// by javaComparableClassOf) in the subtree rooted at x and returns its
// node, or nil. It moves down one side for as long as the hashes differ or
// compareTo between two keys of one Comparable class tells the side. When
// neither does, the tree order fell back on fallbackOrder, which a lookup
// cannot reproduce (a fresh probe key has its own identity), so both
// children are searched. Keys are unique, so which child is searched first
// cannot be observed.
func (x *hashMapNode[K, V]) searchTree(h int32, k K, kc reflect.Type) *hashMapNode[K, V] {
	for x != nil {
		if h != x.hash {
			if h < x.hash {
				x = x.rb.left
			} else {
				x = x.rb.right
			}
			continue
		}
		if hashMapKeysEqual(k, x.key) {
			return x
		}
		if c, ok := javaCompareSameClass(kc, any(k), any(x.key)); ok && c != 0 {
			if c < 0 {
				x = x.rb.left
			} else {
				x = x.rb.right
			}
			continue
		}
		if found := x.rb.left.searchTree(h, k, kc); found != nil {
			return found
		}
		x = x.rb.right
	}
	return nil
}

// javaHashMapIsRed reports whether x is a red tree node (nil is black).
func javaHashMapIsRed[K any, V any](x *hashMapNode[K, V]) bool {
	return x != nil && x.rb.red
}

// treeOrder decides on which side of tree node p a new key (hash, key, kc)
// goes: negative for the left. See the observed tree order in the type
// comment: hash, then compareTo within one Comparable class, then class
// name, then identity hash code.
func treeOrder[K any, V any](hash int32, key K, kc reflect.Type, p *hashMapNode[K, V]) int {
	if hash != p.hash {
		if hash < p.hash {
			return -1
		}
		return 1
	}
	if c, ok := javaCompareSameClass(kc, any(key), any(p.key)); ok && c != 0 {
		return c
	}
	return fallbackOrder(any(key), any(p.key))
}

// fallbackOrder orders two keys with equal hashes that compareTo cannot
// order: by class name, then by identity hash code (never 0).
func fallbackOrder(a, b any) int {
	if !IsNil(a) && !IsNil(b) {
		if d := StringCompareTo(javaClassName(a), javaClassName(b)); d != 0 {
			return d
		}
	}
	if IdentityHashCode(a) <= IdentityHashCode(b) {
		return -1
	}
	return 1
}

// javaClassName is getClass().getName() for the Go types that stand for
// Java's boxed types and String.
// Go-only: language-adapter — other types have no Java class name at run
// time; their Go type name stands in (it only orders keys of different
// classes with equal hash codes, which Closure never relies on).
func javaClassName(v any) string {
	switch v.(type) {
	case string, *string:
		return "java.lang.String"
	case int32, int:
		return "java.lang.Integer"
	case int64:
		return "java.lang.Long"
	case int16:
		return "java.lang.Short"
	case int8:
		return "java.lang.Byte"
	case uint16:
		return "java.lang.Character"
	case bool:
		return "java.lang.Boolean"
	case float64:
		return "java.lang.Double"
	case float32:
		return "java.lang.Float"
	}
	return reflect.TypeOf(v).String()
}

// javaComparableClassOf is the type a sorted bucket may use compareTo for:
// v's own type when values of that type can be compared with each other
// (the Go stand-ins for String and the boxed primitives, and any type with
// a method CompareTo(itself) int), and nil for every other v, null included.
func javaComparableClassOf(v any) reflect.Type {
	switch v.(type) {
	case nil:
		return nil
	case string, int32, int, int64, int16, int8, uint16, bool, float64, float32:
		return reflect.TypeOf(v)
	case *string:
		if IsNil(v) {
			return nil
		}
		return reflect.TypeOf(v)
	}
	t := reflect.TypeOf(v)
	m, ok := t.MethodByName("CompareTo")
	if !ok || m.Type.NumIn() != 2 || m.Type.In(1) != t || m.Type.NumOut() != 1 || m.Type.Out(0).Kind() != reflect.Int {
		return nil
	}
	if IsNil(v) {
		return nil
	}
	return t
}

// javaCompareSameClass compares a with b by compareTo, where kc is a's type
// as javaComparableClassOf gives it. It reports ok = false, without
// comparing, when kc is nil or b is null or of another type.
func javaCompareSameClass(kc reflect.Type, a, b any) (int, bool) {
	if kc == nil || IsNil(b) || reflect.TypeOf(b) != kc {
		return 0, false
	}
	switch x := a.(type) {
	case string:
		return StringCompareTo(x, b.(string)), true
	case *string:
		return StringCompareTo(*x, *b.(*string)), true
	case int32:
		return javaCompareInts(x, b.(int32)), true
	case int:
		return javaCompareInts(x, b.(int)), true
	case int64:
		return javaCompareInts(x, b.(int64)), true
	case int16:
		return javaCompareInts(x, b.(int16)), true
	case int8:
		return javaCompareInts(x, b.(int8)), true
	case uint16:
		return javaCompareInts(x, b.(uint16)), true
	case bool:
		y := b.(bool)
		switch {
		case x == y:
			return 0, true
		case x:
			return 1, true
		}
		return -1, true
	case float64:
		return DoubleCompare(x, b.(float64)), true
	case float32:
		return DoubleCompare(float64(x), float64(b.(float32))), true
	}
	out := reflect.ValueOf(a).MethodByName("CompareTo").Call([]reflect.Value{reflect.ValueOf(b)})
	return int(out[0].Int()), true
}

// growOrConvertBucket makes list bucket i a sorted bucket, or, while the
// table has fewer than 64 buckets, doubles the table instead (observed).
func (h *JavaHashMap[K, V]) growOrConvertBucket(i int) {
	if len(h.table) < javaHashMapMinTreeTable {
		h.grow()
		return
	}
	h.buildTree(h.table, i)
}

// buildTree makes the chain of bucket i of tab a sorted bucket: the keys
// enter the tree in chain order, the chain keeps its order and the root
// moves to its head.
func (h *JavaHashMap[K, V]) buildTree(tab []*hashMapNode[K, V], i int) {
	var root, prev *hashMapNode[K, V]
	for x := tab[i]; x != nil; x = x.next {
		x.rb = &hashTreeLinks[K, V]{prev: prev, red: true}
		prev = x
		if root == nil {
			x.rb.red = false
			root = x
			continue
		}
		kc := javaComparableClassOf(any(x.key))
		p := root
		for {
			var next **hashMapNode[K, V]
			if treeOrder(x.hash, x.key, kc, p) <= 0 {
				next = &p.rb.left
			} else {
				next = &p.rb.right
			}
			if *next == nil {
				*next = x
				x.rb.parent = p
				break
			}
			p = *next
		}
		root = rbInsertFixup(root, x)
	}
	promoteRootToHead(tab, i, root)
}

// treePut adds key to sorted bucket i, or returns the node that already
// holds it. The new node goes right after its tree parent in the chain.
func (h *JavaHashMap[K, V]) treePut(i int, hash int32, key K, value V) *hashMapNode[K, V] {
	root := h.table[i].root()
	kc := javaComparableClassOf(any(key))
	if e := root.searchTree(hash, key, kc); e != nil {
		return e
	}
	p := root
	x := &hashMapNode[K, V]{hash: hash, key: key, value: value}
	for {
		var next **hashMapNode[K, V]
		if treeOrder(hash, key, kc, p) <= 0 {
			next = &p.rb.left
		} else {
			next = &p.rb.right
		}
		if *next == nil {
			*next = x
			break
		}
		p = *next
	}
	x.rb = &hashTreeLinks[K, V]{parent: p, prev: p, red: true}
	x.next = p.next
	if p.next != nil {
		p.next.rb.prev = x
	}
	p.next = x
	promoteRootToHead(h.table, i, rbInsertFixup(root, x))
	return nil
}

// treeRemove unlinks node p of sorted bucket i (see the observed removal
// rules in the header).
func (h *JavaHashMap[K, V]) treeRemove(i int, p *hashMapNode[K, V], byIterator bool) {
	tab := h.table
	pred, succ := p.rb.prev, p.next
	if pred == nil {
		tab[i] = succ
	} else {
		pred.next = succ
	}
	if succ != nil {
		succ.rb.prev = pred
	}
	if tab[i] == nil {
		return
	}
	root := tab[i].root()
	if !byIterator && (root.rb.right == nil || root.rb.left == nil || root.rb.left.rb.left == nil) {
		for x := tab[i]; x != nil; x = x.next {
			x.rb = nil
		}
		return
	}
	root = rbDelete(root, p)
	p.rb = nil
	if !byIterator {
		promoteRootToHead(tab, i, root)
	}
}

// splitSortedBucket distributes sorted bucket j (chain head) of the old
// table (of prevBins buckets) over buckets j and j+prevBins of tab.
func (h *JavaHashMap[K, V]) splitSortedBucket(tab []*hashMapNode[K, V], head *hashMapNode[K, V], j, prevBins int) {
	var heads, tails [2]*hashMapNode[K, V]
	var counts [2]int
	for e := head; e != nil; {
		next := e.next
		e.next = nil
		half := 0
		if int(e.hash)&prevBins != 0 {
			half = 1
		}
		e.rb.prev = tails[half]
		if tails[half] == nil {
			heads[half] = e
		} else {
			tails[half].next = e
		}
		tails[half] = e
		counts[half]++
		e = next
	}
	for half, idx := range [2]int{j, j + prevBins} {
		head := heads[half]
		tab[idx] = head
		switch {
		case head == nil:
		case counts[half] <= javaHashMapListBelow:
			for x := head; x != nil; x = x.next {
				x.rb = nil
			}
		case heads[1-half] != nil:
			h.buildTree(tab, idx)
		}
	}
}

// promoteRootToHead moves root to the head of bucket i's chain, unlinking
// it from where it was.
func promoteRootToHead[K any, V any](tab []*hashMapNode[K, V], i int, root *hashMapNode[K, V]) {
	first := tab[i]
	if root == first {
		return
	}
	before, after := root.rb.prev, root.next
	before.next = after
	if after != nil {
		after.rb.prev = before
	}
	root.rb.prev = nil
	root.next = first
	first.rb.prev = root
	tab[i] = root
}

// Textbook red-black tree operations (CLRS, "Introduction to Algorithms",
// chapter 13) on sorted-bucket nodes. Each returns the tree's root
// afterwards.

func rbRotateLeft[K any, V any](root, x *hashMapNode[K, V]) *hashMapNode[K, V] {
	y := x.rb.right
	x.rb.right = y.rb.left
	if y.rb.left != nil {
		y.rb.left.rb.parent = x
	}
	y.rb.parent = x.rb.parent
	switch xp := x.rb.parent; {
	case xp == nil:
		root = y
	case xp.rb.left == x:
		xp.rb.left = y
	default:
		xp.rb.right = y
	}
	y.rb.left = x
	x.rb.parent = y
	return root
}

func rbRotateRight[K any, V any](root, x *hashMapNode[K, V]) *hashMapNode[K, V] {
	y := x.rb.left
	x.rb.left = y.rb.right
	if y.rb.right != nil {
		y.rb.right.rb.parent = x
	}
	y.rb.parent = x.rb.parent
	switch xp := x.rb.parent; {
	case xp == nil:
		root = y
	case xp.rb.right == x:
		xp.rb.right = y
	default:
		xp.rb.left = y
	}
	y.rb.right = x
	x.rb.parent = y
	return root
}

// rbInsertFixup restores the red-black properties after red leaf z was
// attached.
func rbInsertFixup[K any, V any](root, z *hashMapNode[K, V]) *hashMapNode[K, V] {
	for javaHashMapIsRed(z.rb.parent) {
		zp := z.rb.parent
		zpp := zp.rb.parent
		if zp == zpp.rb.left {
			if y := zpp.rb.right; javaHashMapIsRed(y) {
				zp.rb.red, y.rb.red, zpp.rb.red = false, false, true
				z = zpp
				continue
			}
			if z == zp.rb.right {
				z = zp
				root = rbRotateLeft(root, z)
				zp = z.rb.parent
			}
			zp.rb.red, zpp.rb.red = false, true
			root = rbRotateRight(root, zpp)
		} else {
			if y := zpp.rb.left; javaHashMapIsRed(y) {
				zp.rb.red, y.rb.red, zpp.rb.red = false, false, true
				z = zpp
				continue
			}
			if z == zp.rb.left {
				z = zp
				root = rbRotateRight(root, z)
				zp = z.rb.parent
			}
			zp.rb.red, zpp.rb.red = false, true
			root = rbRotateLeft(root, zpp)
		}
	}
	root.rb.red = false
	return root
}

// rbTransplant puts subtree v (possibly nil) where subtree u was.
func rbTransplant[K any, V any](root, u, v *hashMapNode[K, V]) *hashMapNode[K, V] {
	switch up := u.rb.parent; {
	case up == nil:
		root = v
	case up.rb.left == u:
		up.rb.left = v
	default:
		up.rb.right = v
	}
	if v != nil {
		v.rb.parent = u.rb.parent
	}
	return root
}

// rbDelete removes z from the tree. A node with two children is replaced by
// its in-order successor node (the node itself moves, so the chain order of
// the other keys is untouched).
func rbDelete[K any, V any](root, z *hashMapNode[K, V]) *hashMapNode[K, V] {
	removedRed := z.rb.red
	var x, xp *hashMapNode[K, V]
	switch {
	case z.rb.left == nil:
		x, xp = z.rb.right, z.rb.parent
		root = rbTransplant(root, z, x)
	case z.rb.right == nil:
		x, xp = z.rb.left, z.rb.parent
		root = rbTransplant(root, z, x)
	default:
		y := z.rb.right
		for y.rb.left != nil {
			y = y.rb.left
		}
		removedRed = y.rb.red
		x = y.rb.right
		if y.rb.parent == z {
			xp = y
		} else {
			xp = y.rb.parent
			root = rbTransplant(root, y, x)
			y.rb.right = z.rb.right
			y.rb.right.rb.parent = y
		}
		root = rbTransplant(root, z, y)
		y.rb.left = z.rb.left
		y.rb.left.rb.parent = y
		y.rb.red = z.rb.red
	}
	if !removedRed {
		root = rbDeleteFixup(root, x, xp)
	}
	return root
}

// rbDeleteFixup restores the red-black properties after a black node was
// removed; x (possibly nil, with parent xp) carries the extra black.
func rbDeleteFixup[K any, V any](root, x, xp *hashMapNode[K, V]) *hashMapNode[K, V] {
	for x != root && !javaHashMapIsRed(x) {
		if x == xp.rb.left {
			w := xp.rb.right
			if javaHashMapIsRed(w) {
				w.rb.red, xp.rb.red = false, true
				root = rbRotateLeft(root, xp)
				w = xp.rb.right
			}
			if !javaHashMapIsRed(w.rb.left) && !javaHashMapIsRed(w.rb.right) {
				w.rb.red = true
				x, xp = xp, xp.rb.parent
				continue
			}
			if !javaHashMapIsRed(w.rb.right) {
				w.rb.left.rb.red, w.rb.red = false, true
				root = rbRotateRight(root, w)
				w = xp.rb.right
			}
			w.rb.red, xp.rb.red = xp.rb.red, false
			w.rb.right.rb.red = false
			root = rbRotateLeft(root, xp)
			x = root
		} else {
			w := xp.rb.left
			if javaHashMapIsRed(w) {
				w.rb.red, xp.rb.red = false, true
				root = rbRotateRight(root, xp)
				w = xp.rb.left
			}
			if !javaHashMapIsRed(w.rb.left) && !javaHashMapIsRed(w.rb.right) {
				w.rb.red = true
				x, xp = xp, xp.rb.parent
				continue
			}
			if !javaHashMapIsRed(w.rb.left) {
				w.rb.right.rb.red, w.rb.red = false, true
				root = rbRotateLeft(root, w)
				w = xp.rb.left
			}
			w.rb.red, xp.rb.red = xp.rb.red, false
			w.rb.left.rb.red = false
			root = rbRotateRight(root, xp)
			x = root
		}
	}
	if x != nil {
		x.rb.red = false
	}
	return root
}
