// Port of com/google/javascript/jscomp/Es6RewriteBlockScopedDeclaration.java @ closure-compiler v20261021

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
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Es6RewriteBlockScopedDeclaration rewrites "let"s and "const"s as "var"s,
// renaming block-scoped declarations and their references when captured by
// loop closures.
//
// Note that this must run after Es6RewriteDestructuring, since it does not
// process destructuring let/const declarations at all.
type Es6RewriteBlockScopedDeclaration struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler          AbstractCompiler
	astFactory        *AstFactory
	letConsts         []rhino.Node
	letConstActive    map[rhino.Node]bool
	loopLetConstNames map[string]struct{}
	uniqueIdSupplier  *UniqueIdSupplier
}

func (e *Es6RewriteBlockScopedDeclaration) noVisitLeaves() {}

var es6RewriteBlockScopedDeclarationTranspiledFeatures = parser.FeatureSet_BARE_MINIMUM.WithFeatures(
	parser.FeatureSet_Feature_LET_DECLARATIONS,
	parser.FeatureSet_Feature_CONST_DECLARATIONS,
)

const (
	es6BlockLoopParamNamePrefix = "$jscomp$loop_param$"
	es6BlockLoopObjectName      = "$jscomp$loop"
)

func NewEs6RewriteBlockScopedDeclaration(compiler AbstractCompiler) *Es6RewriteBlockScopedDeclaration {
	return &Es6RewriteBlockScopedDeclaration{
		compiler:         compiler,
		uniqueIdSupplier: compiler.GetUniqueIdSupplier(),
		astFactory:       compiler.CreateAstFactory(),
	}
}

func (e *Es6RewriteBlockScopedDeclaration) addLetConst(n rhino.Node) {
	if e.letConstActive == nil {
		e.letConstActive = make(map[rhino.Node]bool, 16)
	}
	if _, seen := e.letConstActive[n]; !seen {
		e.letConsts = append(e.letConsts, n)
	}
	e.letConstActive[n] = true
}

func (e *Es6RewriteBlockScopedDeclaration) removeLetConst(n rhino.Node) {
	if e.letConstActive != nil {
		e.letConstActive[n] = false
	}
}

func (e *Es6RewriteBlockScopedDeclaration) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	jutil.CheckStateMsg1(parent.IsNil() || !parent.IsForOf(), "%s", parent)

	if n.IsLet() || n.IsConst() {
		e.addLetConst(n)
		for nameNode := n.GetFirstChild(); !nameNode.IsNil(); nameNode = nameNode.GetNext() {
			e.visitBlockScopedNameDeclaration(n, nameNode)
		}
	}
}

func (e *Es6RewriteBlockScopedDeclaration) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(e.compiler, root, e)
	if len(e.letConsts) > 0 {
		if len(e.loopLetConstNames) > 0 {
			sc := NewSyntacticScopeCreator(e.compiler)
			transformer := &es6BlockLoopClosureTransformer{pass: e}
			NodeTraversal_Builder_().SetCompiler(e.compiler).SetScopeCreator(sc).SetCallback(transformer).Traverse(root)
			transformer.transformLoopClosure()
		}
		e.rewriteDeclsToVars()
	}
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(e.compiler, root, es6RewriteBlockScopedDeclarationTranspiledFeatures)
}

func (e *Es6RewriteBlockScopedDeclaration) visitBlockScopedNameDeclaration(decl rhino.Node, nameNode rhino.Node) {
	if !e.inLoop(decl) {
		return
	}
	if nm, ok := ast.AsName(nameNode); ok {
		if e.loopLetConstNames == nil {
			e.loopLetConstNames = make(map[string]struct{}, 8)
		}
		e.loopLetConstNames[nm.Ident()] = struct{}{}
	}
	parent := decl.GetParent()
	// Normalize "let x;" to "let x = undefined;" if in a loop, since we later convert x
	// to $jscomp$loop$0.x and want to reset the property to undefined every loop iteration.
	if !nameNode.HasChildren() && (parent.IsNil() || !parent.IsForIn()) {
		undefined := e.astFactory.CreateUndefinedValue().SrcrefTree(nameNode)
		nameNode.AddChildToFront(undefined)
		e.compiler.ReportChangeToEnclosingScope(undefined)
	}
}

func (e *Es6RewriteBlockScopedDeclaration) inLoop(n rhino.Node) bool {
	enclosingNode := NodeUtil_GetEnclosingNode(n, func(curr rhino.Node) bool {
		_, isFn := ast.AsFunction(curr)
		return isFn || NodeUtil_IsLoopStructure(curr)
	})
	return !enclosingNode.IsNil() && !enclosingNode.IsFunction()
}

func es6BlockExtractInlineJSDoc(srcDeclaration rhino.Node, srcName rhino.Node, destDeclaration rhino.Node) {
	existingInfo := srcDeclaration.GetJSDocInfo()
	if existingInfo == nil {
		existingInfo = srcName.GetJSDocInfo()
		srcName.SetJSDocInfo(nil)
	}
	builder := rhino.JSDocInfo_Builder_MaybeCopyFrom(existingInfo)
	destDeclaration.SetJSDocInfo(builder.Build())
}

func es6BlockMaybeAddConstJSDoc(srcDeclaration rhino.Node, srcName rhino.Node, destDeclaration rhino.Node) {
	if srcDeclaration.IsConst() {
		es6BlockExtractInlineJSDoc(srcDeclaration, srcName, destDeclaration)
		builder := rhino.JSDocInfo_Builder_MaybeCopyFrom(destDeclaration.GetJSDocInfo())
		builder.RecordConstancy()
		destDeclaration.SetJSDocInfo(builder.Build())
	}
}

func (e *Es6RewriteBlockScopedDeclaration) handleLetConstDeclarationList(declarationList rhino.Node, parent rhino.Node) {
	if declarationList.IsVar() {
		jutil.CheckStateMsg(
			!declarationList.HasChildren() || declarationList.HasOneChild(),
			"var declaration can have maximum one child post normalization.")
		return
	}

	if forLoop, ok := ast.AsFor(parent); ok {
		e.handleLetConstInVanillaForInitializer(declarationList, forLoop)
		return
	}

	if forIn, ok := ast.AsForIn(parent); ok {
		e.handleLetConstInForInInitializer(declarationList, forIn)
		return
	}

	for declarationList.HasChildren() {
		name := declarationList.GetLastChild()
		newDeclaration := rhino.IR_Var(name.Detach()).Srcref(declarationList)
		es6BlockExtractInlineJSDoc(declarationList, name, newDeclaration)
		es6BlockMaybeAddConstJSDoc(declarationList, name, newDeclaration)
		newDeclaration.InsertAfter(declarationList)
		e.compiler.ReportChangeToEnclosingScope(parent)
	}

	declarationList.Detach()
	e.compiler.ReportChangeToEnclosingScope(parent)
}

func (e *Es6RewriteBlockScopedDeclaration) handleLetConstInForInInitializer(declarationList rhino.Node, forIn ast.ForIn) {
	loopNode := forIn.Node()
	lhs := declarationList.GetFirstChild()
	if dlhs, ok := ast.AsDestructuringLhs(lhs); ok {
		jutil.CheckStateMsg(
			false, "Destructuring syntax is unsupported in ES6RewriteBlockScopedDeclarations pass")
		NodeUtil_VisitLhsNodesInNode(
			lhs,
			func(name rhino.Node) {
				nm, ok := ast.AsName(name)
				jutil.CheckStateMsg1(ok, "lhs in destructuring declaration should be a simple name. (%s)", name)
				newVar := rhino.IR_Var(rhino.IR_Name(nm.Ident()).Srcref(name)).Srcref(name)
				es6BlockExtractInlineJSDoc(declarationList, name, newVar)
				e.addNodeBeforeLoop(newVar, loopNode)
			})
		declarationList.ReplaceWith(dlhs.Pattern().Node().Detach())
	} else {
		newStatement := declarationList.CloneTree()
		name := newStatement.GetFirstChild().CloneNode()
		newStatement.SetToken(rhino.Token_VAR)
		es6BlockExtractInlineJSDoc(declarationList, lhs, newStatement)
		declarationList.ReplaceWith(name)
		e.addNodeBeforeLoop(newStatement, loopNode)
	}
}

func (e *Es6RewriteBlockScopedDeclaration) handleLetConstInVanillaForInitializer(declarationList rhino.Node, forLoop ast.For) {
	parent := forLoop.Node()
	insertSpot := e.insertSpotBeforeLoop(parent)
	for declarationList.HasChildren() {
		name := declarationList.GetLastChild()
		newDeclaration := rhino.IR_Var(name.Detach()).Srcref(declarationList)
		es6BlockExtractInlineJSDoc(declarationList, name, newDeclaration)
		newDeclaration.InsertBefore(insertSpot)
		insertSpot = newDeclaration
		e.compiler.ReportChangeToEnclosingScope(parent)
	}
	empty := e.astFactory.CreateEmpty().Srcref(declarationList)
	declarationList.ReplaceWith(empty)
	e.compiler.ReportChangeToEnclosingScope(empty)
}

func (e *Es6RewriteBlockScopedDeclaration) addNodeBeforeLoop(newNode rhino.Node, loopNode rhino.Node) {
	insertSpot := e.insertSpotBeforeLoop(loopNode)
	newNode.InsertBefore(insertSpot)
	e.compiler.ReportChangeToEnclosingScope(newNode)
}

func (e *Es6RewriteBlockScopedDeclaration) insertSpotBeforeLoop(loopNode rhino.Node) rhino.Node {
	insertSpot := loopNode
	for _, ok := ast.AsLabel(insertSpot.GetParent()); ok; _, ok = ast.AsLabel(insertSpot.GetParent()) {
		insertSpot = insertSpot.GetParent()
	}
	return insertSpot
}

func (e *Es6RewriteBlockScopedDeclaration) rewriteDeclsToVars() {
	for _, n := range e.letConsts {
		if e.letConstActive[n] {
			e.handleLetConstDeclarationList(n, n.GetParent())
		}
	}
}

type es6LoopClosureHandledKey struct {
	node rhino.Node
	name string
}

type es6BlockLoopObject struct {
	name string
	vars []*Var
}

type es6WrapTarget struct {
	node    rhino.Node
	objects []*es6BlockLoopObject
}

// es6BlockLoopClosureTransformer transforms let/const declarations captured by loop closures.
type es6BlockLoopClosureTransformer struct {
	NodeTraversal_AbstractPreOrderCallbackBase
	pass *Es6RewriteBlockScopedDeclaration

	loopOrder                        []rhino.Node
	loopObjects                      map[rhino.Node]*es6BlockLoopObject
	wrapTargets                      []es6WrapTarget
	wrapTargetIdx                    map[rhino.Node]int
	nodesHandledForLoopObjectClosure map[es6LoopClosureHandledKey]struct{}
	references                       map[*Var][]rhino.Node
}

func (l *es6BlockLoopClosureTransformer) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if !NodeUtil_IsReferenceName(n) {
		return true
	}

	name := n.GetString()
	if _, ok := l.pass.loopLetConstNames[name]; !ok {
		return true
	}
	referencedIn := t.GetScope()
	v := referencedIn.GetVar(name)
	if v == nil || (!v.IsLet() && !v.IsConst()) {
		return true
	}

	declaredIn := v.GetScope().(*Scope)
	var loopNode rhino.Node
	for s := declaredIn; ; s = s.GetParent() {
		scopeRoot := s.GetRootNode()
		if NodeUtil_IsLoopStructure(scopeRoot) {
			loopNode = scopeRoot
			break
		} else if scopeRoot.HasParent() && NodeUtil_IsLoopStructure(scopeRoot.GetParent()) {
			loopNode = scopeRoot.GetParent()
			break
		} else if s.IsFunctionBlockScope() || s.IsGlobal() {
			return true
		}
	}

	if l.references == nil {
		l.references = make(map[*Var][]rhino.Node, 8)
	}
	l.references[v] = append(l.references[v], n)

	var outerMostFunctionScope *Scope
	for s := referencedIn; s != declaredIn && s.GetRootNode() != loopNode; s = s.GetParent() {
		if s.IsFunctionScope() {
			outerMostFunctionScope = s
		}
	}

	if outerMostFunctionScope != nil {
		enclosingFunction := outerMostFunctionScope.GetRootNode()
		nodeToWrapInClosure := enclosingFunction
		if _, isAccessor := ast.AsGetterOrSetterDef(enclosingFunction.GetParent()); isAccessor {
			nodeToWrapInClosure = enclosingFunction.GetGrandparent()
			jutil.CheckState(nodeToWrapInClosure.IsObjectLit())
		}

		key := es6LoopClosureHandledKey{node: nodeToWrapInClosure, name: name}
		if _, ok := l.nodesHandledForLoopObjectClosure[key]; ok {
			return true
		}
		if l.nodesHandledForLoopObjectClosure == nil {
			l.nodesHandledForLoopObjectClosure = make(map[es6LoopClosureHandledKey]struct{}, 8)
		}
		l.nodesHandledForLoopObjectClosure[key] = struct{}{}

		if l.loopObjects == nil {
			l.loopObjects = make(map[rhino.Node]*es6BlockLoopObject, 4)
		}
		object, ok := l.loopObjects[loopNode]
		if !ok {
			object = &es6BlockLoopObject{
				name: es6BlockLoopObjectName + "$" + l.pass.uniqueIdSupplier.GetUniqueId(t.GetInput()),
			}
			l.loopObjects[loopNode] = object
			l.loopOrder = append(l.loopOrder, loopNode)
		}
		if !slices.Contains(object.vars, v) {
			object.vars = append(object.vars, v)
		}

		if l.wrapTargetIdx == nil {
			l.wrapTargetIdx = make(map[rhino.Node]int, 4)
		}
		if idx, exists := l.wrapTargetIdx[nodeToWrapInClosure]; exists {
			wt := &l.wrapTargets[idx]
			if !slices.Contains(wt.objects, object) {
				wt.objects = append(wt.objects, object)
			}
		} else {
			l.wrapTargetIdx[nodeToWrapInClosure] = len(l.wrapTargets)
			l.wrapTargets = append(l.wrapTargets, es6WrapTarget{
				node:    nodeToWrapInClosure,
				objects: []*es6BlockLoopObject{object},
			})
		}
	}
	return true
}

func (l *es6BlockLoopClosureTransformer) transformLoopClosure() {
	if len(l.loopOrder) == 0 {
		return
	}

	wrapperFunctions := l.createWrapperFunctions()

	for _, loopNode := range l.loopOrder {
		loopObject := l.loopObjects[loopNode]
		objectLitNextIteration := l.pass.astFactory.CreateObjectLit()
		l.renameVarsToProperties(loopObject, objectLitNextIteration, loopNode)

		updateLoopObject := l.pass.astFactory.CreateAssign(
			l.createLoopObjectNameNode(loopObject), objectLitNextIteration)
		objectLit := rhino.IR_VarWithValue(
			l.createLoopObjectNameNode(loopObject), l.pass.astFactory.CreateObjectLit()).
			SrcrefTree(loopNode)
		l.pass.addNodeBeforeLoop(objectLit, loopNode)
		if forLoop, ok := ast.AsFor(loopNode); ok {
			l.changeVanillaForLoopHeader(forLoop, updateLoopObject)
		} else {
			loopBody := NodeUtil_GetLoopCodeBlock(loopNode)
			loopBody.AddChildToFront(rhino.IR_ExprResult(updateLoopObject).SrcrefTreeIfMissing(loopNode))
		}
		l.pass.compiler.ReportChangeToEnclosingScope(loopNode)

		l.changeLoopLocalVariablesToProperties(loopNode, loopObject)
	}

	for _, fn := range wrapperFunctions {
		body := fn.Body().Node()
		for param := range fn.Params().Params() {
			paramName := ast.MustName(param.Node())
			loopObjectName, _, _ := strings.Cut(paramName.Ident(), es6BlockLoopParamNamePrefix)
			l.updateNames(body, paramName.Ident(), loopObjectName)
		}
	}
}

func (l *es6BlockLoopClosureTransformer) createWrapperFunctions() []ast.Function {
	wrapperFunctions := make([]ast.Function, 0, len(l.wrapTargets))
	for _, wt := range l.wrapTargets {
		functionOrObjectLit := wt.node
		if _, isObjLit := ast.AsObjectLit(functionOrObjectLit); isObjLit {
			enclosingFunction := NodeUtil_GetEnclosingFunction(functionOrObjectLit)
			var scopeBody rhino.Node
			isConstructor := false
			if fn, ok := ast.AsFunction(enclosingFunction); ok {
				scopeBody = fn.Body().Node()
				isConstructor = fn.IsEs6Constructor()
			} else {
				scopeBody = NodeUtil_GetEnclosingScript(functionOrObjectLit)
			}
			script := NodeUtil_GetEnclosingScript(functionOrObjectLit)
			input := l.pass.compiler.GetInput(script.GetInputId())
			thisContext := NewThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext(
				scopeBody, isConstructor, l.pass.uniqueIdSupplier.GetUniqueId(input))
			updater := NewThisAndArgumentsReferenceUpdater(l.pass.compiler, thisContext, l.pass.astFactory)
			NodeTraversal_Traverse(l.pass.compiler, functionOrObjectLit, updater)
			thisContext.AddVarDeclarations(l.pass.compiler, l.pass.astFactory, script)
			for child := scopeBody.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
				if child.IsLet() || child.IsConst() {
					l.pass.addLetConst(child)
				}
			}
		}
		returnNode := rhino.IR_ReturnNode()
		parameterNames := make([]rhino.Node, len(wt.objects))
		objectNamesForCall := make([]rhino.Node, len(wt.objects))
		for i, object := range wt.objects {
			parameterNames[i] = l.createUniqueParameterNameToUseWithinLoop(functionOrObjectLit, object)
			objectNamesForCall[i] = l.createLoopObjectNameNode(object)
		}

		iife := l.pass.astFactory.CreateFunction(
			"",
			rhino.IR_ParamList(parameterNames...),
			rhino.IR_BlockOfNode(returnNode),
			AstFactory_Type_OfColor(colors.StandardColors_TOP_OBJECT))
		l.pass.compiler.ReportChangeToChangeScope(iife)
		call := l.pass.astFactory.CreateCall(iife, AstFactory_Type_(functionOrObjectLit), objectNamesForCall...)
		ast.MustCall(call).SetFreeCall(true)
		jutil.CheckStateMsg1(
			!NodeUtil_IsFunctionDeclaration(functionOrObjectLit),
			"block-scoped function declarations should not exist now: %s",
			functionOrObjectLit)
		replacement := call.SrcrefTreeIfMissing(functionOrObjectLit)
		functionOrObjectLit.ReplaceWith(replacement)
		wrapperFunctions = append(wrapperFunctions, ast.MustFunction(iife))
		returnNode.AddChildToFront(functionOrObjectLit)
		l.pass.compiler.ReportChangeToEnclosingScope(replacement)
	}
	return wrapperFunctions
}

func (l *es6BlockLoopClosureTransformer) updateNames(node rhino.Node, paramIdent string, loopObjectName string) {
	if nm, ok := ast.AsName(node); ok && nm.Ident() == loopObjectName {
		nm.SetIdent(paramIdent)
		return
	}
	for child := node.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		l.updateNames(child, paramIdent, loopObjectName)
	}
}

func (l *es6BlockLoopClosureTransformer) changeVanillaForLoopHeader(forLoop ast.For, updateLoopObject rhino.Node) {
	loopNode := forLoop.Node()
	if init, hasInit := forLoop.Init(); hasInit {
		initializer := init.Node()
		initializer.ReplaceWith(rhino.IR_Empty())
		if !NodeUtil_IsNameDeclaration(initializer) {
			initializer = rhino.IR_ExprResult(initializer).Srcref(initializer)
		}
		l.pass.addNodeBeforeLoop(initializer, loopNode)
	}

	if inc, hasUpdate := forLoop.Update(); !hasUpdate {
		forLoop.SetUpdate(ast.MustExpr(updateLoopObject.SrcrefTreeIfMissing(loopNode)))
	} else {
		increment := inc.Node()
		placeHolder := rhino.IR_Empty()
		increment.ReplaceWith(placeHolder)
		placeHolder.ReplaceWith(
			l.pass.astFactory.CreateComma(updateLoopObject, increment).SrcrefTreeIfMissing(loopNode))
	}
}

func (l *es6BlockLoopClosureTransformer) changeLoopLocalVariablesToProperties(loopNode rhino.Node, loopObject *es6BlockLoopObject) {
	for _, v := range loopObject.vars {
		newPropertyName := v.GetNameNode().GetString()
		for _, reference := range l.references[v] {
			jutil.CheckStateMsg1(!loopNode.IsForOf(), "%s", loopNode)
			if forIn, ok := ast.AsForIn(loopNode); ok && forIn.Target().Node() == reference.GetParent() {
				l.assignLoopVarToLoopObjectProperty(forIn, loopObject, v, newPropertyName, reference)
			} else if _, isDecl := ast.AsNameDeclaration(reference.GetParent()); isDecl {
				l.replaceDeclarationWithProperty(loopObject, newPropertyName, reference)
			} else {
				referenceParent := reference.GetParent()
				reference.ReplaceWith(l.createLoopVarReferenceReplacement(loopObject, reference, newPropertyName))
				l.pass.compiler.ReportChangeToEnclosingScope(referenceParent)
			}
		}
	}
}

func (l *es6BlockLoopClosureTransformer) replaceDeclarationWithProperty(
	loopObject *es6BlockLoopObject, newPropertyName string, reference rhino.Node,
) {
	declaration := reference.GetParent()
	grandParent := declaration.GetParent()
	l.pass.removeLetConst(declaration)
	l.pass.handleLetConstDeclarationList(declaration, grandParent)

	declaration = reference.GetParent()
	if reference.HasChildren() {
		newReference := l.createLoopVarReferenceReplacement(loopObject, reference, newPropertyName)
		assign := l.pass.astFactory.CreateAssign(newReference, reference.RemoveFirstChild())
		es6BlockExtractInlineJSDoc(declaration, reference, declaration)
		es6BlockMaybeAddConstJSDoc(declaration, reference, declaration)
		assign.SetJSDocInfo(declaration.GetJSDocInfo())

		declaration.ReplaceWith(rhino.IR_ExprResult(assign).SrcrefTreeIfMissing(declaration))
	} else {
		declaration.Detach()
	}
	l.pass.compiler.ReportChangeToEnclosingScope(grandParent)
}

func (l *es6BlockLoopClosureTransformer) assignLoopVarToLoopObjectProperty(
	forIn ast.ForIn, loopObject *es6BlockLoopObject, v *Var, newPropertyName string, reference rhino.Node,
) {
	jutil.CheckStateMsg1(reference == v.GetNameNode(), "%s", reference)
	referenceParent := reference.GetParent()
	jutil.CheckStateMsg1(NodeUtil_IsNameDeclaration(referenceParent), "%s", referenceParent)
	jutil.CheckStateMsg1(reference.IsName(), "%s", reference)

	forInPropAssignmentStatement := rhino.IR_ExprResult(
		l.pass.astFactory.CreateAssign(
			l.createLoopVarReferenceReplacement(loopObject, reference, newPropertyName),
			reference.CloneNode())).
		SrcrefTreeIfMissing(reference)

	loopUpdateStatement := forIn.Body().Node().GetFirstChild()
	forInPropAssignmentStatement.InsertAfter(loopUpdateStatement)
}

func (l *es6BlockLoopClosureTransformer) renameVarsToProperties(
	loopObject *es6BlockLoopObject, objectLitNextIteration rhino.Node, scopeRoot rhino.Node,
) {
	for _, v := range loopObject.vars {
		newPropertyName := v.GetNameNode().GetString()
		var newPropertyValue rhino.Node
		if v.GetScopeRoot() == scopeRoot {
			newPropertyValue = l.createLoopVarReferenceReplacement(loopObject, v.GetNameNode(), newPropertyName)
		} else {
			newPropertyValue = l.pass.astFactory.CreateUndefinedValue()
		}
		objectLitNextIteration.AddChildToBack(
			l.pass.astFactory.CreateStringKey(newPropertyName, newPropertyValue))
	}
}

func (l *es6BlockLoopClosureTransformer) createLoopVarReferenceReplacement(
	loopObject *es6BlockLoopObject, reference rhino.Node, propertyName string,
) rhino.Node {
	return l.pass.astFactory.
		CreateGetProp(l.createLoopObjectNameNode(loopObject), propertyName, AstFactory_Type_(reference)).
		SrcrefTree(reference)
}

func (l *es6BlockLoopClosureTransformer) createUniqueParameterNameToUseWithinLoop(node rhino.Node, loopObject *es6BlockLoopObject) rhino.Node {
	return l.pass.astFactory.CreateName(
		loopObject.name+
			es6BlockLoopParamNamePrefix+
			l.pass.uniqueIdSupplier.GetUniqueId(l.pass.compiler.GetInput(NodeUtil_GetInputId(node))),
		AstFactory_Type_OfColor(colors.StandardColors_TOP_OBJECT))
}

func (l *es6BlockLoopClosureTransformer) createLoopObjectNameNode(loopObject *es6BlockLoopObject) rhino.Node {
	return l.pass.astFactory.CreateName(loopObject.name, AstFactory_Type_OfColor(colors.StandardColors_TOP_OBJECT))
}
