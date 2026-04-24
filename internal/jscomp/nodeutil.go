// Port of com/google/javascript/jscomp/NodeUtil.java @ closure-compiler v20260415 (subset; core helpers only)
//
// Includes: getBooleanValue (Tri-valued), getConditionExpression,
// getEnclosingFunction, getEnclosingScript, getEnclosingType,
// isStatement, isStatementBlock, isControlStructure, isImmutableValue,
// isLiteralValue, getFunctionName, getBestLValueName, getFunctionParameters,
// mayHaveSideEffects, isCall, isNew, isFunctionDeclaration, canBeSideEffected,
// isThis, isSuperOrSuperProperty, isGet, isPrototypeProperty,
// getAssignedValue, isBleedingFunctionName, isConstantDeclaration,
// isExpressionResultUsed, isFunctionExpression, isInFunction, isInSyntheticScript,
// isLhsByDestructuring, isObjectLitKey, mayBeObjectLitKey, isPrototypeMethod,
// mayBeString, mayEffectMutableState, newQualifiedNameNode, removeChildren,
// getBestLValue, getBestLValueOwner, getBestJSDocInfo, isNameDeclaration,
// isExprAssign, tryGetStringFromBoolean, isFunctionObjectLit, isGetPropOnThis,
// isNonEs6ClassExpression, isNameNodeInsideExtends.
//
// All helpers live on a zero-value nodeUtil receiver, exposed via the
// package-level NU sentinel. Call as NU.GetBooleanValue(n), etc.
// This avoids collision with any future package-level functions of the
// same names (the task spec explicitly prohibits such collisions).

package jscomp

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// TriV is a tri-state boolean mirroring com.google.javascript.jscomp.base.Tri.
// FALSE < UNKNOWN < TRUE (integer order matches Java's internal value).
type TriV int8

const (
	TriVFalse   TriV = -1
	TriVUnknown TriV = 0
	TriVTrue    TriV = 1
)

// ForBoolean converts a plain bool to TriV.  Mirrors Tri.forBoolean.
func ForBoolean(b bool) TriV {
	if b {
		return TriVTrue
	}
	return TriVFalse
}

// Not mirrors Tri.not().
func (t TriV) Not() TriV { return TriV(-int8(t)) }

// And mirrors Tri.and(): returns whichever is smaller (more false).
func (t TriV) And(other TriV) TriV {
	if t < other {
		return t
	}
	return other
}

// Or mirrors Tri.or(): returns whichever is larger (more true).
func (t TriV) Or(other TriV) TriV {
	if t > other {
		return t
	}
	return other
}

// ---------------------------------------------------------------------------
// nodeUtil namespace
// ---------------------------------------------------------------------------

// nodeUtil is the receiver for all NodeUtil helper methods. The zero value is
// ready to use — no fields, no state.
type nodeUtil struct{}

// NU is the package-level sentinel. Use NU.Method(n) at call sites.
var NU = nodeUtil{}

// ---------------------------------------------------------------------------
// getBooleanValue — Tri-valued
// ---------------------------------------------------------------------------

// GetBooleanValue returns the static boolean value of an expression node, or
// TriVUnknown when it cannot be determined without running the code.
// Mirrors NodeUtil.getBooleanValue (Java).
func (nodeUtil) GetBooleanValue(n *rhino.Node) TriV {
	if n == nil {
		return TriVUnknown
	}
	switch n.Token() {
	case rhino.NULL, rhino.FALSE, rhino.VOID:
		return TriVFalse
	case rhino.TRUE, rhino.FUNCTION, rhino.NEW, rhino.ARRAYLIT, rhino.OBJECTLIT:
		return TriVTrue
	case rhino.STRINGLIT:
		return ForBoolean(len(n.String_()) > 0)
	case rhino.NUMBER:
		v := n.Number()
		return ForBoolean(v != 0 && v == v)
	case rhino.UNDEFINED:
		// `undefined` is falsy.
		return TriVFalse
	case rhino.NAME:
		switch n.String_() {
		case "undefined", "NaN":
			return TriVFalse
		case "Infinity":
			return TriVTrue
		default:
			return TriVUnknown
		}
	case rhino.NOT:
		if fc := n.FirstChild(); fc != nil {
			return NU.GetBooleanValue(fc).Not()
		}
		return TriVUnknown
	case rhino.ASSIGN:
		// Value is the RHS (last child).  Java also handles COMMA here;
		// the Go rhino package does not yet define a COMMA token so we
		// skip that branch.
		kids := n.Children()
		if len(kids) >= 2 {
			return NU.GetBooleanValue(kids[len(kids)-1])
		}
		return TriVUnknown
	case rhino.AND:
		kids := n.Children()
		if len(kids) == 2 {
			return NU.GetBooleanValue(kids[0]).And(NU.GetBooleanValue(kids[1]))
		}
		return TriVUnknown
	case rhino.OR:
		kids := n.Children()
		if len(kids) == 2 {
			return NU.GetBooleanValue(kids[0]).Or(NU.GetBooleanValue(kids[1]))
		}
		return TriVUnknown
	case rhino.HOOK:
		kids := n.Children()
		if len(kids) == 3 {
			tv := NU.GetBooleanValue(kids[1])
			fv := NU.GetBooleanValue(kids[2])
			if tv == fv {
				return tv
			}
		}
		return TriVUnknown
	case rhino.COALESCE:
		kids := n.Children()
		if len(kids) == 2 {
			lhs := NU.GetBooleanValue(kids[0])
			rhs := NU.GetBooleanValue(kids[1])
			if lhs == TriVTrue || lhs == rhs {
				return lhs
			}
		}
		return TriVUnknown
	}
	return TriVUnknown
}

// ---------------------------------------------------------------------------
// getConditionExpression
// ---------------------------------------------------------------------------

// GetConditionExpression returns the condition node for IF/WHILE/DO/FOR nodes,
// or nil for FOR_IN / FOR_OF / CASE (no static condition node).
// Panics for any other token.  Mirrors NodeUtil.getConditionExpression (Java).
func (nodeUtil) GetConditionExpression(n *rhino.Node) *rhino.Node {
	if n == nil {
		return nil
	}
	kids := n.Children()
	switch n.Token() {
	case rhino.IF, rhino.WHILE:
		// condition is the first child
		if len(kids) >= 1 {
			return kids[0]
		}
		return nil
	case rhino.DO_WHILE:
		// condition is the last child
		if len(kids) >= 1 {
			return kids[len(kids)-1]
		}
		return nil
	case rhino.FOR:
		// init ; cond ; update — condition is second child
		if len(kids) >= 2 {
			return kids[1]
		}
		return nil
	case rhino.FOR_IN, rhino.FOR_OF, rhino.CASE:
		return nil
	}
	// Java throws here; we return nil to avoid panics in the port.
	return nil
}

// ---------------------------------------------------------------------------
// getEnclosingType / getEnclosingFunction / getEnclosingScript
// ---------------------------------------------------------------------------

// GetEnclosingType walks up the parent chain until it finds a node whose
// token matches `typ`, or returns nil.  Mirrors NodeUtil.getEnclosingType.
func (nodeUtil) GetEnclosingType(n *rhino.Node, typ rhino.Token) *rhino.Node {
	curr := n
	for curr != nil {
		if curr.Token() == typ {
			return curr
		}
		curr = curr.Parent()
	}
	return nil
}

// GetEnclosingFunction returns the nearest enclosing FUNCTION node, or nil.
// Mirrors NodeUtil.getEnclosingFunction.
func (nodeUtil) GetEnclosingFunction(n *rhino.Node) *rhino.Node {
	return NU.GetEnclosingType(n, rhino.FUNCTION)
}

// GetEnclosingScript returns the nearest enclosing SCRIPT node, or nil.
// Mirrors NodeUtil.getEnclosingScript.
func (nodeUtil) GetEnclosingScript(n *rhino.Node) *rhino.Node {
	return NU.GetEnclosingType(n, rhino.SCRIPT)
}

// ---------------------------------------------------------------------------
// isStatement / isStatementBlock / isControlStructure
// ---------------------------------------------------------------------------

// IsStatementBlock reports whether `n` is a node kind that contains statements.
// Mirrors NodeUtil.isStatementBlock.
func (nodeUtil) IsStatementBlock(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	switch n.Token() {
	case rhino.SCRIPT, rhino.BLOCK:
		return true
	}
	return false
}

// isStatementParent mirrors the Java IS_STATEMENT_PARENT set.
func isStatementParent(p *rhino.Node) bool {
	if p == nil {
		return false
	}
	switch p.Token() {
	case rhino.SCRIPT, rhino.BLOCK, rhino.LABEL:
		return true
	}
	return false
}

// IsStatement reports whether `n` is used as a statement.
// Mirrors NodeUtil.isStatement.
func (nodeUtil) IsStatement(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	if n.Token() == rhino.SCRIPT {
		return false
	}
	return isStatementParent(n.Parent())
}

// IsControlStructure reports whether `n` is a control-structure node
// (IF, WHILE, DO_WHILE, FOR, FOR_IN, FOR_OF, LABEL, TRY, CATCH,
// SWITCH, CASE, DEFAULT_CASE).
// Mirrors NodeUtil.isControlStructure.
func (nodeUtil) IsControlStructure(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	switch n.Token() {
	case rhino.FOR, rhino.FOR_IN, rhino.FOR_OF,
		rhino.DO_WHILE, rhino.WHILE, rhino.IF,
		rhino.LABEL, rhino.TRY, rhino.CATCH,
		rhino.SWITCH, rhino.CASE, rhino.DEFAULT_CASE:
		return true
	}
	return false
}

// ---------------------------------------------------------------------------
// isImmutableValue / isLiteralValue
// ---------------------------------------------------------------------------

// IsImmutableValue reports whether `n` is a node whose value is a compile-time
// constant that cannot be mutated at run time.
// Mirrors NodeUtil.isImmutableValue (Java).
func (nodeUtil) IsImmutableValue(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	switch n.Token() {
	case rhino.STRINGLIT, rhino.NUMBER, rhino.NULL, rhino.TRUE, rhino.FALSE:
		return true
	case rhino.NOT, rhino.VOID, rhino.NEG:
		if fc := n.FirstChild(); fc != nil {
			return NU.IsImmutableValue(fc)
		}
		return false
	case rhino.UNDEFINED:
		return true
	case rhino.NAME:
		s := n.String_()
		return s == "undefined" || s == "Infinity" || s == "NaN"
	}
	return false
}

// IsLiteralValue reports whether `n` is a literal value.  When
// includeFunctions is true, function expressions also qualify.
// Mirrors NodeUtil.isLiteralValue (Java).
func (nodeUtil) IsLiteralValue(n *rhino.Node, includeFunctions bool) bool {
	if n == nil {
		return false
	}
	switch n.Token() {
	case rhino.ARRAYLIT:
		for _, c := range n.Children() {
			if !NU.IsLiteralValue(c, includeFunctions) {
				return false
			}
		}
		return true
	case rhino.OBJECTLIT:
		for _, c := range n.Children() {
			switch c.Token() {
			case rhino.STRING_KEY:
				if fc := c.FirstChild(); fc != nil && !NU.IsLiteralValue(fc, includeFunctions) {
					return false
				}
			default:
				// computed / spread / getter / setter — conservatively false
				return false
			}
		}
		return true
	case rhino.FUNCTION:
		return includeFunctions && !NU.IsFunctionDeclaration(n)
	}
	return NU.IsImmutableValue(n)
}

// ---------------------------------------------------------------------------
// getFunctionName / getBestLValueName
// ---------------------------------------------------------------------------

// GetFunctionName returns the inferred name of a FUNCTION node, or "".
// Mirrors NodeUtil.getNearestFunctionName (Java) — the task calls this
// getFunctionName for the subset we expose.
func (nodeUtil) GetFunctionName(n *rhino.Node) string {
	if n == nil || n.Token() != rhino.FUNCTION {
		return ""
	}
	// Named function expression / declaration: first child is NAME.
	kids := n.Children()
	if len(kids) >= 1 {
		nameNode := kids[0]
		if nameNode.Token() == rhino.NAME && nameNode.String_() != "" {
			return nameNode.String_()
		}
	}
	// Infer from the parent: `var f = function(){}` or `{f: function(){}}`
	parent := n.Parent()
	if parent == nil {
		return ""
	}
	switch parent.Token() {
	case rhino.NAME:
		return parent.String_()
	case rhino.STRING_KEY:
		return parent.String_()
	case rhino.ASSIGN, rhino.GETPROP:
		// `Foo.bar = function(){}` — return last GETPROP string
		if parent.Token() == rhino.GETPROP {
			return parent.String_()
		}
		lhs := parent.FirstChild()
		if lhs != nil && lhs.Token() == rhino.GETPROP {
			return lhs.String_()
		}
	}
	return ""
}

// GetBestLValueName returns the best inferred name for an l-value node.
// For NAME nodes this is the identifier; for GETPROP chains it is the
// qualified name.  Returns "" when no name can be determined.
// Mirrors NodeUtil.getBestLValueName (Java).
func (nodeUtil) GetBestLValueName(n *rhino.Node) string {
	if n == nil || n.Parent() == nil {
		return ""
	}
	switch n.Token() {
	case rhino.NAME:
		return n.String_()
	case rhino.GETPROP:
		// qualified.name — reconstruct from the chain
		return nuQualifiedName(n)
	case rhino.STRING_KEY:
		// { key: fn } — owner name + "." + key
		owner := n.Parent()
		if owner == nil {
			return ""
		}
		ownerName := NU.GetBestLValueName(owner)
		if ownerName != "" {
			return ownerName + "." + n.String_()
		}
	}
	return ""
}

// nuQualifiedName reconstructs a.b.c from a GETPROP chain.
func nuQualifiedName(n *rhino.Node) string {
	if n == nil {
		return ""
	}
	if n.Token() == rhino.NAME {
		return n.String_()
	}
	if n.Token() != rhino.GETPROP {
		return ""
	}
	fc := n.FirstChild()
	base := nuQualifiedName(fc)
	if base == "" {
		return ""
	}
	return base + "." + n.String_()
}

// ---------------------------------------------------------------------------
// getFunctionParameters
// ---------------------------------------------------------------------------

// GetFunctionParameters returns the PARAM_LIST node of a FUNCTION node.
// Mirrors NodeUtil.getFunctionParameters (Java).
func (nodeUtil) GetFunctionParameters(fnNode *rhino.Node) *rhino.Node {
	if fnNode == nil || fnNode.Token() != rhino.FUNCTION {
		return nil
	}
	kids := fnNode.Children()
	// FUNCTION children: [NAME, PARAM_LIST, BLOCK]
	if len(kids) >= 2 {
		return kids[1]
	}
	return nil
}

// ---------------------------------------------------------------------------
// mayHaveSideEffects
// ---------------------------------------------------------------------------

// MayHaveSideEffects reports conservatively whether evaluating `n` can
// produce observable side effects.  Mirrors the narrow form of
// AstAnalyzer.mayHaveSideEffects (Java) — no @nosideeffects JSDoc support.
func (nodeUtil) MayHaveSideEffects(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	switch n.Token() {
	// Definitely side-effecting.
	case rhino.CALL, rhino.NEW:
		return true
	case rhino.THROW:
		return true
	case rhino.DELPROP:
		return true
	// Assignments and compound-assigns mutate state.
	case rhino.ASSIGN, rhino.ASSIGN_ADD, rhino.ASSIGN_SUB, rhino.ASSIGN_MUL,
		rhino.ASSIGN_DIV, rhino.ASSIGN_MOD,
		rhino.ASSIGN_BITAND, rhino.ASSIGN_BITOR, rhino.ASSIGN_BITXOR,
		rhino.ASSIGN_LSH, rhino.ASSIGN_RSH, rhino.ASSIGN_URSH:
		return true
	// Declarations change the namespace.
	case rhino.VAR, rhino.LET, rhino.CONST:
		return true
	// FOR_IN / FOR_OF are almost always side-effecting (they represent
	// an assignment on each iteration).
	case rhino.FOR_IN, rhino.FOR_OF:
		return true
	// Pure leaves.
	case rhino.NAME:
		// NAME with children is a var/let/const initialiser — side-effecting.
		return len(n.Children()) > 0
	case rhino.FUNCTION:
		// Function expressions are side-effect-free; declarations change scope.
		return NU.IsFunctionDeclaration(n)
	case rhino.THIS, rhino.STRINGLIT, rhino.NUMBER, rhino.TRUE, rhino.FALSE,
		rhino.NULL, rhino.UNDEFINED:
		return false
	case rhino.ARRAYLIT, rhino.OBJECTLIT:
		// Creating new mutable objects is conservative-pure when we only care
		// about externally-visible side effects.
		for _, c := range n.Children() {
			if NU.MayHaveSideEffects(c) {
				return true
			}
		}
		return false
	}
	// Compound: side-effecting if any child is.
	for _, c := range n.Children() {
		if NU.MayHaveSideEffects(c) {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------------------------
// isCall / isNew / isFunctionDeclaration
// ---------------------------------------------------------------------------

// IsCall reports whether `n` is a CALL node.
// Mirrors Node.isCall() in Java.
func (nodeUtil) IsCall(n *rhino.Node) bool {
	return n != nil && n.Token() == rhino.CALL
}

// IsNew reports whether `n` is a NEW node.
// Mirrors Node.isNew() in Java.
func (nodeUtil) IsNew(n *rhino.Node) bool {
	return n != nil && n.Token() == rhino.NEW
}

// IsFunctionDeclaration reports whether `n` is a function declaration —
// a FUNCTION node that is a direct child of a statement-parent and has
// a non-empty name.  Mirrors NodeUtil.isFunctionDeclaration (Java).
func (nodeUtil) IsFunctionDeclaration(n *rhino.Node) bool {
	if n == nil || n.Token() != rhino.FUNCTION {
		return false
	}
	// Must have a non-empty name (first child NAME with non-empty string).
	kids := n.Children()
	if len(kids) < 1 {
		return false
	}
	nameNode := kids[0]
	if nameNode.Token() != rhino.NAME || nameNode.String_() == "" {
		return false
	}
	// Parent must be a declaration parent (statement parent).
	parent := n.Parent()
	return isStatementParent(parent)
}

// ---------------------------------------------------------------------------
// canBeSideEffected
// ---------------------------------------------------------------------------

// CanBeSideEffected reports whether evaluating `n` can be *affected* (not
// just affect) by concurrent side effects — i.e. whether its value might
// change between two evaluations separated by effectful code.
// Mirrors NodeUtil.canBeSideEffected (Java), with knownConstants = {}.
func (nodeUtil) CanBeSideEffected(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	switch n.Token() {
	case rhino.CALL, rhino.NEW:
		return true
	case rhino.NAME:
		// Non-constant names can be changed between reads.
		s := n.String_()
		return s != "undefined" && s != "Infinity" && s != "NaN"
	case rhino.GETPROP:
		// Property access on anything other than Symbol is potentially affected.
		fc := n.FirstChild()
		if fc != nil && fc.Token() == rhino.NAME && fc.String_() == "Symbol" {
			return false
		}
		return true
	case rhino.GETELEM:
		return true
	case rhino.FUNCTION:
		// Function expressions are not changed by side-effects.
		return false
	}
	for _, c := range n.Children() {
		if NU.CanBeSideEffected(c) {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------------------------
// isThis / isSuperOrSuperProperty / isGet / isPrototypeProperty
// ---------------------------------------------------------------------------

// IsThis reports whether `n` is a THIS node.
// Mirrors Node.isThis() in Java.
func (nodeUtil) IsThis(n *rhino.Node) bool {
	return n != nil && n.Token() == rhino.THIS
}

// IsSuperOrSuperProperty reports whether `n` is a SUPER node or a property
// access whose object is SUPER (e.g. `super.method`).
// There is no direct Java counterpart by this exact name; the predicate
// is spelled out inline in Java's NodeUtil in several places.
func (nodeUtil) IsSuperOrSuperProperty(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	// TODO: SUPER token not yet in the Go rhino package; always false
	// until rhino.SUPER is added.
	if n.Token() == rhino.GETPROP {
		fc := n.FirstChild()
		// fc is SUPER when the Go rhino package exposes it.
		_ = fc
	}
	return false
}

// IsGet reports whether `n` is a property-access node (GETPROP or GETELEM).
// Mirrors NodeUtil.isNormalGet (Java).
func (nodeUtil) IsGet(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	return n.Token() == rhino.GETPROP || n.Token() == rhino.GETELEM
}

// IsPrototypeProperty reports whether `n` is a qualified prototype property
// — i.e. a GETPROP whose object is itself a GETPROP named "prototype".
// Example: `Foo.prototype.bar`.
// Mirrors NodeUtil.isPrototypeProperty (Java).
func (nodeUtil) IsPrototypeProperty(n *rhino.Node) bool {
	if n == nil || n.Token() != rhino.GETPROP {
		return false
	}
	recv := n.FirstChild()
	if recv == nil {
		return false
	}
	return recv.Token() == rhino.GETPROP && recv.String_() == "prototype"
}

// MayBeString reports whether n could evaluate to a string.
// Mirrors NodeUtil.mayBeString(Node) (Java).
func (nodeUtil) MayBeString(n *rhino.Node) bool {
	if n == nil {
		return true
	}
	switch n.Token() {
	case rhino.NUMBER, rhino.URSH, rhino.POS, rhino.INC, rhino.DEC,
		rhino.BITNOT, rhino.NEG, rhino.BITOR, rhino.BITXOR, rhino.BITAND,
		rhino.LSH, rhino.RSH, rhino.SUB, rhino.MUL, rhino.MOD, rhino.DIV,
		rhino.TRUE, rhino.FALSE, rhino.EQ, rhino.NE, rhino.SHEQ, rhino.SHNE,
		rhino.LT, rhino.GT, rhino.LE, rhino.GE, rhino.IN, rhino.INSTANCEOF,
		rhino.NOT, rhino.DELPROP, rhino.NULL, rhino.VOID:
		return false
	case rhino.NAME:
		s := n.String_()
		return s != "undefined" && s != "NaN" && s != "Infinity"
	}
	return true
}

// IsNameDeclaration reports whether n is VAR, LET, or CONST. Mirrors NodeUtil.isNameDeclaration.
func (nodeUtil) IsNameDeclaration(n *rhino.Node) bool {
	return n != nil && (n.Token() == rhino.VAR || n.Token() == rhino.LET || n.Token() == rhino.CONST)
}

// IsExprAssign reports whether n is EXPR_RESULT wrapping ASSIGN. Mirrors NodeUtil.isExprAssign.
func (nodeUtil) IsExprAssign(n *rhino.Node) bool {
	if n == nil || n.Token() != rhino.EXPR_RESULT {
		return false
	}
	fc := n.FirstChild()
	return fc != nil && fc.Token() == rhino.ASSIGN
}

// IsFunctionExpression reports whether n is a non-declaration FUNCTION. Mirrors NodeUtil.isFunctionExpression.
func (nodeUtil) IsFunctionExpression(n *rhino.Node) bool {
	return n != nil && n.Token() == rhino.FUNCTION && !NU.IsFunctionDeclaration(n)
}

// GetAssignedValue returns the RHS assigned to NAME/GETPROP lvalue n, or nil.
// Mirrors NodeUtil.getAssignedValue (Java).
func (nodeUtil) GetAssignedValue(n *rhino.Node) *rhino.Node {
	if n == nil {
		return nil
	}
	p := n.Parent()
	if p == nil {
		return nil
	}
	if NU.IsNameDeclaration(p) {
		return n.FirstChild()
	}
	if p.Token() == rhino.ASSIGN && p.FirstChild() == n {
		return n.Next()
	}
	return nil
}

// IsBleedingFunctionName reports whether n is the bleeding NAME of a named
// function expression. Mirrors NodeUtil.isBleedingFunctionName (Java).
func (nodeUtil) IsBleedingFunctionName(n *rhino.Node) bool {
	if n == nil || n.Token() != rhino.NAME || n.String_() == "" {
		return false
	}
	p := n.Parent()
	return p != nil && NU.IsFunctionExpression(p) && p.FirstChild() == n
}

// IsConstantDeclaration reports whether n is a constant lvalue.
// Mirrors NodeUtil.isConstantDeclaration (Java).
func (nodeUtil) IsConstantDeclaration(info *rhino.JSDocInfo, n *rhino.Node) bool {
	if n == nil {
		return false
	}
	p := n.Parent()
	if NU.IsObjectLitKey(n) ||
		(p != nil && p.Token() == rhino.ASSIGN && p.FirstChild() == n) ||
		(p != nil && p.Token() == rhino.EXPR_RESULT && NU.IsGet(n)) {
		return info != nil && info.IsConstant()
	}
	if n.Token() != rhino.NAME {
		return false
	}
	if p != nil && p.Token() == rhino.CONST {
		return true
	}
	return (info != nil && info.IsConstant()) || n.IsInferredConstantVar()
}

// IsExpressionResultUsed reports whether the result of expr is consumed.
// Mirrors NodeUtil.isExpressionResultUsed (Java).
func (nodeUtil) IsExpressionResultUsed(expr *rhino.Node) bool {
	if expr == nil {
		return false
	}
	p := expr.Parent()
	if p == nil {
		return true
	}
	switch p.Token() {
	case rhino.BLOCK, rhino.EXPR_RESULT:
		return false
	case rhino.CAST:
		return NU.IsExpressionResultUsed(p)
	case rhino.HOOK, rhino.AND, rhino.OR, rhino.COALESCE:
		return expr == p.FirstChild() || NU.IsExpressionResultUsed(p)
	case rhino.FOR:
		return p.GetSecondChild() == expr
	}
	return true
}

// IsInFunction reports whether n is nested inside a FUNCTION node.
// Mirrors NodeUtil.isInFunction (Java).
func (nodeUtil) IsInFunction(n *rhino.Node) bool { return NU.GetEnclosingFunction(n) != nil }

// IsInSyntheticScript always returns false; rhino.Node has no source-file info.
// Mirrors NodeUtil.isInSyntheticScript (Java).
func (nodeUtil) IsInSyntheticScript(_ *rhino.Node) bool { return false }

// IsLhsByDestructuring reports whether n is an LHS target in a destructuring pattern.
// Mirrors NodeUtil.isLhsByDestructuring (Java).
func (nodeUtil) IsLhsByDestructuring(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	switch n.Token() {
	case rhino.NAME, rhino.GETPROP, rhino.GETELEM:
		p, gp := n.Parent(), n.GetGrandparent()
		if p == nil {
			return false
		}
		switch p.Token() {
		case rhino.ITER_REST:
			return true
		case rhino.COMPUTED_PROP:
			return !n.IsFirstChildOf(p) && gp != nil && gp.Token() == rhino.OBJECTLIT
		case rhino.STRING_KEY:
			return gp != nil && gp.Token() == rhino.OBJECTLIT
		}
	}
	return false
}

// MayBeObjectLitKey reports whether n could be an object-literal key.
// Mirrors NodeUtil.mayBeObjectLitKey (Java).
func (nodeUtil) MayBeObjectLitKey(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	switch n.Token() {
	case rhino.STRING_KEY, rhino.MEMBER_FUNCTION_DEF, rhino.COMPUTED_PROP:
		return true
	}
	return n.IsGetterDef() || n.IsSetterDef()
}

// IsObjectLitKey reports whether n is a key inside an OBJECTLIT node.
// Mirrors NodeUtil.isObjectLitKey (Java).
func (nodeUtil) IsObjectLitKey(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	p := n.Parent()
	return p != nil && p.Token() == rhino.OBJECTLIT && NU.MayBeObjectLitKey(n)
}

// IsPrototypeMethod reports whether n is a FUNCTION assigned as a prototype method.
// Mirrors NodeUtil.isPrototypeMethod (Java).
func (nodeUtil) IsPrototypeMethod(n *rhino.Node) bool {
	if n == nil || n.Token() != rhino.FUNCTION {
		return false
	}
	assign := n.Parent()
	if assign == nil || assign.Token() != rhino.ASSIGN {
		return false
	}
	er := assign.Parent()
	lhs := er.FirstChild()
	return NU.IsExprAssign(er) && lhs != nil && NU.IsPrototypeProperty(lhs.FirstChild())
}

// MayEffectMutableState reports whether n may create or mutate state.
// Mirrors AstAnalyzer.mayEffectMutableState (Java).
func (nodeUtil) MayEffectMutableState(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	switch n.Token() {
	case rhino.CALL, rhino.NEW, rhino.THROW, rhino.FOR_IN, rhino.FOR_OF,
		rhino.VAR, rhino.LET, rhino.CONST,
		rhino.ASSIGN, rhino.ASSIGN_ADD, rhino.ASSIGN_SUB, rhino.ASSIGN_MUL,
		rhino.ASSIGN_DIV, rhino.ASSIGN_MOD,
		rhino.ASSIGN_BITAND, rhino.ASSIGN_BITOR, rhino.ASSIGN_BITXOR,
		rhino.ASSIGN_LSH, rhino.ASSIGN_RSH, rhino.ASSIGN_URSH,
		rhino.ARRAYLIT, rhino.OBJECTLIT, rhino.REGEXP:
		return true
	case rhino.NAME:
		return n.HasChildren()
	case rhino.FUNCTION:
		return NU.IsFunctionDeclaration(n)
	}
	for _, c := range n.Children() {
		if NU.MayEffectMutableState(c) {
			return true
		}
	}
	return false
}

// NewQualifiedNameNode builds a NAME or GETPROP chain for a dotted name.
// Mirrors NodeUtil.newQName (Java) without compiler context.
func (nodeUtil) NewQualifiedNameNode(name string) *rhino.Node {
	dot := strings.IndexByte(name, '.')
	if dot == -1 {
		return rhino.NewString(rhino.NAME, name)
	}
	root := rhino.NewString(rhino.NAME, name[:dot])
	for rest := name[dot+1:]; ; {
		d := strings.IndexByte(rest, '.')
		var part string
		if d == -1 {
			part, rest = rest, ""
		} else {
			part, rest = rest[:d], rest[d+1:]
		}
		prop := rhino.NewString(rhino.GETPROP, part)
		prop.AddChildToBack(root)
		root = prop
		if d == -1 {
			break
		}
	}
	return root
}

// RemoveChildren detaches and returns all children of n. Mirrors Node.removeChildren().
func (nodeUtil) RemoveChildren(n *rhino.Node) []*rhino.Node {
	if n == nil {
		return nil
	}
	return n.RemoveChildren()
}

// GetBestLValue returns the l-value node that n (r-value) is being assigned to.
// Mirrors NodeUtil.getBestLValue (Java).
func (nodeUtil) GetBestLValue(n *rhino.Node) *rhino.Node {
	if n == nil {
		return nil
	}
	if NU.IsFunctionDeclaration(n) {
		return n.FirstChild()
	}
	p := n.Parent()
	if p == nil {
		return nil
	}
	switch p.Token() {
	case rhino.NAME:
		return p
	case rhino.ASSIGN:
		return p.FirstChild()
	case rhino.HOOK, rhino.AND, rhino.OR, rhino.COALESCE:
		if p.FirstChild() != n {
			return NU.GetBestLValue(p)
		}
	case rhino.CAST:
		return NU.GetBestLValue(p)
	}
	if NU.MayBeObjectLitKey(p) || p.Token() == rhino.COMPUTED_PROP {
		return p
	}
	return nil
}

// GetBestLValueOwner returns the owner node of an l-value.
// Mirrors NodeUtil.getBestLValueOwner (Java).
func (nodeUtil) GetBestLValueOwner(lv *rhino.Node) *rhino.Node {
	if lv == nil || lv.Parent() == nil {
		return nil
	}
	if NU.MayBeObjectLitKey(lv) || lv.Token() == rhino.COMPUTED_PROP {
		return NU.GetBestLValue(lv.Parent())
	}
	if NU.IsGet(lv) {
		return lv.FirstChild()
	}
	return nil
}

// GetBestJSDocInfo returns the JSDocInfo most closely associated with n.
// Mirrors NodeUtil.getBestJSDocInfo (Java).
func (nodeUtil) GetBestJSDocInfo(n *rhino.Node) *rhino.JSDocInfo {
	if n == nil {
		return nil
	}
	if n.Token() == rhino.EXPR_RESULT {
		if fc := n.FirstChild(); fc != nil {
			return NU.GetBestJSDocInfo(fc)
		}
		return nil
	}
	if info := n.JSDocInfo(); info != nil {
		return info
	}
	p := n.Parent()
	if p == nil {
		return nil
	}
	switch p.Token() {
	case rhino.NAME, rhino.VAR, rhino.LET, rhino.CONST, rhino.DECLARE, rhino.ASSIGN:
		return NU.GetBestJSDocInfo(p)
	case rhino.FUNCTION, rhino.CLASS:
		if n == p.FirstChild() {
			return NU.GetBestJSDocInfo(p)
		}
	}
	if NU.MayBeObjectLitKey(p) || p.Token() == rhino.COMPUTED_PROP {
		return p.JSDocInfo()
	}
	return nil
}

// TryGetStringFromBoolean returns "true"/"false" for TRUE/FALSE nodes, else "".
func (nodeUtil) TryGetStringFromBoolean(n *rhino.Node) string {
	if n == nil {
		return ""
	}
	switch n.Token() {
	case rhino.TRUE:
		return "true"
	case rhino.FALSE:
		return "false"
	}
	return ""
}

// IsFunctionObjectLit reports whether n is a FUNCTION whose sole body statement
// is EXPR_RESULT(OBJECTLIT).
func (nodeUtil) IsFunctionObjectLit(n *rhino.Node) bool {
	if n == nil || n.Token() != rhino.FUNCTION {
		return false
	}
	block := n.LastChild()
	if block == nil || block.Token() != rhino.BLOCK {
		return false
	}
	kids := block.Children()
	return len(kids) == 1 && kids[0].Token() == rhino.EXPR_RESULT &&
		kids[0].FirstChild() != nil && kids[0].FirstChild().Token() == rhino.OBJECTLIT
}

// IsGetPropOnThis reports whether n is GETPROP with THIS as receiver.
func (nodeUtil) IsGetPropOnThis(n *rhino.Node) bool {
	if n == nil || n.Token() != rhino.GETPROP {
		return false
	}
	fc := n.FirstChild()
	return fc != nil && fc.Token() == rhino.THIS
}

// IsNonEs6ClassExpression reports whether n is a class expression (not a class declaration).
func (nodeUtil) IsNonEs6ClassExpression(n *rhino.Node) bool {
	if n == nil || n.Token() != rhino.CLASS {
		return false
	}
	kids := n.Children()
	if len(kids) >= 1 && kids[0].Token() == rhino.NAME && kids[0].String_() != "" {
		return !isStatementParent(n.Parent())
	}
	return true
}

// IsNameNodeInsideExtends reports whether n is the NAME in a class extends clause.
func (nodeUtil) IsNameNodeInsideExtends(n *rhino.Node) bool {
	if n == nil || n.Token() != rhino.NAME {
		return false
	}
	p := n.Parent()
	return p != nil && p.Token() == rhino.CLASS && p.GetSecondChild() == n
}
