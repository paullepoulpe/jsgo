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

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

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
	compiler AbstractCompiler

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

	astAnalyzer := p.outer.compiler.GetAstAnalyzer()
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

// Java: FlowSensitiveInlineVariables.enterScope(NodeTraversal)
func (f *FlowSensitiveInlineVariables) EnterScope(t *NodeTraversal) {

	if t.InGlobalScope() {
		return // Don't even brother. All global variables are likely escaped.
	}

	if !t.InFunctionBlockScope() {
		return // Only want to do the following if its a function block scope.
	}

	functionScopeRoot := t.GetScopeRoot().GetParent()

	if !f.IsCandidateFunction(functionScopeRoot) {
		return
	}

	scope := t.GetScope()
	if LiveVariablesAnalysis_MAX_VARIABLES_TO_ANALYZE < scope.GetVarCount() {
		return
	}

	scopeCreator := t.GetScopeCreator().(*SyntacticScopeCreator)
	fnScope := scope.GetParent()
	var allVarsInFn *jutil.LinkedHashMap[string, *Var]
	if !fsivHasNestedBlockScopeDecls(functionScopeRoot, t.GetScopeRoot()) {
		f.allVarsMap.Clear()
		for v := range fnScope.GetVarIterable() {
			f.allVarsMap.Put(v.GetName(), v)
		}
		for v := range scope.GetVarIterable() {
			f.allVarsMap.Put(v.GetName(), v)
		}
		allVarsInFn = &f.allVarsMap
	} else {
		allVarsDeclaredInFunction := NodeUtil_GetAllVarsDeclaredInFunction(f.compiler, scopeCreator, fnScope)
		allVarsInFn = allVarsDeclaredInFunction.GetAllVariables()
	}
	if allVarsInFn.IsEmpty() {
		return
	}

	f.escapedSet.Clear()
	escaped := &f.escapedSet
	DataFlowAnalysis_ComputeEscaped(
		fnScope, escaped, f.compiler, scopeCreator, allVarsInFn)
	if escaped.Size() >= allVarsInFn.Size() {
		hasUnescaped := false
		for v := range allVarsInFn.Values() {
			if !escaped.Contains(v) {
				hasUnescaped = true
				break
			}
		}
		if !hasUnescaped {
			return
		}
	}

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
	fsivTraverseGatherCandidates(f.compiler, t.GetScopeRoot(), &f.gatherCandidates)

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

// Java: FlowSensitiveInlineVariables.process(Node, Node)
func (f *FlowSensitiveInlineVariables) Process(externs *rhino.Node, root *rhino.Node) {
	sc := NewSyntacticScopeCreator(f.compiler)
	sc.EnableArenaPool()
	rc := acquireReusableCfg()
	f.reusableCfg = rc
	defer func() {
		if f.reusableCfg != nil {
			releaseReusableCfg(f.reusableCfg)
			f.reusableCfg = nil
		}
		f.escapedSet.Clear()
		sc.ReleaseArena()
	}()
	b := NodeTraversal_Builder{
		compiler:     f.compiler,
		callback:     f,
		scopeCreator: sc,
	}
	b.TraverseRoots(externs, root)
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
			!cb.outer.compiler.GetCodingConvention().IsExportedWithLocal(name /* local= */, true) &&
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
	case rhino.Token_WHILE, rhino.Token_WITH, rhino.Token_SWITCH, rhino.Token_CASE, rhino.Token_CATCH, rhino.Token_LABEL:
		if second := n.GetSecondChild(); second != nil {
			fsivWalkGatherCandidates(second, n, cb, currentNode)
		}
	case rhino.Token_DO, rhino.Token_DEFAULT_CASE:
		if first := n.GetFirstChild(); first != nil {
			fsivWalkGatherCandidates(first, n, cb, currentNode)
		}
	case rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
		if first := n.GetFirstChild(); first != nil {
			fsivWalkGatherCandidates(first, n, cb, currentNode)
			if last := n.GetLastChild(); last != nil && last != first {
				fsivWalkGatherCandidates(last, n, cb, currentNode)
			}
		}
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
	if c.outer.compiler.GetAstAnalyzer().MayHaveSideEffects(c.def.GetLastChild()) {
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
		c.outer.compiler.ReportChangeToEnclosingScope(defParent)
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
	c.outer.compiler.ReportChangeToEnclosingScope(useParent)
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

func fsivHasOutOfScopeName(n *rhino.Node, usageScope *Scope) bool {
	if n.IsName() {
		name := n.GetString()
		if name != "" && !usageScope.HasSlot(name) {
			return true
		}
	}
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		if fsivHasOutOfScopeName(c, usageScope) {
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
	if fsivHasOutOfScopeName(rhs, usageScope) {
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
