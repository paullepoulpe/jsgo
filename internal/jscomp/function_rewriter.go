// Port of com/google/javascript/jscomp/FunctionRewriter.java @ closure-compiler v20260415

/*
 * Copyright 2009 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// FunctionRewriter reduces the size of common function expressions.
//
// This pass will rewrite:
//
//	C.prototype.getA = function() { return this.a_ };
//	C.prototype.setA = function(newValue) { this.a_ = newValue };
//
// as:
//
//	C.prototype.getA = JSCompiler_get("a");
//	C.prototype.setA = JSCompiler_set("a");
//
// if by doing so we will save bytes, after the helper functions are added
// and renaming is done.
//
// NOTE: JSCompiler_get and JSCompiler_set turn dotted accesses to computed
// accesses, which causes JS engines to use dictionary lookups. Because of
// this perf regression, this pass is off by default in advanced mode even
// though it improves code size.
type FunctionRewriter struct {
	compiler AbstractCompiler
}

// Safety margin used to avoid growing simple programs by a few bytes.
// Selected arbitrarily.
const FunctionRewriter_SAVINGS_THRESHOLD = 16

// Java: FunctionRewriter(AbstractCompiler)
func NewFunctionRewriter(compiler AbstractCompiler) *FunctionRewriter {
	return &FunctionRewriter{compiler: compiler}
}

// Java: FunctionRewriter.process(Node, Node)
func (f *FunctionRewriter) Process(externs *rhino.Node, root *rhino.Node) {
	reducers := []FunctionRewriter_Reducer{
		NewFunctionRewriter_ReturnConstantReducer(),
		NewFunctionRewriter_GetterReducer(),
		NewFunctionRewriter_SetterReducer(),
		NewFunctionRewriter_EmptyFunctionReducer(),
		NewFunctionRewriter_IdentityReducer(),
	}

	// Go-only: language-adapter — a HashMultimap keyed by identity-hashed
	// Reducers with identity-hashed Reduction values; Java's iteration
	// order there is not deterministic, so insertion order is used.
	reductionMap := jutil.NewSetMultimap[FunctionRewriter_Reducer, *FunctionRewriter_Reduction]()

	// Accumulate possible reductions in the reduction multi-map.  They
	// will be applied in the loop below.
	NodeTraversal_Traverse(f.compiler, root,
		NewFunctionRewriter_ReductionGatherer(f, reducers, reductionMap))

	// Apply reductions iff they will provide some savings.
	for _, reducer := range reducers {
		reductions := reductionMap.Get(reducer)
		if reductions.IsEmpty() {
			continue
		}

		helperCode := f.ParseHelperCode(reducer)
		if helperCode == nil {
			continue
		}

		helperCodeCost := InlineCostEstimator_GetCost(helperCode)

		// Estimate savings
		savings := 0
		for reduction := range reductions.All() {
			savings += reduction.EstimateSavings()
		}

		// Compare estimated savings against the helper cost.  Apply
		// reductions if doing so will result in some savings.
		if savings > (helperCodeCost + FunctionRewriter_SAVINGS_THRESHOLD) {
			for reduction := range reductions.All() {
				reduction.Apply()
			}

			addingRoot := f.compiler.GetNodeForCodeInsertion(nil)
			NodeUtil_MarkNewScopesChanged(helperCode, f.compiler)
			addingRoot.AddChildToFront(helperCode)
			f.compiler.ReportChangeToEnclosingScope(addingRoot)
		}
	}
}

// ParseHelperCode parses helper code needed by a reducer.
//
// Returns the helper code root. If parse fails, return null.
//
// Java: FunctionRewriter.parseHelperCode(Reducer)
func (f *FunctionRewriter) ParseHelperCode(reducer FunctionRewriter_Reducer) *rhino.Node {
	root :=
		f.compiler.ParseSyntheticCode(functionRewriter_classString(reducer)+":helper", reducer.GetHelperSource())
	if root != nil {
		return root.RemoveFirstChild()
	}
	return nil
}

// functionRewriter_classString is `reducer.getClass()` in a string
// concatenation (Class.toString()).
//
// Go-only: language-adapter — Go has no Class objects; the names are the
// Java binary names of the reducer classes.
func functionRewriter_classString(reducer FunctionRewriter_Reducer) string {
	switch reducer.(type) {
	case *FunctionRewriter_EmptyFunctionReducer:
		return "class com.google.javascript.jscomp.FunctionRewriter$EmptyFunctionReducer"
	case *FunctionRewriter_IdentityReducer:
		return "class com.google.javascript.jscomp.FunctionRewriter$IdentityReducer"
	case *FunctionRewriter_ReturnConstantReducer:
		return "class com.google.javascript.jscomp.FunctionRewriter$ReturnConstantReducer"
	case *FunctionRewriter_GetterReducer:
		return "class com.google.javascript.jscomp.FunctionRewriter$GetterReducer"
	case *FunctionRewriter_SetterReducer:
		return "class com.google.javascript.jscomp.FunctionRewriter$SetterReducer"
	}
	panic(jutil.NewIllegalStateExceptionMsg("unknown reducer"))
}

// Java: FunctionRewriter.isReduceableFunctionExpression(Node)
func FunctionRewriter_IsReduceableFunctionExpression(n *rhino.Node) bool {
	parent := n.GetParent()
	return NodeUtil_IsFunctionExpression(n) &&
		!NodeUtil_IsGetOrSetKey(parent) &&
		!parent.IsMemberFunctionDef() &&
		!n.IsAsyncFunction() &&
		!n.IsGeneratorFunction()
}

// FunctionRewriter_Reduction is information needed to apply a reduction.
type FunctionRewriter_Reduction struct {
	outer    *FunctionRewriter
	oldChild *rhino.Node
	newChild *rhino.Node
}

// Java: FunctionRewriter.Reduction(Node, Node, Node)
func NewFunctionRewriter_Reduction(outer *FunctionRewriter, parent *rhino.Node, oldChild *rhino.Node, newChild *rhino.Node) *FunctionRewriter_Reduction {
	return &FunctionRewriter_Reduction{outer: outer, oldChild: oldChild, newChild: newChild}
}

// Apply applies the reduction by replacing the old child with the new
// child.
//
// Java: FunctionRewriter.Reduction.apply()
func (r *FunctionRewriter_Reduction) Apply() {
	r.oldChild.ReplaceWith(r.newChild)
	NodeUtil_MarkFunctionsDeleted(r.oldChild, r.outer.compiler)
	r.outer.compiler.ReportChangeToEnclosingScope(r.newChild)
}

// EstimateSavings estimates the number of bytes saved by applying this
// reduction.
//
// Java: FunctionRewriter.Reduction.estimateSavings()
func (r *FunctionRewriter_Reduction) EstimateSavings() int {
	return InlineCostEstimator_GetCost(r.oldChild) -
		InlineCostEstimator_GetCost(r.newChild)
}

// FunctionRewriter_ReductionGatherer gathers a list of reductions to apply
// later by doing an in-order AST traversal. If a suitable reduction is
// found, stop traversal in that branch.
type FunctionRewriter_ReductionGatherer struct {
	outer      *FunctionRewriter
	reducers   []FunctionRewriter_Reducer
	reductions *jutil.SetMultimap[FunctionRewriter_Reducer, *FunctionRewriter_Reduction]
}

// NewFunctionRewriter_ReductionGatherer: reducers is the list of reducers
// to apply during traversal; reductions the Reducer -> Reduction multimap,
// populated during traversal.
//
// Java: FunctionRewriter.ReductionGatherer(List<Reducer>, Multimap<Reducer, Reduction>)
func NewFunctionRewriter_ReductionGatherer(outer *FunctionRewriter, reducers []FunctionRewriter_Reducer,
	reductions *jutil.SetMultimap[FunctionRewriter_Reducer, *FunctionRewriter_Reduction]) *FunctionRewriter_ReductionGatherer {
	// Use a native array to avoid creating an iterator for every node in the AST.
	// which accord to pprof accounts for ~20% of the runtime of this pass.
	return &FunctionRewriter_ReductionGatherer{
		outer:      outer,
		reducers:   append([]FunctionRewriter_Reducer(nil), reducers...),
		reductions: reductions,
	}
}

// Java: FunctionRewriter.ReductionGatherer.shouldTraverse(NodeTraversal, Node, Node)
func (g *FunctionRewriter_ReductionGatherer) ShouldTraverse(raversal *NodeTraversal,
	node *rhino.Node,
	parent *rhino.Node) bool {
	for _, reducer := range g.reducers {
		replacement := reducer.Reduce(node)
		if replacement != node {
			g.reductions.Put(reducer, NewFunctionRewriter_Reduction(g.outer, parent, node, replacement))
			return false
		}
	}
	return true
}

// Java: FunctionRewriter.ReductionGatherer.visit(NodeTraversal, Node, Node)
func (g *FunctionRewriter_ReductionGatherer) Visit(traversal *NodeTraversal, node *rhino.Node, parent *rhino.Node) {
}

// FunctionRewriter_Reducer is the interface implemented by the
// strength-reduction optimizers below.
type FunctionRewriter_Reducer interface {
	// GetHelperSource returns the JS source for helper methods used by this
	// reduction.
	//
	// Java: FunctionRewriter.Reducer.getHelperSource()
	GetHelperSource() string

	// Reduce returns the root of the reduced subtree if a reduction was
	// applied; otherwise returns the node argument.
	//
	// Java: FunctionRewriter.Reducer.reduce(Node)
	Reduce(node *rhino.Node) *rhino.Node

	GetFunctionRewriter_ReducerBase() *FunctionRewriter_ReducerBase
}

// FunctionRewriter_ReducerBase holds Reducer's final method.
type FunctionRewriter_ReducerBase struct {
	// Go-only: language-adapter — distinct identity (reducers are multimap
	// keys).
	_ byte
}

// GetFunctionRewriter_ReducerBase gives access to the base through the
// interface.
func (b *FunctionRewriter_ReducerBase) GetFunctionRewriter_ReducerBase() *FunctionRewriter_ReducerBase {
	return b
}

// BuildCallNode builds a method call based on the the given method name,
// argument and history.
//
// methodName is the method to call; argumentNode the method argument.
//
// Java: FunctionRewriter.Reducer.buildCallNode(String, Node, Node)
func (b *FunctionRewriter_ReducerBase) BuildCallNode(
	methodName string, argumentNode *rhino.Node, srcref *rhino.Node) *rhino.Node {
	call := rhino.IR_Call(rhino.IR_Name(methodName)).SrcrefTree(srcref)
	call.PutBooleanProp(rhino.Node_FREE_CALL, true)
	if argumentNode != nil {
		call.AddChildToBack(argumentNode.CloneTree())
	}
	return call
}

// FunctionRewriter_EmptyFunctionReducer reduces return immutable constant
// literal methods declarations with calls to a constant return method
// factory.
//
// Example:
//
//	a.prototype.b = function() {}
//
// is reduced to:
//
//	a.prototype.b = emptyFn();
type FunctionRewriter_EmptyFunctionReducer struct {
	FunctionRewriter_ReducerBase
}

const FunctionRewriter_EmptyFunctionReducer_FACTORY_METHOD_NAME = "JSCompiler_emptyFn"

const FunctionRewriter_EmptyFunctionReducer_HELPER_SOURCE = "function " + FunctionRewriter_EmptyFunctionReducer_FACTORY_METHOD_NAME + "() {" +
	"  return function() {}" +
	"}"

// Java: FunctionRewriter.EmptyFunctionReducer()
func NewFunctionRewriter_EmptyFunctionReducer() *FunctionRewriter_EmptyFunctionReducer {
	return &FunctionRewriter_EmptyFunctionReducer{}
}

// Java: FunctionRewriter.EmptyFunctionReducer.getHelperSource()
func (r *FunctionRewriter_EmptyFunctionReducer) GetHelperSource() string {
	return FunctionRewriter_EmptyFunctionReducer_HELPER_SOURCE
}

// Java: FunctionRewriter.EmptyFunctionReducer.reduce(Node)
func (r *FunctionRewriter_EmptyFunctionReducer) Reduce(node *rhino.Node) *rhino.Node {
	if FunctionRewriter_IsReduceableFunctionExpression(node) && NodeUtil_IsEmptyFunctionExpression(node) {
		return r.BuildCallNode(FunctionRewriter_EmptyFunctionReducer_FACTORY_METHOD_NAME, nil, node)
	} else {
		return node
	}
}

// FunctionRewriter_SingleReturnStatementReducer is the base class for
// reducers that match functions that contain a single return statement.
type FunctionRewriter_SingleReturnStatementReducer interface {
	FunctionRewriter_Reducer

	GetFunctionRewriter_SingleReturnStatementReducerBase() *FunctionRewriter_SingleReturnStatementReducerBase
}

// FunctionRewriter_SingleReturnStatementReducerBase holds
// SingleReturnStatementReducer's final method.
type FunctionRewriter_SingleReturnStatementReducerBase struct {
	FunctionRewriter_ReducerBase
}

// GetFunctionRewriter_SingleReturnStatementReducerBase gives access to the
// base through the interface.
func (b *FunctionRewriter_SingleReturnStatementReducerBase) GetFunctionRewriter_SingleReturnStatementReducerBase() *FunctionRewriter_SingleReturnStatementReducerBase {
	return b
}

// MaybeGetSingleReturnRValue returns the function return value node if
// function body contains a single return statement. Otherwise, null.
//
// Java: FunctionRewriter.SingleReturnStatementReducer.maybeGetSingleReturnRValue(Node)
func (b *FunctionRewriter_SingleReturnStatementReducerBase) MaybeGetSingleReturnRValue(functionNode *rhino.Node) *rhino.Node {
	body := functionNode.GetLastChild()
	if !body.HasOneChild() {
		return nil
	}

	statement := body.GetFirstChild()
	if statement.IsReturn() {
		return statement.GetFirstChild()
	}
	return nil
}

// FunctionRewriter_IdentityReducer reduces property getter method
// declarations with calls to a getter method factory.
//
// Example:
//
//	a.prototype.b = function(a) {return a}
//
// is reduced to:
//
//	a.prototype.b = getter(a);
type FunctionRewriter_IdentityReducer struct {
	FunctionRewriter_SingleReturnStatementReducerBase
}

const FunctionRewriter_IdentityReducer_FACTORY_METHOD_NAME = "JSCompiler_identityFn"

const FunctionRewriter_IdentityReducer_HELPER_SOURCE = "function " + FunctionRewriter_IdentityReducer_FACTORY_METHOD_NAME + "() {" +
	"  return function(" + FunctionRewriter_IdentityReducer_FACTORY_METHOD_NAME + "_value) {" +
	"return " + FunctionRewriter_IdentityReducer_FACTORY_METHOD_NAME + "_value}" +
	"}"

// Java: FunctionRewriter.IdentityReducer()
func NewFunctionRewriter_IdentityReducer() *FunctionRewriter_IdentityReducer {
	return &FunctionRewriter_IdentityReducer{}
}

// Java: FunctionRewriter.IdentityReducer.getHelperSource()
func (r *FunctionRewriter_IdentityReducer) GetHelperSource() string {
	return FunctionRewriter_IdentityReducer_HELPER_SOURCE
}

// Java: FunctionRewriter.IdentityReducer.reduce(Node)
func (r *FunctionRewriter_IdentityReducer) Reduce(node *rhino.Node) *rhino.Node {
	if !FunctionRewriter_IsReduceableFunctionExpression(node) {
		return node
	}

	if r.IsIdentityFunction(node) {
		return r.BuildCallNode(FunctionRewriter_IdentityReducer_FACTORY_METHOD_NAME, nil, node)
	} else {
		return node
	}
}

// IsIdentityFunction checks if the function matches the pattern:
//
//	function(<value>, <rest>) {return <value>}
//
// Returns whether the function matches the pattern.
//
// Java: FunctionRewriter.IdentityReducer.isIdentityFunction(Node)
func (r *FunctionRewriter_IdentityReducer) IsIdentityFunction(functionNode *rhino.Node) bool {
	argList := functionNode.GetSecondChild()
	paramNode := argList.GetFirstChild()
	if paramNode == nil || !paramNode.IsName() {
		// no parameters, or first parameter uses destructuring or a default value
		return false
	}

	value := r.MaybeGetSingleReturnRValue(functionNode)
	return value != nil && value.IsName() && value.GetString() == paramNode.GetString()
}

// FunctionRewriter_ReturnConstantReducer reduces return immutable constant
// literal methods declarations with calls to a constant return method
// factory.
//
// Example:
//
//	a.prototype.b = function() {return 10}
//
// is reduced to:
//
//	a.prototype.b = returnconst(10);
type FunctionRewriter_ReturnConstantReducer struct {
	FunctionRewriter_SingleReturnStatementReducerBase
}

const FunctionRewriter_ReturnConstantReducer_FACTORY_METHOD_NAME = "JSCompiler_returnArg"

const FunctionRewriter_ReturnConstantReducer_HELPER_SOURCE = "function " + FunctionRewriter_ReturnConstantReducer_FACTORY_METHOD_NAME +
	"(" + FunctionRewriter_ReturnConstantReducer_FACTORY_METHOD_NAME + "_value) {" +
	"  return function() {return " + FunctionRewriter_ReturnConstantReducer_FACTORY_METHOD_NAME + "_value}" +
	"}"

// Java: FunctionRewriter.ReturnConstantReducer()
func NewFunctionRewriter_ReturnConstantReducer() *FunctionRewriter_ReturnConstantReducer {
	return &FunctionRewriter_ReturnConstantReducer{}
}

// Java: FunctionRewriter.ReturnConstantReducer.getHelperSource()
func (r *FunctionRewriter_ReturnConstantReducer) GetHelperSource() string {
	return FunctionRewriter_ReturnConstantReducer_HELPER_SOURCE
}

// Java: FunctionRewriter.ReturnConstantReducer.reduce(Node)
func (r *FunctionRewriter_ReturnConstantReducer) Reduce(node *rhino.Node) *rhino.Node {
	if !FunctionRewriter_IsReduceableFunctionExpression(node) {
		return node
	}

	valueNode := r.GetValueNode(node)
	if valueNode != nil {
		return r.BuildCallNode(FunctionRewriter_ReturnConstantReducer_FACTORY_METHOD_NAME, valueNode, node)
	} else {
		return node
	}
}

// GetValueNode checks if the function matches the pattern: function(<args>)
// {return <immutable value>} and returns <immutable value> if a match is
// found.
//
// Returns the immutable value node; or null.
//
// Java: FunctionRewriter.ReturnConstantReducer.getValueNode(Node)
func (r *FunctionRewriter_ReturnConstantReducer) GetValueNode(functionNode *rhino.Node) *rhino.Node {
	value := r.MaybeGetSingleReturnRValue(functionNode)
	if value != nil &&
		NodeUtil_IsImmutableValue(value) {
		return value
	}
	return nil
}

// FunctionRewriter_GetterReducer reduces property getter method
// declarations with calls to a getter method factory.
//
// Example:
//
//	a.prototype.b = function() {return this.b_}
//
// is reduced to:
//
//	a.prototype.b = getter("b_");
type FunctionRewriter_GetterReducer struct {
	FunctionRewriter_SingleReturnStatementReducerBase
}

const FunctionRewriter_GetterReducer_FACTORY_METHOD_NAME = "JSCompiler_get"

const FunctionRewriter_GetterReducer_HELPER_SOURCE = "function " + FunctionRewriter_GetterReducer_FACTORY_METHOD_NAME + "(" +
	FunctionRewriter_GetterReducer_FACTORY_METHOD_NAME + "_name) {" +
	"  return function() {return this[" + FunctionRewriter_GetterReducer_FACTORY_METHOD_NAME + "_name]}" +
	"}"

// Java: FunctionRewriter.GetterReducer()
func NewFunctionRewriter_GetterReducer() *FunctionRewriter_GetterReducer {
	return &FunctionRewriter_GetterReducer{}
}

// Java: FunctionRewriter.GetterReducer.getHelperSource()
func (r *FunctionRewriter_GetterReducer) GetHelperSource() string {
	return FunctionRewriter_GetterReducer_HELPER_SOURCE
}

// Java: FunctionRewriter.GetterReducer.reduce(Node)
func (r *FunctionRewriter_GetterReducer) Reduce(node *rhino.Node) *rhino.Node {
	if !FunctionRewriter_IsReduceableFunctionExpression(node) || node.IsArrowFunction() {
		return node
	}

	getprop := r.GetGetPropertyName(node)
	if getprop != nil {
		if !getprop.IsGetProp() {
			panic(jutil.NewIllegalStateExceptionMsg("Expected GETPROP, got " + getprop.GetToken().String()))
		}

		arg := rhino.IR_String(getprop.GetString()).Srcref(getprop)
		return r.BuildCallNode(FunctionRewriter_GetterReducer_FACTORY_METHOD_NAME, arg, node)
	} else {
		return node
	}
}

// GetGetPropertyName checks if the function matches the pattern:
// function(<args>) {return this.<name>} and returns <name> if a match is
// found.
//
// Returns the STRING node that is the RHS of a this property get; or null.
//
// Java: FunctionRewriter.GetterReducer.getGetPropertyName(Node)
func (r *FunctionRewriter_GetterReducer) GetGetPropertyName(functionNode *rhino.Node) *rhino.Node {
	value := r.MaybeGetSingleReturnRValue(functionNode)
	if value != nil &&
		value.IsGetProp() &&
		value.GetFirstChild().IsThis() {
		return value
	}
	return nil
}

// FunctionRewriter_SetterReducer reduces property setter method
// declarations with calls to a setter method factory.
//
// Example:
//
//	a.prototype.setB = function(value) {this.b_ = value}
//
// reduces to:
//
//	a.prototype.setB = getter("b_");
type FunctionRewriter_SetterReducer struct {
	FunctionRewriter_ReducerBase
}

const FunctionRewriter_SetterReducer_FACTORY_METHOD_NAME = "JSCompiler_set"

const FunctionRewriter_SetterReducer_HELPER_SOURCE = "function " + FunctionRewriter_SetterReducer_FACTORY_METHOD_NAME + "(" +
	FunctionRewriter_SetterReducer_FACTORY_METHOD_NAME + "_name) {" +
	"  return function(" + FunctionRewriter_SetterReducer_FACTORY_METHOD_NAME + "_value) {" +
	"this[" + FunctionRewriter_SetterReducer_FACTORY_METHOD_NAME + "_name] = " +
	FunctionRewriter_SetterReducer_FACTORY_METHOD_NAME + "_value}" +
	"}"

// Java: FunctionRewriter.SetterReducer()
func NewFunctionRewriter_SetterReducer() *FunctionRewriter_SetterReducer {
	return &FunctionRewriter_SetterReducer{}
}

// Java: FunctionRewriter.SetterReducer.getHelperSource()
func (r *FunctionRewriter_SetterReducer) GetHelperSource() string {
	return FunctionRewriter_SetterReducer_HELPER_SOURCE
}

// Java: FunctionRewriter.SetterReducer.reduce(Node)
func (r *FunctionRewriter_SetterReducer) Reduce(node *rhino.Node) *rhino.Node {
	if !FunctionRewriter_IsReduceableFunctionExpression(node) || node.IsArrowFunction() {
		return node
	}

	propName := r.GetSetPropertyName(node)
	if propName != nil {
		if !propName.IsStringLit() {
			panic(jutil.NewIllegalStateExceptionMsg("Expected STRING, got " + propName.GetToken().String()))
		}

		return r.BuildCallNode(FunctionRewriter_SetterReducer_FACTORY_METHOD_NAME, propName, node)
	} else {
		return node
	}
}

// GetSetPropertyName checks if the function matches the pattern:
// function(<value>, <rest>) {this.<name> = <value>} and returns <name> if a
// match is found.
//
// Returns the STRING node that is the RHS of a this property get; or null.
//
// Java: FunctionRewriter.SetterReducer.getSetPropertyName(Node)
func (r *FunctionRewriter_SetterReducer) GetSetPropertyName(functionNode *rhino.Node) *rhino.Node {
	body := functionNode.GetLastChild()
	if !body.HasOneChild() {
		return nil
	}

	argList := functionNode.GetSecondChild()
	paramNode := argList.GetFirstChild()
	if paramNode == nil || !paramNode.IsName() {
		// no parameters, or first parameter uses destructuring or a default value
		return nil
	}

	statement := body.GetFirstChild()
	if !NodeUtil_IsExprAssign(statement) {
		return nil
	}

	assign := statement.GetFirstChild()
	lhs := assign.GetFirstChild()
	if lhs.IsGetProp() && lhs.GetFirstChild().IsThis() {
		rhs := assign.GetLastChild()
		if rhs.IsName() && rhs.GetString() == paramNode.GetString() {
			return rhino.IR_String(lhs.GetString()).ClonePropsFrom(lhs).Srcref(lhs)
		}
	}
	return nil
}
