// Port of com/google/javascript/jscomp/DataFlowAnalysis.java
// @ closure-compiler v20260415
//
// DataFlowAnalysis<N, L extends LatticeElement> is the Java framework for
// writing static program analysis. This Go port uses generics (Go 1.18+) to
// mirror the Java parameterization.
//
// Architecture:
//   - LatticeElement: marker interface for lattice values.
//   - FlowJoiner[L]: reducer that merges incoming flow states at join points.
//   - FlowBrancher[L]: produces distinct output states per outgoing CFG edge.
//   - LinearFlowState[L]: (in, out) annotation stored at each CFG node.
//   - GenericDFA[N, L]: the generic fixed-point engine.
//   - DataFlowAnalysis (non-generic): concrete forward analysis over FlowScope,
//     kept for backward-compatibility with TypeInferencePass callers.

package jscomp

import (
	"fmt"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ---------------------------------------------------------------------------
// Lattice interfaces — mirrors LatticeElement, FlowJoiner, FlowBrancher
// ---------------------------------------------------------------------------

// LatticeElement is the marker interface for values in the analysis lattice.
// Any type used as the lattice type L must implement LatticeElement.
// Mirrors com.google.javascript.jscomp.graph.LatticeElement in Java.
type LatticeElement[L any] interface {
	// Equals reports whether this lattice value equals other.
	// Used by the fixed-point test to detect when an output has changed.
	// Mirrors LatticeElement.equals() (inherited from Object) in Java.
	Equals(other L) bool
}

// FlowJoiner[L] accumulates flow states from distinct incoming edges and
// produces a single merged state for the join point.
// Mirrors DataFlowAnalysis.FlowJoiner<L> in Java.
type FlowJoiner[L any] interface {
	// JoinFlow feeds one incoming flow state into this joiner.
	// Called once per incoming edge.
	JoinFlow(input L)
	// Finish returns the accumulated join of all inputs fed so far.
	Finish() L
}

// FlowBrancher[L] produces a distinct output state for each outgoing CFG
// edge at a branching node. Created by GenericDFABase.CreateFlowBrancher.
// Mirrors DataFlowAnalysis.FlowBrancher<L> in Java.
type FlowBrancher[L any] interface {
	// BranchFlow returns the output state for the given branch.
	BranchFlow(b Branch) L
}

// ---------------------------------------------------------------------------
// LinearFlowState — (in, out, stepCount) per CFG node
// ---------------------------------------------------------------------------

// LinearFlowState stores the IN and OUT lattice values at a single CFG node,
// together with a step counter for divergence detection.
// Mirrors DataFlowAnalysis.LinearFlowState<L> in Java.
type LinearFlowState[L any] struct {
	In        L
	Out       L
	StepCount int
}

// GetIn returns the IN lattice value.
// Mirrors LinearFlowState.getIn() in Java.
func (s *LinearFlowState[L]) GetIn() L { return s.In }

// SetIn sets the IN lattice value.
// Mirrors LinearFlowState.setIn() in Java.
func (s *LinearFlowState[L]) SetIn(v L) { s.In = v }

// GetOut returns the OUT lattice value.
// Mirrors LinearFlowState.getOut() in Java.
func (s *LinearFlowState[L]) GetOut() L { return s.Out }

// SetOut sets the OUT lattice value.
// Mirrors LinearFlowState.setOut() in Java.
func (s *LinearFlowState[L]) SetOut(v L) { s.Out = v }

// GetStepCount returns the iteration step count for this node.
// Mirrors LinearFlowState.getStepCount() in Java.
func (s *LinearFlowState[L]) GetStepCount() int { return s.StepCount }

// String returns a human-readable representation of this flow state.
// Mirrors LinearFlowState.toString() in Java.
func (s *LinearFlowState[L]) String() string {
	return fmt.Sprintf("IN: %v OUT: %v", s.In, s.Out)
}

// ---------------------------------------------------------------------------
// GenericDFABase — abstract base for concrete analyses
// ---------------------------------------------------------------------------

// MaxStepsPerNode is the maximum number of times a single node may be visited
// before the analysis is presumed to diverge.
// Mirrors DataFlowAnalysis.MAX_STEPS_PER_NODE in Java.
const MaxStepsPerNode = 20000

// GenericDFABase is the callback interface that concrete analyses must
// satisfy to be driven by GenericDFA. It corresponds to the abstract methods
// of Java's DataFlowAnalysis<N, L>.
//
// N is the node-value type stored in each CfgNode (use *rhino.Node or any
// user-defined struct).
// L is the lattice type.
type GenericDFABase[N any, L LatticeElement[L]] interface {
	// IsForward reports whether this is a forward (true) or backward (false)
	// analysis. Mirrors DataFlowAnalysis.isForward() in Java.
	IsForward() bool

	// IsBranched reports whether CreateFlowBrancher should be used.
	// Return true when ON_TRUE/ON_FALSE edges need distinct lattice values.
	// Mirrors DataFlowAnalysis.isBranched() in Java.
	IsBranched() bool

	// CreateFlowJoiner returns a fresh joiner for one fixed-point step.
	// Mirrors DataFlowAnalysis.createFlowJoiner() in Java.
	CreateFlowJoiner() FlowJoiner[L]

	// CreateFlowBrancher returns a brancher for the given node's output state.
	// Called only when IsBranched() returns true.
	// Mirrors DataFlowAnalysis.createFlowBrancher(N, L) in Java.
	CreateFlowBrancher(node N, output L) FlowBrancher[L]

	// FlowThrough computes the output state for a node given its input state.
	// Mirrors DataFlowAnalysis.flowThrough(N, L) in Java.
	FlowThrough(node N, input L) L

	// CreateEntryLattice returns the lattice value for the entry node.
	// Mirrors DataFlowAnalysis.createEntryLattice() in Java.
	CreateEntryLattice() L

	// CreateInitialEstimateLattice returns the bottom (initial estimate) value.
	// Mirrors DataFlowAnalysis.createInitialEstimateLattice() in Java.
	CreateInitialEstimateLattice() L
}

// ---------------------------------------------------------------------------
// GenericDFA — the generic fixed-point engine
// ---------------------------------------------------------------------------

// GenericDFA[N, L] drives a concrete DataFlowAnalysis over a ControlFlowGraph
// whose nodes carry values of type N. After Analyze() returns, each node's
// LinearFlowState is accessible via State.
//
// Usage:
//
//	type MyAnalysis struct { ... }
//	func (a *MyAnalysis) IsForward() bool { return true }
//	...
//	dfa := NewGenericDFA[*rhino.Node, MyLattice](cfg, nodeValFunc, &MyAnalysis{})
//	dfa.Analyze()
//	state := dfa.State[someNode]
//
// Mirrors DataFlowAnalysis<N, L> in Java.
type GenericDFA[N any, L LatticeElement[L]] struct {
	// Cfg is the control-flow graph being analysed.
	Cfg *ControlFlowGraph
	// nodeVal extracts the N-typed value from a CfgNode. For CFGs whose
	// nodes wrap *rhino.Node, use `func(n *CfgNode) *rhino.Node { return n.AST }`.
	nodeVal func(*CfgNode) N
	// impl is the concrete analysis implementation.
	impl GenericDFABase[N, L]
	// State maps each CfgNode to its LinearFlowState after Analyze.
	State map[*CfgNode]*LinearFlowState[L]
	// edgeAnnotation stores per-edge lattice values when impl.IsBranched().
	edgeAnnotation map[edgeKey]L
}

// NewGenericDFA constructs a GenericDFA.
// nodeVal must return the N-typed payload stored in a CfgNode.
// Mirrors the DataFlowAnalysis(ControlFlowGraph) constructor in Java.
// Go-only: language-adapter Go generics require an explicit constructor; Java uses the class constructor.
func NewGenericDFA[N any, L LatticeElement[L]](
	cfg *ControlFlowGraph,
	nodeVal func(*CfgNode) N,
	impl GenericDFABase[N, L],
) *GenericDFA[N, L] {
	return &GenericDFA[N, L]{
		Cfg:            cfg,
		nodeVal:        nodeVal,
		impl:           impl,
		State:          make(map[*CfgNode]*LinearFlowState[L]),
		edgeAnnotation: make(map[edgeKey]L),
	}
}

// initialize seeds every node with the initial-estimate lattice and fills
// the work-queue with all non-exit nodes.
// Mirrors DataFlowAnalysis.initialize() in Java.
func (d *GenericDFA[N, L]) initialize(queue *uniqueQueue[*CfgNode]) {
	bottom := d.impl.CreateInitialEstimateLattice()
	for _, n := range d.Cfg.Nodes {
		d.State[n] = &LinearFlowState[L]{In: bottom, Out: bottom}
		if n != d.Cfg.Exit {
			queue.add(n)
		}
	}
	if d.impl.IsBranched() {
		for _, n := range d.Cfg.Nodes {
			for _, e := range n.Successors {
				key := edgeKey{from: n, to: e.To, branch: e.Branch}
				d.edgeAnnotation[key] = bottom
			}
		}
	}
}

// joinInputs computes the merged IN (forward) or OUT (backward) from all
// incoming edges and stores it on the node's state.
// Mirrors DataFlowAnalysis.joinInputs() in Java.
func (d *GenericDFA[N, L]) joinInputs(n *CfgNode) {
	state := d.State[n]
	if d.impl.IsForward() && d.Cfg.Entry == n {
		state.In = d.impl.CreateEntryLattice()
		return
	}
	var inEdges []CfgEdge
	if d.impl.IsForward() {
		inEdges = n.inEdges()
	} else {
		// backward: "inputs" are the successor edges
		inEdges = n.Successors
	}
	switch len(inEdges) {
	case 0:
		return
	case 1:
		e := inEdges[0]
		if d.impl.IsForward() {
			state.In = d.getInputFromEdge(e, e.To)
		} else {
			state.Out = d.getInputFromEdge(e, e.To)
		}
	default:
		joiner := d.impl.CreateFlowJoiner()
		for _, e := range inEdges {
			var src *CfgNode
			if d.impl.IsForward() {
				src = e.To // inEdge's source (we stored them as src→n)
			} else {
				src = e.To
			}
			joiner.JoinFlow(d.getInputFromEdge(e, src))
		}
		result := joiner.Finish()
		if d.impl.IsForward() {
			state.In = result
		} else {
			state.Out = result
		}
	}
}

// getInputFromEdge retrieves the lattice value to feed into a join at n from
// the given edge. Mirrors DataFlowAnalysis.getInputFromEdge() in Java.
func (d *GenericDFA[N, L]) getInputFromEdge(e CfgEdge, src *CfgNode) L {
	if d.impl.IsBranched() {
		// branched forward: edge annotation carries the split value
		key := edgeKey{from: src, to: e.To, branch: e.Branch}
		return d.edgeAnnotation[key]
	}
	if d.impl.IsForward() {
		if srcState, ok := d.State[src]; ok {
			return srcState.Out
		}
		return d.impl.CreateInitialEstimateLattice()
	}
	// backward: input at a node is its IN (which will become the successor's OUT)
	if e.To == d.Cfg.Exit {
		return d.impl.CreateEntryLattice()
	}
	if destState, ok := d.State[e.To]; ok {
		return destState.In
	}
	return d.impl.CreateInitialEstimateLattice()
}

// flow executes one transfer at `n` and returns true if the state changed.
// Mirrors DataFlowAnalysis.flow() in Java.
func (d *GenericDFA[N, L]) flow(n *CfgNode) bool {
	state := d.State[n]
	nVal := d.nodeVal(n)
	if d.impl.IsForward() {
		outBefore := state.Out
		state.Out = d.impl.FlowThrough(nVal, state.In)
		changed := !outBefore.Equals(state.Out)
		if d.impl.IsBranched() {
			brancher := d.impl.CreateFlowBrancher(nVal, state.Out)
			for _, e := range n.Successors {
				key := edgeKey{from: n, to: e.To, branch: e.Branch}
				outBranchBefore := d.edgeAnnotation[key]
				d.edgeAnnotation[key] = brancher.BranchFlow(e.Branch)
				if !changed {
					changed = !outBranchBefore.Equals(d.edgeAnnotation[key])
				}
			}
		}
		return changed
	}
	// backward
	inBefore := state.In
	state.In = d.impl.FlowThrough(nVal, state.Out)
	return !inBefore.Equals(state.In)
}

// Analyze runs the fixed-point iteration to convergence.
// After Analyze returns, d.State holds the final IN/OUT for every node.
// Mirrors DataFlowAnalysis.analyze() in Java.
func (d *GenericDFA[N, L]) Analyze() {
	queue := newUniqueQueue[*CfgNode]()
	d.initialize(queue)
	for !queue.isEmpty() {
		cur := queue.removeFirst()
		state := d.State[cur]
		state.StepCount++
		if state.StepCount > MaxStepsPerNode {
			panic(fmt.Sprintf("dataflow analysis appears to diverge around: %v", cur))
		}
		d.joinInputs(cur)
		if d.flow(cur) {
			// Propagate to successors (forward) or predecessors (backward).
			var nextNodes []*CfgNode
			if d.impl.IsForward() {
				for _, e := range cur.Successors {
					nextNodes = append(nextNodes, e.To)
				}
			} else {
				nextNodes = cur.Predecessors
			}
			for _, next := range nextNodes {
				if next != d.Cfg.Exit {
					queue.add(next)
				}
			}
		}
	}
	// Final join at the implicit return (forward only).
	if d.impl.IsForward() {
		d.joinInputs(d.Cfg.Exit)
	}
}

// IsBranched reports whether this analysis uses branched flow (per-edge states).
// The default is false; analyses that need ON_TRUE/ON_FALSE distinctions should
// override IsBranched via GenericDFABase to return true.
// Mirrors DataFlowAnalysis.isBranched() in Java.
func (d *GenericDFA[N, L]) IsBranched() bool {
	return d.impl.IsBranched()
}

// IsForward reports whether this is a forward analysis.
// Mirrors DataFlowAnalysis.isForward() in Java.
func (d *GenericDFA[N, L]) IsForward() bool {
	return d.impl.IsForward()
}

// CreateEntryLattice returns the entry lattice value.
// Mirrors DataFlowAnalysis.createEntryLattice() in Java.
func (d *GenericDFA[N, L]) CreateEntryLattice() L {
	return d.impl.CreateEntryLattice()
}

// CreateInitialEstimateLattice returns the bottom lattice value.
// Mirrors DataFlowAnalysis.createInitialEstimateLattice() in Java.
func (d *GenericDFA[N, L]) CreateInitialEstimateLattice() L {
	return d.impl.CreateInitialEstimateLattice()
}

// FlowThrough computes the output lattice for node nVal given input.
// Mirrors DataFlowAnalysis.flowThrough(N, L) in Java.
func (d *GenericDFA[N, L]) FlowThrough(nVal N, input L) L {
	return d.impl.FlowThrough(nVal, input)
}

// Join merges two lattice values using this analysis's FlowJoiner.
// Mirrors DataFlowAnalysis.join(L, L) in Java.
func (d *GenericDFA[N, L]) Join(latticeA, latticeB L) L {
	joiner := d.impl.CreateFlowJoiner()
	joiner.JoinFlow(latticeA)
	joiner.JoinFlow(latticeB)
	return joiner.Finish()
}

// GetCfg returns the control flow graph being analysed.
// Mirrors DataFlowAnalysis.getCfg() in Java.
func (d *GenericDFA[N, L]) GetCfg() *ControlFlowGraph {
	return d.Cfg
}

// CreateFlowJoiner returns a fresh joiner for an analysis step.
// Mirrors DataFlowAnalysis.createFlowJoiner() in Java.
func (d *GenericDFA[N, L]) CreateFlowJoiner() FlowJoiner[L] {
	return d.impl.CreateFlowJoiner()
}

// CreateFlowBrancher returns a brancher for the given node's output state.
// Mirrors DataFlowAnalysis.createFlowBrancher(N, L) in Java.
func (d *GenericDFA[N, L]) CreateFlowBrancher(node N, output L) FlowBrancher[L] {
	return d.impl.CreateFlowBrancher(node, output)
}

// GetFlowState returns the LinearFlowState for node n, or nil if n is not in
// the CFG or Analyze has not been called.
// Mirrors DiGraphNode.getAnnotation() usage in Java.
// Go-only: language-adapter Java uses DiGraphNode.getAnnotation(); Go uses a typed map lookup.
func (d *GenericDFA[N, L]) GetFlowState(n *CfgNode) *LinearFlowState[L] {
	return d.State[n]
}

// ---------------------------------------------------------------------------
// uniqueQueue — deduplicating FIFO queue (mirrors UniqueQueue<T> in Java)
// ---------------------------------------------------------------------------

// uniqueQueue is a FIFO queue that ignores duplicate inserts while the
// element is already present in the queue.
// Mirrors the private UniqueQueue<T> in DataFlowAnalysis.java.
// Go-only: language-adapter Java's UniqueQueue uses PriorityQueue+LinkedHashSet; Go uses a slice+map for simplicity.
type uniqueQueue[T comparable] struct {
	queue  []T
	inQ    map[T]bool
}

// Go-only: language-adapter Go generics need an explicit factory; Java uses constructor.
func newUniqueQueue[T comparable]() *uniqueQueue[T] {
	return &uniqueQueue[T]{inQ: make(map[T]bool)}
}

func (q *uniqueQueue[T]) isEmpty() bool { return len(q.queue) == 0 }

func (q *uniqueQueue[T]) removeFirst() T {
	t := q.queue[0]
	q.queue = q.queue[1:]
	delete(q.inQ, t)
	return t
}

func (q *uniqueQueue[T]) add(t T) {
	if !q.inQ[t] {
		q.inQ[t] = true
		q.queue = append(q.queue, t)
	}
}

func (q *uniqueQueue[T]) clear() {
	q.queue = q.queue[:0]
	for k := range q.inQ {
		delete(q.inQ, k)
	}
}

// ---------------------------------------------------------------------------
// inEdges helper on CfgNode
// ---------------------------------------------------------------------------

// inEdges returns the set of incoming edges to n as CfgEdge values.
// The CfgEdge.To field is the source node (inverted so the join loop
// can call getInputFromEdge uniformly).
// Used only by joinInputs when isForward==true.
// Go-only: language-adapter Java uses DiGraph.getInEdges() built into the graph; Go reconstructs from Predecessors.
func (n *CfgNode) inEdges() []CfgEdge {
	edges := make([]CfgEdge, 0, len(n.Predecessors))
	for _, pred := range n.Predecessors {
		for _, e := range pred.Successors {
			if e.To == n {
				edges = append(edges, CfgEdge{To: pred, Branch: e.Branch})
			}
		}
	}
	return edges
}

// ---------------------------------------------------------------------------
// computeEscaped — escaped-variable computation
// ---------------------------------------------------------------------------

// EscapedResult records which Vars escape the closest CFG-root scope.
// Mirrors the escaped set populated by DataFlowAnalysis.computeEscaped() in Java.
type EscapedResult struct {
	// Escaped contains each Var that is referenced outside the scope in
	// which it is declared (within a nested function) or is explicitly
	// exported by name convention.
	Escaped map[string]bool
}

// ComputeEscaped finds variables declared in jsScope that escape into
// nested function bodies. jsScope must be a function scope.
//
// A variable is considered escaped when:
//  1. It is referenced inside a nested function (closure capture), OR
//  2. It is declared as a catch parameter (any JS engine may expose them), OR
//  3. Its name is exported by the compiler's coding convention.
//
// Mirrors DataFlowAnalysis.computeEscaped() in Java.
func ComputeEscaped(
	rootAST *rhino.Node,
	allVarsInFn map[string]*rhino.Node, // name → nameNode
) *EscapedResult {
	result := &EscapedResult{Escaped: make(map[string]bool)}
	if rootAST == nil {
		return result
	}
	// Walk the AST, collecting name-references inside nested functions.
	// The outer function scope is identified as any FUNCTION directly under root.
	var walk func(n *rhino.Node, inNestedFn bool)
	walk = func(n *rhino.Node, inNestedFn bool) {
		if n == nil {
			return
		}
		tok := n.Token()
		// A NAME node (not the function's own name position) references a variable.
		if tok == rhino.NAME {
			name := n.String_()
			if inNestedFn {
				if _, declared := allVarsInFn[name]; declared {
					result.Escaped[name] = true
				}
			}
		}
		// Entering a nested FUNCTION means we're now in a nested scope.
		for _, c := range n.Children() {
			childInNested := inNestedFn
			if c.Token() == rhino.FUNCTION {
				childInNested = true
			}
			walk(c, childInNested)
		}
	}
	walk(rootAST, false)
	// Catch variables always escape (any engine may expose them).
	var markCatch func(n *rhino.Node)
	markCatch = func(n *rhino.Node) {
		if n == nil {
			return
		}
		if n.Token() == rhino.CATCH {
			kids := n.Children()
			if len(kids) > 0 && kids[0].Token() == rhino.NAME {
				result.Escaped[kids[0].String_()] = true
			}
		}
		for _, c := range n.Children() {
			markCatch(c)
		}
	}
	markCatch(rootAST)
	return result
}

// ---------------------------------------------------------------------------
// Forward DataFlowAnalysis — concrete FlowScope-based (kept for callers)
// ---------------------------------------------------------------------------

// TransferFunction computes a node's OUT FlowScope given its IN and
// (for branching nodes) the branch kind along which control reached
// it. Returns multiple outputs when a node refines differently along
// BranchOnTrue vs BranchOnFalse successors — the key is the Branch,
// the value is the FlowScope to attach to that outgoing edge.
type TransferFunction func(node *CfgNode, in *FlowScope) map[Branch]*FlowScope

// DataFlowAnalysis is the forward-worklist runner. It stores the
// computed IN scope for each CfgNode keyed by pointer — callers
// consume this after Run returns.
type DataFlowAnalysis struct {
	Cfg      *ControlFlowGraph
	In       map[*CfgNode]*FlowScope
	edgeOut  map[edgeKey]*FlowScope
	transfer TransferFunction
	initial  *FlowScope
}

// edgeKey identifies a directed CFG edge (from→to with branch label).
// Go-only: language-adapter Java uses DiGraphEdge identity; Go uses a struct key for map lookups.
type edgeKey struct {
	from   *CfgNode
	to     *CfgNode
	branch Branch
}

// NewDataFlowAnalysis wires a DataFlowAnalysis over `cfg` with the
// given transfer function. `initial` is the IN scope attached to the
// entry node (typically an empty FlowScope).
// Go-only: language-adapter Go needs a constructor function; Java DataFlowAnalysis uses subclassing.
func NewDataFlowAnalysis(cfg *ControlFlowGraph, initial *FlowScope, transfer TransferFunction) *DataFlowAnalysis {
	return &DataFlowAnalysis{
		Cfg:      cfg,
		In:       map[*CfgNode]*FlowScope{},
		edgeOut:  map[edgeKey]*FlowScope{},
		transfer: transfer,
		initial:  initial,
	}
}

// Run executes the forward-worklist fixed-point loop. Termination is
// guaranteed because FlowScope entries only widen (each name either
// keeps its refined type or joins to a union that subsumes it), and
// NewUnion deduplicates by name — so the number of distinct union
// shapes per name is bounded by the type lattice.
// Go-only: language-adapter Mirrors DataFlowAnalysis.analyze(); renamed Run to avoid collision with the generic Analyze.
func (d *DataFlowAnalysis) Run() {
	d.In[d.Cfg.Entry] = d.initial
	worklist := []*CfgNode{d.Cfg.Entry}
	// Each visit bound is a safety belt — the lattice itself is finite so
	// we shouldn't need it, but wide unions produced by a pathological
	// transfer could otherwise loop.
	const maxIterations = 20000
	iterations := 0
	inWorklist := map[*CfgNode]bool{d.Cfg.Entry: true}
	for len(worklist) > 0 && iterations < maxIterations {
		iterations++
		n := worklist[0]
		worklist = worklist[1:]
		delete(inWorklist, n)
		in := d.In[n]
		// Recompute IN from predecessors' edge-out. On the first
		// visit for the entry, we skip this step.
		if n != d.Cfg.Entry {
			var merged *FlowScope
			for _, p := range n.Predecessors {
				for _, e := range p.Successors {
					if e.To != n {
						continue
					}
					edgeIn := d.edgeOut[edgeKey{from: p, to: n, branch: e.Branch}]
					if edgeIn == nil {
						continue
					}
					if merged == nil {
						merged = edgeIn
					} else {
						merged = JoinFlowScopes(merged, edgeIn)
					}
				}
			}
			if merged == nil {
				merged = EmptyFlowScope()
			}
			if in != nil && in.Equal(merged) {
				continue
			}
			in = merged
			d.In[n] = in
		}
		outs := d.transfer(n, in)
		// Default: if the transfer didn't produce a per-branch
		// distinction, the same OUT scope flows to every successor
		// regardless of branch label.
		defaultOut, hasDefault := outs[BranchUnconditional]
		for _, e := range n.Successors {
			key := edgeKey{from: n, to: e.To, branch: e.Branch}
			newOut, ok := outs[e.Branch]
			if !ok {
				if hasDefault {
					newOut = defaultOut
				} else {
					newOut = in
				}
			}
			prevOut := d.edgeOut[key]
			if prevOut != nil && prevOut.Equal(newOut) {
				continue
			}
			d.edgeOut[key] = newOut
			if !inWorklist[e.To] {
				worklist = append(worklist, e.To)
				inWorklist[e.To] = true
			}
		}
	}
}

// InOf returns the IN FlowScope for the CfgNode wrapping `ast`. Nil
// if the node was never reached or the AST isn't a statement in the
// analysed CFG.
// Go-only: language-adapter Kept for interface-parity; real callers use InOfAst.
func (d *DataFlowAnalysis) InOf(ast interface{ Token() Tokenish }) *FlowScope {
	// Real callers use InOfAst with a *CfgNode; this overload exists for interface parity.
	_ = ast
	return nil
}

// InOfAst is the concrete lookup used by TypeCheck — given a CFG and
// an AST statement node, return its IN FlowScope.
// Go-only: language-adapter Java callers use DiGraphNode.getAnnotation(); Go uses a typed map lookup via InOfAst.
func (d *DataFlowAnalysis) InOfAst(node *CfgNode) *FlowScope {
	if node == nil {
		return nil
	}
	return d.In[node]
}

// Tokenish is a tiny interface needed only to keep the (unused)
// InOf signature stable across refactors. Real callers don't touch
// it.
type Tokenish interface{ String() string }
