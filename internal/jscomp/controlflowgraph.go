// Port of com/google/javascript/jscomp/ControlFlowGraph.java
// @ closure-compiler v20260415 (narrow; no syntactic exceptions, no
// finally interactions, no for-in / for-of iterator edges).
//
// A ControlFlowGraph is a directed graph whose nodes are AST
// statements (plus a synthetic ENTRY and EXIT). Edges carry a Branch
// label so DataFlowAnalysis can distinguish the true/false arms of
// an IF condition or the fall-through edge after a return. One CFG
// is built per function body and per top-level script.

package jscomp

import "github.com/paullepoulpe/jsgo/internal/rhino"

// Branch mirrors ControlFlowGraph.Branch in Java.
type Branch int

const (
	// BranchUnconditional: straight-line successor.
	BranchUnconditional Branch = iota
	// BranchOnTrue: condition is truthy.
	BranchOnTrue
	// BranchOnFalse: condition is falsy.
	BranchOnFalse
	// BranchException: exceptional exit — used by THROW and by
	// body-to-catch edges in TRY statements.
	BranchException
	// BranchSyn: synthetic edge (e.g. switch-case match test, loop
	// back-edge). Doesn't carry a semantic refinement but keeps the
	// graph connected so reachability and dataflow work.
	BranchSyn
)

// String returns the DOT/debug label for a branch kind.
// Mirrors ControlFlowGraph.Branch.toString() in Java.
func (b Branch) String() string {
	switch b {
	case BranchUnconditional:
		return "UNCOND"
	case BranchOnTrue:
		return "ON_TRUE"
	case BranchOnFalse:
		return "ON_FALSE"
	case BranchException:
		return "ON_EX"
	case BranchSyn:
		return "SYN"
	default:
		return "UNKNOWN"
	}
}

// CfgEdge is a directed successor with a branch label.
type CfgEdge struct {
	To     *CfgNode
	Branch Branch
}

// CfgNode wraps an AST statement (or a synthetic entry/exit marker)
// and its CFG successors. Predecessors are not materialised — the
// DataFlowAnalysis walker rebuilds them from the successor lists.
type CfgNode struct {
	AST         *rhino.Node // nil for synthetic ENTRY / EXIT
	Kind        string      // "ENTRY", "EXIT", or the AST token name
	Successors  []CfgEdge
	Predecessors []*CfgNode // populated by ControlFlowGraph.computePreds
}

// ControlFlowGraph is the outer container: an entry + exit node plus
// the collected statement nodes. Maps AST → CFG so consumers can
// query the CfgNode for a given statement.
type ControlFlowGraph struct {
	Entry *CfgNode
	Exit  *CfgNode
	Nodes []*CfgNode
	byAst map[*rhino.Node]*CfgNode
}

// NewControlFlowGraph builds an empty CFG with synthetic entry/exit.
func NewControlFlowGraph() *ControlFlowGraph {
	cfg := &ControlFlowGraph{
		Entry: &CfgNode{Kind: "ENTRY"},
		Exit:  &CfgNode{Kind: "EXIT"},
		byAst: map[*rhino.Node]*CfgNode{},
	}
	cfg.Nodes = append(cfg.Nodes, cfg.Entry, cfg.Exit)
	return cfg
}

// node returns the CfgNode wrapping `ast`, creating one if needed.
// ENTRY / EXIT are reserved; pass nil to target the synthetic EXIT
// from a return/throw.
func (g *ControlFlowGraph) node(ast *rhino.Node) *CfgNode {
	if ast == nil {
		return g.Exit
	}
	if existing, ok := g.byAst[ast]; ok {
		return existing
	}
	n := &CfgNode{AST: ast, Kind: ast.Token().String()}
	g.byAst[ast] = n
	g.Nodes = append(g.Nodes, n)
	return n
}

// addEdge records `from → to` with the given branch label.
func (g *ControlFlowGraph) addEdge(from, to *CfgNode, b Branch) {
	if from == nil || to == nil {
		return
	}
	from.Successors = append(from.Successors, CfgEdge{To: to, Branch: b})
}

// Get returns the CfgNode wrapping `ast` or nil if absent.
func (g *ControlFlowGraph) Get(ast *rhino.Node) *CfgNode {
	return g.byAst[ast]
}

// computePreds fills Predecessors in one pass after construction so
// DataFlowAnalysis can walk the graph both ways.
func (g *ControlFlowGraph) computePreds() {
	for _, n := range g.Nodes {
		n.Predecessors = n.Predecessors[:0]
	}
	for _, n := range g.Nodes {
		for _, e := range n.Successors {
			e.To.Predecessors = append(e.To.Predecessors, n)
		}
	}
}
