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
// Port of com/google/javascript/jscomp/CustomPassExecutionTime.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java enum: com.google.javascript.jscomp.CustomPassExecutionTime
type CustomPassExecutionTime uint8

const (
	enum_CustomPassExecutionTime_BEFORE_CHECKS CustomPassExecutionTime = 1
	enum_CustomPassExecutionTime_BEFORE_OPTIMIZATIONS CustomPassExecutionTime = 2
	enum_CustomPassExecutionTime_BEFORE_OPTIMIZATION_LOOP CustomPassExecutionTime = 3
	enum_CustomPassExecutionTime_AFTER_OPTIMIZATION_LOOP CustomPassExecutionTime = 4
)

var enumBaseTable_CustomPassExecutionTime = [...]j2gort.EnumBase{
	{Name_: "BEFORE_CHECKS", Ordinal_: 0},
	{Name_: "BEFORE_OPTIMIZATIONS", Ordinal_: 1},
	{Name_: "BEFORE_OPTIMIZATION_LOOP", Ordinal_: 2},
	{Name_: "AFTER_OPTIMIZATION_LOOP", Ordinal_: 3},
}
var enumValues_CustomPassExecutionTime = []CustomPassExecutionTime{1, 2, 3, 4}

func StaticGet_CustomPassExecutionTime_BEFORE_CHECKS() CustomPassExecutionTime {
	return enum_CustomPassExecutionTime_BEFORE_CHECKS
}
func StaticGet_CustomPassExecutionTime_BEFORE_OPTIMIZATIONS() CustomPassExecutionTime {
	return enum_CustomPassExecutionTime_BEFORE_OPTIMIZATIONS
}
func StaticGet_CustomPassExecutionTime_BEFORE_OPTIMIZATION_LOOP() CustomPassExecutionTime {
	return enum_CustomPassExecutionTime_BEFORE_OPTIMIZATION_LOOP
}
func StaticGet_CustomPassExecutionTime_AFTER_OPTIMIZATION_LOOP() CustomPassExecutionTime {
	return enum_CustomPassExecutionTime_AFTER_OPTIMIZATION_LOOP
}

func CustomPassExecutionTime_Values() []CustomPassExecutionTime {
	return j2gort.CloneSlice(enumValues_CustomPassExecutionTime)
}

func CustomPassExecutionTime_ValueOf(name string) CustomPassExecutionTime {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_CustomPassExecutionTime {
		if enumBaseTable_CustomPassExecutionTime[i_].Name_ == name {
			return CustomPassExecutionTime(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.CustomPassExecutionTime." + name))
}

func (this_ CustomPassExecutionTime) IsJavaNil_() bool { return this_ == 0 }
func (CustomPassExecutionTime) HasIdentityEquals_() {}
func (CustomPassExecutionTime) Implements_j2gort_Ext_java_lang_Enum() {}
func (CustomPassExecutionTime) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ CustomPassExecutionTime) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase { return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()} }
func (CustomPassExecutionTime) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.CustomPassExecutionTime", "CustomPassExecutionTime" }
func (this_ CustomPassExecutionTime) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_CustomPassExecutionTime[this_-1]
}
func (this_ CustomPassExecutionTime) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ CustomPassExecutionTime) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ CustomPassExecutionTime) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CustomPassExecutionTime[this_-1].Name_
}
func (this_ CustomPassExecutionTime) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_CustomPassExecutionTime[this_-1].Name_
}
func (this_ CustomPassExecutionTime) Equals(other_ any) bool {
	o_, ok_ := other_.(CustomPassExecutionTime)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ CustomPassExecutionTime) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_CustomPassExecutionTime[this_-1].HashCode()
}
func (this_ CustomPassExecutionTime) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CustomPassExecutionTime[this_-1].CompareTo(other_)
}
func (this_ CustomPassExecutionTime) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_CustomPassExecutionTime[this_-1].Name_)
}
func init() {
	vals_ := []any{CustomPassExecutionTime(1), CustomPassExecutionTime(2), CustomPassExecutionTime(3), CustomPassExecutionTime(4)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CustomPassExecutionTime", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CustomPassExecutionTime", "com.google.javascript.jscomp.CustomPassExecutionTime", func(v_ any) bool {
		return j2gort.InstanceOf[CustomPassExecutionTime](v_)
	})
}

