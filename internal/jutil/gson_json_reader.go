// Port of com/google/gson/stream/JsonReader.java (Gson 2.9.1), with
// com/google/gson/stream/JsonToken.java and JsonScope.java (all three files
// carry the same header)

/*
 * Copyright (C) 2010 Google Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jutil

import "unicode/utf16"

// Go-only: language-adapter — Gson is a library dependency of Closure; it is
// ported here as Java-semantics runtime. The reader keeps Gson's 1024-char
// buffer and refill logic because they are observable (a number longer than
// the buffer is read as an unquoted literal, and positions in error messages
// are computed from the buffer). nextDouble/nextLong/nextInt (which need
// Java's Double.parseDouble) and promoteNameToValue are not ported.

// JsonToken is com.google.gson.stream.JsonToken.
type JsonToken int32

const (
	JsonToken_BEGIN_ARRAY JsonToken = iota
	JsonToken_END_ARRAY
	JsonToken_BEGIN_OBJECT
	JsonToken_END_OBJECT
	JsonToken_NAME
	JsonToken_STRING
	JsonToken_NUMBER
	JsonToken_BOOLEAN
	JsonToken_NULL
	JsonToken_END_DOCUMENT
)

var jsonTokenNames = [...]string{
	"BEGIN_ARRAY", "END_ARRAY", "BEGIN_OBJECT", "END_OBJECT", "NAME",
	"STRING", "NUMBER", "BOOLEAN", "NULL", "END_DOCUMENT",
}

// String is the enum constant's name().
func (t JsonToken) String() string { return jsonTokenNames[t] }

// com.google.gson.stream.JsonScope
const (
	jsonScope_EMPTY_ARRAY       = 1
	jsonScope_NONEMPTY_ARRAY    = 2
	jsonScope_EMPTY_OBJECT      = 3
	jsonScope_DANGLING_NAME     = 4
	jsonScope_NONEMPTY_OBJECT   = 5
	jsonScope_EMPTY_DOCUMENT    = 6
	jsonScope_NONEMPTY_DOCUMENT = 7
	jsonScope_CLOSED            = 8
)

// JavaReader is java.io.Reader's read(char[], int, int) / close().
type JavaReader interface {
	// Read is read(char[] cbuf, int off, int len): the number of chars read,
	// or -1 at the end of the stream.
	Read(cbuf []uint16, off, length int) int
	// Close is close().
	Close()
}

// StringReader is java.io.StringReader.
type StringReader struct {
	next int

	// Go-only: performance — the string is decoded incrementally in Read
	// (s holds it from the start), so that CharStreams_ToString can return
	// an unread StringReader's string as it is, and Read never allocates a
	// full []uint16 slice.
	s            string
	converted    bool
	pendingLo    uint16
	hasPendingLo bool
}

// NewStringReader is new StringReader(String).
func NewStringReader(s string) *StringReader {
	return &StringReader{s: s}
}

// Read is StringReader.read(char[], int, int).
func (r *StringReader) Read(cbuf []uint16, off, length int) int {
	if off < 0 || length < 0 || off+length > len(cbuf) {
		panic(NewIndexOutOfBoundsException("Range [" + itoa(off) + ", " + itoa(off) + " + " + itoa(length) + ") out of bounds for length " + itoa(len(cbuf))))
	}
	if length == 0 {
		return 0
	}
	r.converted = true
	if !r.hasPendingLo && r.next >= len(r.s) {
		return -1
	}
	dst := cbuf[off : off+length]
	written := 0
	if r.hasPendingLo {
		dst[0] = r.pendingLo
		r.hasPendingLo = false
		written = 1
	}
	s := r.s
	next := r.next
	for written < length && next < len(s) {
		if s[next] < 0x80 {
			rem := min(length-written, len(s)-next)
			i := 0
			for ; i < rem && s[next+i] < 0x80; i++ {
				dst[written+i] = uint16(s[next+i])
			}
			written += i
			next += i
			continue
		}
		ru, sz := decodeWTF8(s[next:])
		next += sz
		if ru >= 0x10000 {
			hi, lo := utf16.EncodeRune(ru)
			dst[written] = uint16(hi)
			written++
			if written < length {
				dst[written] = uint16(lo)
				written++
			} else {
				r.pendingLo = uint16(lo)
				r.hasPendingLo = true
			}
		} else {
			dst[written] = uint16(ru)
			written++
		}
	}
	r.next = next
	if !r.hasPendingLo && next >= len(s) {
		r.s = ""
	}
	return written
}

// Close is StringReader.close().
func (r *StringReader) Close() {}

const jsonReader_MIN_INCOMPLETE_INTEGER = -0x8000000000000000 / 10

const (
	jsonReader_PEEKED_NONE               = 0
	jsonReader_PEEKED_BEGIN_OBJECT       = 1
	jsonReader_PEEKED_END_OBJECT         = 2
	jsonReader_PEEKED_BEGIN_ARRAY        = 3
	jsonReader_PEEKED_END_ARRAY          = 4
	jsonReader_PEEKED_TRUE               = 5
	jsonReader_PEEKED_FALSE              = 6
	jsonReader_PEEKED_NULL               = 7
	jsonReader_PEEKED_SINGLE_QUOTED      = 8
	jsonReader_PEEKED_DOUBLE_QUOTED      = 9
	jsonReader_PEEKED_UNQUOTED           = 10
	jsonReader_PEEKED_BUFFERED           = 11
	jsonReader_PEEKED_SINGLE_QUOTED_NAME = 12
	jsonReader_PEEKED_DOUBLE_QUOTED_NAME = 13
	jsonReader_PEEKED_UNQUOTED_NAME      = 14
	jsonReader_PEEKED_LONG               = 15
	jsonReader_PEEKED_NUMBER             = 16
	jsonReader_PEEKED_EOF                = 17
)

// State machine when parsing numbers
const (
	jsonReader_NUMBER_CHAR_NONE           = 0
	jsonReader_NUMBER_CHAR_SIGN           = 1
	jsonReader_NUMBER_CHAR_DIGIT          = 2
	jsonReader_NUMBER_CHAR_DECIMAL        = 3
	jsonReader_NUMBER_CHAR_FRACTION_DIGIT = 4
	jsonReader_NUMBER_CHAR_EXP_E          = 5
	jsonReader_NUMBER_CHAR_EXP_SIGN       = 6
	jsonReader_NUMBER_CHAR_EXP_DIGIT      = 7
)

const jsonReader_BUFFER_SIZE = 1024

// JsonReader is com.google.gson.stream.JsonReader: reads a JSON encoded
// value as a stream of tokens.
type JsonReader struct {
	in JavaReader

	// True to accept non-spec compliant JSON.
	lenient bool

	// Use a manual buffer to easily read and unread upcoming characters, and
	// also so we can create strings without an intermediate StringBuilder.
	buffer  []uint16
	builder []uint16
	pos     int
	limit   int

	lineNumber int
	lineStart  int

	peeked int

	// A peeked value that was composed entirely of digits with an optional
	// leading dash. Positive values may not have a leading 0.
	peekedLong int64

	// The number of characters in a peeked number literal.
	peekedNumberLength int

	// A peeked string that should be parsed on the next double, long or
	// string. This is populated before a numeric value is parsed and used if
	// that parsing fails.
	peekedString string

	// The nesting stack.
	stack     []int
	stackSize int

	// The path members. It corresponds directly to stack: At indices where
	// the stack contains an object (EMPTY_OBJECT, DANGLING_NAME or
	// NONEMPTY_OBJECT), pathNames contains the name at this scope. Where it
	// contains an array (EMPTY_ARRAY, NONEMPTY_ARRAY) pathIndices contains
	// the current index in that array. Otherwise the value is undefined, and
	// we take advantage of that by incrementing pathIndices when doing so
	// isn't useful.
	pathNames   []*string
	pathIndices []int
}

// NewJsonReader creates a new instance that reads a JSON-encoded stream
// from in.
//
// Java: JsonReader(Reader)
func NewJsonReader(in JavaReader) *JsonReader {
	r := &JsonReader{
		buffer:      make([]uint16, jsonReader_BUFFER_SIZE),
		stack:       make([]int, 32),
		pathNames:   make([]*string, 32),
		pathIndices: make([]int, 32),
	}
	r.stack[r.stackSize] = jsonScope_EMPTY_DOCUMENT
	r.stackSize++
	if in == nil {
		panic(NewNullPointerExceptionMsg("in == null"))
	}
	r.in = in
	return r
}

// SetLenient configures this parser to be liberal in what it accepts.
//
// Java: JsonReader.setLenient(boolean)
func (r *JsonReader) SetLenient(lenient bool) {
	r.lenient = lenient
}

// IsLenient returns true if this parser is liberal in what it accepts.
//
// Java: JsonReader.isLenient()
func (r *JsonReader) IsLenient() bool {
	return r.lenient
}

// BeginArray consumes the next token from the JSON stream and asserts that
// it is the beginning of a new array.
//
// Java: JsonReader.beginArray()
func (r *JsonReader) BeginArray() {
	p := r.peeked
	if p == jsonReader_PEEKED_NONE {
		p = r.doPeek()
	}
	if p == jsonReader_PEEKED_BEGIN_ARRAY {
		r.push(jsonScope_EMPTY_ARRAY)
		r.pathIndices[r.stackSize-1] = 0
		r.peeked = jsonReader_PEEKED_NONE
	} else {
		panic(NewIllegalStateExceptionMsg("Expected BEGIN_ARRAY but was " + r.Peek().String() + r.locationString()))
	}
}

// EndArray consumes the next token from the JSON stream and asserts that it
// is the end of the current array.
//
// Java: JsonReader.endArray()
func (r *JsonReader) EndArray() {
	p := r.peeked
	if p == jsonReader_PEEKED_NONE {
		p = r.doPeek()
	}
	if p == jsonReader_PEEKED_END_ARRAY {
		r.stackSize--
		r.pathIndices[r.stackSize-1]++
		r.peeked = jsonReader_PEEKED_NONE
	} else {
		panic(NewIllegalStateExceptionMsg("Expected END_ARRAY but was " + r.Peek().String() + r.locationString()))
	}
}

// BeginObject consumes the next token from the JSON stream and asserts that
// it is the beginning of a new object.
//
// Java: JsonReader.beginObject()
func (r *JsonReader) BeginObject() {
	p := r.peeked
	if p == jsonReader_PEEKED_NONE {
		p = r.doPeek()
	}
	if p == jsonReader_PEEKED_BEGIN_OBJECT {
		r.push(jsonScope_EMPTY_OBJECT)
		r.peeked = jsonReader_PEEKED_NONE
	} else {
		panic(NewIllegalStateExceptionMsg("Expected BEGIN_OBJECT but was " + r.Peek().String() + r.locationString()))
	}
}

// EndObject consumes the next token from the JSON stream and asserts that
// it is the end of the current object.
//
// Java: JsonReader.endObject()
func (r *JsonReader) EndObject() {
	p := r.peeked
	if p == jsonReader_PEEKED_NONE {
		p = r.doPeek()
	}
	if p == jsonReader_PEEKED_END_OBJECT {
		r.stackSize--
		r.pathNames[r.stackSize] = nil // Free the last path name so that it can be garbage collected!
		r.pathIndices[r.stackSize-1]++
		r.peeked = jsonReader_PEEKED_NONE
	} else {
		panic(NewIllegalStateExceptionMsg("Expected END_OBJECT but was " + r.Peek().String() + r.locationString()))
	}
}

// HasNext returns true if the current array or object has another element.
//
// Java: JsonReader.hasNext()
func (r *JsonReader) HasNext() bool {
	p := r.peeked
	if p == jsonReader_PEEKED_NONE {
		p = r.doPeek()
	}
	return p != jsonReader_PEEKED_END_OBJECT && p != jsonReader_PEEKED_END_ARRAY && p != jsonReader_PEEKED_EOF
}

// Peek returns the type of the next token without consuming it.
//
// Java: JsonReader.peek()
func (r *JsonReader) Peek() JsonToken {
	p := r.peeked
	if p == jsonReader_PEEKED_NONE {
		p = r.doPeek()
	}

	switch p {
	case jsonReader_PEEKED_BEGIN_OBJECT:
		return JsonToken_BEGIN_OBJECT
	case jsonReader_PEEKED_END_OBJECT:
		return JsonToken_END_OBJECT
	case jsonReader_PEEKED_BEGIN_ARRAY:
		return JsonToken_BEGIN_ARRAY
	case jsonReader_PEEKED_END_ARRAY:
		return JsonToken_END_ARRAY
	case jsonReader_PEEKED_SINGLE_QUOTED_NAME, jsonReader_PEEKED_DOUBLE_QUOTED_NAME, jsonReader_PEEKED_UNQUOTED_NAME:
		return JsonToken_NAME
	case jsonReader_PEEKED_TRUE, jsonReader_PEEKED_FALSE:
		return JsonToken_BOOLEAN
	case jsonReader_PEEKED_NULL:
		return JsonToken_NULL
	case jsonReader_PEEKED_SINGLE_QUOTED, jsonReader_PEEKED_DOUBLE_QUOTED, jsonReader_PEEKED_UNQUOTED, jsonReader_PEEKED_BUFFERED:
		return JsonToken_STRING
	case jsonReader_PEEKED_LONG, jsonReader_PEEKED_NUMBER:
		return JsonToken_NUMBER
	case jsonReader_PEEKED_EOF:
		return JsonToken_END_DOCUMENT
	default:
		panic(NewAssertionErrorNoMessage())
	}
}

func (r *JsonReader) doPeek() int {
	peekStack := r.stack[r.stackSize-1]
	if peekStack == jsonScope_EMPTY_ARRAY {
		r.stack[r.stackSize-1] = jsonScope_NONEMPTY_ARRAY
	} else if peekStack == jsonScope_NONEMPTY_ARRAY {
		// Look for a comma before the next element.
		c := r.nextNonWhitespace(true)
		switch c {
		case ']':
			r.peeked = jsonReader_PEEKED_END_ARRAY
			return r.peeked
		case ';':
			r.checkLenient() // fall-through
		case ',':
		default:
			panic(r.syntaxError("Unterminated array"))
		}
	} else if peekStack == jsonScope_EMPTY_OBJECT || peekStack == jsonScope_NONEMPTY_OBJECT {
		r.stack[r.stackSize-1] = jsonScope_DANGLING_NAME
		// Look for a comma before the next element.
		if peekStack == jsonScope_NONEMPTY_OBJECT {
			c := r.nextNonWhitespace(true)
			switch c {
			case '}':
				r.peeked = jsonReader_PEEKED_END_OBJECT
				return r.peeked
			case ';':
				r.checkLenient() // fall-through
			case ',':
			default:
				panic(r.syntaxError("Unterminated object"))
			}
		}
		c := r.nextNonWhitespace(true)
		switch c {
		case '"':
			r.peeked = jsonReader_PEEKED_DOUBLE_QUOTED_NAME
			return r.peeked
		case '\'':
			r.checkLenient()
			r.peeked = jsonReader_PEEKED_SINGLE_QUOTED_NAME
			return r.peeked
		case '}':
			if peekStack != jsonScope_NONEMPTY_OBJECT {
				r.peeked = jsonReader_PEEKED_END_OBJECT
				return r.peeked
			} else {
				panic(r.syntaxError("Expected name"))
			}
		default:
			r.checkLenient()
			r.pos-- // Don't consume the first character in an unquoted string.
			if r.isLiteral(uint16(c)) {
				r.peeked = jsonReader_PEEKED_UNQUOTED_NAME
				return r.peeked
			} else {
				panic(r.syntaxError("Expected name"))
			}
		}
	} else if peekStack == jsonScope_DANGLING_NAME {
		r.stack[r.stackSize-1] = jsonScope_NONEMPTY_OBJECT
		// Look for a colon before the value.
		c := r.nextNonWhitespace(true)
		switch c {
		case ':':
		case '=':
			r.checkLenient()
			if (r.pos < r.limit || r.fillBuffer(1)) && r.buffer[r.pos] == '>' {
				r.pos++
			}
		default:
			panic(r.syntaxError("Expected ':'"))
		}
	} else if peekStack == jsonScope_EMPTY_DOCUMENT {
		if r.lenient {
			r.consumeNonExecutePrefix()
		}
		r.stack[r.stackSize-1] = jsonScope_NONEMPTY_DOCUMENT
	} else if peekStack == jsonScope_NONEMPTY_DOCUMENT {
		c := r.nextNonWhitespace(false)
		if c == -1 {
			r.peeked = jsonReader_PEEKED_EOF
			return r.peeked
		} else {
			r.checkLenient()
			r.pos--
		}
	} else if peekStack == jsonScope_CLOSED {
		panic(NewIllegalStateExceptionMsg("JsonReader is closed"))
	}

	c := r.nextNonWhitespace(true)
	switch c {
	case ']':
		if peekStack == jsonScope_EMPTY_ARRAY {
			r.peeked = jsonReader_PEEKED_END_ARRAY
			return r.peeked
		}
		// fall-through to handle ",]"
		fallthrough
	case ';', ',':
		// In lenient mode, a 0-length literal in an array means 'null'.
		if peekStack == jsonScope_EMPTY_ARRAY || peekStack == jsonScope_NONEMPTY_ARRAY {
			r.checkLenient()
			r.pos--
			r.peeked = jsonReader_PEEKED_NULL
			return r.peeked
		} else {
			panic(r.syntaxError("Unexpected value"))
		}
	case '\'':
		r.checkLenient()
		r.peeked = jsonReader_PEEKED_SINGLE_QUOTED
		return r.peeked
	case '"':
		r.peeked = jsonReader_PEEKED_DOUBLE_QUOTED
		return r.peeked
	case '[':
		r.peeked = jsonReader_PEEKED_BEGIN_ARRAY
		return r.peeked
	case '{':
		r.peeked = jsonReader_PEEKED_BEGIN_OBJECT
		return r.peeked
	default:
		r.pos-- // Don't consume the first character in a literal value.
	}

	result := r.peekKeyword()
	if result != jsonReader_PEEKED_NONE {
		return result
	}

	result = r.peekNumber()
	if result != jsonReader_PEEKED_NONE {
		return result
	}

	if !r.isLiteral(r.buffer[r.pos]) {
		panic(r.syntaxError("Expected value"))
	}

	r.checkLenient()
	r.peeked = jsonReader_PEEKED_UNQUOTED
	return r.peeked
}

func (r *JsonReader) peekKeyword() int {
	// Figure out which keyword we're matching against by its first character.
	c := r.buffer[r.pos]
	var keyword, keywordUpper string
	var peeking int
	if c == 't' || c == 'T' {
		keyword = "true"
		keywordUpper = "TRUE"
		peeking = jsonReader_PEEKED_TRUE
	} else if c == 'f' || c == 'F' {
		keyword = "false"
		keywordUpper = "FALSE"
		peeking = jsonReader_PEEKED_FALSE
	} else if c == 'n' || c == 'N' {
		keyword = "null"
		keywordUpper = "NULL"
		peeking = jsonReader_PEEKED_NULL
	} else {
		return jsonReader_PEEKED_NONE
	}

	// Confirm that chars [1..length) match the keyword.
	length := len(keyword)
	for i := 1; i < length; i++ {
		if r.pos+i >= r.limit && !r.fillBuffer(i+1) {
			return jsonReader_PEEKED_NONE
		}
		c = r.buffer[r.pos+i]
		if c != uint16(keyword[i]) && c != uint16(keywordUpper[i]) {
			return jsonReader_PEEKED_NONE
		}
	}

	if (r.pos+length < r.limit || r.fillBuffer(length+1)) && r.isLiteral(r.buffer[r.pos+length]) {
		return jsonReader_PEEKED_NONE // Don't match trues, falsey or nullsoft!
	}

	// We've found the keyword followed either by EOF or by a non-literal
	// character.
	r.pos += length
	r.peeked = peeking
	return r.peeked
}

func (r *JsonReader) peekNumber() int {
	// Like nextNonWhitespace, this uses locals 'p' and 'l' to save
	// inner-loop field access.
	buffer := r.buffer
	p := r.pos
	l := r.limit

	var value int64 // Negative to accommodate Long.MIN_VALUE more easily.
	negative := false
	fitsInLong := true
	last := jsonReader_NUMBER_CHAR_NONE

	i := 0

charactersOfNumber:
	for ; true; i++ {
		if p+i == l {
			if i == len(buffer) {
				// Though this looks like a well-formed number, it's too long
				// to continue reading. Give up and let the application handle
				// this as an unquoted literal.
				return jsonReader_PEEKED_NONE
			}
			if !r.fillBuffer(i + 1) {
				break
			}
			p = r.pos
			l = r.limit
		}

		c := buffer[p+i]
		switch c {
		case '-':
			if last == jsonReader_NUMBER_CHAR_NONE {
				negative = true
				last = jsonReader_NUMBER_CHAR_SIGN
				continue
			} else if last == jsonReader_NUMBER_CHAR_EXP_E {
				last = jsonReader_NUMBER_CHAR_EXP_SIGN
				continue
			}
			return jsonReader_PEEKED_NONE

		case '+':
			if last == jsonReader_NUMBER_CHAR_EXP_E {
				last = jsonReader_NUMBER_CHAR_EXP_SIGN
				continue
			}
			return jsonReader_PEEKED_NONE

		case 'e', 'E':
			if last == jsonReader_NUMBER_CHAR_DIGIT || last == jsonReader_NUMBER_CHAR_FRACTION_DIGIT {
				last = jsonReader_NUMBER_CHAR_EXP_E
				continue
			}
			return jsonReader_PEEKED_NONE

		case '.':
			if last == jsonReader_NUMBER_CHAR_DIGIT {
				last = jsonReader_NUMBER_CHAR_DECIMAL
				continue
			}
			return jsonReader_PEEKED_NONE

		default:
			if c < '0' || c > '9' {
				if !r.isLiteral(c) {
					break charactersOfNumber
				}
				return jsonReader_PEEKED_NONE
			}
			if last == jsonReader_NUMBER_CHAR_SIGN || last == jsonReader_NUMBER_CHAR_NONE {
				value = -int64(c - '0')
				last = jsonReader_NUMBER_CHAR_DIGIT
			} else if last == jsonReader_NUMBER_CHAR_DIGIT {
				if value == 0 {
					return jsonReader_PEEKED_NONE // Leading '0' prefix is not allowed (since it could be octal).
				}
				newValue := value*10 - int64(c-'0')
				fitsInLong = fitsInLong && (value > jsonReader_MIN_INCOMPLETE_INTEGER ||
					(value == jsonReader_MIN_INCOMPLETE_INTEGER && newValue < value))
				value = newValue
			} else if last == jsonReader_NUMBER_CHAR_DECIMAL {
				last = jsonReader_NUMBER_CHAR_FRACTION_DIGIT
			} else if last == jsonReader_NUMBER_CHAR_EXP_E || last == jsonReader_NUMBER_CHAR_EXP_SIGN {
				last = jsonReader_NUMBER_CHAR_EXP_DIGIT
			}
		}
	}

	// We've read a complete number. Decide if it's a PEEKED_LONG or a
	// PEEKED_NUMBER.
	if last == jsonReader_NUMBER_CHAR_DIGIT && fitsInLong && (value != -0x8000000000000000 || negative) && (value != 0 || !negative) {
		if negative {
			r.peekedLong = value
		} else {
			r.peekedLong = -value
		}
		r.pos += i
		r.peeked = jsonReader_PEEKED_LONG
		return r.peeked
	} else if last == jsonReader_NUMBER_CHAR_DIGIT || last == jsonReader_NUMBER_CHAR_FRACTION_DIGIT ||
		last == jsonReader_NUMBER_CHAR_EXP_DIGIT {
		r.peekedNumberLength = i
		r.peeked = jsonReader_PEEKED_NUMBER
		return r.peeked
	} else {
		return jsonReader_PEEKED_NONE
	}
}

func (r *JsonReader) isLiteral(c uint16) bool {
	switch c {
	case '/', '\\', ';', '#', '=':
		r.checkLenient() // fall-through
		return false
	case '{', '}', '[', ']', ':', ',', ' ', '\t', '\f', '\r', '\n':
		return false
	default:
		return true
	}
}
