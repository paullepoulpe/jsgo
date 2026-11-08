// Port of com/google/javascript/jscomp/lint/CheckInterfaces.java @ closure-compiler v20261104

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
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

var checkInterfacesVisitTokens = makeVisitTokenSet(rhino.Token_FUNCTION, rhino.Token_CLASS)

var (
	_ nodeTraversal_visitTokens   = (*CheckInterfaces)(nil)
	_ nodeTraversal_noVisitLeaves = (*CheckInterfaces)(nil)
)

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

// Java: CheckInterfaces.visit(NodeTraversal, Node, Node)
func (c *CheckInterfaces) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if fn, ok := ast.AsFunction(n); ok {
		if jsdoc := NodeUtil_GetBestJSDocInfo(n); jsdoc != nil && jsdoc.IsInterface() {
			checkInterfaceConstructorArgs(t, fn)
			checkConstructorBlock(t, fn)
		}
	} else if cls, ok := ast.AsClass(n); ok {
		if jsdoc := NodeUtil_GetBestJSDocInfo(n); jsdoc != nil && jsdoc.IsInterface() {
			if ext, hasExt := cls.Extends(); hasExt && ext.Token() == rhino.Token_NAME {
				t.Report(n, CheckInterfaces_INTERFACE_DEFINED_WITH_EXTENDS)
			}
			ctorDef, hasCtor := cls.Constructor()
			if hasCtor {
				ctor := ctorDef.Function()
				checkInterfaceConstructorArgs(t, ctor)
				checkConstructorBlock(t, ctor)
			}
			checkClassMethods(t, cls, ctorDef.Node())
		}
	}
}

// Java: CheckInterfaces.checkInterfaceConstructorArgs(NodeTraversal, Node)
func checkInterfaceConstructorArgs(t *NodeTraversal, fn ast.Function) {
	if firstParam, ok := fn.Params().FirstParam(); ok {
		t.Report(firstParam.Node(), CheckInterfaces_INTERFACE_CONSTRUCTOR_SHOULD_NOT_TAKE_ARGS)
	}
}

// checkClassMethods: non-static class methods must be empty for `@record` and
// `@interface` as per the style guide.
//
// Java: CheckInterfaces.checkClassMethods(NodeTraversal, Node, Node)
func checkClassMethods(t *NodeTraversal, cls ast.Class, ctorDef rhino.Node) {
	for member := range cls.Body().Members() {
		memberNode := member.Node()
		if memberNode == ctorDef {
			continue // constructor was already checked; don't check here.
		}
		if memberNode.IsStaticMember() {
			name := checkInterfaces_ANONYMOUS_CLASSNAME
			if className, ok := NodeUtil_GetName(cls.Node()); ok {
				name = className
			}
			t.Report(memberNode, CheckInterfaces_STATIC_MEMBER_FUNCTION_IN_INTERFACE_CLASS, name, memberNode.GetString())
		} else if fn, ok := ast.AsFunction(memberNode.GetLastChild()); ok {
			if first := fn.Body().Node().GetFirstChild(); !first.IsNil() {
				t.Report(first, CheckInterfaces_INTERFACE_CLASS_NONSTATIC_METHOD_NOT_EMPTY)
			}
		}
	}
}

// checkConstructorBlock: `@record` constructors can be non-empty, check that
// only field declarations exist in them.
//
// Java: CheckInterfaces.checkConstructorBlock(NodeTraversal, Node)
func checkConstructorBlock(t *NodeTraversal, fn ast.Function) {
	block, ok := fn.Body().AsBlock()
	if !ok {
		return
	}
	for stmt := range block.Statements() {
		stmtNode := stmt.Node()
		exprStmt, ok := stmt.AsExprStmt()
		if !ok {
			t.Report(stmtNode, CheckInterfaces_NON_DECLARATION_STATEMENT_IN_INTERFACE)
			break
		}
		gp, ok := exprStmt.Expr().AsGetProp()
		if !ok || gp.Object().Token() != rhino.Token_THIS {
			t.Report(stmtNode, CheckInterfaces_NON_DECLARATION_STATEMENT_IN_INTERFACE)
			break
		}
		if gp.Node().GetJSDocInfo() == nil {
			t.Report(stmtNode, CheckInterfaces_MISSING_JSDOC_IN_DECLARATION_STATEMENT)
			break
		}
	}
}
