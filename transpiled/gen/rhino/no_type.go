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
 *   Bob Jervis
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */
// Port of com/google/javascript/rhino/jstype/NoType.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.rhino.jstype.NoType
type NoType struct {
	idHash_ atomic.Int32
}

func (this_ *NoType) BindSelf_(s any) {
	_ = s
}
func (this_ *NoType) IsJavaNil_() bool { return this_ == nil }
func (*NoType) JavaClassInfo_() (string, string) { return "com.google.javascript.rhino.jstype.NoType", "NoType" }
func (this_ *NoType) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*NoType) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.NoType", "com.google.javascript.rhino.jstype.NoType", func(v_ any) bool {
		return j2gort.InstanceOf[*NoType](v_)
	})
}

func (this_ *NoType) initFields_NoType() {
	_ = this_
}

// Java constructor: NoType()
func NewNoType() *NoType {
	this_ := &NoType{}
	InitNoType(this_)
	return this_
}

func InitNoType(this_ *NoType) {
	this_.initFields_NoType()
}

