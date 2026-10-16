// Port of com/google/common/base/Ascii.java (Guava 33.4.6-jre)

/*
 * Copyright (C) 2010 The Guava Authors
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

// Go-only: language-adapter — com.google.common.base.Ascii. Its case
// conversions change only the ASCII letters A-Z / a-z and leave every other
// char unchanged.

// asciiCaseMask is Ascii.CASE_MASK.
const asciiCaseMask = 0x20

// AsciiEqualsIgnoreCase is Ascii.equalsIgnoreCase(CharSequence,
// CharSequence): equality of UTF-16 code units, where only the ASCII letters
// a-z and A-Z are compared case-insensitively.
func AsciiEqualsIgnoreCase(s1, s2 string) bool {
	if s1 == s2 {
		return true
	}
	u1 := StringToUTF16(s1)
	u2 := StringToUTF16(s2)
	if len(u1) != len(u2) {
		return false
	}
	for i := range u1 {
		c1 := u1[i]
		c2 := u2[i]
		if c1 == c2 {
			continue
		}
		alphaIndex := asciiGetAlphaIndex(c1)
		if alphaIndex < 26 && alphaIndex == asciiGetAlphaIndex(c2) {
			continue
		}
		return false
	}
	return true
}

// asciiGetAlphaIndex is Ascii.getAlphaIndex(char): folds upper-case ASCII to
// lower-case and makes it zero-indexed and unsigned (by casting to char).
func asciiGetAlphaIndex(c uint16) uint16 {
	return (c | asciiCaseMask) - 'a'
}

// Ascii_ToLowerCase is Ascii.toLowerCase(String).
func Ascii_ToLowerCase(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] = c + ('a' - 'A')
		}
	}
	return string(b)
}

// Ascii_ToUpperCase is Ascii.toUpperCase(String).
func Ascii_ToUpperCase(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'a' && c <= 'z' {
			b[i] = c - ('a' - 'A')
		}
	}
	return string(b)
}
