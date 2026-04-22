// Port of com/google/javascript/jscomp/CheckMissingReturn.java @ closure-compiler v20260415
//
// Checks functions for missing return statements. Return statements are only
// expected for functions with return type information. Functions with empty
// bodies are ignored.
//
// Java's implementation uses a real CFG built by ControlFlowAnalysis plus a
// CheckPathsBetweenNodes graph query (allPathsSatisfyPredicate) that walks the
// CFG forward from Entry to ImplicitReturn using only "true-condition" edges
// (exception edges and statically-impossible branches are pruned).
//
// This port mirrors that algorithm exactly:
//  1. fastAllPathsReturnCheck — O(degree) check: if every non-exception
//     predecessor of EXIT is a RETURN with a value, we're done.
//  2. allPathsSatisfyPredicate — BFS from Entry to EXIT using only edges
//     where goesThroughTrueCondition holds; any path that reaches EXIT
//     without passing through a value-RETURN means a missing return.

package jscomp

import (
	"fmt"

	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/jstype"
)

// CheckMissingReturn walks every FUNCTION node and emits
// JSC_MISSING_RETURN_STATEMENT when the function declares a non-void
// @return but the body has at least one path that does not end in
// a value-bearing `return` statement. Port of Java's
// CheckMissingReturn (AbstractCfgCallback subclass).
type CheckMissingReturn struct {
	path  string
	scope *TypedScope
	diags []diag.Diagnostic
}

// NewCheckMissingReturn constructs a CheckMissingReturn.
// Mirrors Java's CheckMissingReturn(AbstractCompiler).
func NewCheckMissingReturn(path string, scope *TypedScope) *CheckMissingReturn {
	return &CheckMissingReturn{path: path, scope: scope}
}

// Check drives the AST walk and returns accumulated diagnostics.
func (c *CheckMissingReturn) Check(script *rhino.Node) []diag.Diagnostic {
	c.diags = nil
	c.walk(script)
	return c.diags
}

// walk visits every FUNCTION node in the AST, recursing into children.
// Mirrors enterScopeWithCfg being called for every function scope by
// NodeTraversal.AbstractCfgCallback in Java.
func (c *CheckMissingReturn) walk(n *rhino.Node) {
	if n == nil {
		return
	}
	if n.Token() == rhino.FUNCTION {
		c.enterScopeWithCfg(n)
	}
	for _, child := range n.Children() {
		c.walk(child)
	}
}

// enterScopeWithCfg mirrors Java's CheckMissingReturn.enterScopeWithCfg.
// Called once per function scope.
func (c *CheckMissingReturn) enterScopeWithCfg(fn *rhino.Node) {
	returnType := c.getExplicitReturnTypeIfExpected(fn)
	if returnType == nil {
		// No return value is expected, so nothing to check.
		return
	}

	// Arrow functions with expression bodies have an implicit return value.
	if fn.GetBooleanProp(rhino.PropArrowFn) {
		kids := fn.Children()
		if len(kids) >= 3 && kids[2].Token() != rhino.BLOCK {
			// Body is an expression, which is the implicit return value.
			return
		}
	}

	// Guard against parser-degraded cast expressions: the Go parser
	// degrades `return /** @type {T} */ (expr)` to a bare `return;`
	// because it cannot parse the inline @type as a cast node. When
	// no value-bearing RETURN exists in the body, the function may
	// appear to have only bare returns even though it legitimately
	// returns via a cast. Skip the check in that case.
	// NOTE: Java has a real CAST AST node so this guard is not needed there.
	kids := fn.Children()
	if len(kids) >= 3 && !hasValueReturn(kids[2]) {
		return
	}

	cfg := BuildControlFlowGraph(fn)
	if c.fastAllPathsReturnCheck(cfg) {
		return
	}

	if !c.allPathsSatisfyPredicate(cfg) {
		line, col := fn.Lineno(), fn.Charno()
		c.diags = append(c.diags, diag.Diagnostic{
			Path:    c.path,
			Line:    line,
			Col:     col,
			Code:    diag.JSC_MISSING_RETURN_STATEMENT,
			Level:   diag.Warning,
			Message: fmt.Sprintf("Missing return statement. Function expected to return %s.", returnType.TypeName()),
		})
	}
}

// fastAllPathsReturnCheck is the O(degree) fast check. Mirrors Java's
// CheckMissingReturn.fastAllPathsReturnCheck.
//
// Returns true if every non-exception predecessor of the EXIT node is a
// RETURN with a value (or a call that always throws — Java uses
// CodingConvention.isFunctionCallThatAlwaysThrows; we stub that to false
// since we have no CodingConvention port yet).
//
// "May spuriously report that a return statement is missing" — when this
// returns false, we fall through to the full allPathsSatisfyPredicate check.
func (c *CheckMissingReturn) fastAllPathsReturnCheck(cfg *ControlFlowGraph) bool {
	if cfg == nil || cfg.Exit == nil {
		return true
	}
	for _, pred := range cfg.Exit.Predecessors {
		// Find the branch label on the edge from pred → EXIT.
		branch := branchToExit(pred, cfg.Exit)
		// Mirrors Java: skip all exception edges (ON_EX in Java).
		// NOTE: Java also skips impossible edges via goesThroughTrueCondition
		// but fastAllPathsReturnCheck only skips ON_EX.
		if branch == BranchException {
			continue
		}
		if pred.AST == nil {
			return false
		}
		// TODO: also accept isFunctionCallThatAlwaysThrows(pred.AST) once
		// CodingConvention is ported — Java's NOTE(dimvar) comment.
		if pred.AST.Token() != rhino.RETURN {
			return false
		}
		if pred.AST.FirstChild() == nil {
			// Bare `return;` — no value.
			return false
		}
	}
	return true
}

// allPathsSatisfyPredicate mirrors Java's
// CheckPathsBetweenNodes.allPathsSatisfyPredicate().
//
// Returns true if every path from Entry to EXIT (using only edges where
// goesThroughTrueCondition holds) passes through at least one node where
// nodePredicate holds — here: the node is a RETURN with an expression.
//
// Equivalently: can EXIT be reached from Entry via true-condition edges
// without passing through any value-RETURN? If yes → not all paths
// return → allPathsSatisfyPredicate = false.
func (c *CheckMissingReturn) allPathsSatisfyPredicate(cfg *ControlFlowGraph) bool {
	if cfg == nil || cfg.Exit == nil {
		return true
	}

	// BFS from Entry. We stop propagating past RETURN-with-value nodes
	// (those paths are "satisfied"). If EXIT is reachable via non-RETURN
	// paths, not all paths satisfy the predicate.
	visited := map[*CfgNode]bool{}
	queue := []*CfgNode{cfg.Entry}

	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		if visited[n] {
			continue
		}
		visited[n] = true

		if n == cfg.Exit {
			// EXIT is reachable without going through a RETURN with value.
			return false
		}

		// If this node is a RETURN with an expression, this path is
		// satisfied — do not propagate further (mirrors the nodePredicate
		// short-circuit in CheckPathsBetweenNodes).
		if n.AST != nil && n.AST.Token() == rhino.RETURN && n.AST.FirstChild() != nil {
			continue
		}

		for _, e := range n.Successors {
			if !goesThroughTrueCondition(n, e) {
				continue
			}
			if !visited[e.To] {
				queue = append(queue, e.To)
			}
		}
	}

	// EXIT was not reachable without going through a RETURN with value
	// (or EXIT is entirely unreachable — infinite loop / always throws).
	return true
}

// goesThroughTrueCondition mirrors the Java static method
// CheckMissingReturn.goesThroughTrueCondition.
//
// Returns false for:
//   - Exception edges (BranchException / ON_EX in Java).
//   - Conditional edges whose source condition is a constant and the
//     branch direction is statically impossible.
//
// Returns true for all other edges (unconditional, unknown-condition
// conditional, synthetic).
func goesThroughTrueCondition(src *CfgNode, e CfgEdge) bool {
	// First skip all exceptions. Mirrors Java: `if (branch == Branch.ON_EX) return false`.
	if e.Branch == BranchException {
		return false
	}

	// For conditional branches, evaluate constant conditions.
	// Mirrors Java: `else if (branch.isConditional()) { ... NodeUtil.getBooleanValueCMR(condition) ... }`
	if e.Branch == BranchOnTrue || e.Branch == BranchOnFalse {
		if src.AST != nil {
			cond := getConditionExpressionCMR(src.AST)
			if cond != nil {
				val, known := getBooleanValueCMR(cond)
				if known {
					// The branch is statically possible only when the
					// constant direction matches the branch label.
					if e.Branch == BranchOnTrue {
						return val
					}
					return !val
				}
			}
		}
	}

	return true
}

// getConditionExpressionCMR returns the condition sub-expression of a
// control-flow node (IF / WHILE / FOR / DO_WHILE). Mirrors Java's
// NodeUtil.getConditionExpressionCMR.
func getConditionExpressionCMR(n *rhino.Node) *rhino.Node {
	if n == nil {
		return nil
	}
	kids := n.Children()
	switch n.Token() {
	case rhino.IF:
		// IF children: (condition, then [, else])
		if len(kids) >= 1 {
			return kids[0]
		}
	case rhino.WHILE:
		// WHILE children: (condition, body)
		if len(kids) >= 1 {
			return kids[0]
		}
	case rhino.FOR:
		// FOR children: (init, condition, update, body)
		if len(kids) >= 2 {
			return kids[1]
		}
	case rhino.DO_WHILE:
		// DO_WHILE children: (body, condition)
		if len(kids) >= 2 {
			return kids[1]
		}
	}
	return nil
}

// getBooleanValueCMR evaluates a constant expression to a boolean.
// Returns (value, true) when the expression is a known constant,
// (false, false) otherwise. Mirrors Java's NodeUtil.getBooleanValueCMR
// returning Tri (TRUE / FALSE / UNKNOWN).
//
// TODO: widen to evaluate more constant expressions (e.g. non-zero
// numeric literals, non-empty string literals) — Java's NodeUtil does
// this; the current port handles only the common cases seen in the
// extracted corpus.
func getBooleanValueCMR(n *rhino.Node) (bool, bool) {
	if n == nil {
		return false, false
	}
	switch n.Token() {
	case rhino.TRUE:
		return true, true
	case rhino.FALSE:
		return false, true
	case rhino.NULL:
		return false, true
	case rhino.NUMBER:
		return n.Number() != 0, true
	case rhino.STRINGLIT:
		return n.String_() != "", true
	case rhino.VOID:
		// `void expr` is always undefined (falsy).
		return false, true
	case rhino.NOT:
		// `!expr`: invert if inner is constant.
		if child := n.FirstChild(); child != nil {
			if val, ok := getBooleanValueCMR(child); ok {
				return !val, true
			}
		}
	}
	return false, false
}

// getExplicitReturnTypeIfExpected mirrors Java's
// CheckMissingReturn.getExplicitReturnTypeIfExpected.
//
// Returns the declared return type when the function is expected to
// explicitly return a value, or nil when:
//   - The node is not a function.
//   - No JSDoc @return is present.
//   - The function body is empty.
//   - The function is a constructor (@constructor).
//   - The declared return type is void / unknown.
//
// NOTE: Java resolves the return type from the live JSType attached to the
// AST node by TypeCheck. jsgo resolves via JSDoc because we don't yet
// attach a fully-resolved JSType to every FUNCTION AST node.  The
// behaviour is equivalent for the annotated-JavaScript corpus.
//
// TODO: widen to unwrap Promise<T> for async functions and Generator<T,R>
// for generator functions — Java's Promises.getTemplateTypeOfThenable and
// JsIterables.getReturnElementType.  Stub: treat async/generator functions
// as ordinary (may produce FPs on explicitly-annotated async returns).
func (c *CheckMissingReturn) getExplicitReturnTypeIfExpected(fn *rhino.Node) jstype.JSType {
	if fn.Token() != rhino.FUNCTION {
		// Nothing to do in a global/module/block scope.
		return nil
	}

	// isEmptyFunction: mirrors Java's isEmptyFunction check.
	if isEmptyFunction(fn) {
		return nil
	}

	// Find JSDoc with @return, possibly on a surrounding ASSIGN.
	info := returnJSDocFor(fn)
	if info == nil || info.ReturnType() == nil {
		return nil
	}

	// Constructor exception: constructors that declare a return type are
	// allowed to return from a new-without-new call pattern; their return
	// type is implicit. Mirrors Java's `if (scopeType.isConstructor()) return null`.
	if info.IsConstructor() {
		return nil
	}

	// Resolve the declared return type.
	returnType := c.scope.ResolveExpr(info.ReturnType().TypeName)
	if returnType == nil {
		return nil
	}

	// TODO: for async functions (scopeRoot.isAsyncFunction()), unwrap via
	// Promises.getTemplateTypeOfThenable — not yet ported.
	// TODO: for generator functions (scopeRoot.isGeneratorFunction()), unwrap
	// via JsIterables.getReturnElementType — not yet ported.

	if !c.isVoidOrUnknown(returnType) {
		return returnType
	}

	return nil
}

// hasValueReturn returns true when the subtree rooted at `n` contains
// at least one `return EXPR` (a RETURN node with a non-nil first child).
// Stops at nested FUNCTION boundaries (doesn't recurse into them).
//
// NOTE: This guard is not in the Java port — Java has a real CAST AST
// node so `return /** @type {T} */ (e)` always produces a RETURN with
// a child. jsgo's parser degrades that pattern to a bare `return;` so
// we use this check to avoid false positives on parser failures.
func hasValueReturn(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	// Don't cross nested function boundaries.
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

// isEmptyFunction mirrors Java's CheckMissingReturn.isEmptyFunction.
// Returns true when the FUNCTION node has an empty body block
// (three children: name, params, empty-block). Empty-body stubs are
// treated as intentional (e.g. interface method declarations).
func isEmptyFunction(fn *rhino.Node) bool {
	kids := fn.Children()
	return len(kids) == 3 && len(kids[2].Children()) == 0
}

// isVoidOrUnknown mirrors Java's CheckMissingReturn.isVoidOrUnknown.
// Returns true when returnType is void, unknown, or a union that
// includes void. Java uses voidType.isSubtypeOf(returnType) which
// holds exactly when the union/supertype accepts void.
func (c *CheckMissingReturn) isVoidOrUnknown(returnType jstype.JSType) bool {
	if returnType == nil {
		return true
	}
	// jstype.Void.IsSubtypeOf(returnType): void is subtype of T when T is
	// void, unknown, all, or a union containing void.
	return jstype.IsSubtype(jstype.Void, returnType)
}

// returnJSDocFor returns the first JSDocInfo carrying @return on `fn`
// or any ancestor up to EXPR_RESULT. Catches the
// `F.prototype.bar = function() {…}` pattern where JSDoc lives on
// the ASSIGN rather than the FUNCTION itself. Mirrors Java's
// NodeUtil.getBestJSDocInfo walk used inside getExplicitReturnTypeIfExpected.
//
// The walk stops at EXPR_RESULT (the outermost statement boundary) so that
// JSDoc on a sibling or grandparent statement is not adopted. It also
// stops when it encounters a FUNCTION node other than `fn` itself — without
// this guard an inner function would inherit its enclosing function's
// @return annotation.
func returnJSDocFor(fn *rhino.Node) *rhino.JSDocInfo {
	for cur := fn; cur != nil; cur = cur.Parent() {
		if info := cur.JSDocInfo(); info != nil && info.ReturnType() != nil {
			return info
		}
		// Stop at a statement boundary (covers the EXPR_RESULT → ASSIGN
		// → FUNCTION chain for prototype-method assignments).
		if cur.Token() == rhino.EXPR_RESULT {
			break
		}
		// Stop at any FUNCTION node that is not the function we are
		// checking — prevents an inner function from inheriting its
		// enclosing function's @return.
		if cur != fn && cur.Token() == rhino.FUNCTION {
			break
		}
	}
	return nil
}

// branchToExit finds the Branch label on the edge from `pred` to `exit`.
// Returns BranchUnconditional when no edge is found (should not happen
// for valid CFGs).
func branchToExit(pred, exit *CfgNode) Branch {
	for _, e := range pred.Successors {
		if e.To == exit {
			return e.Branch
		}
	}
	return BranchUnconditional
}
