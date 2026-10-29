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
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type ifPhase1Worker struct {
	outer   *InlineFunctions
	scratch inlineCandidateScratch
	cands   []inlineCandidateAnalysis
}

func (w *ifPhase1Worker) walk(n rhino.Node, parent rhino.Node, chunk *JSChunk, inLocalHoist bool) {
	inlineFunctions_walkCandidates(w.outer.compiler, n, parent, chunk, inLocalHoist, w.visitWalkNode)
}

func (w *ifPhase1Worker) visitWalkNode(n, parent rhino.Node, chunk *JSChunk, inLocalHoist bool) {
	if target, ok := inlineFunctions_matchCandidateNode(w.outer.reach, n, parent, inLocalHoist); ok {
		w.cands = append(w.cands, w.scratch.analyzeCandidate(w.outer, target, chunk, false, true))
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

	candsByScript := make([][]inlineCandidateAnalysis, len(scripts))
	runParallelWorkersIndexed(
		len(scripts),
		numWorkers,
		func(_ int) ifPhase1Worker {
			return ifPhase1Worker{outer: f}
		},
		func(w *ifPhase1Worker, idx int) {
			w.cands = nil
			w.scratch.lastCheckedBlock = rhino.NilNode
			w.walk(scripts[idx], root, nil, false)
			candsByScript[idx] = w.cands
		},
		nil,
	)

	for idx := range scripts {
		for i := range candsByScript[idx] {
			cand := &candsByScript[idx][i]
			var fs *InlineFunctions_FunctionState
			if cand.target.kind == fnKindExpr {
				fs = f.allocAnonFunctionState(cand.target.fnNode)
			} else {
				fs = f.getOrCreateNamedFunctionState(cand.target.nameNode.GetStringID())
			}
			cand.applyTo(f, fs, nil)
			jutil.CheckState(!(InlineFunctions_HasRequireInliningAnnotation(cand.target.fnNode) && !fs.CanInline()))
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

var (
	_ nodeTraversal_visitTokens      = (*ifPhase2RefCallback)(nil)
	_ nodeTraversal_onlyNameLeaves   = (*ifPhase2RefCallback)(nil)
	_ nodeTraversal_noShouldTraverse = (*ifPhase2RefCallback)(nil)
)

func (*ifPhase2RefCallback) ShouldTraverse(*NodeTraversal, rhino.Node, rhino.Node) bool { return true }
func (*ifPhase2RefCallback) nodeTraversalNoShouldTraverse()                             {}
func (*ifPhase2RefCallback) onlyNameLeaves()                                            {}
func (*ifPhase2RefCallback) visitTokens() *ast.TokenSet                                 { return &inlineCallAndNameVisitTokens }

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
	return inlineTargetSizeExceedsLimit(
		t, fs, r.outer.maxSizeAfterInlining, &r.lastContainingFn, &r.lastContainingFnSize, false)
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
	ok, scope, reqDecomp := inlineCheckCallSiteWithMode(
		t, r.injector, &r.scratchRef, functionState, functionState.readOnlySummary(), callNode, chunk, mode)
	return ok, scope, reqDecomp, mode
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
