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
	nameID  uint32
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

	cfg            *ControlFlowGraph[rhino.Node]
	candidates     []FlowSensitiveInlineVariables_Candidate
	reachingDef    *MustBeReachingVariableDef
	reachingUses   *MaybeReachingVariableUse
	activeVarTable *fnVarTable

	reusableCfg         ReusableCfg
	reachingDefStorage  MustBeReachingVariableDef
	reachingUsesStorage MaybeReachingVariableUse
	varTable            fnVarTable
	candVarSeen         inlineBitSet
	candVarIDs          []uint32
	gatherCfgNode       rhino.Node
	gatherInState       *MustBeReachingVariableDef_MustDef
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
	parallelMode        bool
	parallelChanges     []rhino.Node
}

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
		table := f.activeVarTable
		f.candVarSeen.clear(len(table.vars) <= 256)
		f.candVarIDs = f.candVarIDs[:0]
		for i := range f.candidates {
			cand := &f.candidates[i]
			if cand.removed {
				continue
			}
			idx, ok := table.lookupIndexByID(cand.nameID)
			if !ok || f.candVarSeen.Get(idx) {
				continue
			}
			f.candVarSeen.Set(idx)
			if !table.escaped.Get(idx) {
				f.candVarIDs = append(f.candVarIDs, table.vars[idx].nameID)
			}
		}
		f.reachingUsesStorage.resetWithVarIDs(f.cfg, f.candVarIDs)
		f.reachingUses = &f.reachingUsesStorage
		f.reachingUses.Analyze()
	}
	return f.reachingUses
}

func (f *FlowSensitiveInlineVariables) sideEffectPredicateFunc() func(rhino.Node) bool {
	if f.sideEffectPredicateFn == nil {
		f.sideEffectPredicateFn = f.sideEffectPredicate.Apply
	}
	return f.sideEffectPredicateFn
}

// Java: FlowSensitiveInlineVariables(AbstractCompiler)
func NewFlowSensitiveInlineVariables(compiler AbstractCompiler) *FlowSensitiveInlineVariables {
	f := &FlowSensitiveInlineVariables{compiler: compiler}
	f.sideEffectPredicate.outer = f
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

func fsivRhsDependsOnOuterOrCatchVars(n, parent rhino.Node, table *fnVarTable) bool {
	if !parent.IsNil() && ControlFlowGraph_IsEnteringNewCfgNode(n) {
		return false
	}
	if n.IsName() {
		v := table.lookupVarByID(n.GetStringID())
		if v == nil || v.GetScope().(*Scope).IsCatchScope() {
			return true
		}
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if fsivRhsDependsOnOuterOrCatchVars(c, n, table) {
			return true
		}
	}
	return false
}

func (f *FlowSensitiveInlineVariables) maybeAddCandidateDef(name ast.Name, rhs ast.Expr) {
	nameID := name.Node().GetStringID()
	if nameID <= rhino.RhinoStringPool_EmptyID {
		return
	}
	unsafe, hasFunc := fsivInspectRhs(rhs.Node())
	if unsafe {
		return
	}
	if f.getCodingConvention().IsExportedWithLocal(rhino.RhinoStringPool_StringFromID(nameID), true) {
		return
	}
	if f.getAstAnalyzer().MayHaveSideEffects(rhs.Node()) {
		return
	}
	f.inlinableDefs = append(f.inlinableDefs, fsivCandidateDef{
		nameID:  nameID,
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
					f.maybeAddCandidateDef(name, value)
				}
			}
		}
	case rhino.Token_EXPR_RESULT:
		if assign, ok := ast.MustExprStmt(n).Expr().AsAssign(); ok {
			if lhs, ok := assign.Target().AsName(); ok {
				f.maybeAddCandidateDef(lhs, assign.Value())
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
			nameID := c.GetStringID()
			for i := range f.inlinableDefs {
				if f.inlinableDefs[i].nameID == nameID && !f.inlinableDefs[i].hasRead {
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

func (f *FlowSensitiveInlineVariables) isInlinableDefVarID(nameID uint32) bool {
	if len(f.inlinableDefs) == 0 {
		return true
	}
	for i := range f.inlinableDefs {
		if f.inlinableDefs[i].nameID == nameID {
			return true
		}
	}
	return false
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

	var table *fnVarTable
	if !fsivHasNestedBlockScopeDecls(fn, body) {
		f.varTable.reset()
		f.varTable.appendScopeVars(fnScope)
		f.varTable.appendScopeVars(scope)
		f.varTable.finalize()
		table = &f.varTable
	} else {
		allVarsDeclaredInFunction := f.varsFinder.collect(f.compiler, scopeCreator, fnScope)
		table = &allVarsDeclaredInFunction.table
	}
	if table.isEmpty() {
		return
	}

	nDefs := 0
	for _, d := range f.inlinableDefs {
		if table.lookupVarByID(d.nameID) != nil && !fsivRhsDependsOnOuterOrCatchVars(d.rhs.Node(), rhino.NilNode, table) {
			f.inlinableDefs[nDefs] = d
			nDefs++
		}
	}
	f.inlinableDefs = f.inlinableDefs[:nDefs]
	if nDefs == 0 {
		return
	}

	table.escaped.clear(len(table.vars) <= 256)
	computeEscapedInTable(fnScope, table, f.compiler, scopeCreator)
	nDefs = 0
	for _, d := range f.inlinableDefs {
		if idx, ok := table.lookupIndexByID(d.nameID); ok && !table.escaped.Get(idx) {
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
		f.activeVarTable = nil
		f.cfg = nil
		f.reusableCfg.release()
	}()

	// Compute the forward reaching definition.
	f.cfg = f.reusableCfg.Compute(
		f.compiler,
		functionScopeRoot,
		/* shouldTraverseFunctions= */ false,
		/* edgeAnnotations= */ true)

	f.reachingDefStorage.resetWithTable(f.cfg, f.compiler, table)
	f.reachingDef = &f.reachingDefStorage
	f.reachingDef.Analyze()
	f.candidates = f.candidates[:0]

	// Using the forward reaching definition search to find all the inline
	// candidates
	f.traverseGatherCandidates(fnBody)

	if len(f.candidates) == 0 {
		return
	}

	// Lazily compute the backward reaching use only if a candidate passes the
	// syntactic and side-effect checks in Candidate.CanInline.
	f.reachingUses = nil
	f.activeVarTable = table
	for i := range f.candidates {
		c := &f.candidates[i]
		if c.removed {
			continue
		}
		candidateVar := jutil.CheckNotNull(table.lookupVarByID(c.nameID))
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
				if scopeVarID, ok := f.reachingDef.ResolveVarID(scope.GetVarByID(c.nameID)); ok {
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

		type fsivWorkerState struct {
			worker *FlowSensitiveInlineVariables
			sc     *SyntacticScopeCreator
		}
		changedByFunc := make([][]rhino.Node, len(parFuncs))
		runParallelWorkersIndexed(
			len(parFuncs),
			numWorkers,
			func(_ int) fsivWorkerState {
				worker := newFlowSensitiveInlineVariablesWorker(
					f, f.compiler, astAnalyzer, codingConvention, f.dummyGlobalScope, f.globalRoot)
				wsc := NewSyntacticScopeCreator(f.compiler)
				worker.scopeCreator = wsc
				return fsivWorkerState{worker: worker, sc: wsc}
			},
			func(ws *fsivWorkerState, idx int) {
				ws.worker.parallelChanges = ws.worker.parallelChanges[:0]
				cf := parFuncs[idx]
				ws.worker.processFunction(cf.fn, cf.body, ws.sc)
				if len(ws.worker.parallelChanges) > 0 {
					nodes := make([]rhino.Node, len(ws.worker.parallelChanges))
					copy(nodes, ws.worker.parallelChanges)
					changedByFunc[idx] = nodes
				}
			},
			func(_ int, ws *fsivWorkerState) {
				ws.worker.scopeCreator = nil
			},
		)
		for _, nodes := range changedByFunc {
			for _, n := range nodes {
				f.compiler.ReportChangeToChangeScope(n)
			}
		}
	}
}

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

func (f *FlowSensitiveInlineVariables) hasEnclosingSlotByID(nameID uint32) bool {
	if f.curFnEnclosingScope == nil {
		f.curFnEnclosingScope = f.buildEnclosingScope()
	}
	return f.curFnEnclosingScope != nil && f.curFnEnclosingScope.GetVarByID(nameID) != nil
}

func (f *FlowSensitiveInlineVariables) visitCandidateName(n rhino.Node, parent rhino.Node) {
	// Make sure that the name node is purely a read.
	if !fsivIsReadName(n, parent) {
		return
	}

	nameID := n.GetStringID()
	if len(f.inlinableDefs) > 0 && !f.isInlinableDefVarID(nameID) {
		return
	}

	var def *MustBeReachingVariableDef_Definition
	if f.gatherInState != nil {
		def = f.reachingDef.getDefFromMustDefByID(nameID, f.gatherInState)
	} else {
		def = f.reachingDef.getDefByID(nameID, f.gatherCfgNode)
	}
	if def != nil && !def.node.IsFunction() && !def.node.IsBlock() &&
		!f.reachingDef.DependsOnOuterScopeVars(def) {
		if !f.getCodingConvention().IsExportedWithLocal(rhino.RhinoStringPool_StringFromID(nameID) /* local= */, true) {
			f.candidates = append(f.candidates, FlowSensitiveInlineVariables_Candidate{
				outer:       f,
				nameID:      nameID,
				defMetadata: def,
				use:         n,
				useCfgNode:  f.gatherCfgNode,
			})
		}
	}
}

func (f *FlowSensitiveInlineVariables) walkGatherCfgNode(n, parent rhino.Node, currentNode *rhino.Node) {
	*currentNode = n
	if !parent.IsNil() && ControlFlowGraph_IsEnteringNewCfgNode(n) {
		return
	}
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		if Platform_IsThreadInterrupted() {
			panic(jutil.NewRuntimeExceptionWithCause(jutil.NewInterruptedException()))
		}
		for child := n.GetFirstChild(); !child.IsNil(); {
			next := child.GetNext()
			f.walkGatherCfgNode(child, n, currentNode)
			child = next
		}
	case rhino.Token_FUNCTION:
		fn := ast.MustFunction(n)
		fnName := fn.NameNode().Node()
		args := fn.Params().Node()
		body := fn.Body().Node()
		f.walkGatherCfgNode(fnName, n, currentNode)
		f.walkGatherCfgNode(args, n, currentNode)
		f.walkGatherCfgNode(body, n, currentNode)
	case rhino.Token_CLASS:
		cls := ast.MustClass(n)
		className := cls.NameNode()
		extendsClause := cls.ExtendsNode()
		body := cls.Body().Node()
		f.walkGatherCfgNode(extendsClause, n, currentNode)
		f.walkGatherCfgNode(className, n, currentNode)
		f.walkGatherCfgNode(body, n, currentNode)
	case rhino.Token_CLASS_MEMBERS:
		for child := n.GetFirstChild(); !child.IsNil(); {
			next := child.GetNext()
			if cm, ok := ast.AsComputedMember(child); ok {
				f.walkGatherCfgNode(cm.Key().Node(), child, currentNode)
			}
			child = next
		}
		for child := n.GetFirstChild(); !child.IsNil(); {
			next := child.GetNext()
			if cm, ok := ast.AsComputedMember(child); ok {
				*currentNode = n
				if !ControlFlowGraph_IsEnteringNewCfgNode(child) {
					if val, hasVal := cm.Value(); hasVal {
						f.walkGatherCfgNode(val.Node(), child, currentNode)
					}
				}
			} else {
				switch child.GetToken() {
				case rhino.Token_MEMBER_FIELD_DEF, rhino.Token_BLOCK, rhino.Token_MEMBER_FUNCTION_DEF,
					rhino.Token_MEMBER_VARIABLE_DEF, rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF:
					f.walkGatherCfgNode(child, n, currentNode)
				default:
					panic(jutil.NewIllegalStateExceptionMsg("Invalid class member: " + child.GetToken().String()))
				}
			}
			child = next
		}
	default:
		for child := n.GetFirstChild(); !child.IsNil(); {
			next := child.GetNext()
			if child.GetFirstChild().IsNil() {
				if child.IsName() {
					*currentNode = child
					f.visitCandidateName(child, n)
				} else if child.GetToken() == rhino.Token_SCRIPT && !ControlFlowGraph_IsEnteringNewCfgNode(child) {
					if Platform_IsThreadInterrupted() {
						panic(jutil.NewRuntimeExceptionWithCause(jutil.NewInterruptedException()))
					}
				}
			} else {
				f.walkGatherCfgNode(child, n, currentNode)
			}
			child = next
		}
	}
	if n.IsName() {
		*currentNode = n
		f.visitCandidateName(n, parent)
	}
}

func (f *FlowSensitiveInlineVariables) visitGatherNode(n rhino.Node, currentNode *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_BLOCK, rhino.Token_TRY, rhino.Token_SWITCH_BODY, rhino.Token_EMPTY:
		return
	}
	graphNode := f.cfg.Nodes.Get(n)
	if graphNode == nil {
		// Not a CFG node.
		return
	}
	f.gatherCfgNode = n
	f.gatherInState = f.reachingDef.inStateAt(graphNode.Index)
	f.walkGatherCfgNode(n, rhino.NilNode, currentNode)
	f.gatherInState = nil
}

func (f *FlowSensitiveInlineVariables) traverseGatherCandidates(root rhino.Node) {
	jutil.CheckNotNull(f.compiler)
	if root.IsNil() {
		throwNilCfgNodeTraversal(f.compiler)
		return
	}
	var currentNode rhino.Node
	defer func() {
		if r := recover(); r != nil {
			reportUnexpectedCfgException(f.compiler, root, rhino.NilNode, currentNode, jutil.AsThrowable(r))
		}
	}()
	if Platform_IsThreadInterrupted() {
		panic(jutil.NewRuntimeExceptionWithCause(jutil.NewInterruptedException()))
	}
	currentNode = root
	f.walkGatherCandidates(root, rhino.NilNode, &currentNode)
}

func (f *FlowSensitiveInlineVariables) walkGatherCandidates(n, parent rhino.Node, currentNode *rhino.Node) {
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
			f.walkGatherCandidates(child, n, currentNode)
			child = next
		}
	case rhino.Token_FUNCTION:
		// Only the function itself is a CFG node of the enclosing function:
		// its name is no CFG node, and its params and body are not entered
		// (above).
	case rhino.Token_BLOCK, rhino.Token_SWITCH_BODY, rhino.Token_TRY:
		for child := n.GetFirstChild(); !child.IsNil(); {
			next := child.GetNext()
			f.walkGatherCandidates(child, n, currentNode)
			child = next
		}
	case rhino.Token_IF:
		ifs := ast.MustIf(n)
		els, hasElse := ifs.Else()
		f.walkGatherCandidates(ifs.Then().Node(), n, currentNode)
		if hasElse {
			f.walkGatherCandidates(els.Node(), n, currentNode)
		}
	case rhino.Token_WITH:
		f.walkGatherCandidates(ast.MustWith(n).Body().Node(), n, currentNode)
	case rhino.Token_SWITCH:
		f.walkGatherCandidates(ast.MustSwitch(n).Body().Node(), n, currentNode)
	case rhino.Token_CASE, rhino.Token_DEFAULT_CASE:
		f.walkGatherCandidates(ast.MustSwitchCase(n).Body().Node(), n, currentNode)
	case rhino.Token_CATCH:
		f.walkGatherCandidates(ast.MustCatch(n).Body().Node(), n, currentNode)
	case rhino.Token_LABEL:
		f.walkGatherCandidates(ast.MustLabel(n).Body().Node(), n, currentNode)
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
			f.walkGatherCandidates(child, n, currentNode)
			child = next
		}
	}
	*currentNode = n
	f.visitGatherNode(n, currentNode)
}
