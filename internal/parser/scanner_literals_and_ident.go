// Port of com/google/javascript/jscomp/parsing/parser/Scanner.java @ closure-compiler v20261021

/*
 * Copyright 2011 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package parser

import (
	"strings"
	"unicode/utf8"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

func fnvHashString(s string) uint32 {
	h := uint32(2166136261)
	for i := 0; i < len(s); i++ {
		h ^= uint32(s[i])
		h *= 16777619
	}
	return h
}

func fnvHashBytes(b []byte) uint32 {
	h := uint32(2166136261)
	for i := 0; i < len(b); i++ {
		h ^= uint32(b[i])
		h *= 16777619
	}
	return h
}

func (s *Scanner) internIdent(str string) strpool.ID {
	slot := fnvHashString(str) & scannerIdentCacheMask
	if id := s.identCache[slot]; id != 0 && strpool.String(id) == str {
		return id
	}
	id := strpool.Intern(str)
	s.identCache[slot] = id
	return id
}

func (s *Scanner) internIdentBytes(b []byte) strpool.ID {
	slot := fnvHashBytes(b) & scannerIdentCacheMask
	if id := s.identCache[slot]; id != 0 && strpool.String(id) == string(b) {
		return id
	}
	id := strpool.InternBytes(b)
	s.identCache[slot] = id
	return id
}

// Java: Scanner.scanNumberPostPeriod(int)
func (s *Scanner) ScanNumberPostPeriod(beginToken int) Token {
	s.sawNumericSeparator = false
	s.SkipDecimalDigits()
	return s.ScanExponentOfNumericLiteral(beginToken)
}

// Java: Scanner.scanPostDigit(int)
func (s *Scanner) ScanPostDigit(beginToken int) Token {
	s.sawNumericSeparator = false
	s.SkipDecimalDigits()
	if s.Peek('n') {
		s.NextChar()
		return s.newNumberOrBigIntToken(TokenType_BIGINT, s.GetTokenRange(beginToken))
	}
	return s.ScanFractionalNumericLiteral(beginToken)
}

// Java: Scanner.scanPostZero(int)
func (s *Scanner) ScanPostZero(beginToken int) Token {
	s.sawNumericSeparator = false
	var isBigInt bool
	switch s.PeekChar() {
	case 'b', 'B':
		// binary
		s.NextChar()
		if !Scanner_IsBinaryDigit(s.PeekChar()) {
			s.ReportError("Binary Integer Literal must contain at least one digit")
		}
		s.SkipBinaryDigits()
		isBigInt = s.Peek('n')
		if isBigInt {
			s.NextChar()
		}
		return s.newNumberOrBigIntToken(bigIntOrNumber(isBigInt), s.GetTokenRange(beginToken))
	case 'o', 'O':
		// octal
		s.NextChar()
		if !Scanner_IsOctalDigit(s.PeekChar()) {
			s.ReportError("Octal Integer Literal must contain at least one digit")
		}
		s.SkipOctalDigits()
		if s.Peek('8') || s.Peek('9') {
			s.ReportError("Invalid octal digit in octal literal.")
		}
		isBigInt = s.Peek('n')
		if isBigInt {
			s.NextChar()
		}
		return s.newNumberOrBigIntToken(bigIntOrNumber(isBigInt), s.GetTokenRange(beginToken))
	case 'x', 'X':
		s.NextChar()
		if !s.PeekHexDigit() {
			s.ReportError("Hex Integer Literal must contain at least one digit")
		}
		s.SkipHexDigits()
		isBigInt = s.Peek('n')
		if isBigInt {
			s.NextChar()
		}
		return s.newNumberOrBigIntToken(bigIntOrNumber(isBigInt), s.GetTokenRange(beginToken))
	case 'e', 'E':
		return s.ScanExponentOfNumericLiteral(beginToken)
	case '.':
		return s.ScanFractionalNumericLiteral(beginToken)
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		s.SkipDecimalDigits()
		if s.Peek('.') {
			s.NextChar()
			s.SkipDecimalDigitsWithHasLeadingDigit( /* hasLeadingDigit= */ false)
		}
		if s.Peek('n') {
			s.ReportError("SyntaxError: nonzero BigInt can't have leading zero")
		}
		return s.newNumberOrBigIntToken(TokenType_NUMBER, s.GetTokenRange(beginToken))
	case 'n':
		s.NextChar()
		return s.newNumberOrBigIntToken(TokenType_BIGINT, s.GetTokenRange(beginToken))
	default:
		return s.newNumberOrBigIntToken(TokenType_NUMBER, s.GetTokenRange(beginToken))
	}
}

// bigIntOrNumber is the Java conditional `isBigInt ? TokenType.BIGINT : TokenType.NUMBER`.
//
// Go-only: language-adapter — Java's conditional operator.
func bigIntOrNumber(isBigInt bool) TokenType {
	if isBigInt {
		return TokenType_BIGINT
	}
	return TokenType_NUMBER
}

// Java: Scanner.scanFractionalNumericLiteral(int)
func (s *Scanner) ScanFractionalNumericLiteral(beginToken int) Token {
	if s.Peek('.') {
		s.NextChar()
		s.SkipDecimalDigitsWithHasLeadingDigit( /* hasLeadingDigit= */ false)
	}
	return s.ScanExponentOfNumericLiteral(beginToken)
}

// Java: Scanner.scanExponentOfNumericLiteral(int)
func (s *Scanner) ScanExponentOfNumericLiteral(beginToken int) Token {
	switch s.PeekChar() {
	case 'e', 'E':
		s.NextChar()
		switch s.PeekChar() {
		case '+', '-':
			s.NextChar()
		default: // fall out
		}
		if !Scanner_IsDecimalDigit(s.PeekChar()) {
			s.ReportError("Exponent part must contain at least one digit")
		}
		s.SkipDecimalDigits()
	default:
	}
	return s.newNumberOrBigIntToken(TokenType_NUMBER, s.GetTokenRange(beginToken))
}

// Java: Scanner.skipDecimalDigits()
func (s *Scanner) SkipDecimalDigits() {
	s.SkipDecimalDigitsWithHasLeadingDigit( /* hasLeadingDigit= */ true)
}

// Java: Scanner.skipDecimalDigits(boolean)
func (s *Scanner) SkipDecimalDigitsWithHasLeadingDigit(hasLeadingDigit bool) {
	hasDigit := hasLeadingDigit
	if s.index < s.asciiPrefix {
		raw := s.rawContents
		limit := s.asciiPrefix
		i := s.index
		for i < limit && raw[i] >= '0' && raw[i] <= '9' {
			i++
		}
		if i > s.index {
			hasDigit = true
			s.index = i
		}
	}
	ch := s.PeekChar()
	for Scanner_IsDecimalDigit(ch) || ch == '_' {
		s.NextChar()
		if ch == '_' {
			s.sawNumericSeparator = true
			if hasDigit && Scanner_IsDecimalDigit(s.PeekChar()) {
				s.NextChar()
			} else {
				s.ReportError("Trailing numeric separator")
			}
		} else {
			hasDigit = true
		}
		ch = s.PeekChar()
	}
}

// Java: Scanner.isDecimalDigit(char)
func Scanner_IsDecimalDigit(ch uint16) bool {
	return ch >= '0' && ch <= '9'
}

// Java: Scanner.peekHexDigit()
func (s *Scanner) PeekHexDigit() bool {
	return Scanner_IsHexDigit(s.PeekChar())
}

// Java: Scanner.isHexDigit(char)
func Scanner_IsHexDigit(ch uint16) bool {
	if ch < 0x80 {
		return (ch >= '0' && ch <= '9') || ((ch|0x20)-'a' < 6)
	}
	return jutil.CharacterDigit(ch, 0x10) >= 0
}

// Java: Scanner.skipHexDigits()
func (s *Scanner) SkipHexDigits() {
	ch := s.PeekChar()
	for Scanner_IsHexDigit(ch) || ch == '_' {
		s.NextChar()
		if ch == '_' {
			s.sawNumericSeparator = true
			if s.PeekHexDigit() {
				s.NextChar()
			} else {
				s.ReportError("Trailing numeric separator")
			}
		}
		ch = s.PeekChar()
	}
}

// Java: Scanner.skipHexDigit()
func (s *Scanner) SkipHexDigit() bool {
	if !s.PeekHexDigit() {
		return false
	}
	s.NextChar()
	return true
}

// Java: Scanner.skipOctalDigits()
func (s *Scanner) SkipOctalDigits() {
	ch := s.PeekChar()
	for Scanner_IsOctalDigit(ch) || ch == '_' {
		s.NextChar()
		if ch == '_' {
			s.sawNumericSeparator = true
			if Scanner_IsOctalDigit(s.PeekChar()) {
				s.NextChar()
			} else {
				s.ReportError("Trailing numeric separator")
			}
		}
		ch = s.PeekChar()
	}
}

// Java: Scanner.isOctalDigit(char)
func Scanner_IsOctalDigit(ch uint16) bool {
	return ch >= '0' && ch <= '7'
}

// Java: Scanner.valueOfOctalDigit(char)
func Scanner_ValueOfOctalDigit(ch uint16) int {
	if ch >= '0' && ch <= '7' {
		return int(ch) - '0'
	}
	return -1
}

// Java: Scanner.skipBinaryDigits()
func (s *Scanner) SkipBinaryDigits() {
	ch := s.PeekChar()
	for Scanner_IsBinaryDigit(ch) || ch == '_' {
		s.NextChar()
		if ch == '_' {
			s.sawNumericSeparator = true
			if Scanner_IsBinaryDigit(s.PeekChar()) {
				s.NextChar()
			} else {
				s.ReportError("Trailing numeric separator")
			}
		}
		ch = s.PeekChar()
	}
}

// Java: Scanner.isBinaryDigit(char)
func Scanner_IsBinaryDigit(ch uint16) bool {
	return ch == '0' || ch == '1'
}

// Java: Scanner.valueOfBinaryDigit(char)
func Scanner_ValueOfBinaryDigit(ch uint16) int {
	switch ch {
	case '0':
		return 0
	case '1':
		return 1
	default:
		return -1
	}
}

// Java: Scanner.scanIdentifierOrKeyword(int, char)
func (s *Scanner) ScanIdentifierOrKeyword(beginToken int, ch uint16) Token {
	valueStartIndex := s.index - 1

	containsUnicodeEscape := ch == '\\'
	bracedUnicodeEscape := false
	isPrivateIdentifier := ch == '#'
	unicodeEscapeLen := 0
	if containsUnicodeEscape {
		unicodeEscapeLen = 1
	} else {
		// Fast path for consecutive non-escaped identifier characters.
		if s.index < s.asciiPrefix {
			raw := s.rawContents
			limit := s.asciiPrefix
			i := s.index
			for i < limit && asciiIdentFlags[raw[i]]&asciiFlagIdentPart != 0 {
				i++
			}
			s.index = i
		}
		contentsLength := s.contentsLength
		for s.index < contentsLength {
			c := s.charAt(s.index)
			if c <= 127 {
				if asciiIdentFlags[c]&asciiFlagIdentPart != 0 {
					s.index++
					continue
				}
			} else if Identifiers_IsIdentifierPart(c) {
				s.index++
				continue
			}
			break
		}
	}

	ch = s.PeekChar()
	for Identifiers_IsIdentifierPart(ch) ||
		ch == '\\' ||
		(ch == '{' && unicodeEscapeLen == 2) ||
		(ch == '}' && bracedUnicodeEscape) {
		if ch == '\\' {
			containsUnicodeEscape = true
		}
		// Update length of current Unicode escape.
		if ch == '\\' || unicodeEscapeLen > 0 {
			unicodeEscapeLen++
		}
		// Enter Unicode point escape.
		if ch == '{' {
			bracedUnicodeEscape = true
		}
		// Exit Unicode escape
		if ch == '}' || (unicodeEscapeLen >= 6 && !bracedUnicodeEscape) {
			bracedUnicodeEscape = false
			unicodeEscapeLen = 0
		}

		// Add character to token
		s.NextChar()
		ch = s.PeekChar()
	}

	if !containsUnicodeEscape {
		if isPrivateIdentifier && s.index == valueStartIndex+1 {
			s.ReportErrorAtPosition(s.GetPositionAt(beginToken), "Invalid usage of #")
			return s.CreateToken(TokenType_ERROR, beginToken)
		}
		start := s.charAt(valueStartIndex)
		if isPrivateIdentifier {
			// Skip the leading # for name validation.
			start = s.charAt(valueStartIndex + 1)
		}
		if !Identifiers_IsIdentifierStart(start) {
			s.ReportErrorAtPosition(
				s.GetPositionAt(beginToken),
				"Character '%c' (U+%04X) is not a valid identifier start char",
				start,
				int(start))
			return s.CreateToken(TokenType_ERROR, beginToken)
		}
		loc := s.GetTokenRange(beginToken)
		if s.index <= s.asciiPrefix {
			value := s.rawContents[valueStartIndex:s.index]
			if !isPrivateIdentifier {
				if k := Keywords_Get(value); k != nil {
					return NewToken(k.Type(), loc)
				}
				return newIdentifierTokenFromID(loc, s.internIdent(value), 0)
			}
			return newIdentifierTokenFromID(loc, s.internIdent(value), tokenFlagPrivateIdent)
		}
		if valueStartIndex >= s.asciiPrefix {
			n := s.index - valueStartIndex
			if n <= 64 {
				tail := s.utf16Tail[valueStartIndex-s.asciiPrefix : s.index-s.asciiPrefix]
				var buf [64]byte
				allASCII := true
				for i, u := range tail {
					if u >= 0x80 {
						allASCII = false
						break
					}
					buf[i] = byte(u)
				}
				if allASCII {
					b := buf[:n]
					if !isPrivateIdentifier {
						if k := keywords_getBytes(b); k != nil {
							return NewToken(k.Type(), loc)
						}
						return newIdentifierTokenFromID(loc, s.internIdentBytes(b), 0)
					}
					return newIdentifierTokenFromID(loc, s.internIdentBytes(b), tokenFlagPrivateIdent)
				}
			}
		}
		value := s.sliceString(valueStartIndex, s.index)
		if !isPrivateIdentifier {
			if k := Keywords_Get(value); k != nil {
				return NewToken(k.Type(), loc)
			}
			return newIdentifierTokenFromID(loc, s.internIdent(value), 0)
		}
		return newIdentifierTokenFromID(loc, s.internIdent(value), tokenFlagPrivateIdent)
	}

	value := s.sliceString(valueStartIndex, s.index)
	if isPrivateIdentifier && value == "#" {
		s.ReportErrorAtPosition(s.GetPositionAt(beginToken), "Invalid usage of #")
		return s.CreateToken(TokenType_ERROR, beginToken)
	}

	processed, ok := scanner_processUnicodeEscapes(value)
	if !ok {
		s.ReportErrorAtPosition(s.GetPositionAt(s.index), "Invalid escape sequence")
		return s.CreateToken(TokenType_ERROR, beginToken)
	}
	value = processed
	start := jutil.StringCharAt(value, 0)
	if isPrivateIdentifier {
		start = jutil.StringCharAt(value, 1)
	}

	// Check to make sure the first character (or the unicode escape at the
	// beginning of the identifier) is a valid identifier start character.
	if !Identifiers_IsIdentifierStart(start) {
		s.ReportErrorAtPosition(
			s.GetPositionAt(beginToken),
			"Character '%c' (U+%04X) is not a valid identifier start char",
			start,
			int(start))
		return s.CreateToken(TokenType_ERROR, beginToken)
	}

	// Keywords containing unicode escape sequences are treated as
	// IdentifierTokens rather than keyword Tokens (ECMA-262 11.6.2).
	var flags uint8
	if isPrivateIdentifier {
		flags |= tokenFlagPrivateIdent
	} else if Keywords_Get(value) != nil {
		flags |= tokenFlagIdentIsKeyword
	}
	return newIdentifierTokenFromID(s.GetTokenRange(beginToken), s.internIdent(value), flags)
}

func scannerNextHexDigit(s string, pos int) (uint32, int, bool) {
	if pos >= len(s) {
		return 0, pos, false
	}
	b := s[pos]
	if b < 0x80 {
		switch {
		case b >= '0' && b <= '9':
			return uint32(b - '0'), pos + 1, true
		case b >= 'a' && b <= 'f':
			return uint32(b - 'a' + 10), pos + 1, true
		case b >= 'A' && b <= 'F':
			return uint32(b - 'A' + 10), pos + 1, true
		default:
			return 0, pos, false
		}
	}
	r, width := utf8.DecodeRuneInString(s[pos:])
	if (r == utf8.RuneError && width == 1) || r > 0xFFFF {
		return 0, pos, false
	}
	d := jutil.CharacterDigit(uint16(r), 16)
	if d < 0 {
		return 0, pos, false
	}
	return uint32(d), pos + width, true
}

func scanner_processUnicodeEscapes(valueString string) (string, bool) {
	firstSlash := strings.IndexByte(valueString, '\\')
	if firstSlash < 0 {
		return valueString, true
	}

	var sb strings.Builder
	sb.Grow(len(valueString))
	pos := 0
	slash := firstSlash
	for {
		sb.WriteString(valueString[pos : pos+slash])
		i := pos + slash
		if i+2 >= len(valueString) || valueString[i+1] != 'u' {
			return "", false
		}

		var codePoint uint32
		if valueString[i+2] != '{' {
			// Simple escape with exactly four hex digits: \uXXXX
			j := i + 2
			for k := 0; k < 4; k++ {
				nib, nextPos, ok := scannerNextHexDigit(valueString, j)
				if !ok {
					return "", false
				}
				codePoint = (codePoint << 4) | nib
				j = nextPos
			}
			pos = j
		} else {
			// Escape with braces can have any number of hex digits: \u{XXXXXXX}
			j := i + 3
			for j < len(valueString) {
				nib, nextPos, ok := scannerNextHexDigit(valueString, j)
				if !ok {
					break
				}
				if codePoint <= 0xFFFF {
					codePoint = (codePoint << 4) | nib
				}
				j = nextPos
			}
			if j == i+3 || j >= len(valueString) || valueString[j] != '}' || codePoint > 0xFFFF {
				return "", false
			}
			pos = j + 1
		}

		ch := uint16(codePoint)
		if !Identifiers_IsIdentifierPart(ch) {
			return "", false
		}
		if ch < 0x80 {
			sb.WriteByte(byte(ch))
		} else {
			sb.WriteRune(rune(ch))
		}

		slash = strings.IndexByte(valueString[pos:], '\\')
		if slash < 0 {
			sb.WriteString(valueString[pos:])
			break
		}
	}
	return sb.String(), true
}

// Java: Scanner.scanStringLiteral(int, char)
func (s *Scanner) ScanStringLiteral(beginIndex int, terminator uint16) Token {
	// String literals might span multiple lines.
	startingPosition := s.GetPositionAt(beginIndex)

	hasUnescapedUnicodeLineOrParagraphSeparator := false
	for {
		if s.index < s.asciiPrefix {
			raw := s.rawContents
			limit := s.asciiPrefix
			term := byte(terminator)
			i := s.index
			for i < limit {
				b := raw[i]
				if b == term || b == '\\' || b == '\n' || b == '\r' {
					break
				}
				i++
			}
			s.index = i
			if i < limit && raw[i] == term {
				s.index = i + 1
				return s.newStringLiteralToken(
					s.GetTokenRangeFromPosition(startingPosition),
					hasUnescapedUnicodeLineOrParagraphSeparator)
			}
		}
		if !s.PeekStringLiteralChar(terminator) {
			break
		}
		c := s.PeekChar()
		hasUnescapedUnicodeLineOrParagraphSeparator =
			hasUnescapedUnicodeLineOrParagraphSeparator || c == 0x2028 || c == 0x2029
		if !s.SkipStringLiteralChar() {
			return s.newStringLiteralToken(
				s.GetTokenRangeFromPosition(startingPosition),
				hasUnescapedUnicodeLineOrParagraphSeparator)
		}
	}
	if s.PeekChar() != terminator {
		s.ReportErrorAtPosition(startingPosition, "Unterminated string literal")
	} else {
		s.NextChar()
	}
	return s.newStringLiteralToken(
		s.GetTokenRangeFromPosition(startingPosition),
		hasUnescapedUnicodeLineOrParagraphSeparator)
}

// Java: Scanner.scanTemplateLiteral(int)
func (s *Scanner) ScanTemplateLiteral(beginIndex int) Token {
	if s.IsAtEnd() {
		s.ReportErrorAtPosition(s.GetPositionAt(beginIndex), "Unterminated template literal")
	}

	return s.NextTemplateLiteralTokenShared(
		TokenType_NO_SUBSTITUTION_TEMPLATE, TokenType_TEMPLATE_HEAD)
}

// Java: Scanner.nextTemplateLiteralTokenShared(TokenType, TokenType)
func (s *Scanner) NextTemplateLiteralTokenShared(endType, middleType TokenType) Token {
	beginIndex := s.index
	// Save the starting position to use with the multi-line safe version of getTokenRange().
	startingPosition := s.GetPositionAt(beginIndex)
	skipResult := s.SkipTemplateCharacters()
	if s.IsAtEnd() {
		s.ReportErrorAtPosition(startingPosition, "Unterminated template literal")
	}

	var flags uint8
	var aux uint16
	if skipResult.hasMsg || skipResult.hasLevel || skipResult.position != startingPosition {
		if skipResult.hasLevel && skipResult.errorLevel == TemplateLiteralToken_ErrorLevel_ERROR {
			flags |= tokenFlagTemplateHasError
		}
		var lvl uint8
		if skipResult.hasLevel {
			lvl = uint8(skipResult.errorLevel) + 1
		}
		s.templateErrors = append(s.templateErrors, TemplateTokenError{
			Message:  skipResult.errorMessage,
			Position: skipResult.position,
			HasMsg:   skipResult.hasMsg,
			Level:    lvl,
		})
		aux = uint16(len(s.templateErrors))
	}

	switch s.PeekChar() {
	case '`':
		s.NextChar()
		flags |= tokenFlagTemplateTrim1
		return newScannerTemplateToken(
			endType,
			s.GetTokenRangeFromPosition(startingPosition),
			flags,
			aux)
	case '$':
		s.NextChar() // $
		s.NextChar() // {
		flags |= tokenFlagTemplateTrim2
		return newScannerTemplateToken(
			middleType,
			s.GetTokenRangeFromPosition(startingPosition),
			flags,
			aux)
	default:
		// Should have reported error already
		return newScannerTemplateToken(
			endType,
			s.GetTokenRangeFromPosition(startingPosition),
			flags,
			aux)
	}
}

// Java: Scanner.peekStringLiteralChar(char)
func (s *Scanner) PeekStringLiteralChar(terminator uint16) bool {
	return !s.IsAtEnd() && s.PeekChar() != terminator && !Scanner_IsStringLineTerminator(s.PeekChar())
}

// Java: Scanner.skipStringLiteralChar()
func (s *Scanner) SkipStringLiteralChar() bool {
	if s.Peek('\\') {
		return s.SkipStringLiteralEscapeSequence()
	}
	s.NextChar()
	return true
}

// Java: Scanner.skipTemplateCharacters()
func (s *Scanner) SkipTemplateCharacters() Scanner_SkipTemplateCharactersResult {
	result := Scanner_SkipTemplateCharactersResult{position: s.GetPosition()}
	for !s.IsAtEnd() {
		if s.index < s.asciiPrefix {
			raw := s.rawContents
			limit := s.asciiPrefix
			i := s.index
			for i < limit {
				b := raw[i]
				if b == '`' || b == '\\' || b == '$' {
					break
				}
				i++
			}
			s.index = i
			if s.IsAtEnd() {
				break
			}
		}
		switch s.PeekChar() {
		case '`':
			return result
		case '\\':
			// There might be multiple errors. Take the first one but continue scanning
			newError, hasErr := s.SkipTemplateLiteralEscapeSequence()
			if hasErr && !result.HasError() {
				result = newError
			}
		case '$':
			if s.PeekCharAt(1) == '{' {
				return result
			}
			// Fall through.
			fallthrough
		default:
			s.NextChar()
		}
	}
	return result
}

// Java: Scanner.skipTemplateLiteralEscapeSequence()
func (s *Scanner) SkipTemplateLiteralEscapeSequence() (Scanner_SkipTemplateCharactersResult, bool) {
	s.NextChar()
	if s.IsAtEnd() {
		s.ReportError("Unterminated template literal escape sequence")
		return Scanner_SkipTemplateCharactersResult{}, false
	}
	if Scanner_IsLineTerminator(s.PeekChar()) {
		s.SkipLineTerminator()
		return Scanner_SkipTemplateCharactersResult{}, false
	}
	next := s.NextChar()
	switch next {
	case '0':
		if Scanner_IsDecimalDigit(s.PeekChar()) {
			return s.createTemplateDiag("Invalid escape sequence", TemplateLiteralToken_ErrorLevel_ERROR), true
		}
		return Scanner_SkipTemplateCharactersResult{}, false
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return s.createTemplateDiag("Invalid escape sequence", TemplateLiteralToken_ErrorLevel_ERROR), true
	case 'x':
		doubleHexDigit := s.SkipHexDigit() && s.SkipHexDigit()
		if !doubleHexDigit {
			return s.createTemplateDiag("Hex digit expected", TemplateLiteralToken_ErrorLevel_ERROR), true
		}
		return Scanner_SkipTemplateCharactersResult{}, false
	case 'u':
		if s.Peek('{') {
			s.NextChar()
			if s.Peek('}') {
				return s.createTemplateDiag("Empty unicode escape", TemplateLiteralToken_ErrorLevel_ERROR), true
			}
			allHexDigits := true
			for !s.Peek('}') && allHexDigits {
				allHexDigits = allHexDigits && s.SkipHexDigit()
			}
			if !allHexDigits {
				return s.createTemplateDiag("Hex digit expected", TemplateLiteralToken_ErrorLevel_ERROR), true
			}
			s.NextChar()
			return Scanner_SkipTemplateCharactersResult{}, false
		} else {
			quadHexDigit :=
				s.SkipHexDigit() && s.SkipHexDigit() && s.SkipHexDigit() && s.SkipHexDigit()
			if !quadHexDigit {
				return s.createTemplateDiag("Hex digit expected", TemplateLiteralToken_ErrorLevel_ERROR), true
			}
			return Scanner_SkipTemplateCharactersResult{}, false
		}
		// https://tc39.es/ecma262/#prod-TemplateEscapeSequence
	case '\\', 'b', 'f', 'n', 'r', 't', 'v', '$', '`':
		// special meaning in template literal
		return Scanner_SkipTemplateCharactersResult{}, false
	case '\'':
		// special the error message for a single quote
		return s.createTemplateDiag(
			jutil.LenientFormat("Unnecessary escape: \"\\%s\" is equivalent to just \"%s\"", next, next),
			TemplateLiteralToken_ErrorLevel_WARNING), true
	default:
		return s.createTemplateDiag(
			jutil.LenientFormat("Unnecessary escape: '\\%s' is equivalent to just '%s'", next, next),
			TemplateLiteralToken_ErrorLevel_WARNING), true
	}
}

// Java: Scanner.skipStringLiteralEscapeSequence()
func (s *Scanner) SkipStringLiteralEscapeSequence() bool {
	s.NextChar()
	if s.IsAtEnd() {
		s.ReportError("Unterminated string literal escape sequence")
		return false
	}
	if Scanner_IsStringLineTerminator(s.PeekChar()) {
		s.SkipLineTerminator()
		return true
	}

	next := s.NextChar()
	switch next {
	case '\'', '"', '`', '\\', 'b', 'f', 'n', 'r', 't', 'v', '0':
		return true
	case 'x':
		doubleHexDigit := s.SkipHexDigit() && s.SkipHexDigit()
		if !doubleHexDigit {
			s.ReportError("Hex digit expected")
		}
		return doubleHexDigit
	case 'u':
		if s.Peek('{') {
			s.NextChar()
			if s.Peek('}') {
				s.ReportError("Empty unicode escape")
				return false
			}
			allHexDigits := true
			for !s.Peek('}') && allHexDigits {
				allHexDigits = allHexDigits && s.SkipHexDigit()
			}
			if !allHexDigits {
				s.ReportError("Hex digit expected")
			}
			s.NextChar()
			return allHexDigits
		} else {
			quadHexDigit :=
				s.SkipHexDigit() && s.SkipHexDigit() && s.SkipHexDigit() && s.SkipHexDigit()
			if !quadHexDigit {
				s.ReportError("Hex digit expected")
			}
			return quadHexDigit
		}
	default:
	}

	if next == '/' {
		// Don't warn for '\/' (for now) since it's common in "<\/script>"
	} else {
		s.ReportWarning("Unnecessary escape: '\\%s' is equivalent to just '%s'", next, next)
	}
	return true
}

func (s *Scanner) createTemplateDiag(message string, errorLevel TemplateLiteralToken_ErrorLevel) Scanner_SkipTemplateCharactersResult {
	return Scanner_SkipTemplateCharactersResult{
		errorMessage: message,
		position:     s.GetPosition(),
		errorLevel:   errorLevel,
		hasMsg:       true,
		hasLevel:     true,
	}
}

// Scanner_SkipTemplateCharactersResult is Scanner.SkipTemplateCharactersResult
// represented as a value to avoid per-template heap allocation.
type Scanner_SkipTemplateCharactersResult struct {
	errorMessage string
	position     SourcePosition
	errorLevel   TemplateLiteralToken_ErrorLevel
	hasMsg       bool
	hasLevel     bool
}

// Java: Scanner.SkipTemplateCharactersResult.hasError()
func (r Scanner_SkipTemplateCharactersResult) HasError() bool {
	return r.hasMsg
}
