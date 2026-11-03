// Port of com/google/javascript/jscomp/ClosureCheckModule.java @ closure-compiler v20261028

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
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: this pass uses the typed AST view (internal/ast) to read the
// calls and declarations it checks: the first argument of goog.module and
// goog.require (ast.Call.FirstArg), the target of an export assignment, and
// the members of a destructuring import (`const {a: b} = goog.require(...)`
// is an object pattern of string keys whose values are names). The
// position checks (is this the first statement, what is the require's
// parent) read parents and siblings of any kind and stay on rhino.

// ClosureCheckModule checks that goog.module() is used correctly.
//
// Note that this file only does checks that can be done per-file. Whole
// program checks happen during goog.module rewriting, in
// ClosureRewriteModule.
type ClosureCheckModule struct {
	moduleTraversal

	currentModuleInfo *closureModuleCheckInfo
}

var _ NodeTraversal_Callback = (*ClosureCheckModule)(nil)

var ClosureCheckModule_AT_EXPORT_IN_GOOG_MODULE = DiagnosticType_Error(
	"JSC_AT_EXPORT_IN_GOOG_MODULE",
	"@export has no effect on top-level names in a goog.module."+
		" Consider using goog.exportSymbol instead.")

var ClosureCheckModule_AT_EXPORT_IN_NON_LEGACY_GOOG_MODULE = DiagnosticType_Error(
	"JSC_AT_EXPORT_IN_NON_LEGACY_GOOG_MODULE",
	"@export is not allowed here in a non-legacy goog.module."+
		" Consider using goog.exportSymbol instead.")

var ClosureCheckModule_GOOG_MODULE_IN_NON_MODULE = DiagnosticType_Error(
	"JSC_GOOG_MODULE_IN_NON_MODULE",
	"goog.module() call must be the first statement in a module.")

var ClosureCheckModule_GOOG_MODULE_MISPLACED = DiagnosticType_Error(
	"JSC_GOOG_MODULE_MISPLACED", "goog.module() call must be the first statement in a file.")

var ClosureCheckModule_DECLARE_LEGACY_NAMESPACE_IN_NON_MODULE = DiagnosticType_Error(
	"JSC_DECLARE_LEGACY_NAMESPACE_IN_NON_MODULE",
	"goog.module.declareLegacyNamespace may only be called in a goog.module.")

var ClosureCheckModule_GOOG_MODULE_REFERENCES_THIS = DiagnosticType_Error(
	"JSC_GOOG_MODULE_REFERENCES_THIS", "The body of a goog.module cannot reference 'this'.")

var ClosureCheckModule_GOOG_MODULE_USES_THROW = DiagnosticType_Error(
	"JSC_GOOG_MODULE_USES_THROW", "The body of a goog.module cannot use 'throw'.")

var ClosureCheckModule_DUPLICATE_NAME_SHORT_REQUIRE = DiagnosticType_Error(
	"JSC_DUPLICATE_NAME_SHORT_REQUIRE",
	"Found multiple goog.require statements importing identifier ''{0}''.")

var ClosureCheckModule_INVALID_DESTRUCTURING_REQUIRE = DiagnosticType_Error(
	"JSC_INVALID_DESTRUCTURING_REQUIRE",
	"Destructuring goog.require must be a simple object pattern.")

var ClosureCheckModule_LET_GOOG_REQUIRE = DiagnosticType_Disabled(
	"JSC_LET_GOOG_REQUIRE",
	"Module imports must be constant. Please use ''const'' instead of ''let''.")

var ClosureCheckModule_MULTIPLE_MODULES_IN_FILE = DiagnosticType_Error(
	"JSC_MULTIPLE_MODULES_IN_FILE",
	"There should only be a single goog.module() statement per file.")

var ClosureCheckModule_ONE_REQUIRE_PER_DECLARATION = DiagnosticType_Error(
	"JSC_ONE_REQUIRE_PER_DECLARATION",
	"There may only be one goog.require() per var/let/const declaration.")

var ClosureCheckModule_INCORRECT_SHORTNAME_CAPITALIZATION = DiagnosticType_Disabled(
	"JSC_INCORRECT_SHORTNAME_CAPITALIZATION",
	"The capitalization of short name {0} is incorrect; it should be {1}.")

var ClosureCheckModule_EXPORT_NOT_AT_MODULE_SCOPE = DiagnosticType_Error(
	"JSC_EXPORT_NOT_AT_MODULE_SCOPE", "Exports must be at the top-level of a module")

var ClosureCheckModule_EXPORT_NOT_A_STATEMENT = DiagnosticType_Error("JSC_EXPORT_NOT_A_STATEMENT", "Exports should be a statement.")

var ClosureCheckModule_EXPORT_REPEATED_ERROR = DiagnosticType_Error(
	"JSC_EXPORT_REPEATED_ERROR",
	"Name cannot be exported multiple times. Previous export on line {0}.")

var ClosureCheckModule_REFERENCE_TO_MODULE_GLOBAL_NAME = DiagnosticType_Error(
	"JSC_REFERENCE_TO_MODULE_GLOBAL_NAME",
	"References to the global name of a module are not allowed. Perhaps you meant exports?")

var ClosureCheckModule_REFERENCE_TO_FULLY_QUALIFIED_IMPORT_NAME = DiagnosticType_Disabled(
	"JSC_REFERENCE_TO_FULLY_QUALIFIED_IMPORT_NAME",
	"Reference to fully qualified import name ''{0}''."+
		" Imports in goog.module should use the return value of"+
		" goog.require / goog.forwardDeclare instead.")

var ClosureCheckModule_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME = DiagnosticType_Disabled(
	"JSC_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME",
	"Reference to fully qualified import name ''{0}''."+
		" Please use the short name ''{1}'' instead.")

var ClosureCheckModule_USE_OF_GOOG_PROVIDE = DiagnosticType_Disabled(
	"JSC_USE_OF_GOOG_PROVIDE",
	"goog.provide is deprecated in favor of goog.module.")

var ClosureCheckModule_REQUIRE_NOT_AT_TOP_LEVEL = DiagnosticType_Error(
	"JSC_REQUIRE_NOT_AT_TOP_LEVEL", "goog.require() must be called at file scope.")

var ClosureCheckModule_LEGACY_NAMESPACE_NOT_AT_TOP_LEVEL = DiagnosticType_Error(
	"JSC_LEGACY_NAMESPACE_NOT_AT_TOP_LEVEL",
	"goog.module.declareLegacyNamespace() does not return a value")

var ClosureCheckModule_LEGACY_NAMESPACE_NOT_AFTER_GOOG_MODULE = DiagnosticType_Error(
	"JSC_LEGACY_NAMESPACE_NOT_AT_TOP_LEVEL",
	"goog.module.declareLegacyNamespace() must be immediately after the"+
		" goog.module('...'); call")

var ClosureCheckModule_LEGACY_NAMESPACE_ARGUMENT = DiagnosticType_Error(
	"JSC_LEGACY_NAMESPACE_ARGUMENT",
	"goog.module.declareLegacyNamespace() takes no arguments")

var closureCheckModule_GOOG_MODULE = rhino.QualifiedName_Of("goog.module")
var closureCheckModule_GOOG_PROVIDE = rhino.QualifiedName_Of("goog.provide")
var closureCheckModule_GOOG_REQUIRE = rhino.QualifiedName_Of("goog.require")
var closureCheckModule_GOOG_REQUIRE_TYPE = rhino.QualifiedName_Of("goog.requireType")
var closureCheckModule_GOOG_MODULE_GET = rhino.QualifiedName_Of("goog.module.get")
var closureCheckModule_GOOG_FORWARD_DECLARE = rhino.QualifiedName_Of("goog.forwardDeclare")
var closureCheckModule_GOOG_MODULE_DECLARE_LEGACY_NAMESPACE = rhino.QualifiedName_Of("goog.module.declareLegacyNamespace")

// closureModuleCheckInfo is ClosureCheckModule.ModuleInfo.
type closureModuleCheckInfo struct {
	// Name of the module in question (i.e. the argument to goog.module)
	name string
	// Mapping from fully qualified goog.required names to the import LHS
	// node. For standalone goog.require()s the value is the EXPR_RESULT
	// node.
	//
	// Go-only: language-adapter — the Java map may receive a null key (the
	// first argument of a malformed goog.require call); it is only ever read
	// with non-null keys, so null keys are not stored.
	importsByLongRequiredName map[string]rhino.Node
	// Module-local short names for goog.required symbols.
	shortImportNames map[string]struct{}
	// A map from the export names "exports.name" to the nodes of those named
	// exports. The default export is keyed with just "exports"
	exportNodesByName map[string]rhino.Node
}

// Java: ClosureCheckModule.ModuleInfo(String)
func newClosureModuleCheckInfo(moduleName string) *closureModuleCheckInfo {
	return &closureModuleCheckInfo{
		name:                      moduleName,
		importsByLongRequiredName: make(map[string]rhino.Node),
		shortImportNames:          make(map[string]struct{}),
		exportNodesByName:         make(map[string]rhino.Node),
	}
}

// Java: ClosureCheckModule(AbstractCompiler, ModuleMetadataMap)
func NewClosureCheckModule(compiler AbstractCompiler, moduleMetadataMap *ModuleMetadataMap) *ClosureCheckModule {
	c := &ClosureCheckModule{}
	c.initModuleTraversal(compiler, moduleMetadataMap)
	return c
}

// Java: ClosureCheckModule.process(Node, Node)
func (c *ClosureCheckModule) Process(externs rhino.Node, root rhino.Node) {
	if c.processParallel(root) {
		return
	}
	NodeTraversal_Traverse(c.compiler, root, c)
}

func (c *ClosureCheckModule) processParallel(root rhino.Node) bool {
	return traverseScriptsParallel(c.compiler, root, parallelScriptOptions[*ClosureCheckModule]{
		useSharedGlobalScope: true,
		initWorker: func(_ int, rc *Compiler) *ClosureCheckModule {
			return NewClosureCheckModule(rc, c.moduleMetadataMap)
		},
	})
}

var closureCheckModuleVisitTokens = makeVisitTokenSet(
	rhino.Token_SCRIPT,
	rhino.Token_BLOCK,
	rhino.Token_CALL,
	rhino.Token_EXPR_RESULT,
	rhino.Token_ASSIGN,
	rhino.Token_CLASS,
	rhino.Token_FUNCTION,
	rhino.Token_VAR,
	rhino.Token_LET,
	rhino.Token_CONST,
	rhino.Token_THIS,
	rhino.Token_THROW,
	rhino.Token_NAME,
	rhino.Token_GETPROP,
)

var (
	_ NodeTraversal_Callback    = (*ClosureCheckModule)(nil)
	_ nodeTraversal_visitTokens = (*ClosureCheckModule)(nil)
	_ nodeTraversal_visitJSDoc  = (*ClosureCheckModule)(nil)
)

func (*ClosureCheckModule) visitTokens() *ast.TokenSet { return &closureCheckModuleVisitTokens }

func (*ClosureCheckModule) visitJSDoc() bool { return true }

// Java: NodeTraversal.AbstractModuleCallback.shouldTraverse(NodeTraversal, Node, Node)
func (c *ClosureCheckModule) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if c.enterNode(t, n) {
		c.enterModule(c.currentModule, c.scopeRoot)
	}
	return true
}

// Java: NodeTraversal.AbstractModuleCallback.visit(NodeTraversal, Node, Node)
func (c *ClosureCheckModule) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if exitedModule, exitedScopeRoot := c.exitNode(t, n); exitedModule != nil {
		c.exitModule(exitedModule, exitedScopeRoot)
	}
	c.visitWithModule(t, n, parent, c.currentModule, c.scopeRoot)
}

// Java: ClosureCheckModule.enterModule(ModuleMetadata, Node)
func (c *ClosureCheckModule) enterModule(currentModule *ModuleMetadataMap_ModuleMetadata, moduleScopeRoot rhino.Node) {
	if !currentModule.IsGoogModule() {
		return
	}

	jutil.CheckState(c.currentModuleInfo == nil)
	jutil.CheckState(!currentModule.GoogNamespaces().IsEmpty())
	first := ""
	for ns := range currentModule.GoogNamespaces().All() {
		first = ns
		break
	}
	c.currentModuleInfo = newClosureModuleCheckInfo(first)
}

// Java: ClosureCheckModule.exitModule(ModuleMetadata, Node)
func (c *ClosureCheckModule) exitModule(currentModule *ModuleMetadataMap_ModuleMetadata, moduleScopeRoot rhino.Node) {
	c.currentModuleInfo = nil
}

// Java: ClosureCheckModule.visit(NodeTraversal, Node, ModuleMetadata, Node)
func (c *ClosureCheckModule) visitWithModule(
	t *NodeTraversal,
	n rhino.Node,
	parent rhino.Node,
	currentModule *ModuleMetadataMap_ModuleMetadata,
	moduleScopeRoot rhino.Node) {
	if c.currentModuleInfo == nil {
		if !n.IsCall() && !n.IsExprResult() {
			return
		}
		if NodeUtil_IsCallToQualifiedName(n, closureCheckModule_GOOG_MODULE) {
			t.Report(n, ClosureCheckModule_GOOG_MODULE_IN_NON_MODULE)
		} else if NodeUtil_IsGoogModuleDeclareLegacyNamespaceCall(n) {
			t.Report(n, ClosureCheckModule_DECLARE_LEGACY_NAMESPACE_IN_NON_MODULE)
		} else if NodeUtil_IsCallToQualifiedName(n, closureCheckModule_GOOG_PROVIDE) {
			// This error is reported here, rather than in a provide-specific pass, because it
			// must be reported prior to ClosureRewriteModule converting legacy goog.modules into
			// goog.provides.
			t.Report(n, ClosureCheckModule_USE_OF_GOOG_PROVIDE)
		}
		return
	}
	if jsDoc := n.GetJSDocInfo(); jsDoc != nil {
		c.checkJSDoc(t, jsDoc)
	}
	switch n.GetToken() {
	case rhino.Token_CALL:
		call := ast.MustCall(n)
		callee := call.Callee().Node()
		if closureCheckModule_GOOG_MODULE.Matches(callee) {
			firstArgumentName, hasFirstArg := extractFirstStringArg(call)
			if !hasFirstArg || c.currentModuleInfo.name != firstArgumentName {
				t.Report(n, ClosureCheckModule_MULTIPLE_MODULES_IN_FILE)
			} else if !isFirstExpressionInGoogModule(parent) {
				t.Report(n, ClosureCheckModule_GOOG_MODULE_MISPLACED)
			}
		} else if closureCheckModule_GOOG_REQUIRE.Matches(callee) ||
			closureCheckModule_GOOG_REQUIRE_TYPE.Matches(callee) ||
			closureCheckModule_GOOG_FORWARD_DECLARE.Matches(callee) {
			c.checkRequireCall(t, call, parent)
		} else if closureCheckModule_GOOG_MODULE_GET.Matches(callee) && t.InModuleHoistScope() {
			t.Report(n, ClosurePrimitiveErrors_MODULE_USES_GOOG_MODULE_GET)
		} else if closureCheckModule_GOOG_MODULE_DECLARE_LEGACY_NAMESPACE.Matches(callee) {
			checkLegacyNamespaceCall(t, call, parent)
		}
	case rhino.Token_ASSIGN:
		assign := ast.MustAssign(n)
		if isClosureExportLhs(assign.Target()) {
			c.checkModuleExport(t, assign, parent)
		}
	case rhino.Token_CLASS, rhino.Token_FUNCTION:
		if !NodeUtil_IsStatement(n) {
			break
		}
		fallthrough
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		if t.InModuleHoistScope() &&
			(n.IsClass() || NodeUtil_GetEnclosingClass(n).IsNil()) &&
			NodeUtil_GetEnclosingType(n, rhino.Token_OBJECTLIT).IsNil() {
			if jsdoc := NodeUtil_GetBestJSDocInfo(n); jsdoc != nil && jsdoc.IsExport() {
				t.Report(n, ClosureCheckModule_AT_EXPORT_IN_GOOG_MODULE)
			}
		}
	case rhino.Token_THIS:
		if t.InModuleHoistScope() {
			t.Report(n, ClosureCheckModule_GOOG_MODULE_REFERENCES_THIS)
		}
	case rhino.Token_THROW:
		if t.InModuleHoistScope() {
			t.Report(n, ClosureCheckModule_GOOG_MODULE_USES_THROW)
		}
	case rhino.Token_NAME, rhino.Token_GETPROP:
		if n.MatchesQualifiedName(c.currentModuleInfo.name) {
			if n.IsGetProp() {
				// This warning makes sense on NAME nodes as well,
				// but it would require a cleanup to land.
				t.Report(n, ClosureCheckModule_REFERENCE_TO_MODULE_GLOBAL_NAME)
			}
		} else if !parent.IsGetProp() && len(c.currentModuleInfo.importsByLongRequiredName) != 0 {
			if qname, ok := n.GetQualifiedNameString(); ok {
				node := n
				var nextPart string
				hasNextPart := false

				for c.checkImproperReferenceToImport(t, node, qname, nextPart, hasNextPart) {
					prop, ok := ast.AsGetProp(node)
					if !ok {
						// Don't get the value of a declaration (i.e. `var name = 2`)
						return
					}
					nextPart = prop.Property()
					hasNextPart = true
					qname = qname[:len(qname)-1-len(nextPart)]
					node = prop.Object().Node()
				}
			}
		}
	default:
	}
}

// Java: ClosureCheckModule.checkJSDoc(NodeTraversal, JSDocInfo)
func (c *ClosureCheckModule) checkJSDoc(t *NodeTraversal, jsDoc *rhino.JSDocInfo) {
	for it := jsDoc.TypeExprIter(); it.Next(); {
		c.checkTypeExpression(t, it.Node())
	}
}

// Java: ClosureCheckModule.checkTypeExpression(NodeTraversal, Node)
func (c *ClosureCheckModule) checkTypeExpression(t *NodeTraversal, typeNode rhino.Node) {
	NodeUtil_VisitPreOrder(
		typeNode,
		func(node rhino.Node) {
			if !node.IsStringLit() {
				return
			}

			qname := node.GetString()
			var nextQnamePart string
			hasNextQnamePart := false

			for c.checkImproperReferenceToImport(t, node, qname, nextQnamePart, hasNextQnamePart) {
				lastDot := strings.LastIndexByte(qname, '.')
				if lastDot < 0 {
					return
				}
				nextQnamePart = qname[lastDot+1:]
				hasNextQnamePart = true
				qname = qname[:lastDot]
			}
		})
}

// Java: ClosureCheckModule.isLocalVar(Var)
func isClosureLocalVar(v *Var) bool {
	return v != nil && v.IsLocal()
}

// checkImproperReferenceToImport returns true if the reference might be
// valid, false means an error was reported.
//
// Java: ClosureCheckModule.checkImproperReferenceToImport(NodeTraversal, Node, String, String)
func (c *ClosureCheckModule) checkImproperReferenceToImport(
	t *NodeTraversal, n rhino.Node, qname string, nextQnamePart string, hasNextQnamePart bool) bool {
	importLhs := c.currentModuleInfo.importsByLongRequiredName[qname]
	if importLhs.IsNil() {
		return true
	}

	if !strings.Contains(qname, ".") && isClosureLocalVar(t.GetScope().GetVar(qname)) {
		return true
	}

	if nameNode, ok := ast.AsName(importLhs); ok {
		c.compiler.Report(
			JSError_MakeWithNode(
				n,
				ClosureCheckModule_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME,
				qname,
				nameNode.Ident()))
		return false
	} else if lhs, ok := ast.AsDestructuringLhs(importLhs); ok && hasNextQnamePart {
		if objPattern, ok := lhs.Pattern().AsObjectPattern(); ok {
			for prop := range objPattern.Props() {
				// const {foo: barFoo} = goog.require('ns.bar');
				// Should use the short name "barFoo" instead of "ns.bar.foo".
				key, ok := prop.AsStringKey()
				if !ok || key.Key() != nextQnamePart {
					continue
				}
				shortName, ok := key.Value().AsName()
				if !ok {
					continue
				}
				parent := n.GetParent()
				reportNode := n
				if parent.IsGetProp() {
					reportNode = parent
				}
				c.compiler.Report(
					JSError_MakeWithNode(
						reportNode,
						ClosureCheckModule_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME,
						qname+"."+nextQnamePart,
						shortName.Ident()))
				return false
			}
		}
	}
	if exprStmt, ok := ast.AsExprStmt(importLhs); ok && NodeUtil_IsGoogForwardDeclareCall(exprStmt.Expr().Node()) {
		// goog.forwardDeclare does not need to be aliased.
		// This is because goog.forwardDeclares that are aliased have stricter semantics than those
		// that are not aliased: aliased goog.forwardDeclares always need to reference a namespace
		// that's actually defined, while non-aliased forwardDeclares may reference namespaces that
		// do not actually exist.
		return true
	}

	c.compiler.Report(JSError_MakeWithNode(n, ClosureCheckModule_REFERENCE_TO_FULLY_QUALIFIED_IMPORT_NAME, qname))
	return false
}

// isClosureExportLhs reports whether target is the LHS of a goog.module
// export, i.e. either "exports" or "exports.name".
//
// Java: ClosureCheckModule.isExportLhs(Node)
func isClosureExportLhs(target ast.Target) bool {
	if target.Node().MatchesName("exports") {
		return true
	}
	prop, ok := target.AsGetProp()
	return ok && prop.Object().Node().MatchesName("exports")
}

// Java: ClosureCheckModule.checkModuleExport(NodeTraversal, Node, Node)
func (c *ClosureCheckModule) checkModuleExport(t *NodeTraversal, assign ast.Assign, parent rhino.Node) {
	n := assign.Node()
	lhs := assign.Target().Node()
	jutil.CheckState(isClosureExportLhs(assign.Target()))
	lhsQName, _ := lhs.GetQualifiedNameString()
	// Check multiple exports of the same name
	if previousDefinition := c.currentModuleInfo.exportNodesByName[lhsQName]; !previousDefinition.IsNil() {
		previousLine := previousDefinition.GetLineno()
		t.Report(n, ClosureCheckModule_EXPORT_REPEATED_ERROR, strconv.Itoa(previousLine))
	}
	// Check exports in invalid program position
	defaultExportNode := c.currentModuleInfo.exportNodesByName["exports"]
	// If we have never seen an `exports =` default export assignment, or this is the
	// default export, then treat this assignment as an export and do the checks it is well formed.
	if defaultExportNode.IsNil() || lhs.MatchesName("exports") {
		c.currentModuleInfo.exportNodesByName[lhsQName] = lhs
		if !t.InModuleScope() {
			t.Report(n, ClosureCheckModule_EXPORT_NOT_AT_MODULE_SCOPE)
		} else if !parent.IsExprResult() {
			t.Report(n, ClosureCheckModule_EXPORT_NOT_A_STATEMENT)
		}
	}
	// Check @export on a module local name
	if !NodeUtil_IsPrototypeProperty(lhs) &&
		!NodeUtil_IsLegacyGoogModuleFile(NodeUtil_GetEnclosingScript(n)) {
		if jsDoc := n.GetJSDocInfo(); jsDoc != nil && jsDoc.IsExport() {
			t.Report(n, ClosureCheckModule_AT_EXPORT_IN_NON_LEGACY_GOOG_MODULE)
		}
	}
}

// Java: ClosureCheckModule.extractFirstArgumentName(Node)
func extractFirstStringArg(call ast.Call) (string, bool) {
	if firstArg, ok := call.FirstArg(); ok {
		if name, ok := firstArg.AsStringLit(); ok {
			return name.Value(), true
		}
	}
	return "", false
}

// Java: ClosureCheckModule.checkRequireCall(NodeTraversal, Node, Node)
func (c *ClosureCheckModule) checkRequireCall(t *NodeTraversal, call ast.Call, parent rhino.Node) {
	namespace, ok := call.LastArg()
	jutil.CheckState(ok && namespace.Token() == rhino.Token_STRINGLIT)
	nsValue := ast.MustStringLit(namespace.Node()).Value()
	switch parent.GetToken() {
	case rhino.Token_EXPR_RESULT:
		if key, ok := extractFirstStringArg(call); ok {
			if _, exists := c.currentModuleInfo.importsByLongRequiredName[key]; !exists {
				c.currentModuleInfo.importsByLongRequiredName[key] = parent
			}
		}
	case rhino.Token_NAME, rhino.Token_DESTRUCTURING_LHS:
		c.checkShortGoogRequireCall(t, call, nsValue, parent.GetParent())
	case rhino.Token_AWAIT:
		grandParent := parent.GetParent()
		if grandParent.IsDestructuringLhs() || grandParent.IsName() {
			c.checkShortGoogRequireCall(t, call, nsValue, grandParent.GetParent())
		}
	default:
		t.Report(call.Node(), ClosureCheckModule_REQUIRE_NOT_AT_TOP_LEVEL)
	}
}

// Java: ClosureCheckModule.checkShortGoogRequireCall(NodeTraversal, Node, Node)
func (c *ClosureCheckModule) checkShortGoogRequireCall(t *NodeTraversal, call ast.Call, namespace string, declaration rhino.Node) {
	isForwardDeclare := closureCheckModule_GOOG_FORWARD_DECLARE.Matches(call.Callee().Node())
	if declaration.IsLet() && !isForwardDeclare {
		t.Report(declaration, ClosureCheckModule_LET_GOOG_REQUIRE)
	}
	decl, isDecl := ast.AsNameDeclaration(declaration)
	if !isDecl {
		t.Report(declaration, ClosureCheckModule_ONE_REQUIRE_PER_DECLARATION)
		return
	}
	lhsDecl, hasOne := decl.OnlyDeclarator()
	if !hasOne {
		t.Report(declaration, ClosureCheckModule_ONE_REQUIRE_PER_DECLARATION)
		return
	}
	lhs := lhsDecl.Node()
	if destructuring, ok := lhsDecl.AsDestructuringLhs(); ok {
		if !ClosureCheckModule_IsValidDestructuringImport(destructuring) {
			t.Report(declaration, ClosureCheckModule_INVALID_DESTRUCTURING_REQUIRE)
		}
		if isForwardDeclare {
			t.Report(lhs, ClosurePrimitiveErrors_INVALID_DESTRUCTURING_FORWARD_DECLARE)
		}
	} else if nameNode, ok := lhsDecl.AsName(); ok {
		checkClosureShortName(t, nameNode, namespace)
	}
	if key, ok := extractFirstStringArg(call); ok {
		c.currentModuleInfo.importsByLongRequiredName[key] = lhs
	}
	NodeUtil_VisitLhsNodesInNode(
		declaration,
		func(nameNode rhino.Node) {
			name := ast.MustName(nameNode).Ident()
			if _, exists := c.currentModuleInfo.shortImportNames[name]; exists {
				t.Report(nameNode, ClosureCheckModule_DUPLICATE_NAME_SHORT_REQUIRE, name)
			} else {
				c.currentModuleInfo.shortImportNames[name] = struct{}{}
			}
		})
}

// closureCheckModule_asciiIsUpperCase is Guava Ascii.isUpperCase(char).
//
// Go-only: language-adapter — Guava Ascii.
func closureCheckModule_asciiIsUpperCase(c uint16) bool {
	return c >= 'A' && c <= 'Z'
}

// closureCheckModule_asciiToLowerCase is Guava Ascii.toLowerCase(char).
//
// Go-only: language-adapter — Guava Ascii.
func closureCheckModule_asciiToLowerCase(c uint16) uint16 {
	if closureCheckModule_asciiIsUpperCase(c) {
		return c ^ 0x20
	}
	return c
}

// closureCheckModule_asciiToUpperCase is Guava Ascii.toUpperCase(char).
//
// Go-only: language-adapter — Guava Ascii.
func closureCheckModule_asciiToUpperCase(c uint16) uint16 {
	if c >= 'a' && c <= 'z' {
		return c ^ 0x20
	}
	return c
}

// Java: ClosureCheckModule.checkShortName(NodeTraversal, Node, String)
func checkClosureShortName(t *NodeTraversal, shortName ast.Name, namespace string) {
	nextQnamePart := shortName.Ident()
	if strings.HasPrefix(namespace, "google3.") {
		// `google3` namespaces don't provide capitalization context for the import name
		return
	}

	lastSegment := namespace[strings.LastIndexByte(namespace, '.')+1:]
	if nextQnamePart == "" || nextQnamePart == lastSegment || lastSegment == "" {
		return
	}

	firstShort := jutil.StringCharAt(nextQnamePart, 0)
	if closureCheckModule_asciiIsUpperCase(firstShort) != closureCheckModule_asciiIsUpperCase(jutil.StringCharAt(lastSegment, 0)) {
		var newStartChar uint16
		if closureCheckModule_asciiIsUpperCase(firstShort) {
			newStartChar = closureCheckModule_asciiToLowerCase(firstShort)
		} else {
			newStartChar = closureCheckModule_asciiToUpperCase(firstShort)
		}
		correctedName := jutil.StringFromChar(newStartChar) + jutil.StringSubstringFrom(nextQnamePart, 1)
		t.Report(shortName.Node(), ClosureCheckModule_INCORRECT_SHORTNAME_CAPITALIZATION, nextQnamePart, correctedName)
	}
}

// Java: ClosureCheckModule.isValidDestructuringImport(Node)
func ClosureCheckModule_IsValidDestructuringImport(destructuringLhs ast.DestructuringLhs) bool {
	objectPattern, ok := destructuringLhs.Pattern().AsObjectPattern()
	if !ok {
		return false
	}
	// Every member is a string key whose value is a name: {a, b: c}.
	for prop := range objectPattern.Props() {
		stringKey, ok := prop.AsStringKey()
		if !ok {
			return false
		}
		if _, ok := stringKey.Value().AsName(); !ok {
			return false
		}
	}
	return true
}

// checkLegacyNamespaceCall validates the position of a
// goog.module.declareLegacyNamespace() call.
//
// Java: ClosureCheckModule.checkLegacyNamespaceCall(NodeTraversal, Node, Node)
func checkLegacyNamespaceCall(t *NodeTraversal, call ast.Call, parent rhino.Node) {
	callNode := call.Node()
	if call.HasArgs() {
		t.Report(callNode, ClosureCheckModule_LEGACY_NAMESPACE_ARGUMENT)
	}
	if !parent.IsExprResult() {
		t.Report(callNode, ClosureCheckModule_LEGACY_NAMESPACE_NOT_AT_TOP_LEVEL)
		return
	}
	prev := parent.GetPrevious()
	if prev.IsNil() || !NodeUtil_IsGoogModuleCall(prev) {
		t.Report(callNode, ClosureCheckModule_LEGACY_NAMESPACE_NOT_AFTER_GOOG_MODULE)
	}
}

// isFirstExpressionInGoogModule reports whether node is the first statement in
// a module.
//
// Java: ClosureCheckModule.isFirstExpressionInGoogModule(Node)
func isFirstExpressionInGoogModule(node rhino.Node) bool {
	if !node.IsExprResult() || !node.GetPrevious().IsNil() {
		return false
	}
	statementBlock := node.GetParent()
	return statementBlock.IsModuleBody() || NodeUtil_IsBundledGoogModuleScopeRoot(statementBlock)
}
