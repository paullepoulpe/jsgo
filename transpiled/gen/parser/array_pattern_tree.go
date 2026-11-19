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
// Port of com/google/javascript/jscomp/parsing/parser/trees/ArrayPatternTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.ArrayPatternTree
type ArrayPatternTree struct {
	ParseTreeBase
	elements j2gort.Ext_com_google_common_collect_ImmutableList `gson:"elements"`
}

func (this_ *ArrayPatternTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *ArrayPatternTree) IsJavaNil_() bool { return this_ == nil }
func (*ArrayPatternTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.ArrayPatternTree", "ArrayPatternTree"
}
func (*ArrayPatternTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.ArrayPatternTree", "com.google.javascript.jscomp.parsing.parser.trees.ArrayPatternTree", func(v_ any) bool {
		return j2gort.InstanceOf[*ArrayPatternTree](v_)
	})
}
func (*ArrayPatternTree) Implements_parser_ParseTree() {}
func (this_ *ArrayPatternTree) FieldGet_elements() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.elements
}

func (this_ *ArrayPatternTree) initFields_ArrayPatternTree() {
	_ = this_
}

// Java constructor: ArrayPatternTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.parsing.parser.trees.ParseTree> elements)
func NewArrayPatternTree(location SourceRange, elements j2gort.Ext_com_google_common_collect_ImmutableList) *ArrayPatternTree {
	this_ := &ArrayPatternTree{}
	this_.BindSelf_(this_)
	InitArrayPatternTree(this_, location, elements)
	return this_
}
func allocPooled_ArrayPatternTree(slab_ *j2gort.PooledSlab[ArrayPatternTree], location SourceRange, elements j2gort.Ext_com_google_common_collect_ImmutableList) *ArrayPatternTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitArrayPatternTree(this_, location, elements)
	return this_
}

func InitArrayPatternTree(this_ *ArrayPatternTree, location SourceRange, elements j2gort.Ext_com_google_common_collect_ImmutableList) {
	_ = location
	_ = elements
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_ARRAY_PATTERN(), location)
	this_.initFields_ArrayPatternTree()
	this_.elements = elements
}
