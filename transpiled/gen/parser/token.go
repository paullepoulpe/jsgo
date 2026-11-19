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
// Port of com/google/javascript/jscomp/parsing/parser/Token.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.parsing.parser.Token
type Token struct {
	value string `gson:"value"`
	templateError *TemplateLiteralToken `gson:"templateError"`
	location SourceRange `gson:"location"`
	type_ TokenType `gson:"type"`
	tokenFlags int8 `gson:"tokenFlags"`
}

func (this_ Token) BindSelf_(s any) {
	_ = s
}
func (this_ Token) IsJavaNil_() bool { return this_.type_ == 0 }
func (Token) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.parsing.parser.Token", "Token" }
func (Token) HasIdentityEquals_() {}
func (this_ Token) IdentityHash_() int32 {
	if this_.type_ == 0 {
		return 0
	}
	var h_ int32 = 1
	h_ = 31*h_ + (this_.location).IdentityHash_()
	h_ = 31*h_ + (this_.type_).HashCode()
	{ hs_3 := this_.value; h_ = 31 * h_; if hs_3 != j2gort.NullString { h_ += javalang.StringHashCode(hs_3) } }
	{ hp_4 := this_.templateError; h_ = 31 * h_; if hp_4 != nil { h_ += j2gort.IdentityHash(hp_4) } }
	h_ = 31*h_ + int32(this_.tokenFlags)
	return h_
}
func (this_ Token) FastValHash_() (int32, bool, bool) {
	if this_.IsJavaNil_() {
		return 0, false, true
	}
	return this_.IdentityHash_(), true, true
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.Token", "com.google.javascript.jscomp.parsing.parser.Token", func(v_ any) bool {
		return j2gort.InstanceOf[Token](v_)
	})
}
func (this_ Token) FieldGet_location() SourceRange { return this_.location }
func (this_ Token) FieldGet_type_() TokenType { return this_.type_ }
func (this_ Token) FieldGet_value() string { return this_.value }

func StaticGet_Token_FLAG_IDENTIFIER() int8 { return int8(1) }
func StaticGet_Token_FLAG_LITERAL() int8 { return int8(2) }
func StaticGet_Token_FLAG_STRING_LITERAL() int8 { return int8(4) }
func StaticGet_Token_FLAG_TEMPLATE_LITERAL() int8 { return int8(8) }
func StaticGet_Token_FLAG_PRIVATE_IDENTIFIER() int8 { return int8(16) }
func StaticGet_Token_FLAG_UNESCAPED_UNICODE_LINE_OR_PARAGRAPH_SEP() int8 { return int8(32) }
func (this_ *Token) initFields_Token() {
	_ = this_
	this_.value = j2gort.NullString
}

// Java constructor: Token(com.google.javascript.jscomp.parsing.parser.TokenType type, com.google.javascript.jscomp.parsing.parser.util.SourceRange location)
func NewToken_2(type_ TokenType, location SourceRange) Token {
	this_ := Token{}
	InitToken_2(&this_, type_, location)
	return this_
}

func InitToken_2(this_ *Token, type_ TokenType, location SourceRange) {
	_ = type_
	_ = location
	this_.initFields_Token()
	this_.type_ = j2gort.CheckNotNullEnum(type_)
	this_.location = location
	this_.value = ""
	this_.templateError = nil
	this_.tokenFlags = int8(int32(0))
}

// Java constructor: Token(com.google.javascript.jscomp.parsing.parser.TokenType type, com.google.javascript.jscomp.parsing.parser.util.SourceRange location, java.lang.String value, com.google.javascript.jscomp.parsing.parser.@org.jspecify.annotations.Nullable TemplateLiteralToken templateError, byte tokenFlags)
func NewToken_5(type_ TokenType, location SourceRange, value string, templateError *TemplateLiteralToken, tokenFlags int8) Token {
	this_ := Token{}
	InitToken_5(&this_, type_, location, value, templateError, tokenFlags)
	return this_
}

func InitToken_5(this_ *Token, type_ TokenType, location SourceRange, value string, templateError *TemplateLiteralToken, tokenFlags int8) {
	_ = type_
	_ = location
	_ = value
	_ = templateError
	_ = tokenFlags
	this_.initFields_Token()
	this_.type_ = j2gort.CheckNotNullEnum(type_)
	this_.location = location
	this_.value = j2gort.CheckNotNullString(value)
	this_.templateError = templateError
	this_.tokenFlags = tokenFlags
}

// Java static method: Token.identifier(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, java.lang.String value)
func Token_Identifier(location SourceRange, value string) Token {
	_ = location
	_ = value
	var flags int8 = int8(((int32(int8(1)) | int32((func() int8 { if javalang.StringStartsWith(j2gort.CheckNotNullString(value), "#") { return int8(16) }; return int8(int32(0)) }())))))
	_ = flags
	return NewToken_5(StaticGet_TokenType_IDENTIFIER(), location, value, nil, flags)
}

// Java static method: Token.literal(com.google.javascript.jscomp.parsing.parser.TokenType type, java.lang.String value, com.google.javascript.jscomp.parsing.parser.util.SourceRange location)
func Token_Literal(type_ TokenType, value string, location SourceRange) Token {
	_ = type_
	_ = value
	_ = location
	return NewToken_5(type_, location, value, nil, int8(2))
}

// Java static method: Token.stringLiteral(java.lang.String value, com.google.javascript.jscomp.parsing.parser.util.SourceRange location, boolean hasUnescapedUnicodeLineOrParagraphSeparator)
func Token_StringLiteral(value string, location SourceRange, hasUnescapedUnicodeLineOrParagraphSeparator bool) Token {
	_ = value
	_ = location
	_ = hasUnescapedUnicodeLineOrParagraphSeparator
	var flags int8 = int8((((int32(int8(2)) | int32(int8(4))) | int32((func() int8 { if hasUnescapedUnicodeLineOrParagraphSeparator { return int8(32) }; return int8(int32(0)) }())))))
	_ = flags
	return NewToken_5(StaticGet_TokenType_STRING(), location, value, nil, flags)
}

// Java static method: Token.templateLiteral(com.google.javascript.jscomp.parsing.parser.TokenType type, java.lang.String value, java.lang.@org.jspecify.annotations.Nullable String errorMsg, com.google.javascript.jscomp.parsing.parser.TemplateLiteralToken.@org.jspecify.annotations.Nullable ErrorLevel errorLevel, com.google.javascript.jscomp.parsing.parser.util.SourcePosition position, com.google.javascript.jscomp.parsing.parser.util.SourceRange location)
func Token_TemplateLiteral(type_ TokenType, value string, errorMsg string, errorLevel TemplateLiteralToken_ErrorLevel, position SourcePosition, location SourceRange) Token {
	_ = type_
	_ = value
	_ = errorMsg
	_ = errorLevel
	_ = position
	_ = location
	var templateError *TemplateLiteralToken
	_ = templateError
	if (((errorMsg == j2gort.NullString) && (errorLevel == 0))) {
		templateError = nil
	} else {
		templateError = NewTemplateLiteralToken(errorMsg, errorLevel, position)
	}
	return NewToken_5(type_, location, value, templateError, int8(10))
}

// Java method: Token.getStart()
func (this_ Token) GetStart() SourcePosition {
	if this_.type_ == 0 {
		j2gort.ThrowNPE()
	}
	return this_.location.start
}

// Java method: Token.isIdentifierToken()
func (this_ Token) IsIdentifierToken() bool {
	if this_.type_ == 0 {
		j2gort.ThrowNPE()
	}
	return (((int32(this_.tokenFlags) & int32(int8(1)))) != int32(0))
}

// Java method: Token.isLiteralToken()
func (this_ Token) IsLiteralToken() bool {
	if this_.type_ == 0 {
		j2gort.ThrowNPE()
	}
	return (((int32(this_.tokenFlags) & int32(int8(2)))) != int32(0))
}

// Java method: Token.isStringLiteralToken()
func (this_ Token) IsStringLiteralToken() bool {
	if this_.type_ == 0 {
		j2gort.ThrowNPE()
	}
	return (((int32(this_.tokenFlags) & int32(int8(4)))) != int32(0))
}

// Java method: Token.isTemplateLiteralToken()
func (this_ Token) IsTemplateLiteralToken() bool {
	if this_.type_ == 0 {
		j2gort.ThrowNPE()
	}
	return (((int32(this_.tokenFlags) & int32(int8(8)))) != int32(0))
}

// Java method: Token.asIdentifier()
func (this_ Token) AsIdentifier() Token {
	if this_.type_ == 0 {
		j2gort.ThrowNPE()
	}
	if !((((int32(this_.tokenFlags) & int32(int8(1)))) != int32(0))) { j2gort.CheckState(false) }
	return this_
}

// Java method: Token.asLiteral()
func (this_ Token) AsLiteral() Token {
	if this_.type_ == 0 {
		j2gort.ThrowNPE()
	}
	if !((((int32(this_.tokenFlags) & int32(int8(2)))) != int32(0))) { j2gort.CheckState(false) }
	return this_
}

// Java method: Token.isKeyword()
func (this_ Token) IsKeyword() bool {
	if this_.type_ == 0 {
		j2gort.ThrowNPE()
	}
	return Keywords_IsKeyword_1_string(this_.value)
}

// Java method: Token.valueEquals(java.lang.String str)
func (this_ Token) ValueEquals(str string) bool {
	if this_.type_ == 0 {
		j2gort.ThrowNPE()
	}
	_ = str
	return (j2gort.CheckNotNullString(this_.value) == str)
}

// Java method: Token.getValue()
func (this_ Token) GetValue() string {
	if this_.type_ == 0 {
		j2gort.ThrowNPE()
	}
	if !((((int32(this_.tokenFlags) & int32(int8(16)))) == int32(0))) { j2gort.CheckState(false) }
	return this_.value
}

// Java method: Token.getMaybePrivateValue()
func (this_ Token) GetMaybePrivateValue() string {
	if this_.type_ == 0 {
		j2gort.ThrowNPE()
	}
	return this_.value
}

// Java method: Token.isPrivateIdentifier()
func (this_ Token) IsPrivateIdentifier() bool {
	if this_.type_ == 0 {
		j2gort.ThrowNPE()
	}
	return (((int32(this_.tokenFlags) & int32(int8(16)))) != int32(0))
}

// Java method: Token.hasUnescapedUnicodeLineOrParagraphSeparator()
func (this_ Token) HasUnescapedUnicodeLineOrParagraphSeparator() bool {
	if this_.type_ == 0 {
		j2gort.ThrowNPE()
	}
	return (((int32(this_.tokenFlags) & int32(int8(32)))) != int32(0))
}

// Java method: Token.hasError()
func (this_ Token) HasError() bool {
	if this_.type_ == 0 {
		j2gort.ThrowNPE()
	}
	return ((this_.templateError != nil) && func() bool { re_1 := StaticGet_TemplateLiteralToken_ErrorLevel_ERROR(); ae_2 := this_.templateError.errorLevel; if re_1 == 0 { j2gort.ThrowNPE() }; return re_1 == ae_2 }())
}

// Java method: Token.getErrorMessage()
func (this_ Token) GetErrorMessage() string {
	if this_.type_ == 0 {
		j2gort.ThrowNPE()
	}
	if (this_.templateError == nil) {
		return j2gort.NullString
	} else {
		return this_.templateError.errorMessage
	}
}

// Java method: Token.getErrorPosition()
func (this_ Token) GetErrorPosition() SourcePosition {
	if this_.type_ == 0 {
		j2gort.ThrowNPE()
	}
	if (this_.templateError == nil) {
		return (SourcePosition{})
	} else {
		return this_.templateError.errorPosition
	}
}

// Java method: Token.getErrorLevel()
func (this_ Token) GetErrorLevel() TemplateLiteralToken_ErrorLevel {
	if this_.type_ == 0 {
		j2gort.ThrowNPE()
	}
	if (this_.templateError == nil) {
		return TemplateLiteralToken_ErrorLevel(0)
	} else {
		return this_.templateError.errorLevel
	}
}

// Java method: Token.toString()
func (this_ Token) String() string {
	if this_.type_ == 0 {
		j2gort.ThrowNPE()
	}
	if ((((int32(this_.tokenFlags) & ((int32(int8(1)) | int32(int8(2)))))) != int32(0))) {
		return this_.value
	}
	if (this_.type_ == 0) {
		return "null"
	} else {
		return this_.type_.String()
	}
}

