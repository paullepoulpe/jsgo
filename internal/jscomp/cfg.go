// Port of com/google/javascript/jscomp/ControlFlowGraph.java
// @ closure-compiler v20260415 (narrow; AST-walk approximation only).
//
// This file exposes the `definitelyReturns` helper used by
// post-guard refinement (an if-then that definitely returns lets the
// type-checker narrow names with the negated condition for the
// statements that follow the if).
//
// The full CheckMissingReturn port (CFG + dataflow) now lives in
// checkmissingreturn.go.

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// definitelyReturns reports whether every execution path through `n`
// ends in a `return EXPR` or `throw`. Walks AST structure rather than
// a real CFG — sufficient for the post-guard narrowing used by TypeCheck.
//
// Statement-level contract:
//
//	RETURN      → true iff carries an expression
//	THROW       → true (parsed as the keyword consumeLoopHeader sink)
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
