// Port of com/google/javascript/jscomp/InstrumentAsyncContext.java @ closure-compiler v20261021

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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// InstrumentAsyncContext instruments `await` and `yield` for the
// `AsyncContext` polyfill.
type InstrumentAsyncContext struct {
	start RuntimeJsLibManager_JsLibField

	compiler         AbstractCompiler
	astFactory       *AstFactory
	uniqueIdSupplier *UniqueIdSupplier

	// Whether `await`s should be instrumented. When we're transpiling them away, we can skip
	// instrumenting them completely (the transpiled `yield`s don't need instrumentation because we
	// can rely on `Promise.then` to take care of it).
	shouldInstrumentAwait bool

	tryFunctionStack     []rhino.Node
	needsInstrumentation map[rhino.Node]*asyncContextFuncInfo
	alreadyInstrumented  map[rhino.Node]bool
	hasSuper             map[rhino.Node]bool
}

// We prefix all internal symbols with a characteristic "ᵃᶜ" (U+1D43, U+1D9C) to
// significantly reduce the chance of conflicts with existing names in the code.
const (
	asyncContextFactoryPrefix = "ᵃᶜfactory"
	asyncContextSuspendPrefix = "ᵃᶜsuspend"
	asyncContextResumePrefix  = "ᵃᶜresume"
)

var (
	asyncContextArgumentsOriginalName = "arguments"
	asyncContextArgumentsID           = uint32(strpool.Intern("arguments"))
)

var asyncContextArgsRenamerVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_NAME)
	return s
}()

func NewInstrumentAsyncContext(compiler AbstractCompiler, shouldInstrumentAwait bool) *InstrumentAsyncContext {
	return &InstrumentAsyncContext{
		needsInstrumentation:  make(map[rhino.Node]*asyncContextFuncInfo),
		alreadyInstrumented:   make(map[rhino.Node]bool),
		hasSuper:              make(map[rhino.Node]bool),
		compiler:              compiler,
		astFactory:            compiler.CreateAstFactory(),
		uniqueIdSupplier:      compiler.GetUniqueIdSupplier(),
		shouldInstrumentAwait: shouldInstrumentAwait,
		start:                 compiler.GetRuntimeJsLibManager().GetJsLibField("$jscomp.asyncContextStart"),
	}
}

func (i *InstrumentAsyncContext) Process(externs rhino.Node, root rhino.Node) {
	jutil.CheckStateMsg1(i.compiler.GetLifeCycleStage().IsNormalized(), "%s", i.compiler.GetLifeCycleStage())
	NodeTraversal_Traverse(i.compiler, root, i)
}

func (i *InstrumentAsyncContext) tryFunctionStackPeek() rhino.Node {
	if len(i.tryFunctionStack) == 0 {
		return rhino.NilNode
	}
	return i.tryFunctionStack[len(i.tryFunctionStack)-1]
}

func (i *InstrumentAsyncContext) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	// The top of tryFunctionStack points to the innermost try or function block so we can find an
	// immediately-enclosing try without walking ancestors across function boundaries.
	if n.IsTry() || n.IsFunction() {
		i.tryFunctionStack = append(i.tryFunctionStack, n)
	}
	return true
}

func (i *InstrumentAsyncContext) isReentrance(n rhino.Node) bool {
	if n.IsYield() {
		return true
	}
	if !i.shouldInstrumentAwait {
		return false
	}
	return n.IsForAwaitOf() || n.IsAwait()
}

func (i *InstrumentAsyncContext) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if nm, ok := ast.AsName(n); ok && strings.HasPrefix(nm.Ident(), asyncContextFactoryPrefix) {
		i.alreadyInstrumented[t.GetEnclosingFunction()] = true
		return
	}
	if n.IsSuper() {
		i.hasSuper[t.GetEnclosingFunction()] = true
		return
	}
	if i.isReentrance(n) {
		i.instrumentReentrance(t, n, parent)
		return
	}
	if n == i.tryFunctionStackPeek() {
		i.tryFunctionStack = i.tryFunctionStack[:len(i.tryFunctionStack)-1]
		if i.alreadyInstrumented[n] {
			return
		}
		ctx := i.needsInstrumentation[n]
		if ctx == nil && n.IsGeneratorFunction() {
			ctx = i.newFunctionContext(t.GetInput())
			i.needsInstrumentation[n] = ctx
		}
		if ctx == nil {
			return
		}
		if tryNode, ok := ast.AsTry(n); ok {
			i.instrumentTry(tryNode, ctx)
		} else {
			fn := ast.MustFunction(n)
			if fn.IsGenerator() {
				i.instrumentGeneratorFunction(t, fn, ctx)
			} else {
				jutil.CheckState(fn.IsAsync())
				i.instrumentAsyncFunction(fn, ctx)
			}
		}
	}
}

// instrumentReentrance wraps `await EXPR` / `yield EXPR` as `ᵃᶜresume(await ᵃᶜsuspend(EXPR))`
// and schedules instrumentation for the enclosing function and try block.
func (i *InstrumentAsyncContext) instrumentReentrance(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	enclosingTry := i.tryFunctionStackPeek()
	if !enclosingTry.IsNil() && !enclosingTry.IsTry() {
		enclosingTry = rhino.NilNode
	}
	enclosingFunction := t.GetEnclosingFunction()
	if enclosingFunction.IsNil() ||
		i.alreadyInstrumented[enclosingFunction] ||
		parent.IsFunction() || // `async (...) => await ...` - no reentrance
		(parent.IsReturn() && enclosingTry.IsNil()) { // `return await ...` outside a `try`
		return
	}

	functionContext := i.needsInstrumentation[enclosingFunction]
	if functionContext == nil {
		functionContext = i.newFunctionContext(t.GetInput())
		i.needsInstrumentation[enclosingFunction] = functionContext
	}
	if !enclosingTry.IsNil() {
		i.needsInstrumentation[enclosingTry] = functionContext
	}

	if fao, ok := ast.AsForAwaitOf(n); ok {
		i.instrumentForAwaitOf(fao, parent, functionContext)
		return
	}

	placeholder := rhino.IR_Empty()
	arg := n.GetFirstChild()
	if arg.IsNil() {
		n.AddChildToBack(functionContext.suspend(i.astFactory).SrcrefTreeIfMissing(n))
	} else {
		arg.ReplaceWith(placeholder)
		placeholder.ReplaceWith(functionContext.suspendWithInner(i.astFactory, arg).SrcrefTreeIfMissing(arg))
	}

	n.ReplaceWith(placeholder)
	placeholder.ReplaceWith(functionContext.resumeWithInner(i.astFactory, n).SrcrefTreeIfMissing(n))
}

func (i *InstrumentAsyncContext) instrumentForAwaitOf(fao ast.ForAwaitOf, parent rhino.Node, functionContext *asyncContextFuncInfo) {
	// for await (const x of expr()) { use(x); }
	//   =>
	// for await (const x of suspend(expr())) {
	//   resume();
	//   try { use(x); } finally { suspend(); }
	// }
	// resume();
	placeholder := rhino.IR_Empty()
	arg := fao.Iterable().Node()
	arg.ReplaceWith(placeholder)
	placeholder.ReplaceWith(functionContext.suspendWithInner(i.astFactory, arg).SrcrefTreeIfMissing(arg))

	body := fao.Body().Node()
	jutil.CheckState(body.IsBlock())
	body.ReplaceWith(placeholder)
	placeholder.ReplaceWith(
		rhino.IR_BlockOfNodes(
			rhino.IR_ExprResult(functionContext.resume(i.astFactory).SrcrefTreeIfMissing(arg)),
			rhino.IR_TryFinally(
				body,
				rhino.IR_BlockOfNode(rhino.IR_ExprResult(functionContext.suspend(i.astFactory).SrcrefTreeIfMissing(arg))))))

	target := fao.Node()
	if !parent.IsBlock() {
		block := rhino.IR_Block()
		target.ReplaceWith(block)
		block.AddChildToFront(target)
		target = block
	}
	rhino.IR_ExprResult(functionContext.resume(i.astFactory).SrcrefTreeIfMissing(arg)).InsertAfter(target)
}

func (i *InstrumentAsyncContext) instrumentTry(tryNode ast.Try, functionContext *asyncContextFuncInfo) {
	var block rhino.Node
	if catchClause, ok := tryNode.Catch(); ok {
		block = catchClause.Body().Node()
	} else if finallyBlock, ok := tryNode.Finally(); ok {
		block = finallyBlock.Node()
	}
	jutil.CheckState(block.IsBlock())
	block.AddChildToFront(rhino.IR_ExprResult(functionContext.resume(i.astFactory).SrcrefTreeIfMissing(block)))
}

func (i *InstrumentAsyncContext) instrumentGeneratorFunction(t *NodeTraversal, fn ast.Function, functionContext *asyncContextFuncInfo) {
	f := fn.Node()
	if !fn.Body().Node().HasChildren() {
		return
	}
	if i.hasSuper[f] {
		i.instrumentGeneratorMethod(t, fn, functionContext)
		return
	}

	fn.SetGenerator(false)
	generatorBody := fn.Body().Node().Detach()
	inner := i.astFactory.CreateEmptyGeneratorFunction(AstFactory_Type_(f))
	innerFn := ast.MustFunction(inner)

	if fn.IsAsync() {
		fn.SetAsync(false)
		innerFn.SetAsync(true)
	}
	innerFn.Body().Node().ReplaceWith(generatorBody)

	generatorBody = i.addFinallySuspend(generatorBody, functionContext)
	generatorBody.AddChildToFront(rhino.IR_ExprResult(functionContext.resume(i.astFactory)))
	newOuterBody := rhino.IR_BlockOfNodes(
		i.createFactorySuspendCall(functionContext.factoryName),
		rhino.IR_ReturnNodeWithExpr(i.astFactory.CreateCallWithUnknownType(inner)))
	i.addSuspendResumeVarsAfter(newOuterBody.GetFirstChild(), functionContext)
	f.AddChildToBack(newOuterBody.SrcrefTreeIfMissing(generatorBody))

	thisContext := NewThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext(
		newOuterBody, false, i.uniqueIdSupplier.GetUniqueId(t.GetInput()))
	updater := NewThisAndArgumentsReferenceUpdater(i.compiler, thisContext, i.astFactory)
	NodeTraversal_Traverse(i.compiler, generatorBody, updater)
	thisContext.AddVarDeclarations(i.compiler, i.astFactory, t.GetCurrentScript())

	i.compiler.ReportChangeToChangeScope(f)
	i.compiler.ReportChangeToChangeScope(inner)
	NodeUtil_AddFeatureToScript(
		t.GetCurrentScript(), parser.FeatureSet_Feature_BLOCK_SCOPED_FUNCTION_DECLARATION, i.compiler)
}

func (i *InstrumentAsyncContext) instrumentGeneratorMethod(t *NodeTraversal, fn ast.Function, functionContext *asyncContextFuncInfo) {
	f := fn.Node()
	fp := f.GetParent()
	innerMethodNamePrefix := ""
	if mfd, ok := ast.AsMemberFunctionDef(fp); ok {
		innerMethodNamePrefix = mfd.Key()
	}
	innerMethodName := innerMethodNamePrefix + "$jscomp$" + i.uniqueIdSupplier.GetUniqueId(t.GetInput())
	outerParams := fn.Params().Node()
	innerParams := outerParams.CloneTree()

	innerFactoryName := functionContext.factoryName
	outerFactoryName := i.astFactory.CreateNameWithUnknownType(i.uniqueIdSupplier.UniqueVarName(asyncContextFactoryPrefix, t.GetInput()))
	innerParams.AddChildToFront(innerFactoryName)
	unknownType := AstFactory_Type_(innerFactoryName)

	generatorBody := i.addFinallySuspend(fn.Body().Node(), functionContext)
	generatorBody.AddChildToFront(rhino.IR_ExprResult(functionContext.resume(i.astFactory)))

	argumentsParam := rhino.NilNode
	argumentsReference := rhino.NilNode
	renamer := asyncContextArgsRenamer{outer: i}
	NodeTraversal_Traverse(i.compiler, generatorBody, &renamer)
	if renamer.argumentsName != "" {
		argumentsReference = i.astFactory.CreateArgumentsReference()
		argumentsParam = i.astFactory.CreateName(renamer.argumentsName, AstFactory_Type_(argumentsReference))
		argumentsParam.InsertAfter(innerFactoryName)
	}

	outerBody := rhino.IR_Block()
	generatorBody.ReplaceWith(outerBody)

	innerCall := i.astFactory.CreateCallWithUnknownType(
		i.astFactory.CreateGetProp(
			i.astFactory.CreateThis(unknownType), innerMethodName, AstFactory_Type_(f)))
	innerParam := innerParams.GetFirstChild()
	innerCall.AddChildToBack(outerFactoryName.CloneTree())
	innerParam = innerParam.GetNext()
	if !innerParam.IsNil() && innerParam == argumentsParam {
		innerCall.AddChildToBack(argumentsReference)
	}
	for outerParam := outerParams.GetFirstChild(); !outerParam.IsNil(); outerParam = outerParam.GetNext() {
		argument, parameter := i.simplifyParameter(outerParam.CloneTree(), t)
		oldOuterParam := outerParam
		outerParam = parameter
		oldOuterParam.ReplaceWith(outerParam)
		innerCall.AddChildToBack(argument)
	}

	newInnerMethod := i.astFactory.CreateFunction("", innerParams, generatorBody, AstFactory_Type_(f))
	newInnerFn := ast.MustFunction(newInnerMethod)
	np := rhino.IR_MemberFunctionDef(innerMethodName, newInnerMethod)
	np.InsertAfter(fp)
	newInnerFn.NameNode().Node().MakeNonIndexable()

	newInnerFn.SetAsync(fn.IsAsync())
	newInnerFn.SetGenerator(true)
	fn.SetGenerator(false)
	fn.SetAsync(false)

	np.SetStaticMember(fp.IsStaticMember())
	if jsdoc := fp.GetJSDocInfo(); jsdoc != nil {
		builder := rhino.JSDocInfo_Builder_()
		if jsdoc.IsNoCollapse() {
			builder.RecordNoCollapse()
		}
		if jsdoc.IsNoInline() {
			builder.RecordNoInline()
		}
		np.SetJSDocInfo(builder.Build())
	}

	outerBody.AddChildToBack(i.createFactorySuspendCall(outerFactoryName))
	i.addSuspendResumeVarsAtStartOf(generatorBody, innerFactoryName, functionContext)
	outerBody.AddChildToBack(rhino.IR_ReturnNodeWithExpr(innerCall))

	f.SrcrefTreeIfMissing(f)
	np.SrcrefTreeIfMissing(f)
	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_MEMBER_DECLARATIONS, i.compiler)
	i.compiler.ReportChangeToChangeScope(f)
	i.compiler.ReportChangeToChangeScope(newInnerMethod)
	i.compiler.ReportChangeToEnclosingScope(fp)
	NodeUtil_AddFeatureToScript(
		t.GetCurrentScript(), parser.FeatureSet_Feature_BLOCK_SCOPED_FUNCTION_DECLARATION, i.compiler)
}

func (i *InstrumentAsyncContext) simplifyParameter(param rhino.Node, t *NodeTraversal) (argument, parameter rhino.Node) {
	switch {
	case param.IsRest():
		NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_SPREAD_EXPRESSIONS, i.compiler)
		arg, p := i.simplifyParameter(param.RemoveFirstChild(), t)
		return rhino.IR_IterSpread(arg).SrcrefTreeIfMissing(param),
			rhino.IR_IterRest(p).SrcrefTreeIfMissing(param)
	case param.IsDefaultValue():
		return i.simplifyParameter(param.RemoveFirstChild(), t)
	case param.IsArrayPattern() || param.IsObjectPattern():
		newParam := i.astFactory.
			CreateName(i.uniqueIdSupplier.UniqueVarName("$jscomp$param$", t.GetInput()), AstFactory_Type_(param)).
			SrcrefTreeIfMissing(param)
		return newParam.CloneTree(), newParam
	case param.IsName():
		newParam := i.astFactory.
			CreateName(i.uniqueIdSupplier.UniqueVarName(param.GetString()+"$jscomp$param$", t.GetInput()), AstFactory_Type_(param)).
			SrcrefTreeIfMissing(param)
		return newParam.CloneTree(), newParam
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Unexpected parameter: " + param.String()))
	}
}

type asyncContextArgsRenamer struct {
	outer         *InstrumentAsyncContext
	argumentsName string
}

var _ nodeTraversal_visitTokens = (*asyncContextArgsRenamer)(nil)

func (*asyncContextArgsRenamer) visitTokens() *ast.TokenSet {
	return &asyncContextArgsRenamerVisitTokens
}

func (r *asyncContextArgsRenamer) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	fn, ok := ast.AsFunction(n)
	return !ok || fn.IsArrow()
}

func (r *asyncContextArgsRenamer) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if nm, ok := ast.AsName(n); ok && n.GetStringID() == asyncContextArgumentsID {
		if r.argumentsName == "" {
			r.argumentsName = r.outer.uniqueIdSupplier.UniqueVarName("$jscomp$arguments$", t.GetInput())
		}
		nm.SetIdent(r.argumentsName)
		if r.outer.compiler.GetOptions().PreservesDetailedSourceInfo() {
			n.SetOriginalName(&asyncContextArgumentsOriginalName)
		}
		n.MakeNonIndexable()
		nm.SetConstantName(true)
	}
}

func (i *InstrumentAsyncContext) instrumentAsyncFunction(fn ast.Function, functionContext *asyncContextFuncInfo) {
	newBlock := i.addFinallySuspend(fn.Body().Node(), functionContext)
	newBlock.AddChildToFront(
		i.createFactoryCall(functionContext.factoryName).SrcrefTreeIfMissing(newBlock))
	i.addSuspendResumeVarsAfter(newBlock.GetFirstChild(), functionContext)
	i.compiler.ReportChangeToChangeScope(fn.Node())
}

func (i *InstrumentAsyncContext) addFinallySuspend(block rhino.Node, functionContext *asyncContextFuncInfo) rhino.Node {
	placeholder := rhino.IR_Empty()
	block.ReplaceWith(placeholder)
	finallyBlock := rhino.IR_BlockOfNode(rhino.IR_ExprResult(functionContext.suspend(i.astFactory)))
	newBlock := rhino.IR_BlockOfNode(rhino.IR_TryFinally(block, finallyBlock)).SrcrefTreeIfMissing(block)
	placeholder.ReplaceWith(newBlock)
	return newBlock
}

func (i *InstrumentAsyncContext) addSuspendResumeVarsAfter(node rhino.Node, functionContext *asyncContextFuncInfo) {
	resumeCall := i.astFactory.CreateCallWithUnknownType(
		functionContext.factoryName.CloneNode(), i.astFactory.CreateNumber(1))
	rhino.IR_VarWithValue(functionContext.resumeName, resumeCall).SrcrefTreeIfMissing(node).InsertAfter(node)
	suspendCall := i.astFactory.CreateCallWithUnknownType(functionContext.factoryName.CloneNode())
	rhino.IR_VarWithValue(functionContext.suspendName, suspendCall).SrcrefTreeIfMissing(node).InsertAfter(node)
}

func (i *InstrumentAsyncContext) addSuspendResumeVarsAtStartOf(
	node rhino.Node, factoryName rhino.Node, functionContext *asyncContextFuncInfo,
) {
	resumeCall := i.astFactory.CreateCallWithUnknownType(factoryName.CloneNode(), i.astFactory.CreateNumber(1))
	node.AddChildToFront(rhino.IR_VarWithValue(functionContext.resumeName, resumeCall).SrcrefTreeIfMissing(node))
	suspendCall := i.astFactory.CreateCallWithUnknownType(factoryName.CloneNode())
	node.AddChildToFront(rhino.IR_VarWithValue(functionContext.suspendName, suspendCall).SrcrefTreeIfMissing(node))
}

func (i *InstrumentAsyncContext) createFactoryCall(factoryName rhino.Node, arg ...rhino.Node) rhino.Node {
	call := i.astFactory.CreateCallWithUnknownType(i.astFactory.CreateQNameWithUnknownTypeOfField(i.start), arg...)
	return rhino.IR_VarWithValue(factoryName, call)
}

func (i *InstrumentAsyncContext) createFactorySuspendCall(factoryName rhino.Node) rhino.Node {
	return i.createFactoryCall(factoryName, i.astFactory.CreateNumber(1))
}

type asyncContextFuncInfo struct {
	suspendName rhino.Node
	resumeName  rhino.Node
	factoryName rhino.Node
}

func (c *asyncContextFuncInfo) suspend(af *AstFactory) rhino.Node {
	return af.CreateCall(c.suspendName.CloneNode(), AstFactory_Type_(c.suspendName))
}

func (c *asyncContextFuncInfo) suspendWithInner(af *AstFactory, inner rhino.Node) rhino.Node {
	return af.CreateCall(c.suspendName.CloneNode(), AstFactory_Type_(inner), inner)
}

func (c *asyncContextFuncInfo) resume(af *AstFactory) rhino.Node {
	return af.CreateCall(c.resumeName.CloneNode(), AstFactory_Type_(c.resumeName))
}

func (c *asyncContextFuncInfo) resumeWithInner(af *AstFactory, inner rhino.Node) rhino.Node {
	return af.CreateCall(c.resumeName.CloneNode(), AstFactory_Type_(inner), inner)
}

func (i *InstrumentAsyncContext) newFunctionContext(input *CompilerInput) *asyncContextFuncInfo {
	uniqueSuffix := i.uniqueIdSupplier.GetUniqueId(input)
	return &asyncContextFuncInfo{
		suspendName: i.astFactory.CreateNameWithUnknownType(asyncContextSuspendPrefix + uniqueSuffix),
		resumeName:  i.astFactory.CreateNameWithUnknownType(asyncContextResumePrefix + uniqueSuffix),
		factoryName: i.astFactory.CreateNameWithUnknownType(asyncContextFactoryPrefix + uniqueSuffix),
	}
}
