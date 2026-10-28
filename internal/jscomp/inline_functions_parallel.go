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
	name                string
	chunk               *JSChunk
	disallowReason      InlineFunctions_DisallowInliningReason
	hasDisallow         bool
	disallowBeforeSetFn bool
	inlineDirectly      bool
	referencesThis      bool
	hasInnerClosures    bool
	namesToAlias        *jutil.LinkedHashSet[string]
	bodyAstSize         int32
	summary             *calleeSummary
}

type ifPhase1Worker struct {
	outer                         *InlineFunctions
	lastCheckedBlock              rhino.Node
	lastCheckedBlockHasLetOrConst bool
	cands                         []ifPhase1Candidate
	scratchSummary                calleeSummary
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
	name := nameNode.GetString()
	cand := ifPhase1Candidate{
		kind:   kind,
		node:   node,
		fnNode: fnNode,
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
			namesToAlias := jutil.NewLinkedHashSet[string]()
			params := s.params()
			for _, idx := range s.modifiedParamIndices() {
				namesToAlias.Add(params[idx].nameStr())
			}
			cand.inlineDirectly = false
			cand.namesToAlias = namesToAlias
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
		paramSet := FunctionArgumentInjector_GetFunctionParameterSet(fnNode)
		if !paramSet.IsEmpty() {
			namesToAlias := FunctionArgumentInjector_FindModifiedParametersIn(body, paramSet, false)
			if !namesToAlias.IsEmpty() {
				cand.inlineDirectly = false
				cand.namesToAlias = namesToAlias
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

	callsSeen := 0
	for idx := range scripts {
		if p := panicsByScript[idx]; p != nil {
			panic(p)
		}
		for i := range candsByScript[idx] {
			cand := &candsByScript[idx][i]
			if cand.kind == ifCandAnonExpr {
				target := newExprFnTarget(cand.fnNode, callsSeen)
				callsSeen++
				f.MaybeAddFunction(target, cand.chunk)
				f.anonFns.Put(cand.fnNode, target.name)
				continue
			}
			fs := f.GetOrCreateFunctionState(cand.name)
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
						if cand.namesToAlias != nil && !cand.namesToAlias.IsEmpty() {
							fs.InlineDirectly(false)
							fs.SetNamesToAlias(cand.namesToAlias)
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
	result := r.injector.canInlineReferenceWithSummary(
		&r.scratchRef,
		functionState.functionNode(),
		functionState.GetNamesToAlias(),
		functionState.GetReferencesThis(),
		functionState.HasInnerClosures(),
		functionState.readOnlySummary())
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

	name := n.GetString()
	if name == "" {
		return
	}
	if InlineFunctions_IsCandidateUsage(n) {
		return
	}
	functionState := r.visitor.functionMap.Get(name)
	if functionState == nil {
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
		for fs := range f.fns.Values() {
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
				cb = &ifPhase2RefCallback{outer: f}
				cb.visitor = inlineCallVisitor{
					functionMap:     f.fns,
					anonFunctionMap: f.anonFns,
					onCallSite:      cb.VisitCallSite,
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
	for key, functionState := range f.fns.All() {
		if !functionState.HasExistingFunctionDefinition() || !functionState.CanInline() {
			f.fns.Remove(key)
		}
	}
}

// TrimCandidatesUsingOnCost removes entries from the list of candidates that
// can't be inlined.
//
// Java: InlineFunctions.trimCandidatesUsingOnCost()
func (f *InlineFunctions) TrimCandidatesUsingOnCost() {
	for key, functionState := range f.fns.All() {
		if functionState.HasReferences() {
			// Only inline function if it decreases the code size.
			lowersCost := f.MinimizeCost(functionState)
			if !lowersCost {
				// It shouldn't be inlined; remove it from the list.
				f.fns.Remove(key)
			}
		} else if !functionState.CanRemove() {
			// Don't bother tracking functions without references that can't be
			// removed.
			f.fns.Remove(key)
		}
	}
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
	p := f.ensureRunState()
	p.refsScratch = p.refsScratch[:0]
	if functionState.references != nil {
		for ref := range functionState.references.Values() {
			p.refsScratch = append(p.refsScratch, ref)
		}
	}
	return f.injector.inliningLowersCostWithSummary(
		functionState.GetChunk(),
		functionState.functionNode(),
		p.refsScratch,
		functionState.GetNamesToAlias(),
		functionState.CanRemove(),
		functionState.GetReferencesThis(),
		functionState.ensureSummary())
}

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
	if f.fns == nil || f.fns.IsEmpty() {
		return
	}
	cands := make([]*InlineFunctions_FunctionState, 0, f.fns.Size())
	for functionState := range f.fns.Values() {
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
	anonFns *jutil.LinkedHashMap[rhino.Node, string],
	fns *jutil.LinkedHashMap[string, *InlineFunctions_FunctionState],
	out *[]*InlineFunctions_FunctionState,
) {
	if node.IsFunction() {
		if anonName, ok := anonFns.Lookup(node); ok {
			if fs := fns.Get(anonName); fs != nil && fs.CanInline() {
				*out = append(*out, fs)
			}
		}
	}
	for c := node.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if !c.GetFirstChild().IsNil() {
			inlineFunctions_collectAnonInlinableStates(c, anonFns, fns, out)
		}
	}
}

func (f *InlineFunctions) resolveInlineConflictsParallel(
	cands []*InlineFunctions_FunctionState,
	numWorkers int,
) {
	hasAnon := f.anonFns != nil && !f.anonFns.IsEmpty()
	// Ensure LinkedHashMap lookup tables are materialized before concurrent reads.
	_ = f.fns.Get("")
	if hasAnon {
		_, _ = f.anonFns.Lookup(rhino.NilNode)
	}

	scanResults := make([]ifConflictScanResult, len(cands))
	var nextIdx atomic.Uint32
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for w := 0; w < numWorkers; w++ {
		go func() {
			defer wg.Done()
			names := jutil.NewLinkedHashSet[string]()
			for {
				idx := int(nextIdx.Add(1) - 1)
				if idx >= len(cands) {
					return
				}
				fs := cands[idx]
				fnNode := fs.functionNode()
				body := ast.MustFunction(fnNode).Body().Node()
				names.Clear()
				inlineFunctions_findCalledFunctions(body, names)
				if names.IsEmpty() {
					continue
				}
				var res ifConflictScanResult
				res.hasCalledNames = true
				if hasAnon {
					inlineFunctions_collectAnonInlinableStates(body, f.anonFns, f.fns, &res.anonStates)
				}
				res.structuralNeedsClone = !inlineFunctions_isIsolatedStmtFn(fnNode) || inlineFunctions_hasParamWithChildren(fnNode)
				for name := range names.All() {
					if fsCalled := f.fns.Get(name); fsCalled != nil {
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
	names := p.calledNames
	names.Clear()
	hasAnonInlinable := inlineFunctions_findCalledFunctionsAndAnon(
		ast.MustFunction(fnNode).Body().Node(), names, f.anonFns, f.fns)
	if !names.IsEmpty() {
		needsSafeClone := hasAnonInlinable || !inlineFunctions_isIsolatedStmtFn(fnNode) || inlineFunctions_hasParamWithChildren(fnNode)
		// Prevent the removal of the referenced functions unless they will be immediately inlined.
		for name := range names.All() {
			fsCalled := f.fns.Get(name)
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
	changed := jutil.NewLinkedHashSet[string]()
	inlineFunctions_findCalledFunctions(ast.MustFunction(node).Body().Node(), changed)
	return changed
}

// InlineFunctions_FindCalledFunctionsInto: see findCalledFunctions(Node).
//
// Java: InlineFunctions.findCalledFunctions(Node, Set<String>)
func InlineFunctions_FindCalledFunctionsInto(node rhino.Node, changed *jutil.LinkedHashSet[string]) {
	jutil.CheckArgument(changed != nil)
	inlineFunctions_findCalledFunctions(node, changed)
}

func inlineFunctions_findCalledFunctions(node rhino.Node, changed *jutil.LinkedHashSet[string]) {
	if node.IsName() && InlineFunctions_IsCandidateUsage(node) {
		changed.Add(node.GetString())
	}

	for c := node.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if c.GetFirstChild().IsNil() {
			if c.IsName() && InlineFunctions_IsCandidateUsage(c) {
				changed.Add(c.GetString())
			}
			continue
		}
		inlineFunctions_findCalledFunctions(c, changed)
	}
}

func inlineFunctions_findCalledFunctionsAndAnon(
	node rhino.Node,
	changed *jutil.LinkedHashSet[string],
	anonFns *jutil.LinkedHashMap[rhino.Node, string],
	fns *jutil.LinkedHashMap[string, *InlineFunctions_FunctionState],
) bool {
	if anonFns == nil || anonFns.IsEmpty() {
		inlineFunctions_findCalledFunctions(node, changed)
		return false
	}
	return inlineFunctions_findCalledFunctionsWithAnon(node, changed, anonFns, fns)
}

func inlineFunctions_findCalledFunctionsWithAnon(
	node rhino.Node,
	changed *jutil.LinkedHashSet[string],
	anonFns *jutil.LinkedHashMap[rhino.Node, string],
	fns *jutil.LinkedHashMap[string, *InlineFunctions_FunctionState],
) bool {
	hasAnon := false
	if node.IsName() && InlineFunctions_IsCandidateUsage(node) {
		changed.Add(node.GetString())
	} else if node.IsFunction() {
		if anonName, ok := anonFns.Lookup(node); ok {
			if fs := fns.Get(anonName); fs != nil && fs.CanInline() {
				hasAnon = true
			}
		}
	}

	for c := node.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if c.GetFirstChild().IsNil() {
			if c.IsName() && InlineFunctions_IsCandidateUsage(c) {
				changed.Add(c.GetString())
			}
			continue
		}
		if inlineFunctions_findCalledFunctionsWithAnon(c, changed, anonFns, fns) {
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
	if f.fns == nil || f.fns.IsEmpty() {
		return
	}
	var workRefs []*FunctionInjector_Reference
	for functionState := range f.fns.Values() {
		if functionState.CanInline() && functionState.references != nil {
			for ref := range functionState.references.Values() {
				if ref.requiresDecomposition {
					workRefs = append(workRefs, ref)
				}
			}
		}
	}
	if len(workRefs) == 0 {
		return
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

	runGroup := func(g *ifDecompScriptGroup, rc *recordingCompiler, workerInj *FunctionInjector, curRefIdx *int, curSeq *int) {
		anyPlaceholder := false
		for _, refIdx := range g.refIndices {
			initRecordingCompiler(rc, f.compiler)
			*curRefIdx = refIdx
			*curSeq = 0
			workerInj.MaybePrepareCall(workRefs[refIdx])
			if *curSeq > 0 {
				seqCounts[refIdx] = *curSeq
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

	newWorkerInjector := func(rc *recordingCompiler) (*FunctionInjector, *int, *int) {
		curRefIdx := new(int)
		curSeq := new(int)
		placeholderSupplier := func() string {
			s := *curSeq
			*curSeq = s + 1
			return ifDecompPlaceholderMarker + strconv.Itoa(*curRefIdx) + "_" + strconv.Itoa(s) + "\x00"
		}
		inj := NewFunctionInjector_Builder(rc).
			SafeNameIdSupplier(placeholderSupplier).
			AllowDecomposition(f.injector.allowDecomposition).
			AssumeStrictThis(f.injector.assumeStrictThis).
			FunctionArgumentInjector(NewFunctionArgumentInjector(f.compiler.GetAstAnalyzer())).
			Build()
		inj.SetKnownConstantFunctions(f.injector.knownConstantFunctions)
		return inj, curRefIdx, curSeq
	}

	runParallelWork(f.compiler, numWorkers, len(parallelGroupIndices), func(_ int) (func(int, *recordingCompiler), func()) {
		var boundRC *recordingCompiler
		var workerInj *FunctionInjector
		var curRefIdx *int
		var curSeq *int
		runItem := func(idx int, rc *recordingCompiler) {
			if workerInj == nil || boundRC != rc {
				boundRC = rc
				workerInj, curRefIdx, curSeq = newWorkerInjector(rc)
			}
			g := &groups[parallelGroupIndices[idx]]
			runGroup(g, rc, workerInj, curRefIdx, curSeq)
		}
		return runItem, nil
	}, false)

	if len(mainGroupIndices) > 0 {
		var mainRC recordingCompiler
		initRecordingCompiler(&mainRC, f.compiler)
		mainInj, curRefIdx, curSeq := newWorkerInjector(&mainRC)
		for _, gIdx := range mainGroupIndices {
			runGroup(&groups[gIdx], &mainRC, mainInj, curRefIdx, curSeq)
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
	for name, functionState := range f.fns.All() {
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
					f.VerifyAllReferencesInlined(name, functionState)
					if target.kind != fnKindExpr {
						f.compiler.GetChangeTracker().notifyChangeHandlers()
					}
					continue
				}
			}
			f.VerifyAllReferencesInlined(name, functionState)
			target.remove(f.compiler)
			NodeUtil_MarkFunctionsDeleted(target.fnNode, f.compiler)
		}
	}
}

// VerifyAllReferencesInlined checks to verify that expression rewriting
// didn't make a call inaccessible.
//
// Java: InlineFunctions.verifyAllReferencesInlined(String, FunctionState)
func (f *InlineFunctions) VerifyAllReferencesInlined(name string, functionState *InlineFunctions_FunctionState) {
	if functionState.references == nil {
		return
	}
	for ref := range functionState.references.Values() {
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
