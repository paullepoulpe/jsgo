// Port of com/google/javascript/jscomp/parsing/parser/Scanner.java
// @ closure-compiler v20260415.
//
// Scans JavaScript source code into tokens. All entrypoints assume the
// caller is not expecting a regular expression literal except for
// NextRegularExpressionLiteralToken.
//
// 7 Lexical Conventions
//
// The token-level API (NextToken / PeekToken / NextRegularExpressionLiteralToken /
// NextTemplateLiteralToken) is a 1:1 port of Scanner.java. The
// character-level helpers (Peek/Advance/ConsumeKeyword/SkipTrivia/etc.)
// are preserved as exported methods so the existing char-level parser
// (parser.go) continues to compile unchanged; the parser rewrite (#141)
// will migrate to the token-level API.

package parser

import (
	"fmt"
	"strings"
	"unicode"
)

// CommentRecorder mirrors Scanner.CommentRecorder.
type CommentRecorder interface {
	RecordComment(typ CommentType, rang SourceRange, value string)
}

// ErrorReporter mirrors parsing.parser.util.ErrorReporter.
type ErrorReporter interface {
	ReportError(position SourcePosition, format string, args ...interface{})
	ReportWarning(position SourcePosition, format string, args ...interface{})
}

// ---- Scanner --------------------------------------------------------

// Scanner mirrors com.google.javascript.jscomp.parsing.parser.Scanner.
// It walks the source string rune-by-rune (index is a rune index, not byte).
// All public token-level methods are 1:1 with Java; all private helpers
// are named identically to their Java counterparts.
type Scanner struct {
	// token-level state (mirrors Java fields)
	errorReporter   ErrorReporter
	commentRecorder CommentRecorder
	contents        []rune
	contentsLength  int
	index           int
	currentTokens   []*Token
	typeParameterLevel int

	// char-level position tracking (used by both APIs)
	// line is 1-based, lineBOL is the rune index of the start of current line.
	line    int
	lineBOL int

	// lastTokenLine: line of the most-recently consumed non-trivia rune.
	// Used by the char-level parser's ASI path.
	lastTokenLine int
}

// NewScanner mirrors Java's `new Scanner(errorReporter, commentRecorder, file, offset)`.
// For the char-level-only use case (existing parser), errorReporter and
// commentRecorder may be nil.
func NewScanner(src string) *Scanner {
	rs := []rune(src)
	return &Scanner{
		contents:       rs,
		contentsLength: len(rs),
		index:          0,
		line:           1,
		lastTokenLine:  1,
	}
}

// NewScannerFull constructs a Scanner with a full error reporter and comment recorder,
// mirroring the Java constructor signature.
func NewScannerFull(er ErrorReporter, cr CommentRecorder, src string, offset int) *Scanner {
	rs := []rune(src)
	return &Scanner{
		errorReporter:   er,
		commentRecorder: cr,
		contents:        rs,
		contentsLength:  len(rs),
		index:           offset,
		line:            1,
		lastTokenLine:   1,
	}
}

// ---- Token-level API (1:1 with Java) --------------------------------

// NextToken mirrors Scanner.nextToken().
func (s *Scanner) NextToken() *Token {
	s.peekToken(0)
	tok := s.currentTokens[0]
	s.currentTokens = s.currentTokens[1:]
	return tok
}

// clearTokenLookahead mirrors Scanner.clearTokenLookahead().
func (s *Scanner) clearTokenLookahead() {
	if len(s.currentTokens) > 0 {
		s.setPosition(s.peekToken(0).Location.Start)
	}
}

// NextRegularExpressionLiteralToken mirrors Scanner.nextRegularExpressionLiteralToken().
func (s *Scanner) NextRegularExpressionLiteralToken() *LiteralToken {
	s.clearTokenLookahead()

	beginToken := s.index

	// leading '/'
	s.nextChar()

	// body
	if !s.skipRegularExpressionBody() {
		return newLiteralToken(TokenTypeRegularExpression, s.getTokenString(beginToken), s.getTokenRangeFromOffset(beginToken))
	}

	// separating '/'
	if s.peekChar() != '/' {
		s.reportError(s.getPositionAt(s.index), "Expected '/' in regular expression literal")
		return newLiteralToken(TokenTypeRegularExpression, s.getTokenString(beginToken), s.getTokenRangeFromOffset(beginToken))
	}
	s.nextChar()

	// flags
	for identifiersIsIdentifierPart(s.peekChar()) {
		s.nextChar()
	}

	return newLiteralToken(TokenTypeRegularExpression, s.getTokenString(beginToken), s.getTokenRangeFromOffset(beginToken))
}

// NextTemplateLiteralToken mirrors Scanner.nextTemplateLiteralToken().
func (s *Scanner) NextTemplateLiteralToken() *TemplateLiteralToken {
	tok := s.NextToken()
	if s.isAtEnd() || tok.Type != TokenTypeCloseCurly {
		s.reportError(s.getPositionAt(s.index), "Expected '}' after expression in template literal")
	}
	return s.nextTemplateLiteralTokenShared(TokenTypeTemplateTail, TokenTypeTemplateMiddle)
}

// PeekToken mirrors Scanner.peekToken().
func (s *Scanner) PeekTokenAt(index int) *Token {
	return s.peekToken(index)
}

// PeekToken returns the next token without consuming it.
func (s *Scanner) PeekToken() *Token {
	return s.peekToken(0)
}

// GetOffset mirrors Scanner.getOffset().
func (s *Scanner) GetOffset() int {
	if len(s.currentTokens) == 0 {
		return s.index
	}
	return s.peekToken(0).Location.Start.Offset
}

// SetPosition mirrors Scanner.setPosition(SourcePosition).
func (s *Scanner) setPosition(position SourcePosition) {
	s.currentTokens = nil
	s.index = position.Offset
	// Rewind line tracking.
	s.rewindLineTo(position)
}

// GetPosition mirrors Scanner.getPosition().
func (s *Scanner) GetPosition() SourcePosition {
	if len(s.currentTokens) == 0 {
		return s.getPositionAt(s.index)
	}
	return s.peekToken(0).Location.Start
}

// IncTypeParameterLevel mirrors Scanner.incTypeParameterLevel().
func (s *Scanner) IncTypeParameterLevel() { s.typeParameterLevel++ }

// DecTypeParameterLevel mirrors Scanner.decTypeParameterLevel().
func (s *Scanner) DecTypeParameterLevel() { s.typeParameterLevel-- }

// ---- Private token-level helpers ------------------------------------

// peekToken mirrors Scanner.peekToken(int).
func (s *Scanner) peekToken(idx int) *Token {
	for len(s.currentTokens) <= idx {
		s.currentTokens = append(s.currentTokens, s.scanToken())
	}
	return s.currentTokens[idx]
}

// isAtEnd mirrors Scanner.isAtEnd().
func (s *Scanner) isAtEnd() bool {
	return !s.isValidIndex(s.index)
}

// isValidIndex mirrors Scanner.isValidIndex(int).
func (s *Scanner) isValidIndex(idx int) bool {
	return idx >= 0 && idx < s.contentsLength
}

// ---- 7.2 White Space ------------------------------------------------

// skipWhitespace mirrors Scanner.skipWhitespace().
// Returns true if the whitespace included any line terminators.
func (s *Scanner) skipWhitespace() bool {
	foundLineTerminator := false
	for !s.isAtEnd() && s.peekWhitespace() {
		if isLineTerminator(s.nextChar()) {
			foundLineTerminator = true
		}
	}
	return foundLineTerminator
}

// peekWhitespace mirrors Scanner.peekWhitespace().
func (s *Scanner) peekWhitespace() bool {
	return scannerIsWhitespace(s.peekChar())
}

// scannerIsWhitespace mirrors Scanner.isWhitespace(char).
func scannerIsWhitespace(ch rune) bool {
	switch ch {
	case '	', // Tab
		'', // Vertical Tab
		'', // Form Feed
		' ', // Space
		' ', // No-break space
		'\ufeff', // Byte Order Mark
		'\n',     // Line Feed
		'\r',     // Carriage Return
		' ', // Line Separator
		' ', // Paragraph Separator
		'　': // Ideographic Space
		return true
	}
	return false
}

// ---- 7.3 Line Terminators -------------------------------------------

// isLineTerminator mirrors Scanner.isLineTerminator(char).
func isLineTerminator(ch rune) bool {
	switch ch {
	case '\n', '\r', ' ', ' ':
		return true
	}
	return false
}

// isStringLineTerminator mirrors Scanner.isStringLineTerminator(char).
// Allows line separator and paragraph separator in string literals per
// https://github.com/tc39/proposal-json-superset
func isStringLineTerminator(ch rune) bool {
	switch ch {
	case ' ', ' ':
		return false
	}
	return isLineTerminator(ch)
}

// ---- 7.4 Comments ---------------------------------------------------

// skipComments mirrors Scanner.skipComments().
func (s *Scanner) skipComments() {
	for s.skipComment() {
	}
}

// skipComment mirrors Scanner.skipComment().
func (s *Scanner) skipComment() bool {
	isStartOfLine := s.skipWhitespace()
	if !s.isAtEnd() {
		switch s.peekChar() {
		case '/':
			switch s.peekCharAt(1) {
			case '/':
				s.skipSingleLineCommentType(CommentLine)
				return true
			case '*':
				s.skipMultiLineComment()
				return true
			}
		case '<':
			// Check if this is the start of an HTML comment ("<!--").
			if s.peekCharAt(1) == '!' && s.peekCharAt(2) == '-' && s.peekCharAt(3) == '-' {
				s.reportHtmlCommentWarning()
				s.skipSingleLineCommentType(CommentLine)
				return true
			}
		case '-':
			// Check if this is the start of an HTML comment ("-->").
			if isStartOfLine && s.peekCharAt(1) == '-' && s.peekCharAt(2) == '>' {
				s.reportHtmlCommentWarning()
				s.skipSingleLineCommentType(CommentLine)
				return true
			}
		case '#':
			if s.index == 0 && s.peekCharAt(1) == '!' {
				s.skipSingleLineCommentType(CommentShebang)
				return true
			}
		}
	}
	return false
}

// reportHtmlCommentWarning mirrors Scanner.reportHtmlCommentWarning().
func (s *Scanner) reportHtmlCommentWarning() {
	s.reportWarningAt(s.getPositionAt(s.index),
		"In some cases, '<!--' and '-->' are treated as a '//' for legacy reasons. Removing this from your code is safe for all browsers currently in use.")
}

// skipSingleLineComment mirrors Scanner.skipSingleLineComment().
func (s *Scanner) skipSingleLineComment() {
	s.skipSingleLineCommentType(CommentLine)
}

// skipSingleLineCommentType mirrors Scanner.skipSingleLineComment(Comment.Type).
func (s *Scanner) skipSingleLineCommentType(typ CommentType) {
	startOffset := s.index
	startPos := s.getPositionAt(startOffset)
	for !s.isAtEnd() && !isLineTerminator(s.peekChar()) {
		s.nextChar()
	}
	rang := SourceRange{Start: startPos, End: s.getPositionAt(s.index)}
	value := string(s.contents[startOffset:s.index])
	s.recordComment(typ, rang, value)
}

// skipMultiLineComment mirrors Scanner.skipMultiLineComment().
func (s *Scanner) skipMultiLineComment() {
	startOffset := s.index
	startPos := s.getPositionAt(startOffset)
	s.nextChar() // '/'
	s.nextChar() // '*'
	for !s.isAtEnd() && !(s.peekChar() == '*' && s.peekCharAt(1) == '/') {
		s.nextChar()
	}
	if !s.isAtEnd() {
		s.nextChar() // '*'
		s.nextChar() // '/'
		typ := CommentBlock
		if s.index-startOffset > 4 {
			if s.contents[startOffset+2] == '*' {
				typ = CommentJSDoc
			} else if s.contents[startOffset+2] == '!' {
				typ = CommentImportant
			}
		}
		rang := SourceRange{Start: startPos, End: s.getPositionAt(s.index)}
		value := string(s.contents[startOffset:s.index])
		s.recordComment(typ, rang, value)
	} else {
		s.reportError(s.getPositionAt(s.index), "unterminated comment")
	}
}

// recordComment mirrors Scanner.recordComment().
func (s *Scanner) recordComment(typ CommentType, rang SourceRange, value string) {
	if s.commentRecorder != nil {
		s.commentRecorder.RecordComment(typ, rang, value)
	}
}

// ---- scanToken (the core dispatch) ----------------------------------

// scanToken mirrors Scanner.scanToken().
func (s *Scanner) scanToken() *Token {
	s.skipComments()
	beginToken := s.index
	if s.isAtEnd() {
		return s.createToken(TokenTypeEndOfFile, beginToken)
	}
	ch := s.nextChar()
	switch ch {
	case '{':
		return s.createToken(TokenTypeOpenCurly, beginToken)
	case '}':
		return s.createToken(TokenTypeCloseCurly, beginToken)
	case '(':
		return s.createToken(TokenTypeOpenParen, beginToken)
	case ')':
		return s.createToken(TokenTypeCloseParen, beginToken)
	case '[':
		return s.createToken(TokenTypeOpenSquare, beginToken)
	case ']':
		return s.createToken(TokenTypeCloseSquare, beginToken)
	case '.':
		if isDecimalDigit(s.peekChar()) {
			return s.scanNumberPostPeriod(beginToken)
		}
		// Harmony spread operator
		if s.peekChar() == '.' && s.peekCharAt(1) == '.' {
			s.nextChar()
			s.nextChar()
			return s.createToken(TokenTypeEllipsis, beginToken)
		}
		return s.createToken(TokenTypePeriod, beginToken)
	case ';':
		return s.createToken(TokenTypeSemiColon, beginToken)
	case ',':
		return s.createToken(TokenTypeComma, beginToken)
	case '~':
		return s.createToken(TokenTypeTilde, beginToken)
	case '?':
		if s.peekChar() == '?' { // see ??
			s.nextChar()
			if s.peekChar() == '=' {
				s.nextChar()
				return s.createToken(TokenTypeQuestionQuestionEqual, beginToken)
			}
			return s.createToken(TokenTypeQuestionQuestion, beginToken)
		}
		if s.peekChar() == '.' { // see ?.
			if !isDecimalDigit(s.peekCharAt(1)) {
				s.nextChar()
				// a?.1:2 should be a ? 0.1 : 2 not a ?. 1 : 2 (syntax error)
				return s.createToken(TokenTypeQuestionDot, beginToken)
			}
		}
		return s.createToken(TokenTypeQuestion, beginToken)
	case ':':
		return s.createToken(TokenTypeColon, beginToken)
	case '<':
		switch s.peekChar() {
		case '<':
			s.nextChar()
			if s.peekChar() == '=' {
				s.nextChar()
				return s.createToken(TokenTypeLeftShiftEqual, beginToken)
			}
			return s.createToken(TokenTypeLeftShift, beginToken)
		case '=':
			s.nextChar()
			return s.createToken(TokenTypeLessEqual, beginToken)
		default:
			return s.createToken(TokenTypeOpenAngle, beginToken)
		}
	case '>':
		if s.typeParameterLevel > 0 {
			return s.createToken(TokenTypeCloseAngle, beginToken)
		}
		switch s.peekChar() {
		case '>':
			s.nextChar()
			switch s.peekChar() {
			case '=':
				s.nextChar()
				return s.createToken(TokenTypeRightShiftEqual, beginToken)
			case '>':
				s.nextChar()
				if s.peekChar() == '=' {
					s.nextChar()
					return s.createToken(TokenTypeUnsignedRightShiftEqual, beginToken)
				}
				return s.createToken(TokenTypeUnsignedRightShift, beginToken)
			default:
				return s.createToken(TokenTypeRightShift, beginToken)
			}
		case '=':
			s.nextChar()
			return s.createToken(TokenTypeGreaterEqual, beginToken)
		default:
			return s.createToken(TokenTypeCloseAngle, beginToken)
		}
	case '=':
		switch s.peekChar() {
		case '=':
			s.nextChar()
			if s.peekChar() == '=' {
				s.nextChar()
				return s.createToken(TokenTypeEqualEqualEqual, beginToken)
			}
			return s.createToken(TokenTypeEqualEqual, beginToken)
		case '>':
			s.nextChar()
			return s.createToken(TokenTypeArrow, beginToken)
		default:
			return s.createToken(TokenTypeEqual, beginToken)
		}
	case '!':
		if s.peekChar() == '=' {
			s.nextChar()
			if s.peekChar() == '=' {
				s.nextChar()
				return s.createToken(TokenTypeNotEqualEqual, beginToken)
			}
			return s.createToken(TokenTypeNotEqual, beginToken)
		}
		return s.createToken(TokenTypeBang, beginToken)
	case '*':
		if s.peekChar() == '=' {
			s.nextChar()
			return s.createToken(TokenTypeStarEqual, beginToken)
		} else if s.peekChar() == '*' {
			s.nextChar()
			if s.peekChar() == '=' {
				s.nextChar()
				return s.createToken(TokenTypeStarStarEqual, beginToken)
			}
			return s.createToken(TokenTypeStarStar, beginToken)
		}
		return s.createToken(TokenTypeStar, beginToken)
	case '%':
		if s.peekChar() == '=' {
			s.nextChar()
			return s.createToken(TokenTypePercentEqual, beginToken)
		}
		return s.createToken(TokenTypePercent, beginToken)
	case '^':
		if s.peekChar() == '=' {
			s.nextChar()
			return s.createToken(TokenTypeCaretEqual, beginToken)
		}
		return s.createToken(TokenTypeCaret, beginToken)
	case '/':
		if s.peekChar() == '=' {
			s.nextChar()
			return s.createToken(TokenTypeSlashEqual, beginToken)
		}
		return s.createToken(TokenTypeSlash, beginToken)
	case '+':
		switch s.peekChar() {
		case '+':
			s.nextChar()
			return s.createToken(TokenTypePlusPlus, beginToken)
		case '=':
			s.nextChar()
			return s.createToken(TokenTypePlusEqual, beginToken)
		default:
			return s.createToken(TokenTypePlus, beginToken)
		}
	case '-':
		switch s.peekChar() {
		case '-':
			s.nextChar()
			return s.createToken(TokenTypeMinusMinus, beginToken)
		case '=':
			s.nextChar()
			return s.createToken(TokenTypeMinusEqual, beginToken)
		default:
			return s.createToken(TokenTypeMinus, beginToken)
		}
	case '&':
		switch s.peekChar() {
		case '&':
			s.nextChar()
			if s.peekChar() == '=' {
				s.nextChar()
				return s.createToken(TokenTypeAndEqual, beginToken)
			}
			return s.createToken(TokenTypeAnd, beginToken)
		case '=':
			s.nextChar()
			return s.createToken(TokenTypeAmpersandEqual, beginToken)
		default:
			return s.createToken(TokenTypeAmpersand, beginToken)
		}
	case '|':
		switch s.peekChar() {
		case '|':
			s.nextChar()
			if s.peekChar() == '=' {
				s.nextChar()
				return s.createToken(TokenTypeOrEqual, beginToken)
			}
			return s.createToken(TokenTypeOr, beginToken)
		case '=':
			s.nextChar()
			return s.createToken(TokenTypeBarEqual, beginToken)
		default:
			return s.createToken(TokenTypeBar, beginToken)
		}
	case '#':
		// Shebang is not actually ever parsed here (when used correctly,
		// it's handled above in the skipComments() call) so its token is an error.
		if s.peekChar() == '!' {
			s.reportError(s.getPositionAt(s.index), "Shebang comment must be at the start of the file")
			return s.createToken(TokenTypeError, beginToken)
		}
		// Handle private identifiers.
		return s.scanIdentifierOrKeyword(beginToken, ch)
	case '0':
		return s.scanPostZero(beginToken)
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return s.scanPostDigit(beginToken)
	case '"', '\'':
		return s.scanStringLiteral(beginToken, ch)
	case '`':
		return s.scanTemplateLiteral(beginToken)
	default:
		return s.scanIdentifierOrKeyword(beginToken, ch)
	}
}

// ---- Numeric literal scanning ---------------------------------------

// scanNumberPostPeriod mirrors Scanner.scanNumberPostPeriod(int).
func (s *Scanner) scanNumberPostPeriod(beginToken int) *Token {
	s.skipDecimalDigits()
	return s.scanExponentOfNumericLiteral(beginToken)
}

// scanPostDigit mirrors Scanner.scanPostDigit(int).
func (s *Scanner) scanPostDigit(beginToken int) *Token {
	s.skipDecimalDigits()
	if s.peekChar() == 'n' {
		s.nextChar()
		return newLiteralToken(TokenTypeBigInt, s.getTokenString(beginToken), s.getTokenRangeFromOffset(beginToken))
	}
	return s.scanFractionalNumericLiteral(beginToken)
}

// scanPostZero mirrors Scanner.scanPostZero(int).
func (s *Scanner) scanPostZero(beginToken int) *Token {
	switch s.peekChar() {
	case 'b', 'B':
		// binary
		s.nextChar()
		if !isBinaryDigit(s.peekChar()) {
			s.reportError(s.getPositionAt(s.index), "Binary Integer Literal must contain at least one digit")
		}
		s.skipBinaryDigits()
		isBigInt := s.peekChar() == 'n'
		if isBigInt {
			s.nextChar()
		}
		typ := TokenTypeNumber
		if isBigInt {
			typ = TokenTypeBigInt
		}
		return newLiteralToken(typ, s.getTokenString(beginToken), s.getTokenRangeFromOffset(beginToken))
	case 'o', 'O':
		// octal
		s.nextChar()
		if !isOctalDigit(s.peekChar()) {
			s.reportError(s.getPositionAt(s.index), "Octal Integer Literal must contain at least one digit")
		}
		s.skipOctalDigits()
		if s.peekChar() == '8' || s.peekChar() == '9' {
			s.reportError(s.getPositionAt(s.index), "Invalid octal digit in octal literal.")
		}
		isBigInt := s.peekChar() == 'n'
		if isBigInt {
			s.nextChar()
		}
		typ := TokenTypeNumber
		if isBigInt {
			typ = TokenTypeBigInt
		}
		return newLiteralToken(typ, s.getTokenString(beginToken), s.getTokenRangeFromOffset(beginToken))
	case 'x', 'X':
		s.nextChar()
		if !s.peekHexDigit() {
			s.reportError(s.getPositionAt(s.index), "Hex Integer Literal must contain at least one digit")
		}
		s.skipHexDigits()
		isBigInt := s.peekChar() == 'n'
		if isBigInt {
			s.nextChar()
		}
		typ := TokenTypeNumber
		if isBigInt {
			typ = TokenTypeBigInt
		}
		return newLiteralToken(typ, s.getTokenString(beginToken), s.getTokenRangeFromOffset(beginToken))
	case 'e', 'E':
		return s.scanExponentOfNumericLiteral(beginToken)
	case '.':
		return s.scanFractionalNumericLiteral(beginToken)
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		s.skipDecimalDigits()
		if s.peekChar() == '.' {
			s.nextChar()
			s.skipDecimalDigits()
		}
		if s.peekChar() == 'n' {
			s.reportError(s.getPositionAt(s.index), "SyntaxError: nonzero BigInt can't have leading zero")
		}
		return newLiteralToken(TokenTypeNumber, s.getTokenString(beginToken), s.getTokenRangeFromOffset(beginToken))
	case 'n':
		s.nextChar()
		return newLiteralToken(TokenTypeBigInt, s.getTokenString(beginToken), s.getTokenRangeFromOffset(beginToken))
	default:
		return newLiteralToken(TokenTypeNumber, s.getTokenString(beginToken), s.getTokenRangeFromOffset(beginToken))
	}
}

// createToken mirrors Scanner.createToken(TokenType, int).
func (s *Scanner) createToken(typ TokenType, beginToken int) *Token {
	return &Token{Type: typ, Location: s.getTokenRangeFromOffset(beginToken)}
}

// ---- Identifier / keyword scanning ---------------------------------

// scanIdentifierOrKeyword mirrors Scanner.scanIdentifierOrKeyword(int, char).
func (s *Scanner) scanIdentifierOrKeyword(beginToken int, ch rune) *Token {
	valueStartIndex := s.index - 1

	containsUnicodeEscape := ch == '\\'
	bracedUnicodeEscape := false
	isPrivateIdentifier := ch == '#'
	unicodeEscapeLen := 0
	if containsUnicodeEscape {
		unicodeEscapeLen = 1
	}

	ch = s.peekChar()
	for identifiersIsIdentifierPart(ch) ||
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
		// Exit Unicode escape.
		if ch == '}' || (unicodeEscapeLen >= 6 && !bracedUnicodeEscape) {
			bracedUnicodeEscape = false
			unicodeEscapeLen = 0
		}

		// Add character to token.
		s.nextChar()
		ch = s.peekChar()
	}

	value := string(s.contents[valueStartIndex:s.index])

	if isPrivateIdentifier && value == "#" {
		s.reportError(s.getPositionAt(beginToken), "Invalid usage of #")
		return s.createToken(TokenTypeError, beginToken)
	}

	// Process unicode escapes.
	if containsUnicodeEscape {
		processed := processUnicodeEscapes(value)
		if processed == nil {
			s.reportError(s.getPositionAt(s.index), "Invalid escape sequence")
			return s.createToken(TokenTypeError, beginToken)
		}
		value = *processed
	}

	// Check first character is a valid identifier start.
	start := rune(value[0])
	if isPrivateIdentifier {
		// Skip the leading # for name validation.
		start = rune(value[1])
	}
	if !identifiersIsIdentifierStart(start) {
		s.reportError(s.getPositionAt(beginToken),
			"Character '%c' (U+%04X) is not a valid identifier start char", start, int(start))
		return s.createToken(TokenTypeError, beginToken)
	}

	if typ, ok := getKeywordType(value); ok {
		return &Token{Type: typ, Location: s.getTokenRangeFromOffset(beginToken)}
	}

	return newIdentifierToken(s.getTokenRangeFromOffset(beginToken), value)
}

// processUnicodeEscapes mirrors Scanner.processUnicodeEscapes(String).
// Returns nil if there is an invalid escape sequence.
func processUnicodeEscapes(value string) *string {
	if !strings.ContainsRune(value, '\\') {
		return &value
	}
	rs := []rune(value)
	for {
		// Find next backslash.
		escapeStart := -1
		for i, c := range rs {
			if c == '\\' {
				escapeStart = i
				break
			}
		}
		if escapeStart < 0 {
			s := string(rs)
			return &s
		}
		if escapeStart+1 >= len(rs) || rs[escapeStart+1] != 'u' {
			return nil
		}
		var hexDigits []rune
		var escapeEnd int
		if escapeStart+2 >= len(rs) {
			return nil
		}
		if rs[escapeStart+2] != '{' {
			// Simple escape: \\uXXXX
			escapeEnd = escapeStart + 6
			if escapeEnd > len(rs) {
				return nil
			}
			hexDigits = rs[escapeStart+2 : escapeEnd]
		} else {
			// Braced escape: \\u{XXXXXXX}
			escapeEnd = escapeStart + 3
			for escapeEnd < len(rs) && isHexDigit(rs[escapeEnd]) {
				escapeEnd++
			}
			if escapeEnd >= len(rs) || rs[escapeEnd] != '}' {
				return nil
			}
			hexDigits = rs[escapeStart+3 : escapeEnd]
			escapeEnd++
		}
		// Parse hex value.
		v := 0
		for _, hd := range hexDigits {
			if !isHexDigit(hd) {
				return nil
			}
			v = v*16 + hexNib(hd)
		}
		ch := rune(v)
		if !identifiersIsIdentifierPart(ch) {
			return nil
		}
		// Replace escape with the character.
		newRs := make([]rune, 0, len(rs)-len(hexDigits)+1)
		newRs = append(newRs, rs[:escapeStart]...)
		newRs = append(newRs, ch)
		newRs = append(newRs, rs[escapeEnd:]...)
		rs = newRs
	}
}

// ---- String literal scanning ----------------------------------------

// scanStringLiteral mirrors Scanner.scanStringLiteral(int, char).
func (s *Scanner) scanStringLiteral(beginIndex int, terminator rune) *Token {
	// String literals might span multiple lines.
	startingPosition := s.getPositionAt(beginIndex)

	hasUnescapedUnicodeLineOrParagraphSeparator := false
	for s.peekStringLiteralChar(terminator) {
		c := s.peekChar()
		if c == ' ' || c == ' ' {
			hasUnescapedUnicodeLineOrParagraphSeparator = true
		}
		if !s.skipStringLiteralChar() {
			return newStringLiteralToken(s.getTokenString(beginIndex), s.getTokenRangeFromPos(startingPosition), hasUnescapedUnicodeLineOrParagraphSeparator)
		}
	}
	if s.peekChar() != terminator {
		s.reportError(startingPosition, "Unterminated string literal")
	} else {
		s.nextChar()
	}
	return newStringLiteralToken(s.getTokenString(beginIndex), s.getTokenRangeFromPos(startingPosition), hasUnescapedUnicodeLineOrParagraphSeparator)
}

// scanTemplateLiteral mirrors Scanner.scanTemplateLiteral(int).
func (s *Scanner) scanTemplateLiteral(beginIndex int) *Token {
	if s.isAtEnd() {
		s.reportError(s.getPositionAt(beginIndex), "Unterminated template literal")
	}
	return s.nextTemplateLiteralTokenShared(TokenTypeNoSubstitutionTemplate, TokenTypeTemplateHead)
}

// nextTemplateLiteralTokenShared mirrors Scanner.nextTemplateLiteralTokenShared(TokenType, TokenType).
func (s *Scanner) nextTemplateLiteralTokenShared(endType, middleType TokenType) *TemplateLiteralToken {
	beginIndex := s.index
	// Save the starting position for the multi-line safe version.
	startingPosition := s.getPositionAt(beginIndex)
	result := s.skipTemplateCharacters()
	if s.isAtEnd() {
		s.reportError(startingPosition, "Unterminated template literal")
	}

	value := s.getTokenString(beginIndex)
	switch s.peekChar() {
	case '`':
		s.nextChar()
		return newTemplateLiteralToken(endType, value, result.errorMessage, result.errorLevel, result.position, s.getTokenRangeFromPos(startingPosition))
	case '$':
		s.nextChar() // $
		s.nextChar() // {
		return newTemplateLiteralToken(middleType, value, result.errorMessage, result.errorLevel, result.position, s.getTokenRangeFromPos(startingPosition))
	default:
		// Should have reported error already.
		return newTemplateLiteralToken(endType, value, result.errorMessage, result.errorLevel, result.position, s.getTokenRangeFromPos(startingPosition))
	}
}

// getTokenString mirrors Scanner.getTokenString(int).
func (s *Scanner) getTokenString(beginIndex int) string {
	return string(s.contents[beginIndex:s.index])
}

// peekStringLiteralChar mirrors Scanner.peekStringLiteralChar(char).
func (s *Scanner) peekStringLiteralChar(terminator rune) bool {
	return !s.isAtEnd() && s.peekChar() != terminator && !isStringLineTerminator(s.peekChar())
}

// skipStringLiteralChar mirrors Scanner.skipStringLiteralChar().
func (s *Scanner) skipStringLiteralChar() bool {
	if s.peekChar() == '\\' {
		return s.skipStringLiteralEscapeSequence()
	}
	s.nextChar()
	return true
}

// skipTemplateCharactersResult mirrors Scanner's inner SkipTemplateCharactersResult class.
type skipTemplateCharactersResult struct {
	errorMessage string
	errorLevel   TemplateLiteralErrorLevel
	position     SourcePosition
}

func (r *skipTemplateCharactersResult) hasError() bool {
	return r.errorMessage != ""
}

// skipTemplateCharacters mirrors Scanner.skipTemplateCharacters().
func (s *Scanner) skipTemplateCharacters() *skipTemplateCharactersResult {
	result := &skipTemplateCharactersResult{position: s.getPositionAt(s.index)}
	for !s.isAtEnd() {
		switch s.peekChar() {
		case '`':
			return result
		case '\\':
			// There might be multiple errors. Take the first one but continue scanning.
			newErr := s.skipTemplateLiteralEscapeSequence()
			if newErr != nil && !result.hasError() {
				result = newErr
			}
		case '$':
			if s.peekCharAt(1) == '{' {
				return result
			}
			s.nextChar()
		default:
			s.nextChar()
		}
	}
	return result
}

// skipTemplateLiteralEscapeSequence mirrors Scanner.skipTemplateLiteralEscapeSequence().
func (s *Scanner) skipTemplateLiteralEscapeSequence() *skipTemplateCharactersResult {
	s.nextChar() // consume '\'
	if s.isAtEnd() {
		s.reportError(s.getPositionAt(s.index), "Unterminated template literal escape sequence")
		return nil
	}
	if isLineTerminator(s.peekChar()) {
		s.skipLineTerminator()
		return nil
	}
	next := s.nextChar()
	switch next {
	case '0':
		if isDecimalDigit(s.peekChar()) {
			return s.createSkipTemplateCharactersResult("Invalid escape sequence", TemplateLiteralError)
		}
		return nil
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return s.createSkipTemplateCharactersResult("Invalid escape sequence", TemplateLiteralError)
	case 'x':
		doubleHexDigit := s.skipHexDigit() && s.skipHexDigit()
		if !doubleHexDigit {
			return s.createSkipTemplateCharactersResult("Hex digit expected", TemplateLiteralError)
		}
		return nil
	case 'u':
		if s.peekChar() == '{' {
			s.nextChar()
			if s.peekChar() == '}' {
				return s.createSkipTemplateCharactersResult("Empty unicode escape", TemplateLiteralError)
			}
			allHexDigits := true
			for s.peekChar() != '}' && allHexDigits {
				allHexDigits = allHexDigits && s.skipHexDigit()
			}
			if !allHexDigits {
				return s.createSkipTemplateCharactersResult("Hex digit expected", TemplateLiteralError)
			}
			s.nextChar() // '}'
			return nil
		}
		// \\uXXXX
		quadHexDigit := s.skipHexDigit() && s.skipHexDigit() && s.skipHexDigit() && s.skipHexDigit()
		if !quadHexDigit {
			return s.createSkipTemplateCharactersResult("Hex digit expected", TemplateLiteralError)
		}
		return nil
	case '\\', 'b', 'f', 'n', 'r', 't', 'v', '$', '`':
		return nil
	case '\'':
		return s.createSkipTemplateCharactersResult(
			fmt.Sprintf("Unnecessary escape: \"\\%c\" is equivalent to just \"%c\"", next, next),
			TemplateLiteralWarning)
	default:
		return s.createSkipTemplateCharactersResult(
			fmt.Sprintf("Unnecessary escape: '\\%c' is equivalent to just '%c'", next, next),
			TemplateLiteralWarning)
	}
}

// skipStringLiteralEscapeSequence mirrors Scanner.skipStringLiteralEscapeSequence().
func (s *Scanner) skipStringLiteralEscapeSequence() bool {
	s.nextChar() // consume '\'
	if s.isAtEnd() {
		s.reportError(s.getPositionAt(s.index), "Unterminated string literal escape sequence")
		return false
	}
	if isStringLineTerminator(s.peekChar()) {
		s.skipLineTerminator()
		return true
	}
	next := s.nextChar()
	switch next {
	case '\'', '"', '`', '\\', 'b', 'f', 'n', 'r', 't', 'v', '0':
		return true
	case 'x':
		doubleHexDigit := s.skipHexDigit() && s.skipHexDigit()
		if !doubleHexDigit {
			s.reportError(s.getPositionAt(s.index), "Hex digit expected")
		}
		return doubleHexDigit
	case 'u':
		if s.peekChar() == '{' {
			s.nextChar()
			if s.peekChar() == '}' {
				s.reportError(s.getPositionAt(s.index), "Empty unicode escape")
				return false
			}
			allHexDigits := true
			for s.peekChar() != '}' && allHexDigits {
				allHexDigits = allHexDigits && s.skipHexDigit()
			}
			if !allHexDigits {
				s.reportError(s.getPositionAt(s.index), "Hex digit expected")
			}
			s.nextChar() // '}'
			return allHexDigits
		}
		quadHexDigit := s.skipHexDigit() && s.skipHexDigit() && s.skipHexDigit() && s.skipHexDigit()
		if !quadHexDigit {
			s.reportError(s.getPositionAt(s.index), "Hex digit expected")
		}
		return quadHexDigit
	}
	if next == '/' {
		// Don't warn for '\/' (for now) since it's common in "<\/script>"
	} else {
		s.reportWarningAt(s.getPositionAt(s.index), fmt.Sprintf("Unnecessary escape: '\\%c' is equivalent to just '%c'", next, next))
	}
	return true
}

// skipHexDigit mirrors Scanner.skipHexDigit().
func (s *Scanner) skipHexDigit() bool {
	if !s.peekHexDigit() {
		return false
	}
	s.nextChar()
	return true
}

// skipLineTerminator mirrors Scanner.skipLineTerminator().
func (s *Scanner) skipLineTerminator() {
	first := s.nextChar()
	if first == '\r' && s.peekChar() == '\n' {
		s.nextChar()
	}
}

// ---- Fractional / exponent numeric literals -------------------------

// scanFractionalNumericLiteral mirrors Scanner.scanFractionalNumericLiteral(int).
func (s *Scanner) scanFractionalNumericLiteral(beginToken int) *LiteralToken {
	if s.peekChar() == '.' {
		s.nextChar()
		s.skipDecimalDigits()
	}
	return s.scanExponentOfNumericLiteral(beginToken)
}

// scanExponentOfNumericLiteral mirrors Scanner.scanExponentOfNumericLiteral(int).
func (s *Scanner) scanExponentOfNumericLiteral(beginToken int) *LiteralToken {
	switch s.peekChar() {
	case 'e', 'E':
		s.nextChar()
		switch s.peekChar() {
		case '+', '-':
			s.nextChar()
		}
		if !isDecimalDigit(s.peekChar()) {
			s.reportError(s.getPositionAt(s.index), "Exponent part must contain at least one digit")
		}
		s.skipDecimalDigits()
	}
	return newLiteralToken(TokenTypeNumber, s.getTokenString(beginToken), s.getTokenRangeFromOffset(beginToken))
}

// skipDecimalDigits mirrors Scanner.skipDecimalDigits().
func (s *Scanner) skipDecimalDigits() {
	ch := s.peekChar()
	for isDecimalDigit(ch) || ch == '_' {
		s.nextChar()
		if ch == '_' {
			if isDecimalDigit(s.peekChar()) {
				s.nextChar()
			} else {
				s.reportError(s.getPositionAt(s.index), "Trailing numeric separator")
			}
		}
		ch = s.peekChar()
	}
}

// isDecimalDigit mirrors Scanner.isDecimalDigit(char).
func isDecimalDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

// peekHexDigit mirrors Scanner.peekHexDigit().
func (s *Scanner) peekHexDigit() bool {
	return isHexDigit(s.peekChar())
}

// isHexDigit mirrors Scanner.isHexDigit(char).
func isHexDigit(ch rune) bool {
	return (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}

// skipHexDigits mirrors Scanner.skipHexDigits().
func (s *Scanner) skipHexDigits() {
	ch := s.peekChar()
	for isHexDigit(ch) || ch == '_' {
		s.nextChar()
		if ch == '_' {
			if s.peekHexDigit() {
				s.nextChar()
			} else {
				s.reportError(s.getPositionAt(s.index), "Trailing numeric separator")
			}
		}
		ch = s.peekChar()
	}
}

// skipOctalDigits mirrors Scanner.skipOctalDigits().
func (s *Scanner) skipOctalDigits() {
	ch := s.peekChar()
	for isOctalDigit(ch) || ch == '_' {
		s.nextChar()
		if ch == '_' {
			if isOctalDigit(s.peekChar()) {
				s.nextChar()
			} else {
				s.reportError(s.getPositionAt(s.index), "Trailing numeric separator")
			}
		}
		ch = s.peekChar()
	}
}

// isOctalDigit mirrors Scanner.isOctalDigit(char).
func isOctalDigit(ch rune) bool {
	return ch >= '0' && ch <= '7'
}

// skipBinaryDigits mirrors Scanner.skipBinaryDigits().
func (s *Scanner) skipBinaryDigits() {
	ch := s.peekChar()
	for isBinaryDigit(ch) || ch == '_' {
		s.nextChar()
		if ch == '_' {
			if isBinaryDigit(s.peekChar()) {
				s.nextChar()
			} else {
				s.reportError(s.getPositionAt(s.index), "Trailing numeric separator")
			}
		}
		ch = s.peekChar()
	}
}

// isBinaryDigit mirrors Scanner.isBinaryDigit(char).
func isBinaryDigit(ch rune) bool {
	return ch == '0' || ch == '1'
}

// ---- Low-level char operations (mirrors Java's nextChar/peekChar/etc.) ---

// nextChar mirrors Scanner.nextChar().
// Updates line tracking.
func (s *Scanner) nextChar() rune {
	if s.isAtEnd() {
		return 0
	}
	c := s.contents[s.index]
	s.index++
	if c == '\n' {
		s.line++
		s.lineBOL = s.index
	}
	return c
}

// peekChar mirrors Scanner.peekChar().
func (s *Scanner) peekChar() rune {
	return s.peekCharAt(0)
}

// peekCharAt mirrors Scanner.peekChar(int offset).
func (s *Scanner) peekCharAt(offset int) rune {
	if !s.isValidIndex(s.index + offset) {
		return 0
	}
	return s.contents[s.index+offset]
}

// ---- Error reporting ------------------------------------------------

// reportError mirrors Scanner.reportError(String, Object...).
func (s *Scanner) reportError(position SourcePosition, format string, args ...interface{}) {
	if s.errorReporter != nil {
		s.errorReporter.ReportError(position, format, args...)
	}
}

// reportWarningAt mirrors Scanner.reportWarning(String, Object...).
func (s *Scanner) reportWarningAt(position SourcePosition, msg string) {
	if s.errorReporter != nil {
		s.errorReporter.ReportWarning(position, "%s", msg)
	}
}

// ---- Position helpers -----------------------------------------------

// getPositionAt computes the SourcePosition for rune index offset.
// This is a simplified version that tracks line/column from current state.
// For accuracy it uses the line counter maintained by nextChar.
func (s *Scanner) getPositionAt(offset int) SourcePosition {
	// Compute line/col by scanning from lineBOL.
	// Since we maintain s.line and s.lineBOL incrementally, we compute
	// column as offset - lineBOL (both are rune indices).
	col := offset - s.lineBOL
	if col < 0 {
		col = 0
	}
	return SourcePosition{Offset: offset, Line: s.line - 1, Column: col}
}

// getTokenRangeFromOffset mirrors Scanner.getTokenRange(int startOffset).
func (s *Scanner) getTokenRangeFromOffset(startOffset int) SourceRange {
	return SourceRange{
		Start: s.getPositionAt(startOffset),
		End:   s.getPositionAt(s.index),
	}
}

// getTokenRangeFromPos mirrors Scanner.getTokenRange(SourcePosition).
func (s *Scanner) getTokenRangeFromPos(start SourcePosition) SourceRange {
	return SourceRange{
		Start: start,
		End:   s.getPositionAt(s.index),
	}
}

// rewindLineTo resets line-tracking state to match a previously saved position.
func (s *Scanner) rewindLineTo(position SourcePosition) {
	// Recompute line and lineBOL from position.line (0-based in SourcePosition).
	s.line = position.Line + 1
	s.lineBOL = position.Offset - position.Column
}

// createSkipTemplateCharactersResult mirrors Scanner.createSkipTemplateCharactersResult.
func (s *Scanner) createSkipTemplateCharactersResult(message string, level TemplateLiteralErrorLevel) *skipTemplateCharactersResult {
	return &skipTemplateCharactersResult{
		errorMessage: message,
		errorLevel:   level,
		position:     s.getPositionAt(s.index),
	}
}

// ---- Regular expression scanning ------------------------------------

// skipRegularExpressionBody mirrors Scanner.skipRegularExpressionBody().
func (s *Scanner) skipRegularExpressionBody() bool {
	if !isRegularExpressionFirstChar(s.peekChar()) {
		s.reportError(s.getPositionAt(s.index), "Expected regular expression first char")
		return false
	}
	if !s.skipRegularExpressionChar() {
		return false
	}
	for !s.isAtEnd() && isRegularExpressionChar(s.peekChar()) {
		if !s.skipRegularExpressionChar() {
			return false
		}
	}
	return true
}

// skipRegularExpressionChar mirrors Scanner.skipRegularExpressionChar().
func (s *Scanner) skipRegularExpressionChar() bool {
	switch s.peekChar() {
	case '\\':
		return s.skipRegularExpressionBackslashSequence()
	case '[':
		return s.skipRegularExpressionClass()
	default:
		s.nextChar()
		return true
	}
}

// skipRegularExpressionBackslashSequence mirrors Scanner.skipRegularExpressionBackslashSequence().
func (s *Scanner) skipRegularExpressionBackslashSequence() bool {
	s.nextChar() // '\'
	if isLineTerminator(s.peekChar()) {
		s.reportError(s.getPositionAt(s.index), "New line not allowed in regular expression literal")
		return false
	}
	s.nextChar()
	return true
}

// skipRegularExpressionClass mirrors Scanner.skipRegularExpressionClass().
func (s *Scanner) skipRegularExpressionClass() bool {
	s.nextChar() // '['
	for !s.isAtEnd() && s.peekRegularExpressionClassChar() {
		if !s.skipRegularExpressionClassChar() {
			return false
		}
	}
	if s.peekChar() != ']' {
		s.reportError(s.getPositionAt(s.index), "']' expected")
		return false
	}
	s.nextChar()
	return true
}

// peekRegularExpressionClassChar mirrors Scanner.peekRegularExpressionClassChar().
func (s *Scanner) peekRegularExpressionClassChar() bool {
	return s.peekChar() != ']' && !isLineTerminator(s.peekChar())
}

// skipRegularExpressionClassChar mirrors Scanner.skipRegularExpressionClassChar().
func (s *Scanner) skipRegularExpressionClassChar() bool {
	if s.peekChar() == '\\' {
		return s.skipRegularExpressionBackslashSequence()
	}
	s.nextChar()
	return true
}

// isRegularExpressionFirstChar mirrors Scanner.isRegularExpressionFirstChar(char).
func isRegularExpressionFirstChar(ch rune) bool {
	return isRegularExpressionChar(ch) && ch != '*'
}

// isRegularExpressionChar mirrors Scanner.isRegularExpressionChar(char).
func isRegularExpressionChar(ch rune) bool {
	switch ch {
	case '/':
		return false
	case '\\', '[':
		return true
	default:
		return !isLineTerminator(ch)
	}
}

// ---- Identifiers (mirrors com.google.javascript.jscomp.parsing.parser.Identifiers) ---

// identifiersIsIdentifierStart mirrors Identifiers.isIdentifierStart(char).
func identifiersIsIdentifierStart(ch rune) bool {
	if ch <= 127 {
		return (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || ch == '_' || ch == '$'
	}
	return unicode.IsLetter(ch)
}

// identifiersIsIdentifierPart mirrors Identifiers.isIdentifierPart(char).
func identifiersIsIdentifierPart(ch rune) bool {
	if ch <= 127 {
		return (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') ||
			(ch >= '0' && ch <= '9') || ch == '_' || ch == '$'
	}
	return identifiersIsIdentifierStart(ch) || unicode.IsDigit(ch)
}

// ---- Utility --------------------------------------------------------

// hexNib returns the hex value of a single hex digit rune.
func hexNib(c rune) int {
	if c >= '0' && c <= '9' {
		return int(c - '0')
	}
	if c >= 'a' && c <= 'f' {
		return int(c-'a') + 10
	}
	return int(c-'A') + 10
}

// =====================================================================
// Char-level API — preserved for the existing char-level parser (parser.go).
// These are 1:1 wrappers/aliases over the private scanner primitives so
// the existing parser continues to compile without modification.
// The parser rewrite (#141) will remove these callers and use the token API.
// =====================================================================

// ScannerMark is the checkpoint type for Save/Restore speculative parsing.
type ScannerMark struct {
	pos, line, lineBOL, lastTokenLine int
	currentTokens                     []*Token
}

// Save returns a checkpoint that can be restored with Restore.
func (s *Scanner) Save() ScannerMark {
	// Copy current tokens slice.
	toks := make([]*Token, len(s.currentTokens))
	copy(toks, s.currentTokens)
	return ScannerMark{
		pos:           s.index,
		line:          s.line,
		lineBOL:       s.lineBOL,
		lastTokenLine: s.lastTokenLine,
		currentTokens: toks,
	}
}

// Restore rewinds the scanner to a previously saved mark.
func (s *Scanner) Restore(m ScannerMark) {
	s.index = m.pos
	s.line = m.line
	s.lineBOL = m.lineBOL
	s.lastTokenLine = m.lastTokenLine
	s.currentTokens = m.currentTokens
}

// Pos returns the current rune index. Mirrors Java's getOffset().
func (s *Scanner) Pos() int { return s.index }

// Line returns the current 1-based line number.
func (s *Scanner) Line() int { return s.line }

// Col returns the current 0-based column (runes from BOL).
func (s *Scanner) Col() int {
	col := s.index - s.lineBOL
	if col < 0 {
		return 0
	}
	return col
}

// AtEnd returns true if the scanner has consumed all input.
func (s *Scanner) AtEnd() bool { return s.isAtEnd() }

// Peek returns the current rune without consuming. Returns 0 at EOF.
func (s *Scanner) Peek() rune { return s.peekChar() }

// PeekAt returns the rune at offset ahead of the cursor, or 0 past end.
func (s *Scanner) PeekAt(offset int) rune { return s.peekCharAt(offset) }

// Advance consumes one rune and updates line/col tracking.
func (s *Scanner) Advance() rune {
	ch := s.nextChar()
	return ch
}

// LastTokenLine returns the line of the most recently consumed non-trivia character.
func (s *Scanner) LastTokenLine() int { return s.lastTokenLine }

// MarkTokenConsumed updates lastTokenLine to the current scanner line.
func (s *Scanner) MarkTokenConsumed() { s.lastTokenLine = s.line }

// Text returns the substring between two rune-positions.
func (s *Scanner) Text(start, end int) string {
	if start < 0 {
		start = 0
	}
	if end > s.contentsLength {
		end = s.contentsLength
	}
	if start >= end {
		return ""
	}
	return string(s.contents[start:end])
}

// ConsumeByte advances past `b` if the current rune matches.
func (s *Scanner) ConsumeByte(b byte) bool {
	if s.isAtEnd() || byte(s.peekChar()) != b {
		return false
	}
	s.Advance()
	s.lastTokenLine = s.line
	return true
}

// ConsumeIdent reads a JS identifier starting at the cursor.
// Returns ("", false) if the cursor isn't on an ident start or if the
// identifier matches a reserved word.
func (s *Scanner) ConsumeIdent() (string, bool) {
	if s.isAtEnd() {
		return "", false
	}
	c := s.peekChar()
	if !isIdentStart(c) {
		return "", false
	}
	start := s.index
	end := s.index
	for end < s.contentsLength && isIdentCont(s.contents[end]) {
		end++
	}
	ident := string(s.contents[start:end])
	if isReservedWord(ident) {
		return "", false
	}
	for s.index < end {
		s.Advance()
	}
	s.lastTokenLine = s.line
	return ident, true
}

// isReservedWord returns true for JS reserved words.
func isReservedWord(str string) bool {
	switch str {
	case "break", "case", "catch", "class", "const", "continue",
		"debugger", "default", "do", "else", "export",
		"extends", "finally", "for", "function", "if", "import",
		"return", "switch", "throw", "try", "var", "while",
		"let":
		return true
	case "true", "false", "null", "undefined", "this":
		return true
	case "new", "typeof", "void", "delete", "in", "instanceof":
		return true
	}
	return false
}

// AtKeyword reports whether `kw` is the next keyword (non-consuming).
func (s *Scanner) AtKeyword(kw string) bool {
	rs := []rune(kw)
	if s.index+len(rs) > s.contentsLength {
		return false
	}
	for i, r := range rs {
		if s.contents[s.index+i] != r {
			return false
		}
	}
	if s.index+len(rs) < s.contentsLength && isIdentCont(s.contents[s.index+len(rs)]) {
		return false
	}
	return true
}

// ConsumeKeywordRaw advances past a keyword without the reserved-word filter.
func (s *Scanner) ConsumeKeywordRaw(kw string) {
	for range []rune(kw) {
		s.Advance()
	}
	s.lastTokenLine = s.line
}

// ConsumeKeyword matches an exact keyword followed by a non-ident boundary.
func (s *Scanner) ConsumeKeyword(kw string) bool {
	rs := []rune(kw)
	if s.index+len(rs) > s.contentsLength {
		return false
	}
	for i, r := range rs {
		if s.contents[s.index+i] != r {
			return false
		}
	}
	if s.index+len(rs) < s.contentsLength {
		nxt := s.contents[s.index+len(rs)]
		if isIdentCont(nxt) {
			return false
		}
	}
	for range rs {
		s.Advance()
	}
	s.lastTokenLine = s.line
	return true
}

// ConsumeBlockComment reads a /* … */ block comment.
func (s *Scanner) ConsumeBlockComment() (string, bool) {
	if s.index+1 >= s.contentsLength || s.peekChar() != '/' || s.peekCharAt(1) != '*' {
		return "", false
	}
	s.Advance() // '/'
	s.Advance() // '*'
	start := s.index
	for s.index+1 < s.contentsLength && !(s.peekChar() == '*' && s.peekCharAt(1) == '/') {
		s.Advance()
	}
	if s.index+1 >= s.contentsLength {
		return "", false
	}
	body := string(s.contents[start:s.index])
	s.Advance() // '*'
	s.Advance() // '/'
	s.lastTokenLine = s.line
	return body, true
}

// ConsumeLineComment reads a // … comment.
func (s *Scanner) ConsumeLineComment() bool {
	if s.index+1 < s.contentsLength && s.peekChar() == '/' && s.peekCharAt(1) == '/' {
		for !s.isAtEnd() && s.peekChar() != '\n' {
			s.Advance()
		}
		return true
	}
	return false
}

// SkipTrivia consumes whitespace plus plain (non-JSDoc) block + line comments.
func (s *Scanner) SkipTrivia() {
	s.skipTriviaTracked()
}

// SkipTriviaSeenNewline is SkipTrivia but reports whether the skipped region crossed a line boundary.
func (s *Scanner) SkipTriviaSeenNewline() bool {
	return s.skipTriviaTracked()
}

func (s *Scanner) skipTriviaTracked() bool {
	startLine := s.line
	for !s.isAtEnd() {
		c := s.peekChar()
		if unicode.IsSpace(c) {
			s.Advance()
			continue
		}
		if s.ConsumeLineComment() {
			continue
		}
		if s.index+2 < s.contentsLength && c == '/' && s.peekCharAt(1) == '*' && s.peekCharAt(2) != '*' {
			s.Advance()
			s.Advance()
			for s.index+1 < s.contentsLength && !(s.peekChar() == '*' && s.peekCharAt(1) == '/') {
				s.Advance()
			}
			if s.index+1 < s.contentsLength {
				s.Advance()
				s.Advance()
			}
			continue
		}
		break
	}
	return s.line != startLine
}

// SkipToStatementBoundary walks forward until just past the next `;` or `}` at top-level, or to EOF.
func (s *Scanner) SkipToStatementBoundary() {
	depth := 0
	for !s.isAtEnd() {
		c := s.peekChar()
		switch c {
		case ';':
			s.Advance()
			if depth == 0 {
				return
			}
		case '{':
			depth++
			s.Advance()
		case '}':
			if depth == 0 {
				return
			}
			depth--
			s.Advance()
			if depth == 0 {
				return
			}
		case '/':
			if s.index+1 < s.contentsLength && s.peekCharAt(1) == '/' {
				for !s.isAtEnd() && s.peekChar() != '\n' {
					s.Advance()
				}
			} else if s.index+1 < s.contentsLength && s.peekCharAt(1) == '*' {
				s.Advance()
				s.Advance()
				for s.index+1 < s.contentsLength && !(s.peekChar() == '*' && s.peekCharAt(1) == '/') {
					s.Advance()
				}
				if s.index+1 < s.contentsLength {
					s.Advance()
					s.Advance()
				}
			} else {
				s.Advance()
			}
		case '"', '\'':
			s.Advance()
			for !s.isAtEnd() && s.peekChar() != c {
				if s.peekChar() == '\\' && s.index+1 < s.contentsLength {
					s.Advance()
					s.Advance()
					continue
				}
				if s.peekChar() == '\n' {
					break
				}
				s.Advance()
			}
			if !s.isAtEnd() {
				s.Advance()
			}
		default:
			s.Advance()
		}
	}
}

// ---- Template literal char-level API (used by parser.go) ------------

// TemplatePartKind mirrors the old char-level template part kind.
type TemplatePartKind int

const (
	TemplateNoSub  TemplatePartKind = iota
	TemplateHead
	TemplateMiddle
	TemplateTail
)

// TemplatePart carries the raw segment text and token kind.
type TemplatePart struct {
	Kind   TemplatePartKind
	Cooked string
}

// ScanTemplateLiteral is called when the scanner is positioned just after the opening backtick.
func (s *Scanner) ScanTemplateLiteral() (TemplatePart, bool) {
	return s.scanTemplateCharacters(TemplateNoSub, TemplateHead)
}

// ScanTemplateMiddleOrTail is called after the parser has consumed the closing `}`.
func (s *Scanner) ScanTemplateMiddleOrTail() (TemplatePart, bool) {
	return s.scanTemplateCharacters(TemplateTail, TemplateMiddle)
}

// scanTemplateCharacters is the char-level template scanner used by parser.go.
func (s *Scanner) scanTemplateCharacters(endKind, middleKind TemplatePartKind) (TemplatePart, bool) {
	start := s.index
	for !s.isAtEnd() {
		c := s.peekChar()
		switch c {
		case '`':
			cooked := string(s.contents[start:s.index])
			s.Advance()
			s.lastTokenLine = s.line
			return TemplatePart{Kind: endKind, Cooked: cooked}, true
		case '\\':
			s.Advance()
			if !s.isAtEnd() {
				s.Advance()
			}
		case '$':
			if s.index+1 < s.contentsLength && s.peekCharAt(1) == '{' {
				cooked := string(s.contents[start:s.index])
				s.Advance() // '$'
				s.Advance() // '{'
				s.lastTokenLine = s.line
				return TemplatePart{Kind: middleKind, Cooked: cooked}, true
			}
			s.Advance()
		default:
			s.Advance()
		}
	}
	return TemplatePart{}, false
}

// ---- Numeric literal char-level API (used by parser.go) -------------

// ScanNumericLiteral mirrors the old char-level numeric scanner.
// Cursor must be on the first digit (or '.' for .5 forms).
// Returns (raw, isFloat, ok); advances the scanner on success.
func (s *Scanner) ScanNumericLiteral() (string, bool, bool) {
	start := s.index
	c := s.Peek()
	if c == '0' {
		s.Advance()
		switch p := s.Peek(); {
		case p == 'b' || p == 'B':
			s.Advance()
			s.skipDigitsFunc(isBinaryDigit)
			if s.Peek() == 'n' {
				s.Advance()
			}
			s.lastTokenLine = s.line
			return s.Text(start, s.index), false, true
		case p == 'o' || p == 'O':
			s.Advance()
			s.skipDigitsFunc(isOctalDigit)
			if s.Peek() == 'n' {
				s.Advance()
			}
			s.lastTokenLine = s.line
			return s.Text(start, s.index), false, true
		case p == 'x' || p == 'X':
			s.Advance()
			s.skipDigitsFunc(isHexDigit)
			if s.Peek() == 'n' {
				s.Advance()
			}
			s.lastTokenLine = s.line
			return s.Text(start, s.index), false, true
		case p == 'n':
			s.Advance()
			s.lastTokenLine = s.line
			return s.Text(start, s.index), false, true
		case p == '.':
			s.Advance()
			s.skipDigitsFunc(isDecimalDigit)
			s.scanExponentSuffix()
			s.lastTokenLine = s.line
			return s.Text(start, s.index), true, true
		case p == 'e' || p == 'E':
			s.scanExponentSuffix()
			s.lastTokenLine = s.line
			return s.Text(start, s.index), true, true
		default:
			if isDecimalDigit(p) {
				s.skipDigitsFunc(isDecimalDigit)
				if s.Peek() == '.' {
					s.Advance()
					s.skipDigitsFunc(isDecimalDigit)
					s.scanExponentSuffix()
					s.lastTokenLine = s.line
					return s.Text(start, s.index), true, true
				}
			}
			s.lastTokenLine = s.line
			return s.Text(start, s.index), false, true
		}
	}
	if c == '.' {
		s.Advance()
		s.skipDigitsFunc(isDecimalDigit)
		s.scanExponentSuffix()
		s.lastTokenLine = s.line
		return s.Text(start, s.index), true, true
	}
	s.skipDigitsFunc(isDecimalDigit)
	if s.Peek() == 'n' {
		s.Advance()
		s.lastTokenLine = s.line
		return s.Text(start, s.index), false, true
	}
	isFloat := false
	if s.Peek() == '.' {
		isFloat = true
		s.Advance()
		s.skipDigitsFunc(isDecimalDigit)
	}
	if s.scanExponentSuffix() {
		isFloat = true
	}
	s.lastTokenLine = s.line
	return s.Text(start, s.index), isFloat, true
}

// scanExponentSuffix advances past an optional e/E exponent; returns true if found.
func (s *Scanner) scanExponentSuffix() bool {
	if p := s.Peek(); p != 'e' && p != 'E' {
		return false
	}
	s.Advance()
	if p := s.Peek(); p == '+' || p == '-' {
		s.Advance()
	}
	s.skipDigitsFunc(isDecimalDigit)
	return true
}

// skipDigitsFunc advances past digit runs matching pred, including '_' separators.
func (s *Scanner) skipDigitsFunc(pred func(rune) bool) {
	for !s.isAtEnd() {
		if c := s.Peek(); pred(c) {
			s.Advance()
		} else if c == '_' {
			s.Advance()
			if pred(s.Peek()) {
				s.Advance()
			}
		} else {
			break
		}
	}
}

// ---- String body char-level API (used by parser.go) -----------------

// ScanStringBody mirrors old char-level ScanStringBody.
// Cursor must be just AFTER the opening quote.
func (s *Scanner) ScanStringBody(terminator rune) (string, bool) {
	start := s.index
	for !s.isAtEnd() && s.Peek() != terminator {
		c := s.Peek()
		if c == '\n' || c == '\r' {
			return "", false
		}
		if c == '\\' {
			s.Advance()
			if s.isAtEnd() {
				return "", false
			}
			nc := s.Peek()
			if nc == '\r' {
				s.Advance()
				if s.Peek() == '\n' {
					s.Advance()
				}
				continue
			}
			if nc == '\n' {
				s.Advance()
				continue
			}
			s.Advance()
			continue
		}
		s.Advance()
	}
	if s.isAtEnd() {
		return "", false
	}
	raw := s.Text(start, s.index)
	s.Advance()
	s.lastTokenLine = s.line
	return raw, true
}

// ---- Regular expression char-level API (used by parser.go) ----------

// ScanRegularExpressionLiteral mirrors old char-level ScanRegularExpressionLiteral.
func (s *Scanner) ScanRegularExpressionLiteral() (string, bool) {
	begin := s.index

	// leading '/'
	s.Advance()

	// body
	if !s.skipRegularExpressionBody() {
		return "", false
	}

	// separating '/'
	if s.Peek() != '/' {
		return "", false
	}
	s.Advance()

	// flags
	for isIdentCont(s.Peek()) {
		s.Advance()
	}

	s.lastTokenLine = s.line
	return string(s.contents[begin:s.index]), true
}

// ---- Shared char-class predicates ----------------------------------

// isIdentStart and isIdentCont are the char-level predicates used by
// the old ConsumeIdent / ConsumeKeyword paths (identical to identifiersIsIdentifierStart/Part
// but kept as local names for the char-level API callers).
func isIdentStart(c rune) bool {
	return unicode.IsLetter(c) || c == '_' || c == '$'
}

func isIdentCont(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsDigit(c) || c == '_' || c == '$'
}

// DecodeStringEscapes decodes the raw body of a string literal.
// Mirrors the old scanner's DecodeStringEscapes helper.
func DecodeStringEscapes(raw string) string {
	if !strings.ContainsRune(raw, '\\') {
		return raw
	}
	var b strings.Builder
	b.Grow(len(raw))
	rs := []rune(raw)
	for i := 0; i < len(rs); {
		c := rs[i]
		i++
		if c != '\\' {
			b.WriteRune(c)
			continue
		}
		if i >= len(rs) {
			b.WriteByte('\\')
			break
		}
		esc := rs[i]
		i++
		switch esc {
		case 'n':
			b.WriteByte('\n')
		case 'r':
			b.WriteByte('\r')
		case 't':
			b.WriteByte('\t')
		case 'b':
			b.WriteByte('\b')
		case 'f':
			b.WriteByte('\f')
		case 'v':
			b.WriteByte('\v')
		case '0':
			b.WriteByte(0)
		case '\\', '\'', '"', '`':
			b.WriteRune(esc)
		case '\r':
			if i < len(rs) && rs[i] == '\n' {
				i++
			}
		case '\n': // line continuation — swallow
		case 'x':
			if i+1 < len(rs) && isHexDigit(rs[i]) && isHexDigit(rs[i+1]) {
				b.WriteRune(rune(hexNib(rs[i])<<4 | hexNib(rs[i+1])))
				i += 2
			} else {
				b.WriteByte('\\')
				b.WriteRune(esc)
			}
		case 'u':
			if i < len(rs) && rs[i] == '{' {
				i++
				j := i
				for i < len(rs) && rs[i] != '}' {
					i++
				}
				if v, ok := parseHexRunes(rs[j:i]); ok {
					b.WriteRune(rune(v))
				} else {
					b.WriteString(`\u{`)
					b.WriteString(string(rs[j:i]))
					b.WriteByte('}')
				}
				if i < len(rs) {
					i++
				}
			} else if i+3 < len(rs) && isHexDigit(rs[i]) && isHexDigit(rs[i+1]) &&
				isHexDigit(rs[i+2]) && isHexDigit(rs[i+3]) {
				b.WriteRune(rune(hexNib(rs[i])<<12 | hexNib(rs[i+1])<<8 |
					hexNib(rs[i+2])<<4 | hexNib(rs[i+3])))
				i += 4
			} else {
				b.WriteByte('\\')
				b.WriteRune(esc)
			}
		default:
			b.WriteRune(esc)
		}
	}
	return b.String()
}

func parseHexRunes(rs []rune) (int, bool) {
	if len(rs) == 0 {
		return 0, false
	}
	v := 0
	for _, c := range rs {
		if !isHexDigit(c) {
			return 0, false
		}
		v = v<<4 | hexNib(c)
	}
	return v, true
}
