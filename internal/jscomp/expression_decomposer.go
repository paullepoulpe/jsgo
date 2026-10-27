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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: conditional, call, property-access, spread, and assignment shapes
// are read through package ast while keeping ExpressionDecomposer's rhino.Node
// API intact.

// ExpressionDecomposer partially or fully decomposes an expression with
// respect to some sub-expression. Initially this is intended to expand the
// locations where inlining can occur, but has other uses as well.
//
// For example: `var x = y() + z();` becomes `var a = y(); var b = z(); var x
// = a + b;`.
//
// Decomposing, in this context does not mean full decomposition to "atomic"
// expressions. While it is possible to iteratively apply decomposition to
// get statements with at most one side-effect, that isn't the intended
// purpose of this class. The focus is on decomposing "just enough" to
// "free" a particular subexpression. For example:
//
//   - Given: `return (alert() + alert()) + z();`
//   - Exposing: `z()`
//   - Sufficient decomposition: `var temp = alert() + alert(); return temp + z();`
type ExpressionDecomposer struct {
	compiler                 AbstractCompiler
	astAnalyzer              *AstAnalyzer
	astFactory               *AstFactory
	safeNameIdSupplier       func() string
	knownConstantFunctions   map[string]bool
	scope                    *Scope
	brokenIE11LocationAssign bool
	unknownType              rhino.JSType

	tempNamePrefix   string
	resultNamePrefix string
}

// ExpressionDecomposer_DecompositionType is the type of decomposition that
// can be performed on an expression.
//
// See CanExposeExpression.
type ExpressionDecomposer_DecompositionType int32

const (
	ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE ExpressionDecomposer_DecompositionType = iota
	ExpressionDecomposer_DecompositionType_MOVABLE
	ExpressionDecomposer_DecompositionType_DECOMPOSABLE
)

var expressionDecomposer_DecompositionType_names = [...]string{
	"UNDECOMPOSABLE", "MOVABLE", "DECOMPOSABLE",
}

// Java: ExpressionDecomposer.DecompositionType.name()
func (t ExpressionDecomposer_DecompositionType) String() string {
	return expressionDecomposer_DecompositionType_names[t]
}

// Java: ExpressionDecomposer.DecompositionType.ordinal()
func (t ExpressionDecomposer_DecompositionType) Ordinal() int { return int(t) }

// NewExpressionDecomposer creates an ExpressionDecomposer.
//
// DO NOT USE THIS METHOD DIRECTLY. Call
// `compiler.createExpressionDecomposer()` instead so the compiler can supply
// the workarounds needed for this compilation.
//
// constFunctionNames is the set of names known to be constant functions.
// Used by InlineFunctions to prevent this pass from breaking bookkeeping
// for functions it's inlining. brokenIE11LocationAssign indicates whether the
// IE11 window.location.assign workaround should be enabled.
//
// Java: ExpressionDecomposer(AbstractCompiler, Supplier<String>, ImmutableSet<String>, Scope, EnumSet<ExpressionDecomposer.Workaround>)
func NewExpressionDecomposer(
	compiler AbstractCompiler,
	safeNameIdSupplier func() string,
	constFunctionNames map[string]bool,
	scope *Scope,
	brokenIE11LocationAssign bool) *ExpressionDecomposer {
	e := &ExpressionDecomposer{
		tempNamePrefix:   "JSCompiler_temp",
		resultNamePrefix: "JSCompiler_inline_result",
	}
	jutil.CheckNotNull(compiler)
	jutil.CheckNotNull(safeNameIdSupplier)
	e.compiler = compiler
	e.astAnalyzer = compiler.GetAstAnalyzer()
	e.astFactory = compiler.CreateAstFactory()
	e.safeNameIdSupplier = safeNameIdSupplier
	e.knownConstantFunctions = constFunctionNames
	e.scope = scope
	e.brokenIE11LocationAssign = brokenIE11LocationAssign
	if compiler.HasTypeCheckingRun() && !compiler.HasOptimizationColors() {
		e.unknownType = compiler.GetTypeRegistry().GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
	} else {
		e.unknownType = nil
	}
	return e
}

// An arbitrary limit to prevent catch infinite recursion.
// Raised from 100->1000 on Jan 22, 2021
const ExpressionDecomposer_MAX_ITERATIONS = 1000

// MaybeExposeExpression performs any rewriting necessary so that the
// specified expression is MOVABLE.
//
// This method is a primary entrypoint into this class. It performs
// expression decomposition such that expression can be moved to a preceding
// statement without changing behaviour.
//
// Exposing expression generally doesn't mean that expression itself will
// moved. An expression is exposed within a larger statement if no preceding
// expression would interact with it.
//
// See CanExposeExpression.
//
// Java: ExpressionDecomposer.maybeExposeExpression(Node)
func (e *ExpressionDecomposer) MaybeExposeExpression(expression rhino.Node) {
	// If the expression needs to exposed.
	i := 0
	for ExpressionDecomposer_DecompositionType_DECOMPOSABLE == e.CanExposeExpression(expression) {
		if !e.ExposeExpression(expression) {
			// If `canExposeExpression` returned `DECOMPOSABLE` but `exposeExpression` returned false`
			// nothing was exposed so there's no point in trying again. Indicates a bug in either
			// `canExposeExpression` or `exposeExpression`.
			panic(jutil.NewIllegalStateExceptionMsg(
				"exposeExpression exposed nothing for:\n" + expression.ToStringTree()))
		}
		i++
		if i > ExpressionDecomposer_MAX_ITERATIONS {
			panic(jutil.NewIllegalStateExceptionMsg(
				"exposeExpression depth exceeded on:\n" + expression.ToStringTree()))
		}
	}
}

// ExposeExpression performs partial decomposition to get the given
// expression closer to being MOVEABLE.
//
// Returns whether any modifications furthering exposure were made to the
// expression.
//
// Java: ExpressionDecomposer.exposeExpression(Node)
func (e *ExpressionDecomposer) ExposeExpression(expression rhino.Node) bool {
	// First rewrite all optional chains containing the expression.
	// This must be done first, because the expression root may be an optional chain, and rewriting
	// it creates a new node to be the expression root.
	e.RewriteAllContainingOptionalChains(expression)
	expressionRoot := ExpressionDecomposer_FindExpressionRoot(expression)
	rhino.CheckNotNull(expressionRoot)
	jutil.CheckStateMsg1(NodeUtil_IsStatement(expressionRoot), "%s", expressionRoot)
	return e.ExposeExpressionWithRoot(expressionRoot, expression)
}

// ExposeExpressionWithRoot rewrites expressionRoot such that subExpression
// is a MOVABLE while maintaining evaluation order.
//
// IMPORTANT: This method assumes there are no optional chain parents of
// subExpression. The single-argument version of this method takes care of
// that and should be the only caller of this method.
//
// Two types of subexpressions are extracted from the source expression:
//
//  1. subexpressions with side-effects
//  2. conditional expressions that contain subExpression, which are
//     transformed into IF statements.
//
// The following terms are used:
//
//   - expressionRoot: The top-level node, before which any extracted
//     expressions should be placed.
//   - nodeWithNonconditionalParent: The node that will be extracted.
//
// expressionRoot is the root of the subtree within which to expose
// subExpression; subExpression a descendant of expressionRoot to be
// exposed. Returns whether any modifications furthering exposure were made
// to the expression.
//
// Java: ExpressionDecomposer.exposeExpression(Node, Node)
func (e *ExpressionDecomposer) ExposeExpressionWithRoot(expressionRoot rhino.Node, subExpression rhino.Node) bool {
	exposedSomething := false

	nodeWithNonconditionalParent := ExpressionDecomposer_FindNonconditionalParent(subExpression, expressionRoot)
	// Before extraction, record whether there are side-effect
	hasFollowingSideEffects := e.astAnalyzer.MayHaveSideEffects(nodeWithNonconditionalParent)

	exprInjectionPoint := ExpressionDecomposer_FindInjectionPoint(nodeWithNonconditionalParent)
	state := NewExpressionDecomposer_DecompositionState()
	state.sideEffects = hasFollowingSideEffects
	state.extractBeforeStatement = exprInjectionPoint

	// Extract expressions in the reverse order of their evaluation. This is roughly, traverse up
	// the AST extracting any preceding expressions that may have side-effects or be side-effected.
	var lastExposedSubexpression rhino.Node = rhino.NilNode
	expressionToExpose := nodeWithNonconditionalParent
	expressionParent := expressionToExpose.GetParent()
	for expressionParent != expressionRoot {
		jutil.CheckStateMsg1(
			!ExpressionDecomposer_IsConditionalOp(expressionParent) || expressionToExpose.IsFirstChildOf(expressionParent),
			"%s",
			expressionParent)

		if assign, isAssign := ast.AsAssign(expressionParent); isAssign {
			if e.IsSafeAssign(expressionParent, state.sideEffects) {
				// It is always safe to inline "foo()" for expressions such as
				// "a = b = c = foo();"
				// As the assignment is unaffected by side effect of "foo()"
				// and the names assigned-to cannot influence the state before
				// the call to foo.
				//
				// This is not true of more complex LHS values, such as
				// a.x = foo();
				// next().x = foo();
				// in these cases the checks below are necessary.
			} else if !assign.Target().Is(expressionToExpose) {
				// Alias "next()" in "next().foo"
				left := assign.Target()
				if get, isGet := left.AsNormalGet(); isGet {
					exposedSomething =
						e.DecomposeSubExpressions(get.Object().Node(), rhino.NilNode, state) || exposedSomething
				} else if _, isPat := left.AsPattern(); !isPat {
					// ArrayPattern / ObjectPattern: e.g. backoff from exposing expression `getObj()` in
					// `[getObj().propName] = ...` or `{a: getObj().propName} = ...` as the RHS must execute first.
					panic(jutil.NewIllegalStateExceptionMsg(
						"Expected a property access or destructuring pattern: " + left.Node().ToStringTree()))
				}
			}
		} else if call, isCall := ast.AsCall(expressionParent); isCall &&
			NodeUtil_IsNormalGet(call.Callee().Node()) {
			callee := call.Callee().MustNormalGet()
			if !callee.Is(expressionToExpose) {
				firstArg, _ := call.FirstArg()
				exposedSomething =
					e.DecomposeSubExpressions(firstArg.Node(), expressionToExpose, state) ||
						exposedSomething
			}

			// Now handle the call expression. We only have to do this if we arrived at decomposing this
			// call through one of the arguments, rather than the callee; otherwise the callee would
			// already be safe.
			if e.IsExpressionTreeUnsafe(callee.Node(), state.sideEffects) &&
				!callee.Object().Is(lastExposedSubexpression) {
				// Either there were preexisting side-effects, or this node has side-effects.
				state.sideEffects = true
				// Rewrite the call so "this" is preserved and continue walking up from there.
				e.RewriteCallExpression(expressionParent, state)
				exposedSomething = true
			}
		} else {
			exposedSomething =
				e.DecomposeSubExpressions(expressionParent.GetFirstChild(), expressionToExpose, state) ||
					exposedSomething
		}

		lastExposedSubexpression = expressionToExpose
		expressionToExpose = expressionParent
		expressionParent = expressionToExpose.GetParent()
	}

	// Now extract the expression that the decomposition is being performed to
	// to allow to be moved.  All expressions that need to be evaluated before
	// this have been extracted, so add the expression statement after the
	// other extracted expressions and the original statement (or replace
	// the original statement.
	if nodeWithNonconditionalParent == subExpression {
		// Don't extract the call, as that introduces an extra constant VAR
		// that will simply need to be inlined back.  It will be handled as
		// an EXPRESSION call site type.
		// Node extractedCall = extractExpression(decomposition, expressionRoot);
	} else {
		parent := nodeWithNonconditionalParent.GetParent()
		needResult := !parent.IsExprResult()
		e.ExtractConditional(nodeWithNonconditionalParent, exprInjectionPoint, needResult)
		exposedSomething = true
	}

	return exposedSomething
}

// RewriteAllContainingOptionalChains rewrites all of the optional chains
// containing the given subExpression.
//
// Java: ExpressionDecomposer.rewriteAllContainingOptionalChains(Node)
func (e *ExpressionDecomposer) RewriteAllContainingOptionalChains(subExpression rhino.Node) {
	optChainRewriterBuilder :=
		OptionalChainRewriter_Builder_(e.compiler).
			SetTmpVarNameCreator(e.GetTempConstantValueName).
			SetScope(e.scope)
	// Rewriting the chains changes the shape of the AST in a way that would interfere
	// with the simple traversal from child to parent done here, so we'll traverse
	// them all first, then rewrite them.
	//
	// Go-only: language-adapter — the ArrayDeque filled with addFirst is a
	// slice filled with append and read back to front.
	var rewriters []*OptionalChainRewriter

	for exprParent := subExpression.GetParent(); !NodeUtil_IsStatement(exprParent); exprParent = exprParent.GetParent() {
		if NodeUtil_IsEndOfFullOptChain(exprParent) {
			// We want to rewrite the outermost chain first, so the last one
			// we find is the first one we rewrite.
			rewriters = append(rewriters, optChainRewriterBuilder.Build(exprParent))
		} else if call, isCall := ast.AsCall(exprParent); isCall {
			// It is possible to make a non-optional call against an optional chain callee by applying
			// parentheses like this.
			// `(obj?.method)(arg)`
			// I don't think there's a good reason to do that, since it will cause a runtime exception
			// if the chain is ever undefined, but it is allowed, so we must handle it.
			// Fortunately the OptionalChainRewriter knows how to fix the call so it will still get
			// made with the right `this` value.
			callee := call.Callee().Node()
			if NodeUtil_IsOptChainGet(callee) {
				// By definition callee must end an optional chain, because it is the first child of a
				// non-optional parent.
				// checkState(NodeUtil.isEndOfFullOptChain(callee))
				rewriters = append(rewriters, optChainRewriterBuilder.Build(callee))
			}
		}
	}
	for i := len(rewriters) - 1; i >= 0; i-- {
		rewriters[i].Rewrite()
	}
}

// MoveExpression extracts the specified expression from its parent
// expression.
//
// See CanExposeExpression.
//
// Java: ExpressionDecomposer.moveExpression(Node)
func (e *ExpressionDecomposer) MoveExpression(expression rhino.Node) {
	// TODO(johnlenz): This is not currently used by the function inliner,
	// as moving the call out of the expression before the actual function call
	// causes additional variables to be introduced.  As the variable
	// inliner is improved, this might be a viable option.

	resultName := e.GetResultValueName()
	injectionPoint := ExpressionDecomposer_FindInjectionPoint(expression)
	rhino.CheckNotNull(injectionPoint)
	injectionPointParent := injectionPoint.GetParent()
	rhino.CheckNotNull(injectionPointParent)
	jutil.CheckState(NodeUtil_IsStatementBlock(injectionPointParent))

	// Replace the expression with a reference to the new name.
	expression.ReplaceWith(rhino.IR_Name(resultName).CopyTypeFrom(expression))

	// Re-add the expression at the appropriate place.
	newExpressionRoot := NodeUtil_NewVarNode(resultName, expression)
	newExpressionRoot.GetFirstChild().CopyTypeFrom(expression)
	newExpressionRoot.InsertBefore(injectionPoint)

	e.compiler.ReportChangeToEnclosingScope(injectionPointParent)
}

// ExpressionDecomposer_FindNonconditionalParent returns the enclosing
// expression to decompose
//
// The intention is to indicate the top-most node that could be rewritten as
// an if-statement in order to better expose subExpression for inlining.
//
// Examples:
//
//	a = (x() && y()) && subExpression; // result is (x() && y()) && subExpression
//	a = x() && (y() && subExpression); // result is x() && (y() && subExpression)
//	a = (x() && subExpression) && y(); // result is x() && subExpression
//	a = x() && (subExpression && y()); // result is x() && (subExpression && y())
//	a = (subExpression && x()) && y(); // result is subExpression
//	a = subExpression && (x() && y()); // result is subExpression
//
// When subExpression is contained within an optional chain, we want to
// treat everything after a `?.` up until the next `?.` as a single
// conditional operation.
//
// Examples:
//
//	a = subExpression.x?.y.z();          // result is subExpression
//	a = x()?.[subExpression].y;          // result is x()?.[subExpression].y
//	a = x()?.y.z?.p(subExpression).q?.r; // result is x()?.y.z?.p(subExpression).q
//	a
//
// subExpression is the expression to consider entire chains; expressionRoot
// a node containing subExpression. The returned node will be a descendent
// of this one.
//
// Java: ExpressionDecomposer.findNonconditionalParent(Node, Node)
func ExpressionDecomposer_FindNonconditionalParent(subExpression rhino.Node, expressionRoot rhino.Node) rhino.Node {
	result := subExpression

	for child, parent := subExpression, subExpression.GetParent(); parent != expressionRoot; child, parent = parent, parent.GetParent() {
		if ExpressionDecomposer_IsConditionalOp(parent) && !child.IsFirstChildOf(parent) {
			// subExpression is not part of the first child (which is always executed), so
			// parent decides whether subExpression will be executed or not
			result = parent
		}
	}
	if NodeUtil_IsOptChainNode(result) {
		// the loop above may have left result pointing into the middle of an optional chain for
		// a case like this.
		// `x?.y.z(subExpression).p.q?.r.s`
		// result is currently `x?.y.z(subExpression)`, but we want it to be the full sub-chain
		// containing subExpression
		// `x?.y.z(subExpression).p.q`
		result = NodeUtil_GetEndOfOptChainSegment(result)
	}

	return result
}

// ExpressionDecomposer_DecompositionState is a simple class to track two
// things: - whether side effects have been seen. - the last statement
// inserted
type ExpressionDecomposer_DecompositionState struct {
	sideEffects            bool
	extractBeforeStatement rhino.Node
}

// Java: ExpressionDecomposer.DecompositionState()
func NewExpressionDecomposer_DecompositionState() *ExpressionDecomposer_DecompositionState {
	return &ExpressionDecomposer_DecompositionState{}
}

// Java: ExpressionDecomposer.DecompositionState.toString()
func (s *ExpressionDecomposer_DecompositionState) String() string {
	return "DecompositionState{sideEffects=" + jutil.ToJavaString(s.sideEffects) +
		", extractBeforeStatement=" + jutil.ToJavaString(s.extractBeforeStatement) + "}"
}

// DecomposeSubExpressions: n is the node with which to start iterating;
// stopNode a node after which to stop iterating. Returns whether any
// modifications furthering decomposition were made to n.
//
// Java: ExpressionDecomposer.decomposeSubExpressions(Node, Node, ExpressionDecomposer.DecompositionState)
func (e *ExpressionDecomposer) DecomposeSubExpressions(
	n rhino.Node, stopNode rhino.Node, state *ExpressionDecomposer_DecompositionState) bool {
	if n.IsNil() || n == stopNode {
		return false
	}

	decomposedSomething := false

	// Decompose the children in reverse evaluation order. This simplifies determining if any of
	// the children following have side-effects. If they do we need to be more aggressive about
	// removing values from the expression. Reverse order also maintains evaluation order as each
	// extracted statemented is inserted on top of the others.
	decomposedSomething =
		e.DecomposeSubExpressions(n.GetNext(), stopNode, state) || decomposedSomething

	// Now this node.

	if cp, isComputedProp := ast.AsComputedProp(n); isComputedProp {
		// TODO(b/111621528): Delete when fixed.
		// If the prop is computed we have to fork the decomposition between the key and value. This
		// is because we can't move the property assignment itself; COMPUTED_PROP must remain a
		// child of OBJECTLIT for example.
		//
		// We decompose the value of the prop first because decomposition is in reverse order of
		// evaluation.
		if val, hasVal := cp.Value(); hasVal {
			decomposedSomething =
				e.DecomposeSubExpressions(val.Node(), stopNode, state) || decomposedSomething
		}
		n = cp.Key().Node()
	} else if sk, isStringKey := ast.AsStringKey(n); isStringKey {
		n = sk.Value().Node()
	} else if md, isMethodDef := ast.AsMethodDef(n); isMethodDef {
		n = md.Function().Node()
	} else if sub, isSub := ast.AsTemplateLitSub(n); isSub {
		// A template literal substitution expression like ${f()} is represented in the AST as
		//   TEMPLATELIT_SUB
		//     CALL
		//       NAME f
		// The TEMPLATELIT_SUB node is not actually an expression and can't be extracted, but we may
		// need to extract the expression inside of it.
		n = sub.Expr().Node()
	} else if n.IsSpread() {
		// SPREADs aren't expression but they can still be extracted using temp variables.
		//
		// Because object-spread can trigger getters we assume all spreads have side-effects.
		// TODO(nickreid): Use `assumeGettersArePure` here. It would have been a pain to pipe it down
		// here and write all the tests. Since there are very few cases, and it doesn't affect code
		// removal, we didn't bother initially. Everything always works one way.
	} else if !rhino.IR_MayBeExpression(n) {
		// If n is not an expression then it can't be extracted. For example if n is the destructuring
		// pattern on the left side of a VAR statement:
		//   var {pattern} = rhs();
		// See test case: testExposeExpression18
		return decomposedSomething
	}

	// TODO(johnlenz): Move "safety" code to a shared class.
	if e.IsExpressionTreeUnsafe(n, state.sideEffects) {
		// Either there were preexisting side-effects, or this node has side-effects.
		state.sideEffects = true
		state.extractBeforeStatement = e.ExtractExpression(n, state.extractBeforeStatement)
		decomposedSomething = true
	}
	return decomposedSomething
}

// ExtractConditional: expr is the conditional expression to extract;
// injectionPoint the node before which the extracted expression would be
// injected; needResult whether the result of the expression is required.
// Returns the node that contains the logic of the expression after
// extraction.
//
// Java: ExpressionDecomposer.extractConditional(Node, Node, boolean)
func (e *ExpressionDecomposer) ExtractConditional(expr rhino.Node, injectionPoint rhino.Node, needResult bool) rhino.Node {
	parent := expr.GetParent()
	tempName := e.GetTempValueName()

	// Transform the conditional to an IF statement.
	var cond rhino.Node = rhino.NilNode
	trueExpr := e.astFactory.CreateBlock().Srcref(expr)
	falseExpr := e.astFactory.CreateBlock().Srcref(expr)
	switch expr.GetToken() {
	case rhino.Token_HOOK:
		// a = x?y:z --> if (x) {a=y} else {a=z}
		hook := ast.MustHook(expr)
		condNode, thenNode, elseNode := hook.Cond().Node(), hook.Then().Node(), hook.Else().Node()
		expr.DetachChildren()
		cond = condNode
		trueExpr.AddChildToFront(
			e.astFactory.ExprResult(ExpressionDecomposer_BuildResultExpression(thenNode, needResult, tempName)))
		falseExpr.AddChildToFront(
			e.astFactory.ExprResult(ExpressionDecomposer_BuildResultExpression(elseNode, needResult, tempName)))
	case rhino.Token_AND:
		// a = x&&y --> if (a=x) {a=y} else {}
		bin := ast.MustBinary(expr)
		left, right := bin.Left().Node(), bin.Right().Node()
		expr.DetachChildren()
		cond = ExpressionDecomposer_BuildResultExpression(left, needResult, tempName)
		trueExpr.AddChildToFront(
			e.astFactory.ExprResult(ExpressionDecomposer_BuildResultExpression(right, needResult, tempName)))
	case rhino.Token_OR:
		// a = x||y --> if (a=x) {} else {a=y}
		bin := ast.MustBinary(expr)
		left, right := bin.Left().Node(), bin.Right().Node()
		expr.DetachChildren()
		cond = ExpressionDecomposer_BuildResultExpression(left, needResult, tempName)
		falseExpr.AddChildToFront(
			e.astFactory.ExprResult(ExpressionDecomposer_BuildResultExpression(right, needResult, tempName)))
	case rhino.Token_COALESCE:
		// a = x ?? y --> if ((temp=x)!=null) {a=temp} else {a=y}
		bin := ast.MustBinary(expr)
		left, right := bin.Left().Node(), bin.Right().Node()
		expr.DetachChildren()
		tempNameAssign := e.GetTempValueName()
		tempVarNodeAssign :=
			e.astFactory.CreateSingleVarNameDeclaration(tempNameAssign).SrcrefTreeIfMissing(expr)
		tempVarNodeAssign.InsertBefore(injectionPoint)

		assignLhs := ExpressionDecomposer_BuildResultExpression(left, true, tempNameAssign)
		nullNode := e.astFactory.CreateNull().Srcref(expr)
		cond = e.astFactory.CreateNe(assignLhs, nullNode).Srcref(expr)
		trueExpr.AddChildToFront(
			e.astFactory.ExprResult(
				ExpressionDecomposer_BuildResultExpression(
					e.astFactory.CreateName(tempNameAssign, AstFactory_Type_(left)).Srcref(expr),
					needResult,
					tempName)))
		falseExpr.AddChildToFront(
			e.astFactory.ExprResult(ExpressionDecomposer_BuildResultExpression(right, needResult, tempName)))
	default:
		// With a valid tree we should never get here.
		panic(jutil.NewIllegalStateExceptionMsg("Unexpected expression: " + expr.String()))
	}

	var ifNode rhino.Node
	if falseExpr.HasChildren() {
		ifNode = e.astFactory.CreateIfWithElse(cond, trueExpr, falseExpr)
	} else {
		ifNode = e.astFactory.CreateIf(cond, trueExpr)
	}
	ifNode.SrcrefIfMissing(expr)

	if needResult {
		tempVarNode :=
			e.astFactory.CreateSingleVarNameDeclaration(tempName).SrcrefTreeIfMissing(expr)
		tempVarNode.InsertBefore(injectionPoint)
		ifNode.InsertAfter(tempVarNode)

		// Replace the expression with the temporary name.
		replacementValueNode := rhino.IR_Name(tempName).CopyTypeFrom(expr)
		expr.ReplaceWith(replacementValueNode)
	} else {
		// Only conditionals that are the direct child of an expression statement
		// don't need results, for those simply replace the expression statement.
		jutil.CheckArgument(parent.IsExprResult())
		parent.ReplaceWith(ifNode)
	}

	return ifNode
}

// ExpressionDecomposer_BuildResultExpression creates an expression tree for
// an expression.
//
// If the result of the expression is needed, then:
//
//	ASSIGN
//	  tempName
//	  expr
//
// otherwise, simply: `expr`
//
// Java: ExpressionDecomposer.buildResultExpression(Node, boolean, String)
func ExpressionDecomposer_BuildResultExpression(expr rhino.Node, needResult bool, tempName string) rhino.Node {
	if needResult {
		return rhino.IR_Assign(rhino.IR_Name(tempName).CopyTypeFrom(expr), expr).
			CopyTypeFrom(expr).
			SrcrefTree(expr)
	} else {
		return expr
	}
}

// Java: ExpressionDecomposer.isConstantNameNode(Node)
func (e *ExpressionDecomposer) IsConstantNameNode(n rhino.Node) bool {
	// Non-constant names values may have been changed.
	return n.IsName() &&
		(NodeUtil_IsConstantVar(n, e.scope) || e.knownConstantFunctions[n.GetString()])
}

// ExtractExpression: expr is the expression to extract; injectionPoint the
// node before which to added the extracted expression. Returns the
// extracted statement node.
//
// Java: ExpressionDecomposer.extractExpression(Node, Node)
func (e *ExpressionDecomposer) ExtractExpression(expr rhino.Node, injectionPoint rhino.Node) rhino.Node {
	// Since all instances of logical assignment ops are normalized into an expression that
	// separates the logical operation from the assignment, logical assignment ops are
	// never seen here, and its logical operation part is instead handled by extractConditional().
	jutil.CheckStateMsg1(!NodeUtil_IsLogicalAssignmentOp(expr), "%s", expr)
	parent := expr.GetParent()

	isLhsOfAssignOp := ast.RoleIn(expr, parent) == ast.RoleAssignOpTarget

	var firstExtractedNode rhino.Node = rhino.NilNode

	// Expressions on the LHS of an assignment-op must have any possible
	// side-effects extracted as the value must be duplicated:
	//    next().foo += 2;
	// becomes:
	//    var t1 = next();
	//    t1.foo = t1.foo + 2;
	if isLhsOfAssignOp && NodeUtil_IsNormalGet(expr) {
		var next rhino.Node
		for n := expr.GetFirstChild(); !n.IsNil(); n = next {
			// Cache next before extracting n because n is mutated or replaced within the loop body
			// during expression extraction.
			next = n.GetNext()
			if !n.IsStringLit() && !e.IsConstantNameNode(n) {
				extractedNode := e.ExtractExpression(n, injectionPoint)
				if firstExtractedNode.IsNil() {
					firstExtractedNode = extractedNode
				}
			}
		}
	}

	// The temp is known to be constant.
	tempName := e.GetTempConstantValueName()
	replacementValueNode := rhino.IR_Name(tempName).CopyTypeFrom(expr).Srcref(expr)

	var tempNameValue rhino.Node

	// If it is ASSIGN_XXX and not a logical assignment operator keep the
	// assignment in place and extract the original value of the LHS operand.
	if isLhsOfAssignOp {
		jutil.CheckStateMsg1(expr.IsName() || NodeUtil_IsNormalGet(expr), "%s", expr)
		assignOp := ast.MustAssignOp(parent)
		// Transform "x += 2" into "x = temp + 2"
		opNode :=
			rhino.NewNode(NodeUtil_GetOpFromAssignmentOp(parent)).
				CopyTypeFrom(parent).
				SrcrefIfMissing(parent)

		rightOperand := assignOp.Value().Node()

		parent.SetToken(rhino.Token_ASSIGN)
		rightOperand.ReplaceWith(opNode)
		opNode.AddChildToFront(replacementValueNode)
		opNode.AddChildToBack(rightOperand)

		// The original expression is still being used, so make a clone.
		tempNameValue = expr.CloneTree()
	} else if expr.IsSpread() {
		// We need to treat spreads differently because unlike other expressions, they can't be
		// directly assigned to new variables. Instead we wrap them in a literal.
		//
		// We make sure to do `var tmp = [...fn()];` rather than `var tmp = fn()` because the
		// execution of a spread on an arbitrary iterable/object can both have side-effects and be
		// side-effected. However, once done we are then sure that spreading `tmp` is isolated.
		var spreadInner rhino.Node
		if sp, ok := ast.AsIterSpread(expr); ok {
			spreadInner = sp.Expr().Node()
		} else {
			spreadInner = ast.MustObjectSpread(expr).Expr().Node()
		}

		// Replace the expression with the spread for the temporary name.
		spreadCopy := expr.CloneNode()
		spreadCopy.AddChildToBack(replacementValueNode)
		expr.ReplaceWith(spreadCopy)

		// Move the original node into a legal context.
		switch parent.GetToken() {
		case rhino.Token_ARRAYLIT, rhino.Token_CALL, rhino.Token_NEW:
			tempNameValue = e.astFactory.CreateArraylit(expr).Srcref(spreadInner)
		case rhino.Token_OBJECTLIT:
			tempNameValue = e.astFactory.CreateObjectLit(expr).Srcref(spreadInner)
		default:
			panic(jutil.NewIllegalStateExceptionMsg(
				"Unexpected parent of SPREAD:" + parent.ToStringTree()))
		}
	} else {
		// Replace the expression with the temporary name.
		expr.ReplaceWith(replacementValueNode)

		// Keep the original node so that CALL expressions can still be found
		// and inlined properly.
		tempNameValue = expr
	}

	// Re-add the expression in the declaration of the temporary name.
	tempVarNode := NodeUtil_NewVarNode(tempName, tempNameValue)
	tempVarNode.GetFirstChild().CopyTypeFrom(tempNameValue)
	tempVarNode.GetFirstChild().SetInferredConstantVar(true)
	containingHoistScope := e.scope.GetClosestHoistScope()
	containingHoistScope.Declare(tempName, tempVarNode.GetFirstChild() /* input= */, nil)

	tempVarNode.InsertBefore(injectionPoint)

	if firstExtractedNode.IsNil() {
		firstExtractedNode = tempVarNode
	}

	jutil.CheckState(firstExtractedNode.IsVar())
	return firstExtractedNode
}

// RewriteCallExpression rewrites the call so "this" is preserved.
//
//	a.b(c);
//
// becomes:
//
//	var temp1 = a; var temp0 = temp1.b;
//	temp0.call(temp1,c);
//
// Java: ExpressionDecomposer.rewriteCallExpression(Node, ExpressionDecomposer.DecompositionState)
func (e *ExpressionDecomposer) RewriteCallExpression(callNode rhino.Node, state *ExpressionDecomposer_DecompositionState) {
	call, isCall := ast.AsCall(callNode)
	jutil.CheckArgumentMsg1(isCall, "%s", callNode)
	calleeGet, isGet := call.Callee().AsNormalGet()
	jutil.CheckArgumentMsg1(isGet, "%s", call.Callee().Node())
	first := calleeGet.Node()

	// Find the type of (fn expression).call
	var fnCallType rhino.JSType = nil
	if e.astFactory.IsAddingTypes() {
		fnType := first.GetJSType()
		if fnType.IsFunctionType() {
			fnCallType = fnType.ToMaybeFunctionType().GetPropertyType("call")
		} else {
			fnCallType = e.unknownType
		}
	}

	// Extracts the expression representing the function to call. For example:
	//   "a['b'].c" from "a['b'].c()"
	getVarNode := e.ExtractExpression(first, state.extractBeforeStatement)
	state.extractBeforeStatement = getVarNode
	varDecl := ast.MustNameDeclaration(getVarNode).MustOnlyDeclarator().MustName()
	functionNameNode := varDecl.Node().CloneNode()

	if call.IsFreeCall() {
		// For a free call, we don't need to extract the receiver
		call.SetCallee(ast.MustExpr(functionNameNode))
		return
	}

	// Extracts the object reference to be used as "this". For example:
	//   "a['b']" from "a['b'].c"
	initExpr, _ := varDecl.Init()
	getExpr, ok := initExpr.AsNormalGet()
	jutil.CheckArgumentMsg1(ok, "%s", initExpr.Node())
	origThisValue := getExpr.Object().Node()
	var receiverNode rhino.Node

	if origThisValue.IsThis() {
		// No need to create a variable for `this`, just clone it.
		receiverNode = origThisValue.CloneNode()
	} else if origThisValue.IsSuper() {
		// Original callee was like `super.prop(args)`.
		// The correct way to call the value `super.prop` from a temporary variable is
		// `tmpVar.call(this, args)`, so just create a `this` here.
		receiverNode = e.astFactory.CreateThis(AstFactory_Type_(origThisValue)).Srcref(origThisValue)
	} else {
		thisVarNode := e.ExtractExpression(origThisValue, state.extractBeforeStatement)
		state.extractBeforeStatement = thisVarNode
		thisDecl := ast.MustNameDeclaration(thisVarNode).MustOnlyDeclarator().MustName()
		receiverNode = thisDecl.Node().CloneNode()
	}

	// CALL
	//   GETPROP "call"
	//     functionName
	//   thisName
	//   original-parameter1
	//   original-parameter2
	//   ...

	// Reuse the existing CALL node instead of creating a new one to avoid breaking InlineFunction's
	// bookkeeping. See b/124253050.
	callNode.RemoveFirstChild()
	callNode.AddChildToFront(receiverNode)
	callNode.AddChildToFront(
		e.astFactory.
			CreateGetProp(functionNameNode, "call", AstFactory_Type_OfJSTypeAndColor(fnCallType, colors.StandardColors_TOP_OBJECT)).
			SrcrefTreeIfMissing(callNode))
	call.SetFreeCall(false)
}

// SetTempNamePrefix allows the temp name to be overridden to make tests
// more readable.
//
// Java: ExpressionDecomposer.setTempNamePrefix(String)
func (e *ExpressionDecomposer) SetTempNamePrefix(prefix string) {
	e.tempNamePrefix = prefix
}

// GetTempValueName creates a unique temp name.
//
// Java: ExpressionDecomposer.getTempValueName()
func (e *ExpressionDecomposer) GetTempValueName() string {
	return e.tempNamePrefix + MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SEPARATOR + e.safeNameIdSupplier()
}

// SetResultNamePrefix allows the temp name to be overridden to make tests
// more readable.
//
// Java: ExpressionDecomposer.setResultNamePrefix(String)
func (e *ExpressionDecomposer) SetResultNamePrefix(prefix string) {
	e.resultNamePrefix = prefix
}

// GetResultValueName creates a unique name for call results.
//
// Java: ExpressionDecomposer.getResultValueName()
func (e *ExpressionDecomposer) GetResultValueName() string {
	return e.resultNamePrefix + MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SEPARATOR + e.safeNameIdSupplier()
}

// GetTempConstantValueName creates a constant unique temp name.
//
// Java: ExpressionDecomposer.getTempConstantValueName()
func (e *ExpressionDecomposer) GetTempConstantValueName() string {
	return e.tempNamePrefix +
		"_const" +
		MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SEPARATOR +
		e.safeNameIdSupplier()
}

// Java: ExpressionDecomposer.isTempConstantValueName(Node)
func (e *ExpressionDecomposer) IsTempConstantValueName(name rhino.Node) bool {
	return name.IsName() &&
		strings.HasPrefix(
			name.GetString(),
			e.tempNamePrefix+"_const"+MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SEPARATOR)
}

// ExpressionDecomposer_FindInjectionPoint returns, for the subExpression,
// the nearest statement Node before which it can be inlined. Null if no
// such location can be found.
//
// Java: ExpressionDecomposer.findInjectionPoint(Node)
func ExpressionDecomposer_FindInjectionPoint(subExpression rhino.Node) rhino.Node {
	expressionRoot := ExpressionDecomposer_FindExpressionRoot(subExpression)
	rhino.CheckNotNull(expressionRoot)

	injectionPoint := expressionRoot

	parent := injectionPoint.GetParent()
	for parent.IsLabel() {
		injectionPoint = parent
		parent = injectionPoint.GetParent()
	}

	jutil.CheckStateMsg1(NodeUtil_IsStatementBlock(parent), "%s", parent)
	return injectionPoint
}

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
	for current := range child.GetAncestors().All() {
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
	for parent := range child.GetAncestors().All() {
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
