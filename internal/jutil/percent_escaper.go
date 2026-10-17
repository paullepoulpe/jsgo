// Port of com/google/common/net/PercentEscaper.java (Guava 33.4.6-jre),
// with escapeSlow and codePointAt of com/google/common/escape/UnicodeEscaper.java
// (both files carry the same header)

/*
 * Copyright (C) 2008 The Guava Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 */

package jutil

import "strconv"

// Go-only: language-adapter — Guava com.google.common.net.PercentEscaper
// (with UnicodeEscaper.escapeSlow/codePointAt): characters outside the safe
// set are UTF-8 encoded and written as %XX with upper-case hex digits.

// PercentEscaper is com.google.common.net.PercentEscaper.
type PercentEscaper struct {
	plusForSpace bool
	safeOctets   []bool
}

var _ Escaper = (*PercentEscaper)(nil)

// NewPercentEscaper is new PercentEscaper(safeChars, plusForSpace).
func NewPercentEscaper(safeChars string, plusForSpace bool) *PercentEscaper {
	for _, c := range StringToUTF16(safeChars) {
		if (c >= '0' && c <= '9') || (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') {
			panic(NewIllegalArgumentExceptionMsg(
				"Alphanumeric characters are always 'safe' and should not be explicitly specified"))
		}
	}
	safeChars += "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	if plusForSpace && StringIndexOfChar(safeChars, ' ') >= 0 {
		panic(NewIllegalArgumentExceptionMsg("plusForSpace cannot be specified when space is a 'safe' character"))
	}
	return &PercentEscaper{plusForSpace: plusForSpace, safeOctets: percentEscaperCreateSafeOctets(safeChars)}
}

func percentEscaperCreateSafeOctets(safeChars string) []bool {
	maxChar := -1
	units := StringToUTF16(safeChars)
	for _, c := range units {
		maxChar = max(int(c), maxChar)
	}
	octets := make([]bool, maxChar+1)
	for _, c := range units {
		octets[c] = true
	}
	return octets
}

const percentEscaperUpperHexDigits = "0123456789ABCDEF"

func (e *PercentEscaper) isSafe(c int) bool { return c < len(e.safeOctets) && e.safeOctets[c] }

// escapeCodePoint is PercentEscaper.escape(int); nil means "unchanged".
func (e *PercentEscaper) escapeCodePoint(cp int) []byte {
	hex := percentEscaperUpperHexDigits
	switch {
	case e.isSafe(cp):
		return nil
	case cp == ' ' && e.plusForSpace:
		return []byte{'+'}
	case cp <= 0x7F:
		return []byte{'%', hex[cp>>4], hex[cp&0xF]}
	case cp <= 0x7ff:
		// Two byte UTF-8 characters [cp >= 0x80 && cp <= 0x7ff]
		dest := []byte{'%', 0, 0, '%', 0, 0}
		dest[5] = hex[cp&0xF]
		cp >>= 4
		dest[4] = hex[0x8|(cp&0x3)]
		cp >>= 2
		dest[2] = hex[cp&0xF]
		cp >>= 4
		dest[1] = hex[0xC|cp]
		return dest
	case cp <= 0xffff:
		// Three byte UTF-8 characters [cp >= 0x800 && cp <= 0xffff]
		dest := []byte{'%', 'E', 0, '%', 0, 0, '%', 0, 0}
		dest[8] = hex[cp&0xF]
		cp >>= 4
		dest[7] = hex[0x8|(cp&0x3)]
		cp >>= 2
		dest[5] = hex[cp&0xF]
		cp >>= 4
		dest[4] = hex[0x8|(cp&0x3)]
		cp >>= 2
		dest[2] = hex[cp]
		return dest
	case cp <= 0x10FFFF:
		dest := []byte{'%', 'F', 0, '%', 0, 0, '%', 0, 0, '%', 0, 0}
		// Four byte UTF-8 characters [cp >= 0xffff && cp <= 0x10ffff]
		dest[11] = hex[cp&0xF]
		cp >>= 4
		dest[10] = hex[0x8|(cp&0x3)]
		cp >>= 2
		dest[8] = hex[cp&0xF]
		cp >>= 4
		dest[7] = hex[0x8|(cp&0x3)]
		cp >>= 2
		dest[5] = hex[cp&0xF]
		cp >>= 4
		dest[4] = hex[0x8|(cp&0x3)]
		cp >>= 2
		dest[2] = hex[cp&0x7]
		return dest
	}
	panic(NewIllegalArgumentExceptionMsg("Invalid unicode character value " + strconv.Itoa(cp)))
}

// Escape is Escaper.escape(String).
func (e *PercentEscaper) Escape(s string) string {
	units := StringToUTF16(s)
	for index, c := range units {
		if !e.isSafe(int(c)) {
			return e.escapeSlow(s, units, index)
		}
	}
	return s
}

// percentEscaperCodePointAt is UnicodeEscaper.codePointAt.
func percentEscaperCodePointAt(s string, seq []uint16, index, end int) int {
	if index < end {
		c1 := seq[index]
		index++
		if c1 < 0xD800 || c1 > 0xDFFF {
			return int(c1)
		} else if c1 <= 0xDBFF {
			if index == end {
				return -int(c1)
			}
			c2 := seq[index]
			if c2 >= 0xDC00 && c2 <= 0xDFFF {
				return 0x10000 + (int(c1)-0xD800)<<10 + (int(c2) - 0xDC00)
			}
			panic(NewIllegalArgumentExceptionMsg("Expected low surrogate but got char '" + StringFromChar(c2) +
				"' with value " + strconv.Itoa(int(c2)) + " at index " + strconv.Itoa(index) + " in '" + s + "'"))
		} else {
			panic(NewIllegalArgumentExceptionMsg("Unexpected low surrogate character '" + StringFromChar(c1) +
				"' with value " + strconv.Itoa(int(c1)) + " at index " + strconv.Itoa(index-1) + " in '" + s + "'"))
		}
	}
	panic(NewIndexOutOfBoundsException("Index exceeds specified range"))
}

func (e *PercentEscaper) escapeSlow(s string, seq []uint16, index int) string {
	end := len(seq)
	var dest []byte
	unescapedChunkStart := 0
	for index < end {
		cp := percentEscaperCodePointAt(s, seq, index, end)
		if cp < 0 {
			panic(NewIllegalArgumentExceptionMsg("Trailing high surrogate at end of input"))
		}
		escaped := e.escapeCodePoint(cp)
		nextIndex := index + 1
		if cp >= 0x10000 {
			nextIndex = index + 2
		}
		if escaped != nil {
			dest = append(dest, StringFromUTF16(seq[unescapedChunkStart:index])...)
			dest = append(dest, escaped...)
			unescapedChunkStart = nextIndex
		}
		// nextEscapeIndex
		index = nextIndex
		for index < end && e.isSafe(int(seq[index])) {
			index++
		}
	}
	dest = append(dest, StringFromUTF16(seq[unescapedChunkStart:end])...)
	return string(dest)
}
