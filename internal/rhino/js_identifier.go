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

// Java: JSIdentifier.isJSIdentifier(String)
func JSIdentifier_IsJSIdentifier(s string) bool {
	// Go-only: performance — an ASCII string of up to 64 bytes (whose bytes
	// are its UTF-16 chars) is copied into a stack buffer; any other string
	// is converted on the heap as before. TypedScopeCreator checks every
	// object-literal key this way.
	var buf [64]uint16
	var chars []uint16 // s.charAt(i) is chars[i]
	if len(s) <= len(buf) {
		chars = buf[:len(s)]
		for i := 0; i < len(s); i++ {
			if s[i] >= 0x80 {
				chars = nil
				break
			}
			chars[i] = uint16(s[i])
		}
	}
	if chars == nil {
		chars = jutil.StringToUTF16(s)
	}
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
