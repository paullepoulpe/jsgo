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
//
// An async function, foo(a, b), will be rewritten as:
//
//	function foo(a, b) {
//	  let $jscomp$async$this = this;
//	  let $jscomp$async$arguments = arguments;
//	  let $jscomp$async$new$target = new.target;
//	  let $jscomp$async$super$get$x = () => super.x;
//	  return $jscomp.asyncExecutePromiseGeneratorFunction(
//	      function* () {
//	        // original body of foo() with:
//	        // - await (x) replaced with yield (x)
//	        // - arguments replaced with $jscomp$async$arguments
//	        // - this replaced with $jscomp$async$this
//	        // - new.target replaced with $jscomp$async$new$target
//	        // - super.x replaced with $jscomp$async$super$get$x()
//	        // - super.x(5) replaced with $jscomp$async$super$get$x().call($jscomp$async$this, 5)
//	      });
//	}
type RewriteAsyncFunctions struct {
	namespace    rhino.StaticScope
	contextStack []*rewriteAsyncLexicalContext
	compiler     AbstractCompiler
	astFactory   *AstFactory
}

const (
	// Java: RewriteAsyncFunctions.ASYNC_ARGUMENTS
	RewriteAsyncFunctions_ASYNC_ARGUMENTS = "$jscomp$async$arguments$"
	// Java: RewriteAsyncFunctions.ASYNC_THIS
	RewriteAsyncFunctions_ASYNC_THIS = "$jscomp$async$this$"
	// Java: RewriteAsyncFunctions.ASYNC_NEW_TARGET
	RewriteAsyncFunctions_ASYNC_NEW_TARGET = "$jscomp$async$new$target$"
	// Java: RewriteAsyncFunctions.ASYNC_SUPER_PROP_GETTER_PREFIX
	RewriteAsyncFunctions_ASYNC_SUPER_PROP_GETTER_PREFIX = "$jscomp$async$super$get$"
)

// asyncSuperPropertyWrapperInfo holds the information needed to replace a
// reference to `super.propertyName` with a call to a wrapper function.
type asyncSuperPropertyWrapperInfo struct {
	// The first `super.property` Node we come across during traversal.
	// Information from this node will be used when creating the wrapper function.
	firstInstanceOfSuperDotProperty rhino.Node
	wrapperFunctionName             string
	// The type to use for the wrapper function.
	// Will be null if type checking has not run.
	wrapperFunctionReturnType AstFactory_Type
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

// rewriteAsyncLexicalContext determines both what to do when visiting a node
// and how to determine the context for its descendants.
type rewriteAsyncLexicalContext struct {
	kind            asyncContextKind
	pass            *RewriteAsyncFunctions
	contextRootNode rhino.Node
	// precomputed unique id to append to names in this context. This is used to ensure that names
	// in different contexts don't collide (e.g. 2 functions don't get the same `let
	// $jscomp$async$this` name declared in their bodies)
	uniqueId string

	// Set when kind == asyncContextParamList.
	functionContext *rewriteAsyncLexicalContext

	// Set when kind == asyncContextFunction: points to the enclosing async FunctionContext whose
	// `this`, `arguments`, `new.target`, and `super` bindings apply here, or nil if none.
	asyncThisAndArgumentsContext *rewriteAsyncLexicalContext
	// Lazily allocated on async contexts that actually reference `super.prop`.
	superPropertyWrappers *jutil.LinkedHashMap[string, asyncSuperPropertyWrapperInfo]

	mustAddAsyncThisVariable bool
	// null if mustAddAsyncThisVariable is false
	typeOfThis                    AstFactory_Type
	mustAddAsyncArgumentsVariable bool
	mustAddAsyncNewTargetVariable bool
	// null if mustAddAsyncNewTargetVariable is false
	typeOfNewTarget AstFactory_Type
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

// outerAsyncContext is the async context that `this`, `arguments`, etc.
// refer to at the point where this function is defined, or nil if there is none.
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
			// An async arrow function inside an outer async context shares outerAsyncContext.
			c.asyncThisAndArgumentsContext = outerAsyncContext
		} else {
			// An async non-arrow function (or top-level async arrow) owns its async context.
			c.asyncThisAndArgumentsContext = c
		}
	} else if fn.IsArrow() {
		// A non-async arrow function inherits outerAsyncContext (which may be nil).
		c.asyncThisAndArgumentsContext = outerAsyncContext
	}
	return c
}

// asyncContextForParameterList returns the async context that `this`,
// `arguments`, `new.target`, and `super` refer to in this parameter list, or
// nil if they should be left unchanged.
//
// When the function owning this parameter list is itself the async context,
// the parameter list stays on the outer function and is not moved into the
// generator body, so references here (including inside nested arrow
// functions) must not be aliased.
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
			// Function defined within a parameter list, e.g. `function someFunc(cb = () => something) {}`
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
	return c.pass.astFactory.CreateName(RewriteAsyncFunctions_ASYNC_THIS+c.uniqueId, typeOfThis)
}

func (c *rewriteAsyncLexicalContext) createNewTargetVariableReference(typeOfNewTarget AstFactory_Type) rhino.Node {
	c.mustAddAsyncNewTargetVariable = true
	c.typeOfNewTarget = typeOfNewTarget
	return c.pass.astFactory.CreateName(RewriteAsyncFunctions_ASYNC_NEW_TARGET+c.uniqueId, typeOfNewTarget)
}

func (c *rewriteAsyncLexicalContext) getOrCreateSuperPropertyWrapperInfo(superDotPropertyNode rhino.Node) asyncSuperPropertyWrapperInfo {
	gp, ok := ast.AsGetProp(superDotPropertyNode)
	jutil.CheckArgumentMsg1(ok, "%s", superDotPropertyNode)
	superNode := gp.Object().Node()
	jutil.CheckArgumentMsg1(superNode.IsSuper(), "%s", superNode)

	propertyName := gp.Property()
	propertyType := superDotPropertyNode.GetColor()
	if c.superPropertyWrappers != nil && c.superPropertyWrappers.ContainsKey(propertyName) {
		info := c.superPropertyWrappers.Get(propertyName)
		// Every reference to `super.propertyName` within a single lexical context should
		// have the same type. Make sure this is true.
		existingColor := info.firstInstanceOfSuperDotProperty.GetColor()
		jutil.CheckStateMsg2(
			jutil.ObjectsEquals(existingColor, propertyType),
			"Previous reference type: %s differs from current reference type: %s",
			existingColor,
			propertyType)
		return info
	}

	if c.superPropertyWrappers == nil {
		c.superPropertyWrappers = jutil.NewLinkedHashMap[string, asyncSuperPropertyWrapperInfo]()
	}
	wrapperFunctionName :=
		RewriteAsyncFunctions_ASYNC_SUPER_PROP_GETTER_PREFIX +
			c.pass.compiler.
				GetUniqueIdSupplier().
				GetUniqueId(c.pass.compiler.GetInput(NodeUtil_GetInputId(superDotPropertyNode))) +
			"$" +
			propertyName
	info := asyncSuperPropertyWrapperInfo{
		firstInstanceOfSuperDotProperty: superDotPropertyNode,
		wrapperFunctionName:             wrapperFunctionName,
		wrapperFunctionReturnType:       AstFactory_Type_(superDotPropertyNode),
	}
	c.superPropertyWrappers.Put(propertyName, info)
	return info
}

func (c *rewriteAsyncLexicalContext) createWrapperArrowFunction(wrapperInfo asyncSuperPropertyWrapperInfo) rhino.Node {
	// () => { return super.propertyName; };
	superDotProperty := wrapperInfo.firstInstanceOfSuperDotProperty.CloneTree()
	return c.pass.astFactory.CreateZeroArgArrowFunctionForExpression(
		c.pass.astFactory.CreateBlock(c.pass.astFactory.CreateReturn(superDotProperty)))
}

func (c *rewriteAsyncLexicalContext) visit(t *NodeTraversal, n rhino.Node) {
	switch c.kind {
	case asyncContextRoot:
		// In root context we haven't entered an async function yet, so there's nothing to do.
	case asyncContextParamList:
		if c.asyncContextForParameterList() != nil {
			// e.g.
			// async function outer(outerT = this) {
			//   // `this` in outer parameter list must remain unchanged
			//   // but inner parameter list must be aliased
			//   const inner = async (t = this) => t;
			// }
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
			// All descendant nodes have been updated as necessary; now convert the async function itself.
			r.convertAsyncFunction(t, c, fn)
		}
		return
	}
	asyncCtx := c.asyncThisAndArgumentsContext
	if asyncCtx == nil {
		return
	}

	// We're in the context of an async function's body, so perform lexical replacements.
	switch n.GetToken() {
	case rhino.Token_NAME:
		nm := ast.MustName(n)
		if nm.Ident() == "arguments" {
			nm.SetIdent(RewriteAsyncFunctions_ASYNC_ARGUMENTS + asyncCtx.uniqueId)
			if r.compiler.GetLifeCycleStage().IsNormalized() {
				// TODO: b/322009741 - Stop depending on lifeCycleStage.isNormalized() to decide constness
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
			// NOTE: `super.prop = x` is valid, and can be useful in overridden setters, but
			// these cannot be async. For now, we don't support this construct in async methods.
			r.compiler.Report(
				JSError_MakeWithNode(
					parent.GetParent(),
					TranspilationUtil_CANNOT_CONVERT_YET,
					"assignment to super property"))
		}
		superDotProperty := parent
		wrapperInfo := asyncCtx.getOrCreateSuperPropertyWrapperInfo(superDotProperty)

		// super.x   =>   $jscomp$super$get$x()
		getPropReplacement := wrapperInfo.createWrapperFunctionCallNode(r.astFactory)
		grandparent := superDotProperty.GetParent()
		if call, ok := ast.AsCall(grandparent); ok && call.Callee().Node() == superDotProperty {
			// $jscomp$super$get$x(args) => $jscomp$super$get$x().call($jscomp$async$this, args)
			getPropReplacement = r.astFactory.CreateGetPropWithUnknownType(getPropReplacement, "call")
			thisAlias := r.astFactory.
				CreateThisAliasReferenceForEs6Class(
					RewriteAsyncFunctions_ASYNC_THIS+asyncCtx.uniqueId,
					NodeUtil_GetEnclosingClass(asyncCtx.contextRootNode)).
				Srcref(superDotProperty)
			thisAlias.InsertAfter(superDotProperty)
			asyncCtx.recordAsyncThisReplacement(AstFactory_Type_(thisAlias))
		} else if ttl, ok := ast.AsTaggedTemplateLit(grandparent); ok && ttl.Tag().Node() == superDotProperty {
			// $jscomp$super$get$x`template` => $jscomp$super$get$x().bind($jscomp$async$this)`template`
			thisAlias := r.astFactory.
				CreateThisAliasReferenceForEs6Class(
					RewriteAsyncFunctions_ASYNC_THIS+asyncCtx.uniqueId,
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
		// Awaits become yields in the converted async function's inner generator function.
		aw := ast.MustAwait(n)
		n.ReplaceWith(r.astFactory.CreateYield(AstFactory_Type_(n), aw.Value().Node().Detach()))
	}
}

// Java: RewriteAsyncFunctions.transpiledFeatures
var RewriteAsyncFunctions_transpiledFeatures = parser.FeatureSet_BARE_MINIMUM.With(parser.FeatureSet_Feature_ASYNC_FUNCTIONS)

// Java: RewriteAsyncFunctions(AbstractCompiler, AstFactory, StaticScope)
func newRewriteAsyncFunctions(
	compiler AbstractCompiler, astFactory *AstFactory, namespace rhino.StaticScope) *RewriteAsyncFunctions {
	return &RewriteAsyncFunctions{
		compiler:   jutil.CheckNotNull(compiler),
		astFactory: jutil.CheckNotNull(astFactory),
		namespace:  jutil.CheckNotNull(namespace),
	}
}

// Java: RewriteAsyncFunctions.create(AbstractCompiler)
func RewriteAsyncFunctions_Create(compiler AbstractCompiler) *RewriteAsyncFunctions {
	astFactory := compiler.CreateAstFactory()
	namespace := compiler.GetTranspilationNamespace()
	return newRewriteAsyncFunctions(compiler, astFactory, namespace)
}

// Java: RewriteAsyncFunctions.process(Node, Node)
func (r *RewriteAsyncFunctions) Process(externs rhino.Node, root rhino.Node) {
	TranspilationPasses_ProcessTranspile(r.compiler, root, RewriteAsyncFunctions_transpiledFeatures, r)
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(r.compiler, root, RewriteAsyncFunctions_transpiledFeatures)
}

func (r *RewriteAsyncFunctions) contextStackPeek() *rewriteAsyncLexicalContext {
	if len(r.contextStack) == 0 {
		return nil
	}
	return r.contextStack[len(r.contextStack)-1]
}

// Java: RewriteAsyncFunctions.shouldTraverse(NodeTraversal, Node, Node)
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

// Java: RewriteAsyncFunctions.visit(NodeTraversal, Node, Node)
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

// Java: RewriteAsyncFunctions.convertAsyncFunction(NodeTraversal, FunctionContext)
func (r *RewriteAsyncFunctions) convertAsyncFunction(
	t *NodeTraversal,
	functionContext *rewriteAsyncLexicalContext,
	originalFunction ast.Function,
) {
	originalFunction.SetAsync(false)
	originalBody := originalFunction.Body().Node()
	if originalFunction.Node().IsFromExterns() {
		// A function defined in externs will never be executed, so we don't need to transpile it.
		// Make sure it has an empty body though so later passes won't trip over uses of `await` or
		// anything like that.
		if !NodeUtil_IsEmptyBlock(originalBody) {
			// TODO(b/119685646): Maybe we should warn for non-empty functions in externs?
			originalBody.ReplaceWith(r.astFactory.CreateBlock())
			NodeUtil_MarkFunctionsDeleted(originalBody, r.compiler)
		}
		return
	}
	newBody := r.astFactory.CreateBlock()
	originalBody.ReplaceWith(newBody)

	asyncCtx := functionContext.asyncThisAndArgumentsContext
	if functionContext.mustAddAsyncThisVariable {
		// const this$ = this;
		newBody.AddChildToBack(
			r.astFactory.CreateSingleConstNameDeclaration(
				RewriteAsyncFunctions_ASYNC_THIS+asyncCtx.uniqueId,
				r.astFactory.CreateThis(functionContext.typeOfThis)))
		NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_CONST_DECLARATIONS, r.compiler)
	}
	if functionContext.mustAddAsyncArgumentsVariable {
		// const arguments$ = arguments;
		newBody.AddChildToBack(
			r.astFactory.CreateArgumentsAliasDeclaration(
				RewriteAsyncFunctions_ASYNC_ARGUMENTS + asyncCtx.uniqueId))
		NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_CONST_DECLARATIONS, r.compiler)
	}
	if functionContext.mustAddAsyncNewTargetVariable {
		// const newTarget$ = new.target;
		newBody.AddChildToBack(
			r.astFactory.CreateSingleConstNameDeclaration(
				RewriteAsyncFunctions_ASYNC_NEW_TARGET+asyncCtx.uniqueId,
				r.astFactory.CreateNewTarget(functionContext.typeOfNewTarget)))
		NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_CONST_DECLARATIONS, r.compiler)
	}
	if functionContext.superPropertyWrappers != nil {
		for wrapperInfo := range functionContext.superPropertyWrappers.Values() {
			arrowFunction := functionContext.createWrapperArrowFunction(wrapperInfo)
			// const super$get$x = () => super.x;
			arrowFunctionDeclarationStatement :=
				r.astFactory.CreateSingleConstNameDeclaration(
					wrapperInfo.wrapperFunctionName, arrowFunction)
			newBody.AddChildToBack(arrowFunctionDeclarationStatement)

			// Make sure the compiler knows about the new arrow function's scope
			r.compiler.ReportChangeToChangeScope(arrowFunction)
			// Record that we've added arrow functions and const declarations to this script,
			// so later transpilations of those features will run, if needed.
			enclosingScript := t.GetCurrentScript()
			NodeUtil_AddFeatureToScript(enclosingScript, parser.FeatureSet_Feature_ARROW_FUNCTIONS, r.compiler)
			NodeUtil_AddFeatureToScript(enclosingScript, parser.FeatureSet_Feature_CONST_DECLARATIONS, r.compiler)
		}
	}

	// Normalize arrow function short body to block body
	if !originalBody.IsBlock() {
		originalBody =
			r.astFactory.
				CreateBlock(r.astFactory.CreateReturn(originalBody)).
				SrcrefTreeIfMissing(originalBody)
	}
	// NOTE: visit() will already have made appropriate replacements in originalBody so it may
	// be used as the generator function body.
	generatorFunction :=
		r.astFactory.CreateZeroArgGeneratorFunction("", originalBody /* returnType= */, nil)
	r.compiler.ReportChangeToChangeScope(generatorFunction)
	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_GENERATORS, r.compiler)

	// return $jscomp.asyncExecutePromiseGeneratorFunction(function* () { ... });
	newBody.AddChildToBack(
		r.astFactory.CreateReturn(
			r.astFactory.CreateJscompAsyncExecutePromiseGeneratorFunctionCall(
				r.namespace, generatorFunction)))

	newBody.SrcrefTreeIfMissing(originalBody)
	r.compiler.ReportChangeToEnclosingScope(newBody)
}
