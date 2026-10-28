// Port of com/google/javascript/jscomp/RewriteAsyncIteration.java @ closure-compiler v20261021

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
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// RewriteAsyncIteration converts async generator functions into a function
// returning a new $jscomp.AsyncGeneratorWrapper around an inner generator and
// lowers `for await (... of ...)` loops.
type RewriteAsyncIteration struct {
	asyncGeneratorWrapper RuntimeJsLibManager_JsLibField
	actionRecordName      RuntimeJsLibManager_JsLibField
	actionEnumAwait       RuntimeJsLibManager_JsLibField
	actionEnumYield       RuntimeJsLibManager_JsLibField
	actionEnumYieldStar   RuntimeJsLibManager_JsLibField

	nextForAwaitId int
	compiler       AbstractCompiler
	contextStack   []asyncIterContext
	astFactory     *AstFactory
	namespace      rhino.StaticScope
}

var rewriteAsyncIterationTranspiledFeatures = parser.FeatureSet_BARE_MINIMUM.WithFeatures(
	parser.FeatureSet_Feature_ASYNC_GENERATORS,
	parser.FeatureSet_Feature_FOR_AWAIT_OF,
)

var RewriteAsyncIteration_CANNOT_CONVERT_ASYNCGEN = DiagnosticType_Error(
	"JSC_CANNOT_CONVERT_ASYNCGEN", "Cannot convert async generator. {0}")

const (
	forAwaitIteratorTempName    = "$jscomp$forAwait$tempIterator"
	forAwaitResultTempName      = "$jscomp$forAwait$tempResult"
	forAwaitErrorResultTempName = "$jscomp$forAwait$errResult"
	forAwaitCatchParamTempName  = "$jscomp$forAwait$catchErrParam"
	forAwaitReturnFnTempName    = "$jscomp$forAwait$retFn"

	asyncIterThisVarName      = "$jscomp$asyncIter$this$"
	asyncIterArgumentsVarName = "$jscomp$asyncIter$arguments"
	asyncIterSuperGetterPref  = "$jscomp$asyncIter$super$get$"
)

// asyncIterContext is a value-type stack frame tracking the enclosing function
// and (when inside an async generator) its this/super/arguments alias state.
type asyncIterContext struct {
	contextRoot rhino.Node
	function    rhino.Node
	scopeState  *asyncIterScopeState
}

func (c asyncIterContext) mustReplaceThisSuperArgs() bool {
	return c.scopeState != nil
}

// asyncIterScopeState is allocated only for async generator functions to record
// references to `this`, `arguments`, and `super.prop` that must be hoisted out
// of the inner generator wrapper.
type asyncIterScopeState struct {
	usedSuperProperties []rhino.Node
	thisNodeToAdd       rhino.Node
	usedArguments       bool
	uniqueId            string
}

func RewriteAsyncIteration_Create(compiler AbstractCompiler) *RewriteAsyncIteration {
	jutil.CheckNotNull(compiler)
	runtimeLibManager := compiler.GetRuntimeJsLibManager()
	return &RewriteAsyncIteration{
		compiler:              compiler,
		astFactory:            jutil.CheckNotNull(compiler.CreateAstFactory()),
		namespace:             jutil.CheckNotNull(compiler.GetTranspilationNamespace()),
		asyncGeneratorWrapper: runtimeLibManager.GetJsLibField("$jscomp.AsyncGeneratorWrapper"),
		actionRecordName:      runtimeLibManager.GetJsLibField("$jscomp.AsyncGeneratorWrapper$ActionRecord"),
		actionEnumAwait:       runtimeLibManager.GetJsLibField("$jscomp.AsyncGeneratorWrapper$ActionEnum.AWAIT_VALUE"),
		actionEnumYield:       runtimeLibManager.GetJsLibField("$jscomp.AsyncGeneratorWrapper$ActionEnum.YIELD_VALUE"),
		actionEnumYieldStar:   runtimeLibManager.GetJsLibField("$jscomp.AsyncGeneratorWrapper$ActionEnum.YIELD_STAR"),
	}
}

func (r *RewriteAsyncIteration) contextTop() asyncIterContext {
	if len(r.contextStack) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	return r.contextStack[len(r.contextStack)-1]
}

func (r *RewriteAsyncIteration) popContext() asyncIterContext {
	top := r.contextTop()
	r.contextStack = r.contextStack[:len(r.contextStack)-1]
	return top
}

func (r *RewriteAsyncIteration) Process(externs rhino.Node, root rhino.Node) {
	jutil.CheckState(len(r.contextStack) == 0)
	r.contextStack = append(r.contextStack, asyncIterContext{contextRoot: rhino.CheckNotNull(root)})
	TranspilationPasses_ProcessTranspile(r.compiler, root, rewriteAsyncIterationTranspiledFeatures, r)
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(r.compiler, root, rewriteAsyncIterationTranspiledFeatures)
	jutil.CheckState(r.contextTop().function.IsNil())
	r.popContext()
	jutil.CheckState(len(r.contextStack) == 0)
}

func (r *RewriteAsyncIteration) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if fn, ok := ast.AsFunction(n); ok {
		prev := r.contextTop()
		var state *asyncIterScopeState
		if fn.IsArrow() {
			state = prev.scopeState
		} else {
			uniqueId := r.compiler.GetUniqueIdSupplier().GetUniqueId(r.compiler.GetInput(NodeUtil_GetInputId(n)))
			if fn.IsAsync() && fn.IsGenerator() {
				state = &asyncIterScopeState{uniqueId: uniqueId}
			}
		}
		r.contextStack = append(r.contextStack, asyncIterContext{
			contextRoot: n,
			function:    n,
			scopeState:  state,
		})
	} else if _, ok := ast.AsParamList(n); ok {
		prev := r.contextTop()
		var state *asyncIterScopeState
		if prev.function.IsArrowFunction() {
			state = prev.scopeState
		}
		r.contextStack = append(r.contextStack, asyncIterContext{
			contextRoot: n,
			function:    prev.function,
			scopeState:  state,
		})
	}
	return true
}

func (r *RewriteAsyncIteration) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	ctx := r.contextTop()
	switch n.GetToken() {
	case rhino.Token_PARAM_LIST:
		jutil.CheckStateMsg1(n == ctx.contextRoot, "%s", n)
		r.popContext()
	case rhino.Token_FUNCTION:
		jutil.CheckState(n == ctx.contextRoot)
		if n.IsAsyncGeneratorFunction() {
			r.convertAsyncGenerator(ast.MustFunction(n))
			r.prependTempVarDeclarations(ctx, t)
		}
		r.popContext()
	case rhino.Token_AWAIT:
		rhino.CheckNotNull(ctx.function)
		if ctx.function.IsAsyncGeneratorFunction() {
			r.convertAwaitOfAsyncGenerator(n)
		}
	case rhino.Token_YIELD:
		rhino.CheckNotNull(ctx.function)
		if ctx.function.IsAsyncGeneratorFunction() {
			r.convertYieldOfAsyncGenerator(n)
		}
	case rhino.Token_RETURN:
		rhino.CheckNotNull(ctx.function)
		if ctx.function.IsAsyncGeneratorFunction() {
			r.convertReturnOfAsyncGenerator(n)
		}
	case rhino.Token_FOR_AWAIT_OF:
		rhino.CheckNotNull(ctx.function)
		jutil.CheckState(ctx.function.IsAsyncFunction())
		r.replaceForAwaitOf(ctx, ast.MustForAwaitOf(n))
		NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_CONST_DECLARATIONS, r.compiler)
	case rhino.Token_THIS:
		if ctx.mustReplaceThisSuperArgs() {
			r.replaceThis(ctx, n)
		}
	case rhino.Token_NAME:
		if ctx.mustReplaceThisSuperArgs() && n.MatchesName("arguments") {
			r.replaceArguments(ctx, n)
		}
	case rhino.Token_SUPER:
		if ctx.mustReplaceThisSuperArgs() {
			r.replaceSuper(ctx, n, parent)
		}
	}
}

func (r *RewriteAsyncIteration) newActionRecord(enumField RuntimeJsLibManager_JsLibField, expr rhino.Node) rhino.Node {
	return r.astFactory.CreateNewNode(
		r.astFactory.CreateQNameOfField(r.namespace, r.actionRecordName),
		r.astFactory.CreateQNameOfField(r.namespace, enumField),
		expr)
}

func (r *RewriteAsyncIteration) convertAsyncGenerator(originalFn ast.Function) {
	originalFunction := originalFn.Node()
	asyncGeneratorWrapperRef := r.astFactory.CreateQNameOfField(r.namespace, r.asyncGeneratorWrapper)
	innerFunction := r.astFactory.CreateEmptyAsyncGeneratorWrapperArgument(nil)

	innerBlock := originalFn.Body().Node().Detach()
	ast.MustFunction(innerFunction).Body().Node().ReplaceWith(innerBlock)

	outerBlock := r.astFactory.CreateBlock(
		r.astFactory.CreateReturn(
			r.astFactory.CreateNewNode(
				asyncGeneratorWrapperRef,
				r.astFactory.CreateCall(innerFunction, AstFactory_Type_OfColorId(colors.StandardColors_GENERATOR_ID)))))
	originalFunction.AddChildToBack(outerBlock)

	originalFn.SetAsync(false)
	originalFn.SetGenerator(false)
	originalFunction.SrcrefTreeIfMissing(originalFunction)
	r.compiler.ReportChangeToChangeScope(originalFunction)
	r.compiler.ReportChangeToChangeScope(innerFunction)
}

func (r *RewriteAsyncIteration) convertAwaitOfAsyncGenerator(awaitNode rhino.Node) {
	expression := awaitNode.RemoveFirstChild()
	rhino.CheckNotNullMsg(expression, "await needs an expression")
	newActionRecord := r.newActionRecord(r.actionEnumAwait, expression).SrcrefTreeIfMissing(awaitNode)
	awaitNode.AddChildToFront(newActionRecord)
	awaitNode.SetToken(rhino.Token_YIELD)
}

func (r *RewriteAsyncIteration) convertYieldOfAsyncGenerator(yieldNode rhino.Node) {
	yd := ast.MustYield(yieldNode)
	expression := yieldNode.RemoveFirstChild()
	enumField := r.actionEnumYield
	if yd.IsYieldAll() {
		rhino.CheckNotNull(expression)
		enumField = r.actionEnumYieldStar
	} else if expression.IsNil() {
		expression = NodeUtil_NewUndefinedNode(rhino.NilNode)
	}
	newActionRecord := r.newActionRecord(enumField, expression).SrcrefTreeIfMissing(yieldNode)
	yieldNode.AddChildToFront(newActionRecord)
	yd.SetYieldAll(false)
}

func (r *RewriteAsyncIteration) convertReturnOfAsyncGenerator(returnNode rhino.Node) {
	expression := returnNode.RemoveFirstChild()
	if expression.IsNil() {
		expression = NodeUtil_NewUndefinedNode(rhino.NilNode)
	}
	newActionRecord := r.newActionRecord(r.actionEnumYield, expression).SrcrefTreeIfMissing(returnNode)
	returnNode.AddChildToFront(newActionRecord)
}

func (r *RewriteAsyncIteration) replaceForAwaitOf(ctx asyncIterContext, fao ast.ForAwaitOf) {
	forAwaitOf := fao.Node()
	id := strconv.Itoa(r.nextForAwaitId)
	r.nextForAwaitId++
	iteratorTempName := forAwaitIteratorTempName + id
	resultTempName := forAwaitResultTempName + id
	errorResultTempName := forAwaitErrorResultTempName + id
	catchErrorParamTempName := forAwaitCatchParamTempName + id
	returnFuncTempName := forAwaitReturnFnTempName + id

	jutil.CheckStateMsg(forAwaitOf.HasParent(), "Cannot replace parentless for-await-of")

	forAwaitOfParent := forAwaitOf.GetParent()
	replacementPoint := forAwaitOf
	if lbl, ok := ast.AsLabel(forAwaitOfParent); ok {
		jutil.CheckStateMsg1(lbl.Body().Node() == forAwaitOf, "%s", forAwaitOfParent)
		replacementPoint = forAwaitOfParent
	}

	lhs := forAwaitOf.RemoveFirstChild()
	rhs := forAwaitOf.RemoveFirstChild()
	originalBody := forAwaitOf.RemoveFirstChild()

	initializer := r.astFactory.
		CreateSingleVarNameDeclarationWithValue(
			iteratorTempName, r.astFactory.CreateJSCompMakeAsyncIteratorCall(rhs, r.namespace)).
		SrcrefTreeIfMissing(rhs)

	iterableResultType := AstFactory_Type_OfColor(colors.StandardColors_TOP_OBJECT)

	breakIfDone := r.astFactory.CreateIf(
		r.astFactory.CreateGetProp(
			r.astFactory.CreateName(resultTempName, iterableResultType),
			"done",
			AstFactory_Type_OfColor(colors.StandardColors_BOOLEAN)),
		r.astFactory.CreateBlock(r.astFactory.CreateBreak()))

	var lhsAssignment rhino.Node
	var resultType AstFactory_Type
	if lhs.IsValidAssignmentTarget() {
		resultType = AstFactory_Type_(lhs)
		lhsAssignment = r.astFactory.ExprResult(
			r.astFactory.CreateAssign(
				lhs,
				r.astFactory.CreateGetProp(
					r.astFactory.CreateName(resultTempName, iterableResultType),
					"value",
					resultType)))
	} else if NodeUtil_IsNameDeclaration(lhs) {
		declarationTarget := lhs.GetFirstChild()
		if declarationTarget.IsName() {
			resultType = AstFactory_Type_(declarationTarget)
		} else {
			dlhs, ok := ast.AsDestructuringLhs(declarationTarget)
			jutil.CheckStateMsg1(ok, "%s", declarationTarget)
			resultType = AstFactory_Type_(dlhs.Pattern().Node())
		}
		declarationTarget.AddChildToBack(
			r.astFactory.CreateGetProp(
				r.astFactory.CreateName(resultTempName, iterableResultType), "value", resultType))
		lhsAssignment = lhs
	} else {
		panic(jutil.NewAssertionError("unexpected for-await-of lhs"))
	}
	lhsAssignment.SrcrefTreeIfMissing(lhs)

	errorResDecl := r.astFactory.CreateSingleVarNameDeclaration(errorResultTempName).SrcrefTreeIfMissing(forAwaitOf)
	tempResultDecl := r.astFactory.CreateSingleVarNameDeclaration(resultTempName).SrcrefTreeIfMissing(forAwaitOf)
	returnFuncDecl := r.astFactory.CreateSingleVarNameDeclaration(returnFuncTempName).SrcrefTreeIfMissing(forAwaitOf)

	resetTempResult := r.astFactory.ExprResult(
		r.astFactory.CreateAssignToName(resultTempName, r.astFactory.CreateUndefinedValue()))

	resultDeclaration := r.astFactory.ExprResult(
		r.astFactory.CreateAssignToName(
			resultTempName,
			r.constructAwaitNextResult(ctx, iteratorTempName, resultType, iterableResultType)))

	wrappedBody := originalBody
	if !wrappedBody.IsBlock() {
		wrappedBody = r.astFactory.CreateBlock(wrappedBody).Srcref(wrappedBody)
	}

	newForLoop := r.astFactory.CreateFor(
		r.astFactory.CreateEmpty(),
		r.astFactory.CreateEmpty(),
		r.astFactory.CreateEmpty(),
		r.astFactory.CreateBlock(
			resetTempResult,
			resultDeclaration,
			breakIfDone,
			lhsAssignment,
			wrappedBody))

	if lbl, ok := ast.AsLabel(replacementPoint); ok {
		newForLoop = r.astFactory.CreateLabel(lbl.Name().Node().CloneNode(), newForLoop)
	}

	tryNode := r.astFactory.CreateBlock(newForLoop)
	initializer.InsertBefore(newForLoop)

	catchNode := r.createOuterCatch(catchErrorParamTempName, errorResultTempName)
	finallyNode := r.createOuterFinally(
		ctx,
		iterableResultType,
		resultType,
		resultTempName,
		returnFuncTempName,
		iteratorTempName,
		errorResultTempName)

	tryCatchFinally := r.astFactory.CreateTryCatchFinally(tryNode, catchNode, finallyNode)
	replacementPoint.ReplaceWith(tryCatchFinally)
	tryCatchFinally.SrcrefTreeIfMissing(replacementPoint)
	errorResDecl.InsertBefore(tryCatchFinally)
	tempResultDecl.InsertBefore(tryCatchFinally)
	returnFuncDecl.InsertBefore(tryCatchFinally)

	r.compiler.ReportChangeToEnclosingScope(tryCatchFinally)
}

func (r *RewriteAsyncIteration) createOuterCatch(catchErrorParamTempName, errorResultTempName string) rhino.Node {
	catchBodyStmt := r.astFactory.ExprResult(
		r.astFactory.CreateAssignToName(
			errorResultTempName,
			r.astFactory.CreateObjectLit(
				r.astFactory.CreateStringKey(
					"error", r.astFactory.CreateNameWithUnknownType(catchErrorParamTempName)))))
	return r.astFactory.CreateCatch(
		r.astFactory.CreateNameWithUnknownType(catchErrorParamTempName),
		r.astFactory.CreateBlock(catchBodyStmt))
}

func (r *RewriteAsyncIteration) createOuterFinally(
	ctx asyncIterContext,
	iterableResultType AstFactory_Type,
	resultType AstFactory_Type,
	resultTempName string,
	returnFuncTempName string,
	iteratorTempName string,
	errorResultTempName string,
) rhino.Node {
	tmpResNameNode := r.astFactory.CreateNameWithUnknownType(resultTempName)
	tmpResDoneGetProp := r.astFactory.CreateGetProp(
		r.astFactory.CreateName(resultTempName, iterableResultType),
		"done",
		AstFactory_Type_OfColor(colors.StandardColors_BOOLEAN))
	and := r.astFactory.CreateAnd(tmpResNameNode, r.astFactory.CreateNot(tmpResDoneGetProp))

	unknownType := AstFactory_Type_OfColor(colors.StandardColors_UNKNOWN)
	assign := r.astFactory.CreateAssign(
		r.astFactory.CreateNameWithUnknownType(returnFuncTempName),
		r.astFactory.CreateGetProp(
			r.astFactory.CreateName(iteratorTempName, resultType),
			"return",
			unknownType))

	ifCond := r.astFactory.CreateAnd(and, assign)
	var awaitOrYieldStmt rhino.Node
	if ctx.function.IsAsyncGeneratorFunction() {
		call := r.astFactory.CreateCall(
			r.astFactory.CreateGetPropWithUnknownType(
				r.astFactory.CreateName(returnFuncTempName, unknownType),
				"call"),
			unknownType,
			r.astFactory.CreateName(iteratorTempName, resultType))
		awaitOrYieldStmt = r.astFactory.ExprResult(
			r.astFactory.CreateYield(iterableResultType, r.newActionRecord(r.actionEnumAwait, call)))
	} else {
		call := r.astFactory.CreateCall(
			r.astFactory.CreateGetPropWithUnknownType(
				r.astFactory.CreateName(returnFuncTempName, unknownType),
				"call"),
			AstFactory_Type_OfColorId(colors.StandardColors_PROMISE_ID),
			r.astFactory.CreateName(iteratorTempName, resultType))
		awaitOrYieldStmt = r.astFactory.ExprResult(r.astFactory.CreateAwait(iterableResultType, call))
	}

	innerTryBlock := r.astFactory.CreateBlock(
		r.astFactory.CreateIf(ifCond, r.astFactory.CreateBlock(awaitOrYieldStmt)))

	throwStmt := r.astFactory.CreateThrow(
		r.astFactory.CreateGetPropWithUnknownType(
			r.astFactory.CreateNameWithUnknownType(errorResultTempName), "error"))
	innerFinallyBlock := r.astFactory.CreateBlock(
		r.astFactory.CreateIf(
			r.astFactory.CreateNameWithUnknownType(errorResultTempName),
			r.astFactory.CreateBlock(throwStmt)))

	return r.astFactory.CreateBlock(r.astFactory.CreateTryFinally(innerTryBlock, innerFinallyBlock))
}

func (r *RewriteAsyncIteration) constructAwaitNextResult(
	ctx asyncIterContext,
	iteratorTempName string,
	iteratorType AstFactory_Type,
	iterableResultType AstFactory_Type,
) rhino.Node {
	iteratorTemp := r.astFactory.CreateName(iteratorTempName, iteratorType)
	if ctx.function.IsAsyncGeneratorFunction() {
		call := r.astFactory.CreateCallWithUnknownType(
			r.astFactory.CreateGetPropWithUnknownType(iteratorTemp, "next"))
		return r.astFactory.CreateYield(iterableResultType, r.newActionRecord(r.actionEnumAwait, call))
	}
	call := r.astFactory.CreateCall(
		r.astFactory.CreateGetPropWithUnknownType(iteratorTemp, "next"),
		AstFactory_Type_OfColorId(colors.StandardColors_PROMISE_ID))
	return r.astFactory.CreateAwait(iterableResultType, call)
}

func (r *RewriteAsyncIteration) replaceThis(ctx asyncIterContext, n rhino.Node) {
	state := ctx.scopeState
	n.ReplaceWith(
		r.astFactory.
			CreateName(asyncIterThisVarName+state.uniqueId, AstFactory_Type_(n)).
			Srcref(n))
	state.thisNodeToAdd = r.astFactory.CreateThis(AstFactory_Type_(n))
	r.compiler.ReportChangeToChangeScope(ctx.function)
}

func (r *RewriteAsyncIteration) replaceArguments(ctx asyncIterContext, n rhino.Node) {
	n.ReplaceWith(r.astFactory.CreateName(asyncIterArgumentsVarName, AstFactory_Type_(n)).Srcref(n))
	ctx.scopeState.usedArguments = true
	r.compiler.ReportChangeToChangeScope(ctx.function)
}

func (r *RewriteAsyncIteration) replaceSuper(ctx asyncIterContext, n rhino.Node, parent rhino.Node) {
	gp, ok := ast.AsGetProp(parent)
	if !ok {
		r.compiler.Report(
			JSError_MakeWithNode(
				parent,
				RewriteAsyncIteration_CANNOT_CONVERT_ASYNCGEN,
				"super only allowed with getprop (like super.foo(), not super['foo']())"))
		return
	}
	state := ctx.scopeState
	propertyName := gp.Property()
	propertyReplacementNameText := asyncIterSuperGetterPref + propertyName

	getPropReplacement := r.astFactory.CreateCall(
		r.astFactory.CreateName(propertyReplacementNameText, AstFactory_Type_OfColor(colors.StandardColors_TOP_OBJECT)),
		AstFactory_Type_(parent))
	grandparent := parent.GetParent()
	if call, ok := ast.AsCall(grandparent); ok && call.Callee().Node() == parent {
		getPropReplacement = r.astFactory.CreateGetPropWithUnknownType(getPropReplacement, "call")
		state.thisNodeToAdd = r.astFactory.CreateThisForEs6ClassMember(ctx.contextRoot.GetParent())
		r.astFactory.
			CreateName(asyncIterThisVarName+state.uniqueId, AstFactory_Type_(state.thisNodeToAdd)).
			Srcref(parent).
			InsertAfter(parent)
	} else if ttl, ok := ast.AsTaggedTemplateLit(grandparent); ok && ttl.Tag().Node() == parent {
		state.thisNodeToAdd = r.astFactory.CreateThisForEs6ClassMember(ctx.contextRoot.GetParent())
		thisAlias := r.astFactory.
			CreateName(asyncIterThisVarName+state.uniqueId, AstFactory_Type_(state.thisNodeToAdd)).
			Srcref(parent)
		getPropReplacement = r.astFactory.CreateCallWithUnknownType(
			r.astFactory.CreateGetPropWithUnknownType(getPropReplacement, "bind"), thisAlias)
		ttl.SetFreeCall(true)
	}
	getPropReplacement.SrcrefTree(parent)
	parent.ReplaceWith(getPropReplacement)
	if !slices.Contains(state.usedSuperProperties, parent) {
		state.usedSuperProperties = append(state.usedSuperProperties, parent)
	}
	r.compiler.ReportChangeToChangeScope(ctx.function)
}

func (r *RewriteAsyncIteration) prependTempVarDeclarations(ctx asyncIterContext, t *NodeTraversal) {
	state := ctx.scopeState
	if state == nil {
		return
	}
	if state.thisNodeToAdd.IsNil() && !state.usedArguments && len(state.usedSuperProperties) == 0 {
		return
	}

	function := ctx.function
	block := ast.MustFunction(function).Body().Node()
	prefixBlock := r.astFactory.CreateBlock()

	if !state.thisNodeToAdd.IsNil() {
		prefixBlock.AddChildToBack(
			r.astFactory.
				CreateSingleConstNameDeclaration(asyncIterThisVarName+state.uniqueId, state.thisNodeToAdd).
				SrcrefTree(block))
	}
	if state.usedArguments {
		prefixBlock.AddChildToBack(
			r.astFactory.
				CreateSingleConstNameDeclaration(asyncIterArgumentsVarName, r.astFactory.CreateArgumentsReference()).
				SrcrefTree(block))
	}
	for _, replacedMethodReference := range state.usedSuperProperties {
		prefixBlock.AddChildToBack(r.createSuperMethodReferenceGetter(replacedMethodReference, t))
	}
	prefixBlock.SrcrefTreeIfMissing(block)
	block.AddChildrenToFront(prefixBlock.RemoveChildren())

	r.compiler.ReportChangeToChangeScope(function)
	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_CONST_DECLARATIONS, r.compiler)
}

func (r *RewriteAsyncIteration) createSuperMethodReferenceGetter(replacedMethodReference rhino.Node, t *NodeTraversal) rhino.Node {
	gp := ast.MustGetProp(replacedMethodReference)
	typeOfSuper := AstFactory_Type_(gp.Object().Node())
	superReference := r.astFactory.CreateSuper(typeOfSuper)
	replacedMethodName := gp.Property()
	arrowFunction := r.astFactory.CreateZeroArgArrowFunctionForExpression(
		r.astFactory.CreateBlock(
			r.astFactory.CreateReturn(
				r.astFactory.CreateGetProp(
					superReference, replacedMethodName, AstFactory_Type_(replacedMethodReference)))))
	r.compiler.ReportChangeToChangeScope(arrowFunction)
	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_ARROW_FUNCTIONS, r.compiler)
	return r.astFactory.CreateSingleConstNameDeclaration(asyncIterSuperGetterPref+replacedMethodName, arrowFunction)
}
