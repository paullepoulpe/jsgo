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
// Port of com/google/javascript/jscomp/parsing/parser/trees/ThisExpressionTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.ThisExpressionTree
type ThisExpressionTree struct {
	ParseTreeBase
}

func (this_ *ThisExpressionTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *ThisExpressionTree) IsJavaNil_() bool { return this_ == nil }
func (*ThisExpressionTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.ThisExpressionTree", "ThisExpressionTree"
}
func (*ThisExpressionTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.ThisExpressionTree", "com.google.javascript.jscomp.parsing.parser.trees.ThisExpressionTree", func(v_ any) bool {
		return j2gort.InstanceOf[*ThisExpressionTree](v_)
	})
}
func (*ThisExpressionTree) Implements_parser_ParseTree() {}

func (this_ *ThisExpressionTree) initFields_ThisExpressionTree() {
	_ = this_
}

// Java constructor: ThisExpressionTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location)
func NewThisExpressionTree(location SourceRange) *ThisExpressionTree {
	this_ := &ThisExpressionTree{}
	this_.BindSelf_(this_)
	InitThisExpressionTree(this_, location)
	return this_
}
func allocPooled_ThisExpressionTree(slab_ *j2gort.PooledSlab[ThisExpressionTree], location SourceRange) *ThisExpressionTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitThisExpressionTree(this_, location)
	return this_
}

func InitThisExpressionTree(this_ *ThisExpressionTree, location SourceRange) {
	_ = location
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_THIS_EXPRESSION(), location)
	this_.initFields_ThisExpressionTree()
}
