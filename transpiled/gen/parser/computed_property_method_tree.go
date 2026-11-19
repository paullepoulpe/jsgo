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
// Port of com/google/javascript/jscomp/parsing/parser/trees/ComputedPropertyMethodTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.ComputedPropertyMethodTree
type ComputedPropertyMethodTree struct {
	ParseTreeBase
	property ParseTree `gson:"property"`
	method   ParseTree `gson:"method"`
}

func (this_ *ComputedPropertyMethodTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *ComputedPropertyMethodTree) IsJavaNil_() bool { return this_ == nil }
func (*ComputedPropertyMethodTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.ComputedPropertyMethodTree", "ComputedPropertyMethodTree"
}
func (*ComputedPropertyMethodTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.ComputedPropertyMethodTree", "com.google.javascript.jscomp.parsing.parser.trees.ComputedPropertyMethodTree", func(v_ any) bool {
		return j2gort.InstanceOf[*ComputedPropertyMethodTree](v_)
	})
}
func (*ComputedPropertyMethodTree) Implements_parser_ParseTree()       {}
func (this_ *ComputedPropertyMethodTree) FieldGet_property() ParseTree { return this_.property }
func (this_ *ComputedPropertyMethodTree) FieldGet_method() ParseTree   { return this_.method }

func (this_ *ComputedPropertyMethodTree) initFields_ComputedPropertyMethodTree() {
	_ = this_
}

// Java constructor: ComputedPropertyMethodTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree property, com.google.javascript.jscomp.parsing.parser.trees.ParseTree method)
func NewComputedPropertyMethodTree(location SourceRange, property ParseTree, method ParseTree) *ComputedPropertyMethodTree {
	this_ := &ComputedPropertyMethodTree{}
	this_.BindSelf_(this_)
	InitComputedPropertyMethodTree(this_, location, property, method)
	return this_
}
func allocPooled_ComputedPropertyMethodTree(slab_ *j2gort.PooledSlab[ComputedPropertyMethodTree], location SourceRange, property ParseTree, method ParseTree) *ComputedPropertyMethodTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitComputedPropertyMethodTree(this_, location, property, method)
	return this_
}

func InitComputedPropertyMethodTree(this_ *ComputedPropertyMethodTree, location SourceRange, property ParseTree, method ParseTree) {
	_ = location
	_ = property
	_ = method
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_COMPUTED_PROPERTY_METHOD(), location)
	this_.initFields_ComputedPropertyMethodTree()
	this_.property = property
	this_.method = method
}
