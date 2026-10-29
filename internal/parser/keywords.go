// Port of com/google/javascript/jscomp/parsing/parser/Keywords.java @ closure-compiler v20261021

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

// Keywords: the JavaScript keywords.
type Keywords int32

const (

	// 7.6.1.1 Keywords
	Keywords_BREAK Keywords = iota
	Keywords_CASE
	Keywords_CATCH
	Keywords_CONTINUE
	Keywords_DEBUGGER
	Keywords_DEFAULT
	Keywords_DELETE
	Keywords_DO
	Keywords_ELSE
	Keywords_FINALLY
	Keywords_FOR
	Keywords_FUNCTION
	Keywords_IF
	Keywords_IN
	Keywords_INSTANCEOF
	Keywords_NEW
	Keywords_RETURN
	Keywords_SWITCH
	Keywords_THIS
	Keywords_THROW
	Keywords_TRY
	Keywords_TYPEOF
	Keywords_VAR
	Keywords_VOID
	Keywords_WHILE
	Keywords_WITH

	// 7.6.1.2 Future Reserved Words
	Keywords_CLASS
	Keywords_CONST
	Keywords_ENUM
	Keywords_EXPORT
	Keywords_EXTENDS
	Keywords_IMPORT
	Keywords_SUPER

	// Future Reserved Words in a strict context
	Keywords_IMPLEMENTS
	Keywords_INTERFACE
	Keywords_LET
	Keywords_PACKAGE
	Keywords_PRIVATE
	Keywords_PROTECTED
	Keywords_PUBLIC
	Keywords_STATIC
	Keywords_YIELD

	// 7.8 Literals
	Keywords_NULL
	Keywords_TRUE
	Keywords_FALSE
)

var keywords_names = [...]string{"BREAK", "CASE", "CATCH", "CONTINUE", "DEBUGGER", "DEFAULT", "DELETE", "DO", "ELSE", "FINALLY", "FOR", "FUNCTION", "IF", "IN", "INSTANCEOF", "NEW", "RETURN", "SWITCH", "THIS", "THROW", "TRY", "TYPEOF", "VAR", "VOID", "WHILE", "WITH", "CLASS", "CONST", "ENUM", "EXPORT", "EXTENDS", "IMPORT", "SUPER", "IMPLEMENTS", "INTERFACE", "LET", "PACKAGE", "PRIVATE", "PROTECTED", "PUBLIC", "STATIC", "YIELD", "NULL", "TRUE", "FALSE"}

// keywords_value is the Java field Keywords.value.
var keywords_value = [...]string{"break", "case", "catch", "continue", "debugger", "default", "delete", "do", "else", "finally", "for", "function", "if", "in", "instanceof", "new", "return", "switch", "this", "throw", "try", "typeof", "var", "void", "while", "with", "class", "const", "enum", "export", "extends", "import", "super", "implements", "interface", "let", "package", "private", "protected", "public", "static", "yield", "null", "true", "false"}

// keywords_type is the Java field Keywords.type.
var keywords_type = [...]TokenType{
	TokenType_BREAK,
	TokenType_CASE,
	TokenType_CATCH,
	TokenType_CONTINUE,
	TokenType_DEBUGGER,
	TokenType_DEFAULT,
	TokenType_DELETE,
	TokenType_DO,
	TokenType_ELSE,
	TokenType_FINALLY,
	TokenType_FOR,
	TokenType_FUNCTION,
	TokenType_IF,
	TokenType_IN,
	TokenType_INSTANCEOF,
	TokenType_NEW,
	TokenType_RETURN,
	TokenType_SWITCH,
	TokenType_THIS,
	TokenType_THROW,
	TokenType_TRY,
	TokenType_TYPEOF,
	TokenType_VAR,
	TokenType_VOID,
	TokenType_WHILE,
	TokenType_WITH,
	TokenType_CLASS,
	TokenType_CONST,
	TokenType_ENUM,
	TokenType_EXPORT,
	TokenType_EXTENDS,
	TokenType_IMPORT,
	TokenType_SUPER,
	TokenType_IMPLEMENTS,
	TokenType_INTERFACE,
	TokenType_LET,
	TokenType_PACKAGE,
	TokenType_PRIVATE,
	TokenType_PROTECTED,
	TokenType_PUBLIC,
	TokenType_STATIC,
	TokenType_YIELD,
	TokenType_NULL,
	TokenType_TRUE,
	TokenType_FALSE,
}
var (
	keywords_KEYWORDS_BY_TYPE [len(tokenType_names)]*Keywords
	// keywords_constants backs the *Keywords values handed out by the
	// lookups (a nil *Keywords is Java null).
	keywords_constants [len(keywords_names)]Keywords
)

// Java: Keywords static initializer
func init() {
	for i := range keywords_constants {
		kw := Keywords(i)
		keywords_constants[i] = kw
		keywords_KEYWORDS_BY_TYPE[keywords_type[i]] = &keywords_constants[i]
	}
}

// Keywords_values is Keywords.values().
//
// Java: Keywords.values()
func Keywords_values() []Keywords {
	values := make([]Keywords, len(keywords_names))
	for i := range values {
		values[i] = Keywords(i)
	}
	return values
}

// Name is Enum.name().
//
// Java: Keywords.name()
func (k Keywords) Name() string {
	return keywords_names[k]
}

// Value is the Java field Keywords.value.
//
// Java: field Keywords.value
func (k Keywords) Value() string {
	return keywords_value[k]
}

// Type is the Java field Keywords.type.
//
// Java: field Keywords.type
func (k Keywords) Type() TokenType {
	return keywords_type[k]
}

// String returns the keyword's text.
//
// Java: Keywords.toString()
func (k Keywords) String() string {
	return k.Value()
}

// Keywords_IsKeyword reports whether value is a keyword.
//
// Java: Keywords.isKeyword(String)
func Keywords_IsKeyword(value string) bool {
	return Keywords_Get(value) != nil
}

// Keywords_IsKeywordTokenType reports whether token is a keyword token.
//
// Java: Keywords.isKeyword(TokenType)
func Keywords_IsKeywordTokenType(token TokenType) bool {
	return Keywords_GetByTokenType(token) != nil
}

// Keywords_IsStrictKeyword returns true if token is a "future reserved word"
// which can be used as a variable identifier, but only in non-strict mode.
//
// Java: Keywords.isStrictKeyword(TokenType)
func Keywords_IsStrictKeyword(token TokenType) bool {
	switch token {
	case TokenType_IMPLEMENTS, TokenType_INTERFACE, TokenType_LET, TokenType_PACKAGE,
		TokenType_PRIVATE, TokenType_PROTECTED, TokenType_PUBLIC, TokenType_STATIC, TokenType_YIELD:
		return true
	default:
		return false
	}
}

// Keywords_Get returns the keyword named value, or nil.
//
// Java: Keywords.get(String)
func Keywords_Get(value string) *Keywords {
	switch value {
	case "break":
		return &keywords_constants[Keywords_BREAK]
	case "case":
		return &keywords_constants[Keywords_CASE]
	case "catch":
		return &keywords_constants[Keywords_CATCH]
	case "continue":
		return &keywords_constants[Keywords_CONTINUE]
	case "debugger":
		return &keywords_constants[Keywords_DEBUGGER]
	case "default":
		return &keywords_constants[Keywords_DEFAULT]
	case "delete":
		return &keywords_constants[Keywords_DELETE]
	case "do":
		return &keywords_constants[Keywords_DO]
	case "else":
		return &keywords_constants[Keywords_ELSE]
	case "finally":
		return &keywords_constants[Keywords_FINALLY]
	case "for":
		return &keywords_constants[Keywords_FOR]
	case "function":
		return &keywords_constants[Keywords_FUNCTION]
	case "if":
		return &keywords_constants[Keywords_IF]
	case "in":
		return &keywords_constants[Keywords_IN]
	case "instanceof":
		return &keywords_constants[Keywords_INSTANCEOF]
	case "new":
		return &keywords_constants[Keywords_NEW]
	case "return":
		return &keywords_constants[Keywords_RETURN]
	case "switch":
		return &keywords_constants[Keywords_SWITCH]
	case "this":
		return &keywords_constants[Keywords_THIS]
	case "throw":
		return &keywords_constants[Keywords_THROW]
	case "try":
		return &keywords_constants[Keywords_TRY]
	case "typeof":
		return &keywords_constants[Keywords_TYPEOF]
	case "var":
		return &keywords_constants[Keywords_VAR]
	case "void":
		return &keywords_constants[Keywords_VOID]
	case "while":
		return &keywords_constants[Keywords_WHILE]
	case "with":
		return &keywords_constants[Keywords_WITH]
	case "class":
		return &keywords_constants[Keywords_CLASS]
	case "const":
		return &keywords_constants[Keywords_CONST]
	case "enum":
		return &keywords_constants[Keywords_ENUM]
	case "export":
		return &keywords_constants[Keywords_EXPORT]
	case "extends":
		return &keywords_constants[Keywords_EXTENDS]
	case "import":
		return &keywords_constants[Keywords_IMPORT]
	case "super":
		return &keywords_constants[Keywords_SUPER]
	case "implements":
		return &keywords_constants[Keywords_IMPLEMENTS]
	case "interface":
		return &keywords_constants[Keywords_INTERFACE]
	case "let":
		return &keywords_constants[Keywords_LET]
	case "package":
		return &keywords_constants[Keywords_PACKAGE]
	case "private":
		return &keywords_constants[Keywords_PRIVATE]
	case "protected":
		return &keywords_constants[Keywords_PROTECTED]
	case "public":
		return &keywords_constants[Keywords_PUBLIC]
	case "static":
		return &keywords_constants[Keywords_STATIC]
	case "yield":
		return &keywords_constants[Keywords_YIELD]
	case "null":
		return &keywords_constants[Keywords_NULL]
	case "true":
		return &keywords_constants[Keywords_TRUE]
	case "false":
		return &keywords_constants[Keywords_FALSE]
	default:
		return nil
	}
}

// Keywords_GetByTokenType returns the keyword of token, or nil.
//
// Java: Keywords.get(TokenType)
func Keywords_GetByTokenType(token TokenType) *Keywords {
	return keywords_KEYWORDS_BY_TYPE[token]
}

// keywords_getBytes returns the keyword matching ASCII bytes b without
// allocating a string.
func keywords_getBytes(b []byte) *Keywords {
	if len(b) < 2 || len(b) > 10 {
		return nil
	}
	switch string(b) {
	case "break":
		return &keywords_constants[Keywords_BREAK]
	case "case":
		return &keywords_constants[Keywords_CASE]
	case "catch":
		return &keywords_constants[Keywords_CATCH]
	case "continue":
		return &keywords_constants[Keywords_CONTINUE]
	case "debugger":
		return &keywords_constants[Keywords_DEBUGGER]
	case "default":
		return &keywords_constants[Keywords_DEFAULT]
	case "delete":
		return &keywords_constants[Keywords_DELETE]
	case "do":
		return &keywords_constants[Keywords_DO]
	case "else":
		return &keywords_constants[Keywords_ELSE]
	case "finally":
		return &keywords_constants[Keywords_FINALLY]
	case "for":
		return &keywords_constants[Keywords_FOR]
	case "function":
		return &keywords_constants[Keywords_FUNCTION]
	case "if":
		return &keywords_constants[Keywords_IF]
	case "in":
		return &keywords_constants[Keywords_IN]
	case "instanceof":
		return &keywords_constants[Keywords_INSTANCEOF]
	case "new":
		return &keywords_constants[Keywords_NEW]
	case "return":
		return &keywords_constants[Keywords_RETURN]
	case "switch":
		return &keywords_constants[Keywords_SWITCH]
	case "this":
		return &keywords_constants[Keywords_THIS]
	case "throw":
		return &keywords_constants[Keywords_THROW]
	case "try":
		return &keywords_constants[Keywords_TRY]
	case "typeof":
		return &keywords_constants[Keywords_TYPEOF]
	case "var":
		return &keywords_constants[Keywords_VAR]
	case "void":
		return &keywords_constants[Keywords_VOID]
	case "while":
		return &keywords_constants[Keywords_WHILE]
	case "with":
		return &keywords_constants[Keywords_WITH]
	case "class":
		return &keywords_constants[Keywords_CLASS]
	case "const":
		return &keywords_constants[Keywords_CONST]
	case "enum":
		return &keywords_constants[Keywords_ENUM]
	case "export":
		return &keywords_constants[Keywords_EXPORT]
	case "extends":
		return &keywords_constants[Keywords_EXTENDS]
	case "import":
		return &keywords_constants[Keywords_IMPORT]
	case "super":
		return &keywords_constants[Keywords_SUPER]
	case "implements":
		return &keywords_constants[Keywords_IMPLEMENTS]
	case "interface":
		return &keywords_constants[Keywords_INTERFACE]
	case "let":
		return &keywords_constants[Keywords_LET]
	case "package":
		return &keywords_constants[Keywords_PACKAGE]
	case "private":
		return &keywords_constants[Keywords_PRIVATE]
	case "protected":
		return &keywords_constants[Keywords_PROTECTED]
	case "public":
		return &keywords_constants[Keywords_PUBLIC]
	case "static":
		return &keywords_constants[Keywords_STATIC]
	case "yield":
		return &keywords_constants[Keywords_YIELD]
	case "null":
		return &keywords_constants[Keywords_NULL]
	case "true":
		return &keywords_constants[Keywords_TRUE]
	case "false":
		return &keywords_constants[Keywords_FALSE]
	default:
		return nil
	}
}
