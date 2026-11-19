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
// Port of com/google/javascript/rhino/jstype/NumberType.java (generated mechanically by j2go).

package rhino

import (
	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.rhino.jstype.NumberType
type NumberType struct {
	ValueTypeBase
}

func (this_ *NumberType) BindSelf_(s any) {
	_ = s
	this_.ValueTypeBase.BindSelf_(s)
}
func (this_ *NumberType) IsJavaNil_() bool { return this_ == nil }
func (*NumberType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.NumberType", "NumberType"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.NumberType", "com.google.javascript.rhino.jstype.NumberType", func(v_ any) bool {
		return j2gort.InstanceOf[*NumberType](v_)
	})
}
func (*NumberType) Implements_rhino_ValueType() {}
func (*NumberType) Implements_rhino_JSType()    {}

func (this_ *NumberType) initFields_NumberType() {
	_ = this_
}

// Java constructor: NumberType(com.google.javascript.rhino.jstype.JSTypeRegistry registry)
func NewNumberType(registry *JSTypeRegistry) *NumberType {
	this_ := &NumberType{}
	this_.BindSelf_(this_)
	InitNumberType(this_, registry)
	return this_
}

func InitNumberType(this_ *NumberType, registry *JSTypeRegistry) {
	_ = registry
	InitValueType(&this_.ValueTypeBase, registry, int32(16384))
	this_.initFields_NumberType()
}

// Java method: NumberType.getTypeClass()
func (this_ *NumberType) GetTypeClass() JSTypeClass {
	return StaticGet_JSTypeClass_NUMBER()
}

// Java method: NumberType.testForEquality(com.google.javascript.rhino.jstype.JSType that)
func (this_ *NumberType) TestForEquality(that JSType) pkg_base.Tri {
	_ = that
	var result pkg_base.Tri = this_.ValueTypeBase.TestForEquality(that)
	_ = result
	if result != 0 {
		return result
	}
	base_that_ := that.GetJSTypeBase()
	if ((((base_that_.IsUnknownType() || base_that_.IsSubtypeOf_1(this_.GetNativeType(StaticGet_JSTypeNative_OBJECT_TYPE()))) || base_that_.IsSubtypeOf_1(this_.GetNativeType(StaticGet_JSTypeNative_NUMBER_TYPE()))) || base_that_.IsSubtypeOf_1(this_.GetNativeType(StaticGet_JSTypeNative_STRING_TYPE()))) || base_that_.IsSubtypeOf_1(this_.GetNativeType(StaticGet_JSTypeNative_BOOLEAN_TYPE()))) || base_that_.IsSubtypeOf_1(this_.GetNativeType(StaticGet_JSTypeNative_BIGINT_TYPE())) {
		return pkg_base.StaticGet_Tri_UNKNOWN()
	}
	return pkg_base.StaticGet_Tri_FALSE()
}

// Java method: NumberType.matchesNumberContext()
func (this_ *NumberType) MatchesNumberContext() bool {
	return true
}

// Java method: NumberType.matchesStringContext()
func (this_ *NumberType) MatchesStringContext() bool {
	return true
}

// Java method: NumberType.matchesObjectContext()
func (this_ *NumberType) MatchesObjectContext() bool {
	return true
}

// Java method: NumberType.getDisplayName()
func (this_ *NumberType) GetDisplayName() string {
	return "number"
}

// Java method: NumberType.getPossibleToBooleanOutcomes()
func (this_ *NumberType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return StaticGet_BooleanLiteralSet_BOTH()
}

// Java method: NumberType.visit(com.google.javascript.rhino.jstype.Visitor<T> visitor)
func (this_ *NumberType) Visit_1(visitor Visitor) any {
	_ = visitor
	return visitor.CaseNumberType()
}

// Java method: NumberType.autoboxesTo()
func (this_ *NumberType) AutoboxesTo() JSType {
	return this_.GetNativeType(StaticGet_JSTypeNative_NUMBER_OBJECT_TYPE())
}
