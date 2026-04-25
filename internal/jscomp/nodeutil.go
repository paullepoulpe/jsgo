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
	"math"
	"strconv"
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

// Go-only: language-adapter — Java exposes getName and getNearestFunctionName
// separately; this consolidates them for the subset callers need.
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

// Go-only: package-bridge — inlines Java's private getQualifiedName logic.
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
	kids := n.Children()
	if len(kids) < 2 {
		return ""
	}
	base := nuQualifiedName(kids[0])
	if base == "" {
		return ""
	}
	// In our AST, GETPROP has the property name as its second child (STRINGLIT).
	// Use its String_() rather than n.String_() (which is empty for GETPROP).
	prop := kids[1].String_()
	if prop == "" {
		return ""
	}
	return base + "." + prop
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

// Go-only: language-adapter — Java's version is on AstAnalyzer, not NodeUtil;
// this provides a simpler inline version without @nosideeffects JSDoc support.
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

// Go-only: language-adapter — Java inlines this predicate; we extract it
// as a named method to avoid repetition across callers.
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

// Go-only: language-adapter — Java's version is on AstAnalyzer, not NodeUtil.
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

// Go-only: language-adapter — Java's newQName takes an AbstractCompiler;
// this version omits the compiler context for callers that don't need it.
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

// Go-only: language-adapter — delegates to Node.removeChildren(); Java callers
// call Node.removeChildren() directly without going through NodeUtil.
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

// Go-only: language-adapter — Java inlines this; Go extracts it for reuse.
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

// Go-only: language-adapter — Java inlines this pattern; Go extracts it.
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

// Go-only: language-adapter — Java inlines this predicate.
func (nodeUtil) IsGetPropOnThis(n *rhino.Node) bool {
	if n == nil || n.Token() != rhino.GETPROP {
		return false
	}
	fc := n.FirstChild()
	return fc != nil && fc.Token() == rhino.THIS
}

// Go-only: language-adapter — Java has isNamedClassExpression but not this variant.
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

// Go-only: language-adapter — Java has no equivalent; used by type inference.
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

// Go-only: package-bridge — Java name is createsScope (package-private); renamed
// with nu prefix to avoid collision with method namespace in Go.
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

// Go-only: language-adapter — Java uses getFunctionParameters().getFirstChild() inline.
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

// Go-only: language-adapter — Java uses visitLhsNodesInNode with a callback;
// Go returns a slice for easier iteration.
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

// Go-only: language-adapter — Java uses visitLhsNodesInNode with a callback;
// Go returns a slice for easier iteration.
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

// Go-only: package-bridge — Java uses inner lambda; Go needs a named helper.
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

// Go-only: language-adapter — alias for GetRValueOfLValue for callers using
// the alternate Java naming variant.
func (nodeUtil) GetRValueFromLValue(n *rhino.Node) *rhino.Node {
	return NU.GetRValueOfLValue(n)
}

// ---------------------------------------------------------------------------
// call / invocation target predicates
// ---------------------------------------------------------------------------

// Go-only: language-adapter — Java inlines this predicate in several passes.
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

// Go-only: language-adapter — extends Java's isInvocationTarget to cover NEW.
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

// Go-only: package-bridge — Java name is matchesQualifiedName or similar; renamed with nu prefix.
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

// Go-only: package-bridge — Java name is isObjectCallMethod (package-private).
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

// Go-only: language-adapter — Java name is isFunctionBlock; renamed to avoid
// ambiguity with the IsFunctionBlock alias below.
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

// Go-only: language-adapter — Java uses isTopLevel which is more general;
// this variant is narrowed to functions for type-checker use.
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

// Go-only: language-adapter — Java has isAnonymousFunction but it is on Node,
// not NodeUtil; this wraps the predicate on the NU helper.
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

// Go-only: package-bridge — Java name is matchesQualifiedName (package-private).
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

// Go-only: language-adapter — Java uses getGoogRequireInfo which returns a
// richer object; this extracts just the name for simple callers.
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

// Go-only: language-adapter — Java uses getGoogRequireInfo; this extracts just the name.
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
	case rhino.VOID, rhino.UNDEFINED:
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
	case rhino.NUMBER:
		return numberToString(n.Number()), true
	case rhino.NEG:
		// Negation of a NUMBER literal: `-1.0` → `"-1"`.
		fc := n.FirstChild()
		if fc != nil && fc.Token() == rhino.NUMBER {
			return numberToString(-fc.Number()), true
		}
		return "", false
	case rhino.ARRAYLIT:
		s, ok := arrayToString(n)
		if !ok {
			return "", false
		}
		return s, true
	}
	return "", false
}

// numberToString converts a float64 to the JavaScript string representation.
// Mirrors com.google.javascript.rhino.jstype.TernaryValue.DToA.numberToString
// for the common cases used by GetStringValue tests.
func numberToString(v float64) string {
	if math.IsNaN(v) {
		return "NaN"
	}
	if math.IsInf(v, 1) {
		return "Infinity"
	}
	if math.IsInf(v, -1) {
		return "-Infinity"
	}
	// Use strconv.FormatFloat with shortest representation.
	// JavaScript's Number.toString() uses scientific notation when the value
	// is >= 1e21 or has many significant digits in the exponent.
	abs := v
	if abs < 0 {
		abs = -abs
	}
	// Use 'e' (scientific) notation for very large / very small numbers,
	// matching JavaScript's Number.toString() cutoff at 1e21.
	if abs >= 1e21 || (abs > 0 && abs < 1e-6) {
		s := strconv.FormatFloat(v, 'e', -1, 64)
		// Go uses e.g. "1.2323919403474454e+21", JS uses same format.
		return s
	}
	// 'f' format gives "10" for 10.0, "1" for 1.0, etc.
	return strconv.FormatFloat(v, 'f', -1, 64)
}

// arrayToString mirrors Java's NodeUtil.arrayToString(Node) — converts an
// ARRAYLIT to its JavaScript string representation (joining elements with ',').
// Returns "", false when any element cannot be converted.
func arrayToString(n *rhino.Node) (string, bool) {
	var sb strings.Builder
	first := true
	for child := n.FirstChild(); child != nil; child = child.Next() {
		if !first {
			sb.WriteByte(',')
		}
		first = false
		// null, undefined, empty → empty string.
		if child.Token() == rhino.NULL || child.Token() == rhino.UNDEFINED ||
			child.Token() == rhino.VOID || child.Token() == rhino.EMPTY {
			continue
		}
		if child.Token() == rhino.NAME && (child.String_() == "undefined") {
			continue
		}
		s, ok := NU.GetStringValue(child)
		if !ok {
			return "", false
		}
		sb.WriteString(s)
	}
	return sb.String(), true
}

// Go-only: language-adapter — Java has no direct equivalent; used for name
// mangling in the type-checker pipeline.
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

// ---------------------------------------------------------------------------
// isDefinedValue / mayBeUndefined
// ---------------------------------------------------------------------------

// IsDefinedValue reports whether the expression is statically known not to
// evaluate to undefined.  Mirrors NodeUtil.isDefinedValue (Java).
func (nodeUtil) IsDefinedValue(value *rhino.Node) bool {
	if value == nil {
		return false
	}
	switch value.Token() {
	case rhino.ASSIGN, rhino.CAST:
		// Only the assigned value matters.
		return NU.IsDefinedValue(value.LastChild())
	case rhino.COALESCE:
		// 'null' is defined, so we can only trust the RHS.
		return NU.IsDefinedValue(value.GetSecondChild())
	case rhino.AND, rhino.OR:
		return NU.IsDefinedValue(value.FirstChild()) && NU.IsDefinedValue(value.LastChild())
	case rhino.HOOK:
		return NU.IsDefinedValue(value.GetSecondChild()) && NU.IsDefinedValue(value.LastChild())
	case rhino.CALL, rhino.OPTCHAIN_CALL,
		rhino.GETELEM, rhino.GETPROP,
		rhino.OPTCHAIN_GETELEM, rhino.OPTCHAIN_GETPROP,
		rhino.THIS, rhino.VOID, rhino.YIELD, rhino.AWAIT:
		return false
	case rhino.DELPROP, rhino.INC, rhino.DEC,
		rhino.CLASS, rhino.FUNCTION, rhino.REGEXP,
		rhino.EMPTY, rhino.ARRAYLIT, rhino.OBJECTLIT,
		rhino.TEMPLATELIT, rhino.STRINGLIT,
		rhino.NUMBER, rhino.NULL, rhino.TRUE, rhino.FALSE, rhino.NEW:
		return true
	case rhino.NAME:
		name := value.String_()
		return name == "Infinity" || name == "NaN"
	default:
		return false
	}
}

// MayBeUndefined mirrors NodeUtil.mayBeUndefined.
func (nodeUtil) MayBeUndefined(n *rhino.Node) bool {
	return !NU.IsDefinedValue(n)
}

// ---------------------------------------------------------------------------
// isSomeCompileTimeConstStringValue
// ---------------------------------------------------------------------------

// IsSomeCompileTimeConstStringValue reports whether the node is a compile-time
// constant string value (string literal, template literal with no substitutions,
// concatenation of such values, or ternary between such values).
// Mirrors NodeUtil.isSomeCompileTimeConstStringValue (Java).
func (nodeUtil) IsSomeCompileTimeConstStringValue(node *rhino.Node) bool {
	if node == nil {
		return false
	}
	if node.IsStringLit() {
		return true
	}
	if node.Token() == rhino.TEMPLATELIT && node.HasOneChild() {
		return true
	}
	if node.Token() == rhino.ADD {
		left := node.FirstChild()
		right := node.LastChild()
		return left != nil && right != nil &&
			NU.IsSomeCompileTimeConstStringValue(left) &&
			NU.IsSomeCompileTimeConstStringValue(right)
	}
	if node.Token() == rhino.HOOK {
		left := node.GetSecondChild()
		right := node.LastChild()
		return left != nil && right != nil &&
			NU.IsSomeCompileTimeConstStringValue(left) &&
			NU.IsSomeCompileTimeConstStringValue(right)
	}
	return false
}

// ---------------------------------------------------------------------------
// isDeclarationParent / isNamedClass / isNamedFunction / isReferenceName
// ---------------------------------------------------------------------------

// isDeclarationParent mirrors the private NodeUtil.isDeclarationParent.
func isDeclarationParent(parent *rhino.Node) bool {
	if parent == nil {
		return false
	}
	switch parent.Token() {
	case rhino.DECLARE, rhino.EXPORT:
		return true
	}
	return isStatementParent(parent)
}

// isReferenceName mirrors NodeUtil.isReferenceName.
func isReferenceName(n *rhino.Node) bool {
	return n != nil && n.Token() == rhino.NAME && n.String_() != ""
}

// isNamedFunction mirrors NodeUtil.isNamedFunction.
func isNamedFunction(n *rhino.Node) bool {
	return n != nil && n.IsFunction() && isReferenceName(n.FirstChild())
}

// isNamedClass mirrors NodeUtil.isNamedClass.
func isNamedClass(n *rhino.Node) bool {
	return n != nil && n.IsClass() && isReferenceName(n.FirstChild())
}

// isSwitchCase mirrors NodeUtil.isSwitchCase.
func isSwitchCase(n *rhino.Node) bool {
	return n != nil && (n.IsCase() || n.IsDefaultCase())
}

// isExportFrom mirrors NodeUtil.isExportFrom — an EXPORT node that has a from-clause (2 children).
func isExportFrom(n *rhino.Node) bool {
	return n != nil && n.IsExport() && n.HasTwoChildren()
}

// hasCatchHandler mirrors NodeUtil.hasCatchHandler.
func hasCatchHandler(catchBlock *rhino.Node) bool {
	return catchBlock != nil && catchBlock.HasChildren()
}

// hasFinally mirrors NodeUtil.hasFinally.
func hasFinally(n *rhino.Node) bool {
	return n != nil && n.IsTry() && n.HasXChildren(3)
}

// getCatchBlock returns the second child of a TRY node (the catch container block).
func getCatchBlock(tryNode *rhino.Node) *rhino.Node {
	return tryNode.GetSecondChild()
}

// isTryFinallyNode mirrors NodeUtil.isTryFinallyNode.
func isTryFinallyNode(parent, child *rhino.Node) bool {
	return parent != nil && parent.IsTry() && parent.HasXChildren(3) && child == parent.LastChild()
}

// isTryCatchNodeContainer mirrors NodeUtil.isTryCatchNodeContainer.
func isTryCatchNodeContainer(n *rhino.Node) bool {
	parent := n.Parent()
	return parent != nil && parent.IsTry() && parent.GetSecondChild() == n
}

// ---------------------------------------------------------------------------
// isMethodDeclaration / isClassDeclaration / isBlockScopedFunctionDeclaration
// ---------------------------------------------------------------------------

// IsMethodDeclaration reports whether n is a method in a class or object
// literal (getter, setter, computed property method, or member function).
// Mirrors NodeUtil.isMethodDeclaration (Java).
func (nodeUtil) IsMethodDeclaration(n *rhino.Node) bool {
	if n == nil || !n.IsFunction() {
		return false
	}
	p := n.Parent()
	if p == nil {
		return false
	}
	switch p.Token() {
	case rhino.GETTER_DEF, rhino.SETTER_DEF, rhino.MEMBER_FUNCTION_DEF:
		return true
	case rhino.COMPUTED_PROP:
		// Java also checks COMPUTED_PROP_METHOD/GETTER/SETTER boolean props;
		// those props are not yet ported to the Go rhino package, so we
		// conservatively check only whether this function is the last child
		// (the method body position in a computed method).
		return p.LastChild() == n
	}
	return false
}

// IsClassDeclaration reports whether n is a class declaration.
// Mirrors NodeUtil.isClassDeclaration (Java).
func (nodeUtil) IsClassDeclaration(n *rhino.Node) bool {
	return n != nil && n.IsClass() && isDeclarationParent(n.Parent()) && isNamedClass(n)
}

// IsClassExpression reports whether n is a class used within an expression.
// Mirrors NodeUtil.isClassExpression (Java).
func (nodeUtil) IsClassExpression(n *rhino.Node) bool {
	return n != nil && n.IsClass() && (!isNamedClass(n) || !isDeclarationParent(n.Parent()))
}

// IsBlockScopedFunctionDeclaration reports whether n is a function declaration
// inside a block scope (not at the top of a function/script/module body).
// Mirrors NodeUtil.isBlockScopedFunctionDeclaration (Java).
func (nodeUtil) IsBlockScopedFunctionDeclaration(n *rhino.Node) bool {
	if !NU.IsFunctionDeclaration(n) {
		return false
	}
	current := n.Parent()
	for current != nil {
		switch current.Token() {
		case rhino.BLOCK:
			p := current.Parent()
			return p == nil || !p.IsFunction()
		case rhino.FUNCTION, rhino.SCRIPT, rhino.DECLARE, rhino.EXPORT, rhino.MODULE_BODY:
			return false
		default:
			// Must be LABEL
			current = current.Parent()
		}
	}
	return false
}

// ---------------------------------------------------------------------------
// isControlStructureCodeBlock
// ---------------------------------------------------------------------------

// IsControlStructureCodeBlock reports whether n is the code block of its
// control-structure parent.  Mirrors NodeUtil.isControlStructureCodeBlock (Java).
func (nodeUtil) IsControlStructureCodeBlock(parent, n *rhino.Node) bool {
	if parent == nil || n == nil {
		return false
	}
	switch parent.Token() {
	case rhino.DO_WHILE:
		return parent.FirstChild() == n
	case rhino.TRY:
		return parent.FirstChild() == n || parent.LastChild() == n
	case rhino.FOR, rhino.FOR_IN, rhino.FOR_OF, rhino.FOR_AWAIT_OF,
		rhino.WHILE, rhino.LABEL, rhino.CATCH:
		return parent.LastChild() == n
	case rhino.IF, rhino.SWITCH, rhino.CASE:
		return parent.FirstChild() != n
	case rhino.DEFAULT_CASE:
		return true
	}
	return false
}

// ---------------------------------------------------------------------------
// isNonlocalModuleExportName
// ---------------------------------------------------------------------------

// IsNonlocalModuleExportName reports whether the NAME node n is the non-local
// part of an import or export specifier.
// Mirrors NodeUtil.isNonlocalModuleExportName (Java).
func (nodeUtil) IsNonlocalModuleExportName(n *rhino.Node) bool {
	if n == nil || n.Token() != rhino.NAME {
		return false
	}
	parent := n.Parent()
	if parent == nil {
		return false
	}
	if parent.IsImportSpec() && n.IsFirstChildOf(parent) {
		// import {nonlocal as x} from './foo.js'
		return true
	}
	if parent.IsExportSpec() {
		if n.IsFirstChildOf(parent) {
			// export {nonlocal as b} from './foo.js';
			gp := parent.GetGrandparent()
			return gp != nil && isExportFrom(gp)
		}
		// export {local as nonlocal};
		return true
	}
	return false
}

// ---------------------------------------------------------------------------
// getNodeTypeReferenceCount / isNameReferenced / getNameReferenceCount
// ---------------------------------------------------------------------------

// GetNodeTypeReferenceCount returns the count of nodes with the given token
// reachable from node, with traversal limited by traverseChildrenPred.
// Mirrors NodeUtil.getNodeTypeReferenceCount (Java).
func (nodeUtil) GetNodeTypeReferenceCount(node *rhino.Node, typ rhino.Token, traverseChildrenPred func(*rhino.Node) bool) int {
	return nuCount(node, func(n *rhino.Node) bool { return n.Token() == typ }, traverseChildrenPred)
}

// IsNameReferenced reports whether the simple name `name` appears anywhere
// reachable from node (with full traversal, no predicate filter).
// Mirrors NodeUtil.isNameReferenced(Node, String) (Java).
func (nodeUtil) IsNameReferenced(node *rhino.Node, name string) bool {
	return nuHas(node, func(n *rhino.Node) bool {
		return n.Token() == rhino.NAME && n.String_() == name
	}, func(_ *rhino.Node) bool { return true })
}

// GetNameReferenceCount returns the number of NAME references to `name` in the
// subtree rooted at node.  Mirrors NodeUtil.getNameReferenceCount (Java).
func (nodeUtil) GetNameReferenceCount(node *rhino.Node, name string) int {
	return nuCount(node, func(n *rhino.Node) bool {
		return n.Token() == rhino.NAME && n.String_() == name
	}, func(_ *rhino.Node) bool { return true })
}

// Go-only: package-bridge — Java name is has (package-private); renamed to nuHas.
func nuHas(node *rhino.Node, pred func(*rhino.Node) bool, traverseChildrenPred func(*rhino.Node) bool) bool {
	if node == nil {
		return false
	}
	if pred(node) {
		return true
	}
	if !traverseChildrenPred(node) {
		return false
	}
	for c := node.FirstChild(); c != nil; c = c.Next() {
		if nuHas(c, pred, traverseChildrenPred) {
			return true
		}
	}
	return false
}

// Go-only: package-bridge — Java name is getCount (package-private); renamed.
func nuCount(node *rhino.Node, pred func(*rhino.Node) bool, traverseChildrenPred func(*rhino.Node) bool) int {
	if node == nil {
		return 0
	}
	count := 0
	if pred(node) {
		count++
	}
	if traverseChildrenPred(node) {
		for c := node.FirstChild(); c != nil; c = c.Next() {
			count += nuCount(c, pred, traverseChildrenPred)
		}
	}
	return count
}

// ---------------------------------------------------------------------------
// getVarsDeclaredInBranch
// ---------------------------------------------------------------------------

// GetVarsDeclaredInBranch collects all VAR-declared NAME nodes in the subtree
// rooted at root, excluding descent into nested functions (which have their
// own scope).  Mirrors NodeUtil.getVarsDeclaredInBranch (Java).
func (nodeUtil) GetVarsDeclaredInBranch(root *rhino.Node) []*rhino.Node {
	var vars []*rhino.Node
	var collect func(n *rhino.Node)
	collect = func(n *rhino.Node) {
		if n == nil {
			return
		}
		if n != root && n.IsFunction() {
			// Don't descend into nested functions.
			return
		}
		if n.IsVar() {
			// Collect all NAME children of VAR.
			for c := n.FirstChild(); c != nil; c = c.Next() {
				if c.Token() == rhino.NAME {
					vars = append(vars, c)
				}
			}
		}
		for c := n.FirstChild(); c != nil; c = c.Next() {
			collect(c)
		}
	}
	collect(root)
	return vars
}

// ---------------------------------------------------------------------------
// removeChild / replaceDeclarationChild / tryMergeBlock / canMergeBlock
// ---------------------------------------------------------------------------

// RemoveChild safely removes node from parent, maintaining a valid AST.
// Mirrors NodeUtil.removeChild (Java).
func (nodeUtil) RemoveChild(parent, node *rhino.Node) {
	if isTryFinallyNode(parent, node) {
		if hasCatchHandler(getCatchBlock(parent)) {
			node.Detach()
		} else {
			node.DetachChildren()
		}
	} else if node.Token() == rhino.CATCH {
		// CATCH can only be removed if there is a finally clause.
		node.Detach()
	} else if isTryCatchNodeContainer(node) {
		// The container node itself can't be removed, but CATCH can if there's a finally.
		node.DetachChildren()
	} else if node.IsBlock() {
		// Simply empty the block.
		node.DetachChildren()
	} else if NU.IsStatementBlock(parent) || isSwitchCase(node) || node.IsMemberFunctionDef() {
		// A statement in a block or a member function can simply be removed.
		node.Detach()
	} else if NU.IsNameDeclaration(parent) || parent.IsExprResult() {
		if parent.HasMoreThanOneChild() {
			node.Detach()
		} else {
			node.Detach()
			NU.RemoveChild(parent.Parent(), parent)
		}
	} else if parent.IsLabel() && node == parent.LastChild() {
		node.Detach()
		NU.RemoveChild(parent.Parent(), parent)
	} else if parent.IsVanillaFor() {
		node.ReplaceWith(rhino.New(rhino.EMPTY))
	} else if parent.IsObjectPattern() {
		node.Detach()
	} else if parent.IsArrayPattern() {
		if node == parent.LastChild() {
			node.Detach()
		} else {
			node.ReplaceWith(rhino.New(rhino.EMPTY))
		}
	} else if parent.IsDestructuringLhs() {
		node.Detach()
		if parent.Parent() != nil && parent.Parent().HasChildren() {
			NU.RemoveChild(parent.Parent(), parent)
		}
	} else if parent.IsRest() {
		parent.Detach()
	} else if parent.IsParamList() {
		node.Detach()
	} else if parent.IsImport() {
		if node == parent.FirstChild() {
			node.ReplaceWith(rhino.New(rhino.EMPTY))
		}
		// else: invalid removal, silently ignore in Go port
	}
	// else: invalid removal, silently ignore in Go port
}

// ReplaceDeclarationChild replaces the child of a var/let/const declaration
// with a new statement, preserving side-effects for other declaration children.
// Mirrors NodeUtil.replaceDeclarationChild (Java).
func (nodeUtil) ReplaceDeclarationChild(declChild, newStatement *rhino.Node) {
	decl := declChild.Parent()
	if decl == nil {
		return
	}
	if decl.HasOneChild() {
		decl.ReplaceWith(newStatement)
	} else if declChild.Next() == nil {
		declChild.Detach()
		newStatement.InsertAfter(decl)
	} else if declChild.Previous() == nil {
		declChild.Detach()
		newStatement.InsertBefore(decl)
	} else {
		// Split the declaration: keep previous children in decl, move
		// subsequent children to a new declaration node.
		newDecl := rhino.New(decl.Token())
		newDecl.Srcref(decl)
		for after := declChild.Next(); after != nil; {
			next := after.Next()
			after.Detach()
			newDecl.AddChildToBack(after)
			after = next
		}
		declChild.Detach()
		newStatement.InsertAfter(decl)
		newDecl.InsertAfter(newStatement)
	}
}

// ---------------------------------------------------------------------------
// Ported from Java NodeUtil: isUndefined, isNullOrUndefined, isEmptyBlock,
// isAssignmentOp, isLogicalAssignmentOp, isCompoundAssignmentOp,
// isHoistScopeRoot, isEnhancedFor, isLoopStructure, getLoopCodeBlock,
// getEnclosingBlock, getEnclosingBlockScopeRoot, getEnclosingHoistScopeRoot,
// getEnclosingClass, getEnclosingStatement, getEnclosingNode,
// isStatementParent, isNormalGet, isNormalOrOptChainGet, isNormalOrOptChainGetProp,
// isNormalOrOptChainCall, isOptChainGet, isOptChainNode,
// isNumericResult, isBooleanResult, isStringResult, isObjectResult, isBigIntResult,
// isLValue, isLhsOfAssign, isImportedName, getClassMembers,
// getEs6ClassConstructorMemberFunctionDef, getName, getNameNode,
// isEs6Constructor, isEs6ConstructorMemberFunctionDef.
// ---------------------------------------------------------------------------

// IsUndefined reports whether n is an undefined value expression
// (VOID node, or NAME "undefined").
// Mirrors NodeUtil.isUndefined.
func (nodeUtil) IsUndefined(n *rhino.Node) bool {
	switch n.Token() {
	case rhino.VOID:
		return true
	case rhino.NAME:
		return n.GetString() == "undefined"
	}
	return false
}

// IsNullOrUndefined reports whether n is a null or undefined expression.
// Mirrors NodeUtil.isNullOrUndefined.
func (nodeUtil) IsNullOrUndefined(n *rhino.Node) bool {
	return n.IsNull() || NU.IsUndefined(n)
}

// IsEmptyBlock reports whether block is a BLOCK node containing only EMPTY
// nodes (or no children at all).
// Mirrors NodeUtil.isEmptyBlock.
func (nodeUtil) IsEmptyBlock(block *rhino.Node) bool {
	if !block.IsBlock() {
		return false
	}
	for c := block.FirstChild(); c != nil; c = c.Next() {
		if !c.IsEmpty() {
			return false
		}
	}
	return true
}

// IsAssignmentOp reports whether n is any assignment operator
// (=, +=, -=, *=, /=, %=, **=, |=, ^=, &=, <<=, >>=, >>>=, ||=, &&=, ??=).
// Mirrors NodeUtil.isAssignmentOp.
func (nodeUtil) IsAssignmentOp(n *rhino.Node) bool {
	switch n.Token() {
	case rhino.ASSIGN, rhino.ASSIGN_BITOR, rhino.ASSIGN_BITXOR, rhino.ASSIGN_BITAND,
		rhino.ASSIGN_LSH, rhino.ASSIGN_RSH, rhino.ASSIGN_URSH,
		rhino.ASSIGN_ADD, rhino.ASSIGN_SUB, rhino.ASSIGN_MUL,
		rhino.ASSIGN_EXPONENT, rhino.ASSIGN_DIV, rhino.ASSIGN_MOD,
		rhino.ASSIGN_OR, rhino.ASSIGN_AND, rhino.ASSIGN_COALESCE:
		return true
	}
	return false
}

// IsLogicalAssignmentOp reports whether n is a logical assignment operator
// (||=, &&=, ??=).
// Mirrors NodeUtil.isLogicalAssignmentOp.
func (nodeUtil) IsLogicalAssignmentOp(n *rhino.Node) bool {
	switch n.Token() {
	case rhino.ASSIGN_OR, rhino.ASSIGN_AND, rhino.ASSIGN_COALESCE:
		return true
	}
	return false
}

// IsCompoundAssignmentOp reports whether n is a compound assignment operator
// (any assignment op except plain =).
// Mirrors NodeUtil.isCompoundAssignmentOp.
func (nodeUtil) IsCompoundAssignmentOp(n *rhino.Node) bool {
	return NU.IsAssignmentOp(n) && !n.IsAssign()
}

// IsHoistScopeRoot reports whether n is a hoist-scope root (FUNCTION, MODULE_BODY,
// or class static block).
// Mirrors NodeUtil.isHoistScopeRoot.
func (nodeUtil) IsHoistScopeRoot(n *rhino.Node) bool {
	return n.IsFunction() || n.IsModuleBody() || nodeIsClassStaticBlock(n)
}

// IsEnhancedFor reports whether n is a for-of, for-await-of, or for-in loop.
// Mirrors NodeUtil.isEnhancedFor.
func (nodeUtil) IsEnhancedFor(n *rhino.Node) bool {
	return n.IsForOf() || n.IsForAwaitOf() || n.IsForIn()
}

// IsLoopStructure reports whether n is a loop (FOR, FOR_IN, FOR_OF,
// FOR_AWAIT_OF, DO, WHILE).
// Mirrors NodeUtil.isLoopStructure.
func (nodeUtil) IsLoopStructure(n *rhino.Node) bool {
	switch n.Token() {
	case rhino.FOR, rhino.FOR_IN, rhino.FOR_OF, rhino.FOR_AWAIT_OF, rhino.DO, rhino.WHILE:
		return true
	}
	return false
}

// GetLoopCodeBlock returns the body BLOCK of a loop node, or nil for non-loops.
// Mirrors NodeUtil.getLoopCodeBlock.
func (nodeUtil) GetLoopCodeBlock(n *rhino.Node) *rhino.Node {
	switch n.Token() {
	case rhino.FOR, rhino.FOR_IN, rhino.FOR_OF, rhino.FOR_AWAIT_OF, rhino.WHILE:
		return n.LastChild()
	case rhino.DO:
		return n.FirstChild()
	}
	return nil
}

// GetEnclosingBlock returns the nearest enclosing BLOCK node, or nil.
// Mirrors NodeUtil.getEnclosingBlock.
func (nodeUtil) GetEnclosingBlock(n *rhino.Node) *rhino.Node {
	return NU.GetEnclosingType(n, rhino.BLOCK)
}

// Go-only: package-bridge — Java name is createsBlockScope (package-private).
func nuCreatesBlockScope(n *rhino.Node) bool {
	switch n.Token() {
	case rhino.FOR, rhino.FOR_IN, rhino.FOR_OF, rhino.FOR_AWAIT_OF, rhino.SWITCH_BODY:
		return true
	case rhino.BLOCK:
		p := n.Parent()
		if p == nil {
			return false
		}
		switch p.Token() {
		case rhino.FUNCTION, rhino.SCRIPT, rhino.MODULE_BODY:
			return false
		}
		return true
	}
	return false
}

// GetEnclosingBlockScopeRoot returns the nearest block-scope root, or nil.
// Mirrors NodeUtil.getEnclosingBlockScopeRoot.
func (nodeUtil) GetEnclosingBlockScopeRoot(n *rhino.Node) *rhino.Node {
	curr := n.Parent()
	for curr != nil {
		if nuCreatesBlockScope(curr) {
			return curr
		}
		curr = curr.Parent()
	}
	return nil
}

// GetEnclosingHoistScopeRoot returns the nearest hoist-scope root, or nil.
// Mirrors NodeUtil.getEnclosingHoistScopeRoot.
func (nodeUtil) GetEnclosingHoistScopeRoot(n *rhino.Node) *rhino.Node {
	curr := n.Parent()
	for curr != nil {
		if NU.IsHoistScopeRoot(curr) {
			return curr
		}
		curr = curr.Parent()
	}
	return nil
}

// GetEnclosingClass returns the nearest enclosing CLASS node, or nil.
// Mirrors NodeUtil.getEnclosingClass.
func (nodeUtil) GetEnclosingClass(n *rhino.Node) *rhino.Node {
	return NU.GetEnclosingType(n, rhino.CLASS)
}

// GetEnclosingStatement returns the nearest ancestor that is a statement, or nil.
// Mirrors NodeUtil.getEnclosingStatement.
func (nodeUtil) GetEnclosingStatement(n *rhino.Node) *rhino.Node {
	curr := n
	for curr != nil {
		if NU.IsStatement(curr) {
			return curr
		}
		curr = curr.Parent()
	}
	return nil
}

// GetEnclosingNode returns the nearest ancestor (including n itself) for which
// pred returns true, or nil.
// Mirrors NodeUtil.getEnclosingNode.
func (nodeUtil) GetEnclosingNode(n *rhino.Node, pred func(*rhino.Node) bool) *rhino.Node {
	curr := n
	for curr != nil {
		if pred(curr) {
			return curr
		}
		curr = curr.Parent()
	}
	return nil
}

// GetEnclosingModuleIfPresent returns the nearest enclosing MODULE_BODY, or nil.
// Mirrors NodeUtil.getEnclosingModuleIfPresent.
func (nodeUtil) GetEnclosingModuleIfPresent(n *rhino.Node) *rhino.Node {
	return NU.GetEnclosingType(n, rhino.MODULE_BODY)
}

// IsNormalGet reports whether n is a GETPROP or GETELEM node.
// Mirrors NodeUtil.isNormalGet.
func (nodeUtil) IsNormalGet(n *rhino.Node) bool {
	return n.IsGetProp() || n.IsGetElem()
}

// IsOptChainGet reports whether n is an OPTCHAIN_GETPROP or OPTCHAIN_GETELEM.
// Mirrors NodeUtil.isOptChainGet.
func (nodeUtil) IsOptChainGet(n *rhino.Node) bool {
	return n.IsOptChainGetProp() || n.IsOptChainGetElem()
}

// IsOptChainNode reports whether n is any optional-chain node
// (OPTCHAIN_GETPROP, OPTCHAIN_GETELEM, or OPTCHAIN_CALL).
// Mirrors NodeUtil.isOptChainNode.
func (nodeUtil) IsOptChainNode(n *rhino.Node) bool {
	return n.IsOptChainGetProp() || n.IsOptChainGetElem() || n.IsOptChainCall()
}

// IsNormalOrOptChainGet reports whether n is any get node
// (GETPROP, GETELEM, OPTCHAIN_GETPROP, or OPTCHAIN_GETELEM).
// Mirrors NodeUtil.isNormalOrOptChainGet.
func (nodeUtil) IsNormalOrOptChainGet(n *rhino.Node) bool {
	return NU.IsNormalGet(n) || NU.IsOptChainGet(n)
}

// IsNormalOrOptChainGetProp reports whether n is a GETPROP or OPTCHAIN_GETPROP.
// Mirrors NodeUtil.isNormalOrOptChainGetProp.
func (nodeUtil) IsNormalOrOptChainGetProp(n *rhino.Node) bool {
	return n.IsGetProp() || n.IsOptChainGetProp()
}

// IsNormalOrOptChainCall reports whether n is a CALL or OPTCHAIN_CALL.
// Mirrors NodeUtil.isNormalOrOptChainCall.
func (nodeUtil) IsNormalOrOptChainCall(n *rhino.Node) bool {
	return n.IsCall() || n.IsOptChainCall()
}

// ---------------------------------------------------------------------------
// Value type predicates (depend on GetKnownValueType — ported below).
// ---------------------------------------------------------------------------

// ValueType mirrors com.google.javascript.jscomp.NodeUtil.ValueType.
type ValueType int

const (
	ValueTypeUndetermined ValueType = iota
	ValueTypeNull
	ValueTypeVoid
	ValueTypeNumber
	ValueTypeBigint
	ValueTypeString
	ValueTypeBoolean
	ValueTypeObject
)

// Go-only: package-bridge — Java's getKnownValueType is directly on NodeUtil;
// this is an unexported wrapper to break the call cycle in Go.
func getKnownValueTypePriv(n *rhino.Node) ValueType {
	switch n.Token() {
	case rhino.CAST:
		return getKnownValueTypePriv(n.FirstChild())
	case rhino.ASSIGN, rhino.COMMA:
		return getKnownValueTypePriv(n.LastChild())
	case rhino.AND, rhino.OR, rhino.COALESCE,
		rhino.ASSIGN_OR, rhino.ASSIGN_AND, rhino.ASSIGN_COALESCE:
		a := getKnownValueTypePriv(n.FirstChild())
		b := getKnownValueTypePriv(n.LastChild())
		if a == b {
			return a
		}
		return ValueTypeUndetermined
	case rhino.HOOK:
		a := getKnownValueTypePriv(n.GetSecondChild())
		b := getKnownValueTypePriv(n.LastChild())
		if a == b {
			return a
		}
		return ValueTypeUndetermined
	case rhino.ADD:
		last := getKnownValueTypePriv(n.LastChild())
		if last == ValueTypeString {
			return ValueTypeString
		}
		first := getKnownValueTypePriv(n.FirstChild())
		if first == ValueTypeString {
			return ValueTypeString
		}
		if first == ValueTypeObject || last == ValueTypeObject {
			return ValueTypeUndetermined
		}
		if nuMayBeString(first) || nuMayBeString(last) {
			return ValueTypeUndetermined
		}
		if first == ValueTypeBigint || last == ValueTypeBigint {
			return ValueTypeBigint
		}
		return ValueTypeNumber
	case rhino.ASSIGN_ADD:
		last := getKnownValueTypePriv(n.LastChild())
		if last == ValueTypeString {
			return ValueTypeString
		}
		return ValueTypeUndetermined
	case rhino.NAME:
		s := n.GetString()
		if s == "undefined" {
			return ValueTypeVoid
		}
		if s == "NaN" || s == "Infinity" {
			return ValueTypeNumber
		}
		return ValueTypeUndetermined
	case rhino.ASSIGN_BITOR, rhino.ASSIGN_BITXOR, rhino.ASSIGN_BITAND,
		rhino.ASSIGN_LSH, rhino.ASSIGN_RSH, rhino.ASSIGN_URSH,
		rhino.ASSIGN_SUB, rhino.ASSIGN_MUL, rhino.ASSIGN_EXPONENT,
		rhino.ASSIGN_DIV, rhino.ASSIGN_MOD:
		if getKnownValueTypePriv(n.LastChild()) == ValueTypeBigint {
			return ValueTypeBigint
		}
		return ValueTypeNumber
	case rhino.BIGINT:
		return ValueTypeBigint
	case rhino.BITOR, rhino.BITXOR, rhino.BITAND,
		rhino.LSH, rhino.RSH, rhino.SUB, rhino.MUL, rhino.MOD, rhino.DIV, rhino.EXPONENT:
		first := getKnownValueTypePriv(n.FirstChild())
		last := getKnownValueTypePriv(n.LastChild())
		if first == ValueTypeBigint || last == ValueTypeBigint {
			return ValueTypeBigint
		}
		return ValueTypeNumber
	case rhino.BITNOT, rhino.NEG:
		if getKnownValueTypePriv(n.FirstChild()) == ValueTypeBigint {
			return ValueTypeBigint
		}
		return ValueTypeNumber
	case rhino.INC, rhino.DEC, rhino.URSH, rhino.POS, rhino.NUMBER:
		return ValueTypeNumber
	case rhino.TRUE, rhino.FALSE,
		rhino.EQ, rhino.NE, rhino.SHEQ, rhino.SHNE,
		rhino.LT, rhino.GT, rhino.LE, rhino.GE,
		rhino.IN, rhino.INSTANCEOF,
		rhino.NOT, rhino.DELPROP:
		return ValueTypeBoolean
	case rhino.TYPEOF, rhino.STRINGLIT, rhino.TEMPLATELIT:
		return ValueTypeString
	case rhino.NULL:
		return ValueTypeNull
	case rhino.VOID:
		return ValueTypeVoid
	case rhino.FUNCTION, rhino.NEW, rhino.ARRAYLIT, rhino.OBJECTLIT, rhino.REGEXP:
		return ValueTypeObject
	default:
		return ValueTypeUndetermined
	}
}

// Go-only: package-bridge — Java name is mayBeString(ValueType) (package-private).
func nuMayBeString(t ValueType) bool {
	switch t {
	case ValueTypeBoolean, ValueTypeNull, ValueTypeNumber, ValueTypeBigint, ValueTypeVoid:
		return false
	}
	return true
}

// GetKnownValueType returns the known primitive value type of n, or
// ValueTypeUndetermined when it cannot be determined statically.
// Mirrors NodeUtil.getKnownValueType.
func (nodeUtil) GetKnownValueType(n *rhino.Node) ValueType {
	return getKnownValueTypePriv(n)
}

// IsNumericResult reports whether n always evaluates to a number.
// Mirrors NodeUtil.isNumericResult.
func (nodeUtil) IsNumericResult(n *rhino.Node) bool {
	return getKnownValueTypePriv(n) == ValueTypeNumber
}

// IsBigIntResult reports whether n always evaluates to a bigint.
// Mirrors NodeUtil.isBigIntResult.
func (nodeUtil) IsBigIntResult(n *rhino.Node) bool {
	return getKnownValueTypePriv(n) == ValueTypeBigint
}

// IsBooleanResult reports whether n always evaluates to a boolean.
// Mirrors NodeUtil.isBooleanResult.
func (nodeUtil) IsBooleanResult(n *rhino.Node) bool {
	return getKnownValueTypePriv(n) == ValueTypeBoolean
}

// IsStringResult reports whether n always evaluates to a string.
// Mirrors NodeUtil.isStringResult.
func (nodeUtil) IsStringResult(n *rhino.Node) bool {
	return getKnownValueTypePriv(n) == ValueTypeString
}

// IsObjectResult reports whether n always evaluates to an object.
// Mirrors NodeUtil.isObjectResult.
func (nodeUtil) IsObjectResult(n *rhino.Node) bool {
	return getKnownValueTypePriv(n) == ValueTypeObject
}

// ---------------------------------------------------------------------------
// LValue predicates
// ---------------------------------------------------------------------------

// IsLValue reports whether n is used as an L-value (assignment target or
// declaration binding).
// Mirrors NodeUtil.isLValue.
func (nodeUtil) IsLValue(n *rhino.Node) bool {
	switch n.Token() {
	case rhino.NAME, rhino.GETPROP, rhino.GETELEM:
		// fall through to parent check
	default:
		return false
	}
	parent := n.Parent()
	if parent == nil {
		return false
	}
	switch parent.Token() {
	case rhino.IMPORT_SPEC:
		return parent.LastChild() == n
	case rhino.VAR, rhino.LET, rhino.CONST,
		rhino.ITER_REST, rhino.OBJECT_REST,
		rhino.PARAM_LIST, rhino.IMPORT,
		rhino.INC, rhino.DEC, rhino.CATCH:
		return true
	case rhino.CLASS, rhino.FUNCTION, rhino.DEFAULT_VALUE,
		rhino.FOR, rhino.FOR_IN, rhino.FOR_OF, rhino.FOR_AWAIT_OF:
		return parent.FirstChild() == n
	case rhino.ARRAY_PATTERN, rhino.STRING_KEY, rhino.COMPUTED_PROP:
		return NU.IsLhsByDestructuring(n)
	default:
		return NU.IsAssignmentOp(parent) && parent.FirstChild() == n
	}
}

// IsLhsOfAssign reports whether n is the left-hand side of a plain assignment.
// Mirrors NodeUtil.isLhsOfAssign.
func (nodeUtil) IsLhsOfAssign(n *rhino.Node) bool {
	parent := n.Parent()
	return parent != nil && parent.IsAssign() && parent.FirstChild() == n
}

// IsImportedName reports whether n is a NAME node that was introduced by an
// import declaration.
// Mirrors NodeUtil.isImportedName.
func (nodeUtil) IsImportedName(n *rhino.Node) bool {
	parent := n.Parent()
	if parent == nil {
		return false
	}
	return parent.IsImport() || (parent.IsImportSpec() && parent.LastChild() == n)
}

// ---------------------------------------------------------------------------
// Class helpers
// ---------------------------------------------------------------------------

// GetClassMembers returns the CLASS_MEMBERS child of a CLASS node.
// Mirrors NodeUtil.getClassMembers.
func (nodeUtil) GetClassMembers(n *rhino.Node) *rhino.Node {
	return n.LastChild()
}

// IsEs6ConstructorMemberFunctionDef reports whether n is a MEMBER_FUNCTION_DEF
// whose string value is "constructor" and that is not static.
// Mirrors NodeUtil.isEs6ConstructorMemberFunctionDef.
func (nodeUtil) IsEs6ConstructorMemberFunctionDef(n *rhino.Node) bool {
	return n.IsMemberFunctionDef() && !n.IsStaticMember() && n.GetString() == "constructor"
}

// IsEs6Constructor reports whether n is the FUNCTION node of an ES6 class
// constructor.
// Mirrors NodeUtil.isEs6Constructor.
func (nodeUtil) IsEs6Constructor(n *rhino.Node) bool {
	return n.IsFunction() && n.Parent() != nil &&
		NU.IsEs6ConstructorMemberFunctionDef(n.Parent())
}

// GetEs6ClassConstructorMemberFunctionDef returns the MEMBER_FUNCTION_DEF for
// the constructor of a class node, or nil if none is present.
// Mirrors NodeUtil.getEs6ClassConstructorMemberFunctionDef.
func (nodeUtil) GetEs6ClassConstructorMemberFunctionDef(classNode *rhino.Node) *rhino.Node {
	classMembers := classNode.LastChild()
	if classMembers == nil {
		return nil
	}
	for m := classMembers.FirstChild(); m != nil; m = m.Next() {
		if NU.IsEs6ConstructorMemberFunctionDef(m) {
			return m
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// getName / getNameNode
// ---------------------------------------------------------------------------

// GetName returns the declared name of a FUNCTION or CLASS node, or "" if
// the node has no discernible name.
// Mirrors NodeUtil.getName.
func (nodeUtil) GetName(n *rhino.Node) string {
	nameNode := NU.GetNameNode(n)
	if nameNode == nil {
		return ""
	}
	if nameNode.IsQualifiedName() {
		return nameNode.GetQualifiedName()
	}
	return nameNode.GetString()
}

// GetNameNode returns the AST node that best represents the name of a FUNCTION
// or CLASS node, or nil when no name can be identified.
// Mirrors NodeUtil.getNameNode.
func (nodeUtil) GetNameNode(n *rhino.Node) *rhino.Node {
	parent := n.Parent()
	if parent == nil {
		return nil
	}
	switch parent.Token() {
	case rhino.NAME:
		return parent
	case rhino.ASSIGN:
		first := parent.FirstChild()
		if first.IsQualifiedName() {
			return first
		}
		return nil
	default:
		funNameNode := n.FirstChild()
		if funNameNode == nil || funNameNode.IsEmpty() {
			return nil
		}
		if funNameNode.IsName() && funNameNode.GetString() == "" {
			return nil
		}
		return funNameNode
	}
}

// ---------------------------------------------------------------------------
// IsStatementParent
// ---------------------------------------------------------------------------

// IsStatementParent reports whether parent is a node that may contain
// statement children.
// Mirrors NodeUtil.isStatementParent.
func (nodeUtil) IsStatementParent(parent *rhino.Node) bool {
	switch parent.Token() {
	case rhino.SCRIPT, rhino.MODULE_BODY, rhino.BLOCK, rhino.LABEL:
		return true
	}
	return false
}

// TryMergeBlock merges a block node with its parent block when safe to do so.
// Returns true if the block was merged.
// Mirrors NodeUtil.tryMergeBlock (Java).
func (nodeUtil) TryMergeBlock(block *rhino.Node, ignoreBlockScopedDeclarations bool) bool {
	if !block.IsBlock() {
		return false
	}
	parent := block.Parent()
	canMerge := ignoreBlockScopedDeclarations || NU.CanMergeBlock(block)
	if NU.IsStatementBlock(parent) && canMerge {
		prev := block
		for block.HasChildren() {
			child := block.RemoveFirstChild()
			child.InsertAfter(prev)
			prev = child
		}
		block.Detach()
		return true
	}
	return false
}

// CanMergeBlock reports whether a block node can be safely merged with its parent
// (i.e., it contains no block-scoped declarations that might cause name collisions).
// Mirrors NodeUtil.canMergeBlock (Java).
func (nodeUtil) CanMergeBlock(block *rhino.Node) bool {
	for c := block.FirstChild(); c != nil; c = c.Next() {
		switch c.Token() {
		case rhino.LABEL:
			if !NU.CanMergeBlock(c) {
				return false
			}
		case rhino.CONST, rhino.LET, rhino.CLASS, rhino.FUNCTION:
			return false
		}
	}
	return true
}

// ---------------------------------------------------------------------------
// Additional Java NodeUtil methods ported to Go
// ---------------------------------------------------------------------------

// IsNameDeclOrSimpleAssignLhs reports whether n is the LHS of a simple
// assignment or a name declaration (VAR/LET/CONST).
// Mirrors NodeUtil.isNameDeclOrSimpleAssignLhs.
func (nodeUtil) IsNameDeclOrSimpleAssignLhs(n, parent *rhino.Node) bool {
	return (parent.IsAssign() && parent.FirstChild() == n) || NU.IsNameDeclaration(parent)
}

// IsDeclarationLValue reports whether n is an l-value that also constitutes a
// binding declaration (e.g. `const x`, `import x`, but NOT `x = 5`).
// Mirrors NodeUtil.isDeclarationLValue.
func (nodeUtil) IsDeclarationLValue(n *rhino.Node) bool {
	if !NU.IsLValue(n) {
		return false
	}
	parent := n.Parent()
	if parent == nil {
		return false
	}
	switch parent.Token() {
	case rhino.IMPORT_SPEC, rhino.VAR, rhino.LET, rhino.CONST,
		rhino.PARAM_LIST, rhino.IMPORT, rhino.CATCH, rhino.CLASS, rhino.FUNCTION:
		return true
	case rhino.STRING_KEY:
		gp := parent.Parent()
		if gp == nil {
			return false
		}
		return NU.IsNameDeclaration(gp.Parent())
	case rhino.OBJECT_PATTERN, rhino.ARRAY_PATTERN:
		gp := parent.Parent()
		if gp == nil {
			return false
		}
		return NU.IsNameDeclaration(gp)
	}
	return false
}

// IsTopLevel reports whether n is a SCRIPT or MODULE_BODY node (a top-level
// scope root in the AST).
// Mirrors NodeUtil.isTopLevel.
func (nodeUtil) IsTopLevel(n *rhino.Node) bool {
	return n.IsScript() || n.IsModuleBody()
}

// IsShallowStatementTree reports whether n is at a "shallow" nesting depth
// for statement trees: control structures, statement blocks, or switch bodies.
// Mirrors NodeUtil.isShallowStatementTree.
func (nodeUtil) IsShallowStatementTree(n *rhino.Node) bool {
	return n == nil || NU.IsControlStructure(n) || NU.IsStatementBlock(n) ||
		n.Token() == rhino.SWITCH_BODY
}

// IsCallTo reports whether n is a CALL to the function identified by
// qualifiedName.
// Mirrors NodeUtil.isCallTo(Node, String).
func (nodeUtil) IsCallTo(n *rhino.Node, qualifiedName string) bool {
	return n.IsCall() && n.FirstChild().MatchesQualifiedNameStr(qualifiedName)
}

// IsGoogModuleGetCall reports whether n is a call to goog.module.get(…).
// Mirrors NodeUtil.isGoogModuleGetCall.
func (nodeUtil) IsGoogModuleGetCall(n *rhino.Node) bool {
	return NU.IsCallTo(n, "goog.module.get")
}

// GetNearestFunctionName returns the nearest inferred name for a FUNCTION
// node, including keys in object literals.
// Mirrors NodeUtil.getNearestFunctionName.
func (nodeUtil) GetNearestFunctionName(n *rhino.Node) string {
	if !n.IsFunction() {
		return ""
	}
	name := NU.GetName(n)
	if name != "" {
		return name
	}
	parent := n.Parent()
	if parent == nil {
		return ""
	}
	switch parent.Token() {
	case rhino.MEMBER_FUNCTION_DEF, rhino.STRING_KEY:
		return parent.GetString()
	case rhino.NUMBER:
		s, _ := NU.GetStringValue(parent)
		return s
	}
	return ""
}

// IsObjLitProperty reports whether node is a property key node in an object
// literal (STRING_KEY, GETTER_DEF, SETTER_DEF, MEMBER_FUNCTION_DEF, or COMPUTED_PROP).
// Mirrors NodeUtil.isObjLitProperty.
func (nodeUtil) IsObjLitProperty(node *rhino.Node) bool {
	return node.IsStringKey() || node.IsGetterDef() || node.IsSetterDef() ||
		node.IsMemberFunctionDef() || node.Token() == rhino.COMPUTED_PROP
}

// IsBlocklessArrowFunctionResult reports whether n is the implicit return
// expression of a blockless arrow function.
// Mirrors NodeUtil.isBlocklessArrowFunctionResult.
func (nodeUtil) IsBlocklessArrowFunctionResult(n *rhino.Node) bool {
	parent := n.Parent()
	return parent != nil && parent.IsFunction() && n == parent.LastChild() && !n.IsBlock()
}

// GetBestJSDocInfoNode returns the node that best carries the JSDoc for n.
// Mirrors NodeUtil.getBestJSDocInfoNode.
func (nodeUtil) GetBestJSDocInfoNode(n *rhino.Node) *rhino.Node {
	if n.IsExprResult() {
		return NU.GetBestJSDocInfoNode(n.FirstChild())
	}
	if n.JSDocInfo() != nil {
		return n
	}
	parent := n.Parent()
	if parent == nil {
		return nil
	}
	switch {
	case parent.IsName() || parent.IsExport() ||
		parent.IsVar() || parent.IsLet() || parent.IsConst():
		return NU.GetBestJSDocInfoNode(parent)
	case parent.IsAssign():
		return NU.GetBestJSDocInfoNode(parent)
	case NU.MayBeObjectLitKey(parent) || parent.Token() == rhino.COMPUTED_PROP:
		if parent.JSDocInfo() != nil {
			return parent
		}
	case parent.IsClass() || parent.IsMemberFunctionDef():
		if parent.JSDocInfo() != nil {
			return parent
		}
	}
	return nil
}

// IsPrototypePropertyDeclaration reports whether n is an assignment to a
// prototype property (e.g. `Foo.prototype.bar = function() {}`).
// Mirrors NodeUtil.isPrototypePropertyDeclaration.
func (nodeUtil) IsPrototypePropertyDeclaration(n *rhino.Node) bool {
	if !NU.IsExprAssign(n) {
		return false
	}
	assign := n.FirstChild()
	return NU.IsPrototypeProperty(assign.FirstChild())
}

// IsNamedExportsLiteral reports whether objectLiteral is an object literal
// where all keys are unquoted string keys whose values are NAME nodes —
// the pattern used for CommonJS-style named exports.
// Mirrors NodeUtil.isNamedExportsLiteral.
func (nodeUtil) IsNamedExportsLiteral(objectLiteral *rhino.Node) bool {
	if !objectLiteral.IsObjectLit() || !objectLiteral.HasChildren() {
		return false
	}
	for key := objectLiteral.FirstChild(); key != nil; key = key.Next() {
		if !key.IsStringKey() || key.IsQuotedStringKey() {
			return false
		}
		first := key.FirstChild()
		if first == nil || !first.IsName() {
			return false
		}
	}
	return true
}

// IsFromTypeSummary reports whether a SCRIPT node comes from a type summary
// file (declared with @typeSummary in its JSDoc).
// Mirrors NodeUtil.isFromTypeSummary.
func (nodeUtil) IsFromTypeSummary(n *rhino.Node) bool {
	if !n.IsScript() {
		return false
	}
	info := n.JSDocInfo()
	return info != nil && info.IsTypeSummary()
}

// Precedence returns the operator precedence (0 = lowest) for a given token
// type. Mirrors NodeUtil.precedence(Token).
func (nodeUtil) Precedence(typ rhino.Token) int {
	switch typ {
	case rhino.COMMA:
		return 0
	case rhino.ASSIGN, rhino.ASSIGN_BITOR, rhino.ASSIGN_BITXOR, rhino.ASSIGN_BITAND,
		rhino.ASSIGN_LSH, rhino.ASSIGN_RSH, rhino.ASSIGN_URSH,
		rhino.ASSIGN_ADD, rhino.ASSIGN_SUB, rhino.ASSIGN_MUL,
		rhino.ASSIGN_EXPONENT, rhino.ASSIGN_DIV, rhino.ASSIGN_MOD,
		rhino.ASSIGN_OR, rhino.ASSIGN_AND, rhino.ASSIGN_COALESCE:
		return 1
	case rhino.YIELD:
		return 2
	case rhino.HOOK:
		return 3
	case rhino.OR:
		return 4
	case rhino.AND:
		return 5
	case rhino.COALESCE:
		return 6
	case rhino.BITOR:
		return 7
	case rhino.BITXOR:
		return 8
	case rhino.BITAND:
		return 9
	case rhino.EQ, rhino.NE, rhino.SHEQ, rhino.SHNE:
		return 10
	case rhino.LT, rhino.GT, rhino.LE, rhino.GE, rhino.INSTANCEOF, rhino.IN:
		return 11
	case rhino.LSH, rhino.RSH, rhino.URSH:
		return 12
	case rhino.SUB, rhino.ADD:
		return 13
	case rhino.MUL, rhino.MOD, rhino.DIV:
		return 14
	case rhino.EXPONENT:
		return 15
	case rhino.AWAIT, rhino.NEW, rhino.DELPROP, rhino.TYPEOF, rhino.VOID,
		rhino.NOT, rhino.BITNOT, rhino.POS, rhino.NEG:
		return 16
	case rhino.INC, rhino.DEC:
		return 17
	case rhino.CALL, rhino.GETELEM, rhino.GETPROP,
		rhino.OPTCHAIN_CALL, rhino.OPTCHAIN_GETELEM, rhino.OPTCHAIN_GETPROP,
		rhino.ARRAYLIT, rhino.ARRAY_PATTERN, rhino.DEFAULT_VALUE,
		rhino.DESTRUCTURING_LHS, rhino.EMPTY, rhino.FALSE,
		rhino.FUNCTION, rhino.CLASS, rhino.NAME, rhino.NULL,
		rhino.NUMBER, rhino.BIGINT, rhino.OBJECTLIT, rhino.OBJECT_PATTERN,
		rhino.REGEXP, rhino.ITER_REST, rhino.OBJECT_REST,
		rhino.ITER_SPREAD, rhino.OBJECT_SPREAD, rhino.STRINGLIT, rhino.STRING_KEY,
		rhino.THIS, rhino.SUPER, rhino.TRUE, rhino.TAGGED_TEMPLATELIT, rhino.TEMPLATELIT,
		rhino.DYNAMIC_IMPORT:
		return 18
	case rhino.CAST:
		return 22
	default:
		return 18 // safe fallback for unknown tokens
	}
}

// GetCallTargetResolvingIndirectCalls returns the effective call target of
// n, looking through (0, x.method)(...) indirect call patterns.
// Mirrors NodeUtil.getCallTargetResolvingIndirectCalls.
func (nodeUtil) GetCallTargetResolvingIndirectCalls(call *rhino.Node) *rhino.Node {
	target := call.FirstChild()
	if target == nil {
		return nil
	}
	if target.IsComma() {
		second := target.GetSecondChild()
		if second != nil && second.IsQualifiedName() {
			return second
		}
	}
	return target
}

// NewDeclaration creates a VAR/LET/CONST declaration node for lhs
// optionally initialized to rhs.
// Mirrors NodeUtil.newDeclaration.
func (nodeUtil) NewDeclaration(lhs, rhs *rhino.Node, declarationType rhino.Token) *rhino.Node {
	if rhs == nil {
		return rhino.IRDeclarationLhs(lhs, declarationType)
	}
	return rhino.IRDeclaration(lhs, rhs, declarationType)
}

// NewQName creates a NAME or GETPROP chain for the given dotted name.
// This is the compiler-context-free version.
// Mirrors NodeUtil.newQName (without compiler).
func (nodeUtil) NewQName(name string) *rhino.Node {
	return NU.NewQualifiedNameNode(name)
}

// NewPropertyAccess creates a GETPROP node for context.name.
// Mirrors NodeUtil.newPropertyAccess (without compiler context).
func (nodeUtil) NewPropertyAccess(context *rhino.Node, name string) *rhino.Node {
	return rhino.IRGetprop(context, name)
}

// NewUndefinedNode creates a VOID(0) node that evaluates to undefined.
// Mirrors NodeUtil.newUndefinedNode.
func (nodeUtil) NewUndefinedNode(srcRef *rhino.Node) *rhino.Node {
	node := rhino.IRVoid(rhino.NewNumber(0))
	if srcRef != nil {
		node.SrcrefTree(srcRef)
	}
	return node
}

// EmptyFunction creates an empty FUNCTION node: function() {}.
// Mirrors NodeUtil.emptyFunction.
func (nodeUtil) EmptyFunction() *rhino.Node {
	return rhino.IRFunction(rhino.NewString(rhino.NAME, ""), rhino.IRParamList(), rhino.IRBlock())
}

// RemoveName sets the first child of a FUNCTION or CLASS node to an empty
// name, effectively anonymizing the function/class.
// Mirrors NodeUtil.removeName.
func (nodeUtil) RemoveName(n *rhino.Node) {
	orig := n.FirstChild()
	if orig == nil {
		return
	}
	var emptyName *rhino.Node
	if n.IsFunction() {
		emptyName = rhino.NewString(rhino.NAME, "")
	} else {
		emptyName = rhino.New(rhino.EMPTY)
	}
	emptyName.Srcref(orig)
	orig.ReplaceWith(emptyName)
}

// GetRootOfQualifiedName returns the root NAME/THIS/SUPER node of a qualified
// name chain (e.g. returns the NAME "foo" for "foo.bar.baz").
// Mirrors NodeUtil.getRootOfQualifiedName(Node).
func (nodeUtil) GetRootOfQualifiedName(qName *rhino.Node) *rhino.Node {
	curr := qName
	for {
		if curr.IsName() || curr.IsThis() || curr.IsSuper() {
			return curr
		}
		first := curr.FirstChild()
		if first == nil {
			return curr
		}
		curr = first
	}
}

// IsObjectDefinePropertiesDefinition reports whether n is a call to
// Object.defineProperties with exactly 3 arguments.
// Mirrors NodeUtil.isObjectDefinePropertiesDefinition.
func (nodeUtil) IsObjectDefinePropertiesDefinition(n *rhino.Node) bool {
	if !n.IsCall() || !n.HasXChildren(3) {
		return false
	}
	getprop := n.FirstChild()
	if !getprop.IsGetProp() {
		return false
	}
	if getprop.GetString() != "defineProperties" {
		return false
	}
	obj := getprop.FirstChild()
	if obj == nil {
		return false
	}
	return (obj.IsName() && obj.GetString() == "Object") ||
		obj.MatchesQualifiedNameStr("$jscomp.global.Object") ||
		obj.MatchesQualifiedNameStr("$jscomp$global.Object")
}

// GetSourceName returns the source file name for n by walking up the parent
// chain until a node with a source file name is found.
// Mirrors NodeUtil.getSourceName.
func (nodeUtil) GetSourceName(n *rhino.Node) string {
	curr := n
	for curr != nil {
		s := curr.GetSourceFileName()
		if s != "" {
			return s
		}
		curr = curr.Parent()
	}
	return ""
}

// ---------------------------------------------------------------------------
// Visitor / iteration methods ported from Java NodeUtil
// ---------------------------------------------------------------------------

// Visitor mirrors com.google.javascript.jscomp.NodeUtil.Visitor: a callback
// invoked for each node during a traversal.
type Visitor func(node *rhino.Node)

// VisitPreOrder visits all nodes in n's subtree in pre-order (parent before
// children), calling v for each node. If traverseChildrenPred returns false
// for a node, its children are skipped.
// Mirrors NodeUtil.visitPreOrder.
func (nodeUtil) VisitPreOrder(node *rhino.Node, v Visitor, traverseChildrenPred func(*rhino.Node) bool) {
	v(node)
	if traverseChildrenPred(node) {
		for c := node.FirstChild(); c != nil; c = c.Next() {
			NU.VisitPreOrder(c, v, traverseChildrenPred)
		}
	}
}

// VisitPostOrder visits all nodes in n's subtree in post-order (children
// before parent). If traverseChildrenPred returns false for a node, its
// children are skipped.
// Mirrors NodeUtil.visitPostOrder.
func (nodeUtil) VisitPostOrder(node *rhino.Node, v Visitor, traverseChildrenPred func(*rhino.Node) bool) {
	if traverseChildrenPred(node) {
		for c := node.FirstChild(); c != nil; {
			next := c.Next()
			NU.VisitPostOrder(c, v, traverseChildrenPred)
			c = next
		}
	}
	v(node)
}

// GetCount counts how many nodes in the subtree of n satisfy pred, subject to
// traversal gating by traverseChildrenPred.
// Mirrors NodeUtil.getCount.
func (nodeUtil) GetCount(n *rhino.Node, pred func(*rhino.Node) bool, traverseChildrenPred func(*rhino.Node) bool) int {
	total := 0
	if pred(n) {
		total++
	}
	if traverseChildrenPred(n) {
		for c := n.FirstChild(); c != nil; c = c.Next() {
			total += NU.GetCount(c, pred, traverseChildrenPred)
		}
	}
	return total
}

// Has reports whether any node in the subtree of n satisfies pred, subject to
// traversal gating by traverseChildrenPred.
// Mirrors NodeUtil.has.
func (nodeUtil) Has(n *rhino.Node, pred func(*rhino.Node) bool, traverseChildrenPred func(*rhino.Node) bool) bool {
	return nuHas(n, pred, traverseChildrenPred)
}

// FindPreorder returns the first node (pre-order) in the subtree of n that
// satisfies pred, or nil. Children of non-matching nodes whose
// IteratesImpureIterable reports whether evaluating n may trigger an
// impure (side-effecting) iterator. Only called for nodes that
// may iterate: ITER_SPREAD and ITER_REST.
// Mirrors NodeUtil.iteratesImpureIterable (subset — YIELD/FOR_OF handled separately).
func (nodeUtil) IteratesImpureIterable(n *rhino.Node) bool {
	var iterable *rhino.Node
	switch n.Token() {
	case rhino.ITER_SPREAD:
		iterable = n.FirstChild()
	case rhino.ITER_REST:
		// Rest in a PARAM_LIST is flat (no iteration); in ARRAY_PATTERN it iterates.
		if p := n.Parent(); p != nil && p.IsParamList() {
			return false
		}
		return true // conservative: assume impure iterable
	default:
		return false
	}
	if iterable == nil {
		return true
	}
	// Only array/string/template literals are guaranteed pure iterables.
	switch iterable.Token() {
	case rhino.ARRAYLIT, rhino.STRINGLIT, rhino.TEMPLATELIT:
		return false
	}
	return true
}

// EvaluatesToLocalValue reports whether n evaluates to a freshly created
// value that cannot be aliased by a non-local name.
// Mirrors NodeUtil.evaluatesToLocalValue (simplified: no AccessorSummary,
// no type information).
func (nodeUtil) EvaluatesToLocalValue(value *rhino.Node) bool {
	if value == nil {
		return false
	}
	switch value.Token() {
	case rhino.ASSIGN:
		return NU.IsImmutableValue(value.LastChild())
	case rhino.COMMA:
		return NU.EvaluatesToLocalValue(value.LastChild())
	case rhino.AND, rhino.OR, rhino.COALESCE:
		return NU.EvaluatesToLocalValue(value.FirstChild()) &&
			NU.EvaluatesToLocalValue(value.LastChild())
	case rhino.HOOK:
		return NU.EvaluatesToLocalValue(value.GetSecondChild()) &&
			NU.EvaluatesToLocalValue(value.LastChild())
	case rhino.DYNAMIC_IMPORT:
		return true // always returns a fresh Promise
	case rhino.THIS, rhino.SUPER:
		return false
	case rhino.NAME:
		return NU.IsImmutableValue(value)
	case rhino.GETELEM, rhino.GETPROP, rhino.OPTCHAIN_GETELEM, rhino.OPTCHAIN_GETPROP:
		return false
	case rhino.CALL, rhino.OPTCHAIN_CALL:
		// Only toString()-style method calls produce local string values;
		// without type info we cannot confirm, so conservatively return false.
		return false
	case rhino.TAGGED_TEMPLATELIT:
		return false
	case rhino.NEW:
		return value.IsOnlyModifiesThisCall()
	case rhino.DELPROP, rhino.INC, rhino.DEC, rhino.CLASS, rhino.FUNCTION,
		rhino.REGEXP, rhino.EMPTY, rhino.ARRAYLIT, rhino.OBJECTLIT, rhino.TEMPLATELIT:
		return true
	case rhino.CAST:
		return NU.EvaluatesToLocalValue(value.FirstChild())
	case rhino.ITER_SPREAD, rhino.OBJECT_SPREAD:
		return false
	case rhino.YIELD, rhino.AWAIT:
		return false
	default:
		if value.Token().IsAssign() {
			return false // logical-assignment ops: lhs may not be local
		}
		if isSimpleOperatorToken(value.Token()) || NU.IsImmutableValue(value) {
			return true
		}
		return false // conservative
	}
}

// AllArgsUnescapedLocal reports whether every argument to a CALL or NEW node
// evaluates to a freshly created local value (cannot be aliased externally).
// Mirrors NodeUtil.allArgsUnescapedLocal.
func (nodeUtil) AllArgsUnescapedLocal(callOrNew *rhino.Node) bool {
	// Arguments start at the second child (first child is the callee).
	for arg := callOrNew.GetSecondChild(); arg != nil; arg = arg.Next() {
		if !NU.EvaluatesToLocalValue(arg) {
			return false
		}
	}
	return true
}

// traverseChildrenPred returns false are skipped.
// Mirrors NodeUtil.findPreorder.
func (nodeUtil) FindPreorder(n *rhino.Node, pred func(*rhino.Node) bool, traverseChildrenPred func(*rhino.Node) bool) *rhino.Node {
	if pred(n) {
		return n
	}
	if !traverseChildrenPred(n) {
		return nil
	}
	for c := n.FirstChild(); c != nil; c = c.Next() {
		if found := NU.FindPreorder(c, pred, traverseChildrenPred); found != nil {
			return found
		}
	}
	return nil
}
