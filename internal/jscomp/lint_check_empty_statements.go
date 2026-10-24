// Port of com/google/javascript/jscomp/lint/CheckEmptyStatements.java @ closure-compiler v20261021

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

// CheckEmptyStatements checks for empty statements (i.e. extraneous
// semicolons)
type CheckEmptyStatements struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler
}

var CheckEmptyStatements_USELESS_EMPTY_STATEMENT = DiagnosticType_Disabled(
	"JSC_USELESS_EMPTY_STATEMENT",
	"Useless empty statement. Remove semicolon.")

// NewCheckEmptyStatements is new CheckEmptyStatements(compiler).
//
// Java: CheckEmptyStatements(AbstractCompiler)
func NewCheckEmptyStatements(compiler AbstractCompiler) *CheckEmptyStatements {
	return &CheckEmptyStatements{compiler: compiler}
}

// Java: CheckEmptyStatements.process(Node, Node)
func (c *CheckEmptyStatements) Process(externs *rhino.Node, root *rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}

// Java: CheckEmptyStatements.visit(NodeTraversal, Node, Node)
func (c *CheckEmptyStatements) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if n.IsEmpty() && NodeUtil_IsStatement(n) {
		t.Report(n, CheckEmptyStatements_USELESS_EMPTY_STATEMENT)
	}
}
