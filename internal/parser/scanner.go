// Port of com/google/javascript/jscomp/parsing/parser/Scanner.java @ closure-compiler v20260415

package parser

import (
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Scanner scans javascript source code into tokens. All entrypoints assume
// the caller is not expecting a regular expression literal except for
// NextRegularExpressionLiteralToken.
//
// 7 Lexical Conventions
type Scanner struct {
	errorReporter     ErrorReporter
	source            *SourceFile
	lineNumberScanner *LineNumberScanner
	// Go-only: language-adapter — the Java String contents as UTF-16 code
	// units, so that every index is a Java char index.
	contents           []uint16
	contentsLength     int
	currentTokens      *jutil.ArrayList[Token]
	index              int
	commentRecorder    Scanner_CommentRecorder
	typeParameterLevel int
}

// NewScanner constructs a Scanner.
//
// Java: Scanner(ErrorReporter, CommentRecorder, SourceFile, int)
func NewScanner(errorReporter ErrorReporter, commentRecorder Scanner_CommentRecorder, file *SourceFile, offset int) *Scanner {
	s := &Scanner{currentTokens: jutil.NewArrayList[Token]()}
	s.errorReporter = errorReporter
	s.commentRecorder = commentRecorder
	s.source = file
	s.lineNumberScanner = NewLineNumberScanner(s.source)
	// To help reason about the expected JVM performance unwrap "file" values.
	// The scanner is key to the parsing speed.
	s.contents = file.ContentsUTF16
	s.contentsLength = len(file.ContentsUTF16)
	s.index = offset
	s.typeParameterLevel = 0
	return s
}

// Scanner_CommentRecorder is Scanner.CommentRecorder.
type Scanner_CommentRecorder interface {
	// Java: Scanner.CommentRecorder.recordComment(Comment.Type, SourceRange, String)
	RecordComment(type_ Comment_Type, range_ *SourceRange, value string)
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
	if s.currentTokens.IsEmpty() {
		return s.index
	}
	return s.PeekToken().GetTokenBase().Location.Start.Offset
}

// SetPosition moves the scanner to position, dropping any lookahead.
//
// Java: Scanner.setPosition(SourcePosition)
func (s *Scanner) SetPosition(position *SourcePosition) {
	s.lineNumberScanner.RewindTo(position)
	s.currentTokens.Clear()
	s.index = position.Offset
}

// GetPosition returns the position of the next token.
//
// Java: Scanner.getPosition()
func (s *Scanner) GetPosition() *SourcePosition {
	if s.currentTokens.IsEmpty() {
		return s.GetPositionAt(s.index)
	}
	return s.PeekToken().GetTokenBase().Location.Start
}

// GetPositionAt returns the source position of offset.
//
// Java: Scanner.getPosition(int)
func (s *Scanner) GetPositionAt(offset int) *SourcePosition {
	return s.lineNumberScanner.GetSourcePosition(offset)
}

// GetTokenRange returns the range from startOffset to the current index.
//
// Java: Scanner.getTokenRange(int)
func (s *Scanner) GetTokenRange(startOffset int) *SourceRange {
	return s.lineNumberScanner.GetSourceRange(startOffset, s.index)
}

// GetTokenRangeFromPosition: prefer this to GetTokenRange(int) when the
// token might span multiple lines.
//
// Java: Scanner.getTokenRange(SourcePosition)
func (s *Scanner) GetTokenRangeFromPosition(position *SourcePosition) *SourceRange {
	s.lineNumberScanner.RewindTo(position)
	return s.lineNumberScanner.GetSourceRange(position.Offset, s.index)
}

// NextToken consumes and returns the next token.
//
// Java: Scanner.nextToken()
func (s *Scanner) NextToken() Token {
	s.PeekToken()
	return s.currentTokens.RemoveAt(0)
}

// ClearTokenLookahead rewinds to the first lookahead token, if any.
//
// Java: Scanner.clearTokenLookahead()
func (s *Scanner) ClearTokenLookahead() {
	if !s.currentTokens.IsEmpty() {
		s.SetPosition(s.PeekToken().GetTokenBase().Location.Start)
	}
}

// NextRegularExpressionLiteralToken scans a regular expression literal.
//
// Java: Scanner.nextRegularExpressionLiteralToken()
func (s *Scanner) NextRegularExpressionLiteralToken() LiteralToken {
	s.ClearTokenLookahead()

	beginToken := s.index

	// leading '/'
	s.NextChar()

	// body
	if !s.SkipRegularExpressionBody() {
		return NewLiteralToken(
			TokenType_REGULAR_EXPRESSION, s.GetTokenString(beginToken), s.GetTokenRange(beginToken))
	}

	// separating '/'
	if s.PeekChar() != '/' {
		s.ReportError("Expected '/' in regular expression literal")
		return NewLiteralToken(
			TokenType_REGULAR_EXPRESSION, s.GetTokenString(beginToken), s.GetTokenRange(beginToken))
	}
	s.NextChar()

	// flags
	for Identifiers_IsIdentifierPart(s.PeekChar()) {
		s.NextChar()
	}

	return NewLiteralToken(
		TokenType_REGULAR_EXPRESSION, s.GetTokenString(beginToken), s.GetTokenRange(beginToken))
}

// NextTemplateLiteralToken scans the template literal portion that follows
// a substitution's closing '}'.
//
// Java: Scanner.nextTemplateLiteralToken()
func (s *Scanner) NextTemplateLiteralToken() *TemplateLiteralToken {
	token := s.NextToken()
	if s.IsAtEnd() || token.GetTokenBase().Type != TokenType_CLOSE_CURLY {
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
	return s.PeekTokenAt(0)
}

// PeekTokenAt returns the index-th next token without consuming any.
//
// Java: Scanner.peekToken(int)
func (s *Scanner) PeekTokenAt(index int) Token {
	for s.currentTokens.Size() <= index {
		s.currentTokens.Add(s.ScanToken())
	}
	return s.currentTokens.Get(index)
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
	for !s.IsAtEnd() && s.PeekWhitespace() {
		if Scanner_IsLineTerminator(s.NextChar()) {
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
	for !s.IsAtEnd() && !Scanner_IsLineTerminator(s.PeekChar()) {
		s.NextChar()
	}
	range_ := s.lineNumberScanner.GetSourceRange(startOffset, s.index)
	value := jutil.StringFromUTF16(s.contents[startOffset:s.index])
	s.RecordComment(type_, range_, value)
}

// Java: Scanner.recordComment(Comment.Type, SourceRange, String)
func (s *Scanner) RecordComment(type_ Comment_Type, range_ *SourceRange, value string) {
	s.commentRecorder.RecordComment(type_, range_, value)
}

// Java: Scanner.skipMultiLineComment()
func (s *Scanner) SkipMultiLineComment() {
	startOffset := s.index
	s.NextChar() // '/'
	s.NextChar() // '*'
	for !s.IsAtEnd() && (s.PeekChar() != '*' || s.PeekCharAt(1) != '/') {
		s.NextChar()
	}
	if !s.IsAtEnd() {
		s.NextChar()
		s.NextChar()
		type_ := Comment_Type_BLOCK
		if s.index-startOffset > 4 {
			if s.contents[startOffset+2] == '*' {
				type_ = Comment_Type_JSDOC
			} else if s.contents[startOffset+2] == '!' {
				type_ = Comment_Type_IMPORTANT
			}
		}
		range_ := s.lineNumberScanner.GetSourceRange(startOffset, s.index)
		value := jutil.StringFromUTF16(s.contents[startOffset:s.index])
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
	s.SkipDecimalDigits()
	return s.ScanExponentOfNumericLiteral(beginToken)
}

// Java: Scanner.scanPostDigit(int)
func (s *Scanner) ScanPostDigit(beginToken int) Token {
	s.SkipDecimalDigits()
	if s.Peek('n') {
		s.NextChar()
		return NewLiteralToken(
			TokenType_BIGINT, s.GetTokenString(beginToken), s.GetTokenRange(beginToken))
	}
	return s.ScanFractionalNumericLiteral(beginToken)
}

// Java: Scanner.scanPostZero(int)
func (s *Scanner) ScanPostZero(beginToken int) Token {
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
		return NewLiteralToken(
			bigIntOrNumber(isBigInt),
			s.GetTokenString(beginToken),
			s.GetTokenRange(beginToken))
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
		return NewLiteralToken(
			bigIntOrNumber(isBigInt),
			s.GetTokenString(beginToken),
			s.GetTokenRange(beginToken))
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
		return NewLiteralToken(
			bigIntOrNumber(isBigInt),
			s.GetTokenString(beginToken),
			s.GetTokenRange(beginToken))
	case 'e', 'E':
		return s.ScanExponentOfNumericLiteral(beginToken)
	case '.':
		return s.ScanFractionalNumericLiteral(beginToken)
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		s.SkipDecimalDigits()
		if s.Peek('.') {
			s.NextChar()
			s.SkipDecimalDigits()
		}
		if s.Peek('n') {
			s.ReportError("SyntaxError: nonzero BigInt can't have leading zero")
		}
		return NewLiteralToken(
			TokenType_NUMBER, s.GetTokenString(beginToken), s.GetTokenRange(beginToken))
	case 'n':
		s.NextChar()
		return NewLiteralToken(
			TokenType_BIGINT, s.GetTokenString(beginToken), s.GetTokenRange(beginToken))
	default:
		return NewLiteralToken(
			TokenType_NUMBER, s.GetTokenString(beginToken), s.GetTokenRange(beginToken))
	}
}

// bigIntOrNumber is the Java conditional `isBigInt ? TokenType.BIGINT : TokenType.NUMBER`.
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
	// NOTE: This code previously used a StringBuilder to collect the characters of the identifier
	// or keyword. Recording the staring position and using contents.substring() below instead was
	// found to eliminate 1.84% of all JVM "frequently collected garbage" in the compilation of a
	// large project.
	valueStartIndex := s.index - 1

	containsUnicodeEscape := ch == '\\'
	bracedUnicodeEscape := false
	isPrivateIdentifier := ch == '#'
	unicodeEscapeLen := 0
	if containsUnicodeEscape {
		unicodeEscapeLen = 1
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

	value := jutil.StringFromUTF16(s.contents[valueStartIndex:s.index])

	if isPrivateIdentifier && value == "#" {
		s.ReportErrorAtPosition(s.GetPositionAt(beginToken), "Invalid usage of #")
		return s.CreateToken(TokenType_ERROR, beginToken)
	}

	// Process unicode escapes.
	if containsUnicodeEscape {
		processed := Scanner_ProcessUnicodeEscapes(value)
		if processed == nil {
			s.ReportErrorAtPosition(s.GetPositionAt(s.index), "Invalid escape sequence")
			return s.CreateToken(TokenType_ERROR, beginToken)
		}
		value = *processed
	}

	// Check to make sure the first character (or the unicode escape at the
	// beginning of the identifier) is a valid identifier start character.
	start := jutil.StringCharAt(value, 0)
	if isPrivateIdentifier {
		// Skip the leading # for name validation.
		start = jutil.StringCharAt(value, 1)
	}
	if !Identifiers_IsIdentifierStart(start) {
		s.ReportErrorAtPosition(
			s.GetPositionAt(beginToken),
			"Character '%c' (U+%04X) is not a valid identifier start char",
			start,
			int(start))
		return s.CreateToken(TokenType_ERROR, beginToken)
	}

	k := Keywords_Get(value)
	if k != nil {
		return NewToken(k.Type(), s.GetTokenRange(beginToken))
	}

	return NewIdentifierToken(s.GetTokenRange(beginToken), value)
}

// Scanner_ProcessUnicodeEscapes converts unicode escapes in the given
// string to the equivalent unicode character. If there are no escapes,
// returns the input unchanged. If there is an invalid escape sequence,
// returns nil (Java null).
//
// Java: Scanner.processUnicodeEscapes(String)
func Scanner_ProcessUnicodeEscapes(valueString string) *string {
	// Go-only: language-adapter — the Java String is handled as UTF-16 code
	// units so that indices are char indices; charAt past the end throws
	// StringIndexOutOfBoundsException like Java's.
	value := jutil.StringToUTF16(valueString)
	charAt := func(i int) uint16 {
		if i < 0 || i >= len(value) {
			panic(jutil.NewStringIndexOutOfBoundsException(i, len(value)))
		}
		return value[i]
	}
	for slices.Contains(value, '\\') {
		escapeStart := slices.Index(value, '\\')
		invalid := false
		jutil.Try(func() {
			if charAt(escapeStart+1) != 'u' {
				invalid = true
				return
			}

			var hexDigits []uint16
			var escapeEnd int
			if charAt(escapeStart+2) != '{' {
				// Simple escape with exactly four hex digits: \\uXXXX
				escapeEnd = escapeStart + 6
				// TODO(b/155480859): Don't trust String#substring to throw on out of bounds. J2CL
				// implements it incorrectly.
				if escapeEnd > len(value) {
					invalid = true
					return
				}
				hexDigits = value[escapeStart+2 : escapeEnd]
			} else {
				// Escape with braces can have any number of hex digits: \\u{XXXXXXX}
				escapeEnd = escapeStart + 3
				for Scanner_IsHexDigit(charAt(escapeEnd)) {
					escapeEnd++
				}
				if charAt(escapeEnd) != '}' {
					invalid = true
					return
				}
				hexDigits = value[escapeStart+3 : escapeEnd]
				escapeEnd++
			}
			// TODO(mattloring): Allow code points >= 0xFFFF (greater than the size of a char).
			ch := uint16(jutil.IntegerParseInt(jutil.StringFromUTF16(hexDigits), 0x10))
			if !Identifiers_IsIdentifierPart(ch) {
				invalid = true
				return
			}
			next := make([]uint16, 0, len(value))
			next = append(next, value[:escapeStart]...)
			next = append(next, ch)
			next = append(next, value[escapeEnd:]...)
			value = next
		}, func(e jutil.IRuntimeException) {
			switch e.(type) {
			case *jutil.NumberFormatException, *jutil.StringIndexOutOfBoundsException:
				invalid = true
			default:
				panic(e)
			}
		})
		if invalid {
			return nil
		}
	}
	result := jutil.StringFromUTF16(value)
	return &result
}

// Java: Scanner.scanStringLiteral(int, char)
func (s *Scanner) ScanStringLiteral(beginIndex int, terminator uint16) Token {
	// String literals might span multiple lines.
	startingPosition := s.GetPositionAt(beginIndex)

	hasUnescapedUnicodeLineOrParagraphSeparator := false
	for s.PeekStringLiteralChar(terminator) {
		c := s.PeekChar()
		hasUnescapedUnicodeLineOrParagraphSeparator =
			hasUnescapedUnicodeLineOrParagraphSeparator || c == 0x2028 || c == 0x2029
		if !s.SkipStringLiteralChar() {
			return NewStringLiteralToken(
				s.GetTokenString(beginIndex),
				s.GetTokenRangeFromPosition(startingPosition),
				hasUnescapedUnicodeLineOrParagraphSeparator)
		}
	}
	if s.PeekChar() != terminator {
		s.ReportErrorAtPosition(startingPosition, "Unterminated string literal")
	} else {
		s.NextChar()
	}
	return NewStringLiteralToken(
		s.GetTokenString(beginIndex),
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
func (s *Scanner) NextTemplateLiteralTokenShared(endType, middleType TokenType) *TemplateLiteralToken {
	beginIndex := s.index
	// Save the starting position to use with the multi-line safe version of getTokenRange().
	startingPosition := s.GetPositionAt(beginIndex)
	skipTemplateCharactersResult := s.SkipTemplateCharacters()
	if s.IsAtEnd() {
		s.ReportErrorAtPosition(startingPosition, "Unterminated template literal")
	}

	value := s.GetTokenString(beginIndex)
	switch s.PeekChar() {
	case '`':
		s.NextChar()
		return NewTemplateLiteralToken(
			endType,
			value,
			skipTemplateCharactersResult.GetErrorMessage(),
			skipTemplateCharactersResult.GetErrorLevel(),
			skipTemplateCharactersResult.GetPosition(),
			s.GetTokenRangeFromPosition(startingPosition))
	case '$':
		s.NextChar() // $
		s.NextChar() // {
		return NewTemplateLiteralToken(
			middleType,
			value,
			skipTemplateCharactersResult.GetErrorMessage(),
			skipTemplateCharactersResult.GetErrorLevel(),
			skipTemplateCharactersResult.GetPosition(),
			s.GetTokenRangeFromPosition(startingPosition))
	default:
		// Should have reported error already
		return NewTemplateLiteralToken(
			endType,
			value,
			skipTemplateCharactersResult.GetErrorMessage(),
			skipTemplateCharactersResult.GetErrorLevel(),
			skipTemplateCharactersResult.GetPosition(),
			s.GetTokenRangeFromPosition(startingPosition))
	}
}

// Java: Scanner.getTokenString(int)
func (s *Scanner) GetTokenString(beginIndex int) string {
	return jutil.StringFromUTF16(s.contents[beginIndex:s.index])
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
func (s *Scanner) SkipTemplateCharacters() *Scanner_SkipTemplateCharactersResult {
	result := s.CreateSkipTemplateCharactersResult(nil, nil)
	for !s.IsAtEnd() {
		switch s.PeekChar() {
		case '`':
			return result
		case '\\':
			// There might be multiple errors. Take the first one but continue scanning
			newError := s.SkipTemplateLiteralEscapeSequence()
			if newError != nil && !result.HasError() {
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
func (s *Scanner) SkipTemplateLiteralEscapeSequence() *Scanner_SkipTemplateCharactersResult {
	s.NextChar()
	if s.IsAtEnd() {
		s.ReportError("Unterminated template literal escape sequence")
		return nil
	}
	if Scanner_IsLineTerminator(s.PeekChar()) {
		s.SkipLineTerminator()
		return nil
	}
	next := s.NextChar()
	switch next {
	case '0':
		if Scanner_IsDecimalDigit(s.PeekChar()) {
			return s.CreateSkipTemplateCharactersResult(stringPtr("Invalid escape sequence"), errorLevelPtr(TemplateLiteralToken_ErrorLevel_ERROR))
		}
		return nil
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return s.CreateSkipTemplateCharactersResult(stringPtr("Invalid escape sequence"), errorLevelPtr(TemplateLiteralToken_ErrorLevel_ERROR))
	case 'x':
		doubleHexDigit := s.SkipHexDigit() && s.SkipHexDigit()
		if !doubleHexDigit {
			return s.CreateSkipTemplateCharactersResult(stringPtr("Hex digit expected"), errorLevelPtr(TemplateLiteralToken_ErrorLevel_ERROR))
		}
		return nil
	case 'u':
		if s.Peek('{') {
			s.NextChar()
			if s.Peek('}') {
				return s.CreateSkipTemplateCharactersResult(stringPtr("Empty unicode escape"), errorLevelPtr(TemplateLiteralToken_ErrorLevel_ERROR))
			}
			allHexDigits := true
			for !s.Peek('}') && allHexDigits {
				allHexDigits = allHexDigits && s.SkipHexDigit()
			}
			if !allHexDigits {
				return s.CreateSkipTemplateCharactersResult(stringPtr("Hex digit expected"), errorLevelPtr(TemplateLiteralToken_ErrorLevel_ERROR))
			}
			s.NextChar()
			return nil
		} else {
			quadHexDigit :=
				s.SkipHexDigit() && s.SkipHexDigit() && s.SkipHexDigit() && s.SkipHexDigit()
			if !quadHexDigit {
				return s.CreateSkipTemplateCharactersResult(stringPtr("Hex digit expected"), errorLevelPtr(TemplateLiteralToken_ErrorLevel_ERROR))
			}
			return nil
		}
		// https://tc39.es/ecma262/#prod-TemplateEscapeSequence
	case '\\', 'b', 'f', 'n', 'r', 't', 'v', '$', '`':
		// special meaning in template literal
		return nil
	case '\'':
		// special the error message for a single quote
		return s.CreateSkipTemplateCharactersResult(
			stringPtr(jutil.LenientFormat("Unnecessary escape: \"\\%s\" is equivalent to just \"%s\"", next, next)),
			errorLevelPtr(TemplateLiteralToken_ErrorLevel_WARNING))
	default:
		return s.CreateSkipTemplateCharactersResult(
			stringPtr(jutil.LenientFormat("Unnecessary escape: '\\%s' is equivalent to just '%s'", next, next)),
			errorLevelPtr(TemplateLiteralToken_ErrorLevel_WARNING))
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
func (s *Scanner) ScanFractionalNumericLiteral(beginToken int) LiteralToken {
	if s.Peek('.') {
		s.NextChar()
		s.SkipDecimalDigits()
	}
	return s.ScanExponentOfNumericLiteral(beginToken)
}

// Java: Scanner.scanExponentOfNumericLiteral(int)
func (s *Scanner) ScanExponentOfNumericLiteral(beginToken int) LiteralToken {
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
	return NewLiteralToken(
		TokenType_NUMBER, s.GetTokenString(beginToken), s.GetTokenRange(beginToken))
}

// Java: Scanner.skipDecimalDigits()
func (s *Scanner) SkipDecimalDigits() {
	ch := s.PeekChar()
	for Scanner_IsDecimalDigit(ch) || ch == '_' {
		s.NextChar()
		if ch == '_' {
			if Scanner_IsDecimalDigit(s.PeekChar()) {
				s.NextChar()
			} else {
				s.ReportError("Trailing numeric separator")
			}
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
	ch := s.contents[s.index]
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
	return s.contents[s.index+offset]
}

// Java: Scanner.reportError(String, Object...)
func (s *Scanner) ReportError(format string, arguments ...any) {
	s.ReportErrorAtPosition(s.GetPosition(), format, arguments...)
}

// Java: Scanner.reportError(SourcePosition, String, Object...)
func (s *Scanner) ReportErrorAtPosition(position *SourcePosition, format string, arguments ...any) {
	s.errorReporter.ReportErrorFormat(position, format, arguments...)
}

// Java: Scanner.reportWarning(String, Object...)
func (s *Scanner) ReportWarning(format string, arguments ...any) {
	s.errorReporter.ReportWarningFormat(s.GetPosition(), format, arguments...)
}

// Java: Scanner.incTypeParameterLevel()
func (s *Scanner) IncTypeParameterLevel() {
	s.typeParameterLevel++
}

// Java: Scanner.decTypeParameterLevel()
func (s *Scanner) DecTypeParameterLevel() {
	s.typeParameterLevel--
}

// Java: Scanner.createSkipTemplateCharactersResult(String, ErrorLevel)
func (s *Scanner) CreateSkipTemplateCharactersResult(message *string, errorLevel *TemplateLiteralToken_ErrorLevel) *Scanner_SkipTemplateCharactersResult {
	return NewScanner_SkipTemplateCharactersResult(message, errorLevel, s.GetPosition())
}

// Scanner_SkipTemplateCharactersResult is Scanner.SkipTemplateCharactersResult.
type Scanner_SkipTemplateCharactersResult struct {
	// Go-only: language-adapter — nil is Java null.
	errorMessage *string
	position     *SourcePosition
	// Go-only: language-adapter — nil is Java null.
	errorLevel *TemplateLiteralToken_ErrorLevel
}

// Java: Scanner.SkipTemplateCharactersResult(String, ErrorLevel, SourcePosition)
func NewScanner_SkipTemplateCharactersResult(message *string, errorLevel *TemplateLiteralToken_ErrorLevel, position *SourcePosition) *Scanner_SkipTemplateCharactersResult {
	r := &Scanner_SkipTemplateCharactersResult{}
	r.errorMessage = message
	r.errorLevel = errorLevel
	r.position = position
	return r
}

// Java: Scanner.SkipTemplateCharactersResult.getErrorMessage()
func (r *Scanner_SkipTemplateCharactersResult) GetErrorMessage() *string {
	return r.errorMessage
}

// Java: Scanner.SkipTemplateCharactersResult.getErrorLevel()
func (r *Scanner_SkipTemplateCharactersResult) GetErrorLevel() *TemplateLiteralToken_ErrorLevel {
	return r.errorLevel
}

// Java: Scanner.SkipTemplateCharactersResult.getPosition()
func (r *Scanner_SkipTemplateCharactersResult) GetPosition() *SourcePosition {
	return r.position
}

// Java: Scanner.SkipTemplateCharactersResult.hasError()
func (r *Scanner_SkipTemplateCharactersResult) HasError() bool {
	return r.errorMessage != nil
}

// Go-only: language-adapter — helpers for Java nullable String / enum
// values, which are pointers here.
func stringPtr(s string) *string { return &s }

func errorLevelPtr(l TemplateLiteralToken_ErrorLevel) *TemplateLiteralToken_ErrorLevel { return &l }
