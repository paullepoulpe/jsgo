// Port of com/google/javascript/jscomp/Es6CheckModule.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/CheckTypeImportCodeReferences.java @ closure-compiler v20261021

/*
 * Copyright 2017 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

var (
	es6CheckModuleVisitTokens = func() ast.TokenSet {
		var s ast.TokenSet
		s.Add(rhino.Token_SCRIPT)
		s.Add(rhino.Token_THIS)
		s.Add(rhino.Token_GETPROP)
		s.Add(rhino.Token_GETELEM)
		s.Add(rhino.Token_NAME)
		return s
	}()
	checkTypeImportCodeReferencesVisitTokens = func() ast.TokenSet {
		var s ast.TokenSet
		s.Add(rhino.Token_NAME)
		return s
	}()
)

// Es6CheckModule checks that ES6 Modules are used correctly, and do not
// reference undefined keywords or features.
type Es6CheckModule struct {
	compiler                 AbstractCompiler
	noImportsInCurrentScript bool
}

var Es6CheckModule_ES6_MODULE_REFERENCES_THIS = DiagnosticType_Warning(
	"ES6_MODULE_REFERENCES_THIS", "The body of an ES6 module cannot reference 'this'.")

var Es6CheckModule_IMPORT_CANNOT_BE_REASSIGNED = DiagnosticType_Error(
	"JSC_IMPORT_CANNOT_BE_REASSIGNED", "Assignment to constant variable \"{0}\".")

// Java: Es6CheckModule(AbstractCompiler)
func NewEs6CheckModule(compiler AbstractCompiler) *Es6CheckModule {
	e := &Es6CheckModule{}
	e.compiler = compiler
	return e
}

// Java: Es6CheckModule.process(Node, Node)
func (e *Es6CheckModule) Process(externs rhino.Node, root rhino.Node) {
	sc := NewSyntacticScopeCreator(e.compiler)
	b := NodeTraversal_Builder{
		compiler:     e.compiler,
		callback:     e,
		scopeCreator: sc,
	}
	b.Traverse(root)
}

func (e *Es6CheckModule) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if n.IsScript() {
		hasModule := false
		hasImports := false
		for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			if c.IsModuleBody() {
				hasModule = true
				for mc := c.GetFirstChild(); !mc.IsNil(); mc = mc.GetNext() {
					if mc.IsImport() {
						hasImports = true
						break
					}
				}
			} else if c.IsImport() {
				hasModule = true
				hasImports = true
			}
		}
		e.noImportsInCurrentScript = !hasImports
		return hasModule
	}
	return true
}

func (*Es6CheckModule) visitTokens() *ast.TokenSet        { return &es6CheckModuleVisitTokens }
func (*Es6CheckModule) combinedVisitTokens() ast.TokenSet { return es6CheckModuleVisitTokens }

// Java: Es6CheckModule.visit(NodeTraversal, Node, Node)
func (e *Es6CheckModule) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		e.noImportsInCurrentScript = false
	case rhino.Token_THIS:
		if t.InModuleHoistScope() {
			t.Report(n, Es6CheckModule_ES6_MODULE_REFERENCES_THIS)
		}
	case rhino.Token_GETPROP, rhino.Token_GETELEM:
		if !e.noImportsInCurrentScript &&
			NodeUtil_IsLValue(n) &&
			!NodeUtil_IsDeclarationLValue(n) &&
			n.GetFirstChild().IsName() {
			var_ := t.GetScope().GetVarByNode(n.GetFirstChild())
			if var_ != nil {
				nameNode := var_.GetNameNode()
				if !nameNode.IsNil() && nameNode.IsImportStar() {
					// import * as M from '';
					// M.x = 2;
					e.compiler.Report(JSError_MakeWithNode(n, Es6CheckModule_IMPORT_CANNOT_BE_REASSIGNED, nameNode.GetString()))
				}
			}
		}
	case rhino.Token_NAME:
		if !e.noImportsInCurrentScript &&
			NodeUtil_IsLValue(n) && !NodeUtil_IsDeclarationLValue(n) {
			var_ := t.GetScope().GetVarByNode(n)
			if var_ != nil {
				nameNode := var_.GetNameNode()
				if !nameNode.IsNil() && nameNode != n {
					if NodeUtil_IsImportedName(nameNode) {
						// import { x } from '';
						// x = 2;
						e.compiler.Report(JSError_MakeWithNode(n, Es6CheckModule_IMPORT_CANNOT_BE_REASSIGNED, nameNode.GetString()))
					}
				}
			}
		}
	default:
	}
}

// ---------------------------------------------------------------------------
// CheckTypeImportCodeReferences
// ---------------------------------------------------------------------------

// CheckTypeImportCodeReferences checks for invalid code references to
// type-only imports (i.e., goog.requireType).
type CheckTypeImportCodeReferences struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	compiler       AbstractCompiler
	requireTypeIDs map[uint32]struct{}
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

func collectRequireTypeIDs(n rhino.Node, ids *map[uint32]struct{}) {
	if NodeUtil_IsNameDeclaration(n) && n.HasOneChild() {
		first := n.GetFirstChild()
		if first.HasChildren() {
			last := first.GetLastChild()
			if last.IsCall() && last.GetFirstChild().MatchesQualifiedNameNode(checkTypeImportCodeReferences_GOOG_REQUIRE_TYPE) {
				if *ids == nil {
					*ids = make(map[uint32]struct{})
				}
				if first.IsName() {
					(*ids)[first.GetStringID()] = struct{}{}
				} else {
					for pat := first.GetFirstChild(); !pat.IsNil(); pat = pat.GetNext() {
						for key := pat.GetFirstChild(); !key.IsNil(); key = key.GetNext() {
							if key.IsStringKey() {
								for child := key.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
									if child.IsName() {
										(*ids)[child.GetStringID()] = struct{}{}
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
		collectRequireTypeIDs(child, ids)
	}
}

// Java: CheckTypeImportCodeReferences.process(Node, Node)
func (c *CheckTypeImportCodeReferences) Process(externs rhino.Node, root rhino.Node) {
	c.requireTypeIDs = nil
	collectRequireTypeIDs(root, &c.requireTypeIDs)
	if len(c.requireTypeIDs) == 0 {
		return
	}
	sc := NewSyntacticScopeCreator(c.compiler)
	b := NodeTraversal_Builder{
		compiler:     c.compiler,
		callback:     c,
		scopeCreator: sc,
	}
	b.Traverse(root)
	c.requireTypeIDs = nil
}

func (*CheckTypeImportCodeReferences) visitTokens() *ast.TokenSet {
	return &checkTypeImportCodeReferencesVisitTokens
}
func (*CheckTypeImportCodeReferences) combinedVisitTokens() ast.TokenSet {
	return checkTypeImportCodeReferencesVisitTokens
}

// Java: CheckTypeImportCodeReferences.visit(NodeTraversal, Node, Node)
func (c *CheckTypeImportCodeReferences) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if !n.IsName() {
		return
	}
	id := n.GetStringID()
	if c.requireTypeIDs != nil {
		if _, ok := c.requireTypeIDs[id]; !ok {
			return
		}
	}

	var_ := t.GetScope().GetVarByID(id)
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
