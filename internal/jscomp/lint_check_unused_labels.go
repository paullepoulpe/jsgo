// Port of com/google/javascript/jscomp/lint/CheckUnusedLabels.java @ closure-compiler v20261015

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
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CheckUnusedLabels checks for unused labels blocks. This can help catching
// errors like:
//
//	() => {a: 2}  // Returns undefined, not an Object
//
// Inspired by ESLint
// (https://github.com/eslint/eslint/blob/master/lib/rules/no-unused-labels.js)
type CheckUnusedLabels struct {
	compiler       AbstractCompiler
	currentContext *CheckUnusedLabels_LabelContext
}

var CheckUnusedLabels_UNUSED_LABEL = DiagnosticType_Disabled(
	"JSC_UNUSED_LABEL", "Unused label {0}.")

// CheckUnusedLabels_LabelContext is the private static class LabelContext.
type CheckUnusedLabels_LabelContext struct {
	name   string
	parent *CheckUnusedLabels_LabelContext
	used   bool
}

// Java: CheckUnusedLabels.LabelContext(String, LabelContext)
func newCheckUnusedLabels_LabelContext(name string, parent *CheckUnusedLabels_LabelContext) *CheckUnusedLabels_LabelContext {
	return &CheckUnusedLabels_LabelContext{name: name, parent: parent}
}

// NewCheckUnusedLabels is new CheckUnusedLabels(compiler).
//
// Java: CheckUnusedLabels(AbstractCompiler)
func NewCheckUnusedLabels(compiler AbstractCompiler) *CheckUnusedLabels {
	return &CheckUnusedLabels{compiler: compiler}
}

// Java: CheckUnusedLabels.process(Node, Node)
func (c *CheckUnusedLabels) Process(externs *rhino.Node, root *rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}

// Java: CheckUnusedLabels.shouldTraverse(NodeTraversal, Node, Node)
func (c *CheckUnusedLabels) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_BREAK, rhino.Token_CONTINUE:
		if n.HasChildren() {
			temp := c.currentContext
			for temp != nil {
				if temp.name == n.GetFirstChild().GetString() {
					temp.used = true
					break
				}
				temp = temp.parent
			}
		}
		return false
	case rhino.Token_LABEL:
		c.currentContext = newCheckUnusedLabels_LabelContext(n.GetFirstChild().GetString(), c.currentContext)
	default:
	}
	return true
}

// Java: CheckUnusedLabels.visit(NodeTraversal, Node, Node)
func (c *CheckUnusedLabels) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if n.IsLabel() && c.currentContext != nil {
		if !c.currentContext.used {
			t.Report(n, CheckUnusedLabels_UNUSED_LABEL, n.GetFirstChild().GetString())
		}
		c.currentContext = c.currentContext.parent
	}
}
