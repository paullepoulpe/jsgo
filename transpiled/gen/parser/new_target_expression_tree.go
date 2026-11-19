/*
 * Copyright 2016 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/parsing/parser/trees/NewTargetExpressionTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.NewTargetExpressionTree
type NewTargetExpressionTree struct {
	ParseTreeBase
}

func (this_ *NewTargetExpressionTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *NewTargetExpressionTree) IsJavaNil_() bool { return this_ == nil }
func (*NewTargetExpressionTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.NewTargetExpressionTree", "NewTargetExpressionTree"
}
func (*NewTargetExpressionTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.NewTargetExpressionTree", "com.google.javascript.jscomp.parsing.parser.trees.NewTargetExpressionTree", func(v_ any) bool {
		return j2gort.InstanceOf[*NewTargetExpressionTree](v_)
	})
}
func (*NewTargetExpressionTree) Implements_parser_ParseTree() {}

func (this_ *NewTargetExpressionTree) initFields_NewTargetExpressionTree() {
	_ = this_
}

// Java constructor: NewTargetExpressionTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location)
func NewNewTargetExpressionTree(location SourceRange) *NewTargetExpressionTree {
	this_ := &NewTargetExpressionTree{}
	this_.BindSelf_(this_)
	InitNewTargetExpressionTree(this_, location)
	return this_
}
func allocPooled_NewTargetExpressionTree(slab_ *j2gort.PooledSlab[NewTargetExpressionTree], location SourceRange) *NewTargetExpressionTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitNewTargetExpressionTree(this_, location)
	return this_
}

func InitNewTargetExpressionTree(this_ *NewTargetExpressionTree, location SourceRange) {
	_ = location
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_NEW_TARGET_EXPRESSION(), location)
	this_.initFields_NewTargetExpressionTree()
}
