// Port of com/google/javascript/jscomp/RescopeGlobalSymbols.java @ closure-compiler v20260415

/*
 * Copyright 2011 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// RescopeGlobalSymbols finds all references to global symbols and rewrites
// them to be property accesses to a special object with the same name as the
// global symbol.
//
// Given the name of the global object is NS
//
//	var a = 1; function b() { return a }
//
// becomes
//
//	NS.a = 1; NS.b = function b() { return NS.a }
//
// This allows splitting code into chunks that depend on each other's global
// symbols, without using polluting JavaScript's global scope with those
// symbols. You typically define just a single global symbol, wrap each chunk
// in a function wrapper, and pass the global symbol around, eg,
//
//	var uniqueNs = uniqueNs || {};
//
//	(function (NS) { ...your chunk code here... })(uniqueNs);
//
// This compile step requires rewriteGlobalDeclarationsForTryCatchWrapping to
// be turned on to guarantee semantics.
//
// For lots of examples, see the unit test.
//
// Java: final class RescopeGlobalSymbols implements CompilerPass
type RescopeGlobalSymbols struct {
	compiler              AbstractCompiler
	globalSymbolNamespace string
	addExtern             bool
	assumeCrossChunkNames bool
	crossChunkNames       *jutil.LinkedHashSet[string]

	// Global identifiers that may be a non-arrow function referencing "this"
	maybeReferencesThis *jutil.LinkedHashSet[string]

	externNames *jutil.LinkedHashSet[string]
}

// Appended to variables names that conflict with globalSymbolNamespace.
const rescopeGlobalSymbols_DISAMBIGUATION_SUFFIX = "$"

// NewRescopeGlobalSymbols is the constructor for the RescopeGlobalSymbols
// compiler pass.
//
// compiler is the JSCompiler, for reporting code changes;
// globalSymbolNamespace is the name of namespace into which all global
// symbols are transferred. If assumeCrossChunkNames is true, all global
// symbols will be assumed cross chunk boundaries and thus require renaming.
//
// Java: RescopeGlobalSymbols(AbstractCompiler, String, boolean)
func NewRescopeGlobalSymbols(
	compiler AbstractCompiler, globalSymbolNamespace string, assumeCrossChunkNames bool) *RescopeGlobalSymbols {
	return NewRescopeGlobalSymbolsWithAddExtern(compiler, globalSymbolNamespace, true, assumeCrossChunkNames)
}

// NewRescopeGlobalSymbolsWithAddExtern is the constructor for the
// RescopeGlobalSymbols compiler pass for use in testing.
//
// compiler is the JSCompiler, for reporting code changes;
// globalSymbolNamespace is the name of namespace into which all global
// symbols are transferred. If addExtern is true, the compiler will consider
// the globalSymbolNamespace an extern name. If assumeCrossChunkNames is true,
// all global symbols will be assumed cross chunk boundaries and thus require
// renaming. VisibleForTesting
//
// Java: RescopeGlobalSymbols(AbstractCompiler, String, boolean, boolean)
func NewRescopeGlobalSymbolsWithAddExtern(
	compiler AbstractCompiler,
	globalSymbolNamespace string,
	addExtern bool,
	assumeCrossChunkNames bool) *RescopeGlobalSymbols {
	r := &RescopeGlobalSymbols{
		crossChunkNames:     jutil.NewLinkedHashSet[string](),
		maybeReferencesThis: jutil.NewLinkedHashSet[string](),
	}
	r.compiler = compiler
	r.globalSymbolNamespace = globalSymbolNamespace
	r.addExtern = addExtern
	r.assumeCrossChunkNames = assumeCrossChunkNames
	return r
}

// Java: RescopeGlobalSymbols.isCrossChunkName(String)
func (r *RescopeGlobalSymbols) IsCrossChunkName(name string) bool {
	return r.assumeCrossChunkNames ||
		r.crossChunkNames.Contains(name) ||
		r.compiler.GetCodingConvention().IsExportedWithLocal(name /* local= */, false)
}

// Java: RescopeGlobalSymbols.isExternVar(String, NodeTraversal)
func (r *RescopeGlobalSymbols) IsExternVar(varname string, t *NodeTraversal) bool {
	if varname == "" {
		return false
	}
	v := t.GetScope().GetVar(varname)
	return v == nil ||
		v.IsExtern() ||
		(v.GetScope().(*Scope).IsGlobal() && r.externNames.Contains(varname))
}

// Java: RescopeGlobalSymbols.addExternForGlobalSymbolNamespace()
func (r *RescopeGlobalSymbols) AddExternForGlobalSymbolNamespace() {
	varNode := rhino.IR_Var(rhino.IR_Name(r.globalSymbolNamespace))
	input := r.compiler.GetSynthesizedExternsInput()
	input.GetAstRoot(r.compiler).AddChildToBack(varNode)
	r.compiler.ReportChangeToEnclosingScope(varNode)
}

// Java: RescopeGlobalSymbols.process(Node, Node)
func (r *RescopeGlobalSymbols) Process(externs *rhino.Node, root *rhino.Node) {
	// Collect variables in externs; they can be shadowed by the same names in global scope.
	r.externNames = NodeUtil_CollectExternVariableNames(r.compiler, externs)

	// Make the name of the globalSymbolNamespace an extern.
	if r.addExtern {
		r.AddExternForGlobalSymbolNamespace()
	}

	// Rewrite all references to global symbols to properties of a single symbol:

	// Turn global named function statements into var assignments.
	NodeTraversal_Traverse(
		r.compiler, root, NewRescopeGlobalSymbols_RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback(r))

	// Find global names that are used in more than one chunk. Those that
	// are have to be rewritten.
	var nonMutatingPasses []NodeTraversal_Callback
	nonMutatingPasses = append(nonMutatingPasses, &RescopeGlobalSymbols_FindCrossChunkNamesCallback{outer: r})

	// And find names that may refer to functions that reference this.
	nonMutatingPasses = append(nonMutatingPasses, &RescopeGlobalSymbols_FindNamesReferencingThis{outer: r})
	CombinedCompilerPass_Traverse(r.compiler, root, nonMutatingPasses)

	// Rewrite all references to be property accesses of the single symbol.
	rewriteScope := &RescopeGlobalSymbols_RewriteScopeCallback{outer: r, preDeclarations: jutil.NewArrayList[*RescopeGlobalSymbols_RewriteScopeCallback_ChunkGlobal]()}
	NodeTraversal_Traverse(r.compiler, root, rewriteScope)

	// Remove the var from statements in global scope if the declared names have been rewritten
	// in the previous pass.
	NodeTraversal_Traverse(r.compiler, root, NewRescopeGlobalSymbols_RemoveGlobalVarCallback(r))
	rewriteScope.DeclareChunkGlobals()
}

// RescopeGlobalSymbols_RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback
// rewrites global function and class declarations to var statements +
// assignment. Ignores non-global function and class declarations.
//
//	function test(){}
//
// becomes
//
//	var test = function (){}
//
//	class A {}
//
// becomes
//
//	var A = class {}
//
// After this traversal, the special case of global class and function
// statements can be ignored.
//
// This is helpful when rewriting simple names to property accesses on the
// global symbol, since class A {} cannot be rewritten directly to
// class NS.A {}.
//
// Java: private class RescopeGlobalSymbols.RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback extends AbstractShallowStatementCallback
type RescopeGlobalSymbols_RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback struct {
	NodeTraversal_AbstractShallowStatementCallbackBase
	outer *RescopeGlobalSymbols
}

// Java: RescopeGlobalSymbols.RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback()
func NewRescopeGlobalSymbols_RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback(outer *RescopeGlobalSymbols) *RescopeGlobalSymbols_RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback {
	c := &RescopeGlobalSymbols_RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback{outer: outer}
	return c
}

// Java: RescopeGlobalSymbols.RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback.visit(NodeTraversal, Node, Node)
func (c *RescopeGlobalSymbols_RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	compiler := c.outer.compiler
	// Ignore block scopes within the global scope, as class and function declarations are
	// block-scoped.
	// Note that we should never find block-scoped function declarations if outputting ES5
	// code. Es6RewriteBlockScopedFunctionDeclaration will have rewritten them.
	if !t.InGlobalScope() {
		return
	}
	// Ignore everything that's not a function or class declaration.
	if !NodeUtil_IsFunctionDeclaration(n) && !NodeUtil_IsClassDeclaration(n) {
		return
	}
	nameNode := NodeUtil_GetNameNode(n)
	name := nameNode.GetString()
	// Remove the class or function name. Anonymous classes have an EMPTY node, while anonymous
	// functions have a NAME node with an empty string.
	if n.IsClass() {
		nameNode.ReplaceWith(rhino.IR_Empty().Srcref(nameNode))
	} else {
		nameNode.SetString("")
		compiler.ReportChangeToEnclosingScope(nameNode)
	}
	prev := n.GetPrevious()
	n.Detach()
	var_ := NodeUtil_NewVarNode(name, n)
	if prev == nil {
		parent.AddChildToFront(var_)
	} else {
		var_.InsertAfter(prev)
	}
	compiler.ReportChangeToEnclosingScope(parent)
}

// RescopeGlobalSymbols_FindCrossChunkNamesCallback finds all global names
// that are used in more than one chunk. The following compiler
// transformations can ignore the globals that are not.
//
// Java: private class RescopeGlobalSymbols.FindCrossChunkNamesCallback extends AbstractPostOrderCallback
type RescopeGlobalSymbols_FindCrossChunkNamesCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer *RescopeGlobalSymbols
}

// Java: RescopeGlobalSymbols.FindCrossChunkNamesCallback.visit(NodeTraversal, Node, Node)
func (c *RescopeGlobalSymbols_FindCrossChunkNamesCallback) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	r := c.outer
	if n.IsName() {
		name := n.GetString()
		if name == "" || r.crossChunkNames.Contains(name) {
			return
		}
		s := t.GetScope()
		v := s.GetVar(name)
		if v == nil || !v.IsGlobal() {
			return
		}
		input := v.GetInput()
		if input == nil {
			// We know nothing. Assume name is used across chunks.
			r.crossChunkNames.Add(name)
			return
		}
		// Compare the chunk where the variable is declared to the current
		// chunk. If they are different, the variable is used across chunks.
		chunk := input.GetChunk()
		if chunk != t.GetChunk() {
			r.crossChunkNames.Add(name)
		}
	}
}

// RescopeGlobalSymbols_FindNamesReferencingThis builds the
// maybeReferencesThis set of names that may reference a function that
// references this. If the function a name references does not reference this
// it can be called as a method call where the this value is not the same as
// in a normal function call.
//
// Java: private class RescopeGlobalSymbols.FindNamesReferencingThis extends AbstractPostOrderCallback
type RescopeGlobalSymbols_FindNamesReferencingThis struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer *RescopeGlobalSymbols
}

// Java: RescopeGlobalSymbols.FindNamesReferencingThis.visit(NodeTraversal, Node, Node)
func (c *RescopeGlobalSymbols_FindNamesReferencingThis) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	r := c.outer
	if n.IsName() {
		name := n.GetString()
		if name == "" {
			return
		}
		var value *rhino.Node
		if parent.IsAssign() && n == parent.GetFirstChild() {
			value = parent.GetLastChild()
		} else if NodeUtil_IsNameDeclaration(parent) {
			value = n.GetFirstChild()
		} else if parent.IsFunction() {
			value = parent
		}
		if value == nil && !NodeUtil_IsLhsByDestructuring(n) {
			// If n is assigned in a destructuring pattern, don't bother finding its value and just
			// assume it may reference this.
			return
		}
		// We already added this symbol. Done after checks above because those
		// are comparatively cheap.
		if r.maybeReferencesThis.Contains(name) {
			return
		}
		s := t.GetScope()
		v := s.GetVar(name)
		if v == nil || !v.IsGlobal() {
			return
		}
		// If anything but a function is assigned we assume that possibly
		// a function referencing this is being assigned. Otherwise we
		// check whether the function assigned is a) an arrow function, which has a
		// lexically-scoped this, or b) a non-arrow function that does not reference this.
		if value == nil || !value.IsFunction() || NodeUtil_ReferencesOwnReceiver(value) {
			r.maybeReferencesThis.Add(name)
		}
	}
}

// RescopeGlobalSymbols_RewriteScopeCallback visits each NAME token and checks
// whether it refers to a global variable. If yes, rewrites the name to be a
// property access on the "globalSymbolNamespace". If the NAME is an extern
// variable, it becomes a property access on window.
//
//	var a = 1, b = 2, c = 3;
//
// becomes
//
//	var NS.a = 1, NS.b = 2, NS.c = 4
//
// (The var token is removed in a later traversal.)
//
//	a + b
//
// becomes
//
//	NS.a + NS.b
//
//	a()
//
// becomes
//
//	(0,NS.a)()
//
// Notice the special syntax here to preserve the *this* semantics in the
// function call.
//
//	var {a: b} = {}
//
// becomes
//
//	var {a: NS.b} = {}
//
// (This is invalid syntax, but the VAR token is removed later).
//
// Java: private class RescopeGlobalSymbols.RewriteScopeCallback implements NodeTraversal.Callback
type RescopeGlobalSymbols_RewriteScopeCallback struct {
	outer           *RescopeGlobalSymbols
	preDeclarations *jutil.ArrayList[*RescopeGlobalSymbols_RewriteScopeCallback_ChunkGlobal]
}

// Java: RescopeGlobalSymbols.RewriteScopeCallback.shouldTraverse(NodeTraversal, Node, Node)
func (c *RescopeGlobalSymbols_RewriteScopeCallback) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	if NodeUtil_IsNameDeclaration(n) {
		c.VisitNameDeclaration(t, n)
	}
	return true
}

// Java: RescopeGlobalSymbols.RewriteScopeCallback.visit(NodeTraversal, Node, Node)
func (c *RescopeGlobalSymbols_RewriteScopeCallback) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if n.IsName() {
		c.VisitName(t, n, parent)
	}
}

// Java: RescopeGlobalSymbols.RewriteScopeCallback.visitNameDeclaration(NodeTraversal, Node)
func (c *RescopeGlobalSymbols_RewriteScopeCallback) VisitNameDeclaration(t *NodeTraversal, declaration *rhino.Node) {
	r := c.outer
	var allLhsNodes []*rhino.Node
	NodeUtil_VisitLhsNodesInNode(declaration, func(n *rhino.Node) { allLhsNodes = append(allLhsNodes, n) })
	if len(allLhsNodes) == 0 {
		return
	}
	hasImportantName := false
	isGlobalDeclaration := t.GetScope().GetVar(allLhsNodes[0].GetString()).IsGlobal()

	// Check if any names are in the externs or are global and cross chunk.
	for _, lhs := range allLhsNodes {
		jutil.CheckStateMsg1(lhs.IsName(), "Unexpected lhs node %s, expected NAME", lhs)
		if (isGlobalDeclaration && r.IsCrossChunkName(lhs.GetString())) ||
			r.IsExternVar(lhs.GetString(), t) {
			hasImportantName = true
			break
		}
	}

	if hasImportantName {
		c.RewriteNameDeclaration(t, declaration, allLhsNodes, isGlobalDeclaration)
	}
}

// RewriteNameDeclaration partially rewrites a declaration as an assignment.
//
// In the post traversal, all global, cross-chunk names and extern name
// references will become property accesses. They will then be invalid as the
// lhs of a declaration, so we need to convert them to assignments. We also
// convert any other names or destructuring patterns in the same declaration
// to assignments and add an earlier declaration.
//
// Java: RescopeGlobalSymbols.RewriteScopeCallback.rewriteNameDeclaration(NodeTraversal, Node, List<Node>, boolean)
func (c *RescopeGlobalSymbols_RewriteScopeCallback) RewriteNameDeclaration(
	t *NodeTraversal, declaration *rhino.Node, allLhsNodes []*rhino.Node, isGlobalDeclaration bool) {
	r := c.outer

	// Add predeclarations for variables that are neither global/cross-chunk names nor externs.
	input := t.GetInput()
	for _, lhs := range allLhsNodes {
		name := lhs.GetString()
		if !(isGlobalDeclaration && r.IsCrossChunkName(name)) && !r.IsExternVar(name, t) {
			c.preDeclarations.Add(
				NewRescopeGlobalSymbols_RewriteScopeCallback_ChunkGlobal(input.GetAstRoot(r.compiler), rhino.IR_Name(name).Srcref(lhs)))
		}
	}

	// Convert all names with an rhs and all destructuring patterns to be assignments. e.g.
	//  VAR
	//    NAME foo
	//      NUMBER 3
	// becomes
	//  VAR
	//    ASSIGN
	//      NAME foo
	//      NUMBER 3
	for child := declaration.GetFirstChild(); child != nil; {
		next := child.GetNext()
		if child.IsName() && child.HasChildren() {
			assign := rhino.IR_Assign(child.CloneNode(), child.RemoveFirstChild())
			child.ReplaceWith(assign)
			assign.SetJSDocInfo(declaration.GetJSDocInfo())
		} else if child.IsDestructuringLhs() {
			if child.HasOneChild() {
				if !NodeUtil_IsEnhancedFor(declaration.GetParent()) {
					jutil.CheckStateMsg1(
						false,
						"DESTRUCTURING_LHS should have two children: %s",
						declaration.ToStringTree())
				}
				// remove the DESTRUCTURING_LHS but leave the actual destructuring pattern
				child.ReplaceWith(child.RemoveFirstChild())
			} else {
				assign := rhino.IR_Assign(child.RemoveFirstChild(), child.RemoveFirstChild())
				child.ReplaceWith(assign)
				assign.SetJSDocInfo(declaration.GetJSDocInfo())
			}
		}
		child = next
	}
	r.compiler.ReportChangeToEnclosingScope(declaration)
}

// Java: RescopeGlobalSymbols.RewriteScopeCallback.visitName(NodeTraversal, Node, Node)
func (c *RescopeGlobalSymbols_RewriteScopeCallback) VisitName(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	r := c.outer
	name := n.GetString()
	// Ignore anonymous functions
	if parent.IsFunction() && name == "" {
		return
	}
	if r.IsExternVar(name, t) {
		return
	}
	// When the globalSymbolNamespace is used as a local variable name
	// add suffix to avoid shadowing the namespace. Also add a suffix
	// if a name starts with the name of the globalSymbolNamespace and
	// the suffix.
	var_ := t.GetScope().GetVar(name)
	if !var_.IsGlobal() &&
		(name == r.globalSymbolNamespace ||
			strings.HasPrefix(name, r.globalSymbolNamespace+rescopeGlobalSymbols_DISAMBIGUATION_SUFFIX)) {
		n.SetString(name + rescopeGlobalSymbols_DISAMBIGUATION_SUFFIX)
		r.compiler.ReportChangeToEnclosingScope(n)
	}
	// We only care about global vars.
	if !(var_.IsGlobal() && r.IsCrossChunkName(name)) {
		return
	}
	c.ReplaceSymbol(n, name)
}

// ReplaceSymbol replaces a global cross-chunk name with an access on the
// global namespace symbol.
//
// Java: RescopeGlobalSymbols.RewriteScopeCallback.replaceSymbol(Node, String)
func (c *RescopeGlobalSymbols_RewriteScopeCallback) ReplaceSymbol(node *rhino.Node, name string) {
	r := c.outer
	parent := node.GetParent()
	replacement := rhino.IR_Getprop(rhino.IR_Name(r.globalSymbolNamespace), name)
	replacement.SrcrefTree(node)

	node.ReplaceWith(replacement)
	r.compiler.ReportChangeToEnclosingScope(replacement)
	if parent.IsCall() && !r.maybeReferencesThis.Contains(name) {
		// Do not write calls like this: (0, _a)() but rather as _.a(). The
		// this inside the function will be wrong, but it doesn't matter
		// because the this is never read.
		parent.PutBooleanProp(rhino.Node_FREE_CALL, false)
	}
	r.compiler.ReportChangeToEnclosingScope(parent)
}

// DeclareChunkGlobals adds back declarations for variables that do not cross
// chunk boundaries. Must be called after RemoveGlobalVarCallback.
//
// Java: RescopeGlobalSymbols.RewriteScopeCallback.declareChunkGlobals()
func (c *RescopeGlobalSymbols_RewriteScopeCallback) DeclareChunkGlobals() {
	for global := range c.preDeclarations.All() {
		if global.root.HasChildren() && global.root.GetFirstChild().IsVar() {
			global.root.GetFirstChild().AddChildToBack(global.name)
		} else {
			global.root.AddChildToFront(rhino.IR_Var(global.name).Srcref(global.name))
		}
		c.outer.compiler.ReportChangeToEnclosingScope(global.root)
	}
}

// RescopeGlobalSymbols_RewriteScopeCallback_ChunkGlobal is a variable that
// doesn't cross chunk boundaries.
//
// Java: private static class RescopeGlobalSymbols.RewriteScopeCallback.ChunkGlobal
type RescopeGlobalSymbols_RewriteScopeCallback_ChunkGlobal struct {
	root *rhino.Node
	name *rhino.Node
}

// Java: RescopeGlobalSymbols.RewriteScopeCallback.ChunkGlobal(Node, Node)
func NewRescopeGlobalSymbols_RewriteScopeCallback_ChunkGlobal(root *rhino.Node, name *rhino.Node) *RescopeGlobalSymbols_RewriteScopeCallback_ChunkGlobal {
	return &RescopeGlobalSymbols_RewriteScopeCallback_ChunkGlobal{root: root, name: name}
}

// RescopeGlobalSymbols_RemoveGlobalVarCallback removes every occurrence of
// var/let/const that declares a global variable.
//
//	var NS.a = 1, NS.b = 2;
//
// becomes
//
//	NS.a = 1; NS.b = 2;
//
//	for (var a = 0, b = 0;;)
//
// becomes
//
//	for (NS.a = 0, NS.b = 0;;)
//
// Declarations without assignments are optimized away:
//
//	var a = 1, b;
//
// becomes
//
//	NS.a = 1
//
// Java: private class RescopeGlobalSymbols.RemoveGlobalVarCallback extends AbstractShallowStatementCallback
type RescopeGlobalSymbols_RemoveGlobalVarCallback struct {
	NodeTraversal_AbstractShallowStatementCallbackBase
	outer *RescopeGlobalSymbols
}

// Java: RescopeGlobalSymbols.RemoveGlobalVarCallback()
func NewRescopeGlobalSymbols_RemoveGlobalVarCallback(outer *RescopeGlobalSymbols) *RescopeGlobalSymbols_RemoveGlobalVarCallback {
	c := &RescopeGlobalSymbols_RemoveGlobalVarCallback{outer: outer}
	return c
}

// Java: RescopeGlobalSymbols.RemoveGlobalVarCallback.visit(NodeTraversal, Node, Node)
func (c *RescopeGlobalSymbols_RemoveGlobalVarCallback) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	compiler := c.outer.compiler
	if !NodeUtil_IsNameDeclaration(n) {
		return
	}

	var commas []*rhino.Node
	var interestingChildren []*rhino.Node
	// Filter out declarations without assignments.
	// As opposed to regular var nodes, there are always assignments
	// because the previous traversal in RewriteScopeCallback creates
	// them.
	allNameOrDestructuring := true
	for ch := n.GetFirstChild(); ch != nil; ch = ch.GetNext() {
		if !ch.IsName() && !ch.IsDestructuringLhs() {
			allNameOrDestructuring = false
		}
		if ch.IsAssign() || NodeUtil_IsAnyFor(parent) {
			interestingChildren = append(interestingChildren, ch)
		}
	}
	// If every child of a var declares a name, it must stay in place.
	// This is the case if none of the declared variables cross chunk
	// boundaries.
	if allNameOrDestructuring {
		return
	}
	for _, ch := range interestingChildren {
		if NodeUtil_IsAnyFor(parent) && parent.GetFirstChild() == n {
			commas = append(commas, ch.CloneTree())
		} else {
			// Var statement outside of for-loop.
			expr := rhino.IR_ExprResult(ch.CloneTree()).Srcref(ch)
			NodeUtil_MarkNewScopesChanged(expr, compiler)
			expr.InsertBefore(n)
		}
	}
	if len(commas) != 0 {
		comma := c.JoinOnComma(commas, n)
		comma.InsertBefore(n)
	}
	// Remove the var/const/let node.
	n.Detach()
	NodeUtil_MarkFunctionsDeleted(n, compiler)
	compiler.ReportChangeToEnclosingScope(parent)
}

// Java: RescopeGlobalSymbols.RemoveGlobalVarCallback.joinOnComma(List<Node>, Node)
func (c *RescopeGlobalSymbols_RemoveGlobalVarCallback) JoinOnComma(commas []*rhino.Node, source *rhino.Node) *rhino.Node {
	comma := commas[0]
	for i := 1; i < len(commas); i++ {
		nextComma := rhino.IR_Comma(comma, commas[i])
		nextComma.SrcrefIfMissing(source)
		comma = nextComma
	}
	return comma
}
