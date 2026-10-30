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
			cand.applyTo(f, fs)
			jutil.CheckState(!(InlineFunctions_HasRequireInliningAnnotation(cand.target.fnNode) && !fs.CanInline()))
		}
	}
	return true
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
		var cb *inlineRefCollector
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
				cb = &inlineRefCollector{}
				localSafeSupplier := func() string {
					id := cb.safeIdsCounter
					cb.safeIdsCounter++
					return strconv.Itoa(id)
				}
				argInj := NewFunctionArgumentInjector(f.compiler.GetAstAnalyzer())
				argInj.readOnlySummaries = true
				inj := NewFunctionInjector(rc, FunctionInjectorOptions{
					SafeNameIdSupplier:       localSafeSupplier,
					AllowDecomposition:       f.injector.allowDecomposition,
					AssumeStrictThis:         f.injector.assumeStrictThis,
					FunctionArgumentInjector: argInj,
				})
				inj.shareKnownConstantFunctionsFrom(f.injector)
				cb.initParallel(f, inj)
				sc = NewSyntacticScopeCreator(rc)
			} else {
				cb.resetForScript()
			}
			b := NodeTraversal_Builder{
				compiler:     rc,
				callback:     &cb.visitor,
				scopeCreator: sc,
			}
			b.traverseScriptUnderRoot(root, root, scripts[idx], lazyGlobal)
			ifEventsByScript[idx] = cb.events
		}
		return runItem, cleanup
	}, false)

	for idx := range ifEventsByScript {
		f.replayPhase2Events(ifEventsByScript[idx])
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
	return f.injector.inliningLowersCost(
		functionState.GetChunk(),
		functionState.functionNode(),
		functionState.refs,
		functionStateAliasCount(functionState),
		functionState.CanRemove(),
		functionState.GetReferencesThis(),
		summary)
}
