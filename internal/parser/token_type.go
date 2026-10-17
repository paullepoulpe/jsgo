// Port of com/google/javascript/jscomp/parsing/parser/TokenType.java @ closure-compiler v20260415

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

package parser

// TokenType: 7.5 Tokens
type TokenType int32

const (
	TokenType_END_OF_FILE TokenType = iota
	TokenType_ERROR

	// 7.6 Identifier Names and Identifiers
	TokenType_IDENTIFIER

	// 7.6.1.1 keywords
	TokenType_BREAK
	TokenType_CASE
	TokenType_CATCH
	TokenType_CONTINUE
	TokenType_DEBUGGER
	TokenType_DEFAULT
	TokenType_DELETE
	TokenType_DO
	TokenType_ELSE
	TokenType_FINALLY
	TokenType_FOR
	TokenType_FUNCTION
	TokenType_IF
	TokenType_IN
	TokenType_INSTANCEOF
	TokenType_NEW
	TokenType_RETURN
	TokenType_SWITCH
	TokenType_THIS
	TokenType_THROW
	TokenType_TRY
	TokenType_TYPEOF
	TokenType_VAR
	TokenType_VOID
	TokenType_WHILE
	TokenType_WITH

	// 7.6.1.2 Future reserved words
	TokenType_CLASS
	TokenType_CONST
	TokenType_ENUM
	TokenType_EXPORT
	TokenType_EXTENDS
	TokenType_IMPORT
	TokenType_SUPER

	// Future reserved words in strict mode
	TokenType_IMPLEMENTS
	TokenType_INTERFACE
	TokenType_LET
	TokenType_PACKAGE
	TokenType_PRIVATE
	TokenType_PROTECTED
	TokenType_PUBLIC
	TokenType_STATIC
	TokenType_YIELD

	// 7.7 Punctuators
	TokenType_OPEN_CURLY
	TokenType_CLOSE_CURLY
	TokenType_OPEN_PAREN
	TokenType_CLOSE_PAREN
	TokenType_OPEN_SQUARE
	TokenType_CLOSE_SQUARE
	TokenType_PERIOD
	TokenType_SEMI_COLON
	TokenType_COMMA
	TokenType_OPEN_ANGLE
	TokenType_CLOSE_ANGLE
	TokenType_LESS_EQUAL
	TokenType_GREATER_EQUAL
	TokenType_ARROW
	TokenType_EQUAL_EQUAL
	TokenType_NOT_EQUAL
	TokenType_EQUAL_EQUAL_EQUAL
	TokenType_NOT_EQUAL_EQUAL
	TokenType_PLUS
	TokenType_MINUS
	TokenType_STAR
	TokenType_STAR_STAR
	TokenType_PERCENT
	TokenType_PLUS_PLUS
	TokenType_MINUS_MINUS
	TokenType_LEFT_SHIFT
	TokenType_RIGHT_SHIFT
	TokenType_UNSIGNED_RIGHT_SHIFT
	TokenType_AMPERSAND
	TokenType_BAR
	TokenType_CARET
	TokenType_BANG
	TokenType_TILDE
	TokenType_AND
	TokenType_OR
	TokenType_QUESTION
	TokenType_QUESTION_QUESTION
	TokenType_QUESTION_DOT
	TokenType_COLON
	TokenType_EQUAL
	TokenType_PLUS_EQUAL
	TokenType_MINUS_EQUAL
	TokenType_STAR_EQUAL
	TokenType_STAR_STAR_EQUAL
	TokenType_PERCENT_EQUAL
	TokenType_LEFT_SHIFT_EQUAL
	TokenType_RIGHT_SHIFT_EQUAL
	TokenType_UNSIGNED_RIGHT_SHIFT_EQUAL
	TokenType_AMPERSAND_EQUAL
	TokenType_BAR_EQUAL
	TokenType_CARET_EQUAL
	TokenType_SLASH
	TokenType_SLASH_EQUAL

	// Logical Assignment Punctuators
	TokenType_AND_EQUAL
	TokenType_OR_EQUAL
	TokenType_QUESTION_QUESTION_EQUAL

	// 7.8 Literals
	TokenType_NULL
	TokenType_TRUE
	TokenType_FALSE
	TokenType_NUMBER
	TokenType_STRING
	TokenType_BIGINT
	TokenType_REGULAR_EXPRESSION

	// Harmony extensions
	TokenType_ELLIPSIS
	// 12.2.9 Template Literals
	// Template literal tokens corresponding to different parts of the literal
	// Eg: `hello` is scanned as a single NO_SUBSTITUTION_TEMPLATE: hello
	// `hello${world}!` is scanned as TEMPLATE_HEAD: hello, TEMPLATE_MIDDLE: world,
	// and TEMPLATE_TAIL: !
	TokenType_TEMPLATE_HEAD
	TokenType_TEMPLATE_MIDDLE
	TokenType_TEMPLATE_TAIL
	TokenType_NO_SUBSTITUTION_TEMPLATE

	// TypeScript
	TokenType_TYPE
	TokenType_DECLARE
	TokenType_MODULE
	TokenType_NAMESPACE
)

var tokenType_names = [...]string{
	"END_OF_FILE",
	"ERROR",
	"IDENTIFIER",
	"BREAK",
	"CASE",
	"CATCH",
	"CONTINUE",
	"DEBUGGER",
	"DEFAULT",
	"DELETE",
	"DO",
	"ELSE",
	"FINALLY",
	"FOR",
	"FUNCTION",
	"IF",
	"IN",
	"INSTANCEOF",
	"NEW",
	"RETURN",
	"SWITCH",
	"THIS",
	"THROW",
	"TRY",
	"TYPEOF",
	"VAR",
	"VOID",
	"WHILE",
	"WITH",
	"CLASS",
	"CONST",
	"ENUM",
	"EXPORT",
	"EXTENDS",
	"IMPORT",
	"SUPER",
	"IMPLEMENTS",
	"INTERFACE",
	"LET",
	"PACKAGE",
	"PRIVATE",
	"PROTECTED",
	"PUBLIC",
	"STATIC",
	"YIELD",
	"OPEN_CURLY",
	"CLOSE_CURLY",
	"OPEN_PAREN",
	"CLOSE_PAREN",
	"OPEN_SQUARE",
	"CLOSE_SQUARE",
	"PERIOD",
	"SEMI_COLON",
	"COMMA",
	"OPEN_ANGLE",
	"CLOSE_ANGLE",
	"LESS_EQUAL",
	"GREATER_EQUAL",
	"ARROW",
	"EQUAL_EQUAL",
	"NOT_EQUAL",
	"EQUAL_EQUAL_EQUAL",
	"NOT_EQUAL_EQUAL",
	"PLUS",
	"MINUS",
	"STAR",
	"STAR_STAR",
	"PERCENT",
	"PLUS_PLUS",
	"MINUS_MINUS",
	"LEFT_SHIFT",
	"RIGHT_SHIFT",
	"UNSIGNED_RIGHT_SHIFT",
	"AMPERSAND",
	"BAR",
	"CARET",
	"BANG",
	"TILDE",
	"AND",
	"OR",
	"QUESTION",
	"QUESTION_QUESTION",
	"QUESTION_DOT",
	"COLON",
	"EQUAL",
	"PLUS_EQUAL",
	"MINUS_EQUAL",
	"STAR_EQUAL",
	"STAR_STAR_EQUAL",
	"PERCENT_EQUAL",
	"LEFT_SHIFT_EQUAL",
	"RIGHT_SHIFT_EQUAL",
	"UNSIGNED_RIGHT_SHIFT_EQUAL",
	"AMPERSAND_EQUAL",
	"BAR_EQUAL",
	"CARET_EQUAL",
	"SLASH",
	"SLASH_EQUAL",
	"AND_EQUAL",
	"OR_EQUAL",
	"QUESTION_QUESTION_EQUAL",
	"NULL",
	"TRUE",
	"FALSE",
	"NUMBER",
	"STRING",
	"BIGINT",
	"REGULAR_EXPRESSION",
	"ELLIPSIS",
	"TEMPLATE_HEAD",
	"TEMPLATE_MIDDLE",
	"TEMPLATE_TAIL",
	"NO_SUBSTITUTION_TEMPLATE",
	"TYPE",
	"DECLARE",
	"MODULE",
	"NAMESPACE",
}

// tokenType_value is the Java field TokenType.value; "" stands for Java
// null (no TokenType has an empty value).
var tokenType_value = [...]string{
	"End of File",                // END_OF_FILE
	"error",                      // ERROR
	"identifier",                 // IDENTIFIER
	"",                           // BREAK
	"",                           // CASE
	"",                           // CATCH
	"",                           // CONTINUE
	"",                           // DEBUGGER
	"",                           // DEFAULT
	"",                           // DELETE
	"",                           // DO
	"",                           // ELSE
	"",                           // FINALLY
	"",                           // FOR
	"",                           // FUNCTION
	"",                           // IF
	"",                           // IN
	"",                           // INSTANCEOF
	"",                           // NEW
	"",                           // RETURN
	"",                           // SWITCH
	"",                           // THIS
	"",                           // THROW
	"",                           // TRY
	"",                           // TYPEOF
	"",                           // VAR
	"",                           // VOID
	"",                           // WHILE
	"",                           // WITH
	"",                           // CLASS
	"",                           // CONST
	"",                           // ENUM
	"",                           // EXPORT
	"",                           // EXTENDS
	"",                           // IMPORT
	"",                           // SUPER
	"",                           // IMPLEMENTS
	"",                           // INTERFACE
	"",                           // LET
	"",                           // PACKAGE
	"",                           // PRIVATE
	"",                           // PROTECTED
	"",                           // PUBLIC
	"",                           // STATIC
	"",                           // YIELD
	"{",                          // OPEN_CURLY
	"}",                          // CLOSE_CURLY
	"(",                          // OPEN_PAREN
	")",                          // CLOSE_PAREN
	"[",                          // OPEN_SQUARE
	"]",                          // CLOSE_SQUARE
	".",                          // PERIOD
	";",                          // SEMI_COLON
	",",                          // COMMA
	"<",                          // OPEN_ANGLE
	">",                          // CLOSE_ANGLE
	"<=",                         // LESS_EQUAL
	">=",                         // GREATER_EQUAL
	"=>",                         // ARROW
	"==",                         // EQUAL_EQUAL
	"!=",                         // NOT_EQUAL
	"===",                        // EQUAL_EQUAL_EQUAL
	"!==",                        // NOT_EQUAL_EQUAL
	"+",                          // PLUS
	"-",                          // MINUS
	"*",                          // STAR
	"**",                         // STAR_STAR
	"%",                          // PERCENT
	"++",                         // PLUS_PLUS
	"--",                         // MINUS_MINUS
	"<<",                         // LEFT_SHIFT
	">>",                         // RIGHT_SHIFT
	">>>",                        // UNSIGNED_RIGHT_SHIFT
	"&",                          // AMPERSAND
	"|",                          // BAR
	"^",                          // CARET
	"!",                          // BANG
	"~",                          // TILDE
	"&&",                         // AND
	"||",                         // OR
	"?",                          // QUESTION
	"??",                         // QUESTION_QUESTION
	"?.",                         // QUESTION_DOT
	":",                          // COLON
	"=",                          // EQUAL
	"+=",                         // PLUS_EQUAL
	"-=",                         // MINUS_EQUAL
	"*=",                         // STAR_EQUAL
	"**=",                        // STAR_STAR_EQUAL
	"%=",                         // PERCENT_EQUAL
	"<<=",                        // LEFT_SHIFT_EQUAL
	">>=",                        // RIGHT_SHIFT_EQUAL
	">>>=",                       // UNSIGNED_RIGHT_SHIFT_EQUAL
	"&=",                         // AMPERSAND_EQUAL
	"|=",                         // BAR_EQUAL
	"^=",                         // CARET_EQUAL
	"/",                          // SLASH
	"/=",                         // SLASH_EQUAL
	"&&=",                        // AND_EQUAL
	"||=",                        // OR_EQUAL
	"??=",                        // QUESTION_QUESTION_EQUAL
	"",                           // NULL
	"",                           // TRUE
	"",                           // FALSE
	"number literal",             // NUMBER
	"string literal",             // STRING
	"bigint literal",             // BIGINT
	"regular expression literal", // REGULAR_EXPRESSION
	"...",                        // ELLIPSIS
	"template head",              // TEMPLATE_HEAD
	"template middle",            // TEMPLATE_MIDDLE
	"template tail",              // TEMPLATE_TAIL
	"no substitution template",   // NO_SUBSTITUTION_TEMPLATE
	"",                           // TYPE
	"",                           // DECLARE
	"",                           // MODULE
	"",                           // NAMESPACE
}

// TokenType_values is TokenType.values().
//
// Java: TokenType.values()
func TokenType_values() []TokenType {
	values := make([]TokenType, len(tokenType_names))
	for i := range values {
		values[i] = TokenType(i)
	}
	return values
}

// Name is Enum.name().
//
// Java: TokenType.name()
func (t TokenType) Name() string {
	return tokenType_names[t]
}

// String returns the token's text, or the keyword for keyword tokens.
//
// Java: TokenType.toString()
func (t TokenType) String() string {
	// TODO: straighten this out
	if tokenType_value[t] == "" {
		return Keywords_GetByTokenType(t).String()
	}
	return tokenType_value[t]
}
