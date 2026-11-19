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
// Port of com/google/javascript/jscomp/parsing/parser/trees/GetAccessorTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.GetAccessorTree
type GetAccessorTree struct {
	ParseTreeBase
	propertyName  Token      `gson:"propertyName"`
	body          *BlockTree `gson:"body"`
	isClassMember bool       `gson:"isClassMember"`
	isStatic      bool       `gson:"isStatic"`
}

func (this_ *GetAccessorTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *GetAccessorTree) IsJavaNil_() bool { return this_ == nil }
func (*GetAccessorTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.GetAccessorTree", "GetAccessorTree"
}
func (*GetAccessorTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.GetAccessorTree", "com.google.javascript.jscomp.parsing.parser.trees.GetAccessorTree", func(v_ any) bool {
		return j2gort.InstanceOf[*GetAccessorTree](v_)
	})
}
func (*GetAccessorTree) Implements_parser_ParseTree()       {}
func (this_ *GetAccessorTree) FieldGet_propertyName() Token { return this_.propertyName }
func (this_ *GetAccessorTree) FieldGet_isClassMember() bool { return this_.isClassMember }
func (this_ *GetAccessorTree) FieldGet_isStatic() bool      { return this_.isStatic }
func (this_ *GetAccessorTree) FieldGet_body() *BlockTree    { return this_.body }

func (this_ *GetAccessorTree) initFields_GetAccessorTree() {
	_ = this_
}

// Java constructor: GetAccessorTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.Token propertyName, boolean isClassMember, boolean isStatic, com.google.javascript.jscomp.parsing.parser.trees.BlockTree body)
func NewGetAccessorTree(location SourceRange, propertyName Token, isClassMember bool, isStatic bool, body *BlockTree) *GetAccessorTree {
	this_ := &GetAccessorTree{}
	this_.BindSelf_(this_)
	InitGetAccessorTree(this_, location, propertyName, isClassMember, isStatic, body)
	return this_
}
func allocPooled_GetAccessorTree(slab_ *j2gort.PooledSlab[GetAccessorTree], location SourceRange, propertyName Token, isClassMember bool, isStatic bool, body *BlockTree) *GetAccessorTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitGetAccessorTree(this_, location, propertyName, isClassMember, isStatic, body)
	return this_
}

func InitGetAccessorTree(this_ *GetAccessorTree, location SourceRange, propertyName Token, isClassMember bool, isStatic bool, body *BlockTree) {
	_ = location
	_ = propertyName
	_ = isClassMember
	_ = isStatic
	_ = body
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_GET_ACCESSOR(), location)
	this_.initFields_GetAccessorTree()
	this_.propertyName = propertyName
	this_.isClassMember = isClassMember
	this_.isStatic = isStatic
	this_.body = body
}
