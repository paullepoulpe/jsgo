// Port of com/google/javascript/jscomp/parsing/JsDocTokenStream.java @ closure-compiler v20260415

package parsing

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// JsDocTokenStream implements the scanner for JsDoc strings.
//
// It is heavily based on Rhino's TokenStream.
type JsDocTokenStream struct {
	// Set this to an initial non-null value so that the Parser has
	// something to retrieve even if an error has occurred and no
	// string is found.  Fosters one class of error, but saves lots of
	// code.
	string_ string

	stringBuffer    []uint16
	stringBufferTop int

	// Room to backtrace from to < on failed match of the last - in <!--
	ungetBuffer [3]int
	ungetCursor int

	lineStart   int
	lineEndChar int
	lineno      int
	charno      int
	initCharno  int
	initLineno  int

	// Go-only: language-adapter — the Java String sourceString, which is
	// only read through charAt, held as its UTF-16 code units.
	sourceString []uint16
	sourceEnd    int

	// sourceCursor is an index into a small buffer that keeps a
	// sliding window of the source stream.
	sourceCursor int

	// cursor is a monotonically increasing index into the original
	// source stream, tracking exactly how far scanning has progressed.
	// Its value is the index of the next character to be scanned.
	cursor int
}

// For chars - because we need something out-of-range
// to check.  (And checking EOF by exception is annoying.)
// Note distinction from EOF token type!
const JsDocTokenStream_EOF_CHAR = -1

// Java: JsDocTokenStream(String)
func NewJsDocTokenStream(sourceString string) *JsDocTokenStream {
	return NewJsDocTokenStreamWithLineno(sourceString, 0)
}

// Java: JsDocTokenStream(String, int)
func NewJsDocTokenStreamWithLineno(sourceString string, lineno int) *JsDocTokenStream {
	return NewJsDocTokenStreamWithLinenoCharno(sourceString, lineno, 0)
}

// Java: JsDocTokenStream(String, int, int)
func NewJsDocTokenStreamWithLinenoCharno(sourceString string, lineno int, initCharno int) *JsDocTokenStream {
	s := &JsDocTokenStream{
		string_:      "",
		stringBuffer: make([]uint16, 128),
		lineStart:    0,
		lineEndChar:  -1,
		charno:       -1,
	}
	// (checkNotNull(sourceString): a Go string is never null.)
	s.lineno = lineno
	s.sourceString = jutil.StringToUTF16(sourceString)
	s.sourceEnd = len(s.sourceString)
	s.cursor = 0
	s.sourceCursor = s.cursor
	s.initLineno = lineno
	s.initCharno = initCharno
	return s
}

// GetJsDocToken tokenizes JSDoc comments.
//
// Java: JsDocTokenStream.getJsDocToken()
func (s *JsDocTokenStream) GetJsDocToken() JsDocToken {
	var c int
	s.stringBufferTop = 0
	for {
		// eat white spaces
		for {
			s.charno = -1
			c = s.GetChar()
			if c == JsDocTokenStream_EOF_CHAR {
				return JsDocToken_EOF
			} else if c == '\n' {
				return JsDocToken_EOL
			} else if !rhino.TokenUtil_IsJSSpace(c) {
				break
			}
		}

		switch c {
		// annotation, e.g. @type or @constructor
		case '@':
			for {
				c = s.GetChar()
				if JsDocTokenStream_IsAlpha(c) {
					s.AddToString(c)
				} else {
					s.UngetChar(c)
					s.string_ = s.GetStringFromBuffer()
					s.stringBufferTop = 0
					return JsDocToken_ANNOTATION
				}
			}

		case '*':
			if s.MatchChar('/') {
				return JsDocToken_EOC
			} else {
				return JsDocToken_STAR
			}

		case ',':
			return JsDocToken_COMMA

		case '>':
			return JsDocToken_RIGHT_ANGLE

		case '(':
			return JsDocToken_LEFT_PAREN

		case ')':
			return JsDocToken_RIGHT_PAREN

		case '{':
			return JsDocToken_LEFT_CURLY

		case '}':
			return JsDocToken_RIGHT_CURLY

		case '[':
			return JsDocToken_LEFT_SQUARE

		case ']':
			return JsDocToken_RIGHT_SQUARE

		case '?':
			return JsDocToken_QMARK

		case '!':
			return JsDocToken_BANG

		case ':':
			return JsDocToken_COLON

		case '=':
			return JsDocToken_EQUALS

		case '|':
			return JsDocToken_PIPE

		case '<':
			return JsDocToken_LEFT_ANGLE

		case '.':
			c = s.GetChar()
			if c == '<' {
				return JsDocToken_LEFT_ANGLE
			} else {
				if c == '.' {
					c = s.GetChar()
					if c == '.' {
						return JsDocToken_ITER_REST
					} else {
						s.AddToString('.')
					}
				}
				// we may backtrack across line boundary
				s.ungetBuffer[s.ungetCursor] = c
				s.ungetCursor++
				c = '.'
			}
			fallthrough

		default:
			{
				// recognize a JsDoc string but discard last . if it is followed by
				// a non-JsDoc comment char, e.g. Array.<
				c1 := c
				s.AddToString(c)
				c2 := s.GetChar()
				if !JsDocTokenStream_IsJSDocString(c2) {
					s.UngetChar(c2)
					s.string_ = s.GetStringFromBuffer()
					s.stringBufferTop = 0
					return JsDocToken_STRING
				} else {
					for {
						c1 = c2
						c2 = s.GetChar()
						if c1 == '.' && c2 == '<' {
							s.UngetChar(c2)
							s.UngetChar(c1)
							s.string_ = s.GetStringFromBuffer()
							s.stringBufferTop = 0
							return JsDocToken_STRING
						} else {
							if JsDocTokenStream_IsJSDocString(c2) {
								s.AddToString(c1)
							} else {
								s.UngetChar(c2)
								s.AddToString(c1)
								s.string_ = s.GetStringFromBuffer()
								s.stringBufferTop = 0
								return JsDocToken_STRING
							}
						}
					}
				}
			}
		}
	}
}

// GetRemainingJSDocLine gets the remaining JSDoc line without the
// JsDocToken_EOL, JsDocToken_EOF or JsDocToken_EOC.
//
// Java: JsDocTokenStream.getRemainingJSDocLine()
func (s *JsDocTokenStream) GetRemainingJSDocLine() string {
	var c int
	for {
		c = s.GetChar()
		switch c {
		case '*':
			if s.PeekChar() != '/' {
				s.AddToString(c)
				break
			}
			fallthrough
		case JsDocTokenStream_EOF_CHAR, '\n':
			s.UngetChar(c)
			s.string_ = s.GetStringFromBuffer()
			s.stringBufferTop = 0
			return s.string_

		default:
			s.AddToString(c)
		}
	}
}

// Java: JsDocTokenStream.getLineno()
func (s *JsDocTokenStream) GetLineno() int { return s.lineno }

// Java: JsDocTokenStream.getCharno()
func (s *JsDocTokenStream) GetCharno() int {
	if s.lineno == s.initLineno {
		return s.initCharno + s.charno
	}
	return s.charno
}

// Java: JsDocTokenStream.getString()
func (s *JsDocTokenStream) GetString() string { return s.string_ }

// Java: JsDocTokenStream.GetStringFromBuffer()
func (s *JsDocTokenStream) GetStringFromBuffer() string {
	return jutil.StringFromUTF16(s.stringBuffer[0:s.stringBufferTop])
}

// Java: JsDocTokenStream.AddToString(int)
func (s *JsDocTokenStream) AddToString(c int) {
	n := s.stringBufferTop
	if n == len(s.stringBuffer) {
		tmp := make([]uint16, len(s.stringBuffer)*2)
		copy(tmp, s.stringBuffer[:n])
		s.stringBuffer = tmp
	}
	s.stringBuffer[n] = uint16(c)
	s.stringBufferTop = n + 1
}

// Java: JsDocTokenStream.ungetChar(int)
func (s *JsDocTokenStream) UngetChar(c int) {
	// can not unread past across line boundary
	// (the Java assert is dropped, PORTING.md §5)
	s.ungetBuffer[s.ungetCursor] = c
	s.ungetCursor++
	s.cursor--
}

// Java: JsDocTokenStream.MatchChar(int)
func (s *JsDocTokenStream) MatchChar(test int) bool {
	c := s.GetCharIgnoreLineEnd()
	if c == test {
		return true
	} else {
		s.UngetCharIgnoreLineEnd(c)
		return false
	}
}

// Java: JsDocTokenStream.isAlpha(int)
func JsDocTokenStream_IsAlpha(c int) bool {
	// Use 'Z' < 'a'
	if c <= 'Z' {
		return 'A' <= c
	} else {
		return 'a' <= c && c <= 'z'
	}
}

// Java: JsDocTokenStream.isJSDocString(int)
func JsDocTokenStream_IsJSDocString(c int) bool {
	switch c {
	case '@',
		'*',
		',',
		'<',
		'>',
		':',
		'(',
		')',
		'{',
		'}',
		'[',
		']',
		'?',
		'!',
		'|',
		'=',
		JsDocTokenStream_EOF_CHAR,
		'\n':
		return false
	default:
		return !rhino.TokenUtil_IsJSSpace(c)
	}
}

// Update allows the JSDocParser to update the character offset so that
// getCharno() returns a valid character position.
//
// Java: JsDocTokenStream.update()
func (s *JsDocTokenStream) Update() {
	s.charno = s.GetLineOffset()
}

// Java: JsDocTokenStream.PeekChar()
func (s *JsDocTokenStream) PeekChar() int {
	c := s.GetChar()
	s.UngetChar(c)
	return c
}

// Java: JsDocTokenStream.getChar()
func (s *JsDocTokenStream) GetChar() int {
	if s.ungetCursor != 0 {
		s.cursor++
		s.ungetCursor--
		if s.charno == -1 {
			s.charno = s.GetLineOffset()
		}
		return s.ungetBuffer[s.ungetCursor]
	}

	for {
		var c int
		if s.sourceCursor == s.sourceEnd {
			if s.charno == -1 {
				s.charno = s.GetLineOffset()
			}
			return JsDocTokenStream_EOF_CHAR
		}
		s.cursor++
		c = int(s.sourceString[s.sourceCursor])
		s.sourceCursor++

		if s.lineEndChar >= 0 {
			if s.lineEndChar == '\r' && c == '\n' {
				s.lineEndChar = '\n'
				continue
			}
			s.lineEndChar = -1
			s.lineStart = s.sourceCursor - 1
			s.lineno++
		}

		if c <= 127 {
			if c == '\n' || c == '\r' {
				s.lineEndChar = c
				c = '\n'
			}
		} else {
			if rhino.TokenUtil_IsJSFormatChar(c) {
				continue
			}
			if JsDocTokenStream_IsJSLineTerminator(c) {
				s.lineEndChar = c
				c = '\n'
			}
		}

		if s.charno == -1 {
			s.charno = s.GetLineOffset()
		}

		return c
	}
}

// Java: JsDocTokenStream.GetCharIgnoreLineEnd()
func (s *JsDocTokenStream) GetCharIgnoreLineEnd() int {
	if s.ungetCursor != 0 {
		s.cursor++
		s.ungetCursor--
		if s.charno == -1 {
			s.charno = s.GetLineOffset()
		}
		return s.ungetBuffer[s.ungetCursor]
	}

	for {
		var c int
		if s.sourceCursor == s.sourceEnd {
			if s.charno == -1 {
				s.charno = s.GetLineOffset()
			}
			return JsDocTokenStream_EOF_CHAR
		}
		s.cursor++
		c = int(s.sourceString[s.sourceCursor])
		s.sourceCursor++

		if c <= 127 {
			if c == '\n' || c == '\r' {
				s.lineEndChar = c
				c = '\n'
			}
		} else {
			if rhino.TokenUtil_IsJSFormatChar(c) {
				continue
			}
			if JsDocTokenStream_IsJSLineTerminator(c) {
				s.lineEndChar = c
				c = '\n'
			}
		}

		if s.charno == -1 {
			s.charno = s.GetLineOffset()
		}

		return c
	}
}

// Java: JsDocTokenStream.isJSLineTerminator(int)
func JsDocTokenStream_IsJSLineTerminator(c int) bool {
	// Optimization for faster check for eol character:
	// they do not have 0xDFD0 bits set
	if (c & 0xDFD0) != 0 {
		return false
	}
	return c == '\n' || c == '\r' || c == 0x2028 || c == 0x2029
}

// Java: JsDocTokenStream.UngetCharIgnoreLineEnd(int)
func (s *JsDocTokenStream) UngetCharIgnoreLineEnd(c int) {
	s.ungetBuffer[s.ungetCursor] = c
	s.ungetCursor++
	s.cursor--
}

// GetLineOffset returns the offset into the current line.
//
// Java: JsDocTokenStream.GetLineOffset()
func (s *JsDocTokenStream) GetLineOffset() int {
	return s.sourceCursor - s.lineStart - s.ungetCursor - 1
}

// Java: JsDocTokenStream.getCursor()
func (s *JsDocTokenStream) GetCursor() int {
	return s.cursor
}
