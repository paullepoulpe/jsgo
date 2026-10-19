package jutil

// Go-only: language-adapter — java.util.BitSet(int nbits) and
// BitSet.size().

// NewBitSetWithSize is new BitSet(nbits): an empty set with room for nbits
// bits. new BitSet() is NewBitSetWithSize(64).
func NewBitSetWithSize(nbits int) *BitSet {
	b := &BitSet{}
	b.InitWithSize(nbits)
	return b
}

// InitWithSize initializes b for nbits bits, using b's inline buffer when it fits.
func (b *BitSet) InitWithSize(nbits int) {
	if nbits < 0 {
		panic(NewNegativeArraySizeException("nbits < 0: " + itoa(nbits)))
	}
	nwords := (nbits + 63) >> 6
	if nwords <= len(b.buf) {
		b.words = b.buf[:0:nwords]
	} else {
		b.words = make([]uint64, 0, nwords)
	}
}

// Size is BitSet.size(): the number of bits of space in use, a multiple of
// 64. It is the capacity reserved by NewBitSetWithSize or, once the set has
// grown past it, the capacity Go's append chose; Java's growth policy
// (doubling) can differ, but both are at least the highest bit index ever
// set plus one.
func (b *BitSet) Size() int { return cap(b.words) * 64 }

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
