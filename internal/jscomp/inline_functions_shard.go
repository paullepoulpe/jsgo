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

const (
	ifUnwrapMinStmts      = 128
	ifMaxUnwrapBlockDepth = 4
)

type ifShardUnitKind uint8

const (
	ifUnitScript ifShardUnitKind = iota + 1
	ifUnitScriptRange
	ifUnitCallChildren
	ifUnitFnParams
	ifUnitFnHeader
	ifUnitFnBodyRange
	ifUnitVisitCall
	ifUnitPostCall
)

type ifShardUnit struct {
	kind              ifShardUnitKind
	script            rhino.Node
	chunk             *JSChunk
	parent            rhino.Node
	first             rhino.Node
	last              rhino.Node
	visitCall         rhino.Node
	visitCallParent   rhino.Node
	parentBlockIdx    int
	unwrappedFnIdx    int
	preCompilerEvents []recordedCompilerEvent
}

type ifUnwrappedBlock struct {
	blkNode        rhino.Node
	parentBlockIdx int
	firstUnitIdx   int
	leafScope      *Scope
	scopeChain     []scopeStackEntry
}

type ifUnwrappedFn struct {
	fnNode           rhino.Node
	fnBody           rhino.Node
	script           rhino.Node
	parentBlockIdx   int
	firstUnitIdx     int
	firstBodyUnitIdx int
	unitIndices      []int
	paramScopeEntry  scopeStackEntry
	blockScopeEntry  scopeStackEntry
	astSize          int
}

type ifUnwrappedCall struct {
	call              rhino.Node
	callParent        rhino.Node
	script            rhino.Node
	chunk             *JSChunk
	target            inlineFnTarget
	calleeIsUnwrapped bool
	calleeUnitIdx     int
	postCallSlot      int
}

type ifShardPlan struct {
	units           []ifShardUnit
	unwrappedBlocks []ifUnwrappedBlock
	unwrappedFns    []ifUnwrappedFn
	unwrappedCalls  []ifUnwrappedCall
}

func ifStmtBatchSize(count, maxThreads int) int {
	target := (count + maxThreads*8 - 1) / (maxThreads * 8)
	return min(max(target, 32), 64)
}

func ifMatchTopWrapperCall(stmt rhino.Node) (call, callParent rhino.Node, ok bool) {
	if !stmt.IsExprResult() {
		return rhino.NilNode, rhino.NilNode, false
	}
	expr := stmt.GetFirstChild()
	if expr.IsNil() {
		return rhino.NilNode, rhino.NilNode, false
	}
	if expr.IsCall() || expr.IsOptChainCall() {
		return expr, stmt, true
	}
	switch expr.GetToken() {
	case rhino.Token_NOT, rhino.Token_VOID, rhino.Token_POS, rhino.Token_NEG, rhino.Token_BITNOT:
		inner := expr.GetFirstChild()
		if !inner.IsNil() && (inner.IsCall() || inner.IsOptChainCall()) {
			return inner, expr, true
		}
	}
	return rhino.NilNode, rhino.NilNode, false
}

func ifUnwrappableCallChild(
	call, ch rhino.Node,
	minStmts int,
	cache *inlineFunctionsCache,
	isPhase1 bool,
) (fn, fnBody rhino.Node, bodyStmtCount int, ok bool) {
	if ch.IsFunction() {
		fn = ch
	} else if ch == call.GetFirstChild() && NodeUtil_IsFunctionObjectCall(call) {
		obj := ch.GetFirstChild()
		if !obj.IsNil() && obj.IsFunction() {
			fn = obj
		} else {
			return rhino.NilNode, rhino.NilNode, 0, false
		}
	} else {
		return rhino.NilNode, rhino.NilNode, 0, false
	}

	fnBody = ast.MustFunction(fn).Body().Node()
	if !fnBody.IsBlock() || fnBody.IsSyntheticBlock() || NodeUtil_IsBundledGoogModuleScopeRoot(fnBody) {
		return rhino.NilNode, rhino.NilNode, 0, false
	}
	if isPhase1 && cache != nil && !inlineFnCacheVerify {
		if prev, found := cache.lookup(fn); found && !fn.IsDeleted() &&
			prev.ownChangeTime == fn.GetChangeTime() &&
			prev.hasFlag(fnCacheHasDirectInnerFns) && prev.innerLen == 0 {
			return rhino.NilNode, rhino.NilNode, 0, false
		}
	}
	count := 0
	for st := fnBody.GetFirstChild(); !st.IsNil(); st = st.GetNext() {
		count++
	}
	if count < minStmts {
		return rhino.NilNode, rhino.NilNode, 0, false
	}
	return fn, fnBody, count, true
}

func ifCallHasUnwrappableChild(call rhino.Node, minStmts int, cache *inlineFunctionsCache, isPhase1 bool) bool {
	for ch := call.GetFirstChild(); !ch.IsNil(); ch = ch.GetNext() {
		if _, _, _, ok := ifUnwrappableCallChild(call, ch, minStmts, cache, isPhase1); ok {
			return true
		}
	}
	return false
}

func ifIsUnwrappableTopBlock(stmt rhino.Node) bool {
	return stmt.IsBlock() && !NodeUtil_IsBundledGoogModuleScopeRoot(stmt)
}

func ifBlockNeedsUnwrap(blk rhino.Node, cache *inlineFunctionsCache, isPhase1 bool, depth int) bool {
	if depth > ifMaxUnwrapBlockDepth || !ifIsUnwrappableTopBlock(blk) {
		return false
	}
	count := 0
	for st := blk.GetFirstChild(); !st.IsNil(); st = st.GetNext() {
		count++
		if count >= ifUnwrapMinStmts {
			return true
		}
		if call, _, ok := ifMatchTopWrapperCall(st); ok &&
			ifCallHasUnwrappableChild(call, ifUnwrapMinStmts, cache, isPhase1) {
			return true
		}
		if depth < ifMaxUnwrapBlockDepth && ifBlockNeedsUnwrap(st, cache, isPhase1, depth+1) {
			return true
		}
	}
	return false
}

func ifScriptNeedsUnwrap(script rhino.Node, cache *inlineFunctionsCache, isPhase1 bool) (topCount int, needsUnwrap bool) {
	first := script.GetFirstChild()
	if first.IsNil() || first.IsModuleBody() {
		return 0, false
	}
	for stmt := first; !stmt.IsNil(); stmt = stmt.GetNext() {
		topCount++
		if !needsUnwrap {
			if call, _, ok := ifMatchTopWrapperCall(stmt); ok &&
				ifCallHasUnwrappableChild(call, ifUnwrapMinStmts, cache, isPhase1) {
				needsUnwrap = true
			} else if ifBlockNeedsUnwrap(stmt, cache, isPhase1, 1) {
				needsUnwrap = true
			}
		}
	}
	if topCount >= ifUnwrapMinStmts {
		needsUnwrap = true
	}
	return topCount, needsUnwrap
}

func ifNodeHasNonLeafChild(n rhino.Node) bool {
	for ch := n.GetFirstChild(); !ch.IsNil(); ch = ch.GetNext() {
		if !ch.GetFirstChild().IsNil() {
			return true
		}
	}
	return false
}

func makeFrozenScopeStackEntry(s *Scope) scopeStackEntry {
	lookup := s
	if s.scanDone && len(s.varList) == 0 && !s.hasImplicitVars && s.lookupParent != nil {
		lookup = s.lookupParent
	}
	return scopeStackEntry{
		node:         s.rootNode,
		scope:        s,
		untypedScope: s,
		lookupScope:  lookup,
	}
}

func (f *InlineFunctions) buildIfShardPlan(root rhino.Node, isPhase1 bool) (ifShardPlan, int, bool) {
	opts := f.compiler.GetOptions()
	if opts == nil || opts.GetNumParallelThreads() <= 1 {
		return ifShardPlan{}, 0, false
	}
	maxThreads := opts.GetNumParallelThreads()
	scripts, ok := collectRootScripts(root)
	if !ok || len(scripts) == 0 {
		return ifShardPlan{}, 0, false
	}
	if len(scripts) == 1 {
		if _, unwrap := ifScriptNeedsUnwrap(scripts[0], f.activeCache, isPhase1); !unwrap {
			return ifShardPlan{}, 0, false
		}
	}

	var plan ifShardPlan
	plan.units = make([]ifShardUnit, 0, max(len(scripts), 16))

	for _, script := range scripts {
		_, unwrap := ifScriptNeedsUnwrap(script, f.activeCache, isPhase1)
		if !unwrap {
			plan.units = append(plan.units, ifShardUnit{
				kind:   ifUnitScript,
				script: script,
			})
			continue
		}
		f.appendUnwrappedScriptUnits(&plan, script, maxThreads, isPhase1)
	}

	numWorkers := parallelThreadCount(f.compiler, len(plan.units))
	if numWorkers <= 1 {
		return ifShardPlan{}, 0, false
	}
	return plan, numWorkers, true
}

func (f *InlineFunctions) appendUnwrappedScriptUnits(
	plan *ifShardPlan,
	script rhino.Node,
	maxThreads int,
	isPhase1 bool,
) {
	var scriptChunk *JSChunk
	if isPhase1 {
		if id := script.GetInputId(); id != nil {
			if inp := f.compiler.GetInput(id); inp != nil {
				scriptChunk = inp.GetChunk()
			}
		}
	}
	f.appendUnwrappedContainerUnits(plan, script, scriptChunk, script, -1, 0, maxThreads, isPhase1)
}

func (f *InlineFunctions) appendUnwrappedContainerUnits(
	plan *ifShardPlan,
	script rhino.Node,
	scriptChunk *JSChunk,
	container rhino.Node,
	parentBlockIdx int,
	depth int,
	maxThreads int,
	isPhase1 bool,
) {
	var runFirst rhino.Node
	runCount := 0
	flushRun := func() {
		if runCount == 0 {
			return
		}
		batchSz := ifStmtBatchSize(runCount, maxThreads)
		cur := runFirst
		rem := runCount
		for rem > 0 {
			take := min(batchSz, rem)
			bFirst := cur
			bLast := cur
			for i := 1; i < take; i++ {
				bLast = bLast.GetNext()
			}
			cur = bLast.GetNext()
			rem -= take
			plan.units = append(plan.units, ifShardUnit{
				kind:           ifUnitScriptRange,
				script:         script,
				chunk:          scriptChunk,
				parent:         container,
				first:          bFirst,
				last:           bLast,
				parentBlockIdx: parentBlockIdx,
			})
		}
		runFirst = rhino.NilNode
		runCount = 0
	}

	for stmt := container.GetFirstChild(); !stmt.IsNil(); stmt = stmt.GetNext() {
		if call, callParent, isCall := ifMatchTopWrapperCall(stmt); isCall &&
			ifCallHasUnwrappableChild(call, ifUnwrapMinStmts, f.activeCache, isPhase1) {
			flushRun()
			f.appendUnwrappedCallUnits(plan, script, scriptChunk, parentBlockIdx, call, callParent, maxThreads, isPhase1)
			continue
		}
		if depth < ifMaxUnwrapBlockDepth && ifBlockNeedsUnwrap(stmt, f.activeCache, isPhase1, depth+1) {
			flushRun()
			blkIdx := len(plan.unwrappedBlocks)
			plan.unwrappedBlocks = append(plan.unwrappedBlocks, ifUnwrappedBlock{
				blkNode:        stmt,
				parentBlockIdx: parentBlockIdx,
				firstUnitIdx:   len(plan.units),
			})
			f.appendUnwrappedContainerUnits(plan, script, scriptChunk, stmt, blkIdx, depth+1, maxThreads, isPhase1)
			continue
		}
		if runCount == 0 {
			runFirst = stmt
		}
		runCount++
	}
	flushRun()
}

func (f *InlineFunctions) appendUnwrappedCallUnits(
	plan *ifShardPlan,
	script rhino.Node,
	scriptChunk *JSChunk,
	parentBlockIdx int,
	call rhino.Node,
	callParent rhino.Node,
	maxThreads int,
	isPhase1 bool,
) {
	uCall := ifUnwrappedCall{
		call:          call,
		callParent:    callParent,
		script:        script,
		chunk:         scriptChunk,
		calleeUnitIdx: -1,
		postCallSlot:  -1,
	}

	var argFirst, argLast rhino.Node
	flushCallArgs := func() {
		if argFirst.IsNil() {
			return
		}
		uIdx := len(plan.units)
		if argFirst == call.GetFirstChild() {
			uCall.calleeUnitIdx = uIdx
		}
		plan.units = append(plan.units, ifShardUnit{
			kind:           ifUnitCallChildren,
			script:         script,
			chunk:          scriptChunk,
			parent:         call,
			first:          argFirst,
			last:           argLast,
			parentBlockIdx: parentBlockIdx,
		})
		argFirst = rhino.NilNode
		argLast = rhino.NilNode
	}

	for ch := call.GetFirstChild(); !ch.IsNil(); ch = ch.GetNext() {
		fn, fnBody, bodyCount, ok := ifUnwrappableCallChild(call, ch, ifUnwrapMinStmts, f.activeCache, isPhase1)
		if !ok {
			if argFirst.IsNil() {
				argFirst = ch
			}
			argLast = ch
			continue
		}

		flushCallArgs()
		if ch == call.GetFirstChild() {
			uCall.calleeIsUnwrapped = true
		}

		fnIdx := len(plan.unwrappedFns)
		uFn := ifUnwrappedFn{
			fnNode:           fn,
			fnBody:           fnBody,
			script:           script,
			parentBlockIdx:   parentBlockIdx,
			firstUnitIdx:     -1,
			firstBodyUnitIdx: -1,
		}
		fnName := fn.GetFirstChild()
		fnParams := fnName.GetNext()

		if isPhase1 {
			if ifNodeHasNonLeafChild(fnParams) {
				uIdx := len(plan.units)
				uFn.firstUnitIdx = uIdx
				uFn.unitIndices = append(uFn.unitIndices, uIdx)
				plan.units = append(plan.units, ifShardUnit{
					kind:           ifUnitFnParams,
					script:         script,
					chunk:          scriptChunk,
					unwrappedFnIdx: fnIdx,
					first:          fnParams,
				})
			}
		} else {
			if fnName.GetStringID() != strpool.EmptyID || !fnParams.GetFirstChild().IsNil() {
				uIdx := len(plan.units)
				uFn.firstUnitIdx = uIdx
				plan.units = append(plan.units, ifShardUnit{
					kind:           ifUnitFnHeader,
					script:         script,
					unwrappedFnIdx: fnIdx,
				})
			}
		}

		batchSz := ifStmtBatchSize(bodyCount, maxThreads)
		cur := fnBody.GetFirstChild()
		rem := bodyCount
		for rem > 0 {
			take := min(batchSz, rem)
			bFirst := cur
			bLast := cur
			for i := 1; i < take; i++ {
				bLast = bLast.GetNext()
			}
			cur = bLast.GetNext()
			rem -= take
			uIdx := len(plan.units)
			if uFn.firstUnitIdx < 0 {
				uFn.firstUnitIdx = uIdx
			}
			if uFn.firstBodyUnitIdx < 0 {
				uFn.firstBodyUnitIdx = uIdx
			}
			if isPhase1 {
				uFn.unitIndices = append(uFn.unitIndices, uIdx)
			}
			plan.units = append(plan.units, ifShardUnit{
				kind:           ifUnitFnBodyRange,
				script:         script,
				chunk:          scriptChunk,
				unwrappedFnIdx: fnIdx,
				parent:         fnBody,
				first:          bFirst,
				last:           bLast,
			})
		}

		plan.unwrappedFns = append(plan.unwrappedFns, uFn)
	}

	if !argFirst.IsNil() {
		if !isPhase1 {
			uIdx := len(plan.units)
			if argFirst == call.GetFirstChild() {
				uCall.calleeUnitIdx = uIdx
			}
			plan.units = append(plan.units, ifShardUnit{
				kind:            ifUnitCallChildren,
				script:          script,
				parent:          call,
				first:           argFirst,
				last:            argLast,
				visitCall:       call,
				visitCallParent: callParent,
				parentBlockIdx:  parentBlockIdx,
			})
		} else {
			flushCallArgs()
		}
	} else if !isPhase1 {
		plan.units = append(plan.units, ifShardUnit{
			kind:            ifUnitVisitCall,
			script:          script,
			visitCall:       call,
			visitCallParent: callParent,
			parentBlockIdx:  parentBlockIdx,
		})
	}

	if isPhase1 && f.reach.IncludesGlobals() {
		if target, match := inlineFunctions_matchCandidateNode(f.reach, call, callParent, false); match {
			uCall.target = target
			uCall.postCallSlot = len(plan.units)
			plan.units = append(plan.units, ifShardUnit{
				kind:   ifUnitPostCall,
				script: script,
				chunk:  scriptChunk,
			})
			plan.unwrappedCalls = append(plan.unwrappedCalls, uCall)
		}
	}
}

func (w *ifPhase1Worker) runShardUnit(
	root rhino.Node,
	u *ifShardUnit,
	unwrappedFns []ifUnwrappedFn,
) (directInners []rhino.Node, allInnerUnchanged bool) {
	defer func() {
		if r := recover(); r != nil {
			script := NodeUtil_GetEnclosingScript(w.current)
			if script.IsNil() {
				script = u.script
			}
			inlineFunctions_throwWalkException(w.outer.compiler, r, w.current, script)
		}
	}()

	switch u.kind {
	case ifUnitScript:
		inlineFunctions_walkCandidates(w.outer.compiler, u.script, root, nil, false, w.scratch, w.visitWalkNode)
	case ifUnitScriptRange:
		for ch := u.first; ; ch = ch.GetNext() {
			if !ch.GetFirstChild().IsNil() {
				inlineFunctions_walkCandidates(w.outer.compiler, ch, u.parent, u.chunk, false, w.scratch, w.visitWalkNode)
			}
			if ch == u.last {
				break
			}
		}
	case ifUnitCallChildren:
		for ch := u.first; ; ch = ch.GetNext() {
			if !ch.GetFirstChild().IsNil() {
				inlineFunctions_walkCandidates(w.outer.compiler, ch, u.parent, u.chunk, false, w.scratch, w.visitWalkNode)
			}
			if ch == u.last {
				break
			}
		}
	case ifUnitFnParams:
		uFn := &unwrappedFns[u.unwrappedFnIdx]
		if w.scratch.cache != nil {
			w.scratch.innerBuf = w.scratch.innerBuf[:0]
			w.scratch.fnFrames = append(w.scratch.fnFrames[:0], inlineFnWalkFrame{
				fnNode:                    uFn.fnNode,
				innerStart:                0,
				allInnerSubtreesUnchanged: true,
			})
		}
		inlineFunctions_walkCandidates(w.outer.compiler, u.first, uFn.fnNode, u.chunk, true, w.scratch, w.visitWalkNode)
		if w.scratch.cache != nil {
			allInnerUnchanged = w.scratch.fnFrames[0].allInnerSubtreesUnchanged
			if len(w.scratch.innerBuf) > 0 {
				directInners = slices.Clone(w.scratch.innerBuf)
			}
			w.scratch.fnFrames = w.scratch.fnFrames[:0]
			w.scratch.innerBuf = w.scratch.innerBuf[:0]
		}
	case ifUnitFnBodyRange:
		uFn := &unwrappedFns[u.unwrappedFnIdx]
		if w.scratch.cache != nil {
			w.scratch.innerBuf = w.scratch.innerBuf[:0]
			w.scratch.fnFrames = append(w.scratch.fnFrames[:0], inlineFnWalkFrame{
				fnNode:                    uFn.fnNode,
				innerStart:                0,
				allInnerSubtreesUnchanged: true,
			})
		}
		for st := u.first; ; st = st.GetNext() {
			if !st.GetFirstChild().IsNil() {
				inlineFunctions_walkCandidates(w.outer.compiler, st, u.parent, u.chunk, true, w.scratch, w.visitWalkNode)
			}
			if st == u.last {
				break
			}
		}
		if w.scratch.cache != nil {
			allInnerUnchanged = w.scratch.fnFrames[0].allInnerSubtreesUnchanged
			if len(w.scratch.innerBuf) > 0 {
				directInners = slices.Clone(w.scratch.innerBuf)
			}
			w.scratch.fnFrames = w.scratch.fnFrames[:0]
			w.scratch.innerBuf = w.scratch.innerBuf[:0]
		}
	case ifUnitPostCall:
		// Evaluated on main goroutine after unwrapped function stamps are finalized.
	}
	return directInners, allInnerUnchanged
}

func (s *inlineCandidateScratch) finishUnwrappedFunction(
	fnNode rhino.Node,
	directInners []rhino.Node,
	allInnerSubtreesUnchanged bool,
) {
	s.liveFns = append(s.liveFns, fnNode)
	prev, ok := s.cache.lookup(fnNode)
	ownTime := fnNode.GetChangeTime()
	ownUnchanged := ok && !fnNode.IsDeleted() && prev.ownChangeTime == ownTime
	sameInners := ok && prev.hasFlag(fnCacheHasDirectInnerFns) && slices.Equal(directInners, s.cache.directInnersOf(prev))
	numInners := len(directInners)

	subtreeUnchanged := ownUnchanged && sameInners && allInnerSubtreesUnchanged
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
}

func (f *InlineFunctions) prebuildPhase2UnwrappedScopes(
	plan *ifShardPlan,
	globalScope *Scope,
	globalEntry scopeStackEntry,
) {
	var mainRC recordingCompiler
	initRecordingCompiler(&mainRC, f.compiler)
	mainSC := NewSyntacticScopeCreator(&mainRC)
	opts := f.compiler.GetOptions()
	mayContainSyntheticBlocks := opts == nil || opts.HasSyntheticBlockStartMarker()
	baseChain := []scopeStackEntry{globalEntry}

	for i := range plan.unwrappedBlocks {
		uBlk := &plan.unwrappedBlocks[i]
		parentScope := globalScope
		parentChain := baseChain
		if uBlk.parentBlockIdx >= 0 {
			pBlk := &plan.unwrappedBlocks[uBlk.parentBlockIdx]
			parentScope = pBlk.leafScope
			parentChain = pBlk.scopeChain
		}
		if mayContainSyntheticBlocks && uBlk.blkNode.IsSyntheticBlock() {
			uBlk.leafScope = parentScope
			uBlk.scopeChain = parentChain
			continue
		}
		resetRecordingCompiler(&mainRC, f.compiler)
		blkScope := mainSC.CreateScopeOfScope(uBlk.blkNode, parentScope)
		freezeGlobalScopeForParallelReads(blkScope)
		uBlk.leafScope = blkScope
		uBlk.scopeChain = append(slices.Clone(parentChain), makeFrozenScopeStackEntry(blkScope))
		if evs := mainRC.takeEvents(); len(evs) > 0 && uBlk.firstUnitIdx < len(plan.units) {
			u := &plan.units[uBlk.firstUnitIdx]
			u.preCompilerEvents = append(u.preCompilerEvents, evs...)
		}
	}

	for i := range plan.unwrappedFns {
		uFn := &plan.unwrappedFns[i]
		parentScope := globalScope
		if uFn.parentBlockIdx >= 0 {
			parentScope = plan.unwrappedBlocks[uFn.parentBlockIdx].leafScope
		}

		resetRecordingCompiler(&mainRC, f.compiler)
		paramScope := mainSC.CreateScopeOfScope(uFn.fnNode, parentScope)
		freezeGlobalScopeForParallelReads(paramScope)
		uFn.paramScopeEntry = makeFrozenScopeStackEntry(paramScope)
		if evs := mainRC.takeEvents(); len(evs) > 0 && uFn.firstUnitIdx >= 0 {
			u := &plan.units[uFn.firstUnitIdx]
			u.preCompilerEvents = append(u.preCompilerEvents, evs...)
		}

		resetRecordingCompiler(&mainRC, f.compiler)
		blockScope := mainSC.CreateScopeOfScope(uFn.fnBody, paramScope)
		freezeGlobalScopeForParallelReads(blockScope)
		uFn.blockScopeEntry = makeFrozenScopeStackEntry(blockScope)
		if evs := mainRC.takeEvents(); len(evs) > 0 && uFn.firstBodyUnitIdx >= 0 {
			u := &plan.units[uFn.firstBodyUnitIdx]
			u.preCompilerEvents = append(u.preCompilerEvents, evs...)
		}

		if f.enforceMaxSizeAfterInlining {
			uFn.astSize = NodeUtil_CountAstSizeUpToLimit(uFn.fnNode, f.maxSizeAfterInlining)
		}
	}
}

func (f *InlineFunctions) newPhase2WorkerCollector(rc *Compiler) *inlineRefCollector {
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
}

func runPhase2ShardUnit(
	t *NodeTraversal,
	cb *inlineRefCollector,
	root rhino.Node,
	lazyGlobal *lazyGlobalScope,
	globalEntry scopeStackEntry,
	hasGlobalEntry bool,
	u *ifShardUnit,
	unwrappedBlocks []ifUnwrappedBlock,
	unwrappedFns []ifUnwrappedFn,
) {
	if u.kind == ifUnitScript {
		traverseScriptUnderRoot(t, root, root, u.script, lazyGlobal)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			t.ThrowUnexpectedException(jutil.AsThrowable(r))
		}
	}()

	t.lazyGlobalScope = lazyGlobal
	t.currentScript = u.script
	t.ClearScriptState()

	pushOuterScopes := func(parentBlockIdx int) {
		clear(t.scopes)
		if parentBlockIdx >= 0 {
			t.scopes = append(t.scopes[:0], unwrappedBlocks[parentBlockIdx].scopeChain...)
			return
		}
		if hasGlobalEntry {
			t.scopes = append(t.scopes[:0], globalEntry)
		} else {
			t.scopes = append(t.scopes[:0], scopeStackEntry{node: root})
		}
	}

	switch u.kind {
	case ifUnitScriptRange:
		t.currentHoistScopeRoot = rhino.NilNode
		t.currentFunction = rhino.NilNode
		t.currentChangeScope = u.script
		pushOuterScopes(u.parentBlockIdx)
		for ch := u.first; ; ch = ch.GetNext() {
			t.TraverseBranch(ch, u.parent)
			if ch == u.last {
				break
			}
		}
		clear(t.scopes)
		t.scopes = t.scopes[:0]

	case ifUnitCallChildren:
		t.currentHoistScopeRoot = rhino.NilNode
		t.currentFunction = rhino.NilNode
		t.currentChangeScope = u.script
		pushOuterScopes(u.parentBlockIdx)
		for ch := u.first; ; ch = ch.GetNext() {
			t.TraverseBranch(ch, u.parent)
			if ch == u.last {
				break
			}
		}
		if !u.visitCall.IsNil() {
			t.currentNode = u.visitCall
			t.visitNode(u.visitCall, u.visitCallParent)
		}
		clear(t.scopes)
		t.scopes = t.scopes[:0]

	case ifUnitVisitCall:
		t.currentHoistScopeRoot = rhino.NilNode
		t.currentFunction = rhino.NilNode
		t.currentChangeScope = u.script
		pushOuterScopes(u.parentBlockIdx)
		t.currentNode = u.visitCall
		t.visitNode(u.visitCall, u.visitCallParent)
		clear(t.scopes)
		t.scopes = t.scopes[:0]

	case ifUnitFnHeader:
		uFn := &unwrappedFns[u.unwrappedFnIdx]
		cb.containingFnSize.wrapperFn = uFn.fnNode
		cb.containingFnSize.wrapperSize = uFn.astSize
		t.currentHoistScopeRoot = uFn.fnNode
		t.currentFunction = uFn.fnNode
		t.currentChangeScope = uFn.fnNode
		pushOuterScopes(uFn.parentBlockIdx)
		t.scopes = append(t.scopes, uFn.paramScopeEntry)
		fnName := uFn.fnNode.GetFirstChild()
		fnParams := fnName.GetNext()
		t.TraverseBranch(fnName, uFn.fnNode)
		t.TraverseBranch(fnParams, uFn.fnNode)
		clear(t.scopes)
		t.scopes = t.scopes[:0]

	case ifUnitFnBodyRange:
		uFn := &unwrappedFns[u.unwrappedFnIdx]
		cb.containingFnSize.wrapperFn = uFn.fnNode
		cb.containingFnSize.wrapperSize = uFn.astSize
		t.currentHoistScopeRoot = uFn.fnNode
		t.currentFunction = uFn.fnNode
		t.currentChangeScope = uFn.fnNode
		pushOuterScopes(uFn.parentBlockIdx)
		t.scopes = append(t.scopes, uFn.paramScopeEntry, uFn.blockScopeEntry)
		for st := u.first; ; st = st.GetNext() {
			t.TraverseBranch(st, u.parent)
			if st == u.last {
				break
			}
		}
		clear(t.scopes)
		t.scopes = t.scopes[:0]
	}
}
