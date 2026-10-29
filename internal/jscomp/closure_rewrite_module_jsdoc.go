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
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func (m *ClosureRewriteModule) rewriteJsdocInfoWithTraversal(info *rhino.JSDocInfo, t *NodeTraversal) {
	for it := info.TypeExprIter(); it.Next(); {
		m.rewriteJsdocTypeTree(t, it.Node())
	}
}

func (m *ClosureRewriteModule) rewriteJsdocTypeTree(t *NodeTraversal, n rhino.Node) {
	m.visitJsdocTypeRef(t, n)
	for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		m.rewriteJsdocTypeTree(t, child)
	}
}

func isClosureLocalShadowVar(t *NodeTraversal, rootOfType string) bool {
	if t.InGlobalScope() {
		return false
	}
	rootVar := t.LookupVar(rootOfType)
	if rootVar == nil {
		return false
	}
	varScope := rootVar.GetScope().(*Scope)
	return !varScope.GetClosestHoistScope().IsGlobal() && !varScope.IsModuleScope()
}

func (m *ClosureRewriteModule) visitJsdocTypeRef(t *NodeTraversal, typeRefNode rhino.Node) {
	if !typeRefNode.IsStringLit() {
		return
	}
	typeName := typeRefNode.GetString()
	dot := strings.IndexByte(typeName, '.')
	rootOfType := typeName
	if dot != -1 {
		rootOfType = typeName[:dot]
	}

	if alias, ok := m.currentScript.namesToInlineByAlias[rootOfType]; ok {
		if isClosureLocalShadowVar(t, rootOfType) {
			return
		}
		if m.preprocessorSymbolTable != nil {
			moduleOnlyNode := typeRefNode.CloneNode()
			m.safeSetString(moduleOnlyNode, rootOfType)
			moduleOnlyNode.SetLength(jutil.StringLength(rootOfType))
			m.maybeAddAliasToSymbolTable(moduleOnlyNode, m.currentScript.namespaceIdOrNull())
		}

		newType := alias.newName
		if dot != -1 {
			newType += typeName[dot:]
		}
		m.safeSetString(typeRefNode, newType)
	} else if m.currentScript.isModule && len(m.currentScript.topLevelNameIDs) > 0 &&
		hasUint32Key(m.currentScript.topLevelNameIDs, rhino.RhinoStringPool_LookupID(rootOfType)) {
		if isClosureLocalShadowVar(t, rootOfType) {
			return
		}
		m.safeSetString(typeRefNode, m.currentScript.contentsPrefixOrNull()+typeName)
	} else if m.currentScript.isModule && rootOfType == "exports" {
		if isClosureLocalShadowVar(t, rootOfType) {
			return
		}
		newType := m.currentScript.binaryNamespaceOrNull()
		if dot != -1 {
			newType += typeName[dot:]
		}
		m.safeSetString(typeRefNode, newType)
	} else if len(m.rewriteState.scriptDescriptionsByGoogModuleNamespace) > 0 {
		m.rewriteIfClosureNamespaceRefForRoot(t, typeName, rootOfType, typeRefNode)
	}
}

func (m *ClosureRewriteModule) rewriteIfClosureNamespaceRefForRoot(
	t *NodeTraversal, typeName, rootOfType string, typeRefNode rhino.Node) {
	prefixTypeName := typeName
	suffix := ""

	for {
		binaryNamespaceIfModule, hasBinary := m.rewriteState.getBinaryNamespace(prefixTypeName)
		if _, isLegacy := m.legacyScriptNamespacesAndPrefixes[prefixTypeName]; isLegacy && !hasBinary {
			return
		}
		if hasBinary {
			if isClosureLocalShadowVar(t, rootOfType) {
				return
			}
			m.safeSetString(typeRefNode, binaryNamespaceIfModule+suffix)
			return
		}
		dot := strings.LastIndexByte(prefixTypeName, '.')
		if dot < 0 {
			break
		}
		prefixTypeName = prefixTypeName[:dot]
		suffix = typeName[len(prefixTypeName):]
	}
}

func hasUint32Key(m map[uint32]struct{}, k uint32) bool {
	if k == 0 {
		return false
	}
	_, ok := m[k]
	return ok
}

func closureRewriteModule_markConst(n rhino.Node) {
	builder := rhino.JSDocInfo_Builder_MaybeCopyFrom(n.GetJSDocInfo())
	builder.RecordConstancy()
	n.SetJSDocInfo(builder.Build())
}

func closureRewriteModule_markConstAndCopyJsDoc(from rhino.Node, target rhino.Node) {
	info := from.GetJSDocInfo()
	builder := rhino.JSDocInfo_Builder_MaybeCopyFrom(info)
	builder.RecordConstancy()
	target.SetJSDocInfo(builder.Build())
}

func (m *ClosureRewriteModule) maybeAddAliasToSymbolTable(n rhino.Node, moduleOrNull string) {
	if m.preprocessorSymbolTable == nil {
		return
	}
	n.PutBooleanProp(rhino.Node_MODULE_ALIAS, true)
	var nodeNameStr string
	if n.IsStringLit() {
		nodeNameStr = n.GetString()
	} else if qn := m.preprocessorSymbolTable.GetQualifiedName(n); qn != nil {
		nodeNameStr = *qn
	} else {
		nodeNameStr = "null"
	}
	name := "alias_" + moduleOrNull + "_" + nodeNameStr
	m.preprocessorSymbolTable.AddReferenceWithName(n, name)
}

func (m *ClosureRewriteModule) updateGoogForwardDeclare(t *NodeTraversal, forwardDeclare ast.Call) {
	call := forwardDeclare.Node()
	namespaceId := closureRewriteModule_namespaceArg(forwardDeclare).GetString()
	if m.isUnrequiredWeakNamespace(call, namespaceId) {
		m.compiler.ReportChangeToEnclosingScope(call)
		call.ReplaceWith(m.astFactory.CreateNull().Srcref(call))
		return
	}
	m.updateGoogRequire(t, forwardDeclare)
}

func (m *ClosureRewriteModule) updateGoogRequireDynamicCall(call ast.Call) {
	parent := call.Node().GetParent()
	if parent.IsNil() {
		return
	}
	if await, ok := ast.AsAwait(parent); ok {
		m.updateGoogRequireDynamicCallInAwait(call, await)
		return
	}
	then, ok := ast.AsGetProp(parent)
	jutil.CheckStateMsg(
		ok && then.Property() == "then",
		"goog.requireDynamic() in only allowed in await/then expression")
	m.updateGoogRequireDynamicCallInThen(call, then)
}

func (m *ClosureRewriteModule) updateGoogRequireDynamicCallInThen(requireDynamic ast.Call, then ast.GetProp) {
	call := requireDynamic.Node()
	namespaceIdNode := closureRewriteModule_firstArg(requireDynamic)
	namespaceId := namespaceIdNode.GetString()
	exportedNamespace, hasExportedNamespace := m.rewriteState.getExportedNamespaceOrScript(namespaceId)
	jutil.CheckStateMsg(
		hasExportedNamespace, "Exported namespace for goog.requireDynamic() canot be null")

	thenCall, ok := ast.AsCall(then.Node().GetParent())
	jutil.CheckStateMsg(ok, "must be a 'then' call expression")

	argNode := m.astFactory.CreateString(m.namespaceIdToXid(namespaceId)).Srcref(namespaceIdNode)
	namespaceIdNode.ReplaceWith(argNode)
	closureRewriteModule_renameRequireDynamic(requireDynamic)

	exportedNamespaceNameNode := m.astFactory.CreateNameInScope(m.globalTypedStaticScope(), exportedNamespace).SrcrefTree(call)
	exportedNamespaceNameNode.SetOriginalName(&namespaceId)

	callback, _ := thenCall.FirstArg()
	functionNode, ok := ast.AsFunction(callback.Node())
	jutil.CheckStateMsg(ok, "must be a function in `then`")

	paramListNode := functionNode.Params().Node()
	param, ok := functionNode.Params().OnlyParam()
	jutil.CheckStateMsg(ok, "function must have only one parameter")

	functionBody := functionNode.Body().Node()

	objectPatternOrNameNode := param.Node()
	_, isObjectPattern := param.AsObjectPattern()
	_, isName := param.AsName()
	jutil.CheckStateMsg(isObjectPattern || isName, "parameter of callback function must be object pattern or name")
	objectPatternOrNameNode.Detach()
	paramListNode.RemoveChildren()

	languageIn := m.compiler.GetOptions().GetLanguageIn()
	jutil.CheckStateMsg2(
		languageIn.ToFeatureSet().ContainsFeature(parser.FeatureSet_Feature_CONST_DECLARATIONS),
		"'%s' does not contain '%s'",
		languageIn,
		parser.FeatureSet_Feature_CONST_DECLARATIONS)

	enclosingScript := NodeUtil_GetEnclosingScript(call)
	NodeUtil_AddFeatureToScript(enclosingScript, parser.FeatureSet_Feature_CONST_DECLARATIONS, m.compiler)

	var declarationNode rhino.Node
	if isObjectPattern {
		declarationNode = m.astFactory.
			CreateSingleConstObjectPatternDeclaration(
				objectPatternOrNameNode, exportedNamespaceNameNode).
			SrcrefTreeIfMissing(call)
	} else {
		declarationNode = m.astFactory.
			CreateSingleConstNameDeclaration(
				objectPatternOrNameNode.GetString(), exportedNamespaceNameNode).
			SrcrefTreeIfMissing(call)
	}

	if block, ok := functionNode.Body().AsBlock(); ok {
		block.Node().AddChildToFront(declarationNode)
		m.compiler.ReportChangeToEnclosingScope(declarationNode)
	} else {
		returnStmt := m.astFactory.CreateReturn(functionBody.Detach())
		newBlock := m.astFactory.CreateBlock(returnStmt).SrcrefTree(call)
		newBlock.InsertAfter(paramListNode)
		newBlock.AddChildToFront(declarationNode)
		m.compiler.ReportChangeToEnclosingScope(newBlock)
	}

	m.compiler.ReportChangeToEnclosingScope(call)
}

func (m *ClosureRewriteModule) updateGoogRequireDynamicCallInAwait(requireDynamic ast.Call, await ast.Await) {
	call := requireDynamic.Node()
	namespaceIdNode := closureRewriteModule_firstArg(requireDynamic)
	namespaceId := namespaceIdNode.GetString()
	exportedNamespace, hasExportedNamespace := m.rewriteState.getExportedNamespaceOrScript(namespaceId)
	jutil.CheckStateMsg(
		hasExportedNamespace, "Exported namespace for goog.requireDynamic() cannot be null")

	existingAwait := await.Node()
	awaitParent := existingAwait.GetParent()
	if awaitParent.IsNil() || (!awaitParent.IsDestructuringLhs() && !awaitParent.IsName()) {
		m.compiler.Report(JSError_MakeWithNode(call, ClosureRewriteModule_ILLEGAL_STMT_OF_GOOG_REQUIRE_DYNAMIC_IN_AWAIT))
	}

	declarationStatement := awaitParent.GetParent()
	if !NodeUtil_IsNameDeclaration(declarationStatement) ||
		!NodeUtil_IsStatement(declarationStatement) {
		m.compiler.Report(JSError_MakeWithNode(call, ClosureRewriteModule_ILLEGAL_STMT_OF_GOOG_REQUIRE_DYNAMIC_IN_AWAIT))
	}

	exportedNamespaceNameNode := m.astFactory.
		CreateQNameUsingJSTypeInfo(m.globalTypedScope, exportedNamespace).
		SrcrefTree(call)
	exportedNamespaceNameNode.SetJSType(m.rewriteState.getGoogModuleNamespaceType(namespaceId))
	exportedNamespaceNameNode.SetOriginalName(&namespaceId)

	argNode := m.astFactory.CreateString(m.namespaceIdToXid(namespaceId)).Srcref(namespaceIdNode)
	namespaceIdNode.ReplaceWith(argNode)
	closureRewriteModule_renameRequireDynamic(requireDynamic)

	existingAwait.ReplaceWith(exportedNamespaceNameNode)

	awaitStatement := m.astFactory.ExprResult(existingAwait).Srcref(existingAwait)
	awaitStatement.InsertBefore(declarationStatement)
	m.compiler.ReportChangeToEnclosingScope(call)
}

func closureRewriteModule_firstArg(call ast.Call) rhino.Node {
	first, _ := call.FirstArg()
	return first.Node()
}

func closureRewriteModule_renameRequireDynamic(call ast.Call) {
	calleeNode := call.Callee().Node()
	jutil.CheckStateMsg1(calleeNode.MatchesQualifiedName(closureRewriteModule_GOOG_REQUIREDYNAMIC_NAME), "%s", calleeNode)
	ast.MustGetProp(calleeNode).SetProperty(closureRewriteModule_IMPORT_HANDLER_NAME)
}

func (m *ClosureRewriteModule) namespaceIdToXid(namespaceId string) string {
	if m.xid == nil {
		if hashFunction := m.compiler.GetOptions().GetChunkIdHashFunction(); hashFunction != nil {
			m.xid = NewXidWithHasher(hashFunction)
		} else {
			m.xid = NewXid()
		}
	}
	return m.xid.Get(namespaceId)
}

func (m *ClosureRewriteModule) updateGoogModuleGetCall(moduleGet ast.Call) {
	call := moduleGet.Node()
	namespaceId := closureRewriteModule_firstArg(moduleGet).GetString()

	if m.isUnrequiredWeakNamespace(call, namespaceId) {
		m.compiler.Report(JSError_MakeWithNode(call, ClosurePrimitiveErrors_GOOG_MODULE_GET_OF_WEAK_MODULE, namespaceId))
		m.compiler.ReportChangeToEnclosingScope(call)
		call.ReplaceWith(m.astFactory.CreateNull().Srcref(call))
		return
	}

	if exportedNamespace, ok := m.rewriteState.getExportedNamespaceOrScript(namespaceId); ok {
		m.compiler.ReportChangeToEnclosingScope(call)
		exportedNamespaceName := m.astFactory.
			CreateQNameUsingJSTypeInfo(m.globalTypedScope, exportedNamespace).
			SrcrefTree(call)
		exportedNamespaceName.SetJSType(m.rewriteState.getGoogModuleNamespaceType(namespaceId))
		exportedNamespaceName.SetOriginalName(&namespaceId)
		call.ReplaceWith(exportedNamespaceName)
	}
}
