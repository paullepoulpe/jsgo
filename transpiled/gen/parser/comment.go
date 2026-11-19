/*
 * Copyright 2014 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/parsing/parser/trees/Comment.java (generated mechanically by j2go).

package parser

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.Comment
type Comment struct {
	value string `gson:"value"`
	location SourceRange `gson:"location"`
	idHash_ atomic.Int32
	type_ Comment_Type `gson:"type"`
	parsed bool `gson:"parsed"`
}

func (this_ *Comment) BindSelf_(s any) {
	_ = s
}
func (this_ *Comment) IsJavaNil_() bool { return this_ == nil }
func (*Comment) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.parsing.parser.trees.Comment", "Comment" }
func (this_ *Comment) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Comment) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.Comment", "com.google.javascript.jscomp.parsing.parser.trees.Comment", func(v_ any) bool {
		return j2gort.InstanceOf[*Comment](v_)
	})
}
func (this_ *Comment) FieldGet_value() string { return this_.value }
func (this_ *Comment) FieldGet_location() SourceRange { return this_.location }
func (this_ *Comment) FieldGet_type_() Comment_Type { return this_.type_ }
func (this_ *Comment) FieldGet_parsed() bool { return this_.parsed }
func (this_ *Comment) FieldSet_parsed(v bool) { this_.parsed = v }

func (this_ *Comment) initFields_Comment() {
	_ = this_
	this_.value = j2gort.NullString
}

// Java constructor: Comment(java.lang.String value, com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.Comment.Type type)
func NewComment(value string, location SourceRange, type_ Comment_Type) *Comment {
	this_ := &Comment{}
	InitComment(this_, value, location, type_)
	return this_
}

func InitComment(this_ *Comment, value string, location SourceRange, type_ Comment_Type) {
	_ = value
	_ = location
	_ = type_
	this_.initFields_Comment()
	this_.value = value
	this_.location = location
	this_.type_ = type_
}

// Java method: Comment.isJsDoc()
func (this_ *Comment) IsJsDoc() bool {
	return (this_.type_ == StaticGet_Comment_Type_JSDOC())
}

// Java method: Comment.getAbsolutePosition()
func (this_ *Comment) GetAbsolutePosition() int32 {
	return this_.location.start.offset
}

// Java method: Comment.getLength()
func (this_ *Comment) GetLength() int32 {
	return (this_.location.end.offset - this_.location.start.offset)
}

// Java enum: com.google.javascript.jscomp.parsing.parser.trees.Comment.Type
type Comment_Type uint8

const (
	enum_Comment_Type_BLOCK Comment_Type = 1
	enum_Comment_Type_LINE Comment_Type = 2
	enum_Comment_Type_JSDOC Comment_Type = 3
	enum_Comment_Type_SHEBANG Comment_Type = 4
	enum_Comment_Type_IMPORTANT Comment_Type = 5
)

var enumBaseTable_Comment_Type = [...]j2gort.EnumBase{
	{Name_: "BLOCK", Ordinal_: 0},
	{Name_: "LINE", Ordinal_: 1},
	{Name_: "JSDOC", Ordinal_: 2},
	{Name_: "SHEBANG", Ordinal_: 3},
	{Name_: "IMPORTANT", Ordinal_: 4},
}
var enumValues_Comment_Type = []Comment_Type{1, 2, 3, 4, 5}

func StaticGet_Comment_Type_BLOCK() Comment_Type {
	return enum_Comment_Type_BLOCK
}
func StaticGet_Comment_Type_LINE() Comment_Type {
	return enum_Comment_Type_LINE
}
func StaticGet_Comment_Type_JSDOC() Comment_Type {
	return enum_Comment_Type_JSDOC
}
func StaticGet_Comment_Type_SHEBANG() Comment_Type {
	return enum_Comment_Type_SHEBANG
}
func StaticGet_Comment_Type_IMPORTANT() Comment_Type {
	return enum_Comment_Type_IMPORTANT
}

func Comment_Type_Values() []Comment_Type {
	return j2gort.CloneSlice(enumValues_Comment_Type)
}

func Comment_Type_ValueOf(name string) Comment_Type {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_Comment_Type {
		if enumBaseTable_Comment_Type[i_].Name_ == name {
			return Comment_Type(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.parsing.parser.trees.Comment.Type." + name))
}

func (this_ Comment_Type) IsJavaNil_() bool { return this_ == 0 }
func (Comment_Type) HasIdentityEquals_() {}
func (Comment_Type) Implements_j2gort_Ext_java_lang_Enum() {}
func (Comment_Type) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ Comment_Type) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase { return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()} }
func (Comment_Type) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.parsing.parser.trees.Comment$Type", "Type" }
func (this_ Comment_Type) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_Comment_Type[this_-1]
}
func (this_ Comment_Type) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ Comment_Type) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ Comment_Type) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Comment_Type[this_-1].Name_
}
func (this_ Comment_Type) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_Comment_Type[this_-1].Name_
}
func (this_ Comment_Type) Equals(other_ any) bool {
	o_, ok_ := other_.(Comment_Type)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ Comment_Type) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_Comment_Type[this_-1].HashCode()
}
func (this_ Comment_Type) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Comment_Type[this_-1].CompareTo(other_)
}
func (this_ Comment_Type) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_Comment_Type[this_-1].Name_)
}
func init() {
	vals_ := []any{Comment_Type(1), Comment_Type(2), Comment_Type(3), Comment_Type(4), Comment_Type(5)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.parser.trees.Comment.Type", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.parser.trees.Comment$Type", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.Comment.Type", "com.google.javascript.jscomp.parsing.parser.trees.Comment$Type", func(v_ any) bool {
		return j2gort.InstanceOf[Comment_Type](v_)
	})
}

