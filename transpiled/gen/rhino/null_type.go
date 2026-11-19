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
// Port of com/google/javascript/rhino/jstype/NullType.java (generated mechanically by j2go).

package rhino

import (
	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.rhino.jstype.NullType
type NullType struct {
	ValueTypeBase
}

func (this_ *NullType) BindSelf_(s any) {
	_ = s
	this_.ValueTypeBase.BindSelf_(s)
}
func (this_ *NullType) IsJavaNil_() bool { return this_ == nil }
func (*NullType) JavaClassInfo_() (string, string) { return "com.google.javascript.rhino.jstype.NullType", "NullType" }
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.NullType", "com.google.javascript.rhino.jstype.NullType", func(v_ any) bool {
		return j2gort.InstanceOf[*NullType](v_)
	})
}
func (*NullType) Implements_rhino_ValueType() {}
func (*NullType) Implements_rhino_JSType() {}

func (this_ *NullType) initFields_NullType() {
	_ = this_
}

// Java constructor: NullType(com.google.javascript.rhino.jstype.JSTypeRegistry registry)
func NewNullType(registry *JSTypeRegistry) *NullType {
	this_ := &NullType{}
	this_.BindSelf_(this_)
	InitNullType(this_, registry)
	return this_
}

func InitNullType(this_ *NullType, registry *JSTypeRegistry) {
	_ = registry
	InitValueType(&this_.ValueTypeBase, registry, int32(131072))
	this_.initFields_NullType()
}

// Java method: NullType.getTypeClass()
func (this_ *NullType) GetTypeClass() JSTypeClass {
	return StaticGet_JSTypeClass_NULL()
}

// Java method: NullType.isNullable()
func (this_ *NullType) IsNullable() bool {
	return true
}

// Java method: NullType.matchesNumberContext()
func (this_ *NullType) MatchesNumberContext() bool {
	return true
}

// Java method: NullType.matchesObjectContext()
func (this_ *NullType) MatchesObjectContext() bool {
	return false
}

// Java method: NullType.matchesStringContext()
func (this_ *NullType) MatchesStringContext() bool {
	return true
}

// Java method: NullType.restrictByNotNullOrUndefined()
func (this_ *NullType) RestrictByNotNullOrUndefined() JSType {
	return this_.GetJSTypeBase().registry.GetNativeType(StaticGet_JSTypeNative_NO_TYPE())
}

// Java method: NullType.restrictByNotNull()
func (this_ *NullType) RestrictByNotNull() JSType {
	return this_.GetJSTypeBase().registry.GetNativeType(StaticGet_JSTypeNative_NO_TYPE())
}

// Java method: NullType.testForEquality(com.google.javascript.rhino.jstype.JSType that)
func (this_ *NullType) TestForEquality(that JSType) pkg_base.Tri {
	_ = that
	var result pkg_base.Tri = this_.ValueTypeBase.TestForEquality(that)
	_ = result
	if ((result != 0)) {
		return result
	}
	if ((that.IsNullType() || that.IsVoidType())) {
		return pkg_base.StaticGet_Tri_TRUE()
	}
	if ((that.IsUnknownType() || that.IsNullable())) {
		return pkg_base.StaticGet_Tri_UNKNOWN()
	}
	return pkg_base.StaticGet_Tri_FALSE()
}

// Java method: NullType.getDisplayName()
func (this_ *NullType) GetDisplayName() string {
	return "null"
}

// Java method: NullType.getPossibleToBooleanOutcomes()
func (this_ *NullType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return StaticGet_BooleanLiteralSet_FALSE()
}

// Java method: NullType.visit(com.google.javascript.rhino.jstype.Visitor<T> visitor)
func (this_ *NullType) Visit_1(visitor Visitor) any {
	_ = visitor
	return visitor.CaseNullType()
}

