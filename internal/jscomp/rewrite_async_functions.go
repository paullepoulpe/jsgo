// Port of com/google/javascript/jscomp/RewriteAsyncFunctions.java @ closure-compiler v20261021

/*
 * Copyright 2016 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// RewriteAsyncFunctions converts async functions to valid ES6 generator
// functions code.
//
// This pass must run before the passes that transpile let declarations,
// arrow functions, and generators.
type RewriteAsyncFunctions struct {
	namespace    rhino.StaticScope
	contextStack []*rewriteAsyncLexicalContext
	compiler     AbstractCompiler
	astFactory   *AstFactory
}

const (
	asyncFuncArgumentsPrefix   = "$jscomp$async$arguments$"
	asyncFuncThisPrefix        = "$jscomp$async$this$"
	asyncFuncNewTargetPrefix   = "$jscomp$async$new$target$"
	asyncFuncSuperGetterPrefix = "$jscomp$async$super$get$"
)

type asyncSuperPropertyWrapperInfo struct {
	propertyName                    string
	firstInstanceOfSuperDotProperty rhino.Node
	wrapperFunctionName             string
	wrapperFunctionReturnType       AstFactory_Type
}

func (i asyncSuperPropertyWrapperInfo) createWrapperFunctionCallNode(astFactory *AstFactory) rhino.Node {
	nameNode := astFactory.CreateName(i.wrapperFunctionName, AstFactory_Type_OfColor(colors.StandardColors_TOP_OBJECT))
	return astFactory.CreateCall(nameNode, i.wrapperFunctionReturnType)
}

type asyncContextKind uint8

const (
	asyncContextRoot asyncContextKind = iota
	asyncContextParamList
	asyncContextFunction
)

type rewriteAsyncLexicalContext struct {
	kind            asyncContextKind
	pass            *RewriteAsyncFunctions
	contextRootNode rhino.Node
	uniqueId        string

	// Set when kind == asyncContextParamList.
	functionContext *rewriteAsyncLexicalContext

	// Set when kind == asyncContextFunction: points to the enclosing async FunctionContext whose
	// `this`, `arguments`, `new.target`, and `super` bindings apply here, or nil if none.
	asyncThisAndArgumentsContext *rewriteAsyncLexicalContext
	// Ordered list of distinct `super.prop` wrappers needed by this async context.
	superPropertyWrappers []asyncSuperPropertyWrapperInfo

	mustAddAsyncThisVariable      bool
	typeOfThis                    AstFactory_Type
	mustAddAsyncArgumentsVariable bool
	mustAddAsyncNewTargetVariable bool
	typeOfNewTarget               AstFactory_Type
}

func newAsyncRootContext(pass *RewriteAsyncFunctions, contextRootNode rhino.Node, uniqueId string) *rewriteAsyncLexicalContext {
	return &rewriteAsyncLexicalContext{
		kind:            asyncContextRoot,
		pass:            pass,
		contextRootNode: rhino.CheckNotNull(contextRootNode),
		uniqueId:        uniqueId,
	}
}

func newAsyncParamListContext(
	pass *RewriteAsyncFunctions,
	functionContext *rewriteAsyncLexicalContext,
	contextRootNode rhino.Node,
	uniqueId string,
) *rewriteAsyncLexicalContext {
	return &rewriteAsyncLexicalContext{
		kind:            asyncContextParamList,
		pass:            pass,
		contextRootNode: rhino.CheckNotNull(contextRootNode),
		uniqueId:        uniqueId,
		functionContext: jutil.CheckNotNull(functionContext),
	}
}

func newAsyncFunctionContext(
	pass *RewriteAsyncFunctions,
	outerAsyncContext *rewriteAsyncLexicalContext,
	fn ast.Function,
	uniqueId string,
) *rewriteAsyncLexicalContext {
	c := &rewriteAsyncLexicalContext{
		kind:            asyncContextFunction,
		pass:            pass,
		contextRootNode: rhino.CheckNotNull(fn.Node()),
		uniqueId:        uniqueId,
	}
	if fn.IsAsync() {
		if fn.IsArrow() && outerAsyncContext != nil {
			c.asyncThisAndArgumentsContext = outerAsyncContext
		} else {
			c.asyncThisAndArgumentsContext = c
		}
	} else if fn.IsArrow() {
		c.asyncThisAndArgumentsContext = outerAsyncContext
	}
	return c
}

func (c *rewriteAsyncLexicalContext) asyncContextForParameterList() *rewriteAsyncLexicalContext {
	asyncContext := c.functionContext.asyncThisAndArgumentsContext
	if asyncContext == c.functionContext {
		return nil
	}
	return asyncContext
}

func (c *rewriteAsyncLexicalContext) contextForNode(t *NodeTraversal, n rhino.Node) *rewriteAsyncLexicalContext {
	switch c.kind {
	case asyncContextRoot:
		if fn, ok := ast.AsFunction(n); ok {
			uniqueId := c.pass.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())
			return newAsyncFunctionContext(c.pass, nil, fn, uniqueId)
		}
		return c
	case asyncContextParamList:
		if fn, ok := ast.AsFunction(n); ok {
			uniqueId := c.pass.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())
			return newAsyncFunctionContext(c.pass, c.asyncContextForParameterList(), fn, uniqueId)
		}
		return c
	case asyncContextFunction:
		if n == c.contextRootNode {
			return c
		}
		if fn, ok := ast.AsFunction(n); ok {
			uniqueId := c.pass.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())
			return newAsyncFunctionContext(c.pass, c.asyncThisAndArgumentsContext, fn, uniqueId)
		}
		if _, ok := ast.AsParamList(n); ok {
			uniqueId := c.pass.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())
			return newAsyncParamListContext(c.pass, c, n, uniqueId)
		}
		return c
	default:
		panic("unreachable asyncContextKind")
	}
}

func (c *rewriteAsyncLexicalContext) recordAsyncThisReplacement(typeOfThis AstFactory_Type) {
	c.mustAddAsyncThisVariable = true
	c.typeOfThis = typeOfThis
}

func (c *rewriteAsyncLexicalContext) createThisVariableReference(typeOfThis AstFactory_Type) rhino.Node {
	c.recordAsyncThisReplacement(typeOfThis)
	return c.pass.astFactory.CreateName(asyncFuncThisPrefix+c.uniqueId, typeOfThis)
}

func (c *rewriteAsyncLexicalContext) createNewTargetVariableReference(typeOfNewTarget AstFactory_Type) rhino.Node {
	c.mustAddAsyncNewTargetVariable = true
	c.typeOfNewTarget = typeOfNewTarget
	return c.pass.astFactory.CreateName(asyncFuncNewTargetPrefix+c.uniqueId, typeOfNewTarget)
}

func (c *rewriteAsyncLexicalContext) getOrCreateSuperPropertyWrapperInfo(superDotPropertyNode rhino.Node) asyncSuperPropertyWrapperInfo {
	gp, ok := ast.AsGetProp(superDotPropertyNode)
	jutil.CheckArgumentMsg1(ok, "%s", superDotPropertyNode)
	superNode := gp.Object().Node()
	jutil.CheckArgumentMsg1(superNode.IsSuper(), "%s", superNode)

	propertyName := gp.Property()
	propertyType := superDotPropertyNode.GetColor()
	for _, info := range c.superPropertyWrappers {
		if info.propertyName == propertyName {
			existingColor := info.firstInstanceOfSuperDotProperty.GetColor()
			jutil.CheckStateMsg2(
				jutil.ObjectsEquals(existingColor, propertyType),
				"Previous reference type: %s differs from current reference type: %s",
				existingColor,
				propertyType)
			return info
		}
	}

	wrapperFunctionName := asyncFuncSuperGetterPrefix +
		c.pass.compiler.GetUniqueIdSupplier().GetUniqueId(c.pass.compiler.GetInput(NodeUtil_GetInputId(superDotPropertyNode))) +
		"$" + propertyName
	info := asyncSuperPropertyWrapperInfo{
		propertyName:                    propertyName,
		firstInstanceOfSuperDotProperty: superDotPropertyNode,
		wrapperFunctionName:             wrapperFunctionName,
		wrapperFunctionReturnType:       AstFactory_Type_(superDotPropertyNode),
	}
	c.superPropertyWrappers = append(c.superPropertyWrappers, info)
	return info
}

func (c *rewriteAsyncLexicalContext) createWrapperArrowFunction(wrapperInfo asyncSuperPropertyWrapperInfo) rhino.Node {
	superDotProperty := wrapperInfo.firstInstanceOfSuperDotProperty.CloneTree()
	return c.pass.astFactory.CreateZeroArgArrowFunctionForExpression(
		c.pass.astFactory.CreateBlock(c.pass.astFactory.CreateReturn(superDotProperty)))
}

func (c *rewriteAsyncLexicalContext) visit(t *NodeTraversal, n rhino.Node) {
	switch c.kind {
	case asyncContextRoot:
		// Outside any async function.
	case asyncContextParamList:
		if c.asyncContextForParameterList() != nil {
			c.functionContext.visitFunctionNode(t, n)
		}
	case asyncContextFunction:
		c.visitFunctionNode(t, n)
	}
}

func (c *rewriteAsyncLexicalContext) visitFunctionNode(t *NodeTraversal, n rhino.Node) {
	r := c.pass
	if c.contextRootNode == n {
		if fn, ok := ast.AsFunction(n); ok && fn.IsAsync() {
			r.convertAsyncFunction(t, c, fn)
		}
		return
	}
	asyncCtx := c.asyncThisAndArgumentsContext
	if asyncCtx == nil {
		return
	}

	switch n.GetToken() {
	case rhino.Token_NAME:
		nm := ast.MustName(n)
		if nm.Ident() == "arguments" {
			nm.SetIdent(asyncFuncArgumentsPrefix + asyncCtx.uniqueId)
			if r.compiler.GetLifeCycleStage().IsNormalized() {
				nm.SetConstantName(true)
			}
			asyncCtx.mustAddAsyncArgumentsVariable = true
			r.compiler.ReportChangeToChangeScope(c.contextRootNode)
		}
	case rhino.Token_THIS:
		n.ReplaceWith(asyncCtx.createThisVariableReference(AstFactory_Type_(n)))
		r.compiler.ReportChangeToChangeScope(c.contextRootNode)
	case rhino.Token_NEW_TARGET:
		n.ReplaceWith(asyncCtx.createNewTargetVariableReference(AstFactory_Type_(n)))
		r.compiler.ReportChangeToChangeScope(c.contextRootNode)
	case rhino.Token_SUPER:
		parent := n.GetParent()
		if !parent.IsGetProp() {
			r.compiler.Report(
				JSError_MakeWithNode(parent, TranspilationUtil_CANNOT_CONVERT_YET, "super expression"))
		} else if NodeUtil_IsLValue(parent) {
			r.compiler.Report(
				JSError_MakeWithNode(
					parent.GetParent(),
					TranspilationUtil_CANNOT_CONVERT_YET,
					"assignment to super property"))
		}
		superDotProperty := parent
		wrapperInfo := asyncCtx.getOrCreateSuperPropertyWrapperInfo(superDotProperty)

		getPropReplacement := wrapperInfo.createWrapperFunctionCallNode(r.astFactory)
		grandparent := superDotProperty.GetParent()
		if call, ok := ast.AsCall(grandparent); ok && call.Callee().Node() == superDotProperty {
			getPropReplacement = r.astFactory.CreateGetPropWithUnknownType(getPropReplacement, "call")
			thisAlias := r.astFactory.
				CreateThisAliasReferenceForEs6Class(
					asyncFuncThisPrefix+asyncCtx.uniqueId,
					NodeUtil_GetEnclosingClass(asyncCtx.contextRootNode)).
				Srcref(superDotProperty)
			thisAlias.InsertAfter(superDotProperty)
			asyncCtx.recordAsyncThisReplacement(AstFactory_Type_(thisAlias))
		} else if ttl, ok := ast.AsTaggedTemplateLit(grandparent); ok && ttl.Tag().Node() == superDotProperty {
			thisAlias := r.astFactory.
				CreateThisAliasReferenceForEs6Class(
					asyncFuncThisPrefix+asyncCtx.uniqueId,
					NodeUtil_GetEnclosingClass(asyncCtx.contextRootNode)).
				Srcref(superDotProperty)
			asyncCtx.recordAsyncThisReplacement(AstFactory_Type_(thisAlias))
			getPropReplacement = r.astFactory.CreateCallWithUnknownType(
				r.astFactory.CreateGetPropWithUnknownType(getPropReplacement, "bind"),
				thisAlias)
			ttl.SetFreeCall(true)
		}
		getPropReplacement.SrcrefTree(superDotProperty)
		superDotProperty.ReplaceWith(getPropReplacement)
		r.compiler.ReportChangeToChangeScope(c.contextRootNode)
	case rhino.Token_AWAIT:
		aw := ast.MustAwait(n)
		n.ReplaceWith(r.astFactory.CreateYield(AstFactory_Type_(n), aw.Value().Node().Detach()))
	}
}

var rewriteAsyncFunctionsTranspiledFeatures = parser.FeatureSet_BARE_MINIMUM.With(parser.FeatureSet_Feature_ASYNC_FUNCTIONS)

func RewriteAsyncFunctions_Create(compiler AbstractCompiler) *RewriteAsyncFunctions {
	return &RewriteAsyncFunctions{
		compiler:   jutil.CheckNotNull(compiler),
		astFactory: jutil.CheckNotNull(compiler.CreateAstFactory()),
		namespace:  jutil.CheckNotNull(compiler.GetTranspilationNamespace()),
	}
}

func (r *RewriteAsyncFunctions) Process(externs rhino.Node, root rhino.Node) {
	TranspilationPasses_ProcessTranspile(r.compiler, root, rewriteAsyncFunctionsTranspiledFeatures, r)
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(r.compiler, root, rewriteAsyncFunctionsTranspiledFeatures)
}

func (r *RewriteAsyncFunctions) contextStackPeek() *rewriteAsyncLexicalContext {
	if len(r.contextStack) == 0 {
		return nil
	}
	return r.contextStack[len(r.contextStack)-1]
}

func (r *RewriteAsyncFunctions) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if parent.IsNil() {
		jutil.CheckState(len(r.contextStack) == 0)
		uniqueId := r.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())
		r.contextStack = append(r.contextStack, newAsyncRootContext(r, n, uniqueId))
	} else {
		parentContext := r.contextStackPeek()
		nodeContext := parentContext.contextForNode(t, n)
		if nodeContext != parentContext {
			r.contextStack = append(r.contextStack, nodeContext)
		}
	}
	return true
}

func (r *RewriteAsyncFunctions) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	context := r.contextStackPeek()
	context.visit(t, n)
	if context.contextRootNode == n {
		if len(r.contextStack) == 0 {
			panic(jutil.NewNoSuchElementException())
		}
		r.contextStack = r.contextStack[:len(r.contextStack)-1]
	}
}

func (r *RewriteAsyncFunctions) convertAsyncFunction(
	t *NodeTraversal,
	functionContext *rewriteAsyncLexicalContext,
	originalFunction ast.Function,
) {
	originalFunction.SetAsync(false)
	originalBody := originalFunction.Body().Node()
	if originalFunction.Node().IsFromExterns() {
		if !NodeUtil_IsEmptyBlock(originalBody) {
			originalBody.ReplaceWith(r.astFactory.CreateBlock())
			NodeUtil_MarkFunctionsDeleted(originalBody, r.compiler)
		}
		return
	}
	newBody := r.astFactory.CreateBlock()
	originalBody.ReplaceWith(newBody)

	enclosingScript := t.GetCurrentScript()
	asyncCtx := functionContext.asyncThisAndArgumentsContext
	if functionContext.mustAddAsyncThisVariable {
		newBody.AddChildToBack(
			r.astFactory.CreateSingleConstNameDeclaration(
				asyncFuncThisPrefix+asyncCtx.uniqueId,
				r.astFactory.CreateThis(functionContext.typeOfThis)))
		NodeUtil_AddFeatureToScript(enclosingScript, parser.FeatureSet_Feature_CONST_DECLARATIONS, r.compiler)
	}
	if functionContext.mustAddAsyncArgumentsVariable {
		newBody.AddChildToBack(
			r.astFactory.CreateArgumentsAliasDeclaration(
				asyncFuncArgumentsPrefix + asyncCtx.uniqueId))
		NodeUtil_AddFeatureToScript(enclosingScript, parser.FeatureSet_Feature_CONST_DECLARATIONS, r.compiler)
	}
	if functionContext.mustAddAsyncNewTargetVariable {
		newBody.AddChildToBack(
			r.astFactory.CreateSingleConstNameDeclaration(
				asyncFuncNewTargetPrefix+asyncCtx.uniqueId,
				r.astFactory.CreateNewTarget(functionContext.typeOfNewTarget)))
		NodeUtil_AddFeatureToScript(enclosingScript, parser.FeatureSet_Feature_CONST_DECLARATIONS, r.compiler)
	}
	for _, wrapperInfo := range functionContext.superPropertyWrappers {
		arrowFunction := functionContext.createWrapperArrowFunction(wrapperInfo)
		newBody.AddChildToBack(
			r.astFactory.CreateSingleConstNameDeclaration(
				wrapperInfo.wrapperFunctionName, arrowFunction))
		r.compiler.ReportChangeToChangeScope(arrowFunction)
		NodeUtil_AddFeatureToScript(enclosingScript, parser.FeatureSet_Feature_ARROW_FUNCTIONS, r.compiler)
		NodeUtil_AddFeatureToScript(enclosingScript, parser.FeatureSet_Feature_CONST_DECLARATIONS, r.compiler)
	}

	if !originalBody.IsBlock() {
		originalBody = r.astFactory.
			CreateBlock(r.astFactory.CreateReturn(originalBody)).
			SrcrefTreeIfMissing(originalBody)
	}
	generatorFunction := r.astFactory.CreateZeroArgGeneratorFunction("", originalBody, nil)
	r.compiler.ReportChangeToChangeScope(generatorFunction)
	NodeUtil_AddFeatureToScript(enclosingScript, parser.FeatureSet_Feature_GENERATORS, r.compiler)

	newBody.AddChildToBack(
		r.astFactory.CreateReturn(
			r.astFactory.CreateJscompAsyncExecutePromiseGeneratorFunctionCall(
				r.namespace, generatorFunction)))

	newBody.SrcrefTreeIfMissing(originalBody)
	r.compiler.ReportChangeToEnclosingScope(newBody)
}
