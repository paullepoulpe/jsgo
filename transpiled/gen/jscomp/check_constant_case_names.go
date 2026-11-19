/*
 * Copyright 2019 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/lint/CheckConstantCaseNames.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.lint.CheckConstantCaseNames
type CheckConstantCaseNames struct {
	compiler              AbstractCompiler                          `gson:"compiler"`
	convention            CodingConvention                          `gson:"convention"`
	invalidNamesPerModule j2gort.OrderedMap[string, pkg_rhino.Node] `gson:"invalidNamesPerModule"`
	reassignedNames       j2gort.OrderedSet[string]                 `gson:"reassignedNames"`
	idHash_               atomic.Int32
}

func (this_ *CheckConstantCaseNames) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckConstantCaseNames) IsJavaNil_() bool { return this_ == nil }
func (*CheckConstantCaseNames) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.lint.CheckConstantCaseNames", "CheckConstantCaseNames"
}
func (this_ *CheckConstantCaseNames) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CheckConstantCaseNames) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.lint.CheckConstantCaseNames", "com.google.javascript.jscomp.lint.CheckConstantCaseNames", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckConstantCaseNames](v_)
	})
}
func (*CheckConstantCaseNames) Implements_jscomp_NodeTraversal_Callback() {}
func (*CheckConstantCaseNames) Implements_jscomp_CompilerPass()           {}

var static_CheckConstantCaseNames_MISSING_CONST_PROPERTY *DiagnosticType

func StaticGet_CheckConstantCaseNames_MISSING_CONST_PROPERTY() *DiagnosticType {
	clinit_CheckConstantCaseNames()
	return static_CheckConstantCaseNames_MISSING_CONST_PROPERTY
}
func StaticSet_CheckConstantCaseNames_MISSING_CONST_PROPERTY(v *DiagnosticType) {
	clinit_CheckConstantCaseNames()
	static_CheckConstantCaseNames_MISSING_CONST_PROPERTY = v
}

var static_CheckConstantCaseNames_REASSIGNED_CONSTANT_CASE_NAME *DiagnosticType

func StaticGet_CheckConstantCaseNames_REASSIGNED_CONSTANT_CASE_NAME() *DiagnosticType {
	clinit_CheckConstantCaseNames()
	return static_CheckConstantCaseNames_REASSIGNED_CONSTANT_CASE_NAME
}
func StaticSet_CheckConstantCaseNames_REASSIGNED_CONSTANT_CASE_NAME(v *DiagnosticType) {
	clinit_CheckConstantCaseNames()
	static_CheckConstantCaseNames_REASSIGNED_CONSTANT_CASE_NAME = v
}

var clinitDone_CheckConstantCaseNames uint32
var clinitState_CheckConstantCaseNames j2gort.ClinitState

func clinit_CheckConstantCaseNames() {
	if atomic.LoadUint32(&clinitDone_CheckConstantCaseNames) != 0 {
		return
	}
	clinitSlow_CheckConstantCaseNames()
}

func clinitSlow_CheckConstantCaseNames() {
	if !j2gort.EnterClinit(&clinitState_CheckConstantCaseNames, &clinitDone_CheckConstantCaseNames) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckConstantCaseNames, &clinitDone_CheckConstantCaseNames)
	static_CheckConstantCaseNames_MISSING_CONST_PROPERTY = DiagnosticType_Disabled("JSC_MISSING_CONST_ON_CONSTANT_CASE", "CONSTANT_CASE name \"{0}\" is constant-by-convention, so must be explicitly `const` or @const")
	static_CheckConstantCaseNames_REASSIGNED_CONSTANT_CASE_NAME = DiagnosticType_Disabled("JSC_REASSIGNED_CONSTANT_CASE_NAME", "CONSTANT_CASE name \"{0}\" is constant-by-convention but is reassigned. Use camelCase instead.")
}

func (this_ *CheckConstantCaseNames) initFields_CheckConstantCaseNames() {
	_ = this_
	this_.invalidNamesPerModule = j2gort.NewOrderedMap[string, pkg_rhino.Node]()
	this_.reassignedNames = j2gort.NewOrderedSet[string]()
}

// Java constructor: CheckConstantCaseNames(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewCheckConstantCaseNames(compiler AbstractCompiler) *CheckConstantCaseNames {
	clinit_CheckConstantCaseNames()
	this_ := &CheckConstantCaseNames{}
	InitCheckConstantCaseNames(this_, compiler)
	return this_
}

func InitCheckConstantCaseNames(this_ *CheckConstantCaseNames, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_CheckConstantCaseNames()
	this_.compiler = compiler
	this_.convention = compiler.GetCodingConvention()
}

// Java method: CheckConstantCaseNames.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *CheckConstantCaseNames) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
}

// Java method: CheckConstantCaseNames.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckConstantCaseNames) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	if n.IsScript() {
		return (n.HasChildren() && n.GetFirstChild().IsModuleBody())
	}
	return true
}

// Java method: CheckConstantCaseNames.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckConstantCaseNames) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsModuleBody() {
		this_.priv_CheckConstantCaseNames_reportWarningsAndClear_6546()
		return
	}
	switch n.TokenSwitchOrd_() {
	case int32(84) /* Token.VAR */, int32(111) /* Token.LET */ :
		if !(t.InModuleScope()) {
			return
		}
		var info *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
		_ = info
		if (info != nil) && info.HasConstAnnotation() {
			break
		}
		NodeUtil_VisitLhsNodesInNode(n, j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_1 any) {
			name := j2gort.CastNode[pkg_rhino.Node](lp_1)
			_ = name
			if this_.convention.IsConstant(name.GetString()) {
				this_.invalidNamesPerModule.Put(name.GetString(), name)
			}
		}))
	case int32(32) /* Token.NAME */ :
		if !(this_.invalidNamesPerModule.ContainsKey(n.GetString())) {
			return
		}
		if !(NodeUtil_IsLValue(n)) {
			return
		}
		var v *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(n.GetString()))
		_ = v
		if v.GetScopeRoot().IsModuleBody() {
			this_.reassignedNames.Add(n.GetString())
		}
	default:
	}
}

// Java method: CheckConstantCaseNames.reportWarningsAndClear()
func (this_ *CheckConstantCaseNames) priv_CheckConstantCaseNames_reportWarningsAndClear_6546() {
	{
		cur_2 := this_.invalidNamesPerModule.IterValues()
		for cur_2.HasNext() {

			nameNode := cur_2.
				Next()
			_ = nameNode
			var name string = nameNode.GetString()
			_ = name
			if this_.reassignedNames.Contains(name) {
				this_.compiler.Report_1(JSError_Make_3(nameNode, StaticGet_CheckConstantCaseNames_REASSIGNED_CONSTANT_CASE_NAME(), []string{name}))
			} else {
				this_.compiler.Report_1(JSError_Make_3(nameNode, StaticGet_CheckConstantCaseNames_MISSING_CONST_PROPERTY(), []string{name}))
			}
		}
	}
	this_.invalidNamesPerModule = j2gort.NewOrderedMap[string, pkg_rhino.Node]()
	this_.reassignedNames = j2gort.NewOrderedSet[string]()
}
func (*CheckConstantCaseNames) CanReleaseSyntacticScopes_() {}
