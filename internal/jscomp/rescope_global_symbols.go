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
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// RescopeGlobalSymbols finds all references to global symbols and rewrites
// them to be property accesses to a special object with the same name as the
// global symbol.
//
// Given the name of the global object is NS:
//
//	var a = 1; function b() { return a }
//	-> NS.a = 1; NS.b = function b() { return NS.a }
type RescopeGlobalSymbols struct {
	compiler              AbstractCompiler
	globalSymbolNamespace string
	addExtern             bool
	assumeCrossChunkNames bool
	optimizeLocalAccess   CompilerOptions_OptimizeLocalAccess
}

// NewRescopeGlobalSymbols creates a new RescopeGlobalSymbols compiler pass.
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

// NewRescopeGlobalSymbolsWithAddExtern creates a RescopeGlobalSymbols pass
// with explicit control over extern synthesis (used in testing).
//
// Java: RescopeGlobalSymbols(AbstractCompiler, String, boolean, boolean, CompilerOptions.OptimizeLocalAccess)
func NewRescopeGlobalSymbolsWithAddExtern(
	compiler AbstractCompiler,
	globalSymbolNamespace string,
	addExtern bool,
	assumeCrossChunkNames bool,
	optimizeLocalAccess CompilerOptions_OptimizeLocalAccess) *RescopeGlobalSymbols {
	return &RescopeGlobalSymbols{
		compiler:              compiler,
		globalSymbolNamespace: globalSymbolNamespace,
		addExtern:             addExtern,
		assumeCrossChunkNames: assumeCrossChunkNames,
		optimizeLocalAccess:   optimizeLocalAccess,
	}
}

func (r *RescopeGlobalSymbols) addExternForGlobalSymbolNamespace() {
	varNode := rhino.IR_Var(rhino.IR_Name(r.globalSymbolNamespace))
	input := r.compiler.GetSynthesizedExternsInput()
	input.GetAstRoot(r.compiler).AddChildToBack(varNode)
	r.compiler.ReportChangeToEnclosingScope(varNode)
}

// Java: RescopeGlobalSymbols.process(Node, Node)
func (r *RescopeGlobalSymbols) Process(externs rhino.Node, root rhino.Node) {
	if r.addExtern {
		r.addExternForGlobalSymbolNamespace()
	}

	// Turn global named function and class statements into var assignments.
	declsCallback := rescopeGlobalDeclsCallback{compiler: r.compiler}
	NodeTraversal_Traverse(r.compiler, root, &declsCallback)

	// Find global names that are used in more than one chunk and names that
	// may refer to functions that reference `this`.
	analyzer := rescopeAnalyzeNamesCallback{
		trackLocalAccessSets: r.optimizeLocalAccess.ShouldOptimize(),
		info: rescopeSymbolInfo{
			crossChunkNames:     make(map[string]struct{}),
			maybeReferencesThis: make(map[string]struct{}),
		},
	}
	if analyzer.trackLocalAccessSets {
		analyzer.info.globalNamesWithReadInDefiningChunk = make(map[string]struct{})
	}
	NodeTraversal_Traverse(r.compiler, root, &analyzer)

	// Rewrite all references to be property accesses of the single symbol.
	rewriteScope := newRescopeRewriteCallback(
		r.compiler,
		r.globalSymbolNamespace,
		r.assumeCrossChunkNames,
		r.optimizeLocalAccess,
		externs,
		&analyzer.info)
	NodeTraversal_Traverse(r.compiler, root, rewriteScope)

	// Remove the var/let/const from statements in global scope if the declared
	// names have been rewritten in the previous pass.
	removeVarCallback := rescopeRemoveGlobalVarCallback{compiler: r.compiler}
	NodeTraversal_Traverse(r.compiler, root, &removeVarCallback)

	rewriteScope.addDeclarations()
}

// rescopeGlobalDeclsCallback rewrites global function and class declarations to
// var statements + assignment. Ignores non-global declarations.
type rescopeGlobalDeclsCallback struct {
	NodeTraversal_AbstractShallowStatementCallbackBase
	compiler AbstractCompiler
}

var rescopeGlobalDeclsVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_FUNCTION)
	s.Add(rhino.Token_CLASS)
	return s
}()

var (
	_ nodeTraversal_visitTokens   = (*rescopeGlobalDeclsCallback)(nil)
	_ nodeTraversal_noVisitLeaves = (*rescopeGlobalDeclsCallback)(nil)
)

func (*rescopeGlobalDeclsCallback) visitTokens() *ast.TokenSet { return &rescopeGlobalDeclsVisitTokens }
func (*rescopeGlobalDeclsCallback) noVisitLeaves()             {}

func (c *rescopeGlobalDeclsCallback) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if !t.InGlobalScope() {
		return
	}
	if !NodeUtil_IsFunctionDeclaration(n) && !NodeUtil_IsClassDeclaration(n) {
		return
	}
	var name string
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
	varNode := NodeUtil_NewVarNode(name, n)
	if prev.IsNil() {
		parent.AddChildToFront(varNode)
	} else {
		varNode.InsertAfter(prev)
	}
	c.compiler.ReportChangeToEnclosingScope(parent)
}

type rescopeOrderedStringSet struct {
	names []string
	seen  map[string]struct{}
}

func (s *rescopeOrderedStringSet) add(name string) {
	if s.seen == nil {
		s.seen = make(map[string]struct{})
	}
	if _, ok := s.seen[name]; !ok {
		s.seen[name] = struct{}{}
		s.names = append(s.names, name)
	}
}

func (s *rescopeOrderedStringSet) contains(name string) bool {
	_, ok := s.seen[name]
	return ok
}

type rescopeSymbolInfo struct {
	crossChunkNames                               map[string]struct{}
	crossChunkNamesWithWriteFromOtherChunk        rescopeOrderedStringSet
	globalNamesWithReadInDefiningChunk            map[string]struct{}
	globalNamesWithInnerScopeWriteInDefiningChunk rescopeOrderedStringSet
	maybeReferencesThis                           map[string]struct{}
}

// rescopeAnalyzeNamesCallback combines FindCrossChunkNamesCallback and
// FindNamesReferencingThis into a single post-order pass over NAME nodes.
type rescopeAnalyzeNamesCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	trackLocalAccessSets bool
	info                 rescopeSymbolInfo
}

var rescopeNameVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_NAME)
	return s
}()

var (
	_ nodeTraversal_visitTokens    = (*rescopeAnalyzeNamesCallback)(nil)
	_ nodeTraversal_onlyNameLeaves = (*rescopeAnalyzeNamesCallback)(nil)
)

func (*rescopeAnalyzeNamesCallback) visitTokens() *ast.TokenSet { return &rescopeNameVisitTokens }
func (*rescopeAnalyzeNamesCallback) onlyNameLeaves()            {}

func (c *rescopeAnalyzeNamesCallback) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	nameNode, ok := ast.AsName(n)
	if !ok {
		return
	}
	id := n.GetStringID()
	if id <= rhino.RhinoStringPool_EmptyID {
		return
	}
	v := t.GetScope().GetVarByID(id)
	if v == nil || !v.IsGlobal() {
		return
	}
	name := v.GetName()
	c.visitCrossChunkName(t, nameNode, n, name, v, parent)
	c.visitMaybeReferencesThis(nameNode, n, name, v, parent)
}

func (c *rescopeAnalyzeNamesCallback) visitCrossChunkName(
	t *NodeTraversal, nameNode ast.Name, n rhino.Node, name string, v *Var, parent rhino.Node) {
	if !c.trackLocalAccessSets {
		if _, exists := c.info.crossChunkNames[name]; exists {
			return
		}
	}
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
		c.info.globalNamesWithInnerScopeWriteInDefiningChunk.add(name)
	}

	if input == nil {
		c.info.crossChunkNames[name] = struct{}{}
		return
	}
	if input.GetChunk() != t.GetChunk() {
		c.info.crossChunkNames[name] = struct{}{}
		if c.trackLocalAccessSets &&
			NodeUtil_IsLValue(n) &&
			(!NodeUtil_IsNameDeclaration(parent) || hasInit) {
			c.info.crossChunkNamesWithWriteFromOtherChunk.add(name)
		}
	} else if c.trackLocalAccessSets && n != v.GetNameNode() {
		c.info.globalNamesWithReadInDefiningChunk[name] = struct{}{}
	}
}

func (c *rescopeAnalyzeNamesCallback) visitMaybeReferencesThis(
	nameNode ast.Name, n rhino.Node, name string, v *Var, parent rhino.Node) {
	if v == nil || !v.IsGlobal() {
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
		return
	}
	if _, already := c.info.maybeReferencesThis[name]; already {
		return
	}
	if value.IsNil() || !value.IsFunction() || NodeUtil_ReferencesOwnReceiver(value) {
		c.info.maybeReferencesThis[name] = struct{}{}
	}
}

// rescopeRemoveGlobalVarCallback removes every occurrence of var/let/const
// that declares a rewritten global variable.
type rescopeRemoveGlobalVarCallback struct {
	NodeTraversal_AbstractShallowStatementCallbackBase
	compiler AbstractCompiler
	commas   []rhino.Node
	children []rhino.Node
}

var rescopeRemoveGlobalVarVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_VAR)
	s.Add(rhino.Token_LET)
	s.Add(rhino.Token_CONST)
	return s
}()

var (
	_ nodeTraversal_visitTokens   = (*rescopeRemoveGlobalVarCallback)(nil)
	_ nodeTraversal_noVisitLeaves = (*rescopeRemoveGlobalVarCallback)(nil)
)

func (*rescopeRemoveGlobalVarCallback) visitTokens() *ast.TokenSet {
	return &rescopeRemoveGlobalVarVisitTokens
}
func (*rescopeRemoveGlobalVarCallback) noVisitLeaves() {}

func (c *rescopeRemoveGlobalVarCallback) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if !NodeUtil_IsNameDeclaration(n) {
		return
	}

	c.commas = c.commas[:0]
	c.children = c.children[:0]
	parentIsAnyFor := NodeUtil_IsAnyFor(parent)
	allNameOrDestructuring := true
	for ch := n.GetFirstChild(); !ch.IsNil(); ch = ch.GetNext() {
		if !ch.IsName() && !ch.IsDestructuringLhs() {
			allNameOrDestructuring = false
		}
		if ch.IsAssign() || parentIsAnyFor {
			c.children = append(c.children, ch)
		}
	}
	if allNameOrDestructuring {
		return
	}
	isForInit := parentIsAnyFor && parent.GetFirstChild() == n
	for _, ch := range c.children {
		if isForInit {
			c.commas = append(c.commas, ch.CloneTree())
		} else {
			expr := rhino.IR_ExprResult(ch.CloneTree()).Srcref(ch)
			NodeUtil_MarkNewScopesChanged(expr, c.compiler)
			expr.InsertBefore(n)
		}
	}
	if len(c.commas) != 0 {
		comma := joinNodesOnComma(c.commas, n)
		comma.InsertBefore(n)
	}
	n.Detach()
	NodeUtil_MarkFunctionsDeleted(n, c.compiler)
	c.compiler.ReportChangeToEnclosingScope(parent)
}

func joinNodesOnComma(commas []rhino.Node, source rhino.Node) rhino.Node {
	comma := commas[0]
	for i := 1; i < len(commas); i++ {
		nextComma := rhino.IR_Comma(comma, commas[i])
		nextComma.SrcrefIfMissing(source)
		comma = nextComma
	}
	return comma
}
