// Port of com/google/javascript/jscomp/RescopeGlobalSymbolsRewriteCallback.java @ closure-compiler v20261104

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

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

type rescopeChunkGlobal struct {
	root rhino.Node
	name rhino.Node
}

// rescopeRewriteCallback is the node traversal callback to rewrite global
// symbols as part of the RescopeGlobalSymbols pass.
//
// Java: RescopeGlobalSymbolsRewriteCallback
type rescopeRewriteCallback struct {
	compiler            AbstractCompiler
	optimizeLocalAccess CompilerOptions_OptimizeLocalAccess
	symbolInfo          *rescopeSymbolInfo
	renameHelper        rescopeRenameHelper

	preDeclarations []rescopeChunkGlobal
	lhsScratch      []rhino.Node

	localAliasesForUnwrappedCrossChunkNames map[*JSChunk]*OrderedSet[string]
	chunksUsingCrossChunkSymbols            map[string]*OrderedSet[*JSChunk]
	chunksUsingWrappedReassignableSymbols   map[string]*OrderedSet[*JSChunk]
	chunksRequiringWrapperInitialization    map[string]*OrderedSet[*JSChunk]
}

var (
	_ nodeTraversal_visitTokens    = (*rescopeRewriteCallback)(nil)
	_ nodeTraversal_onlyNameLeaves = (*rescopeRewriteCallback)(nil)
)

func (*rescopeRewriteCallback) visitTokens() *ast.TokenSet { return &rescopeNameVisitTokens }
func (*rescopeRewriteCallback) onlyNameLeaves()            {}

func newRescopeRewriteCallback(
	compiler AbstractCompiler,
	optimizeLocalAccess CompilerOptions_OptimizeLocalAccess,
	symbolInfo *rescopeSymbolInfo,
	renameHelper rescopeRenameHelper,
) *rescopeRewriteCallback {
	return &rescopeRewriteCallback{
		compiler:            compiler,
		optimizeLocalAccess: optimizeLocalAccess,
		symbolInfo:          symbolInfo,
		renameHelper:        renameHelper,
	}
}

func (c *rescopeRewriteCallback) isCrossChunkName(name string) bool {
	return c.symbolInfo.crossChunkNames.Has(name)
}

func (c *rescopeRewriteCallback) isExternVar(id strpool.ID, scope *Scope) bool {
	if id <= strpool.EmptyID {
		return false
	}
	return c.isExternVarResolved(id, scope.GetVarByID(uint32(id)))
}

func (c *rescopeRewriteCallback) isExternVarResolved(id strpool.ID, v *Var) bool {
	if id <= strpool.EmptyID {
		return false
	}
	if v == nil || v.IsExtern() {
		return true
	}
	if v.IsGlobal() {
		_, ok := c.symbolInfo.externNameIDs[id]
		return ok
	}
	return false
}

func (c *rescopeRewriteCallback) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if NodeUtil_IsNameDeclaration(n) {
		c.visitNameDeclaration(t, n)
	}
	return true
}

func (c *rescopeRewriteCallback) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsName() {
		c.visitName(t, n, parent)
	}
}

func (c *rescopeRewriteCallback) visitNameDeclaration(t *NodeTraversal, declaration rhino.Node) {
	c.lhsScratch = c.lhsScratch[:0]
	NodeUtil_VisitLhsNodesInNode(declaration, func(n rhino.Node) { c.lhsScratch = append(c.lhsScratch, n) })
	if len(c.lhsScratch) == 0 {
		return
	}
	hasImportantName := false
	scope := t.GetScope()
	isGlobalDeclaration := scope.GetVarByID(c.lhsScratch[0].GetStringID()).IsGlobal()

	for _, lhs := range c.lhsScratch {
		jutil.CheckStateMsg1(lhs.IsName(), "Unexpected lhs node %s, expected NAME", lhs)
		name := lhs.GetString()
		if (isGlobalDeclaration && c.isCrossChunkName(name)) || c.isExternVar(nodeStringID(lhs), scope) {
			hasImportantName = true
			break
		}
	}

	if hasImportantName {
		c.rewriteNameDeclaration(t, scope, declaration, c.lhsScratch, isGlobalDeclaration)
	}
}

func (c *rescopeRewriteCallback) globalVarStaysLocalInDefiningChunk(id strpool.ID, name string, scope *Scope) bool {
	v := scope.GetVarByID(uint32(id))
	jutil.CheckStateMsg1(v != nil && v.IsGlobal(), "Expected global variable: %s", name)
	if c.isExternVarResolved(id, v) {
		return false
	}
	if !c.isCrossChunkName(name) {
		return true
	}
	if !c.optimizeLocalAccess.ShouldOptimize() {
		return false
	}
	if c.optimizeLocalAccess.ShouldWrapReassignableSymbols() &&
		c.symbolInfo.wrappedCrossChunkNames.Has(name) {
		return c.renameHelper.hasUnwrappedLocalNameInDefiningChunk(name)
	}
	_, readInDefining := c.symbolInfo.globalNamesWithReadInDefiningChunk[name]
	return !c.symbolInfo.crossChunkNamesWithWriteFromOtherChunk.Has(name) && readInDefining
}

func (c *rescopeRewriteCallback) rewriteNameDeclaration(
	t *NodeTraversal, scope *Scope, declaration rhino.Node, allLhsNodes []rhino.Node, isGlobalDeclaration bool) {
	input := t.GetInput()
	var astRoot rhino.Node
	maybeAddPreDeclaration := func(lhs rhino.Node) {
		id := nodeStringID(lhs)
		name := lhs.GetString()
		if c.globalVarStaysLocalInDefiningChunk(id, name, scope) {
			localName := c.renameHelper.getLocalNameInDefiningChunk(name, isGlobalDeclaration)
			if astRoot.IsNil() {
				astRoot = input.GetAstRoot(c.compiler)
			}
			c.preDeclarations = append(c.preDeclarations,
				rescopeChunkGlobal{root: astRoot, name: rhino.IR_Name(localName).Srcref(lhs)})
		}
	}

	if c.optimizeLocalAccess.ShouldPreserveLocalVarDeclarations() &&
		!NodeUtil_IsAnyFor(declaration.GetParent()) {
		changed := false
		if isGlobalDeclaration &&
			!declaration.IsVar() &&
			c.containsNonLocalOrInitializedVars(scope, declaration, isGlobalDeclaration) {
			declaration.SetToken(rhino.Token_VAR)
			changed = true
		}
		for decl := range ast.ReadAhead(ast.MustNameDeclaration(declaration).FirstDeclarator()) {
			if nameDecl, ok := decl.AsName(); ok {
				id := nodeStringID(nameDecl.Node())
				if !c.globalVarStaysLocalInDefiningChunk(id, nameDecl.Ident(), scope) {
					if init, hasInit := nameDecl.Init(); hasInit {
						assign := rhino.IR_Assign(nameDecl.Node().CloneNode(), init.Detach().Node())
						nameDecl.Node().ReplaceWith(assign)
						assign.SetJSDocInfo(declaration.GetJSDocInfo())
						changed = true
					} else if c.isExternVar(id, scope) {
						nameDecl.Detach()
						changed = true
					}
				}
			} else if dlhs, ok := decl.AsDestructuringLhs(); ok {
				NodeUtil_VisitLhsNodesInNode(dlhs.Node(), maybeAddPreDeclaration)
				pattern := dlhs.Pattern().Node()
				init, hasInit := dlhs.Init()
				if !hasInit {
					panic(jutil.NewIllegalStateExceptionMsg(
						"DESTRUCTURING_LHS should have two children: " + declaration.ToStringTree()))
				}
				rhs := init.Node()
				pattern.Detach()
				rhs.Detach()
				assign := rhino.IR_Assign(pattern, rhs)
				dlhs.Node().ReplaceWith(assign)
				assign.SetJSDocInfo(declaration.GetJSDocInfo())
				changed = true
			}
		}
		if changed {
			c.compiler.ReportChangeToEnclosingScope(declaration)
		}
		return
	}

	if !c.optimizeLocalAccess.ShouldOptimize() ||
		c.containsNonLocalOrInitializedVars(scope, declaration, isGlobalDeclaration) {
		for _, lhs := range allLhsNodes {
			maybeAddPreDeclaration(lhs)
		}
	}

	for decl := range ast.ReadAhead(ast.MustNameDeclaration(declaration).FirstDeclarator()) {
		if nameDecl, ok := decl.AsName(); ok {
			if init, hasInit := nameDecl.Init(); hasInit {
				assign := rhino.IR_Assign(nameDecl.Node().CloneNode(), init.Detach().Node())
				nameDecl.Node().ReplaceWith(assign)
				assign.SetJSDocInfo(declaration.GetJSDocInfo())
			}
		} else if dlhs, ok := decl.AsDestructuringLhs(); ok {
			pattern := dlhs.Pattern().Node()
			init, hasInit := dlhs.Init()
			if !hasInit {
				if !NodeUtil_IsEnhancedFor(declaration.GetParent()) {
					panic(jutil.NewIllegalStateExceptionMsg(
						"DESTRUCTURING_LHS should have two children: " + declaration.ToStringTree()))
				}
				dlhs.Node().ReplaceWith(pattern.Detach())
			} else {
				rhs := init.Node()
				pattern.Detach()
				rhs.Detach()
				assign := rhino.IR_Assign(pattern, rhs)
				dlhs.Node().ReplaceWith(assign)
				assign.SetJSDocInfo(declaration.GetJSDocInfo())
			}
		}
	}
	c.compiler.ReportChangeToEnclosingScope(declaration)
}

func (c *rescopeRewriteCallback) containsNonLocalOrInitializedVars(
	scope *Scope, declaration rhino.Node, isGlobalDeclaration bool) bool {
	for child := declaration.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		if isGlobalDeclaration &&
			child.IsName() &&
			!c.globalVarStaysLocalInDefiningChunk(nodeStringID(child), child.GetString(), scope) {
			return true
		}
		if child.HasChildren() || child.IsDestructuringLhs() {
			return true
		}
	}
	return false
}

func (c *rescopeRewriteCallback) ensureLocalAliasesForChunk(chunk *JSChunk) *OrderedSet[string] {
	if s, ok := c.localAliasesForUnwrappedCrossChunkNames[chunk]; ok {
		return s
	}
	if c.localAliasesForUnwrappedCrossChunkNames == nil {
		c.localAliasesForUnwrappedCrossChunkNames = make(map[*JSChunk]*OrderedSet[string])
	}
	s := &OrderedSet[string]{}
	c.localAliasesForUnwrappedCrossChunkNames[chunk] = s
	return s
}

func ensureChunkSetForSymbol(m *map[string]*OrderedSet[*JSChunk], name string) *OrderedSet[*JSChunk] {
	if *m != nil {
		if s, ok := (*m)[name]; ok {
			return s
		}
	} else {
		*m = make(map[string]*OrderedSet[*JSChunk])
	}
	s := &OrderedSet[*JSChunk]{}
	(*m)[name] = s
	return s
}

func (c *rescopeRewriteCallback) visitName(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	id := nodeStringID(n)
	if parent.IsFunction() && id <= strpool.EmptyID {
		return
	}
	scope := t.GetScope()
	var var_ *Var
	if id > strpool.EmptyID {
		var_ = scope.GetVarByID(uint32(id))
		if c.isExternVarResolved(id, var_) {
			return
		}
	}

	name := n.GetString()
	c.renameHelper.maybeDisambiguateLocalVar(n, var_)

	if !(var_.IsGlobal() && c.isCrossChunkName(name)) {
		return
	}

	currentChunk := t.GetChunk()
	if c.optimizeLocalAccess.ShouldWrapReassignableSymbols() {
		ensureChunkSetForSymbol(&c.chunksUsingCrossChunkSymbols, name).Add(currentChunk)
	}

	if c.optimizeLocalAccess.ShouldOptimize() &&
		!c.symbolInfo.crossChunkNamesWithWriteFromOtherChunk.Has(name) {
		definingChunk := var_.GetInput().GetChunk()
		if currentChunk == definingChunk {
			if c.globalVarStaysLocalInDefiningChunk(id, name, scope) {
				localName := c.renameHelper.maybeRenameGlobalInDefiningChunk(n)
				if NodeUtil_IsLValueWrite(n) {
					c.addGlobalNamespaceAlias(n, name, localName, currentChunk)
				}
				return
			}
			if !c.optimizeLocalAccess.ShouldWrapReassignableSymbols() ||
				!c.symbolInfo.wrappedCrossChunkNames.Has(name) {
				replacement := rhino.IR_Getprop(rhino.IR_Name(c.symbolInfo.globalSymbolNamespace), name)
				c.replaceSymbol(n, name, replacement)
				return
			}
		} else if c.optimizeLocalAccess.ShouldOptimizeAllChunks() &&
			(!c.optimizeLocalAccess.ShouldWrapReassignableSymbols() ||
				!c.symbolInfo.wrappedCrossChunkNames.Has(name)) {
			if !c.symbolInfo.globalNamesWithInnerScopeWriteInDefiningChunk.Has(name) &&
				c.compiler.GetChunkGraph().DependsOn(currentChunk, definingChunk) {
				c.ensureLocalAliasesForChunk(currentChunk).Add(name)
				return
			}
		}
	}

	if !c.optimizeLocalAccess.ShouldWrapReassignableSymbols() {
		replacement := rhino.IR_Getprop(rhino.IR_Name(c.symbolInfo.globalSymbolNamespace), name)
		c.replaceSymbol(n, name, replacement)
		return
	}

	jutil.CheckStateMsg1(
		c.symbolInfo.wrappedCrossChunkNames.Has(name),
		"Unexpected unaliased cross-chunk symbol under ALL_CHUNKS_WITH_WRAPPED_REASSIGNABLE_SYMBOLS: %s",
		name)
	ensureChunkSetForSymbol(&c.chunksUsingWrappedReassignableSymbols, name).Add(currentChunk)
	ensureChunkSetForSymbol(&c.chunksRequiringWrapperInitialization, name).Add(currentChunk)
	replacement := rhino.IR_Getprop(rhino.IR_Name(name), "_")
	c.replaceSymbol(n, name, replacement)
}

func (c *rescopeRewriteCallback) replaceSymbol(node rhino.Node, name string, replacement rhino.Node) {
	parent := node.GetParent()
	replacement.SrcrefTree(node)
	node.ReplaceWith(replacement)
	c.compiler.ReportChangeToEnclosingScope(replacement)
	if call, ok := ast.AsCall(parent); ok {
		if _, maybeThis := c.symbolInfo.maybeReferencesThis[name]; !maybeThis {
			call.SetFreeCall(false)
		}
	}
	c.compiler.ReportChangeToEnclosingScope(parent)
}

func (c *rescopeRewriteCallback) addGlobalNamespaceAlias(
	n rhino.Node, name string, localName string, currentChunk *JSChunk) {
	stmt := NodeUtil_GetEnclosingStatement(n)
	if stmt.IsNil() {
		return
	}
	globalTarget := rhino.IR_Getprop(rhino.IR_Name(c.symbolInfo.globalSymbolNamespace), name)
	if c.optimizeLocalAccess.ShouldWrapReassignableSymbols() &&
		c.symbolInfo.wrappedCrossChunkNames.Has(name) {
		ensureChunkSetForSymbol(&c.chunksRequiringWrapperInitialization, name).Add(currentChunk)
		globalTarget = rhino.IR_Getprop(globalTarget, "_")
	}
	parent := n.GetParent()
	if assign, ok := ast.AsAssign(parent); ok && assign.Target().Is(n) {
		rhs := assign.Value().Detach().Node()
		globalNamespaceAssign := rhino.IR_Assign(globalTarget, rhs)
		globalNamespaceAssign.SrcrefTree(n)
		parent.AddChildToBack(globalNamespaceAssign)
		c.compiler.ReportChangeToEnclosingScope(parent)
	} else if NodeUtil_IsNameDeclaration(parent) {
		jutil.CheckStateMsg1(n.HasChildren(), "Expected initialized declaration for %s", n)
		init, _ := ast.MustName(n).Init()
		globalNamespaceAssign := rhino.IR_Assign(globalTarget, init.Detach().Node())
		globalNamespaceAssign.SrcrefTree(n)
		n.AddChildToBack(globalNamespaceAssign)
		c.compiler.ReportChangeToEnclosingScope(n)
	} else if c.optimizeLocalAccess.ShouldRewriteCompoundAndIncrementWritesInPlace() &&
		NodeUtil_IsCompoundAssignmentOp(parent) &&
		parent.GetFirstChild() == n {
		placeholder := rhino.IR_Empty()
		parent.ReplaceWith(placeholder)
		globalNamespaceAssign := rhino.IR_Assign(globalTarget.SrcrefTree(n), parent).Srcref(n)
		placeholder.ReplaceWith(globalNamespaceAssign)
		c.compiler.ReportChangeToEnclosingScope(globalNamespaceAssign)
	} else if c.optimizeLocalAccess.ShouldRewriteCompoundAndIncrementWritesInPlace() &&
		(parent.IsInc() || parent.IsDec()) {
		if !parent.GetBooleanProp(rhino.Node_INCRDECR_PROP) {
			placeholder := rhino.IR_Empty()
			parent.ReplaceWith(placeholder)
			globalNamespaceAssign := rhino.IR_Assign(globalTarget.SrcrefTree(n), parent).Srcref(n)
			placeholder.ReplaceWith(globalNamespaceAssign)
			c.compiler.ReportChangeToEnclosingScope(globalNamespaceAssign)
		} else {
			placeholder := rhino.IR_Empty()
			parent.ReplaceWith(placeholder)
			var globalUpdate rhino.Node
			if parent.IsInc() {
				globalUpdate = rhino.IR_Inc(globalTarget, false).SrcrefTree(n)
			} else {
				globalUpdate = rhino.IR_Dec(globalTarget, false).SrcrefTree(n)
			}
			comma := rhino.IR_Comma(globalUpdate, parent).Srcref(n)
			placeholder.ReplaceWith(comma)
			c.compiler.ReportChangeToEnclosingScope(comma)
		}
	} else if !c.optimizeLocalAccess.ShouldRewriteCompoundAndIncrementWritesInPlace() ||
		rescopeIsSafeDestructuringWrite(n, stmt) {
		globalNamespaceAssign := rhino.IR_Assign(globalTarget, rhino.IR_Name(localName))
		aliasStatement := rhino.IR_ExprResult(globalNamespaceAssign)
		aliasStatement.SrcrefTree(n)
		aliasStatement.InsertAfter(stmt)
		c.compiler.ReportChangeToEnclosingScope(stmt.GetParent())
	} else {
		panic(jutil.NewIllegalStateExceptionMsg("Unexpected unsafe write to local alias: " + n.String()))
	}
}

func rescopeIsSafeDestructuringWrite(n rhino.Node, stmt rhino.Node) bool {
	if !NodeUtil_IsLhsByDestructuring(n) || !NodeUtil_IsStatementBlock(stmt.GetParent()) {
		return false
	}
	rootTarget := NodeUtil_GetRootTarget(n)
	rootTargetParent := rootTarget.GetParent()
	if !rootTargetParent.IsAssign() || rootTargetParent.GetParent() != stmt {
		return false
	}
	return stmt.IsExprResult() || NodeUtil_IsNameDeclaration(stmt)
}

func (c *rescopeRewriteCallback) addDeclarations() {
	c.declareChunkGlobals()
	if c.optimizeLocalAccess.ShouldOptimizeAllChunks() {
		c.declareLocalAliasesAndWrappers()
	}
}

func (c *rescopeRewriteCallback) declareChunkGlobals() {
	for i := range c.preDeclarations {
		global := &c.preDeclarations[i]
		if first, ok := ast.MustScript(global.root).FirstItem(); ok && first.Token() == rhino.Token_VAR {
			first.Node().AddChildToBack(global.name)
		} else {
			global.root.AddChildToFront(rhino.IR_Var(global.name).Srcref(global.name))
		}
		c.compiler.ReportChangeToEnclosingScope(global.root)
	}
}

func (c *rescopeRewriteCallback) declareLocalAliasesAndWrappers() {
	chunkGraph := c.compiler.GetChunkGraph()
	chunks := chunkGraph.GetAllChunks()
	if len(chunks) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	defaultRootChunk := chunks[0]

	var wrapperAssignmentChunks map[string]*JSChunk
	if c.optimizeLocalAccess.ShouldWrapReassignableSymbols() && !c.symbolInfo.wrappedCrossChunkNames.IsEmpty() {
		wrapperAssignmentChunks = make(map[string]*JSChunk, c.symbolInfo.wrappedCrossChunkNames.Len())
		for name := range c.symbolInfo.wrappedCrossChunkNames.All() {
			chunksRequiringWrapper := c.chunksRequiringWrapperInitialization[name]
			if chunksRequiringWrapper.IsEmpty() {
				wrapperAssignmentChunks[name] = defaultRootChunk
			} else {
				wrapperAssignmentChunks[name] = chunkGraph.GetDeepestCommonDependencyInclusiveOfCollection(chunksRequiringWrapper.ToSlice())
			}
		}
	}

	cleanupHelperName := c.getCleanupHelperPropertyName()
	needsCleanupHelper := false
	var allLocalAliasesScratch []string
	var aliasesToKeepScratch []string
	var aliasesToDeleteScratch []string
	for _, chunk := range chunks {
		var usedCleanup bool
		allLocalAliasesScratch, aliasesToKeepScratch, aliasesToDeleteScratch, usedCleanup = c.declareLocalAliasesAndWrappersForChunk(
			chunk,
			wrapperAssignmentChunks,
			cleanupHelperName,
			allLocalAliasesScratch[:0],
			aliasesToKeepScratch[:0],
			aliasesToDeleteScratch[:0])
		if usedCleanup {
			needsCleanupHelper = true
		}
	}
	if needsCleanupHelper {
		c.declareCleanupHelper(cleanupHelperName)
	}
}

func (c *rescopeRewriteCallback) getCleanupHelperPropertyName() string {
	name := "_"
	for {
		if _, exists := c.chunksUsingCrossChunkSymbols[name]; !exists {
			return name
		}
		name += rescopeDisambiguationSuffix
	}
}

func (c *rescopeRewriteCallback) declareCleanupHelper(cleanupHelperName string) {
	script := c.compiler.GetNodeForCodeInsertion(nil)
	globalSymbolNamespace := c.symbolInfo.globalSymbolNamespace
	paramName := "a"
	if globalSymbolNamespace == "a" {
		paramName = "b"
	}
	loopVar := "c"
	if paramName != "b" && globalSymbolNamespace != "b" {
		loopVar = "b"
	}
	splitCall := rhino.IR_Call(rhino.IR_Getprop(rhino.IR_Name(paramName), "split"), rhino.IR_String(","))
	deleteStmt := rhino.IR_ExprResult(
		rhino.IR_Delprop(rhino.IR_Getelem(rhino.IR_Name(globalSymbolNamespace), rhino.IR_Name(loopVar))))
	forOf := rhino.NewNodeWith3Children(
		rhino.Token_FOR_OF,
		rhino.IR_Let(rhino.IR_Name(loopVar)),
		splitCall,
		rhino.IR_BlockOfNode(deleteStmt))
	helperFn := rhino.IR_Function(
		rhino.IR_Name(""),
		rhino.IR_ParamList(rhino.IR_Name(paramName)),
		rhino.IR_BlockOfNode(forOf))
	helperDecl := rhino.IR_ExprResult(
		rhino.IR_Assign(rhino.IR_Getprop(rhino.IR_Name(globalSymbolNamespace), cleanupHelperName), helperFn))
	rescopeAddDeclaration(script, script.GetFirstChild(), helperDecl)
	NodeUtil_MarkNewScopesChanged(helperFn, c.compiler)
	NodeUtil_AddFeatureToScript(script, parser.FeatureSet_Feature_FOR_OF, c.compiler)
	NodeUtil_AddFeatureToScript(script, parser.FeatureSet_Feature_LET_DECLARATIONS, c.compiler)
	c.compiler.ReportChangeToEnclosingScope(script)
}

func (c *rescopeRewriteCallback) declareLocalAliasesAndWrappersForChunk(
	chunk *JSChunk,
	wrapperAssignmentChunks map[string]*JSChunk,
	cleanupHelperName string,
	allLocalAliases []string,
	aliasesToKeep []string,
	aliasesToDelete []string,
) ([]string, []string, []string, bool) {
	inputs := chunk.GetInputs()
	if len(inputs) == 0 {
		return allLocalAliases, aliasesToKeep, aliasesToDelete, false
	}
	script := inputs[0].GetAstRoot(c.compiler)
	insertionPoint := script.GetFirstChild()
	changed := false
	usedCleanupHelper := false
	globalSymbolNamespace := c.symbolInfo.globalSymbolNamespace

	for name := range c.localAliasesForUnwrappedCrossChunkNames[chunk].All() {
		allLocalAliases = append(allLocalAliases, name)
	}

	if c.optimizeLocalAccess.ShouldWrapReassignableSymbols() {
		for name := range c.symbolInfo.wrappedCrossChunkNames.All() {
			assignmentChunk := wrapperAssignmentChunks[name]
			usingChunks := c.chunksUsingWrappedReassignableSymbols[name]
			usedHere := usingChunks.Has(chunk)
			if chunk == assignmentChunk {
				if usedHere {
					rescopeAddDeclaration(
						script,
						insertionPoint,
						rhino.IR_VarWithValue(
							rhino.IR_Name(name),
							rhino.IR_Assign(rhino.IR_Getprop(rhino.IR_Name(globalSymbolNamespace), name), rhino.IR_Objectlit())))
				} else {
					rescopeAddDeclaration(
						script,
						insertionPoint,
						rhino.IR_ExprResult(
							rhino.IR_Assign(rhino.IR_Getprop(rhino.IR_Name(globalSymbolNamespace), name), rhino.IR_Objectlit())))
				}
				changed = true
			} else if usedHere {
				allLocalAliases = append(allLocalAliases, name)
			}
		}
	}

	if len(allLocalAliases) != 0 {
		for _, name := range allLocalAliases {
			if c.optimizeLocalAccess.ShouldDeleteGlobalsOnLastAccess() &&
				c.canDeleteFromGlobalNamespace(chunk, name) {
				aliasesToDelete = append(aliasesToDelete, name)
			} else {
				aliasesToKeep = append(aliasesToKeep, name)
			}
		}
		outputFeatures := c.compiler.GetOptions().GetOutputFeatureSet()
		if outputFeatures.ContainsFeature(parser.FeatureSet_Feature_OBJECT_DESTRUCTURING) &&
			outputFeatures.ContainsFeature(parser.FeatureSet_Feature_FOR_OF) &&
			outputFeatures.ContainsFeature(parser.FeatureSet_Feature_LET_DECLARATIONS) {
			objectPattern := rhino.IR_ObjectPattern()
			for _, name := range aliasesToKeep {
				stringKey := rhino.IR_StringKeyWithValue(name, rhino.IR_Name(name))
				stringKey.SetShorthandProperty(true)
				objectPattern.AddChildToBack(stringKey)
			}
			for _, name := range aliasesToDelete {
				stringKey := rhino.IR_StringKeyWithValue(name, rhino.IR_Name(name))
				stringKey.SetShorthandProperty(true)
				objectPattern.AddChildToBack(stringKey)
			}
			rescopeAddDeclaration(
				script, insertionPoint, rhino.IR_VarWithValue(objectPattern, rhino.IR_Name(globalSymbolNamespace)))
			NodeUtil_AddFeatureToScript(script, parser.FeatureSet_Feature_OBJECT_DESTRUCTURING, c.compiler)
			if len(aliasesToDelete) > 0 {
				cleanupCall := rhino.IR_ExprResult(
					rhino.IR_Call(
						rhino.IR_Getprop(rhino.IR_Name(globalSymbolNamespace), cleanupHelperName),
						rhino.IR_String(strings.Join(aliasesToDelete, ","))))
				rescopeAddDeclaration(script, insertionPoint, cleanupCall)
				usedCleanupHelper = true
			}
		} else {
			for _, name := range aliasesToKeep {
				rescopeAddDeclaration(
					script,
					insertionPoint,
					rhino.IR_VarWithValue(rhino.IR_Name(name), rhino.IR_Getprop(rhino.IR_Name(globalSymbolNamespace), name)))
			}
			for _, name := range aliasesToDelete {
				rescopeAddDeclaration(
					script,
					insertionPoint,
					rhino.IR_VarWithValue(rhino.IR_Name(name), rhino.IR_Getprop(rhino.IR_Name(globalSymbolNamespace), name)))
				rescopeAddDeclaration(
					script,
					insertionPoint,
					rhino.IR_ExprResult(rhino.IR_Delprop(rhino.IR_Getprop(rhino.IR_Name(globalSymbolNamespace), name))))
			}
		}
		changed = true
	}

	if changed {
		c.compiler.ReportChangeToEnclosingScope(script)
	}
	return allLocalAliases, aliasesToKeep, aliasesToDelete, usedCleanupHelper
}

func (c *rescopeRewriteCallback) canDeleteFromGlobalNamespace(currentChunk *JSChunk, name string) bool {
	if c.compiler.GetCodingConvention().IsExportedWithLocal(name /* local= */, false) {
		return false
	}
	if c.renameHelper.hasUnwrappedLocalNameInDefiningChunk(name) &&
		c.symbolInfo.globalNamesWithInnerScopeWriteInDefiningChunk.Has(name) {
		return false
	}
	usingChunks := c.chunksUsingCrossChunkSymbols[name]
	if !usingChunks.Has(currentChunk) {
		return false
	}
	chunkGraph := c.compiler.GetChunkGraph()
	for otherChunk := range usingChunks.All() {
		if otherChunk != currentChunk && !chunkGraph.DependsOn(currentChunk, otherChunk) {
			return false
		}
	}
	return true
}

func rescopeAddDeclaration(script rhino.Node, insertionPoint rhino.Node, decl rhino.Node) {
	decl.SrcrefTree(script)
	if insertionPoint.IsNil() {
		script.AddChildToBack(decl)
	} else {
		decl.InsertBefore(insertionPoint)
	}
}
