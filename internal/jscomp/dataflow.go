// Port of com/google/javascript/jscomp/DataFlowAnalysis.java
// @ closure-compiler v20260415 (narrow; forward only, no lattice-
// annotation side table, no fixed-point limit heuristic).
//
// The Java framework is abstract over the lattice type and direction
// of analysis. jsgo's use case is a single concrete analysis — the
// TypeInferencePass — so the port inlines the forward-worklist loop
// against FlowScope directly. Backward / max-fix-point variants land
// with the next caller that needs them.

package jscomp

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
	Cfg       *ControlFlowGraph
	In        map[*CfgNode]*FlowScope
	edgeOut   map[edgeKey]*FlowScope
	transfer  TransferFunction
	initial   *FlowScope
}

type edgeKey struct {
	from   *CfgNode
	to     *CfgNode
	branch Branch
}

// NewDataFlowAnalysis wires a DataFlowAnalysis over `cfg` with the
// given transfer function. `initial` is the IN scope attached to the
// entry node (typically an empty FlowScope).
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
func (d *DataFlowAnalysis) Run() {
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
