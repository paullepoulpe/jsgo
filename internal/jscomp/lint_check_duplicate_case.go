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
	"github.com/paullepoulpe/jsgo/internal/jutil"
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
func (c *CheckDuplicateCase) Process(externs *rhino.Node, root *rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}

// Java: CheckDuplicateCase.visit(NodeTraversal, Node, Node)
func (c *CheckDuplicateCase) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if n.IsSwitch() {
		switchBody := n.GetSecondChild()
		first := switchBody.GetFirstChild()
		if first == nil || first.GetNext() == nil {
			return
		}
		if checkDuplicateCaseFast(first) {
			return
		}
		cases := jutil.NewLinkedHashSet[string]()
		for curr := first; curr != nil; curr = curr.GetNext() {
			source := c.compiler.ToSourceOfNode(curr.GetFirstChild())
			if !cases.Add(source) {
				t.Report(curr, CheckDuplicateCase_DUPLICATE_CASE)
			}
		}
	}
}

// checkDuplicateCaseFast returns true if all switch cases are simple distinct
// integer or string literals (with at most one default case), avoiding
// CodePrinter invocations on the common path.
func checkDuplicateCaseFast(first *rhino.Node) bool {
	var intBuf [32]int64
	intCount := 0
	var strBuf [32]string
	strCount := 0
	hasDefault := false

	for curr := first; curr != nil; curr = curr.GetNext() {
		if curr.IsDefaultCase() {
			if hasDefault {
				return false
			}
			hasDefault = true
			continue
		}
		if !curr.IsCase() {
			return false
		}
		expr := curr.GetFirstChild()
		if expr == nil {
			return false
		}
		switch expr.GetToken() {
		case rhino.Token_NUMBER:
			d := expr.GetDouble()
			iv := int64(d)
			if float64(iv) != d || intCount >= len(intBuf) {
				return false
			}
			for i := 0; i < intCount; i++ {
				if intBuf[i] == iv {
					return false
				}
			}
			intBuf[intCount] = iv
			intCount++
		case rhino.Token_STRINGLIT:
			if strCount >= len(strBuf) {
				return false
			}
			s := expr.GetString()
			for i := 0; i < strCount; i++ {
				if strBuf[i] == s {
					return false
				}
			}
			strBuf[strCount] = s
			strCount++
		default:
			return false
		}
	}
	return true
}
