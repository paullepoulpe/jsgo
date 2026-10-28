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

import (
	"math"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// SourceMapJsonLexer is a JSON lexer optimized for large string extraction.
// By using String.substring() on the source contents, it avoids the massive
// GC churn caused by JsonReader's internal StringBuilder resizing when
// parsing large Base64 VLQ mapping strings.
//
// Internal callers use the (..., error) methods directly; exported methods
// preserve the SourceMapParseException panic contract at the public boundary.
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
	units := jutil.StringToUTF16(json)
	return &SourceMapJsonLexer{
		json:   units,
		length: len(units),
	}
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

func (l *SourceMapJsonLexer) beginObject() error {
	l.SkipWhitespace()
	if l.pos < l.length && l.json[l.pos] == '{' {
		l.pos++
		return nil
	}
	return NewSourceMapParseException("Expected '{'")
}

// Java: SourceMapJsonLexer.beginObject()
func (l *SourceMapJsonLexer) BeginObject() {
	if err := l.beginObject(); err != nil {
		panic(err)
	}
}

func (l *SourceMapJsonLexer) endObject() error {
	l.SkipWhitespace()
	if l.pos < l.length && l.json[l.pos] == '}' {
		l.pos++
		return nil
	}
	return NewSourceMapParseException("Expected '}'")
}

// Java: SourceMapJsonLexer.endObject()
func (l *SourceMapJsonLexer) EndObject() {
	if err := l.endObject(); err != nil {
		panic(err)
	}
}

func (l *SourceMapJsonLexer) beginArray() error {
	l.SkipWhitespace()
	if l.pos < l.length && l.json[l.pos] == '[' {
		l.pos++
		return nil
	}
	return NewSourceMapParseException("Expected '['")
}

// Java: SourceMapJsonLexer.beginArray()
func (l *SourceMapJsonLexer) BeginArray() {
	if err := l.beginArray(); err != nil {
		panic(err)
	}
}

func (l *SourceMapJsonLexer) endArray() error {
	l.SkipWhitespace()
	if l.pos < l.length && l.json[l.pos] == ']' {
		l.pos++
		return nil
	}
	return NewSourceMapParseException("Expected ']'")
}

// Java: SourceMapJsonLexer.endArray()
func (l *SourceMapJsonLexer) EndArray() {
	if err := l.endArray(); err != nil {
		panic(err)
	}
}

func (l *SourceMapJsonLexer) nextName() (string, error) {
	name, err := l.nextString()
	if err != nil {
		return "", err
	}
	l.SkipWhitespace()
	if l.pos < l.length && l.json[l.pos] == ':' {
		l.pos++
		return name, nil
	}
	return "", NewSourceMapParseException("Expected ':'")
}

// Java: SourceMapJsonLexer.nextName()
func (l *SourceMapJsonLexer) NextName() string {
	name, err := l.nextName()
	if err != nil {
		panic(err)
	}
	return name
}

func (l *SourceMapJsonLexer) nextStringOrNull() (*string, error) {
	l.SkipWhitespace()
	if l.pos < l.length && l.startsWith("null", l.pos) {
		l.pos += 4
		return nil, nil
	}
	s, err := l.nextString()
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// NextStringOrNull returns the next string, or nil for a null literal.
//
// Java: SourceMapJsonLexer.nextStringOrNull()
func (l *SourceMapJsonLexer) NextStringOrNull() *string {
	s, err := l.nextStringOrNull()
	if err != nil {
		panic(err)
	}
	return s
}

func (l *SourceMapJsonLexer) nextString() (string, error) {
	l.SkipWhitespace()
	if l.pos >= l.length || l.json[l.pos] != '"' {
		return "", NewSourceMapParseException("Expected string")
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
		return "", NewSourceMapParseException("Unterminated string")
	}
	var result string
	if !hasEscape {
		result = l.substring(start, l.pos)
	} else {
		var err error
		result, err = sourceMapJsonLexer_unescape(l.json, start, l.pos)
		if err != nil {
			return "", err
		}
	}
	l.pos++ // skip quote
	return result, nil
}

// Java: SourceMapJsonLexer.nextString()
func (l *SourceMapJsonLexer) NextString() string {
	s, err := l.nextString()
	if err != nil {
		panic(err)
	}
	return s
}

func (l *SourceMapJsonLexer) nextInt() (int, error) {
	l.SkipWhitespace()
	start := l.pos
	if l.pos < l.length && l.json[l.pos] == '-' {
		l.pos++
	}
	for l.pos < l.length && l.json[l.pos] >= '0' && l.json[l.pos] <= '9' {
		l.pos++
	}
	if start == l.pos {
		return 0, NewSourceMapParseException("Expected integer")
	}
	result, ok := sourceMapJsonLexer_parseInt(l.json, start, l.pos, 10)
	if !ok {
		return 0, NewSourceMapParseException("Invalid integer: " + l.substring(start, l.pos))
	}
	return int(result), nil
}

// Java: SourceMapJsonLexer.nextInt()
func (l *SourceMapJsonLexer) NextInt() int {
	v, err := l.nextInt()
	if err != nil {
		panic(err)
	}
	return v
}

func (l *SourceMapJsonLexer) nextRawValue() (string, error) {
	l.SkipWhitespace()
	start := l.pos
	if err := l.skipValue(); err != nil {
		return "", err
	}
	return l.substring(start, l.pos), nil
}

// Java: SourceMapJsonLexer.nextRawValue()
func (l *SourceMapJsonLexer) NextRawValue() string {
	s, err := l.nextRawValue()
	if err != nil {
		panic(err)
	}
	return s
}

func (l *SourceMapJsonLexer) skipValue() error {
	l.SkipWhitespace()
	if l.pos >= l.length {
		return NewSourceMapParseException("Unexpected end of input")
	}
	c := l.json[l.pos]
	if c == '{' {
		l.pos++
		for l.HasNext() {
			if _, err := l.nextName(); err != nil {
				return err
			}
			if err := l.skipValue(); err != nil {
				return err
			}
			if err := l.checkComma(); err != nil {
				return err
			}
		}
		return l.endObject()
	} else if c == '[' {
		l.pos++
		for l.HasNext() {
			if err := l.skipValue(); err != nil {
				return err
			}
			if err := l.checkComma(); err != nil {
				return err
			}
		}
		return l.endArray()
	} else if c == '"' {
		_, err := l.nextString()
		return err
	} else if c == 't' && l.startsWith("true", l.pos) {
		l.pos += 4
		return nil
	} else if c == 'f' && l.startsWith("false", l.pos) {
		l.pos += 5
		return nil
	} else if c == 'n' && l.startsWith("null", l.pos) {
		l.pos += 4
		return nil
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
		return nil
	}
	return NewSourceMapParseException("Unexpected character: " + jutil.StringFromChar(c))
}

// Java: SourceMapJsonLexer.skipValue()
func (l *SourceMapJsonLexer) SkipValue() {
	if err := l.skipValue(); err != nil {
		panic(err)
	}
}

func (l *SourceMapJsonLexer) checkComma() error {
	hasComma := l.TryConsumeComma()
	if hasComma && !l.HasNext() {
		return NewSourceMapParseException("Unexpected trailing comma")
	}
	if !hasComma && l.HasNext() {
		return NewSourceMapParseException("Expected ','")
	}
	return nil
}

// Java: SourceMapJsonLexer.checkComma()
func (l *SourceMapJsonLexer) CheckComma() {
	if err := l.checkComma(); err != nil {
		panic(err)
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

func sourceMapJsonLexer_unescape(s []uint16, start, end int) (string, error) {
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
					return "", NewSourceMapParseException("Invalid unicode escape sequence")
				}
				v, ok := sourceMapJsonLexer_parseInt(s, i+1, i+5, 16)
				if !ok {
					return "", NewSourceMapParseException(
						"Invalid unicode escape: " + jutil.StringFromUTF16(s[i+1:i+5]))
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
	return jutil.StringFromUTF16(sb), nil
}

// SourceMapJsonLexer_Unescape decodes the JSON escapes in s[start:end].
//
// Java: SourceMapJsonLexer.unescape(String, int, int)
func SourceMapJsonLexer_Unescape(s []uint16, start, end int) string {
	out, err := sourceMapJsonLexer_unescape(s, start, end)
	if err != nil {
		panic(err)
	}
	return out
}

// sourceMapJsonLexer_parseInt is Integer.parseInt(CharSequence, int, int,
// int) on s[beginIndex:endIndex] (for radix 10 or 16); ok is false where Java
// throws NumberFormatException.
func sourceMapJsonLexer_parseInt(s []uint16, beginIndex, endIndex, radix int) (int32, bool) {
	i := beginIndex
	neg := false
	if i < endIndex && (s[i] == '-' || s[i] == '+') {
		neg = s[i] == '-'
		i++
	}
	if i == endIndex {
		return 0, false
	}
	limit := uint64(math.MaxInt32)
	if neg {
		limit++
	}
	base := uint64(radix)
	var mag uint64
	for ; i < endIndex; i++ {
		c := s[i]
		var d int
		if c >= '0' && c <= '9' {
			d = int(c - '0')
		} else if radix == 16 && c >= 'a' && c <= 'f' {
			d = int(c - 'a' + 10)
		} else if radix == 16 && c >= 'A' && c <= 'F' {
			d = int(c - 'A' + 10)
		} else if c < 0x80 {
			return 0, false
		} else {
			d = jutil.CharacterDigit(c, radix)
			if d < 0 {
				return 0, false
			}
		}
		if mag > (limit-uint64(d))/base {
			return 0, false
		}
		mag = mag*base + uint64(d)
	}
	if neg {
		return int32(-int64(mag)), true
	}
	return int32(mag), true
}
