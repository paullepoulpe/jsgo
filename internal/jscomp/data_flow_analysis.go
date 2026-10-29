// Port of com/google/javascript/jscomp/DataFlowAnalysis.java @ closure-compiler v20261021

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

	"github.com/paullepoulpe/jsgo/internal/graph"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// DataFlowAnalysis_MAX_STEPS_PER_NODE is the maximum number of steps per
// individual CFG node before we assume the analysis is divergent.
//
// TODO(b/196398705): This is way too high. Find traversal ordering
// heurisitc that reduces it.
const DataFlowAnalysis_MAX_STEPS_PER_NODE = 20000

// dfaReuseMaxNodes is the largest flow state slice or queue that an analysis
// keeps for its next CFG.
const dfaReuseMaxNodes = 1024

// cfgWorklist is a reusable, non-generic worklist over AST CFG nodes (min-heap
// on Priority when forward is true and the CFG assigns node priorities, max-heap
// when forward is false, else FIFO).
type cfgWorklist struct {
	queue       []*graph.LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch]
	queueHead   int
	usePriority bool
	forward     bool
}

func (w *cfgWorklist) reset(cfg *ControlFlowGraph[rhino.Node], numNodes int, forward bool) {
	clear(w.queue)
	w.queue = w.queue[:0]
	w.queueHead = 0
	w.usePriority = cfg.hasNodePriorities
	w.forward = forward
	if cap(w.queue) < numNodes {
		w.queue = make([]*graph.LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch], 0, numNodes)
	}
}

func (w *cfgWorklist) releaseLarge(maxCap int) {
	if cap(w.queue) > maxCap {
		w.queue = nil
		w.queueHead = 0
	}
}

func (w *cfgWorklist) notEmpty() bool {
	return len(w.queue) > w.queueHead
}

func (w *cfgWorklist) enqueue(node *graph.LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch]) {
	if node == nil {
		panic(jutil.NewNullPointerException())
	}
	if node.DfaInQueue {
		return
	}
	node.DfaInQueue = true
	w.queue = append(w.queue, node)
	if w.usePriority {
		i := len(w.queue) - 1
		xPrio := node.Priority
		if xPrio < 0 {
			panic(jutil.NewIllegalStateExceptionMsg("priority not set"))
		}
		if w.forward {
			for i > 0 {
				parent := (i - 1) / 2
				p := w.queue[parent]
				if xPrio >= p.Priority {
					break
				}
				w.queue[i] = p
				i = parent
			}
		} else {
			for i > 0 {
				parent := (i - 1) / 2
				p := w.queue[parent]
				if xPrio <= p.Priority {
					break
				}
				w.queue[i] = p
				i = parent
			}
		}
		w.queue[i] = node
	}
}

func (w *cfgWorklist) dequeue() *graph.LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch] {
	var curNode *graph.LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch]
	if w.usePriority {
		n := len(w.queue)
		curNode = w.queue[0]
		last := w.queue[n-1]
		w.queue[n-1] = nil
		w.queue = w.queue[:n-1]
		if n > 1 {
			m := n - 1
			xPrio := last.Priority
			i := 0
			if w.forward {
				for {
					pick := 2*i + 1
					if pick >= m {
						break
					}
					smaller := w.queue[pick]
					smallerPrio := smaller.Priority
					if sibling := pick + 1; sibling < m {
						sib := w.queue[sibling]
						if sibPrio := sib.Priority; smallerPrio > sibPrio {
							pick = sibling
							smaller = sib
							smallerPrio = sibPrio
						}
					}
					if xPrio <= smallerPrio {
						break
					}
					w.queue[i] = smaller
					i = pick
				}
			} else {
				for {
					pick := 2*i + 1
					if pick >= m {
						break
					}
					larger := w.queue[pick]
					largerPrio := larger.Priority
					if sibling := pick + 1; sibling < m {
						sib := w.queue[sibling]
						if sibPrio := sib.Priority; largerPrio < sibPrio {
							pick = sibling
							larger = sib
							largerPrio = sibPrio
						}
					}
					if xPrio >= largerPrio {
						break
					}
					w.queue[i] = larger
					i = pick
				}
			}
			w.queue[i] = last
		}
	} else {
		curNode = w.queue[w.queueHead]
		w.queue[w.queueHead] = nil
		w.queueHead++
		if w.queueHead == len(w.queue) {
			w.queue = w.queue[:0]
			w.queueHead = 0
		}
	}
	curNode.DfaInQueue = false
	return curNode
}

type forwardCfgWorklist struct {
	cfgWorklist
}

func (w *forwardCfgWorklist) reset(cfg *ControlFlowGraph[rhino.Node], numNodes int) {
	w.cfgWorklist.reset(cfg, numNodes, true)
}

func (w *forwardCfgWorklist) releaseLarge() {
	w.cfgWorklist.releaseLarge(dfaReuseMaxNodes)
}

// dfaBitSet is a compact bitset indexed by dense local variable ID.
// Functions with <= 128 variables stay strictly in w0/w1 with zero heap allocation.
type dfaBitSet struct {
	w0    uint64
	w1    uint64
	extra []uint64
}

func (b *dfaBitSet) clear(releaseLarge bool) {
	b.w0 = 0
	b.w1 = 0
	if releaseLarge && cap(b.extra) > 16 {
		b.extra = nil
	} else if len(b.extra) > 0 {
		clear(b.extra)
		b.extra = b.extra[:0]
	}
}

func (b *dfaBitSet) copyFrom(src *dfaBitSet) {
	b.w0 = src.w0
	b.w1 = src.w1
	if len(src.extra) == 0 {
		if len(b.extra) > 0 {
			clear(b.extra)
			b.extra = b.extra[:0]
		}
		return
	}
	if cap(b.extra) >= len(src.extra) {
		b.extra = b.extra[:len(src.extra)]
	} else {
		b.extra = make([]uint64, len(src.extra))
	}
	copy(b.extra, src.extra)
}

func (b *dfaBitSet) set(idx int) {
	if idx < 64 {
		b.w0 |= uint64(1) << uint(idx)
		return
	}
	if idx < 128 {
		b.w1 |= uint64(1) << uint(idx-64)
		return
	}
	wi := (idx - 128) >> 6
	if wi >= len(b.extra) {
		if cap(b.extra) > wi {
			oldLen := len(b.extra)
			b.extra = b.extra[:wi+1]
			clear(b.extra[oldLen:])
		} else {
			extra := make([]uint64, wi+1, max(wi+1, 4))
			copy(extra, b.extra)
			b.extra = extra
		}
	}
	b.extra[wi] |= uint64(1) << uint((idx-128)&63)
}

func (b *dfaBitSet) has(idx int) bool {
	if idx < 64 {
		return (b.w0 & (uint64(1) << uint(idx))) != 0
	}
	if idx < 128 {
		return (b.w1 & (uint64(1) << uint(idx-64))) != 0
	}
	wi := (idx - 128) >> 6
	return wi < len(b.extra) && (b.extra[wi]&(uint64(1)<<uint((idx-128)&63))) != 0
}

func (b *dfaBitSet) isEmpty() bool {
	if (b.w0 | b.w1) != 0 {
		return false
	}
	for _, w := range b.extra {
		if w != 0 {
			return false
		}
	}
	return true
}

func (b *dfaBitSet) count() int {
	c := bits.OnesCount64(b.w0) + bits.OnesCount64(b.w1)
	for _, w := range b.extra {
		c += bits.OnesCount64(w)
	}
	return c
}

func (b *dfaBitSet) equals(o *dfaBitSet) bool {
	if b.w0 != o.w0 || b.w1 != o.w1 {
		return false
	}
	if len(b.extra) == 0 && len(o.extra) == 0 {
		return true
	}
	n := min(len(b.extra), len(o.extra))
	for i := 0; i < n; i++ {
		if b.extra[i] != o.extra[i] {
			return false
		}
	}
	for _, w := range b.extra[n:] {
		if w != 0 {
			return false
		}
	}
	for _, w := range o.extra[n:] {
		if w != 0 {
			return false
		}
	}
	return true
}

func (b *dfaBitSet) orWith(o *dfaBitSet) {
	b.w0 |= o.w0
	b.w1 |= o.w1
	if len(o.extra) == 0 {
		return
	}
	if len(b.extra) < len(o.extra) {
		grown := make([]uint64, len(o.extra))
		copy(grown, b.extra)
		b.extra = grown
	}
	for i, w := range o.extra {
		b.extra[i] |= w
	}
}

func (b *dfaBitSet) clone() dfaBitSet {
	out := dfaBitSet{w0: b.w0, w1: b.w1}
	if len(b.extra) > 0 {
		out.extra = slices.Clone(b.extra)
	}
	return out
}

func (b *dfaBitSet) intersects(other *dfaBitSet) bool {
	if (b.w0&other.w0)|(b.w1&other.w1) != 0 {
		return true
	}
	n := min(len(b.extra), len(other.extra))
	for i := 0; i < n; i++ {
		if b.extra[i]&other.extra[i] != 0 {
			return true
		}
	}
	return false
}

func (b *dfaBitSet) containsAll(other *dfaBitSet) bool {
	if (other.w0&^b.w0)|(other.w1&^b.w1) != 0 {
		return false
	}
	for wi, ow := range other.extra {
		if ow != 0 && (wi >= len(b.extra) || ow&^b.extra[wi] != 0) {
			return false
		}
	}
	return true
}

func (b *dfaBitSet) forEach(fn func(idx int)) {
	for w := b.w0; w != 0; w &= w - 1 {
		fn(bits.TrailingZeros64(w))
	}
	for w := b.w1; w != 0; w &= w - 1 {
		fn(64 + bits.TrailingZeros64(w))
	}
	for wi, word := range b.extra {
		base := 128 + wi*64
		for w := word; w != 0; w &= w - 1 {
			fn(base + bits.TrailingZeros64(w))
		}
	}
}

func (b *dfaBitSet) nextSetBit(fromIndex int) int {
	if fromIndex < 64 {
		if w0 := b.w0 & (^uint64(0) << uint(fromIndex)); w0 != 0 {
			return bits.TrailingZeros64(w0)
		}
		fromIndex = 64
	}
	if fromIndex < 128 {
		if w1 := b.w1 & (^uint64(0) << uint(fromIndex-64)); w1 != 0 {
			return 64 + bits.TrailingZeros64(w1)
		}
		fromIndex = 128
	}
	if len(b.extra) == 0 {
		return -1
	}
	wi := (fromIndex - 128) >> 6
	if wi >= len(b.extra) {
		return -1
	}
	if w := b.extra[wi] & (^uint64(0) << uint((fromIndex-128)&63)); w != 0 {
		return 128 + wi*64 + bits.TrailingZeros64(w)
	}
	for wi++; wi < len(b.extra); wi++ {
		if w := b.extra[wi]; w != 0 {
			return 128 + wi*64 + bits.TrailingZeros64(w)
		}
	}
	return -1
}

const fnVarTableInlineCap = 16

type fnVarTable struct {
	vars              []*Var
	idToIndex         map[uint32]int
	escaped           dfaBitSet
	uniqueCount       int
	hasDuplicateNames bool
}

func (t *fnVarTable) reset() {
	clear(t.vars)
	if cap(t.vars) > dfaReuseMaxNodes {
		t.vars = nil
	} else {
		t.vars = t.vars[:0]
	}
	t.escaped.clear(true)
	t.uniqueCount = 0
	t.hasDuplicateNames = false
}

func (t *fnVarTable) isEmpty() bool {
	return len(t.vars) == 0
}

func (t *fnVarTable) appendScopeVars(currentScope *Scope) {
	n := currentScope.GetVarCount()
	if n == 0 {
		return
	}
	if need := len(t.vars) + n; need > cap(t.vars) {
		newCap := max(need, 2*cap(t.vars), 16)
		buf := make([]*Var, len(t.vars), newCap)
		copy(buf, t.vars)
		t.vars = buf
	}
	for _, v := range currentScope.varList {
		if v != nil {
			t.vars = append(t.vars, v)
		}
	}
}

func (t *fnVarTable) finalize() {
	n := len(t.vars)
	if n == 0 {
		if len(t.idToIndex) > dfaReuseMaxNodes {
			t.idToIndex = nil
		} else if len(t.idToIndex) > 0 {
			clear(t.idToIndex)
		}
		t.uniqueCount = 0
		t.hasDuplicateNames = false
		return
	}
	if n <= fnVarTableInlineCap {
		unique := n
		for i := 1; i < n; i++ {
			id := t.vars[i].nameID
			for j := 0; j < i; j++ {
				if t.vars[j].nameID == id {
					unique--
					break
				}
			}
		}
		t.uniqueCount = unique
		t.hasDuplicateNames = unique < n
		return
	}
	if t.idToIndex == nil || len(t.idToIndex) > dfaReuseMaxNodes {
		t.idToIndex = make(map[uint32]int, max(n, 32))
	} else {
		clear(t.idToIndex)
	}
	for i, v := range t.vars {
		t.idToIndex[v.nameID] = i
	}
	t.uniqueCount = len(t.idToIndex)
	t.hasDuplicateNames = t.uniqueCount < n
}

func (t *fnVarTable) lookupIndexByID(id uint32) (int, bool) {
	if id == 0 || len(t.vars) == 0 {
		return 0, false
	}
	if len(t.vars) <= fnVarTableInlineCap {
		for i := len(t.vars) - 1; i >= 0; i-- {
			if t.vars[i].nameID == id {
				return i, true
			}
		}
		return 0, false
	}
	idx, ok := t.idToIndex[id]
	return idx, ok
}

func (t *fnVarTable) lookupVarByID(id uint32) *Var {
	idx, ok := t.lookupIndexByID(id)
	if !ok {
		return nil
	}
	return t.vars[idx]
}

func (t *fnVarTable) lookupVar(name string) *Var {
	return t.lookupVarByID(rhino.RhinoStringPool_LookupID(name))
}

func (t *fnVarTable) resolveExactVarIndex(v *Var) (int, bool) {
	if v == nil {
		return -1, false
	}
	idx, ok := t.lookupIndexByID(v.nameID)
	if !ok {
		return -1, false
	}
	if t.vars[idx] == v || t.vars[idx].Equals(v) {
		return idx, true
	}
	if t.hasDuplicateNames {
		for i, cand := range t.vars {
			if cand == v || cand.Equals(v) {
				return i, true
			}
		}
	}
	return -1, false
}

// dataFlowAnalysis_linearFlowStateOf is the unchecked cast
// `LinearFlowState<L> state = node.getAnnotation()`.
//
// Go-only: language-adapter — Java's generic getAnnotation() casts at the
// call site; a null annotation stays null.
func dataFlowAnalysis_linearFlowStateOf[L any](a graph.Annotatable) *DataFlowAnalysis_LinearFlowState[L] {
	annotation := a.GetAnnotation()
	if annotation == nil {
		return nil
	}
	return annotation.(*DataFlowAnalysis_LinearFlowState[L])
}

// DataFlowAnalysis_LinearFlowState holds the in and out states of a node.
type DataFlowAnalysis_LinearFlowState[L any] struct {
	stepCount int
	in        L
	out       L
}

// Java: DataFlowAnalysis.LinearFlowState(L, L)
func newDataFlowAnalysis_LinearFlowState[L any](in L, out L) *DataFlowAnalysis_LinearFlowState[L] {
	jutil.CheckNotNull(in)
	jutil.CheckNotNull(out)
	return &DataFlowAnalysis_LinearFlowState[L]{stepCount: 0, in: in, out: out}
}

// Java: DataFlowAnalysis.LinearFlowState.getStepCount()
func (s *DataFlowAnalysis_LinearFlowState[L]) GetStepCount() int {
	return s.stepCount
}

// Java: DataFlowAnalysis.LinearFlowState.getIn()
func (s *DataFlowAnalysis_LinearFlowState[L]) GetIn() L {
	return s.in
}

// Java: DataFlowAnalysis.LinearFlowState.setIn(L)
func (s *DataFlowAnalysis_LinearFlowState[L]) SetIn(in L) {
	jutil.CheckNotNull(in)
	s.in = in
}

// Java: DataFlowAnalysis.LinearFlowState.getOut()
func (s *DataFlowAnalysis_LinearFlowState[L]) GetOut() L {
	return s.out
}

// Java: DataFlowAnalysis.LinearFlowState.setOut(L)
func (s *DataFlowAnalysis_LinearFlowState[L]) SetOut(out L) {
	jutil.CheckNotNull(out)
	s.out = out
}

// Java: DataFlowAnalysis.LinearFlowState.toString()
func (s *DataFlowAnalysis_LinearFlowState[L]) String() string {
	return jutil.StringFormat("IN: %s OUT: %s", s.in, s.out)
}

func dfaHasNestedFunctionOrClass(n rhino.Node) bool {
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		tok := c.GetToken()
		if tok == rhino.Token_FUNCTION || tok == rhino.Token_CLASS {
			return true
		}
		if dfaHasNestedFunctionOrClass(c) {
			return true
		}
	}
	return false
}

// DataFlowAnalysis_ComputeEscaped computes the set of escaped variables.
// When a variable is escaped in a dataflow analysis, it can be referenced
// outside of the code that we are analyzing. A variable is escaped if any
// of the following is true:
//
// 1. Exported variables as they can be needed after the script terminates.
// 2. Names of named functions because in JavaScript, function foo(){} does
// not kill foo in the dataflow.
//
// jsScope must be a function scope.
//
// Java: DataFlowAnalysis.computeEscaped(Scope, Set<Var>, AbstractCompiler, ScopeCreator, Map<String, Var>)
func DataFlowAnalysis_ComputeEscaped(
	jsScope *Scope,
	escaped *jutil.ObjectLinkedHashSet[*Var],
	compiler AbstractCompiler,
	scopeCreator ScopeCreator,
	allVarsInFn *jutil.LinkedHashMap[string, *Var]) {

	// Optimize for Class Static Blocks.
	jutil.CheckArgument(jsScope.IsFunctionScope())
	if allVarsInFn != nil && allVarsInFn.IsEmpty() {
		return
	}

	escapedStartedEmpty := escaped.IsEmpty()
	hasNestedScopeVars := false

	// TODO (simranarora) catch variables should not be considered escaped in
	// ES6. Getting rid of the catch check is causing breakages however
	codingConvention := compiler.GetCodingConvention()
	if allVarsInFn != nil {
		for var_ := range allVarsInFn.Values() {
			if var_.scope != jsScope && (var_.scope == nil || !var_.scope.IsFunctionBlockScope()) {
				hasNestedScopeVars = true
			}
			if var_.GetParentNode().IsCatch() ||
				codingConvention.IsExported(var_.GetName()) {
				escaped.Add(var_)
			}
		}
	}

	fnRoot := jsScope.GetRootNode()
	fnNameNode := fnRoot.GetFirstChild()
	fnNameID := fnNameNode.GetStringID()
	fnNameShadowed := false
	if fnNameID > rhino.RhinoStringPool_EmptyID && allVarsInFn != nil && jsScope.hasOwnSlotByID(fnNameID) {
		if v := allVarsInFn.Get(fnNameNode.GetString()); v == nil || v.scope != jsScope {
			fnNameShadowed = true
		}
	}

	if allVarsInFn != nil && !hasNestedScopeVars && !fnNameShadowed && escapedStartedEmpty && escaped.Size() >= allVarsInFn.Size() {
		return
	}

	args := fnNameNode.GetNext()
	fnBody := args.GetNext()
	argsHasFnOrClass := dfaHasNestedFunctionOrClass(args)

	var firstEscapingStmt rhino.Node
	if !argsHasFnOrClass && fnBody.IsBlock() {
		for stmt := fnBody.GetFirstChild(); !stmt.IsNil(); stmt = stmt.GetNext() {
			if stmt.IsFunction() || stmt.IsClass() || dfaHasNestedFunctionOrClass(stmt) {
				firstEscapingStmt = stmt
				break
			}
		}
		if firstEscapingStmt.IsNil() {
			return
		}
	} else if !argsHasFnOrClass && !dfaHasNestedFunctionOrClass(fnBody) {
		return
	}

	hasDuplicateVarNames := allVarsInFn == nil || hasNestedScopeVars || fnNameShadowed || argsHasFnOrClass
	nonEscapingDepth := jsScope.GetDepth()
	var skipArgs rhino.Node
	var filterFnBody rhino.Node
	if !argsHasFnOrClass {
		skipArgs = args
		if fnBody.IsBlock() {
			nonEscapingDepth = jsScope.GetDepth() + 1
			filterFnBody = fnBody
		}
	}

	finder := dataFlowAnalysis_escapedFinder{
		jsScope:              jsScope,
		nonEscapingDepth:     nonEscapingDepth,
		escaped:              escaped,
		allVarsInFn:          allVarsInFn,
		skipArgs:             skipArgs,
		fnBody:               filterFnBody,
		firstEscapingStmt:    firstEscapingStmt,
		hasDuplicateVarNames: hasDuplicateVarNames,
		escapedStartedEmpty:  escapedStartedEmpty,
	}

	b := NodeTraversal_Builder{
		compiler:     compiler,
		callback:     &finder,
		scopeCreator: scopeCreator,
	}
	b.TraverseAtScope(jsScope)
}

func computeEscapedInTable(
	jsScope *Scope,
	table *fnVarTable,
	compiler AbstractCompiler,
	scopeCreator ScopeCreator,
) {
	jutil.CheckArgument(jsScope.IsFunctionScope())
	if table.isEmpty() {
		return
	}

	escapedStartedEmpty := table.escaped.isEmpty()
	hasNestedScopeVars := false

	codingConvention := compiler.GetCodingConvention()
	for i, var_ := range table.vars {
		if table.hasDuplicateNames {
			if canonIdx, _ := table.lookupIndexByID(var_.nameID); canonIdx != i {
				continue
			}
		}
		if var_.scope != jsScope && (var_.scope == nil || !var_.scope.IsFunctionBlockScope()) {
			hasNestedScopeVars = true
		}
		if var_.GetParentNode().IsCatch() ||
			codingConvention.IsExported(var_.GetName()) {
			table.escaped.set(i)
		}
	}

	fnRoot := jsScope.GetRootNode()
	fnNameNode := fnRoot.GetFirstChild()
	fnNameID := fnNameNode.GetStringID()
	fnNameShadowed := false
	if fnNameID > rhino.RhinoStringPool_EmptyID && jsScope.hasOwnSlotByID(fnNameID) {
		if v := table.lookupVarByID(fnNameID); v == nil || v.scope != jsScope {
			fnNameShadowed = true
		}
	}

	if !hasNestedScopeVars && !fnNameShadowed && escapedStartedEmpty && table.escaped.count() >= table.uniqueCount {
		return
	}

	args := fnNameNode.GetNext()
	fnBody := args.GetNext()
	argsHasFnOrClass := dfaHasNestedFunctionOrClass(args)

	var firstEscapingStmt rhino.Node
	if !argsHasFnOrClass && fnBody.IsBlock() {
		for stmt := fnBody.GetFirstChild(); !stmt.IsNil(); stmt = stmt.GetNext() {
			if stmt.IsFunction() || stmt.IsClass() || dfaHasNestedFunctionOrClass(stmt) {
				firstEscapingStmt = stmt
				break
			}
		}
		if firstEscapingStmt.IsNil() {
			return
		}
	} else if !argsHasFnOrClass && !dfaHasNestedFunctionOrClass(fnBody) {
		return
	}

	hasDuplicateVarNames := table.hasDuplicateNames || hasNestedScopeVars || fnNameShadowed || argsHasFnOrClass
	nonEscapingDepth := jsScope.GetDepth()
	var skipArgs rhino.Node
	var filterFnBody rhino.Node
	if !argsHasFnOrClass {
		skipArgs = args
		if fnBody.IsBlock() {
			nonEscapingDepth = jsScope.GetDepth() + 1
			filterFnBody = fnBody
		}
	}

	finder := dataFlowAnalysis_escapedFinder{
		jsScope:              jsScope,
		nonEscapingDepth:     nonEscapingDepth,
		table:                table,
		skipArgs:             skipArgs,
		fnBody:               filterFnBody,
		firstEscapingStmt:    firstEscapingStmt,
		hasDuplicateVarNames: hasDuplicateVarNames,
		escapedStartedEmpty:  escapedStartedEmpty,
	}

	b := NodeTraversal_Builder{
		compiler:     compiler,
		callback:     &finder,
		scopeCreator: scopeCreator,
	}
	b.TraverseAtScope(jsScope)
}

// dataFlowAnalysis_escapedFinder is the anonymous AbstractPostOrderCallback
// of computeEscaped.
type dataFlowAnalysis_escapedFinder struct {
	jsScope              *Scope
	nonEscapingDepth     int
	table                *fnVarTable
	escaped              *jutil.ObjectLinkedHashSet[*Var]
	allVarsInFn          *jutil.LinkedHashMap[string, *Var]
	skipArgs             rhino.Node
	fnBody               rhino.Node
	firstEscapingStmt    rhino.Node
	reachedFirstEscaping bool
	hasDuplicateVarNames bool
	escapedStartedEmpty  bool
}

func (f *dataFlowAnalysis_escapedFinder) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if !f.hasDuplicateVarNames && f.escapedStartedEmpty {
		if f.table != nil {
			if f.table.escaped.count() >= f.table.uniqueCount {
				return false
			}
		} else if f.escaped.Size() >= f.allVarsInFn.Size() {
			return false
		}
	}
	if n == f.skipArgs {
		return false
	}
	if !parent.IsNil() && parent == f.fnBody {
		if !f.reachedFirstEscaping {
			if n != f.firstEscapingStmt {
				return false
			}
			f.reachedFirstEscaping = true
			return true
		}
		return n.IsFunction() || n.IsClass() || dfaHasNestedFunctionOrClass(n)
	}
	return true
}

func (f *dataFlowAnalysis_escapedFinder) isInsideEscapingScope(t *NodeTraversal) bool {
	for i := len(t.scopes) - 1; i > f.nonEscapingDepth; i-- {
		root := t.scopes[i].node
		switch root.GetToken() {
		case rhino.Token_FUNCTION, rhino.Token_SCRIPT, rhino.Token_MODULE_BODY, rhino.Token_ROOT:
			return true
		case rhino.Token_BLOCK:
			if NodeUtil_IsClassStaticBlock(root) {
				return true
			}
		case rhino.Token_MEMBER_FIELD_DEF, rhino.Token_COMPUTED_FIELD_DEF:
			if !root.IsStaticMember() {
				return true
			}
		}
	}
	return false
}

func (*dataFlowAnalysis_escapedFinder) onlyNameLeaves() {}

// Java: DataFlowAnalysis.computeEscaped(...).new AbstractPostOrderCallback() {...}.visit(NodeTraversal, Node, Node)
func (f *dataFlowAnalysis_escapedFinder) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {

	if !n.IsName() || parent.IsFunction() || len(t.scopes)-1 <= f.nonEscapingDepth || !f.isInsideEscapingScope(t) {
		return
	}

	id := n.GetStringID()
	if id <= rhino.RhinoStringPool_EmptyID {
		return
	}
	if f.table != nil {
		candIdx, ok := f.table.lookupIndexByID(id)
		if !ok {
			return
		}
		if !f.hasDuplicateVarNames && f.table.escaped.has(candIdx) {
			return
		}
	} else if f.allVarsInFn != nil {
		targetVar := f.allVarsInFn.Get(n.GetString())
		if targetVar == nil {
			return
		}
		if !f.hasDuplicateVarNames && f.escaped.Contains(targetVar) {
			return
		}
	}

	var_ := t.LookupVarByID(id)
	if var_ == nil {
		return
	}

	// AbstractVar.getScope() returns S (here Scope); see abstract_scope.go.
	varScope, _ := var_.GetScope().(*Scope)
	variableCfgScope := varScope.GetClosestCfgRootScope()

	// If the current examined scope and var decleration scope are different
	// then var is 'not escaped'.
	if variableCfgScope != f.jsScope {
		return
	}

	// Variable is referenced outside of declaration scope it is 'escaped'.
	// Non-static class field initializers are evaluated at instance creation time
	// (like closures), so variables referenced in non-static class fields must also be
	// marked as escaped.
	if f.table != nil {
		if idx, ok := f.table.resolveExactVarIndex(var_); ok {
			f.table.escaped.set(idx)
		}
	} else {
		f.escaped.Add(var_)
	}
}

// Java: DataFlowAnalysis.isReferencedInNonStaticClassField(Scope, Scope)
func DataFlowAnalysis_IsReferencedInNonStaticClassField(refScope *Scope, targetScope *Scope) bool {
	for s := refScope; s != nil && s != targetScope; s = s.GetParent() {
		root := s.GetRootNode()
		if (root.IsMemberFieldDef() || root.IsComputedFieldDef()) && !root.IsStaticMember() {
			return true
		}
	}
	return false
}
