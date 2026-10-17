package jutil

// Go-only: language-adapter — java.util.BitSet.flip(int) and
// BitSet.flip(int, int).

// Flip is BitSet.flip(int): sets the bit at bitIndex to the complement of
// its current value.
func (b *BitSet) Flip(bitIndex int) {
	bitSetCheckIndex(bitIndex, "bitIndex")
	w := bitIndex >> 6
	b.ensure(w)
	b.words[w] ^= 1 << (uint(bitIndex) & 63)
	b.trim()
}

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
