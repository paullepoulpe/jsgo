// Port of com/google/javascript/jscomp/InlineFunctions.java @ closure-compiler v20261028

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

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

type inlineFnKey struct {
	nameID   uint32
	anonNode rhino.Node
}

// Go-only: the shapes this pass reads (a var holding a function, the
// function a call calls, how a name is used, the params and body of a
// function, the members of a class body) go through the typed view
// (internal/ast); the subtree walks and the parallel bookkeeping stay on
// rhino.Node.

// inlineFunctionsRunState is the state of one run of InlineFunctions.Process,
// dropped when the run ends.
type inlineFunctionsRunState struct {
	fnKeys      []inlineFnKey
	anonOrder   []rhino.Node
	constFnIDs  []uint32
	calledIDs   smallIDSet
	refsScratch []*FunctionInjector_Reference
	movedArgs   map[rhino.Node]struct{}
	references  jutil.Slab[FunctionInjector_Reference]    // one run's references die together
	fnStates    jutil.Slab[InlineFunctions_FunctionState] // one run's function states die together
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
	fns      []*InlineFunctions_FunctionState
	namedFns map[uint32]*InlineFunctions_FunctionState
	anonFns  map[rhino.Node]*InlineFunctions_FunctionState

	compiler AbstractCompiler

	injector                 *FunctionInjector
	functionArgumentInjector *FunctionArgumentInjector

	reach CompilerOptions_Reach

	enforceMaxSizeAfterInlining bool
	maxSizeAfterInlining        int

	runState *inlineFunctionsRunState

	candScratch inlineCandidateScratch
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
	f.injector = NewFunctionInjector(compiler, FunctionInjectorOptions{
		SafeNameIdSupplier:       safeNameIdSupplier,
		AssumeStrictThis:         assumeStrictThis,
		AllowDecomposition:       true,
		FunctionArgumentInjector: f.functionArgumentInjector,
	})
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
		f.runState = &inlineFunctionsRunState{}
	}
	return f.runState
}

func (f *InlineFunctions) removeKey(k inlineFnKey) {
	if !k.anonNode.IsNil() {
		delete(f.anonFns, k.anonNode)
	} else {
		delete(f.namedFns, k.nameID)
	}
}

func (f *InlineFunctions) functionStateName(fs *InlineFunctions_FunctionState) string {
	if fs.target.name != "" {
		return fs.target.name
	}
	if fs.target.kind == fnKindExpr && f.runState != nil {
		if idx := slices.Index(f.runState.anonOrder, fs.target.fnNode); idx >= 0 {
			return strconv.Itoa(idx)
		}
	}
	return ""
}

func functionStateAliasCount(fs *InlineFunctions_FunctionState) int {
	if fs.summaryValid {
		return int(fs.summary.numModified)
	}
	return int(fs.aliasCount)
}

func (f *InlineFunctions) allocReference(callNode rhino.Node, scope *Scope, chunk *JSChunk, mode FunctionInjector_InliningMode) *FunctionInjector_Reference {
	r := f.ensureRunState().references.New()
	r.init(callNode, scope, chunk, mode, false)
	return r
}

// Java: InlineFunctions.getOrCreateFunctionState(String)
func (f *InlineFunctions) getOrCreateNamedFunctionState(nameID uint32) *InlineFunctions_FunctionState {
	p := f.ensureRunState()
	if fs := f.namedFns[nameID]; fs != nil {
		return fs
	}
	if f.namedFns == nil {
		f.namedFns = make(map[uint32]*InlineFunctions_FunctionState, 64)
	}
	fs := p.fnStates.New()
	fs.inline = true
	fs.remove = true
	f.namedFns[nameID] = fs
	f.fns = append(f.fns, fs)
	p.fnKeys = append(p.fnKeys, inlineFnKey{nameID: nameID})
	return fs
}

func (f *InlineFunctions) allocAnonFunctionState(fnNode rhino.Node) *InlineFunctions_FunctionState {
	p := f.ensureRunState()
	p.anonOrder = append(p.anonOrder, fnNode)
	if f.anonFns == nil {
		f.anonFns = make(map[rhino.Node]*InlineFunctions_FunctionState, 8)
	}
	fs := p.fnStates.New()
	fs.inline = true
	fs.remove = true
	f.anonFns[fnNode] = fs
	f.fns = append(f.fns, fs)
	p.fnKeys = append(p.fnKeys, inlineFnKey{anonNode: fnNode})
	return fs
}

func (f *InlineFunctions) addAnonFunction(fnNode rhino.Node, chunk *JSChunk) {
	fs := f.allocAnonFunctionState(fnNode)
	target := inlineFnTarget{kind: fnKindExpr, fnNode: fnNode}
	f.UpdateFunctionStateForInlining(target, chunk, fs)
	jutil.CheckState(!(InlineFunctions_HasRequireInliningAnnotation(fnNode) && !fs.CanInline()))
}

func inlineFunctions_traverseWithScopeCreator(compiler AbstractCompiler, root rhino.Node, cb NodeTraversal_Callback, sc ScopeCreator) {
	b := NodeTraversal_Builder{compiler: compiler, callback: cb, scopeCreator: sc}
	b.Traverse(root)
}

// Java: InlineFunctions.process(Node, Node)
func (f *InlineFunctions) Process(externs rhino.Node, root rhino.Node) {
	jutil.CheckState(f.compiler.GetLifeCycleStage().IsNormalized())
	if fastFault(f.compiler.GetOptions(), faultBenignNoInline) {
		return
	}

	p := f.ensureRunState()
	defer func() {
		f.injector.fnInfoCache = nil
		f.injector.clearKnownConstantFunctions()
		f.functionArgumentInjector.sideEffectCache = nil
		f.runState = nil
		f.fns = nil
		f.namedFns = nil
		f.anonFns = nil
		f.candScratch.lastCheckedBlock = rhino.NilNode
	}()

	if !root.IsNil() {
		if !f.findCandidateFunctionsParallel(root) {
			findCandidates := inlineCandidateFinder{outer: f}
			findCandidates.walk(root, rhino.NilNode, nil, false)
		}
	}
	if len(f.fns) == 0 {
		return // Nothing left to do.
	}
	if !f.findCandidatesReferencesParallel(root) {
		sc := NewSyntacticScopeCreator(f.compiler)
		var collector inlineRefCollector
		collector.initSerial(f)
		inlineFunctions_traverseWithScopeCreator(f.compiler, root, &collector.visitor, sc)
	}
	f.TrimCandidatesNotMeetingMinimumRequirements()
	if len(f.fns) == 0 {
		return // Nothing left to do.
	}

	// Store the set of function names eligible for inlining and use this to
	// prevent function names from being moved into temporaries during
	// expression decomposition. If this movement were allowed it would prevent
	// the Inline callback from finding the function calls.
	//
	// This pass already assumes these are constants, so this is safe for anyone
	// using function inlining.
	p.constFnIDs = make([]uint32, 0, len(f.namedFns))
	for _, k := range p.fnKeys {
		if k.anonNode.IsNil() {
			p.constFnIDs = append(p.constFnIDs, k.nameID)
		}
	}

	f.TrimCandidatesUsingOnCost()
	if len(f.fns) == 0 {
		return // Nothing left to do.
	}
	f.ResolveInlineConflicts()
	f.DecomposeExpressions()
	if !root.IsNil() {
		callVisitor := inlineCallVisitor{
			namedFns:   f.namedFns,
			anonFns:    f.anonFns,
			onCallSite: f.inlineCallSite,
		}
		callVisitor.inlineWalkRoot(f.compiler, root)
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

func inlineFunctions_walkCandidates(
	compiler AbstractCompiler,
	n rhino.Node,
	parent rhino.Node,
	chunk *JSChunk,
	inLocalHoist bool,
	visit func(n, parent rhino.Node, chunk *JSChunk, inLocalHoist bool),
) {
	tok := n.GetToken()
	childLocalHoist := inLocalHoist
	switch tok {
	case rhino.Token_SCRIPT:
		if id := n.GetInputId(); id != nil {
			if inp := compiler.GetInput(id); inp != nil {
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
		for ch, p := range ast.MustClassBody(n).EvalOrder() {
			hoist := inLocalHoist || (p.IsComputedFieldDef() && ch != p.GetFirstChild())
			inlineFunctions_walkCandidates(compiler, ch, p, chunk, hoist, visit)
		}
		return
	}
	for ch := n.GetFirstChild(); !ch.IsNil(); ch = ch.GetNext() {
		if !ch.GetFirstChild().IsNil() {
			inlineFunctions_walkCandidates(compiler, ch, n, chunk, childLocalHoist, visit)
		}
	}
	visit(n, parent, chunk, inLocalHoist)
}

// inlineFunctions_throwWalkException rethrows r, thrown while one of the
// walks that stand for the NodeTraversals of InlineFunctions.process
// (FindCandidateFunctions, CallVisitor) visited n in script, the way
// NodeTraversal.traverse rethrows what its callback throws:
// throwUnexpectedException wraps it in an INTERNAL COMPILER ERROR that names
// n and its parent.
func inlineFunctions_throwWalkException(compiler AbstractCompiler, r any, n, script rhino.Node) {
	t := NodeTraversal{compiler: compiler, currentNode: n, currentScript: script}
	t.ThrowUnexpectedException(jutil.AsThrowable(r))
}

// recoverWalkException recovers from panics during custom AST walks, clears the
// tracked node/script pointers, and reports unexpected exceptions with source
// context via NodeTraversal.ThrowUnexpectedException. Must be invoked directly
// via `defer recoverWalkException(...)`.
func recoverWalkException(compiler AbstractCompiler, root rhino.Node, nodePtr, scriptPtr *rhino.Node) {
	r := recover()
	var node, script rhino.Node
	if nodePtr != nil {
		node = *nodePtr
		*nodePtr = rhino.NilNode
	}
	if scriptPtr != nil {
		script = *scriptPtr
		*scriptPtr = rhino.NilNode
	}
	if r == nil {
		return
	}
	if script.IsNil() && !node.IsNil() {
		script = NodeUtil_GetEnclosingScript(node)
	}
	if script.IsNil() {
		script = NodeUtil_GetEnclosingScript(root)
	}
	inlineFunctions_throwWalkException(compiler, r, node, script)
}

func inlineFunctions_matchCandidateNode(reach CompilerOptions_Reach, n, parent rhino.Node, inLocalHoist bool) (inlineFnTarget, bool) {
	if !reach.IncludesGlobals() && !inLocalHoist {
		return inlineFnTarget{}, false
	}
	switch n.GetToken() {
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		if !NodeUtil_IsStatement(n) {
			return inlineFnTarget{}, false
		}
		// Functions expressions in the form of:
		//   var fooFn = function(x) { return ... }
		decl := ast.MustNameDeclaration(n)
		jutil.CheckStateMsg1(decl.NumDeclarators() == 1, "%s", n)
		if name, fn, ok := inlineFunctions_varFunction(decl); ok {
			return inlineFnTarget{
				kind:     fnKindVar,
				name:     name.Ident(),
				nameNode: name.Node(),
				fnNode:   fn.Node(),
				declNode: n,
			}, true
		}
	case rhino.Token_FUNCTION:
		if !NodeUtil_IsStatement(n) {
			return inlineFnTarget{}, false
		}
		jutil.CheckState(NodeUtil_IsStatementBlock(parent) || parent.IsLabel())
		if NodeUtil_IsFunctionDeclaration(n) {
			name, ok := ast.MustFunction(n).Name()
			jutil.CheckState(ok)
			return inlineFnTarget{
				kind:     fnKindDecl,
				name:     name.Ident(),
				nameNode: name.Node(),
				fnNode:   n,
				declNode: n,
			}, true
		}
	case rhino.Token_OPTCHAIN_CALL, rhino.Token_CALL:
		if fn, ok := inlineFunctions_calledFunction(ast.MustNormalOrOptChainCall(n)); ok {
			return inlineFnTarget{
				kind:   fnKindExpr,
				fnNode: fn.Node(),
			}, true
		}
	}
	return inlineFnTarget{}, false
}

// inlineCandidateFinder finds functions that might be inlined.
type inlineCandidateFinder struct {
	outer *InlineFunctions
	// current is the node being visited (NodeTraversal.currentNode), for
	// recoverWalkException.
	current rhino.Node
}

func (c *inlineCandidateFinder) walk(n rhino.Node, parent rhino.Node, chunk *JSChunk, inLocalHoist bool) {
	defer recoverWalkException(c.outer.compiler, n, &c.current, nil)
	inlineFunctions_walkCandidates(c.outer.compiler, n, parent, chunk, inLocalHoist, c.visitWalkNode)
}

func (c *inlineCandidateFinder) visitWalkNode(n, parent rhino.Node, chunk *JSChunk, inLocalHoist bool) {
	c.current = n
	if target, ok := inlineFunctions_matchCandidateNode(c.outer.reach, n, parent, inLocalHoist); ok {
		c.outer.MaybeAddFunction(target, chunk)
	}
}

type inlineCandidateAnalysis struct {
	target inlineFnTarget
	chunk  *JSChunk
	// visited is the node whose visit found the candidate, for
	// inlineFunctions_throwWalkException (set by the parallel walk only).
	visited             rhino.Node
	summary             calleeSummary
	disallowReasons     [5]InlineFunctions_DisallowInliningReason
	bodyAstSize         int32
	aliasCount          int32
	disallowCount       uint8
	hasSummary          bool
	disallowBeforeSetFn bool
	inlineDirectly      bool
	referencesThis      bool
	hasInnerClosures    bool
}

func (a *inlineCandidateAnalysis) addDisallow(reason InlineFunctions_DisallowInliningReason) {
	a.disallowReasons[a.disallowCount] = reason
	a.disallowCount++
}

func (a *inlineCandidateAnalysis) canInline() bool {
	return a.disallowCount == 0
}

func (a *inlineCandidateAnalysis) applyTo(
	f *InlineFunctions,
	fs *InlineFunctions_FunctionState,
) {
	if fs.HasExistingFunctionDefinition() {
		fs.DisallowInlining(f.compiler, InlineFunctions_DisallowInliningReason_MULTIPLE_FUNCTION_DEFINITIONS)
		return
	}
	if !fs.CanInline() {
		return
	}
	if a.disallowBeforeSetFn {
		fs.DisallowInlining(f.compiler, a.disallowReasons[0])
		return
	}
	fs.setTarget(a.target)
	if a.inlineDirectly {
		fs.InlineDirectly(true)
	}
	if a.disallowCount > 0 {
		for i := uint8(0); i < a.disallowCount; i++ {
			fs.DisallowInlining(f.compiler, a.disallowReasons[i])
		}
		return
	}
	fs.SetChunk(a.chunk)
	fs.aliasCount = a.aliasCount
	if a.referencesThis {
		fs.SetReferencesThis(true)
	}
	if a.hasInnerClosures {
		fs.SetHasInnerClosures(true)
	}
	fs.bodyAstSize = a.bodyAstSize
	if a.hasSummary {
		fs.summary = a.summary
		fs.summaryValid = true
	}
}

type inlineCandidateScratch struct {
	lastCheckedBlock              rhino.Node
	lastCheckedBlockHasLetOrConst bool
	scratchSummary                calleeSummary
}

func (s *inlineCandidateScratch) blockHasLetOrConst(block rhino.Node) bool {
	if block == s.lastCheckedBlock {
		return s.lastCheckedBlockHasLetOrConst
	}
	res := inlineFunctions_hasLetOrConst(block)
	s.lastCheckedBlock = block
	s.lastCheckedBlockHasLetOrConst = res
	return res
}

func (s *inlineCandidateScratch) checkCandidateAndScanSummary(f *InlineFunctions, target inlineFnTarget) bool {
	if target.kind != fnKindExpr {
		conv := f.compiler.GetCodingConvention()
		if conv.IsExported(target.name) || (!target.nameNode.IsNil() && conv.IsPropertyRenameFunction(target.nameNode)) {
			return false
		}
	}
	scanCalleeSummary(target.name, target.fnNode, &s.scratchSummary)
	return !s.scratchSummary.blocksMinReqs
}

func (s *inlineCandidateScratch) analyzeCandidate(
	f *InlineFunctions,
	target inlineFnTarget,
	chunk *JSChunk,
	debugLogging bool,
	parallel bool,
) inlineCandidateAnalysis {
	res := inlineCandidateAnalysis{
		target: target,
		chunk:  chunk,
	}
	fnNode := target.fnNode
	if f.HasNoInlineAnnotation(fnNode) {
		res.disallowBeforeSetFn = true
		res.addDisallow(InlineFunctions_DisallowInliningReason_NO_INLINE_ANNOTATION)
		return res
	}
	if f.enforceMaxSizeAfterInlining &&
		!InlineFunctions_IsAlwaysInlinable(fnNode) &&
		f.maxSizeAfterInlining <= NodeUtil_CountAstSizeUpToLimit(fnNode, f.maxSizeAfterInlining) {
		res.disallowBeforeSetFn = true
		res.addDisallow(InlineFunctions_DisallowInliningReason_TARGET_SIZE_EXCEEDS_LIMIT)
		return res
	}
	if FunctionInjector_IsDirectCallNodeReplacementPossible(fnNode) {
		res.inlineDirectly = true
	}
	if !debugLogging {
		if fnNode.IsGeneratorFunction() {
			res.addDisallow(InlineFunctions_DisallowInliningReason_GENERATOR_FUNCTION)
			return res
		}
		if fnNode.IsAsyncFunction() {
			res.addDisallow(InlineFunctions_DisallowInliningReason_ASYNC_FUNCTION)
			return res
		}
	}
	fn := ast.MustFunction(fnNode)
	if InlineFunctions_HasNonInlinableParam(fn.Params().Node()) {
		res.addDisallow(InlineFunctions_DisallowInliningReason_NON_INLINABLE_PARAM)
		if !debugLogging {
			return res
		}
	}

	if res.canInline() || debugLogging {
		if !s.checkCandidateAndScanSummary(f, target) {
			res.addDisallow(InlineFunctions_DisallowInliningReason_NOT_CANDIDATE_FUNCTION)
			if !debugLogging {
				return res
			}
		}
	}

	if res.canInline() {
		sum := &s.scratchSummary
		res.aliasCount = int32(sum.numModified)
		if sum.numModified > 0 {
			res.inlineDirectly = false
		}
		res.referencesThis = sum.referencesThis
		if sum.hasInnerClosures {
			res.hasInnerClosures = true
			if sum.numParams > 0 || sum.hasShallowDeclarations {
				res.addDisallow(InlineFunctions_DisallowInliningReason_HAS_LOCAL_NAMES)
				if !debugLogging {
					return res
				}
			}
		}
		if res.canInline() && f.enforceMaxSizeAfterInlining && !InlineFunctions_IsAlwaysInlinable(fnNode) {
			res.bodyAstSize = int32(min(int(sum.bodyAstSize), f.maxSizeAfterInlining))
		}
	}

	if fnNode.GetGrandparent().IsVar() &&
		(res.canInline() || (debugLogging && InlineFunctions_HasRequireInliningAnnotation(fnNode))) {
		block := target.declaringBlock()
		if block.IsBlock() && !block.GetParent().IsFunction() && s.blockHasLetOrConst(block) {
			res.addDisallow(InlineFunctions_DisallowInliningReason_CAPTURED_VARIABLE_NOT_IN_SCOPE)
			if !debugLogging {
				return res
			}
		}
	}

	if debugLogging {
		if fnNode.IsGeneratorFunction() {
			res.addDisallow(InlineFunctions_DisallowInliningReason_GENERATOR_FUNCTION)
		}
		if fnNode.IsAsyncFunction() {
			res.addDisallow(InlineFunctions_DisallowInliningReason_ASYNC_FUNCTION)
		}
	}

	if res.canInline() && s.scratchSummary.paramsPositional {
		res.summary = s.scratchSummary
		res.hasSummary = true
	}
	return res
}

// MaybeAddFunction updates the FunctionState object for the given function.
// Checks if the given function matches the criteria for an inlinable
// function.
//
// Java: InlineFunctions.maybeAddFunction(Function, JSChunk)
func (f *InlineFunctions) MaybeAddFunction(target inlineFnTarget, chunk *JSChunk) {
	if target.kind == fnKindExpr && target.name == "" {
		f.addAnonFunction(target.fnNode, chunk)
		return
	}
	var nameID uint32
	if !target.nameNode.IsNil() {
		nameID = target.nameNode.GetStringID()
	} else {
		nameID = uint32(strpool.Intern(target.name))
	}
	functionState := f.getOrCreateNamedFunctionState(nameID)
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
	if !functionState.CanInline() {
		return
	}
	analysis := f.candScratch.analyzeCandidate(
		f, target, chunk, f.compiler.IsDebugLoggingEnabled(), false)
	analysis.applyTo(f, functionState)
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

// TrimCandidatesNotMeetingMinimumRequirements removes entries that aren't a
// valid inline candidates, from the list of encountered names.
//
// Java: InlineFunctions.trimCandidatesNotMeetingMinimumRequirements()
func (f *InlineFunctions) TrimCandidatesNotMeetingMinimumRequirements() {
	p := f.runState
	w := 0
	for i, functionState := range f.fns {
		k := p.fnKeys[i]
		if !functionState.HasExistingFunctionDefinition() || !functionState.CanInline() {
			f.removeKey(k)
			continue
		}
		f.fns[w] = functionState
		p.fnKeys[w] = k
		w++
	}
	clear(f.fns[w:])
	f.fns = f.fns[:w]
	p.fnKeys = p.fnKeys[:w]
}

// TrimCandidatesUsingOnCost removes entries from the list of candidates that
// can't be inlined.
//
// Java: InlineFunctions.trimCandidatesUsingOnCost()
func (f *InlineFunctions) TrimCandidatesUsingOnCost() {
	p := f.runState
	w := 0
	for i, functionState := range f.fns {
		k := p.fnKeys[i]
		if functionState.HasReferences() {
			// Only inline function if it decreases the code size.
			lowersCost := f.MinimizeCost(functionState)
			if !lowersCost {
				// It shouldn't be inlined; remove it from the list.
				f.removeKey(k)
				continue
			}
		} else if !functionState.CanRemove() {
			// Don't bother tracking functions without references that can't be
			// removed.
			f.removeKey(k)
			continue
		}
		f.fns[w] = functionState
		p.fnKeys[w] = k
		w++
	}
	clear(f.fns[w:])
	f.fns = f.fns[:w]
	p.fnKeys = p.fnKeys[:w]
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
	summary := functionState.ensureSummary()
	return f.injector.inliningLowersCost(
		functionState.GetChunk(),
		functionState.functionNode(),
		functionState.GetReferences(),
		functionStateAliasCount(functionState),
		functionState.CanRemove(),
		functionState.GetReferencesThis(),
		summary)
}
