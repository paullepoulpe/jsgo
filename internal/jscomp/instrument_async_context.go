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

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// InstrumentAsyncContext instruments `await` and `yield` for the
// `AsyncContext` polyfill.
type InstrumentAsyncContext struct {
	start RuntimeJsLibManager_JsLibField

	compiler         AbstractCompiler
	astFactory       *AstFactory
	uniqueIdSupplier *UniqueIdSupplier

	// Whether `await`s should be instrumented.  When we're transpiling them away, we can skip
	// instrumenting them completely (the transpiled `yield`s don't need instrumentation because we
	// can rely on `Promise.then` to take care of it).
	shouldInstrumentAwait bool

	// TODO(sdh): Investigate whether async generator instrumentation can be skipped in ES2017
	// output.  They are transpiled to ordinary generators with Promise.then, so we may be able
	// to rely on runtime patching in that case as well.  If so, we'll add an additional boolean
	// constructor parameter for shouldInstrumentAsyncGenerators (and figure out whether that means
	// that both yield _and_ await can be skipped, or just await).

	// Go-only: language-adapter — the ArrayDeque used as a stack is a slice
	// whose last element is the top (Deque.peek/push/pop).
	tryFunctionStack     []*rhino.Node
	needsInstrumentation *jutil.LinkedHashMap[*rhino.Node, *InstrumentAsyncContext_FunctionContext]
	alreadyInstrumented  *jutil.LinkedHashSet[*rhino.Node]
	hasSuper             *jutil.LinkedHashSet[*rhino.Node]
}

// NOTE: we prefix all internal symbols with a characteristic "ᵃᶜ" (U+1D43, U+1D9C) to
// significantly reduce the chance of conflicts with existing names in the code.  This isn't
// perfect, but it results in much more readable transpiled code compared to always generating a
// longer name, keeps tests simple and deterministic compared to using an ID generator, and is
// more performant than scanning for conflicting symbols to only disambiguate when necessary.
const (
	// Java: InstrumentAsyncContext.FACTORY
	InstrumentAsyncContext_FACTORY = "ᵃᶜfactory"
	// Java: InstrumentAsyncContext.SUSPEND
	InstrumentAsyncContext_SUSPEND = "ᵃᶜsuspend"
	// Java: InstrumentAsyncContext.RESUME
	InstrumentAsyncContext_RESUME = "ᵃᶜresume"
)

// Java: InstrumentAsyncContext(AbstractCompiler, boolean)
func NewInstrumentAsyncContext(compiler AbstractCompiler, shouldInstrumentAwait bool) *InstrumentAsyncContext {
	i := &InstrumentAsyncContext{}
	i.tryFunctionStack = nil
	i.needsInstrumentation = jutil.NewLinkedHashMap[*rhino.Node, *InstrumentAsyncContext_FunctionContext]()
	i.alreadyInstrumented = jutil.NewLinkedHashSet[*rhino.Node]()
	i.hasSuper = jutil.NewLinkedHashSet[*rhino.Node]()
	i.compiler = compiler
	i.astFactory = compiler.CreateAstFactory()
	i.uniqueIdSupplier = compiler.GetUniqueIdSupplier()
	i.shouldInstrumentAwait = shouldInstrumentAwait
	i.start = compiler.GetRuntimeJsLibManager().GetJsLibField("$jscomp.asyncContextStart")
	return i
}

// Java: InstrumentAsyncContext.process(Node, Node)
func (i *InstrumentAsyncContext) Process(externs *rhino.Node, root *rhino.Node) {
	jutil.CheckStateMsg1(i.compiler.GetLifeCycleStage().IsNormalized(), "%s", i.compiler.GetLifeCycleStage())
	// Scan through code and find uses of AsyncContext. Look specifically for Variable.run.
	NodeTraversal_Traverse(i.compiler, root, i)
}

// tryFunctionStackPeek is Deque.peek(): the top, or null when empty.
//
// Go-only: language-adapter — the ArrayDeque operation on the slice stack.
func (i *InstrumentAsyncContext) tryFunctionStackPeek() *rhino.Node {
	if len(i.tryFunctionStack) == 0 {
		return nil
	}
	return i.tryFunctionStack[len(i.tryFunctionStack)-1]
}

// Java: InstrumentAsyncContext.shouldTraverse(NodeTraversal, Node, Node)
func (i *InstrumentAsyncContext) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	// The top of tryFunctionStack should always point to the innermost try or function block.
	// This is a minor optimization to avoid recursing up the ancestor chain to determine if a
	// reentrance is within a try block, in which case a "resume" call needs to be added to the
	// subsequent catch or finally.  We add functions to the stack as well because functions have
	// their own separate scope, so any try block surrounding the function body should _not_ be
	// instrumented.
	if n.IsTry() || n.IsFunction() {
		i.tryFunctionStack = append(i.tryFunctionStack, n)
	}
	return true
}

// Java: InstrumentAsyncContext.isReentrance(Node)
func (i *InstrumentAsyncContext) IsReentrance(n *rhino.Node) bool {
	if n.IsYield() {
		return true
	} else if !i.shouldInstrumentAwait {
		return false
	}
	return n.IsForAwaitOf() || n.IsAwait()
}

// Java: InstrumentAsyncContext.visit(NodeTraversal, Node, Node)
func (i *InstrumentAsyncContext) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if n.IsName() && strings.HasPrefix(n.GetString(), InstrumentAsyncContext_FACTORY) {
		i.alreadyInstrumented.Add(t.GetEnclosingFunction())
	} else if n.IsSuper() {
		i.hasSuper.Add(t.GetEnclosingFunction())
	} else if i.IsReentrance(n) {
		i.InstrumentReentrance(t, n, parent)
	} else if n == i.tryFunctionStackPeek() {
		i.tryFunctionStack = i.tryFunctionStack[:len(i.tryFunctionStack)-1]
		if i.alreadyInstrumented.Contains(n) {
			return
		}
		var context *InstrumentAsyncContext_FunctionContext
		if n.IsGeneratorFunction() {
			context = i.needsInstrumentation.ComputeIfAbsent(
				n, func(unused *rhino.Node) *InstrumentAsyncContext_FunctionContext {
					return i.CreateNewFunctionContext(t.GetInput())
				})
		} else {
			context = i.needsInstrumentation.Get(n)
		}
		if context == nil {
			return
		}
		if n.IsTry() {
			i.InstrumentTry(n, context)
		} else if n.IsGeneratorFunction() {
			i.InstrumentGeneratorFunction(t, n, context)
		} else {
			jutil.CheckState(n.IsAsyncFunction())
			i.InstrumentAsyncFunction(n, t.GetInput(), context)
		}
	}
}

// InstrumentReentrance: when we see a reentrance node, we need to do several
// things:
//  1. replace `await EXPR` with `ᵃᶜresume(await ᵃᶜsuspend(EXPR))`
//  2. ensure the enclosing function has been instrumented with
//     `$jscomp.asyncContextStart` and a `try`-`finally`
//  3. look for an immediately-enclosing `try` block and instrument it with
//     `ᵃᶜresume`
//
// Java: InstrumentAsyncContext.instrumentReentrance(NodeTraversal, Node, Node)
func (i *InstrumentAsyncContext) InstrumentReentrance(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	enclosingTry := i.tryFunctionStackPeek()
	if enclosingTry != nil && !enclosingTry.IsTry() {
		enclosingTry = nil // not a try
	}
	enclosingFunction := t.GetEnclosingFunction()
	if enclosingFunction == nil {
		// NOTE: This is a top-level await, which is currently not supported.  If it becomes supported
		// then we should consider instrumenting the module body as if it were a function.
		// See https://github.com/google/closure-compiler/issues/3835.
		return
	} else if i.alreadyInstrumented.Contains(enclosingFunction) {
		// This function is already instrumented, so don't do any further instrumentation.
		return
	} else if parent.IsFunction() {
		// `async (...) => await ...` - no reentrance, so no instrumentation needed.
		return
	} else if parent.IsReturn() && enclosingTry == nil {
		// `return await ...` - no reentrance unless we're inside a `try`.
		return
	}

	// At this point, instrumentation is required, both for the immediate node, as well as for any
	// enclosing function or try block.
	functionContext :=
		i.needsInstrumentation.ComputeIfAbsent(
			enclosingFunction, func(unused *rhino.Node) *InstrumentAsyncContext_FunctionContext {
				return i.CreateNewFunctionContext(t.GetInput())
			})
	if enclosingTry != nil {
		i.needsInstrumentation.Put(enclosingTry, functionContext)
	}

	if n.IsForAwaitOf() {
		i.InstrumentForAwaitOf(t, n, parent, functionContext)
		return
	}

	// Wrap yielded/awaited expression with a ᵃᶜsuspend(...) call
	placeholder := rhino.IR_Empty()
	arg := n.GetFirstChild()
	if arg == nil {
		n.AddChildToBack(functionContext.Suspend().SrcrefTreeIfMissing(n))
	} else {
		arg.ReplaceWith(placeholder)
		placeholder.ReplaceWith(functionContext.SuspendWithInner(arg).SrcrefTreeIfMissing(arg))
	}

	// Wrap the entire yield/await with a ᵃᶜresume(...) call
	n.ReplaceWith(placeholder)
	placeholder.ReplaceWith(functionContext.ResumeWithInner(n).SrcrefTreeIfMissing(n))
}

// Java: InstrumentAsyncContext.instrumentForAwaitOf(NodeTraversal, Node, Node, FunctionContext)
func (i *InstrumentAsyncContext) InstrumentForAwaitOf(t *NodeTraversal, n *rhino.Node, parent *rhino.Node, functionContext *InstrumentAsyncContext_FunctionContext) {
	// for await (const x of expr()) {
	//   use(x);
	// }
	//     becomes
	// for await (const x of suspend(expr())) {
	//   resume();
	//   try {
	//     use(x);
	//   } finally {
	//     suspend();
	//   }
	// }
	// resume();

	// First wrap the iterator argument
	placeholder := rhino.IR_Empty()
	arg := n.GetSecondChild()
	arg.ReplaceWith(placeholder)
	placeholder.ReplaceWith(functionContext.SuspendWithInner(arg).SrcrefTreeIfMissing(arg))

	body := n.GetLastChild()
	jutil.CheckState(body.IsBlock()) // NOTE: IRFactory normalizes non-block `for` bodies
	body.ReplaceWith(placeholder)
	placeholder.ReplaceWith(
		rhino.IR_BlockOfNodes(
			rhino.IR_ExprResult(functionContext.Resume().SrcrefTreeIfMissing(arg)),
			rhino.IR_TryFinally(
				body, //
				rhino.IR_BlockOfNode(rhino.IR_ExprResult(functionContext.Suspend().SrcrefTreeIfMissing(arg))))))

	// Add resume call after for-await-of
	if !parent.IsBlock() {
		// NOTE: if the for-await-of is _not_ the child of a block, we could be in trouble.
		block := rhino.IR_Block()
		n.ReplaceWith(block)
		block.AddChildToFront(n)
		n = block
	}
	rhino.IR_ExprResult(functionContext.Resume().SrcrefTreeIfMissing(arg)).InsertAfter(n)
}

// Java: InstrumentAsyncContext.instrumentTry(Node, FunctionContext)
func (i *InstrumentAsyncContext) InstrumentTry(n *rhino.Node, functionContext *InstrumentAsyncContext_FunctionContext) {
	// Find either a catch or a finally block (if there is no catch).
	block := n.GetSecondChild()
	if block.HasChildren() {
		// There's a catch block: descend into the actual block
		block = block.GetFirstChild().GetSecondChild()
	} else {
		// Look for a finally block instead
		block = block.GetNext()
	}
	jutil.CheckState(block.IsBlock())
	// Prepend an empty resume call: ᵃᶜresume()
	block.AddChildToFront(rhino.IR_ExprResult(functionContext.Resume().SrcrefTreeIfMissing(block)))
}

// Java: InstrumentAsyncContext.instrumentGeneratorFunction(NodeTraversal, Node, FunctionContext)
func (i *InstrumentAsyncContext) InstrumentGeneratorFunction(t *NodeTraversal, f *rhino.Node, functionContext *InstrumentAsyncContext_FunctionContext) {
	// If a function is completely empty then there's no need to instrument it
	// because there's nothing within that can observe the state of any variables.
	if !f.GetLastChild().HasChildren() {
		return
	}
	// NOTE: if `f` is a generator _method_ then we need to be more careful.
	// Because generators are not allowed to be arrow functions, we need to
	// use the ThisAndArgumentsReferenceUpdater to extract any references to
	// those symbols within the vanilla IIFE generator.  But this doesn't help
	// with any references to `super`.  Instead, for methods we have a clear
	// method body and can move the generator body into a newly synthesized
	// method that we can directly.  If "arguments" is used, then we need to
	// call it with .apply and forward the arguments.  This is dangerous for
	// certain methods that might be collapsed (e.g. static methods), so we
	// only do when necessary (i.e. `super` is used), which corresponds to a
	// known-safe case (since `super` can't be collapsed).
	if i.hasSuper.Contains(f) {
		i.InstrumentGeneratorMethod(t, f, functionContext)
		return
	}

	// Turn the generator into an ordinary function with an IIFE generator
	f.SetIsGeneratorFunction(false)
	generatorBody := f.GetLastChild().Detach()
	inner := i.astFactory.CreateEmptyGeneratorFunction(AstFactory_Type_(f))

	if f.IsAsyncFunction() {
		f.SetIsAsyncFunction(false)
		inner.SetIsAsyncFunction(true)
	}
	inner.GetLastChild().ReplaceWith(generatorBody)

	generatorBody = i.AddFinallySuspend(generatorBody, functionContext)
	generatorBody.AddChildToFront(rhino.IR_ExprResult(functionContext.Resume()))
	newOuterBody :=
		rhino.IR_BlockOfNodes(
			i.CreateFactorySuspendCall(functionContext.FactoryName()),
			rhino.IR_ReturnNodeWithExpr(i.astFactory.CreateCallWithUnknownType(inner)))
	i.AddSuspendResumeVarsAfter(newOuterBody.GetFirstChild(), functionContext)
	f.AddChildToBack(newOuterBody.SrcrefTreeIfMissing(generatorBody))

	// NOTE: if the IIFE generator refers to `this` or `arguments` then we need to extract these
	// references into local variables.
	thisContext :=
		NewThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext(
			newOuterBody, false, i.uniqueIdSupplier.GetUniqueId(t.GetInput()))
	updater :=
		NewThisAndArgumentsReferenceUpdater(i.compiler, thisContext, i.astFactory)
	NodeTraversal_Traverse(i.compiler, generatorBody, updater)
	// Check the context to see if we found this or arguments.
	thisContext.AddVarDeclarations(i.compiler, i.astFactory, t.GetCurrentScript())

	i.compiler.ReportChangeToChangeScope(f)
	i.compiler.ReportChangeToChangeScope(inner)
	// Instrumentation may move functions into a new nested scope.
	NodeUtil_AddFeatureToScript(
		t.GetCurrentScript(), parser.FeatureSet_Feature_BLOCK_SCOPED_FUNCTION_DECLARATION, i.compiler)
}

// Java: InstrumentAsyncContext.instrumentGeneratorMethod(NodeTraversal, Node, FunctionContext)
func (i *InstrumentAsyncContext) InstrumentGeneratorMethod(t *NodeTraversal, f *rhino.Node, functionContext *InstrumentAsyncContext_FunctionContext) {
	// This requires special handling compared to ordinary generator functions because
	// it may reference `super`, which would break in a vanilla IIFE.  Instead, copy the
	// body to a new method with the same parameters, plus one or two more for the
	// context and possibly `arguments`.

	// Specifically:
	//     *m(a, {b: c=1}, d=2, ...e) { ... }
	// would become
	//     *m$jscomp$123(ᵃᶜfactory, a, {b: c=1}, d=2, ...e) {
	//       var ᵃᶜsuspend = ᵃᶜfactory();
	//       var ᵃᶜresume = ᵃᶜfactory(1);
	//       try { ... } finally { ᵃᶜsuspend(); }
	//     }
	//     m(a, $jscomp$0, d, ...e) {
	//       const ᵃᶜfactory = $jscomp.asyncContextStart();
	//       return this.m$jscomp$123(ᵃᶜfactory, a, $jscomp$0, d, ...e);
	//     }
	// If the function refers to `arguments`, then an additional
	// `ᵃᶜarguments` parameter is added to the synthesized method after
	// the context argument.

	innerMethodNamePrefix := ""
	if f.GetParent().IsMemberFunctionDef() {
		innerMethodNamePrefix = f.GetParent().GetString()
	}
	innerMethodName :=
		innerMethodNamePrefix +
			"$jscomp$" +
			i.uniqueIdSupplier.GetUniqueId(t.GetInput())
	outerParams := f.GetSecondChild()
	innerParams := outerParams.CloneTree()

	// Add the ᵃᶜfactory parameter to the front of the inner parameter list
	innerFactoryName := functionContext.FactoryName()
	outerFactoryName :=
		i.astFactory.CreateNameWithUnknownType(InstrumentAsyncContext_FACTORY + i.uniqueIdSupplier.GetUniqueId(t.GetInput()))
	innerParams.AddChildToFront(innerFactoryName)
	unknownType := AstFactory_Type_(innerFactoryName)

	// Add a finally-suspend wrapper and an initial resume call around the original generator body
	// (which has already had its yields/awaits instrumented).
	generatorBody := i.AddFinallySuspend(f.GetLastChild(), functionContext)
	generatorBody.AddChildToFront(rhino.IR_ExprResult(functionContext.Resume()))

	// Look for references to `arguments`.  If found, add an extra parameter to the inner method.
	// We need to save a reference to the parameter in order to replace it in the function call
	// arguments.
	var argumentsParam *rhino.Node = nil
	var argumentsReference *rhino.Node = nil
	renamer := NewInstrumentAsyncContext_ArgumentsRenamer(i)
	NodeTraversal_Traverse(i.compiler, generatorBody, renamer)
	if renamer.argumentsName != nil {
		argumentsReference = i.astFactory.CreateArgumentsReference()
		argumentsParam =
			i.astFactory.CreateName(*renamer.argumentsName, AstFactory_Type_(argumentsReference))
		argumentsParam.InsertAfter(innerFactoryName)
	}

	// Detach the generator body and replace it with a new empty block.  We'll fill the empty block
	// later, and will insert the instrumented generator body into the new inner method.
	outerBody := rhino.IR_Block()
	generatorBody.ReplaceWith(outerBody)

	// Build the call to the inner method.  Iterate over innerParams to build up the arguments list,
	// replacing any non-simple arguments with simple names.
	innerCall :=
		i.astFactory.CreateCallWithUnknownType(
			i.astFactory.CreateGetProp(
				i.astFactory.CreateThis(unknownType), innerMethodName, AstFactory_Type_(f)))
	// Copy the context parameter.
	innerParam := innerParams.GetFirstChild()
	innerCall.AddChildToBack(outerFactoryName.CloneTree())
	innerParam = innerParam.GetNext()
	// Copy the arguments parameter if it's there.
	if innerParam != nil && innerParam == argumentsParam {
		innerCall.AddChildToBack(argumentsReference)
	}
	// Copy the remaining parameters.
	for outerParam := outerParams.GetFirstChild(); outerParam != nil; outerParam = outerParam.GetNext() {
		argParamPair := i.SimplifyParameter(outerParam.CloneTree(), t)
		// Go-only: language-adapter — Java's `outerParam.replaceWith(outerParam = p)` calls
		// replaceWith on the old outerParam, then continues the loop from p.
		oldOuterParam := outerParam
		outerParam = argParamPair.parameter
		oldOuterParam.ReplaceWith(outerParam)
		innerCall.AddChildToBack(argParamPair.argument)
	}

	// Make a new inner method with the instrumented generator body, and insert it after the
	// original method.
	newInnerMethod :=
		i.astFactory.CreateFunction("", innerParams, generatorBody, AstFactory_Type_(f))
	rhino.IR_MemberFunctionDef(innerMethodName, newInnerMethod).InsertAfter(f.GetParent())
	newInnerMethod.GetFirstChild().MakeNonIndexable()

	// Move the async/generator flags from the original method to the new one.
	newInnerMethod.SetIsAsyncFunction(f.IsAsyncFunction())
	newInnerMethod.SetIsGeneratorFunction(true)
	f.SetIsGeneratorFunction(false)
	f.SetIsAsyncFunction(false)

	// Copy static and/or @nocollapse from the parent (member_function_def) node.
	fp := f.GetParent()
	np := newInnerMethod.GetParent()
	np.SetStaticMember(fp.IsStaticMember())
	jsdoc := fp.GetJSDocInfo()
	if jsdoc != nil {
		builder := rhino.JSDocInfo_Builder_()
		if jsdoc.IsNoCollapse() {
			builder.RecordNoCollapse()
		}
		if jsdoc.IsNoInline() {
			builder.RecordNoInline()
		}
		np.SetJSDocInfo(builder.Build())
	}

	// Populate the outer (non-generator) method body with an start() and call to the inner method.
	outerBody.AddChildToBack(i.CreateFactorySuspendCall(outerFactoryName))
	i.AddSuspendResumeVarsAtStartOf(generatorBody, innerFactoryName, functionContext)
	outerBody.AddChildToBack(rhino.IR_ReturnNodeWithExpr(innerCall))

	// Various cleanups
	f.SrcrefTreeIfMissing(f)
	newInnerMethod.GetParent().SrcrefTreeIfMissing(f)
	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_MEMBER_DECLARATIONS, i.compiler)
	i.compiler.ReportChangeToChangeScope(f)
	i.compiler.ReportChangeToChangeScope(newInnerMethod)
	i.compiler.ReportChangeToEnclosingScope(f.GetParent())
	// Instrumentation may move functions into a new nested scope.
	NodeUtil_AddFeatureToScript(
		t.GetCurrentScript(), parser.FeatureSet_Feature_BLOCK_SCOPED_FUNCTION_DECLARATION, i.compiler)
}

// Java: InstrumentAsyncContext.simplifyParameter(Node, NodeTraversal)
func (i *InstrumentAsyncContext) SimplifyParameter(param *rhino.Node, t *NodeTraversal) *InstrumentAsyncContext_ArgParamPair {
	if param.IsRest() {
		// Replace REST with SPREAD.  For arguments, we know it has to be ITER_REST and ITER_SPREAD.
		NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_SPREAD_EXPRESSIONS, i.compiler)
		paramElement := param.RemoveFirstChild()
		simplified := i.SimplifyParameter(paramElement, t)
		return NewInstrumentAsyncContext_ArgParamPair(
			rhino.IR_IterSpread(simplified.argument).SrcrefTreeIfMissing(param),
			rhino.IR_IterRest(simplified.parameter).SrcrefTreeIfMissing(param))
	} else if param.IsDefaultValue() {
		// Ignore any initializers; replace with just the name.
		return i.SimplifyParameter(param.RemoveFirstChild(), t)
	} else if param.IsArrayPattern() || param.IsObjectPattern() {
		// Replace destructuring patterns with a synthesized name.
		newParam :=
			i.astFactory.
				CreateName(
					"$jscomp$param$"+i.uniqueIdSupplier.GetUniqueId(t.GetInput()),
					AstFactory_Type_(param)).
				SrcrefTreeIfMissing(param)
		return NewInstrumentAsyncContext_ArgParamPair(newParam.CloneTree(), newParam)
	} else if param.IsName() {
		// Replace destructuring patterns with a synthesized name.
		newParam :=
			i.astFactory.
				CreateName(
					param.GetString()+"$jscomp$param$"+i.uniqueIdSupplier.GetUniqueId(t.GetInput()),
					AstFactory_Type_(param)).
				SrcrefTreeIfMissing(param)
		return NewInstrumentAsyncContext_ArgParamPair(newParam.CloneTree(), newParam)
	}
	panic(jutil.NewIllegalStateExceptionMsg("Unexpected parameter: " + param.String()))
}

// InstrumentAsyncContext_ArgParamPair is InstrumentAsyncContext.ArgParamPair.
type InstrumentAsyncContext_ArgParamPair struct {
	argument  *rhino.Node
	parameter *rhino.Node
}

// Java: ArgParamPair(Node, Node)
func NewInstrumentAsyncContext_ArgParamPair(argument *rhino.Node, parameter *rhino.Node) *InstrumentAsyncContext_ArgParamPair {
	p := &InstrumentAsyncContext_ArgParamPair{}
	p.argument = argument
	p.parameter = parameter
	return p
}

// InstrumentAsyncContext_ArgumentsRenamer renames all instances of
// `arguments` with a new unique name.
type InstrumentAsyncContext_ArgumentsRenamer struct {
	outer         *InstrumentAsyncContext
	argumentsName *string
}

// Java: ArgumentsRenamer()
func NewInstrumentAsyncContext_ArgumentsRenamer(outer *InstrumentAsyncContext) *InstrumentAsyncContext_ArgumentsRenamer {
	return &InstrumentAsyncContext_ArgumentsRenamer{outer: outer, argumentsName: nil}
}

// Java: ArgumentsRenamer.shouldTraverse(NodeTraversal, Node, Node)
func (r *InstrumentAsyncContext_ArgumentsRenamer) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	return !n.IsFunction() || n.IsArrowFunction()
}

// Java: ArgumentsRenamer.visit(NodeTraversal, Node, Node)
func (r *InstrumentAsyncContext_ArgumentsRenamer) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if n.IsName() && n.GetString() == "arguments" {
		if r.argumentsName == nil {
			name := "$jscomp$arguments$" + r.outer.uniqueIdSupplier.GetUniqueId(t.GetInput())
			r.argumentsName = &name
		}
		n.SetString(*r.argumentsName)
		if r.outer.compiler.GetOptions().PreservesDetailedSourceInfo() {
			originalName := "arguments"
			n.SetOriginalName(&originalName)
		}
		n.MakeNonIndexable()
		n.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
	}
}

// InstrumentAsyncFunction instruments the given async function node by
// wrapping the body in a try-finally and adding an "start" call to the front.
//
// Java: InstrumentAsyncContext.instrumentAsyncFunction(Node, CompilerInput, FunctionContext)
func (i *InstrumentAsyncContext) InstrumentAsyncFunction(f *rhino.Node, input *CompilerInput, functionContext *InstrumentAsyncContext_FunctionContext) {
	// Add a variable to the front of the body
	block := f.GetLastChild()
	newBlock := i.AddFinallySuspend(block, functionContext)
	newBlock.AddChildToFront(
		i.CreateFactoryCall(functionContext.factoryName).SrcrefTreeIfMissing(newBlock))
	i.AddSuspendResumeVarsAfter(newBlock.GetFirstChild(), functionContext)
	i.compiler.ReportChangeToChangeScope(f)

	// TODO(sdh): Need to instrument async functions for unhandled rejections.
	// Ideally this would happen after devirtualization and then we can extract a
	// helper function to call rather than making a new closure.  Would also be nice
	// to not even bother with it if we're not worrying about that case.
}

// AddFinallySuspend replaces a function body with a new body containing only
// try (original body) finally (suspend).
//
// Java: InstrumentAsyncContext.addFinallySuspend(Node, FunctionContext)
func (i *InstrumentAsyncContext) AddFinallySuspend(block *rhino.Node, functionContext *InstrumentAsyncContext_FunctionContext) *rhino.Node {
	placeholder := rhino.IR_Empty()
	block.ReplaceWith(placeholder)
	finallyBlock := rhino.IR_BlockOfNode(rhino.IR_ExprResult(functionContext.Suspend()))
	newBlock := rhino.IR_BlockOfNode(rhino.IR_TryFinally(block, finallyBlock)).SrcrefTreeIfMissing(block)
	placeholder.ReplaceWith(newBlock)
	return newBlock
}

// Java: InstrumentAsyncContext.addSuspendResumeVarsAfter(Node, FunctionContext)
func (i *InstrumentAsyncContext) AddSuspendResumeVarsAfter(node *rhino.Node, functionContext *InstrumentAsyncContext_FunctionContext) {
	// var ᵃᶜsuspend = ᵃᶜfactory();
	// var ᵃᶜresume = ᵃᶜfactory(1);
	resumeCall :=
		i.astFactory.CreateCallWithUnknownType(
			functionContext.factoryName.CloneNode(), i.astFactory.CreateNumber(1))
	rhino.IR_VarWithValue(functionContext.resumeName, resumeCall).SrcrefTreeIfMissing(node).InsertAfter(node)
	suspendCall :=
		i.astFactory.CreateCallWithUnknownType(functionContext.factoryName.CloneNode())
	rhino.IR_VarWithValue(functionContext.suspendName, suspendCall).SrcrefTreeIfMissing(node).InsertAfter(node)
}

// Java: InstrumentAsyncContext.addSuspendResumeVarsAtStartOf(Node, Node, FunctionContext)
func (i *InstrumentAsyncContext) AddSuspendResumeVarsAtStartOf(
	node *rhino.Node, factoryName *rhino.Node, functionContext *InstrumentAsyncContext_FunctionContext) {
	// var ᵃᶜsuspend = ᵃᶜfactory();
	// var ᵃᶜresume = ᵃᶜfactory(1);
	resumeCall :=
		i.astFactory.CreateCallWithUnknownType(factoryName.CloneNode(), i.astFactory.CreateNumber(1))
	node.AddChildToFront(rhino.IR_VarWithValue(functionContext.resumeName, resumeCall).SrcrefTreeIfMissing(node))
	suspendCall := i.astFactory.CreateCallWithUnknownType(factoryName.CloneNode())
	node.AddChildToFront(
		rhino.IR_VarWithValue(functionContext.suspendName, suspendCall).SrcrefTreeIfMissing(node))
}

// Java: InstrumentAsyncContext.createVar(Node, Node)
func InstrumentAsyncContext_CreateVar(name *rhino.Node, value *rhino.Node) *rhino.Node {
	return rhino.IR_VarWithValue(name, value)
}

// CreateFactoryCall creates a statement
// `const ᵃᶜfactory = $jscomp$asyncContextStart();`.
//
// Java: InstrumentAsyncContext.createFactoryCall(Node, Node...)
func (i *InstrumentAsyncContext) CreateFactoryCall(factoryName *rhino.Node, arg ...*rhino.Node) *rhino.Node {
	call :=
		i.astFactory.CreateCallWithUnknownType(i.astFactory.CreateQNameWithUnknownTypeOfField(i.start), arg...)
	return InstrumentAsyncContext_CreateVar(factoryName, call)
}

// CreateFactorySuspendCall creates a statement
// `const ᵃᶜfactory = $jscomp$asyncContextStart(1);`.
//
// Java: InstrumentAsyncContext.createFactorySuspendCall(Node)
func (i *InstrumentAsyncContext) CreateFactorySuspendCall(factoryName *rhino.Node) *rhino.Node {
	return i.CreateFactoryCall(factoryName, i.astFactory.CreateNumber(1))
}

// InstrumentAsyncContext_FunctionContext is the record
// InstrumentAsyncContext.FunctionContext.
type InstrumentAsyncContext_FunctionContext struct {
	uniqueSuffix string
	suspendName  *rhino.Node
	resumeName   *rhino.Node
	factoryName  *rhino.Node
	astFactory   *AstFactory
}

// Java: FunctionContext(String, Node, Node, Node, AstFactory)
func NewInstrumentAsyncContext_FunctionContext(
	uniqueSuffix string, suspendName *rhino.Node, resumeName *rhino.Node, factoryName *rhino.Node, astFactory *AstFactory) *InstrumentAsyncContext_FunctionContext {
	return &InstrumentAsyncContext_FunctionContext{
		uniqueSuffix: uniqueSuffix,
		suspendName:  suspendName,
		resumeName:   resumeName,
		factoryName:  factoryName,
		astFactory:   astFactory,
	}
}

// Java: FunctionContext.uniqueSuffix()
func (c *InstrumentAsyncContext_FunctionContext) UniqueSuffix() string { return c.uniqueSuffix }

// Java: FunctionContext.suspendName()
func (c *InstrumentAsyncContext_FunctionContext) SuspendName() *rhino.Node { return c.suspendName }

// Java: FunctionContext.resumeName()
func (c *InstrumentAsyncContext_FunctionContext) ResumeName() *rhino.Node { return c.resumeName }

// Java: FunctionContext.factoryName()
func (c *InstrumentAsyncContext_FunctionContext) FactoryName() *rhino.Node { return c.factoryName }

// Java: FunctionContext.astFactory()
func (c *InstrumentAsyncContext_FunctionContext) AstFactory() *AstFactory { return c.astFactory }

// Suspend creates an empty suspend call: `ᵃᶜsuspend()`.
//
// Java: FunctionContext.suspend()
func (c *InstrumentAsyncContext_FunctionContext) Suspend() *rhino.Node {
	return c.astFactory.CreateCall(c.suspendName.CloneNode(), AstFactory_Type_(c.suspendName))
}

// SuspendWithInner wraps the given node in a suspend call: `ᵃᶜsuspend(NODE)`.
//
// Java: FunctionContext.suspend(Node)
func (c *InstrumentAsyncContext_FunctionContext) SuspendWithInner(inner *rhino.Node) *rhino.Node {
	return c.astFactory.CreateCall(c.suspendName.CloneNode(), AstFactory_Type_(inner), inner)
}

// Resume creates an empty resume call: `ᵃᶜresume()`.
//
// Java: FunctionContext.resume()
func (c *InstrumentAsyncContext_FunctionContext) Resume() *rhino.Node {
	return c.astFactory.CreateCall(c.resumeName.CloneNode(), AstFactory_Type_(c.resumeName))
}

// ResumeWithInner wraps the given node in a resume call: `ᵃᶜresume(NODE)`.
//
// Java: FunctionContext.resume(Node)
func (c *InstrumentAsyncContext_FunctionContext) ResumeWithInner(inner *rhino.Node) *rhino.Node {
	return c.astFactory.CreateCall(c.resumeName.CloneNode(), AstFactory_Type_(inner), inner)
}

// Java: InstrumentAsyncContext.createNewFunctionContext(CompilerInput)
func (i *InstrumentAsyncContext) CreateNewFunctionContext(input *CompilerInput) *InstrumentAsyncContext_FunctionContext {
	uniqueSuffix := i.uniqueIdSupplier.GetUniqueId(input)
	// Creates a unique name, prefixed with `ᵃᶜsuspend`
	suspendNode := i.astFactory.CreateNameWithUnknownType(InstrumentAsyncContext_SUSPEND + uniqueSuffix)

	// Creates a unique name, prefixed with `ᵃᶜresume`
	resumeNode := i.astFactory.CreateNameWithUnknownType(InstrumentAsyncContext_RESUME + uniqueSuffix)

	// Creates a unique name, prefixed with `ᵃᶜfactory`
	factoryName := i.astFactory.CreateNameWithUnknownType(InstrumentAsyncContext_FACTORY + uniqueSuffix)

	return NewInstrumentAsyncContext_FunctionContext(uniqueSuffix, suspendNode, resumeNode, factoryName, i.astFactory)
}
