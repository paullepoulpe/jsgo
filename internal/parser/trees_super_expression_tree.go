// Port of com/google/javascript/jscomp/parsing/parser/trees/SuperExpressionTree.java @ closure-compiler v20260415

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

// SuperExpressionTree is the parse tree for ParseTreeType SUPER_EXPRESSION.
type SuperExpressionTree struct {
	ParseTreeBase
}

// NewSuperExpressionTree constructs a SuperExpressionTree.
//
// Java: SuperExpressionTree(SourceRange)
func NewSuperExpressionTree(location *SourceRange) *SuperExpressionTree {
	t := &SuperExpressionTree{}
	t.InitParseTreeBase(t, ParseTreeType_SUPER_EXPRESSION, location)
	return t
}
