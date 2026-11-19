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
// Port of com/google/javascript/rhino/jstype/ContainsUpperBoundSuperTypeVisitor.java (generated mechanically by j2go).

package rhino

import (
	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.rhino.jstype.ContainsUpperBoundSuperTypeVisitor
type ContainsUpperBoundSuperTypeVisitor struct {
	Visitor_WithDefaultCaseBase
	target JSType                                  `gson:"target"`
	seen   *pkg_base.LinkedIdentityHashSet[JSType] `gson:"seen"`
}

func (this_ *ContainsUpperBoundSuperTypeVisitor) BindSelf_(s any) {
	_ = s
	this_.Visitor_WithDefaultCaseBase.BindSelf_(s)
}
func (this_ *ContainsUpperBoundSuperTypeVisitor) IsJavaNil_() bool { return this_ == nil }
func (*ContainsUpperBoundSuperTypeVisitor) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.ContainsUpperBoundSuperTypeVisitor", "ContainsUpperBoundSuperTypeVisitor"
}
func (this_ *ContainsUpperBoundSuperTypeVisitor) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.Visitor_WithDefaultCaseBase.IdentityHash_()
}
func (*ContainsUpperBoundSuperTypeVisitor) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.ContainsUpperBoundSuperTypeVisitor", "com.google.javascript.rhino.jstype.ContainsUpperBoundSuperTypeVisitor", func(v_ any) bool {
		return j2gort.InstanceOf[*ContainsUpperBoundSuperTypeVisitor](v_)
	})
}
func (*ContainsUpperBoundSuperTypeVisitor) Implements_rhino_Visitor_WithDefaultCase() {}
func (*ContainsUpperBoundSuperTypeVisitor) Implements_rhino_Visitor()                 {}

func (this_ *ContainsUpperBoundSuperTypeVisitor) initFields_ContainsUpperBoundSuperTypeVisitor() {
	_ = this_
	this_.seen = pkg_base.NewLinkedIdentityHashSet[JSType]()
}

// Java constructor: ContainsUpperBoundSuperTypeVisitor(com.google.javascript.rhino.jstype.JSType target)
func NewContainsUpperBoundSuperTypeVisitor(target JSType) *ContainsUpperBoundSuperTypeVisitor {
	this_ := &ContainsUpperBoundSuperTypeVisitor{}
	this_.BindSelf_(this_)
	InitContainsUpperBoundSuperTypeVisitor(this_, target)
	return this_
}

func InitContainsUpperBoundSuperTypeVisitor(this_ *ContainsUpperBoundSuperTypeVisitor, target JSType) {
	_ = target
	InitVisitor_WithDefaultCase(&this_.Visitor_WithDefaultCaseBase)
	this_.initFields_ContainsUpperBoundSuperTypeVisitor()
	this_.target = target
}

// Java method: ContainsUpperBoundSuperTypeVisitor.caseDefault(com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType type)
func (this_ *ContainsUpperBoundSuperTypeVisitor) CaseDefault(type_ JSType) any {
	_ = type_
	if this_.target == nil || this_.target.IsJavaNil_() {
		return j2gort.BoxNullableVal(StaticGet_ContainsUpperBoundSuperTypeVisitor_Result_ABSENT())
	}
	if type_ == this_.target || ((type_ == nil || type_.IsJavaNil_()) && (this_.target == nil || this_.target.IsJavaNil_())) {
		return j2gort.BoxNullableVal(StaticGet_ContainsUpperBoundSuperTypeVisitor_Result_PRESENT())
	} else {
		return j2gort.BoxNullableVal(StaticGet_ContainsUpperBoundSuperTypeVisitor_Result_ABSENT())
	}
}

// Java method: ContainsUpperBoundSuperTypeVisitor.caseTemplateType(com.google.javascript.rhino.jstype.TemplateType type)
func (this_ *ContainsUpperBoundSuperTypeVisitor) CaseTemplateType(type_ *TemplateType) any {
	_ = type_
	return j2gort.BoxNullableVal(this_.priv_ContainsUpperBoundSuperTypeVisitor_caseForwardingType_8832(ToIface_JSType(type_), type_.GetBound()))
}

// Java method: ContainsUpperBoundSuperTypeVisitor.caseNamedType(com.google.javascript.rhino.jstype.NamedType type)
func (this_ *ContainsUpperBoundSuperTypeVisitor) CaseNamedType(type_ *NamedType) any {
	_ = type_
	return j2gort.BoxNullableVal(this_.priv_ContainsUpperBoundSuperTypeVisitor_caseForwardingType_8832(ToIface_JSType(type_), type_.GetReferencedType()))
}

// Java method: ContainsUpperBoundSuperTypeVisitor.caseTemplatizedType(com.google.javascript.rhino.jstype.TemplatizedType type)
func (this_ *ContainsUpperBoundSuperTypeVisitor) CaseTemplatizedType(type_ *TemplatizedType) any {
	_ = type_
	return j2gort.BoxNullableVal(this_.priv_ContainsUpperBoundSuperTypeVisitor_caseForwardingType_8832(ToIface_JSType(type_), j2gort.AsIface[JSType](type_.GetReferencedType())))
}

// Java method: ContainsUpperBoundSuperTypeVisitor.caseUnionType(com.google.javascript.rhino.jstype.UnionType type)
func (this_ *ContainsUpperBoundSuperTypeVisitor) CaseUnionType(type_ *UnionType) any {
	_ = type_
	if type_ == this_.target || ((type_ == nil || type_.IsJavaNil_()) && (this_.target == nil || this_.target.IsJavaNil_())) {
		return j2gort.BoxNullableVal(StaticGet_ContainsUpperBoundSuperTypeVisitor_Result_PRESENT())
	}
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = type_.GetAlternatesCount()
		_ = n
		for ; i < n; i = i + 1 {
			var foundInAlt ContainsUpperBoundSuperTypeVisitor_Result = j2gort.Cast[ContainsUpperBoundSuperTypeVisitor_Result](type_.GetAlternateAt(i).Visit_1(this_))
			_ = foundInAlt
			if foundInAlt != StaticGet_ContainsUpperBoundSuperTypeVisitor_Result_ABSENT() {
				return j2gort.BoxNullableVal(foundInAlt)
			}
		}
	}
	return j2gort.BoxNullableVal(StaticGet_ContainsUpperBoundSuperTypeVisitor_Result_ABSENT())
}

// Java method: ContainsUpperBoundSuperTypeVisitor.caseForwardingType(com.google.javascript.rhino.jstype.JSType type, com.google.javascript.rhino.jstype.JSType reference)
func (this_ *ContainsUpperBoundSuperTypeVisitor) priv_ContainsUpperBoundSuperTypeVisitor_caseForwardingType_8832(type_ JSType, reference JSType) ContainsUpperBoundSuperTypeVisitor_Result {
	_ = type_
	_ = reference
	if type_ == this_.target || ((type_ == nil || type_.IsJavaNil_()) && (this_.target == nil || this_.target.IsJavaNil_())) {
		return StaticGet_ContainsUpperBoundSuperTypeVisitor_Result_PRESENT()
	} else if this_.seen.Contains(type_) {
		return StaticGet_ContainsUpperBoundSuperTypeVisitor_Result_CYCLE()
	} else {
		this_.seen.Add(type_)
		return j2gort.Cast[ContainsUpperBoundSuperTypeVisitor_Result](reference.Visit_1(this_))
	}
}

// Java enum: com.google.javascript.rhino.jstype.ContainsUpperBoundSuperTypeVisitor.Result
type ContainsUpperBoundSuperTypeVisitor_Result uint8

const (
	enum_ContainsUpperBoundSuperTypeVisitor_Result_PRESENT ContainsUpperBoundSuperTypeVisitor_Result = 1
	enum_ContainsUpperBoundSuperTypeVisitor_Result_ABSENT  ContainsUpperBoundSuperTypeVisitor_Result = 2
	enum_ContainsUpperBoundSuperTypeVisitor_Result_CYCLE   ContainsUpperBoundSuperTypeVisitor_Result = 3
)

var enumBaseTable_ContainsUpperBoundSuperTypeVisitor_Result = [...]j2gort.EnumBase{
	{Name_: "PRESENT", Ordinal_: 0},
	{Name_: "ABSENT", Ordinal_: 1},
	{Name_: "CYCLE", Ordinal_: 2},
}
var enumValues_ContainsUpperBoundSuperTypeVisitor_Result = []ContainsUpperBoundSuperTypeVisitor_Result{1, 2, 3}

func StaticGet_ContainsUpperBoundSuperTypeVisitor_Result_PRESENT() ContainsUpperBoundSuperTypeVisitor_Result {
	return enum_ContainsUpperBoundSuperTypeVisitor_Result_PRESENT
}
func StaticGet_ContainsUpperBoundSuperTypeVisitor_Result_ABSENT() ContainsUpperBoundSuperTypeVisitor_Result {
	return enum_ContainsUpperBoundSuperTypeVisitor_Result_ABSENT
}
func StaticGet_ContainsUpperBoundSuperTypeVisitor_Result_CYCLE() ContainsUpperBoundSuperTypeVisitor_Result {
	return enum_ContainsUpperBoundSuperTypeVisitor_Result_CYCLE
}

func ContainsUpperBoundSuperTypeVisitor_Result_Values() []ContainsUpperBoundSuperTypeVisitor_Result {
	return j2gort.CloneSlice(enumValues_ContainsUpperBoundSuperTypeVisitor_Result)
}

func ContainsUpperBoundSuperTypeVisitor_Result_ValueOf(name string) ContainsUpperBoundSuperTypeVisitor_Result {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_ContainsUpperBoundSuperTypeVisitor_Result {
		if enumBaseTable_ContainsUpperBoundSuperTypeVisitor_Result[i_].Name_ == name {
			return ContainsUpperBoundSuperTypeVisitor_Result(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.rhino.jstype.ContainsUpperBoundSuperTypeVisitor.Result." + name))
}

func (this_ ContainsUpperBoundSuperTypeVisitor_Result) IsJavaNil_() bool                { return this_ == 0 }
func (ContainsUpperBoundSuperTypeVisitor_Result) HasIdentityEquals_()                   {}
func (ContainsUpperBoundSuperTypeVisitor_Result) Implements_j2gort_Ext_java_lang_Enum() {}
func (ContainsUpperBoundSuperTypeVisitor_Result) Implements_j2gort_Ext_java_lang_constant_Constable() {
}
func (this_ ContainsUpperBoundSuperTypeVisitor_Result) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (ContainsUpperBoundSuperTypeVisitor_Result) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.ContainsUpperBoundSuperTypeVisitor$Result", "Result"
}
func (this_ ContainsUpperBoundSuperTypeVisitor_Result) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_ContainsUpperBoundSuperTypeVisitor_Result[this_-1]
}
func (this_ ContainsUpperBoundSuperTypeVisitor_Result) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ ContainsUpperBoundSuperTypeVisitor_Result) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ ContainsUpperBoundSuperTypeVisitor_Result) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ContainsUpperBoundSuperTypeVisitor_Result[this_-1].Name_
}
func (this_ ContainsUpperBoundSuperTypeVisitor_Result) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_ContainsUpperBoundSuperTypeVisitor_Result[this_-1].Name_
}
func (this_ ContainsUpperBoundSuperTypeVisitor_Result) Equals(other_ any) bool {
	o_, ok_ := other_.(ContainsUpperBoundSuperTypeVisitor_Result)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ ContainsUpperBoundSuperTypeVisitor_Result) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_ContainsUpperBoundSuperTypeVisitor_Result[this_-1].HashCode()
}
func (this_ ContainsUpperBoundSuperTypeVisitor_Result) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ContainsUpperBoundSuperTypeVisitor_Result[this_-1].CompareTo(other_)
}
func (this_ ContainsUpperBoundSuperTypeVisitor_Result) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_ContainsUpperBoundSuperTypeVisitor_Result[this_-1].Name_)
}
func init() {
	vals_ := []any{ContainsUpperBoundSuperTypeVisitor_Result(1), ContainsUpperBoundSuperTypeVisitor_Result(2), ContainsUpperBoundSuperTypeVisitor_Result(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.jstype.ContainsUpperBoundSuperTypeVisitor.Result", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.jstype.ContainsUpperBoundSuperTypeVisitor$Result", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.ContainsUpperBoundSuperTypeVisitor.Result", "com.google.javascript.rhino.jstype.ContainsUpperBoundSuperTypeVisitor$Result", func(v_ any) bool {
		return j2gort.InstanceOf[ContainsUpperBoundSuperTypeVisitor_Result](v_)
	})
}
