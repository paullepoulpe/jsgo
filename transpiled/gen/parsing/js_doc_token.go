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
// Port of com/google/javascript/jscomp/parsing/JsDocToken.java (generated mechanically by j2go).

package parsing

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java enum: com.google.javascript.jscomp.parsing.JsDocToken
type JsDocToken uint8

const (
	enum_JsDocToken_EOF JsDocToken = 1
	enum_JsDocToken_EOL JsDocToken = 2
	enum_JsDocToken_LEFT_ANGLE JsDocToken = 3
	enum_JsDocToken_RIGHT_ANGLE JsDocToken = 4
	enum_JsDocToken_STRING JsDocToken = 5
	enum_JsDocToken_LEFT_SQUARE JsDocToken = 6
	enum_JsDocToken_RIGHT_SQUARE JsDocToken = 7
	enum_JsDocToken_LEFT_CURLY JsDocToken = 8
	enum_JsDocToken_RIGHT_CURLY JsDocToken = 9
	enum_JsDocToken_LEFT_PAREN JsDocToken = 10
	enum_JsDocToken_RIGHT_PAREN JsDocToken = 11
	enum_JsDocToken_COMMA JsDocToken = 12
	enum_JsDocToken_COLON JsDocToken = 13
	enum_JsDocToken_ITER_REST JsDocToken = 14
	enum_JsDocToken_ANNOTATION JsDocToken = 15
	enum_JsDocToken_PIPE JsDocToken = 16
	enum_JsDocToken_STAR JsDocToken = 17
	enum_JsDocToken_EOC JsDocToken = 18
	enum_JsDocToken_QMARK JsDocToken = 19
	enum_JsDocToken_BANG JsDocToken = 20
	enum_JsDocToken_EQUALS JsDocToken = 21
)

var enumBaseTable_JsDocToken = [...]j2gort.EnumBase{
	{Name_: "EOF", Ordinal_: 0},
	{Name_: "EOL", Ordinal_: 1},
	{Name_: "LEFT_ANGLE", Ordinal_: 2},
	{Name_: "RIGHT_ANGLE", Ordinal_: 3},
	{Name_: "STRING", Ordinal_: 4},
	{Name_: "LEFT_SQUARE", Ordinal_: 5},
	{Name_: "RIGHT_SQUARE", Ordinal_: 6},
	{Name_: "LEFT_CURLY", Ordinal_: 7},
	{Name_: "RIGHT_CURLY", Ordinal_: 8},
	{Name_: "LEFT_PAREN", Ordinal_: 9},
	{Name_: "RIGHT_PAREN", Ordinal_: 10},
	{Name_: "COMMA", Ordinal_: 11},
	{Name_: "COLON", Ordinal_: 12},
	{Name_: "ITER_REST", Ordinal_: 13},
	{Name_: "ANNOTATION", Ordinal_: 14},
	{Name_: "PIPE", Ordinal_: 15},
	{Name_: "STAR", Ordinal_: 16},
	{Name_: "EOC", Ordinal_: 17},
	{Name_: "QMARK", Ordinal_: 18},
	{Name_: "BANG", Ordinal_: 19},
	{Name_: "EQUALS", Ordinal_: 20},
}
var enumValues_JsDocToken = []JsDocToken{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21}

func StaticGet_JsDocToken_EOF() JsDocToken {
	return enum_JsDocToken_EOF
}
func StaticGet_JsDocToken_EOL() JsDocToken {
	return enum_JsDocToken_EOL
}
func StaticGet_JsDocToken_LEFT_ANGLE() JsDocToken {
	return enum_JsDocToken_LEFT_ANGLE
}
func StaticGet_JsDocToken_RIGHT_ANGLE() JsDocToken {
	return enum_JsDocToken_RIGHT_ANGLE
}
func StaticGet_JsDocToken_STRING() JsDocToken {
	return enum_JsDocToken_STRING
}
func StaticGet_JsDocToken_LEFT_SQUARE() JsDocToken {
	return enum_JsDocToken_LEFT_SQUARE
}
func StaticGet_JsDocToken_RIGHT_SQUARE() JsDocToken {
	return enum_JsDocToken_RIGHT_SQUARE
}
func StaticGet_JsDocToken_LEFT_CURLY() JsDocToken {
	return enum_JsDocToken_LEFT_CURLY
}
func StaticGet_JsDocToken_RIGHT_CURLY() JsDocToken {
	return enum_JsDocToken_RIGHT_CURLY
}
func StaticGet_JsDocToken_LEFT_PAREN() JsDocToken {
	return enum_JsDocToken_LEFT_PAREN
}
func StaticGet_JsDocToken_RIGHT_PAREN() JsDocToken {
	return enum_JsDocToken_RIGHT_PAREN
}
func StaticGet_JsDocToken_COMMA() JsDocToken {
	return enum_JsDocToken_COMMA
}
func StaticGet_JsDocToken_COLON() JsDocToken {
	return enum_JsDocToken_COLON
}
func StaticGet_JsDocToken_ITER_REST() JsDocToken {
	return enum_JsDocToken_ITER_REST
}
func StaticGet_JsDocToken_ANNOTATION() JsDocToken {
	return enum_JsDocToken_ANNOTATION
}
func StaticGet_JsDocToken_PIPE() JsDocToken {
	return enum_JsDocToken_PIPE
}
func StaticGet_JsDocToken_STAR() JsDocToken {
	return enum_JsDocToken_STAR
}
func StaticGet_JsDocToken_EOC() JsDocToken {
	return enum_JsDocToken_EOC
}
func StaticGet_JsDocToken_QMARK() JsDocToken {
	return enum_JsDocToken_QMARK
}
func StaticGet_JsDocToken_BANG() JsDocToken {
	return enum_JsDocToken_BANG
}
func StaticGet_JsDocToken_EQUALS() JsDocToken {
	return enum_JsDocToken_EQUALS
}

func JsDocToken_Values() []JsDocToken {
	return j2gort.CloneSlice(enumValues_JsDocToken)
}

func JsDocToken_ValueOf(name string) JsDocToken {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_JsDocToken {
		if enumBaseTable_JsDocToken[i_].Name_ == name {
			return JsDocToken(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.parsing.JsDocToken." + name))
}

func (this_ JsDocToken) IsJavaNil_() bool { return this_ == 0 }
func (JsDocToken) HasIdentityEquals_() {}
func (JsDocToken) Implements_j2gort_Ext_java_lang_Enum() {}
func (JsDocToken) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ JsDocToken) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase { return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()} }
func (JsDocToken) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.parsing.JsDocToken", "JsDocToken" }
func (this_ JsDocToken) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_JsDocToken[this_-1]
}
func (this_ JsDocToken) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ JsDocToken) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ JsDocToken) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_JsDocToken[this_-1].Name_
}
func (this_ JsDocToken) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_JsDocToken[this_-1].Name_
}
func (this_ JsDocToken) Equals(other_ any) bool {
	o_, ok_ := other_.(JsDocToken)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ JsDocToken) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_JsDocToken[this_-1].HashCode()
}
func (this_ JsDocToken) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_JsDocToken[this_-1].CompareTo(other_)
}
func (this_ JsDocToken) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_JsDocToken[this_-1].Name_)
}
func init() {
	vals_ := []any{JsDocToken(1), JsDocToken(2), JsDocToken(3), JsDocToken(4), JsDocToken(5), JsDocToken(6), JsDocToken(7), JsDocToken(8), JsDocToken(9), JsDocToken(10), JsDocToken(11), JsDocToken(12), JsDocToken(13), JsDocToken(14), JsDocToken(15), JsDocToken(16), JsDocToken(17), JsDocToken(18), JsDocToken(19), JsDocToken(20), JsDocToken(21)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.JsDocToken", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.JsDocToken", "com.google.javascript.jscomp.parsing.JsDocToken", func(v_ any) bool {
		return j2gort.InstanceOf[JsDocToken](v_)
	})
}

