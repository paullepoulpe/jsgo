// Port of com/google/javascript/jscomp/RewriteCallerCodeLocation.java @ closure-compiler v20261028

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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
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
	// Map of function name to a callerLocationFuncInfo value containing:
	// 1. the function variable
	// 2. the 1-based position of the param which has goog.callerLocation as default value
	// E.g:
	// `function signal(val, here: goog.CodeLocation = goog.callerLocation()) {}`
	// callerLocationFunctionNames["signal"] = {functionVar: Var(signal), paramPosition: 2}
	callerLocationFunctionNames map[string]callerLocationFuncInfo
	// Set of candidate property suffixes for GETPROP callees (e.g. "signal" for
	// "module$contents$main_signal") so non-matching GETPROP calls avoid
	// allocating qualified-name strings or running strings.ReplaceAll.
	callerLocationPropNames map[string]struct{}
}

// callerLocationFuncInfo stores the function variable and the 1-based position
// of the parameter that has goog.callerLocation() as its default value.
type callerLocationFuncInfo struct {
	functionVar   *Var
	paramPosition int
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
	return &RewriteCallerCodeLocation{
		compiler:                    compiler,
		astFactory:                  compiler.CreateAstFactory(),
		callerLocationFunctionNames: make(map[string]callerLocationFuncInfo),
	}
}

// Java: RewriteCallerCodeLocation.process(Node, Node)
func (r *RewriteCallerCodeLocation) Process(externs rhino.Node, root rhino.Node) {
	// Find all functions that have goog.callerLocation as a default parameter. We want to traverse
	// the root node as well as the externs because the function that contains `goog.callerLocation`
	// parameter could be an .i.js file in the externs node subtree.
	find := findCallerLocationFunctions{r: r}
	NodeTraversal_Traverse(r.compiler, r.compiler.GetRoot(), &find)
	if len(r.callerLocationFunctionNames) > 0 {
		// Rewrite call-sites of functions that have goog.callerLocation as a default parameter. We
		// want to traverse the root node as well as the externs to ensure the scope of the call-site
		// is the same scope as the function with the `goog.callerLocation` parameter.
		rewrite := rewriteCallerLocationCalls{r: r}
		NodeTraversal_Traverse(r.compiler, r.compiler.GetRoot(), &rewrite)
	}
}

var (
	findCallerLocationFunctionsVisitTokens = func() ast.TokenSet {
		var s ast.TokenSet
		s.Add(rhino.Token_GETPROP)
		s.Add(rhino.Token_PARAM_LIST)
		return s
	}()
	rewriteCallerLocationFunctionCallsVisitTokens = func() ast.TokenSet {
		var s ast.TokenSet
		s.Add(rhino.Token_CALL)
		return s
	}()
)

type findCallerLocationFunctions struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	r *RewriteCallerCodeLocation
}

func (*findCallerLocationFunctions) visitTokens() *ast.TokenSet {
	return &findCallerLocationFunctionsVisitTokens
}

func (f *findCallerLocationFunctions) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if isGoogCallerLocationMisused(n, parent) {
		f.r.compiler.Report(JSError_MakeWithNode(parent.GetParent(), RewriteCallerCodeLocation_JSC_CALLER_LOCATION_MISUSE_ERROR))
		return
	}

	if paramList, ok := ast.AsParamList(n); ok {
		f.r.visitParamListAndAddCallerLocationFunctionNames(paramList, t)
	}
}

// isGoogCallerLocationMisused checks if `goog.callerLocation` is misused.
// `goog.callerLocation` should only be used as a default parameter
// initializer.
func isGoogCallerLocationMisused(n rhino.Node, parent rhino.Node) bool {
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
			return true
		}

		// `goog.callerLocation` is used correctly as a default value in a function's parameter
		// list.
		return false
	}

	if assign, ok := ast.AsAssign(parent); ok && assign.Target().Node() == n && assign.Value().Node().IsFunction() {
		// This is the definition of goog.callerLocation, so this is not a misuse.
		// i.e: `goog.callerLocation = function(...) { ... }`
		return false
	}

	// `goog.callerLocation` is NOT used as a default value in a function's parameter list.
	return true
}

func (r *RewriteCallerCodeLocation) recordCallerLocationFunction(name string, info callerLocationFuncInfo) {
	r.callerLocationFunctionNames[name] = info
	for i := 0; i < len(name); i++ {
		if name[i] == '_' {
			if r.callerLocationPropNames == nil {
				r.callerLocationPropNames = make(map[string]struct{})
			}
			r.callerLocationPropNames[name[i+1:]] = struct{}{}
		}
	}
}

func (r *RewriteCallerCodeLocation) visitParamListAndAddCallerLocationFunctionNames(paramList ast.ParamList, t *NodeTraversal) {
	// function foo(<params>) ...
	// Each item in <params> is scanned for having a default value.
	// If there exists a default value, check if it is `goog.callerLocation`.
	// Check if there is another optional argument that comes before the goog.callerLocation arg.
	defaultValuesCount := 0
	paramPosition := 0 // 1-based position of the param in the param list.
	for param := range paramList.Params() {
		paramPosition++
		dv, ok := param.AsDefaultValue()
		if !ok {
			continue
		}
		defaultValuesCount++
		call, ok := dv.Value().AsCall()
		if !ok {
			continue
		}
		if RewriteCallerCodeLocation_GOOG_CALLER_LOCATION_QUALIFIED_NAME.Matches(call.Callee().Node()) {
			if defaultValuesCount > 1 {
				r.compiler.Report(JSError_MakeWithNode(dv.Node(), RewriteCallerCodeLocation_JSC_CALLER_LOCATION_POSITION_ERROR))
			}
			fn := ast.MustFunction(paramList.Node().GetParent())
			fnName, hasName := fn.Name()
			if !hasName {
				// Anonymous functions are not allowed to use goog.callerLocation.
				r.compiler.Report(
					JSError_MakeWithNode(
						fn.NameNode().Node(), RewriteCallerCodeLocation_JSC_ANONYMOUS_FUNCTION_CODE_LOCATION_ERROR))
				return
			}
			ident := fnName.Ident()
			r.recordCallerLocationFunction(ident, callerLocationFuncInfo{
				functionVar:   t.GetScope().GetVar(ident),
				paramPosition: paramPosition,
			})
		}
	}
}

type rewriteCallerLocationCalls struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	r *RewriteCallerCodeLocation
}

func (*rewriteCallerLocationCalls) noVisitLeaves() {}

func (*rewriteCallerLocationCalls) visitTokens() *ast.TokenSet {
	return &rewriteCallerLocationFunctionCallsVisitTokens
}

func (w *rewriteCallerLocationCalls) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if call, ok := ast.AsCall(n); ok {
		w.r.visitCallNodeAndRewrite(call, t)
	}
}

func (r *RewriteCallerCodeLocation) resolveCallerLocationCallee(callee ast.Expr) (string, callerLocationFuncInfo, bool) {
	if nameNode, ok := callee.AsName(); ok {
		moduleContentsName := nameNode.Ident()
		if strings.Contains(moduleContentsName, "module$exports$") {
			moduleContentsName = strings.ReplaceAll(moduleContentsName, "module$exports$", "module$contents$")
		}
		info, ok := r.callerLocationFunctionNames[moduleContentsName]
		return moduleContentsName, info, ok
	}

	getProp, ok := callee.AsGetProp()
	if !ok || len(r.callerLocationPropNames) == 0 {
		return "", callerLocationFuncInfo{}, false
	}
	prop := getProp.Property()
	if !strings.Contains(prop, "module$exports$") {
		if _, candidate := r.callerLocationPropNames[prop]; !candidate {
			return "", callerLocationFuncInfo{}, false
		}
	}
	qname, ok := getProp.Node().GetQualifiedNameString()
	if !ok {
		return "", callerLocationFuncInfo{}, false
	}
	// ClosureRewriteModule pass runs before this pass and renames exported function names
	// and aliases in goog.modules, e.g.:
	//   function name = module$contents$google3$javascript$apps$wiz$signals$signal_signal
	//   call node name = module$exports$google3$javascript$apps$wiz$signals$signal.signal
	moduleContentsName := strings.ReplaceAll(qname, ".", "_")
	if strings.Contains(moduleContentsName, "module$exports$") {
		moduleContentsName = strings.ReplaceAll(moduleContentsName, "module$exports$", "module$contents$")
	}
	info, ok := r.callerLocationFunctionNames[moduleContentsName]
	return moduleContentsName, info, ok
}

func (r *RewriteCallerCodeLocation) visitCallNodeAndRewrite(call ast.Call, t *NodeTraversal) {
	callee := call.Callee()
	moduleContentsName, info, ok := r.resolveCallerLocationCallee(callee)
	if !ok {
		return
	}

	callerLocationFunction := info.functionVar
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
	positionOfCallerLocationArg := info.paramPosition
	if call.NumArgs() >= positionOfCallerLocationArg {
		// goog.CodeLocation is provided as an argument, we will not rewrite this call-site.
		// If `undefined` is passed in as an argument to goog.CodeLocation, we will throw an error.
		// Otherwise continue without rewriting.
		arg, _ := call.AsCallOrNew().ArgAt(positionOfCallerLocationArg - 1)
		jsType := arg.Node().GetJSType()
		if jsType != nil && jsType.IsExplicitlyVoidable() {
			r.compiler.Report(JSError_MakeWithNode(callee.Node(), RewriteCallerCodeLocation_JSC_UNDEFINED_CODE_LOCATION_ERROR))
		}
		return
	}

	// create the goog.callerLocationIdInternalDoNotCallOrElse(path/to/file.ts:lineno:charno) and
	// add it as the last parameter.
	callNode := call.Node()
	xidCall := r.createGoogXidFilePathNode(callNode)
	r.compiler.ReportChangeToEnclosingScope(callNode)
	callNode.AddChildToBack(xidCall)
}

func (r *RewriteCallerCodeLocation) createGoogXidFilePathNode(n rhino.Node) rhino.Node {
	googNode := rhino.IR_Name("goog")
	googNode.SrcrefIfMissing(n)

	googXid := r.astFactory.CreateGetPropWithUnknownType(
		googNode, "callerLocationIdInternalDoNotCallOrElse")
	googXid.SrcrefIfMissing(n)

	callNode := r.astFactory.CreateCallWithUnknownType(googXid)
	callNode.SrcrefIfMissing(n)

	srcFile, ok := n.SourceFileName()
	if !ok {
		srcFile = "null"
	}
	stringNode := r.astFactory.CreateString(
		srcFile + ":" + strconv.Itoa(n.GetLineno()) + ":" + strconv.Itoa(n.GetCharno()))
	stringNode.SrcrefIfMissing(n)

	callNode.AddChildToBack(stringNode)
	return callNode
}
