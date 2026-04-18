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

import "unicode"

// Scanner walks a source-string by rune. All positions are reported
// as 1-indexed line, 0-indexed col to match Java's getLineno /
// getCharno convention.
type Scanner struct {
	src     []rune
	pos     int
	line    int
	lineBOL int
}

// NewScanner mirrors Java's `new Scanner(source)` constructor.
func NewScanner(src string) *Scanner {
	return &Scanner{src: []rune(src), line: 1}
}

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
	pos, line, lineBOL int
}

func (s *Scanner) Save() ScannerMark     { return ScannerMark{s.pos, s.line, s.lineBOL} }
func (s *Scanner) Restore(m ScannerMark) { s.pos, s.line, s.lineBOL = m.pos, m.line, m.lineBOL }

// ---- byte-level conditionals ------------------------------------- //

// ConsumeByte advances past `b` if the current rune matches. Returns
// whether it consumed.
func (s *Scanner) ConsumeByte(b byte) bool {
	if s.AtEnd() || byte(s.src[s.pos]) != b {
		return false
	}
	s.Advance()
	return true
}

// ---- identifier / keyword --------------------------------------- //

// ConsumeIdent reads a JS identifier starting at the cursor.
// Returns ("", false) if the cursor isn't on an ident start.
func (s *Scanner) ConsumeIdent() (string, bool) {
	if s.AtEnd() {
		return "", false
	}
	c := s.src[s.pos]
	if !isIdentStart(c) {
		return "", false
	}
	start := s.pos
	for !s.AtEnd() && isIdentCont(s.src[s.pos]) {
		s.Advance()
	}
	return string(s.src[start:s.pos]), true
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

// ---- character-class predicates --------------------------------- //

func isIdentStart(c rune) bool {
	return unicode.IsLetter(c) || c == '_' || c == '$'
}

func isIdentCont(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsDigit(c) || c == '_' || c == '$'
}
