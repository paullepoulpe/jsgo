// Port of com/google/gson/stream/JsonReader.java (Gson 2.9.1)

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

import "strings"

// NextName returns the next token, a property name, and consumes it.
//
// Java: JsonReader.nextName()
func (r *JsonReader) NextName() string {
	p := r.peeked
	if p == jsonReader_PEEKED_NONE {
		p = r.doPeek()
	}
	var result string
	if p == jsonReader_PEEKED_UNQUOTED_NAME {
		result = r.nextUnquotedValue()
	} else if p == jsonReader_PEEKED_SINGLE_QUOTED_NAME {
		result = r.nextQuotedValue('\'')
	} else if p == jsonReader_PEEKED_DOUBLE_QUOTED_NAME {
		result = r.nextQuotedValue('"')
	} else {
		panic(NewIllegalStateExceptionMsg("Expected a name but was " + r.Peek().String() + r.locationString()))
	}
	r.peeked = jsonReader_PEEKED_NONE
	r.pathNames[r.stackSize-1] = &result
	return result
}

// NextString returns the string value of the next token, consuming it. If
// the next token is a number, this method will return its string form.
//
// Java: JsonReader.nextString()
func (r *JsonReader) NextString() string {
	p := r.peeked
	if p == jsonReader_PEEKED_NONE {
		p = r.doPeek()
	}
	var result string
	if p == jsonReader_PEEKED_UNQUOTED {
		result = r.nextUnquotedValue()
	} else if p == jsonReader_PEEKED_SINGLE_QUOTED {
		result = r.nextQuotedValue('\'')
	} else if p == jsonReader_PEEKED_DOUBLE_QUOTED {
		result = r.nextQuotedValue('"')
	} else if p == jsonReader_PEEKED_BUFFERED {
		result = r.peekedString
		r.peekedString = ""
	} else if p == jsonReader_PEEKED_LONG {
		result = LongToString(r.peekedLong)
	} else if p == jsonReader_PEEKED_NUMBER {
		result = StringFromUTF16(r.buffer[r.pos : r.pos+r.peekedNumberLength])
		r.pos += r.peekedNumberLength
	} else {
		panic(NewIllegalStateExceptionMsg("Expected a string but was " + r.Peek().String() + r.locationString()))
	}
	r.peeked = jsonReader_PEEKED_NONE
	r.pathIndices[r.stackSize-1]++
	return result
}

// NextBoolean returns the boolean value of the next token, consuming it.
//
// Java: JsonReader.nextBoolean()
func (r *JsonReader) NextBoolean() bool {
	p := r.peeked
	if p == jsonReader_PEEKED_NONE {
		p = r.doPeek()
	}
	if p == jsonReader_PEEKED_TRUE {
		r.peeked = jsonReader_PEEKED_NONE
		r.pathIndices[r.stackSize-1]++
		return true
	} else if p == jsonReader_PEEKED_FALSE {
		r.peeked = jsonReader_PEEKED_NONE
		r.pathIndices[r.stackSize-1]++
		return false
	}
	panic(NewIllegalStateExceptionMsg("Expected a boolean but was " + r.Peek().String() + r.locationString()))
}

// NextNull consumes the next token from the JSON stream and asserts that it
// is a literal null.
//
// Java: JsonReader.nextNull()
func (r *JsonReader) NextNull() {
	p := r.peeked
	if p == jsonReader_PEEKED_NONE {
		p = r.doPeek()
	}
	if p == jsonReader_PEEKED_NULL {
		r.peeked = jsonReader_PEEKED_NONE
		r.pathIndices[r.stackSize-1]++
	} else {
		panic(NewIllegalStateExceptionMsg("Expected null but was " + r.Peek().String() + r.locationString()))
	}
}

// nextQuotedValue returns the string up to but not including quote,
// unescaping any character escape sequences encountered along the way. The
// opening quote should have already been read. This consumes the closing
// quote, but does not include it in the returned string.
func (r *JsonReader) nextQuotedValue(quote uint16) string {
	// Like nextNonWhitespace, this uses locals 'p' and 'l' to save
	// inner-loop field access.
	buffer := r.buffer
	var builder []uint16
	haveBuilder := false
	for {
		p := r.pos
		l := r.limit
		// the index of the first character not yet appended to the builder.
		start := p
		for p < l {
			c := buffer[p]
			p++

			if c == quote {
				r.pos = p
				length := p - start - 1
				if !haveBuilder {
					return StringFromUTF16(buffer[start : start+length])
				} else {
					builder = jsonReader_appendUint16s(builder, buffer[start:start+length])
					r.builder = builder[:0]
					return StringFromUTF16(builder)
				}
			} else if c == '\\' {
				r.pos = p
				length := p - start - 1
				if !haveBuilder {
					estimatedLength := max((length+1)*2, 16)
					if cap(r.builder) < estimatedLength {
						r.builder = make([]uint16, 0, max(cap(r.builder)*2, max(estimatedLength, 64)))
					}
					builder = r.builder[:0]
					haveBuilder = true
				}
				builder = jsonReader_appendUint16s(builder, buffer[start:start+length])
				builder = jsonReader_appendUint16(builder, r.readEscapeCharacter())
				p = r.pos
				l = r.limit
				start = p
			} else if c == '\n' {
				r.lineNumber++
				r.lineStart = p
			}
		}

		if !haveBuilder {
			if sr, ok := r.in.(*StringReader); ok && !sr.hasPendingLo && sr.next < len(sr.s) {
				bufASCII := true
				for _, u := range buffer[start:p] {
					if u >= 0x80 {
						bufASCII = false
						break
					}
				}
				if bufASCII {
					rem := sr.s[sr.next:]
					if qIdx := strings.IndexByte(rem, byte(quote)); qIdx >= 0 {
						remASCII := true
						for i := 0; i < qIdx; i++ {
							c := rem[i]
							if c >= 0x80 || c == '\\' || c == '\n' {
								remASCII = false
								break
							}
						}
						if remASCII {
							var sb strings.Builder
							sb.Grow((p - start) + qIdx)
							for _, u := range buffer[start:p] {
								sb.WriteByte(byte(u))
							}
							sb.WriteString(rem[:qIdx])
							fullChunks := qIdx / jsonReader_BUFFER_SIZE
							if fullChunks > 0 {
								r.lineStart -= p + (fullChunks-1)*jsonReader_BUFFER_SIZE
								r.pos = jsonReader_BUFFER_SIZE
								r.limit = jsonReader_BUFFER_SIZE
								sr.next += fullChunks * jsonReader_BUFFER_SIZE
								sr.converted = true
							} else {
								r.pos = p
							}
							r.fillBuffer(1)
							r.pos = (qIdx % jsonReader_BUFFER_SIZE) + 1
							return sb.String()
						}
					}
				}
			}
			estimatedLength := max((p-start)*2, 16)
			if cap(r.builder) < estimatedLength {
				r.builder = make([]uint16, 0, max(cap(r.builder)*2, max(estimatedLength, 64)))
			}
			builder = r.builder[:0]
			haveBuilder = true
		}
		builder = jsonReader_appendUint16s(builder, buffer[start:p])
		r.pos = p
		if !r.fillBuffer(1) {
			r.builder = builder[:0]
			panic(r.syntaxError("Unterminated string"))
		}
	}
}

func jsonReader_appendUint16s(builder, chunk []uint16) []uint16 {
	if len(builder)+len(chunk) > cap(builder) {
		newCap := max(cap(builder)*2, max(len(builder)+len(chunk), 64))
		next := make([]uint16, len(builder), newCap)
		copy(next, builder)
		builder = next
	}
	return append(builder, chunk...)
}

func jsonReader_appendUint16(builder []uint16, c uint16) []uint16 {
	if len(builder)+1 > cap(builder) {
		newCap := max(cap(builder)*2, max(len(builder)+1, 64))
		next := make([]uint16, len(builder), newCap)
		copy(next, builder)
		builder = next
	}
	return append(builder, c)
}

// nextUnquotedValue returns an unquoted value as a string.
func (r *JsonReader) nextUnquotedValue() string {
	var builder []uint16
	haveBuilder := false
	i := 0

findNonLiteralCharacter:
	for {
		for ; r.pos+i < r.limit; i++ {
			switch r.buffer[r.pos+i] {
			case '/', '\\', ';', '#', '=':
				r.checkLenient() // fall-through
				break findNonLiteralCharacter
			case '{', '}', '[', ']', ':', ',', ' ', '\t', '\f', '\r', '\n':
				break findNonLiteralCharacter
			}
		}

		// Attempt to load the entire literal into the buffer at once.
		if i < len(r.buffer) {
			if r.fillBuffer(i + 1) {
				continue
			} else {
				break
			}
		}

		// use a StringBuilder when the value is too long. This is too long to
		// be a number!
		if !haveBuilder {
			estimatedLength := max(i, 16)
			if cap(r.builder) < estimatedLength {
				r.builder = make([]uint16, 0, estimatedLength)
			}
			builder = r.builder[:0]
			haveBuilder = true
		}
		builder = append(builder, r.buffer[r.pos:r.pos+i]...)
		r.pos += i
		i = 0
		if !r.fillBuffer(1) {
			break
		}
	}

	var result string
	if !haveBuilder {
		result = StringFromUTF16(r.buffer[r.pos : r.pos+i])
	} else {
		builder = append(builder, r.buffer[r.pos:r.pos+i]...)
		r.builder = builder[:0]
		result = StringFromUTF16(builder)
	}
	r.pos += i
	return result
}

func (r *JsonReader) skipQuotedValue(quote uint16) {
	// Like nextNonWhitespace, this uses locals 'p' and 'l' to save
	// inner-loop field access.
	buffer := r.buffer
	for {
		p := r.pos
		l := r.limit
		// the index of the first character not yet appended to the builder.
		for p < l {
			c := buffer[p]
			p++
			if c == quote {
				r.pos = p
				return
			} else if c == '\\' {
				r.pos = p
				r.readEscapeCharacter()
				p = r.pos
				l = r.limit
			} else if c == '\n' {
				r.lineNumber++
				r.lineStart = p
			}
		}
		r.pos = p
		if !r.fillBuffer(1) {
			break
		}
	}
	panic(r.syntaxError("Unterminated string"))
}

func (r *JsonReader) skipUnquotedValue() {
	for {
		i := 0
		for ; r.pos+i < r.limit; i++ {
			switch r.buffer[r.pos+i] {
			case '/', '\\', ';', '#', '=':
				r.checkLenient() // fall-through
				r.pos += i
				return
			case '{', '}', '[', ']', ':', ',', ' ', '\t', '\f', '\r', '\n':
				r.pos += i
				return
			}
		}
		r.pos += i
		if !r.fillBuffer(1) {
			break
		}
	}
}

// Close closes this JSON reader and the underlying Reader.
//
// Java: JsonReader.close()
func (r *JsonReader) Close() {
	r.peeked = jsonReader_PEEKED_NONE
	r.stack[0] = jsonScope_CLOSED
	r.stackSize = 1
	r.in.Close()
}

// SkipValue skips the next value recursively. If it is an object or array,
// all nested elements are skipped.
//
// Java: JsonReader.skipValue()
func (r *JsonReader) SkipValue() {
	count := 0
	for {
		p := r.peeked
		if p == jsonReader_PEEKED_NONE {
			p = r.doPeek()
		}

		if p == jsonReader_PEEKED_BEGIN_ARRAY {
			r.push(jsonScope_EMPTY_ARRAY)
			count++
		} else if p == jsonReader_PEEKED_BEGIN_OBJECT {
			r.push(jsonScope_EMPTY_OBJECT)
			count++
		} else if p == jsonReader_PEEKED_END_ARRAY {
			r.stackSize--
			count--
		} else if p == jsonReader_PEEKED_END_OBJECT {
			r.stackSize--
			count--
		} else if p == jsonReader_PEEKED_UNQUOTED_NAME || p == jsonReader_PEEKED_UNQUOTED {
			r.skipUnquotedValue()
		} else if p == jsonReader_PEEKED_SINGLE_QUOTED || p == jsonReader_PEEKED_SINGLE_QUOTED_NAME {
			r.skipQuotedValue('\'')
		} else if p == jsonReader_PEEKED_DOUBLE_QUOTED || p == jsonReader_PEEKED_DOUBLE_QUOTED_NAME {
			r.skipQuotedValue('"')
		} else if p == jsonReader_PEEKED_NUMBER {
			r.pos += r.peekedNumberLength
		}
		r.peeked = jsonReader_PEEKED_NONE
		if count == 0 {
			break
		}
	}

	r.pathIndices[r.stackSize-1]++
	null := "null"
	r.pathNames[r.stackSize-1] = &null
}

func (r *JsonReader) push(newTop int) {
	if r.stackSize == len(r.stack) {
		newLength := r.stackSize * 2
		r.stack = append(r.stack, make([]int, newLength-len(r.stack))...)
		r.pathIndices = append(r.pathIndices, make([]int, newLength-len(r.pathIndices))...)
		r.pathNames = append(r.pathNames, make([]*string, newLength-len(r.pathNames))...)
	}
	r.stack[r.stackSize] = newTop
	r.stackSize++
}

// fillBuffer returns true once limit - pos >= minimum. If the data is
// exhausted before that many characters are available, this returns false.
func (r *JsonReader) fillBuffer(minimum int) bool {
	buffer := r.buffer
	r.lineStart -= r.pos
	if r.limit != r.pos {
		r.limit -= r.pos
		copy(buffer[0:r.limit], buffer[r.pos:r.pos+r.limit])
	} else {
		r.limit = 0
	}

	r.pos = 0
	var total int
	for {
		total = r.in.Read(buffer, r.limit, len(buffer)-r.limit)
		if total == -1 {
			break
		}
		r.limit += total

		// if this is the first read, consume an optional byte order mark
		// (BOM) if it exists
		if r.lineNumber == 0 && r.lineStart == 0 && r.limit > 0 && buffer[0] == 0xfeff {
			r.pos++
			r.lineStart++
			minimum++
		}

		if r.limit >= minimum {
			return true
		}
	}
	return false
}

// nextNonWhitespace returns the next character in the stream that is
// neither whitespace nor a part of a comment. When this returns, the
// returned character is always at buffer[pos-1]; this means the caller can
// always push back the returned character by decrementing pos.
func (r *JsonReader) nextNonWhitespace(throwOnEof bool) int {
	// This code uses ugly local variables 'p' and 'l' representing the 'pos'
	// and 'limit' fields respectively. Using locals rather than fields saves
	// a few field reads for each whitespace character in a pretty-printed
	// document, resulting in a 5% speedup. We need to flush 'p' to its field
	// before any (potentially indirect) call to fillBuffer() and reread both
	// 'p' and 'l' after any (potentially indirect) call to the same method.
	buffer := r.buffer
	p := r.pos
	l := r.limit
	for {
		if p == l {
			r.pos = p
			if !r.fillBuffer(1) {
				break
			}
			p = r.pos
			l = r.limit
		}

		c := int(buffer[p])
		p++
		if c == '\n' {
			r.lineNumber++
			r.lineStart = p
			continue
		} else if c == ' ' || c == '\r' || c == '\t' {
			continue
		}

		if c == '/' {
			r.pos = p
			if p == l {
				r.pos-- // push back '/' so it's still in the buffer when this method returns
				charsLoaded := r.fillBuffer(2)
				r.pos++ // consume the '/' again
				if !charsLoaded {
					return c
				}
			}

			r.checkLenient()
			peek := buffer[r.pos]
			switch peek {
			case '*':
				// skip a /* c-style comment */
				r.pos++
				if !r.skipTo("*/") {
					panic(r.syntaxError("Unterminated comment"))
				}
				p = r.pos + 2
				l = r.limit
				continue

			case '/':
				// skip a // end-of-line comment
				r.pos++
				r.skipToEndOfLine()
				p = r.pos
				l = r.limit
				continue

			default:
				return c
			}
		} else if c == '#' {
			r.pos = p
			// Skip a # hash end-of-line comment. The JSON RFC doesn't specify
			// this behaviour, but it's required to parse existing documents.
			// See http://b/2571423.
			r.checkLenient()
			r.skipToEndOfLine()
			p = r.pos
			l = r.limit
		} else {
			r.pos = p
			return c
		}
	}
	if throwOnEof {
		panic(newEOFException("End of input" + r.locationString()))
	} else {
		return -1
	}
}

func (r *JsonReader) checkLenient() {
	if !r.lenient {
		panic(r.syntaxError("Use JsonReader.setLenient(true) to accept malformed JSON"))
	}
}

// skipToEndOfLine advances the position until after the next newline
// character. If the line is terminated by "\r\n", the '\n' must be consumed
// as whitespace by the caller.
func (r *JsonReader) skipToEndOfLine() {
	for r.pos < r.limit || r.fillBuffer(1) {
		c := r.buffer[r.pos]
		r.pos++
		if c == '\n' {
			r.lineNumber++
			r.lineStart = r.pos
			break
		} else if c == '\r' {
			break
		}
	}
}

// skipTo advances to the beginning of toFind (ASCII).
func (r *JsonReader) skipTo(toFind string) bool {
	length := len(toFind)
outer:
	for ; r.pos+length <= r.limit || r.fillBuffer(length); r.pos++ {
		if r.buffer[r.pos] == '\n' {
			r.lineNumber++
			r.lineStart = r.pos + 1
			continue
		}
		for c := 0; c < length; c++ {
			if r.buffer[r.pos+c] != uint16(toFind[c]) {
				continue outer
			}
		}
		return true
	}
	return false
}

// Java: JsonReader.toString()
func (r *JsonReader) String() string {
	return "JsonReader" + r.locationString()
}

func (r *JsonReader) locationString() string {
	line := r.lineNumber + 1
	column := r.pos - r.lineStart + 1
	return " at line " + itoa(line) + " column " + itoa(column) + " path " + r.GetPath()
}

func (r *JsonReader) getPath(usePreviousPath bool) string {
	var result strings.Builder
	result.WriteByte('$')
	for i := 0; i < r.stackSize; i++ {
		switch r.stack[i] {
		case jsonScope_EMPTY_ARRAY, jsonScope_NONEMPTY_ARRAY:
			pathIndex := r.pathIndices[i]
			// If index is last path element it points to next array element;
			// have to decrement
			if usePreviousPath && pathIndex > 0 && i == r.stackSize-1 {
				pathIndex--
			}
			result.WriteByte('[')
			result.WriteString(itoa(pathIndex))
			result.WriteByte(']')
		case jsonScope_EMPTY_OBJECT, jsonScope_DANGLING_NAME, jsonScope_NONEMPTY_OBJECT:
			result.WriteByte('.')
			if r.pathNames[i] != nil {
				result.WriteString(*r.pathNames[i])
			}
		case jsonScope_NONEMPTY_DOCUMENT, jsonScope_EMPTY_DOCUMENT, jsonScope_CLOSED:
		}
	}
	return result.String()
}

// getPreviousPath returns a JsonPath to the element just read.
//
// Java: JsonReader.getPreviousPath()
func (r *JsonReader) getPreviousPath() string {
	return r.getPath(true)
}

// GetPath returns a JsonPath to the next (to be read) element.
//
// Java: JsonReader.getPath()
func (r *JsonReader) GetPath() string {
	return r.getPath(false)
}

// readEscapeCharacter unescapes the character identified by the
// character or characters that immediately follow a backslash. The backslash
// '\' should have already been read. This supports both unicode escapes "u000A"
// and two-character escapes "\n".
func (r *JsonReader) readEscapeCharacter() uint16 {
	if r.pos == r.limit && !r.fillBuffer(1) {
		panic(r.syntaxError("Unterminated escape sequence"))
	}

	escaped := r.buffer[r.pos]
	r.pos++
	switch escaped {
	case 'u':
		if r.pos+4 > r.limit && !r.fillBuffer(4) {
			panic(r.syntaxError("Unterminated escape sequence"))
		}
		// Equivalent to Integer.parseInt(stringPool.get(buffer, pos, 4), 16);
		var result uint16
		for i, end := r.pos, r.pos+4; i < end; i++ {
			c := r.buffer[i]
			result <<= 4
			if c >= '0' && c <= '9' {
				result += c - '0'
			} else if c >= 'a' && c <= 'f' {
				result += c - 'a' + 10
			} else if c >= 'A' && c <= 'F' {
				result += c - 'A' + 10
			} else {
				panic(NewNumberFormatException("\\u" + StringFromUTF16(r.buffer[r.pos:r.pos+4])))
			}
		}
		r.pos += 4
		return result

	case 't':
		return '\t'

	case 'b':
		return '\b'

	case 'n':
		return '\n'

	case 'r':
		return '\r'

	case 'f':
		return '\f'

	case '\n':
		r.lineNumber++
		r.lineStart = r.pos
		// fall-through
		return escaped

	case '\'', '"', '\\', '/':
		return escaped
	default:
		// throw error when none of the above cases are matched
		panic(r.syntaxError("Invalid escape sequence"))
	}
}

// syntaxError builds the MalformedJsonException with the given message and
// a context snippet with this reader's content; the caller throws it.
func (r *JsonReader) syntaxError(message string) *malformedJsonException {
	return newMalformedJsonException(message + r.locationString())
}

// consumeNonExecutePrefix consumes the non-execute prefix if it exists.
func (r *JsonReader) consumeNonExecutePrefix() {
	// fast forward through the leading whitespace
	r.nextNonWhitespace(true)
	r.pos--

	if r.pos+5 > r.limit && !r.fillBuffer(5) {
		return
	}

	p := r.pos
	buf := r.buffer
	if buf[p] != ')' || buf[p+1] != ']' || buf[p+2] != '}' || buf[p+3] != '\'' || buf[p+4] != '\n' {
		return // not a security token!
	}

	// we consumed a security token!
	r.pos += 5
}
