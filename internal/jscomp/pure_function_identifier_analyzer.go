// Port of com/google/javascript/jscomp/PureFunctionIdentifier.java @ closure-compiler v20261021

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
	"slices"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// PureFunctionIdentifier_FunctionBodyAnalyzer inspects function bodies for
// side effects and applies them to the associated AmbiguatedFunctionSummary.
//
// This callback also fills allFunctionCalls
type PureFunctionIdentifier_FunctionBodyAnalyzer struct {
	outer *PureFunctionIdentifier

	// Preloaded with an entry to represent the global scope.
	//
	// Go-only: language-adapter — the ArrayDeque is only used at its end
	// (getLast/addLast/removeLast), so it is a slice. Go-only: performance —
	// the slots past the end keep the entries popped last, and enterScope
	// reuses one unless a SideEffectPropagation holds it (heldByEdge).
	functionScopeStack []*PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry
}

// Java: PureFunctionIdentifier.FunctionBodyAnalyzer()
func NewPureFunctionIdentifier_FunctionBodyAnalyzer(outer *PureFunctionIdentifier) *PureFunctionIdentifier_FunctionBodyAnalyzer {
	b := &PureFunctionIdentifier_FunctionBodyAnalyzer{
		outer:              outer,
		functionScopeStack: make([]*PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry, 1, 8),
	}
	b.functionScopeStack[0] = NewPureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry(rhino.NilNode, false)
	return b
}

// functionScopeStackGetLast is ArrayDeque.getLast().
func (b *PureFunctionIdentifier_FunctionBodyAnalyzer) functionScopeStackGetLast() *PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry {
	if len(b.functionScopeStack) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	return b.functionScopeStack[len(b.functionScopeStack)-1]
}

// PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry is a single
// function literal definition and associated side-effect-related
// information.
//
// We define separate objects for function literal to simplify handling of
// nested function literal definitions - we push a new entry onto the stack
// when entering a new function literal, and pop the entry when leaving the
// function literal.
type PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry struct {
	root                rhino.Node
	skiplistedVars      []*Var
	skiplistedNameNodes []rhino.Node
	taintedVars         []*Var
	// Whether this function was marked as artificially pure, i.e. had a
	// `@nosideeffects` annotation. This tells the compiler to ignore side
	// effects from within the function body when deciding whether the given
	// function as a whole is pure. Note that within the function body, we
	// still record side effects on individual variables accesses; so e.g.
	// given
	//   /** @nosideeffects */
	//   function foo() { console.log(1); }
	//   foo();
	// `foo()` will be marked as pure, but `console.log` will still be marked
	// as impure. Also note that this doesn't extend to nested function
	// definitions: a nested function does not inherit artificial purity from
	// its enclosing function, which is why we need to track this on the
	// function stack.
	isArtificiallyPure bool
	catchDepth         int // The number of try-catch blocks around the current node.

	// Go-only: performance — whether a SideEffectPropagation of an invocation
	// in this function holds the entry (as its callerEntry), so that it may not
	// be reused for another function.
	heldByEdge bool
}

// Java: PureFunctionIdentifier.FunctionBodyAnalyzer.FunctionStackEntry(Node, boolean)
func NewPureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry(root rhino.Node, isArtificiallyPure bool) *PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry {
	return &PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry{
		catchDepth:         0,
		root:               root,
		isArtificiallyPure: isArtificiallyPure,
	}
}

func (e *PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry) addSkiplistedVar(v *Var) {
	if slices.Contains(e.skiplistedVars, v) {
		return
	}
	e.skiplistedVars = append(e.skiplistedVars, v)
}

func (e *PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry) isSkiplistedVar(v *Var) bool {
	if slices.Contains(e.skiplistedVars, v) {
		return true
	}
	for _, node := range e.skiplistedNameNodes {
		if node == v.GetNameNode() {
			return true
		}
		if node.GetStringID() == v.nameID && node.GetParent().IsVar() && v.scope.IsHoistScope() {
			return true
		}
	}
	return false
}

func (e *PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry) addTaintedVar(v *Var) {
	if slices.Contains(e.taintedVars, v) {
		return
	}
	e.taintedVars = append(e.taintedVars, v)
}

func (e *PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry) isTaintedVar(v *Var) bool {
	return slices.Contains(e.taintedVars, v)
}

// IsLocalVar reports whether v, as seen from scope, is a variable of this
// function that only holds values local to the current invocation, so that
// mutating them is invisible to the callers.
//
// Java: PureFunctionIdentifier.FunctionBodyAnalyzer.FunctionStackEntry.isLocalVar(Var, Scope)
func (e *PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry) IsLocalVar(v *Var, scope *Scope) bool {
	return PureFunctionIdentifier_FunctionBodyAnalyzer_IsVarDeclaredInSameContainerScope(v, scope) &&
		!e.isSkiplistedVar(v) &&
		// Parameters and catch values can come from other scopes.
		!v.IsParam() &&
		!v.IsCatch()
}

// Java: PureFunctionIdentifier.FunctionBodyAnalyzer.shouldTraverse(NodeTraversal, Node, Node)
func (b *PureFunctionIdentifier_FunctionBodyAnalyzer) ShouldTraverse(traversal *NodeTraversal, node rhino.Node, parent rhino.Node) bool {
	b.AddToCatchDepthIfTryBlock(node, 1)
	return true
}

// Java: PureFunctionIdentifier.FunctionBodyAnalyzer.visit(NodeTraversal, Node, Node)
func (b *PureFunctionIdentifier_FunctionBodyAnalyzer) Visit(traversal *NodeTraversal, node rhino.Node, parent rhino.Node) {
	b.AddToCatchDepthIfTryBlock(node, -1)

	if NodeUtil_IsInvocation(node) {
		// We collect these after filtering for side-effects because there's no point re-processing
		// a known pure call. This analysis is run multiple times, but no optimization will make a
		// pure function impure.
		b.outer.allFunctionCalls = append(b.outer.allFunctionCalls, node)
	}

	enclosingFunction := b.functionScopeStackGetLast()
	root := enclosingFunction.root
	if root.IsNil() || enclosingFunction.isArtificiallyPure {
		// Within artificially pure function literal bodies, we do not ever need to call
		// updateSideEffectsForNode because we intentionally ignore any side effects, and still
		// consider the function as a whole to be pure.
		// Note: this doesn't mean that individual invocations within the function body are always
		// considered pure - those will still be marked in a later stage of PureFunctionIdentifier.
		// We just don't propagate any impurity to the enclosing function literal.
		return
	}

	for _, summary := range b.outer.summariesForAllNamesOfFunctionByNode[root] {
		b.UpdateSideEffectsForNode(jutil.CheckNotNull(summary), traversal, node)
	}
}

// UpdateSideEffectsForNode updates the side effects of summary based on a
// given node.
//
// This node should be known to (possibly have) side effects. This method
// does not check if the node (possibly) has side effects.
//
// Java: PureFunctionIdentifier.FunctionBodyAnalyzer.updateSideEffectsForNode(AmbiguatedFunctionSummary, NodeTraversal, Node)
func (b *PureFunctionIdentifier_FunctionBodyAnalyzer) UpdateSideEffectsForNode(
	encloserSummary *PureFunctionIdentifier_AmbiguatedFunctionSummary, traversal *NodeTraversal, node rhino.Node) {
	if encloserSummary.MutatesGlobalState() {
		return // Functions with MUTATES_GLOBAL_STATE already have all side-effects set.
	}

	switch node.GetToken() {
	case rhino.Token_ASSIGN:
		// e.g.
		// lhs = rhs;
		// ({x, y} = object);
		lhs := ast.MustAssign(node).Target()
		// Consider destructured properties or values to be nonlocal.
		if _, isPat := lhs.AsPattern(); isPat {
			NodeUtil_VisitLhsNodesInNode(
				node,
				func(lhsNode rhino.Node) {
					b.visitLhsNode(encloserSummary, traversal, lhsNode, PureFunctionIdentifier_RHS_IS_NEVER_LOCAL)
				})
		} else {
			b.visitLhsNode(encloserSummary, traversal, lhs.Node(), PureFunctionIdentifier_FIND_RHS_AND_CHECK_FOR_LOCAL_VALUE)
		}

	case rhino.Token_INC, rhino.Token_DEC: // e.g. x++;
		b.visitLhsNode(
			encloserSummary,
			traversal,
			ast.MustIncDec(node).Target().Node(),
			// The value assigned by a unary op is always local.
			PureFunctionIdentifier_RHS_IS_ALWAYS_LOCAL)

	case rhino.Token_DELPROP:
		b.visitLhsNode(
			encloserSummary,
			traversal,
			ast.MustUnary(node).Operand().Node(),
			// The value assigned by a unary op is always local.
			PureFunctionIdentifier_RHS_IS_ALWAYS_LOCAL)

	case rhino.Token_FOR_AWAIT_OF:
		b.DeprecatedSetSideEffectsForControlLoss(encloserSummary) // Control is lost during await.
		// Fall through.
		fallthrough
	case rhino.Token_FOR_OF:
		// e.g.
		// for (const {prop1, prop2} of iterable) {...}
		// for ({prop1: x.p1, prop2: x.p2} of iterable) {...}
		NodeUtil_VisitLhsNodesInNode(
			node,
			func(lhsNode rhino.Node) {
				b.visitLhsNode(
					encloserSummary,
					traversal,
					lhsNode,
					// The RHS of a for-of must always be an iterable, making it a container, so
					// we can't
					// consider its contents to be local
					PureFunctionIdentifier_RHS_IS_NEVER_LOCAL)
			})
		b.CheckIteratesImpureIterable(node, encloserSummary)

	case rhino.Token_FOR_IN:
		// e.g.
		// for (prop in obj) {...}
		// Also this, though not very useful or readable.
		// for ([char1, char2, ...x.rest] in obj) {...}
		NodeUtil_VisitLhsNodesInNode(
			node,
			func(lhsNode rhino.Node) {
				b.visitLhsNode(
					encloserSummary,
					traversal,
					lhsNode,
					// A for-in always assigns a string, which is a local value by definition.
					PureFunctionIdentifier_RHS_IS_ALWAYS_LOCAL)
			})

	case rhino.Token_OPTCHAIN_CALL, rhino.Token_CALL, rhino.Token_NEW, rhino.Token_TAGGED_TEMPLATELIT:
		b.VisitCall(encloserSummary, traversal.GetScope(), node)

	case rhino.Token_DESTRUCTURING_LHS:
		if NodeUtil_IsAnyFor(node.GetParent()) {
			// This case is handled when visiting the enclosing for loop.
			break
		}
		// Assume the value assigned to each item is potentially global state. This is overly
		// conservative but necessary because in the common case the rhs is not a literal.
		NodeUtil_VisitLhsNodesInNode(
			node.GetParent(),
			func(lhsNode rhino.Node) {
				b.visitLhsNode(encloserSummary, traversal, lhsNode, PureFunctionIdentifier_RHS_IS_NEVER_LOCAL)
			})

	case rhino.Token_NAME:
		// Local variable declarations are not a side-effect, but we do want to track them.
		if _, ok := ast.AsNameDeclaration(node.GetParent()); ok {
			// Assignment to local, if the value isn't a safe local value,
			// new object creation or literal or known primitive result
			// value, add it to the local skiplist.
			if value, hasInit := ast.MustName(node).Init(); hasInit && !NodeUtil_EvaluatesToLocalValue(value.Node()) {
				e := b.functionScopeStackGetLast()
				e.skiplistedNameNodes = append(e.skiplistedNameNodes, node)
			}
		}

	case rhino.Token_THROW:
		b.RecordThrowsBasedOnContext(encloserSummary)

	case rhino.Token_YIELD:
		b.CheckIteratesImpureIterable(node, encloserSummary) // `yield*` triggers iteration.
		// 'yield' throws if the caller calls `.throw` on the generator object.
		b.DeprecatedSetSideEffectsForControlLoss(encloserSummary)

	case rhino.Token_AWAIT:
		// 'await' throws if the promise it's waiting on is rejected.
		b.DeprecatedSetSideEffectsForControlLoss(encloserSummary)

	case rhino.Token_OBJECT_REST, rhino.Token_OBJECT_SPREAD:
		if !b.outer.assumeGettersArePure {
			// May trigger a getter.
			encloserSummary.SetMutatesGlobalStateAndAllOtherFlags(node)
		}

	case rhino.Token_ITER_REST, rhino.Token_ITER_SPREAD:
		b.CheckIteratesImpureIterable(node, encloserSummary)

	case rhino.Token_STRING_KEY:
		if _, ok := ast.AsObjectPattern(node.GetParent()); ok {
			// This is an l-value STRING_KEY.
			// Assumption: GETELEM (via a COMPUTED_PROP) is never side-effectful.
			if b.outer.GetPropertyKind(ast.MustStringKey(node).Key()).HasGetter() {
				encloserSummary.SetMutatesGlobalStateAndAllOtherFlags(node)
			}
		}

	case rhino.Token_OPTCHAIN_GETPROP, rhino.Token_GETPROP:
		// Assumption: GETELEM and OPTCHAIN_GETELEM are never side-effectful.
		if b.outer.GetPropertyKind(ast.MustNormalOrOptChainGetProp(node).Property()).HasGetterOrSetter() {
			encloserSummary.SetMutatesGlobalStateAndAllOtherFlags(node)
		}

	case rhino.Token_DYNAMIC_IMPORT:
		// Modules may be imported for side-effects only. This is frequently
		// a pattern used to load polyfills.
		encloserSummary.SetMutatesGlobalStateAndAllOtherFlags(node)

	default:
		if NodeUtil_IsCompoundAssignmentOp(node) {
			// e.g.
			// x += 3;
			b.visitLhsNode(
				encloserSummary,
				traversal,
				ast.MustAssignOp(node).Target().Node(),
				// The update assignments (e.g. `+=) always assign primitive, and therefore local,
				// values.
				PureFunctionIdentifier_RHS_IS_ALWAYS_LOCAL)
			break
		}

		if b.outer.astAnalyzer.NodeTypeMayHaveSideEffects(node) {
			panic(jutil.NewIllegalArgumentExceptionMsg("Unhandled side effect node type " + node.String()))
		}
	}
}

// CheckIteratesImpureIterable inspects node for impure iteration and
// assigns the appropriate side-effects to encloserSummary if so.
//
// Java: PureFunctionIdentifier.FunctionBodyAnalyzer.checkIteratesImpureIterable(Node, AmbiguatedFunctionSummary)
func (b *PureFunctionIdentifier_FunctionBodyAnalyzer) CheckIteratesImpureIterable(node rhino.Node, encloserSummary *PureFunctionIdentifier_AmbiguatedFunctionSummary) {
	if !NodeUtil_IteratesImpureIterable(node) {
		return
	}
	encloserSummary.SetMutatesGlobalStateAndAllOtherFlags(node)
}

// DeprecatedSetSideEffectsForControlLoss assigns the set of side-effects
// associated with an arbitrary loss of control flow to encloserSummary.
//
// This function is kept to retain behaviour but marks places where the
// analysis is inaccurate.
//
// See b/135475880.
//
// Java: PureFunctionIdentifier.FunctionBodyAnalyzer.deprecatedSetSideEffectsForControlLoss(AmbiguatedFunctionSummary)
func (b *PureFunctionIdentifier_FunctionBodyAnalyzer) DeprecatedSetSideEffectsForControlLoss(encloserSummary *PureFunctionIdentifier_AmbiguatedFunctionSummary) {
	b.RecordThrowsBasedOnContext(encloserSummary)
}

// Java: PureFunctionIdentifier.FunctionBodyAnalyzer.recordThrowsBasedOnContext(AmbiguatedFunctionSummary)
func (b *PureFunctionIdentifier_FunctionBodyAnalyzer) RecordThrowsBasedOnContext(encloserSummary *PureFunctionIdentifier_AmbiguatedFunctionSummary) {
	if b.functionScopeStackGetLast().catchDepth == 0 {
		encloserSummary.SetThrows(b.functionScopeStackGetLast().root)
	}
}

// Java: PureFunctionIdentifier.FunctionBodyAnalyzer.enterScope(NodeTraversal)
func (b *PureFunctionIdentifier_FunctionBodyAnalyzer) EnterScope(t *NodeTraversal) {
	function := t.GetScopeRoot()
	if !function.IsFunction() {
		if b.functionScopeStackGetLast().root.IsNil() {
			// Variables declared outside all functions can never share a container
			// scope with any scope inside a function. Pre-populating an empty
			// Scope avoids scanning the global/module AST in SyntacticScopeCreator.
			if len(t.scopes) == 1 {
				t.scopes[0].untypedScope = Scope_CreateGlobalScope(function)
			} else if parentScope := t.scopes[len(t.scopes)-2].untypedScope; parentScope != nil {
				t.scopes[len(t.scopes)-1].untypedScope = Scope_CreateChildScope(parentScope, function)
			}
		}
		return
	}
	jsdoc := NodeUtil_GetBestJSDocInfo(function)
	isArtificiallyPure := jsdoc != nil && jsdoc.IsNoSideEffects()

	n := len(b.functionScopeStack)
	if n < cap(b.functionScopeStack) {
		b.functionScopeStack = b.functionScopeStack[:n+1]
		if e := b.functionScopeStack[n]; e != nil && !e.heldByEdge {
			e.root = function
			e.isArtificiallyPure = isArtificiallyPure
			e.catchDepth = 0
			e.skiplistedVars = e.skiplistedVars[:0]
			e.skiplistedNameNodes = e.skiplistedNameNodes[:0]
			e.taintedVars = e.taintedVars[:0]
		} else {
			b.functionScopeStack[n] = NewPureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry(function, isArtificiallyPure)
		}
	} else {
		b.functionScopeStack = append(b.functionScopeStack, NewPureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry(function, isArtificiallyPure))
	}
	if len(b.outer.summariesForAllNamesOfFunctionByNode[function]) == 0 {
		// This function was not part of a definition which is why it was not created by
		// {@link populateDatastructuresForAnalysisTraversal}. For example, an anonymous
		// function.
		summary := b.outer.createSummaryInGraph("<anonymous>")
		b.outer.putSummaryForFunctionNode(function, summary)
	}
}

// Java: PureFunctionIdentifier.FunctionBodyAnalyzer.exitScope(NodeTraversal)
func (b *PureFunctionIdentifier_FunctionBodyAnalyzer) ExitScope(t *NodeTraversal) {
	idx := len(b.functionScopeStack) - 1
	functionEntry := b.functionScopeStack[idx]
	if functionEntry.root.IsNil() {
		return
	}
	if t.GetScopeRoot() == functionEntry.root {
		// ArrayDeque.removeLast()
		b.functionScopeStack = b.functionScopeStack[:idx]
	}

	if len(functionEntry.taintedVars) == 0 {
		return
	}

	// If the current scope was never materialized, none of its variables could
	// have been looked up or added to taintedVars.
	scope := t.scopes[len(t.scopes)-1].untypedScope
	if scope == nil {
		return
	}
	functionScope := scope.GetClosestContainerScope()
	if !functionScope.IsFunctionScope() {
		return
	}
	jutil.CheckStateMsg1(
		functionScope.GetRootNode() == functionEntry.root, "%s", functionScope.GetRootNode())

	// Handle deferred local variable modifications:
	for _, sideEffectInfo := range b.outer.summariesForAllNamesOfFunctionByNode[functionEntry.root] {
		jutil.CheckNotNullMsg1(sideEffectInfo, "%s has no side effect info.", functionEntry.root)

		if sideEffectInfo.MutatesGlobalState() {
			continue
		}

		for _, v := range scope.varList {
			if v == nil || !functionEntry.isTaintedVar(v) {
				continue
			}
			isFromDestructuring := NodeUtil_IsLhsByDestructuring(v.GetNameNode())
			if v.IsParam() &&
				// Ignore destructuring parameters because they don't directly correspond to an
				// argument passed to the function for the purposes of "setMutatesArguments"
				!isFromDestructuring &&
				!functionEntry.isSkiplistedVar(v) {
				sideEffectInfo.SetMutatesArgumentAt(PureFunctionIdentifier_GetParamIndex(v), v.GetNode())
				continue
			}

			// Take care of locals that might have been tainted.
			if !functionEntry.IsLocalVar(v, scope) {
				// If the function has global side-effects
				// don't bother with the local side-effects.
				sideEffectInfo.SetMutatesGlobalStateAndAllOtherFlags(v.GetNode())
				break
			}
		}
	}
}

// Java: PureFunctionIdentifier.FunctionBodyAnalyzer.isVarDeclaredInSameContainerScope(Var, Scope)
func PureFunctionIdentifier_FunctionBodyAnalyzer_IsVarDeclaredInSameContainerScope(v *Var, scope *Scope) bool {
	if v == nil {
		return false
	}
	varScope, _ := v.GetScope().(*Scope)
	return varScope.HasSameContainerScope(scope)
}

// visitLhsNode records information about the side effects caused by assigning
// a value to a given LHS, materializing traversal's Scope only when needed.
//
// If the operation modifies this or taints global state, mark the enclosing
// function as having those side effects.
//
// Java: PureFunctionIdentifier.FunctionBodyAnalyzer.visitLhsNode(AmbiguatedFunctionSummary, Scope, Node, Predicate<Node>)
func (b *PureFunctionIdentifier_FunctionBodyAnalyzer) visitLhsNode(
	encloserSummary *PureFunctionIdentifier_AmbiguatedFunctionSummary,
	traversal *NodeTraversal,
	lhs rhino.Node,
	hasLocalRhs func(rhino.Node) bool) {
	if get, ok := ast.AsNormalOrOptChainGet(lhs); ok {
		// Although OPTCHAIN_GETPROP can not be an LHS of an assign, it can be a child to DELPROP.
		// e.g. `delete obj?.prop` <==> `obj == null ?  true : delete obj.prop;`
		// Hence the enclosing function's side effects must be recorded.
		object := get.Object()
		if _, isThis := object.AsThis(); isThis {
			if ast.MustFunction(b.functionScopeStackGetLast().root).IsArrow() {
				// `this` in an arrow function is the `this` of the enclosing function, which is not
				// the receiver the arrow function is called with.
				encloserSummary.SetMutatesGlobalStateAndAllOtherFlags(lhs)
			} else {
				encloserSummary.SetMutatesThis(lhs)
			}
		} else if objName, isName := object.AsName(); isName {
			scope := traversal.GetScope()
			var_ := scope.GetVar(objName.Ident())
			if PureFunctionIdentifier_FunctionBodyAnalyzer_IsVarDeclaredInSameContainerScope(var_, scope) {
				// Maybe a local object modification. We won't know for sure until
				// we exit the scope and can validate the value of the local.
				b.functionScopeStackGetLast().addTaintedVar(var_)
			} else {
				encloserSummary.SetMutatesGlobalStateAndAllOtherFlags(lhs)
			}
		} else {
			// Don't track multi level locals: local.prop.prop2++;
			encloserSummary.SetMutatesGlobalStateAndAllOtherFlags(lhs)
		}
	} else {
		name, ok := ast.AsName(lhs)
		jutil.CheckStateMsg1(ok, "%s", lhs)
		scope := traversal.GetScope()
		var_ := scope.GetVar(name.Ident())
		if PureFunctionIdentifier_FunctionBodyAnalyzer_IsVarDeclaredInSameContainerScope(var_, scope) {
			if !hasLocalRhs(lhs) {
				// Assigned value is not guaranteed to be a local value,
				// so if we see any property assignments on this variable,
				// they could be tainting a non-local value.
				b.functionScopeStackGetLast().addSkiplistedVar(var_)
			}
		} else {
			encloserSummary.SetMutatesGlobalStateAndAllOtherFlags(lhs)
		}
	}
}

// VisitCall records information about a call site.
//
// Java: PureFunctionIdentifier.FunctionBodyAnalyzer.visitCall(AmbiguatedFunctionSummary, Scope, Node)
func (b *PureFunctionIdentifier_FunctionBodyAnalyzer) VisitCall(callerInfo *PureFunctionIdentifier_AmbiguatedFunctionSummary, scope *Scope, invocation rhino.Node) {
	// Handle special cases (Math, RegExp)
	// TODO: This logic can probably be replaced with @nosideeffects annotations in externs.
	if invocation.IsCall() && !b.outer.astAnalyzer.FunctionCallHasSideEffects(invocation) {
		return
	}

	// Handle known cases now (Object, Date, RegExp, etc)
	if invocation.IsNew() && !b.outer.astAnalyzer.ConstructorCallHasSideEffects(invocation) {
		return
	}

	calleeSummaries := b.outer.GetSummariesForCallee(invocation)
	if len(calleeSummaries) == 0 {
		callerInfo.SetMutatesGlobalStateAndAllOtherFlags(rhino.NilNode)
		return
	}

	enclosingFunction := b.functionScopeStackGetLast()
	propagatesThrows := enclosingFunction.catchDepth == 0
	enclosingFunction.heldByEdge = true
	for _, calleeInfo := range calleeSummaries {
		edge :=
			b.outer.allocPropagationForInvocation(
				invocation,
				scope,
				enclosingFunction,
				b.ArgumentsMatchParameters(invocation),
				propagatesThrows)
		b.outer.reverseCallGraph.ConnectNodes(calleeInfo.graphNode, edge, callerInfo.graphNode)
	}
}

// ArgumentsMatchParameters reports whether the i-th argument of invocation
// is passed as the i-th parameter of the callee, in the absence of spread
// arguments.
//
// Java: PureFunctionIdentifier.FunctionBodyAnalyzer.argumentsMatchParameters(Node)
func (b *PureFunctionIdentifier_FunctionBodyAnalyzer) ArgumentsMatchParameters(invocation rhino.Node) bool {
	return !invocation.IsTaggedTemplateLit() &&
		!PureFunctionIdentifier_IsInvocationViaCallOrApply(invocation) &&
		// For caching calls like `goog.reflect.cache(cacheObj, key, valueFn)`, the analyzed
		// callee is `valueFn(key)`, whose 0th parameter is the 2nd argument (`key`).
		b.outer.compiler.GetCodingConvention().DescribeCachingCall(invocation) == nil
}

// Java: PureFunctionIdentifier.FunctionBodyAnalyzer.addToCatchDepthIfTryBlock(Node, int)
func (b *PureFunctionIdentifier_FunctionBodyAnalyzer) AddToCatchDepthIfTryBlock(n rhino.Node, delta int) {
	parent := n.GetParent()
	if ast.RoleIn(n, parent) != ast.RoleTryBody {
		return
	}
	if _, hasCatch := ast.MustTry(parent).Catch(); !hasCatch {
		return
	}
	b.functionScopeStackGetLast().catchDepth += delta
}

// PureFunctionIdentifier_GetParamIndex returns the index of param in the
// parameter list of its function.
//
// param must not be destructured (which includes rest parameters).
//
// Java: PureFunctionIdentifier.getParamIndex(Var)
func PureFunctionIdentifier_GetParamIndex(param *Var) int {
	var node rhino.Node
	if param.IsDefaultParam() {
		node = param.GetParentNode()
	} else {
		node = param.GetNameNode()
	}
	_, isParamList := ast.AsParamList(node.GetParent())
	jutil.CheckStateMsg1(isParamList, "%s", node)
	return node.GetParent().GetIndexOfChild(node)
}

// Java: PureFunctionIdentifier.isInvocationViaCallOrApply(Node)
func PureFunctionIdentifier_IsInvocationViaCallOrApply(callSite rhino.Node) bool {
	call, ok := ast.AsNormalOrOptChainCall(callSite)
	if !ok {
		return false
	}
	calleeGet, ok := call.Callee().AsNormalOrOptChainGet()
	if !ok {
		return false
	}
	receiver := calleeGet.Object()
	_, isName := receiver.AsName()
	_, isGetProp := receiver.AsNormalOrOptChainGetProp()
	if !isName && !isGetProp {
		return false
	}
	return NodeUtil_IsFunctionObjectCall(callSite) || NodeUtil_IsFunctionObjectApply(callSite)
}

// Java: PureFunctionIdentifier.isCallOrTaggedTemplateLit(Node)
func PureFunctionIdentifier_IsCallOrTaggedTemplateLit(invocation rhino.Node) bool {
	return invocation.IsCall() || invocation.IsOptChainCall() || invocation.IsTaggedTemplateLit()
}

// PureFunctionIdentifier_NameForReference returns the unqualified name
// associated with an R-value.
//
// For NAMEs this is the name. For GETPROPs this is the last segment
// including a leading dot.
//
// Java: PureFunctionIdentifier.nameForReference(Node)
func PureFunctionIdentifier_NameForReference(nameRef rhino.Node) string {
	if name, ok := ast.AsName(nameRef); ok {
		return name.Ident()
	}
	if gp, ok := ast.AsNormalOrOptChainGetProp(nameRef); ok {
		return PureFunctionIdentifier_PROP_NAME_PREFIX + gp.Property()
	}
	panic(jutil.NewIllegalStateExceptionMsg("Unexpected name reference: " + nameRef.String()))
}

// Java: PureFunctionIdentifier.getPropertyKind(String)
func (p *PureFunctionIdentifier) GetPropertyKind(name string) AccessorSummary_PropertyAccessKind {
	if p.assumeGettersArePure {
		return AccessorSummary_PropertyAccessKind_NORMAL
	}
	return p.compiler.GetAccessorSummary().GetKind(name)
}

// PureFunctionIdentifier_SideEffectPropagation stores all the information
// about a connection between functions needed to propagate side effects from
// one instance of AmbiguatedFunctionSummary to another.
type PureFunctionIdentifier_SideEffectPropagation struct {
	// Whether this propagation represents an aliasing of one name by another.
	// In that case, all side effects of the "callee" just need to be copied
	// onto the "caller".
	callerIsAlias bool

	// If you call a function with apply or call, one of the arguments at the
	// call site will be used as 'this' inside the implementation. If this is
	// pass into apply like so: function.apply(this, ...) then 'this' in the
	// caller is tainted.
	calleeThisEqualsCallerThis bool

	// Whether this propagation includes the "throws" bit.
	//
	// In some contexts, such as an invocation inside a "try", the caller is
	// uneffected by the callee throwing.
	propagateThrows bool

	// The token used to invoke the callee by the caller.
	invocation rhino.Node

	// The scope containing the invocation.
	callScope *Scope

	// The function containing the invocation.
	callerEntry *PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry

	// Whether the i-th argument of the invocation is passed as the i-th parameter of the callee,
	// in the absence of spread arguments.
	argumentsMatchParameters bool
}

// Java: PureFunctionIdentifier.SideEffectPropagation(boolean, boolean, boolean, Node, Scope, FunctionStackEntry, boolean)
func NewPureFunctionIdentifier_SideEffectPropagation(
	callerIsAlias bool,
	calleeThisEqualsCallerThis bool,
	propagateThrows bool,
	invocation rhino.Node,
	callScope *Scope,
	callerEntry *PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry,
	argumentsMatchParameters bool) *PureFunctionIdentifier_SideEffectPropagation {
	jutil.CheckArgumentMsg1(invocation.IsNil() || NodeUtil_IsInvocation(invocation), "%s", invocation)

	return &PureFunctionIdentifier_SideEffectPropagation{
		callerIsAlias:              callerIsAlias,
		calleeThisEqualsCallerThis: calleeThisEqualsCallerThis,
		propagateThrows:            propagateThrows,
		invocation:                 invocation,
		callScope:                  callScope,
		callerEntry:                callerEntry,
		argumentsMatchParameters:   argumentsMatchParameters,
	}
}

var pureFunctionIdentifier_forAliasSingleton = &PureFunctionIdentifier_SideEffectPropagation{
	callerIsAlias:   true,
	propagateThrows: true,
}

// Java: PureFunctionIdentifier.SideEffectPropagation.forAlias()
func PureFunctionIdentifier_SideEffectPropagation_ForAlias() *PureFunctionIdentifier_SideEffectPropagation {
	return pureFunctionIdentifier_forAliasSingleton
}

// Java: PureFunctionIdentifier.SideEffectPropagation.forInvocation(Node, Scope, FunctionStackEntry, boolean, boolean)
func PureFunctionIdentifier_SideEffectPropagation_ForInvocation(
	invocation rhino.Node,
	callScope *Scope,
	callerEntry *PureFunctionIdentifier_FunctionBodyAnalyzer_FunctionStackEntry,
	argumentsMatchParameters bool,
	propagateThrows bool) *PureFunctionIdentifier_SideEffectPropagation {
	jutil.CheckArgumentMsg1(NodeUtil_IsInvocation(invocation), "%s", invocation)

	return NewPureFunctionIdentifier_SideEffectPropagation(
		false,
		// `this` in an arrow function is the `this` of the enclosing function, which is not
		// tracked.
		PureFunctionIdentifier_SideEffectPropagation_CalleeAndCallerShareThis(invocation) && !callerEntry.root.IsArrowFunction(),
		propagateThrows,
		invocation,
		callScope,
		callerEntry,
		argumentsMatchParameters)
}

// holdsLocalValue reports whether expr, passed to the callee, holds a value
// that is local to the caller.
//
// Must only be called once the caller has been fully traversed.
//
// Java: PureFunctionIdentifier.SideEffectPropagation.holdsLocalValue(Node)
func (e *PureFunctionIdentifier_SideEffectPropagation) holdsLocalValue(expr rhino.Node) bool {
	if NodeUtil_EvaluatesToLocalValue(expr) {
		return true
	}
	name, isName := ast.AsName(expr)
	return isName && e.callerEntry.IsLocalVar(e.callScope.GetVar(name.Ident()), e.callScope)
}

// mutatedArgumentsAreLocal reports whether all the arguments that may be
// mutated by the callee are local to the caller.
//
// Java: PureFunctionIdentifier.SideEffectPropagation.mutatedArgumentsAreLocal(long)
func (e *PureFunctionIdentifier_SideEffectPropagation) mutatedArgumentsAreLocal(mutatedParamMask int64) bool {
	if !e.argumentsMatchParameters {
		return NodeUtil_AllArgsUnescapedLocal(e.invocation)
	}
	index := 0
	for arg := range ast.MustCallOrNew(e.invocation).Args() {
		if _, isSpread := arg.AsIterSpread(); isSpread {
			return false
		}
		if (index >= 64 /* Long.SIZE */ || (mutatedParamMask&(int64(1)<<index)) != 0) &&
			!e.holdsLocalValue(arg.Node()) {
			return false
		}
		index++
	}
	return true
}

// receiverIsLocal reports whether the receiver of the invocation is local to
// the caller, e.g. `o` in `o.m()`.
//
// Java: PureFunctionIdentifier.SideEffectPropagation.receiverIsLocal()
func (e *PureFunctionIdentifier_SideEffectPropagation) receiverIsLocal() bool {
	callee := pureFunctionIdentifier_invocationCallee(e.invocation)
	getProp, isGetProp := callee.AsNormalOrOptChainGetProp()
	return !PureFunctionIdentifier_IsInvocationViaCallOrApply(e.invocation) &&
		isGetProp &&
		e.holdsLocalValue(getProp.Object().Node())
}

// Java: PureFunctionIdentifier.SideEffectPropagation.calleeAndCallerShareThis(Node)
func PureFunctionIdentifier_SideEffectPropagation_CalleeAndCallerShareThis(invocation rhino.Node) bool {
	if !PureFunctionIdentifier_IsCallOrTaggedTemplateLit(invocation) {
		return false // Calling a constructor creates a new object bound to `this`.
	}

	callee := pureFunctionIdentifier_invocationCallee(invocation)
	if _, isSuper := callee.AsSuper(); isSuper {
		return true
	}

	var thisArg rhino.Node
	if PureFunctionIdentifier_IsInvocationViaCallOrApply(invocation) {
		// If the call site is actually a `.call` or `.apply`, then `this` will be an argument.
		if firstArg, ok := ast.MustCallOrNew(invocation).FirstArg(); ok {
			thisArg = firstArg.Node()
		}
	} else if getProp, ok := callee.AsNormalOrOptChainGetProp(); ok {
		thisArg = getProp.Object().Node()
	}

	if thisArg.IsNil() {
		return false // No `this` is being passed.
	} else if thisArg.IsThis() || thisArg.IsSuper() {
		return true
	}

	// TODO(nickreid): If `thisArg` is a known local or known arg we could say something more
	// specific about the effect of the callee mutating `this`.

	// We're not sure what `this` is being passed, so make a conservative choice.
	return false
}

// Propagate propagates the side effects from the callee to the caller.
//
// callee: propagate from; caller: propagate to. Returns true if the
// propagation changed the side effects on the caller.
//
// Java: PureFunctionIdentifier.SideEffectPropagation.propagate(AmbiguatedFunctionSummary, AmbiguatedFunctionSummary)
func (e *PureFunctionIdentifier_SideEffectPropagation) Propagate(callee *PureFunctionIdentifier_AmbiguatedFunctionSummary, caller *PureFunctionIdentifier_AmbiguatedFunctionSummary) bool {
	initialCallerFlags := caller.bitmask
	initialCallerParamMask := caller.mutatedParamMask

	if e.callerIsAlias {
		caller.SetMask(callee.bitmask)
		caller.mutatedParamMask |= callee.mutatedParamMask
	} else if callee.MutatesGlobalState() {
		// If the callee modifies global state then so does that caller.
		caller.SetMutatesGlobalStateAndAllOtherFlags(e.invocation)
	} else {
		if e.propagateThrows && callee.FunctionThrows() {
			// If the callee throws an exception then so does the caller.
			caller.SetThrows(e.invocation)
		}
		if callee.MutatesArguments() && !e.mutatedArgumentsAreLocal(callee.mutatedParamMask) {
			// If the callee mutates its input arguments and the arguments escape the caller then it
			// has unbounded side effects.
			caller.SetMutatesGlobalStateAndAllOtherFlags(e.invocation)
		}
		if callee.MutatesThis() {
			if e.invocation.IsNew() {
				// NEWing a constructor provide a unescaped "this" making side-effects impossible.
			} else if e.calleeThisEqualsCallerThis {
				caller.SetMutatesThis(e.invocation)
			} else if e.receiverIsLocal() {
				// The callee mutates an object local to the caller, e.g. `const o = new Foo(); o.m();`.
			} else {
				caller.SetMutatesGlobalStateAndAllOtherFlags(e.invocation)
			}
		}
	}

	return caller.bitmask != initialCallerFlags ||
		caller.mutatedParamMask != initialCallerParamMask
}
