// Port of com/google/javascript/jscomp/ExpressionDecomposer.java @ closure-compiler v20261021

/*
 * Copyright 2009 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ExpressionDecomposer_IsConditionalOp returns whether the node is a
// conditional op.
//
// Java: ExpressionDecomposer.isConditionalOp(Node)
func ExpressionDecomposer_IsConditionalOp(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_HOOK, rhino.Token_AND, rhino.Token_OR, rhino.Token_COALESCE,
		rhino.Token_OPTCHAIN_GETELEM, rhino.Token_OPTCHAIN_GETPROP, rhino.Token_OPTCHAIN_CALL:
		return true
	default:
		return false
	}
}

// ExpressionDecomposer_FindExpressionRoot finds the statement containing
// subExpression.
//
// If subExpression is not contained by a statement where inlining is known
// to be possible, null is returned. For example, the condition expression
// of a WHILE loop.
//
// Java: ExpressionDecomposer.findExpressionRoot(Node)
func ExpressionDecomposer_FindExpressionRoot(subExpression rhino.Node) rhino.Node {
	child := subExpression
	for current := child.GetParent(); !current.IsNil(); current = current.GetParent() {
		parent := current.GetParent()
		switch current.GetToken() {
		// Supported expression roots:
		// SWITCH and IF can have multiple children, but the CASE, DEFAULT,
		// or BLOCK will be encountered first for any of the children other
		// than the condition.
		case rhino.Token_EXPR_RESULT,
			rhino.Token_IF,
			rhino.Token_SWITCH,
			rhino.Token_RETURN,
			rhino.Token_THROW:
			jutil.CheckState(child.IsFirstChildOf(current))
			return current

		case rhino.Token_VAR,
			// Normalization will remove LABELs from VARs.
			rhino.Token_LET,
			rhino.Token_CONST:
			if NodeUtil_IsAnyFor(parent) {
				break // Name declarations may not be roots if they're for-loop initializers.
			}
			return current

		// Any of these indicate an unsupported expression:
		case rhino.Token_FOR:
			if ast.RoleIn(child, current) == ast.RoleForInit {
				// Only the initializer of a for-loop could possibly be decomposed since the other
				// statements need to execute each iteration.
				return current
			}
			fallthrough
		case rhino.Token_FOR_IN,
			rhino.Token_FOR_OF,
			rhino.Token_FOR_AWAIT_OF,
			rhino.Token_DO,
			rhino.Token_WHILE,
			rhino.Token_SCRIPT,
			rhino.Token_BLOCK,
			rhino.Token_LABEL,
			rhino.Token_CASE,
			rhino.Token_SWITCH_BODY,
			rhino.Token_DEFAULT_CASE,
			rhino.Token_DEFAULT_VALUE,
			rhino.Token_PARAM_LIST,
			// For top-level class declarations, without this, we would have eventually returned `null`
			// due to `case SCRIPT` above. This just makes class expressions behave the same way. Future
			// optimizations may want to better handle class members and static blocks.
			rhino.Token_CLASS:
			return rhino.NilNode

		default:
		}
		child = current
	}

	panic(jutil.NewIllegalStateExceptionMsg("Unexpected AST structure."))
}

// CanExposeExpression determines if subExpression can be moved before
// expressionRoot without changing the behaviour of the code, or if there is
// a rewriting that would make such motion possible.
//
// Walks the AST from subExpression to expressionRoot and verifies that the
// portions of the expressionRoot subtree that are evaluated before
// subExpression:
//
//  1. are unaffected by the side-effects, if any, of the subExpression.
//  2. have no side-effects that may influence the subExpression.
//  3. have a syntactically legal rewriting.
//
// Examples:
//
//   - expressionRoot = `a = 1 + x();`, subExpression = `x()`, has
//     side-effects: MOVABLE because the final value of `1` cannot be
//     influenced by `x()`.
//   - expressionRoot = `a = b + x();`, subExpression = `x()`, has
//     side-effects: DECOMPOSABLE because `b` may be modified by `x()`, but
//     `b` can be cached.
//   - expressionRoot = `a = b + x();`, subExpression = `x()`, no
//     side-effects: MOVABLE because `x()` can be computed before or after
//     `b` is resolved.
//   - expressionRoot = `a = (b = c) + x();`, subExpression = `x()`, no
//     side-effects, is side-effected: DECOMPOSABLE because `x()` may read
//     `b`.
//
// Returns MOVABLE if subExpression can already be moved; DECOMPOSABLE if
// the expressionRoot subtree could be rewritten such that subExpression
// would be made movable; UNDECOMPOSABLE otherwise.
//
// Java: ExpressionDecomposer.canExposeExpression(Node)
func (e *ExpressionDecomposer) CanExposeExpression(subExpression rhino.Node) ExpressionDecomposer_DecompositionType {
	expressionRoot := ExpressionDecomposer_FindExpressionRoot(subExpression)
	if !expressionRoot.IsNil() {
		return e.IsSubexpressionMovable(expressionRoot, subExpression)
	}
	return ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE
}

// IsSubexpressionMovable: see CanExposeExpression.
//
// Java: ExpressionDecomposer.isSubexpressionMovable(Node, Node)
func (e *ExpressionDecomposer) IsSubexpressionMovable(
	expressionRoot rhino.Node, subExpression rhino.Node) ExpressionDecomposer_DecompositionType {
	requiresDecomposition := false
	seenSideEffects := e.astAnalyzer.MayHaveSideEffects(subExpression)

	if NodeUtil_IsOptChainNode(subExpression) && !NodeUtil_IsEndOfFullOptChain(subExpression) {
		// e.g `sub?.expression.rest?.of.expression`
		// It is always necessary to decompose the prefix of an optional chain.
		requiresDecomposition = true
	}
	child := subExpression
	for parent := child.GetParent(); !parent.IsNil(); parent = parent.GetParent() {
		if decl, isDecl := ast.AsNameDeclaration(parent); isDecl && !decl.MustFirstDeclarator().Is(child) {
			// Most declarations are split by `Normalize` but to do so in loops would change the
			// behavior of the code.  We don't current handle this case but it is possible:
			// For this case: `for (let x = 5, y = 2 * x; ...` where `child = y`.
			// As later expressions may reference the earlier expression, it would be necessary to
			// rewrite references when extracting the expressions:
			// `var temp1 = 5; var temp2 = 2 * temp1; for (let x = temp1, y = temp2; ...`
			return ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE
		}

		if tt, isTagged := ast.AsTaggedTemplateLit(parent); isTagged && !tt.IsFreeCall() {
			// We're looking at something like: something.method`${subExpression()}`
			// You can't use the `.call(something, ...)` trick for a tagged template literal.
			// TODO(b/251958225): Implement decomposition of this case.
			return ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE
		}

		if parent == expressionRoot {
			// Done. The walk back to the root of the expression is complete, and
			// nothing was encountered that blocks the call from being moved.
			if requiresDecomposition {
				return ExpressionDecomposer_DecompositionType_DECOMPOSABLE
			}
			return ExpressionDecomposer_DecompositionType_MOVABLE
		}

		if ExpressionDecomposer_IsConditionalOp(parent) {
			// Only the first child is always executed, otherwise it must be
			// decomposed.
			if child != parent.GetFirstChild() {
				requiresDecomposition = true
			}
		} else {
			// Only inline the call if none of the preceding siblings in the
			// expression have side-effects, and are unaffected by the side-effects,
			// if any, of the call in question.
			// NOTE: The siblings are not always in the order in which they are evaluated, so we call
			// getEvaluationDirection to see in which order to traverse the siblings.

			// SPECIAL CASE: Assignment to a simple name
			if e.IsSafeAssign(parent, seenSideEffects) {
				// It is always safe to inline "foo()" for expressions such as
				//   "a = b = c = foo();"
				// As the assignment is unaffected by side effect of "foo()"
				// and the names assigned-to cannot influence the state before
				// the call to foo.
				//
				// This is not true of more complex LHS values, such as
				//    a.x = foo();
				//    next().x = foo();
				// in these cases the checks below are necessary.
			} else {
				// Everything else.
				direction := ExpressionDecomposer_GetEvaluationDirection(parent)
				for n := e.GetFirstEvaluatedChild(parent, direction); !n.IsNil(); n = e.GetNextEvaluatedSibling(n, direction) {
					if n == child {
						// None of the preceding siblings have side-effects.
						// This is OK.
						break
					}

					if e.IsExpressionTreeUnsafe(n, seenSideEffects) {
						seenSideEffects = true
						requiresDecomposition = true
					}
				}

				if call, isCall := ast.AsCall(parent); requiresDecomposition && isCall {
					if _, isGet := call.Callee().AsNormalGet(); isGet {
						return ExpressionDecomposer_DecompositionType_DECOMPOSABLE
					}
				}
			}
		}
		// Continue looking up the expression tree.
		child = parent
	}

	// With a valid tree we should never get here.
	panic(jutil.NewIllegalStateExceptionMsg("Unexpected."))
}

// ExpressionDecomposer_EvaluationDirection is
// ExpressionDecomposer.EvaluationDirection.
type ExpressionDecomposer_EvaluationDirection int32

const (
	ExpressionDecomposer_EvaluationDirection_FORWARD ExpressionDecomposer_EvaluationDirection = iota
	ExpressionDecomposer_EvaluationDirection_REVERSE
)

var expressionDecomposer_EvaluationDirection_names = [...]string{"FORWARD", "REVERSE"}

// Java: ExpressionDecomposer.EvaluationDirection.name()
func (d ExpressionDecomposer_EvaluationDirection) String() string {
	return expressionDecomposer_EvaluationDirection_names[d]
}

// Java: ExpressionDecomposer.EvaluationDirection.ordinal()
func (d ExpressionDecomposer_EvaluationDirection) Ordinal() int { return int(d) }

// ExpressionDecomposer_GetEvaluationDirection returns the order in which
// the given node's children should be evaluated.
//
// In most cases, this is EvaluationDirection.FORWARD because the AST order
// matches the actual evaluation order. A few nodes require reversed
// evaluation instead.
//
// Java: ExpressionDecomposer.getEvaluationDirection(Node)
func ExpressionDecomposer_GetEvaluationDirection(node rhino.Node) ExpressionDecomposer_EvaluationDirection {
	// The lhs of a destructuring assignment is evaluated AFTER the rhs. This is only true for
	// destructuring, though, not assignments like "first().x = second()" where "first()" is
	// evaluated first.
	switch node.GetToken() {
	case rhino.Token_DESTRUCTURING_LHS:
		return ExpressionDecomposer_EvaluationDirection_REVERSE
	case rhino.Token_ASSIGN:
		if _, isPat := ast.MustAssign(node).Target().AsPattern(); isPat {
			return ExpressionDecomposer_EvaluationDirection_REVERSE
		}
	case rhino.Token_DEFAULT_VALUE:
		if _, isPat := ast.MustDefaultValue(node).Target().AsPattern(); isPat {
			return ExpressionDecomposer_EvaluationDirection_REVERSE
		}
	default:
	}
	return ExpressionDecomposer_EvaluationDirection_FORWARD
}

// Java: ExpressionDecomposer.getFirstEvaluatedChild(Node, ExpressionDecomposer.EvaluationDirection)
func (e *ExpressionDecomposer) GetFirstEvaluatedChild(
	parent rhino.Node, direction ExpressionDecomposer_EvaluationDirection) rhino.Node {
	if direction == ExpressionDecomposer_EvaluationDirection_FORWARD {
		return parent.GetFirstChild()
	}
	return parent.GetLastChild()
}

// Java: ExpressionDecomposer.getNextEvaluatedSibling(Node, ExpressionDecomposer.EvaluationDirection)
func (e *ExpressionDecomposer) GetNextEvaluatedSibling(
	node rhino.Node, direction ExpressionDecomposer_EvaluationDirection) rhino.Node {
	if direction == ExpressionDecomposer_EvaluationDirection_FORWARD {
		return node.GetNext()
	}
	return node.GetPrevious()
}

// IsSafeAssign: it is always safe to inline "foo()" for expressions such as
// "a = b = c = foo();" As the assignment is unaffected by side effect of
// "foo()" and the names assigned-to cannot influence the state before the
// call to foo.
//
// It is also safe in cases where the object is constant:
//
//	CONST_NAME.a = foo()
//	CONST_NAME[CONST_VALUE] = foo();
//
// This is not true of more complex LHS values, such as
//
//	a.x = foo();
//	next().x = foo();
//
// in these cases the checks below are necessary.
//
// seenSideEffects: if true, check to see if node-tree maybe affected by
// side-effects, otherwise if the tree has side-effects (see
// IsExpressionTreeUnsafe). Returns whether the assignment is safe from
// side-effects.
//
// Java: ExpressionDecomposer.isSafeAssign(Node, boolean)
func (e *ExpressionDecomposer) IsSafeAssign(n rhino.Node, seenSideEffects bool) bool {
	if assign, ok := ast.AsAssign(n); ok {
		lhs := assign.Target()
		if _, isName := lhs.AsName(); isName {
			return true
		}
		if gp, isGetProp := lhs.AsGetProp(); isGetProp {
			return !e.IsExpressionTreeUnsafe(gp.Object().Node(), seenSideEffects)
		}
		if ge, isGetElem := lhs.AsGetElem(); isGetElem {
			return !e.IsExpressionTreeUnsafe(ge.Object().Node(), seenSideEffects) &&
				!e.IsExpressionTreeUnsafe(ge.Index().Node(), seenSideEffects)
		}
	}
	return false
}

// We must not decompose this method qname.
//
// See usage location for more information.
var ExpressionDecomposer_WINDOW_LOCATION_ASSIGN = rhino.QualifiedName_Of("window.location.assign")

// IsExpressionTreeUnsafe determines if there is any subexpression below
// tree that would make it incorrect for some expression that follows tree,
// E, to be executed before tree.
//
// followingSideEffectsExist is whether E causes side-effects. Returns true
// if tree contains any subexpressions that would make movement incorrect.
//
// Java: ExpressionDecomposer.isExpressionTreeUnsafe(Node, boolean)
func (e *ExpressionDecomposer) IsExpressionTreeUnsafe(tree rhino.Node, followingSideEffectsExist bool) bool {
	// Spread expressions would cause recursive rewriting if not special cased here.
	// When extracted, spreads can't be assigned to a single variable and instead are put into
	// a literal. However, that literal must be spread again at the original site. This
	// check is what prevents the original spread from triggering recursion.
	if sp, ok := ast.AsIterSpread(tree); ok {
		if e.IsTempConstantValueName(sp.Expr().Node()) {
			return false
		}
	} else if sp, ok := ast.AsObjectSpread(tree); ok {
		if e.IsTempConstantValueName(sp.Expr().Node()) {
			return false
		}
	}

	if followingSideEffectsExist {
		// If the call to be inlined has side-effects, check to see if this
		// expression tree can be affected by any side-effects.

		// Assume that "tmp1.call(...)" is safe (where tmp1 is a const temp variable created by
		// ExpressionDecomposer) otherwise we end up trying to decompose the same tree
		// an infinite number of times. Also assume that "fn.call" is safe when "fn" is a known
		// constant function, as decomposing it may mess up InlineFunction's bookkeeping if it is
		// attempting to inline "fn".
		parent := tree.GetParent()
		if NodeUtil_IsObjectCallMethod(parent, "call") || parent.GetBooleanProp(rhino.Node_FREE_CALL) {
			var callee rhino.Node
			if gp, ok := ast.AsGetProp(tree); ok && gp.Property() == "call" {
				callee = gp.Object().Node()
			} else {
				callee = tree
			}
			if tree.IsFirstChildOf(parent) {
				if e.IsTempConstantValueName(callee) {
					return false
				}
				if len(e.knownConstantFunctions) > 0 {
					if qn, ok := callee.GetQualifiedNameString(); ok && e.knownConstantFunctions[qn] {
						return false
					}
				}
			}
		}

		if e.brokenIE11LocationAssign &&
			ExpressionDecomposer_WINDOW_LOCATION_ASSIGN.Matches(tree) {
			// IE11 throws an exception if we decompose a call to `window.location.assign`
			// e.g.
			// ```js
			// obj = window.location, method = obj.assign, method.call(obj, url); // exception on IE11
			// ```
			// So we will consider it to be a constant value that doesn't need to be decomposed to
			// protect it against side effects.
			// This is a hack to work around a specific case we've seen in real world code.
			// We know it won't catch some cases, but we would rather miss those than add more
			// sophisticated logic that would also encounter false-positives (e.g. is `x.assign` a
			// reference to `window.location.assign`?).
			//
			// We cannot have canBeSideEffected() (below) do this check for us, because it only accepts
			// simple names as constants.
			return false
		}
		// This is a superset of "AstAnalyzer.mayHaveSideEffects".
		return NodeUtil_CanBeSideEffectedWithScope(tree, e.knownConstantFunctions, e.scope)
	} else {
		// The function called doesn't have side-effects but check to see if there
		// are side-effects that that may affect it.
		return e.astAnalyzer.MayHaveSideEffects(tree)
	}
}
