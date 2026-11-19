/*
 * Copyright 2011 The Closure Compiler Authors.
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
// Port of com/google/debugging/sourcemap/SourceMapFormat.java (generated mechanically by j2go).

package sourcemap

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java enum: com.google.debugging.sourcemap.SourceMapFormat
type SourceMapFormat uint8

const (
	enum_SourceMapFormat_DEFAULT SourceMapFormat = 1
	enum_SourceMapFormat_V3 SourceMapFormat = 2
)

var enumBaseTable_SourceMapFormat = [...]j2gort.EnumBase{
	{Name_: "DEFAULT", Ordinal_: 0},
	{Name_: "V3", Ordinal_: 1},
}
var enumValues_SourceMapFormat = []SourceMapFormat{1, 2}

func StaticGet_SourceMapFormat_DEFAULT() SourceMapFormat {
	return enum_SourceMapFormat_DEFAULT
}
func StaticGet_SourceMapFormat_V3() SourceMapFormat {
	return enum_SourceMapFormat_V3
}

func SourceMapFormat_Values() []SourceMapFormat {
	return j2gort.CloneSlice(enumValues_SourceMapFormat)
}

func SourceMapFormat_ValueOf(name string) SourceMapFormat {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_SourceMapFormat {
		if enumBaseTable_SourceMapFormat[i_].Name_ == name {
			return SourceMapFormat(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.debugging.sourcemap.SourceMapFormat." + name))
}

func (this_ SourceMapFormat) IsJavaNil_() bool { return this_ == 0 }
func (SourceMapFormat) HasIdentityEquals_() {}
func (SourceMapFormat) Implements_j2gort_Ext_java_lang_Enum() {}
func (SourceMapFormat) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ SourceMapFormat) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase { return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()} }
func (SourceMapFormat) JavaClassInfo_() (string, string) { return "com.google.debugging.sourcemap.SourceMapFormat", "SourceMapFormat" }
func (this_ SourceMapFormat) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_SourceMapFormat[this_-1]
}
func (this_ SourceMapFormat) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ SourceMapFormat) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ SourceMapFormat) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_SourceMapFormat[this_-1].Name_
}
func (this_ SourceMapFormat) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_SourceMapFormat[this_-1].Name_
}
func (this_ SourceMapFormat) Equals(other_ any) bool {
	o_, ok_ := other_.(SourceMapFormat)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ SourceMapFormat) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_SourceMapFormat[this_-1].HashCode()
}
func (this_ SourceMapFormat) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_SourceMapFormat[this_-1].CompareTo(other_)
}
func (this_ SourceMapFormat) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_SourceMapFormat[this_-1].Name_)
}
func init() {
	vals_ := []any{SourceMapFormat(1), SourceMapFormat(2)}
	j2gort.RegisterEnumConstants("com.google.debugging.sourcemap.SourceMapFormat", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.debugging.sourcemap.SourceMapFormat", "com.google.debugging.sourcemap.SourceMapFormat", func(v_ any) bool {
		return j2gort.InstanceOf[SourceMapFormat](v_)
	})
}

