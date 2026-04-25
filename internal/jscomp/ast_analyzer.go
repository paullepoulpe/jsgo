// Port of com/google/javascript/jscomp/AstAnalyzer.java @ closure-compiler v20260415
// (subset: MayHaveSideEffects, MayEffectMutableState, NodeTypeMayHaveSideEffects,
//  FunctionCallHasSideEffects, ConstructorCallHasSideEffects)
//
// Implements the astAnalyzer interface declared in nodeiterators.go.
// Tokens absent from rhino/token.go default to true (conservative).
// assumeGettersArePure=true, assumeKnownBuiltinsArePure=true (no AccessorSummary).

package jscomp

import "github.com/paullepoulpe/jsgo/internal/rhino"

var builtinConstructorsNoSideEffects = map[string]bool{
	"Array": true, "Date": true, "Error": true,
	"Object": true, "RegExp": true, "XMLHttpRequest": true,
}

var builtinFunctionsNoSideEffects = map[string]bool{
	"Object": true, "Array": true, "String": true, "Number": true,
	"BigInt": true, "Boolean": true, "RegExp": true, "Error": true,
}

// AstAnalyzer answers whether AST nodes have side effects.
// Mirrors com.google.javascript.jscomp.AstAnalyzer.
type AstAnalyzer struct{}

// NewAstAnalyzer constructs a zero-value AstAnalyzer.
func NewAstAnalyzer() *AstAnalyzer { return &AstAnalyzer{} }

// MayEffectMutableState mirrors AstAnalyzer.mayEffectMutableState.
func (a *AstAnalyzer) MayEffectMutableState(n *rhino.Node) bool {
	return a.checkForStateChangeHelper(n, true)
}

// MayHaveSideEffects mirrors AstAnalyzer.mayHaveSideEffects.
func (a *AstAnalyzer) MayHaveSideEffects(n *rhino.Node) bool {
	return a.checkForStateChangeHelper(n, false)
}

// FunctionCallHasSideEffects mirrors AstAnalyzer.functionCallHasSideEffects.
// Simplified: skips @nosideeffects and regexp/string heuristics.
func (a *AstAnalyzer) FunctionCallHasSideEffects(callNode *rhino.Node) bool {
	callee := callNode.FirstChild()
	if callee != nil && callee.Token() == rhino.NAME {
		if builtinFunctionsNoSideEffects[callee.String_()] {
			return false
		}
	}
	return true
}

// ConstructorCallHasSideEffects mirrors AstAnalyzer.constructorCallHasSideEffects.
func (a *AstAnalyzer) ConstructorCallHasSideEffects(newNode *rhino.Node) bool {
	nameNode := newNode.FirstChild()
	if nameNode == nil || nameNode.Token() != rhino.NAME {
		return true
	}
	return !builtinConstructorsNoSideEffects[nameNode.String_()]
}

// NodeTypeMayHaveSideEffects mirrors AstAnalyzer.nodeTypeMayHaveSideEffects.
// Tokens not in rhino/token.go fall to default → true (conservative).
func (a *AstAnalyzer) NodeTypeMayHaveSideEffects(n *rhino.Node) bool {
	if n == nil {
		return false
	}
	if n.Token().IsAssign() {
		return true
	}
	switch n.Token() {
	case rhino.DELPROP, rhino.THROW,
		rhino.FOR_IN, rhino.FOR_OF, rhino.FOR_AWAIT_OF:
		return true
	case rhino.CALL, rhino.TAGGED_TEMPLATELIT:
		return a.FunctionCallHasSideEffects(n)
	case rhino.NEW:
		return a.ConstructorCallHasSideEffects(n)
	case rhino.NAME:
		return n.FirstChild() != nil
	case rhino.ITER_REST:
		return true // may iterate impure iterable; conservative
	case rhino.GETPROP, rhino.STRING_KEY:
		return false // assumeGettersArePure=true
	default:
		if isSimpleOperatorToken(n.Token()) {
			return false
		}
		return true // unknown token — conservative
	}
}

// checkForStateChangeHelper is the recursive core used by MayHaveSideEffects
// and MayEffectMutableState. Mirrors AstAnalyzer.checkForStateChangeHelper.
func (a *AstAnalyzer) checkForStateChangeHelper(n *rhino.Node, checkForNewObjects bool) bool {
	if n == nil {
		return false
	}
	switch n.Token() {
	case rhino.THROW, rhino.FOR_AWAIT_OF, rhino.FOR_OF, rhino.FOR_IN,
		rhino.VAR, rhino.LET, rhino.CONST:
		return true
	case rhino.SUPER:
		return false
	case rhino.OBJECTLIT, rhino.ARRAYLIT, rhino.REGEXP:
		return checkForNewObjects || a.checkForChangeStateInChildren(n, checkForNewObjects)
	case rhino.ITER_REST:
		return a.checkForChangeStateInChildren(n, checkForNewObjects)
	case rhino.NAME:
		return n.FirstChild() != nil
	case rhino.FUNCTION:
		return checkForNewObjects || NU.IsFunctionDeclaration(n)
	case rhino.MEMBER_FUNCTION_DEF:
		return false
	case rhino.COMPUTED_PROP:
		parent := n.Parent()
		if parent != nil {
			switch parent.Token() {
			case rhino.OBJECTLIT:
				return a.checkForChangeStateInChildren(n, checkForNewObjects)
			case rhino.CLASS:
				if fc := n.FirstChild(); fc != nil {
					return a.checkForStateChangeHelper(fc, checkForNewObjects)
				}
				return false
			}
		}
		return true // conservative for OBJECT_PATTERN etc.
	case rhino.CLASS_FIELD_DEF:
		return true // conservative: may run initialiser
	case rhino.CLASS:
		if checkForNewObjects || NU.IsStatement(n) {
			return true
		}
		kids := n.Children()
		if len(kids) >= 2 && a.checkForStateChangeHelper(kids[1], checkForNewObjects) {
			return true
		}
		if len(kids) >= 3 && a.checkForStateChangeHelper(kids[2], checkForNewObjects) {
			return true
		}
		return false
	case rhino.NEW:
		return checkForNewObjects || a.ConstructorCallHasSideEffects(n) ||
			a.checkForChangeStateInChildren(n, checkForNewObjects)
	case rhino.CALL, rhino.TAGGED_TEMPLATELIT:
		return a.FunctionCallHasSideEffects(n) ||
			a.checkForChangeStateInChildren(n, checkForNewObjects)
	case rhino.AND, rhino.BLOCK, rhino.EXPR_RESULT, rhino.HOOK, rhino.IF,
		rhino.PARAM_LIST, rhino.NUMBER, rhino.OR, rhino.COALESCE,
		rhino.THIS, rhino.TRUE, rhino.FALSE, rhino.NULL, rhino.STRINGLIT,
		rhino.SWITCH, rhino.TEMPLATELIT_SUB, rhino.TRY,
		rhino.TEMPLATELIT, rhino.TEMPLATELIT_STRING,
		rhino.STRING_KEY, rhino.GETELEM, rhino.GETPROP:
		return a.checkForChangeStateInChildren(n, checkForNewObjects)
	default:
		if n.Token().IsAssign() {
			return a.checkAssignmentForChangeState(n, checkForNewObjects)
		}
		if isSimpleOperatorToken(n.Token()) {
			return a.checkForChangeStateInChildren(n, checkForNewObjects)
		}
		return true // unknown token — conservative
	}
}

func (a *AstAnalyzer) checkForChangeStateInChildren(n *rhino.Node, checkForNewObjects bool) bool {
	for c := n.FirstChild(); c != nil; c = c.Next() {
		if a.checkForStateChangeHelper(c, checkForNewObjects) {
			return true
		}
	}
	return false
}

// checkAssignmentForChangeState mirrors AstAnalyzer.checkAssignmentForChangeState.
// Simplified: without evaluatesToLocalValue, non-NAME assignments are conservative.
func (a *AstAnalyzer) checkAssignmentForChangeState(assignOp *rhino.Node, checkForNewObjects bool) bool {
	target := assignOp.FirstChild()
	if target == nil || target.Token() == rhino.NAME {
		return true
	}
	value := assignOp.LastChild()
	if a.checkForStateChangeHelper(target, checkForNewObjects) ||
		(value != nil && a.checkForStateChangeHelper(value, checkForNewObjects)) {
		return true
	}
	// Conservative: without evaluatesToLocalValue, assume external state change.
	return true
}

// isSimpleOperatorToken reports whether t is an operator that cannot itself
// cause a side effect. Mirrors NodeUtil.isSimpleOperator in Java.
func isSimpleOperatorToken(t rhino.Token) bool {
	switch t {
	case rhino.ADD, rhino.SUB, rhino.MUL, rhino.DIV, rhino.MOD,
		rhino.LT, rhino.LE, rhino.GT, rhino.GE, rhino.EQ, rhino.NE, rhino.SHEQ, rhino.SHNE,
		rhino.BITAND, rhino.BITOR, rhino.BITXOR, rhino.LSH, rhino.RSH, rhino.URSH,
		rhino.NOT, rhino.BITNOT, rhino.POS, rhino.NEG, rhino.TYPEOF, rhino.VOID,
		rhino.INSTANCEOF, rhino.IN, rhino.HOOK:
		return true
	}
	return false
}
