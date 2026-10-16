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
	r := rune(ch)
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
	return unicode.IsLetter(rune(ch)) && !isSurrogate(ch)
}

// CharacterIsDigit is Character.isDigit(char): general category Nd.
func CharacterIsDigit(ch uint16) bool {
	if ch < 0x80 {
		return ch >= '0' && ch <= '9'
	}
	return unicode.Is(unicode.Nd, rune(ch))
}

// CharacterIsUpperCase is Character.isUpperCase(char): per its Javadoc,
// general category Lu or the Unicode property Other_Uppercase.
func CharacterIsUpperCase(ch uint16) bool {
	if ch < 0x80 {
		return ch >= 'A' && ch <= 'Z'
	}
	return unicode.In(rune(ch), unicode.Lu, unicode.Other_Uppercase)
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
// (messages: checkParseRadix, numberFormatForInputString).
func IntegerParseInt(s string, radix int) int32 {
	checkParseRadix(radix)
	v, ok := parseRadixInteger(s, radix, math.MaxInt32)
	if !ok {
		panic(numberFormatForInputString(s, radix))
	}
	return int32(v)
}

// numberFormatForInputString is the NumberFormatException that
// Integer.parseInt, Long.parseLong and new BigInteger(String, int) throw for
// a string that is not a number in radix.
//
// Observed behaviour (black-box): the message quotes the string as given and
// names the radix unless it is 10.
func numberFormatForInputString(s string, radix int) *NumberFormatException {
	msg := "For input string: \"" + s + "\""
	if radix != 10 {
		msg += " under radix " + itoa(radix)
	}
	return NewNumberFormatException(msg)
}

func itoa(i int) string { return ToJavaString(i) }
