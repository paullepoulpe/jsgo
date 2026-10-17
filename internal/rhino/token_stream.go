// Port of com/google/javascript/rhino/TokenStream.java @ closure-compiler v20261015

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
// Since isKeyword() is used repeatedly in scanning, for performance reasons,
// it is implemented to do the minimum number of character comparisons to
// ascertain whether the given name is a keyword, instead of doing string
// search of the name against JS keywords.
//
// Java: TokenStream.isKeyword(String)
func TokenStream_IsKeyword(name string) bool {
	id := false
	s := name
	// Go-only: performance — the switch below reads chars only for lengths
	// up to 12, the longest keyword's. An ASCII name (whose bytes are its
	// UTF-16 chars) that short is copied into a stack buffer, and a longer
	// one is not converted at all; any other name is converted as before.
	// This ran for every NAME the IRFactory builds.
	var buf [12]uint16
	var chars []uint16
	length := len(s)
	ascii := true
	for i := 0; i < len(s); i++ {
		if s[i] >= 0x80 {
			ascii = false
			break
		}
	}
	if !ascii {
		chars = jutil.StringToUTF16(s)
		length = len(chars)
	} else if length <= len(buf) {
		for i := 0; i < length; i++ {
			buf[i] = uint16(s[i])
		}
		chars = buf[:length]
	}
	charAt := func(i int) uint16 { return chars[i] } // s.charAt(i)
complete:
	for {
		var X *string
		var c uint16
		set := func(x string) { X = &x; id = true }
	partial:
		switch length {
		case 2:
			c = charAt(1)
			if c == 'f' {
				if charAt(0) == 'i' {
					id = true
					break complete
				}
			} else if c == 'n' {
				if charAt(0) == 'i' {
					id = true
					break complete
				}
			} else if c == 'o' {
				if charAt(0) == 'd' {
					id = true
					break complete
				}
			}
			break partial
		case 3:
			switch charAt(0) {
			case 'f':
				if charAt(2) == 'r' && charAt(1) == 'o' {
					id = true
					break complete
				}
				break partial
			case 'i':
				if charAt(2) == 't' && charAt(1) == 'n' {
					id = true
					break complete
				}
				break partial
			case 'n':
				if charAt(2) == 'w' && charAt(1) == 'e' {
					id = true
					break complete
				}
				break partial
			case 't':
				if charAt(2) == 'y' && charAt(1) == 'r' {
					id = true
					break complete
				}
				break partial
			case 'v':
				if charAt(2) == 'r' && charAt(1) == 'a' {
					id = true
					break complete
				}
				break partial
			}
			break partial
		case 4:
			switch charAt(0) {
			case 'b':
				set("byte")
				break partial
			case 'c':
				c = charAt(3)
				if c == 'e' {
					if charAt(2) == 's' && charAt(1) == 'a' {
						id = true
						break complete
					}
				} else if c == 'r' {
					if charAt(2) == 'a' && charAt(1) == 'h' {
						id = true
						break complete
					}
				}
				break partial
			case 'e':
				c = charAt(3)
				if c == 'e' {
					if charAt(2) == 's' && charAt(1) == 'l' {
						id = true
						break complete
					}
				} else if c == 'm' {
					if charAt(2) == 'u' && charAt(1) == 'n' {
						id = true
						break complete
					}
				}
				break partial
			case 'g':
				set("goto")
				break partial
			case 'l':
				set("long")
				break partial
			case 'n':
				set("null")
				break partial
			case 't':
				c = charAt(3)
				if c == 'e' {
					if charAt(2) == 'u' && charAt(1) == 'r' {
						id = true
						break complete
					}
				} else if c == 's' {
					if charAt(2) == 'i' && charAt(1) == 'h' {
						id = true
						break complete
					}
				}
				break partial
			case 'v':
				set("void")
				break partial
			case 'w':
				set("with")
				break partial
			}
			break partial
		case 5:
			switch charAt(2) {
			case 'a':
				set("class")
				break partial
			case 'e':
				set("break")
				break partial
			case 'i':
				set("while")
				break partial
			case 'l':
				set("false")
				break partial
			case 'n':
				c = charAt(0)
				if c == 'c' {
					set("const")
				} else if c == 'f' {
					set("final")
				}
				break partial
			case 'o':
				c = charAt(0)
				if c == 'f' {
					set("float")
				} else if c == 's' {
					set("short")
				}
				break partial
			case 'p':
				set("super")
				break partial
			case 'r':
				set("throw")
				break partial
			case 't':
				set("catch")
				break partial
			}
			break partial
		case 6:
			switch charAt(1) {
			case 'a':
				set("native")
				break partial
			case 'e':
				c = charAt(0)
				if c == 'd' {
					set("delete")
				} else if c == 'r' {
					set("return")
				}
				break partial
			case 'h':
				set("throws")
				break partial
			case 'm':
				set("import")
				break partial
			case 'o':
				set("double")
				break partial
			case 't':
				set("static")
				break partial
			case 'u':
				set("public")
				break partial
			case 'w':
				set("switch")
				break partial
			case 'x':
				set("export")
				break partial
			case 'y':
				set("typeof")
				break partial
			}
			break partial
		case 7:
			switch charAt(1) {
			case 'a':
				set("package")
				break partial
			case 'e':
				set("default")
				break partial
			case 'i':
				set("finally")
				break partial
			case 'o':
				set("boolean")
				break partial
			case 'r':
				set("private")
				break partial
			case 'x':
				set("extends")
				break partial
			}
			break partial
		case 8:
			switch charAt(0) {
			case 'a':
				set("abstract")
				break partial
			case 'c':
				set("continue")
				break partial
			case 'd':
				set("debugger")
				break partial
			case 'f':
				set("function")
				break partial
			case 'v':
				set("volatile")
				break partial
			}
			break partial
		case 9:
			c = charAt(0)
			if c == 'i' {
				set("interface")
			} else if c == 'p' {
				set("protected")
			} else if c == 't' {
				set("transient")
			}
			break partial
		case 10:
			c = charAt(1)
			if c == 'm' {
				set("implements")
			} else if c == 'n' {
				set("instanceof")
			}
			break partial
		case 12:
			set("synchronized")
			break partial
		}
		// partial match validate the entire string the one possibility
		if X != nil && *X != s {
			return false
		}
		break complete
	}
	return id
}

// Java: TokenStream.isJSIdentifier(String)
func TokenStream_IsJSIdentifier(s string) bool {
	return JSIdentifier_IsJSIdentifier(s)
}

// Java: TokenStream()
func NewTokenStream() *TokenStream { return &TokenStream{} }
