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
// Port of com/google/javascript/jscomp/parsing/parser/trees/ComputedPropertyGetterTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.ComputedPropertyGetterTree
type ComputedPropertyGetterTree struct {
	ParseTreeBase
	property ParseTree  `gson:"property"`
	body     *BlockTree `gson:"body"`
	isStatic bool       `gson:"isStatic"`
}

func (this_ *ComputedPropertyGetterTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *ComputedPropertyGetterTree) IsJavaNil_() bool { return this_ == nil }
func (*ComputedPropertyGetterTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.ComputedPropertyGetterTree", "ComputedPropertyGetterTree"
}
func (*ComputedPropertyGetterTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.ComputedPropertyGetterTree", "com.google.javascript.jscomp.parsing.parser.trees.ComputedPropertyGetterTree", func(v_ any) bool {
		return j2gort.InstanceOf[*ComputedPropertyGetterTree](v_)
	})
}
func (*ComputedPropertyGetterTree) Implements_parser_ParseTree()       {}
func (this_ *ComputedPropertyGetterTree) FieldGet_property() ParseTree { return this_.property }
func (this_ *ComputedPropertyGetterTree) FieldGet_isStatic() bool      { return this_.isStatic }
func (this_ *ComputedPropertyGetterTree) FieldGet_body() *BlockTree    { return this_.body }

func (this_ *ComputedPropertyGetterTree) initFields_ComputedPropertyGetterTree() {
	_ = this_
}

// Java constructor: ComputedPropertyGetterTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree property, boolean isStatic, com.google.javascript.jscomp.parsing.parser.trees.BlockTree body)
func NewComputedPropertyGetterTree(location SourceRange, property ParseTree, isStatic bool, body *BlockTree) *ComputedPropertyGetterTree {
	this_ := &ComputedPropertyGetterTree{}
	this_.BindSelf_(this_)
	InitComputedPropertyGetterTree(this_, location, property, isStatic, body)
	return this_
}
func allocPooled_ComputedPropertyGetterTree(slab_ *j2gort.PooledSlab[ComputedPropertyGetterTree], location SourceRange, property ParseTree, isStatic bool, body *BlockTree) *ComputedPropertyGetterTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitComputedPropertyGetterTree(this_, location, property, isStatic, body)
	return this_
}

func InitComputedPropertyGetterTree(this_ *ComputedPropertyGetterTree, location SourceRange, property ParseTree, isStatic bool, body *BlockTree) {
	_ = location
	_ = property
	_ = isStatic
	_ = body
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_COMPUTED_PROPERTY_GETTER(), location)
	this_.initFields_ComputedPropertyGetterTree()
	this_.property = property
	this_.isStatic = isStatic
	this_.body = body
}
