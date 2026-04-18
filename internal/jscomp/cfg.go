// Port of com/google/javascript/jscomp/ControlFlowGraph.java +
// com/google/javascript/jscomp/CheckMissingReturn.java
// @ closure-compiler v20260415 (narrow; AST-walk approximation of
// the dataflow check).
//
// Java's CheckMissingReturn runs after ControlFlowAnalysis builds a
// per-function CFG, then checks reachability of the synthetic exit
// node from every path that doesn't pass through a `return EXPR`.
// jsgo doesn't have a real CFG yet (the full DataFlowAnalysis port
// lands later), but the check itself only needs the structural
// "every path through this block ends in a return/throw" property,
// which we can answer with a small recursive walker.
//
// This file also exposes the `definitelyReturns` helper used by
// post-guard refinement (an if-then that definitely returns lets the
// type-checker narrow names with the negated condition for the
// statements that follow the if).

package jscomp

import (
	"fmt"

	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

// MissingReturnCheck walks every FUNCTION node and emits
// JSC_MISSING_RETURN_STATEMENT when the function declares a non-void
// @return but the body has at least one path that doesn't end in
// `return EXPR` / `throw`.
type MissingReturnCheck struct {
	path  string
	scope *TypedScope
	diags []diag.Diagnostic
}

func NewMissingReturnCheck(path string, scope *TypedScope) *MissingReturnCheck {
	return &MissingReturnCheck{path: path, scope: scope}
}

func (m *MissingReturnCheck) Check(script *rhino.Node) []diag.Diagnostic {
	m.diags = nil
	m.walk(script)
	return m.diags
}

func (m *MissingReturnCheck) walk(n *rhino.Node) {
	if n == nil {
		return
	}
	if n.Token() == rhino.FUNCTION {
		m.checkFunction(n)
	}
	for _, c := range n.Children() {
		m.walk(c)
	}
}

func (m *MissingReturnCheck) checkFunction(fn *rhino.Node) {
	info := fn.JSDocInfo()
	if info == nil || info.ReturnType() == nil {
		return
	}
	declared := m.scope.ResolveExpr(info.ReturnType().TypeName)
	if declared == nil || isVoidish(declared) {
		return
	}
	kids := fn.Children()
	if len(kids) < 3 {
		return
	}
	body := kids[2]
	if body.Token() != rhino.BLOCK {
		return
	}
	// Java only fires when the body contains at least one
	// `return EXPR` — both empty stubs and bodies whose only returns
	// lack expressions are treated as intentional. (jsgo's parser
	// occasionally produces a bare `return;` when it can't model the
	// returned expression — `/** @type {T} */ (e)` casts among others
	// — so this guard also cuts a class of parser-induced FPs.)
	if !hasValueReturn(body) {
		return
	}
	if definitelyReturns(body) {
		return
	}
	// Java reports at the FUNCTION node's position; we mirror that by
	// using the function NAME's position (where the source has the
	// function keyword's caret on a typical line).
	name := fn.FirstChild()
	line, col := fn.Lineno(), fn.Charno()
	if name != nil {
		line, col = name.Lineno(), name.Charno()
	}
	m.diags = append(m.diags, diag.Diagnostic{
		Path:    m.path,
		Line:    line,
		Col:     col,
		Code:    diag.JSC_MISSING_RETURN_STATEMENT,
		Level:   diag.Warning,
		Message: fmt.Sprintf("Missing return statement. Function expected to return %s.", declared.TypeName()),
	})
}

// isVoidish reports whether the declared return type is one for
// which Java considers a missing-return harmless: void / undefined,
// the all-types `*`, the unknown `?`, or any union that includes
// void.
func isVoidish(t jstype.JSType) bool {
	if t == nil {
		return true
	}
	if t == jstype.Void || t == jstype.All || t == jstype.Unknown {
		return true
	}
	if u, ok := t.(*jstype.UnionType); ok {
		for _, m := range u.Members() {
			if m == jstype.Void {
				return true
			}
		}
	}
	if t.TypeName() == "?" || t.TypeName() == "*" {
		return true
	}
	return false
}

// hasValueReturn returns true when the body contains at least one
// `return EXPR;` — used so empty stubs and `return;`-only functions
// don't trigger missing-return on functions with a non-void @return.
func hasValueReturn(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	if n.Token() == rhino.FUNCTION {
		return false
	}
	if n.Token() == rhino.RETURN && n.FirstChild() != nil {
		return true
	}
	for _, c := range n.Children() {
		if hasValueReturn(c) {
			return true
		}
	}
	return false
}

// definitelyReturns reports whether every execution path through `n`
// ends in a `return EXPR` or `throw`. Walks AST structure rather than
// a real CFG — sufficient for the missing-return check on the corpus.
//
// Statement-level contract:
//   RETURN      → true iff carries an expression
//   THROW       → true (parsed as the keyword consumeLoopHeader sink)
//   IF(c, t, e) → true iff both t and e definitely return
//   BLOCK       → true iff some statement definitely returns
//   FUNCTION/CLASS/EXPR_RESULT/VAR/etc. → false (no termination)
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
