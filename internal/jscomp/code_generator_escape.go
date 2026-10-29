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
	"unicode/utf16"
	"unicode/utf8"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

const codeGenerator_LT_ESCAPED = "\\x3c"
const codeGenerator_GT_ESCAPED = "\\x3e"
const codeGenerator_HEX_CHARS = "0123456789abcdef"

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

var codeGenerator_templateSafeASCII = func() [256]bool {
	var t [256]bool
	t['\b'] = true
	t['\f'] = true
	t['\n'] = true
	t['\r'] = true
	t['\t'] = true
	for c := 0x20; c <= 0x7e; c++ {
		switch byte(c) {
		case '=', '&', '>', '<':
		default:
			t[c] = true
		}
	}
	return t
}()

// escBuf buffers escaped output on the stack before flushing chunks to
// either a *CodeConsumer or a *strings.Builder, keeping zero pointers in the
// struct so Go's escape analysis never spills it to the heap.
type escBuf struct {
	buf     [128]byte
	n       int
	isASCII bool
}

func (b *escBuf) flush(cc *CodeConsumer, sb *strings.Builder) {
	if b.n == 0 {
		return
	}
	if cc != nil {
		cc.appendContinuationChunk(b.buf[:b.n], b.isASCII)
	} else {
		sb.Write(b.buf[:b.n])
	}
	b.n = 0
	b.isASCII = true
}

func (b *escBuf) writeByte(cc *CodeConsumer, sb *strings.Builder, ch byte) {
	if b.n == len(b.buf) {
		b.flush(cc, sb)
	}
	b.buf[b.n] = ch
	b.n++
}

func (b *escBuf) writeASCII(cc *CodeConsumer, sb *strings.Builder, s string) {
	for len(s) > 0 {
		avail := len(b.buf) - b.n
		if avail == 0 {
			b.flush(cc, sb)
			avail = len(b.buf)
		}
		if len(s) <= avail {
			copy(b.buf[b.n:], s)
			b.n += len(s)
			return
		}
		copy(b.buf[b.n:], s[:avail])
		b.n += avail
		s = s[avail:]
		b.flush(cc, sb)
	}
}

func (b *escBuf) writeEscapeParam(cc *CodeConsumer, sb *strings.Builder, s string) {
	if cc != nil || jutil.IsASCII(s) {
		b.writeASCII(cc, sb, s)
		return
	}
	b.flush(cc, sb)
	sb.WriteString(s)
}

func (b *escBuf) writeHexJS(cc *CodeConsumer, sb *strings.Builder, c uint16) {
	if b.n+6 > len(b.buf) {
		b.flush(cc, sb)
	}
	dst := b.buf[b.n : b.n+6]
	dst[0] = '\\'
	dst[1] = 'u'
	dst[2] = codeGenerator_HEX_CHARS[(c>>12)&0xf]
	dst[3] = codeGenerator_HEX_CHARS[(c>>8)&0xf]
	dst[4] = codeGenerator_HEX_CHARS[(c>>4)&0xf]
	dst[5] = codeGenerator_HEX_CHARS[c&0xf]
	b.n += 6
}

func (b *escBuf) writeChar(cc *CodeConsumer, sb *strings.Builder, c uint16) {
	if c < utf8.RuneSelf {
		b.writeByte(cc, sb, byte(c))
		return
	}
	if b.n+3 > len(b.buf) {
		b.flush(cc, sb)
	}
	if c < 0x800 {
		b.buf[b.n] = byte(0xC0 | (c >> 6))
		b.buf[b.n+1] = byte(0x80 | (c & 0x3F))
		b.n += 2
		b.isASCII = false
		return
	}
	// Encodes both valid BMP runes (0x800..0xD7FF, 0xE000..0xFFFF) and WTF-8
	// lone surrogates (0xD800..0xDFFF, only reachable when sb != nil).
	b.buf[b.n] = byte(0xE0 | (c >> 12))
	b.buf[b.n+1] = byte(0x80 | ((c >> 6) & 0x3F))
	b.buf[b.n+2] = byte(0x80 | (c & 0x3F))
	b.n += 3
	b.isASCII = false
}

// decodeWTF8ToBuffer decodes WTF-8 string s into UTF-16 code units in dst
// (caller ensures cap(dst) >= len(s), which always suffices since UTF-16 code
// unit count <= WTF-8 byte length).
func codeGenerator_decodeWTF8ToBuffer(s string, dst []uint16) []uint16 {
	out := dst[:0]
	prefix := jutil.ASCIIPrefixLen(s)
	for i := 0; i < prefix; i++ {
		out = append(out, uint16(s[i]))
	}
	for i := prefix; i < len(s); {
		c := s[i]
		if c < 0x80 {
			out = append(out, uint16(c))
			i++
			continue
		}
		if c == 0xED && i+3 <= len(s) && s[i+1] >= 0xA0 && s[i+1] <= 0xBF && s[i+2]&0xC0 == 0x80 {
			u := uint16(c&0x0F)<<12 | uint16(s[i+1]&0x3F)<<6 | uint16(s[i+2]&0x3F)
			out = append(out, u)
			i += 3
			continue
		}
		r, sz := utf8.DecodeRuneInString(s[i:])
		i += sz
		if r >= 0x10000 {
			hi, lo := utf16.EncodeRune(r)
			out = append(out, uint16(hi), uint16(lo))
		} else {
			out = append(out, uint16(r))
		}
	}
	return out
}

// decodeWTF8Units decodes s into UTF-16 code units using stackBuf when s fits
// on the stack, or g.scratchUnits otherwise.
func (g *CodeGenerator) decodeWTF8Units(s string, stackBuf []uint16) []uint16 {
	if cap(stackBuf) >= len(s) {
		return codeGenerator_decodeWTF8ToBuffer(s, stackBuf)
	}
	if cap(g.scratchUnits) < len(s) {
		g.scratchUnits = make([]uint16, 0, max(len(s), 256))
	}
	units := codeGenerator_decodeWTF8ToBuffer(s, g.scratchUnits)
	g.scratchUnits = units
	return units
}

// codeGenerator_foldChar folds a UTF-16 code unit identically to Java's
// String.regionMatches(ignoreCase=true) (Character.toLowerCase(Character.toUpperCase(c))).
func codeGenerator_foldChar(c uint16) uint16 {
	if c < 0x80 {
		if c >= 'A' && c <= 'Z' {
			return c | 0x20
		}
		return c
	}
	return jutil.CharacterToLowerCase(jutil.CharacterToUpperCase(c))
}

// codeGenerator_matchesLtHazardASCII checks whether pure-ASCII s[afterLt:]
// starts with "/script" (ignoreCase), "/style" (ignoreCase), or "!--".
func codeGenerator_matchesLtHazardASCII(s string, afterLt int) bool {
	rem := len(s) - afterLt
	if rem < 3 {
		return false
	}
	if s[afterLt] == '!' && s[afterLt+1] == '-' && s[afterLt+2] == '-' {
		return true
	}
	if rem >= 6 && s[afterLt] == '/' && (s[afterLt+1]|0x20) == 's' {
		if rem >= 7 &&
			(s[afterLt+2]|0x20) == 'c' &&
			(s[afterLt+3]|0x20) == 'r' &&
			(s[afterLt+4]|0x20) == 'i' &&
			(s[afterLt+5]|0x20) == 'p' &&
			(s[afterLt+6]|0x20) == 't' {
			return true
		}
		if (s[afterLt+2]|0x20) == 't' &&
			(s[afterLt+3]|0x20) == 'y' &&
			(s[afterLt+4]|0x20) == 'l' &&
			(s[afterLt+5]|0x20) == 'e' {
			return true
		}
	}
	return false
}

// codeGenerator_matchesLtHazardUnits checks whether UTF-16 units[afterLt:]
// starts with "/script" (ignoreCase), "/style" (ignoreCase), or "!--" with
// full Java Character case-folding parity and zero allocations.
func codeGenerator_matchesLtHazardUnits(s []uint16, afterLt int) bool {
	rem := len(s) - afterLt
	if rem < 3 {
		return false
	}
	if s[afterLt] == '!' && s[afterLt+1] == '-' && s[afterLt+2] == '-' {
		return true
	}
	if rem >= 6 && s[afterLt] == '/' && codeGenerator_foldChar(s[afterLt+1]) == 's' {
		if rem >= 7 &&
			codeGenerator_foldChar(s[afterLt+2]) == 'c' &&
			codeGenerator_foldChar(s[afterLt+3]) == 'r' &&
			codeGenerator_foldChar(s[afterLt+4]) == 'i' &&
			codeGenerator_foldChar(s[afterLt+5]) == 'p' &&
			codeGenerator_foldChar(s[afterLt+6]) == 't' {
			return true
		}
		if codeGenerator_foldChar(s[afterLt+2]) == 't' &&
			codeGenerator_foldChar(s[afterLt+3]) == 'y' &&
			codeGenerator_foldChar(s[afterLt+4]) == 'l' &&
			codeGenerator_foldChar(s[afterLt+5]) == 'e' {
			return true
		}
	}
	return false
}

func (g *CodeGenerator) analyzeJsString(s string) (quote byte, doublequote, singlequote string, plain, isASCII bool) {
	singleq := 0
	doubleq := 0
	plain = true
	isASCII = true
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !codeGenerator_plainASCII[c] {
			plain = false
			if c >= utf8.RuneSelf {
				isASCII = false
			} else if c == '"' {
				doubleq++
			} else if c == '\'' {
				singleq++
			}
		}
	}
	if plain {
		if g.preferSingleQuotes {
			return '\'', "", "", true, true
		}
		return '"', "", "", true, true
	}
	useSingle := singleq < doubleq
	if g.preferSingleQuotes {
		useSingle = singleq <= doubleq
	}
	if useSingle {
		return '\'', "\"", "\\'", false, isASCII
	}
	return '"', "\\\"", "'", false, isASCII
}

func (g *CodeGenerator) escapeInto(
	b *escBuf,
	cc *CodeConsumer,
	sb *strings.Builder,
	s string,
	isASCII bool,
	doublequoteEscape string,
	singlequoteEscape string,
	backtickEscape string,
	backslashEscape string,
	dollarEscape string,
	isRegexp bool,
) {
	if isASCII {
		g.strEscapeASCII(b, cc, sb, s, doublequoteEscape, singlequoteEscape, backtickEscape, backslashEscape, dollarEscape, isRegexp)
		return
	}
	var stackUnits [128]uint16
	units := g.decodeWTF8Units(s, stackUnits[:0])
	g.strEscapeUnits(b, cc, sb, units, doublequoteEscape, singlequoteEscape, backtickEscape, backslashEscape, dollarEscape, isRegexp)
}

// AddJsString outputs a JS string, using the optimal (single/double) quote
// character, streaming directly into the CodeConsumer without intermediate
// string allocations.
//
// Java: CodeGenerator.addJsString(Node)
func (g *CodeGenerator) AddJsString(n rhino.Node) {
	s := n.GetString()
	quote, doublequote, singlequote, plain, isASCII := g.analyzeJsString(s)
	if plain {
		g.cc.addPlainQuotedASCII(quote, s)
		return
	}
	cc := &g.cc
	cc.beginToken(quote)
	b := escBuf{isASCII: true}
	b.writeByte(cc, nil, quote)
	g.escapeInto(&b, cc, nil, s, isASCII, doublequote, singlequote, "`", "\\\\", "$", false)
	b.writeByte(cc, nil, quote)
	b.flush(cc, nil)
}

// addRegexp streams /pattern/flags directly into g.cc without temporary string
// or []uint16 allocations.
func (g *CodeGenerator) addRegexp(pattern string, flags string) {
	if flags != "" && (!jutil.IsASCII(flags) || strings.IndexByte(flags, '\n') >= 0) {
		g.Add(g.RegexpEscape(pattern) + flags)
		return
	}
	cc := &g.cc
	cc.beginToken('/')
	b := escBuf{isASCII: true}
	b.writeByte(cc, nil, '/')
	g.escapeInto(&b, cc, nil, pattern, jutil.IsASCII(pattern), "\"", "'", "`", "\\", "$", true)
	b.writeByte(cc, nil, '/')
	if flags != "" {
		b.writeASCII(cc, nil, flags)
	}
	b.flush(cc, nil)
}

// addNonLatinIdentifier streams an escaped non-Latin identifier directly into
// g.cc without allocating an intermediate string or []uint16.
func (g *CodeGenerator) addNonLatinIdentifier(s string) {
	var stackUnits [128]uint16
	units := g.decodeWTF8Units(s, stackUnits[:0])
	firstByte := byte('\\')
	if units[0] > 0x1F && units[0] < 0x7F {
		firstByte = byte(units[0])
	}
	cc := &g.cc
	cc.beginToken(firstByte)
	b := escBuf{isASCII: true}
	for _, c := range units {
		if c > 0x1F && c < 0x7F {
			b.writeByte(cc, nil, byte(c))
		} else {
			b.writeHexJS(cc, nil, c)
		}
	}
	b.flush(cc, nil)
}

// Java: CodeGenerator.jsString(String)
func (g *CodeGenerator) JsString(s string) string {
	quote, doublequote, singlequote, plain, isASCII := g.analyzeJsString(s)
	if plain {
		var sb strings.Builder
		sb.Grow(len(s) + 2)
		sb.WriteByte(quote)
		sb.WriteString(s)
		sb.WriteByte(quote)
		return sb.String()
	}
	var sb strings.Builder
	sb.Grow(len(s) + 4)
	b := escBuf{isASCII: true}
	b.writeByte(nil, &sb, quote)
	g.escapeInto(&b, nil, &sb, s, isASCII, doublequote, singlequote, "`", "\\\\", "$", false)
	b.writeByte(nil, &sb, quote)
	b.flush(nil, &sb)
	return sb.String()
}

// RegexpEscape escapes regular expression.
//
// Java: CodeGenerator.regexpEscape(String)
func (g *CodeGenerator) RegexpEscape(s string) string {
	var sb strings.Builder
	sb.Grow(len(s) + 2)
	b := escBuf{isASCII: true}
	b.writeByte(nil, &sb, '/')
	g.escapeInto(&b, nil, &sb, s, jutil.IsASCII(s), "\"", "'", "`", "\\", "$", true)
	b.writeByte(nil, &sb, '/')
	b.flush(nil, &sb)
	return sb.String()
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
	isRegexp bool,
) string {
	plain := true
	isASCII := true
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !codeGenerator_plainASCII[c] {
			plain = false
			if c >= utf8.RuneSelf {
				isASCII = false
				break
			}
		}
	}
	if plain {
		return s
	}
	var sb strings.Builder
	sb.Grow(len(s) + 2)
	b := escBuf{isASCII: true}
	g.escapeInto(&b, nil, &sb, s, isASCII, doublequoteEscape, singlequoteEscape, backtickEscape, backslashEscape, dollarEscape, isRegexp)
	b.flush(nil, &sb)
	return sb.String()
}

func (g *CodeGenerator) writeCommonEscape(
	b *escBuf,
	cc *CodeConsumer,
	sb *strings.Builder,
	c byte,
	doublequoteEscape string,
	singlequoteEscape string,
	backtickEscape string,
	backslashEscape string,
	dollarEscape string,
	isRegexp bool,
) bool {
	switch c {
	case 0:
		b.writeASCII(cc, sb, "\\x00")
	case '\u000B':
		if !isRegexp {
			b.writeASCII(cc, sb, "\\v")
		} else {
			b.writeASCII(cc, sb, "\\x0B")
		}
	case '\b':
		b.writeASCII(cc, sb, "\\b")
	case '\f':
		b.writeASCII(cc, sb, "\\f")
	case '\n':
		b.writeASCII(cc, sb, "\\n")
	case '\r':
		b.writeASCII(cc, sb, "\\r")
	case '\t':
		b.writeASCII(cc, sb, "\\t")
	case '\\':
		b.writeEscapeParam(cc, sb, backslashEscape)
	case '"':
		b.writeEscapeParam(cc, sb, doublequoteEscape)
	case '\'':
		b.writeEscapeParam(cc, sb, singlequoteEscape)
	case '$':
		b.writeEscapeParam(cc, sb, dollarEscape)
	case '`':
		b.writeEscapeParam(cc, sb, backtickEscape)
	case '=':
		if g.trustedStrings || isRegexp {
			b.writeByte(cc, sb, '=')
		} else {
			b.writeASCII(cc, sb, "\\x3d")
		}
	case '&':
		if g.trustedStrings || isRegexp {
			b.writeByte(cc, sb, '&')
		} else {
			b.writeASCII(cc, sb, "\\x26")
		}
	default:
		return false
	}
	return true
}

func (g *CodeGenerator) strEscapeASCII(
	b *escBuf,
	cc *CodeConsumer,
	sb *strings.Builder,
	s string,
	doublequoteEscape string,
	singlequoteEscape string,
	backtickEscape string,
	backslashEscape string,
	dollarEscape string,
	isRegexp bool,
) {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if codeGenerator_plainASCII[c] {
			b.writeByte(cc, sb, c)
			continue
		}
		if g.writeCommonEscape(b, cc, sb, c, doublequoteEscape, singlequoteEscape, backtickEscape, backslashEscape, dollarEscape, isRegexp) {
			continue
		}
		switch c {
		case '>':
			if !g.trustedStrings && !isRegexp {
				b.writeASCII(cc, sb, codeGenerator_GT_ESCAPED)
			} else if i >= 2 && ((s[i-1] == '-' && s[i-2] == '-') || (s[i-1] == ']' && s[i-2] == ']')) {
				b.writeASCII(cc, sb, codeGenerator_GT_ESCAPED)
			} else {
				b.writeByte(cc, sb, '>')
			}
		case '<':
			if !g.trustedStrings && !isRegexp {
				b.writeASCII(cc, sb, codeGenerator_LT_ESCAPED)
			} else if codeGenerator_matchesLtHazardASCII(s, i+1) {
				b.writeASCII(cc, sb, codeGenerator_LT_ESCAPED)
			} else {
				b.writeByte(cc, sb, '<')
			}
		default:
			if (g.outputCharsetEncoder != nil && g.outputCharsetEncoder.CanEncode(uint16(c))) ||
				(c > 0x1f && c < 0x7f) {
				b.writeByte(cc, sb, c)
			} else {
				b.writeHexJS(cc, sb, uint16(c))
			}
		}
	}
}

func (g *CodeGenerator) strEscapeUnits(
	b *escBuf,
	cc *CodeConsumer,
	sb *strings.Builder,
	units []uint16,
	doublequoteEscape string,
	singlequoteEscape string,
	backtickEscape string,
	backslashEscape string,
	dollarEscape string,
	isRegexp bool,
) {
	escapeLineSep := isRegexp ||
		!g.outputFeatureSet.ContainsFeature(parser.FeatureSet_Feature_UNESCAPED_UNICODE_LINE_OR_PARAGRAPH_SEP)

	for i := 0; i < len(units); i++ {
		c := units[i]
		if c < 0x80 {
			bc := byte(c)
			if codeGenerator_plainASCII[bc] {
				b.writeByte(cc, sb, bc)
				continue
			}
			if g.writeCommonEscape(b, cc, sb, bc, doublequoteEscape, singlequoteEscape, backtickEscape, backslashEscape, dollarEscape, isRegexp) {
				continue
			}
		}
		switch c {
		case '>':
			if !g.trustedStrings && !isRegexp {
				b.writeASCII(cc, sb, codeGenerator_GT_ESCAPED)
			} else if i >= 2 && ((units[i-1] == '-' && units[i-2] == '-') || (units[i-1] == ']' && units[i-2] == ']')) {
				b.writeASCII(cc, sb, codeGenerator_GT_ESCAPED)
			} else {
				b.writeByte(cc, sb, '>')
			}
		case '<':
			if !g.trustedStrings && !isRegexp {
				b.writeASCII(cc, sb, codeGenerator_LT_ESCAPED)
			} else if codeGenerator_matchesLtHazardUnits(units, i+1) {
				b.writeASCII(cc, sb, codeGenerator_LT_ESCAPED)
			} else {
				b.writeByte(cc, sb, '<')
			}
		default:
			if escapeLineSep {
				// In 2019 these characters (line and paragraph separators)
				// are valid in strings but not regular expressions.
				// https://github.com/tc39/proposal-json-superset
				if c == '\u2028' {
					b.writeASCII(cc, sb, "\\u2028")
					break
				}
				if c == '\u2029' {
					b.writeASCII(cc, sb, "\\u2029")
					break
				}
			}

			if (g.outputCharsetEncoder != nil && g.outputCharsetEncoder.CanEncode(c)) ||
				(c > 0x1f && c < 0x7f) {
				b.writeChar(cc, sb, c)
			} else {
				b.writeHexJS(cc, sb, c)
			}
		}
	}
}

// EscapeUnrecognizedCharacters is a helper to escape the characters that
// might be misinterpreted.
//
// Java: CodeGenerator.escapeUnrecognizedCharacters(String)
func (g *CodeGenerator) EscapeUnrecognizedCharacters(s string) string {
	unchanged := true
	for i := 0; i < len(s); i++ {
		c := s[i]
		if codeGenerator_templateSafeASCII[c] {
			continue
		}
		if c < utf8.RuneSelf && g.trustedStrings {
			switch c {
			case '=', '&':
				continue
			case '>':
				if !(i >= 2 && ((s[i-1] == '-' && s[i-2] == '-') || (s[i-1] == ']' && s[i-2] == ']'))) {
					continue
				}
			case '<':
				if !codeGenerator_matchesLtHazardASCII(s, i+1) {
					continue
				}
			}
		}
		unchanged = false
		break
	}
	if unchanged {
		return s
	}

	var stackUnits [128]uint16
	units := g.decodeWTF8Units(s, stackUnits[:0])
	var sb strings.Builder
	sb.Grow(len(s))
	b := escBuf{isASCII: true}
	for i := 0; i < len(units); i++ {
		c := units[i]
		switch c {
		case '\b', '\f', '\n', '\r', '\t', '\\', '"', '\'', '$', '`', '\u2028', '\u2029':
			b.writeChar(nil, &sb, c)
		case '=':
			if g.trustedStrings {
				b.writeByte(nil, &sb, '=')
			} else {
				b.writeASCII(nil, &sb, "\\x3d")
			}
		case '&':
			if g.trustedStrings {
				b.writeByte(nil, &sb, '&')
			} else {
				b.writeASCII(nil, &sb, "\\x26")
			}
		case '>':
			if !g.trustedStrings ||
				(i >= 2 && ((units[i-1] == '-' && units[i-2] == '-') || (units[i-1] == ']' && units[i-2] == ']'))) {
				b.writeASCII(nil, &sb, codeGenerator_GT_ESCAPED)
			} else {
				b.writeByte(nil, &sb, '>')
			}
		case '<':
			if !g.trustedStrings || codeGenerator_matchesLtHazardUnits(units, i+1) {
				b.writeASCII(nil, &sb, codeGenerator_LT_ESCAPED)
			} else {
				b.writeByte(nil, &sb, '<')
			}
		default:
			if (g.outputCharsetEncoder != nil && g.outputCharsetEncoder.CanEncode(c)) ||
				(c > 0x1f && c < 0x7f) {
				b.writeChar(nil, &sb, c)
			} else {
				b.writeHexJS(nil, &sb, c)
			}
		}
	}
	b.flush(nil, &sb)
	return sb.String()
}

// Java: CodeGenerator.identifierEscape(String)
func CodeGenerator_IdentifierEscape(s string) string {
	// First check if escaping is needed at all -- in most cases it isn't.
	if NodeUtil_IsLatin(s) {
		return s
	}

	var stackUnits [128]uint16
	var units []uint16
	if len(s) <= len(stackUnits) {
		units = codeGenerator_decodeWTF8ToBuffer(s, stackUnits[:0])
	} else {
		units = jutil.StringToUTF16(s)
	}
	var sb strings.Builder
	sb.Grow(len(s) + 6)
	b := escBuf{isASCII: true}
	for _, c := range units {
		if c > 0x1F && c < 0x7F {
			b.writeByte(nil, &sb, byte(c))
		} else {
			b.writeHexJS(nil, &sb, c)
		}
	}
	b.flush(nil, &sb)
	return sb.String()
}
