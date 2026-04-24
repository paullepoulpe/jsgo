// Port of com/google/javascript/jscomp/CheckSideEffects.java
// @ closure-compiler v20260415.
//
// Walks every EXPR_RESULT statement and fires
// JSC_USELESS_CODE whenever the contained expression has no
// observable side effect. Side-effect-free expressions are the
// pure / read-only ones: literals, bare identifiers, property
// access, arithmetic / comparison / logical / ternary / typeof /
// unary-negate of pure operands. Calls, `new`, assignment,
// increment / decrement, and `delete` are treated as side-
// effecting — unless the callee carries @nosideeffects (mirroring
// Java's GetNoSideEffectExterns + @nosideeffects-aware CALL
// classification). AND / OR / COALESCE / HOOK with a pure RHS
// recurse so we report the inner useless piece at its own location
// (matches Java). Sub-expression `a[b].c;` is reported at the
// outermost expression position per the Java emitOpAt convention.

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CheckSideEffects implements the JSC_USELESS_CODE diagnostic.
type CheckSideEffects struct {
	path     string
	diags    []diag.Diagnostic
	noSEFns  map[string]bool // qualified names of @nosideeffects functions
	analyzer *AstAnalyzer
}

func NewCheckSideEffects(path string) *CheckSideEffects {
	return &CheckSideEffects{path: path, analyzer: NewAstAnalyzer()}
}

// Check walks the AST firing JSC_USELESS_CODE at every EXPR_RESULT
// whose contained expression is observably pure.
// Mirrors Java's process(Node externs, Node root):
//  1. GetNoSideEffectExterns sub-pass collects @nosideeffects names.
//  2. Main traversal fires diagnostics.
func (c *CheckSideEffects) Check(script *rhino.Node) []diag.Diagnostic {
	c.diags = nil
	c.noSEFns = make(map[string]bool)
	c.collectNoSideEffectFns(script)
	c.walk(script)
	return c.diags
}

// collectNoSideEffectFns mirrors Java's GetNoSideEffectExterns —
// scans the AST for FUNCTION nodes annotated @nosideeffects and
// records their qualified names so CALL classification can treat
// them as pure.
func (c *CheckSideEffects) collectNoSideEffectFns(n *rhino.Node) {
	if n == nil {
		return
	}
	if n.Token() == rhino.FUNCTION {
		info := n.JSDocInfo()
		if info == nil {
			if p := n.Parent(); p != nil {
				info = p.JSDocInfo()
			}
		}
		if info != nil && info.IsNoSideEffects() {
			if name := NU.GetFunctionName(n); name != "" {
				c.noSEFns[name] = true
			}
		}
	}
	for _, child := range n.Children() {
		c.collectNoSideEffectFns(child)
	}
}

func (c *CheckSideEffects) walk(n *rhino.Node) {
	if n == nil {
		return
	}
	if n.Token() == rhino.EXPR_RESULT {
		if !isDeclarationStatement(n) {
			if inner := n.FirstChild(); inner != nil {
				c.reportUseless(n, inner, true /*top-level*/)
			}
		}
	}
	for _, child := range n.Children() {
		c.walk(child)
	}
}

// reportUseless fires JSC_USELESS_CODE at `expr`'s position when the
// expression is pure. Mirrors the double-dispatch in Java's
// CheckSideEffects — when the outer form is a short-circuit
// operator or HOOK, the pure sub-expressions are ALSO reported
// individually, in addition to the whole-statement diagnostic when
// the outer form is fully pure.
func (c *CheckSideEffects) reportUseless(stmt, expr *rhino.Node, top bool) {
	pure := !c.hasSideEffects(expr)
	// Per-operator "result not used" diagnostics fire even when the
	// operand has side effects — the operator itself is dropped.
	// Java emits at specific positions for each operator kind.
	if top {
		switch expr.Token() {
		case rhino.IN:
			c.emitOpAt(stmt, "in")
			return
		case rhino.NOT:
			c.emitOpAt(expr, "not")
			return
		case rhino.BITNOT:
			c.emitOpAt(expr, "bitnot")
			return
		case rhino.GETELEM:
			c.emitOpAt(stmt, "getelem")
			if pure {
				return
			}
			// fall through to recursion for ternary/and/or branches
		}
	}
	if pure {
		if expr.Token() == rhino.GETPROP {
			kids := expr.Children()
			if len(kids) >= 2 {
				c.emitOpAt(kids[1], "getprop")
			} else {
				c.emitAt(expr)
			}
		} else {
			c.emitAt(expr)
		}
	}
	_ = stmt
	// Recurse into short-circuit / ternary branches so Java's
	// per-branch pure-RHS warnings still fire even when the outer
	// statement was already flagged.
	switch expr.Token() {
	case rhino.AND, rhino.OR, rhino.COALESCE:
		kids := expr.Children()
		if len(kids) == 2 && !c.hasSideEffects(kids[1]) {
			c.emitAt(kids[1])
		}
	case rhino.HOOK:
		kids := expr.Children()
		if len(kids) == 3 {
			if !c.hasSideEffects(kids[1]) {
				c.emitAt(kids[1])
			}
			if !c.hasSideEffects(kids[2]) {
				c.emitAt(kids[2])
			}
		}
	}
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

// hasSideEffects returns true for expression nodes that Closure
// considers observably effectful. The classification consults
// AstAnalyzer for constructor calls (so new Date() is pure) and
// the noSEFns set for @nosideeffects-annotated function calls —
// mirroring Java's AstAnalyzer.mayHaveSideEffects + the
// GetNoSideEffectExterns pass that populates noSideEffectExterns.
func (c *CheckSideEffects) hasSideEffects(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	switch n.Token() {
	// Definitely-effectful: assignments mutate a binding;
	// increment / decrement mutate a LHS; delete mutates a property
	// set; throw interrupts control flow; dynamic import is async.
	case rhino.ASSIGN, rhino.ASSIGN_ADD, rhino.ASSIGN_SUB, rhino.ASSIGN_MUL,
		rhino.ASSIGN_DIV, rhino.ASSIGN_MOD, rhino.ASSIGN_BITAND,
		rhino.ASSIGN_BITOR, rhino.ASSIGN_BITXOR, rhino.ASSIGN_LSH,
		rhino.ASSIGN_RSH, rhino.ASSIGN_URSH,
		rhino.ASSIGN_AND, rhino.ASSIGN_OR, rhino.ASSIGN_COALESCE,
		rhino.DELPROP, rhino.THROW, rhino.DYNAMIC_IMPORT:
		return true
	// CALL / OPTCHAIN_CALL: pure only if callee is @nosideeffects.
	// Mirrors Java's noSideEffectExterns set + AstAnalyzer check.
	// OPTCHAIN_CALL is treated like CALL — it may call a function.
	case rhino.CALL, rhino.OPTCHAIN_CALL:
		callee := n.FirstChild()
		if callee != nil {
			name := callee.GetQualifiedName()
			if name != "" && c.noSEFns[name] {
				break // callee is @nosideeffects → check children
			}
		}
		return true
	// NEW: use AstAnalyzer.ConstructorCallHasSideEffects so that
	// known-pure constructors like new Date() are treated as pure
	// (mirrors Java's ModifiedStateChecker / AstAnalyzer).
	case rhino.NEW:
		if c.analyzer.ConstructorCallHasSideEffects(n) {
			return true
		}
		// Pure constructor: side-effecting only if args are.
	// Pure leaves.
	case rhino.NAME, rhino.STRINGLIT, rhino.NUMBER, rhino.TRUE, rhino.FALSE,
		rhino.NULL, rhino.UNDEFINED, rhino.THIS, rhino.GETPROP,
		rhino.ARRAYLIT, rhino.OBJECTLIT, rhino.FUNCTION:
		return false
	}
	// Compound: side-effecting iff any child is.
	for _, child := range n.Children() {
		if c.hasSideEffects(child) {
			return true
		}
	}
	return false
}
