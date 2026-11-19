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
// Port of com/google/javascript/jscomp/parsing/parser/trees/ObjectRestTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.ObjectRestTree
type ObjectRestTree struct {
	ParseTreeBase
	assignmentTarget ParseTree `gson:"assignmentTarget"`
}

func (this_ *ObjectRestTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *ObjectRestTree) IsJavaNil_() bool { return this_ == nil }
func (*ObjectRestTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.ObjectRestTree", "ObjectRestTree"
}
func (*ObjectRestTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.ObjectRestTree", "com.google.javascript.jscomp.parsing.parser.trees.ObjectRestTree", func(v_ any) bool {
		return j2gort.InstanceOf[*ObjectRestTree](v_)
	})
}
func (*ObjectRestTree) Implements_parser_ParseTree()               {}
func (this_ *ObjectRestTree) FieldGet_assignmentTarget() ParseTree { return this_.assignmentTarget }

func (this_ *ObjectRestTree) initFields_ObjectRestTree() {
	_ = this_
}

// Java constructor: ObjectRestTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree assignmentTarget)
func NewObjectRestTree(location SourceRange, assignmentTarget ParseTree) *ObjectRestTree {
	this_ := &ObjectRestTree{}
	this_.BindSelf_(this_)
	InitObjectRestTree(this_, location, assignmentTarget)
	return this_
}
func allocPooled_ObjectRestTree(slab_ *j2gort.PooledSlab[ObjectRestTree], location SourceRange, assignmentTarget ParseTree) *ObjectRestTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitObjectRestTree(this_, location, assignmentTarget)
	return this_
}

func InitObjectRestTree(this_ *ObjectRestTree, location SourceRange, assignmentTarget ParseTree) {
	_ = location
	_ = assignmentTarget
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_OBJECT_REST(), location)
	this_.initFields_ObjectRestTree()
	this_.assignmentTarget = assignmentTarget
}
