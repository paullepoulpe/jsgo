// Port of com/google/javascript/jscomp/lint/CheckPrimitiveAsObject.java @ closure-compiler v20260415

/*
 * Copyright 2016 The Closure Compiler Authors.
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

// CheckPrimitiveAsObject checks for explicit creation of the object
// equivalents of primitive types (e.g. Boolean instead of boolean) and their
// use in type declarations.
//
// Using these is confusing and gives no benefit. For example, the result of
// typeof (new Boolean(true)) is "object". and the result of (new
// Boolean(false)) ? "true" : "false" is "true".
type CheckPrimitiveAsObject struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler
}

var CheckPrimitiveAsObject_NEW_PRIMITIVE_OBJECT = DiagnosticType_Warning("JSC_PRIMITIVE_OBJECT", "Explicit creation of a {0} object.")

var CheckPrimitiveAsObject_PRIMITIVE_OBJECT_DECLARATION = DiagnosticType_Warning(
	"JSC_PRIMITIVE_OBJECT_DECLARATION",
	"Declaration of {0} object instead of primitive type.")

var checkPrimitiveAsObject_PRIMITIVE_OBJECT_CONSTRUCTORS = jutil.NewLinkedHashSet("Boolean", "Number", "String")

// NewCheckPrimitiveAsObject is new CheckPrimitiveAsObject(compiler).
//
// Java: CheckPrimitiveAsObject(AbstractCompiler)
func NewCheckPrimitiveAsObject(compiler AbstractCompiler) *CheckPrimitiveAsObject {
	return &CheckPrimitiveAsObject{compiler: compiler}
}

// Java: CheckPrimitiveAsObject.process(Node, Node)
func (c *CheckPrimitiveAsObject) Process(externs *rhino.Node, root *rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}

// Java: CheckPrimitiveAsObject.visit(NodeTraversal, Node, Node)
func (c *CheckPrimitiveAsObject) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	c.CheckForPrimitiveObjectConstructor(t, n)
	c.CheckForPrimitiveObjectDeclaration(t, n)
}

// Java: CheckPrimitiveAsObject.checkForPrimitiveObjectDeclaration(NodeTraversal, Node)
func (c *CheckPrimitiveAsObject) CheckForPrimitiveObjectDeclaration(t *NodeTraversal, n *rhino.Node) {
	jsDocInfo := n.GetJSDocInfo()

	if jsDocInfo != nil {
		for _, typeRoot := range jsDocInfo.GetTypeNodes() {
			c.CheckTypeNodeForPrimitiveObjectDeclaration(t, typeRoot)
		}
	}
}

// Java: CheckPrimitiveAsObject.checkTypeNodeForPrimitiveObjectDeclaration(NodeTraversal, Node)
func (c *CheckPrimitiveAsObject) CheckTypeNodeForPrimitiveObjectDeclaration(t *NodeTraversal, typeRoot *rhino.Node) {
	NodeUtil_VisitPreOrder(
		typeRoot,
		NodeUtil_VisitorFunc(func(node *rhino.Node) {
			if node.IsStringLit() {
				typeName := node.GetString()
				if checkPrimitiveAsObject_PRIMITIVE_OBJECT_CONSTRUCTORS.Contains(typeName) {
					t.Report(node, CheckPrimitiveAsObject_PRIMITIVE_OBJECT_DECLARATION, typeName)
				}
			}
		}))
}

// Java: CheckPrimitiveAsObject.checkForPrimitiveObjectConstructor(NodeTraversal, Node)
func (c *CheckPrimitiveAsObject) CheckForPrimitiveObjectConstructor(t *NodeTraversal, n *rhino.Node) {
	if n.IsNew() {
		constructorFunction := n.GetFirstChild()
		if constructorFunction.IsName() {
			constructorName := constructorFunction.GetString()
			if checkPrimitiveAsObject_PRIMITIVE_OBJECT_CONSTRUCTORS.Contains(constructorName) {
				t.Report(n, CheckPrimitiveAsObject_NEW_PRIMITIVE_OBJECT, constructorName)
			}
		}
	}
}
