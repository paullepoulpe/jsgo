package jutil

import (
	"math"
	"unicode"
)

// Go-only: language-adapter — java.lang.Character classification of UTF-16
// code units (Java `char`). Java 21 uses Unicode 15.0, as does Go 1.24's
// unicode package; testdata/character.tsv (generated with JDK 21) checks
// every code unit.

// Java Character.getType() general category constants.
const (
	character_UNASSIGNED                = 0
	character_UPPERCASE_LETTER          = 1
	character_LOWERCASE_LETTER          = 2
	character_TITLECASE_LETTER          = 3
	character_MODIFIER_LETTER           = 4
	character_OTHER_LETTER              = 5
	Character_NON_SPACING_MARK          = 6
	character_ENCLOSING_MARK            = 7
	character_COMBINING_SPACING_MARK    = 8
	character_DECIMAL_DIGIT_NUMBER      = 9
	character_LETTER_NUMBER             = 10
	character_OTHER_NUMBER              = 11
	Character_SPACE_SEPARATOR           = 12
	character_LINE_SEPARATOR            = 13
	character_PARAGRAPH_SEPARATOR       = 14
	character_CONTROL                   = 15
	Character_FORMAT                    = 16
	character_PRIVATE_USE               = 18
	character_SURROGATE                 = 19
	character_DASH_PUNCTUATION          = 20
	character_START_PUNCTUATION         = 21
	character_END_PUNCTUATION           = 22
	character_CONNECTOR_PUNCTUATION     = 23
	character_OTHER_PUNCTUATION         = 24
	character_MATH_SYMBOL               = 25
	character_CURRENCY_SYMBOL           = 26
	character_MODIFIER_SYMBOL           = 27
	character_OTHER_SYMBOL              = 28
	character_INITIAL_QUOTE_PUNCTUATION = 29
	character_FINAL_QUOTE_PUNCTUATION   = 30
)

var characterCategories = []struct {
	table *unicode.RangeTable
	typ   int
}{
	{unicode.Lu, character_UPPERCASE_LETTER},
	{unicode.Ll, character_LOWERCASE_LETTER},
	{unicode.Lt, character_TITLECASE_LETTER},
	{unicode.Lm, character_MODIFIER_LETTER},
	{unicode.Lo, character_OTHER_LETTER},
	{unicode.Mn, Character_NON_SPACING_MARK},
	{unicode.Me, character_ENCLOSING_MARK},
	{unicode.Mc, character_COMBINING_SPACING_MARK},
	{unicode.Nd, character_DECIMAL_DIGIT_NUMBER},
	{unicode.Nl, character_LETTER_NUMBER},
	{unicode.No, character_OTHER_NUMBER},
	{unicode.Zs, Character_SPACE_SEPARATOR},
	{unicode.Zl, character_LINE_SEPARATOR},
	{unicode.Zp, character_PARAGRAPH_SEPARATOR},
	{unicode.Cc, character_CONTROL},
	{unicode.Cf, Character_FORMAT},
	{unicode.Co, character_PRIVATE_USE},
	{unicode.Cs, character_SURROGATE},
	{unicode.Pd, character_DASH_PUNCTUATION},
	{unicode.Ps, character_START_PUNCTUATION},
	{unicode.Pe, character_END_PUNCTUATION},
	{unicode.Pc, character_CONNECTOR_PUNCTUATION},
	{unicode.Po, character_OTHER_PUNCTUATION},
	{unicode.Sm, character_MATH_SYMBOL},
	{unicode.Sc, character_CURRENCY_SYMBOL},
	{unicode.Sk, character_MODIFIER_SYMBOL},
	{unicode.So, character_OTHER_SYMBOL},
	{unicode.Pi, character_INITIAL_QUOTE_PUNCTUATION},
	{unicode.Pf, character_FINAL_QUOTE_PUNCTUATION},
}

// CharacterGetType is Character.getType(char).
func CharacterGetType(ch uint16) int {
	if ch == 0x0295 {
		return character_LOWERCASE_LETTER
	}
	r := rune(ch)
	if isUnassignedUnicode15(r) {
		return character_UNASSIGNED
	}
	for _, c := range characterCategories {
		if unicode.Is(c.table, r) {
			return c.typ
		}
	}
	return character_UNASSIGNED
}

// CharacterIsLetter is Character.isLetter(char): general category Lu, Ll,
// Lt, Lm or Lo.
func CharacterIsLetter(ch uint16) bool {
	if ch < 0x80 {
		return (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z')
	}
	r := rune(ch)
	return !isUnassignedUnicode15(r) && unicode.IsLetter(r) && !isSurrogate(ch)
}

// CharacterIsDigit is Character.isDigit(char): general category Nd.
func CharacterIsDigit(ch uint16) bool {
	if ch < 0x80 {
		return ch >= '0' && ch <= '9'
	}
	r := rune(ch)
	return !isUnassignedUnicode15(r) && unicode.Is(unicode.Nd, r)
}

// CharacterIsUpperCase is Character.isUpperCase(char): per its Javadoc,
// general category Lu or the Unicode property Other_Uppercase.
func CharacterIsUpperCase(ch uint16) bool {
	if ch < 0x80 {
		return ch >= 'A' && ch <= 'Z'
	}
	r := rune(ch)
	return !isUnassignedUnicode15(r) && unicode.In(r, unicode.Lu, unicode.Other_Uppercase)
}

func isSurrogate(ch uint16) bool { return ch >= 0xD800 && ch <= 0xDFFF }

// character_MIN_RADIX / character_MAX_RADIX are Character.MIN_RADIX/MAX_RADIX.
const (
	character_MIN_RADIX = 2
	character_MAX_RADIX = 36
)

// CharacterDigit is Character.digit(char, int): the value of ch in radix, or
// -1. Decimal digits (category Nd, which Unicode assigns in contiguous runs
// of 0..9), ASCII Latin letters and their fullwidth forms (U+FF21..U+FF3A,
// U+FF41..U+FF5A) are digits.
func CharacterDigit(ch uint16, radix int) int {
	if radix < character_MIN_RADIX || radix > character_MAX_RADIX {
		return -1
	}
	value := -1
	switch {
	case ch >= 'A' && ch <= 'Z':
		value = int(ch-'A') + 10
	case ch >= 'a' && ch <= 'z':
		value = int(ch-'a') + 10
	case ch >= 0xFF21 && ch <= 0xFF3A:
		value = int(ch-0xFF21) + 10
	case ch >= 0xFF41 && ch <= 0xFF5A:
		value = int(ch-0xFF41) + 10
	case CharacterIsDigit(ch):
		value = decimalDigitValue(rune(ch))
	}
	if value < radix {
		return value
	}
	return -1
}

// decimalDigitValue returns the numeric value of an Nd code point: its
// offset from the start of its contiguous Nd run, modulo 10 (Unicode
// guarantees Nd characters come in complete ascending 0..9 sequences).
func decimalDigitValue(r rune) int {
	for _, rg := range unicode.Nd.R16 {
		if r >= rune(rg.Lo) && r <= rune(rg.Hi) {
			if rg.Stride != 1 {
				break
			}
			return int(r-rune(rg.Lo)) % 10
		}
	}
	for _, rg := range unicode.Nd.R32 {
		if r >= rune(rg.Lo) && r <= rune(rg.Hi) {
			if rg.Stride != 1 {
				break
			}
			return int(r-rune(rg.Lo)) % 10
		}
	}
	return -1
}

// IntegerParseInt is Integer.parseInt(String s, int radix), from its
// Javadoc: an optional leading '-' or '+' (only when more characters
// follow), then one or more characters that Character.digit accepts in
// radix; the value must fit in an int. Failures throw NumberFormatException
// (messages: checkParseRadix, badNumberInput).
func IntegerParseInt(s string, radix int) int32 {
	v, ok := TryIntegerParseInt(s, radix)
	if !ok {
		panic(badNumberInput(s, radix))
	}
	return v
}

// TryIntegerParseInt parses s in radix with exact Java Integer.parseInt
// semantics without panicking on malformed number strings.
func TryIntegerParseInt(s string, radix int) (int32, bool) {
	checkParseRadix(radix)
	v, ok := parseRadixInteger(s, radix, math.MaxInt32)
	if !ok {
		return 0, false
	}
	return int32(v), true
}

// badNumberInput builds the NumberFormatException with which
// Integer.parseInt, Long.parseLong and new BigInteger(String, int) reject s
// as a number in radix. Observed (black-box): the message quotes s
// unchanged, and adds the radix when it is not 10.
func badNumberInput(s string, radix int) *NumberFormatException {
	msg := "For input string: \"" + s + "\""
	if radix != 10 {
		msg += " under radix " + itoa(radix)
	}
	return NewNumberFormatException(msg)
}

func itoa(i int) string { return ToJavaString(i) }

// Go-only: language-adapter — JDK 21 implements Unicode 15.0, whereas newer Go
// toolchains update the standard library's unicode tables to Unicode 16.0+.

var unicode15UnassignedAfter15 = &unicode.RangeTable{
	R16: []unicode.Range16{
		{0x088F, 0x088F, 1},
		{0x0897, 0x0897, 1},
		{0x0C5C, 0x0C5C, 1},
		{0x0CDC, 0x0CDC, 1},
		{0x1ACF, 0x1ADD, 1},
		{0x1AE0, 0x1AEB, 1},
		{0x1B4E, 0x1B4F, 1},
		{0x1B7F, 0x1B7F, 1},
		{0x1C89, 0x1C8A, 1},
		{0x20C1, 0x20C1, 1},
		{0x2427, 0x2429, 1},
		{0x2B96, 0x2B96, 1},
		{0x2FFC, 0x2FFF, 1},
		{0x31E4, 0x31E5, 1},
		{0x31EF, 0x31EF, 1},
		{0xA7CB, 0xA7CF, 1},
		{0xA7D2, 0xA7D2, 1},
		{0xA7D4, 0xA7D4, 1},
		{0xA7DA, 0xA7DC, 1},
		{0xA7F1, 0xA7F1, 1},
		{0xFBC3, 0xFBD2, 1},
		{0xFD90, 0xFD91, 1},
		{0xFDC8, 0xFDCE, 1},
	},
	R32: []unicode.Range32{
		{0x105C0, 0x105F3, 1},
		{0x10940, 0x10959, 1},
		{0x10D40, 0x10D65, 1},
		{0x10D69, 0x10D85, 1},
		{0x10D8E, 0x10D8F, 1},
		{0x10EC2, 0x10EC7, 1},
		{0x10ED0, 0x10ED8, 1},
		{0x10EFA, 0x10EFC, 1},
		{0x11380, 0x11389, 1},
		{0x1138B, 0x1138B, 1},
		{0x1138E, 0x1138E, 1},
		{0x11390, 0x113B5, 1},
		{0x113B7, 0x113C0, 1},
		{0x113C2, 0x113C2, 1},
		{0x113C5, 0x113C5, 1},
		{0x113C7, 0x113CA, 1},
		{0x113CC, 0x113D5, 1},
		{0x113D7, 0x113D8, 1},
		{0x113E1, 0x113E2, 1},
		{0x116D0, 0x116E3, 1},
		{0x11B60, 0x11B67, 1},
		{0x11BC0, 0x11BE1, 1},
		{0x11BF0, 0x11BF9, 1},
		{0x11DB0, 0x11DDB, 1},
		{0x11DE0, 0x11DE9, 1},
		{0x11F5A, 0x11F5A, 1},
		{0x13460, 0x143FA, 1},
		{0x16100, 0x16139, 1},
		{0x16D40, 0x16D79, 1},
		{0x16EA0, 0x16EB8, 1},
		{0x16EBB, 0x16ED3, 1},
		{0x16FF2, 0x16FF6, 1},
		{0x187F8, 0x187FF, 1},
		{0x18CFF, 0x18CFF, 1},
		{0x18D09, 0x18D1E, 1},
		{0x18D80, 0x18DF2, 1},
		{0x1CC00, 0x1CCFC, 1},
		{0x1CD00, 0x1CEB3, 1},
		{0x1CEBA, 0x1CED0, 1},
		{0x1CEE0, 0x1CEF0, 1},
		{0x1E5D0, 0x1E5FA, 1},
		{0x1E5FF, 0x1E5FF, 1},
		{0x1E6C0, 0x1E6DE, 1},
		{0x1E6E0, 0x1E6F5, 1},
		{0x1E6FE, 0x1E6FF, 1},
		{0x1F6D8, 0x1F6D8, 1},
		{0x1F777, 0x1F77A, 1},
		{0x1F8B2, 0x1F8BB, 1},
		{0x1F8C0, 0x1F8C1, 1},
		{0x1F8D0, 0x1F8D8, 1},
		{0x1FA54, 0x1FA57, 1},
		{0x1FA89, 0x1FA8A, 1},
		{0x1FA8E, 0x1FA8F, 1},
		{0x1FABE, 0x1FABE, 1},
		{0x1FAC6, 0x1FAC6, 1},
		{0x1FAC8, 0x1FAC8, 1},
		{0x1FACD, 0x1FACD, 1},
		{0x1FADC, 0x1FADC, 1},
		{0x1FADF, 0x1FADF, 1},
		{0x1FAE9, 0x1FAEA, 1},
		{0x1FAEF, 0x1FAEF, 1},
		{0x1FBCB, 0x1FBEF, 1},
		{0x1FBFA, 0x1FBFA, 1},
		{0x2B73A, 0x2B73F, 1},
		{0x2CEA2, 0x2CEAD, 1},
		{0x2EBF0, 0x2EE5D, 1},
		{0x323B0, 0x33479, 1},
	},
}

func isUnassignedUnicode15(r rune) bool {
	if r < 0x088F {
		return false
	}
	return unicode.Is(unicode15UnassignedAfter15, r)
}

// Unicode15ToUpper returns Character.toUpperCase(int) under JDK 21 (Unicode 15.0).
func Unicode15ToUpper(r rune) rune {
	return unicode15ToUpper(r)
}

func unicode15ToUpper(r rune) rune {
	switch r {
	case 0x019B, 0x0264, 0xA7D3, 0xA7D5:
		return r
	}
	if isUnassignedUnicode15(r) {
		return r
	}
	return unicode.ToUpper(r)
}

func unicode15ToLower(r rune) rune {
	if isUnassignedUnicode15(r) {
		return r
	}
	return unicode.ToLower(r)
}

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
	u := unicode15ToUpper(rune(ch))
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
	l := unicode15ToLower(rune(ch))
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
	return unicode15ToLower(unicode15ToUpper(cp))
}

func isHighSurrogateUnit(c uint16) bool { return c >= 0xD800 && c <= 0xDBFF }

func isLowSurrogateUnit(c uint16) bool { return c >= 0xDC00 && c <= 0xDFFF }

// Go-only: language-adapter — the java.lang.Character methods that take an
// int code point (character.go has the char overloads), and the Java
// identifier predicates. testdata/character_codepoint.tsv (JDK 21) checks
// every code point.

// CharacterGetTypeOfCodePoint is Character.getType(int codePoint).
func CharacterGetTypeOfCodePoint(codePoint int) int {
	if codePoint < 0 || codePoint > unicode.MaxRune {
		return character_UNASSIGNED
	}
	if codePoint <= 0xFFFF {
		return CharacterGetType(uint16(codePoint))
	}
	if codePoint == 0x1171E {
		return Character_NON_SPACING_MARK
	}
	r := rune(codePoint)
	if isUnassignedUnicode15(r) {
		return character_UNASSIGNED
	}
	for _, c := range characterCategories {
		if unicode.Is(c.table, r) {
			return c.typ
		}
	}
	return character_UNASSIGNED
}

// CharacterIsWhitespaceCodePoint is Character.isWhitespace(int): a Unicode
// space, line or paragraph separator other than the non-breaking spaces
// U+00A0, U+2007 and U+202F, or one of \t \n \u000B \f \r \u001C-\u001F.
func CharacterIsWhitespaceCodePoint(codePoint int) bool {
	switch codePoint {
	case '\t', '\n', 0x0B, '\f', '\r', 0x1C, 0x1D, 0x1E, 0x1F:
		return true
	case 0x00A0, 0x2007, 0x202F:
		return false
	}
	switch CharacterGetTypeOfCodePoint(codePoint) {
	case Character_SPACE_SEPARATOR, character_LINE_SEPARATOR, character_PARAGRAPH_SEPARATOR:
		return true
	}
	return false
}

// CharacterIsIdentifierIgnorable is Character.isIdentifierIgnorable(char):
// ISO control characters that are not whitespace (U+0000-U+0008,
// U+000E-U+001B, U+007F-U+009F) and FORMAT characters.
func CharacterIsIdentifierIgnorable(ch uint16) bool {
	if ch <= 0x08 || (ch >= 0x0E && ch <= 0x1B) || (ch >= 0x7F && ch <= 0x9F) {
		return true
	}
	return CharacterGetType(ch) == Character_FORMAT
}

// CharacterIsJavaIdentifierStart is Character.isJavaIdentifierStart(char):
// a letter, a letter number, a currency symbol or a connector punctuation.
func CharacterIsJavaIdentifierStart(ch uint16) bool {
	switch CharacterGetType(ch) {
	case character_UPPERCASE_LETTER, character_LOWERCASE_LETTER, character_TITLECASE_LETTER,
		character_MODIFIER_LETTER, character_OTHER_LETTER, character_LETTER_NUMBER,
		character_CURRENCY_SYMBOL, character_CONNECTOR_PUNCTUATION:
		return true
	}
	return false
}

// CharacterIsJavaIdentifierPart is Character.isJavaIdentifierPart(char): an
// identifier start, a digit, a combining or non-spacing mark, or an
// identifier-ignorable character.
func CharacterIsJavaIdentifierPart(ch uint16) bool {
	switch CharacterGetType(ch) {
	case character_UPPERCASE_LETTER, character_LOWERCASE_LETTER, character_TITLECASE_LETTER,
		character_MODIFIER_LETTER, character_OTHER_LETTER, character_LETTER_NUMBER,
		character_CURRENCY_SYMBOL, character_CONNECTOR_PUNCTUATION,
		character_DECIMAL_DIGIT_NUMBER, character_COMBINING_SPACING_MARK, Character_NON_SPACING_MARK:
		return true
	}
	return CharacterIsIdentifierIgnorable(ch)
}
