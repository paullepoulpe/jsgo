package jutil

import "unicode"

// Go-only: language-adapter — the java.lang.Character methods that take an
// int code point (character.go has the char overloads), and the Java
// identifier predicates. testdata/character_codepoint.tsv (JDK 21) checks
// every code point.

// CharacterGetTypeOfCodePoint is Character.getType(int codePoint).
func CharacterGetTypeOfCodePoint(codePoint int) int {
	if codePoint < 0 || codePoint > unicode.MaxRune {
		return Character_UNASSIGNED
	}
	if codePoint <= 0xFFFF {
		return CharacterGetType(uint16(codePoint))
	}
	r := rune(codePoint)
	for _, c := range characterCategories {
		if unicode.Is(c.table, r) {
			return c.typ
		}
	}
	return Character_UNASSIGNED
}

// CharacterIsLetterCodePoint is Character.isLetter(int).
func CharacterIsLetterCodePoint(codePoint int) bool {
	switch CharacterGetTypeOfCodePoint(codePoint) {
	case Character_UPPERCASE_LETTER, Character_LOWERCASE_LETTER, Character_TITLECASE_LETTER,
		Character_MODIFIER_LETTER, Character_OTHER_LETTER:
		return true
	}
	return false
}

// CharacterIsDigitCodePoint is Character.isDigit(int).
func CharacterIsDigitCodePoint(codePoint int) bool {
	return CharacterGetTypeOfCodePoint(codePoint) == Character_DECIMAL_DIGIT_NUMBER
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
	case Character_SPACE_SEPARATOR, Character_LINE_SEPARATOR, Character_PARAGRAPH_SEPARATOR:
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
	case Character_UPPERCASE_LETTER, Character_LOWERCASE_LETTER, Character_TITLECASE_LETTER,
		Character_MODIFIER_LETTER, Character_OTHER_LETTER, Character_LETTER_NUMBER,
		Character_CURRENCY_SYMBOL, Character_CONNECTOR_PUNCTUATION:
		return true
	}
	return false
}

// CharacterIsJavaIdentifierPart is Character.isJavaIdentifierPart(char): an
// identifier start, a digit, a combining or non-spacing mark, or an
// identifier-ignorable character.
func CharacterIsJavaIdentifierPart(ch uint16) bool {
	switch CharacterGetType(ch) {
	case Character_UPPERCASE_LETTER, Character_LOWERCASE_LETTER, Character_TITLECASE_LETTER,
		Character_MODIFIER_LETTER, Character_OTHER_LETTER, Character_LETTER_NUMBER,
		Character_CURRENCY_SYMBOL, Character_CONNECTOR_PUNCTUATION,
		Character_DECIMAL_DIGIT_NUMBER, Character_COMBINING_SPACING_MARK, Character_NON_SPACING_MARK:
		return true
	}
	return CharacterIsIdentifierIgnorable(ch)
}
