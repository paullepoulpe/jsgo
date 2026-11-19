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
// Port of com/google/javascript/jscomp/parsing/parser/trees/ComputedPropertyDefinitionTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.ComputedPropertyDefinitionTree
type ComputedPropertyDefinitionTree struct {
	ParseTreeBase
	property ParseTree `gson:"property"`
	value    ParseTree `gson:"value"`
}

func (this_ *ComputedPropertyDefinitionTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *ComputedPropertyDefinitionTree) IsJavaNil_() bool { return this_ == nil }
func (*ComputedPropertyDefinitionTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.ComputedPropertyDefinitionTree", "ComputedPropertyDefinitionTree"
}
func (*ComputedPropertyDefinitionTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.ComputedPropertyDefinitionTree", "com.google.javascript.jscomp.parsing.parser.trees.ComputedPropertyDefinitionTree", func(v_ any) bool {
		return j2gort.InstanceOf[*ComputedPropertyDefinitionTree](v_)
	})
}
func (*ComputedPropertyDefinitionTree) Implements_parser_ParseTree()       {}
func (this_ *ComputedPropertyDefinitionTree) FieldGet_property() ParseTree { return this_.property }
func (this_ *ComputedPropertyDefinitionTree) FieldGet_value() ParseTree    { return this_.value }

func (this_ *ComputedPropertyDefinitionTree) initFields_ComputedPropertyDefinitionTree() {
	_ = this_
}

// Java constructor: ComputedPropertyDefinitionTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree property, com.google.javascript.jscomp.parsing.parser.trees.ParseTree value)
func NewComputedPropertyDefinitionTree(location SourceRange, property ParseTree, value ParseTree) *ComputedPropertyDefinitionTree {
	this_ := &ComputedPropertyDefinitionTree{}
	this_.BindSelf_(this_)
	InitComputedPropertyDefinitionTree(this_, location, property, value)
	return this_
}
func allocPooled_ComputedPropertyDefinitionTree(slab_ *j2gort.PooledSlab[ComputedPropertyDefinitionTree], location SourceRange, property ParseTree, value ParseTree) *ComputedPropertyDefinitionTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitComputedPropertyDefinitionTree(this_, location, property, value)
	return this_
}

func InitComputedPropertyDefinitionTree(this_ *ComputedPropertyDefinitionTree, location SourceRange, property ParseTree, value ParseTree) {
	_ = location
	_ = property
	_ = value
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_COMPUTED_PROPERTY_DEFINITION(), location)
	this_.initFields_ComputedPropertyDefinitionTree()
	this_.property = property
	this_.value = value
}
