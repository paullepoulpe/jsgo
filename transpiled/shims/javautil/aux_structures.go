// Clean-room specification based on java.util.ArrayDeque, java.util.BitSet, java.util.Optional Javadoc.

package javautil

import (
	"iter"
	"math/bits"
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// ArrayDeque is an O(1) amortized circular resizable double-ended queue.
type ArrayDeque[T any] struct {
	buf  []T
	head int
	size int
}

// NewArrayDeque creates an empty ArrayDeque.
func NewArrayDeque[T any]() *ArrayDeque[T] {
	return &ArrayDeque[T]{}
}

// NewArrayDequeWithCapacity creates an ArrayDeque pre-sized for capacity.
func NewArrayDequeWithCapacity[T any](capacity int) *ArrayDeque[T] {
	n := 8
	for n < capacity {
		n <<= 1
	}
	return &ArrayDeque[T]{buf: make([]T, n)}
}

func (q *ArrayDeque[T]) Size() int {
	if q == nil {
		return 0
	}
	return q.size
}

func (q *ArrayDeque[T]) IsEmpty() bool {
	return q.Size() == 0
}

func (q *ArrayDeque[T]) grow() {
	newCap := 8
	if len(q.buf) > 0 {
		newCap = len(q.buf) << 1
	}
	nb := make([]T, newCap)
	if q.size > 0 {
		mask := len(q.buf) - 1
		for i := 0; i < q.size; i++ {
			nb[i] = q.buf[(q.head+i)&mask]
		}
	}
	q.buf = nb
	q.head = 0
}

// AddFirst inserts an element at the front in O(1) amortized time.
func (q *ArrayDeque[T]) AddFirst(v T) {
	if q.size == len(q.buf) {
		q.grow()
	}
	mask := len(q.buf) - 1
	q.head = (q.head - 1) & mask
	q.buf[q.head] = v
	q.size++
}

// AddLast appends an element at the end in O(1) amortized time.
func (q *ArrayDeque[T]) AddLast(v T) {
	if q.size == len(q.buf) {
		q.grow()
	}
	mask := len(q.buf) - 1
	q.buf[(q.head+q.size)&mask] = v
	q.size++
}

// Add appends an element at the end and returns true.
func (q *ArrayDeque[T]) Add(v T) bool {
	q.AddLast(v)
	return true
}

// OfferFirst inserts an element at the front and returns true.
func (q *ArrayDeque[T]) OfferFirst(v T) bool {
	q.AddFirst(v)
	return true
}

// OfferLast appends an element at the end and returns true.
func (q *ArrayDeque[T]) OfferLast(v T) bool {
	q.AddLast(v)
	return true
}

// Offer appends an element at the end and returns true.
func (q *ArrayDeque[T]) Offer(v T) bool {
	q.AddLast(v)
	return true
}

// Push pushes an element onto the stack (front).
func (q *ArrayDeque[T]) Push(v T) {
	q.AddFirst(v)
}

// Pop pops an element from the stack (front), panicking if empty.
func (q *ArrayDeque[T]) Pop() T {
	v, ok := q.PollFirst()
	if !ok {
		panic(javalang.NewNoSuchElementException("empty deque on pop"))
	}
	return v
}

// RemoveFirst removes and returns the first element, panicking if empty.
func (q *ArrayDeque[T]) RemoveFirst() T {
	return q.Pop()
}

// RemoveLast removes and returns the last element, panicking if empty.
func (q *ArrayDeque[T]) RemoveLast() T {
	v, ok := q.PollLast()
	if !ok {
		panic(javalang.NewNoSuchElementException("empty deque on removeLast"))
	}
	return v
}

// Peek returns the front element without removing it, or zero value if empty.
func (q *ArrayDeque[T]) Peek() T {
	return q.PeekFirst()
}

// PeekFirst returns the front element without removing it, or zero value if empty.
func (q *ArrayDeque[T]) PeekFirst() T {
	if q.IsEmpty() {
		var zero T
		return zero
	}
	return q.buf[q.head]
}

// PeekLast returns the back element without removing it, or zero value if empty.
func (q *ArrayDeque[T]) PeekLast() T {
	if q.IsEmpty() {
		var zero T
		return zero
	}
	mask := len(q.buf) - 1
	return q.buf[(q.head+q.size-1)&mask]
}

// Poll retrieves and removes the head of the queue, or (zero, false) if empty.
func (q *ArrayDeque[T]) Poll() (T, bool) {
	return q.PollFirst()
}

// PollFirst retrieves and removes the first element, or (zero, false) if empty.
func (q *ArrayDeque[T]) PollFirst() (T, bool) {
	if q.IsEmpty() {
		var zero T
		return zero, false
	}
	var zero T
	v := q.buf[q.head]
	q.buf[q.head] = zero
	q.head = (q.head + 1) & (len(q.buf) - 1)
	q.size--
	return v, true
}

// PollLast retrieves and removes the last element, or (zero, false) if empty.
func (q *ArrayDeque[T]) PollLast() (T, bool) {
	if q.IsEmpty() {
		var zero T
		return zero, false
	}
	var zero T
	idx := (q.head + q.size - 1) & (len(q.buf) - 1)
	v := q.buf[idx]
	q.buf[idx] = zero
	q.size--
	return v, true
}

// Clear removes all elements from this deque.
func (q *ArrayDeque[T]) Clear() {
	if q == nil || q.size == 0 {
		return
	}
	clear(q.buf)
	q.head = 0
	q.size = 0
}

// Contains reports whether elem is in the deque.
func (q *ArrayDeque[T]) Contains(elem T) bool {
	if q == nil || q.size == 0 {
		return false
	}
	mask := len(q.buf) - 1
	for i := 0; i < q.size; i++ {
		if javalang.ObjectsEquals(q.buf[(q.head+i)&mask], elem) {
			return true
		}
	}
	return false
}

// ToSlice returns a slice of elements from front to back.
func (q *ArrayDeque[T]) ToSlice() []T {
	if q == nil || q.size == 0 {
		return nil
	}
	out := make([]T, q.size)
	mask := len(q.buf) - 1
	for i := 0; i < q.size; i++ {
		out[i] = q.buf[(q.head+i)&mask]
	}
	return out
}

// All returns a front-to-back iterator over the deque.
func (q *ArrayDeque[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		if q == nil || q.size == 0 {
			return
		}
		mask := len(q.buf) - 1
		for i := 0; i < q.size; i++ {
			if !yield(q.buf[(q.head+i)&mask]) {
				return
			}
		}
	}
}

// ---------------------------------------------------------------------------
// BitSet
// ---------------------------------------------------------------------------

// BitSet implements a vector of bits that grows as needed (java.util.BitSet).
type BitSet struct {
	words []uint64
	buf   [2]uint64
}

// NewBitSet creates a new empty BitSet.
func NewBitSet() *BitSet {
	b := &BitSet{}
	b.words = b.buf[:0:len(b.buf)]
	return b
}

// NewBitSetWithCapacity creates a BitSet with enough initial capacity for nbits.
func NewBitSetWithCapacity(nbits int) *BitSet {
	if nbits < 0 {
		panic(javalang.NewIndexOutOfBoundsException("nbits < 0: " + strconv.Itoa(nbits)))
	}
	if nbits == 0 {
		return NewBitSet()
	}
	words := (nbits + 63) / 64
	b := &BitSet{}
	if words <= len(b.buf) {
		b.words = b.buf[:0:len(b.buf)]
	} else {
		b.words = make([]uint64, 0, words)
	}
	return b
}

func bitSetCheckIndex(bitIndex int, name string) {
	if bitIndex < 0 {
		panic(javalang.NewIndexOutOfBoundsException(name + " < 0: " + strconv.Itoa(bitIndex)))
	}
}

func bitSetCheckRange(fromIndex, toIndex int) {
	bitSetCheckIndex(fromIndex, "fromIndex")
	bitSetCheckIndex(toIndex, "toIndex")
	if fromIndex > toIndex {
		panic(javalang.NewIndexOutOfBoundsException("fromIndex: " + strconv.Itoa(fromIndex) +
			" > toIndex: " + strconv.Itoa(toIndex)))
	}
}

func (b *BitSet) ensure(wordIdx int) {
	if wordIdx < len(b.words) {
		return
	}
	if b.words == nil && wordIdx < len(b.buf) {
		clear(b.buf[:wordIdx+1])
		b.words = b.buf[:wordIdx+1]
		return
	}
	if wordIdx < cap(b.words) {
		prev := len(b.words)
		b.words = b.words[:wordIdx+1]
		clear(b.words[prev:])
		return
	}
	newCap := max(wordIdx+1, 2*cap(b.words))
	newWords := make([]uint64, wordIdx+1, newCap)
	copy(newWords, b.words)
	b.words = newWords
}

func (b *BitSet) trim() {
	n := len(b.words)
	for n > 0 && b.words[n-1] == 0 {
		n--
	}
	b.words = b.words[:n]
}

// Size returns the number of bits of space actually in use (a multiple of 64).
func (b *BitSet) Size() int {
	return cap(b.words) * 64
}

// Length returns the logical size of this BitSet: highest set bit index + 1, or 0 if empty.
func (b *BitSet) Length() int {
	b.trim()
	if len(b.words) == 0 {
		return 0
	}
	last := b.words[len(b.words)-1]
	return 64*(len(b.words)-1) + (64 - bits.LeadingZeros64(last))
}

// IsEmpty reports whether this BitSet contains no bits set to true.
func (b *BitSet) IsEmpty() bool {
	b.trim()
	return len(b.words) == 0
}

// Set sets the bit at bitIndex to true.
func (b *BitSet) Set(bitIndex int) {
	bitSetCheckIndex(bitIndex, "bitIndex")
	w := bitIndex >> 6
	b.ensure(w)
	b.words[w] |= (1 << (uint(bitIndex) & 63))
}

// SetVal sets the bit at bitIndex to value.
func (b *BitSet) SetVal(bitIndex int, value bool) {
	if value {
		b.Set(bitIndex)
	} else {
		b.Clear(bitIndex)
	}
}

// SetRange sets bits from fromIndex (inclusive) to toIndex (exclusive) to true.
func (b *BitSet) SetRange(fromIndex, toIndex int) {
	bitSetCheckRange(fromIndex, toIndex)
	if fromIndex == toIndex {
		return
	}
	startWord := fromIndex >> 6
	endWord := (toIndex - 1) >> 6
	b.ensure(endWord)
	firstWordMask := ^uint64(0) << (uint(fromIndex) & 63)
	lastWordMask := ^uint64(0) >> (63 - (uint(toIndex-1) & 63))
	if startWord == endWord {
		b.words[startWord] |= firstWordMask & lastWordMask
	} else {
		b.words[startWord] |= firstWordMask
		for i := startWord + 1; i < endWord; i++ {
			b.words[i] = ^uint64(0)
		}
		b.words[endWord] |= lastWordMask
	}
}

// Clear clears the bit at bitIndex.
func (b *BitSet) Clear(bitIndex int) {
	bitSetCheckIndex(bitIndex, "bitIndex")
	w := bitIndex >> 6
	if w < len(b.words) {
		b.words[w] &^= (1 << (uint(bitIndex) & 63))
		b.trim()
	}
}

// ClearAll clears all bits in this BitSet.
func (b *BitSet) ClearAll() {
	b.words = b.words[:0]
}

// Flip complements the bit at bitIndex.
func (b *BitSet) Flip(bitIndex int) {
	bitSetCheckIndex(bitIndex, "bitIndex")
	w := bitIndex >> 6
	b.ensure(w)
	b.words[w] ^= (1 << (uint(bitIndex) & 63))
	b.trim()
}

// FlipRange complements bits in [fromIndex, toIndex).
func (b *BitSet) FlipRange(fromIndex, toIndex int) {
	bitSetCheckRange(fromIndex, toIndex)
	if fromIndex == toIndex {
		return
	}
	startWord := fromIndex >> 6
	endWord := (toIndex - 1) >> 6
	b.ensure(endWord)
	firstWordMask := ^uint64(0) << (uint(fromIndex) & 63)
	lastWordMask := ^uint64(0) >> (63 - (uint(toIndex-1) & 63))
	if startWord == endWord {
		b.words[startWord] ^= firstWordMask & lastWordMask
	} else {
		b.words[startWord] ^= firstWordMask
		for i := startWord + 1; i < endWord; i++ {
			b.words[i] = ^b.words[i]
		}
		b.words[endWord] ^= lastWordMask
	}
	b.trim()
}

// Get returns the value of the bit with the specified index.
func (b *BitSet) Get(bitIndex int) bool {
	bitSetCheckIndex(bitIndex, "bitIndex")
	w := bitIndex >> 6
	return w < len(b.words) && (b.words[w]&(1<<(uint(bitIndex)&63))) != 0
}

// NextSetBit returns the index of the first bit set to true on or after fromIndex, or -1.
func (b *BitSet) NextSetBit(fromIndex int) int {
	bitSetCheckIndex(fromIndex, "fromIndex")
	u := fromIndex >> 6
	if u >= len(b.words) {
		return -1
	}
	word := b.words[u] & (^uint64(0) << (uint(fromIndex) & 63))
	for {
		if word != 0 {
			return u*64 + bits.TrailingZeros64(word)
		}
		u++
		if u >= len(b.words) {
			return -1
		}
		word = b.words[u]
	}
}

// NextClearBit returns the index of the first bit set to false on or after fromIndex.
func (b *BitSet) NextClearBit(fromIndex int) int {
	bitSetCheckIndex(fromIndex, "fromIndex")
	u := fromIndex >> 6
	if u >= len(b.words) {
		return fromIndex
	}
	word := ^b.words[u] & (^uint64(0) << (uint(fromIndex) & 63))
	for {
		if word != 0 {
			return u*64 + bits.TrailingZeros64(word)
		}
		u++
		if u >= len(b.words) {
			return len(b.words) * 64
		}
		word = ^b.words[u]
	}
}

// PreviousSetBit returns the index of the nearest bit set to true on or before fromIndex, or -1.
func (b *BitSet) PreviousSetBit(fromIndex int) int {
	if fromIndex < 0 {
		if fromIndex == -1 {
			return -1
		}
		panic(javalang.NewIndexOutOfBoundsException("fromIndex < -1: " + strconv.Itoa(fromIndex)))
	}
	b.trim()
	w := fromIndex >> 6
	if w >= len(b.words) {
		return b.Length() - 1
	}
	word := b.words[w] & (^uint64(0) >> (63 - (uint(fromIndex) & 63)))
	for {
		if word != 0 {
			return (w+1)*64 - 1 - bits.LeadingZeros64(word)
		}
		if w == 0 {
			return -1
		}
		w--
		word = b.words[w]
	}
}

// Cardinality returns the number of bits set to true.
func (b *BitSet) Cardinality() int {
	cnt := 0
	for _, w := range b.words {
		cnt += bits.OnesCount64(w)
	}
	return cnt
}

// Intersects reports whether other has any bits set to true that are also set in b.
func (b *BitSet) Intersects(other *BitSet) bool {
	for i := 0; i < len(b.words) && i < len(other.words); i++ {
		if b.words[i]&other.words[i] != 0 {
			return true
		}
	}
	return false
}

// And performs a logical AND with other in-place.
func (b *BitSet) And(other *BitSet) {
	if b == other {
		return
	}
	n := min(len(b.words), len(other.words))
	for i := 0; i < n; i++ {
		b.words[i] &= other.words[i]
	}
	b.words = b.words[:n]
	b.trim()
}

// Or performs a logical OR with other in-place.
func (b *BitSet) Or(other *BitSet) {
	if b == other {
		return
	}
	other.trim()
	if len(other.words) == 0 {
		return
	}
	b.ensure(len(other.words) - 1)
	for i, w := range other.words {
		b.words[i] |= w
	}
}

// Xor performs a logical XOR with other in-place.
func (b *BitSet) Xor(other *BitSet) {
	if b == other {
		b.ClearAll()
		return
	}
	other.trim()
	if len(other.words) == 0 {
		return
	}
	b.ensure(len(other.words) - 1)
	for i, w := range other.words {
		b.words[i] ^= w
	}
	b.trim()
}

// AndNot clears all bits in b whose corresponding bit is set in other.
func (b *BitSet) AndNot(other *BitSet) {
	n := min(len(b.words), len(other.words))
	for i := 0; i < n; i++ {
		b.words[i] &^= other.words[i]
	}
	b.trim()
}

// Clone returns a deep copy of this BitSet.
func (b *BitSet) Clone() *BitSet {
	b.trim()
	dst := &BitSet{}
	n := len(b.words)
	if n <= len(dst.buf) {
		dst.words = dst.buf[:n:len(dst.buf)]
	} else {
		dst.words = make([]uint64, n)
	}
	copy(dst.words, b.words)
	return dst
}

// CopyFrom overwrites b with the bits of src.
func (b *BitSet) CopyFrom(src *BitSet) {
	if b == src {
		return
	}
	src.trim()
	n := len(src.words)
	if n <= cap(b.words) {
		b.words = b.words[:n]
	} else if b.words == nil && n <= len(b.buf) {
		b.words = b.buf[:n:len(b.buf)]
	} else {
		b.words = make([]uint64, n)
	}
	copy(b.words, src.words)
}

// Equals implements Java BitSet.equals(Object).
func (b *BitSet) Equals(o any) bool {
	other, ok := o.(*BitSet)
	if !ok || other == nil {
		return false
	}
	if b == other {
		return true
	}
	b.trim()
	other.trim()
	if len(b.words) != len(other.words) {
		return false
	}
	for i := range b.words {
		if b.words[i] != other.words[i] {
			return false
		}
	}
	return true
}

// HashCode implements Java BitSet.hashCode().
func (b *BitSet) HashCode() int32 {
	b.trim()
	h := int64(1234)
	for i := len(b.words) - 1; i >= 0; i-- {
		h ^= int64(b.words[i]) * int64(i+1)
	}
	return int32((h >> 32) ^ h)
}

// String implements Java BitSet.toString(), e.g. "{1, 3, 64}".
func (b *BitSet) String() string {
	var sb strings.Builder
	sb.WriteByte('{')
	first := true
	for i := b.NextSetBit(0); i >= 0; i = b.NextSetBit(i + 1) {
		if !first {
			sb.WriteString(", ")
		}
		first = false
		sb.WriteString(strconv.Itoa(i))
	}
	sb.WriteByte('}')
	return sb.String()
}

// ---------------------------------------------------------------------------
// Optional
// ---------------------------------------------------------------------------

// Optional is a container object which may or may not contain a non-null value.
type Optional[T any] struct {
	val     T
	present bool
}

// Optional_Of returns an Optional describing the given non-null value.
func Optional_Of[T any](v T) Optional[T] {
	if javalang.IsNil(v) {
		panic(javalang.NewNullPointerExceptionMsg("Optional.of received null value"))
	}
	return Optional[T]{val: v, present: true}
}

// Optional_OfNullable returns an Optional describing the given value if non-null, or empty.
func Optional_OfNullable[T any](v T) Optional[T] {
	if javalang.IsNil(v) {
		return Optional[T]{present: false}
	}
	return Optional[T]{val: v, present: true}
}

// Optional_Empty returns an empty Optional instance.
func Optional_Empty[T any]() Optional[T] {
	return Optional[T]{present: false}
}

// IsPresent reports whether there is a value present.
func (o Optional[T]) IsPresent() bool {
	return o.present
}

// IsEmpty reports whether there is no value present.
func (o Optional[T]) IsEmpty() bool {
	return !o.present
}

// Get returns the value if present, or panics with NoSuchElementException.
func (o Optional[T]) Get() T {
	if !o.present {
		panic(javalang.NewNoSuchElementException("No value present in Optional"))
	}
	return o.val
}

// OrElse returns the value if present, otherwise returns other.
func (o Optional[T]) OrElse(other T) T {
	if o.present {
		return o.val
	}
	return other
}

// OrNull returns the value if present, otherwise returns the zero value (Guava Optional.orNull).
func (o Optional[T]) OrNull() T {
	if o.present {
		return o.val
	}
	var zero T
	return zero
}
