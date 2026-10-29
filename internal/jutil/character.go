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
	Character_UNASSIGNED                = 0
	Character_UPPERCASE_LETTER          = 1
	Character_LOWERCASE_LETTER          = 2
	Character_TITLECASE_LETTER          = 3
	Character_MODIFIER_LETTER           = 4
	Character_OTHER_LETTER              = 5
	Character_NON_SPACING_MARK          = 6
	Character_ENCLOSING_MARK            = 7
	Character_COMBINING_SPACING_MARK    = 8
	Character_DECIMAL_DIGIT_NUMBER      = 9
	Character_LETTER_NUMBER             = 10
	Character_OTHER_NUMBER              = 11
	Character_SPACE_SEPARATOR           = 12
	Character_LINE_SEPARATOR            = 13
	Character_PARAGRAPH_SEPARATOR       = 14
	Character_CONTROL                   = 15
	Character_FORMAT                    = 16
	Character_PRIVATE_USE               = 18
	Character_SURROGATE                 = 19
	Character_DASH_PUNCTUATION          = 20
	Character_START_PUNCTUATION         = 21
	Character_END_PUNCTUATION           = 22
	Character_CONNECTOR_PUNCTUATION     = 23
	Character_OTHER_PUNCTUATION         = 24
	Character_MATH_SYMBOL               = 25
	Character_CURRENCY_SYMBOL           = 26
	Character_MODIFIER_SYMBOL           = 27
	Character_OTHER_SYMBOL              = 28
	Character_INITIAL_QUOTE_PUNCTUATION = 29
	Character_FINAL_QUOTE_PUNCTUATION   = 30
)

var characterCategories = []struct {
	table *unicode.RangeTable
	typ   int
}{
	{unicode.Lu, Character_UPPERCASE_LETTER},
	{unicode.Ll, Character_LOWERCASE_LETTER},
	{unicode.Lt, Character_TITLECASE_LETTER},
	{unicode.Lm, Character_MODIFIER_LETTER},
	{unicode.Lo, Character_OTHER_LETTER},
	{unicode.Mn, Character_NON_SPACING_MARK},
	{unicode.Me, Character_ENCLOSING_MARK},
	{unicode.Mc, Character_COMBINING_SPACING_MARK},
	{unicode.Nd, Character_DECIMAL_DIGIT_NUMBER},
	{unicode.Nl, Character_LETTER_NUMBER},
	{unicode.No, Character_OTHER_NUMBER},
	{unicode.Zs, Character_SPACE_SEPARATOR},
	{unicode.Zl, Character_LINE_SEPARATOR},
	{unicode.Zp, Character_PARAGRAPH_SEPARATOR},
	{unicode.Cc, Character_CONTROL},
	{unicode.Cf, Character_FORMAT},
	{unicode.Co, Character_PRIVATE_USE},
	{unicode.Cs, Character_SURROGATE},
	{unicode.Pd, Character_DASH_PUNCTUATION},
	{unicode.Ps, Character_START_PUNCTUATION},
	{unicode.Pe, Character_END_PUNCTUATION},
	{unicode.Pc, Character_CONNECTOR_PUNCTUATION},
	{unicode.Po, Character_OTHER_PUNCTUATION},
	{unicode.Sm, Character_MATH_SYMBOL},
	{unicode.Sc, Character_CURRENCY_SYMBOL},
	{unicode.Sk, Character_MODIFIER_SYMBOL},
	{unicode.So, Character_OTHER_SYMBOL},
	{unicode.Pi, Character_INITIAL_QUOTE_PUNCTUATION},
	{unicode.Pf, Character_FINAL_QUOTE_PUNCTUATION},
}

// CharacterGetType is Character.getType(char).
func CharacterGetType(ch uint16) int {
	if ch == 0x0295 {
		return Character_LOWERCASE_LETTER
	}
	r := rune(ch)
	if isUnassignedUnicode15(r) {
		return Character_UNASSIGNED
	}
	for _, c := range characterCategories {
		if unicode.Is(c.table, r) {
			return c.typ
		}
	}
	return Character_UNASSIGNED
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

// Character_MIN_RADIX / Character_MAX_RADIX are Character.MIN_RADIX/MAX_RADIX.
const (
	Character_MIN_RADIX = 2
	Character_MAX_RADIX = 36
)

// CharacterDigit is Character.digit(char, int): the value of ch in radix, or
// -1. Decimal digits (category Nd, which Unicode assigns in contiguous runs
// of 0..9), ASCII Latin letters and their fullwidth forms (U+FF21..U+FF3A,
// U+FF41..U+FF5A) are digits.
func CharacterDigit(ch uint16, radix int) int {
	if radix < Character_MIN_RADIX || radix > Character_MAX_RADIX {
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
