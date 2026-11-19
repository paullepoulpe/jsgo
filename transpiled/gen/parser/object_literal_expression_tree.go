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
// Port of com/google/javascript/jscomp/parsing/parser/trees/ObjectLiteralExpressionTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.ObjectLiteralExpressionTree
type ObjectLiteralExpressionTree struct {
	ParseTreeBase
	propertyNameAndValues j2gort.Ext_com_google_common_collect_ImmutableList `gson:"propertyNameAndValues"`
	hasTrailingComma      bool                                               `gson:"hasTrailingComma"`
}

func (this_ *ObjectLiteralExpressionTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *ObjectLiteralExpressionTree) IsJavaNil_() bool { return this_ == nil }
func (*ObjectLiteralExpressionTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.ObjectLiteralExpressionTree", "ObjectLiteralExpressionTree"
}
func (*ObjectLiteralExpressionTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.ObjectLiteralExpressionTree", "com.google.javascript.jscomp.parsing.parser.trees.ObjectLiteralExpressionTree", func(v_ any) bool {
		return j2gort.InstanceOf[*ObjectLiteralExpressionTree](v_)
	})
}
func (*ObjectLiteralExpressionTree) Implements_parser_ParseTree() {}
func (this_ *ObjectLiteralExpressionTree) FieldGet_propertyNameAndValues() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.propertyNameAndValues
}
func (this_ *ObjectLiteralExpressionTree) FieldGet_hasTrailingComma() bool {
	return this_.hasTrailingComma
}

func (this_ *ObjectLiteralExpressionTree) initFields_ObjectLiteralExpressionTree() {
	_ = this_
}

// Java constructor: ObjectLiteralExpressionTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.parsing.parser.trees.ParseTree> propertyNameAndValues, boolean hasTrailingComma)
func NewObjectLiteralExpressionTree(location SourceRange, propertyNameAndValues j2gort.Ext_com_google_common_collect_ImmutableList, hasTrailingComma bool) *ObjectLiteralExpressionTree {
	this_ := &ObjectLiteralExpressionTree{}
	this_.BindSelf_(this_)
	InitObjectLiteralExpressionTree(this_, location, propertyNameAndValues, hasTrailingComma)
	return this_
}
func allocPooled_ObjectLiteralExpressionTree(slab_ *j2gort.PooledSlab[ObjectLiteralExpressionTree], location SourceRange, propertyNameAndValues j2gort.Ext_com_google_common_collect_ImmutableList, hasTrailingComma bool) *ObjectLiteralExpressionTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitObjectLiteralExpressionTree(this_, location, propertyNameAndValues, hasTrailingComma)
	return this_
}

func InitObjectLiteralExpressionTree(this_ *ObjectLiteralExpressionTree, location SourceRange, propertyNameAndValues j2gort.Ext_com_google_common_collect_ImmutableList, hasTrailingComma bool) {
	_ = location
	_ = propertyNameAndValues
	_ = hasTrailingComma
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_OBJECT_LITERAL_EXPRESSION(), location)
	this_.initFields_ObjectLiteralExpressionTree()
	this_.propertyNameAndValues = propertyNameAndValues
	this_.hasTrailingComma = hasTrailingComma
}
