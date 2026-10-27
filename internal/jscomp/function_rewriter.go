// Port of com/google/javascript/jscomp/FunctionRewriter.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/ast"
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
const functionRewriterSavingsThreshold = 16

// Java: FunctionRewriter(AbstractCompiler)
func NewFunctionRewriter(compiler AbstractCompiler) *FunctionRewriter {
	return &FunctionRewriter{compiler: compiler}
}

type functionRewriterReduction struct {
	oldChild rhino.Node
	newChild rhino.Node
}

type functionRewriterReducer struct {
	helperSourceName string
	helperSource     string
	reduce           func(fn ast.Function) (rhino.Node, bool)
}

// Order matches Java FunctionRewriter.process:
// [ReturnConstantReducer, GetterReducer, SetterReducer, EmptyFunctionReducer, IdentityReducer].
var functionRewriterReducers = [...]functionRewriterReducer{
	{
		helperSourceName: "class com.google.javascript.jscomp.FunctionRewriter$ReturnConstantReducer:helper",
		helperSource: "function JSCompiler_returnArg(JSCompiler_returnArg_value) {" +
			"  return function() {return JSCompiler_returnArg_value}" +
			"}",
		reduce: functionRewriter_reduceReturnConstant,
	},
	{
		helperSourceName: "class com.google.javascript.jscomp.FunctionRewriter$GetterReducer:helper",
		helperSource: "function JSCompiler_get(JSCompiler_get_name) {" +
			"  return function() {return this[JSCompiler_get_name]}" +
			"}",
		reduce: functionRewriter_reduceGetter,
	},
	{
		helperSourceName: "class com.google.javascript.jscomp.FunctionRewriter$SetterReducer:helper",
		helperSource: "function JSCompiler_set(JSCompiler_set_name) {" +
			"  return function(JSCompiler_set_value) {" +
			"this[JSCompiler_set_name] = JSCompiler_set_value}" +
			"}",
		reduce: functionRewriter_reduceSetter,
	},
	{
		helperSourceName: "class com.google.javascript.jscomp.FunctionRewriter$EmptyFunctionReducer:helper",
		helperSource: "function JSCompiler_emptyFn() {" +
			"  return function() {}" +
			"}",
		reduce: functionRewriter_reduceEmptyFunction,
	},
	{
		helperSourceName: "class com.google.javascript.jscomp.FunctionRewriter$IdentityReducer:helper",
		helperSource: "function JSCompiler_identityFn() {" +
			"  return function(JSCompiler_identityFn_value) {" +
			"return JSCompiler_identityFn_value}" +
			"}",
		reduce: functionRewriter_reduceIdentity,
	},
}

// Java: FunctionRewriter.process(Node, Node)
func (f *FunctionRewriter) Process(externs rhino.Node, root rhino.Node) {
	gatherer := &functionRewriterGatherer{}
	NodeTraversal_Traverse(f.compiler, root, gatherer)

	// Apply reductions iff they will provide some savings.
	for i := range functionRewriterReducers {
		reductions := gatherer.reductions[i]
		if len(reductions) == 0 {
			continue
		}

		helperCode := f.parseHelperCode(&functionRewriterReducers[i])
		if helperCode.IsNil() {
			continue
		}

		helperCodeCost := InlineCostEstimator_GetCost(helperCode)

		savings := 0
		for _, r := range reductions {
			savings += InlineCostEstimator_GetCost(r.oldChild) - InlineCostEstimator_GetCost(r.newChild)
		}

		// Compare estimated savings against the helper cost. Apply
		// reductions if doing so will result in some savings.
		if savings > helperCodeCost+functionRewriterSavingsThreshold {
			for _, r := range reductions {
				r.oldChild.ReplaceWith(r.newChild)
				NodeUtil_MarkFunctionsDeleted(r.oldChild, f.compiler)
				f.compiler.ReportChangeToEnclosingScope(r.newChild)
			}

			addingRoot := f.compiler.GetNodeForCodeInsertion(nil)
			NodeUtil_MarkNewScopesChanged(helperCode, f.compiler)
			addingRoot.AddChildToFront(helperCode)
			f.compiler.ReportChangeToEnclosingScope(addingRoot)
		}
	}
}

// Java: FunctionRewriter.parseHelperCode(Reducer)
func (f *FunctionRewriter) parseHelperCode(reducer *functionRewriterReducer) rhino.Node {
	root := f.compiler.ParseSyntheticCode(reducer.helperSourceName, reducer.helperSource)
	if !root.IsNil() {
		return root.RemoveFirstChild()
	}
	return rhino.NilNode
}

// Java: FunctionRewriter.isReduceableFunctionExpression(Node)
func functionRewriter_isReduceableFunctionExpression(fn ast.Function) bool {
	n := fn.Node()
	parent := n.GetParent()
	return NodeUtil_IsFunctionExpression(n) &&
		!NodeUtil_IsGetOrSetKey(parent) &&
		!parent.IsMemberFunctionDef() &&
		!fn.IsAsync() &&
		!fn.IsGenerator() &&
		!NodeUtil_HasNonSimpleParameters(n)
}

// functionRewriterGatherer gathers candidate reductions during pre-order
// traversal, stopping traversal inside reduced function subtrees.
type functionRewriterGatherer struct {
	reductions [len(functionRewriterReducers)][]functionRewriterReduction
}

func (*functionRewriterGatherer) noVisit() {}

func (g *functionRewriterGatherer) ShouldTraverse(_ *NodeTraversal, node rhino.Node, _ rhino.Node) bool {
	fn, ok := ast.AsFunction(node)
	if !ok || !functionRewriter_isReduceableFunctionExpression(fn) {
		return true
	}
	for i := range functionRewriterReducers {
		if replacement, reduced := functionRewriterReducers[i].reduce(fn); reduced {
			g.reductions[i] = append(g.reductions[i], functionRewriterReduction{
				oldChild: node,
				newChild: replacement,
			})
			return false
		}
	}
	return true
}

func (g *functionRewriterGatherer) Visit(*NodeTraversal, rhino.Node, rhino.Node) {}

func functionRewriter_buildCallNode(methodName string, arg rhino.Node, srcref rhino.Node) rhino.Node {
	call := ast.NewCall(ast.NewName(methodName).AsExpr())
	call.Node().SrcrefTree(srcref)
	call.SetFreeCall(true)
	if !arg.IsNil() {
		call.Node().AddChildToBack(arg)
	}
	return call.Node()
}

// Java: FunctionRewriter.SingleReturnStatementReducer.maybeGetSingleReturnRValue(Node)
func functionRewriter_singleReturnValue(fn ast.Function) (ast.Expr, bool) {
	block, ok := fn.Body().AsBlock()
	if !ok {
		return ast.Expr{}, false
	}
	stmt, ok := block.OnlyStatement()
	if !ok {
		return ast.Expr{}, false
	}
	ret, ok := stmt.AsReturn()
	if !ok {
		return ast.Expr{}, false
	}
	return ret.Value()
}

// Java: FunctionRewriter.EmptyFunctionReducer.reduce(Node)
func functionRewriter_reduceEmptyFunction(fn ast.Function) (rhino.Node, bool) {
	if NodeUtil_IsEmptyFunctionExpression(fn.Node()) {
		return functionRewriter_buildCallNode("JSCompiler_emptyFn", rhino.NilNode, fn.Node()), true
	}
	return rhino.NilNode, false
}

// Java: FunctionRewriter.IdentityReducer.reduce(Node)
func functionRewriter_reduceIdentity(fn ast.Function) (rhino.Node, bool) {
	firstParam, ok := fn.Params().FirstParam()
	if !ok {
		return rhino.NilNode, false
	}
	paramName, ok := firstParam.AsName()
	if !ok {
		return rhino.NilNode, false
	}
	retVal, ok := functionRewriter_singleReturnValue(fn)
	if !ok {
		return rhino.NilNode, false
	}
	retName, ok := retVal.AsName()
	if !ok || retName.Ident() != paramName.Ident() {
		return rhino.NilNode, false
	}
	return functionRewriter_buildCallNode("JSCompiler_identityFn", rhino.NilNode, fn.Node()), true
}

// Java: FunctionRewriter.ReturnConstantReducer.reduce(Node)
func functionRewriter_reduceReturnConstant(fn ast.Function) (rhino.Node, bool) {
	retVal, ok := functionRewriter_singleReturnValue(fn)
	if !ok || !NodeUtil_IsImmutableValue(retVal.Node()) {
		return rhino.NilNode, false
	}
	return functionRewriter_buildCallNode("JSCompiler_returnArg", retVal.Node().CloneTree(), fn.Node()), true
}

// Java: FunctionRewriter.GetterReducer.reduce(Node)
func functionRewriter_reduceGetter(fn ast.Function) (rhino.Node, bool) {
	if fn.IsArrow() {
		return rhino.NilNode, false
	}
	retVal, ok := functionRewriter_singleReturnValue(fn)
	if !ok {
		return rhino.NilNode, false
	}
	getprop, ok := retVal.AsGetProp()
	if !ok || !getprop.Object().Node().IsThis() {
		return rhino.NilNode, false
	}
	arg := rhino.IR_String(getprop.Property()).Srcref(getprop.Node())
	return functionRewriter_buildCallNode("JSCompiler_get", arg, fn.Node()), true
}

// Java: FunctionRewriter.SetterReducer.reduce(Node)
func functionRewriter_reduceSetter(fn ast.Function) (rhino.Node, bool) {
	if fn.IsArrow() {
		return rhino.NilNode, false
	}
	block, ok := fn.Body().AsBlock()
	if !ok {
		return rhino.NilNode, false
	}
	stmt, ok := block.OnlyStatement()
	if !ok {
		return rhino.NilNode, false
	}
	firstParam, ok := fn.Params().FirstParam()
	if !ok {
		return rhino.NilNode, false
	}
	paramName, ok := firstParam.AsName()
	if !ok {
		return rhino.NilNode, false
	}
	exprStmt, ok := stmt.AsExprStmt()
	if !ok {
		return rhino.NilNode, false
	}
	assign, ok := exprStmt.Expr().AsAssign()
	if !ok {
		return rhino.NilNode, false
	}
	lhs, ok := assign.Target().AsGetProp()
	if !ok || !lhs.Object().Node().IsThis() {
		return rhino.NilNode, false
	}
	rhs, ok := assign.Value().AsName()
	if !ok || rhs.Ident() != paramName.Ident() {
		return rhino.NilNode, false
	}
	propName := rhino.IR_String(lhs.Property()).ClonePropsFrom(lhs.Node()).Srcref(lhs.Node())
	return functionRewriter_buildCallNode("JSCompiler_set", propName, fn.Node()), true
}
