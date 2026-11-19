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
// Port of com/google/javascript/jscomp/CheckTypeImportCodeReferences.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.CheckTypeImportCodeReferences
type CheckTypeImportCodeReferences struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler `gson:"compiler"`
}

func (this_ *CheckTypeImportCodeReferences) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckTypeImportCodeReferences) IsJavaNil_() bool { return this_ == nil }
func (*CheckTypeImportCodeReferences) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CheckTypeImportCodeReferences", "CheckTypeImportCodeReferences"
}
func (this_ *CheckTypeImportCodeReferences) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*CheckTypeImportCodeReferences) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CheckTypeImportCodeReferences", "com.google.javascript.jscomp.CheckTypeImportCodeReferences", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckTypeImportCodeReferences](v_)
	})
}
func (*CheckTypeImportCodeReferences) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*CheckTypeImportCodeReferences) Implements_jscomp_CompilerPass()                            {}
func (*CheckTypeImportCodeReferences) Implements_jscomp_NodeTraversal_Callback()                  {}

var static_CheckTypeImportCodeReferences_TYPE_IMPORT_CODE_REFERENCE *DiagnosticType

func StaticGet_CheckTypeImportCodeReferences_TYPE_IMPORT_CODE_REFERENCE() *DiagnosticType {
	clinit_CheckTypeImportCodeReferences()
	return static_CheckTypeImportCodeReferences_TYPE_IMPORT_CODE_REFERENCE
}
func StaticSet_CheckTypeImportCodeReferences_TYPE_IMPORT_CODE_REFERENCE(v *DiagnosticType) {
	clinit_CheckTypeImportCodeReferences()
	static_CheckTypeImportCodeReferences_TYPE_IMPORT_CODE_REFERENCE = v
}

var static_CheckTypeImportCodeReferences_GOOG_REQUIRE_TYPE pkg_rhino.Node

func StaticGet_CheckTypeImportCodeReferences_GOOG_REQUIRE_TYPE() pkg_rhino.Node {
	clinit_CheckTypeImportCodeReferences()
	return static_CheckTypeImportCodeReferences_GOOG_REQUIRE_TYPE
}
func StaticSet_CheckTypeImportCodeReferences_GOOG_REQUIRE_TYPE(v pkg_rhino.Node) {
	clinit_CheckTypeImportCodeReferences()
	static_CheckTypeImportCodeReferences_GOOG_REQUIRE_TYPE = v
}

var clinitDone_CheckTypeImportCodeReferences uint32
var clinitState_CheckTypeImportCodeReferences j2gort.ClinitState

func clinit_CheckTypeImportCodeReferences() {
	if atomic.LoadUint32(&clinitDone_CheckTypeImportCodeReferences) != 0 {
		return
	}
	clinitSlow_CheckTypeImportCodeReferences()
}

func clinitSlow_CheckTypeImportCodeReferences() {
	if !j2gort.EnterClinit(&clinitState_CheckTypeImportCodeReferences, &clinitDone_CheckTypeImportCodeReferences) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckTypeImportCodeReferences, &clinitDone_CheckTypeImportCodeReferences)
	static_CheckTypeImportCodeReferences_TYPE_IMPORT_CODE_REFERENCE = DiagnosticType_Error("JSC_TYPE_IMPORT_CODE_REFERENCE", "Cannot reference goog.requireType()''d name {0} outside of a type annotation.")
	static_CheckTypeImportCodeReferences_GOOG_REQUIRE_TYPE = pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name("goog"), "requireType")
}

func (this_ *CheckTypeImportCodeReferences) initFields_CheckTypeImportCodeReferences() {
	_ = this_
}

// Java constructor: CheckTypeImportCodeReferences(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewCheckTypeImportCodeReferences(compiler AbstractCompiler) *CheckTypeImportCodeReferences {
	clinit_CheckTypeImportCodeReferences()
	this_ := &CheckTypeImportCodeReferences{}
	InitCheckTypeImportCodeReferences(this_, compiler)
	return this_
}

func InitCheckTypeImportCodeReferences(this_ *CheckTypeImportCodeReferences, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_CheckTypeImportCodeReferences()
	this_.compiler = compiler
}

// Java method: CheckTypeImportCodeReferences.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *CheckTypeImportCodeReferences) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
}

// Java method: CheckTypeImportCodeReferences.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckTypeImportCodeReferences) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if !(n.IsName()) {
		return
	}
	var var_ *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(n.GetString()))
	_ = var_
	if var_ == nil {
		return
	}
	var declNameNode pkg_rhino.Node = var_.GetNameNode()
	_ = declNameNode
	if ((declNameNode == 0) || (declNameNode == n)) || !(NodeUtil_IsDeclarationLValue(declNameNode)) {
		return
	}
	var decl pkg_rhino.Node
	_ = decl
	if declNameNode.GetParent().IsStringKey() {
		decl = declNameNode.GetGrandparent().GetGrandparent()
	} else {
		decl = declNameNode.GetParent()
	}
	if ((!(NodeUtil_IsNameDeclaration(decl)) || !(decl.HasOneChild())) || !(decl.GetFirstChild().HasChildren())) || !(decl.GetFirstChild().GetLastChild().IsCall()) {
		return
	}
	var callNode pkg_rhino.Node = decl.GetFirstChild().GetLastChild()
	_ = callNode
	if !(callNode.GetFirstChild().MatchesQualifiedName_1_Node(StaticGet_CheckTypeImportCodeReferences_GOOG_REQUIRE_TYPE())) {
		return
	}
	t.Report_3(n, StaticGet_CheckTypeImportCodeReferences_TYPE_IMPORT_CODE_REFERENCE(), []string{n.GetString()})
}
func (*CheckTypeImportCodeReferences) AlwaysShouldTraverse_()      {}
func (*CheckTypeImportCodeReferences) CanReleaseSyntacticScopes_() {}
