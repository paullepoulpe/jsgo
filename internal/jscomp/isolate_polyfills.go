// Port of com/google/javascript/jscomp/IsolatePolyfills.java @ closure-compiler v20261104

/*
 * Copyright 2020 The Closure Compiler Authors.
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

var (
	ipJscompPolyfillQName = rhino.QualifiedName_Of("$jscomp.polyfill")
	ipJscompPolyfillColID = uint32(strpool.Intern("$jscomp$polyfill"))
)

const isolatePolyfills_LOOKUP_METHOD = "$jscomp$lookupPolyfilledValue"

// IsolatePolyfills rewrites potential polyfill usages to use the hidden
// JSCompiler polyfills instead of the global.
//
// When $jscomp.ISOLATE_POLYFILLS is enabled, the $jscomp.polyfill library
// function does not add polyfills to the global scope or as properties on
// the native type. Instead, classes like Map are added to the
// $jscomp.polyfills object. Methods like `String.prototype.includes` are
// defined under a unique Symbol on String.prototype.
//
// This pass rewrites polyfill usages so that they access the actual
// polyfills instead of trying to access the native types. For example,
// `new Map()` becomes `new $jscomp.polyfills['Map']`.
//
// Limitations of this pass: a) it ignores destructuring and b) it does not
// support rewriting writes to polyfilled methods.
type IsolatePolyfills struct {
	compiler              AbstractCompiler
	polyfills             *PolyfillUsageFinder_Polyfills
	jscompPolyfillsObject rhino.Node

	jscompLookupMethod rhino.Node

	usedPolyfillMethodLookup bool
	isTempVarInitialized     bool
}

// Java: IsolatePolyfills.POLYFILL_TEMP
const isolatePolyfills_POLYFILL_TEMP = "$jscomp$polyfillTmp"

// Java: IsolatePolyfills(AbstractCompiler)
func NewIsolatePolyfills(compiler AbstractCompiler) *IsolatePolyfills {
	return NewIsolatePolyfillsWithPolyfills(compiler, defaultPolyfillsTable())
}

// Java: IsolatePolyfills(AbstractCompiler, Polyfills)
func NewIsolatePolyfillsWithPolyfills(compiler AbstractCompiler, polyfills *PolyfillUsageFinder_Polyfills) *IsolatePolyfills {
	i := &IsolatePolyfills{
		jscompLookupMethod: ast.NewName(isolatePolyfills_LOOKUP_METHOD).Node(),
		compiler:           compiler,
		polyfills:          polyfills,
	}
	if compiler.GetOptions().GetPropertyCollapseLevel() == CompilerOptions_PropertyCollapseLevel_ALL {
		i.jscompPolyfillsObject = IsolatePolyfills_CreateCollapsedName()
	} else {
		i.jscompPolyfillsObject = IsolatePolyfills_CreateJSCompPolyfillsAccess()
	}
	return i
}

// IsolatePolyfills_CreateCollapsedName returns a name `$jscomp$polyfills`
//
// Java: IsolatePolyfills.createCollapsedName()
func IsolatePolyfills_CreateCollapsedName() rhino.Node {
	collapsedName := ast.NewName("$jscomp$polyfills")
	collapsedName.SetConstantName(true)
	return collapsedName.Node()
}

// IsolatePolyfills_CreateJSCompPolyfillsAccess returns a getprop
// `$jscomp.polyfills`
//
// Java: IsolatePolyfills.createJSCompPolyfillsAccess()
func IsolatePolyfills_CreateJSCompPolyfillsAccess() rhino.Node {
	jscomp := ast.NewName("$jscomp")
	jscomp.SetConstantName(true)
	return ast.NewGetProp(jscomp.AsExpr(), "polyfills").Node()
}

// Java: IsolatePolyfills.process(Node, Node)
func (i *IsolatePolyfills) Process(externs rhino.Node, root rhino.Node) {
	// Calculate the set of polyfills that are actually present in the AST. It may be a subset of
	// the potential polyfills which PolyfillFindingCallback finds (it's fine if it's a superset.)
	injectedPolyfills := i.FindAllInjectedPolyfills()

	var polyfillUsages []*PolyfillUsageFinder_PolyfillUsage
	shouldDisambiguateProperties := i.compiler.GetOptions().ShouldDisambiguateProperties()
	NewPolyfillUsageFinder(i.compiler, i.polyfills, shouldDisambiguateProperties).
		TraverseIncludingGuarded(root, func(u *PolyfillUsageFinder_PolyfillUsage) {
			polyfillUsages = append(polyfillUsages, u)
		})

	visitedNodes := make(map[rhino.Node]struct{}, len(polyfillUsages))
	for _, usage := range polyfillUsages {
		node := usage.Node()
		_, visited := visitedNodes[node]
		if
		// Some nodes map to more than one polyfill usage. For example, `x.includes` maps to both
		// Array.prototype.includes and String.prototype.includes, but only needs to be isolated once.
		visited ||
			// Skip visiting nodes whose 'polyfill.library' is empty. This is true for language
			// features like `Proxy` and `String.raw` that have no associated polyfill, and hence are
			// unnecessary to isolate.
			usage.Polyfill().library == "" ||
			// The PolyfillFindingCallback may detect possible polyfill usages that are not
			// in fact injected. (possibly because RemoveUnusedCode deleted the polyfill.)
			!injectedPolyfills[usage.Polyfill().nativeSymbol] {
			continue
		}
		i.RewritePolyfill(usage)
		visitedNodes[node] = struct{}{}
	}

	i.CleanUpJscompLookupPolyfilledValue()
}

// FindAllInjectedPolyfills searches the AST for all calls to
// $jscomp.polyfill and returns the polyfilled symbol names.
//
// At the moment we don't track anywhere the set of all polyfills that have
// been injected. That set may be modified by RewritePolyfills,
// Es6InjectRuntimeLibraries, and RemoveUnusedCode. If desired, we could
// delete this method by making any passes that add/remove polyfill calls
// responsible for tracking their presence.
//
// Note: this set cannot be injected into the constructor because it is not
// known until the polyfill injection pass actually runs.
//
// Java: IsolatePolyfills.findAllInjectedPolyfills()
func (i *IsolatePolyfills) FindAllInjectedPolyfills() map[string]bool {
	actualPolyfills := make(map[string]bool)

	lastInjectedNode := i.compiler.GetNodeForCodeInsertion(nil)

	NodeTraversal_Traverse(
		i.compiler,
		lastInjectedNode,
		&isolatePolyfills_findAllInjectedPolyfillsCallback{outer: i, actualPolyfills: actualPolyfills})

	return actualPolyfills
}

// isolatePolyfills_findAllInjectedPolyfillsCallback is the anonymous
// NodeTraversal.AbstractShallowCallback of findAllInjectedPolyfills.
type isolatePolyfills_findAllInjectedPolyfillsCallback struct {
	NodeTraversal_AbstractShallowCallbackBase
	outer           *IsolatePolyfills
	actualPolyfills map[string]bool
}

// Java: IsolatePolyfills.findAllInjectedPolyfills().new AbstractShallowCallback() {...}.visit(NodeTraversal, Node, Node)
func (c *isolatePolyfills_findAllInjectedPolyfillsCallback) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if c.outer.IsJSCompPolyfillCall(n) {
		// CALL
		//  GETPROP/NAME $jscomp.polyfill
		//  STRING NativeSymbol.prototype.method
		//  [...]
		polyfilledSymbol := ast.MustCall(n).MustFirstArg().MustStringLit().Value()
		c.actualPolyfills[polyfilledSymbol] = true
	}
}

// Java: IsolatePolyfills.isJSCompPolyfillCall(Node)
func (i *IsolatePolyfills) IsJSCompPolyfillCall(n rhino.Node) bool {
	call, ok := ast.AsCall(n)
	if !ok {
		return false
	}
	callee := call.Callee()
	if i.compiler.GetOptions().GetPropertyCollapseLevel() == CompilerOptions_PropertyCollapseLevel_ALL {
		name, ok := callee.AsName()
		return ok && name.Node().GetStringID() == ipJscompPolyfillColID
	}
	return ipJscompPolyfillQName.Matches(callee.Node())
}

// RewritePolyfill rewrites a potential access of a polyfilled class or
// method to first look for the non-global, polyfilled version.
//
// Java: IsolatePolyfills.rewritePolyfill(PolyfillUsage)
func (i *IsolatePolyfills) RewritePolyfill(polyfillUsage *PolyfillUsageFinder_PolyfillUsage) {
	polyfill := polyfillUsage.Polyfill()

	// If the output FeatureSet includes all of the features of the language version in which
	// the polyfilled symbol was introduced, we can assume that the intended platform has
	// the symbol defined natively, so the compiler won't include the polyfill in its output and it
	// won't need to be isolated.
	if i.compiler.
		GetOptions().
		GetOutputFeatureSet().
		Contains(parser.FeatureSet_ValueOf(polyfill.nativeVersion)) {
		return
	}

	polyfillAccess := polyfillUsage.Node()
	parent := polyfillAccess.GetParent()

	if sourceFileName, ok := polyfillAccess.SourceFileName(); ok &&
		isolatePolyfills_isAllowedUnqualifiedPolyfillFile(sourceFileName) {
		return
	}

	// For now we need to assume that these lvalues are unrelated to the polyfills, as we are not
	// using any type information. If we wanted, we could support assignments to properties that are
	// already present, e.g.
	//   Array.includes = intercept;
	// to
	//   tmp = Array;
	//   tmp[$maybePolyfillProp(tmp, 'includes)] = intercept;
	role := ast.RoleIn(polyfillAccess, parent)
	if NodeUtil_IsLValue(polyfillAccess) || role == ast.RoleAssignTarget {
		return
	}

	name := polyfillUsage.Name()
	isGlobalClass := strings.IndexByte(name, '.') == -1 && polyfill.kind == PolyfillUsageFinder_Polyfill_Kind_STATIC

	if isGlobalClass {
		// e.g. Symbol, Map, window.Map, or goog.global.Map
		// Optional chaining is forbidden on global classes, hence producing `getelem` for property
		// access would suffice.
		polyfillAccess.ReplaceWith(
			// $jscomp.polyfills['Map']
			ast.NewGetElem(ast.MustExpr(i.jscompPolyfillsObject.CloneTree()), ast.NewStringLit(name).AsExpr()).
				Node().SrcrefTree(polyfillAccess))
	} else if role == ast.RoleCallCallee || role == ast.RoleOptChainCallCallee {
		// e.g. `getStr().includes('x')` or `getStr()?.includes('x')`
		i.RewritePolyfillInCall(polyfillAccess)
	} else {
		// e.g. [].includes.call(myIter, 0)
		propAccess := ast.MustNormalOrOptChainGetProp(polyfillAccess)
		methodName := ast.NewStringLit(propAccess.Property()).Node().Srcref(polyfillAccess)
		receiver := propAccess.Object().Detach().Node()
		// The `$jscomp$lookupPolyfilledValue` can handle both normal prop access as well as optional
		// by checking whether the lhs (receiver) is null or undefined first.
		polyfillAccess.ReplaceWith(
			i.CreatePolyfillMethodLookup(receiver, methodName, NodeUtil_IsOptChainNode(polyfillAccess)).
				SrcrefTree(polyfillAccess))
	}

	i.compiler.ReportChangeToEnclosingScope(parent)
}

// Code in the runtime libraries that may execute before any polyfills are injected and needs
// to opt out from polyfill isolation.
//   - util/global.js looks for `globalThis`
//   - util/shouldpolyfill.js checks whether Symbol is native
//   - es6/util/construct.js looks for the native Reflect.construct
//
// If desired we could programmatically detect these cases instead of having this allowlist of
// polyfill accesses, but that might silently allow other usages of third-party polyfills
// into the codebase.
// TODO(b/156776817): crash on early references to polyfills that are not in our allowlist.
//
// Java: IsolatePolyfills.FILES_ALLOWED_UNQUALIFIED_POLYFILL_ACCESSES.contains(String)
func isolatePolyfills_isAllowedUnqualifiedPolyfillFile(name string) bool {
	dir := GetRuntimeLibDir()
	switch name {
	case dir + "util/global.js",
		dir + "util/shouldpolyfill.js",
		dir + "es6/util/construct.js":
		return true
	default:
		return false
	}
}

// RewritePolyfillInCall rewrites a call where the receiver is a potential
// polyfilled method
//
// Before: `receiver.method(arg)`
//
// After: `lookupPolyfilledValue(receiver, 'method').call(receiver, arg)`
//
// Or, if evaluating the receiver may have side effects, we store the
// receiver in a temporary variable to avoid evaluating it twice:
//
// After: `(tmpNode = receiver, lookupPolyfilledValue(tmpNode, 'method')).call(tmpNode, arg)`
//
// Java: IsolatePolyfills.rewritePolyfillInCall(Node)
func (i *IsolatePolyfills) RewritePolyfillInCall(callee rhino.Node) {
	propAccess := ast.MustNormalOrOptChainGetProp(callee)
	methodName := ast.NewStringLit(propAccess.Property()).Node().Srcref(callee)
	receiver := propAccess.Object().Detach().Node()
	isCalleeOptChain := NodeUtil_IsOptChainNode(callee)

	var polyfilledMethod rhino.Node
	var thisNode rhino.Node

	if i.compiler.GetAstAnalyzer().MayEffectMutableState(receiver) {
		// e.g. `sideEffects().includes(arg)`
		thisNode = i.CreateTempName(callee)
		// (tmpNode = sideEffects(), lookupMethod(tmpNode, 'includes'))
		polyfilledMethod = rhino.IR_Comma(
			rhino.IR_Assign(thisNode.CloneTree(), receiver),
			i.CreatePolyfillMethodLookup(thisNode.CloneTree(), methodName, isCalleeOptChain))
	} else {
		thisNode = receiver
		polyfilledMethod = i.CreatePolyfillMethodLookup(receiver.CloneTree(), methodName, isCalleeOptChain)
	}

	// Fix the `this` type by using .call:
	//   lookupMethod(receiver, 'includes', isOptChainNode).call(receiver, arg)
	var receiverDotCall rhino.Node
	if isCalleeOptChain {
		// If optional chaining exists at this point, we can have it in the output
		receiverDotCall = rhino.IR_StartOptChainGetprop(polyfilledMethod, "call").SrcrefTree(callee)
	} else {
		receiverDotCall = ast.NewGetProp(ast.MustExpr(polyfilledMethod), "call").Node().SrcrefTree(callee)
	}
	callee.ReplaceWith(receiverDotCall)
	thisNode.InsertAfter(receiverDotCall)
}

// Java: IsolatePolyfills.createTempName(Node)
func (i *IsolatePolyfills) CreateTempName(srcref rhino.Node) rhino.Node {
	if !i.isTempVarInitialized {
		i.isTempVarInitialized = true
		decl := ast.NewNameDeclaration(rhino.Token_VAR, ast.NewName(isolatePolyfills_POLYFILL_TEMP).AsDeclarator()).Node().SrcrefTree(srcref)
		i.compiler.GetNodeForCodeInsertion(nil).AddChildToFront(decl)
		i.compiler.ReportChangeToEnclosingScope(decl)
	}
	// The same temporary variable is always used for every polyfill invocation. This is believed
	// to be safe and makes the code easier to generate and smaller. There's a change it will make
	// it harder for V8 to optimize, though. If proves to be a problem we could introduce unique tmp
	// variables.
	return ast.NewName(isolatePolyfills_POLYFILL_TEMP).Node().Srcref(srcref)
}

// CreatePolyfillMethodLookup returns a call
// `$jscomp$lookupPolyfilledValue(receiver, 'methodName', isOptChainNode)`
//
// Java: IsolatePolyfills.createPolyfillMethodLookup(Node, Node, boolean)
func (i *IsolatePolyfills) CreatePolyfillMethodLookup(receiver rhino.Node, methodName rhino.Node, isOptChainNode bool) rhino.Node {
	i.usedPolyfillMethodLookup = true
	callee := ast.MustExpr(i.jscompLookupMethod.CloneTree())
	var call ast.Call
	if isOptChainNode {
		call = ast.NewCall(callee, ast.MustArg(receiver), ast.MustArg(methodName), ast.NewTrue().AsArg())
	} else {
		// 3rd param of `jscompLookupMethod` is optional
		call = ast.NewCall(callee, ast.MustArg(receiver), ast.MustArg(methodName))
	}
	call.SetFreeCall(true)
	return call.Node()
}

// CleanUpJscompLookupPolyfilledValue deletes the dummy externs declaration
// of $jscomp$lookupPolyfilledValue and the function itself if unused.
//
// The RewritePolyfills pass injected a definition of
// $jscomp$lookupPolyfilledValue into the externs. This prevented dead code
// elimination, since the function is never unused until this pass runs.
// However, now we need to delete the externs definition so that variable
// renaming can actually rename $jscomp$lookupPolyfilledValue.
//
// Java: IsolatePolyfills.cleanUpJscompLookupPolyfilledValue()
func (i *IsolatePolyfills) CleanUpJscompLookupPolyfilledValue() {
	syntheticExternsRoot := i.compiler.GetSynthesizedExternsInput().GetAstRoot(i.compiler)
	syntheticExternsScope :=
		NewSyntacticScopeCreator(i.compiler).CreateScopeOfScope(syntheticExternsRoot /* parent= */, nil)
	externVar :=
		jutil.CheckNotNullMsg(
			syntheticExternsScope.GetVar(isolatePolyfills_LOOKUP_METHOD),
			"Failed to find synthetic $jscomp$lookupPolyfilledValue extern")
	NodeUtil_DeleteNode(externVar.GetParentNode(), i.compiler)

	if i.usedPolyfillMethodLookup {
		return
	}

	syntheticCodeScope :=
		NewSyntacticScopeCreator(i.compiler).
			CreateScopeOfScope(i.compiler.GetNodeForCodeInsertion( /* chunk= */ nil) /* parent= */, nil)
	syntheticVar := syntheticCodeScope.GetVar(isolatePolyfills_LOOKUP_METHOD)
	// Don't error if we can't find a definition for jscompLookupMethod. It's possible that we are
	// running in transpileOnly mode and are not injecting runtime libraries.
	if syntheticVar != nil {
		NodeUtil_DeleteNode(syntheticVar.GetParentNode(), i.compiler)
	}
}
