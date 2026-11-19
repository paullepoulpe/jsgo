/*
 * Copyright 2014 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/parsing/parser/trees/TemplateSubstitutionTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.TemplateSubstitutionTree
type TemplateSubstitutionTree struct {
	ParseTreeBase
	expression ParseTree `gson:"expression"`
}

func (this_ *TemplateSubstitutionTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *TemplateSubstitutionTree) IsJavaNil_() bool { return this_ == nil }
func (*TemplateSubstitutionTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.TemplateSubstitutionTree", "TemplateSubstitutionTree"
}
func (*TemplateSubstitutionTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.TemplateSubstitutionTree", "com.google.javascript.jscomp.parsing.parser.trees.TemplateSubstitutionTree", func(v_ any) bool {
		return j2gort.InstanceOf[*TemplateSubstitutionTree](v_)
	})
}
func (*TemplateSubstitutionTree) Implements_parser_ParseTree()         {}
func (this_ *TemplateSubstitutionTree) FieldGet_expression() ParseTree { return this_.expression }

func (this_ *TemplateSubstitutionTree) initFields_TemplateSubstitutionTree() {
	_ = this_
}

// Java constructor: TemplateSubstitutionTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree expression)
func NewTemplateSubstitutionTree(location SourceRange, expression ParseTree) *TemplateSubstitutionTree {
	this_ := &TemplateSubstitutionTree{}
	this_.BindSelf_(this_)
	InitTemplateSubstitutionTree(this_, location, expression)
	return this_
}
func allocPooled_TemplateSubstitutionTree(slab_ *j2gort.PooledSlab[TemplateSubstitutionTree], location SourceRange, expression ParseTree) *TemplateSubstitutionTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitTemplateSubstitutionTree(this_, location, expression)
	return this_
}

func InitTemplateSubstitutionTree(this_ *TemplateSubstitutionTree, location SourceRange, expression ParseTree) {
	_ = location
	_ = expression
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_TEMPLATE_SUBSTITUTION(), location)
	this_.initFields_TemplateSubstitutionTree()
	this_.expression = expression
}
