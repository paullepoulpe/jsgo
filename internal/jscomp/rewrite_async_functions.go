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
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

var asyncArgumentsID = uint32(strpool.Intern("arguments"))

var rewriteAsyncVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_SCRIPT)
	s.Add(rhino.Token_ROOT)
	s.Add(rhino.Token_FUNCTION)
	s.Add(rhino.Token_PARAM_LIST)
	s.Add(rhino.Token_NAME)
	s.Add(rhino.Token_THIS)
	s.Add(rhino.Token_NEW_TARGET)
	s.Add(rhino.Token_SUPER)
	s.Add(rhino.Token_AWAIT)
	return s
}()

// RewriteAsyncFunctions converts async functions to valid ES6 generator
// functions code.
//
// This pass must run before the passes that transpile let declarations,
// arrow functions, and generators.
type RewriteAsyncFunctions struct {
	namespace       rhino.StaticScope
	contextStack    []rewriteAsyncLexicalContext
	contextStackBuf [8]rewriteAsyncLexicalContext
	compiler        AbstractCompiler
	astFactory      *AstFactory
}

const (
	asyncFuncArgumentsPrefix   = "$jscomp$async$arguments$"
	asyncFuncThisPrefix        = "$jscomp$async$this$"
	asyncFuncNewTargetPrefix   = "$jscomp$async$new$target$"
	asyncFuncSuperGetterPrefix = "$jscomp$async$super$get$"
)

type asyncSuperPropertyWrapperInfo struct {
	propertyID                      uint32
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

// rewriteAsyncLexicalContext stores lexical context entries by value in
// RewriteAsyncFunctions.contextStack, referencing ancestor contexts via stack
// indices (-1 when absent) so slice growth never invalidates interior pointers.
type rewriteAsyncLexicalContext struct {
	kind            asyncContextKind
	contextRootNode rhino.Node
	uniqueId        string

	// Set when kind == asyncContextParamList (-1 otherwise).
	functionContextIdx int

	// Set when kind == asyncContextFunction (-1 if none): stack index of the
	// enclosing async FunctionContext whose `this`, `arguments`, `new.target`,
	// and `super` bindings apply here.
	asyncThisAndArgumentsContextIdx int
	// Ordered list of distinct `super.prop` wrappers needed by this async context.
	superPropertyWrappers []asyncSuperPropertyWrapperInfo

	mustAddAsyncThisVariable      bool
	typeOfThis                    AstFactory_Type
	mustAddAsyncArgumentsVariable bool
	mustAddAsyncNewTargetVariable bool
	typeOfNewTarget               AstFactory_Type
}

func newAsyncRootContext(contextRootNode rhino.Node, uniqueId string) rewriteAsyncLexicalContext {
	return rewriteAsyncLexicalContext{
		kind:                            asyncContextRoot,
		contextRootNode:                 rhino.CheckNotNull(contextRootNode),
		uniqueId:                        uniqueId,
		functionContextIdx:              -1,
		asyncThisAndArgumentsContextIdx: -1,
	}
}

func newAsyncParamListContext(
	functionContextIdx int,
	contextRootNode rhino.Node,
	uniqueId string,
) rewriteAsyncLexicalContext {
	jutil.CheckState(functionContextIdx >= 0)
	return rewriteAsyncLexicalContext{
		kind:                            asyncContextParamList,
		contextRootNode:                 rhino.CheckNotNull(contextRootNode),
		uniqueId:                        uniqueId,
		functionContextIdx:              functionContextIdx,
		asyncThisAndArgumentsContextIdx: -1,
	}
}

func newAsyncFunctionContext(
	selfIdx int,
	outerAsyncIdx int,
	fn ast.Function,
	uniqueId string,
) rewriteAsyncLexicalContext {
	asyncIdx := -1
	if fn.IsAsync() {
		if fn.IsArrow() && outerAsyncIdx >= 0 {
			asyncIdx = outerAsyncIdx
		} else {
			asyncIdx = selfIdx
		}
	} else if fn.IsArrow() {
		asyncIdx = outerAsyncIdx
	}
	return rewriteAsyncLexicalContext{
		kind:                            asyncContextFunction,
		contextRootNode:                 rhino.CheckNotNull(fn.Node()),
		uniqueId:                        uniqueId,
		functionContextIdx:              -1,
		asyncThisAndArgumentsContextIdx: asyncIdx,
	}
}

func (r *RewriteAsyncFunctions) asyncContextIdxForParameterList(paramCtxIdx int) int {
	fnIdx := r.contextStack[paramCtxIdx].functionContextIdx
	asyncIdx := r.contextStack[fnIdx].asyncThisAndArgumentsContextIdx
	if asyncIdx == fnIdx {
		return -1
	}
	return asyncIdx
}

func (c *rewriteAsyncLexicalContext) recordAsyncThisReplacement(typeOfThis AstFactory_Type) {
	c.mustAddAsyncThisVariable = true
	c.typeOfThis = typeOfThis
}

func (r *RewriteAsyncFunctions) createThisVariableReference(c *rewriteAsyncLexicalContext, typeOfThis AstFactory_Type) rhino.Node {
	c.recordAsyncThisReplacement(typeOfThis)
	return r.astFactory.CreateName(asyncFuncThisPrefix+c.uniqueId, typeOfThis)
}

func (r *RewriteAsyncFunctions) createNewTargetVariableReference(c *rewriteAsyncLexicalContext, typeOfNewTarget AstFactory_Type) rhino.Node {
	c.mustAddAsyncNewTargetVariable = true
	c.typeOfNewTarget = typeOfNewTarget
	return r.astFactory.CreateName(asyncFuncNewTargetPrefix+c.uniqueId, typeOfNewTarget)
}

func (r *RewriteAsyncFunctions) getOrCreateSuperPropertyWrapperInfo(c *rewriteAsyncLexicalContext, superDotPropertyNode rhino.Node) asyncSuperPropertyWrapperInfo {
	gp, ok := ast.AsGetProp(superDotPropertyNode)
	jutil.CheckArgumentMsg1(ok, "%s", superDotPropertyNode)
	superNode := gp.Object().Node()
	jutil.CheckArgumentMsg1(superNode.IsSuper(), "%s", superNode)

	propID := superDotPropertyNode.GetStringID()
	propertyType := superDotPropertyNode.GetColor()
	for _, info := range c.superPropertyWrappers {
		if info.propertyID == propID {
			existingColor := info.firstInstanceOfSuperDotProperty.GetColor()
			jutil.CheckStateMsg2(
				jutil.ObjectsEquals(existingColor, propertyType),
				"Previous reference type: %s differs from current reference type: %s",
				existingColor,
				propertyType)
			return info
		}
	}

	propertyName := gp.Property()
	wrapperFunctionName := asyncFuncSuperGetterPrefix +
		r.compiler.GetUniqueIdSupplier().GetUniqueId(r.compiler.GetInput(NodeUtil_GetInputId(superDotPropertyNode))) +
		"$" + propertyName
	info := asyncSuperPropertyWrapperInfo{
		propertyID:                      propID,
		propertyName:                    propertyName,
		firstInstanceOfSuperDotProperty: superDotPropertyNode,
		wrapperFunctionName:             wrapperFunctionName,
		wrapperFunctionReturnType:       AstFactory_Type_(superDotPropertyNode),
	}
	c.superPropertyWrappers = append(c.superPropertyWrappers, info)
	return info
}

func (r *RewriteAsyncFunctions) createWrapperArrowFunction(wrapperInfo asyncSuperPropertyWrapperInfo) rhino.Node {
	superDotProperty := wrapperInfo.firstInstanceOfSuperDotProperty.CloneTree()
	return r.astFactory.CreateZeroArgArrowFunctionForExpression(
		r.astFactory.CreateBlock(r.astFactory.CreateReturn(superDotProperty)))
}

func (r *RewriteAsyncFunctions) visitContext(t *NodeTraversal, ctxIdx int, n rhino.Node) {
	switch r.contextStack[ctxIdx].kind {
	case asyncContextRoot:
		// Outside any async function.
	case asyncContextParamList:
		if r.asyncContextIdxForParameterList(ctxIdx) >= 0 {
			r.visitFunctionNode(t, r.contextStack[ctxIdx].functionContextIdx, n)
		}
	case asyncContextFunction:
		r.visitFunctionNode(t, ctxIdx, n)
	}
}

func (r *RewriteAsyncFunctions) visitFunctionNode(t *NodeTraversal, fnCtxIdx int, n rhino.Node) {
	c := &r.contextStack[fnCtxIdx]
	if c.contextRootNode == n {
		if fn, ok := ast.AsFunction(n); ok && fn.IsAsync() {
			r.convertAsyncFunction(t, c, fn)
		}
		return
	}
	if c.asyncThisAndArgumentsContextIdx < 0 {
		return
	}
	asyncCtx := &r.contextStack[c.asyncThisAndArgumentsContextIdx]

	switch n.GetToken() {
	case rhino.Token_NAME:
		if n.GetStringID() == asyncArgumentsID {
			nm := ast.MustName(n)
			nm.SetIdent(asyncFuncArgumentsPrefix + asyncCtx.uniqueId)
			if r.compiler.GetLifeCycleStage().IsNormalized() {
				nm.SetConstantName(true)
			}
			asyncCtx.mustAddAsyncArgumentsVariable = true
			r.compiler.ReportChangeToChangeScope(c.contextRootNode)
		}
	case rhino.Token_THIS:
		n.ReplaceWith(r.createThisVariableReference(asyncCtx, AstFactory_Type_(n)))
		r.compiler.ReportChangeToChangeScope(c.contextRootNode)
	case rhino.Token_NEW_TARGET:
		n.ReplaceWith(r.createNewTargetVariableReference(asyncCtx, AstFactory_Type_(n)))
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
		wrapperInfo := r.getOrCreateSuperPropertyWrapperInfo(asyncCtx, superDotProperty)

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
	r := &RewriteAsyncFunctions{
		compiler:   jutil.CheckNotNull(compiler),
		astFactory: jutil.CheckNotNull(compiler.CreateAstFactory()),
		namespace:  jutil.CheckNotNull(compiler.GetTranspilationNamespace()),
	}
	r.contextStack = r.contextStackBuf[:0]
	return r
}

func (*RewriteAsyncFunctions) visitTokens() *ast.TokenSet {
	return &rewriteAsyncVisitTokens
}

func (r *RewriteAsyncFunctions) Process(externs rhino.Node, root rhino.Node) {
	TranspilationPasses_ProcessTranspile(r.compiler, root, rewriteAsyncFunctionsTranspiledFeatures, r)
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(r.compiler, root, rewriteAsyncFunctionsTranspiledFeatures)
}

func (r *RewriteAsyncFunctions) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if parent.IsNil() {
		jutil.CheckState(len(r.contextStack) == 0)
		uniqueId := r.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())
		r.contextStack = append(r.contextStack, newAsyncRootContext(n, uniqueId))
		return true
	}
	tok := n.GetToken()
	if tok != rhino.Token_FUNCTION && tok != rhino.Token_PARAM_LIST {
		return true
	}
	parentIdx := len(r.contextStack) - 1
	parentCtx := &r.contextStack[parentIdx]
	switch parentCtx.kind {
	case asyncContextRoot:
		if fn, ok := ast.AsFunction(n); ok {
			uniqueId := r.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())
			r.contextStack = append(r.contextStack, newAsyncFunctionContext(len(r.contextStack), -1, fn, uniqueId))
		}
	case asyncContextParamList:
		if fn, ok := ast.AsFunction(n); ok {
			outerAsyncIdx := r.asyncContextIdxForParameterList(parentIdx)
			uniqueId := r.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())
			r.contextStack = append(r.contextStack, newAsyncFunctionContext(len(r.contextStack), outerAsyncIdx, fn, uniqueId))
		}
	case asyncContextFunction:
		if n == parentCtx.contextRootNode {
			return true
		}
		if fn, ok := ast.AsFunction(n); ok {
			outerAsyncIdx := parentCtx.asyncThisAndArgumentsContextIdx
			uniqueId := r.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())
			r.contextStack = append(r.contextStack, newAsyncFunctionContext(len(r.contextStack), outerAsyncIdx, fn, uniqueId))
		} else if _, ok := ast.AsParamList(n); ok {
			uniqueId := r.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())
			r.contextStack = append(r.contextStack, newAsyncParamListContext(parentIdx, n, uniqueId))
		}
	default:
		panic("unreachable asyncContextKind")
	}
	return true
}

func (r *RewriteAsyncFunctions) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if len(r.contextStack) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	topIdx := len(r.contextStack) - 1
	r.visitContext(t, topIdx, n)
	if r.contextStack[topIdx].contextRootNode == n {
		r.contextStack[topIdx] = rewriteAsyncLexicalContext{}
		r.contextStack = r.contextStack[:topIdx]
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
	var asyncCtx *rewriteAsyncLexicalContext
	if functionContext.asyncThisAndArgumentsContextIdx >= 0 {
		asyncCtx = &r.contextStack[functionContext.asyncThisAndArgumentsContextIdx]
	}
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
		arrowFunction := r.createWrapperArrowFunction(wrapperInfo)
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
