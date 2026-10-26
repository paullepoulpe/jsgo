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
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// inlineFunctionsRunState is the state of one run of InlineFunctions.Process,
// dropped when the run ends.
type inlineFunctionsRunState struct {
	fns          *jutil.LinkedHashMap[string, *InlineFunctions_FunctionState]
	anonFns      *jutil.LinkedHashMap[rhino.Node, string]
	fnNames      map[string]bool
	calledNames  *jutil.LinkedHashSet[string]
	refsScratch  []FunctionInjector_Reference
	paramScratch []string
	movedArgs    map[rhino.Node]struct{}
	references   jutil.Slab[InlineFunctions_Reference] // one run's references die together
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
	anonFns *jutil.LinkedHashMap[rhino.Node, string]

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
			fns:         jutil.NewLinkedHashMap[string, *InlineFunctions_FunctionState](),
			anonFns:     jutil.NewLinkedHashMap[rhino.Node, string](),
			calledNames: jutil.NewLinkedHashSet[string](),
		}
		f.runState = p
		f.fns = p.fns
		f.anonFns = p.anonFns
	}
	return f.runState
}

func (f *InlineFunctions) allocReference(callNode rhino.Node, scope *Scope, chunk *JSChunk, mode FunctionInjector_InliningMode) *InlineFunctions_Reference {
	r := f.ensureRunState().references.New()
	r.InitFunctionInjector_ReferenceBase(callNode, scope, chunk, mode)
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
		findRefs := InlineFunctions_FindCandidatesReferences{outer: f}
		findRefs.InitInlineFunctions_CallVisitorBase(f.fns, f.anonFns, &findRefs)
		inlineFunctions_traverseWithScopeCreator(f.compiler, root, &findRefs, sc)
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
	inlineCb := InlineFunctions_Inline{injector: f.injector}
	callVisitor := InlineFunctions_CallVisitorBase{}
	callVisitor.InitInlineFunctions_CallVisitorBase(f.fns, f.anonFns, &inlineCb)
	if !root.IsNil() {
		callVisitor.inlineWalk(root, rhino.NilNode)
	}

	f.RemoveInlinedFunctions()
}

// Java: InlineFunctions.isAlwaysInlinable(Node)
func InlineFunctions_IsAlwaysInlinable(fn rhino.Node) bool {
	jutil.CheckArgument(fn.IsFunction())
	body := NodeUtil_GetFunctionBody(fn)
	return (!body.HasChildren()) || (body.HasOneChild() && body.GetFirstChild().IsReturn())
}

// Java: InlineFunctions.targetSizeAfterInlineExceedsLimit(NodeTraversal, FunctionState)
func (f *InlineFunctions) TargetSizeAfterInlineExceedsLimit(t *NodeTraversal, functionState *InlineFunctions_FunctionState) bool {
	containingFunction := t.GetEnclosingFunction()
	// Always inline at the top level,
	// unless maybeAddFunction has marked functionState as not inlinable.
	if containingFunction.IsNil() {
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
		for ch := n.GetFirstChild(); !ch.IsNil(); ch = ch.GetNext() {
			if ch.IsComputedProp() || ch.IsComputedFieldDef() {
				c.walk(ch.GetFirstChild(), ch, chunk, inLocalHoist)
			}
		}
		for ch := n.GetFirstChild(); !ch.IsNil(); ch = ch.GetNext() {
			switch ch.GetToken() {
			case rhino.Token_COMPUTED_PROP:
				c.walk(ch.GetLastChild(), ch, chunk, inLocalHoist)
			case rhino.Token_COMPUTED_FIELD_DEF:
				if ch.HasTwoChildren() {
					c.walk(ch.GetLastChild(), ch, chunk, true)
				}
			default:
				c.walk(ch, n, chunk, inLocalHoist)
			}
		}
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
		jutil.CheckStateMsg1(n.HasOneChild(), "%s", n)
		nameNode := n.GetFirstChild()
		if nameNode.IsName() &&
			nameNode.HasChildren() &&
			nameNode.GetFirstChild().IsFunction() {
			c.outer.MaybeAddFunction(NewInlineFunctions_FunctionVar(c.outer, n), chunk)
		}
		// Named functions
		// function Foo(x) { return ... }
	case rhino.Token_FUNCTION:
		if !NodeUtil_IsStatement(n) {
			return
		}
		jutil.CheckState(NodeUtil_IsStatementBlock(parent) || parent.IsLabel())
		if NodeUtil_IsFunctionDeclaration(n) {
			c.outer.MaybeAddFunction(NewInlineFunctions_NamedFunction(c.outer, n), chunk)
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
		var fnNode rhino.Node
		first := n.GetFirstChild()
		if first.IsFunction() {
			fnNode = first
		} else if first.HasChildren() && first.GetFirstChild().IsFunction() && NodeUtil_IsFunctionObjectCall(n) {
			fnNode = first.GetFirstChild()
		}

		// If an interesting function was discovered, add it.
		if !fnNode.IsNil() {
			callsSeen := c.callsSeen
			c.callsSeen++
			fn := NewInlineFunctions_FunctionExpression(fnNode, callsSeen)
			c.outer.MaybeAddFunction(fn, chunk)
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
				if InlineFunctions_HasLocalNames(fnNode) {
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

func (f *InlineFunctions) blockHasLetOrConst(block rhino.Node) bool {
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

func (f *InlineFunctions) analyzeFunctionBody(fnNode rhino.Node) (*jutil.LinkedHashSet[string], bool, bool) {
	return inlineFunctions_analyzeBody(fnNode, &f.ensureRunState().paramScratch)
}

func inlineFunctions_analyzeBody(fnNode rhino.Node, paramScratch *[]string) (*jutil.LinkedHashSet[string], bool, bool) {
	params := NodeUtil_GetFunctionParameters(fnNode)
	names := (*paramScratch)[:0]
	for n := params.GetFirstChild(); !n.IsNil(); n = n.GetNext() {
		if n.IsRest() {
			names = append(names, FunctionArgumentInjector_REST_MARKER)
		} else if n.IsDefaultValue() || n.IsObjectPattern() || n.IsArrayPattern() {
			panic(jutil.NewIllegalStateExceptionMsg("Not supported: " + n.String()))
		} else {
			names = append(names, n.GetString())
		}
	}
	*paramScratch = names
	var res inlineFunctions_bodyAnalysis
	inlineFunctions_analyzeBodyNode(NodeUtil_GetFunctionBody(fnNode), names, false, true, &res)
	return res.namesToAlias, res.referencesThis, res.hasInnerClosures
}

func inlineFunctions_analyzeBodyNode(
	n rhino.Node,
	paramNames []string,
	inInnerClosure bool,
	inReceiverScope bool,
	res *inlineFunctions_bodyAnalysis,
) {
	switch n.GetToken() {
	case rhino.Token_FUNCTION:
		inInnerClosure = true
		res.hasInnerClosures = true
		if !n.IsArrowFunction() {
			inReceiverScope = false
		}
	case rhino.Token_MEMBER_FIELD_DEF, rhino.Token_COMPUTED_FIELD_DEF:
		if !n.IsStaticMember() {
			inInnerClosure = true
			res.hasInnerClosures = true
		}
	case rhino.Token_THIS, rhino.Token_SUPER:
		if inReceiverScope {
			res.referencesThis = true
		}
	case rhino.Token_NAME:
		if len(paramNames) != 0 {
			name := n.GetString()
			if slices.Contains(paramNames, name) && (inInnerClosure || FunctionArgumentInjector_CanNameValueChange(n)) {
				if res.namesToAlias == nil {
					res.namesToAlias = jutil.NewLinkedHashSet[string]()
				}
				res.namesToAlias.Add(name)
			}
		}
	}

	if (!inReceiverScope || res.referencesThis) && res.hasInnerClosures && len(paramNames) == 0 {
		return
	}

	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		inlineFunctions_analyzeBodyNode(c, paramNames, inInnerClosure, inReceiverScope, res)
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
	VisitCallSite(t *NodeTraversal, callNode rhino.Node, functionState *InlineFunctions_FunctionState)
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

// InlineFunctions_CallVisitorBase holds CallVisitor's fields and its visit.
type InlineFunctions_CallVisitorBase struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	callback        InlineFunctions_CallVisitorCallback
	functionMap     *jutil.LinkedHashMap[string, *InlineFunctions_FunctionState]
	anonFunctionMap *jutil.LinkedHashMap[rhino.Node, string]
}

// Java: InlineFunctions.CallVisitor(Map<String, FunctionState>, Map<Node, String>, CallVisitorCallback)
func NewInlineFunctions_CallVisitor(
	fns *jutil.LinkedHashMap[string, *InlineFunctions_FunctionState],
	anonFns *jutil.LinkedHashMap[rhino.Node, string],
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
	anonFns *jutil.LinkedHashMap[rhino.Node, string],
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

func (*InlineFunctions_CallVisitorBase) noVisitLeaves() {}

func (v *InlineFunctions_CallVisitorBase) inlineWalk(n rhino.Node, parent rhino.Node) {
	if n.GetToken() == rhino.Token_CLASS_MEMBERS {
		for c := n.GetFirstChild(); !c.IsNil(); {
			next := c.GetNext()
			if c.IsComputedProp() || c.IsComputedFieldDef() {
				v.inlineWalk(c.GetFirstChild(), c)
			}
			c = next
		}
		for c := n.GetFirstChild(); !c.IsNil(); {
			next := c.GetNext()
			switch c.GetToken() {
			case rhino.Token_COMPUTED_PROP:
				v.inlineWalk(c.GetLastChild(), c)
			case rhino.Token_COMPUTED_FIELD_DEF:
				if c.HasTwoChildren() {
					v.inlineWalk(c.GetLastChild(), c)
				}
			default:
				v.inlineWalk(c, n)
			}
			c = next
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
		v.Visit(nil, n, parent)
	}
}

// Java: InlineFunctions.CallVisitor.visit(NodeTraversal, Node, Node)
func (v *InlineFunctions_CallVisitorBase) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
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
			functionState := v.functionMap.Get(name)

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
func InlineFunctions_IsCandidateUsage(name rhino.Node) bool {
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
	anonFns *jutil.LinkedHashMap[rhino.Node, string]) *InlineFunctions_FindCandidatesReferences {
	r := &InlineFunctions_FindCandidatesReferences{outer: outer}
	r.InitInlineFunctions_CallVisitorBase(fns, anonFns, nil)
	r.callback = r
	return r
}

func (*InlineFunctions_FindCandidatesReferences) onlyNameLeaves() {}

// Java: InlineFunctions.FindCandidatesReferences.visit(NodeTraversal, Node, Node)
func (r *InlineFunctions_FindCandidatesReferences) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsName() {
		r.CheckNameUsage(n, parent)
		return
	}
	r.InlineFunctions_CallVisitorBase.Visit(t, n, parent)
}

// Java: InlineFunctions.FindCandidatesReferences.visitCallSite(NodeTraversal, Node, FunctionState)
func (r *InlineFunctions_FindCandidatesReferences) VisitCallSite(t *NodeTraversal, callNode rhino.Node, functionState *InlineFunctions_FunctionState) {
	r.MaybeAddReference(t, functionState, callNode, t.GetChunk())
}

// Java: InlineFunctions.FindCandidatesReferences.maybeAddReference(NodeTraversal, FunctionState, Node, JSChunk)
func (r *InlineFunctions_FindCandidatesReferences) MaybeAddReference(
	t *NodeTraversal, functionState *InlineFunctions_FunctionState, callNode rhino.Node, chunk *JSChunk) {
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

	r.scratchRef.requiresDecomposition = false
	r.scratchRef.inlined = false
	r.scratchRef.InitFunctionInjector_ReferenceBase(callNode, nil, chunk, mode)
	r.scratchRef.t = t
	result := r.outer.injector.CanInlineReferenceToFunction(
		&r.scratchRef,
		functionState.GetFn().GetFunctionNode(),
		functionState.GetNamesToAlias(),
		functionState.GetReferencesThis(),
		functionState.HasInnerClosures())
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

// CheckNameUsage finds functions that can be inlined.
//
// Java: InlineFunctions.FindCandidatesReferences.checkNameUsage(Node, Node)
func (r *InlineFunctions_FindCandidatesReferences) CheckNameUsage(n rhino.Node, parent rhino.Node) {
	jutil.CheckStateMsg1(n.IsName(), "%s", n)

	name := n.GetString()
	if name == "" {
		return
	}

	if InlineFunctions_IsCandidateUsage(n) {
		return
	}

	// Other refs to a function name remove its candidacy for inlining
	functionState := r.functionMap.Get(name)
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
		(!parent.GetParent().IsNil() &&
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

type ifCandKind uint8

const (
	ifCandVar ifCandKind = iota + 1
	ifCandNamedFn
	ifCandAnonExpr
)

type ifPhase1Candidate struct {
	kind                ifCandKind
	node                rhino.Node
	fnNode              rhino.Node
	name                string
	chunk               *JSChunk
	disallowReason      InlineFunctions_DisallowInliningReason
	hasDisallow         bool
	disallowBeforeSetFn bool
	inlineDirectly      bool
	referencesThis      bool
	hasInnerClosures    bool
	namesToAlias        *jutil.LinkedHashSet[string]
	bodyAstSize         int32
}

type ifPhase1Worker struct {
	outer                         *InlineFunctions
	paramScratch                  []string
	lastCheckedBlock              rhino.Node
	lastCheckedBlockHasLetOrConst bool
	cands                         []ifPhase1Candidate
}

func (w *ifPhase1Worker) blockHasLetOrConst(block rhino.Node) bool {
	if block == w.lastCheckedBlock {
		return w.lastCheckedBlockHasLetOrConst
	}
	res := inlineFunctions_hasLetOrConst(block)
	w.lastCheckedBlock = block
	w.lastCheckedBlockHasLetOrConst = res
	return res
}

func (w *ifPhase1Worker) analyzeBody(fnNode rhino.Node) (*jutil.LinkedHashSet[string], bool, bool) {
	return inlineFunctions_analyzeBody(fnNode, &w.paramScratch)
}

func (w *ifPhase1Worker) buildNamedOrVarCandidate(kind ifCandKind, node rhino.Node, nameNode rhino.Node, fnNode rhino.Node, chunk *JSChunk) ifPhase1Candidate {
	f := w.outer
	name := nameNode.GetString()
	cand := ifPhase1Candidate{
		kind:   kind,
		node:   node,
		fnNode: fnNode,
		name:   name,
		chunk:  chunk,
	}
	if f.HasNoInlineAnnotation(fnNode) {
		cand.hasDisallow = true
		cand.disallowBeforeSetFn = true
		cand.disallowReason = InlineFunctions_DisallowInliningReason_NO_INLINE_ANNOTATION
		return cand
	}
	if f.enforceMaxSizeAfterInlining &&
		!InlineFunctions_IsAlwaysInlinable(fnNode) &&
		f.maxSizeAfterInlining <= NodeUtil_CountAstSizeUpToLimit(fnNode, f.maxSizeAfterInlining) {
		cand.hasDisallow = true
		cand.disallowBeforeSetFn = true
		cand.disallowReason = InlineFunctions_DisallowInliningReason_TARGET_SIZE_EXCEEDS_LIMIT
		return cand
	}
	if FunctionInjector_IsDirectCallNodeReplacementPossible(fnNode) {
		cand.inlineDirectly = true
	}
	if fnNode.IsGeneratorFunction() {
		cand.hasDisallow = true
		cand.disallowReason = InlineFunctions_DisallowInliningReason_GENERATOR_FUNCTION
		return cand
	}
	if fnNode.IsAsyncFunction() {
		cand.hasDisallow = true
		cand.disallowReason = InlineFunctions_DisallowInliningReason_ASYNC_FUNCTION
		return cand
	}
	if InlineFunctions_HasNonInlinableParam(NodeUtil_GetFunctionParameters(fnNode)) {
		cand.hasDisallow = true
		cand.disallowReason = InlineFunctions_DisallowInliningReason_NON_INLINABLE_PARAM
		return cand
	}
	conv := f.compiler.GetCodingConvention()
	if conv.IsExported(name) || conv.IsPropertyRenameFunction(nameNode) || !f.injector.DoesFunctionMeetMinimumRequirements(name, fnNode) {
		cand.hasDisallow = true
		cand.disallowReason = InlineFunctions_DisallowInliningReason_NOT_CANDIDATE_FUNCTION
		return cand
	}
	namesToAlias, referencesThis, hasInnerClosures := w.analyzeBody(fnNode)
	if namesToAlias != nil && !namesToAlias.IsEmpty() {
		cand.inlineDirectly = false
		cand.namesToAlias = namesToAlias
	}
	cand.referencesThis = referencesThis
	if hasInnerClosures {
		cand.hasInnerClosures = true
		if InlineFunctions_HasLocalNames(fnNode) {
			cand.hasDisallow = true
			cand.disallowReason = InlineFunctions_DisallowInliningReason_HAS_LOCAL_NAMES
			return cand
		}
	}
	if fnNode.GetGrandparent().IsVar() {
		var block rhino.Node
		if kind == ifCandVar {
			block = node.GetParent()
		} else {
			block = fnNode.GetParent()
		}
		if block.IsBlock() && !block.GetParent().IsFunction() && w.blockHasLetOrConst(block) {
			cand.hasDisallow = true
			cand.disallowReason = InlineFunctions_DisallowInliningReason_CAPTURED_VARIABLE_NOT_IN_SCOPE
			return cand
		}
	}
	if f.enforceMaxSizeAfterInlining && !InlineFunctions_IsAlwaysInlinable(fnNode) {
		cand.bodyAstSize = int32(NodeUtil_CountAstSizeUpToLimit(NodeUtil_GetFunctionBody(fnNode), f.maxSizeAfterInlining))
	}
	return cand
}

func (w *ifPhase1Worker) walk(n rhino.Node, parent rhino.Node, chunk *JSChunk, inLocalHoist bool) {
	tok := n.GetToken()
	childLocalHoist := inLocalHoist
	switch tok {
	case rhino.Token_SCRIPT:
		if id := n.GetInputId(); id != nil {
			if inp := w.outer.compiler.GetInput(id); inp != nil {
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
		for ch := n.GetFirstChild(); !ch.IsNil(); ch = ch.GetNext() {
			if ch.IsComputedProp() || ch.IsComputedFieldDef() {
				w.walk(ch.GetFirstChild(), ch, chunk, inLocalHoist)
			}
		}
		for ch := n.GetFirstChild(); !ch.IsNil(); ch = ch.GetNext() {
			switch ch.GetToken() {
			case rhino.Token_COMPUTED_PROP:
				w.walk(ch.GetLastChild(), ch, chunk, inLocalHoist)
			case rhino.Token_COMPUTED_FIELD_DEF:
				if ch.HasTwoChildren() {
					w.walk(ch.GetLastChild(), ch, chunk, true)
				}
			default:
				w.walk(ch, n, chunk, inLocalHoist)
			}
		}
		return
	}
	for ch := n.GetFirstChild(); !ch.IsNil(); ch = ch.GetNext() {
		if !ch.GetFirstChild().IsNil() {
			w.walk(ch, n, chunk, childLocalHoist)
		}
	}
	switch tok {
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		if (w.outer.reach.IncludesGlobals() || inLocalHoist) && NodeUtil_IsStatement(n) {
			jutil.CheckStateMsg1(n.HasOneChild(), "%s", n)
			nameNode := n.GetFirstChild()
			if nameNode.IsName() && nameNode.HasChildren() && nameNode.GetFirstChild().IsFunction() {
				w.cands = append(w.cands, w.buildNamedOrVarCandidate(ifCandVar, n, nameNode, nameNode.GetFirstChild(), chunk))
			}
		}
	case rhino.Token_FUNCTION:
		if (w.outer.reach.IncludesGlobals() || inLocalHoist) && NodeUtil_IsStatement(n) {
			jutil.CheckState(NodeUtil_IsStatementBlock(parent) || parent.IsLabel())
			if NodeUtil_IsFunctionDeclaration(n) {
				w.cands = append(w.cands, w.buildNamedOrVarCandidate(ifCandNamedFn, n, n.GetFirstChild(), n, chunk))
			}
		}
	case rhino.Token_OPTCHAIN_CALL, rhino.Token_CALL:
		if w.outer.reach.IncludesGlobals() || inLocalHoist {
			var fnNode rhino.Node
			first := n.GetFirstChild()
			if first.IsFunction() {
				fnNode = first
			} else if first.HasChildren() && first.GetFirstChild().IsFunction() && NodeUtil_IsFunctionObjectCall(n) {
				fnNode = first.GetFirstChild()
			}
			if !fnNode.IsNil() {
				w.cands = append(w.cands, ifPhase1Candidate{
					kind:   ifCandAnonExpr,
					node:   fnNode,
					fnNode: fnNode,
					chunk:  chunk,
				})
			}
		}
	}
}

func (f *InlineFunctions) findCandidateFunctionsParallel(root rhino.Node) bool {
	if f.compiler.IsDebugLoggingEnabled() {
		return false
	}
	scripts, ok := collectRootScripts(root)
	if !ok {
		return false
	}
	numWorkers := parallelThreadCount(f.compiler, len(scripts))
	if numWorkers <= 1 {
		return false
	}
	prewarmCompilerForParallelReads(f.compiler)

	candsByScript := make([][]ifPhase1Candidate, len(scripts))
	panicsByScript := make([]any, len(scripts))

	var nextIdx atomic.Uint32
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for w := 0; w < numWorkers; w++ {
		go func() {
			defer wg.Done()
			worker := ifPhase1Worker{outer: f}
			for {
				idx := int(nextIdx.Add(1) - 1)
				if idx >= len(scripts) {
					return
				}
				worker.cands = nil
				worker.lastCheckedBlock = rhino.NilNode
				func() {
					defer func() {
						if r := recover(); r != nil {
							panicsByScript[idx] = r
						}
					}()
					worker.walk(scripts[idx], root, nil, false)
				}()
				candsByScript[idx] = worker.cands
				if panicsByScript[idx] != nil {
					return
				}
			}
		}()
	}
	wg.Wait()

	callsSeen := 0
	for idx := range scripts {
		if p := panicsByScript[idx]; p != nil {
			panic(p)
		}
		for i := range candsByScript[idx] {
			cand := &candsByScript[idx][i]
			if cand.kind == ifCandAnonExpr {
				fn := NewInlineFunctions_FunctionExpression(cand.fnNode, callsSeen)
				callsSeen++
				f.MaybeAddFunction(fn, cand.chunk)
				f.anonFns.Put(cand.fnNode, fn.GetName())
				continue
			}
			fs := f.GetOrCreateFunctionState(cand.name)
			if fs.HasExistingFunctionDefinition() {
				fs.DisallowInlining(f.compiler, InlineFunctions_DisallowInliningReason_MULTIPLE_FUNCTION_DEFINITIONS)
			} else if fs.CanInline() {
				if cand.disallowBeforeSetFn {
					fs.DisallowInlining(f.compiler, cand.disallowReason)
				} else {
					var fn InlineFunctions_Function
					if cand.kind == ifCandVar {
						fn = NewInlineFunctions_FunctionVar(f, cand.node)
					} else {
						fn = NewInlineFunctions_NamedFunction(f, cand.node)
					}
					fs.SetFn(fn)
					if cand.hasDisallow {
						fs.DisallowInlining(f.compiler, cand.disallowReason)
					} else {
						if cand.inlineDirectly {
							fs.InlineDirectly(true)
						}
						fs.SetChunk(cand.chunk)
						if cand.namesToAlias != nil && !cand.namesToAlias.IsEmpty() {
							fs.InlineDirectly(false)
							fs.SetNamesToAlias(cand.namesToAlias)
						}
						if cand.referencesThis {
							fs.SetReferencesThis(true)
						}
						if cand.hasInnerClosures {
							fs.SetHasInnerClosures(true)
						}
						fs.bodyAstSize = cand.bodyAstSize
					}
				}
			}
			jutil.CheckState(!(InlineFunctions_HasRequireInliningAnnotation(cand.fnNode) && !fs.CanInline()))
		}
	}
	return true
}

type ifPhase2EventKind uint8

const (
	ifEvAddRef ifPhase2EventKind = iota + 1
	ifEvUninlinableCall
	ifEvDisallowReassigned
	ifEvUnremovableConst
	ifEvUnremovableContext
)

type ifPhase2Event struct {
	kind                  ifPhase2EventKind
	fs                    *InlineFunctions_FunctionState
	node                  rhino.Node
	scope                 *Scope
	chunk                 *JSChunk
	mode                  FunctionInjector_InliningMode
	requiresDecomposition bool
	safeIdsConsumed       uint32
}

type ifPhase2RefCallback struct {
	InlineFunctions_CallVisitorBase

	outer                *InlineFunctions
	injector             *FunctionInjector
	safeIdsCounter       int
	scratchRef           InlineFunctions_Reference
	lastContainingFn     rhino.Node
	lastContainingFnSize int
	localDisallowed      map[*InlineFunctions_FunctionState]struct{}
	events               []ifPhase2Event
}

func (*ifPhase2RefCallback) onlyNameLeaves() {}

func (r *ifPhase2RefCallback) isDisallowed(fs *InlineFunctions_FunctionState) bool {
	if !fs.CanInline() {
		return true
	}
	if r.localDisallowed != nil {
		_, ok := r.localDisallowed[fs]
		return ok
	}
	return false
}

func (r *ifPhase2RefCallback) targetSizeExceedsLimit(t *NodeTraversal, fs *InlineFunctions_FunctionState) bool {
	containingFunction := t.GetEnclosingFunction()
	if containingFunction.IsNil() {
		return false
	}
	inlinedFun := fs.GetFn().GetFunctionNode()
	if InlineFunctions_IsAlwaysInlinable(inlinedFun) {
		return false
	}
	inlinedFunSize := int(fs.bodyAstSize)
	if inlinedFunSize == 0 {
		inlinedFunSize = NodeUtil_CountAstSizeUpToLimit(NodeUtil_GetFunctionBody(inlinedFun), r.outer.maxSizeAfterInlining)
	}
	var targetFunSize int
	if containingFunction == r.lastContainingFn {
		targetFunSize = r.lastContainingFnSize
	} else {
		targetFunSize = NodeUtil_CountAstSizeUpToLimit(containingFunction, r.outer.maxSizeAfterInlining)
		r.lastContainingFn = containingFunction
		r.lastContainingFnSize = targetFunSize
	}
	return inlinedFunSize+targetFunSize > r.outer.maxSizeAfterInlining
}

func (r *ifPhase2RefCallback) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsName() {
		r.checkNameUsage(n, parent)
		return
	}
	r.InlineFunctions_CallVisitorBase.Visit(t, n, parent)
}

func (r *ifPhase2RefCallback) VisitCallSite(t *NodeTraversal, callNode rhino.Node, functionState *InlineFunctions_FunctionState) {
	if r.isDisallowed(functionState) {
		return
	}
	chunk := t.GetChunk()
	mode := FunctionInjector_InliningMode_BLOCK
	if functionState.CanInlineDirectly() {
		mode = FunctionInjector_InliningMode_DIRECT
	}
	beforeIds := r.safeIdsCounter
	added, scope, reqDecomp, usedMode := r.tryMode(t, functionState, callNode, chunk, mode)
	if !added && mode == FunctionInjector_InliningMode_DIRECT {
		added, scope, reqDecomp, usedMode = r.tryMode(t, functionState, callNode, chunk, FunctionInjector_InliningMode_BLOCK)
	}
	consumed := uint32(r.safeIdsCounter - beforeIds)
	if added {
		r.events = append(r.events, ifPhase2Event{
			kind:                  ifEvAddRef,
			fs:                    functionState,
			node:                  callNode,
			scope:                 scope,
			chunk:                 chunk,
			mode:                  usedMode,
			requiresDecomposition: reqDecomp,
			safeIdsConsumed:       consumed,
		})
	} else {
		r.events = append(r.events, ifPhase2Event{
			kind:            ifEvUninlinableCall,
			fs:              functionState,
			node:            callNode,
			safeIdsConsumed: consumed,
		})
	}
}

func (r *ifPhase2RefCallback) tryMode(
	t *NodeTraversal,
	functionState *InlineFunctions_FunctionState,
	callNode rhino.Node,
	chunk *JSChunk,
	mode FunctionInjector_InliningMode,
) (bool, *Scope, bool, FunctionInjector_InliningMode) {
	if r.outer.enforceMaxSizeAfterInlining && r.targetSizeExceedsLimit(t, functionState) {
		return false, nil, false, mode
	}
	r.scratchRef.requiresDecomposition = false
	r.scratchRef.inlined = false
	r.scratchRef.InitFunctionInjector_ReferenceBase(callNode, nil, chunk, mode)
	r.scratchRef.t = t
	result := r.injector.CanInlineReferenceToFunction(
		&r.scratchRef,
		functionState.GetFn().GetFunctionNode(),
		functionState.GetNamesToAlias(),
		functionState.GetReferencesThis(),
		functionState.HasInnerClosures())
	if result != FunctionInjector_CanInlineResult_NO {
		var scope *Scope
		if result == FunctionInjector_CanInlineResult_AFTER_PREPARATION {
			scope = t.GetScope()
		} else if mode != FunctionInjector_InliningMode_DIRECT && !NodeUtil_IsExprCall(callNode.GetParent()) {
			scope = t.resolveLookupScope(len(t.scopes) - 1)
		}
		return true, scope, result == FunctionInjector_CanInlineResult_AFTER_PREPARATION, mode
	}
	return false, nil, false, mode
}

func (r *ifPhase2RefCallback) checkNameUsage(n rhino.Node, parent rhino.Node) {
	jutil.CheckStateMsg1(n.IsName(), "%s", n)

	name := n.GetString()
	if name == "" {
		return
	}
	if InlineFunctions_IsCandidateUsage(n) {
		return
	}
	functionState := r.functionMap.Get(name)
	if functionState == nil {
		return
	}
	if NodeUtil_IsLValue(n) {
		if r.localDisallowed == nil {
			r.localDisallowed = make(map[*InlineFunctions_FunctionState]struct{}, 4)
		}
		r.localDisallowed[functionState] = struct{}{}
		r.events = append(r.events, ifPhase2Event{
			kind: ifEvDisallowReassigned,
			fs:   functionState,
		})
	} else if (parent.IsAssign() &&
		parent.GetJSDocInfo() != nil &&
		parent.GetJSDocInfo().IsConstant() &&
		parent.GetSecondChild() == n) ||
		(!parent.GetParent().IsNil() &&
			parent.GetParent().IsConst() &&
			parent.GetFirstChild() == n) {
		r.events = append(r.events, ifPhase2Event{
			kind: ifEvUnremovableConst,
			fs:   functionState,
		})
	} else {
		r.events = append(r.events, ifPhase2Event{
			kind: ifEvUnremovableContext,
			fs:   functionState,
			node: parent,
		})
	}
}

func (f *InlineFunctions) findCandidatesReferencesParallel(root rhino.Node) bool {
	scripts, ok := collectRootScripts(root)
	if !ok {
		return false
	}
	numWorkers := parallelThreadCount(f.compiler, len(scripts))
	if numWorkers <= 1 {
		return false
	}

	prewarmCompilerForParallelReads(f.compiler)
	if f.enforceMaxSizeAfterInlining {
		for fs := range f.fns.Values() {
			if fs.CanInline() && fs.bodyAstSize == 0 {
				fnNode := fs.GetFn().GetFunctionNode()
				if !InlineFunctions_IsAlwaysInlinable(fnNode) {
					fs.bodyAstSize = int32(NodeUtil_CountAstSizeUpToLimit(NodeUtil_GetFunctionBody(fnNode), f.maxSizeAfterInlining))
				}
			}
		}
	}

	lazyGlobal := newLazyGlobalScope(f.compiler, root)
	workerFnCaches := make([]map[rhino.Node]functionInjectorFnInfo, numWorkers)

	ifEventsByScript := make([][]ifPhase2Event, len(scripts))
	runParallelWork(f.compiler, numWorkers, len(scripts), func(workerID int) (func(int, *recordingCompiler), func()) {
		var workerRC *recordingCompiler
		var cb *ifPhase2RefCallback
		var sc *SyntacticScopeCreator

		cleanup := func() {
			if cb != nil && cb.injector != nil && len(cb.injector.fnInfoCache) != 0 {
				workerFnCaches[workerID] = cb.injector.fnInfoCache
			}
		}
		runItem := func(idx int, rc *recordingCompiler) {
			if cb == nil || workerRC != rc {
				workerRC = rc
				cb = &ifPhase2RefCallback{outer: f}
				cb.InitInlineFunctions_CallVisitorBase(f.fns, f.anonFns, cb)
				localSafeSupplier := func() string {
					id := cb.safeIdsCounter
					cb.safeIdsCounter++
					return strconv.Itoa(id)
				}
				injBuilder := NewFunctionInjector_Builder(rc).
					SafeNameIdSupplier(localSafeSupplier).
					AllowDecomposition(f.injector.allowDecomposition).
					AssumeStrictThis(f.injector.assumeStrictThis).
					FunctionArgumentInjector(NewFunctionArgumentInjector(f.compiler.GetAstAnalyzer()))
				cb.injector = injBuilder.Build()
				cb.injector.SetKnownConstantFunctions(f.injector.knownConstantFunctions)
				sc = NewSyntacticScopeCreator(rc)
			} else {
				cb.events = nil
				cb.lastContainingFn = rhino.NilNode
				if cb.localDisallowed != nil {
					clear(cb.localDisallowed)
				}
			}
			b := NodeTraversal_Builder{
				compiler:     rc,
				callback:     cb,
				scopeCreator: sc,
			}
			b.traverseScriptUnderRoot(root, root, scripts[idx], lazyGlobal)
			ifEventsByScript[idx] = cb.events
		}
		return runItem, cleanup
	}, false)

	for idx := range ifEventsByScript {
		events := ifEventsByScript[idx]
		for i := range events {
			ev := &events[i]
			switch ev.kind {
			case ifEvAddRef:
				if !ev.fs.CanInline() {
					continue
				}
				for k := uint32(0); k < ev.safeIdsConsumed; k++ {
					_ = f.injector.safeNameIdSupplier()
				}
				candidate := f.allocReference(ev.node, ev.scope, ev.chunk, ev.mode)
				candidate.SetRequiresDecomposition(ev.requiresDecomposition)
				ev.fs.AddReference(candidate)
			case ifEvUninlinableCall:
				if !ev.fs.CanInline() {
					continue
				}
				for k := uint32(0); k < ev.safeIdsConsumed; k++ {
					_ = f.injector.safeNameIdSupplier()
				}
				ev.fs.SetRemoveWithContext(
					false,
					f.compiler,
					InlineFunctions_CannotRemoveReason_UNINLINABLE_REFERENCE,
					InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_NO,
					ev.node)
			case ifEvDisallowReassigned:
				ev.fs.DisallowInlining(f.compiler, InlineFunctions_DisallowInliningReason_REASSIGNED)
			case ifEvUnremovableConst:
				ev.fs.SetRemove(
					false,
					f.compiler,
					InlineFunctions_CannotRemoveReason_UNREMOVABLE_REFERENCE,
					InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_NO)
			case ifEvUnremovableContext:
				ev.fs.SetRemoveWithContext(
					false,
					f.compiler,
					InlineFunctions_CannotRemoveReason_UNREMOVABLE_REFERENCE,
					InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_YES,
					ev.node)
			}
		}
	}

	for _, cache := range workerFnCaches {
		if len(cache) == 0 {
			continue
		}
		if f.injector.fnInfoCache == nil {
			f.injector.fnInfoCache = make(map[rhino.Node]functionInjectorFnInfo, len(cache)*2)
		}
		// order-independent: merges per-function facts, which every worker computes identically from fnNode
		for fnNode, info := range cache {
			if existing, ok := f.injector.fnInfoCache[fnNode]; ok {
				if existing.inner.IsNil() && !info.inner.IsNil() {
					existing.inner = info.inner
				}
				if existing.evalState == 0 && info.evalState != 0 {
					existing.evalState = info.evalState
				}
				f.injector.fnInfoCache[fnNode] = existing
			} else {
				f.injector.fnInfoCache[fnNode] = info
			}
		}
	}

	return true
}

// Java: InlineFunctions.Inline(FunctionInjector)
func NewInlineFunctions_Inline(injector *FunctionInjector) *InlineFunctions_Inline {
	return &InlineFunctions_Inline{injector: injector}
}

func (*InlineFunctions_Inline) noVisitLeaves() {}

// Java: InlineFunctions.Inline.visitCallSite(NodeTraversal, Node, FunctionState)
func (i *InlineFunctions_Inline) VisitCallSite(t *NodeTraversal, callNode rhino.Node, functionState *InlineFunctions_FunctionState) {
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
		if t != nil {
			t.GetCompiler().ReportChangeToEnclosingScope(newExpr)
		} else {
			i.injector.compiler.ReportChangeToEnclosingScope(newExpr)
		}
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
	p := f.ensureRunState()
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
	if f.fns == nil || f.fns.IsEmpty() {
		return
	}
	cands := make([]*InlineFunctions_FunctionState, 0, f.fns.Size())
	for functionState := range f.fns.Values() {
		if functionState.HasReferences() && functionState.CanInline() {
			cands = append(cands, functionState)
		}
	}
	numWorkers := parallelThreadCount(f.compiler, len(cands))
	if numWorkers <= 1 {
		for _, functionState := range cands {
			f.ResolveInlineConflictsForFunction(functionState)
		}
		return
	}
	f.resolveInlineConflictsParallel(cands, numWorkers)
}

type ifConflictScanResult struct {
	calledStates         []*InlineFunctions_FunctionState
	anonStates           []*InlineFunctions_FunctionState
	hasCalledNames       bool
	structuralNeedsClone bool
}

func inlineFunctions_collectAnonInlinableStates(
	node rhino.Node,
	anonFns *jutil.LinkedHashMap[rhino.Node, string],
	fns *jutil.LinkedHashMap[string, *InlineFunctions_FunctionState],
	out *[]*InlineFunctions_FunctionState,
) {
	if node.IsFunction() {
		if anonName, ok := anonFns.Lookup(node); ok {
			if fs := fns.Get(anonName); fs != nil && fs.CanInline() {
				*out = append(*out, fs)
			}
		}
	}
	for c := node.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if !c.GetFirstChild().IsNil() {
			inlineFunctions_collectAnonInlinableStates(c, anonFns, fns, out)
		}
	}
}

func (f *InlineFunctions) resolveInlineConflictsParallel(
	cands []*InlineFunctions_FunctionState,
	numWorkers int,
) {
	hasAnon := f.anonFns != nil && !f.anonFns.IsEmpty()
	// Ensure LinkedHashMap lookup tables are materialized before concurrent reads.
	_ = f.fns.Get("")
	if hasAnon {
		_, _ = f.anonFns.Lookup(rhino.NilNode)
	}

	scanResults := make([]ifConflictScanResult, len(cands))
	var nextIdx atomic.Uint32
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for w := 0; w < numWorkers; w++ {
		go func() {
			defer wg.Done()
			names := jutil.NewLinkedHashSet[string]()
			for {
				idx := int(nextIdx.Add(1) - 1)
				if idx >= len(cands) {
					return
				}
				fs := cands[idx]
				fnNode := fs.GetFn().GetFunctionNode()
				body := NodeUtil_GetFunctionBody(fnNode)
				names.Clear()
				inlineFunctions_findCalledFunctions(body, names)
				if names.IsEmpty() {
					continue
				}
				var res ifConflictScanResult
				res.hasCalledNames = true
				if hasAnon {
					inlineFunctions_collectAnonInlinableStates(body, f.anonFns, f.fns, &res.anonStates)
				}
				structural := !inlineFunctions_isIsolatedStmtFn(fnNode)
				if !structural {
					for param := NodeUtil_GetFunctionParameters(fnNode).GetFirstChild(); !param.IsNil(); param = param.GetNext() {
						if !param.GetFirstChild().IsNil() {
							structural = true
							break
						}
					}
				}
				res.structuralNeedsClone = structural
				for name := range names.All() {
					if fsCalled := f.fns.Get(name); fsCalled != nil {
						res.calledStates = append(res.calledStates, fsCalled)
					}
				}
				scanResults[idx] = res
			}
		}()
	}
	wg.Wait()

	var toClone []*InlineFunctions_FunctionState
	for i, functionState := range cands {
		if !functionState.HasReferences() || !functionState.CanInline() {
			continue
		}
		res := &scanResults[i]
		if !res.hasCalledNames {
			continue
		}
		needsSafeClone := res.structuralNeedsClone
		if !needsSafeClone {
			for _, afs := range res.anonStates {
				if afs.CanInline() {
					needsSafeClone = true
					break
				}
			}
		}
		for _, fsCalled := range res.calledStates {
			if fsCalled.CanRemove() {
				fsCalled.SetRemove(
					false,
					f.compiler,
					InlineFunctions_CannotRemoveReason_CALLED_BY_INLINED_FUNCTION,
					InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_NO)
				if !f.MinimizeCost(fsCalled) &&
					!fsCalled.RequireInlining() &&
					!fsCalled.EncourageInlining() {
					fsCalled.DisallowInlining(f.compiler, InlineFunctions_DisallowInliningReason_COST)
				}
			}
			if fsCalled.CanInline() {
				needsSafeClone = true
			}
		}
		if needsSafeClone {
			toClone = append(toClone, functionState)
		}
	}

	cloneWorkers := parallelThreadCount(f.compiler, len(toClone))
	if cloneWorkers <= 1 {
		for _, fs := range toClone {
			fs.SetSafeFnNode(fs.GetFn().GetFunctionNode().CloneTree())
		}
		return
	}
	var nextClone atomic.Uint32
	wg.Add(cloneWorkers)
	for w := 0; w < cloneWorkers; w++ {
		go func() {
			defer wg.Done()
			for {
				idx := int(nextClone.Add(1) - 1)
				if idx >= len(toClone) {
					return
				}
				fs := toClone[idx]
				fs.SetSafeFnNode(fs.GetFn().GetFunctionNode().CloneTree())
			}
		}()
	}
	wg.Wait()
}

func inlineFunctions_hasNonInlinableParamNode(n rhino.Node) bool {
	if n.IsDefaultValue() || n.IsDestructuringPattern() {
		return true
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
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
func InlineFunctions_HasNonInlinableParam(node rhino.Node) bool {
	rhino.CheckNotNull(node)
	for c := node.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
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
	p := f.ensureRunState()
	names := p.calledNames
	names.Clear()
	hasAnonInlinable := inlineFunctions_findCalledFunctionsAndAnon(
		NodeUtil_GetFunctionBody(fnNode), names, f.anonFns, f.fns)
	if !names.IsEmpty() {
		needsSafeClone := hasAnonInlinable || !inlineFunctions_isIsolatedStmtFn(fnNode)
		if !needsSafeClone {
			for param := NodeUtil_GetFunctionParameters(fnNode).GetFirstChild(); !param.IsNil(); param = param.GetNext() {
				if !param.GetFirstChild().IsNil() {
					needsSafeClone = true
					break
				}
			}
		}
		// Prevent the removal of the referenced functions unless they will be immediately inlined.
		for name := range names.All() {
			fsCalled := f.fns.Get(name)
			if fsCalled != nil {
				if fsCalled.CanRemove() {
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
				if fsCalled.CanInline() {
					needsSafeClone = true
				}
			}
		}

		if needsSafeClone {
			// Make a copy of the Node, so it isn't changed by other inlines.
			functionState.SetSafeFnNode(fnNode.CloneTree())
		}
	}
}

func inlineFunctions_isIsolatedStmtFn(fnNode rhino.Node) bool {
	for anc := fnNode.GetParent(); !anc.IsNil(); anc = anc.GetParent() {
		switch anc.GetToken() {
		case rhino.Token_SCRIPT, rhino.Token_MODULE_BODY:
			return true
		case rhino.Token_BLOCK, rhino.Token_FUNCTION,
			rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST,
			rhino.Token_NAME, rhino.Token_EXPR_RESULT, rhino.Token_LABEL:
		case rhino.Token_ASSIGN:
			if p := anc.GetParent(); p.IsNil() || !p.IsExprResult() {
				return false
			}
		default:
			return false
		}
	}
	return false
}

// FindCalledFunctions: this functions that may be called directly.
//
// Java: InlineFunctions.findCalledFunctions(Node)
func (f *InlineFunctions) FindCalledFunctions(node rhino.Node) *jutil.LinkedHashSet[string] {
	changed := jutil.NewLinkedHashSet[string]()
	inlineFunctions_findCalledFunctions(NodeUtil_GetFunctionBody(node), changed)
	return changed
}

// InlineFunctions_FindCalledFunctionsInto: see findCalledFunctions(Node).
//
// Java: InlineFunctions.findCalledFunctions(Node, Set<String>)
func InlineFunctions_FindCalledFunctionsInto(node rhino.Node, changed *jutil.LinkedHashSet[string]) {
	jutil.CheckArgument(changed != nil)
	inlineFunctions_findCalledFunctions(node, changed)
}

func inlineFunctions_findCalledFunctions(node rhino.Node, changed *jutil.LinkedHashSet[string]) {
	if node.IsName() && InlineFunctions_IsCandidateUsage(node) {
		changed.Add(node.GetString())
	}

	for c := node.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if c.GetFirstChild().IsNil() {
			if c.IsName() && InlineFunctions_IsCandidateUsage(c) {
				changed.Add(c.GetString())
			}
			continue
		}
		inlineFunctions_findCalledFunctions(c, changed)
	}
}

func inlineFunctions_findCalledFunctionsAndAnon(
	node rhino.Node,
	changed *jutil.LinkedHashSet[string],
	anonFns *jutil.LinkedHashMap[rhino.Node, string],
	fns *jutil.LinkedHashMap[string, *InlineFunctions_FunctionState],
) bool {
	if anonFns == nil || anonFns.IsEmpty() {
		inlineFunctions_findCalledFunctions(node, changed)
		return false
	}
	return inlineFunctions_findCalledFunctionsWithAnon(node, changed, anonFns, fns)
}

func inlineFunctions_findCalledFunctionsWithAnon(
	node rhino.Node,
	changed *jutil.LinkedHashSet[string],
	anonFns *jutil.LinkedHashMap[rhino.Node, string],
	fns *jutil.LinkedHashMap[string, *InlineFunctions_FunctionState],
) bool {
	hasAnon := false
	if node.IsName() && InlineFunctions_IsCandidateUsage(node) {
		changed.Add(node.GetString())
	} else if node.IsFunction() {
		if anonName, ok := anonFns.Lookup(node); ok {
			if fs := fns.Get(anonName); fs != nil && fs.CanInline() {
				hasAnon = true
			}
		}
	}

	for c := node.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if c.GetFirstChild().IsNil() {
			if c.IsName() && InlineFunctions_IsCandidateUsage(c) {
				changed.Add(c.GetString())
			}
			continue
		}
		if inlineFunctions_findCalledFunctionsWithAnon(c, changed, anonFns, fns) {
			hasAnon = true
		}
	}
	return hasAnon
}

// DecomposeExpressions: for any call-site that needs it, prepare the
// call-site for inlining by rewriting the containing expression.
//
// Java: InlineFunctions.decomposeExpressions()
func (f *InlineFunctions) DecomposeExpressions() {
	if f.fns == nil || f.fns.IsEmpty() {
		return
	}
	var workRefs []*InlineFunctions_Reference
	for functionState := range f.fns.Values() {
		if functionState.CanInline() && functionState.references != nil {
			for ref := range functionState.references.Values() {
				if ref.requiresDecomposition {
					workRefs = append(workRefs, ref)
				}
			}
		}
	}
	if len(workRefs) == 0 {
		return
	}
	if len(workRefs) < 2 || f.compiler == nil || f.compiler.GetOptions() == nil || f.compiler.GetOptions().GetNumParallelThreads() <= 1 {
		for _, ref := range workRefs {
			f.injector.MaybePrepareCall(ref)
		}
		return
	}
	if !f.decomposeExpressionsParallel(workRefs) {
		for _, ref := range workRefs {
			f.injector.MaybePrepareCall(ref)
		}
	}
}

const ifDecompPlaceholderMarker = "\x00ifd_"

type ifDecompPlaceholder struct {
	node   rhino.Node
	prefix string
	seqIdx int
}

type ifDecompScriptGroup struct {
	script     rhino.Node
	refIndices []int
	mainOnly   bool
}

func parseIfDecompPlaceholder(s string) (prefix string, refIdx int, seqIdx int, ok bool) {
	n := len(s)
	if n < 8 || s[n-1] != '\x00' {
		return "", 0, 0, false
	}
	pos := strings.Index(s, ifDecompPlaceholderMarker)
	if pos < 0 {
		return "", 0, 0, false
	}
	body := s[pos+len(ifDecompPlaceholderMarker) : n-1]
	sep := strings.IndexByte(body, '_')
	if sep <= 0 || sep+1 >= len(body) {
		return "", 0, 0, false
	}
	ri, err1 := strconv.Atoi(body[:sep])
	si, err2 := strconv.Atoi(body[sep+1:])
	if err1 != nil || err2 != nil || ri < 0 || si < 0 {
		return "", 0, 0, false
	}
	return s[:pos], ri, si, true
}

func collectIfDecompPlaceholders(
	n rhino.Node,
	refIndices []int,
	placeholdersByRef [][]ifDecompPlaceholder,
) {
	if n.IsName() {
		s := n.GetString()
		if len(s) > 0 && s[len(s)-1] == '\x00' {
			if prefix, refIdx, seqIdx, ok := parseIfDecompPlaceholder(s); ok {
				if _, belongs := slices.BinarySearch(refIndices, refIdx); belongs {
					placeholdersByRef[refIdx] = append(placeholdersByRef[refIdx], ifDecompPlaceholder{
						node:   n,
						prefix: prefix,
						seqIdx: seqIdx,
					})
				}
			}
		}
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if c.GetFirstChild().IsNil() {
			if c.IsName() {
				s := c.GetString()
				if len(s) > 0 && s[len(s)-1] == '\x00' {
					if prefix, refIdx, seqIdx, ok := parseIfDecompPlaceholder(s); ok {
						if _, belongs := slices.BinarySearch(refIndices, refIdx); belongs {
							placeholdersByRef[refIdx] = append(placeholdersByRef[refIdx], ifDecompPlaceholder{
								node:   c,
								prefix: prefix,
								seqIdx: seqIdx,
							})
						}
					}
				}
			}
			continue
		}
		collectIfDecompPlaceholders(c, refIndices, placeholdersByRef)
	}
}

func (f *InlineFunctions) decomposeExpressionsParallel(workRefs []*InlineFunctions_Reference) bool {
	groupIdxByScript := make(map[rhino.Node]int, 32)
	var groups []ifDecompScriptGroup
	var anyScope *Scope
	for refIdx, ref := range workRefs {
		script := NodeUtil_GetEnclosingScript(ref.callNode)
		s := ref.lookupScope()
		if anyScope == nil && s != nil {
			anyScope = s
		}
		refMainOnly := script.IsNil() || s == nil || s.IsGlobal()
		if !refMainOnly {
			hs := s.GetClosestHoistScope()
			if hs == nil || hs.IsGlobal() {
				refMainOnly = true
			}
		}
		gIdx, exists := groupIdxByScript[script]
		if !exists {
			gIdx = len(groups)
			groupIdxByScript[script] = gIdx
			groups = append(groups, ifDecompScriptGroup{
				script:   script,
				mainOnly: refMainOnly,
			})
		} else if refMainOnly {
			groups[gIdx].mainOnly = true
		}
		groups[gIdx].refIndices = append(groups[gIdx].refIndices, refIdx)
	}

	var parallelGroupIndices []int
	var mainGroupIndices []int
	for i := range groups {
		if groups[i].mainOnly {
			mainGroupIndices = append(mainGroupIndices, i)
		} else {
			parallelGroupIndices = append(parallelGroupIndices, i)
		}
	}
	numWorkers := parallelThreadCount(f.compiler, len(parallelGroupIndices))
	if numWorkers <= 1 {
		return false
	}

	prewarmCompilerForParallelReads(f.compiler)
	restoreReg := enableParallelTypeRegistryReads(f.compiler)
	defer restoreReg()
	if anyScope != nil {
		freezeGlobalScopeForParallelReads(anyScope.GetGlobalScope())
	}

	seqCounts := make([]int, len(workRefs))
	eventsByRef := make([][]recordedCompilerEvent, len(workRefs))
	placeholdersByRef := make([][]ifDecompPlaceholder, len(workRefs))

	runGroup := func(g *ifDecompScriptGroup, rc *recordingCompiler, workerInj *FunctionInjector, curRefIdx *int, curSeq *int) {
		anyPlaceholder := false
		for _, refIdx := range g.refIndices {
			initRecordingCompiler(rc, f.compiler)
			*curRefIdx = refIdx
			*curSeq = 0
			workerInj.MaybePrepareCall(workRefs[refIdx])
			if *curSeq > 0 {
				seqCounts[refIdx] = *curSeq
				anyPlaceholder = true
			}
			eventsByRef[refIdx] = rc.takeEvents()
		}
		if anyPlaceholder {
			scanRoot := g.script
			if scanRoot.IsNil() {
				scanRoot = f.compiler.GetJsRoot()
			}
			if !scanRoot.IsNil() {
				collectIfDecompPlaceholders(scanRoot, g.refIndices, placeholdersByRef)
			}
		}
	}

	newWorkerInjector := func(rc *recordingCompiler) (*FunctionInjector, *int, *int) {
		curRefIdx := new(int)
		curSeq := new(int)
		placeholderSupplier := func() string {
			s := *curSeq
			*curSeq = s + 1
			return ifDecompPlaceholderMarker + strconv.Itoa(*curRefIdx) + "_" + strconv.Itoa(s) + "\x00"
		}
		inj := NewFunctionInjector_Builder(rc).
			SafeNameIdSupplier(placeholderSupplier).
			AllowDecomposition(f.injector.allowDecomposition).
			AssumeStrictThis(f.injector.assumeStrictThis).
			FunctionArgumentInjector(NewFunctionArgumentInjector(f.compiler.GetAstAnalyzer())).
			Build()
		inj.SetKnownConstantFunctions(f.injector.knownConstantFunctions)
		return inj, curRefIdx, curSeq
	}

	runParallelWork(f.compiler, numWorkers, len(parallelGroupIndices), func(_ int) (func(int, *recordingCompiler), func()) {
		var boundRC *recordingCompiler
		var workerInj *FunctionInjector
		var curRefIdx *int
		var curSeq *int
		runItem := func(idx int, rc *recordingCompiler) {
			if workerInj == nil || boundRC != rc {
				boundRC = rc
				workerInj, curRefIdx, curSeq = newWorkerInjector(rc)
			}
			g := &groups[parallelGroupIndices[idx]]
			runGroup(g, rc, workerInj, curRefIdx, curSeq)
		}
		return runItem, nil
	}, false)

	if len(mainGroupIndices) > 0 {
		var mainRC recordingCompiler
		initRecordingCompiler(&mainRC, f.compiler)
		mainInj, curRefIdx, curSeq := newWorkerInjector(&mainRC)
		for _, gIdx := range mainGroupIndices {
			runGroup(&groups[gIdx], &mainRC, mainInj, curRefIdx, curSeq)
		}
	}

	var realIDs []string
	var renamedSeq []bool
	for localOrGlobalRefIdx, ref := range workRefs {
		cnt := seqCounts[localOrGlobalRefIdx]
		if cnt > 0 {
			realIDs = slices.Grow(realIDs[:0], cnt)[:cnt]
			for k := 0; k < cnt; k++ {
				realIDs[k] = f.injector.safeNameIdSupplier()
			}
			renamedSeq = slices.Grow(renamedSeq[:0], cnt)[:cnt]
			clear(renamedSeq)
			s := ref.lookupScope()
			var hs *Scope
			if s != nil {
				hs = s.GetClosestHoistScope()
			}
			for _, pn := range placeholdersByRef[localOrGlobalRefIdx] {
				if pn.seqIdx >= cnt {
					continue
				}
				oldName := pn.node.GetString()
				newName := pn.prefix + realIDs[pn.seqIdx]
				pn.node.SetString(newName)
				if !renamedSeq[pn.seqIdx] {
					renamedSeq[pn.seqIdx] = true
					if s != nil {
						s.renameOwnVar(oldName, newName)
						if hs != nil && hs != s {
							hs.renameOwnVar(oldName, newName)
						}
					}
				}
			}
		}
		if evs := eventsByRef[localOrGlobalRefIdx]; len(evs) > 0 {
			replayCompilerEvents(f.compiler, evs)
		}
	}
	return true
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
			if f.runState != nil && len(f.runState.movedArgs) > 0 {
				inMovedArg := false
				for anc := fn.GetFunctionNode(); !anc.IsNil(); anc = anc.GetParent() {
					if _, ok := f.runState.movedArgs[anc]; ok {
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
			if !parent.IsNil() {
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
	safeFnNode       rhino.Node
	inline           bool
	remove           bool
	inlineDirectly   bool
	referencesThis   bool
	hasInnerClosures bool
	references       *jutil.LinkedHashMap[rhino.Node, *InlineFunctions_Reference]
	refMapVal        jutil.LinkedHashMap[rhino.Node, *InlineFunctions_Reference]
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
		safeFnNode:     rhino.NilNode,
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
func (s *InlineFunctions_FunctionState) GetSafeFnNode() rhino.Node {
	if !s.safeFnNode.IsNil() {
		return s.safeFnNode
	}
	return s.fn.GetFunctionNode()
}

// Java: InlineFunctions.FunctionState.setSafeFnNode(Node)
func (s *InlineFunctions_FunctionState) SetSafeFnNode(safeFnNode rhino.Node) {
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
	if s.refMapVal.Size() > 0 {
		s.refMapVal.Clear()
	}
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
	s.SetRemoveWithContext(remove, compiler, reason, shouldWarnIfRequireInlining /* contextNode= */, rhino.NilNode)
}

// Java: InlineFunctions.FunctionState.setRemove(boolean, AbstractCompiler, CannotRemoveReason, ShouldWarnWhenRequireInliningCannotInline, Node)
func (s *InlineFunctions_FunctionState) SetRemoveWithContext(
	remove bool,
	compiler AbstractCompiler,
	reason InlineFunctions_CannotRemoveReason,
	shouldWarnIfRequireInlining InlineFunctions_ShouldWarnWhenRequireInliningCannotInline,
	contextNode rhino.Node) {
	if !remove &&
		(shouldWarnIfRequireInlining != InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_NO) &&
		s.RequireInlining() &&
		compiler.IsDebugLoggingEnabled() {
		contextString := "unknown"
		if !contextNode.IsNil() {
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
func (s *InlineFunctions_FunctionState) GetReferencesInternal() *jutil.LinkedHashMap[rhino.Node, *InlineFunctions_Reference] {
	if s.references == nil {
		return jutil.NewLinkedHashMap[rhino.Node, *InlineFunctions_Reference]()
	}
	return s.references
}

// Java: InlineFunctions.FunctionState.addReference(Reference)
func (s *InlineFunctions_FunctionState) AddReference(ref *InlineFunctions_Reference) {
	if s.references == nil {
		s.references = &s.refMapVal
	}
	s.references.PutNew(ref.callNode, ref)
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
func (s *InlineFunctions_FunctionState) GetReference(n rhino.Node) *InlineFunctions_Reference {
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
	GetNameNode() rhino.Node

	// GetFunctionNode gets the function node.
	//
	// Java: InlineFunctions.Function.getFunctionNode()
	GetFunctionNode() rhino.Node

	// Remove removes itself from the JavaScript.
	//
	// Java: InlineFunctions.Function.remove()
	Remove()

	// Java: InlineFunctions.Function.getDeclaringBlock()
	GetDeclaringBlock() rhino.Node
}

// InlineFunctions_NamedFunction is the NamedFunction implementation of the
// Function interface.
type InlineFunctions_NamedFunction struct {
	outer *InlineFunctions

	fn rhino.Node
}

// Java: InlineFunctions.NamedFunction(Node)
func NewInlineFunctions_NamedFunction(outer *InlineFunctions, fn rhino.Node) *InlineFunctions_NamedFunction {
	return &InlineFunctions_NamedFunction{outer: outer, fn: fn}
}

// Java: InlineFunctions.NamedFunction.getName()
func (n *InlineFunctions_NamedFunction) GetName() string {
	return n.fn.GetFirstChild().GetString()
}

// Java: InlineFunctions.NamedFunction.getNameNode()
func (n *InlineFunctions_NamedFunction) GetNameNode() rhino.Node {
	return n.fn.GetFirstChild()
}

// Java: InlineFunctions.NamedFunction.getFunctionNode()
func (n *InlineFunctions_NamedFunction) GetFunctionNode() rhino.Node {
	return n.fn
}

// Java: InlineFunctions.NamedFunction.remove()
func (n *InlineFunctions_NamedFunction) Remove() {
	n.outer.compiler.ReportChangeToEnclosingScope(n.fn)
	NodeUtil_RemoveChild(n.fn.GetParent(), n.fn)
	NodeUtil_MarkFunctionsDeleted(n.fn, n.outer.compiler)
}

// Java: InlineFunctions.NamedFunction.getDeclaringBlock()
func (n *InlineFunctions_NamedFunction) GetDeclaringBlock() rhino.Node {
	return n.fn.GetParent()
}

// InlineFunctions_FunctionVar is the FunctionVar implementation of the
// Function interface.
type InlineFunctions_FunctionVar struct {
	outer *InlineFunctions

	var_ rhino.Node
}

// Java: InlineFunctions.FunctionVar(Node)
func NewInlineFunctions_FunctionVar(outer *InlineFunctions, var_ rhino.Node) *InlineFunctions_FunctionVar {
	return &InlineFunctions_FunctionVar{outer: outer, var_: var_}
}

// Java: InlineFunctions.FunctionVar.getName()
func (v *InlineFunctions_FunctionVar) GetName() string {
	return v.var_.GetFirstChild().GetString()
}

// Java: InlineFunctions.FunctionVar.getNameNode()
func (v *InlineFunctions_FunctionVar) GetNameNode() rhino.Node {
	return v.var_.GetFirstChild()
}

// Java: InlineFunctions.FunctionVar.getFunctionNode()
func (v *InlineFunctions_FunctionVar) GetFunctionNode() rhino.Node {
	return v.var_.GetFirstFirstChild()
}

// Java: InlineFunctions.FunctionVar.remove()
func (v *InlineFunctions_FunctionVar) Remove() {
	v.outer.compiler.ReportChangeToEnclosingScope(v.var_)
	NodeUtil_RemoveChild(v.var_.GetParent(), v.var_)
	NodeUtil_MarkFunctionsDeleted(v.var_, v.outer.compiler)
}

// Java: InlineFunctions.FunctionVar.getDeclaringBlock()
func (v *InlineFunctions_FunctionVar) GetDeclaringBlock() rhino.Node {
	return v.var_.GetParent()
}

// InlineFunctions_FunctionExpression is the FunctionExpression
// implementation of the Function interface.
type InlineFunctions_FunctionExpression struct {
	fn           rhino.Node
	fakeName     string
	fakeNameNode rhino.Node
}

// Java: InlineFunctions.FunctionExpression(Node, int)
func NewInlineFunctions_FunctionExpression(fn rhino.Node, index int) *InlineFunctions_FunctionExpression {
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
func (e *InlineFunctions_FunctionExpression) GetNameNode() rhino.Node {
	if e.fakeNameNode.IsNil() {
		e.fakeNameNode = rhino.IR_Name(e.fakeName)
	}
	return e.fakeNameNode
}

// Java: InlineFunctions.FunctionExpression.getFunctionNode()
func (e *InlineFunctions_FunctionExpression) GetFunctionNode() rhino.Node {
	return e.fn
}

// Java: InlineFunctions.FunctionExpression.remove()
func (e *InlineFunctions_FunctionExpression) Remove() {
	// Nothing to do. The function is removed with the call.
}

// Java: InlineFunctions.FunctionExpression.getDeclaringBlock()
func (e *InlineFunctions_FunctionExpression) GetDeclaringBlock() rhino.Node {
	return rhino.NilNode
}

// InlineFunctions_Reference is InlineFunctions.Reference, a subclass of
// FunctionInjector.Reference.
type InlineFunctions_Reference struct {
	FunctionInjector_ReferenceBase

	requiresDecomposition bool
	inlined               bool
}

// Java: InlineFunctions.Reference(Node, Scope, JSChunk, InliningMode)
func NewInlineFunctions_Reference(callNode rhino.Node, scope *Scope, chunk *JSChunk, mode FunctionInjector_InliningMode) *InlineFunctions_Reference {
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
