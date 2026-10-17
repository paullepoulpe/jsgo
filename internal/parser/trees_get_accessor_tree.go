// Port of com/google/javascript/jscomp/parsing/parser/trees/GetAccessorTree.java @ closure-compiler v20261015

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

// GetAccessorTree is the parse tree for ParseTreeType GET_ACCESSOR.
type GetAccessorTree struct {
	ParseTreeBase

	PropertyName  Token
	IsClassMember bool
	IsStatic      bool
	Body          *BlockTree
}

// NewGetAccessorTree constructs a GetAccessorTree.
//
// Java: GetAccessorTree(SourceRange, Token, boolean, boolean, BlockTree)
func NewGetAccessorTree(location *SourceRange, propertyName Token, isClassMember bool, isStatic bool, body *BlockTree) *GetAccessorTree {
	t := &GetAccessorTree{}
	t.InitParseTreeBase(t, ParseTreeType_GET_ACCESSOR, location)
	t.PropertyName = propertyName
	t.IsClassMember = isClassMember
	t.IsStatic = isStatic
	t.Body = body
	return t
}
