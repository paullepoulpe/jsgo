// Port of com/google/javascript/jscomp/ControlFlowAnalysis.java
// @ closure-compiler v20260415 (narrow; no try/catch/finally edges,
// no for-in / for-of iterator modeling, no labeled break/continue).
//
// Builds a ControlFlowGraph for a FUNCTION body or top-level SCRIPT.
// Walks the AST producing one CfgNode per statement, then connects
// sequential siblings with unconditional edges and IF branches with
// the expected BranchOnTrue / BranchOnFalse labels. RETURN statements
// and the function's fall-through tail both flow to the synthetic
// EXIT node.

package jscomp

import "github.com/paullepoulpe/jsgo/internal/rhino"

// BuildControlFlowGraph runs on a FUNCTION or SCRIPT node and returns
// its CFG. For a FUNCTION, the body is kids[2]; for SCRIPT, the body
// is the node itself. Nested FUNCTION / CLASS nodes are NOT descended
// into — each has its own CFG, built lazily by the caller.
func BuildControlFlowGraph(root *rhino.Node) *ControlFlowGraph {
	var body *rhino.Node
	switch root.Token() {
	case rhino.FUNCTION:
		kids := root.Children()
		if len(kids) < 3 {
			cfg := NewControlFlowGraph()
			cfg.addEdge(cfg.Entry, cfg.Exit, BranchUnconditional)
			cfg.computePreds()
			return cfg
		}
		body = kids[2]
	case rhino.SCRIPT:
		body = root
	default:
		cfg := NewControlFlowGraph()
		cfg.addEdge(cfg.Entry, cfg.Exit, BranchUnconditional)
		cfg.computePreds()
		return cfg
	}
	b := &cfgBuilder{cfg: NewControlFlowGraph()}
	// exitSink is the node every "fall through past the last statement"
	// edge lands on. For a top-level run, that's the CFG exit.
	b.visitStatements(body.Children(), b.cfg.Entry, b.cfg.Exit)
	b.cfg.computePreds()
	return b.cfg
}

// cfgBuilder carries the mutable state of a single CFG build.
type cfgBuilder struct {
	cfg *ControlFlowGraph
}

// visitStatements walks a sequence of statements, chaining them with
// unconditional edges. Returns nothing — the edges are written onto
// `b.cfg`. `prev` is the node whose fall-through edge should be the
// first statement; `fallthrough` is where the last statement's
// implicit fall-through goes.
func (b *cfgBuilder) visitStatements(stmts []*rhino.Node, prev *CfgNode, fallthrough_ *CfgNode) {
	current := prev
	for i, stmt := range stmts {
		// Dead statements after a non-falling-through control-flow
		// construct (return/throw) are still linked so the walker
		// doesn't leave orphans, but they never receive an IN scope
		// from the dataflow analysis.
		next := b.visitStatement(stmt, current)
		if next == nil {
			// Statement terminates flow (return / throw). Any
			// following statements start fresh from a null "prev";
			// they're reachable only as jump targets, but at this
			// grammar scope there aren't labelled jump targets, so
			// we chain them straight from the exit marker and let
			// DataFlow drop them on the floor.
			current = nil
			_ = i
			continue
		}
		current = next
	}
	if current != nil {
		b.cfg.addEdge(current, fallthrough_, BranchUnconditional)
	}
}

// visitStatement processes a single statement and returns the CfgNode
// that the following statement's unconditional predecessor should be.
// Returns nil when the statement terminates flow (return / throw).
// `prev` is the CfgNode of the previous statement or the entry node;
// visitStatement connects `prev → thisNode` with an unconditional
// edge before processing children. Pass `prev == nil` to skip that
// edge (used when the caller has already attached a branch-labeled
// predecessor, e.g. from an IF's true / false arm).
func (b *cfgBuilder) visitStatement(stmt *rhino.Node, prev *CfgNode) *CfgNode {
	n := b.cfg.node(stmt)
	if prev != nil {
		b.cfg.addEdge(prev, n, BranchUnconditional)
	}
	switch stmt.Token() {
	case rhino.IF:
		return b.visitIf(stmt, n)
	case rhino.RETURN:
		// Fall-through from RETURN goes to EXIT.
		b.cfg.addEdge(n, b.cfg.Exit, BranchUnconditional)
		return nil
	case rhino.BLOCK:
		// Nested block — flatten. The block node itself is still a
		// CFG node (so transfer functions can associate a FlowScope
		// with it), but its children are chained individually.
		if kids := stmt.Children(); len(kids) > 0 {
			after := &CfgNode{Kind: "BLOCK_JOIN"}
			b.cfg.Nodes = append(b.cfg.Nodes, after)
			b.visitStatements(kids, n, after)
			return after
		}
		return n
	}
	// Everything else is a linear statement — the CfgNode alone is
	// enough, and fall-through goes to the next sibling.
	return n
}

// visitIf handles an `if (cond) then [else]` statement. Emits
// BranchOnTrue / BranchOnFalse edges from the IF node to the heads
// of the then / else blocks, then merges both tails at a shared
// join node that becomes the predecessor of the following sibling.
func (b *cfgBuilder) visitIf(stmt *rhino.Node, n *CfgNode) *CfgNode {
	kids := stmt.Children()
	if len(kids) < 2 {
		return n
	}
	thenBlock := kids[1]
	var elseBlock *rhino.Node
	if len(kids) >= 3 {
		elseBlock = kids[2]
	}
	join := &CfgNode{Kind: "IF_JOIN"}
	b.cfg.Nodes = append(b.cfg.Nodes, join)
	b.visitBranch(n, thenBlock, BranchOnTrue, join)
	if elseBlock != nil {
		b.visitBranch(n, elseBlock, BranchOnFalse, join)
	} else {
		b.cfg.addEdge(n, join, BranchOnFalse)
	}
	return join
}

// visitBranch wires one arm of an if/else. The first statement
// receives a branch-labeled predecessor edge (OnTrue / OnFalse);
// subsequent statements chain unconditionally. A non-terminating
// tail joins `join`.
func (b *cfgBuilder) visitBranch(ifNode *CfgNode, block *rhino.Node, branch Branch, join *CfgNode) {
	stmts := flattenBlock(block)
	if len(stmts) == 0 {
		b.cfg.addEdge(ifNode, join, branch)
		return
	}
	// Process each branch statement — the first inherits the branch
	// label via a synthetic "dummy" that's really `ifNode` with the
	// correct edge kind; subsequent statements are unconditional.
	first := stmts[0]
	firstNode := b.cfg.node(first)
	b.cfg.addEdge(ifNode, firstNode, branch)
	// Drive visitStatement directly so control-terminating shapes
	// (RETURN / nested IF) record their real successors. We don't
	// use visitStatements here because that would add a *second*
	// unconditional edge from ifNode → firstNode.
	next := b.visitStatement(first, nil /* skip the prev-edge add, we already did it */)
	if next != nil {
		b.visitStatements(stmts[1:], next, join)
	} else if len(stmts) > 1 {
		// Dead statements after a terminator — link them linearly
		// from a null prev so they get FlowScope-less CFG nodes but
		// no incoming flow.
		b.visitStatements(stmts[1:], nil, join)
	}
}

// flattenBlock returns the statement list inside `n`. If `n` is a
// BLOCK, returns its children; otherwise returns a single-element
// slice with `n` itself so `if (c) x;` parses like `if (c) { x; }`.
func flattenBlock(n *rhino.Node) []*rhino.Node {
	if n == nil {
		return nil
	}
	if n.Token() == rhino.BLOCK {
		return n.Children()
	}
	return []*rhino.Node{n}
}
