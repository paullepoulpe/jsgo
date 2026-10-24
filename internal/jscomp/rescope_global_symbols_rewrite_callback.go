// Port of com/google/javascript/jscomp/RescopeGlobalSymbolsRewriteCallback.java @ closure-compiler v20261015

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
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// RescopeGlobalSymbolsRewriteCallback is the node traversal callback to
// rewrite global symbols as part of the RescopeGlobalSymbols pass.
//
// Visits each NAME token and checks whether it refers to a global variable.
// If yes, rewrites the name to be a property access on the
// "globalSymbolNamespace". If the NAME is an extern variable, it becomes a
// property access on window.
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
// Java: final class RescopeGlobalSymbolsRewriteCallback implements NodeTraversal.Callback
type RescopeGlobalSymbolsRewriteCallback struct {
	compiler              AbstractCompiler
	globalSymbolNamespace string
	assumeCrossChunkNames bool
	optimizeLocalAccess   CompilerOptions_OptimizeLocalAccess
	externNames           *jutil.LinkedHashSet[string]
	symbolInfo            *RescopeGlobalSymbolsRewriteCallback_SymbolInformation

	preDeclarations []*RescopeGlobalSymbolsRewriteCallback_ChunkGlobal

	// Map from chunks to the set of global symbols to be aliased directly for
	// that chunk (e.g. `var {a} = _`). Only populated if optimizeLocalAccess is
	// an ALL_CHUNKS option. These symbols are used in that chunk, to avoid
	// declaring unnecessary aliases. They can only be reassigned during static
	// execution of the defining chunk. If
	// ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS is set, then reassignable
	// symbols that should be wrapped and aliased will be tracked separately in
	// wrappedReassignableCrossChunkNames.
	localAliasesForUnwrappedCrossChunkNames *jutil.LinkedHashMap[*JSChunk, *jutil.LinkedHashSet[string]]

	// Map from global wrapped reassignable symbols to the chunks that use them.
	// Only populated if optimizeLocalAccess is
	// ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS. This is used to determine
	// which chunks to define the wrapper objects in. These symbols are wrapped
	// in holder objects in that chunk (e.g. `_.a = {}`). Subsequent assignments
	// look like `a._ = value`, where `a` is a local alias for `_.a`). These
	// symbols may be reassigned anywhere.
	chunksUsingWrappedReassignableSymbols *jutil.LinkedHashMap[string, *jutil.LinkedHashSet[*JSChunk]]

	// Global symbols that are reassignable and will be wrapped in an object.
	// This allows for constant tracking and inlining by the JavaScript runtime,
	// and enables safe aliasing in all chunks. Only populated if
	// optimizeLocalAccess is ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS.
	wrappedReassignableCrossChunkNames *jutil.LinkedHashSet[string]
}

// Appended to variables names that conflict with globalSymbolNamespace.
//
// Java: RescopeGlobalSymbolsRewriteCallback.DISAMBIGUATION_SUFFIX
const rescopeGlobalSymbolsRewriteCallback_DISAMBIGUATION_SUFFIX = "$"

// RescopeGlobalSymbolsRewriteCallback_SymbolInformation is the record
// SymbolInformation.
//
// Java: static record RescopeGlobalSymbolsRewriteCallback.SymbolInformation
type RescopeGlobalSymbolsRewriteCallback_SymbolInformation struct {
	crossChunkNames                               *jutil.LinkedHashSet[string]
	crossChunkNamesWithWriteFromOtherChunk        *jutil.LinkedHashSet[string]
	globalNamesWithReadInDefiningChunk            *jutil.LinkedHashSet[string]
	globalNamesWithInnerScopeWriteInDefiningChunk *jutil.LinkedHashSet[string]
	maybeReferencesThis                           *jutil.LinkedHashSet[string]
}

// NewRescopeGlobalSymbolsRewriteCallback_SymbolInformation is the record's
// canonical constructor.
//
// Java: RescopeGlobalSymbolsRewriteCallback.SymbolInformation(ImmutableSet<String>, ImmutableSet<String>, ImmutableSet<String>, ImmutableSet<String>, ImmutableSet<String>)
func NewRescopeGlobalSymbolsRewriteCallback_SymbolInformation(
	crossChunkNames *jutil.LinkedHashSet[string],
	crossChunkNamesWithWriteFromOtherChunk *jutil.LinkedHashSet[string],
	globalNamesWithReadInDefiningChunk *jutil.LinkedHashSet[string],
	globalNamesWithInnerScopeWriteInDefiningChunk *jutil.LinkedHashSet[string],
	maybeReferencesThis *jutil.LinkedHashSet[string]) *RescopeGlobalSymbolsRewriteCallback_SymbolInformation {
	return &RescopeGlobalSymbolsRewriteCallback_SymbolInformation{
		crossChunkNames:                               crossChunkNames,
		crossChunkNamesWithWriteFromOtherChunk:        crossChunkNamesWithWriteFromOtherChunk,
		globalNamesWithReadInDefiningChunk:            globalNamesWithReadInDefiningChunk,
		globalNamesWithInnerScopeWriteInDefiningChunk: globalNamesWithInnerScopeWriteInDefiningChunk,
		maybeReferencesThis:                           maybeReferencesThis,
	}
}

// NewRescopeGlobalSymbolsRewriteCallback_SymbolInformationFromSets copies
// each set (ImmutableSet.copyOf).
//
// Java: RescopeGlobalSymbolsRewriteCallback.SymbolInformation(Set<String>, Set<String>, Set<String>, Set<String>, Set<String>)
func NewRescopeGlobalSymbolsRewriteCallback_SymbolInformationFromSets(
	crossChunkNames *jutil.LinkedHashSet[string],
	crossChunkNamesWithWriteFromOtherChunk *jutil.LinkedHashSet[string],
	globalNamesWithReadInDefiningChunk *jutil.LinkedHashSet[string],
	globalNamesWithInnerScopeWriteInDefiningChunk *jutil.LinkedHashSet[string],
	maybeReferencesThis *jutil.LinkedHashSet[string]) *RescopeGlobalSymbolsRewriteCallback_SymbolInformation {
	copyOf := func(s *jutil.LinkedHashSet[string]) *jutil.LinkedHashSet[string] {
		return jutil.NewLinkedHashSet(s.ToSlice()...)
	}
	return NewRescopeGlobalSymbolsRewriteCallback_SymbolInformation(
		copyOf(crossChunkNames),
		copyOf(crossChunkNamesWithWriteFromOtherChunk),
		copyOf(globalNamesWithReadInDefiningChunk),
		copyOf(globalNamesWithInnerScopeWriteInDefiningChunk),
		copyOf(maybeReferencesThis))
}

// Java: RescopeGlobalSymbolsRewriteCallback.SymbolInformation.crossChunkNames()
func (s *RescopeGlobalSymbolsRewriteCallback_SymbolInformation) CrossChunkNames() *jutil.LinkedHashSet[string] {
	return s.crossChunkNames
}

// Java: RescopeGlobalSymbolsRewriteCallback.SymbolInformation.crossChunkNamesWithWriteFromOtherChunk()
func (s *RescopeGlobalSymbolsRewriteCallback_SymbolInformation) CrossChunkNamesWithWriteFromOtherChunk() *jutil.LinkedHashSet[string] {
	return s.crossChunkNamesWithWriteFromOtherChunk
}

// Java: RescopeGlobalSymbolsRewriteCallback.SymbolInformation.globalNamesWithReadInDefiningChunk()
func (s *RescopeGlobalSymbolsRewriteCallback_SymbolInformation) GlobalNamesWithReadInDefiningChunk() *jutil.LinkedHashSet[string] {
	return s.globalNamesWithReadInDefiningChunk
}

// Java: RescopeGlobalSymbolsRewriteCallback.SymbolInformation.globalNamesWithInnerScopeWriteInDefiningChunk()
func (s *RescopeGlobalSymbolsRewriteCallback_SymbolInformation) GlobalNamesWithInnerScopeWriteInDefiningChunk() *jutil.LinkedHashSet[string] {
	return s.globalNamesWithInnerScopeWriteInDefiningChunk
}

// Java: RescopeGlobalSymbolsRewriteCallback.SymbolInformation.maybeReferencesThis()
func (s *RescopeGlobalSymbolsRewriteCallback_SymbolInformation) MaybeReferencesThis() *jutil.LinkedHashSet[string] {
	return s.maybeReferencesThis
}

// Java: RescopeGlobalSymbolsRewriteCallback(AbstractCompiler, String, boolean, CompilerOptions.OptimizeLocalAccess, Node, SymbolInformation)
func NewRescopeGlobalSymbolsRewriteCallback(
	compiler AbstractCompiler,
	globalSymbolNamespace string,
	assumeCrossChunkNames bool,
	optimizeLocalAccess CompilerOptions_OptimizeLocalAccess,
	externs rhino.Node,
	symbolInfo *RescopeGlobalSymbolsRewriteCallback_SymbolInformation) *RescopeGlobalSymbolsRewriteCallback {
	c := &RescopeGlobalSymbolsRewriteCallback{
		localAliasesForUnwrappedCrossChunkNames: jutil.NewLinkedHashMap[*JSChunk, *jutil.LinkedHashSet[string]](),
		chunksUsingWrappedReassignableSymbols:   jutil.NewLinkedHashMap[string, *jutil.LinkedHashSet[*JSChunk]](),
	}
	c.compiler = compiler
	c.globalSymbolNamespace = globalSymbolNamespace
	c.assumeCrossChunkNames = assumeCrossChunkNames
	c.optimizeLocalAccess = optimizeLocalAccess
	c.externNames = NodeUtil_CollectExternVariableNames(compiler, externs)
	c.symbolInfo = symbolInfo
	c.wrappedReassignableCrossChunkNames = c.CollectWrappedReassignableCrossChunkNames()
	return c
}

// Java: RescopeGlobalSymbolsRewriteCallback.isCrossChunkName(String)
func (c *RescopeGlobalSymbolsRewriteCallback) IsCrossChunkName(name string) bool {
	return c.assumeCrossChunkNames ||
		c.symbolInfo.CrossChunkNames().Contains(name) ||
		c.compiler.GetCodingConvention().IsExportedWithLocal(name /* local= */, false)
}

// Java: RescopeGlobalSymbolsRewriteCallback.isExternVar(String, NodeTraversal)
func (c *RescopeGlobalSymbolsRewriteCallback) IsExternVar(varname string, t *NodeTraversal) bool {
	if varname == "" {
		return false
	}
	v := t.GetScope().GetVar(varname)
	return v == nil ||
		v.IsExtern() ||
		(v.GetScope().(*Scope).IsGlobal() && c.externNames.Contains(varname))
}

// Java: RescopeGlobalSymbolsRewriteCallback.collectWrappedReassignableCrossChunkNames()
func (c *RescopeGlobalSymbolsRewriteCallback) CollectWrappedReassignableCrossChunkNames() *jutil.LinkedHashSet[string] {
	if c.optimizeLocalAccess != CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS {
		return jutil.NewLinkedHashSet[string]()
	}

	builder := jutil.NewLinkedHashSet[string]()
	builder.AddAll(c.symbolInfo.CrossChunkNamesWithWriteFromOtherChunk().ToSlice()...)
	for name := range c.symbolInfo.GlobalNamesWithInnerScopeWriteInDefiningChunk().All() {
		if c.IsCrossChunkName(name) {
			builder.Add(name)
		}
	}
	return builder
}

// Java: RescopeGlobalSymbolsRewriteCallback.shouldTraverse(NodeTraversal, Node, Node)
func (c *RescopeGlobalSymbolsRewriteCallback) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if NodeUtil_IsNameDeclaration(n) {
		c.VisitNameDeclaration(t, n)
	}
	return true
}

// Java: RescopeGlobalSymbolsRewriteCallback.visit(NodeTraversal, Node, Node)
func (c *RescopeGlobalSymbolsRewriteCallback) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsName() {
		c.VisitName(t, n, parent)
	}
}

// Java: RescopeGlobalSymbolsRewriteCallback.visitNameDeclaration(NodeTraversal, Node)
func (c *RescopeGlobalSymbolsRewriteCallback) VisitNameDeclaration(t *NodeTraversal, declaration rhino.Node) {
	var allLhsNodes []rhino.Node
	NodeUtil_VisitLhsNodesInNode(declaration, func(n rhino.Node) { allLhsNodes = append(allLhsNodes, n) })
	if len(allLhsNodes) == 0 {
		return
	}
	hasImportantName := false
	isGlobalDeclaration := t.GetScope().GetVar(allLhsNodes[0].GetString()).IsGlobal()

	// Check if any names are in the externs or are global and cross chunk.
	for _, lhs := range allLhsNodes {
		jutil.CheckStateMsg1(lhs.IsName(), "Unexpected lhs node %s, expected NAME", lhs)
		if (isGlobalDeclaration && c.IsCrossChunkName(lhs.GetString())) ||
			c.IsExternVar(lhs.GetString(), t) {
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
// Java: RescopeGlobalSymbolsRewriteCallback.rewriteNameDeclaration(NodeTraversal, Node, List<Node>, boolean)
func (c *RescopeGlobalSymbolsRewriteCallback) RewriteNameDeclaration(
	t *NodeTraversal, declaration rhino.Node, allLhsNodes []rhino.Node, isGlobalDeclaration bool) {
	input := t.GetInput()

	// Add pre-declarations for all LHS variables that are neither global/cross-chunk names nor
	// externs.
	if c.optimizeLocalAccess != CompilerOptions_OptimizeLocalAccess_DISABLED {
		// If we are optimizing local accesses, we only want to pre-declare variables if the
		// declaration will be converted to assignments/property accesses and removed by
		// RemoveGlobalVarCallback. We don't want to add a new var declaration otherwise to avoid
		// redeclaration errors (e.g. "var a; let a;").
		if c.ContainsRescopedOrInitializedVars(declaration, isGlobalDeclaration) {
			for _, lhs := range allLhsNodes {
				name := lhs.GetString()
				if !c.symbolInfo.CrossChunkNamesWithWriteFromOtherChunk().Contains(name) &&
					!c.wrappedReassignableCrossChunkNames.Contains(name) &&
					!c.IsExternVar(name, t) &&
					(!c.IsCrossChunkName(name) ||
						c.symbolInfo.GlobalNamesWithReadInDefiningChunk().Contains(name)) {
					c.preDeclarations = append(c.preDeclarations,
						NewRescopeGlobalSymbolsRewriteCallback_ChunkGlobal(input.GetAstRoot(c.compiler), rhino.IR_Name(name).Srcref(lhs)))
				}
			}
		}
	} else {
		for _, lhs := range allLhsNodes {
			name := lhs.GetString()
			if !(isGlobalDeclaration && c.IsCrossChunkName(name)) && !c.IsExternVar(name, t) {
				c.preDeclarations = append(c.preDeclarations,
					NewRescopeGlobalSymbolsRewriteCallback_ChunkGlobal(input.GetAstRoot(c.compiler), rhino.IR_Name(name).Srcref(lhs)))
			}
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
	for child := declaration.GetFirstChild(); !child.IsNil(); {
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
	c.compiler.ReportChangeToEnclosingScope(declaration)
}

// ContainsRescopedOrInitializedVars determines whether a variable
// declaration statement contains any variables that will be rescoped
// (replaced with property accesses) or initialized (replaced with
// assignments). Assumes that optimizeLocalAccess is true.
//
// Examples:
//
//   - Rescoped: `var a;` where `a` is a cross-chunk name with a nonlocal
//     write (becomes `_.a`), or `a` is a reassignable cross-chunk name
//     (becomes `a._`).
//   - Initialized: `var a = 1;` (has children) or `var [a] = [1];`
//     (destructuring LHS).
//
// If a declaration contains any such variables, the statement will be
// removed by RemoveGlobalVarCallback, and we must pre-declare any other
// variables in it.
//
// Java: RescopeGlobalSymbolsRewriteCallback.containsRescopedOrInitializedVars(Node, boolean)
func (c *RescopeGlobalSymbolsRewriteCallback) ContainsRescopedOrInitializedVars(declaration rhino.Node, isGlobalDeclaration bool) bool {
	for child := declaration.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		if isGlobalDeclaration &&
			child.IsName() &&
			(c.symbolInfo.CrossChunkNamesWithWriteFromOtherChunk().Contains(child.GetString()) ||
				c.wrappedReassignableCrossChunkNames.Contains(child.GetString())) {
			return true
		}
		if child.HasChildren() || child.IsDestructuringLhs() {
			return true
		}
	}
	return false
}

// Java: RescopeGlobalSymbolsRewriteCallback.visitName(NodeTraversal, Node, Node)
func (c *RescopeGlobalSymbolsRewriteCallback) VisitName(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	name := n.GetString()

	// Ignore anonymous functions
	if parent.IsFunction() && name == "" {
		return
	}

	if c.IsExternVar(name, t) {
		return
	}

	// When the globalSymbolNamespace is used as a local variable name
	// add suffix to avoid shadowing the namespace. Also add a suffix
	// if a name starts with the name of the globalSymbolNamespace and
	// the suffix.
	var_ := t.GetScope().GetVar(name)
	if !var_.IsGlobal() &&
		(name == c.globalSymbolNamespace ||
			strings.HasPrefix(name, c.globalSymbolNamespace+rescopeGlobalSymbolsRewriteCallback_DISAMBIGUATION_SUFFIX)) {
		n.SetString(name + rescopeGlobalSymbolsRewriteCallback_DISAMBIGUATION_SUFFIX)
		c.compiler.ReportChangeToEnclosingScope(n)
	}

	// We only care about global vars.
	if !(var_.IsGlobal() && c.IsCrossChunkName(name)) {
		return
	}

	currentChunk := t.GetChunk()

	if c.optimizeLocalAccess != CompilerOptions_OptimizeLocalAccess_DISABLED {
		definingChunk := var_.GetInput().GetChunk()
		if currentChunk == definingChunk {
			if !c.symbolInfo.CrossChunkNamesWithWriteFromOtherChunk().Contains(name) &&
				!c.wrappedReassignableCrossChunkNames.Contains(name) &&
				c.symbolInfo.GlobalNamesWithReadInDefiningChunk().Contains(name) {
				// If the cross-chunk variable is defined in this chunk, not re-assigned in any other
				// chunk, and not rewritten to wrapper access, then we skip the replacement and keep the
				// local name.
				if NodeUtil_IsLValue(n) && (!NodeUtil_IsNameDeclaration(parent) || n.HasChildren()) {
					// Any assignment needs to be also set on the global namespace symbol.
					c.AddGlobalNamespaceAlias(n, name)
				}
				// Early return to skip replacing the symbol.
				return
			}
		} else if c.optimizeLocalAccess == CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS ||
			c.optimizeLocalAccess ==
				CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS {
			// Variables defined in a different chunk are still safe for local aliasing if all of the
			// following conditions are met.
			// 1. The variable is not written to in any other chunk than in the defining chunk.
			// 2. The variable is not written to in an inner scope in the defining chunk, so all
			//     assignments are guaranteed to happen during the initial execution of the chunk.
			// 3. The current chunk depends on the defining chunk, so all writes are guaranteed to
			//     have occurred already when the current chunk is loaded.
			if !c.symbolInfo.CrossChunkNamesWithWriteFromOtherChunk().Contains(name) &&
				!c.symbolInfo.GlobalNamesWithInnerScopeWriteInDefiningChunk().Contains(name) &&
				c.compiler.GetChunkGraph().DependsOn(currentChunk, definingChunk) {
				c.localAliasesForUnwrappedCrossChunkNames.
					ComputeIfAbsent(currentChunk, func(*JSChunk) *jutil.LinkedHashSet[string] { return jutil.NewLinkedHashSet[string]() }).
					Add(name)
				// Early return to skip replacing the symbol.
				return
			}
		}
	}

	if c.wrappedReassignableCrossChunkNames.Contains(name) {
		// If the symbol is a wrapper for a reassignable symbol, record the chunk using the symbol
		// to ensure the wrapper is initialized before use, and replace the symbol with an access to
		// the wrapper.
		c.chunksUsingWrappedReassignableSymbols.
			ComputeIfAbsent(name, func(string) *jutil.LinkedHashSet[*JSChunk] { return jutil.NewLinkedHashSet[*JSChunk]() }).
			Add(currentChunk)
		replacement := rhino.IR_Getprop(rhino.IR_Name(name), "_")
		c.ReplaceSymbol(n, name, replacement)
	} else {
		// Otherwise replace the symbol with an access on the global namespace symbol.
		replacement := rhino.IR_Getprop(rhino.IR_Name(c.globalSymbolNamespace), name)
		c.ReplaceSymbol(n, name, replacement)
	}
}

// ReplaceSymbol replaces a symbol with an access to the replacement.
//
// Java: RescopeGlobalSymbolsRewriteCallback.replaceSymbol(Node, String, Node)
func (c *RescopeGlobalSymbolsRewriteCallback) ReplaceSymbol(node rhino.Node, name string, replacement rhino.Node) {
	parent := node.GetParent()
	replacement.SrcrefTree(node)
	node.ReplaceWith(replacement)
	c.compiler.ReportChangeToEnclosingScope(replacement)
	if parent.IsCall() && !c.symbolInfo.MaybeReferencesThis().Contains(name) {
		// Do not write calls like this: (0, _a)() but rather as _.a(). The
		// this inside the function will be wrong, but it doesn't matter
		// because the this is never read.
		parent.PutBooleanProp(rhino.Node_FREE_CALL, false)
	}
	c.compiler.ReportChangeToEnclosingScope(parent)
}

// AddGlobalNamespaceAlias aliases the local variable's value to the global
// namespace, so that it can be accessed from other chunks.
//
// For simple assignments, this chains the assignment to avoid an extra
// statement (e.g. converts `a = x` to `a = _.a = x`). For other writes (e.g.
// compound assignments, destructuring, increments), it inserts a separate
// assignment statement after the current one (e.g. appends `_.a = a;`).
//
// Java: RescopeGlobalSymbolsRewriteCallback.addGlobalNamespaceAlias(Node, String)
func (c *RescopeGlobalSymbolsRewriteCallback) AddGlobalNamespaceAlias(n rhino.Node, name string) {
	stmt := NodeUtil_GetEnclosingStatement(n)
	if stmt.IsNil() {
		return
	}
	parent := n.GetParent()
	if parent.IsAssign() && parent.GetFirstChild() == n {
		// Optimized path for simple assignments (e.g. `a = x`).
		// We chain the assignment to the global namespace to reduce output size:
		// `a = _.a = x` instead of `a = x; _.a = a;`.
		rhs := parent.GetSecondChild()
		globalNamespaceAssign :=
			rhino.IR_Assign(rhino.IR_Getprop(rhino.IR_Name(c.globalSymbolNamespace), name), rhs.Detach())
		globalNamespaceAssign.SrcrefTree(n)
		parent.AddChildToBack(globalNamespaceAssign)
		c.compiler.ReportChangeToEnclosingScope(parent)
	} else {
		// Fallback path for other writes (e.g. `a++`, `a += x`, destructuring).
		// In these cases, we cannot easily chain the assignment without changing semantics
		// or making the code overly complex. Instead, we append a separate assignment
		// statement (`_.a = a;`).
		globalNamespaceAssign :=
			rhino.IR_Assign(rhino.IR_Getprop(rhino.IR_Name(c.globalSymbolNamespace), name), rhino.IR_Name(name))
		aliasStatement := rhino.IR_ExprResult(globalNamespaceAssign)
		aliasStatement.SrcrefTree(n)
		aliasStatement.InsertAfter(stmt)
		c.compiler.ReportChangeToEnclosingScope(stmt.GetParent())
	}
}

// AddDeclarations adds back declarations for variables that do not cross
// chunk boundaries, and declares local aliases and wrappers. Must be called
// after RemoveGlobalVarCallback.
//
// Java: RescopeGlobalSymbolsRewriteCallback.addDeclarations()
func (c *RescopeGlobalSymbolsRewriteCallback) AddDeclarations() {
	c.DeclareChunkGlobals()
	if c.optimizeLocalAccess == CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS ||
		c.optimizeLocalAccess ==
			CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS {
		c.DeclareLocalAliasesAndWrappers()
	}
}

// DeclareChunkGlobals adds back declarations for variables that do not cross
// chunk boundaries. Must be called after RemoveGlobalVarCallback.
//
// Java: RescopeGlobalSymbolsRewriteCallback.declareChunkGlobals()
func (c *RescopeGlobalSymbolsRewriteCallback) DeclareChunkGlobals() {
	for _, global := range c.preDeclarations {
		if global.root.HasChildren() && global.root.GetFirstChild().IsVar() {
			global.root.GetFirstChild().AddChildToBack(global.name)
		} else {
			global.root.AddChildToFront(rhino.IR_Var(global.name).Srcref(global.name))
		}
		c.compiler.ReportChangeToEnclosingScope(global.root)
	}
}

// DeclareLocalAliasesAndWrappers declares local aliases (e.g. `var {a} = _;`
// or `var a = _.a;`) and wrapper objects (e.g. `_.a = {};` or
// `var a = _.a = {};`) at the beginning of chunks for cross-chunk variables
// that can be safely accessed locally. Must be called after
// RemoveGlobalVarCallback. Only applies if optimizeLocalAccess is an
// ALL_CHUNKS option.
//
// Java: RescopeGlobalSymbolsRewriteCallback.declareLocalAliasesAndWrappers()
func (c *RescopeGlobalSymbolsRewriteCallback) DeclareLocalAliasesAndWrappers() {
	chunkGraph := c.compiler.GetChunkGraph()
	chunks := chunkGraph.GetAllChunks()
	if len(chunks) == 0 {
		// Iterator.next() on an empty Iterable.
		panic(jutil.NewNoSuchElementException())
	}
	defaultRootChunk := chunks[0]

	// Compute the chunk in which to define the wrapper for each wrapped reassignable symbol.
	wrapperAssignmentChunks := jutil.NewLinkedHashMap[string, *JSChunk]()
	if !c.wrappedReassignableCrossChunkNames.IsEmpty() {
		for name := range c.wrappedReassignableCrossChunkNames.All() {
			usingChunks :=
				c.chunksUsingWrappedReassignableSymbols.GetOrDefault(name, jutil.NewLinkedHashSet[*JSChunk]())
			var wrapperAssignmentChunk *JSChunk
			if usingChunks.IsEmpty() {
				wrapperAssignmentChunk = defaultRootChunk
			} else {
				wrapperAssignmentChunk = chunkGraph.GetDeepestCommonDependencyInclusiveOfCollection(usingChunks.ToSlice())
			}
			wrapperAssignmentChunks.Put(name, wrapperAssignmentChunk)
		}
	}

	// Prepend local aliases and wrapper assignments to each chunk.
	for _, chunk := range chunkGraph.GetAllChunks() {
		c.DeclareLocalAliasesAndWrappersForChunk(chunk, wrapperAssignmentChunks)
	}
}

// DeclareLocalAliasesAndWrappersForChunk declares all local aliases and
// wrapper assignments that need to be prepended to the given chunk.
//
// Java: RescopeGlobalSymbolsRewriteCallback.declareLocalAliasesAndWrappersForChunk(JSChunk, Map<String, JSChunk>)
func (c *RescopeGlobalSymbolsRewriteCallback) DeclareLocalAliasesAndWrappersForChunk(
	chunk *JSChunk, wrapperAssignmentChunks *jutil.LinkedHashMap[string, *JSChunk]) {
	inputs := chunk.GetInputs()
	if len(inputs) == 0 {
		return
	}
	script := inputs[0].GetAstRoot(c.compiler)
	insertionPoint := script.GetFirstChild()
	changed := false

	wrappedReassignableLocalAliases := jutil.NewLinkedHashSet[string]()
	if c.optimizeLocalAccess ==
		CompilerOptions_OptimizeLocalAccess_ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS {
		wrappedReassignableLocalAliases = jutil.NewLinkedHashSet[string]()

		for name := range c.wrappedReassignableCrossChunkNames.All() {
			assignmentChunk := wrapperAssignmentChunks.Get(name)
			usingChunks := c.chunksUsingWrappedReassignableSymbols.Get(name)
			if chunk == assignmentChunk {
				// Define the wrapper object for the reassignable symbol in this chunk.
				if usingChunks != nil && usingChunks.Contains(chunk) {
					// The symbol is used in this chunk, so we also declare a local alias:
					// `var name = _.name = {};`.
					rescopeGlobalSymbolsRewriteCallback_addDeclaration(
						script,
						insertionPoint,
						rhino.IR_VarWithValue(
							rhino.IR_Name(name),
							rhino.IR_Assign(rhino.IR_Getprop(rhino.IR_Name(c.globalSymbolNamespace), name), rhino.IR_Objectlit())))
				} else {
					// The symbol is not used in this chunk, so we only define the wrapper object:
					// `_.name = {};`.
					rescopeGlobalSymbolsRewriteCallback_addDeclaration(
						script,
						insertionPoint,
						rhino.IR_ExprResult(
							rhino.IR_Assign(rhino.IR_Getprop(rhino.IR_Name(c.globalSymbolNamespace), name), rhino.IR_Objectlit())))
				}
				changed = true
			} else if usingChunks != nil && usingChunks.Contains(chunk) {
				// If the wrapper is defined in a different chunk, but the symbol is used in this
				// chunk, then only declare a local alias.
				wrappedReassignableLocalAliases.Add(name)
			}
		}
	}

	// Combine all local aliases into a single destructuring assignment
	// (e.g. `var {a, b, c} = _;`) if output is ES2015+, otherwise declare
	// individual local aliases (e.g. `var a = _.a;`).
	//
	// Go-only: language-adapter — Iterables.concat of the two sets, as a
	// slice.
	localAliases :=
		append(
			c.localAliasesForUnwrappedCrossChunkNames.GetOrDefault(chunk, jutil.NewLinkedHashSet[string]()).ToSlice(),
			wrappedReassignableLocalAliases.ToSlice()...)
	if len(localAliases) != 0 {
		if c.compiler.GetOptions().GetOutputFeatureSet().ContainsFeature(parser.FeatureSet_Feature_OBJECT_DESTRUCTURING) {
			objectPattern := rhino.IR_ObjectPattern()
			for _, name := range localAliases {
				stringKey := rhino.IR_StringKeyWithValue(name, rhino.IR_Name(name))
				stringKey.SetShorthandProperty(true)
				objectPattern.AddChildToBack(stringKey)
			}
			rescopeGlobalSymbolsRewriteCallback_addDeclaration(
				script, insertionPoint, rhino.IR_VarWithValue(objectPattern, rhino.IR_Name(c.globalSymbolNamespace)))
			NodeUtil_AddFeatureToScript(script, parser.FeatureSet_Feature_OBJECT_DESTRUCTURING, c.compiler)
		} else {
			for _, name := range localAliases {
				rescopeGlobalSymbolsRewriteCallback_addDeclaration(
					script,
					insertionPoint,
					rhino.IR_VarWithValue(rhino.IR_Name(name), rhino.IR_Getprop(rhino.IR_Name(c.globalSymbolNamespace), name)))
			}
		}
		changed = true
	}

	if changed {
		c.compiler.ReportChangeToEnclosingScope(script)
	}
}

// Java: RescopeGlobalSymbolsRewriteCallback.addDeclaration(Node, Node, Node)
func rescopeGlobalSymbolsRewriteCallback_addDeclaration(script rhino.Node, insertionPoint rhino.Node, decl rhino.Node) {
	decl.SrcrefTree(script)
	if insertionPoint.IsNil() {
		script.AddChildToBack(decl)
	} else {
		decl.InsertBefore(insertionPoint)
	}
}

// RescopeGlobalSymbolsRewriteCallback_ChunkGlobal is a variable that doesn't
// cross chunk boundaries.
//
// Java: private static class RescopeGlobalSymbolsRewriteCallback.ChunkGlobal
type RescopeGlobalSymbolsRewriteCallback_ChunkGlobal struct {
	root rhino.Node
	name rhino.Node
}

// Java: RescopeGlobalSymbolsRewriteCallback.ChunkGlobal(Node, Node)
func NewRescopeGlobalSymbolsRewriteCallback_ChunkGlobal(root rhino.Node, name rhino.Node) *RescopeGlobalSymbolsRewriteCallback_ChunkGlobal {
	return &RescopeGlobalSymbolsRewriteCallback_ChunkGlobal{root: root, name: name}
}
