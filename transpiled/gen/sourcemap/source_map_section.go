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
// Port of com/google/debugging/sourcemap/SourceMapSection.java (generated mechanically by j2go).

package sourcemap

import (
	"fmt"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java record: com.google.debugging.sourcemap.SourceMapSection
type SourceMapSection struct {
	value string `gson:"value"`
	line int32 `gson:"line"`
	column int32 `gson:"column"`
	type_ SourceMapSection_SectionType `gson:"type"`
}

func (this_ *SourceMapSection) BindSelf_(s any) {
	_ = s
}
func (this_ *SourceMapSection) IsJavaNil_() bool { return this_ == nil }
func (*SourceMapSection) JavaClassInfo_() (string, string) { return "com.google.debugging.sourcemap.SourceMapSection", "SourceMapSection" }
func init() {
	j2gort.RegisterTypePredicate("com.google.debugging.sourcemap.SourceMapSection", "com.google.debugging.sourcemap.SourceMapSection", func(v_ any) bool {
		return j2gort.InstanceOf[*SourceMapSection](v_)
	})
}
func (*SourceMapSection) Implements_j2gort_Ext_java_lang_Record() {}

func (this_ *SourceMapSection) initFields_SourceMapSection() {
	_ = this_
}

// Java constructor: SourceMapSection(com.google.debugging.sourcemap.SourceMapSection.SectionType type, java.lang.String value, int line, int column)
func NewSourceMapSection(type_ SourceMapSection_SectionType, value string, line int32, column int32) *SourceMapSection {
	this_ := &SourceMapSection{}
	InitSourceMapSection(this_, type_, value, line, column)
	return this_
}

func InitSourceMapSection(this_ *SourceMapSection, type_ SourceMapSection_SectionType, value string, line int32, column int32) {
	_ = type_
	_ = value
	_ = line
	_ = column
	this_.initFields_SourceMapSection()
	this_.type_ = type_
	this_.value = value
	this_.line = line
	this_.column = column
}

// Java static method: SourceMapSection.forMap(java.lang.String value, int line, int column)
func SourceMapSection_ForMap(value string, line int32, column int32) *SourceMapSection {
	_ = value
	_ = line
	_ = column
	return NewSourceMapSection(StaticGet_SourceMapSection_SectionType_MAP(), value, line, column)
}

// Java static method: SourceMapSection.forURL(java.lang.String value, int line, int column)
func SourceMapSection_ForURL(value string, line int32, column int32) *SourceMapSection {
	_ = value
	_ = line
	_ = column
	return NewSourceMapSection(StaticGet_SourceMapSection_SectionType_URL(), value, line, column)
}

// Java method: SourceMapSection.getSectionType()
func (this_ *SourceMapSection) GetSectionType() SourceMapSection_SectionType {
	return this_.type_
}

// Java method: SourceMapSection.getSectionValue()
func (this_ *SourceMapSection) GetSectionValue() string {
	return this_.value
}

// Java method: SourceMapSection.getLine()
func (this_ *SourceMapSection) GetLine() int32 {
	return this_.line
}

// Java method: SourceMapSection.getColumn()
func (this_ *SourceMapSection) GetColumn() int32 {
	return this_.column
}

func (this_ *SourceMapSection) Type() SourceMapSection_SectionType { return this_.type_ }
func (this_ *SourceMapSection) Value() string { return this_.value }
func (this_ *SourceMapSection) Line() int32 { return this_.line }
func (this_ *SourceMapSection) Column() int32 { return this_.column }
func (this_ *SourceMapSection) Equals(other_ any) bool {
	o_, ok_ := other_.(*SourceMapSection)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if this_.type_ != o_.type_ { return false }
	if this_.value != o_.value { return false }
	if this_.line != o_.line { return false }
	if this_.column != o_.column { return false }
	return true
}
func (this_ *SourceMapSection) HashCode() int32 {
	var h_ int32 = 0
	h_ = 31*h_ + (this_.type_).HashCode()
	{ hs_1 := this_.value; h_ = 31 * h_; if hs_1 != j2gort.NullString { h_ += javalang.StringHashCode(hs_1) } }
	h_ = 31*h_ + (this_.line)
	h_ = 31*h_ + (this_.column)
	return h_
}
func (this_ *SourceMapSection) String() string {
	_ = fmt.Sprint
	return "SourceMapSection[" + "type=" + j2gort.StrVal(j2gort.BoxNullableVal(this_.type_)) + ", " + "value=" + j2gort.StrVal(j2gort.BoxString(this_.value)) + ", " + "line=" + j2gort.StrVal(j2gort.ValueOfInt(this_.line)) + ", " + "column=" + j2gort.StrVal(j2gort.ValueOfInt(this_.column)) + "]"
}

// Java enum: com.google.debugging.sourcemap.SourceMapSection.SectionType
type SourceMapSection_SectionType uint8

const (
	enum_SourceMapSection_SectionType_URL SourceMapSection_SectionType = 1
	enum_SourceMapSection_SectionType_MAP SourceMapSection_SectionType = 2
)

var enumBaseTable_SourceMapSection_SectionType = [...]j2gort.EnumBase{
	{Name_: "URL", Ordinal_: 0},
	{Name_: "MAP", Ordinal_: 1},
}
var enumValues_SourceMapSection_SectionType = []SourceMapSection_SectionType{1, 2}

func StaticGet_SourceMapSection_SectionType_URL() SourceMapSection_SectionType {
	return enum_SourceMapSection_SectionType_URL
}
func StaticGet_SourceMapSection_SectionType_MAP() SourceMapSection_SectionType {
	return enum_SourceMapSection_SectionType_MAP
}

func SourceMapSection_SectionType_Values() []SourceMapSection_SectionType {
	return j2gort.CloneSlice(enumValues_SourceMapSection_SectionType)
}

func SourceMapSection_SectionType_ValueOf(name string) SourceMapSection_SectionType {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_SourceMapSection_SectionType {
		if enumBaseTable_SourceMapSection_SectionType[i_].Name_ == name {
			return SourceMapSection_SectionType(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.debugging.sourcemap.SourceMapSection.SectionType." + name))
}

func (this_ SourceMapSection_SectionType) IsJavaNil_() bool { return this_ == 0 }
func (SourceMapSection_SectionType) HasIdentityEquals_() {}
func (SourceMapSection_SectionType) Implements_j2gort_Ext_java_lang_Enum() {}
func (SourceMapSection_SectionType) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ SourceMapSection_SectionType) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase { return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()} }
func (SourceMapSection_SectionType) JavaClassInfo_() (string, string) { return "com.google.debugging.sourcemap.SourceMapSection$SectionType", "SectionType" }
func (this_ SourceMapSection_SectionType) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_SourceMapSection_SectionType[this_-1]
}
func (this_ SourceMapSection_SectionType) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ SourceMapSection_SectionType) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ SourceMapSection_SectionType) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_SourceMapSection_SectionType[this_-1].Name_
}
func (this_ SourceMapSection_SectionType) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_SourceMapSection_SectionType[this_-1].Name_
}
func (this_ SourceMapSection_SectionType) Equals(other_ any) bool {
	o_, ok_ := other_.(SourceMapSection_SectionType)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ SourceMapSection_SectionType) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_SourceMapSection_SectionType[this_-1].HashCode()
}
func (this_ SourceMapSection_SectionType) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_SourceMapSection_SectionType[this_-1].CompareTo(other_)
}
func (this_ SourceMapSection_SectionType) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_SourceMapSection_SectionType[this_-1].Name_)
}
func init() {
	vals_ := []any{SourceMapSection_SectionType(1), SourceMapSection_SectionType(2)}
	j2gort.RegisterEnumConstants("com.google.debugging.sourcemap.SourceMapSection.SectionType", vals_)
	j2gort.RegisterEnumConstants("com.google.debugging.sourcemap.SourceMapSection$SectionType", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.debugging.sourcemap.SourceMapSection.SectionType", "com.google.debugging.sourcemap.SourceMapSection$SectionType", func(v_ any) bool {
		return j2gort.InstanceOf[SourceMapSection_SectionType](v_)
	})
}

