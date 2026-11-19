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
// Port of com/google/javascript/jscomp/lint/CheckEmptyStatements.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.lint.CheckEmptyStatements
type CheckEmptyStatements struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler `gson:"compiler"`
}

func (this_ *CheckEmptyStatements) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckEmptyStatements) IsJavaNil_() bool { return this_ == nil }
func (*CheckEmptyStatements) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.lint.CheckEmptyStatements", "CheckEmptyStatements"
}
func (this_ *CheckEmptyStatements) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*CheckEmptyStatements) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.lint.CheckEmptyStatements", "com.google.javascript.jscomp.lint.CheckEmptyStatements", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckEmptyStatements](v_)
	})
}
func (*CheckEmptyStatements) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*CheckEmptyStatements) Implements_jscomp_CompilerPass()                            {}
func (*CheckEmptyStatements) Implements_jscomp_NodeTraversal_Callback()                  {}

var static_CheckEmptyStatements_USELESS_EMPTY_STATEMENT *DiagnosticType

func StaticGet_CheckEmptyStatements_USELESS_EMPTY_STATEMENT() *DiagnosticType {
	clinit_CheckEmptyStatements()
	return static_CheckEmptyStatements_USELESS_EMPTY_STATEMENT
}
func StaticSet_CheckEmptyStatements_USELESS_EMPTY_STATEMENT(v *DiagnosticType) {
	clinit_CheckEmptyStatements()
	static_CheckEmptyStatements_USELESS_EMPTY_STATEMENT = v
}

var clinitDone_CheckEmptyStatements uint32
var clinitState_CheckEmptyStatements j2gort.ClinitState

func clinit_CheckEmptyStatements() {
	if atomic.LoadUint32(&clinitDone_CheckEmptyStatements) != 0 {
		return
	}
	clinitSlow_CheckEmptyStatements()
}

func clinitSlow_CheckEmptyStatements() {
	if !j2gort.EnterClinit(&clinitState_CheckEmptyStatements, &clinitDone_CheckEmptyStatements) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckEmptyStatements, &clinitDone_CheckEmptyStatements)
	static_CheckEmptyStatements_USELESS_EMPTY_STATEMENT = DiagnosticType_Disabled("JSC_USELESS_EMPTY_STATEMENT", "Useless empty statement. Remove semicolon.")
}

func (this_ *CheckEmptyStatements) initFields_CheckEmptyStatements() {
	_ = this_
}

// Java constructor: CheckEmptyStatements(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewCheckEmptyStatements(compiler AbstractCompiler) *CheckEmptyStatements {
	clinit_CheckEmptyStatements()
	this_ := &CheckEmptyStatements{}
	InitCheckEmptyStatements(this_, compiler)
	return this_
}

func InitCheckEmptyStatements(this_ *CheckEmptyStatements, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_CheckEmptyStatements()
	this_.compiler = compiler
}

// Java method: CheckEmptyStatements.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *CheckEmptyStatements) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
}

// Java method: CheckEmptyStatements.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckEmptyStatements) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsEmpty() && NodeUtil_IsStatement(n) {
		t.Report_3(n, StaticGet_CheckEmptyStatements_USELESS_EMPTY_STATEMENT(), []string{})
	}
}
func (*CheckEmptyStatements) AlwaysShouldTraverse_()      {}
func (*CheckEmptyStatements) CanReleaseSyntacticScopes_() {}
