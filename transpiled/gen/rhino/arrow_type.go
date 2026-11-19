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
 *   Norris Boyd
 *   Roger Lawrence
 *   Mike McCabe
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

// Port of com/google/javascript/rhino/jstype/ArrowType.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.rhino.jstype.ArrowType
type ArrowType struct {
	idHash_ atomic.Int32
}

func (this_ *ArrowType) BindSelf_(s any) {
	_ = s
}
func (this_ *ArrowType) IsJavaNil_() bool { return this_ == nil }
func (*ArrowType) JavaClassInfo_() (string, string) { return "com.google.javascript.rhino.jstype.ArrowType", "ArrowType" }
func (this_ *ArrowType) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ArrowType) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.ArrowType", "com.google.javascript.rhino.jstype.ArrowType", func(v_ any) bool {
		return j2gort.InstanceOf[*ArrowType](v_)
	})
}

func (this_ *ArrowType) initFields_ArrowType() {
	_ = this_
}

// Java constructor: ArrowType()
func NewArrowType() *ArrowType {
	this_ := &ArrowType{}
	InitArrowType(this_)
	return this_
}

func InitArrowType(this_ *ArrowType) {
	this_.initFields_ArrowType()
}

