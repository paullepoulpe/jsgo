// Port of com/google/javascript/jscomp/lint/CheckVar.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/lint/CheckEmptyStatements.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/lint/CheckDuplicateCase.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/lint/CheckUnusedLabels.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/lint/CheckMissingSemicolon.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/lint/CheckUselessBlocks.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ---------------------------------------------------------------------------
// CheckVar
// ---------------------------------------------------------------------------

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
func (c *CheckVar) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}

// Java: CheckVar.visit(NodeTraversal, Node, Node)
func (c *CheckVar) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsVar() {
		t.Report(n, CheckVar_VAR)
	}
}

// ---------------------------------------------------------------------------
// CheckEmptyStatements
// ---------------------------------------------------------------------------

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
func (c *CheckEmptyStatements) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}

// Java: CheckEmptyStatements.visit(NodeTraversal, Node, Node)
func (c *CheckEmptyStatements) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if _, ok := ast.AsEmpty(n); ok && NodeUtil_IsStatement(n) {
		t.Report(n, CheckEmptyStatements_USELESS_EMPTY_STATEMENT)
	}
}

// ---------------------------------------------------------------------------
// CheckDuplicateCase
// ---------------------------------------------------------------------------

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
	sw, ok := ast.AsSwitch(n)
	if !ok {
		return
	}
	body := sw.Body()
	first, ok := body.FirstCase()
	if !ok {
		return
	}
	if _, hasSecond := first.Next(); !hasSecond {
		return
	}
	cases := make(map[string]bool)
	for curr := range body.Cases() {
		var child rhino.Node
		if caseClause, isCase := curr.AsCase(); isCase {
			child = caseClause.Value().Node()
		} else if defClause, isDef := curr.AsDefaultCase(); isDef {
			child = defClause.Body().Node()
		} else {
			child = curr.Node().GetFirstChild()
		}
		source := c.compiler.ToSourceOfNode(child)
		if cases[source] {
			t.Report(curr.Node(), CheckDuplicateCase_DUPLICATE_CASE)
		} else {
			cases[source] = true
		}
	}
}

// ---------------------------------------------------------------------------
// CheckUnusedLabels
// ---------------------------------------------------------------------------

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
func (c *CheckUnusedLabels) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}

// Java: CheckUnusedLabels.shouldTraverse(NodeTraversal, Node, Node)
func (c *CheckUnusedLabels) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if jmp, ok := ast.AsJump(n); ok {
		if labelNode, hasLabel := jmp.Label(); hasLabel {
			target := labelNode.Ident()
			for temp := c.currentContext; temp != nil; temp = temp.parent {
				if temp.name == target {
					temp.used = true
					break
				}
			}
		}
		return false
	}
	if lbl, ok := ast.AsLabel(n); ok {
		c.currentContext = newCheckUnusedLabels_LabelContext(lbl.Name().Ident(), c.currentContext)
	}
	return true
}

// Java: CheckUnusedLabels.visit(NodeTraversal, Node, Node)
func (c *CheckUnusedLabels) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if lbl, ok := ast.AsLabel(n); ok && c.currentContext != nil {
		if !c.currentContext.used {
			t.Report(n, CheckUnusedLabels_UNUSED_LABEL, lbl.Name().Ident())
		}
		c.currentContext = c.currentContext.parent
	}
}

// ---------------------------------------------------------------------------
// CheckMissingSemicolon
// ---------------------------------------------------------------------------

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
func (c *CheckMissingSemicolon) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}

// Java: CheckMissingSemicolon.visit(NodeTraversal, Node, Node)
func (c *CheckMissingSemicolon) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if !n.IsScript() && NodeUtil_IsStatement(n) && c.ShouldHaveSemicolon(n) {
		c.CheckSemicolon(t, n)
	}
}

// Java: CheckMissingSemicolon.shouldHaveSemicolon(Node)
func (c *CheckMissingSemicolon) ShouldHaveSemicolon(statement rhino.Node) bool {
	if statement.IsFunction() ||
		statement.IsClass() ||
		statement.IsBlock() ||
		statement.IsLabelName() ||
		statement.IsModuleBody() ||
		(NodeUtil_IsControlStructure(statement) && !statement.IsDo()) {
		return false
	}
	if exp, ok := ast.AsExport(statement); ok {
		return c.ShouldHaveSemicolon(exp.Body().Node())
	}
	return true
}

// Java: CheckMissingSemicolon.checkSemicolon(NodeTraversal, Node)
func (c *CheckMissingSemicolon) CheckSemicolon(t *NodeTraversal, n rhino.Node) {
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
		if c.codeCharAt(code, position) != ';' {
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

// ---------------------------------------------------------------------------
// CheckUselessBlocks
// ---------------------------------------------------------------------------

// CheckUselessBlocks checks for useless blocks. A block is considered useful
// if it is part of a control structure like if / else / while / switch / etc.
// or if it contains any block-scoped variables (let, const, class, or
// function declarations). Otherwise there is no reason to use it so it is
// likely a mistake. This would catch the classic error:
//
// return {foo: 'bar'};
//
// or more contrived cases like:
//
// if (denied) { showAccessDenied(); } { grantAccess(); }
//
// Inspired by ESLint (https://github.com/eslint/eslint/blob/master/lib/rules/no-lone-blocks.js)
type CheckUselessBlocks struct {
	compiler AbstractCompiler
	// Go-only: language-adapter — the ArrayDeque used as a stack (push, pop
	// and peek at its head) is a slice whose last element is the head.
	loneBlocks []rhino.Node
}

var CheckUselessBlocks_USELESS_BLOCK = DiagnosticType_Disabled(
	"JSC_USELESS_BLOCK", "Useless block.")

// NewCheckUselessBlocks is new CheckUselessBlocks(compiler).
//
// Java: CheckUselessBlocks(AbstractCompiler)
func NewCheckUselessBlocks(compiler AbstractCompiler) *CheckUselessBlocks {
	return &CheckUselessBlocks{compiler: compiler}
}

// Java: CheckUselessBlocks.process(Node, Node)
func (c *CheckUselessBlocks) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}

// IsLoneBlock: a lone block is a non-synthetic, not-added BLOCK that is a
// direct child of another non-synthetic, not-added BLOCK or a SCRIPT node.
//
// Java: CheckUselessBlocks.isLoneBlock(Node)
func (c *CheckUselessBlocks) IsLoneBlock(n rhino.Node) bool {
	parent := n.GetParent()
	if !parent.IsNil() && (parent.IsScript() ||
		(parent.IsBlock() && !parent.IsSyntheticBlock() && !parent.IsAddedBlock())) {
		return !n.IsSyntheticBlock() && !n.IsAddedBlock()
	}
	return false
}

// AllowLoneBlock removes the enclosing block of a block-scoped declaration
// from the loneBlocks stack.
//
// Java: CheckUselessBlocks.allowLoneBlock(Node)
func (c *CheckUselessBlocks) AllowLoneBlock(parent rhino.Node) {
	if len(c.loneBlocks) == 0 {
		return
	}
	if c.loneBlocks[len(c.loneBlocks)-1] == parent {
		c.popLoneBlock()
	}
}

// popLoneBlock is loneBlocks.pop() on a non-empty deque.
//
// Go-only: language-adapter — see loneBlocks.
func (c *CheckUselessBlocks) popLoneBlock() {
	last := len(c.loneBlocks) - 1
	c.loneBlocks[last] = rhino.NilNode
	c.loneBlocks = c.loneBlocks[:last]
}

// Java: CheckUselessBlocks.shouldTraverse(NodeTraversal, Node, Node)
func (c *CheckUselessBlocks) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_BLOCK:
		if c.IsLoneBlock(n) {
			c.loneBlocks = append(c.loneBlocks, n)
		}
	case rhino.Token_LET, rhino.Token_CONST:
		c.AllowLoneBlock(parent)
	case rhino.Token_CLASS:
		if NodeUtil_IsClassDeclaration(n) {
			c.AllowLoneBlock(parent)
		}
	case rhino.Token_FUNCTION:
		if NodeUtil_IsFunctionDeclaration(n) {
			c.AllowLoneBlock(parent)
		}
	default:
	}
	return true
}

// Java: CheckUselessBlocks.visit(NodeTraversal, Node, Node)
func (c *CheckUselessBlocks) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if _, ok := ast.AsBlock(n); ok && len(c.loneBlocks) != 0 && c.loneBlocks[len(c.loneBlocks)-1] == n {
		c.popLoneBlock()
		t.Report(n, CheckUselessBlocks_USELESS_BLOCK)
	}
}
