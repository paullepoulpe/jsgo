// Port of com/google/javascript/jscomp/ProcessCommonJSModules.java @ closure-compiler v20261021

/*
 * Copyright 2011 The Closure Compiler Authors.
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
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func sameInputId(a, b *rhino.InputId) bool {
	if a == b {
		return true
	}
	return a != nil && b != nil && a.Equals(b)
}

// commonJsRewriteModule traverses a file and rewrites all references to
// imported names directly to the targeted module name, and rewrites CommonJS
// export statements.
type commonJsRewriteModule struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	outer *ProcessCommonJSModules

	allowFullRewrite          bool
	exports                   []commonJsExportInfo
	imports                   []rhino.Node
	rewrittenClassExpressions []rhino.Node
	functionsToHoist          []rhino.Node
	defaultExportIsConst      bool
}

func (r *commonJsRewriteModule) reset(allowFullRewrite bool, defaultExportIsConst bool) {
	r.allowFullRewrite = allowFullRewrite
	r.defaultExportIsConst = defaultExportIsConst
	r.exports = r.exports[:0]
	r.imports = r.imports[:0]
	r.rewrittenClassExpressions = r.rewrittenClassExpressions[:0]
	r.functionsToHoist = r.functionsToHoist[:0]
}

func (r *commonJsRewriteModule) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	compiler := r.outer.compiler
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		for _, clazz := range r.rewrittenClassExpressions {
			clazz.GetFirstChild().ReplaceWith(rhino.IR_Empty().Srcref(clazz.GetFirstChild()))
			t.ReportCodeChange()
		}

		ci := compiler.GetInput(n.GetInputId())
		moduleName := commonJsModuleName(ci)

		if len(r.functionsToHoist) > 1 {
			baseImport := r.outer.getBasePropertyImport(moduleName)
			for i := 1; i < len(r.functionsToHoist); i++ {
				if r.functionsToHoist[i].GetFirstFirstChild().MatchesQualifiedName(baseImport) {
					fncVar := r.functionsToHoist[i]
					r.functionsToHoist = slices.Delete(r.functionsToHoist, i, i+1)
					r.functionsToHoist = slices.Insert(r.functionsToHoist, 0, fncVar)
					break
				}
			}
		}

		for i := len(r.functionsToHoist) - 1; i >= 0; i-- {
			functionExpr := r.functionsToHoist[i]
			scopeRoot := t.GetClosestHoistScopeRoot()
			insertionPoint := scopeRoot.GetFirstChild()
			if insertionPoint.IsNil() ||
				!(insertionPoint.IsVar() && insertionPoint.GetFirstChild().GetString() == moduleName) {
				insertionPoint = rhino.NilNode
			}

			if insertionPoint.IsNil() {
				if scopeRoot.GetFirstChild() != functionExpr {
					scopeRoot.AddChildToFront(functionExpr.Detach())
				}
			} else if insertionPoint != functionExpr && insertionPoint.GetNext() != functionExpr {
				functionExpr.Detach().InsertAfter(insertionPoint)
			}
		}

		for i := range r.exports {
			r.visitExport(t, r.exports[i])
		}

		for _, require := range r.imports {
			r.visitRequireCall(t, require)
		}
	case rhino.Token_CALL:
		if r.outer.isCommonJsImport(n) {
			r.imports = append(r.imports, n)
		}
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		if n.HasMoreThanOneChild() && !NodeUtil_IsAnyFor(parent) && !parent.IsExport() {
			vars := splitMultipleDeclarations(n)
			t.ReportCodeChange()
			for _, varNode := range vars {
				r.Visit(t, varNode.GetFirstChild(), varNode)
			}
		}

		if processCommonJSModules_isSelfAlias(ast.MustNameDeclaration(n)) {
			n.Detach()
			t.ReportCodeChange()
			return
		}
	case rhino.Token_NAME:
		if NodeUtil_IsNameDeclaration(n.GetParent()) &&
			n.GetParent().HasMoreThanOneChild() &&
			!NodeUtil_IsAnyFor(n.GetGrandparent()) {
			break
		}

		qName, ok := n.GetQualifiedNameString()
		if !ok {
			break
		}
		nameDeclaration := t.GetScope().GetVar(qName)
		if nameDeclaration != nil {
			if NodeUtil_IsLhsByDestructuring(n) {
				r.maybeUpdateName(t, n, nameDeclaration)
			} else if !nameDeclaration.GetNode().IsNil() &&
				sameInputId(nameDeclaration.GetNode().GetInputId(), n.GetInputId()) {
				declNameNode := nameDeclaration.GetNameNode()
				enclosingDeclaration := NodeUtil_GetEnclosingNode(
					n, func(node rhino.Node) bool { return node == declNameNode })

				if enclosingDeclaration.IsNil() ||
					enclosingDeclaration == n ||
					nameDeclaration.GetScope().(*Scope) != t.GetScope() {
					r.maybeUpdateName(t, n, nameDeclaration)
				}
			}
		} else if r.allowFullRewrite &&
			n.GetString() == processCommonJSModules_MODULE &&
			!(n.GetParent().IsGetProp() || n.GetParent().IsGetElem() || n.GetParent().IsTypeOf()) {
			objectLit := rhino.IR_Objectlit().Srcref(n)
			n.ReplaceWith(objectLit)
			t.ReportCodeChangeWithNode(objectLit)
		}
	case rhino.Token_GETPROP:
		if n.MatchesQualifiedName(processCommonJSModules_MODULE+".id") ||
			n.MatchesQualifiedName(processCommonJSModules_MODULE+".filename") {
			v := t.GetScope().GetVar(processCommonJSModules_MODULE)
			if v == nil || v.IsExtern() {
				n.ReplaceWith(rhino.IR_String(t.GetInput().GetPath().String()).Srcref(n))
			}
		} else if object, ok := ast.MustGetProp(n).Object().AsName(); r.allowFullRewrite &&
			ok && object.Ident() == processCommonJSModules_MODULE &&
			!n.MatchesQualifiedName(processCommonJSModules_MODULE+"."+processCommonJSModules_EXPORTS) {
			v := t.GetScope().GetVar(processCommonJSModules_MODULE)
			if v == nil || v.IsExtern() {
				object.Node().ReplaceWith(rhino.IR_Objectlit().Srcref(object.Node()))
			}
		}
	case rhino.Token_TYPEOF:
		if operand, ok := ast.MustUnary(n).Operand().AsName(); r.allowFullRewrite &&
			ok && (operand.Ident() == processCommonJSModules_MODULE ||
			operand.Ident() == processCommonJSModules_EXPORTS) {
			v := t.GetScope().GetVar(operand.Ident())
			if v == nil || v.IsExtern() {
				n.ReplaceWith(rhino.IR_String("object"))
			}
		}
	}

	r.fixTypeAnnotationsForNode(t, n)
}

func processCommonJSModules_isSelfAlias(declaration ast.NameDeclaration) bool {
	first, _ := declaration.FirstDeclarator()
	name, ok := first.AsName()
	if !ok {
		return false
	}
	init, ok := name.Init()
	if !ok {
		return false
	}
	value, ok := init.AsName()
	return ok && value.Ident() == name.Ident()
}

func (r *commonJsRewriteModule) fixTypeAnnotationsForNode(t *NodeTraversal, n rhino.Node) {
	if info := n.GetJSDocInfo(); info != nil {
		for it := info.TypeExprIter(); it.Next(); {
			r.fixTypeNode(t, it.Node())
		}
	}
}

func (r *commonJsRewriteModule) visitRequireCall(t *NodeTraversal, require rhino.Node) {
	moduleName := r.outer.getImportedModuleName(t, require)
	firstArg, _ := ast.MustCall(require).FirstArg()
	moduleRef := NodeUtil_NewQName(r.outer.compiler, r.outer.getBasePropertyImportWithRequireCall(moduleName, require)).
		SrcrefTree(firstArg.Node())
	require.ReplaceWith(moduleRef)

	t.ReportCodeChange()
}

func (r *commonJsRewriteModule) visitExport(t *NodeTraversal, export commonJsExportInfo) {
	compiler := r.outer.compiler
	root := baseQualifiedNameNode(export.node)
	rValue := NodeUtil_GetRValueOfLValue(root)

	if processCommonJSModules_MODULE_EXPORTS.Matches(root) &&
		!rValue.IsNil() &&
		rValue.IsObjectLit() &&
		root.GetParent().IsAssign() &&
		root.GetGrandparent().IsExprResult() {
		if r.expandObjectLitAssignment(t, root, export.scope) {
			return
		}
	}

	inputModuleName := commonJsModuleName(t.GetInput())
	moduleInitialization := t.GetScope().GetVar(inputModuleName)

	var rValueVar *Var
	if rValQName, ok := rValue.GetQualifiedNameString(); ok {
		rValueVar = export.scope.GetVar(rValQName)
	}

	if root.GetParent().IsAssign() &&
		root.GetGrandparent().IsExprResult() &&
		(!root.GetNext().IsNil() && (root.GetNext().IsName() || root.GetNext().IsGetProp())) &&
		rValueVar != nil &&
		(NodeUtil_GetEnclosingScript(rValueVar.GetNameNode()).IsNil() ||
			(rValueVar.GetNameNode().HasParent() && !rValueVar.IsParam())) &&
		export.isInSupportedScope &&
		(rValueVar.GetNameNode().GetParent().IsNil() ||
			!NodeUtil_IsLhsByDestructuring(rValueVar.GetNameNode())) {
		root.GetGrandparent().Detach()
		t.ReportCodeChange()
		return
	}

	defaultModuleName := inputModuleName + "." + processCommonJSModules_EXPORT_PROPERTY_NAME
	exportNodeQName, _ := export.node.GetQualifiedNameString()
	updatedExport := NodeUtil_NewQNameWithOriginalName(compiler, defaultModuleName, export.node, exportNodeQName)
	updatedExport.PutBooleanProp(rhino.Node_MODULE_EXPORT, true)
	exportIsConst := r.defaultExportIsConst &&
		updatedExport.MatchesQualifiedName(r.outer.getBasePropertyImport(inputModuleName)) &&
		root == export.node &&
		NodeUtil_IsLValue(export.node)

	var changeScope rhino.Node

	if processCommonJSModules_MODULE_EXPORTS.Matches(root) &&
		!rValue.IsNil() &&
		export.scope.GetVar("module.exports") == nil &&
		root.GetParent().IsAssign() &&
		root.GetGrandparent().IsExprResult() &&
		moduleInitialization == nil {
		parent := root.GetParent()
		exportName := rhino.IR_ExprResult(rhino.IR_Assign(updatedExport, rValue.Detach()))
		if exportIsConst {
			info := rhino.JSDocInfo_Builder_()
			info.RecordConstancy()
			exportName.GetFirstChild().SetJSDocInfo(info.Build())
		}
		parent.GetParent().ReplaceWith(exportName.SrcrefTree(parent))
		changeScope = ChangeTracker_GetEnclosingChangeScopeRoot(parent)
	} else if !root.GetNext().IsNil() &&
		root.GetNext().IsName() &&
		rValueVar != nil &&
		rValueVar.IsGlobal() &&
		export.isInSupportedScope &&
		(rValueVar.GetNameNode().GetParent().IsNil() ||
			!NodeUtil_IsLhsByDestructuring(rValueVar.GetNameNode())) {
		parent := root.GetParent()
		root.Detach()
		parent.ReplaceWith(root)
		if root == export.node {
			root = updatedExport
		}
		export.node.ReplaceWith(updatedExport)
		changeScope = ChangeTracker_GetEnclosingChangeScopeRoot(root)
	} else {
		export.node.ReplaceWith(updatedExport)
		if updatedExport.GetParent().IsAssign() && exportIsConst {
			infoBuilder := rhino.JSDocInfo_Builder_MaybeCopyFrom(updatedExport.GetParent().GetJSDocInfo())
			infoBuilder.RecordConstancy()
			updatedExport.GetParent().SetJSDocInfo(infoBuilder.Build())
		}
		changeScope = ChangeTracker_GetEnclosingChangeScopeRoot(updatedExport)
	}

	if !changeScope.IsNil() {
		compiler.ReportChangeToChangeScope(changeScope)
	}
}

func (r *commonJsRewriteModule) expandObjectLitAssignment(t *NodeTraversal, export rhino.Node, scope *Scope) bool {
	jutil.CheckState(export.GetParent().IsAssign())
	insertionRef := export.GetGrandparent()
	jutil.CheckState(insertionRef.IsExprResult())
	rhino.CheckNotNull(insertionRef.GetParent())

	rValue := ast.MustObjectLit(NodeUtil_GetRValueOfLValue(export))
	member, ok := rValue.FirstMember()

	removedNodes := false
	for ok {
		if !processCommonJSModules_isNamespaceMember(member) {
			member, ok = member.Next()
			continue
		}
		key := member.Node()

		lhs := rhino.IR_Getprop(export.CloneTree(), key.GetString())
		var value rhino.Node
		if stringKey, ok := member.AsStringKey(); ok {
			value = stringKey.Value().Detach().Node()
		} else {
			value = ast.MustMemberFunctionDef(key).Function().Detach().Node()
		}

		expr := rhino.IR_ExprResult(rhino.IR_Assign(lhs, value)).SrcrefTreeIfMissing(key)
		expr.InsertAfter(insertionRef)
		newExport := newCommonJsExportInfo(lhs.GetFirstChild(), scope)
		r.visitExport(t, newExport)

		if !expr.IsNil() && expr.HasParent() {
			insertionRef = expr
		}

		current := member
		member, ok = member.Next()
		current.Detach()
		removedNodes = true
	}

	if !rValue.HasMembers() {
		export.GetGrandparent().Detach()
		return true
	}

	if removedNodes {
		t.ReportCodeChangeWithNode(rValue.Node())
	}
	return false
}

func (r *commonJsRewriteModule) maybeUpdateName(t *NodeTraversal, n rhino.Node, var_ *Var) {
	jutil.CheckNotNull(var_)
	jutil.CheckState(n.IsName() || n.IsGetProp())
	jutil.CheckState(n.HasParent())
	importedModuleName, isImportAlias := r.getModuleImportName(t, var_.GetNode())
	name, _ := n.GetQualifiedNameString()

	if isImportAlias && n != var_.GetNode() {
		r.updateNameReference(t, n, name, importedModuleName, false, false)
	} else if r.allowFullRewrite {
		exportedName, hasExportedName := r.getExportedName(t, n, var_)

		if (n != var_.GetNode() || n.GetParent().IsClass()) &&
			!hasExportedName &&
			(var_ == nil ||
				var_.GetNameNode().GetParent().IsNil() ||
				!NodeUtil_IsLhsByDestructuring(var_.GetNameNode())) {
			if n.GetParent().IsClass() && n.GetParent().GetFirstChild() == n {
				r.rewrittenClassExpressions = append(r.rewrittenClassExpressions, n.GetParent())
			}
			return
		}

		if !isImportAlias &&
			hasExportedName &&
			exportedName != name &&
			!var_.IsParam() &&
			(var_.GetNameNode().GetParent().IsNil() ||
				!NodeUtil_IsLhsByDestructuring(var_.GetNameNode())) {
			exportPropIsConst := r.defaultExportIsConst &&
				r.outer.getBasePropertyImport(commonJsModuleName(t.GetInput())) == exportedName &&
				baseQualifiedNameNode(n) == n &&
				NodeUtil_IsLValue(n)
			r.updateNameReference(t, n, name, exportedName, true, exportPropIsConst)
		} else if var_.IsGlobal() {
			currentModuleName := commonJsModuleName(t.GetInput())
			if currentModuleName == name || processCommonJSModules_EXPORTS == name {
				return
			}
			if r.outer.compiler.GetOptions().ShouldExportTestFunctions() &&
				strings.HasPrefix(currentModuleName, "test") {
				return
			}
			newName := name + "$$" + currentModuleName
			r.updateNameReference(t, n, name, newName, false, false)
		}
	}
}

func (r *commonJsRewriteModule) updateNameReference(
	t *NodeTraversal,
	nameRef rhino.Node,
	originalName string,
	newName string,
	requireFunctionExpressions bool,
	qualifiedNameIsConst bool) {
	compiler := r.outer.compiler
	parent := nameRef.GetParent()
	rhino.CheckNotNull(parent)
	newNameIsQualified := strings.IndexByte(newName, '.') >= 0
	newNameIsModuleExport := newName == r.outer.getBasePropertyImport(commonJsModuleName(t.GetInput()))

	newNameDeclaration := t.GetScope().GetVar(newName)

	switch parent.GetToken() {
	case rhino.Token_CLASS:
		cls := ast.MustClass(parent)
		className, hasName := cls.Name()
		superClass, hasExtends := cls.Extends()
		if hasName && className.Is(nameRef) &&
			(newNameIsQualified || requireFunctionExpressions) {
			r.rewrittenClassExpressions = append(r.rewrittenClassExpressions, parent)

			newNameRef := NodeUtil_NewQNameWithOriginalName(compiler, newName, nameRef, originalName)
			if newNameIsModuleExport {
				newNameRef.PutBooleanProp(rhino.Node_MODULE_EXPORT, true)
			}

			var expr rhino.Node
			if !newNameIsQualified && newNameDeclaration == nil {
				expr = rhino.IR_LetWithValue(newNameRef, rhino.IR_NullNode()).SrcrefTreeIfMissing(nameRef)
			} else {
				expr = rhino.IR_ExprResult(rhino.IR_Assign(newNameRef, rhino.IR_NullNode())).SrcrefTreeIfMissing(nameRef)
				info := rhino.JSDocInfo_Builder_MaybeCopyFrom(parent.GetJSDocInfo())
				parent.SetJSDocInfo(nil)
				if qualifiedNameIsConst {
					info.RecordConstancy()
				}
				expr.GetFirstChild().SetJSDocInfo(info.Build())
				r.fixTypeAnnotationsForNode(t, expr.GetFirstChild())
			}
			parent.ReplaceWith(expr)
			if expr.IsLet() {
				expr.GetFirstFirstChild().ReplaceWith(parent)
			} else {
				expr.GetFirstChild().GetSecondChild().ReplaceWith(parent)
			}
		} else if hasExtends && superClass.Is(nameRef) {
			newNameRef := NodeUtil_NewQNameWithOriginalName(compiler, newName, nameRef, originalName)
			if newNameIsModuleExport {
				newNameRef.PutBooleanProp(rhino.Node_MODULE_EXPORT, true)
			}
			nameRef.ReplaceWith(newNameRef)
		} else {
			nameRef.SetString(newName)
			nameRef.SetOriginalName(&originalName)
		}
	case rhino.Token_FUNCTION:
		if newNameIsQualified || requireFunctionExpressions {
			if NodeUtil_IsFunctionExpression(parent) {
				return
			}
			newNameRef := NodeUtil_NewQNameWithOriginalName(compiler, newName, nameRef, originalName)
			if newNameIsModuleExport {
				newNameRef.PutBooleanProp(rhino.Node_MODULE_EXPORT, true)
			}
			nameRef.SetString("")

			var expr rhino.Node
			if !newNameIsQualified && newNameDeclaration == nil {
				expr = rhino.IR_VarWithValue(newNameRef, rhino.IR_NullNode()).SrcrefTreeIfMissing(nameRef)
			} else {
				expr = rhino.IR_ExprResult(rhino.IR_Assign(newNameRef, rhino.IR_NullNode())).SrcrefTreeIfMissing(nameRef)
			}
			parent.ReplaceWith(expr)
			if expr.IsVar() {
				expr.GetFirstFirstChild().ReplaceWith(parent)
			} else {
				expr.GetFirstChild().GetSecondChild().ReplaceWith(parent)
				info := rhino.JSDocInfo_Builder_MaybeCopyFrom(parent.GetJSDocInfo())
				parent.SetJSDocInfo(nil)
				if qualifiedNameIsConst {
					info.RecordConstancy()
				}
				expr.GetFirstChild().SetJSDocInfo(info.Build())
				r.fixTypeAnnotationsForNode(t, expr.GetFirstChild())
			}
			r.functionsToHoist = append(r.functionsToHoist, expr)
		} else {
			nameRef.SetString(newName)
			nameRef.SetOriginalName(&originalName)
		}
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		if parent.HasMoreThanOneChild() && !NodeUtil_IsAnyFor(parent.GetParent()) {
			splitMultipleDeclarations(parent)
			parent = nameRef.GetParent()
			newNameDeclaration = t.GetScope().GetVar(newName)
		}

		init, hasInit := ast.MustName(nameRef).Init()
		if newNameIsQualified {
			if !hasInit && parent.GetJSDocInfo() == nil {
				parent.Detach()
				break
			}

			getProp := NodeUtil_NewQNameWithOriginalName(compiler, newName, nameRef, originalName)
			if newNameIsModuleExport {
				getProp.PutBooleanProp(rhino.Node_MODULE_EXPORT, true)
			}
			info := parent.GetJSDocInfo()
			parent.SetJSDocInfo(nil)
			if hasInit {
				assign := rhino.IR_Assign(getProp, init.Detach().Node())
				expr := rhino.IR_ExprResult(assign).SrcrefTreeIfMissing(nameRef)
				parent.ReplaceWith(expr)
				infoBuilder := rhino.JSDocInfo_Builder_MaybeCopyFrom(info)
				if qualifiedNameIsConst {
					infoBuilder.RecordConstancy()
				}
				assign.SetJSDocInfo(infoBuilder.Build())
				r.fixTypeAnnotationsForNode(t, assign)
			} else {
				getProp.SetJSDocInfo(info)
				parent.ReplaceWith(rhino.IR_ExprResult(getProp).Srcref(getProp))
			}
		} else if newNameDeclaration != nil && newNameDeclaration.GetNameNode() != nameRef {
			if !hasInit {
				parent.Detach()
				break
			}

			name := NodeUtil_NewNameWithOriginalName(compiler, newName, nameRef, originalName)
			assign := rhino.IR_Assign(name, init.Detach().Node())
			if info := parent.GetJSDocInfo(); info != nil {
				parent.SetJSDocInfo(nil)
				assign.SetJSDocInfo(info)
			}
			parent.ReplaceWith(rhino.IR_ExprResult(assign).SrcrefTree(nameRef))
		} else {
			nameRef.SetString(newName)
			nameRef.SetOriginalName(&originalName)
		}
	default:
		if !newNameIsQualified && nameRef.IsName() {
			nameRef.SetString(newName)
			nameRef.SetOriginalName(&originalName)
		} else {
			var name rhino.Node
			if newNameIsQualified {
				name = NodeUtil_NewQNameWithOriginalName(compiler, newName, nameRef, originalName)
			} else {
				name = NodeUtil_NewNameWithOriginalName(compiler, newName, nameRef, originalName)
			}

			if newNameIsModuleExport {
				name.PutBooleanProp(rhino.Node_MODULE_EXPORT, true)
			}
			if info := nameRef.GetJSDocInfo(); info != nil {
				nameRef.SetJSDocInfo(nil)
				name.SetJSDocInfo(info)
			}
			name.SrcrefTree(nameRef)
			nameRef.ReplaceWith(name)
			if nameRef.HasChildren() {
				name.AddChildrenToFront(nameRef.RemoveChildren())
			}
		}
	}

	t.ReportCodeChange()
}

func (r *commonJsRewriteModule) getExportedName(t *NodeTraversal, n rhino.Node, var_ *Var) (string, bool) {
	if var_ == nil || !sameInputId(var_.GetNode().GetInputId(), n.GetInputId()) {
		return n.GetQualifiedNameString()
	}

	baseExportName := r.outer.getBasePropertyImport(commonJsModuleName(t.GetInput()))

	for i := range r.exports {
		export := r.exports[i]
		exportBase := baseQualifiedNameNode(export.node)
		exportRValue := NodeUtil_GetRValueOfLValue(exportBase)
		if exportRValue.IsNil() {
			continue
		}

		exportedName := getExportedNameNode(export)
		if export.isInSupportedScope &&
			(exportRValue == n ||
				((NodeUtil_IsClassExpression(exportRValue) || NodeUtil_IsFunctionExpression(exportRValue)) &&
					exportedName == n)) {
			return "", false
		}

		exportBaseQName, hasExportBaseQName := exportBase.GetQualifiedNameString()

		if exportRValue.IsObjectLit() {
			if !hasExportBaseQName || "module.exports" != exportBaseQName {
				return n.GetQualifiedNameString()
			}

			featureSet := r.outer.compiler.GetOptions().GetLanguageIn().ToFeatureSet()
			for member := range ast.MustObjectLit(exportRValue).Members() {
				key, ok := member.AsStringKey()
				if !ok || key.IsQuoted() || !NodeUtil_IsValidPropertyName(featureSet, key.Key()) {
					continue
				}
				value := key.Value().Node()
				valQName, isQName := value.GetQualifiedNameString()
				if !isQName {
					continue
				}
				if value == n {
					return "", false
				}

				valVar := t.GetScope().GetVar(valQName)
				if valVar != nil && valVar.GetNameNode() == var_.GetNameNode() {
					if export.isInSupportedScope {
						return baseExportName + "." + key.Key(), true
					}
					return n.GetQualifiedNameString()
				}
			}
		} else {
			if !export.isInSupportedScope {
				return n.GetQualifiedNameString()
			}
			if var_.GetNameNode() == exportedName {
				exportPrefix := processCommonJSModules_EXPORTS
				if strings.HasPrefix(exportBaseQName, processCommonJSModules_MODULE) {
					exportPrefix = processCommonJSModules_MODULE + "." + processCommonJSModules_EXPORTS
				}
				if len(exportBaseQName) == len(exportPrefix) {
					return baseExportName, true
				}
				return baseExportName + exportBaseQName[len(exportPrefix):], true
			}
		}
	}
	return n.GetQualifiedNameString()
}

func getExportedNameNode(info commonJsExportInfo) rhino.Node {
	qNameBase := baseQualifiedNameNode(info.node)
	rValue := NodeUtil_GetRValueOfLValue(qNameBase)
	if rValue.IsNil() {
		return rhino.NilNode
	}
	if NodeUtil_IsFunctionExpression(rValue) {
		return rValue.GetFirstChild()
	}
	if NodeUtil_IsClassExpression(rValue) && rValue.GetFirstChild() == qNameBase {
		return rValue.GetFirstChild()
	}
	if !rValue.IsName() {
		return rhino.NilNode
	}
	var_ := info.scope.GetVar(rValue.GetString())
	if var_ == nil {
		return rhino.NilNode
	}
	return var_.GetNameNode()
}

func (r *commonJsRewriteModule) getModuleImportName(t *NodeTraversal, n rhino.Node) (string, bool) {
	var rValue rhino.Node
	propSuffix := ""
	parent := n.GetParent()
	if key, ok := ast.AsStringKey(parent); ok {
		if lhs, ok := ast.AsDestructuringLhs(parent.GetGrandparent()); ok && lhs.Pattern().Is(parent.GetParent()) {
			if _, isObjectPattern := lhs.Pattern().AsObjectPattern(); isObjectPattern {
				init, _ := lhs.Init()
				rValue = init.Node()
				propSuffix = "." + key.Key()
			}
		}
	}
	if propSuffix == "" && !parent.IsNil() {
		rValue = NodeUtil_GetRValueOfLValue(n)
	}

	if rValue.IsNil() {
		return "", false
	}

	if rValue.IsCall() && r.outer.isCommonJsImport(rValue) {
		return r.outer.getBasePropertyImportWithRequireCall(r.outer.getImportedModuleName(t, rValue), rValue) + propSuffix, true
	}
	if getProp, ok := ast.AsGetProp(rValue); ok && r.outer.isCommonJsImport(getProp.Object().Node()) {
		require := getProp.Object().Node()
		importName := r.outer.getBasePropertyImportWithRequireCall(r.outer.getImportedModuleName(t, require), require)
		return importName + "." + getProp.Property() + propSuffix, true
	}

	return "", false
}

func (r *commonJsRewriteModule) fixTypeNode(t *NodeTraversal, typeNode rhino.Node) {
	if typeNode.IsStringLit() {
		name := typeNode.GetString()
		if ModuleLoader_IsPathIdentifier(name) {
			from := max(strings.LastIndexByte(name, '/'), 0)
			moduleName := name
			localTypeName := ""
			if rel := strings.IndexByte(name[from:], '.'); rel >= 0 {
				endIndex := from + rel
				moduleName = name[:endIndex]
				localTypeName = name[endIndex:]
			}

			globalModuleName := r.outer.getImportedModuleNameWithImportPath(t, typeNode, moduleName)
			baseImportProperty := r.outer.getBasePropertyImport(globalModuleName)
			typeNode.SetString(baseImportProperty + localTypeName)
		} else {
			wasRewritten := false
			typeInputId := typeNode.GetInputId()
			scope := t.GetScope()
			for start := 0; start <= len(name); {
				endIndex := len(name)
				if rel := strings.IndexByte(name[start:], '.'); rel >= 0 {
					endIndex = start + rel
				}
				baseName := name[:endIndex]
				suffix := name[endIndex:]
				typeDeclaration := scope.GetVar(baseName)

				if typeDeclaration != nil && !typeDeclaration.GetNode().IsNil() &&
					sameInputId(typeDeclaration.GetNode().GetInputId(), typeInputId) {
					if importedModuleName, ok := r.getModuleImportName(t, typeDeclaration.GetNode()); ok {
						typeNode.SetString(importedModuleName + suffix)
						typeNode.SetOriginalName(&name)
						wasRewritten = true
						break
					} else if r.allowFullRewrite {
						if exportedName, ok := r.getExportedName(t, typeNode, typeDeclaration); ok && exportedName != name {
							typeNode.SetString(exportedName + suffix)
							typeNode.SetOriginalName(&name)
							wasRewritten = true
							break
						}
					}
				}
				if endIndex == len(name) {
					break
				}
				start = endIndex + 1
			}

			if !wasRewritten && r.allowFullRewrite {
				baseName, rest, hasDot := strings.Cut(name, ".")
				if typeDeclaration := scope.GetVar(baseName); typeDeclaration != nil && typeDeclaration.IsGlobal() {
					moduleName := commonJsModuleName(t.GetInput())
					newName := baseName + "$$" + moduleName
					if hasDot {
						newName += "." + rest
					}
					typeNode.SetString(newName)
					typeNode.SetOriginalName(&name)
				}
			}
		}
	}

	for child := typeNode.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		r.fixTypeNode(t, child)
	}
}

func splitMultipleDeclarations(varNode rhino.Node) []rhino.Node {
	declaration, ok := ast.AsNameDeclaration(varNode)
	jutil.CheckState(ok)
	var vars []rhino.Node
	info := varNode.GetJSDocInfo()
	for first, _ := declaration.FirstDeclarator(); !first.IsLast(); first, _ = declaration.FirstDeclarator() {
		newVar := rhino.NewNodeWithChild(varNode.GetToken(), first.Detach().Node())
		if info != nil {
			newVar.SetJSDocInfo(info.Clone())
		}
		newVar.Srcref(varNode)
		newVar.InsertBefore(varNode)
		vars = append(vars, newVar)
	}
	vars = append(vars, varNode)
	return vars
}
