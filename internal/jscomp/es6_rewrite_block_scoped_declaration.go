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
	letConsts         *jutil.LinkedHashSet[rhino.Node]
	loopLetConstNames map[string]struct{}
	uniqueIdSupplier  *UniqueIdSupplier
}

func (e *Es6RewriteBlockScopedDeclaration) noVisitLeaves() {}

var es6RewriteBlockScopedDeclaration_transpiledFeatures = parser.FeatureSet_BARE_MINIMUM.WithFeatures(
	parser.FeatureSet_Feature_LET_DECLARATIONS,
	parser.FeatureSet_Feature_CONST_DECLARATIONS,
)

const es6BlockLoopParamNamePrefix = "$jscomp$loop_param$"

// Java: Es6RewriteBlockScopedDeclaration(AbstractCompiler)
func NewEs6RewriteBlockScopedDeclaration(compiler AbstractCompiler) *Es6RewriteBlockScopedDeclaration {
	return &Es6RewriteBlockScopedDeclaration{
		letConsts:        jutil.NewLinkedHashSet[rhino.Node](),
		compiler:         compiler,
		uniqueIdSupplier: compiler.GetUniqueIdSupplier(),
		astFactory:       compiler.CreateAstFactory(),
	}
}

// Java: Es6RewriteBlockScopedDeclaration.visit(NodeTraversal, Node, Node)
func (e *Es6RewriteBlockScopedDeclaration) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	jutil.CheckStateMsg1(parent.IsNil() || !parent.IsForOf(), "%s", parent)

	if n.IsLet() || n.IsConst() {
		e.letConsts.Add(n)
		for nameNode := n.GetFirstChild(); !nameNode.IsNil(); nameNode = nameNode.GetNext() {
			e.visitBlockScopedNameDeclaration(n, nameNode)
		}
	}
}

// Java: Es6RewriteBlockScopedDeclaration.process(Node, Node)
func (e *Es6RewriteBlockScopedDeclaration) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(e.compiler, root, e)
	if !e.letConsts.IsEmpty() {
		if len(e.loopLetConstNames) > 0 {
			sc := NewSyntacticScopeCreator(e.compiler)
			transformer := newEs6BlockLoopClosureTransformer(e)
			NodeTraversal_Builder_().SetCompiler(e.compiler).SetScopeCreator(sc).SetCallback(transformer).Traverse(root)
			transformer.transformLoopClosure()
		}
		e.rewriteDeclsToVars()
	}
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(e.compiler, root, es6RewriteBlockScopedDeclaration_transpiledFeatures)
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

	if parent.IsVanillaFor() {
		e.handleLetConstInVanillaForInitializer(declarationList, parent)
		return
	}

	if parent.IsForIn() {
		e.handleLetConstInForInInitializer(declarationList, parent)
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

func (e *Es6RewriteBlockScopedDeclaration) handleLetConstInForInInitializer(declarationList rhino.Node, parent rhino.Node) {
	jutil.CheckState(parent.IsForIn())
	first := declarationList
	lhs := first.GetFirstChild()
	loopNode := parent
	if lhs.IsDestructuringLhs() {
		jutil.CheckStateMsg(
			false, "Destructuring syntax is unsupported in ES6RewriteBlockScopedDeclarations pass")
		NodeUtil_VisitLhsNodesInNode(
			lhs,
			func(name rhino.Node) {
				nm, ok := ast.AsName(name)
				jutil.CheckStateMsg1(ok, "lhs in destructuring declaration should be a simple name. (%s)", name)
				newName := rhino.IR_Name(nm.Ident()).Srcref(name)
				newVar := rhino.IR_Var(newName).Srcref(name)
				es6BlockExtractInlineJSDoc(declarationList, name, newVar)
				e.addNodeBeforeLoop(newVar, loopNode)
			})

		destructuringPattern := lhs.RemoveFirstChild()
		first.ReplaceWith(destructuringPattern)
	} else {
		newStatement := first.CloneTree()
		name := newStatement.GetFirstChild().CloneNode()
		newStatement.SetToken(rhino.Token_VAR)
		es6BlockExtractInlineJSDoc(first, first.GetFirstChild(), newStatement)
		first.ReplaceWith(name)
		e.addNodeBeforeLoop(newStatement, loopNode)
	}
}

func (e *Es6RewriteBlockScopedDeclaration) handleLetConstInVanillaForInitializer(
	declarationList rhino.Node, parent rhino.Node) {
	jutil.CheckState(parent.IsVanillaFor())
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
	for n := range e.letConsts.All() {
		e.handleLetConstDeclarationList(n, n.GetParent())
	}
}

type es6LoopClosureHandledKey struct {
	node rhino.Node
	name string
}

// es6BlockLoopClosureTransformer transforms let/const declarations captured by loop closures.
type es6BlockLoopClosureTransformer struct {
	NodeTraversal_AbstractPreOrderCallbackBase
	pass *Es6RewriteBlockScopedDeclaration

	loopObjectMap                       *jutil.LinkedHashMap[rhino.Node, *es6BlockLoopObject]
	nodesRequiringLoopObjectsClosureMap *jutil.SetMultimap[rhino.Node, *es6BlockLoopObject]
	nodesHandledForLoopObjectClosure    map[es6LoopClosureHandledKey]struct{}
	referenceMap                        *jutil.SetMultimap[*Var, rhino.Node]
}

const es6BlockLoopObjectName = "$jscomp$loop"

func newEs6BlockLoopClosureTransformer(pass *Es6RewriteBlockScopedDeclaration) *es6BlockLoopClosureTransformer {
	return &es6BlockLoopClosureTransformer{
		pass:                                pass,
		loopObjectMap:                       jutil.NewLinkedHashMap[rhino.Node, *es6BlockLoopObject](),
		nodesRequiringLoopObjectsClosureMap: jutil.NewSetMultimap[rhino.Node, *es6BlockLoopObject](),
		referenceMap:                        jutil.NewSetMultimap[*Var, rhino.Node](),
	}
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

	l.referenceMap.Put(v, n)

	var outerMostFunctionScope *Scope
	for s := referencedIn; s != declaredIn && s.GetRootNode() != loopNode; s = s.GetParent() {
		if s.IsFunctionScope() {
			outerMostFunctionScope = s
		}
	}

	if outerMostFunctionScope != nil {
		enclosingFunction := outerMostFunctionScope.GetRootNode()

		var nodeToWrapInClosure rhino.Node
		if _, isAccessor := ast.AsGetterOrSetterDef(enclosingFunction.GetParent()); isAccessor {
			nodeToWrapInClosure = enclosingFunction.GetGrandparent()
			jutil.CheckState(nodeToWrapInClosure.IsObjectLit())
		} else {
			nodeToWrapInClosure = enclosingFunction
		}
		key := es6LoopClosureHandledKey{node: nodeToWrapInClosure, name: name}
		if _, ok := l.nodesHandledForLoopObjectClosure[key]; ok {
			return true
		}
		if l.nodesHandledForLoopObjectClosure == nil {
			l.nodesHandledForLoopObjectClosure = make(map[es6LoopClosureHandledKey]struct{}, 8)
		}
		l.nodesHandledForLoopObjectClosure[key] = struct{}{}

		object := l.loopObjectMap.ComputeIfAbsent(
			loopNode, func(k rhino.Node) *es6BlockLoopObject {
				return &es6BlockLoopObject{
					name: es6BlockLoopObjectName + "$" + l.pass.uniqueIdSupplier.GetUniqueId(t.GetInput()),
					vars: jutil.NewLinkedHashSet[*Var](),
				}
			})
		object.vars.Add(v)
		l.nodesRequiringLoopObjectsClosureMap.Put(nodeToWrapInClosure, object)
	}
	return true
}

func (l *es6BlockLoopClosureTransformer) loopObjPropName(v *Var) string {
	return v.GetNameNode().GetString()
}

func (l *es6BlockLoopClosureTransformer) transformLoopClosure() {
	if l.loopObjectMap.IsEmpty() {
		return
	}

	wrapperFunctions := l.createWrapperFunctions()

	for loopNode := range l.loopObjectMap.Keys() {
		loopObject := l.loopObjectMap.Get(loopNode)
		objectLitNextIteration := l.pass.astFactory.CreateObjectLit()
		l.renameVarsToProperties(loopObject, objectLitNextIteration, loopNode)

		updateLoopObject :=
			l.pass.astFactory.CreateAssign(l.createLoopObjectNameNode(loopObject), objectLitNextIteration)
		objectLit :=
			rhino.IR_VarWithValue(l.createLoopObjectNameNode(loopObject), l.pass.astFactory.CreateObjectLit()).
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

	l.updateNamesInWrapperFunctions(wrapperFunctions)
}

func (l *es6BlockLoopClosureTransformer) updateNamesInWrapperFunctions(wrapperFunctions *jutil.LinkedHashSet[rhino.Node]) {
	for fnNode := range wrapperFunctions.All() {
		fn := ast.MustFunction(fnNode)
		body := fn.Body().Node()
		for param := range fn.Params().Params() {
			paramName := ast.MustName(param.Node())
			loopObjectName, _, _ := strings.Cut(paramName.Ident(), es6BlockLoopParamNamePrefix)
			l.updateNames(body, paramName.Ident(), loopObjectName)
		}
	}
}

func (l *es6BlockLoopClosureTransformer) createWrapperFunctions() *jutil.LinkedHashSet[rhino.Node] {
	wrapperFunctions := jutil.NewLinkedHashSet[rhino.Node]()
	for functionOrObjectLit := range l.nodesRequiringLoopObjectsClosureMap.KeySet() {
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
			thisContext :=
				NewThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext(
					scopeBody, isConstructor, l.pass.uniqueIdSupplier.GetUniqueId(input))
			updater :=
				NewThisAndArgumentsReferenceUpdater(l.pass.compiler, thisContext, l.pass.astFactory)
			NodeTraversal_Traverse(l.pass.compiler, functionOrObjectLit, updater)
			thisContext.AddVarDeclarations(l.pass.compiler, l.pass.astFactory, script)
			for child := scopeBody.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
				if child.IsLet() || child.IsConst() {
					l.pass.letConsts.Add(child)
				}
			}
		}
		returnNode := rhino.IR_ReturnNode()
		objects := l.nodesRequiringLoopObjectsClosureMap.Get(functionOrObjectLit)
		parameterNames := make([]rhino.Node, objects.Size())
		objectNamesForCall := make([]rhino.Node, objects.Size())
		loopObjectIndex := 0
		for object := range objects.All() {
			parameterNames[loopObjectIndex] =
				l.createUniqueParameterNameToUseWithinLoop(functionOrObjectLit, object)
			objectNamesForCall[loopObjectIndex] = l.createLoopObjectNameNode(object)
			loopObjectIndex++
		}

		iife :=
			l.pass.astFactory.CreateFunction(
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
		wrapperFunctions.Add(iife)
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

	increment := loopNode.GetChildAtIndex(2)
	if increment.IsEmpty() {
		increment.ReplaceWith(updateLoopObject.SrcrefTreeIfMissing(loopNode))
	} else {
		placeHolder := rhino.IR_Empty()
		increment.ReplaceWith(placeHolder)
		placeHolder.ReplaceWith(
			l.pass.astFactory.CreateComma(updateLoopObject, increment).SrcrefTreeIfMissing(loopNode))
	}
}

func (l *es6BlockLoopClosureTransformer) changeLoopLocalVariablesToProperties(loopNode rhino.Node, loopObject *es6BlockLoopObject) {
	for v := range loopObject.vars.All() {
		newPropertyName := l.loopObjPropName(v)
		for reference := range l.referenceMap.Get(v).All() {
			jutil.CheckStateMsg1(!loopNode.IsForOf(), "%s", loopNode)
			if forIn, ok := ast.AsForIn(loopNode); ok && forIn.Target().Node() == reference.GetParent() {
				l.assignLoopVarToLoopObjectProperty(
					forIn, loopObject, v, newPropertyName, reference)
			} else if _, isDecl := ast.AsNameDeclaration(reference.GetParent()); isDecl {
				l.replaceDeclarationWithProperty(loopObject, newPropertyName, reference)
			} else {
				l.replaceReferenceWithProperty(loopObject, newPropertyName, reference)
			}
		}
	}
}

func (l *es6BlockLoopClosureTransformer) replaceDeclarationWithProperty(
	loopObject *es6BlockLoopObject, newPropertyName string, reference rhino.Node) {
	declaration := reference.GetParent()
	grandParent := declaration.GetParent()
	l.pass.letConsts.Remove(declaration)
	l.pass.handleLetConstDeclarationList(declaration, grandParent)

	declaration = reference.GetParent()
	if reference.HasChildren() {
		newReference := l.createLoopVarReferenceReplacement(loopObject, reference, newPropertyName)
		assign := l.pass.astFactory.CreateAssign(newReference, reference.RemoveFirstChild())
		es6BlockExtractInlineJSDoc(declaration, reference, declaration)
		es6BlockMaybeAddConstJSDoc(declaration, reference, declaration)
		assign.SetJSDocInfo(declaration.GetJSDocInfo())

		replacement := rhino.IR_ExprResult(assign).SrcrefTreeIfMissing(declaration)
		declaration.ReplaceWith(replacement)
	} else {
		declaration.Detach()
	}
	l.pass.compiler.ReportChangeToEnclosingScope(grandParent)
}

func (l *es6BlockLoopClosureTransformer) replaceReferenceWithProperty(
	loopObject *es6BlockLoopObject, newPropertyName string, reference rhino.Node) {
	referenceParent := reference.GetParent()
	reference.ReplaceWith(l.createLoopVarReferenceReplacement(loopObject, reference, newPropertyName))
	l.pass.compiler.ReportChangeToEnclosingScope(referenceParent)
}

func (l *es6BlockLoopClosureTransformer) assignLoopVarToLoopObjectProperty(
	forIn ast.ForIn, loopObject *es6BlockLoopObject, v *Var, newPropertyName string, reference rhino.Node) {
	jutil.CheckStateMsg1(reference == v.GetNameNode(), "%s", reference)
	referenceParent := reference.GetParent()
	jutil.CheckStateMsg1(NodeUtil_IsNameDeclaration(referenceParent), "%s", referenceParent)
	jutil.CheckStateMsg1(reference.IsName(), "%s", reference)

	loopVarReference := reference.CloneNode()
	forInPropAssignmentStatement :=
		rhino.IR_ExprResult(
			l.pass.astFactory.CreateAssign(
				l.createLoopVarReferenceReplacement(loopObject, reference, newPropertyName),
				loopVarReference)).
			SrcrefTreeIfMissing(reference)

	loopUpdateStatement := forIn.Body().Node().GetFirstChild()
	forInPropAssignmentStatement.InsertAfter(loopUpdateStatement)
}

func (l *es6BlockLoopClosureTransformer) renameVarsToProperties(
	loopObject *es6BlockLoopObject, objectLitNextIteration rhino.Node, scopeRoot rhino.Node) {
	for v := range loopObject.vars.All() {
		newPropertyName := l.loopObjPropName(v)
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
	loopObject *es6BlockLoopObject, reference rhino.Node, propertyName string) rhino.Node {
	replacement :=
		l.pass.astFactory.CreateGetProp(
			l.createLoopObjectNameNode(loopObject), propertyName, AstFactory_Type_(reference))
	replacement.SrcrefTree(reference)
	return replacement
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

type es6BlockLoopObject struct {
	name string
	vars *jutil.LinkedHashSet[*Var]
}
