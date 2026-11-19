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
 *   Brock Smickley
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */
// Port of com/google/javascript/rhino/jstype/BigIntType.java (generated mechanically by j2go).

package rhino

import (
	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.rhino.jstype.BigIntType
type BigIntType struct {
	ValueTypeBase
}

func (this_ *BigIntType) BindSelf_(s any) {
	_ = s
	this_.ValueTypeBase.BindSelf_(s)
}
func (this_ *BigIntType) IsJavaNil_() bool { return this_ == nil }
func (*BigIntType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.BigIntType", "BigIntType"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.BigIntType", "com.google.javascript.rhino.jstype.BigIntType", func(v_ any) bool {
		return j2gort.InstanceOf[*BigIntType](v_)
	})
}
func (*BigIntType) Implements_rhino_ValueType() {}
func (*BigIntType) Implements_rhino_JSType()    {}

func (this_ *BigIntType) initFields_BigIntType() {
	_ = this_
}

// Java constructor: BigIntType(com.google.javascript.rhino.jstype.JSTypeRegistry registry)
func NewBigIntType(registry *JSTypeRegistry) *BigIntType {
	this_ := &BigIntType{}
	this_.BindSelf_(this_)
	InitBigIntType(this_, registry)
	return this_
}

func InitBigIntType(this_ *BigIntType, registry *JSTypeRegistry) {
	_ = registry
	InitValueType(&this_.ValueTypeBase, registry, int32(524288))
	this_.initFields_BigIntType()
}

// Java method: BigIntType.getTypeClass()
func (this_ *BigIntType) GetTypeClass() JSTypeClass {
	return StaticGet_JSTypeClass_BIGINT()
}

// Java method: BigIntType.testForEquality(com.google.javascript.rhino.jstype.JSType that)
func (this_ *BigIntType) TestForEquality(that JSType) pkg_base.Tri {
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

// Java method: BigIntType.matchesNumberContext()
func (this_ *BigIntType) MatchesNumberContext() bool {
	return false
}

// Java method: BigIntType.matchesStringContext()
func (this_ *BigIntType) MatchesStringContext() bool {
	return true
}

// Java method: BigIntType.matchesObjectContext()
func (this_ *BigIntType) MatchesObjectContext() bool {
	return true
}

// Java method: BigIntType.getDisplayName()
func (this_ *BigIntType) GetDisplayName() string {
	return "bigint"
}

// Java method: BigIntType.getPossibleToBooleanOutcomes()
func (this_ *BigIntType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return StaticGet_BooleanLiteralSet_BOTH()
}

// Java method: BigIntType.visit(com.google.javascript.rhino.jstype.Visitor<T> visitor)
func (this_ *BigIntType) Visit_1(visitor Visitor) any {
	_ = visitor
	return visitor.CaseBigIntType()
}

// Java method: BigIntType.autoboxesTo()
func (this_ *BigIntType) AutoboxesTo() JSType {
	return this_.GetNativeType(StaticGet_JSTypeNative_BIGINT_OBJECT_TYPE())
}
