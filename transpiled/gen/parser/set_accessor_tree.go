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
// Port of com/google/javascript/jscomp/parsing/parser/trees/SetAccessorTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.SetAccessorTree
type SetAccessorTree struct {
	ParseTreeBase
	propertyName  Token                    `gson:"propertyName"`
	parameter     *FormalParameterListTree `gson:"parameter"`
	body          *BlockTree               `gson:"body"`
	isClassMember bool                     `gson:"isClassMember"`
	isStatic      bool                     `gson:"isStatic"`
}

func (this_ *SetAccessorTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *SetAccessorTree) IsJavaNil_() bool { return this_ == nil }
func (*SetAccessorTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.SetAccessorTree", "SetAccessorTree"
}
func (*SetAccessorTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.SetAccessorTree", "com.google.javascript.jscomp.parsing.parser.trees.SetAccessorTree", func(v_ any) bool {
		return j2gort.InstanceOf[*SetAccessorTree](v_)
	})
}
func (*SetAccessorTree) Implements_parser_ParseTree()                       {}
func (this_ *SetAccessorTree) FieldGet_propertyName() Token                 { return this_.propertyName }
func (this_ *SetAccessorTree) FieldGet_parameter() *FormalParameterListTree { return this_.parameter }
func (this_ *SetAccessorTree) FieldGet_isClassMember() bool                 { return this_.isClassMember }
func (this_ *SetAccessorTree) FieldGet_isStatic() bool                      { return this_.isStatic }
func (this_ *SetAccessorTree) FieldGet_body() *BlockTree                    { return this_.body }

func (this_ *SetAccessorTree) initFields_SetAccessorTree() {
	_ = this_
}

// Java constructor: SetAccessorTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.Token propertyName, boolean isClassMember, boolean isStatic, com.google.javascript.jscomp.parsing.parser.trees.FormalParameterListTree parameter, com.google.javascript.jscomp.parsing.parser.trees.BlockTree body)
func NewSetAccessorTree(location SourceRange, propertyName Token, isClassMember bool, isStatic bool, parameter *FormalParameterListTree, body *BlockTree) *SetAccessorTree {
	this_ := &SetAccessorTree{}
	this_.BindSelf_(this_)
	InitSetAccessorTree(this_, location, propertyName, isClassMember, isStatic, parameter, body)
	return this_
}
func allocPooled_SetAccessorTree(slab_ *j2gort.PooledSlab[SetAccessorTree], location SourceRange, propertyName Token, isClassMember bool, isStatic bool, parameter *FormalParameterListTree, body *BlockTree) *SetAccessorTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitSetAccessorTree(this_, location, propertyName, isClassMember, isStatic, parameter, body)
	return this_
}

func InitSetAccessorTree(this_ *SetAccessorTree, location SourceRange, propertyName Token, isClassMember bool, isStatic bool, parameter *FormalParameterListTree, body *BlockTree) {
	_ = location
	_ = propertyName
	_ = isClassMember
	_ = isStatic
	_ = parameter
	_ = body
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_SET_ACCESSOR(), location)
	this_.initFields_SetAccessorTree()
	this_.propertyName = propertyName
	this_.isClassMember = isClassMember
	this_.isStatic = isStatic
	this_.parameter = parameter
	this_.body = body
}
