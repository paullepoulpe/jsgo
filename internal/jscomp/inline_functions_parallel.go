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
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type ifPhase1Worker struct {
	outer   *InlineFunctions
	scratch inlineCandidateScratch
	cands   []inlineCandidateAnalysis
	current rhino.Node // as InlineFunctions_FindCandidateFunctions.current
}

func (w *ifPhase1Worker) walk(n rhino.Node, parent rhino.Node, chunk *JSChunk, inLocalHoist bool) {
	defer func() {
		if r := recover(); r != nil {
			inlineFunctions_throwWalkException(w.outer.compiler, r, w.current, NodeUtil_GetEnclosingScript(w.current))
		}
	}()
	inlineFunctions_walkCandidates(w.outer.compiler, n, parent, chunk, inLocalHoist, w.visitWalkNode)
}

func (w *ifPhase1Worker) visitWalkNode(n, parent rhino.Node, chunk *JSChunk, inLocalHoist bool) {
	w.current = n
	if target, ok := inlineFunctions_matchCandidateNode(w.outer.reach, n, parent, inLocalHoist); ok {
		w.cands = append(w.cands, w.scratch.analyzeCandidate(w.outer, target, chunk, false, true))
		w.cands[len(w.cands)-1].visited = n
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

	f.applyParallelCandidates(candsByScript)
	return true
}

// applyParallelCandidates is the rest of Java's maybeAddFunction for the
// candidates of findCandidateFunctionsParallel, in script order.
func (f *InlineFunctions) applyParallelCandidates(candsByScript [][]inlineCandidateAnalysis) {
	// Java runs it in the visit that found the candidate: what it throws is
	// rethrown as the walk rethrows it.
	var visited rhino.Node
	defer func() {
		if r := recover(); r != nil {
			inlineFunctions_throwWalkException(f.compiler, r, visited, NodeUtil_GetEnclosingScript(visited))
		}
	}()
	for _, cands := range candsByScript {
		for i := range cands {
			cand := &cands[i]
			visited = cand.visited
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
