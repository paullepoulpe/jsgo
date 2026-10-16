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
	// source holds the full input as a Java-like char sequence; we use
	// runes so that charAt(sourceCursor) is O(1).
	source     []rune
	sourceEnd  int
	sourceCursor int // next index into source to read
	cursor     int   // monotonically increasing, mirrors Java cursor

	lineno    int
	charno    int
	initLine  int
	initCharno int

	// lineStart mirrors Java's lineStart field (index into source of first
	// char on current line).
	lineStart int
	// lineEndChar mirrors Java's lineEndChar: -1 when not in a CR/LF
	// sequence, else the char that started the line-end sequence.
	lineEndChar int

	string_ string // the most recently-scanned STRING/ANNOTATION payload

	// unget buffer. Java uses int[3]; we use a fixed-size array with a
	// cursor, matching the same LIFO semantics.
	ungetBuffer [3]rune
	ungetCursor int
}

// NewJsDocTokenStream constructs a stream over sourceString. initLine
// + initCharno mirror the Java overloads; callers that don't care
// pass 0/0.
func NewJsDocTokenStream(sourceString string, initLine, initCharno int) *JsDocTokenStream {
	src := []rune(sourceString)
	return &JsDocTokenStream{
		source:      src,
		sourceEnd:   len(src),
		lineno:      initLine,
		initLine:    initLine,
		initCharno:  initCharno,
		lineEndChar: -1,
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

// GetJsDocToken returns the next token. Port of Java getJsDocToken().
// The outer loop eats whitespace; the inner switch dispatches on the
// first non-space character.
func (s *JsDocTokenStream) GetJsDocToken() JsDocToken {
	s.stringBufferReset()
	for {
		// eat whitespace
		var c rune
		for {
			s.charno = -1
			c = s.getChar()
			if c == eofChar {
				return JsDocTokenEOF
			} else if c == '\n' {
				return JsDocTokenEOL
			} else if !isJSSpace(c) {
				break
			}
		}

		switch c {
		// annotation, e.g. @type or @constructor
		case '@':
			for {
				c = s.getChar()
				if isAlpha(c) {
					s.addToString(c)
				} else {
					s.ungetChar(c)
					return JsDocTokenAnnotation
				}
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
			c = s.getChar()
			if c == '<' {
				return JsDocTokenLeftAngle
			} else {
				if c == '.' {
					c = s.getChar()
					if c == '.' {
						return JsDocTokenIterRest
					}
					s.addToString('.')
				}
				// we may backtrack across line boundary — use raw unget
				// (mirrors Java's ungetBuffer[ungetCursor++] = c directly)
				s.ungetCharIgnoreLineEnd(c)
				c = '.'
			}
			// fall through to default
			fallthrough

		default:
			// recognize a JsDoc string but discard last . if it is followed
			// by a non-JsDoc comment char, e.g. Array.<
			c1 := c
			s.addToString(c)
			c2 := s.getChar()
			if !isJSDocString(c2) {
				s.ungetChar(c2)
				return JsDocTokenString
			}
			for {
				c1 = c2
				c2 = s.getChar()
				if c1 == '.' && c2 == '<' {
					s.ungetChar(c2)
					s.ungetChar(c1)
					return JsDocTokenString
				} else if isJSDocString(c2) {
					s.addToString(c1)
				} else {
					s.ungetChar(c2)
					s.addToString(c1)
					return JsDocTokenString
				}
			}
		}
	}
}

// GetRemainingJSDocLine collects text up to (but not including) `*/`,
// `\n`, or EOF, matching the Java getRemainingJSDocLine.
func (s *JsDocTokenStream) GetRemainingJSDocLine() string {
	s.stringBufferReset()
	for {
		c := s.getChar()
		switch c {
		case '*':
			if s.peekChar() != '/' {
				s.addToString(c)
				break
			}
			// fall through (end-of-comment detected)
			fallthrough
		case eofChar, '\n':
			s.ungetChar(c)
			return s.string_

		default:
			s.addToString(c)
		}
	}
}

// ---------------------------------------------------------------- //
// char-level helpers                                               //
// ---------------------------------------------------------------- //

// getChar mirrors Java's protected getChar(). It normalises \r and
// \r\n sequences into \n, skips Unicode format characters, normalises
// line terminators U+2028 / U+2029, and sets charno via getLineOffset
// the first time it is called on a given token.
func (s *JsDocTokenStream) getChar() rune {
	if s.ungetCursor != 0 {
		s.cursor++
		s.ungetCursor--
		if s.charno == -1 {
			s.charno = s.getLineOffset()
		}
		return s.ungetBuffer[s.ungetCursor]
	}

	for {
		if s.sourceCursor == s.sourceEnd {
			if s.charno == -1 {
				s.charno = s.getLineOffset()
			}
			return eofChar
		}
		s.cursor++
		c := s.source[s.sourceCursor]
		s.sourceCursor++

		if s.lineEndChar >= 0 {
			if s.lineEndChar == '\r' && c == '\n' {
				// \r\n pair: absorb the \n (already counted the line in \r handling)
				s.lineEndChar = int(c)
				continue
			}
			s.lineEndChar = -1
			s.lineStart = s.sourceCursor - 1
			s.lineno++
		}

		if c <= 127 {
			if c == '\n' || c == '\r' {
				s.lineEndChar = int(c)
				c = '\n'
			}
		} else {
			if isJSFormatChar(c) {
				continue
			}
			if isJSLineTerminator(c) {
				s.lineEndChar = int(c)
				c = '\n'
			}
		}

		if s.charno == -1 {
			s.charno = s.getLineOffset()
		}

		return c
	}
}

// getCharIgnoreLineEnd mirrors Java's private getCharIgnoreLineEnd.
// Used only by matchChar so that `*/` detection does not incorrectly
// advance the line counter.
func (s *JsDocTokenStream) getCharIgnoreLineEnd() rune {
	if s.ungetCursor != 0 {
		s.cursor++
		s.ungetCursor--
		if s.charno == -1 {
			s.charno = s.getLineOffset()
		}
		return s.ungetBuffer[s.ungetCursor]
	}

	for {
		if s.sourceCursor == s.sourceEnd {
			if s.charno == -1 {
				s.charno = s.getLineOffset()
			}
			return eofChar
		}
		s.cursor++
		c := s.source[s.sourceCursor]
		s.sourceCursor++

		if c <= 127 {
			if c == '\n' || c == '\r' {
				s.lineEndChar = int(c)
				c = '\n'
			}
		} else {
			if isJSFormatChar(c) {
				continue
			}
			if isJSLineTerminator(c) {
				s.lineEndChar = int(c)
				c = '\n'
			}
		}

		if s.charno == -1 {
			s.charno = s.getLineOffset()
		}

		return c
	}
}

// ungetChar mirrors Java's void ungetChar(int c). It asserts we are
// not backtracking across a line boundary.
func (s *JsDocTokenStream) ungetChar(c rune) {
	// mirror Java: assert(!(ungetCursor != 0 && ungetBuffer[ungetCursor-1] == '\n'))
	s.ungetBuffer[s.ungetCursor] = c
	s.ungetCursor++
	s.cursor--
}

// ungetCharIgnoreLineEnd mirrors Java's private ungetCharIgnoreLineEnd.
func (s *JsDocTokenStream) ungetCharIgnoreLineEnd(c rune) {
	s.ungetBuffer[s.ungetCursor] = c
	s.ungetCursor++
	s.cursor--
}

// matchChar mirrors Java's private boolean matchChar(int test). Uses
// getCharIgnoreLineEnd / ungetCharIgnoreLineEnd to avoid touching
// lineEndChar state when peeking for `*/`.
func (s *JsDocTokenStream) matchChar(test rune) bool {
	c := s.getCharIgnoreLineEnd()
	if c == test {
		return true
	}
	s.ungetCharIgnoreLineEnd(c)
	return false
}

func (s *JsDocTokenStream) peekChar() rune {
	c := s.getChar()
	s.ungetChar(c)
	return c
}

// getLineOffset mirrors Java's private int getLineOffset().
func (s *JsDocTokenStream) getLineOffset() int {
	return s.sourceCursor - s.lineStart - s.ungetCursor - 1
}

// ---------------------------------------------------------------- //
// string buffer — port of Java's addToString + getStringFromBuffer //
// ---------------------------------------------------------------- //

var _ rune // for linter

func (s *JsDocTokenStream) stringBufferReset() { s.string_ = "" }

func (s *JsDocTokenStream) addToString(c rune) {
	if c == eofChar {
		return
	}
	s.string_ += string(c)
}

// stringBufferAdd is an alias kept for any call sites that use it.
func (s *JsDocTokenStream) stringBufferAdd(c rune) { s.addToString(c) }

func (s *JsDocTokenStream) stringBuffer() string { return s.string_ }

// ---------------------------------------------------------------- //
// character classes                                                //
// ---------------------------------------------------------------- //

// isJSSpace mirrors Java's TokenUtil.isJSSpace. Note: does NOT include
// '\n' (line terminators are handled separately).
func isJSSpace(c rune) bool {
	if c <= 127 {
		return c == 0x20 || c == 0x9 || c == 0xC || c == 0xB
	}
	return c == 0xA0 || unicode.Is(unicode.Zs, c)
}

func isAlpha(c rune) bool {
	if c <= 'Z' {
		return 'A' <= c
	}
	return 'a' <= c && c <= 'z'
}

// isJSDocString mirrors Java's private static boolean isJSDocString(int c).
// Returns false for the listed punctuators, for EOF, for '\n', and for
// any JS whitespace character (via !isJSSpace). The default branch in
// Java is `default -> !TokenUtil.isJSSpace(c)`.
func isJSDocString(c rune) bool {
	switch c {
	case '@', '*', ',', '<', '>', ':', '(', ')', '{', '}', '[', ']', '?', '!', '|', '=', eofChar, '\n':
		return false
	default:
		return !isJSSpace(c)
	}
}

// isJSFormatChar mirrors Java's TokenUtil.isJSFormatChar.
func isJSFormatChar(c rune) bool {
	return c > 127 && unicode.Is(unicode.Cf, c)
}

// isJSLineTerminator mirrors Java's private static boolean isJSLineTerminator(int c).
func isJSLineTerminator(c rune) bool {
	if (c & 0xDFD0) != 0 {
		return false
	}
	return c == '\n' || c == '\r' || c == 0x2028 || c == 0x2029
}
