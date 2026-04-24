// Port of com/google/javascript/jscomp/DataFlowAnalysis.java
// @ closure-compiler v20260415 (forward + backward, LinearFlowState
// side-table, per-node iteration limit, FlowJoiner interface).
//
// The Java framework is abstract over the lattice type and direction
// of analysis. jsgo's use case today is a single concrete analysis —
// the TypeInferencePass (forward). Backward support lands here so the
// next caller (e.g. liveness / must-def) can reuse the same runner
// without re-implementing the worklist.

package jscomp

// FlowJoiner is the Go equivalent of Java's DataFlowAnalysis.FlowJoiner<L>.
// A new instance is created at each join point; JoinFlow is called once per
// incoming edge, then Finish returns the merged lattice element.
//
// For FlowScope the concrete implementation is FlowScopeJoiner.
type FlowJoiner interface {
	JoinFlow(s *FlowScope)
	Finish() *FlowScope
}

// FlowScopeJoiner implements FlowJoiner for the FlowScope lattice.
// It mirrors LinkedFlowScope.FlowScopeJoinOp in Java.
type FlowScopeJoiner struct {
	merged *FlowScope
}

// JoinFlow folds one incoming FlowScope into the running merge.
func (j *FlowScopeJoiner) JoinFlow(s *FlowScope) {
	if j.merged == nil {
		j.merged = s
	} else {
		j.merged = JoinFlowScopes(j.merged, s)
	}
}

// Finish returns the merged FlowScope, or an empty scope if no edge fed in.
func (j *FlowScopeJoiner) Finish() *FlowScope {
	if j.merged == nil {
		return EmptyFlowScope()
	}
	return j.merged
}

// TransferFunction computes a node's OUT FlowScope given its IN and
// (for branching nodes) the branch kind along which control reached
// it. Returns multiple outputs when a node refines differently along
// BranchOnTrue vs BranchOnFalse successors — the key is the Branch,
// the value is the FlowScope to attach to that outgoing edge.
type TransferFunction func(node *CfgNode, in *FlowScope) map[Branch]*FlowScope

// linearFlowState holds the IN, OUT and per-node visit count for one
// CfgNode. This is the Go equivalent of Java's LinearFlowState<L>
// annotation attached to each graph node.
type linearFlowState struct {
	in        *FlowScope
	out       *FlowScope
	stepCount int
}

// edgeKey identifies a directed CFG edge (source, destination, branch).
type edgeKey struct {
	from   *CfgNode
	to     *CfgNode
	branch Branch
}

// maxStepsPerNode mirrors DataFlowAnalysis.MAX_STEPS_PER_NODE.
// Analysis is considered divergent when any node exceeds this many visits.
const maxStepsPerNode = 20000

// DataFlowAnalysis is the worklist-based fixed-point runner. It supports
// both forward and backward analyses. The direction is set at construction
// time via NewDataFlowAnalysis (forward) or NewBackwardDataFlowAnalysis.
//
// After Run the computed IN FlowScope for each CfgNode is available via
// InOfAst. Backward analyses additionally populate the OUT table via
// OutOfAst.
type DataFlowAnalysis struct {
	Cfg     *ControlFlowGraph
	In      map[*CfgNode]*FlowScope
	edgeOut map[edgeKey]*FlowScope // forward branched edge annotations
	states  map[*CfgNode]*linearFlowState
	transfer TransferFunction
	initial  *FlowScope
	forward  bool
}

// NewDataFlowAnalysis wires a forward DataFlowAnalysis over `cfg` with
// the given transfer function. `initial` is the IN scope for the entry
// node (typically EmptyFlowScope()). This is the constructor used by all
// existing forward callers (TypeInferencePass).
func NewDataFlowAnalysis(cfg *ControlFlowGraph, initial *FlowScope, transfer TransferFunction) *DataFlowAnalysis {
	return &DataFlowAnalysis{
		Cfg:     cfg,
		In:      map[*CfgNode]*FlowScope{},
		edgeOut: map[edgeKey]*FlowScope{},
		states:  map[*CfgNode]*linearFlowState{},
		transfer: transfer,
		initial: initial,
		forward: true,
	}
}

// NewBackwardDataFlowAnalysis wires a backward DataFlowAnalysis.
// `exitLattice` is the OUT scope attached to the EXIT node (typically
// EmptyFlowScope() for must-analyses). The transfer function receives
// a node's OUT state and returns its IN state.
func NewBackwardDataFlowAnalysis(cfg *ControlFlowGraph, exitLattice *FlowScope, transfer TransferFunction) *DataFlowAnalysis {
	return &DataFlowAnalysis{
		Cfg:     cfg,
		In:      map[*CfgNode]*FlowScope{},
		edgeOut: map[edgeKey]*FlowScope{},
		states:  map[*CfgNode]*linearFlowState{},
		transfer: transfer,
		initial: exitLattice,
		forward: false,
	}
}

// Run executes the fixed-point worklist loop.
//
// Forward: seeds ENTRY with the initial lattice, propagates OUT→IN along
// successor edges, re-queues changed successors.
//
// Backward: seeds all non-EXIT nodes, propagates IN→OUT along predecessor
// edges (joining successor INs), runs transfer OUT→IN, re-queues changed
// predecessors.
//
// Panics when any node exceeds maxStepsPerNode visits (divergence guard,
// mirrors Java's IllegalStateException).
func (d *DataFlowAnalysis) Run() {
	empty := EmptyFlowScope()

	// Initialize side-table with bottom (empty) for every node.
	for _, n := range d.Cfg.Nodes {
		d.states[n] = &linearFlowState{in: empty, out: empty}
	}

	var worklist []*CfgNode
	inWorklist := map[*CfgNode]bool{}

	if d.forward {
		d.states[d.Cfg.Entry].in = d.initial
		worklist = []*CfgNode{d.Cfg.Entry}
		inWorklist[d.Cfg.Entry] = true
	} else {
		// Backward: EXIT holds the exit lattice; enqueue all other nodes.
		d.states[d.Cfg.Exit].out = d.initial
		for _, n := range d.Cfg.Nodes {
			if n == d.Cfg.Exit {
				continue
			}
			worklist = append(worklist, n)
			inWorklist[n] = true
		}
	}

	for len(worklist) > 0 {
		n := worklist[0]
		worklist = worklist[1:]
		delete(inWorklist, n)

		st := d.states[n]
		st.stepCount++
		if st.stepCount > maxStepsPerNode {
			panic("DataFlowAnalysis: divergent around node " + n.Kind)
		}

		if d.forward {
			worklist = d.stepForward(n, st, worklist, inWorklist)
		} else {
			worklist = d.stepBackward(n, st, worklist, inWorklist)
		}
	}

	// Expose converged IN scopes via the public map.
	for n, st := range d.states {
		d.In[n] = st.in
	}
}

// stepForward implements one forward-analysis worklist step (Java's flow +
// joinInputs for a forward node). Returns the (possibly extended) worklist.
func (d *DataFlowAnalysis) stepForward(n *CfgNode, st *linearFlowState, worklist []*CfgNode, inWorklist map[*CfgNode]bool) []*CfgNode {
	// Recompute IN by joining incoming edge annotations.
	if n == d.Cfg.Entry {
		st.in = d.initial
	} else {
		joiner := &FlowScopeJoiner{}
		for _, p := range n.Predecessors {
			for _, e := range p.Successors {
				if e.To != n {
					continue
				}
				if edgeIn := d.edgeOut[edgeKey{from: p, to: n, branch: e.Branch}]; edgeIn != nil {
					joiner.JoinFlow(edgeIn)
				}
			}
		}
		merged := joiner.Finish()
		if st.in != nil && st.in.Equal(merged) {
			return worklist // IN unchanged; no need to re-transfer
		}
		st.in = merged
	}

	// Run transfer and update per-edge OUT annotations.
	outs := d.transfer(n, st.in)
	defaultOut, hasDefault := outs[BranchUnconditional]
	for _, e := range n.Successors {
		newOut, ok := outs[e.Branch]
		if !ok {
			if hasDefault {
				newOut = defaultOut
			} else {
				newOut = st.in
			}
		}
		key := edgeKey{from: n, to: e.To, branch: e.Branch}
		if prev := d.edgeOut[key]; prev != nil && prev.Equal(newOut) {
			continue
		}
		d.edgeOut[key] = newOut
		if !inWorklist[e.To] && e.To != d.Cfg.Exit {
			worklist = append(worklist, e.To)
			inWorklist[e.To] = true
		}
	}
	return worklist
}

// stepBackward implements one backward-analysis worklist step. It joins
// successor INs to form a new OUT, runs the transfer OUT→IN, then
// re-enqueues predecessors whose OUT may have changed.
func (d *DataFlowAnalysis) stepBackward(n *CfgNode, st *linearFlowState, worklist []*CfgNode, inWorklist map[*CfgNode]bool) []*CfgNode {
	// Compute new OUT by joining INs of each successor.
	joiner := &FlowScopeJoiner{}
	for _, e := range n.Successors {
		succSt := d.states[e.To]
		if succSt == nil {
			continue
		}
		if e.To == d.Cfg.Exit {
			joiner.JoinFlow(d.initial) // exit lattice
		} else {
			joiner.JoinFlow(succSt.in)
		}
	}
	newOut := joiner.Finish()
	outChanged := !st.out.Equal(newOut)
	st.out = newOut

	// Run backward transfer: OUT → IN.
	transferResult := d.transfer(n, newOut)
	var newIn *FlowScope
	if v, ok := transferResult[BranchUnconditional]; ok {
		newIn = v
	} else {
		newIn = newOut
	}
	inChanged := !st.in.Equal(newIn)
	st.in = newIn

	// Propagate: enqueue predecessors if anything changed.
	if outChanged || inChanged {
		for _, p := range n.Predecessors {
			if p != d.Cfg.Exit && !inWorklist[p] {
				worklist = append(worklist, p)
				inWorklist[p] = true
			}
		}
	}
	return worklist
}

// InOf returns the IN FlowScope for the CfgNode wrapping `ast`. Nil
// if the node was never reached or the AST isn't a statement in the
// analysed CFG.
func (d *DataFlowAnalysis) InOf(ast interface{ Token() Tokenish }) *FlowScope {
	return nil // kept for interface-parity documentation; real callers use InOfAst
}

// InOfAst is the concrete lookup used by TypeCheck — given a CfgNode,
// return its IN FlowScope after Run has converged.
func (d *DataFlowAnalysis) InOfAst(node *CfgNode) *FlowScope {
	if node == nil {
		return nil
	}
	return d.In[node]
}

// OutOfAst returns the OUT FlowScope for a backward-analysis node.
// For forward analyses the OUT is per-edge (in edgeOut); this method
// exposes the merged OUT stored in the side-table.
func (d *DataFlowAnalysis) OutOfAst(node *CfgNode) *FlowScope {
	if node == nil {
		return nil
	}
	if st, ok := d.states[node]; ok {
		return st.out
	}
	return nil
}

// Tokenish is a tiny interface needed only to keep the (unused) InOf
// signature stable across refactors. Real callers don't touch it.
type Tokenish interface{ String() string }
