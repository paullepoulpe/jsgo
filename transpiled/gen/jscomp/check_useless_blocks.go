/*
 * Copyright 2016 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/lint/CheckUselessBlocks.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.lint.CheckUselessBlocks
type CheckUselessBlocks struct {
	compiler   AbstractCompiler             `gson:"compiler"`
	loneBlocks j2gort.Deque[pkg_rhino.Node] `gson:"loneBlocks"`
	idHash_    atomic.Int32
}

func (this_ *CheckUselessBlocks) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckUselessBlocks) IsJavaNil_() bool { return this_ == nil }
func (*CheckUselessBlocks) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.lint.CheckUselessBlocks", "CheckUselessBlocks"
}
func (this_ *CheckUselessBlocks) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CheckUselessBlocks) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.lint.CheckUselessBlocks", "com.google.javascript.jscomp.lint.CheckUselessBlocks", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckUselessBlocks](v_)
	})
}
func (*CheckUselessBlocks) Implements_jscomp_NodeTraversal_Callback() {}
func (*CheckUselessBlocks) Implements_jscomp_CompilerPass()           {}

var static_CheckUselessBlocks_USELESS_BLOCK *DiagnosticType

func StaticGet_CheckUselessBlocks_USELESS_BLOCK() *DiagnosticType {
	clinit_CheckUselessBlocks()
	return static_CheckUselessBlocks_USELESS_BLOCK
}
func StaticSet_CheckUselessBlocks_USELESS_BLOCK(v *DiagnosticType) {
	clinit_CheckUselessBlocks()
	static_CheckUselessBlocks_USELESS_BLOCK = v
}

var clinitDone_CheckUselessBlocks uint32
var clinitState_CheckUselessBlocks j2gort.ClinitState

func clinit_CheckUselessBlocks() {
	if atomic.LoadUint32(&clinitDone_CheckUselessBlocks) != 0 {
		return
	}
	clinitSlow_CheckUselessBlocks()
}

func clinitSlow_CheckUselessBlocks() {
	if !j2gort.EnterClinit(&clinitState_CheckUselessBlocks, &clinitDone_CheckUselessBlocks) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckUselessBlocks, &clinitDone_CheckUselessBlocks)
	static_CheckUselessBlocks_USELESS_BLOCK = DiagnosticType_Disabled("JSC_USELESS_BLOCK", "Useless block.")
}

func (this_ *CheckUselessBlocks) initFields_CheckUselessBlocks() {
	_ = this_
}

// Java constructor: CheckUselessBlocks(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewCheckUselessBlocks(compiler AbstractCompiler) *CheckUselessBlocks {
	clinit_CheckUselessBlocks()
	this_ := &CheckUselessBlocks{}
	InitCheckUselessBlocks(this_, compiler)
	return this_
}

func InitCheckUselessBlocks(this_ *CheckUselessBlocks, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_CheckUselessBlocks()
	this_.compiler = compiler
	this_.loneBlocks = j2gort.Deque[pkg_rhino.Node]{}
}

// Java method: CheckUselessBlocks.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *CheckUselessBlocks) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
}

// Java method: CheckUselessBlocks.isLoneBlock(com.google.javascript.rhino.Node n)
func (this_ *CheckUselessBlocks) priv_CheckUselessBlocks_isLoneBlock_6568(n pkg_rhino.Node) bool {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	if (parent != 0) && (parent.IsScript() || ((parent.IsBlock() && !(parent.IsSyntheticBlock())) && !(parent.IsAddedBlock()))) {
		return (!(n.IsSyntheticBlock()) && !(n.IsAddedBlock()))
	}
	return false
}

// Java method: CheckUselessBlocks.allowLoneBlock(com.google.javascript.rhino.Node parent)
func (this_ *CheckUselessBlocks) priv_CheckUselessBlocks_allowLoneBlock_6569(parent pkg_rhino.Node) {
	_ = parent
	if this_.loneBlocks.IsEmpty() {
		return
	}
	if this_.loneBlocks.PeekFirst() == parent {
		this_.loneBlocks.Pop()
	}
}

// Java method: CheckUselessBlocks.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckUselessBlocks) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(90) /* Token.BLOCK */ :
		if this_.priv_CheckUselessBlocks_isLoneBlock_6568(n) {
			this_.loneBlocks.Push(j2gort.CheckNotNullNode(n))
		}
	case int32(111) /* Token.LET */, int32(97) /* Token.CONST */ :
		this_.priv_CheckUselessBlocks_allowLoneBlock_6569(parent)
	case int32(105) /* Token.CLASS */ :
		if NodeUtil_IsClassDeclaration(n) {
			this_.priv_CheckUselessBlocks_allowLoneBlock_6569(parent)
		}
	case int32(73) /* Token.FUNCTION */ :
		if NodeUtil_IsFunctionDeclaration(n) {
			this_.priv_CheckUselessBlocks_allowLoneBlock_6569(parent)
		}
	default:
	}
	return true
}

// Java method: CheckUselessBlocks.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckUselessBlocks) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if (n.IsBlock() && !(this_.loneBlocks.IsEmpty())) && (this_.loneBlocks.PeekFirst() == n) {
		this_.loneBlocks.Pop()
		t.Report_3(n, StaticGet_CheckUselessBlocks_USELESS_BLOCK(), []string{})
	}
}
func (*CheckUselessBlocks) CanReleaseSyntacticScopes_() {}
