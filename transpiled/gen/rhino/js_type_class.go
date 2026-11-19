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
// Port of com/google/javascript/rhino/jstype/JSTypeClass.java (generated mechanically by j2go).

package rhino

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java enum: com.google.javascript.rhino.jstype.JSTypeClass
type JSTypeClass uint8

const (
	enum_JSTypeClass_ALL JSTypeClass = 1
	enum_JSTypeClass_ARROW JSTypeClass = 2
	enum_JSTypeClass_BOOLEAN JSTypeClass = 3
	enum_JSTypeClass_BIGINT JSTypeClass = 4
	enum_JSTypeClass_ENUM JSTypeClass = 5
	enum_JSTypeClass_ENUM_ELEMENT JSTypeClass = 6
	enum_JSTypeClass_FUNCTION JSTypeClass = 7
	enum_JSTypeClass_INSTANCE_OBJECT JSTypeClass = 8
	enum_JSTypeClass_NAMED JSTypeClass = 9
	enum_JSTypeClass_NO JSTypeClass = 10
	enum_JSTypeClass_NO_OBJECT JSTypeClass = 11
	enum_JSTypeClass_NO_RESOLVED JSTypeClass = 12
	enum_JSTypeClass_NULL JSTypeClass = 13
	enum_JSTypeClass_NUMBER JSTypeClass = 14
	enum_JSTypeClass_PROTOTYPE_OBJECT JSTypeClass = 15
	enum_JSTypeClass_PROXY_OBJECT JSTypeClass = 16
	enum_JSTypeClass_RECORD JSTypeClass = 17
	enum_JSTypeClass_STRING JSTypeClass = 18
	enum_JSTypeClass_SYMBOL JSTypeClass = 19
	enum_JSTypeClass_WELL_KNOWN_SYMBOL JSTypeClass = 20
	enum_JSTypeClass_TEMPLATE JSTypeClass = 21
	enum_JSTypeClass_TEMPLATIZED JSTypeClass = 22
	enum_JSTypeClass_UNION JSTypeClass = 23
	enum_JSTypeClass_UNKNOWN JSTypeClass = 24
	enum_JSTypeClass_VOID JSTypeClass = 25
)

var enumBaseTable_JSTypeClass = [...]j2gort.EnumBase{
	{Name_: "ALL", Ordinal_: 0},
	{Name_: "ARROW", Ordinal_: 1},
	{Name_: "BOOLEAN", Ordinal_: 2},
	{Name_: "BIGINT", Ordinal_: 3},
	{Name_: "ENUM", Ordinal_: 4},
	{Name_: "ENUM_ELEMENT", Ordinal_: 5},
	{Name_: "FUNCTION", Ordinal_: 6},
	{Name_: "INSTANCE_OBJECT", Ordinal_: 7},
	{Name_: "NAMED", Ordinal_: 8},
	{Name_: "NO", Ordinal_: 9},
	{Name_: "NO_OBJECT", Ordinal_: 10},
	{Name_: "NO_RESOLVED", Ordinal_: 11},
	{Name_: "NULL", Ordinal_: 12},
	{Name_: "NUMBER", Ordinal_: 13},
	{Name_: "PROTOTYPE_OBJECT", Ordinal_: 14},
	{Name_: "PROXY_OBJECT", Ordinal_: 15},
	{Name_: "RECORD", Ordinal_: 16},
	{Name_: "STRING", Ordinal_: 17},
	{Name_: "SYMBOL", Ordinal_: 18},
	{Name_: "WELL_KNOWN_SYMBOL", Ordinal_: 19},
	{Name_: "TEMPLATE", Ordinal_: 20},
	{Name_: "TEMPLATIZED", Ordinal_: 21},
	{Name_: "UNION", Ordinal_: 22},
	{Name_: "UNKNOWN", Ordinal_: 23},
	{Name_: "VOID", Ordinal_: 24},
}
var enumValues_JSTypeClass = []JSTypeClass{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25}

func StaticGet_JSTypeClass_ALL() JSTypeClass {
	return enum_JSTypeClass_ALL
}
func StaticGet_JSTypeClass_ARROW() JSTypeClass {
	return enum_JSTypeClass_ARROW
}
func StaticGet_JSTypeClass_BOOLEAN() JSTypeClass {
	return enum_JSTypeClass_BOOLEAN
}
func StaticGet_JSTypeClass_BIGINT() JSTypeClass {
	return enum_JSTypeClass_BIGINT
}
func StaticGet_JSTypeClass_ENUM() JSTypeClass {
	return enum_JSTypeClass_ENUM
}
func StaticGet_JSTypeClass_ENUM_ELEMENT() JSTypeClass {
	return enum_JSTypeClass_ENUM_ELEMENT
}
func StaticGet_JSTypeClass_FUNCTION() JSTypeClass {
	return enum_JSTypeClass_FUNCTION
}
func StaticGet_JSTypeClass_INSTANCE_OBJECT() JSTypeClass {
	return enum_JSTypeClass_INSTANCE_OBJECT
}
func StaticGet_JSTypeClass_NAMED() JSTypeClass {
	return enum_JSTypeClass_NAMED
}
func StaticGet_JSTypeClass_NO() JSTypeClass {
	return enum_JSTypeClass_NO
}
func StaticGet_JSTypeClass_NO_OBJECT() JSTypeClass {
	return enum_JSTypeClass_NO_OBJECT
}
func StaticGet_JSTypeClass_NO_RESOLVED() JSTypeClass {
	return enum_JSTypeClass_NO_RESOLVED
}
func StaticGet_JSTypeClass_NULL() JSTypeClass {
	return enum_JSTypeClass_NULL
}
func StaticGet_JSTypeClass_NUMBER() JSTypeClass {
	return enum_JSTypeClass_NUMBER
}
func StaticGet_JSTypeClass_PROTOTYPE_OBJECT() JSTypeClass {
	return enum_JSTypeClass_PROTOTYPE_OBJECT
}
func StaticGet_JSTypeClass_PROXY_OBJECT() JSTypeClass {
	return enum_JSTypeClass_PROXY_OBJECT
}
func StaticGet_JSTypeClass_RECORD() JSTypeClass {
	return enum_JSTypeClass_RECORD
}
func StaticGet_JSTypeClass_STRING() JSTypeClass {
	return enum_JSTypeClass_STRING
}
func StaticGet_JSTypeClass_SYMBOL() JSTypeClass {
	return enum_JSTypeClass_SYMBOL
}
func StaticGet_JSTypeClass_WELL_KNOWN_SYMBOL() JSTypeClass {
	return enum_JSTypeClass_WELL_KNOWN_SYMBOL
}
func StaticGet_JSTypeClass_TEMPLATE() JSTypeClass {
	return enum_JSTypeClass_TEMPLATE
}
func StaticGet_JSTypeClass_TEMPLATIZED() JSTypeClass {
	return enum_JSTypeClass_TEMPLATIZED
}
func StaticGet_JSTypeClass_UNION() JSTypeClass {
	return enum_JSTypeClass_UNION
}
func StaticGet_JSTypeClass_UNKNOWN() JSTypeClass {
	return enum_JSTypeClass_UNKNOWN
}
func StaticGet_JSTypeClass_VOID() JSTypeClass {
	return enum_JSTypeClass_VOID
}

func JSTypeClass_Values() []JSTypeClass {
	return j2gort.CloneSlice(enumValues_JSTypeClass)
}

func JSTypeClass_ValueOf(name string) JSTypeClass {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_JSTypeClass {
		if enumBaseTable_JSTypeClass[i_].Name_ == name {
			return JSTypeClass(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.rhino.jstype.JSTypeClass." + name))
}

// Java method: JSTypeClass.isTypeOf(com.google.javascript.rhino.jstype.JSType t)
func (this_ JSTypeClass) IsTypeOf(t JSType) bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	_ = t
	return func() bool { re_1 := this_; ae_2 := t.GetTypeClass(); if re_1 == 0 { j2gort.ThrowNPE() }; return re_1 == ae_2 }()
}

func (this_ JSTypeClass) IsJavaNil_() bool { return this_ == 0 }
func (JSTypeClass) HasIdentityEquals_() {}
func (JSTypeClass) Implements_j2gort_Ext_java_lang_Enum() {}
func (JSTypeClass) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ JSTypeClass) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase { return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()} }
func (JSTypeClass) JavaClassInfo_() (string, string) { return "com.google.javascript.rhino.jstype.JSTypeClass", "JSTypeClass" }
func (this_ JSTypeClass) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_JSTypeClass[this_-1]
}
func (this_ JSTypeClass) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ JSTypeClass) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ JSTypeClass) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_JSTypeClass[this_-1].Name_
}
func (this_ JSTypeClass) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_JSTypeClass[this_-1].Name_
}
func (this_ JSTypeClass) Equals(other_ any) bool {
	o_, ok_ := other_.(JSTypeClass)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ JSTypeClass) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_JSTypeClass[this_-1].HashCode()
}
func (this_ JSTypeClass) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_JSTypeClass[this_-1].CompareTo(other_)
}
func (this_ JSTypeClass) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_JSTypeClass[this_-1].Name_)
}
func init() {
	vals_ := []any{JSTypeClass(1), JSTypeClass(2), JSTypeClass(3), JSTypeClass(4), JSTypeClass(5), JSTypeClass(6), JSTypeClass(7), JSTypeClass(8), JSTypeClass(9), JSTypeClass(10), JSTypeClass(11), JSTypeClass(12), JSTypeClass(13), JSTypeClass(14), JSTypeClass(15), JSTypeClass(16), JSTypeClass(17), JSTypeClass(18), JSTypeClass(19), JSTypeClass(20), JSTypeClass(21), JSTypeClass(22), JSTypeClass(23), JSTypeClass(24), JSTypeClass(25)}
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.jstype.JSTypeClass", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.JSTypeClass", "com.google.javascript.rhino.jstype.JSTypeClass", func(v_ any) bool {
		return j2gort.InstanceOf[JSTypeClass](v_)
	})
}

