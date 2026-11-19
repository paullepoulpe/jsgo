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
 *   Roger Lawrence
 *   Mike McCabe
 *   Igor Bukanov
 *   Milen Nankov
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */
// Port of com/google/javascript/rhino/Token.java (generated mechanically by j2go).

package rhino

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java enum: com.google.javascript.rhino.Token
type Token uint8

const (
	enum_Token_RETURN              Token = 1
	enum_Token_BITOR               Token = 2
	enum_Token_BITXOR              Token = 3
	enum_Token_BITAND              Token = 4
	enum_Token_EQ                  Token = 5
	enum_Token_NE                  Token = 6
	enum_Token_LT                  Token = 7
	enum_Token_LE                  Token = 8
	enum_Token_GT                  Token = 9
	enum_Token_GE                  Token = 10
	enum_Token_LSH                 Token = 11
	enum_Token_RSH                 Token = 12
	enum_Token_URSH                Token = 13
	enum_Token_ADD                 Token = 14
	enum_Token_SUB                 Token = 15
	enum_Token_MUL                 Token = 16
	enum_Token_DIV                 Token = 17
	enum_Token_MOD                 Token = 18
	enum_Token_EXPONENT            Token = 19
	enum_Token_NOT                 Token = 20
	enum_Token_BITNOT              Token = 21
	enum_Token_POS                 Token = 22
	enum_Token_NEG                 Token = 23
	enum_Token_NEW                 Token = 24
	enum_Token_DELPROP             Token = 25
	enum_Token_TYPEOF              Token = 26
	enum_Token_GETPROP             Token = 27
	enum_Token_GETELEM             Token = 28
	enum_Token_CALL                Token = 29
	enum_Token_OPTCHAIN_GETPROP    Token = 30
	enum_Token_OPTCHAIN_GETELEM    Token = 31
	enum_Token_OPTCHAIN_CALL       Token = 32
	enum_Token_NAME                Token = 33
	enum_Token_NUMBER              Token = 34
	enum_Token_BIGINT              Token = 35
	enum_Token_STRINGLIT           Token = 36
	enum_Token_NULL                Token = 37
	enum_Token_THIS                Token = 38
	enum_Token_FALSE               Token = 39
	enum_Token_TRUE                Token = 40
	enum_Token_SHEQ                Token = 41
	enum_Token_SHNE                Token = 42
	enum_Token_REGEXP              Token = 43
	enum_Token_THROW               Token = 44
	enum_Token_IN                  Token = 45
	enum_Token_INSTANCEOF          Token = 46
	enum_Token_ARRAYLIT            Token = 47
	enum_Token_OBJECTLIT           Token = 48
	enum_Token_TRY                 Token = 49
	enum_Token_PARAM_LIST          Token = 50
	enum_Token_COMMA               Token = 51
	enum_Token_ASSIGN              Token = 52
	enum_Token_ASSIGN_BITOR        Token = 53
	enum_Token_ASSIGN_BITXOR       Token = 54
	enum_Token_ASSIGN_BITAND       Token = 55
	enum_Token_ASSIGN_LSH          Token = 56
	enum_Token_ASSIGN_RSH          Token = 57
	enum_Token_ASSIGN_URSH         Token = 58
	enum_Token_ASSIGN_ADD          Token = 59
	enum_Token_ASSIGN_SUB          Token = 60
	enum_Token_ASSIGN_MUL          Token = 61
	enum_Token_ASSIGN_DIV          Token = 62
	enum_Token_ASSIGN_MOD          Token = 63
	enum_Token_ASSIGN_EXPONENT     Token = 64
	enum_Token_ASSIGN_OR           Token = 65
	enum_Token_ASSIGN_AND          Token = 66
	enum_Token_ASSIGN_COALESCE     Token = 67
	enum_Token_HOOK                Token = 68
	enum_Token_OR                  Token = 69
	enum_Token_AND                 Token = 70
	enum_Token_COALESCE            Token = 71
	enum_Token_INC                 Token = 72
	enum_Token_DEC                 Token = 73
	enum_Token_FUNCTION            Token = 74
	enum_Token_IF                  Token = 75
	enum_Token_SWITCH              Token = 76
	enum_Token_CASE                Token = 77
	enum_Token_DEFAULT_CASE        Token = 78
	enum_Token_WHILE               Token = 79
	enum_Token_DO                  Token = 80
	enum_Token_FOR                 Token = 81
	enum_Token_FOR_IN              Token = 82
	enum_Token_BREAK               Token = 83
	enum_Token_CONTINUE            Token = 84
	enum_Token_VAR                 Token = 85
	enum_Token_WITH                Token = 86
	enum_Token_CATCH               Token = 87
	enum_Token_VOID                Token = 88
	enum_Token_EMPTY               Token = 89
	enum_Token_ROOT                Token = 90
	enum_Token_BLOCK               Token = 91
	enum_Token_SWITCH_BODY         Token = 92
	enum_Token_LABEL               Token = 93
	enum_Token_EXPR_RESULT         Token = 94
	enum_Token_SCRIPT              Token = 95
	enum_Token_GETTER_DEF          Token = 96
	enum_Token_SETTER_DEF          Token = 97
	enum_Token_CONST               Token = 98
	enum_Token_DEBUGGER            Token = 99
	enum_Token_LABEL_NAME          Token = 100
	enum_Token_STRING_KEY          Token = 101
	enum_Token_CAST                Token = 102
	enum_Token_ARRAY_PATTERN       Token = 103
	enum_Token_OBJECT_PATTERN      Token = 104
	enum_Token_DESTRUCTURING_LHS   Token = 105
	enum_Token_CLASS               Token = 106
	enum_Token_CLASS_MEMBERS       Token = 107
	enum_Token_MEMBER_FUNCTION_DEF Token = 108
	enum_Token_MEMBER_FIELD_DEF    Token = 109
	enum_Token_COMPUTED_FIELD_DEF  Token = 110
	enum_Token_SUPER               Token = 111
	enum_Token_LET                 Token = 112
	enum_Token_FOR_OF              Token = 113
	enum_Token_FOR_AWAIT_OF        Token = 114
	enum_Token_YIELD               Token = 115
	enum_Token_AWAIT               Token = 116
	enum_Token_IMPORT              Token = 117
	enum_Token_IMPORT_SPECS        Token = 118
	enum_Token_IMPORT_SPEC         Token = 119
	enum_Token_IMPORT_STAR         Token = 120
	enum_Token_EXPORT              Token = 121
	enum_Token_EXPORT_SPECS        Token = 122
	enum_Token_EXPORT_SPEC         Token = 123
	enum_Token_MODULE_BODY         Token = 124
	enum_Token_DYNAMIC_IMPORT      Token = 125
	enum_Token_ITER_REST           Token = 126
	enum_Token_OBJECT_REST         Token = 127
	enum_Token_ITER_SPREAD         Token = 128
	enum_Token_OBJECT_SPREAD       Token = 129
	enum_Token_COMPUTED_PROP       Token = 130
	enum_Token_TAGGED_TEMPLATELIT  Token = 131
	enum_Token_TEMPLATELIT         Token = 132
	enum_Token_TEMPLATELIT_SUB     Token = 133
	enum_Token_TEMPLATELIT_STRING  Token = 134
	enum_Token_DEFAULT_VALUE       Token = 135
	enum_Token_NEW_TARGET          Token = 136
	enum_Token_IMPORT_META         Token = 137
	enum_Token_STRING_TYPE         Token = 138
	enum_Token_BOOLEAN_TYPE        Token = 139
	enum_Token_NUMBER_TYPE         Token = 140
	enum_Token_FUNCTION_TYPE       Token = 141
	enum_Token_PARAMETERIZED_TYPE  Token = 142
	enum_Token_UNION_TYPE          Token = 143
	enum_Token_ANY_TYPE            Token = 144
	enum_Token_NULLABLE_TYPE       Token = 145
	enum_Token_VOID_TYPE           Token = 146
	enum_Token_REST_PARAMETER_TYPE Token = 147
	enum_Token_NAMED_TYPE          Token = 148
	enum_Token_OPTIONAL_PARAMETER  Token = 149
	enum_Token_RECORD_TYPE         Token = 150
	enum_Token_UNDEFINED_TYPE      Token = 151
	enum_Token_ARRAY_TYPE          Token = 152
	enum_Token_GENERIC_TYPE        Token = 153
	enum_Token_GENERIC_TYPE_LIST   Token = 154
	enum_Token_ANNOTATION          Token = 155
	enum_Token_PIPE                Token = 156
	enum_Token_STAR                Token = 157
	enum_Token_EOC                 Token = 158
	enum_Token_QMARK               Token = 159
	enum_Token_BANG                Token = 160
	enum_Token_EQUALS              Token = 161
	enum_Token_LB                  Token = 162
	enum_Token_LC                  Token = 163
	enum_Token_COLON               Token = 164
	enum_Token_INTERFACE           Token = 165
	enum_Token_INTERFACE_EXTENDS   Token = 166
	enum_Token_INTERFACE_MEMBERS   Token = 167
	enum_Token_ENUM                Token = 168
	enum_Token_ENUM_MEMBERS        Token = 169
	enum_Token_IMPLEMENTS          Token = 170
	enum_Token_TYPE_ALIAS          Token = 171
	enum_Token_DECLARE             Token = 172
	enum_Token_MEMBER_VARIABLE_DEF Token = 173
	enum_Token_INDEX_SIGNATURE     Token = 174
	enum_Token_CALL_SIGNATURE      Token = 175
	enum_Token_NAMESPACE           Token = 176
	enum_Token_NAMESPACE_ELEMENTS  Token = 177
	enum_Token_PLACEHOLDER1        Token = 178
	enum_Token_PLACEHOLDER2        Token = 179
	enum_Token_PLACEHOLDER3        Token = 180
)

var enumBaseTable_Token = [...]j2gort.EnumBase{
	{Name_: "RETURN", Ordinal_: 0},
	{Name_: "BITOR", Ordinal_: 1},
	{Name_: "BITXOR", Ordinal_: 2},
	{Name_: "BITAND", Ordinal_: 3},
	{Name_: "EQ", Ordinal_: 4},
	{Name_: "NE", Ordinal_: 5},
	{Name_: "LT", Ordinal_: 6},
	{Name_: "LE", Ordinal_: 7},
	{Name_: "GT", Ordinal_: 8},
	{Name_: "GE", Ordinal_: 9},
	{Name_: "LSH", Ordinal_: 10},
	{Name_: "RSH", Ordinal_: 11},
	{Name_: "URSH", Ordinal_: 12},
	{Name_: "ADD", Ordinal_: 13},
	{Name_: "SUB", Ordinal_: 14},
	{Name_: "MUL", Ordinal_: 15},
	{Name_: "DIV", Ordinal_: 16},
	{Name_: "MOD", Ordinal_: 17},
	{Name_: "EXPONENT", Ordinal_: 18},
	{Name_: "NOT", Ordinal_: 19},
	{Name_: "BITNOT", Ordinal_: 20},
	{Name_: "POS", Ordinal_: 21},
	{Name_: "NEG", Ordinal_: 22},
	{Name_: "NEW", Ordinal_: 23},
	{Name_: "DELPROP", Ordinal_: 24},
	{Name_: "TYPEOF", Ordinal_: 25},
	{Name_: "GETPROP", Ordinal_: 26},
	{Name_: "GETELEM", Ordinal_: 27},
	{Name_: "CALL", Ordinal_: 28},
	{Name_: "OPTCHAIN_GETPROP", Ordinal_: 29},
	{Name_: "OPTCHAIN_GETELEM", Ordinal_: 30},
	{Name_: "OPTCHAIN_CALL", Ordinal_: 31},
	{Name_: "NAME", Ordinal_: 32},
	{Name_: "NUMBER", Ordinal_: 33},
	{Name_: "BIGINT", Ordinal_: 34},
	{Name_: "STRINGLIT", Ordinal_: 35},
	{Name_: "NULL", Ordinal_: 36},
	{Name_: "THIS", Ordinal_: 37},
	{Name_: "FALSE", Ordinal_: 38},
	{Name_: "TRUE", Ordinal_: 39},
	{Name_: "SHEQ", Ordinal_: 40},
	{Name_: "SHNE", Ordinal_: 41},
	{Name_: "REGEXP", Ordinal_: 42},
	{Name_: "THROW", Ordinal_: 43},
	{Name_: "IN", Ordinal_: 44},
	{Name_: "INSTANCEOF", Ordinal_: 45},
	{Name_: "ARRAYLIT", Ordinal_: 46},
	{Name_: "OBJECTLIT", Ordinal_: 47},
	{Name_: "TRY", Ordinal_: 48},
	{Name_: "PARAM_LIST", Ordinal_: 49},
	{Name_: "COMMA", Ordinal_: 50},
	{Name_: "ASSIGN", Ordinal_: 51},
	{Name_: "ASSIGN_BITOR", Ordinal_: 52},
	{Name_: "ASSIGN_BITXOR", Ordinal_: 53},
	{Name_: "ASSIGN_BITAND", Ordinal_: 54},
	{Name_: "ASSIGN_LSH", Ordinal_: 55},
	{Name_: "ASSIGN_RSH", Ordinal_: 56},
	{Name_: "ASSIGN_URSH", Ordinal_: 57},
	{Name_: "ASSIGN_ADD", Ordinal_: 58},
	{Name_: "ASSIGN_SUB", Ordinal_: 59},
	{Name_: "ASSIGN_MUL", Ordinal_: 60},
	{Name_: "ASSIGN_DIV", Ordinal_: 61},
	{Name_: "ASSIGN_MOD", Ordinal_: 62},
	{Name_: "ASSIGN_EXPONENT", Ordinal_: 63},
	{Name_: "ASSIGN_OR", Ordinal_: 64},
	{Name_: "ASSIGN_AND", Ordinal_: 65},
	{Name_: "ASSIGN_COALESCE", Ordinal_: 66},
	{Name_: "HOOK", Ordinal_: 67},
	{Name_: "OR", Ordinal_: 68},
	{Name_: "AND", Ordinal_: 69},
	{Name_: "COALESCE", Ordinal_: 70},
	{Name_: "INC", Ordinal_: 71},
	{Name_: "DEC", Ordinal_: 72},
	{Name_: "FUNCTION", Ordinal_: 73},
	{Name_: "IF", Ordinal_: 74},
	{Name_: "SWITCH", Ordinal_: 75},
	{Name_: "CASE", Ordinal_: 76},
	{Name_: "DEFAULT_CASE", Ordinal_: 77},
	{Name_: "WHILE", Ordinal_: 78},
	{Name_: "DO", Ordinal_: 79},
	{Name_: "FOR", Ordinal_: 80},
	{Name_: "FOR_IN", Ordinal_: 81},
	{Name_: "BREAK", Ordinal_: 82},
	{Name_: "CONTINUE", Ordinal_: 83},
	{Name_: "VAR", Ordinal_: 84},
	{Name_: "WITH", Ordinal_: 85},
	{Name_: "CATCH", Ordinal_: 86},
	{Name_: "VOID", Ordinal_: 87},
	{Name_: "EMPTY", Ordinal_: 88},
	{Name_: "ROOT", Ordinal_: 89},
	{Name_: "BLOCK", Ordinal_: 90},
	{Name_: "SWITCH_BODY", Ordinal_: 91},
	{Name_: "LABEL", Ordinal_: 92},
	{Name_: "EXPR_RESULT", Ordinal_: 93},
	{Name_: "SCRIPT", Ordinal_: 94},
	{Name_: "GETTER_DEF", Ordinal_: 95},
	{Name_: "SETTER_DEF", Ordinal_: 96},
	{Name_: "CONST", Ordinal_: 97},
	{Name_: "DEBUGGER", Ordinal_: 98},
	{Name_: "LABEL_NAME", Ordinal_: 99},
	{Name_: "STRING_KEY", Ordinal_: 100},
	{Name_: "CAST", Ordinal_: 101},
	{Name_: "ARRAY_PATTERN", Ordinal_: 102},
	{Name_: "OBJECT_PATTERN", Ordinal_: 103},
	{Name_: "DESTRUCTURING_LHS", Ordinal_: 104},
	{Name_: "CLASS", Ordinal_: 105},
	{Name_: "CLASS_MEMBERS", Ordinal_: 106},
	{Name_: "MEMBER_FUNCTION_DEF", Ordinal_: 107},
	{Name_: "MEMBER_FIELD_DEF", Ordinal_: 108},
	{Name_: "COMPUTED_FIELD_DEF", Ordinal_: 109},
	{Name_: "SUPER", Ordinal_: 110},
	{Name_: "LET", Ordinal_: 111},
	{Name_: "FOR_OF", Ordinal_: 112},
	{Name_: "FOR_AWAIT_OF", Ordinal_: 113},
	{Name_: "YIELD", Ordinal_: 114},
	{Name_: "AWAIT", Ordinal_: 115},
	{Name_: "IMPORT", Ordinal_: 116},
	{Name_: "IMPORT_SPECS", Ordinal_: 117},
	{Name_: "IMPORT_SPEC", Ordinal_: 118},
	{Name_: "IMPORT_STAR", Ordinal_: 119},
	{Name_: "EXPORT", Ordinal_: 120},
	{Name_: "EXPORT_SPECS", Ordinal_: 121},
	{Name_: "EXPORT_SPEC", Ordinal_: 122},
	{Name_: "MODULE_BODY", Ordinal_: 123},
	{Name_: "DYNAMIC_IMPORT", Ordinal_: 124},
	{Name_: "ITER_REST", Ordinal_: 125},
	{Name_: "OBJECT_REST", Ordinal_: 126},
	{Name_: "ITER_SPREAD", Ordinal_: 127},
	{Name_: "OBJECT_SPREAD", Ordinal_: 128},
	{Name_: "COMPUTED_PROP", Ordinal_: 129},
	{Name_: "TAGGED_TEMPLATELIT", Ordinal_: 130},
	{Name_: "TEMPLATELIT", Ordinal_: 131},
	{Name_: "TEMPLATELIT_SUB", Ordinal_: 132},
	{Name_: "TEMPLATELIT_STRING", Ordinal_: 133},
	{Name_: "DEFAULT_VALUE", Ordinal_: 134},
	{Name_: "NEW_TARGET", Ordinal_: 135},
	{Name_: "IMPORT_META", Ordinal_: 136},
	{Name_: "STRING_TYPE", Ordinal_: 137},
	{Name_: "BOOLEAN_TYPE", Ordinal_: 138},
	{Name_: "NUMBER_TYPE", Ordinal_: 139},
	{Name_: "FUNCTION_TYPE", Ordinal_: 140},
	{Name_: "PARAMETERIZED_TYPE", Ordinal_: 141},
	{Name_: "UNION_TYPE", Ordinal_: 142},
	{Name_: "ANY_TYPE", Ordinal_: 143},
	{Name_: "NULLABLE_TYPE", Ordinal_: 144},
	{Name_: "VOID_TYPE", Ordinal_: 145},
	{Name_: "REST_PARAMETER_TYPE", Ordinal_: 146},
	{Name_: "NAMED_TYPE", Ordinal_: 147},
	{Name_: "OPTIONAL_PARAMETER", Ordinal_: 148},
	{Name_: "RECORD_TYPE", Ordinal_: 149},
	{Name_: "UNDEFINED_TYPE", Ordinal_: 150},
	{Name_: "ARRAY_TYPE", Ordinal_: 151},
	{Name_: "GENERIC_TYPE", Ordinal_: 152},
	{Name_: "GENERIC_TYPE_LIST", Ordinal_: 153},
	{Name_: "ANNOTATION", Ordinal_: 154},
	{Name_: "PIPE", Ordinal_: 155},
	{Name_: "STAR", Ordinal_: 156},
	{Name_: "EOC", Ordinal_: 157},
	{Name_: "QMARK", Ordinal_: 158},
	{Name_: "BANG", Ordinal_: 159},
	{Name_: "EQUALS", Ordinal_: 160},
	{Name_: "LB", Ordinal_: 161},
	{Name_: "LC", Ordinal_: 162},
	{Name_: "COLON", Ordinal_: 163},
	{Name_: "INTERFACE", Ordinal_: 164},
	{Name_: "INTERFACE_EXTENDS", Ordinal_: 165},
	{Name_: "INTERFACE_MEMBERS", Ordinal_: 166},
	{Name_: "ENUM", Ordinal_: 167},
	{Name_: "ENUM_MEMBERS", Ordinal_: 168},
	{Name_: "IMPLEMENTS", Ordinal_: 169},
	{Name_: "TYPE_ALIAS", Ordinal_: 170},
	{Name_: "DECLARE", Ordinal_: 171},
	{Name_: "MEMBER_VARIABLE_DEF", Ordinal_: 172},
	{Name_: "INDEX_SIGNATURE", Ordinal_: 173},
	{Name_: "CALL_SIGNATURE", Ordinal_: 174},
	{Name_: "NAMESPACE", Ordinal_: 175},
	{Name_: "NAMESPACE_ELEMENTS", Ordinal_: 176},
	{Name_: "PLACEHOLDER1", Ordinal_: 177},
	{Name_: "PLACEHOLDER2", Ordinal_: 178},
	{Name_: "PLACEHOLDER3", Ordinal_: 179},
}
var enumValues_Token = []Token{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127, 128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143, 144, 145, 146, 147, 148, 149, 150, 151, 152, 153, 154, 155, 156, 157, 158, 159, 160, 161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 171, 172, 173, 174, 175, 176, 177, 178, 179, 180}

func StaticGet_Token_RETURN() Token {
	return enum_Token_RETURN
}
func StaticGet_Token_BITOR() Token {
	return enum_Token_BITOR
}
func StaticGet_Token_BITXOR() Token {
	return enum_Token_BITXOR
}
func StaticGet_Token_BITAND() Token {
	return enum_Token_BITAND
}
func StaticGet_Token_EQ() Token {
	return enum_Token_EQ
}
func StaticGet_Token_NE() Token {
	return enum_Token_NE
}
func StaticGet_Token_LT() Token {
	return enum_Token_LT
}
func StaticGet_Token_LE() Token {
	return enum_Token_LE
}
func StaticGet_Token_GT() Token {
	return enum_Token_GT
}
func StaticGet_Token_GE() Token {
	return enum_Token_GE
}
func StaticGet_Token_LSH() Token {
	return enum_Token_LSH
}
func StaticGet_Token_RSH() Token {
	return enum_Token_RSH
}
func StaticGet_Token_URSH() Token {
	return enum_Token_URSH
}
func StaticGet_Token_ADD() Token {
	return enum_Token_ADD
}
func StaticGet_Token_SUB() Token {
	return enum_Token_SUB
}
func StaticGet_Token_MUL() Token {
	return enum_Token_MUL
}
func StaticGet_Token_DIV() Token {
	return enum_Token_DIV
}
func StaticGet_Token_MOD() Token {
	return enum_Token_MOD
}
func StaticGet_Token_EXPONENT() Token {
	return enum_Token_EXPONENT
}
func StaticGet_Token_NOT() Token {
	return enum_Token_NOT
}
func StaticGet_Token_BITNOT() Token {
	return enum_Token_BITNOT
}
func StaticGet_Token_POS() Token {
	return enum_Token_POS
}
func StaticGet_Token_NEG() Token {
	return enum_Token_NEG
}
func StaticGet_Token_NEW() Token {
	return enum_Token_NEW
}
func StaticGet_Token_DELPROP() Token {
	return enum_Token_DELPROP
}
func StaticGet_Token_TYPEOF() Token {
	return enum_Token_TYPEOF
}
func StaticGet_Token_GETPROP() Token {
	return enum_Token_GETPROP
}
func StaticGet_Token_GETELEM() Token {
	return enum_Token_GETELEM
}
func StaticGet_Token_CALL() Token {
	return enum_Token_CALL
}
func StaticGet_Token_OPTCHAIN_GETPROP() Token {
	return enum_Token_OPTCHAIN_GETPROP
}
func StaticGet_Token_OPTCHAIN_GETELEM() Token {
	return enum_Token_OPTCHAIN_GETELEM
}
func StaticGet_Token_OPTCHAIN_CALL() Token {
	return enum_Token_OPTCHAIN_CALL
}
func StaticGet_Token_NAME() Token {
	return enum_Token_NAME
}
func StaticGet_Token_NUMBER() Token {
	return enum_Token_NUMBER
}
func StaticGet_Token_BIGINT() Token {
	return enum_Token_BIGINT
}
func StaticGet_Token_STRINGLIT() Token {
	return enum_Token_STRINGLIT
}
func StaticGet_Token_NULL() Token {
	return enum_Token_NULL
}
func StaticGet_Token_THIS() Token {
	return enum_Token_THIS
}
func StaticGet_Token_FALSE() Token {
	return enum_Token_FALSE
}
func StaticGet_Token_TRUE() Token {
	return enum_Token_TRUE
}
func StaticGet_Token_SHEQ() Token {
	return enum_Token_SHEQ
}
func StaticGet_Token_SHNE() Token {
	return enum_Token_SHNE
}
func StaticGet_Token_REGEXP() Token {
	return enum_Token_REGEXP
}
func StaticGet_Token_THROW() Token {
	return enum_Token_THROW
}
func StaticGet_Token_IN() Token {
	return enum_Token_IN
}
func StaticGet_Token_INSTANCEOF() Token {
	return enum_Token_INSTANCEOF
}
func StaticGet_Token_ARRAYLIT() Token {
	return enum_Token_ARRAYLIT
}
func StaticGet_Token_OBJECTLIT() Token {
	return enum_Token_OBJECTLIT
}
func StaticGet_Token_TRY() Token {
	return enum_Token_TRY
}
func StaticGet_Token_PARAM_LIST() Token {
	return enum_Token_PARAM_LIST
}
func StaticGet_Token_COMMA() Token {
	return enum_Token_COMMA
}
func StaticGet_Token_ASSIGN() Token {
	return enum_Token_ASSIGN
}
func StaticGet_Token_ASSIGN_BITOR() Token {
	return enum_Token_ASSIGN_BITOR
}
func StaticGet_Token_ASSIGN_BITXOR() Token {
	return enum_Token_ASSIGN_BITXOR
}
func StaticGet_Token_ASSIGN_BITAND() Token {
	return enum_Token_ASSIGN_BITAND
}
func StaticGet_Token_ASSIGN_LSH() Token {
	return enum_Token_ASSIGN_LSH
}
func StaticGet_Token_ASSIGN_RSH() Token {
	return enum_Token_ASSIGN_RSH
}
func StaticGet_Token_ASSIGN_URSH() Token {
	return enum_Token_ASSIGN_URSH
}
func StaticGet_Token_ASSIGN_ADD() Token {
	return enum_Token_ASSIGN_ADD
}
func StaticGet_Token_ASSIGN_SUB() Token {
	return enum_Token_ASSIGN_SUB
}
func StaticGet_Token_ASSIGN_MUL() Token {
	return enum_Token_ASSIGN_MUL
}
func StaticGet_Token_ASSIGN_DIV() Token {
	return enum_Token_ASSIGN_DIV
}
func StaticGet_Token_ASSIGN_MOD() Token {
	return enum_Token_ASSIGN_MOD
}
func StaticGet_Token_ASSIGN_EXPONENT() Token {
	return enum_Token_ASSIGN_EXPONENT
}
func StaticGet_Token_ASSIGN_OR() Token {
	return enum_Token_ASSIGN_OR
}
func StaticGet_Token_ASSIGN_AND() Token {
	return enum_Token_ASSIGN_AND
}
func StaticGet_Token_ASSIGN_COALESCE() Token {
	return enum_Token_ASSIGN_COALESCE
}
func StaticGet_Token_HOOK() Token {
	return enum_Token_HOOK
}
func StaticGet_Token_OR() Token {
	return enum_Token_OR
}
func StaticGet_Token_AND() Token {
	return enum_Token_AND
}
func StaticGet_Token_COALESCE() Token {
	return enum_Token_COALESCE
}
func StaticGet_Token_INC() Token {
	return enum_Token_INC
}
func StaticGet_Token_DEC() Token {
	return enum_Token_DEC
}
func StaticGet_Token_FUNCTION() Token {
	return enum_Token_FUNCTION
}
func StaticGet_Token_IF() Token {
	return enum_Token_IF
}
func StaticGet_Token_SWITCH() Token {
	return enum_Token_SWITCH
}
func StaticGet_Token_CASE() Token {
	return enum_Token_CASE
}
func StaticGet_Token_DEFAULT_CASE() Token {
	return enum_Token_DEFAULT_CASE
}
func StaticGet_Token_WHILE() Token {
	return enum_Token_WHILE
}
func StaticGet_Token_DO() Token {
	return enum_Token_DO
}
func StaticGet_Token_FOR() Token {
	return enum_Token_FOR
}
func StaticGet_Token_FOR_IN() Token {
	return enum_Token_FOR_IN
}
func StaticGet_Token_BREAK() Token {
	return enum_Token_BREAK
}
func StaticGet_Token_CONTINUE() Token {
	return enum_Token_CONTINUE
}
func StaticGet_Token_VAR() Token {
	return enum_Token_VAR
}
func StaticGet_Token_WITH() Token {
	return enum_Token_WITH
}
func StaticGet_Token_CATCH() Token {
	return enum_Token_CATCH
}
func StaticGet_Token_VOID() Token {
	return enum_Token_VOID
}
func StaticGet_Token_EMPTY() Token {
	return enum_Token_EMPTY
}
func StaticGet_Token_ROOT() Token {
	return enum_Token_ROOT
}
func StaticGet_Token_BLOCK() Token {
	return enum_Token_BLOCK
}
func StaticGet_Token_SWITCH_BODY() Token {
	return enum_Token_SWITCH_BODY
}
func StaticGet_Token_LABEL() Token {
	return enum_Token_LABEL
}
func StaticGet_Token_EXPR_RESULT() Token {
	return enum_Token_EXPR_RESULT
}
func StaticGet_Token_SCRIPT() Token {
	return enum_Token_SCRIPT
}
func StaticGet_Token_GETTER_DEF() Token {
	return enum_Token_GETTER_DEF
}
func StaticGet_Token_SETTER_DEF() Token {
	return enum_Token_SETTER_DEF
}
func StaticGet_Token_CONST() Token {
	return enum_Token_CONST
}
func StaticGet_Token_DEBUGGER() Token {
	return enum_Token_DEBUGGER
}
func StaticGet_Token_LABEL_NAME() Token {
	return enum_Token_LABEL_NAME
}
func StaticGet_Token_STRING_KEY() Token {
	return enum_Token_STRING_KEY
}
func StaticGet_Token_CAST() Token {
	return enum_Token_CAST
}
func StaticGet_Token_ARRAY_PATTERN() Token {
	return enum_Token_ARRAY_PATTERN
}
func StaticGet_Token_OBJECT_PATTERN() Token {
	return enum_Token_OBJECT_PATTERN
}
func StaticGet_Token_DESTRUCTURING_LHS() Token {
	return enum_Token_DESTRUCTURING_LHS
}
func StaticGet_Token_CLASS() Token {
	return enum_Token_CLASS
}
func StaticGet_Token_CLASS_MEMBERS() Token {
	return enum_Token_CLASS_MEMBERS
}
func StaticGet_Token_MEMBER_FUNCTION_DEF() Token {
	return enum_Token_MEMBER_FUNCTION_DEF
}
func StaticGet_Token_MEMBER_FIELD_DEF() Token {
	return enum_Token_MEMBER_FIELD_DEF
}
func StaticGet_Token_COMPUTED_FIELD_DEF() Token {
	return enum_Token_COMPUTED_FIELD_DEF
}
func StaticGet_Token_SUPER() Token {
	return enum_Token_SUPER
}
func StaticGet_Token_LET() Token {
	return enum_Token_LET
}
func StaticGet_Token_FOR_OF() Token {
	return enum_Token_FOR_OF
}
func StaticGet_Token_FOR_AWAIT_OF() Token {
	return enum_Token_FOR_AWAIT_OF
}
func StaticGet_Token_YIELD() Token {
	return enum_Token_YIELD
}
func StaticGet_Token_AWAIT() Token {
	return enum_Token_AWAIT
}
func StaticGet_Token_IMPORT() Token {
	return enum_Token_IMPORT
}
func StaticGet_Token_IMPORT_SPECS() Token {
	return enum_Token_IMPORT_SPECS
}
func StaticGet_Token_IMPORT_SPEC() Token {
	return enum_Token_IMPORT_SPEC
}
func StaticGet_Token_IMPORT_STAR() Token {
	return enum_Token_IMPORT_STAR
}
func StaticGet_Token_EXPORT() Token {
	return enum_Token_EXPORT
}
func StaticGet_Token_EXPORT_SPECS() Token {
	return enum_Token_EXPORT_SPECS
}
func StaticGet_Token_EXPORT_SPEC() Token {
	return enum_Token_EXPORT_SPEC
}
func StaticGet_Token_MODULE_BODY() Token {
	return enum_Token_MODULE_BODY
}
func StaticGet_Token_DYNAMIC_IMPORT() Token {
	return enum_Token_DYNAMIC_IMPORT
}
func StaticGet_Token_ITER_REST() Token {
	return enum_Token_ITER_REST
}
func StaticGet_Token_OBJECT_REST() Token {
	return enum_Token_OBJECT_REST
}
func StaticGet_Token_ITER_SPREAD() Token {
	return enum_Token_ITER_SPREAD
}
func StaticGet_Token_OBJECT_SPREAD() Token {
	return enum_Token_OBJECT_SPREAD
}
func StaticGet_Token_COMPUTED_PROP() Token {
	return enum_Token_COMPUTED_PROP
}
func StaticGet_Token_TAGGED_TEMPLATELIT() Token {
	return enum_Token_TAGGED_TEMPLATELIT
}
func StaticGet_Token_TEMPLATELIT() Token {
	return enum_Token_TEMPLATELIT
}
func StaticGet_Token_TEMPLATELIT_SUB() Token {
	return enum_Token_TEMPLATELIT_SUB
}
func StaticGet_Token_TEMPLATELIT_STRING() Token {
	return enum_Token_TEMPLATELIT_STRING
}
func StaticGet_Token_DEFAULT_VALUE() Token {
	return enum_Token_DEFAULT_VALUE
}
func StaticGet_Token_NEW_TARGET() Token {
	return enum_Token_NEW_TARGET
}
func StaticGet_Token_IMPORT_META() Token {
	return enum_Token_IMPORT_META
}
func StaticGet_Token_STRING_TYPE() Token {
	return enum_Token_STRING_TYPE
}
func StaticGet_Token_BOOLEAN_TYPE() Token {
	return enum_Token_BOOLEAN_TYPE
}
func StaticGet_Token_NUMBER_TYPE() Token {
	return enum_Token_NUMBER_TYPE
}
func StaticGet_Token_FUNCTION_TYPE() Token {
	return enum_Token_FUNCTION_TYPE
}
func StaticGet_Token_PARAMETERIZED_TYPE() Token {
	return enum_Token_PARAMETERIZED_TYPE
}
func StaticGet_Token_UNION_TYPE() Token {
	return enum_Token_UNION_TYPE
}
func StaticGet_Token_ANY_TYPE() Token {
	return enum_Token_ANY_TYPE
}
func StaticGet_Token_NULLABLE_TYPE() Token {
	return enum_Token_NULLABLE_TYPE
}
func StaticGet_Token_VOID_TYPE() Token {
	return enum_Token_VOID_TYPE
}
func StaticGet_Token_REST_PARAMETER_TYPE() Token {
	return enum_Token_REST_PARAMETER_TYPE
}
func StaticGet_Token_NAMED_TYPE() Token {
	return enum_Token_NAMED_TYPE
}
func StaticGet_Token_OPTIONAL_PARAMETER() Token {
	return enum_Token_OPTIONAL_PARAMETER
}
func StaticGet_Token_RECORD_TYPE() Token {
	return enum_Token_RECORD_TYPE
}
func StaticGet_Token_UNDEFINED_TYPE() Token {
	return enum_Token_UNDEFINED_TYPE
}
func StaticGet_Token_ARRAY_TYPE() Token {
	return enum_Token_ARRAY_TYPE
}
func StaticGet_Token_GENERIC_TYPE() Token {
	return enum_Token_GENERIC_TYPE
}
func StaticGet_Token_GENERIC_TYPE_LIST() Token {
	return enum_Token_GENERIC_TYPE_LIST
}
func StaticGet_Token_ANNOTATION() Token {
	return enum_Token_ANNOTATION
}
func StaticGet_Token_PIPE() Token {
	return enum_Token_PIPE
}
func StaticGet_Token_STAR() Token {
	return enum_Token_STAR
}
func StaticGet_Token_EOC() Token {
	return enum_Token_EOC
}
func StaticGet_Token_QMARK() Token {
	return enum_Token_QMARK
}
func StaticGet_Token_BANG() Token {
	return enum_Token_BANG
}
func StaticGet_Token_EQUALS() Token {
	return enum_Token_EQUALS
}
func StaticGet_Token_LB() Token {
	return enum_Token_LB
}
func StaticGet_Token_LC() Token {
	return enum_Token_LC
}
func StaticGet_Token_COLON() Token {
	return enum_Token_COLON
}
func StaticGet_Token_INTERFACE() Token {
	return enum_Token_INTERFACE
}
func StaticGet_Token_INTERFACE_EXTENDS() Token {
	return enum_Token_INTERFACE_EXTENDS
}
func StaticGet_Token_INTERFACE_MEMBERS() Token {
	return enum_Token_INTERFACE_MEMBERS
}
func StaticGet_Token_ENUM() Token {
	return enum_Token_ENUM
}
func StaticGet_Token_ENUM_MEMBERS() Token {
	return enum_Token_ENUM_MEMBERS
}
func StaticGet_Token_IMPLEMENTS() Token {
	return enum_Token_IMPLEMENTS
}
func StaticGet_Token_TYPE_ALIAS() Token {
	return enum_Token_TYPE_ALIAS
}
func StaticGet_Token_DECLARE() Token {
	return enum_Token_DECLARE
}
func StaticGet_Token_MEMBER_VARIABLE_DEF() Token {
	return enum_Token_MEMBER_VARIABLE_DEF
}
func StaticGet_Token_INDEX_SIGNATURE() Token {
	return enum_Token_INDEX_SIGNATURE
}
func StaticGet_Token_CALL_SIGNATURE() Token {
	return enum_Token_CALL_SIGNATURE
}
func StaticGet_Token_NAMESPACE() Token {
	return enum_Token_NAMESPACE
}
func StaticGet_Token_NAMESPACE_ELEMENTS() Token {
	return enum_Token_NAMESPACE_ELEMENTS
}
func StaticGet_Token_PLACEHOLDER1() Token {
	return enum_Token_PLACEHOLDER1
}
func StaticGet_Token_PLACEHOLDER2() Token {
	return enum_Token_PLACEHOLDER2
}
func StaticGet_Token_PLACEHOLDER3() Token {
	return enum_Token_PLACEHOLDER3
}

func Token_Values() []Token {
	return j2gort.CloneSlice(enumValues_Token)
}

func Token_ValueOf(name string) Token {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_Token {
		if enumBaseTable_Token[i_].Name_ == name {
			return Token(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.rhino.Token." + name))
}

// Java static method: Token.arity(com.google.javascript.rhino.Token token)
func Token_Arity(token Token) int32 {
	_ = token
	switch (token).Ordinal() {
	case int32(154) /* Token.ANNOTATION */, int32(46) /* Token.ARRAYLIT */, int32(159) /* Token.BANG */, int32(90) /* Token.BLOCK */, int32(89) /* Token.ROOT */, int32(82) /* Token.BREAK */, int32(28) /* Token.CALL */, int32(31) /* Token.OPTCHAIN_CALL */, int32(163) /* Token.COLON */, int32(97) /* Token.CONST */, int32(83) /* Token.CONTINUE */, int32(98) /* Token.DEBUGGER */, int32(157) /* Token.EOC */, int32(160) /* Token.EQUALS */, int32(80) /* Token.FOR */, int32(74) /* Token.IF */, int32(161) /* Token.LB */, int32(162) /* Token.LC */, int32(23) /* Token.NEW */, int32(47) /* Token.OBJECTLIT */, int32(49) /* Token.PARAM_LIST */, int32(155) /* Token.PIPE */, int32(158) /* Token.QMARK */, int32(42) /* Token.REGEXP */, int32(0) /* Token.RETURN */, int32(94) /* Token.SCRIPT */, int32(156) /* Token.STAR */, int32(100) /* Token.STRING_KEY */, int32(91) /* Token.SWITCH_BODY */, int32(131) /* Token.TEMPLATELIT */, int32(48) /* Token.TRY */, int32(84) /* Token.VAR */, int32(114) /* Token.YIELD */, int32(108) /* Token.MEMBER_FIELD_DEF */, int32(109) /* Token.COMPUTED_FIELD_DEF */ :
		return int32(-1)
	case int32(88) /* Token.EMPTY */, int32(38) /* Token.FALSE */, int32(119) /* Token.IMPORT_STAR */, int32(99) /* Token.LABEL_NAME */, int32(172) /* Token.MEMBER_VARIABLE_DEF */, int32(32) /* Token.NAME */, int32(36) /* Token.NULL */, int32(33) /* Token.NUMBER */, int32(34) /* Token.BIGINT */, int32(35) /* Token.STRINGLIT */, int32(133) /* Token.TEMPLATELIT_STRING */, int32(37) /* Token.THIS */, int32(39) /* Token.TRUE */ :
		return int32(0)
	case int32(115) /* Token.AWAIT */, int32(20) /* Token.BITNOT */, int32(174) /* Token.CALL_SIGNATURE */, int32(101) /* Token.CAST */, int32(72) /* Token.DEC */, int32(77) /* Token.DEFAULT_CASE */, int32(24) /* Token.DELPROP */, int32(93) /* Token.EXPR_RESULT */, int32(26) /* Token.GETPROP */, int32(95) /* Token.GETTER_DEF */, int32(71) /* Token.INC */, int32(173) /* Token.INDEX_SIGNATURE */, int32(125) /* Token.ITER_REST */, int32(127) /* Token.ITER_SPREAD */, int32(107) /* Token.MEMBER_FUNCTION_DEF */, int32(147) /* Token.NAMED_TYPE */, int32(22) /* Token.NEG */, int32(19) /* Token.NOT */, int32(126) /* Token.OBJECT_REST */, int32(128) /* Token.OBJECT_SPREAD */, int32(29) /* Token.OPTCHAIN_GETPROP */, int32(21) /* Token.POS */, int32(96) /* Token.SETTER_DEF */, int32(132) /* Token.TEMPLATELIT_SUB */, int32(43) /* Token.THROW */, int32(25) /* Token.TYPEOF */, int32(170) /* Token.TYPE_ALIAS */, int32(87) /* Token.VOID */ :
		return int32(1)
	case int32(13) /* Token.ADD */, int32(69) /* Token.AND */, int32(51) /* Token.ASSIGN */, int32(58) /* Token.ASSIGN_ADD */, int32(54) /* Token.ASSIGN_BITAND */, int32(52) /* Token.ASSIGN_BITOR */, int32(53) /* Token.ASSIGN_BITXOR */, int32(61) /* Token.ASSIGN_DIV */, int32(55) /* Token.ASSIGN_LSH */, int32(62) /* Token.ASSIGN_MOD */, int32(60) /* Token.ASSIGN_MUL */, int32(63) /* Token.ASSIGN_EXPONENT */, int32(56) /* Token.ASSIGN_RSH */, int32(59) /* Token.ASSIGN_SUB */, int32(57) /* Token.ASSIGN_URSH */, int32(64) /* Token.ASSIGN_OR */, int32(65) /* Token.ASSIGN_AND */, int32(66) /* Token.ASSIGN_COALESCE */, int32(3) /* Token.BITAND */, int32(1) /* Token.BITOR */, int32(2) /* Token.BITXOR */, int32(76) /* Token.CASE */, int32(70) /* Token.COALESCE */, int32(86) /* Token.CATCH */, int32(50) /* Token.COMMA */, int32(129) /* Token.COMPUTED_PROP */, int32(134) /* Token.DEFAULT_VALUE */, int32(16) /* Token.DIV */, int32(79) /* Token.DO */, int32(167) /* Token.ENUM */, int32(4) /* Token.EQ */, int32(18) /* Token.EXPONENT */, int32(9) /* Token.GE */, int32(27) /* Token.GETELEM */, int32(30) /* Token.OPTCHAIN_GETELEM */, int32(8) /* Token.GT */, int32(44) /* Token.IN */, int32(45) /* Token.INSTANCEOF */, int32(92) /* Token.LABEL */, int32(7) /* Token.LE */, int32(10) /* Token.LSH */, int32(6) /* Token.LT */, int32(17) /* Token.MOD */, int32(15) /* Token.MUL */, int32(175) /* Token.NAMESPACE */, int32(5) /* Token.NE */, int32(68) /* Token.OR */, int32(11) /* Token.RSH */, int32(75) /* Token.SWITCH */, int32(40) /* Token.SHEQ */, int32(41) /* Token.SHNE */, int32(14) /* Token.SUB */, int32(130) /* Token.TAGGED_TEMPLATELIT */, int32(12) /* Token.URSH */, int32(78) /* Token.WHILE */, int32(85) /* Token.WITH */ :
		return int32(2)
	case int32(105) /* Token.CLASS */, int32(81) /* Token.FOR_IN */, int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */, int32(73) /* Token.FUNCTION */, int32(67) /* Token.HOOK */, int32(116) /* Token.IMPORT */, int32(164) /* Token.INTERFACE */ :
		return int32(3)
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "No arity defined for "+j2gort.StrVal(token))))
	}
	panic("unreachable switch expression")
}

func (this_ Token) IsJavaNil_() bool                              { return this_ == 0 }
func (Token) HasIdentityEquals_()                                 {}
func (Token) Implements_j2gort_Ext_java_lang_Enum()               {}
func (Token) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ Token) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (Token) JavaClassInfo_() (string, string) { return "com.google.javascript.rhino.Token", "Token" }
func (this_ Token) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_Token[this_-1]
}
func (this_ Token) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ Token) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ Token) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Token[this_-1].Name_
}
func (this_ Token) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_Token[this_-1].Name_
}
func (this_ Token) Equals(other_ any) bool {
	o_, ok_ := other_.(Token)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ Token) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_Token[this_-1].HashCode()
}
func (this_ Token) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Token[this_-1].CompareTo(other_)
}
func (this_ Token) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_Token[this_-1].Name_)
}
func init() {
	vals_ := []any{Token(1), Token(2), Token(3), Token(4), Token(5), Token(6), Token(7), Token(8), Token(9), Token(10), Token(11), Token(12), Token(13), Token(14), Token(15), Token(16), Token(17), Token(18), Token(19), Token(20), Token(21), Token(22), Token(23), Token(24), Token(25), Token(26), Token(27), Token(28), Token(29), Token(30), Token(31), Token(32), Token(33), Token(34), Token(35), Token(36), Token(37), Token(38), Token(39), Token(40), Token(41), Token(42), Token(43), Token(44), Token(45), Token(46), Token(47), Token(48), Token(49), Token(50), Token(51), Token(52), Token(53), Token(54), Token(55), Token(56), Token(57), Token(58), Token(59), Token(60), Token(61), Token(62), Token(63), Token(64), Token(65), Token(66), Token(67), Token(68), Token(69), Token(70), Token(71), Token(72), Token(73), Token(74), Token(75), Token(76), Token(77), Token(78), Token(79), Token(80), Token(81), Token(82), Token(83), Token(84), Token(85), Token(86), Token(87), Token(88), Token(89), Token(90), Token(91), Token(92), Token(93), Token(94), Token(95), Token(96), Token(97), Token(98), Token(99), Token(100), Token(101), Token(102), Token(103), Token(104), Token(105), Token(106), Token(107), Token(108), Token(109), Token(110), Token(111), Token(112), Token(113), Token(114), Token(115), Token(116), Token(117), Token(118), Token(119), Token(120), Token(121), Token(122), Token(123), Token(124), Token(125), Token(126), Token(127), Token(128), Token(129), Token(130), Token(131), Token(132), Token(133), Token(134), Token(135), Token(136), Token(137), Token(138), Token(139), Token(140), Token(141), Token(142), Token(143), Token(144), Token(145), Token(146), Token(147), Token(148), Token(149), Token(150), Token(151), Token(152), Token(153), Token(154), Token(155), Token(156), Token(157), Token(158), Token(159), Token(160), Token(161), Token(162), Token(163), Token(164), Token(165), Token(166), Token(167), Token(168), Token(169), Token(170), Token(171), Token(172), Token(173), Token(174), Token(175), Token(176), Token(177), Token(178), Token(179), Token(180)}
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.Token", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.Token", "com.google.javascript.rhino.Token", func(v_ any) bool {
		return j2gort.InstanceOf[Token](v_)
	})
}
