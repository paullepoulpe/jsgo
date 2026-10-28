// Port of com/google/javascript/jscomp/parsing/JsDocTokenStream.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/JsDocToken.java @ closure-compiler v20261021

/*
 * Copyright 2009 The Closure Compiler Authors.
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
	stringBufferBuf [128]uint16
	asciiBuffer     []byte
	asciiBufferBuf  [128]byte
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
	// only read through charAt, held as its UTF-16 code units when non-ASCII
	// or directly as rawSource when pure ASCII.
	sourceString []uint16
	rawSource    string
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
	s := &JsDocTokenStream{}
	s.Init(sourceString, lineno, initCharno)
	return s
}

// Init resets the stream for a new JSDoc comment without allocating.
func (s *JsDocTokenStream) Init(sourceString string, lineno int, initCharno int) {
	s.string_ = ""
	s.stringBufferTop = 0
	s.ungetCursor = 0
	s.lineStart = 0
	s.lineEndChar = -1
	s.charno = -1
	// (checkNotNull(sourceString): a Go string is never null.)
	s.lineno = lineno
	if jutil.IsASCII(sourceString) {
		s.rawSource = sourceString
		s.sourceString = nil
		s.sourceEnd = len(sourceString)
		s.asciiBuffer = s.asciiBufferBuf[:]
	} else {
		s.rawSource = ""
		s.sourceString = jutil.StringToUTF16(sourceString)
		s.sourceEnd = len(s.sourceString)
		s.stringBuffer = s.stringBufferBuf[:]
	}
	s.cursor = 0
	s.sourceCursor = 0
	s.initLineno = lineno
	s.initCharno = initCharno
}

var asciiJSDocStringTable = func() [128]bool {
	var tbl [128]bool
	for c := 0; c < 128; c++ {
		if c != '\r' && JsDocTokenStream_IsJSDocString(c) {
			tbl[c] = true
		}
	}
	return tbl
}()

func isASCIIAlpha(b byte) bool {
	return (b|0x20)-'a' < 26
}

// GetJsDocToken tokenizes JSDoc comments.
//
// Java: JsDocTokenStream.getJsDocToken()
func (s *JsDocTokenStream) GetJsDocToken() JsDocToken {
	var c int
	s.stringBufferTop = 0
	for {
		// Fast-skip ASCII horizontal whitespace before per-char GetChar().
		if s.sourceString == nil && s.ungetCursor == 0 && s.lineEndChar < 0 {
			src := s.rawSource
			end := s.sourceEnd
			cur := s.sourceCursor
			for cur < end {
				b := src[cur]
				if b != ' ' && b != '\t' && b != '\f' && b != '\v' {
					break
				}
				cur++
			}
			s.cursor += cur - s.sourceCursor
			s.sourceCursor = cur
		}
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
			if s.sourceString == nil && s.ungetCursor == 0 {
				src := s.rawSource
				start := s.sourceCursor
				cur := start
				end := s.sourceEnd
				for cur < end && isASCIIAlpha(src[cur]) {
					cur++
				}
				s.sourceCursor = cur
				s.cursor += cur - start
				s.string_ = src[start:cur]
				return JsDocToken_ANNOTATION
			}
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
				if s.sourceString == nil && s.ungetCursor == 0 {
					src := s.rawSource
					start := s.sourceCursor - 1
					cur := s.sourceCursor
					end := s.sourceEnd
					for cur < end && asciiJSDocStringTable[src[cur]] {
						cur++
					}
					if cur > start+1 && src[cur-1] == '.' && cur < end && src[cur] == '<' {
						cur--
					}
					s.cursor += cur - s.sourceCursor
					s.sourceCursor = cur
					s.string_ = src[start:cur]
					return JsDocToken_STRING
				}
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
	if s.sourceString == nil && s.ungetCursor == 0 && s.stringBufferTop == 0 {
		src := s.rawSource
		end := s.sourceEnd
		if s.lineEndChar >= 0 {
			if s.sourceCursor < end && s.lineEndChar == '\r' && src[s.sourceCursor] == '\n' {
				s.lineEndChar = '\n'
				s.cursor++
				s.sourceCursor++
			}
			if s.sourceCursor < end {
				s.lineEndChar = -1
				s.lineStart = s.sourceCursor
				s.lineno++
			}
		}
		start := s.sourceCursor
		cur := start
		for cur < end {
			b := src[cur]
			if b == '\n' || b == '\r' || (b == '*' && cur+1 < end && src[cur+1] == '/') {
				break
			}
			cur++
		}
		if cur > start {
			if s.charno == -1 {
				s.charno = start - s.lineStart
			}
			s.cursor += cur - start
			s.sourceCursor = cur
		}
		stopChar := s.GetChar()
		s.UngetChar(stopChar)
		s.string_ = src[start:cur]
		return s.string_
	}
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
	n := s.stringBufferTop
	if n == 0 {
		return ""
	}
	if s.sourceString == nil {
		end := s.sourceCursor - s.ungetCursor
		start := end - n
		if start >= 0 && end <= len(s.rawSource) && s.rawSource[start:end] == string(s.asciiBuffer[:n]) {
			return s.rawSource[start:end]
		}
		return string(s.asciiBuffer[:n])
	}
	return jutil.StringFromUTF16(s.stringBuffer[0:n])
}

// Java: JsDocTokenStream.AddToString(int)
func (s *JsDocTokenStream) AddToString(c int) {
	n := s.stringBufferTop
	if s.sourceString == nil {
		if n == len(s.asciiBuffer) {
			tmp := make([]byte, len(s.asciiBuffer)*2)
			copy(tmp, s.asciiBuffer[:n])
			s.asciiBuffer = tmp
		}
		s.asciiBuffer[n] = byte(c)
		s.stringBufferTop = n + 1
		return
	}
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
	return s.getChar(false)
}

// Java: JsDocTokenStream.GetCharIgnoreLineEnd()
func (s *JsDocTokenStream) GetCharIgnoreLineEnd() int {
	return s.getChar(true)
}

func (s *JsDocTokenStream) getChar(ignoreLineEnd bool) int {
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
		if s.sourceString == nil {
			c = int(s.rawSource[s.sourceCursor])
		} else {
			c = int(s.sourceString[s.sourceCursor])
		}
		s.sourceCursor++

		if !ignoreLineEnd && s.lineEndChar >= 0 {
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
	s.UngetChar(c)
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

// JsDocToken holds the JSDoc-specific tokens.
//
// This class is based on Rhino's Token.
type JsDocToken int32

const (
	// Tokens recycled from Rhino

	JsDocToken_EOF JsDocToken = iota // end of file token - (not EOF_CHAR)
	JsDocToken_EOL                   // end of line
	JsDocToken_LEFT_ANGLE
	JsDocToken_RIGHT_ANGLE
	JsDocToken_STRING
	JsDocToken_LEFT_SQUARE
	JsDocToken_RIGHT_SQUARE
	JsDocToken_LEFT_CURLY
	JsDocToken_RIGHT_CURLY
	JsDocToken_LEFT_PAREN
	JsDocToken_RIGHT_PAREN
	JsDocToken_COMMA // comma operator
	JsDocToken_COLON
	JsDocToken_ITER_REST

	// JsDoc-only tokens

	JsDocToken_ANNOTATION
	JsDocToken_PIPE
	JsDocToken_STAR
	JsDocToken_EOC
	JsDocToken_QMARK
	JsDocToken_BANG
	JsDocToken_EQUALS
)

var jsDocToken_names = [...]string{
	"EOF",
	"EOL",
	"LEFT_ANGLE",
	"RIGHT_ANGLE",
	"STRING",
	"LEFT_SQUARE",
	"RIGHT_SQUARE",
	"LEFT_CURLY",
	"RIGHT_CURLY",
	"LEFT_PAREN",
	"RIGHT_PAREN",
	"COMMA",
	"COLON",
	"ITER_REST",
	"ANNOTATION",
	"PIPE",
	"STAR",
	"EOC",
	"QMARK",
	"BANG",
	"EQUALS",
}

// Name is Enum.name().
//
// Java: JsDocToken.name()
func (t JsDocToken) Name() string { return jsDocToken_names[t] }

// String is Enum.toString().
//
// Java: JsDocToken.toString()
func (t JsDocToken) String() string { return t.Name() }
