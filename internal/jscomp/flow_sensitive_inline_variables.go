// Port of com/google/javascript/jscomp/FlowSensitiveInlineVariables.java @ closure-compiler v20261021

/*
 * Copyright 2009 The Closure Compiler Authors.
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
	"sync"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/graph"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: the shapes this pass reads (declarations and assignments with
// their values, statement bodies, class members, function parameters,
// whether a name is written or read) go through the typed view
// (internal/ast); the subtree walks that look at every child stay on
// rhino.Node. The Go-only prefilters (scanStatementsForCandidateDefs,
// collectCandidateFuncs, fsivWalkGatherCandidates) follow the statement
// structure that ControlFlowAnalysis builds a graph for.

type fsivCandidateDef struct {
	varName string
	rhs     ast.Expr
	hasRead bool
	hasFunc bool
}

type fsivCandFunc struct {
	fn         ast.Function
	body       ast.Block
	depth      int
	hasFuncRhs bool
}

// FlowSensitiveInlineVariables inlines variables when possible. Using the
// information from MaybeReachingVariableUse and MustBeReachingVariableDef,
// this pass attempts to inline a variable by placing the value at the
// definition where the variable is used. The basic requirements for
// inlining are the following:
//
//   - There is exactly one reaching definition at the use of that variable
//   - There is exactly one use for that definition of the variable
//
// Other requirements can be found in Candidate#canInline. Currently this
// pass does not operate on the global scope due to compilation time.
type FlowSensitiveInlineVariables struct {
	// Implementation:
	//
	// This pass first performs a traversal to gather a list of Candidates that
	// could be inlined using GatherCandidates.
	//
	// The second step involves verifying that each candidate is actually safe
	// to inline with Candidate#canInline(Scope) and finally performing
	// inlining using Candidate#inlineVariable().
	//
	// The reason for the delayed evaluation of the candidates is because we
	// need two separate dataflow results.
	compiler         AbstractCompiler
	astAnalyzer      *AstAnalyzer
	codingConvention CodingConvention

	sideEffectPredicate   FlowSensitiveInlineVariables_SideEffectPredicate
	sideEffectPredicateFn func(rhino.Node) bool

	cfg          *ControlFlowGraph[rhino.Node]
	candidates   []FlowSensitiveInlineVariables_Candidate
	reachingDef  *MustBeReachingVariableDef
	reachingUses *MaybeReachingVariableUse
	escaped      *jutil.ObjectLinkedHashSet[*Var]
	allVarsInFn  *jutil.LinkedHashMap[string, *Var]

	reusableCfg         ReusableCfg
	reachingDefStorage  MustBeReachingVariableDef
	reachingUsesStorage MaybeReachingVariableUse
	allVarsMap          jutil.LinkedHashMap[string, *Var]
	candidateVarsMap    jutil.LinkedHashMap[string, *Var]
	namesToCheckSet     jutil.LinkedHashSet[string]
	escapedSet          jutil.ObjectLinkedHashSet[*Var]
	gatherCandidates    FlowSensitiveInlineVariables_GatherCandidates
	gatherCb            FlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback
	pathChecker         graph.CheckPathsBetweenNodes[rhino.Node, ControlFlowGraph_Branch]

	varsFinder          lvaAllVarsFinder
	inlinableDefs       []fsivCandidateDef
	unreadDefCount      int
	candFuncs           []fsivCandFunc
	maxCandDepth        int
	dummyGlobalScope    *Scope
	globalRoot          rhino.Node
	globalScopeOnce     sync.Once
	realGlobalScope     *Scope
	globalScopeCreator  *SyntacticScopeCreator
	sharedParent        *FlowSensitiveInlineVariables
	scopeCreator        *SyntacticScopeCreator
	curFn               rhino.Node
	curFnEnclosingScope *Scope
	curTraversal        *NodeTraversal
	parallelMode        bool
	parallelChanges     []rhino.Node
}

func (*FlowSensitiveInlineVariables) nodeTraversalNoVisit() {}

func (f *FlowSensitiveInlineVariables) getAstAnalyzer() *AstAnalyzer {
	if f.astAnalyzer == nil {
		f.astAnalyzer = f.compiler.GetAstAnalyzer()
	}
	return f.astAnalyzer
}

func (f *FlowSensitiveInlineVariables) getCodingConvention() CodingConvention {
	if f.codingConvention == nil {
		f.codingConvention = f.compiler.GetCodingConvention()
	}
	return f.codingConvention
}

func (f *FlowSensitiveInlineVariables) reportChange(n rhino.Node) {
	if f.parallelMode {
		if changeRoot := ChangeTracker_GetEnclosingChangeScopeRoot(n); !changeRoot.IsNil() {
			f.parallelChanges = append(f.parallelChanges, changeRoot)
		}
		return
	}
	f.compiler.ReportChangeToEnclosingScope(n)
}

func (f *FlowSensitiveInlineVariables) getReachingUses() *MaybeReachingVariableUse {
	if f.reachingUses == nil {
		f.candidateVarsMap.Clear()
		for i := range f.candidates {
			cand := &f.candidates[i]
			if cand.removed {
				continue
			}
			if v := f.allVarsInFn.Get(cand.varName); v != nil {
				f.candidateVarsMap.Put(cand.varName, v)
			}
		}
		f.reachingUsesStorage.Reset(f.cfg, f.escaped, &f.candidateVarsMap)
		f.reachingUses = &f.reachingUsesStorage
		f.reachingUses.Analyze()
	}
	return f.reachingUses
}

// Java: FlowSensitiveInlineVariables(AbstractCompiler)
func NewFlowSensitiveInlineVariables(compiler AbstractCompiler) *FlowSensitiveInlineVariables {
	f := &FlowSensitiveInlineVariables{}
	f.compiler = compiler
	f.sideEffectPredicate = FlowSensitiveInlineVariables_SideEffectPredicate{outer: f}
	f.sideEffectPredicateFn = f.sideEffectPredicate.Apply
	f.gatherCb = FlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback{outer: f}
	f.gatherCandidates = FlowSensitiveInlineVariables_GatherCandidates{
		outer:    f,
		gatherCb: &f.gatherCb,
	}
	return f
}

func newFlowSensitiveInlineVariablesWorker(parent *FlowSensitiveInlineVariables, compiler AbstractCompiler, astAnalyzer *AstAnalyzer, codingConvention CodingConvention, dummyGlobalScope *Scope, globalRoot rhino.Node) *FlowSensitiveInlineVariables {
	f := NewFlowSensitiveInlineVariables(compiler)
	f.astAnalyzer = astAnalyzer
	f.codingConvention = codingConvention
	f.dummyGlobalScope = dummyGlobalScope
	f.globalRoot = globalRoot
	f.sharedParent = parent
	f.parallelMode = true
	return f
}

// Java: FlowSensitiveInlineVariables.shouldTraverse(NodeTraversal, Node, Node)
func (f *FlowSensitiveInlineVariables) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	return !n.IsScript() || !t.GetInput().IsExtern()
}

func fsivSubtreeHasBlockScopeDecls(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_LET, rhino.Token_CONST, rhino.Token_CLASS, rhino.Token_CATCH:
		return true
	case rhino.Token_FUNCTION:
		return !NodeUtil_IsFunctionExpression(n)
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if fsivSubtreeHasBlockScopeDecls(c) {
			return true
		}
	}
	return false
}

func fsivHasNestedBlockScopeDecls(fn ast.Function, fnBody ast.Block) bool {
	for p := range fn.Params().Params() {
		if _, ok := p.AsName(); !ok && fsivSubtreeHasBlockScopeDecls(p.Node()) {
			return true
		}
	}
	for stmt := range fnBody.Statements() {
		switch stmt.Token() {
		case rhino.Token_FUNCTION:
			continue
		case rhino.Token_CLASS, rhino.Token_CATCH:
			return true
		default:
			for c := stmt.Node().GetFirstChild(); !c.IsNil(); c = c.GetNext() {
				if fsivSubtreeHasBlockScopeDecls(c) {
					return true
				}
			}
		}
	}
	return false
}

func fsivInspectRhs(n rhino.Node) (unsafe bool, hasFunc bool) {
	switch n.GetToken() {
	case rhino.Token_GETELEM,
		rhino.Token_GETPROP,
		rhino.Token_OPTCHAIN_GETPROP,
		rhino.Token_OPTCHAIN_GETELEM,
		rhino.Token_CLASS,
		rhino.Token_ARRAYLIT,
		rhino.Token_OBJECTLIT,
		rhino.Token_REGEXP,
		rhino.Token_NEW:
		return true, false
	case rhino.Token_FUNCTION:
		return false, true
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		u, hf := fsivInspectRhs(c)
		if u {
			return true, false
		}
		if hf {
			hasFunc = true
		}
	}
	return false, hasFunc
}

// fsivIsReadName reports whether the name n (whose parent is parent) is
// purely a read: not assigned, declared, incremented, a parameter, a catch
// binding or a destructuring target.
func fsivIsReadName(n, parent rhino.Node) bool {
	if parent.IsNil() {
		return false
	}
	switch ast.RoleOf(n) {
	case ast.RoleAssignTarget, ast.RoleAssignOpTarget,
		ast.RoleNameDeclarationDeclarators,
		ast.RoleIncDecTarget,
		ast.RoleParamListParams,
		ast.RoleCatchParam:
		return false
	}
	return !NodeUtil_IsLhsByDestructuring(n)
}

func fsivRhsDependsOnOuterOrCatchVars(n, parent rhino.Node, allVarsInFn *jutil.LinkedHashMap[string, *Var]) bool {
	if !parent.IsNil() && ControlFlowGraph_IsEnteringNewCfgNode(n) {
		return false
	}
	if n.IsName() {
		v := allVarsInFn.Get(n.GetString())
		if v == nil || v.GetScope().(*Scope).IsCatchScope() {
			return true
		}
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if fsivRhsDependsOnOuterOrCatchVars(c, n, allVarsInFn) {
			return true
		}
	}
	return false
}

func (f *FlowSensitiveInlineVariables) maybeAddCandidateDef(varName string, rhs ast.Expr) {
	if varName == "" {
		return
	}
	unsafe, hasFunc := fsivInspectRhs(rhs.Node())
	if unsafe {
		return
	}
	if f.getCodingConvention().IsExportedWithLocal(varName, true) {
		return
	}
	if f.getAstAnalyzer().MayHaveSideEffects(rhs.Node()) {
		return
	}
	f.inlinableDefs = append(f.inlinableDefs, fsivCandidateDef{
		varName: varName,
		rhs:     rhs,
		hasFunc: hasFunc,
	})
}

// scanStatementsForCandidateDefs collects the definitions under n that
// the pass could inline: a declared name with a value (var x = v) or an
// assignment statement to a name (x = v). It walks the statements that
// get CFG nodes, and of a for loop its init or target, where a var can be.
func (f *FlowSensitiveInlineVariables) scanStatementsForCandidateDefs(n rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_BLOCK, rhino.Token_SWITCH_BODY, rhino.Token_TRY:
		for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			f.scanStatementsForCandidateDefs(c)
		}
	case rhino.Token_IF:
		ifs := ast.MustIf(n)
		f.scanStatementsForCandidateDefs(ifs.Then().Node())
		if els, ok := ifs.Else(); ok {
			f.scanStatementsForCandidateDefs(els.Node())
		}
	case rhino.Token_WHILE, rhino.Token_DO:
		f.scanStatementsForCandidateDefs(ast.MustLoop(n).Body().Node())
	case rhino.Token_WITH:
		f.scanStatementsForCandidateDefs(ast.MustWith(n).Body().Node())
	case rhino.Token_SWITCH:
		f.scanStatementsForCandidateDefs(ast.MustSwitch(n).Body().Node())
	case rhino.Token_CASE, rhino.Token_DEFAULT_CASE:
		f.scanStatementsForCandidateDefs(ast.MustSwitchCase(n).Body().Node())
	case rhino.Token_CATCH:
		f.scanStatementsForCandidateDefs(ast.MustCatch(n).Body().Node())
	case rhino.Token_LABEL:
		f.scanStatementsForCandidateDefs(ast.MustLabel(n).Body().Node())
	case rhino.Token_FOR:
		loop := ast.MustFor(n)
		if init, ok := loop.Init(); ok {
			f.scanStatementsForCandidateDefs(init.Node())
		}
		f.scanStatementsForCandidateDefs(loop.Body().Node())
	case rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
		loop := ast.MustEnhancedFor(n)
		f.scanStatementsForCandidateDefs(loop.Target().Node())
		f.scanStatementsForCandidateDefs(loop.AsLoop().Body().Node())
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		for d := range ast.MustNameDeclaration(n).Declarators() {
			if name, ok := d.AsName(); ok {
				if value, ok := name.Init(); ok {
					f.maybeAddCandidateDef(name.Ident(), value)
				}
			}
		}
	case rhino.Token_EXPR_RESULT:
		if assign, ok := ast.MustExprStmt(n).Expr().AsAssign(); ok {
			if lhs, ok := assign.Target().AsName(); ok {
				f.maybeAddCandidateDef(lhs.Ident(), assign.Value())
			}
		}
	case rhino.Token_FUNCTION, rhino.Token_RETURN, rhino.Token_THROW,
		rhino.Token_BREAK, rhino.Token_CONTINUE, rhino.Token_DEBUGGER,
		rhino.Token_EMPTY, rhino.Token_CLASS, rhino.Token_EXPORT, rhino.Token_IMPORT:
		// ControlFlowAnalysis does not enter expressions, functions, or classes for statements.
	default:
		for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			f.scanStatementsForCandidateDefs(c)
		}
	}
}

func (f *FlowSensitiveInlineVariables) markReadsOutsideLoop(n rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_FUNCTION,
		rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF,
		rhino.Token_DO, rhino.Token_WHILE,
		rhino.Token_CLASS, rhino.Token_IMPORT:
		return
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if c.IsName() && fsivIsReadName(c, n) {
			name := c.GetString()
			for i := range f.inlinableDefs {
				if f.inlinableDefs[i].varName == name && !f.inlinableDefs[i].hasRead {
					f.inlinableDefs[i].hasRead = true
					f.unreadDefCount--
				}
			}
			if f.unreadDefCount == 0 {
				return
			}
		}
		if c.HasChildren() {
			f.markReadsOutsideLoop(c)
			if f.unreadDefCount == 0 {
				return
			}
		}
	}
}

func (f *FlowSensitiveInlineVariables) collectInlinableDefsWithReads(fnBody ast.Block) (ok bool, hasFuncRhs bool) {
	f.inlinableDefs = f.inlinableDefs[:0]
	f.scanStatementsForCandidateDefs(fnBody.Node())
	if len(f.inlinableDefs) == 0 {
		return false, false
	}
	f.unreadDefCount = len(f.inlinableDefs)
	f.markReadsOutsideLoop(fnBody.Node())
	nDefs := 0
	for _, d := range f.inlinableDefs {
		if d.hasRead {
			if d.hasFunc {
				hasFuncRhs = true
			}
			f.inlinableDefs[nDefs] = d
			nDefs++
		}
	}
	f.inlinableDefs = f.inlinableDefs[:nDefs]
	return nDefs > 0, hasFuncRhs
}

func (f *FlowSensitiveInlineVariables) isInlinableDefVar(name string) bool {
	if len(f.inlinableDefs) == 0 {
		return true
	}
	for i := range f.inlinableDefs {
		if f.inlinableDefs[i].varName == name {
			return true
		}
	}
	return false
}

// Java: FlowSensitiveInlineVariables.enterScope(NodeTraversal)
func (f *FlowSensitiveInlineVariables) EnterScope(t *NodeTraversal) {
	if len(t.scopes) <= 1 {
		return // Don't even bother. All global variables are likely escaped.
	}
	fnBody, ok := ast.AsBlock(t.GetScopeRoot())
	if !ok {
		return
	}
	functionScopeRoot, ok := ast.AsFunction(fnBody.Node().GetParent())
	if !ok {
		return // Only want to do the following if it's a function block scope.
	}

	f.curTraversal = t
	defer func() {
		f.curTraversal = nil
	}()
	scopeCreator := t.GetScopeCreator().(*SyntacticScopeCreator)
	f.processFunction(functionScopeRoot, fnBody, scopeCreator)
}

func (f *FlowSensitiveInlineVariables) processFunction(fn ast.Function, body ast.Block, scopeCreator *SyntacticScopeCreator) {
	if ok, _ := f.collectInlinableDefsWithReads(body); !ok {
		return
	}
	functionScopeRoot, fnBody := fn.Node(), body.Node()

	var fnScope, scope *Scope
	if f.dummyGlobalScope != nil {
		fnScope = scopeCreator.CreateScopeOfScope(functionScopeRoot, f.dummyGlobalScope)
		scope = scopeCreator.CreateScopeOfScope(fnBody, fnScope)
	} else if f.curTraversal != nil {
		scope = f.curTraversal.GetScope()
		fnScope = scope.GetParent()
	} else {
		root := functionScopeRoot
		for !root.GetParent().IsNil() {
			root = root.GetParent()
		}
		f.dummyGlobalScope = Scope_CreateGlobalScope(root)
		fnScope = scopeCreator.CreateScopeOfScope(functionScopeRoot, f.dummyGlobalScope)
		scope = scopeCreator.CreateScopeOfScope(fnBody, fnScope)
	}
	if LiveVariablesAnalysis_MAX_VARIABLES_TO_ANALYZE < scope.GetVarCount() {
		return
	}

	var allVarsInFn *jutil.LinkedHashMap[string, *Var]
	if !fsivHasNestedBlockScopeDecls(fn, body) {
		f.allVarsMap.Clear()
		for _, v := range fnScope.varList {
			if v != nil {
				f.allVarsMap.Put(v.GetName(), v)
			}
		}
		for _, v := range scope.varList {
			if v != nil {
				f.allVarsMap.Put(v.GetName(), v)
			}
		}
		allVarsInFn = &f.allVarsMap
	} else {
		allVarsDeclaredInFunction := f.varsFinder.collect(f.compiler, scopeCreator, fnScope)
		allVarsInFn = allVarsDeclaredInFunction.GetAllVariables()
	}
	if allVarsInFn.IsEmpty() {
		return
	}

	nDefs := 0
	for _, d := range f.inlinableDefs {
		if allVarsInFn.Get(d.varName) != nil && !fsivRhsDependsOnOuterOrCatchVars(d.rhs.Node(), rhino.NilNode, allVarsInFn) {
			f.inlinableDefs[nDefs] = d
			nDefs++
		}
	}
	f.inlinableDefs = f.inlinableDefs[:nDefs]
	if nDefs == 0 {
		return
	}

	f.escapedSet.Clear()
	escaped := &f.escapedSet
	DataFlowAnalysis_ComputeEscaped(
		fnScope, escaped, f.compiler, scopeCreator, allVarsInFn)
	nDefs = 0
	for _, d := range f.inlinableDefs {
		if v := allVarsInFn.Get(d.varName); v != nil && !escaped.Contains(v) {
			f.inlinableDefs[nDefs] = d
			nDefs++
		}
	}
	f.inlinableDefs = f.inlinableDefs[:nDefs]
	if nDefs == 0 {
		return
	}

	f.curFn = functionScopeRoot
	f.curFnEnclosingScope = nil
	f.scopeCreator = scopeCreator
	defer func() {
		clear(f.candidates)
		f.candidates = f.candidates[:0]
		f.curFn = rhino.NilNode
		f.curFnEnclosingScope = nil
		f.cfg = nil
		f.reusableCfg.release()
	}()

	// Compute the forward reaching definition.
	f.cfg = f.reusableCfg.Compute(
		f.compiler,
		functionScopeRoot,
		/* shouldTraverseFunctions= */ false,
		/* edgeAnnotations= */ true)

	f.reachingDefStorage.Reset(f.cfg, f.compiler, escaped, allVarsInFn)
	f.reachingDef = &f.reachingDefStorage
	f.reachingDef.Analyze()
	f.candidates = f.candidates[:0]

	// Using the forward reaching definition search to find all the inline
	// candidates
	fsivTraverseGatherCandidates(f.compiler, fnBody, &f.gatherCandidates)

	if len(f.candidates) == 0 {
		return
	}

	// Lazily compute the backward reaching use only if a candidate passes the
	// syntactic and side-effect checks in Candidate.CanInline.
	f.reachingUses = nil
	f.escaped = escaped
	f.allVarsInFn = allVarsInFn
	for i := range f.candidates {
		c := &f.candidates[i]
		if c.removed {
			continue
		}
		candidateVar := jutil.CheckNotNull(allVarsInFn.Get(c.varName))
		if c.CanInline(candidateVar.GetScope().(*Scope)) {
			c.InlineVariable()
			c.removed = true

			// If candidate "c" has dependencies, then inlining it may have introduced new dependencies
			// for our other inlining candidates. MustBeReachingVariableDef uses a dependency graph in
			// its analysis. Generating a new dependency graph will need another CFG computation.
			// Ideally we should iterate to a fixed point, but that can be costly. Therefore, we use
			// a conservative heuristic here: For each candidate "other", we back off if its set of
			// dependencies cannot contain all of "c"'s dependencies.
			if c.defMetadata.HasDependencies() {
				if scopeVarID, ok := f.reachingDef.ResolveVarID(scope.GetVar(c.varName)); ok {
					for j := i + 1; j < len(f.candidates); j++ {
						other := &f.candidates[j]
						if !other.removed &&
							other.defMetadata.DependsOnVarID(scopeVarID) &&
							!other.defMetadata.ContainsAllDependencies(c.defMetadata) {
							other.removed = true
						}
					}
				}
			}
		} else {
			c.removed = true
		}
	}
}

// Java: FlowSensitiveInlineVariables.isCandidateFunction(Node)
func (f *FlowSensitiveInlineVariables) IsCandidateFunction(fn rhino.Node) bool {
	fnBody := ast.MustFunction(fn).Body()
	return FlowSensitiveInlineVariables_ContainsCandidateExpressions(fnBody.Node())
}

// Java: FlowSensitiveInlineVariables.containsCandidateExpressions(Node)
func FlowSensitiveInlineVariables_ContainsCandidateExpressions(n rhino.Node) bool {
	if n.IsFunction() {
		// don't recurse into inner functions or into expressions that can't contain declarations.
		return false
	}

	if decl, ok := ast.AsNameDeclaration(n); ok {
		if first, ok := decl.FirstDeclarator(); ok && first.Token() == rhino.Token_NAME {
			return true
		}
	} else if FlowSensitiveInlineVariables_IsAssignmentToName(n) {
		return true
	}

	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if FlowSensitiveInlineVariables_ContainsCandidateExpressions(c) {
			return true
		}
	}
	return false
}

// Java: FlowSensitiveInlineVariables.isAssignmentToName(Node)
func FlowSensitiveInlineVariables_IsAssignmentToName(n rhino.Node) bool {
	var target ast.Target
	if assign, ok := ast.AsAssignment(n); ok {
		target = assign.Target()
	} else if incDec, ok := ast.AsIncDec(n); ok {
		target = incDec.Target()
	} else {
		return false
	}
	return target.Token() == rhino.Token_NAME
}

// Java: FlowSensitiveInlineVariables.exitScope(NodeTraversal)
func (f *FlowSensitiveInlineVariables) ExitScope(t *NodeTraversal) {}

// collectCandidateFuncs collects the functions under n that have a
// definition the pass could inline, with their nesting depth among such
// functions, in the order NodeTraversal enters their scopes (a class's
// extends before its body, the computed keys of a class body before its
// members). The names of functions and classes hold nothing to collect.
func (f *FlowSensitiveInlineVariables) collectCandidateFuncs(n rhino.Node, depth int) {
	switch n.GetToken() {
	case rhino.Token_FUNCTION:
		fn := ast.MustFunction(n)
		if fnBody, ok := fn.Body().AsBlock(); ok {
			if ok, hasFuncRhs := f.collectInlinableDefsWithReads(fnBody); ok {
				if depth > f.maxCandDepth {
					f.maxCandDepth = depth
				}
				f.candFuncs = append(f.candFuncs, fsivCandFunc{
					fn:         fn,
					body:       fnBody,
					depth:      depth,
					hasFuncRhs: hasFuncRhs,
				})
				depth++
			}
		}
		f.collectCandidateFuncs(fn.Params().Node(), depth)
		f.collectCandidateFuncs(fn.Body().Node(), depth)
		return
	case rhino.Token_CLASS:
		class := ast.MustClass(n)
		if extends, ok := class.Extends(); ok {
			f.collectCandidateFuncs(extends.Node(), depth)
		}
		f.collectCandidateFuncs(class.Body().Node(), depth)
		return
	case rhino.Token_CLASS_MEMBERS:
		members := ast.MustClassBody(n)
		for m := range members.Members() {
			switch m.Token() {
			case rhino.Token_COMPUTED_PROP:
				f.collectCandidateFuncs(m.MustComputedProp().Key().Node(), depth)
			case rhino.Token_COMPUTED_FIELD_DEF:
				f.collectCandidateFuncs(m.MustComputedFieldDef().Key().Node(), depth)
			}
		}
		for m := range members.Members() {
			switch m.Token() {
			case rhino.Token_COMPUTED_PROP:
				if value, ok := m.MustComputedProp().Value(); ok {
					f.collectCandidateFuncs(value.Node(), depth)
				}
			case rhino.Token_COMPUTED_FIELD_DEF:
				if value, ok := m.MustComputedFieldDef().Value(); ok {
					f.collectCandidateFuncs(value.Node(), depth)
				}
			default:
				f.collectCandidateFuncs(m.Node(), depth)
			}
		}
		return
	}
	for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		f.collectCandidateFuncs(child, depth)
	}
}

// Java: FlowSensitiveInlineVariables.process(Node, Node)
func (f *FlowSensitiveInlineVariables) Process(externs rhino.Node, root rhino.Node) {
	// Pre-populate cached AstAnalyzer and CodingConvention on the main goroutine
	// so worker goroutines never race on Compiler.cachedAstAnalyzer initialization.
	astAnalyzer := f.getAstAnalyzer()
	codingConvention := f.getCodingConvention()

	f.candFuncs = f.candFuncs[:0]
	f.maxCandDepth = 0
	if !root.IsNil() {
		f.collectCandidateFuncs(root, 0)
	}
	if len(f.candFuncs) == 0 {
		return
	}

	globalRoot := root
	if !externs.IsNil() && !externs.GetParent().IsNil() {
		globalRoot = externs.GetParent()
	} else {
		for !globalRoot.GetParent().IsNil() {
			globalRoot = globalRoot.GetParent()
		}
	}
	f.globalRoot = globalRoot
	f.dummyGlobalScope = Scope_CreateGlobalScope(globalRoot)
	f.globalScopeOnce = sync.Once{}
	f.realGlobalScope = nil
	f.globalScopeCreator = nil

	sc := NewSyntacticScopeCreator(f.compiler)
	f.scopeCreator = sc
	defer func() {
		if f.globalScopeCreator != nil {
			f.globalScopeCreator = nil
		}
		f.realGlobalScope = nil
		f.scopeCreator = nil
		f.dummyGlobalScope = nil
		f.globalRoot = rhino.NilNode
		f.escapedSet.Clear()
	}()

	const parallelThreshold = 8
	numWorkers := 1
	if opts := f.compiler.GetOptions(); opts != nil && opts.GetNumParallelThreads() > 1 {
		numWorkers = opts.GetNumParallelThreads()
	}

	var levelFuncs []fsivCandFunc
	for depth := 0; depth <= f.maxCandDepth; depth++ {
		levelFuncs = levelFuncs[:0]
		for _, cf := range f.candFuncs {
			if cf.depth == depth {
				levelFuncs = append(levelFuncs, cf)
			}
		}
		if len(levelFuncs) == 0 {
			continue
		}

		// Count how many functions at this depth can run in parallel (those without FUNCTION in RHS).
		parallelCount := 0
		for _, cf := range levelFuncs {
			if !cf.hasFuncRhs {
				parallelCount++
			}
		}

		if parallelCount < parallelThreshold || numWorkers == 1 {
			f.parallelMode = false
			for _, cf := range levelFuncs {
				f.processFunction(cf.fn, cf.body, sc)
			}
			continue
		}

		// Run any hasFuncRhs functions sequentially first on the main goroutine.
		f.parallelMode = false
		for _, cf := range levelFuncs {
			if cf.hasFuncRhs {
				f.processFunction(cf.fn, cf.body, sc)
			}
		}

		// Collect parallelizable functions at this depth.
		parFuncs := make([]fsivCandFunc, 0, parallelCount)
		for _, cf := range levelFuncs {
			if !cf.hasFuncRhs {
				parFuncs = append(parFuncs, cf)
			}
		}

		workers := numWorkers
		if workers > len(parFuncs) {
			workers = len(parFuncs)
		}
		changedByFunc := make([][]rhino.Node, len(parFuncs))
		var nextIdx atomic.Int64
		var wg sync.WaitGroup
		var panicVal any
		var panicOnce sync.Once

		wg.Add(workers)
		for w := 0; w < workers; w++ {
			go func() {
				defer wg.Done()
				defer func() {
					if r := recover(); r != nil {
						panicOnce.Do(func() {
							panicVal = r
						})
					}
				}()
				worker := newFlowSensitiveInlineVariablesWorker(
					f, f.compiler, astAnalyzer, codingConvention, f.dummyGlobalScope, f.globalRoot)
				wsc := NewSyntacticScopeCreator(f.compiler)
				worker.scopeCreator = wsc
				defer func() {
					worker.scopeCreator = nil
					worker.escapedSet.Clear()
				}()

				for {
					idx := int(nextIdx.Add(1) - 1)
					if idx >= len(parFuncs) {
						break
					}
					worker.parallelChanges = worker.parallelChanges[:0]
					cf := parFuncs[idx]
					worker.processFunction(cf.fn, cf.body, wsc)
					if len(worker.parallelChanges) > 0 {
						nodes := make([]rhino.Node, len(worker.parallelChanges))
						copy(nodes, worker.parallelChanges)
						changedByFunc[idx] = nodes
					}
				}
			}()
		}
		wg.Wait()
		if panicVal != nil {
			panic(panicVal)
		}
		for _, nodes := range changedByFunc {
			for _, n := range nodes {
				f.compiler.ReportChangeToChangeScope(n)
			}
		}
	}
}

// Java: FlowSensitiveInlineVariables.visit(NodeTraversal, Node, Node)
func (f *FlowSensitiveInlineVariables) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {}

func (f *FlowSensitiveInlineVariables) getRealGlobalScope(root rhino.Node) *Scope {
	if f.sharedParent != nil {
		return f.sharedParent.getRealGlobalScope(root)
	}
	if !f.globalRoot.IsNil() && root == f.globalRoot {
		f.globalScopeOnce.Do(func() {
			gsc := NewSyntacticScopeCreator(f.compiler)
			f.globalScopeCreator = gsc
			f.realGlobalScope = gsc.CreateScopeOfScope(f.globalRoot, nil)
		})
		return f.realGlobalScope
	}
	return f.scopeCreator.CreateScopeOfScope(root, nil)
}

func (f *FlowSensitiveInlineVariables) buildEnclosingScope() *Scope {
	if f.curFn.IsNil() || f.scopeCreator == nil {
		return nil
	}
	var rootsBuf [16]rhino.Node
	roots := rootsBuf[:0]
	prev := f.curFn
	for p := prev.GetParent(); !p.IsNil(); prev, p = p, p.GetParent() {
		switch ast.RoleOf(prev) {
		case ast.RoleClassExtends, ast.RoleComputedFieldDefKey:
			// Evaluated in the scope outside the class.
			continue
		}
		if (p.IsRoot() && p.GetParent().IsNil()) || p.GetParent().IsNil() || p.IsFunction() || NodeUtil_CreatesBlockScope(p) {
			if !p.GetParent().IsNil() &&
				((p.IsBlock() && !NodeUtil_IsFunctionBlock(p) && !p.IsClassMembers() && !NodeUtil_IsClassStaticBlock(p) && !NodeUtil_IsBundledGoogModuleScopeRoot(p)) || NodeUtil_IsAnyFor(p)) &&
				!syntacticScopeHasBlockDecls(p) {
				continue
			}
			roots = append(roots, p)
		}
	}
	var parentScope *Scope
	for i := len(roots) - 1; i >= 0; i-- {
		r := roots[i]
		if parentScope == nil && r.GetParent().IsNil() {
			parentScope = f.getRealGlobalScope(r)
			continue
		}
		parentScope = f.scopeCreator.CreateScopeOfScope(r, parentScope)
	}
	return parentScope
}

func (f *FlowSensitiveInlineVariables) hasEnclosingSlot(name string) bool {
	if f.curFnEnclosingScope == nil {
		f.curFnEnclosingScope = f.buildEnclosingScope()
	}
	return f.curFnEnclosingScope != nil && f.curFnEnclosingScope.HasSlot(name)
}

// FlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback is an
// AbstractCfgNodeTraversalCallback.
type FlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback struct {
	ControlFlowGraph_AbstractCfgNodeTraversalCallbackBase
	outer *FlowSensitiveInlineVariables

	cfgNode rhino.Node
	inState *MustBeReachingVariableDef_MustDef
}

// Java: FlowSensitiveInlineVariables.GatherCandidatesCfgNodeCallback()
func NewFlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback(outer *FlowSensitiveInlineVariables) *FlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback {
	return &FlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback{outer: outer, cfgNode: rhino.NilNode}
}

// Java: FlowSensitiveInlineVariables.GatherCandidatesCfgNodeCallback.setCfgNode(Node)
func (cb *FlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback) SetCfgNode(cfgNode rhino.Node) {
	cb.cfgNode = cfgNode
	cb.inState = nil
}

// Java: FlowSensitiveInlineVariables.GatherCandidatesCfgNodeCallback.visit(NodeTraversal, Node, Node)
func (cb *FlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsName() {
		// Make sure that the name node is purely a read.
		if !fsivIsReadName(n, parent) {
			return
		}

		name := n.GetString()
		if len(cb.outer.inlinableDefs) > 0 && !cb.outer.isInlinableDefVar(name) {
			return
		}

		var def *MustBeReachingVariableDef_Definition
		if cb.inState != nil {
			def = cb.outer.reachingDef.getDefFromMustDef(name, cb.inState)
		} else {
			def = cb.outer.reachingDef.GetDef(name, cb.cfgNode)
		}
		if def != nil && !def.node.IsFunction() && !def.node.IsBlock() &&
			!cb.outer.getCodingConvention().IsExportedWithLocal(name /* local= */, true) &&
			!cb.outer.reachingDef.DependsOnOuterScopeVars(def) {
			cb.outer.candidates = append(cb.outer.candidates, FlowSensitiveInlineVariables_Candidate{
				outer:       cb.outer,
				varName:     name,
				defMetadata: def,
				use:         n,
				useCfgNode:  cb.cfgNode,
			})
		}
	}
}

// FlowSensitiveInlineVariables_GatherCandidates gathers a list of possible
// candidates for inlining based only on information from
// MustBeReachingVariableDef. The list will be stored in candidates and the
// validity of each inlining Candidate should be later verified with
// Candidate#canInline(Scope) when MaybeReachingVariableUse has been
// performed.
type FlowSensitiveInlineVariables_GatherCandidates struct {
	NodeTraversal_AbstractShallowCallbackBase
	outer *FlowSensitiveInlineVariables

	gatherCb *FlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback
}

// Java: FlowSensitiveInlineVariables.GatherCandidates()
func NewFlowSensitiveInlineVariables_GatherCandidates(outer *FlowSensitiveInlineVariables) *FlowSensitiveInlineVariables_GatherCandidates {
	return &FlowSensitiveInlineVariables_GatherCandidates{
		outer:    outer,
		gatherCb: NewFlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback(outer),
	}
}

// Java: FlowSensitiveInlineVariables.GatherCandidates.visit(NodeTraversal, Node, Node)
func (g *FlowSensitiveInlineVariables_GatherCandidates) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_BLOCK, rhino.Token_TRY, rhino.Token_SWITCH_BODY, rhino.Token_EMPTY:
		return
	}
	graphNode := g.outer.cfg.Nodes.Get(n)
	if graphNode == nil {
		// Not a CFG node.
		return
	}
	cfgNode := n

	g.gatherCb.cfgNode = cfgNode
	g.gatherCb.inState = g.outer.reachingDef.FlowStateAt(graphNode.Index).GetIn()
	var currentNode rhino.Node
	walkCfgNodeBranch(cfgNode, rhino.NilNode, g.gatherCb, &currentNode)
	g.gatherCb.inState = nil
}

func fsivTraverseGatherCandidates(compiler AbstractCompiler, root rhino.Node, cb *FlowSensitiveInlineVariables_GatherCandidates) {
	jutil.CheckNotNull(compiler)
	if root.IsNil() {
		NodeTraversal_Traverse(compiler, root, cb)
		return
	}
	var currentNode rhino.Node
	jutil.Try(func() {
		if Platform_IsThreadInterrupted() {
			panic(jutil.NewRuntimeExceptionWithCause(jutil.NewInterruptedException()))
		}
		currentNode = root
		fsivWalkGatherCandidates(root, rhino.NilNode, cb, &currentNode)
	}, func(unexpectedException jutil.Throwable) {
		b := NodeTraversal_Builder{compiler: compiler, callback: cb}
		t := newNodeTraversal(&b)
		t.InitTraversal(root)
		t.currentNode = currentNode
		t.ThrowUnexpectedException(unexpectedException)
	})
}

func fsivWalkGatherCandidates(n, parent rhino.Node, cb *FlowSensitiveInlineVariables_GatherCandidates, currentNode *rhino.Node) {
	*currentNode = n
	if !parent.IsNil() && parent.IsFunction() && ast.RoleOf(n) != ast.RoleFunctionName {
		return
	}
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		if Platform_IsThreadInterrupted() {
			panic(jutil.NewRuntimeExceptionWithCause(jutil.NewInterruptedException()))
		}
		for child := n.GetFirstChild(); !child.IsNil(); {
			next := child.GetNext()
			fsivWalkGatherCandidates(child, n, cb, currentNode)
			child = next
		}
	case rhino.Token_FUNCTION:
		// Only the function itself is a CFG node of the enclosing function:
		// its name is no CFG node, and its params and body are not entered
		// (above).
	case rhino.Token_BLOCK, rhino.Token_SWITCH_BODY, rhino.Token_TRY:
		for child := n.GetFirstChild(); !child.IsNil(); {
			next := child.GetNext()
			fsivWalkGatherCandidates(child, n, cb, currentNode)
			child = next
		}
	case rhino.Token_IF:
		ifs := ast.MustIf(n)
		els, hasElse := ifs.Else()
		fsivWalkGatherCandidates(ifs.Then().Node(), n, cb, currentNode)
		if hasElse {
			fsivWalkGatherCandidates(els.Node(), n, cb, currentNode)
		}
	case rhino.Token_WITH:
		fsivWalkGatherCandidates(ast.MustWith(n).Body().Node(), n, cb, currentNode)
	case rhino.Token_SWITCH:
		fsivWalkGatherCandidates(ast.MustSwitch(n).Body().Node(), n, cb, currentNode)
	case rhino.Token_CASE, rhino.Token_DEFAULT_CASE:
		fsivWalkGatherCandidates(ast.MustSwitchCase(n).Body().Node(), n, cb, currentNode)
	case rhino.Token_CATCH:
		fsivWalkGatherCandidates(ast.MustCatch(n).Body().Node(), n, cb, currentNode)
	case rhino.Token_LABEL:
		fsivWalkGatherCandidates(ast.MustLabel(n).Body().Node(), n, cb, currentNode)
	case rhino.Token_WHILE, rhino.Token_DO, rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
		// Any use inside a loop node has NodeUtil_IsWithinLoop(use) == true,
		// which CanInline always rejects. Skip loop subtrees entirely.
		return
	case rhino.Token_EXPR_RESULT, rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST,
		rhino.Token_RETURN, rhino.Token_THROW, rhino.Token_BREAK, rhino.Token_CONTINUE,
		rhino.Token_DEBUGGER, rhino.Token_EMPTY, rhino.Token_CLASS, rhino.Token_EXPORT, rhino.Token_IMPORT:
		// ControlFlowAnalysis does not enter expressions or classes, so no child of
		// these nodes is in the CFG.
	default:
		for child := n.GetFirstChild(); !child.IsNil(); {
			next := child.GetNext()
			fsivWalkGatherCandidates(child, n, cb, currentNode)
			child = next
		}
	}
	*currentNode = n
	cb.Visit(nil, n, parent)
}
