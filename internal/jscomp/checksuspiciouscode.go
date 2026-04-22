// Port of com/google/javascript/jscomp/CheckSuspiciousCode.java @ closure-compiler v20260415
// (narrow; no JSType on nodes so type-based boolean-value branch is stubbed;
// no OPTCHAIN / FOR_AWAIT_OF tokens so those sub-checks are no-ops).

package jscomp

import (
	"fmt"
	"math"

	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Tri mirrors com/google/javascript/jscomp/base/Tri.java — a triboolean
// used by getBooleanValueWithTypes. Only the subset needed here is
// included. Java's Tri is an enum with TRUE / FALSE / UNKNOWN values.
type Tri int

const (
	// TriUnknown mirrors Tri.UNKNOWN.
	TriUnknown Tri = iota
	// TriTrue mirrors Tri.TRUE.
	TriTrue
	// TriFalse mirrors Tri.FALSE.
	TriFalse
)

// not mirrors Tri.not() in Java: TRUE→FALSE, FALSE→TRUE, UNKNOWN→UNKNOWN.
func (t Tri) not() Tri {
	switch t {
	case TriTrue:
		return TriFalse
	case TriFalse:
		return TriTrue
	}
	return TriUnknown
}

// and mirrors Tri.and() in Java. Short-circuit: FALSE and _ = FALSE;
// TRUE and x = x; UNKNOWN and TRUE = UNKNOWN; UNKNOWN and _ = UNKNOWN.
func (t Tri) and(other Tri) Tri {
	if t == TriFalse || other == TriFalse {
		return TriFalse
	}
	if t == TriTrue {
		return other
	}
	return TriUnknown
}

// or mirrors Tri.or() in Java. Short-circuit: TRUE or _ = TRUE;
// FALSE or x = x; UNKNOWN or FALSE = UNKNOWN; UNKNOWN or _ = UNKNOWN.
func (t Tri) or(other Tri) Tri {
	if t == TriTrue || other == TriTrue {
		return TriTrue
	}
	if t == TriFalse {
		return other
	}
	return TriUnknown
}

// CheckSuspiciousCode checks for common errors, such as misplaced semicolons:
//
//	if (x); act_now();
//
// or comparison against NaN:
//
//	if (x === NaN) act();
//
// and generates warnings. Port of Java's CheckSuspiciousCode which
// extends AbstractPostOrderCallback.
type CheckSuspiciousCode struct {
	path  string
	diags []diag.Diagnostic
}

// NewCheckSuspiciousCode builds a new pass for the given source file.
func NewCheckSuspiciousCode(path string) *CheckSuspiciousCode {
	return &CheckSuspiciousCode{path: path}
}

// Check runs the pass over the given script and returns accumulated
// diagnostics. Drives a simple post-order walk (mirrors Java's
// AbstractPostOrderCallback).
func (c *CheckSuspiciousCode) Check(script *rhino.Node) []diag.Diagnostic {
	c.diags = nil
	c.walk(script)
	return c.diags
}

func (c *CheckSuspiciousCode) walk(n *rhino.Node) {
	if n == nil {
		return
	}
	for _, child := range n.Children() {
		c.walk(child)
	}
	c.visit(n)
}

// visit mirrors the Java visit(NodeTraversal t, Node n, Node parent)
// method, dispatching to each sub-check in the same order as Java.
func (c *CheckSuspiciousCode) visit(n *rhino.Node) {
	c.checkMissingSemicolon(n)
	c.checkNaN(n)
	c.checkInvalidIn(n)
	c.checkNonObjectInstanceOf(n)
	c.checkNegatedLeftOperandOfInOperator(n)
	c.checkLeftOperandOfLogicalOperator(n)
	// checkSuspiciousBreakingOutOfOptionalChain: OPTCHAIN tokens not in
	// our token set (optional chaining not yet parsed), so this is a
	// no-op. Java line 182.
}

// checkMissingSemicolon fires JSC_SUSPICIOUS_SEMICOLON when a control
// statement has a synthesized empty body (i.e. `if (x);`). Mirrors
// Java's checkMissingSemicolon.
func (c *CheckSuspiciousCode) checkMissingSemicolon(n *rhino.Node) {
	switch n.Token() {
	case rhino.IF:
		trueCase := n.FirstChild()
		if trueCase == nil {
			return
		}
		trueCase = trueCase.Next() // second child is the thenBlock
		if trueCase == nil {
			return
		}
		c.reportIfWasEmpty(trueCase)
		elseCase := trueCase.Next()
		if elseCase != nil {
			c.reportIfWasEmpty(elseCase)
		}
	case rhino.WHILE, rhino.FOR, rhino.FOR_IN, rhino.FOR_OF:
		body := loopCodeBlock(n)
		if body != nil {
			c.reportIfWasEmpty(body)
		}
	}
}

// reportIfWasEmpty fires JSC_SUSPICIOUS_SEMICOLON on a BLOCK node that
// was synthesized (added block) and has no children. Mirrors Java's
// reportIfWasEmpty(NodeTraversal t, Node block).
func (c *CheckSuspiciousCode) reportIfWasEmpty(block *rhino.Node) {
	if block.Token() != rhino.BLOCK {
		return
	}
	// A semicolon is distinguished from a block without children by
	// annotating it with PropAddedBlock. Blocks without children are
	// usually intentional, especially with loops.
	if !block.HasChildren() && block.IsAddedBlock() {
		c.emit(block, diag.JSC_SUSPICIOUS_SEMICOLON, diag.Warning,
			"If this if/for/while really shouldn't have a body, use '{}'")
	}
}

// checkNaN fires JSC_SUSPICIOUS_NAN when a comparison has NaN as an
// operand. Mirrors Java's checkNaN.
func (c *CheckSuspiciousCode) checkNaN(n *rhino.Node) {
	switch n.Token() {
	case rhino.EQ, rhino.GE, rhino.GT, rhino.LE, rhino.LT, rhino.NE, rhino.SHEQ, rhino.SHNE:
		c.reportIfNaN(n.FirstChild())
		c.reportIfNaN(lastChild(n))
	}
}

// reportIfNaN fires JSC_SUSPICIOUS_NAN on the parent node when n is NaN.
// Mirrors Java's reportIfNaN(NodeTraversal t, Node n).
func (c *CheckSuspiciousCode) reportIfNaN(n *rhino.Node) {
	if isNaN(n) {
		parent := n.Parent()
		if parent == nil {
			return
		}
		c.emit(parent, diag.JSC_SUSPICIOUS_NAN, diag.Warning,
			"Comparison against NaN is always false. Did you mean isNaN()?")
	}
}

// checkInvalidIn fires JSC_SUSPICIOUS_IN when the `in` operator is used
// with a non-object RHS. Mirrors Java's checkInvalidIn.
func (c *CheckSuspiciousCode) checkInvalidIn(n *rhino.Node) {
	if n.Token() == rhino.IN {
		c.reportIfNonObject(lastChild(n), diag.JSC_SUSPICIOUS_IN, diag.Warning,
			"Use of the \"in\" keyword on non-object types throws an exception.")
	}
}

// checkNonObjectInstanceOf fires JSC_SUSPICIOUS_INSTANCEOF_LEFT when the
// left operand of `instanceof` is not an object. Mirrors Java's
// checkNonObjectInstanceOf.
func (c *CheckSuspiciousCode) checkNonObjectInstanceOf(n *rhino.Node) {
	if n.Token() == rhino.INSTANCEOF {
		c.reportIfNonObject(n.FirstChild(), diag.JSC_SUSPICIOUS_INSTANCEOF_LEFT, diag.Warning,
			"\"instanceof\" with left non-object operand is always false.")
	}
}

// reportIfNonObject fires the given diagnostic on n.Parent() when n
// cannot be an object. Mirrors Java's reportIfNonObject.
func (c *CheckSuspiciousCode) reportIfNonObject(n *rhino.Node, code string, level diag.Level, msg string) bool {
	if n == nil {
		return false
	}
	if n.Token() == rhino.ADD || !mayBeObject(n) {
		parent := n.Parent()
		if parent == nil {
			return false
		}
		c.emit(parent, code, level, msg)
		return true
	}
	return false
}

// checkNegatedLeftOperandOfInOperator fires
// JSC_SUSPICIOUS_NEGATED_LEFT_OPERAND_OF_IN_OPERATOR when the `in`
// operator's LHS is a NOT expression. Mirrors Java's
// checkNegatedLeftOperandOfInOperator.
func (c *CheckSuspiciousCode) checkNegatedLeftOperandOfInOperator(n *rhino.Node) {
	if n.Token() == rhino.IN {
		lhs := n.FirstChild()
		if lhs != nil && lhs.Token() == rhino.NOT {
			c.emit(lhs, diag.JSC_SUSPICIOUS_NEGATED_LEFT_OPERAND_OF_IN_OPERATOR, diag.Warning,
				"Suspicious negated left operand of 'in' operator.")
		}
	}
}

// checkLeftOperandOfLogicalOperator fires
// JSC_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR when the LHS of `&&`
// or `||` is deterministically truthy or falsy. Mirrors Java's
// checkLeftOperandOfLogicalOperator.
//
// Check for the LHS of a logical operator (&& and ||) being
// deterministically truthy or falsy, using both syntactic and type
// information. This is always suspicious (though for different reasons:
// "truthy and" means the LHS is always ignored and should be removed,
// "falsy and" means the RHS is dead code, and vice versa for "or".
// However, there are a number of legitimate use cases where we need to
// back off from using type information (see getBooleanValueWithTypes for
// more details on these back offs).
func (c *CheckSuspiciousCode) checkLeftOperandOfLogicalOperator(n *rhino.Node) {
	if n.Token() == rhino.OR || n.Token() == rhino.AND {
		var operator string
		if n.Token() == rhino.OR {
			operator = "||"
		} else {
			operator = "&&"
		}
		v := c.getBooleanValueWithTypes(n.FirstChild())
		if v != TriUnknown {
			var result string
			if v == TriTrue {
				result = "truthy"
			} else {
				result = "falsy"
			}
			c.emit(n, diag.JSC_SUSPICIOUS_LEFT_OPERAND_OF_LOGICAL_OPERATOR, diag.Warning,
				fmt.Sprintf("Left operand of %s operator is always %s.", operator, result))
		}
	}
}

// getBooleanValueWithTypes returns the possible boolean value of a node,
// combining syntactic checks (NodeUtil.getBooleanValue) with type
// information where available. Mirrors Java's getBooleanValueWithTypes.
//
// Returns the possible boolean values of a node. This is a combination of
// NodeUtil#getBooleanValue and JSType#getPossibleToBooleanOutcomes, with
// some additional backoff for situations that are known to be less
// accurate (to wit, qualified names and truthy return types, which return
// UNKNOWN even if the type information seems conclusive). Another
// difference from the NodeUtil and JSType methods is that we include some
// specific optimization to avoid quadratic behavior of large nested
// logical operations.
//
// Since nodes in this port do not carry a JSType annotation, the type-
// based branch always returns UNKNOWN (same as if type is nil in Java).
// TODO(java:277): attach JSType to nodes and re-enable type-based branch.
func (c *CheckSuspiciousCode) getBooleanValueWithTypes(n *rhino.Node) Tri {
	if n == nil {
		return TriUnknown
	}
	switch n.Token() {
	case rhino.ASSIGN:
		// Java also handles COMMA (comma operator) here; that token is not
		// in our subset (parser does not emit it). Java line 219.
		return c.getBooleanValueWithTypes(lastChild(n))
	case rhino.NOT:
		return c.getBooleanValueWithTypes(lastChild(n)).not()
	case rhino.AND:
		// Assume the left-hand side is unknown. If it's not then we'll
		// report it elsewhere. This prevents revisiting deeper nodes
		// repeatedly, which would result in O(n^2) performance.
		return TriUnknown.and(c.getBooleanValueWithTypes(lastChild(n)))
	case rhino.OR:
		// Assume the left-hand side is unknown. If it's not then we'll
		// report it elsewhere. This prevents revisiting deeper nodes
		// repeatedly, which would result in O(n^2) performance.
		return TriUnknown.or(c.getBooleanValueWithTypes(lastChild(n)))
	case rhino.HOOK:
		trueValue := c.getBooleanValueWithTypes(secondChild(n))
		falseValue := c.getBooleanValueWithTypes(lastChild(n))
		if trueValue == falseValue {
			return trueValue
		}
		return TriUnknown
	case rhino.FUNCTION, rhino.CLASS, rhino.NEW, rhino.ARRAYLIT, rhino.OBJECTLIT:
		return TriTrue
	case rhino.VOID:
		return TriFalse
	case rhino.GETPROP, rhino.GETELEM:
		// Assume that type information on getprops and getelems are likely
		// to be wrong. This prevents spurious warnings from not including
		// undefined in getelem's return value, from existence checks of
		// symbols the externs define as certainly true, or from default
		// initialization of globals (x.y = x.y || {}).
		// OPTCHAIN_GETELEM / OPTCHAIN_GETPROP are not in our token set.
		return TriUnknown
	}
	// If we reach this point then all the composite structures that we
	// can decompose have already been handled, leaving only qualified
	// names and type-aware checks to handle below.
	literalValue := getBooleanValue(n)
	if literalValue != TriUnknown || n.Token() == rhino.NAME {
		// If the truthiness is deterministic from the syntax then return
		// that immediately. Alternatively, NAME nodes also get a pass since
		// we don't trust the type information.
		return literalValue
	}
	// type-based check: nodes carry no JSType in this port, so always
	// UNKNOWN. Java: Tri.UNKNOWN.and(type.getPossibleToBooleanOutcomes().toTri())
	return TriUnknown
}

// ---- NodeUtil helpers (inlined, narrow port) -----------------------

// isNaN mirrors NodeUtil.isNaN(n) from Java. Returns true when the node
// is a NAME "NaN" or a numeric NaN literal.
func isNaN(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	if n.Token() == rhino.NAME && n.String_() == "NaN" {
		return true
	}
	if n.Token() == rhino.NUMBER && math.IsNaN(n.Number()) {
		return true
	}
	return false
}

// mayBeObject mirrors NodeUtil.mayBeObject(n) from Java. Returns false
// only when we can prove statically that n cannot be an object.
// Conservative: returns true for anything that might be an object.
func mayBeObject(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	switch n.Token() {
	// Definitely-non-object primitives.
	case rhino.NUMBER, rhino.STRINGLIT, rhino.TRUE, rhino.FALSE,
		rhino.NULL, rhino.VOID, rhino.TYPEOF:
		return false
	// NAME "undefined" is also not an object.
	case rhino.NAME:
		return n.String_() != "undefined"
	}
	// Anything else (CALL, GETPROP, GETELEM, THIS, NEW, OBJECTLIT,
	// ARRAYLIT, AND/OR results, etc.) may be an object.
	return true
}

// getBooleanValue mirrors NodeUtil.getBooleanValue(n) from Java. Returns
// the syntactic boolean value of a pure expression, or TriUnknown if
// indeterminate. Does not recurse into complex expressions beyond NOT.
func getBooleanValue(n *rhino.Node) Tri {
	if n == nil {
		return TriUnknown
	}
	switch n.Token() {
	case rhino.TRUE:
		return TriTrue
	case rhino.FALSE:
		return TriFalse
	case rhino.NULL:
		return TriFalse
	case rhino.VOID:
		return TriFalse
	case rhino.NUMBER:
		v := n.Number()
		if v == 0 || math.IsNaN(v) {
			return TriFalse
		}
		return TriTrue
	case rhino.STRINGLIT:
		if n.String_() == "" {
			return TriFalse
		}
		return TriTrue
	case rhino.NAME:
		switch n.String_() {
		case "undefined":
			return TriFalse
		case "NaN":
			return TriFalse
		case "Infinity":
			return TriTrue
		}
		return TriUnknown
	case rhino.NOT:
		inner := n.FirstChild()
		return getBooleanValue(inner).not()
	case rhino.FUNCTION, rhino.CLASS, rhino.ARRAYLIT, rhino.OBJECTLIT:
		return TriTrue
	}
	return TriUnknown
}

// loopCodeBlock mirrors NodeUtil.getLoopCodeBlock(n) from Java.
// Returns the BLOCK body of a loop statement.
//
// Child layout per token:
//   - WHILE(cond, body): body is child[1]
//   - FOR(init, cond, upd, body): body is last child
//   - FOR_IN(decl, iterable, body): body is last child
//   - FOR_OF(decl, iterable, body): body is last child
//   - FOR_AWAIT_OF: not in our token set, not handled
func loopCodeBlock(n *rhino.Node) *rhino.Node {
	switch n.Token() {
	case rhino.WHILE:
		kids := n.Children()
		if len(kids) >= 2 {
			return kids[1]
		}
	case rhino.FOR, rhino.FOR_IN, rhino.FOR_OF:
		return lastChild(n)
	}
	return nil
}

// lastChild returns the last child of n, or nil when n has no children.
// Mirrors Node.getLastChild() in Java.
func lastChild(n *rhino.Node) *rhino.Node {
	kids := n.Children()
	if len(kids) == 0 {
		return nil
	}
	return kids[len(kids)-1]
}

// secondChild returns the second child of n, or nil when n has fewer
// than 2 children. Mirrors Node.getSecondChild() in Java.
func secondChild(n *rhino.Node) *rhino.Node {
	kids := n.Children()
	if len(kids) < 2 {
		return nil
	}
	return kids[1]
}

// emit appends a diagnostic at node n's position.
func (c *CheckSuspiciousCode) emit(n *rhino.Node, code string, level diag.Level, msg string) {
	c.diags = append(c.diags, diag.Diagnostic{
		Path:    c.path,
		Line:    n.Lineno(),
		Col:     n.Charno(),
		Code:    code,
		Level:   level,
		Message: msg,
	})
}
