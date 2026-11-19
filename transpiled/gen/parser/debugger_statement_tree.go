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
// Port of com/google/javascript/jscomp/parsing/parser/trees/DebuggerStatementTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.DebuggerStatementTree
type DebuggerStatementTree struct {
	ParseTreeBase
}

func (this_ *DebuggerStatementTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *DebuggerStatementTree) IsJavaNil_() bool { return this_ == nil }
func (*DebuggerStatementTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.DebuggerStatementTree", "DebuggerStatementTree"
}
func (*DebuggerStatementTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.DebuggerStatementTree", "com.google.javascript.jscomp.parsing.parser.trees.DebuggerStatementTree", func(v_ any) bool {
		return j2gort.InstanceOf[*DebuggerStatementTree](v_)
	})
}
func (*DebuggerStatementTree) Implements_parser_ParseTree() {}

func (this_ *DebuggerStatementTree) initFields_DebuggerStatementTree() {
	_ = this_
}

// Java constructor: DebuggerStatementTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location)
func NewDebuggerStatementTree(location SourceRange) *DebuggerStatementTree {
	this_ := &DebuggerStatementTree{}
	this_.BindSelf_(this_)
	InitDebuggerStatementTree(this_, location)
	return this_
}
func allocPooled_DebuggerStatementTree(slab_ *j2gort.PooledSlab[DebuggerStatementTree], location SourceRange) *DebuggerStatementTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitDebuggerStatementTree(this_, location)
	return this_
}

func InitDebuggerStatementTree(this_ *DebuggerStatementTree, location SourceRange) {
	_ = location
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_DEBUGGER_STATEMENT(), location)
	this_.initFields_DebuggerStatementTree()
}
