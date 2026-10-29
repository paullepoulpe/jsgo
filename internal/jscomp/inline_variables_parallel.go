// Port of com/google/javascript/jscomp/InlineVariables.java @ closure-compiler v20261021

/*
 * Copyright 2008 The Closure Compiler Authors.
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

import "github.com/paullepoulpe/jsgo/internal/rhino"

func (iv *InlineVariables) processParallel(externs rhino.Node, root rhino.Node) bool {
	scopeRoot, items, ok := collectExternAndRootScripts(externs, root)
	if !ok {
		return false
	}
	numWorkers := parallelThreadCount(iv.compiler, len(items))
	if numWorkers <= 1 {
		return false
	}

	prewarmCompilerForParallelReads(iv.compiler)
	restoreReg := enableParallelTypeRegistryReads(iv.compiler)
	defer restoreReg()

	mainSC := NewSyntacticScopeCreator(iv.compiler)
	globalScope := mainSC.CreateScopeOfScope(scopeRoot, nil)
	freezeGlobalScopeForParallelReads(globalScope)
	lazyGlobal := &lazyGlobalScope{
		compiler: iv.compiler,
		root:     scopeRoot,
		scope:    globalScope,
	}
	lazyGlobal.once.Do(func() {})
	globalBlock := NewBasicBlock(nil, scopeRoot)
	varFilter := inlineVariables_Mode_varPredicate[iv.mode]

	statesByItem := make([]*rcParallelScriptState, len(items))
	runParallelWork(iv.compiler, numWorkers, len(items), func(_ int) (func(int, *recordingCompiler), func()) {
		var reusableBeh *InlineVariables
		var reusableCollector *ReferenceCollector
		var reusableSC *SyntacticScopeCreator

		runItem := func(idx int, rc *recordingCompiler) {
			if reusableSC == nil || reusableSC.compiler != rc {
				reusableSC = NewSyntacticScopeCreator(rc)
			}
			var beh *InlineVariables
			var col *ReferenceCollector
			if reusableBeh != nil && !reusableBeh.hasGlobalAliasRetries && len(reusableBeh.aliasRetryHandlers) == 0 {
				beh = reusableBeh
				beh.compiler = rc
				col = reusableCollector
				col.compiler = rc
				col.scopeCreator = reusableSC
				col.clearLocalBase()
				col.maybeReclaimLocals()
				if col.collectedHoistedFunctions != nil {
					clear(col.collectedHoistedFunctions)
				}
			} else {
				beh = NewInlineVariables(rc, iv.mode)
				beh.resetBehavior()
				col = NewReferenceCollectorWithVarFilter(rc, beh, reusableSC, varFilter)
				bRef := beh
				cRef := col
				col.enableLocalReclamation(func() bool {
					if len(bRef.aliasRetryHandlers) == 0 && !bRef.hasGlobalAliasRetries {
						bRef.unpinnedRetries = 0
						cRef.clearLocalBase()
						return true
					}
					if bRef.unpinnedRetries == 0 {
						return true
					}
					bRef.unpinnedRetries = 0
					return false
				})
				reusableBeh = beh
				reusableCollector = col
			}

			ps := &rcParallelScriptState{
				globalScope: globalScope,
				globalBlock: globalBlock,
				rc:          rc,
			}
			ps.curSeg = &ps.rootSeg
			col.parallelState = ps
			beh.onGlobalAliasRetry = func(aliasedVar *Var, retryHandler *InlineVariables_InliningBehavior_AliasInlineRetryHandler) {
				ps.flushCompilerEvents()
				ps.curSeg.ops = append(ps.curSeg.ops, rcStreamOp{
					kind:  rcOpGlobalAliasRetry,
					v:     aliasedVar,
					retry: retryHandler,
				})
			}

			item := items[idx]
			tb := col.builderVal()
			tb.traverseScriptUnderRoot(scopeRoot, item.parent, item.script, lazyGlobal)
			ps.flushCompilerEvents()

			col.parallelState = nil
			beh.onGlobalAliasRetry = nil
			if beh.hasGlobalAliasRetries || len(beh.aliasRetryHandlers) > 0 {
				// Keep beh and col slabs alive until globalScope finishes on main goroutine.
				reusableBeh = nil
				reusableCollector = nil
			}
			statesByItem[idx] = ps
		}
		return runItem, nil
	}, false)

	var hoistedByFn map[rhino.Node]*rcStreamSegment
	for _, ps := range statesByItem {
		if len(ps.hoistedSegs) > 0 {
			if hoistedByFn == nil {
				hoistedByFn = make(map[rhino.Node]*rcStreamSegment, len(ps.hoistedSegs)*2)
			}
			for _, hSeg := range ps.hoistedSegs {
				hoistedByFn[hSeg.fnNode] = hSeg
			}
		}
	}

	iv.resetBehavior()
	globalCollector := NewReferenceCollectorWithVarFilter(iv.compiler, iv, mainSC, varFilter)
	var collectedHoisted map[rhino.Node]struct{}
	if len(hoistedByFn) > 0 {
		collectedHoisted = make(map[rhino.Node]struct{}, len(hoistedByFn))
	}

	var replaySeg func(seg *rcStreamSegment)
	replaySeg = func(seg *rcStreamSegment) {
		if seg == nil {
			return
		}
		for i := range seg.ops {
			op := &seg.ops[i]
			switch op.kind {
			case rcOpCompilerEvents:
				replayCompilerEvents(iv.compiler, op.events)
			case rcOpGlobalRef:
				globalCollector.addGlobalReferenceDirect(op.v, op.node, op.block, op.scope, op.inputId)
			case rcOpGlobalAliasRetry:
				op.retry.outer.outer = iv
				iv.putAliasRetryHandler(op.v, op.retry)
			case rcOpJumpHoisted, rcOpDeclHoisted:
				if _, done := collectedHoisted[op.node]; done {
					continue
				}
				if collectedHoisted == nil {
					collectedHoisted = make(map[rhino.Node]struct{})
				}
				collectedHoisted[op.node] = struct{}{}
				globalCollector.addGlobalReferenceDirect(op.v, op.node.GetFirstChild(), globalBlock, op.scope, op.inputId)
				if hoistedByFn != nil {
					replaySeg(hoistedByFn[op.node])
				}
			}
		}
	}

	for _, ps := range statesByItem {
		replaySeg(&ps.rootSeg)
	}

	gb := globalCollector.builderVal()
	globalTraversal := acquireNodeTraversal(&gb)
	defer releaseNodeTraversal(globalTraversal)
	globalTraversal.InitTraversal(globalScope.GetRootNode())
	globalTraversal.currentNode = globalScope.GetRootNode()
	globalTraversal.PushScopeOfScopeWithQuietly(globalScope, true)
	iv.AfterExitScope(globalTraversal, globalCollector)
	return true
}
