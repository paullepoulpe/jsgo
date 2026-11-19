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
// Port of com/google/javascript/rhino/JSIdentifier.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.rhino.JSIdentifier
type JSIdentifier struct {
	idHash_ atomic.Int32
}

func (this_ *JSIdentifier) BindSelf_(s any) {
	_ = s
}
func (this_ *JSIdentifier) IsJavaNil_() bool { return this_ == nil }
func (*JSIdentifier) JavaClassInfo_() (string, string) { return "com.google.javascript.rhino.JSIdentifier", "JSIdentifier" }
func (this_ *JSIdentifier) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JSIdentifier) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.JSIdentifier", "com.google.javascript.rhino.JSIdentifier", func(v_ any) bool {
		return j2gort.InstanceOf[*JSIdentifier](v_)
	})
}

func (this_ *JSIdentifier) initFields_JSIdentifier() {
	_ = this_
}

// Java constructor: JSIdentifier()
func NewJSIdentifier() *JSIdentifier {
	this_ := &JSIdentifier{}
	InitJSIdentifier(this_)
	return this_
}

func InitJSIdentifier(this_ *JSIdentifier) {
	this_.initFields_JSIdentifier()
}

// Java static method: JSIdentifier.isJSIdentifier(java.lang.String s)
func JSIdentifier_IsJSIdentifier(s string) bool {
	_ = s
	var length int32 = j2gort.JStrLen(s)
	_ = length
	if ((((length == int32(0)) || j2gort.BoxChar_IsIdentifierIgnorable_1_uint16(j2gort.JStrCharAt(s, int32(0)))) || !(j2gort.BoxChar_IsJavaIdentifierStart_1_uint16(j2gort.JStrCharAt(s, int32(0)))))) {
		return false
	}
	{
		var i int32 = int32(1)
		_ = i
		for ; (i < length); i = i + 1 {
			if ((j2gort.BoxChar_IsIdentifierIgnorable_1_uint16(j2gort.JStrCharAt(s, i)) || !(j2gort.BoxChar_IsJavaIdentifierPart_1_uint16(j2gort.JStrCharAt(s, i))))) {
				return false
			}
		}
	}
	return true
}

