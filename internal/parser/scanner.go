// Port of com/google/javascript/jscomp/parsing/parser/Scanner.java
// @ closure-compiler v20260415 (subset).
//
// Java's Scanner is a full ES tokenizer that emits TokenType-bearing
// Token records (IDENTIFIER, NUMBER_LITERAL, STRING_LITERAL, OPEN_CURLY,
// etc.) from a char stream. This file carries the character-level
// scanning primitives the narrow jsgo grammar currently uses —
// SkipTrivia, ConsumeIdent, ConsumeKeyword, ConsumeBlockComment,
// position tracking — but does not yet emit a full Token type. A
// future expansion adds a TokenType enum + Peek/Consume(TokenType)
// interface so the Parser stops peeking at raw chars.

package parser

import (
	"strings"
	"unicode"
)

// Scanner walks a source-string by rune. All positions are reported
// as 1-indexed line, 0-indexed col to match Java's getLineno /
// getCharno convention.
type Scanner struct {
	src     []rune
	pos     int
	line    int
	lineBOL int
	// lastTokenLine is the line of the most recently consumed
	// non-trivia character. Used by the parser's ASI path to decide
	// whether the next token sits on a later line. Trivia-only
	// consumption (SkipTrivia, ConsumeLineComment, non-JSDoc block
	// comments) does not touch this field.
	lastTokenLine int
}

// NewScanner mirrors Java's `new Scanner(source)` constructor.
func NewScanner(src string) *Scanner {
	return &Scanner{src: []rune(src), line: 1, lastTokenLine: 1}
}

// LastTokenLine returns the line of the most recently consumed
// non-trivia character. The parser compares this to Line() to
// decide whether ES5 §7.9 ASI should fire after a statement that's
// missing its terminating `;`.
func (s *Scanner) LastTokenLine() int { return s.lastTokenLine }

// MarkTokenConsumed updates lastTokenLine to the current scanner
// line. The Parser calls this after any Advance-based token
// consumption that doesn't go through one of Scanner's Consume*
// helpers (e.g. parseNumberLiteral / parseStringLiteral loops).
func (s *Scanner) MarkTokenConsumed() { s.lastTokenLine = s.line }

// ---- basic position / peek / advance ----------------------------- //

// Text returns the substring between two byte-positions previously
// observed via Pos(). Useful for literal-value extraction
// (parseStringLiteral, parseNumberLiteral) without callers reaching
// into the underlying rune slice.
func (s *Scanner) Text(start, end int) string {
	if start < 0 {
		start = 0
	}
	if end > len(s.src) {
		end = len(s.src)
	}
	if start >= end {
		return ""
	}
	return string(s.src[start:end])
}

func (s *Scanner) Pos() int    { return s.pos }
func (s *Scanner) Line() int   { return s.line }
func (s *Scanner) Col() int    { return s.pos - s.lineBOL }
func (s *Scanner) AtEnd() bool { return s.pos >= len(s.src) }


// Peek returns the current rune without consuming. Returns 0 at EOF.
func (s *Scanner) Peek() rune {
	if s.AtEnd() {
		return 0
	}
	return s.src[s.pos]
}

// PeekAt returns the rune at offset ahead of the cursor, or 0 past end.
func (s *Scanner) PeekAt(offset int) rune {
	i := s.pos + offset
	if i < 0 || i >= len(s.src) {
		return 0
	}
	return s.src[i]
}

// Advance consumes one rune and updates line/col tracking.
func (s *Scanner) Advance() rune {
	if s.AtEnd() {
		return 0
	}
	c := s.src[s.pos]
	s.pos++
	if c == '\n' {
		s.line++
		s.lineBOL = s.pos
	}
	return c
}

// Save / Restore give the parser a checkpoint it can rewind to for
// speculative parse paths. Mirrors how Java's Scanner is coupled with
// the parser via Token-level markSavePoint / resetTo methods.
type ScannerMark struct {
	pos, line, lineBOL, lastTokenLine int
}

func (s *Scanner) Save() ScannerMark {
	return ScannerMark{s.pos, s.line, s.lineBOL, s.lastTokenLine}
}
func (s *Scanner) Restore(m ScannerMark) {
	s.pos, s.line, s.lineBOL, s.lastTokenLine = m.pos, m.line, m.lineBOL, m.lastTokenLine
}

// ---- byte-level conditionals ------------------------------------- //

// ConsumeByte advances past `b` if the current rune matches. Returns
// whether it consumed.
func (s *Scanner) ConsumeByte(b byte) bool {
	if s.AtEnd() || byte(s.src[s.pos]) != b {
		return false
	}
	s.Advance()
	s.lastTokenLine = s.line
	return true
}

// ---- identifier / keyword --------------------------------------- //

// ConsumeIdent reads a JS identifier starting at the cursor.
// Returns ("", false) if the cursor isn't on an ident start or if the
// identifier matches a reserved word (Java's Scanner keeps keyword +
// identifier disjoint via its TokenType enum; we emulate that by
// refusing to return reserved names, which prevents `if`, `return`,
// etc. from being wrongly consumed as identifiers in speculative
// statement parses).
func (s *Scanner) ConsumeIdent() (string, bool) {
	if s.AtEnd() {
		return "", false
	}
	c := s.src[s.pos]
	if !isIdentStart(c) {
		return "", false
	}
	start := s.pos
	end := s.pos
	for end < len(s.src) && isIdentCont(s.src[end]) {
		end++
	}
	ident := string(s.src[start:end])
	if isReservedWord(ident) {
		return "", false
	}
	for s.pos < end {
		s.Advance()
	}
	s.lastTokenLine = s.line
	return ident, true
}

// isReservedWord returns true for JS reserved words + the subset of
// contextual keywords we must not misread as identifiers. Keeping the
// list in sync with Parser.java's ReservedWord handling is important
// for parser correctness; anything missing here risks a bad AST (see
// `if (x) { … }` being parsed as `if(x) { … };`).
func isReservedWord(s string) bool {
	switch s {
	// Control flow / declarations the parser might speculatively
	// treat as identifiers in the LHS slot.
	case "break", "case", "catch", "class", "const", "continue",
		"debugger", "default", "do", "else", "export",
		"extends", "finally", "for", "function", "if", "import",
		"return", "switch", "throw", "try", "var", "while",
		"let":
		return true
	// Literal-shaped keywords that are parsed in specific
	// expression branches (tryPrimary consumes them before the
	// identifier fallback runs).
	case "true", "false", "null", "undefined", "this":
		return true
	// Operator-shaped keywords handled by tryUnary / peekBinOp.
	case "new", "typeof", "void", "delete", "in", "instanceof":
		return true
	}
	// Contextual keywords (`super`, `yield`, `async`, `await`,
	// `static`, `enum`, `with`) intentionally pass through as
	// identifiers — the scope / checker layers give them meaning,
	// and some corpus fixtures use them unadorned.
	return false
}

// AtKeyword reports whether `kw` is the next keyword (followed by a
// non-ident boundary) at the cursor — non-consuming. Used by the
// binary-op parser to decide whether to consume `instanceof` / `in`
// without advancing the scanner.
func (s *Scanner) AtKeyword(kw string) bool {
	rs := []rune(kw)
	if s.pos+len(rs) > len(s.src) {
		return false
	}
	for i, r := range rs {
		if s.src[s.pos+i] != r {
			return false
		}
	}
	if s.pos+len(rs) < len(s.src) && isIdentCont(s.src[s.pos+len(rs)]) {
		return false
	}
	return true
}

// ConsumeKeywordRaw advances past a keyword that the caller has
// already verified via AtKeyword. Skips the reserved-word filter that
// ConsumeIdent applies, so keywords known to the binary-op parser
// (instanceof, in) can still be consumed.
func (s *Scanner) ConsumeKeywordRaw(kw string) {
	for i := 0; i < len([]rune(kw)); i++ {
		s.Advance()
	}
	s.lastTokenLine = s.line
}

// ConsumeKeyword matches an exact keyword followed by a non-ident
// boundary. Returns whether the keyword was consumed.
func (s *Scanner) ConsumeKeyword(kw string) bool {
	rs := []rune(kw)
	if s.pos+len(rs) > len(s.src) {
		return false
	}
	for i, r := range rs {
		if s.src[s.pos+i] != r {
			return false
		}
	}
	if s.pos+len(rs) < len(s.src) {
		nxt := s.src[s.pos+len(rs)]
		if isIdentCont(nxt) {
			return false
		}
	}
	for i := 0; i < len(rs); i++ {
		s.Advance()
	}
	s.lastTokenLine = s.line
	return true
}

// ---- trivia (whitespace + comments) ----------------------------- //

// ConsumeBlockComment reads a `/* … */` block comment (including the
// `/**` JSDoc variant). Returns the comment body (the text between
// the delimiters) and whether a comment was consumed.
func (s *Scanner) ConsumeBlockComment() (string, bool) {
	if s.pos+1 >= len(s.src) || s.src[s.pos] != '/' || s.src[s.pos+1] != '*' {
		return "", false
	}
	s.Advance()
	s.Advance()
	start := s.pos
	for s.pos+1 < len(s.src) && !(s.src[s.pos] == '*' && s.src[s.pos+1] == '/') {
		s.Advance()
	}
	if s.pos+1 >= len(s.src) {
		return "", false
	}
	body := string(s.src[start:s.pos])
	s.Advance()
	s.Advance()
	s.lastTokenLine = s.line
	return body, true
}

// ConsumeLineComment reads a `// …\n` comment. Returns whether it
// consumed one.
func (s *Scanner) ConsumeLineComment() bool {
	if s.pos+1 < len(s.src) && s.src[s.pos] == '/' && s.src[s.pos+1] == '/' {
		for !s.AtEnd() && s.src[s.pos] != '\n' {
			s.Advance()
		}
		return true
	}
	return false
}

// SkipTrivia consumes whitespace plus plain (non-JSDoc) block + line
// comments. JSDoc block comments (`/**`) are left for the parser to
// attach to the following declaration.
func (s *Scanner) SkipTrivia() {
	s.skipTriviaTracked()
}

// SkipTriviaSeenNewline is SkipTrivia but reports whether the skipped
// region crossed a line boundary. Used by the parser's ASI path —
// when a statement isn't terminated by `;` but the next token sits on
// a later line, ES5 §7.9 says to insert a semicolon.
func (s *Scanner) SkipTriviaSeenNewline() bool {
	return s.skipTriviaTracked()
}

func (s *Scanner) skipTriviaTracked() bool {
	startLine := s.line
	for !s.AtEnd() {
		c := s.src[s.pos]
		if unicode.IsSpace(c) {
			s.Advance()
			continue
		}
		if s.ConsumeLineComment() {
			continue
		}
		if s.pos+2 < len(s.src) && c == '/' && s.src[s.pos+1] == '*' && s.src[s.pos+2] != '*' {
			s.Advance()
			s.Advance()
			for s.pos+1 < len(s.src) && !(s.src[s.pos] == '*' && s.src[s.pos+1] == '/') {
				s.Advance()
			}
			if s.pos+1 < len(s.src) {
				s.Advance()
				s.Advance()
			}
			continue
		}
		break
	}
	return s.line != startLine
}

// SkipToStatementBoundary walks forward until just past the next `;`
// or `}` at top-level, or to EOF. Used by the Parser to recover from
// unrecognised input without losing subsequent valid statements.
// Strings/comments are coarsely handled so the scan doesn't mistake
// content for a boundary.
func (s *Scanner) SkipToStatementBoundary() {
	depth := 0
	for !s.AtEnd() {
		c := s.src[s.pos]
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
			if s.pos+1 < len(s.src) && s.src[s.pos+1] == '/' {
				for !s.AtEnd() && s.src[s.pos] != '\n' {
					s.Advance()
				}
			} else if s.pos+1 < len(s.src) && s.src[s.pos+1] == '*' {
				s.Advance()
				s.Advance()
				for s.pos+1 < len(s.src) && !(s.src[s.pos] == '*' && s.src[s.pos+1] == '/') {
					s.Advance()
				}
				if s.pos+1 < len(s.src) {
					s.Advance()
					s.Advance()
				}
			} else {
				s.Advance()
			}
		case '"', '\'':
			s.Advance()
			for !s.AtEnd() && s.src[s.pos] != c {
				if s.src[s.pos] == '\\' && s.pos+1 < len(s.src) {
					s.Advance()
					s.Advance()
					continue
				}
				if s.src[s.pos] == '\n' {
					break
				}
				s.Advance()
			}
			if !s.AtEnd() {
				s.Advance()
			}
		default:
			s.Advance()
		}
	}
}

// ---- template literal scanning ---------------------------------- //

// TemplatePart is the result of one call to ScanTemplatePart. It mirrors
// the TemplateLiteralToken produced by Java's
// nextTemplateLiteralTokenShared (Scanner.java lines 938-981).
//
// Kind is one of:
//   - TemplateNoSub  — the whole template has no substitutions: `abc`
//   - TemplateHead   — the text before the first ${
//   - TemplateMiddle — the text between two ${ … } substitutions
//   - TemplateTail   — the text after the last } up to the closing `
type TemplatePartKind int

const (
	TemplateNoSub  TemplatePartKind = iota // NO_SUBSTITUTION_TEMPLATE
	TemplateHead                           // TEMPLATE_HEAD
	TemplateMiddle                         // TEMPLATE_MIDDLE
	TemplateTail                           // TEMPLATE_TAIL
)

// TemplatePart carries the cooked text segment and the kind of token.
type TemplatePart struct {
	Kind TemplatePartKind
	// Cooked is the raw source characters between the delimiters
	// (backtick/`${`/`}` boundaries), not including those delimiters.
	// Escape sequences are preserved verbatim (no decoding) to match the
	// Java scanner which stores the raw string and lets IR nodes carry it.
	Cooked string
}

// ScanTemplateLiteral is called when the scanner is positioned just after
// the opening backtick. It mirrors Java's scanTemplateLiteral (line 929)
// which immediately delegates to nextTemplateLiteralTokenShared.
//
// Returns the first TemplatePart (either TemplateNoSub or TemplateHead)
// and whether scanning succeeded. On failure (unterminated literal) it
// returns false and the scanner position is left just past the last
// character consumed.
func (s *Scanner) ScanTemplateLiteral() (TemplatePart, bool) {
	return s.scanTemplateCharacters(TemplateNoSub, TemplateHead)
}

// ScanTemplateMiddleOrTail is called after the parser has finished
// parsing a ${…} substitution expression and consumed the closing `}`.
// It mirrors Java's nextTemplateLiteral (line 140-150), which calls
// nextTemplateLiteralTokenShared(TEMPLATE_TAIL, TEMPLATE_MIDDLE).
func (s *Scanner) ScanTemplateMiddleOrTail() (TemplatePart, bool) {
	return s.scanTemplateCharacters(TemplateTail, TemplateMiddle)
}

// scanTemplateCharacters ports nextTemplateLiteralTokenShared (lines
// 938-981). It collects raw characters until a backtick (→ endKind) or
// a `${` sequence (→ middleKind), skipping over escape sequences.
func (s *Scanner) scanTemplateCharacters(endKind, middleKind TemplatePartKind) (TemplatePart, bool) {
	start := s.pos
	for !s.AtEnd() {
		c := s.src[s.pos]
		switch c {
		case '`':
			// End of template (NO_SUBSTITUTION or TAIL).
			cooked := string(s.src[start:s.pos])
			s.Advance() // consume closing backtick
			s.lastTokenLine = s.line
			return TemplatePart{Kind: endKind, Cooked: cooked}, true
		case '\\':
			// Skip escape sequence: consume '\' plus the next char.
			// For line-terminators the Java scanner calls skipLineTerminator
			// which handles \r\n too; we match that by just advancing once
			// more if the next char is \n (covers the common case).
			s.Advance() // '\'
			if !s.AtEnd() {
				s.Advance() // escaped char
			}
		case '$':
			if s.pos+1 < len(s.src) && s.src[s.pos+1] == '{' {
				// Start of a substitution.
				cooked := string(s.src[start:s.pos])
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
	// Unterminated template literal — report failure.
	return TemplatePart{}, false
}


// ---- numeric literal scanning ----------------------------------- //

// ScanNumericLiteral mirrors Java's scanPostZero / scanPostDigit /
// scanFractionalNumericLiteral / scanExponentOfNumericLiteral.
// Cursor must be on the first digit (or '.' for .5 forms).
// Returns (raw, isFloat, ok); advances the scanner on success.
func (s *Scanner) ScanNumericLiteral() (string, bool, bool) {
	start := s.pos
	c := s.Peek()
	if c == '0' {
		s.Advance()
		switch p := s.Peek(); {
		case p == 'b' || p == 'B':
			s.Advance(); s.skipDigits(func(r rune) bool { return r == '0' || r == '1' })
			if s.Peek() == 'n' { s.Advance() }
			s.lastTokenLine = s.line; return s.Text(start, s.pos), false, true
		case p == 'o' || p == 'O':
			s.Advance(); s.skipDigits(isOctalDigit)
			if s.Peek() == 'n' { s.Advance() }
			s.lastTokenLine = s.line; return s.Text(start, s.pos), false, true
		case p == 'x' || p == 'X':
			s.Advance(); s.skipDigits(isHexDigit)
			if s.Peek() == 'n' { s.Advance() }
			s.lastTokenLine = s.line; return s.Text(start, s.pos), false, true
		case p == 'n':
			s.Advance()
			s.lastTokenLine = s.line; return s.Text(start, s.pos), false, true
		case p == '.':
			s.Advance(); s.skipDigits(isDecimalDigit); s.scanExponent()
			s.lastTokenLine = s.line; return s.Text(start, s.pos), true, true
		case p == 'e' || p == 'E':
			s.scanExponent()
			s.lastTokenLine = s.line; return s.Text(start, s.pos), true, true
		default:
			if isDecimalDigit(p) {
				s.skipDigits(isDecimalDigit)
				if s.Peek() == '.' {
					s.Advance(); s.skipDigits(isDecimalDigit); s.scanExponent()
					s.lastTokenLine = s.line; return s.Text(start, s.pos), true, true
				}
			}
			s.lastTokenLine = s.line; return s.Text(start, s.pos), false, true
		}
	}
	if c == '.' {
		s.Advance(); s.skipDigits(isDecimalDigit); s.scanExponent()
		s.lastTokenLine = s.line; return s.Text(start, s.pos), true, true
	}
	s.skipDigits(isDecimalDigit)
	if s.Peek() == 'n' {
		s.Advance(); s.lastTokenLine = s.line; return s.Text(start, s.pos), false, true
	}
	isFloat := false
	if s.Peek() == '.' { isFloat = true; s.Advance(); s.skipDigits(isDecimalDigit) }
	if s.scanExponent() { isFloat = true }
	s.lastTokenLine = s.line; return s.Text(start, s.pos), isFloat, true
}

// scanExponent advances past an optional e/E exponent; returns true if found.
func (s *Scanner) scanExponent() bool {
	if p := s.Peek(); p != 'e' && p != 'E' { return false }
	s.Advance()
	if p := s.Peek(); p == '+' || p == '-' { s.Advance() }
	s.skipDigits(isDecimalDigit)
	return true
}

// skipDigits advances past digit runs matching pred, including '_' separators
// (mirrors Java's skipDecimalDigits / skipHexDigits / skipOctalDigits / skipBinaryDigits).
func (s *Scanner) skipDigits(pred func(rune) bool) {
	for !s.AtEnd() {
		if c := s.Peek(); pred(c) {
			s.Advance()
		} else if c == '_' {
			s.Advance(); if pred(s.Peek()) { s.Advance() }
		} else {
			break
		}
	}
}

func isDecimalDigit(c rune) bool { return c >= '0' && c <= '9' }
func isHexDigit(c rune) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}
func isOctalDigit(c rune) bool { return c >= '0' && c <= '7' }

// ---- string literal scanning ------------------------------------ //

// ScanStringBody mirrors Java's scanStringLiteral + skipStringLiteralChar.
// Cursor must be just AFTER the opening quote. Returns the raw content
// (no quotes) and whether scanning succeeded; consumes the closing quote.
func (s *Scanner) ScanStringBody(terminator rune) (string, bool) {
	start := s.pos
	for !s.AtEnd() && s.Peek() != terminator {
		c := s.Peek()
		if c == '\n' || c == '\r' { return "", false }
		if c == '\\' {
			s.Advance()
			if s.AtEnd() { return "", false }
			nc := s.Peek()
			if nc == '\r' { s.Advance(); if s.Peek() == '\n' { s.Advance() }; continue }
			if nc == '\n' { s.Advance(); continue }
			s.Advance(); continue
		}
		s.Advance()
	}
	if s.AtEnd() { return "", false }
	raw := s.Text(start, s.pos)
	s.Advance(); s.lastTokenLine = s.line
	return raw, true
}

// DecodeStringEscapes decodes the raw body of a string literal (text
// between, not including, quotes). Mirrors Java's
// skipStringLiteralEscapeSequence for \n \r \t \\ \' \" \0 \b \f \v
// \xXX \uXXXX \u{XXXXX} and line continuations.
// Unknown escapes pass through the character verbatim.
func DecodeStringEscapes(raw string) string {
	if !strings.ContainsRune(raw, '\\') { return raw }
	var b strings.Builder
	b.Grow(len(raw))
	rs := []rune(raw)
	for i := 0; i < len(rs); {
		c := rs[i]; i++
		if c != '\\' { b.WriteRune(c); continue }
		if i >= len(rs) { b.WriteByte('\\'); break }
		esc := rs[i]; i++
		switch esc {
		case 'n': b.WriteByte('\n')
		case 'r': b.WriteByte('\r')
		case 't': b.WriteByte('\t')
		case 'b': b.WriteByte('\b')
		case 'f': b.WriteByte('\f')
		case 'v': b.WriteByte('\v')
		case '0': b.WriteByte(0)
		case '\\', '\'', '"', '`': b.WriteRune(esc)
		case '\r': if i < len(rs) && rs[i] == '\n' { i++ } // line continuation
		case '\n': // line continuation — swallow
		case 'x':
			if i+1 < len(rs) && isHexDigit(rs[i]) && isHexDigit(rs[i+1]) {
				b.WriteRune(rune(hexNib(rs[i])<<4 | hexNib(rs[i+1]))); i += 2
			} else { b.WriteByte('\\'); b.WriteRune(esc) }
		case 'u':
			if i < len(rs) && rs[i] == '{' {
				i++; j := i
				for i < len(rs) && rs[i] != '}' { i++ }
				if v, ok := parseHexRunes(rs[j:i]); ok { b.WriteRune(rune(v)) } else {
					b.WriteString(`\u{`); b.WriteString(string(rs[j:i])); b.WriteByte('}')
				}
				if i < len(rs) { i++ }
			} else if i+3 < len(rs) && isHexDigit(rs[i]) && isHexDigit(rs[i+1]) &&
				isHexDigit(rs[i+2]) && isHexDigit(rs[i+3]) {
				b.WriteRune(rune(hexNib(rs[i])<<12 | hexNib(rs[i+1])<<8 |
					hexNib(rs[i+2])<<4 | hexNib(rs[i+3]))); i += 4
			} else { b.WriteByte('\\'); b.WriteRune(esc) }
		default: b.WriteRune(esc)
		}
	}
	return b.String()
}

func hexNib(c rune) int {
	if c >= '0' && c <= '9' { return int(c - '0') }
	if c >= 'a' && c <= 'f' { return int(c-'a') + 10 }
	return int(c-'A') + 10
}

func parseHexRunes(rs []rune) (int, bool) {
	if len(rs) == 0 { return 0, false }
	v := 0
	for _, c := range rs { if !isHexDigit(c) { return 0, false }; v = v<<4 | hexNib(c) }
	return v, true
}

// isStringLineTerminator mirrors Java's isStringLineTerminator (ES5 §7.3).
func isStringLineTerminator(c rune) bool { return c == '\n' || c == '\r' }

// ---- character-class predicates --------------------------------- //

func isIdentStart(c rune) bool {
	return unicode.IsLetter(c) || c == '_' || c == '$'
}

func isIdentCont(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsDigit(c) || c == '_' || c == '$'
}

// ---- regular-expression literal ---------------------------------- //

// ScanRegularExpressionLiteral mirrors Scanner.java's
// nextRegularExpressionLiteralToken. The caller must have confirmed
// that the cursor is on the opening `/` and that context requires a
// regex (not division). Returns ("", false) on scan failure.
//
// Grammar (ES5 §7.8.5):
//
//	RegularExpressionLiteral :: `/` RegularExpressionBody `/` RegularExpressionFlags
func (s *Scanner) ScanRegularExpressionLiteral() (string, bool) {
	begin := s.pos

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

	// flags — zero or more IdentifierPart characters
	for isIdentCont(s.Peek()) {
		s.Advance()
	}

	s.lastTokenLine = s.line
	return string(s.src[begin:s.pos]), true
}

// skipRegularExpressionBody mirrors Scanner.java's skipRegularExpressionBody.
func (s *Scanner) skipRegularExpressionBody() bool {
	if !isRegularExpressionFirstChar(s.Peek()) {
		return false
	}
	if !s.skipRegularExpressionChar() {
		return false
	}
	for !s.AtEnd() && isRegularExpressionChar(s.Peek()) {
		if !s.skipRegularExpressionChar() {
			return false
		}
	}
	return true
}

// skipRegularExpressionChar mirrors Scanner.java's skipRegularExpressionChar.
func (s *Scanner) skipRegularExpressionChar() bool {
	switch s.Peek() {
	case '\\':
		return s.skipRegularExpressionBackslashSequence()
	case '[':
		return s.skipRegularExpressionClass()
	default:
		s.Advance()
		return true
	}
}

// skipRegularExpressionBackslashSequence mirrors
// Scanner.java's skipRegularExpressionBackslashSequence.
func (s *Scanner) skipRegularExpressionBackslashSequence() bool {
	s.Advance() // consume '\'
	if isLineTerminator(s.Peek()) {
		return false
	}
	s.Advance()
	return true
}

// skipRegularExpressionClass mirrors Scanner.java's skipRegularExpressionClass.
func (s *Scanner) skipRegularExpressionClass() bool {
	s.Advance() // consume '['
	for !s.AtEnd() && peekRegularExpressionClassChar(s.Peek()) {
		if !s.skipRegularExpressionClassChar() {
			return false
		}
	}
	if s.Peek() != ']' {
		return false
	}
	s.Advance()
	return true
}

// peekRegularExpressionClassChar mirrors
// Scanner.java's peekRegularExpressionClassChar.
func peekRegularExpressionClassChar(ch rune) bool {
	return ch != ']' && !isLineTerminator(ch)
}

// skipRegularExpressionClassChar mirrors
// Scanner.java's skipRegularExpressionClassChar.
func (s *Scanner) skipRegularExpressionClassChar() bool {
	if s.Peek() == '\\' {
		return s.skipRegularExpressionBackslashSequence()
	}
	s.Advance()
	return true
}

// isRegularExpressionFirstChar mirrors Scanner.java's isRegularExpressionFirstChar.
func isRegularExpressionFirstChar(ch rune) bool {
	return isRegularExpressionChar(ch) && ch != '*'
}

// isRegularExpressionChar mirrors Scanner.java's isRegularExpressionChar.
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

// isLineTerminator mirrors Scanner.java's isLineTerminator (ES5 §7.3).
func isLineTerminator(ch rune) bool {
	switch ch {
	case '\n', '\r', ' ', ' ':
		return true
	default:
		return false
	}
}
