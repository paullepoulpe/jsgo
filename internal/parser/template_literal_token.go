// Port of com/google/javascript/jscomp/parsing/parser/TemplateLiteralToken.java @ closure-compiler v20261015

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

package parser

// TemplateLiteralToken is a token representing a javascript template
// literal substring.
//
// The value of the Token is the raw string. The token also stores whether
// this token contains any error messages that should be passed to the parser
// due to invalid escapes or unnecessary escapes. The parser, not the
// scanner, reports these errors because the errors are suppressed in tagged
// template literals. The scanner does not know if it's in a tagged or
// untagged template lit.
type TemplateLiteralToken struct {
	LiteralTokenBase

	// Go-only: language-adapter — nil is Java null.
	ErrorMessage  *string
	ErrorPosition *SourcePosition
	// Go-only: language-adapter — nil is Java null (no error).
	ErrorLevel *TemplateLiteralToken_ErrorLevel
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
	return []TemplateLiteralToken_ErrorLevel{TemplateLiteralToken_ErrorLevel_WARNING, TemplateLiteralToken_ErrorLevel_ERROR}
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

// NewTemplateLiteralToken constructs a TemplateLiteralToken.
//
// Java: TemplateLiteralToken(TokenType, String, String, ErrorLevel, SourcePosition, SourceRange)
func NewTemplateLiteralToken(type_ TokenType, value string, errorMsg *string, errorLevel *TemplateLiteralToken_ErrorLevel, position *SourcePosition, location *SourceRange) *TemplateLiteralToken {
	t := &TemplateLiteralToken{}
	t.InitLiteralTokenBase(t, type_, value, location)
	t.ErrorMessage = errorMsg
	t.ErrorLevel = errorLevel
	t.ErrorPosition = position
	return t
}

// String returns the raw string.
//
// Java: TemplateLiteralToken.toString()
func (t *TemplateLiteralToken) String() string {
	return t.Value
}

// HasError reports whether the token carries an error (not a warning).
//
// Java: TemplateLiteralToken.hasError()
func (t *TemplateLiteralToken) HasError() bool {
	return t.ErrorLevel != nil && *t.ErrorLevel == TemplateLiteralToken_ErrorLevel_ERROR
}
