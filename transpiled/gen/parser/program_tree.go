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
// Port of com/google/javascript/jscomp/parsing/parser/trees/ProgramTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.ProgramTree
type ProgramTree struct {
	ParseTreeBase
	sourceElements j2gort.Ext_com_google_common_collect_ImmutableList `gson:"sourceElements"`
	sourceComments j2gort.Ext_com_google_common_collect_ImmutableList `gson:"sourceComments"`
}

func (this_ *ProgramTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *ProgramTree) IsJavaNil_() bool { return this_ == nil }
func (*ProgramTree) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.parsing.parser.trees.ProgramTree", "ProgramTree" }
func (*ProgramTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.ProgramTree", "com.google.javascript.jscomp.parsing.parser.trees.ProgramTree", func(v_ any) bool {
		return j2gort.InstanceOf[*ProgramTree](v_)
	})
}
func (*ProgramTree) Implements_parser_ParseTree() {}
func (this_ *ProgramTree) FieldGet_sourceElements() j2gort.Ext_com_google_common_collect_ImmutableList { return this_.sourceElements }
func (this_ *ProgramTree) FieldGet_sourceComments() j2gort.Ext_com_google_common_collect_ImmutableList { return this_.sourceComments }

func (this_ *ProgramTree) initFields_ProgramTree() {
	_ = this_
}

// Java constructor: ProgramTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.parsing.parser.trees.ParseTree> sourceElements, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.parsing.parser.trees.Comment> sourceComments)
func NewProgramTree(location SourceRange, sourceElements j2gort.Ext_com_google_common_collect_ImmutableList, sourceComments j2gort.Ext_com_google_common_collect_ImmutableList) *ProgramTree {
	this_ := &ProgramTree{}
	this_.BindSelf_(this_)
	InitProgramTree(this_, location, sourceElements, sourceComments)
	return this_
}

func InitProgramTree(this_ *ProgramTree, location SourceRange, sourceElements j2gort.Ext_com_google_common_collect_ImmutableList, sourceComments j2gort.Ext_com_google_common_collect_ImmutableList) {
	_ = location
	_ = sourceElements
	_ = sourceComments
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_PROGRAM(), location)
	this_.initFields_ProgramTree()
	this_.sourceElements = sourceElements
	this_.sourceComments = sourceComments
}

