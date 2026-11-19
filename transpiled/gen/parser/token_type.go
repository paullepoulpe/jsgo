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
// Port of com/google/javascript/jscomp/parsing/parser/TokenType.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java enum: com.google.javascript.jscomp.parsing.parser.TokenType
type TokenType uint8

const (
	enum_TokenType_END_OF_FILE TokenType = 1
	enum_TokenType_ERROR TokenType = 2
	enum_TokenType_IDENTIFIER TokenType = 3
	enum_TokenType_BREAK TokenType = 4
	enum_TokenType_CASE TokenType = 5
	enum_TokenType_CATCH TokenType = 6
	enum_TokenType_CONTINUE TokenType = 7
	enum_TokenType_DEBUGGER TokenType = 8
	enum_TokenType_DEFAULT TokenType = 9
	enum_TokenType_DELETE TokenType = 10
	enum_TokenType_DO TokenType = 11
	enum_TokenType_ELSE TokenType = 12
	enum_TokenType_FINALLY TokenType = 13
	enum_TokenType_FOR TokenType = 14
	enum_TokenType_FUNCTION TokenType = 15
	enum_TokenType_IF TokenType = 16
	enum_TokenType_IN TokenType = 17
	enum_TokenType_INSTANCEOF TokenType = 18
	enum_TokenType_NEW TokenType = 19
	enum_TokenType_RETURN TokenType = 20
	enum_TokenType_SWITCH TokenType = 21
	enum_TokenType_THIS TokenType = 22
	enum_TokenType_THROW TokenType = 23
	enum_TokenType_TRY TokenType = 24
	enum_TokenType_TYPEOF TokenType = 25
	enum_TokenType_VAR TokenType = 26
	enum_TokenType_VOID TokenType = 27
	enum_TokenType_WHILE TokenType = 28
	enum_TokenType_WITH TokenType = 29
	enum_TokenType_CLASS TokenType = 30
	enum_TokenType_CONST TokenType = 31
	enum_TokenType_ENUM TokenType = 32
	enum_TokenType_EXPORT TokenType = 33
	enum_TokenType_EXTENDS TokenType = 34
	enum_TokenType_IMPORT TokenType = 35
	enum_TokenType_SUPER TokenType = 36
	enum_TokenType_IMPLEMENTS TokenType = 37
	enum_TokenType_INTERFACE TokenType = 38
	enum_TokenType_LET TokenType = 39
	enum_TokenType_PACKAGE TokenType = 40
	enum_TokenType_PRIVATE TokenType = 41
	enum_TokenType_PROTECTED TokenType = 42
	enum_TokenType_PUBLIC TokenType = 43
	enum_TokenType_STATIC TokenType = 44
	enum_TokenType_YIELD TokenType = 45
	enum_TokenType_OPEN_CURLY TokenType = 46
	enum_TokenType_CLOSE_CURLY TokenType = 47
	enum_TokenType_OPEN_PAREN TokenType = 48
	enum_TokenType_CLOSE_PAREN TokenType = 49
	enum_TokenType_OPEN_SQUARE TokenType = 50
	enum_TokenType_CLOSE_SQUARE TokenType = 51
	enum_TokenType_PERIOD TokenType = 52
	enum_TokenType_SEMI_COLON TokenType = 53
	enum_TokenType_COMMA TokenType = 54
	enum_TokenType_OPEN_ANGLE TokenType = 55
	enum_TokenType_CLOSE_ANGLE TokenType = 56
	enum_TokenType_LESS_EQUAL TokenType = 57
	enum_TokenType_GREATER_EQUAL TokenType = 58
	enum_TokenType_ARROW TokenType = 59
	enum_TokenType_EQUAL_EQUAL TokenType = 60
	enum_TokenType_NOT_EQUAL TokenType = 61
	enum_TokenType_EQUAL_EQUAL_EQUAL TokenType = 62
	enum_TokenType_NOT_EQUAL_EQUAL TokenType = 63
	enum_TokenType_PLUS TokenType = 64
	enum_TokenType_MINUS TokenType = 65
	enum_TokenType_STAR TokenType = 66
	enum_TokenType_STAR_STAR TokenType = 67
	enum_TokenType_PERCENT TokenType = 68
	enum_TokenType_PLUS_PLUS TokenType = 69
	enum_TokenType_MINUS_MINUS TokenType = 70
	enum_TokenType_LEFT_SHIFT TokenType = 71
	enum_TokenType_RIGHT_SHIFT TokenType = 72
	enum_TokenType_UNSIGNED_RIGHT_SHIFT TokenType = 73
	enum_TokenType_AMPERSAND TokenType = 74
	enum_TokenType_BAR TokenType = 75
	enum_TokenType_CARET TokenType = 76
	enum_TokenType_BANG TokenType = 77
	enum_TokenType_TILDE TokenType = 78
	enum_TokenType_AND TokenType = 79
	enum_TokenType_OR TokenType = 80
	enum_TokenType_QUESTION TokenType = 81
	enum_TokenType_QUESTION_QUESTION TokenType = 82
	enum_TokenType_QUESTION_DOT TokenType = 83
	enum_TokenType_COLON TokenType = 84
	enum_TokenType_EQUAL TokenType = 85
	enum_TokenType_PLUS_EQUAL TokenType = 86
	enum_TokenType_MINUS_EQUAL TokenType = 87
	enum_TokenType_STAR_EQUAL TokenType = 88
	enum_TokenType_STAR_STAR_EQUAL TokenType = 89
	enum_TokenType_PERCENT_EQUAL TokenType = 90
	enum_TokenType_LEFT_SHIFT_EQUAL TokenType = 91
	enum_TokenType_RIGHT_SHIFT_EQUAL TokenType = 92
	enum_TokenType_UNSIGNED_RIGHT_SHIFT_EQUAL TokenType = 93
	enum_TokenType_AMPERSAND_EQUAL TokenType = 94
	enum_TokenType_BAR_EQUAL TokenType = 95
	enum_TokenType_CARET_EQUAL TokenType = 96
	enum_TokenType_SLASH TokenType = 97
	enum_TokenType_SLASH_EQUAL TokenType = 98
	enum_TokenType_AND_EQUAL TokenType = 99
	enum_TokenType_OR_EQUAL TokenType = 100
	enum_TokenType_QUESTION_QUESTION_EQUAL TokenType = 101
	enum_TokenType_NULL TokenType = 102
	enum_TokenType_TRUE TokenType = 103
	enum_TokenType_FALSE TokenType = 104
	enum_TokenType_NUMBER TokenType = 105
	enum_TokenType_STRING TokenType = 106
	enum_TokenType_BIGINT TokenType = 107
	enum_TokenType_REGULAR_EXPRESSION TokenType = 108
	enum_TokenType_ELLIPSIS TokenType = 109
	enum_TokenType_TEMPLATE_HEAD TokenType = 110
	enum_TokenType_TEMPLATE_MIDDLE TokenType = 111
	enum_TokenType_TEMPLATE_TAIL TokenType = 112
	enum_TokenType_NO_SUBSTITUTION_TEMPLATE TokenType = 113
	enum_TokenType_TYPE TokenType = 114
	enum_TokenType_DECLARE TokenType = 115
	enum_TokenType_MODULE TokenType = 116
	enum_TokenType_NAMESPACE TokenType = 117
)

var enumBaseTable_TokenType = [...]j2gort.EnumBase{
	{Name_: "END_OF_FILE", Ordinal_: 0},
	{Name_: "ERROR", Ordinal_: 1},
	{Name_: "IDENTIFIER", Ordinal_: 2},
	{Name_: "BREAK", Ordinal_: 3},
	{Name_: "CASE", Ordinal_: 4},
	{Name_: "CATCH", Ordinal_: 5},
	{Name_: "CONTINUE", Ordinal_: 6},
	{Name_: "DEBUGGER", Ordinal_: 7},
	{Name_: "DEFAULT", Ordinal_: 8},
	{Name_: "DELETE", Ordinal_: 9},
	{Name_: "DO", Ordinal_: 10},
	{Name_: "ELSE", Ordinal_: 11},
	{Name_: "FINALLY", Ordinal_: 12},
	{Name_: "FOR", Ordinal_: 13},
	{Name_: "FUNCTION", Ordinal_: 14},
	{Name_: "IF", Ordinal_: 15},
	{Name_: "IN", Ordinal_: 16},
	{Name_: "INSTANCEOF", Ordinal_: 17},
	{Name_: "NEW", Ordinal_: 18},
	{Name_: "RETURN", Ordinal_: 19},
	{Name_: "SWITCH", Ordinal_: 20},
	{Name_: "THIS", Ordinal_: 21},
	{Name_: "THROW", Ordinal_: 22},
	{Name_: "TRY", Ordinal_: 23},
	{Name_: "TYPEOF", Ordinal_: 24},
	{Name_: "VAR", Ordinal_: 25},
	{Name_: "VOID", Ordinal_: 26},
	{Name_: "WHILE", Ordinal_: 27},
	{Name_: "WITH", Ordinal_: 28},
	{Name_: "CLASS", Ordinal_: 29},
	{Name_: "CONST", Ordinal_: 30},
	{Name_: "ENUM", Ordinal_: 31},
	{Name_: "EXPORT", Ordinal_: 32},
	{Name_: "EXTENDS", Ordinal_: 33},
	{Name_: "IMPORT", Ordinal_: 34},
	{Name_: "SUPER", Ordinal_: 35},
	{Name_: "IMPLEMENTS", Ordinal_: 36},
	{Name_: "INTERFACE", Ordinal_: 37},
	{Name_: "LET", Ordinal_: 38},
	{Name_: "PACKAGE", Ordinal_: 39},
	{Name_: "PRIVATE", Ordinal_: 40},
	{Name_: "PROTECTED", Ordinal_: 41},
	{Name_: "PUBLIC", Ordinal_: 42},
	{Name_: "STATIC", Ordinal_: 43},
	{Name_: "YIELD", Ordinal_: 44},
	{Name_: "OPEN_CURLY", Ordinal_: 45},
	{Name_: "CLOSE_CURLY", Ordinal_: 46},
	{Name_: "OPEN_PAREN", Ordinal_: 47},
	{Name_: "CLOSE_PAREN", Ordinal_: 48},
	{Name_: "OPEN_SQUARE", Ordinal_: 49},
	{Name_: "CLOSE_SQUARE", Ordinal_: 50},
	{Name_: "PERIOD", Ordinal_: 51},
	{Name_: "SEMI_COLON", Ordinal_: 52},
	{Name_: "COMMA", Ordinal_: 53},
	{Name_: "OPEN_ANGLE", Ordinal_: 54},
	{Name_: "CLOSE_ANGLE", Ordinal_: 55},
	{Name_: "LESS_EQUAL", Ordinal_: 56},
	{Name_: "GREATER_EQUAL", Ordinal_: 57},
	{Name_: "ARROW", Ordinal_: 58},
	{Name_: "EQUAL_EQUAL", Ordinal_: 59},
	{Name_: "NOT_EQUAL", Ordinal_: 60},
	{Name_: "EQUAL_EQUAL_EQUAL", Ordinal_: 61},
	{Name_: "NOT_EQUAL_EQUAL", Ordinal_: 62},
	{Name_: "PLUS", Ordinal_: 63},
	{Name_: "MINUS", Ordinal_: 64},
	{Name_: "STAR", Ordinal_: 65},
	{Name_: "STAR_STAR", Ordinal_: 66},
	{Name_: "PERCENT", Ordinal_: 67},
	{Name_: "PLUS_PLUS", Ordinal_: 68},
	{Name_: "MINUS_MINUS", Ordinal_: 69},
	{Name_: "LEFT_SHIFT", Ordinal_: 70},
	{Name_: "RIGHT_SHIFT", Ordinal_: 71},
	{Name_: "UNSIGNED_RIGHT_SHIFT", Ordinal_: 72},
	{Name_: "AMPERSAND", Ordinal_: 73},
	{Name_: "BAR", Ordinal_: 74},
	{Name_: "CARET", Ordinal_: 75},
	{Name_: "BANG", Ordinal_: 76},
	{Name_: "TILDE", Ordinal_: 77},
	{Name_: "AND", Ordinal_: 78},
	{Name_: "OR", Ordinal_: 79},
	{Name_: "QUESTION", Ordinal_: 80},
	{Name_: "QUESTION_QUESTION", Ordinal_: 81},
	{Name_: "QUESTION_DOT", Ordinal_: 82},
	{Name_: "COLON", Ordinal_: 83},
	{Name_: "EQUAL", Ordinal_: 84},
	{Name_: "PLUS_EQUAL", Ordinal_: 85},
	{Name_: "MINUS_EQUAL", Ordinal_: 86},
	{Name_: "STAR_EQUAL", Ordinal_: 87},
	{Name_: "STAR_STAR_EQUAL", Ordinal_: 88},
	{Name_: "PERCENT_EQUAL", Ordinal_: 89},
	{Name_: "LEFT_SHIFT_EQUAL", Ordinal_: 90},
	{Name_: "RIGHT_SHIFT_EQUAL", Ordinal_: 91},
	{Name_: "UNSIGNED_RIGHT_SHIFT_EQUAL", Ordinal_: 92},
	{Name_: "AMPERSAND_EQUAL", Ordinal_: 93},
	{Name_: "BAR_EQUAL", Ordinal_: 94},
	{Name_: "CARET_EQUAL", Ordinal_: 95},
	{Name_: "SLASH", Ordinal_: 96},
	{Name_: "SLASH_EQUAL", Ordinal_: 97},
	{Name_: "AND_EQUAL", Ordinal_: 98},
	{Name_: "OR_EQUAL", Ordinal_: 99},
	{Name_: "QUESTION_QUESTION_EQUAL", Ordinal_: 100},
	{Name_: "NULL", Ordinal_: 101},
	{Name_: "TRUE", Ordinal_: 102},
	{Name_: "FALSE", Ordinal_: 103},
	{Name_: "NUMBER", Ordinal_: 104},
	{Name_: "STRING", Ordinal_: 105},
	{Name_: "BIGINT", Ordinal_: 106},
	{Name_: "REGULAR_EXPRESSION", Ordinal_: 107},
	{Name_: "ELLIPSIS", Ordinal_: 108},
	{Name_: "TEMPLATE_HEAD", Ordinal_: 109},
	{Name_: "TEMPLATE_MIDDLE", Ordinal_: 110},
	{Name_: "TEMPLATE_TAIL", Ordinal_: 111},
	{Name_: "NO_SUBSTITUTION_TEMPLATE", Ordinal_: 112},
	{Name_: "TYPE", Ordinal_: 113},
	{Name_: "DECLARE", Ordinal_: 114},
	{Name_: "MODULE", Ordinal_: 115},
	{Name_: "NAMESPACE", Ordinal_: 116},
}
var enumValues_TokenType = []TokenType{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117}

var enumField_TokenType_value = [...]string{
	"End of File",
	"error",
	"identifier",
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	"{",
	"}",
	"(",
	")",
	"[",
	"]",
	".",
	";",
	",",
	"<",
	">",
	"<=",
	">=",
	"=>",
	"==",
	"!=",
	"===",
	"!==",
	"+",
	"-",
	"*",
	"**",
	"%",
	"++",
	"--",
	"<<",
	">>",
	">>>",
	"&",
	"|",
	"^",
	"!",
	"~",
	"&&",
	"||",
	"?",
	"??",
	"?.",
	":",
	"=",
	"+=",
	"-=",
	"*=",
	"**=",
	"%=",
	"<<=",
	">>=",
	">>>=",
	"&=",
	"|=",
	"^=",
	"/",
	"/=",
	"&&=",
	"||=",
	"??=",
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	"number literal",
	"string literal",
	"bigint literal",
	"regular expression literal",
	"...",
	"template head",
	"template middle",
	"template tail",
	"no substitution template",
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
	j2gort.NullString,
}
func (this_ TokenType) FieldGet_value() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumField_TokenType_value[this_-1]
}

func StaticGet_TokenType_END_OF_FILE() TokenType {
	return enum_TokenType_END_OF_FILE
}
func StaticGet_TokenType_ERROR() TokenType {
	return enum_TokenType_ERROR
}
func StaticGet_TokenType_IDENTIFIER() TokenType {
	return enum_TokenType_IDENTIFIER
}
func StaticGet_TokenType_BREAK() TokenType {
	return enum_TokenType_BREAK
}
func StaticGet_TokenType_CASE() TokenType {
	return enum_TokenType_CASE
}
func StaticGet_TokenType_CATCH() TokenType {
	return enum_TokenType_CATCH
}
func StaticGet_TokenType_CONTINUE() TokenType {
	return enum_TokenType_CONTINUE
}
func StaticGet_TokenType_DEBUGGER() TokenType {
	return enum_TokenType_DEBUGGER
}
func StaticGet_TokenType_DEFAULT() TokenType {
	return enum_TokenType_DEFAULT
}
func StaticGet_TokenType_DELETE() TokenType {
	return enum_TokenType_DELETE
}
func StaticGet_TokenType_DO() TokenType {
	return enum_TokenType_DO
}
func StaticGet_TokenType_ELSE() TokenType {
	return enum_TokenType_ELSE
}
func StaticGet_TokenType_FINALLY() TokenType {
	return enum_TokenType_FINALLY
}
func StaticGet_TokenType_FOR() TokenType {
	return enum_TokenType_FOR
}
func StaticGet_TokenType_FUNCTION() TokenType {
	return enum_TokenType_FUNCTION
}
func StaticGet_TokenType_IF() TokenType {
	return enum_TokenType_IF
}
func StaticGet_TokenType_IN() TokenType {
	return enum_TokenType_IN
}
func StaticGet_TokenType_INSTANCEOF() TokenType {
	return enum_TokenType_INSTANCEOF
}
func StaticGet_TokenType_NEW() TokenType {
	return enum_TokenType_NEW
}
func StaticGet_TokenType_RETURN() TokenType {
	return enum_TokenType_RETURN
}
func StaticGet_TokenType_SWITCH() TokenType {
	return enum_TokenType_SWITCH
}
func StaticGet_TokenType_THIS() TokenType {
	return enum_TokenType_THIS
}
func StaticGet_TokenType_THROW() TokenType {
	return enum_TokenType_THROW
}
func StaticGet_TokenType_TRY() TokenType {
	return enum_TokenType_TRY
}
func StaticGet_TokenType_TYPEOF() TokenType {
	return enum_TokenType_TYPEOF
}
func StaticGet_TokenType_VAR() TokenType {
	return enum_TokenType_VAR
}
func StaticGet_TokenType_VOID() TokenType {
	return enum_TokenType_VOID
}
func StaticGet_TokenType_WHILE() TokenType {
	return enum_TokenType_WHILE
}
func StaticGet_TokenType_WITH() TokenType {
	return enum_TokenType_WITH
}
func StaticGet_TokenType_CLASS() TokenType {
	return enum_TokenType_CLASS
}
func StaticGet_TokenType_CONST() TokenType {
	return enum_TokenType_CONST
}
func StaticGet_TokenType_ENUM() TokenType {
	return enum_TokenType_ENUM
}
func StaticGet_TokenType_EXPORT() TokenType {
	return enum_TokenType_EXPORT
}
func StaticGet_TokenType_EXTENDS() TokenType {
	return enum_TokenType_EXTENDS
}
func StaticGet_TokenType_IMPORT() TokenType {
	return enum_TokenType_IMPORT
}
func StaticGet_TokenType_SUPER() TokenType {
	return enum_TokenType_SUPER
}
func StaticGet_TokenType_IMPLEMENTS() TokenType {
	return enum_TokenType_IMPLEMENTS
}
func StaticGet_TokenType_INTERFACE() TokenType {
	return enum_TokenType_INTERFACE
}
func StaticGet_TokenType_LET() TokenType {
	return enum_TokenType_LET
}
func StaticGet_TokenType_PACKAGE() TokenType {
	return enum_TokenType_PACKAGE
}
func StaticGet_TokenType_PRIVATE() TokenType {
	return enum_TokenType_PRIVATE
}
func StaticGet_TokenType_PROTECTED() TokenType {
	return enum_TokenType_PROTECTED
}
func StaticGet_TokenType_PUBLIC() TokenType {
	return enum_TokenType_PUBLIC
}
func StaticGet_TokenType_STATIC() TokenType {
	return enum_TokenType_STATIC
}
func StaticGet_TokenType_YIELD() TokenType {
	return enum_TokenType_YIELD
}
func StaticGet_TokenType_OPEN_CURLY() TokenType {
	return enum_TokenType_OPEN_CURLY
}
func StaticGet_TokenType_CLOSE_CURLY() TokenType {
	return enum_TokenType_CLOSE_CURLY
}
func StaticGet_TokenType_OPEN_PAREN() TokenType {
	return enum_TokenType_OPEN_PAREN
}
func StaticGet_TokenType_CLOSE_PAREN() TokenType {
	return enum_TokenType_CLOSE_PAREN
}
func StaticGet_TokenType_OPEN_SQUARE() TokenType {
	return enum_TokenType_OPEN_SQUARE
}
func StaticGet_TokenType_CLOSE_SQUARE() TokenType {
	return enum_TokenType_CLOSE_SQUARE
}
func StaticGet_TokenType_PERIOD() TokenType {
	return enum_TokenType_PERIOD
}
func StaticGet_TokenType_SEMI_COLON() TokenType {
	return enum_TokenType_SEMI_COLON
}
func StaticGet_TokenType_COMMA() TokenType {
	return enum_TokenType_COMMA
}
func StaticGet_TokenType_OPEN_ANGLE() TokenType {
	return enum_TokenType_OPEN_ANGLE
}
func StaticGet_TokenType_CLOSE_ANGLE() TokenType {
	return enum_TokenType_CLOSE_ANGLE
}
func StaticGet_TokenType_LESS_EQUAL() TokenType {
	return enum_TokenType_LESS_EQUAL
}
func StaticGet_TokenType_GREATER_EQUAL() TokenType {
	return enum_TokenType_GREATER_EQUAL
}
func StaticGet_TokenType_ARROW() TokenType {
	return enum_TokenType_ARROW
}
func StaticGet_TokenType_EQUAL_EQUAL() TokenType {
	return enum_TokenType_EQUAL_EQUAL
}
func StaticGet_TokenType_NOT_EQUAL() TokenType {
	return enum_TokenType_NOT_EQUAL
}
func StaticGet_TokenType_EQUAL_EQUAL_EQUAL() TokenType {
	return enum_TokenType_EQUAL_EQUAL_EQUAL
}
func StaticGet_TokenType_NOT_EQUAL_EQUAL() TokenType {
	return enum_TokenType_NOT_EQUAL_EQUAL
}
func StaticGet_TokenType_PLUS() TokenType {
	return enum_TokenType_PLUS
}
func StaticGet_TokenType_MINUS() TokenType {
	return enum_TokenType_MINUS
}
func StaticGet_TokenType_STAR() TokenType {
	return enum_TokenType_STAR
}
func StaticGet_TokenType_STAR_STAR() TokenType {
	return enum_TokenType_STAR_STAR
}
func StaticGet_TokenType_PERCENT() TokenType {
	return enum_TokenType_PERCENT
}
func StaticGet_TokenType_PLUS_PLUS() TokenType {
	return enum_TokenType_PLUS_PLUS
}
func StaticGet_TokenType_MINUS_MINUS() TokenType {
	return enum_TokenType_MINUS_MINUS
}
func StaticGet_TokenType_LEFT_SHIFT() TokenType {
	return enum_TokenType_LEFT_SHIFT
}
func StaticGet_TokenType_RIGHT_SHIFT() TokenType {
	return enum_TokenType_RIGHT_SHIFT
}
func StaticGet_TokenType_UNSIGNED_RIGHT_SHIFT() TokenType {
	return enum_TokenType_UNSIGNED_RIGHT_SHIFT
}
func StaticGet_TokenType_AMPERSAND() TokenType {
	return enum_TokenType_AMPERSAND
}
func StaticGet_TokenType_BAR() TokenType {
	return enum_TokenType_BAR
}
func StaticGet_TokenType_CARET() TokenType {
	return enum_TokenType_CARET
}
func StaticGet_TokenType_BANG() TokenType {
	return enum_TokenType_BANG
}
func StaticGet_TokenType_TILDE() TokenType {
	return enum_TokenType_TILDE
}
func StaticGet_TokenType_AND() TokenType {
	return enum_TokenType_AND
}
func StaticGet_TokenType_OR() TokenType {
	return enum_TokenType_OR
}
func StaticGet_TokenType_QUESTION() TokenType {
	return enum_TokenType_QUESTION
}
func StaticGet_TokenType_QUESTION_QUESTION() TokenType {
	return enum_TokenType_QUESTION_QUESTION
}
func StaticGet_TokenType_QUESTION_DOT() TokenType {
	return enum_TokenType_QUESTION_DOT
}
func StaticGet_TokenType_COLON() TokenType {
	return enum_TokenType_COLON
}
func StaticGet_TokenType_EQUAL() TokenType {
	return enum_TokenType_EQUAL
}
func StaticGet_TokenType_PLUS_EQUAL() TokenType {
	return enum_TokenType_PLUS_EQUAL
}
func StaticGet_TokenType_MINUS_EQUAL() TokenType {
	return enum_TokenType_MINUS_EQUAL
}
func StaticGet_TokenType_STAR_EQUAL() TokenType {
	return enum_TokenType_STAR_EQUAL
}
func StaticGet_TokenType_STAR_STAR_EQUAL() TokenType {
	return enum_TokenType_STAR_STAR_EQUAL
}
func StaticGet_TokenType_PERCENT_EQUAL() TokenType {
	return enum_TokenType_PERCENT_EQUAL
}
func StaticGet_TokenType_LEFT_SHIFT_EQUAL() TokenType {
	return enum_TokenType_LEFT_SHIFT_EQUAL
}
func StaticGet_TokenType_RIGHT_SHIFT_EQUAL() TokenType {
	return enum_TokenType_RIGHT_SHIFT_EQUAL
}
func StaticGet_TokenType_UNSIGNED_RIGHT_SHIFT_EQUAL() TokenType {
	return enum_TokenType_UNSIGNED_RIGHT_SHIFT_EQUAL
}
func StaticGet_TokenType_AMPERSAND_EQUAL() TokenType {
	return enum_TokenType_AMPERSAND_EQUAL
}
func StaticGet_TokenType_BAR_EQUAL() TokenType {
	return enum_TokenType_BAR_EQUAL
}
func StaticGet_TokenType_CARET_EQUAL() TokenType {
	return enum_TokenType_CARET_EQUAL
}
func StaticGet_TokenType_SLASH() TokenType {
	return enum_TokenType_SLASH
}
func StaticGet_TokenType_SLASH_EQUAL() TokenType {
	return enum_TokenType_SLASH_EQUAL
}
func StaticGet_TokenType_AND_EQUAL() TokenType {
	return enum_TokenType_AND_EQUAL
}
func StaticGet_TokenType_OR_EQUAL() TokenType {
	return enum_TokenType_OR_EQUAL
}
func StaticGet_TokenType_QUESTION_QUESTION_EQUAL() TokenType {
	return enum_TokenType_QUESTION_QUESTION_EQUAL
}
func StaticGet_TokenType_NULL() TokenType {
	return enum_TokenType_NULL
}
func StaticGet_TokenType_TRUE() TokenType {
	return enum_TokenType_TRUE
}
func StaticGet_TokenType_FALSE() TokenType {
	return enum_TokenType_FALSE
}
func StaticGet_TokenType_NUMBER() TokenType {
	return enum_TokenType_NUMBER
}
func StaticGet_TokenType_STRING() TokenType {
	return enum_TokenType_STRING
}
func StaticGet_TokenType_BIGINT() TokenType {
	return enum_TokenType_BIGINT
}
func StaticGet_TokenType_REGULAR_EXPRESSION() TokenType {
	return enum_TokenType_REGULAR_EXPRESSION
}
func StaticGet_TokenType_ELLIPSIS() TokenType {
	return enum_TokenType_ELLIPSIS
}
func StaticGet_TokenType_TEMPLATE_HEAD() TokenType {
	return enum_TokenType_TEMPLATE_HEAD
}
func StaticGet_TokenType_TEMPLATE_MIDDLE() TokenType {
	return enum_TokenType_TEMPLATE_MIDDLE
}
func StaticGet_TokenType_TEMPLATE_TAIL() TokenType {
	return enum_TokenType_TEMPLATE_TAIL
}
func StaticGet_TokenType_NO_SUBSTITUTION_TEMPLATE() TokenType {
	return enum_TokenType_NO_SUBSTITUTION_TEMPLATE
}
func StaticGet_TokenType_TYPE() TokenType {
	return enum_TokenType_TYPE
}
func StaticGet_TokenType_DECLARE() TokenType {
	return enum_TokenType_DECLARE
}
func StaticGet_TokenType_MODULE() TokenType {
	return enum_TokenType_MODULE
}
func StaticGet_TokenType_NAMESPACE() TokenType {
	return enum_TokenType_NAMESPACE
}

func TokenType_Values() []TokenType {
	return j2gort.CloneSlice(enumValues_TokenType)
}

func TokenType_ValueOf(name string) TokenType {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_TokenType {
		if enumBaseTable_TokenType[i_].Name_ == name {
			return TokenType(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.parsing.parser.TokenType." + name))
}

// Java method: TokenType.toString()
func (this_ TokenType) String() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	if (this_.FieldGet_value() == j2gort.NullString) {
		return Keywords_Get_1_TokenType(this_).String()
	} else {
		return this_.FieldGet_value()
	}
}

func (this_ TokenType) IsJavaNil_() bool { return this_ == 0 }
func (TokenType) HasIdentityEquals_() {}
func (TokenType) Implements_j2gort_Ext_java_lang_Enum() {}
func (TokenType) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ TokenType) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase { return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()} }
func (TokenType) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.parsing.parser.TokenType", "TokenType" }
func (this_ TokenType) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_TokenType[this_-1]
}
func (this_ TokenType) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ TokenType) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ TokenType) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_TokenType[this_-1].Name_
}
func (this_ TokenType) Equals(other_ any) bool {
	o_, ok_ := other_.(TokenType)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ TokenType) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_TokenType[this_-1].HashCode()
}
func (this_ TokenType) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_TokenType[this_-1].CompareTo(other_)
}
func (this_ TokenType) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_TokenType[this_-1].Name_)
}
func init() {
	vals_ := []any{TokenType(1), TokenType(2), TokenType(3), TokenType(4), TokenType(5), TokenType(6), TokenType(7), TokenType(8), TokenType(9), TokenType(10), TokenType(11), TokenType(12), TokenType(13), TokenType(14), TokenType(15), TokenType(16), TokenType(17), TokenType(18), TokenType(19), TokenType(20), TokenType(21), TokenType(22), TokenType(23), TokenType(24), TokenType(25), TokenType(26), TokenType(27), TokenType(28), TokenType(29), TokenType(30), TokenType(31), TokenType(32), TokenType(33), TokenType(34), TokenType(35), TokenType(36), TokenType(37), TokenType(38), TokenType(39), TokenType(40), TokenType(41), TokenType(42), TokenType(43), TokenType(44), TokenType(45), TokenType(46), TokenType(47), TokenType(48), TokenType(49), TokenType(50), TokenType(51), TokenType(52), TokenType(53), TokenType(54), TokenType(55), TokenType(56), TokenType(57), TokenType(58), TokenType(59), TokenType(60), TokenType(61), TokenType(62), TokenType(63), TokenType(64), TokenType(65), TokenType(66), TokenType(67), TokenType(68), TokenType(69), TokenType(70), TokenType(71), TokenType(72), TokenType(73), TokenType(74), TokenType(75), TokenType(76), TokenType(77), TokenType(78), TokenType(79), TokenType(80), TokenType(81), TokenType(82), TokenType(83), TokenType(84), TokenType(85), TokenType(86), TokenType(87), TokenType(88), TokenType(89), TokenType(90), TokenType(91), TokenType(92), TokenType(93), TokenType(94), TokenType(95), TokenType(96), TokenType(97), TokenType(98), TokenType(99), TokenType(100), TokenType(101), TokenType(102), TokenType(103), TokenType(104), TokenType(105), TokenType(106), TokenType(107), TokenType(108), TokenType(109), TokenType(110), TokenType(111), TokenType(112), TokenType(113), TokenType(114), TokenType(115), TokenType(116), TokenType(117)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.parser.TokenType", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.TokenType", "com.google.javascript.jscomp.parsing.parser.TokenType", func(v_ any) bool {
		return j2gort.InstanceOf[TokenType](v_)
	})
}

