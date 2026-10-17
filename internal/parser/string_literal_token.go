// Port of com/google/javascript/jscomp/parsing/parser/StringLiteralToken.java @ closure-compiler v20260415

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

// StringLiteralToken is a single or double quoted JavaScript string literal.
type StringLiteralToken struct {
	LiteralTokenBase

	hasUnescapedUnicodeLineOrParagraphSeparator bool
}

// NewStringLiteralToken constructs a StringLiteralToken.
//
// Java: StringLiteralToken(String, SourceRange, boolean)
func NewStringLiteralToken(value string, location *SourceRange, hasUnescapedUnicodeLineOrParagraphSeparator bool) *StringLiteralToken {
	t := &StringLiteralToken{}
	t.InitLiteralTokenBase(t, TokenType_STRING, value, location)
	t.hasUnescapedUnicodeLineOrParagraphSeparator = hasUnescapedUnicodeLineOrParagraphSeparator
	return t
}

// HasUnescapedUnicodeLineOrParagraphSeparator reports whether the literal
// contains an unescaped U+2028 or U+2029.
//
// Java: StringLiteralToken.hasUnescapedUnicodeLineOrParagraphSeparator()
func (t *StringLiteralToken) HasUnescapedUnicodeLineOrParagraphSeparator() bool {
	return t.hasUnescapedUnicodeLineOrParagraphSeparator
}
