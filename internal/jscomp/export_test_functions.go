// Port of com/google/javascript/jscomp/ExportTestFunctions.java @ closure-compiler v20261015

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
	"regexp"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ExportTestFunctions generates goog.exportSymbol for test functions, so they
// can be recognized by the test runner, even if the code is compiled.
type ExportTestFunctions struct {
	compiler               AbstractCompiler
	exportSymbolFunction   *string
	exportPropertyFunction *string
}

// Go-only: language-adapter — the Java pattern is anchored with ^...$ and
// only used with Matcher.matches(), where Java's `$` can only match at the
// end of the input; Go's regexp gives the same answer (\w is ASCII in both),
// and jutil.Pattern does not accept `$`.
var exportTestFunctions_TEST_FUNCTIONS_NAME_PATTERN = regexp.MustCompile(
	`^(?:((\w+\.)+prototype\.||window\.)*` +
		`(setUpPage|setUp|shouldRunTests|tearDown|tearDownPage|test[\w\$]+))$`)

const exportTestFunctions_GOOG_TESTING_TEST_SUITE = "goog.testing.testSuite"

// NewExportTestFunctions creates a new export test functions compiler pass.
//
// exportSymbolFunction is the function name used to export symbols in JS;
// exportPropertyFunction is the function name used to export properties in
// JS.
//
// Java: ExportTestFunctions(AbstractCompiler, String, String)
func NewExportTestFunctions(
	compiler AbstractCompiler, exportSymbolFunction *string, exportPropertyFunction *string) *ExportTestFunctions {
	e := &ExportTestFunctions{}

	jutil.CheckNotNull(compiler)
	e.compiler = compiler
	e.exportSymbolFunction = exportSymbolFunction
	e.exportPropertyFunction = exportPropertyFunction
	return e
}

// ExportTestFunctions_ExportTestFunctionsNodes is the traversal callback.
type ExportTestFunctions_ExportTestFunctionsNodes struct {
	NodeTraversal_AbstractShallowCallbackBase

	outer *ExportTestFunctions
}

// Java: ExportTestFunctions.ExportTestFunctionsNodes()
func newExportTestFunctions_ExportTestFunctionsNodes(outer *ExportTestFunctions) *ExportTestFunctions_ExportTestFunctionsNodes {
	return &ExportTestFunctions_ExportTestFunctionsNodes{outer: outer}
}

// Java: ExportTestFunctions.ExportTestFunctionsNodes.visit(NodeTraversal, Node, Node)
func (x *ExportTestFunctions_ExportTestFunctionsNodes) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if parent.IsNil() {
		return
	}

	if parent.IsScript() || parent.IsModuleBody() {
		if NodeUtil_IsFunctionDeclaration(n) {
			// Check for a test function statement.
			functionName := NodeUtil_GetName(n)
			if ExportTestFunctions_IsTestFunction(functionName) {
				x.outer.ExportTestFunctionAsSymbol(*functionName, n)
			}
		} else if x.IsNameDeclaredFunction(n) {
			// Check for a test function expression.
			functionNode := n.GetFirstFirstChild()
			functionName := NodeUtil_GetName(functionNode)
			if ExportTestFunctions_IsTestFunction(functionName) {
				x.outer.ExportTestFunctionAsSymbol(*functionName, n)
			}
		} else if x.IsNameDeclaredClass(n) {
			classNode := n.GetFirstFirstChild()
			className := NodeUtil_GetName(classNode)
			x.ExportClassWithName(classNode, className, n)
		} else if n.IsClass() {
			x.ExportClass(n)
		}
	} else if NodeUtil_IsExprAssign(parent) {
		// Check for a test method assignment.
		grandparent := parent.GetParent()
		if !grandparent.IsNil() && (grandparent.IsScript() || grandparent.IsModuleBody()) {
			//                                    NAME/(GETPROP -> ... -> NAME)
			// SCRIPT -> EXPR_RESULT -> ASSIGN ->
			//                                    FUNCTION/CLASS
			firstChild := n.GetFirstChild()
			lastChild := n.GetLastChild()
			nodeName := firstChild.GetQualifiedName()

			if lastChild.IsFunction() {
				if ExportTestFunctions_IsTestFunction(nodeName) {
					if n.GetFirstChild().IsName() {
						x.outer.ExportTestFunctionAsSymbol(*nodeName, parent)
					} else {
						x.outer.ExportTestFunctionAsProperty(firstChild, n)
					}
				}
			} else if lastChild.IsClass() {
				x.ExportClassWithName(lastChild, nodeName, parent)
			}
		}
	} else if ExportTestFunctions_IsTestSuiteArgument(n, t) {
		for c := n.GetFirstChild(); !c.IsNil(); {
			next := c.GetNext()
			if c.IsStringKey() && !c.IsQuotedStringKey() {
				c.SetQuotedStringKey()
				x.outer.compiler.ReportChangeToEnclosingScope(c)
			} else if c.IsMemberFunctionDef() {
				x.RewriteMemberDefInObjLit(c, n)
			}
			c = next
		}
	}
}

// Java: ExportTestFunctions.ExportTestFunctionsNodes.exportClass(Node)
func (x *ExportTestFunctions_ExportTestFunctionsNodes) ExportClass(classNode rhino.Node) {
	className := NodeUtil_GetName(classNode)
	x.ExportClassWithName(classNode, className, classNode)
}

// ExportClassWithName: className is a nullable Java String (Java string
// concatenation renders null as "null").
//
// Java: ExportTestFunctions.ExportTestFunctionsNodes.exportClass(Node, String, Node)
func (x *ExportTestFunctions_ExportTestFunctionsNodes) ExportClassWithName(classNode rhino.Node, className *string, addAfter rhino.Node) {
	classMembers := classNode.GetLastChild()
	for maybeMemberFunctionDef := classMembers.GetFirstChild(); !maybeMemberFunctionDef.IsNil(); maybeMemberFunctionDef = maybeMemberFunctionDef.GetNext() {
		if maybeMemberFunctionDef.IsMemberFunctionDef() {
			methodName := maybeMemberFunctionDef.GetString()
			if ExportTestFunctions_IsTestFunction(&methodName) {
				functionRef := jutil.ToJavaString(className) + ".prototype." + methodName
				classRef := jutil.ToJavaString(className) + ".prototype"

				exportCallTarget :=
					NodeUtil_NewQNameWithOriginalName(
						x.outer.compiler, *x.outer.exportPropertyFunction, maybeMemberFunctionDef, methodName)
				call := rhino.IR_Call(exportCallTarget)
				if exportCallTarget.IsName() {
					call.PutBooleanProp(rhino.Node_FREE_CALL, true)
				}

				call.AddChildToBack(
					NodeUtil_NewQNameWithOriginalName(x.outer.compiler, classRef, maybeMemberFunctionDef, classRef))
				call.AddChildToBack(rhino.IR_String(methodName))
				call.AddChildToBack(
					NodeUtil_NewQNameWithOriginalName(x.outer.compiler, functionRef, maybeMemberFunctionDef, functionRef))

				expression := rhino.IR_ExprResult(call).SrcrefTreeIfMissing(maybeMemberFunctionDef)

				expression.InsertAfter(addAfter)
				x.outer.compiler.ReportChangeToEnclosingScope(expression)
				addAfter = expression
			}
		}
	}
}

// RewriteMemberDefInObjLit converts a member function into a quoted string
// key to avoid property renaming
//
// Java: ExportTestFunctions.ExportTestFunctionsNodes.rewriteMemberDefInObjLit(Node, Node)
func (x *ExportTestFunctions_ExportTestFunctionsNodes) RewriteMemberDefInObjLit(memberDef rhino.Node, objLit rhino.Node) {
	name := memberDef.GetString()
	stringKey := rhino.IR_StringKeyWithValue(name, memberDef.RemoveFirstChild())
	memberDef.ReplaceWith(stringKey)
	stringKey.SetQuotedStringKey()
	stringKey.SetJSDocInfo(memberDef.GetJSDocInfo())
	x.outer.compiler.ReportChangeToEnclosingScope(objLit)
}

// GetNameDeclaredGrandchild gets the node that corresponds to an expression
// declared with var, let or const. This has the AST structure VAR/LET/CONST
// -> NAME -> NODE
//
// Java: ExportTestFunctions.ExportTestFunctionsNodes.getNameDeclaredGrandchild(Node)
func (x *ExportTestFunctions_ExportTestFunctionsNodes) GetNameDeclaredGrandchild(node rhino.Node) rhino.Node {
	if !NodeUtil_IsNameDeclaration(node) {
		return rhino.Node{}
	}
	return node.GetFirstFirstChild()
}

// IsNameDeclaredFunction is whether node corresponds to a function
// expression declared with var, let or const which is of the form:
//
//	var/let/const functionName = function() {
//	  // Implementation
//	};
//
// This has the AST structure VAR/LET/CONST -> NAME -> FUNCTION
//
// Java: ExportTestFunctions.ExportTestFunctionsNodes.isNameDeclaredFunction(Node)
func (x *ExportTestFunctions_ExportTestFunctionsNodes) IsNameDeclaredFunction(node rhino.Node) bool {
	grandchild := x.GetNameDeclaredGrandchild(node)
	return !grandchild.IsNil() && grandchild.IsFunction()
}

// IsNameDeclaredClass is whether node corresponds to a class declared with
// var, let or const which is of the form:
//
//	var/let/const className = class {
//	  // Implementation
//	};
//
// This has the AST structure VAR/LET/CONST -> NAME -> CLASS
//
// Java: ExportTestFunctions.ExportTestFunctionsNodes.isNameDeclaredClass(Node)
func (x *ExportTestFunctions_ExportTestFunctionsNodes) IsNameDeclaredClass(node rhino.Node) bool {
	grandchild := x.GetNameDeclaredGrandchild(node)
	return !grandchild.IsNil() && grandchild.IsClass()
}

// Java: ExportTestFunctions.process(Node, Node)
func (e *ExportTestFunctions) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(e.compiler, root, newExportTestFunctions_ExportTestFunctionsNodes(e))
}

// ExportTestFunctionAsSymbol adds exportSymbol(testFunctionName,
// testFunction);
//
// Java: ExportTestFunctions.exportTestFunctionAsSymbol(String, Node)
func (e *ExportTestFunctions) ExportTestFunctionAsSymbol(testFunctionName string, node rhino.Node) {

	exportCallTarget := NodeUtil_NewQNameWithOriginalName(e.compiler,
		*e.exportSymbolFunction, node, testFunctionName)
	call := rhino.IR_Call(exportCallTarget)
	if exportCallTarget.IsName() {
		call.PutBooleanProp(rhino.Node_FREE_CALL, true)
	}
	call.AddChildToBack(rhino.IR_String(testFunctionName))
	call.AddChildToBack(NodeUtil_NewQNameWithOriginalName(e.compiler,
		testFunctionName, node, testFunctionName))

	expression := rhino.IR_ExprResult(call).SrcrefTreeIfMissing(node)

	expression.InsertAfter(node)
	e.compiler.ReportChangeToEnclosingScope(expression)
}

// ExportTestFunctionAsProperty adds exportProperty() of the test function
// name on the prototype object
//
// Java: ExportTestFunctions.exportTestFunctionAsProperty(Node, Node)
func (e *ExportTestFunctions) ExportTestFunctionAsProperty(fullyQualifiedFunctionName rhino.Node, node rhino.Node) {
	jutil.CheckStateMsg1(fullyQualifiedFunctionName.IsGetProp(), "%s", fullyQualifiedFunctionName)

	testFunctionName :=
		NodeUtil_GetPrototypePropertyName(node.GetFirstChild())
	if jutil.StringStartsWith(*node.GetFirstChild().GetQualifiedName(), "window.") {
		testFunctionName = jutil.StringSubstringFrom(*node.GetFirstChild().GetQualifiedName(), jutil.StringLength("window."))
	}

	exportCall :=
		rhino.IR_Call(
			NodeUtil_NewQName(e.compiler, *e.exportPropertyFunction),
			fullyQualifiedFunctionName.GetOnlyChild().CloneTree(),
			rhino.IR_String(testFunctionName),
			fullyQualifiedFunctionName.CloneTree())
	exportCall.PutBooleanProp(rhino.Node_FREE_CALL, exportCall.GetFirstChild().IsName())

	export := rhino.IR_ExprResult(exportCall).SrcrefTree(node)

	export.InsertAfter(node.GetParent())
	e.compiler.ReportChangeToEnclosingScope(export)
}

// ExportTestFunctions_IsTestSuiteArgument is whether this is an object
// literal containing test methods passed to goog.testing.testSuite
//
// t: we only need the scope from this NodeTraversal but want to lazily
// create the scope.
//
// Java: ExportTestFunctions.isTestSuiteArgument(Node, NodeTraversal)
func ExportTestFunctions_IsTestSuiteArgument(n rhino.Node, t *NodeTraversal) bool {
	return n.IsObjectLit() &&
		n.GetParent().IsCall() &&
		n.IsSecondChildOf(n.GetParent()) &&
		ExportTestFunctions_IsGoogTestingTestSuite(t, n.GetPrevious())
}

// ExportTestFunctions_IsGoogTestingTestSuite is whether this node is a
// reference to goog.testing.testSuite, either through a goog.require or by
// its fully qualified name.
//
// If this becomes more broadly useful consider branching it into its own
// class and making more robust to handle forwardDeclares, destructuring
// requires, etc.
//
// Java: ExportTestFunctions.isGoogTestingTestSuite(NodeTraversal, Node)
func ExportTestFunctions_IsGoogTestingTestSuite(t *NodeTraversal, qname rhino.Node) bool {
	qname = NodeUtil_GetCallTargetResolvingIndirectCalls(qname.GetParent())
	if !qname.IsQualifiedName() {
		return false
	}
	root := NodeUtil_GetRootOfQualifiedName(qname)
	rootName := root.GetString()
	rootVar, _ := t.GetScope().GetSlot(rootName).(*Var)
	if rootVar == nil || rootVar.IsGlobal() {
		return qname.MatchesQualifiedName(exportTestFunctions_GOOG_TESTING_TEST_SUITE)
	} else if rootVar.GetScope().(*Scope).IsModuleScope() {
		originalValue := rootVar.GetInitialValue()
		return !originalValue.IsNil() &&
			NodeUtil_IsGoogRequireCall(originalValue) &&
			originalValue.HasTwoChildren() &&
			originalValue.GetSecondChild().GetString() == exportTestFunctions_GOOG_TESTING_TEST_SUITE
	}
	return false
}

// ExportTestFunctions_IsTestFunction is whether a function is recognized as
// a test function. We follow the JsUnit convention for naming (functions
// should start with "test"), and we also check if it has no parameters
// declared.
//
// functionName is the name of the function; returns true if the function
// is recognized as a test function.
//
// Java: ExportTestFunctions.isTestFunction(String)
func ExportTestFunctions_IsTestFunction(functionName *string) bool {
	return functionName != nil &&
		exportTestFunctions_TEST_FUNCTIONS_NAME_PATTERN.MatchString(*functionName)
}
