package jscomp

import (
	"encoding/hex"
	"math/bits"
	"slices"
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// smallBitSet is a compact 32-byte bitset storing bits 0..127 inline in w0/w1
// and bits >= 128 in ext (ext[k] holds bits (k+2)*64 .. (k+3)*64-1).
// Invariant: whenever len(ext) > 0, ext[len(ext)-1] != 0.
type smallBitSet struct {
	w0  uint64
	w1  uint64
	ext []uint64
}

func (b *smallBitSet) trim() {
	i := len(b.ext)
	for i > 0 && b.ext[i-1] == 0 {
		i--
	}
	b.ext = b.ext[:i]
}

func (b *smallBitSet) ensureExt(extWords int) {
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

func (b *smallBitSet) wordPtr(w int) *uint64 {
	if w == 0 {
		return &b.w0
	}
	if w == 1 {
		return &b.w1
	}
	return &b.ext[w-2]
}

func (b *smallBitSet) Set(bitIndex int) {
	if uint(bitIndex) < 64 {
		b.w0 |= uint64(1) << uint(bitIndex)
		return
	}
	if uint(bitIndex) < 128 {
		b.w1 |= uint64(1) << uint(bitIndex-64)
		return
	}
	if bitIndex < 0 {
		panic(jutil.NewIndexOutOfBoundsException("bitIndex < 0: " + strconv.Itoa(bitIndex)))
	}
	w := (bitIndex - 128) >> 6
	b.ensureExt(w + 1)
	b.ext[w] |= uint64(1) << uint(bitIndex&63)
}

func (b *smallBitSet) Get(bitIndex int) bool {
	if uint(bitIndex) < 64 {
		return (b.w0 & (uint64(1) << uint(bitIndex))) != 0
	}
	if uint(bitIndex) < 128 {
		return (b.w1 & (uint64(1) << uint(bitIndex-64))) != 0
	}
	if bitIndex < 0 {
		panic(jutil.NewIndexOutOfBoundsException("bitIndex < 0: " + strconv.Itoa(bitIndex)))
	}
	w := (bitIndex - 128) >> 6
	return w < len(b.ext) && (b.ext[w]&(uint64(1)<<uint(bitIndex&63))) != 0
}

// Reset clears all bits in b, releasing large heap-allocated extension storage
// when release is true.
func (b *smallBitSet) Reset(release bool) {
	b.w0 = 0
	b.w1 = 0
	if release && cap(b.ext) > 16 {
		b.ext = nil
	} else if len(b.ext) > 0 {
		clear(b.ext)
		b.ext = b.ext[:0]
	}
}

func (b *smallBitSet) ClearAll() {
	b.Reset(false)
}

func (b *smallBitSet) IsEmpty() bool {
	return (b.w0|b.w1) == 0 && len(b.ext) == 0
}

func (b *smallBitSet) WordCapacityBits() int {
	if len(b.ext) > 0 {
		return (len(b.ext) + 2) * 64
	}
	if b.w1 != 0 {
		return 128
	}
	if b.w0 != 0 {
		return 64
	}
	return 0
}

func (b *smallBitSet) Length() int {
	if n := len(b.ext); n > 0 {
		return (n+1)*64 + (64 - bits.LeadingZeros64(b.ext[n-1]))
	}
	if b.w1 != 0 {
		return 128 - bits.LeadingZeros64(b.w1)
	}
	if b.w0 != 0 {
		return 64 - bits.LeadingZeros64(b.w0)
	}
	return 0
}

func (b *smallBitSet) Cardinality() int {
	n := bits.OnesCount64(b.w0) + bits.OnesCount64(b.w1)
	for _, w := range b.ext {
		n += bits.OnesCount64(w)
	}
	return n
}

func (b *smallBitSet) SetRange(fromIndex, toIndex int) {
	if fromIndex < 0 || toIndex < fromIndex {
		panic(jutil.NewIndexOutOfBoundsException("invalid range"))
	}
	if fromIndex == toIndex {
		return
	}
	endWord := (toIndex - 1) >> 6
	if endWord >= 2 {
		b.ensureExt(endWord - 1)
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

func (b *smallBitSet) FlipRange(fromIndex, toIndex int) {
	if fromIndex < 0 || toIndex < fromIndex {
		panic(jutil.NewIndexOutOfBoundsException("invalid range"))
	}
	if fromIndex == toIndex {
		return
	}
	endWord := (toIndex - 1) >> 6
	if endWord >= 2 {
		b.ensureExt(endWord - 1)
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

func (b *smallBitSet) NextSetBit(fromIndex int) int {
	if fromIndex < 0 {
		panic(jutil.NewIndexOutOfBoundsException("fromIndex < 0: " + strconv.Itoa(fromIndex)))
	}
	if fromIndex < 64 {
		if w := b.w0 & (^uint64(0) << uint(fromIndex)); w != 0 {
			return bits.TrailingZeros64(w)
		}
		fromIndex = 64
	}
	if fromIndex < 128 {
		if w := b.w1 & (^uint64(0) << uint(fromIndex-64)); w != 0 {
			return 64 + bits.TrailingZeros64(w)
		}
		fromIndex = 128
	}
	u := (fromIndex - 128) >> 6
	if u >= len(b.ext) {
		return -1
	}
	w := b.ext[u] & (^uint64(0) << uint(fromIndex&63))
	for {
		if w != 0 {
			return (u+2)*64 + bits.TrailingZeros64(w)
		}
		u++
		if u >= len(b.ext) {
			return -1
		}
		w = b.ext[u]
	}
}

func (b *smallBitSet) PreviousSetBit(fromIndex int) int {
	if fromIndex < 0 {
		if fromIndex == -1 {
			return -1
		}
		panic(jutil.NewIndexOutOfBoundsException("fromIndex < -1: " + strconv.Itoa(fromIndex)))
	}
	if fromIndex < 64 {
		w := b.w0 & (^uint64(0) >> uint(63-fromIndex))
		if w == 0 {
			return -1
		}
		return 63 - bits.LeadingZeros64(w)
	}
	if fromIndex < 128 {
		if w := b.w1 & (^uint64(0) >> uint(127-fromIndex)); w != 0 {
			return 127 - bits.LeadingZeros64(w)
		}
		if b.w0 == 0 {
			return -1
		}
		return 63 - bits.LeadingZeros64(b.w0)
	}
	u := (fromIndex - 128) >> 6
	if u >= len(b.ext) {
		return b.Length() - 1
	}
	w := b.ext[u] & (^uint64(0) >> uint(63-(fromIndex&63)))
	for {
		if w != 0 {
			return (u+2)*64 + 63 - bits.LeadingZeros64(w)
		}
		u--
		if u < 0 {
			break
		}
		w = b.ext[u]
	}
	if b.w1 != 0 {
		return 127 - bits.LeadingZeros64(b.w1)
	}
	if b.w0 != 0 {
		return 63 - bits.LeadingZeros64(b.w0)
	}
	return -1
}

func (b *smallBitSet) Or(other *smallBitSet) {
	if b == other {
		return
	}
	b.w0 |= other.w0
	b.w1 |= other.w1
	if n := len(other.ext); n > 0 {
		b.ensureExt(n)
		for i, w := range other.ext {
			b.ext[i] |= w
		}
	}
}

func (b *smallBitSet) And(other *smallBitSet) {
	if b == other {
		return
	}
	b.w0 &= other.w0
	b.w1 &= other.w1
	n := min(len(b.ext), len(other.ext))
	clear(b.ext[n:])
	b.ext = b.ext[:n]
	for i := 0; i < n; i++ {
		b.ext[i] &= other.ext[i]
	}
	b.trim()
}

func (b *smallBitSet) AndNot(other *smallBitSet) {
	if b == other {
		b.ClearAll()
		return
	}
	b.w0 &^= other.w0
	b.w1 &^= other.w1
	n := min(len(b.ext), len(other.ext))
	for i := 0; i < n; i++ {
		b.ext[i] &^= other.ext[i]
	}
	b.trim()
}

func (b *smallBitSet) Intersects(other *smallBitSet) bool {
	if (b.w0&other.w0)|(b.w1&other.w1) != 0 {
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

func (b *smallBitSet) ContainsAll(other *smallBitSet) bool {
	if (other.w0&^b.w0)|(other.w1&^b.w1) != 0 {
		return false
	}
	if len(other.ext) > len(b.ext) {
		return false
	}
	for i, ow := range other.ext {
		if ow&^b.ext[i] != 0 {
			return false
		}
	}
	return true
}

func (b *smallBitSet) Equals(other *smallBitSet) bool {
	if b.w0 != other.w0 || b.w1 != other.w1 || len(b.ext) != len(other.ext) {
		return false
	}
	for i, w := range b.ext {
		if w != other.ext[i] {
			return false
		}
	}
	return true
}

func (b *smallBitSet) CopyFrom(src *smallBitSet) {
	if b == src {
		return
	}
	b.w0 = src.w0
	b.w1 = src.w1
	n := len(src.ext)
	if n == 0 {
		if len(b.ext) > 0 {
			clear(b.ext)
			b.ext = b.ext[:0]
		}
		return
	}
	if n <= cap(b.ext) {
		b.ext = b.ext[:n]
	} else {
		b.ext = make([]uint64, n)
	}
	copy(b.ext, src.ext)
}

func (b *smallBitSet) Clone() smallBitSet {
	out := smallBitSet{w0: b.w0, w1: b.w1}
	if len(b.ext) > 0 {
		out.ext = slices.Clone(b.ext)
	}
	return out
}

func (b *smallBitSet) ForEach(fn func(idx int)) {
	for w := b.w0; w != 0; w &= w - 1 {
		fn(bits.TrailingZeros64(w))
	}
	for w := b.w1; w != 0; w &= w - 1 {
		fn(64 + bits.TrailingZeros64(w))
	}
	for wi, word := range b.ext {
		base := 128 + wi*64
		for w := word; w != 0; w &= w - 1 {
			fn(base + bits.TrailingZeros64(w))
		}
	}
}

func (b *smallBitSet) CopyFromJutil(src *jutil.BitSet) {
	b.ClearAll()
	b.w0, b.w1 = src.Words128()
	if src.Length() > 128 {
		for i := src.NextSetBit(128); i >= 0; i = src.NextSetBit(i + 1) {
			b.Set(i)
		}
	}
}

func (b *smallBitSet) HexString() string {
	if len(b.ext) == 0 {
		if b.w1 == 0 {
			if b.w0 == 0 {
				return ""
			}
			var raw [8]byte
			var dst [16]byte
			length := 0
			for x := b.w0; x != 0; x >>= 8 {
				raw[length] = byte(x)
				length++
			}
			hex.Encode(dst[:length*2], raw[:length])
			return string(dst[:length*2])
		}
		var raw [16]byte
		var dst [32]byte
		for i := 0; i < 8; i++ {
			raw[i] = byte(b.w0 >> (8 * uint(i)))
		}
		length := 8
		for x := b.w1; x != 0; x >>= 8 {
			raw[length] = byte(x)
			length++
		}
		hex.Encode(dst[:length*2], raw[:length])
		return string(dst[:length*2])
	}
	lastIdx := len(b.ext) - 1
	length := 16 + 8*lastIdx
	for x := b.ext[lastIdx]; x != 0; x >>= 8 {
		length++
	}
	raw := make([]byte, length)
	for i := 0; i < 8; i++ {
		raw[i] = byte(b.w0 >> (8 * uint(i)))
	}
	for i := 8; i < 16; i++ {
		raw[i] = byte(b.w1 >> (8 * uint(i-8)))
	}
	for i := 16; i < length; i++ {
		wIdx := (i / 8) - 2
		raw[i] = byte(b.ext[wIdx] >> (8 * uint(i%8)))
	}
	return hex.EncodeToString(raw)
}

func (b *smallBitSet) String() string {
	if b.IsEmpty() {
		return "{}"
	}
	buf := make([]byte, 0, 32)
	buf = append(buf, '{')
	first := true
	b.ForEach(func(idx int) {
		if !first {
			buf = append(buf, ',', ' ')
		}
		first = false
		buf = strconv.AppendInt(buf, int64(idx), 10)
	})
	buf = append(buf, '}')
	return string(buf)
}

// SetAndNotOr sets b to (src &^ andNot) | or, reusing b.ext capacity when large enough.
func (b *smallBitSet) SetAndNotOr(src, andNot, or *smallBitSet) {
	b.w0 = (src.w0 &^ andNot.w0) | or.w0
	b.w1 = (src.w1 &^ andNot.w1) | or.w1
	n := max(len(src.ext), len(or.ext))
	if n == 0 {
		if len(b.ext) > 0 {
			clear(b.ext)
			b.ext = b.ext[:0]
		}
		return
	}
	if n <= cap(b.ext) {
		b.ext = b.ext[:n]
	} else {
		b.ext = make([]uint64, n)
	}
	for i := 0; i < n; i++ {
		var sw, kw, gw uint64
		if i < len(src.ext) {
			sw = src.ext[i]
		}
		if i < len(andNot.ext) {
			kw = andNot.ext[i]
		}
		if i < len(or.ext) {
			gw = or.ext[i]
		}
		b.ext[i] = (sw &^ kw) | gw
	}
	b.trim()
}

// JavaHashCode returns java.util.BitSet.hashCode() for b.
func (b *smallBitSet) JavaHashCode() int32 {
	var h int64 = 1234
	h ^= int64(b.w0)
	h ^= int64(b.w1) * 2
	for i, w := range b.ext {
		h ^= int64(w) * int64(i+3)
	}
	return int32((h >> 32) ^ h)
}

// bits128 returns the low 128 bits of b as a fixed-width value bitset.
func (b *smallBitSet) bits128() bitSet128 {
	return bitSet128{b.w0, b.w1}
}

// bitSet128 is a 128-bit value bitset for bounded (< 128-variable) analyses.
type bitSet128 [2]uint64

func (b *bitSet128) set(i int) {
	if i < 64 {
		b[0] |= uint64(1) << uint(i)
	} else {
		b[1] |= uint64(1) << uint(i-64)
	}
}

func (b bitSet128) has(i int) bool {
	if i < 64 {
		return b[0]&(uint64(1)<<uint(i)) != 0
	}
	return b[1]&(uint64(1)<<uint(i-64)) != 0
}

func (b bitSet128) or(o bitSet128) bitSet128 {
	return bitSet128{b[0] | o[0], b[1] | o[1]}
}

func (b *bitSet128) orWith(o bitSet128) {
	b[0] |= o[0]
	b[1] |= o[1]
}

func (b *bitSet128) andWith(o bitSet128) {
	b[0] &= o[0]
	b[1] &= o[1]
}

func (b bitSet128) onesCount() int {
	return bits.OnesCount64(b[0]) + bits.OnesCount64(b[1])
}

func (b bitSet128) forEach(fn func(i int)) {
	for w := b[0]; w != 0; w &= w - 1 {
		fn(bits.TrailingZeros64(w))
	}
	for w := b[1]; w != 0; w &= w - 1 {
		fn(64 + bits.TrailingZeros64(w))
	}
}
