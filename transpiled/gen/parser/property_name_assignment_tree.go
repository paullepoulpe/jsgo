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
// Port of com/google/javascript/jscomp/parsing/parser/trees/PropertyNameAssignmentTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.PropertyNameAssignmentTree
type PropertyNameAssignmentTree struct {
	ParseTreeBase
	name  Token     `gson:"name"`
	value ParseTree `gson:"value"`
}

func (this_ *PropertyNameAssignmentTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *PropertyNameAssignmentTree) IsJavaNil_() bool { return this_ == nil }
func (*PropertyNameAssignmentTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.PropertyNameAssignmentTree", "PropertyNameAssignmentTree"
}
func (*PropertyNameAssignmentTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.PropertyNameAssignmentTree", "com.google.javascript.jscomp.parsing.parser.trees.PropertyNameAssignmentTree", func(v_ any) bool {
		return j2gort.InstanceOf[*PropertyNameAssignmentTree](v_)
	})
}
func (*PropertyNameAssignmentTree) Implements_parser_ParseTree()    {}
func (this_ *PropertyNameAssignmentTree) FieldGet_name() Token      { return this_.name }
func (this_ *PropertyNameAssignmentTree) FieldGet_value() ParseTree { return this_.value }

func (this_ *PropertyNameAssignmentTree) initFields_PropertyNameAssignmentTree() {
	_ = this_
}

// Java constructor: PropertyNameAssignmentTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.Token name, com.google.javascript.jscomp.parsing.parser.trees.ParseTree value)
func NewPropertyNameAssignmentTree(location SourceRange, name Token, value ParseTree) *PropertyNameAssignmentTree {
	this_ := &PropertyNameAssignmentTree{}
	this_.BindSelf_(this_)
	InitPropertyNameAssignmentTree(this_, location, name, value)
	return this_
}
func allocPooled_PropertyNameAssignmentTree(slab_ *j2gort.PooledSlab[PropertyNameAssignmentTree], location SourceRange, name Token, value ParseTree) *PropertyNameAssignmentTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitPropertyNameAssignmentTree(this_, location, name, value)
	return this_
}

func InitPropertyNameAssignmentTree(this_ *PropertyNameAssignmentTree, location SourceRange, name Token, value ParseTree) {
	_ = location
	_ = name
	_ = value
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_PROPERTY_NAME_ASSIGNMENT(), location)
	this_.initFields_PropertyNameAssignmentTree()
	this_.name = name
	this_.value = value
}
