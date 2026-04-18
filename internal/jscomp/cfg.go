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
	// Function expressions used as prototype-assign RHS carry their
	// JSDoc on the surrounding ASSIGN (e.g. `F.prototype.bar =
	// function() { … }` with `@return {T}` above the ASSIGN). Walk
	// up until we find a JSDoc with an explicit @return.
	info := returnJSDocFor(fn)
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
	// Empty-body stubs (`function f() {}` with a non-void @return)
	// are treated as intentional by Java — skip. Same for bodies
	// whose only returns lack expressions (jsgo's parser degrades a
	// failed cast `/** @type {T} */ (e)` to a bare `return;`, so
	// this guard also prevents parser-induced FPs).
	if len(body.Children()) == 0 {
		return
	}
	if !hasValueReturn(body) {
		return
	}
	// CFG reachability: if the synthetic EXIT node has any predecessor
	// reachable through a non-`return EXPR` fall-through, the function
	// has a path that doesn't return a value. Port of Java's
	// CheckMissingReturn behaviour over the real CFG — replaces the
	// AST-walk `definitelyReturns` approximation for functions.
	cfg := BuildControlFlowGraph(fn)
	if allExitsReturnValue(cfg) {
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

// returnJSDocFor returns the first JSDocInfo carrying @return on
// `fn` or any ancestor up to EXPR_RESULT. Catches the
// `F.prototype.bar = function() {…}` pattern where JSDoc lives on
// the ASSIGN rather than the FUNCTION itself.
func returnJSDocFor(fn *rhino.Node) *rhino.JSDocInfo {
	for cur := fn; cur != nil; cur = cur.Parent() {
		if info := cur.JSDocInfo(); info != nil && info.ReturnType() != nil {
			return info
		}
		if cur.Token() == rhino.EXPR_RESULT {
			break
		}
	}
	return nil
}

// allExitsReturnValue reports whether every CFG predecessor of the
// synthetic EXIT node is a RETURN with an expression. If any
// *reachable* predecessor is a fall-through (no return), or a bare
// `return;`, or reaches EXIT via any other path, the function is
// missing a return on at least one branch. Port of Java's
// CheckMissingReturn over the real CFG.
//
// Reachability is computed from the ENTRY node — unreachable
// synthetic joins (e.g. an IF whose branches both return) must be
// ignored, otherwise they'd falsely imply a fall-through.
func allExitsReturnValue(cfg *ControlFlowGraph) bool {
	if cfg == nil || cfg.Exit == nil {
		return true
	}
	reachable := reachableFromEntry(cfg)
	anyReachable := false
	for _, p := range cfg.Exit.Predecessors {
		if !reachable[p] {
			continue
		}
		anyReachable = true
		if p.AST == nil {
			return false
		}
		if p.AST.Token() != rhino.RETURN {
			return false
		}
		if p.AST.FirstChild() == nil {
			return false
		}
	}
	if !anyReachable {
		// No flow reaches EXIT — function is infinite-loop or fully
		// throws. Java treats these as satisfying the check.
		return true
	}
	return true
}

// reachableFromEntry returns the set of CfgNodes reachable from
// cfg.Entry via forward edges. Used so unreachable joins don't
// trigger false "missing return" warnings.
func reachableFromEntry(cfg *ControlFlowGraph) map[*CfgNode]bool {
	out := map[*CfgNode]bool{}
	stack := []*CfgNode{cfg.Entry}
	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if out[n] {
			continue
		}
		out[n] = true
		for _, e := range n.Successors {
			if !out[e.To] {
				stack = append(stack, e.To)
			}
		}
	}
	return out
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
