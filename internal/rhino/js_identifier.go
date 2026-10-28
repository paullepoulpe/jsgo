// Port of com/google/javascript/rhino/JSIdentifier.java @ closure-compiler v20261021

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Roger Lawrence
 *   Mike McCabe
 *   Igor Bukanov
 *   Ethan Hugg
 *   Bob Jervis
 *   Terry Lucas
 *   Milen Nankov
 *   Pascal-Louis Perez
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino

import "github.com/paullepoulpe/jsgo/internal/jutil"

// JSIdentifier is a utility class to hold isJSIdentifier.
//
// IMPORTANT: As of 2018-03-09 it is still not possible to use Java 8
// features in this file due to limitations on some internal Google projects
// that depend on it.
type JSIdentifier struct{}

var asciiIdentStart = [128]uint8{
	'$': 1,
	'A': 1, 'B': 1, 'C': 1, 'D': 1, 'E': 1, 'F': 1, 'G': 1, 'H': 1, 'I': 1, 'J': 1, 'K': 1, 'L': 1, 'M': 1,
	'N': 1, 'O': 1, 'P': 1, 'Q': 1, 'R': 1, 'S': 1, 'T': 1, 'U': 1, 'V': 1, 'W': 1, 'X': 1, 'Y': 1, 'Z': 1,
	'_': 1,
	'a': 1, 'b': 1, 'c': 1, 'd': 1, 'e': 1, 'f': 1, 'g': 1, 'h': 1, 'i': 1, 'j': 1, 'k': 1, 'l': 1, 'm': 1,
	'n': 1, 'o': 1, 'p': 1, 'q': 1, 'r': 1, 's': 1, 't': 1, 'u': 1, 'v': 1, 'w': 1, 'x': 1, 'y': 1, 'z': 1,
}

var asciiIdentPart = [128]uint8{
	'$': 1,
	'0': 1, '1': 1, '2': 1, '3': 1, '4': 1, '5': 1, '6': 1, '7': 1, '8': 1, '9': 1,
	'A': 1, 'B': 1, 'C': 1, 'D': 1, 'E': 1, 'F': 1, 'G': 1, 'H': 1, 'I': 1, 'J': 1, 'K': 1, 'L': 1, 'M': 1,
	'N': 1, 'O': 1, 'P': 1, 'Q': 1, 'R': 1, 'S': 1, 'T': 1, 'U': 1, 'V': 1, 'W': 1, 'X': 1, 'Y': 1, 'Z': 1,
	'_': 1,
	'a': 1, 'b': 1, 'c': 1, 'd': 1, 'e': 1, 'f': 1, 'g': 1, 'h': 1, 'i': 1, 'j': 1, 'k': 1, 'l': 1, 'm': 1,
	'n': 1, 'o': 1, 'p': 1, 'q': 1, 'r': 1, 's': 1, 't': 1, 'u': 1, 'v': 1, 'w': 1, 'x': 1, 'y': 1, 'z': 1,
}

// Java: JSIdentifier.isJSIdentifier(String)
func JSIdentifier_IsJSIdentifier(s string) bool {
	if len(s) == 0 {
		return false
	}
	c0 := s[0]
	if c0 >= 0x80 {
		return isJSIdentifierUTF16(s)
	}
	if asciiIdentStart[c0] == 0 {
		return false
	}
	for i := 1; i < len(s); i++ {
		c := s[i]
		if c >= 0x80 {
			return isJSIdentifierUTF16(s)
		}
		if asciiIdentPart[c] == 0 {
			return false
		}
	}
	return true
}

//go:noinline
func isJSIdentifierUTF16(s string) bool {
	chars := jutil.StringToUTF16(s)
	length := len(chars)
	if length == 0 ||
		jutil.CharacterIsIdentifierIgnorable(chars[0]) ||
		!jutil.CharacterIsJavaIdentifierStart(chars[0]) {
		return false
	}
	for i := 1; i < length; i++ {
		if jutil.CharacterIsIdentifierIgnorable(chars[i]) ||
			!jutil.CharacterIsJavaIdentifierPart(chars[i]) {
			return false
		}
	}
	return true
}

// Java: JSIdentifier()
func NewJSIdentifier() *JSIdentifier { return &JSIdentifier{} }
