// Port of com/google/javascript/jscomp/AstManipulations.java
// @ closure-compiler v20260415 (subset; FuseExpressions stub pending COMMA token).
//
// AstManipulations provides utilities for common AST rewrites.
// FuseExpressions is the public API called by transform passes; its
// body requires a COMMA token not yet present in rhino/token.go — the
// stub will be filled in when the token is added.

package jscomp

import "github.com/paullepoulpe/jsgo/legacy/internal/rhino"

// FuseExpressions returns a single node equivalent to executing
// <exp1, exp2>. Both exp1 and exp2 must be detached (no parent).
// If exp2 has no children (a stand-in for Java's EMPTY), exp1 is
// returned unchanged.
//
// Mirrors AstManipulations.fuseExpressions 1:1. Body is a stub
// until the COMMA token lands in internal/rhino/token.go.
func FuseExpressions(exp1, exp2 *rhino.Node) *rhino.Node {
	if !exp2.HasChildren() && exp2.Token() == rhino.ERROR {
		return exp1
	}
	// TODO: construct comma node once rhino.COMMA token is defined.
	_ = exp2
	return exp1
}
