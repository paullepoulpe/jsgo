package jscomp

import (
	"cmp"
	"os"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Experimental direct CFG construction for consumers that need no public graph
// objects. IDs are assigned in the original graph's node creation order; linked
// integer edge chains retain independent incoming/outgoing encounter order.
// Freezing produces the indexed solver's topology without allocating even a
// temporary LinkedDirectedGraph, graph node, graph edge, or annotation adapter.
type nativeCFGBuilder struct {
	ids      map[*rhino.Node]int32
	nodes    []nativeCFGNode
	edges    []nativeCFGEdge
	entry    int32
	implicit int32
}

type nativeCFGNode struct {
	value                              *rhino.Node
	inFirst, inLast, outFirst, outLast int32
	inCount, outCount                  int32
	priority                           int32
}

type nativeCFGEdge struct {
	indexedFlowEdge
	nextIn, nextOut int32
}

func nativeCFGEnabled() bool { return os.Getenv("JSGO_EXPERIMENT_NATIVE_CFG") == "1" }

func (g *nativeCFGBuilder) initialize(entry *rhino.Node) {
	g.ids = make(map[*rhino.Node]int32)
	g.implicit = g.node(nil)
	g.entry = g.node(entry)
}

func (g *nativeCFGBuilder) node(value *rhino.Node) int32 {
	if id, ok := g.ids[value]; ok {
		return id
	}
	id := int32(len(g.nodes))
	g.ids[value] = id
	g.nodes = append(g.nodes, nativeCFGNode{value: value, inFirst: -1, inLast: -1, outFirst: -1, outLast: -1, priority: -1})
	return id
}

func (g *nativeCFGBuilder) connect(fromValue *rhino.Node, branch ControlFlowGraph_Branch, toValue *rhino.Node) {
	from, to := g.node(fromValue), g.node(toValue)
	// Match the generic graph's smaller-adjacency search. Always scanning
	// outgoing edges would make large switches quadratic during construction.
	if g.nodes[from].outCount < g.nodes[to].inCount {
		for edge := g.nodes[from].outFirst; edge >= 0; edge = g.edges[edge].nextOut {
			if previous := g.edges[edge]; previous.to == to && previous.branch == branch {
				return
			}
		}
	} else {
		for edge := g.nodes[to].inFirst; edge >= 0; edge = g.edges[edge].nextIn {
			if previous := g.edges[edge]; previous.from == from && previous.branch == branch {
				return
			}
		}
	}
	id := int32(len(g.edges))
	g.edges = append(g.edges, nativeCFGEdge{indexedFlowEdge: indexedFlowEdge{from: from, to: to, branch: branch}, nextIn: -1, nextOut: -1})
	if last := g.nodes[from].outLast; last >= 0 {
		g.edges[last].nextOut = id
	} else {
		g.nodes[from].outFirst = id
	}
	g.nodes[from].outLast = id
	if last := g.nodes[to].inLast; last >= 0 {
		g.edges[last].nextIn = id
	} else {
		g.nodes[to].inFirst = id
	}
	g.nodes[to].inLast = id
	g.nodes[from].outCount++
	g.nodes[to].inCount++
}

func (c *ControlFlowAnalysis) cfgEntryValue() *rhino.Node {
	if c.nativeCFG != nil {
		return c.nativeCFG.nodes[c.nativeCFG.entry].value
	}
	return c.cfg.GetEntry().GetValue()
}

// ComputeIndexedCfg is an explicit native construction entry point. Ordinary
// graph consumers retain ComputeCfg; no unsupported interface methods are
// stubbed and no generic graph is kept behind an adapter.
func (b *ControlFlowAnalysis_Builder) ComputeIndexedCfg() *indexedFlowTopology[*rhino.Node] {
	jutil.CheckNotNullMsg(b.compiler, "Need to call setCompiler()")
	jutil.CheckNotNullMsg(b.cfgRoot, "Need to call setCfgRoot()")
	c := newControlFlowAnalysis(b.compiler, b.shouldTraverseFunctions, b.edgeAnnotations)
	c.nativeCFG = &nativeCFGBuilder{}
	c.initializeCfgTraversal(b.cfgRoot)
	g := c.nativeCFG
	c.priorityCounter = 0
	c.prioritizeNativeFromEntry(g.entry)
	if c.shouldTraverseFunctions {
		for id := range g.nodes {
			if value := g.nodes[id].value; value != nil && value.IsFunction() {
				c.prioritizeNativeFromEntry(int32(id))
			}
		}
	}
	c.astPosition = nil
	for id := range g.nodes {
		if g.nodes[id].priority < 0 {
			c.priorityCounter++
			g.nodes[id].priority = int32(c.priorityCounter)
		}
	}
	c.priorityCounter++
	g.nodes[g.implicit].priority = int32(c.priorityCounter)
	return g.freeze()
}

func (c *ControlFlowAnalysis) prioritizeNativeFromEntry(entry int32) {
	g := c.nativeCFG
	queue := jutil.NewPriorityQueueWithCapacity(10, func(a, b int32) int {
		return cmp.Compare(c.astPositionOfValue(g.nodes[a].value), c.astPositionOfValue(g.nodes[b].value))
	})
	queue.Add(entry)
	for !queue.IsEmpty() {
		id := queue.Remove()
		if g.nodes[id].priority >= 0 {
			continue
		}
		c.priorityCounter++
		g.nodes[id].priority = int32(c.priorityCounter)
		for edge := g.nodes[id].outFirst; edge >= 0; edge = g.edges[edge].nextOut {
			queue.Add(g.edges[edge].to)
		}
	}
}

func (g *nativeCFGBuilder) freeze() *indexedFlowTopology[*rhino.Node] {
	t := &indexedFlowTopology[*rhino.Node]{
		values: make([]*rhino.Node, len(g.nodes)), records: make([]indexedFlowNode, len(g.nodes)),
		edges: make([]indexedFlowEdge, len(g.edges)), incoming: make([]int32, 0, len(g.edges)),
		outgoing: make([]int32, 0, len(g.edges)), entry: g.entry, implicit: g.implicit,
	}
	ranks := make([]int32, len(g.nodes))
	// Generic getDiGraphEdges groups edges by source-node encounter order,
	// which can differ from construction order. Preserve that ordering for
	// branch initialization and remap both CSR lists to those exact IDs.
	edgeIDs := make([]int32, len(g.edges))
	next := int32(0)
	for _, node := range g.nodes {
		for edge := node.outFirst; edge >= 0; edge = g.edges[edge].nextOut {
			edgeIDs[edge] = next
			t.edges[next] = g.edges[edge].indexedFlowEdge
			next++
		}
	}
	for id, node := range g.nodes {
		t.values[id], ranks[id] = node.value, node.priority
		record := &t.records[id]
		record.inStart = int32(len(t.incoming))
		for edge := node.inFirst; edge >= 0; edge = g.edges[edge].nextIn {
			t.incoming = append(t.incoming, edgeIDs[edge])
		}
		record.inEnd = int32(len(t.incoming))
		record.outStart = int32(len(t.outgoing))
		for edge := node.outFirst; edge >= 0; edge = g.edges[edge].nextOut {
			t.outgoing = append(t.outgoing, edgeIDs[edge])
		}
		record.outEnd = int32(len(t.outgoing))
	}
	t.forwardOrder = func(a, b int32) int { return cmp.Compare(ranks[a], ranks[b]) }
	t.backwardOrder = func(a, b int32) int { return cmp.Compare(ranks[b], ranks[a]) }
	return t
}
