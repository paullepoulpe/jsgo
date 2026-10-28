// Port of com/google/javascript/jscomp/MaybeReachingVariableUse.java @ closure-compiler v20261021

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
	"math/bits"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/graph"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// MaybeReachingVariableUse computes "may be" reaching use for all
// definitions of each variable.
//
// A use of A in `alert(A)` is a "may be" reaching use of the definition of A
// at `A = foo()` if at least one path from the definition node to the end
// node reaches that use and it is the last definition before the use on
// that path.
//
// Example:
//
//	D1: var A = foo();
//	U1: alert(A);
//	    if(....) {
//	       D2: A = bar();
//	       U2: alert(A);
//	    }
//	U3: alert(A);
//
// Here, MaybeReachingUses[D1] = {U1, U3} and MaybeReachingUses[D2]={U2, U3}.
// The use U3 is not guaranteed to use def D1: this is a "may-be" analysis.
//
// Go-only: performance — instead of re-walking each CFG node's AST and
// rebuilding persistent HAMT maps on every fixed-point step, Analyze() walks
// each CFG block once to compute local upward-exposed (varIdx, cfgNode) uses
// and unconditional variable kills, assigns consecutive bit positions grouped
// by variable, and solves the backward fixed-point over flat uint64 bitvectors
// (single uint64 per node when total uses <= 64, []uint64 when > 64).
type MaybeReachingVariableUse struct {
	cfg         ControlFlowGraph[rhino.Node]
	cfgBase     *ControlFlowGraphBase[rhino.Node]
	hasPriority bool

	numVars       int
	varNameIDs    []uint32
	varNameIDsBuf [16]uint32
	varIdxByID    map[uint32]int // stores varIdx + 1
	useVarMap     bool

	// Per-block scratch during the single AST pre-pass.
	curGen64        uint64
	curKill64       uint64
	curGenVarWords  []uint64
	curKillVarWords []uint64

	blockKillVars64   []uint64
	blockKillVarWords []uint64
	rawUses           []mrvuRawUse
	varUseCount       []int
	varBitStart       []int
	varBitEnd         []int
	varWritePos       []int

	totalUses int
	useWords  int
	useNodes  []rhino.Node

	// Fast path (totalUses <= 64).
	varMask64 []uint64
	gen64     []uint64
	kill64    []uint64
	in64      []uint64
	out64     []uint64

	// General path (totalUses > 64).
	genWords  []uint64
	killWords []uint64
	inWords   []uint64
	outWords  []uint64

	// Worklist (max-heap on Priority when hasPriority, else FIFO).
	queue     []*graph.LinkedDirectedGraph_LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch]
	queueHead int
	steps     []int
}

type mrvuRawUse struct {
	varIdx  uint32
	cfgIdx  uint32
	cfgNode rhino.Node
}

const mrvuReuseMaxCap = 4096

// Java: MaybeReachingVariableUse(ControlFlowGraph<Node>, Set<Var>, Map<String, Var>)
func NewMaybeReachingVariableUse(
	cfg ControlFlowGraph[rhino.Node],
	escaped *jutil.ObjectLinkedHashSet[*Var],
	allVarsInFn *jutil.LinkedHashMap[string, *Var],
) *MaybeReachingVariableUse {
	u := &MaybeReachingVariableUse{}
	u.Reset(cfg, escaped, allVarsInFn)
	return u
}

// Reset resets u for reuse on a new function scope while retaining its
// storage.
func (u *MaybeReachingVariableUse) Reset(
	cfg ControlFlowGraph[rhino.Node],
	escaped *jutil.ObjectLinkedHashSet[*Var],
	allVarsInFn *jutil.LinkedHashMap[string, *Var],
) {
	u.releaseLargeStorage()
	u.cfg = cfg
	u.cfgBase = cfg.GetControlFlowGraphBase()
	u.hasPriority = cfg.GetOptionalNodeComparator(false) != nil
	u.totalUses = 0
	u.useWords = 0

	u.varNameIDs = u.varNameIDsBuf[:0]
	u.useVarMap = false
	if allVarsInFn != nil {
		for v := range allVarsInFn.Values() {
			if escaped == nil || !escaped.Contains(v) {
				id := v.nameID
				if id == 0 {
					id = rhino.RhinoStringPool_InternID(v.GetName())
				}
				u.varNameIDs = append(u.varNameIDs, id)
			}
		}
		if len(u.varNameIDs) > len(u.varNameIDsBuf) {
			u.useVarMap = true
			if u.varIdxByID == nil {
				u.varIdxByID = make(map[uint32]int, max(len(u.varNameIDs), 32))
			} else {
				clear(u.varIdxByID)
			}
			for i, id := range u.varNameIDs {
				u.varIdxByID[id] = i + 1
			}
		}
	}
	u.numVars = len(u.varNameIDs)
}

func (u *MaybeReachingVariableUse) releaseLargeStorage() {
	clear(u.useNodes)
	u.useNodes = u.useNodes[:0]
	if cap(u.useNodes) > mrvuReuseMaxCap {
		u.useNodes = nil
	}
	clear(u.rawUses)
	u.rawUses = u.rawUses[:0]
	if cap(u.rawUses) > mrvuReuseMaxCap {
		u.rawUses = nil
	}
	if cap(u.gen64) > mrvuReuseMaxCap {
		u.blockKillVars64 = nil
		u.gen64 = nil
		u.kill64 = nil
		u.in64 = nil
		u.out64 = nil
		u.steps = nil
	}
	if cap(u.genWords) > mrvuReuseMaxCap {
		u.blockKillVarWords = nil
		u.genWords = nil
		u.killWords = nil
		u.inWords = nil
		u.outWords = nil
	}
	if cap(u.queue) > mrvuReuseMaxCap {
		u.queue = nil
		u.queueHead = 0
	}
}

// GetCfg returns the control flow graph that this analysis was performed on.
//
// Java: DataFlowAnalysis.getCfg()
func (u *MaybeReachingVariableUse) GetCfg() ControlFlowGraph[rhino.Node] {
	return u.cfg
}

func (u *MaybeReachingVariableUse) lookupVarIdxByID(id uint32) int {
	if id == 0 {
		return -1
	}
	if !u.useVarMap {
		for i, vid := range u.varNameIDs {
			if vid == id {
				return i
			}
		}
		return -1
	}
	return u.varIdxByID[id] - 1
}

func (u *MaybeReachingVariableUse) lookupVarIdxByName(name string) int {
	return u.lookupVarIdxByID(rhino.RhinoStringPool_LookupID(name))
}

func mrvuResetZeroSlice[T any](s []T, n int) []T {
	if cap(s) < n {
		return make([]T, n)
	}
	clear(s)
	return s[:n]
}

func mrvuOrBitRange(dst []uint64, start, end int) {
	if start >= end {
		return
	}
	sw, sb := start>>6, uint(start&63)
	ew, eb := (end-1)>>6, uint((end-1)&63)
	if sw == ew {
		dst[sw] |= (^uint64(0) >> (63 - (eb - sb))) << sb
		return
	}
	dst[sw] |= ^uint64(0) << sb
	for w := sw + 1; w < ew; w++ {
		dst[w] = ^uint64(0)
	}
	dst[ew] |= ^uint64(0) >> (63 - eb)
}

// Analyze computes the fixed-point reaching uses across the CFG.
//
// Java: DataFlowAnalysis.analyze()
func (u *MaybeReachingVariableUse) Analyze() {
	numVars := u.numVars
	u.totalUses = 0
	if numVars == 0 {
		return
	}

	cfgBase := u.cfgBase
	implicitReturn, _ := cfgBase.implicitReturn.(*graph.LinkedDirectedGraph_LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch])
	nodeCap := cfgBase.NodeIndexCap()

	clear(u.rawUses)
	u.rawUses = u.rawUses[:0]
	u.varUseCount = mrvuResetZeroSlice(u.varUseCount, numVars)

	varWords := (numVars + 63) >> 6
	if numVars <= 64 {
		u.blockKillVars64 = mrvuResetZeroSlice(u.blockKillVars64, nodeCap)
		for node := range cfgBase.Nodes.Values() {
			if node == implicitReturn {
				continue
			}
			conditional := false
			for _, edge := range node.OutEdgeList.Slice() {
				if edge.Value == ControlFlowGraph_Branch_ON_EX {
					conditional = true
					break
				}
			}
			u.curGen64 = 0
			u.curKill64 = 0
			u.computeMayUse(node.Value, conditional)
			u.blockKillVars64[node.Index] = u.curKill64
			for m := u.curGen64; m != 0; m &= m - 1 {
				v := bits.TrailingZeros64(m)
				u.varUseCount[v]++
				u.rawUses = append(u.rawUses, mrvuRawUse{
					varIdx:  uint32(v),
					cfgIdx:  node.Index,
					cfgNode: node.Value,
				})
			}
		}
	} else {
		u.blockKillVarWords = mrvuResetZeroSlice(u.blockKillVarWords, nodeCap*varWords)
		u.curGenVarWords = mrvuResetZeroSlice(u.curGenVarWords, varWords)
		u.curKillVarWords = mrvuResetZeroSlice(u.curKillVarWords, varWords)
		for node := range cfgBase.Nodes.Values() {
			if node == implicitReturn {
				continue
			}
			conditional := false
			for _, edge := range node.OutEdgeList.Slice() {
				if edge.Value == ControlFlowGraph_Branch_ON_EX {
					conditional = true
					break
				}
			}
			clear(u.curGenVarWords)
			clear(u.curKillVarWords)
			u.computeMayUse(node.Value, conditional)
			base := int(node.Index) * varWords
			copy(u.blockKillVarWords[base:base+varWords], u.curKillVarWords)
			for wIdx, w := range u.curGenVarWords {
				for m := w; m != 0; m &= m - 1 {
					v := (wIdx << 6) + bits.TrailingZeros64(m)
					u.varUseCount[v]++
					u.rawUses = append(u.rawUses, mrvuRawUse{
						varIdx:  uint32(v),
						cfgIdx:  node.Index,
						cfgNode: node.Value,
					})
				}
			}
		}
	}

	totalUses := len(u.rawUses)
	u.totalUses = totalUses
	if totalUses == 0 {
		return
	}

	u.varBitStart = mrvuResetZeroSlice(u.varBitStart, numVars)
	u.varBitEnd = mrvuResetZeroSlice(u.varBitEnd, numVars)
	u.varWritePos = mrvuResetZeroSlice(u.varWritePos, numVars)
	offset := 0
	for v := range numVars {
		u.varBitStart[v] = offset
		u.varWritePos[v] = offset
		offset += u.varUseCount[v]
		u.varBitEnd[v] = offset
	}

	u.useNodes = mrvuResetZeroSlice(u.useNodes, totalUses)
	u.steps = mrvuResetZeroSlice(u.steps, nodeCap)

	if totalUses <= 64 {
		u.varMask64 = mrvuResetZeroSlice(u.varMask64, numVars)
		for v := range numVars {
			s, e := u.varBitStart[v], u.varBitEnd[v]
			if s < e {
				u.varMask64[v] = (^uint64(0) >> uint(64-(e-s))) << uint(s)
			}
		}
		u.gen64 = mrvuResetZeroSlice(u.gen64, nodeCap)
		u.kill64 = mrvuResetZeroSlice(u.kill64, nodeCap)
		u.in64 = mrvuResetZeroSlice(u.in64, nodeCap)
		u.out64 = mrvuResetZeroSlice(u.out64, nodeCap)

		for _, ru := range u.rawUses {
			bit := u.varWritePos[ru.varIdx]
			u.varWritePos[ru.varIdx]++
			u.useNodes[bit] = ru.cfgNode
			u.gen64[ru.cfgIdx] |= uint64(1) << uint(bit)
		}
		if numVars <= 64 {
			for node := range cfgBase.Nodes.Values() {
				if node == implicitReturn {
					continue
				}
				var k uint64
				for m := u.blockKillVars64[node.Index]; m != 0; m &= m - 1 {
					v := bits.TrailingZeros64(m)
					k |= u.varMask64[v]
				}
				u.kill64[node.Index] = k
			}
		} else {
			for node := range cfgBase.Nodes.Values() {
				if node == implicitReturn {
					continue
				}
				var k uint64
				base := int(node.Index) * varWords
				for wIdx, w := range u.blockKillVarWords[base : base+varWords] {
					for m := w; m != 0; m &= m - 1 {
						v := (wIdx << 6) + bits.TrailingZeros64(m)
						k |= u.varMask64[v]
					}
				}
				u.kill64[node.Index] = k
			}
		}
		u.solve64(cfgBase, implicitReturn)
		return
	}

	words := (totalUses + 63) >> 6
	u.useWords = words
	totalWords := nodeCap * words
	u.genWords = mrvuResetZeroSlice(u.genWords, totalWords)
	u.killWords = mrvuResetZeroSlice(u.killWords, totalWords)
	u.inWords = mrvuResetZeroSlice(u.inWords, totalWords)
	u.outWords = mrvuResetZeroSlice(u.outWords, totalWords)

	for _, ru := range u.rawUses {
		bit := u.varWritePos[ru.varIdx]
		u.varWritePos[ru.varIdx]++
		u.useNodes[bit] = ru.cfgNode
		u.genWords[int(ru.cfgIdx)*words+(bit>>6)] |= uint64(1) << uint(bit&63)
	}
	if numVars <= 64 {
		for node := range cfgBase.Nodes.Values() {
			if node == implicitReturn {
				continue
			}
			killRow := u.killWords[int(node.Index)*words : (int(node.Index)+1)*words]
			for m := u.blockKillVars64[node.Index]; m != 0; m &= m - 1 {
				v := bits.TrailingZeros64(m)
				mrvuOrBitRange(killRow, u.varBitStart[v], u.varBitEnd[v])
			}
		}
	} else {
		for node := range cfgBase.Nodes.Values() {
			if node == implicitReturn {
				continue
			}
			killRow := u.killWords[int(node.Index)*words : (int(node.Index)+1)*words]
			base := int(node.Index) * varWords
			for wIdx, w := range u.blockKillVarWords[base : base+varWords] {
				for m := w; m != 0; m &= m - 1 {
					v := (wIdx << 6) + bits.TrailingZeros64(m)
					mrvuOrBitRange(killRow, u.varBitStart[v], u.varBitEnd[v])
				}
			}
		}
	}
	u.solveWords(cfgBase, implicitReturn, words)
}

func (u *MaybeReachingVariableUse) initQueue(
	cfgBase *ControlFlowGraphBase[rhino.Node],
	implicitReturn *graph.LinkedDirectedGraph_LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch],
) {
	clear(u.queue)
	u.queue = u.queue[:0]
	u.queueHead = 0
	numNodes := cfgBase.Nodes.Size()
	if cap(u.queue) < numNodes {
		u.queue = make([]*graph.LinkedDirectedGraph_LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch], 0, numNodes)
	}
	for node := range cfgBase.Nodes.Values() {
		node.DfaInQueue = false
		if node != implicitReturn {
			u.enqueue(node)
		}
	}
}

func (u *MaybeReachingVariableUse) enqueue(node *graph.LinkedDirectedGraph_LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch]) {
	node.DfaInQueue = true
	u.queue = append(u.queue, node)
	if u.hasPriority {
		i := len(u.queue) - 1
		xPrio := node.Priority
		if xPrio < 0 {
			panic(jutil.NewIllegalStateExceptionMsg("priority not set"))
		}
		for i > 0 {
			parent := (i - 1) / 2
			p := u.queue[parent]
			if xPrio <= p.Priority {
				break
			}
			u.queue[i] = p
			i = parent
		}
		u.queue[i] = node
	}
}

func (u *MaybeReachingVariableUse) dequeue() *graph.LinkedDirectedGraph_LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch] {
	var cur *graph.LinkedDirectedGraph_LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch]
	if u.hasPriority {
		n := len(u.queue)
		cur = u.queue[0]
		last := u.queue[n-1]
		u.queue[n-1] = nil
		u.queue = u.queue[:n-1]
		if n > 1 {
			i := 0
			m := n - 1
			xPrio := last.Priority
			for {
				pick := 2*i + 1
				if pick >= m {
					break
				}
				larger := u.queue[pick]
				largerPrio := larger.Priority
				if sibling := pick + 1; sibling < m {
					sib := u.queue[sibling]
					if sibPrio := sib.Priority; largerPrio < sibPrio {
						pick = sibling
						larger = sib
						largerPrio = sibPrio
					}
				}
				if xPrio >= largerPrio {
					break
				}
				u.queue[i] = larger
				i = pick
			}
			u.queue[i] = last
		}
	} else {
		cur = u.queue[u.queueHead]
		u.queue[u.queueHead] = nil
		u.queueHead++
		if u.queueHead == len(u.queue) {
			u.queue = u.queue[:0]
			u.queueHead = 0
		}
	}
	cur.DfaInQueue = false
	return cur
}

func (u *MaybeReachingVariableUse) solve64(
	cfgBase *ControlFlowGraphBase[rhino.Node],
	implicitReturn *graph.LinkedDirectedGraph_LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch],
) {
	u.initQueue(cfgBase, implicitReturn)
	for len(u.queue) > u.queueHead {
		curNode := u.dequeue()
		idx := curNode.Index
		stepCount := u.steps[idx]
		u.steps[idx] = stepCount + 1
		if stepCount > DataFlowAnalysis_MAX_STEPS_PER_NODE {
			panic(jutil.NewIllegalStateExceptionMsg("Dataflow analysis appears to diverge around: " + jutil.ToJavaString(curNode)))
		}

		outEdges := curNode.OutEdgeList.Slice()
		if len(outEdges) > 0 {
			var joined uint64
			for _, edge := range outEdges {
				dest := edge.DestNode
				if dest != implicitReturn {
					joined |= u.in64[dest.Index]
				}
			}
			u.out64[idx] = joined
		}

		inAfter := (u.out64[idx] &^ u.kill64[idx]) | u.gen64[idx]
		if inAfter != u.in64[idx] {
			u.in64[idx] = inAfter
			for _, edge := range curNode.InEdgeList.Slice() {
				prev := edge.SourceNode
				if !prev.DfaInQueue && prev != implicitReturn {
					u.enqueue(prev)
				}
			}
		}
	}
}

func (u *MaybeReachingVariableUse) solveWords(
	cfgBase *ControlFlowGraphBase[rhino.Node],
	implicitReturn *graph.LinkedDirectedGraph_LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch],
	words int,
) {
	u.initQueue(cfgBase, implicitReturn)
	for len(u.queue) > u.queueHead {
		curNode := u.dequeue()
		idx := int(curNode.Index)
		stepCount := u.steps[idx]
		u.steps[idx] = stepCount + 1
		if stepCount > DataFlowAnalysis_MAX_STEPS_PER_NODE {
			panic(jutil.NewIllegalStateExceptionMsg("Dataflow analysis appears to diverge around: " + jutil.ToJavaString(curNode)))
		}

		base := idx * words
		outRow := u.outWords[base : base+words]
		outEdges := curNode.OutEdgeList.Slice()
		if len(outEdges) > 0 {
			clear(outRow)
			for _, edge := range outEdges {
				dest := edge.DestNode
				if dest != implicitReturn {
					destBase := int(dest.Index) * words
					destIn := u.inWords[destBase : destBase+words]
					for w := range words {
						outRow[w] |= destIn[w]
					}
				}
			}
		}

		inRow := u.inWords[base : base+words]
		killRow := u.killWords[base : base+words]
		genRow := u.genWords[base : base+words]
		changed := false
		for w := range words {
			after := (outRow[w] &^ killRow[w]) | genRow[w]
			if after != inRow[w] {
				inRow[w] = after
				changed = true
			}
		}
		if changed {
			for _, edge := range curNode.InEdgeList.Slice() {
				prev := edge.SourceNode
				if !prev.DfaInQueue && prev != implicitReturn {
					u.enqueue(prev)
				}
			}
		}
	}
}

func (u *MaybeReachingVariableUse) genVarIdx(v int) {
	if u.numVars <= 64 {
		u.curGen64 |= uint64(1) << uint(v)
	} else {
		u.curGenVarWords[v>>6] |= uint64(1) << uint(v&63)
	}
}

func (u *MaybeReachingVariableUse) killVarIdx(v int) {
	if u.numVars <= 64 {
		bit := uint64(1) << uint(v)
		u.curKill64 |= bit
		u.curGen64 &^= bit
	} else {
		w, b := v>>6, uint(v&63)
		bit := uint64(1) << b
		u.curKillVarWords[w] |= bit
		u.curGenVarWords[w] &^= bit
	}
}

func (u *MaybeReachingVariableUse) genLocalID(id uint32) {
	if v := u.lookupVarIdxByID(id); v >= 0 {
		u.genVarIdx(v)
	}
}

func (u *MaybeReachingVariableUse) killLocalID(id uint32) {
	if v := u.lookupVarIdxByID(id); v >= 0 {
		u.killVarIdx(v)
	}
}

// computeMayUse finds all upward-exposed uses and unconditional kills of
// unescaped local variables within a single CFG node, walking subexpressions
// in reverse evaluation order.
//
// Java: MaybeReachingVariableUse.computeMayUse(Node, Node, ReachingUses, boolean)
func (u *MaybeReachingVariableUse) computeMayUse(n rhino.Node, conditional bool) {
	if n.IsNil() {
		return
	}
	switch n.GetToken() {
	case rhino.Token_BLOCK, rhino.Token_ROOT, rhino.Token_FUNCTION:
		return
	case rhino.Token_NAME:
		id := n.GetStringID()
		if NodeUtil_IsLhsByDestructuring(n) {
			if !conditional {
				u.killLocalID(id)
			}
		} else {
			u.genLocalID(id)
		}
		return
	case rhino.Token_WHILE, rhino.Token_DO, rhino.Token_IF, rhino.Token_FOR:
		condExpr := NodeUtil_GetConditionExpression(n)
		u.computeMayUse(condExpr, conditional)
		return
	case rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
		// for(x in y) {...}
		loop := ast.MustEnhancedFor(n)
		target := loop.Target()
		var lhs rhino.Node
		if decl, ok := target.AsNameDeclaration(); ok {
			if last, hasLast := decl.LastDeclarator(); hasLast { // for(var x in y) {...}
				if dlhs, isDlhs := last.AsDestructuringLhs(); isDlhs {
					lhs = dlhs.Pattern().Node() // for (let [x] of obj) {...}
				} else {
					lhs = last.Node()
				}
			}
		} else {
			lhs = target.Node()
		}
		if _, isPat := ast.AsPattern(lhs); isPat {
			u.computeMayUse(lhs, true)
		}
		u.computeMayUse(loop.Collection().Node(), conditional)
		return
	case rhino.Token_OPTCHAIN_GETPROP:
		// Go-only: OPTCHAIN_GETPROP has only one child (Object; property is its string),
		// so visiting Object once matches Java without O(2^k) re-traversal on chains.
		u.computeMayUse(ast.MustOptChainGetProp(n).Object().Node(), conditional)
		return
	case rhino.Token_AND, rhino.Token_OR, rhino.Token_COALESCE:
		bin := ast.MustBinary(n)
		u.computeMayUse(bin.Right().Node() /* conditional= */, true)
		u.computeMayUse(bin.Left().Node(), conditional)
		return
	case rhino.Token_OPTCHAIN_GETELEM:
		ge := ast.MustOptChainGetElem(n)
		u.computeMayUse(ge.Index().Node() /* conditional= */, true)
		u.computeMayUse(ge.Object().Node(), conditional)
		return
	case rhino.Token_OPTCHAIN_CALL:
		call := ast.MustOptChainCall(n)
		callee := call.Callee().Node()
		// As args are evaluated in AST order, we traverse in reverse AST order for backward
		// dataflow analysis.
		for c := n.GetLastChild(); c != callee; c = c.GetPrevious() {
			u.computeMayUse(c /* conditional= */, true)
		}
		u.computeMayUse(callee, conditional)
		return
	case rhino.Token_HOOK:
		hook := ast.MustHook(n)
		u.computeMayUse(hook.Else().Node() /* conditional= */, true)
		u.computeMayUse(hook.Then().Node() /* conditional= */, true)
		u.computeMayUse(hook.Cond().Node(), conditional)
		return
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		decl := ast.MustNameDeclaration(n)
		jutil.CheckStateMsg1(decl.HasDeclarators(), "AST should be normalized (%s)", n)
		d := decl.MustFirstDeclarator()
		if dlhs, ok := d.AsDestructuringLhs(); ok {
			// Note: since destructuring is evaluated in reverse AST order, we traverse the pattern
			// before the initializer in order to do our backwards data flow analysis.
			u.computeMayUse(dlhs.Pattern().Node(), conditional)
			if init, hasInit := dlhs.Init(); hasInit {
				u.computeMayUse(init.Node(), conditional)
			}
		} else {
			varName := d.MustName()
			if init, hasInit := varName.Init(); hasInit {
				u.computeMayUse(init.Node(), conditional)
				if !conditional {
					u.killLocalID(varName.Node().GetStringID())
				}
			} // else var name declaration with no initial value
		}
		return
	case rhino.Token_DEFAULT_VALUE:
		dv := ast.MustDefaultValue(n)
		target := dv.Target()
		value := dv.Value().Node()
		if pat, ok := target.AsPattern(); ok {
			u.computeMayUse(pat.Node(), conditional)
			u.computeMayUse(value, true)
		} else if name, ok := target.AsName(); ok {
			// assigning to the name occurs after evaluating the default value
			if !conditional {
				u.killLocalID(name.Node().GetStringID())
			}
			u.computeMayUse(value, true)
		} else {
			u.computeMayUse(value, true)
			u.computeMayUse(target.Node(), conditional)
		}
	default:
		if assign, ok := ast.AsAssignment(n); ok {
			if name, isName := assign.Target().AsName(); isName {
				jutil.CheckState(!NodeUtil_IsLogicalAssignmentOp(n))
				if v := u.lookupVarIdxByID(name.Node().GetStringID()); v >= 0 {
					if !conditional {
						u.killVarIdx(v)
					}
					// In case of a += "Hello". There is a read of a.
					if assign.Token() != rhino.Token_ASSIGN {
						u.genVarIdx(v)
					}
				}
				u.computeMayUse(assign.Value().Node(), conditional)
				return
			}
			if assign.Token() == rhino.Token_ASSIGN {
				if pat, isPat := assign.Target().AsPattern(); isPat {
					// Note: the rhs of destructuring is evaluated before the lhs
					u.computeMayUse(pat.Node(), conditional)
					u.computeMayUse(assign.Value().Node(), conditional)
					return
				}
			}
		}
		/*
		 * We want to traverse in reverse order because we want the LAST
		 * definition in the sub-tree.
		 */
		for c := n.GetLastChild(); !c.IsNil(); c = c.GetPrevious() {
			u.computeMayUse(c, conditional)
		}
	}
}

func mrvuEmptyUses(func(rhino.Node) bool) {}

// HasExactlyOneUse reports whether defNode has exactly one upward-exposed
// reaching use CFG node for name, without allocating an iterator closure.
func (u *MaybeReachingVariableUse) HasExactlyOneUse(name string, defNode rhino.Node) bool {
	n := u.cfgBase.Nodes.Get(defNode)
	if n == nil {
		panic(jutil.NewNullPointerException())
	}
	if u.totalUses == 0 {
		return false
	}
	v := u.lookupVarIdxByName(name)
	if v < 0 {
		return false
	}
	if u.totalUses <= 64 {
		m := u.out64[n.Index] & u.varMask64[v]
		return m != 0 && (m&(m-1)) == 0
	}
	start, end := u.varBitStart[v], u.varBitEnd[v]
	if start >= end {
		return false
	}
	words := u.useWords
	outRow := u.outWords[int(n.Index)*words : (int(n.Index)+1)*words]
	sw, sb := start>>6, uint(start&63)
	ew, eb := (end-1)>>6, uint((end-1)&63)
	if sw == ew {
		w := outRow[sw] & ((^uint64(0) >> (63 - (eb - sb))) << sb)
		return w != 0 && (w&(w-1)) == 0
	}
	count := 0
	if w := outRow[sw] & (^uint64(0) << sb); w != 0 {
		if (w & (w - 1)) != 0 {
			return false
		}
		count++
	}
	for wi := sw + 1; wi < ew; wi++ {
		if w := outRow[wi]; w != 0 {
			if count > 0 || (w&(w-1)) != 0 {
				return false
			}
			count++
		}
	}
	if w := outRow[ew] & (^uint64(0) >> (63 - eb)); w != 0 {
		if count > 0 || (w&(w-1)) != 0 {
			return false
		}
		count++
	}
	return count == 1
}

// GetUses gets a list of nodes that may be using the value assigned to name
// in defNode. defNode must be one of the control flow graph nodes.
//
// name is the name of the variable. It can only be names of local variable
// that are not function parameters, escaped variables or variables declared
// in catch. defNode is the control flow graph node that may assign a value
// to name. Returns the list of upward exposed uses of the variable name at
// defNode.
//
// Java: MaybeReachingVariableUse.getUses(String, Node)
func (u *MaybeReachingVariableUse) GetUses(name string, defNode rhino.Node) iter.Seq[rhino.Node] {
	n := u.cfgBase.Nodes.Get(defNode)
	if n == nil {
		panic(jutil.NewNullPointerException())
	}
	if u.totalUses == 0 {
		return mrvuEmptyUses
	}
	v := u.lookupVarIdxByName(name)
	if v < 0 {
		return mrvuEmptyUses
	}
	if u.totalUses <= 64 {
		m0 := u.out64[n.Index] & u.varMask64[v]
		if m0 == 0 {
			return mrvuEmptyUses
		}
		useNodes := u.useNodes
		return func(yield func(rhino.Node) bool) {
			for m := m0; m != 0; m &= m - 1 {
				if !yield(useNodes[bits.TrailingZeros64(m)]) {
					return
				}
			}
		}
	}
	start, end := u.varBitStart[v], u.varBitEnd[v]
	if start >= end {
		return mrvuEmptyUses
	}
	words := u.useWords
	outRow := u.outWords[int(n.Index)*words : (int(n.Index)+1)*words]
	useNodes := u.useNodes
	return func(yield func(rhino.Node) bool) {
		sw, sb := start>>6, uint(start&63)
		ew, eb := (end-1)>>6, uint((end-1)&63)
		if sw == ew {
			w := outRow[sw] & ((^uint64(0) >> (63 - (eb - sb))) << sb)
			base := sw << 6
			for m := w; m != 0; m &= m - 1 {
				if !yield(useNodes[base+bits.TrailingZeros64(m)]) {
					return
				}
			}
			return
		}
		if w := outRow[sw] & (^uint64(0) << sb); w != 0 {
			base := sw << 6
			for m := w; m != 0; m &= m - 1 {
				if !yield(useNodes[base+bits.TrailingZeros64(m)]) {
					return
				}
			}
		}
		for wi := sw + 1; wi < ew; wi++ {
			if w := outRow[wi]; w != 0 {
				base := wi << 6
				for m := w; m != 0; m &= m - 1 {
					if !yield(useNodes[base+bits.TrailingZeros64(m)]) {
						return
					}
				}
			}
		}
		if w := outRow[ew] & (^uint64(0) >> (63 - eb)); w != 0 {
			base := ew << 6
			for m := w; m != 0; m &= m - 1 {
				if !yield(useNodes[base+bits.TrailingZeros64(m)]) {
					return
				}
			}
		}
	}
}
