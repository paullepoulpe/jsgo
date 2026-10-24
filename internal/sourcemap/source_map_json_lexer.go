// Port of com/google/debugging/sourcemap/SourceMapJsonLexer.java @ closure-compiler v20261021

/*
 * Copyright 2026 The Closure Compiler Authors.
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

package sourcemap

import "github.com/paullepoulpe/jsgo/internal/jutil"

// SourceMapJsonLexer is a JSON lexer optimized for large string extraction.
// By using String.substring() on the source contents, it avoids the massive
// GC churn caused by JsonReader's internal StringBuilder resizing when
// parsing large Base64 VLQ mapping strings.
//
// Every method that throws SourceMapParseException in Java panics with a
// *SourceMapParseException.
type SourceMapJsonLexer struct {
	// Go-only: language-adapter — the Java String json as UTF-16 code units,
	// so that charAt is O(1) and pos is a Java char index.
	json   []uint16
	pos    int
	length int
}

// NewSourceMapJsonLexer constructs a lexer over json.
//
// Java: SourceMapJsonLexer(String)
func NewSourceMapJsonLexer(json string) *SourceMapJsonLexer {
	l := &SourceMapJsonLexer{}
	l.json = jutil.StringToUTF16(json)
	l.length = len(l.json)
	return l
}

// startsWith is json.startsWith(prefix, toffset) for an ASCII prefix.
func (l *SourceMapJsonLexer) startsWith(prefix string, toffset int) bool {
	if toffset < 0 || toffset > l.length-len(prefix) {
		return false
	}
	for i := 0; i < len(prefix); i++ {
		if l.json[toffset+i] != uint16(prefix[i]) {
			return false
		}
	}
	return true
}

// substring is json.substring(beginIndex, endIndex).
func (l *SourceMapJsonLexer) substring(beginIndex, endIndex int) string {
	if beginIndex < 0 || endIndex > l.length || beginIndex > endIndex {
		panic(jutil.NewStringIndexOutOfBoundsException(beginIndex, endIndex))
	}
	return jutil.StringFromUTF16(l.json[beginIndex:endIndex])
}

// Java: SourceMapJsonLexer.skipWhitespace()
func (l *SourceMapJsonLexer) SkipWhitespace() {
	for l.pos < l.length {
		c := l.json[l.pos]
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			l.pos++
		} else {
			break
		}
	}
}

// Java: SourceMapJsonLexer.hasNext()
func (l *SourceMapJsonLexer) HasNext() bool {
	l.SkipWhitespace()
	if l.pos >= l.length {
		return false
	}
	c := l.json[l.pos]
	return c != ']' && c != '}'
}

// Java: SourceMapJsonLexer.beginObject()
func (l *SourceMapJsonLexer) BeginObject() {
	l.SkipWhitespace()
	if l.pos < l.length && l.json[l.pos] == '{' {
		l.pos++
	} else {
		panic(NewSourceMapParseException("Expected '{'"))
	}
}

// Java: SourceMapJsonLexer.endObject()
func (l *SourceMapJsonLexer) EndObject() {
	l.SkipWhitespace()
	if l.pos < l.length && l.json[l.pos] == '}' {
		l.pos++
	} else {
		panic(NewSourceMapParseException("Expected '}'"))
	}
}

// Java: SourceMapJsonLexer.beginArray()
func (l *SourceMapJsonLexer) BeginArray() {
	l.SkipWhitespace()
	if l.pos < l.length && l.json[l.pos] == '[' {
		l.pos++
	} else {
		panic(NewSourceMapParseException("Expected '['"))
	}
}

// Java: SourceMapJsonLexer.endArray()
func (l *SourceMapJsonLexer) EndArray() {
	l.SkipWhitespace()
	if l.pos < l.length && l.json[l.pos] == ']' {
		l.pos++
	} else {
		panic(NewSourceMapParseException("Expected ']'"))
	}
}

// Java: SourceMapJsonLexer.nextName()
func (l *SourceMapJsonLexer) NextName() string {
	name := l.NextString()
	l.SkipWhitespace()
	if l.pos < l.length && l.json[l.pos] == ':' {
		l.pos++
	} else {
		panic(NewSourceMapParseException("Expected ':'"))
	}
	return name
}

// NextStringOrNull returns the next string, or nil for a null literal.
//
// Java: SourceMapJsonLexer.nextStringOrNull()
func (l *SourceMapJsonLexer) NextStringOrNull() *string {
	l.SkipWhitespace()
	if l.pos < l.length && l.startsWith("null", l.pos) {
		l.pos += 4
		return nil
	}
	s := l.NextString()
	return &s
}

// Java: SourceMapJsonLexer.nextString()
func (l *SourceMapJsonLexer) NextString() string {
	l.SkipWhitespace()
	if l.pos >= l.length || l.json[l.pos] != '"' {
		panic(NewSourceMapParseException("Expected string"))
	}
	l.pos++ // skip quote
	start := l.pos
	hasEscape := false
	for l.pos < l.length {
		c := l.json[l.pos]
		if c == '"' {
			break
		}
		if c == '\\' {
			hasEscape = true
			l.pos++
		}
		l.pos++
	}
	if l.pos >= l.length {
		panic(NewSourceMapParseException("Unterminated string"))
	}
	var result string
	if !hasEscape {
		result = l.substring(start, l.pos)
	} else {
		result = SourceMapJsonLexer_Unescape(l.json, start, l.pos)
	}
	l.pos++ // skip quote
	return result
}

// Java: SourceMapJsonLexer.nextInt()
func (l *SourceMapJsonLexer) NextInt() int {
	l.SkipWhitespace()
	start := l.pos
	if l.pos < l.length && l.json[l.pos] == '-' {
		l.pos++
	}
	for l.pos < l.length && l.json[l.pos] >= '0' && l.json[l.pos] <= '9' {
		l.pos++
	}
	if start == l.pos {
		panic(NewSourceMapParseException("Expected integer"))
	}
	result, ok := sourceMapJsonLexer_parseInt(l.json, start, l.pos, 10)
	if !ok {
		// catch (NumberFormatException e)
		panic(NewSourceMapParseException("Invalid integer: " + l.substring(start, l.pos)))
	}
	return int(result)
}

// Java: SourceMapJsonLexer.nextRawValue()
func (l *SourceMapJsonLexer) NextRawValue() string {
	l.SkipWhitespace()
	start := l.pos
	l.SkipValue()
	return l.substring(start, l.pos)
}

// Java: SourceMapJsonLexer.skipValue()
func (l *SourceMapJsonLexer) SkipValue() {
	l.SkipWhitespace()
	if l.pos >= l.length {
		panic(NewSourceMapParseException("Unexpected end of input"))
	}
	c := l.json[l.pos]
	if c == '{' {
		l.pos++
		for l.HasNext() {
			l.NextName()
			l.SkipValue()
			l.CheckComma()
		}
		l.EndObject()
	} else if c == '[' {
		l.pos++
		for l.HasNext() {
			l.SkipValue()
			l.CheckComma()
		}
		l.EndArray()
	} else if c == '"' {
		l.NextString()
	} else if c == 't' && l.startsWith("true", l.pos) {
		l.pos += 4
	} else if c == 'f' && l.startsWith("false", l.pos) {
		l.pos += 5
	} else if c == 'n' && l.startsWith("null", l.pos) {
		l.pos += 4
	} else if c == '-' || (c >= '0' && c <= '9') {
		if c == '-' {
			l.pos++
		}
		for l.pos < l.length &&
			((l.json[l.pos] >= '0' && l.json[l.pos] <= '9') ||
				l.json[l.pos] == '.' ||
				l.json[l.pos] == 'e' ||
				l.json[l.pos] == 'E' ||
				l.json[l.pos] == '+' ||
				l.json[l.pos] == '-') {
			l.pos++
		}
	} else {
		panic(NewSourceMapParseException("Unexpected character: " + jutil.StringFromChar(c)))
	}
}

// Java: SourceMapJsonLexer.checkComma()
func (l *SourceMapJsonLexer) CheckComma() {
	hasComma := l.TryConsumeComma()
	if hasComma && !l.HasNext() {
		panic(NewSourceMapParseException("Unexpected trailing comma"))
	}
	if !hasComma && l.HasNext() {
		panic(NewSourceMapParseException("Expected ','"))
	}
}

// Java: SourceMapJsonLexer.tryConsumeComma()
func (l *SourceMapJsonLexer) TryConsumeComma() bool {
	l.SkipWhitespace()
	if l.pos < l.length && l.json[l.pos] == ',' {
		l.pos++
		return true
	}
	return false
}

// SourceMapJsonLexer_Unescape decodes the JSON escapes in s[start:end].
//
// Java: SourceMapJsonLexer.unescape(String, int, int)
func SourceMapJsonLexer_Unescape(s []uint16, start, end int) string {
	sb := make([]uint16, 0, end-start)
	for i := start; i < end; i++ {
		c := s[i]
		if c == '\\' && i+1 < end {
			i++
			next := s[i]
			switch next {
			case '"':
				sb = append(sb, '"')
			case '\\':
				sb = append(sb, '\\')
			case '/':
				sb = append(sb, '/')
			case 'b':
				sb = append(sb, '\b')
			case 'f':
				sb = append(sb, '\f')
			case 'n':
				sb = append(sb, '\n')
			case 'r':
				sb = append(sb, '\r')
			case 't':
				sb = append(sb, '\t')
			case 'u':
				if i+4 >= end {
					panic(NewSourceMapParseException("Invalid unicode escape sequence"))
				}
				v, ok := sourceMapJsonLexer_parseInt(s, i+1, i+5, 16)
				if !ok {
					// catch (NumberFormatException e)
					panic(NewSourceMapParseException(
						"Invalid unicode escape: " + jutil.StringFromUTF16(s[i+1:i+5])))
				}
				sb = append(sb, uint16(v))
				i += 4
			default:
				sb = append(sb, next)
			}
		} else {
			sb = append(sb, c)
		}
	}
	return jutil.StringFromUTF16(sb)
}

// sourceMapJsonLexer_parseInt is Integer.parseInt(CharSequence, int, int,
// int) on s[beginIndex:endIndex]; ok is false where Java throws
// NumberFormatException. By its Javadoc it parses that range exactly as
// Integer.parseInt(String, int) parses the same chars, and the callers
// replace the exception's message with their own, so jutil.IntegerParseInt
// on the substring is equivalent. (The callers only pass in-range indices.)
func sourceMapJsonLexer_parseInt(s []uint16, beginIndex, endIndex, radix int) (v int32, ok bool) {
	jutil.Try(func() {
		v, ok = jutil.IntegerParseInt(jutil.StringFromUTF16(s[beginIndex:endIndex]), radix), true
	}, func(e jutil.INumberFormatException) {})
	return v, ok
}
