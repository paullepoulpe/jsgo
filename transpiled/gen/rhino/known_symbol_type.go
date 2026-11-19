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
 *   Nick Santos
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */
// Port of com/google/javascript/rhino/jstype/KnownSymbolType.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.rhino.jstype.KnownSymbolType
type KnownSymbolType struct {
	SymbolTypeBase
	name string `gson:"name"`
}

func (this_ *KnownSymbolType) BindSelf_(s any) {
	_ = s
	this_.SymbolTypeBase.BindSelf_(s)
}
func (this_ *KnownSymbolType) IsJavaNil_() bool { return this_ == nil }
func (*KnownSymbolType) JavaClassInfo_() (string, string) { return "com.google.javascript.rhino.jstype.KnownSymbolType", "KnownSymbolType" }
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.KnownSymbolType", "com.google.javascript.rhino.jstype.KnownSymbolType", func(v_ any) bool {
		return j2gort.InstanceOf[*KnownSymbolType](v_)
	})
}
func (*KnownSymbolType) Implements_rhino_SymbolType() {}
func (*KnownSymbolType) Implements_rhino_ValueType() {}
func (*KnownSymbolType) Implements_rhino_JSType() {}

var static_KnownSymbolType_TYPE_CLASS JSTypeClass
func StaticGet_KnownSymbolType_TYPE_CLASS() JSTypeClass {
	clinit_KnownSymbolType()
	return static_KnownSymbolType_TYPE_CLASS
}
func StaticSet_KnownSymbolType_TYPE_CLASS(v JSTypeClass) {
	clinit_KnownSymbolType()
	static_KnownSymbolType_TYPE_CLASS = v
}
var clinitDone_KnownSymbolType uint32
var clinitState_KnownSymbolType j2gort.ClinitState
func clinit_KnownSymbolType() {
	if atomic.LoadUint32(&clinitDone_KnownSymbolType) != 0 {
		return
	}
	clinitSlow_KnownSymbolType()
}

func clinitSlow_KnownSymbolType() {
	if !j2gort.EnterClinit(&clinitState_KnownSymbolType, &clinitDone_KnownSymbolType) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_KnownSymbolType, &clinitDone_KnownSymbolType)
	static_KnownSymbolType_TYPE_CLASS = StaticGet_JSTypeClass_WELL_KNOWN_SYMBOL()
}

func (this_ *KnownSymbolType) initFields_KnownSymbolType() {
	_ = this_
	this_.name = j2gort.NullString
}

// Java constructor: KnownSymbolType(com.google.javascript.rhino.jstype.JSTypeRegistry registry, java.lang.String name)
func NewKnownSymbolType(registry *JSTypeRegistry, name string) *KnownSymbolType {
	this_ := &KnownSymbolType{}
	this_.BindSelf_(this_)
	InitKnownSymbolType(this_, registry, name)
	return this_
}

func InitKnownSymbolType(this_ *KnownSymbolType, registry *JSTypeRegistry, name string) {
	_ = registry
	_ = name
	InitSymbolType_2(&this_.SymbolTypeBase, registry, (int32(1048576) | int32(2097152)))
	this_.initFields_KnownSymbolType()
	this_.name = name
}

// Java method: KnownSymbolType.getTypeClass()
func (this_ *KnownSymbolType) GetTypeClass() JSTypeClass {
	return StaticGet_KnownSymbolType_TYPE_CLASS()
}

// Java method: KnownSymbolType.getDisplayName()
func (this_ *KnownSymbolType) GetDisplayName() string {
	return this_.name
}

// Java method: KnownSymbolType.toMaybeKnownSymbolType()
func (this_ *KnownSymbolType) ToMaybeKnownSymbolType() *KnownSymbolType {
	return this_
}

