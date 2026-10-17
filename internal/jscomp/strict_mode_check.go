// Port of com/google/javascript/jscomp/StrictModeCheck.java @ closure-compiler v20261015

/*
 * Copyright 2009 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// StrictModeCheck checks that the code obeys the static restrictions of
// strict mode:
//
//  1. No use of "with".
//  2. No deleting variables, functions, or arguments.
//  3. No re-declarations or assignments of "eval" or arguments.
//  4. No use of arguments.callee
//  5. No use of arguments.caller
//  6. Class: Always under strict mode
//  7. In addition, no duplicate class method names
type StrictModeCheck struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	compiler     AbstractCompiler
	defaultLevel CheckLevel
}

// Java: StrictModeCheck.USE_OF_WITH
var StrictModeCheck_USE_OF_WITH = DiagnosticType_Error(
	"JSC_USE_OF_WITH", "The 'with' statement cannot be used in strict mode.")

// Java: StrictModeCheck.EVAL_DECLARATION
var StrictModeCheck_EVAL_DECLARATION = DiagnosticType_Error("JSC_EVAL_DECLARATION", "\"eval\" cannot be redeclared in strict mode")

// Java: StrictModeCheck.EVAL_ASSIGNMENT
var StrictModeCheck_EVAL_ASSIGNMENT = DiagnosticType_Error(
	"JSC_EVAL_ASSIGNMENT", "the \"eval\" object cannot be reassigned in strict mode")

// Java: StrictModeCheck.ARGUMENTS_DECLARATION
var StrictModeCheck_ARGUMENTS_DECLARATION = DiagnosticType_Error(
	"JSC_ARGUMENTS_DECLARATION", "\"arguments\" cannot be redeclared in strict mode")

// Java: StrictModeCheck.ARGUMENTS_ASSIGNMENT
var StrictModeCheck_ARGUMENTS_ASSIGNMENT = DiagnosticType_Error(
	"JSC_ARGUMENTS_ASSIGNMENT",
	"the \"arguments\" object cannot be reassigned in strict mode")

// Java: StrictModeCheck.ARGUMENTS_CALLEE_FORBIDDEN
var StrictModeCheck_ARGUMENTS_CALLEE_FORBIDDEN = DiagnosticType_Error(
	"JSC_ARGUMENTS_CALLEE_FORBIDDEN", "\"arguments.callee\" cannot be used in strict mode")

// Java: StrictModeCheck.ARGUMENTS_CALLER_FORBIDDEN
var StrictModeCheck_ARGUMENTS_CALLER_FORBIDDEN = DiagnosticType_Error(
	"JSC_ARGUMENTS_CALLER_FORBIDDEN", "\"arguments.caller\" cannot be used in strict mode")

// Java: StrictModeCheck.FUNCTION_CALLER_FORBIDDEN
var StrictModeCheck_FUNCTION_CALLER_FORBIDDEN = DiagnosticType_Error(
	"JSC_FUNCTION_CALLER_FORBIDDEN",
	"A function''s \"caller\" property cannot be used in strict mode")

// Java: StrictModeCheck.FUNCTION_ARGUMENTS_PROP_FORBIDDEN
var StrictModeCheck_FUNCTION_ARGUMENTS_PROP_FORBIDDEN = DiagnosticType_Error(
	"JSC_FUNCTION_ARGUMENTS_PROP_FORBIDDEN",
	"A function''s \"arguments\" property cannot be used in strict mode")

// Java: StrictModeCheck.DELETE_VARIABLE
var StrictModeCheck_DELETE_VARIABLE = DiagnosticType_Error(
	"JSC_DELETE_VARIABLE",
	"variables, functions, and arguments cannot be deleted in strict mode")

// Java: StrictModeCheck.DUPLICATE_MEMBER
var StrictModeCheck_DUPLICATE_MEMBER = DiagnosticType_Warning(
	"JSC_DUPLICATE_MEMBER",
	"Class or object literal contains duplicate member \"{0}\". In non-strict code, the last"+
		" duplicate will overwrite the others.")

// Java: StrictModeCheck(AbstractCompiler, CheckLevel)
func NewStrictModeCheck(compiler AbstractCompiler, defaultLevel CheckLevel) *StrictModeCheck {
	c := &StrictModeCheck{}
	c.compiler = compiler
	c.defaultLevel = defaultLevel
	return c
}

// Java: StrictModeCheck.process(Node, Node)
func (c *StrictModeCheck) Process(externs *rhino.Node, root *rhino.Node) {
	NodeTraversal_TraverseRoots(c.compiler, c, externs, root)
	NodeTraversal_Traverse(c.compiler, root, NewStrictModeCheck_NonExternChecks(c))
}

// Java: StrictModeCheck.visit(NodeTraversal, Node, Node)
func (c *StrictModeCheck) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if n.IsAssign() {
		c.CheckAssignment(n)
	} else if n.IsDelProp() {
		c.CheckDelete(t, n)
	} else if n.IsObjectLit() {
		c.CheckObjectLiteralOrClass(n)
	} else if n.IsClass() {
		c.CheckObjectLiteralOrClass(n.GetLastChild())
	} else if n.IsWith() {
		c.CheckWith(n)
	}
}

// CheckWith reports a warning for with statements.
//
// Java: StrictModeCheck.checkWith(Node)
func (c *StrictModeCheck) CheckWith(n *rhino.Node) {
	info := n.GetJSDocInfo()
	allowWith := info != nil && info.GetSuppressions().Contains("with")
	if !allowWith {
		c.Report(n, StrictModeCheck_USE_OF_WITH)
	}
}

// StrictModeCheck_IsDeclaration determines if the given name is a
// declaration, which can be a declaration of a variable, function, or
// argument.
//
// Java: StrictModeCheck.isDeclaration(Node)
func StrictModeCheck_IsDeclaration(n *rhino.Node) bool {
	switch n.GetParent().GetToken() {
	case rhino.Token_LET, rhino.Token_CONST, rhino.Token_VAR, rhino.Token_CATCH:
		return true
	case rhino.Token_FUNCTION:
		return n == n.GetParent().GetFirstChild()
	case rhino.Token_PARAM_LIST:
		return n.GetGrandparent().IsFunction()
	default:
		return false
	}
}

// CheckAssignment checks that an assignment is not to the "arguments" object.
//
// Java: StrictModeCheck.checkAssignment(Node)
func (c *StrictModeCheck) CheckAssignment(n *rhino.Node) {
	if n.GetFirstChild().IsName() {
		if "arguments" == n.GetFirstChild().GetString() {
			c.Report(n, StrictModeCheck_ARGUMENTS_ASSIGNMENT)
		} else if "eval" == n.GetFirstChild().GetString() {
			// Note that assignment to eval is already illegal because any use of
			// that name is illegal.
			c.Report(n, StrictModeCheck_EVAL_ASSIGNMENT)
		}
	}
}

// CheckDelete checks that variables, functions, and arguments are not
// deleted.
//
// Java: StrictModeCheck.checkDelete(NodeTraversal, Node)
func (c *StrictModeCheck) CheckDelete(t *NodeTraversal, n *rhino.Node) {
	if n.GetFirstChild().IsName() {
		v := t.GetScope().GetVar(n.GetFirstChild().GetString())
		if v != nil {
			c.Report(n, StrictModeCheck_DELETE_VARIABLE)
		}
	}
}

// CheckObjectLiteralOrClass checks that object literal keys or class method
// names are valid.
//
// Java: StrictModeCheck.checkObjectLiteralOrClass(Node)
func (c *StrictModeCheck) CheckObjectLiteralOrClass(n *rhino.Node) {
	getters := jutil.NewLinkedHashSet[string]()
	setters := jutil.NewLinkedHashSet[string]()
	staticGetters := jutil.NewLinkedHashSet[string]()
	staticSetters := jutil.NewLinkedHashSet[string]()

	/*
	 * Iterate backwards because the last duplicate is the one that will be used in sloppy or ES6
	 * code. The earlier duplicates are the ones that should be removed.
	 */
	for key := n.GetLastChild(); key != nil; key = key.GetPrevious() {
		if key.IsEmpty() ||
			key.IsComputedProp() ||
			key.IsSpread() ||
			key.IsComputedFieldDef() ||
			// Computed properties cannot be computed at compile time
			NodeUtil_IsClassStaticBlock(key) {
			// Will not check since whether duplicates are declared/assigned cannot be determined at
			// compile time since we do not know which code will be run
			continue
		}

		keyName := key.GetString()
		if !key.IsSetterDef() {
			// normal property and getter cases
			set := getters
			if key.IsStaticMember() {
				set = staticGetters
			}
			if !set.Add(keyName) {
				c.Report(key, StrictModeCheck_DUPLICATE_MEMBER, keyName)
			}
		}
		if !key.IsGetterDef() {
			// normal property and setter cases
			set := setters
			if key.IsStaticMember() {
				set = staticSetters
			}
			if !set.Add(keyName) {
				c.Report(key, StrictModeCheck_DUPLICATE_MEMBER, keyName)
			}
		}
	}
}

// StrictModeCheck_NonExternChecks holds checks that are performed on
// non-extern code only.
//
// Java: StrictModeCheck.NonExternChecks
type StrictModeCheck_NonExternChecks struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	// Go-only: language-adapter — the enclosing StrictModeCheck instance of
	// this inner class.
	outer *StrictModeCheck
}

// Java: StrictModeCheck.NonExternChecks()
func NewStrictModeCheck_NonExternChecks(outer *StrictModeCheck) *StrictModeCheck_NonExternChecks {
	return &StrictModeCheck_NonExternChecks{outer: outer}
}

// Java: StrictModeCheck.NonExternChecks.visit(NodeTraversal, Node, Node)
func (c *StrictModeCheck_NonExternChecks) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if n.IsName() && StrictModeCheck_IsDeclaration(n) {
		c.CheckDeclaration(n)
	} else if n.IsGetProp() {
		c.CheckGetProp(n)
	}
}

// CheckDeclaration checks for illegal declarations.
//
// Java: StrictModeCheck.NonExternChecks.checkDeclaration(Node)
func (c *StrictModeCheck_NonExternChecks) CheckDeclaration(n *rhino.Node) {
	if "eval" == n.GetString() {
		c.outer.Report(n, StrictModeCheck_EVAL_DECLARATION)
	} else if "arguments" == n.GetString() {
		c.outer.Report(n, StrictModeCheck_ARGUMENTS_DECLARATION)
	}
}

// CheckGetProp checks that the arguments.callee is not used.
//
// Java: StrictModeCheck.NonExternChecks.checkGetProp(Node)
func (c *StrictModeCheck_NonExternChecks) CheckGetProp(n *rhino.Node) {
	target := n.GetFirstChild()
	name := n.GetString()
	if name == "callee" {
		if target.IsName() && target.GetString() == "arguments" {
			c.outer.Report(n, StrictModeCheck_ARGUMENTS_CALLEE_FORBIDDEN)
		}
	} else if name == "caller" {
		if target.IsName() && target.GetString() == "arguments" {
			c.outer.Report(n, StrictModeCheck_ARGUMENTS_CALLER_FORBIDDEN)
		} else if StrictModeCheck_IsFunctionType(target) {
			c.outer.Report(n, StrictModeCheck_FUNCTION_CALLER_FORBIDDEN)
		}
	} else if name == "arguments" && StrictModeCheck_IsFunctionType(target) {
		c.outer.Report(n, StrictModeCheck_FUNCTION_ARGUMENTS_PROP_FORBIDDEN)
	}
}

// Java: StrictModeCheck.isFunctionType(Node)
func StrictModeCheck_IsFunctionType(n *rhino.Node) bool {
	type_ := n.GetJSType()
	return type_ != nil && type_.IsFunctionType()
}

// Java: StrictModeCheck.report(Node, DiagnosticType, String...)
func (c *StrictModeCheck) Report(n *rhino.Node, diagnostic *DiagnosticType, args ...string) {
	c.compiler.Report(
		JSError_Builder_(diagnostic, args...).SetLevel(c.defaultLevel).SetNode(n).Build())
}
