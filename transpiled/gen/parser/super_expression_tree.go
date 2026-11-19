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
// Port of com/google/javascript/jscomp/parsing/parser/trees/SuperExpressionTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.SuperExpressionTree
type SuperExpressionTree struct {
	ParseTreeBase
}

func (this_ *SuperExpressionTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *SuperExpressionTree) IsJavaNil_() bool { return this_ == nil }
func (*SuperExpressionTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.SuperExpressionTree", "SuperExpressionTree"
}
func (*SuperExpressionTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.SuperExpressionTree", "com.google.javascript.jscomp.parsing.parser.trees.SuperExpressionTree", func(v_ any) bool {
		return j2gort.InstanceOf[*SuperExpressionTree](v_)
	})
}
func (*SuperExpressionTree) Implements_parser_ParseTree() {}

func (this_ *SuperExpressionTree) initFields_SuperExpressionTree() {
	_ = this_
}

// Java constructor: SuperExpressionTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location)
func NewSuperExpressionTree(location SourceRange) *SuperExpressionTree {
	this_ := &SuperExpressionTree{}
	this_.BindSelf_(this_)
	InitSuperExpressionTree(this_, location)
	return this_
}
func allocPooled_SuperExpressionTree(slab_ *j2gort.PooledSlab[SuperExpressionTree], location SourceRange) *SuperExpressionTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitSuperExpressionTree(this_, location)
	return this_
}

func InitSuperExpressionTree(this_ *SuperExpressionTree, location SourceRange) {
	_ = location
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_SUPER_EXPRESSION(), location)
	this_.initFields_SuperExpressionTree()
}
