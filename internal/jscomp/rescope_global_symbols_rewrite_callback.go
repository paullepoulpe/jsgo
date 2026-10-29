// Port of com/google/javascript/jscomp/RescopeGlobalSymbolsRewriteCallback.java @ closure-compiler v20261021

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
)

const rescopeDisambiguationSuffix = "$"

type rescopeOrderedChunkSet struct {
	chunks []*JSChunk
	seen   map[*JSChunk]struct{}
}

func (s *rescopeOrderedChunkSet) add(chunk *JSChunk) {
	if s.seen == nil {
		s.seen = make(map[*JSChunk]struct{})
	}
	if _, ok := s.seen[chunk]; !ok {
		s.seen[chunk] = struct{}{}
		s.chunks = append(s.chunks, chunk)
	}
}

func (s *rescopeOrderedChunkSet) contains(chunk *JSChunk) bool {
	if s == nil {
		return false
	}
	_, ok := s.seen[chunk]
	return ok
}

type rescopeChunkGlobal struct {
	root rhino.Node
	name rhino.Node
}

// rescopeRewriteCallback is the node traversal callback to rewrite global
// symbols as part of the RescopeGlobalSymbols pass.
type rescopeRewriteCallback struct {
	compiler              AbstractCompiler
	globalSymbolNamespace string
	disambigPrefix        string
	assumeCrossChunkNames bool
	optimizeLocalAccess   CompilerOptions_OptimizeLocalAccess
	externNames           map[uint32]struct{}
	symbolInfo            *rescopeSymbolInfo

	preDeclarations []rescopeChunkGlobal
	lhsScratch      []rhino.Node

	localAliasesForUnwrappedCrossChunkNames map[*JSChunk]*rescopeOrderedStringSet
	chunksUsingWrappedReassignableSymbols   map[string]*rescopeOrderedChunkSet
	wrappedReassignableCrossChunkNames      rescopeOrderedStringSet
}

var (
	_ nodeTraversal_visitTokens    = (*rescopeRewriteCallback)(nil)
	_ nodeTraversal_onlyNameLeaves = (*rescopeRewriteCallback)(nil)
)

func (*rescopeRewriteCallback) visitTokens() *ast.TokenSet { return &rescopeNameVisitTokens }
func (*rescopeRewriteCallback) onlyNameLeaves()            {}

func rescopeCollectExternVariableNameIDs(compiler AbstractCompiler, externs rhino.Node) map[uint32]struct{} {
	externsRefs := NewReferenceCollector(
		compiler, ReferenceCollector_DO_NOTHING_BEHAVIOR, NewSyntacticScopeCreator(compiler))
	externsRefs.ProcessRoot(externs)
	externsNames := make(map[uint32]struct{})
	for v := range externsRefs.GetAllSymbols() {
		if !v.IsParam() {
			externsNames[v.nameID] = struct{}{}
		}
	}
	return externsNames
}

func newRescopeRewriteCallback(
	compiler AbstractCompiler,
	globalSymbolNamespace string,
	assumeCrossChunkNames bool,
	optimizeLocalAccess CompilerOptions_OptimizeLocalAccess,
	externs rhino.Node,
	symbolInfo *rescopeSymbolInfo) *rescopeRewriteCallback {
	c := &rescopeRewriteCallback{
		compiler:              compiler,
		globalSymbolNamespace: globalSymbolNamespace,
		disambigPrefix:        globalSymbolNamespace + rescopeDisambiguationSuffix,
		assumeCrossChunkNames: assumeCrossChunkNames,
		optimizeLocalAccess:   optimizeLocalAccess,
		externNames:           rescopeCollectExternVariableNameIDs(compiler, externs),
		symbolInfo:            symbolInfo,
	}
	c.collectWrappedReassignableCrossChunkNames()
	return c
}

func (c *rescopeRewriteCallback) isCrossChunkName(name string) bool {
	if c.assumeCrossChunkNames {
		return true
	}
	if _, ok := c.symbolInfo.crossChunkNames[name]; ok {
		return true
	}
	return c.compiler.GetCodingConvention().IsExportedWithLocal(name /* local= */, false)
}

func (c *rescopeRewriteCallback) isExternVar(id uint32, scope *Scope) bool {
	if id <= rhino.RhinoStringPool_EmptyID {
		return false
	}
	return c.isExternVarResolved(id, scope.GetVarByID(id))
}

func (c *rescopeRewriteCallback) isExternVarResolved(id uint32, v *Var) bool {
	if id <= rhino.RhinoStringPool_EmptyID {
		return false
	}
	if v == nil || v.IsExtern() {
		return true
	}
	if v.IsGlobal() {
		_, ok := c.externNames[id]
		return ok
	}
	return false
}

func (c *rescopeRewriteCallback) collectWrappedReassignableCrossChunkNames() {
	if !c.optimizeLocalAccess.ShouldWrapReassignableSymbols() {
		return
	}
	for _, name := range c.symbolInfo.crossChunkNamesWithWriteFromOtherChunk.names {
		c.wrappedReassignableCrossChunkNames.add(name)
	}
	for _, name := range c.symbolInfo.globalNamesWithInnerScopeWriteInDefiningChunk.names {
		if c.isCrossChunkName(name) {
			c.wrappedReassignableCrossChunkNames.add(name)
		}
	}
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
		if (isGlobalDeclaration && c.isCrossChunkName(name)) || c.isExternVar(lhs.GetStringID(), scope) {
			hasImportantName = true
			break
		}
	}

	if hasImportantName {
		c.rewriteNameDeclaration(t, scope, declaration, c.lhsScratch, isGlobalDeclaration)
	}
}

func (c *rescopeRewriteCallback) rewriteNameDeclaration(
	t *NodeTraversal, scope *Scope, declaration rhino.Node, allLhsNodes []rhino.Node, isGlobalDeclaration bool) {
	input := t.GetInput()

	if c.optimizeLocalAccess.ShouldOptimize() {
		if c.containsRescopedOrInitializedVars(declaration, isGlobalDeclaration) {
			var astRoot rhino.Node
			for _, lhs := range allLhsNodes {
				name := lhs.GetString()
				_, readInDefining := c.symbolInfo.globalNamesWithReadInDefiningChunk[name]
				if !c.symbolInfo.crossChunkNamesWithWriteFromOtherChunk.contains(name) &&
					!c.wrappedReassignableCrossChunkNames.contains(name) &&
					!c.isExternVar(lhs.GetStringID(), scope) &&
					(!c.isCrossChunkName(name) || readInDefining) {
					if astRoot.IsNil() {
						astRoot = input.GetAstRoot(c.compiler)
					}
					c.preDeclarations = append(c.preDeclarations,
						rescopeChunkGlobal{root: astRoot, name: rhino.IR_Name(name).Srcref(lhs)})
				}
			}
		}
	} else {
		var astRoot rhino.Node
		for _, lhs := range allLhsNodes {
			name := lhs.GetString()
			if !(isGlobalDeclaration && c.isCrossChunkName(name)) && !c.isExternVar(lhs.GetStringID(), scope) {
				if astRoot.IsNil() {
					astRoot = input.GetAstRoot(c.compiler)
				}
				c.preDeclarations = append(c.preDeclarations,
					rescopeChunkGlobal{root: astRoot, name: rhino.IR_Name(name).Srcref(lhs)})
			}
		}
	}

	for decl := range ast.ReadAhead(ast.MustNameDeclaration(declaration).FirstDeclarator()) {
		if name, ok := decl.AsName(); ok {
			if init, hasInit := name.Init(); hasInit {
				assign := rhino.IR_Assign(name.Node().CloneNode(), init.Detach().Node())
				name.Node().ReplaceWith(assign)
				assign.SetJSDocInfo(declaration.GetJSDocInfo())
			}
		} else if dlhs, ok := decl.AsDestructuringLhs(); ok {
			pattern := dlhs.Pattern().Node()
			init, hasInit := dlhs.Init()
			if !hasInit {
				if !NodeUtil_IsEnhancedFor(declaration.GetParent()) {
					jutil.CheckStateMsg1(
						false,
						"DESTRUCTURING_LHS should have two children: %s",
						declaration.ToStringTree())
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

func (c *rescopeRewriteCallback) containsRescopedOrInitializedVars(declaration rhino.Node, isGlobalDeclaration bool) bool {
	for decl := range ast.MustNameDeclaration(declaration).Declarators() {
		if name, ok := decl.AsName(); ok {
			ident := name.Ident()
			if isGlobalDeclaration &&
				(c.symbolInfo.crossChunkNamesWithWriteFromOtherChunk.contains(ident) ||
					c.wrappedReassignableCrossChunkNames.contains(ident)) {
				return true
			}
			if _, hasInit := name.Init(); hasInit {
				return true
			}
		} else if _, ok := decl.AsDestructuringLhs(); ok {
			return true
		}
	}
	return false
}

func (c *rescopeRewriteCallback) ensureLocalAliasesForChunk(chunk *JSChunk) *rescopeOrderedStringSet {
	if s, ok := c.localAliasesForUnwrappedCrossChunkNames[chunk]; ok {
		return s
	}
	if c.localAliasesForUnwrappedCrossChunkNames == nil {
		c.localAliasesForUnwrappedCrossChunkNames = make(map[*JSChunk]*rescopeOrderedStringSet)
	}
	s := &rescopeOrderedStringSet{}
	c.localAliasesForUnwrappedCrossChunkNames[chunk] = s
	return s
}

func (c *rescopeRewriteCallback) ensureChunksUsingWrappedSymbol(name string) *rescopeOrderedChunkSet {
	if s, ok := c.chunksUsingWrappedReassignableSymbols[name]; ok {
		return s
	}
	if c.chunksUsingWrappedReassignableSymbols == nil {
		c.chunksUsingWrappedReassignableSymbols = make(map[string]*rescopeOrderedChunkSet)
	}
	s := &rescopeOrderedChunkSet{}
	c.chunksUsingWrappedReassignableSymbols[name] = s
	return s
}

func (c *rescopeRewriteCallback) visitName(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	id := n.GetStringID()
	if parent.IsFunction() && id <= rhino.RhinoStringPool_EmptyID {
		return
	}
	var var_ *Var
	if id > rhino.RhinoStringPool_EmptyID {
		var_ = t.GetScope().GetVarByID(id)
		if c.isExternVarResolved(id, var_) {
			return
		}
	}

	name := n.GetString()
	if !var_.IsGlobal() &&
		(name == c.globalSymbolNamespace || strings.HasPrefix(name, c.disambigPrefix)) {
		n.SetString(name + rescopeDisambiguationSuffix)
		c.compiler.ReportChangeToEnclosingScope(n)
	}

	if !(var_.IsGlobal() && c.isCrossChunkName(name)) {
		return
	}

	currentChunk := t.GetChunk()

	if c.optimizeLocalAccess.ShouldOptimize() {
		definingChunk := var_.GetInput().GetChunk()
		if currentChunk == definingChunk {
			_, readInDefining := c.symbolInfo.globalNamesWithReadInDefiningChunk[name]
			if !c.symbolInfo.crossChunkNamesWithWriteFromOtherChunk.contains(name) &&
				!c.wrappedReassignableCrossChunkNames.contains(name) &&
				readInDefining {
				if NodeUtil_IsLValue(n) && (!NodeUtil_IsNameDeclaration(parent) || n.HasChildren()) {
					c.addGlobalNamespaceAlias(n, name)
				}
				return
			}
		} else if c.optimizeLocalAccess.ShouldOptimizeAllChunks() {
			if !c.symbolInfo.crossChunkNamesWithWriteFromOtherChunk.contains(name) &&
				!c.symbolInfo.globalNamesWithInnerScopeWriteInDefiningChunk.contains(name) &&
				c.compiler.GetChunkGraph().DependsOn(currentChunk, definingChunk) {
				c.ensureLocalAliasesForChunk(currentChunk).add(name)
				return
			}
		}
	}

	if c.wrappedReassignableCrossChunkNames.contains(name) {
		c.ensureChunksUsingWrappedSymbol(name).add(currentChunk)
		replacement := rhino.IR_Getprop(rhino.IR_Name(name), "_")
		c.replaceSymbol(n, name, replacement)
	} else {
		replacement := rhino.IR_Getprop(rhino.IR_Name(c.globalSymbolNamespace), name)
		c.replaceSymbol(n, name, replacement)
	}
}

func (c *rescopeRewriteCallback) replaceSymbol(node rhino.Node, name string, replacement rhino.Node) {
	parent := node.GetParent()
	replacement.SrcrefTree(node)
	node.ReplaceWith(replacement)
	c.compiler.ReportChangeToEnclosingScope(replacement)
	if parent.IsCall() {
		if _, maybeThis := c.symbolInfo.maybeReferencesThis[name]; !maybeThis {
			parent.PutBooleanProp(rhino.Node_FREE_CALL, false)
		}
	}
	c.compiler.ReportChangeToEnclosingScope(parent)
}

func (c *rescopeRewriteCallback) addGlobalNamespaceAlias(n rhino.Node, name string) {
	stmt := NodeUtil_GetEnclosingStatement(n)
	if stmt.IsNil() {
		return
	}
	parent := n.GetParent()
	if assign, ok := ast.AsAssign(parent); ok && assign.Target().Is(n) {
		rhs := assign.Value().Detach().Node()
		globalNamespaceAssign := rhino.IR_Assign(rhino.IR_Getprop(rhino.IR_Name(c.globalSymbolNamespace), name), rhs)
		globalNamespaceAssign.SrcrefTree(n)
		parent.AddChildToBack(globalNamespaceAssign)
		c.compiler.ReportChangeToEnclosingScope(parent)
	} else {
		globalNamespaceAssign := rhino.IR_Assign(
			rhino.IR_Getprop(rhino.IR_Name(c.globalSymbolNamespace), name), rhino.IR_Name(name))
		aliasStatement := rhino.IR_ExprResult(globalNamespaceAssign)
		aliasStatement.SrcrefTree(n)
		aliasStatement.InsertAfter(stmt)
		c.compiler.ReportChangeToEnclosingScope(stmt.GetParent())
	}
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
		if global.root.HasChildren() && global.root.GetFirstChild().IsVar() {
			global.root.GetFirstChild().AddChildToBack(global.name)
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
	if len(c.wrappedReassignableCrossChunkNames.names) > 0 {
		wrapperAssignmentChunks = make(map[string]*JSChunk, len(c.wrappedReassignableCrossChunkNames.names))
		for _, name := range c.wrappedReassignableCrossChunkNames.names {
			usingChunks := c.chunksUsingWrappedReassignableSymbols[name]
			if usingChunks == nil || len(usingChunks.chunks) == 0 {
				wrapperAssignmentChunks[name] = defaultRootChunk
			} else {
				wrapperAssignmentChunks[name] = chunkGraph.GetDeepestCommonDependencyInclusiveOfCollection(usingChunks.chunks)
			}
		}
	}

	var localAliasesScratch []string
	for _, chunk := range chunks {
		localAliasesScratch = c.declareLocalAliasesAndWrappersForChunk(chunk, wrapperAssignmentChunks, localAliasesScratch[:0])
	}
}

func (c *rescopeRewriteCallback) declareLocalAliasesAndWrappersForChunk(
	chunk *JSChunk,
	wrapperAssignmentChunks map[string]*JSChunk,
	localAliases []string) []string {
	inputs := chunk.GetInputs()
	if len(inputs) == 0 {
		return localAliases
	}
	script := inputs[0].GetAstRoot(c.compiler)
	insertionPoint := script.GetFirstChild()
	changed := false

	if unwrapped := c.localAliasesForUnwrappedCrossChunkNames[chunk]; unwrapped != nil {
		localAliases = append(localAliases, unwrapped.names...)
	}

	if c.optimizeLocalAccess.ShouldWrapReassignableSymbols() {
		for _, name := range c.wrappedReassignableCrossChunkNames.names {
			assignmentChunk := wrapperAssignmentChunks[name]
			usingChunks := c.chunksUsingWrappedReassignableSymbols[name]
			usedHere := usingChunks.contains(chunk)
			if chunk == assignmentChunk {
				if usedHere {
					rescopeAddDeclaration(
						script,
						insertionPoint,
						rhino.IR_VarWithValue(
							rhino.IR_Name(name),
							rhino.IR_Assign(rhino.IR_Getprop(rhino.IR_Name(c.globalSymbolNamespace), name), rhino.IR_Objectlit())))
				} else {
					rescopeAddDeclaration(
						script,
						insertionPoint,
						rhino.IR_ExprResult(
							rhino.IR_Assign(rhino.IR_Getprop(rhino.IR_Name(c.globalSymbolNamespace), name), rhino.IR_Objectlit())))
				}
				changed = true
			} else if usedHere {
				localAliases = append(localAliases, name)
			}
		}
	}

	if len(localAliases) != 0 {
		if c.compiler.GetOptions().GetOutputFeatureSet().ContainsFeature(parser.FeatureSet_Feature_OBJECT_DESTRUCTURING) {
			objectPattern := rhino.IR_ObjectPattern()
			for _, name := range localAliases {
				stringKey := rhino.IR_StringKeyWithValue(name, rhino.IR_Name(name))
				stringKey.SetShorthandProperty(true)
				objectPattern.AddChildToBack(stringKey)
			}
			rescopeAddDeclaration(
				script, insertionPoint, rhino.IR_VarWithValue(objectPattern, rhino.IR_Name(c.globalSymbolNamespace)))
			NodeUtil_AddFeatureToScript(script, parser.FeatureSet_Feature_OBJECT_DESTRUCTURING, c.compiler)
		} else {
			for _, name := range localAliases {
				rescopeAddDeclaration(
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
	return localAliases
}

func rescopeAddDeclaration(script rhino.Node, insertionPoint rhino.Node, decl rhino.Node) {
	decl.SrcrefTree(script)
	if insertionPoint.IsNil() {
		script.AddChildToBack(decl)
	} else {
		decl.InsertBefore(insertionPoint)
	}
}
