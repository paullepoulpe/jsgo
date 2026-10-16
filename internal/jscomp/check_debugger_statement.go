// Port of com/google/javascript/jscomp/CheckDebuggerStatement.java @ closure-compiler v20260415

/*
 * Copyright 2011 The Closure Compiler Authors.
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

// CheckDebuggerStatement checks for the presence of the "debugger" statement
// in JavaScript code. It is appropriate to use this statement while developing
// JavaScript; however, it is generally undesirable to include it in production
// code.
type CheckDebuggerStatement struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	compiler AbstractCompiler
}

// Java: CheckDebuggerStatement.DEBUGGER_STATEMENT_PRESENT
var CheckDebuggerStatement_DEBUGGER_STATEMENT_PRESENT = DiagnosticType_Disabled(
	"JSC_DEBUGGER_STATEMENT_PRESENT",
	"Using the debugger statement can halt your application if the user has a JavaScript"+
		" debugger running.\n"+
		"To disable this check when you want to do debugging, you can suppress this message"+
		" like this:\n"+
		// NOTE: Single quotes around curly braces are necessary to prevent them being
		// interpreted as a placeholder for replacement. See MessageFormat javadoc.
		"  /** @suppress '{'checkDebuggerStatement'}' */\n"+
		"  debugger;")

// Java: CheckDebuggerStatement(AbstractCompiler)
func NewCheckDebuggerStatement(compiler AbstractCompiler) *CheckDebuggerStatement {
	c := &CheckDebuggerStatement{}
	c.compiler = compiler
	return c
}

// Java: CheckDebuggerStatement.process(Node, Node)
func (c *CheckDebuggerStatement) Process(externs *rhino.Node, root *rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}

// Java: CheckDebuggerStatement.visit(NodeTraversal, Node, Node)
func (c *CheckDebuggerStatement) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if n.IsDebugger() {
		t.Report(n, CheckDebuggerStatement_DEBUGGER_STATEMENT_PRESENT)
	}
}
