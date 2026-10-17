// Port of com/google/javascript/jscomp/parsing/parser/Keywords.java @ closure-compiler v20261015

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
	keywords_KEYWORDS_BY_NAME map[string]*Keywords
	keywords_KEYWORDS_BY_TYPE [len(tokenType_names)]*Keywords
	// keywords_constants backs the *Keywords values handed out by the
	// lookups (a nil *Keywords is Java null).
	keywords_constants [len(keywords_names)]Keywords
)

// Java: Keywords static initializer
func init() {
	// Go-only: language-adapter — KEYWORDS_BY_NAME is only ever looked up,
	// never iterated, so a Go map is order-safe; KEYWORDS_BY_TYPE is an
	// EnumMap, i.e. an array indexed by TokenType ordinal.
	keywordsByName := make(map[string]*Keywords)
	for _, kw := range Keywords_values() {
		keywords_constants[kw] = kw
		p := &keywords_constants[kw]
		if _, dup := keywordsByName[kw.Value()]; dup {
			panic("Multiple entries with same key: " + kw.Value())
		}
		keywordsByName[kw.Value()] = p
		keywords_KEYWORDS_BY_TYPE[kw.Type()] = p
	}
	keywords_KEYWORDS_BY_NAME = keywordsByName
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

// Keywords_GetTokenType returns the token type of keyword value (a Go nil
// dereference, like Java's NPE, if value is not a keyword).
//
// Java: Keywords.getTokenType(String)
func Keywords_GetTokenType(value string) TokenType {
	return keywords_KEYWORDS_BY_NAME[value].Type()
}

// Keywords_Get returns the keyword named value, or nil.
//
// Java: Keywords.get(String)
func Keywords_Get(value string) *Keywords {
	return keywords_KEYWORDS_BY_NAME[value]
}

// Keywords_GetByTokenType returns the keyword of token, or nil.
//
// Java: Keywords.get(TokenType)
func Keywords_GetByTokenType(token TokenType) *Keywords {
	return keywords_KEYWORDS_BY_TYPE[token]
}
