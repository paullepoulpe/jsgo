// Port of com/google/javascript/jscomp/parsing/parser/trees/TemplateLiteralPortionTree.java @ closure-compiler v20261015

/*
 * Copyright 2014 The Closure Compiler Authors.
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

// A production representing a literal portion of a template literal.
// Eg: `hello`.
type TemplateLiteralPortionTree struct {
	ParseTreeBase

	Value Token
}

// NewTemplateLiteralPortionTree constructs a TemplateLiteralPortionTree.
//
// Java: TemplateLiteralPortionTree(SourceRange, Token)
func NewTemplateLiteralPortionTree(location SourceRange, value Token) *TemplateLiteralPortionTree {
	t := &TemplateLiteralPortionTree{}
	t.InitParseTreeBase(t, ParseTreeType_TEMPLATE_LITERAL_PORTION, location)
	t.Value = value
	return t
}
