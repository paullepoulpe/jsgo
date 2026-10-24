// Port of com/google/javascript/rhino/jstype/JSTypeNative.java @ closure-compiler v20261021

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

package rhino

// JSTypeNative is the set of constants corresponding to types that are
// built into a JavaScript engine and other types that occur very often in
// the type system. See JSTypeRegistry#getNativeType(JSTypeNative).
//
// CHECKED_UNKNOWN_TYPE: a checked unknown type is a type that we know
// something about, but we're not really sure what we know about it.
// Examples of checked unknown types include
//
//	if (x) { // x is unknown
//	  alert(x); // x is checked unknown
//	}
//
// and a parameter of a forward-declared type that has not been included in
// this binary. This is useful for missing property warnings, where we don't
// want to emit warnings on things that have been checked.
type JSTypeNative int32

const (
	JSTypeNative_ARGUMENTS_TYPE JSTypeNative = iota
	JSTypeNative_ARGUMENTS_FUNCTION_TYPE
	JSTypeNative_ARRAY_TYPE
	JSTypeNative_ARRAY_FUNCTION_TYPE
	JSTypeNative_ASYNC_ITERABLE_FUNCTION_TYPE
	JSTypeNative_ASYNC_ITERABLE_TYPE
	JSTypeNative_ASYNC_ITERATOR_FUNCTION_TYPE
	JSTypeNative_ASYNC_ITERATOR_TYPE
	JSTypeNative_ASYNC_ITERATOR_ITERABLE_FUNCTION_TYPE
	JSTypeNative_ASYNC_ITERATOR_ITERABLE_TYPE
	JSTypeNative_ASYNC_GENERATOR_FUNCTION_TYPE
	JSTypeNative_ASYNC_GENERATOR_TYPE
	JSTypeNative_BIGINT_TYPE
	JSTypeNative_BIGINT_OBJECT_TYPE
	JSTypeNative_BIGINT_OBJECT_FUNCTION_TYPE
	JSTypeNative_BOOLEAN_TYPE
	JSTypeNative_BOOLEAN_OBJECT_TYPE
	JSTypeNative_BOOLEAN_OBJECT_FUNCTION_TYPE
	JSTypeNative_CHECKED_UNKNOWN_TYPE
	JSTypeNative_DATE_TYPE
	JSTypeNative_DATE_FUNCTION_TYPE
	JSTypeNative_FUNCTION_TYPE
	JSTypeNative_FUNCTION_FUNCTION_TYPE
	JSTypeNative_FUNCTION_PROTOTYPE
	JSTypeNative_FUNCTION_INSTANCE_PROTOTYPE
	JSTypeNative_GBIGINT_TYPE
	JSTypeNative_GENERATOR_FUNCTION_TYPE
	JSTypeNative_GENERATOR_TYPE
	JSTypeNative_I_ITERABLE_RESULT_FUNCTION_TYPE
	JSTypeNative_I_ITERABLE_RESULT_TYPE
	JSTypeNative_ITERABLE_FUNCTION_TYPE
	JSTypeNative_ITERABLE_TYPE
	JSTypeNative_ITERATOR_LIKE_FUNCTION_TYPE
	JSTypeNative_ITERATOR_LIKE_TYPE
	JSTypeNative_ITERATOR_FUNCTION_TYPE
	JSTypeNative_ITERATOR_TYPE
	JSTypeNative_ITERATOR_ITERABLE_FUNCTION_TYPE
	JSTypeNative_ITERATOR_ITERABLE_TYPE
	JSTypeNative_I_ARRAY_LIKE_FUNCTION_TYPE
	JSTypeNative_I_ARRAY_LIKE_TYPE
	JSTypeNative_I_TEMPLATE_ARRAY_TYPE
	JSTypeNative_I_OBJECT_FUNCTION_TYPE
	JSTypeNative_I_OBJECT_TYPE
	JSTypeNative_I_THENABLE_FUNCTION_TYPE
	JSTypeNative_I_THENABLE_TYPE
	JSTypeNative_NULL_TYPE
	JSTypeNative_NUMBER_TYPE
	JSTypeNative_NUMBER_OBJECT_TYPE
	JSTypeNative_NUMBER_OBJECT_FUNCTION_TYPE
	JSTypeNative_PROMISE_TYPE
	JSTypeNative_PROMISE_FUNCTION_TYPE
	JSTypeNative_OBJECT_TYPE
	JSTypeNative_OBJECT_FUNCTION_TYPE
	JSTypeNative_OBJECT_PROTOTYPE
	JSTypeNative_READONLY_ARRAY_TYPE
	JSTypeNative_READONLY_ARRAY_FUNCTION_TYPE
	JSTypeNative_READONLY_MAP_TYPE
	JSTypeNative_READONLY_MAP_FUNCTION_TYPE
	JSTypeNative_MAP_TYPE
	JSTypeNative_MAP_FUNCTION_TYPE
	JSTypeNative_REGEXP_TYPE
	JSTypeNative_REGEXP_FUNCTION_TYPE
	JSTypeNative_STRING_OBJECT_TYPE
	JSTypeNative_STRING_OBJECT_FUNCTION_TYPE
	JSTypeNative_STRING_TYPE
	JSTypeNative_SYMBOL_OBJECT_TYPE
	JSTypeNative_SYMBOL_OBJECT_FUNCTION_TYPE
	JSTypeNative_SYMBOL_TYPE
	JSTypeNative_THENABLE_TYPE
	JSTypeNative_UNKNOWN_TYPE
	JSTypeNative_VOID_TYPE
	JSTypeNative_ALL_TYPE
	JSTypeNative_NO_TYPE
	JSTypeNative_NO_OBJECT_TYPE
	JSTypeNative_GLOBAL_THIS
	JSTypeNative_LEAST_FUNCTION_TYPE
	JSTypeNative_GREATEST_FUNCTION_TYPE
	JSTypeNative_NULL_VOID
	JSTypeNative_NUMBER_STRING_BOOLEAN
	JSTypeNative_VALUE_TYPES
	JSTypeNative_NUMBER_SYMBOL
	JSTypeNative_STRING_SYMBOL
	JSTypeNative_NUMBER_STRING
	JSTypeNative_NUMBER_STRING_SYMBOL
	JSTypeNative_BIGINT_NUMBER
	JSTypeNative_BIGINT_NUMBER_OBJECT
	JSTypeNative_BIGINT_NUMBER_STRING
	JSTypeNative_BIGINT_NUMBER_STRING_OBJECT
	JSTypeNative_NUMBER_ADDITION_SUPERTYPE
)

var jsTypeNative_names = [...]string{
	"ARGUMENTS_TYPE",
	"ARGUMENTS_FUNCTION_TYPE",
	"ARRAY_TYPE",
	"ARRAY_FUNCTION_TYPE",
	"ASYNC_ITERABLE_FUNCTION_TYPE",
	"ASYNC_ITERABLE_TYPE",
	"ASYNC_ITERATOR_FUNCTION_TYPE",
	"ASYNC_ITERATOR_TYPE",
	"ASYNC_ITERATOR_ITERABLE_FUNCTION_TYPE",
	"ASYNC_ITERATOR_ITERABLE_TYPE",
	"ASYNC_GENERATOR_FUNCTION_TYPE",
	"ASYNC_GENERATOR_TYPE",
	"BIGINT_TYPE",
	"BIGINT_OBJECT_TYPE",
	"BIGINT_OBJECT_FUNCTION_TYPE",
	"BOOLEAN_TYPE",
	"BOOLEAN_OBJECT_TYPE",
	"BOOLEAN_OBJECT_FUNCTION_TYPE",
	"CHECKED_UNKNOWN_TYPE",
	"DATE_TYPE",
	"DATE_FUNCTION_TYPE",
	"FUNCTION_TYPE",
	"FUNCTION_FUNCTION_TYPE",
	"FUNCTION_PROTOTYPE",
	"FUNCTION_INSTANCE_PROTOTYPE",
	"GBIGINT_TYPE",
	"GENERATOR_FUNCTION_TYPE",
	"GENERATOR_TYPE",
	"I_ITERABLE_RESULT_FUNCTION_TYPE",
	"I_ITERABLE_RESULT_TYPE",
	"ITERABLE_FUNCTION_TYPE",
	"ITERABLE_TYPE",
	"ITERATOR_LIKE_FUNCTION_TYPE",
	"ITERATOR_LIKE_TYPE",
	"ITERATOR_FUNCTION_TYPE",
	"ITERATOR_TYPE",
	"ITERATOR_ITERABLE_FUNCTION_TYPE",
	"ITERATOR_ITERABLE_TYPE",
	"I_ARRAY_LIKE_FUNCTION_TYPE",
	"I_ARRAY_LIKE_TYPE",
	"I_TEMPLATE_ARRAY_TYPE",
	"I_OBJECT_FUNCTION_TYPE",
	"I_OBJECT_TYPE",
	"I_THENABLE_FUNCTION_TYPE",
	"I_THENABLE_TYPE",
	"NULL_TYPE",
	"NUMBER_TYPE",
	"NUMBER_OBJECT_TYPE",
	"NUMBER_OBJECT_FUNCTION_TYPE",
	"PROMISE_TYPE",
	"PROMISE_FUNCTION_TYPE",
	"OBJECT_TYPE",
	"OBJECT_FUNCTION_TYPE",
	"OBJECT_PROTOTYPE",
	"READONLY_ARRAY_TYPE",
	"READONLY_ARRAY_FUNCTION_TYPE",
	"READONLY_MAP_TYPE",
	"READONLY_MAP_FUNCTION_TYPE",
	"MAP_TYPE",
	"MAP_FUNCTION_TYPE",
	"REGEXP_TYPE",
	"REGEXP_FUNCTION_TYPE",
	"STRING_OBJECT_TYPE",
	"STRING_OBJECT_FUNCTION_TYPE",
	"STRING_TYPE",
	"SYMBOL_OBJECT_TYPE",
	"SYMBOL_OBJECT_FUNCTION_TYPE",
	"SYMBOL_TYPE",
	"THENABLE_TYPE",
	"UNKNOWN_TYPE",
	"VOID_TYPE",
	"ALL_TYPE",
	"NO_TYPE",
	"NO_OBJECT_TYPE",
	"GLOBAL_THIS",
	"LEAST_FUNCTION_TYPE",
	"GREATEST_FUNCTION_TYPE",
	"NULL_VOID",
	"NUMBER_STRING_BOOLEAN",
	"VALUE_TYPES",
	"NUMBER_SYMBOL",
	"STRING_SYMBOL",
	"NUMBER_STRING",
	"NUMBER_STRING_SYMBOL",
	"BIGINT_NUMBER",
	"BIGINT_NUMBER_OBJECT",
	"BIGINT_NUMBER_STRING",
	"BIGINT_NUMBER_STRING_OBJECT",
	"NUMBER_ADDITION_SUPERTYPE",
}

// JSTypeNative_values is JSTypeNative.values().
func JSTypeNative_values() []JSTypeNative {
	out := make([]JSTypeNative, len(jsTypeNative_names))
	for i := range out {
		out[i] = JSTypeNative(i)
	}
	return out
}

// String is Enum.name().
func (n JSTypeNative) String() string { return jsTypeNative_names[n] }

// Ordinal is Enum.ordinal().
func (n JSTypeNative) Ordinal() int { return int(n) }
