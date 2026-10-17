// Port of com/google/javascript/jscomp/parsing/parser/trees/SetAccessorTree.java @ closure-compiler v20260415

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

// SetAccessorTree is the parse tree for ParseTreeType SET_ACCESSOR.
type SetAccessorTree struct {
	ParseTreeBase

	PropertyName  Token
	Parameter     *FormalParameterListTree
	IsClassMember bool
	IsStatic      bool
	Body          *BlockTree
}

// NewSetAccessorTree constructs a SetAccessorTree.
//
// Java: SetAccessorTree(SourceRange, Token, boolean, boolean, FormalParameterListTree, BlockTree)
func NewSetAccessorTree(location *SourceRange, propertyName Token, isClassMember bool, isStatic bool, parameter *FormalParameterListTree, body *BlockTree) *SetAccessorTree {
	t := &SetAccessorTree{}
	t.InitParseTreeBase(t, ParseTreeType_SET_ACCESSOR, location)
	t.PropertyName = propertyName
	t.IsClassMember = isClassMember
	t.IsStatic = isStatic
	t.Parameter = parameter
	t.Body = body
	return t
}
