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
// Port of com/google/javascript/jscomp/lint/CheckDuplicateCase.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.lint.CheckDuplicateCase
type CheckDuplicateCase struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler `gson:"compiler"`
}

func (this_ *CheckDuplicateCase) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckDuplicateCase) IsJavaNil_() bool { return this_ == nil }
func (*CheckDuplicateCase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.lint.CheckDuplicateCase", "CheckDuplicateCase"
}
func (this_ *CheckDuplicateCase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*CheckDuplicateCase) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.lint.CheckDuplicateCase", "com.google.javascript.jscomp.lint.CheckDuplicateCase", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckDuplicateCase](v_)
	})
}
func (*CheckDuplicateCase) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*CheckDuplicateCase) Implements_jscomp_CompilerPass()                            {}
func (*CheckDuplicateCase) Implements_jscomp_NodeTraversal_Callback()                  {}

var static_CheckDuplicateCase_DUPLICATE_CASE *DiagnosticType

func StaticGet_CheckDuplicateCase_DUPLICATE_CASE() *DiagnosticType {
	clinit_CheckDuplicateCase()
	return static_CheckDuplicateCase_DUPLICATE_CASE
}
func StaticSet_CheckDuplicateCase_DUPLICATE_CASE(v *DiagnosticType) {
	clinit_CheckDuplicateCase()
	static_CheckDuplicateCase_DUPLICATE_CASE = v
}

var clinitDone_CheckDuplicateCase uint32
var clinitState_CheckDuplicateCase j2gort.ClinitState

func clinit_CheckDuplicateCase() {
	if atomic.LoadUint32(&clinitDone_CheckDuplicateCase) != 0 {
		return
	}
	clinitSlow_CheckDuplicateCase()
}

func clinitSlow_CheckDuplicateCase() {
	if !j2gort.EnterClinit(&clinitState_CheckDuplicateCase, &clinitDone_CheckDuplicateCase) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckDuplicateCase, &clinitDone_CheckDuplicateCase)
	static_CheckDuplicateCase_DUPLICATE_CASE = DiagnosticType_Warning("JSC_DUPLICATE_CASE", "Duplicate case in a switch statement.")
}

func (this_ *CheckDuplicateCase) initFields_CheckDuplicateCase() {
	_ = this_
}

// Java constructor: CheckDuplicateCase(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewCheckDuplicateCase(compiler AbstractCompiler) *CheckDuplicateCase {
	clinit_CheckDuplicateCase()
	this_ := &CheckDuplicateCase{}
	InitCheckDuplicateCase(this_, compiler)
	return this_
}

func InitCheckDuplicateCase(this_ *CheckDuplicateCase, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_CheckDuplicateCase()
	this_.compiler = compiler
}

// Java method: CheckDuplicateCase.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *CheckDuplicateCase) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
}

// Java method: CheckDuplicateCase.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckDuplicateCase) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsSwitch() {
		var switchBody pkg_rhino.Node = n.GetSecondChild()
		_ = switchBody
		var cases j2gort.LocalOrderedSet[string] = j2gort.NewLocalOrderedSet[string]()
		_ = cases
		{
			var curr pkg_rhino.Node = switchBody.GetFirstChild()
			_ = curr
			for ; curr != 0; curr = curr.GetNext() {
				var source string = this_.compiler.ToSource_1_Node(curr.GetFirstChild())
				_ = source
				if !(cases.Add(source)) {
					t.Report_3(curr, StaticGet_CheckDuplicateCase_DUPLICATE_CASE(), []string{})
				}
			}
		}
	}
}
func (*CheckDuplicateCase) AlwaysShouldTraverse_()      {}
func (*CheckDuplicateCase) CanReleaseSyntacticScopes_() {}
