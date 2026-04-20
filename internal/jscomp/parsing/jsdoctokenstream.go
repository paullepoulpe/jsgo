// Port of com/google/javascript/jscomp/parsing/JsDocTokenStream.java @ closure-compiler v20260415
//
// Mirrors the Java tokenizer structure: getJsDocToken returns the
// next JsDocToken and updates lineno/charno; getRemainingJSDocLine is
// used by JsDocInfoParser to collect free text following a tag (e.g.
// the name after `@param {T}`).

package parsing

import "unicode"

// eofChar is returned by getChar at end of input. Distinct from
// JsDocTokenEOF (which is what the tokenizer emits on the token
// stream). Mirrors the same Java constant.
const eofChar = -1

type JsDocTokenStream struct {
	source    []rune
	pos       int
	lineno    int
	charno    int
	initLine  int
	initCharno int

	string_ string // the most recently-scanned STRING/ANNOTATION payload

	// tokenStartLine / tokenStartCol capture the (line, col) at which
	// the last emitted non-EOL / non-EOF token began — the first
	// non-whitespace char after the previous token was consumed.
	// Mirrors JsDocTokenStream.getStartLineno / getStartCharno — not
	// part of the Java API we previously ported, but added for
	// position-aware type diagnostics (see jstype_resolver).
	tokenStartLine int
	tokenStartCol  int

	// unget buffer. Closure uses a small fixed stack; we use a slice
	// with the same LIFO semantics.
	ungetBuf []rune
}

// NewJsDocTokenStream constructs a stream over sourceString. initLine
// + initCharno mirror the Java overloads; callers that don't care
// pass 0/0.
func NewJsDocTokenStream(sourceString string, initLine, initCharno int) *JsDocTokenStream {
	return &JsDocTokenStream{
		source:     []rune(sourceString),
		lineno:     initLine,
		initLine:   initLine,
		initCharno: initCharno,
	}
}

// String returns the last-scanned STRING / ANNOTATION payload. Named
// with the trailing underscore to avoid colliding with Go's default
// String() method; the Java side exposes this as getString().
func (s *JsDocTokenStream) String_() string { return s.string_ }

func (s *JsDocTokenStream) Lineno() int { return s.lineno }

// Charno mirrors Java's getCharno.
func (s *JsDocTokenStream) Charno() int {
	if s.lineno == s.initLine {
		return s.initCharno + s.charno
	}
	return s.charno
}

// GetJsDocToken returns the next token. Empty-bodied EOF is idempotent.
// Reproduces the Java switch statement exactly: whitespace is eaten,
// then a single dispatch covers all punctuator tokens, and the
// default branch collects a STRING (an identifier or keyword).
func (s *JsDocTokenStream) GetJsDocToken() JsDocToken {
	s.stringBufferReset()
	for {
		s.charno = -1
		// Snapshot the absolute body offset BEFORE consuming a char so
		// the saved (line, col) points at the first char of the next
		// token. Whitespace skipped below loops back to refresh it.
		s.tokenStartLine, s.tokenStartCol = s.computeAbsoluteLineCol(s.peekAbsolutePos())
		c := s.getChar()
		if c == eofChar {
			return JsDocTokenEOF
		}
		if c == '\n' {
			return JsDocTokenEOL
		}
		if !isJSSpace(c) {
			return s.dispatch(c)
		}
	}
}

// TokenStartLine returns the absolute line (0-indexed from comment
// start) of the most recently emitted non-EOL token. Used by the JSDoc
// parser to tag parsed type names with their origin position.
func (s *JsDocTokenStream) TokenStartLine() int { return s.tokenStartLine }

// TokenStartCol returns the column offset of the most recently emitted
// non-EOL token. Combined with TokenStartLine, this lets the resolver
// project a JSDoc-internal position onto the surrounding source file.
func (s *JsDocTokenStream) TokenStartCol() int { return s.tokenStartCol }

// peekAbsolutePos returns the absolute offset of the next character
// that would be returned by getChar, accounting for the unget buffer.
func (s *JsDocTokenStream) peekAbsolutePos() int {
	// Characters in the unget buffer were consumed from positions
	// before s.pos; conceptually the next absolute pos is s.pos -
	// len(ungetBuf). The body offset maps 1:1 onto the rune slice so
	// this is accurate as long as we don't skip any runes.
	return s.pos - len(s.ungetBuf)
}

// computeAbsoluteLineCol walks `source` from the start through `pos`
// counting newlines to derive a (line, col) pair. O(pos) per call; the
// body length is bounded by the enclosing comment so this is fine in
// practice. Mirrors how Java's JsDocTokenStream recomputes position on
// demand for diagnostics.
func (s *JsDocTokenStream) computeAbsoluteLineCol(pos int) (int, int) {
	line, col := 0, 0
	if pos > len(s.source) {
		pos = len(s.source)
	}
	for i := 0; i < pos; i++ {
		if s.source[i] == '\n' {
			line++
			col = 0
			continue
		}
		col++
	}
	return line, col
}

func (s *JsDocTokenStream) dispatch(c rune) JsDocToken {
	switch c {
	case '@':
		// annotation: `@` followed by [A-Za-z]+
		for {
			c2 := s.getChar()
			if isAlpha(c2) {
				s.stringBufferAdd(c2)
				continue
			}
			s.ungetChar(c2)
			s.string_ = s.stringBuffer()
			return JsDocTokenAnnotation
		}
	case '*':
		if s.matchChar('/') {
			return JsDocTokenEOC
		}
		return JsDocTokenStar
	case ',':
		return JsDocTokenComma
	case '>':
		return JsDocTokenRightAngle
	case '(':
		return JsDocTokenLeftParen
	case ')':
		return JsDocTokenRightParen
	case '{':
		return JsDocTokenLeftCurly
	case '}':
		return JsDocTokenRightCurly
	case '[':
		return JsDocTokenLeftSquare
	case ']':
		return JsDocTokenRightSquare
	case '?':
		return JsDocTokenQMark
	case '!':
		return JsDocTokenBang
	case ':':
		return JsDocTokenColon
	case '=':
		return JsDocTokenEquals
	case '|':
		return JsDocTokenPipe
	case '<':
		return JsDocTokenLeftAngle
	case '.':
		c2 := s.getChar()
		if c2 == '<' {
			return JsDocTokenLeftAngle
		}
		if c2 == '.' {
			c3 := s.getChar()
			if c3 == '.' {
				return JsDocTokenIterRest
			}
			s.stringBufferAdd('.')
			s.ungetChar(c3)
		}
		s.ungetChar(c2)
		// Fall through into the default STRING collector starting with '.'.
		return s.collectString('.')
	default:
		return s.collectString(c)
	}
}

// collectString reads a JSDoc STRING token starting with the given
// already-consumed first char. Mirrors the Java default branch: stops
// at non-JSDoc-string chars AND handles the special `.< ` lookahead
// that signals a generic-arg `<` (used in `Array.<T>`).
func (s *JsDocTokenStream) collectString(first rune) JsDocToken {
	s.stringBufferAdd(first)
	c1 := first
	c2 := s.getChar()
	if !isJSDocString(c2) {
		s.ungetChar(c2)
		s.string_ = s.stringBuffer()
		return JsDocTokenString
	}
	for {
		c1 = c2
		c2 = s.getChar()
		if c1 == '.' && c2 == '<' {
			s.ungetChar(c2)
			s.ungetChar(c1)
			s.string_ = s.stringBuffer()
			return JsDocTokenString
		}
		if isJSDocString(c2) {
			s.stringBufferAdd(c1)
			continue
		}
		s.ungetChar(c2)
		s.stringBufferAdd(c1)
		s.string_ = s.stringBuffer()
		return JsDocTokenString
	}
}

// GetRemainingJSDocLine collects text up to (but not including) `*/`,
// `\n`, or EOF, matching the Java variant used after `@param {T}` to
// consume the parameter name / description.
func (s *JsDocTokenStream) GetRemainingJSDocLine() string {
	s.stringBufferReset()
	for {
		c := s.getChar()
		switch {
		case c == '*' && s.peekChar() == '/':
			s.ungetChar(c)
			return s.stringBuffer()
		case c == eofChar, c == '\n':
			s.ungetChar(c)
			return s.stringBuffer()
		default:
			s.stringBufferAdd(c)
		}
	}
}

// ---------------------------------------------------------------- //
// char-level helpers (getChar / ungetChar / peekChar / matchChar)  //
// ---------------------------------------------------------------- //

func (s *JsDocTokenStream) getChar() rune {
	if len(s.ungetBuf) > 0 {
		c := s.ungetBuf[len(s.ungetBuf)-1]
		s.ungetBuf = s.ungetBuf[:len(s.ungetBuf)-1]
		s.charno++
		if c == '\n' {
			s.lineno++
			s.charno = 0
		}
		return c
	}
	if s.pos >= len(s.source) {
		return eofChar
	}
	c := s.source[s.pos]
	s.pos++
	s.charno++
	if c == '\n' {
		s.lineno++
		s.charno = 0
	}
	return c
}

func (s *JsDocTokenStream) ungetChar(c rune) {
	if c == eofChar {
		return
	}
	s.ungetBuf = append(s.ungetBuf, c)
	s.charno--
	if c == '\n' {
		s.lineno--
	}
}

func (s *JsDocTokenStream) peekChar() rune {
	c := s.getChar()
	s.ungetChar(c)
	return c
}

func (s *JsDocTokenStream) matchChar(c rune) bool {
	if s.peekChar() == c {
		s.getChar()
		return true
	}
	return false
}

// ---------------------------------------------------------------- //
// string buffer — port of Java's addToString + getStringFromBuffer //
// ---------------------------------------------------------------- //

var _ rune // for linter

func (s *JsDocTokenStream) stringBufferReset() { s.string_ = "" }

func (s *JsDocTokenStream) stringBufferAdd(c rune) {
	if c == eofChar {
		return
	}
	s.string_ += string(c)
}

func (s *JsDocTokenStream) stringBuffer() string { return s.string_ }

// ---------------------------------------------------------------- //
// character classes (match Java TokenUtil.isJSSpace + helpers)     //
// ---------------------------------------------------------------- //

func isJSSpace(c rune) bool {
	if c == '\t' || c == '\r' || c == ' ' || c == '\f' {
		return true
	}
	// Match Java's Character.isWhitespace for remaining cases.
	return unicode.IsSpace(c) && c != '\n'
}

func isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

// isJSDocString is the continuation predicate for the STRING token's
// default-branch accumulator. Mirrors the Java implementation.
func isJSDocString(c rune) bool {
	if c == eofChar {
		return false
	}
	switch c {
	case '\n', '\r', '\t', '\f', ' ', '@', '*', ',', '>', '<', '(', ')', '{', '}', '[', ']', '?', '!', ':', '=', '|':
		return false
	}
	return true
}
