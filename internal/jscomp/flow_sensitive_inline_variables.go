// Port of com/google/javascript/jscomp/FlowSensitiveInlineVariables.java @ closure-compiler v20261015

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
	"iter"
	"sync"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type fsivCandidateDef struct {
	varName string
	rhs     *rhino.Node
	hasRead bool
	hasFunc bool
}

type fsivCandFunc struct {
	fn         *rhino.Node
	body       *rhino.Node
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
	// This pass first perform a traversal to gather a list of Candidates that
	// could be inlined using GatherCandidates.
	//
	// The second step involves verifying that each candidate is actually safe
	// to inline with Candidate#canInline(Scope) and finally perform
	// inlining using Candidate#inlineVariable().
	//
	// The reason for the delayed evaluation of the candidates is because we
	// need two separate dataflow result.
	compiler         AbstractCompiler
	astAnalyzer      *AstAnalyzer
	codingConvention CodingConvention

	sideEffectPredicate   FlowSensitiveInlineVariables_SideEffectPredicate
	sideEffectPredicateFn func(*rhino.Node) bool

	// These two pieces of data is persistent in the whole execution of enter
	// scope.
	cfg          ControlFlowGraph[*rhino.Node]
	candidates   *jutil.LinkedHashSet[*FlowSensitiveInlineVariables_Candidate]
	reachingDef  *MustBeReachingVariableDef
	reachingUses *MaybeReachingVariableUse
	escaped      *jutil.ObjectLinkedHashSet[*Var]
	allVarsInFn  *jutil.LinkedHashMap[string, *Var]

	reusableCfg         *ReusableCfg
	reachingDefStorage  MustBeReachingVariableDef
	reachingUsesStorage MaybeReachingVariableUse
	allVarsMap          jutil.LinkedHashMap[string, *Var]
	candidateVarsMap    jutil.LinkedHashMap[string, *Var]
	namesToCheckSet     jutil.LinkedHashSet[string]
	escapedSet          jutil.ObjectLinkedHashSet[*Var]
	candidatesSet       jutil.LinkedHashSet[*FlowSensitiveInlineVariables_Candidate]
	gatherCandidates    FlowSensitiveInlineVariables_GatherCandidates
	gatherCb            FlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback
	candidateSlab       []FlowSensitiveInlineVariables_Candidate
	candidateBuf        [8]FlowSensitiveInlineVariables_Candidate

	varsFinder          lvaAllVarsFinder
	inlinableDefs       []fsivCandidateDef
	inlinableDefsBuf    [16]fsivCandidateDef
	unreadDefCount      int
	candFuncs           []fsivCandFunc
	maxCandDepth        int
	dummyGlobalScope    *Scope
	globalRoot          *rhino.Node
	globalScopeOnce     sync.Once
	realGlobalScope     *Scope
	globalScopeCreator  *SyntacticScopeCreator
	sharedParent        *FlowSensitiveInlineVariables
	scopeCreator        *SyntacticScopeCreator
	curFn               *rhino.Node
	curFnEnclosingScope *Scope
	curTraversal        *NodeTraversal
	parallelMode        bool
	parallelChanges     []*rhino.Node
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

func (f *FlowSensitiveInlineVariables) reportChange(n *rhino.Node) {
	if f.parallelMode {
		if changeRoot := ChangeTracker_GetEnclosingChangeScopeRoot(n); changeRoot != nil {
			f.parallelChanges = append(f.parallelChanges, changeRoot)
		}
		return
	}
	f.compiler.ReportChangeToEnclosingScope(n)
}

func (f *FlowSensitiveInlineVariables) allocCandidate(
	varName string,
	defMetadata *MustBeReachingVariableDef_Definition,
	use *rhino.Node,
	useCfgNode *rhino.Node) *FlowSensitiveInlineVariables_Candidate {
	jutil.CheckArgument(use.IsName())
	if len(f.candidateSlab) == cap(f.candidateSlab) {
		newCap := 16
		if cap(f.candidateSlab) >= 16 {
			newCap = cap(f.candidateSlab) * 2
		}
		f.candidateSlab = make([]FlowSensitiveInlineVariables_Candidate, 0, newCap)
	}
	idx := len(f.candidateSlab)
	f.candidateSlab = f.candidateSlab[:idx+1]
	c := &f.candidateSlab[idx]
	*c = FlowSensitiveInlineVariables_Candidate{
		outer:       f,
		varName:     varName,
		defMetadata: defMetadata,
		use:         use,
		useCfgNode:  useCfgNode,
	}
	return c
}

func (f *FlowSensitiveInlineVariables) getReachingUses() *MaybeReachingVariableUse {
	if f.reachingUses == nil {
		f.candidateVarsMap.Clear()
		for cand := range f.candidates.All() {
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

// FlowSensitiveInlineVariables_SideEffectPredicate is a Predicate<Node>.
type FlowSensitiveInlineVariables_SideEffectPredicate struct {
	outer *FlowSensitiveInlineVariables

	// Check if there are side effects affecting the value of any of these
	// names (but not properties defined on that name)
	namesToCheck *jutil.LinkedHashSet[string]
}

// Java: FlowSensitiveInlineVariables.SideEffectPredicate()
func NewFlowSensitiveInlineVariables_SideEffectPredicate(outer *FlowSensitiveInlineVariables) *FlowSensitiveInlineVariables_SideEffectPredicate {
	return &FlowSensitiveInlineVariables_SideEffectPredicate{outer: outer, namesToCheck: nil}
}

// Java: FlowSensitiveInlineVariables.SideEffectPredicate(Set<String>)
func NewFlowSensitiveInlineVariables_SideEffectPredicateWithNames(outer *FlowSensitiveInlineVariables, names *jutil.LinkedHashSet[string]) *FlowSensitiveInlineVariables_SideEffectPredicate {
	return &FlowSensitiveInlineVariables_SideEffectPredicate{outer: outer, namesToCheck: names}
}

// Java: FlowSensitiveInlineVariables.SideEffectPredicate.apply(Node)
func (p *FlowSensitiveInlineVariables_SideEffectPredicate) Apply(n *rhino.Node) bool {
	// When the node is null it means, we reached the implicit return
	// where the function returns (possibly without an return statement)
	if n == nil {
		return false
	}

	if p.namesToCheck != nil &&
		n.IsName() &&
		p.namesToCheck.Contains(n.GetString()) &&
		NodeUtil_IsLValue(n) {
		// the name is being written to. this is a problem, unless it is part of a top-level assign
		// chain and the write will take place after all CFG node subexpressions are evaluated
		return !FlowSensitiveInlineVariables_IsTopLevelAssignTarget(n)
	}

	astAnalyzer := p.outer.getAstAnalyzer()
	// TODO(user): We only care about calls to functions that
	// passes one of the dependent variable to a non-side-effect free
	// function.
	if (n.IsCall() || n.IsOptChainCall() || n.IsTaggedTemplateLit()) &&
		astAnalyzer.FunctionCallHasSideEffects(n) {
		return true
	}

	if n.IsNew() && astAnalyzer.ConstructorCallHasSideEffects(n) {
		return true
	}

	if n.IsDelProp() {
		return true
	}

	if (n.IsGetProp() || n.IsGetElem()) && NodeUtil_IsLValue(n) {
		return p.namesToCheck == nil || !FlowSensitiveInlineVariables_IsTopLevelAssignTarget(n)
	}

	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		if !ControlFlowGraph_IsEnteringNewCfgNode(c) && p.Apply(c) {
			return true
		}
	}
	return false
}

// FlowSensitiveInlineVariables_IsTopLevelAssignTarget reports whether the
// given node is the target of a (possibly chained) assignment.
//
// Java: FlowSensitiveInlineVariables.isTopLevelAssignTarget(Node)
func FlowSensitiveInlineVariables_IsTopLevelAssignTarget(n *rhino.Node) bool {
	ancestor := n.GetParent()
	for ancestor.IsAssign() {
		ancestor = ancestor.GetParent()
	}
	return ancestor.IsExprResult()
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

func newFlowSensitiveInlineVariablesWorker(parent *FlowSensitiveInlineVariables, compiler AbstractCompiler, astAnalyzer *AstAnalyzer, codingConvention CodingConvention, dummyGlobalScope *Scope, globalRoot *rhino.Node) *FlowSensitiveInlineVariables {
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
func (f *FlowSensitiveInlineVariables) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	return !n.IsScript() || !t.GetInput().IsExtern()
}

func fsivSubtreeHasBlockScopeDecls(n *rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_LET, rhino.Token_CONST, rhino.Token_CLASS, rhino.Token_CATCH:
		return true
	case rhino.Token_FUNCTION:
		return !NodeUtil_IsFunctionExpression(n)
	}
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		if fsivSubtreeHasBlockScopeDecls(c) {
			return true
		}
	}
	return false
}

func fsivHasNestedBlockScopeDecls(fnNode *rhino.Node, fnBody *rhino.Node) bool {
	for p := NodeUtil_GetFunctionParameters(fnNode).GetFirstChild(); p != nil; p = p.GetNext() {
		if !p.IsName() && fsivSubtreeHasBlockScopeDecls(p) {
			return true
		}
	}
	for stmt := fnBody.GetFirstChild(); stmt != nil; stmt = stmt.GetNext() {
		switch stmt.GetToken() {
		case rhino.Token_FUNCTION:
			continue
		case rhino.Token_CLASS, rhino.Token_CATCH:
			return true
		default:
			for c := stmt.GetFirstChild(); c != nil; c = c.GetNext() {
				if fsivSubtreeHasBlockScopeDecls(c) {
					return true
				}
			}
		}
	}
	return false
}

func fsivInspectRhs(n *rhino.Node) (unsafe bool, hasFunc bool) {
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
	default:
	}
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
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

func fsivIsReadName(n, parent *rhino.Node) bool {
	if parent == nil {
		return false
	}
	if (NodeUtil_IsAssignmentOp(parent) && parent.GetFirstChild() == n) ||
		NodeUtil_IsNameDeclaration(parent) ||
		parent.IsInc() ||
		parent.IsDec() ||
		parent.IsParamList() ||
		parent.IsCatch() ||
		NodeUtil_IsLhsByDestructuring(n) {
		return false
	}
	return true
}

func fsivRhsDependsOnOuterOrCatchVars(n, parent *rhino.Node, allVarsInFn *jutil.LinkedHashMap[string, *Var]) bool {
	if parent != nil && ControlFlowGraph_IsEnteringNewCfgNode(n) {
		return false
	}
	if n.IsName() {
		v := allVarsInFn.Get(n.GetString())
		if v == nil || v.GetScope().(*Scope).IsCatchScope() {
			return true
		}
	}
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		if fsivRhsDependsOnOuterOrCatchVars(c, n, allVarsInFn) {
			return true
		}
	}
	return false
}

func (f *FlowSensitiveInlineVariables) maybeAddCandidateDef(varName string, rhs *rhino.Node) {
	if varName == "" {
		return
	}
	unsafe, hasFunc := fsivInspectRhs(rhs)
	if unsafe {
		return
	}
	if f.getCodingConvention().IsExportedWithLocal(varName, true) {
		return
	}
	if f.getAstAnalyzer().MayHaveSideEffects(rhs) {
		return
	}
	f.inlinableDefs = append(f.inlinableDefs, fsivCandidateDef{
		varName: varName,
		rhs:     rhs,
		hasFunc: hasFunc,
	})
}

func (f *FlowSensitiveInlineVariables) scanStatementsForCandidateDefs(n *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_BLOCK, rhino.Token_SWITCH_BODY, rhino.Token_TRY:
		for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
			f.scanStatementsForCandidateDefs(c)
		}
	case rhino.Token_IF:
		for c := n.GetSecondChild(); c != nil; c = c.GetNext() {
			f.scanStatementsForCandidateDefs(c)
		}
	case rhino.Token_WHILE, rhino.Token_WITH, rhino.Token_SWITCH, rhino.Token_CASE, rhino.Token_CATCH, rhino.Token_LABEL:
		if second := n.GetSecondChild(); second != nil {
			f.scanStatementsForCandidateDefs(second)
		}
	case rhino.Token_DO, rhino.Token_DEFAULT_CASE:
		if first := n.GetFirstChild(); first != nil {
			f.scanStatementsForCandidateDefs(first)
		}
	case rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
		if first := n.GetFirstChild(); first != nil {
			f.scanStatementsForCandidateDefs(first)
			if last := n.GetLastChild(); last != nil && last != first {
				f.scanStatementsForCandidateDefs(last)
			}
		}
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		for lhs := n.GetFirstChild(); lhs != nil; lhs = lhs.GetNext() {
			if lhs.IsName() && lhs.HasChildren() {
				f.maybeAddCandidateDef(lhs.GetString(), lhs.GetLastChild())
			}
		}
	case rhino.Token_EXPR_RESULT:
		if c := n.GetFirstChild(); c != nil && c.IsAssign() {
			if lhs := c.GetFirstChild(); lhs.IsName() {
				f.maybeAddCandidateDef(lhs.GetString(), c.GetLastChild())
			}
		}
	case rhino.Token_FUNCTION, rhino.Token_RETURN, rhino.Token_THROW,
		rhino.Token_BREAK, rhino.Token_CONTINUE, rhino.Token_DEBUGGER,
		rhino.Token_EMPTY, rhino.Token_CLASS, rhino.Token_EXPORT, rhino.Token_IMPORT:
		// ControlFlowAnalysis does not enter expressions, functions, or classes for statements.
	default:
		for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
			f.scanStatementsForCandidateDefs(c)
		}
	}
}

func (f *FlowSensitiveInlineVariables) markReadsOutsideLoop(n *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_FUNCTION,
		rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF,
		rhino.Token_DO, rhino.Token_WHILE,
		rhino.Token_CLASS, rhino.Token_IMPORT:
		return
	}
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
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

func (f *FlowSensitiveInlineVariables) collectInlinableDefsWithReads(fnBody *rhino.Node) (ok bool, hasFuncRhs bool) {
	if cap(f.inlinableDefs) == 0 {
		f.inlinableDefs = f.inlinableDefsBuf[:0]
	} else {
		f.inlinableDefs = f.inlinableDefs[:0]
	}
	f.scanStatementsForCandidateDefs(fnBody)
	if len(f.inlinableDefs) == 0 {
		return false, false
	}
	f.unreadDefCount = len(f.inlinableDefs)
	f.markReadsOutsideLoop(fnBody)
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
	fnBody := t.GetScopeRoot()
	if !fnBody.IsBlock() {
		return
	}
	functionScopeRoot := fnBody.GetParent()
	if functionScopeRoot == nil || !functionScopeRoot.IsFunction() {
		return // Only want to do the following if it's a function block scope.
	}

	f.curTraversal = t
	defer func() {
		f.curTraversal = nil
	}()
	scopeCreator := t.GetScopeCreator().(*SyntacticScopeCreator)
	f.processFunction(functionScopeRoot, fnBody, scopeCreator)
}

func (f *FlowSensitiveInlineVariables) processFunction(functionScopeRoot, fnBody *rhino.Node, scopeCreator *SyntacticScopeCreator) {
	if ok, _ := f.collectInlinableDefsWithReads(fnBody); !ok {
		return
	}

	var fnScope, scope *Scope
	if f.dummyGlobalScope != nil {
		fnScope = scopeCreator.CreateScopeOfScope(functionScopeRoot, f.dummyGlobalScope)
		scope = scopeCreator.CreateScopeOfScope(fnBody, fnScope)
	} else if f.curTraversal != nil {
		scope = f.curTraversal.GetScope()
		fnScope = scope.GetParent()
	} else {
		root := functionScopeRoot
		for root.GetParent() != nil {
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
	if !fsivHasNestedBlockScopeDecls(functionScopeRoot, fnBody) {
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
		if allVarsInFn.Get(d.varName) != nil && !fsivRhsDependsOnOuterOrCatchVars(d.rhs, nil, allVarsInFn) {
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
		f.curFn = nil
		f.curFnEnclosingScope = nil
	}()

	// Compute the forward reaching definition.
	rc := f.reusableCfg
	if rc == nil {
		rc = acquireReusableCfg()
		f.reusableCfg = rc
	}
	f.cfg = rc.cfa.ComputeReusableCfg(
		f.compiler,
		functionScopeRoot,
		/* shouldTraverseFunctions= */ false,
		/* edgeAnnotations= */ true,
		&rc.cfg)

	f.reachingDefStorage.Reset(f.cfg, f.compiler, escaped, allVarsInFn)
	f.reachingDef = &f.reachingDefStorage
	f.reachingDef.Analyze()
	f.candidatesSet.Clear()
	f.candidates = &f.candidatesSet
	if cap(f.candidateSlab) == 0 {
		f.candidateSlab = f.candidateBuf[:0]
	} else {
		clear(f.candidateSlab)
		f.candidateSlab = f.candidateSlab[:0]
	}

	// Using the forward reaching definition search to find all the inline
	// candidates
	fsivTraverseGatherCandidates(f.compiler, fnBody, &f.gatherCandidates)

	if f.candidates.IsEmpty() {
		return
	}

	// Lazily compute the backward reaching use only if a candidate passes the
	// syntactic and side-effect checks in Candidate.CanInline.
	f.reachingUses = nil
	f.escaped = escaped
	f.allVarsInFn = allVarsInFn
	for !f.candidates.IsEmpty() {
		var c *FlowSensitiveInlineVariables_Candidate
		for first := range f.candidates.All() {
			c = first
			break
		}
		candidateVar := jutil.CheckNotNull(allVarsInFn.Get(c.varName))
		if c.CanInline(candidateVar.GetScope().(*Scope)) {
			c.InlineVariable()
			f.candidates.Remove(c)

			// If candidate "c" has dependencies, then inlining it may have introduced new dependencies
			// for our other inlining candidates. MustBeReachingVariableDef uses a dependency graph in
			// its analysis. Generating a new dependency graph will need another CFG computation.
			// Ideally we should iterate to a fixed point, but that can be costly. Therefore, we use
			// a conservative heuristic here: For each candidate "other", we back off if its set of
			// dependencies cannot contain all of "c"'s dependencies.
			if !c.defMetadata.depends.IsEmpty() {
				scopeVar := scope.GetVar(c.varName)
				for other := range f.candidates.All() {
					if other.defMetadata.depends.Contains(scopeVar) &&
						!flowSensitiveInlineVariables_containsAll(other.defMetadata.depends, c.defMetadata.depends) {
						f.candidates.Remove(other)
					}
				}
			}
		} else {
			f.candidates.Remove(c)
		}
	}
}

// flowSensitiveInlineVariables_containsAll is Set.containsAll.
//
// Go-only: language-adapter — jutil.ObjectLinkedHashSet has no containsAll.
func flowSensitiveInlineVariables_containsAll(set *jutil.ObjectLinkedHashSet[*Var], elems *jutil.ObjectLinkedHashSet[*Var]) bool {
	for e := range elems.All() {
		if !set.Contains(e) {
			return false
		}
	}
	return true
}

// Java: FlowSensitiveInlineVariables.isCandidateFunction(Node)
func (f *FlowSensitiveInlineVariables) IsCandidateFunction(fn *rhino.Node) bool {
	fnBody := fn.GetLastChild()
	return FlowSensitiveInlineVariables_ContainsCandidateExpressions(fnBody)
}

// Java: FlowSensitiveInlineVariables.containsCandidateExpressions(Node)
func FlowSensitiveInlineVariables_ContainsCandidateExpressions(n *rhino.Node) bool {
	if n.IsFunction() {
		// don't recurse into inner functions or into expressions the can't contain declarations.
		return false
	}

	if NodeUtil_IsNameDeclaration(n) || FlowSensitiveInlineVariables_IsAssignmentToName(n) {
		// if it is a simple assignment
		if n.GetFirstChild().IsName() {
			return true
		}
	}

	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		if FlowSensitiveInlineVariables_ContainsCandidateExpressions(c) {
			return true
		}
	}
	return false
}

// Java: FlowSensitiveInlineVariables.isAssignmentToName(Node)
func FlowSensitiveInlineVariables_IsAssignmentToName(n *rhino.Node) bool {
	if NodeUtil_IsAssignmentOp(n) || n.IsDec() || n.IsInc() {
		// if it is a simple assignment
		return n.GetFirstChild().IsName()
	}
	return false
}

// Java: FlowSensitiveInlineVariables.exitScope(NodeTraversal)
func (f *FlowSensitiveInlineVariables) ExitScope(t *NodeTraversal) {}

func (f *FlowSensitiveInlineVariables) collectCandidateFuncs(n, parent *rhino.Node, depth int) {
	if n == nil {
		return
	}
	if n.IsFunction() {
		fnBody := n.GetLastChild()
		if fnBody != nil && fnBody.IsBlock() {
			if ok, hasFuncRhs := f.collectInlinableDefsWithReads(fnBody); ok {
				if depth > f.maxCandDepth {
					f.maxCandDepth = depth
				}
				f.candFuncs = append(f.candFuncs, fsivCandFunc{
					fn:         n,
					body:       fnBody,
					depth:      depth,
					hasFuncRhs: hasFuncRhs,
				})
				depth++
			}
		}
		fnName := n.GetFirstChild()
		isFunctionDeclaration := parent != nil && NodeUtil_IsFunctionDeclaration(n)
		if isFunctionDeclaration {
			f.collectCandidateFuncs(fnName, n, depth)
		}
		if !isFunctionDeclaration {
			f.collectCandidateFuncs(fnName, n, depth)
		}
		args := fnName.GetNext()
		body := args.GetNext()
		f.collectCandidateFuncs(args, n, depth)
		f.collectCandidateFuncs(body, n, depth)
		return
	}
	if n.IsClass() {
		className := n.GetFirstChild()
		extendsClause := className.GetNext()
		body := extendsClause.GetNext()
		isClassExpression := NodeUtil_IsClassExpression(n)
		f.collectCandidateFuncs(extendsClause, n, depth)
		if !isClassExpression {
			f.collectCandidateFuncs(className, n, depth)
		}
		if isClassExpression {
			f.collectCandidateFuncs(className, n, depth)
		}
		f.collectCandidateFuncs(body, n, depth)
		return
	}
	if n.IsClassMembers() {
		for child := n.GetFirstChild(); child != nil; child = child.GetNext() {
			if child.IsComputedProp() || child.IsComputedFieldDef() {
				f.collectCandidateFuncs(child.GetFirstChild(), child, depth)
			}
		}
		for child := n.GetFirstChild(); child != nil; child = child.GetNext() {
			switch child.GetToken() {
			case rhino.Token_COMPUTED_PROP:
				f.collectCandidateFuncs(child.GetLastChild(), child, depth)
			case rhino.Token_COMPUTED_FIELD_DEF:
				if child.HasTwoChildren() {
					f.collectCandidateFuncs(child.GetLastChild(), child, depth)
				}
			default:
				f.collectCandidateFuncs(child, n, depth)
			}
		}
		return
	}
	for child := n.GetFirstChild(); child != nil; child = child.GetNext() {
		f.collectCandidateFuncs(child, n, depth)
	}
}

// Java: FlowSensitiveInlineVariables.process(Node, Node)
func (f *FlowSensitiveInlineVariables) Process(externs *rhino.Node, root *rhino.Node) {
	// Pre-populate cached AstAnalyzer and CodingConvention on the main goroutine
	// so worker goroutines never race on Compiler.cachedAstAnalyzer initialization.
	astAnalyzer := f.getAstAnalyzer()
	codingConvention := f.getCodingConvention()

	f.candFuncs = f.candFuncs[:0]
	f.maxCandDepth = 0
	if root != nil {
		f.collectCandidateFuncs(root, root.GetParent(), 0)
	}
	if len(f.candFuncs) == 0 {
		return
	}

	globalRoot := root
	if externs != nil && externs.GetParent() != nil {
		globalRoot = externs.GetParent()
	} else {
		for globalRoot.GetParent() != nil {
			globalRoot = globalRoot.GetParent()
		}
	}
	f.globalRoot = globalRoot
	f.dummyGlobalScope = Scope_CreateGlobalScope(globalRoot)
	f.globalScopeOnce = sync.Once{}
	f.realGlobalScope = nil
	f.globalScopeCreator = nil

	sc := NewSyntacticScopeCreator(f.compiler)
	sc.EnableArenaPool()
	rc := acquireReusableCfg()
	f.reusableCfg = rc
	f.scopeCreator = sc
	defer func() {
		if f.reusableCfg != nil {
			releaseReusableCfg(f.reusableCfg)
			f.reusableCfg = nil
		}
		if f.globalScopeCreator != nil {
			f.globalScopeCreator.ReleaseArena()
			f.globalScopeCreator = nil
		}
		f.realGlobalScope = nil
		f.scopeCreator = nil
		f.dummyGlobalScope = nil
		f.globalRoot = nil
		f.escapedSet.Clear()
		sc.ReleaseArena()
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
		changedByFunc := make([][]*rhino.Node, len(parFuncs))
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
				wsc.EnableArenaPool()
				wrc := acquireReusableCfg()
				worker.reusableCfg = wrc
				worker.scopeCreator = wsc
				defer func() {
					releaseReusableCfg(wrc)
					worker.reusableCfg = nil
					worker.scopeCreator = nil
					worker.escapedSet.Clear()
					wsc.ReleaseArena()
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
						nodes := make([]*rhino.Node, len(worker.parallelChanges))
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
func (f *FlowSensitiveInlineVariables) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	// TODO(user): While the helpers do a subtree traversal on the AST, the
	// compiler pass itself only traverse the AST to look for function
	// declarations to perform dataflow analysis on. We could combine
	// the traversal in DataFlowAnalysis's computeEscaped later to save some
	// time.
}

// FlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback is an
// AbstractCfgNodeTraversalCallback.
type FlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback struct {
	ControlFlowGraph_AbstractCfgNodeTraversalCallbackBase
	outer *FlowSensitiveInlineVariables

	cfgNode *rhino.Node
	inState *MustBeReachingVariableDef_MustDef
}

// Java: FlowSensitiveInlineVariables.GatherCandidatesCfgNodeCallback()
func NewFlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback(outer *FlowSensitiveInlineVariables) *FlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback {
	return &FlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback{outer: outer, cfgNode: nil}
}

// Java: FlowSensitiveInlineVariables.GatherCandidatesCfgNodeCallback.setCfgNode(Node)
func (cb *FlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback) SetCfgNode(cfgNode *rhino.Node) {
	cb.cfgNode = cfgNode
	cb.inState = nil
}

// Java: FlowSensitiveInlineVariables.GatherCandidatesCfgNodeCallback.visit(NodeTraversal, Node, Node)
func (cb *FlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if n.IsName() {

		// n.getParent() isn't null. This just the case where n is the root
		// node that gatherCb started at.
		if parent == nil {
			return
		}

		// Make sure that the name node is purely a read.
		if (NodeUtil_IsAssignmentOp(parent) && parent.GetFirstChild() == n) ||
			NodeUtil_IsNameDeclaration(parent) ||
			parent.IsInc() ||
			parent.IsDec() ||
			parent.IsParamList() ||
			parent.IsCatch() ||
			NodeUtil_IsLhsByDestructuring(n) {
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
		// TODO(nicksantos): We need to add some notion of @const outer
		// scope vars. We can inline those just fine.
		if def != nil && !def.node.IsFunction() && !def.node.IsBlock() &&
			// This pass only runs on local scopes.
			!cb.outer.getCodingConvention().IsExportedWithLocal(name /* local= */, true) &&
			!cb.outer.reachingDef.DependsOnOuterScopeVars(def) {
			cb.outer.candidates.Add(cb.outer.allocCandidate(name, def, n, cb.cfgNode))
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
func (g *FlowSensitiveInlineVariables_GatherCandidates) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_BLOCK, rhino.Token_TRY, rhino.Token_SWITCH_BODY, rhino.Token_EMPTY:
		return
	}
	graphNode := g.outer.cfg.GetDiGraphNode(n)
	if graphNode == nil {
		// Not a CFG node.
		return
	}
	cfgNode := n

	g.gatherCb.cfgNode = cfgNode
	state := graphNode.GetAnnotation().(*DataFlowAnalysis_LinearFlowState[*MustBeReachingVariableDef_MustDef])
	g.gatherCb.inState = state.GetIn()
	var currentNode *rhino.Node
	walkCfgNodeBranch(cfgNode, nil, g.gatherCb, &currentNode)
	g.gatherCb.inState = nil
}

func fsivTraverseGatherCandidates(compiler AbstractCompiler, root *rhino.Node, cb *FlowSensitiveInlineVariables_GatherCandidates) {
	jutil.CheckNotNull(compiler)
	if root == nil {
		NodeTraversal_Traverse(compiler, root, cb)
		return
	}
	var currentNode *rhino.Node
	jutil.Try(func() {
		if Platform_IsThreadInterrupted() {
			panic(jutil.NewRuntimeExceptionWithCause(jutil.NewInterruptedException()))
		}
		currentNode = root
		fsivWalkGatherCandidates(root, nil, cb, &currentNode)
	}, func(unexpectedException jutil.Throwable) {
		b := NodeTraversal_Builder{compiler: compiler, callback: cb}
		t := newNodeTraversal(&b)
		t.InitTraversal(root)
		t.currentNode = currentNode
		t.ThrowUnexpectedException(unexpectedException)
	})
}

func fsivWalkGatherCandidates(n, parent *rhino.Node, cb *FlowSensitiveInlineVariables_GatherCandidates, currentNode **rhino.Node) {
	*currentNode = n
	if parent != nil && parent.IsFunction() && n != parent.GetFirstChild() {
		return
	}
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		if Platform_IsThreadInterrupted() {
			panic(jutil.NewRuntimeExceptionWithCause(jutil.NewInterruptedException()))
		}
		for child := n.GetFirstChild(); child != nil; {
			next := child.GetNext()
			fsivWalkGatherCandidates(child, n, cb, currentNode)
			child = next
		}
	case rhino.Token_FUNCTION:
		if parent == nil {
			fnName := n.GetFirstChild()
			args := fnName.GetNext()
			body := args.GetNext()
			fsivWalkGatherCandidates(fnName, n, cb, currentNode)
			fsivWalkGatherCandidates(args, n, cb, currentNode)
			fsivWalkGatherCandidates(body, n, cb, currentNode)
		}
	case rhino.Token_BLOCK, rhino.Token_SWITCH_BODY, rhino.Token_TRY:
		for child := n.GetFirstChild(); child != nil; {
			next := child.GetNext()
			fsivWalkGatherCandidates(child, n, cb, currentNode)
			child = next
		}
	case rhino.Token_IF:
		for child := n.GetSecondChild(); child != nil; {
			next := child.GetNext()
			fsivWalkGatherCandidates(child, n, cb, currentNode)
			child = next
		}
	case rhino.Token_WITH, rhino.Token_SWITCH, rhino.Token_CASE, rhino.Token_CATCH, rhino.Token_LABEL:
		if second := n.GetSecondChild(); second != nil {
			fsivWalkGatherCandidates(second, n, cb, currentNode)
		}
	case rhino.Token_DEFAULT_CASE:
		if first := n.GetFirstChild(); first != nil {
			fsivWalkGatherCandidates(first, n, cb, currentNode)
		}
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
		for child := n.GetFirstChild(); child != nil; {
			next := child.GetNext()
			fsivWalkGatherCandidates(child, n, cb, currentNode)
			child = next
		}
	}
	*currentNode = n
	cb.Visit(nil, n, parent)
}

// FlowSensitiveInlineVariables_Candidate models the connection between a
// definition and a use of that definition.
type FlowSensitiveInlineVariables_Candidate struct {
	outer *FlowSensitiveInlineVariables

	// Name of the variable.
	varName string

	// Nodes related to the definition.
	def         *rhino.Node
	defMetadata *MustBeReachingVariableDef_Definition

	// Nodes related to the use.
	use        *rhino.Node
	useCfgNode *rhino.Node

	// Number of uses of the variable within the current CFG node.
	numUsesWithinCfgNode int
}

// Java: FlowSensitiveInlineVariables.Candidate(String, Definition, Node, Node)
func NewFlowSensitiveInlineVariables_Candidate(
	outer *FlowSensitiveInlineVariables,
	varName string,
	defMetadata *MustBeReachingVariableDef_Definition,
	use *rhino.Node,
	useCfgNode *rhino.Node) *FlowSensitiveInlineVariables_Candidate {
	jutil.CheckArgument(use.IsName())
	return &FlowSensitiveInlineVariables_Candidate{
		outer:       outer,
		varName:     varName,
		defMetadata: defMetadata,
		use:         use,
		useCfgNode:  useCfgNode,
	}
}

// Java: FlowSensitiveInlineVariables.Candidate.getDefCfgNode()
func (c *FlowSensitiveInlineVariables_Candidate) GetDefCfgNode() *rhino.Node {
	return c.defMetadata.node
}

func fsivAlwaysTrueEdgePredicate(DiGraph_DiGraphEdge[*rhino.Node, ControlFlowGraph_Branch]) bool {
	return true
}

// Java: FlowSensitiveInlineVariables.Candidate.canInline(Scope)
func (c *FlowSensitiveInlineVariables_Candidate) CanInline(scope *Scope) bool {
	// Cannot inline a parameter.
	if c.GetDefCfgNode().IsFunction() {
		return false
	}

	// Make sure that the name is not within a loop
	if NodeUtil_IsWithinLoop(c.use) {
		return false
	}

	c.GetDefinition(c.GetDefCfgNode())

	// Definition was not found.
	if c.def == nil {
		return false
	}

	// Check that the assignment isn't used as a R-Value.
	// TODO(user): Certain cases we can still inline.
	if c.def.IsAssign() && !NodeUtil_IsExprAssign(c.def.GetParent()) {
		return false
	}

	if !c.IsRhsSafeToInline(scope) {
		return false
	}

	// TODO(user): Side-effect is OK sometimes. As long as there are no
	// side-effect function down all paths to the use. Once we have all the
	// side-effect analysis tool.
	if c.outer.getAstAnalyzer().MayHaveSideEffects(c.def.GetLastChild()) {
		return false
	}

	// TODO(user): We could inline all the uses if the expression is short.

	// Finally we have to make sure that there are no more than one use
	// in the program and in the CFG node. Even when it is semantically
	// correctly inlining twice increases code size.
	c.GetNumUseInUseCfgNode(c.useCfgNode)
	if c.numUsesWithinCfgNode != 1 {
		return false
	}

	namesToCheck := &c.outer.namesToCheckSet
	namesToCheck.Clear()
	if c.defMetadata.depends != nil {
		for var_ := range c.defMetadata.depends.All() {
			namesToCheck.Add(var_.GetName())
		}
	}

	sideEffectPredicateWithNames := FlowSensitiveInlineVariables_SideEffectPredicate{
		outer:        c.outer,
		namesToCheck: namesToCheck,
	}

	// A subexpression evaluated after the variable has a side effect.
	// Example, for x:
	// x = readProp(b), modifyProp(b); print(x);
	if fsivCheckPostExpressions(c.def, c.GetDefCfgNode(), &sideEffectPredicateWithNames) {
		return false
	}

	// Similar check as the above but this time, all the sub-expressions
	// evaluated before the variable.
	// x = readProp(b); modifyProp(b), print(x);
	if fsivCheckPreExpressions(c.use, c.useCfgNode, &sideEffectPredicateWithNames) {
		return false
	}

	if !c.HasExactlyOne(c.outer.getReachingUses().GetUses(c.varName, c.GetDefCfgNode())) {
		return false
	}

	// We can skip the side effect check along the paths of two nodes if
	// they are just next to each other.
	if NodeUtil_IsStatementBlock(c.GetDefCfgNode().GetParent()) &&
		c.GetDefCfgNode().GetNext() != c.useCfgNode {
		// Similar side effect check as above but this time the side effect is
		// else where along the path.
		// x = readProp(b); while(modifyProp(b)) {}; print(x);
		pathCheck := NewCheckPathsBetweenNodes[*rhino.Node, ControlFlowGraph_Branch](
			c.outer.cfg,
			c.outer.cfg.GetDiGraphNode(c.GetDefCfgNode()),
			c.outer.cfg.GetDiGraphNode(c.useCfgNode),
			c.outer.sideEffectPredicateFn,
			fsivAlwaysTrueEdgePredicate,
			false)
		if pathCheck.SomePathsSatisfyPredicate() {
			return false
		}
	}

	return true
}

// Java: FlowSensitiveInlineVariables.Candidate.hasExactlyOne(Iterable<Node>)
func (c *FlowSensitiveInlineVariables_Candidate) HasExactlyOne(iterable iter.Seq[*rhino.Node]) bool {
	count := 0
	for range iterable {
		count++
		if count > 1 {
			break
		}
	}
	return count == 1
}

// InlineVariable performs the actual transformation.
//
// Java: FlowSensitiveInlineVariables.Candidate.inlineVariable()
func (c *FlowSensitiveInlineVariables_Candidate) InlineVariable() {
	defParent := c.def.GetParent()
	useParent := c.use.GetParent()
	if c.def.IsAssign() {
		rhs := c.def.GetLastChild()
		rhs.Detach()
		// Oh yes! I have grandparent to remove this.
		jutil.CheckState(defParent.IsExprResult())
		for defParent.GetParent().IsLabel() {
			defParent = defParent.GetParent()
		}
		c.outer.reportChange(defParent)
		defParent.Detach()
		c.use.ReplaceWith(rhs)
	} else if NodeUtil_IsNameDeclaration(defParent) {
		rhs := c.def.GetLastChild()
		if defParent.IsConst() {
			// If it is a const var we don't want to remove the rhs of the variable
			rhs.ReplaceWith(rhino.Node_NewStringWithToken(rhino.Token_NAME, "undefined"))
			c.use.ReplaceWith(rhs)
		} else {
			rhs.Detach()
			c.use.ReplaceWith(rhs)
		}
	} else {
		panic(jutil.NewIllegalStateExceptionMsg("No other definitions can be inlined."))
	}
	c.outer.reportChange(useParent)
}

func (c *FlowSensitiveInlineVariables_Candidate) findDefWalk(n, parent *rhino.Node) {
	if parent != nil && ControlFlowGraph_IsEnteringNewCfgNode(n) {
		return
	}
	switch n.GetToken() {
	case rhino.Token_CLASS:
		className := n.GetFirstChild()
		extendsClause := className.GetNext()
		body := extendsClause.GetNext()
		c.findDefWalk(extendsClause, n)
		c.findDefWalk(className, n)
		c.findDefWalk(body, n)
	case rhino.Token_CLASS_MEMBERS:
		for child := n.GetFirstChild(); child != nil; child = child.GetNext() {
			if child.IsComputedProp() || child.IsComputedFieldDef() {
				c.findDefWalk(child.GetFirstChild(), child)
			}
		}
		for child := n.GetFirstChild(); child != nil; child = child.GetNext() {
			switch child.GetToken() {
			case rhino.Token_COMPUTED_PROP, rhino.Token_COMPUTED_FIELD_DEF:
				if !ControlFlowGraph_IsEnteringNewCfgNode(child) && child.HasTwoChildren() {
					c.findDefWalk(child.GetLastChild(), child)
				}
			default:
				c.findDefWalk(child, n)
			}
		}
	default:
		for child := n.GetFirstChild(); child != nil; child = child.GetNext() {
			c.findDefWalk(child, n)
		}
	}
	switch n.GetToken() {
	case rhino.Token_NAME:
		if n.GetString() == c.varName && n.HasChildren() {
			c.def = n
		}
	case rhino.Token_ASSIGN:
		lhs := n.GetFirstChild()
		if lhs.IsName() && lhs.GetString() == c.varName {
			c.def = n
		}
	}
}

// GetDefinition sets the def node.
//
// n is a node that has a corresponding CFG node in the CFG.
//
// Java: FlowSensitiveInlineVariables.Candidate.getDefinition(Node)
func (c *FlowSensitiveInlineVariables_Candidate) GetDefinition(n *rhino.Node) {
	jutil.CheckNotNull(c.outer.compiler)
	if n == nil {
		gatherCb := flowSensitiveInlineVariables_definitionCallback{candidate: c}
		traverseCfgNode(c.outer.compiler, n, &gatherCb)
		return
	}
	c.findDefWalk(n, nil)
}

// flowSensitiveInlineVariables_definitionCallback is the anonymous
// AbstractCfgNodeTraversalCallback of Candidate.getDefinition.
type flowSensitiveInlineVariables_definitionCallback struct {
	ControlFlowGraph_AbstractCfgNodeTraversalCallbackBase
	candidate *FlowSensitiveInlineVariables_Candidate
}

// Java: FlowSensitiveInlineVariables.Candidate.getDefinition(...).new AbstractCfgNodeTraversalCallback() {...}.visit(NodeTraversal, Node, Node)
func (cb *flowSensitiveInlineVariables_definitionCallback) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_NAME:
		if n.GetString() == cb.candidate.varName && n.HasChildren() {
			cb.candidate.def = n
		}
		return
	case rhino.Token_ASSIGN:
		lhs := n.GetFirstChild()
		if lhs.IsName() && lhs.GetString() == cb.candidate.varName {
			cb.candidate.def = n
		}
		return
	default:
	}
}

func (c *FlowSensitiveInlineVariables_Candidate) countUseWalk(n, parent, cfgNode *rhino.Node) {
	if parent != nil && ControlFlowGraph_IsEnteringNewCfgNode(n) {
		return
	}
	if n.GetToken() == rhino.Token_CLASS_MEMBERS {
		for child := n.GetFirstChild(); child != nil; child = child.GetNext() {
			if child.IsComputedProp() || child.IsComputedFieldDef() {
				c.countUseWalk(child.GetFirstChild(), child, cfgNode)
			}
		}
		for child := n.GetFirstChild(); child != nil; child = child.GetNext() {
			switch child.GetToken() {
			case rhino.Token_COMPUTED_PROP, rhino.Token_COMPUTED_FIELD_DEF:
				if !ControlFlowGraph_IsEnteringNewCfgNode(child) && child.HasTwoChildren() {
					c.countUseWalk(child.GetLastChild(), child, cfgNode)
				}
			default:
				c.countUseWalk(child, n, cfgNode)
			}
		}
		return
	}
	for child := n.GetFirstChild(); child != nil; child = child.GetNext() {
		c.countUseWalk(child, n, cfgNode)
	}
	if n.IsName() && n.GetString() == c.varName {
		if parent != nil && parent.IsAssign() && (parent.GetFirstChild() == n) &&
			fsivIsAssignChain(parent, cfgNode) {
			return
		}
		c.numUsesWithinCfgNode++
	}
}

func fsivIsAssignChain(child *rhino.Node, ancestor *rhino.Node) bool {
	for n := child; n != ancestor; n = n.GetParent() {
		if !n.IsAssign() {
			return false
		}
	}
	return true
}

// GetNumUseInUseCfgNode computes the number of uses of the variable varName
// and store it in numUseWithinUseCfgNode.
//
// Java: FlowSensitiveInlineVariables.Candidate.getNumUseInUseCfgNode(Node)
func (c *FlowSensitiveInlineVariables_Candidate) GetNumUseInUseCfgNode(cfgNode *rhino.Node) {
	c.numUsesWithinCfgNode = 0
	jutil.CheckNotNull(c.outer.compiler)
	if cfgNode == nil {
		gatherCb := flowSensitiveInlineVariables_numUseCallback{candidate: c, cfgNode: cfgNode}
		traverseCfgNode(c.outer.compiler, cfgNode, &gatherCb)
		return
	}
	c.countUseWalk(cfgNode, nil, cfgNode)
}

// flowSensitiveInlineVariables_numUseCallback is the anonymous
// AbstractCfgNodeTraversalCallback of Candidate.getNumUseInUseCfgNode.
type flowSensitiveInlineVariables_numUseCallback struct {
	ControlFlowGraph_AbstractCfgNodeTraversalCallbackBase
	candidate *FlowSensitiveInlineVariables_Candidate
	cfgNode   *rhino.Node
}

// Java: FlowSensitiveInlineVariables.Candidate.getNumUseInUseCfgNode(...).new AbstractCfgNodeTraversalCallback() {...}.visit(NodeTraversal, Node, Node)
func (cb *flowSensitiveInlineVariables_numUseCallback) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if n.IsName() && n.GetString() == cb.candidate.varName {
		// We make a special exception when the entire cfgNode is a chain
		// of assignments, since in that case the assignment statements
		// will happen after the inlining of the right hand side.
		// TODO(lharker): We can probably remove the isAssignChain check, and instead use
		// the SideEffectPredicate to look for dangerous assignments in the same CFG node
		if parent.IsAssign() && (parent.GetFirstChild() == n) &&
			cb.IsAssignChain(parent, cb.cfgNode) {
			// Don't count lhs of top-level assignment chain
			return
		} else {
			cb.candidate.numUsesWithinCfgNode++
		}
	}
}

// Java: FlowSensitiveInlineVariables.Candidate.getNumUseInUseCfgNode(...).new AbstractCfgNodeTraversalCallback() {...}.isAssignChain(Node, Node)
func (cb *flowSensitiveInlineVariables_numUseCallback) IsAssignChain(child *rhino.Node, ancestor *rhino.Node) bool {
	for n := child; n != ancestor; n = n.GetParent() {
		if !n.IsAssign() {
			return false
		}
	}
	return true
}

func fsivHasUnsafeRhsNode(n *rhino.Node) bool {
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
		return true
	case rhino.Token_FUNCTION:
		return false
	default:
	}
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		if fsivHasUnsafeRhsNode(c) {
			return true
		}
	}
	return false
}

func (f *FlowSensitiveInlineVariables) getRealGlobalScope(root *rhino.Node) *Scope {
	if f.sharedParent != nil {
		return f.sharedParent.getRealGlobalScope(root)
	}
	if f.globalRoot != nil && root == f.globalRoot {
		f.globalScopeOnce.Do(func() {
			gsc := NewSyntacticScopeCreator(f.compiler)
			gsc.EnableArenaPool()
			f.globalScopeCreator = gsc
			f.realGlobalScope = gsc.CreateScopeOfScope(f.globalRoot, nil)
		})
		return f.realGlobalScope
	}
	return f.scopeCreator.CreateScopeOfScope(root, nil)
}

func (f *FlowSensitiveInlineVariables) buildEnclosingScope() *Scope {
	if f.curFn == nil || f.scopeCreator == nil {
		return nil
	}
	var rootsBuf [16]*rhino.Node
	roots := rootsBuf[:0]
	prev := f.curFn
	for p := prev.GetParent(); p != nil; prev, p = p, p.GetParent() {
		if p.IsClass() && prev == p.GetSecondChild() {
			continue
		}
		if p.IsComputedFieldDef() && prev == p.GetFirstChild() {
			continue
		}
		if (p.IsRoot() && p.GetParent() == nil) || p.GetParent() == nil || p.IsFunction() || NodeUtil_CreatesBlockScope(p) {
			if p.GetParent() != nil &&
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
		if parentScope == nil && r.GetParent() == nil {
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

func (c *FlowSensitiveInlineVariables_Candidate) hasOutOfScopeName(n *rhino.Node, usageScope *Scope) bool {
	if n.IsName() {
		name := n.GetString()
		if name != "" && !usageScope.HasSlot(name) {
			if c.outer == nil || !c.outer.hasEnclosingSlot(name) {
				return true
			}
		}
	}
	for child := n.GetFirstChild(); child != nil; child = child.GetNext() {
		if c.hasOutOfScopeName(child, usageScope) {
			return true
		}
	}
	return false
}

// IsRhsSafeToInline checks if the definition we're considering inline has
// anything that makes inlining unsafe (that hasn't already been caught).
//
// usageScope is the scope we will inline the variable into.
//
// Java: FlowSensitiveInlineVariables.Candidate.isRhsSafeToInline(Scope)
func (c *FlowSensitiveInlineVariables_Candidate) IsRhsSafeToInline(usageScope *Scope) bool {
	rhs := c.def.GetLastChild()
	if fsivHasUnsafeRhsNode(rhs) {
		return false
	}
	if c.hasOutOfScopeName(rhs, usageScope) {
		return false
	}
	return true
}

func fsivCheckPostExpressions(
	n *rhino.Node, expressionRoot *rhino.Node, predicate *FlowSensitiveInlineVariables_SideEffectPredicate) bool {
	for p := n; p != expressionRoot; p = p.GetParent() {
		for cur := p.GetNext(); cur != nil; cur = cur.GetNext() {
			if predicate.Apply(cur) {
				return true
			}
		}
	}
	return false
}

func fsivCheckPreExpressions(
	n *rhino.Node, expressionRoot *rhino.Node, predicate *FlowSensitiveInlineVariables_SideEffectPredicate) bool {
	for p := n; p != expressionRoot; p = p.GetParent() {
		oldestSibling := p.GetParent().GetFirstChild()
		if oldestSibling.IsDestructuringPattern() {
			if p.IsDestructuringPattern() {
				if p.GetNext() != nil && predicate.Apply(p.GetNext()) {
					return true
				}
			}
			continue
		}
		for cur := oldestSibling; cur != p; cur = cur.GetNext() {
			if predicate.Apply(cur) {
				return true
			}
		}
	}
	return false
}

// FlowSensitiveInlineVariables_CheckPostExpressions returns true, given an
// expression by its root and sub-expression n, if the predicate is true for
// some expression evaluated after n.
//
// NOTE: this doesn't correctly check destructuring patterns, because their
// order of evaluation is different from AST traversal order, but currently
// this is ok because FlowSensitiveInlineVariables never inlines variable
// assignments inside destructuring.
//
// Example:
//
// NotChecked(), NotChecked(), n, Checked(), Checked();
//
// Java: FlowSensitiveInlineVariables.checkPostExpressions(Node, Node, Predicate<Node>)
func FlowSensitiveInlineVariables_CheckPostExpressions(
	n *rhino.Node, expressionRoot *rhino.Node, predicate func(*rhino.Node) bool) bool {
	for p := n; p != expressionRoot; p = p.GetParent() {
		for cur := p.GetNext(); cur != nil; cur = cur.GetNext() {
			if predicate(cur) {
				return true
			}
		}
	}
	return false
}

// FlowSensitiveInlineVariables_CheckPreExpressions returns true, given an
// expression by its root and sub-expression n, if the predicate is true for
// some expression evaluated before n.
//
// In most cases evaluation order follows left-to-right AST order.
// Destructuring pattern evaluation is an exception.
//
// Example:
//
// Checked(), Checked(), n, NotChecked(), NotChecked();
//
// Java: FlowSensitiveInlineVariables.checkPreExpressions(Node, Node, Predicate<Node>)
func FlowSensitiveInlineVariables_CheckPreExpressions(
	n *rhino.Node, expressionRoot *rhino.Node, predicate func(*rhino.Node) bool) bool {
	for p := n; p != expressionRoot; p = p.GetParent() {
		oldestSibling := p.GetParent().GetFirstChild()
		// Evaluate a destructuring assignment right-to-left.
		if oldestSibling.IsDestructuringPattern() {
			if p.IsDestructuringPattern() {
				if p.GetNext() != nil && predicate(p.GetNext()) {
					return true
				}
			}
			continue
		}
		for cur := oldestSibling; cur != p; cur = cur.GetNext() {
			if predicate(cur) {
				return true
			}
		}
	}
	return false
}
