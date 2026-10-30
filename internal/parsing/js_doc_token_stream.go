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
	"unicode/utf8"

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

	stringBuffer    []byte
	stringScratch   [128]byte
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

	// Go-only: language-adapter — the Java String sourceString is scanned
	// directly in WTF-8 bytes (rawSource / byteCursor) while tracking exact
	// UTF-16 code-unit positions in sourceCursor and cursor.
	rawSource  string
	byteCursor int

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
	s.rawSource = sourceString
	s.byteCursor = 0
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
		if s.ungetCursor == 0 && s.lineEndChar < 0 {
			src := s.rawSource
			end := len(src)
			cur := s.byteCursor
			for cur < end {
				b := src[cur]
				if b != ' ' && b != '\t' && b != '\f' && b != '\v' {
					break
				}
				cur++
			}
			if delta := cur - s.byteCursor; delta > 0 {
				s.byteCursor = cur
				s.sourceCursor += delta
				s.cursor += delta
			}
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
			if s.ungetCursor == 0 {
				src := s.rawSource
				start := s.byteCursor
				cur := start
				end := len(src)
				for cur < end && isASCIIAlpha(src[cur]) {
					cur++
				}
				if cur == end || src[cur] < 0x80 {
					delta := cur - start
					s.byteCursor = cur
					s.sourceCursor += delta
					s.cursor += delta
					s.string_ = src[start:cur]
					return JsDocToken_ANNOTATION
				}
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
			fastDot := false
			if s.ungetCursor == 0 && s.lineEndChar < 0 && s.byteCursor < len(s.rawSource) {
				b1 := s.rawSource[s.byteCursor]
				if b1 == '<' {
					s.byteCursor++
					s.sourceCursor++
					s.cursor++
					return JsDocToken_LEFT_ANGLE
				}
				if b1 == '.' && s.byteCursor+1 < len(s.rawSource) && s.rawSource[s.byteCursor+1] == '.' {
					s.byteCursor += 2
					s.sourceCursor += 2
					s.cursor += 2
					return JsDocToken_ITER_REST
				}
				if b1 < 0x80 && b1 != '.' && b1 != '\r' && b1 != '\n' {
					c = '.'
					fastDot = true
				}
			}
			if !fastDot {
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
			}
			fallthrough

		default:
			{
				// recognize a JsDoc string but discard last . if it is followed by
				// a non-JsDoc comment char, e.g. Array.<
				if c < 128 && s.ungetCursor == 0 {
					src := s.rawSource
					start := s.byteCursor - 1
					cur := s.byteCursor
					end := len(src)
					for cur < end && src[cur] < 0x80 && asciiJSDocStringTable[src[cur]] {
						cur++
					}
					if cur == end || src[cur] < 0x80 {
						if cur > start+1 && src[cur-1] == '.' && cur < end && src[cur] == '<' {
							cur--
						}
						delta := cur - s.byteCursor
						s.byteCursor = cur
						s.sourceCursor += delta
						s.cursor += delta
						s.string_ = src[start:cur]
						return JsDocToken_STRING
					}
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
	if s.ungetCursor == 0 && s.stringBufferTop == 0 && s.lineEndChar < 0 {
		src := s.rawSource
		end := len(src)
		start := s.byteCursor
		cur := start
		utf16Len := 0
		canSlice := true
		for cur < end {
			b := src[cur]
			if b < 0x80 {
				if b == '\n' || b == '\r' {
					break
				}
				if b == '*' && cur+1 < end {
					if src[cur+1] == '/' {
						break
					}
					if src[cur+1] >= 0x80 {
						canSlice = false
						break
					}
				}
				cur++
				utf16Len++
				continue
			}
			r, width := decodeWTF8Rune(src[cur:])
			if rhino.TokenUtil_IsJSFormatChar(int(r)) || JsDocTokenStream_IsJSLineTerminator(int(r)) {
				canSlice = false
				break
			}
			cur += width
			if r >= 0x10000 {
				utf16Len += 2
			} else {
				utf16Len++
			}
		}
		if canSlice {
			if utf16Len > 0 {
				if s.charno == -1 {
					s.charno = s.sourceCursor - s.lineStart
				}
				s.byteCursor = cur
				s.sourceCursor += utf16Len
				s.cursor += utf16Len
			}
			stopChar := s.GetChar()
			s.UngetChar(stopChar)
			s.string_ = src[start:cur]
			return s.string_
		}
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
	b := s.stringBuffer[:n]
	if s.ungetCursor <= 1 {
		end := s.byteCursor
		if s.ungetCursor == 1 {
			uc := s.ungetBuffer[0]
			if uc >= 0 && uc < 0x80 && end > 0 && (s.rawSource[end-1] == byte(uc) || (uc == '\n' && s.rawSource[end-1] == '\r')) {
				end--
			}
		}
		start := end - len(b)
		if start >= 0 && end <= len(s.rawSource) && s.rawSource[start:end] == string(b) {
			return s.rawSource[start:end]
		}
	}
	return string(b)
}

func decodeWTF8Rune(s string) (rune, int) {
	c := s[0]
	if c == 0xED && len(s) >= 3 && s[1] >= 0xA0 && s[1] <= 0xBF && s[2]&0xC0 == 0x80 {
		return rune(c&0x0F)<<12 | rune(s[1]&0x3F)<<6 | rune(s[2]&0x3F), 3
	}
	return utf8.DecodeRuneInString(s)
}

func charUTF16Width(c int) int {
	if c >= 0x10000 {
		return 2
	}
	return 1
}

func (s *JsDocTokenStream) ungetUTF16Len() int {
	n := 0
	for i := 0; i < s.ungetCursor; i++ {
		n += charUTF16Width(s.ungetBuffer[i])
	}
	return n
}

// Java: JsDocTokenStream.AddToString(int)
func (s *JsDocTokenStream) AddToString(c int) {
	if s.stringBuffer == nil {
		s.stringBuffer = s.stringScratch[:0]
	}
	b := s.stringBuffer[:s.stringBufferTop]
	if c < 0x80 {
		b = append(b, byte(c))
	} else if c < 0x800 {
		b = append(b, byte(0xC0|(c>>6)), byte(0x80|(c&0x3F)))
	} else if c < 0x10000 {
		if c >= 0xDC00 && c <= 0xDFFF && len(b) >= 3 && b[len(b)-3] == 0xED && b[len(b)-2] >= 0xA0 && b[len(b)-2] <= 0xAF {
			hi := int(b[len(b)-3]&0x0F)<<12 | int(b[len(b)-2]&0x3F)<<6 | int(b[len(b)-1]&0x3F)
			cp := 0x10000 + ((hi - 0xD800) << 10) + (c - 0xDC00)
			b = utf8.AppendRune(b[:len(b)-3], rune(cp))
		} else {
			b = append(b, byte(0xE0|(c>>12)), byte(0x80|((c>>6)&0x3F)), byte(0x80|(c&0x3F)))
		}
	} else {
		b = utf8.AppendRune(b, rune(c))
	}
	s.stringBuffer = b
	s.stringBufferTop = len(b)
}

// Java: JsDocTokenStream.ungetChar(int)
func (s *JsDocTokenStream) UngetChar(c int) {
	// can not unread past across line boundary
	// (the Java assert is dropped, PORTING.md §5)
	s.ungetBuffer[s.ungetCursor] = c
	s.ungetCursor++
	s.cursor -= charUTF16Width(c)
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
		s.ungetCursor--
		c := s.ungetBuffer[s.ungetCursor]
		w := charUTF16Width(c)
		s.cursor += w
		if s.charno == -1 {
			s.charno = s.GetLineOffset() - (w - 1)
		}
		return c
	}

	for {
		if s.byteCursor == len(s.rawSource) {
			if s.charno == -1 {
				s.charno = s.GetLineOffset()
			}
			return JsDocTokenStream_EOF_CHAR
		}
		var c int
		w := 1
		b := s.rawSource[s.byteCursor]
		if b < 0x80 {
			c = int(b)
			s.byteCursor++
		} else {
			r, width := decodeWTF8Rune(s.rawSource[s.byteCursor:])
			s.byteCursor += width
			c = int(r)
			if c >= 0x10000 {
				w = 2
			}
		}
		s.cursor += w
		s.sourceCursor += w

		if !ignoreLineEnd && s.lineEndChar >= 0 {
			if s.lineEndChar == '\r' && c == '\n' {
				s.lineEndChar = '\n'
				continue
			}
			s.lineEndChar = -1
			s.lineStart = s.sourceCursor - w
			s.lineno++
		}

		if c <= 127 {
			if c == '\n' || c == '\r' {
				s.lineEndChar = c
				c = '\n'
			}
		} else if c <= 0xFFFF {
			if rhino.TokenUtil_IsJSFormatChar(c) {
				continue
			}
			if JsDocTokenStream_IsJSLineTerminator(c) {
				s.lineEndChar = c
				c = '\n'
			}
		}

		if s.charno == -1 {
			s.charno = s.GetLineOffset() - (w - 1)
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
	return s.sourceCursor - s.lineStart - s.ungetUTF16Len() - 1
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
