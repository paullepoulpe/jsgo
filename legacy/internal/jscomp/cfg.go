// Port of com/google/javascript/jscomp/ControlFlowGraph.java
// @ closure-compiler v20260415 (narrow; structural shell + AST-walk helpers).
//
// This file ports the public interface of ControlFlowGraph: Branch enum,
// getImplicitReturn, getEntry, isImplicitReturn, getOptionalNodeComparator,
// and the static isEnteringNewCfgNode helper. The full CFG-build and
// dataflow machinery live in controlflowgraph.go and controlflowanalysis.go.
//
// The definitelyReturns helper (a simple AST-walk approximation used by
// post-guard narrowing in TypeCheck) is also kept here because it exposes the
// same ControlFlowGraph.java Java-source attribution.

package jscomp

import (
	"github.com/paullepoulpe/jsgo/legacy/internal/rhino"
)

// isConditional mirrors ControlFlowGraph.Branch.isConditional() in Java.
// Returns true for ON_TRUE and ON_FALSE branch types.
func (b Branch) isConditional() bool {
	return b == BranchOnTrue || b == BranchOnFalse
}

// getImplicitReturn returns the synthetic EXIT node of the CFG.
// Mirrors ControlFlowGraph.getImplicitReturn().
func (g *ControlFlowGraph) getImplicitReturn() *CfgNode {
	return g.Exit
}

// getEntry returns the synthetic ENTRY node of the CFG.
// Mirrors ControlFlowGraph.getEntry().
func (g *ControlFlowGraph) getEntry() *CfgNode {
	return g.Entry
}

// isImplicitReturn reports whether `node` is the synthetic exit (implicit
// return). Mirrors ControlFlowGraph.isImplicitReturn(DiGraphNode).
func (g *ControlFlowGraph) isImplicitReturn(node *CfgNode) bool {
	return node == g.Exit
}

// getOptionalNodeComparator returns nil — the default implementation in
// Java returns null, meaning no ordering is enforced. Overriding subclasses
// in Java return a comparator; jsgo's CFG does not require ordering.
// Mirrors ControlFlowGraph.getOptionalNodeComparator(boolean).
func (g *ControlFlowGraph) getOptionalNodeComparator(_ bool) func(a, b *CfgNode) int {
	return nil
}

// isEnteringNewCfgNode reports whether `n` should be represented by a
// new CFG node separate from its parent. Mirrors
// ControlFlowGraph.isEnteringNewCfgNode(Node) in Java.
//
// Rules (matching Java switch):
//
//	BLOCK, ROOT, SCRIPT, TRY → always a new CFG node
//	FUNCTION → new node unless `n` is the function's second child (body)
//	WHILE, DO, IF → new node when `n` is NOT the condition expression
//	FOR(;;) → new node when `n` is NOT the condition expression
//	FOR_IN → new node when `n` is NOT the first child (LHS)
//	CASE, CATCH, WITH → new node when `n` is NOT the first child
//	default → false
func isEnteringNewCfgNode(n *rhino.Node) bool {
	parent := n.Parent()
	if parent == nil {
		return false
	}
	switch parent.Token() {
	case rhino.BLOCK, rhino.ROOT, rhino.SCRIPT, rhino.TRY:
		return true
	case rhino.FUNCTION:
		// In Java: n != parent.getSecondChild()
		kids := parent.Children()
		if len(kids) >= 2 && kids[1] == n {
			return false
		}
		return true
	case rhino.WHILE, rhino.DO, rhino.IF:
		// New CFG node when n is NOT the condition expression.
		// getConditionExpression(IF) = first child.
		return parent.FirstChild() != n
	case rhino.FOR:
		// The FOR(;;) node represents only the condition check.
		// init and increment have their own CFG nodes.
		// Condition is the second child of FOR.
		kids := parent.Children()
		if len(kids) >= 2 && kids[1] == n {
			return false // this IS the condition expression
		}
		return true
	case rhino.FOR_IN, rhino.FOR_OF:
		// New node unless n is the first child (the LHS variable).
		return parent.FirstChild() != n
	case rhino.CASE, rhino.CATCH, rhino.WITH:
		// New node unless n is the first child (condition / binding).
		return parent.FirstChild() != n
	}
	return false
}

// definitelyReturns reports whether every execution path through `n`
// ends in a `return EXPR` or `throw`. Walks AST structure rather than
// a real CFG — sufficient for the post-guard narrowing used by TypeCheck.
//
// Statement-level contract:
//
//	RETURN      → true iff carries an expression
//	THROW       → true
//	IF(c, t, e) → true iff both t and e definitely return
//	BLOCK       → true iff some statement definitely returns
//	FUNCTION/CLASS/EXPR_RESULT/VAR/etc. → false (no termination)
func definitelyReturns(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	switch n.Token() {
	case rhino.RETURN:
		return n.FirstChild() != nil
	case rhino.THROW:
		return true
	case rhino.IF:
		kids := n.Children()
		if len(kids) < 3 {
			return false
		}
		return definitelyReturns(kids[1]) && definitelyReturns(kids[2])
	case rhino.BLOCK:
		for _, c := range n.Children() {
			if definitelyReturns(c) {
				return true
			}
		}
		return false
	}
	return false
}
