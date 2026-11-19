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
// Port of com/google/javascript/jscomp/CheckDebuggerStatement.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.CheckDebuggerStatement
type CheckDebuggerStatement struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler `gson:"compiler"`
}

func (this_ *CheckDebuggerStatement) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckDebuggerStatement) IsJavaNil_() bool { return this_ == nil }
func (*CheckDebuggerStatement) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CheckDebuggerStatement", "CheckDebuggerStatement"
}
func (this_ *CheckDebuggerStatement) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*CheckDebuggerStatement) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CheckDebuggerStatement", "com.google.javascript.jscomp.CheckDebuggerStatement", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckDebuggerStatement](v_)
	})
}
func (*CheckDebuggerStatement) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*CheckDebuggerStatement) Implements_jscomp_CompilerPass()                            {}
func (*CheckDebuggerStatement) Implements_jscomp_NodeTraversal_Callback()                  {}

var static_CheckDebuggerStatement_DEBUGGER_STATEMENT_PRESENT *DiagnosticType

func StaticGet_CheckDebuggerStatement_DEBUGGER_STATEMENT_PRESENT() *DiagnosticType {
	clinit_CheckDebuggerStatement()
	return static_CheckDebuggerStatement_DEBUGGER_STATEMENT_PRESENT
}
func StaticSet_CheckDebuggerStatement_DEBUGGER_STATEMENT_PRESENT(v *DiagnosticType) {
	clinit_CheckDebuggerStatement()
	static_CheckDebuggerStatement_DEBUGGER_STATEMENT_PRESENT = v
}

var clinitDone_CheckDebuggerStatement uint32
var clinitState_CheckDebuggerStatement j2gort.ClinitState

func clinit_CheckDebuggerStatement() {
	if atomic.LoadUint32(&clinitDone_CheckDebuggerStatement) != 0 {
		return
	}
	clinitSlow_CheckDebuggerStatement()
}

func clinitSlow_CheckDebuggerStatement() {
	if !j2gort.EnterClinit(&clinitState_CheckDebuggerStatement, &clinitDone_CheckDebuggerStatement) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckDebuggerStatement, &clinitDone_CheckDebuggerStatement)
	static_CheckDebuggerStatement_DEBUGGER_STATEMENT_PRESENT = DiagnosticType_Disabled("JSC_DEBUGGER_STATEMENT_PRESENT", "Using the debugger statement can halt your application if the user has a JavaScript debugger running.\nTo disable this check when you want to do debugging, you can suppress this message like this:\n  /** @suppress '{'checkDebuggerStatement'}' */\n  debugger;")
}

func (this_ *CheckDebuggerStatement) initFields_CheckDebuggerStatement() {
	_ = this_
}

// Java constructor: CheckDebuggerStatement(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewCheckDebuggerStatement(compiler AbstractCompiler) *CheckDebuggerStatement {
	clinit_CheckDebuggerStatement()
	this_ := &CheckDebuggerStatement{}
	InitCheckDebuggerStatement(this_, compiler)
	return this_
}

func InitCheckDebuggerStatement(this_ *CheckDebuggerStatement, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_CheckDebuggerStatement()
	this_.compiler = compiler
}

// Java method: CheckDebuggerStatement.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *CheckDebuggerStatement) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
}

// Java method: CheckDebuggerStatement.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckDebuggerStatement) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsDebugger() {
		t.Report_3(n, StaticGet_CheckDebuggerStatement_DEBUGGER_STATEMENT_PRESENT(), []string{})
	}
}
func (*CheckDebuggerStatement) AlwaysShouldTraverse_()      {}
func (*CheckDebuggerStatement) CanReleaseSyntacticScopes_() {}
