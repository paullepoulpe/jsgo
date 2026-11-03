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
	"fmt"
	"slices"
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type ifPhase1Worker struct {
	outer        *InlineFunctions
	scratch      *inlineCandidateScratch
	localScratch inlineCandidateScratch
	cands        []inlineCandidateAnalysis
	current      rhino.Node // as InlineFunctions_FindCandidateFunctions.current
}

func (w *ifPhase1Worker) walk(n rhino.Node, parent rhino.Node, chunk *JSChunk, inLocalHoist bool) {
	defer func() {
		if r := recover(); r != nil {
			inlineFunctions_throwWalkException(w.outer.compiler, r, w.current, NodeUtil_GetEnclosingScript(w.current))
		}
	}()
	inlineFunctions_walkCandidates(w.outer.compiler, n, parent, chunk, inLocalHoist, w.scratch, w.visitWalkNode)
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
	var scratches []inlineCandidateScratch
	if f.activeCache != nil {
		scratches = f.activeCache.prepareWorkerScratches(numWorkers, f.sideEffectCfg)
	}
	runParallelWorkersIndexed(
		len(scripts),
		numWorkers,
		func(wIdx int) ifPhase1Worker {
			w := ifPhase1Worker{outer: f}
			if scratches != nil {
				w.scratch = &scratches[wIdx]
			} else {
				w.scratch = &w.localScratch
			}
			return w
		},
		func(w *ifPhase1Worker, idx int) {
			w.cands = nil
			w.scratch.lastCheckedBlock = rhino.NilNode
			w.walk(scripts[idx], root, nil, false)
			candsByScript[idx] = w.cands
		},
		nil,
	)

	if f.activeCache != nil {
		totalLive := 0
		for i := range scratches {
			totalLive += len(scratches[i].liveFns)
		}
		f.activeCache.beginPassMerge(totalLive, f.sideEffectCfg)
		for i := range scratches {
			f.activeCache.applyScratch(&scratches[i])
		}
		f.activeCache.maybePruneDead(totalLive, scratches)
		for i := range scratches {
			scratches[i].clearWalkState()
		}
	}

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
	var workerFnCaches []map[rhino.Node]functionInjectorFnInfo
	var workerSideEffects []map[rhino.Node]calleeSideEffectSummary
	var ifEventsByScript [][]ifPhase2Event

	ok := traverseScriptsParallel(f.compiler, root, parallelScriptOptions[*inlineRefCollector]{
		useSharedGlobalScope: true,
		prepare: func(numWorkers, scriptCount int) bool {
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
			workerFnCaches = make([]map[rhino.Node]functionInjectorFnInfo, numWorkers)
			workerSideEffects = make([]map[rhino.Node]calleeSideEffectSummary, numWorkers)
			ifEventsByScript = make([][]ifPhase2Event, scriptCount)
			return true
		},
		initWorker: func(_ int, rc *Compiler) *inlineRefCollector {
			cb := &inlineRefCollector{}
			argInj := NewFunctionArgumentInjector(f.compiler.GetAstAnalyzer())
			argInj.readOnlySummaries = true
			inj := NewFunctionInjector(rc, FunctionInjectorOptions{
				SafeNameIdSupplier: func() string {
					id := cb.safeIdsCounter
					cb.safeIdsCounter++
					return strconv.Itoa(id)
				},
				AllowDecomposition:       f.injector.allowDecomposition,
				AssumeStrictThis:         f.injector.assumeStrictThis,
				FunctionArgumentInjector: argInj,
			})
			inj.seedCache = f.activeCache
			inj.shareKnownConstantFunctionsFrom(f.injector)
			cb.initParallel(f, inj)
			return cb
		},
		callback: func(cb *inlineRefCollector, _ int, _, _ rhino.Node, _ *Compiler) NodeTraversal_Callback {
			cb.resetForScript()
			return &cb.visitor
		},
		afterScript: func(cb *inlineRefCollector, idx int, _, _ rhino.Node, _ *Compiler) {
			ifEventsByScript[idx] = cb.events
		},
		finishWorker: func(workerID int, cb *inlineRefCollector) {
			cb.injector.seedCache = nil
			if len(cb.injector.fnInfoCache) != 0 {
				workerFnCaches[workerID] = cb.injector.fnInfoCache
			}
			if len(cb.injector.functionArgumentInjector.sideEffectCache) != 0 {
				workerSideEffects[workerID] = cb.injector.functionArgumentInjector.sideEffectCache
			}
		},
	})
	if !ok {
		return false
	}

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
				if existing.innerState == fnInnerUnknown && info.innerState != fnInnerUnknown {
					existing.innerState = info.innerState
					existing.inner = info.inner
				} else if existing.inner.IsNil() && !info.inner.IsNil() {
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

func (s *inlineCandidateScratch) resetForWalk(cache *inlineFunctionsCache, seCfg inlineSideEffectConfig) {
	s.lastCheckedBlock = rhino.NilNode
	s.cache = cache
	s.sideEffectCfg = seCfg
	s.fnFrames = s.fnFrames[:0]
	s.innerBuf = s.innerBuf[:0]
	s.liveFns = s.liveFns[:0]
	s.unchangedFns = s.unchangedFns[:0]
	s.stamps = s.stamps[:0]
	s.stampInners = s.stampInners[:0]
	s.candUpdates = s.candUpdates[:0]
	if len(s.candSummaries) > 0 {
		clear(s.candSummaries)
		s.candSummaries = s.candSummaries[:0]
	}
	if len(s.nonLeafSubtree) > 0 {
		clear(s.nonLeafSubtree)
	}
	s.candidateHits = 0
	s.candidateMisses = 0
	s.leafWalkSkips = 0
	s.sideEffectHits = 0
}

func (s *inlineCandidateScratch) clearWalkState() {
	s.resetForWalk(nil, inlineSideEffectConfig{})
}

func (s *inlineCandidateScratch) trySkipLeafFunctionChildren(fnNode rhino.Node) bool {
	if len(s.fnFrames) > 0 {
		s.innerBuf = append(s.innerBuf, fnNode)
	}
	s.liveFns = append(s.liveFns, fnNode)
	prev, ok := s.cache.lookup(fnNode)
	ownTime := fnNode.GetChangeTime()
	ownUnchanged := ok && !fnNode.IsDeleted() && prev.ownChangeTime == ownTime
	if ownUnchanged && prev.hasFlag(fnCacheHasDirectInnerFns) && prev.innerLen == 0 && !inlineFnCacheVerify {
		s.leafWalkSkips++
		s.unchangedFns = append(s.unchangedFns, fnNode)
		return true
	}
	s.fnFrames = append(s.fnFrames, inlineFnWalkFrame{
		fnNode:                    fnNode,
		innerStart:                len(s.innerBuf),
		allInnerSubtreesUnchanged: true,
	})
	return false
}

func (s *inlineCandidateScratch) finishFunctionWalk(fnNode rhino.Node) {
	last := len(s.fnFrames) - 1
	frame := s.fnFrames[last]
	s.fnFrames = s.fnFrames[:last]

	directInners := s.innerBuf[frame.innerStart:]
	prev, ok := s.cache.lookup(fnNode)
	ownTime := fnNode.GetChangeTime()
	ownUnchanged := ok && !fnNode.IsDeleted() && prev.ownChangeTime == ownTime
	sameInners := ok && prev.hasFlag(fnCacheHasDirectInnerFns) && slices.Equal(directInners, s.cache.directInnersOf(prev))
	numInners := len(directInners)

	subtreeUnchanged := ownUnchanged && sameInners && frame.allInnerSubtreesUnchanged
	if !subtreeUnchanged && len(s.fnFrames) > 0 {
		s.fnFrames[len(s.fnFrames)-1].allInnerSubtreesUnchanged = false
	}
	if numInners > 0 || inlineFnCacheVerify {
		if s.nonLeafSubtree == nil {
			s.nonLeafSubtree = make(map[rhino.Node]bool, 16)
		}
		s.nonLeafSubtree[fnNode] = subtreeUnchanged
	}

	if subtreeUnchanged {
		s.unchangedFns = append(s.unchangedFns, fnNode)
	} else {
		var off, ln uint32
		if !sameInners && numInners > 0 {
			off = uint32(len(s.stampInners))
			ln = uint32(numInners)
			s.stampInners = append(s.stampInners, directInners...)
		}
		s.stamps = append(s.stamps, inlineFnStampUpdate{
			fnNode:        fnNode,
			ownChangeTime: ownTime,
			innerOff:      off,
			innerLen:      ln,
			ownUnchanged:  ownUnchanged,
			sameInners:    sameInners,
		})
	}
	s.innerBuf = s.innerBuf[:frame.innerStart]
}

func makeCandUpdateFlags(alwaysInlinable, directReplace, isGen, isAsync, hasNonInlinableParam, localClosureDisallowed bool) uint8 {
	var f uint8
	if alwaysInlinable {
		f |= fnCacheAlwaysInlinable
	}
	if directReplace {
		f |= fnCacheDirectReplace
	}
	if isGen {
		f |= fnCacheIsGenerator
	}
	if isAsync {
		f |= fnCacheIsAsync
	}
	if hasNonInlinableParam {
		f |= fnCacheHasNonInlinableParam
	}
	if localClosureDisallowed {
		f |= fnCacheShallowLocalClosureDisallowed
	}
	return f
}

func (s *inlineCandidateScratch) analyzeCandidateCached(
	f *InlineFunctions,
	target inlineFnTarget,
	chunk *JSChunk,
	debugLogging bool,
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

	prev, hasPrev := s.cache.lookup(fnNode)
	ownUnchanged := hasPrev && !fnNode.IsDeleted() && prev.ownChangeTime == fnNode.GetChangeTime()
	subtreeUnchanged := false
	if ownUnchanged && prev.hasFlag(fnCacheHasDirectInnerFns) {
		if prev.innerLen == 0 && !inlineFnCacheVerify {
			subtreeUnchanged = true
		} else if s.nonLeafSubtree != nil {
			subtreeUnchanged = s.nonLeafSubtree[fnNode]
		}
	}

	var countedSize int32
	var countedLimit int32
	var alwaysInlinable bool
	haveAlwaysInlinable := false
	if ownUnchanged && prev.hasFlag(fnCacheShallowReady) {
		alwaysInlinable = prev.hasFlag(fnCacheAlwaysInlinable)
		haveAlwaysInlinable = true
	}

	if f.enforceMaxSizeAfterInlining {
		if !haveAlwaysInlinable {
			alwaysInlinable = InlineFunctions_IsAlwaysInlinable(fnNode)
			haveAlwaysInlinable = true
		}
		if !alwaysInlinable {
			limit := f.maxSizeAfterInlining
			if subtreeUnchanged && prev.fullAstSizeLimit > 0 &&
				(int(prev.fullAstSizeLimit) == limit || prev.fullAstSize < prev.fullAstSizeLimit) {
				if limit <= int(prev.fullAstSize) {
					res.disallowBeforeSetFn = true
					res.addDisallow(InlineFunctions_DisallowInliningReason_TARGET_SIZE_EXCEEDS_LIMIT)
					s.candidateHits++
					if inlineFnCacheVerify {
						verifyCachedCandidateAnalysis(f, target, chunk, debugLogging, &res)
					}
					return res
				}
			} else {
				sz := NodeUtil_CountAstSizeUpToLimit(fnNode, limit)
				countedSize = int32(sz)
				countedLimit = int32(limit)
				if limit <= sz {
					res.disallowBeforeSetFn = true
					res.addDisallow(InlineFunctions_DisallowInliningReason_TARGET_SIZE_EXCEEDS_LIMIT)
					s.candidateMisses++
					s.candUpdates = append(s.candUpdates, inlineFnCandUpdate{
						fnNode:           fnNode,
						fullAstSize:      countedSize,
						fullAstSizeLimit: countedLimit,
						flags: makeCandUpdateFlags(
							alwaysInlinable,
							FunctionInjector_IsDirectCallNodeReplacementPossible(fnNode),
							fnNode.IsGeneratorFunction(),
							fnNode.IsAsyncFunction(),
							InlineFunctions_HasNonInlinableParam(ast.MustFunction(fnNode).Params().Node()),
							false,
						),
					})
					return res
				}
			}
		}
	}

	targetNameID := effectiveCalleeFnNameID(target)

	// Fast Path 1: Shallow scope cache hit for non-debug disallowed candidates
	// (including large IIFE/UMD outer wrappers whose inner closures mutated).
	if ownUnchanged && prev.hasFlag(fnCacheShallowReady) && !debugLogging && countedLimit == 0 {
		if prev.hasFlag(fnCacheDirectReplace) {
			res.inlineDirectly = true
		}
		if prev.hasFlag(fnCacheIsGenerator) {
			res.addDisallow(InlineFunctions_DisallowInliningReason_GENERATOR_FUNCTION)
			s.candidateHits++
			if inlineFnCacheVerify {
				verifyCachedCandidateAnalysis(f, target, chunk, debugLogging, &res)
			}
			return res
		}
		if prev.hasFlag(fnCacheIsAsync) {
			res.addDisallow(InlineFunctions_DisallowInliningReason_ASYNC_FUNCTION)
			s.candidateHits++
			if inlineFnCacheVerify {
				verifyCachedCandidateAnalysis(f, target, chunk, debugLogging, &res)
			}
			return res
		}
		if prev.hasFlag(fnCacheHasNonInlinableParam) {
			res.addDisallow(InlineFunctions_DisallowInliningReason_NON_INLINABLE_PARAM)
			s.candidateHits++
			if inlineFnCacheVerify {
				verifyCachedCandidateAnalysis(f, target, chunk, debugLogging, &res)
			}
			return res
		}
		if isExportedOrRenameCandidate(f, target) {
			res.addDisallow(InlineFunctions_DisallowInliningReason_NOT_CANDIDATE_FUNCTION)
			s.candidateHits++
			if inlineFnCacheVerify {
				verifyCachedCandidateAnalysis(f, target, chunk, debugLogging, &res)
			}
			return res
		}
		if prev.hasFlag(fnCacheShallowLocalClosureDisallowed) {
			res.addDisallow(InlineFunctions_DisallowInliningReason_HAS_LOCAL_NAMES)
			s.candidateHits++
			if inlineFnCacheVerify {
				verifyCachedCandidateAnalysis(f, target, chunk, debugLogging, &res)
			}
			return res
		}
	}

	// Fast Path 2: Full compositional subtree cache hit.
	if sumEntry := s.cache.summaryOf(prev); subtreeUnchanged && prev.hasFlag(fnCacheShallowReady) && sumEntry != nil && sumEntry.summaryFnNameID == targetNameID && countedLimit == 0 {
		if prev.hasFlag(fnCacheDirectReplace) {
			res.inlineDirectly = true
		}
		if !debugLogging {
			if prev.hasFlag(fnCacheIsGenerator) {
				res.addDisallow(InlineFunctions_DisallowInliningReason_GENERATOR_FUNCTION)
				s.candidateHits++
				if inlineFnCacheVerify {
					verifyCachedCandidateAnalysis(f, target, chunk, debugLogging, &res)
				}
				return res
			}
			if prev.hasFlag(fnCacheIsAsync) {
				res.addDisallow(InlineFunctions_DisallowInliningReason_ASYNC_FUNCTION)
				s.candidateHits++
				if inlineFnCacheVerify {
					verifyCachedCandidateAnalysis(f, target, chunk, debugLogging, &res)
				}
				return res
			}
		}
		if prev.hasFlag(fnCacheHasNonInlinableParam) {
			res.addDisallow(InlineFunctions_DisallowInliningReason_NON_INLINABLE_PARAM)
			if !debugLogging {
				s.candidateHits++
				if inlineFnCacheVerify {
					verifyCachedCandidateAnalysis(f, target, chunk, debugLogging, &res)
				}
				return res
			}
		}
		if res.canInline() || debugLogging {
			if isExportedOrRenameCandidate(f, target) || sumEntry.summary.blocksMinReqs {
				res.addDisallow(InlineFunctions_DisallowInliningReason_NOT_CANDIDATE_FUNCTION)
				if !debugLogging {
					s.candidateHits++
					if inlineFnCacheVerify {
						verifyCachedCandidateAnalysis(f, target, chunk, debugLogging, &res)
					}
					return res
				}
			}
		}
		if res.canInline() {
			sum := &sumEntry.summary
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
						s.candidateHits++
						if inlineFnCacheVerify {
							verifyCachedCandidateAnalysis(f, target, chunk, debugLogging, &res)
						}
						return res
					}
				}
			}
			if res.canInline() && f.enforceMaxSizeAfterInlining && !prev.hasFlag(fnCacheAlwaysInlinable) {
				res.bodyAstSize = int32(min(int(sum.bodyAstSize), f.maxSizeAfterInlining))
			}
		}
		if fnNode.GetGrandparent().IsVar() &&
			(res.canInline() || (debugLogging && InlineFunctions_HasRequireInliningAnnotation(fnNode))) {
			block := target.declaringBlock()
			if block.IsBlock() && !block.GetParent().IsFunction() && s.blockHasLetOrConst(block) {
				res.addDisallow(InlineFunctions_DisallowInliningReason_CAPTURED_VARIABLE_NOT_IN_SCOPE)
				if !debugLogging {
					s.candidateHits++
					if inlineFnCacheVerify {
						verifyCachedCandidateAnalysis(f, target, chunk, debugLogging, &res)
					}
					return res
				}
			}
		}
		if debugLogging {
			if prev.hasFlag(fnCacheIsGenerator) {
				res.addDisallow(InlineFunctions_DisallowInliningReason_GENERATOR_FUNCTION)
			}
			if prev.hasFlag(fnCacheIsAsync) {
				res.addDisallow(InlineFunctions_DisallowInliningReason_ASYNC_FUNCTION)
			}
		}
		if res.canInline() && sumEntry.summary.paramsPositional {
			res.summary = sumEntry.summary
			if sumEntry.sideEffects.computed && s.cache.seCfgGen != 0 && sumEntry.seCfgGen == s.cache.seCfgGen && s.cache.seCfg == s.sideEffectCfg {
				res.summary.sideEffects = sumEntry.sideEffects
				s.sideEffectHits++
			} else {
				res.summary.sideEffects = calleeSideEffectSummary{}
			}
			res.hasSummary = true
		}
		s.candidateHits++
		if inlineFnCacheVerify {
			verifyCachedCandidateAnalysis(f, target, chunk, debugLogging, &res)
		}
		return res
	}

	// Cache Miss: compute fresh intrinsic facts and stage an update.
	s.candidateMisses++
	if !haveAlwaysInlinable {
		alwaysInlinable = InlineFunctions_IsAlwaysInlinable(fnNode)
	}
	directReplace := FunctionInjector_IsDirectCallNodeReplacementPossible(fnNode)
	if directReplace {
		res.inlineDirectly = true
	}
	isGen := fnNode.IsGeneratorFunction()
	isAsync := fnNode.IsAsyncFunction()
	fn := ast.MustFunction(fnNode)
	hasNonInlinableParam := InlineFunctions_HasNonInlinableParam(fn.Params().Node())

	upd := inlineFnCandUpdate{
		fnNode:           fnNode,
		fullAstSize:      countedSize,
		fullAstSizeLimit: countedLimit,
		flags:            makeCandUpdateFlags(alwaysInlinable, directReplace, isGen, isAsync, hasNonInlinableParam, false),
	}

	if !debugLogging {
		if isGen {
			res.addDisallow(InlineFunctions_DisallowInliningReason_GENERATOR_FUNCTION)
			s.candUpdates = append(s.candUpdates, upd)
			return res
		}
		if isAsync {
			res.addDisallow(InlineFunctions_DisallowInliningReason_ASYNC_FUNCTION)
			s.candUpdates = append(s.candUpdates, upd)
			return res
		}
	}
	if hasNonInlinableParam {
		res.addDisallow(InlineFunctions_DisallowInliningReason_NON_INLINABLE_PARAM)
		if !debugLogging {
			s.candUpdates = append(s.candUpdates, upd)
			return res
		}
	}

	var scanned, abortedLocal bool
	if res.canInline() || debugLogging {
		var candOk bool
		candOk, scanned, abortedLocal = s.checkCandidateAndScanSummary(f, target, !debugLogging)
		if scanned && !hasNonInlinableParam &&
			(abortedLocal || (s.scratchSummary.hasInnerClosures && (s.scratchSummary.numParams > 0 || s.scratchSummary.hasShallowDeclarations))) {
			upd.flags |= fnCacheShallowLocalClosureDisallowed
		}
		if scanned && !abortedLocal {
			s.candSummaries = append(s.candSummaries, s.scratchSummary)
			upd.summaryIdx = uint32(len(s.candSummaries))
			upd.summaryFnNameID = targetNameID
		}
		if !candOk {
			res.addDisallow(InlineFunctions_DisallowInliningReason_NOT_CANDIDATE_FUNCTION)
			if !debugLogging {
				s.candUpdates = append(s.candUpdates, upd)
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
					s.candUpdates = append(s.candUpdates, upd)
					return res
				}
			}
		}
		if res.canInline() && f.enforceMaxSizeAfterInlining && !alwaysInlinable {
			res.bodyAstSize = int32(min(int(sum.bodyAstSize), f.maxSizeAfterInlining))
		}
	}

	if fnNode.GetGrandparent().IsVar() &&
		(res.canInline() || (debugLogging && InlineFunctions_HasRequireInliningAnnotation(fnNode))) {
		block := target.declaringBlock()
		if block.IsBlock() && !block.GetParent().IsFunction() && s.blockHasLetOrConst(block) {
			res.addDisallow(InlineFunctions_DisallowInliningReason_CAPTURED_VARIABLE_NOT_IN_SCOPE)
			if !debugLogging {
				s.candUpdates = append(s.candUpdates, upd)
				return res
			}
		}
	}

	if debugLogging {
		if isGen {
			res.addDisallow(InlineFunctions_DisallowInliningReason_GENERATOR_FUNCTION)
		}
		if isAsync {
			res.addDisallow(InlineFunctions_DisallowInliningReason_ASYNC_FUNCTION)
		}
	}

	if res.canInline() && s.scratchSummary.paramsPositional {
		res.summary = s.scratchSummary
		res.hasSummary = true
	}
	s.candUpdates = append(s.candUpdates, upd)
	return res
}

func verifySeededFnInfo(fnNode rhino.Node, info functionInjectorFnInfo) {
	if info.evalState != 0 {
		wantEval := functionInjector_referencesEval(fnNode, fnNode)
		gotEval := info.evalState == 2
		if gotEval != wantEval {
			panic(fmt.Sprintf("inlineFnCache verify evalState mismatch on %v: got %v, want %v", fnNode, gotEval, wantEval))
		}
	}
	if info.innerState != fnInnerUnknown {
		firstInner, count := functionInjector_findInnerFunctions(fnNode, fnNode, rhino.NilNode, 0)
		var wantState uint8
		var wantInner rhino.Node
		switch count {
		case 0:
			wantState = fnInnerNone
		case 1:
			wantState = fnInnerOne
			wantInner = firstInner
		default:
			wantState = fnInnerMultiple
		}
		if info.innerState != wantState || info.inner != wantInner {
			panic(fmt.Sprintf("inlineFnCache verify innerState mismatch on %v: got (%d, %v), want (%d, %v)",
				fnNode, info.innerState, info.inner, wantState, wantInner))
		}
	}
}

func verifyCachedCandidateAnalysis(
	f *InlineFunctions,
	target inlineFnTarget,
	chunk *JSChunk,
	debugLogging bool,
	got *inlineCandidateAnalysis,
) {
	var fresh inlineCandidateScratch
	want := fresh.analyzeCandidateUncached(f, target, chunk, debugLogging)
	if got.canInline() != want.canInline() {
		panic(fmt.Sprintf("inlineFnCache verify canInline mismatch on %q (%v): got %v, want %v",
			target.name, target.fnNode, got.canInline(), want.canInline()))
	}
	if !got.canInline() {
		if got.disallowBeforeSetFn != want.disallowBeforeSetFn {
			panic(fmt.Sprintf("inlineFnCache verify disallowBeforeSetFn mismatch on %q (%v): got %v, want %v",
				target.name, target.fnNode, got.disallowBeforeSetFn, want.disallowBeforeSetFn))
		}
		if debugLogging && (got.disallowCount != want.disallowCount || got.disallowReasons != want.disallowReasons) {
			panic(fmt.Sprintf("inlineFnCache verify disallowReasons mismatch on %q (%v): got %v/%v, want %v/%v",
				target.name, target.fnNode, got.disallowCount, got.disallowReasons, want.disallowCount, want.disallowReasons))
		}
		return
	}
	if got.inlineDirectly != want.inlineDirectly ||
		got.referencesThis != want.referencesThis ||
		got.hasInnerClosures != want.hasInnerClosures ||
		got.aliasCount != want.aliasCount ||
		got.bodyAstSize != want.bodyAstSize ||
		got.hasSummary != want.hasSummary {
		panic(fmt.Sprintf("inlineFnCache verify candidate flags mismatch on %q (%v): got %+v, want %+v",
			target.name, target.fnNode, *got, want))
	}
	if got.hasSummary {
		gs, ws := &got.summary, &want.summary
		if gs.fnNode != ws.fnNode ||
			gs.paramBloom != ws.paramBloom ||
			gs.bodyAstSize != ws.bodyAstSize ||
			gs.shallowReturnCount != ws.shallowReturnCount ||
			gs.numParams != ws.numParams ||
			gs.numModified != ws.numModified ||
			gs.paramsPositional != ws.paramsPositional ||
			gs.isTrivialBody != ws.isTrivialBody ||
			gs.hasInnerClosures != ws.hasInnerClosures ||
			gs.hasShallowDeclarations != ws.hasShallowDeclarations ||
			gs.hasLabels != ws.hasLabels ||
			gs.hasInnerFunctionDecls != ws.hasInnerFunctionDecls ||
			gs.hasVarOrLetDecls != ws.hasVarOrLetDecls ||
			gs.referencesThis != ws.referencesThis ||
			gs.blocksMinReqs != ws.blocksMinReqs ||
			gs.hasParamShadowing != ws.hasParamShadowing ||
			gs.hasSubScopeDecls != ws.hasSubScopeDecls ||
			!slices.Equal(gs.params(), ws.params()) ||
			!slices.Equal(gs.modifiedParamIndices(), ws.modifiedParamIndices()) {
			panic(fmt.Sprintf("inlineFnCache verify summary mismatch on %q (%v)", target.name, target.fnNode))
		}
		if gs.sideEffects.computed {
			argInj := NewFunctionArgumentInjector(f.compiler.GetAstAnalyzer())
			wantSE, ok := argInj.ensureCalleeSideEffects(ws)
			if !ok || gs.sideEffects != wantSE {
				panic(fmt.Sprintf("inlineFnCache verify sideEffects mismatch on %q (%v): got %+v, want %+v (ok=%v)",
					target.name, target.fnNode, gs.sideEffects, wantSE, ok))
			}
		}
	}
}
