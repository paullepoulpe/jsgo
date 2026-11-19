/*
 * Copyright 2015 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/WhitespaceWrapGoogModules.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.WhitespaceWrapGoogModules
type WhitespaceWrapGoogModules struct {
	compiler AbstractCompiler `gson:"compiler"`
	idHash_ atomic.Int32
}

func (this_ *WhitespaceWrapGoogModules) BindSelf_(s any) {
	_ = s
}
func (this_ *WhitespaceWrapGoogModules) IsJavaNil_() bool { return this_ == nil }
func (*WhitespaceWrapGoogModules) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.WhitespaceWrapGoogModules", "WhitespaceWrapGoogModules" }
func (this_ *WhitespaceWrapGoogModules) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*WhitespaceWrapGoogModules) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.WhitespaceWrapGoogModules", "com.google.javascript.jscomp.WhitespaceWrapGoogModules", func(v_ any) bool {
		return j2gort.InstanceOf[*WhitespaceWrapGoogModules](v_)
	})
}
func (*WhitespaceWrapGoogModules) Implements_jscomp_CompilerPass() {}

func (this_ *WhitespaceWrapGoogModules) initFields_WhitespaceWrapGoogModules() {
	_ = this_
}

// Java constructor: WhitespaceWrapGoogModules(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewWhitespaceWrapGoogModules(compiler AbstractCompiler) *WhitespaceWrapGoogModules {
	this_ := &WhitespaceWrapGoogModules{}
	InitWhitespaceWrapGoogModules(this_, compiler)
	return this_
}

func InitWhitespaceWrapGoogModules(this_ *WhitespaceWrapGoogModules, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_WhitespaceWrapGoogModules()
	this_.compiler = compiler
}

// Java method: WhitespaceWrapGoogModules.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *WhitespaceWrapGoogModules) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	{
		var c pkg_rhino.Node = root.GetFirstChild()
		_ = c
		for ; (c != 0); c = c.GetNext() {
			if !(c.IsScript()) { j2gort.CheckState(false) }
			this_.priv_WhitespaceWrapGoogModules_visitScript_5881(c)
		}
	}
}

// Java method: WhitespaceWrapGoogModules.visitScript(com.google.javascript.rhino.Node scriptRoot)
func (this_ *WhitespaceWrapGoogModules) priv_WhitespaceWrapGoogModules_visitScript_5881(scriptRoot pkg_rhino.Node) {
	_ = scriptRoot
	if (!(NodeUtil_IsGoogModuleFile(scriptRoot))) {
		return
	}
	var moduleBody pkg_rhino.Node = scriptRoot.GetFirstChild()
	_ = moduleBody
	moduleBody.SetToken(pkg_rhino.StaticGet_Token_BLOCK())
	NodeUtil_TryMergeBlock(moduleBody, true)
	this_.compiler.ReportChangeToEnclosingScope(scriptRoot)
	var block pkg_rhino.Node = pkg_rhino.IR_Block_0()
	_ = block
	block.AddChildToBack(pkg_rhino.IR_ExprResult(pkg_rhino.IR_String("use strict")))
	var function pkg_rhino.Node = pkg_rhino.IR_Function(pkg_rhino.IR_Name(""), pkg_rhino.IR_ParamList([]pkg_rhino.Node{pkg_rhino.IR_Name("exports")}), block)
	_ = function
	this_.compiler.ReportChangeToChangeScope(function)
	var loadMod pkg_rhino.Node = pkg_rhino.IR_ExprResult(pkg_rhino.IR_Call(pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name("goog"), "loadModule"), []pkg_rhino.Node{function})).SrcrefTree(scriptRoot)
	_ = loadMod
	if (scriptRoot.HasChildren()) {
		block.AddChildrenToBack(scriptRoot.RemoveChildren())
	}
	block.AddChildToBack(pkg_rhino.IR_ReturnNode_1(pkg_rhino.IR_Name("exports")).SrcrefTree(scriptRoot))
	scriptRoot.AddChildToBack(loadMod)
	this_.compiler.ReportChangeToChangeScope(scriptRoot)
}

