// Port of com/google/javascript/jscomp/parsing/parser/trees/NewTargetExpressionTree.java @ closure-compiler v20261015

/*
 * Copyright 2016 The Closure Compiler Authors.
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

// Represents 'new.target'.
type NewTargetExpressionTree struct {
	ParseTreeBase
}

// NewNewTargetExpressionTree constructs a NewTargetExpressionTree.
//
// Java: NewTargetExpressionTree(SourceRange)
func NewNewTargetExpressionTree(location SourceRange) *NewTargetExpressionTree {
	t := &NewTargetExpressionTree{}
	t.InitParseTreeBase(t, ParseTreeType_NEW_TARGET_EXPRESSION, location)
	return t
}
