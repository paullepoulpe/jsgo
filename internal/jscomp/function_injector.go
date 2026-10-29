// Port of com/google/javascript/jscomp/FunctionInjector.java @ closure-compiler v20261021

/*
 * Copyright 2008 The Closure Compiler Authors.
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

// In order to estimate the cost of lining, we make the assumption that
// Identifiers are reduced 2 characters. For the call arguments, the
// important thing is that the cost is assumed to be the same in the call and
// the function, so the actual length doesn't matter in most cases.
const FunctionInjector_NAME_COST_ESTIMATE = InlineCostEstimator_ESTIMATED_IDENTIFIER_COST

// The cost of a argument separator (a comma).
const FunctionInjector_COMMA_COST = 1

// The cost of the parentheses needed to make a call.
const FunctionInjector_PAREN_COST = 2

func functionInjector_blocksMinimumRequirements(n rhino.Node, fnNameID, fnRecursionNameID uint32, checkArguments bool) bool {
	switch n.GetToken() {
	case rhino.Token_NAME:
		nameID := n.GetStringID()
		if (checkArguments && nameID == calleeArgumentsIdentID) ||
			nameID == calleeEvalIdentID ||
			(fnNameID != 0 && nameID == fnNameID) ||
			(fnRecursionNameID != 0 && nameID == fnRecursionNameID) {
			return true
		}
	case rhino.Token_SUPER, rhino.Token_NEW_TARGET:
		return true
	case rhino.Token_FUNCTION:
		if !n.IsArrowFunction() {
			checkArguments = false
		}
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if c.GetFirstChild().IsNil() {
			switch c.GetToken() {
			case rhino.Token_NAME:
				nameID := c.GetStringID()
				if (checkArguments && nameID == calleeArgumentsIdentID) ||
					nameID == calleeEvalIdentID ||
					(fnNameID != 0 && nameID == fnNameID) ||
					(fnRecursionNameID != 0 && nameID == fnRecursionNameID) {
					return true
				}
			case rhino.Token_SUPER, rhino.Token_NEW_TARGET:
				return true
			}
			continue
		}
		if functionInjector_blocksMinimumRequirements(c, fnNameID, fnRecursionNameID, checkArguments) {
			return true
		}
	}
	return false
}

// DoesFunctionMeetMinimumRequirements: fnName is the name of this function.
// This either the name of the variable to which the function is assigned or
// the name from the FUNCTION node. fnNode is the FUNCTION node of the
// function to inspect. Returns whether the function node meets the minimum
// requirements for inlining.
//
// Java: FunctionInjector.doesFunctionMeetMinimumRequirements(String, Node)
func (f *FunctionInjector) DoesFunctionMeetMinimumRequirements(fnName string, fnNode rhino.Node) bool {
	fn := ast.MustFunction(fnNode)
	block := fn.Body().Node()

	// Basic restrictions on functions that can be inlined:
	// 1) It contains a reference to itself.
	// 2) It uses its parameters indirectly using "arguments" (it isn't
	//    handled yet.
	// 3) It references "eval". Inline a function containing eval can have
	//    large performance implications.

	var fnNameID uint32
	if fnName != "" {
		fnNameID = rhino.RhinoStringPool_LookupID(fnName)
	}
	// The name of an anonymous function is "": no recursion through it.
	var fnRecursionNameID uint32
	if name, ok := fn.Name(); ok {
		if recID := name.Node().GetStringID(); recID != fnNameID {
			fnRecursionNameID = recID
		}
	}

	return !functionInjector_blocksMinimumRequirements(block, fnNameID, fnRecursionNameID, true)
}

func linkedHashSetSize(s *jutil.LinkedHashSet[string]) int {
	if s == nil {
		return 0
	}
	return s.Size()
}

// CanInlineReferenceToFunction: fnNode is the function to evaluate for
// inlining; needAliases a set of function parameter names that can not be
// used without aliasing (returned by getUnsafeParameterNames());
// referencesThis whether fnNode contains references to its this object;
// containsClosures whether fnNode contains inner functions or non-static
// class field initializers. Returns whether the inlining can occur.
//
// Java: FunctionInjector.canInlineReferenceToFunction(Reference, Node, ImmutableSet<String>, boolean, boolean)
func (f *FunctionInjector) CanInlineReferenceToFunction(
	ref *FunctionInjector_Reference,
	fnNode rhino.Node,
	needAliases *jutil.LinkedHashSet[string],
	referencesThis bool,
	containsClosures bool) FunctionInjector_CanInlineResult {
	return f.canInlineReferenceWithSummary(
		ref, fnNode, linkedHashSetSize(needAliases), referencesThis, containsClosures, nil)
}

func (f *FunctionInjector) canInlineReferenceWithSummary(
	ref *FunctionInjector_Reference,
	fnNode rhino.Node,
	aliasCount int,
	referencesThis bool,
	containsClosures bool,
	summary *calleeSummary) FunctionInjector_CanInlineResult {
	// TODO(johnlenz): This function takes too many parameter, without
	// context.  Modify the API to take a structure describing the function.

	// Allow direct function calls or "fn.call" style calls.
	callNode := ref.callNode
	if !f.IsSupportedCallType(callNode, fnNode) {
		return FunctionInjector_CanInlineResult_NO
	}

	if FunctionInjector_HasSpreadCallArgument(callNode) {
		return FunctionInjector_CanInlineResult_NO
	}

	// Limit where functions that contain functions can be inline.  Introducing
	// an inner function into another function can capture a variable and cause
	// a memory leak.  This isn't a problem in the global scope as those values
	// last until explicitly cleared.
	if containsClosures {
		if !ref.inGlobalScope() {
			// TODO(johnlenz): Allow inlining into any scope without local names or inner functions.
			return FunctionInjector_CanInlineResult_NO
		} else if NodeUtil_IsWithinLoop(callNode) {
			// An inner closure maybe relying on a local value holding a value for a
			// single iteration through a loop.
			return FunctionInjector_CanInlineResult_NO
		}
	}

	// TODO(johnlenz): Add support for 'apply'
	if referencesThis && !NodeUtil_IsFunctionObjectCall(callNode) {
		// TODO(johnlenz): Allow 'this' references to be replaced with a
		// global 'this' object.
		return FunctionInjector_CanInlineResult_NO
	}

	if ref.mode == FunctionInjector_InliningMode_DIRECT {
		return f.canInlineReferenceDirectlyWithSummary(ref, fnNode, aliasCount, summary)
	} else {
		return f.canInlineReferenceAsStatementBlockWithSummary(ref, fnNode, aliasCount, summary)
	}
}

// IsSupportedCallType: only ".call" calls and direct calls to functions are
// supported.
//
// callNode is the call evaluate. Returns whether the call is of a type that
// is supported.
//
// Java: FunctionInjector.isSupportedCallType(Node, Node)
func (f *FunctionInjector) IsSupportedCallType(callNode rhino.Node, fnNode rhino.Node) bool {
	call := ast.MustNormalOrOptChainCall(callNode)
	if call.Callee().Token() != rhino.Token_NAME {
		if NodeUtil_IsFunctionObjectCall(callNode) {
			if !f.assumeStrictThis && !fnNode.IsArrowFunction() {
				thisValue, ok := call.FirstArg()
				if !ok || thisValue.Token() != rhino.Token_THIS {
					return false
				}
			}
		} else if NodeUtil_IsFunctionObjectApply(callNode) {
			return false
		}
	}

	return true
}

// Java: FunctionInjector.hasSpreadCallArgument(Node)
func FunctionInjector_HasSpreadCallArgument(callNode rhino.Node) bool {
	call, ok := ast.AsNormalOrOptChainCall(callNode)
	jutil.CheckArgumentMsg1(ok, "%s", callNode)
	for arg := range call.Args() {
		if arg.Token() == rhino.Token_ITER_SPREAD {
			return true
		}
	}
	return false
}

// FunctionInjector_IsDirectCallNodeReplacementPossible checks if the given
// function matches the criteria for an inlinable function, and if so, adds
// it to our set of inlinable functions.
//
// Java: FunctionInjector.isDirectCallNodeReplacementPossible(Node)
func FunctionInjector_IsDirectCallNodeReplacementPossible(fnNode rhino.Node) bool {
	// Only inline single-statement functions
	body := ast.MustFunction(fnNode).Body()
	if block, ok := body.AsBlock(); ok {
		// Check if this function is suitable for direct replacement of a CALL node:
		// a function that consists of single return that returns an expression.
		if !block.HasStatements() {
			// special case empty functions.
			return true
		}
		if onlyStmt, ok := block.OnlyStatement(); ok {
			// Only inline functions that return something.
			if ret, ok := onlyStmt.AsReturn(); ok {
				_, hasValue := ret.Value()
				return hasValue
			}
		}
		return false
	}
	return !body.Node().HasChildren()
}

// FunctionInjector_CanInlineResult is the enum FunctionInjector.CanInlineResult.
type FunctionInjector_CanInlineResult int32

const (
	FunctionInjector_CanInlineResult_YES FunctionInjector_CanInlineResult = iota
	FunctionInjector_CanInlineResult_AFTER_PREPARATION
	FunctionInjector_CanInlineResult_NO
)

var functionInjector_CanInlineResult_names = [...]string{"YES", "AFTER_PREPARATION", "NO"}

// Java: FunctionInjector.CanInlineResult.name()
func (r FunctionInjector_CanInlineResult) String() string {
	return functionInjector_CanInlineResult_names[r]
}

// Java: FunctionInjector.CanInlineResult.ordinal()
func (r FunctionInjector_CanInlineResult) Ordinal() int { return int(r) }

// CanInlineReferenceAsStatementBlock determines whether a function can be
// inlined at a particular call site. There are several criteria that the
// function and reference must hold in order for the functions to be
// inlined: - It must be a simple call, or assignment, or var
// initialization.
//
//	f();
//	a = foo();
//	var a = foo();
//
// Java: FunctionInjector.canInlineReferenceAsStatementBlock(Reference, Node, ImmutableSet<String>)
func (f *FunctionInjector) CanInlineReferenceAsStatementBlock(
	ref *FunctionInjector_Reference, fnNode rhino.Node, namesToAlias *jutil.LinkedHashSet[string]) FunctionInjector_CanInlineResult {
	return f.canInlineReferenceAsStatementBlockWithSummary(ref, fnNode, linkedHashSetSize(namesToAlias), nil)
}

func (f *FunctionInjector) canInlineReferenceAsStatementBlockWithSummary(
	ref *FunctionInjector_Reference,
	fnNode rhino.Node,
	aliasCount int,
	summary *calleeSummary) FunctionInjector_CanInlineResult {
	callSiteType := f.classifyCallSite(ref, f.allowDecomposition)
	if callSiteType == FunctionInjector_CallSiteType_UNSUPPORTED {
		return FunctionInjector_CanInlineResult_NO
	}

	if !f.allowDecomposition &&
		(callSiteType == FunctionInjector_CallSiteType_DECOMPOSABLE_EXPRESSION ||
			callSiteType == FunctionInjector_CallSiteType_EXPRESSION) {
		return FunctionInjector_CanInlineResult_NO
	}

	if !f.callMeetsBlockInliningRequirementsWithSummary(ref, fnNode, aliasCount, summary) {
		return FunctionInjector_CanInlineResult_NO
	}

	if callSiteType == FunctionInjector_CallSiteType_DECOMPOSABLE_EXPRESSION ||
		callSiteType == FunctionInjector_CallSiteType_EXPRESSION {
		return FunctionInjector_CanInlineResult_AFTER_PREPARATION
	} else {
		return FunctionInjector_CanInlineResult_YES
	}
}

func functionInjector_isEvalName(n rhino.Node) bool {
	return n.IsName() && n.GetStringID() == calleeEvalIdentID
}

func functionInjector_isShallowStatement(n rhino.Node) bool {
	return n.IsRoot() || n.IsBlock() || (!n.IsFunction() && NodeUtil_IsShallowStatementTree(n.GetParent()))
}

func functionInjector_hasShallowDeclaration(n rhino.Node) bool {
	if NodeUtil_IsDeclaration(n) {
		return true
	}
	if functionInjector_isShallowStatement(n) {
		for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			if functionInjector_hasShallowDeclaration(c) {
				return true
			}
		}
	}
	return false
}

func functionInjector_countShallowReturns(n rhino.Node) int {
	count := 0
	if n.IsReturn() {
		count++
	}
	if functionInjector_isShallowStatement(n) {
		for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			count += functionInjector_countShallowReturns(c)
		}
	}
	return count
}

func functionInjector_referencesEval(n rhino.Node, rootFn rhino.Node) bool {
	if n.IsName() && n.GetStringID() == calleeEvalIdentID {
		return true
	}
	if n.IsFunction() && n != rootFn {
		return false
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if c.GetFirstChild().IsNil() {
			if c.IsName() && c.GetStringID() == calleeEvalIdentID {
				return true
			}
			continue
		}
		if functionInjector_referencesEval(c, rootFn) {
			return true
		}
	}
	return false
}

// ReferencesEval returns whether or not fn includes a call to `eval` in its
// scope.
//
// Results are cached to make subsequent calls faster.
//
// Java: FunctionInjector.referencesEval(Node)
func (f *FunctionInjector) ReferencesEval(fn rhino.Node) bool {
	jutil.CheckState(fn.IsFunction())

	info := f.fnInfoCache[fn]
	if info.evalState != 0 {
		return info.evalState == 2
	}

	result := functionInjector_referencesEval(fn, fn)
	if f.fnInfoCache == nil {
		f.fnInfoCache = make(map[rhino.Node]functionInjectorFnInfo, 64)
	}
	if result {
		info.evalState = 2
	} else {
		info.evalState = 1
	}
	f.fnInfoCache[fn] = info
	return result
}

func functionInjector_findInnerFunctions(n rhino.Node, containerFn rhino.Node, first rhino.Node, count int) (rhino.Node, int) {
	if n != containerFn && NodeUtil_IsFunctionOrNonStaticClassFieldInitializer(n) {
		if count == 0 {
			first = n
		}
		count++
		if count > 1 {
			return first, count
		}
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if c.GetFirstChild().IsNil() {
			if c.IsMemberFieldDef() && !c.IsStaticMember() {
				if count == 0 {
					first = c
				}
				count++
				if count > 1 {
					return first, count
				}
			}
			continue
		}
		first, count = functionInjector_findInnerFunctions(c, containerFn, first, count)
		if count > 1 {
			return first, count
		}
	}
	return first, count
}

// InnerFunctionOf returns any inner function or non-static class field
// initializer of containerFn.
//
// Non-static class field initializers are treated like inner functions
// because they are evaluated later, once per instantiation, and may capture
// locals.
//
// If there are no inner functions, or multilple inner functions, the
// sentinel values NO_FUNCTIONS, MULTIPLE_FUNCTIONS are returned
// respectively.
//
// Java: FunctionInjector.innerFunctionOf(Node)
func (f *FunctionInjector) InnerFunctionOf(containerFn rhino.Node) rhino.Node {
	jutil.CheckState(containerFn.IsFunction())

	info := f.fnInfoCache[containerFn]
	if !info.inner.IsNil() {
		return info.inner
	}

	firstInner, count := functionInjector_findInnerFunctions(containerFn, containerFn, rhino.NilNode, 0)

	var cached rhino.Node
	switch count {
	case 0:
		cached = FunctionInjector_NO_FUNCTIONS
	case 1:
		cached = firstInner
	default:
		cached = FunctionInjector_MULTIPLE_FUNCTIONS
	}

	if f.fnInfoCache == nil {
		f.fnInfoCache = make(map[rhino.Node]functionInjectorFnInfo, 64)
	}
	info.inner = cached
	f.fnInfoCache[containerFn] = info
	return cached
}

// CallMeetsBlockInliningRequirements determines whether a function can be
// inlined at a particular call site. - Don't inline if the calling function
// contains an inner function and inlining would introduce new globals.
//
// Java: FunctionInjector.callMeetsBlockInliningRequirements(Reference, Node, ImmutableSet<String>)
func (f *FunctionInjector) CallMeetsBlockInliningRequirements(
	callRef *FunctionInjector_Reference, calleeFn rhino.Node, namesToAlias *jutil.LinkedHashSet[string]) bool {
	return f.callMeetsBlockInliningRequirementsWithSummary(callRef, calleeFn, linkedHashSetSize(namesToAlias), nil)
}

func (f *FunctionInjector) callMeetsBlockInliningRequirementsWithSummary(
	callRef *FunctionInjector_Reference,
	calleeFn rhino.Node,
	aliasCount int,
	summary *calleeSummary) bool {
	// Note: functions that contain function definitions are filtered out
	// in isCandidateFunction.

	// TODO(johnlenz): Determining if the called function contains VARs
	// or if the caller contains inner functions accounts for 20% of the
	// run-time cost of this pass.
	//
	// Note that as of 2019-10-11, "eval" and function checking are cached so this may be less
	// of an issue.

	forbidTemps := false
	if hoistRoot := callRef.closestHoistScopeRoot(); !hoistRoot.IsNil() {
		callerFn := hoistRoot.GetParent() // the function of the hoist scope's block

		// Don't allow any new vars into a scope that contains eval or one
		// that contains functions (excluding the function being inlined).
		if f.ReferencesEval(callerFn) {
			forbidTemps = true
		} else {
			innerFn := f.InnerFunctionOf(callerFn)
			calleeIsOnlyInnerFn := innerFn == FunctionInjector_NO_FUNCTIONS || innerFn == calleeFn
			forbidTemps = !calleeIsOnlyInnerFn
		}
	}

	// If the caller contains functions or evals, verify we aren't adding any
	// additional VAR declarations because aliasing is needed.
	if forbidTemps {
		if aliasCount > 0 {
			return false
		}

		// Don't inline functions with var declarations into a scope with inner
		// functions as the new vars would leak into the inner function and
		// cause memory leaks.
		var calleeContainsVars bool
		if summary != nil && summary.fnNode == calleeFn {
			calleeContainsVars = summary.hasShallowDeclarations
		} else {
			calleeContainsVars = functionInjector_hasShallowDeclaration(ast.MustFunction(calleeFn).Body().Node())
		}
		if calleeContainsVars {
			return false
		}

		bindings := f.functionArgumentInjector.getCallParamBindingsScratch(
			calleeFn, callRef.callNode, f.safeNameIdSupplier)
		if len(bindings) > 0 {
			if f.functionArgumentInjector.hasCallArgumentsNeedingTempsForBindings(
				f.compiler, calleeFn, bindings, nil, f.compiler.GetCodingConvention(), summary) {
				return false
			}
		}
	}

	return true
}

// CanInlineReferenceDirectly determines whether a function can be inlined at
// a particular call site. There are several criteria that the function and
// reference must hold in order for the functions to be inlined: 1) If a
// call's arguments have side effects, the corresponding argument in the
// function must only be referenced once. For instance, this will not be
// inlined:
//
//	function foo(a) { return a + a }
//	x = foo(i++);
//
// Java: FunctionInjector.canInlineReferenceDirectly(Reference, Node, ImmutableSet<String>)
func (f *FunctionInjector) CanInlineReferenceDirectly(
	ref *FunctionInjector_Reference, fnNode rhino.Node, namesToAlias *jutil.LinkedHashSet[string]) FunctionInjector_CanInlineResult {
	return f.canInlineReferenceDirectlyWithSummary(ref, fnNode, linkedHashSetSize(namesToAlias), nil)
}

func (f *FunctionInjector) canInlineReferenceDirectlyWithSummary(
	ref *FunctionInjector_Reference,
	fnNode rhino.Node,
	aliasCount int,
	summary *calleeSummary) FunctionInjector_CanInlineResult {
	if !FunctionInjector_IsDirectCallNodeReplacementPossible(fnNode) {
		return FunctionInjector_CanInlineResult_NO
	}

	// CALL NODE: [ NAME, ARG1, ARG2, ... ]
	callNode := ref.callNode
	call := ast.MustNormalOrOptChainCall(callNode)

	// Functions called via 'call' and 'apply' have a this-object as
	// the first parameter, but this is not part of the called function's
	// parameter list.
	if call.Callee().Token() != rhino.Token_NAME {
		if NodeUtil_IsFunctionObjectCall(callNode) {
			// TODO(johnlenz): Support replace this with a value.
			thisArg, ok := call.FirstArg()
			if !fnNode.IsArrowFunction() && (!ok || thisArg.Token() != rhino.Token_THIS) {
				return FunctionInjector_CanInlineResult_NO
			}
		} else {
			// ".apply" call should be filtered before this.
			jutil.CheckState(!NodeUtil_IsFunctionObjectApply(callNode))
		}
	}

	if aliasCount > 0 {
		return FunctionInjector_CanInlineResult_NO
	}

	bindings := f.functionArgumentInjector.getCallParamBindingsScratch(
		fnNode, callNode, f.throwawayNameSupplier)
	if len(bindings) > 0 {
		if f.functionArgumentInjector.hasCallArgumentsNeedingTempsForBindings(
			f.compiler, fnNode, bindings, nil, f.compiler.GetCodingConvention(), summary) {
			return FunctionInjector_CanInlineResult_NO
		}
	}

	return FunctionInjector_CanInlineResult_YES
}

// InliningLowersCost determines if inlining the function is likely to reduce
// the code size.
//
// Java: FunctionInjector.inliningLowersCost(JSChunk, Node, Collection<? extends Reference>, Set<String>, boolean, boolean)
func (f *FunctionInjector) InliningLowersCost(
	fnChunk *JSChunk,
	fnNode rhino.Node,
	refs []*FunctionInjector_Reference,
	namesToAlias *jutil.LinkedHashSet[string],
	isRemovable bool,
	referencesThis bool) bool {
	return f.inliningLowersCostWithSummary(
		fnChunk, fnNode, refs, linkedHashSetSize(namesToAlias), isRemovable, referencesThis, nil)
}

func (f *FunctionInjector) inliningLowersCostWithSummary(
	fnChunk *JSChunk,
	fnNode rhino.Node,
	refs []*FunctionInjector_Reference,
	aliasCount int,
	isRemovable bool,
	referencesThis bool,
	summary *calleeSummary) bool {
	referenceCount := len(refs)
	if referenceCount == 0 {
		return true
	}

	referencesUsingBlockInlining := 0

	checkChunks := isRemovable && fnChunk != nil
	chunkGraph := f.compiler.GetChunkGraph()

	for _, ref := range refs {
		if ref.mode == FunctionInjector_InliningMode_BLOCK {
			referencesUsingBlockInlining++
		}

		// Check if any of the references cross the chunk boundaries.
		if checkChunks && ref.chunk != nil {
			if ref.chunk != fnChunk && !chunkGraph.DependsOn(ref.chunk, fnChunk) {
				// Calculate the cost as if the function were non-removable,
				// if it still lowers the cost inline it.
				isRemovable = false
				checkChunks = false // no need to check additional chunks.
			}
		}
	}

	referencesUsingDirectInlining := referenceCount - referencesUsingBlockInlining

	// Don't bother calculating the cost of function for simple functions where
	// possible.
	// However, when inlining a complex function, even a single reference may be
	// larger than the original function if there are many returns (resulting
	// in additional assignments) or many parameters that need to be aliased
	// so use the cost estimating.
	if referenceCount == 1 && isRemovable && referencesUsingDirectInlining == 1 {
		return true
	}

	callCost := FunctionInjector_EstimateCallCost(fnNode, referencesThis)
	overallCallCost := callCost * referenceCount

	costDeltaDirect := functionInjector_inlineCostDeltaWithSummary(fnNode, aliasCount, FunctionInjector_InliningMode_DIRECT, summary)
	costDeltaBlock := functionInjector_inlineCostDeltaWithSummary(fnNode, aliasCount, FunctionInjector_InliningMode_BLOCK, summary)

	return FunctionInjector_DoesLowerCost(
		fnNode,
		overallCallCost,
		referencesUsingDirectInlining,
		costDeltaDirect,
		referencesUsingBlockInlining,
		costDeltaBlock,
		isRemovable)
}

// FunctionInjector_DoesLowerCost returns whether inlining will lower cost.
//
// Java: FunctionInjector.doesLowerCost(Node, int, int, int, int, int, boolean)
func FunctionInjector_DoesLowerCost(
	fnNode rhino.Node,
	callCost int,
	directInlines int,
	costDeltaDirect int,
	blockInlines int,
	costDeltaBlock int,
	removable bool) bool {

	// Determine the threshold value for this inequality:
	//     inline_cost < call_cost
	// But solve it for the function declaration size so the size of it
	// is only calculated once and terminated early if possible.

	removableCount := 0
	if removable {
		removableCount = 1
	}
	fnInstanceCount := directInlines + blockInlines - removableCount
	// Prevent division by zero.
	if fnInstanceCount == 0 {
		// Special case single reference function that are being block inlined:
		// If the cost of the inline is greater than the function definition size,
		// don't inline.
		return blockInlines <= 0 || costDeltaBlock <= 0
	}

	costDelta := (directInlines * -costDeltaDirect) + (blockInlines * -costDeltaBlock)
	threshold := (callCost + costDelta) / fnInstanceCount

	return InlineCostEstimator_GetCostWithThreshold(fnNode, threshold+1) <= threshold
}

// FunctionInjector_EstimateCallCost gets an estimate of the cost in
// characters of making the function call: the sum of the identifiers and the
// separators.
//
// Java: FunctionInjector.estimateCallCost(Node, boolean)
func FunctionInjector_EstimateCallCost(fnNode rhino.Node, referencesThis bool) int {
	numArgs := ast.MustFunction(fnNode).Params().NumParams()

	callCost := FunctionInjector_NAME_COST_ESTIMATE + FunctionInjector_PAREN_COST
	if numArgs > 0 {
		callCost += (numArgs * FunctionInjector_NAME_COST_ESTIMATE) + ((numArgs - 1) * FunctionInjector_COMMA_COST)
	}

	if referencesThis {
		// TODO(johnlenz): Update this if we start supporting inlining
		// other functions that reference this.
		// The only functions that reference this that are currently inlined
		// are those that are called via ".call" with an explicit "this".
		callCost += 5 + 5 // ".call" + "this,"
	}

	return callCost
}

// FunctionInjector_InlineCostDelta returns the difference between the
// function definition cost and inline cost.
//
// Java: FunctionInjector.inlineCostDelta(Node, Set<String>, InliningMode)
func FunctionInjector_InlineCostDelta(fnNode rhino.Node, namesToAlias *jutil.LinkedHashSet[string], mode FunctionInjector_InliningMode) int {
	return functionInjector_inlineCostDeltaWithSummary(fnNode, linkedHashSetSize(namesToAlias), mode, nil)
}

func functionInjector_inlineCostDeltaWithSummary(
	fnNode rhino.Node,
	aliasCount int,
	mode FunctionInjector_InliningMode,
	summary *calleeSummary) int {
	// The part of the function that is never inlined:
	//    "function xx(xx,xx){}" (15 + (param count * 3) -1;
	fn := ast.MustFunction(fnNode)
	paramCount := fn.Params().NumParams()
	commaCount := 0
	if paramCount > 1 {
		commaCount = paramCount - 1
	}
	costDeltaFunctionOverhead :=
		15 + commaCount + (paramCount * InlineCostEstimator_ESTIMATED_IDENTIFIER_COST)

	block := fn.Body().Node()
	if !block.HasChildren() {
		// Assume the inline cost is zero for empty functions.
		return -costDeltaFunctionOverhead
	}

	if mode == FunctionInjector_InliningMode_DIRECT {
		// The part of the function that is inlined using direct inlining:
		//    "return " (7)
		return -(costDeltaFunctionOverhead + 7)
	} else {
		// Originally, we estimated purely base on the function code size, relying
		// on later optimizations. But that did not produce good results, so here
		// we try to estimate the something closer to the actual inlined coded.

		// NOTE 1: Result overhead is only if there is an assignment, but
		// getting that information would require some refactoring.
		// NOTE 2: The aliasing overhead is currently an under-estimate,
		// as some parameters are aliased because of the parameters used.
		// Perhaps we should just assume all parameters will be aliased?
		const inlineBlockOverhead = 4     // "X:{}"
		const perReturnOverhead = 2       // "return" --> "break X"
		const perReturnResultOverhead = 3 // "XX="
		const perAliasOverhead = 3        // "XX="

		var returnCount int
		if summary != nil && summary.fnNode == fnNode {
			returnCount = int(summary.shallowReturnCount)
		} else {
			returnCount = functionInjector_countShallowReturns(block)
		}
		resultCount := 0
		if returnCount > 0 {
			resultCount = returnCount - 1
		}
		baseOverhead := 0
		if returnCount > 0 {
			baseOverhead = inlineBlockOverhead
		}

		overhead := baseOverhead +
			returnCount*perReturnOverhead +
			resultCount*perReturnResultOverhead +
			aliasCount*perAliasOverhead

		return overhead - costDeltaFunctionOverhead
	}
}
