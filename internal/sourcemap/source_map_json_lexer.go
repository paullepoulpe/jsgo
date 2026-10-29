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
	"strings"
	"unicode/utf16"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// SourceMapJsonLexer is a JSON lexer optimized for large string extraction.
// By slicing the UTF-8 source contents directly, it avoids allocating UTF-16
// buffers or copying unescaped Base64 VLQ mapping strings.
//
// Internal callers use the (..., error) methods directly; exported methods
// preserve the SourceMapParseException panic contract at the public boundary.
type SourceMapJsonLexer struct {
	json   string
	pos    int
	length int
}

// NewSourceMapJsonLexer constructs a lexer over json.
//
// Java: SourceMapJsonLexer(String)
func NewSourceMapJsonLexer(json string) *SourceMapJsonLexer {
	return &SourceMapJsonLexer{
		json:   json,
		length: len(json),
	}
}

// startsWith is json.startsWith(prefix, toffset) for an ASCII prefix.
func (l *SourceMapJsonLexer) startsWith(prefix string, toffset int) bool {
	if toffset < 0 || toffset > l.length-len(prefix) {
		return false
	}
	return l.json[toffset:toffset+len(prefix)] == prefix
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

func (l *SourceMapJsonLexer) endObject() error {
	l.SkipWhitespace()
	if l.pos < l.length && l.json[l.pos] == '}' {
		l.pos++
		return nil
	}
	return NewSourceMapParseException("Expected '}'")
}

func (l *SourceMapJsonLexer) beginArray() error {
	l.SkipWhitespace()
	if l.pos < l.length && l.json[l.pos] == '[' {
		l.pos++
		return nil
	}
	return NewSourceMapParseException("Expected '['")
}

func (l *SourceMapJsonLexer) endArray() error {
	l.SkipWhitespace()
	if l.pos < l.length && l.json[l.pos] == ']' {
		l.pos++
		return nil
	}
	return NewSourceMapParseException("Expected ']'")
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
		result = l.json[start:l.pos]
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
		return 0, NewSourceMapParseException("Invalid integer: " + l.json[start:l.pos])
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

// Java: SourceMapJsonLexer.tryConsumeComma()
func (l *SourceMapJsonLexer) TryConsumeComma() bool {
	l.SkipWhitespace()
	if l.pos < l.length && l.json[l.pos] == ',' {
		l.pos++
		return true
	}
	return false
}

func parseHex4ASCII(s string) (uint16, bool) {
	if len(s) < 4 {
		return 0, false
	}
	var v uint16
	for i := 0; i < 4; i++ {
		c := s[i]
		var d uint16
		switch {
		case c >= '0' && c <= '9':
			d = uint16(c - '0')
		case c >= 'a' && c <= 'f':
			d = uint16(c - 'a' + 10)
		case c >= 'A' && c <= 'F':
			d = uint16(c - 'A' + 10)
		default:
			return 0, false
		}
		v = (v << 4) | d
	}
	return v, true
}

func sourceMapJsonLexer_unescape(s string, start, end int) (string, error) {
	raw := s[start:end]
	var sb strings.Builder
	sb.Grow(len(raw))
	segStart := 0
	for i := 0; i < len(raw); {
		c := raw[i]
		if c == 0xED && i+2 < len(raw) && raw[i+1] >= 0xA0 && raw[i+1] <= 0xBF {
			return sourceMapJsonLexer_unescapeUTF16(jutil.StringToUTF16(raw))
		}
		if c != '\\' {
			i++
			continue
		}
		if i+1 >= len(raw) {
			i++
			continue
		}
		next := raw[i+1]
		if next >= 0x80 {
			return sourceMapJsonLexer_unescapeUTF16(jutil.StringToUTF16(raw))
		}
		if segStart < i {
			sb.WriteString(raw[segStart:i])
		}
		switch next {
		case '"':
			sb.WriteByte('"')
			i += 2
		case '\\':
			sb.WriteByte('\\')
			i += 2
		case '/':
			sb.WriteByte('/')
			i += 2
		case 'b':
			sb.WriteByte('\b')
			i += 2
		case 'f':
			sb.WriteByte('\f')
			i += 2
		case 'n':
			sb.WriteByte('\n')
			i += 2
		case 'r':
			sb.WriteByte('\r')
			i += 2
		case 't':
			sb.WriteByte('\t')
			i += 2
		case 'u':
			if i+6 > len(raw) {
				return sourceMapJsonLexer_unescapeUTF16(jutil.StringToUTF16(raw))
			}
			v, ok := parseHex4ASCII(raw[i+2 : i+6])
			if !ok {
				return sourceMapJsonLexer_unescapeUTF16(jutil.StringToUTF16(raw))
			}
			if v < 0xD800 || v > 0xDFFF {
				sb.WriteRune(rune(v))
				i += 6
			} else if v <= 0xDBFF && i+12 <= len(raw) && raw[i+6] == '\\' && raw[i+7] == 'u' {
				lo, okLo := parseHex4ASCII(raw[i+8 : i+12])
				if !okLo || lo < 0xDC00 || lo > 0xDFFF {
					return sourceMapJsonLexer_unescapeUTF16(jutil.StringToUTF16(raw))
				}
				sb.WriteRune(utf16.DecodeRune(rune(v), rune(lo)))
				i += 12
			} else {
				return sourceMapJsonLexer_unescapeUTF16(jutil.StringToUTF16(raw))
			}
		default:
			sb.WriteByte(next)
			i += 2
		}
		segStart = i
	}
	if segStart < len(raw) {
		sb.WriteString(raw[segStart:])
	}
	return sb.String(), nil
}

func sourceMapJsonLexer_unescapeUTF16(s []uint16) (string, error) {
	end := len(s)
	sb := make([]uint16, 0, end)
	for i := 0; i < end; i++ {
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
				esc := jutil.StringFromUTF16(s[i+1 : i+5])
				v, ok := sourceMapJsonLexer_parseInt(esc, 0, len(esc), 16)
				if !ok {
					return "", NewSourceMapParseException("Invalid unicode escape: " + esc)
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

// sourceMapJsonLexer_parseInt is Integer.parseInt(CharSequence, int, int,
// int) on s[beginIndex:endIndex] (for radix 10 or 16); ok is false where Java
// throws NumberFormatException.
func sourceMapJsonLexer_parseInt(s string, beginIndex, endIndex, radix int) (int32, bool) {
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
	for i < endIndex {
		c := s[i]
		var d int
		if c >= '0' && c <= '9' {
			d = int(c - '0')
			i++
		} else if radix == 16 && c >= 'a' && c <= 'f' {
			d = int(c - 'a' + 10)
			i++
		} else if radix == 16 && c >= 'A' && c <= 'F' {
			d = int(c - 'A' + 10)
			i++
		} else if c < 0x80 {
			return 0, false
		} else {
			r, width := util_decodeWTF8(s[i:endIndex])
			if r >= 0x10000 {
				return 0, false
			}
			d = jutil.CharacterDigit(uint16(r), radix)
			if d < 0 {
				return 0, false
			}
			i += width
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
