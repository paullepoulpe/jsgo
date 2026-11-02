// Port of com/google/javascript/jscomp/LiveVariablesAnalysis.java @ closure-compiler v20261028

/*
 * Copyright 2017 The Closure Compiler Authors.
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
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/graph"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

var lvaArgumentsIdentID = uint32(strpool.Intern("arguments"))

const LiveVariablesAnalysis_MAX_VARIABLES_TO_ANALYZE = 100

// LiveVariablesAnalysis computes the "liveness" of all local variables. A
// variable is "live" at a point of a program if the value it is currently
// holding might be read later. Otherwise, the variable is considered "dead"
// if we know for sure that it will no longer be read. Dead variables are
// candidates for dead assignment elimination and variable name sharing. The
// worst case safe assumption is to assume that all variables are live. In
// that case, we will have no opportunity for optimizations. This is
// especially the case within a TRY block when an assignment is not
// guaranteed to take place. We bail out by assuming that all variables are
// live.
//
// Due to the possibility of inner functions and closures, certain "local"
// variables can escape the function. These variables will be considered as
// global.
type LiveVariablesAnalysis struct {
	cfg      *ControlFlowGraph[rhino.Node]
	worklist cfgWorklist
	states   []DataFlowAnalysis_LinearFlowState[*LiveVariablesAnalysis_LiveVariableLattice]

	// The scope of the function that we are analyzing.
	jsScope *Scope

	// The scope of the body of the function that we are analyzing.
	jsScopeChild *Scope

	table *fnVarTable

	// escaped marks indices in table.vars that are escaped.
	escaped smallBitSet
	// untracked marks indices in table.vars that are !local || escaped.
	untracked           smallBitSet
	genScratch          smallBitSet
	killScratch         smallBitSet
	outScratch          smallBitSet
	paramsMarkedEscaped bool

	emptyLattice LiveVariablesAnalysis_LiveVariableLattice
	lattices     jutil.ReuseSlab[LiveVariablesAnalysis_LiveVariableLattice]
}

// LiveVariablesAnalysis_LiveVariableLattice is the lattice that stores the
// liveness of all local variables at a given point in the program. The whole
// lattice is the power set of all local variables and a variable is live if
// it is in the set.
type LiveVariablesAnalysis_LiveVariableLattice struct {
	liveSet smallBitSet
}

// numVars is the number of all local variables.
//
// Java: LiveVariablesAnalysis.LiveVariableLattice(int)
func NewLiveVariablesAnalysis_LiveVariableLattice(numVars int) *LiveVariablesAnalysis_LiveVariableLattice {
	if numVars < 0 {
		panic(jutil.NewNegativeArraySizeException("nbits < 0: " + strconv.Itoa(numVars)))
	}
	return &LiveVariablesAnalysis_LiveVariableLattice{}
}

// Java: LiveVariablesAnalysis.LiveVariableLattice(LiveVariableLattice)
func NewLiveVariablesAnalysis_LiveVariableLatticeOf(other *LiveVariablesAnalysis_LiveVariableLattice) *LiveVariablesAnalysis_LiveVariableLattice {
	jutil.CheckNotNull(other)
	return &LiveVariablesAnalysis_LiveVariableLattice{liveSet: other.liveSet.Clone()}
}

// Java: LiveVariablesAnalysis.LiveVariableLattice.equals(Object)
func (l *LiveVariablesAnalysis_LiveVariableLattice) Equals(other any) bool {
	jutil.CheckNotNull(other)
	liveVariableLattice, ok := other.(*LiveVariablesAnalysis_LiveVariableLattice)
	return ok && (l == liveVariableLattice || l.liveSet.Equals(&liveVariableLattice.liveSet))
}

// IsLive: there is only a version of this function with index since
// var.index will return the wrong one. Use an instantiation of
// LiveVariablesAnalysis and getVarIndex(var) to get the right index.
//
// Java: LiveVariablesAnalysis.LiveVariableLattice.isLive(int)
func (l *LiveVariablesAnalysis_LiveVariableLattice) IsLive(index int) bool {
	return l.liveSet.Get(index)
}

// Java: LiveVariablesAnalysis.LiveVariableLattice.toString()
func (l *LiveVariablesAnalysis_LiveVariableLattice) String() string {
	return l.liveSet.String()
}

// NextSetBit returns the index of the first bit that is set to true that
// occurs on or after the specified starting index.
//
// Java: LiveVariablesAnalysis.LiveVariableLattice.nextSetBit(int)
func (l *LiveVariablesAnalysis_LiveVariableLattice) NextSetBit(fromIndex int) int {
	return l.liveSet.NextSetBit(fromIndex)
}

// Java: LiveVariablesAnalysis.LiveVariableLattice.hashCode()
func (l *LiveVariablesAnalysis_LiveVariableLattice) HashCode() int32 {
	return l.liveSet.JavaHashCode()
}

// Reset resets l for reuse on a new function scope while retaining slab,
// queue, and bitset capacities.
func (l *LiveVariablesAnalysis) Reset(
	cfg *ControlFlowGraph[rhino.Node],
	jsScope *Scope,
	jsScopeChild *Scope,
	compiler AbstractCompiler,
	scopeCreator ScopeCreator,
	allVarsDeclaredInFunction *NodeUtil_AllVarsDeclaredInFunction) {
	l.cfg = cfg
	l.worklist.releaseLarge(dfaReuseMaxNodes)
	if cap(l.states) > dfaReuseMaxNodes {
		l.states = nil
	}
	jutil.CheckStateMsg1(jsScope.IsFunctionScope(), "%s", jsScope)

	table := &allVarsDeclaredInFunction.table
	releaseLarge := len(table.vars) <= 256
	table.escaped.clear(releaseLarge)
	computeEscapedInTable(jsScope, table, compiler, scopeCreator)

	l.jsScope = jsScope
	l.jsScopeChild = jsScopeChild
	l.table = table
	l.escaped.CopyFrom(&table.escaped)
	l.untracked.clear(releaseLarge)
	l.genScratch.clear(releaseLarge)
	l.killScratch.clear(releaseLarge)
	l.outScratch.clear(releaseLarge)
	l.paramsMarkedEscaped = false
	l.emptyLattice = LiveVariablesAnalysis_LiveVariableLattice{}
	l.lattices.Reset(nil)

	l.AddScopeVariables()
}

type lvaAllVarsFinder struct {
	finder nodeUtil_getAllVarsDeclaredInFunction_finder
}

func (f *lvaAllVarsFinder) EnterScope(t *NodeTraversal) {
	f.finder.EnterScope(t)
}

func (f *lvaAllVarsFinder) ExitScope(t *NodeTraversal) {}

func (f *lvaAllVarsFinder) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	return f.finder.ShouldTraverse(t, n, parent)
}

func (f *lvaAllVarsFinder) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {}

func (*lvaAllVarsFinder) nodeTraversalNoVisit() {}

func (f *lvaAllVarsFinder) collect(
	compiler AbstractCompiler, scopeCreator ScopeCreator, scope *Scope,
) *NodeUtil_AllVarsDeclaredInFunction {
	return f.collectWithLimit(compiler, scopeCreator, scope, 0)
}

func (f *lvaAllVarsFinder) collectWithLimit(
	compiler AbstractCompiler, scopeCreator ScopeCreator, scope *Scope, maxVars int,
) *NodeUtil_AllVarsDeclaredInFunction {
	jutil.CheckStateMsg1(scope.IsFunctionScope(), "%s", scope)
	if maxVars > 0 && scope.GetVarCount() >= maxVars {
		return nil
	}
	f.finder.result.allVarsInFn = nil
	f.finder.result.table.reset()
	f.finder.scope = scope

	fn := ast.MustFunction(scope.GetRootNode())
	if body, ok := fn.Body().AsBlock(); ok && !fsivHasNestedBlockScopeDecls(fn, body) {
		var bodyScope *Scope
		if body.HasStatements() {
			bodyScope = scopeCreator.CreateScope(body.Node(), scope).(*Scope)
			if maxVars > 0 && scope.GetVarCount()+bodyScope.GetVarCount() >= maxVars {
				return nil
			}
		}
		f.finder.appendScopeVars(scope)
		if bodyScope != nil {
			f.finder.appendScopeVars(bodyScope)
		}
		f.finder.result.table.finalize()
		jutil.CheckState((f.finder.result.table.uniqueCount == 0) == (len(f.finder.result.table.vars) == 0))
		return &f.finder.result
	}

	b := NodeTraversal_Builder{
		compiler:     compiler,
		callback:     f,
		scopeCreator: scopeCreator,
	}
	b.TraverseAtScope(scope)

	if maxVars > 0 && len(f.finder.result.table.vars) >= maxVars {
		return nil
	}
	f.finder.result.table.finalize()
	jutil.CheckState((f.finder.result.table.uniqueCount == 0) == (len(f.finder.result.table.vars) == 0))
	return &f.finder.result
}

func (l *LiveVariablesAnalysis) allocLattice() *LiveVariablesAnalysis_LiveVariableLattice {
	lat := l.lattices.Alloc()
	*lat = LiveVariablesAnalysis_LiveVariableLattice{}
	return lat
}

// AddScopeVariables: parameters belong to the function scope, but variables
// defined in the function body belong to the function body scope. Assign a
// unique index to each variable, regardless of which scope it's in.
//
// Java: LiveVariablesAnalysis.addScopeVariables()
func (l *LiveVariablesAnalysis) AddScopeVariables() {
	for num, v := range l.table.vars {
		id := v.nameID
		localScope, _ := v.GetScope().(*Scope)
		var local bool
		if localScope.IsFunctionBlockScope() {
			local = isDeclaredInFunctionBlockOrParameterByID(localScope, id)
		} else if localScope == l.jsScope && l.jsScopeChild != nil {
			local = isDeclaredInFunctionBlockOrParameterByID(l.jsScopeChild, id)
		} else {
			local = localScope.hasOwnSlotByID(id)
		}
		if !local || l.escaped.Get(num) {
			l.untracked.Set(num)
		}
	}
}

func (l *LiveVariablesAnalysis) lookupVarIndexByID(id uint32) (int, bool) {
	return l.table.lookupIndexByID(id)
}

func (l *LiveVariablesAnalysis) lookupVarIndex(name string) (int, bool) {
	return l.table.lookupIndexByID(rhino.RhinoStringPool_LookupID(name))
}

func (l *LiveVariablesAnalysis) getVarIndexByID(id uint32) int {
	index, ok := l.table.lookupIndexByID(id)
	if !ok {
		panic(jutil.NewNullPointerException())
	}
	return index
}

func (l *LiveVariablesAnalysis) isEscapedIndex(idx int) bool {
	return l.escaped.Get(idx)
}

func (l *LiveVariablesAnalysis) allInFnVarsEscaped() bool {
	if l.escaped.Cardinality() < l.table.uniqueCount {
		return false
	}
	if !l.table.hasDuplicateNames {
		return true
	}
	for i, v := range l.table.vars {
		if canonIdx, _ := l.table.lookupIndexByID(v.nameID); canonIdx == i && !l.escaped.Get(i) {
			return false
		}
	}
	return true
}

// Java: LiveVariablesAnalysis.getAllVariablesInOrder()
func (l *LiveVariablesAnalysis) GetAllVariablesInOrder() []*Var {
	return l.table.vars
}

// Java: LiveVariablesAnalysis.getVarIndex(String)
func (l *LiveVariablesAnalysis) GetVarIndex(var_ string) int {
	index, ok := l.lookupVarIndex(var_)
	if !ok {
		// Unboxing the null Integer that Map.get returns.
		panic(jutil.NewNullPointerException())
	}
	return index
}

// Java: DataFlowAnalysis.getCfg()
func (l *LiveVariablesAnalysis) GetCfg() *ControlFlowGraph[rhino.Node] {
	return l.cfg
}

func (l *LiveVariablesAnalysis) FlowStateAt(index uint32) *DataFlowAnalysis_LinearFlowState[*LiveVariablesAnalysis_LiveVariableLattice] {
	return &l.states[index]
}

// Java: LiveVariablesAnalysis.isForward()
func (l *LiveVariablesAnalysis) IsForward() bool {
	return false
}

// Java: DataFlowAnalysis.isBranched()
func (l *LiveVariablesAnalysis) IsBranched() bool {
	return false
}

// Java: LiveVariablesAnalysis.createEntryLattice()
func (l *LiveVariablesAnalysis) CreateEntryLattice() *LiveVariablesAnalysis_LiveVariableLattice {
	return &l.emptyLattice
}

// Java: LiveVariablesAnalysis.createInitialEstimateLattice()
func (l *LiveVariablesAnalysis) CreateInitialEstimateLattice() *LiveVariablesAnalysis_LiveVariableLattice {
	return &l.emptyLattice
}

func (l *LiveVariablesAnalysis) EqualsLattice(a, b *LiveVariablesAnalysis_LiveVariableLattice) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.liveSet.Equals(&b.liveSet)
}

// Java: DataFlowAnalysis.join(L, L)
func (l *LiveVariablesAnalysis) Join(latticeA, latticeB *LiveVariablesAnalysis_LiveVariableLattice) *LiveVariablesAnalysis_LiveVariableLattice {
	jutil.CheckNotNull(latticeA)
	jutil.CheckNotNull(latticeB)
	result := l.allocLattice()
	result.liveSet.Or(&latticeA.liveSet)
	result.liveSet.Or(&latticeB.liveSet)
	return result
}

// Java: DataFlowAnalysis.analyze()
func (l *LiveVariablesAnalysis) Analyze() {
	cfg := l.cfg
	implicitReturn := cfg.implicitReturn
	cfgNodes := &cfg.Nodes
	numNodes := cfgNodes.Size()
	nodeCap := cfg.NodeIndexCap()

	l.worklist.reset(cfg, numNodes, false)
	if cap(l.states) < nodeCap {
		l.states = make([]DataFlowAnalysis_LinearFlowState[*LiveVariablesAnalysis_LiveVariableLattice], nodeCap)
	} else {
		clear(l.states[:max(len(l.states), nodeCap)])
		l.states = l.states[:nodeCap]
	}

	initLat := &l.emptyLattice
	states := l.states
	for node := range cfgNodes.Values() {
		node.DfaInQueue = false
		state := &states[node.Index]
		state.in = initLat
		state.out = initLat
		if node != implicitReturn {
			l.worklist.enqueue(node)
		}
	}

	for l.worklist.notEmpty() {
		curNode := l.worklist.dequeue()
		curState := &states[curNode.Index]
		stepCount := curState.stepCount
		curState.stepCount++
		if stepCount > DataFlowAnalysis_MAX_STEPS_PER_NODE {
			panic(jutil.NewIllegalStateExceptionMsg("Dataflow analysis appears to diverge around: " + jutil.ToJavaString(curNode)))
		}

		l.joinNodeInputs(curNode, curState, implicitReturn)
		inBefore := curState.in
		inAfter := l.flowThroughEdges(curNode.Value, curNode.OutEdgeList.Slice(), curState.out, inBefore)
		curState.in = jutil.CheckNotNull(inAfter)
		if !l.EqualsLattice(inBefore, inAfter) {
			for _, edge := range curNode.InEdgeList.Slice() {
				nextNode := edge.SourceNode
				if !nextNode.DfaInQueue && nextNode != implicitReturn {
					l.worklist.enqueue(nextNode)
				}
			}
		}
	}
}

func (l *LiveVariablesAnalysis) joinNodeInputs(
	node *graph.LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch],
	state *DataFlowAnalysis_LinearFlowState[*LiveVariablesAnalysis_LiveVariableLattice],
	implicitReturn *graph.LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch],
) {
	outEdges := node.OutEdgeList.Slice()
	switch len(outEdges) {
	case 0:
		return
	case 1:
		dest := outEdges[0].DestNode
		if dest == implicitReturn {
			state.out = &l.emptyLattice
		} else {
			state.out = l.states[dest.Index].in
		}
		return
	}

	var first *LiveVariablesAnalysis_LiveVariableLattice
	distinct := false
	l.outScratch.ClearAll()
	for _, edge := range outEdges {
		dest := edge.DestNode
		if dest == implicitReturn {
			continue
		}
		inLat := l.states[dest.Index].in
		if inLat == &l.emptyLattice {
			continue
		}
		l.outScratch.Or(&inLat.liveSet)
		if first == nil {
			first = inLat
		} else if !distinct && first != inLat && !first.liveSet.Equals(&inLat.liveSet) {
			distinct = true
		}
	}
	if first == nil {
		state.out = &l.emptyLattice
		return
	}
	if !distinct || first.liveSet.Equals(&l.outScratch) {
		state.out = first
		return
	}
	if prev := state.out; prev != nil && prev.liveSet.Equals(&l.outScratch) {
		return
	}
	joined := l.allocLattice()
	joined.liveSet = l.outScratch.Clone()
	state.out = joined
}

// Java: LiveVariablesAnalysis.flowThrough(Node, LiveVariableLattice)
func (l *LiveVariablesAnalysis) FlowThrough(node rhino.Node, input *LiveVariablesAnalysis_LiveVariableLattice) *LiveVariablesAnalysis_LiveVariableLattice {
	return l.flowThroughEdges(node, l.cfg.GetOutEdges(node), input, nil)
}

func (l *LiveVariablesAnalysis) flowThroughEdges(
	node rhino.Node,
	edgeList []*graph.LinkedDiGraphEdge[rhino.Node, ControlFlowGraph_Branch],
	input *LiveVariablesAnalysis_LiveVariableLattice,
	prevIn *LiveVariablesAnalysis_LiveVariableLattice,
) *LiveVariablesAnalysis_LiveVariableLattice {
	l.genScratch.ClearAll()
	l.killScratch.ClearAll()

	// Make kills conditional if the node can end abruptly by an exception.
	conditional := false
	for _, edge := range edgeList {
		if ControlFlowGraph_Branch_ON_EX == edge.GetValue() {
			conditional = true
			break
		}
	}
	l.computeGenKillBits(node, &l.genScratch, &l.killScratch, conditional)
	// L_in = L_out - Kill + Gen
	l.outScratch.SetAndNotOr(&input.liveSet, &l.killScratch, &l.genScratch)
	if l.outScratch.Equals(&input.liveSet) {
		return input
	}
	if l.outScratch.IsEmpty() {
		return &l.emptyLattice
	}
	if prevIn != nil && prevIn.liveSet.Equals(&l.outScratch) {
		return prevIn
	}
	result := l.allocLattice()
	result.liveSet = l.outScratch.Clone()
	return result
}

func (l *LiveVariablesAnalysis) computeGenKillBits(n rhino.Node, gen *smallBitSet, kill *smallBitSet, conditional bool) {
	switch n.GetToken() {
	case rhino.Token_SCRIPT, rhino.Token_ROOT, rhino.Token_FUNCTION, rhino.Token_BLOCK:
		return
	case rhino.Token_WHILE, rhino.Token_DO, rhino.Token_IF, rhino.Token_FOR:
		l.computeGenKillBits(NodeUtil_GetConditionExpression(n), gen, kill, conditional)
		return
	case rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF, rhino.Token_FOR_IN:
		{
			// for (x in y) {...}
			lhs := ast.MustEnhancedFor(n).Target().Node()
			if decl, ok := ast.AsNameDeclaration(lhs); ok {
				// for (var x in y) {...}
				if last, hasLast := decl.LastDeclarator(); hasLast {
					lhs = last.Node()
				} else {
					lhs = rhino.NilNode
				}
			}

			// Note that the LHS may never be assigned to or evaluated, like in:
			//   for (x in []) {}
			// so should not be killed.
			l.computeGenKillBits(lhs, gen, kill, conditional)

			// rhs is executed only once so we don't go into it every loop.
			return
		}
	case rhino.Token_LET, rhino.Token_CONST, rhino.Token_VAR:
		decl := ast.MustNameDeclaration(n)
		for d, ok := decl.FirstDeclarator(); ok; d, ok = d.Next() {
			if name, isName := d.AsName(); isName {
				if init, hasInit := name.Init(); hasInit {
					l.computeGenKillBits(init.Node(), gen, kill, conditional)
					if !conditional {
						l.addToSetIfLocalBits(name.Node(), kill)
					}
				}
			} else {
				dlhs := d.MustDestructuringLhs()
				if !conditional {
					NodeUtil_VisitLhsNodesInNode(dlhs.Node(), func(lhsNode rhino.Node) { l.addToSetIfLocalBits(lhsNode, kill) })
				}
				l.computeGenKillBits(dlhs.Pattern().Node(), gen, kill, conditional)
				if init, hasInit := dlhs.Init(); hasInit {
					l.computeGenKillBits(init.Node(), gen, kill, conditional)
				}
			}
		}
		return
	case rhino.Token_OPTCHAIN_GETPROP:
		// Go-only: OPTCHAIN_GETPROP has only one child (Object; property is its string),
		// so visiting Object once matches Java without O(2^k) re-traversal on chains.
		l.computeGenKillBits(ast.MustOptChainGetProp(n).Object().Node(), gen, kill, conditional)
		return
	case rhino.Token_AND, rhino.Token_OR, rhino.Token_COALESCE:
		bin := ast.MustBinary(n)
		l.computeGenKillBits(bin.Left().Node(), gen, kill, conditional)
		// May short circuit.
		l.computeGenKillBits(bin.Right().Node(), gen, kill, true)
		return
	case rhino.Token_OPTCHAIN_GETELEM:
		ge := ast.MustOptChainGetElem(n)
		l.computeGenKillBits(ge.Object().Node(), gen, kill, conditional)
		// May short circuit.
		l.computeGenKillBits(ge.Index().Node(), gen, kill, true)
		return
	case rhino.Token_OPTCHAIN_CALL:
		call := ast.MustOptChainCall(n)
		l.computeGenKillBits(call.Callee().Node(), gen, kill, conditional)
		// Unlike OPTCHAIN_GETPROP and OPTCHAIN_GETELEM, the OPTCHAIN_CALLs can have multiple
		// children on rhs which get executed conditionally
		for arg, ok := call.FirstArg(); ok; arg, ok = arg.Next() {
			l.computeGenKillBits(arg.Node(), gen, kill, true)
		}
		return
	case rhino.Token_HOOK:
		hook := ast.MustHook(n)
		l.computeGenKillBits(hook.Cond().Node(), gen, kill, conditional)
		// Assume both sides are conditional.
		l.computeGenKillBits(hook.Then().Node(), gen, kill, true)
		l.computeGenKillBits(hook.Else().Node(), gen, kill, true)
		return
	case rhino.Token_NAME:
		if n.GetStringID() == lvaArgumentsIdentID {
			l.MarkAllParametersEscaped()
		} else if !NodeUtil_IsLhsByDestructuring(n) {
			// Only add names in destructuring patterns if they're not lvalues.
			// e.g. "x" in "const {foo = x} = obj;"
			l.addToSetIfLocalBits(n, gen)
		}
		return
	default:
		if assign, ok := ast.AsAssignment(n); ok {
			if lhs, isName := assign.Target().AsName(); isName {
				if !conditional {
					l.addToSetIfLocalBits(lhs.Node(), kill)
				}
				if assign.Token() != rhino.Token_ASSIGN {
					// assignments such as a += 1 reads a.
					l.addToSetIfLocalBits(lhs.Node(), gen)
				}
				l.computeGenKillBits(assign.Value().Node(), gen, kill, conditional)
				return
			}
			if assign.Token() == rhino.Token_ASSIGN {
				if pat, isPat := assign.Target().AsPattern(); isPat {
					if !conditional {
						NodeUtil_VisitLhsNodesInNode(
							n,
							func(child rhino.Node) {
								if child.IsName() {
									l.addToSetIfLocalBits(child, kill)
								}
							})
					}
					l.computeGenKillBits(pat.Node(), gen, kill, conditional)
					l.computeGenKillBits(assign.Value().Node(), gen, kill, conditional)
					return
				}
			}
		}
		for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			l.computeGenKillBits(c, gen, kill, conditional)
		}
		return
	}
}

func (l *LiveVariablesAnalysis) addToSetIfLocalBits(node rhino.Node, set *smallBitSet) {
	jutil.CheckStateMsg1(node.IsName(), "%s", node)
	idx, found := l.table.lookupIndexByID(node.GetStringID())
	if !found || l.untracked.Get(idx) {
		return
	}
	set.Set(idx)
}

func isDeclaredInFunctionBlockOrParameterByID(scope *Scope, id uint32) bool {
	jutil.CheckState(scope.IsFunctionBlockScope())
	return scope.hasOwnSlotByID(id) || scope.GetParent().hasOwnSlotByID(id)
}

// MarkAllParametersEscaped gives up computing liveness of formal parameters
// by putting all the simple parameters in the escaped set.
//
// This only applies to simple parameters, that is NAMEs, because other
// parameter syntaxes never need to be escaped in this way. The known
// applications of this method are for uses of `arguments`, and for IE8. In a
// function with non-simple paremeters, `arguments` is not parameter-mapped,
// and so referencing it doesn't escape paremeters. IE8 just doess't support
// non-simple parameters.
//
// We could actaully continue tracking simple parameters if any parameter is
// non-simple, but it wasn't worth the complexity or cost to do so.
//
// See https://tc39.github.io/ecma262/#sec-functiondeclarationinstantiation
//
// Java: LiveVariablesAnalysis.markAllParametersEscaped()
func (l *LiveVariablesAnalysis) MarkAllParametersEscaped() {
	if l.paramsMarkedEscaped {
		return
	}
	l.paramsMarkedEscaped = true
	params := ast.MustFunction(l.jsScope.GetRootNode()).Params()
	for param, ok := params.FirstParam(); ok; param, ok = param.Next() {
		if param.Token() == rhino.Token_NAME {
			id := param.Node().GetStringID()
			paramVar := l.jsScope.GetVarByID(id)
			if paramVar != nil {
				if idx, ok := l.table.resolveExactVarIndex(paramVar); ok {
					l.escaped.Set(idx)
					if canonIdx, _ := l.table.lookupIndexByID(id); canonIdx == idx {
						l.untracked.Set(idx)
					}
				}
			} else if idx, ok := l.table.lookupIndexByID(id); ok {
				l.escaped.Set(idx)
				l.untracked.Set(idx)
			}
		}
	}
}
