package jutil

import "math/bits"

// Go-only: language-adapter — java.util.BitSet.

// BitSet is java.util.BitSet: a growable vector of bits. Size() reports the
// capacity in bits (words * 64) and Length() the index of the highest set bit
// plus one, with the JDK's growth policy, so code that inspects either sees
// the same values as in Java.
type BitSet struct {
	words []uint64
	// wordsInUse is the number of words in the logical size of this set
	// (the JDK's invariant: words[wordsInUse-1] != 0 when wordsInUse > 0).
	wordsInUse int
}

func bitSetWordIndex(bitIndex int) int { return bitIndex >> 6 }

// NewBitSet is new BitSet(): 64 bits of initial capacity.
func NewBitSet() *BitSet {
	return &BitSet{words: make([]uint64, 1)}
}

// NewBitSetWithSize is new BitSet(nbits).
func NewBitSetWithSize(nbits int) *BitSet {
	if nbits < 0 {
		panic(NewNegativeArraySizeException("nbits < 0: " + itoa(nbits)))
	}
	return &BitSet{words: make([]uint64, bitSetWordIndex(nbits-1)+1)}
}

func (b *BitSet) recalculateWordsInUse() {
	i := b.wordsInUse - 1
	for ; i >= 0; i-- {
		if b.words[i] != 0 {
			break
		}
	}
	b.wordsInUse = i + 1
}

func (b *BitSet) ensureCapacity(wordsRequired int) {
	if len(b.words) < wordsRequired {
		request := max(2*len(b.words), wordsRequired)
		nw := make([]uint64, request)
		copy(nw, b.words)
		b.words = nw
	}
}

func (b *BitSet) expandTo(wordIndex int) {
	wordsRequired := wordIndex + 1
	if b.wordsInUse < wordsRequired {
		b.ensureCapacity(wordsRequired)
		b.wordsInUse = wordsRequired
	}
}

func checkBitIndex(bitIndex int) {
	if bitIndex < 0 {
		panic(NewIndexOutOfBoundsException("bitIndex < 0: " + itoa(bitIndex)))
	}
}

// Set is BitSet.set(int).
func (b *BitSet) Set(bitIndex int) {
	checkBitIndex(bitIndex)
	wordIndex := bitSetWordIndex(bitIndex)
	b.expandTo(wordIndex)
	b.words[wordIndex] |= uint64(1) << (uint(bitIndex) & 63)
}

// SetValue is BitSet.set(int, boolean).
func (b *BitSet) SetValue(bitIndex int, value bool) {
	if value {
		b.Set(bitIndex)
	} else {
		b.Clear(bitIndex)
	}
}

// Clear is BitSet.clear(int).
func (b *BitSet) Clear(bitIndex int) {
	checkBitIndex(bitIndex)
	wordIndex := bitSetWordIndex(bitIndex)
	if wordIndex >= b.wordsInUse {
		return
	}
	b.words[wordIndex] &^= uint64(1) << (uint(bitIndex) & 63)
	b.recalculateWordsInUse()
}

// Get is BitSet.get(int).
func (b *BitSet) Get(bitIndex int) bool {
	checkBitIndex(bitIndex)
	wordIndex := bitSetWordIndex(bitIndex)
	return wordIndex < b.wordsInUse && b.words[wordIndex]&(uint64(1)<<(uint(bitIndex)&63)) != 0
}

// Or is BitSet.or(BitSet).
func (b *BitSet) Or(set *BitSet) {
	if b == set {
		return
	}
	wordsInCommon := min(b.wordsInUse, set.wordsInUse)
	if b.wordsInUse < set.wordsInUse {
		b.ensureCapacity(set.wordsInUse)
		b.wordsInUse = set.wordsInUse
	}
	for i := 0; i < wordsInCommon; i++ {
		b.words[i] |= set.words[i]
	}
	if wordsInCommon < set.wordsInUse {
		copy(b.words[wordsInCommon:b.wordsInUse], set.words[wordsInCommon:set.wordsInUse])
	}
}

// Length is BitSet.length(): the index of the highest set bit plus one.
func (b *BitSet) Length() int {
	if b.wordsInUse == 0 {
		return 0
	}
	return 64*(b.wordsInUse-1) + (64 - bits.LeadingZeros64(b.words[b.wordsInUse-1]))
}

// Size is BitSet.size(): the number of bits of space actually in use.
func (b *BitSet) Size() int { return len(b.words) * 64 }

// NegativeArraySizeException is java.lang.NegativeArraySizeException.
type NegativeArraySizeException struct{ RuntimeException }

func (*NegativeArraySizeException) negativeArraySizeException() {}

// INegativeArraySizeException marks NegativeArraySizeException and its
// subclasses for jutil.Try.
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
