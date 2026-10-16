// Port of com/google/javascript/rhino/Token.java @ closure-compiler v20260415

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1/GPL 2.0
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
 * Alternatively, the contents of this file may be used under the terms of
 * the GNU General Public License Version 2 or later (the "GPL"), in which
 * case the provisions of the GPL are applicable instead of those above. If
 * you wish to allow use of your version of this file only under the terms of
 * the GPL and not to allow others to use your version of this file under the
 * MPL, indicate your decision by deleting the provisions above and replacing
 * them with the notice and other provisions required by the GPL. If you do
 * not delete the provisions above, a recipient may use your version of this
 * file under either the MPL or the GPL.
 *
 * ***** END LICENSE BLOCK ***** */

// Package rhino is the port of com.google.javascript.rhino (and, later,
// com.google.javascript.rhino.jstype): the AST and JSDoc model.
package rhino

import "github.com/paullepoulpe/jsgo/internal/jutil"

// Token implements the JavaScript scanner's token kinds.
//
// It is based on the C source files jsscan.c and jsscan.h in the jsref package.
type Token int32

const (
	Token_RETURN Token = iota
	Token_BITOR
	Token_BITXOR
	Token_BITAND
	Token_EQ
	Token_NE
	Token_LT
	Token_LE
	Token_GT
	Token_GE
	Token_LSH
	Token_RSH
	Token_URSH
	Token_ADD
	Token_SUB
	Token_MUL
	Token_DIV
	Token_MOD
	Token_EXPONENT
	Token_NOT
	Token_BITNOT
	Token_POS
	Token_NEG
	Token_NEW
	Token_DELPROP
	Token_TYPEOF
	Token_GETPROP
	Token_GETELEM
	Token_CALL

	// Part of optional chain (?.)
	Token_OPTCHAIN_GETPROP
	Token_OPTCHAIN_GETELEM
	Token_OPTCHAIN_CALL

	Token_NAME
	Token_NUMBER
	Token_BIGINT
	Token_STRINGLIT
	Token_NULL
	Token_THIS
	Token_FALSE
	Token_TRUE
	Token_SHEQ // shallow equality (===)
	Token_SHNE // shallow inequality (!==)
	Token_REGEXP
	Token_THROW
	Token_IN
	Token_INSTANCEOF
	Token_ARRAYLIT  // array literal
	Token_OBJECTLIT // object literal

	Token_TRY
	Token_PARAM_LIST
	Token_COMMA // comma operator

	Token_ASSIGN          // simple assignment  (=)
	Token_ASSIGN_BITOR    // |=
	Token_ASSIGN_BITXOR   // ^=
	Token_ASSIGN_BITAND   // &=
	Token_ASSIGN_LSH      // <<=
	Token_ASSIGN_RSH      // >>=
	Token_ASSIGN_URSH     // >>>=
	Token_ASSIGN_ADD      // +=
	Token_ASSIGN_SUB      // -=
	Token_ASSIGN_MUL      // *=
	Token_ASSIGN_DIV      // /=
	Token_ASSIGN_MOD      // %=
	Token_ASSIGN_EXPONENT // **=

	Token_ASSIGN_OR       // logical or assignment (||=)
	Token_ASSIGN_AND      // logical and assignment (&&=)
	Token_ASSIGN_COALESCE // logical nullish coalesce assignment (??=)

	Token_HOOK         // conditional (?:)
	Token_OR           // logical or (||)
	Token_AND          // logical and (&&)
	Token_COALESCE     // Nullish coalesce (??)
	Token_INC          // increment (++)
	Token_DEC          // decrement (--)
	Token_FUNCTION     // function keyword
	Token_IF           // if keyword
	Token_SWITCH       // switch keyword
	Token_CASE         // case keyword
	Token_DEFAULT_CASE // default keyword
	Token_WHILE        // while keyword
	Token_DO           // do keyword
	Token_FOR          // for(;;) statement
	Token_FOR_IN       // for-in
	Token_BREAK        // break keyword
	Token_CONTINUE     // continue keyword
	Token_VAR          // var keyword
	Token_WITH         // with keyword
	Token_CATCH        // catch keyword
	Token_VOID         // void keyword

	Token_EMPTY

	Token_ROOT        // Used only for the 3 root nodes of the AST: externsRoot, jsRoot, and externsAndJsRoot
	Token_BLOCK       // statement block
	Token_SWITCH_BODY // holds case statements within a switch
	Token_LABEL       // label
	Token_EXPR_RESULT // expression statement in scripts
	Token_SCRIPT      // top-level node for entire script

	Token_GETTER_DEF
	Token_SETTER_DEF

	Token_CONST // JS 1.5 const keyword
	Token_DEBUGGER

	// JSCompiler introduced tokens
	Token_LABEL_NAME
	Token_STRING_KEY // object literal key
	Token_CAST

	// ES6
	Token_ARRAY_PATTERN // destructuring patterns
	Token_OBJECT_PATTERN
	Token_DESTRUCTURING_LHS // The node inside a var/let/const with a destructuring LHS

	Token_CLASS         // classes
	Token_CLASS_MEMBERS // class member container
	Token_MEMBER_FUNCTION_DEF
	Token_MEMBER_FIELD_DEF   // class field general case
	Token_COMPUTED_FIELD_DEF // class field computed property case
	Token_SUPER

	Token_LET // block scoped vars

	Token_FOR_OF       // for-of
	Token_FOR_AWAIT_OF // for-await-of

	Token_YIELD // generators

	Token_AWAIT // async functions

	Token_IMPORT // modules
	Token_IMPORT_SPECS
	Token_IMPORT_SPEC
	Token_IMPORT_STAR // "* as name", called NameSpaceImport in the spec.
	Token_EXPORT
	Token_EXPORT_SPECS
	Token_EXPORT_SPEC
	Token_MODULE_BODY
	Token_DYNAMIC_IMPORT

	Token_ITER_REST     // Rests that use the iterator protocol.
	Token_OBJECT_REST   // Rests that get object properties.
	Token_ITER_SPREAD   // Spreads that use the iterator protocol.
	Token_OBJECT_SPREAD // Spreads that get object properties.

	Token_COMPUTED_PROP

	Token_TAGGED_TEMPLATELIT // tagged template literal, e.g. foo`bar`
	Token_TEMPLATELIT        // template literal, e.g: `bar`
	Token_TEMPLATELIT_SUB    // template literal substitution
	Token_TEMPLATELIT_STRING // template literal string

	Token_DEFAULT_VALUE // Formal parameter or destructuring element with a default value
	Token_NEW_TARGET    // new.target
	Token_IMPORT_META   // import.meta

	// Used by type declaration ASTs
	Token_STRING_TYPE
	Token_BOOLEAN_TYPE
	Token_NUMBER_TYPE
	Token_FUNCTION_TYPE
	Token_PARAMETERIZED_TYPE
	Token_UNION_TYPE
	Token_ANY_TYPE
	Token_NULLABLE_TYPE
	Token_VOID_TYPE
	Token_REST_PARAMETER_TYPE
	Token_NAMED_TYPE
	Token_OPTIONAL_PARAMETER
	Token_RECORD_TYPE
	Token_UNDEFINED_TYPE
	Token_ARRAY_TYPE
	Token_GENERIC_TYPE
	Token_GENERIC_TYPE_LIST

	// JSDoc-only tokens
	Token_ANNOTATION
	Token_PIPE
	Token_STAR
	Token_EOC
	Token_QMARK // type is nullable or unknown
	Token_BANG
	Token_EQUALS
	Token_LB // left brackets
	Token_LC // left curly braces
	Token_COLON

	// TypeScript
	Token_INTERFACE
	Token_INTERFACE_EXTENDS
	Token_INTERFACE_MEMBERS
	Token_ENUM
	Token_ENUM_MEMBERS
	Token_IMPLEMENTS
	Token_TYPE_ALIAS
	Token_DECLARE
	Token_MEMBER_VARIABLE_DEF
	Token_INDEX_SIGNATURE
	Token_CALL_SIGNATURE
	Token_NAMESPACE
	Token_NAMESPACE_ELEMENTS

	// Tokens to use for internal bookkeeping,
	// an AST is invalid while these are present.
	Token_PLACEHOLDER1
	Token_PLACEHOLDER2
	Token_PLACEHOLDER3
)

var token_names = [...]string{
	"RETURN",
	"BITOR",
	"BITXOR",
	"BITAND",
	"EQ",
	"NE",
	"LT",
	"LE",
	"GT",
	"GE",
	"LSH",
	"RSH",
	"URSH",
	"ADD",
	"SUB",
	"MUL",
	"DIV",
	"MOD",
	"EXPONENT",
	"NOT",
	"BITNOT",
	"POS",
	"NEG",
	"NEW",
	"DELPROP",
	"TYPEOF",
	"GETPROP",
	"GETELEM",
	"CALL",
	"OPTCHAIN_GETPROP",
	"OPTCHAIN_GETELEM",
	"OPTCHAIN_CALL",
	"NAME",
	"NUMBER",
	"BIGINT",
	"STRINGLIT",
	"NULL",
	"THIS",
	"FALSE",
	"TRUE",
	"SHEQ",
	"SHNE",
	"REGEXP",
	"THROW",
	"IN",
	"INSTANCEOF",
	"ARRAYLIT",
	"OBJECTLIT",
	"TRY",
	"PARAM_LIST",
	"COMMA",
	"ASSIGN",
	"ASSIGN_BITOR",
	"ASSIGN_BITXOR",
	"ASSIGN_BITAND",
	"ASSIGN_LSH",
	"ASSIGN_RSH",
	"ASSIGN_URSH",
	"ASSIGN_ADD",
	"ASSIGN_SUB",
	"ASSIGN_MUL",
	"ASSIGN_DIV",
	"ASSIGN_MOD",
	"ASSIGN_EXPONENT",
	"ASSIGN_OR",
	"ASSIGN_AND",
	"ASSIGN_COALESCE",
	"HOOK",
	"OR",
	"AND",
	"COALESCE",
	"INC",
	"DEC",
	"FUNCTION",
	"IF",
	"SWITCH",
	"CASE",
	"DEFAULT_CASE",
	"WHILE",
	"DO",
	"FOR",
	"FOR_IN",
	"BREAK",
	"CONTINUE",
	"VAR",
	"WITH",
	"CATCH",
	"VOID",
	"EMPTY",
	"ROOT",
	"BLOCK",
	"SWITCH_BODY",
	"LABEL",
	"EXPR_RESULT",
	"SCRIPT",
	"GETTER_DEF",
	"SETTER_DEF",
	"CONST",
	"DEBUGGER",
	"LABEL_NAME",
	"STRING_KEY",
	"CAST",
	"ARRAY_PATTERN",
	"OBJECT_PATTERN",
	"DESTRUCTURING_LHS",
	"CLASS",
	"CLASS_MEMBERS",
	"MEMBER_FUNCTION_DEF",
	"MEMBER_FIELD_DEF",
	"COMPUTED_FIELD_DEF",
	"SUPER",
	"LET",
	"FOR_OF",
	"FOR_AWAIT_OF",
	"YIELD",
	"AWAIT",
	"IMPORT",
	"IMPORT_SPECS",
	"IMPORT_SPEC",
	"IMPORT_STAR",
	"EXPORT",
	"EXPORT_SPECS",
	"EXPORT_SPEC",
	"MODULE_BODY",
	"DYNAMIC_IMPORT",
	"ITER_REST",
	"OBJECT_REST",
	"ITER_SPREAD",
	"OBJECT_SPREAD",
	"COMPUTED_PROP",
	"TAGGED_TEMPLATELIT",
	"TEMPLATELIT",
	"TEMPLATELIT_SUB",
	"TEMPLATELIT_STRING",
	"DEFAULT_VALUE",
	"NEW_TARGET",
	"IMPORT_META",
	"STRING_TYPE",
	"BOOLEAN_TYPE",
	"NUMBER_TYPE",
	"FUNCTION_TYPE",
	"PARAMETERIZED_TYPE",
	"UNION_TYPE",
	"ANY_TYPE",
	"NULLABLE_TYPE",
	"VOID_TYPE",
	"REST_PARAMETER_TYPE",
	"NAMED_TYPE",
	"OPTIONAL_PARAMETER",
	"RECORD_TYPE",
	"UNDEFINED_TYPE",
	"ARRAY_TYPE",
	"GENERIC_TYPE",
	"GENERIC_TYPE_LIST",
	"ANNOTATION",
	"PIPE",
	"STAR",
	"EOC",
	"QMARK",
	"BANG",
	"EQUALS",
	"LB",
	"LC",
	"COLON",
	"INTERFACE",
	"INTERFACE_EXTENDS",
	"INTERFACE_MEMBERS",
	"ENUM",
	"ENUM_MEMBERS",
	"IMPLEMENTS",
	"TYPE_ALIAS",
	"DECLARE",
	"MEMBER_VARIABLE_DEF",
	"INDEX_SIGNATURE",
	"CALL_SIGNATURE",
	"NAMESPACE",
	"NAMESPACE_ELEMENTS",
	"PLACEHOLDER1",
	"PLACEHOLDER2",
	"PLACEHOLDER3",
}

// Token_values is Token.values().
func Token_values() []Token {
	out := make([]Token, len(token_names))
	for i := range out {
		out[i] = Token(i)
	}
	return out
}

// String is Enum.name() / toString().
func (t Token) String() string { return token_names[t] }

// Name is Enum.name().
func (t Token) Name() string { return token_names[t] }

// Ordinal is Enum.ordinal().
func (t Token) Ordinal() int { return int(t) }

// Token_ValueOf is Token.valueOf(String).
func Token_ValueOf(name string) Token {
	for i, n := range token_names {
		if n == name {
			return Token(i)
		}
	}
	panic(jutil.NewIllegalArgumentExceptionMsg("No enum constant com.google.javascript.rhino.Token." + name))
}

// Token_Arity returns the arity of the token; if the arity isn't always the
// same, this function returns -1.
//
// Java: Token.arity(Token)
func Token_Arity(token Token) int {
	switch token {
	case Token_ANNOTATION,
		Token_ARRAYLIT,
		Token_BANG,
		Token_BLOCK,
		Token_ROOT,
		Token_BREAK,
		Token_CALL,
		Token_OPTCHAIN_CALL,
		Token_COLON,
		Token_CONST,
		Token_CONTINUE,
		Token_DEBUGGER,
		Token_EOC,
		Token_EQUALS,
		Token_FOR,
		Token_IF,
		Token_LB,
		Token_LC,
		Token_NEW,
		Token_OBJECTLIT,
		Token_PARAM_LIST,
		Token_PIPE,
		Token_QMARK,
		Token_REGEXP,
		Token_RETURN,
		Token_SCRIPT,
		Token_STAR,
		Token_STRING_KEY,
		Token_SWITCH_BODY,
		Token_TEMPLATELIT,
		Token_TRY,
		Token_VAR,
		Token_YIELD,
		Token_MEMBER_FIELD_DEF,
		Token_COMPUTED_FIELD_DEF:
		return -1
	case Token_EMPTY,
		Token_FALSE,
		Token_IMPORT_STAR,
		Token_LABEL_NAME,
		Token_MEMBER_VARIABLE_DEF,
		Token_NAME,
		Token_NULL,
		Token_NUMBER,
		Token_BIGINT,
		Token_STRINGLIT,
		Token_TEMPLATELIT_STRING,
		Token_THIS,
		Token_TRUE:
		return 0
	case Token_AWAIT,
		Token_BITNOT,
		Token_CALL_SIGNATURE,
		Token_CAST,
		Token_DEC,
		Token_DEFAULT_CASE,
		Token_DELPROP,
		Token_EXPR_RESULT,
		Token_GETPROP,
		Token_GETTER_DEF,
		Token_INC,
		Token_INDEX_SIGNATURE,
		Token_ITER_REST,
		Token_ITER_SPREAD,
		Token_MEMBER_FUNCTION_DEF,
		Token_NAMED_TYPE,
		Token_NEG,
		Token_NOT,
		Token_OBJECT_REST,
		Token_OBJECT_SPREAD,
		Token_OPTCHAIN_GETPROP,
		Token_POS,
		Token_SETTER_DEF,
		Token_TEMPLATELIT_SUB,
		Token_THROW,
		Token_TYPEOF,
		Token_TYPE_ALIAS,
		Token_VOID:
		return 1
	case Token_ADD,
		Token_AND,
		Token_ASSIGN,
		Token_ASSIGN_ADD,
		Token_ASSIGN_BITAND,
		Token_ASSIGN_BITOR,
		Token_ASSIGN_BITXOR,
		Token_ASSIGN_DIV,
		Token_ASSIGN_LSH,
		Token_ASSIGN_MOD,
		Token_ASSIGN_MUL,
		Token_ASSIGN_EXPONENT,
		Token_ASSIGN_RSH,
		Token_ASSIGN_SUB,
		Token_ASSIGN_URSH,
		Token_ASSIGN_OR,
		Token_ASSIGN_AND,
		Token_ASSIGN_COALESCE,
		Token_BITAND,
		Token_BITOR,
		Token_BITXOR,
		Token_CASE,
		Token_COALESCE,
		Token_CATCH,
		Token_COMMA,
		Token_COMPUTED_PROP,
		Token_DEFAULT_VALUE,
		Token_DIV,
		Token_DO,
		Token_ENUM,
		Token_EQ,
		Token_EXPONENT,
		Token_GE,
		Token_GETELEM,
		Token_OPTCHAIN_GETELEM,
		Token_GT,
		Token_IN,
		Token_INSTANCEOF,
		Token_LABEL,
		Token_LE,
		Token_LSH,
		Token_LT,
		Token_MOD,
		Token_MUL,
		Token_NAMESPACE,
		Token_NE,
		Token_OR,
		Token_RSH,
		Token_SWITCH,
		Token_SHEQ,
		Token_SHNE,
		Token_SUB,
		Token_TAGGED_TEMPLATELIT,
		Token_URSH,
		Token_WHILE,
		Token_WITH:
		return 2
	case Token_CLASS,
		Token_FOR_IN,
		Token_FOR_OF,
		Token_FOR_AWAIT_OF,
		Token_FUNCTION,
		Token_HOOK,
		Token_IMPORT,
		Token_INTERFACE:
		return 3
	}
	panic(jutil.NewIllegalStateExceptionMsg("No arity defined for " + token.String()))
}
