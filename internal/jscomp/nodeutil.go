// Port of com/google/javascript/jscomp/NodeUtil.java @ closure-compiler v20260415 (subset; core helpers only)
//
// Includes: getBooleanValue (Tri-valued), getConditionExpression,
// getEnclosingFunction, getEnclosingScript, getEnclosingType, getEnclosingScopeRoot,
// isStatement, isStatementBlock, isControlStructure, isImmutableValue,
// isLiteralValue, getFunctionName, getBestLValueName, getFunctionParameters,
// getFunctionArgs, getFunctionBody, mayHaveSideEffects, isCall, isNew,
// isFunctionDeclaration, canBeSideEffected,
// isThis, isSuperOrSuperProperty, isGet, isPrototypeProperty,
// getAssignedValue, isBleedingFunctionName, isConstantDeclaration,
// isExpressionResultUsed, isFunctionExpression, isInFunction, isInSyntheticScript,
// isLhsByDestructuring, isObjectLitKey, mayBeObjectLitKey, isPrototypeMethod,
// mayBeString, mayEffectMutableState, newQualifiedNameNode, removeChildren,
// getBestLValue, getBestLValueOwner, getBestJSDocInfo, isNameDeclaration,
// isExprAssign, tryGetStringFromBoolean, isFunctionObjectLit, isGetPropOnThis,
// isNonEs6ClassExpression, isNameNodeInsideExtends,
// getLValuesOfAssign, getLValuesOfVar, getRValueOfLValue, getRValueFromLValue,
// isCallTarget, isCallOrNewTarget, isInvocation, isInvocationTarget,
// isNamespaceDecl, isFunctionObjectCall, isFunctionObjectApply,
// isHoistedFunctionDeclaration, isExprCall, isCallOrNew, isAnyFor,
// isTopLevelFunction, isFunctionBody, isFunctionBlock,
// isAnonymousFunction,
// isGoogProvideCall, isGoogRequireCall, isGoogModuleCall,
// getGoogProvideName, getGoogRequireName,
// getStringValue, getPrefixedName,
// newName, newDecl.
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

// ---------------------------------------------------------------------------
// scope navigation
// ---------------------------------------------------------------------------

// nuCreatesScope mirrors NodeUtil.createsScope (package-private in Java).
// Returns true when n introduces a new lexical scope.
func nuCreatesScope(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	switch n.Token() {
	case rhino.BLOCK:
		// Synthetic blocks (added by the compiler) do not create scopes.
		// Mirror Java's Node.isSyntheticBlock() which checks Node.SYNTHETIC.
		if n.GetBooleanProp(rhino.PropSynthetic) {
			return false
		}
		p := n.Parent()
		if p == nil {
			return false
		}
		// Switch cases (CASE/DEFAULT_CASE) and catch bodies handled by parent.
		switch p.Token() {
		case rhino.CASE, rhino.DEFAULT_CASE, rhino.CATCH:
			return false
		}
		return true
	case rhino.FOR, rhino.FOR_IN, rhino.FOR_OF, rhino.FOR_AWAIT_OF,
		rhino.SWITCH_BODY, rhino.CLASS:
		return true
	case rhino.FUNCTION, rhino.MODULE_BODY:
		return true
	case rhino.ROOT:
		return n.Parent() == nil
	}
	return false
}

// GetEnclosingScopeRoot walks up the parent chain and returns the first node
// that creates a scope (function, block, module body, root), or nil.
// Mirrors NodeUtil.getEnclosingScopeRoot.
func (nodeUtil) GetEnclosingScopeRoot(n *rhino.Node) *rhino.Node {
	curr := n
	for curr != nil {
		if nuCreatesScope(curr) {
			return curr
		}
		curr = curr.Parent()
	}
	return nil
}

// ---------------------------------------------------------------------------
// function body / function args / getFunctionBody
// ---------------------------------------------------------------------------

// GetFunctionBody returns the BLOCK node that is the body of fn.
// fn must be a FUNCTION node. Mirrors NodeUtil.getFunctionBody.
func (nodeUtil) GetFunctionBody(fn *rhino.Node) *rhino.Node {
	if fn == nil || fn.Token() != rhino.FUNCTION {
		return nil
	}
	return fn.LastChild()
}

// GetFunctionArgs returns the first child of the PARAM_LIST of fn (i.e. the
// first formal parameter node), or nil.  Companion to GetFunctionParameters.
func (nodeUtil) GetFunctionArgs(fn *rhino.Node) *rhino.Node {
	params := NU.GetFunctionParameters(fn)
	if params == nil {
		return nil
	}
	return params.FirstChild()
}

// ---------------------------------------------------------------------------
// l-value / r-value helpers
// ---------------------------------------------------------------------------

// GetLValuesOfVar collects the NAME nodes that a VAR/LET/CONST declaration
// binds.  Simple (non-destructuring) declarations yield one NAME; destructuring
// yields all leaf NAME nodes.  Mirrors the spirit of NodeUtil.visitLhsNodesInNode
// for name-declaration nodes.
func (nodeUtil) GetLValuesOfVar(decl *rhino.Node) []*rhino.Node {
	if decl == nil {
		return nil
	}
	switch decl.Token() {
	case rhino.VAR, rhino.LET, rhino.CONST:
	default:
		return nil
	}
	var out []*rhino.Node
	for child := decl.FirstChild(); child != nil; child = child.Next() {
		nuCollectLValues(child, &out)
	}
	return out
}

// GetLValuesOfAssign collects the NAME/GETPROP/GETELEM nodes on the left-hand
// side of an assignment or destructuring assignment.  Mirrors the spirit of
// NodeUtil.visitLhsNodesInNode for assignment nodes.
func (nodeUtil) GetLValuesOfAssign(assign *rhino.Node) []*rhino.Node {
	if assign == nil {
		return nil
	}
	switch assign.Token() {
	case rhino.ASSIGN, rhino.ASSIGN_ADD, rhino.ASSIGN_SUB, rhino.ASSIGN_MUL,
		rhino.ASSIGN_DIV, rhino.ASSIGN_MOD,
		rhino.ASSIGN_BITAND, rhino.ASSIGN_BITOR, rhino.ASSIGN_BITXOR,
		rhino.ASSIGN_LSH, rhino.ASSIGN_RSH, rhino.ASSIGN_URSH,
		rhino.ASSIGN_EXPONENT, rhino.ASSIGN_OR, rhino.ASSIGN_AND,
		rhino.ASSIGN_COALESCE:
	default:
		return nil
	}
	lhs := assign.FirstChild()
	if lhs == nil {
		return nil
	}
	var out []*rhino.Node
	nuCollectLValues(lhs, &out)
	return out
}

// nuCollectLValues recursively collects leaf l-value nodes from a potentially
// destructuring pattern or a simple name/getprop/getelem.
func nuCollectLValues(n *rhino.Node, out *[]*rhino.Node) {
	if n == nil {
		return
	}
	switch n.Token() {
	case rhino.NAME, rhino.GETPROP, rhino.GETELEM:
		*out = append(*out, n)
	case rhino.ARRAY_PATTERN:
		for child := n.FirstChild(); child != nil; child = child.Next() {
			nuCollectLValues(child, out)
		}
	case rhino.OBJECT_PATTERN:
		for child := n.FirstChild(); child != nil; child = child.Next() {
			// STRING_KEY holds the actual binding in its last child.
			if child.Token() == rhino.STRING_KEY {
				nuCollectLValues(child.LastChild(), out)
			} else {
				nuCollectLValues(child, out)
			}
		}
	case rhino.DEFAULT_VALUE:
		nuCollectLValues(n.FirstChild(), out)
	case rhino.ITER_REST, rhino.OBJECT_REST:
		nuCollectLValues(n.FirstChild(), out)
	}
}

// GetRValueOfLValue returns the right-hand side expression for a given l-value
// node n, or nil if n is not a recognisable l-value context.
// Mirrors NodeUtil.getRValueOfLValue.
func (nodeUtil) GetRValueOfLValue(n *rhino.Node) *rhino.Node {
	if n == nil {
		return nil
	}
	p := n.Parent()
	if p == nil {
		return nil
	}
	switch p.Token() {
	case rhino.ASSIGN, rhino.ASSIGN_ADD, rhino.ASSIGN_SUB, rhino.ASSIGN_MUL,
		rhino.ASSIGN_DIV, rhino.ASSIGN_MOD,
		rhino.ASSIGN_BITAND, rhino.ASSIGN_BITOR, rhino.ASSIGN_BITXOR,
		rhino.ASSIGN_LSH, rhino.ASSIGN_RSH, rhino.ASSIGN_URSH,
		rhino.ASSIGN_EXPONENT, rhino.ASSIGN_OR, rhino.ASSIGN_AND,
		rhino.ASSIGN_COALESCE:
		// n is the LHS; the RHS is the next sibling.
		return n.Next()
	case rhino.VAR, rhino.LET, rhino.CONST:
		// NAME node inside a declaration: initialiser is last child.
		return n.LastChild()
	case rhino.OBJECTLIT:
		// A STRING_KEY / COMPUTED_PROP / member inside an object literal.
		return n.LastChild()
	case rhino.CLASS_MEMBERS:
		// Method/field inside a class body.
		return n.LastChild()
	case rhino.FUNCTION, rhino.CLASS:
		// The name node of a named function/class expression — the parent IS the value.
		return p
	}
	return nil
}

// GetRValueFromLValue is an alias for GetRValueOfLValue that accepts the
// l-value node directly.  Provided as a convenience companion matching the
// Java naming variant used in some callers.
func (nodeUtil) GetRValueFromLValue(n *rhino.Node) *rhino.Node {
	return NU.GetRValueOfLValue(n)
}

// ---------------------------------------------------------------------------
// call / invocation target predicates
// ---------------------------------------------------------------------------

// IsCallTarget reports whether n is the target (callee) of a CALL expression,
// including the pattern `x.call(...)` where n is `x`.
// Mirrors OptimizeCalls.ReferenceMap.isCallTarget.
func (nodeUtil) IsCallTarget(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	p := n.Parent()
	if p == nil {
		return false
	}
	// Direct call: n is the first child of a CALL.
	if p.Token() == rhino.CALL && p.FirstChild() == n {
		return true
	}
	// x.call(...) pattern: n is the receiver of a GETPROP("call") that is
	// itself the callee of a CALL.
	gp := p.Parent()
	return p.Token() == rhino.GETPROP &&
		p.String_() == "call" &&
		gp != nil && gp.Token() == rhino.CALL &&
		p.IsFirstChildOf(gp)
}

// IsCallOrNewTarget reports whether n is the target of a CALL or NEW expression
// (or an optional-chain call), including the `.call(...)` pattern.
// Mirrors NodeUtil's isInvocationTarget extended to cover NEW.
func (nodeUtil) IsCallOrNewTarget(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	p := n.Parent()
	if p == nil {
		return false
	}
	switch p.Token() {
	case rhino.CALL, rhino.NEW, rhino.OPTCHAIN_CALL, rhino.TAGGED_TEMPLATELIT:
		return p.FirstChild() == n
	}
	// x.call(...) / x.apply(...) pattern.
	gp := p.Parent()
	if p.Token() == rhino.GETPROP && gp != nil &&
		(gp.Token() == rhino.CALL || gp.Token() == rhino.OPTCHAIN_CALL) &&
		p.IsFirstChildOf(gp) {
		m := p.String_()
		return m == "call" || m == "apply"
	}
	return false
}

// IsInvocation reports whether n is a call-like invocation node: CALL, NEW,
// OPTCHAIN_CALL, or TAGGED_TEMPLATELIT.
// Mirrors NodeUtil.isInvocation.
func (nodeUtil) IsInvocation(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	switch n.Token() {
	case rhino.CALL, rhino.NEW, rhino.OPTCHAIN_CALL, rhino.TAGGED_TEMPLATELIT:
		return true
	}
	return false
}

// IsInvocationTarget reports whether n is the target (first child) of an
// invocation node (CALL, NEW, OPTCHAIN_CALL, or TAGGED_TEMPLATELIT).
// Mirrors NodeUtil.isInvocationTarget.
func (nodeUtil) IsInvocationTarget(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	p := n.Parent()
	return p != nil && NU.IsInvocation(p) && p.FirstChild() == n
}

// ---------------------------------------------------------------------------
// namespace / goog.* helpers
// ---------------------------------------------------------------------------

// IsNamespaceDecl reports whether n looks like a namespace declaration:
// a @const assignment or declaration of an object literal (possibly with
// the `x = x || {}` pattern).  Mirrors NodeUtil.isNamespaceDecl.
func (nodeUtil) IsNamespaceDecl(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	jsdoc := NU.GetBestJSDocInfo(n)
	// If the JSDoc has a non-typedef type annotation it is NOT a namespace.
	// Mirrors: jsdoc != null && !jsdoc.getTypeNodes().isEmpty() && !jsdoc.hasTypedefType()
	if jsdoc != nil && jsdoc.HasType() && jsdoc.TypedefType() == nil {
		return false
	}
	isMarkedConst := false
	if p := n.Parent(); p != nil && p.Token() == rhino.CONST {
		isMarkedConst = true
	}
	if !isMarkedConst && jsdoc != nil && jsdoc.IsConstant() {
		isMarkedConst = true
	}
	// Note: Java also allows namespace declarations in externs without @const.
	// In this port we do not yet track externs provenance on nodes, so we
	// conservatively require @const (or 'const' keyword) for all inputs.
	if !isMarkedConst {
		return false
	}
	// Determine qnameNode and initializer depending on the declaration shape.
	var qnameNode, initializer *rhino.Node
	p := n.Parent()
	switch {
	case p != nil && NU.IsNameDeclaration(p):
		// var/let/const name = <init>
		qnameNode = n
		initializer = n.FirstChild()
	case n.Token() == rhino.EXPR_RESULT:
		expr := n.FirstChild()
		if expr == nil || expr.Token() != rhino.ASSIGN {
			return false
		}
		lhs := expr.FirstChild()
		if lhs == nil || lhs.Token() != rhino.GETPROP {
			return false
		}
		qnameNode = lhs
		initializer = expr.LastChild()
	case n.Token() == rhino.GETPROP:
		pp := n.Parent()
		if pp == nil || pp.Token() != rhino.ASSIGN {
			return false
		}
		gpp := pp.Parent()
		if gpp == nil || gpp.Token() != rhino.EXPR_RESULT {
			return false
		}
		qnameNode = n
		initializer = pp.LastChild()
	default:
		return false
	}
	if initializer == nil || qnameNode == nil {
		return false
	}
	if initializer.Token() == rhino.OBJECTLIT {
		return true
	}
	// `x = x || {}` pattern.
	if initializer.Token() == rhino.OR {
		lhs := initializer.FirstChild()
		rhs := initializer.LastChild()
		if lhs == nil || rhs == nil {
			return false
		}
		return rhs.Token() == rhino.OBJECTLIT &&
			nuQualifiedNameStr(qnameNode) != "" &&
			nuQualifiedNameStr(qnameNode) == nuQualifiedNameStr(lhs)
	}
	return false
}

// nuQualifiedNameStr returns the dotted qualified name string for a NAME or
// GETPROP chain, or "" if n is not a qualified name.
func nuQualifiedNameStr(n *rhino.Node) string {
	if n == nil {
		return ""
	}
	switch n.Token() {
	case rhino.NAME:
		return n.String_()
	case rhino.GETPROP:
		kids := n.Children()
		if len(kids) >= 2 {
			left := nuQualifiedNameStr(kids[0])
			if left == "" {
				return ""
			}
			return left + "." + kids[1].String_()
		}
	}
	return ""
}

// IsFunctionObjectCall reports whether callNode is a `x.call(...)` expression.
// Mirrors NodeUtil.isFunctionObjectCall.
func (nodeUtil) IsFunctionObjectCall(callNode *rhino.Node) bool {
	return nuIsObjectCallMethod(callNode, "call")
}

// IsFunctionObjectApply reports whether callNode is a `x.apply(...)` expression.
// Mirrors NodeUtil.isFunctionObjectApply.
func (nodeUtil) IsFunctionObjectApply(callNode *rhino.Node) bool {
	return nuIsObjectCallMethod(callNode, "apply")
}

// nuIsObjectCallMethod mirrors NodeUtil.isObjectCallMethod: returns true when
// callNode is a CALL/OPTCHAIN_CALL whose callee is a GETPROP/OPTCHAIN_GETPROP
// with the given method name.
func nuIsObjectCallMethod(callNode *rhino.Node, methodName string) bool {
	if callNode == nil {
		return false
	}
	switch callNode.Token() {
	case rhino.CALL, rhino.OPTCHAIN_CALL:
	default:
		return false
	}
	callee := callNode.FirstChild()
	if callee == nil {
		return false
	}
	switch callee.Token() {
	case rhino.GETPROP, rhino.OPTCHAIN_GETPROP:
		return callee.String_() == methodName
	case rhino.GETELEM, rhino.OPTCHAIN_GETELEM:
		last := callee.LastChild()
		return last != nil && last.Token() == rhino.STRINGLIT && last.String_() == methodName
	}
	return false
}

// ---------------------------------------------------------------------------
// hoisting / function-block predicates
// ---------------------------------------------------------------------------

// IsHoistedFunctionDeclaration reports whether n is a function declaration
// that is hoisted to the top of its enclosing function or script scope.
// Mirrors NodeUtil.isHoistedFunctionDeclaration.
func (nodeUtil) IsHoistedFunctionDeclaration(n *rhino.Node) bool {
	if !NU.IsFunctionDeclaration(n) {
		return false
	}
	p := n.Parent()
	if p == nil {
		return false
	}
	switch p.Token() {
	case rhino.SCRIPT, rhino.MODULE_BODY:
		return true
	}
	// Inside a function body (parent of parent is FUNCTION).
	gp := p.Parent()
	if gp != nil && gp.Token() == rhino.FUNCTION {
		return true
	}
	return false
}

// IsFunctionBody reports whether n is the BLOCK node that forms the body of a
// function expression or declaration.  Mirrors NodeUtil.isFunctionBlock.
func (nodeUtil) IsFunctionBody(n *rhino.Node) bool {
	return n != nil && n.Token() == rhino.BLOCK &&
		n.HasParent() && n.Parent().Token() == rhino.FUNCTION
}

// IsFunctionBlock is an alias for IsFunctionBody; the Java source uses both
// names interchangeably (isFunctionBlock is package-private).
// Mirrors NodeUtil.isFunctionBlock.
func (nodeUtil) IsFunctionBlock(n *rhino.Node) bool {
	return NU.IsFunctionBody(n)
}

// ---------------------------------------------------------------------------
// expression predicates
// ---------------------------------------------------------------------------

// IsExprCall reports whether n is an EXPR_RESULT whose sole child is a CALL.
// Mirrors NodeUtil.isExprCall.
func (nodeUtil) IsExprCall(n *rhino.Node) bool {
	if n == nil || n.Token() != rhino.EXPR_RESULT {
		return false
	}
	fc := n.FirstChild()
	return fc != nil && fc.Token() == rhino.CALL
}

// IsCallOrNew reports whether node is a CALL, NEW, or OPTCHAIN_CALL.
// Mirrors NodeUtil.isCallOrNew.
func (nodeUtil) IsCallOrNew(node *rhino.Node) bool {
	if node == nil {
		return false
	}
	switch node.Token() {
	case rhino.CALL, rhino.NEW, rhino.OPTCHAIN_CALL:
		return true
	}
	return false
}

// IsAnyFor reports whether n is any for-loop variant: FOR, FOR_IN, FOR_OF, or
// FOR_AWAIT_OF.  Mirrors NodeUtil.isAnyFor.
func (nodeUtil) IsAnyFor(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	switch n.Token() {
	case rhino.FOR, rhino.FOR_IN, rhino.FOR_OF, rhino.FOR_AWAIT_OF:
		return true
	}
	return false
}

// IsTopLevelFunction reports whether n is a FUNCTION node whose parent is a
// SCRIPT or MODULE_BODY (top-level in the file).
func (nodeUtil) IsTopLevelFunction(n *rhino.Node) bool {
	if n == nil || n.Token() != rhino.FUNCTION {
		return false
	}
	p := n.Parent()
	if p == nil {
		return false
	}
	switch p.Token() {
	case rhino.SCRIPT, rhino.MODULE_BODY:
		return true
	}
	// A function declaration wrapped in an EXPR_RESULT at top level (unusual
	// but guard anyway).
	if p.Token() == rhino.EXPR_RESULT {
		gp := p.Parent()
		if gp != nil {
			switch gp.Token() {
			case rhino.SCRIPT, rhino.MODULE_BODY:
				return true
			}
		}
	}
	return false
}

// IsAnonymousFunction reports whether n is a FUNCTION node with no name (or an
// empty-string name), i.e. a function expression without a declared name.
// Mirrors Java's treatment of anonymous functions in NodeUtil.getNameNode.
func (nodeUtil) IsAnonymousFunction(n *rhino.Node) bool {
	if n == nil || n.Token() != rhino.FUNCTION {
		return false
	}
	nameNode := n.FirstChild()
	if nameNode == nil {
		return true
	}
	// Arrow functions and function expressions may have a NAME node with an
	// empty string as their name component.
	return nameNode.Token() == rhino.NAME && nameNode.String_() == ""
}

// ---------------------------------------------------------------------------
// goog.* call predicates and namespace extractors
// ---------------------------------------------------------------------------

// nuMatchesQualifiedName returns true when n spells out the given dotted name
// as a GETPROP chain or bare NAME.
func nuMatchesQualifiedName(n *rhino.Node, name string) bool {
	return nuQualifiedNameStr(n) == name
}

// IsGoogProvideCall reports whether n is an EXPR_RESULT wrapping a
// goog.provide(…) call.  Mirrors NodeUtil.isGoogProvideCall.
func (nodeUtil) IsGoogProvideCall(n *rhino.Node) bool {
	if !NU.IsExprCall(n) {
		return false
	}
	target := n.FirstChild().FirstChild()
	return nuMatchesQualifiedName(target, "goog.provide")
}

// IsGoogRequireCall reports whether n is a CALL to goog.require(…).
// Mirrors NodeUtil.isGoogRequireCall.
func (nodeUtil) IsGoogRequireCall(n *rhino.Node) bool {
	if n == nil || n.Token() != rhino.CALL {
		return false
	}
	return nuMatchesQualifiedName(n.FirstChild(), "goog.require")
}

// IsGoogModuleCall reports whether n is an EXPR_RESULT wrapping a
// goog.module(…) call.  Mirrors NodeUtil.isGoogModuleCall.
func (nodeUtil) IsGoogModuleCall(n *rhino.Node) bool {
	if !NU.IsExprCall(n) {
		return false
	}
	target := n.FirstChild().FirstChild()
	return nuMatchesQualifiedName(target, "goog.module")
}

// GetGoogProvideName returns the string argument of a goog.provide(…) call,
// or "" if n is not a goog.provide EXPR_RESULT or has no string argument.
func (nodeUtil) GetGoogProvideName(n *rhino.Node) string {
	if !NU.IsGoogProvideCall(n) {
		return ""
	}
	// n is EXPR_RESULT → CALL → [callee, arg0]
	call := n.FirstChild()
	arg := call.GetSecondChild()
	if arg == nil || arg.Token() != rhino.STRINGLIT {
		return ""
	}
	return arg.String_()
}

// GetGoogRequireName returns the string argument of a goog.require(…) call,
// or "" if n is not a goog.require CALL or has no string argument.
func (nodeUtil) GetGoogRequireName(n *rhino.Node) string {
	if !NU.IsGoogRequireCall(n) {
		return ""
	}
	arg := n.GetSecondChild()
	if arg == nil || arg.Token() != rhino.STRINGLIT {
		return ""
	}
	return arg.String_()
}

// ---------------------------------------------------------------------------
// string-value extraction
// ---------------------------------------------------------------------------

// GetStringValue attempts to fold n to a compile-time string.  Returns the
// string and true on success; returns "", false when folding is not possible.
// Mirrors NodeUtil.getStringValue (Java) — template literals are supported
// only when all substitution expressions are themselves foldable.
func (nodeUtil) GetStringValue(n *rhino.Node) (string, bool) {
	if n == nil {
		return "", false
	}
	switch n.Token() {
	case rhino.STRINGLIT, rhino.STRING_KEY:
		return n.String_(), true
	case rhino.TEMPLATELIT:
		// Fold only if every sub-expression is foldable.
		var sb strings.Builder
		for child := n.FirstChild(); child != nil; child = child.Next() {
			expr := child
			if child.Token() == rhino.TEMPLATELIT_SUB {
				expr = child.FirstChild()
			}
			s, ok := NU.GetStringValue(expr)
			if !ok {
				return "", false
			}
			sb.WriteString(s)
		}
		return sb.String(), true
	case rhino.TEMPLATELIT_STRING:
		// The cooked string is stored in String_().
		return n.String_(), true
	case rhino.NAME:
		name := n.String_()
		switch name {
		case "undefined", "Infinity", "NaN":
			return name, true
		}
		return "", false
	case rhino.FALSE:
		return "false", true
	case rhino.TRUE:
		return "true", true
	case rhino.NULL:
		return "null", true
	case rhino.VOID:
		return "undefined", true
	case rhino.NOT:
		fc := n.FirstChild()
		if fc == nil {
			return "", false
		}
		bv := NU.GetBooleanValue(fc)
		if bv == TriVUnknown {
			return "", false
		}
		if bv == TriVTrue {
			return "false", true
		}
		return "true", true
	case rhino.OBJECTLIT:
		return "[object Object]", true
	}
	return "", false
}

// GetPrefixedName returns the string "prefix + name" for a NAME or GETPROP
// node, prepending prefix to the qualified name string.  Returns "" when n is
// not a recognisable name node.  Useful for synthesising mangled names.
func (nodeUtil) GetPrefixedName(n *rhino.Node, prefix string) string {
	qname := nuQualifiedNameStr(n)
	if qname == "" {
		return ""
	}
	return prefix + qname
}

// ---------------------------------------------------------------------------
// node-creation helpers
// ---------------------------------------------------------------------------

// NewName creates a new NAME node carrying the given identifier string.
// Mirrors NodeUtil.newName (the package-private overloads without compiler arg).
func (nodeUtil) NewName(name string) *rhino.Node {
	return rhino.NewString(rhino.NAME, name)
}

// NewDecl creates a VAR/LET/CONST declaration node for lhs, optionally
// initialised with rhs (may be nil).
// Mirrors NodeUtil.newDeclaration.
func (nodeUtil) NewDecl(lhs *rhino.Node, rhs *rhino.Node, declType rhino.Token) *rhino.Node {
	decl := rhino.New(declType)
	if rhs == nil {
		decl.AddChildToBack(lhs)
	} else {
		lhs.AddChildToBack(rhs)
		decl.AddChildToBack(lhs)
	}
	return decl
}
