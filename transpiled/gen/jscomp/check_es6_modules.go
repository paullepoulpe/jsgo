/*
 * Copyright 2018 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/lint/CheckEs6Modules.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.lint.CheckEs6Modules
type CheckEs6Modules struct {
	compiler         AbstractCompiler                          `gson:"compiler"`
	importSpecifiers j2gort.OrderedMap[string, pkg_rhino.Node] `gson:"importSpecifiers"`
	idHash_          atomic.Int32
}

func (this_ *CheckEs6Modules) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckEs6Modules) IsJavaNil_() bool { return this_ == nil }
func (*CheckEs6Modules) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.lint.CheckEs6Modules", "CheckEs6Modules"
}
func (this_ *CheckEs6Modules) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CheckEs6Modules) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.lint.CheckEs6Modules", "com.google.javascript.jscomp.lint.CheckEs6Modules", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckEs6Modules](v_)
	})
}
func (*CheckEs6Modules) Implements_jscomp_NodeTraversal_Callback() {}
func (*CheckEs6Modules) Implements_jscomp_CompilerPass()           {}

var static_CheckEs6Modules_DUPLICATE_IMPORT *DiagnosticType

func StaticGet_CheckEs6Modules_DUPLICATE_IMPORT() *DiagnosticType {
	clinit_CheckEs6Modules()
	return static_CheckEs6Modules_DUPLICATE_IMPORT
}
func StaticSet_CheckEs6Modules_DUPLICATE_IMPORT(v *DiagnosticType) {
	clinit_CheckEs6Modules()
	static_CheckEs6Modules_DUPLICATE_IMPORT = v
}

var static_CheckEs6Modules_NO_DEFAULT_EXPORT *DiagnosticType

func StaticGet_CheckEs6Modules_NO_DEFAULT_EXPORT() *DiagnosticType {
	clinit_CheckEs6Modules()
	return static_CheckEs6Modules_NO_DEFAULT_EXPORT
}
func StaticSet_CheckEs6Modules_NO_DEFAULT_EXPORT(v *DiagnosticType) {
	clinit_CheckEs6Modules()
	static_CheckEs6Modules_NO_DEFAULT_EXPORT = v
}

var clinitDone_CheckEs6Modules uint32
var clinitState_CheckEs6Modules j2gort.ClinitState

func clinit_CheckEs6Modules() {
	if atomic.LoadUint32(&clinitDone_CheckEs6Modules) != 0 {
		return
	}
	clinitSlow_CheckEs6Modules()
}

func clinitSlow_CheckEs6Modules() {
	if !j2gort.EnterClinit(&clinitState_CheckEs6Modules, &clinitDone_CheckEs6Modules) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckEs6Modules, &clinitDone_CheckEs6Modules)
	static_CheckEs6Modules_DUPLICATE_IMPORT = DiagnosticType_Warning("JSC_DUPLICATE_IMPORT", "The module \"{0}\" has already been imported at {1}, {2}.")
	static_CheckEs6Modules_NO_DEFAULT_EXPORT = DiagnosticType_Warning("JSC_DEFAULT_EXPORT", "Do not use the default export. There is no way to force consistent naming when imported.")
}

func (this_ *CheckEs6Modules) initFields_CheckEs6Modules() {
	_ = this_
	this_.importSpecifiers = j2gort.NewOrderedMap[string, pkg_rhino.Node]()
}

// Java constructor: CheckEs6Modules(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewCheckEs6Modules(compiler AbstractCompiler) *CheckEs6Modules {
	clinit_CheckEs6Modules()
	this_ := &CheckEs6Modules{}
	InitCheckEs6Modules(this_, compiler)
	return this_
}

func InitCheckEs6Modules(this_ *CheckEs6Modules, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_CheckEs6Modules()
	this_.compiler = compiler
}

// Java method: CheckEs6Modules.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckEs6Modules) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsModuleBody() {
		this_.importSpecifiers.Clear()
	}
}

// Java method: CheckEs6Modules.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckEs6Modules) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(89) /* Token.ROOT */, int32(123) /* Token.MODULE_BODY */ :
		return true
	case int32(94) /* Token.SCRIPT */ :
		return n.GetBooleanProp(pkg_rhino.StaticGet_Node_ES6_MODULE())
	case int32(116) /* Token.IMPORT */ :
		this_.priv_CheckEs6Modules_visitImport_6572(t, n)
		return false
	case int32(120) /* Token.EXPORT */ :
		this_.priv_CheckEs6Modules_visitExport_6573(t, n)
		return false
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: CheckEs6Modules.visitImport(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node importNode)
func (this_ *CheckEs6Modules) priv_CheckEs6Modules_visitImport_6572(t *NodeTraversal, importNode pkg_rhino.Node) {
	_ = t
	_ = importNode
	var specifier string = importNode.GetLastChild().GetString()
	_ = specifier
	var duplicateImport pkg_rhino.Node = this_.importSpecifiers.PutIfAbsent(specifier, importNode)
	_ = duplicateImport
	if duplicateImport != 0 {
		t.Report_3(importNode, StaticGet_CheckEs6Modules_DUPLICATE_IMPORT(), []string{specifier, ("" + j2gort.StrValInt32(duplicateImport.GetLineno())), ("" + j2gort.StrValInt32(duplicateImport.GetCharno()))})
	}
}

// Java method: CheckEs6Modules.visitExport(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node export)
func (this_ *CheckEs6Modules) priv_CheckEs6Modules_visitExport_6573(t *NodeTraversal, export pkg_rhino.Node) {
	_ = t
	_ = export
	if export.GetBooleanProp(pkg_rhino.StaticGet_Node_EXPORT_DEFAULT()) {
		t.Report_3(export, StaticGet_CheckEs6Modules_NO_DEFAULT_EXPORT(), []string{})
	}
}

// Java method: CheckEs6Modules.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *CheckEs6Modules) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
}
func (*CheckEs6Modules) CanReleaseSyntacticScopes_() {}
