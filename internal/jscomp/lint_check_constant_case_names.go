// Port of com/google/javascript/jscomp/lint/CheckConstantCaseNames.java @ closure-compiler v20261015

/*
 * Copyright 2019 The Closure Compiler Authors.
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

// CheckConstantCaseNames looks for module-level variable declarations that
// use CONSTANT_CASE, according to the Google style guide, and verifies that
// they are also annotated @const or are in a const clause.
//
// This pass could be extended to check CONSTANT_CASE properties in the
// future.
//
// Non-module-level variables should always use camel case according to the
// Google style guide. In order to not confuse users, this pass does not warn
// that they should be @const. (A more correct lint check could warn that
// non-module-locals should not be constant case.)
type CheckConstantCaseNames struct {
	compiler   AbstractCompiler
	convention CodingConvention
	// Maps CONSTANT_CASE module-level names to their initializing NAME node
	invalidNamesPerModule *jutil.LinkedHashMap[string, rhino.Node]
	// Subset of variables in `invalidNamesPerModule` that are mutated post-declaration.
	reassignedNames *jutil.LinkedHashSet[string]
}

var CheckConstantCaseNames_MISSING_CONST_PROPERTY = DiagnosticType_Disabled(
	"JSC_MISSING_CONST_ON_CONSTANT_CASE",
	"CONSTANT_CASE name \"{0}\" is constant-by-convention, so must be explicitly"+
		" `const` or @const")

var CheckConstantCaseNames_REASSIGNED_CONSTANT_CASE_NAME = DiagnosticType_Disabled(
	"JSC_REASSIGNED_CONSTANT_CASE_NAME",
	"CONSTANT_CASE name \"{0}\" is constant-by-convention but is reassigned. "+
		"Use camelCase instead.")

// NewCheckConstantCaseNames is new CheckConstantCaseNames(compiler).
//
// Java: CheckConstantCaseNames(AbstractCompiler)
func NewCheckConstantCaseNames(compiler AbstractCompiler) *CheckConstantCaseNames {
	c := &CheckConstantCaseNames{
		invalidNamesPerModule: jutil.NewLinkedHashMap[string, rhino.Node](),
		reassignedNames:       jutil.NewLinkedHashSet[string](),
	}
	c.compiler = compiler
	c.convention = compiler.GetCodingConvention()
	return c
}

// Java: CheckConstantCaseNames.process(Node, Node)
func (c *CheckConstantCaseNames) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}

// Java: CheckConstantCaseNames.shouldTraverse(NodeTraversal, Node, Node)
func (c *CheckConstantCaseNames) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	// Only need to warn for module-level names, so don't visit other files.
	if n.IsScript() {
		return n.HasChildren() && n.GetFirstChild().IsModuleBody()
	}
	return true
}

// Java: CheckConstantCaseNames.visit(NodeTraversal, Node, Node)
func (c *CheckConstantCaseNames) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsModuleBody() {
		c.ReportWarningsAndClear()
		return
	}
	switch n.GetToken() {
	case rhino.Token_VAR, rhino.Token_LET:
		// Skip CONST as it automatically meets the criteria and only look for module-level vars.
		if !t.InModuleScope() {
			return
		}
		info := n.GetJSDocInfo()
		if info != nil && info.HasConstAnnotation() {
			break
		}
		NodeUtil_VisitLhsNodesInNode(
			n,
			func(name rhino.Node) {
				if c.convention.IsConstant(name.GetString()) {
					c.invalidNamesPerModule.Put(name.GetString(), name)
				}
			})
	case rhino.Token_NAME:
		if !c.invalidNamesPerModule.ContainsKey(n.GetString()) {
			return
		}
		if !NodeUtil_IsLValue(n) {
			return
		}
		// Verify this name is referring to the actual module-level var and not a local shadow.
		v := t.GetScope().GetVar(n.GetString())
		if v.GetScopeRoot().IsModuleBody() {
			c.reassignedNames.Add(n.GetString())
		}
	default:
	}
}

// Java: CheckConstantCaseNames.reportWarningsAndClear()
func (c *CheckConstantCaseNames) ReportWarningsAndClear() {
	for nameNode := range c.invalidNamesPerModule.Values() {
		name := nameNode.GetString()
		if c.reassignedNames.Contains(name) {
			c.compiler.Report(JSError_MakeWithNode(nameNode, CheckConstantCaseNames_REASSIGNED_CONSTANT_CASE_NAME, name))
		} else {
			c.compiler.Report(JSError_MakeWithNode(nameNode, CheckConstantCaseNames_MISSING_CONST_PROPERTY, name))
		}
	}
	c.invalidNamesPerModule = jutil.NewLinkedHashMap[string, rhino.Node]()
	c.reassignedNames = jutil.NewLinkedHashSet[string]()
}
