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
}

// NewBitSet is new BitSet().
func NewBitSet() *BitSet { return &BitSet{} }

func bitSetCheckIndex(bitIndex int, name string) {
	if bitIndex < 0 {
		panic(NewIndexOutOfBoundsException(name + " < 0: " + strconv.Itoa(bitIndex)))
	}
}

func (b *BitSet) ensure(wordIndex int) {
	if wordIndex >= len(b.words) {
		b.words = append(b.words, make([]uint64, wordIndex+1-len(b.words))...)
	}
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

// SetTo is BitSet.set(int, boolean).
func (b *BitSet) SetTo(bitIndex int, value bool) {
	if value {
		b.Set(bitIndex)
	} else {
		b.Clear(bitIndex)
	}
}

// SetRange is BitSet.set(int fromIndex, int toIndex).
func (b *BitSet) SetRange(fromIndex, toIndex int) {
	bitSetCheckRange(fromIndex, toIndex)
	for i := fromIndex; i < toIndex; i++ {
		b.Set(i)
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

// NextClearBit is BitSet.nextClearBit(int).
func (b *BitSet) NextClearBit(fromIndex int) int {
	bitSetCheckIndex(fromIndex, "fromIndex")
	w := fromIndex >> 6
	if w >= len(b.words) {
		return fromIndex
	}
	word := ^b.words[w] & (^uint64(0) << (uint(fromIndex) & 63))
	for {
		if word != 0 {
			return w*64 + bits.TrailingZeros64(word)
		}
		w++
		if w == len(b.words) {
			return len(b.words) * 64
		}
		word = ^b.words[w]
	}
}

// Or is BitSet.or(BitSet).
func (b *BitSet) Or(set *BitSet) {
	if b == set {
		return
	}
	b.ensure(len(set.words) - 1)
	for i, w := range set.words {
		b.words[i] |= w
	}
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

// ToByteArray is BitSet.toByteArray(): little-endian bytes, without
// trailing zero bytes.
func (b *BitSet) ToByteArray() []byte {
	n := len(b.words)
	if n == 0 {
		return []byte{}
	}
	length := 8 * (n - 1)
	for x := b.words[n-1]; x != 0; x >>= 8 {
		length++
	}
	out := make([]byte, length)
	for i := 0; i < length; i++ {
		out[i] = byte(b.words[i/8] >> (8 * uint(i%8)))
	}
	return out
}

// Clone is BitSet.clone().
func (b *BitSet) Clone() *BitSet {
	return &BitSet{words: append([]uint64(nil), b.words...)}
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
