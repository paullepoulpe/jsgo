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
// Port of com/google/javascript/rhino/jstype/VoidType.java (generated mechanically by j2go).

package rhino

import (
	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.rhino.jstype.VoidType
type VoidType struct {
	ValueTypeBase
}

func (this_ *VoidType) BindSelf_(s any) {
	_ = s
	this_.ValueTypeBase.BindSelf_(s)
}
func (this_ *VoidType) IsJavaNil_() bool { return this_ == nil }
func (*VoidType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.VoidType", "VoidType"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.VoidType", "com.google.javascript.rhino.jstype.VoidType", func(v_ any) bool {
		return j2gort.InstanceOf[*VoidType](v_)
	})
}
func (*VoidType) Implements_rhino_ValueType() {}
func (*VoidType) Implements_rhino_JSType()    {}

func (this_ *VoidType) initFields_VoidType() {
	_ = this_
}

// Java constructor: VoidType(com.google.javascript.rhino.jstype.JSTypeRegistry registry)
func NewVoidType(registry *JSTypeRegistry) *VoidType {
	this_ := &VoidType{}
	this_.BindSelf_(this_)
	InitVoidType(this_, registry)
	return this_
}

func InitVoidType(this_ *VoidType, registry *JSTypeRegistry) {
	_ = registry
	InitValueType(&this_.ValueTypeBase, registry, int32(262144))
	this_.initFields_VoidType()
}

// Java method: VoidType.getTypeClass()
func (this_ *VoidType) GetTypeClass() JSTypeClass {
	return StaticGet_JSTypeClass_VOID()
}

// Java method: VoidType.restrictByNotNullOrUndefined()
func (this_ *VoidType) RestrictByNotNullOrUndefined() JSType {
	return this_.GetJSTypeBase().registry.GetNativeType(StaticGet_JSTypeNative_NO_TYPE())
}

// Java method: VoidType.restrictByNotUndefined()
func (this_ *VoidType) RestrictByNotUndefined() JSType {
	return this_.GetJSTypeBase().registry.GetNativeType(StaticGet_JSTypeNative_NO_TYPE())
}

// Java method: VoidType.testForEquality(com.google.javascript.rhino.jstype.JSType that)
func (this_ *VoidType) TestForEquality(that JSType) pkg_base.Tri {
	_ = that
	{
		var iifeCond_1 bool
		{

			re_1 := pkg_base.StaticGet_Tri_UNKNOWN()
			ae_2 := this_.ValueTypeBase.TestForEquality(that)
			if re_1 == 0 {
				j2gort.ThrowNPE()
			}
			iifeCond_1 = re_1 == ae_2
		}
		if iifeCond_1 {
			return pkg_base.StaticGet_Tri_UNKNOWN()
		}
	}
	if that.IsSubtypeOf_1(this_) || that.IsSubtypeOf_1(this_.GetNativeType(StaticGet_JSTypeNative_NULL_TYPE())) {
		return pkg_base.StaticGet_Tri_TRUE()
	}
	return pkg_base.StaticGet_Tri_FALSE()
}

// Java method: VoidType.matchesNumberContext()
func (this_ *VoidType) MatchesNumberContext() bool {
	return false
}

// Java method: VoidType.matchesObjectContext()
func (this_ *VoidType) MatchesObjectContext() bool {
	return false
}

// Java method: VoidType.matchesStringContext()
func (this_ *VoidType) MatchesStringContext() bool {
	return true
}

// Java method: VoidType.isVoidable()
func (this_ *VoidType) IsVoidable() bool {
	return true
}

// Java method: VoidType.isExplicitlyVoidable()
func (this_ *VoidType) IsExplicitlyVoidable() bool {
	return true
}

// Java method: VoidType.getDisplayName()
func (this_ *VoidType) GetDisplayName() string {
	return "undefined"
}

// Java method: VoidType.getPossibleToBooleanOutcomes()
func (this_ *VoidType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return StaticGet_BooleanLiteralSet_FALSE()
}

// Java method: VoidType.visit(com.google.javascript.rhino.jstype.Visitor<T> visitor)
func (this_ *VoidType) Visit_1(visitor Visitor) any {
	_ = visitor
	return visitor.CaseVoidType()
}
