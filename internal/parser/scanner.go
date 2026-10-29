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

	"github.com/paullepoulpe/jsgo/internal/strpool"
)

const scannerIdentCacheMask = 127

// Scanner scans javascript source code into tokens. All entrypoints assume
// the caller is not expecting a regular expression literal except for
// NextRegularExpressionLiteralToken.
//
// 7 Lexical Conventions
type Scanner struct {
	errorReporter ErrorReporter
	source        *SourceFile
	// Go-only: language-adapter — the Java String contents, unwrapped from
	// the SourceFile (see SourceFile.CharAt).
	rawContents         string
	asciiPrefix         int
	utf16Tail           []uint16
	contentsLength      int
	lookahead           []Token
	lookaheadBuf        [8]Token
	index               int
	commentRecorder     Scanner_CommentRecorder
	typeParameterLevel  int
	sawNumericSeparator bool

	// L1 per-Scanner identifier intern cache avoiding strpool RWMutex/maphash
	// on repeated identifier names within a file.
	identCache [scannerIdentCacheMask + 1]strpool.ID

	// Side table for rare TemplateLiteralToken diagnostics / non-start
	// positions (indexed 1-based by Token.aux).
	templateErrors []TemplateTokenError
}

// TemplateError returns the diagnostic / error-position record for a
// TemplateLiteralToken.
func (s *Scanner) TemplateError(tok Token) TemplateTokenError {
	if tok.aux == 0 || s == nil {
		return TemplateTokenError{Position: tok.Location.Start}
	}
	return s.templateErrors[tok.aux-1]
}

func (s *Scanner) newStringLiteralToken(location SourceRange, hasUnescapedUnicodeLineOrParagraphSeparator bool) Token {
	var flags uint8
	if hasUnescapedUnicodeLineOrParagraphSeparator {
		flags = tokenFlagUnescapedLineOrParagraphSep
	}
	return newScannerLiteralToken(TokenType_STRING, location, flags)
}

func (s *Scanner) newNumberOrBigIntToken(type_ TokenType, location SourceRange) Token {
	var flags uint8
	if s.sawNumericSeparator {
		flags = tokenFlagNumericSeparator
	}
	return newScannerLiteralToken(type_, location, flags)
}

func (s *Scanner) newRegexToken(location SourceRange) Token {
	return newScannerLiteralToken(TokenType_REGULAR_EXPRESSION, location, 0)
}

// charAt is contents.charAt(index), for a valid index (SourceFile.CharAt).
func (s *Scanner) charAt(index int) uint16 {
	if index < s.asciiPrefix {
		return uint16(s.rawContents[index])
	}
	return s.utf16Tail[index-s.asciiPrefix]
}

// sliceString is contents.substring(start, end).
func (s *Scanner) sliceString(start, end int) string {
	return s.source.Substring(start, end)
}

// NewScanner constructs a Scanner.
//
// Java: Scanner(ErrorReporter, CommentRecorder, SourceFile, int)
func NewScanner(errorReporter ErrorReporter, commentRecorder Scanner_CommentRecorder, file *SourceFile, offset int) *Scanner {
	s := &Scanner{}
	s.lookahead = s.lookaheadBuf[:0]
	s.errorReporter = errorReporter
	s.commentRecorder = commentRecorder
	s.source = file
	// To help reason about the expected JVM performance unwrap "file" values.
	// The scanner is key to the parsing speed.
	s.rawContents = file.Contents
	s.asciiPrefix = file.asciiPrefix
	s.utf16Tail = file.utf16Tail
	s.contentsLength = file.Length()
	s.index = offset
	s.typeParameterLevel = 0
	return s
}

// Scanner_CommentRecorder is Scanner.CommentRecorder.
type Scanner_CommentRecorder interface {
	// Java: Scanner.CommentRecorder.recordComment(Comment.Type, SourceRange, String)
	RecordComment(type_ Comment_Type, range_ LineRange, value string)
}

// GetFile returns the source file.
//
// Java: Scanner.getFile()
func (s *Scanner) GetFile() *SourceFile {
	return s.source
}

// GetOffset returns the offset of the next token.
//
// Java: Scanner.getOffset()
func (s *Scanner) GetOffset() int {
	if len(s.lookahead) == 0 {
		return s.index
	}
	return int(s.lookahead[0].Location.Start.Offset)
}

// SetPosition moves the scanner to position, dropping any lookahead.
//
// Java: Scanner.setPosition(SourcePosition)
func (s *Scanner) SetPosition(position SourcePosition) {
	s.lookahead = s.lookahead[:0]
	s.index = int(position.Offset)
}

// GetPosition returns the position of the next token.
//
// Java: Scanner.getPosition()
func (s *Scanner) GetPosition() SourcePosition {
	if len(s.lookahead) == 0 {
		return s.GetPositionAt(s.index)
	}
	return s.lookahead[0].Location.Start
}

// GetPositionAt returns the source position of offset.
//
// Java: Scanner.getPosition(int)
func (s *Scanner) GetPositionAt(offset int) SourcePosition {
	return NewSourcePosition(offset)
}

// GetTokenRange returns the range from startOffset to the current index.
//
// Java: Scanner.getTokenRange(int)
func (s *Scanner) GetTokenRange(startOffset int) SourceRange {
	return NewSourceRange(NewSourcePosition(startOffset), NewSourcePosition(s.index))
}

// GetTokenRangeFromPosition: prefer this to GetTokenRange(int) when the
// token might span multiple lines.
//
// Java: Scanner.getTokenRange(SourcePosition)
func (s *Scanner) GetTokenRangeFromPosition(position SourcePosition) SourceRange {
	return NewSourceRange(position, NewSourcePosition(s.index))
}

// NextToken consumes and returns the next token.
//
// Java: Scanner.nextToken()
func (s *Scanner) NextToken() Token {
	n := len(s.lookahead)
	if n == 0 {
		return s.ScanToken()
	}
	tok := s.lookahead[0]
	if n == 1 {
		s.lookahead = s.lookahead[:0]
		return tok
	}
	copy(s.lookahead, s.lookahead[1:])
	s.lookahead = s.lookahead[:n-1]
	return tok
}

// ClearTokenLookahead rewinds to the first lookahead token, if any.
//
// Java: Scanner.clearTokenLookahead()
func (s *Scanner) ClearTokenLookahead() {
	if len(s.lookahead) > 0 {
		s.SetPosition(s.lookahead[0].Location.Start)
	}
}

// NextRegularExpressionLiteralToken scans a regular expression literal.
//
// Java: Scanner.nextRegularExpressionLiteralToken()
func (s *Scanner) NextRegularExpressionLiteralToken() Token {
	s.ClearTokenLookahead()

	beginToken := s.index

	// leading '/'
	s.NextChar()

	// body
	if !s.SkipRegularExpressionBody() {
		return s.newRegexToken(s.GetTokenRange(beginToken))
	}

	// separating '/'
	if s.PeekChar() != '/' {
		s.ReportError("Expected '/' in regular expression literal")
		return s.newRegexToken(s.GetTokenRange(beginToken))
	}
	s.NextChar()

	// flags
	for Identifiers_IsIdentifierPart(s.PeekChar()) {
		s.NextChar()
	}

	return s.newRegexToken(s.GetTokenRange(beginToken))
}

// NextTemplateLiteralToken scans the template literal portion that follows
// a substitution's closing '}'.
//
// Java: Scanner.nextTemplateLiteralToken()
func (s *Scanner) NextTemplateLiteralToken() Token {
	token := s.NextToken()
	if s.IsAtEnd() || token.Type() != TokenType_CLOSE_CURLY {
		s.ReportErrorAtPosition(s.GetPositionAt(s.index), "Expected '}' after expression in template literal")
	}

	return s.NextTemplateLiteralTokenShared(TokenType_TEMPLATE_TAIL, TokenType_TEMPLATE_MIDDLE)
}

// Java: Scanner.skipRegularExpressionBody()
func (s *Scanner) SkipRegularExpressionBody() bool {
	if !Scanner_IsRegularExpressionFirstChar(s.PeekChar()) {
		s.ReportError("Expected regular expression first char")
		return false
	}
	if !s.SkipRegularExpressionChar() {
		return false
	}
	for !s.IsAtEnd() && Scanner_IsRegularExpressionChar(s.PeekChar()) {
		if !s.SkipRegularExpressionChar() {
			return false
		}
	}
	return true
}

// Java: Scanner.skipRegularExpressionChar()
func (s *Scanner) SkipRegularExpressionChar() bool {
	switch s.PeekChar() {
	case '\\':
		return s.SkipRegularExpressionBackslashSequence()
	case '[':
		return s.SkipRegularExpressionClass()
	default:
		s.NextChar()
		return true
	}
}

// Java: Scanner.skipRegularExpressionBackslashSequence()
func (s *Scanner) SkipRegularExpressionBackslashSequence() bool {
	// TODO(tbreisacher): Warn if this is an unnecessary escape, like we do for string literals.
	s.NextChar()
	if Scanner_IsLineTerminator(s.PeekChar()) {
		s.ReportError("New line not allowed in regular expression literal")
		return false
	}
	s.NextChar()
	return true
}

// Java: Scanner.skipRegularExpressionClass()
func (s *Scanner) SkipRegularExpressionClass() bool {
	s.NextChar()
	for !s.IsAtEnd() && s.PeekRegularExpressionClassChar() {
		if !s.SkipRegularExpressionClassChar() {
			return false
		}
	}
	if s.PeekChar() != ']' {
		s.ReportError("']' expected")
		return false
	}
	s.NextChar()
	return true
}

// Java: Scanner.peekRegularExpressionClassChar()
func (s *Scanner) PeekRegularExpressionClassChar() bool {
	return s.PeekChar() != ']' && !Scanner_IsLineTerminator(s.PeekChar())
}

// Java: Scanner.skipRegularExpressionClassChar()
func (s *Scanner) SkipRegularExpressionClassChar() bool {
	if s.Peek('\\') {
		return s.SkipRegularExpressionBackslashSequence()
	}
	s.NextChar()
	return true
}

// Java: Scanner.isRegularExpressionFirstChar(char)
func Scanner_IsRegularExpressionFirstChar(ch uint16) bool {
	return Scanner_IsRegularExpressionChar(ch) && ch != '*'
}

// Java: Scanner.isRegularExpressionChar(char)
func Scanner_IsRegularExpressionChar(ch uint16) bool {
	switch ch {
	case '/':
		return false
	case '\\', '[':
		return true
	default:
		return !Scanner_IsLineTerminator(ch)
	}
}

// PeekToken returns the next token without consuming it.
//
// Java: Scanner.peekToken()
func (s *Scanner) PeekToken() Token {
	if len(s.lookahead) == 0 {
		tok := s.ScanToken()
		s.lookahead = append(s.lookahead, tok)
		return tok
	}
	return s.lookahead[0]
}

// PeekTokenAt returns the index-th next token without consuming any.
//
// Java: Scanner.peekToken(int)
func (s *Scanner) PeekTokenAt(index int) Token {
	for len(s.lookahead) <= index {
		s.lookahead = append(s.lookahead, s.ScanToken())
	}
	return s.lookahead[index]
}

// Java: Scanner.isAtEnd()
func (s *Scanner) IsAtEnd() bool {
	return !s.IsValidIndex(s.index)
}

// Java: Scanner.isValidIndex(int)
func (s *Scanner) IsValidIndex(index int) bool {
	return index >= 0 && index < s.contentsLength
}

// 7.2 White Space

// SkipWhitespace returns true if the whitespace that was skipped included
// any line terminators.
//
// Java: Scanner.skipWhitespace()
func (s *Scanner) SkipWhitespace() bool {
	foundLineTerminator := false
	if s.index < s.asciiPrefix {
		raw := s.rawContents
		limit := s.asciiPrefix
		i := s.index
		for i < limit {
			switch raw[i] {
			case ' ', '\t', 0x0B, 0x0C:
				i++
			case '\n', '\r':
				foundLineTerminator = true
				i++
			default:
				s.index = i
				return foundLineTerminator
			}
		}
		s.index = i
	}
	contentsLength := s.contentsLength
	for s.index < contentsLength {
		ch := s.charAt(s.index)
		if ch == ' ' || ch == '\t' {
			s.index++
			continue
		}
		if !Scanner_IsWhitespace(ch) {
			break
		}
		s.index++
		if Scanner_IsLineTerminator(ch) {
			foundLineTerminator = true
		}
	}
	return foundLineTerminator
}

// Java: Scanner.peekWhitespace()
func (s *Scanner) PeekWhitespace() bool {
	return Scanner_IsWhitespace(s.PeekChar())
}

// Java: Scanner.isWhitespace(char)
func Scanner_IsWhitespace(ch uint16) bool {
	switch ch {
	case 0x0009, // Tab
		0x000B, // Vertical Tab
		0x000C, // Form Feed
		0x0020, // Space
		0x00A0, // No-break space
		0xFEFF, // Byte Order Mark
		'\n',   // Line Feed
		'\r',   // Carriage Return
		0x2028, // Line Separator
		0x2029, // Paragraph Separator
		0x3000: // Ideographic Space
		// TODO: there are other Unicode Category 'Zs' chars that should go here.
		return true
	default:
		return false
	}
}

// 7.3 Line Terminators

// Java: Scanner.isLineTerminator(char)
func Scanner_IsLineTerminator(ch uint16) bool {
	switch ch {
	case '\n', // Line Feed
		'\r',   // Carriage Return
		0x2028, // Line Separator
		0x2029: // Paragraph Separator
		return true
	default:
		return false
	}
}

// Scanner_IsStringLineTerminator: allow line separator and paragraph
// separator in string literals.
// https://github.com/tc39/proposal-json-superset
//
// Java: Scanner.isStringLineTerminator(char)
func Scanner_IsStringLineTerminator(ch uint16) bool {
	switch ch {
	case 0x2028, // Line Separator
		0x2029: // Paragraph Separator
		return false
	default:
		return Scanner_IsLineTerminator(ch)
	}
}

// Java: Scanner.skipLineTerminator()
func (s *Scanner) SkipLineTerminator() {
	first := s.NextChar()
	if first == '\r' && s.Peek('\n') {
		s.NextChar()
	}
}

// 7.4 Comments

// Java: Scanner.skipComments()
func (s *Scanner) SkipComments() {
	for s.SkipComment() {
	}
}

// Java: Scanner.skipComment()
func (s *Scanner) SkipComment() bool {
	isStartOfLine := s.SkipWhitespace()
	if !s.IsAtEnd() {
		switch s.PeekCharAt(0) {
		case '/':
			switch s.PeekCharAt(1) {
			case '/':
				s.SkipSingleLineComment()
				return true
			case '*':
				s.SkipMultiLineComment()
				return true
			default: // fall out
			}
		case '<':
			// Check if this is the start of an HTML comment ("<!--").
			// http://www.w3.org/TR/REC-html40/interact/scripts.html#h-18.3.2
			if s.PeekCharAt(1) == '!' && s.PeekCharAt(2) == '-' && s.PeekCharAt(3) == '-' {
				s.ReportHtmlCommentWarning()
				s.SkipSingleLineComment()
				return true
			}
		case '-':
			// Check if this is the start of an HTML comment ("-->").
			// Note that the spec does not require us to check for this case,
			// but there is some legacy code that depends on this behavior.
			if isStartOfLine && s.PeekCharAt(1) == '-' && s.PeekCharAt(2) == '>' {
				s.ReportHtmlCommentWarning()
				s.SkipSingleLineComment()
				return true
			}
		case '#':
			if s.index == 0 && s.PeekCharAt(1) == '!' {
				s.SkipSingleLineCommentOfType(Comment_Type_SHEBANG)
				return true
			}
		default:
		}
	}
	return false
}

// Java: Scanner.reportHtmlCommentWarning()
func (s *Scanner) ReportHtmlCommentWarning() {
	s.ReportWarning(
		"In some cases, '<!--' and '-->' are treated as a '//' " +
			"for legacy reasons. Removing this from your code is " +
			"safe for all browsers currently in use.")
}

// Java: Scanner.skipSingleLineComment()
func (s *Scanner) SkipSingleLineComment() {
	s.SkipSingleLineCommentOfType(Comment_Type_LINE)
}

// Java: Scanner.skipSingleLineComment(Comment.Type)
func (s *Scanner) SkipSingleLineCommentOfType(type_ Comment_Type) {
	startOffset := s.index
	if s.index < s.asciiPrefix {
		raw := s.rawContents
		limit := s.asciiPrefix
		i := s.index
		for i < limit && raw[i] != '\n' && raw[i] != '\r' {
			i++
		}
		s.index = i
	}
	for !s.IsAtEnd() && !Scanner_IsLineTerminator(s.PeekChar()) {
		s.NextChar()
	}
	range_ := s.source.LineRange(startOffset, s.index)
	value := s.sliceString(startOffset, s.index)
	s.RecordComment(type_, range_, value)
}

// Java: Scanner.recordComment(Comment.Type, SourceRange, String)
func (s *Scanner) RecordComment(type_ Comment_Type, range_ LineRange, value string) {
	s.commentRecorder.RecordComment(type_, range_, value)
}

// Java: Scanner.skipMultiLineComment()
func (s *Scanner) SkipMultiLineComment() {
	startOffset := s.index
	s.NextChar() // '/'
	s.NextChar() // '*'
	for s.index+1 < s.asciiPrefix {
		rel := strings.IndexByte(s.rawContents[s.index:s.asciiPrefix-1], '*')
		if rel < 0 {
			s.index = s.asciiPrefix - 1
			break
		}
		s.index += rel
		if s.rawContents[s.index+1] == '/' {
			break
		}
		s.index++
	}
	for !s.IsAtEnd() && (s.PeekChar() != '*' || s.PeekCharAt(1) != '/') {
		s.NextChar()
	}
	if !s.IsAtEnd() {
		s.NextChar()
		s.NextChar()
		type_ := Comment_Type_BLOCK
		if s.index-startOffset > 4 {
			if s.charAt(startOffset+2) == '*' {
				type_ = Comment_Type_JSDOC
			} else if s.charAt(startOffset+2) == '!' {
				type_ = Comment_Type_IMPORTANT
			}
		}
		range_ := s.source.LineRange(startOffset, s.index)
		value := s.sliceString(startOffset, s.index)
		s.RecordComment(type_, range_, value)
	} else {
		s.ReportError("unterminated comment")
	}
}

// Java: Scanner.scanToken()
func (s *Scanner) ScanToken() Token {
	s.SkipComments()
	beginToken := s.index
	if s.IsAtEnd() {
		return s.CreateToken(TokenType_END_OF_FILE, beginToken)
	}
	ch := s.NextChar()
	switch ch {
	case '{':
		return s.CreateToken(TokenType_OPEN_CURLY, beginToken)
	case '}':
		return s.CreateToken(TokenType_CLOSE_CURLY, beginToken)
	case '(':
		return s.CreateToken(TokenType_OPEN_PAREN, beginToken)
	case ')':
		return s.CreateToken(TokenType_CLOSE_PAREN, beginToken)
	case '[':
		return s.CreateToken(TokenType_OPEN_SQUARE, beginToken)
	case ']':
		return s.CreateToken(TokenType_CLOSE_SQUARE, beginToken)
	case '.':
		if Scanner_IsDecimalDigit(s.PeekChar()) {
			return s.ScanNumberPostPeriod(beginToken)
		}

		// Harmony spread operator
		if s.Peek('.') && s.PeekCharAt(1) == '.' {
			s.NextChar()
			s.NextChar()
			return s.CreateToken(TokenType_ELLIPSIS, beginToken)
		}

		return s.CreateToken(TokenType_PERIOD, beginToken)
	case ';':
		return s.CreateToken(TokenType_SEMI_COLON, beginToken)
	case ',':
		return s.CreateToken(TokenType_COMMA, beginToken)
	case '~':
		return s.CreateToken(TokenType_TILDE, beginToken)
	case '?':
		if s.Peek('?') { // see ??
			s.NextChar()
			if s.Peek('=') {
				s.NextChar()
				return s.CreateToken(TokenType_QUESTION_QUESTION_EQUAL, beginToken)
			}
			return s.CreateToken(TokenType_QUESTION_QUESTION, beginToken)
		}
		if s.Peek('.') { // see ?.
			if !Scanner_IsDecimalDigit(s.PeekCharAt(1)) {
				s.NextChar()
				// a?.1:2 should be a ? 0.1 : 2 not a ?. 1 : 2 (syntax error)
				return s.CreateToken(TokenType_QUESTION_DOT, beginToken)
			}
		}
		return s.CreateToken(TokenType_QUESTION, beginToken)
	case ':':
		return s.CreateToken(TokenType_COLON, beginToken)
	case '<':
		switch s.PeekChar() {
		case '<':
			s.NextChar()
			if s.Peek('=') {
				s.NextChar()
				return s.CreateToken(TokenType_LEFT_SHIFT_EQUAL, beginToken)
			}
			return s.CreateToken(TokenType_LEFT_SHIFT, beginToken)
		case '=':
			s.NextChar()
			return s.CreateToken(TokenType_LESS_EQUAL, beginToken)
		default:
			return s.CreateToken(TokenType_OPEN_ANGLE, beginToken)
		}
	case '>':
		if s.typeParameterLevel > 0 {
			return s.CreateToken(TokenType_CLOSE_ANGLE, beginToken)
		}
		switch s.PeekChar() {
		case '>':
			s.NextChar()
			switch s.PeekChar() {
			case '=':
				s.NextChar()
				return s.CreateToken(TokenType_RIGHT_SHIFT_EQUAL, beginToken)
			case '>':
				s.NextChar()
				if s.Peek('=') {
					s.NextChar()
					return s.CreateToken(TokenType_UNSIGNED_RIGHT_SHIFT_EQUAL, beginToken)
				}
				return s.CreateToken(TokenType_UNSIGNED_RIGHT_SHIFT, beginToken)
			default:
				return s.CreateToken(TokenType_RIGHT_SHIFT, beginToken)
			}
		case '=':
			s.NextChar()
			return s.CreateToken(TokenType_GREATER_EQUAL, beginToken)
		default:
			return s.CreateToken(TokenType_CLOSE_ANGLE, beginToken)
		}
	case '=':
		switch s.PeekChar() {
		case '=':
			s.NextChar()
			if s.Peek('=') {
				s.NextChar()
				return s.CreateToken(TokenType_EQUAL_EQUAL_EQUAL, beginToken)
			}
			return s.CreateToken(TokenType_EQUAL_EQUAL, beginToken)
		case '>':
			s.NextChar()
			return s.CreateToken(TokenType_ARROW, beginToken)
		default:
			return s.CreateToken(TokenType_EQUAL, beginToken)
		}
	case '!':
		if s.Peek('=') {
			s.NextChar()
			if s.Peek('=') {
				s.NextChar()
				return s.CreateToken(TokenType_NOT_EQUAL_EQUAL, beginToken)
			}
			return s.CreateToken(TokenType_NOT_EQUAL, beginToken)
		}
		return s.CreateToken(TokenType_BANG, beginToken)
	case '*':
		if s.Peek('=') {
			s.NextChar()
			return s.CreateToken(TokenType_STAR_EQUAL, beginToken)
		} else if s.Peek('*') {
			s.NextChar()
			// '**' seen so far
			if s.Peek('=') {
				s.NextChar()
				return s.CreateToken(TokenType_STAR_STAR_EQUAL, beginToken)
			} else {
				return s.CreateToken(TokenType_STAR_STAR, beginToken)
			}
		}
		return s.CreateToken(TokenType_STAR, beginToken)
	case '%':
		if s.Peek('=') {
			s.NextChar()
			return s.CreateToken(TokenType_PERCENT_EQUAL, beginToken)
		}
		return s.CreateToken(TokenType_PERCENT, beginToken)
	case '^':
		if s.Peek('=') {
			s.NextChar()
			return s.CreateToken(TokenType_CARET_EQUAL, beginToken)
		}
		return s.CreateToken(TokenType_CARET, beginToken)
	case '/':
		if s.Peek('=') {
			s.NextChar()
			return s.CreateToken(TokenType_SLASH_EQUAL, beginToken)
		}
		return s.CreateToken(TokenType_SLASH, beginToken)
	case '+':
		switch s.PeekChar() {
		case '+':
			s.NextChar()
			return s.CreateToken(TokenType_PLUS_PLUS, beginToken)
		case '=':
			s.NextChar()
			return s.CreateToken(TokenType_PLUS_EQUAL, beginToken)
		default:
			return s.CreateToken(TokenType_PLUS, beginToken)
		}
	case '-':
		switch s.PeekChar() {
		case '-':
			s.NextChar()
			return s.CreateToken(TokenType_MINUS_MINUS, beginToken)
		case '=':
			s.NextChar()
			return s.CreateToken(TokenType_MINUS_EQUAL, beginToken)
		default:
			return s.CreateToken(TokenType_MINUS, beginToken)
		}
	case '&':
		switch s.PeekChar() {
		case '&':
			s.NextChar()
			if s.Peek('=') {
				s.NextChar()
				return s.CreateToken(TokenType_AND_EQUAL, beginToken)
			}
			return s.CreateToken(TokenType_AND, beginToken)
		case '=':
			s.NextChar()
			return s.CreateToken(TokenType_AMPERSAND_EQUAL, beginToken)
		default:
			return s.CreateToken(TokenType_AMPERSAND, beginToken)
		}
	case '|':
		switch s.PeekChar() {
		case '|':
			s.NextChar()
			if s.Peek('=') {
				s.NextChar()
				return s.CreateToken(TokenType_OR_EQUAL, beginToken)
			}
			return s.CreateToken(TokenType_OR, beginToken)
		case '=':
			s.NextChar()
			return s.CreateToken(TokenType_BAR_EQUAL, beginToken)
		default:
			return s.CreateToken(TokenType_BAR, beginToken)
		}
	case '#':
		// Shebang is not actually ever parsed here (when used correctly, it's handled above in the
		// skipComments() call) so its token is an error.
		if s.Peek('!') {
			s.ReportErrorAtPosition(s.GetPositionAt(s.index), "Shebang comment must be at the start of the file")
			return s.CreateToken(TokenType_ERROR, beginToken)
		}
		// Handle private identifiers.
		return s.ScanIdentifierOrKeyword(beginToken, ch)
		// TODO: add NumberToken
		// TODO: character following NumericLiteral must not be an IdentifierStart or DecimalDigit
	case '0':
		return s.ScanPostZero(beginToken)
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return s.ScanPostDigit(beginToken)
	case '"', '\'':
		return s.ScanStringLiteral(beginToken, ch)
	case '`':
		return s.ScanTemplateLiteral(beginToken)
	default:
		return s.ScanIdentifierOrKeyword(beginToken, ch)
	}
}

// Java: Scanner.createToken(TokenType, int)
func (s *Scanner) CreateToken(type_ TokenType, beginToken int) Token {
	return NewToken(type_, s.GetTokenRange(beginToken))
}

// Java: Scanner.getTokenString(int)
func (s *Scanner) GetTokenString(beginIndex int) string {
	return s.sliceString(beginIndex, s.index)
}

// Java: Scanner.nextChar()
func (s *Scanner) NextChar() uint16 {
	if s.IsAtEnd() {
		return 0
	}
	ch := s.charAt(s.index)
	s.index++
	return ch
}

// Java: Scanner.peek(char)
func (s *Scanner) Peek(ch uint16) bool {
	return s.PeekChar() == ch
}

// Java: Scanner.peekChar()
func (s *Scanner) PeekChar() uint16 {
	return s.PeekCharAt(0)
}

// Java: Scanner.peekChar(int)
func (s *Scanner) PeekCharAt(offset int) uint16 {
	if !s.IsValidIndex(s.index + offset) {
		return 0
	}
	return s.charAt(s.index + offset)
}

// Java: Scanner.reportError(String, Object...)
func (s *Scanner) ReportError(format string, arguments ...any) {
	s.ReportErrorAtPosition(s.GetPosition(), format, arguments...)
}

// Java: Scanner.reportError(SourcePosition, String, Object...)
func (s *Scanner) ReportErrorAtPosition(position SourcePosition, format string, arguments ...any) {
	reportErrorFormat(s.errorReporter, s.source.LinePosition(position), format, arguments...)
}

// Java: Scanner.reportWarning(String, Object...)
func (s *Scanner) ReportWarning(format string, arguments ...any) {
	reportWarningFormat(s.errorReporter, s.source.LinePosition(s.GetPosition()), format, arguments...)
}

// Java: Scanner.incTypeParameterLevel()
func (s *Scanner) IncTypeParameterLevel() {
	s.typeParameterLevel++
}

// Java: Scanner.decTypeParameterLevel()
func (s *Scanner) DecTypeParameterLevel() {
	s.typeParameterLevel--
}
