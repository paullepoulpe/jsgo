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
// Port of com/google/javascript/jscomp/parsing/parser/trees/ThrowStatementTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.ThrowStatementTree
type ThrowStatementTree struct {
	ParseTreeBase
	value ParseTree `gson:"value"`
}

func (this_ *ThrowStatementTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *ThrowStatementTree) IsJavaNil_() bool { return this_ == nil }
func (*ThrowStatementTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.ThrowStatementTree", "ThrowStatementTree"
}
func (*ThrowStatementTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.ThrowStatementTree", "com.google.javascript.jscomp.parsing.parser.trees.ThrowStatementTree", func(v_ any) bool {
		return j2gort.InstanceOf[*ThrowStatementTree](v_)
	})
}
func (*ThrowStatementTree) Implements_parser_ParseTree()    {}
func (this_ *ThrowStatementTree) FieldGet_value() ParseTree { return this_.value }

func (this_ *ThrowStatementTree) initFields_ThrowStatementTree() {
	_ = this_
}

// Java constructor: ThrowStatementTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree value)
func NewThrowStatementTree(location SourceRange, value ParseTree) *ThrowStatementTree {
	this_ := &ThrowStatementTree{}
	this_.BindSelf_(this_)
	InitThrowStatementTree(this_, location, value)
	return this_
}
func allocPooled_ThrowStatementTree(slab_ *j2gort.PooledSlab[ThrowStatementTree], location SourceRange, value ParseTree) *ThrowStatementTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitThrowStatementTree(this_, location, value)
	return this_
}

func InitThrowStatementTree(this_ *ThrowStatementTree, location SourceRange, value ParseTree) {
	_ = location
	_ = value
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_THROW_STATEMENT(), location)
	this_.initFields_ThrowStatementTree()
	this_.value = value
}
