// Port of com/google/javascript/jscomp/lint/CheckPrototypeProperties.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CheckPrototypeProperties checks when a mutable property is assigned to a
// prototype. This is generally undesirable because it can lead to the
// following unexpected situation.
//
//	/** @constructor * /
//	function MyClass() {}
//	MyClass.prototype.prop = [];
//	x = new MyClass;
//	y = new MyClass;
//	x.prop.push(1);
//	console.log(y.prop) // [1]
type CheckPrototypeProperties struct {
	compiler AbstractCompiler
}

var CheckPrototypeProperties_ILLEGAL_PROTOTYPE_MEMBER = DiagnosticType_Disabled(
	"JSC_ILLEGAL_PROTOTYPE_MEMBER",
	"Prototype property {0} should be a primitive, not an Array or Object.")

// NewCheckPrototypeProperties is new CheckPrototypeProperties(compiler).
//
// Java: CheckPrototypeProperties(AbstractCompiler)
func NewCheckPrototypeProperties(compiler AbstractCompiler) *CheckPrototypeProperties {
	return &CheckPrototypeProperties{compiler: compiler}
}

// Java: CheckPrototypeProperties.process(Node, Node)
func (c *CheckPrototypeProperties) Process(externs *rhino.Node, root *rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}

// Java: CheckPrototypeProperties.visit(NodeTraversal, Node, Node)
func (c *CheckPrototypeProperties) Visit(unused *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if NodeUtil_IsPrototypePropertyDeclaration(n) {
		assign := n.GetFirstChild()
		rhs := assign.GetLastChild()
		if rhs.IsArrayLit() || rhs.IsObjectLit() {
			jsDoc := NodeUtil_GetBestJSDocInfo(rhs)
			if jsDoc != nil && jsDoc.HasEnumParameterType() {
				// Don't report for @enum's on the prototype. Sometimes this is necessary, for example,
				// to expose the enum values to an Angular template.
				return
			}
			propName := assign.GetFirstChild().GetLastChild().GetString()
			c.compiler.Report(JSError_MakeWithNode(assign, CheckPrototypeProperties_ILLEGAL_PROTOTYPE_MEMBER, propName))
		}
	}
}

// Java: CheckPrototypeProperties.shouldTraverse(NodeTraversal, Node, Node)
func (c *CheckPrototypeProperties) ShouldTraverse(nodeTraversal *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	return true
}
