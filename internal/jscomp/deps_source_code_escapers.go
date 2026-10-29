// Port of com/google/javascript/jscomp/deps/SourceCodeEscapers.java @ closure-compiler v20261021

/*
 * Copyright 2014 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// SourceCodeEscapers is a factory for Escaper instances used to escape
// strings for safe use in various common programming languages.
//
// NOTE: This class is cribbed from the Guava libraries SourceCodeEscapers
// which is not part of the current Guava release.
// https://github.com/google/guava/issues/1620
type SourceCodeEscapers struct{}

const (
	// From: http://en.wikipedia.org/wiki/ASCII#ASCII_printable_characters
	sourceCodeEscapers_PRINTABLE_ASCII_MIN uint16 = 0x20 // ' '
	sourceCodeEscapers_PRINTABLE_ASCII_MAX uint16 = 0x7E // '~'
)

var sourceCodeEscapers_HEX_DIGITS = jutil.StringToUTF16("0123456789abcdef")

var sourceCodeEscapers_JAVASCRIPT_ESCAPER = newSourceCodeEscapers_JavaScriptEscaper()

// SourceCodeEscapers_JavascriptEscaper returns an Escaper instance that
// replaces non-ASCII characters in a string with their equivalent
// Javascript UTF-16 escape sequences "\unnnn", "\xnn" or special
// replacement sequences "\b", "\t", "\n", "\f", "\r" or "\\".
//
// Warning: This escaper is not suitable for JSON. JSON users may wish to
// use GSON or other high-level APIs when possible.
//
// Java: SourceCodeEscapers.javascriptEscaper()
func SourceCodeEscapers_JavascriptEscaper() *SourceCodeEscapers_JavaScriptEscaper {
	return sourceCodeEscapers_JAVASCRIPT_ESCAPER
}

// SourceCodeEscapers_AppendWithJavascriptEscaper uses javascriptEscaper()
// to append the escaped transformation of c to to. Using this method can be
// more memory efficient than calling Escaper#escape(String).
//
// Java: SourceCodeEscapers.appendWithJavascriptEscaper(CharSequence, Appendable)
func SourceCodeEscapers_AppendWithJavascriptEscaper(c string, to jutil.Appendable) {
	sourceCodeEscapers_JAVASCRIPT_ESCAPER.AppendTo(c, to)
}

// SourceCodeEscapers_JavaScriptEscaper is an Escaper for javascript
// strings. Turns all non-ASCII characters into ASCII javascript escape
// sequences.
type SourceCodeEscapers_JavaScriptEscaper struct{}

var sourceCodeEscapers_JavaScriptEscaper_REPLACEMENT_CHARS [][]uint16

func init() {
	hexDigits := sourceCodeEscapers_HEX_DIGITS
	replacement := make([][]uint16, 0x100)
	for i := uint16(0); i < sourceCodeEscapers_PRINTABLE_ASCII_MIN; i++ {
		c := i
		r := make([]uint16, 4)
		r[3] = hexDigits[c&0xF]
		c = c >> 4
		r[2] = hexDigits[c&0xF]
		r[1] = 'x'
		r[0] = '\\'
		replacement[i] = r
	}
	for i := sourceCodeEscapers_PRINTABLE_ASCII_MAX + 1; i < 0x100; i++ {
		c := i
		r := make([]uint16, 4)
		r[3] = hexDigits[c&0xF]
		c = c >> 4
		r[2] = hexDigits[c&0xF]
		r[1] = 'x'
		r[0] = '\\'
		replacement[i] = r
	}
	replacement['\''] = jutil.StringToUTF16("\\x27")
	replacement['"'] = jutil.StringToUTF16("\\x22")
	replacement['<'] = jutil.StringToUTF16("\\x3c")
	replacement['='] = jutil.StringToUTF16("\\x3d")
	replacement['>'] = jutil.StringToUTF16("\\x3e")
	replacement['&'] = jutil.StringToUTF16("\\x26")
	replacement['\b'] = jutil.StringToUTF16("\\b")
	replacement['\t'] = jutil.StringToUTF16("\\t")
	replacement['\n'] = jutil.StringToUTF16("\\n")
	replacement['\f'] = jutil.StringToUTF16("\\f")
	replacement['\r'] = jutil.StringToUTF16("\\r")
	replacement['\\'] = jutil.StringToUTF16("\\\\")
	sourceCodeEscapers_JavaScriptEscaper_REPLACEMENT_CHARS = replacement
}

func newSourceCodeEscapers_JavaScriptEscaper() *SourceCodeEscapers_JavaScriptEscaper {
	return &SourceCodeEscapers_JavaScriptEscaper{}
}

// Java: SourceCodeEscapers.JavaScriptEscaper.appendTo(CharSequence, Appendable)
func (e *SourceCodeEscapers_JavaScriptEscaper) AppendTo(cs string, to jutil.Appendable) {
	units := jutil.StringToUTF16(cs)
	last := 0
	length := len(units)
	for i := 0; i < length; i++ {
		c := units[i]
		var replacement []uint16
		if c < 0x100 {
			replacement = sourceCodeEscapers_JavaScriptEscaper_REPLACEMENT_CHARS[c]
			if replacement == nil {
				continue
			}
		} else {
			replacement = sourceCodeEscapers_JavaScriptEscaper_asUnicodeHexEscape(c)
		}
		if last < i {
			to.AppendCharSequence(jutil.StringFromUTF16(units[last:i]))
		}
		for _, r := range replacement {
			to.AppendCharacter(r)
		}
		last = i + 1
	}
	if last < length {
		to.AppendCharSequence(jutil.StringFromUTF16(units[last:length]))
	}
}

// sourceCodeEscapers_JavaScriptEscaper_asUnicodeHexEscape is the helper for
// the common case of escaping a single char.
//
// Java: SourceCodeEscapers.JavaScriptEscaper.asUnicodeHexEscape(char)
func sourceCodeEscapers_JavaScriptEscaper_asUnicodeHexEscape(c uint16) []uint16 {
	hexDigits := sourceCodeEscapers_HEX_DIGITS
	// Equivalent to String.format("\\u%04x", (int) c);
	r := make([]uint16, 6)
	r[0] = '\\'
	r[1] = 'u'
	r[5] = hexDigits[c&0xF]
	c = c >> 4
	r[4] = hexDigits[c&0xF]
	c = c >> 4
	r[3] = hexDigits[c&0xF]
	c = c >> 4
	r[2] = hexDigits[c&0xF]
	return r
}

// Java: SourceCodeEscapers.JavaScriptEscaper.escape(String)
func (e *SourceCodeEscapers_JavaScriptEscaper) Escape(s string) string {
	sb := jutil.NewStringBuilder()
	// StringBuilder.append doesn't actually throw IOException.
	e.AppendTo(s, sb)
	return sb.String()
}
