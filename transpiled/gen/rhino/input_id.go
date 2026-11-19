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
 *   John Lenz
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */
// Port of com/google/javascript/rhino/InputId.java (generated mechanically by j2go).

package rhino

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.rhino.InputId
type InputId struct {
	id string `gson:"id"`
	hashCode int32 `gson:"-"`
}

func (this_ *InputId) BindSelf_(s any) {
	_ = s
}
func (this_ *InputId) IsJavaNil_() bool { return this_ == nil }
func (*InputId) JavaClassInfo_() (string, string) { return "com.google.javascript.rhino.InputId", "InputId" }
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.InputId", "com.google.javascript.rhino.InputId", func(v_ any) bool {
		return j2gort.InstanceOf[*InputId](v_)
	})
}

func StaticGet_InputId_SerialVersionUID() int64 { return int64(1) }
func (this_ *InputId) initFields_InputId() {
	_ = this_
	this_.id = j2gort.NullString
}

// Java constructor: InputId(java.lang.String id)
func NewInputId(id string) *InputId {
	this_ := &InputId{}
	InitInputId(this_, id)
	return this_
}

func InitInputId(this_ *InputId, id string) {
	_ = id
	this_.initFields_InputId()
	this_.id = id
	if (id != j2gort.NullString) {
		this_.hashCode = javalang.StringHashCode(j2gort.CheckNotNullString(id))
	} else {
		this_.hashCode = int32(0)
	}
}

// Java method: InputId.getIdName()
func (this_ *InputId) GetIdName() string {
	return this_.id
}

// Java method: InputId.hashCode()
func (this_ *InputId) HashCode() int32 {
	if ((this_.id == j2gort.NullString)) {
		panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException")))
	}
	return this_.hashCode
}

// Java method: InputId.equals(java.lang.Object obj)
func (this_ *InputId) Equals(arg_1 any) bool {
	obj := arg_1
	_ = obj
	if (j2gort.RefEquals(this_, obj)) {
		return true
	}
	var inputId *InputId
	_ = inputId
	var ok_2 bool
	inputId, ok_2 = j2gort.TryCast[*InputId](obj)
	if !ok_2 {
		return false
	}
	return func() bool { sr_3 := this_.id; se_4 := inputId.id; return j2gort.CheckNotNullString(sr_3) == se_4 }()
}

// Java method: InputId.toString()
func (this_ *InputId) String() string {
	return ("InputId: " + j2gort.StrValString(this_.GetIdName()))
}

