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
// Port of com/google/javascript/jscomp/parsing/parser/trees/CatchTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.CatchTree
type CatchTree struct {
	ParseTreeBase
	exception ParseTree `gson:"exception"`
	catchBody ParseTree `gson:"catchBody"`
}

func (this_ *CatchTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *CatchTree) IsJavaNil_() bool { return this_ == nil }
func (*CatchTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.CatchTree", "CatchTree"
}
func (*CatchTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.CatchTree", "com.google.javascript.jscomp.parsing.parser.trees.CatchTree", func(v_ any) bool {
		return j2gort.InstanceOf[*CatchTree](v_)
	})
}
func (*CatchTree) Implements_parser_ParseTree()        {}
func (this_ *CatchTree) FieldGet_exception() ParseTree { return this_.exception }
func (this_ *CatchTree) FieldGet_catchBody() ParseTree { return this_.catchBody }

func (this_ *CatchTree) initFields_CatchTree() {
	_ = this_
}

// Java constructor: CatchTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree exception, com.google.javascript.jscomp.parsing.parser.trees.ParseTree catchBody)
func NewCatchTree(location SourceRange, exception ParseTree, catchBody ParseTree) *CatchTree {
	this_ := &CatchTree{}
	this_.BindSelf_(this_)
	InitCatchTree(this_, location, exception, catchBody)
	return this_
}
func allocPooled_CatchTree(slab_ *j2gort.PooledSlab[CatchTree], location SourceRange, exception ParseTree, catchBody ParseTree) *CatchTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitCatchTree(this_, location, exception, catchBody)
	return this_
}

func InitCatchTree(this_ *CatchTree, location SourceRange, exception ParseTree, catchBody ParseTree) {
	_ = location
	_ = exception
	_ = catchBody
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_CATCH(), location)
	this_.initFields_CatchTree()
	this_.exception = exception
	this_.catchBody = catchBody
}
