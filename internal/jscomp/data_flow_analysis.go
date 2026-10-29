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
	"github.com/paullepoulpe/jsgo/internal/graph"
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
type DataFlowAnalysis[N comparable, L graph.LatticeElement] interface {
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

	EqualsLattice(a, b L) bool

	GetDataFlowAnalysisBase() *DataFlowAnalysisBase[N, L]
}

// DataFlowAnalysisBase holds DataFlowAnalysis's fields and its base and
// final implementations.
type DataFlowAnalysisBase[N comparable, L graph.LatticeElement] struct {
	self DataFlowAnalysis[N, L]

	cfg           *ControlFlowGraph[N]
	priority      func(a, b *graph.LinkedDiGraphNode[N, ControlFlowGraph_Branch]) int
	isForward     bool
	isBranched    bool
	curLinkedNode *graph.LinkedDiGraphNode[N, ControlFlowGraph_Branch]

	// The work queue: a binary heap ordered like Java's PriorityQueue when
	// priority is non-nil, otherwise a FIFO whose head is queue[queueHead].
	queue     []*graph.LinkedDiGraphNode[N, ControlFlowGraph_Branch]
	queueHead int

	// Go-only: performance — the LinearFlowStates that Initialize sets as node
	// annotations are the elements of one slice indexed by node.Index, and
	// branched edge states are indexed by edge.Index, both reused across CFGs
	// (see dfaReuseMaxNodes).
	flowStates []DataFlowAnalysis_LinearFlowState[L]
	edgeStates []L
}

// DataFlowAnalysis_MAX_STEPS_PER_NODE is the maximum number of steps per
// individual CFG node before we assume the analysis is divergent.
//
// TODO(b/196398705): This is way too high. Find traversal ordering
// heurisitc that reduces it.
const DataFlowAnalysis_MAX_STEPS_PER_NODE = 20000

// dfaReuseMaxNodes is the largest flow state slice or queue that an analysis
// keeps for its next CFG.
const dfaReuseMaxNodes = 1024

// InitDataFlowAnalysisBase constructs a data flow analysis; self is the
// analysis being constructed. An analysis object may be constructed again
// for another CFG; it then keeps its storage.
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
func (b *DataFlowAnalysisBase[N, L]) InitDataFlowAnalysisBase(self DataFlowAnalysis[N, L], cfg *ControlFlowGraph[N]) {
	b.self = self
	b.cfg = cfg
	b.isForward = b.self.IsForward()
	b.isBranched = b.self.IsBranched()
	b.priority = cfg.GetOptionalNodeComparator(b.isForward)
	b.curLinkedNode = nil

	if b.isBranched {
		jutil.CheckState(b.isForward)
	}
}

// releaseLargeStorage drops storage that is too large to keep for reuse.
func (b *DataFlowAnalysisBase[N, L]) releaseLargeStorage() {
	if cap(b.flowStates) > dfaReuseMaxNodes {
		b.flowStates = nil
	}
	if cap(b.edgeStates) > dfaReuseMaxNodes*2 {
		b.edgeStates = nil
	}
	if cap(b.queue) > dfaReuseMaxNodes {
		b.queue = nil
		b.queueHead = 0
	}
}

// forwardCfgWorklist is a reusable, non-generic forward worklist over AST CFG
// nodes (min-heap on Priority when the CFG assigns node priorities, else FIFO).
type forwardCfgWorklist struct {
	queue       []*graph.LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch]
	queueHead   int
	usePriority bool
}

func (w *forwardCfgWorklist) reset(cfg *ControlFlowGraph[rhino.Node], numNodes int) {
	clear(w.queue)
	w.queue = w.queue[:0]
	w.queueHead = 0
	w.usePriority = cfg.hasNodePriorities
	if cap(w.queue) < numNodes {
		w.queue = make([]*graph.LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch], 0, numNodes)
	}
}

func (w *forwardCfgWorklist) releaseLarge() {
	if cap(w.queue) > dfaReuseMaxNodes {
		w.queue = nil
		w.queueHead = 0
	}
}

func (w *forwardCfgWorklist) notEmpty() bool {
	return len(w.queue) > w.queueHead
}

func (w *forwardCfgWorklist) enqueue(node *graph.LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch]) {
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
		for i > 0 {
			parent := (i - 1) / 2
			p := w.queue[parent]
			if xPrio >= p.Priority {
				break
			}
			w.queue[i] = p
			i = parent
		}
		w.queue[i] = node
	}
}

func (w *forwardCfgWorklist) dequeue() *graph.LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch] {
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

// FlowStateAt returns the LinearFlowState for the node with the given dense Index.
func (b *DataFlowAnalysisBase[N, L]) FlowStateAt(index uint32) *DataFlowAnalysis_LinearFlowState[L] {
	return &b.flowStates[index]
}

func (b *DataFlowAnalysisBase[N, L]) dfaSiftUp(i int, x *graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch]) {
	xPrio := x.Priority
	if xPrio < 0 {
		panic(jutil.NewIllegalStateExceptionMsg("priority not set"))
	}
	if b.isForward {
		for i > 0 {
			parent := (i - 1) / 2
			p := b.queue[parent]
			if xPrio >= p.Priority {
				break
			}
			b.queue[i] = p
			i = parent
		}
	} else {
		for i > 0 {
			parent := (i - 1) / 2
			p := b.queue[parent]
			if xPrio <= p.Priority {
				break
			}
			b.queue[i] = p
			i = parent
		}
	}
	b.queue[i] = x
}

func (b *DataFlowAnalysisBase[N, L]) dfaSiftDown(i int, x *graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch]) {
	n := len(b.queue)
	xPrio := x.Priority
	if b.isForward {
		for {
			pick := 2*i + 1
			if pick >= n {
				break
			}
			smaller := b.queue[pick]
			smallerPrio := smaller.Priority
			if sibling := pick + 1; sibling < n {
				sib := b.queue[sibling]
				if sibPrio := sib.Priority; smallerPrio > sibPrio {
					pick = sibling
					smaller = sib
					smallerPrio = sibPrio
				}
			}
			if xPrio <= smallerPrio {
				break
			}
			b.queue[i] = smaller
			i = pick
		}
	} else {
		for {
			pick := 2*i + 1
			if pick >= n {
				break
			}
			smaller := b.queue[pick]
			smallerPrio := smaller.Priority
			if sibling := pick + 1; sibling < n {
				sib := b.queue[sibling]
				if sibPrio := sib.Priority; smallerPrio < sibPrio {
					pick = sibling
					smaller = sib
					smallerPrio = sibPrio
				}
			}
			if xPrio >= smallerPrio {
				break
			}
			b.queue[i] = smaller
			i = pick
		}
	}
	b.queue[i] = x
}

func (b *DataFlowAnalysisBase[N, L]) enqueueNode(node *graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch], state *DataFlowAnalysis_LinearFlowState[L]) {
	if state.inQueue {
		return
	}
	if node == nil {
		panic(jutil.NewNullPointerException())
	}
	state.inQueue = true
	node.DfaInQueue = true
	b.queue = append(b.queue, node)
	if b.priority != nil {
		b.dfaSiftUp(len(b.queue)-1, node)
	}
}

func (b *DataFlowAnalysisBase[N, L]) dequeueNode() (*graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch], *DataFlowAnalysis_LinearFlowState[L]) {
	var curNode *graph.LinkedDirectedGraph_LinkedDiGraphNode[N, ControlFlowGraph_Branch]
	if b.priority != nil {
		n := len(b.queue)
		curNode = b.queue[0]
		last := b.queue[n-1]
		b.queue[n-1] = nil
		b.queue = b.queue[:n-1]
		if n > 1 {
			b.dfaSiftDown(0, last)
		}
	} else {
		curNode = b.queue[b.queueHead]
		b.queue[b.queueHead] = nil
		b.queueHead++
		if b.queueHead == len(b.queue) {
			b.queue = b.queue[:0]
			b.queueHead = 0
		}
	}
	curNode.DfaInQueue = false
	curState := &b.flowStates[curNode.Index]
	curState.inQueue = false
	return curNode, curState
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
func (b *DataFlowAnalysisBase[N, L]) GetCfg() *ControlFlowGraph[N] {
	return b.cfg
}

// Java: DataFlowAnalysis.isBranched()
func (b *DataFlowAnalysisBase[N, L]) IsBranched() bool {
	return false
}

func (b *DataFlowAnalysisBase[N, L]) EqualsLattice(a, bL L) bool {
	return jutil.ObjectsEquals(a, bL)
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
	isForward := b.isForward
	entry := b.cfg.entry
	implicitReturn := b.cfg.implicitReturn
	for len(b.queue) > b.queueHead {
		curNode, curState := b.dequeueNode()
		stepCount := curState.stepCount
		curState.stepCount++
		if stepCount > DataFlowAnalysis_MAX_STEPS_PER_NODE {
			panic(jutil.NewIllegalStateExceptionMsg("Dataflow analysis appears to diverge around: " + jutil.ToJavaString(curNode)))
		}

		b.joinInputsLinked(curNode, curState, entry, implicitReturn)
		if b.flowLinked(curNode, curState) {
			// If there is a change in the current node, we want to grab the
			// list of nodes that this node affects.
			if isForward {
				for _, edge := range curNode.OutEdgeList.Slice() {
					nextNode := edge.DestNode
					if !nextNode.DfaInQueue && nextNode != implicitReturn {
						b.enqueueNode(nextNode, &b.flowStates[nextNode.Index])
					}
				}
			} else {
				for _, edge := range curNode.InEdgeList.Slice() {
					nextNode := edge.SourceNode
					if !nextNode.DfaInQueue && nextNode != implicitReturn {
						b.enqueueNode(nextNode, &b.flowStates[nextNode.Index])
					}
				}
			}
		}
	}
	b.curLinkedNode = nil
	if isForward && implicitReturn != nil {
		b.joinInputsLinked(implicitReturn, &b.flowStates[implicitReturn.Index], entry, implicitReturn)
	}
	if b.isBranched {
		for node := range b.cfg.Nodes.Values() {
			for _, edge := range node.OutEdgeList.Slice() {
				edge.SetAnnotation(b.edgeStates[edge.Index])
			}
		}
	}
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

// dataFlowAnalysis_latticeOf is the unchecked cast `L l =
// edge.getAnnotation()`.
//
// Go-only: language-adapter — Java's generic getAnnotation() casts at the
// call site; a null annotation stays null.
func dataFlowAnalysis_latticeOf[L any](a graph.Annotatable) L {
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
	clear(b.queue)
	b.queue = b.queue[:0]
	b.queueHead = 0

	implicitReturn := b.cfg.implicitReturn
	cfgNodes := &b.cfg.Nodes
	numNodes := cfgNodes.Size()
	nodeCap := b.cfg.NodeIndexCap()
	if cap(b.queue) < numNodes {
		b.queue = make([]*graph.LinkedDiGraphNode[N, ControlFlowGraph_Branch], 0, numNodes)
	}
	if cap(b.flowStates) < nodeCap {
		b.flowStates = make([]DataFlowAnalysis_LinearFlowState[L], nodeCap)
	} else {
		clear(b.flowStates)
		b.flowStates = b.flowStates[:nodeCap]
	}
	for node := range cfgNodes.Values() {
		node.DfaInQueue = false
		// new LinearFlowState<>(createInitialEstimateLattice(), createInitialEstimateLattice())
		in, out := b.self.CreateInitialEstimateLattice(), b.self.CreateInitialEstimateLattice()
		state := &b.flowStates[node.Index]
		state.in = jutil.CheckNotNull(in)
		state.out = jutil.CheckNotNull(out)
		node.SetAnnotation(state)
		if node != implicitReturn {
			b.enqueueNode(node, state)
		}
	}
	if b.isBranched {
		edgeCap := b.cfg.EdgeIndexCap()
		if cap(b.edgeStates) < edgeCap {
			b.edgeStates = make([]L, edgeCap)
		} else {
			clear(b.edgeStates)
			b.edgeStates = b.edgeStates[:edgeCap]
		}
		for node := range cfgNodes.Values() {
			for _, edge := range node.OutEdgeList.Slice() {
				b.edgeStates[edge.Index] = b.self.CreateInitialEstimateLattice()
			}
		}
	}
}

// Flow performs a single flow through a node. Returns true if the flow
// state differs from the previous state.
//
// Java: DataFlowAnalysis.flow(DiGraphNode<N, Branch>)
func (b *DataFlowAnalysisBase[N, L]) Flow(node *graph.LinkedDiGraphNode[N, ControlFlowGraph_Branch]) bool {
	if int(node.Index) < len(b.flowStates) && (!b.isBranched || len(b.edgeStates) > 0) {
		return b.flowLinked(node, &b.flowStates[node.Index])
	}
	return b.flowWithState(node, dataFlowAnalysis_linearFlowStateOf[L](node))
}

func (b *DataFlowAnalysisBase[N, L]) flowLinked(node *graph.LinkedDiGraphNode[N, ControlFlowGraph_Branch], state *DataFlowAnalysis_LinearFlowState[L]) bool {
	b.curLinkedNode = node
	if b.isForward {
		outBefore := state.out
		outAfter := b.self.FlowThrough(node.Value, state.in)
		state.SetOut(outAfter)
		changed := !b.self.EqualsLattice(outBefore, outAfter)

		if b.isBranched {
			brancher := b.self.CreateFlowBrancher(node.Value, outAfter)
			for _, le := range node.OutEdgeList.Slice() {
				outBranchBefore := b.edgeStates[le.Index]
				outBranchAfter := brancher(le.Value)
				b.edgeStates[le.Index] = outBranchAfter
				if !changed {
					changed = !b.self.EqualsLattice(outBranchBefore, outBranchAfter)
				}
			}
		}

		return changed
	} else {
		inBefore := state.in
		inAfter := b.self.FlowThrough(node.Value, state.out)
		state.SetIn(inAfter)
		return !b.self.EqualsLattice(inBefore, inAfter)
	}
}

func (b *DataFlowAnalysisBase[N, L]) flowWithState(node *graph.LinkedDiGraphNode[N, ControlFlowGraph_Branch], state *DataFlowAnalysis_LinearFlowState[L]) bool {
	b.curLinkedNode = node
	if b.isForward {
		outBefore := state.GetOut()
		outAfter := b.self.FlowThrough(node.Value, state.GetIn())
		state.SetOut(outAfter)
		changed := !b.self.EqualsLattice(outBefore, outAfter)

		if b.isBranched {
			brancher := b.self.CreateFlowBrancher(node.Value, outAfter)
			for _, outEdge := range node.OutEdgeList.Slice() {
				outBranchBefore := dataFlowAnalysis_latticeOf[L](outEdge)
				outBranchAfter := brancher(outEdge.Value)
				outEdge.SetAnnotation(outBranchAfter)
				if !changed {
					changed = !b.self.EqualsLattice(outBranchBefore, outBranchAfter)
				}
			}
		}

		return changed
	} else {
		inBefore := state.GetIn()
		inAfter := b.self.FlowThrough(node.Value, state.GetOut())
		state.SetIn(inAfter)
		return !b.self.EqualsLattice(inBefore, inAfter)
	}
}

// JoinInputs computes the new flow state at a given node's entry by merging
// the output (input) lattice of the node's predecessor (successor).
//
// node is the node to compute new join.
//
// Java: DataFlowAnalysis.joinInputs(DiGraphNode<N, Branch>)
func (b *DataFlowAnalysisBase[N, L]) JoinInputs(node *graph.LinkedDiGraphNode[N, ControlFlowGraph_Branch]) {
	if int(node.Index) < len(b.flowStates) && (!b.isBranched || len(b.edgeStates) > 0) {
		b.joinInputsLinked(node, &b.flowStates[node.Index], b.cfg.entry, b.cfg.implicitReturn)
		return
	}
	b.joinInputsWithState(node, dataFlowAnalysis_linearFlowStateOf[L](node))
}

func (b *DataFlowAnalysisBase[N, L]) joinInputsLinked(
	node *graph.LinkedDiGraphNode[N, ControlFlowGraph_Branch],
	state *DataFlowAnalysis_LinearFlowState[L],
	entry *graph.LinkedDiGraphNode[N, ControlFlowGraph_Branch],
	implicitReturn *graph.LinkedDiGraphNode[N, ControlFlowGraph_Branch],
) {
	if b.isForward && entry == node {
		state.SetIn(b.self.CreateEntryLattice())
		return
	}

	var inEdges []*graph.LinkedDiGraphEdge[N, ControlFlowGraph_Branch]
	if b.isForward {
		inEdges = node.InEdgeList.Slice()
	} else {
		inEdges = node.OutEdgeList.Slice()
	}

	var result L
	switch len(inEdges) {
	case 0:
		return
	case 1:
		result = b.getInputFromLinkedEdge(inEdges[0], implicitReturn)
	default:
		joiner := b.self.CreateFlowJoiner()
		for _, inEdge := range inEdges {
			joiner.JoinFlow(b.getInputFromLinkedEdge(inEdge, implicitReturn))
		}
		result = joiner.Finish()
	}

	if b.isForward {
		state.SetIn(result)
	} else {
		state.SetOut(result)
	}
}

func (b *DataFlowAnalysisBase[N, L]) getInputFromLinkedEdge(
	edge *graph.LinkedDiGraphEdge[N, ControlFlowGraph_Branch],
	implicitReturn *graph.LinkedDiGraphNode[N, ControlFlowGraph_Branch],
) L {
	if b.isBranched {
		return b.edgeStates[edge.Index]
	} else if b.isForward {
		return b.flowStates[edge.SourceNode.Index].out
	} else {
		node := edge.DestNode
		if node == implicitReturn {
			return b.self.CreateEntryLattice()
		}
		return b.flowStates[node.Index].in
	}
}

func (b *DataFlowAnalysisBase[N, L]) joinInputsWithState(node *graph.LinkedDiGraphNode[N, ControlFlowGraph_Branch], state *DataFlowAnalysis_LinearFlowState[L]) {
	if b.isForward && b.cfg.entry == node {
		state.SetIn(b.self.CreateEntryLattice())
		return
	}

	var inEdges []*graph.LinkedDiGraphEdge[N, ControlFlowGraph_Branch]
	if b.isForward {
		inEdges = node.InEdgeList.Slice()
	} else {
		inEdges = node.OutEdgeList.Slice()
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

	if b.isForward {
		state.SetIn(result)
	} else {
		state.SetOut(result)
	}
}

// Java: DataFlowAnalysis.getInputFromEdge(DiGraphEdge<N, Branch>)
func (b *DataFlowAnalysisBase[N, L]) GetInputFromEdge(edge *graph.LinkedDiGraphEdge[N, ControlFlowGraph_Branch]) L {
	if b.isBranched {
		return dataFlowAnalysis_latticeOf[L](edge)
	} else if b.isForward {
		state := dataFlowAnalysis_linearFlowStateOf[L](edge.SourceNode)
		return state.GetOut()
	} else {
		node := edge.DestNode
		if node == b.cfg.implicitReturn {
			return b.self.CreateEntryLattice()
		}
		state := dataFlowAnalysis_linearFlowStateOf[L](node)
		return state.GetIn()
	}
}

// DataFlowAnalysis_LinearFlowState holds the in and out states of a node.
type DataFlowAnalysis_LinearFlowState[L any] struct {
	stepCount int
	in        L
	out       L
	inQueue   bool
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
	fnName := fnNameNode.GetString()
	fnNameShadowed := false
	if fnName != "" && allVarsInFn != nil && jsScope.HasOwnSlot(fnName) {
		if v := allVarsInFn.Get(fnName); v == nil || v.scope != jsScope {
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

// dataFlowAnalysis_escapedFinder is the anonymous AbstractPostOrderCallback
// of computeEscaped.
type dataFlowAnalysis_escapedFinder struct {
	jsScope              *Scope
	nonEscapingDepth     int
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
	if !f.hasDuplicateVarNames && f.escapedStartedEmpty && f.escaped.Size() >= f.allVarsInFn.Size() {
		return false
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

	name := n.GetString()
	if f.allVarsInFn != nil {
		targetVar := f.allVarsInFn.Get(name)
		if targetVar == nil {
			return
		}
		if !f.hasDuplicateVarNames && f.escaped.Contains(targetVar) {
			return
		}
	}

	var_ := t.LookupVar(name)
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
	f.escaped.Add(var_)
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
