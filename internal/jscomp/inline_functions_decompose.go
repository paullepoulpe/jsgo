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
	"strings"
	"sync"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
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

	scanResults := make([]ifConflictScanResult, len(cands))
	var nextIdx atomic.Uint32
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for w := 0; w < numWorkers; w++ {
		go func() {
			defer wg.Done()
			var calledIDs inlineIDSet
			for {
				idx := int(nextIdx.Add(1) - 1)
				if idx >= len(cands) {
					return
				}
				fs := cands[idx]
				fnNode := fs.functionNode()
				body := ast.MustFunction(fnNode).Body().Node()
				calledIDs.clear()
				inlineFunctions_collectCalledFunctionIDs(body, &calledIDs)
				if calledIDs.isEmpty() {
					continue
				}
				var res ifConflictScanResult
				res.hasCalledNames = true
				if hasAnon {
					inlineFunctions_collectAnonInlinableStates(body, f.anonFns, &res.anonStates)
				}
				res.structuralNeedsClone = !inlineFunctions_isIsolatedStmtFn(fnNode) || inlineFunctions_hasParamWithChildren(fnNode)
				for _, id := range calledIDs.items {
					if fsCalled := fnByID[id]; fsCalled != nil {
						res.calledStates = append(res.calledStates, fsCalled)
					}
				}
				scanResults[idx] = res
			}
		}()
	}
	wg.Wait()

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
	if cloneWorkers <= 1 {
		for _, fs := range toClone {
			fs.SetSafeFnNode(fs.functionNode().CloneTree())
		}
		return
	}
	var nextClone atomic.Uint32
	wg.Add(cloneWorkers)
	for w := 0; w < cloneWorkers; w++ {
		go func() {
			defer wg.Done()
			for {
				idx := int(nextClone.Add(1) - 1)
				if idx >= len(toClone) {
					return
				}
				fs := toClone[idx]
				fs.SetSafeFnNode(fs.functionNode().CloneTree())
			}
		}()
	}
	wg.Wait()
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

// FindCalledFunctions: this functions that may be called directly.
//
// Java: InlineFunctions.findCalledFunctions(Node)
func (f *InlineFunctions) FindCalledFunctions(node rhino.Node) *jutil.LinkedHashSet[string] {
	var calledIDs inlineIDSet
	inlineFunctions_collectCalledFunctionIDs(ast.MustFunction(node).Body().Node(), &calledIDs)
	return calledIDs.toLinkedHashSet()
}

// InlineFunctions_FindCalledFunctionsInto: see findCalledFunctions(Node).
//
// Java: InlineFunctions.findCalledFunctions(Node, Set<String>)
func InlineFunctions_FindCalledFunctionsInto(node rhino.Node, changed *jutil.LinkedHashSet[string]) {
	jutil.CheckArgument(changed != nil)
	var calledIDs inlineIDSet
	inlineFunctions_collectCalledFunctionIDs(node, &calledIDs)
	for _, id := range calledIDs.items {
		changed.Add(rhino.RhinoStringPool_StringFromID(id))
	}
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
	if f.runState != nil && len(f.injector.knownConstantFunctions) == 0 && len(f.runState.constFnIDs) > 0 {
		fnNames := make(map[string]bool, len(f.runState.constFnIDs))
		for _, id := range f.runState.constFnIDs {
			fnNames[rhino.RhinoStringPool_StringFromID(id)] = true
		}
		f.injector.SetKnownConstantFunctions(fnNames)
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

const ifDecompPlaceholderMarker = "\x00ifd_"

type ifDecompPlaceholder struct {
	node   rhino.Node
	prefix string
	seqIdx int
}

type ifDecompScriptGroup struct {
	script     rhino.Node
	refIndices []int
	mainOnly   bool
}

func parseIfDecompPlaceholder(s string) (prefix string, refIdx int, seqIdx int, ok bool) {
	n := len(s)
	if n < 8 || s[n-1] != '\x00' {
		return "", 0, 0, false
	}
	pos := strings.Index(s, ifDecompPlaceholderMarker)
	if pos < 0 {
		return "", 0, 0, false
	}
	body := s[pos+len(ifDecompPlaceholderMarker) : n-1]
	sep := strings.IndexByte(body, '_')
	if sep <= 0 || sep+1 >= len(body) {
		return "", 0, 0, false
	}
	ri, err1 := strconv.Atoi(body[:sep])
	si, err2 := strconv.Atoi(body[sep+1:])
	if err1 != nil || err2 != nil || ri < 0 || si < 0 {
		return "", 0, 0, false
	}
	return s[:pos], ri, si, true
}

func collectIfDecompPlaceholders(
	n rhino.Node,
	refIndices []int,
	placeholdersByRef [][]ifDecompPlaceholder,
) {
	if n.IsName() {
		s := n.GetString()
		if len(s) > 0 && s[len(s)-1] == '\x00' {
			if prefix, refIdx, seqIdx, ok := parseIfDecompPlaceholder(s); ok {
				if _, belongs := slices.BinarySearch(refIndices, refIdx); belongs {
					placeholdersByRef[refIdx] = append(placeholdersByRef[refIdx], ifDecompPlaceholder{
						node:   n,
						prefix: prefix,
						seqIdx: seqIdx,
					})
				}
			}
		}
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if c.GetFirstChild().IsNil() {
			if c.IsName() {
				s := c.GetString()
				if len(s) > 0 && s[len(s)-1] == '\x00' {
					if prefix, refIdx, seqIdx, ok := parseIfDecompPlaceholder(s); ok {
						if _, belongs := slices.BinarySearch(refIndices, refIdx); belongs {
							placeholdersByRef[refIdx] = append(placeholdersByRef[refIdx], ifDecompPlaceholder{
								node:   c,
								prefix: prefix,
								seqIdx: seqIdx,
							})
						}
					}
				}
			}
			continue
		}
		collectIfDecompPlaceholders(c, refIndices, placeholdersByRef)
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
	var mainGroupIndices []int
	for i := range groups {
		if groups[i].mainOnly {
			mainGroupIndices = append(mainGroupIndices, i)
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
		boundRC   *recordingCompiler
		inj       *FunctionInjector
		curRefIdx int
		curSeq    int
	}

	initWorkerInjector := func(w *ifDecompWorker, rc *recordingCompiler) {
		w.boundRC = rc
		placeholderSupplier := func() string {
			s := w.curSeq
			w.curSeq = s + 1
			return ifDecompPlaceholderMarker + strconv.Itoa(w.curRefIdx) + "_" + strconv.Itoa(s) + "\x00"
		}
		inj := NewFunctionInjector_Builder(rc).
			SafeNameIdSupplier(placeholderSupplier).
			AllowDecomposition(f.injector.allowDecomposition).
			AssumeStrictThis(f.injector.assumeStrictThis).
			FunctionArgumentInjector(NewFunctionArgumentInjector(f.compiler.GetAstAnalyzer())).
			Build()
		inj.SetKnownConstantFunctions(f.injector.knownConstantFunctions)
		w.inj = inj
	}

	runGroup := func(g *ifDecompScriptGroup, rc *recordingCompiler, w *ifDecompWorker) {
		anyPlaceholder := false
		for _, refIdx := range g.refIndices {
			initRecordingCompiler(rc, f.compiler)
			w.curRefIdx = refIdx
			w.curSeq = 0
			w.inj.MaybePrepareCall(workRefs[refIdx])
			if w.curSeq > 0 {
				seqCounts[refIdx] = w.curSeq
				anyPlaceholder = true
			}
			eventsByRef[refIdx] = rc.takeEvents()
		}
		if anyPlaceholder {
			scanRoot := g.script
			if scanRoot.IsNil() {
				scanRoot = f.compiler.GetJsRoot()
			}
			if !scanRoot.IsNil() {
				collectIfDecompPlaceholders(scanRoot, g.refIndices, placeholdersByRef)
			}
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

	if len(mainGroupIndices) > 0 {
		var mainRC recordingCompiler
		initRecordingCompiler(&mainRC, f.compiler)
		var mainW ifDecompWorker
		initWorkerInjector(&mainW, &mainRC)
		for _, gIdx := range mainGroupIndices {
			runGroup(&groups[gIdx], &mainRC, &mainW)
		}
	}

	var realIDs []string
	var renamedSeq []bool
	for localOrGlobalRefIdx, ref := range workRefs {
		cnt := seqCounts[localOrGlobalRefIdx]
		if cnt > 0 {
			realIDs = slices.Grow(realIDs[:0], cnt)[:cnt]
			for k := 0; k < cnt; k++ {
				realIDs[k] = f.injector.safeNameIdSupplier()
			}
			renamedSeq = slices.Grow(renamedSeq[:0], cnt)[:cnt]
			clear(renamedSeq)
			s := ref.lookupScope()
			var hs *Scope
			if s != nil {
				hs = s.GetClosestHoistScope()
			}
			for _, pn := range placeholdersByRef[localOrGlobalRefIdx] {
				if pn.seqIdx >= cnt {
					continue
				}
				oldName := pn.node.GetString()
				newName := pn.prefix + realIDs[pn.seqIdx]
				pn.node.SetString(newName)
				if !renamedSeq[pn.seqIdx] {
					renamedSeq[pn.seqIdx] = true
					if s != nil {
						s.renameOwnVar(oldName, newName)
						if hs != nil && hs != s {
							hs.renameOwnVar(oldName, newName)
						}
					}
				}
			}
		}
		if evs := eventsByRef[localOrGlobalRefIdx]; len(evs) > 0 {
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
