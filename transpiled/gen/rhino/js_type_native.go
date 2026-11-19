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
 *   Bob Jervis
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */
// Port of com/google/javascript/rhino/jstype/JSTypeNative.java (generated mechanically by j2go).

package rhino

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java enum: com.google.javascript.rhino.jstype.JSTypeNative
type JSTypeNative uint8

const (
	enum_JSTypeNative_ARGUMENTS_TYPE JSTypeNative = 1
	enum_JSTypeNative_ARGUMENTS_FUNCTION_TYPE JSTypeNative = 2
	enum_JSTypeNative_ARRAY_TYPE JSTypeNative = 3
	enum_JSTypeNative_ARRAY_FUNCTION_TYPE JSTypeNative = 4
	enum_JSTypeNative_ASYNC_ITERABLE_FUNCTION_TYPE JSTypeNative = 5
	enum_JSTypeNative_ASYNC_ITERABLE_TYPE JSTypeNative = 6
	enum_JSTypeNative_ASYNC_ITERATOR_FUNCTION_TYPE JSTypeNative = 7
	enum_JSTypeNative_ASYNC_ITERATOR_TYPE JSTypeNative = 8
	enum_JSTypeNative_ASYNC_ITERATOR_ITERABLE_FUNCTION_TYPE JSTypeNative = 9
	enum_JSTypeNative_ASYNC_ITERATOR_ITERABLE_TYPE JSTypeNative = 10
	enum_JSTypeNative_ASYNC_GENERATOR_FUNCTION_TYPE JSTypeNative = 11
	enum_JSTypeNative_ASYNC_GENERATOR_TYPE JSTypeNative = 12
	enum_JSTypeNative_BIGINT_TYPE JSTypeNative = 13
	enum_JSTypeNative_BIGINT_OBJECT_TYPE JSTypeNative = 14
	enum_JSTypeNative_BIGINT_OBJECT_FUNCTION_TYPE JSTypeNative = 15
	enum_JSTypeNative_BOOLEAN_TYPE JSTypeNative = 16
	enum_JSTypeNative_BOOLEAN_OBJECT_TYPE JSTypeNative = 17
	enum_JSTypeNative_BOOLEAN_OBJECT_FUNCTION_TYPE JSTypeNative = 18
	enum_JSTypeNative_CHECKED_UNKNOWN_TYPE JSTypeNative = 19
	enum_JSTypeNative_DATE_TYPE JSTypeNative = 20
	enum_JSTypeNative_DATE_FUNCTION_TYPE JSTypeNative = 21
	enum_JSTypeNative_FUNCTION_TYPE JSTypeNative = 22
	enum_JSTypeNative_FUNCTION_FUNCTION_TYPE JSTypeNative = 23
	enum_JSTypeNative_FUNCTION_PROTOTYPE JSTypeNative = 24
	enum_JSTypeNative_FUNCTION_INSTANCE_PROTOTYPE JSTypeNative = 25
	enum_JSTypeNative_GBIGINT_TYPE JSTypeNative = 26
	enum_JSTypeNative_GENERATOR_FUNCTION_TYPE JSTypeNative = 27
	enum_JSTypeNative_GENERATOR_TYPE JSTypeNative = 28
	enum_JSTypeNative_I_ITERABLE_RESULT_FUNCTION_TYPE JSTypeNative = 29
	enum_JSTypeNative_I_ITERABLE_RESULT_TYPE JSTypeNative = 30
	enum_JSTypeNative_ITERABLE_FUNCTION_TYPE JSTypeNative = 31
	enum_JSTypeNative_ITERABLE_TYPE JSTypeNative = 32
	enum_JSTypeNative_ITERATOR_LIKE_FUNCTION_TYPE JSTypeNative = 33
	enum_JSTypeNative_ITERATOR_LIKE_TYPE JSTypeNative = 34
	enum_JSTypeNative_ITERATOR_FUNCTION_TYPE JSTypeNative = 35
	enum_JSTypeNative_ITERATOR_TYPE JSTypeNative = 36
	enum_JSTypeNative_ITERATOR_ITERABLE_FUNCTION_TYPE JSTypeNative = 37
	enum_JSTypeNative_ITERATOR_ITERABLE_TYPE JSTypeNative = 38
	enum_JSTypeNative_I_ARRAY_LIKE_FUNCTION_TYPE JSTypeNative = 39
	enum_JSTypeNative_I_ARRAY_LIKE_TYPE JSTypeNative = 40
	enum_JSTypeNative_I_TEMPLATE_ARRAY_TYPE JSTypeNative = 41
	enum_JSTypeNative_I_OBJECT_FUNCTION_TYPE JSTypeNative = 42
	enum_JSTypeNative_I_OBJECT_TYPE JSTypeNative = 43
	enum_JSTypeNative_I_THENABLE_FUNCTION_TYPE JSTypeNative = 44
	enum_JSTypeNative_I_THENABLE_TYPE JSTypeNative = 45
	enum_JSTypeNative_NULL_TYPE JSTypeNative = 46
	enum_JSTypeNative_NUMBER_TYPE JSTypeNative = 47
	enum_JSTypeNative_NUMBER_OBJECT_TYPE JSTypeNative = 48
	enum_JSTypeNative_NUMBER_OBJECT_FUNCTION_TYPE JSTypeNative = 49
	enum_JSTypeNative_PROMISE_TYPE JSTypeNative = 50
	enum_JSTypeNative_PROMISE_FUNCTION_TYPE JSTypeNative = 51
	enum_JSTypeNative_OBJECT_TYPE JSTypeNative = 52
	enum_JSTypeNative_OBJECT_FUNCTION_TYPE JSTypeNative = 53
	enum_JSTypeNative_OBJECT_PROTOTYPE JSTypeNative = 54
	enum_JSTypeNative_READONLY_ARRAY_TYPE JSTypeNative = 55
	enum_JSTypeNative_READONLY_ARRAY_FUNCTION_TYPE JSTypeNative = 56
	enum_JSTypeNative_READONLY_MAP_TYPE JSTypeNative = 57
	enum_JSTypeNative_READONLY_MAP_FUNCTION_TYPE JSTypeNative = 58
	enum_JSTypeNative_MAP_TYPE JSTypeNative = 59
	enum_JSTypeNative_MAP_FUNCTION_TYPE JSTypeNative = 60
	enum_JSTypeNative_REGEXP_TYPE JSTypeNative = 61
	enum_JSTypeNative_REGEXP_FUNCTION_TYPE JSTypeNative = 62
	enum_JSTypeNative_STRING_OBJECT_TYPE JSTypeNative = 63
	enum_JSTypeNative_STRING_OBJECT_FUNCTION_TYPE JSTypeNative = 64
	enum_JSTypeNative_STRING_TYPE JSTypeNative = 65
	enum_JSTypeNative_SYMBOL_OBJECT_TYPE JSTypeNative = 66
	enum_JSTypeNative_SYMBOL_OBJECT_FUNCTION_TYPE JSTypeNative = 67
	enum_JSTypeNative_SYMBOL_TYPE JSTypeNative = 68
	enum_JSTypeNative_THENABLE_TYPE JSTypeNative = 69
	enum_JSTypeNative_UNKNOWN_TYPE JSTypeNative = 70
	enum_JSTypeNative_VOID_TYPE JSTypeNative = 71
	enum_JSTypeNative_ALL_TYPE JSTypeNative = 72
	enum_JSTypeNative_NO_TYPE JSTypeNative = 73
	enum_JSTypeNative_NO_OBJECT_TYPE JSTypeNative = 74
	enum_JSTypeNative_GLOBAL_THIS JSTypeNative = 75
	enum_JSTypeNative_LEAST_FUNCTION_TYPE JSTypeNative = 76
	enum_JSTypeNative_GREATEST_FUNCTION_TYPE JSTypeNative = 77
	enum_JSTypeNative_NULL_VOID JSTypeNative = 78
	enum_JSTypeNative_NUMBER_STRING_BOOLEAN JSTypeNative = 79
	enum_JSTypeNative_VALUE_TYPES JSTypeNative = 80
	enum_JSTypeNative_NUMBER_SYMBOL JSTypeNative = 81
	enum_JSTypeNative_STRING_SYMBOL JSTypeNative = 82
	enum_JSTypeNative_NUMBER_STRING JSTypeNative = 83
	enum_JSTypeNative_NUMBER_STRING_SYMBOL JSTypeNative = 84
	enum_JSTypeNative_BIGINT_NUMBER JSTypeNative = 85
	enum_JSTypeNative_BIGINT_NUMBER_OBJECT JSTypeNative = 86
	enum_JSTypeNative_BIGINT_NUMBER_STRING JSTypeNative = 87
	enum_JSTypeNative_BIGINT_NUMBER_STRING_OBJECT JSTypeNative = 88
	enum_JSTypeNative_NUMBER_ADDITION_SUPERTYPE JSTypeNative = 89
)

var enumBaseTable_JSTypeNative = [...]j2gort.EnumBase{
	{Name_: "ARGUMENTS_TYPE", Ordinal_: 0},
	{Name_: "ARGUMENTS_FUNCTION_TYPE", Ordinal_: 1},
	{Name_: "ARRAY_TYPE", Ordinal_: 2},
	{Name_: "ARRAY_FUNCTION_TYPE", Ordinal_: 3},
	{Name_: "ASYNC_ITERABLE_FUNCTION_TYPE", Ordinal_: 4},
	{Name_: "ASYNC_ITERABLE_TYPE", Ordinal_: 5},
	{Name_: "ASYNC_ITERATOR_FUNCTION_TYPE", Ordinal_: 6},
	{Name_: "ASYNC_ITERATOR_TYPE", Ordinal_: 7},
	{Name_: "ASYNC_ITERATOR_ITERABLE_FUNCTION_TYPE", Ordinal_: 8},
	{Name_: "ASYNC_ITERATOR_ITERABLE_TYPE", Ordinal_: 9},
	{Name_: "ASYNC_GENERATOR_FUNCTION_TYPE", Ordinal_: 10},
	{Name_: "ASYNC_GENERATOR_TYPE", Ordinal_: 11},
	{Name_: "BIGINT_TYPE", Ordinal_: 12},
	{Name_: "BIGINT_OBJECT_TYPE", Ordinal_: 13},
	{Name_: "BIGINT_OBJECT_FUNCTION_TYPE", Ordinal_: 14},
	{Name_: "BOOLEAN_TYPE", Ordinal_: 15},
	{Name_: "BOOLEAN_OBJECT_TYPE", Ordinal_: 16},
	{Name_: "BOOLEAN_OBJECT_FUNCTION_TYPE", Ordinal_: 17},
	{Name_: "CHECKED_UNKNOWN_TYPE", Ordinal_: 18},
	{Name_: "DATE_TYPE", Ordinal_: 19},
	{Name_: "DATE_FUNCTION_TYPE", Ordinal_: 20},
	{Name_: "FUNCTION_TYPE", Ordinal_: 21},
	{Name_: "FUNCTION_FUNCTION_TYPE", Ordinal_: 22},
	{Name_: "FUNCTION_PROTOTYPE", Ordinal_: 23},
	{Name_: "FUNCTION_INSTANCE_PROTOTYPE", Ordinal_: 24},
	{Name_: "GBIGINT_TYPE", Ordinal_: 25},
	{Name_: "GENERATOR_FUNCTION_TYPE", Ordinal_: 26},
	{Name_: "GENERATOR_TYPE", Ordinal_: 27},
	{Name_: "I_ITERABLE_RESULT_FUNCTION_TYPE", Ordinal_: 28},
	{Name_: "I_ITERABLE_RESULT_TYPE", Ordinal_: 29},
	{Name_: "ITERABLE_FUNCTION_TYPE", Ordinal_: 30},
	{Name_: "ITERABLE_TYPE", Ordinal_: 31},
	{Name_: "ITERATOR_LIKE_FUNCTION_TYPE", Ordinal_: 32},
	{Name_: "ITERATOR_LIKE_TYPE", Ordinal_: 33},
	{Name_: "ITERATOR_FUNCTION_TYPE", Ordinal_: 34},
	{Name_: "ITERATOR_TYPE", Ordinal_: 35},
	{Name_: "ITERATOR_ITERABLE_FUNCTION_TYPE", Ordinal_: 36},
	{Name_: "ITERATOR_ITERABLE_TYPE", Ordinal_: 37},
	{Name_: "I_ARRAY_LIKE_FUNCTION_TYPE", Ordinal_: 38},
	{Name_: "I_ARRAY_LIKE_TYPE", Ordinal_: 39},
	{Name_: "I_TEMPLATE_ARRAY_TYPE", Ordinal_: 40},
	{Name_: "I_OBJECT_FUNCTION_TYPE", Ordinal_: 41},
	{Name_: "I_OBJECT_TYPE", Ordinal_: 42},
	{Name_: "I_THENABLE_FUNCTION_TYPE", Ordinal_: 43},
	{Name_: "I_THENABLE_TYPE", Ordinal_: 44},
	{Name_: "NULL_TYPE", Ordinal_: 45},
	{Name_: "NUMBER_TYPE", Ordinal_: 46},
	{Name_: "NUMBER_OBJECT_TYPE", Ordinal_: 47},
	{Name_: "NUMBER_OBJECT_FUNCTION_TYPE", Ordinal_: 48},
	{Name_: "PROMISE_TYPE", Ordinal_: 49},
	{Name_: "PROMISE_FUNCTION_TYPE", Ordinal_: 50},
	{Name_: "OBJECT_TYPE", Ordinal_: 51},
	{Name_: "OBJECT_FUNCTION_TYPE", Ordinal_: 52},
	{Name_: "OBJECT_PROTOTYPE", Ordinal_: 53},
	{Name_: "READONLY_ARRAY_TYPE", Ordinal_: 54},
	{Name_: "READONLY_ARRAY_FUNCTION_TYPE", Ordinal_: 55},
	{Name_: "READONLY_MAP_TYPE", Ordinal_: 56},
	{Name_: "READONLY_MAP_FUNCTION_TYPE", Ordinal_: 57},
	{Name_: "MAP_TYPE", Ordinal_: 58},
	{Name_: "MAP_FUNCTION_TYPE", Ordinal_: 59},
	{Name_: "REGEXP_TYPE", Ordinal_: 60},
	{Name_: "REGEXP_FUNCTION_TYPE", Ordinal_: 61},
	{Name_: "STRING_OBJECT_TYPE", Ordinal_: 62},
	{Name_: "STRING_OBJECT_FUNCTION_TYPE", Ordinal_: 63},
	{Name_: "STRING_TYPE", Ordinal_: 64},
	{Name_: "SYMBOL_OBJECT_TYPE", Ordinal_: 65},
	{Name_: "SYMBOL_OBJECT_FUNCTION_TYPE", Ordinal_: 66},
	{Name_: "SYMBOL_TYPE", Ordinal_: 67},
	{Name_: "THENABLE_TYPE", Ordinal_: 68},
	{Name_: "UNKNOWN_TYPE", Ordinal_: 69},
	{Name_: "VOID_TYPE", Ordinal_: 70},
	{Name_: "ALL_TYPE", Ordinal_: 71},
	{Name_: "NO_TYPE", Ordinal_: 72},
	{Name_: "NO_OBJECT_TYPE", Ordinal_: 73},
	{Name_: "GLOBAL_THIS", Ordinal_: 74},
	{Name_: "LEAST_FUNCTION_TYPE", Ordinal_: 75},
	{Name_: "GREATEST_FUNCTION_TYPE", Ordinal_: 76},
	{Name_: "NULL_VOID", Ordinal_: 77},
	{Name_: "NUMBER_STRING_BOOLEAN", Ordinal_: 78},
	{Name_: "VALUE_TYPES", Ordinal_: 79},
	{Name_: "NUMBER_SYMBOL", Ordinal_: 80},
	{Name_: "STRING_SYMBOL", Ordinal_: 81},
	{Name_: "NUMBER_STRING", Ordinal_: 82},
	{Name_: "NUMBER_STRING_SYMBOL", Ordinal_: 83},
	{Name_: "BIGINT_NUMBER", Ordinal_: 84},
	{Name_: "BIGINT_NUMBER_OBJECT", Ordinal_: 85},
	{Name_: "BIGINT_NUMBER_STRING", Ordinal_: 86},
	{Name_: "BIGINT_NUMBER_STRING_OBJECT", Ordinal_: 87},
	{Name_: "NUMBER_ADDITION_SUPERTYPE", Ordinal_: 88},
}
var enumValues_JSTypeNative = []JSTypeNative{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89}

func StaticGet_JSTypeNative_ARGUMENTS_TYPE() JSTypeNative {
	return enum_JSTypeNative_ARGUMENTS_TYPE
}
func StaticGet_JSTypeNative_ARGUMENTS_FUNCTION_TYPE() JSTypeNative {
	return enum_JSTypeNative_ARGUMENTS_FUNCTION_TYPE
}
func StaticGet_JSTypeNative_ARRAY_TYPE() JSTypeNative {
	return enum_JSTypeNative_ARRAY_TYPE
}
func StaticGet_JSTypeNative_ARRAY_FUNCTION_TYPE() JSTypeNative {
	return enum_JSTypeNative_ARRAY_FUNCTION_TYPE
}
func StaticGet_JSTypeNative_ASYNC_ITERABLE_FUNCTION_TYPE() JSTypeNative {
	return enum_JSTypeNative_ASYNC_ITERABLE_FUNCTION_TYPE
}
func StaticGet_JSTypeNative_ASYNC_ITERABLE_TYPE() JSTypeNative {
	return enum_JSTypeNative_ASYNC_ITERABLE_TYPE
}
func StaticGet_JSTypeNative_ASYNC_ITERATOR_FUNCTION_TYPE() JSTypeNative {
	return enum_JSTypeNative_ASYNC_ITERATOR_FUNCTION_TYPE
}
func StaticGet_JSTypeNative_ASYNC_ITERATOR_TYPE() JSTypeNative {
	return enum_JSTypeNative_ASYNC_ITERATOR_TYPE
}
func StaticGet_JSTypeNative_ASYNC_ITERATOR_ITERABLE_FUNCTION_TYPE() JSTypeNative {
	return enum_JSTypeNative_ASYNC_ITERATOR_ITERABLE_FUNCTION_TYPE
}
func StaticGet_JSTypeNative_ASYNC_ITERATOR_ITERABLE_TYPE() JSTypeNative {
	return enum_JSTypeNative_ASYNC_ITERATOR_ITERABLE_TYPE
}
func StaticGet_JSTypeNative_ASYNC_GENERATOR_FUNCTION_TYPE() JSTypeNative {
	return enum_JSTypeNative_ASYNC_GENERATOR_FUNCTION_TYPE
}
func StaticGet_JSTypeNative_ASYNC_GENERATOR_TYPE() JSTypeNative {
	return enum_JSTypeNative_ASYNC_GENERATOR_TYPE
}
func StaticGet_JSTypeNative_BIGINT_TYPE() JSTypeNative {
	return enum_JSTypeNative_BIGINT_TYPE
}
func StaticGet_JSTypeNative_BIGINT_OBJECT_TYPE() JSTypeNative {
	return enum_JSTypeNative_BIGINT_OBJECT_TYPE
}
func StaticGet_JSTypeNative_BIGINT_OBJECT_FUNCTION_TYPE() JSTypeNative {
	return enum_JSTypeNative_BIGINT_OBJECT_FUNCTION_TYPE
}
func StaticGet_JSTypeNative_BOOLEAN_TYPE() JSTypeNative {
	return enum_JSTypeNative_BOOLEAN_TYPE
}
func StaticGet_JSTypeNative_BOOLEAN_OBJECT_TYPE() JSTypeNative {
	return enum_JSTypeNative_BOOLEAN_OBJECT_TYPE
}
func StaticGet_JSTypeNative_BOOLEAN_OBJECT_FUNCTION_TYPE() JSTypeNative {
	return enum_JSTypeNative_BOOLEAN_OBJECT_FUNCTION_TYPE
}
func StaticGet_JSTypeNative_CHECKED_UNKNOWN_TYPE() JSTypeNative {
	return enum_JSTypeNative_CHECKED_UNKNOWN_TYPE
}
func StaticGet_JSTypeNative_DATE_TYPE() JSTypeNative {
	return enum_JSTypeNative_DATE_TYPE
}
func StaticGet_JSTypeNative_DATE_FUNCTION_TYPE() JSTypeNative {
	return enum_JSTypeNative_DATE_FUNCTION_TYPE
}
func StaticGet_JSTypeNative_FUNCTION_TYPE() JSTypeNative {
	return enum_JSTypeNative_FUNCTION_TYPE
}
func StaticGet_JSTypeNative_FUNCTION_FUNCTION_TYPE() JSTypeNative {
	return enum_JSTypeNative_FUNCTION_FUNCTION_TYPE
}
func StaticGet_JSTypeNative_FUNCTION_PROTOTYPE() JSTypeNative {
	return enum_JSTypeNative_FUNCTION_PROTOTYPE
}
func StaticGet_JSTypeNative_FUNCTION_INSTANCE_PROTOTYPE() JSTypeNative {
	return enum_JSTypeNative_FUNCTION_INSTANCE_PROTOTYPE
}
func StaticGet_JSTypeNative_GBIGINT_TYPE() JSTypeNative {
	return enum_JSTypeNative_GBIGINT_TYPE
}
func StaticGet_JSTypeNative_GENERATOR_FUNCTION_TYPE() JSTypeNative {
	return enum_JSTypeNative_GENERATOR_FUNCTION_TYPE
}
func StaticGet_JSTypeNative_GENERATOR_TYPE() JSTypeNative {
	return enum_JSTypeNative_GENERATOR_TYPE
}
func StaticGet_JSTypeNative_I_ITERABLE_RESULT_FUNCTION_TYPE() JSTypeNative {
	return enum_JSTypeNative_I_ITERABLE_RESULT_FUNCTION_TYPE
}
func StaticGet_JSTypeNative_I_ITERABLE_RESULT_TYPE() JSTypeNative {
	return enum_JSTypeNative_I_ITERABLE_RESULT_TYPE
}
func StaticGet_JSTypeNative_ITERABLE_FUNCTION_TYPE() JSTypeNative {
	return enum_JSTypeNative_ITERABLE_FUNCTION_TYPE
}
func StaticGet_JSTypeNative_ITERABLE_TYPE() JSTypeNative {
	return enum_JSTypeNative_ITERABLE_TYPE
}
func StaticGet_JSTypeNative_ITERATOR_LIKE_FUNCTION_TYPE() JSTypeNative {
	return enum_JSTypeNative_ITERATOR_LIKE_FUNCTION_TYPE
}
func StaticGet_JSTypeNative_ITERATOR_LIKE_TYPE() JSTypeNative {
	return enum_JSTypeNative_ITERATOR_LIKE_TYPE
}
func StaticGet_JSTypeNative_ITERATOR_FUNCTION_TYPE() JSTypeNative {
	return enum_JSTypeNative_ITERATOR_FUNCTION_TYPE
}
func StaticGet_JSTypeNative_ITERATOR_TYPE() JSTypeNative {
	return enum_JSTypeNative_ITERATOR_TYPE
}
func StaticGet_JSTypeNative_ITERATOR_ITERABLE_FUNCTION_TYPE() JSTypeNative {
	return enum_JSTypeNative_ITERATOR_ITERABLE_FUNCTION_TYPE
}
func StaticGet_JSTypeNative_ITERATOR_ITERABLE_TYPE() JSTypeNative {
	return enum_JSTypeNative_ITERATOR_ITERABLE_TYPE
}
func StaticGet_JSTypeNative_I_ARRAY_LIKE_FUNCTION_TYPE() JSTypeNative {
	return enum_JSTypeNative_I_ARRAY_LIKE_FUNCTION_TYPE
}
func StaticGet_JSTypeNative_I_ARRAY_LIKE_TYPE() JSTypeNative {
	return enum_JSTypeNative_I_ARRAY_LIKE_TYPE
}
func StaticGet_JSTypeNative_I_TEMPLATE_ARRAY_TYPE() JSTypeNative {
	return enum_JSTypeNative_I_TEMPLATE_ARRAY_TYPE
}
func StaticGet_JSTypeNative_I_OBJECT_FUNCTION_TYPE() JSTypeNative {
	return enum_JSTypeNative_I_OBJECT_FUNCTION_TYPE
}
func StaticGet_JSTypeNative_I_OBJECT_TYPE() JSTypeNative {
	return enum_JSTypeNative_I_OBJECT_TYPE
}
func StaticGet_JSTypeNative_I_THENABLE_FUNCTION_TYPE() JSTypeNative {
	return enum_JSTypeNative_I_THENABLE_FUNCTION_TYPE
}
func StaticGet_JSTypeNative_I_THENABLE_TYPE() JSTypeNative {
	return enum_JSTypeNative_I_THENABLE_TYPE
}
func StaticGet_JSTypeNative_NULL_TYPE() JSTypeNative {
	return enum_JSTypeNative_NULL_TYPE
}
func StaticGet_JSTypeNative_NUMBER_TYPE() JSTypeNative {
	return enum_JSTypeNative_NUMBER_TYPE
}
func StaticGet_JSTypeNative_NUMBER_OBJECT_TYPE() JSTypeNative {
	return enum_JSTypeNative_NUMBER_OBJECT_TYPE
}
func StaticGet_JSTypeNative_NUMBER_OBJECT_FUNCTION_TYPE() JSTypeNative {
	return enum_JSTypeNative_NUMBER_OBJECT_FUNCTION_TYPE
}
func StaticGet_JSTypeNative_PROMISE_TYPE() JSTypeNative {
	return enum_JSTypeNative_PROMISE_TYPE
}
func StaticGet_JSTypeNative_PROMISE_FUNCTION_TYPE() JSTypeNative {
	return enum_JSTypeNative_PROMISE_FUNCTION_TYPE
}
func StaticGet_JSTypeNative_OBJECT_TYPE() JSTypeNative {
	return enum_JSTypeNative_OBJECT_TYPE
}
func StaticGet_JSTypeNative_OBJECT_FUNCTION_TYPE() JSTypeNative {
	return enum_JSTypeNative_OBJECT_FUNCTION_TYPE
}
func StaticGet_JSTypeNative_OBJECT_PROTOTYPE() JSTypeNative {
	return enum_JSTypeNative_OBJECT_PROTOTYPE
}
func StaticGet_JSTypeNative_READONLY_ARRAY_TYPE() JSTypeNative {
	return enum_JSTypeNative_READONLY_ARRAY_TYPE
}
func StaticGet_JSTypeNative_READONLY_ARRAY_FUNCTION_TYPE() JSTypeNative {
	return enum_JSTypeNative_READONLY_ARRAY_FUNCTION_TYPE
}
func StaticGet_JSTypeNative_READONLY_MAP_TYPE() JSTypeNative {
	return enum_JSTypeNative_READONLY_MAP_TYPE
}
func StaticGet_JSTypeNative_READONLY_MAP_FUNCTION_TYPE() JSTypeNative {
	return enum_JSTypeNative_READONLY_MAP_FUNCTION_TYPE
}
func StaticGet_JSTypeNative_MAP_TYPE() JSTypeNative {
	return enum_JSTypeNative_MAP_TYPE
}
func StaticGet_JSTypeNative_MAP_FUNCTION_TYPE() JSTypeNative {
	return enum_JSTypeNative_MAP_FUNCTION_TYPE
}
func StaticGet_JSTypeNative_REGEXP_TYPE() JSTypeNative {
	return enum_JSTypeNative_REGEXP_TYPE
}
func StaticGet_JSTypeNative_REGEXP_FUNCTION_TYPE() JSTypeNative {
	return enum_JSTypeNative_REGEXP_FUNCTION_TYPE
}
func StaticGet_JSTypeNative_STRING_OBJECT_TYPE() JSTypeNative {
	return enum_JSTypeNative_STRING_OBJECT_TYPE
}
func StaticGet_JSTypeNative_STRING_OBJECT_FUNCTION_TYPE() JSTypeNative {
	return enum_JSTypeNative_STRING_OBJECT_FUNCTION_TYPE
}
func StaticGet_JSTypeNative_STRING_TYPE() JSTypeNative {
	return enum_JSTypeNative_STRING_TYPE
}
func StaticGet_JSTypeNative_SYMBOL_OBJECT_TYPE() JSTypeNative {
	return enum_JSTypeNative_SYMBOL_OBJECT_TYPE
}
func StaticGet_JSTypeNative_SYMBOL_OBJECT_FUNCTION_TYPE() JSTypeNative {
	return enum_JSTypeNative_SYMBOL_OBJECT_FUNCTION_TYPE
}
func StaticGet_JSTypeNative_SYMBOL_TYPE() JSTypeNative {
	return enum_JSTypeNative_SYMBOL_TYPE
}
func StaticGet_JSTypeNative_THENABLE_TYPE() JSTypeNative {
	return enum_JSTypeNative_THENABLE_TYPE
}
func StaticGet_JSTypeNative_UNKNOWN_TYPE() JSTypeNative {
	return enum_JSTypeNative_UNKNOWN_TYPE
}
func StaticGet_JSTypeNative_VOID_TYPE() JSTypeNative {
	return enum_JSTypeNative_VOID_TYPE
}
func StaticGet_JSTypeNative_ALL_TYPE() JSTypeNative {
	return enum_JSTypeNative_ALL_TYPE
}
func StaticGet_JSTypeNative_NO_TYPE() JSTypeNative {
	return enum_JSTypeNative_NO_TYPE
}
func StaticGet_JSTypeNative_NO_OBJECT_TYPE() JSTypeNative {
	return enum_JSTypeNative_NO_OBJECT_TYPE
}
func StaticGet_JSTypeNative_GLOBAL_THIS() JSTypeNative {
	return enum_JSTypeNative_GLOBAL_THIS
}
func StaticGet_JSTypeNative_LEAST_FUNCTION_TYPE() JSTypeNative {
	return enum_JSTypeNative_LEAST_FUNCTION_TYPE
}
func StaticGet_JSTypeNative_GREATEST_FUNCTION_TYPE() JSTypeNative {
	return enum_JSTypeNative_GREATEST_FUNCTION_TYPE
}
func StaticGet_JSTypeNative_NULL_VOID() JSTypeNative {
	return enum_JSTypeNative_NULL_VOID
}
func StaticGet_JSTypeNative_NUMBER_STRING_BOOLEAN() JSTypeNative {
	return enum_JSTypeNative_NUMBER_STRING_BOOLEAN
}
func StaticGet_JSTypeNative_VALUE_TYPES() JSTypeNative {
	return enum_JSTypeNative_VALUE_TYPES
}
func StaticGet_JSTypeNative_NUMBER_SYMBOL() JSTypeNative {
	return enum_JSTypeNative_NUMBER_SYMBOL
}
func StaticGet_JSTypeNative_STRING_SYMBOL() JSTypeNative {
	return enum_JSTypeNative_STRING_SYMBOL
}
func StaticGet_JSTypeNative_NUMBER_STRING() JSTypeNative {
	return enum_JSTypeNative_NUMBER_STRING
}
func StaticGet_JSTypeNative_NUMBER_STRING_SYMBOL() JSTypeNative {
	return enum_JSTypeNative_NUMBER_STRING_SYMBOL
}
func StaticGet_JSTypeNative_BIGINT_NUMBER() JSTypeNative {
	return enum_JSTypeNative_BIGINT_NUMBER
}
func StaticGet_JSTypeNative_BIGINT_NUMBER_OBJECT() JSTypeNative {
	return enum_JSTypeNative_BIGINT_NUMBER_OBJECT
}
func StaticGet_JSTypeNative_BIGINT_NUMBER_STRING() JSTypeNative {
	return enum_JSTypeNative_BIGINT_NUMBER_STRING
}
func StaticGet_JSTypeNative_BIGINT_NUMBER_STRING_OBJECT() JSTypeNative {
	return enum_JSTypeNative_BIGINT_NUMBER_STRING_OBJECT
}
func StaticGet_JSTypeNative_NUMBER_ADDITION_SUPERTYPE() JSTypeNative {
	return enum_JSTypeNative_NUMBER_ADDITION_SUPERTYPE
}

func JSTypeNative_Values() []JSTypeNative {
	return j2gort.CloneSlice(enumValues_JSTypeNative)
}

func JSTypeNative_ValueOf(name string) JSTypeNative {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_JSTypeNative {
		if enumBaseTable_JSTypeNative[i_].Name_ == name {
			return JSTypeNative(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.rhino.jstype.JSTypeNative." + name))
}

func (this_ JSTypeNative) IsJavaNil_() bool { return this_ == 0 }
func (JSTypeNative) HasIdentityEquals_() {}
func (JSTypeNative) Implements_j2gort_Ext_java_lang_Enum() {}
func (JSTypeNative) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ JSTypeNative) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase { return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()} }
func (JSTypeNative) JavaClassInfo_() (string, string) { return "com.google.javascript.rhino.jstype.JSTypeNative", "JSTypeNative" }
func (this_ JSTypeNative) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_JSTypeNative[this_-1]
}
func (this_ JSTypeNative) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ JSTypeNative) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ JSTypeNative) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_JSTypeNative[this_-1].Name_
}
func (this_ JSTypeNative) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_JSTypeNative[this_-1].Name_
}
func (this_ JSTypeNative) Equals(other_ any) bool {
	o_, ok_ := other_.(JSTypeNative)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ JSTypeNative) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_JSTypeNative[this_-1].HashCode()
}
func (this_ JSTypeNative) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_JSTypeNative[this_-1].CompareTo(other_)
}
func (this_ JSTypeNative) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_JSTypeNative[this_-1].Name_)
}
func init() {
	vals_ := []any{JSTypeNative(1), JSTypeNative(2), JSTypeNative(3), JSTypeNative(4), JSTypeNative(5), JSTypeNative(6), JSTypeNative(7), JSTypeNative(8), JSTypeNative(9), JSTypeNative(10), JSTypeNative(11), JSTypeNative(12), JSTypeNative(13), JSTypeNative(14), JSTypeNative(15), JSTypeNative(16), JSTypeNative(17), JSTypeNative(18), JSTypeNative(19), JSTypeNative(20), JSTypeNative(21), JSTypeNative(22), JSTypeNative(23), JSTypeNative(24), JSTypeNative(25), JSTypeNative(26), JSTypeNative(27), JSTypeNative(28), JSTypeNative(29), JSTypeNative(30), JSTypeNative(31), JSTypeNative(32), JSTypeNative(33), JSTypeNative(34), JSTypeNative(35), JSTypeNative(36), JSTypeNative(37), JSTypeNative(38), JSTypeNative(39), JSTypeNative(40), JSTypeNative(41), JSTypeNative(42), JSTypeNative(43), JSTypeNative(44), JSTypeNative(45), JSTypeNative(46), JSTypeNative(47), JSTypeNative(48), JSTypeNative(49), JSTypeNative(50), JSTypeNative(51), JSTypeNative(52), JSTypeNative(53), JSTypeNative(54), JSTypeNative(55), JSTypeNative(56), JSTypeNative(57), JSTypeNative(58), JSTypeNative(59), JSTypeNative(60), JSTypeNative(61), JSTypeNative(62), JSTypeNative(63), JSTypeNative(64), JSTypeNative(65), JSTypeNative(66), JSTypeNative(67), JSTypeNative(68), JSTypeNative(69), JSTypeNative(70), JSTypeNative(71), JSTypeNative(72), JSTypeNative(73), JSTypeNative(74), JSTypeNative(75), JSTypeNative(76), JSTypeNative(77), JSTypeNative(78), JSTypeNative(79), JSTypeNative(80), JSTypeNative(81), JSTypeNative(82), JSTypeNative(83), JSTypeNative(84), JSTypeNative(85), JSTypeNative(86), JSTypeNative(87), JSTypeNative(88), JSTypeNative(89)}
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.jstype.JSTypeNative", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.JSTypeNative", "com.google.javascript.rhino.jstype.JSTypeNative", func(v_ any) bool {
		return j2gort.InstanceOf[JSTypeNative](v_)
	})
}

