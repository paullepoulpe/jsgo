// Clean-room specification based on Unicode WTF-8, RFC 3629, and Java Language Specification.

package javalang

import (
	"encoding/binary"
	"hash/maphash"
	"math/bits"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"unicode/utf16"
	"unicode/utf8"
)

// StringFromUTF16 converts UTF-16 code units (Java chars) to a WTF-8 Go string.
// Surrogate pairs become valid 4-byte UTF-8; lone surrogates are preserved
// as 3-byte generalized UTF-8 sequences without data loss.
func StringFromUTF16(units []uint16) string {
	if len(units) == 0 {
		return ""
	}
	ascii := true
	for _, u := range units {
		if u >= 0x80 {
			ascii = false
			break
		}
	}
	if ascii {
		if len(units) <= 128 {
			var tmp [128]byte
			for i, u := range units {
				tmp[i] = byte(u)
			}
			return string(tmp[:len(units)])
		}
		var sb strings.Builder
		sb.Grow(len(units))
		var tmp [128]byte
		for i := 0; i < len(units); i += len(tmp) {
			chunk := units[i:min(i+len(tmp), len(units))]
			for j, u := range chunk {
				tmp[j] = byte(u)
			}
			sb.Write(tmp[:len(chunk)])
		}
		return sb.String()
	}
	var sb strings.Builder
	sb.Grow(len(units) * 2)
	for i := 0; i < len(units); i++ {
		u := rune(units[i])
		if utf16.IsSurrogate(u) && u < 0xDC00 && i+1 < len(units) {
			if lo := rune(units[i+1]); lo >= 0xDC00 && lo <= 0xDFFF {
				sb.WriteRune(utf16.DecodeRune(u, lo))
				i++
				continue
			}
		}
		appendWTF8(&sb, u)
	}
	return sb.String()
}

func appendWTF8(sb *strings.Builder, r rune) {
	if r >= 0xD800 && r <= 0xDFFF {
		sb.WriteByte(byte(0xE0 | (r >> 12)))
		sb.WriteByte(byte(0x80 | ((r >> 6) & 0x3F)))
		sb.WriteByte(byte(0x80 | (r & 0x3F)))
		return
	}
	sb.WriteRune(r)
}

const ascii1Bytes = "\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0a\x0b\x0c\x0d\x0e\x0f" +
	"\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f" +
	" !\"#$%&'()*+,-./0123456789:;<=>?" +
	"@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_" +
	"`abcdefghijklmnopqrstuvwxyz{|}~\x7f"

var ascii1Strings = func() [128]string {
	var tbl [128]string
	for i := range tbl {
		tbl[i] = ascii1Bytes[i : i+1]
	}
	return tbl
}()

// StringFromChar returns the one-char Java string for a UTF-16 code unit.
func StringFromChar(c uint16) string {
	if c < 0x80 {
		return ascii1Strings[c]
	}
	var sb strings.Builder
	appendWTF8(&sb, rune(c))
	return sb.String()
}

// StringFromCodePoints implements Java `new String(int[] codePoints, int offset, int count)`.
func StringFromCodePoints(cps []int32, offset, count int32) string {
	if cps == nil {
		panic(NewNullPointerException())
	}
	if offset < 0 || count < 0 || int(offset) > len(cps)-int(count) {
		panic(NewStringIndexOutOfBoundsException("String index out of range"))
	}
	if count == 0 {
		return ""
	}
	sub := cps[offset : offset+count]
	ascii := true
	for _, cp := range sub {
		if uint32(cp) >= 0x80 {
			ascii = false
			break
		}
	}
	if ascii {
		var sb strings.Builder
		sb.Grow(len(sub))
		for _, cp := range sub {
			sb.WriteByte(byte(cp))
		}
		return sb.String()
	}
	units := make([]uint16, 0, len(sub))
	for _, cp := range sub {
		if cp < 0 || cp > 0x10FFFF {
			panic(NewIllegalArgumentException("Invalid Unicode code point"))
		}
		if cp <= 0xFFFF {
			units = append(units, uint16(cp))
		} else {
			hi, lo := utf16.EncodeRune(rune(cp))
			units = append(units, uint16(hi), uint16(lo))
		}
	}
	return StringFromUTF16(units)
}

func decodeWTF8(s string) (r rune, size int) {
	c := s[0]
	if c < 0x80 {
		return rune(c), 1
	}
	if c == 0xED && len(s) >= 3 && s[1] >= 0xA0 && s[1] <= 0xBF && s[2]&0xC0 == 0x80 {
		return rune(c&0x0F)<<12 | rune(s[1]&0x3F)<<6 | rune(s[2]&0x3F), 3
	}
	return utf8.DecodeRuneInString(s)
}

const (
	strCacheMinLen          = 48
	strCacheMinASCIILen     = 512
	strSubstringUnpinMinLen = 1024
	strCacheMaxLen          = 32 << 20
	strCacheShards          = 64
	strCacheWays            = 4
	strMetaBlockUnits       = 64

	unpinDedupMaxLen      = 96
	unpinDedupBuckets     = 16384
	unpinDedupBucketMask  = unpinDedupBuckets - 1
	shortASCIICacheSlots  = 512
	shortASCIICacheMask   = shortASCIICacheSlots - 1

	cpASCIIRunFlag uint32 = 1 << 31
	cpOffsetMask   uint32 = (1 << 31) - 1
)

type unpinBucket struct {
	s0, s1 string
}

// unpinDedup is a bounded 16,384-bucket lock-free 2-way set-associative deduplication
// table (32,768 cached substrings = 128 KiB BSS) for short substrings (<= 96 bytes)
// sliced out of >= 1 KB source buffers, hashed with hardware maphash.String.
var (
	unpinSeed  = maphash.MakeSeed()
	unpinDedup [unpinDedupBuckets]atomic.Pointer[unpinBucket]
)

// unpinSubstring clones small substrings sliced out of >= 1 KB source buffers
// (< 1 KB or < 1/4 of parentLen), preventing identifiers, JSDoc tags/comments,
// or module paths from pinning entire JS SourceFile backing byte arrays in live heap,
// while deduplicating repeated short tokens (<= 96 bytes) lock-free.
func unpinSubstring(parentLen int, sub string) string {
	if len(sub) == 0 {
		return ""
	}
	if parentLen >= strSubstringUnpinMinLen {
		return unpinSubstringFromParent(parentLen, sub)
	}
	return sub
}

//go:noinline
func unpinSubstringFromParent(parentLen int, sub string) string {
	if len(sub) == 1 && sub[0] < 0x80 {
		return ascii1Strings[sub[0]]
	}
	if len(sub) < parentLen && (len(sub) < strSubstringUnpinMinLen || len(sub)*4 < parentLen) {
		if len(sub) > unpinDedupMaxLen {
			return strings.Clone(sub)
		}
		idx := maphash.String(unpinSeed, sub) & unpinDedupBucketMask
		b := unpinDedup[idx].Load()
		if b != nil {
			if b.s0 == sub {
				return b.s0
			}
			if b.s1 == sub {
				return b.s1
			}
		}
		return unpinDedupMiss(idx, b, sub)
	}
	return sub
}

//go:noinline
func unpinDedupMiss(idx uint64, prev *unpinBucket, sub string) string {
	cloned := strings.Clone(sub)
	var prevS0 string
	if prev != nil {
		prevS0 = prev.s0
	}
	unpinDedup[idx].Store(&unpinBucket{s0: cloned, s1: prevS0})
	return cloned
}

type strCheckpoint struct {
	uOff uint32 // UTF-16 code unit offset; cpASCIIRunFlag indicates pure ASCII until next checkpoint
	bOff uint32 // byte offset in s (always at a rune boundary)
}

type strMeta struct {
	s           string
	hash        uint64
	uAscii      uint            // uint(asciiLen) for zero-cast inline bounds check
	asciiLen    int             // leading pure-ASCII byte/UTF-16 count (== len(s) iff pure ASCII)
	asciiTailU  int             // UTF-16 offset where trailing pure-ASCII suffix begins
	asciiTailB  int             // byte offset where trailing pure-ASCII suffix begins
	utf16Len    int             // total UTF-16 code unit length
	checkpoints []strCheckpoint // sparse checkpoints across [asciiLen, asciiTailB); nil when pure ASCII or <= 64 units
	cleanUTF8   bool            // true iff s is valid canonical UTF-8 (no lone surrogates and no invalid UTF-8 bytes)
}

type strMetaShard struct {
	mu    sync.Mutex
	slots [strCacheWays]atomic.Pointer[strMeta]
}

type shortASCIIBucket struct {
	s0, s1, s2, s3 string
}

var (
	sentinelStrMeta = strMeta{cleanUTF8: true}
	strMetaCache    [strCacheShards]strMetaShard
	hotStrMeta0     atomic.Pointer[strMeta]
	hotStrMeta1     atomic.Pointer[strMeta]

	shortASCIICache [shortASCIICacheSlots]atomic.Pointer[shortASCIIBucket]
	shortASCIISeen  [shortASCIICacheSlots]atomic.Uint64
)

func init() {
	hotStrMeta0.Store(&sentinelStrMeta)
	hotStrMeta1.Store(&sentinelStrMeta)
}

const asciiHighMask64 uint64 = 0x8080808080808080

// isASCIIUpTo32 reports whether s (with len(s) <= 32) is pure ASCII using
// overlapping word loads with zero loops.
func isASCIIUpTo32(s string) bool {
	n := len(s)
	if n >= 16 {
		w := binary.LittleEndian.Uint64([]byte(s[:8])) |
			binary.LittleEndian.Uint64([]byte(s[8:16])) |
			binary.LittleEndian.Uint64([]byte(s[n-16:n-8])) |
			binary.LittleEndian.Uint64([]byte(s[n-8:]))
		return (w & asciiHighMask64) == 0
	}
	if n >= 8 {
		w := binary.LittleEndian.Uint64([]byte(s[:8])) |
			binary.LittleEndian.Uint64([]byte(s[n-8:]))
		return (w & asciiHighMask64) == 0
	}
	if n >= 4 {
		w := binary.LittleEndian.Uint32([]byte(s[:4])) |
			binary.LittleEndian.Uint32([]byte(s[n-4:]))
		return (w & 0x80808080) == 0
	}
	if n > 0 {
		return (s[0] | s[n>>1] | s[n-1]) < 0x80
	}
	return true
}

// isASCIIUpTo64 reports whether s (with len(s) <= 64) is pure ASCII using
// overlapping 32-byte head/tail loads with zero loops.
func isASCIIUpTo64(s string) bool {
	n := len(s)
	if n <= 32 {
		return isASCIIUpTo32(s)
	}
	wHead := binary.LittleEndian.Uint64([]byte(s[:8])) |
		binary.LittleEndian.Uint64([]byte(s[8:16])) |
		binary.LittleEndian.Uint64([]byte(s[16:24])) |
		binary.LittleEndian.Uint64([]byte(s[24:32]))
	wTail := binary.LittleEndian.Uint64([]byte(s[n-32:n-24])) |
		binary.LittleEndian.Uint64([]byte(s[n-24:n-16])) |
		binary.LittleEndian.Uint64([]byte(s[n-16:n-8])) |
		binary.LittleEndian.Uint64([]byte(s[n-8:]))
	return ((wHead | wTail) & asciiHighMask64) == 0
}

// shortASCIIHash computes a fast 64-bit fingerprint for s with len(s) > 64.
func shortASCIIHash(s string) uint64 {
	n := len(s)
	mid := n >> 1
	h := uint64(n)*0x9e3779b97f4a7c15 ^
		binary.LittleEndian.Uint64([]byte(s[:8])) ^
		(binary.LittleEndian.Uint64([]byte(s[mid:mid+8])) * 0xbf58476d1ce4e5b9) ^
		(binary.LittleEndian.Uint64([]byte(s[n-8:])) * 0x94d049bb133111eb)
	return h ^ (h >> 32)
}

func checkShortASCII(s string) bool {
	b := shortASCIICache[shortASCIIHash(s)&shortASCIICacheMask].Load()
	return b != nil && (b.s0 == s || b.s1 == s || b.s2 == s || b.s3 == s)
}

// isShortASCII reports whether s (with 64 < len(s) < strCacheMinASCIILen) is pure ASCII,
// using a lock-free 4-way set-associative atomic.Pointer cache so O(N) character loops
// (e.g. JsDocTokenStream and Wtf8.encodeToWtf8) hit O(1) string-header equality with
// zero mutexes, zero atomic RMWs, and zero allocations on hits (and zero allocations
// on single-probe strings via shortASCIISeen).
func isShortASCII(s string) bool {
	h := shortASCIIHash(s)
	idx := h & shortASCIICacheMask
	b := shortASCIICache[idx].Load()
	if b != nil && (b.s0 == s || b.s1 == s || b.s2 == s || b.s3 == s) {
		return true
	}
	return isShortASCIIMiss(s, h, idx, b)
}

//go:noinline
func isShortASCIIMiss(s string, h, idx uint64, prev *shortASCIIBucket) bool {
	if asciiPrefixLenRaw(s) != len(s) {
		return false
	}
	tag := h | 1
	if shortASCIISeen[idx].Load() != tag {
		shortASCIISeen[idx].Store(tag)
		return true
	}
	nb := &shortASCIIBucket{s0: s}
	if prev != nil {
		nb.s1 = prev.s0
		nb.s2 = prev.s1
		nb.s3 = prev.s2
	}
	shortASCIICache[idx].Store(nb)
	return true
}

func (m *strMeta) isASCII() bool {
	return m.asciiLen == len(m.s)
}

func asciiPrefixLenRaw(s string) int {
	i := 0
	for ; i+32 <= len(s); i += 32 {
		w0 := binary.LittleEndian.Uint64([]byte(s[i : i+8]))
		w1 := binary.LittleEndian.Uint64([]byte(s[i+8 : i+16]))
		w2 := binary.LittleEndian.Uint64([]byte(s[i+16 : i+24]))
		w3 := binary.LittleEndian.Uint64([]byte(s[i+24 : i+32]))
		if ((w0 | w1 | w2 | w3) & asciiHighMask64) != 0 {
			if m := w0 & asciiHighMask64; m != 0 {
				return i + (bits.TrailingZeros64(m) >> 3)
			}
			if m := w1 & asciiHighMask64; m != 0 {
				return i + 8 + (bits.TrailingZeros64(m) >> 3)
			}
			if m := w2 & asciiHighMask64; m != 0 {
				return i + 16 + (bits.TrailingZeros64(m) >> 3)
			}
			return i + 24 + (bits.TrailingZeros64(w3&asciiHighMask64) >> 3)
		}
	}
	for ; i+8 <= len(s); i += 8 {
		if m := binary.LittleEndian.Uint64([]byte(s[i:i+8])) & asciiHighMask64; m != 0 {
			return i + (bits.TrailingZeros64(m) >> 3)
		}
	}
	for ; i < len(s); i++ {
		if s[i] >= 0x80 {
			return i
		}
	}
	return len(s)
}

func stringToUTF16NonASCII(s string, prefix int) []uint16 {
	out := make([]uint16, prefix, len(s))
	for i := 0; i < prefix; i++ {
		out[i] = uint16(s[i])
	}
	for i := prefix; i < len(s); {
		r, n := decodeWTF8(s[i:])
		i += n
		if r >= 0x10000 {
			hi, lo := utf16.EncodeRune(r)
			out = append(out, uint16(hi), uint16(lo))
		} else {
			out = append(out, uint16(r))
		}
	}
	return out
}

func buildNonASCIIMeta(s string, h uint64, prefix int) *strMeta {
	tailB := len(s)
	for tailB-8 > prefix {
		if (s[tailB-8] | s[tailB-7] | s[tailB-6] | s[tailB-5] | s[tailB-4] | s[tailB-3] | s[tailB-2] | s[tailB-1]) >= 0x80 {
			break
		}
		tailB -= 8
	}
	for tailB > prefix && s[tailB-1] < 0x80 {
		tailB--
	}

	var cps []strCheckpoint
	u := prefix
	lastCheckU := prefix
	cleanUTF8 := true
	for i := prefix; i < tailB; {
		if s[i] < 0x80 {
			j := i + 1
			for j+8 <= tailB && (s[j]|s[j+1]|s[j+2]|s[j+3]|s[j+4]|s[j+5]|s[j+6]|s[j+7]) < 0x80 {
				j += 8
			}
			for j < tailB && s[j] < 0x80 {
				j++
			}
			asciiRun := j - i
			if asciiRun >= strMetaBlockUnits {
				cps = append(cps, strCheckpoint{
					uOff: uint32(u) | cpASCIIRunFlag,
					bOff: uint32(i),
				})
				u += asciiRun
				i = j
				cps = append(cps, strCheckpoint{
					uOff: uint32(u),
					bOff: uint32(i),
				})
				lastCheckU = u
				continue
			}
			for i < j {
				if u-lastCheckU >= strMetaBlockUnits {
					cps = append(cps, strCheckpoint{uOff: uint32(u), bOff: uint32(i)})
					lastCheckU = u
				}
				step := min(j-i, strMetaBlockUnits-(u-lastCheckU))
				u += step
				i += step
			}
			continue
		}
		if u-lastCheckU >= strMetaBlockUnits {
			cps = append(cps, strCheckpoint{uOff: uint32(u), bOff: uint32(i)})
			lastCheckU = u
		}
		r, n := decodeWTF8(s[i:])
		if (r >= 0xD800 && r <= 0xDFFF) || (r == utf8.RuneError && n == 1) {
			cleanUTF8 = false
		}
		i += n
		if r >= 0x10000 {
			u += 2
		} else {
			u++
		}
	}

	tailU := u
	return &strMeta{
		s:           s,
		hash:        h,
		uAscii:      uint(prefix),
		asciiLen:    prefix,
		asciiTailU:  tailU,
		asciiTailB:  tailB,
		utf16Len:    tailU + (len(s) - tailB),
		checkpoints: cps,
		cleanUTF8:   cleanUTF8,
	}
}

func (m *strMeta) seekUTF16Checkpoint(uIdx int) (curU, curB int, asciiRun bool) {
	cps := m.checkpoints
	if len(cps) == 0 || uint32(uIdx) < (cps[0].uOff&cpOffsetMask) {
		return m.asciiLen, m.asciiLen, false
	}
	lo, hi := 0, len(cps)-1
	target := uint32(uIdx)
	for lo < hi {
		mid := int(uint(lo+hi+1) >> 1)
		if (cps[mid].uOff & cpOffsetMask) <= target {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	cp := cps[lo]
	return int(cp.uOff & cpOffsetMask), int(cp.bOff), cp.uOff&cpASCIIRunFlag != 0
}

func (m *strMeta) charAt(uIdx int) uint16 {
	if uint(uIdx) >= uint(m.utf16Len) {
		panic(NewStringIndexOutOfBoundsExceptionForIndex(uIdx, m.utf16Len))
	}
	s := m.s
	if uIdx < m.asciiLen {
		return uint16(s[uIdx])
	}
	if uIdx >= m.asciiTailU {
		return uint16(s[m.asciiTailB+(uIdx-m.asciiTailU)])
	}
	curU, curB, asciiRun := m.seekUTF16Checkpoint(uIdx)
	if asciiRun {
		return uint16(s[curB+(uIdx-curU)])
	}
	for {
		r, n := decodeWTF8(s[curB:])
		if r >= 0x10000 {
			if uIdx == curU {
				hi, _ := utf16.EncodeRune(r)
				return uint16(hi)
			}
			if uIdx == curU+1 {
				_, lo := utf16.EncodeRune(r)
				return uint16(lo)
			}
			curU += 2
		} else {
			if uIdx == curU {
				return uint16(r)
			}
			curU++
		}
		curB += n
	}
}

// resolveUTF16Boundary resolves UTF-16 boundary index uIdx (0 <= uIdx <= m.utf16Len):
//   - If uIdx falls on a rune boundary, splitHi == 0 and splitLo == 0, and runeStartB == runeEndB.
//   - If uIdx splits a 4-byte supplementary rune [runeStartB : runeEndB], splitHi and splitLo
//     are that rune's high and low UTF-16 surrogates.
func (m *strMeta) resolveUTF16Boundary(uIdx int) (runeStartB, runeEndB int, splitHi, splitLo uint16) {
	if uIdx <= m.asciiLen {
		return uIdx, uIdx, 0, 0
	}
	if uIdx >= m.asciiTailU {
		b := m.asciiTailB + (uIdx - m.asciiTailU)
		return b, b, 0, 0
	}
	curU, curB, asciiRun := m.seekUTF16Checkpoint(uIdx)
	if asciiRun {
		b := curB + (uIdx - curU)
		return b, b, 0, 0
	}
	s := m.s
	if m.cleanUTF8 {
		for curU < uIdx {
			b := s[curB]
			switch {
			case b < 0x80:
				curB++
				curU++
			case b < 0xE0:
				curB += 2
				curU++
			case b < 0xF0:
				curB += 3
				curU++
			default:
				if uIdx == curU+1 {
					r, sz := utf8.DecodeRuneInString(s[curB:])
					hi, lo := utf16.EncodeRune(r)
					return curB, curB + sz, uint16(hi), uint16(lo)
				}
				curB += 4
				curU += 2
			}
		}
		return curB, curB, 0, 0
	}
	for curU < uIdx {
		r, sz := decodeWTF8(s[curB:])
		if r >= 0x10000 {
			if uIdx == curU+1 {
				hi, lo := utf16.EncodeRune(r)
				return curB, curB + sz, uint16(hi), uint16(lo)
			}
			curU += 2
		} else {
			curU++
		}
		curB += sz
	}
	return curB, curB, 0, 0
}

func (m *strMeta) substring(begin, end int) string {
	if begin < 0 || begin > end || end > m.utf16Len {
		panic(NewStringIndexOutOfBoundsExceptionForRange(begin, end, m.utf16Len))
	}
	if begin == end {
		return ""
	}
	if end <= m.asciiLen {
		return unpinSubstring(len(m.s), m.s[begin:end])
	}
	if begin >= m.asciiTailU {
		b0 := m.asciiTailB + (begin - m.asciiTailU)
		return unpinSubstring(len(m.s), m.s[b0:b0+(end-begin)])
	}
	if m.cleanUTF8 {
		_, bAfterBegin, _, beginLo := m.resolveUTF16Boundary(begin)
		bBeforeEnd, _, endHi, _ := m.resolveUTF16Boundary(end)
		if beginLo == 0 && endHi == 0 {
			return unpinSubstring(len(m.s), m.s[bAfterBegin:bBeforeEnd])
		}
		if bAfterBegin > bBeforeEnd {
			// Cannot happen when begin < end: same supplementary rune has begin==curU (beginLo==0) or end==curU+2 (endHi==0).
			return ""
		}
		extra := 0
		if beginLo != 0 {
			extra += 3
		}
		if endHi != 0 {
			extra += 3
		}
		var sb strings.Builder
		sb.Grow(extra + (bBeforeEnd - bAfterBegin))
		if beginLo != 0 {
			appendWTF8(&sb, rune(beginLo))
		}
		sb.WriteString(m.s[bAfterBegin:bBeforeEnd])
		if endHi != 0 {
			appendWTF8(&sb, rune(endHi))
		}
		return sb.String()
	}
	units := make([]uint16, 0, end-begin)
	_, bPos, _, beginLo := m.resolveUTF16Boundary(begin)
	if beginLo != 0 {
		units = append(units, beginLo)
	}
	s := m.s
	for len(units) < end-begin {
		r, sz := decodeWTF8(s[bPos:])
		bPos += sz
		if r >= 0x10000 {
			hi, lo := utf16.EncodeRune(r)
			units = append(units, uint16(hi))
			if len(units) < end-begin {
				units = append(units, uint16(lo))
			}
		} else {
			units = append(units, uint16(r))
		}
	}
	return StringFromUTF16(units)
}

// utf16IndexAtByte returns the UTF-16 code unit offset corresponding to rune boundary bIdx.
func (m *strMeta) utf16IndexAtByte(bIdx int) int {
	if bIdx <= m.asciiLen {
		return bIdx
	}
	if bIdx >= m.asciiTailB {
		if bIdx >= len(m.s) {
			return m.utf16Len
		}
		return m.asciiTailU + (bIdx - m.asciiTailB)
	}
	cps := m.checkpoints
	curU, curB := m.asciiLen, m.asciiLen
	if len(cps) > 0 && uint32(bIdx) >= cps[0].bOff {
		lo, hi := 0, len(cps)-1
		target := uint32(bIdx)
		for lo < hi {
			mid := int(uint(lo+hi+1) >> 1)
			if cps[mid].bOff <= target {
				lo = mid
			} else {
				hi = mid - 1
			}
		}
		cp := cps[lo]
		curU = int(cp.uOff & cpOffsetMask)
		curB = int(cp.bOff)
		if cp.uOff&cpASCIIRunFlag != 0 {
			return curU + (bIdx - curB)
		}
	}
	s := m.s
	if m.cleanUTF8 {
		for curB < bIdx {
			b := s[curB]
			switch {
			case b < 0x80:
				curB++
				curU++
			case b < 0xE0:
				curB += 2
				curU++
			case b < 0xF0:
				curB += 3
				curU++
			default:
				curB += 4
				curU += 2
			}
		}
		return curU
	}
	for curB < bIdx {
		r, sz := decodeWTF8(s[curB:])
		curB += sz
		if r >= 0x10000 {
			curU += 2
		} else {
			curU++
		}
	}
	return curU
}

func hashStrMetaKey(s string) uint64 {
	n := len(s)
	h := uint64(n) * 0x9e3779b97f4a7c15
	if n >= 16 {
		h = (h ^ binary.LittleEndian.Uint64([]byte(s[:8]))) * 0xbf58476d1ce4e5b9
		h = (h ^ binary.LittleEndian.Uint64([]byte(s[8:16]))) * 0x94d049bb133111eb
		h = (h ^ binary.LittleEndian.Uint64([]byte(s[n-16:n-8]))) * 0xbf58476d1ce4e5b9
		h = (h ^ binary.LittleEndian.Uint64([]byte(s[n-8:]))) * 0x94d049bb133111eb
		mid := n >> 1
		h = (h ^ binary.LittleEndian.Uint64([]byte(s[mid:mid+8]))) * 0x9e3779b97f4a7c15
		return h ^ (h >> 30) ^ (h >> 53)
	}
	for i := 0; i < n; i++ {
		h = (h ^ uint64(s[i])) * 0xbf58476d1ce4e5b9
	}
	return h ^ (h >> 30) ^ (h >> 53)
}

func promoteHotStrMeta(m *strMeta) *strMeta {
	if prev := hotStrMeta0.Load(); prev != m {
		hotStrMeta1.Store(prev)
		hotStrMeta0.Store(m)
	}
	return m
}

func lookupStrMetaCache(s string, h uint64) *strMeta {
	sh := &strMetaCache[h&(strCacheShards-1)]
	for i := 0; i < strCacheWays; i++ {
		if m := sh.slots[i].Load(); m != nil && m.hash == h && m.s == s {
			return promoteHotStrMeta(m)
		}
	}
	return nil
}

func storeStrMetaCache(m *strMeta) *strMeta {
	sh := &strMetaCache[m.hash&(strCacheShards-1)]
	sh.mu.Lock()
	for i := 0; i < strCacheWays; i++ {
		if ex := sh.slots[i].Load(); ex != nil && ex.hash == m.hash && ex.s == m.s {
			sh.mu.Unlock()
			return promoteHotStrMeta(ex)
		}
	}
	for i := strCacheWays - 1; i > 0; i-- {
		sh.slots[i].Store(sh.slots[i-1].Load())
	}
	sh.slots[0].Store(m)
	sh.mu.Unlock()
	return promoteHotStrMeta(m)
}

// tryGetStrMeta returns cached metadata for s when len(s) is in [strCacheMinLen, strCacheMaxLen].
// Pure-ASCII strings shorter than strCacheMinASCIILen return nil after a fast register scan so
// qualified type/identifier names never allocate *strMeta or evict source-file entries.
func tryGetStrMeta(s string) *strMeta {
	n := len(s)
	if n < strCacheMinLen || n > strCacheMaxLen {
		return nil
	}
	if m0 := hotStrMeta0.Load(); m0.s == s {
		return m0
	}
	if m1 := hotStrMeta1.Load(); m1.s == s {
		return m1
	}
	if n < strCacheMinASCIILen {
		if n <= 64 {
			if isASCIIUpTo64(s) {
				return nil
			}
		} else if checkShortASCII(s) {
			return nil
		}
		prefix := asciiPrefixLenRaw(s)
		if prefix == n {
			return nil
		}
		h := hashStrMetaKey(s)
		if m := lookupStrMetaCache(s, h); m != nil {
			return m
		}
		return storeStrMetaCache(buildNonASCIIMeta(s, h, prefix))
	}
	h := hashStrMetaKey(s)
	if m := lookupStrMetaCache(s, h); m != nil {
		return m
	}
	prefix := asciiPrefixLenRaw(s)
	var m *strMeta
	if prefix == n {
		m = &strMeta{
			s:          s,
			hash:       h,
			uAscii:     uint(n),
			asciiLen:   n,
			asciiTailU: n,
			asciiTailB: n,
			utf16Len:   n,
			cleanUTF8:  true,
		}
	} else {
		m = buildNonASCIIMeta(s, h, prefix)
	}
	return storeStrMetaCache(m)
}

func stringUTF16View(s string) []uint16 {
	return StringToUTF16(s)
}

// ASCIIPrefixLen returns the length of the leading ASCII bytes in s.
func ASCIIPrefixLen(s string) int {
	if len(s) == 0 || s[0] >= 0x80 {
		return 0
	}
	if m := hotStrMeta0.Load(); m.s == s {
		return m.asciiLen
	}
	if len(s) >= strCacheMinASCIILen && len(s) <= strCacheMaxLen {
		return tryGetStrMeta(s).asciiLen
	}
	return asciiPrefixLenRaw(s)
}

// IsASCII reports whether s contains only ASCII characters.
func IsASCII(s string) bool {
	n := len(s)
	if n <= 64 {
		return isASCIIUpTo64(s)
	}
	if s[0] >= 0x80 {
		return false
	}
	if n < strCacheMinASCIILen {
		if checkShortASCII(s) {
			return true
		}
		return asciiPrefixLenRaw(s) == n
	}
	if m := hotStrMeta0.Load(); m.s == s {
		return m.asciiLen == n
	}
	if n <= strCacheMaxLen {
		return tryGetStrMeta(s).isASCII()
	}
	return asciiPrefixLenRaw(s) == n
}

// StringToUTF16 converts a WTF-8 Go string to UTF-16 code units.
func StringToUTF16(s string) []uint16 {
	prefix := asciiPrefixLenRaw(s)
	if prefix == len(s) {
		out := make([]uint16, len(s))
		for i := 0; i < len(s); i++ {
			out[i] = uint16(s[i])
		}
		return out
	}
	return stringToUTF16NonASCII(s, prefix)
}

func utf16LenRaw(s string) int {
	prefix := asciiPrefixLenRaw(s)
	if prefix == len(s) {
		return len(s)
	}
	n := prefix
	for i := prefix; i < len(s); {
		r, size := decodeWTF8(s[i:])
		i += size
		if r >= 0x10000 {
			n += 2
		} else {
			n++
		}
	}
	return n
}

// StringLength returns the number of UTF-16 code units in s (Java String.length()).
func StringLength(s string) int {
	if len(s) <= 1 {
		return len(s)
	}
	return stringLengthSlow(s)
}

//go:noinline
func stringLengthSlow(s string) int {
	if m := hotStrMeta0.Load(); m.s == s {
		return m.utf16Len
	}
	n := len(s)
	if n <= 64 {
		if isASCIIUpTo64(s) {
			return n
		}
	} else if n < strCacheMinASCIILen {
		if checkShortASCII(s) {
			return n
		}
	}
	if n < strCacheMinASCIILen {
		prefix := asciiPrefixLenRaw(s)
		if prefix == n {
			return n
		}
		if n >= strCacheMinLen {
			h := hashStrMetaKey(s)
			if m := lookupStrMetaCache(s, h); m != nil {
				return m.utf16Len
			}
			return storeStrMetaCache(buildNonASCIIMeta(s, h, prefix)).utf16Len
		}
		u := prefix
		for i := prefix; i < n; {
			r, size := decodeWTF8(s[i:])
			i += size
			if r >= 0x10000 {
				u += 2
			} else {
				u++
			}
		}
		return u
	}
	if m := tryGetStrMeta(s); m != nil {
		return m.utf16Len
	}
	return utf16LenRaw(s)
}

// StringCharAt returns the UTF-16 code unit at index (Java String.charAt(int)).
func StringCharAt(s string, index int) uint16 {
	if m := hotStrMeta0.Load(); uint(index) < m.uAscii && m.s == s {
		return uint16(s[index])
	}
	if uint(index) < uint(len(s)) && s[index] < 0x80 {
		if index == 0 {
			return uint16(s[0])
		}
		if len(s) < strCacheMinASCIILen {
			if index < 32 {
				if isASCIIUpTo32(s[:index+1]) {
					return uint16(s[index])
				}
			} else if len(s) <= 64 {
				if isASCIIUpTo64(s[:index+1]) {
					return uint16(s[index])
				}
			} else {
				b := shortASCIICache[shortASCIIHash(s)&shortASCIICacheMask].Load()
				if (b != nil && (b.s0 == s || b.s1 == s || b.s2 == s || b.s3 == s)) || isShortASCII(s) {
					return uint16(s[index])
				}
			}
		}
	}
	return stringCharAtSlow(s, index)
}

//go:noinline
func stringCharAtSlow(s string, index int) uint16 {
	if m := tryGetStrMeta(s); m != nil {
		return m.charAt(index)
	}
	if len(s) >= strCacheMinLen && len(s) < strCacheMinASCIILen {
		if uint(index) >= uint(len(s)) {
			panic(NewStringIndexOutOfBoundsExceptionForIndex(index, len(s)))
		}
		return uint16(s[index])
	}
	if index < 0 {
		panic(NewStringIndexOutOfBoundsExceptionForIndex(index, StringLength(s)))
	}
	prefix := asciiPrefixLenRaw(s)
	if index < prefix {
		return uint16(s[index])
	}
	if prefix == len(s) {
		panic(NewStringIndexOutOfBoundsExceptionForIndex(index, len(s)))
	}
	uIdx := prefix
	for i := prefix; i < len(s); {
		r, size := decodeWTF8(s[i:])
		if r >= 0x10000 {
			if index == uIdx {
				hi, _ := utf16.EncodeRune(r)
				return uint16(hi)
			}
			if index == uIdx+1 {
				_, lo := utf16.EncodeRune(r)
				return uint16(lo)
			}
			uIdx += 2
		} else {
			if index == uIdx {
				return uint16(r)
			}
			uIdx++
		}
		i += size
	}
	panic(NewStringIndexOutOfBoundsExceptionForIndex(index, uIdx))
}

// StringCodePointAt returns the Unicode code point at UTF-16 index (Java String.codePointAt(int)).
func StringCodePointAt(s string, index int) int32 {
	if m := hotStrMeta0.Load(); uint(index) < m.uAscii && m.s == s {
		return int32(s[index])
	}
	if uint(index) < uint(len(s)) && s[index] < 0x80 {
		if index == 0 {
			return int32(s[0])
		}
		if len(s) < strCacheMinASCIILen {
			if index < 32 {
				if isASCIIUpTo32(s[:index+1]) {
					return int32(s[index])
				}
			} else if len(s) <= 64 {
				if isASCIIUpTo64(s[:index+1]) {
					return int32(s[index])
				}
			} else {
				b := shortASCIICache[shortASCIIHash(s)&shortASCIICacheMask].Load()
				if (b != nil && (b.s0 == s || b.s1 == s || b.s2 == s || b.s3 == s)) || isShortASCII(s) {
					return int32(s[index])
				}
			}
		}
	}
	return stringCodePointAtSlow(s, index)
}

//go:noinline
func stringCodePointAtSlow(s string, index int) int32 {
	if m := tryGetStrMeta(s); m != nil {
		c1 := m.charAt(index)
		if c1 >= 0xD800 && c1 <= 0xDBFF && index+1 < m.utf16Len {
			c2 := m.charAt(index + 1)
			if c2 >= 0xDC00 && c2 <= 0xDFFF {
				return 0x10000 + ((int32(c1) - 0xD800) << 10) + (int32(c2) - 0xDC00)
			}
		}
		return int32(c1)
	}
	if len(s) >= strCacheMinLen && len(s) < strCacheMinASCIILen {
		if uint(index) >= uint(len(s)) {
			panic(NewStringIndexOutOfBoundsExceptionForIndex(index, len(s)))
		}
		return int32(s[index])
	}
	if index < 0 {
		panic(NewStringIndexOutOfBoundsExceptionForIndex(index, StringLength(s)))
	}
	prefix := asciiPrefixLenRaw(s)
	if index < prefix {
		return int32(s[index])
	}
	if prefix == len(s) {
		panic(NewStringIndexOutOfBoundsExceptionForIndex(index, len(s)))
	}
	uIdx := prefix
	for i := prefix; i < len(s); {
		r, size := decodeWTF8(s[i:])
		if r >= 0x10000 {
			if index == uIdx {
				return int32(r)
			}
			if index == uIdx+1 {
				_, lo := utf16.EncodeRune(r)
				if i+size < len(s) {
					r2, _ := decodeWTF8(s[i+size:])
					if lo >= 0xD800 && lo <= 0xDBFF && r2 >= 0xDC00 && r2 <= 0xDFFF {
						return 0x10000 + ((int32(lo) - 0xD800) << 10) + (int32(r2) - 0xDC00)
					}
				}
				return int32(lo)
			}
			uIdx += 2
		} else {
			if index == uIdx {
				if r >= 0xD800 && r <= 0xDBFF && i+size < len(s) {
					r2, _ := decodeWTF8(s[i+size:])
					if r2 >= 0xDC00 && r2 <= 0xDFFF {
						return 0x10000 + ((int32(r) - 0xD800) << 10) + (int32(r2) - 0xDC00)
					}
				}
				return int32(r)
			}
			uIdx++
		}
		i += size
	}
	panic(NewStringIndexOutOfBoundsExceptionForIndex(index, uIdx))
}

// StringSubstring returns the substring from begin to end in UTF-16 code unit offsets.
func StringSubstring(s string, begin, end int) string {
	if m := hotStrMeta0.Load(); uint(end) <= m.uAscii && uint(begin) <= uint(end) && m.s == s {
		if end-begin <= 1 {
			if begin == end {
				return ""
			}
			return ascii1Strings[s[begin]]
		}
		return unpinSubstring(len(s), s[begin:end])
	}
	if begin >= 0 && begin <= end && end <= len(s) && len(s) < strCacheMinASCIILen {
		if (end <= 64 && isASCIIUpTo64(s[:end])) || (len(s) > 64 && checkShortASCII(s)) || asciiPrefixLenRaw(s[:end]) == end {
			if begin == end {
				return ""
			}
			return s[begin:end]
		}
	}
	if m := tryGetStrMeta(s); m != nil {
		return m.substring(begin, end)
	}
	if len(s) >= strCacheMinLen && len(s) < strCacheMinASCIILen {
		if begin < 0 || begin > end || end > len(s) {
			panic(NewStringIndexOutOfBoundsExceptionForRange(begin, end, len(s)))
		}
		return s[begin:end]
	}
	if begin < 0 || begin > end {
		panic(NewStringIndexOutOfBoundsExceptionForRange(begin, end, StringLength(s)))
	}
	prefix := asciiPrefixLenRaw(s)
	if end <= prefix {
		return unpinSubstring(len(s), s[begin:end])
	}
	if prefix == len(s) {
		panic(NewStringIndexOutOfBoundsExceptionForRange(begin, end, len(s)))
	}
	var (
		byteStart    int
		haveBegin    bool
		splitBeginLo rune
		uIdx         = prefix
	)
	if begin <= prefix {
		byteStart = begin
		haveBegin = true
	}
	for i := prefix; i < len(s); {
		r, size := decodeWTF8(s[i:])
		if !haveBegin {
			if r >= 0x10000 {
				if begin == uIdx {
					byteStart = i
					haveBegin = true
				} else if begin == uIdx+1 {
					if end == begin {
						return ""
					}
					_, splitBeginLo = utf16.EncodeRune(r)
					byteStart = i + size
					haveBegin = true
					uIdx += 2
					i += size
					if end == uIdx {
						return StringFromChar(uint16(splitBeginLo))
					}
					continue
				}
			} else if begin == uIdx {
				byteStart = i
				haveBegin = true
			}
		}
		if haveBegin {
			if end == uIdx {
				if splitBeginLo == 0 {
					return unpinSubstring(len(s), s[byteStart:i])
				}
				var sb strings.Builder
				sb.Grow(3 + i - byteStart)
				appendWTF8(&sb, splitBeginLo)
				sb.WriteString(s[byteStart:i])
				return sb.String()
			}
			if r >= 0x10000 && end == uIdx+1 {
				hi, _ := utf16.EncodeRune(r)
				var sb strings.Builder
				extra := 3
				if splitBeginLo != 0 {
					extra = 6
				}
				sb.Grow(extra + i - byteStart)
				if splitBeginLo != 0 {
					appendWTF8(&sb, splitBeginLo)
				}
				sb.WriteString(s[byteStart:i])
				appendWTF8(&sb, hi)
				return sb.String()
			}
		}
		if r >= 0x10000 {
			uIdx += 2
		} else {
			uIdx++
		}
		i += size
	}
	if end == uIdx {
		if !haveBegin {
			return ""
		}
		if splitBeginLo == 0 {
			return unpinSubstring(len(s), s[byteStart:])
		}
		var sb strings.Builder
		sb.Grow(3 + len(s) - byteStart)
		appendWTF8(&sb, splitBeginLo)
		sb.WriteString(s[byteStart:])
		return sb.String()
	}
	panic(NewStringIndexOutOfBoundsExceptionForRange(begin, end, uIdx))
}

// StringSubstringFrom returns the substring from begin to the end of s.
func StringSubstringFrom(s string, begin int) string {
	if m := hotStrMeta0.Load(); m.asciiLen == len(s) && uint(begin) <= uint(len(s)) && m.s == s {
		if begin == len(s) {
			return ""
		}
		return unpinSubstring(len(s), s[begin:])
	}
	if begin >= 0 && begin <= len(s) && len(s) < strCacheMinASCIILen {
		if (len(s) <= 64 && isASCIIUpTo64(s)) || (len(s) > 64 && checkShortASCII(s)) || asciiPrefixLenRaw(s) == len(s) {
			if begin == len(s) {
				return ""
			}
			return s[begin:]
		}
	}
	if m := tryGetStrMeta(s); m != nil {
		return m.substring(begin, m.utf16Len)
	}
	if len(s) >= strCacheMinLen && len(s) < strCacheMinASCIILen {
		if begin < 0 || begin > len(s) {
			panic(NewStringIndexOutOfBoundsExceptionForRange(begin, len(s), len(s)))
		}
		return s[begin:]
	}
	if begin < 0 {
		n := StringLength(s)
		panic(NewStringIndexOutOfBoundsExceptionForRange(begin, n, n))
	}
	prefix := asciiPrefixLenRaw(s)
	if begin <= prefix {
		return unpinSubstring(len(s), s[begin:])
	}
	if prefix == len(s) {
		panic(NewStringIndexOutOfBoundsExceptionForRange(begin, len(s), len(s)))
	}
	uIdx := prefix
	for i := prefix; i < len(s); {
		r, size := decodeWTF8(s[i:])
		if r >= 0x10000 {
			if begin == uIdx {
				return unpinSubstring(len(s), s[i:])
			}
			if begin == uIdx+1 {
				_, lo := utf16.EncodeRune(r)
				var sb strings.Builder
				sb.Grow(3 + len(s) - (i + size))
				appendWTF8(&sb, lo)
				sb.WriteString(s[i+size:])
				return sb.String()
			}
			uIdx += 2
		} else {
			if begin == uIdx {
				return unpinSubstring(len(s), s[i:])
			}
			uIdx++
		}
		i += size
	}
	if begin == uIdx {
		return ""
	}
	panic(NewStringIndexOutOfBoundsExceptionForRange(begin, uIdx, uIdx))
}

func hasSurrogateByte(s string) bool {
	return strings.IndexByte(s, 0xED) >= 0
}

func isCleanUTF8(s string) bool {
	if len(s) < strCacheMinASCIILen {
		prefix := asciiPrefixLenRaw(s)
		if prefix == len(s) {
			return true
		}
		if len(s) >= strCacheMinLen {
			return tryGetStrMeta(s).cleanUTF8
		}
		return utf8.ValidString(s[prefix:])
	}
	if m := tryGetStrMeta(s); m != nil {
		return m.cleanUTF8
	}
	prefix := asciiPrefixLenRaw(s)
	return prefix == len(s) || utf8.ValidString(s[prefix:])
}

func byteIndexToUTF16(s string, bIdx int) int {
	if bIdx == 0 {
		return 0
	}
	if len(s) < strCacheMinASCIILen || bIdx < strCacheMinLen {
		return utf16LenRaw(s[:bIdx])
	}
	if m := tryGetStrMeta(s); m != nil {
		return m.utf16IndexAtByte(bIdx)
	}
	return utf16LenRaw(s[:bIdx])
}

// IsCanonicalWTF8 reports whether s is valid canonical UTF-8 (no lone surrogates or invalid UTF-8).
func IsCanonicalWTF8(s string) bool {
	return isCleanUTF8(s)
}

func utf16Index(us, usub []uint16) int {
	if len(usub) == 0 {
		return 0
	}
	limit := len(us) - len(usub)
	first := usub[0]
	for i := 0; i <= limit; i++ {
		if us[i] == first && slices.Equal(us[i+1:i+len(usub)], usub[1:]) {
			return i
		}
	}
	return -1
}

// StringIndexOf returns the first UTF-16 index of sub in s, or -1.
func StringIndexOf(s, sub string) int {
	if sub == "" {
		return 0
	}
	if asciiPrefixLenRaw(sub) == len(sub) {
		i := strings.Index(s, sub)
		if i < 0 {
			return -1
		}
		return byteIndexToUTF16(s, i)
	}
	if isCleanUTF8(sub) && isCleanUTF8(s) {
		i := strings.Index(s, sub)
		if i < 0 {
			return -1
		}
		return byteIndexToUTF16(s, i)
	}
	return utf16Index(stringUTF16View(s), stringUTF16View(sub))
}

// StringIndexOfFrom returns the first UTF-16 index of sub in s starting at fromIndex, or -1.
func StringIndexOfFrom(s, sub string, fromIndex int) int {
	if fromIndex <= 0 {
		return StringIndexOf(s, sub)
	}
	if sub != "" && asciiPrefixLenRaw(sub) == len(sub) {
		if fromIndex > len(s) {
			return -1
		}
		if asciiPrefixLenRaw(s[:fromIndex]) == fromIndex {
			i := strings.Index(s[fromIndex:], sub)
			if i < 0 {
				return -1
			}
			return byteIndexToUTF16(s, fromIndex+i)
		}
	}
	if IsASCII(s) && IsASCII(sub) {
		if fromIndex > len(s) {
			if sub == "" {
				return len(s)
			}
			return -1
		}
		i := strings.Index(s[fromIndex:], sub)
		if i < 0 {
			return -1
		}
		return fromIndex + i
	}
	if m := tryGetStrMeta(s); m != nil && m.cleanUTF8 && isCleanUTF8(sub) {
		if fromIndex > m.utf16Len {
			if sub == "" {
				return m.utf16Len
			}
			return -1
		}
		if sub == "" {
			return fromIndex
		}
		_, bStart, _, _ := m.resolveUTF16Boundary(fromIndex)
		i := strings.Index(s[bStart:], sub)
		if i < 0 {
			return -1
		}
		return m.utf16IndexAtByte(bStart + i)
	}
	u16s := stringUTF16View(s)
	u16sub := stringUTF16View(sub)
	if fromIndex > len(u16s) {
		if len(u16sub) == 0 {
			return len(u16s)
		}
		return -1
	}
	if len(u16sub) == 0 {
		return fromIndex
	}
	i := utf16Index(u16s[fromIndex:], u16sub)
	if i < 0 {
		return -1
	}
	return fromIndex + i
}

// StringLastIndexOf returns the last UTF-16 index of sub in s, or -1.
func StringLastIndexOf(s, sub string) int {
	if sub == "" {
		return StringLength(s)
	}
	if asciiPrefixLenRaw(sub) == len(sub) {
		i := strings.LastIndex(s, sub)
		if i < 0 {
			return -1
		}
		return byteIndexToUTF16(s, i)
	}
	if isCleanUTF8(sub) && isCleanUTF8(s) {
		i := strings.LastIndex(s, sub)
		if i < 0 {
			return -1
		}
		return byteIndexToUTF16(s, i)
	}
	us, usub := stringUTF16View(s), stringUTF16View(sub)
	for i := len(us) - len(usub); i >= 0; i-- {
		if us[i] == usub[0] && slices.Equal(us[i+1:i+len(usub)], usub[1:]) {
			return i
		}
	}
	return -1
}

// StringLastIndexOfFrom returns the last UTF-16 index of sub in s starting backwards from fromIndex.
func StringLastIndexOfFrom(s, sub string, fromIndex int) int {
	if fromIndex < 0 {
		return -1
	}
	if sub != "" && asciiPrefixLenRaw(sub) == len(sub) {
		if fromIndex >= len(s)-len(sub) {
			return StringLastIndexOf(s, sub)
		}
		if asciiPrefixLenRaw(s[:fromIndex+len(sub)]) == fromIndex+len(sub) {
			return strings.LastIndex(s[:fromIndex+len(sub)], sub)
		}
	}
	if IsASCII(s) && IsASCII(sub) {
		if fromIndex < len(s)-len(sub) {
			s = s[:fromIndex+len(sub)]
		}
		return strings.LastIndex(s, sub)
	}
	if m := tryGetStrMeta(s); m != nil && m.cleanUTF8 && isCleanUTF8(sub) {
		subLen := StringLength(sub)
		k := min(m.utf16Len-subLen, fromIndex)
		if k < 0 {
			return -1
		}
		if sub == "" {
			return k
		}
		bBeforeEnd, _, _, _ := m.resolveUTF16Boundary(k + subLen)
		i := strings.LastIndex(s[:bBeforeEnd], sub)
		if i < 0 {
			return -1
		}
		return m.utf16IndexAtByte(i)
	}
	us, usub := stringUTF16View(s), stringUTF16View(sub)
	k := len(us) - len(usub)
	if k > fromIndex {
		k = fromIndex
	}
	for ; k >= 0; k-- {
		if slices.Equal(us[k:k+len(usub)], usub) {
			return k
		}
	}
	return -1
}

// StringIndexOfChar returns the first UTF-16 index of code unit ch in s.
func StringIndexOfChar(s string, ch uint16) int {
	if ch < 0x80 {
		i := strings.IndexByte(s, byte(ch))
		if i < 0 {
			return -1
		}
		return byteIndexToUTF16(s, i)
	}
	prefix := ASCIIPrefixLen(s)
	if prefix == len(s) {
		return -1
	}
	if (ch < 0xD800 || (ch > 0xDFFF && ch != 0xFFFD)) && isCleanUTF8(s) {
		i := strings.IndexRune(s[prefix:], rune(ch))
		if i < 0 {
			return -1
		}
		return byteIndexToUTF16(s, prefix+i)
	}
	uIdx := prefix
	for i := prefix; i < len(s); {
		if s[i] < 0x80 {
			i++
			uIdx++
			continue
		}
		r, sz := decodeWTF8(s[i:])
		i += sz
		if r >= 0x10000 {
			hi, lo := utf16.EncodeRune(r)
			if uint16(hi) == ch {
				return uIdx
			}
			if uint16(lo) == ch {
				return uIdx + 1
			}
			uIdx += 2
		} else {
			if uint16(r) == ch {
				return uIdx
			}
			uIdx++
		}
	}
	return -1
}

// StringLastIndexOfCharFrom returns the last UTF-16 index of code unit ch in s at or before fromIndex.
func StringLastIndexOfCharFrom(s string, ch uint16, fromIndex int) int {
	if fromIndex < 0 {
		return -1
	}
	if ch < 0x80 {
		if fromIndex >= len(s) {
			b := strings.LastIndexByte(s, byte(ch))
			if b < 0 {
				return -1
			}
			return byteIndexToUTF16(s, b)
		}
		if len(s) < strCacheMinASCIILen && asciiPrefixLenRaw(s[:fromIndex+1]) == fromIndex+1 {
			return strings.LastIndexByte(s[:fromIndex+1], byte(ch))
		}
		if IsASCII(s) {
			if fromIndex >= len(s) {
				fromIndex = len(s) - 1
			}
			if fromIndex < 0 {
				return -1
			}
			return strings.LastIndexByte(s[:fromIndex+1], byte(ch))
		}
		if m := tryGetStrMeta(s); m != nil {
			idx := min(fromIndex, m.utf16Len-1)
			if idx < 0 {
				return -1
			}
			bBeforeEnd, _, _, _ := m.resolveUTF16Boundary(idx + 1)
			b := strings.LastIndexByte(s[:bBeforeEnd], byte(ch))
			if b < 0 {
				return -1
			}
			return m.utf16IndexAtByte(b)
		}
	}
	prefix := ASCIIPrefixLen(s)
	if prefix == len(s) {
		return -1
	}
	lastMatch := -1
	uIdx := prefix
	for i := prefix; i < len(s) && uIdx <= fromIndex; {
		if s[i] < 0x80 {
			if uint16(s[i]) == ch {
				lastMatch = uIdx
			}
			i++
			uIdx++
			continue
		}
		r, sz := decodeWTF8(s[i:])
		i += sz
		if r >= 0x10000 {
			hi, lo := utf16.EncodeRune(r)
			if uint16(hi) == ch {
				lastMatch = uIdx
			}
			if uIdx+1 <= fromIndex && uint16(lo) == ch {
				lastMatch = uIdx + 1
			}
			uIdx += 2
		} else {
			if uint16(r) == ch {
				lastMatch = uIdx
			}
			uIdx++
		}
	}
	if ch < 0x80 && lastMatch < 0 {
		endASCII := min(prefix-1, fromIndex)
		if endASCII >= 0 {
			return strings.LastIndexByte(s[:endASCII+1], byte(ch))
		}
	}
	return lastMatch
}

// StringCompareTo compares two strings in UTF-16 code unit order (Java String.compareTo).
func StringCompareTo(a, b string) int {
	if a == b {
		return 0
	}
	n := min(len(a), len(b))
	i := 0
	for i+8 <= n && a[i:i+8] == b[i:i+8] {
		i += 8
	}
	for i < n && a[i] == b[i] {
		i++
	}
	switch {
	case i < n:
		if a[i] < 0x80 && b[i] < 0x80 {
			return int(a[i]) - int(b[i])
		}
	case len(a) < len(b):
		if b[n]&0xC0 != 0x80 {
			return -utf16LenRaw(b[n:])
		}
	default:
		if a[n]&0xC0 != 0x80 {
			return utf16LenRaw(a[n:])
		}
	}
	return stringCompareUTF16From(a, b, i)
}

func stringCompareUTF16From(a, b string, diff int) int {
	start := diff
	if start > 0 && a[start-1] >= 0x80 {
		start--
		for start > 0 && a[start]&0xC0 == 0x80 {
			start--
		}
	}
	ia, ib := start, start
	pendA, pendB := int32(-1), int32(-1)
	for {
		var ua uint16
		hasA := true
		if pendA >= 0 {
			ua = uint16(pendA)
			pendA = -1
		} else if ia < len(a) {
			r, sz := decodeWTF8(a[ia:])
			ia += sz
			if r >= 0x10000 {
				hi, lo := utf16.EncodeRune(r)
				ua = uint16(hi)
				pendA = int32(uint16(lo))
			} else {
				ua = uint16(r)
			}
		} else {
			hasA = false
		}

		var ub uint16
		hasB := true
		if pendB >= 0 {
			ub = uint16(pendB)
			pendB = -1
		} else if ib < len(b) {
			r, sz := decodeWTF8(b[ib:])
			ib += sz
			if r >= 0x10000 {
				hi, lo := utf16.EncodeRune(r)
				ub = uint16(hi)
				pendB = int32(uint16(lo))
			} else {
				ub = uint16(r)
			}
		} else {
			hasB = false
		}

		if !hasA {
			if !hasB {
				return 0
			}
			rem := 1 + utf16LenRaw(b[ib:])
			if pendB >= 0 {
				rem++
			}
			return -rem
		}
		if !hasB {
			rem := 1 + utf16LenRaw(a[ia:])
			if pendA >= 0 {
				rem++
			}
			return rem
		}
		if ua != ub {
			return int(ua) - int(ub)
		}
	}
}

// StringHashCode returns the 32-bit polynomial hash code of s (Java String.hashCode()).
func StringHashCode(s string) int32 {
	var h int32
	i := 0
	for ; i+4 <= len(s); i += 4 {
		b0, b1, b2, b3 := s[i], s[i+1], s[i+2], s[i+3]
		if (b0 | b1 | b2 | b3) >= 0x80 {
			break
		}
		h = 923521*h + 29791*int32(b0) + 961*int32(b1) + 31*int32(b2) + int32(b3)
	}
	for ; i < len(s); i++ {
		c := s[i]
		if c >= 0x80 {
			return stringHashCodeNonASCII(s, i, h)
		}
		h = 31*h + int32(c)
	}
	return h
}

//go:noinline
func stringHashCodeNonASCII(s string, i int, h int32) int32 {
	for i < len(s) {
		r, size := decodeWTF8(s[i:])
		i += size
		if r >= 0x10000 {
			hi, lo := utf16.EncodeRune(r)
			h = 31*(31*h+int32(hi)) + int32(lo)
		} else {
			h = 31*h + int32(r)
		}
	}
	return h
}

// rebindCachedStrMeta updates the cached strMeta's string header to s when a
// SourceFile is freshly loaded (checked once per file via BOM StringStartsWith),
// ensuring subsequent O(N) StringCharAt/StringLength scans match m.s == s by
// pointer equality in 2 instructions with zero per-character atomics or allocs.
func rebindCachedStrMeta(s string) {
	if len(s) < strCacheMinASCIILen || len(s) > strCacheMaxLen {
		return
	}
	if m0 := hotStrMeta0.Load(); m0.s == s {
		cp := *m0
		cp.s = s
		hotStrMeta0.Store(&cp)
		sh := &strMetaCache[cp.hash&(strCacheShards-1)]
		for i := 0; i < strCacheWays; i++ {
			if sh.slots[i].CompareAndSwap(m0, &cp) {
				break
			}
		}
		return
	}
	h := hashStrMetaKey(s)
	sh := &strMetaCache[h&(strCacheShards-1)]
	for i := 0; i < strCacheWays; i++ {
		if m := sh.slots[i].Load(); m != nil && m.hash == h && m.s == s {
			cp := *m
			cp.s = s
			sh.slots[i].CompareAndSwap(m, &cp)
			promoteHotStrMeta(&cp)
			return
		}
	}
}

// StringStartsWith reports whether s begins with prefix in UTF-16 code units.
func StringStartsWith(s, prefix string) bool {
	if prefix == "" {
		return true
	}
	if len(prefix) > len(s) {
		return false
	}
	if len(prefix) == 3 && prefix == "\ufeff" && len(s) >= strCacheMinASCIILen {
		rebindCachedStrMeta(s)
	}
	if (s[0] < 0x80 || prefix[0] < 0x80) && s[0] != prefix[0] {
		return false
	}
	if asciiPrefixLenRaw(prefix) == len(prefix) || (isCleanUTF8(prefix) && isCleanUTF8(s)) {
		return s[:len(prefix)] == prefix
	}
	us, up := stringUTF16View(s), stringUTF16View(prefix)
	return len(us) >= len(up) && slices.Equal(us[:len(up)], up)
}

// StringEndsWith reports whether s ends with suffix in UTF-16 code units.
func StringEndsWith(s, suffix string) bool {
	if suffix == "" {
		return true
	}
	if len(suffix) > len(s) {
		return false
	}
	if (s[len(s)-1]|suffix[len(suffix)-1]) < 0x80 && s[len(s)-1] != suffix[len(suffix)-1] {
		return false
	}
	if asciiPrefixLenRaw(suffix) == len(suffix) || (isCleanUTF8(suffix) && isCleanUTF8(s)) {
		return s[len(s)-len(suffix):] == suffix
	}
	us, usuf := stringUTF16View(s), stringUTF16View(suffix)
	return len(us) >= len(usuf) && slices.Equal(us[len(us)-len(usuf):], usuf)
}

// StringRegionMatches implements Java String.regionMatches(toffset, other, ooffset, len).
func StringRegionMatches(s string, toffset int, other string, ooffset int, length int) bool {
	if toffset < 0 || ooffset < 0 {
		return false
	}
	if length > 0 && int64(toffset)+int64(length) <= int64(len(s)) && int64(ooffset)+int64(length) <= int64(len(other)) {
		if asciiPrefixLenRaw(s[:toffset+length]) == toffset+length && asciiPrefixLenRaw(other[:ooffset+length]) == ooffset+length {
			return s[toffset:toffset+length] == other[ooffset:ooffset+length]
		}
	}
	if IsASCII(s) && IsASCII(other) {
		if int64(toffset) > int64(len(s))-int64(length) || int64(ooffset) > int64(len(other))-int64(length) {
			return false
		}
		if length <= 0 {
			return true
		}
		return s[toffset:toffset+length] == other[ooffset:ooffset+length]
	}
	sLen := StringLength(s)
	oLen := StringLength(other)
	if int64(toffset) > int64(sLen)-int64(length) || int64(ooffset) > int64(oLen)-int64(length) {
		return false
	}
	if length <= 0 {
		return true
	}
	for i := 0; i < length; i++ {
		if StringCharAt(s, toffset+i) != StringCharAt(other, ooffset+i) {
			return false
		}
	}
	return true
}

// StringReplace replaces all occurrences of target with replacement in UTF-16 code units.
func StringReplace(s, target, replacement string) string {
	if !hasSurrogateByte(s) && !hasSurrogateByte(target) && !hasSurrogateByte(replacement) && (target != "" || IsASCII(s)) {
		return strings.ReplaceAll(s, target, replacement)
	}
	us, ut := stringUTF16View(s), stringUTF16View(target)
	if len(ut) == 0 {
		ur := stringUTF16View(replacement)
		out := make([]uint16, 0, len(us)+len(ur)*(len(us)+1))
		out = append(out, ur...)
		for _, c := range us {
			out = append(out, c)
			out = append(out, ur...)
		}
		return StringFromUTF16(out)
	}
	idx := utf16Index(us, ut)
	if idx < 0 {
		if !hasSurrogateByte(s) {
			return s
		}
		return StringFromUTF16(us)
	}
	ur := stringUTF16View(replacement)
	out := make([]uint16, 0, len(us)+len(ur))
	pos := 0
	for {
		out = append(out, us[pos:pos+idx]...)
		out = append(out, ur...)
		pos += idx + len(ut)
		if pos > len(us)-len(ut) {
			break
		}
		idx = utf16Index(us[pos:], ut)
		if idx < 0 {
			break
		}
	}
	out = append(out, us[pos:]...)
	return StringFromUTF16(out)
}

// StringTrim trims leading and trailing whitespace characters (<= 0x20) per Java String.trim().
func StringTrim(s string) string {
	start := 0
	for start < len(s) && s[start] <= ' ' {
		start++
	}
	end := len(s)
	for end > start && s[end-1] <= ' ' {
		end--
	}
	return unpinSubstring(len(s), s[start:end])
}

// StringSplitLiteral splits s around literal sep (limit == 0: trailing empty strings discarded).
func StringSplitLiteral(s, sep string) []string {
	if sep == "" {
		panic(NewIllegalArgumentException("StringSplitLiteral: empty separator"))
	}
	if !strings.Contains(s, sep) {
		return []string{s}
	}
	parts := strings.Split(s, sep)
	n := len(parts)
	for n > 0 && parts[n-1] == "" {
		n--
	}
	return parts[:n]
}

// StringSplitLiteralLimit splits s around literal sep with Java String.split(regex, limit) semantics.
func StringSplitLiteralLimit(s, sep string, limit int) []string {
	if limit == 0 {
		return StringSplitLiteral(s, sep)
	}
	if sep == "" {
		panic(NewIllegalArgumentException("StringSplitLiteralLimit: empty separator"))
	}
	if limit < 0 {
		return strings.Split(s, sep)
	}
	return strings.SplitN(s, sep, limit)
}

// StringToLowerCaseASCII converts ASCII uppercase characters to lowercase.
func StringToLowerCaseASCII(s string) string {
	if !hasUpper(s) {
		return s
	}
	b := []byte(s)
	for i := range b {
		if b[i] >= 'A' && b[i] <= 'Z' {
			b[i] += 'a' - 'A'
		}
	}
	return string(b)
}

// StringToUpperCaseASCII converts ASCII lowercase characters to uppercase.
func StringToUpperCaseASCII(s string) string {
	if !hasLower(s) {
		return s
	}
	b := []byte(s)
	for i := range b {
		if b[i] >= 'a' && b[i] <= 'z' {
			b[i] -= 'a' - 'A'
		}
	}
	return string(b)
}

func hasUpper(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= 'A' && s[i] <= 'Z' {
			return true
		}
	}
	return false
}

func hasLower(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= 'a' && s[i] <= 'z' {
			return true
		}
	}
	return false
}
