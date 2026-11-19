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
// Port of com/google/javascript/jscomp/parsing/parser/trees/ObjectSpreadTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.ObjectSpreadTree
type ObjectSpreadTree struct {
	ParseTreeBase
	expression ParseTree `gson:"expression"`
}

func (this_ *ObjectSpreadTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *ObjectSpreadTree) IsJavaNil_() bool { return this_ == nil }
func (*ObjectSpreadTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.ObjectSpreadTree", "ObjectSpreadTree"
}
func (*ObjectSpreadTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.ObjectSpreadTree", "com.google.javascript.jscomp.parsing.parser.trees.ObjectSpreadTree", func(v_ any) bool {
		return j2gort.InstanceOf[*ObjectSpreadTree](v_)
	})
}
func (*ObjectSpreadTree) Implements_parser_ParseTree()         {}
func (this_ *ObjectSpreadTree) FieldGet_expression() ParseTree { return this_.expression }

func (this_ *ObjectSpreadTree) initFields_ObjectSpreadTree() {
	_ = this_
}

// Java constructor: ObjectSpreadTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree expression)
func NewObjectSpreadTree(location SourceRange, expression ParseTree) *ObjectSpreadTree {
	this_ := &ObjectSpreadTree{}
	this_.BindSelf_(this_)
	InitObjectSpreadTree(this_, location, expression)
	return this_
}
func allocPooled_ObjectSpreadTree(slab_ *j2gort.PooledSlab[ObjectSpreadTree], location SourceRange, expression ParseTree) *ObjectSpreadTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitObjectSpreadTree(this_, location, expression)
	return this_
}

func InitObjectSpreadTree(this_ *ObjectSpreadTree, location SourceRange, expression ParseTree) {
	_ = location
	_ = expression
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_OBJECT_SPREAD(), location)
	this_.initFields_ObjectSpreadTree()
	this_.expression = expression
}
