// Port of com/google/javascript/jscomp/CheckSuspiciousCode.java @ closure-compiler v20260415 (narrow subset — SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR + SUSPICIOUS_EXPR only).
//
// Java's CheckSuspiciousCode is a grab-bag of lint checks against
// common developer mistakes. This port covers just two:
//
//   - JSC_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR — flags literals
//     on the LHS of `||` / `&&` whose boolean value is fixed (`null`,
//     `undefined`, `0`, `NaN`, `false`, `""`). Java's version uses
//     type information in addition to syntactic literal detection; we
//     restrict to the syntactic half to stay narrow.
//
//   - JSC_SUSPICIOUS_EXPR — flags a bare ASSIGN (`x = y`) used
//     directly as the condition of IF / WHILE / DO / FOR / HOOK. Java
//     fires this from a broader dispatch (also covers SEMICOLON,
//     NaN-comparison, bad `in` / `instanceof`); the narrow port keeps
//     only the assign-in-condition form.
//
// Other diagnostics in the Java class (SUSPICIOUS_SEMICOLON,
// SUSPICIOUS_NAN, SUSPICIOUS_IN, SUSPICIOUS_INSTANCEOF_LEFT,
// SUSPICIOUS_NEGATED_LEFT_OPERAND_OF_IN_OPERATOR) are deferred.

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CheckSuspiciousCode implements the narrow subset described above.
type CheckSuspiciousCode struct {
	path  string
	diags []diag.Diagnostic
}

func NewCheckSuspiciousCode(path string) *CheckSuspiciousCode {
	return &CheckSuspiciousCode{path: path}
}

// Check walks the AST emitting the two suspicious-code diagnostics.
func (c *CheckSuspiciousCode) Check(script *rhino.Node) []diag.Diagnostic {
	c.diags = nil
	c.walk(script)
	return c.diags
}

func (c *CheckSuspiciousCode) walk(n *rhino.Node) {
	if n == nil {
		return
	}
	switch n.Token() {
	case rhino.AND, rhino.OR:
		c.checkLeftOperandOfLogicalOperator(n)
	case rhino.IF, rhino.WHILE, rhino.DO_WHILE, rhino.HOOK:
		// The condition is the first child for IF / WHILE / HOOK;
		// for DO_WHILE Closure lays out the body first then the test,
		// but our narrow check simply inspects every child for a
		// bare ASSIGN directly in the conditional slot. Rather than
		// hard-code per-token offsets, scan the ASSIGN candidates
		// here — since the body is a BLOCK, it won't match anyway.
		for _, child := range n.Children() {
			if child.Token() == rhino.ASSIGN {
				c.emitSuspiciousExpr(child)
			}
		}
	case rhino.FOR:
		// FOR(init, cond, update, body). Only the cond slot (index 1)
		// gets flagged; init and update legitimately assign.
		kids := n.Children()
		if len(kids) >= 2 && kids[1].Token() == rhino.ASSIGN {
			c.emitSuspiciousExpr(kids[1])
		}
	}
	for _, child := range n.Children() {
		c.walk(child)
	}
}

// checkLeftOperandOfLogicalOperator fires when the LHS of `&&` / `||`
// is a literal with a deterministic boolean value. Java's version
// also consults the static type to reject always-truthy / always-falsy
// expressions; the narrow port stays syntactic so it can run without
// type information.
func (c *CheckSuspiciousCode) checkLeftOperandOfLogicalOperator(n *rhino.Node) {
	kids := n.Children()
	if len(kids) < 2 {
		return
	}
	lhs := kids[0]
	if !isAlwaysFalsyLiteral(lhs) {
		return
	}
	op := "&&"
	if n.Token() == rhino.OR {
		op = "||"
	}
	c.diags = append(c.diags, diag.Diagnostic{
		Path:    c.path,
		Line:    lhs.Lineno(),
		Col:     lhs.Charno(),
		Code:    diag.JSC_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR,
		Level:   diag.Warning,
		Message: "Suspicious code. The result of the '" + op + "' operator is constant because the left operand is always falsy.",
	})
}

// isAlwaysFalsyLiteral reports whether `n` is one of the literal
// falsy values `null`, `undefined`, `0`, `NaN`, `false`, or `""`.
// NaN is represented as a bare NAME with value "NaN" in Closure's IR.
func isAlwaysFalsyLiteral(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	switch n.Token() {
	case rhino.NULL, rhino.UNDEFINED, rhino.FALSE:
		return true
	case rhino.NUMBER:
		return n.Number() == 0
	case rhino.STRINGLIT:
		return n.String_() == ""
	case rhino.NAME:
		// NaN is the only always-falsy NAME we recognise. `undefined`
		// could also appear as a NAME in Java's older IR, but our
		// parser surfaces it as an UNDEFINED token — handled above.
		return n.String_() == "NaN"
	}
	return false
}

func (c *CheckSuspiciousCode) emitSuspiciousExpr(n *rhino.Node) {
	c.diags = append(c.diags, diag.Diagnostic{
		Path:    c.path,
		Line:    n.Lineno(),
		Col:     n.Charno(),
		Code:    diag.JSC_SUSPICIOUS_EXPR,
		Level:   diag.Warning,
		Message: "Suspicious code. Assignment used as a conditional; did you mean '=='?",
	})
}
