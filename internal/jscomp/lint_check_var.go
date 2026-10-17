// Port of com/google/javascript/jscomp/lint/CheckVar.java @ closure-compiler v20261015

/*
 * Copyright 2015 The Closure Compiler Authors.
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

// CheckVar checks for `var` (prefer `const` or `let`).
type CheckVar struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler
}

var CheckVar_VAR = DiagnosticType_Disabled(
	"JSC_VAR",
	"Using `var` (prefer `const` or `let`).")

// NewCheckVar is new CheckVar(compiler).
//
// Java: CheckVar(AbstractCompiler)
func NewCheckVar(compiler AbstractCompiler) *CheckVar {
	return &CheckVar{compiler: compiler}
}

// Java: CheckVar.process(Node, Node)
func (c *CheckVar) Process(externs *rhino.Node, root *rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}

// Java: CheckVar.visit(NodeTraversal, Node, Node)
func (c *CheckVar) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if n.IsVar() {
		t.Report(n, CheckVar_VAR)
	}
}
