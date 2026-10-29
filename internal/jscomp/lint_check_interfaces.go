// Port of com/google/javascript/jscomp/lint/CheckInterfaces.java @ closure-compiler v20261021

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

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

var checkInterfacesVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_FUNCTION)
	s.Add(rhino.Token_CLASS)
	return s
}()

var _ nodeTraversal_visitTokens = (*CheckInterfaces)(nil)

func (*CheckInterfaces) noVisitLeaves()             {}
func (*CheckInterfaces) visitTokens() *ast.TokenSet { return &checkInterfacesVisitTokens }

// CheckInterfaces checks for errors related to interfaces.
type CheckInterfaces struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler
}

// Placeholder class name for error reporting on anonymous classes.
const checkInterfaces_ANONYMOUS_CLASSNAME = "<anonymous>"

var CheckInterfaces_NON_DECLARATION_STATEMENT_IN_INTERFACE = DiagnosticType_Disabled(
	"JSC_NON_DECLARATION_STATEMENT_IN_INTERFACE",
	"@interface or @record functions should not contain statements other than field"+
		" declarations")

var CheckInterfaces_MISSING_JSDOC_IN_DECLARATION_STATEMENT = DiagnosticType_Disabled(
	"JSC_MISSING_JSDOC_IN_DECLARATION_STATEMENT",
	"@interface or @record functions must contain JSDoc for each field declaration.")

var CheckInterfaces_INTERFACE_CLASS_NONSTATIC_METHOD_NOT_EMPTY = DiagnosticType_Disabled(
	"JSC_INTERFACE_CLASS_NONSTATIC_METHOD_NOT_EMPTY",
	"interface methods must have an empty body")

var CheckInterfaces_INTERFACE_CONSTRUCTOR_SHOULD_NOT_TAKE_ARGS = DiagnosticType_Disabled(
	"JSC_INTERFACE_CONSTRUCTOR_SHOULD_NOT_TAKE_ARGS",
	"Interface constructors should not take any arguments")

var CheckInterfaces_STATIC_MEMBER_FUNCTION_IN_INTERFACE_CLASS = DiagnosticType_Disabled(
	"JSC_STATIC_MEMBER_FUNCTION_IN_INTERFACE_CLASS",
	"Interface class should not have static member functions. "+
		"Consider pulling out the static method into a flat name as {0}_{1}")

var CheckInterfaces_INTERFACE_DEFINED_WITH_EXTENDS = DiagnosticType_Disabled(
	"JSC_INTERFACE_DEFINED_WITH_EXTENDS",
	"Interface/Record class should use the `@extends` annotation instead of extends"+
		" keyword.")

// NewCheckInterfaces is new CheckInterfaces(compiler).
//
// Java: CheckInterfaces(AbstractCompiler)
func NewCheckInterfaces(compiler AbstractCompiler) *CheckInterfaces {
	return &CheckInterfaces{compiler: compiler}
}

func (*CheckInterfaces) cloneForParallel(c *Compiler) NodeTraversal_Callback {
	return NewCheckInterfaces(c)
}

// Java: CheckInterfaces.process(Node, Node)
func (c *CheckInterfaces) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}

// CheckInterfaces_IsInterface: whether jsDoc is present and has an
// @interface or @record annotation
//
// Java: CheckInterfaces.isInterface(JSDocInfo)
func CheckInterfaces_IsInterface(jsDoc *rhino.JSDocInfo) bool {
	return jsDoc != nil && jsDoc.IsInterface()
}

// Java: CheckInterfaces.visit(NodeTraversal, Node, Node)
func (c *CheckInterfaces) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_FUNCTION:
		jsdoc := NodeUtil_GetBestJSDocInfo(n)
		if CheckInterfaces_IsInterface(jsdoc) {
			CheckInterfaces_CheckInterfaceConstructorArgs(t, n)
			CheckInterfaces_CheckConstructorBlock(t, n)
		}
	case rhino.Token_CLASS:
		jsdoc := NodeUtil_GetBestJSDocInfo(n)
		if CheckInterfaces_IsInterface(jsdoc) {
			if !n.GetSecondChild().IsNil() && n.GetSecondChild().IsName() {
				t.Report(n, CheckInterfaces_INTERFACE_DEFINED_WITH_EXTENDS)
			}
			ctorDef := NodeUtil_GetEs6ClassConstructorMemberFunctionDef(n)
			if !ctorDef.IsNil() {
				ctor := ctorDef.GetFirstChild()
				CheckInterfaces_CheckInterfaceConstructorArgs(t, ctor)
				CheckInterfaces_CheckConstructorBlock(t, ctor)
			}
			CheckInterfaces_CheckClassMethods(t, n, ctorDef)
		}
	default:
		return
	}
}

// Java: CheckInterfaces.checkInterfaceConstructorArgs(NodeTraversal, Node)
func CheckInterfaces_CheckInterfaceConstructorArgs(t *NodeTraversal, funcNode rhino.Node) {
	args := funcNode.GetSecondChild()
	if args.HasChildren() {
		t.Report(args.GetFirstChild(), CheckInterfaces_INTERFACE_CONSTRUCTOR_SHOULD_NOT_TAKE_ARGS)
	}
}

// CheckInterfaces_CheckClassMethods: non-static class methods must be empty
// for `@record` and `@interface` as per the style guide.
//
// Java: CheckInterfaces.checkClassMethods(NodeTraversal, Node, Node)
func CheckInterfaces_CheckClassMethods(t *NodeTraversal, classNode rhino.Node, ctorDef rhino.Node) {
	classMembers := classNode.GetLastChild()
	jutil.CheckStateMsg1(classMembers.IsClassMembers(), "%s", classMembers)
	for memberFuncDef := classMembers.GetFirstChild(); !memberFuncDef.IsNil(); memberFuncDef = memberFuncDef.GetNext() {
		// Node does not override equals: memberFuncDef.equals(ctorDef) is identity.
		if memberFuncDef == ctorDef {
			continue // constructor was already checked; don't check here.
		}
		if memberFuncDef.IsStaticMember() {
			// `static foo() {...}`
			name := checkInterfaces_ANONYMOUS_CLASSNAME
			if className := NodeUtil_GetName(classNode); className != nil {
				name = *className
			}
			funcName := memberFuncDef.GetString()
			t.Report(memberFuncDef, CheckInterfaces_STATIC_MEMBER_FUNCTION_IN_INTERFACE_CLASS, name, funcName)
		} else {
			block := memberFuncDef.GetLastChild().GetLastChild()
			if block.HasChildren() {
				t.Report(block.GetFirstChild(), CheckInterfaces_INTERFACE_CLASS_NONSTATIC_METHOD_NOT_EMPTY)
			}
		}
	}
}

// CheckInterfaces_CheckConstructorBlock: `@record` constructors can be
// non-empty, check that only field declarations exist in them.
//
// Java: CheckInterfaces.checkConstructorBlock(NodeTraversal, Node)
func CheckInterfaces_CheckConstructorBlock(t *NodeTraversal, funcNode rhino.Node) {
	block := funcNode.GetLastChild()
	if !block.HasChildren() {
		return
	}

	for stmt := block.GetFirstChild(); !stmt.IsNil(); stmt = stmt.GetNext() {
		if !CheckInterfaces_IsThisPropAccess(stmt) {
			// Only field declarations are expected inside @record and @interface.
			t.Report(stmt, CheckInterfaces_NON_DECLARATION_STATEMENT_IN_INTERFACE)
			break
		} else if stmt.GetFirstChild().GetJSDocInfo() == nil {
			// A field declaration that's missing a JSDoc.
			t.Report(stmt, CheckInterfaces_MISSING_JSDOC_IN_DECLARATION_STATEMENT)
			break
		}
	}
}

// Java: CheckInterfaces.isThisPropAccess(Node)
func CheckInterfaces_IsThisPropAccess(stmt rhino.Node) bool {
	return stmt.IsExprResult() &&
		stmt.GetFirstChild().IsGetProp() &&
		stmt.GetFirstFirstChild().IsThis()
}
