/*
 * Copyright 2009 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
// Port of com/google/javascript/jscomp/PropertyRenamingPolicy.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java enum: com.google.javascript.jscomp.PropertyRenamingPolicy
type PropertyRenamingPolicy uint8

const (
	enum_PropertyRenamingPolicy_OFF PropertyRenamingPolicy = 1
	enum_PropertyRenamingPolicy_ALL_UNQUOTED PropertyRenamingPolicy = 2
)

var enumBaseTable_PropertyRenamingPolicy = [...]j2gort.EnumBase{
	{Name_: "OFF", Ordinal_: 0},
	{Name_: "ALL_UNQUOTED", Ordinal_: 1},
}
var enumValues_PropertyRenamingPolicy = []PropertyRenamingPolicy{1, 2}

func StaticGet_PropertyRenamingPolicy_OFF() PropertyRenamingPolicy {
	return enum_PropertyRenamingPolicy_OFF
}
func StaticGet_PropertyRenamingPolicy_ALL_UNQUOTED() PropertyRenamingPolicy {
	return enum_PropertyRenamingPolicy_ALL_UNQUOTED
}

func PropertyRenamingPolicy_Values() []PropertyRenamingPolicy {
	return j2gort.CloneSlice(enumValues_PropertyRenamingPolicy)
}

func PropertyRenamingPolicy_ValueOf(name string) PropertyRenamingPolicy {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_PropertyRenamingPolicy {
		if enumBaseTable_PropertyRenamingPolicy[i_].Name_ == name {
			return PropertyRenamingPolicy(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.PropertyRenamingPolicy." + name))
}

func (this_ PropertyRenamingPolicy) IsJavaNil_() bool { return this_ == 0 }
func (PropertyRenamingPolicy) HasIdentityEquals_() {}
func (PropertyRenamingPolicy) Implements_j2gort_Ext_java_lang_Enum() {}
func (PropertyRenamingPolicy) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ PropertyRenamingPolicy) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase { return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()} }
func (PropertyRenamingPolicy) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.PropertyRenamingPolicy", "PropertyRenamingPolicy" }
func (this_ PropertyRenamingPolicy) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_PropertyRenamingPolicy[this_-1]
}
func (this_ PropertyRenamingPolicy) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ PropertyRenamingPolicy) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ PropertyRenamingPolicy) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_PropertyRenamingPolicy[this_-1].Name_
}
func (this_ PropertyRenamingPolicy) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_PropertyRenamingPolicy[this_-1].Name_
}
func (this_ PropertyRenamingPolicy) Equals(other_ any) bool {
	o_, ok_ := other_.(PropertyRenamingPolicy)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ PropertyRenamingPolicy) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_PropertyRenamingPolicy[this_-1].HashCode()
}
func (this_ PropertyRenamingPolicy) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_PropertyRenamingPolicy[this_-1].CompareTo(other_)
}
func (this_ PropertyRenamingPolicy) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_PropertyRenamingPolicy[this_-1].Name_)
}
func init() {
	vals_ := []any{PropertyRenamingPolicy(1), PropertyRenamingPolicy(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.PropertyRenamingPolicy", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PropertyRenamingPolicy", "com.google.javascript.jscomp.PropertyRenamingPolicy", func(v_ any) bool {
		return j2gort.InstanceOf[PropertyRenamingPolicy](v_)
	})
}

