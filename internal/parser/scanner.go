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
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
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

// TemplateError returns the diagnostic / error-position record for a
// TemplateLiteralToken.
func (s *Scanner) TemplateError(tok Token) TemplateTokenError {
	if tok.aux == 0 {
		return TemplateTokenError{Position: tok.Location.Start}
	}
	idx := int(tok.aux - 1)
	if s != nil && idx < len(s.templateErrors) {
		return s.templateErrors[idx]
	}
	return standaloneTemplateErrors[idx]
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

// Java: Scanner.createToken(TokenType, int)
func (s *Scanner) CreateToken(type_ TokenType, beginToken int) Token {
	return NewToken(type_, s.GetTokenRange(beginToken))
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

// Scanner_ProcessUnicodeEscapes converts unicode escapes in the given
// string to the equivalent unicode character. If there are no escapes,
// returns the input unchanged. If there is an invalid escape sequence,
// returns nil (Java null).
//
// Java: Scanner.processUnicodeEscapes(String)
func Scanner_ProcessUnicodeEscapes(valueString string) *string {
	result, ok := scanner_processUnicodeEscapes(valueString)
	if !ok {
		return nil
	}
	return &result
}

func scanner_processUnicodeEscapes(valueString string) (string, bool) {
	// Go-only: language-adapter — the Java String is handled as UTF-16 code
	// units so that indices are char indices.
	value := jutil.StringToUTF16(valueString)
	for {
		escapeStart := slices.Index(value, '\\')
		if escapeStart < 0 {
			break
		}
		if escapeStart+2 >= len(value) || value[escapeStart+1] != 'u' {
			return "", false
		}

		var hexDigits []uint16
		var escapeEnd int
		if value[escapeStart+2] != '{' {
			// Simple escape with exactly four hex digits: \\uXXXX
			escapeEnd = escapeStart + 6
			if escapeEnd > len(value) {
				return "", false
			}
			hexDigits = value[escapeStart+2 : escapeEnd]
		} else {
			// Escape with braces can have any number of hex digits: \\u{XXXXXXX}
			escapeEnd = escapeStart + 3
			for escapeEnd < len(value) && Scanner_IsHexDigit(value[escapeEnd]) {
				escapeEnd++
			}
			if escapeEnd >= len(value) || value[escapeEnd] != '}' {
				return "", false
			}
			hexDigits = value[escapeStart+3 : escapeEnd]
			escapeEnd++
		}
		codePoint, ok := jutil.TryIntegerParseInt(jutil.StringFromUTF16(hexDigits), 0x10)
		if !ok || codePoint > 0xFFFF {
			return "", false
		}
		ch := uint16(codePoint)
		if !Identifiers_IsIdentifierPart(ch) {
			return "", false
		}
		value[escapeStart] = ch
		value = slices.Delete(value, escapeStart+1, escapeEnd)
	}
	return jutil.StringFromUTF16(value), true
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

// Java: Scanner.getTokenString(int)
func (s *Scanner) GetTokenString(beginIndex int) string {
	return s.sliceString(beginIndex, s.index)
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

// Java: Scanner.skipHexDigit()
func (s *Scanner) SkipHexDigit() bool {
	if !s.PeekHexDigit() {
		return false
	}
	s.NextChar()
	return true
}

// Java: Scanner.skipLineTerminator()
func (s *Scanner) SkipLineTerminator() {
	first := s.NextChar()
	if first == '\r' && s.Peek('\n') {
		s.NextChar()
	}
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
	switch ch {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true
	default:
		return false
	}
}

// Java: Scanner.peekHexDigit()
func (s *Scanner) PeekHexDigit() bool {
	return Scanner_IsHexDigit(s.PeekChar())
}

// Java: Scanner.isHexDigit(char)
func Scanner_IsHexDigit(ch uint16) bool {
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
	return Scanner_ValueOfOctalDigit(ch) >= 0
}

// Java: Scanner.valueOfOctalDigit(char)
func Scanner_ValueOfOctalDigit(ch uint16) int {
	switch ch {
	case '0', '1', '2', '3', '4', '5', '6', '7':
		return int(ch) - '0'
	default:
		return -1
	}
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
	return Scanner_ValueOfBinaryDigit(ch) >= 0
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
	s.errorReporter.ReportErrorFormat(s.source.LinePosition(position), format, arguments...)
}

// Java: Scanner.reportWarning(String, Object...)
func (s *Scanner) ReportWarning(format string, arguments ...any) {
	s.errorReporter.ReportWarningFormat(s.source.LinePosition(s.GetPosition()), format, arguments...)
}

// Java: Scanner.incTypeParameterLevel()
func (s *Scanner) IncTypeParameterLevel() {
	s.typeParameterLevel++
}

// Java: Scanner.decTypeParameterLevel()
func (s *Scanner) DecTypeParameterLevel() {
	s.typeParameterLevel--
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
