// Port of com/google/javascript/jscomp/RewriteCallerCodeLocation.java @ closure-compiler v20261021

/*
 * Copyright 2024 The Closure Compiler Authors.
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

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// RewriteCallerCodeLocation rewrites call-sites of functions that have
// goog.callerLocation as a default parameter.
//
// E.g: function signal(val, here: goog.CodeLocation = goog.callerLocation()) {}
//
// When a function is called without providing the optional `here` argument,
// we will rewrite the function call to include the code location.
//
// i.e. `signal(0, goog.callerLocationIdInternalDoNotCallOrElse(path/to/file.ts:lineno:charno))`
type RewriteCallerCodeLocation struct {
	compiler   AbstractCompiler
	astFactory *AstFactory
	// Map of function name to a FunctionVarAndParamPosition object which contains:
	// 1. the function variable
	// 2. the position of the param which has goog.callerLocation as default value
	// E.g:
	// `function signal(val, here: goog.CodeLocation = goog.callerLocation()) {}`
	// callerLocationFunctionNames.put("signal", {functionVar: `Var signal @ NAME signal 1:9 ...`,
	// paramPosition: 2}); // 2 because "here" is the second param
	callerLocationFunctionNames *jutil.LinkedHashMap[string, *RewriteCallerCodeLocation_FunctionVarAndParamPosition]
}

// Java: RewriteCallerCodeLocation.JSC_CALLER_LOCATION_POSITION_ERROR
var RewriteCallerCodeLocation_JSC_CALLER_LOCATION_POSITION_ERROR = DiagnosticType_Error(
	"JSC_CALLER_LOCATION_POSITION_ERROR",
	"Please make sure there is only one goog.callerLocation argument in your function's"+
		" parameter list, and it is the first optional argument in the list")

// Java: RewriteCallerCodeLocation.JSC_CALLER_LOCATION_MISUSE_ERROR
var RewriteCallerCodeLocation_JSC_CALLER_LOCATION_MISUSE_ERROR = DiagnosticType_Error(
	"JSC_CALLER_LOCATION_MISUSE_ERROR",
	"goog.callerLocation should only be used as a default parameter initializer")

// Java: RewriteCallerCodeLocation.JSC_UNDEFINED_CODE_LOCATION_ERROR
var RewriteCallerCodeLocation_JSC_UNDEFINED_CODE_LOCATION_ERROR = DiagnosticType_Error(
	"JSC_UNDEFINED_CODE_LOCATION_ERROR",
	"Do not pass in undefined as an argument to goog.CodeLocation parameter")

// Java: RewriteCallerCodeLocation.JSC_ANONYMOUS_FUNCTION_CODE_LOCATION_ERROR
var RewriteCallerCodeLocation_JSC_ANONYMOUS_FUNCTION_CODE_LOCATION_ERROR = DiagnosticType_Error(
	"JSC_ANONYMOUS_FUNCTION_CODE_LOCATION_ERROR",
	"Do not use goog.callerLocation in an anonymous functions. Functions that use"+
		" goog.callerLocation should be named.")

// Java: RewriteCallerCodeLocation.GOOG_CALLER_LOCATION_QUALIFIED_NAME
var RewriteCallerCodeLocation_GOOG_CALLER_LOCATION_QUALIFIED_NAME = rhino.QualifiedName_Of("goog.callerLocation")

// NewRewriteCallerCodeLocation creates an instance.
//
// Java: RewriteCallerCodeLocation(AbstractCompiler)
func NewRewriteCallerCodeLocation(compiler AbstractCompiler) *RewriteCallerCodeLocation {
	r := &RewriteCallerCodeLocation{}
	r.compiler = compiler
	r.astFactory = compiler.CreateAstFactory()
	r.callerLocationFunctionNames = jutil.NewLinkedHashMap[string, *RewriteCallerCodeLocation_FunctionVarAndParamPosition]()
	return r
}

// Java: RewriteCallerCodeLocation.process(Node, Node)
func (r *RewriteCallerCodeLocation) Process(externs rhino.Node, root rhino.Node) {
	// Find all functions that have goog.callerLocation as a default parameter. We want to traverse
	// the root node as well as the externs because the function that contains `goog.callerLocation`
	// parameter could be an
	// .i.js file in the externs node subtree
	NodeTraversal_Traverse(r.compiler, r.compiler.GetRoot(), NewRewriteCallerCodeLocation_FindCallerLocationFunctions(r))
	if !r.callerLocationFunctionNames.IsEmpty() {
		// Rewrite call-sites of functions that have goog.callerLocation as a default parameter. We
		// want to traverse the root node as well as the externs to ensure the scope of the call-site
		// is the same scope as the function with the `goog.callerLocation` parameter.
		NodeTraversal_Traverse(
			r.compiler, r.compiler.GetRoot(), NewRewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls(r))
	}
}

// RewriteCallerCodeLocation_FindCallerLocationFunctions is the inner class
// FindCallerLocationFunctions.
type RewriteCallerCodeLocation_FindCallerLocationFunctions struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	this_ *RewriteCallerCodeLocation // the enclosing instance
}

// Java: FindCallerLocationFunctions()
func NewRewriteCallerCodeLocation_FindCallerLocationFunctions(outer *RewriteCallerCodeLocation) *RewriteCallerCodeLocation_FindCallerLocationFunctions {
	return &RewriteCallerCodeLocation_FindCallerLocationFunctions{this_: outer}
}

// Java: FindCallerLocationFunctions.visit(NodeTraversal, Node, Node)
func (f *RewriteCallerCodeLocation_FindCallerLocationFunctions) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if f.IsGoogCallerLocationMisused(n, parent) {
		f.this_.compiler.Report(JSError_MakeWithNode(parent.GetParent(), RewriteCallerCodeLocation_JSC_CALLER_LOCATION_MISUSE_ERROR))
		return
	}

	if n.IsParamList() {
		f.VisitParamListAndAddCallerLocationFunctionNames(n, t)
	}
}

// IsGoogCallerLocationMisused checks if `goog.callerLocation` is misused.
// `goog.callerLocation` should only be used as a default parameter
// initializer.
//
// n is the node to check, parent its parent. Returns true if
// goog.callerLocation is misused, false otherwise.
//
// Java: FindCallerLocationFunctions.isGoogCallerLocationMisused(Node, Node)
func (f *RewriteCallerCodeLocation_FindCallerLocationFunctions) IsGoogCallerLocationMisused(n rhino.Node, parent rhino.Node) bool {
	if !RewriteCallerCodeLocation_GOOG_CALLER_LOCATION_QUALIFIED_NAME.Matches(n) {
		// not a `goog.callerLocation` node
		return false
	}

	// Check for misuse of goog.callerLocation.
	if parent.IsCall() && parent.GetParent().IsDefaultValue() {
		if parent.GetGrandparent().IsStringKey() {
			// Throw an error when `goog.callerLocation` is used in an object literal.
			// E.g:
			// function foo({val1, val2, here = goog.callerLocation()}) {}
			// The AST for `here = goog.callerLocation()`looks like:
			// STRING_KEY here (This node tells us we are in an object literal)
			//   DEFAULT_VALUE
			//     NAME here 1:26
			//     CALL 1:33
			//       GETPROP callerLocation (This is the node `n` we are currently at)
			//         NAME goog 1:33
			return true
		}

		// `goog.callerLocation` is used correctly as a default value in a function's parameter
		// list.
		return false
	}

	if parent.IsAssign() && !n.GetNext().IsNil() && n.GetNext().IsFunction() {
		// This is the definition of goog.callerLocation, so this is not a misuse.
		// i.e: `goog.callerLocation = function(...) { ... }`
		return false
	}

	// `goog.callerLocation` is NOT used as a default value in a function's parameter list.
	return true
}

// VisitParamListAndAddCallerLocationFunctionNames visits the param list of a
// function and checks if it contains a default value of
// goog.callerLocation. If it does, `visitParamList` stores the function name
// and the index of the param in `callerLocationFunctionNames` map.
//
// Example: function signal(val, here: goog.CodeLocation = goog.callerLocation()) {}
//
// `callerLocationFunctionNames` will have an entry for "signal" -> 2
// (because "here" is the second param)
//
// n is the param list node.
//
// Java: FindCallerLocationFunctions.visitParamListAndAddCallerLocationFunctionNames(Node, NodeTraversal)
func (f *RewriteCallerCodeLocation_FindCallerLocationFunctions) VisitParamListAndAddCallerLocationFunctionNames(n rhino.Node, t *NodeTraversal) {
	jutil.CheckStateMsg1(n.IsParamList(), "%s", n)
	// function foo(<params>) ...
	// Each item in <params> is scanned for having a default value.
	// If there exists a default value, check if it is `goog.callerLocation`.
	// Check if there is another optional argument that comes before the goog.callerLocation arg.
	defaultValuesCount := 0
	paramPosition := 0 // keep track of the position of the param in the param list.
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		paramPosition++
		if !c.IsDefaultValue() {
			continue
		}
		defaultValuesCount++
		call := c.GetSecondChild() // CALL node (E.g: goog.callerLocation())
		if !call.IsCall() {
			continue
		}
		getProp := call.GetFirstChild()
		if RewriteCallerCodeLocation_GOOG_CALLER_LOCATION_QUALIFIED_NAME.Matches(getProp) {
			if defaultValuesCount > 1 {
				f.this_.compiler.Report(JSError_MakeWithNode(c, RewriteCallerCodeLocation_JSC_CALLER_LOCATION_POSITION_ERROR))
			}
			// E.g: function signal(val, here: goog.CodeLocation = goog.callerLocation()) {}
			// Add ("signal" -> 2) to `callerLocationFunctionNames`
			functionName := n.GetParent().GetFirstChild().GetQualifiedName()
			if functionName == nil {
				// Anonymous functions are not allowed to use goog.callerLocation.
				f.this_.compiler.Report(
					JSError_MakeWithNode(
						n.GetParent().GetFirstChild(), RewriteCallerCodeLocation_JSC_ANONYMOUS_FUNCTION_CODE_LOCATION_ERROR))
				return
			}

			functionVarAndPosition :=
				NewRewriteCallerCodeLocation_FunctionVarAndParamPosition(t.GetScope().GetVar(*functionName), paramPosition)
			f.this_.callerLocationFunctionNames.Put(*functionName, functionVarAndPosition)
		}
	}
}

// RewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls is the inner
// class RewriteCallerLocationFunctionCalls.
type RewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	this_ *RewriteCallerCodeLocation // the enclosing instance
}

// Java: RewriteCallerLocationFunctionCalls()
func NewRewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls(outer *RewriteCallerCodeLocation) *RewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls {
	return &RewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls{this_: outer}
}

// Java: RewriteCallerLocationFunctionCalls.visit(NodeTraversal, Node, Node)
func (w *RewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsCall() {
		w.VisitCallNodeAndRewrite(n, t)
	}
}

// VisitCallNodeAndRewrite visits call expression nodes and checks if they
// require transformations. If they do, it rewrites the call expression node
// to include the code location.
//
// E.g: `signal(0)` will be rewritten to `signal(0,
// goog.callerLocationIdInternalDoNotCallOrElse(path/to/file.ts:lineno:charno))`
//
// n is the call node.
//
// Java: RewriteCallerLocationFunctionCalls.visitCallNodeAndRewrite(Node, NodeTraversal)
func (w *RewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls) VisitCallNodeAndRewrite(n rhino.Node, t *NodeTraversal) {
	firstChild := n.GetFirstChild()
	if firstChild.IsNil() ||
		(!firstChild.IsName() && !firstChild.IsGetProp()) ||
		firstChild.GetQualifiedName() == nil {
		return
	}

	// Check if the call-site is calling a callerLocation function.
	// E.g: signal(0); name = "signal"
	name := *firstChild.GetQualifiedName()

	// ClosureRewriteModule pass will run before this pass and will rename exported function names
	// and aliases in goog.modules.
	// E.g:
	// function name = module$contents$google3$javascript$apps$wiz$signals$signal_signal
	// call node name = module$exports$google3$javascript$apps$wiz$signals$signal.signal
	moduleContentsName :=
		jutil.StringReplace(jutil.StringReplace(name, ".", "_"), "module$exports$", "module$contents$")

	if !w.this_.callerLocationFunctionNames.ContainsKey(moduleContentsName) {
		return
	}

	functionVarAndPosition := w.this_.callerLocationFunctionNames.Get(moduleContentsName)

	callerLocationFunction := functionVarAndPosition.GetFunctionVar()
	calleeFunction := t.GetScope().GetVar(moduleContentsName)
	if callerLocationFunction == nil ||
		calleeFunction == nil ||
		callerLocationFunction.GetNameNode() != calleeFunction.GetNameNode() {
		return
	}

	// Check if the argument is provided.
	// E.g:
	// function signal(val, here: goog.CodeLocation = goog.callerLocation()) {}
	// signal(0, xid('path/to/file.ts:25')); // goog.CodeLocation is provided
	positionOfCallerLocationArg := functionVarAndPosition.GetParamPosition()
	numberOfArgs := n.GetChildCount() - 1 // -1 for the firstChild which is the function name.
	if numberOfArgs >= positionOfCallerLocationArg {
		// goog.CodeLocation is provided as an argument, we will not rewrite this call-site.
		// If `undefined` is passed in as an argument to goog.CodeLocation, we will throw an error.
		// Otherwise continue without rewriting.
		jsType := n.GetChildAtIndex(positionOfCallerLocationArg).GetJSType()
		if jsType != nil && jsType.IsExplicitlyVoidable() {
			// user is passing in undefined as an argument to goog.CodeLocation
			w.this_.compiler.Report(JSError_MakeWithNode(firstChild, RewriteCallerCodeLocation_JSC_UNDEFINED_CODE_LOCATION_ERROR))
		}

		return
	}

	// create the goog.callerLocationIdInternalDoNotCallOrElse(path/to/file.ts:lineno:charno) and
	// add it as the last parameter.
	xidCall := w.CreateGoogXidFilePathNode(n)
	w.this_.compiler.ReportChangeToEnclosingScope(n)
	n.AddChildToBack(xidCall)
}

// CreateGoogXidFilePathNode creates a call node for
// "goog.callerLocationIdInternalDoNotCallOrElse(path/to/file.ts:lineno:charno)"
//
// n is the call node of a function that needs to be rewritten to include the
// code location. Returns the call node including code location i.e.
// "goog.callerLocationIdInternalDoNotCallOrElse(path/to/file.ts:lineno:charno)"
//
// Java: RewriteCallerLocationFunctionCalls.createGoogXidFilePathNode(Node)
func (w *RewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls) CreateGoogXidFilePathNode(n rhino.Node) rhino.Node {
	// googNode is "goog"
	googNode := rhino.IR_Name("goog")
	googNode.SrcrefIfMissing(n)

	// googXid is "goog.callerLocationIdInternalDoNotCallOrElse" node
	googXid := w.this_.astFactory.CreateGetPropWithUnknownType(
		googNode, "callerLocationIdInternalDoNotCallOrElse")
	googXid.SrcrefIfMissing(n)

	// callNode is "goog.callerLocationIdInternalDoNotCallOrElse()" node
	callNode := w.this_.astFactory.CreateCallWithUnknownType(googXid)
	callNode.SrcrefIfMissing(n)

	// stringNode is "path/to/file.ts:lineno:charno" node
	stringNode := w.this_.astFactory.CreateString(
		jutil.ToJavaString(n.GetSourceFileName()) + ":" + strconv.Itoa(n.GetLineno()) + ":" + strconv.Itoa(n.GetCharno()))
	stringNode.SrcrefIfMissing(n)

	// turns callNode from "goog.callerLocationIdInternalDoNotCallOrElse()" to
	// "goog.callerLocationIdInternalDoNotCallOrElse(path/to/file.ts:lineno:charno)"
	callNode.AddChildToBack(stringNode)

	return callNode
}

// RewriteCallerCodeLocation_FunctionVarAndParamPosition stores the function
// variable and the position of the param which has goog.callerLocation as
// default value.
//
// Since it is not guaranteed that names are unique at this point in
// compilation, we could have something like (1) a function using
// goog.callerLocation that's actually nested within another function, and
// not available globally or (2) a function using goog.callerLocation that's
// then shadowed by another function locally. We do not want to rewrite the
// call-site in these cases, so we'll check the scope of the call-site
// against the function with the `goog.callerLocation` default parameter.
//
// E.g: `function signal(val, here: goog.CodeLocation = goog.callerLocation()) {}`
//
// This class will have:
//
// 1. functionVar: `Var signal @ NAME signal 1:9 ...`
//
// 2. paramPosition: 2 (because "here" is the second param)
//
// Storing the functionVar is useful for checking if the call-site is calling
// the same function as the one with the `goog.callerLocation` default
// parameter.
type RewriteCallerCodeLocation_FunctionVarAndParamPosition struct {
	functionVar   *Var
	paramPosition int
}

// Java: FunctionVarAndParamPosition(Var, int)
func NewRewriteCallerCodeLocation_FunctionVarAndParamPosition(functionVar *Var, paramPosition int) *RewriteCallerCodeLocation_FunctionVarAndParamPosition {
	return &RewriteCallerCodeLocation_FunctionVarAndParamPosition{functionVar: functionVar, paramPosition: paramPosition}
}

// Java: FunctionVarAndParamPosition.getFunctionVar()
func (p *RewriteCallerCodeLocation_FunctionVarAndParamPosition) GetFunctionVar() *Var {
	return p.functionVar
}

// Java: FunctionVarAndParamPosition.getParamPosition()
func (p *RewriteCallerCodeLocation_FunctionVarAndParamPosition) GetParamPosition() int {
	return p.paramPosition
}
