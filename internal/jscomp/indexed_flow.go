package jscomp

import (
	"cmp"
	"math"
	"os"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: explicit experimental opt-in, checked once per analysis rather
// than in the fixed-point loop. Both backends share one build for differential
// tests and factorial combination measurements.
func dataFlowAnalysisUseIndexed() bool {
	return os.Getenv("JSGO_EXPERIMENT_INDEXED_FLOW") == "1"
}

// Custom analysis callbacks may inspect live graph annotations during a
// transfer. The shipped analyses below do not; preserve the generic engine for
// other extensions rather than changing their annotation visibility contract.
func dataFlowAnalysisSupportsIndexed(analysis, cfg any) bool {
	switch analysis.(type) {
	case *MustBeReachingVariableDef, *MaybeReachingVariableUse, *LiveVariablesAnalysis, *TypeInference:
		// Comparator factories on custom CFGs can have side effects. Restrict
		// automatic dispatch to the shipped immutable-comparator CFG types;
		// explicit sidecar users/tests can choose their own snapshot contract.
		switch cfg.(type) {
		case *ControlFlowAnalysis_AstControlFlowGraph, *ControlFlowGraphBase[*rhino.Node]:
			return true
		default:
			return false
		}
	default:
		return false
	}
}

// Go-only: experimental compiler-native CFG sidecar. IDs follow graph node
// encounter order. Both CSR edge lists preserve their original edge order,
// including duplicate edges and exception branches. A topology is immutable
// after construction, owns no analysis state, and can be shared by independent
// analyses while the AST/CFG is unchanged. The Java-compatible graph remains
// the construction/publication boundary and generic fallback.
type indexedFlowTopology[N comparable] struct {
	nodes              []DiGraph_DiGraphNode[N, ControlFlowGraph_Branch]
	values             []N
	records            []indexedFlowNode
	edges              []indexedFlowEdge
	originalEdges      []DiGraph_DiGraphEdge[N, ControlFlowGraph_Branch]
	incoming, outgoing []int32
	entry, implicit    int32
	forwardOrder       func(int32, int32) int
	backwardOrder      func(int32, int32) int
}

type indexedFlowNode struct {
	inStart, inEnd, outStart, outEnd int32
}

type indexedFlowEdge struct {
	from, to int32
	branch   ControlFlowGraph_Branch
}

func newIndexedFlowTopology[N comparable](cfg ControlFlowGraph[N]) *indexedFlowTopology[N] {
	nodes, edges := cfg.GetDiGraphNodes(), cfg.GetDiGraphEdges()
	if len(nodes) > math.MaxInt32 || len(edges) > math.MaxInt32 {
		return nil
	}
	t := &indexedFlowTopology[N]{nodes: nodes, originalEdges: edges,
		values: make([]N, len(nodes)), records: make([]indexedFlowNode, len(nodes)),
		edges:    make([]indexedFlowEdge, len(edges)),
		incoming: make([]int32, 0, len(edges)), outgoing: make([]int32, 0, len(edges))}
	ids := make(map[DiGraph_DiGraphNode[N, ControlFlowGraph_Branch]]int32, len(nodes))
	for id, node := range nodes {
		ids[node], t.values[id] = int32(id), node.GetValue()
	}
	var ok bool
	if t.entry, ok = ids[cfg.GetEntry()]; !ok {
		return nil
	}
	if t.implicit, ok = ids[cfg.GetImplicitReturn()]; !ok {
		return nil
	}
	edgeIDs := make(map[DiGraph_DiGraphEdge[N, ControlFlowGraph_Branch]]int32, len(edges))
	for id, edge := range edges {
		from, fromOK := ids[edge.GetSource()]
		to, toOK := ids[edge.GetDestination()]
		if !fromOK || !toOK {
			return nil
		}
		t.edges[id] = indexedFlowEdge{from: from, to: to, branch: edge.GetValue()}
		edgeIDs[edge] = int32(id)
	}
	for id, node := range nodes {
		record := &t.records[id]
		record.inStart = int32(len(t.incoming))
		for _, edge := range node.GetInEdges() {
			edgeID, ok := edgeIDs[edge]
			if !ok || len(t.incoming) == math.MaxInt32 {
				return nil
			}
			t.incoming = append(t.incoming, edgeID)
		}
		record.inEnd = int32(len(t.incoming))
		record.outStart = int32(len(t.outgoing))
		for _, edge := range node.GetOutEdges() {
			edgeID, ok := edgeIDs[edge]
			if !ok || len(t.outgoing) == math.MaxInt32 {
				return nil
			}
			t.outgoing = append(t.outgoing, edgeID)
		}
		record.outEnd = int32(len(t.outgoing))
	}
	// The native AST CFG's comparator reads immutable integer priorities.
	// Copy those once rather than dispatch through graph interfaces per heap
	// comparison. Other comparators retain their exact implementation.
	if _, native := any(cfg).(*ControlFlowAnalysis_AstControlFlowGraph); native {
		ranks := make([]int, len(nodes))
		for id, node := range nodes {
			ranks[id] = node.GetPriority()
		}
		t.forwardOrder = func(a, b int32) int { return cmp.Compare(ranks[a], ranks[b]) }
		t.backwardOrder = func(a, b int32) int { return cmp.Compare(ranks[b], ranks[a]) }
	} else {
		if compare := cfg.GetOptionalNodeComparator(true); compare != nil {
			t.forwardOrder = func(a, b int32) int { return compare(nodes[a], nodes[b]) }
		}
		if compare := cfg.GetOptionalNodeComparator(false); compare != nil {
			t.backwardOrder = func(a, b int32) int { return compare(nodes[a], nodes[b]) }
		}
	}
	return t
}

// indexedFlowResult keeps lattice state in contiguous, analysis-owned arrays.
// Compute touches no graph annotations; publication is a separate serial step.
// Concrete transfer functions retain their established lattice semantics:
// MustDef continues to use its canonical variable IDs and alias fallback.
type indexedFlowResult[N comparable, L LatticeElement] struct {
	topology *indexedFlowTopology[N]
	analysis DataFlowAnalysis[N, L]
	states   []DataFlowAnalysis_LinearFlowState[L]
	branches []L
	forward  bool
	branched bool
}

func newIndexedFlowResult[N comparable, L LatticeElement](topology *indexedFlowTopology[N], analysis DataFlowAnalysis[N, L]) *indexedFlowResult[N, L] {
	r := &indexedFlowResult[N, L]{topology: topology, analysis: analysis,
		states:  make([]DataFlowAnalysis_LinearFlowState[L], len(topology.values)),
		forward: analysis.IsForward(), branched: analysis.IsBranched()}
	for id := range r.states {
		in, out := analysis.CreateInitialEstimateLattice(), analysis.CreateInitialEstimateLattice()
		jutil.CheckNotNull(in)
		jutil.CheckNotNull(out)
		r.states[id].in, r.states[id].out = in, out
	}
	if r.branched {
		jutil.CheckState(r.forward)
		r.branches = make([]L, len(topology.edges))
		for id := range r.branches {
			r.branches[id] = analysis.CreateInitialEstimateLattice()
		}
	}
	return r
}

// AnalyzeIndexed computes and publishes an indexed analysis. A false result
// leaves annotations untouched and asks the caller to use the generic engine.
// An experimental environment toggle selects it without changing CLI flags.
func (b *DataFlowAnalysisBase[N, L]) AnalyzeIndexed() bool {
	topology := newIndexedFlowTopology(b.cfg)
	if topology == nil {
		return false
	}
	result := newIndexedFlowResult(topology, b.self)
	// Publish partial state on a divergent-analysis exception too, as the
	// generic engine leaves its initialized/partially updated annotations.
	defer result.publish()
	result.compute()
	return true
}

func (r *indexedFlowResult[N, L]) publish() {
	for id, node := range r.topology.nodes {
		node.SetAnnotation(&r.states[id])
	}
	if r.branched {
		for id, edge := range r.topology.originalEdges {
			edge.SetAnnotation(r.branches[id])
		}
	}
}

func (r *indexedFlowResult[N, L]) inputFromEdge(id int32) L {
	if r.branched {
		return r.branches[id]
	}
	edge := r.topology.edges[id]
	if r.forward {
		return r.states[edge.from].out
	}
	if edge.to == r.topology.implicit {
		return r.analysis.CreateEntryLattice()
	}
	return r.states[edge.to].in
}

func (r *indexedFlowResult[N, L]) joinInputs(id int32) {
	state := &r.states[id]
	if r.forward && id == r.topology.entry {
		state.SetIn(r.analysis.CreateEntryLattice())
		return
	}
	record := r.topology.records[id]
	edges := r.topology.incoming[record.inStart:record.inEnd]
	if !r.forward {
		edges = r.topology.outgoing[record.outStart:record.outEnd]
	}
	var result L
	switch len(edges) {
	case 0:
		return
	case 1:
		result = r.inputFromEdge(edges[0])
	default:
		joiner := r.analysis.CreateFlowJoiner()
		for _, edge := range edges {
			joiner.JoinFlow(r.inputFromEdge(edge))
		}
		result = joiner.Finish()
	}
	if r.forward {
		state.SetIn(result)
	} else {
		state.SetOut(result)
	}
}

func (r *indexedFlowResult[N, L]) flow(id int32) bool {
	state := &r.states[id]
	if !r.forward {
		before := state.in
		state.SetIn(r.analysis.FlowThrough(r.topology.values[id], state.out))
		return !jutil.ObjectsEquals(before, state.in)
	}
	before := state.out
	state.SetOut(r.analysis.FlowThrough(r.topology.values[id], state.in))
	changed := !jutil.ObjectsEquals(before, state.out)
	if r.branched {
		brancher := r.analysis.CreateFlowBrancher(r.topology.values[id], state.out)
		record := r.topology.records[id]
		for _, edgeID := range r.topology.outgoing[record.outStart:record.outEnd] {
			before := r.branches[edgeID]
			r.branches[edgeID] = brancher(r.topology.edges[edgeID].branch)
			if !changed {
				changed = !jutil.ObjectsEquals(before, r.branches[edgeID])
			}
		}
	}
	return changed
}

func (r *indexedFlowResult[N, L]) compute() {
	compare := r.topology.forwardOrder
	if !r.forward {
		compare = r.topology.backwardOrder
	}
	queue := indexedFlowQueue{compare: compare, queued: make([]bool, len(r.states))}
	for id := range r.states {
		if int32(id) != r.topology.implicit {
			queue.add(int32(id))
		}
	}
	for !queue.empty() {
		id := queue.remove()
		state := &r.states[id]
		previous := state.stepCount
		state.stepCount++
		if previous > DataFlowAnalysis_MAX_STEPS_PER_NODE {
			var display any = r.topology.values[id]
			if len(r.topology.nodes) != 0 {
				display = r.topology.nodes[id]
			}
			panic(jutil.NewIllegalStateExceptionMsg("Dataflow analysis appears to diverge around: " + jutil.ToJavaString(display)))
		}
		r.joinInputs(id)
		if r.flow(id) {
			record := r.topology.records[id]
			if r.forward {
				for _, edge := range r.topology.outgoing[record.outStart:record.outEnd] {
					if next := r.topology.edges[edge].to; next != r.topology.implicit {
						queue.add(next)
					}
				}
			} else {
				for _, edge := range r.topology.incoming[record.inStart:record.inEnd] {
					if next := r.topology.edges[edge].from; next != r.topology.implicit {
						queue.add(next)
					}
				}
			}
		}
	}
	if r.forward {
		r.joinInputs(r.topology.implicit)
	}
}

// The integer heap uses the same strictly-less sift-up and left-child-on-tie
// sift-down rules as the Java-compatible priority queue. FIFO mode uses the
// same encounter order. Only IDs and a dense membership array are stored.
type indexedFlowQueue struct {
	values  []int32
	head    int
	queued  []bool
	compare func(int32, int32) int
}

func (q *indexedFlowQueue) empty() bool { return q.head == len(q.values) }

func (q *indexedFlowQueue) add(id int32) {
	if q.queued[id] {
		return
	}
	q.queued[id] = true
	q.values = append(q.values, id)
	if q.compare != nil {
		index := len(q.values) - 1
		for index > 0 {
			parent := (index - 1) / 2
			if q.compare(id, q.values[parent]) >= 0 {
				break
			}
			q.values[index] = q.values[parent]
			index = parent
		}
		q.values[index] = id
	}
}

func (q *indexedFlowQueue) remove() int32 {
	id := q.values[q.head]
	q.queued[id] = false
	if q.compare == nil {
		q.head++
		if q.head == len(q.values) {
			q.head, q.values = 0, q.values[:0]
		}
		return id
	}
	last := q.values[len(q.values)-1]
	q.values = q.values[:len(q.values)-1]
	if len(q.values) != 0 {
		index := 0
		for {
			child := 2*index + 1
			if child >= len(q.values) {
				break
			}
			if sibling := child + 1; sibling < len(q.values) && q.compare(q.values[child], q.values[sibling]) > 0 {
				child = sibling
			}
			if q.compare(last, q.values[child]) <= 0 {
				break
			}
			q.values[index] = q.values[child]
			index = child
		}
		q.values[index] = last
	}
	return id
}
