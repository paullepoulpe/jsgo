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
// Port of com/google/javascript/jscomp/parsing/parser/trees/IterSpreadTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.IterSpreadTree
type IterSpreadTree struct {
	ParseTreeBase
	expression ParseTree `gson:"expression"`
}

func (this_ *IterSpreadTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *IterSpreadTree) IsJavaNil_() bool { return this_ == nil }
func (*IterSpreadTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.IterSpreadTree", "IterSpreadTree"
}
func (*IterSpreadTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.IterSpreadTree", "com.google.javascript.jscomp.parsing.parser.trees.IterSpreadTree", func(v_ any) bool {
		return j2gort.InstanceOf[*IterSpreadTree](v_)
	})
}
func (*IterSpreadTree) Implements_parser_ParseTree()         {}
func (this_ *IterSpreadTree) FieldGet_expression() ParseTree { return this_.expression }

func (this_ *IterSpreadTree) initFields_IterSpreadTree() {
	_ = this_
}

// Java constructor: IterSpreadTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree expression)
func NewIterSpreadTree(location SourceRange, expression ParseTree) *IterSpreadTree {
	this_ := &IterSpreadTree{}
	this_.BindSelf_(this_)
	InitIterSpreadTree(this_, location, expression)
	return this_
}
func allocPooled_IterSpreadTree(slab_ *j2gort.PooledSlab[IterSpreadTree], location SourceRange, expression ParseTree) *IterSpreadTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitIterSpreadTree(this_, location, expression)
	return this_
}

func InitIterSpreadTree(this_ *IterSpreadTree, location SourceRange, expression ParseTree) {
	_ = location
	_ = expression
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_ITER_SPREAD(), location)
	this_.initFields_IterSpreadTree()
	this_.expression = expression
}
