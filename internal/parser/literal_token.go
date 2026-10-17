// Port of com/google/javascript/jscomp/parsing/parser/LiteralToken.java @ closure-compiler v20260415

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

// LiteralToken is a token representing a javascript literal. Includes
// string, regexp, and number literals. Boolean and null literals are
// represented as regular keyword tokens.
//
// The value just includes the raw lexeme. For string literals it includes
// the beginning and ending delimiters.
//
// TODO: Regexp literals should have their own token type.
// TODO: A way to get the processed value, rather than the raw value.
//
// Java: class LiteralToken (with subclasses StringLiteralToken and
// TemplateLiteralToken).
type LiteralToken interface {
	Token
	GetLiteralTokenBase() *LiteralTokenBase
}

// LiteralTokenBase holds LiteralToken's fields; plain literal tokens are
// *LiteralTokenBase values.
type LiteralTokenBase struct {
	TokenBase

	Value string
}

// NewLiteralToken constructs a plain LiteralToken.
//
// Java: LiteralToken(TokenType, String, SourceRange)
func NewLiteralToken(type_ TokenType, value string, location *SourceRange) *LiteralTokenBase {
	t := &LiteralTokenBase{}
	t.InitLiteralTokenBase(t, type_, value, location)
	return t
}

// InitLiteralTokenBase is the Go form of the LiteralToken constructor for
// subclasses.
//
// Java: LiteralToken(TokenType, String, SourceRange)
func (t *LiteralTokenBase) InitLiteralTokenBase(self LiteralToken, type_ TokenType, value string, location *SourceRange) {
	t.InitTokenBase(self, type_, location)
	t.Value = value
}

// GetLiteralTokenBase returns the base struct (field access through the
// interface).
//
// Go-only: language-adapter — field access through the interface (PORTING.md §4).
func (t *LiteralTokenBase) GetLiteralTokenBase() *LiteralTokenBase {
	return t
}

// String returns the raw lexeme.
//
// Java: LiteralToken.toString()
func (t *LiteralTokenBase) String() string {
	return t.Value
}
