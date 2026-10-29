// Port of com/google/javascript/jscomp/CoalesceVariableNames.java @ closure-compiler v20261021

/*
 * Copyright 2008 The Closure Compiler Authors.
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
	"math/bits"
	"slices"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/graph"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type coalesceScopeState struct {
	varsFinder    lvaAllVarsFinder
	liveness      LiveVariablesAnalysis
	orderedVars   []*Var
	varColor      [LiveVariablesAnalysis_MAX_VARIABLES_TO_ANALYZE]int16
	superNode     [LiveVariablesAnalysis_MAX_VARIABLES_TO_ANALYZE]*Var
	degrees       [LiveVariablesAnalysis_MAX_VARIABLES_TO_ANALYZE]int
	worklistBuf   [LiveVariablesAnalysis_MAX_VARIABLES_TO_ANALYZE]uint8
	hasCoalescing bool
}

// CoalesceVariableNames reuses variable names if possible.
//
// For example, from
//
//	var x = 1; print(x); var y = 2; print(y);
//
// to
//
//	var x = 1; print(x); x = 2; print(x)
//
// The benefits are slightly shorter code because of the removed var
// declaration, less unique variables in hope for better renaming, and
// finally better gzip compression.
//
// The pass operates similar to a typical register allocator found in an
// optimizing compiler by first computing live ranges with
// LiveVariablesAnalysis and a variable interference graph. Then it uses
// greedy graph coloring to determine which two variables can be merged
// together safely.
type CoalesceVariableNames struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	cfgScopeStack

	compiler       AbstractCompiler
	scopeCreator   *MemoizedScopeCreator
	activeStates   []*coalesceScopeState
	activeStateBuf [8]*coalesceScopeState
	usePseudoNames bool
	astFactory     *AstFactory
	liveness       *LiveVariablesAnalysis

	// A stack of shouldOptimizeScope results (top is the last element).
	shouldOptimizeScopeStack    []bool
	shouldOptimizeScopeStackBuf [16]bool

	scopeStates    []*coalesceScopeState
	scopeStatesBuf [4]*coalesceScopeState

	interferenceBitsBuf [LiveVariablesAnalysis_MAX_VARIABLES_TO_ANALYZE][2]uint64
	liveRangeChecker    CoalesceVariableNames_LiveRangeChecker
	assignToBuf         [16]int
	readFromBuf         [16]int

	deferChanges     bool
	changedScopes    []rhino.Node
	changedScopesBuf [8]rhino.Node
}

// NewCoalesceVariableNames: usePseudoNames is for debug purposes, when
// merging variable foo and bar to foo, rename both variable to foo_bar.
//
// Java: CoalesceVariableNames(AbstractCompiler, boolean)
func NewCoalesceVariableNames(compiler AbstractCompiler, usePseudoNames bool) *CoalesceVariableNames {
	// The code is normalized at this point in the compilation process. This allows us to use the
	// fact that all variables have been given unique names. We can hoist coalesced variables to
	// VARS because we know that shadowing can't occur.
	jutil.CheckState(compiler.GetLifeCycleStage().IsNormalized())

	c := &CoalesceVariableNames{
		compiler:       compiler,
		usePseudoNames: usePseudoNames,
		astFactory:     compiler.CreateAstFactory(),
		scopeCreator:   NewMemoizedScopeCreator(NewSyntacticScopeCreator(compiler)),
	}
	c.cfgs = c.cfgsBuf[:0]
	c.shouldOptimizeScopeStack = c.shouldOptimizeScopeStackBuf[:0]
	c.activeStates = c.activeStateBuf[:0]
	c.scopeStates = c.scopeStatesBuf[:0]
	c.changedScopes = c.changedScopesBuf[:0]
	return c
}

func cvnHasLocalDecl(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST, rhino.Token_CATCH:
		return true
	default:
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if cvnHasLocalDecl(c) {
			return true
		}
	}
	return false
}

func cvnCollectTopFunctions(n rhino.Node, out *[]rhino.Node) {
	if n.IsFunction() {
		if cvnHasLocalDecl(n) {
			*out = append(*out, n)
		}
		return
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		cvnCollectTopFunctions(c, out)
	}
}

func (c *CoalesceVariableNames) processTopFunction(topFn rhino.Node, dummyGlobal *Scope) {
	fnScope := c.scopeCreator.CreateScope(topFn, dummyGlobal)
	b := NodeTraversal_Builder{
		compiler:     c.compiler,
		callback:     c,
		scopeCreator: c.scopeCreator,
	}
	b.TraverseAtScope(fnScope)
}

func (c *CoalesceVariableNames) reportChange(n rhino.Node) {
	if c.deferChanges {
		c.changedScopes = append(c.changedScopes, ChangeTracker_GetEnclosingChangeScopeRoot(n.GetParent()))
		return
	}
	c.compiler.ReportChangeToEnclosingScope(n)
}

// Java: CoalesceVariableNames.process(Node, Node)
func (c *CoalesceVariableNames) Process(externs rhino.Node, root rhino.Node) {
	rhino.CheckNotNull(externs)
	rhino.CheckNotNull(root)

	if c.usePseudoNames {
		b := NodeTraversal_Builder{
			compiler:     c.compiler,
			callback:     c,
			scopeCreator: c.scopeCreator,
		}
		b.Traverse(root)
		c.compiler.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_RAW)
		return
	}

	var topFns []rhino.Node
	cvnCollectTopFunctions(root, &topFns)
	if len(topFns) == 0 {
		c.compiler.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_RAW)
		return
	}

	dummyGlobal := Scope_CreateGlobalScope(root)
	numThreads := 1
	if opts := c.compiler.GetOptions(); opts != nil && opts.GetNumParallelThreads() > 1 {
		numThreads = opts.GetNumParallelThreads()
	}

	if numThreads <= 1 || len(topFns) <= 1 {
		c.pushCfgRoot(root)
		defer c.resetCfgStack()
		for _, topFn := range topFns {
			c.processTopFunction(topFn, dummyGlobal)
		}
		c.compiler.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_RAW)
		return
	}

	workers := min(numThreads, len(topFns))
	changesByFn := make([][]rhino.Node, len(topFns))
	failures := make([]jutil.Throwable, len(topFns))
	var next atomic.Int64
	var wg sync.WaitGroup
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			wc := NewCoalesceVariableNames(c.compiler, false)
			wc.deferChanges = true
			wc.pushCfgRoot(root)
			defer wc.resetCfgStack()
			curIdx := -1
			defer func() {
				if r := recover(); r != nil && curIdx >= 0 {
					failures[curIdx] = jutil.AsThrowable(r)
				}
			}()
			for {
				idx := int(next.Add(1) - 1)
				if idx >= len(topFns) {
					return
				}
				curIdx = idx
				wc.changedScopes = wc.changedScopes[:0]
				wc.processTopFunction(topFns[idx], dummyGlobal)
				if len(wc.changedScopes) > 0 {
					changesByFn[idx] = slices.Clone(wc.changedScopes)
				}
				curIdx = -1
			}
		}()
	}
	wg.Wait()
	for _, f := range failures {
		if f != nil {
			panic(f)
		}
	}
	for _, scopes := range changesByFn {
		for _, scope := range scopes {
			c.compiler.ReportChangeToChangeScope(scope)
		}
	}
	c.compiler.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_RAW)
}

// CoalesceVariableNames_ShouldOptimizeScope returns populated
// AllVarsDeclaredInFunction object iff shouldOptimizeScope is true.
//
// Java: CoalesceVariableNames.shouldOptimizeScope(NodeTraversal)
func CoalesceVariableNames_ShouldOptimizeScope(t *NodeTraversal) *NodeUtil_AllVarsDeclaredInFunction {
	if t.GetScopeRoot().IsFunction() {
		allVarsDeclaredInFunction :=
			NodeUtil_GetAllVarsDeclaredInFunction(t.GetCompiler(), t.GetScopeCreator(), t.GetScope())
		if LiveVariablesAnalysis_MAX_VARIABLES_TO_ANALYZE >
			len(allVarsDeclaredInFunction.GetAllVariablesInOrder()) {
			return allVarsDeclaredInFunction
		}
	}
	return nil
}

func (c *CoalesceVariableNames) EnterScope(t *NodeTraversal) {
	c.pushScopeRoot(t.GetScopeRoot())
	c.EnterScopeWithCfg(t)
}

func (c *CoalesceVariableNames) ExitScope(t *NodeTraversal) {
	c.ExitScopeWithCfg(t)
	c.popScopeRoot(t.GetScopeRoot())
}

// Java: CoalesceVariableNames.enterScopeWithCfg(NodeTraversal)
func (c *CoalesceVariableNames) EnterScopeWithCfg(t *NodeTraversal) {
	if !t.GetScopeRoot().IsFunction() {
		c.shouldOptimizeScopeStack = append(c.shouldOptimizeScopeStack, false)
		return
	}
	fnScope := t.GetScope()
	if fnScope.GetVarCount() >= LiveVariablesAnalysis_MAX_VARIABLES_TO_ANALYZE {
		c.shouldOptimizeScopeStack = append(c.shouldOptimizeScopeStack, false)
		return
	}
	depth := len(c.activeStates)
	for len(c.scopeStates) <= depth {
		c.scopeStates = append(c.scopeStates, &coalesceScopeState{})
	}
	state := c.scopeStates[depth]
	allVarsDeclaredInFunction := state.varsFinder.collectWithLimit(
		t.GetCompiler(), t.GetScopeCreator(), fnScope, LiveVariablesAnalysis_MAX_VARIABLES_TO_ANALYZE)
	if allVarsDeclaredInFunction == nil {
		c.shouldOptimizeScopeStack = append(c.shouldOptimizeScopeStack, false)
		return
	}
	orderedVars := allVarsDeclaredInFunction.GetAllVariablesInOrder()
	c.shouldOptimizeScopeStack = append(c.shouldOptimizeScopeStack, true)

	candidateCount := 0
	nonParamCount := 0
	for _, v := range orderedVars {
		p := v.GetParentNode()
		if p.IsFunction() || p.IsClass() {
			continue
		}
		candidateCount++
		if !v.IsParam() {
			nonParamCount++
		}
	}
	if candidateCount < 2 || nonParamCount == 0 {
		c.liveness = nil
		c.activeStates = append(c.activeStates, nil)
		return
	}

	jutil.CheckStateMsg1(fnScope.IsFunctionScope(), "%s", fnScope)

	// Live variables analysis is based off of the control flow graph.
	cfg := c.GetControlFlowGraph(c.compiler)

	state.liveness.Reset(cfg, fnScope, nil, c.compiler, c.scopeCreator, allVarsDeclaredInFunction)
	c.liveness = &state.liveness

	if parser.FeatureSet_ES3.Contains(c.compiler.GetOptions().GetOutputFeatureSet()) {
		// If the function has exactly 2 params, mark them as escaped. This is a work-around for a
		// bug in IE 8 and below, where it throws an exception if you write to the parameters of the
		// callback in a sort(). See http://blickly.github.io/closure-compiler-issues/#58 and
		// https://www.zachleat.com/web/array-sort/
		if ast.MustFunction(fnScope.GetRootNode()).Params().NumParams() == 2 {
			c.liveness.MarkAllParametersEscaped()
		}
	}

	unescapedCandidates := 0
	unescapedNonParams := 0
	for vIdx, v := range orderedVars {
		if c.liveness.isEscapedIndex(vIdx) {
			continue
		}
		p := v.GetParentNode()
		if p.IsFunction() || p.IsClass() || c.IsInMultipleLvalueDecl(v) {
			continue
		}
		unescapedCandidates++
		if !v.IsParam() {
			unescapedNonParams++
		}
	}
	if unescapedCandidates < 2 || unescapedNonParams == 0 {
		c.liveness = nil
		c.activeStates = append(c.activeStates, nil)
		return
	}

	c.liveness.Analyze()
	c.computeAndColorScope(cfg, state)
	c.activeStates = append(c.activeStates, state)
}

// Java: CoalesceVariableNames.exitScopeWithCfg(NodeTraversal)
func (c *CoalesceVariableNames) ExitScopeWithCfg(t *NodeTraversal) {
	if len(c.shouldOptimizeScopeStack) == 0 {
		panic(jutil.NewNoSuchElementException()) // ArrayDeque.pop()
	}
	shouldOptimize := c.shouldOptimizeScopeStack[len(c.shouldOptimizeScopeStack)-1]
	c.shouldOptimizeScopeStack = c.shouldOptimizeScopeStack[:len(c.shouldOptimizeScopeStack)-1]
	if !shouldOptimize {
		return
	}
	if len(c.activeStates) == 0 {
		panic(jutil.NewNoSuchElementException()) // ArrayDeque.pop()
	}
	last := len(c.activeStates) - 1
	if st := c.activeStates[last]; st != nil {
		st.orderedVars = nil
		clear(st.superNode[:])
	}
	c.activeStates[last] = nil
	c.activeStates = c.activeStates[:last]
	if len(c.activeStates) == 0 || c.activeStates[len(c.activeStates)-1] == nil {
		c.liveness = nil
	} else {
		c.liveness = &c.activeStates[len(c.activeStates)-1].liveness
	}
}

var coalesceVariableNamesVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_NAME)
	return s
}()

var (
	_ nodeTraversal_visitTokens    = (*CoalesceVariableNames)(nil)
	_ nodeTraversal_onlyNameLeaves = (*CoalesceVariableNames)(nil)
)

func (*CoalesceVariableNames) visitTokens() *ast.TokenSet { return &coalesceVariableNamesVisitTokens }
func (*CoalesceVariableNames) onlyNameLeaves()            {}

// Java: CoalesceVariableNames.visit(NodeTraversal, Node, Node)
func (c *CoalesceVariableNames) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	nameNode, isName := ast.AsName(n)
	if len(c.activeStates) == 0 || !isName || parent.IsFunction() {
		// Don't rename named functions.
		return
	}

	state := c.activeStates[len(c.activeStates)-1]
	if state == nil || !state.hasCoalescing {
		return
	}
	vIdx, ok := c.liveness.lookupVarIndexByID(n.GetStringID())
	if !ok {
		return
	}
	color := state.varColor[vIdx]
	if color < 0 {
		// Not a coalescing candidate (escaped, function/class decl, or multi-lvalue decl).
		return
	}
	var_ := state.orderedVars[vIdx]

	if !c.usePseudoNames {
		coalescedVar := state.superNode[color]
		if coalescedVar == nil {
			state.superNode[color] = var_
			return
		}
		if var_.Equals(coalescedVar) {
			// The coalesced name is itself, nothing to do.
			return
		}

		// Rename.
		nameNode.SetIdent(coalescedVar.GetName())
		c.reportChange(n)
		c.UpdateDeclarationsPostCoalescing(n, coalescedVar, var_, parent)
	} else {
		coalescedVar := state.superNode[color]
		if coalescedVar == nil {
			coalescedVar = var_
			state.superNode[color] = var_
		}
		// This code block is slow but since usePseudoName is for debugging,
		// we should not sacrifice performance for non-debugging compilation to
		// make this fast.
		var allMergedNames []string
		for i, iVar := range state.orderedVars {
			if state.varColor[i] == color {
				allMergedNames = append(allMergedNames, iVar.GetName())
			}
		}
		slices.Sort(allMergedNames)
		allMergedNames = slices.Compact(allMergedNames)

		// Keep its original name.
		if len(allMergedNames) == 1 {
			return
		}

		pseudoName := strings.Join(allMergedNames, "_")

		for t.GetScope().HasSlot(pseudoName) {
			pseudoName += "$"
		}

		// Rename.
		nameNode.SetIdent(pseudoName)
		c.compiler.ReportChangeToEnclosingScope(n)

		if var_.Equals(coalescedVar) {
			return
		}
		c.UpdateDeclarationsPostCoalescing(n, coalescedVar, var_, parent)
	}
}

// UpdateDeclarationsPostCoalescing updates declarations of the given name
// node and the coalesced variable. Firstly, it converts the coalesced
// variable's declaration into a `var` if it is a `const` or a `let`.
// Secondly, it removes the declaration of the given name node if it is under
// a parent which is a name declaration or if the given name is a
// destructuring LHS name under a destructuring LHS declaration.
//
// Java: CoalesceVariableNames.updateDeclarationsPostCoalescing(Node, Var, Var, Node)
func (c *CoalesceVariableNames) UpdateDeclarationsPostCoalescing(
	n rhino.Node, coalescedVar *Var, originalVar *Var, parent rhino.Node) {
	jutil.CheckStateMsg(n.IsName(), "trying to update the declaration of a non-name node")
	if NodeUtil_IsNameDeclaration(parent) ||
		(!NodeUtil_GetEnclosingType(n, rhino.Token_DESTRUCTURING_LHS).IsNil() &&
			NodeUtil_IsLhsByDestructuring(n) &&
			// Consider this to be a destructuring LHS declaration only if we don't hit an
			// assignment like `[a,b] = [3,4]` first.
			!c.IsNameInsideDestructuringAssignment(n)) {
		// convert the coalesced variable's declaration into a `var` if it is a `const` or a `let`
		c.MakeDeclarationVar(coalescedVar)
		// remove the declaration of the given name node as it has been coalesced with coalescedVar
		c.RemoveVarDeclaration(n, originalVar)
	}
}

// IsNameInsideDestructuringAssignment: is the given name node inside the LHS
// of a destructuring pattern assignment like `[a,b,c] = [1,2,3]`.
//
// Importantly, for this pass, it means that the name is not inside a
// declaration like `const [a,b,c] = [1,2,3]` and must not be undeclared.
//
// Java: CoalesceVariableNames.isNameInsideDestructuringAssignment(Node)
func (c *CoalesceVariableNames) IsNameInsideDestructuringAssignment(n rhino.Node) bool {
	enclosingTarget := NodeUtil_GetRootTarget(n)
	_, isPat := ast.AsPattern(enclosingTarget)
	return isPat && ast.RoleOf(enclosingTarget) == ast.RoleAssignTarget
}

// computeAndColorScope builds the 128-bit interference matrix for the
// candidate variables in orderedVariables (guaranteed < MAX_VARIABLES_TO_ANALYZE = 100)
// and greedily colors them with exact LinkedUndirectedGraph degree + tie-break order.
func (c *CoalesceVariableNames) computeAndColorScope(
	cfg *ControlFlowGraph[rhino.Node],
	state *coalesceScopeState,
) {
	orderedVariables := c.liveness.GetAllVariablesInOrder()
	numVars := len(orderedVariables)
	for i := 0; i < numVars; i++ {
		state.varColor[i] = -1
	}

	var candW0, candW1 uint64
	var paramW0, paramW1 uint64
	worklist := state.worklistBuf[:0]
	for vIndex, v := range orderedVariables {
		if c.liveness.isEscapedIndex(vIndex) {
			continue
		}
		// NOTE(user): In theory, we CAN coalesce function names just like any variables,
		// and class declarations if rewritten as expressions, but Closure skips both due to
		// block-scoping and gzip size regressions.
		p := v.GetParentNode()
		if p.IsFunction() || p.IsClass() {
			continue
		}
		// Skip lets, consts, and vars that have multiple variables declared in them
		// (see testCapture and https://github.com/google/closure-compiler/issues/3164).
		if c.IsInMultipleLvalueDecl(v) {
			continue
		}
		if vIndex < 64 {
			bit := uint64(1) << uint(vIndex)
			candW0 |= bit
			if v.IsParam() {
				paramW0 |= bit
			}
		} else {
			bit := uint64(1) << uint(vIndex-64)
			candW1 |= bit
			if v.IsParam() {
				paramW1 |= bit
			}
		}
		worklist = append(worklist, uint8(vIndex))
	}
	numCandidates := len(worklist)

	interferenceBits := c.interferenceBitsBuf[:numVars]
	clear(interferenceBits)

	lrc := &c.liveRangeChecker
	lrc.orderedVariables = orderedVariables
	lrc.liveness = c.liveness
	if lrc.isAssignToList == nil {
		lrc.isAssignToList = c.assignToBuf[:0]
	}
	if lrc.isReadFromList == nil {
		lrc.isReadFromList = c.readFromBuf[:0]
	}

	implicitReturn := cfg.implicitReturn
	for cfgNode := range cfg.Nodes.Values() {
		if cfgNode != implicitReturn {
			c.processCfgNodeBits128(cfgNode, interferenceBits)
		}
	}
	lrc.root = rhino.NilNode
	lrc.state = nil
	lrc.orderedVariables = nil
	lrc.liveness = nil

	// Finalize candidate adjacency masks (adding param-param interference) and
	// compute exact LinkedUndirectedGraph node degrees (where self-edges count twice).
	for _, rawIdx := range worklist {
		vIdx := int(rawIdx)
		w0 := interferenceBits[vIdx][0]
		w1 := interferenceBits[vIdx][1]
		if orderedVariables[vIdx].IsParam() {
			w0 |= paramW0
			w1 |= paramW1
		}
		w0 &= candW0
		w1 &= candW1
		interferenceBits[vIdx][0] = w0
		interferenceBits[vIdx][1] = w1

		deg := bits.OnesCount64(w0) + bits.OnesCount64(w1)
		if vIdx < 64 {
			if w0&(uint64(1)<<uint(vIdx)) != 0 {
				deg++
			}
		} else if w1&(uint64(1)<<uint(vIdx-64)) != 0 {
			deg++
		}
		state.degrees[vIdx] = deg
	}

	// Sort candidate indices by descending degree; ties retain ascending vIdx order.
	degrees := &state.degrees
	slices.SortStableFunc(worklist, func(left, right uint8) int {
		return degrees[right] - degrees[left]
	})

	colorCount := 0
	if len(worklist) == 0 {
		colorCount = 1
	} else {
		for len(worklist) > 0 {
			var forb0, forb1 uint64
			kept := worklist[:0]
			for _, rawIdx := range worklist {
				vIdx := int(rawIdx)
				var independent bool
				if vIdx < 64 {
					independent = forb0&(uint64(1)<<uint(vIdx)) == 0
				} else {
					independent = forb1&(uint64(1)<<uint(vIdx-64)) == 0
				}
				if independent {
					forb0 |= interferenceBits[vIdx][0]
					forb1 |= interferenceBits[vIdx][1]
					state.varColor[vIdx] = int16(colorCount)
				} else {
					kept = append(kept, rawIdx)
				}
			}
			worklist = kept
			colorCount++
		}
	}

	clear(state.superNode[:colorCount])
	state.orderedVars = orderedVariables
	state.hasCoalescing = c.usePseudoNames || colorCount < numCandidates
}

func (c *CoalesceVariableNames) processCfgNodeBits128(
	cfgNode *graph.LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch],
	interferenceBits [][2]uint64,
) {
	state := c.liveness.FlowStateAt(cfgNode.Index)

	inSet := &state.GetIn().liveSet
	inW0, inW1 := inSet.w0, inSet.w1
	for w := inW0; w != 0; w &= w - 1 {
		i := bits.TrailingZeros64(w)
		interferenceBits[i][0] |= inW0
		interferenceBits[i][1] |= inW1
	}
	for w := inW1; w != 0; w &= w - 1 {
		i := 64 + bits.TrailingZeros64(w)
		interferenceBits[i][0] |= inW0
		interferenceBits[i][1] |= inW1
	}

	outSet := &state.GetOut().liveSet
	outW0, outW1 := outSet.w0, outSet.w1
	for w := outW0; w != 0; w &= w - 1 {
		i := bits.TrailingZeros64(w)
		interferenceBits[i][0] |= outW0
		interferenceBits[i][1] |= outW1
	}
	for w := outW1; w != 0; w &= w - 1 {
		i := 64 + bits.TrailingZeros64(w)
		interferenceBits[i][0] |= outW0
		interferenceBits[i][1] |= outW1
	}

	lrc := &c.liveRangeChecker
	val := cfgNode.Value
	lrc.root = val
	lrc.state = state
	lrc.isAssignToList = lrc.isAssignToList[:0]
	lrc.isReadFromList = lrc.isReadFromList[:0]
	lrc.Check(val)
	if len(lrc.isAssignToList) > 0 {
		lrc.setCrossingVariablesBits128(interferenceBits, outW0, outW1)
	}
}

// IsInMultipleLvalueDecl returns whether this variable's declaration also
// declares other names.
//
// For example, this would return true for `x` in `let [x, y, z] = []`;
//
// Java: CoalesceVariableNames.isInMultipleLvalueDecl(Var)
func (c *CoalesceVariableNames) IsInMultipleLvalueDecl(v *Var) bool {
	declarationType, ok := v.DeclarationToken()
	if !ok {
		// switch on a null enum throws.
		panic(jutil.NewNullPointerException())
	}
	switch declarationType {
	case rhino.Token_LET, rhino.Token_CONST, rhino.Token_VAR:
		nameDecl := NodeUtil_GetEnclosingNode(v.GetNode(), NodeUtil_IsNameDeclaration)
		count := 0
		NodeUtil_VisitLhsNodesInNode(nameDecl, func(lhs rhino.Node) { count++ })
		return count > 1
	default:
		return false
	}
}

// RemoveVarDeclaration removes variable declaration if the variable has been
// coalesced with another variable that has already been declared.
//
// Java: CoalesceVariableNames.removeVarDeclaration(Node, Var)
func (c *CoalesceVariableNames) RemoveVarDeclaration(name rhino.Node, originalVar *Var) {
	decl := ast.MustNameDeclaration(NodeUtil_GetEnclosingNode(name, NodeUtil_IsNameDeclaration))
	var_ := decl.Node()
	parent := var_.GetParent()

	if dlhs, ok := decl.MustFirstDeclarator().AsDestructuringLhs(); ok {
		// convert `const [x] = arr` to `[x] = arr`
		// a precondition for this method is that `x` is the only lvalue in the destructuring pattern
		if _, isEnhFor := ast.AsEnhancedFor(parent); isEnhFor {
			var_.ReplaceWith(dlhs.Pattern().Detach().Node())
		} else {
			init, _ := dlhs.Init()
			rvalue := init.Detach()
			pattern := dlhs.Pattern().Detach()
			assign := ast.NewAssign(pattern.AsTarget(), rvalue).Node().Srcref(var_)
			var_.ReplaceWith(NodeUtil_NewExpr(assign))
		}
	} else if _, isEnhFor := ast.AsEnhancedFor(parent); isEnhFor {
		// convert `for (let x of ...` to `for (x of ...`
		var_.ReplaceWith(name.Detach())
	} else {
		// either `var x = 0;` or `var x;`
		onlyDecl, isOnly := decl.OnlyDeclarator()
		nameView, isOnlyName := onlyDecl.AsName()
		jutil.CheckStateMsg1(isOnly && isOnlyName && nameView.Is(name), "%s", var_)
		if value, hasInit := nameView.Init(); hasInit {
			// convert `let x = 0;` to `x = 0;`
			value.Detach()
			nameView.Detach()
			assign := ast.NewAssign(nameView.AsTarget(), value).Node().Srcref(name)

			// We don't need to wrap it with EXPR node if it is within a FOR.
			if !parent.IsVanillaFor() {
				assign = NodeUtil_NewExpr(assign)
			}
			var_.ReplaceWith(assign)
		} else if c.IsUninitializedLetNameInLoopBodyAndLive(name, originalVar) {
			nameView.Detach()
			undefinedValue := ast.MustExpr(c.astFactory.CreateUndefinedValue().SrcrefTree(name))
			assign := ast.NewAssign(nameView.AsTarget(), undefinedValue).Node().Srcref(name)
			if !parent.IsVanillaFor() {
				assign = NodeUtil_NewExpr(assign)
			}
			var_.ReplaceWith(assign)
		} else {
			// convert `let x;` to ``
			// and `for (let x;;) {}` to `for (;;) {}`
			NodeUtil_RemoveChild(parent, var_)
		}
	}
}

// MakeDeclarationVar converts `const` or `let` declarations to `var`
// declarations.
//
// Java: CoalesceVariableNames.makeDeclarationVar(Var)
func (c *CoalesceVariableNames) MakeDeclarationVar(coalescedName *Var) {
	if coalescedName.IsConst() || coalescedName.IsLet() {
		nameNode := rhino.CheckNotNullMsg1(coalescedName.GetNameNode(), "%s", coalescedName)
		if CoalesceVariableNames_IsUninitializedLetNameInLoopBody(nameNode) {
			undefinedValue := ast.MustExpr(c.astFactory.CreateUndefinedValue().SrcrefTree(nameNode))
			ast.MustName(nameNode).SetInit(undefinedValue)
		}
		// find the declaration node in a way that works normal and destructuring declarations.
		declNode := NodeUtil_GetEnclosingNode(nameNode.GetParent(), NodeUtil_IsNameDeclaration)
		// normalization ensures that all variables in a function are uniquely named, so it's OK
		// to turn a `const` or `let` into a `var`.
		declNode.SetToken(rhino.Token_VAR)
	}
}

// Java: CoalesceVariableNames.isUninitializedLetNameInLoopBody(Node)
func CoalesceVariableNames_IsUninitializedLetNameInLoopBody(nameNode rhino.Node) bool {
	name, isName := ast.AsName(nameNode)
	jutil.CheckStateMsg1(isName, "%s", nameNode)
	letDecl, isLet := ast.AsNameDeclarationOf(nameNode.GetParent(), rhino.Token_LET)
	if !isLet {
		return false
	}
	if _, hasInit := name.Init(); hasInit {
		return false
	}

	letParent := letDecl.Node().GetParent()
	if _, isLoop := ast.AsLoop(letParent); isLoop {
		return false
	}
	// Inside a loop body, but not the loop control node itself
	return NodeUtil_IsWithinLoop(letParent)
}

// Java: CoalesceVariableNames.isUninitializedLetNameInLoopBodyAndLive(Node, Var)
func (c *CoalesceVariableNames) IsUninitializedLetNameInLoopBodyAndLive(nameNode rhino.Node, originalVar *Var) bool {
	if !CoalesceVariableNames_IsUninitializedLetNameInLoopBody(nameNode) {
		return false
	}
	jutil.CheckState(c.liveness != nil)
	letNode := nameNode.GetParent()
	cfgNode := jutil.CheckNotNull(c.liveness.GetCfg().Nodes.Get(letNode))
	state := c.liveness.FlowStateAt(cfgNode.Index)
	varIndex := c.liveness.getVarIndexByID(originalVar.nameID)
	return varIndex >= 0 && state.GetOut().IsLive(varIndex)
}

// CoalesceVariableNames_LiveRangeChecker is used to find written and read
// variables in the same CFG node so that the variable pairs can be marked as
// interfering in an interference bit map.
type CoalesceVariableNames_LiveRangeChecker struct {
	root             rhino.Node
	state            *DataFlowAnalysis_LinearFlowState[*LiveVariablesAnalysis_LiveVariableLattice]
	orderedVariables []*Var
	liveness         *LiveVariablesAnalysis
	isAssignToList   []int // indices of written variables
	isReadFromList   []int // indices of read variables
}

// Java: CoalesceVariableNames.LiveRangeChecker(Node, Var[], LinearFlowState<LiveVariableLattice>)
func NewCoalesceVariableNames_LiveRangeChecker(
	root rhino.Node, orderedVariables []*Var, state *DataFlowAnalysis_LinearFlowState[*LiveVariablesAnalysis_LiveVariableLattice]) *CoalesceVariableNames_LiveRangeChecker {
	return &CoalesceVariableNames_LiveRangeChecker{
		root:             root,
		orderedVariables: orderedVariables,
		state:            state,
	}
}

func (l *CoalesceVariableNames_LiveRangeChecker) varIndexByID(id uint32) int {
	if id == 0 {
		return -1
	}
	if l.liveness != nil {
		if idx, ok := l.liveness.lookupVarIndexByID(id); ok {
			return idx
		}
		return -1
	}
	for i := len(l.orderedVariables) - 1; i >= 0; i-- {
		if l.orderedVariables[i].nameID == id {
			return i
		}
	}
	return -1
}

// Java: CoalesceVariableNames.LiveRangeChecker.check(Node)
func (l *CoalesceVariableNames_LiveRangeChecker) Check(n rhino.Node) {
	// For most AST nodes, traverse the subtree in postorder because that's how the expressions
	// are evaluated.
	if n == l.root || !ControlFlowGraph_IsEnteringNewCfgNode(n) {
		if dlhs, ok := ast.AsDestructuringLhs(n); ok {
			if init, hasInit := dlhs.Init(); hasInit {
				// Evaluate the rhs of a destructuring declaration before the lhs.
				l.Check(init.Node())
				l.Check(dlhs.Pattern().Node())
			} else {
				l.Check(dlhs.Pattern().Node())
			}
		} else if assign, ok := ast.AsAssign(n); ok {
			if pat, isPat := assign.Target().AsPattern(); isPat {
				// Evaluate the rhs of a destructuring assignment before the lhs.
				l.Check(assign.Value().Node())
				l.Check(pat.Node())
			} else {
				l.Check(assign.Target().Node())
				l.Check(assign.Value().Node())
			}
		} else if dv, ok := ast.AsDefaultValue(n); ok {
			l.Check(dv.Value().Node())
			l.Check(dv.Target().Node())
		} else {
			for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
				l.Check(c)
			}
		}
		if CoalesceVariableNames_LiveRangeChecker_ShouldVisit(n) {
			l.Visit(n, n.GetParent())
		}
	}
}

// Java: CoalesceVariableNames.LiveRangeChecker.visit(Node, Node)
func (l *CoalesceVariableNames_LiveRangeChecker) Visit(n rhino.Node, parent rhino.Node) {
	var assignedID uint32
	if name, isName := ast.AsName(n); isName {
		_, hasInit := name.Init()
		if parent.IsParamList() ||
			(NodeUtil_IsNameDeclaration(parent) && hasInit) ||
			NodeUtil_IsLhsByDestructuring(n) {
			assignedID = n.GetStringID()
		}
	} else if assign, isAssign := ast.AsAssignment(n); isAssign {
		if target := assign.Target().Node(); target.IsName() {
			assignedID = target.GetStringID()
		}
	}
	if assignedID != 0 {
		if iVar := l.varIndexByID(assignedID); iVar >= 0 {
			l.isAssignToList = append(l.isAssignToList, iVar)
		}
	}
	if len(l.isAssignToList) != 0 {
		if n.IsName() && !NodeUtil_IsNameDeclOrSimpleAssignLhs(n, parent) {
			if iVar := l.varIndexByID(n.GetStringID()); iVar >= 0 {
				l.isReadFromList = append(l.isReadFromList, iVar)
			}
		}
	}
}

func (l *CoalesceVariableNames_LiveRangeChecker) setCrossingVariablesBits128(
	interferenceBits [][2]uint64,
	outW0, outW1 uint64,
) {
	var writtenW0, writtenW1 uint64
	var readW0, readW1 uint64
	for _, iReadVar := range l.isReadFromList {
		if iReadVar < 64 {
			readW0 |= uint64(1) << uint(iReadVar)
		} else {
			readW1 |= uint64(1) << uint(iReadVar-64)
		}
	}
	for _, iWrittenVar := range l.isAssignToList {
		if iWrittenVar < 64 {
			writtenW0 |= uint64(1) << uint(iWrittenVar)
		} else {
			writtenW1 |= uint64(1) << uint(iWrittenVar-64)
		}
		interferenceBits[iWrittenVar][0] |= outW0 | readW0
		interferenceBits[iWrittenVar][1] |= outW1 | readW1
	}
	unionRead0 := outW0 | readW0
	unionRead1 := outW1 | readW1
	for w := unionRead0; w != 0; w &= w - 1 {
		iReadVar := bits.TrailingZeros64(w)
		interferenceBits[iReadVar][0] |= writtenW0
		interferenceBits[iReadVar][1] |= writtenW1
	}
	for w := unionRead1; w != 0; w &= w - 1 {
		iReadVar := 64 + bits.TrailingZeros64(w)
		interferenceBits[iReadVar][0] |= writtenW0
		interferenceBits[iReadVar][1] |= writtenW1
	}
}

// CoalesceVariableNames_LiveRangeChecker_ShouldVisit returns whether any
// LiveRangeChecker would be interested in the node.
//
// Java: CoalesceVariableNames.LiveRangeChecker.shouldVisit(Node)
func CoalesceVariableNames_LiveRangeChecker_ShouldVisit(n rhino.Node) bool {
	return n.IsName() || (n.HasChildren() && n.GetFirstChild().IsName())
}
