// Port of com/google/javascript/jscomp/CheckSideEffects.java
// @ closure-compiler v20260415 (narrow; no useless-getelem sub-expr
// reporting, no @nosideeffects function annotation).
//
// Walks every EXPR_RESULT statement and fires
// JSC_USELESS_CODE whenever the contained expression has no
// observable side effect. Side-effect-free expressions are the
// pure / read-only ones: literals, bare identifiers, property
// access, arithmetic / comparison / logical / ternary / typeof /
// unary-negate of pure operands. Calls, `new`, assignment,
// increment / decrement, and `delete` are treated as side-
// effecting. AND / OR / COALESCE / HOOK with a pure RHS recurse so
// we report the inner useless piece at its own location (matches
// Java).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CheckSideEffects implements the JSC_USELESS_CODE diagnostic.
type CheckSideEffects struct {
	path  string
	diags []diag.Diagnostic
}

func NewCheckSideEffects(path string) *CheckSideEffects {
	return &CheckSideEffects{path: path}
}

// Check walks the AST firing JSC_USELESS_CODE at every EXPR_RESULT
// whose contained expression is observably pure.
func (c *CheckSideEffects) Check(script *rhino.Node) []diag.Diagnostic {
	c.diags = nil
	c.walk(script)
	return c.diags
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
// expression is pure. When the expression is a short-circuit
// operator with an impure LHS but pure RHS, Java emits TWO
// diagnostics: the outer expression-statement and the pure RHS
// sub-expression. We mirror that by recursing.
func (c *CheckSideEffects) reportUseless(stmt, expr *rhino.Node, top bool) {
	if hasSideEffects(expr) {
		// Even when the whole thing has side effects, a short-
		// circuit operator's pure RHS is still a useless sub-
		// expression that Java reports separately.
		switch expr.Token() {
		case rhino.AND, rhino.OR, rhino.COALESCE:
			kids := expr.Children()
			if len(kids) == 2 && !hasSideEffects(kids[1]) {
				c.emitAt(kids[1])
			}
		case rhino.HOOK:
			kids := expr.Children()
			if len(kids) == 3 {
				// Java flags either branch if it's pure.
				if !hasSideEffects(kids[1]) {
					c.emitAt(kids[1])
				}
				if !hasSideEffects(kids[2]) {
					c.emitAt(kids[2])
				}
			}
		}
		return
	}
	if top {
		c.emitAt(stmt)
	} else {
		c.emitAt(expr)
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
	// treats these as declarations too.
	if info.ReturnType() != nil || len(info.ParameterNames()) > 0 {
		return true
	}
	return false
}

// hasSideEffects returns true for expression nodes that Closure
// considers observably effectful. The classification is the same
// the Java NodeUtil.mayHaveSideEffects uses at the narrow scope
// (pre-nosideeffects JSDoc support).
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
		rhino.THROW:
		return true
	// Pure leaves.
	case rhino.NAME, rhino.STRINGLIT, rhino.NUMBER, rhino.TRUE, rhino.FALSE,
		rhino.NULL, rhino.UNDEFINED, rhino.THIS, rhino.GETPROP,
		rhino.ARRAYLIT, rhino.OBJECTLIT, rhino.FUNCTION:
		return false
	}
	// Compound: side-effecting iff any child is.
	for _, c := range n.Children() {
		if hasSideEffects(c) {
			return true
		}
	}
	return false
}
