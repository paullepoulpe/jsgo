// Port of com/google/javascript/jscomp/InlineFunctions.java @ closure-compiler v20261015

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
	"strconv"
	"strings"
	"sync"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

const (
	ifStateChunkSize = 512
	ifRefChunkSize   = 512
	ifFnChunkSize    = 512
)

type inlineFunctionsPoolState struct {
	fns          *jutil.LinkedHashMap[string, *InlineFunctions_FunctionState]
	anonFns      *jutil.LinkedHashMap[*rhino.Node, string]
	fnNames      *jutil.LinkedHashSet[string]
	calledNames  *jutil.LinkedHashSet[string]
	refsScratch  []FunctionInjector_Reference
	paramScratch []string

	fnStateChunks [][]InlineFunctions_FunctionState
	fnStateIdx    int
	fnStateUsed   int

	refChunks [][]InlineFunctions_Reference
	refIdx    int
	refUsed   int

	namedFnChunks [][]InlineFunctions_NamedFunction
	namedFnIdx    int
	namedFnUsed   int

	fnVarChunks [][]InlineFunctions_FunctionVar
	fnVarIdx    int
	fnVarUsed   int

	fnExprChunks [][]InlineFunctions_FunctionExpression
	fnExprIdx    int
	fnExprUsed   int

	movedArgs map[*rhino.Node]struct{}
}

var ifPool sync.Pool

func acquireIfPoolState() *inlineFunctionsPoolState {
	if v := ifPool.Get(); v != nil {
		return v.(*inlineFunctionsPoolState)
	}
	return &inlineFunctionsPoolState{
		fns:         jutil.NewLinkedHashMapWithCapacity[string, *InlineFunctions_FunctionState](512),
		anonFns:     jutil.NewLinkedHashMap[*rhino.Node, string](),
		fnNames:     jutil.NewLinkedHashSet[string](),
		calledNames: jutil.NewLinkedHashSet[string](),
	}
}

func releaseIfPoolState(p *inlineFunctionsPoolState) {
	if p.fns.Size() > 262144 {
		return
	}
	p.fns.Clear()
	p.anonFns.Clear()
	p.fnNames.Clear()
	p.calledNames.Clear()
	clear(p.refsScratch)
	p.refsScratch = p.refsScratch[:0]
	clear(p.paramScratch)
	p.paramScratch = p.paramScratch[:0]
	clear(p.movedArgs)

	for ci := 0; ci <= p.fnStateIdx && ci < len(p.fnStateChunks); ci++ {
		chunk := p.fnStateChunks[ci]
		limit := len(chunk)
		if ci == p.fnStateIdx {
			limit = p.fnStateUsed
		}
		clear(chunk[:limit])
	}
	p.fnStateIdx = 0
	p.fnStateUsed = 0

	for ci := 0; ci <= p.refIdx && ci < len(p.refChunks); ci++ {
		chunk := p.refChunks[ci]
		limit := len(chunk)
		if ci == p.refIdx {
			limit = p.refUsed
		}
		clear(chunk[:limit])
	}
	p.refIdx = 0
	p.refUsed = 0

	for ci := 0; ci <= p.namedFnIdx && ci < len(p.namedFnChunks); ci++ {
		chunk := p.namedFnChunks[ci]
		limit := len(chunk)
		if ci == p.namedFnIdx {
			limit = p.namedFnUsed
		}
		clear(chunk[:limit])
	}
	p.namedFnIdx = 0
	p.namedFnUsed = 0

	for ci := 0; ci <= p.fnVarIdx && ci < len(p.fnVarChunks); ci++ {
		chunk := p.fnVarChunks[ci]
		limit := len(chunk)
		if ci == p.fnVarIdx {
			limit = p.fnVarUsed
		}
		clear(chunk[:limit])
	}
	p.fnVarIdx = 0
	p.fnVarUsed = 0

	for ci := 0; ci <= p.fnExprIdx && ci < len(p.fnExprChunks); ci++ {
		chunk := p.fnExprChunks[ci]
		limit := len(chunk)
		if ci == p.fnExprIdx {
			limit = p.fnExprUsed
		}
		clear(chunk[:limit])
	}
	p.fnExprIdx = 0
	p.fnExprUsed = 0

	ifPool.Put(p)
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
	anonFns *jutil.LinkedHashMap[*rhino.Node, string]

	compiler AbstractCompiler

	injector                 *FunctionInjector
	functionArgumentInjector *FunctionArgumentInjector

	reach                CompilerOptions_Reach
	assumeMinimumCapture bool

	enforceMaxSizeAfterInlining bool
	maxSizeAfterInlining        int

	poolState *inlineFunctionsPoolState

	lastCheckedBlock              *rhino.Node
	lastCheckedBlockHasLetOrConst bool
	lastContainingFn              *rhino.Node
	lastContainingFnSize          int
}

// Java: InlineFunctions.MISSED_REQUIRED_INLINING
var InlineFunctions_MISSED_REQUIRED_INLINING = DiagnosticType_Warning(
	"JSC_MISSED_REQUIRED_INLINING",
	"function {1} annotated @requireInlining could not be inlined here")

// Java: InlineFunctions(AbstractCompiler, Supplier<String>, Reach, boolean, boolean, int)
func NewInlineFunctions(
	compiler AbstractCompiler,
	safeNameIdSupplier func() string,
	reach CompilerOptions_Reach,
	assumeStrictThis bool,
	assumeMinimumCapture bool,
	maxSizeAfterInlining int) *InlineFunctions {
	jutil.CheckArgument(compiler != nil)
	jutil.CheckArgument(safeNameIdSupplier != nil)
	jutil.CheckArgument(reach != CompilerOptions_Reach_NONE)

	f := &InlineFunctions{}
	f.compiler = compiler

	f.reach = reach
	f.assumeMinimumCapture = assumeMinimumCapture

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
		AssumeMinimumCapture(assumeMinimumCapture).
		FunctionArgumentInjector(f.functionArgumentInjector).
		Build()
	f.injector.mutator.onMovedArg = func(arg *rhino.Node) {
		if f.poolState != nil {
			if f.poolState.movedArgs == nil {
				f.poolState.movedArgs = make(map[*rhino.Node]struct{}, 16)
			}
			f.poolState.movedArgs[arg] = struct{}{}
		}
	}
	return f
}

func (f *InlineFunctions) ensurePoolState() *inlineFunctionsPoolState {
	if f.poolState == nil {
		p := acquireIfPoolState()
		f.poolState = p
		f.fns = p.fns
		f.anonFns = p.anonFns
	}
	return f.poolState
}

func (f *InlineFunctions) allocFunctionState() *InlineFunctions_FunctionState {
	p := f.ensurePoolState()
	if len(p.fnStateChunks) == 0 {
		p.fnStateChunks = append(p.fnStateChunks, make([]InlineFunctions_FunctionState, ifStateChunkSize))
		p.fnStateIdx = 0
		p.fnStateUsed = 0
	} else if p.fnStateUsed == len(p.fnStateChunks[p.fnStateIdx]) {
		p.fnStateIdx++
		if p.fnStateIdx == len(p.fnStateChunks) {
			p.fnStateChunks = append(p.fnStateChunks, make([]InlineFunctions_FunctionState, ifStateChunkSize))
		}
		p.fnStateUsed = 0
	}
	fs := &p.fnStateChunks[p.fnStateIdx][p.fnStateUsed]
	p.fnStateUsed++
	fs.inline = true
	fs.remove = true
	return fs
}

func (f *InlineFunctions) allocReference(callNode *rhino.Node, scope *Scope, chunk *JSChunk, mode FunctionInjector_InliningMode) *InlineFunctions_Reference {
	p := f.ensurePoolState()
	if len(p.refChunks) == 0 {
		p.refChunks = append(p.refChunks, make([]InlineFunctions_Reference, ifRefChunkSize))
		p.refIdx = 0
		p.refUsed = 0
	} else if p.refUsed == len(p.refChunks[p.refIdx]) {
		p.refIdx++
		if p.refIdx == len(p.refChunks) {
			p.refChunks = append(p.refChunks, make([]InlineFunctions_Reference, ifRefChunkSize))
		}
		p.refUsed = 0
	}
	r := &p.refChunks[p.refIdx][p.refUsed]
	p.refUsed++
	r.InitFunctionInjector_ReferenceBase(callNode, scope, chunk, mode)
	return r
}

func (f *InlineFunctions) allocNamedFunction(fn *rhino.Node) *InlineFunctions_NamedFunction {
	p := f.ensurePoolState()
	if len(p.namedFnChunks) == 0 {
		p.namedFnChunks = append(p.namedFnChunks, make([]InlineFunctions_NamedFunction, ifFnChunkSize))
		p.namedFnIdx = 0
		p.namedFnUsed = 0
	} else if p.namedFnUsed == len(p.namedFnChunks[p.namedFnIdx]) {
		p.namedFnIdx++
		if p.namedFnIdx == len(p.namedFnChunks) {
			p.namedFnChunks = append(p.namedFnChunks, make([]InlineFunctions_NamedFunction, ifFnChunkSize))
		}
		p.namedFnUsed = 0
	}
	nf := &p.namedFnChunks[p.namedFnIdx][p.namedFnUsed]
	p.namedFnUsed++
	nf.outer = f
	nf.fn = fn
	return nf
}

func (f *InlineFunctions) allocFunctionVar(var_ *rhino.Node) *InlineFunctions_FunctionVar {
	p := f.ensurePoolState()
	if len(p.fnVarChunks) == 0 {
		p.fnVarChunks = append(p.fnVarChunks, make([]InlineFunctions_FunctionVar, ifFnChunkSize))
		p.fnVarIdx = 0
		p.fnVarUsed = 0
	} else if p.fnVarUsed == len(p.fnVarChunks[p.fnVarIdx]) {
		p.fnVarIdx++
		if p.fnVarIdx == len(p.fnVarChunks) {
			p.fnVarChunks = append(p.fnVarChunks, make([]InlineFunctions_FunctionVar, ifFnChunkSize))
		}
		p.fnVarUsed = 0
	}
	fv := &p.fnVarChunks[p.fnVarIdx][p.fnVarUsed]
	p.fnVarUsed++
	fv.outer = f
	fv.var_ = var_
	return fv
}

func (f *InlineFunctions) allocFunctionExpression(fn *rhino.Node, index int) *InlineFunctions_FunctionExpression {
	p := f.ensurePoolState()
	if len(p.fnExprChunks) == 0 {
		p.fnExprChunks = append(p.fnExprChunks, make([]InlineFunctions_FunctionExpression, 64))
		p.fnExprIdx = 0
		p.fnExprUsed = 0
	} else if p.fnExprUsed == len(p.fnExprChunks[p.fnExprIdx]) {
		p.fnExprIdx++
		if p.fnExprIdx == len(p.fnExprChunks) {
			p.fnExprChunks = append(p.fnExprChunks, make([]InlineFunctions_FunctionExpression, 64))
		}
		p.fnExprUsed = 0
	}
	fe := &p.fnExprChunks[p.fnExprIdx][p.fnExprUsed]
	p.fnExprUsed++
	fe.fn = fn
	fe.fakeName = strconv.Itoa(index)
	return fe
}

// Java: InlineFunctions.getOrCreateFunctionState(String)
func (f *InlineFunctions) GetOrCreateFunctionState(fnName string) *InlineFunctions_FunctionState {
	f.ensurePoolState()
	if fs := f.fns.Get(fnName); fs != nil {
		return fs
	}
	fs := f.allocFunctionState()
	f.fns.PutNew(fnName, fs)
	return fs
}

func inlineFunctions_traverseWithScopeCreator(compiler AbstractCompiler, root *rhino.Node, cb NodeTraversal_Callback, sc ScopeCreator) {
	b := NodeTraversal_Builder{compiler: compiler, callback: cb, scopeCreator: sc}
	t := acquireNodeTraversal(&b)
	t.Traverse(root)
	releaseNodeTraversal(t)
}

// Java: InlineFunctions.process(Node, Node)
func (f *InlineFunctions) Process(externs *rhino.Node, root *rhino.Node) {
	jutil.CheckState(f.compiler.GetLifeCycleStage().IsNormalized())

	p := f.ensurePoolState()
	defer func() {
		f.poolState = nil
		f.fns = nil
		f.anonFns = nil
		f.lastCheckedBlock = nil
		f.lastContainingFn = nil
		releaseIfPoolState(p)
	}()

	sc := NewSyntacticScopeCreator(f.compiler)
	sc.EnableArenaPool()
	defer sc.ReleaseArena()

	findCandidates := InlineFunctions_FindCandidateFunctions{outer: f}
	inlineFunctions_traverseWithScopeCreator(f.compiler, root, &findCandidates, sc)
	if f.fns.IsEmpty() {
		return // Nothing left to do.
	}
	findRefs := InlineFunctions_FindCandidatesReferences{outer: f}
	findRefs.InitInlineFunctions_CallVisitorBase(f.fns, f.anonFns, &findRefs)
	inlineFunctions_traverseWithScopeCreator(f.compiler, root, &findRefs, sc)
	f.lastContainingFn = nil
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
	fnNames := p.fnNames
	fnNames.Clear()
	for k := range f.fns.Keys() {
		fnNames.Add(k)
	}
	f.injector.SetKnownConstantFunctions(fnNames)

	f.TrimCandidatesUsingOnCost()
	if f.fns.IsEmpty() {
		return // Nothing left to do.
	}
	f.ResolveInlineConflicts()
	f.DecomposeExpressions()
	inlineCb := InlineFunctions_Inline{injector: f.injector}
	callVisitor := InlineFunctions_CallVisitorBase{}
	callVisitor.InitInlineFunctions_CallVisitorBase(f.fns, f.anonFns, &inlineCb)
	inlineFunctions_traverseWithScopeCreator(f.compiler, root, &callVisitor, sc)

	f.RemoveInlinedFunctions()
}

// Java: InlineFunctions.isAlwaysInlinable(Node)
func InlineFunctions_IsAlwaysInlinable(fn *rhino.Node) bool {
	jutil.CheckArgument(fn.IsFunction())
	body := NodeUtil_GetFunctionBody(fn)
	return (!body.HasChildren()) || (body.HasOneChild() && body.GetFirstChild().IsReturn())
}

// Java: InlineFunctions.targetSizeAfterInlineExceedsLimit(NodeTraversal, FunctionState)
func (f *InlineFunctions) TargetSizeAfterInlineExceedsLimit(t *NodeTraversal, functionState *InlineFunctions_FunctionState) bool {
	containingFunction := t.GetEnclosingFunction()
	// Always inline at the top level,
	// unless maybeAddFunction has marked functionState as not inlinable.
	if containingFunction == nil {
		return false
	}
	inlinedFun := functionState.GetFn().GetFunctionNode()
	if InlineFunctions_IsAlwaysInlinable(inlinedFun) {
		return false
	}

	inlinedFunSize := int(functionState.bodyAstSize)
	if inlinedFunSize == 0 {
		inlinedFunSize = NodeUtil_CountAstSizeUpToLimit(NodeUtil_GetFunctionBody(inlinedFun), f.maxSizeAfterInlining)
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

// Java: InlineFunctions.FindCandidateFunctions.visit(NodeTraversal, Node, Node)
func (c *InlineFunctions_FindCandidateFunctions) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
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
func (c *InlineFunctions_FindCandidateFunctions) FindNamedFunctions(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		if !NodeUtil_IsStatement(n) {
			return
		}
		// Functions expressions in the form of:
		//   var fooFn = function(x) { return ... }
		jutil.CheckStateMsg1(n.HasOneChild(), "%s", n)
		nameNode := n.GetFirstChild()
		if nameNode.IsName() &&
			nameNode.HasChildren() &&
			nameNode.GetFirstChild().IsFunction() {
			c.outer.MaybeAddFunction(c.outer.allocFunctionVar(n), t.GetChunk())
		}
		// Named functions
		// function Foo(x) { return ... }
	case rhino.Token_FUNCTION:
		if !NodeUtil_IsStatement(n) {
			return
		}
		jutil.CheckState(NodeUtil_IsStatementBlock(parent) || parent.IsLabel())
		if NodeUtil_IsFunctionDeclaration(n) {
			c.outer.MaybeAddFunction(c.outer.allocNamedFunction(n), t.GetChunk())
		}
	default:
	}
}

// FindFunctionExpressions finds function expressions that are called
// directly in the form of (function(a,b,...){...})(a,b,...) or
// (function(a,b,...){...}).call(this,a,b, ...)
//
// Java: InlineFunctions.FindCandidateFunctions.findFunctionExpressions(NodeTraversal, Node)
func (c *InlineFunctions_FindCandidateFunctions) FindFunctionExpressions(t *NodeTraversal, n *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_OPTCHAIN_CALL, rhino.Token_CALL:
		// Functions expressions in the form of:
		//   (function(){})();
		var fnNode *rhino.Node
		first := n.GetFirstChild()
		if first.IsFunction() {
			fnNode = first
		} else if first.HasChildren() && first.GetFirstChild().IsFunction() && NodeUtil_IsFunctionObjectCall(n) {
			fnNode = first.GetFirstChild()
		}

		// If an interesting function was discovered, add it.
		if fnNode != nil {
			callsSeen := c.callsSeen
			c.callsSeen++
			fn := c.outer.allocFunctionExpression(fnNode, callsSeen)
			c.outer.MaybeAddFunction(fn, t.GetChunk())
			c.outer.anonFns.Put(fnNode, fn.GetName())
		}
	default:
	}
}

// MaybeAddFunction updates the FunctionState object for the given function.
// Checks if the given function matches the criteria for an inlinable
// function.
//
// Java: InlineFunctions.maybeAddFunction(Function, JSChunk)
func (f *InlineFunctions) MaybeAddFunction(fn InlineFunctions_Function, chunk *JSChunk) {
	name := fn.GetName()
	functionState := f.GetOrCreateFunctionState(name)
	f.UpdateFunctionStateForInlining(fn, chunk, functionState)
	jutil.CheckState(!(InlineFunctions_HasRequireInliningAnnotation(fn.GetFunctionNode()) && !functionState.CanInline()))
}

// UpdateFunctionStateForInlining updates the FunctionState object for the
// given function. Checks if the given function matches the criteria for an
// inlinable function.
//
// Java: InlineFunctions.updateFunctionStateForInlining(Function, JSChunk, FunctionState)
func (f *InlineFunctions) UpdateFunctionStateForInlining(
	fn InlineFunctions_Function, chunk *JSChunk, functionState *InlineFunctions_FunctionState) {
	// TODO(johnlenz): Maybe "smarten" FunctionState by adding this logic to it?

	// If the function has multiple definitions, don't inline it.
	if functionState.HasExistingFunctionDefinition() {
		functionState.DisallowInlining(
			f.compiler, InlineFunctions_DisallowInliningReason_MULTIPLE_FUNCTION_DEFINITIONS)
		return
	}
	fnNode := fn.GetFunctionNode()

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
		functionState.SetFn(fn)
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

		if InlineFunctions_HasNonInlinableParam(NodeUtil_GetFunctionParameters(fnNode)) {
			functionState.DisallowInlining(f.compiler, InlineFunctions_DisallowInliningReason_NON_INLINABLE_PARAM)
		}

		// verify the function meets all the requirements.
		// TODO(johnlenz): Minimum requirement checks are about 5% of the
		// run-time cost of this pass.
		if (functionState.CanInline() || debugLogging) && !f.IsCandidateFunction(fn) {
			// It doesn't meet the requirements.
			functionState.DisallowInlining(f.compiler, InlineFunctions_DisallowInliningReason_NOT_CANDIDATE_FUNCTION)
		}

		// Set the chunk and gather names that need temporaries.
		if functionState.CanInline() {
			functionState.SetChunk(chunk)

			namesToAlias, referencesThis, hasInnerClosures := f.analyzeFunctionBody(fnNode)
			if namesToAlias != nil && !namesToAlias.IsEmpty() {
				functionState.InlineDirectly(false)
				functionState.SetNamesToAlias(namesToAlias)
			}

			if referencesThis {
				functionState.SetReferencesThis(true)
			}

			if hasInnerClosures {
				functionState.SetHasInnerClosures(true)
				// If there are inner functions, we can inline into global scope
				// if there are no local vars or named functions.
				// TODO(johnlenz): this can be improved by looking at the possible
				// values for locals.  If there are simple values, or constants
				// we could still inline.
				if !f.assumeMinimumCapture && InlineFunctions_HasLocalNames(fnNode) {
					functionState.DisallowInlining(f.compiler, InlineFunctions_DisallowInliningReason_HAS_LOCAL_NAMES)
				}
			}
		}

		if fnNode.GetGrandparent().IsVar() &&
			(functionState.CanInline() || (functionState.RequireInlining() && debugLogging)) {
			block := functionState.GetFn().GetDeclaringBlock()
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
	}
}

func (f *InlineFunctions) blockHasLetOrConst(block *rhino.Node) bool {
	if block == f.lastCheckedBlock {
		return f.lastCheckedBlockHasLetOrConst
	}
	res := inlineFunctions_hasLetOrConst(block)
	f.lastCheckedBlock = block
	f.lastCheckedBlockHasLetOrConst = res
	return res
}

type inlineFunctions_bodyAnalysis struct {
	namesToAlias     *jutil.LinkedHashSet[string]
	referencesThis   bool
	hasInnerClosures bool
}

func (f *InlineFunctions) analyzeFunctionBody(fnNode *rhino.Node) (*jutil.LinkedHashSet[string], bool, bool) {
	p := f.ensurePoolState()
	params := NodeUtil_GetFunctionParameters(fnNode)
	p.paramScratch = p.paramScratch[:0]
	for n := params.GetFirstChild(); n != nil; n = n.GetNext() {
		if n.IsRest() {
			p.paramScratch = append(p.paramScratch, FunctionArgumentInjector_REST_MARKER)
		} else if n.IsDefaultValue() || n.IsObjectPattern() || n.IsArrayPattern() {
			panic(jutil.NewIllegalStateExceptionMsg("Not supported: " + n.String()))
		} else {
			p.paramScratch = append(p.paramScratch, n.GetString())
		}
	}
	var res inlineFunctions_bodyAnalysis
	inlineFunctions_analyzeBodyNode(NodeUtil_GetFunctionBody(fnNode), p.paramScratch, false, true, &res)
	return res.namesToAlias, res.referencesThis, res.hasInnerClosures
}

func inlineFunctions_analyzeBodyNode(
	n *rhino.Node,
	paramNames []string,
	inInnerClosure bool,
	inReceiverScope bool,
	res *inlineFunctions_bodyAnalysis,
) {
	if NodeUtil_IsFunctionOrNonStaticClassFieldInitializer(n) {
		inInnerClosure = true
		res.hasInnerClosures = true
		if NodeUtil_IsNonArrowFunction(n) {
			inReceiverScope = false
		}
	}

	if inReceiverScope && !res.referencesThis && (n.IsThis() || n.IsSuper()) {
		res.referencesThis = true
	}

	if len(paramNames) > 0 && n.IsName() {
		name := n.GetString()
		for _, p := range paramNames {
			if p == name {
				if inInnerClosure || FunctionArgumentInjector_CanNameValueChange(n) {
					if res.namesToAlias == nil {
						res.namesToAlias = jutil.NewLinkedHashSet[string]()
					}
					res.namesToAlias.Add(name)
				}
				break
			}
		}
	}

	if (!inReceiverScope || res.referencesThis) && res.hasInnerClosures && len(paramNames) == 0 {
		return
	}

	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		inlineFunctions_analyzeBodyNode(c, paramNames, inInnerClosure, inReceiverScope, res)
	}
}

// Java: InlineFunctions.hasNoInlineAnnotation(Node)
func (f *InlineFunctions) HasNoInlineAnnotation(fnNode *rhino.Node) bool {
	jsDocInfo := NodeUtil_GetBestJSDocInfo(fnNode)
	return jsDocInfo != nil && jsDocInfo.IsNoInline()
}

// Java: InlineFunctions.hasEncourageInliningAnnotation(Node)
func InlineFunctions_HasEncourageInliningAnnotation(fnNode *rhino.Node) bool {
	jsDocInfo := NodeUtil_GetBestJSDocInfo(fnNode)
	return jsDocInfo != nil && jsDocInfo.IsEncourageInlining()
}

// Java: InlineFunctions.hasRequireInliningAnnotation(Node)
func InlineFunctions_HasRequireInliningAnnotation(fnNode *rhino.Node) bool {
	jsDocInfo := NodeUtil_GetBestJSDocInfo(fnNode)
	return jsDocInfo != nil && jsDocInfo.IsRequireInlining()
}

func inlineFunctions_hasLetOrConst(n *rhino.Node) bool {
	if n.IsLet() || n.IsConst() {
		return true
	}
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		if inlineFunctions_hasLetOrConst(c) {
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
func InlineFunctions_HasLocalNames(fnNode *rhino.Node) bool {
	block := NodeUtil_GetFunctionBody(fnNode)
	return NodeUtil_GetFunctionParameters(fnNode).HasChildren() ||
		functionInjector_hasShallowDeclaration(block)
}

// IsCandidateFunction checks if the given function matches the criteria for
// an inlinable function.
//
// Java: InlineFunctions.isCandidateFunction(Function)
func (f *InlineFunctions) IsCandidateFunction(fn InlineFunctions_Function) bool {
	// Don't inline exported functions.
	fnName := fn.GetName()
	if f.compiler.GetCodingConvention().IsExported(fnName) {
		// TODO(johnlenz): Should we allow internal references to be inlined?
		// An exported name can be replaced externally, any inlined instance
		// would not reflect this change.
		// To allow inlining we need to be able to distinguish between exports
		// that are used in a read-only fashion and those that can be replaced
		// by external definitions.
		return false
	}

	// Don't inline this special function
	if f.compiler.GetCodingConvention().IsPropertyRenameFunction(fn.GetNameNode()) {
		return false
	}

	fnNode := fn.GetFunctionNode()
	return f.injector.DoesFunctionMeetMinimumRequirements(fnName, fnNode)
}

// InlineFunctions_CallVisitorCallback: see CallVisitor.
type InlineFunctions_CallVisitorCallback interface {
	// Java: InlineFunctions.CallVisitorCallback.visitCallSite(NodeTraversal, Node, FunctionState)
	VisitCallSite(t *NodeTraversal, callNode *rhino.Node, functionState *InlineFunctions_FunctionState)
}

// InlineFunctions_CallVisitor visits call sites for functions in
// functionMap.
//
// It has a subclass (FindCandidatesReferences), so it is an interface
// (PORTING.md §4); InlineFunctions_CallVisitorBase holds its fields and
// implementation.
type InlineFunctions_CallVisitor interface {
	NodeTraversal_Callback

	GetInlineFunctions_CallVisitorBase() *InlineFunctions_CallVisitorBase
}

type ifLookupCacheEntry struct {
	name  string
	state *InlineFunctions_FunctionState
	valid bool
}

// InlineFunctions_CallVisitorBase holds CallVisitor's fields and its visit.
type InlineFunctions_CallVisitorBase struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	callback        InlineFunctions_CallVisitorCallback
	functionMap     *jutil.LinkedHashMap[string, *InlineFunctions_FunctionState]
	anonFunctionMap *jutil.LinkedHashMap[*rhino.Node, string]
	lookupCache     [512]ifLookupCacheEntry
	lastMissName    string
	hasLastMiss     bool
}

func (v *InlineFunctions_CallVisitorBase) lookupFunctionState(name string) *InlineFunctions_FunctionState {
	if v.hasLastMiss && name == v.lastMissName {
		return nil
	}
	if name != "" {
		idx := ocCacheIdx(name)
		if e := &v.lookupCache[idx]; e.valid && e.name == name {
			if e.state == nil {
				v.lastMissName = name
				v.hasLastMiss = true
			}
			return e.state
		}
		fs := v.functionMap.Get(name)
		v.lookupCache[idx] = ifLookupCacheEntry{name: name, state: fs, valid: true}
		if fs == nil {
			v.lastMissName = name
			v.hasLastMiss = true
		}
		return fs
	}
	return v.functionMap.Get(name)
}

// Java: InlineFunctions.CallVisitor(Map<String, FunctionState>, Map<Node, String>, CallVisitorCallback)
func NewInlineFunctions_CallVisitor(
	fns *jutil.LinkedHashMap[string, *InlineFunctions_FunctionState],
	anonFns *jutil.LinkedHashMap[*rhino.Node, string],
	callback InlineFunctions_CallVisitorCallback) *InlineFunctions_CallVisitorBase {
	v := &InlineFunctions_CallVisitorBase{}
	v.InitInlineFunctions_CallVisitorBase(fns, anonFns, callback)
	return v
}

// InitInlineFunctions_CallVisitorBase is the CallVisitor constructor body
// (super(...) for subclasses).
//
// Java: InlineFunctions.CallVisitor(Map<String, FunctionState>, Map<Node, String>, CallVisitorCallback)
func (v *InlineFunctions_CallVisitorBase) InitInlineFunctions_CallVisitorBase(
	fns *jutil.LinkedHashMap[string, *InlineFunctions_FunctionState],
	anonFns *jutil.LinkedHashMap[*rhino.Node, string],
	callback InlineFunctions_CallVisitorCallback) {
	v.functionMap = fns
	v.anonFunctionMap = anonFns
	v.callback = callback
}

// GetInlineFunctions_CallVisitorBase gives access to the base through the
// interface.
func (v *InlineFunctions_CallVisitorBase) GetInlineFunctions_CallVisitorBase() *InlineFunctions_CallVisitorBase {
	return v
}

// Java: InlineFunctions.CallVisitor.visit(NodeTraversal, Node, Node)
func (v *InlineFunctions_CallVisitorBase) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_OPTCHAIN_CALL, rhino.Token_CALL:
		// Function calls
		child := n.GetFirstChild()
		var name string
		var hasName bool
		// NOTE: The normalization pass ensures that local names do not collide with global names.
		if child.IsName() {
			name, hasName = child.GetString(), true
		} else if child.IsFunction() {
			name, hasName = v.anonFunctionMap.Lookup(child)
		} else if NodeUtil_IsFunctionObjectCall(n) {
			jutil.CheckState(NodeUtil_IsNormalOrOptChainGet(child))
			fnIdentifyingNode := child.GetFirstChild()
			if fnIdentifyingNode.IsName() {
				name, hasName = fnIdentifyingNode.GetString(), true
			} else if fnIdentifyingNode.IsFunction() {
				name, hasName = v.anonFunctionMap.Lookup(fnIdentifyingNode)
			}
		}

		if hasName {
			functionState := v.lookupFunctionState(name)

			// Only visit call-sites for functions that can be inlined.
			if functionState != nil {
				v.callback.VisitCallSite(t, n, functionState)
			}
		}
	default:
	}
}

// InlineFunctions_IsCandidateUsage returns whether the name is used in a way
// that might be a candidate for inlining.
//
// Java: InlineFunctions.isCandidateUsage(Node)
func InlineFunctions_IsCandidateUsage(name *rhino.Node) bool {
	parent := name.GetParent()
	jutil.CheckState(name.IsName())
	if NodeUtil_IsNameDeclaration(parent) || parent.IsFunction() {
		// This is a declaration.  Duplicate declarations are handle during
		// function candidate gathering.
		return true
	}

	if NodeUtil_IsNormalOrOptChainCall(parent) && parent.GetFirstChild() == name {
		// This is a normal reference to the function.
		return true
	}

	// Check for a ".call" to the named function:
	//   CALL
	//     GETPROP/GETELEM
	//       NAME
	//       STRING == "call"
	//     This-Value
	//     Function-parameter-1
	//     ...
	if (parent.IsGetElem() &&
		name == parent.GetFirstChild() &&
		parent.GetSecondChild().IsStringLit() &&
		parent.GetSecondChild().GetString() == "call") ||
		(parent.IsGetProp() && parent.GetString() == "call") {
		grandparent := name.GetAncestor(2)
		if grandparent.IsCall() && grandparent.GetFirstChild() == parent {
			// Yep, a ".call".
			return true
		}
	}
	return false
}

// InlineFunctions_FindCandidatesReferences finds references to functions
// that are inlinable.
type InlineFunctions_FindCandidatesReferences struct {
	InlineFunctions_CallVisitorBase
	outer      *InlineFunctions
	scratchRef InlineFunctions_Reference
}

// Java: InlineFunctions.FindCandidatesReferences(Map<String, FunctionState>, Map<Node, String>)
func NewInlineFunctions_FindCandidatesReferences(
	outer *InlineFunctions,
	fns *jutil.LinkedHashMap[string, *InlineFunctions_FunctionState],
	anonFns *jutil.LinkedHashMap[*rhino.Node, string]) *InlineFunctions_FindCandidatesReferences {
	r := &InlineFunctions_FindCandidatesReferences{outer: outer}
	r.InitInlineFunctions_CallVisitorBase(fns, anonFns, nil)
	r.callback = r
	return r
}

// Java: InlineFunctions.FindCandidatesReferences.visit(NodeTraversal, Node, Node)
func (r *InlineFunctions_FindCandidatesReferences) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	r.InlineFunctions_CallVisitorBase.Visit(t, n, parent)
	if n.IsName() {
		r.CheckNameUsage(n, parent)
	}
}

// Java: InlineFunctions.FindCandidatesReferences.visitCallSite(NodeTraversal, Node, FunctionState)
func (r *InlineFunctions_FindCandidatesReferences) VisitCallSite(t *NodeTraversal, callNode *rhino.Node, functionState *InlineFunctions_FunctionState) {
	r.MaybeAddReference(t, functionState, callNode, t.GetChunk())
}

// Java: InlineFunctions.FindCandidatesReferences.maybeAddReference(NodeTraversal, FunctionState, Node, JSChunk)
func (r *InlineFunctions_FindCandidatesReferences) MaybeAddReference(
	t *NodeTraversal, functionState *InlineFunctions_FunctionState, callNode *rhino.Node, chunk *JSChunk) {
	if !functionState.CanInline() {
		return
	}

	mode := FunctionInjector_InliningMode_BLOCK
	if functionState.CanInlineDirectly() {
		mode = FunctionInjector_InliningMode_DIRECT
	}
	referenceAdded := r.MaybeAddReferenceUsingMode(t, functionState, callNode, chunk, mode)
	if !referenceAdded && mode == FunctionInjector_InliningMode_DIRECT {
		// This reference can not be directly inlined, see if
		// block replacement inlining is possible.
		mode = FunctionInjector_InliningMode_BLOCK
		referenceAdded = r.MaybeAddReferenceUsingMode(t, functionState, callNode, chunk, mode)
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

// Java: InlineFunctions.FindCandidatesReferences.maybeAddReferenceUsingMode(NodeTraversal, FunctionState, Node, JSChunk, InliningMode)
func (r *InlineFunctions_FindCandidatesReferences) MaybeAddReferenceUsingMode(
	t *NodeTraversal,
	functionState *InlineFunctions_FunctionState,
	callNode *rhino.Node,
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

	var scope *Scope
	if mode != FunctionInjector_InliningMode_DIRECT || (!r.outer.assumeMinimumCapture && functionState.HasInnerClosures()) {
		scope = t.GetScope()
	}
	r.scratchRef.requiresDecomposition = false
	r.scratchRef.inlined = false
	r.scratchRef.InitFunctionInjector_ReferenceBase(callNode, scope, chunk, mode)
	result := r.outer.injector.CanInlineReferenceToFunction(
		&r.scratchRef,
		functionState.GetFn().GetFunctionNode(),
		functionState.GetNamesToAlias(),
		functionState.GetReferencesThis(),
		functionState.HasInnerClosures())
	if result != FunctionInjector_CanInlineResult_NO {
		// Yeah!
		candidate := r.outer.allocReference(callNode, scope, chunk, mode)
		candidate.SetRequiresDecomposition(result == FunctionInjector_CanInlineResult_AFTER_PREPARATION)
		functionState.AddReference(candidate)
		return true
	}

	return false
}

// CheckNameUsage finds functions that can be inlined.
//
// Java: InlineFunctions.FindCandidatesReferences.checkNameUsage(Node, Node)
func (r *InlineFunctions_FindCandidatesReferences) CheckNameUsage(n *rhino.Node, parent *rhino.Node) {
	jutil.CheckStateMsg1(n.IsName(), "%s", n)

	name := n.GetString()
	if name == "" || (r.hasLastMiss && name == r.lastMissName) {
		return
	}

	if InlineFunctions_IsCandidateUsage(n) {
		return
	}

	// Other refs to a function name remove its candidacy for inlining
	functionState := r.lookupFunctionState(name)
	if functionState == nil {
		return
	}

	// If the name is being assigned to or modified as an L-value (e.g. `bar = something;`,
	// destructuring `[bar] = ...`, or `for (bar of ...)`), it cannot be inlined.
	if NodeUtil_IsLValue(n) {
		// Mark the function as uninlinable.
		functionState.DisallowInlining(r.outer.compiler, InlineFunctions_DisallowInliningReason_REASSIGNED)
	} else if (parent.IsAssign() &&
		parent.GetJSDocInfo() != nil &&
		parent.GetJSDocInfo().IsConstant() &&
		parent.GetSecondChild() == n) ||
		(parent.GetParent() != nil &&
			parent.GetParent().IsConst() &&
			parent.GetFirstChild() == n) {
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

// InlineFunctions_Inline inlines functions at the call sites.
type InlineFunctions_Inline struct {
	injector *FunctionInjector
}

// Java: InlineFunctions.Inline(FunctionInjector)
func NewInlineFunctions_Inline(injector *FunctionInjector) *InlineFunctions_Inline {
	return &InlineFunctions_Inline{injector: injector}
}

// Java: InlineFunctions.Inline.visitCallSite(NodeTraversal, Node, FunctionState)
func (i *InlineFunctions_Inline) VisitCallSite(t *NodeTraversal, callNode *rhino.Node, functionState *InlineFunctions_FunctionState) {
	jutil.CheckState(functionState.HasExistingFunctionDefinition())
	if functionState.CanInline() {
		ref := functionState.GetReference(callNode)

		// There are two cases ref can be null: if the call site was introduced
		// because it was part of a function that was inlined during this pass
		// or if the call site was trimmed from the list of references because
		// the function couldn't be inlined at this location.
		if ref != nil {
			i.InlineFunction(t, ref, functionState)
			// Keep track of references that have been inlined so that
			// we can verify that none have been missed.
			ref.inlined = true
		}
	}
}

// InlineFunction inlines a function into the call site.
//
// Java: InlineFunctions.Inline.inlineFunction(NodeTraversal, Reference, FunctionState)
func (i *InlineFunctions_Inline) InlineFunction(t *NodeTraversal, ref *InlineFunctions_Reference, functionState *InlineFunctions_FunctionState) {
	fn := functionState.GetFn()
	fnName := fn.GetName()
	fnNode := functionState.GetSafeFnNode()

	newExpr := i.injector.Inline(ref, fnName, fnNode)
	if newExpr != ref.callNode {
		t.GetCompiler().ReportChangeToEnclosingScope(newExpr)
	}
}

// TrimCandidatesNotMeetingMinimumRequirements removes entries that aren't a
// valid inline candidates, from the list of encountered names.
//
// Java: InlineFunctions.trimCandidatesNotMeetingMinimumRequirements()
func (f *InlineFunctions) TrimCandidatesNotMeetingMinimumRequirements() {
	for key, functionState := range f.fns.All() {
		if !functionState.HasExistingFunctionDefinition() || !functionState.CanInline() {
			f.fns.Remove(key)
		}
	}
}

// TrimCandidatesUsingOnCost removes entries from the list of candidates that
// can't be inlined.
//
// Java: InlineFunctions.trimCandidatesUsingOnCost()
func (f *InlineFunctions) TrimCandidatesUsingOnCost() {
	for key, functionState := range f.fns.All() {
		if functionState.HasReferences() {
			// Only inline function if it decreases the code size.
			lowersCost := f.MinimizeCost(functionState)
			if !lowersCost {
				// It shouldn't be inlined; remove it from the list.
				f.fns.Remove(key)
			}
		} else if !functionState.CanRemove() {
			// Don't bother tracking functions without references that can't be
			// removed.
			f.fns.Remove(key)
		}
	}
}

// MinimizeCost determines if the function is worth inlining and potentially
// trims references that increase the cost.
//
// Returns whether inlining the references lowers the overall cost.
//
// Java: InlineFunctions.minimizeCost(FunctionState)
func (f *InlineFunctions) MinimizeCost(functionState *InlineFunctions_FunctionState) bool {
	if !f.InliningLowersCost(functionState) {
		// Try again without Block inlining references
		if functionState.HasBlockInliningReferences() {
			functionState.SetRemove(
				false,
				f.compiler,
				InlineFunctions_CannotRemoveReason_BLOCK_INLINING_REFERENCE,
				InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_YES)
			functionState.RemoveBlockInliningReferences()
			if !functionState.HasReferences() || !f.InliningLowersCost(functionState) {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

// InliningLowersCost returns whether inlining the function reduces code
// size.
//
// Java: InlineFunctions.inliningLowersCost(FunctionState)
func (f *InlineFunctions) InliningLowersCost(functionState *InlineFunctions_FunctionState) bool {
	if functionState.RequireInlining() || functionState.EncourageInlining() {
		return true
	}
	p := f.ensurePoolState()
	p.refsScratch = p.refsScratch[:0]
	if functionState.references != nil {
		for ref := range functionState.references.Values() {
			p.refsScratch = append(p.refsScratch, ref)
		}
	}
	return f.injector.InliningLowersCost(
		functionState.GetChunk(),
		functionState.GetFn().GetFunctionNode(),
		p.refsScratch,
		functionState.GetNamesToAlias(),
		functionState.CanRemove(),
		functionState.GetReferencesThis())
}

// ResolveInlineConflicts: size base inlining calculations are thrown off
// when a function that is being inlined also contains calls to functions
// that are slated for inlining.
//
// Specifically, a clone of the FUNCTION node tree is used when the function
// is inlined. Calls in this new tree are not included in the list of
// function references so they won't be inlined (which is what we want).
// Here we mark those functions as non-removable (as they will have new
// references in the cloned node trees).
//
// This prevents a function that would only be inlined because it is
// referenced once from being inlined into multiple call sites because the
// calling function has been inlined in multiple locations or the function
// being removed while there are still references.
//
// Java: InlineFunctions.resolveInlineConflicts()
func (f *InlineFunctions) ResolveInlineConflicts() {
	for functionState := range f.fns.Values() {
		f.ResolveInlineConflictsForFunction(functionState)
	}
}

func inlineFunctions_hasNonInlinableParamNode(n *rhino.Node) bool {
	if n.IsDefaultValue() || n.IsDestructuringPattern() {
		return true
	}
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		if inlineFunctions_hasNonInlinableParamNode(c) {
			return true
		}
	}
	return false
}

// InlineFunctions_HasNonInlinableParam returns whether the function has any
// parameters that would stop the compiler from inlining. Currently this
// includes object patterns, array patterns, and default values.
//
// Java: InlineFunctions.hasNonInlinableParam(Node)
func InlineFunctions_HasNonInlinableParam(node *rhino.Node) bool {
	jutil.CheckNotNull(node)
	for c := node.GetFirstChild(); c != nil; c = c.GetNext() {
		if !c.IsName() || c.HasChildren() {
			return inlineFunctions_hasNonInlinableParamNode(node)
		}
	}
	return false
}

// ResolveInlineConflictsForFunction: see resolveInlineConflicts.
//
// Java: InlineFunctions.resolveInlineConflictsForFunction(FunctionState)
func (f *InlineFunctions) ResolveInlineConflictsForFunction(functionState *InlineFunctions_FunctionState) {
	// Functions that aren't referenced don't cause conflicts.
	if !functionState.HasReferences() || !functionState.CanInline() {
		return
	}

	fnNode := functionState.GetFn().GetFunctionNode()
	p := f.ensurePoolState()
	names := p.calledNames
	names.Clear()
	inlineFunctions_findCalledFunctions(NodeUtil_GetFunctionBody(fnNode), names)
	if !names.IsEmpty() {
		// Prevent the removal of the referenced functions unless they will be immediately inlined.
		for name := range names.All() {
			fsCalled := f.fns.Get(name)
			if fsCalled != nil && fsCalled.CanRemove() {
				fsCalled.SetRemove(
					false,
					f.compiler,
					InlineFunctions_CannotRemoveReason_CALLED_BY_INLINED_FUNCTION,
					InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_NO)
				// For functions that can no longer be removed, check if they should
				// still be inlined.
				if !f.MinimizeCost(fsCalled) &&
					!fsCalled.RequireInlining() &&
					!fsCalled.EncourageInlining() {
					// It can't be inlined remove it from the list.
					fsCalled.DisallowInlining(f.compiler, InlineFunctions_DisallowInliningReason_COST)
				}
			}
		}

		// Make a copy of the Node, so it isn't changed by other inlines.
		functionState.SetSafeFnNode(functionState.GetFn().GetFunctionNode().CloneTree())
	}
}

// FindCalledFunctions: this functions that may be called directly.
//
// Java: InlineFunctions.findCalledFunctions(Node)
func (f *InlineFunctions) FindCalledFunctions(node *rhino.Node) *jutil.LinkedHashSet[string] {
	changed := jutil.NewLinkedHashSet[string]()
	inlineFunctions_findCalledFunctions(NodeUtil_GetFunctionBody(node), changed)
	return changed
}

// InlineFunctions_FindCalledFunctionsInto: see findCalledFunctions(Node).
//
// Java: InlineFunctions.findCalledFunctions(Node, Set<String>)
func InlineFunctions_FindCalledFunctionsInto(node *rhino.Node, changed *jutil.LinkedHashSet[string]) {
	jutil.CheckArgument(changed != nil)
	inlineFunctions_findCalledFunctions(node, changed)
}

func inlineFunctions_findCalledFunctions(node *rhino.Node, changed *jutil.LinkedHashSet[string]) {
	if node.IsName() && InlineFunctions_IsCandidateUsage(node) {
		changed.Add(node.GetString())
	}

	for c := node.GetFirstChild(); c != nil; c = c.GetNext() {
		if c.HasChildren() || c.IsName() {
			inlineFunctions_findCalledFunctions(c, changed)
		}
	}
}

// DecomposeExpressions: for any call-site that needs it, prepare the
// call-site for inlining by rewriting the containing expression.
//
// Java: InlineFunctions.decomposeExpressions()
func (f *InlineFunctions) DecomposeExpressions() {
	for functionState := range f.fns.Values() {
		if functionState.CanInline() && functionState.references != nil {
			for ref := range functionState.references.Values() {
				if ref.requiresDecomposition {
					f.injector.MaybePrepareCall(ref)
				}
			}
		}
	}
}

// RemoveInlinedFunctions removes inlined functions that no longer have any
// references.
//
// Java: InlineFunctions.removeInlinedFunctions()
func (f *InlineFunctions) RemoveInlinedFunctions() {
	for name, functionState := range f.fns.All() {
		if functionState.CanRemove() {
			fn := functionState.GetFn()
			jutil.CheckState(functionState.CanInline())
			jutil.CheckState(fn != nil)
			if f.poolState != nil && len(f.poolState.movedArgs) > 0 {
				inMovedArg := false
				for anc := fn.GetFunctionNode(); anc != nil; anc = anc.GetParent() {
					if _, ok := f.poolState.movedArgs[anc]; ok {
						inMovedArg = true
						break
					}
				}
				if inMovedArg {
					f.VerifyAllReferencesInlined(name, functionState)
					if _, isExpr := fn.(*InlineFunctions_FunctionExpression); !isExpr {
						f.compiler.GetChangeTracker().notifyChangeHandlers()
					}
					continue
				}
			}
			f.VerifyAllReferencesInlined(name, functionState)
			fn.Remove()
			NodeUtil_MarkFunctionsDeleted(fn.GetFunctionNode(), f.compiler)
		}
	}
}

// VerifyAllReferencesInlined checks to verify that expression rewriting
// didn't make a call inaccessible.
//
// Java: InlineFunctions.verifyAllReferencesInlined(String, FunctionState)
func (f *InlineFunctions) VerifyAllReferencesInlined(name string, functionState *InlineFunctions_FunctionState) {
	if functionState.references == nil {
		return
	}
	for ref := range functionState.references.Values() {
		if !ref.inlined {
			parent := ref.callNode.GetParent()
			parentString := "null"
			if parent != nil {
				parentString = parent.ToStringTree()
			}
			panic(jutil.NewIllegalStateExceptionMsg(
				"Call site missed (" +
					name +
					").\n call: " +
					ref.callNode.ToStringTree() +
					"\n parent:  " +
					parentString))
		}
	}
}

// InlineFunctions_FunctionState is used to track the decisions that have
// been made about a function.
type InlineFunctions_FunctionState struct {
	fn               InlineFunctions_Function
	safeFnNode       *rhino.Node
	inline           bool
	remove           bool
	inlineDirectly   bool
	referencesThis   bool
	hasInnerClosures bool
	references       *jutil.LinkedHashMap[*rhino.Node, *InlineFunctions_Reference]
	refMapVal        jutil.LinkedHashMap[*rhino.Node, *InlineFunctions_Reference]
	refMapBuf        [2]jutil.LhmSlot[*rhino.Node, *InlineFunctions_Reference]
	chunk            *JSChunk
	namesToAlias     *jutil.LinkedHashSet[string]

	hasRequireInliningAnnotation            bool
	hasRequireInliningAnnotationInitialized bool

	hasEncourageInliningAnnotation            bool
	hasEncourageInliningAnnotationInitialized bool
	bodyAstSize                               int32
}

// Java: InlineFunctions.FunctionState()
func NewInlineFunctions_FunctionState() *InlineFunctions_FunctionState {
	return &InlineFunctions_FunctionState{
		fn:             nil,
		safeFnNode:     nil,
		inline:         true,
		remove:         true,
		inlineDirectly: false,
		referencesThis: false,
	}
}

// Java: InlineFunctions.FunctionState.hasExistingFunctionDefinition()
func (s *InlineFunctions_FunctionState) HasExistingFunctionDefinition() bool {
	return s.fn != nil
}

// Java: InlineFunctions.FunctionState.setReferencesThis(boolean)
func (s *InlineFunctions_FunctionState) SetReferencesThis(referencesThis bool) {
	s.referencesThis = referencesThis
}

// Java: InlineFunctions.FunctionState.getReferencesThis()
func (s *InlineFunctions_FunctionState) GetReferencesThis() bool {
	return s.referencesThis
}

// SetHasInnerClosures: whether the function contains inner functions or
// non-static class field initializers.
//
// Java: InlineFunctions.FunctionState.setHasInnerClosures(boolean)
func (s *InlineFunctions_FunctionState) SetHasInnerClosures(hasInnerClosures bool) {
	s.hasInnerClosures = hasInnerClosures
}

// Java: InlineFunctions.FunctionState.hasInnerClosures()
func (s *InlineFunctions_FunctionState) HasInnerClosures() bool {
	return s.hasInnerClosures
}

// Java: InlineFunctions.FunctionState.removeBlockInliningReferences()
func (s *InlineFunctions_FunctionState) RemoveBlockInliningReferences() {
	if s.references == nil {
		return
	}
	for key, value := range s.references.All() {
		if value.mode == FunctionInjector_InliningMode_BLOCK {
			s.references.Remove(key)
		}
	}
}

// Java: InlineFunctions.FunctionState.hasBlockInliningReferences()
func (s *InlineFunctions_FunctionState) HasBlockInliningReferences() bool {
	if s.references == nil {
		return false
	}
	for r := range s.references.Values() {
		if r.mode == FunctionInjector_InliningMode_BLOCK {
			return true
		}
	}
	return false
}

// Java: InlineFunctions.FunctionState.getFn()
func (s *InlineFunctions_FunctionState) GetFn() InlineFunctions_Function {
	return s.fn
}

// Java: InlineFunctions.FunctionState.setFn(Function)
func (s *InlineFunctions_FunctionState) SetFn(fn InlineFunctions_Function) {
	jutil.CheckState(s.fn == nil)
	s.fn = fn
}

// Java: InlineFunctions.FunctionState.getSafeFnNode()
func (s *InlineFunctions_FunctionState) GetSafeFnNode() *rhino.Node {
	if s.safeFnNode != nil {
		return s.safeFnNode
	}
	return s.fn.GetFunctionNode()
}

// Java: InlineFunctions.FunctionState.setSafeFnNode(Node)
func (s *InlineFunctions_FunctionState) SetSafeFnNode(safeFnNode *rhino.Node) {
	s.safeFnNode = safeFnNode
}

// Java: InlineFunctions.FunctionState.canInline()
func (s *InlineFunctions_FunctionState) CanInline() bool {
	return s.inline
}

// Java: InlineFunctions.FunctionState.encourageInlining()
func (s *InlineFunctions_FunctionState) EncourageInlining() bool {
	if s.fn == nil {
		return false
	}
	if !s.hasEncourageInliningAnnotationInitialized {
		s.hasEncourageInliningAnnotationInitialized = true
		s.hasEncourageInliningAnnotation = InlineFunctions_HasEncourageInliningAnnotation(s.fn.GetFunctionNode())
	}
	return s.hasEncourageInliningAnnotation
}

// Java: InlineFunctions.FunctionState.requireInlining()
func (s *InlineFunctions_FunctionState) RequireInlining() bool {
	if s.fn == nil {
		return false
	}
	if !s.hasRequireInliningAnnotationInitialized {
		s.hasRequireInliningAnnotationInitialized = true
		s.hasRequireInliningAnnotation = InlineFunctions_HasRequireInliningAnnotation(s.fn.GetFunctionNode())
	}
	return s.hasRequireInliningAnnotation
}

// Java: InlineFunctions.FunctionState.disallowInlining(AbstractCompiler, DisallowInliningReason)
func (s *InlineFunctions_FunctionState) DisallowInlining(compiler AbstractCompiler, reason InlineFunctions_DisallowInliningReason) {
	s.DisallowInliningWithWarning(compiler, reason, InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_YES)
}

// Java: InlineFunctions.FunctionState.disallowInlining(AbstractCompiler, DisallowInliningReason, ShouldWarnWhenRequireInliningCannotInline)
func (s *InlineFunctions_FunctionState) DisallowInliningWithWarning(
	compiler AbstractCompiler,
	reason InlineFunctions_DisallowInliningReason,
	shouldWarnWhenRequireInliningCannotInline InlineFunctions_ShouldWarnWhenRequireInliningCannotInline) {
	if s.RequireInlining() &&
		compiler.IsDebugLoggingEnabled() &&
		shouldWarnWhenRequireInliningCannotInline ==
			InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_YES {
		compiler.Report(
			JSError_MakeWithNode(s.fn.GetFunctionNode(), InlineFunctions_MISSED_REQUIRED_INLINING, reason.String()))
	}
	s.inline = false

	// No need to keep references to function that can't be inlined.
	s.references = nil
	// Don't remove functions that we aren't inlining.
	s.remove = false
}

// Java: InlineFunctions.FunctionState.canRemove()
func (s *InlineFunctions_FunctionState) CanRemove() bool {
	return s.remove
}

// Java: InlineFunctions.FunctionState.setRemove(boolean, AbstractCompiler, CannotRemoveReason, ShouldWarnWhenRequireInliningCannotInline)
func (s *InlineFunctions_FunctionState) SetRemove(
	remove bool,
	compiler AbstractCompiler,
	reason InlineFunctions_CannotRemoveReason,
	shouldWarnIfRequireInlining InlineFunctions_ShouldWarnWhenRequireInliningCannotInline) {
	s.SetRemoveWithContext(remove, compiler, reason, shouldWarnIfRequireInlining /* contextNode= */, nil)
}

// Java: InlineFunctions.FunctionState.setRemove(boolean, AbstractCompiler, CannotRemoveReason, ShouldWarnWhenRequireInliningCannotInline, Node)
func (s *InlineFunctions_FunctionState) SetRemoveWithContext(
	remove bool,
	compiler AbstractCompiler,
	reason InlineFunctions_CannotRemoveReason,
	shouldWarnIfRequireInlining InlineFunctions_ShouldWarnWhenRequireInliningCannotInline,
	contextNode *rhino.Node) {
	if !remove &&
		(shouldWarnIfRequireInlining != InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_NO) &&
		s.RequireInlining() &&
		compiler.IsDebugLoggingEnabled() {
		contextString := "unknown"
		if contextNode != nil {
			contextString = contextNode.String()
		}
		compiler.Report(
			JSError_MakeWithNode(
				s.fn.GetFunctionNode(),
				InlineFunctions_MISSED_REQUIRED_INLINING,
				reason.String(),
				contextString))
	}
	s.remove = remove
}

// Java: InlineFunctions.FunctionState.canInlineDirectly()
func (s *InlineFunctions_FunctionState) CanInlineDirectly() bool {
	return s.inlineDirectly
}

// Java: InlineFunctions.FunctionState.inlineDirectly(boolean)
func (s *InlineFunctions_FunctionState) InlineDirectly(directReplacement bool) {
	s.inlineDirectly = directReplacement
}

// Java: InlineFunctions.FunctionState.hasReferences()
func (s *InlineFunctions_FunctionState) HasReferences() bool {
	return s.references != nil && !s.references.IsEmpty()
}

// Java: InlineFunctions.FunctionState.getReferencesInternal()
func (s *InlineFunctions_FunctionState) GetReferencesInternal() *jutil.LinkedHashMap[*rhino.Node, *InlineFunctions_Reference] {
	if s.references == nil {
		return jutil.NewLinkedHashMap[*rhino.Node, *InlineFunctions_Reference]()
	}
	return s.references
}

// Java: InlineFunctions.FunctionState.addReference(Reference)
func (s *InlineFunctions_FunctionState) AddReference(ref *InlineFunctions_Reference) {
	if s.references == nil {
		s.refMapVal.InitWithBuffer(s.refMapBuf[:])
		s.references = &s.refMapVal
	}
	s.references.Put(ref.callNode, ref)
}

var inlineFunctions_emptyStringSet = jutil.NewLinkedHashSet[string]()

// GetReferences returns the references (the values of the references map).
//
// Go-only: language-adapter — Java's live values() view becomes a snapshot
// slice; no caller modifies the map while iterating it.
//
// Java: InlineFunctions.FunctionState.getReferences()
func (s *InlineFunctions_FunctionState) GetReferences() []*InlineFunctions_Reference {
	if s.references == nil || s.references.IsEmpty() {
		return nil
	}
	out := make([]*InlineFunctions_Reference, 0, s.references.Size())
	for ref := range s.references.Values() {
		out = append(out, ref)
	}
	return out
}

// Java: InlineFunctions.FunctionState.getReference(Node)
func (s *InlineFunctions_FunctionState) GetReference(n *rhino.Node) *InlineFunctions_Reference {
	if s.references == nil {
		return nil
	}
	return s.references.Get(n)
}

// Java: InlineFunctions.FunctionState.getNamesToAlias()
func (s *InlineFunctions_FunctionState) GetNamesToAlias() *jutil.LinkedHashSet[string] {
	if s.namesToAlias == nil {
		return inlineFunctions_emptyStringSet
	}
	return s.namesToAlias
}

// Java: InlineFunctions.FunctionState.setNamesToAlias(Set<String>)
func (s *InlineFunctions_FunctionState) SetNamesToAlias(names *jutil.LinkedHashSet[string]) {
	s.namesToAlias = names
}

// Java: InlineFunctions.FunctionState.setChunk(JSChunk)
func (s *InlineFunctions_FunctionState) SetChunk(chunk *JSChunk) {
	s.chunk = chunk
}

// Java: InlineFunctions.FunctionState.getChunk()
func (s *InlineFunctions_FunctionState) GetChunk() *JSChunk {
	return s.chunk
}

// InlineFunctions_Function is an interface for dealing with function
// declarations and function expressions equally.
type InlineFunctions_Function interface {
	// GetName gets the name of the function.
	//
	// Java: InlineFunctions.Function.getName()
	GetName() string

	// GetNameNode gets the name node of the function.
	//
	// Java: InlineFunctions.Function.getNameNode()
	GetNameNode() *rhino.Node

	// GetFunctionNode gets the function node.
	//
	// Java: InlineFunctions.Function.getFunctionNode()
	GetFunctionNode() *rhino.Node

	// Remove removes itself from the JavaScript.
	//
	// Java: InlineFunctions.Function.remove()
	Remove()

	// Java: InlineFunctions.Function.getDeclaringBlock()
	GetDeclaringBlock() *rhino.Node
}

// InlineFunctions_NamedFunction is the NamedFunction implementation of the
// Function interface.
type InlineFunctions_NamedFunction struct {
	outer *InlineFunctions

	fn *rhino.Node
}

// Java: InlineFunctions.NamedFunction(Node)
func NewInlineFunctions_NamedFunction(outer *InlineFunctions, fn *rhino.Node) *InlineFunctions_NamedFunction {
	return &InlineFunctions_NamedFunction{outer: outer, fn: fn}
}

// Java: InlineFunctions.NamedFunction.getName()
func (n *InlineFunctions_NamedFunction) GetName() string {
	return n.fn.GetFirstChild().GetString()
}

// Java: InlineFunctions.NamedFunction.getNameNode()
func (n *InlineFunctions_NamedFunction) GetNameNode() *rhino.Node {
	return n.fn.GetFirstChild()
}

// Java: InlineFunctions.NamedFunction.getFunctionNode()
func (n *InlineFunctions_NamedFunction) GetFunctionNode() *rhino.Node {
	return n.fn
}

// Java: InlineFunctions.NamedFunction.remove()
func (n *InlineFunctions_NamedFunction) Remove() {
	n.outer.compiler.ReportChangeToEnclosingScope(n.fn)
	NodeUtil_RemoveChild(n.fn.GetParent(), n.fn)
	NodeUtil_MarkFunctionsDeleted(n.fn, n.outer.compiler)
}

// Java: InlineFunctions.NamedFunction.getDeclaringBlock()
func (n *InlineFunctions_NamedFunction) GetDeclaringBlock() *rhino.Node {
	return n.fn.GetParent()
}

// InlineFunctions_FunctionVar is the FunctionVar implementation of the
// Function interface.
type InlineFunctions_FunctionVar struct {
	outer *InlineFunctions

	var_ *rhino.Node
}

// Java: InlineFunctions.FunctionVar(Node)
func NewInlineFunctions_FunctionVar(outer *InlineFunctions, var_ *rhino.Node) *InlineFunctions_FunctionVar {
	return &InlineFunctions_FunctionVar{outer: outer, var_: var_}
}

// Java: InlineFunctions.FunctionVar.getName()
func (v *InlineFunctions_FunctionVar) GetName() string {
	return v.var_.GetFirstChild().GetString()
}

// Java: InlineFunctions.FunctionVar.getNameNode()
func (v *InlineFunctions_FunctionVar) GetNameNode() *rhino.Node {
	return v.var_.GetFirstChild()
}

// Java: InlineFunctions.FunctionVar.getFunctionNode()
func (v *InlineFunctions_FunctionVar) GetFunctionNode() *rhino.Node {
	return v.var_.GetFirstFirstChild()
}

// Java: InlineFunctions.FunctionVar.remove()
func (v *InlineFunctions_FunctionVar) Remove() {
	v.outer.compiler.ReportChangeToEnclosingScope(v.var_)
	NodeUtil_RemoveChild(v.var_.GetParent(), v.var_)
	NodeUtil_MarkFunctionsDeleted(v.var_, v.outer.compiler)
}

// Java: InlineFunctions.FunctionVar.getDeclaringBlock()
func (v *InlineFunctions_FunctionVar) GetDeclaringBlock() *rhino.Node {
	return v.var_.GetParent()
}

// InlineFunctions_FunctionExpression is the FunctionExpression
// implementation of the Function interface.
type InlineFunctions_FunctionExpression struct {
	fn           *rhino.Node
	fakeName     string
	fakeNameNode *rhino.Node
}

// Java: InlineFunctions.FunctionExpression(Node, int)
func NewInlineFunctions_FunctionExpression(fn *rhino.Node, index int) *InlineFunctions_FunctionExpression {
	return &InlineFunctions_FunctionExpression{
		fn:       fn,
		fakeName: strconv.Itoa(index),
	}
}

// Java: InlineFunctions.FunctionExpression.getName()
func (e *InlineFunctions_FunctionExpression) GetName() string {
	return e.fakeName
}

// Java: InlineFunctions.FunctionExpression.getNameNode()
func (e *InlineFunctions_FunctionExpression) GetNameNode() *rhino.Node {
	if e.fakeNameNode == nil {
		e.fakeNameNode = rhino.IR_Name(e.fakeName)
	}
	return e.fakeNameNode
}

// Java: InlineFunctions.FunctionExpression.getFunctionNode()
func (e *InlineFunctions_FunctionExpression) GetFunctionNode() *rhino.Node {
	return e.fn
}

// Java: InlineFunctions.FunctionExpression.remove()
func (e *InlineFunctions_FunctionExpression) Remove() {
	// Nothing to do. The function is removed with the call.
}

// Java: InlineFunctions.FunctionExpression.getDeclaringBlock()
func (e *InlineFunctions_FunctionExpression) GetDeclaringBlock() *rhino.Node {
	return nil
}

// InlineFunctions_Reference is InlineFunctions.Reference, a subclass of
// FunctionInjector.Reference.
type InlineFunctions_Reference struct {
	FunctionInjector_ReferenceBase

	requiresDecomposition bool
	inlined               bool
}

// Java: InlineFunctions.Reference(Node, Scope, JSChunk, InliningMode)
func NewInlineFunctions_Reference(callNode *rhino.Node, scope *Scope, chunk *JSChunk, mode FunctionInjector_InliningMode) *InlineFunctions_Reference {
	r := &InlineFunctions_Reference{requiresDecomposition: false, inlined: false}
	r.InitFunctionInjector_ReferenceBase(callNode, scope, chunk, mode)
	return r
}

// Java: InlineFunctions.Reference.setRequiresDecomposition(boolean)
func (r *InlineFunctions_Reference) SetRequiresDecomposition(newVal bool) {
	r.requiresDecomposition = newVal
}

// InlineFunctions_ShouldWarnWhenRequireInliningCannotInline is the enum
// InlineFunctions.ShouldWarnWhenRequireInliningCannotInline.
type InlineFunctions_ShouldWarnWhenRequireInliningCannotInline int32

const (
	InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_YES InlineFunctions_ShouldWarnWhenRequireInliningCannotInline = iota
	InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_NO
)

var inlineFunctions_ShouldWarnWhenRequireInliningCannotInline_names = [...]string{"YES", "NO"}

// Java: InlineFunctions.ShouldWarnWhenRequireInliningCannotInline.values()
func InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_values() []InlineFunctions_ShouldWarnWhenRequireInliningCannotInline {
	return []InlineFunctions_ShouldWarnWhenRequireInliningCannotInline{
		InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_YES,
		InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_NO,
	}
}

// Java: InlineFunctions.ShouldWarnWhenRequireInliningCannotInline.name()
func (e InlineFunctions_ShouldWarnWhenRequireInliningCannotInline) String() string {
	return inlineFunctions_ShouldWarnWhenRequireInliningCannotInline_names[e]
}

// Java: InlineFunctions.ShouldWarnWhenRequireInliningCannotInline.ordinal()
func (e InlineFunctions_ShouldWarnWhenRequireInliningCannotInline) Ordinal() int { return int(e) }

// InlineFunctions_CannotRemoveReason is the enum
// InlineFunctions.CannotRemoveReason.
type InlineFunctions_CannotRemoveReason int32

const (
	InlineFunctions_CannotRemoveReason_BLOCK_INLINING_REFERENCE InlineFunctions_CannotRemoveReason = iota
	InlineFunctions_CannotRemoveReason_CALLED_BY_INLINED_FUNCTION
	InlineFunctions_CannotRemoveReason_UNINLINABLE_REFERENCE
	InlineFunctions_CannotRemoveReason_UNREMOVABLE_REFERENCE
)

var inlineFunctions_CannotRemoveReason_names = [...]string{
	"BLOCK_INLINING_REFERENCE",
	"CALLED_BY_INLINED_FUNCTION",
	"UNINLINABLE_REFERENCE",
	"UNREMOVABLE_REFERENCE",
}

// Java: InlineFunctions.CannotRemoveReason.values()
func InlineFunctions_CannotRemoveReason_values() []InlineFunctions_CannotRemoveReason {
	return []InlineFunctions_CannotRemoveReason{
		InlineFunctions_CannotRemoveReason_BLOCK_INLINING_REFERENCE,
		InlineFunctions_CannotRemoveReason_CALLED_BY_INLINED_FUNCTION,
		InlineFunctions_CannotRemoveReason_UNINLINABLE_REFERENCE,
		InlineFunctions_CannotRemoveReason_UNREMOVABLE_REFERENCE,
	}
}

// Java: InlineFunctions.CannotRemoveReason.name()
func (e InlineFunctions_CannotRemoveReason) String() string {
	return inlineFunctions_CannotRemoveReason_names[e]
}

// Java: InlineFunctions.CannotRemoveReason.ordinal()
func (e InlineFunctions_CannotRemoveReason) Ordinal() int { return int(e) }

// InlineFunctions_DisallowInliningReason is the enum
// InlineFunctions.DisallowInliningReason.
type InlineFunctions_DisallowInliningReason int32

const (
	InlineFunctions_DisallowInliningReason_ASYNC_FUNCTION InlineFunctions_DisallowInliningReason = iota
	InlineFunctions_DisallowInliningReason_CAPTURED_VARIABLE_NOT_IN_SCOPE
	InlineFunctions_DisallowInliningReason_COST
	InlineFunctions_DisallowInliningReason_GENERATOR_FUNCTION
	InlineFunctions_DisallowInliningReason_HAS_LOCAL_NAMES
	InlineFunctions_DisallowInliningReason_MULTIPLE_FUNCTION_DEFINITIONS
	InlineFunctions_DisallowInliningReason_NON_INLINABLE_PARAM
	InlineFunctions_DisallowInliningReason_NOT_CANDIDATE_FUNCTION
	InlineFunctions_DisallowInliningReason_NO_INLINE_ANNOTATION
	InlineFunctions_DisallowInliningReason_REASSIGNED
	InlineFunctions_DisallowInliningReason_TARGET_SIZE_EXCEEDS_LIMIT
)

var inlineFunctions_DisallowInliningReason_names = [...]string{
	"ASYNC_FUNCTION",
	"CAPTURED_VARIABLE_NOT_IN_SCOPE",
	"COST",
	"GENERATOR_FUNCTION",
	"HAS_LOCAL_NAMES",
	"MULTIPLE_FUNCTION_DEFINITIONS",
	"NON_INLINABLE_PARAM",
	"NOT_CANDIDATE_FUNCTION",
	"NO_INLINE_ANNOTATION",
	"REASSIGNED",
	"TARGET_SIZE_EXCEEDS_LIMIT",
}

// Java: InlineFunctions.DisallowInliningReason.values()
func InlineFunctions_DisallowInliningReason_values() []InlineFunctions_DisallowInliningReason {
	out := make([]InlineFunctions_DisallowInliningReason, len(inlineFunctions_DisallowInliningReason_names))
	for i := range out {
		out[i] = InlineFunctions_DisallowInliningReason(i)
	}
	return out
}

// Java: InlineFunctions.DisallowInliningReason.name()
func (e InlineFunctions_DisallowInliningReason) String() string {
	return inlineFunctions_DisallowInliningReason_names[e]
}

// Java: InlineFunctions.DisallowInliningReason.ordinal()
func (e InlineFunctions_DisallowInliningReason) Ordinal() int { return int(e) }

// InlineFunctions_ShouldRequireInlining is the record
// InlineFunctions.ShouldRequireInlining(boolean shouldValidate, String
// optMessage).
type InlineFunctions_ShouldRequireInlining struct {
	shouldValidate bool
	optMessage     string
}

// Java: InlineFunctions.ShouldRequireInlining(boolean, String)
func NewInlineFunctions_ShouldRequireInlining(shouldValidate bool, optMessage string) *InlineFunctions_ShouldRequireInlining {
	return &InlineFunctions_ShouldRequireInlining{shouldValidate: shouldValidate, optMessage: optMessage}
}

// Java: InlineFunctions.ShouldRequireInlining.shouldValidate()
func (r *InlineFunctions_ShouldRequireInlining) ShouldValidate() bool { return r.shouldValidate }

// Java: InlineFunctions.ShouldRequireInlining.optMessage()
func (r *InlineFunctions_ShouldRequireInlining) OptMessage() string { return r.optMessage }

// Java: InlineFunctions.ShouldRequireInlining.equals(Object)
func (r *InlineFunctions_ShouldRequireInlining) Equals(o any) bool {
	other, ok := o.(*InlineFunctions_ShouldRequireInlining)
	return ok && other != nil && r.shouldValidate == other.shouldValidate && r.optMessage == other.optMessage
}

// Java: InlineFunctions.ShouldRequireInlining.hashCode()
func (r *InlineFunctions_ShouldRequireInlining) HashCode() int32 {
	return 31*jutil.ObjectHashCode(r.shouldValidate) + jutil.StringHashCode(r.optMessage)
}

// Java: InlineFunctions.ShouldRequireInlining.toString()
func (r *InlineFunctions_ShouldRequireInlining) String() string {
	return "ShouldRequireInlining[shouldValidate=" + strconv.FormatBool(r.shouldValidate) + ", optMessage=" + r.optMessage + "]"
}

// Java: InlineFunctions.checkOptimizationLevelSupportsRequireInlining(CompilerOptions)
func InlineFunctions_CheckOptimizationLevelSupportsRequireInlining(
	options *CompilerOptions) *InlineFunctions_ShouldRequireInlining {
	// Non-optimized code might not perform all required inlinings.
	var problems []string
	if !options.ShouldOptimize() {
		problems = append(problems, "non-optimized code")
	}
	// Required inlinings don't make sense without inlining. If we don't inline
	// properties we might not be able to inline functions defined on namespaces.
	if !options.ShouldInlineProperties() {
		problems = append(problems, "missing shouldInlineProperties")
	}
	if !options.ShouldInlineConstantVars() {
		problems = append(problems, "missing inlineConstantVars")
	}
	if !options.GetSmartNameRemoval() {
		problems = append(problems, "missing smartNameRemoval")
	}
	if !options.GetInlineFunctionsLevel().IsOn() {
		problems = append(problems, "missing inlineFunctions")
	}
	if !options.GetInlineFunctionsLevel().IncludesGlobals() {
		problems = append(problems, "inlineFunctions does not include globals")
	}
	// Coverage runs may break some inlining.
	if options.GetInstrumentForCoverageOption() != CompilerOptions_InstrumentOption_NONE {
		problems = append(problems, "instrumentForCoverageOption enabled")
	}
	// We need to collapse properties or functions might be retained due to
	// being written on an exported namespace.
	if options.GetPropertyCollapseLevel() != CompilerOptions_PropertyCollapseLevel_ALL {
		problems = append(problems, "missing property collapsing")
	}
	// Avoid inlining failures due to asserts.
	if !options.ShouldRemoveClosureAsserts() {
		problems = append(problems, "missing removeClosureAsserts")
	}
	// Avoid inlining failures due to being unable to devirtualize a method.
	if !options.ShouldDevirtualizeMethods() {
		problems = append(problems, "missing devirtualizeMethods")
	}
	// We need more than one inlining pass for most required inlinings to work. This
	// eliminates 'fast' mode.
	if options.GetMaxOptimizationLoopIterations() == 1 {
		problems = append(problems,
			"optimizationLoopMaxIterations should be >= 2, but is: "+
				strconv.Itoa(options.GetMaxOptimizationLoopIterations()))
	}
	if len(problems) == 0 {
		return NewInlineFunctions_ShouldRequireInlining(true, "")
	}
	return NewInlineFunctions_ShouldRequireInlining(false, "["+strings.Join(problems, ", ")+"]")
}
