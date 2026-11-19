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
// Port of com/google/javascript/rhino/jstype/BooleanType.java (generated mechanically by j2go).

package rhino

import (
	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.rhino.jstype.BooleanType
type BooleanType struct {
	ValueTypeBase
}

func (this_ *BooleanType) BindSelf_(s any) {
	_ = s
	this_.ValueTypeBase.BindSelf_(s)
}
func (this_ *BooleanType) IsJavaNil_() bool { return this_ == nil }
func (*BooleanType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.BooleanType", "BooleanType"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.BooleanType", "com.google.javascript.rhino.jstype.BooleanType", func(v_ any) bool {
		return j2gort.InstanceOf[*BooleanType](v_)
	})
}
func (*BooleanType) Implements_rhino_ValueType() {}
func (*BooleanType) Implements_rhino_JSType()    {}

func (this_ *BooleanType) initFields_BooleanType() {
	_ = this_
}

// Java constructor: BooleanType(com.google.javascript.rhino.jstype.JSTypeRegistry registry)
func NewBooleanType(registry *JSTypeRegistry) *BooleanType {
	this_ := &BooleanType{}
	this_.BindSelf_(this_)
	InitBooleanType(this_, registry)
	return this_
}

func InitBooleanType(this_ *BooleanType, registry *JSTypeRegistry) {
	_ = registry
	InitValueType(&this_.ValueTypeBase, registry, int32(65536))
	this_.initFields_BooleanType()
}

// Java method: BooleanType.getTypeClass()
func (this_ *BooleanType) GetTypeClass() JSTypeClass {
	return StaticGet_JSTypeClass_BOOLEAN()
}

// Java method: BooleanType.testForEquality(com.google.javascript.rhino.jstype.JSType that)
func (this_ *BooleanType) TestForEquality(that JSType) pkg_base.Tri {
	_ = that
	var result pkg_base.Tri = this_.ValueTypeBase.TestForEquality(that)
	_ = result
	if result != 0 {
		return result
	}
	base_that_ := that.GetJSTypeBase()
	if ((base_that_.IsUnknownType() || base_that_.IsSubtypeOf_1(this_.GetNativeType(StaticGet_JSTypeNative_NUMBER_STRING_BOOLEAN()))) || base_that_.IsSubtypeOf_1(this_.GetNativeType(StaticGet_JSTypeNative_BIGINT_TYPE()))) || that.IsObject() {
		return pkg_base.StaticGet_Tri_UNKNOWN()
	}
	return pkg_base.StaticGet_Tri_FALSE()
}

// Java method: BooleanType.matchesNumberContext()
func (this_ *BooleanType) MatchesNumberContext() bool {
	return true
}

// Java method: BooleanType.matchesStringContext()
func (this_ *BooleanType) MatchesStringContext() bool {
	return true
}

// Java method: BooleanType.matchesObjectContext()
func (this_ *BooleanType) MatchesObjectContext() bool {
	return true
}

// Java method: BooleanType.autoboxesTo()
func (this_ *BooleanType) AutoboxesTo() JSType {
	return this_.GetNativeType(StaticGet_JSTypeNative_BOOLEAN_OBJECT_TYPE())
}

// Java method: BooleanType.getDisplayName()
func (this_ *BooleanType) GetDisplayName() string {
	return "boolean"
}

// Java method: BooleanType.getPossibleToBooleanOutcomes()
func (this_ *BooleanType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return StaticGet_BooleanLiteralSet_BOTH()
}

// Java method: BooleanType.visit(com.google.javascript.rhino.jstype.Visitor<T> visitor)
func (this_ *BooleanType) Visit_1(visitor Visitor) any {
	_ = visitor
	return visitor.CaseBooleanType()
}
