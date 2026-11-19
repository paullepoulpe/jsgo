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
// Port of com/google/javascript/rhino/TokenStream.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.rhino.TokenStream
type TokenStream struct {
	idHash_ atomic.Int32
}

func (this_ *TokenStream) BindSelf_(s any) {
	_ = s
}
func (this_ *TokenStream) IsJavaNil_() bool { return this_ == nil }
func (*TokenStream) JavaClassInfo_() (string, string) { return "com.google.javascript.rhino.TokenStream", "TokenStream" }
func (this_ *TokenStream) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TokenStream) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.TokenStream", "com.google.javascript.rhino.TokenStream", func(v_ any) bool {
		return j2gort.InstanceOf[*TokenStream](v_)
	})
}

func (this_ *TokenStream) initFields_TokenStream() {
	_ = this_
}

// Java constructor: TokenStream()
func NewTokenStream() *TokenStream {
	this_ := &TokenStream{}
	InitTokenStream(this_)
	return this_
}

func InitTokenStream(this_ *TokenStream) {
	this_.initFields_TokenStream()
}

// Java static method: TokenStream.isKeyword(java.lang.String name)
func TokenStream_IsKeyword(name string) bool {
	_ = name
	var id bool = false
	_ = id
	var s string = name
	_ = s
	complete:
	switch {
	default:
		var l_X string = j2gort.NullString
		_ = l_X
		var c int32 = int32(0)
		_ = c
		partial:
		switch {
		default:
			switch j2gort.JStrLen(s) {
			case int32(2):
				c = int32(j2gort.JStrCharAt(s, int32(1)))
				if ((c == int32(uint16(102)))) {
					if ((int32(j2gort.JStrCharAt(s, int32(0))) == int32(uint16(105)))) {
						id = true
						break complete
					}
				} else if ((c == int32(uint16(110)))) {
					if ((int32(j2gort.JStrCharAt(s, int32(0))) == int32(uint16(105)))) {
						id = true
						break complete
					}
				} else if ((c == int32(uint16(111)))) {
					if ((int32(j2gort.JStrCharAt(s, int32(0))) == int32(uint16(100)))) {
						id = true
						break complete
					}
				}
				break partial
			case int32(3):
				switch j2gort.JStrCharAt(s, int32(0)) {
				case uint16(102):
					if (((int32(j2gort.JStrCharAt(s, int32(2))) == int32(uint16(114))) && (int32(j2gort.JStrCharAt(s, int32(1))) == int32(uint16(111))))) {
						id = true
						break complete
					}
					break partial
				case uint16(105):
					if (((int32(j2gort.JStrCharAt(s, int32(2))) == int32(uint16(116))) && (int32(j2gort.JStrCharAt(s, int32(1))) == int32(uint16(110))))) {
						id = true
						break complete
					}
					break partial
				case uint16(110):
					if (((int32(j2gort.JStrCharAt(s, int32(2))) == int32(uint16(119))) && (int32(j2gort.JStrCharAt(s, int32(1))) == int32(uint16(101))))) {
						id = true
						break complete
					}
					break partial
				case uint16(116):
					if (((int32(j2gort.JStrCharAt(s, int32(2))) == int32(uint16(121))) && (int32(j2gort.JStrCharAt(s, int32(1))) == int32(uint16(114))))) {
						id = true
						break complete
					}
					break partial
				case uint16(118):
					if (((int32(j2gort.JStrCharAt(s, int32(2))) == int32(uint16(114))) && (int32(j2gort.JStrCharAt(s, int32(1))) == int32(uint16(97))))) {
						id = true
						break complete
					}
					break partial
				}
				break partial
			case int32(4):
				switch j2gort.JStrCharAt(s, int32(0)) {
				case uint16(98):
					l_X = "byte"
					id = true
					break partial
				case uint16(99):
					c = int32(j2gort.JStrCharAt(s, int32(3)))
					if ((c == int32(uint16(101)))) {
						if (((int32(j2gort.JStrCharAt(s, int32(2))) == int32(uint16(115))) && (int32(j2gort.JStrCharAt(s, int32(1))) == int32(uint16(97))))) {
							id = true
							break complete
						}
					} else if ((c == int32(uint16(114)))) {
						if (((int32(j2gort.JStrCharAt(s, int32(2))) == int32(uint16(97))) && (int32(j2gort.JStrCharAt(s, int32(1))) == int32(uint16(104))))) {
							id = true
							break complete
						}
					}
					break partial
				case uint16(101):
					c = int32(j2gort.JStrCharAt(s, int32(3)))
					if ((c == int32(uint16(101)))) {
						if (((int32(j2gort.JStrCharAt(s, int32(2))) == int32(uint16(115))) && (int32(j2gort.JStrCharAt(s, int32(1))) == int32(uint16(108))))) {
							id = true
							break complete
						}
					} else if ((c == int32(uint16(109)))) {
						if (((int32(j2gort.JStrCharAt(s, int32(2))) == int32(uint16(117))) && (int32(j2gort.JStrCharAt(s, int32(1))) == int32(uint16(110))))) {
							id = true
							break complete
						}
					}
					break partial
				case uint16(103):
					l_X = "goto"
					id = true
					break partial
				case uint16(108):
					l_X = "long"
					id = true
					break partial
				case uint16(110):
					l_X = "null"
					id = true
					break partial
				case uint16(116):
					c = int32(j2gort.JStrCharAt(s, int32(3)))
					if ((c == int32(uint16(101)))) {
						if (((int32(j2gort.JStrCharAt(s, int32(2))) == int32(uint16(117))) && (int32(j2gort.JStrCharAt(s, int32(1))) == int32(uint16(114))))) {
							id = true
							break complete
						}
					} else if ((c == int32(uint16(115)))) {
						if (((int32(j2gort.JStrCharAt(s, int32(2))) == int32(uint16(105))) && (int32(j2gort.JStrCharAt(s, int32(1))) == int32(uint16(104))))) {
							id = true
							break complete
						}
					}
					break partial
				case uint16(118):
					l_X = "void"
					id = true
					break partial
				case uint16(119):
					l_X = "with"
					id = true
					break partial
				}
				break partial
			case int32(5):
				switch j2gort.JStrCharAt(s, int32(2)) {
				case uint16(97):
					l_X = "class"
					id = true
					break partial
				case uint16(101):
					l_X = "break"
					id = true
					break partial
				case uint16(105):
					l_X = "while"
					id = true
					break partial
				case uint16(108):
					l_X = "false"
					id = true
					break partial
				case uint16(110):
					c = int32(j2gort.JStrCharAt(s, int32(0)))
					if ((c == int32(uint16(99)))) {
						l_X = "const"
						id = true
					} else if ((c == int32(uint16(102)))) {
						l_X = "final"
						id = true
					}
					break partial
				case uint16(111):
					c = int32(j2gort.JStrCharAt(s, int32(0)))
					if ((c == int32(uint16(102)))) {
						l_X = "float"
						id = true
					} else if ((c == int32(uint16(115)))) {
						l_X = "short"
						id = true
					}
					break partial
				case uint16(112):
					l_X = "super"
					id = true
					break partial
				case uint16(114):
					l_X = "throw"
					id = true
					break partial
				case uint16(116):
					l_X = "catch"
					id = true
					break partial
				}
				break partial
			case int32(6):
				switch j2gort.JStrCharAt(s, int32(1)) {
				case uint16(97):
					l_X = "native"
					id = true
					break partial
				case uint16(101):
					c = int32(j2gort.JStrCharAt(s, int32(0)))
					if ((c == int32(uint16(100)))) {
						l_X = "delete"
						id = true
					} else if ((c == int32(uint16(114)))) {
						l_X = "return"
						id = true
					}
					break partial
				case uint16(104):
					l_X = "throws"
					id = true
					break partial
				case uint16(109):
					l_X = "import"
					id = true
					break partial
				case uint16(111):
					l_X = "double"
					id = true
					break partial
				case uint16(116):
					l_X = "static"
					id = true
					break partial
				case uint16(117):
					l_X = "public"
					id = true
					break partial
				case uint16(119):
					l_X = "switch"
					id = true
					break partial
				case uint16(120):
					l_X = "export"
					id = true
					break partial
				case uint16(121):
					l_X = "typeof"
					id = true
					break partial
				}
				break partial
			case int32(7):
				switch j2gort.JStrCharAt(s, int32(1)) {
				case uint16(97):
					l_X = "package"
					id = true
					break partial
				case uint16(101):
					l_X = "default"
					id = true
					break partial
				case uint16(105):
					l_X = "finally"
					id = true
					break partial
				case uint16(111):
					l_X = "boolean"
					id = true
					break partial
				case uint16(114):
					l_X = "private"
					id = true
					break partial
				case uint16(120):
					l_X = "extends"
					id = true
					break partial
				}
				break partial
			case int32(8):
				switch j2gort.JStrCharAt(s, int32(0)) {
				case uint16(97):
					l_X = "abstract"
					id = true
					break partial
				case uint16(99):
					l_X = "continue"
					id = true
					break partial
				case uint16(100):
					l_X = "debugger"
					id = true
					break partial
				case uint16(102):
					l_X = "function"
					id = true
					break partial
				case uint16(118):
					l_X = "volatile"
					id = true
					break partial
				}
				break partial
			case int32(9):
				c = int32(j2gort.JStrCharAt(s, int32(0)))
				if ((c == int32(uint16(105)))) {
					l_X = "interface"
					id = true
				} else if ((c == int32(uint16(112)))) {
					l_X = "protected"
					id = true
				} else if ((c == int32(uint16(116)))) {
					l_X = "transient"
					id = true
				}
				break partial
			case int32(10):
				c = int32(j2gort.JStrCharAt(s, int32(1)))
				if ((c == int32(uint16(109)))) {
					l_X = "implements"
					id = true
				} else if ((c == int32(uint16(110)))) {
					l_X = "instanceof"
					id = true
				}
				break partial
			case int32(12):
				l_X = "synchronized"
				id = true
				break partial
			}
		}
		if ((((l_X != j2gort.NullString) && (l_X != s)) && !((j2gort.CheckNotNullString(l_X) == s)))) {
			return false
		}
	}
	return id
}

// Java static method: TokenStream.isJSIdentifier(java.lang.String s)
func TokenStream_IsJSIdentifier(s string) bool {
	_ = s
	return JSIdentifier_IsJSIdentifier(s)
}

