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
// Port of com/google/javascript/jscomp/lint/CheckVar.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.lint.CheckVar
type CheckVar struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler `gson:"compiler"`
}

func (this_ *CheckVar) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckVar) IsJavaNil_() bool { return this_ == nil }
func (*CheckVar) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.lint.CheckVar", "CheckVar"
}
func (this_ *CheckVar) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*CheckVar) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.lint.CheckVar", "com.google.javascript.jscomp.lint.CheckVar", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckVar](v_)
	})
}
func (*CheckVar) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*CheckVar) Implements_jscomp_CompilerPass()                            {}
func (*CheckVar) Implements_jscomp_NodeTraversal_Callback()                  {}

var static_CheckVar_VAR *DiagnosticType

func StaticGet_CheckVar_VAR() *DiagnosticType {
	clinit_CheckVar()
	return static_CheckVar_VAR
}
func StaticSet_CheckVar_VAR(v *DiagnosticType) {
	clinit_CheckVar()
	static_CheckVar_VAR = v
}

var clinitDone_CheckVar uint32
var clinitState_CheckVar j2gort.ClinitState

func clinit_CheckVar() {
	if atomic.LoadUint32(&clinitDone_CheckVar) != 0 {
		return
	}
	clinitSlow_CheckVar()
}

func clinitSlow_CheckVar() {
	if !j2gort.EnterClinit(&clinitState_CheckVar, &clinitDone_CheckVar) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckVar, &clinitDone_CheckVar)
	static_CheckVar_VAR = DiagnosticType_Disabled("JSC_VAR", "Using `var` (prefer `const` or `let`).")
}

func (this_ *CheckVar) initFields_CheckVar() {
	_ = this_
}

// Java constructor: CheckVar(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewCheckVar(compiler AbstractCompiler) *CheckVar {
	clinit_CheckVar()
	this_ := &CheckVar{}
	InitCheckVar(this_, compiler)
	return this_
}

func InitCheckVar(this_ *CheckVar, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_CheckVar()
	this_.compiler = compiler
}

// Java method: CheckVar.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *CheckVar) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
}

// Java method: CheckVar.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckVar) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsVar() {
		t.Report_3(n, StaticGet_CheckVar_VAR(), []string{})
	}
}
func (*CheckVar) AlwaysShouldTraverse_()      {}
func (*CheckVar) CanReleaseSyntacticScopes_() {}
