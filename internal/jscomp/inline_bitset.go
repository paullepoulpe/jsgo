package jscomp

import (
	"encoding/hex"
	"math/bits"
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// inlineBitSet is a compact 24-byte bitset storing bits 0..63 inline in lo and
// bits >= 64 in ext (ext[k] holds bits (k+1)*64 .. (k+2)*64-1).
// Invariant: whenever len(ext) > 0, ext[len(ext)-1] != 0.
type inlineBitSet struct {
	lo  uint64
	ext []uint64
}

func (b *inlineBitSet) trim() {
	i := len(b.ext)
	for i > 0 && b.ext[i-1] == 0 {
		i--
	}
	b.ext = b.ext[:i]
}

func (b *inlineBitSet) ensureExt(extWords int) {
	if extWords <= len(b.ext) {
		return
	}
	if extWords <= cap(b.ext) {
		prev := len(b.ext)
		b.ext = b.ext[:extWords]
		clear(b.ext[prev:])
		return
	}
	newCap := max(extWords, cap(b.ext)*2, 4)
	grown := make([]uint64, extWords, newCap)
	copy(grown, b.ext)
	b.ext = grown
}

func (b *inlineBitSet) wordPtr(w int) *uint64 {
	if w == 0 {
		return &b.lo
	}
	return &b.ext[w-1]
}

func (b *inlineBitSet) Set(bitIndex int) {
	if bitIndex < 0 {
		panic(jutil.NewIndexOutOfBoundsException("bitIndex < 0: " + strconv.Itoa(bitIndex)))
	}
	if bitIndex < 64 {
		b.lo |= uint64(1) << uint(bitIndex)
		return
	}
	w := (bitIndex >> 6) - 1
	b.ensureExt(w + 1)
	b.ext[w] |= uint64(1) << uint(bitIndex&63)
}

func (b *inlineBitSet) Get(bitIndex int) bool {
	if bitIndex < 0 {
		panic(jutil.NewIndexOutOfBoundsException("bitIndex < 0: " + strconv.Itoa(bitIndex)))
	}
	if bitIndex < 64 {
		return (b.lo & (uint64(1) << uint(bitIndex))) != 0
	}
	w := (bitIndex >> 6) - 1
	if w >= len(b.ext) {
		return false
	}
	return (b.ext[w] & (uint64(1) << uint(bitIndex&63))) != 0
}

func (b *inlineBitSet) ClearAll() {
	b.lo = 0
	clear(b.ext)
	b.ext = b.ext[:0]
}

func (b *inlineBitSet) IsEmpty() bool {
	return b.lo == 0 && len(b.ext) == 0
}

func (b *inlineBitSet) WordCapacityBits() int {
	if len(b.ext) > 0 {
		return (len(b.ext) + 1) * 64
	}
	if b.lo != 0 {
		return 64
	}
	return 0
}

func (b *inlineBitSet) Length() int {
	if n := len(b.ext); n > 0 {
		return n*64 + (64 - bits.LeadingZeros64(b.ext[n-1]))
	}
	if b.lo != 0 {
		return 64 - bits.LeadingZeros64(b.lo)
	}
	return 0
}

func (b *inlineBitSet) Cardinality() int {
	n := bits.OnesCount64(b.lo)
	for _, w := range b.ext {
		n += bits.OnesCount64(w)
	}
	return n
}

func (b *inlineBitSet) SetRange(fromIndex, toIndex int) {
	if fromIndex < 0 || toIndex < fromIndex {
		panic(jutil.NewIndexOutOfBoundsException("invalid range"))
	}
	if fromIndex == toIndex {
		return
	}
	endWord := (toIndex - 1) >> 6
	if endWord > 0 {
		b.ensureExt(endWord)
	}
	startWord := fromIndex >> 6
	firstMask := ^uint64(0) << uint(fromIndex&63)
	lastMask := ^uint64(0) >> uint(63-((toIndex-1)&63))
	if startWord == endWord {
		*b.wordPtr(startWord) |= firstMask & lastMask
		return
	}
	*b.wordPtr(startWord) |= firstMask
	for w := startWord + 1; w < endWord; w++ {
		*b.wordPtr(w) = ^uint64(0)
	}
	*b.wordPtr(endWord) |= lastMask
}

func (b *inlineBitSet) FlipRange(fromIndex, toIndex int) {
	if fromIndex < 0 || toIndex < fromIndex {
		panic(jutil.NewIndexOutOfBoundsException("invalid range"))
	}
	if fromIndex == toIndex {
		return
	}
	endWord := (toIndex - 1) >> 6
	if endWord > 0 {
		b.ensureExt(endWord)
	}
	startWord := fromIndex >> 6
	firstMask := ^uint64(0) << uint(fromIndex&63)
	lastMask := ^uint64(0) >> uint(63-((toIndex-1)&63))
	if startWord == endWord {
		*b.wordPtr(startWord) ^= firstMask & lastMask
	} else {
		*b.wordPtr(startWord) ^= firstMask
		for w := startWord + 1; w < endWord; w++ {
			*b.wordPtr(w) = ^(*b.wordPtr(w))
		}
		*b.wordPtr(endWord) ^= lastMask
	}
	b.trim()
}

func (b *inlineBitSet) NextSetBit(fromIndex int) int {
	if fromIndex < 0 {
		panic(jutil.NewIndexOutOfBoundsException("fromIndex < 0: " + strconv.Itoa(fromIndex)))
	}
	if fromIndex < 64 {
		if w := b.lo & (^uint64(0) << uint(fromIndex)); w != 0 {
			return bits.TrailingZeros64(w)
		}
		fromIndex = 64
	}
	u := (fromIndex >> 6) - 1
	if u >= len(b.ext) {
		return -1
	}
	w := b.ext[u] & (^uint64(0) << uint(fromIndex&63))
	for {
		if w != 0 {
			return (u+1)*64 + bits.TrailingZeros64(w)
		}
		u++
		if u >= len(b.ext) {
			return -1
		}
		w = b.ext[u]
	}
}

func (b *inlineBitSet) PreviousSetBit(fromIndex int) int {
	if fromIndex < 0 {
		if fromIndex == -1 {
			return -1
		}
		panic(jutil.NewIndexOutOfBoundsException("fromIndex < -1: " + strconv.Itoa(fromIndex)))
	}
	if fromIndex < 64 {
		w := b.lo & (^uint64(0) >> uint(63-fromIndex))
		if w == 0 {
			return -1
		}
		return 63 - bits.LeadingZeros64(w)
	}
	u := (fromIndex >> 6) - 1
	if u >= len(b.ext) {
		return b.Length() - 1
	}
	w := b.ext[u] & (^uint64(0) >> uint(63-(fromIndex&63)))
	for {
		if w != 0 {
			return (u+1)*64 + 63 - bits.LeadingZeros64(w)
		}
		u--
		if u < 0 {
			break
		}
		w = b.ext[u]
	}
	if b.lo == 0 {
		return -1
	}
	return 63 - bits.LeadingZeros64(b.lo)
}

func (b *inlineBitSet) Or(other *inlineBitSet) {
	if b == other {
		return
	}
	b.lo |= other.lo
	if n := len(other.ext); n > 0 {
		b.ensureExt(n)
		for i, w := range other.ext {
			b.ext[i] |= w
		}
	}
}

func (b *inlineBitSet) And(other *inlineBitSet) {
	if b == other {
		return
	}
	b.lo &= other.lo
	n := min(len(b.ext), len(other.ext))
	clear(b.ext[n:])
	b.ext = b.ext[:n]
	for i := 0; i < n; i++ {
		b.ext[i] &= other.ext[i]
	}
	b.trim()
}

func (b *inlineBitSet) AndNot(other *inlineBitSet) {
	if b == other {
		b.ClearAll()
		return
	}
	b.lo &^= other.lo
	n := min(len(b.ext), len(other.ext))
	for i := 0; i < n; i++ {
		b.ext[i] &^= other.ext[i]
	}
	b.trim()
}

func (b *inlineBitSet) Intersects(other *inlineBitSet) bool {
	if b.lo&other.lo != 0 {
		return true
	}
	n := min(len(b.ext), len(other.ext))
	for i := 0; i < n; i++ {
		if b.ext[i]&other.ext[i] != 0 {
			return true
		}
	}
	return false
}

func (b *inlineBitSet) CopyFrom(src *inlineBitSet) {
	if b == src {
		return
	}
	b.lo = src.lo
	n := len(src.ext)
	if n == 0 {
		clear(b.ext)
		b.ext = b.ext[:0]
		return
	}
	if n <= cap(b.ext) {
		b.ext = b.ext[:n]
	} else {
		b.ext = make([]uint64, n)
	}
	copy(b.ext, src.ext)
}

func (b *inlineBitSet) CopyFromJutil(src *jutil.BitSet) {
	b.ClearAll()
	w0, _ := src.Words128()
	b.lo = w0
	if src.Length() > 64 {
		for i := src.NextSetBit(64); i >= 0; i = src.NextSetBit(i + 1) {
			b.Set(i)
		}
	}
}

func (b *inlineBitSet) HexString() string {
	if len(b.ext) == 0 {
		if b.lo == 0 {
			return ""
		}
		var raw [8]byte
		var dst [16]byte
		length := 0
		for x := b.lo; x != 0; x >>= 8 {
			raw[length] = byte(x)
			length++
		}
		hex.Encode(dst[:length*2], raw[:length])
		return string(dst[:length*2])
	}
	lastIdx := len(b.ext) - 1
	length := 8 * (len(b.ext))
	for x := b.ext[lastIdx]; x != 0; x >>= 8 {
		length++
	}
	raw := make([]byte, length)
	for i := 0; i < 8; i++ {
		raw[i] = byte(b.lo >> (8 * uint(i)))
	}
	for i := 8; i < length; i++ {
		wIdx := (i / 8) - 1
		raw[i] = byte(b.ext[wIdx] >> (8 * uint(i%8)))
	}
	return hex.EncodeToString(raw)
}

func (b *inlineBitSet) String() string {
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
