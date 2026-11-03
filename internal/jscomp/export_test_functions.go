// Port of com/google/javascript/jscomp/ExportTestFunctions.java @ closure-compiler v20261028

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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ExportTestFunctions generates goog.exportSymbol for test functions, so they
// can be recognized by the test runner, even if the code is compiled.
type ExportTestFunctions struct {
	compiler               AbstractCompiler
	exportSymbolFunction   string
	exportPropertyFunction string
}

// Go-only: language-adapter — the Java pattern is anchored with ^...$ and
// only used with Matcher.matches(), where Java's `$` can only match at the
// end of the input; Go's regexp gives the same answer (\w is ASCII in both),
// and jutil.Pattern does not accept `$`.
var testFunctionsNamePattern = regexp.MustCompile(
	`^(?:((\w+\.)+prototype\.||window\.)*` +
		`(setUpPage|setUp|shouldRunTests|tearDown|tearDownPage|test[\w\$]+))$`)

const googTestingTestSuite = "goog.testing.testSuite"

// NewExportTestFunctions creates a new export test functions compiler pass.
//
// exportSymbolFunction is the function name used to export symbols in JS;
// exportPropertyFunction is the function name used to export properties in
// JS.
//
// Java: ExportTestFunctions(AbstractCompiler, String, String)
func NewExportTestFunctions(
	compiler AbstractCompiler, exportSymbolFunction string, exportPropertyFunction string) *ExportTestFunctions {
	jutil.CheckNotNull(compiler)
	return &ExportTestFunctions{
		compiler:               compiler,
		exportSymbolFunction:   exportSymbolFunction,
		exportPropertyFunction: exportPropertyFunction,
	}
}

type exportTestFunctionsNodes struct {
	NodeTraversal_AbstractShallowCallbackBase

	outer *ExportTestFunctions
}

var exportTestFunctionsVisitTokens = makeVisitTokenSet(
	rhino.Token_FUNCTION,
	rhino.Token_VAR,
	rhino.Token_LET,
	rhino.Token_CONST,
	rhino.Token_CLASS,
	rhino.Token_ASSIGN,
	rhino.Token_OBJECTLIT,
)

var _ nodeTraversal_visitTokens = (*exportTestFunctionsNodes)(nil)

func (*exportTestFunctionsNodes) visitTokens() *ast.TokenSet {
	return &exportTestFunctionsVisitTokens
}

// Java: ExportTestFunctions.ExportTestFunctionsNodes.visit(NodeTraversal, Node, Node)
func (x *exportTestFunctionsNodes) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if parent.IsNil() {
		return
	}

	if parent.IsScript() || parent.IsModuleBody() {
		if NodeUtil_IsFunctionDeclaration(n) {
			// Check for a test function statement.
			if functionName, ok := NodeUtil_GetNameString(n); ok && isTestFunctionName(functionName) {
				x.outer.exportTestFunctionAsSymbol(functionName, n)
			}
		} else if init, ok := nameDeclaredInit(n); ok {
			if init.IsFunction() {
				// Check for a test function expression.
				if functionName, hasName := NodeUtil_GetNameString(init); hasName && isTestFunctionName(functionName) {
					x.outer.exportTestFunctionAsSymbol(functionName, n)
				}
			} else if cls, isCls := ast.AsClass(init); isCls {
				className, hasClassName := NodeUtil_GetNameString(init)
				x.exportClassWithName(cls, className, hasClassName, n)
			}
		} else if cls, ok := ast.AsClass(n); ok {
			x.exportClass(cls)
		}
	} else if NodeUtil_IsExprAssign(parent) {
		// Check for a test method assignment.
		grandparent := parent.GetParent()
		if !grandparent.IsNil() && (grandparent.IsScript() || grandparent.IsModuleBody()) {
			//                                    NAME/(GETPROP -> ... -> NAME)
			// SCRIPT -> EXPR_RESULT -> ASSIGN ->
			//                                    FUNCTION/CLASS
			assign := ast.MustAssign(n)
			lhs := assign.Target().Node()
			rhs := assign.Value().Node()
			nodeName, hasNodeName := lhs.GetQualifiedNameString()

			if rhs.IsFunction() {
				if hasNodeName && isTestFunctionName(nodeName) {
					if lhs.IsName() {
						x.outer.exportTestFunctionAsSymbol(nodeName, parent)
					} else {
						x.outer.exportTestFunctionAsProperty(ast.MustGetProp(lhs), assign)
					}
				}
			} else if cls, ok := ast.AsClass(rhs); ok {
				x.exportClassWithName(cls, nodeName, hasNodeName, parent)
			}
		}
	} else if objLit, ok := asTestSuiteArgument(n, t); ok {
		for member := range ast.ReadAhead(objLit.FirstMember()) {
			if sk, isSk := member.AsStringKey(); isSk && !sk.IsQuoted() {
				sk.SetQuoted(true)
				x.outer.compiler.ReportChangeToEnclosingScope(sk.Node())
			} else if mfd, isMfd := member.AsMemberFunctionDef(); isMfd {
				x.rewriteMemberDefInObjLit(mfd, objLit)
			}
		}
	}
}

// Java: ExportTestFunctions.ExportTestFunctionsNodes.exportClass(Node)
func (x *exportTestFunctionsNodes) exportClass(cls ast.Class) {
	className, hasClassName := NodeUtil_GetNameString(cls.Node())
	x.exportClassWithName(cls, className, hasClassName, cls.Node())
}

// exportClassWithName exports test methods on an ES6 class. When !hasClassName,
// className renders as "null" to match Java string concatenation.
//
// Java: ExportTestFunctions.ExportTestFunctionsNodes.exportClass(Node, String, Node)
func (x *exportTestFunctionsNodes) exportClassWithName(cls ast.Class, className string, hasClassName bool, addAfter rhino.Node) {
	if !hasClassName {
		className = "null"
	}
	for member := range cls.Body().Members() {
		mfd, ok := member.AsMemberFunctionDef()
		if !ok {
			continue
		}
		methodName := mfd.Key()
		if !isTestFunctionName(methodName) {
			continue
		}
		memberNode := mfd.Node()
		functionRef := className + ".prototype." + methodName
		classRef := className + ".prototype"

		exportCallTarget := NodeUtil_NewQNameWithOriginalName(
			x.outer.compiler, x.outer.exportPropertyFunction, memberNode, methodName)
		call := rhino.IR_Call(exportCallTarget)
		if exportCallTarget.IsName() {
			ast.MustCall(call).SetFreeCall(true)
		}

		call.AddChildToBack(
			NodeUtil_NewQNameWithOriginalName(x.outer.compiler, classRef, memberNode, classRef))
		call.AddChildToBack(rhino.IR_String(methodName))
		call.AddChildToBack(
			NodeUtil_NewQNameWithOriginalName(x.outer.compiler, functionRef, memberNode, functionRef))

		expression := rhino.IR_ExprResult(call).SrcrefTreeIfMissing(memberNode)

		expression.InsertAfter(addAfter)
		x.outer.compiler.ReportChangeToEnclosingScope(expression)
		addAfter = expression
	}
}

// rewriteMemberDefInObjLit converts a member function into a quoted string
// key to avoid property renaming.
//
// Java: ExportTestFunctions.ExportTestFunctionsNodes.rewriteMemberDefInObjLit(Node, Node)
func (x *exportTestFunctionsNodes) rewriteMemberDefInObjLit(mfd ast.MemberFunctionDef, objLit ast.ObjectLit) {
	name := mfd.Key()
	fnNode := mfd.Function().Detach().Node()
	stringKey := rhino.IR_StringKeyWithValue(name, fnNode)
	mfd.Node().ReplaceWith(stringKey)
	ast.MustStringKey(stringKey).SetQuoted(true)
	stringKey.SetJSDocInfo(mfd.Node().GetJSDocInfo())
	x.outer.compiler.ReportChangeToEnclosingScope(objLit.Node())
}

// nameDeclaredInit gets the initializer node that corresponds to an expression
// declared with var, let or const (`VAR/LET/CONST -> NAME -> NODE`).
//
// Java: ExportTestFunctions.ExportTestFunctionsNodes.getNameDeclaredGrandchild(Node)
func nameDeclaredInit(node rhino.Node) (rhino.Node, bool) {
	decl, ok := ast.AsNameDeclaration(node)
	if !ok {
		return rhino.NilNode, false
	}
	first, ok := decl.FirstDeclarator()
	if !ok {
		return rhino.NilNode, false
	}
	name, ok := first.AsName()
	if !ok {
		return rhino.NilNode, false
	}
	init, ok := name.Init()
	if !ok {
		return rhino.NilNode, false
	}
	return init.Node(), true
}

// Java: ExportTestFunctions.process(Node, Node)
func (e *ExportTestFunctions) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(e.compiler, root, &exportTestFunctionsNodes{outer: e})
}

// exportTestFunctionAsSymbol adds exportSymbol(testFunctionName, testFunction);
//
// Java: ExportTestFunctions.exportTestFunctionAsSymbol(String, Node)
func (e *ExportTestFunctions) exportTestFunctionAsSymbol(testFunctionName string, node rhino.Node) {
	exportCallTarget := NodeUtil_NewQNameWithOriginalName(e.compiler,
		e.exportSymbolFunction, node, testFunctionName)
	call := rhino.IR_Call(exportCallTarget)
	if exportCallTarget.IsName() {
		ast.MustCall(call).SetFreeCall(true)
	}
	call.AddChildToBack(rhino.IR_String(testFunctionName))
	call.AddChildToBack(NodeUtil_NewQNameWithOriginalName(e.compiler,
		testFunctionName, node, testFunctionName))

	expression := rhino.IR_ExprResult(call).SrcrefTreeIfMissing(node)

	expression.InsertAfter(node)
	e.compiler.ReportChangeToEnclosingScope(expression)
}

// exportTestFunctionAsProperty adds exportProperty() of the test function
// name on the prototype object.
//
// Java: ExportTestFunctions.exportTestFunctionAsProperty(Node, Node)
func (e *ExportTestFunctions) exportTestFunctionAsProperty(fullyQualifiedFunctionName ast.GetProp, assign ast.Assign) {
	lhs := assign.Target().Node()
	testFunctionName := NodeUtil_GetPrototypePropertyName(lhs)
	qname, _ := lhs.GetQualifiedNameString()
	if after, ok := strings.CutPrefix(qname, "window."); ok {
		testFunctionName = after
	}

	callee := NodeUtil_NewQName(e.compiler, e.exportPropertyFunction)
	exportCall := rhino.IR_Call(
		callee,
		fullyQualifiedFunctionName.Object().Node().CloneTree(),
		rhino.IR_String(testFunctionName),
		fullyQualifiedFunctionName.Node().CloneTree())
	ast.MustCall(exportCall).SetFreeCall(callee.IsName())

	export := rhino.IR_ExprResult(exportCall).SrcrefTree(assign.Node())

	export.InsertAfter(assign.Node().GetParent())
	e.compiler.ReportChangeToEnclosingScope(export)
}

// asTestSuiteArgument reports whether n is an object literal containing test
// methods passed to goog.testing.testSuite.
//
// t: we only need the scope from this NodeTraversal but want to lazily
// create the scope.
//
// Java: ExportTestFunctions.isTestSuiteArgument(Node, NodeTraversal)
func asTestSuiteArgument(n rhino.Node, t *NodeTraversal) (ast.ObjectLit, bool) {
	objLit, ok := ast.AsObjectLit(n)
	if !ok {
		return ast.ObjectLit{}, false
	}
	call, ok := ast.AsCall(n.GetParent())
	if !ok {
		return ast.ObjectLit{}, false
	}
	firstArg, ok := call.FirstArg()
	if !ok || !firstArg.Is(n) || !isGoogTestingTestSuiteCall(t, call) {
		return ast.ObjectLit{}, false
	}
	return objLit, true
}

// isGoogTestingTestSuiteCall reports whether call invokes
// goog.testing.testSuite, either through a goog.require or by its fully
// qualified name.
//
// If this becomes more broadly useful consider branching it into its own
// class and making more robust to handle forwardDeclares, destructuring
// requires, etc.
//
// Java: ExportTestFunctions.isGoogTestingTestSuite(NodeTraversal, Node)
func isGoogTestingTestSuiteCall(t *NodeTraversal, call ast.Call) bool {
	qname := NodeUtil_GetCallTargetResolvingIndirectCalls(call.Node())
	if !qname.IsQualifiedName() {
		return false
	}
	root := NodeUtil_GetRootOfQualifiedName(qname)
	rootVar, _ := t.GetScope().GetSlot(root.GetString()).(*Var)
	if rootVar == nil || rootVar.IsGlobal() {
		return qname.MatchesQualifiedName(googTestingTestSuite)
	}
	if rootVar.GetScope().(*Scope).IsModuleScope() {
		originalValue := rootVar.GetInitialValue()
		if reqCall, ok := ast.AsCall(originalValue); ok && NodeUtil_IsGoogRequireCall(originalValue) {
			if arg, hasArg := reqCall.OnlyArg(); hasArg {
				if str, isStr := arg.AsStringLit(); isStr {
					return str.Value() == googTestingTestSuite
				}
			}
		}
	}
	return false
}

func isTestFunctionName(functionName string) bool {
	return testFunctionsNamePattern.MatchString(functionName)
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
func ExportTestFunctions_IsTestFunction(functionName string) bool {
	return isTestFunctionName(functionName)
}
