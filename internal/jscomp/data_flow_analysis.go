// Port of com/google/javascript/jscomp/DataFlowAnalysis.java @ closure-compiler v20261015

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// DataFlowAnalysis is a framework to help writing static program analysis.
//
// A subclass of this framework should specify how a single node changes
// the state of a program. This class finds a safe estimate (a fixed-point)
// for the whole program. The proven facts about the program will be
// annotated with GraphNode#setAnnotation to the given control flow graph's
// nodes in form of LatticeElement after calling Analyze().
//
// As a guideline, the following is a list of behaviors that any analysis
// can take:
//
//  1. Flow Direction: Is the analysis a forward or backward analysis?
//  2. Lattice Elements: How does the analysis represent the state of the
//     program at any given point?
//  3. Branching: Does the analysis propagate a different lattice element
//     along each branch exiting a node?
//  4. JOIN Operation: Given two incoming paths and a lattice state value,
//     what can the compiler conclude at the join point?
//  5. Flow Equations: How does an instruction modify the state of program
//     in terms of lattice values?
//  6. Initial Entry Value: What can the compiler assume at the beginning of
//     the program?
//  7. Initial Estimate: What can the compiler assume at each point of the
//     program? (What is the BOTTOM value of the lattice) By definition this
//     lattice JOIN x for any x must also be x.
//  8. (Optional) Branch Operation: How should the flow branch along edges?
//
// To make these behaviors known to the framework, the following steps must
// be taken.
//
//  1. Flow Direction: Implement IsForward().
//  2. Lattice Elements: Implement LatticeElement.
//  3. JOIN Operation: Implement JoinOp#apply.
//  4. Flow Equations: Implement FlowThrough(Object, LatticeElement).
//  5. Initial Entry Value: Implement CreateEntryLattice().
//  6. Initial Estimate: Implement CreateInitialEstimateLattice().
//  7. (Optional) Branch Operation: Return true from IsBranched() and
//     implement CreateFlowBrancher.
//
// Upon execution of the Analyze() method, nodes of the input control flow
// graph will be annotated with a FlowState object that represents maximum
// fixed point solution. Any previous annotations at the nodes of the
// control flow graph will be lost.
//
// N is the control flow graph's node value type; L the lattice element
// type.
type DataFlowAnalysis[N comparable, L LatticeElement] interface {
	// IsForward checks whether the analysis is a forward flow analysis or
	// backward flow analysis. Returns true if it is a forward analysis.
	//
	// Java: DataFlowAnalysis.isForward()
	IsForward() bool

	// IsBranched reports whether or not CreateFlowBrancher should be used.
	//
	// Java: DataFlowAnalysis.isBranched()
	IsBranched() bool

	// CreateFlowJoiner gets a new joiner for an analysis step.
	//
	// The joiner is invoked once for each input edge and then the final
	// joined result is retrieved. No joiner will be created for a single
	// input.
	//
	// Java: DataFlowAnalysis.createFlowJoiner()
	CreateFlowJoiner() DataFlowAnalysis_FlowJoiner[L]

	// CreateFlowBrancher gets a new brancher for an analysis step.
	//
	// The brancher is invoked for each output edge. The returned states are
	// annotated onto each edge and compared against the previous state to
	// determine whether a fixed-point has been reached.
	//
	// node is the node at which to branch; output the output state of
	// FlowThrough at node.
	//
	// Java: DataFlowAnalysis.createFlowBrancher(N, L)
	CreateFlowBrancher(node N, output L) DataFlowAnalysis_FlowBrancher[L]

	// FlowThrough computes the output state for a given node given its input
	// state.
	//
	// node is the node; input the input lattice that should be read-only.
	// Returns the output lattice.
	//
	// Java: DataFlowAnalysis.flowThrough(N, L)
	FlowThrough(node N, input L) L

	// CreateInitialEstimateLattice gets the state of the initial estimation
	// at each node.
	//
	// Java: DataFlowAnalysis.createInitialEstimateLattice()
	CreateInitialEstimateLattice() L

	// CreateEntryLattice gets the incoming state of the entry node.
	//
	// Java: DataFlowAnalysis.createEntryLattice()
	CreateEntryLattice() L

	GetDataFlowAnalysisBase() *DataFlowAnalysisBase[N, L]
}

// DataFlowAnalysisBase holds DataFlowAnalysis's fields and its base and
// final implementations.
type DataFlowAnalysisBase[N comparable, L LatticeElement] struct {
	self DataFlowAnalysis[N, L]

	cfg       ControlFlowGraph[N]
	workQueue *dataFlowAnalysis_UniqueQueue[DiGraph_DiGraphNode[N, ControlFlowGraph_Branch]]
}

// DataFlowAnalysis_MAX_STEPS_PER_NODE is the maximum number of steps per
// individual CFG node before we assume the analysis is divergent.
//
// TODO(b/196398705): This is way too high. Find traversal ordering
// heurisitc that reduces it.
const DataFlowAnalysis_MAX_STEPS_PER_NODE = 20000

// InitDataFlowAnalysisBase constructs a data flow analysis; self is the
// analysis being constructed.
//
// Typical usage
//
//	DataFlowAnalysis dfa = ...
//	dfa.analyze();
//
// Analyze() annotates the result to the control flow graph by means of
// DiGraphNode#setAnnotation without any modification of the graph itself.
// Additional calls to Analyze() recomputes the analysis which can be useful
// if the control flow graph has been modified.
//
// cfg is the control flow graph object that this object performs on.
// Modification of the graph requires a separate call to Analyze().
//
// Java: DataFlowAnalysis(ControlFlowGraph<N>)
func (b *DataFlowAnalysisBase[N, L]) InitDataFlowAnalysisBase(self DataFlowAnalysis[N, L], cfg ControlFlowGraph[N]) {
	b.self = self
	b.cfg = cfg
	// Go-only: language-adapter — each analysis owns the CFG annotations while
	// it runs. Keep queue membership beside the existing per-node flow state
	// instead of allocating a second pointer-keyed ordered set. The queue's
	// comparator and encounter order are unchanged.
	b.workQueue = newDataFlowAnalysis_UniqueQueueWithMembership(
		cfg.GetOptionalNodeComparator(b.self.IsForward()),
		func(node DiGraph_DiGraphNode[N, ControlFlowGraph_Branch]) *bool {
			return &dataFlowAnalysis_linearFlowStateOf[L](node).queued
		})

	if b.self.IsBranched() {
		jutil.CheckState(b.self.IsForward())
	}
}

// GetDataFlowAnalysisBase gives access to the base through the interface.
func (b *DataFlowAnalysisBase[N, L]) GetDataFlowAnalysisBase() *DataFlowAnalysisBase[N, L] {
	return b
}

// GetCfg returns the control flow graph that this analysis was performed
// on. Modifications can be done on this graph, however, the only time that
// the annotations are correct is after Analyze() is called and before the
// graph has been modified.
//
// Java: DataFlowAnalysis.getCfg()
func (b *DataFlowAnalysisBase[N, L]) GetCfg() ControlFlowGraph[N] {
	return b.cfg
}

// Java: DataFlowAnalysis.isBranched()
func (b *DataFlowAnalysisBase[N, L]) IsBranched() bool {
	return false
}

// Java: DataFlowAnalysis.join(L, L)
func (b *DataFlowAnalysisBase[N, L]) Join(latticeA L, latticeB L) L {
	joiner := b.self.CreateFlowJoiner()
	joiner.JoinFlow(latticeA)
	joiner.JoinFlow(latticeB)
	return joiner.Finish()
}

// DataFlowAnalysis_FlowJoiner is a reducer that joins flow states from
// distinct input states into a single input state.
type DataFlowAnalysis_FlowJoiner[L any] interface {
	// Java: DataFlowAnalysis.FlowJoiner.joinFlow(L)
	JoinFlow(input L)

	// Java: DataFlowAnalysis.FlowJoiner.finish()
	Finish() L
}

// Java: DataFlowAnalysis.createFlowBrancher(N, L)
func (b *DataFlowAnalysisBase[N, L]) CreateFlowBrancher(node N, output L) DataFlowAnalysis_FlowBrancher[L] {
	panic(jutil.NewUnsupportedOperationException())
}

// DataFlowAnalysis_FlowBrancher is a callback that branches a flow state
// into a distinct state for each output edge.
//
// Go-only: language-adapter — a @FunctionalInterface, so a Go func type;
// Java's FlowBrancher.branchFlow(Branch) is a call of the func.
type DataFlowAnalysis_FlowBrancher[L any] func(branch ControlFlowGraph_Branch) L

// Analyze finds a fixed-point solution. The function has the side effect of
// replacing the existing node annotations with the computed solutions using
// GraphNode#setAnnotation(Annotation).
//
// Initially, each node's input and output flow state contains the value
// given by CreateInitialEstimateLattice() (with the exception of the entry
// node of the graph which takes on the CreateEntryLattice() value. Each node
// will use the output state of its predecessor and compute an output state
// according to the instruction. At that time, any nodes that depend on the
// node's newly modified output value will need to recompute their output
// state again. Each step will perform a computation at one node until no
// extra computation will modify any existing output state anymore.
//
// Java: DataFlowAnalysis.analyze()
func (b *DataFlowAnalysisBase[N, L]) Analyze() {
	b.Initialize()
	for !b.workQueue.IsEmpty() {
		curNode := b.workQueue.RemoveFirst()
		curState := dataFlowAnalysis_linearFlowStateOf[L](curNode)
		stepCount := curState.stepCount
		curState.stepCount++
		if stepCount > DataFlowAnalysis_MAX_STEPS_PER_NODE {
			panic(jutil.NewIllegalStateExceptionMsg("Dataflow analysis appears to diverge around: " + jutil.ToJavaString(curNode)))
		}

		b.JoinInputs(curNode)
		if b.Flow(curNode) {
			// If there is a change in the current node, we want to grab the
			// list of nodes that this node affects.
			// Go-only: language-adapter — the CFG's predecessor/successor
			// lists are exactly its edge lists projected to their endpoints.
			// Iterate those edges directly to avoid a temporary slice on each
			// changed node, preserving duplicate edges and encounter order.
			if b.self.IsForward() {
				for _, edge := range curNode.GetOutEdges() {
					nextNode := edge.GetDestination()
					if nextNode != b.cfg.GetImplicitReturn() {
						b.workQueue.Add(nextNode)
					}
				}
			} else {
				for _, edge := range curNode.GetInEdges() {
					nextNode := edge.GetSource()
					if nextNode != b.cfg.GetImplicitReturn() {
						b.workQueue.Add(nextNode)
					}
				}
			}
		}
	}
	if b.self.IsForward() {
		b.JoinInputs(b.GetCfg().GetImplicitReturn())
	}
}

// dataFlowAnalysis_linearFlowStateOf is the unchecked cast
// `LinearFlowState<L> state = node.getAnnotation()`.
//
// Go-only: language-adapter — Java's generic getAnnotation() casts at the
// call site; a null annotation stays null.
func dataFlowAnalysis_linearFlowStateOf[L any](a Annotatable) *DataFlowAnalysis_LinearFlowState[L] {
	annotation := a.GetAnnotation()
	if annotation == nil {
		return nil
	}
	return annotation.(*DataFlowAnalysis_LinearFlowState[L])
}

// dataFlowAnalysis_latticeOf is the unchecked cast `L l =
// edge.getAnnotation()`.
//
// Go-only: language-adapter — Java's generic getAnnotation() casts at the
// call site; a null annotation stays null.
func dataFlowAnalysis_latticeOf[L any](a Annotatable) L {
	annotation := a.GetAnnotation()
	if annotation == nil {
		var zero L
		return zero
	}
	return annotation.(L)
}

// Initialize initializes the work list and the control flow graph.
//
// Java: DataFlowAnalysis.initialize()
func (b *DataFlowAnalysisBase[N, L]) Initialize() {
	// Go-only: language-adapter — an interrupted analysis can leave queued
	// nodes whose annotations have since been cleared or replaced by another
	// analysis. Discard queue storage without consulting those annotations;
	// every membership flag is recreated with its flow state below.
	b.workQueue.clearStorage()
	for _, node := range b.cfg.GetDiGraphNodes() {
		node.SetAnnotation(
			newDataFlowAnalysis_LinearFlowState(b.self.CreateInitialEstimateLattice(), b.self.CreateInitialEstimateLattice()))
		if node != b.cfg.GetImplicitReturn() {
			b.workQueue.Add(node)
		}
	}
	if b.self.IsBranched() {
		for _, edge := range b.cfg.GetDiGraphEdges() {
			edge.SetAnnotation(b.self.CreateInitialEstimateLattice())
		}
	}
}

// Flow performs a single flow through a node. Returns true if the flow
// state differs from the previous state.
//
// Java: DataFlowAnalysis.flow(DiGraphNode<N, Branch>)
func (b *DataFlowAnalysisBase[N, L]) Flow(node DiGraph_DiGraphNode[N, ControlFlowGraph_Branch]) bool {
	state := dataFlowAnalysis_linearFlowStateOf[L](node)
	if b.self.IsForward() {
		outBefore := state.GetOut()
		state.SetOut(b.self.FlowThrough(node.GetValue(), state.GetIn()))
		changed := !jutil.ObjectsEquals(outBefore, state.GetOut())

		if b.self.IsBranched() {
			brancher := b.self.CreateFlowBrancher(node.GetValue(), state.GetOut())
			for _, outEdge := range node.GetOutEdges() {
				outBranchBefore := dataFlowAnalysis_latticeOf[L](outEdge)
				outEdge.SetAnnotation(brancher(outEdge.GetValue()))
				if !changed {
					changed = !jutil.ObjectsEquals(outBranchBefore, dataFlowAnalysis_latticeOf[L](outEdge))
				}
			}
		}

		return changed
	} else {
		inBefore := state.GetIn()
		state.SetIn(b.self.FlowThrough(node.GetValue(), state.GetOut()))
		return !jutil.ObjectsEquals(inBefore, state.GetIn())
	}
}

// JoinInputs computes the new flow state at a given node's entry by merging
// the output (input) lattice of the node's predecessor (successor).
//
// node is the node to compute new join.
//
// Java: DataFlowAnalysis.joinInputs(DiGraphNode<N, Branch>)
func (b *DataFlowAnalysisBase[N, L]) JoinInputs(node DiGraph_DiGraphNode[N, ControlFlowGraph_Branch]) {
	state := dataFlowAnalysis_linearFlowStateOf[L](node)
	if b.self.IsForward() && b.cfg.GetEntry() == node {
		state.SetIn(b.self.CreateEntryLattice())
		return
	}

	var inEdges []DiGraph_DiGraphEdge[N, ControlFlowGraph_Branch]
	if b.self.IsForward() {
		inEdges = node.GetInEdges()
	} else {
		inEdges = node.GetOutEdges()
	}

	var result L
	switch len(inEdges) {
	case 0:
		return
	case 1:
		result = b.GetInputFromEdge(inEdges[0])
	default:
		joiner := b.self.CreateFlowJoiner()
		for _, inEdge := range inEdges {
			joiner.JoinFlow(b.GetInputFromEdge(inEdge))
		}
		result = joiner.Finish()
	}

	if b.self.IsForward() {
		state.SetIn(result)
	} else {
		state.SetOut(result)
	}
}

// Java: DataFlowAnalysis.getInputFromEdge(DiGraphEdge<N, Branch>)
func (b *DataFlowAnalysisBase[N, L]) GetInputFromEdge(edge DiGraph_DiGraphEdge[N, ControlFlowGraph_Branch]) L {
	if b.self.IsBranched() {
		return dataFlowAnalysis_latticeOf[L](edge)
	} else if b.self.IsForward() {
		state := dataFlowAnalysis_linearFlowStateOf[L](edge.GetSource())
		return state.GetOut()
	} else {
		node := edge.GetDestination()
		if node == b.cfg.GetImplicitReturn() {
			return b.self.CreateEntryLattice()
		}
		state := dataFlowAnalysis_linearFlowStateOf[L](node)
		return state.GetIn()
	}
}

// DataFlowAnalysis_LinearFlowState holds the in and out states of a node.
type DataFlowAnalysis_LinearFlowState[L any] struct {
	// Go-only: language-adapter — the divergence check stops at 20002, so
	// int32 holds every observable count. queued uses the remaining padding
	// of this word rather than increasing the per-node annotation size.
	stepCount int32
	queued    bool
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
	return int(s.stepCount)
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

	finder := &dataFlowAnalysis_escapedFinder{jsScope: jsScope, escaped: escaped}

	NodeTraversal_Builder_().
		SetCompiler(compiler).
		SetCallback(finder).
		SetScopeCreator(scopeCreator).
		TraverseAtScope(jsScope)

	// TODO (simranarora) catch variables should not be considered escaped in
	// ES6. Getting rid of the catch check is causing breakages however
	for var_ := range allVarsInFn.Values() {
		if var_.GetParentNode().IsCatch() ||
			compiler.GetCodingConvention().IsExported(var_.GetName()) {
			escaped.Add(var_)
		}
	}
}

// dataFlowAnalysis_escapedFinder is the anonymous AbstractPostOrderCallback
// of computeEscaped.
type dataFlowAnalysis_escapedFinder struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	jsScope *Scope
	escaped *jutil.ObjectLinkedHashSet[*Var]
}

// Java: DataFlowAnalysis.computeEscaped(...).new AbstractPostOrderCallback() {...}.visit(NodeTraversal, Node, Node)
func (f *dataFlowAnalysis_escapedFinder) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {

	if !n.IsName() || parent.IsFunction() {
		return
	}

	name := n.GetString()
	var_ := t.GetScope().GetVar(name)

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

	referenceCfgScope := t.GetScope().GetClosestCfgRootScope()

	// Variable is referenced outside of declaration scope it is 'escaped'.
	// Non-static class field initializers are evaluated at instance creation time
	// (like closures), so variables referenced in non-static class fields must also be
	// marked as escaped.
	if referenceCfgScope != variableCfgScope ||
		DataFlowAnalysis_IsReferencedInNonStaticClassField(t.GetScope(), variableCfgScope) {
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

// dataFlowAnalysis_UniqueQueue is a work queue that holds each element at
// most once.
type dataFlowAnalysis_UniqueQueue[T comparable] struct {
	seenSet *jutil.LinkedHashSet[T]
	// membership returns an exclusively owned queue-membership flag. The
	// standalone queue uses seenSet when no annotation storage is supplied.
	membership func(T) *bool
	// Go-only: language-adapter — the Queue<T> is an ArrayDeque (deque,
	// polled at the front) without a priority, else a PriorityQueue
	// (priorityQueue).
	deque         []T
	priorityQueue *jutil.PriorityQueue[T]
}

// Java: DataFlowAnalysis.UniqueQueue(Comparator<T>)
func newDataFlowAnalysis_UniqueQueue[T comparable](priority func(a, b T) int) *dataFlowAnalysis_UniqueQueue[T] {
	return newDataFlowAnalysis_UniqueQueueWithMembership(priority, nil)
}

// Go-only: language-adapter — reuse analysis-owned per-node membership
// storage without changing the FIFO/priority queue's insertion rules.
func newDataFlowAnalysis_UniqueQueueWithMembership[T comparable](priority func(a, b T) int, membership func(T) *bool) *dataFlowAnalysis_UniqueQueue[T] {
	q := &dataFlowAnalysis_UniqueQueue[T]{membership: membership}
	if membership == nil {
		q.seenSet = jutil.NewLinkedHashSet[T]()
	}
	if priority != nil {
		q.priorityQueue = jutil.NewPriorityQueue(priority)
	}
	return q
}

// Java: DataFlowAnalysis.UniqueQueue.isEmpty()
func (q *dataFlowAnalysis_UniqueQueue[T]) IsEmpty() bool {
	if q.priorityQueue != nil {
		return q.priorityQueue.IsEmpty()
	}
	return len(q.deque) == 0
}

// Java: DataFlowAnalysis.UniqueQueue.removeFirst()
func (q *dataFlowAnalysis_UniqueQueue[T]) RemoveFirst() T {
	var t T
	removed := false
	if q.priorityQueue != nil {
		removed = !q.priorityQueue.IsEmpty()
		t = q.priorityQueue.Poll()
	} else if len(q.deque) > 0 {
		removed = true
		// ArrayDeque.poll(): null when empty.
		t = q.deque[0]
		var zero T
		q.deque[0] = zero
		q.deque = q.deque[1:]
	}
	if q.membership != nil {
		if removed {
			*q.membership(t) = false
		}
	} else {
		q.seenSet.Remove(t)
	}
	return t
}

// Java: DataFlowAnalysis.UniqueQueue.add(T)
func (q *dataFlowAnalysis_UniqueQueue[T]) Add(t T) {
	var added bool
	if q.membership != nil {
		jutil.CheckNotNull(t)
		member := q.membership(t)
		added = !*member
		*member = true
	} else {
		added = q.seenSet.Add(t)
	}
	if added {
		if q.priorityQueue != nil {
			q.priorityQueue.Add(t)
		} else {
			// ArrayDeque.add(E) is null-hostile.
			jutil.CheckNotNull(t)
			q.deque = append(q.deque, t)
		}
	}
}

// Java: DataFlowAnalysis.UniqueQueue.clear()
func (q *dataFlowAnalysis_UniqueQueue[T]) Clear() {
	if q.membership != nil {
		if q.priorityQueue != nil {
			for node := range q.priorityQueue.All() {
				*q.membership(node) = false
			}
		} else {
			for _, node := range q.deque {
				*q.membership(node) = false
			}
		}
	} else {
		q.seenSet.Clear()
	}
	q.clearStorage()
}

// Go-only: language-adapter — discard only the queue entries. The caller
// must reset membership separately, or replace all analysis annotations.
func (q *dataFlowAnalysis_UniqueQueue[T]) clearStorage() {
	if q.priorityQueue != nil {
		q.priorityQueue.Clear()
	} else {
		clear(q.deque)
		q.deque = q.deque[:0]
	}
}
