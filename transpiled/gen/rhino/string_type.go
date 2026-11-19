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
// Port of com/google/javascript/rhino/jstype/StringType.java (generated mechanically by j2go).

package rhino

import (
	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.rhino.jstype.StringType
type StringType struct {
	ValueTypeBase
}

func (this_ *StringType) BindSelf_(s any) {
	_ = s
	this_.ValueTypeBase.BindSelf_(s)
}
func (this_ *StringType) IsJavaNil_() bool { return this_ == nil }
func (*StringType) JavaClassInfo_() (string, string) { return "com.google.javascript.rhino.jstype.StringType", "StringType" }
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.StringType", "com.google.javascript.rhino.jstype.StringType", func(v_ any) bool {
		return j2gort.InstanceOf[*StringType](v_)
	})
}
func (*StringType) Implements_rhino_ValueType() {}
func (*StringType) Implements_rhino_JSType() {}

func (this_ *StringType) initFields_StringType() {
	_ = this_
}

// Java constructor: StringType(com.google.javascript.rhino.jstype.JSTypeRegistry registry)
func NewStringType(registry *JSTypeRegistry) *StringType {
	this_ := &StringType{}
	this_.BindSelf_(this_)
	InitStringType(this_, registry)
	return this_
}

func InitStringType(this_ *StringType, registry *JSTypeRegistry) {
	_ = registry
	InitValueType(&this_.ValueTypeBase, registry, int32(32768))
	this_.initFields_StringType()
}

// Java method: StringType.getTypeClass()
func (this_ *StringType) GetTypeClass() JSTypeClass {
	return StaticGet_JSTypeClass_STRING()
}

// Java method: StringType.testForEquality(com.google.javascript.rhino.jstype.JSType that)
func (this_ *StringType) TestForEquality(that JSType) pkg_base.Tri {
	_ = that
	var result pkg_base.Tri = this_.ValueTypeBase.TestForEquality(that)
	_ = result
	if ((result != 0)) {
		return result
	}
	if ((((((that.IsUnknownType() || that.IsSubtypeOf_1(this_.GetNativeType(StaticGet_JSTypeNative_OBJECT_TYPE()))) || that.IsSubtypeOf_1(this_.GetNativeType(StaticGet_JSTypeNative_NUMBER_TYPE()))) || that.IsSubtypeOf_1(this_.GetNativeType(StaticGet_JSTypeNative_STRING_TYPE()))) || that.IsSubtypeOf_1(this_.GetNativeType(StaticGet_JSTypeNative_BOOLEAN_TYPE()))) || that.IsSubtypeOf_1(this_.GetNativeType(StaticGet_JSTypeNative_BIGINT_TYPE())))) {
		return pkg_base.StaticGet_Tri_UNKNOWN()
	}
	return pkg_base.StaticGet_Tri_FALSE()
}

// Java method: StringType.matchesNumberContext()
func (this_ *StringType) MatchesNumberContext() bool {
	return true
}

// Java method: StringType.matchesStringContext()
func (this_ *StringType) MatchesStringContext() bool {
	return true
}

// Java method: StringType.matchesObjectContext()
func (this_ *StringType) MatchesObjectContext() bool {
	return true
}

// Java method: StringType.getDisplayName()
func (this_ *StringType) GetDisplayName() string {
	return "string"
}

// Java method: StringType.autoboxesTo()
func (this_ *StringType) AutoboxesTo() JSType {
	return this_.GetNativeType(StaticGet_JSTypeNative_STRING_OBJECT_TYPE())
}

// Java method: StringType.getPossibleToBooleanOutcomes()
func (this_ *StringType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return StaticGet_BooleanLiteralSet_BOTH()
}

// Java method: StringType.visit(com.google.javascript.rhino.jstype.Visitor<T> visitor)
func (this_ *StringType) Visit_1(visitor Visitor) any {
	_ = visitor
	return visitor.CaseStringType()
}

