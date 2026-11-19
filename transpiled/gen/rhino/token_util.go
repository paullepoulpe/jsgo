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
// Port of com/google/javascript/rhino/TokenUtil.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.rhino.TokenUtil
type TokenUtil struct {
	idHash_ atomic.Int32
}

func (this_ *TokenUtil) BindSelf_(s any) {
	_ = s
}
func (this_ *TokenUtil) IsJavaNil_() bool { return this_ == nil }
func (*TokenUtil) JavaClassInfo_() (string, string) { return "com.google.javascript.rhino.TokenUtil", "TokenUtil" }
func (this_ *TokenUtil) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TokenUtil) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.TokenUtil", "com.google.javascript.rhino.TokenUtil", func(v_ any) bool {
		return j2gort.InstanceOf[*TokenUtil](v_)
	})
}

func (this_ *TokenUtil) initFields_TokenUtil() {
	_ = this_
}

// Java constructor: TokenUtil()
func NewTokenUtil() *TokenUtil {
	this_ := &TokenUtil{}
	InitTokenUtil(this_)
	return this_
}

func InitTokenUtil(this_ *TokenUtil) {
	this_.initFields_TokenUtil()
}

// Java static method: TokenUtil.isJSSpace(int c)
func TokenUtil_IsJSSpace(c int32) bool {
	_ = c
	if ((c <= int32(127))) {
		return ((((c == int32(32)) || (c == int32(9))) || (c == int32(12))) || (c == int32(11)))
	} else {
		return ((c == int32(160)) || (j2gort.BoxChar_GetType_1_uint16(uint16(c)) == int32(int8(12))))
	}
}

// Java static method: TokenUtil.isJSFormatChar(int c)
func TokenUtil_IsJSFormatChar(c int32) bool {
	_ = c
	return ((c > int32(127)) && (j2gort.BoxChar_GetType_1_int32(c) == int32(int8(16))))
}

// Java static method: TokenUtil.isWhitespace(int c)
func TokenUtil_IsWhitespace(c int32) bool {
	_ = c
	return j2gort.BoxChar_IsWhitespace_1_int32(c)
}

// Java static method: TokenUtil.isStrWhiteSpaceChar(int c)
func TokenUtil_IsStrWhiteSpaceChar(c int32) pkg_base.Tri {
	_ = c
	switch c {
	case int32(uint16(11)):
		return pkg_base.StaticGet_Tri_UNKNOWN()
	case int32(uint16(32)), int32(uint16(10)), int32(uint16(13)), int32(uint16(9)), int32(uint16(160)), int32(uint16(12)), int32(uint16(8232)), int32(uint16(8233)), int32(uint16(65279)):
		return pkg_base.StaticGet_Tri_TRUE()
	default:
		return func() pkg_base.Tri { if ((j2gort.BoxChar_GetType_1_int32(c) == int32(int8(12)))) { return pkg_base.StaticGet_Tri_TRUE() }; return pkg_base.StaticGet_Tri_FALSE() }()
	}
	panic("unreachable switch expression")
}

