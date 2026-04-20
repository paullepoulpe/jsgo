// Port of com/google/javascript/jscomp/CheckSideEffects.java @ closure-compiler v20260415 (P0-A2 expansion)
//
// Walks every EXPR_RESULT statement and fires JSC_USELESS_CODE whenever
// the contained expression's value is observably discarded. Behaviour
// mirrors the Java class:
//
//   - Side-effect-free roots at statement level (literals, NAMEs,
//     GETPROP reads, arithmetic / comparison / logical / unary / in /
//     instanceof / typeof / void) are reported.
//   - GETPROP reads at statement level are always reported, even when
//     the base has side effects (`foo().bar;` → emit on `.bar`); Java's
//     diagnostic lands on the property-name token, so we emit at the
//     second child's position.
//   - AND / OR / COALESCE / HOOK with a pure RHS (or pure branches for
//     HOOK) yield a second diagnostic on the pure operand — matching
//     Java's "useless && value" / "useless ?? value" double-report.
//   - CALL is treated as effectful by default, but a CALL whose callee
//     resolves to a named function annotated `@nosideeffects` (or its
//     alias `@noside`) is reported when its result is discarded and no
//     argument has its own side effect. Mirrors Java's reliance on
//     FunctionType.isNoSideEffects().
//
// Shortcuts for `++` / `--` / ASSIGN family / `throw` / `new` / plain
// CALL are kept from the earlier narrow port.

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CheckSideEffects implements the JSC_USELESS_CODE diagnostic.
type CheckSideEffects struct {
	path  string
	diags []diag.Diagnostic
	// pureCallees maps a function-declaration NAME to the JSDocInfo
	// attached to that function's declaration. Populated in a pre-walk
	// so call-site classification can look up @nosideeffects on a
	// simple-name callee without a scope handle. Mirrors how Java
	// CheckSideEffects consults FunctionType.isNoSideEffects() on the
	// resolved callee.
	pureCallees map[string]bool
}

func NewCheckSideEffects(path string) *CheckSideEffects {
	return &CheckSideEffects{path: path}
}

// Check walks the AST firing JSC_USELESS_CODE at every EXPR_RESULT
// whose contained expression is observably pure.
func (c *CheckSideEffects) Check(script *rhino.Node) []diag.Diagnostic {
	c.diags = nil
	c.pureCallees = map[string]bool{}
	c.collectPureFunctions(script)
	c.walk(script)
	return c.diags
}

// collectPureFunctions pre-walks the script and records every named
// function declaration whose JSDocInfo carries @nosideeffects. Mirrors
// the subset of TypedScopeCreator.recordFunction Java uses to stamp the
// FunctionType's NoSideEffects bit. Closure's FUNCTION layout puts the
// binding name in the first child (empty string when anonymous).
func (c *CheckSideEffects) collectPureFunctions(n *rhino.Node) {
	if n == nil {
		return
	}
	if n.Token() == rhino.FUNCTION {
		if info := n.JSDocInfo(); info != nil && info.IsNoSideEffects() {
			if name := n.FirstChild(); name != nil && name.String_() != "" {
				c.pureCallees[name.String_()] = true
			}
		}
	}
	for _, child := range n.Children() {
		c.collectPureFunctions(child)
	}
}

func (c *CheckSideEffects) walk(n *rhino.Node) {
	if n == nil {
		return
	}
	if n.Token() == rhino.EXPR_RESULT && !isDeclarationStatement(n) {
		if inner := n.FirstChild(); inner != nil {
			c.reportUseless(n, inner)
		}
	}
	for _, child := range n.Children() {
		c.walk(child)
	}
}

// reportUseless fires JSC_USELESS_CODE at `expr`'s position when the
// expression is pure. Mirrors the double-dispatch in Java's
// CheckSideEffects — per-operator positions for IN / NOT / GETELEM /
// GETPROP, @nosideeffects-call classification for CALL, and short-
// circuit / HOOK branches get a second pass so Java's "useless &&
// value" / "useless ?? value" double-report still fires.
func (c *CheckSideEffects) reportUseless(stmt, expr *rhino.Node) {
	pure := !hasSideEffects(expr)
	switch expr.Token() {
	case rhino.IN:
		c.emitOpAt(stmt, "in")
		return
	case rhino.NOT:
		c.emitOpAt(expr, "not")
		return
	case rhino.GETELEM:
		c.emitOpAt(stmt, "getelem")
		if pure {
			return
		}
	case rhino.GETPROP:
		if kids := expr.Children(); len(kids) >= 2 {
			c.emitOpAt(kids[1], "getprop")
		} else {
			c.emitAt(expr)
		}
		return
	case rhino.CALL:
		if c.isPureCall(expr) {
			c.emitAt(stmt)
		}
		c.reportShortCircuitRHS(expr)
		return
	}
	if pure {
		c.emitAt(expr)
	}
	c.reportShortCircuitRHS(expr)
}

// reportShortCircuitRHS implements the Java "useless && value"
// double-report: for AND / OR / COALESCE / HOOK, when a non-first
// operand is side-effect-free, emit on that operand.
func (c *CheckSideEffects) reportShortCircuitRHS(expr *rhino.Node) {
	kids := expr.Children()
	switch expr.Token() {
	case rhino.AND, rhino.OR, rhino.COALESCE:
		if len(kids) == 2 && !hasSideEffects(kids[1]) {
			c.emitAt(kids[1])
		}
	case rhino.HOOK:
		if len(kids) == 3 {
			if !hasSideEffects(kids[1]) {
				c.emitAt(kids[1])
			}
			if !hasSideEffects(kids[2]) {
				c.emitAt(kids[2])
			}
		}
	}
}

// isPureCall reports whether `call` is a call to a function resolved to
// @nosideeffects with no effectful arguments. Matches the Java fast
// path in NodeUtil.functionCallHasSideEffects: a CALL is pure iff the
// callee has NoSideEffects on its FunctionType / JSDocInfo AND none of
// the arguments individually have side effects.
//
// Guard: only classify on named callees; dynamic / anonymous callees
// can't be traced back to a JSDoc declaration.
func (c *CheckSideEffects) isPureCall(call *rhino.Node) bool {
	kids := call.Children()
	if len(kids) == 0 {
		return false
	}
	callee := kids[0]
	if callee.Token() != rhino.NAME {
		return false
	}
	if !c.pureCallees[callee.String_()] {
		return false
	}
	for _, arg := range kids[1:] {
		if hasSideEffects(arg) {
			return false
		}
	}
	return true
}

func (c *CheckSideEffects) emitAt(n *rhino.Node) {
	c.diags = append(c.diags, diag.Diagnostic{
		Path:    c.path,
		Line:    n.Lineno(),
		Col:     n.Charno(),
		Code:    diag.JSC_USELESS_CODE,
		Level:   diag.Warning,
		Message: "Suspicious code. This code lacks side-effects. Is there a bug?",
	})
}

// emitOpAt emits the operator-specific "The result of the 'foo'
// operator is not being used." variant that Java uses for
// getprop / in / not / getelem top-level expressions.
func (c *CheckSideEffects) emitOpAt(n *rhino.Node, op string) {
	c.diags = append(c.diags, diag.Diagnostic{
		Path:    c.path,
		Line:    n.Lineno(),
		Col:     n.Charno(),
		Code:    diag.JSC_USELESS_CODE,
		Level:   diag.Warning,
		Message: "Suspicious code. The result of the '" + op + "' operator is not being used.",
	})
}

// isDeclarationStatement reports whether `n` (an EXPR_RESULT) is
// actually a prototype-slot / property-type declaration rather than
// a real expression statement. Closure treats `Foo.prototype.bar;`
// (with an @type JSDoc) as a declaration of a slot — Java's
// CheckSideEffects skips these. We do the same by checking for an
// @type / @const / @record-style JSDoc on the statement and a
// GETPROP inner.
func isDeclarationStatement(stmt *rhino.Node) bool {
	info := stmt.JSDocInfo()
	if info == nil {
		if inner := stmt.FirstChild(); inner != nil {
			info = inner.JSDocInfo()
		}
	}
	if info == nil {
		return false
	}
	if info.HasType() || info.IsConstructor() || info.IsInterface() || info.IsRecord() || info.TypedefType() != nil {
		return true
	}
	// `/** @param {T} x */ obj.prop;` / `/** @return {T} */
	// obj.prop;` declare the callable signature of a prototype
	// slot without an explicit value. Java's CheckSideEffects
	// treats these as declarations too — but only when the inner
	// expression is a property access. A top-level FUNCTION /
	// paren-wrapped FUNCTION with the same JSDoc is still a
	// useless-expression statement (the function is never called).
	if info.ReturnType() != nil || len(info.ParameterNames()) > 0 {
		if inner := stmt.FirstChild(); inner != nil && inner.Token() == rhino.GETPROP {
			return true
		}
	}
	return false
}

// hasSideEffects mirrors NodeUtil.mayHaveSideEffects: a node is
// effectful iff its own operator writes / calls / throws, or any
// child is. GETPROP is pure at its own level but recurses into its
// base (so `foo().bar` still reports as effectful because of the
// CALL).
func hasSideEffects(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	switch n.Token() {
	// Definitely-effectful: calls may do anything; new constructs an
	// object; assignments mutate a binding; increment / decrement
	// mutate a LHS; delete mutates a property set.
	case rhino.CALL, rhino.NEW, rhino.DELPROP,
		rhino.ASSIGN, rhino.ASSIGN_ADD, rhino.ASSIGN_SUB, rhino.ASSIGN_MUL,
		rhino.ASSIGN_DIV, rhino.ASSIGN_MOD, rhino.ASSIGN_BITAND,
		rhino.ASSIGN_BITOR, rhino.ASSIGN_BITXOR, rhino.ASSIGN_LSH,
		rhino.ASSIGN_RSH, rhino.ASSIGN_URSH,
		rhino.ASSIGN_AND, rhino.ASSIGN_OR, rhino.ASSIGN_COALESCE,
		rhino.THROW:
		return true
	// Pure leaves.
	case rhino.NAME, rhino.STRINGLIT, rhino.NUMBER, rhino.TRUE, rhino.FALSE,
		rhino.NULL, rhino.UNDEFINED, rhino.THIS,
		rhino.ARRAYLIT, rhino.OBJECTLIT, rhino.FUNCTION:
		return false
	}
	// Compound (including GETPROP): effectful iff any child is.
	for _, c := range n.Children() {
		if hasSideEffects(c) {
			return true
		}
	}
	return false
}
