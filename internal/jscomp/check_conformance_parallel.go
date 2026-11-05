/*
 * Copyright 2014 The Closure Compiler Authors.
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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type banUnknownThisGlobalHit struct {
	flatIdx  int
	eventIdx int
	node     rhino.Node
}

type checkConformanceWorkerState struct {
	rc                  recordingCompiler
	sc                  *SyntacticScopeCreator
	cc                  *CheckConformance
	banUnknownThisRules []*conformanceBanUnknownThisRule
	curIdx              int
}

func (b *conformanceAbstractRule) rebindForWorker(
	rc *recordingCompiler,
	check func(*NodeTraversal, rhino.Node) *ConformanceRules_ConformanceResult,
) {
	b.compiler = rc
	b.checkConformance = check
	b.cachedRawPath = ""
	b.cachedRawTsPath = ""
	b.cachedNormPath = ""
	b.cachedNormTsPath = ""
	b.hasCachedPath = false
	b.hasCachedRawTsPath = false
	b.hasCachedNormTsPath = false
}

func canCloneConformanceRule(rule CheckConformance_Rule) (isBanUnknownThis bool, ok bool) {
	switch rule.(type) {
	case *conformanceNoOpRule,
		*conformanceInferredConstRule,
		*conformanceBannedDependencyRule,
		*conformanceBannedDependencyRegexRule,
		*conformanceBannedNameRule,
		*ConformanceRules_BannedProperty,
		*conformanceBannedStringRegexRule,
		*conformanceRestrictedNameCallRule,
		*conformanceRestrictedMethodCallRule,
		*conformanceRestrictedPropertyWriteRule,
		*conformanceBannedCodePatternRule,
		*conformanceBanForOfRule,
		*conformanceRequireUseStrictRule,
		*conformanceBanThrowOfNonErrorTypesRule,
		*conformanceBanNullDerefRule,
		*conformanceBanUnknownDirectThisPropsRule,
		*conformanceBanUnknownTypedClassPropsRule,
		*conformanceBanUnresolvedTypeRule,
		*conformanceStrictBanUnresolvedTypeRule,
		*conformanceBanGlobalVarsRule,
		*conformanceBannedEnhanceRule,
		*conformanceBannedModsRegexRule,
		*conformanceBanCreateElementRule,
		*conformanceBanCreateDomRule,
		*conformanceBanElementSetAttributeRule,
		*conformanceBanSettingAttributesRule,
		*conformanceBanExecCommandRule,
		*conformanceBanStaticThisRule:
		return false, true
	case *conformanceBanUnknownThisRule:
		return true, true
	default:
		return false, false
	}
}

func prepareConformanceParallelRules(flatRules []flatConformanceRule) {
	for i := range flatRules {
		if r, ok := flatRules[i].rule.(*conformanceBanElementSetAttributeRule); ok {
			if gns := r.getGlobalNamespace(); gns != nil {
				gns.EnsureGenerated()
			}
		}
	}
}

func cloneConformanceRuleForWorker(
	flatIdx int,
	rule CheckConformance_Rule,
	rc *recordingCompiler,
	onGlobalHit func(flatIdx int, n rhino.Node),
) (CheckConformance_Rule, *ConformanceRules_AbstractRule, *conformanceBanUnknownThisRule) {
	switch r := rule.(type) {
	case *conformanceNoOpRule:
		clone := *r
		clone.rebindForWorker(rc, clone.checkNode)
		return &clone, &clone.conformanceAbstractRule, nil
	case *conformanceInferredConstRule:
		clone := *r
		clone.rebindForWorker(rc, clone.checkNode)
		return &clone, &clone.conformanceAbstractRule, nil
	case *conformanceBannedDependencyRule:
		clone := *r
		clone.rebindForWorker(rc, clone.checkNode)
		return &clone, &clone.conformanceAbstractRule, nil
	case *conformanceBannedDependencyRegexRule:
		clone := *r
		clone.rebindForWorker(rc, clone.checkNode)
		return &clone, &clone.conformanceAbstractRule, nil
	case *conformanceBannedNameRule:
		clone := *r
		clone.rebindForWorker(rc, clone.checkNode)
		return &clone, &clone.conformanceAbstractRule, nil
	case *ConformanceRules_BannedProperty:
		clone := *r
		clone.rebindForWorker(rc, clone.checkNode)
		return &clone, &clone.conformanceAbstractRule, nil
	case *conformanceBannedStringRegexRule:
		clone := *r
		clone.rebindForWorker(rc, clone.checkNode)
		return &clone, &clone.conformanceAbstractRule, nil
	case *conformanceRestrictedNameCallRule:
		clone := *r
		clone.rebindForWorker(rc, clone.checkNode)
		return &clone, &clone.conformanceAbstractRule, nil
	case *conformanceRestrictedMethodCallRule:
		clone := *r
		clone.rebindForWorker(rc, clone.checkNode)
		return &clone, &clone.conformanceAbstractRule, nil
	case *conformanceRestrictedPropertyWriteRule:
		clone := *r
		clone.rebindForWorker(rc, clone.checkNode)
		return &clone, &clone.conformanceAbstractRule, nil
	case *conformanceBannedCodePatternRule:
		clone := *r
		clone.restrictions = make([]*TemplateAstMatcher, len(r.restrictions))
		for i, m := range r.restrictions {
			clone.restrictions[i] = m.cloneForWorker()
		}
		clone.rebindForWorker(rc, clone.checkNode)
		return &clone, &clone.conformanceAbstractRule, nil
	case *conformanceBanForOfRule:
		clone := *r
		clone.rebindForWorker(rc, clone.checkNode)
		return &clone, &clone.conformanceAbstractRule, nil
	case *conformanceRequireUseStrictRule:
		clone := *r
		clone.rebindForWorker(rc, clone.checkNode)
		return &clone, &clone.conformanceAbstractRule, nil
	case *conformanceBanThrowOfNonErrorTypesRule:
		clone := *r
		clone.rebindForWorker(rc, clone.checkNode)
		return &clone, &clone.conformanceAbstractRule, nil
	case *conformanceBanNullDerefRule:
		clone := *r
		clone.rebindForWorker(rc, clone.checkNode)
		return &clone, &clone.conformanceAbstractRule, nil
	case *conformanceBanUnknownThisRule:
		clone := *r
		clone.reports = make(map[rhino.Node]struct{})
		if onGlobalHit != nil {
			idx := flatIdx
			clone.onGlobalScopeHit = func(n rhino.Node) {
				onGlobalHit(idx, n)
			}
		}
		clone.rebindForWorker(rc, clone.checkNode)
		return &clone, &clone.conformanceAbstractRule, &clone
	case *conformanceBanUnknownDirectThisPropsRule:
		clone := *r
		clone.rebindForWorker(rc, clone.checkNode)
		return &clone, &clone.conformanceAbstractRule, nil
	case *conformanceBanUnknownTypedClassPropsRule:
		clone := *r
		clone.rebindForWorker(rc, clone.checkNode)
		return &clone, &clone.conformanceAbstractRule, nil
	case *conformanceBanUnresolvedTypeRule:
		clone := *r
		clone.rebindForWorker(rc, clone.checkNode)
		return &clone, &clone.conformanceAbstractRule, nil
	case *conformanceStrictBanUnresolvedTypeRule:
		clone := *r
		clone.rebindForWorker(rc, clone.checkNode)
		return &clone, &clone.conformanceAbstractRule, nil
	case *conformanceBanGlobalVarsRule:
		clone := *r
		clone.rebindForWorker(rc, clone.checkNode)
		return &clone, &clone.conformanceAbstractRule, nil
	case *conformanceBannedEnhanceRule:
		clone := *r
		clone.rebindForWorker(rc, clone.checkNode)
		return &clone, &clone.conformanceAbstractRule, nil
	case *conformanceBannedModsRegexRule:
		clone := *r
		clone.rebindForWorker(rc, clone.checkNode)
		return &clone, &clone.conformanceAbstractRule, nil
	case *conformanceBanCreateElementRule:
		clone := *r
		clone.rebindForWorker(rc, clone.checkNode)
		return &clone, &clone.conformanceAbstractRule, nil
	case *conformanceBanCreateDomRule:
		clone := *r
		clone.rebindForWorker(rc, clone.checkNode)
		return &clone, &clone.conformanceAbstractRule, nil
	case *conformanceBanElementSetAttributeRule:
		clone := *r
		gns := r.globalNamespace
		clone.performGlobalNamespaceAnalysis = gns != nil
		ssa := *r.securitySensitiveAttributes
		ssa.globalNamespaceSupplier = func() *GlobalNamespace { return gns }
		clone.securitySensitiveAttributes = &ssa
		clone.rebindForWorker(rc, clone.checkNode)
		return &clone, &clone.conformanceAbstractRule, nil
	case *conformanceBanSettingAttributesRule:
		clone := *r
		clone.rebindForWorker(rc, clone.checkNode)
		return &clone, &clone.conformanceAbstractRule, nil
	case *conformanceBanExecCommandRule:
		clone := *r
		clone.rebindForWorker(rc, clone.checkNode)
		return &clone, &clone.conformanceAbstractRule, nil
	case *conformanceBanStaticThisRule:
		clone := *r
		clone.rebindForWorker(rc, clone.checkNode)
		return &clone, &clone.conformanceAbstractRule, nil
	default:
		panic(jutil.NewIllegalStateException())
	}
}

func (c *CheckConformance) cloneForWorker(
	rc *recordingCompiler,
	onGlobalHit func(flatIdx int, n rhino.Node),
) (*CheckConformance, []*conformanceBanUnknownThisRule) {
	cc := &CheckConformance{
		compiler:     rc,
		categories:   c.categories,
		activeTokens: c.activeTokens,
		hasAllTokens: c.hasAllTokens,
		flatRules:    make([]flatConformanceRule, len(c.flatRules)),
		useMask64:    c.useMask64,
		tokenMask0:   c.tokenMask0,
		tokenMasks:   c.tokenMasks,
	}
	if len(c.scriptSkipMask) != 0 {
		cc.scriptSkipMask = make([]uint64, len(c.scriptSkipMask))
	}
	var banUnknownThis []*conformanceBanUnknownThisRule
	for i := range c.flatRules {
		orig := &c.flatRules[i]
		clonedRule, clonedAbs, unknownThisRule := cloneConformanceRuleForWorker(i, orig.rule, rc, onGlobalHit)
		entry := *orig
		entry.rule = clonedRule
		entry.absRule = clonedAbs
		cc.flatRules[i] = entry
		if unknownThisRule != nil {
			banUnknownThis = append(banUnknownThis, unknownThisRule)
		}
	}
	return cc, banUnknownThis
}

func (c *CheckConformance) processParallel(externs rhino.Node, root rhino.Node) bool {
	compiler := c.compiler
	if compiler == nil ||
		compiler.IsDebugLoggingEnabled() ||
		!hasDefaultConformanceReporting(compiler) ||
		CustomBannedNameEvaluator != nil {
		return false
	}
	if compiler.typedScopeCreator != nil &&
		(compiler.typedScopeCreator.stage != TypedScopeCreator_Stage_FROZEN || compiler.typedScopeCreator.hasUnevaluatedRecordEnumAlias) {
		return false
	}
	hasBanUnknownThis := false
	for i := range c.flatRules {
		isUnknownThis, ok := canCloneConformanceRule(c.flatRules[i].rule)
		if !ok {
			return false
		}
		if isUnknownThis {
			hasBanUnknownThis = true
		}
	}

	extCount, ok1 := countRootScripts(externs)
	rootCount, ok2 := countRootScripts(root)
	if !ok1 || !ok2 {
		return false
	}
	scopeRoot := externs.GetParent()
	if scopeRoot.IsNil() || root.GetParent() != scopeRoot {
		return false
	}

	items := make([]scriptWithParent, 0, rootCount)
	skipFiltered := !Platform_IsThreadInterrupted()
	for _, parent := range [2]rhino.Node{externs, root} {
		for s := parent.GetFirstChild(); !s.IsNil(); s = s.GetNext() {
			if skipFiltered {
				if inp := compiler.GetInput(s.GetInputId()); inp != nil {
					if sf := inp.GetSourceFile(); sf != nil {
						if !c.IsScriptOfInterest(sf) || strings.HasSuffix(sf.GetName(), "tsmes.closure.js") {
							continue
						}
					}
				}
			}
			items = append(items, scriptWithParent{parent: parent, script: s})
		}
	}

	if len(items) == 0 && !c.hasAllTokens && extCount+rootCount > 0 {
		return true
	}
	numWorkers := parallelThreadCount(compiler, len(items))
	if numWorkers <= 1 {
		return false
	}

	prepareConformanceParallelRules(c.flatRules)
	prewarmCompilerForParallelReads(compiler)
	restoreReg := enableParallelTypeRegistryReads(compiler)
	defer restoreReg()

	lazyGlobal := newLazyGlobalScope(compiler, scopeRoot)
	eventsByItem := make([][]recordedCompilerEvent, len(items))
	var globalHitsByItem [][]banUnknownThisGlobalHit
	if hasBanUnknownThis {
		globalHitsByItem = make([][]banUnknownThisGlobalHit, len(items))
	}

	panicsByItem, workerPanic := collectParallelWorkersDeferredPins(
		len(items),
		numWorkers,
		func(_ int, w *checkConformanceWorkerState) {
			initRecordingCompiler(&w.rc, compiler)
			w.sc = NewSyntacticScopeCreator(&w.rc)
			var onGlobalHit func(flatIdx int, n rhino.Node)
			if hasBanUnknownThis {
				onGlobalHit = func(flatIdx int, n rhino.Node) {
					evIdx := len(w.rc.parallelRecorder.events)
					globalHitsByItem[w.curIdx] = append(globalHitsByItem[w.curIdx], banUnknownThisGlobalHit{
						flatIdx:  flatIdx,
						eventIdx: evIdx,
						node:     n,
					})
				}
			}
			w.cc, w.banUnknownThisRules = c.cloneForWorker(&w.rc, onGlobalHit)
		},
		func(w *checkConformanceWorkerState, idx int) {
			resetRecordingCompiler(&w.rc, compiler)
			w.curIdx = idx
			for _, r := range w.banUnknownThisRules {
				clear(r.reports)
			}
			item := items[idx]
			wt := acquireNodeTraversalDirect(&w.rc, w.cc, w.sc, false)
			defer releaseNodeTraversal(wt)
			func() {
				defer func() {
					if r := recover(); r != nil {
						wt.ThrowUnexpectedException(jutil.AsThrowable(r))
					}
				}()
				wt.lazyGlobalScope = lazyGlobal
				wt.InitTraversal(scopeRoot)
				wt.currentNode = scopeRoot
				wt.PushScope(scopeRoot)
				wt.TraverseBranch(item.script, item.parent)
				wt.PopScope()
			}()
		},
		func(w *checkConformanceWorkerState, idx int) {
			eventsByItem[idx] = w.rc.takeEvents()
		},
		nil,
		false,
	)

	var rootTraversal *NodeTraversal
	if c.hasAllTokens {
		rootTraversal = acquireNodeTraversalDirect(compiler, c, nil, false)
		defer releaseNodeTraversal(rootTraversal)
		rootTraversal.lazyGlobalScope = lazyGlobal
		rootTraversal.InitTraversal(scopeRoot)
		rootTraversal.PushScope(scopeRoot)
		rootTraversal.currentNode = externs
		c.Visit(rootTraversal, externs, scopeRoot)
	}

	var seenGlobalUnknownThis []bool
	if hasBanUnknownThis {
		seenGlobalUnknownThis = make([]bool, len(c.flatRules))
	}
	for i := range items {
		evs := eventsByItem[i]
		if hasBanUnknownThis && len(globalHitsByItem[i]) != 0 {
			var dropEventIdx map[int]struct{}
			for _, hit := range globalHitsByItem[i] {
				if seenGlobalUnknownThis[hit.flatIdx] {
					if hit.eventIdx < len(evs) &&
						evs[hit.eventIdx].kind == recordedEventReport &&
						evs[hit.eventIdx].err != nil &&
						evs[hit.eventIdx].err.Node() == hit.node {
						if dropEventIdx == nil {
							dropEventIdx = make(map[int]struct{}, len(globalHitsByItem[i]))
						}
						dropEventIdx[hit.eventIdx] = struct{}{}
					}
				} else {
					seenGlobalUnknownThis[hit.flatIdx] = true
				}
			}
			if len(dropEventIdx) != 0 {
				filtered := evs[:0]
				for evIdx, ev := range evs {
					if _, drop := dropEventIdx[evIdx]; !drop {
						filtered = append(filtered, ev)
					}
				}
				evs = filtered
			}
		}
		if len(evs) != 0 {
			replayCompilerEvents(compiler, evs)
		}
		if p := panicsByItem[i]; p != nil {
			rethrowParallelPanic(compiler, p)
		}
	}
	if workerPanic != nil {
		rethrowParallelPanic(compiler, workerPanic)
	}

	if rootTraversal != nil {
		rootTraversal.currentNode = root
		c.Visit(rootTraversal, root, scopeRoot)
		rootTraversal.currentNode = scopeRoot
		c.Visit(rootTraversal, scopeRoot, rhino.NilNode)
		rootTraversal.PopScope()
	}
	return true
}
