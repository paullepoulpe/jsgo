// Port of jscomp/AstAnalyzer.java @ closure-compiler v20260415
// (subset: MayHaveSideEffects, MayEffectMutableState, NodeTypeMayHaveSideEffects,
//  FunctionCallHasSideEffects, ConstructorCallHasSideEffects)
//
// Implements the astAnalyzer interface declared in nodeiterators.go.
// Tokens absent from rhino/token.go default to true (conservative).
// assumeGettersArePure=true, assumeKnownBuiltinsArePure=true (no AccessorSummary).
// hasRegexpGlobalReferences=false (no global regexp state).

package jscomp

import "github.com/paullepoulpe/jsgo/legacy/internal/rhino"

var builtinConstructorsNoSideEffects = map[string]bool{
	"Array": true, "Date": true, "Error": true,
	"Object": true, "RegExp": true, "XMLHttpRequest": true,
}

var builtinFunctionsNoSideEffects = map[string]bool{
	"Object": true, "Array": true, "String": true, "Number": true,
	"BigInt": true, "Boolean": true, "RegExp": true, "Error": true,
}

var objectMethodsNoSideEffects = map[string]bool{
	"toString": true, "valueOf": true,
}

var regexpMethods = map[string]bool{
	"test": true, "exec": true,
}

var stringRegexpMethods = map[string]bool{
	"match": true, "replace": true, "search": true, "split": true,
}

var mathFunctionsNoSideEffects = map[string]bool{
	"abs": true, "acos": true, "acosh": true, "asin": true, "asinh": true,
	"atan": true, "atanh": true, "atan2": true, "cbrt": true, "ceil": true,
	"cos": true, "cosh": true, "exp": true, "expm1": true, "floor": true,
	"hypot": true, "log": true, "log10": true, "log1p": true, "log2": true,
	"max": true, "min": true, "pow": true, "round": true, "sign": true,
	"sin": true, "sinh": true, "sqrt": true, "tan": true, "tanh": true,
	"trunc": true,
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
// assumeGettersArePure=true, assumeKnownBuiltinsArePure=true, hasRegexpGlobalReferences=false.
func (a *AstAnalyzer) FunctionCallHasSideEffects(callNode *rhino.Node) bool {
	// @nosideeffects annotation on the call site.
	if callNode.IsNoSideEffectsCall() {
		return false
	}

	// If the call only modifies its arguments and all arguments are local,
	// there are no externally observable side effects.
	if callNode.IsOnlyModifiesArgumentsCall() && NU.AllArgsUnescapedLocal(callNode) {
		return false
	}

	callee := callNode.FirstChild()
	if callee == nil {
		return true
	}

	// Built-in functions with no side effects (NAME callee).
	if callee.IsName() {
		if builtinFunctionsNoSideEffects[callee.String_()] {
			return false
		}
	} else if callee.IsGetProp() || callee.IsOptChainGetProp() {
		method := callee.String_()

		// Object.prototype methods: toString / valueOf with no extra args.
		if callNode.HasOneChild() && objectMethodsNoSideEffects[method] {
			return false
		}

		// @modifiesThis call on a local receiver is safe.
		if callNode.IsOnlyModifiesThisCall() &&
			NU.EvaluatesToLocalValue(callee.FirstChild()) {
			return false
		}

		// Math.* functions with no side effects.
		// assumeKnownBuiltinsArePure=true.
		if fc := callee.FirstChild(); fc != nil && fc.IsName() &&
			fc.String_() == "Math" && callee.IsQualifiedName() {
			if mathFunctionsNoSideEffects[method] {
				return false
			}
			if method == "random" {
				// Math.random() takes no args; hasOneChild means callee only.
				return !callNode.HasOneChild()
			}
			// Unknown Math.* — fall through.
		}

		// Regexp / string method heuristics.
		// hasRegexpGlobalReferences=false, assumeKnownBuiltinsArePure=true.
		receiver := callee.FirstChild()
		if receiver != nil {
			if receiver.IsRegExp() && regexpMethods[method] {
				return false
			}
			// String methods: receiver is a string literal or typed as string.
			if receiver.IsStringLit() {
				param := callee.Next()
				if param != nil {
					if param.IsStringLit() && stringRegexpMethods[method] {
						return false
					}
					if param.IsRegExp() {
						if method == "replace" {
							// Assume anything but a string constant has side-effects.
							return param.Next() == nil || !param.Next().IsStringLit()
						}
						if stringRegexpMethods[method] {
							return false
						}
					}
				}
			}
		}
	}

	return true
}

// ConstructorCallHasSideEffects mirrors AstAnalyzer.constructorCallHasSideEffects.
func (a *AstAnalyzer) ConstructorCallHasSideEffects(newNode *rhino.Node) bool {
	if newNode.IsNoSideEffectsCall() {
		return false
	}
	if newNode.IsOnlyModifiesArgumentsCall() && NU.AllArgsUnescapedLocal(newNode) {
		return false
	}
	nameNode := newNode.FirstChild()
	if nameNode == nil || !nameNode.IsName() {
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
	case rhino.DELPROP, rhino.DEC, rhino.INC, rhino.YIELD, rhino.THROW,
		rhino.AWAIT, rhino.DYNAMIC_IMPORT:
		return true
	case rhino.FOR_IN, rhino.FOR_OF, rhino.FOR_AWAIT_OF:
		return true
	case rhino.OPTCHAIN_CALL, rhino.CALL, rhino.TAGGED_TEMPLATELIT:
		return a.FunctionCallHasSideEffects(n)
	case rhino.NEW:
		return a.ConstructorCallHasSideEffects(n)
	case rhino.NAME:
		return n.FirstChild() != nil
	case rhino.DESTRUCTURING_LHS:
		return true
	// assumeGettersArePure=true → OBJECT_REST/OBJECT_SPREAD have no side effects.
	case rhino.OBJECT_REST, rhino.OBJECT_SPREAD:
		return false
	case rhino.ITER_REST, rhino.ITER_SPREAD:
		return NU.IteratesImpureIterable(n)
	case rhino.STRING_KEY:
		// assumeGettersArePure=true → reading a named key never triggers a getter.
		return false
	case rhino.GETPROP, rhino.OPTCHAIN_GETPROP:
		// assumeGettersArePure=true → property reads are pure.
		return false
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
	// Throw is a side-effect by definition.
	// Context switches can conceal side-effects.
	// Enhanced for loops are almost always side-effectful.
	// Variable declarations are side-effects.
	case rhino.THROW, rhino.YIELD, rhino.AWAIT, rhino.FOR_AWAIT_OF,
		rhino.FOR_OF, rhino.FOR_IN,
		rhino.VAR, rhino.LET, rhino.CONST,
		rhino.EXPORT, rhino.DYNAMIC_IMPORT:
		return true

	// The super keyword is a noop on its own.
	case rhino.SUPER:
		return false

	case rhino.OBJECTLIT, rhino.ARRAYLIT, rhino.REGEXP:
		return checkForNewObjects || a.checkForChangeStateInChildren(n, checkForNewObjects)

	// Object-rest and object-spread may trigger a getter.
	// assumeGettersArePure=true → treat as pure structural ops.
	case rhino.OBJECT_REST, rhino.OBJECT_SPREAD:
		// assumeGettersArePure=true: !assumeGettersArePure is false.
		return a.checkForChangeStateInChildren(n, checkForNewObjects)

	case rhino.ITER_REST, rhino.ITER_SPREAD:
		return NU.IteratesImpureIterable(n) ||
			a.checkForChangeStateInChildren(n, checkForNewObjects)

	case rhino.NAME:
		// NAMEs have children if the left side of a var/let/const.
		return n.FirstChild() != nil

	case rhino.FUNCTION:
		// Function expressions don't have side-effects, but function
		// declarations change the namespace.
		return checkForNewObjects || NU.IsFunctionDeclaration(n)

	// Simply defining a member function, getter, or setter has no side effects.
	case rhino.GETTER_DEF, rhino.SETTER_DEF, rhino.MEMBER_FUNCTION_DEF:
		return false

	case rhino.COMPUTED_PROP:
		parent := n.Parent()
		if parent != nil {
			switch parent.Token() {
			case rhino.CLASS_MEMBERS:
				// Computed method: only the key expression can have side effects.
				if fc := n.FirstChild(); fc != nil {
					return a.checkForStateChangeHelper(fc, checkForNewObjects)
				}
				return false
			case rhino.OBJECT_PATTERN:
				// Only the last child can be an OBJECT_REST.
				return parent.LastChild().IsObjectRest() ||
					a.checkForChangeStateInChildren(n, checkForNewObjects)
			case rhino.OBJECTLIT:
				// Assume COMPUTED_PROP keys in OBJECTLIT never trigger getters.
				return a.checkForChangeStateInChildren(n, checkForNewObjects)
			}
		}
		return true // conservative for unknown parent

	case rhino.MEMBER_FIELD_DEF:
		// Static fields with initializers run the initializer at class definition time.
		return n.IsStaticMember() && n.HasChildren() &&
			a.checkForStateChangeHelper(n.FirstChild(), checkForNewObjects)

	case rhino.CLASS:
		return checkForNewObjects ||
			NU.IsClassDeclaration(n) ||
			// Check the extends clause for side effects.
			a.checkForStateChangeHelper(n.GetSecondChild(), checkForNewObjects) ||
			// Check for class members that are computed properties with side effects.
			a.checkForStateChangeHelper(n.LastChild(), checkForNewObjects)

	case rhino.NEW:
		return checkForNewObjects ||
			a.ConstructorCallHasSideEffects(n) ||
			a.checkForChangeStateInChildren(n, checkForNewObjects)

	case rhino.CALL, rhino.OPTCHAIN_CALL, rhino.TAGGED_TEMPLATELIT:
		return a.FunctionCallHasSideEffects(n) ||
			a.checkForChangeStateInChildren(n, checkForNewObjects)

	case rhino.CAST, rhino.AND, rhino.BLOCK, rhino.ROOT, rhino.EXPR_RESULT,
		rhino.HOOK, rhino.IF, rhino.CLASS_MEMBERS,
		rhino.PARAM_LIST, rhino.DEFAULT_VALUE,
		rhino.NUMBER, rhino.BIGINT, rhino.OR, rhino.COALESCE,
		rhino.THIS, rhino.TRUE, rhino.FALSE, rhino.NULL, rhino.STRINGLIT,
		rhino.SWITCH, rhino.TEMPLATELIT_SUB, rhino.TRY, rhino.EMPTY,
		rhino.TEMPLATELIT, rhino.TEMPLATELIT_STRING:
		return a.checkForChangeStateInChildren(n, checkForNewObjects)

	case rhino.STRING_KEY:
		// assumeGettersArePure=true: reading a named key never triggers a getter;
		// OBJECT_REST handling also only relevant for OBJECT_PATTERN parent.
		// For OBJECT_PATTERN with trailing OBJECT_REST, the key presence matters —
		// but without AccessorSummary, we are conservative only on the getter check.
		parent := n.Parent()
		if parent != nil && parent.IsObjectPattern() {
			// No AccessorSummary: assume getPropertyKind returns NORMAL (no getter).
			// Check for trailing OBJECT_REST: presence of this key affects rest spread.
			if parent.LastChild().IsObjectRest() {
				return true
			}
		}
		return a.checkForChangeStateInChildren(n, checkForNewObjects)

	// GETELEM / OPTCHAIN_GETELEM: we cannot tell which property is accessed,
	// but we assume no getter is triggered (matches Java's comment about
	// code-size cost of conservative treatment).
	case rhino.GETELEM, rhino.OPTCHAIN_GETELEM:
		return a.checkForChangeStateInChildren(n, checkForNewObjects)

	case rhino.GETPROP, rhino.OPTCHAIN_GETPROP:
		// assumeGettersArePure=true → getPropertyKind returns NORMAL (no getter/setter).
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
func (a *AstAnalyzer) checkAssignmentForChangeState(assignOp *rhino.Node, checkForNewObjects bool) bool {
	assignTarget := assignOp.FirstChild()
	if assignTarget == nil {
		return true
	}

	// NAME assignment always has side effects (modifies a binding).
	if assignTarget.IsName() {
		return true
	}

	assignValue := assignOp.LastChild()
	// Side effects if LHS or RHS has side effects.
	if a.checkForStateChangeHelper(assignTarget, checkForNewObjects) ||
		(assignValue != nil && a.checkForStateChangeHelper(assignValue, checkForNewObjects)) {
		return true
	}

	if !NU.IsNormalGet(assignTarget) {
		// Legacy: preserve existing behavior for non-normal-get LHS.
		return !NU.IsLiteralValue(assignTarget, true)
	}

	// If the object being assigned to is a local object, this is not a
	// visible side effect (can't be referenced elsewhere).
	current := assignTarget.FirstChild()
	if current != nil && NU.EvaluatesToLocalValue(current) {
		return false
	}

	// Walk down to the root object; if it's a literal it can't be aliased.
	for current != nil && NU.IsNormalGet(current) {
		current = current.FirstChild()
	}
	if current == nil {
		return true
	}
	return !NU.IsLiteralValue(current, true)
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
