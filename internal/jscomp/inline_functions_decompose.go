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

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

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
	if len(f.fns) == 0 {
		return
	}
	cands := make([]*InlineFunctions_FunctionState, 0, len(f.fns))
	for _, functionState := range f.fns {
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
	anonFns map[rhino.Node]*InlineFunctions_FunctionState,
	out *[]*InlineFunctions_FunctionState,
) {
	if node.IsFunction() {
		if fs := anonFns[node]; fs != nil && fs.CanInline() {
			*out = append(*out, fs)
		}
	}
	for c := node.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if !c.GetFirstChild().IsNil() {
			inlineFunctions_collectAnonInlinableStates(c, anonFns, out)
		}
	}
}

func (f *InlineFunctions) resolveInlineConflictsParallel(
	cands []*InlineFunctions_FunctionState,
	numWorkers int,
) {
	hasAnon := len(f.anonFns) > 0
	fnByID := f.namedFns

	type ifConflictScanWorker struct {
		calledIDs inlineIDSet
	}

	scanResults := make([]ifConflictScanResult, len(cands))
	runParallelWorkersIndexed(
		len(cands),
		numWorkers,
		func(_ int) ifConflictScanWorker { return ifConflictScanWorker{} },
		func(w *ifConflictScanWorker, idx int) {
			fs := cands[idx]
			fnNode := fs.functionNode()
			body := ast.MustFunction(fnNode).Body().Node()
			w.calledIDs.clear()
			inlineFunctions_collectCalledFunctionIDs(body, &w.calledIDs)
			if w.calledIDs.isEmpty() {
				return
			}
			var res ifConflictScanResult
			res.hasCalledNames = true
			if hasAnon {
				inlineFunctions_collectAnonInlinableStates(body, f.anonFns, &res.anonStates)
			}
			res.structuralNeedsClone = !inlineFunctions_isIsolatedStmtFn(fnNode) || inlineFunctions_hasParamWithChildren(fnNode)
			for _, id := range w.calledIDs.items {
				if fsCalled := fnByID[id]; fsCalled != nil {
					res.calledStates = append(res.calledStates, fsCalled)
				}
			}
			scanResults[idx] = res
		},
		nil,
	)

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
	runParallelWorkersIndexed(
		len(toClone),
		cloneWorkers,
		func(_ int) struct{} { return struct{}{} },
		func(_ *struct{}, idx int) {
			fs := toClone[idx]
			fs.SetSafeFnNode(fs.functionNode().CloneTree())
		},
		nil,
	)
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
	for p := range ast.MustParamList(node).Params() {
		if p.Token() != rhino.Token_NAME {
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

	fnNode := functionState.functionNode()
	p := f.ensureRunState()
	calledIDs := &p.calledIDs
	calledIDs.clear()
	hasAnonInlinable := inlineFunctions_findCalledFunctionsAndAnon(
		ast.MustFunction(fnNode).Body().Node(), calledIDs, f.anonFns)
	if !calledIDs.isEmpty() {
		needsSafeClone := hasAnonInlinable || !inlineFunctions_isIsolatedStmtFn(fnNode) || inlineFunctions_hasParamWithChildren(fnNode)
		// Prevent the removal of the referenced functions unless they will be immediately inlined.
		for _, id := range calledIDs.items {
			fsCalled := f.namedFns[id]
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

// inlineFunctions_hasParamWithChildren reports whether a parameter of the
// function has children: a default value, a rest or a non-empty pattern.
func inlineFunctions_hasParamWithChildren(fnNode rhino.Node) bool {
	for p := range ast.MustFunction(fnNode).Params().Params() {
		if p.Node().HasChildren() {
			return true
		}
	}
	return false
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
			if ast.RoleOf(anc) != ast.RoleExprStmtExpr {
				return false
			}
		default:
			return false
		}
	}
	return false
}

func inlineFunctions_collectCalledFunctionIDs(node rhino.Node, changed *inlineIDSet) {
	if node.IsName() && InlineFunctions_IsCandidateUsage(node) {
		changed.add(node.GetStringID())
	}

	for c := node.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if c.GetFirstChild().IsNil() {
			if c.IsName() && InlineFunctions_IsCandidateUsage(c) {
				changed.add(c.GetStringID())
			}
			continue
		}
		inlineFunctions_collectCalledFunctionIDs(c, changed)
	}
}

func inlineFunctions_findCalledFunctionsAndAnon(
	node rhino.Node,
	changed *inlineIDSet,
	anonFns map[rhino.Node]*InlineFunctions_FunctionState,
) bool {
	if len(anonFns) == 0 {
		inlineFunctions_collectCalledFunctionIDs(node, changed)
		return false
	}
	return inlineFunctions_findCalledFunctionsWithAnon(node, changed, anonFns)
}

func inlineFunctions_findCalledFunctionsWithAnon(
	node rhino.Node,
	changed *inlineIDSet,
	anonFns map[rhino.Node]*InlineFunctions_FunctionState,
) bool {
	hasAnon := false
	if node.IsName() && InlineFunctions_IsCandidateUsage(node) {
		changed.add(node.GetStringID())
	} else if node.IsFunction() {
		if fs := anonFns[node]; fs != nil && fs.CanInline() {
			hasAnon = true
		}
	}

	for c := node.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if c.GetFirstChild().IsNil() {
			if c.IsName() && InlineFunctions_IsCandidateUsage(c) {
				changed.add(c.GetStringID())
			}
			continue
		}
		if inlineFunctions_findCalledFunctionsWithAnon(c, changed, anonFns) {
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
	if len(f.fns) == 0 {
		return
	}
	var workRefs []*FunctionInjector_Reference
	for _, functionState := range f.fns {
		if functionState.CanInline() {
			for _, ref := range functionState.refs {
				if ref.requiresDecomposition {
					workRefs = append(workRefs, ref)
				}
			}
		}
	}
	if len(workRefs) == 0 {
		return
	}
	if f.runState != nil && !f.injector.hasKnownConstantFunctions() && len(f.runState.constFnIDs) > 0 {
		f.injector.setKnownConstantFunctionIDs(f.runState.constFnIDs)
	}
	for _, ref := range workRefs {
		f.invalidateEnclosingSummaries(ref.callNode)
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

type ifDecompPrefixKind uint8

const (
	ifDecompPrefixTemp ifDecompPrefixKind = iota
	ifDecompPrefixConst
	ifDecompPrefixResult
)

func (k ifDecompPrefixKind) prefix() string {
	switch k {
	case ifDecompPrefixTemp:
		return "JSCompiler_temp" + MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SEPARATOR
	case ifDecompPrefixConst:
		return "JSCompiler_temp_const" + MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SEPARATOR
	default:
		return "JSCompiler_inline_result" + MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SEPARATOR
	}
}

type ifDecompOrdEntry struct {
	suffix    string
	constName string
	tempID    uint32
	constID   uint32
	resultID  uint32
}

var (
	ifDecompOrdOnce    sync.Once
	ifDecompOrdMu      sync.RWMutex
	ifDecompOrdEntries []ifDecompOrdEntry
)

func makeIfDecompOrdEntry(ord int) ifDecompOrdEntry {
	suffix := "\x00ifd_" + strconv.Itoa(ord)
	tempName := "JSCompiler_temp" + MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SEPARATOR + suffix
	constName := "JSCompiler_temp_const" + MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SEPARATOR + suffix
	resultName := "JSCompiler_inline_result" + MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SEPARATOR + suffix
	return ifDecompOrdEntry{
		suffix:    suffix,
		constName: constName,
		tempID:    rhino.RhinoStringPool_InternID(tempName),
		constID:   rhino.RhinoStringPool_InternID(constName),
		resultID:  rhino.RhinoStringPool_InternID(resultName),
	}
}

func ifDecompGetOrdEntry(ord int) ifDecompOrdEntry {
	ifDecompOrdOnce.Do(func() {
		ifDecompOrdEntries = make([]ifDecompOrdEntry, 16)
		for i := range ifDecompOrdEntries {
			ifDecompOrdEntries[i] = makeIfDecompOrdEntry(i)
		}
	})
	ifDecompOrdMu.RLock()
	if ord < len(ifDecompOrdEntries) {
		e := ifDecompOrdEntries[ord]
		ifDecompOrdMu.RUnlock()
		return e
	}
	ifDecompOrdMu.RUnlock()

	ifDecompOrdMu.Lock()
	defer ifDecompOrdMu.Unlock()
	if ord >= len(ifDecompOrdEntries) {
		target := max(ord+1, len(ifDecompOrdEntries)*2)
		ifDecompOrdEntries = slices.Grow(ifDecompOrdEntries, target-len(ifDecompOrdEntries))
		for i := len(ifDecompOrdEntries); i < target; i++ {
			ifDecompOrdEntries = append(ifDecompOrdEntries, makeIfDecompOrdEntry(i))
		}
	}
	return ifDecompOrdEntries[ord]
}

type ifDecompSlot struct {
	refIdx int32
	seqIdx int32
	ord    int32
	kind   ifDecompPrefixKind
}

type ifDecompPlaceholder struct {
	node   rhino.Node
	ord    int32
	seqIdx int32
	kind   ifDecompPrefixKind
}

type ifDecompScriptGroup struct {
	script     rhino.Node
	refIndices []int
	mainOnly   bool
}

func collectIfDecompPlaceholdersByID(
	n rhino.Node,
	minID, maxID uint32,
	slotByStrID map[uint32]ifDecompSlot,
	placeholdersByRef [][]ifDecompPlaceholder,
) {
	if n.IsName() {
		id := n.GetStringID()
		if id >= minID && id <= maxID {
			if slot, ok := slotByStrID[id]; ok {
				placeholdersByRef[slot.refIdx] = append(placeholdersByRef[slot.refIdx], ifDecompPlaceholder{
					node:   n,
					ord:    slot.ord,
					seqIdx: slot.seqIdx,
					kind:   slot.kind,
				})
			}
		}
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if c.GetFirstChild().IsNil() {
			if c.IsName() {
				id := c.GetStringID()
				if id >= minID && id <= maxID {
					if slot, ok := slotByStrID[id]; ok {
						placeholdersByRef[slot.refIdx] = append(placeholdersByRef[slot.refIdx], ifDecompPlaceholder{
							node:   c,
							ord:    slot.ord,
							seqIdx: slot.seqIdx,
							kind:   slot.kind,
						})
					}
				}
			}
			continue
		}
		collectIfDecompPlaceholdersByID(c, minID, maxID, slotByStrID, placeholdersByRef)
	}
}

func (f *InlineFunctions) decomposeExpressionsParallel(workRefs []*FunctionInjector_Reference) bool {
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
	var isMainRef []bool
	for i := range groups {
		if groups[i].mainOnly {
			if isMainRef == nil {
				isMainRef = make([]bool, len(workRefs))
			}
			for _, refIdx := range groups[i].refIndices {
				isMainRef[refIdx] = true
			}
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

	type ifDecompWorker struct {
		boundRC         *recordingCompiler
		inj             *FunctionInjector
		curRefIdx       int
		curSeq          int
		scriptOrd       int
		minStrID        uint32
		maxStrID        uint32
		slotByStrID     map[uint32]ifDecompSlot
		touchedBlocks   []rhino.Node
		touchedBlockSet map[rhino.Node]struct{}
	}

	initWorkerInjector := func(w *ifDecompWorker, rc *recordingCompiler) {
		w.boundRC = rc
		placeholderSupplier := func() string {
			s := w.curSeq
			w.curSeq = s + 1
			ord := w.scriptOrd
			w.scriptOrd = ord + 1
			entry := ifDecompGetOrdEntry(ord)
			if w.slotByStrID == nil {
				w.slotByStrID = make(map[uint32]ifDecompSlot, 32)
			}
			w.slotByStrID[entry.tempID] = ifDecompSlot{refIdx: int32(w.curRefIdx), seqIdx: int32(s), ord: int32(ord), kind: ifDecompPrefixTemp}
			w.slotByStrID[entry.constID] = ifDecompSlot{refIdx: int32(w.curRefIdx), seqIdx: int32(s), ord: int32(ord), kind: ifDecompPrefixConst}
			w.slotByStrID[entry.resultID] = ifDecompSlot{refIdx: int32(w.curRefIdx), seqIdx: int32(s), ord: int32(ord), kind: ifDecompPrefixResult}
			lo := min(entry.tempID, entry.constID, entry.resultID)
			hi := max(entry.tempID, entry.constID, entry.resultID)
			if lo < w.minStrID {
				w.minStrID = lo
			}
			if hi > w.maxStrID {
				w.maxStrID = hi
			}
			return entry.suffix
		}
		inj := NewFunctionInjector_Builder(rc).
			SafeNameIdSupplier(placeholderSupplier).
			AllowDecomposition(f.injector.allowDecomposition).
			AssumeStrictThis(f.injector.assumeStrictThis).
			FunctionArgumentInjector(NewFunctionArgumentInjector(f.compiler.GetAstAnalyzer())).
			Build()
		inj.shareKnownConstantFunctionsFrom(f.injector)
		w.inj = inj
	}

	runGroup := func(g *ifDecompScriptGroup, rc *recordingCompiler, w *ifDecompWorker) {
		w.scriptOrd = 0
		w.minStrID = ^uint32(0)
		w.maxStrID = 0
		if w.slotByStrID != nil {
			clear(w.slotByStrID)
		}
		w.touchedBlocks = w.touchedBlocks[:0]
		if w.touchedBlockSet != nil {
			clear(w.touchedBlockSet)
		}
		scanFullScript := false
		for _, refIdx := range g.refIndices {
			var stmtBlock rhino.Node
			if exprRoot := ExpressionDecomposer_FindExpressionRoot(workRefs[refIdx].callNode); !exprRoot.IsNil() {
				stmtBlock = expressionDecomposer_injectionPointFromRoot(exprRoot).GetParent()
			}
			initRecordingCompiler(rc, f.compiler)
			w.curRefIdx = refIdx
			w.curSeq = 0
			w.inj.MaybePrepareCall(workRefs[refIdx])
			if w.curSeq > 0 {
				seqCounts[refIdx] = w.curSeq
				if stmtBlock.IsNil() {
					scanFullScript = true
				} else {
					if w.touchedBlockSet == nil {
						w.touchedBlockSet = make(map[rhino.Node]struct{}, 8)
					}
					if _, exists := w.touchedBlockSet[stmtBlock]; !exists {
						w.touchedBlockSet[stmtBlock] = struct{}{}
						w.touchedBlocks = append(w.touchedBlocks, stmtBlock)
					}
				}
			}
			eventsByRef[refIdx] = rc.takeEvents()
		}
		if w.scriptOrd == 0 {
			return
		}
		if scanFullScript {
			collectIfDecompPlaceholdersByID(g.script, w.minStrID, w.maxStrID, w.slotByStrID, placeholdersByRef)
			return
		}
		roots := w.touchedBlocks[:0]
		for _, b := range w.touchedBlocks {
			if b == g.script {
				collectIfDecompPlaceholdersByID(g.script, w.minStrID, w.maxStrID, w.slotByStrID, placeholdersByRef)
				return
			}
			covered := false
			attached := false
			for anc := b.GetParent(); !anc.IsNil(); anc = anc.GetParent() {
				if _, exists := w.touchedBlockSet[anc]; exists {
					covered = true
				}
				if anc == g.script {
					attached = true
					break
				}
			}
			if !attached {
				collectIfDecompPlaceholdersByID(g.script, w.minStrID, w.maxStrID, w.slotByStrID, placeholdersByRef)
				return
			}
			if !covered {
				roots = append(roots, b)
			}
		}
		for _, b := range roots {
			collectIfDecompPlaceholdersByID(b, w.minStrID, w.maxStrID, w.slotByStrID, placeholdersByRef)
		}
	}

	runParallelWork(f.compiler, numWorkers, len(parallelGroupIndices), func(_ int) (func(int, *recordingCompiler), func()) {
		var w ifDecompWorker
		runItem := func(idx int, rc *recordingCompiler) {
			if w.inj == nil || w.boundRC != rc {
				initWorkerInjector(&w, rc)
			}
			g := &groups[parallelGroupIndices[idx]]
			runGroup(g, rc, &w)
		}
		return runItem, nil
	}, false)

	var realIDs []string
	var realStrIDs []strpool.ID
	var renamedSeq []bool
	for refIdx, ref := range workRefs {
		if isMainRef != nil && isMainRef[refIdx] {
			f.injector.MaybePrepareCall(ref)
			continue
		}
		cnt := seqCounts[refIdx]
		if cnt > 0 {
			realIDs = slices.Grow(realIDs[:0], cnt)[:cnt]
			for k := 0; k < cnt; k++ {
				realIDs[k] = f.injector.safeNameIdSupplier()
			}
			realStrIDs = slices.Grow(realStrIDs[:0], cnt)[:cnt]
			clear(realStrIDs)
			renamedSeq = slices.Grow(renamedSeq[:0], cnt)[:cnt]
			clear(renamedSeq)
			s := ref.lookupScope()
			var hs *Scope
			if s != nil {
				hs = s.GetClosestHoistScope()
			}
			for _, pn := range placeholdersByRef[refIdx] {
				seq := int(pn.seqIdx)
				if seq >= cnt {
					continue
				}
				targetID := realStrIDs[seq]
				if targetID == 0 {
					newName := pn.kind.prefix() + realIDs[seq]
					targetID = strpool.Intern(newName)
					realStrIDs[seq] = targetID
					if pn.kind == ifDecompPrefixConst && !renamedSeq[seq] {
						renamedSeq[seq] = true
						oldName := ifDecompGetOrdEntry(int(pn.ord)).constName
						if s != nil {
							s.renameOwnVar(oldName, newName)
							if hs != nil && hs != s {
								hs.renameOwnVar(oldName, newName)
							}
						}
					}
				}
				setNodeStringID(pn.node, targetID)
			}
		}
		if evs := eventsByRef[refIdx]; len(evs) > 0 {
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
	for _, functionState := range f.fns {
		if functionState.CanRemove() {
			target := functionState.target
			jutil.CheckState(functionState.CanInline())
			jutil.CheckState(target.kind != fnKindNone)
			if f.runState != nil && len(f.runState.movedArgs) > 0 {
				inMovedArg := false
				for anc := target.fnNode; !anc.IsNil(); anc = anc.GetParent() {
					if _, ok := f.runState.movedArgs[anc]; ok {
						inMovedArg = true
						break
					}
				}
				if inMovedArg {
					f.verifyAllReferencesInlined(functionState)
					if target.kind != fnKindExpr {
						f.compiler.GetChangeTracker().notifyChangeHandlers()
					}
					continue
				}
			}
			f.verifyAllReferencesInlined(functionState)
			target.remove(f.compiler)
			NodeUtil_MarkFunctionsDeleted(target.fnNode, f.compiler)
		}
	}
}

func (f *InlineFunctions) verifyAllReferencesInlined(functionState *InlineFunctions_FunctionState) {
	for _, ref := range functionState.refs {
		if !ref.inlined {
			f.VerifyAllReferencesInlined(f.functionStateName(functionState), functionState)
			return
		}
	}
}

// VerifyAllReferencesInlined checks to verify that expression rewriting
// didn't make a call inaccessible.
//
// Java: InlineFunctions.verifyAllReferencesInlined(String, FunctionState)
func (f *InlineFunctions) VerifyAllReferencesInlined(name string, functionState *InlineFunctions_FunctionState) {
	for _, ref := range functionState.refs {
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
