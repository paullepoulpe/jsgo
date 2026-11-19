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
// Port of com/google/javascript/rhino/jstype/EquivalenceMethod.java (generated mechanically by j2go).

package rhino

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java enum: com.google.javascript.rhino.jstype.EquivalenceMethod
type EquivalenceMethod uint8

const (
	enum_EquivalenceMethod_IDENTITY EquivalenceMethod = 1
	enum_EquivalenceMethod_DATA_FLOW EquivalenceMethod = 2
)

var enumBaseTable_EquivalenceMethod = [...]j2gort.EnumBase{
	{Name_: "IDENTITY", Ordinal_: 0},
	{Name_: "DATA_FLOW", Ordinal_: 1},
}
var enumValues_EquivalenceMethod = []EquivalenceMethod{1, 2}

func StaticGet_EquivalenceMethod_IDENTITY() EquivalenceMethod {
	return enum_EquivalenceMethod_IDENTITY
}
func StaticGet_EquivalenceMethod_DATA_FLOW() EquivalenceMethod {
	return enum_EquivalenceMethod_DATA_FLOW
}

func EquivalenceMethod_Values() []EquivalenceMethod {
	return j2gort.CloneSlice(enumValues_EquivalenceMethod)
}

func EquivalenceMethod_ValueOf(name string) EquivalenceMethod {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_EquivalenceMethod {
		if enumBaseTable_EquivalenceMethod[i_].Name_ == name {
			return EquivalenceMethod(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.rhino.jstype.EquivalenceMethod." + name))
}

func (this_ EquivalenceMethod) IsJavaNil_() bool { return this_ == 0 }
func (EquivalenceMethod) HasIdentityEquals_() {}
func (EquivalenceMethod) Implements_j2gort_Ext_java_lang_Enum() {}
func (EquivalenceMethod) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ EquivalenceMethod) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase { return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()} }
func (EquivalenceMethod) JavaClassInfo_() (string, string) { return "com.google.javascript.rhino.jstype.EquivalenceMethod", "EquivalenceMethod" }
func (this_ EquivalenceMethod) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_EquivalenceMethod[this_-1]
}
func (this_ EquivalenceMethod) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ EquivalenceMethod) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ EquivalenceMethod) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_EquivalenceMethod[this_-1].Name_
}
func (this_ EquivalenceMethod) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_EquivalenceMethod[this_-1].Name_
}
func (this_ EquivalenceMethod) Equals(other_ any) bool {
	o_, ok_ := other_.(EquivalenceMethod)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ EquivalenceMethod) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_EquivalenceMethod[this_-1].HashCode()
}
func (this_ EquivalenceMethod) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_EquivalenceMethod[this_-1].CompareTo(other_)
}
func (this_ EquivalenceMethod) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_EquivalenceMethod[this_-1].Name_)
}
func init() {
	vals_ := []any{EquivalenceMethod(1), EquivalenceMethod(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.jstype.EquivalenceMethod", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.EquivalenceMethod", "com.google.javascript.rhino.jstype.EquivalenceMethod", func(v_ any) bool {
		return j2gort.InstanceOf[EquivalenceMethod](v_)
	})
}

