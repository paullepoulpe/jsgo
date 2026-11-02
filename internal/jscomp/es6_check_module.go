// Port of com/google/javascript/jscomp/Es6CheckModule.java @ closure-compiler v20261028
// Port of com/google/javascript/jscomp/CheckTypeImportCodeReferences.java @ closure-compiler v20261028

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
	script, ok := ast.AsScript(n)
	if !ok {
		return true
	}
	hasModule, hasImports := false, false
	for item := range script.Items() {
		if mod, ok := item.AsModuleBody(); ok {
			hasModule = true
			for stmt := range mod.Statements() {
				if _, ok := stmt.AsImport(); ok {
					hasImports = true
					break
				}
			}
		} else if _, ok := item.AsImport(); ok {
			hasModule = true
			hasImports = true
		}
	}
	e.noImportsInCurrentScript = !hasImports
	return hasModule
}

func (*Es6CheckModule) visitTokens() *ast.TokenSet { return &es6CheckModuleVisitTokens }
func (*Es6CheckModule) cloneForParallel(c *Compiler) NodeTraversal_Callback {
	return NewEs6CheckModule(c)
}

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
			!NodeUtil_IsDeclarationLValue(n) {
			if obj, ok := ast.MustNormalGet(n).Object().AsName(); ok {
				if v := t.GetScope().GetVarByNode(obj.Node()); v != nil {
					if star, ok := ast.AsImportStar(v.GetNameNode()); ok {
						// import * as M from '';
						// M.x = 2;
						e.compiler.Report(JSError_MakeWithNode(n, Es6CheckModule_IMPORT_CANNOT_BE_REASSIGNED, star.Ident()))
					}
				}
			}
		}
	case rhino.Token_NAME:
		if !e.noImportsInCurrentScript &&
			NodeUtil_IsLValue(n) && !NodeUtil_IsDeclarationLValue(n) {
			if v := t.GetScope().GetVarByNode(n); v != nil {
				if nameNode := v.GetNameNode(); !nameNode.IsNil() && nameNode != n && NodeUtil_IsImportedName(nameNode) {
					// import { x } from '';
					// x = 2;
					e.compiler.Report(JSError_MakeWithNode(n, Es6CheckModule_IMPORT_CANNOT_BE_REASSIGNED, nameNode.GetString()))
				}
			}
		}
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
var checkTypeImportCodeReferences_GOOG_REQUIRE_TYPE = rhino.QualifiedName_Of("goog.requireType")

// Java: CheckTypeImportCodeReferences(AbstractCompiler)
func NewCheckTypeImportCodeReferences(compiler AbstractCompiler) *CheckTypeImportCodeReferences {
	return &CheckTypeImportCodeReferences{compiler: compiler}
}

func googRequireTypeDeclarator(n rhino.Node) (ast.Declarator, bool) {
	decl, ok := ast.AsNameDeclaration(n)
	if !ok {
		return ast.Declarator{}, false
	}
	d, ok := decl.OnlyDeclarator()
	if !ok {
		return ast.Declarator{}, false
	}
	var init ast.Expr
	if name, isName := d.AsName(); isName {
		init, ok = name.Init()
	} else if dlhs, isDlhs := d.AsDestructuringLhs(); isDlhs {
		init, ok = dlhs.Init()
	}
	if !ok {
		return ast.Declarator{}, false
	}
	call, ok := init.AsCall()
	if !ok || !checkTypeImportCodeReferences_GOOG_REQUIRE_TYPE.Matches(call.Callee().Node()) {
		return ast.Declarator{}, false
	}
	return d, true
}

func collectRequireTypeIDs(n rhino.Node, ids *map[uint32]struct{}) {
	if d, ok := googRequireTypeDeclarator(n); ok {
		if *ids == nil {
			*ids = make(map[uint32]struct{})
		}
		if name, ok := d.AsName(); ok {
			(*ids)[name.Node().GetStringID()] = struct{}{}
		} else if objPat, ok := d.MustDestructuringLhs().Pattern().AsObjectPattern(); ok {
			for prop := range objPat.Props() {
				if sk, ok := prop.AsStringKey(); ok {
					if target, ok := sk.Value().AsName(); ok {
						(*ids)[target.Node().GetStringID()] = struct{}{}
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
	if !c.processParallel(root) {
		sc := NewSyntacticScopeCreator(c.compiler)
		b := NodeTraversal_Builder{
			compiler:     c.compiler,
			callback:     c,
			scopeCreator: sc,
		}
		b.Traverse(root)
	}
	c.requireTypeIDs = nil
}

// processParallel checks the program's scripts on workers. Go-only: the pass
// only reads the tree (each NAME of a goog.requireType'd name, looked up
// through the scopes to its declaration), so each script is checked against
// the shared global scope of the program, and the reports are replayed in
// script order.
func (c *CheckTypeImportCodeReferences) processParallel(root rhino.Node) bool {
	ids := c.requireTypeIDs
	return traverseScriptsParallel(c.compiler, root, parallelScriptOptions[*CheckTypeImportCodeReferences]{
		useSharedGlobalScope: true,
		initWorker: func(_ int, rc *Compiler) *CheckTypeImportCodeReferences {
			w := NewCheckTypeImportCodeReferences(rc)
			w.requireTypeIDs = ids
			return w
		},
	})
}

func (*CheckTypeImportCodeReferences) visitTokens() *ast.TokenSet {
	return &checkTypeImportCodeReferencesVisitTokens
}

// Java: CheckTypeImportCodeReferences.visit(NodeTraversal, Node, Node)
func (c *CheckTypeImportCodeReferences) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	name, ok := ast.AsName(n)
	if !ok {
		return
	}
	id := n.GetStringID()
	if c.requireTypeIDs != nil {
		if _, ok := c.requireTypeIDs[id]; !ok {
			return
		}
	}

	v := t.GetScope().GetVarByID(id)
	if v == nil {
		return
	}

	declNameNode := v.GetNameNode()
	if declNameNode.IsNil() || declNameNode == n || !NodeUtil_IsDeclarationLValue(declNameNode) {
		return
	}

	decl := declNameNode.GetParent()
	if _, ok := ast.AsStringKey(decl); ok {
		decl = decl.GetGrandparent().GetParent()
	}
	if _, ok := googRequireTypeDeclarator(decl); !ok {
		return
	}

	t.Report(n, CheckTypeImportCodeReferences_TYPE_IMPORT_CODE_REFERENCE, name.Ident())
}
