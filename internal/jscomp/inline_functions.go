// Port of com/google/javascript/jscomp/InlineFunctions.java @ closure-compiler v20261021

/*
 * Copyright 2005 The Closure Compiler Authors.
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

// Go-only: the shapes this pass reads (a var holding a function, the
// function a call calls, how a name is used, the params and body of a
// function, the members of a class body) go through the typed view
// (internal/ast); the subtree walks and the parallel bookkeeping stay on
// rhino.Node.

// inlineFunctionsRunState is the state of one run of InlineFunctions.Process,
// dropped when the run ends.
type inlineFunctionsRunState struct {
	fns         *jutil.LinkedHashMap[string, *InlineFunctions_FunctionState]
	anonFns     map[rhino.Node]string
	fnNames     map[string]bool
	calledIDs   inlineIDSet
	refsScratch []*FunctionInjector_Reference
	movedArgs   map[rhino.Node]struct{}
	references  jutil.Slab[FunctionInjector_Reference] // one run's references die together
}

// InlineFunctions inlines functions that are divided into two types: "direct
// call node replacement" (aka "direct") and as a block of statements (aka
// block). Function that can be inlined "directly" functions consist of a
// single return statement, everything else is must be inlined as a "block".
// These functions must meet these general requirements: - it is not
// recursive - the function does not contain another function -- these may
// be intentional to to limit the scope of closures. - function is called
// only once OR the size of the inline function is smaller than the call
// itself. - the function name is not referenced in any other manner
//
// "directly" inlined functions must meet these additional requirements: -
// consists of a single return statement
type InlineFunctions struct {
	// TODO(nicksantos): This needs to be completely rewritten to use scopes
	// to do variable lookups. Right now, it assumes that all functions are
	// uniquely named variables. There's currently a stopgap scope-check
	// to ensure that this doesn't produce invalid code. But in the long run,
	// this needs a major refactor.
	fns     *jutil.LinkedHashMap[string, *InlineFunctions_FunctionState]
	anonFns map[rhino.Node]string

	compiler AbstractCompiler

	injector                 *FunctionInjector
	functionArgumentInjector *FunctionArgumentInjector

	reach CompilerOptions_Reach

	enforceMaxSizeAfterInlining bool
	maxSizeAfterInlining        int

	runState *inlineFunctionsRunState

	lastCheckedBlock              rhino.Node
	lastCheckedBlockHasLetOrConst bool
	lastContainingFn              rhino.Node
	scratchSummary                calleeSummary
	lastContainingFnSize          int
}

// Java: InlineFunctions.MISSED_REQUIRED_INLINING
var InlineFunctions_MISSED_REQUIRED_INLINING = DiagnosticType_Warning(
	"JSC_MISSED_REQUIRED_INLINING",
	"function {1} annotated @requireInlining could not be inlined here")

// Java: InlineFunctions(AbstractCompiler, Supplier<String>, Reach, boolean, int)
func NewInlineFunctions(
	compiler AbstractCompiler,
	safeNameIdSupplier func() string,
	reach CompilerOptions_Reach,
	assumeStrictThis bool,
	maxSizeAfterInlining int) *InlineFunctions {
	jutil.CheckArgument(compiler != nil)
	jutil.CheckArgument(safeNameIdSupplier != nil)
	jutil.CheckArgument(reach != CompilerOptions_Reach_NONE)

	f := &InlineFunctions{}
	f.compiler = compiler

	f.reach = reach

	f.maxSizeAfterInlining = maxSizeAfterInlining
	f.enforceMaxSizeAfterInlining =
		maxSizeAfterInlining != CompilerOptions_UNLIMITED_FUN_SIZE_AFTER_INLINING

	// TODO(b/124253050): Update bookkeeping logic and reenable method call inliing.
	// Method call decomposition creates new call nodes after all the analysis
	// is done, which would cause such calls to function to be left behind after
	// the function itself is removed.  The function inliner need to be made
	// aware of these new calls in order to enable it.

	f.functionArgumentInjector = NewFunctionArgumentInjector(compiler.GetAstAnalyzer())
	f.injector = NewFunctionInjector_Builder(compiler).
		SafeNameIdSupplier(safeNameIdSupplier).
		AssumeStrictThis(assumeStrictThis).
		FunctionArgumentInjector(f.functionArgumentInjector).
		Build()
	f.injector.mutator.onMovedArg = func(arg rhino.Node) {
		if f.runState != nil {
			if f.runState.movedArgs == nil {
				f.runState.movedArgs = make(map[rhino.Node]struct{}, 16)
			}
			f.runState.movedArgs[arg] = struct{}{}
		}
	}
	return f
}

func (f *InlineFunctions) ensureRunState() *inlineFunctionsRunState {
	if f.runState == nil {
		p := &inlineFunctionsRunState{
			fns: jutil.NewLinkedHashMap[string, *InlineFunctions_FunctionState](),
		}
		f.runState = p
		f.fns = p.fns
		f.anonFns = nil
	}
	return f.runState
}

func (f *InlineFunctions) putAnonFn(fnNode rhino.Node, name string) {
	p := f.ensureRunState()
	if p.anonFns == nil {
		p.anonFns = make(map[rhino.Node]string, 8)
		f.anonFns = p.anonFns
	}
	p.anonFns[fnNode] = name
}

func functionStateAliasCount(fs *InlineFunctions_FunctionState) int {
	if fs.summary != nil {
		return int(fs.summary.numModified)
	}
	return fs.namesToAliasIDs.size()
}

func (f *InlineFunctions) allocReference(callNode rhino.Node, scope *Scope, chunk *JSChunk, mode FunctionInjector_InliningMode) *FunctionInjector_Reference {
	r := f.ensureRunState().references.New()
	r.init(callNode, scope, chunk, mode, false)
	return r
}

// Java: InlineFunctions.getOrCreateFunctionState(String)
func (f *InlineFunctions) GetOrCreateFunctionState(fnName string) *InlineFunctions_FunctionState {
	f.ensureRunState()
	if fs := f.fns.Get(fnName); fs != nil {
		return fs
	}
	fs := &InlineFunctions_FunctionState{inline: true, remove: true}
	f.fns.PutNew(fnName, fs)
	return fs
}

func inlineFunctions_traverseWithScopeCreator(compiler AbstractCompiler, root rhino.Node, cb NodeTraversal_Callback, sc ScopeCreator) {
	b := NodeTraversal_Builder{compiler: compiler, callback: cb, scopeCreator: sc}
	b.Traverse(root)
}

// Java: InlineFunctions.process(Node, Node)
func (f *InlineFunctions) Process(externs rhino.Node, root rhino.Node) {
	jutil.CheckState(f.compiler.GetLifeCycleStage().IsNormalized())

	p := f.ensureRunState()
	defer func() {
		f.injector.fnInfoCache = nil
		f.injector.knownConstantFunctions = nil
		f.functionArgumentInjector.sideEffectCache = nil
		f.runState = nil
		f.fns = nil
		f.anonFns = nil
		f.lastCheckedBlock = rhino.NilNode
		f.lastContainingFn = rhino.NilNode
	}()

	if !root.IsNil() {
		if !f.findCandidateFunctionsParallel(root) {
			findCandidates := InlineFunctions_FindCandidateFunctions{outer: f}
			findCandidates.walk(root, rhino.NilNode, nil, false)
		}
	}
	if f.fns.IsEmpty() {
		return // Nothing left to do.
	}
	if !f.findCandidatesReferencesParallel(root) {
		sc := NewSyntacticScopeCreator(f.compiler)
		finder := inlineRefFinder{outer: f}
		visitor := inlineCallVisitor{
			functionMap:     f.fns,
			anonFunctionMap: f.anonFns,
			onCallSite:      finder.visitCallSite,
			onName:          finder.checkNameUsage,
		}
		inlineFunctions_traverseWithScopeCreator(f.compiler, root, &visitor, sc)
	}
	f.lastContainingFn = rhino.NilNode
	f.TrimCandidatesNotMeetingMinimumRequirements()
	if f.fns.IsEmpty() {
		return // Nothing left to do.
	}

	// Store the set of function names eligible for inlining and use this to
	// prevent function names from being moved into temporaries during
	// expression decomposition. If this movement were allowed it would prevent
	// the Inline callback from finding the function calls.
	//
	// This pass already assumes these are constants, so this is safe for anyone
	// using function inlining.
	//
	fnNames := make(map[string]bool, f.fns.Size())
	for k := range f.fns.Keys() {
		fnNames[k] = true
	}
	p.fnNames = fnNames
	f.injector.SetKnownConstantFunctions(fnNames)

	f.TrimCandidatesUsingOnCost()
	if f.fns.IsEmpty() {
		return // Nothing left to do.
	}
	f.ResolveInlineConflicts()
	f.DecomposeExpressions()
	if !root.IsNil() {
		callVisitor := inlineCallVisitor{
			functionMap:     f.fns,
			anonFunctionMap: f.anonFns,
			onCallSite:      f.inlineCallSite,
		}
		callVisitor.inlineWalk(root, rhino.NilNode)
	}

	f.RemoveInlinedFunctions()
}

// Java: InlineFunctions.isAlwaysInlinable(Node)
func InlineFunctions_IsAlwaysInlinable(fn rhino.Node) bool {
	body := ast.MustFunction(fn).Body()
	block, ok := body.AsBlock()
	if !ok {
		// Java reads an expression body (an arrow function) with the same
		// child tests: one without children (x => y) passes.
		return !body.Node().HasChildren()
	}
	only, ok := block.OnlyStatement()
	return !block.HasStatements() || (ok && only.Token() == rhino.Token_RETURN)
}

// inlineFunctions_varFunction returns the name and the function of a
// declaration whose first declarator is a name with a function value (var
// fooFn = function(x) {...}); ok is false for any other declaration.
func inlineFunctions_varFunction(decl ast.NameDeclaration) (name ast.Name, fn ast.Function, ok bool) {
	d, ok := decl.FirstDeclarator()
	if !ok {
		return ast.Name{}, ast.Function{}, false
	}
	if name, ok = d.AsName(); !ok {
		return ast.Name{}, ast.Function{}, false
	}
	value, ok := name.Init()
	if !ok {
		return ast.Name{}, ast.Function{}, false
	}
	fn, ok = value.AsFunction()
	return name, fn, ok
}

// inlineFunctions_calledFunction returns the function expression that the
// call calls directly: (function(){})() or (function(){}).call(...); ok is
// false otherwise.
func inlineFunctions_calledFunction(call ast.NormalOrOptChainCall) (ast.Function, bool) {
	callee := call.Callee()
	if fn, ok := callee.AsFunction(); ok {
		return fn, true
	}
	if NodeUtil_IsFunctionObjectCall(call.Node()) {
		return callee.MustNormalOrOptChainGet().Object().AsFunction()
	}
	return ast.Function{}, false
}

// inlineFunctions_walkClassMembers walks the members of a class body as
// NodeTraversal does: the computed keys first, then the values of computed
// members (a computed field's value is in a local hoist scope) and the
// other members.
func inlineFunctions_walkClassMembers(body ast.ClassBody, inLocalHoist bool, walk func(n, parent rhino.Node, inLocalHoist bool)) {
	for m := range body.Members() {
		switch m.Token() {
		case rhino.Token_COMPUTED_PROP:
			walk(m.MustComputedProp().Key().Node(), m.Node(), inLocalHoist)
		case rhino.Token_COMPUTED_FIELD_DEF:
			walk(m.MustComputedFieldDef().Key().Node(), m.Node(), inLocalHoist)
		}
	}
	for m := range body.Members() {
		switch m.Token() {
		case rhino.Token_COMPUTED_PROP:
			if value, ok := m.MustComputedProp().Value(); ok {
				walk(value.Node(), m.Node(), inLocalHoist)
			}
		case rhino.Token_COMPUTED_FIELD_DEF:
			if value, ok := m.MustComputedFieldDef().Value(); ok {
				walk(value.Node(), m.Node(), true)
			}
		default:
			walk(m.Node(), body.Node(), inLocalHoist)
		}
	}
}

// Java: InlineFunctions.targetSizeAfterInlineExceedsLimit(NodeTraversal, FunctionState)
func (f *InlineFunctions) TargetSizeAfterInlineExceedsLimit(t *NodeTraversal, functionState *InlineFunctions_FunctionState) bool {
	containingFunction := t.GetEnclosingFunction()
	// Always inline at the top level,
	// unless maybeAddFunction has marked functionState as not inlinable.
	if containingFunction.IsNil() {
		return false
	}
	inlinedFun := functionState.functionNode()
	if InlineFunctions_IsAlwaysInlinable(inlinedFun) {
		return false
	}

	inlinedFunSize := int(functionState.bodyAstSize)
	if inlinedFunSize == 0 {
		inlinedFunSize = NodeUtil_CountAstSizeUpToLimit(ast.MustFunction(inlinedFun).Body().Node(), f.maxSizeAfterInlining)
		functionState.bodyAstSize = int32(inlinedFunSize)
	}
	var targetFunSize int
	if containingFunction == f.lastContainingFn {
		targetFunSize = f.lastContainingFnSize
	} else {
		targetFunSize = NodeUtil_CountAstSizeUpToLimit(containingFunction, f.maxSizeAfterInlining)
		f.lastContainingFn = containingFunction
		f.lastContainingFnSize = targetFunSize
	}
	return inlinedFunSize+targetFunSize > f.maxSizeAfterInlining
}

// InlineFunctions_FindCandidateFunctions finds functions that might be
// inlined.
type InlineFunctions_FindCandidateFunctions struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer *InlineFunctions

	callsSeen int
}

// Java: InlineFunctions.FindCandidateFunctions()
func NewInlineFunctions_FindCandidateFunctions(outer *InlineFunctions) *InlineFunctions_FindCandidateFunctions {
	return &InlineFunctions_FindCandidateFunctions{outer: outer, callsSeen: 0}
}

func (*InlineFunctions_FindCandidateFunctions) noVisitLeaves() {}

func (c *InlineFunctions_FindCandidateFunctions) walk(n rhino.Node, parent rhino.Node, chunk *JSChunk, inLocalHoist bool) {
	tok := n.GetToken()
	childLocalHoist := inLocalHoist
	switch tok {
	case rhino.Token_SCRIPT:
		if id := n.GetInputId(); id != nil {
			if inp := c.outer.compiler.GetInput(id); inp != nil {
				chunk = inp.GetChunk()
			}
		}
	case rhino.Token_FUNCTION, rhino.Token_MODULE_BODY, rhino.Token_MEMBER_FIELD_DEF:
		childLocalHoist = true
	case rhino.Token_BLOCK:
		if !parent.IsNil() && parent.IsClassMembers() {
			childLocalHoist = true
		}
	case rhino.Token_CLASS_MEMBERS:
		inlineFunctions_walkClassMembers(ast.MustClassBody(n), inLocalHoist, func(n, parent rhino.Node, inLocalHoist bool) {
			c.walk(n, parent, chunk, inLocalHoist)
		})
		return
	}
	for ch := n.GetFirstChild(); !ch.IsNil(); ch = ch.GetNext() {
		if !ch.GetFirstChild().IsNil() {
			c.walk(ch, n, chunk, childLocalHoist)
		}
	}
	switch tok {
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST, rhino.Token_FUNCTION:
		if c.outer.reach.IncludesGlobals() || inLocalHoist {
			c.findNamedFunctionsWithChunk(n, parent, chunk)
		}
	case rhino.Token_OPTCHAIN_CALL, rhino.Token_CALL:
		if c.outer.reach.IncludesGlobals() || inLocalHoist {
			c.findFunctionExpressionsWithChunk(n, chunk)
		}
	}
}

// Java: InlineFunctions.FindCandidateFunctions.visit(NodeTraversal, Node, Node)
func (c *InlineFunctions_FindCandidateFunctions) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST, rhino.Token_FUNCTION:
		if c.outer.reach.IncludesGlobals() || !t.InGlobalHoistScope() {
			c.FindNamedFunctions(t, n, parent)
		}
	case rhino.Token_OPTCHAIN_CALL, rhino.Token_CALL:
		if c.outer.reach.IncludesGlobals() || !t.InGlobalHoistScope() {
			c.FindFunctionExpressions(t, n)
		}
	default:
	}
}

// Java: InlineFunctions.FindCandidateFunctions.findNamedFunctions(NodeTraversal, Node, Node)
func (c *InlineFunctions_FindCandidateFunctions) FindNamedFunctions(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	c.findNamedFunctionsWithChunk(n, parent, t.GetChunk())
}

func (c *InlineFunctions_FindCandidateFunctions) findNamedFunctionsWithChunk(n rhino.Node, parent rhino.Node, chunk *JSChunk) {
	switch n.GetToken() {
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		if !NodeUtil_IsStatement(n) {
			return
		}
		// Functions expressions in the form of:
		//   var fooFn = function(x) { return ... }
		decl := ast.MustNameDeclaration(n)
		jutil.CheckStateMsg1(decl.NumDeclarators() == 1, "%s", n)
		if _, _, ok := inlineFunctions_varFunction(decl); ok {
			c.outer.MaybeAddFunction(newVarFnTarget(n), chunk)
		}
		// Named functions
		// function Foo(x) { return ... }
	case rhino.Token_FUNCTION:
		if !NodeUtil_IsStatement(n) {
			return
		}
		jutil.CheckState(NodeUtil_IsStatementBlock(parent) || parent.IsLabel())
		if NodeUtil_IsFunctionDeclaration(n) {
			c.outer.MaybeAddFunction(newNamedFnTarget(n), chunk)
		}
	default:
	}
}

// FindFunctionExpressions finds function expressions that are called
// directly in the form of (function(a,b,...){...})(a,b,...) or
// (function(a,b,...){...}).call(this,a,b, ...)
//
// Java: InlineFunctions.FindCandidateFunctions.findFunctionExpressions(NodeTraversal, Node)
func (c *InlineFunctions_FindCandidateFunctions) FindFunctionExpressions(t *NodeTraversal, n rhino.Node) {
	c.findFunctionExpressionsWithChunk(n, t.GetChunk())
}

func (c *InlineFunctions_FindCandidateFunctions) findFunctionExpressionsWithChunk(n rhino.Node, chunk *JSChunk) {
	switch n.GetToken() {
	case rhino.Token_OPTCHAIN_CALL, rhino.Token_CALL:
		// Functions expressions in the form of:
		//   (function(){})();
		// If an interesting function was discovered, add it.
		if fn, ok := inlineFunctions_calledFunction(ast.MustNormalOrOptChainCall(n)); ok {
			fnNode := fn.Node()
			callsSeen := c.callsSeen
			c.callsSeen++
			target := newExprFnTarget(fnNode, callsSeen)
			c.outer.MaybeAddFunction(target, chunk)
			c.outer.putAnonFn(fnNode, target.name)
		}
	default:
	}
}

// MaybeAddFunction updates the FunctionState object for the given function.
// Checks if the given function matches the criteria for an inlinable
// function.
//
// Java: InlineFunctions.maybeAddFunction(Function, JSChunk)
func (f *InlineFunctions) MaybeAddFunction(target inlineFnTarget, chunk *JSChunk) {
	functionState := f.GetOrCreateFunctionState(target.name)
	f.UpdateFunctionStateForInlining(target, chunk, functionState)
	jutil.CheckState(!(InlineFunctions_HasRequireInliningAnnotation(target.fnNode) && !functionState.CanInline()))
}

// UpdateFunctionStateForInlining updates the FunctionState object for the
// given function. Checks if the given function matches the criteria for an
// inlinable function.
//
// Java: InlineFunctions.updateFunctionStateForInlining(Function, JSChunk, FunctionState)
func (f *InlineFunctions) UpdateFunctionStateForInlining(
	target inlineFnTarget, chunk *JSChunk, functionState *InlineFunctions_FunctionState) {
	// TODO(johnlenz): Maybe "smarten" FunctionState by adding this logic to it?

	// If the function has multiple definitions, don't inline it.
	if functionState.HasExistingFunctionDefinition() {
		functionState.DisallowInlining(
			f.compiler, InlineFunctions_DisallowInliningReason_MULTIPLE_FUNCTION_DEFINITIONS)
		return
	}
	fnNode := target.fnNode

	if f.HasNoInlineAnnotation(fnNode) {
		functionState.DisallowInlining(f.compiler, InlineFunctions_DisallowInliningReason_NO_INLINE_ANNOTATION)
		return
	}

	if f.enforceMaxSizeAfterInlining &&
		!InlineFunctions_IsAlwaysInlinable(fnNode) &&
		f.maxSizeAfterInlining <= NodeUtil_CountAstSizeUpToLimit(fnNode, f.maxSizeAfterInlining) {
		functionState.DisallowInlining(f.compiler, InlineFunctions_DisallowInliningReason_TARGET_SIZE_EXCEEDS_LIMIT)
		return
	}

	// verify the function hasn't already been marked as "don't inline"
	if functionState.CanInline() {
		// store it for use when inlining.
		functionState.setTarget(target)
		if FunctionInjector_IsDirectCallNodeReplacementPossible(fnNode) {
			functionState.InlineDirectly(true)
		}

		debugLogging := f.compiler.IsDebugLoggingEnabled()
		if !debugLogging {
			if fnNode.IsGeneratorFunction() {
				functionState.DisallowInlining(f.compiler, InlineFunctions_DisallowInliningReason_GENERATOR_FUNCTION)
				return
			}
			if fnNode.IsAsyncFunction() {
				functionState.DisallowInlining(f.compiler, InlineFunctions_DisallowInliningReason_ASYNC_FUNCTION)
				return
			}
		}

		if InlineFunctions_HasNonInlinableParam(ast.MustFunction(fnNode).Params().Node()) {
			functionState.DisallowInlining(f.compiler, InlineFunctions_DisallowInliningReason_NON_INLINABLE_PARAM)
		}

		// verify the function meets all the requirements.
		summaryValid := false
		if functionState.CanInline() || debugLogging {
			var candOk bool
			candOk, summaryValid = f.checkCandidateAndScanSummary(target, &f.scratchSummary)
			if !candOk {
				// It doesn't meet the requirements.
				functionState.DisallowInlining(f.compiler, InlineFunctions_DisallowInliningReason_NOT_CANDIDATE_FUNCTION)
			}
		}

		// Set the chunk and gather names that need temporaries.
		if functionState.CanInline() {
			functionState.SetChunk(chunk)
			f.applyBodyAnalysis(fnNode, functionState, summaryValid)
		}

		if fnNode.GetGrandparent().IsVar() &&
			(functionState.CanInline() || (functionState.RequireInlining() && debugLogging)) {
			block := functionState.target.declaringBlock()
			if block.IsBlock() &&
				!block.GetParent().IsFunction() &&
				f.blockHasLetOrConst(block) {
				// The function might capture a variable that's not in scope at the call site,
				// so don't inline.
				functionState.DisallowInlining(
					f.compiler, InlineFunctions_DisallowInliningReason_CAPTURED_VARIABLE_NOT_IN_SCOPE)
			}
		}

		if debugLogging {
			if fnNode.IsGeneratorFunction() {
				functionState.DisallowInlining(f.compiler, InlineFunctions_DisallowInliningReason_GENERATOR_FUNCTION)
			}

			if fnNode.IsAsyncFunction() {
				functionState.DisallowInlining(f.compiler, InlineFunctions_DisallowInliningReason_ASYNC_FUNCTION)
			}
		}

		if summaryValid && functionState.CanInline() {
			sumCopy := f.scratchSummary
			functionState.summary = &sumCopy
			functionState.summaryValid = true
		}
	}
}

func (f *InlineFunctions) blockHasLetOrConst(block rhino.Node) bool {
	if block == f.lastCheckedBlock {
		return f.lastCheckedBlockHasLetOrConst
	}
	res := inlineFunctions_hasLetOrConst(block)
	f.lastCheckedBlock = block
	f.lastCheckedBlockHasLetOrConst = res
	return res
}

func inlineFunctions_alwaysTrue(rhino.Node) bool { return true }

func (f *InlineFunctions) applyBodyAnalysis(fnNode rhino.Node, functionState *InlineFunctions_FunctionState, summaryValid bool) {
	if summaryValid {
		s := &f.scratchSummary
		if s.numModified > 0 {
			functionState.InlineDirectly(false)
		}
		if s.referencesThis {
			functionState.SetReferencesThis(true)
		}
		if s.hasInnerClosures {
			functionState.SetHasInnerClosures(true)
			// If there are inner functions, we can inline into global scope
			// if there are no local vars or named functions.
			// TODO(johnlenz): this can be improved by looking at the possible
			// values for locals.  If there are simple values, or constants
			// we could still inline.
			if s.numParams > 0 || s.hasShallowDeclarations {
				functionState.DisallowInlining(f.compiler, InlineFunctions_DisallowInliningReason_HAS_LOCAL_NAMES)
			}
		}
		if functionState.CanInline() && f.enforceMaxSizeAfterInlining && !InlineFunctions_IsAlwaysInlinable(fnNode) {
			functionState.bodyAstSize = int32(min(int(s.bodyAstSize), f.maxSizeAfterInlining))
		}
		return
	}

	var paramIDs, modifiedIDs inlineIDSet
	f.functionArgumentInjector.findModifiedParametersByIDInto(fnNode, &paramIDs, &modifiedIDs)
	if !modifiedIDs.isEmpty() {
		functionState.InlineDirectly(false)
		functionState.setNamesToAliasIDs(&modifiedIDs)
	}
	body := ast.MustFunction(fnNode).Body().Node()
	if NodeUtil_ReferencesEnclosingReceiver(body) {
		functionState.SetReferencesThis(true)
	}
	if NodeUtil_Has(body, NodeUtil_IsFunctionOrNonStaticClassFieldInitializer, inlineFunctions_alwaysTrue) {
		functionState.SetHasInnerClosures(true)
		if InlineFunctions_HasLocalNames(fnNode) {
			functionState.DisallowInlining(f.compiler, InlineFunctions_DisallowInliningReason_HAS_LOCAL_NAMES)
		}
	}
}

// Java: InlineFunctions.hasNoInlineAnnotation(Node)
func (f *InlineFunctions) HasNoInlineAnnotation(fnNode rhino.Node) bool {
	jsDocInfo := NodeUtil_GetBestJSDocInfo(fnNode)
	return jsDocInfo != nil && jsDocInfo.IsNoInline()
}

// Java: InlineFunctions.hasEncourageInliningAnnotation(Node)
func InlineFunctions_HasEncourageInliningAnnotation(fnNode rhino.Node) bool {
	jsDocInfo := NodeUtil_GetBestJSDocInfo(fnNode)
	return jsDocInfo != nil && jsDocInfo.IsEncourageInlining()
}

// Java: InlineFunctions.hasRequireInliningAnnotation(Node)
func InlineFunctions_HasRequireInliningAnnotation(fnNode rhino.Node) bool {
	jsDocInfo := NodeUtil_GetBestJSDocInfo(fnNode)
	return jsDocInfo != nil && jsDocInfo.IsRequireInlining()
}

func inlineFunctions_hasLetOrConst(n rhino.Node) bool {
	if n.IsLet() || n.IsConst() {
		return true
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if !c.GetFirstChild().IsNil() && inlineFunctions_hasLetOrConst(c) {
			return true
		}
	}
	return false
}

// InlineFunctions_HasLocalNames: fnNode is the function to inspect. Returns
// whether the function has parameters, var/const/let, class, or function
// declarations.
//
// Java: InlineFunctions.hasLocalNames(Node)
func InlineFunctions_HasLocalNames(fnNode rhino.Node) bool {
	fn := ast.MustFunction(fnNode)
	return fn.Params().HasParams() || functionInjector_hasShallowDeclaration(fn.Body().Node())
}

// checkCandidateAndScanSummary checks if target meets minimum candidate
// requirements while populating dst in a single body pass when possible.
func (f *InlineFunctions) checkCandidateAndScanSummary(target inlineFnTarget, dst *calleeSummary) (bool, bool) {
	conv := f.compiler.GetCodingConvention()
	if conv.IsExported(target.name) || (!target.nameNode.IsNil() && conv.IsPropertyRenameFunction(target.nameNode)) {
		return false, false
	}
	if scanCalleeSummary(target.name, target.fnNode, dst) {
		return !dst.blocksMinReqs, true
	}
	return f.injector.DoesFunctionMeetMinimumRequirements(target.name, target.fnNode), false
}

// inlineCallVisitor visits candidate call sites (and optionally NAME nodes)
// without interface or base-struct wrapper ceremony.
type inlineCallVisitor struct {
	functionMap     *jutil.LinkedHashMap[string, *InlineFunctions_FunctionState]
	anonFunctionMap map[rhino.Node]string
	onCallSite      func(t *NodeTraversal, callNode rhino.Node, fs *InlineFunctions_FunctionState)
	onName          func(n, parent rhino.Node)
}

func (*inlineCallVisitor) ShouldTraverse(*NodeTraversal, rhino.Node, rhino.Node) bool { return true }
func (*inlineCallVisitor) onlyNameLeaves()                                            {}

func (v *inlineCallVisitor) inlineWalk(n rhino.Node, parent rhino.Node) {
	if n.GetToken() == rhino.Token_CLASS_MEMBERS {
		body := ast.MustClassBody(n)
		for m := range ast.ReadAhead(body.FirstMember()) {
			switch m.Token() {
			case rhino.Token_COMPUTED_PROP:
				v.inlineWalk(m.MustComputedProp().Key().Node(), m.Node())
			case rhino.Token_COMPUTED_FIELD_DEF:
				v.inlineWalk(m.MustComputedFieldDef().Key().Node(), m.Node())
			}
		}
		for m := range ast.ReadAhead(body.FirstMember()) {
			switch m.Token() {
			case rhino.Token_COMPUTED_PROP:
				if value, ok := m.MustComputedProp().Value(); ok {
					v.inlineWalk(value.Node(), m.Node())
				}
			case rhino.Token_COMPUTED_FIELD_DEF:
				if value, ok := m.MustComputedFieldDef().Value(); ok {
					v.inlineWalk(value.Node(), m.Node())
				}
			default:
				v.inlineWalk(m.Node(), n)
			}
		}
		return
	}
	for c := n.GetFirstChild(); !c.IsNil(); {
		next := c.GetNext()
		if !c.GetFirstChild().IsNil() {
			v.inlineWalk(c, n)
		}
		c = next
	}
	switch n.GetToken() {
	case rhino.Token_OPTCHAIN_CALL, rhino.Token_CALL:
		v.visitCallNode(nil, n)
	}
}

func (v *inlineCallVisitor) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_NAME:
		if v.onName != nil {
			v.onName(n, parent)
		}
	case rhino.Token_OPTCHAIN_CALL, rhino.Token_CALL:
		v.visitCallNode(t, n)
	}
}

func (v *inlineCallVisitor) visitCallNode(t *NodeTraversal, n rhino.Node) {
	child := ast.MustNormalOrOptChainCall(n).Callee()
	var name string
	var hasName bool
	// NOTE: The normalization pass ensures that local names do not collide with global names.
	if callee, ok := child.AsName(); ok {
		name, hasName = callee.Ident(), true
	} else if child.Token() == rhino.Token_FUNCTION {
		name, hasName = v.anonFunctionMap[child.Node()]
	} else if NodeUtil_IsFunctionObjectCall(n) {
		fnIdentifyingNode := child.MustNormalOrOptChainGet().Object()
		if fnName, ok := fnIdentifyingNode.AsName(); ok {
			name, hasName = fnName.Ident(), true
		} else if fnIdentifyingNode.Token() == rhino.Token_FUNCTION {
			name, hasName = v.anonFunctionMap[fnIdentifyingNode.Node()]
		}
	}

	if hasName {
		if functionState := v.functionMap.Get(name); functionState != nil {
			v.onCallSite(t, n, functionState)
		}
	}
}

// InlineFunctions_IsCandidateUsage returns whether the name is used in a way
// that might be a candidate for inlining.
//
// Java: InlineFunctions.isCandidateUsage(Node)
func InlineFunctions_IsCandidateUsage(name rhino.Node) bool {
	jutil.CheckState(name.IsName())
	switch ast.RoleOf(name) {
	case ast.RoleNameDeclarationDeclarators, ast.RoleFunctionName, ast.RoleFunctionBody:
		// This is a declaration.  Duplicate declarations are handle during
		// function candidate gathering.
		// (Java tests that the parent is a function, so the expression body
		// of an arrow function, () => fn, counts too.)
		return true
	case ast.RoleCallCallee, ast.RoleOptChainCallCallee:
		// This is a normal reference to the function.
		return true
	case ast.RoleGetElemObject:
		// Check for a ".call" to the named function:
		//   CALL
		//     GETPROP/GETELEM
		//       NAME
		//       STRING == "call"
		//     This-Value
		//     Function-parameter-1
		//     ...
		index, ok := ast.MustGetElem(name.GetParent()).Index().AsStringLit()
		if !ok || index.Value() != "call" {
			return false
		}
	case ast.RoleGetPropObject:
		if ast.MustGetProp(name.GetParent()).Property() != "call" {
			return false
		}
	default:
		return false
	}
	// Yep, a ".call".
	return ast.RoleOf(name.GetParent()) == ast.RoleCallCallee
}

// inlineFunctions_isConstAlias reports whether the name n (whose parent is
// parent) is the value of a constant: a declared const (const bar = foo)
// or an assignment marked @const (/** @const */ exports.bar = foo).
func inlineFunctions_isConstAlias(n, parent rhino.Node) bool {
	switch ast.RoleOf(n) {
	case ast.RoleAssignValue:
		info := parent.GetJSDocInfo()
		return info != nil && info.IsConstant()
	case ast.RoleNameInit:
		return parent.GetParent().IsConst()
	}
	return false
}

// inlineRefFinder finds references to functions that are inlinable during
// sequential Phase 2.
type inlineRefFinder struct {
	outer      *InlineFunctions
	scratchRef FunctionInjector_Reference
}

func (r *inlineRefFinder) visitCallSite(t *NodeTraversal, callNode rhino.Node, functionState *InlineFunctions_FunctionState) {
	r.maybeAddReference(t, functionState, callNode, t.GetChunk())
}

func (r *inlineRefFinder) maybeAddReference(
	t *NodeTraversal, functionState *InlineFunctions_FunctionState, callNode rhino.Node, chunk *JSChunk) {
	if !functionState.CanInline() {
		return
	}

	mode := FunctionInjector_InliningMode_BLOCK
	if functionState.CanInlineDirectly() {
		mode = FunctionInjector_InliningMode_DIRECT
	}
	referenceAdded := r.maybeAddReferenceUsingMode(t, functionState, callNode, chunk, mode)
	if !referenceAdded && mode == FunctionInjector_InliningMode_DIRECT {
		// This reference can not be directly inlined, see if
		// block replacement inlining is possible.
		mode = FunctionInjector_InliningMode_BLOCK
		referenceAdded = r.maybeAddReferenceUsingMode(t, functionState, callNode, chunk, mode)
	}

	if !referenceAdded {
		// Don't try to remove a function if we can't inline all
		// the references.
		functionState.SetRemoveWithContext(
			false,
			r.outer.compiler,
			InlineFunctions_CannotRemoveReason_UNINLINABLE_REFERENCE,
			InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_NO,
			/* contextNode= */ callNode)
	}
}

func (r *inlineRefFinder) maybeAddReferenceUsingMode(
	t *NodeTraversal,
	functionState *InlineFunctions_FunctionState,
	callNode rhino.Node,
	chunk *JSChunk,
	mode FunctionInjector_InliningMode) bool {

	// If many functions are inlined into the same function F in the same
	// inlining round, then the size of F may exceed the max size.
	// This could be avoided if we bail later, during the inlining phase, eg,
	// in Inline#visitCallSite. However, that is not safe, because at that
	// point expression decomposition has already run, and we want to
	// decompose expressions only for the calls that are actually inlined.
	if r.outer.enforceMaxSizeAfterInlining && r.outer.TargetSizeAfterInlineExceedsLimit(t, functionState) {
		return false
	}

	r.scratchRef.init(callNode, nil, chunk, mode, false)
	r.scratchRef.t = t
	summary := functionState.ensureSummary()
	result := r.outer.injector.canInlineReferenceWithSummary(
		&r.scratchRef,
		functionState.functionNode(),
		functionStateAliasCount(functionState),
		functionState.GetReferencesThis(),
		functionState.HasInnerClosures(),
		summary)
	if result != FunctionInjector_CanInlineResult_NO {
		// Yeah!
		var scope *Scope
		if result == FunctionInjector_CanInlineResult_AFTER_PREPARATION {
			scope = t.GetScope()
		} else if mode != FunctionInjector_InliningMode_DIRECT && !NodeUtil_IsExprCall(callNode.GetParent()) {
			scope = t.resolveLookupScope(len(t.scopes) - 1)
		}
		candidate := r.outer.allocReference(callNode, scope, chunk, mode)
		candidate.SetRequiresDecomposition(result == FunctionInjector_CanInlineResult_AFTER_PREPARATION)
		functionState.AddReference(candidate)
		return true
	}

	return false
}

func (r *inlineRefFinder) checkNameUsage(n rhino.Node, parent rhino.Node) {
	jutil.CheckStateMsg1(n.IsName(), "%s", n)

	if n.GetStringID() == rhino.RhinoStringPool_EmptyID {
		return
	}

	if InlineFunctions_IsCandidateUsage(n) {
		return
	}

	// Other refs to a function name remove its candidacy for inlining
	functionState := r.outer.fns.Get(n.GetString())
	if functionState == nil {
		return
	}

	// If the name is being assigned to or modified as an L-value (e.g. `bar = something;`,
	// destructuring `[bar] = ...`, or `for (bar of ...)`), it cannot be inlined.
	if NodeUtil_IsLValue(n) {
		// Mark the function as uninlinable.
		functionState.DisallowInlining(r.outer.compiler, InlineFunctions_DisallowInliningReason_REASSIGNED)
	} else if inlineFunctions_isConstAlias(n, parent) {
		// e.g. const bar = foo; exports.bar = foo;
		// We can't see through this reference right now, but we will be able
		// to once `bar` is inlined; so we shouldn't remove `foo` right now but
		// this is probably okay with @requireInlining.
		functionState.SetRemove(
			false,
			r.outer.compiler,
			InlineFunctions_CannotRemoveReason_UNREMOVABLE_REFERENCE,
			InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_NO)
	} else {
		// e.g. var fn = bar; <== we can't inline "bar"
		// As this reference can't be inlined mark the function as
		// unremovable.
		functionState.SetRemoveWithContext(
			false,
			r.outer.compiler,
			InlineFunctions_CannotRemoveReason_UNREMOVABLE_REFERENCE,
			InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_YES,
			/* contextNode= */ parent)
	}
}

func (f *InlineFunctions) inlineCallSite(t *NodeTraversal, callNode rhino.Node, functionState *InlineFunctions_FunctionState) {
	jutil.CheckState(functionState.HasExistingFunctionDefinition())
	if !functionState.CanInline() {
		return
	}
	ref := functionState.GetReference(callNode)
	// There are two cases ref can be null: if the call site was introduced
	// because it was part of a function that was inlined during this pass
	// or if the call site was trimmed from the list of references because
	// the function couldn't be inlined at this location.
	if ref == nil {
		return
	}
	fnName := functionState.functionName()
	fnNode := functionState.GetSafeFnNode()
	newExpr := f.injector.inlineWithSummary(ref, fnName, fnNode, functionState.ensureSummary())
	if newExpr != ref.callNode {
		f.invalidateEnclosingSummaries(newExpr)
		if t != nil {
			t.GetCompiler().ReportChangeToEnclosingScope(newExpr)
		} else {
			f.injector.compiler.ReportChangeToEnclosingScope(newExpr)
		}
	}
	// Keep track of references that have been inlined so that
	// we can verify that none have been missed.
	ref.inlined = true
}

func (f *InlineFunctions) invalidateEnclosingSummaries(node rhino.Node) {
	if f.fns == nil || f.fns.IsEmpty() {
		return
	}
	hasAnon := len(f.anonFns) > 0
	for anc := node.GetParent(); !anc.IsNil(); anc = anc.GetParent() {
		if !anc.IsFunction() {
			continue
		}
		if name, ok := ast.MustFunction(anc).Name(); ok {
			if fs := f.fns.Get(name.Ident()); fs != nil && fs.GetSafeFnNode() == anc {
				fs.summaryValid = false
			}
		}
		if p := anc.GetParent(); p.IsName() {
			if fs := f.fns.Get(p.GetString()); fs != nil && fs.GetSafeFnNode() == anc {
				fs.summaryValid = false
			}
		}
		if hasAnon {
			if anonName, ok := f.anonFns[anc]; ok {
				if fs := f.fns.Get(anonName); fs != nil && fs.GetSafeFnNode() == anc {
					fs.summaryValid = false
				}
			}
		}
	}
}
