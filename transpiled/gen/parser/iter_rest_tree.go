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
// Port of com/google/javascript/jscomp/parsing/parser/trees/IterRestTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.IterRestTree
type IterRestTree struct {
	ParseTreeBase
	assignmentTarget ParseTree `gson:"assignmentTarget"`
}

func (this_ *IterRestTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *IterRestTree) IsJavaNil_() bool { return this_ == nil }
func (*IterRestTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.IterRestTree", "IterRestTree"
}
func (*IterRestTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.IterRestTree", "com.google.javascript.jscomp.parsing.parser.trees.IterRestTree", func(v_ any) bool {
		return j2gort.InstanceOf[*IterRestTree](v_)
	})
}
func (*IterRestTree) Implements_parser_ParseTree()               {}
func (this_ *IterRestTree) FieldGet_assignmentTarget() ParseTree { return this_.assignmentTarget }

func (this_ *IterRestTree) initFields_IterRestTree() {
	_ = this_
}

// Java constructor: IterRestTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree assignmentTarget)
func NewIterRestTree(location SourceRange, assignmentTarget ParseTree) *IterRestTree {
	this_ := &IterRestTree{}
	this_.BindSelf_(this_)
	InitIterRestTree(this_, location, assignmentTarget)
	return this_
}
func allocPooled_IterRestTree(slab_ *j2gort.PooledSlab[IterRestTree], location SourceRange, assignmentTarget ParseTree) *IterRestTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitIterRestTree(this_, location, assignmentTarget)
	return this_
}

func InitIterRestTree(this_ *IterRestTree, location SourceRange, assignmentTarget ParseTree) {
	_ = location
	_ = assignmentTarget
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_ITER_REST(), location)
	this_.initFields_IterRestTree()
	this_.assignmentTarget = assignmentTarget
}
