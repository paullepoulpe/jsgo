// Port of com/google/javascript/jscomp/Es6ConvertSuperConstructorCalls.java @ closure-compiler v20261021

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

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Es6ConvertSuperConstructorCalls converts `super()` calls.
type Es6ConvertSuperConstructorCalls struct {
	jscompInherits  RuntimeJsLibManager_JsLibField
	jscompConstruct RuntimeJsLibManager_JsLibField

	compiler                 AbstractCompiler
	constructorDataStack     []es6ConstructorData
	constructorDataBuf       [8]es6ConstructorData
	astFactory               *AstFactory
	globalNamespace          *GlobalNamespace
	transpilationNamespace   rhino.StaticScope
	uniqueIdSupplier         *UniqueIdSupplier
	es6SubclassTranspilation CompilerOptions_Es6SubclassTranspilation
}

var (
	es6ConvertSuperConstructorCallsVisitTokens = func() ast.TokenSet {
		var s ast.TokenSet
		s.Add(rhino.Token_FUNCTION)
		return s
	}()

	es6UndefinedReturnCheckVisitTokens = func() ast.TokenSet {
		var s ast.TokenSet
		s.Add(rhino.Token_RETURN)
		return s
	}()

	es6ReplaceThisWithSuperThisVisitTokens = func() ast.TokenSet {
		var s ast.TokenSet
		s.Add(rhino.Token_THIS)
		s.Add(rhino.Token_RETURN)
		return s
	}()
)

const (
	es6SuperConstructorTmpError  = "$jscomp$tmp$error"
	es6SuperConstructorSuperThis = "$jscomp$super$this"
)

type es6ConstructorData struct {
	constructor rhino.Node
	superCalls  []rhino.Node
}

// Java: Es6ConvertSuperConstructorCalls(AbstractCompiler, Es6SubclassTranspilation)
func NewEs6ConvertSuperConstructorCalls(
	compiler AbstractCompiler, es6SubclassTranspilation CompilerOptions_Es6SubclassTranspilation) *Es6ConvertSuperConstructorCalls {
	runtimeJsLibManager := compiler.GetRuntimeJsLibManager()
	e := &Es6ConvertSuperConstructorCalls{
		compiler:                 compiler,
		astFactory:               compiler.CreateAstFactory(),
		transpilationNamespace:   compiler.GetTranspilationNamespace(),
		uniqueIdSupplier:         compiler.GetUniqueIdSupplier(),
		jscompInherits:           runtimeJsLibManager.GetJsLibField("$jscomp.inherits"),
		jscompConstruct:          runtimeJsLibManager.GetJsLibField("$jscomp.construct"),
		es6SubclassTranspilation: es6SubclassTranspilation,
	}
	e.constructorDataStack = e.constructorDataBuf[:0]
	return e
}

func (*Es6ConvertSuperConstructorCalls) noVisitLeaves() {}
func (*Es6ConvertSuperConstructorCalls) visitTokens() *ast.TokenSet {
	return &es6ConvertSuperConstructorCallsVisitTokens
}

func (e *Es6ConvertSuperConstructorCalls) peekConstructorData() *es6ConstructorData {
	if len(e.constructorDataStack) == 0 {
		return nil
	}
	return &e.constructorDataStack[len(e.constructorDataStack)-1]
}

// Java: Es6ConvertSuperConstructorCalls.shouldTraverse(NodeTraversal, Node, Node)
func (e *Es6ConvertSuperConstructorCalls) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if _, ok := ast.AsScript(n); ok && !NodeUtil_GetFeatureSetOfScript(n).ContainsFeature(parser.FeatureSet_Feature_SUPER) {
		return false
	}
	if _, ok := ast.AsFunction(n); ok {
		if e.constructorDataStack == nil {
			e.constructorDataStack = e.constructorDataBuf[:0]
		}
		e.constructorDataStack = append(e.constructorDataStack, es6ConstructorData{constructor: n})
	} else if n.IsSuper() {
		jutil.CheckStateMsg1(n.IsFirstChildOf(parent), "%s", parent)
		if _, ok := ast.AsGetProp(parent); ok {
			t.Report(n, TranspilationUtil_CANNOT_CONVERT_YET, "super access with no extends clause")
			return false
		}
		jutil.CheckStateMsg1(parent.IsCall(), "%s", parent)
		constructorData := jutil.CheckNotNull(e.peekConstructorData())
		constructorData.superCalls = append(constructorData.superCalls, parent)
	}
	return true
}

// Java: Es6ConvertSuperConstructorCalls.visit(NodeTraversal, Node, Node)
func (e *Es6ConvertSuperConstructorCalls) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if top := e.peekConstructorData(); top != nil && n == top.constructor {
		data := *top
		*top = es6ConstructorData{}
		e.constructorDataStack = e.constructorDataStack[:len(e.constructorDataStack)-1]
		e.visitSuper(t, data)
	}
}

func (e *Es6ConvertSuperConstructorCalls) visitSuper(t *NodeTraversal, constructorData es6ConstructorData) {
	constructor := constructorData.constructor
	superCalls := constructorData.superCalls
	if len(superCalls) == 0 {
		return
	}

	if constructor.IsFromExterns() {
		for _, superCall := range superCalls {
			enclosingStatement := NodeUtil_GetEnclosingStatement(superCall)
			enclosingScope := enclosingStatement.GetParent()
			enclosingStatement.Detach()
			e.compiler.ReportChangeToEnclosingScope(enclosingScope)
		}
		return
	}
	uniqueSuperThisName :=
		es6SuperConstructorSuperThis +
			"$" +
			e.uniqueIdSupplier.GetUniqueId(e.compiler.GetInput(NodeUtil_GetInputId(constructor)))
	superClassNameNode := e.superClassQNameNode(constructor)
	superClassQName := *superClassNameNode.GetQualifiedName()
	thisType := e.typeOfThisForConstructor(constructor)

	if e.isNativeObjectClass(t, superClassQName) {
		e.replaceSuperCallsWithThis(superCalls, thisType)
	} else if e.isKnownNativeClass(t, superClassQName) {
		e.rewriteSuperCallsToJsCompConstructCalls(
			constructor, superCalls, superClassNameNode, thisType, uniqueSuperThisName)
	} else if e.isNativeErrorClass(t, superClassQName) {
		e.rewriteNativeErrorSuperCalls(t, superCalls, superClassNameNode, thisType)
	} else if e.isKnownToReturnOnlyUndefined(superClassQName) {
		e.rewriteSuperCalls(superCalls, superClassNameNode, thisType)
	} else {
		e.rewriteSuperAndPreserveReturnedThis(
			constructor, superCalls, superClassNameNode, thisType, uniqueSuperThisName)
	}
}

func (e *Es6ConvertSuperConstructorCalls) replaceSuperCallsWithThis(superCalls []rhino.Node, thisType AstFactory_Type) {
	for _, superCall := range superCalls {
		var sideEffectArgs []rhino.Node
		callee := superCall.RemoveFirstChild()
		jutil.CheckStateMsg1(callee.IsSuper(), "%s", callee)
		for superCall.HasChildren() {
			arg := superCall.RemoveFirstChild()
			if _, ok := ast.AsSpread(arg); ok {
				arg = arg.RemoveFirstChild()
			}
			if e.compiler.GetAstAnalyzer().MayHaveSideEffects(arg) {
				sideEffectArgs = append(sideEffectArgs, arg)
			}
		}

		replacement := e.astFactory.CreateThis(thisType).Srcref(superCall)
		if len(sideEffectArgs) != 0 {
			exprs := sideEffectArgs[0]
			for i := 1; i < len(sideEffectArgs); i++ {
				exprs = e.astFactory.CreateComma(exprs, sideEffectArgs[i]).SrcrefTreeIfMissing(superCall)
			}
			replacement = e.astFactory.CreateComma(exprs, replacement).SrcrefTreeIfMissing(superCall)
		}
		superCall.ReplaceWith(replacement)
		e.compiler.ReportChangeToEnclosingScope(replacement)
	}
}

func (e *Es6ConvertSuperConstructorCalls) rewriteSuperCallsToJsCompConstructCalls(
	constructor rhino.Node,
	superCalls []rhino.Node,
	superClassNameNode rhino.Node,
	thisType AstFactory_Type,
	uniqueSuperThisName string,
) {
	constructorBody := rhino.CheckNotNull(ast.MustFunction(constructor).Body().Node())
	firstStatement := constructorBody.GetFirstChild()
	rhino.CheckNotNullMsg(firstStatement, "Empty constructor body")
	firstSuperCall := superCalls[0]

	if exprStmt, ok := ast.AsExprStmt(firstStatement); ok &&
		constructorBody.HasOneChild() &&
		exprStmt.Expr().Node() == firstSuperCall {
		jutil.CheckStateMsg1(len(superCalls) == 1, "%s", constructor)
		firstStatement.ReplaceWith(
			e.astFactory.CreateReturn(
				e.createJSCompConstructorCall(superClassNameNode, firstSuperCall, thisType)))
	} else {
		typeOfThis := e.typeOfThisForConstructor(constructor)
		e.updateThisToSuperThis(typeOfThis, constructorBody, superCalls, uniqueSuperThisName)
		declaration := e.astFactory.
			CreateSingleVarNameDeclaration(uniqueSuperThisName).
			SrcrefTree(constructorBody)
		insertBeforePoint := NodeUtil_GetInsertionPointAfterAllInnerFunctionDeclarations(constructorBody)
		if !insertBeforePoint.IsNil() {
			declaration.InsertBefore(insertBeforePoint)
		} else {
			constructorBody.AddChildToBack(declaration)
		}

		constructorBody.AddChildToBack(
			e.astFactory.
				CreateReturn(e.astFactory.CreateName(uniqueSuperThisName, typeOfThis)).
				SrcrefTree(constructorBody))
		for _, superCall := range superCalls {
			superCall.ReplaceWith(
				e.astFactory.
					CreateAssign(
						e.astFactory.CreateName(uniqueSuperThisName, typeOfThis).Srcref(superCall),
						e.createJSCompConstructorCall(superClassNameNode, superCall, thisType)).
					Srcref(superCall))
		}
	}
	e.compiler.ReportChangeToEnclosingScope(constructorBody)
}

func (e *Es6ConvertSuperConstructorCalls) rewriteNativeErrorSuperCalls(
	t *NodeTraversal, superCalls []rhino.Node, superClassNameNode rhino.Node, thisType AstFactory_Type) {
	for _, superCall := range superCalls {
		newSuperCall :=
			e.createNewSuperCall(superClassNameNode, superCall, thisType, AstFactory_Type_(superCall))
		e.replaceNativeErrorSuperCall(superCall, newSuperCall, t.GetInput())
	}
}

func (e *Es6ConvertSuperConstructorCalls) rewriteSuperCalls(
	superCalls []rhino.Node, superClassNameNode rhino.Node, thisType AstFactory_Type) {
	for _, superCall := range superCalls {
		newSuperCall :=
			e.createNewSuperCall(
				superClassNameNode, superCall, thisType, AstFactory_Type_OfColor(colors.StandardColors_NULL_OR_VOID))
		superCallParent := superCall.GetParent()
		if superCallParent.HasOneChild() && NodeUtil_IsStatement(superCallParent) {
			superCall.ReplaceWith(newSuperCall)
		} else {
			superCall.ReplaceWith(
				e.astFactory.
					CreateComma(newSuperCall, e.astFactory.CreateThis(thisType)).
					SrcrefTreeIfMissing(superCall))
		}
		e.compiler.ReportChangeToEnclosingScope(superCallParent)
	}
}

func (e *Es6ConvertSuperConstructorCalls) rewriteSuperAndPreserveReturnedThis(
	constructor rhino.Node,
	superCalls []rhino.Node,
	superClassNameNode rhino.Node,
	thisType AstFactory_Type,
	uniqueSuperThisName string,
) {
	if e.es6SubclassTranspilation == CompilerOptions_Es6SubclassTranspilation_SAFE_REFLECT_CONSTRUCT {
		e.rewriteSuperCallsToJsCompConstructCalls(
			constructor, superCalls, superClassNameNode, thisType, uniqueSuperThisName)
		return
	}
	constructorBody := rhino.CheckNotNull(ast.MustFunction(constructor).Body().Node())
	firstStatement := constructorBody.GetFirstChild()
	firstSuperCall := superCalls[0]

	if exprStmt, ok := ast.AsExprStmt(firstStatement); ok &&
		constructorBody.HasOneChild() &&
		exprStmt.Expr().Node() == firstSuperCall {
		jutil.CheckStateMsg1(len(superCalls) == 1, "%s", constructor)
		newReturn :=
			e.astFactory.CreateOr(
				e.createNewSuperCall(
					superClassNameNode,
					firstSuperCall,
					AstFactory_Type_(firstSuperCall),
					AstFactory_Type_OfColor(colors.StandardColors_UNKNOWN)),
				e.astFactory.CreateThis(thisType))
		firstStatement.ReplaceWith(rhino.IR_ReturnNodeWithExpr(newReturn).SrcrefTreeIfMissing(firstStatement))
	} else {
		typeOfThis := e.typeOfThisForConstructor(constructor)
		declaration :=
			rhino.IR_Var(e.astFactory.CreateName(uniqueSuperThisName, typeOfThis)).
				SrcrefTree(constructorBody)
		insertBeforePoint := NodeUtil_GetInsertionPointAfterAllInnerFunctionDeclarations(constructorBody)
		if !insertBeforePoint.IsNil() {
			declaration.InsertBefore(insertBeforePoint)
		} else {
			constructorBody.AddChildToBack(declaration)
		}
		constructorBody.AddChildToBack(
			rhino.IR_ReturnNodeWithExpr(e.astFactory.CreateName(uniqueSuperThisName, typeOfThis)).
				SrcrefTree(constructorBody))

		e.updateThisToSuperThis(typeOfThis, constructorBody, superCalls, uniqueSuperThisName)

		for _, superCall := range superCalls {
			newSuperCall :=
				e.createNewSuperCall(
					superClassNameNode, superCall, typeOfThis, AstFactory_Type_OfColor(colors.StandardColors_UNKNOWN))
			superCall.ReplaceWith(
				e.astFactory.
					CreateAssign(
						e.astFactory.CreateName(uniqueSuperThisName, typeOfThis),
						e.astFactory.CreateOr(newSuperCall, e.astFactory.CreateThis(typeOfThis))).
					SrcrefTreeIfMissing(superCall))
		}
	}
	e.compiler.ReportChangeToEnclosingScope(constructorBody)
}

func (e *Es6ConvertSuperConstructorCalls) isKnownToReturnOnlyUndefined(functionQName string) bool {
	if e.globalNamespace == nil {
		return false
	}
	slot := e.globalNamespace.GetSlot(functionQName)
	if slot == nil {
		return false
	}
	globalName := slot.(*GlobalNamespace_Name)

	var declarationRef *GlobalNamespace_Ref
	if d := globalName.GetDeclaration(); d != nil {
		declarationRef = d.(*GlobalNamespace_Ref)
	}
	if declarationRef == nil {
		for ref := range globalName.Refs() {
			if ref.IsSet() {
				declarationRef = ref
			}
		}
	}
	if declarationRef == nil {
		return false
	}

	declaredVarOrProp := declarationRef.GetNode()
	if declaredVarOrProp.IsFromExterns() {
		return false
	}

	declaration := declaredVarOrProp.GetParent()
	var declaredValue rhino.Node
	if _, ok := ast.AsFunction(declaration); ok {
		declaredValue = declaration
	} else if _, ok := ast.AsNameDeclaration(declaration); ok && declaredVarOrProp.IsName() {
		init, hasInit := ast.MustName(declaredVarOrProp).Init()
		if !hasInit {
			return false
		}
		declaredValue = rhino.CheckNotNull(init.Node())
	} else if assign, ok := ast.AsAssign(declaration); ok && assign.Target().Node() == declaredVarOrProp {
		declaredValue = rhino.CheckNotNull(assign.Value().Node())
	} else if declaration.IsObjectLit() && declaredVarOrProp.HasOneChild() {
		declaredValue = rhino.CheckNotNull(declaredVarOrProp.GetFirstChild())
	} else {
		panic(jutil.NewIllegalStateExceptionMsg(
			"Unexpected declaration format:\n" + declaration.ToStringTree()))
	}
	return e.isNodeKnownToOnlyReturnUndefined(declaredValue)
}

func (e *Es6ConvertSuperConstructorCalls) isNodeKnownToOnlyReturnUndefined(node rhino.Node) bool {
	if fn, ok := ast.AsFunction(node); ok {
		functionBody := rhino.CheckNotNull(fn.Body().Node())
		var check es6UndefinedReturnCheck
		NodeTraversal_Traverse(e.compiler, functionBody, &check)
		return !check.foundNonEmptyReturn
	}
	if node.IsQualifiedName() {
		return e.isKnownToReturnOnlyUndefined(*node.GetQualifiedName())
	}
	if hk, ok := ast.AsHook(node); ok {
		return e.isNodeKnownToOnlyReturnUndefined(hk.Then().Node()) &&
			e.isNodeKnownToOnlyReturnUndefined(hk.Else().Node())
	}
	return false
}

type es6UndefinedReturnCheck struct {
	NodeTraversal_AbstractShallowCallbackBase
	foundNonEmptyReturn bool
}

func (*es6UndefinedReturnCheck) visitTokens() *ast.TokenSet {
	return &es6UndefinedReturnCheckVisitTokens
}

func (c *es6UndefinedReturnCheck) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if !c.foundNonEmptyReturn {
		if ret, ok := ast.AsReturn(n); ok {
			if val, hasVal := ret.Value(); hasVal && !val.Node().MatchesName("undefined") {
				c.foundNonEmptyReturn = true
			}
		}
	}
}

func (e *Es6ConvertSuperConstructorCalls) createNewSuperCall(
	superClassQNameNode rhino.Node,
	superCall rhino.Node,
	thisType AstFactory_Type,
	resultType AstFactory_Type,
) rhino.Node {
	jutil.CheckArgumentMsg1(superClassQNameNode.IsQualifiedName(), "%s", superClassQNameNode)
	jutil.CheckArgumentMsg1(superCall.IsCall(), "%s", superCall)

	callee := superCall.RemoveFirstChild()
	jutil.CheckStateMsg1(callee.IsSuper(), "%s", callee)

	var args []rhino.Node
	hasSpreadArg := false
	for superCall.HasChildren() {
		arg := superCall.RemoveFirstChild()
		_, isSpread := ast.AsSpread(arg)
		hasSpreadArg = hasSpreadArg || isSpread
		args = append(args, arg)
	}

	if hasSpreadArg {
		superClassDotApply := e.astFactory.
			CreateGetPropWithUnknownType(superClassQNameNode.CloneTree(), "apply").
			SrcrefTree(callee)
		newSuperCall := e.astFactory.CreateCall(superClassDotApply, resultType).Srcref(superCall)
		newSuperCall.AddChildToBack(e.astFactory.CreateThis(thisType).Srcref(callee))
		ast.MustCall(newSuperCall).SetFreeCall(false)
		if isSingleSpreadOfArguments(args) {
			newSuperCall.AddChildToBack(args[0].GetOnlyChild().Detach())
		} else {
			newSuperCall.AddChildToBack(e.astFactory.CreateArraylitOfIterable(slices.Values(args)).Srcref(superCall))
		}
		return newSuperCall
	}

	superClassDotCall := e.astFactory.
		CreateGetProp(
			superClassQNameNode.CloneTree(), "call", AstFactory_Type_OfColor(colors.StandardColors_TOP_OBJECT)).
		SrcrefTree(callee)
	newSuperCall := e.astFactory.CreateCall(superClassDotCall, resultType).Srcref(superCall)
	newSuperCall.AddChildToBack(e.astFactory.CreateThis(thisType).Srcref(callee))
	ast.MustCall(newSuperCall).SetFreeCall(false)
	for _, arg := range args {
		newSuperCall.AddChildToBack(arg)
	}
	return newSuperCall
}

func (e *Es6ConvertSuperConstructorCalls) createJSCompConstructorCall(
	superClassQNameNode rhino.Node, superCall rhino.Node, thisType AstFactory_Type) rhino.Node {
	jutil.CheckArgumentMsg1(superClassQNameNode.IsQualifiedName(), "%s", superClassQNameNode)
	jutil.CheckArgumentMsg1(superCall.IsCall(), "%s", superCall)

	callee := rhino.CheckNotNullMsg1(superCall.RemoveFirstChild(), "%s", superCall)
	jutil.CheckStateMsg1(callee.IsSuper(), "%s", callee)

	jscompDotConstruct :=
		e.astFactory.CreateQNameOfField(e.transpilationNamespace, e.jscompConstruct).SrcrefTree(callee)
	superClassQName := superClassQNameNode.CloneTree()

	var superCallArgList []rhino.Node
	for superCall.HasChildren() {
		superCallArgList = append(superCallArgList, superCall.RemoveFirstChild())
	}
	var superArgs rhino.Node
	if isSingleSpreadOfArguments(superCallArgList) {
		superArgs = superCallArgList[0].GetOnlyChild().Detach()
	} else {
		superArgs = e.astFactory.CreateArraylitOfIterable(slices.Values(superCallArgList)).Srcref(superCall)
	}

	thisDotConstructor := e.astFactory.
		CreateGetProp(
			e.astFactory.CreateThis(thisType), "constructor", AstFactory_Type_(superClassQNameNode)).
		SrcrefTree(superCall)

	return e.astFactory.
		CreateCall(
			jscompDotConstruct, AstFactory_Type_(superCall), superClassQName, superArgs, thisDotConstructor).
		Srcref(superCall)
}

func isSingleSpreadOfArguments(nodeList []rhino.Node) bool {
	if len(nodeList) != 1 {
		return false
	}
	_, isSpread := ast.AsSpread(nodeList[0])
	return isSpread && nodeList[0].GetOnlyChild().MatchesName("arguments")
}

func (e *Es6ConvertSuperConstructorCalls) replaceNativeErrorSuperCall(
	superCall rhino.Node, newSuperCall rhino.Node, compilerInput *CompilerInput) {
	superStatement := NodeUtil_GetEnclosingStatement(superCall)
	body := superStatement.GetParent()
	jutil.CheckStateMsg1(body.IsBlock(), "%s", body)

	thisType := AstFactory_Type_(newSuperCall)
	tmpErrorName := es6SuperConstructorTmpError + "$" + e.uniqueIdSupplier.GetUniqueId(compilerInput)

	getError :=
		rhino.IR_Var(e.astFactory.CreateName(tmpErrorName, thisType)).SrcrefTreeIfMissing(superCall)
	getError.InsertBefore(superStatement)

	getTmpError :=
		e.astFactory.CreateAssign(e.astFactory.CreateName(tmpErrorName, thisType), newSuperCall)
	copyMessage :=
		e.astFactory.CreateAssign(
			e.astFactory.CreateGetProp(
				e.astFactory.CreateThis(thisType), "message", AstFactory_Type_OfColor(colors.StandardColors_STRING)),
			e.astFactory.CreateGetProp(
				e.astFactory.CreateName(tmpErrorName, thisType),
				"message",
				AstFactory_Type_OfColor(colors.StandardColors_STRING)))

	setStack :=
		e.astFactory.CreateAnd(
			e.astFactory.CreateIn(
				e.astFactory.CreateString("stack"), e.astFactory.CreateName(tmpErrorName, thisType)),
			e.astFactory.CreateAssign(
				e.astFactory.CreateGetProp(
					e.astFactory.CreateThis(thisType), "stack", AstFactory_Type_OfColor(colors.StandardColors_STRING)),
				e.astFactory.CreateGetProp(
					e.astFactory.CreateName(tmpErrorName, thisType),
					"stack",
					AstFactory_Type_OfColor(colors.StandardColors_STRING))))
	superErrorExpr := e.astFactory.
		CreateCommas(getTmpError, copyMessage, setStack, e.astFactory.CreateThis(thisType)).
		SrcrefTreeIfMissing(superCall)
	superCall.ReplaceWith(superErrorExpr)
	e.compiler.ReportChangeToEnclosingScope(superErrorExpr)
}

func (e *Es6ConvertSuperConstructorCalls) isNativeObjectClass(t *NodeTraversal, className string) bool {
	return className == "Object" && !e.isDefinedInSources(t, className)
}

func (e *Es6ConvertSuperConstructorCalls) isNativeErrorClass(t *NodeTraversal, superClassName string) bool {
	switch superClassName {
	case "AggregateError",
		"Error",
		"EvalError",
		"RangeError",
		"ReferenceError",
		"SyntaxError",
		"TypeError",
		"URIError":
		return !e.isDefinedInSources(t, superClassName)
	default:
		return false
	}
}

func (e *Es6ConvertSuperConstructorCalls) isKnownNativeClass(t *NodeTraversal, className string) bool {
	switch className {
	case "Array",
		"ArrayBuffer",
		"Boolean",
		"DataView",
		"Date",
		"Float32Array",
		"Function",
		"Generator",
		"GeneratorFunction",
		"Int16Array",
		"Int32Array",
		"Int8Array",
		"InternalError",
		"Iterator",
		"Map",
		"Number",
		"Object",
		"Promise",
		"Proxy",
		"RegExp",
		"Set",
		"String",
		"Symbol",
		"TypedArray",
		"Uint16Array",
		"Uint32Array",
		"Uint8Array",
		"Uint8ClampedArray",
		"WeakMap",
		"WeakSet":
		return !e.isDefinedInSources(t, className)
	default:
		return false
	}
}

func (e *Es6ConvertSuperConstructorCalls) isDefinedInSources(t *NodeTraversal, varName string) bool {
	objectVar := t.GetScope().GetVar(varName)
	return objectVar != nil && !objectVar.IsExtern()
}

func (e *Es6ConvertSuperConstructorCalls) updateThisToSuperThis(
	typeOfThis AstFactory_Type,
	constructorBody rhino.Node,
	superCalls []rhino.Node,
	uniqueSuperThisName string,
) {
	replacer := es6ReplaceThisWithSuperThis{
		pass:                e,
		typeOfThis:          typeOfThis,
		superCalls:          superCalls,
		uniqueSuperThisName: uniqueSuperThisName,
	}
	NodeTraversal_Traverse(e.compiler, constructorBody, &replacer)
}

type es6ReplaceThisWithSuperThis struct {
	pass                *Es6ConvertSuperConstructorCalls
	typeOfThis          AstFactory_Type
	superCalls          []rhino.Node
	uniqueSuperThisName string
}

func (*es6ReplaceThisWithSuperThis) visitTokens() *ast.TokenSet {
	return &es6ReplaceThisWithSuperThisVisitTokens
}

func (c *es6ReplaceThisWithSuperThis) ShouldTraverse(nodeTraversal *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if slices.Contains(c.superCalls, n) {
		return false
	}
	if fn, ok := ast.AsFunction(n); ok && !fn.IsArrow() {
		return false
	}
	return true
}

func (c *es6ReplaceThisWithSuperThis) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsThis() {
		superThis := c.pass.astFactory.CreateName(c.uniqueSuperThisName, AstFactory_Type_(n)).Srcref(n)
		n.ReplaceWith(superThis)
	} else if ret, ok := ast.AsReturn(n); ok {
		if _, hasVal := ret.Value(); !hasVal {
			n.AddChildToFront(c.pass.astFactory.CreateName(c.uniqueSuperThisName, c.typeOfThis).Srcref(n))
		}
	}
}

func (e *Es6ConvertSuperConstructorCalls) typeOfThisForConstructor(constructor rhino.Node) AstFactory_Type {
	jutil.CheckArgumentMsg1(constructor.IsFunction(), "%s", constructor)
	constructorType := constructor.GetColor()
	if constructorType != nil && len(constructorType.GetInstanceColors()) != 0 {
		return AstFactory_Type_OfColor(colors.Color_CreateUnionSlice(constructorType.GetInstanceColors()))
	}
	return AstFactory_Type_OfColor(colors.StandardColors_UNKNOWN)
}

func (e *Es6ConvertSuperConstructorCalls) superClassQNameNode(constructor rhino.Node) rhino.Node {
	className := NodeUtil_GetNameNode(constructor).GetQualifiedName()
	constructorStatement := rhino.CheckNotNull(NodeUtil_GetEnclosingStatement(constructor))

	var superClassNameNode rhino.Node
	for statement := constructorStatement.GetNext(); !statement.IsNil(); statement = statement.GetNext() {
		superClassNameNode = e.superClassNameNodeIfInheritsStatement(statement, className)
		if !superClassNameNode.IsNil() {
			break
		}
	}

	return rhino.CheckNotNullMsg(superClassNameNode, "$jscomp.inherits() call not found.")
}

func (e *Es6ConvertSuperConstructorCalls) superClassNameNodeIfInheritsStatement(
	statement rhino.Node, className *string) rhino.Node {
	exprStmt, ok := ast.AsExprStmt(statement)
	if !ok {
		return rhino.NilNode
	}
	call, ok := ast.AsCall(exprStmt.Expr().Node())
	if !ok {
		return rhino.NilNode
	}
	jscompDotInherits := call.Callee().Node()
	if !e.jscompInherits.Matches(jscompDotInherits) {
		return rhino.NilNode
	}
	classNameNode := rhino.CheckNotNull(jscompDotInherits.GetNext())
	if classNameNode.MatchesQualifiedName(*className) {
		return rhino.CheckNotNull(classNameNode.GetNext())
	}
	return rhino.NilNode
}

// Java: Es6ConvertSuperConstructorCalls.setGlobalNamespace(GlobalNamespace)
func (e *Es6ConvertSuperConstructorCalls) SetGlobalNamespace(globalNamespace *GlobalNamespace) {
	e.globalNamespace = globalNamespace
}
