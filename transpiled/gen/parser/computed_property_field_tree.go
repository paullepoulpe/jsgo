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
// Port of com/google/javascript/jscomp/parsing/parser/trees/ComputedPropertyFieldTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.ComputedPropertyFieldTree
type ComputedPropertyFieldTree struct {
	ParseTreeBase
	property    ParseTree `gson:"property"`
	initializer ParseTree `gson:"initializer"`
	isStatic    bool      `gson:"isStatic"`
}

func (this_ *ComputedPropertyFieldTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *ComputedPropertyFieldTree) IsJavaNil_() bool { return this_ == nil }
func (*ComputedPropertyFieldTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.ComputedPropertyFieldTree", "ComputedPropertyFieldTree"
}
func (*ComputedPropertyFieldTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.ComputedPropertyFieldTree", "com.google.javascript.jscomp.parsing.parser.trees.ComputedPropertyFieldTree", func(v_ any) bool {
		return j2gort.InstanceOf[*ComputedPropertyFieldTree](v_)
	})
}
func (*ComputedPropertyFieldTree) Implements_parser_ParseTree()          {}
func (this_ *ComputedPropertyFieldTree) FieldGet_property() ParseTree    { return this_.property }
func (this_ *ComputedPropertyFieldTree) FieldGet_isStatic() bool         { return this_.isStatic }
func (this_ *ComputedPropertyFieldTree) FieldGet_initializer() ParseTree { return this_.initializer }

func (this_ *ComputedPropertyFieldTree) initFields_ComputedPropertyFieldTree() {
	_ = this_
}

// Java constructor: ComputedPropertyFieldTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree property, boolean isStatic, com.google.javascript.jscomp.parsing.parser.trees.ParseTree initializer)
func NewComputedPropertyFieldTree(location SourceRange, property ParseTree, isStatic bool, initializer ParseTree) *ComputedPropertyFieldTree {
	this_ := &ComputedPropertyFieldTree{}
	this_.BindSelf_(this_)
	InitComputedPropertyFieldTree(this_, location, property, isStatic, initializer)
	return this_
}
func allocPooled_ComputedPropertyFieldTree(slab_ *j2gort.PooledSlab[ComputedPropertyFieldTree], location SourceRange, property ParseTree, isStatic bool, initializer ParseTree) *ComputedPropertyFieldTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitComputedPropertyFieldTree(this_, location, property, isStatic, initializer)
	return this_
}

func InitComputedPropertyFieldTree(this_ *ComputedPropertyFieldTree, location SourceRange, property ParseTree, isStatic bool, initializer ParseTree) {
	_ = location
	_ = property
	_ = isStatic
	_ = initializer
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_COMPUTED_PROPERTY_FIELD(), location)
	this_.initFields_ComputedPropertyFieldTree()
	this_.property = property
	this_.isStatic = isStatic
	this_.initializer = initializer
}
