// Port of com/google/javascript/jscomp/ClosureRewriteModule.java @ closure-compiler v20261021

/*
 * Copyright 2014 The Closure Compiler Authors.
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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type closureScriptUpdater struct {
	outer              *ClosureRewriteModule
	scriptDescriptions []*closureScriptDescription
}

func (u *closureScriptUpdater) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	m := u.outer

	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		if len(u.scriptDescriptions) == 0 {
			panic(jutil.NewNoSuchElementException())
		}
		currentDescription := u.scriptDescriptions[0]
		u.scriptDescriptions[0] = nil
		u.scriptDescriptions = u.scriptDescriptions[1:]
		jutil.CheckState(currentDescription.rootNode == n)
		if n.IsFromExterns() && !NodeUtil_IsFromTypeSummary(n) {
			return false
		}

		jutil.CheckState(len(m.scriptStack) == 0)
		m.pushScript(currentDescription)

		// Capture the scope before doing any rewriting to the scope.
		t.GetScope()

		// Capture the global scope for later reference.
		if m.globalScope == nil {
			m.globalScope = t.GetScope().GetGlobalScope()
		}
	case rhino.Token_MODULE_BODY:
		if parent.GetBooleanProp(rhino.Node_GOOG_MODULE) {
			m.updateModuleBodyEarly(n)
		} else {
			return false
		}
	case rhino.Token_CALL:
		call := ast.MustCall(n)
		getProp, ok := call.Callee().AsGetProp()
		if !ok {
			break
		}
		method := getProp.Node()
		if method.MatchesQualifiedNameNode(closureRewriteModule_GOOG_MODULE) {
			m.updateGoogModule(t, call)
		} else if method.MatchesQualifiedNameNode(closureRewriteModule_GOOG_MODULE_DECLARELEGACYNAMESPACE) {
			NodeUtil_GetEnclosingStatement(n).Detach()
		} else if method.MatchesQualifiedNameNode(closureRewriteModule_GOOG_REQUIRE) ||
			method.MatchesQualifiedNameNode(closureRewriteModule_GOOG_REQUIRETYPE) {
			m.updateGoogRequire(t, call)
		} else if method.MatchesQualifiedNameNode(closureRewriteModule_GOOG_FORWARDDECLARE) && !parent.IsExprResult() {
			m.updateGoogForwardDeclare(t, call)
		} else if closureRewriteModule_GOOG_MODULE_PREVENTMODULEEXPORTSEALING.Matches(method) {
			NodeUtil_GetEnclosingStatement(n).Detach()
		}
	case rhino.Token_GETPROP:
		if getProp := ast.MustGetProp(n); m.currentScript.isModule && closureRewriteModule_isExportPropertyAssignment(t, getProp) {
			m.updateExportsPropertyAssignment(getProp, t)
		}
	}

	if info := n.GetJSDocInfo(); info != nil &&
		(m.currentScript.isModule ||
			len(m.currentScript.namesToInlineByAlias) > 0 ||
			len(m.rewriteState.scriptDescriptionsByGoogModuleNamespace) > 0) &&
		info.HasTypeNodes() {
		m.rewriteJsdocInfoWithTraversal(info, t)
	}

	return true
}

var closureScriptUpdaterVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_MODULE_BODY)
	s.Add(rhino.Token_NAME)
	s.Add(rhino.Token_SCRIPT)
	return s
}()

func (*closureScriptUpdater) visitTokens() *ast.TokenSet { return &closureScriptUpdaterVisitTokens }

var _ nodeTraversal_visitTokens = (*closureScriptUpdater)(nil)

func (u *closureScriptUpdater) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	m := u.outer
	switch n.GetToken() {
	case rhino.Token_MODULE_BODY:
		m.updateModuleBody(n)
	case rhino.Token_NAME:
		if m.currentScript.isModule {
			m.maybeUpdateTopLevelName(t, n)
			m.maybeUpdateExportDeclaration(t, n)
			t.GetScope() // Creating the scope here has load-bearing side-effects.
			m.maybeUpdateExportNameRef(t, n)
		}
	case rhino.Token_SCRIPT:
		jutil.CheckState(m.currentScript.rootNode == n)
		m.popScript()
	}
}

func (m *ClosureRewriteModule) updateModuleBodyEarly(moduleScopeRoot rhino.Node) {
	m.pushScript(m.currentScript.removeFirstChildScript())
	m.currentScript.rootNode = moduleScopeRoot
}

func (m *ClosureRewriteModule) updateGoogModule(t *NodeTraversal, googModule ast.Call) {
	call := googModule.Node()
	if !m.currentScript.isModule {
		m.compiler.ReportChangeToEnclosingScope(call)
		undefined := m.astFactory.CreateVoid(m.astFactory.CreateNumber(0)).SrcrefTree(call)
		call.ReplaceWith(undefined)
		return
	}

	if m.currentScript.declareLegacyNamespace {
		// Rewrite "goog.module('Foo');" as "goog.provide('Foo');".
		ast.MustGetProp(googModule.Callee().Node()).SetProperty("provide")
		m.compiler.ReportChangeToEnclosingScope(call)
	}

	if !m.currentScript.willCreateExportsObject {
		jutil.CheckState(!m.currentScript.hasCreatedExportObject)
		m.exportTheEmptyBinaryNamespaceAt(NodeUtil_GetEnclosingStatement(call), closureAddAfter, t)
	}

	if !m.currentScript.declareLegacyNamespace && !m.preserveSugar {
		m.compiler.ReportChangeToEnclosingScope(call)
		NodeUtil_GetEnclosingStatement(call).Detach()
	}
}

func (m *ClosureRewriteModule) updateGoogRequire(t *NodeTraversal, require ast.Call) {
	call := require.Node()
	namespaceIdNode := closureRewriteModule_namespaceArg(require)
	statementNode := NodeUtil_GetEnclosingStatement(call)
	namespaceId := namespaceIdNode.GetString()

	targetIsNonLegacyGoogModule := m.rewriteState.containsModule(namespaceId) && !m.rewriteState.isLegacyModule(namespaceId)
	decl, importHasAlias := ast.AsNameDeclaration(statementNode)
	isDestructuring := false
	if importHasAlias {
		if first, ok := decl.FirstDeclarator(); ok {
			_, isDestructuring = first.AsDestructuringLhs()
		}
	}

	currentScriptIsAModule := m.currentScript.isModule
	requireDirectlyStoredInAlias := NodeUtil_IsNameDeclaration(call.GetGrandparent())
	if currentScriptIsAModule &&
		requireDirectlyStoredInAlias &&
		m.isTopLevel(t, statementNode, closureScopeExecContext) {
		lhs := call.GetParent()
		exportedNamespace, hasExportedNamespace := m.rewriteState.getExportedNamespaceOrScript(namespaceId)
		if !hasExportedNamespace {
			// Missing provide/module will be reported elsewhere.
		} else if _, ok := ast.AsName(lhs); ok {
			aliasNameNode := decl.MustFirstDeclarator().MustName()
			aliasName := aliasNameNode.Ident()
			m.recordNameToInline(aliasName, exportedNamespace, namespaceId, true)
			m.maybeAddAliasToSymbolTable(aliasNameNode.Node(), m.currentScript.namespaceIdOrNull())
		} else if pattern, ok := closureRewriteModule_destructuringImport(lhs); ok {
			m.maybeWarnForInvalidDestructuring(t, lhs.GetParent(), namespaceId)
			moduleNs := m.currentScript.namespaceIdOrNull()
			prefix := m.currentScript.contentsPrefixOrNull()
			for prop := range pattern.Props() {
				importSpec := prop.Node()
				jutil.CheckStateMsg1(importSpec.HasChildren(), "%s", importSpec)
				var importedProperty string
				var aliasNode rhino.Node
				if key, ok := prop.AsStringKey(); ok {
					importedProperty = key.Key()
					aliasNode = key.Value().Node()
				} else {
					importedProperty = importSpec.GetString()
					aliasNode = importSpec.GetFirstChild()
				}
				aliasName := aliasNode.GetString()
				fullName := exportedNamespace + "." + importedProperty
				m.recordNameToInline(aliasName, fullName, "", false)

				m.maybeAddAliasToSymbolTable(aliasNode, moduleNs)
				m.safeSetString(aliasNode, prefix+aliasName)
			}
		} else {
			panic(jutil.NewRuntimeException("Illegal goog.module import: " + lhs.String()))
		}
	}

	if m.currentScript.isModule || targetIsNonLegacyGoogModule {
		if isDestructuring {
			if !m.preserveSugar {
				m.compiler.ReportChangeToEnclosingScope(statementNode)
				statementNode.Detach()
			}
		} else if targetIsNonLegacyGoogModule {
			jutil.CheckStateMsg1(
				m.isTopLevel(t, statementNode, closureScopeExecContext),
				"Unexpected non-top-level require at %s",
				call)
			if importHasAlias || !m.rewriteState.isLegacyModule(namespaceId) {
				if !m.preserveSugar {
					m.compiler.ReportChangeToEnclosingScope(statementNode)
					statementNode.Detach()
				}
			}
		} else {
			call.Detach()
			statementNode.ReplaceWith(rhino.IR_ExprResult(call))
			m.compiler.ReportChangeToEnclosingScope(call)
		}
	}
}

func (m *ClosureRewriteModule) maybeWarnForInvalidDestructuring(
	t *NodeTraversal, importNode rhino.Node, importedNamespace string) {
	objPattern, ok := closureRewriteModule_destructuringImport(importNode.GetFirstChild())
	jutil.CheckArgumentMsg1(ok, "%s", importNode)
	importedModule := m.rewriteState.scriptDescriptionsByGoogModuleNamespace[importedNamespace]
	if importedModule == nil {
		return
	}
	if importedModule.hasDefaultExport {
		t.Report(importNode, ClosureRewriteModule_ILLEGAL_DESTRUCTURING_DEFAULT_EXPORT)
		return
	}
	for prop := range objPattern.Props() {
		exportName := prop.Node().GetString()
		if _, ok := importedModule.namedExports[exportName]; !ok {
			t.Report(importNode, ClosureRewriteModule_ILLEGAL_DESTRUCTURING_NOT_EXPORTED, exportName, importedNamespace)
		}
	}
}

func closureRewriteModule_destructuringImport(lhs rhino.Node) (ast.ObjectPattern, bool) {
	destructuring, ok := ast.AsDestructuringLhs(lhs)
	if !ok {
		return ast.ObjectPattern{}, false
	}
	return destructuring.Pattern().AsObjectPattern()
}

func (m *ClosureRewriteModule) recordExportsPropertyAssignment(t *NodeTraversal, getProp ast.GetProp) {
	if !m.currentScript.isModule {
		return
	}

	parent := getProp.Node().GetParent()
	jutil.CheckStateMsg1(parent.IsAssign() || parent.IsExprResult(), "%s", parent)

	exportsNameNode := getProp.Object().Node()
	jutil.CheckStateMsg1(exportsNameNode.GetString() == "exports", "%s", exportsNameNode)

	if t.InModuleScope() {
		exportName := getProp.Property()
		if m.currentScript.namedExports == nil {
			m.currentScript.namedExports = make(map[string]struct{})
		}
		m.currentScript.namedExports[exportName] = struct{}{}
		var exportRhs rhino.Node
		if assign, ok := ast.AsAssign(parent); ok {
			exportRhs = assign.Value().Node()
		}
		namedExport := newNamedExportDefinition(t, exportName, exportRhs)
		if !m.currentScript.declareLegacyNamespace &&
			!m.currentScript.hasDefaultExport &&
			namedExport.hasInlinableName(m.currentScript.exportsToInlineVars) {
			m.recordExportToInline(namedExport)
			parent.GetParent().Detach()
		}
	}
}

func (m *ClosureRewriteModule) updateExportsPropertyAssignment(getProp ast.GetProp, t *NodeTraversal) {
	if !m.currentScript.isModule {
		return
	}

	getpropNode := getProp.Node()
	parent := getpropNode.GetParent()
	jutil.CheckStateMsg1(parent.IsAssign() || parent.IsExprResult(), "%s", parent)

	exportsNameNode := getProp.Object().Node()
	jutil.CheckState(exportsNameNode.GetString() == "exports")
	exportedNamespace, ok := m.currentScript.getExportedNamespace()
	if !ok {
		panic(jutil.NewNullPointerException())
	}
	m.safeSetMaybeQualifiedString(exportsNameNode, exportedNamespace /* isModuleNamespace= */, false)

	jsdocNode := getpropNode
	if parent.IsAssign() {
		jsdocNode = parent
	}
	closureRewriteModule_markConstAndCopyJsDoc(jsdocNode, jsdocNode)

	if !m.currentScript.hasCreatedExportObject {
		m.exportTheEmptyBinaryNamespaceAt(NodeUtil_GetEnclosingStatement(parent), closureAddBefore, t)
	}
}

func (m *ClosureRewriteModule) maybeUpdateTopLevelName(t *NodeTraversal, nameNode rhino.Node) {
	if !m.currentScript.isModule {
		return
	}
	if _, ok := m.currentScript.topLevelNameIDs[nameNode.GetStringID()]; !ok {
		return
	}
	name := nameNode.GetString()
	var_ := t.LookupVar(name)
	if var_ == nil || var_.GetScope().(*Scope).GetRootNode() != m.currentScript.rootNode {
		return
	}

	if var_.GetNameNode() == nameNode &&
		nameNode.GetParent().IsStringKey() &&
		nameNode.GetGrandparent().IsObjectPattern() {
		if dlhs, ok := ast.AsDestructuringLhs(nameNode.GetGrandparent().GetParent()); ok {
			if init, hasInit := dlhs.Init(); hasInit &&
				(NodeUtil_IsCallToNode(init.Node(), closureRewriteModule_GOOG_REQUIRE) ||
					NodeUtil_IsCallToNode(init.Node(), closureRewriteModule_GOOG_REQUIRETYPE)) {
				return
			}
		}
	}

	inline, nameIsAnAlias := m.currentScript.namesToInlineByAlias[name]
	if nameIsAnAlias && var_.GetNode() != nameNode {
		m.maybeAddAliasToSymbolTable(nameNode, m.currentScript.namespaceIdOrNull())

		namespaceToInline := inline.newName
		if binaryNamespace, hasBinary := m.currentScript.getBinaryNamespace(); hasBinary && namespaceToInline == binaryNamespace {
			m.currentScript.hasCreatedExportObject = true
		}
		var importedMod *closureScriptDescription
		if inline.hasNamespace {
			importedMod = m.rewriteState.scriptDescriptionsByGoogModuleNamespace[inline.namespaceId]
		}
		isModuleNamespace := importedMod != nil && !importedMod.willCreateExportsObject
		m.safeSetMaybeQualifiedString(nameNode, namespaceToInline, isModuleNamespace)

		if dot := strings.IndexByte(namespaceToInline, '.'); dot != -1 {
			firstQualifiedName := namespaceToInline[:dot]
			shadowedVar := t.LookupVar(firstQualifiedName)
			if shadowedVar == nil ||
				shadowedVar.IsGlobal() ||
				shadowedVar.GetScope().(*Scope).IsModuleScope() {
				return
			}
			t.Report(
				shadowedVar.GetNode(),
				ClosureRewriteModule_IMPORT_INLINING_SHADOWS_VAR,
				shadowedVar.GetName(),
				namespaceToInline)
		}
		return
	}

	m.safeSetString(nameNode, m.currentScript.contentsPrefixOrNull()+name)
}

func (m *ClosureRewriteModule) maybeUpdateExportObjectLiteral(t *NodeTraversal, n rhino.Node) {
	if !m.currentScript.isModule {
		return
	}

	rhs := n.GetParent().GetLastChild()
	if obj, ok := ast.AsObjectLit(rhs); ok {
		for member := range obj.Members() {
			if computed, ok := member.AsComputedProp(); ok {
				t.Report(computed.Node(), ClosureRewriteModule_INVALID_EXPORT_COMPUTED_PROPERTY)
			} else if key, ok := member.AsStringKey(); ok {
				m.maybeUpdateExportDeclToNode(t, key.Node(), key.Value().Node())
			}
		}
	}
}

func (m *ClosureRewriteModule) maybeUpdateExportDeclToNode(t *NodeTraversal, target rhino.Node, value rhino.Node) {
	if !m.currentScript.isModule {
		return
	}

	if value.IsName() {
		currentScope := t.GetScope()
		if v := currentScope.GetVarByNode(value); v != nil {
			varScope := v.GetScope().(*Scope)
			if varScope.GetDepth() == currentScope.GetDepth() {
				if info := v.GetJSDocInfo(); info != nil && info.HasTypedefType() {
					builder := rhino.JSDocInfo_Builder_CopyFrom(info)
					target.SetJSDocInfo(builder.Build())
					return
				}
			}
		}
	}

	closureRewriteModule_markConstAndCopyJsDoc(target, target)
}

func (m *ClosureRewriteModule) maybeUpdateExportDeclaration(t *NodeTraversal, n rhino.Node) {
	if !m.currentScript.isModule || !closureRewriteModule_isGoogModuleExportsRef(t, n) {
		return
	}
	assign, ok := closureRewriteModule_assignToTarget(n)
	if !ok {
		return
	}

	assignNode := assign.Node()
	rhs := assign.Value().Node()
	if rhs != m.currentScript.defaultExport.rhs {
		// Duplicate 'exports = ' assignment: preserve rhs as an expression.
		assignNode.ReplaceWith(rhs.Detach())
		return
	}

	if !m.currentScript.declareLegacyNamespace && m.currentScript.hasDefaultExportLocalName {
		assignNode.GetParent().Detach()

		binaryNs, _ := m.currentScript.getBinaryNamespace()
		binaryNamespaceName := m.astFactory.CreateName(binaryNs, AstFactory_Type_(n))
		m.declareGlobalVariable(binaryNamespaceName, t)
		return
	}

	var jsdocNode rhino.Node
	if m.currentScript.declareLegacyNamespace {
		if !m.currentScript.hasNamespaceId {
			panic(jutil.NewNullPointerException())
		}
		legacyQname := m.astFactory.
			CreateQNameUsingJSTypeInfo(m.globalTypedScope, m.currentScript.namespaceId).
			SrcrefTree(n)
		legacyQname.SetJSType(n.GetJSType())
		n.ReplaceWith(legacyQname)
		jsdocNode = assignNode
	} else {
		rhs.Detach()
		exprResultNode := assignNode.GetParent()
		binaryNs, _ := m.currentScript.getBinaryNamespace()
		binaryNamespaceName := m.astFactory.CreateName(binaryNs, AstFactory_Type_(n))
		exportsName := "exports"
		binaryNamespaceName.SetOriginalName(&exportsName)
		m.declareGlobalVariable(binaryNamespaceName, t)

		exportsObjectCreationNode := rhino.IR_VarWithValue(binaryNamespaceName, rhs)
		exportsObjectCreationNode.SrcrefTreeIfMissing(exprResultNode)
		exportsObjectCreationNode.PutBooleanProp(rhino.Node_IS_NAMESPACE, true)
		exprResultNode.ReplaceWith(exportsObjectCreationNode)
		jsdocNode = exportsObjectCreationNode
		m.currentScript.hasCreatedExportObject = true
	}
	closureRewriteModule_markConstAndCopyJsDoc(assignNode, jsdocNode)
	m.compiler.ReportChangeToEnclosingScope(jsdocNode)

	m.maybeUpdateExportObjectLiteral(t, rhs)
}

func (m *ClosureRewriteModule) maybeUpdateExportNameRef(t *NodeTraversal, n rhino.Node) {
	if !m.currentScript.isModule || !closureRewriteModule_isGoogModuleExportsRef(t, n) || n.GetParent().IsNil() {
		return
	}
	if n.GetParent().IsParamList() {
		return
	}

	if m.currentScript.declareLegacyNamespace {
		if !m.currentScript.hasNamespaceId {
			panic(jutil.NewNullPointerException())
		}
		legacyQname := m.astFactory.
			CreateQNameUsingJSTypeInfo(m.globalTypedScope, m.currentScript.namespaceId).
			SrcrefTree(n)
		legacyQname.SetJSType(n.GetJSType())
		n.ReplaceWith(legacyQname)
		m.compiler.ReportChangeToEnclosingScope(legacyQname)
		return
	}

	m.safeSetString(n, m.currentScript.binaryNamespaceOrNull())
	jutil.CheckState(m.currentScript.willCreateExportsObject || m.currentScript.hasCreatedExportObject)
}

func (m *ClosureRewriteModule) updateModuleBody(moduleBody rhino.Node) {
	jutil.CheckArgumentMsg1(
		moduleBody.IsModuleBody() && moduleBody.GetParent().GetBooleanProp(rhino.Node_GOOG_MODULE),
		"%s",
		moduleBody)
	moduleBody.SetToken(rhino.Token_BLOCK)
	NodeUtil_TryMergeBlock(moduleBody, true)

	if len(m.currentScript.exportsToInline) > 0 {
		binaryNs := m.currentScript.binaryNamespaceOrNull()
		for i := range m.currentScript.exportsToInline {
			export := &m.currentScript.exportsToInline[i]
			nameNode := export.nameDecl.GetNameNode()
			m.safeSetMaybeQualifiedString(
				nameNode, binaryNs+export.exportPostfix(), false)
		}
	}
	jutil.CheckState(m.currentScript.isModule)
	jutil.CheckState(m.currentScript.declareLegacyNamespace || m.currentScript.hasCreatedExportObject)

	m.popScript()
}

func (m *ClosureRewriteModule) pushScript(newCurrentScript *closureScriptDescription) {
	m.currentScript = newCurrentScript
	if len(m.scriptStack) > 0 {
		parentScript := m.scriptStack[len(m.scriptStack)-1]
		parentScript.addChildScript(m.currentScript)
	}
	m.scriptStack = append(m.scriptStack, m.currentScript)
}

func (m *ClosureRewriteModule) popScript() {
	if len(m.scriptStack) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	last := len(m.scriptStack) - 1
	m.scriptStack[last] = nil
	m.scriptStack = m.scriptStack[:last]
	if len(m.scriptStack) == 0 {
		m.currentScript = nil
	} else {
		m.currentScript = m.scriptStack[len(m.scriptStack)-1]
	}
}

func (m *ClosureRewriteModule) exportTheEmptyBinaryNamespaceAt(atNode rhino.Node, addAt closureAddAt, t *NodeTraversal) {
	if m.currentScript.declareLegacyNamespace {
		return
	}

	binaryNamespaceString, _ := m.currentScript.getBinaryNamespace()
	moduleType := AstFactory_Type_(m.currentScript.rootNode)
	binaryNamespaceName := m.astFactory.CreateName(binaryNamespaceString, moduleType)
	if m.currentScript.hasNamespaceId {
		binaryNamespaceName.SetOriginalName(&m.currentScript.namespaceId)
	}
	m.declareGlobalVariable(binaryNamespaceName, t)

	binaryNamespaceExportNode := rhino.IR_VarWithValue(binaryNamespaceName, m.astFactory.CreateObjectLit())
	if addAt == closureAddBefore {
		binaryNamespaceExportNode.InsertBefore(atNode)
	} else if addAt == closureAddAfter {
		binaryNamespaceExportNode.InsertAfter(atNode)
	}
	binaryNamespaceExportNode.PutBooleanProp(rhino.Node_IS_NAMESPACE, true)
	binaryNamespaceExportNode.SrcrefTree(atNode)
	closureRewriteModule_markConst(binaryNamespaceExportNode)
	m.compiler.ReportChangeToEnclosingScope(binaryNamespaceExportNode)
	m.currentScript.hasCreatedExportObject = true
}

// Java: ClosureRewriteModule.checkAndSetStrictModeDirective(NodeTraversal, Node)
func ClosureRewriteModule_CheckAndSetStrictModeDirective(t *NodeTraversal, n rhino.Node) {
	jutil.CheckStateMsg1(n.IsScript(), "%s", n)

	if n.IsUseStrict() {
		t.Report(n, ClosureRewriteModule_USELESS_USE_STRICT_DIRECTIVE)
	} else {
		n.SetUseStrict(true)
	}
}

func closureRewriteModule_maybeSplitMultiVar(rhsNode rhino.Node) {
	statementNode := rhsNode.GetGrandparent()
	if !statementNode.IsVar() || !statementNode.HasMoreThanOneChild() {
		return
	}

	nameNode := rhsNode.GetParent()
	nameNode.Detach()
	rhsNode.Detach()
	rhino.IR_VarWithValue(nameNode, rhsNode).InsertBefore(statementNode)
}

func (m *ClosureRewriteModule) recordExportToInline(exportDefinition closureExportDefinition) {
	jutil.CheckState(exportDefinition.hasInlinableName(m.currentScript.exportsToInlineVars))
	if m.currentScript.exportsToInlineVars == nil {
		m.currentScript.exportsToInlineVars = make(map[*Var]struct{})
	}
	_, exists := m.currentScript.exportsToInlineVars[exportDefinition.nameDecl]
	jutil.CheckStateMsg1(
		!exists,
		"Already found a mapping for inlining export: %s",
		exportDefinition.nameDecl)
	m.currentScript.exportsToInlineVars[exportDefinition.nameDecl] = struct{}{}
	m.currentScript.exportsToInline = append(m.currentScript.exportsToInline, exportDefinition)

	localName, _ := exportDefinition.localName()
	fullExportedName := m.currentScript.binaryNamespaceOrNull() + exportDefinition.exportPostfix()
	m.recordNameToInline(localName, fullExportedName, "", false)
}

func (m *ClosureRewriteModule) recordNameToInline(aliasName, newName, namespaceId string, hasNamespace bool) {
	if m.currentScript.namesToInlineByAlias == nil {
		m.currentScript.namesToInlineByAlias = make(map[string]closureAliasName)
	}
	m.currentScript.namesToInlineByAlias[aliasName] = closureAliasName{
		newName:      newName,
		namespaceId:  namespaceId,
		hasNamespace: hasNamespace,
	}
}

func (m *ClosureRewriteModule) reportUnrecognizedRequires() {
	for i := range m.unrecognizedRequires {
		unrecognizedRequire := &m.unrecognizedRequires[i]
		namespaceId := unrecognizedRequire.namespaceId

		requireNode := unrecognizedRequire.requireNode
		targetGoogModuleExists := m.rewriteState.containsModule(namespaceId)
		_, targetLegacyScriptExists := m.rewriteState.providedNamespaces[namespaceId]

		if targetGoogModuleExists || targetLegacyScriptExists || m.preserveSugar {
			continue
		}

		if NodeUtil_GetEnclosingScript(requireNode).IsNil() {
			continue
		}

		m.compiler.ReportChangeToEnclosingScope(requireNode)
		enclosingStatement := NodeUtil_GetEnclosingStatement(requireNode)

		if !NodeUtil_IsNameDeclaration(enclosingStatement) {
			requireNode.ReplaceWith(m.astFactory.CreateNull().Srcref(requireNode))
			continue
		}

		enclosingStatement.Detach()
		NodeUtil_VisitLhsNodesInNode(
			enclosingStatement, func(lhs rhino.Node) { m.addSyntheticExternIfAbsent(lhs) })
	}
}

func (m *ClosureRewriteModule) safeSetString(n rhino.Node, newString string) {
	originalName := n.GetString()
	if originalName == newString {
		return
	}

	n.SetString(newString)
	if n.GetOriginalName() == nil {
		n.SetOriginalName(&originalName)
	}
	changeScope := ChangeTracker_GetEnclosingChangeScopeRoot(n)
	if !changeScope.IsNil() {
		m.compiler.ReportChangeToChangeScope(changeScope)
	}
}

func (m *ClosureRewriteModule) safeSetMaybeQualifiedString(
	nameNode rhino.Node, newString string, isModuleNamespace bool) {
	if strings.IndexByte(newString, '.') == -1 {
		m.safeSetString(nameNode, newString)
		parent := nameNode.GetParent()
		if call, ok := ast.AsCall(nameNode.GetGrandparent()); isModuleNamespace &&
			parent.IsGetProp() &&
			ok && call.Callee().Is(parent) {
			nameNode.GetGrandparent().PutBooleanProp(rhino.Node_FREE_CALL, true)
		}
		return
	}

	nameParent := nameNode.GetParent()
	newQualifiedName := m.astFactory.
		CreateQNameUsingJSTypeInfo(m.globalTypedScope, newString).
		SrcrefTree(nameNode)
	newQualifiedName.SetJSType(nameNode.GetJSType())

	if closureRewriteModule_safeSetStringIfDeclaration(nameParent, nameNode, newQualifiedName) {
		return
	}

	nameNode.ReplaceWith(newQualifiedName)
	if newQualifiedName.HasChildren() {
		newQualifiedName.GetFirstChild().MakeNonIndexableRecursive()
	}
	m.compiler.ReportChangeToEnclosingScope(newQualifiedName)
}

func closureRewriteModule_safeSetStringIfDeclaration(
	nameParent rhino.Node, nameNode rhino.Node, newQualifiedName rhino.Node) bool {
	jsdoc := nameParent.GetJSDocInfo()

	switch nameParent.GetToken() {
	case rhino.Token_FUNCTION, rhino.Token_CLASS:
		if !NodeUtil_IsStatement(nameParent) || nameParent.GetFirstChild() != nameNode {
			return false
		}

		placeholder := rhino.IR_Empty()
		nameParent.ReplaceWith(placeholder)
		newDeclaration := NodeUtil_GetDeclarationFromName(newQualifiedName, nameParent, rhino.Token_VAR, jsdoc)
		if NodeUtil_IsExprAssign(newDeclaration) {
			assign := ast.MustExprStmt(newDeclaration).Expr().Node()
			assign.SetJSType(nameNode.GetJSType())
			closureRewriteModule_updateSourceInfoForExportedTopLevelVariable(assign, nameNode)
		}
		nameParent.SetJSDocInfo(nil)
		newDeclaration.SrcrefTreeIfMissing(nameParent)
		placeholder.ReplaceWith(newDeclaration)
		NodeUtil_RemoveName(nameParent)
		return true
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		var rhs rhino.Node
		if init, ok := ast.MustName(nameNode).Init(); ok {
			rhs = init.Detach().Node()
		}
		if jsdoc == nil {
			jsdoc = nameNode.GetJSDocInfo()
		}
		newStatement := NodeUtil_GetDeclarationFromName(newQualifiedName, rhs, rhino.Token_VAR, jsdoc)
		if NodeUtil_IsExprAssign(newStatement) {
			assign := ast.MustExprStmt(newStatement).Expr().Node()
			assign.SetJSType(nameNode.GetJSType())
			closureRewriteModule_updateSourceInfoForExportedTopLevelVariable(assign, nameNode)
			if nameParent.IsConst() {
				jsdocBuilder := rhino.JSDocInfo_Builder_MaybeCopyFrom(jsdoc)
				jsdocBuilder.RecordConstancy()
				assign.SetJSDocInfo(jsdocBuilder.Build())
			}
		}
		newStatement.SrcrefTreeIfMissing(nameParent)
		NodeUtil_ReplaceDeclarationChild(nameNode, newStatement)
		return true
	case rhino.Token_OBJECT_PATTERN, rhino.Token_ARRAY_PATTERN, rhino.Token_PARAM_LIST:
		panic(jutil.NewRuntimeException("Not supported"))
	}
	return false
}

func closureRewriteModule_updateSourceInfoForExportedTopLevelVariable(assignNode rhino.Node, sourceName rhino.Node) {
	assign, ok := ast.AsAssign(assignNode)
	jutil.CheckState(ok)
	jutil.CheckState(sourceName.IsName())

	getProp, ok := assign.Target().AsGetProp()
	if !ok {
		return
	}

	var name string
	if originalName := sourceName.GetOriginalName(); originalName != nil {
		name = *originalName
	} else {
		name = sourceName.GetString()
	}

	sourceName = sourceName.CloneNode()
	sourceName.SetLength(jutil.StringLength(name))
	getProp.Node().SrcrefTree(sourceName)
}

func (m *ClosureRewriteModule) isTopLevel(t *NodeTraversal, n rhino.Node, scopeType closureScopeType) bool {
	if scopeType == closureScopeExecContext {
		return t.InGlobalScope() || t.GetClosestHoistScopeRoot() == m.currentScript.rootNode
	}
	return n.GetParent() == m.currentScript.rootNode
}

// Java: ClosureRewriteModule.isModuleExport(String)
func ClosureRewriteModule_IsModuleExport(name string) bool {
	return strings.HasPrefix(name, ClosureRewriteModule_MODULE_EXPORTS_PREFIX)
}

// Java: ClosureRewriteModule.isModuleContent(String)
func ClosureRewriteModule_IsModuleContent(name string) bool {
	return strings.HasPrefix(name, closureRewriteModule_MODULE_CONTENTS_PREFIX)
}

func closureRewriteModule_isGoogModuleExportsRef(t *NodeTraversal, target rhino.Node) bool {
	if !target.IsName() || target.GetStringID() != closureModuleExportsID {
		return false
	}
	exportsVar := t.LookupVar("exports")
	return exportsVar != nil && exportsVar.IsGoogModuleExports()
}

func closureRewriteModule_isExportPropertyAssignment(t *NodeTraversal, getProp ast.GetProp) bool {
	n := getProp.Node()
	return (closureRewriteModule_isAssignTarget(n) || closureRewriteModule_isTypedefTarget(n)) &&
		closureRewriteModule_isGoogModuleExportsRef(t, getProp.Object().Node())
}

func closureRewriteModule_isAssignTarget(n rhino.Node) bool {
	_, ok := closureRewriteModule_assignToTarget(n)
	return ok
}

func closureRewriteModule_assignToTarget(n rhino.Node) (ast.Assign, bool) {
	assign, ok := ast.AsAssign(n.GetParent())
	return assign, ok && assign.Target().Is(n)
}

func closureRewriteModule_isTypedefTarget(n rhino.Node) bool {
	stmt, ok := ast.AsExprStmt(n.GetParent())
	return ok && stmt.Expr().Is(n)
}

func (m *ClosureRewriteModule) declareGlobalVariable(n rhino.Node, t *NodeTraversal) {
	jutil.CheckState(n.IsName())
	if m.globalTypedScope == nil {
		return
	}

	name := n.GetString()
	if m.globalTypedScope.HasOwnSlot(name) {
		original := m.globalTypedScope.GetOwnSlot(name).(*TypedVar).GetNode()
		originalString := "<unknown>"
		if !original.IsNil() {
			originalString = original.String()
		}
		t.Report(
			t.GetCurrentScript(),
			ClosureRewriteModule_ILLEGAL_MODULE_RENAMING_CONFLICT,
			name,
			originalString)
	} else {
		type_ := jutil.CheckNotNull(n.GetJSType())
		m.globalTypedScope.Declare(name, n, type_, t.GetInput(), false)
	}
}
