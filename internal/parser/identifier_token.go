// Port of com/google/javascript/jscomp/parsing/parser/IdentifierToken.java @ closure-compiler v20260415

package parser

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// IdentifierToken is a token representing an identifier.
type IdentifierToken struct {
	TokenBase

	value             string
	privateIdentifier bool
}

// NewIdentifierToken constructs an IdentifierToken.
//
// Java: IdentifierToken(SourceRange, String)
func NewIdentifierToken(location *SourceRange, value string) *IdentifierToken {
	t := &IdentifierToken{}
	t.InitTokenBase(t, TokenType_IDENTIFIER, location)
	t.value = value
	t.privateIdentifier = strings.HasPrefix(value, "#")
	return t
}

// String returns the identifier's value.
//
// Java: IdentifierToken.toString()
func (t *IdentifierToken) String() string {
	return t.value
}

// IsKeyword reports whether the value is a keyword.
//
// Java: IdentifierToken.isKeyword()
func (t *IdentifierToken) IsKeyword() bool {
	return Keywords_IsKeyword(t.value)
}

// ValueEquals reports whether the value equals str.
//
// Java: IdentifierToken.valueEquals(String)
func (t *IdentifierToken) ValueEquals(str string) bool {
	return t.value == str
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
	return t.value
}

// GetMaybePrivateValue gets the value of the identifier, allowing it to be
// a private identifier.
//
// Java: IdentifierToken.getMaybePrivateValue()
func (t *IdentifierToken) GetMaybePrivateValue() string {
	return t.value
}

// IsPrivateIdentifier reports whether the value starts with a #.
//
// Java: IdentifierToken.isPrivateIdentifier()
func (t *IdentifierToken) IsPrivateIdentifier() bool {
	return t.privateIdentifier
}
