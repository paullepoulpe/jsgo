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
// Port of com/google/javascript/jscomp/parsing/parser/trees/ComputedPropertySetterTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.ComputedPropertySetterTree
type ComputedPropertySetterTree struct {
	ParseTreeBase
	property  ParseTree                `gson:"property"`
	parameter *FormalParameterListTree `gson:"parameter"`
	body      *BlockTree               `gson:"body"`
	isStatic  bool                     `gson:"isStatic"`
}

func (this_ *ComputedPropertySetterTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *ComputedPropertySetterTree) IsJavaNil_() bool { return this_ == nil }
func (*ComputedPropertySetterTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.ComputedPropertySetterTree", "ComputedPropertySetterTree"
}
func (*ComputedPropertySetterTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.ComputedPropertySetterTree", "com.google.javascript.jscomp.parsing.parser.trees.ComputedPropertySetterTree", func(v_ any) bool {
		return j2gort.InstanceOf[*ComputedPropertySetterTree](v_)
	})
}
func (*ComputedPropertySetterTree) Implements_parser_ParseTree()       {}
func (this_ *ComputedPropertySetterTree) FieldGet_property() ParseTree { return this_.property }
func (this_ *ComputedPropertySetterTree) FieldGet_parameter() *FormalParameterListTree {
	return this_.parameter
}
func (this_ *ComputedPropertySetterTree) FieldGet_isStatic() bool   { return this_.isStatic }
func (this_ *ComputedPropertySetterTree) FieldGet_body() *BlockTree { return this_.body }

func (this_ *ComputedPropertySetterTree) initFields_ComputedPropertySetterTree() {
	_ = this_
}

// Java constructor: ComputedPropertySetterTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree property, boolean isStatic, com.google.javascript.jscomp.parsing.parser.trees.FormalParameterListTree parameter, com.google.javascript.jscomp.parsing.parser.trees.BlockTree body)
func NewComputedPropertySetterTree(location SourceRange, property ParseTree, isStatic bool, parameter *FormalParameterListTree, body *BlockTree) *ComputedPropertySetterTree {
	this_ := &ComputedPropertySetterTree{}
	this_.BindSelf_(this_)
	InitComputedPropertySetterTree(this_, location, property, isStatic, parameter, body)
	return this_
}
func allocPooled_ComputedPropertySetterTree(slab_ *j2gort.PooledSlab[ComputedPropertySetterTree], location SourceRange, property ParseTree, isStatic bool, parameter *FormalParameterListTree, body *BlockTree) *ComputedPropertySetterTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitComputedPropertySetterTree(this_, location, property, isStatic, parameter, body)
	return this_
}

func InitComputedPropertySetterTree(this_ *ComputedPropertySetterTree, location SourceRange, property ParseTree, isStatic bool, parameter *FormalParameterListTree, body *BlockTree) {
	_ = location
	_ = property
	_ = isStatic
	_ = parameter
	_ = body
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_COMPUTED_PROPERTY_SETTER(), location)
	this_.initFields_ComputedPropertySetterTree()
	this_.property = property
	this_.isStatic = isStatic
	this_.parameter = parameter
	this_.body = body
}
