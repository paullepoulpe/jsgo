// Port of com/google/javascript/jscomp/LiveVariablesAnalysis.java @ closure-compiler v20261021

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
	"math/bits"
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/graph"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

var lvaArgumentsIdentID = uint32(strpool.Intern("arguments"))

const (
	LiveVariablesAnalysis_MAX_VARIABLES_TO_ANALYZE = 100
	lvaInlineVarCap                                = 16
)

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
// global and they can be retrieved with GetEscapedLocals().
type LiveVariablesAnalysis struct {
	cfg      *ControlFlowGraph[rhino.Node]
	worklist cfgWorklist
	states   []DataFlowAnalysis_LinearFlowState[*LiveVariablesAnalysis_LiveVariableLattice]

	// The scope of the function that we are analyzing.
	jsScope *Scope

	// The scope of the body of the function that we are analyzing.
	jsScopeChild *Scope

	// escaped marks indices in orderedVars that are escaped.
	escaped dfaBitSet
	// untracked marks indices in orderedVars that are !local || escaped.
	untracked           dfaBitSet
	escapedSet          jutil.ObjectLinkedHashSet[*Var]
	escapedSetInit      bool
	paramsMarkedEscaped bool

	// varIDs holds interned variable strpool IDs inline when numVars <= 16;
	// idToIndex is reused across Reset calls when numVars > 16.
	varIDs            [lvaInlineVarCap]uint32
	idToIndex         map[uint32]int
	numVars           int
	uniqueVarCount    int
	hasDuplicateNames bool

	// obtain variables in the order in which they appear in the code
	orderedVars []*Var

	allVarsSrc *NodeUtil_AllVarsDeclaredInFunction

	emptyLattice LiveVariablesAnalysis_LiveVariableLattice
	lattices     jutil.ReuseSlab[LiveVariablesAnalysis_LiveVariableLattice]
}

// LiveVariablesAnalysis_LiveVariableJoinOp is the inner class
// LiveVariableJoinOp.
type LiveVariablesAnalysis_LiveVariableJoinOp struct {
	result *LiveVariablesAnalysis_LiveVariableLattice
}

// Java: LiveVariablesAnalysis.LiveVariableJoinOp()
func NewLiveVariablesAnalysis_LiveVariableJoinOp(outer *LiveVariablesAnalysis) *LiveVariablesAnalysis_LiveVariableJoinOp {
	return &LiveVariablesAnalysis_LiveVariableJoinOp{
		result: NewLiveVariablesAnalysis_LiveVariableLattice(len(outer.orderedVars)),
	}
}

// Java: LiveVariablesAnalysis.LiveVariableJoinOp.joinFlow(LiveVariableLattice)
func (j *LiveVariablesAnalysis_LiveVariableJoinOp) JoinFlow(x *LiveVariablesAnalysis_LiveVariableLattice) {
	j.result.liveSet.orWith(&x.liveSet)
}

// Java: LiveVariablesAnalysis.LiveVariableJoinOp.finish()
func (j *LiveVariablesAnalysis_LiveVariableJoinOp) Finish() *LiveVariablesAnalysis_LiveVariableLattice {
	return j.result
}

// LiveVariablesAnalysis_LiveVariableLattice is the lattice that stores the
// liveness of all local variables at a given point in the program. The whole
// lattice is the power set of all local variables and a variable is live if
// it is in the set.
type LiveVariablesAnalysis_LiveVariableLattice struct {
	liveSet dfaBitSet
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
	return &LiveVariablesAnalysis_LiveVariableLattice{liveSet: other.liveSet.clone()}
}

// Java: LiveVariablesAnalysis.LiveVariableLattice.equals(Object)
func (l *LiveVariablesAnalysis_LiveVariableLattice) Equals(other any) bool {
	jutil.CheckNotNull(other)
	liveVariableLattice, ok := other.(*LiveVariablesAnalysis_LiveVariableLattice)
	return ok && (l == liveVariableLattice || l.liveSet.equals(&liveVariableLattice.liveSet))
}

// IsLive: there is only a version of this function with index since
// var.index will return the wrong one. Use an instantiation of
// LiveVariablesAnalysis and getVarIndex(var) to get the right index.
//
// Java: LiveVariablesAnalysis.LiveVariableLattice.isLive(int)
func (l *LiveVariablesAnalysis_LiveVariableLattice) IsLive(index int) bool {
	if index < 0 {
		panic(jutil.NewIndexOutOfBoundsException("bitIndex < 0: " + strconv.Itoa(index)))
	}
	return l.liveSet.has(index)
}

// Java: LiveVariablesAnalysis.LiveVariableLattice.toString()
func (l *LiveVariablesAnalysis_LiveVariableLattice) String() string {
	if l.liveSet.isEmpty() {
		return "{}"
	}
	buf := make([]byte, 0, 32)
	buf = append(buf, '{')
	first := true
	appendWord := func(base int, word uint64) {
		for w := word; w != 0; w &= w - 1 {
			if !first {
				buf = append(buf, ',', ' ')
			}
			first = false
			buf = strconv.AppendInt(buf, int64(base+bits.TrailingZeros64(w)), 10)
		}
	}
	appendWord(0, l.liveSet.w0)
	appendWord(64, l.liveSet.w1)
	for i, w := range l.liveSet.extra {
		appendWord(128+i*64, w)
	}
	buf = append(buf, '}')
	return string(buf)
}

// NextSetBit returns the index of the first bit that is set to true that
// occurs on or after the specified starting index.
//
// Java: LiveVariablesAnalysis.LiveVariableLattice.nextSetBit(int)
func (l *LiveVariablesAnalysis_LiveVariableLattice) NextSetBit(fromIndex int) int {
	if fromIndex < 0 {
		panic(jutil.NewIndexOutOfBoundsException("fromIndex < 0: " + strconv.Itoa(fromIndex)))
	}
	return l.liveSet.nextSetBit(fromIndex)
}

// Java: LiveVariablesAnalysis.LiveVariableLattice.hashCode()
func (l *LiveVariablesAnalysis_LiveVariableLattice) HashCode() int32 {
	var h int64 = 1234
	h ^= int64(l.liveSet.w0)
	h ^= int64(l.liveSet.w1) * 2
	for i, w := range l.liveSet.extra {
		h ^= int64(w) * int64(i+3)
	}
	return int32((h >> 32) ^ h)
}

// NewLiveVariablesAnalysis is Live Variables Analysis using the ES6 scope
// creator. This analysis should only be done on function where jsScope is
// the function scope. If we call LiveVariablesAnalysis from the function
// scope of our pass, we can pass a null value for the JsScopeChild, but if
// we call it from the function block scope, then JsScopeChild will be the
// function block scope.
//
// We call from the function scope when the pass requires us to traverse
// nodes beginning at the function parameters, and it from the function block
// scope when we are ignoring function parameters.
//
// jsScope is the function scope; jsScopeChild null or function block scope;
// scopeCreator Es6 Scope creator; allVarsDeclaredInFunction mapping of names
// to vars of everything reachable in a function.
//
// Java: LiveVariablesAnalysis(ControlFlowGraph<Node>, Scope, Scope, AbstractCompiler, ScopeCreator, AllVarsDeclaredInFunction)
func NewLiveVariablesAnalysis(
	cfg *ControlFlowGraph[rhino.Node],
	jsScope *Scope,
	jsScopeChild *Scope,
	compiler AbstractCompiler,
	scopeCreator ScopeCreator,
	allVarsDeclaredInFunction *NodeUtil_AllVarsDeclaredInFunction) *LiveVariablesAnalysis {
	l := &LiveVariablesAnalysis{}
	l.Reset(cfg, jsScope, jsScopeChild, compiler, scopeCreator, allVarsDeclaredInFunction)
	return l
}

// Reset resets l for reuse on a new function scope while retaining slab,
// queue, set, and map capacities.
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
	l.escaped.copyFrom(&table.escaped)
	l.untracked.clear(releaseLarge)
	l.paramsMarkedEscaped = false
	l.escapedSetInit = false
	l.allVarsSrc = allVarsDeclaredInFunction
	l.orderedVars = table.vars
	l.numVars = len(table.vars)
	l.uniqueVarCount = table.uniqueCount
	l.hasDuplicateNames = table.hasDuplicateNames
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
	l.numVars = len(l.orderedVars)
	useInline := l.numVars <= lvaInlineVarCap
	reuseTableMap := !useInline && l.allVarsSrc != nil && l.allVarsSrc.table.idToIndex != nil
	if reuseTableMap {
		l.idToIndex = l.allVarsSrc.table.idToIndex
	} else if !useInline {
		if l.idToIndex == nil {
			l.idToIndex = make(map[uint32]int, l.numVars)
		} else {
			clear(l.idToIndex)
		}
	}
	for num, v := range l.orderedVars {
		id := v.nameID
		if useInline {
			l.varIDs[num] = id
		} else if !reuseTableMap {
			l.idToIndex[id] = num
		}
		localScope, _ := v.GetScope().(*Scope)
		var local bool
		if localScope.IsFunctionBlockScope() {
			local = isDeclaredInFunctionBlockOrParameterByID(localScope, id)
		} else if localScope == l.jsScope && l.jsScopeChild != nil {
			local = isDeclaredInFunctionBlockOrParameterByID(l.jsScopeChild, id)
		} else {
			local = localScope.hasOwnSlotByID(id)
		}
		if !local || l.escaped.has(num) {
			l.untracked.set(num)
		}
	}
}

func (l *LiveVariablesAnalysis) lookupVarIndexByID(id uint32) (int, bool) {
	if id == 0 || l.numVars == 0 {
		return 0, false
	}
	if l.numVars <= lvaInlineVarCap {
		for i := l.numVars - 1; i >= 0; i-- {
			if l.varIDs[i] == id {
				return i, true
			}
		}
		return 0, false
	}
	idx, ok := l.idToIndex[id]
	return idx, ok
}

func (l *LiveVariablesAnalysis) lookupVarIndex(name string) (int, bool) {
	id := rhino.RhinoStringPool_LookupID(name)
	if id == 0 {
		return 0, false
	}
	return l.lookupVarIndexByID(id)
}

func (l *LiveVariablesAnalysis) getVarIndexByID(id uint32) int {
	index, ok := l.lookupVarIndexByID(id)
	if !ok {
		panic(jutil.NewNullPointerException())
	}
	return index
}

func (l *LiveVariablesAnalysis) isEscapedIndex(idx int) bool {
	return l.escaped.has(idx)
}

func (l *LiveVariablesAnalysis) allInFnVarsEscaped() bool {
	if l.escaped.count() < l.uniqueVarCount {
		return false
	}
	if !l.hasDuplicateNames {
		return true
	}
	for i, v := range l.orderedVars {
		if canonIdx, _ := l.lookupVarIndexByID(v.nameID); canonIdx == i && !l.escaped.has(i) {
			return false
		}
	}
	return true
}

// Java: LiveVariablesAnalysis.getEscapedLocals()
func (l *LiveVariablesAnalysis) GetEscapedLocals() *jutil.ObjectLinkedHashSet[*Var] {
	if !l.escapedSetInit {
		l.escapedSet.Clear()
		l.escaped.forEach(func(idx int) {
			if idx < len(l.orderedVars) {
				l.escapedSet.Add(l.orderedVars[idx])
			}
		})
		if l.paramsMarkedEscaped {
			params := ast.MustFunction(l.jsScope.GetRootNode()).Params()
			for param, ok := params.FirstParam(); ok; param, ok = param.Next() {
				if param.Token() == rhino.Token_NAME {
					if pv := l.jsScope.GetVarByID(param.Node().GetStringID()); pv != nil {
						l.escapedSet.Add(pv)
					}
				}
			}
		}
		l.escapedSetInit = true
	}
	return &l.escapedSet
}

// Java: LiveVariablesAnalysis.getAllVariables()
func (l *LiveVariablesAnalysis) GetAllVariables() *jutil.LinkedHashMap[string, *Var] {
	return l.allVarsSrc.GetAllVariables()
}

// Java: LiveVariablesAnalysis.getAllVariablesInOrder()
func (l *LiveVariablesAnalysis) GetAllVariablesInOrder() []*Var {
	return l.orderedVars
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
	return a.liveSet.equals(&b.liveSet)
}

// Java: DataFlowAnalysis.join(L, L)
func (l *LiveVariablesAnalysis) Join(latticeA, latticeB *LiveVariablesAnalysis_LiveVariableLattice) *LiveVariablesAnalysis_LiveVariableLattice {
	jutil.CheckNotNull(latticeA)
	jutil.CheckNotNull(latticeB)
	result := l.allocLattice()
	result.liveSet.orWith(&latticeA.liveSet)
	result.liveSet.orWith(&latticeB.liveSet)
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
		node.SetAnnotation(state)
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
	var w0, w1 uint64
	maxExtra := 0
	distinct := false
	for _, edge := range outEdges {
		dest := edge.DestNode
		if dest == implicitReturn {
			continue
		}
		inLat := l.states[dest.Index].in
		if inLat == &l.emptyLattice {
			continue
		}
		w0 |= inLat.liveSet.w0
		w1 |= inLat.liveSet.w1
		if n := len(inLat.liveSet.extra); n > maxExtra {
			maxExtra = n
		}
		if first == nil {
			first = inLat
		} else if !distinct && first != inLat && !first.liveSet.equals(&inLat.liveSet) {
			distinct = true
		}
	}
	if first == nil {
		state.out = &l.emptyLattice
		return
	}
	if !distinct {
		state.out = first
		return
	}
	if maxExtra == 0 {
		if len(first.liveSet.extra) == 0 && first.liveSet.w0 == w0 && first.liveSet.w1 == w1 {
			state.out = first
			return
		}
		if prev := state.out; prev != nil && len(prev.liveSet.extra) == 0 && prev.liveSet.w0 == w0 && prev.liveSet.w1 == w1 {
			return
		}
		joined := l.allocLattice()
		joined.liveSet.w0 = w0
		joined.liveSet.w1 = w1
		state.out = joined
		return
	}
	extra := make([]uint64, maxExtra)
	for _, edge := range outEdges {
		dest := edge.DestNode
		if dest == implicitReturn {
			continue
		}
		inLat := l.states[dest.Index].in
		for i, ew := range inLat.liveSet.extra {
			extra[i] |= ew
		}
	}
	joined := l.allocLattice()
	joined.liveSet.w0 = w0
	joined.liveSet.w1 = w1
	joined.liveSet.extra = extra
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
	var gen, kill dfaBitSet

	// Make kills conditional if the node can end abruptly by an exception.
	conditional := false
	for _, edge := range edgeList {
		if ControlFlowGraph_Branch_ON_EX == edge.GetValue() {
			conditional = true
			break
		}
	}
	l.computeGenKillBits(node, &gen, &kill, conditional)
	// L_in = L_out - Kill + Gen
	out0 := (input.liveSet.w0 &^ kill.w0) | gen.w0
	out1 := (input.liveSet.w1 &^ kill.w1) | gen.w1
	if len(input.liveSet.extra) == 0 && len(gen.extra) == 0 {
		if out0 == input.liveSet.w0 && out1 == input.liveSet.w1 {
			return input
		}
		if (out0 | out1) == 0 {
			return &l.emptyLattice
		}
		if prevIn != nil && len(prevIn.liveSet.extra) == 0 && prevIn.liveSet.w0 == out0 && prevIn.liveSet.w1 == out1 {
			return prevIn
		}
		result := l.allocLattice()
		result.liveSet.w0 = out0
		result.liveSet.w1 = out1
		return result
	}
	maxExtra := max(len(input.liveSet.extra), len(gen.extra))
	changed := out0 != input.liveSet.w0 || out1 != input.liveSet.w1
	if !changed {
		for i := 0; i < maxExtra; i++ {
			var inW, killW, genW uint64
			if i < len(input.liveSet.extra) {
				inW = input.liveSet.extra[i]
			}
			if i < len(kill.extra) {
				killW = kill.extra[i]
			}
			if i < len(gen.extra) {
				genW = gen.extra[i]
			}
			if (inW&^killW)|genW != inW {
				changed = true
				break
			}
		}
		if !changed {
			return input
		}
	}
	result := l.allocLattice()
	result.liveSet.w0 = out0
	result.liveSet.w1 = out1
	if maxExtra > 0 {
		extra := make([]uint64, maxExtra)
		for i := 0; i < maxExtra; i++ {
			var inW, killW, genW uint64
			if i < len(input.liveSet.extra) {
				inW = input.liveSet.extra[i]
			}
			if i < len(kill.extra) {
				killW = kill.extra[i]
			}
			if i < len(gen.extra) {
				genW = gen.extra[i]
			}
			extra[i] = (inW &^ killW) | genW
		}
		result.liveSet.extra = extra
	}
	return result
}

func lvaOrIntoBitSet(dst *jutil.BitSet, src *dfaBitSet) {
	for w := src.w0; w != 0; w &= w - 1 {
		dst.Set(bits.TrailingZeros64(w))
	}
	for w := src.w1; w != 0; w &= w - 1 {
		dst.Set(64 + bits.TrailingZeros64(w))
	}
	for i, word := range src.extra {
		base := 128 + i*64
		for w := word; w != 0; w &= w - 1 {
			dst.Set(base + bits.TrailingZeros64(w))
		}
	}
}

// ComputeGenKill computes the GEN and KILL set.
//
// n is the root node. gen: local variables that are live because of the
// instruction at n will be added to this set. kill: local variables that are
// killed because of the instruction at n will be added to this set.
// conditional: true if any assignments encountered are conditionally
// executed. These assignments might not kill a variable.
//
// Java: LiveVariablesAnalysis.computeGenKill(Node, BitSet, BitSet, boolean)
func (l *LiveVariablesAnalysis) ComputeGenKill(n rhino.Node, gen *jutil.BitSet, kill *jutil.BitSet, conditional bool) {
	var genBits, killBits dfaBitSet
	l.computeGenKillBits(n, &genBits, &killBits, conditional)
	lvaOrIntoBitSet(gen, &genBits)
	lvaOrIntoBitSet(kill, &killBits)
}

func (l *LiveVariablesAnalysis) computeGenKillBits(n rhino.Node, gen *dfaBitSet, kill *dfaBitSet, conditional bool) {
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

func (l *LiveVariablesAnalysis) addToSetIfLocalBits(node rhino.Node, set *dfaBitSet) {
	jutil.CheckStateMsg1(node.IsName(), "%s", node)
	idx, found := l.lookupVarIndexByID(node.GetStringID())
	if !found || l.untracked.has(idx) {
		return
	}
	set.set(idx)
}

// Java: LiveVariablesAnalysis.addToSetIfLocal(Node, BitSet)
func (l *LiveVariablesAnalysis) AddToSetIfLocal(node rhino.Node, set *jutil.BitSet) {
	jutil.CheckStateMsg1(node.IsName(), "%s", node)
	idx, found := l.lookupVarIndexByID(node.GetStringID())
	if !found || l.untracked.has(idx) {
		return
	}
	set.Set(idx)
}

func isDeclaredInFunctionBlockOrParameterByID(scope *Scope, id uint32) bool {
	jutil.CheckState(scope.IsFunctionBlockScope())
	return scope.hasOwnSlotByID(id) || scope.GetParent().hasOwnSlotByID(id)
}

// Java: LiveVariablesAnalysis.isDeclaredInFunctionBlockOrParameter(Scope, String)
func LiveVariablesAnalysis_IsDeclaredInFunctionBlockOrParameter(scope *Scope, name string) bool {
	// In ES6, we create a separate container scope above the function block scope to handle
	// default parameters. Since nothing in the function block scope is allowed to shadow
	// the variables in the function scope, we treat the two scopes as one in this method.
	jutil.CheckState(scope.IsFunctionBlockScope())
	id := rhino.RhinoStringPool_LookupID(name)
	return id != 0 && (scope.hasOwnSlotByID(id) || scope.GetParent().hasOwnSlotByID(id))
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
			if l.escapedSetInit && paramVar != nil {
				l.escapedSet.Add(paramVar)
			}
			if idx, ok := l.lookupVarIndexByID(id); ok {
				if paramVar == nil || l.orderedVars[idx] == paramVar || l.orderedVars[idx].Equals(paramVar) {
					l.escaped.set(idx)
					l.untracked.set(idx)
				} else {
					for i, cand := range l.orderedVars {
						if cand == paramVar || cand.Equals(paramVar) {
							l.escaped.set(i)
							break
						}
					}
				}
			}
		}
	}
}
