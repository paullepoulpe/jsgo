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
// Port of com/google/javascript/jscomp/parsing/parser/trees/FinallyTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.FinallyTree
type FinallyTree struct {
	ParseTreeBase
	block ParseTree `gson:"block"`
}

func (this_ *FinallyTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *FinallyTree) IsJavaNil_() bool { return this_ == nil }
func (*FinallyTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.FinallyTree", "FinallyTree"
}
func (*FinallyTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.FinallyTree", "com.google.javascript.jscomp.parsing.parser.trees.FinallyTree", func(v_ any) bool {
		return j2gort.InstanceOf[*FinallyTree](v_)
	})
}
func (*FinallyTree) Implements_parser_ParseTree()    {}
func (this_ *FinallyTree) FieldGet_block() ParseTree { return this_.block }

func (this_ *FinallyTree) initFields_FinallyTree() {
	_ = this_
}

// Java constructor: FinallyTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree block)
func NewFinallyTree(location SourceRange, block ParseTree) *FinallyTree {
	this_ := &FinallyTree{}
	this_.BindSelf_(this_)
	InitFinallyTree(this_, location, block)
	return this_
}
func allocPooled_FinallyTree(slab_ *j2gort.PooledSlab[FinallyTree], location SourceRange, block ParseTree) *FinallyTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitFinallyTree(this_, location, block)
	return this_
}

func InitFinallyTree(this_ *FinallyTree, location SourceRange, block ParseTree) {
	_ = location
	_ = block
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_FINALLY(), location)
	this_.initFields_FinallyTree()
	this_.block = block
}
