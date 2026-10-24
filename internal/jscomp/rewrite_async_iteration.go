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
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// RewriteAsyncIteration converts async generator functions into a function
// returning a new $jscomp.AsyncGenWrapper around the original block and
// awaits/yields converted to yields of ActionRecords.
//
//	async function* foo() {
//	  let res = await myPromise;
//	  yield res + 1;
//	}
//
// becomes (prefixes trimmed for clarity)
//
//	function foo() {
//	  return new $jscomp.AsyncGeneratorWrapper((function*(){
//	    let res = yield new $ActionRecord($ActionEnum.AWAIT_VALUE, myPromise);
//	    yield new $ActionRecord($ActionEnum.YIELD_VALUE, res + 1);
//	  })());
//	}
type RewriteAsyncIteration struct {
	asyncGeneratorWrapper RuntimeJsLibManager_JsLibField
	actionRecordName      RuntimeJsLibManager_JsLibField
	actionEnumAwait       RuntimeJsLibManager_JsLibField
	actionEnumYield       RuntimeJsLibManager_JsLibField
	actionEnumYieldStar   RuntimeJsLibManager_JsLibField

	nextForAwaitId int

	compiler AbstractCompiler

	// Go-only: language-adapter — the ArrayDeque used as a stack is a slice
	// whose last element is its first (push/pop/getFirst/remove work on it).
	contextStack []*RewriteAsyncIteration_LexicalContext
	astFactory   *AstFactory
	namespace    rhino.StaticScope
}

// Java: RewriteAsyncIteration.transpiledFeatures
var RewriteAsyncIteration_transpiledFeatures = parser.FeatureSet_BARE_MINIMUM.WithFeatures(parser.FeatureSet_Feature_ASYNC_GENERATORS, parser.FeatureSet_Feature_FOR_AWAIT_OF)

// Java: RewriteAsyncIteration.CANNOT_CONVERT_ASYNCGEN
var RewriteAsyncIteration_CANNOT_CONVERT_ASYNCGEN = DiagnosticType_Error("JSC_CANNOT_CONVERT_ASYNCGEN", "Cannot convert async generator. {0}")

// Variables with these names get created when rewriting for-await-of loops
const (
	// Java: RewriteAsyncIteration.FOR_AWAIT_ITERATOR_TEMP_NAME
	RewriteAsyncIteration_FOR_AWAIT_ITERATOR_TEMP_NAME = "$jscomp$forAwait$tempIterator"
	// Java: RewriteAsyncIteration.FOR_AWAIT_RESULT_TEMP_NAME
	RewriteAsyncIteration_FOR_AWAIT_RESULT_TEMP_NAME = "$jscomp$forAwait$tempResult"
	// Java: RewriteAsyncIteration.FOR_AWAIT_ERROR_RESULT_TEMP_NAME
	RewriteAsyncIteration_FOR_AWAIT_ERROR_RESULT_TEMP_NAME = "$jscomp$forAwait$errResult"
	// Java: RewriteAsyncIteration.FOR_AWAIT_CATCH_PARAM_TEMP_NAME
	RewriteAsyncIteration_FOR_AWAIT_CATCH_PARAM_TEMP_NAME = "$jscomp$forAwait$catchErrParam"
	// Java: RewriteAsyncIteration.FOR_AWAIT_RETURN_FN_TEMP_NAME
	RewriteAsyncIteration_FOR_AWAIT_RETURN_FN_TEMP_NAME = "$jscomp$forAwait$retFn"
)

const (
	// Java: RewriteAsyncIteration.THIS_VAR_NAME
	RewriteAsyncIteration_THIS_VAR_NAME = "$jscomp$asyncIter$this$"
	// Java: RewriteAsyncIteration.ARGUMENTS_VAR_NAME
	RewriteAsyncIteration_ARGUMENTS_VAR_NAME = "$jscomp$asyncIter$arguments"
	// Java: RewriteAsyncIteration.SUPER_PROP_GETTER_PREFIX
	RewriteAsyncIteration_SUPER_PROP_GETTER_PREFIX = "$jscomp$asyncIter$super$get$"
)

// RewriteAsyncIteration_LexicalContext tracks a function and its context of
// this/arguments/super, if such a context exists.
type RewriteAsyncIteration_LexicalContext struct {
	// Node that creates the context
	contextRoot rhino.Node
	// The current function, or null if root scope where we are not in a function.
	function rhino.Node
	// The context of the most recent definition of this/super/arguments
	thisSuperArgsContext *RewriteAsyncIteration_ThisSuperArgsContext
}

// newRewriteAsyncIteration_LexicalContext represents the global/root scope.
// Should only exist on the bottom of the contextStack.
//
// Java: LexicalContext(Node)
func newRewriteAsyncIteration_LexicalContext(contextRoot rhino.Node) *RewriteAsyncIteration_LexicalContext {
	c := &RewriteAsyncIteration_LexicalContext{}
	c.contextRoot = rhino.CheckNotNull(contextRoot)
	c.function = rhino.NilNode
	c.thisSuperArgsContext =
		nil // no need for global context to have a this/super/args context
	return c
}

// newRewriteAsyncIteration_LexicalContextWithParent represents the context of
// a function or its parameter list.
//
// parent is the enclosing context; contextRoot is the FUNCTION or PARAM_LIST
// node; function is the same as contextRoot or the FUNCTION containing the
// PARAM_LIST.
//
// Java: LexicalContext(LexicalContext, Node, Node, AbstractCompiler)
func newRewriteAsyncIteration_LexicalContextWithParent(
	parent *RewriteAsyncIteration_LexicalContext, contextRoot rhino.Node, function rhino.Node, compiler AbstractCompiler) *RewriteAsyncIteration_LexicalContext {
	c := &RewriteAsyncIteration_LexicalContext{}
	jutil.CheckNotNull(parent)
	rhino.CheckNotNull(contextRoot)
	jutil.CheckArgumentMsg1(contextRoot == function || contextRoot.IsParamList(), "%s", contextRoot)
	rhino.CheckNotNull(function)
	jutil.CheckArgumentMsg1(function.IsFunction(), "%s", function)
	c.contextRoot = contextRoot
	c.function = function

	if function.IsArrowFunction() {
		// Use the parent context to inherit this, arguments, and super for an arrow function or its
		// parameter list.
		c.thisSuperArgsContext = parent.thisSuperArgsContext
	} else if contextRoot.IsFunction() {
		// Non-arrow function gets its own context defining `this`, `arguments`, and `super`.
		newUniqueId :=
			compiler.
				GetUniqueIdSupplier().
				GetUniqueId(compiler.GetInput(NodeUtil_GetInputId(contextRoot)))
		c.thisSuperArgsContext = newRewriteAsyncIteration_ThisSuperArgsContext(c, newUniqueId)
	} else {
		// contextRoot is a parameter list.
		// Never alias `this`, `arguments`, or `super` for normal function parameter lists.
		// They are implicitly defined there.
		c.thisSuperArgsContext = nil
	}
	return c
}

// Java: LexicalContext.newGlobalContext(Node)
func RewriteAsyncIteration_LexicalContext_NewGlobalContext(contextRoot rhino.Node) *RewriteAsyncIteration_LexicalContext {
	return newRewriteAsyncIteration_LexicalContext(contextRoot)
}

// Java: LexicalContext.newContextForFunction(LexicalContext, Node, AbstractCompiler)
func RewriteAsyncIteration_LexicalContext_NewContextForFunction(
	parent *RewriteAsyncIteration_LexicalContext, function rhino.Node, compiler AbstractCompiler) *RewriteAsyncIteration_LexicalContext {
	// Functions need their own context because:
	//     - async generator functions must be transpiled
	//     - non-async generator functions must NOT be transpiled
	//     - arrow functions inside of async generator functions need to have
	//       `this`, `arguments`, and `super` references aliased, including in their
	//       parameter lists
	return newRewriteAsyncIteration_LexicalContextWithParent(parent, function, function, compiler)
}

// Java: LexicalContext.newContextForParamList(LexicalContext, Node, AbstractCompiler)
func RewriteAsyncIteration_LexicalContext_NewContextForParamList(
	parent *RewriteAsyncIteration_LexicalContext, paramList rhino.Node, compiler AbstractCompiler) *RewriteAsyncIteration_LexicalContext {
	// Parameter lists need their own context because `this`, `arguments`, and `super` must NOT be
	// aliased for non-arrow function parameter lists, even for async generator functions.
	return newRewriteAsyncIteration_LexicalContextWithParent(parent, paramList, parent.function, compiler)
}

// Java: LexicalContext.getFunctionDeclaringThisArgsSuper()
func (c *RewriteAsyncIteration_LexicalContext) GetFunctionDeclaringThisArgsSuper() rhino.Node {
	return c.thisSuperArgsContext.ctx.function
}

// MustReplaceThisSuperArgs: is it necessary to replace `this`, `super`, and
// `arguments` with aliases in this context?
//
// Java: LexicalContext.mustReplaceThisSuperArgs()
func (c *RewriteAsyncIteration_LexicalContext) MustReplaceThisSuperArgs() bool {
	return c.thisSuperArgsContext != nil &&
		c.GetFunctionDeclaringThisArgsSuper().IsAsyncGeneratorFunction()
}

// RewriteAsyncIteration_ThisSuperArgsContext tracks how this/arguments/super
// were used in the function so declarations of replacement variables can be
// prepended
type RewriteAsyncIteration_ThisSuperArgsContext struct {
	// The LexicalContext representing the function that declared this/super/args
	ctx *RewriteAsyncIteration_LexicalContext

	usedSuperProperties *jutil.LinkedHashSet[rhino.Node]
	thisNodeToAdd       rhino.Node
	usedArguments       bool
	// unique id to append to names in this context. This is used to ensure that names
	// in different contexts don't collide (e.g. 2 functions don't get the same `let
	// $jscomp$async$this` name declared in their bodies)
	uniqueId string
}

// Java: ThisSuperArgsContext(LexicalContext, String)
func newRewriteAsyncIteration_ThisSuperArgsContext(ctx *RewriteAsyncIteration_LexicalContext, uniqueId string) *RewriteAsyncIteration_ThisSuperArgsContext {
	c := &RewriteAsyncIteration_ThisSuperArgsContext{}
	c.usedSuperProperties = jutil.NewLinkedHashSet[rhino.Node]()
	c.thisNodeToAdd = rhino.NilNode
	c.usedArguments = false
	c.ctx = ctx
	c.uniqueId = uniqueId
	return c
}

// Java: RewriteAsyncIteration(AbstractCompiler, AstFactory, StaticScope)
func newRewriteAsyncIteration(
	compiler AbstractCompiler, astFactory *AstFactory, namespace rhino.StaticScope) *RewriteAsyncIteration {
	r := &RewriteAsyncIteration{}
	r.nextForAwaitId = 0
	r.compiler = jutil.CheckNotNull(compiler)
	r.astFactory = jutil.CheckNotNull(astFactory)
	r.namespace = jutil.CheckNotNull(namespace)
	r.contextStack = nil

	runtimeLibManager := compiler.GetRuntimeJsLibManager()
	r.asyncGeneratorWrapper = runtimeLibManager.GetJsLibField("$jscomp.AsyncGeneratorWrapper")
	r.actionRecordName =
		runtimeLibManager.GetJsLibField("$jscomp.AsyncGeneratorWrapper$ActionRecord")
	r.actionEnumAwait =
		runtimeLibManager.GetJsLibField("$jscomp.AsyncGeneratorWrapper$ActionEnum.AWAIT_VALUE")
	r.actionEnumYield =
		runtimeLibManager.GetJsLibField("$jscomp.AsyncGeneratorWrapper$ActionEnum.YIELD_VALUE")
	r.actionEnumYieldStar =
		runtimeLibManager.GetJsLibField("$jscomp.AsyncGeneratorWrapper$ActionEnum.YIELD_STAR")
	return r
}

// Java: RewriteAsyncIteration.create(AbstractCompiler)
func RewriteAsyncIteration_Create(compiler AbstractCompiler) *RewriteAsyncIteration {
	astFactory := compiler.CreateAstFactory()
	namespace := compiler.GetTranspilationNamespace()
	return newRewriteAsyncIteration(compiler, astFactory, namespace)
}

// contextStackGetFirst is ArrayDeque.getFirst(): the top, or
// NoSuchElementException.
//
// Go-only: language-adapter — the ArrayDeque operation on the slice stack.
func (r *RewriteAsyncIteration) contextStackGetFirst() *RewriteAsyncIteration_LexicalContext {
	if len(r.contextStack) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	return r.contextStack[len(r.contextStack)-1]
}

// contextStackPop is ArrayDeque.pop() (and remove()): removes the top, or
// NoSuchElementException.
//
// Go-only: language-adapter — the ArrayDeque operation on the slice stack.
func (r *RewriteAsyncIteration) contextStackPop() *RewriteAsyncIteration_LexicalContext {
	top := r.contextStackGetFirst()
	r.contextStack = r.contextStack[:len(r.contextStack)-1]
	return top
}

// Java: RewriteAsyncIteration.process(Node, Node)
func (r *RewriteAsyncIteration) Process(externs rhino.Node, root rhino.Node) {
	jutil.CheckState(len(r.contextStack) == 0)
	r.contextStack = append(r.contextStack, RewriteAsyncIteration_LexicalContext_NewGlobalContext(root))
	TranspilationPasses_ProcessTranspile(r.compiler, root, RewriteAsyncIteration_transpiledFeatures, r)
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(r.compiler, root, RewriteAsyncIteration_transpiledFeatures)
	jutil.CheckState(r.contextStackGetFirst().function.IsNil())
	r.contextStackPop()
	jutil.CheckState(len(r.contextStack) == 0)
}

// Java: RewriteAsyncIteration.shouldTraverse(NodeTraversal, Node, Node)
func (r *RewriteAsyncIteration) ShouldTraverse(nodeTraversal *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if n.IsFunction() {
		r.contextStack = append(r.contextStack,
			RewriteAsyncIteration_LexicalContext_NewContextForFunction(r.contextStackGetFirst(), n, r.compiler))
	} else if n.IsParamList() {
		r.contextStack = append(r.contextStack,
			RewriteAsyncIteration_LexicalContext_NewContextForParamList(r.contextStackGetFirst(), n, r.compiler))
	}
	return true
}

// Java: RewriteAsyncIteration.visit(NodeTraversal, Node, Node)
func (r *RewriteAsyncIteration) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	ctx := r.contextStackGetFirst()
	switch n.GetToken() {
	case rhino.Token_PARAM_LIST:
		// Async Generators (and popping contexts)
		// Done handling parameter list, so pop its context
		jutil.CheckStateMsg1(n == ctx.contextRoot, "%s", n)
		r.contextStackPop()
	case rhino.Token_FUNCTION:
		jutil.CheckState(n == ctx.contextRoot)
		if n.IsAsyncGeneratorFunction() {
			r.ConvertAsyncGenerator(n)
			r.PrependTempVarDeclarations(ctx, t)
		}
		// Done handling function, so pop its context
		r.contextStackPop()
	case rhino.Token_AWAIT:
		rhino.CheckNotNull(ctx.function)
		if ctx.function.IsAsyncGeneratorFunction() {
			r.ConvertAwaitOfAsyncGenerator(ctx, n)
		}
	case rhino.Token_YIELD:
		// Includes yield*
		rhino.CheckNotNull(ctx.function)
		if ctx.function.IsAsyncGeneratorFunction() {
			r.ConvertYieldOfAsyncGenerator(ctx, n)
		}
	case rhino.Token_RETURN:
		rhino.CheckNotNull(ctx.function)
		if ctx.function.IsAsyncGeneratorFunction() {
			r.ConvertReturnOfAsyncGenerator(ctx, n)
		}
		// For-Await-Of loops
	case rhino.Token_FOR_AWAIT_OF:
		rhino.CheckNotNull(ctx.function)
		jutil.CheckState(ctx.function.IsAsyncFunction())
		r.ReplaceForAwaitOf(ctx, n)
		NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_CONST_DECLARATIONS, r.compiler)
		// Maintaining references to this/arguments/super
	case rhino.Token_THIS:
		if ctx.MustReplaceThisSuperArgs() {
			r.ReplaceThis(ctx, n)
		}
	case rhino.Token_NAME:
		if ctx.MustReplaceThisSuperArgs() && n.MatchesName("arguments") {
			r.ReplaceArguments(ctx, n)
		}
	case rhino.Token_SUPER:
		if ctx.MustReplaceThisSuperArgs() {
			r.ReplaceSuper(ctx, n, parent)
		}
	default:
	}
}

// ConvertAsyncGenerator moves the body of an async generator function into a
// nested generator function and removes the async and generator props from
// the original function.
//
//	async function* foo() {
//	  bar();
//	}
//
// becomes
//
//	function foo() {
//	  return new $jscomp.AsyncGeneratorWrapper((function*(){
//	    bar();
//	  })())
//	}
//
// originalFunction is the original AsyncGeneratorFunction Node to be
// converted.
//
// Java: RewriteAsyncIteration.convertAsyncGenerator(Node)
func (r *RewriteAsyncIteration) ConvertAsyncGenerator(originalFunction rhino.Node) {
	rhino.CheckNotNull(originalFunction)
	jutil.CheckState(originalFunction.IsAsyncGeneratorFunction())

	asyncGeneratorWrapperRef := r.astFactory.CreateQNameOfField(r.namespace, r.asyncGeneratorWrapper)
	innerFunction := r.astFactory.CreateEmptyAsyncGeneratorWrapperArgument(nil)

	innerBlock := originalFunction.GetLastChild()
	innerBlock.Detach()
	innerFunction.GetLastChild().ReplaceWith(innerBlock)

	// Body should be:
	// return new $jscomp.AsyncGeneratorWrapper((new function with original block here)());
	outerBlock :=
		r.astFactory.CreateBlock(
			r.astFactory.CreateReturn(
				r.astFactory.CreateNewNode(
					asyncGeneratorWrapperRef,
					r.astFactory.CreateCall(innerFunction, AstFactory_Type_OfColorId(colors.StandardColors_GENERATOR_ID)))))
	originalFunction.AddChildToBack(outerBlock)

	originalFunction.SetIsAsyncFunction(false)
	originalFunction.SetIsGeneratorFunction(false)
	originalFunction.SrcrefTreeIfMissing(originalFunction)
	// Both the inner and original functions should be marked as changed.
	r.compiler.ReportChangeToChangeScope(originalFunction)
	r.compiler.ReportChangeToChangeScope(innerFunction)
}

// ConvertAwaitOfAsyncGenerator converts an await into a yield of an
// ActionRecord to perform "AWAIT".
//
//	await myPromise
//
// becomes
//
//	yield new ActionRecord(ActionEnum.AWAIT_VALUE, myPromise)
//
// awaitNode is the original await Node to be converted
//
// Java: RewriteAsyncIteration.convertAwaitOfAsyncGenerator(LexicalContext, Node)
func (r *RewriteAsyncIteration) ConvertAwaitOfAsyncGenerator(ctx *RewriteAsyncIteration_LexicalContext, awaitNode rhino.Node) {
	rhino.CheckNotNull(awaitNode)
	jutil.CheckState(awaitNode.IsAwait())
	jutil.CheckState(ctx != nil && !ctx.function.IsNil())
	jutil.CheckState(ctx.function.IsAsyncGeneratorFunction())

	expression := awaitNode.RemoveFirstChild()
	rhino.CheckNotNullMsg(expression, "await needs an expression")
	newActionRecord :=
		r.astFactory.CreateNewNode(
			r.astFactory.CreateQNameOfField(r.namespace, r.actionRecordName),
			r.astFactory.CreateQNameOfField(r.namespace, r.actionEnumAwait),
			expression)
	newActionRecord.SrcrefTreeIfMissing(awaitNode)
	awaitNode.AddChildToFront(newActionRecord)
	awaitNode.SetToken(rhino.Token_YIELD)
}

// ConvertYieldOfAsyncGenerator converts a yield into a yield of an
// ActionRecord to perform "YIELD" or "YIELD_STAR".
//
//	yield;
//	yield first;
//	yield* second;
//
// becomes
//
//	yield new ActionRecord(ActionEnum.YIELD_VALUE, undefined);
//	yield new ActionRecord(ActionEnum.YIELD_VALUE, first);
//	yield new ActionRecord(ActionEnum.YIELD_STAR, second);
//
// yieldNode is the Node to be converted
//
// Java: RewriteAsyncIteration.convertYieldOfAsyncGenerator(LexicalContext, Node)
func (r *RewriteAsyncIteration) ConvertYieldOfAsyncGenerator(ctx *RewriteAsyncIteration_LexicalContext, yieldNode rhino.Node) {
	rhino.CheckNotNull(yieldNode)
	jutil.CheckState(yieldNode.IsYield())
	jutil.CheckState(ctx != nil && !ctx.function.IsNil())
	jutil.CheckState(ctx.function.IsAsyncGeneratorFunction())

	expression := yieldNode.RemoveFirstChild()
	newActionRecord :=
		r.astFactory.CreateNewNode(r.astFactory.CreateQNameOfField(r.namespace, r.actionRecordName))

	if yieldNode.IsYieldAll() {
		rhino.CheckNotNull(expression)
		// yield* expression becomes new ActionRecord(YIELD_STAR, expression)
		newActionRecord.AddChildToBack(r.astFactory.CreateQNameOfField(r.namespace, r.actionEnumYieldStar))
		newActionRecord.AddChildToBack(expression)
	} else {
		if expression.IsNil() {
			expression = NodeUtil_NewUndefinedNode(rhino.NilNode)
		}
		// yield expression becomes new ActionRecord(YIELD, expression)
		newActionRecord.AddChildToBack(r.astFactory.CreateQNameOfField(r.namespace, r.actionEnumYield))
		newActionRecord.AddChildToBack(expression)
	}

	newActionRecord.SrcrefTreeIfMissing(yieldNode)
	yieldNode.AddChildToFront(newActionRecord)
	yieldNode.PutBooleanProp(rhino.Node_YIELD_ALL, false)
}

// ConvertReturnOfAsyncGenerator converts a return into a return of an
// ActionRecord.
//
//	return;
//	return value;
//
// becomes
//
//	return new ActionRecord(ActionEnum.YIELD_VALUE, undefined);
//	return new ActionRecord(ActionEnum.YIELD_VALUE, value);
//
// returnNode is the Node to be converted
//
// Java: RewriteAsyncIteration.convertReturnOfAsyncGenerator(LexicalContext, Node)
func (r *RewriteAsyncIteration) ConvertReturnOfAsyncGenerator(ctx *RewriteAsyncIteration_LexicalContext, returnNode rhino.Node) {
	rhino.CheckNotNull(returnNode)
	jutil.CheckState(returnNode.IsReturn())
	jutil.CheckState(ctx != nil && !ctx.function.IsNil())
	jutil.CheckState(ctx.function.IsAsyncGeneratorFunction())

	expression := returnNode.RemoveFirstChild()
	newActionRecord :=
		r.astFactory.CreateNewNode(r.astFactory.CreateQNameOfField(r.namespace, r.actionRecordName))

	if expression.IsNil() {
		expression = NodeUtil_NewUndefinedNode(rhino.NilNode)
	}
	// return expression becomes new ActionRecord(YIELD, expression)
	newActionRecord.AddChildToBack(r.astFactory.CreateQNameOfField(r.namespace, r.actionEnumYield))
	newActionRecord.AddChildToBack(expression)
	newActionRecord.SrcrefTreeIfMissing(returnNode)
	returnNode.AddChildToFront(newActionRecord)
}

// ReplaceForAwaitOf rewrites for await of loop.
//
//	for await (lhs of rhs) { block(); }
//
// ...becomes...
//
//	var errorRes, retFn, tmpRes;
//	try {
//	  for (var tmpIterator = makeAsyncIterator(rhs);;) {
//	     tmpRes = void 0;
//	     tmpRes = await tmpIterator.next();
//	     if (tmpRes.done) {
//	       break;
//	     }
//	     lhs = $tmpRes.value;
//	     {
//	       block(); // Wrapped in a block in case block re-declares lhs variable.
//	     }
//	  }
//	} catch(e) {
//	  errorRes = { error: e };
//	} finally {
//	  try {
//	    if (tmpRes && !tmpRes.done && (retFn = _tmpIterator.return)) await retFn.call(tmpIterator);
//	  }
//	  finally { if (errorRes) throw errorRes.error; }
//	}
//
// Java: RewriteAsyncIteration.replaceForAwaitOf(LexicalContext, Node)
func (r *RewriteAsyncIteration) ReplaceForAwaitOf(ctx *RewriteAsyncIteration_LexicalContext, forAwaitOf rhino.Node) {
	forAwaitId := r.nextForAwaitId
	r.nextForAwaitId++
	id := strconv.Itoa(forAwaitId)
	iteratorTempName := RewriteAsyncIteration_FOR_AWAIT_ITERATOR_TEMP_NAME + id
	resultTempName := RewriteAsyncIteration_FOR_AWAIT_RESULT_TEMP_NAME + id
	errorResultTempName := RewriteAsyncIteration_FOR_AWAIT_ERROR_RESULT_TEMP_NAME + id
	catchErrorParamTempName := RewriteAsyncIteration_FOR_AWAIT_CATCH_PARAM_TEMP_NAME + id
	returnFuncTempName := RewriteAsyncIteration_FOR_AWAIT_RETURN_FN_TEMP_NAME + id

	jutil.CheckStateMsg(forAwaitOf.HasParent(), "Cannot replace parentless for-await-of")

	forAwaitOfParent := forAwaitOf.GetParent()
	var replacementPoint rhino.Node
	if forAwaitOfParent.IsLabel() {
		// If the forAwaitOf is a label's statement child, then the label must move with the for upon
		// rewriting.
		jutil.CheckStateMsg1(forAwaitOf.IsSecondChildOf(forAwaitOfParent), "%s", forAwaitOfParent)
		replacementPoint = forAwaitOfParent
	} else {
		replacementPoint = forAwaitOf
	}

	lhs := forAwaitOf.RemoveFirstChild()
	rhs := forAwaitOf.RemoveFirstChild()
	originalBody := forAwaitOf.RemoveFirstChild()

	// Generate `var tmpIterator = makeAsyncIterator(rhs);`
	initializer :=
		r.astFactory.
			CreateSingleVarNameDeclarationWithValue(
				iteratorTempName, r.astFactory.CreateJSCompMakeAsyncIteratorCall(rhs, r.namespace)).
			SrcrefTreeIfMissing(rhs)

	// IIterableResult<VALUE> - it's a structural type so optimizations treat it as Object
	iterableResultType := AstFactory_Type_OfColor(colors.StandardColors_TOP_OBJECT)

	// Create code `if (tmpRes.done) {break;}`
	breakIfDone :=
		r.astFactory.CreateIf(
			r.astFactory.CreateGetProp(
				r.astFactory.CreateName(resultTempName, iterableResultType),
				"done",
				AstFactory_Type_OfColor(colors.StandardColors_BOOLEAN)),
			r.astFactory.CreateBlock(r.astFactory.CreateBreak()))

	// Assignment statement to be moved from lhs into body of new for-loop
	var lhsAssignment rhino.Node
	var resultType AstFactory_Type
	if lhs.IsValidAssignmentTarget() {
		// In case of "for await (x of _)" just assign into the lhs.
		// Generate `lhs = $tmpRes.value;`
		resultType = AstFactory_Type_(lhs)
		lhsAssignment =
			r.astFactory.ExprResult(
				r.astFactory.CreateAssign(
					lhs,
					r.astFactory.CreateGetProp(
						r.astFactory.CreateName(resultTempName, iterableResultType),
						"value",
						resultType)))
	} else if NodeUtil_IsNameDeclaration(lhs) {
		declarationTarget := lhs.GetFirstChild()
		if declarationTarget.IsName() {
			// `for await (let x of _)`
			// Add a child to the `NAME` node to create `let x = res.value`
			resultType = AstFactory_Type_(declarationTarget)
			declarationTarget.AddChildToBack(
				r.astFactory.CreateGetProp(
					r.astFactory.CreateName(resultTempName, iterableResultType), "value", resultType))
		} else {
			// Generate `for await (let [x, y] of _)`
			// Add a child to the DESTRUCTURING_LHS node to create `[x, y] = res.value`
			jutil.CheckStateMsg1(declarationTarget.IsDestructuringLhs(), "%s", declarationTarget)
			destructuringPattern := declarationTarget.GetOnlyChild()
			resultType = AstFactory_Type_(destructuringPattern)
			declarationTarget.AddChildToBack(
				r.astFactory.CreateGetProp(
					r.astFactory.CreateName(resultTempName, iterableResultType), "value", resultType))
		}
		lhsAssignment = lhs
	} else {
		panic(jutil.NewAssertionError("unexpected for-await-of lhs"))
	}
	lhsAssignment.SrcrefTreeIfMissing(lhs)

	// Generate `var errorRes;`
	errorResDecl :=
		r.astFactory.
			CreateSingleVarNameDeclaration(errorResultTempName).
			SrcrefTreeIfMissing(forAwaitOf)

	// Generate `var tmpRes;`
	tempResultDecl :=
		r.astFactory.CreateSingleVarNameDeclaration(resultTempName).SrcrefTreeIfMissing(forAwaitOf)

	// Generate `var returnFunc;`
	returnFuncDecl :=
		r.astFactory.
			CreateSingleVarNameDeclaration(returnFuncTempName).
			SrcrefTreeIfMissing(forAwaitOf)

	// Generate `tmpRes = void 0;`
	resetTempResult :=
		r.astFactory.ExprResult(
			r.astFactory.CreateAssignToName(resultTempName, r.astFactory.CreateUndefinedValue()))

	// Generate `tmpRes = await tmpIterator.next()`
	resultDeclaration :=
		r.astFactory.ExprResult(
			r.astFactory.CreateAssignToName(
				resultTempName,
				r.ConstructAwaitNextResult(ctx, iteratorTempName, resultType, iterableResultType)))

	newForLoop :=
		r.astFactory.CreateFor(
			r.astFactory.CreateEmpty(),
			r.astFactory.CreateEmpty(),
			r.astFactory.CreateEmpty(),
			r.astFactory.CreateBlock(
				resetTempResult,
				resultDeclaration,
				breakIfDone,
				lhsAssignment,
				r.EnsureBlock(originalBody)))

	if replacementPoint.IsLabel() {
		newForLoop = r.astFactory.CreateLabel(replacementPoint.GetFirstChild().CloneNode(), newForLoop)
	}

	// Generates code `try { .. newForLoop .. }`
	tryNode := r.CreateOuterTry(newForLoop)
	initializer.InsertBefore(newForLoop)

	// Generate code `catch(e) { errorRes = { error: e }; }`
	catchNode := r.CreateOuterCatch(catchErrorParamTempName, errorResultTempName)

	// Generate the finally code block.
	finallyNode :=
		r.CreateOuterFinally(
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

// CreateOuterTry generates code `try { .. newForLoop .. }`
//
// Java: RewriteAsyncIteration.createOuterTry(Node)
func (r *RewriteAsyncIteration) CreateOuterTry(newForLoop rhino.Node) rhino.Node {
	tryNode := r.astFactory.CreateBlock()
	tryNode.AddChildToBack(newForLoop)
	return tryNode
}

// CreateOuterCatch generates code `catch(e) { errorRes = { error: e }; }`
//
// Java: RewriteAsyncIteration.createOuterCatch(String, String)
func (r *RewriteAsyncIteration) CreateOuterCatch(catchErrorParamTempName string, errorResultTempName string) rhino.Node {
	// Generate `errorRes = { error: e };`
	catchBodyStmt :=
		r.astFactory.ExprResult(
			r.astFactory.CreateAssignToName(
				errorResultTempName,
				r.astFactory.CreateObjectLit(
					r.astFactory.CreateStringKey(
						"error", r.astFactory.CreateNameWithUnknownType(catchErrorParamTempName)))))
	// Generate `{ errorRes = { error: e }; }`
	wrapperCatchBlockNode := r.astFactory.CreateBlock()
	wrapperCatchBlockNode.AddChildToBack(catchBodyStmt)

	// Generate `catch(e) { errorRes = { error: e }; }`
	return r.astFactory.CreateCatch(
		r.astFactory.CreateNameWithUnknownType(catchErrorParamTempName), wrapperCatchBlockNode)
}

// CreateOuterFinally generates the outer finally code of the rewriting.
//
//	finally {
//	  try {
//	    if (tmpRes && !tmpRes.done && (retFn = _tmpIterator.return)) await retFn.call(tmpIterator);
//	  }
//	  finally { if (errorRes) throw errorRes.error; }
//	}
//
// Java: RewriteAsyncIteration.createOuterFinally(LexicalContext, AstFactory.Type, AstFactory.Type, String, String, String, String)
func (r *RewriteAsyncIteration) CreateOuterFinally(
	ctx *RewriteAsyncIteration_LexicalContext,
	iterableResultType AstFactory_Type,
	resultType AstFactory_Type,
	resultTempName string,
	returnFuncTempName string,
	iteratorTempName string,
	errorResultTempName string) rhino.Node {
	finallyNode := r.astFactory.CreateBlock()

	// Generate `tmpRes`
	tmpResNameNode := r.astFactory.CreateNameWithUnknownType(resultTempName)
	tmpResDoneGetProp :=
		r.astFactory.CreateGetProp(
			r.astFactory.CreateName(resultTempName, iterableResultType),
			"done",
			AstFactory_Type_OfColor(colors.StandardColors_BOOLEAN))

	// Generate `tmpRes && !tmpRes.done`
	and := r.astFactory.CreateAnd(tmpResNameNode, r.astFactory.CreateNot(tmpResDoneGetProp))

	// Generate `(retFn = _tmpIterator.return)`
	assign :=
		r.astFactory.CreateAssign(
			r.astFactory.CreateNameWithUnknownType(returnFuncTempName),
			r.astFactory.CreateGetProp(
				r.astFactory.CreateName(iteratorTempName, resultType),
				"return",
				AstFactory_Type_OfColor(colors.StandardColors_UNKNOWN)))

	// Generate `(tmpRes && !tmpRes.done && (retFn = _tmpIterator.return))`
	ifCond := r.astFactory.CreateAnd(and, assign)
	var awaitOrYieldStmt rhino.Node = rhino.NilNode
	if ctx.function.IsAsyncGeneratorFunction() {
		// We are in an AsyncGenerator and must instead yield an "await" ActionRecord
		awaitOrYieldStmt =
			r.astFactory.ExprResult(
				r.astFactory.CreateYield(
					iterableResultType,
					r.astFactory.CreateNewNode(
						r.astFactory.CreateQNameOfField(r.namespace, r.actionRecordName),
						r.astFactory.CreateQNameOfField(r.namespace, r.actionEnumAwait),
						r.astFactory.CreateCall(
							r.astFactory.CreateGetPropWithUnknownType(
								r.astFactory.CreateName(
									returnFuncTempName, AstFactory_Type_OfColor(colors.StandardColors_UNKNOWN)),
								"call"),
							AstFactory_Type_OfColor(colors.StandardColors_UNKNOWN),
							r.astFactory.CreateName(iteratorTempName, resultType)))))
	} else {
		//  Generate `await retFn.call(tmpIterator);`
		awaitOrYieldStmt =
			r.astFactory.ExprResult(
				r.astFactory.CreateAwait(
					iterableResultType,
					r.astFactory.CreateCall(
						r.astFactory.CreateGetPropWithUnknownType(
							r.astFactory.CreateName(returnFuncTempName, AstFactory_Type_OfColor(colors.StandardColors_UNKNOWN)),
							"call"),
						AstFactory_Type_OfColorId(colors.StandardColors_PROMISE_ID),
						r.astFactory.CreateName(iteratorTempName, resultType))))
	}

	ifBody := r.astFactory.CreateBlock()
	ifBody.AddChildToBack(awaitOrYieldStmt)

	ifBlock := r.astFactory.CreateIf(ifCond, ifBody)

	innerTryBlock := r.astFactory.CreateBlock()
	innerTryBlock.AddChildToBack(ifBlock)

	//  `finally { if (errorRes) throw errorRes.error; }`
	innerFinallyBlock := r.astFactory.CreateBlock()

	// if (errorRes) throw errorRes.error;
	secondIfBody := r.astFactory.CreateBlock()
	throwStmt :=
		r.astFactory.CreateThrow(
			r.astFactory.CreateGetPropWithUnknownType(
				r.astFactory.CreateNameWithUnknownType(errorResultTempName), "error"))
	secondIfBody.AddChildToBack(throwStmt)
	secondIfCond := r.astFactory.CreateNameWithUnknownType(errorResultTempName)
	secondIfBlock := r.astFactory.CreateIf(secondIfCond, secondIfBody)
	innerFinallyBlock.AddChildToBack(secondIfBlock)

	finallyBody := r.astFactory.CreateTryFinally(innerTryBlock, innerFinallyBlock)
	finallyNode.AddChildToBack(finallyBody)
	return finallyNode
}

// Java: RewriteAsyncIteration.ensureBlock(Node)
func (r *RewriteAsyncIteration) EnsureBlock(possiblyBlock rhino.Node) rhino.Node {
	if possiblyBlock.IsBlock() {
		return possiblyBlock
	}
	return r.astFactory.CreateBlock(possiblyBlock).Srcref(possiblyBlock)
}

// Java: RewriteAsyncIteration.constructAwaitNextResult(LexicalContext, String, AstFactory.Type, AstFactory.Type)
func (r *RewriteAsyncIteration) ConstructAwaitNextResult(
	ctx *RewriteAsyncIteration_LexicalContext,
	iteratorTempName string,
	iteratorType AstFactory_Type,
	iterableResultType AstFactory_Type) rhino.Node {
	rhino.CheckNotNull(ctx.function)
	var result rhino.Node

	iteratorTemp := r.astFactory.CreateName(iteratorTempName, iteratorType)

	if ctx.function.IsAsyncGeneratorFunction() {
		// We are in an AsyncGenerator and must instead yield an "await" ActionRecord
		result =
			r.astFactory.CreateYield(
				iterableResultType,
				r.astFactory.CreateNewNode(
					r.astFactory.CreateQNameOfField(r.namespace, r.actionRecordName),
					r.astFactory.CreateQNameOfField(r.namespace, r.actionEnumAwait),
					r.astFactory.CreateCallWithUnknownType(
						r.astFactory.CreateGetPropWithUnknownType(iteratorTemp, "next"))))
	} else {
		result =
			r.astFactory.CreateAwait(
				iterableResultType,
				r.astFactory.CreateCall(
					r.astFactory.CreateGetPropWithUnknownType(iteratorTemp, "next"),
					AstFactory_Type_OfColorId(colors.StandardColors_PROMISE_ID)))
	}

	return result
}

// Java: RewriteAsyncIteration.replaceThis(LexicalContext, Node)
func (r *RewriteAsyncIteration) ReplaceThis(ctx *RewriteAsyncIteration_LexicalContext, n rhino.Node) {
	jutil.CheckArgument(n.IsThis())
	jutil.CheckArgument(ctx != nil && ctx.MustReplaceThisSuperArgs())
	jutil.CheckArgumentMsg(!ctx.function.IsNil(), "Cannot prepend declarations to root scope")
	jutil.CheckNotNull(ctx.thisSuperArgsContext)

	n.ReplaceWith(
		r.astFactory.
			CreateName(RewriteAsyncIteration_THIS_VAR_NAME+ctx.thisSuperArgsContext.uniqueId, AstFactory_Type_(n)).
			Srcref(n))
	ctx.thisSuperArgsContext.thisNodeToAdd = r.astFactory.CreateThis(AstFactory_Type_(n))
	r.compiler.ReportChangeToChangeScope(ctx.function)
}

// Java: RewriteAsyncIteration.replaceArguments(LexicalContext, Node)
func (r *RewriteAsyncIteration) ReplaceArguments(ctx *RewriteAsyncIteration_LexicalContext, n rhino.Node) {
	jutil.CheckArgument(n.IsName() && "arguments" == n.GetString())
	jutil.CheckArgument(ctx != nil && ctx.MustReplaceThisSuperArgs())
	jutil.CheckArgumentMsg(!ctx.function.IsNil(), "Cannot prepend declarations to root scope")
	jutil.CheckNotNull(ctx.thisSuperArgsContext)

	n.ReplaceWith(r.astFactory.CreateName(RewriteAsyncIteration_ARGUMENTS_VAR_NAME, AstFactory_Type_(n)).Srcref(n))
	ctx.thisSuperArgsContext.usedArguments = true
	r.compiler.ReportChangeToChangeScope(ctx.function)
}

// Java: RewriteAsyncIteration.replaceSuper(LexicalContext, Node, Node)
func (r *RewriteAsyncIteration) ReplaceSuper(ctx *RewriteAsyncIteration_LexicalContext, n rhino.Node, parent rhino.Node) {
	if !parent.IsGetProp() {
		r.compiler.Report(
			JSError_MakeWithNode(
				parent,
				RewriteAsyncIteration_CANNOT_CONVERT_ASYNCGEN,
				"super only allowed with getprop (like super.foo(), not super['foo']())"))
		return
	}
	jutil.CheckArgument(n.IsSuper())
	jutil.CheckArgument(ctx != nil && ctx.MustReplaceThisSuperArgs())
	jutil.CheckArgumentMsg(!ctx.function.IsNil(), "Cannot prepend declarations to root scope")
	jutil.CheckNotNull(ctx.thisSuperArgsContext)

	propertyName := parent.GetString()
	propertyReplacementNameText := RewriteAsyncIteration_SUPER_PROP_GETTER_PREFIX + propertyName

	// super.x   =>   $super$get$x()
	getPropReplacement :=
		r.astFactory.CreateCall(
			r.astFactory.CreateName(propertyReplacementNameText, AstFactory_Type_OfColor(colors.StandardColors_TOP_OBJECT)),
			AstFactory_Type_(parent))
	grandparent := parent.GetParent()
	if grandparent.IsCall() && grandparent.GetFirstChild() == parent {
		// super.x(args) => super.x.call($this, args)
		getPropReplacement = r.astFactory.CreateGetPropWithUnknownType(getPropReplacement, "call")
		ctx.thisSuperArgsContext.thisNodeToAdd =
			r.astFactory.CreateThisForEs6ClassMember(ctx.contextRoot.GetParent())
		r.astFactory.
			CreateName(
				RewriteAsyncIteration_THIS_VAR_NAME+ctx.thisSuperArgsContext.uniqueId,
				AstFactory_Type_(ctx.thisSuperArgsContext.thisNodeToAdd)).
			Srcref(parent).
			InsertAfter(parent)
	} else if grandparent.IsTaggedTemplateLit() && grandparent.GetFirstChild() == parent {
		// super.x`template` => super.x().bind($this)`template`
		ctx.thisSuperArgsContext.thisNodeToAdd =
			r.astFactory.CreateThisForEs6ClassMember(ctx.contextRoot.GetParent())
		thisAlias :=
			r.astFactory.
				CreateName(
					RewriteAsyncIteration_THIS_VAR_NAME+ctx.thisSuperArgsContext.uniqueId,
					AstFactory_Type_(ctx.thisSuperArgsContext.thisNodeToAdd)).
				Srcref(parent)
		getPropReplacement =
			r.astFactory.CreateCallWithUnknownType(
				r.astFactory.CreateGetPropWithUnknownType(getPropReplacement, "bind"), thisAlias)
		grandparent.PutBooleanProp(rhino.Node_FREE_CALL, true)
	}
	getPropReplacement.SrcrefTree(parent)
	parent.ReplaceWith(getPropReplacement)
	ctx.thisSuperArgsContext.usedSuperProperties.Add(parent)
	r.compiler.ReportChangeToChangeScope(ctx.function)
}

// PrependTempVarDeclarations prepends this/super/argument replacement
// variables to the top of the context's block
//
//	function() {
//	  return new AsyncGenWrapper(function*() {
//	    // code using replacements for this and super.foo
//	  }())
//	}
//
// will be converted to
//
//	function() {
//	  const $jscomp$asyncIter$this = this;
//	  const $jscomp$asyncIter$super$get$foo = () => super.foo;
//	  return new AsyncGenWrapper(function*() {
//	    // code using replacements for this and super.foo
//	  }())
//	}
//
// Java: RewriteAsyncIteration.prependTempVarDeclarations(LexicalContext, NodeTraversal)
func (r *RewriteAsyncIteration) PrependTempVarDeclarations(ctx *RewriteAsyncIteration_LexicalContext, t *NodeTraversal) {
	jutil.CheckArgument(ctx != nil)
	jutil.CheckArgumentMsg(!ctx.function.IsNil(), "Cannot prepend declarations to root scope")
	jutil.CheckNotNull(ctx.thisSuperArgsContext)

	thisSuperArgsCtx := ctx.thisSuperArgsContext
	function := ctx.function
	block := function.GetLastChild()
	rhino.CheckNotNullMsg1(block, "%s", function)
	prefixBlock := r.astFactory.CreateBlock() // Temporary block to hold all declarations

	if !thisSuperArgsCtx.thisNodeToAdd.IsNil() {
		// { // prefixBlock
		//   const $jscomp$asyncIter$this = this;
		// }
		prefixBlock.AddChildToBack(
			r.astFactory.
				CreateSingleConstNameDeclaration(
					RewriteAsyncIteration_THIS_VAR_NAME+thisSuperArgsCtx.uniqueId, thisSuperArgsCtx.thisNodeToAdd).
				SrcrefTree(block))
	}
	if thisSuperArgsCtx.usedArguments {
		// { // prefixBlock
		//   const $jscomp$asyncIter$this = this;
		//   const $jscomp$asyncIter$arguments = arguments;
		// }
		prefixBlock.AddChildToBack(
			r.astFactory.
				CreateSingleConstNameDeclaration(
					RewriteAsyncIteration_ARGUMENTS_VAR_NAME, r.astFactory.CreateArgumentsReference()).
				SrcrefTree(block))
	}
	for replacedMethodReference := range thisSuperArgsCtx.usedSuperProperties.All() {
		prefixBlock.AddChildToBack(r.CreateSuperMethodReferenceGetter(replacedMethodReference, t))
	}
	prefixBlock.SrcrefTreeIfMissing(block)
	// Pulls all declarations out of prefixBlock and prepends in block
	// block: {
	//   // declarations
	//   // code using this/super/args
	// }
	block.AddChildrenToFront(prefixBlock.RemoveChildren())

	if !thisSuperArgsCtx.thisNodeToAdd.IsNil() ||
		thisSuperArgsCtx.usedArguments ||
		!thisSuperArgsCtx.usedSuperProperties.IsEmpty() {
		r.compiler.ReportChangeToChangeScope(function)
		NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_CONST_DECLARATIONS, r.compiler)
	}
}

// Java: RewriteAsyncIteration.createSuperMethodReferenceGetter(Node, NodeTraversal)
func (r *RewriteAsyncIteration) CreateSuperMethodReferenceGetter(replacedMethodReference rhino.Node, t *NodeTraversal) rhino.Node {
	// const super$get$x = () => { return super.x; };
	typeOfSuper := AstFactory_Type_(replacedMethodReference.GetFirstChild())
	superReference := r.astFactory.CreateSuper(typeOfSuper)
	replacedMethodName := replacedMethodReference.GetString()
	arrowFunction :=
		r.astFactory.CreateZeroArgArrowFunctionForExpression(
			r.astFactory.CreateBlock(
				r.astFactory.CreateReturn(
					r.astFactory.CreateGetProp(
						superReference, replacedMethodName, AstFactory_Type_(replacedMethodReference)))))
	r.compiler.ReportChangeToChangeScope(arrowFunction)
	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_ARROW_FUNCTIONS, r.compiler)
	superReplacementName := RewriteAsyncIteration_SUPER_PROP_GETTER_PREFIX + replacedMethodName
	return r.astFactory.CreateSingleConstNameDeclaration(superReplacementName, arrowFunction)
}
