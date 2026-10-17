// Port of com/google/javascript/jscomp/lint/CheckMissingSemicolon.java @ closure-compiler v20260415

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

// CheckMissingSemicolon checks for statements that should end with a
// semicolon according to the Google style guide.
type CheckMissingSemicolon struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler

	// Go-only: language-adapter — Java's String.charAt is O(1) on its UTF-16
	// array; jutil.StringCharAt converts the whole WTF-8 string per call,
	// which is quadratic when done for every statement of a file. These
	// cache the UTF-16 code units of the last code string read (the same
	// string, so the same chars as code.charAt).
	lastCode      string
	lastCodeUnits []uint16
}

var CheckMissingSemicolon_MISSING_SEMICOLON = DiagnosticType_Disabled(
	"JSC_MISSING_SEMICOLON",
	"Missing semicolon")

// NewCheckMissingSemicolon is new CheckMissingSemicolon(compiler).
//
// Java: CheckMissingSemicolon(AbstractCompiler)
func NewCheckMissingSemicolon(compiler AbstractCompiler) *CheckMissingSemicolon {
	return &CheckMissingSemicolon{compiler: compiler}
}

// Java: CheckMissingSemicolon.process(Node, Node)
func (c *CheckMissingSemicolon) Process(externs *rhino.Node, root *rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}

// Java: CheckMissingSemicolon.visit(NodeTraversal, Node, Node)
func (c *CheckMissingSemicolon) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if !n.IsScript() && NodeUtil_IsStatement(n) && c.ShouldHaveSemicolon(n) {
		c.CheckSemicolon(t, n)
	}
}

// Java: CheckMissingSemicolon.shouldHaveSemicolon(Node)
func (c *CheckMissingSemicolon) ShouldHaveSemicolon(statement *rhino.Node) bool {
	if statement.IsFunction() ||
		statement.IsClass() ||
		statement.IsBlock() ||
		statement.IsLabelName() ||
		statement.IsModuleBody() ||
		(NodeUtil_IsControlStructure(statement) && !statement.IsDo()) {
		return false
	}
	if statement.IsExport() {
		return c.ShouldHaveSemicolon(statement.GetFirstChild())
	}
	return true
}

// Java: CheckMissingSemicolon.checkSemicolon(NodeTraversal, Node)
func (c *CheckMissingSemicolon) CheckSemicolon(t *NodeTraversal, n *rhino.Node) {
	staticSourceFile := n.GetStaticSourceFile()
	if sourceFile, ok := staticSourceFile.(*SourceFile); ok {

		var code string
		failed := false
		jutil.Try(func() {
			code = sourceFile.GetCode()
		}, func(e jutil.IIOException) {
			// We can't read the original source file. Just skip this check.
			failed = true
		})
		if failed {
			return
		}

		length := n.GetLength()
		if length == 0 {
			// This check needs node lengths to work correctly. If we're not in IDE mode, we don't have
			// that information, so just skip the check.
			return
		}
		position := n.GetSourceOffset() + length - 1
		endsWithSemicolon := c.codeCharAt(code, position) == ';'
		if !endsWithSemicolon {
			t.Report(n, CheckMissingSemicolon_MISSING_SEMICOLON)
		}
	}
}

// codeCharAt is code.charAt(position).
//
// Go-only: language-adapter — see lastCodeUnits.
func (c *CheckMissingSemicolon) codeCharAt(code string, position int) uint16 {
	if c.lastCodeUnits == nil || code != c.lastCode {
		c.lastCode = code
		c.lastCodeUnits = jutil.StringToUTF16(code)
	}
	if position < 0 || position >= len(c.lastCodeUnits) {
		panic(jutil.NewStringIndexOutOfBoundsExceptionForIndex(position, len(c.lastCodeUnits)))
	}
	return c.lastCodeUnits[position]
}
