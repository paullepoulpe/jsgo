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
// Port of com/google/javascript/jscomp/parsing/parser/trees/TemplateLiteralPortionTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.TemplateLiteralPortionTree
type TemplateLiteralPortionTree struct {
	ParseTreeBase
	value Token `gson:"value"`
}

func (this_ *TemplateLiteralPortionTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *TemplateLiteralPortionTree) IsJavaNil_() bool { return this_ == nil }
func (*TemplateLiteralPortionTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.TemplateLiteralPortionTree", "TemplateLiteralPortionTree"
}
func (*TemplateLiteralPortionTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.TemplateLiteralPortionTree", "com.google.javascript.jscomp.parsing.parser.trees.TemplateLiteralPortionTree", func(v_ any) bool {
		return j2gort.InstanceOf[*TemplateLiteralPortionTree](v_)
	})
}
func (*TemplateLiteralPortionTree) Implements_parser_ParseTree() {}
func (this_ *TemplateLiteralPortionTree) FieldGet_value() Token  { return this_.value }

func (this_ *TemplateLiteralPortionTree) initFields_TemplateLiteralPortionTree() {
	_ = this_
}

// Java constructor: TemplateLiteralPortionTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.Token value)
func NewTemplateLiteralPortionTree(location SourceRange, value Token) *TemplateLiteralPortionTree {
	this_ := &TemplateLiteralPortionTree{}
	this_.BindSelf_(this_)
	InitTemplateLiteralPortionTree(this_, location, value)
	return this_
}
func allocPooled_TemplateLiteralPortionTree(slab_ *j2gort.PooledSlab[TemplateLiteralPortionTree], location SourceRange, value Token) *TemplateLiteralPortionTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitTemplateLiteralPortionTree(this_, location, value)
	return this_
}

func InitTemplateLiteralPortionTree(this_ *TemplateLiteralPortionTree, location SourceRange, value Token) {
	_ = location
	_ = value
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_TEMPLATE_LITERAL_PORTION(), location)
	this_.initFields_TemplateLiteralPortionTree()
	this_.value = value
}
