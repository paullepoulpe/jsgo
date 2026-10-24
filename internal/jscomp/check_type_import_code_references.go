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

	compiler         AbstractCompiler
	requireTypeNames map[string]struct{}
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

func collectRequireTypeNames(n rhino.Node, names *map[string]struct{}) {
	if NodeUtil_IsNameDeclaration(n) && n.HasOneChild() {
		first := n.GetFirstChild()
		if first.HasChildren() {
			last := first.GetLastChild()
			if last.IsCall() && last.GetFirstChild().MatchesQualifiedNameNode(checkTypeImportCodeReferences_GOOG_REQUIRE_TYPE) {
				if *names == nil {
					*names = make(map[string]struct{})
				}
				if first.IsName() {
					(*names)[first.GetString()] = struct{}{}
				} else {
					for pat := first.GetFirstChild(); !pat.IsNil(); pat = pat.GetNext() {
						for key := pat.GetFirstChild(); !key.IsNil(); key = key.GetNext() {
							if key.IsStringKey() {
								for child := key.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
									if child.IsName() {
										(*names)[child.GetString()] = struct{}{}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		collectRequireTypeNames(child, names)
	}
}

// Java: CheckTypeImportCodeReferences.process(Node, Node)
func (c *CheckTypeImportCodeReferences) Process(externs rhino.Node, root rhino.Node) {
	c.requireTypeNames = nil
	collectRequireTypeNames(root, &c.requireTypeNames)
	if len(c.requireTypeNames) == 0 {
		return
	}
	sc := NewSyntacticScopeCreator(c.compiler)
	b := NodeTraversal_Builder{
		compiler:     c.compiler,
		callback:     c,
		scopeCreator: sc,
	}
	b.Traverse(root)
	c.requireTypeNames = nil
}

// Java: CheckTypeImportCodeReferences.visit(NodeTraversal, Node, Node)
func (c *CheckTypeImportCodeReferences) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if !n.IsName() {
		return
	}
	if c.requireTypeNames != nil {
		if _, ok := c.requireTypeNames[n.GetString()]; !ok {
			return
		}
	}

	var_ := t.GetScope().GetVar(n.GetString())
	if var_ == nil {
		return
	}

	declNameNode := var_.GetNameNode()

	if declNameNode.IsNil() || declNameNode == n || !NodeUtil_IsDeclarationLValue(declNameNode) {
		return
	}

	var decl rhino.Node
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
