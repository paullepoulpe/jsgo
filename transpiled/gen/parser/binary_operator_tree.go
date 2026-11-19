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
// Port of com/google/javascript/jscomp/parsing/parser/trees/BinaryOperatorTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.BinaryOperatorTree
type BinaryOperatorTree struct {
	ParseTreeBase
	left     ParseTree `gson:"left"`
	operator Token     `gson:"operator"`
	right    ParseTree `gson:"right"`
}

func (this_ *BinaryOperatorTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *BinaryOperatorTree) IsJavaNil_() bool { return this_ == nil }
func (*BinaryOperatorTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.BinaryOperatorTree", "BinaryOperatorTree"
}
func (*BinaryOperatorTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.BinaryOperatorTree", "com.google.javascript.jscomp.parsing.parser.trees.BinaryOperatorTree", func(v_ any) bool {
		return j2gort.InstanceOf[*BinaryOperatorTree](v_)
	})
}
func (*BinaryOperatorTree) Implements_parser_ParseTree()    {}
func (this_ *BinaryOperatorTree) FieldGet_left() ParseTree  { return this_.left }
func (this_ *BinaryOperatorTree) FieldGet_operator() Token  { return this_.operator }
func (this_ *BinaryOperatorTree) FieldGet_right() ParseTree { return this_.right }

func (this_ *BinaryOperatorTree) initFields_BinaryOperatorTree() {
	_ = this_
}

// Java constructor: BinaryOperatorTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree left, com.google.javascript.jscomp.parsing.parser.Token operator, com.google.javascript.jscomp.parsing.parser.trees.ParseTree right)
func NewBinaryOperatorTree(location SourceRange, left ParseTree, operator Token, right ParseTree) *BinaryOperatorTree {
	this_ := &BinaryOperatorTree{}
	this_.BindSelf_(this_)
	InitBinaryOperatorTree(this_, location, left, operator, right)
	return this_
}
func allocPooled_BinaryOperatorTree(slab_ *j2gort.PooledSlab[BinaryOperatorTree], location SourceRange, left ParseTree, operator Token, right ParseTree) *BinaryOperatorTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitBinaryOperatorTree(this_, location, left, operator, right)
	return this_
}

func InitBinaryOperatorTree(this_ *BinaryOperatorTree, location SourceRange, left ParseTree, operator Token, right ParseTree) {
	_ = location
	_ = left
	_ = operator
	_ = right
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_BINARY_OPERATOR(), location)
	this_.initFields_BinaryOperatorTree()
	this_.left = left
	this_.operator = operator
	this_.right = right
}
