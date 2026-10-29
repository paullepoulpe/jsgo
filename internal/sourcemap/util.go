// Port of com/google/debugging/sourcemap/Util.java @ closure-compiler v20261021

/*
 * Copyright 2011 The Closure Compiler Authors.
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
	"io"
	"strings"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Util holds the source map string escaper.
type Util struct{}

// Util_HEX_CHARS is the hex digit table.
var Util_HEX_CHARS = [...]uint16{'0', '1', '2', '3', '4', '5', '6', '7',
	'8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}

const util_hexDigits = "0123456789abcdef"

var util_endScriptUnits = [...]uint16{'/', 's', 'c', 'r', 'i', 'p', 't'}

func util_needsEscape(s string) bool {
	return util_firstEscape(s) >= 0
}

func util_firstEscape(s string) int {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c <= 0x1f || c > 0x7f || c == '\\' || c == '"' || c == '>' || c == '<' {
			return i
		}
	}
	return -1
}

// util_decodeWTF8 decodes one code point (including lone WTF-8 surrogates)
// identically to jutil.StringToUTF16.
func util_decodeWTF8(s string) (rune, int) {
	c := s[0]
	if c == 0xED && len(s) >= 3 && s[1] >= 0xA0 && s[1] <= 0xBF && s[2]&0xC0 == 0x80 {
		return rune(c&0x0F)<<12 | rune(s[1]&0x3F)<<6 | rune(s[2]&0x3F), 3
	}
	return utf8.DecodeRuneInString(s)
}

func util_matchesEndScript(rest string) bool {
	if len(rest) < 7 || rest[0] != '/' {
		return false
	}
	if (rest[1] | rest[2] | rest[3] | rest[4] | rest[5] | rest[6]) < 0x80 {
		return (rest[1]|0x20) == 's' &&
			(rest[2]|0x20) == 'c' &&
			(rest[3]|0x20) == 'r' &&
			(rest[4]|0x20) == 'i' &&
			(rest[5]|0x20) == 'p' &&
			(rest[6]|0x20) == 't'
	}
	var units [7]uint16
	units[0] = '/'
	n := 1
	pos := 1
	for n < 7 && pos < len(rest) {
		var r rune
		var w int
		if rest[pos] < 0x80 {
			r, w = rune(rest[pos]), 1
		} else {
			r, w = util_decodeWTF8(rest[pos:])
		}
		pos += w
		if r >= 0x10000 {
			hi, lo := utf16.EncodeRune(r)
			units[n] = uint16(hi)
			n++
			if n < 7 {
				units[n] = uint16(lo)
				n++
			}
		} else {
			units[n] = uint16(r)
			n++
		}
	}
	if n < 7 {
		return false
	}
	return jutil.UTF16RegionMatches(true, units[:], 0, util_endScriptUnits[:], 0, 7)
}

func util_putHexUnit(dst []byte, c uint16) {
	_ = dst[5]
	dst[0] = '\\'
	dst[1] = 'u'
	dst[2] = util_hexDigits[(c>>12)&0xf]
	dst[3] = util_hexDigits[(c>>8)&0xf]
	dst[4] = util_hexDigits[(c>>4)&0xf]
	dst[5] = util_hexDigits[c&0xf]
}

// Util_EscapeString escapes the given string to a double quoted (")
// JavaScript/JSON string.
//
// Java: Util.escapeString(String)
func Util_EscapeString(s string) string {
	first := util_firstEscape(s)
	if first < 0 {
		return `"` + s + `"`
	}

	var sb strings.Builder
	sb.Grow(len(s) + 16)
	sb.WriteByte('"')

	var hexBuf [12]byte
	start := 0
	for i := first; i < len(s); {
		c := s[i]
		if c >= 0x80 {
			if start < i {
				sb.WriteString(s[start:i])
			}
			r, width := util_decodeWTF8(s[i:])
			if r >= 0x10000 {
				hi, lo := utf16.EncodeRune(r)
				util_putHexUnit(hexBuf[0:6], uint16(hi))
				util_putHexUnit(hexBuf[6:12], uint16(lo))
				sb.Write(hexBuf[:12])
			} else {
				util_putHexUnit(hexBuf[0:6], uint16(r))
				sb.Write(hexBuf[:6])
			}
			i += width
			start = i
			continue
		}

		var esc string
		switch c {
		case '\n':
			esc = `\n`
		case '\r':
			esc = `\r`
		case '\t':
			esc = `\t`
		case '\\':
			esc = `\\`
		case '"':
			esc = `\"`
		case '>':
			// Unicode-escape the '>' in '-->' and ']]>'
			if i >= 2 &&
				((s[i-1] == '-' && s[i-2] == '-') ||
					(s[i-1] == ']' && s[i-2] == ']')) {
				esc = `\u003e`
			} else {
				i++
				continue
			}
		case '<':
			// Unicode-escape the '<' in '</script' and '<!--'
			if util_matchesEndScript(s[i+1:]) || strings.HasPrefix(s[i+1:], "!--") {
				esc = `\u003c`
			} else {
				i++
				continue
			}
		default:
			if c > 0x1f {
				i++
				continue
			}
			if start < i {
				sb.WriteString(s[start:i])
			}
			util_putHexUnit(hexBuf[0:6], uint16(c))
			sb.Write(hexBuf[:6])
			i++
			start = i
			continue
		}

		if start < i {
			sb.WriteString(s[start:i])
		}
		sb.WriteString(esc)
		i++
		start = i
	}
	if start < len(s) {
		sb.WriteString(s[start:])
	}
	sb.WriteByte('"')
	return sb.String()
}

// Util_AppendHexJavaScriptRepresentation appends the hex escape of c.
//
// Java: Util.appendHexJavaScriptRepresentation(StringBuilder, char)
func Util_AppendHexJavaScriptRepresentation(sb *strings.Builder, c uint16) {
	var buf [6]byte
	util_putHexUnit(buf[:], c)
	sb.Write(buf[:])
}

// Util_AppendHexJavaScriptRepresentationToAppendable returns a JavaScript
// representation of the character in a hex escaped format.
//
// codePoint is the code point to append; out is the buffer to which the hex
// representation should be appended.
//
// Java: Util.appendHexJavaScriptRepresentation(int, Appendable)
func Util_AppendHexJavaScriptRepresentationToAppendable(codePoint int, out io.StringWriter) {
	var buf [12]byte
	if jutil.CharacterIsSupplementaryCodePoint(codePoint) {
		// Handle supplementary Unicode values which are not representable in
		// JavaScript. We deal with these by escaping them as two 4B
		// sequences so that they will round-trip properly when sent from Java
		// to JavaScript and back.
		hi, lo := utf16.EncodeRune(rune(codePoint))
		util_putHexUnit(buf[0:6], uint16(hi))
		util_putHexUnit(buf[6:12], uint16(lo))
		appendBytes(out, buf[:12])
		return
	}
	util_putHexUnit(buf[0:6], uint16(uint32(codePoint)))
	appendBytes(out, buf[:6])
}
