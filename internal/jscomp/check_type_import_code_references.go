// Port of com/google/javascript/jscomp/CheckTypeImportCodeReferences.java @ closure-compiler v20261015

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
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CheckTypeImportCodeReferences checks for invalid code references to
// type-only imports (i.e., goog.requireType).
type CheckTypeImportCodeReferences struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	compiler AbstractCompiler
}

// Java: CheckTypeImportCodeReferences.TYPE_IMPORT_CODE_REFERENCE
var CheckTypeImportCodeReferences_TYPE_IMPORT_CODE_REFERENCE = DiagnosticType_Error(
	"JSC_TYPE_IMPORT_CODE_REFERENCE",
	"Cannot reference goog.requireType()''d name {0} outside of a type annotation.")

// Java: CheckTypeImportCodeReferences.GOOG_REQUIRE_TYPE
var checkTypeImportCodeReferences_GOOG_REQUIRE_TYPE = rhino.IR_Getprop(rhino.IR_Name("goog"), "requireType")

// Java: CheckTypeImportCodeReferences(AbstractCompiler)
func NewCheckTypeImportCodeReferences(compiler AbstractCompiler) *CheckTypeImportCodeReferences {
	c := &CheckTypeImportCodeReferences{}
	c.compiler = compiler
	return c
}

// Java: CheckTypeImportCodeReferences.process(Node, Node)
func (c *CheckTypeImportCodeReferences) Process(externs *rhino.Node, root *rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}

// Java: CheckTypeImportCodeReferences.visit(NodeTraversal, Node, Node)
func (c *CheckTypeImportCodeReferences) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if !n.IsName() {
		return
	}

	var_ := t.GetScope().GetVar(n.GetString())
	if var_ == nil {
		return
	}

	declNameNode := var_.GetNameNode()

	if declNameNode == nil || declNameNode == n || !NodeUtil_IsDeclarationLValue(declNameNode) {
		return
	}

	var decl *rhino.Node
	if declNameNode.GetParent().IsStringKey() {
		decl = declNameNode.GetGrandparent().GetGrandparent()
	} else {
		decl = declNameNode.GetParent()
	}
	if !NodeUtil_IsNameDeclaration(decl) ||
		!decl.HasOneChild() ||
		!decl.GetFirstChild().HasChildren() ||
		!decl.GetFirstChild().GetLastChild().IsCall() {
		return
	}

	callNode := decl.GetFirstChild().GetLastChild()
	if !callNode.GetFirstChild().MatchesQualifiedNameNode(checkTypeImportCodeReferences_GOOG_REQUIRE_TYPE) {
		return
	}

	t.Report(n, CheckTypeImportCodeReferences_TYPE_IMPORT_CODE_REFERENCE, n.GetString())
}
