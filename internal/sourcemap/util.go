// Port of com/google/debugging/sourcemap/Util.java @ closure-compiler v20261015

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

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Util holds the source map string escaper.
type Util struct{}

// Util_HEX_CHARS is the hex digit table.
var Util_HEX_CHARS = [...]uint16{'0', '1', '2', '3', '4', '5', '6', '7',
	'8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}

func util_needsEscape(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c <= 0x1f || c > 0x7f || c == '\\' || c == '"' || c == '>' || c == '<' {
			return true
		}
	}
	return false
}

// Util_EscapeString escapes the given string to a double quoted (")
// JavaScript/JSON string.
//
// Java: Util.escapeString(String)
func Util_EscapeString(s string) string {
	if !util_needsEscape(s) {
		return `"` + s + `"`
	}
	const quote = '"'
	const doublequoteEscape = "\\\""
	const singlequoteEscape = "'"
	const backslashEscape = "\\\\"

	if jutil.IsASCII(s) {
		var sb strings.Builder
		sb.Grow(len(s) + 16)
		sb.WriteByte(quote)
		region := util_UnescapedRegion{sb: &sb, ascii: s}
		length := len(s)
		for i := 0; i < length; i++ {
			c := s[i]
			switch c {
			case '\n':
				region.AppendForEscapedChar("\\n")
			case '\r':
				region.AppendForEscapedChar("\\r")
			case '\t':
				region.AppendForEscapedChar("\\t")
			case '\\':
				region.AppendForEscapedChar(backslashEscape)
			case '"':
				region.AppendForEscapedChar(doublequoteEscape)
			case '\'':
				region.AppendForEscapedChar(singlequoteEscape)
			case '>':
				if i >= 2 &&
					((s[i-1] == '-' && s[i-2] == '-') ||
						(s[i-1] == ']' && s[i-2] == ']')) {
					region.AppendForEscapedChar("\\u003e")
				} else {
					region.IncrementForNormalChar()
				}
			case '<':
				if i+1+len("/script") <= length && strings.EqualFold(s[i+1:i+1+len("/script")], "/script") {
					region.AppendForEscapedChar("\\u003c")
				} else if i+1+len("!--") <= length && s[i+1:i+1+len("!--")] == "!--" {
					region.AppendForEscapedChar("\\u003c")
				} else {
					region.IncrementForNormalChar()
				}
			default:
				if c > 0x1f && c <= 0x7f {
					region.IncrementForNormalChar()
				} else {
					region.AppendUnescaped()
					region.IncrementForEscapedChar()
					Util_AppendHexJavaScriptRepresentation(&sb, uint16(c))
				}
			}
		}
		region.AppendUnescaped()
		sb.WriteByte(quote)
		return sb.String()
	}

	// Go-only: language-adapter — the Java String as UTF-16 code units, so
	// that charAt is O(1) and every index is a Java char index.
	units := jutil.StringToUTF16(s)

	var sb strings.Builder
	sb.Grow(len(units) + 2)
	sb.WriteByte(quote)

	region := util_UnescapedRegion{sb: &sb, units: units}

	length := len(units)
	for i := 0; i < length; i++ {
		c := units[i]
		switch c {
		case '\n':
			region.AppendForEscapedChar("\\n")
		case '\r':
			region.AppendForEscapedChar("\\r")
		case '\t':
			region.AppendForEscapedChar("\\t")
		case '\\':
			region.AppendForEscapedChar(backslashEscape)
		case '"':
			region.AppendForEscapedChar(doublequoteEscape)
		case '\'':
			region.AppendForEscapedChar(singlequoteEscape)
		case '>':
			// Unicode-escape the '>' in '-->' and ']]>'
			if i >= 2 &&
				((units[i-1] == '-' && units[i-2] == '-') ||
					(units[i-1] == ']' && units[i-2] == ']')) {
				region.AppendForEscapedChar("\\u003e")
			} else {
				region.IncrementForNormalChar()
			}
		case '<':
			// Unicode-escape the '<' in '</script' and '<!--'
			if jutil.UTF16RegionMatches(true, units, i+1, util_END_SCRIPT, 0, len(util_END_SCRIPT)) {
				region.AppendForEscapedChar("\\u003c")
			} else if jutil.UTF16RegionMatches(false, units, i+1, util_START_COMMENT, 0, len(util_START_COMMENT)) {
				region.AppendForEscapedChar("\\u003c")
			} else {
				region.IncrementForNormalChar()
			}
		default:
			// No charsetEncoder provided - pass straight Latin characters
			// through, and escape the rest. Doing the explicit character
			// check is measurably faster than using the CharsetEncoder.
			if c > 0x1f && c <= 0x7f {
				region.IncrementForNormalChar()
			} else {
				// Other characters can be misinterpreted by some JS parsers,
				// or perhaps mangled by proxies along the way, so we play it
				// safe and Unicode escape them.
				region.AppendUnescaped()
				region.IncrementForEscapedChar()
				Util_AppendHexJavaScriptRepresentation(&sb, c)
			}
		}
	}
	region.AppendUnescaped()
	sb.WriteByte(quote)
	return sb.String()
}

// The '<' escapes' lookaheads (local constants END_SCRIPT and START_COMMENT
// of escapeString), as UTF-16 code units.
var (
	util_END_SCRIPT    = jutil.StringToUTF16("/script")
	util_START_COMMENT = jutil.StringToUTF16("!--")
)

// util_UnescapedRegion is escapeString's local class UnescapedRegion: the
// run of chars that need no escaping and have not been appended yet.
type util_UnescapedRegion struct {
	// Go-only: language-adapter — the captured locals sb and s.
	sb    *strings.Builder
	ascii string
	units []uint16

	unescapedRegionStart int
	unescapedRegionEnd   int
}

// Java: Util.escapeString.UnescapedRegion.appendUnescaped()
func (r *util_UnescapedRegion) AppendUnescaped() {
	if r.unescapedRegionStart != r.unescapedRegionEnd {
		// Note: we want to use the "append(String)" override for performance
		// reasons and the ErrorProne suggestion to use the CharSequence
		// override is inferior for that.
		if r.units == nil {
			r.sb.WriteString(r.ascii[r.unescapedRegionStart:r.unescapedRegionEnd])
		} else {
			for _, u := range r.units[r.unescapedRegionStart:r.unescapedRegionEnd] {
				r.sb.WriteByte(byte(u))
			}
		}
	}
	r.unescapedRegionStart = r.unescapedRegionEnd
}

// Java: Util.escapeString.UnescapedRegion.incrementForNormalChar()
func (r *util_UnescapedRegion) IncrementForNormalChar() {
	r.unescapedRegionEnd++
}

// Java: Util.escapeString.UnescapedRegion.incrementForEscapedChar()
func (r *util_UnescapedRegion) IncrementForEscapedChar() {
	if r.unescapedRegionStart != r.unescapedRegionEnd {
		panic(jutil.NewIllegalStateException())
	}
	r.unescapedRegionStart++
	r.unescapedRegionEnd++
}

// Java: Util.escapeString.UnescapedRegion.appendForEscapedChar(String)
func (r *util_UnescapedRegion) AppendForEscapedChar(escaped string) {
	r.AppendUnescaped()
	r.IncrementForEscapedChar()
	r.sb.WriteString(escaped)
}

// Util_AppendHexJavaScriptRepresentation appends the hex escape of c.
//
// Java: Util.appendHexJavaScriptRepresentation(StringBuilder, char)
func Util_AppendHexJavaScriptRepresentation(sb *strings.Builder, c uint16) {
	var buf [6]byte
	buf[0] = '\\'
	buf[1] = 'u'
	buf[2] = byte(Util_HEX_CHARS[(c>>12)&0xf])
	buf[3] = byte(Util_HEX_CHARS[(c>>8)&0xf])
	buf[4] = byte(Util_HEX_CHARS[(c>>4)&0xf])
	buf[5] = byte(Util_HEX_CHARS[c&0xf])
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
	if jutil.CharacterIsSupplementaryCodePoint(codePoint) {
		// Handle supplementary Unicode values which are not representable in
		// JavaScript. We deal with these by escaping them as two 4B
		// sequences so that they will round-trip properly when sent from Java
		// to JavaScript and back.
		surrogates := jutil.CharacterToChars(codePoint)
		Util_AppendHexJavaScriptRepresentationToAppendable(int(surrogates[0]), out)
		Util_AppendHexJavaScriptRepresentationToAppendable(int(surrogates[1]), out)
		return
	}
	cp := int32(codePoint)
	appendString(out, "\\u")
	appendChar(out, Util_HEX_CHARS[jutil.Ushr32(cp, 12)&0xf])
	appendChar(out, Util_HEX_CHARS[jutil.Ushr32(cp, 8)&0xf])
	appendChar(out, Util_HEX_CHARS[jutil.Ushr32(cp, 4)&0xf])
	appendChar(out, Util_HEX_CHARS[cp&0xf])
}
