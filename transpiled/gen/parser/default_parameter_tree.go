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
// Port of com/google/javascript/jscomp/parsing/parser/trees/DefaultParameterTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.DefaultParameterTree
type DefaultParameterTree struct {
	ParseTreeBase
	lhs          ParseTree `gson:"lhs"`
	defaultValue ParseTree `gson:"defaultValue"`
}

func (this_ *DefaultParameterTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *DefaultParameterTree) IsJavaNil_() bool { return this_ == nil }
func (*DefaultParameterTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.DefaultParameterTree", "DefaultParameterTree"
}
func (*DefaultParameterTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.DefaultParameterTree", "com.google.javascript.jscomp.parsing.parser.trees.DefaultParameterTree", func(v_ any) bool {
		return j2gort.InstanceOf[*DefaultParameterTree](v_)
	})
}
func (*DefaultParameterTree) Implements_parser_ParseTree()           {}
func (this_ *DefaultParameterTree) FieldGet_lhs() ParseTree          { return this_.lhs }
func (this_ *DefaultParameterTree) FieldGet_defaultValue() ParseTree { return this_.defaultValue }

func (this_ *DefaultParameterTree) initFields_DefaultParameterTree() {
	_ = this_
}

// Java constructor: DefaultParameterTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree lhs, com.google.javascript.jscomp.parsing.parser.trees.ParseTree defaultValue)
func NewDefaultParameterTree(location SourceRange, lhs ParseTree, defaultValue ParseTree) *DefaultParameterTree {
	this_ := &DefaultParameterTree{}
	this_.BindSelf_(this_)
	InitDefaultParameterTree(this_, location, lhs, defaultValue)
	return this_
}
func allocPooled_DefaultParameterTree(slab_ *j2gort.PooledSlab[DefaultParameterTree], location SourceRange, lhs ParseTree, defaultValue ParseTree) *DefaultParameterTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitDefaultParameterTree(this_, location, lhs, defaultValue)
	return this_
}

func InitDefaultParameterTree(this_ *DefaultParameterTree, location SourceRange, lhs ParseTree, defaultValue ParseTree) {
	_ = location
	_ = lhs
	_ = defaultValue
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_DEFAULT_PARAMETER(), location)
	this_.initFields_DefaultParameterTree()
	this_.lhs = lhs
	this_.defaultValue = defaultValue
}
