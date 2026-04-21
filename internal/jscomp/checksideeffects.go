// Port of com/google/javascript/jscomp/CheckSideEffects.java @ closure-compiler v20260415
//
// Mirrors the Java post-order visitor: at every AST node whose value is
// not consumed by its enclosing context, the pass fires JSC_USELESS_CODE
// (via NodeUtil.isExpressionResultUsed). Two message variants exist:
//
//   - "Suspicious code. The result of the '<op>' operator is not being
//     used." — emitted when the discarded node is a "simple operator"
//     (Java's NodeUtil.isSimpleOperator union, plus GETPROP / GETELEM).
//     The operator name is the lower-cased Java Token enum name (e.g.
//     EQ → "eq", BITNOT → "bitnot", GETPROP → "getprop").
//   - "Suspicious code. This code lacks side-effects. Is there a bug?"
//     — emitted when the discarded node is not a simple operator and
//     has no side effects.
//
// The CALL branch additionally consults @nosideeffects: a CALL whose
// named callee was annotated `@nosideeffects` (or `@noside`) produces
// the operator-specific "extern function call" warning when its result
// is discarded — Java's GetNoSideEffectExterns visitor only inspects
// the externs tree; jsgo's checker doesn't have a separate externs AST
// at this layer, so we honour `@nosideeffects` on any function
// declaration in scope, matching the spirit of the Java pass even
// though the source of the JSDoc is broader.
//
// Skipped branches that don't apply at the checker layer:
//
//   - protectSideEffects (rewrites pure expressions into JSCOMPILER_PRESERVE
//     calls so optimizers don't drop them) — purely a transformation.
//   - StripProtection (post-optimization undo) — symmetric to the above.
//   - The runtime-library `"require <name>"` exemption — jsgo doesn't
//     bundle Closure's builtin runtime files, so the heuristic cannot
//     trigger.
//   - The COMMA / CAST / EMPTY tokens are not produced by the current
//     parser; their handling will land when the parser grows them.

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CheckSideEffects implements the JSC_USELESS_CODE diagnostic.
type CheckSideEffects struct {
	path  string
	diags []diag.Diagnostic
	// noSideEffectFns maps a function-declaration NAME to true when
	// that function's JSDocInfo carried `@nosideeffects` (or its alias
	// `@noside`). Populated in a pre-walk so call-site classification
	// can look up the marker on a simple-name callee without a scope
	// handle. Mirrors the role of CheckSideEffects.noSideEffectExterns
	// in Java (whose set is sourced exclusively from the externs tree).
	noSideEffectFns map[string]bool
}

func NewCheckSideEffects(path string) *CheckSideEffects {
	return &CheckSideEffects{path: path}
}

// Check walks the AST post-order, firing JSC_USELESS_CODE wherever a
// node's value is structurally discarded. Mirrors Java's
// CheckSideEffects.process: scan the externs (here, the script) for
// `@nosideeffects` markers, then traverse the source tree.
func (c *CheckSideEffects) Check(script *rhino.Node) []diag.Diagnostic {
	c.diags = nil
	c.noSideEffectFns = map[string]bool{}
	c.collectNoSideEffectFns(script)
	c.walk(script, nil, nil)
	return c.diags
}

// collectNoSideEffectFns pre-walks the script and records every named
// function declaration whose JSDocInfo carries `@nosideeffects`. Mirrors
// Java's CheckSideEffects.GetNoSideEffectExterns inner visitor (which
// only walks externs); jsgo runs the same scan over the script because
// the externs and source share a single AST at this layer.
func (c *CheckSideEffects) collectNoSideEffectFns(n *rhino.Node) {
	if n == nil {
		return
	}
	if n.Token() == rhino.FUNCTION {
		if info := n.JSDocInfo(); info != nil && info.IsNoSideEffects() {
			if name := n.FirstChild(); name != nil && name.String_() != "" {
				c.noSideEffectFns[name.String_()] = true
			}
		}
	}
	for _, child := range n.Children() {
		c.collectNoSideEffectFns(child)
	}
}

// walk runs the Java post-order traversal. The Java visitor is invoked
// after children are visited; we mirror that by descending first, then
// calling visit. Parent and grandparent are threaded explicitly because
// rhino.Node's parent pointer isn't always set when we descend into
// freshly-built subtrees in tests, and the class-member guard needs the
// grandparent to distinguish a FUNCTION in a class body from a FUNCTION
// in a block.
func (c *CheckSideEffects) walk(n, parent, grandparent *rhino.Node) {
	if n == nil {
		return
	}
	for _, child := range n.Children() {
		c.walk(child, n, parent)
	}
	c.visit(n, parent, grandparent)
}

// visit mirrors CheckSideEffects.visit(NodeTraversal, Node, Node).
// Returns silently for nodes whose semantics rule them out (EXPR_RESULT
// / BLOCK shells, qualified-name declarations, root nodes) and for
// nodes whose value is consumed.
func (c *CheckSideEffects) visit(n, parent, grandparent *rhino.Node) {
	if parent == nil {
		return
	}
	switch n.Token() {
	case rhino.SCRIPT, rhino.EXPR_RESULT, rhino.BLOCK:
		// Java skips EMPTY / COMMA at the top, then EXPR_RESULT / BLOCK
		// shells. SCRIPT isn't in Java's switch but its children carry
		// EXPR_RESULTs which we never want to flag at the script level.
		return
	}
	// Class-body methods (FUNCTION nodes whose parent is the class BLOCK
	// and whose grandparent is a CLASS) are declarations, not expression
	// statements. Java's AST wraps them in MEMBER_FUNCTION_DEF so the
	// check never fires; jsgo's parser lowers them as bare FUNCTIONs
	// under the class body BLOCK, so we need a structural guard here.
	if n.Token() == rhino.FUNCTION && parent.Token() == rhino.BLOCK &&
		grandparent != nil && grandparent.Token() == rhino.CLASS {
		return
	}
	// `Foo.prototype.bar;` (or any qualified-name read) carrying its
	// own JSDoc is a slot-declaration statement; Java skips it via
	// `n.isQualifiedName() && n.getJSDocInfo() != null`. Apply the same
	// guard at the GETPROP/NAME level. We also keep the pre-existing
	// statement-level helper for the `obj.prop;` case where the JSDoc
	// hangs off the EXPR_RESULT shell instead of the inner expression.
	if isQualifiedName(n) && n.JSDocInfo() != nil {
		return
	}
	if parent.Token() == rhino.EXPR_RESULT && isDeclarationStatement(parent) {
		return
	}
	if isExpressionResultUsed(n, parent) {
		return
	}
	if n.Token() == rhino.VOID {
		// Java explicitly tolerates `void <expr>;` — extra semicolons
		// at the BLOCK level emit VOID nodes that nobody considers
		// suspicious.
		return
	}
	simple := isSimpleOperator(n.Token())
	if simple || !mayHaveSideEffects(n) {
		c.report(n, simple)
		return
	}
	if n.Token() == rhino.CALL {
		c.maybeReportNoSideEffectCall(n)
	}
}

// report emits the JSC_USELESS_CODE diagnostic at `n`'s position. When
// `simple` is true the operator-specific message variant is used; when
// `n` is a string literal the "missing '+'" variant fires (mirrors
// Java's `n.isStringLit() || n.isTemplateLit()` branch). Otherwise the
// generic "lacks side-effects" variant is emitted.
//
// GETPROP gets a position fix-up: Java's IRFactory anchors a GETPROP
// node at the property-name token, while jsgo's parser anchors it at
// the receiver. We compensate so Java's diagnostic column is preserved.
func (c *CheckSideEffects) report(n *rhino.Node, simple bool) {
	msg := "Suspicious code. This code lacks side-effects. Is there a bug?"
	switch {
	case n.Token() == rhino.STRINGLIT:
		msg = "Suspicious code. Is there a missing '+' on the previous line?"
	case simple:
		msg = "Suspicious code. The result of the '" + operatorName(n.Token()) + "' operator is not being used."
	}
	line, col := n.Lineno(), n.Charno()
	if n.Token() == rhino.GETPROP {
		if kids := n.Children(); len(kids) >= 2 {
			line, col = kids[1].Lineno(), kids[1].Charno()
		}
	}
	c.diags = append(c.diags, diag.Diagnostic{
		Path:    c.path,
		Line:    line,
		Col:     col,
		Code:    diag.JSC_USELESS_CODE,
		Level:   diag.Warning,
		Message: msg,
	})
}

// maybeReportNoSideEffectCall mirrors the Java fallback branch: when a
// CALL has side effects (because CALL is treated as effectful by
// default) but the callee resolves to a `@nosideeffects` function, the
// pass still emits a useless-code warning — the value is discarded and
// the @nosideeffects pledge tells us the call was pointless.
func (c *CheckSideEffects) maybeReportNoSideEffectCall(call *rhino.Node) {
	kids := call.Children()
	if len(kids) == 0 {
		return
	}
	callee := kids[0]
	var name string
	switch callee.Token() {
	case rhino.NAME:
		name = callee.String_()
	case rhino.GETPROP:
		// Java's CheckSideEffects accepts GETPROP / NAME / STRINGLIT
		// callees and resolves the qualified name. We currently track
		// only top-level NAME bindings, so a `a.b()` callee never
		// matches and the warning isn't emitted.
		return
	default:
		return
	}
	if !c.noSideEffectFns[name] {
		return
	}
	c.diags = append(c.diags, diag.Diagnostic{
		Path:    c.path,
		Line:    call.Lineno(),
		Col:     call.Charno(),
		Code:    diag.JSC_USELESS_CODE,
		Level:   diag.Warning,
		Message: "Suspicious code. The result of the extern function call '" + name + "' is not being used.",
	})
}

// isQualifiedName mirrors Java's Node.isQualifiedName: a NAME, THIS,
// SUPER, or a GETPROP whose root walks down to one of those. Used by
// the JSDoc-on-qname guard so prototype-slot declarations like
// `/** @type {number} */ Foo.prototype.x;` aren't flagged as useless.
func isQualifiedName(n *rhino.Node) bool {
	switch n.Token() {
	case rhino.NAME, rhino.THIS:
		return true
	case rhino.GETPROP:
		kids := n.Children()
		if len(kids) == 0 {
			return false
		}
		return isQualifiedName(kids[0])
	}
	return false
}

// isExpressionResultUsed ports NodeUtil.isExpressionResultUsed. Returns
// true when `expr`'s value is consumed by the enclosing context. The
// switch covers the contexts where the answer is structurally false
// (EXPR_RESULT / BLOCK children, COMMA's left operand, FOR's init /
// update); the default case is true (a function arg, an operator
// operand, an assignment RHS, etc.).
func isExpressionResultUsed(expr, parent *rhino.Node) bool {
	if parent == nil {
		return false
	}
	switch parent.Token() {
	case rhino.BLOCK, rhino.EXPR_RESULT:
		return false
	case rhino.HOOK, rhino.AND, rhino.OR, rhino.COALESCE:
		// First child of HOOK / AND / OR / COALESCE is always
		// evaluated for control-flow; later children's result-used
		// status mirrors the parent's.
		if kids := parent.Children(); len(kids) > 0 && kids[0] == expr {
			return true
		}
		return isExpressionResultUsed(parent, parent.Parent())
	case rhino.FOR:
		// In Closure's FOR(init, cond, upd, body), only the condition
		// (second child) has its value consumed. Init / update / body
		// run for side effect.
		kids := parent.Children()
		return len(kids) > 1 && kids[1] == expr
	}
	return true
}

// isSimpleOperator ports NodeUtil.isSimpleOperatorType — the union of
// pure binary / unary operators whose result-discard yields the
// operator-specific JSC_USELESS_CODE message variant. GETPROP and
// GETELEM are added by CheckSideEffects on top of this list.
func isSimpleOperator(t rhino.Token) bool {
	switch t {
	case rhino.ADD, rhino.SUB, rhino.MUL, rhino.DIV, rhino.MOD,
		rhino.BITAND, rhino.BITOR, rhino.BITXOR, rhino.BITNOT,
		rhino.LSH, rhino.RSH, rhino.URSH,
		rhino.EQ, rhino.NE, rhino.SHEQ, rhino.SHNE,
		rhino.LT, rhino.LE, rhino.GT, rhino.GE,
		rhino.IN, rhino.INSTANCEOF,
		rhino.NOT, rhino.POS, rhino.NEG, rhino.TYPEOF, rhino.VOID,
		rhino.GETPROP, rhino.GETELEM:
		return true
	}
	return false
}

// operatorName returns the lowercase Java Token enum name used in the
// "result of the '<op>' operator" message. Mirrors Java's
// `Ascii.toLowerCase(n.getToken().toString())`.
func operatorName(t rhino.Token) string {
	switch t {
	case rhino.ADD:
		return "add"
	case rhino.SUB:
		return "sub"
	case rhino.MUL:
		return "mul"
	case rhino.DIV:
		return "div"
	case rhino.MOD:
		return "mod"
	case rhino.BITAND:
		return "bitand"
	case rhino.BITOR:
		return "bitor"
	case rhino.BITXOR:
		return "bitxor"
	case rhino.BITNOT:
		return "bitnot"
	case rhino.LSH:
		return "lsh"
	case rhino.RSH:
		return "rsh"
	case rhino.URSH:
		return "ursh"
	case rhino.EQ:
		return "eq"
	case rhino.NE:
		return "ne"
	case rhino.SHEQ:
		return "sheq"
	case rhino.SHNE:
		return "shne"
	case rhino.LT:
		return "lt"
	case rhino.LE:
		return "le"
	case rhino.GT:
		return "gt"
	case rhino.GE:
		return "ge"
	case rhino.IN:
		return "in"
	case rhino.INSTANCEOF:
		return "instanceof"
	case rhino.NOT:
		return "not"
	case rhino.POS:
		return "pos"
	case rhino.NEG:
		return "neg"
	case rhino.TYPEOF:
		return "typeof"
	case rhino.VOID:
		return "void"
	case rhino.GETPROP:
		return "getprop"
	case rhino.GETELEM:
		return "getelem"
	}
	return ""
}

// isDeclarationStatement reports whether `n` (an EXPR_RESULT) is
// actually a prototype-slot / property-type declaration rather than
// a real expression statement. Closure treats `Foo.prototype.bar;`
// (with an @type JSDoc) as a declaration of a slot — Java's
// CheckSideEffects skips these via the n.isQualifiedName + JSDocInfo
// check on the inner node. We honour the same shape on the EXPR_RESULT
// shell because jsgo's parser hangs the JSDoc on the statement, not
// the inner GETPROP.
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

// mayHaveSideEffects ports the subset of AstAnalyzer.checkForStateChangeHelper
// the checker layer needs. The Java method is large (it special-cases
// every token); we mirror the dispatch shape so the boundary between
// "always effectful" / "transparent compound" / "leaf" matches Java's
// classification:
//
//   - VAR / LET / CONST / THROW / FOR_IN / FOR_OF and the assignment
//     family are always effectful.
//   - CALL / NEW are always effectful here; the @nosideeffects fallback
//     in visit() handles the discardable case after-the-fact.
//   - HOOK / AND / OR / COALESCE / IF / BLOCK / EXPR_RESULT / SWITCH /
//     TRY / GETPROP / GETELEM and the simple operators recurse into
//     children — they're transparent.
//   - Pure leaves (NAME / NUMBER / STRINGLIT / TRUE / FALSE / NULL /
//     UNDEFINED / THIS / ARRAYLIT / OBJECTLIT / FUNCTION-expression)
//     are side-effect-free.
//   - Everything else (statements like WHILE / DO_WHILE / RETURN /
//     BREAK / CONTINUE / LABEL / etc.) defaults to effectful, matching
//     the `default → true` arm of Java's switch.
func mayHaveSideEffects(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	switch n.Token() {
	// Definitely-effectful tokens.
	case rhino.CALL, rhino.NEW, rhino.DELPROP,
		rhino.ASSIGN, rhino.ASSIGN_ADD, rhino.ASSIGN_SUB, rhino.ASSIGN_MUL,
		rhino.ASSIGN_DIV, rhino.ASSIGN_MOD, rhino.ASSIGN_BITAND,
		rhino.ASSIGN_BITOR, rhino.ASSIGN_BITXOR, rhino.ASSIGN_LSH,
		rhino.ASSIGN_RSH, rhino.ASSIGN_URSH,
		rhino.ASSIGN_AND, rhino.ASSIGN_OR, rhino.ASSIGN_COALESCE,
		rhino.THROW, rhino.VAR, rhino.LET, rhino.CONST,
		rhino.FOR_IN, rhino.FOR_OF:
		return true
	// Pure leaves.
	case rhino.NAME, rhino.STRINGLIT, rhino.NUMBER, rhino.TRUE, rhino.FALSE,
		rhino.NULL, rhino.UNDEFINED, rhino.THIS,
		rhino.ARRAYLIT, rhino.OBJECTLIT, rhino.FUNCTION:
		return false
	// Transparent compound expressions: side-effect iff any child is.
	// Includes the simple operators (which by definition never
	// side-effect on their own), plus HOOK / short-circuit / GETPROP /
	// GETELEM, plus the structural shells BLOCK / EXPR_RESULT / SCRIPT
	// the visit walks past on its own.
	case rhino.HOOK, rhino.AND, rhino.OR, rhino.COALESCE,
		rhino.GETPROP, rhino.GETELEM,
		rhino.BLOCK, rhino.EXPR_RESULT, rhino.SCRIPT,
		rhino.ADD, rhino.SUB, rhino.MUL, rhino.DIV, rhino.MOD,
		rhino.BITAND, rhino.BITOR, rhino.BITXOR, rhino.BITNOT,
		rhino.LSH, rhino.RSH, rhino.URSH,
		rhino.EQ, rhino.NE, rhino.SHEQ, rhino.SHNE,
		rhino.LT, rhino.LE, rhino.GT, rhino.GE,
		rhino.IN, rhino.INSTANCEOF,
		rhino.NOT, rhino.POS, rhino.NEG, rhino.TYPEOF, rhino.VOID,
		rhino.PARAM_LIST, rhino.STRING_KEY:
		for _, c := range n.Children() {
			if mayHaveSideEffects(c) {
				return true
			}
		}
		return false
	}
	// Default: statements (WHILE / DO_WHILE / IF / SWITCH / TRY / FOR /
	// RETURN / BREAK / etc.) and any token we haven't classified are
	// treated as effectful — Java's checkForStateChangeHelper falls
	// through to `yield true` here.
	return true
}
