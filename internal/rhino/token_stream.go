// Port of com/google/javascript/rhino/TokenStream.java @ closure-compiler v20261021

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

// TokenStream implements the JavaScript scanner.
//
// It is based on the C source files jsscan.c and jsscan.h in the jsref
// package.
//
// IMPORTANT: As of 2018-03-09 it is still not possible to use Java 8
// features in this file due to limitations on some internal Google projects
// that depend on it.
type TokenStream struct{}

// TokenStream_IsKeyword reports whether the string is an ES3 keyword.
//
// Keywords for versions of JavaScript after ES3 are not included, because
// the parser would reject valid ES3 code that happened to newer keywords if
// we did that.
//
// Java: TokenStream.isKeyword(String)
func TokenStream_IsKeyword(name string) bool {
	switch name {
	case "do", "if", "in",
		"for", "int", "new", "try", "var",
		"byte", "case", "char", "else", "enum", "goto", "long", "null", "this", "true", "void", "with",
		"break", "catch", "class", "const", "false", "final", "float", "short", "super", "throw", "while",
		"delete", "double", "export", "import", "native", "public", "return", "static", "switch", "throws", "typeof",
		"boolean", "default", "extends", "finally", "package", "private",
		"abstract", "continue", "debugger", "function", "volatile",
		"interface", "protected", "transient",
		"implements", "instanceof",
		"synchronized":
		return true
	default:
		return false
	}
}

// Java: TokenStream.isJSIdentifier(String)
func TokenStream_IsJSIdentifier(s string) bool {
	return JSIdentifier_IsJSIdentifier(s)
}

// Java: TokenStream()
func NewTokenStream() *TokenStream { return &TokenStream{} }
