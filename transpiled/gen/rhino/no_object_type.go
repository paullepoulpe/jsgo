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
// Port of com/google/javascript/rhino/jstype/NoObjectType.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.rhino.jstype.NoObjectType
type NoObjectType struct {
	idHash_ atomic.Int32
}

func (this_ *NoObjectType) BindSelf_(s any) {
	_ = s
}
func (this_ *NoObjectType) IsJavaNil_() bool { return this_ == nil }
func (*NoObjectType) JavaClassInfo_() (string, string) { return "com.google.javascript.rhino.jstype.NoObjectType", "NoObjectType" }
func (this_ *NoObjectType) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*NoObjectType) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.NoObjectType", "com.google.javascript.rhino.jstype.NoObjectType", func(v_ any) bool {
		return j2gort.InstanceOf[*NoObjectType](v_)
	})
}

func (this_ *NoObjectType) initFields_NoObjectType() {
	_ = this_
}

// Java constructor: NoObjectType()
func NewNoObjectType() *NoObjectType {
	this_ := &NoObjectType{}
	InitNoObjectType(this_)
	return this_
}

func InitNoObjectType(this_ *NoObjectType) {
	this_.initFields_NoObjectType()
}

