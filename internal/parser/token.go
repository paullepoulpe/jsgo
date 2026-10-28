// Port of com/google/javascript/jscomp/parsing/parser/Token.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/IdentifierToken.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/LiteralToken.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/StringLiteralToken.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/parsing/parser/TemplateLiteralToken.java @ closure-compiler v20261021

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

/*
 * Copyright 2018 The Closure Compiler Authors.
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

/*
 * Copyright 2019 The Closure Compiler Authors.
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

// Token flags packed into Token.flags.
const (
	tokenFlagPrivateIdent               uint8 = 1 << 0
	tokenFlagIdentIsKeyword             uint8 = 1 << 1
	tokenFlagUnescapedLineOrParagraphSep uint8 = 1 << 2
	tokenFlagTemplateTrim1              uint8 = 1 << 3 // template closed with `
	tokenFlagTemplateTrim2              uint8 = 1 << 4 // template closed with ${
	tokenFlagNumericSeparator           uint8 = 1 << 5 // _ in NUMBER or BIGINT
	tokenFlagTemplateHasError           uint8 = 1 << 6 // TemplateLiteralToken.ErrorLevel == ERROR
)

// Token is a 16-byte pointer-free value representing any JavaScript token
// (punctuation, keyword, identifier, number, bigint, string, regex, or
// template literal portion). The zero value Token{} represents an absent
// (Java null) token.
//
// Java: Token, IdentifierToken, LiteralToken, StringLiteralToken, TemplateLiteralToken
type Token struct {
	Location SourceRange // 8 bytes (Start, End SourcePosition)
	identID  strpool.ID  // 4 bytes (interned identifier ID, or standalone literal fallback)
	aux      uint16      // 2 bytes (1-based index into Scanner.templateErrors when != 0)
	rawType  uint8       // 1 byte  (0 = null/absent; uint8(TokenType)+1 otherwise)
	flags    uint8       // 1 byte
}

// TemplateLiteralToken_ErrorLevel is TemplateLiteralToken.ErrorLevel.
type TemplateLiteralToken_ErrorLevel int32

const (
	TemplateLiteralToken_ErrorLevel_WARNING TemplateLiteralToken_ErrorLevel = iota
	TemplateLiteralToken_ErrorLevel_ERROR
)

var templateLiteralToken_ErrorLevel_names = [...]string{"WARNING", "ERROR"}

// TemplateLiteralToken_ErrorLevel_values is ErrorLevel.values().
//
// Java: TemplateLiteralToken.ErrorLevel.values()
func TemplateLiteralToken_ErrorLevel_values() []TemplateLiteralToken_ErrorLevel {
	return []TemplateLiteralToken_ErrorLevel{
		TemplateLiteralToken_ErrorLevel_WARNING,
		TemplateLiteralToken_ErrorLevel_ERROR,
	}
}

// Name is Enum.name().
//
// Java: TemplateLiteralToken.ErrorLevel.name()
func (e TemplateLiteralToken_ErrorLevel) Name() string {
	return templateLiteralToken_ErrorLevel_names[e]
}

// String is Enum.toString().
//
// Java: TemplateLiteralToken.ErrorLevel.toString()
func (e TemplateLiteralToken_ErrorLevel) String() string { return e.Name() }

// TemplateTokenError stores the rare diagnostic or non-start position metadata
// of a TemplateLiteralToken in a side slice on Scanner (indexed 1-based by
// Token.aux).
type TemplateTokenError struct {
	Message  string
	Position SourcePosition
	HasMsg   bool
	Level    uint8 // 0 = null, 1 = WARNING, 2 = ERROR
}

// ErrorLevelPtr returns the nullable ErrorLevel of this entry.
func (e TemplateTokenError) ErrorLevelPtr() *TemplateLiteralToken_ErrorLevel {
	if e.Level == 0 {
		return nil
	}
	lvl := TemplateLiteralToken_ErrorLevel(e.Level - 1)
	return &lvl
}

// ErrorMessagePtr returns the nullable ErrorMessage of this entry.
func (e TemplateTokenError) ErrorMessagePtr() *string {
	if !e.HasMsg {
		return nil
	}
	msg := e.Message
	return &msg
}

// Standalone template diagnostic table for NewTemplateLiteralToken in tests
// (never used by Scanner, which stores diagnostics in Scanner.templateErrors).
var standaloneTemplateErrors []TemplateTokenError

// NewToken constructs a plain Token.
//
// Java: Token(TokenType, SourceRange)
func NewToken(type_ TokenType, location SourceRange) Token {
	return Token{
		Location: location,
		rawType:  uint8(type_) + 1,
	}
}

// NewIdentifierToken constructs an IdentifierToken.
//
// Java: IdentifierToken(SourceRange, String)
func NewIdentifierToken(location SourceRange, value string) Token {
	var flags uint8
	if len(value) > 0 && value[0] == '#' {
		flags |= tokenFlagPrivateIdent
	} else if Keywords_IsKeyword(value) {
		flags |= tokenFlagIdentIsKeyword
	}
	return Token{
		Location: location,
		identID:  strpool.Intern(value),
		rawType:  uint8(TokenType_IDENTIFIER) + 1,
		flags:    flags,
	}
}

func newIdentifierTokenFromID(location SourceRange, id strpool.ID, flags uint8) Token {
	return Token{
		Location: location,
		identID:  id,
		rawType:  uint8(TokenType_IDENTIFIER) + 1,
		flags:    flags,
	}
}

// NewLiteralToken constructs a standalone LiteralToken (used in tests; Scanner
// uses newScannerLiteralToken to avoid slicing/interning raw lexemes).
//
// Java: LiteralToken(TokenType, String, SourceRange)
func NewLiteralToken(type_ TokenType, value string, location SourceRange) Token {
	var flags uint8
	if (type_ == TokenType_NUMBER || type_ == TokenType_BIGINT) && strings.IndexByte(value, '_') >= 0 {
		flags |= tokenFlagNumericSeparator
	}
	return Token{
		Location: location,
		identID:  strpool.Intern(value),
		rawType:  uint8(type_) + 1,
		flags:    flags,
	}
}

// NewStringLiteralToken constructs a standalone StringLiteralToken.
//
// Java: StringLiteralToken(String, SourceRange, boolean)
func NewStringLiteralToken(value string, location SourceRange, hasUnescapedUnicodeLineOrParagraphSeparator bool) Token {
	var flags uint8
	if hasUnescapedUnicodeLineOrParagraphSeparator {
		flags |= tokenFlagUnescapedLineOrParagraphSep
	}
	return Token{
		Location: location,
		identID:  strpool.Intern(value),
		rawType:  uint8(TokenType_STRING) + 1,
		flags:    flags,
	}
}

// NewTemplateLiteralToken constructs a standalone TemplateLiteralToken.
//
// Java: TemplateLiteralToken(TokenType, String, String, ErrorLevel, SourcePosition, SourceRange)
func NewTemplateLiteralToken(
	type_ TokenType,
	value string,
	errorMsg *string,
	errorLevel *TemplateLiteralToken_ErrorLevel,
	position SourcePosition,
	location SourceRange,
) Token {
	var flags uint8
	if errorLevel != nil && *errorLevel == TemplateLiteralToken_ErrorLevel_ERROR {
		flags |= tokenFlagTemplateHasError
	}
	var aux uint16
	if errorMsg != nil || errorLevel != nil || position != location.Start {
		var msg string
		if errorMsg != nil {
			msg = *errorMsg
		}
		var lvl uint8
		if errorLevel != nil {
			lvl = uint8(*errorLevel) + 1
		}
		standaloneTemplateErrors = append(standaloneTemplateErrors, TemplateTokenError{
			Message:  msg,
			Position: position,
			HasMsg:   errorMsg != nil,
			Level:    lvl,
		})
		aux = uint16(len(standaloneTemplateErrors))
	}
	return Token{
		Location: location,
		identID:  strpool.Intern(value),
		aux:      aux,
		rawType:  uint8(type_) + 1,
		flags:    flags,
	}
}

func newScannerLiteralToken(type_ TokenType, location SourceRange, flags uint8) Token {
	return Token{
		Location: location,
		rawType:  uint8(type_) + 1,
		flags:    flags,
	}
}

func newScannerTemplateToken(type_ TokenType, location SourceRange, flags uint8, aux uint16) Token {
	return Token{
		Location: location,
		aux:      aux,
		rawType:  uint8(type_) + 1,
		flags:    flags,
	}
}

// IsNull reports whether t is the zero (absent / Java null) Token.
func (t Token) IsNull() bool {
	return t.rawType == 0
}

// IsPresent reports whether t is non-null.
func (t Token) IsPresent() bool {
	return t.rawType != 0
}

// Type returns the token's TokenType.
func (t Token) Type() TokenType {
	return TokenType(t.rawType - 1)
}

// GetStart returns the start position.
//
// Java: Token.getStart()
func (t Token) GetStart() SourcePosition {
	return t.Location.Start
}

// GetEnd returns the end position.
func (t Token) GetEnd() SourcePosition {
	return t.Location.End
}

// IsPrivateIdentifier reports whether an identifier token starts with '#'.
//
// Java: IdentifierToken.isPrivateIdentifier()
func (t Token) IsPrivateIdentifier() bool {
	return t.flags&tokenFlagPrivateIdent != 0
}

// IsKeyword reports whether an identifier token's value is a keyword.
//
// Java: IdentifierToken.isKeyword()
func (t Token) IsKeyword() bool {
	return t.flags&tokenFlagIdentIsKeyword != 0
}

// HasUnescapedUnicodeLineOrParagraphSeparator reports whether a string literal
// token contains an unescaped U+2028 or U+2029.
//
// Java: StringLiteralToken.hasUnescapedUnicodeLineOrParagraphSeparator()
func (t Token) HasUnescapedUnicodeLineOrParagraphSeparator() bool {
	return t.flags&tokenFlagUnescapedLineOrParagraphSep != 0
}

// HasNumericSeparator reports whether a number or bigint token contains '_'.
func (t Token) HasNumericSeparator() bool {
	return t.flags&tokenFlagNumericSeparator != 0
}

// HasError reports whether a template literal token carries an ERROR (not a
// WARNING).
//
// Java: TemplateLiteralToken.hasError()
func (t Token) HasError() bool {
	return t.flags&tokenFlagTemplateHasError != 0
}

// IsTemplateLiteral reports whether t is a template literal token.
func (t Token) IsTemplateLiteral() bool {
	switch t.Type() {
	case TokenType_NO_SUBSTITUTION_TEMPLATE,
		TokenType_TEMPLATE_HEAD,
		TokenType_TEMPLATE_MIDDLE,
		TokenType_TEMPLATE_TAIL:
		return true
	default:
		return false
	}
}

// ValueEquals reports whether an identifier token's value equals str.
//
// Java: IdentifierToken.valueEquals(String)
func (t Token) ValueEquals(str string) bool {
	return strpool.String(t.identID) == str
}

// GetValue gets the value of an identifier token, assuring that it is not a
// private identifier.
//
// Java: IdentifierToken.getValue()
func (t Token) GetValue() string {
	jutil.CheckState(!t.IsPrivateIdentifier())
	return strpool.String(t.identID)
}

// ValueID returns the interned string pool ID of a non-private identifier.
func (t Token) ValueID() strpool.ID {
	jutil.CheckState(!t.IsPrivateIdentifier())
	return t.identID
}

// MaybePrivateValueID returns the interned string pool ID of an identifier,
// allowing private identifiers.
func (t Token) MaybePrivateValueID() strpool.ID {
	return t.identID
}

// GetMaybePrivateValue gets the value of an identifier token, allowing private
// identifiers.
//
// Java: IdentifierToken.getMaybePrivateValue()
func (t Token) GetMaybePrivateValue() string {
	return strpool.String(t.identID)
}

// ValueRange returns the UTF-16 [start, end) range of a literal or template
// token's raw value in its SourceFile (for template tokens, excluding the
// closing '`' or '${' delimiter).
func (t Token) ValueRange() (start, end int) {
	start = int(t.Location.Start.Offset)
	end = int(t.Location.End.Offset)
	if t.flags&tokenFlagTemplateTrim1 != 0 {
		end--
	} else if t.flags&tokenFlagTemplateTrim2 != 0 {
		end -= 2
	}
	return start, end
}

// RawValue returns the raw lexeme string for an identifier, literal, string,
// or template literal token.
func (t Token) RawValue(file *SourceFile) string {
	if t.identID != 0 {
		return strpool.String(t.identID)
	}
	start, end := t.ValueRange()
	return file.Substring(start, end)
}

// JavaSimpleClassName returns the simple name of the corresponding Java Token
// class (Token, IdentifierToken, LiteralToken, StringLiteralToken, or
// TemplateLiteralToken).
func (t Token) JavaSimpleClassName() string {
	switch t.Type() {
	case TokenType_IDENTIFIER:
		return "IdentifierToken"
	case TokenType_STRING:
		return "StringLiteralToken"
	case TokenType_NO_SUBSTITUTION_TEMPLATE,
		TokenType_TEMPLATE_HEAD,
		TokenType_TEMPLATE_MIDDLE,
		TokenType_TEMPLATE_TAIL:
		return "TemplateLiteralToken"
	case TokenType_NUMBER,
		TokenType_BIGINT,
		TokenType_REGULAR_EXPRESSION:
		return "LiteralToken"
	default:
		return "Token"
	}
}

// String returns the token's string representation.
//
// Java: Token.toString() / IdentifierToken.toString() / LiteralToken.toString()
func (t Token) String() string {
	if t.identID != 0 {
		return strpool.String(t.identID)
	}
	return t.Type().String()
}
