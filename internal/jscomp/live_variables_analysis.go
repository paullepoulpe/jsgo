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
	lvaMaxLatticeBits                              = 128
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
	DataFlowAnalysisBase[rhino.Node, *LiveVariablesAnalysis_LiveVariableLattice]

	// The scope of the function that we are analyzing.
	jsScope *Scope

	// The scope of the body of the function that we are analyzing.
	jsScopeChild *Scope

	// escaped marks indices in orderedVars that are in escapedSet.
	escaped [2]uint64
	// untracked marks indices in orderedVars that are !local || escaped.
	untracked           [2]uint64
	escapedSet          jutil.ObjectLinkedHashSet[*Var]
	paramsMarkedEscaped bool

	// varIDs holds interned variable strpool IDs inline when numVars <= 16;
	// idToIndex is reused across Reset calls when numVars > 16.
	varIDs    [lvaInlineVarCap]uint32
	idToIndex map[uint32]uint8
	numVars   int

	// obtain variables in the order in which they appear in the code
	orderedVars []*Var

	allVarsInFn *jutil.LinkedHashMap[string, *Var]

	emptyLattice LiveVariablesAnalysis_LiveVariableLattice
	joinOp       LiveVariablesAnalysis_LiveVariableJoinOp
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
	j.result.liveSet[0] |= x.liveSet[0]
	j.result.liveSet[1] |= x.liveSet[1]
}

// Java: LiveVariablesAnalysis.LiveVariableJoinOp.finish()
func (j *LiveVariablesAnalysis_LiveVariableJoinOp) Finish() *LiveVariablesAnalysis_LiveVariableLattice {
	return j.result
}

// LiveVariablesAnalysis_LiveVariableLattice is the lattice that stores the
// liveness of all local variables at a given point in the program. The whole
// lattice is the power set of all local variables and a variable is live if
// it is in the set.
//
// Go-only: liveSet is a pointer-free [2]uint64 value so ReuseSlab chunks are
// 100% noscan for the Go GC.
type LiveVariablesAnalysis_LiveVariableLattice struct {
	liveSet [2]uint64
}

// numVars is the number of all local variables.
//
// Java: LiveVariablesAnalysis.LiveVariableLattice(int)
func NewLiveVariablesAnalysis_LiveVariableLattice(numVars int) *LiveVariablesAnalysis_LiveVariableLattice {
	if numVars < 0 {
		panic(jutil.NewNegativeArraySizeException("nbits < 0: " + strconv.Itoa(numVars)))
	}
	jutil.CheckStateMsg1(numVars <= lvaMaxLatticeBits, "LiveVariablesAnalysis supports at most 128 variables, got %d", numVars)
	return &LiveVariablesAnalysis_LiveVariableLattice{}
}

// Java: LiveVariablesAnalysis.LiveVariableLattice(LiveVariableLattice)
func NewLiveVariablesAnalysis_LiveVariableLatticeOf(other *LiveVariablesAnalysis_LiveVariableLattice) *LiveVariablesAnalysis_LiveVariableLattice {
	jutil.CheckNotNull(other)
	return &LiveVariablesAnalysis_LiveVariableLattice{liveSet: other.liveSet}
}

// Java: LiveVariablesAnalysis.LiveVariableLattice.equals(Object)
func (l *LiveVariablesAnalysis_LiveVariableLattice) Equals(other any) bool {
	jutil.CheckNotNull(other)
	liveVariableLattice, ok := other.(*LiveVariablesAnalysis_LiveVariableLattice)
	return ok && (l == liveVariableLattice || l.liveSet == liveVariableLattice.liveSet)
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
	if index < 64 {
		return (l.liveSet[0] & (uint64(1) << uint(index))) != 0
	}
	if index < lvaMaxLatticeBits {
		return (l.liveSet[1] & (uint64(1) << uint(index-64))) != 0
	}
	return false
}

// Java: LiveVariablesAnalysis.LiveVariableLattice.toString()
func (l *LiveVariablesAnalysis_LiveVariableLattice) String() string {
	w0, w1 := l.liveSet[0], l.liveSet[1]
	if w0 == 0 && w1 == 0 {
		return "{}"
	}
	buf := make([]byte, 0, 32)
	buf = append(buf, '{')
	first := true
	for w := w0; w != 0; w &= w - 1 {
		if !first {
			buf = append(buf, ',', ' ')
		}
		first = false
		buf = strconv.AppendInt(buf, int64(bits.TrailingZeros64(w)), 10)
	}
	for w := w1; w != 0; w &= w - 1 {
		if !first {
			buf = append(buf, ',', ' ')
		}
		first = false
		buf = strconv.AppendInt(buf, int64(64+bits.TrailingZeros64(w)), 10)
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
	if fromIndex < 64 {
		if w0 := l.liveSet[0] & (^uint64(0) << uint(fromIndex)); w0 != 0 {
			return bits.TrailingZeros64(w0)
		}
		if w1 := l.liveSet[1]; w1 != 0 {
			return 64 + bits.TrailingZeros64(w1)
		}
		return -1
	}
	if fromIndex < lvaMaxLatticeBits {
		if w1 := l.liveSet[1] & (^uint64(0) << uint(fromIndex-64)); w1 != 0 {
			return 64 + bits.TrailingZeros64(w1)
		}
		return -1
	}
	return -1
}

// Java: LiveVariablesAnalysis.LiveVariableLattice.hashCode()
func (l *LiveVariablesAnalysis_LiveVariableLattice) HashCode() int32 {
	var h int64 = 1234
	if l.liveSet[1] != 0 {
		h ^= int64(l.liveSet[1]) * 2
		h ^= int64(l.liveSet[0])
	} else if l.liveSet[0] != 0 {
		h ^= int64(l.liveSet[0])
	}
	return int32((h >> 32) ^ h)
}

func lvaSetBit128(words *[2]uint64, idx int) {
	if idx < 64 {
		words[0] |= uint64(1) << uint(idx)
	} else {
		words[1] |= uint64(1) << uint(idx-64)
	}
}

func lvaGetBit128(words [2]uint64, idx int) bool {
	if idx < 64 {
		return (words[0] & (uint64(1) << uint(idx))) != 0
	}
	return (words[1] & (uint64(1) << uint(idx-64))) != 0
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
	l.InitDataFlowAnalysisBase(l, cfg)
	jutil.CheckStateMsg1(jsScope.IsFunctionScope(), "%s", jsScope)

	l.jsScope = jsScope
	l.jsScopeChild = jsScopeChild
	l.escaped = [2]uint64{}
	l.untracked = [2]uint64{}
	l.paramsMarkedEscaped = false
	l.escapedSet.Clear()
	l.orderedVars = allVarsDeclaredInFunction.GetAllVariablesInOrder()
	l.numVars = len(l.orderedVars)
	jutil.CheckStateMsg1(l.numVars <= lvaMaxLatticeBits, "LiveVariablesAnalysis supports at most 128 variables, got %d", l.numVars)
	l.allVarsInFn = allVarsDeclaredInFunction.GetAllVariables()
	l.emptyLattice = LiveVariablesAnalysis_LiveVariableLattice{}
	l.lattices.Reset(nil)

	DataFlowAnalysis_ComputeEscaped(jsScope, &l.escapedSet, compiler, scopeCreator, l.allVarsInFn)

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
	f.finder.result.mapStorage.Clear()
	f.finder.result.allVarsInFn = &f.finder.result.mapStorage
	clear(f.finder.result.orderedVars)
	f.finder.result.orderedVars = f.finder.result.orderedVars[:0]
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
		jutil.CheckState(f.finder.result.allVarsInFn.IsEmpty() == (len(f.finder.result.orderedVars) == 0))
		return &f.finder.result
	}

	b := NodeTraversal_Builder{
		compiler:     compiler,
		callback:     f,
		scopeCreator: scopeCreator,
	}
	b.TraverseAtScope(scope)

	if maxVars > 0 && len(f.finder.result.orderedVars) >= maxVars {
		return nil
	}
	jutil.CheckState(f.finder.result.allVarsInFn.IsEmpty() == (len(f.finder.result.orderedVars) == 0))
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
	if !useInline {
		if l.idToIndex == nil {
			l.idToIndex = make(map[uint32]uint8, l.numVars)
		} else {
			clear(l.idToIndex)
		}
	}
	for num, v := range l.orderedVars {
		id := v.nameID
		if useInline {
			l.varIDs[num] = id
		} else {
			l.idToIndex[id] = uint8(num)
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
		isEsc := l.escapedSet.Contains(v)
		if isEsc {
			lvaSetBit128(&l.escaped, num)
		}
		if !local || isEsc {
			lvaSetBit128(&l.untracked, num)
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
	return int(idx), ok
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
	return lvaGetBit128(l.escaped, idx)
}

func (l *LiveVariablesAnalysis) allInFnVarsEscaped() bool {
	total := l.allVarsInFn.Size()
	if l.escapedSet.Size() < total {
		return false
	}
	if l.numVars == total && l.numVars > 0 {
		var mask0, mask1 uint64
		if l.numVars < 64 {
			mask0 = (uint64(1) << uint(l.numVars)) - 1
		} else if l.numVars == 64 {
			mask0 = ^uint64(0)
		} else if l.numVars < lvaMaxLatticeBits {
			mask0 = ^uint64(0)
			mask1 = (uint64(1) << uint(l.numVars-64)) - 1
		} else {
			mask0 = ^uint64(0)
			mask1 = ^uint64(0)
		}
		return (l.escaped[0]&mask0) == mask0 && (l.escaped[1]&mask1) == mask1
	}
	for v := range l.allVarsInFn.Values() {
		if !l.escapedSet.Contains(v) {
			return false
		}
	}
	return true
}

// Java: LiveVariablesAnalysis.getEscapedLocals()
func (l *LiveVariablesAnalysis) GetEscapedLocals() *jutil.ObjectLinkedHashSet[*Var] {
	return &l.escapedSet
}

// Java: LiveVariablesAnalysis.getAllVariables()
func (l *LiveVariablesAnalysis) GetAllVariables() *jutil.LinkedHashMap[string, *Var] {
	return l.allVarsInFn
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

// Java: LiveVariablesAnalysis.isForward()
func (l *LiveVariablesAnalysis) IsForward() bool {
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
	return a.liveSet == b.liveSet
}

// Java: LiveVariablesAnalysis.createFlowJoiner()
func (l *LiveVariablesAnalysis) CreateFlowJoiner() DataFlowAnalysis_FlowJoiner[*LiveVariablesAnalysis_LiveVariableLattice] {
	l.joinOp.result = l.allocLattice()
	return &l.joinOp
}

// Java: LiveVariablesAnalysis.flowThrough(Node, LiveVariableLattice)
func (l *LiveVariablesAnalysis) FlowThrough(node rhino.Node, input *LiveVariablesAnalysis_LiveVariableLattice) *LiveVariablesAnalysis_LiveVariableLattice {
	var gen, kill [2]uint64

	// Make kills conditional if the node can end abruptly by an exception.
	conditional := false
	var edgeList []*graph.LinkedDiGraphEdge[rhino.Node, ControlFlowGraph_Branch]
	if l.curLinkedNode != nil && l.curLinkedNode.Value == node {
		edgeList = l.curLinkedNode.OutEdgeList.Slice()
	} else {
		edgeList = l.GetCfg().GetOutEdges(node)
	}
	for _, edge := range edgeList {
		if ControlFlowGraph_Branch_ON_EX == edge.GetValue() {
			conditional = true
			break
		}
	}
	l.computeGenKill128(node, &gen, &kill, conditional)
	// L_in = L_out - Kill + Gen
	out0 := (input.liveSet[0] &^ kill[0]) | gen[0]
	out1 := (input.liveSet[1] &^ kill[1]) | gen[1]
	if out0 == input.liveSet[0] && out1 == input.liveSet[1] {
		return input
	}
	result := l.allocLattice()
	result.liveSet[0] = out0
	result.liveSet[1] = out1
	return result
}

func lvaOrIntoBitSet(dst *jutil.BitSet, src [2]uint64) {
	for w := src[0]; w != 0; w &= w - 1 {
		dst.Set(bits.TrailingZeros64(w))
	}
	for w := src[1]; w != 0; w &= w - 1 {
		dst.Set(64 + bits.TrailingZeros64(w))
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
	var gen128, kill128 [2]uint64
	l.computeGenKill128(n, &gen128, &kill128, conditional)
	lvaOrIntoBitSet(gen, gen128)
	lvaOrIntoBitSet(kill, kill128)
}

func (l *LiveVariablesAnalysis) computeGenKill128(n rhino.Node, gen *[2]uint64, kill *[2]uint64, conditional bool) {
	switch n.GetToken() {
	case rhino.Token_SCRIPT, rhino.Token_ROOT, rhino.Token_FUNCTION, rhino.Token_BLOCK:
		return
	case rhino.Token_WHILE, rhino.Token_DO, rhino.Token_IF, rhino.Token_FOR:
		l.computeGenKill128(NodeUtil_GetConditionExpression(n), gen, kill, conditional)
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
			l.computeGenKill128(lhs, gen, kill, conditional)

			// rhs is executed only once so we don't go into it every loop.
			return
		}
	case rhino.Token_LET, rhino.Token_CONST, rhino.Token_VAR:
		decl := ast.MustNameDeclaration(n)
		for d, ok := decl.FirstDeclarator(); ok; d, ok = d.Next() {
			if name, isName := d.AsName(); isName {
				if init, hasInit := name.Init(); hasInit {
					l.computeGenKill128(init.Node(), gen, kill, conditional)
					if !conditional {
						l.addToSetIfLocal128(name.Node(), kill)
					}
				}
			} else {
				dlhs := d.MustDestructuringLhs()
				if !conditional {
					NodeUtil_VisitLhsNodesInNode(dlhs.Node(), func(lhsNode rhino.Node) { l.addToSetIfLocal128(lhsNode, kill) })
				}
				l.computeGenKill128(dlhs.Pattern().Node(), gen, kill, conditional)
				if init, hasInit := dlhs.Init(); hasInit {
					l.computeGenKill128(init.Node(), gen, kill, conditional)
				}
			}
		}
		return
	case rhino.Token_OPTCHAIN_GETPROP:
		// Go-only: OPTCHAIN_GETPROP has only one child (Object; property is its string),
		// so visiting Object once matches Java without O(2^k) re-traversal on chains.
		l.computeGenKill128(ast.MustOptChainGetProp(n).Object().Node(), gen, kill, conditional)
		return
	case rhino.Token_AND, rhino.Token_OR, rhino.Token_COALESCE:
		bin := ast.MustBinary(n)
		l.computeGenKill128(bin.Left().Node(), gen, kill, conditional)
		// May short circuit.
		l.computeGenKill128(bin.Right().Node(), gen, kill, true)
		return
	case rhino.Token_OPTCHAIN_GETELEM:
		ge := ast.MustOptChainGetElem(n)
		l.computeGenKill128(ge.Object().Node(), gen, kill, conditional)
		// May short circuit.
		l.computeGenKill128(ge.Index().Node(), gen, kill, true)
		return
	case rhino.Token_OPTCHAIN_CALL:
		call := ast.MustOptChainCall(n)
		l.computeGenKill128(call.Callee().Node(), gen, kill, conditional)
		// Unlike OPTCHAIN_GETPROP and OPTCHAIN_GETELEM, the OPTCHAIN_CALLs can have multiple
		// children on rhs which get executed conditionally
		for arg, ok := call.FirstArg(); ok; arg, ok = arg.Next() {
			l.computeGenKill128(arg.Node(), gen, kill, true)
		}
		return
	case rhino.Token_HOOK:
		hook := ast.MustHook(n)
		l.computeGenKill128(hook.Cond().Node(), gen, kill, conditional)
		// Assume both sides are conditional.
		l.computeGenKill128(hook.Then().Node(), gen, kill, true)
		l.computeGenKill128(hook.Else().Node(), gen, kill, true)
		return
	case rhino.Token_NAME:
		if n.GetStringID() == lvaArgumentsIdentID {
			l.MarkAllParametersEscaped()
		} else if !NodeUtil_IsLhsByDestructuring(n) {
			// Only add names in destructuring patterns if they're not lvalues.
			// e.g. "x" in "const {foo = x} = obj;"
			l.addToSetIfLocal128(n, gen)
		}
		return
	default:
		if assign, ok := ast.AsAssignment(n); ok {
			if lhs, isName := assign.Target().AsName(); isName {
				if !conditional {
					l.addToSetIfLocal128(lhs.Node(), kill)
				}
				if assign.Token() != rhino.Token_ASSIGN {
					// assignments such as a += 1 reads a.
					l.addToSetIfLocal128(lhs.Node(), gen)
				}
				l.computeGenKill128(assign.Value().Node(), gen, kill, conditional)
				return
			}
			if assign.Token() == rhino.Token_ASSIGN {
				if pat, isPat := assign.Target().AsPattern(); isPat {
					if !conditional {
						NodeUtil_VisitLhsNodesInNode(
							n,
							func(child rhino.Node) {
								if child.IsName() {
									l.addToSetIfLocal128(child, kill)
								}
							})
					}
					l.computeGenKill128(pat.Node(), gen, kill, conditional)
					l.computeGenKill128(assign.Value().Node(), gen, kill, conditional)
					return
				}
			}
		}
		for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			l.computeGenKill128(c, gen, kill, conditional)
		}
		return
	}
}

func (l *LiveVariablesAnalysis) addToSetIfLocal128(node rhino.Node, set *[2]uint64) {
	jutil.CheckStateMsg1(node.IsName(), "%s", node)
	idx, found := l.lookupVarIndexByID(node.GetStringID())
	if !found || lvaGetBit128(l.untracked, idx) {
		return
	}
	lvaSetBit128(set, idx)
}

// Java: LiveVariablesAnalysis.addToSetIfLocal(Node, BitSet)
func (l *LiveVariablesAnalysis) AddToSetIfLocal(node rhino.Node, set *jutil.BitSet) {
	jutil.CheckStateMsg1(node.IsName(), "%s", node)
	idx, found := l.lookupVarIndexByID(node.GetStringID())
	if !found || lvaGetBit128(l.untracked, idx) {
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
			l.escapedSet.Add(l.jsScope.GetVarByID(id))
			if idx, ok := l.lookupVarIndexByID(id); ok {
				lvaSetBit128(&l.escaped, idx)
				lvaSetBit128(&l.untracked, idx)
			}
		}
	}
}
