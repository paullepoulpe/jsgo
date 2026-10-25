// Port of com/google/javascript/jscomp/parsing/parser/IdentifierToken.java @ closure-compiler v20261021

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

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// IdentifierToken is a token representing an identifier.
type IdentifierToken struct {
	TokenBase
	// Go-only: performance — before the other fields, in the base's padding.
	privateIdentifier bool

	// Go-only: performance — the value interned in the string pool (4
	// bytes instead of a 16-byte string header, in hundreds of thousands
	// of tokens; the AST interns the same names anyway).
	value strpool.ID
}

// NewIdentifierToken constructs an IdentifierToken.
//
// Java: IdentifierToken(SourceRange, String)
func NewIdentifierToken(location SourceRange, value string) *IdentifierToken {
	t := &IdentifierToken{}
	t.initIdentifierToken(location, value)
	return t
}

func (t *IdentifierToken) initIdentifierToken(location SourceRange, value string) {
	t.InitTokenBase(TokenType_IDENTIFIER, location)
	t.value = strpool.Intern(value)
	t.privateIdentifier = strings.HasPrefix(value, "#")
}

// String returns the identifier's value.
//
// Java: IdentifierToken.toString()
func (t *IdentifierToken) String() string {
	return strpool.String(t.value)
}

// IsKeyword reports whether the value is a keyword.
//
// Java: IdentifierToken.isKeyword()
func (t *IdentifierToken) IsKeyword() bool {
	return Keywords_IsKeyword(strpool.String(t.value))
}

// ValueEquals reports whether the value equals str.
//
// Java: IdentifierToken.valueEquals(String)
func (t *IdentifierToken) ValueEquals(str string) bool {
	return strpool.String(t.value) == str
}

// GetValue gets the value of the identifier assuring that it is not a
// private identifier.
//
// You must verify privateIdentifier is false (and presumably error if it is
// true) before calling this method.
//
// Prefer calling IsKeyword() or ValueEquals(String) if those methods meet
// your needs.
//
// Java: IdentifierToken.getValue()
func (t *IdentifierToken) GetValue() string {
	jutil.CheckState(!t.privateIdentifier)
	return strpool.String(t.value)
}

// ValueID returns the value's ID in the string pool (see GetValue).
//
// Go-only: performance — the IR factory makes nodes from it without
// interning the name again.
func (t *IdentifierToken) ValueID() strpool.ID {
	jutil.CheckState(!t.privateIdentifier)
	return t.value
}

// MaybePrivateValueID returns the ID in the string pool of
// GetMaybePrivateValue.
//
// Go-only: performance (see ValueID).
func (t *IdentifierToken) MaybePrivateValueID() strpool.ID {
	return t.value
}

// GetMaybePrivateValue gets the value of the identifier, allowing it to be
// a private identifier.
//
// Java: IdentifierToken.getMaybePrivateValue()
func (t *IdentifierToken) GetMaybePrivateValue() string {
	return strpool.String(t.value)
}

// IsPrivateIdentifier reports whether the value starts with a #.
//
// Java: IdentifierToken.isPrivateIdentifier()
func (t *IdentifierToken) IsPrivateIdentifier() bool {
	return t.privateIdentifier
}
