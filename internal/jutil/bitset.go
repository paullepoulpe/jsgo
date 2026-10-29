package jutil

import (
	"math/bits"
	"strconv"
	"strings"
)

// Go-only: language-adapter — java.util.BitSet.

// BitSet is java.util.BitSet: a growable vector of bits indexed by
// non-negative ints. The zero value is an empty set.
type BitSet struct {
	words []uint64
	buf   [2]uint64
}

// NewBitSet is new BitSet().
func NewBitSet() *BitSet {
	b := &BitSet{}
	b.words = b.buf[:0:len(b.buf)]
	return b
}

// Size is BitSet.size(): the number of bits of space in use, a multiple of
// 64.
func (b *BitSet) Size() int { return cap(b.words) * 64 }

func bitSetCheckIndex(bitIndex int, name string) {
	if bitIndex < 0 {
		panic(NewIndexOutOfBoundsException(name + " < 0: " + strconv.Itoa(bitIndex)))
	}
}

func (b *BitSet) ensure(wordIndex int) {
	if wordIndex < len(b.words) {
		return
	}
	if b.words == nil && wordIndex < len(b.buf) {
		clear(b.buf[:wordIndex+1])
		b.words = b.buf[:wordIndex+1]
		return
	}
	if wordIndex < cap(b.words) {
		prev := len(b.words)
		b.words = b.words[:wordIndex+1]
		clear(b.words[prev:])
		return
	}
	newCap := max(wordIndex+1, 2*cap(b.words))
	newWords := make([]uint64, wordIndex+1, newCap)
	copy(newWords, b.words)
	b.words = newWords
}

// trim drops trailing zero words so that len(words) is Java's wordsInUse.
func (b *BitSet) trim() {
	n := len(b.words)
	for n > 0 && b.words[n-1] == 0 {
		n--
	}
	b.words = b.words[:n]
}

// Set is BitSet.set(int).
func (b *BitSet) Set(bitIndex int) {
	bitSetCheckIndex(bitIndex, "bitIndex")
	w := bitIndex >> 6
	b.ensure(w)
	b.words[w] |= 1 << (uint(bitIndex) & 63)
}

// SetRange is BitSet.set(int fromIndex, int toIndex).
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

func bitSetCheckRange(fromIndex, toIndex int) {
	bitSetCheckIndex(fromIndex, "fromIndex")
	bitSetCheckIndex(toIndex, "toIndex")
	if fromIndex > toIndex {
		panic(NewIndexOutOfBoundsException("fromIndex: " + strconv.Itoa(fromIndex) +
			" > toIndex: " + strconv.Itoa(toIndex)))
	}
}

// Clear is BitSet.clear(int).
func (b *BitSet) Clear(bitIndex int) {
	bitSetCheckIndex(bitIndex, "bitIndex")
	w := bitIndex >> 6
	if w >= len(b.words) {
		return
	}
	b.words[w] &^= 1 << (uint(bitIndex) & 63)
	b.trim()
}

// ClearAll is BitSet.clear().
func (b *BitSet) ClearAll() { b.words = b.words[:0] }

// FlipRange is BitSet.flip(int fromIndex, int toIndex): sets each bit from
// fromIndex (inclusive) to toIndex (exclusive) to the complement of its
// current value.
func (b *BitSet) FlipRange(fromIndex, toIndex int) {
	bitSetCheckRange(fromIndex, toIndex)
	if fromIndex == toIndex {
		return
	}
	b.ensure((toIndex - 1) >> 6)
	for i := fromIndex; i < toIndex; i++ {
		b.words[i>>6] ^= 1 << (uint(i) & 63)
	}
	b.trim()
}

// NegativeArraySizeException is java.lang.NegativeArraySizeException.
type NegativeArraySizeException struct{ RuntimeException }

func (*NegativeArraySizeException) negativeArraySizeException() {}

// INegativeArraySizeException matches NegativeArraySizeException.
type INegativeArraySizeException interface {
	IRuntimeException
	negativeArraySizeException()
}

// NewNegativeArraySizeException constructs the exception with a message.
func NewNegativeArraySizeException(message string) *NegativeArraySizeException {
	e := &NegativeArraySizeException{}
	InitThrowable(e, "java.lang.NegativeArraySizeException", msgPtr(message), nil)
	return e
}

// Get is BitSet.get(int).
func (b *BitSet) Get(bitIndex int) bool {
	bitSetCheckIndex(bitIndex, "bitIndex")
	w := bitIndex >> 6
	return w < len(b.words) && b.words[w]&(1<<(uint(bitIndex)&63)) != 0
}

// Length is BitSet.length(): the index of the highest set bit plus one.
func (b *BitSet) Length() int {
	if len(b.words) == 0 {
		return 0
	}
	last := b.words[len(b.words)-1]
	return 64*(len(b.words)-1) + (64 - bits.LeadingZeros64(last))
}

// IsEmpty is BitSet.isEmpty().
func (b *BitSet) IsEmpty() bool { return len(b.words) == 0 }

// Cardinality is BitSet.cardinality().
func (b *BitSet) Cardinality() int {
	n := 0
	for _, w := range b.words {
		n += bits.OnesCount64(w)
	}
	return n
}

// NextSetBit is BitSet.nextSetBit(int): -1 if there is none.
func (b *BitSet) NextSetBit(fromIndex int) int {
	bitSetCheckIndex(fromIndex, "fromIndex")
	w := fromIndex >> 6
	if w >= len(b.words) {
		return -1
	}
	word := b.words[w] & (^uint64(0) << (uint(fromIndex) & 63))
	for {
		if word != 0 {
			return w*64 + bits.TrailingZeros64(word)
		}
		w++
		if w == len(b.words) {
			return -1
		}
		word = b.words[w]
	}
}

// PreviousSetBit is BitSet.previousSetBit(int): the index of the nearest
// set bit at or before fromIndex, or -1 if there is none (or fromIndex is
// -1). It panics with IndexOutOfBoundsException if fromIndex < -1.
func (b *BitSet) PreviousSetBit(fromIndex int) int {
	if fromIndex < 0 {
		if fromIndex == -1 {
			return -1
		}
		panic(NewIndexOutOfBoundsException("fromIndex < -1: " + strconv.Itoa(fromIndex)))
	}
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

// Or is BitSet.or(BitSet).
func (b *BitSet) Or(set *BitSet) {
	if b == set || len(set.words) == 0 {
		return
	}
	b.ensure(len(set.words) - 1)
	for i, w := range set.words {
		b.words[i] |= w
	}
}

// Words128 returns the low 128 bits of b as two uint64 words.
func (b *BitSet) Words128() (w0, w1 uint64) {
	if len(b.words) > 0 {
		w0 = b.words[0]
		if len(b.words) > 1 {
			w1 = b.words[1]
		}
	}
	return w0, w1
}

// And is BitSet.and(BitSet).
func (b *BitSet) And(set *BitSet) {
	if b == set {
		return
	}
	for i := range b.words {
		if i < len(set.words) {
			b.words[i] &= set.words[i]
		} else {
			b.words[i] = 0
		}
	}
	b.trim()
}

// AndNot is BitSet.andNot(BitSet).
func (b *BitSet) AndNot(set *BitSet) {
	for i := 0; i < len(b.words) && i < len(set.words); i++ {
		b.words[i] &^= set.words[i]
	}
	b.trim()
}

// Intersects is BitSet.intersects(BitSet).
func (b *BitSet) Intersects(set *BitSet) bool {
	for i := 0; i < len(b.words) && i < len(set.words); i++ {
		if b.words[i]&set.words[i] != 0 {
			return true
		}
	}
	return false
}

// Clone is BitSet.clone().
func (b *BitSet) Clone() *BitSet {
	dst := &BitSet{}
	n := len(b.words)
	c := cap(b.words)
	if c <= len(dst.buf) {
		dst.words = dst.buf[:n:c]
	} else {
		dst.words = make([]uint64, n, c)
	}
	copy(dst.words, b.words)
	return dst
}

// CopyFrom copies src into b, reusing b's existing slice capacity when possible.
func (b *BitSet) CopyFrom(src *BitSet) {
	if b == src {
		return
	}
	n := len(src.words)
	if n <= cap(b.words) {
		b.words = b.words[:n]
	} else if b.words == nil && n <= len(b.buf) {
		b.words = b.buf[:n]
	} else {
		b.words = make([]uint64, n)
	}
	copy(b.words, src.words)
}

// Equals is BitSet.equals(Object).
func (b *BitSet) Equals(o any) bool {
	other, ok := o.(*BitSet)
	if !ok || other == nil {
		return false
	}
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

// HashCode is BitSet.hashCode().
func (b *BitSet) HashCode() int32 {
	h := int64(1234)
	for i := len(b.words) - 1; i >= 0; i-- {
		h ^= int64(b.words[i]) * int64(i+1)
	}
	return int32((h >> 32) ^ h)
}

// String is BitSet.toString(): "{1, 3, 64}".
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
