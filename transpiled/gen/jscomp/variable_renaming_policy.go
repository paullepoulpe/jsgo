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
// Port of com/google/javascript/jscomp/VariableRenamingPolicy.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java enum: com.google.javascript.jscomp.VariableRenamingPolicy
type VariableRenamingPolicy uint8

const (
	enum_VariableRenamingPolicy_OFF VariableRenamingPolicy = 1
	enum_VariableRenamingPolicy_LOCAL VariableRenamingPolicy = 2
	enum_VariableRenamingPolicy_ALL VariableRenamingPolicy = 3
)

var enumBaseTable_VariableRenamingPolicy = [...]j2gort.EnumBase{
	{Name_: "OFF", Ordinal_: 0},
	{Name_: "LOCAL", Ordinal_: 1},
	{Name_: "ALL", Ordinal_: 2},
}
var enumValues_VariableRenamingPolicy = []VariableRenamingPolicy{1, 2, 3}

func StaticGet_VariableRenamingPolicy_OFF() VariableRenamingPolicy {
	return enum_VariableRenamingPolicy_OFF
}
func StaticGet_VariableRenamingPolicy_LOCAL() VariableRenamingPolicy {
	return enum_VariableRenamingPolicy_LOCAL
}
func StaticGet_VariableRenamingPolicy_ALL() VariableRenamingPolicy {
	return enum_VariableRenamingPolicy_ALL
}

func VariableRenamingPolicy_Values() []VariableRenamingPolicy {
	return j2gort.CloneSlice(enumValues_VariableRenamingPolicy)
}

func VariableRenamingPolicy_ValueOf(name string) VariableRenamingPolicy {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_VariableRenamingPolicy {
		if enumBaseTable_VariableRenamingPolicy[i_].Name_ == name {
			return VariableRenamingPolicy(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.VariableRenamingPolicy." + name))
}

func (this_ VariableRenamingPolicy) IsJavaNil_() bool { return this_ == 0 }
func (VariableRenamingPolicy) HasIdentityEquals_() {}
func (VariableRenamingPolicy) Implements_j2gort_Ext_java_lang_Enum() {}
func (VariableRenamingPolicy) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ VariableRenamingPolicy) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase { return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()} }
func (VariableRenamingPolicy) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.VariableRenamingPolicy", "VariableRenamingPolicy" }
func (this_ VariableRenamingPolicy) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_VariableRenamingPolicy[this_-1]
}
func (this_ VariableRenamingPolicy) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ VariableRenamingPolicy) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ VariableRenamingPolicy) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_VariableRenamingPolicy[this_-1].Name_
}
func (this_ VariableRenamingPolicy) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_VariableRenamingPolicy[this_-1].Name_
}
func (this_ VariableRenamingPolicy) Equals(other_ any) bool {
	o_, ok_ := other_.(VariableRenamingPolicy)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ VariableRenamingPolicy) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_VariableRenamingPolicy[this_-1].HashCode()
}
func (this_ VariableRenamingPolicy) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_VariableRenamingPolicy[this_-1].CompareTo(other_)
}
func (this_ VariableRenamingPolicy) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_VariableRenamingPolicy[this_-1].Name_)
}
func init() {
	vals_ := []any{VariableRenamingPolicy(1), VariableRenamingPolicy(2), VariableRenamingPolicy(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.VariableRenamingPolicy", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.VariableRenamingPolicy", "com.google.javascript.jscomp.VariableRenamingPolicy", func(v_ any) bool {
		return j2gort.InstanceOf[VariableRenamingPolicy](v_)
	})
}

