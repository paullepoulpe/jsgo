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
// Port of com/google/javascript/rhino/jstype/SimpleSlot.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java polymorphic class interface: com.google.javascript.rhino.jstype.SimpleSlot
type SimpleSlot interface {
	Implements_rhino_SimpleSlot()
	IsJavaNil_() bool
	GetSimpleSlotBase() *SimpleSlotBase
	StaticTypedSlot
}
func ToIface_SimpleSlot[T interface {
	SimpleSlot
	comparable
}](v T) SimpleSlot {
	var z T
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.rhino.jstype.SimpleSlot
type SimpleSlotBase struct {
	name string `gson:"name"`
	type_ JSType `gson:"type"`
	idHash_ atomic.Int32
	inferred bool `gson:"inferred"`
}

func (this_ *SimpleSlotBase) BindSelf_(s any) {
	_ = s
}
func (this_ *SimpleSlotBase) IsJavaNil_() bool { return this_ == nil }
func (*SimpleSlotBase) JavaClassInfo_() (string, string) { return "com.google.javascript.rhino.jstype.SimpleSlot", "SimpleSlot" }
func (this_ *SimpleSlotBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.SimpleSlot", "com.google.javascript.rhino.jstype.SimpleSlot", func(v_ any) bool {
		return j2gort.InstanceOf[SimpleSlot](v_)
	})
}
func (b *SimpleSlotBase) GetSimpleSlotBase() *SimpleSlotBase { return b }
func (*SimpleSlotBase) Implements_rhino_SimpleSlot() {}
func (*SimpleSlotBase) Implements_rhino_StaticTypedSlot() {}
func (*SimpleSlotBase) Implements_rhino_StaticSlot() {}

func StaticGet_SimpleSlot_SerialVersionUID() int64 { return int64(1) }
func (this_ *SimpleSlotBase) initFields_SimpleSlot() {
	_ = this_
	this_.name = j2gort.NullString
}

// Java constructor: SimpleSlot(java.lang.String name, com.google.javascript.rhino.jstype.JSType type, boolean inferred)
func NewSimpleSlot(name string, type_ JSType, inferred bool) SimpleSlot {
	this_ := &SimpleSlotBase{}
	InitSimpleSlot(this_, name, type_, inferred)
	return this_
}

func InitSimpleSlot(this_ *SimpleSlotBase, name string, type_ JSType, inferred bool) {
	_ = name
	_ = type_
	_ = inferred
	this_.initFields_SimpleSlot()
	this_.name = name
	this_.type_ = type_
	this_.inferred = inferred
}

// Java method: SimpleSlot.getName()
func (this_ *SimpleSlotBase) GetName() string {
	return this_.name
}

// Java method: SimpleSlot.getType()
func (this_ *SimpleSlotBase) GetType() JSType {
	return this_.type_
}

// Java method: SimpleSlot.isTypeInferred()
func (this_ *SimpleSlotBase) IsTypeInferred() bool {
	return this_.inferred
}

// Java method: SimpleSlot.getDeclaration()
func (this_ *SimpleSlotBase) GetDeclaration() any {
	return nil
}

// Java method: SimpleSlot.getJSDocInfo()
func (this_ *SimpleSlotBase) GetJSDocInfo() *JSDocInfo {
	return nil
}

// Java method: SimpleSlot.getScope()
func (this_ *SimpleSlotBase) GetScope() any {
	return nil
}

