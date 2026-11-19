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
// Port of com/google/javascript/jscomp/parsing/parser/trees/ObjectPatternTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.ObjectPatternTree
type ObjectPatternTree struct {
	ParseTreeBase
	fields j2gort.Ext_com_google_common_collect_ImmutableList `gson:"fields"`
}

func (this_ *ObjectPatternTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *ObjectPatternTree) IsJavaNil_() bool { return this_ == nil }
func (*ObjectPatternTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.ObjectPatternTree", "ObjectPatternTree"
}
func (*ObjectPatternTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.ObjectPatternTree", "com.google.javascript.jscomp.parsing.parser.trees.ObjectPatternTree", func(v_ any) bool {
		return j2gort.InstanceOf[*ObjectPatternTree](v_)
	})
}
func (*ObjectPatternTree) Implements_parser_ParseTree() {}
func (this_ *ObjectPatternTree) FieldGet_fields() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.fields
}

func (this_ *ObjectPatternTree) initFields_ObjectPatternTree() {
	_ = this_
}

// Java constructor: ObjectPatternTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.parsing.parser.trees.ParseTree> fields)
func NewObjectPatternTree(location SourceRange, fields j2gort.Ext_com_google_common_collect_ImmutableList) *ObjectPatternTree {
	this_ := &ObjectPatternTree{}
	this_.BindSelf_(this_)
	InitObjectPatternTree(this_, location, fields)
	return this_
}
func allocPooled_ObjectPatternTree(slab_ *j2gort.PooledSlab[ObjectPatternTree], location SourceRange, fields j2gort.Ext_com_google_common_collect_ImmutableList) *ObjectPatternTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitObjectPatternTree(this_, location, fields)
	return this_
}

func InitObjectPatternTree(this_ *ObjectPatternTree, location SourceRange, fields j2gort.Ext_com_google_common_collect_ImmutableList) {
	_ = location
	_ = fields
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_OBJECT_PATTERN(), location)
	this_.initFields_ObjectPatternTree()
	this_.fields = fields
}
