package jutil

import "unicode"

// Go-only: language-adapter — java.lang.Character case mapping and
// isSupplementaryCodePoint, and String.regionMatches on UTF-16 code units.
// Java's Character.toUpperCase(char)/toLowerCase(char) are the simple
// UnicodeData mappings; Go's unicode package uses the same data (Unicode
// 15.0 in both JDK 21 and Go 1.24), checked on all 65536 chars by
// TestCharacterCaseMatchesJava.

// CharacterToUpperCase is Character.toUpperCase(char).
func CharacterToUpperCase(ch uint16) uint16 {
	if ch < 0x80 {
		if ch >= 'a' && ch <= 'z' {
			return ch - 'a' + 'A'
		}
		return ch
	}
	if isSurrogate(ch) {
		return ch
	}
	u := unicode.ToUpper(rune(ch))
	if u > 0xFFFF {
		return ch
	}
	return uint16(u)
}

// CharacterToLowerCase is Character.toLowerCase(char).
func CharacterToLowerCase(ch uint16) uint16 {
	if ch < 0x80 {
		if ch >= 'A' && ch <= 'Z' {
			return ch - 'A' + 'a'
		}
		return ch
	}
	if isSurrogate(ch) {
		return ch
	}
	l := unicode.ToLower(rune(ch))
	if l > 0xFFFF {
		return ch
	}
	return uint16(l)
}

// CharacterIsSupplementaryCodePoint is
// Character.isSupplementaryCodePoint(int).
func CharacterIsSupplementaryCodePoint(codePoint int) bool {
	return codePoint >= 0x10000 && codePoint < 0x110000
}

// UTF16RegionMatches is String.regionMatches(boolean ignoreCase, int
// toffset, String other, int ooffset, int len) on UTF-16 code units (see its
// Javadoc). The result is false when an offset is negative or a region runs
// past the end of its string; otherwise the len units starting at each
// offset are compared. With ignoreCase, two code points are equal when they
// are the same or when Character.toLowerCase(Character.toUpperCase(cp)) is
// the same for both.
//
// Observed behaviour (black-box): the case-insensitive comparison works on
// code points, so a surrogate pair inside both regions is mapped as one
// supplementary character (U+10400 matches U+10428), while a half of a pair
// cut off by a region boundary is compared as a lone surrogate. No case
// mapping crosses between the BMP and the supplementary planes (checked on
// every code point), so walking both regions in step is exact. A negative
// len with in-range offsets compares nothing and gives true.
func UTF16RegionMatches(ignoreCase bool, value []uint16, toffset int, other []uint16, ooffset int, length int) bool {
	if toffset < 0 || ooffset < 0 || toffset > len(value)-length || ooffset > len(other)-length {
		return false
	}
	for i := 0; i < length; i++ {
		a, b := value[toffset+i], other[ooffset+i]
		if a == b && !(ignoreCase && isHighSurrogateUnit(a)) {
			continue
		}
		if !ignoreCase {
			return false
		}
		cpA, cpB := rune(a), rune(b)
		if i+1 < length && isHighSurrogateUnit(a) && isHighSurrogateUnit(b) {
			la, lb := value[toffset+i+1], other[ooffset+i+1]
			if isLowSurrogateUnit(la) && isLowSurrogateUnit(lb) {
				cpA = (cpA-0xD800)<<10 + rune(la) - 0xDC00 + 0x10000
				cpB = (cpB-0xD800)<<10 + rune(lb) - 0xDC00 + 0x10000
				i++
			}
		}
		if cpA != cpB && caseKey(cpA) != caseKey(cpB) {
			return false
		}
	}
	return true
}

// caseKey is Character.toLowerCase(Character.toUpperCase(cp)), the
// equivalence used by regionMatches(true, ...). Surrogates map to
// themselves.
func caseKey(cp rune) rune {
	if cp < 0x10000 {
		return rune(CharacterToLowerCase(CharacterToUpperCase(uint16(cp))))
	}
	return unicode.ToLower(unicode.ToUpper(cp))
}

func isHighSurrogateUnit(c uint16) bool { return c >= 0xD800 && c <= 0xDBFF }

func isLowSurrogateUnit(c uint16) bool { return c >= 0xDC00 && c <= 0xDFFF }
