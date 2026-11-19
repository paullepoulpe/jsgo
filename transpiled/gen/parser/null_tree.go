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
// Port of com/google/javascript/jscomp/parsing/parser/trees/NullTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.NullTree
type NullTree struct {
	ParseTreeBase
}

func (this_ *NullTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *NullTree) IsJavaNil_() bool { return this_ == nil }
func (*NullTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.NullTree", "NullTree"
}
func (*NullTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.NullTree", "com.google.javascript.jscomp.parsing.parser.trees.NullTree", func(v_ any) bool {
		return j2gort.InstanceOf[*NullTree](v_)
	})
}
func (*NullTree) Implements_parser_ParseTree() {}

func (this_ *NullTree) initFields_NullTree() {
	_ = this_
}

// Java constructor: NullTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location)
func NewNullTree(location SourceRange) *NullTree {
	this_ := &NullTree{}
	this_.BindSelf_(this_)
	InitNullTree(this_, location)
	return this_
}
func allocPooled_NullTree(slab_ *j2gort.PooledSlab[NullTree], location SourceRange) *NullTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitNullTree(this_, location)
	return this_
}

func InitNullTree(this_ *NullTree, location SourceRange) {
	_ = location
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_NULL(), location)
	this_.initFields_NullTree()
}
