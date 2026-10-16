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
	"github.com/paullepoulpe/jsgo/legacy/internal/diag"
	"github.com/paullepoulpe/jsgo/legacy/internal/rhino"
)

// CheckSideEffects implements the JSC_USELESS_CODE diagnostic.
type CheckSideEffects struct {
	path            string
	diags           []diag.Diagnostic
	noSEFns         map[string]bool // qualified names of @nosideeffects functions
	undeclaredNames map[string]bool // VarCheck-reported undeclared root names
	analyzer        *AstAnalyzer
}

// Go-only: language-adapter — constructor (Java uses new CheckSideEffects(compiler, report, protectSideEffectFreeCode))
func NewCheckSideEffects(path string) *CheckSideEffects {
	return &CheckSideEffects{path: path, analyzer: NewAstAnalyzer()}
}

// SetUndeclaredNames mirrors TypeCheck.SetUndeclaredNames: when VarCheck has
// reported JSC_UNDEFINED_VARIABLE for a NAME, downstream passes treat the
// root name as "already flagged" and suppress their own diagnostics on
// chains rooted in that name. In Java this is implicit because the
// diagnostic group ordering / pass halts prevent the duplicate warning;
// the Go pipeline runs all passes, so we replicate the suppression here.
func (c *CheckSideEffects) SetUndeclaredNames(names map[string]bool) {
	c.undeclaredNames = names
}

// process mirrors Java's CompilerPass.process(Node externs, Node root):
//  1. GetNoSideEffectExterns sub-pass collects @nosideeffects names.
//  2. Main traversal fires diagnostics.
//  3. Optionally protect side-effect free code for dead-code elimination.
func (c *CheckSideEffects) process(script *rhino.Node) {
	c.noSEFns = make(map[string]bool)
	c.collectNoSideEffectFns(script)
	c.visit(script)
	// protectSideEffects is skipped in this port (requires full compiler infrastructure).
}

// Check walks the AST firing JSC_USELESS_CODE at every EXPR_RESULT
// whose contained expression is observably pure.
// Go-only: language-adapter — Go entry point that wraps process() and returns diagnostics.
func (c *CheckSideEffects) Check(script *rhino.Node) []diag.Diagnostic {
	c.diags = nil
	c.process(script)
	return c.diags
}

// protectSideEffects mirrors Java's CheckSideEffects.protectSideEffects.
// Wraps side-effect-free expressions in a JSCOMPILER_PRESERVE call so they
// survive dead-code elimination. Skipped in this port — requires full compiler
// node manipulation infrastructure.
func (c *CheckSideEffects) protectSideEffects() {
	// Not yet ported: requires compiler.reportChangeToEnclosingScope and IR.call infrastructure.
	// Java: problemNodes are wrapped in JSCOMPILER_PRESERVE(expr) calls.
	// protectSideEffectFreeCode flag must also be set in the Java constructor.
	_ = c.noSEFns // suppress unused-in-stub warning
}

// addExtern mirrors Java's CheckSideEffects.addExtern (static method).
// Injects the JSCOMPILER_PRESERVE extern into the synthetic externs so
// optimizer passes know about it. Skipped in this port — requires a full
// AbstractCompiler instance.
func addExtern(c *CheckSideEffects) {
	// Not yet ported: requires NodeUtil.createSynthesizedExternsSymbol(compiler, PROTECTOR_FN).
	_ = c
}

// Go-only: language-adapter — merged scan (Java uses separate GetNoSideEffectExterns traversal pass)
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

// visit mirrors Java's AbstractPostOrderCallback.visit — called for each node in post-order.
// Go uses pre-order walk since the port applies EXPR_RESULT detection directly.
func (c *CheckSideEffects) visit(n *rhino.Node) {
	if n == nil {
		return
	}
	if n.Token() == rhino.EXPR_RESULT {
		if !isDeclarationStatement(n) {
			if inner := n.FirstChild(); inner != nil {
				if c.rootedInUndeclared(inner) {
					// Suppress USELESS_CODE when the chain is rooted in a
					// name VarCheck already flagged as JSC_UNDEFINED_VARIABLE.
					// Java's downstream passes don't re-fire on these because
					// the compilation halts; we replicate by short-circuit.
				} else {
					c.reportUseless(n, inner, true /*top-level*/)
				}
			}
		}
	}
	// Mirror Java's CheckSideEffects.visit for IF statements whose entire
	// body / else / condition subtree is side-effect-free. Java's logic
	// reaches every node and fires when isExpressionResultUsed=false and
	// !mayHaveSideEffects(n). For an IF inside a BLOCK, the result is
	// never "used" (NodeUtil.isExpressionResultUsed returns false on BLOCK
	// parent), so a fully-pure IF (e.g. `if (x.foo) {}`) becomes
	// USELESS_CODE at the IF's own position. Top-level IFs (parent SCRIPT)
	// fall into NodeUtil.isExpressionResultUsed's default branch which
	// returns true, so Java does NOT fire — we restrict the same way.
	if n.Token() == rhino.IF {
		parent := n.Parent()
		if parent != nil && parent.Token() == rhino.BLOCK {
			cond := n.FirstChild()
			if cond != nil && !c.rootedInUndeclared(cond) && !c.analyzer.MayHaveSideEffects(n) {
				c.emitAt(n)
			}
		}
	}
	for _, child := range n.Children() {
		c.visit(child)
	}
}

// rootedInUndeclared reports whether the qualified-name chain (NAME or
// GETPROP chain) at n is rooted at a NAME that VarCheck flagged as
// undeclared. Mirrors Java's rootNameNode check in CheckSideEffects
// when deciding whether to emit noSideEffectExterns-related warnings.
func (c *CheckSideEffects) rootedInUndeclared(n *rhino.Node) bool {
	if c.undeclaredNames == nil {
		return false
	}
	for n != nil {
		switch n.Token() {
		case rhino.NAME:
			return c.undeclaredNames[n.String_()]
		case rhino.GETPROP, rhino.GETELEM, rhino.OPTCHAIN_GETPROP, rhino.OPTCHAIN_GETELEM:
			n = n.FirstChild()
			continue
		default:
			return false
		}
	}
	return false
}

// Go-only: language-adapter — helper for visit() that handles both full-statement and sub-expression reporting
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
		case rhino.EQ:
			c.emitOpAt(stmt, "eq")
			return
		case rhino.NE:
			c.emitOpAt(stmt, "ne")
			return
		case rhino.SHEQ:
			c.emitOpAt(stmt, "sheq")
			return
		case rhino.SHNE:
			c.emitOpAt(stmt, "shne")
			return
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

// Go-only: language-adapter — helper to emit diagnostic at a node position
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

// Go-only: language-adapter — operator-specific variant of emitAt (Java builds message inline in visit())
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

// Go-only: language-adapter — checks for Closure slot declarations (Java checks isQualifiedName() + JSDoc inline in visit())
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

// Go-only: language-adapter — combined side-effect checker (Java uses AstAnalyzer.mayHaveSideEffects() + noSideEffectExterns set)
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
	// Logical-assignment operators (&&=, ||=, ??=) also mutate a
	// binding — they are as effectful as plain ASSIGN.
	case rhino.ASSIGN, rhino.ASSIGN_ADD, rhino.ASSIGN_SUB, rhino.ASSIGN_MUL,
		rhino.ASSIGN_DIV, rhino.ASSIGN_MOD, rhino.ASSIGN_BITAND,
		rhino.ASSIGN_BITOR, rhino.ASSIGN_BITXOR, rhino.ASSIGN_LSH,
		rhino.ASSIGN_RSH, rhino.ASSIGN_URSH, rhino.ASSIGN_EXPONENT,
		rhino.ASSIGN_AND, rhino.ASSIGN_OR, rhino.ASSIGN_COALESCE,
		rhino.INC, rhino.DEC,
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
