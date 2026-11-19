/*
 * Copyright 2018 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/parsing/parser/TemplateLiteralToken.java (generated mechanically by j2go).

package parser

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.parsing.parser.TemplateLiteralToken
type TemplateLiteralToken struct {
	errorMessage string `gson:"errorMessage"`
	errorPosition SourcePosition `gson:"errorPosition"`
	idHash_ atomic.Int32
	errorLevel TemplateLiteralToken_ErrorLevel `gson:"errorLevel"`
}

func (this_ *TemplateLiteralToken) BindSelf_(s any) {
	_ = s
}
func (this_ *TemplateLiteralToken) IsJavaNil_() bool { return this_ == nil }
func (*TemplateLiteralToken) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.parsing.parser.TemplateLiteralToken", "TemplateLiteralToken" }
func (this_ *TemplateLiteralToken) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TemplateLiteralToken) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.TemplateLiteralToken", "com.google.javascript.jscomp.parsing.parser.TemplateLiteralToken", func(v_ any) bool {
		return j2gort.InstanceOf[*TemplateLiteralToken](v_)
	})
}

func (this_ *TemplateLiteralToken) initFields_TemplateLiteralToken() {
	_ = this_
	this_.errorMessage = j2gort.NullString
}

// Java constructor: TemplateLiteralToken(java.lang.@org.jspecify.annotations.Nullable String errorMsg, com.google.javascript.jscomp.parsing.parser.TemplateLiteralToken.@org.jspecify.annotations.Nullable ErrorLevel errorLevel, com.google.javascript.jscomp.parsing.parser.util.SourcePosition position)
func NewTemplateLiteralToken(errorMsg string, errorLevel TemplateLiteralToken_ErrorLevel, position SourcePosition) *TemplateLiteralToken {
	this_ := &TemplateLiteralToken{}
	InitTemplateLiteralToken(this_, errorMsg, errorLevel, position)
	return this_
}

func InitTemplateLiteralToken(this_ *TemplateLiteralToken, errorMsg string, errorLevel TemplateLiteralToken_ErrorLevel, position SourcePosition) {
	_ = errorMsg
	_ = errorLevel
	_ = position
	this_.initFields_TemplateLiteralToken()
	this_.errorMessage = errorMsg
	this_.errorLevel = errorLevel
	this_.errorPosition = position
}

// Java enum: com.google.javascript.jscomp.parsing.parser.TemplateLiteralToken.ErrorLevel
type TemplateLiteralToken_ErrorLevel uint8

const (
	enum_TemplateLiteralToken_ErrorLevel_WARNING TemplateLiteralToken_ErrorLevel = 1
	enum_TemplateLiteralToken_ErrorLevel_ERROR TemplateLiteralToken_ErrorLevel = 2
)

var enumBaseTable_TemplateLiteralToken_ErrorLevel = [...]j2gort.EnumBase{
	{Name_: "WARNING", Ordinal_: 0},
	{Name_: "ERROR", Ordinal_: 1},
}
var enumValues_TemplateLiteralToken_ErrorLevel = []TemplateLiteralToken_ErrorLevel{1, 2}

func StaticGet_TemplateLiteralToken_ErrorLevel_WARNING() TemplateLiteralToken_ErrorLevel {
	return enum_TemplateLiteralToken_ErrorLevel_WARNING
}
func StaticGet_TemplateLiteralToken_ErrorLevel_ERROR() TemplateLiteralToken_ErrorLevel {
	return enum_TemplateLiteralToken_ErrorLevel_ERROR
}

func TemplateLiteralToken_ErrorLevel_Values() []TemplateLiteralToken_ErrorLevel {
	return j2gort.CloneSlice(enumValues_TemplateLiteralToken_ErrorLevel)
}

func TemplateLiteralToken_ErrorLevel_ValueOf(name string) TemplateLiteralToken_ErrorLevel {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_TemplateLiteralToken_ErrorLevel {
		if enumBaseTable_TemplateLiteralToken_ErrorLevel[i_].Name_ == name {
			return TemplateLiteralToken_ErrorLevel(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.parsing.parser.TemplateLiteralToken.ErrorLevel." + name))
}

func (this_ TemplateLiteralToken_ErrorLevel) IsJavaNil_() bool { return this_ == 0 }
func (TemplateLiteralToken_ErrorLevel) HasIdentityEquals_() {}
func (TemplateLiteralToken_ErrorLevel) Implements_j2gort_Ext_java_lang_Enum() {}
func (TemplateLiteralToken_ErrorLevel) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ TemplateLiteralToken_ErrorLevel) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase { return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()} }
func (TemplateLiteralToken_ErrorLevel) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.parsing.parser.TemplateLiteralToken$ErrorLevel", "ErrorLevel" }
func (this_ TemplateLiteralToken_ErrorLevel) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_TemplateLiteralToken_ErrorLevel[this_-1]
}
func (this_ TemplateLiteralToken_ErrorLevel) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ TemplateLiteralToken_ErrorLevel) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ TemplateLiteralToken_ErrorLevel) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_TemplateLiteralToken_ErrorLevel[this_-1].Name_
}
func (this_ TemplateLiteralToken_ErrorLevel) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_TemplateLiteralToken_ErrorLevel[this_-1].Name_
}
func (this_ TemplateLiteralToken_ErrorLevel) Equals(other_ any) bool {
	o_, ok_ := other_.(TemplateLiteralToken_ErrorLevel)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ TemplateLiteralToken_ErrorLevel) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_TemplateLiteralToken_ErrorLevel[this_-1].HashCode()
}
func (this_ TemplateLiteralToken_ErrorLevel) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_TemplateLiteralToken_ErrorLevel[this_-1].CompareTo(other_)
}
func (this_ TemplateLiteralToken_ErrorLevel) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_TemplateLiteralToken_ErrorLevel[this_-1].Name_)
}
func init() {
	vals_ := []any{TemplateLiteralToken_ErrorLevel(1), TemplateLiteralToken_ErrorLevel(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.parser.TemplateLiteralToken.ErrorLevel", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.parser.TemplateLiteralToken$ErrorLevel", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.TemplateLiteralToken.ErrorLevel", "com.google.javascript.jscomp.parsing.parser.TemplateLiteralToken$ErrorLevel", func(v_ any) bool {
		return j2gort.InstanceOf[TemplateLiteralToken_ErrorLevel](v_)
	})
}

