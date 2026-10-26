// Port of com/google/javascript/jscomp/lint/CheckDuplicateCase.java @ closure-compiler v20261021

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

// CheckDuplicateCase checks for duplicate case labels in a switch statement
// Eg: switch (foo) { case 1: case 1: }
//
// This is normally an indication of a programmer error.
//
// Inspired by ESLint
// (https://github.com/eslint/eslint/blob/master/lib/rules/no-duplicate-case.js)
type CheckDuplicateCase struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler
}

var CheckDuplicateCase_DUPLICATE_CASE = DiagnosticType_Warning("JSC_DUPLICATE_CASE", "Duplicate case in a switch statement.")

// NewCheckDuplicateCase is new CheckDuplicateCase(compiler).
//
// Java: CheckDuplicateCase(AbstractCompiler)
func NewCheckDuplicateCase(compiler AbstractCompiler) *CheckDuplicateCase {
	return &CheckDuplicateCase{compiler: compiler}
}

// Java: CheckDuplicateCase.process(Node, Node)
func (c *CheckDuplicateCase) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}

// Java: CheckDuplicateCase.visit(NodeTraversal, Node, Node)
func (c *CheckDuplicateCase) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsSwitch() {
		switchBody := n.GetSecondChild()
		first := switchBody.GetFirstChild()
		if first.IsNil() || first.GetNext().IsNil() {
			return
		}
		cases := make(map[string]bool)
		for curr := first; !curr.IsNil(); curr = curr.GetNext() {
			source := c.compiler.ToSourceOfNode(curr.GetFirstChild())
			if cases[source] {
				t.Report(curr, CheckDuplicateCase_DUPLICATE_CASE)
			} else {
				cases[source] = true
			}
		}
	}
}
