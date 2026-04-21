// Port of com/google/javascript/jscomp/DataFlowAnalysis.java
// @ closure-compiler v20260415 (forward + backward worklists; no
// lattice-annotation side table, no fixed-point limit heuristic).
//
// The Java framework is abstract over the lattice type and direction
// of analysis. jsgo's use case is a single concrete analysis — the
// TypeInferencePass (forward) — but the backward worklist is also
// provided so live-variable / liveness analyses have a direct entry
// point when they land.

package jscomp

// TransferFunction computes a node's OUT FlowScope given its IN and
// (for branching nodes) the branch kind along which control reached
// it. Returns multiple outputs when a node refines differently along
// BranchOnTrue vs BranchOnFalse successors — the key is the Branch,
// the value is the FlowScope to attach to that outgoing edge.
type TransferFunction func(node *CfgNode, in *FlowScope) map[Branch]*FlowScope

// Direction mirrors DataFlowAnalysis.Direction in Java — forward
// analyses propagate IN → OUT along successor edges; backward analyses
// propagate OUT → IN along predecessor edges.
type Direction int

const (
	DirectionForward Direction = iota
	DirectionBackward
)

// DataFlowAnalysis is the worklist runner. It stores the computed IN
// scope for each CfgNode keyed by pointer — callers consume this after
// Run returns.
type DataFlowAnalysis struct {
	Cfg       *ControlFlowGraph
	In        map[*CfgNode]*FlowScope
	edgeOut   map[edgeKey]*FlowScope
	transfer  TransferFunction
	initial   *FlowScope
	direction Direction
}

type edgeKey struct {
	from   *CfgNode
	to     *CfgNode
	branch Branch
}

// NewDataFlowAnalysis wires a forward DataFlowAnalysis over `cfg`
// with the given transfer function. `initial` is the IN scope attached
// to the entry node (typically an empty FlowScope).
func NewDataFlowAnalysis(cfg *ControlFlowGraph, initial *FlowScope, transfer TransferFunction) *DataFlowAnalysis {
	return &DataFlowAnalysis{
		Cfg:       cfg,
		In:        map[*CfgNode]*FlowScope{},
		edgeOut:   map[edgeKey]*FlowScope{},
		transfer:  transfer,
		initial:   initial,
		direction: DirectionForward,
	}
}

// NewBackwardDataFlowAnalysis wires a backward DataFlowAnalysis —
// `initial` is the OUT scope attached to the exit node. Transfer
// functions see the node's OUT scope via `in` (the name is kept for
// symmetry with the forward API) and produce the node's IN scope
// which is then propagated along predecessor edges.
func NewBackwardDataFlowAnalysis(cfg *ControlFlowGraph, initial *FlowScope, transfer TransferFunction) *DataFlowAnalysis {
	return &DataFlowAnalysis{
		Cfg:       cfg,
		In:        map[*CfgNode]*FlowScope{},
		edgeOut:   map[edgeKey]*FlowScope{},
		transfer:  transfer,
		initial:   initial,
		direction: DirectionBackward,
	}
}

// Run executes the worklist fixed-point loop in whichever direction
// was chosen at construction. Termination is guaranteed because
// FlowScope entries only widen (each name either keeps its refined
// type or joins to a union that subsumes it), and NewUnion
// deduplicates by name — so the number of distinct union shapes per
// name is bounded by the type lattice.
func (d *DataFlowAnalysis) Run() {
	if d.direction == DirectionBackward {
		d.runBackward()
		return
	}
	d.runForward()
}

func (d *DataFlowAnalysis) runForward() {
	d.In[d.Cfg.Entry] = d.initial
	worklist := []*CfgNode{d.Cfg.Entry}
	// Each visit bounds is mostly a safety belt — the lattice itself
	// is finite so we shouldn't need it, but wide unions produced by
	// a pathological transfer could otherwise loop.
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

// runBackward propagates the OUT lattice from EXIT to ENTRY along
// predecessor edges. For a node N, OUT = join(IN of successors) via
// edgeOut (interpreted as "OUT attached to the successor edge that
// flows back to N"). The transfer function's returned map[Branch] is
// keyed by the incoming-edge branch as seen from the predecessor's
// perspective — matching Java's BackwardAnalysis.flowThrough contract.
func (d *DataFlowAnalysis) runBackward() {
	d.In[d.Cfg.Exit] = d.initial
	worklist := []*CfgNode{d.Cfg.Exit}
	const maxIterations = 20000
	iterations := 0
	inWorklist := map[*CfgNode]bool{d.Cfg.Exit: true}
	for len(worklist) > 0 && iterations < maxIterations {
		iterations++
		n := worklist[0]
		worklist = worklist[1:]
		delete(inWorklist, n)
		in := d.In[n]
		// For backward flow, the node's OUT is computed by joining
		// the edge-contributions from each of its *successors* (since
		// edgeOut maps from→to→IN-of-to). We skip this step on the
		// exit node.
		if n != d.Cfg.Exit {
			var merged *FlowScope
			for _, e := range n.Successors {
				edgeIn := d.edgeOut[edgeKey{from: n, to: e.To, branch: e.Branch}]
				if edgeIn == nil {
					continue
				}
				if merged == nil {
					merged = edgeIn
				} else {
					merged = JoinFlowScopes(merged, edgeIn)
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
		defaultOut, hasDefault := outs[BranchUnconditional]
		// Push the transfer result back along predecessor edges.
		for _, p := range n.Predecessors {
			for _, e := range p.Successors {
				if e.To != n {
					continue
				}
				key := edgeKey{from: p, to: n, branch: e.Branch}
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
				if !inWorklist[p] {
					worklist = append(worklist, p)
					inWorklist[p] = true
				}
			}
		}
	}
}

// InOf returns the IN FlowScope for the CfgNode wrapping `ast`. Nil
// if the node was never reached or the AST isn't a statement in the
// analysed CFG.
func (d *DataFlowAnalysis) InOf(ast interface{ Token() Tokenish }) *FlowScope {
	return nil // kept for interface-parity documentation; real callers use InOfAst below
}

// InOfAst is the concrete lookup used by TypeCheck — given a CFG and
// an AST statement node, return its IN FlowScope.
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
