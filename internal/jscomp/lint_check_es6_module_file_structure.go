// Port of com/google/javascript/jscomp/lint/CheckEs6ModuleFileStructure.java @ closure-compiler v20261015

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

package jscomp

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CheckEs6ModuleFileStructure checks the file structure of ES6 modules.
type CheckEs6ModuleFileStructure struct {
	NodeTraversal_AbstractPreOrderCallbackBase
	compiler          AbstractCompiler
	orderedStatements *jutil.TreeSet[CheckEs6ModuleFileStructure_OrderedStatement]
}

var CheckEs6ModuleFileStructure_MUST_COME_BEFORE = DiagnosticType_Warning(
	"JSC_MUST_COME_BEFORE_IN_ES6_MODULE", "In ES6 modules, {0} should come before {1}.")

// CheckEs6ModuleFileStructure_OrderedStatement: statements that must appear
// in a certain order within ES6 modules (for the sake of style).
type CheckEs6ModuleFileStructure_OrderedStatement int32

const (
	CheckEs6ModuleFileStructure_OrderedStatement_IMPORT CheckEs6ModuleFileStructure_OrderedStatement = iota
	CheckEs6ModuleFileStructure_OrderedStatement_DECLARE_MODULE_ID
	CheckEs6ModuleFileStructure_OrderedStatement_GOOG_REQUIRE
	CheckEs6ModuleFileStructure_OrderedStatement_OTHER
)

var checkEs6ModuleFileStructure_OrderedStatement_names = [...]string{
	"IMPORT", "DECLARE_MODULE_ID", "GOOG_REQUIRE", "OTHER",
}

// The description field of each constant, indexed by ordinal.
var checkEs6ModuleFileStructure_OrderedStatement_description = [...]string{
	"import statements",
	"a call to goog.declareModuleId()",
	"calls to goog.require()",
	"other statements",
}

// CheckEs6ModuleFileStructure_OrderedStatement_values is OrderedStatement.values().
func CheckEs6ModuleFileStructure_OrderedStatement_values() []CheckEs6ModuleFileStructure_OrderedStatement {
	return []CheckEs6ModuleFileStructure_OrderedStatement{
		CheckEs6ModuleFileStructure_OrderedStatement_IMPORT,
		CheckEs6ModuleFileStructure_OrderedStatement_DECLARE_MODULE_ID,
		CheckEs6ModuleFileStructure_OrderedStatement_GOOG_REQUIRE,
		CheckEs6ModuleFileStructure_OrderedStatement_OTHER,
	}
}

// Name is Enum.name().
func (s CheckEs6ModuleFileStructure_OrderedStatement) Name() string {
	return checkEs6ModuleFileStructure_OrderedStatement_names[s]
}

// Ordinal is Enum.ordinal().
func (s CheckEs6ModuleFileStructure_OrderedStatement) Ordinal() int { return int(s) }

// String is the overridden toString(): the description.
//
// Java: CheckEs6ModuleFileStructure.OrderedStatement.toString()
func (s CheckEs6ModuleFileStructure_OrderedStatement) String() string {
	return checkEs6ModuleFileStructure_OrderedStatement_description[s]
}

// checkEs6ModuleFileStructure_OrderedStatement_compareTo is Enum.compareTo:
// the natural (ordinal) order of the TreeSet.
func checkEs6ModuleFileStructure_OrderedStatement_compareTo(a, b CheckEs6ModuleFileStructure_OrderedStatement) int {
	return int(a) - int(b)
}

// NewCheckEs6ModuleFileStructure is new CheckEs6ModuleFileStructure(compiler).
//
// Java: CheckEs6ModuleFileStructure(AbstractCompiler)
func NewCheckEs6ModuleFileStructure(compiler AbstractCompiler) *CheckEs6ModuleFileStructure {
	return &CheckEs6ModuleFileStructure{
		compiler:          compiler,
		orderedStatements: jutil.NewTreeSet(checkEs6ModuleFileStructure_OrderedStatement_compareTo),
	}
}

// Java: CheckEs6ModuleFileStructure.shouldTraverse(NodeTraversal, Node, Node)
func (c *CheckEs6ModuleFileStructure) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_ROOT:
		return true
	case rhino.Token_SCRIPT:
		return n.GetBooleanProp(rhino.Node_ES6_MODULE)
	case rhino.Token_MODULE_BODY:
		c.orderedStatements.Clear()
		return true
	case rhino.Token_IMPORT:
		c.VisitImport(t, n)
		return false
	case rhino.Token_EXPR_RESULT:
		return c.VisitExprResult(t, n, parent)
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		return c.VisitDeclaration(t, n, parent)
	default:
		if parent != nil && parent.IsModuleBody() {
			c.CheckOrder(t, n, CheckEs6ModuleFileStructure_OrderedStatement_OTHER)
		}
		return false
	}
}

// Java: CheckEs6ModuleFileStructure.checkOrder(NodeTraversal, Node, OrderedStatement)
func (c *CheckEs6ModuleFileStructure) CheckOrder(t *NodeTraversal, n *rhino.Node, statement CheckEs6ModuleFileStructure_OrderedStatement) {
	c.orderedStatements.Add(statement)
	// Go-only: language-adapter — orderedStatements.tailSet(statement, /* inclusive= */ false):
	// the elements of the sorted set strictly greater than statement, in order.
	var outOfOrder []CheckEs6ModuleFileStructure_OrderedStatement
	for s := range c.orderedStatements.All() {
		if checkEs6ModuleFileStructure_OrderedStatement_compareTo(s, statement) > 0 {
			outOfOrder = append(outOfOrder, s)
		}
	}
	if len(outOfOrder) != 0 {
		parts := make([]string, len(outOfOrder))
		for i, s := range outOfOrder {
			parts[i] = s.String()
		}
		t.Report(n, CheckEs6ModuleFileStructure_MUST_COME_BEFORE, statement.String(), strings.Join(parts, ", "))
	}
}

var checkEs6ModuleFileStructure_GOOG_DECLAREMODULEID = rhino.QualifiedName_Of("goog.declareModuleId")

var checkEs6ModuleFileStructure_GOOG_REQUIRE = rhino.QualifiedName_Of("goog.require")

// Java: CheckEs6ModuleFileStructure.visitExprResult(NodeTraversal, Node, Node)
func (c *CheckEs6ModuleFileStructure) VisitExprResult(t *NodeTraversal, exprResult *rhino.Node, parent *rhino.Node) bool {
	if parent.IsModuleBody() && exprResult.GetFirstChild().IsCall() {
		call := exprResult.GetFirstChild()
		if checkEs6ModuleFileStructure_GOOG_DECLAREMODULEID.Matches(call.GetFirstChild()) {
			c.CheckOrder(t, call, CheckEs6ModuleFileStructure_OrderedStatement_DECLARE_MODULE_ID)
			return false
		} else if checkEs6ModuleFileStructure_GOOG_REQUIRE.Matches(call.GetFirstChild()) {
			c.CheckOrder(t, call, CheckEs6ModuleFileStructure_OrderedStatement_GOOG_REQUIRE)
			return false
		}
	}

	c.CheckOrder(t, exprResult, CheckEs6ModuleFileStructure_OrderedStatement_OTHER)
	return true
}

// Java: CheckEs6ModuleFileStructure.visitDeclaration(NodeTraversal, Node, Node)
func (c *CheckEs6ModuleFileStructure) VisitDeclaration(t *NodeTraversal, declaration *rhino.Node, parent *rhino.Node) bool {
	if parent.IsModuleBody() &&
		declaration.HasOneChild() &&
		declaration.GetFirstChild().HasOneChild() &&
		declaration.GetFirstFirstChild().IsCall() {
		call := declaration.GetFirstFirstChild()
		if checkEs6ModuleFileStructure_GOOG_REQUIRE.Matches(call.GetFirstChild()) {
			c.CheckOrder(t, call, CheckEs6ModuleFileStructure_OrderedStatement_GOOG_REQUIRE)
			return false
		}
	}

	c.CheckOrder(t, declaration, CheckEs6ModuleFileStructure_OrderedStatement_OTHER)
	return true
}

// Java: CheckEs6ModuleFileStructure.visitImport(NodeTraversal, Node)
func (c *CheckEs6ModuleFileStructure) VisitImport(t *NodeTraversal, importNode *rhino.Node) {
	c.CheckOrder(t, importNode, CheckEs6ModuleFileStructure_OrderedStatement_IMPORT)
}

// Java: CheckEs6ModuleFileStructure.process(Node, Node)
func (c *CheckEs6ModuleFileStructure) Process(externs *rhino.Node, root *rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}
