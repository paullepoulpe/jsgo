/*
 * Copyright 2014 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/lint/CheckInterfaces.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.lint.CheckInterfaces
type CheckInterfaces struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler `gson:"compiler"`
}

func (this_ *CheckInterfaces) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckInterfaces) IsJavaNil_() bool { return this_ == nil }
func (*CheckInterfaces) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.lint.CheckInterfaces", "CheckInterfaces"
}
func (this_ *CheckInterfaces) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*CheckInterfaces) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.lint.CheckInterfaces", "com.google.javascript.jscomp.lint.CheckInterfaces", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckInterfaces](v_)
	})
}
func (*CheckInterfaces) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*CheckInterfaces) Implements_jscomp_CompilerPass()                            {}
func (*CheckInterfaces) Implements_jscomp_NodeTraversal_Callback()                  {}

func StaticGet_CheckInterfaces_ANONYMOUS_CLASSNAME() string { return "<anonymous>" }

var static_CheckInterfaces_NON_DECLARATION_STATEMENT_IN_INTERFACE *DiagnosticType

func StaticGet_CheckInterfaces_NON_DECLARATION_STATEMENT_IN_INTERFACE() *DiagnosticType {
	clinit_CheckInterfaces()
	return static_CheckInterfaces_NON_DECLARATION_STATEMENT_IN_INTERFACE
}
func StaticSet_CheckInterfaces_NON_DECLARATION_STATEMENT_IN_INTERFACE(v *DiagnosticType) {
	clinit_CheckInterfaces()
	static_CheckInterfaces_NON_DECLARATION_STATEMENT_IN_INTERFACE = v
}

var static_CheckInterfaces_MISSING_JSDOC_IN_DECLARATION_STATEMENT *DiagnosticType

func StaticGet_CheckInterfaces_MISSING_JSDOC_IN_DECLARATION_STATEMENT() *DiagnosticType {
	clinit_CheckInterfaces()
	return static_CheckInterfaces_MISSING_JSDOC_IN_DECLARATION_STATEMENT
}
func StaticSet_CheckInterfaces_MISSING_JSDOC_IN_DECLARATION_STATEMENT(v *DiagnosticType) {
	clinit_CheckInterfaces()
	static_CheckInterfaces_MISSING_JSDOC_IN_DECLARATION_STATEMENT = v
}

var static_CheckInterfaces_INTERFACE_CLASS_NONSTATIC_METHOD_NOT_EMPTY *DiagnosticType

func StaticGet_CheckInterfaces_INTERFACE_CLASS_NONSTATIC_METHOD_NOT_EMPTY() *DiagnosticType {
	clinit_CheckInterfaces()
	return static_CheckInterfaces_INTERFACE_CLASS_NONSTATIC_METHOD_NOT_EMPTY
}
func StaticSet_CheckInterfaces_INTERFACE_CLASS_NONSTATIC_METHOD_NOT_EMPTY(v *DiagnosticType) {
	clinit_CheckInterfaces()
	static_CheckInterfaces_INTERFACE_CLASS_NONSTATIC_METHOD_NOT_EMPTY = v
}

var static_CheckInterfaces_INTERFACE_CONSTRUCTOR_SHOULD_NOT_TAKE_ARGS *DiagnosticType

func StaticGet_CheckInterfaces_INTERFACE_CONSTRUCTOR_SHOULD_NOT_TAKE_ARGS() *DiagnosticType {
	clinit_CheckInterfaces()
	return static_CheckInterfaces_INTERFACE_CONSTRUCTOR_SHOULD_NOT_TAKE_ARGS
}
func StaticSet_CheckInterfaces_INTERFACE_CONSTRUCTOR_SHOULD_NOT_TAKE_ARGS(v *DiagnosticType) {
	clinit_CheckInterfaces()
	static_CheckInterfaces_INTERFACE_CONSTRUCTOR_SHOULD_NOT_TAKE_ARGS = v
}

var static_CheckInterfaces_STATIC_MEMBER_FUNCTION_IN_INTERFACE_CLASS *DiagnosticType

func StaticGet_CheckInterfaces_STATIC_MEMBER_FUNCTION_IN_INTERFACE_CLASS() *DiagnosticType {
	clinit_CheckInterfaces()
	return static_CheckInterfaces_STATIC_MEMBER_FUNCTION_IN_INTERFACE_CLASS
}
func StaticSet_CheckInterfaces_STATIC_MEMBER_FUNCTION_IN_INTERFACE_CLASS(v *DiagnosticType) {
	clinit_CheckInterfaces()
	static_CheckInterfaces_STATIC_MEMBER_FUNCTION_IN_INTERFACE_CLASS = v
}

var static_CheckInterfaces_INTERFACE_DEFINED_WITH_EXTENDS *DiagnosticType

func StaticGet_CheckInterfaces_INTERFACE_DEFINED_WITH_EXTENDS() *DiagnosticType {
	clinit_CheckInterfaces()
	return static_CheckInterfaces_INTERFACE_DEFINED_WITH_EXTENDS
}
func StaticSet_CheckInterfaces_INTERFACE_DEFINED_WITH_EXTENDS(v *DiagnosticType) {
	clinit_CheckInterfaces()
	static_CheckInterfaces_INTERFACE_DEFINED_WITH_EXTENDS = v
}

var clinitDone_CheckInterfaces uint32
var clinitState_CheckInterfaces j2gort.ClinitState

func clinit_CheckInterfaces() {
	if atomic.LoadUint32(&clinitDone_CheckInterfaces) != 0 {
		return
	}
	clinitSlow_CheckInterfaces()
}

func clinitSlow_CheckInterfaces() {
	if !j2gort.EnterClinit(&clinitState_CheckInterfaces, &clinitDone_CheckInterfaces) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckInterfaces, &clinitDone_CheckInterfaces)
	static_CheckInterfaces_NON_DECLARATION_STATEMENT_IN_INTERFACE = DiagnosticType_Disabled("JSC_NON_DECLARATION_STATEMENT_IN_INTERFACE", "@interface or @record functions should not contain statements other than field declarations")
	static_CheckInterfaces_MISSING_JSDOC_IN_DECLARATION_STATEMENT = DiagnosticType_Disabled("JSC_MISSING_JSDOC_IN_DECLARATION_STATEMENT", "@interface or @record functions must contain JSDoc for each field declaration.")
	static_CheckInterfaces_INTERFACE_CLASS_NONSTATIC_METHOD_NOT_EMPTY = DiagnosticType_Disabled("JSC_INTERFACE_CLASS_NONSTATIC_METHOD_NOT_EMPTY", "interface methods must have an empty body")
	static_CheckInterfaces_INTERFACE_CONSTRUCTOR_SHOULD_NOT_TAKE_ARGS = DiagnosticType_Disabled("JSC_INTERFACE_CONSTRUCTOR_SHOULD_NOT_TAKE_ARGS", "Interface constructors should not take any arguments")
	static_CheckInterfaces_STATIC_MEMBER_FUNCTION_IN_INTERFACE_CLASS = DiagnosticType_Disabled("JSC_STATIC_MEMBER_FUNCTION_IN_INTERFACE_CLASS", "Interface class should not have static member functions. Consider pulling out the static method into a flat name as {0}_{1}")
	static_CheckInterfaces_INTERFACE_DEFINED_WITH_EXTENDS = DiagnosticType_Disabled("JSC_INTERFACE_DEFINED_WITH_EXTENDS", "Interface/Record class should use the `@extends` annotation instead of extends keyword.")
}

func (this_ *CheckInterfaces) initFields_CheckInterfaces() {
	_ = this_
}

// Java constructor: CheckInterfaces(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewCheckInterfaces(compiler AbstractCompiler) *CheckInterfaces {
	clinit_CheckInterfaces()
	this_ := &CheckInterfaces{}
	InitCheckInterfaces(this_, compiler)
	return this_
}

func InitCheckInterfaces(this_ *CheckInterfaces, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_CheckInterfaces()
	this_.compiler = compiler
}

// Java method: CheckInterfaces.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *CheckInterfaces) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
}

// Java static method: CheckInterfaces.isInterface(com.google.javascript.rhino.JSDocInfo jsDoc)
func CheckInterfaces_IsInterface(jsDoc *pkg_rhino.JSDocInfo) bool {
	clinit_CheckInterfaces()
	_ = jsDoc
	return ((jsDoc != nil) && jsDoc.IsInterface())
}

// Java method: CheckInterfaces.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckInterfaces) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(73) /* Token.FUNCTION */ :
		var jsdoc *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(n)
		_ = jsdoc
		if CheckInterfaces_IsInterface(jsdoc) {
			CheckInterfaces_CheckInterfaceConstructorArgs(t, n)
			CheckInterfaces_CheckConstructorBlock(t, n)
		}
	case int32(105) /* Token.CLASS */ :
		var jsdoc_2 *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(n)
		_ = jsdoc_2
		if CheckInterfaces_IsInterface(jsdoc_2) {
			if (n.GetSecondChild() != 0) && n.GetSecondChild().IsName() {
				t.Report_3(n, StaticGet_CheckInterfaces_INTERFACE_DEFINED_WITH_EXTENDS(), []string{})
			}
			var ctorDef pkg_rhino.Node = NodeUtil_GetEs6ClassConstructorMemberFunctionDef(n)
			_ = ctorDef
			if ctorDef != 0 {
				var ctor pkg_rhino.Node = ctorDef.GetFirstChild()
				_ = ctor
				CheckInterfaces_CheckInterfaceConstructorArgs(t, ctor)
				CheckInterfaces_CheckConstructorBlock(t, ctor)
			}
			CheckInterfaces_CheckClassMethods(t, n, ctorDef)
		}
	default:
		return
	}
}

// Java static method: CheckInterfaces.checkInterfaceConstructorArgs(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node funcNode)
func CheckInterfaces_CheckInterfaceConstructorArgs(t *NodeTraversal, funcNode pkg_rhino.Node) {
	clinit_CheckInterfaces()
	_ = t
	_ = funcNode
	var args pkg_rhino.Node = funcNode.GetSecondChild()
	_ = args
	if args.HasChildren() {
		t.Report_3(args.GetFirstChild(), StaticGet_CheckInterfaces_INTERFACE_CONSTRUCTOR_SHOULD_NOT_TAKE_ARGS(), []string{})
	}
}

// Java static method: CheckInterfaces.checkClassMethods(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node classNode, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node ctorDef)
func CheckInterfaces_CheckClassMethods(t *NodeTraversal, classNode pkg_rhino.Node, ctorDef pkg_rhino.Node) {
	clinit_CheckInterfaces()
	_ = t
	_ = classNode
	_ = ctorDef
	var classMembers pkg_rhino.Node = classNode.GetLastChild()
	_ = classMembers
	if !(classMembers.IsClassMembers()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(classMembers))
	}
	{
		var memberFuncDef pkg_rhino.Node = classMembers.GetFirstChild()
		_ = memberFuncDef
		for ; memberFuncDef != 0; memberFuncDef = memberFuncDef.GetNext() {
			if j2gort.CheckNotNullNode(memberFuncDef) == ctorDef {
				continue
			}
			if memberFuncDef.IsStaticMember() {
				var className string = NodeUtil_GetName(classNode)
				_ = className
				if className == j2gort.NullString {
					className = "<anonymous>"
				}
				var funcName string = memberFuncDef.GetString()
				_ = funcName
				t.Report_3(memberFuncDef, StaticGet_CheckInterfaces_STATIC_MEMBER_FUNCTION_IN_INTERFACE_CLASS(), []string{className, funcName})
			} else {
				var block pkg_rhino.Node = memberFuncDef.GetLastChild().GetLastChild()
				_ = block
				if block.HasChildren() {
					t.Report_3(block.GetFirstChild(), StaticGet_CheckInterfaces_INTERFACE_CLASS_NONSTATIC_METHOD_NOT_EMPTY(), []string{})
				}
			}
		}
	}
}

// Java static method: CheckInterfaces.checkConstructorBlock(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node funcNode)
func CheckInterfaces_CheckConstructorBlock(t *NodeTraversal, funcNode pkg_rhino.Node) {
	clinit_CheckInterfaces()
	_ = t
	_ = funcNode
	var block pkg_rhino.Node = funcNode.GetLastChild()
	_ = block
	if !(block.HasChildren()) {
		return
	}
	{
		var stmt pkg_rhino.Node = block.GetFirstChild()
		_ = stmt
		for ; stmt != 0; stmt = stmt.GetNext() {
			if !(CheckInterfaces_IsThisPropAccess(stmt)) {
				t.Report_3(stmt, StaticGet_CheckInterfaces_NON_DECLARATION_STATEMENT_IN_INTERFACE(), []string{})
				break
			} else if stmt.GetFirstChild().GetJSDocInfo() == nil {
				t.Report_3(stmt, StaticGet_CheckInterfaces_MISSING_JSDOC_IN_DECLARATION_STATEMENT(), []string{})
				break
			}
		}
	}
}

// Java static method: CheckInterfaces.isThisPropAccess(com.google.javascript.rhino.Node stmt)
func CheckInterfaces_IsThisPropAccess(stmt pkg_rhino.Node) bool {
	clinit_CheckInterfaces()
	_ = stmt
	return ((stmt.IsExprResult() && stmt.GetFirstChild().IsGetProp()) && stmt.GetFirstFirstChild().IsThis())
}
func (*CheckInterfaces) AlwaysShouldTraverse_()      {}
func (*CheckInterfaces) CanReleaseSyntacticScopes_() {}
