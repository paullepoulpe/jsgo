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
// Port of com/google/javascript/jscomp/parsing/parser/trees/ParseTreeType.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java enum: com.google.javascript.jscomp.parsing.parser.trees.ParseTreeType
type ParseTreeType uint8

const (
	enum_ParseTreeType_PROGRAM ParseTreeType = 1
	enum_ParseTreeType_FUNCTION_DECLARATION ParseTreeType = 2
	enum_ParseTreeType_BLOCK ParseTreeType = 3
	enum_ParseTreeType_VARIABLE_STATEMENT ParseTreeType = 4
	enum_ParseTreeType_VARIABLE_DECLARATION ParseTreeType = 5
	enum_ParseTreeType_EMPTY_STATEMENT ParseTreeType = 6
	enum_ParseTreeType_EXPRESSION_STATEMENT ParseTreeType = 7
	enum_ParseTreeType_IF_STATEMENT ParseTreeType = 8
	enum_ParseTreeType_DO_WHILE_STATEMENT ParseTreeType = 9
	enum_ParseTreeType_WHILE_STATEMENT ParseTreeType = 10
	enum_ParseTreeType_FOR_IN_STATEMENT ParseTreeType = 11
	enum_ParseTreeType_FOR_STATEMENT ParseTreeType = 12
	enum_ParseTreeType_VARIABLE_DECLARATION_LIST ParseTreeType = 13
	enum_ParseTreeType_CONTINUE_STATEMENT ParseTreeType = 14
	enum_ParseTreeType_BREAK_STATEMENT ParseTreeType = 15
	enum_ParseTreeType_RETURN_STATEMENT ParseTreeType = 16
	enum_ParseTreeType_WITH_STATEMENT ParseTreeType = 17
	enum_ParseTreeType_CASE_CLAUSE ParseTreeType = 18
	enum_ParseTreeType_DEFAULT_CLAUSE ParseTreeType = 19
	enum_ParseTreeType_SWITCH_STATEMENT ParseTreeType = 20
	enum_ParseTreeType_LABELLED_STATEMENT ParseTreeType = 21
	enum_ParseTreeType_THROW_STATEMENT ParseTreeType = 22
	enum_ParseTreeType_CATCH ParseTreeType = 23
	enum_ParseTreeType_TRY_STATEMENT ParseTreeType = 24
	enum_ParseTreeType_DEBUGGER_STATEMENT ParseTreeType = 25
	enum_ParseTreeType_THIS_EXPRESSION ParseTreeType = 26
	enum_ParseTreeType_IDENTIFIER_EXPRESSION ParseTreeType = 27
	enum_ParseTreeType_LITERAL_EXPRESSION ParseTreeType = 28
	enum_ParseTreeType_ARRAY_LITERAL_EXPRESSION ParseTreeType = 29
	enum_ParseTreeType_OBJECT_LITERAL_EXPRESSION ParseTreeType = 30
	enum_ParseTreeType_GET_ACCESSOR ParseTreeType = 31
	enum_ParseTreeType_SET_ACCESSOR ParseTreeType = 32
	enum_ParseTreeType_PROPERTY_NAME_ASSIGNMENT ParseTreeType = 33
	enum_ParseTreeType_COMPUTED_PROPERTY_DEFINITION ParseTreeType = 34
	enum_ParseTreeType_COMPUTED_PROPERTY_GETTER ParseTreeType = 35
	enum_ParseTreeType_COMPUTED_PROPERTY_METHOD ParseTreeType = 36
	enum_ParseTreeType_COMPUTED_PROPERTY_SETTER ParseTreeType = 37
	enum_ParseTreeType_MISSING_PRIMARY_EXPRESSION ParseTreeType = 38
	enum_ParseTreeType_COMMA_EXPRESSION ParseTreeType = 39
	enum_ParseTreeType_BINARY_OPERATOR ParseTreeType = 40
	enum_ParseTreeType_CONDITIONAL_EXPRESSION ParseTreeType = 41
	enum_ParseTreeType_UNARY_EXPRESSION ParseTreeType = 42
	enum_ParseTreeType_UPDATE_EXPRESSION ParseTreeType = 43
	enum_ParseTreeType_MEMBER_EXPRESSION ParseTreeType = 44
	enum_ParseTreeType_OPT_CHAIN_MEMBER_EXPRESSION ParseTreeType = 45
	enum_ParseTreeType_NEW_EXPRESSION ParseTreeType = 46
	enum_ParseTreeType_ARGUMENT_LIST ParseTreeType = 47
	enum_ParseTreeType_CALL_EXPRESSION ParseTreeType = 48
	enum_ParseTreeType_OPT_CHAIN_CALL_EXPRESSION ParseTreeType = 49
	enum_ParseTreeType_CLASS_DECLARATION ParseTreeType = 50
	enum_ParseTreeType_MEMBER_LOOKUP_EXPRESSION ParseTreeType = 51
	enum_ParseTreeType_OPT_CHAIN_MEMBER_LOOKUP_EXPRESSION ParseTreeType = 52
	enum_ParseTreeType_PAREN_EXPRESSION ParseTreeType = 53
	enum_ParseTreeType_FINALLY ParseTreeType = 54
	enum_ParseTreeType_SUPER_EXPRESSION ParseTreeType = 55
	enum_ParseTreeType_ARRAY_PATTERN ParseTreeType = 56
	enum_ParseTreeType_OBJECT_PATTERN ParseTreeType = 57
	enum_ParseTreeType_FORMAL_PARAMETER_LIST ParseTreeType = 58
	enum_ParseTreeType_ITER_SPREAD ParseTreeType = 59
	enum_ParseTreeType_OBJECT_SPREAD ParseTreeType = 60
	enum_ParseTreeType_NULL ParseTreeType = 61
	enum_ParseTreeType_ITER_REST ParseTreeType = 62
	enum_ParseTreeType_OBJECT_REST ParseTreeType = 63
	enum_ParseTreeType_MODULE_IMPORT ParseTreeType = 64
	enum_ParseTreeType_EXPORT_DECLARATION ParseTreeType = 65
	enum_ParseTreeType_EXPORT_SPECIFIER ParseTreeType = 66
	enum_ParseTreeType_IMPORT_DECLARATION ParseTreeType = 67
	enum_ParseTreeType_IMPORT_SPECIFIER ParseTreeType = 68
	enum_ParseTreeType_FOR_OF_STATEMENT ParseTreeType = 69
	enum_ParseTreeType_FOR_AWAIT_OF_STATEMENT ParseTreeType = 70
	enum_ParseTreeType_YIELD_EXPRESSION ParseTreeType = 71
	enum_ParseTreeType_DEFAULT_PARAMETER ParseTreeType = 72
	enum_ParseTreeType_TEMPLATE_LITERAL_EXPRESSION ParseTreeType = 73
	enum_ParseTreeType_TEMPLATE_LITERAL_PORTION ParseTreeType = 74
	enum_ParseTreeType_TEMPLATE_SUBSTITUTION ParseTreeType = 75
	enum_ParseTreeType_OPTIONAL_PARAMETER ParseTreeType = 76
	enum_ParseTreeType_NEW_TARGET_EXPRESSION ParseTreeType = 77
	enum_ParseTreeType_AWAIT_EXPRESSION ParseTreeType = 78
	enum_ParseTreeType_DYNAMIC_IMPORT_EXPRESSION ParseTreeType = 79
	enum_ParseTreeType_IMPORT_META_EXPRESSION ParseTreeType = 80
	enum_ParseTreeType_COMPUTED_PROPERTY_FIELD ParseTreeType = 81
	enum_ParseTreeType_FIELD_DECLARATION ParseTreeType = 82
)

var enumBaseTable_ParseTreeType = [...]j2gort.EnumBase{
	{Name_: "PROGRAM", Ordinal_: 0},
	{Name_: "FUNCTION_DECLARATION", Ordinal_: 1},
	{Name_: "BLOCK", Ordinal_: 2},
	{Name_: "VARIABLE_STATEMENT", Ordinal_: 3},
	{Name_: "VARIABLE_DECLARATION", Ordinal_: 4},
	{Name_: "EMPTY_STATEMENT", Ordinal_: 5},
	{Name_: "EXPRESSION_STATEMENT", Ordinal_: 6},
	{Name_: "IF_STATEMENT", Ordinal_: 7},
	{Name_: "DO_WHILE_STATEMENT", Ordinal_: 8},
	{Name_: "WHILE_STATEMENT", Ordinal_: 9},
	{Name_: "FOR_IN_STATEMENT", Ordinal_: 10},
	{Name_: "FOR_STATEMENT", Ordinal_: 11},
	{Name_: "VARIABLE_DECLARATION_LIST", Ordinal_: 12},
	{Name_: "CONTINUE_STATEMENT", Ordinal_: 13},
	{Name_: "BREAK_STATEMENT", Ordinal_: 14},
	{Name_: "RETURN_STATEMENT", Ordinal_: 15},
	{Name_: "WITH_STATEMENT", Ordinal_: 16},
	{Name_: "CASE_CLAUSE", Ordinal_: 17},
	{Name_: "DEFAULT_CLAUSE", Ordinal_: 18},
	{Name_: "SWITCH_STATEMENT", Ordinal_: 19},
	{Name_: "LABELLED_STATEMENT", Ordinal_: 20},
	{Name_: "THROW_STATEMENT", Ordinal_: 21},
	{Name_: "CATCH", Ordinal_: 22},
	{Name_: "TRY_STATEMENT", Ordinal_: 23},
	{Name_: "DEBUGGER_STATEMENT", Ordinal_: 24},
	{Name_: "THIS_EXPRESSION", Ordinal_: 25},
	{Name_: "IDENTIFIER_EXPRESSION", Ordinal_: 26},
	{Name_: "LITERAL_EXPRESSION", Ordinal_: 27},
	{Name_: "ARRAY_LITERAL_EXPRESSION", Ordinal_: 28},
	{Name_: "OBJECT_LITERAL_EXPRESSION", Ordinal_: 29},
	{Name_: "GET_ACCESSOR", Ordinal_: 30},
	{Name_: "SET_ACCESSOR", Ordinal_: 31},
	{Name_: "PROPERTY_NAME_ASSIGNMENT", Ordinal_: 32},
	{Name_: "COMPUTED_PROPERTY_DEFINITION", Ordinal_: 33},
	{Name_: "COMPUTED_PROPERTY_GETTER", Ordinal_: 34},
	{Name_: "COMPUTED_PROPERTY_METHOD", Ordinal_: 35},
	{Name_: "COMPUTED_PROPERTY_SETTER", Ordinal_: 36},
	{Name_: "MISSING_PRIMARY_EXPRESSION", Ordinal_: 37},
	{Name_: "COMMA_EXPRESSION", Ordinal_: 38},
	{Name_: "BINARY_OPERATOR", Ordinal_: 39},
	{Name_: "CONDITIONAL_EXPRESSION", Ordinal_: 40},
	{Name_: "UNARY_EXPRESSION", Ordinal_: 41},
	{Name_: "UPDATE_EXPRESSION", Ordinal_: 42},
	{Name_: "MEMBER_EXPRESSION", Ordinal_: 43},
	{Name_: "OPT_CHAIN_MEMBER_EXPRESSION", Ordinal_: 44},
	{Name_: "NEW_EXPRESSION", Ordinal_: 45},
	{Name_: "ARGUMENT_LIST", Ordinal_: 46},
	{Name_: "CALL_EXPRESSION", Ordinal_: 47},
	{Name_: "OPT_CHAIN_CALL_EXPRESSION", Ordinal_: 48},
	{Name_: "CLASS_DECLARATION", Ordinal_: 49},
	{Name_: "MEMBER_LOOKUP_EXPRESSION", Ordinal_: 50},
	{Name_: "OPT_CHAIN_MEMBER_LOOKUP_EXPRESSION", Ordinal_: 51},
	{Name_: "PAREN_EXPRESSION", Ordinal_: 52},
	{Name_: "FINALLY", Ordinal_: 53},
	{Name_: "SUPER_EXPRESSION", Ordinal_: 54},
	{Name_: "ARRAY_PATTERN", Ordinal_: 55},
	{Name_: "OBJECT_PATTERN", Ordinal_: 56},
	{Name_: "FORMAL_PARAMETER_LIST", Ordinal_: 57},
	{Name_: "ITER_SPREAD", Ordinal_: 58},
	{Name_: "OBJECT_SPREAD", Ordinal_: 59},
	{Name_: "NULL", Ordinal_: 60},
	{Name_: "ITER_REST", Ordinal_: 61},
	{Name_: "OBJECT_REST", Ordinal_: 62},
	{Name_: "MODULE_IMPORT", Ordinal_: 63},
	{Name_: "EXPORT_DECLARATION", Ordinal_: 64},
	{Name_: "EXPORT_SPECIFIER", Ordinal_: 65},
	{Name_: "IMPORT_DECLARATION", Ordinal_: 66},
	{Name_: "IMPORT_SPECIFIER", Ordinal_: 67},
	{Name_: "FOR_OF_STATEMENT", Ordinal_: 68},
	{Name_: "FOR_AWAIT_OF_STATEMENT", Ordinal_: 69},
	{Name_: "YIELD_EXPRESSION", Ordinal_: 70},
	{Name_: "DEFAULT_PARAMETER", Ordinal_: 71},
	{Name_: "TEMPLATE_LITERAL_EXPRESSION", Ordinal_: 72},
	{Name_: "TEMPLATE_LITERAL_PORTION", Ordinal_: 73},
	{Name_: "TEMPLATE_SUBSTITUTION", Ordinal_: 74},
	{Name_: "OPTIONAL_PARAMETER", Ordinal_: 75},
	{Name_: "NEW_TARGET_EXPRESSION", Ordinal_: 76},
	{Name_: "AWAIT_EXPRESSION", Ordinal_: 77},
	{Name_: "DYNAMIC_IMPORT_EXPRESSION", Ordinal_: 78},
	{Name_: "IMPORT_META_EXPRESSION", Ordinal_: 79},
	{Name_: "COMPUTED_PROPERTY_FIELD", Ordinal_: 80},
	{Name_: "FIELD_DECLARATION", Ordinal_: 81},
}
var enumValues_ParseTreeType = []ParseTreeType{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82}

func StaticGet_ParseTreeType_PROGRAM() ParseTreeType {
	return enum_ParseTreeType_PROGRAM
}
func StaticGet_ParseTreeType_FUNCTION_DECLARATION() ParseTreeType {
	return enum_ParseTreeType_FUNCTION_DECLARATION
}
func StaticGet_ParseTreeType_BLOCK() ParseTreeType {
	return enum_ParseTreeType_BLOCK
}
func StaticGet_ParseTreeType_VARIABLE_STATEMENT() ParseTreeType {
	return enum_ParseTreeType_VARIABLE_STATEMENT
}
func StaticGet_ParseTreeType_VARIABLE_DECLARATION() ParseTreeType {
	return enum_ParseTreeType_VARIABLE_DECLARATION
}
func StaticGet_ParseTreeType_EMPTY_STATEMENT() ParseTreeType {
	return enum_ParseTreeType_EMPTY_STATEMENT
}
func StaticGet_ParseTreeType_EXPRESSION_STATEMENT() ParseTreeType {
	return enum_ParseTreeType_EXPRESSION_STATEMENT
}
func StaticGet_ParseTreeType_IF_STATEMENT() ParseTreeType {
	return enum_ParseTreeType_IF_STATEMENT
}
func StaticGet_ParseTreeType_DO_WHILE_STATEMENT() ParseTreeType {
	return enum_ParseTreeType_DO_WHILE_STATEMENT
}
func StaticGet_ParseTreeType_WHILE_STATEMENT() ParseTreeType {
	return enum_ParseTreeType_WHILE_STATEMENT
}
func StaticGet_ParseTreeType_FOR_IN_STATEMENT() ParseTreeType {
	return enum_ParseTreeType_FOR_IN_STATEMENT
}
func StaticGet_ParseTreeType_FOR_STATEMENT() ParseTreeType {
	return enum_ParseTreeType_FOR_STATEMENT
}
func StaticGet_ParseTreeType_VARIABLE_DECLARATION_LIST() ParseTreeType {
	return enum_ParseTreeType_VARIABLE_DECLARATION_LIST
}
func StaticGet_ParseTreeType_CONTINUE_STATEMENT() ParseTreeType {
	return enum_ParseTreeType_CONTINUE_STATEMENT
}
func StaticGet_ParseTreeType_BREAK_STATEMENT() ParseTreeType {
	return enum_ParseTreeType_BREAK_STATEMENT
}
func StaticGet_ParseTreeType_RETURN_STATEMENT() ParseTreeType {
	return enum_ParseTreeType_RETURN_STATEMENT
}
func StaticGet_ParseTreeType_WITH_STATEMENT() ParseTreeType {
	return enum_ParseTreeType_WITH_STATEMENT
}
func StaticGet_ParseTreeType_CASE_CLAUSE() ParseTreeType {
	return enum_ParseTreeType_CASE_CLAUSE
}
func StaticGet_ParseTreeType_DEFAULT_CLAUSE() ParseTreeType {
	return enum_ParseTreeType_DEFAULT_CLAUSE
}
func StaticGet_ParseTreeType_SWITCH_STATEMENT() ParseTreeType {
	return enum_ParseTreeType_SWITCH_STATEMENT
}
func StaticGet_ParseTreeType_LABELLED_STATEMENT() ParseTreeType {
	return enum_ParseTreeType_LABELLED_STATEMENT
}
func StaticGet_ParseTreeType_THROW_STATEMENT() ParseTreeType {
	return enum_ParseTreeType_THROW_STATEMENT
}
func StaticGet_ParseTreeType_CATCH() ParseTreeType {
	return enum_ParseTreeType_CATCH
}
func StaticGet_ParseTreeType_TRY_STATEMENT() ParseTreeType {
	return enum_ParseTreeType_TRY_STATEMENT
}
func StaticGet_ParseTreeType_DEBUGGER_STATEMENT() ParseTreeType {
	return enum_ParseTreeType_DEBUGGER_STATEMENT
}
func StaticGet_ParseTreeType_THIS_EXPRESSION() ParseTreeType {
	return enum_ParseTreeType_THIS_EXPRESSION
}
func StaticGet_ParseTreeType_IDENTIFIER_EXPRESSION() ParseTreeType {
	return enum_ParseTreeType_IDENTIFIER_EXPRESSION
}
func StaticGet_ParseTreeType_LITERAL_EXPRESSION() ParseTreeType {
	return enum_ParseTreeType_LITERAL_EXPRESSION
}
func StaticGet_ParseTreeType_ARRAY_LITERAL_EXPRESSION() ParseTreeType {
	return enum_ParseTreeType_ARRAY_LITERAL_EXPRESSION
}
func StaticGet_ParseTreeType_OBJECT_LITERAL_EXPRESSION() ParseTreeType {
	return enum_ParseTreeType_OBJECT_LITERAL_EXPRESSION
}
func StaticGet_ParseTreeType_GET_ACCESSOR() ParseTreeType {
	return enum_ParseTreeType_GET_ACCESSOR
}
func StaticGet_ParseTreeType_SET_ACCESSOR() ParseTreeType {
	return enum_ParseTreeType_SET_ACCESSOR
}
func StaticGet_ParseTreeType_PROPERTY_NAME_ASSIGNMENT() ParseTreeType {
	return enum_ParseTreeType_PROPERTY_NAME_ASSIGNMENT
}
func StaticGet_ParseTreeType_COMPUTED_PROPERTY_DEFINITION() ParseTreeType {
	return enum_ParseTreeType_COMPUTED_PROPERTY_DEFINITION
}
func StaticGet_ParseTreeType_COMPUTED_PROPERTY_GETTER() ParseTreeType {
	return enum_ParseTreeType_COMPUTED_PROPERTY_GETTER
}
func StaticGet_ParseTreeType_COMPUTED_PROPERTY_METHOD() ParseTreeType {
	return enum_ParseTreeType_COMPUTED_PROPERTY_METHOD
}
func StaticGet_ParseTreeType_COMPUTED_PROPERTY_SETTER() ParseTreeType {
	return enum_ParseTreeType_COMPUTED_PROPERTY_SETTER
}
func StaticGet_ParseTreeType_MISSING_PRIMARY_EXPRESSION() ParseTreeType {
	return enum_ParseTreeType_MISSING_PRIMARY_EXPRESSION
}
func StaticGet_ParseTreeType_COMMA_EXPRESSION() ParseTreeType {
	return enum_ParseTreeType_COMMA_EXPRESSION
}
func StaticGet_ParseTreeType_BINARY_OPERATOR() ParseTreeType {
	return enum_ParseTreeType_BINARY_OPERATOR
}
func StaticGet_ParseTreeType_CONDITIONAL_EXPRESSION() ParseTreeType {
	return enum_ParseTreeType_CONDITIONAL_EXPRESSION
}
func StaticGet_ParseTreeType_UNARY_EXPRESSION() ParseTreeType {
	return enum_ParseTreeType_UNARY_EXPRESSION
}
func StaticGet_ParseTreeType_UPDATE_EXPRESSION() ParseTreeType {
	return enum_ParseTreeType_UPDATE_EXPRESSION
}
func StaticGet_ParseTreeType_MEMBER_EXPRESSION() ParseTreeType {
	return enum_ParseTreeType_MEMBER_EXPRESSION
}
func StaticGet_ParseTreeType_OPT_CHAIN_MEMBER_EXPRESSION() ParseTreeType {
	return enum_ParseTreeType_OPT_CHAIN_MEMBER_EXPRESSION
}
func StaticGet_ParseTreeType_NEW_EXPRESSION() ParseTreeType {
	return enum_ParseTreeType_NEW_EXPRESSION
}
func StaticGet_ParseTreeType_ARGUMENT_LIST() ParseTreeType {
	return enum_ParseTreeType_ARGUMENT_LIST
}
func StaticGet_ParseTreeType_CALL_EXPRESSION() ParseTreeType {
	return enum_ParseTreeType_CALL_EXPRESSION
}
func StaticGet_ParseTreeType_OPT_CHAIN_CALL_EXPRESSION() ParseTreeType {
	return enum_ParseTreeType_OPT_CHAIN_CALL_EXPRESSION
}
func StaticGet_ParseTreeType_CLASS_DECLARATION() ParseTreeType {
	return enum_ParseTreeType_CLASS_DECLARATION
}
func StaticGet_ParseTreeType_MEMBER_LOOKUP_EXPRESSION() ParseTreeType {
	return enum_ParseTreeType_MEMBER_LOOKUP_EXPRESSION
}
func StaticGet_ParseTreeType_OPT_CHAIN_MEMBER_LOOKUP_EXPRESSION() ParseTreeType {
	return enum_ParseTreeType_OPT_CHAIN_MEMBER_LOOKUP_EXPRESSION
}
func StaticGet_ParseTreeType_PAREN_EXPRESSION() ParseTreeType {
	return enum_ParseTreeType_PAREN_EXPRESSION
}
func StaticGet_ParseTreeType_FINALLY() ParseTreeType {
	return enum_ParseTreeType_FINALLY
}
func StaticGet_ParseTreeType_SUPER_EXPRESSION() ParseTreeType {
	return enum_ParseTreeType_SUPER_EXPRESSION
}
func StaticGet_ParseTreeType_ARRAY_PATTERN() ParseTreeType {
	return enum_ParseTreeType_ARRAY_PATTERN
}
func StaticGet_ParseTreeType_OBJECT_PATTERN() ParseTreeType {
	return enum_ParseTreeType_OBJECT_PATTERN
}
func StaticGet_ParseTreeType_FORMAL_PARAMETER_LIST() ParseTreeType {
	return enum_ParseTreeType_FORMAL_PARAMETER_LIST
}
func StaticGet_ParseTreeType_ITER_SPREAD() ParseTreeType {
	return enum_ParseTreeType_ITER_SPREAD
}
func StaticGet_ParseTreeType_OBJECT_SPREAD() ParseTreeType {
	return enum_ParseTreeType_OBJECT_SPREAD
}
func StaticGet_ParseTreeType_NULL() ParseTreeType {
	return enum_ParseTreeType_NULL
}
func StaticGet_ParseTreeType_ITER_REST() ParseTreeType {
	return enum_ParseTreeType_ITER_REST
}
func StaticGet_ParseTreeType_OBJECT_REST() ParseTreeType {
	return enum_ParseTreeType_OBJECT_REST
}
func StaticGet_ParseTreeType_MODULE_IMPORT() ParseTreeType {
	return enum_ParseTreeType_MODULE_IMPORT
}
func StaticGet_ParseTreeType_EXPORT_DECLARATION() ParseTreeType {
	return enum_ParseTreeType_EXPORT_DECLARATION
}
func StaticGet_ParseTreeType_EXPORT_SPECIFIER() ParseTreeType {
	return enum_ParseTreeType_EXPORT_SPECIFIER
}
func StaticGet_ParseTreeType_IMPORT_DECLARATION() ParseTreeType {
	return enum_ParseTreeType_IMPORT_DECLARATION
}
func StaticGet_ParseTreeType_IMPORT_SPECIFIER() ParseTreeType {
	return enum_ParseTreeType_IMPORT_SPECIFIER
}
func StaticGet_ParseTreeType_FOR_OF_STATEMENT() ParseTreeType {
	return enum_ParseTreeType_FOR_OF_STATEMENT
}
func StaticGet_ParseTreeType_FOR_AWAIT_OF_STATEMENT() ParseTreeType {
	return enum_ParseTreeType_FOR_AWAIT_OF_STATEMENT
}
func StaticGet_ParseTreeType_YIELD_EXPRESSION() ParseTreeType {
	return enum_ParseTreeType_YIELD_EXPRESSION
}
func StaticGet_ParseTreeType_DEFAULT_PARAMETER() ParseTreeType {
	return enum_ParseTreeType_DEFAULT_PARAMETER
}
func StaticGet_ParseTreeType_TEMPLATE_LITERAL_EXPRESSION() ParseTreeType {
	return enum_ParseTreeType_TEMPLATE_LITERAL_EXPRESSION
}
func StaticGet_ParseTreeType_TEMPLATE_LITERAL_PORTION() ParseTreeType {
	return enum_ParseTreeType_TEMPLATE_LITERAL_PORTION
}
func StaticGet_ParseTreeType_TEMPLATE_SUBSTITUTION() ParseTreeType {
	return enum_ParseTreeType_TEMPLATE_SUBSTITUTION
}
func StaticGet_ParseTreeType_OPTIONAL_PARAMETER() ParseTreeType {
	return enum_ParseTreeType_OPTIONAL_PARAMETER
}
func StaticGet_ParseTreeType_NEW_TARGET_EXPRESSION() ParseTreeType {
	return enum_ParseTreeType_NEW_TARGET_EXPRESSION
}
func StaticGet_ParseTreeType_AWAIT_EXPRESSION() ParseTreeType {
	return enum_ParseTreeType_AWAIT_EXPRESSION
}
func StaticGet_ParseTreeType_DYNAMIC_IMPORT_EXPRESSION() ParseTreeType {
	return enum_ParseTreeType_DYNAMIC_IMPORT_EXPRESSION
}
func StaticGet_ParseTreeType_IMPORT_META_EXPRESSION() ParseTreeType {
	return enum_ParseTreeType_IMPORT_META_EXPRESSION
}
func StaticGet_ParseTreeType_COMPUTED_PROPERTY_FIELD() ParseTreeType {
	return enum_ParseTreeType_COMPUTED_PROPERTY_FIELD
}
func StaticGet_ParseTreeType_FIELD_DECLARATION() ParseTreeType {
	return enum_ParseTreeType_FIELD_DECLARATION
}

func ParseTreeType_Values() []ParseTreeType {
	return j2gort.CloneSlice(enumValues_ParseTreeType)
}

func ParseTreeType_ValueOf(name string) ParseTreeType {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_ParseTreeType {
		if enumBaseTable_ParseTreeType[i_].Name_ == name {
			return ParseTreeType(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.parsing.parser.trees.ParseTreeType." + name))
}

func (this_ ParseTreeType) IsJavaNil_() bool { return this_ == 0 }
func (ParseTreeType) HasIdentityEquals_() {}
func (ParseTreeType) Implements_j2gort_Ext_java_lang_Enum() {}
func (ParseTreeType) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ ParseTreeType) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase { return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()} }
func (ParseTreeType) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.parsing.parser.trees.ParseTreeType", "ParseTreeType" }
func (this_ ParseTreeType) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_ParseTreeType[this_-1]
}
func (this_ ParseTreeType) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ ParseTreeType) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ ParseTreeType) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ParseTreeType[this_-1].Name_
}
func (this_ ParseTreeType) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_ParseTreeType[this_-1].Name_
}
func (this_ ParseTreeType) Equals(other_ any) bool {
	o_, ok_ := other_.(ParseTreeType)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ ParseTreeType) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_ParseTreeType[this_-1].HashCode()
}
func (this_ ParseTreeType) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ParseTreeType[this_-1].CompareTo(other_)
}
func (this_ ParseTreeType) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_ParseTreeType[this_-1].Name_)
}
func init() {
	vals_ := []any{ParseTreeType(1), ParseTreeType(2), ParseTreeType(3), ParseTreeType(4), ParseTreeType(5), ParseTreeType(6), ParseTreeType(7), ParseTreeType(8), ParseTreeType(9), ParseTreeType(10), ParseTreeType(11), ParseTreeType(12), ParseTreeType(13), ParseTreeType(14), ParseTreeType(15), ParseTreeType(16), ParseTreeType(17), ParseTreeType(18), ParseTreeType(19), ParseTreeType(20), ParseTreeType(21), ParseTreeType(22), ParseTreeType(23), ParseTreeType(24), ParseTreeType(25), ParseTreeType(26), ParseTreeType(27), ParseTreeType(28), ParseTreeType(29), ParseTreeType(30), ParseTreeType(31), ParseTreeType(32), ParseTreeType(33), ParseTreeType(34), ParseTreeType(35), ParseTreeType(36), ParseTreeType(37), ParseTreeType(38), ParseTreeType(39), ParseTreeType(40), ParseTreeType(41), ParseTreeType(42), ParseTreeType(43), ParseTreeType(44), ParseTreeType(45), ParseTreeType(46), ParseTreeType(47), ParseTreeType(48), ParseTreeType(49), ParseTreeType(50), ParseTreeType(51), ParseTreeType(52), ParseTreeType(53), ParseTreeType(54), ParseTreeType(55), ParseTreeType(56), ParseTreeType(57), ParseTreeType(58), ParseTreeType(59), ParseTreeType(60), ParseTreeType(61), ParseTreeType(62), ParseTreeType(63), ParseTreeType(64), ParseTreeType(65), ParseTreeType(66), ParseTreeType(67), ParseTreeType(68), ParseTreeType(69), ParseTreeType(70), ParseTreeType(71), ParseTreeType(72), ParseTreeType(73), ParseTreeType(74), ParseTreeType(75), ParseTreeType(76), ParseTreeType(77), ParseTreeType(78), ParseTreeType(79), ParseTreeType(80), ParseTreeType(81), ParseTreeType(82)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.parser.trees.ParseTreeType", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.ParseTreeType", "com.google.javascript.jscomp.parsing.parser.trees.ParseTreeType", func(v_ any) bool {
		return j2gort.InstanceOf[ParseTreeType](v_)
	})
}

