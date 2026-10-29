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
	"sync"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

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
	nameID              uint32
	name                string
	chunk               *JSChunk
	disallowReason      InlineFunctions_DisallowInliningReason
	hasDisallow         bool
	disallowBeforeSetFn bool
	inlineDirectly      bool
	referencesThis      bool
	hasInnerClosures    bool
	namesToAliasIDs     []uint32
	bodyAstSize         int32
	summary             *calleeSummary
}

type ifPhase1Worker struct {
	outer                         *InlineFunctions
	lastCheckedBlock              rhino.Node
	lastCheckedBlockHasLetOrConst bool
	cands                         []ifPhase1Candidate
	scratchSummary                calleeSummary
	scratchParamIDs               inlineIDSet
	scratchModifiedIDs            inlineIDSet
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

func (w *ifPhase1Worker) buildNamedOrVarCandidate(kind ifCandKind, node rhino.Node, nameNode rhino.Node, fnNode rhino.Node, chunk *JSChunk) ifPhase1Candidate {
	f := w.outer
	nameID := nameNode.GetStringID()
	name := rhino.RhinoStringPool_StringFromID(nameID)
	cand := ifPhase1Candidate{
		kind:   kind,
		node:   node,
		fnNode: fnNode,
		nameID: nameID,
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
	fn := ast.MustFunction(fnNode)
	if InlineFunctions_HasNonInlinableParam(fn.Params().Node()) {
		cand.hasDisallow = true
		cand.disallowReason = InlineFunctions_DisallowInliningReason_NON_INLINABLE_PARAM
		return cand
	}
	conv := f.compiler.GetCodingConvention()
	if conv.IsExported(name) || conv.IsPropertyRenameFunction(nameNode) {
		cand.hasDisallow = true
		cand.disallowReason = InlineFunctions_DisallowInliningReason_NOT_CANDIDATE_FUNCTION
		return cand
	}
	summaryValid := scanCalleeSummary(name, fnNode, &w.scratchSummary)
	if summaryValid {
		s := &w.scratchSummary
		if s.blocksMinReqs {
			cand.hasDisallow = true
			cand.disallowReason = InlineFunctions_DisallowInliningReason_NOT_CANDIDATE_FUNCTION
			return cand
		}
		if s.numModified > 0 {
			cand.inlineDirectly = false
		}
		cand.referencesThis = s.referencesThis
		if s.hasInnerClosures {
			cand.hasInnerClosures = true
			if s.numParams > 0 || s.hasShallowDeclarations {
				cand.hasDisallow = true
				cand.disallowReason = InlineFunctions_DisallowInliningReason_HAS_LOCAL_NAMES
				return cand
			}
		}
	} else {
		if !f.injector.DoesFunctionMeetMinimumRequirements(name, fnNode) {
			cand.hasDisallow = true
			cand.disallowReason = InlineFunctions_DisallowInliningReason_NOT_CANDIDATE_FUNCTION
			return cand
		}
		body := fn.Body().Node()
		collectFunctionParameterIDs(fnNode, &w.scratchParamIDs)
		if !w.scratchParamIDs.isEmpty() {
			w.scratchModifiedIDs.clear()
			collectModifiedParameterIDs(body, &w.scratchParamIDs, false, &w.scratchModifiedIDs)
			if !w.scratchModifiedIDs.isEmpty() {
				cand.inlineDirectly = false
				cand.namesToAliasIDs = slices.Clone(w.scratchModifiedIDs.items)
			}
		}
		cand.referencesThis = NodeUtil_ReferencesEnclosingReceiver(body)
		if NodeUtil_Has(body, NodeUtil_IsFunctionOrNonStaticClassFieldInitializer, inlineFunctions_alwaysTrue) {
			cand.hasInnerClosures = true
			if InlineFunctions_HasLocalNames(fnNode) {
				cand.hasDisallow = true
				cand.disallowReason = InlineFunctions_DisallowInliningReason_HAS_LOCAL_NAMES
				return cand
			}
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
		if summaryValid {
			cand.bodyAstSize = int32(min(int(w.scratchSummary.bodyAstSize), f.maxSizeAfterInlining))
		} else {
			cand.bodyAstSize = int32(NodeUtil_CountAstSizeUpToLimit(fn.Body().Node(), f.maxSizeAfterInlining))
		}
	}
	if summaryValid {
		sumCopy := w.scratchSummary
		cand.summary = &sumCopy
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
		inlineFunctions_walkClassMembers(ast.MustClassBody(n), inLocalHoist, func(n, parent rhino.Node, inLocalHoist bool) {
			w.walk(n, parent, chunk, inLocalHoist)
		})
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
			decl := ast.MustNameDeclaration(n)
			jutil.CheckStateMsg1(decl.NumDeclarators() == 1, "%s", n)
			if name, fn, ok := inlineFunctions_varFunction(decl); ok {
				w.cands = append(w.cands, w.buildNamedOrVarCandidate(ifCandVar, n, name.Node(), fn.Node(), chunk))
			}
		}
	case rhino.Token_FUNCTION:
		if (w.outer.reach.IncludesGlobals() || inLocalHoist) && NodeUtil_IsStatement(n) {
			jutil.CheckState(NodeUtil_IsStatementBlock(parent) || parent.IsLabel())
			if NodeUtil_IsFunctionDeclaration(n) {
				name, _ := ast.MustFunction(n).Name()
				w.cands = append(w.cands, w.buildNamedOrVarCandidate(ifCandNamedFn, n, name.Node(), n, chunk))
			}
		}
	case rhino.Token_OPTCHAIN_CALL, rhino.Token_CALL:
		if w.outer.reach.IncludesGlobals() || inLocalHoist {
			if fn, ok := inlineFunctions_calledFunction(ast.MustNormalOrOptChainCall(n)); ok {
				fnNode := fn.Node()
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

	for idx := range scripts {
		if p := panicsByScript[idx]; p != nil {
			panic(p)
		}
		for i := range candsByScript[idx] {
			cand := &candsByScript[idx][i]
			if cand.kind == ifCandAnonExpr {
				f.addAnonFunction(cand.fnNode, cand.chunk)
				continue
			}
			fs := f.getOrCreateNamedFunctionState(cand.nameID)
			if fs.HasExistingFunctionDefinition() {
				fs.DisallowInlining(f.compiler, InlineFunctions_DisallowInliningReason_MULTIPLE_FUNCTION_DEFINITIONS)
			} else if fs.CanInline() {
				if cand.disallowBeforeSetFn {
					fs.DisallowInlining(f.compiler, cand.disallowReason)
				} else {
					var target inlineFnTarget
					if cand.kind == ifCandVar {
						target = newVarFnTarget(cand.node)
					} else {
						target = newNamedFnTarget(cand.node)
					}
					fs.setTarget(target)
					if cand.hasDisallow {
						fs.DisallowInlining(f.compiler, cand.disallowReason)
					} else {
						if cand.summary != nil {
							fs.summary = cand.summary
							fs.summaryValid = true
						}
						if cand.inlineDirectly {
							fs.InlineDirectly(true)
						}
						fs.SetChunk(cand.chunk)
						if len(cand.namesToAliasIDs) > 0 {
							fs.InlineDirectly(false)
							fs.setNamesToAliasIDSlice(cand.namesToAliasIDs)
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
	visitor              inlineCallVisitor
	fnByID               map[uint32]*InlineFunctions_FunctionState
	outer                *InlineFunctions
	injector             *FunctionInjector
	safeIdsCounter       int
	scratchRef           FunctionInjector_Reference
	lastContainingFn     rhino.Node
	lastContainingFnSize int
	localDisallowed      map[*InlineFunctions_FunctionState]struct{}
	events               []ifPhase2Event
}

func (*ifPhase2RefCallback) ShouldTraverse(*NodeTraversal, rhino.Node, rhino.Node) bool { return true }
func (*ifPhase2RefCallback) onlyNameLeaves()                                            {}

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
	inlinedFun := fs.functionNode()
	if InlineFunctions_IsAlwaysInlinable(inlinedFun) {
		return false
	}
	inlinedFunSize := int(fs.bodyAstSize)
	if inlinedFunSize == 0 {
		inlinedFunSize = NodeUtil_CountAstSizeUpToLimit(ast.MustFunction(inlinedFun).Body().Node(), r.outer.maxSizeAfterInlining)
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
	switch n.GetToken() {
	case rhino.Token_NAME:
		r.checkNameUsage(n, parent)
	case rhino.Token_OPTCHAIN_CALL, rhino.Token_CALL:
		r.visitor.visitCallNode(t, n)
	}
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
	r.scratchRef.init(callNode, nil, chunk, mode, false)
	r.scratchRef.t = t
	summary := functionState.readOnlySummary()
	result := r.injector.canInlineReferenceWithSummary(
		&r.scratchRef,
		functionState.functionNode(),
		functionStateAliasCount(functionState),
		functionState.GetReferencesThis(),
		functionState.HasInnerClosures(),
		summary)
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

	id := n.GetStringID()
	if id == rhino.RhinoStringPool_EmptyID {
		return
	}
	functionState := r.fnByID[id]
	if functionState == nil || InlineFunctions_IsCandidateUsage(n) {
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
	} else if inlineFunctions_isConstAlias(n, parent) {
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
		for _, fs := range f.fns {
			if fs.CanInline() && fs.bodyAstSize == 0 {
				fnNode := fs.functionNode()
				if !InlineFunctions_IsAlwaysInlinable(fnNode) {
					fs.bodyAstSize = int32(NodeUtil_CountAstSizeUpToLimit(ast.MustFunction(fnNode).Body().Node(), f.maxSizeAfterInlining))
				}
			}
		}
	}

	lazyGlobal := newLazyGlobalScope(f.compiler, root)
	workerFnCaches := make([]map[rhino.Node]functionInjectorFnInfo, numWorkers)
	workerSideEffects := make([]map[rhino.Node]calleeSideEffectSummary, numWorkers)

	ifEventsByScript := make([][]ifPhase2Event, len(scripts))
	runParallelWork(f.compiler, numWorkers, len(scripts), func(workerID int) (func(int, *recordingCompiler), func()) {
		var workerRC *recordingCompiler
		var cb *ifPhase2RefCallback
		var sc *SyntacticScopeCreator

		cleanup := func() {
			if cb != nil && cb.injector != nil {
				if len(cb.injector.fnInfoCache) != 0 {
					workerFnCaches[workerID] = cb.injector.fnInfoCache
				}
				if len(cb.injector.functionArgumentInjector.sideEffectCache) != 0 {
					workerSideEffects[workerID] = cb.injector.functionArgumentInjector.sideEffectCache
				}
			}
		}
		runItem := func(idx int, rc *recordingCompiler) {
			if cb == nil || workerRC != rc {
				workerRC = rc
				cb = &ifPhase2RefCallback{outer: f, fnByID: f.namedFns}
				cb.visitor = inlineCallVisitor{
					namedFns:   f.namedFns,
					anonFns:    f.anonFns,
					onCallSite: cb.VisitCallSite,
				}
				localSafeSupplier := func() string {
					id := cb.safeIdsCounter
					cb.safeIdsCounter++
					return strconv.Itoa(id)
				}
				argInj := NewFunctionArgumentInjector(f.compiler.GetAstAnalyzer())
				argInj.readOnlySummaries = true
				injBuilder := NewFunctionInjector_Builder(rc).
					SafeNameIdSupplier(localSafeSupplier).
					AllowDecomposition(f.injector.allowDecomposition).
					AssumeStrictThis(f.injector.assumeStrictThis).
					FunctionArgumentInjector(argInj)
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

	for _, seCache := range workerSideEffects {
		if len(seCache) == 0 {
			continue
		}
		if f.functionArgumentInjector.sideEffectCache == nil {
			f.functionArgumentInjector.sideEffectCache = make(map[rhino.Node]calleeSideEffectSummary, len(seCache)*2)
		}
		// order-independent: merges deterministic per-fnNode side-effect bitmasks computed identically by workers
		for fnNode, se := range seCache {
			f.functionArgumentInjector.sideEffectCache[fnNode] = se
		}
	}

	return true
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
	return f.injector.inliningLowersCostWithSummary(
		functionState.GetChunk(),
		functionState.functionNode(),
		functionState.refs,
		functionStateAliasCount(functionState),
		functionState.CanRemove(),
		functionState.GetReferencesThis(),
		summary)
}
