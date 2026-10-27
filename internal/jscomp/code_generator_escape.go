// Port of com/google/javascript/jscomp/CodeGenerator.java @ closure-compiler v20261021

/*
 * Copyright 2004 The Closure Compiler Authors.
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

package jscomp

import (
	"strings"
	"unicode/utf8"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

const codeGenerator_LT_ESCAPED = "\\x3c"
const codeGenerator_GT_ESCAPED = "\\x3e"

// utf16Units is the UTF-16 view of s, in the reused scratch buffer when s
// is ASCII.
//
// Go-only: language-adapter — see scratchUnits.
func (g *CodeGenerator) utf16Units(s string) []uint16 {
	units := g.scratchUnits[:0]
	for i := 0; i < len(s); i++ {
		if s[i] >= utf8.RuneSelf {
			return jutil.StringToUTF16(s)
		}
		units = append(units, uint16(s[i]))
	}
	g.scratchUnits = units
	return units
}

// AddJsString outputs a JS string, using the optimal (single/double) quote
// character.
//
// Java: CodeGenerator.addJsString(Node)
func (g *CodeGenerator) AddJsString(n rhino.Node) {
	g.Add(g.JsString(n.GetString()))
}

var codeGenerator_plainASCII = func() [256]bool {
	var t [256]bool
	for c := 0x20; c <= 0x7e; c++ {
		switch byte(c) {
		case '\\', '"', '\'', '$', '`', '=', '&', '>', '<':
		default:
			t[c] = true
		}
	}
	return t
}()

const codeGenerator_HEX_CHARS = "0123456789abcdef"

func codeGenerator_appendHexJS(sb *strings.Builder, c uint16) {
	var buf [6]byte
	buf[0] = '\\'
	buf[1] = 'u'
	buf[2] = codeGenerator_HEX_CHARS[(c>>12)&0xf]
	buf[3] = codeGenerator_HEX_CHARS[(c>>8)&0xf]
	buf[4] = codeGenerator_HEX_CHARS[(c>>4)&0xf]
	buf[5] = codeGenerator_HEX_CHARS[c&0xf]
	sb.Write(buf[:])
}

// Java: CodeGenerator.jsString(String)
func (g *CodeGenerator) JsString(s string) string {
	singleq := 0
	doubleq := 0
	plain := true

	// could count the quotes and pick the optimal quote character
	//
	// Go-only: language-adapter — both quotes are ASCII, and no byte of a
	// multi-byte WTF-8 sequence is ASCII, so counting bytes counts the same
	// chars as Java's UTF-16 loop.
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !codeGenerator_plainASCII[c] {
			plain = false
			switch c {
			case '"':
				doubleq++
			case '\'':
				singleq++
			}
		}
	}

	if plain {
		quote := byte('"')
		if g.preferSingleQuotes {
			quote = '\''
		}
		var sb strings.Builder
		sb.Grow(len(s) + 2)
		sb.WriteByte(quote)
		sb.WriteString(s)
		sb.WriteByte(quote)
		return sb.String()
	}

	var doublequote string
	var singlequote string
	var quote byte
	useSingle := singleq < doubleq
	if g.preferSingleQuotes {
		useSingle = singleq <= doubleq
	}
	if useSingle {
		// more double quotes so enclose in single quotes.
		quote = '\''
		doublequote = "\""
		singlequote = "\\'"
	} else {
		// more single quotes so escape the doubles
		quote = '"'
		doublequote = "\\\""
		singlequote = "'"
	}

	var sb strings.Builder
	sb.Grow(len(s) + 4)
	sb.WriteByte(quote)
	g.strEscapeInto(&sb, s, doublequote, singlequote, "`", "\\\\", "$", false)
	sb.WriteByte(quote)
	return sb.String()
}

// RegexpEscape escapes regular expression.
//
// Java: CodeGenerator.regexpEscape(String)
func (g *CodeGenerator) RegexpEscape(s string) string {
	var sb strings.Builder
	sb.Grow(len(s) + 2)
	sb.WriteByte('/')
	g.strEscapeInto(&sb, s, "\"", "'", "`", "\\", "$", true)
	sb.WriteByte('/')
	return sb.String()
}

// codeGenerator_appendChar is sb.append(char) for a char that is not
// escaped.
//
// Go-only: language-adapter — appends one UTF-16 char to WTF-8 output.
func codeGenerator_appendChar(sb *strings.Builder, c uint16) {
	if c < utf8.RuneSelf {
		sb.WriteByte(byte(c))
	} else if c >= 0xD800 && c <= 0xDFFF {
		// A surrogate half on its own: kept as WTF-8.
		sb.WriteString(jutil.StringFromChar(c))
	} else {
		sb.WriteRune(rune(c))
	}
}

// codeGenerator_regionMatches is String.regionMatches(boolean, int, String,
// int, int) on the UTF-16 units of the string being escaped; other is
// ASCII.
//
// Go-only: language-adapter — Java's regionMatches on the UTF-16 view.
func codeGenerator_regionMatches(units []uint16, ignoreCase bool, toffset int, other string) bool {
	return jutil.UTF16RegionMatches(ignoreCase, units, toffset, jutil.StringToUTF16(other), 0, len(other))
}

// StrEscape is a helper to escape JavaScript string as well as regular
// expression.
//
// Java: CodeGenerator.strEscape(String, String, String, String, String, String, boolean)
func (g *CodeGenerator) StrEscape(
	s string,
	doublequoteEscape string,
	singlequoteEscape string,
	backtickEscape string,
	backslashEscape string,
	dollarEscape string,
	isRegexp bool) string {
	plain := true
	for i := 0; i < len(s); i++ {
		if !codeGenerator_plainASCII[s[i]] {
			plain = false
			break
		}
	}
	if plain {
		return s
	}
	var sb strings.Builder
	sb.Grow(len(s) + 2)
	g.strEscapeInto(&sb, s, doublequoteEscape, singlequoteEscape, backtickEscape, backslashEscape, dollarEscape, isRegexp)
	return sb.String()
}

func (g *CodeGenerator) strEscapeInto(
	sb *strings.Builder,
	s string,
	doublequoteEscape string,
	singlequoteEscape string,
	backtickEscape string,
	backslashEscape string,
	dollarEscape string,
	isRegexp bool) {
	// Go-only: language-adapter — Java indexes the UTF-16 chars of s.
	units := g.utf16Units(s)
	for i := 0; i < len(units); i++ {
		c := units[i]
		switch c {
		case 0:
			sb.WriteString("\\x00")
		case '\u000B':
			if !isRegexp {
				sb.WriteString("\\v")
			} else {
				sb.WriteString("\\x0B")
			}
			// From the SingleEscapeCharacter grammar production.
		case '\b':
			sb.WriteString("\\b")
		case '\f':
			sb.WriteString("\\f")
		case '\n':
			sb.WriteString("\\n")
		case '\r':
			sb.WriteString("\\r")
		case '\t':
			sb.WriteString("\\t")
		case '\\':
			sb.WriteString(backslashEscape)
		case '"':
			sb.WriteString(doublequoteEscape)
		case '\'':
			sb.WriteString(singlequoteEscape)
		case '$':
			sb.WriteString(dollarEscape)
		case '`':
			sb.WriteString(backtickEscape)
		case '=':
			g.AppendEscapedEq(sb, isRegexp)
		case '&':
			g.AppendEscapedAmp(sb, isRegexp)
		case '>':
			g.AppendEscapedGt(sb, units, i, isRegexp)
		case '<':
			g.AppendEscapedLt(sb, units, i, isRegexp)
		default:
			if isRegexp ||
				!g.outputFeatureSet.ContainsFeature(parser.FeatureSet_Feature_UNESCAPED_UNICODE_LINE_OR_PARAGRAPH_SEP) {
				// In 2019 these characters (line and paragraph separators)
				// are valid in strings but not regular expressions.
				// https://github.com/tc39/proposal-json-superset
				if c == ' ' {
					sb.WriteString("\\u2028")
					break
				}
				if c == ' ' {
					sb.WriteString("\\u2029")
					break
				}
			}

			if (g.outputCharsetEncoder != nil && g.outputCharsetEncoder.CanEncode(c)) ||
				(c > 0x1f && c < 0x7f) {
				// If we're given an outputCharsetEncoder, then check if the
				// character can be represented in this character set. If no
				// charsetEncoder provided - pass straight Latin characters
				// through, and escape the rest. Doing the explicit character
				// check is measurably faster than using the CharsetEncoder.
				codeGenerator_appendChar(sb, c)
			} else {
				// Other characters can be misinterpreted by some JS parsers,
				// or perhaps mangled by proxies along the way,
				// so we play it safe and Unicode escape them.
				codeGenerator_appendHexJS(sb, c)
			}
		}
	}
}

// Java: CodeGenerator.appendEscapedEq(StringBuilder, boolean)
func (g *CodeGenerator) AppendEscapedEq(sb *strings.Builder, isRegexp bool) {
	// '=' is a syntactically significant regexp character.
	if g.trustedStrings || isRegexp {
		sb.WriteByte('=')
	} else {
		sb.WriteString("\\x3d")
	}
}

// Java: CodeGenerator.appendEscapedAmp(StringBuilder, boolean)
func (g *CodeGenerator) AppendEscapedAmp(sb *strings.Builder, isRegexp bool) {
	if g.trustedStrings || isRegexp {
		sb.WriteByte('&')
	} else {
		sb.WriteString("\\x26")
	}
}

// AppendEscapedGt: s is the string being escaped, as its UTF-16 units.
//
// Java: CodeGenerator.appendEscapedGt(StringBuilder, String, int, boolean)
func (g *CodeGenerator) AppendEscapedGt(sb *strings.Builder, s []uint16, index int, isRegexp bool) {
	if !g.trustedStrings && !isRegexp {
		sb.WriteString(codeGenerator_GT_ESCAPED)
		return
	}

	// Break --> into --\> or ]]> into ]]\>
	//
	// This is just to prevent developers from shooting themselves in the
	// foot, and does not provide the level of security that you get with
	// trustedString == false.
	if index >= 2 &&
		((s[index-1] == '-' && s[index-2] == '-') ||
			(s[index-1] == ']' && s[index-2] == ']')) {
		sb.WriteString(codeGenerator_GT_ESCAPED)
	} else {
		sb.WriteByte('>')
	}
}

// AppendEscapedLt: s is the string being escaped, as its UTF-16 units.
//
// Java: CodeGenerator.appendEscapedLt(StringBuilder, String, int, boolean)
func (g *CodeGenerator) AppendEscapedLt(sb *strings.Builder, s []uint16, index int, isRegexp bool) {
	if !g.trustedStrings && !isRegexp {
		sb.WriteString(codeGenerator_LT_ESCAPED)
		return
	}

	// Break </script into <\/script and </style into <\/style
	// As above, this is just to prevent developers from doing this accidentally.
	const endScript = "/script"
	const endStyle = "/style"

	// Break <!-- into <\!--
	const startComment = "!--"

	if codeGenerator_regionMatches(s, true, index+1, endScript) ||
		codeGenerator_regionMatches(s, true, index+1, endStyle) ||
		codeGenerator_regionMatches(s, false, index+1, startComment) {
		sb.WriteString(codeGenerator_LT_ESCAPED)
	} else {
		sb.WriteByte('<')
	}
}

// EscapeUnrecognizedCharacters is a helper to escape the characters that
// might be misinterpreted.
//
// s is the string to modify; returns the string with unrecognizable
// characters escaped.
//
// Java: CodeGenerator.escapeUnrecognizedCharacters(String)
func (g *CodeGenerator) EscapeUnrecognizedCharacters(s string) string {
	// TODO(yitingwang) Move this method to a suitable place
	plain := true
	for i := 0; i < len(s); i++ {
		if !codeGenerator_plainASCII[s[i]] {
			plain = false
			break
		}
	}
	if plain {
		return s
	}
	// Go-only: language-adapter — Java indexes the UTF-16 chars of s.
	units := g.utf16Units(s)
	var sb strings.Builder
	sb.Grow(len(s))
	for i := 0; i < len(units); i++ {
		c := units[i]
		switch c {
		case '\b', '\f', '\n', '\r', '\t', '\\', '"', '\'', '$', '`', ' ', ' ':
			// From the SingleEscapeCharacter grammar production.
			codeGenerator_appendChar(&sb, c)
		case '=':
			g.AppendEscapedEq(&sb /* isRegexp= */, false)
		case '&':
			g.AppendEscapedAmp(&sb /* isRegexp= */, false)
		case '>':
			g.AppendEscapedGt(&sb, units, i /* isRegexp= */, false)
		case '<':
			g.AppendEscapedLt(&sb, units, i /* isRegexp= */, false)
		default:
			if (g.outputCharsetEncoder != nil && g.outputCharsetEncoder.CanEncode(c)) ||
				(c > 0x1f && c < 0x7f) {
				// If we're given an outputCharsetEncoder, then check if the
				// character can be represented in this character set. If no
				// charsetEncoder provided - pass straight Latin characters
				// through, and escape the rest. Doing the explicit character
				// check is measurably faster than using the CharsetEncoder.
				codeGenerator_appendChar(&sb, c)
			} else {
				// Other characters can be misinterpreted by some JS parsers,
				// or perhaps mangled by proxies along the way,
				// so we play it safe and Unicode escape them.
				codeGenerator_appendHexJS(&sb, c)
			}
		}
	}
	return sb.String()
}

// Java: CodeGenerator.identifierEscape(String)
func CodeGenerator_IdentifierEscape(s string) string {
	// First check if escaping is needed at all -- in most cases it isn't.
	if NodeUtil_IsLatin(s) {
		return s
	}

	// Now going through the string to escape non-Latin characters if needed.
	//
	// Go-only: language-adapter — Java indexes the UTF-16 chars of s.
	units := jutil.StringToUTF16(s)
	var sb strings.Builder
	sb.Grow(len(s) + 6)
	for i := 0; i < len(units); i++ {
		c := units[i]
		// Identifiers should always go to Latin1/ ASCII characters because
		// different browser's rules for valid identifier characters are
		// crazy.
		if c > 0x1F && c < 0x7F {
			sb.WriteByte(byte(c))
		} else {
			codeGenerator_appendHexJS(&sb, c)
		}
	}
	return sb.String()
}
