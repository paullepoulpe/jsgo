// Port of com/google/javascript/jscomp/RescopeGlobalSymbols.java @ closure-compiler v20261104

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
	"github.com/paullepoulpe/jsgo/internal/strpool"
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
		compiler:              r.compiler,
		assumeCrossChunkNames: r.assumeCrossChunkNames,
		trackLocalAccessSets:  r.optimizeLocalAccess.ShouldOptimize(),
		renameHelperBuilder:   newRescopeRenameHelperBuilder(r.compiler, r.globalSymbolNamespace, r.optimizeLocalAccess),
		info: rescopeSymbolInfo{
			globalSymbolNamespace: r.globalSymbolNamespace,
			externNameIDs:         collectExternVariableNameIDs(r.compiler, externs),
			maybeReferencesThis:   make(map[string]struct{}),
		},
	}
	if analyzer.trackLocalAccessSets {
		analyzer.info.globalNamesWithReadInDefiningChunk = make(map[string]struct{})
		analyzer.info.globalNamesWithUnsafeWriteForLocalAliasInDefiningChunk = make(map[string]struct{})
	}
	NodeTraversal_Traverse(r.compiler, root, &analyzer)
	analyzer.info.populateWrappedCrossChunkNames()
	renameHelper := analyzer.renameHelperBuilder.build(&analyzer.info)

	// Rewrite all references to be property accesses of the single symbol.
	rewriteScope := newRescopeRewriteCallback(
		r.compiler,
		r.optimizeLocalAccess,
		&analyzer.info,
		renameHelper)
	NodeTraversal_Traverse(r.compiler, root, rewriteScope)

	// Remove the var/let/const from statements in global scope if the declared
	// names have been rewritten in the previous pass.
	removeVarCallback := rescopeRemoveGlobalVarCallback{
		compiler:                     r.compiler,
		preserveLocalVarDeclarations: r.optimizeLocalAccess.ShouldPreserveLocalVarDeclarations(),
	}
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

// rescopeAnalyzeNamesCallback combines FindCrossChunkNamesCallback and
// FindNamesReferencingThis into a single post-order pass over NAME nodes.
type rescopeAnalyzeNamesCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler              AbstractCompiler
	assumeCrossChunkNames bool
	trackLocalAccessSets  bool
	renameHelperBuilder   rescopeRenameHelperBuilder
	info                  rescopeSymbolInfo
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
	if id <= strpool.EmptyID {
		return
	}
	name := n.GetString()
	var v *Var
	var resolved bool
	getVar := func() *Var {
		if !resolved {
			v = t.GetScope().GetVarByID(id)
			resolved = true
		}
		return v
	}
	c.visitCrossChunkName(t, n, name, getVar, parent)
	c.visitMaybeReferencesThis(nameNode, n, name, getVar, parent)
}

func (c *rescopeAnalyzeNamesCallback) visitCrossChunkName(
	t *NodeTraversal, n rhino.Node, name string, getVar func() *Var, parent rhino.Node) {
	if !c.trackLocalAccessSets && c.info.crossChunkNames.Has(name) {
		return
	}
	c.renameHelperBuilder.recordExistingName(name)
	v := getVar()
	if v == nil || !v.IsGlobal() {
		return
	}

	if c.assumeCrossChunkNames ||
		c.compiler.GetCodingConvention().IsExportedWithLocal(name /* local= */, false) {
		c.info.crossChunkNames.Add(name)
	}

	input := v.GetInput()

	if c.trackLocalAccessSets &&
		!t.InGlobalScope() &&
		(input == nil || input.GetChunk() == t.GetChunk()) &&
		NodeUtil_IsLValueWrite(n) {
		c.info.globalNamesWithInnerScopeWriteInDefiningChunk.Add(name)
	}

	if input == nil {
		c.info.crossChunkNames.Add(name)
		return
	}
	if input.GetChunk() != t.GetChunk() {
		c.info.crossChunkNames.Add(name)
		if c.trackLocalAccessSets {
			if NodeUtil_IsLValueWrite(n) {
				c.info.crossChunkNamesWithWriteFromOtherChunk.Add(name)
			}
			if !c.compiler.GetChunkGraph().DependsOn(t.GetChunk(), input.GetChunk()) {
				c.info.crossChunkNamesUsedInNonDependentChunk.Add(name)
			}
		}
	} else if c.trackLocalAccessSets {
		if rescopeIsReadReference(n) {
			c.info.globalNamesWithReadInDefiningChunk[name] = struct{}{}
		}
		if NodeUtil_IsLValueWrite(n) && !rescopeIsWriteSafeForLocalAliasMirroring(n, parent) {
			c.info.globalNamesWithUnsafeWriteForLocalAliasInDefiningChunk[name] = struct{}{}
		}
	}
}

func rescopeIsReadReference(n rhino.Node) bool {
	if !NodeUtil_IsLValue(n) {
		return true
	}
	parent := n.GetParent()
	return NodeUtil_IsCompoundAssignmentOp(parent) || parent.IsInc() || parent.IsDec()
}

func rescopeIsWriteSafeForLocalAliasMirroring(n rhino.Node, parent rhino.Node) bool {
	if (NodeUtil_IsAssignmentOp(parent) && parent.GetFirstChild() == n) ||
		parent.IsInc() ||
		parent.IsDec() {
		return true
	}
	if NodeUtil_IsNameDeclaration(parent) {
		return !NodeUtil_IsEnhancedFor(parent.GetParent())
	}
	if NodeUtil_IsLhsByDestructuring(n) {
		rootTarget := NodeUtil_GetRootTarget(n)
		rootTargetParent := rootTarget.GetParent()
		if rootTargetParent.IsDestructuringLhs() {
			declaration := rootTargetParent.GetParent()
			return NodeUtil_IsStatementBlock(declaration.GetParent())
		}
		if rootTargetParent.IsAssign() && rootTargetParent.GetFirstChild() == rootTarget {
			exprResult := rootTargetParent.GetParent()
			return exprResult.IsExprResult() && NodeUtil_IsStatementBlock(exprResult.GetParent())
		}
	}
	return false
}

func (c *rescopeAnalyzeNamesCallback) visitMaybeReferencesThis(
	nameNode ast.Name, n rhino.Node, name string, getVar func() *Var, parent rhino.Node) {
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
	v := getVar()
	if v == nil || !v.IsGlobal() {
		return
	}
	if value.IsNil() || !value.IsFunction() || NodeUtil_ReferencesOwnReceiver(value) {
		c.info.maybeReferencesThis[name] = struct{}{}
	}
}

// rescopeRemoveGlobalVarCallback removes var/let/const declarations whose
// global variables have been rewritten to property assignments by
// rescopeRewriteCallback, while preserving var declarations for variables that
// remain local to the chunk.
type rescopeRemoveGlobalVarCallback struct {
	NodeTraversal_AbstractShallowStatementCallbackBase
	compiler                     AbstractCompiler
	preserveLocalVarDeclarations bool
	commas                       []rhino.Node
}

var rescopeRemoveGlobalVarVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_VAR)
	s.Add(rhino.Token_LET)
	s.Add(rhino.Token_CONST)
	return s
}()

var _ nodeTraversal_visitTokens = (*rescopeRemoveGlobalVarCallback)(nil)

func (*rescopeRemoveGlobalVarCallback) visitTokens() *ast.TokenSet {
	return &rescopeRemoveGlobalVarVisitTokens
}

func (c *rescopeRemoveGlobalVarCallback) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if _, ok := ast.AsNameDeclaration(n); !ok {
		return
	}

	// If every child of a non-empty declaration is still a variable binding
	// (NAME or DESTRUCTURING_LHS) rather than a rewritten assignment (ASSIGN) or
	// property access (GETPROP), then no variables in the declaration were
	// rewritten into standalone expressions, and the declaration node must stay
	// in place as-is.
	if rescopeHasOnlyNameOrDestructuringChildren(n) {
		return
	}

	switch ast.RoleIn(n, parent) {
	case ast.RoleForInit, ast.RoleForInTarget, ast.RoleForOfTarget, ast.RoleForAwaitOfTarget:
		c.commas = c.commas[:0]
		for ch := n.GetFirstChild(); !ch.IsNil(); ch = ch.GetNext() {
			c.commas = append(c.commas, ch.CloneTree())
		}
		if len(c.commas) != 0 {
			comma := joinNodesOnComma(c.commas, n)
			comma.InsertBefore(n)
		}
	default:
		var currentVar rhino.Node
		for ch := n.GetFirstChild(); !ch.IsNil(); ch = ch.GetNext() {
			jutil.CheckStateMsg1(!ch.IsDestructuringLhs(), "Unexpected DESTRUCTURING_LHS child: %s", ch)
			if c.preserveLocalVarDeclarations && ch.IsName() {
				if currentVar.IsNil() {
					currentVar = rhino.IR_Var(ch.CloneTree()).Srcref(n)
					currentVar.SetJSDocInfo(n.GetJSDocInfo())
					currentVar.InsertBefore(n)
				} else {
					currentVar.AddChildToBack(ch.CloneTree())
				}
			} else if ch.IsAssign() {
				if !currentVar.IsNil() {
					NodeUtil_MarkNewScopesChanged(currentVar, c.compiler)
					currentVar = rhino.NilNode
				}
				expr := rhino.IR_ExprResult(ch.CloneTree()).Srcref(ch)
				NodeUtil_MarkNewScopesChanged(expr, c.compiler)
				expr.InsertBefore(n)
			}
		}
		if !currentVar.IsNil() {
			NodeUtil_MarkNewScopesChanged(currentVar, c.compiler)
		}
	}
	n.Detach()
	NodeUtil_MarkFunctionsDeleted(n, c.compiler)
	c.compiler.ReportChangeToEnclosingScope(parent)
}

func rescopeHasOnlyNameOrDestructuringChildren(n rhino.Node) bool {
	for ch := n.GetFirstChild(); !ch.IsNil(); ch = ch.GetNext() {
		if _, ok := ast.AsDeclarator(ch); !ok {
			return false
		}
	}
	return n.HasChildren()
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
