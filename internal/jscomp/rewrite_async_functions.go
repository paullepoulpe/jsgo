// Port of com/google/javascript/jscomp/RewriteAsyncFunctions.java @ closure-compiler v20260415

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
	"iter"

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
//	  let $jscomp$async$super$get$x = () => super.x;
//	  return $jscomp.asyncExecutePromiseGeneratorFunction(
//	      function* () {
//	        // original body of foo() with:
//	        // - await (x) replaced with yield (x)
//	        // - arguments replaced with $jscomp$async$arguments
//	        // - this replaced with $jscomp$async$this
//	        // - super.x replaced with $jscomp$async$super$get$x()
//	        // - super.x(5) replaced with $jscomp$async$super$get$x().call($jscomp$async$this, 5)
//	      });
//	}
type RewriteAsyncFunctions struct {
	namespace rhino.StaticScope

	// Go-only: language-adapter — the ArrayDeque used as a stack is a slice
	// whose last element is the top (Deque.peek/push/pop).
	contextStack []RewriteAsyncFunctions_LexicalContext
	compiler     AbstractCompiler

	astFactory *AstFactory
}

const (
	// Java: RewriteAsyncFunctions.ASYNC_ARGUMENTS
	RewriteAsyncFunctions_ASYNC_ARGUMENTS = "$jscomp$async$arguments$"
	// Java: RewriteAsyncFunctions.ASYNC_THIS
	RewriteAsyncFunctions_ASYNC_THIS = "$jscomp$async$this$"
	// Java: RewriteAsyncFunctions.ASYNC_SUPER_PROP_GETTER_PREFIX
	RewriteAsyncFunctions_ASYNC_SUPER_PROP_GETTER_PREFIX = "$jscomp$async$super$get$"
)

// RewriteAsyncFunctions_SuperPropertyWrapperInfo is the information needed
// to replace a reference to `super.propertyName` with a call to a wrapper
// function.
type RewriteAsyncFunctions_SuperPropertyWrapperInfo struct {
	outer *RewriteAsyncFunctions
	// The first `super.property` Node we come across during traversal.
	// Information from this node will be used when creating the wrapper function.
	firstInstanceOfSuperDotProperty *rhino.Node
	wrapperFunctionName             string
	// The type to use for the wrapper function.
	// Will be null if type checking has not run.
	wrapperFunctionReturnType AstFactory_Type
}

// Java: SuperPropertyWrapperInfo(Node, String, AstFactory.Type)
func newRewriteAsyncFunctions_SuperPropertyWrapperInfo(
	outer *RewriteAsyncFunctions,
	firstSuperDotPropertyNode *rhino.Node,
	wrapperFunctionName string,
	wrapperFunctionReturnType AstFactory_Type) *RewriteAsyncFunctions_SuperPropertyWrapperInfo {
	i := &RewriteAsyncFunctions_SuperPropertyWrapperInfo{outer: outer}
	i.firstInstanceOfSuperDotProperty = firstSuperDotPropertyNode
	i.wrapperFunctionName = wrapperFunctionName
	i.wrapperFunctionReturnType = wrapperFunctionReturnType
	return i
}

// Java: SuperPropertyWrapperInfo.getPropertyType()
func (i *RewriteAsyncFunctions_SuperPropertyWrapperInfo) GetPropertyType() *colors.Color {
	return i.firstInstanceOfSuperDotProperty.GetColor()
}

// Java: SuperPropertyWrapperInfo.createWrapperFunctionNameNode()
func (i *RewriteAsyncFunctions_SuperPropertyWrapperInfo) CreateWrapperFunctionNameNode() *rhino.Node {
	return i.outer.astFactory.CreateName(i.wrapperFunctionName, AstFactory_Type_OfColor(colors.StandardColors_TOP_OBJECT))
}

// Java: SuperPropertyWrapperInfo.createWrapperFunctionCallNode()
func (i *RewriteAsyncFunctions_SuperPropertyWrapperInfo) CreateWrapperFunctionCallNode() *rhino.Node {
	return i.outer.astFactory.CreateCall(i.CreateWrapperFunctionNameNode(), i.wrapperFunctionReturnType)
}

// RewriteAsyncFunctions_SuperPropertyWrappers is used to collect information
// about properties referenced via `super.propertyName` within an async
// function.
//
// We'll have to replace these references with calls to a wrapper function.
type RewriteAsyncFunctions_SuperPropertyWrappers struct {
	outer *RewriteAsyncFunctions
	// Use LinkedHashMap in order to ensure the ordering is the same on every compile of the same
	// source code.
	// Note that the Colors will be null if type checking hasn't run.
	propertyNameToTypeMap *jutil.LinkedHashMap[string, *RewriteAsyncFunctions_SuperPropertyWrapperInfo]
}

// Java: SuperPropertyWrappers()
func newRewriteAsyncFunctions_SuperPropertyWrappers(outer *RewriteAsyncFunctions) *RewriteAsyncFunctions_SuperPropertyWrappers {
	w := &RewriteAsyncFunctions_SuperPropertyWrappers{outer: outer}
	w.propertyNameToTypeMap = jutil.NewLinkedHashMap[string, *RewriteAsyncFunctions_SuperPropertyWrapperInfo]()
	return w
}

// Java: SuperPropertyWrappers.getOrCreateSuperPropertyWrapperInfo(Node)
func (w *RewriteAsyncFunctions_SuperPropertyWrappers) GetOrCreateSuperPropertyWrapperInfo(
	superDotPropertyNode *rhino.Node) *RewriteAsyncFunctions_SuperPropertyWrapperInfo {
	jutil.CheckArgumentMsg1(superDotPropertyNode.IsGetProp(), "%s", superDotPropertyNode)
	superNode := superDotPropertyNode.GetFirstChild()
	jutil.CheckArgumentMsg1(superNode.IsSuper(), "%s", superNode)

	propertyName := superDotPropertyNode.GetString()
	propertyType := superDotPropertyNode.GetColor()
	var superPropertyWrapperInfo *RewriteAsyncFunctions_SuperPropertyWrapperInfo
	if w.propertyNameToTypeMap.ContainsKey(propertyName) {
		superPropertyWrapperInfo = w.propertyNameToTypeMap.Get(propertyName)
		// Every reference to `super.propertyName` within a single lexical context should
		// have the same type.  Make sure this is true.
		existingColor := superPropertyWrapperInfo.GetPropertyType()
		jutil.CheckStateMsg2(
			jutil.ObjectsEquals(existingColor, propertyType),
			"Previous reference type: %s differs from current reference type: %s",
			existingColor,
			propertyType)
	} else {
		superPropertyWrapperInfo = w.CreateNewInfo(superDotPropertyNode)
		w.propertyNameToTypeMap.Put(propertyName, superPropertyWrapperInfo)
	}
	return superPropertyWrapperInfo
}

// Java: SuperPropertyWrappers.createNewInfo(Node)
func (w *RewriteAsyncFunctions_SuperPropertyWrappers) CreateNewInfo(firstSuperDotPropertyNode *rhino.Node) *RewriteAsyncFunctions_SuperPropertyWrapperInfo {
	jutil.CheckArgumentMsg1(firstSuperDotPropertyNode.IsGetProp(), "%s", firstSuperDotPropertyNode)
	propertyName := firstSuperDotPropertyNode.GetString()
	wrapperFunctionName :=
		RewriteAsyncFunctions_ASYNC_SUPER_PROP_GETTER_PREFIX +
			w.outer.compiler.
				GetUniqueIdSupplier().
				GetUniqueId(w.outer.compiler.GetInput(NodeUtil_GetInputId(firstSuperDotPropertyNode))) +
			"$" +
			propertyName
	return newRewriteAsyncFunctions_SuperPropertyWrapperInfo(
		w.outer, firstSuperDotPropertyNode, wrapperFunctionName, AstFactory_Type_(firstSuperDotPropertyNode))
}

// Java: SuperPropertyWrappers.asCollection()
func (w *RewriteAsyncFunctions_SuperPropertyWrappers) AsCollection() iter.Seq[*RewriteAsyncFunctions_SuperPropertyWrapperInfo] {
	return w.propertyNameToTypeMap.Values()
}

// RewriteAsyncFunctions_LexicalContext determines both what to do when
// visiting a node and how to determine the context for its descendents.
type RewriteAsyncFunctions_LexicalContext interface {
	// GetRewriteAsyncFunctions_LexicalContextBase gives access to the fields
	// of the base class.
	GetRewriteAsyncFunctions_LexicalContextBase() *RewriteAsyncFunctions_LexicalContextBase

	// Java: LexicalContext.getContextRootNode()
	GetContextRootNode() *rhino.Node

	// GetContextForNode returns the LexicalContext to use for visiting a
	// node: n is this context's root node or one of its descendents; the
	// result is this context or a new one for a child context.
	//
	// Java: LexicalContext.getContextForNode(NodeTraversal, Node)
	GetContextForNode(t *NodeTraversal, n *rhino.Node) RewriteAsyncFunctions_LexicalContext

	// Java: LexicalContext.visit(NodeTraversal, Node)
	Visit(t *NodeTraversal, n *rhino.Node)
}

// RewriteAsyncFunctions_LexicalContextBase holds the fields of
// LexicalContext.
type RewriteAsyncFunctions_LexicalContextBase struct {
	contextRootNode *rhino.Node
	// precomputed unique id to append to names in this context. This is used to ensure that names
	// in different contexts don't collide (e.g. 2 functions don't get the same `let
	// $jscomp$async$this` name declared in their bodies)
	uniqueId string
}

// Java: LexicalContext(Node, String)
func (b *RewriteAsyncFunctions_LexicalContextBase) InitRewriteAsyncFunctions_LexicalContextBase(contextRootNode *rhino.Node, uniqueId string) {
	b.contextRootNode = jutil.CheckNotNull(contextRootNode)
	b.uniqueId = uniqueId
}

// GetRewriteAsyncFunctions_LexicalContextBase gives access to the fields of
// the base class.
//
// Go-only: language-adapter — field access through the interface.
func (b *RewriteAsyncFunctions_LexicalContextBase) GetRewriteAsyncFunctions_LexicalContextBase() *RewriteAsyncFunctions_LexicalContextBase {
	return b
}

// Java: LexicalContext.getContextRootNode()
func (b *RewriteAsyncFunctions_LexicalContextBase) GetContextRootNode() *rhino.Node {
	return b.contextRootNode
}

// RewriteAsyncFunctions_RootContext defines behavior for nodes in the root
// scope, outside of any functions.
type RewriteAsyncFunctions_RootContext struct {
	RewriteAsyncFunctions_LexicalContextBase
	outer *RewriteAsyncFunctions
}

// Java: RootContext(Node, String)
func newRewriteAsyncFunctions_RootContext(outer *RewriteAsyncFunctions, contextRootNode *rhino.Node, uniqueId string) *RewriteAsyncFunctions_RootContext {
	c := &RewriteAsyncFunctions_RootContext{outer: outer}
	c.InitRewriteAsyncFunctions_LexicalContextBase(contextRootNode, uniqueId)
	return c
}

// Java: RootContext.getContextForNode(NodeTraversal, Node)
func (c *RewriteAsyncFunctions_RootContext) GetContextForNode(t *NodeTraversal, n *rhino.Node) RewriteAsyncFunctions_LexicalContext {
	if n.IsFunction() {
		uniqueId := c.outer.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())
		return newRewriteAsyncFunctions_FunctionContext(c.outer, n, uniqueId)
	} else {
		return c
	}
}

// Java: RootContext.visit(NodeTraversal, Node)
func (c *RewriteAsyncFunctions_RootContext) Visit(t *NodeTraversal, n *rhino.Node) {
	// In root context we haven't entered an async function yet, so there's nothing to do.
}

// RewriteAsyncFunctions_ParameterListContext defines the behavior for
// function definition parameter lists and their contents.
type RewriteAsyncFunctions_ParameterListContext struct {
	RewriteAsyncFunctions_LexicalContextBase
	outer           *RewriteAsyncFunctions
	functionContext *RewriteAsyncFunctions_FunctionContext
}

// Java: ParameterListContext(FunctionContext, Node, String)
func newRewriteAsyncFunctions_ParameterListContext(
	outer *RewriteAsyncFunctions, functionContext *RewriteAsyncFunctions_FunctionContext, contextRootNode *rhino.Node, uniqueId string) *RewriteAsyncFunctions_ParameterListContext {
	c := &RewriteAsyncFunctions_ParameterListContext{outer: outer}
	c.InitRewriteAsyncFunctions_LexicalContextBase(contextRootNode, uniqueId)
	c.functionContext = jutil.CheckNotNull(functionContext)
	return c
}

// Java: ParameterListContext.getContextForNode(NodeTraversal, Node)
func (c *RewriteAsyncFunctions_ParameterListContext) GetContextForNode(t *NodeTraversal, n *rhino.Node) RewriteAsyncFunctions_LexicalContext {
	if n.IsFunction() {
		// Function defined within a parameter list.
		// e.g. `() => something`
		// function someFunc(callback = () => something) {}
		uniqueId := c.outer.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())
		return newRewriteAsyncFunctions_FunctionContextWithOuter(c.outer, c.functionContext, n, uniqueId)
	} else {
		return c
	}
}

// Java: ParameterListContext.visit(NodeTraversal, Node)
func (c *RewriteAsyncFunctions_ParameterListContext) Visit(t *NodeTraversal, n *rhino.Node) {
	if c.functionContext.asyncThisAndArgumentsContext != nil &&
		c.functionContext.asyncThisAndArgumentsContext != c.functionContext {
		// e.g.
		// async function outer(outerT = this) {
		//   // `this` in outer parameter list must remain unchanged
		//   // but inner parameter list must be aliased
		//   const inner = async (t = this) => t;
		// }
		c.functionContext.Visit(t, n)
	}
}

// RewriteAsyncFunctions_FunctionContext defines behavior for replacing
// references to `this`, `arguments`, and `super` within async functions and
// rewriting the async functions themselves.
type RewriteAsyncFunctions_FunctionContext struct {
	RewriteAsyncFunctions_LexicalContextBase
	outer *RewriteAsyncFunctions

	// If references to `this`, `arguments`, and `super` should be considered in the context of
	// an async function, this will point to that function's FunctionContext.
	// Otherwise, it will be `null`.
	// TODO(bradfordcsmith): It would cost less memory if we defined a separate object to hold
	// the data for async context accounting instead of having the booleans and super property
	// wrapper fields on every FunctionContext.
	asyncThisAndArgumentsContext *RewriteAsyncFunctions_FunctionContext
	superPropertyWrappers        *RewriteAsyncFunctions_SuperPropertyWrappers

	mustAddAsyncThisVariable bool
	// null if mustAddAsyncThisVariable is false
	typeOfThis                    AstFactory_Type
	mustAddAsyncArgumentsVariable bool
}

// Java: FunctionContext(Node, String)
func newRewriteAsyncFunctions_FunctionContext(outer *RewriteAsyncFunctions, contextRootNode *rhino.Node, uniqueId string) *RewriteAsyncFunctions_FunctionContext {
	c := &RewriteAsyncFunctions_FunctionContext{outer: outer}
	c.InitRewriteAsyncFunctions_LexicalContextBase(contextRootNode, uniqueId)
	c.superPropertyWrappers = newRewriteAsyncFunctions_SuperPropertyWrappers(outer)
	c.mustAddAsyncThisVariable = false
	c.mustAddAsyncArgumentsVariable = false
	if contextRootNode.IsAsyncFunction() {
		c.asyncThisAndArgumentsContext = c
	} else {
		c.asyncThisAndArgumentsContext = nil
	}
	return c
}

// Java: FunctionContext(FunctionContext, Node, String)
func newRewriteAsyncFunctions_FunctionContextWithOuter(
	outerInstance *RewriteAsyncFunctions, outer *RewriteAsyncFunctions_FunctionContext, contextRootNode *rhino.Node, uniqueId string) *RewriteAsyncFunctions_FunctionContext {
	c := &RewriteAsyncFunctions_FunctionContext{outer: outerInstance}
	c.InitRewriteAsyncFunctions_LexicalContextBase(contextRootNode, uniqueId)
	c.superPropertyWrappers = newRewriteAsyncFunctions_SuperPropertyWrappers(outerInstance)
	c.mustAddAsyncThisVariable = false
	c.mustAddAsyncArgumentsVariable = false
	jutil.CheckStateMsg1(contextRootNode.IsFunction(), "%s", contextRootNode)
	if contextRootNode.IsAsyncFunction() {
		if contextRootNode.IsArrowFunction() {
			// An async arrow function context points to outer.asyncThisAndArgumentsContext
			// if non-null, otherwise to itself.
			if outer.asyncThisAndArgumentsContext == nil {
				c.asyncThisAndArgumentsContext = c
			} else {
				c.asyncThisAndArgumentsContext = outer.asyncThisAndArgumentsContext
			}
		} else {
			// An async non-arrow function context always points to itself
			c.asyncThisAndArgumentsContext = c
		}
	} else {
		if contextRootNode.IsArrowFunction() {
			// A non-async arrow function context always points to
			// outer.asyncThisAndArgumentsContext
			c.asyncThisAndArgumentsContext = outer.asyncThisAndArgumentsContext
		} else {
			// A non-async, non-arrow function has no async context.
			c.asyncThisAndArgumentsContext = nil
		}
	}
	return c
}

// Java: FunctionContext.getContextForNode(NodeTraversal, Node)
func (c *RewriteAsyncFunctions_FunctionContext) GetContextForNode(t *NodeTraversal, n *rhino.Node) RewriteAsyncFunctions_LexicalContext {
	if n == c.contextRootNode {
		return c
	} else if n.IsFunction() {
		uniqueId := c.outer.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())
		return newRewriteAsyncFunctions_FunctionContextWithOuter(c.outer, c, n, uniqueId)
	} else if n.IsParamList() {
		uniqueId := c.outer.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())
		return newRewriteAsyncFunctions_ParameterListContext(c.outer, c, n, uniqueId)
	} else {
		return c
	}
}

// Java: FunctionContext.recordAsyncThisReplacementWasDone(AstFactory.Type)
func (c *RewriteAsyncFunctions_FunctionContext) RecordAsyncThisReplacementWasDone(typeOfThis AstFactory_Type) {
	c.asyncThisAndArgumentsContext.mustAddAsyncThisVariable = true
	c.asyncThisAndArgumentsContext.typeOfThis = typeOfThis
}

// Java: FunctionContext.getOrCreateSuperPropertyWrapperInfo(Node)
func (c *RewriteAsyncFunctions_FunctionContext) GetOrCreateSuperPropertyWrapperInfo(
	superDotPropertyNode *rhino.Node) *RewriteAsyncFunctions_SuperPropertyWrapperInfo {
	return c.asyncThisAndArgumentsContext.superPropertyWrappers.GetOrCreateSuperPropertyWrapperInfo(
		superDotPropertyNode)
}

// Java: FunctionContext.recordAsyncArgumentsReplacementWasDone()
func (c *RewriteAsyncFunctions_FunctionContext) RecordAsyncArgumentsReplacementWasDone() {
	c.asyncThisAndArgumentsContext.mustAddAsyncArgumentsVariable = true
}

// CreateThisVariableReference creates a new reference to the variable used to
// hold the value of `this` for async functions.
//
// Java: FunctionContext.createThisVariableReference(AstFactory.Type)
func (c *RewriteAsyncFunctions_FunctionContext) CreateThisVariableReference(typeOfThis AstFactory_Type) *rhino.Node {
	c.RecordAsyncThisReplacementWasDone(typeOfThis)
	return c.outer.astFactory.CreateName(RewriteAsyncFunctions_ASYNC_THIS+c.asyncThisAndArgumentsContext.uniqueId, typeOfThis)
}

// Java: FunctionContext.createWrapperArrowFunction(SuperPropertyWrapperInfo)
func (c *RewriteAsyncFunctions_FunctionContext) CreateWrapperArrowFunction(wrapperInfo *RewriteAsyncFunctions_SuperPropertyWrapperInfo) *rhino.Node {
	// super.propertyName
	superDotProperty := wrapperInfo.firstInstanceOfSuperDotProperty.CloneTree()

	// () => { return super.propertyName; };
	return c.outer.astFactory.CreateZeroArgArrowFunctionForExpression(
		c.outer.astFactory.CreateBlock(c.outer.astFactory.CreateReturn(superDotProperty)))
}

// Java: FunctionContext.visit(NodeTraversal, Node)
func (c *RewriteAsyncFunctions_FunctionContext) Visit(t *NodeTraversal, n *rhino.Node) {
	r := c.outer
	if c.contextRootNode == n && c.contextRootNode.IsAsyncFunction() {
		// We're visiting an async function.
		// All of its descendent nodes will have been updated as necessary, so now we just need to
		// convert the function itself.
		r.ConvertAsyncFunction(t, c)
	} else if c.asyncThisAndArgumentsContext != nil {
		// We're in the context of an async function's body, so we need to do some replacements.
		switch n.GetToken() {
		case rhino.Token_NAME:
			if n.MatchesName("arguments") {
				n.SetString(RewriteAsyncFunctions_ASYNC_ARGUMENTS + c.asyncThisAndArgumentsContext.uniqueId)
				if r.compiler.GetLifeCycleStage().IsNormalized() {
					// TODO: b/322009741 - Stop depending on lifeCycleStage.isNormalized() to decide
					// constness
					n.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
				}
				c.asyncThisAndArgumentsContext.RecordAsyncArgumentsReplacementWasDone()
				r.compiler.ReportChangeToChangeScope(c.contextRootNode)
			}
		case rhino.Token_THIS:
			n.ReplaceWith(c.asyncThisAndArgumentsContext.CreateThisVariableReference(AstFactory_Type_(n)))
			r.compiler.ReportChangeToChangeScope(c.contextRootNode)
		case rhino.Token_SUPER:
			{
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
				// different name for parent for better readability
				superDotProperty := parent

				superPropertyWrapperInfo :=
					c.asyncThisAndArgumentsContext.GetOrCreateSuperPropertyWrapperInfo(
						superDotProperty)

				// super.x   =>   $jscomp$super$get$x()
				getPropReplacement := superPropertyWrapperInfo.CreateWrapperFunctionCallNode()
				grandparent := superDotProperty.GetParent()
				if grandparent.IsCall() && grandparent.GetFirstChild() == superDotProperty {
					// $jscomp$super$get$x(...)   =>   $jscomp$super$get$x().call($jscomp$async$this,
					// ...)
					getPropReplacement =
						r.astFactory.CreateGetPropWithUnknownType(getPropReplacement, "call")
					thisAlias :=
						r.astFactory.
							CreateThisAliasReferenceForEs6Class(
								RewriteAsyncFunctions_ASYNC_THIS+c.asyncThisAndArgumentsContext.uniqueId,
								NodeUtil_GetEnclosingClass(
									c.asyncThisAndArgumentsContext.GetContextRootNode())).
							Srcref(superDotProperty)
					thisAlias.InsertAfter(superDotProperty)
					c.asyncThisAndArgumentsContext.RecordAsyncThisReplacementWasDone(AstFactory_Type_(thisAlias))
				}
				getPropReplacement.SrcrefTree(superDotProperty)
				superDotProperty.ReplaceWith(getPropReplacement)
				r.compiler.ReportChangeToChangeScope(c.contextRootNode)
			}
		case rhino.Token_AWAIT:
			// Awaits become yields in the converted async function's inner generator function.
			n.ReplaceWith(r.astFactory.CreateYield(AstFactory_Type_(n), n.RemoveFirstChild()))
		default:
		}
	}
}

// Java: RewriteAsyncFunctions.transpiledFeatures
var RewriteAsyncFunctions_transpiledFeatures = parser.FeatureSet_BARE_MINIMUM.With(parser.FeatureSet_Feature_ASYNC_FUNCTIONS)

// Java: RewriteAsyncFunctions(AbstractCompiler, AstFactory, StaticScope)
func newRewriteAsyncFunctions(
	compiler AbstractCompiler, astFactory *AstFactory, namespace rhino.StaticScope) *RewriteAsyncFunctions {
	r := &RewriteAsyncFunctions{}
	r.compiler = jutil.CheckNotNull(compiler)
	r.astFactory = jutil.CheckNotNull(astFactory)
	r.namespace = jutil.CheckNotNull(namespace)
	r.contextStack = nil
	return r
}

// Java: RewriteAsyncFunctions.create(AbstractCompiler)
func RewriteAsyncFunctions_Create(compiler AbstractCompiler) *RewriteAsyncFunctions {
	astFactory := compiler.CreateAstFactory()
	namespace := compiler.GetTranspilationNamespace()
	return newRewriteAsyncFunctions(compiler, astFactory, namespace)
}

// Java: RewriteAsyncFunctions.process(Node, Node)
func (r *RewriteAsyncFunctions) Process(externs *rhino.Node, root *rhino.Node) {
	TranspilationPasses_ProcessTranspile(r.compiler, root, RewriteAsyncFunctions_transpiledFeatures, r)
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(r.compiler, root, RewriteAsyncFunctions_transpiledFeatures)
}

// contextStackPeek is Deque.peek(): the top, or null when empty.
//
// Go-only: language-adapter — the ArrayDeque operation on the slice stack.
func (r *RewriteAsyncFunctions) contextStackPeek() RewriteAsyncFunctions_LexicalContext {
	if len(r.contextStack) == 0 {
		return nil
	}
	return r.contextStack[len(r.contextStack)-1]
}

// Java: RewriteAsyncFunctions.shouldTraverse(NodeTraversal, Node, Node)
func (r *RewriteAsyncFunctions) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	if parent == nil {
		jutil.CheckState(len(r.contextStack) == 0)
		uniqueId := r.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())
		r.contextStack = append(r.contextStack, newRewriteAsyncFunctions_RootContext(r, n, uniqueId))
	} else {
		parentContext := r.contextStackPeek()
		nodeContext := parentContext.GetContextForNode(t, n)
		if nodeContext != parentContext {
			r.contextStack = append(r.contextStack, nodeContext)
		}
	}
	return true
}

// Java: RewriteAsyncFunctions.visit(NodeTraversal, Node, Node)
func (r *RewriteAsyncFunctions) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	context := r.contextStackPeek()
	context.Visit(t, n)
	if context.GetContextRootNode() == n {
		if len(r.contextStack) == 0 {
			panic(jutil.NewNoSuchElementException()) // ArrayDeque.pop()
		}
		r.contextStack = r.contextStack[:len(r.contextStack)-1]
	}
}

// Java: RewriteAsyncFunctions.convertAsyncFunction(NodeTraversal, FunctionContext)
func (r *RewriteAsyncFunctions) ConvertAsyncFunction(t *NodeTraversal, functionContext *RewriteAsyncFunctions_FunctionContext) {
	originalFunction := functionContext.GetContextRootNode()
	originalFunction.SetIsAsyncFunction(false)
	originalBody := originalFunction.GetLastChild()
	if originalFunction.IsFromExterns() {
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

	if functionContext.mustAddAsyncThisVariable {
		// const this$ = this;
		newBody.AddChildToBack(
			r.astFactory.CreateSingleConstNameDeclaration(
				RewriteAsyncFunctions_ASYNC_THIS+functionContext.asyncThisAndArgumentsContext.uniqueId,
				r.astFactory.CreateThis(functionContext.typeOfThis)))
		NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_CONST_DECLARATIONS, r.compiler)
	}
	if functionContext.mustAddAsyncArgumentsVariable {
		// const arguments$ = arguments;
		newBody.AddChildToBack(
			r.astFactory.CreateArgumentsAliasDeclaration(
				RewriteAsyncFunctions_ASYNC_ARGUMENTS + functionContext.asyncThisAndArgumentsContext.uniqueId))
		NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_CONST_DECLARATIONS, r.compiler)
	}
	for superPropertyWrapperInfo := range functionContext.superPropertyWrappers.AsCollection() {
		arrowFunction := functionContext.CreateWrapperArrowFunction(superPropertyWrapperInfo)
		// const super$get$x = () => super.x;
		arrowFunctionDeclarationStatement :=
			r.astFactory.CreateSingleConstNameDeclaration(
				superPropertyWrapperInfo.wrapperFunctionName, arrowFunction)
		newBody.AddChildToBack(arrowFunctionDeclarationStatement)

		// Make sure the compiler knows about the new arrow function's scope
		r.compiler.ReportChangeToChangeScope(arrowFunction)
		// Record that we've added arrow functions and const declarations to this script,
		// so later transpilations of those features will run, if needed.
		enclosingScript := t.GetCurrentScript()
		NodeUtil_AddFeatureToScript(enclosingScript, parser.FeatureSet_Feature_ARROW_FUNCTIONS, r.compiler)
		NodeUtil_AddFeatureToScript(enclosingScript, parser.FeatureSet_Feature_CONST_DECLARATIONS, r.compiler)
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
