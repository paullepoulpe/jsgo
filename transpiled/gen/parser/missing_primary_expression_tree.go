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
// Port of com/google/javascript/jscomp/parsing/parser/trees/MissingPrimaryExpressionTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.MissingPrimaryExpressionTree
type MissingPrimaryExpressionTree struct {
	ParseTreeBase
}

func (this_ *MissingPrimaryExpressionTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *MissingPrimaryExpressionTree) IsJavaNil_() bool { return this_ == nil }
func (*MissingPrimaryExpressionTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.MissingPrimaryExpressionTree", "MissingPrimaryExpressionTree"
}
func (*MissingPrimaryExpressionTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.MissingPrimaryExpressionTree", "com.google.javascript.jscomp.parsing.parser.trees.MissingPrimaryExpressionTree", func(v_ any) bool {
		return j2gort.InstanceOf[*MissingPrimaryExpressionTree](v_)
	})
}
func (*MissingPrimaryExpressionTree) Implements_parser_ParseTree() {}

func (this_ *MissingPrimaryExpressionTree) initFields_MissingPrimaryExpressionTree() {
	_ = this_
}

// Java constructor: MissingPrimaryExpressionTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location)
func NewMissingPrimaryExpressionTree(location SourceRange) *MissingPrimaryExpressionTree {
	this_ := &MissingPrimaryExpressionTree{}
	this_.BindSelf_(this_)
	InitMissingPrimaryExpressionTree(this_, location)
	return this_
}
func allocPooled_MissingPrimaryExpressionTree(slab_ *j2gort.PooledSlab[MissingPrimaryExpressionTree], location SourceRange) *MissingPrimaryExpressionTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitMissingPrimaryExpressionTree(this_, location)
	return this_
}

func InitMissingPrimaryExpressionTree(this_ *MissingPrimaryExpressionTree, location SourceRange) {
	_ = location
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_MISSING_PRIMARY_EXPRESSION(), location)
	this_.initFields_MissingPrimaryExpressionTree()
}
