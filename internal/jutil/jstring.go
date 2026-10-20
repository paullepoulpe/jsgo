package jutil

import (
	"slices"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
	"unsafe"
)

// Go-only: language-adapter — Java String semantics (UTF-16 code units) on Go
// strings that hold WTF-8.

// StringFromUTF16 converts UTF-16 code units (Java chars) to a WTF-8 Go
// string. Valid surrogate pairs become one 4-byte sequence; lone surrogates
// are encoded as their 3-byte generalized UTF-8 form so that no information
// is lost (Java strings, and therefore Closure's AST, can hold them).
func StringFromUTF16(units []uint16) string {
	ascii := true
	for _, u := range units {
		if u >= 0x80 {
			ascii = false
			break
		}
	}
	if ascii {
		if len(units) == 0 {
			return ""
		}
		b := make([]byte, len(units))
		for i, u := range units {
			b[i] = byte(u)
		}
		// b is never written again, so it can back the string directly
		// (one allocation instead of two).
		return unsafe.String(unsafe.SliceData(b), len(b))
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

// StringFromChar returns the one-char Java string for a UTF-16 code unit.
func StringFromChar(c uint16) string {
	if c < 0x80 {
		return string(rune(c))
	}
	var sb strings.Builder
	appendWTF8(&sb, rune(c))
	return sb.String()
}

// decodeWTF8 decodes one code point (possibly a lone surrogate) from s.
func decodeWTF8(s string) (r rune, size int) {
	c := s[0]
	if c < 0x80 {
		return rune(c), 1
	}
	// Generalized UTF-8 surrogate: ED A0..BF xx
	if c == 0xED && len(s) >= 3 && s[1] >= 0xA0 && s[1] <= 0xBF && s[2]&0xC0 == 0x80 {
		return rune(c&0x0F)<<12 | rune(s[1]&0x3F)<<6 | rune(s[2]&0x3F), 3
	}
	return utf8.DecodeRuneInString(s)
}

// IsASCII reports whether s contains only ASCII bytes (< 0x80).
func IsASCII(s string) bool {
	return isASCII(s)
}

// ASCIIPrefixLen returns the byte length of the leading pure-ASCII prefix of s
// (len(s) if s is entirely ASCII).
func ASCIIPrefixLen(s string) int {
	i := 0
	if len(s) >= 8 {
		p := unsafe.Pointer(unsafe.StringData(s))
		for ; i+8 <= len(s); i += 8 {
			if *(*uint64)(unsafe.Add(p, i))&0x8080808080808080 != 0 {
				break
			}
		}
	}
	for ; i < len(s); i++ {
		if s[i] >= 0x80 {
			return i
		}
	}
	return len(s)
}

func isASCII(s string) bool {
	// Go-only: performance — eight bytes at a time (unaligned loads are
	// fine on the platforms Go supports); this check runs before most
	// String operations.
	return ASCIIPrefixLen(s) == len(s)
}

// StringToUTF16 converts a WTF-8 Go string to UTF-16 code units.
func StringToUTF16(s string) []uint16 {
	out, _ := StringToUTF16AndASCIIPrefix(s)
	return out
}

// StringToUTF16AndASCIIPrefix converts s to UTF-16 code units and returns
// the length of its leading pure-ASCII prefix.
func StringToUTF16AndASCIIPrefix(s string) ([]uint16, int) {
	prefix := ASCIIPrefixLen(s)
	if prefix == len(s) {
		out := make([]uint16, len(s))
		for i := 0; i < len(s); i++ {
			out[i] = uint16(s[i])
		}
		return out, prefix
	}
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
	return out, prefix
}

// StringLength is Java's String.length(): the number of UTF-16 code units.
func StringLength(s string) int {
	if isASCII(s) {
		return len(s)
	}
	n := 0
	for i := 0; i < len(s); {
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

// StringCharAt is Java's String.charAt(int).
func StringCharAt(s string, index int) uint16 {
	if isASCII(s) {
		if index < 0 || index >= len(s) {
			panic(NewStringIndexOutOfBoundsExceptionForIndex(index, len(s)))
		}
		return uint16(s[index])
	}
	u := StringToUTF16(s)
	if index < 0 || index >= len(u) {
		panic(NewStringIndexOutOfBoundsExceptionForIndex(index, len(u)))
	}
	return u[index]
}

// StringSubstring is Java's String.substring(begin, end) with UTF-16 indices.
func StringSubstring(s string, begin, end int) string {
	if isASCII(s) {
		if begin < 0 || end > len(s) || begin > end {
			panic(NewStringIndexOutOfBoundsExceptionForRange(begin, end, len(s)))
		}
		return s[begin:end]
	}
	u := StringToUTF16(s)
	if begin < 0 || end > len(u) || begin > end {
		panic(NewStringIndexOutOfBoundsExceptionForRange(begin, end, len(u)))
	}
	return StringFromUTF16(u[begin:end])
}

// StringSubstringFrom is Java's String.substring(begin).
func StringSubstringFrom(s string, begin int) string {
	return StringSubstring(s, begin, StringLength(s))
}

// StringIndexOf is Java's String.indexOf(String): a UTF-16 index or -1.
func StringIndexOf(s, sub string) int {
	i := strings.Index(s, sub)
	if i < 0 {
		return -1
	}
	return StringLength(s[:i])
}

// StringIndexOfFrom is Java's String.indexOf(String, int fromIndex).
func StringIndexOfFrom(s, sub string, from int) int {
	if from < 0 {
		from = 0
	}
	if isASCII(s) && isASCII(sub) {
		if from > len(s) {
			if sub == "" {
				return len(s)
			}
			return -1
		}
		i := strings.Index(s[from:], sub)
		if i < 0 {
			return -1
		}
		return from + i
	}
	n := StringLength(s)
	if from > n {
		if sub == "" {
			return n
		}
		return -1
	}
	rest := StringSubstringFrom(s, from)
	i := StringIndexOf(rest, sub)
	if i < 0 {
		return -1
	}
	return from + i
}

// StringLastIndexOf is Java's String.lastIndexOf(String).
func StringLastIndexOf(s, sub string) int {
	i := strings.LastIndex(s, sub)
	if i < 0 {
		return -1
	}
	return StringLength(s[:i])
}

// StringCompareTo is Java's String.compareTo: lexicographic over UTF-16 code
// units (which differs from Go's byte order for supplementary characters).
func StringCompareTo(a, b string) int {
	// Go-only: performance — decide on the bytes when the strings first
	// differ at an ASCII byte in both, or when one is a prefix of the other,
	// instead of checking both strings for non-ASCII first. Either way the
	// common bytes decode to the same UTF-16 units: a decode that started in
	// the common prefix cannot reach past it, because the next byte is ASCII
	// in both strings, or the end of one of them and not a continuation byte
	// in the other.
	n := min(len(a), len(b))
	i := 0
	for i < n && a[i] == b[i] {
		i++
	}
	switch {
	case i < n:
		if a[i] < 0x80 && b[i] < 0x80 {
			return int(a[i]) - int(b[i])
		}
	case len(a) == len(b):
		return 0
	case len(a) < len(b):
		if b[n]&0xC0 != 0x80 {
			return -StringLength(b[n:])
		}
	default:
		if a[n]&0xC0 != 0x80 {
			return StringLength(a[n:])
		}
	}
	ua, ub := StringToUTF16(a), StringToUTF16(b)
	m := min(len(ua), len(ub))
	for k := 0; k < m; k++ {
		if ua[k] != ub[k] {
			return int(ua[k]) - int(ub[k])
		}
	}
	return len(ua) - len(ub)
}

// StringHashCode is Java's String.hashCode().
func StringHashCode(s string) int32 {
	// Go-only: performance — one pass, without converting the string: the
	// ASCII prefix hashes byte by byte, and the rest code point by code point
	// (an ASCII byte ends any code point before it, so decoding from there
	// gives the same UTF-16 units as converting the whole string).
	var h int32
	i := 0
	for ; i < len(s) && s[i] < 0x80; i++ {
		h = 31*h + int32(s[i])
	}
	for i < len(s) {
		r, size := decodeWTF8(s[i:])
		i += size
		if r >= 0x10000 {
			hi, lo := utf16.EncodeRune(r)
			h = 31*h + int32(uint16(hi))
			h = 31*h + int32(uint16(lo))
		} else {
			h = 31*h + int32(uint16(r))
		}
	}
	return h
}

// StringToLowerCaseASCII is Guava's Ascii.toLowerCase(String).
func StringToLowerCaseASCII(s string) string {
	for i := 0; i < len(s); i++ {
		if c := s[i]; c >= 'A' && c <= 'Z' {
			b := []byte(s)
			for ; i < len(b); i++ {
				if c := b[i]; c >= 'A' && c <= 'Z' {
					b[i] = c + ('a' - 'A')
				}
			}
			return unsafe.String(unsafe.SliceData(b), len(b))
		}
	}
	return s
}

// StringToUpperCaseASCII is Guava's Ascii.toUpperCase(String).
func StringToUpperCaseASCII(s string) string {
	for i := 0; i < len(s); i++ {
		if c := s[i]; c >= 'a' && c <= 'z' {
			b := []byte(s)
			for ; i < len(b); i++ {
				if c := b[i]; c >= 'a' && c <= 'z' {
					b[i] = c - ('a' - 'A')
				}
			}
			return unsafe.String(unsafe.SliceData(b), len(b))
		}
	}
	return s
}

// CaseFormat_UpperUnderscoreToLowerCamel is Guava's
// CaseFormat.UPPER_UNDERSCORE.to(CaseFormat.LOWER_CAMEL, s) for ASCII input:
// "INLINE_TYPE" -> "inlineType".
func CaseFormat_UpperUnderscoreToLowerCamel(s string) string {
	var sb strings.Builder
	words := strings.Split(s, "_")
	for i, w := range words {
		if i == 0 {
			sb.WriteString(StringToLowerCaseASCII(w))
			continue
		}
		if w == "" {
			continue
		}
		sb.WriteString(StringToUpperCaseASCII(w[:1]))
		sb.WriteString(StringToLowerCaseASCII(w[1:]))
	}
	return sb.String()
}

// StringReplace is String.replace(CharSequence target, CharSequence
// replacement). Both arguments are whole code-unit sequences, so replacing
// on the WTF-8 bytes gives the same result as Java's UTF-16 replacement.
func StringReplace(s, target, replacement string) string {
	return strings.ReplaceAll(s, target, replacement)
}

// StringStartsWith is String.startsWith(String prefix).
func StringStartsWith(s, prefix string) bool { return strings.HasPrefix(s, prefix) }

// StringEndsWith is String.endsWith(String suffix).
func StringEndsWith(s, suffix string) bool { return strings.HasSuffix(s, suffix) }

// StringRegionMatches is String.regionMatches(int toffset, String other, int
// ooffset, int len) with UTF-16 offsets.
func StringRegionMatches(s string, toffset int, other string, ooffset int, length int) bool {
	if toffset < 0 || ooffset < 0 {
		return false
	}
	if isASCII(s) && isASCII(other) {
		if int64(toffset) > int64(len(s))-int64(length) || int64(ooffset) > int64(len(other))-int64(length) {
			return false
		}
		if length <= 0 {
			return true
		}
		return s[toffset:toffset+length] == other[ooffset:ooffset+length]
	}
	us, uo := StringToUTF16(s), StringToUTF16(other)
	if int64(toffset) > int64(len(us))-int64(length) || int64(ooffset) > int64(len(uo))-int64(length) {
		return false
	}
	for i := 0; i < length; i++ {
		if us[toffset+i] != uo[ooffset+i] {
			return false
		}
	}
	return true
}

// StringLastIndexOfCharFrom is String.lastIndexOf(int ch, int fromIndex)
// for a BMP character, with UTF-16 indices.
func StringLastIndexOfCharFrom(s string, ch uint16, fromIndex int) int {
	// Go-only: performance — on an ASCII string, UTF-16 indices are byte
	// offsets.
	if ch < 0x80 && isASCII(s) {
		if fromIndex >= len(s) {
			fromIndex = len(s) - 1
		}
		if fromIndex < 0 {
			return -1
		}
		return strings.LastIndexByte(s[:fromIndex+1], byte(ch))
	}
	u := StringToUTF16(s)
	if fromIndex >= len(u) {
		fromIndex = len(u) - 1
	}
	for i := fromIndex; i >= 0; i-- {
		if u[i] == ch {
			return i
		}
	}
	return -1
}

// StringIndexOfChar is String.indexOf(int ch) for a BMP character.
func StringIndexOfChar(s string, ch uint16) int {
	if ch < 0x80 {
		i := strings.IndexByte(s, byte(ch))
		if i < 0 {
			return -1
		}
		return StringLength(s[:i])
	}
	for i, c := range StringToUTF16(s) {
		if c == ch {
			return i
		}
	}
	return -1
}

// StringIndexOfCharFrom is String.indexOf(int ch, int fromIndex) for a BMP
// character.
func StringIndexOfCharFrom(s string, ch uint16, fromIndex int) int {
	// Go-only: performance — for an ASCII char, search the bytes instead of
	// converting all of s to UTF-16 on every call (callers loop over the
	// occurrences, which made them quadratic). An ASCII byte is always a
	// code point of its own, so its UTF-16 index is that of the code point
	// where the search starts plus the length of what lies in between.
	if ch < 0x80 {
		if fromIndex < 0 {
			fromIndex = 0
		}
		// Skip the code points before fromIndex. A code point straddling
		// fromIndex (a surrogate pair) is skipped too: its low surrogate is
		// not ch.
		i, u := 0, 0
		for i < len(s) && u < fromIndex {
			if s[i] < 0x80 {
				i++
				u++
				continue
			}
			r, size := decodeWTF8(s[i:])
			i += size
			if r >= 0x10000 {
				u += 2
			} else {
				u++
			}
		}
		j := strings.IndexByte(s[i:], byte(ch))
		if j < 0 {
			return -1
		}
		return u + StringLength(s[i:i+j])
	}
	u := StringToUTF16(s)
	if fromIndex < 0 {
		fromIndex = 0
	}
	for i := fromIndex; i < len(u); i++ {
		if u[i] == ch {
			return i
		}
	}
	return -1
}

// StringCodePointAt is String.codePointAt(int) with a UTF-16 index.
func StringCodePointAt(s string, index int) int {
	if isASCII(s) {
		if index < 0 || index >= len(s) {
			panic(NewStringIndexOutOfBoundsException(index, len(s)))
		}
		return int(s[index])
	}
	u := StringToUTF16(s)
	if index < 0 || index >= len(u) {
		panic(NewStringIndexOutOfBoundsException(index, len(u)))
	}
	c := u[index]
	if c >= 0xD800 && c <= 0xDBFF && index+1 < len(u) {
		if lo := u[index+1]; lo >= 0xDC00 && lo <= 0xDFFF {
			return int(utf16.DecodeRune(rune(c), rune(lo)))
		}
	}
	return int(c)
}

// StringLastIndexOfFrom is String.lastIndexOf(String str, int fromIndex)
// with UTF-16 indices (see its Javadoc): the greatest k <= fromIndex at
// which sub occurs in s, or -1. A negative fromIndex therefore always gives
// -1, a fromIndex past the end is no limit at all, and the empty string
// occurs at every index from 0 to StringLength(s).
func StringLastIndexOfFrom(s, sub string, fromIndex int) int {
	if fromIndex < 0 {
		return -1
	}
	if isASCII(s) && isASCII(sub) {
		// Byte offsets are UTF-16 indices; an occurrence at k <= fromIndex
		// ends at or before fromIndex+len(sub).
		if fromIndex < len(s)-len(sub) {
			s = s[:fromIndex+len(sub)]
		}
		return strings.LastIndex(s, sub)
	}
	// Compare code units, not WTF-8 bytes: in Java a lone surrogate in sub
	// can match one half of a surrogate pair of s.
	us, usub := StringToUTF16(s), StringToUTF16(sub)
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
