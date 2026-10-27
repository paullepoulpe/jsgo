// Port of com/google/javascript/jscomp/RescopeGlobalSymbols.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/ast"
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
	optimizeLocalAccess   CompilerOptions_OptimizeLocalAccess
}

// NewRescopeGlobalSymbols is the constructor for the RescopeGlobalSymbols
// compiler pass.
//
// compiler is the JSCompiler, for reporting code changes;
// globalSymbolNamespace is the name of namespace into which all global
// symbols are transferred. If assumeCrossChunkNames is true, all global
// symbols will be assumed cross chunk boundaries and thus require renaming.
// If optimizeLocalAccess is true, write aliases for global symbols in the
// chunk where they are defined, for more efficient access.
//
// Java: RescopeGlobalSymbols(AbstractCompiler, String, boolean, CompilerOptions.OptimizeLocalAccess)
func NewRescopeGlobalSymbols(
	compiler AbstractCompiler,
	globalSymbolNamespace string,
	assumeCrossChunkNames bool,
	optimizeLocalAccess CompilerOptions_OptimizeLocalAccess) *RescopeGlobalSymbols {
	return NewRescopeGlobalSymbolsWithAddExtern(
		compiler, globalSymbolNamespace, true, assumeCrossChunkNames, optimizeLocalAccess)
}

// NewRescopeGlobalSymbolsWithAddExtern is the constructor for the
// RescopeGlobalSymbols compiler pass for use in testing.
//
// compiler is the JSCompiler, for reporting code changes;
// globalSymbolNamespace is the name of namespace into which all global
// symbols are transferred. If addExtern is true, the compiler will consider
// the globalSymbolNamespace an extern name. If assumeCrossChunkNames is true,
// all global symbols will be assumed cross chunk boundaries and thus require
// renaming. VisibleForTesting. optimizeLocalAccess is the mode for optimizing
// local access to global symbols.
//
// Java: RescopeGlobalSymbols(AbstractCompiler, String, boolean, boolean, CompilerOptions.OptimizeLocalAccess)
func NewRescopeGlobalSymbolsWithAddExtern(
	compiler AbstractCompiler,
	globalSymbolNamespace string,
	addExtern bool,
	assumeCrossChunkNames bool,
	optimizeLocalAccess CompilerOptions_OptimizeLocalAccess) *RescopeGlobalSymbols {
	r := &RescopeGlobalSymbols{}
	r.compiler = compiler
	r.globalSymbolNamespace = globalSymbolNamespace
	r.addExtern = addExtern
	r.assumeCrossChunkNames = assumeCrossChunkNames
	r.optimizeLocalAccess = optimizeLocalAccess
	return r
}

// Java: RescopeGlobalSymbols.addExternForGlobalSymbolNamespace()
func (r *RescopeGlobalSymbols) AddExternForGlobalSymbolNamespace() {
	varNode := rhino.IR_Var(rhino.IR_Name(r.globalSymbolNamespace))
	input := r.compiler.GetSynthesizedExternsInput()
	input.GetAstRoot(r.compiler).AddChildToBack(varNode)
	r.compiler.ReportChangeToEnclosingScope(varNode)
}

// Java: RescopeGlobalSymbols.process(Node, Node)
func (r *RescopeGlobalSymbols) Process(externs rhino.Node, root rhino.Node) {
	// Make the name of the globalSymbolNamespace an extern.
	if r.addExtern {
		r.AddExternForGlobalSymbolNamespace()
	}

	// Rewrite all references to global symbols to properties of a single symbol:

	// Turn global named function statements into var assignments.
	NodeTraversal_Traverse(
		r.compiler,
		root,
		newRescopeGlobalDeclsCallback(r.compiler))

	// Find global names that are used in more than one chunk. Those that
	// are have to be rewritten.
	findCrossChunkNames :=
		newRescopeFindCrossChunkNamesCallback(r.optimizeLocalAccess.ShouldOptimize())

	// And find names that may refer to functions that reference this.
	findNamesReferencingThis := newRescopeFindNamesReferencingThis()

	CombinedCompilerPass_Traverse(
		r.compiler, root, []NodeTraversal_Callback{findCrossChunkNames, findNamesReferencingThis})

	symbolInfo :=
		NewRescopeGlobalSymbolsRewriteCallback_SymbolInformationFromSets(
			findCrossChunkNames.crossChunkNames,
			findCrossChunkNames.crossChunkNamesWithWriteFromOtherChunk,
			findCrossChunkNames.globalNamesWithReadInDefiningChunk,
			findCrossChunkNames.globalNamesWithInnerScopeWriteInDefiningChunk,
			findNamesReferencingThis.maybeReferencesThis)

	// Rewrite all references to be property accesses of the single symbol.
	rewriteScope :=
		NewRescopeGlobalSymbolsRewriteCallback(
			r.compiler,
			r.globalSymbolNamespace,
			r.assumeCrossChunkNames,
			r.optimizeLocalAccess,
			externs,
			symbolInfo)
	NodeTraversal_Traverse(r.compiler, root, rewriteScope)

	// Remove the var from statements in global scope if the declared names have been rewritten
	// in the previous pass.
	NodeTraversal_Traverse(r.compiler, root, newRescopeRemoveGlobalVarCallback(r.compiler))

	rewriteScope.AddDeclarations()
}

// rescopeGlobalDeclsCallback rewrites global function and class declarations to
// var statements + assignment. Ignores non-global function and class
// declarations.
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
// Java: private static class RescopeGlobalSymbols.RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback extends AbstractShallowStatementCallback
type rescopeGlobalDeclsCallback struct {
	NodeTraversal_AbstractShallowStatementCallbackBase
	compiler AbstractCompiler
}

// Java: RescopeGlobalSymbols.RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback(AbstractCompiler)
func newRescopeGlobalDeclsCallback(compiler AbstractCompiler) *rescopeGlobalDeclsCallback {
	return &rescopeGlobalDeclsCallback{compiler: compiler}
}

// Java: RescopeGlobalSymbols.RewriteGlobalClassFunctionDeclarationsToVarAssignmentsCallback.visit(NodeTraversal, Node, Node)
func (c *rescopeGlobalDeclsCallback) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
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
	var name string
	// Remove the class or function name. Anonymous classes have an EMPTY node, while anonymous
	// functions have a NAME node with an empty string.
	if cls, ok := ast.AsClass(n); ok {
		nameNode, _ := cls.Name()
		name = nameNode.Ident()
		cls.RemoveName()
	} else {
		fn := ast.MustFunction(n)
		nameNode := fn.NameNode()
		name = nameNode.Ident()
		fn.RemoveName()
		c.compiler.ReportChangeToEnclosingScope(nameNode.Node())
	}
	prev := n.GetPrevious()
	n.Detach()
	var_ := NodeUtil_NewVarNode(name, n)
	if prev.IsNil() {
		parent.AddChildToFront(var_)
	} else {
		var_.InsertAfter(prev)
	}
	c.compiler.ReportChangeToEnclosingScope(parent)
}

// rescopeFindCrossChunkNamesCallback finds all global names that are used in
// more than one chunk. The following compiler transformations can ignore the
// globals that are not.
//
// Java: private static class RescopeGlobalSymbols.FindCrossChunkNamesCallback extends AbstractPostOrderCallback
type rescopeFindCrossChunkNamesCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	trackLocalAccessSets bool
	crossChunkNames      *jutil.LinkedHashSet[string]

	// Global cross-chunk identifiers that are written to from a chunk other
	// than the defining one. Only populated if trackLocalAccessSets is true.
	crossChunkNamesWithWriteFromOtherChunk *jutil.LinkedHashSet[string]

	// Global identifiers that are read in the chunk where they are defined.
	// Only populated if trackLocalAccessSets is true.
	globalNamesWithReadInDefiningChunk *jutil.LinkedHashSet[string]

	// Global identifiers that are written to from a nested scope (i.e. not
	// the global scope) in their defining chunk. Only populated if
	// trackLocalAccessSets is true.
	globalNamesWithInnerScopeWriteInDefiningChunk *jutil.LinkedHashSet[string]
}

// Java: RescopeGlobalSymbols.FindCrossChunkNamesCallback(boolean)
func newRescopeFindCrossChunkNamesCallback(trackLocalAccessSets bool) *rescopeFindCrossChunkNamesCallback {
	return &rescopeFindCrossChunkNamesCallback{
		trackLocalAccessSets:                          trackLocalAccessSets,
		crossChunkNames:                               jutil.NewLinkedHashSet[string](),
		crossChunkNamesWithWriteFromOtherChunk:        jutil.NewLinkedHashSet[string](),
		globalNamesWithReadInDefiningChunk:            jutil.NewLinkedHashSet[string](),
		globalNamesWithInnerScopeWriteInDefiningChunk: jutil.NewLinkedHashSet[string](),
	}
}

// Java: RescopeGlobalSymbols.FindCrossChunkNamesCallback.visit(NodeTraversal, Node, Node)
func (c *rescopeFindCrossChunkNamesCallback) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if nameNode, ok := ast.AsName(n); ok {
		name := nameNode.Ident()
		if name == "" || (!c.trackLocalAccessSets && c.crossChunkNames.Contains(name)) {
			return
		}
		s := t.GetScope()
		v := s.GetVar(name)
		if v == nil || !v.IsGlobal() {
			return
		}

		input := v.GetInput()
		_, hasInit := nameNode.Init()

		if c.trackLocalAccessSets &&
			!t.InGlobalScope() &&
			(input == nil || input.GetChunk() == t.GetChunk()) &&
			NodeUtil_IsLValue(n) &&
			(!NodeUtil_IsNameDeclaration(parent) || hasInit) {
			c.globalNamesWithInnerScopeWriteInDefiningChunk.Add(name)
		}

		if input == nil {
			// We know nothing. Assume name is used across chunks.
			c.crossChunkNames.Add(name)
			return
		}
		// Compare the chunk where the variable is declared to the current
		// chunk. If they are different, the variable is used across chunks.
		chunk := input.GetChunk()
		if chunk != t.GetChunk() {
			c.crossChunkNames.Add(name)

			// If tracking local access, track names with non-local writes separately. This includes
			// all L-value names except for declarations without children.
			if c.trackLocalAccessSets &&
				NodeUtil_IsLValue(n) &&
				(!NodeUtil_IsNameDeclaration(parent) || hasInit) {
				c.crossChunkNamesWithWriteFromOtherChunk.Add(name)
			}
		} else if c.trackLocalAccessSets && n != v.GetNameNode() {
			c.globalNamesWithReadInDefiningChunk.Add(name)
		}
	}
}

// rescopeFindNamesReferencingThis builds the maybeReferencesThis set of names
// that may reference a function that references this. If the function a name
// references does not reference this it can be called as a method call where
// the this value is not the same as in a normal function call.
//
// Java: private static class RescopeGlobalSymbols.FindNamesReferencingThis extends AbstractPostOrderCallback
type rescopeFindNamesReferencingThis struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	maybeReferencesThis *jutil.LinkedHashSet[string]
}

// Java: RescopeGlobalSymbols.FindNamesReferencingThis()
func newRescopeFindNamesReferencingThis() *rescopeFindNamesReferencingThis {
	return &rescopeFindNamesReferencingThis{
		maybeReferencesThis: jutil.NewLinkedHashSet[string](),
	}
}

// Java: RescopeGlobalSymbols.FindNamesReferencingThis.visit(NodeTraversal, Node, Node)
func (c *rescopeFindNamesReferencingThis) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if nameNode, ok := ast.AsName(n); ok {
		name := nameNode.Ident()
		if name == "" {
			return
		}
		var value rhino.Node
		switch ast.RoleIn(n, parent) {
		case ast.RoleAssignTarget:
			value = ast.MustAssign(parent).Value().Node()
		case ast.RoleNameDeclarationDeclarators:
			if init, hasInit := nameNode.Init(); hasInit {
				value = init.Node()
			}
		case ast.RoleFunctionName:
			value = parent
		}
		if value.IsNil() && !NodeUtil_IsLhsByDestructuring(n) {
			// If n is assigned in a destructuring pattern, don't bother finding its value and just
			// assume it may reference this.
			return
		}
		// We already added this symbol. Done after checks above because those
		// are comparatively cheap.
		if c.maybeReferencesThis.Contains(name) {
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
		if value.IsNil() || !value.IsFunction() || NodeUtil_ReferencesOwnReceiver(value) {
			c.maybeReferencesThis.Add(name)
		}
	}
}

// rescopeRemoveGlobalVarCallback removes every occurrence of var/let/const
// that declares a global variable.
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
// Java: private static class RescopeGlobalSymbols.RemoveGlobalVarCallback extends AbstractShallowStatementCallback
type rescopeRemoveGlobalVarCallback struct {
	NodeTraversal_AbstractShallowStatementCallbackBase
	compiler AbstractCompiler
}

// Java: RescopeGlobalSymbols.RemoveGlobalVarCallback(AbstractCompiler)
func newRescopeRemoveGlobalVarCallback(compiler AbstractCompiler) *rescopeRemoveGlobalVarCallback {
	return &rescopeRemoveGlobalVarCallback{compiler: compiler}
}

// Java: RescopeGlobalSymbols.RemoveGlobalVarCallback.visit(NodeTraversal, Node, Node)
func (c *rescopeRemoveGlobalVarCallback) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if !NodeUtil_IsNameDeclaration(n) {
		return
	}

	var commas []rhino.Node
	var interestingChildren []rhino.Node
	// Filter out declarations without assignments.
	// As opposed to regular var nodes, there are always assignments
	// because the previous traversal in RescopeGlobalSymbolsRewriteCallback creates
	// them.
	allNameOrDestructuring := true
	for ch := n.GetFirstChild(); !ch.IsNil(); ch = ch.GetNext() {
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
			NodeUtil_MarkNewScopesChanged(expr, c.compiler)
			expr.InsertBefore(n)
		}
	}
	if len(commas) != 0 {
		comma := c.JoinOnComma(commas, n)
		comma.InsertBefore(n)
	}
	// Remove the var/const/let node.
	n.Detach()
	NodeUtil_MarkFunctionsDeleted(n, c.compiler)
	c.compiler.ReportChangeToEnclosingScope(parent)
}

// Java: RescopeGlobalSymbols.RemoveGlobalVarCallback.joinOnComma(List<Node>, Node)
func (c *rescopeRemoveGlobalVarCallback) JoinOnComma(commas []rhino.Node, source rhino.Node) rhino.Node {
	comma := commas[0]
	for i := 1; i < len(commas); i++ {
		nextComma := rhino.IR_Comma(comma, commas[i])
		nextComma.SrcrefIfMissing(source)
		comma = nextComma
	}
	return comma
}
