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

import (
	"slices"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: this pass reads the declarations, calls and property accesses
// around a reference through the typed AST view (internal/ast:
// name.Init(), call.Callee(), assign.Target()...) instead of child
// positions.

// InlineVariables uses the infrastructure provided by ReferenceCollector to
// identify variables that are used in a way that is safe to move, and then
// inline them.
//
// This pass has two "modes." One mode only inlines variables declared as
// constants, for legacy compiler clients. The second mode inlines any
// variable that we can provably inline. Note that the second mode is a
// superset of the first mode. We only support the first mode for
// backwards-compatibility with compiler clients that don't want
// --inline_variables.
//
// The basic structure of this class is as follows.
//
//  1. ReferenceCollector#process is invoked on the AST with an instance of
//     InliningBehavior
//  2. InliningBehavior#afterExitScope gets invoked on each scope in DFS
//     order
//  3. It iterates through the variables defined in that scope:
//  1. For each variable a VarExpert is created. This object is
//     responsible for determining:
//     - Can it be inlined?
//     - If so, how should the inlining be done?
//     - If not, is it safe to inline aliases of the variable? (e.g. `const
//     b = doSomething();` isn't safe to inline due to side-effects, but it
//     should be OK to inline `const aliasB = b;` in many cases.
//  2. The VarExpert creates an InlineVarAnalysis object containing its
//     decisions and possibly a method to invoke to do the inlining.
//  3. The InlineVarAnalysis may indicate that a variable is an alias of
//     another variable. If so, InliningBehavior will wait until the aliased
//     variable has been analyzed, then pass that analysis back to the
//     VarExpert for the aliasing variable, so it can complete its decision.
type InlineVariables struct {
	compiler AbstractCompiler

	mode InlineVariables_Mode
}

// InlineVariables_Mode is the enum InlineVariables.Mode.
type InlineVariables_Mode int32

const (
	// Only inline things explicitly marked as constant.
	InlineVariables_Mode_CONSTANTS_ONLY InlineVariables_Mode = iota
	// Locals only
	InlineVariables_Mode_LOCALS_ONLY
	InlineVariables_Mode_ALL
)

var inlineVariables_Mode_names = [...]string{"CONSTANTS_ONLY", "LOCALS_ONLY", "ALL"}

// inlineVariables_Mode_varPredicate is the enum field varPredicate, indexed
// by ordinal.
var inlineVariables_Mode_varPredicate = [...]func(*Var) bool{
	inlineVariables_isDeclaredOrInferredConst,
	func(v *Var) bool { return v.IsLocal() },
	nil,
}

func inlineVariables_isDeclaredOrInferredConst(v *Var) bool {
	n := v.nameNode
	if n.IsNil() {
		return false
	}
	if n.GetConstantVarFlags()&(rhino.Node_ConstantVarFlags_DECLARED|rhino.Node_ConstantVarFlags_INFERRED) != 0 {
		return true
	}
	return n.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME)
}

// Java: InlineVariables.Mode.name()
func (m InlineVariables_Mode) String() string { return inlineVariables_Mode_names[m] }

// Java: InlineVariables.Mode.ordinal()
func (m InlineVariables_Mode) Ordinal() int { return int(m) }

// Java: InlineVariables(AbstractCompiler, Mode)
func NewInlineVariables(compiler AbstractCompiler, mode InlineVariables_Mode) *InlineVariables {
	return &InlineVariables{compiler: compiler, mode: mode}
}

// Java: InlineVariables.process(Node, Node)
func (iv *InlineVariables) Process(externs rhino.Node, root rhino.Node) {
	if iv.processParallel(externs, root) {
		return
	}
	behavior := NewInlineVariables_InliningBehavior(iv)
	callback := NewReferenceCollectorWithVarFilter(
		iv.compiler,
		behavior,
		NewSyntacticScopeCreator(iv.compiler),
		inlineVariables_Mode_varPredicate[iv.mode])
	callback.enableLocalReclamation(func() bool {
		if len(behavior.aliasRetryHandlers) == 0 {
			behavior.unpinnedRetries = 0
			callback.clearLocalBase()
			return true
		}
		if behavior.unpinnedRetries == 0 {
			return true
		}
		behavior.unpinnedRetries = 0
		return false
	})
	callback.Process(externs, root)
}

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
		var reusableBeh *InlineVariables_InliningBehavior
		var reusableCollector *ReferenceCollector
		var reusableSC *SyntacticScopeCreator

		runItem := func(idx int, rc *recordingCompiler) {
			if reusableSC == nil || reusableSC.compiler != rc {
				reusableSC = NewSyntacticScopeCreator(rc)
			}
			var beh *InlineVariables_InliningBehavior
			var col *ReferenceCollector
			if reusableBeh != nil && !reusableBeh.hasGlobalAliasRetries && len(reusableBeh.aliasRetryHandlers) == 0 {
				beh = reusableBeh
				beh.outer.compiler = rc
				col = reusableCollector
				col.compiler = rc
				col.scopeCreator = reusableSC
				col.clearLocalBase()
				col.maybeReclaimLocals()
				if col.collectedHoistedFunctions != nil {
					clear(col.collectedHoistedFunctions)
				}
			} else {
				workerIV := &InlineVariables{compiler: rc, mode: iv.mode}
				beh = NewInlineVariables_InliningBehavior(workerIV)
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

	mainBehavior := NewInlineVariables_InliningBehavior(iv)
	globalCollector := NewReferenceCollectorWithVarFilter(iv.compiler, mainBehavior, mainSC, varFilter)
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
				mainBehavior.putAliasRetryHandler(op.v, op.retry)
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
	mainBehavior.AfterExitScope(globalTraversal, &globalCollector.referenceMapWrapper)
	return true
}

type inlineVariables_handledState uint8

const (
	handledVarUnseen inlineVariables_handledState = iota
	handledVarUnsafe
	handledVarSafe
	handledVarLazy
	handledVarLazyWellDefined
)

// inlineVariables_handledVar records whether an already-analyzed variable's
// name may safely be inlined into aliases of it, deferring expensive
// well-defined / assigned-once checks until an alias actually asks.
type inlineVariables_handledVar struct {
	referenceInfo *ReferenceCollection
	state         inlineVariables_handledState
}

func (h *inlineVariables_handledVar) isSafeToInlineAliases() bool {
	switch h.state {
	case handledVarSafe:
		return true
	case handledVarLazy:
		if h.referenceInfo.IsWellDefined() && h.referenceInfo.IsAssignedOnceInLifetime() {
			h.state = handledVarSafe
			return true
		}
		h.state = handledVarUnsafe
		return false
	case handledVarLazyWellDefined:
		if h.referenceInfo.IsAssignedOnceInLifetime() {
			h.state = handledVarSafe
			return true
		}
		h.state = handledVarUnsafe
		return false
	default:
		return false
	}
}

// InlineVariables_InliningBehavior builds up information about nodes in each
// scope. When exiting the scope, inspects all variables in that scope, and
// inlines any that we can.
type InlineVariables_InliningBehavior struct {
	outer            *InlineVariables
	codingConvention CodingConvention

	// Records the analyses of variables that have already been handled in
	// the current scope (indexed by scope.varList slot) so aliases of those
	// variables can determine whether they may be inlined.
	currentScope     *Scope
	handledScopeVars []inlineVariables_handledVar

	// Records alias variables that are waiting for the original variables to
	// be handled.
	aliasRetryHandlers    map[*Var][]*InlineVariables_InliningBehavior_AliasInlineRetryHandler
	unpinnedRetries       int
	hasGlobalAliasRetries bool
	onGlobalAliasRetry    func(*Var, *InlineVariables_InliningBehavior_AliasInlineRetryHandler)

	scratchExpert    InlineVariables_InliningBehavior_StandardVarExpert
	retryExpertSlab  jutil.ReuseSlab[InlineVariables_InliningBehavior_StandardVarExpert]
	retryHandlerSlab jutil.ReuseSlab[InlineVariables_InliningBehavior_AliasInlineRetryHandler]
	motionScratch    NodeIterators_LocalVarMotion
}

// Java: InlineVariables.InliningBehavior()
func NewInlineVariables_InliningBehavior(outer *InlineVariables) *InlineVariables_InliningBehavior {
	var cc CodingConvention
	if outer != nil && outer.compiler != nil {
		cc = outer.compiler.GetCodingConvention()
	}
	return &InlineVariables_InliningBehavior{
		outer:            outer,
		codingConvention: cc,
	}
}

func scopeVarSlot(scope *Scope, v *Var) int {
	if idx := v.GetIndex(); uint(idx) < uint(len(scope.varList)) && scope.varList[idx] == v {
		return idx
	}
	return slices.Index(scope.varList, v)
}

func (b *InlineVariables_InliningBehavior) recordHandledVar(v *Var, hv inlineVariables_handledVar) {
	if b.outer.mode == InlineVariables_Mode_CONSTANTS_ONLY {
		return
	}
	handledPtr := &hv
	if scope := b.currentScope; v != nil && scope != nil && v.scope == scope {
		if slot := scopeVarSlot(scope, v); slot >= 0 {
			b.handledScopeVars[slot] = hv
			handledPtr = &b.handledScopeVars[slot]
		}
	}
	// Retry any aliases we saw before handling this Var
	b.RetryAliases(v, handledPtr)
}

func (b *InlineVariables_InliningBehavior) getHandledVar(scope *Scope, aliasedVar *Var) *inlineVariables_handledVar {
	if aliasedVar == nil || aliasedVar.scope != scope || scope != b.currentScope {
		return nil
	}
	if slot := scopeVarSlot(scope, aliasedVar); slot >= 0 && b.handledScopeVars[slot].state != handledVarUnseen {
		return &b.handledScopeVars[slot]
	}
	return nil
}

func (b *InlineVariables_InliningBehavior) initScratchVarExpert(
	v *Var,
	referenceInfo *ReferenceCollection,
	isDeclaredOrInferredConstant bool,
	mayBeAParameterModifiedViaArguments bool,
) *InlineVariables_InliningBehavior_StandardVarExpert {
	b.scratchExpert = InlineVariables_InliningBehavior_StandardVarExpert{
		outer:                               b,
		v:                                   v,
		referenceInfo:                       referenceInfo,
		declaration:                         referenceInfo.getRef(0),
		isDeclaredOrInferredConstant:        isDeclaredOrInferredConstant,
		mayBeAParameterModifiedViaArguments: mayBeAParameterModifiedViaArguments,
	}
	return &b.scratchExpert
}

func (b *InlineVariables_InliningBehavior) allocRetryExpert(
	src *InlineVariables_InliningBehavior_StandardVarExpert,
) *InlineVariables_InliningBehavior_StandardVarExpert {
	e := b.retryExpertSlab.Alloc()
	*e = *src
	e.initValAnalysis.outer = e
	return e
}

func (b *InlineVariables_InliningBehavior) allocRetryHandler(
	v *Var,
	expert *InlineVariables_InliningBehavior_StandardVarExpert,
) *InlineVariables_InliningBehavior_AliasInlineRetryHandler {
	h := b.retryHandlerSlab.Alloc()
	*h = InlineVariables_InliningBehavior_AliasInlineRetryHandler{
		outer:  b,
		v:      v,
		expert: expert,
	}
	return h
}

func (b *InlineVariables_InliningBehavior) putAliasRetryHandler(
	aliasedVar *Var,
	retryHandler *InlineVariables_InliningBehavior_AliasInlineRetryHandler,
) {
	if b.onGlobalAliasRetry != nil && !aliasedVar.IsLocal() {
		b.hasGlobalAliasRetries = true
		b.unpinnedRetries++
		b.onGlobalAliasRetry(aliasedVar, retryHandler)
		return
	}
	if b.aliasRetryHandlers == nil {
		b.aliasRetryHandlers = make(map[*Var][]*InlineVariables_InliningBehavior_AliasInlineRetryHandler, 8)
	}
	b.aliasRetryHandlers[aliasedVar] = append(b.aliasRetryHandlers[aliasedVar], retryHandler)
	b.unpinnedRetries++
}

// Java: InlineVariables.InliningBehavior.afterExitScope(NodeTraversal, ReferenceMap)
func (b *InlineVariables_InliningBehavior) AfterExitScope(t *NodeTraversal, referenceMap ReferenceMap) {
	b.DoInlinesForScope(t, referenceMap)
}

// DoInlinesForScope: for all variables in this scope, see if they are only
// used once. If it looks safe to do so, inline them.
//
// Java: InlineVariables.InliningBehavior.doInlinesForScope(NodeTraversal, ReferenceMap)
func (b *InlineVariables_InliningBehavior) DoInlinesForScope(t *NodeTraversal, referenceMap ReferenceMap) {
	// Any variables we completed in an earlier scope are now out of scope,
	// so we can clear this slice.
	clear(b.handledScopeVars)
	b.handledScopeVars = b.handledScopeVars[:0]
	b.currentScope = nil
	if len(b.aliasRetryHandlers) == 0 && !b.hasGlobalAliasRetries {
		b.retryExpertSlab.Reset(nil)
		b.retryHandlerSlab.Reset(nil)
	}
	if len(t.scopes) > 0 && t.scopes[len(t.scopes)-1].scope == nil {
		return
	}
	scope := t.GetScope()
	if scope.GetVarCount() == 0 {
		return
	}
	isLocalScope := scope.IsLocal()
	if b.outer.mode == InlineVariables_Mode_LOCALS_ONLY && !isLocalScope && len(b.aliasRetryHandlers) == 0 {
		return
	}
	if b.outer.mode != InlineVariables_Mode_CONSTANTS_ONLY {
		b.currentScope = scope
		nVars := len(scope.varList)
		b.handledScopeVars = slices.Grow(b.handledScopeVars, nVars)[:nVars]
	}
	mayBeAParameterModifiedViaArguments :=
		b.VarsInThisScopeMayBeModifiedUsingArguments(scope, referenceMap)
	for i := 0; i < len(scope.varList); i++ {
		v := scope.varList[i]
		if v == nil {
			continue
		}
		if b.outer.mode == InlineVariables_Mode_CONSTANTS_ONLY && !inlineVariables_isDeclaredOrInferredConst(v) {
			continue
		}
		referenceInfo := referenceMap.GetReferences(v)
		isConst, bailoutState, bailedOut := b.checkScopeVarEligibility(v, referenceInfo, isLocalScope)
		if bailedOut {
			b.recordHandledVar(v, inlineVariables_handledVar{state: bailoutState})
			continue
		}
		expert := b.initScratchVarExpert(v, referenceInfo, isConst, mayBeAParameterModifiedViaArguments)
		expert.Analyze()
		if expert.ShouldWaitForAliasedVar() {
			// We need to wait for the aliased variable to be handled.
			aliasedVar := expert.GetAliasedVar()
			if aliasedHandled := b.getHandledVar(scope, aliasedVar); aliasedHandled != nil {
				// We've actually already completed the aliased Var
				b.handleAliasedVarCompletion(v, expert, aliasedVar, aliasedHandled)
			} else {
				// wait for completion of the aliased Var
				if expert == &b.scratchExpert {
					expert = b.allocRetryExpert(&b.scratchExpert)
				}
				retryHandler := b.allocRetryHandler(v, expert)
				b.putAliasRetryHandler(aliasedVar, retryHandler)
			}
		} else {
			b.completeVarExpert(v, expert)
		}
	}
}

func (b *InlineVariables_InliningBehavior) completeVarExpert(
	v *Var,
	expert *InlineVariables_InliningBehavior_StandardVarExpert,
) {
	var hv inlineVariables_handledVar
	if expert.ShouldInline() {
		expert.PerformInline()
		hv = inlineVariables_handledVar{state: handledVarSafe}
	} else {
		hv = expert.negativeHandledVar()
	}
	b.recordHandledVar(v, hv)
}

func (b *InlineVariables_InliningBehavior) handleAliasedVarCompletion(
	v *Var,
	expert *InlineVariables_InliningBehavior_StandardVarExpert,
	aliasedVar *Var,
	aliasedHandled *inlineVariables_handledVar,
) {
	expert.ReanalyzeAfterAliasedVar(aliasedVar, aliasedHandled)
	jutil.CheckStateMsg1(
		!expert.ShouldWaitForAliasedVar(), "expert for %s asked to wait a second time", v)
	b.completeVarExpert(v, expert)
}

// VarsInThisScopeMayBeModifiedUsingArguments returns true for function
// scopes (the whole function, not the body), if the function uses
// `arguments` in some way other than a few that are known not to change the
// values of parameter variables.
//
// TODO(bradfordcsmith): In strict mode `arguments` is a copy of the
// parameters, so modifying it cannot change the parameters. Sloppy mode is
// now so rare, that we should probably just ignore the possibility of
// `arguments` being used to modify a parameter value. We ignore loose mode
// issues or simply don't support them ("with" for instance).
//
// Java: InlineVariables.InliningBehavior.varsInThisScopeMayBeModifiedUsingArguments(Scope, ReferenceMap)
func (b *InlineVariables_InliningBehavior) VarsInThisScopeMayBeModifiedUsingArguments(
	scope *Scope, referenceMap ReferenceMap) bool {
	if scope.IsFunctionScope() && !scope.GetRootNode().IsArrowFunction() {
		var refs *ReferenceCollection
		if w, ok := referenceMap.(*ReferenceCollector_ReferenceMapWrapper); ok {
			refs = w.getReferencesByScopedName(scope.GetRootNode(), "arguments")
		} else {
			refs = referenceMap.GetReferences(scope.GetArgumentsVar())
		}
		if refs != nil && len(refs.references) > 0 {
			for _, ref := range refs.references {
				if !b.IsSafeUseOfArguments(ref.GetNode()) {
					return true
				}
			}
		}
	}
	return false
}

// Java: InlineVariables.InliningBehavior.retryAliases(Var, InlineVarAnalysis)
func (b *InlineVariables_InliningBehavior) RetryAliases(aliasedVar *Var, aliasedHandled *inlineVariables_handledVar) {
	if len(b.aliasRetryHandlers) == 0 {
		return
	}
	handlers, ok := b.aliasRetryHandlers[aliasedVar]
	if !ok {
		return
	}
	delete(b.aliasRetryHandlers, aliasedVar)
	if len(b.aliasRetryHandlers) == 0 {
		b.unpinnedRetries = 0
	} else if b.unpinnedRetries >= len(handlers) {
		b.unpinnedRetries -= len(handlers)
	} else {
		b.unpinnedRetries = 0
	}
	for _, aliasInlineRetryHandler := range handlers {
		aliasInlineRetryHandler.HandleAliasedVarCompletion(aliasedVar, aliasedHandled)
	}
}

// InlineVariables_InliningBehavior_AliasInlineRetryHandler retries inlining
// an alias variable after the original variable has been dealt with.
type InlineVariables_InliningBehavior_AliasInlineRetryHandler struct {
	outer *InlineVariables_InliningBehavior

	v      *Var
	expert *InlineVariables_InliningBehavior_StandardVarExpert
}

// Java: InlineVariables.InliningBehavior.AliasInlineRetryHandler.handleAliasedVarCompletion(Var, InlineVarAnalysis)
func (h *InlineVariables_InliningBehavior_AliasInlineRetryHandler) HandleAliasedVarCompletion(
	aliasedVar *Var, aliasedHandled *inlineVariables_handledVar) {
	h.outer.handleAliasedVarCompletion(h.v, h.expert, aliasedVar, aliasedHandled)
}

type standardVarAnalysisKind uint8

const (
	analysisKindNone standardVarAnalysisKind = iota
	analysisKindNegativeUnsafe
	analysisKindNegativeSafe
	analysisKindDelayedNegative
	analysisKindWaitForAlias
	analysisKindInlineWellDefined
	analysisKindInlineNeverAssigned
	analysisKindInlineSingleRef
)

// InlineVariables_InliningBehavior_StandardVarExpert knows how to analyze a
// variable to determine whether it should be inlined, how to do the
// inlining, and whether it's OK to inline aliases of the variable.
type InlineVariables_InliningBehavior_StandardVarExpert struct {
	outer *InlineVariables_InliningBehavior

	v                                   *Var
	referenceInfo                       *ReferenceCollection
	declaration                         *Reference
	isDeclaredOrInferredConstant        bool
	mayBeAParameterModifiedViaArguments bool

	isNeverAssigned           InlineVariables_InitiallyUnknown[bool]
	isWellDefined             InlineVariables_InitiallyUnknown[bool]
	isAssignedOnceInLifetime  InlineVariables_InitiallyUnknown[bool]
	isReferencedWeakly        InlineVariables_InitiallyUnknown[bool]
	isWellDefinedAssignedOnce InlineVariables_InitiallyUnknown[bool]
	initializationReference   InlineVariables_InitiallyUnknown[*Reference]
	allReferencesAreValid     InlineVariables_InitiallyUnknown[bool]

	initValAnalysis       InlineVariables_InliningBehavior_StandardVarExpert_InitialValueAnalysis
	analysisKind          standardVarAnalysisKind
	analysisAliasedVar    *Var
	analysisInitValue     rhino.Node
	analysisInitRef       *Reference
	analysisSingleReadRef *Reference
}

func (e *InlineVariables_InliningBehavior_StandardVarExpert) ShouldInline() bool {
	return e.analysisKind >= analysisKindInlineWellDefined
}

func (e *InlineVariables_InliningBehavior_StandardVarExpert) ShouldWaitForAliasedVar() bool {
	return e.analysisKind == analysisKindWaitForAlias
}

func (e *InlineVariables_InliningBehavior_StandardVarExpert) GetAliasedVar() *Var {
	if e.analysisKind == analysisKindWaitForAlias {
		return e.analysisAliasedVar
	}
	panic(jutil.NewUnsupportedOperationExceptionMsg("no aliased Var"))
}

func (e *InlineVariables_InliningBehavior_StandardVarExpert) PerformInline() {
	switch e.analysisKind {
	case analysisKindInlineWellDefined:
		e.outer.InlineWellDefinedVariable(e.v, e.analysisInitValue, e.referenceInfo.references)
	case analysisKindInlineNeverAssigned:
		srcLocation := e.declaration.GetNode()
		undefinedNode := NodeUtil_NewUndefinedNode(srcLocation)
		e.outer.InlineWellDefinedVariable(e.v, undefinedNode, e.referenceInfo.references)
	case analysisKindInlineSingleRef:
		e.outer.Inline(e.declaration, e.analysisInitRef, e.analysisSingleReadRef)
	default:
		panic(jutil.NewUnsupportedOperationExceptionMsg("cannot inline"))
	}
}

// Java: InlineVariables.InliningBehavior.StandardVarExpert.isNeverAssigned()
func (e *InlineVariables_InliningBehavior_StandardVarExpert) IsNeverAssigned() bool {
	if e.isNeverAssigned.IsKnown() {
		return e.isNeverAssigned.GetKnownValue()
	} else if e.initializationReference.IsKnownNotNull() ||
		e.isAssignedOnceInLifetime.IsKnownToBe(true) {
		return e.isNeverAssigned.SetKnownValueOnce(false)
	} else {
		return e.isNeverAssigned.SetKnownValueOnce(e.referenceInfo.IsNeverAssigned())
	}
}

// Java: InlineVariables.InliningBehavior.StandardVarExpert.isWellDefined()
func (e *InlineVariables_InliningBehavior_StandardVarExpert) IsWellDefined() bool {
	if e.isWellDefined.IsKnown() {
		return e.isWellDefined.GetKnownValue()
	} else {
		return e.isWellDefined.SetKnownValueOnce(e.referenceInfo.IsWellDefined())
	}
}

// Java: InlineVariables.InliningBehavior.StandardVarExpert.isAssignedOnceInLifetime()
func (e *InlineVariables_InliningBehavior_StandardVarExpert) IsAssignedOnceInLifetime() bool {
	if e.isAssignedOnceInLifetime.IsKnown() {
		return e.isAssignedOnceInLifetime.GetKnownValue()
	} else {
		return e.isAssignedOnceInLifetime.SetKnownValueOnce(
			e.referenceInfo.IsAssignedOnceInLifetime())
	}
}

// Java: InlineVariables.InliningBehavior.StandardVarExpert.isReferencedWeakly()
func (e *InlineVariables_InliningBehavior_StandardVarExpert) IsReferencedWeakly() bool {
	if e.isReferencedWeakly.IsKnown() {
		return e.isReferencedWeakly.GetKnownValue()
	} else {
		return e.isReferencedWeakly.SetKnownValueOnce(e.referenceInfo.IsReferencedWeakly())
	}
}

// IsWellDefinedAssignedOnce is a more efficient way to ask if the
// variable is both well defined and assigned exactly once in its lifetime.
//
// Java: InlineVariables.InliningBehavior.StandardVarExpert.isWellDefinedAssignedOnce()
func (e *InlineVariables_InliningBehavior_StandardVarExpert) IsWellDefinedAssignedOnce() bool {
	if e.isWellDefinedAssignedOnce.IsKnown() {
		return e.isWellDefinedAssignedOnce.GetKnownValue()
	} else {
		// Check first to see if either is already known to be false before calculating anything.
		return e.isWellDefinedAssignedOnce.SetKnownValueOnce(
			!e.isWellDefined.IsKnownToBe(false) &&
				!e.isAssignedOnceInLifetime.IsKnownToBe(false) &&
				// isWellDefined is generally less expensive to calculate
				e.IsWellDefined() &&
				e.IsAssignedOnceInLifetime())
	}
}

// Java: InlineVariables.InliningBehavior.StandardVarExpert.getInitialization()
func (e *InlineVariables_InliningBehavior_StandardVarExpert) GetInitialization() *Reference {
	if e.initializationReference.IsKnown() {
		return e.initializationReference.GetKnownValue()
	} else {
		var initRef *Reference
		if e.isDeclaredOrInferredConstant {
			initRef = e.referenceInfo.GetInitializingReferenceForConstants()
		} else {
			initRef = e.referenceInfo.GetInitializingReference()
		}
		return e.initializationReference.SetKnownValueOnce(initRef)
	}
}

// Java: InlineVariables.InliningBehavior.StandardVarExpert.hasValidDeclaration()
func (e *InlineVariables_InliningBehavior_StandardVarExpert) HasValidDeclaration() bool {
	return e.outer.IsValidDeclaration(e.declaration)
}

// Java: InlineVariables.InliningBehavior.StandardVarExpert.hasValidInitialization()
func (e *InlineVariables_InliningBehavior_StandardVarExpert) HasValidInitialization() bool {
	return e.outer.IsValidInitialization(e.GetInitialization())
}

// Java: InlineVariables.InliningBehavior.StandardVarExpert.allReferencesAreValid()
func (e *InlineVariables_InliningBehavior_StandardVarExpert) AllReferencesAreValid() bool {
	if e.allReferencesAreValid.IsKnown() {
		return e.allReferencesAreValid.GetKnownValue()
	} else {
		if e.HasValidDeclaration() && (e.HasValidInitialization() || e.IsNeverAssigned()) {
			initialization := e.GetInitialization()
			for _, ref := range e.referenceInfo.references {
				if ref != e.declaration && ref != initialization && !e.outer.IsValidReference(ref) {
					return e.allReferencesAreValid.SetKnownValueOnce(false)
				}
			}
			return e.allReferencesAreValid.SetKnownValueOnce(true)
		}
		return e.allReferencesAreValid.SetKnownValueOnce(false)
	}
}

// Java: InlineVariables.InliningBehavior.StandardVarExpert.analyze()
func (e *InlineVariables_InliningBehavior_StandardVarExpert) Analyze() {
	if InlineVariables_HasNoInlineAnnotation(e.v) || e.IsReferencedWeakly() {
		e.markNegativeAnalysis()
		return
	}
	// Go-only: a variable whose only reference is its declaration can only
	// be removed as an alias (var a = b; see ReanalyzeAfterAliasedVar), so
	// the others skip the analysis.
	if len(e.referenceInfo.references) == 1 && !e.isDeclaredOrInferredConstant {
		if e.outer.outer.mode == InlineVariables_Mode_CONSTANTS_ONLY ||
			!inlineVariables_isDeclaredAsAlias(e.declaration.nameNode) {
			e.markNegativeAnalysis()
			return
		}
	}
	initialization := e.GetInitialization()
	var initValue rhino.Node
	if initialization != nil {
		initValue = initialization.GetAssignedValue()
	}
	e.initValAnalysis = InlineVariables_InliningBehavior_StandardVarExpert_InitialValueAnalysis{
		outer: e,
		value: initValue,
	}
	initialValueAnalysis := &e.initValAnalysis

	if e.isDeclaredOrInferredConstant &&
		initialValueAnalysis.IsImmutableValue() &&
		// isAssignedOnceInLifetime() is much more expensive than the other checks
		e.IsAssignedOnceInLifetime() {
		e.inlineWellDefinedVarAnalysis(initValue)
		return
	}

	if e.outer.outer.mode == InlineVariables_Mode_CONSTANTS_ONLY {
		// If we're in constants-only mode, don't run more aggressive
		// inlining heuristics. See InlineConstantsTest.
		e.markNegativeAnalysis()
		return
	}

	if initialValueAnalysis.IsAlias() {
		e.analysisKind = analysisKindWaitForAlias
		e.analysisAliasedVar = initialValueAnalysis.GetAliasedVar()
		return
	}

	e.AnalyzeWithInitialValue(initialization, initValue, initialValueAnalysis)
}

// inlineVariables_isDeclaredAsAlias reports whether n is a name declared
// with another name as its value: the a of var a = b.
func inlineVariables_isDeclaredAsAlias(n rhino.Node) bool {
	name, ok := ast.AsName(n)
	if !ok {
		return false
	}
	value, ok := name.Init()
	if !ok {
		return false
	}
	_, ok = value.AsName()
	return ok
}

// Java: InlineVariables.InliningBehavior.StandardVarExpert.analyzeWithInitialValue(Reference, Node, InitialValueAnalysis)
func (e *InlineVariables_InliningBehavior_StandardVarExpert) AnalyzeWithInitialValue(
	initialization *Reference, initValue rhino.Node, initialValueAnalysis *InlineVariables_InliningBehavior_StandardVarExpert_InitialValueAnalysis) {
	refCount := len(e.referenceInfo.references)

	// TODO(bradfordcsmith): We could remove the `refCount > 1` here, but:
	//  1. That will require some additional logic to handle stuff like named function
	//     expressions and avoiding the removal of side-effects.
	//  2. RemoveUnusedCode will remove those cases anyway.
	//  3. The unit tests will have to be more verbose to make sure stuff we don't want to be
	//  removed is used.
	if refCount > 1 && e.AllReferencesAreValid() {
		if e.IsNeverAssigned() {
			e.analysisKind = analysisKindInlineNeverAssigned
			return
		}
		if e.IsWellDefined() &&
			(initialValueAnalysis.IsImmutableValue() ||
				(initValue.IsThis() && !e.referenceInfo.IsThisRebound())) {
			// if the variable is referenced more than once, we can only
			// inline it if it's immutable and never defined before referenced.
			e.inlineWellDefinedVarAnalysis(initValue)
			return
		}

		firstReadRefIndex := 2
		if e.declaration == initialization {
			firstReadRefIndex = 1
		}
		numReadRefs := refCount - firstReadRefIndex
		if numReadRefs == 0 {
			// The only reference is the initialization.
			// Remove the assignment and the variable declaration.
			e.inlineWellDefinedVarAnalysis(initValue)
			return
		}

		if numReadRefs == 1 {
			// The variable is likely only read once, so we can try some more complex inlining
			// heuristics.
			singleReadReference := e.referenceInfo.references[firstReadRefIndex]
			if e.outer.CanInline(e.declaration, initialization, singleReadReference, initValue) {
				// A custom inline method is needed for this case.
				e.analysisKind = analysisKindInlineSingleRef
				e.analysisInitRef = initialization
				e.analysisSingleReadRef = singleReadReference
				return
			}
		}
	}

	e.markNegativeAnalysis()
}

// Java: InlineVariables.InliningBehavior.StandardVarExpert.getNegativeInlineVarAnalysis()
func (e *InlineVariables_InliningBehavior_StandardVarExpert) markNegativeAnalysis() {
	// If we already know whether it is safe to inline aliases of this variable,
	// record the resolved state immediately.
	if e.mayBeAParameterModifiedViaArguments ||
		e.outer.outer.mode == InlineVariables_Mode_CONSTANTS_ONLY ||
		len(e.referenceInfo.references) <= 1 ||
		e.isWellDefinedAssignedOnce.IsKnownToBe(false) ||
		e.isWellDefined.IsKnownToBe(false) ||
		e.isAssignedOnceInLifetime.IsKnownToBe(false) {
		e.analysisKind = analysisKindNegativeUnsafe
		return
	}
	if e.isWellDefinedAssignedOnce.IsKnownToBe(true) ||
		(e.isWellDefined.IsKnownToBe(true) && e.isAssignedOnceInLifetime.IsKnownToBe(true)) {
		e.analysisKind = analysisKindNegativeSafe
		return
	}

	// Delay calculating safety until we're actually asked.
	e.analysisKind = analysisKindDelayedNegative
}

func (e *InlineVariables_InliningBehavior_StandardVarExpert) negativeHandledVar() inlineVariables_handledVar {
	switch e.analysisKind {
	case analysisKindNegativeSafe:
		return inlineVariables_handledVar{state: handledVarSafe}
	case analysisKindDelayedNegative:
		if e.isWellDefined.IsKnownToBe(true) {
			return inlineVariables_handledVar{referenceInfo: e.referenceInfo, state: handledVarLazyWellDefined}
		}
		return inlineVariables_handledVar{referenceInfo: e.referenceInfo, state: handledVarLazy}
	default:
		return inlineVariables_handledVar{state: handledVarUnsafe}
	}
}

func (e *InlineVariables_InliningBehavior_StandardVarExpert) inlineWellDefinedVarAnalysis(initValue rhino.Node) {
	e.analysisKind = analysisKindInlineWellDefined
	e.analysisInitValue = initValue
}

// Java: InlineVariables.InliningBehavior.StandardVarExpert.reanalyzeAfterAliasedVar(Var, InlineVarAnalysis)
func (e *InlineVariables_InliningBehavior_StandardVarExpert) ReanalyzeAfterAliasedVar(
	aliasedVar *Var, aliasedHandled *inlineVariables_handledVar) {
	initialization := e.GetInitialization()
	var initValue rhino.Node
	if initialization != nil {
		initValue = initialization.GetAssignedValue()
	}
	e.initValAnalysis = InlineVariables_InliningBehavior_StandardVarExpert_InitialValueAnalysis{
		outer: e,
		value: initValue,
	}
	initialValueAnalysis := &e.initValAnalysis

	if e.isDeclaredOrInferredConstant &&
		initialValueAnalysis.IsImmutableValue() &&
		// isAssignedOnceInLifetime() is much more expensive than the other checks
		e.IsAssignedOnceInLifetime() {
		e.inlineWellDefinedVarAnalysis(initValue)
		return
	}

	if initialValueAnalysis.IsAlias() &&
		aliasedVar == initialValueAnalysis.GetAliasedVar() &&
		aliasedHandled.isSafeToInlineAliases() &&
		e.IsWellDefinedAssignedOnce() {
		// The variable we aliased couldn't be inlined itself, or it was an alias for another
		// variable that got inlined in its place.
		// However, it is safe to inline the name assigned to this variable now.
		e.inlineWellDefinedVarAnalysis(initValue)
		return
	}
	e.AnalyzeWithInitialValue(initialization, initValue, initialValueAnalysis)
}

// InlineVariables_InliningBehavior_StandardVarExpert_InitialValueAnalysis is
// information about a value used to initialize a variable.
type InlineVariables_InliningBehavior_StandardVarExpert_InitialValueAnalysis struct {
	outer *InlineVariables_InliningBehavior_StandardVarExpert

	value rhino.Node

	isImmutableValue InlineVariables_InitiallyUnknown[bool]
	aliasedVar       InlineVariables_InitiallyUnknown[*Var]
}

// Java: InlineVariables.InliningBehavior.StandardVarExpert.InitialValueAnalysis.isImmutableValue()
func (a *InlineVariables_InliningBehavior_StandardVarExpert_InitialValueAnalysis) IsImmutableValue() bool {
	if a.isImmutableValue.IsKnown() {
		return a.isImmutableValue.GetKnownValue()
	} else {
		return a.isImmutableValue.SetKnownValueOnce(
			!a.value.IsNil() && NodeUtil_IsImmutableValue(a.value))
	}
}

// Java: InlineVariables.InliningBehavior.StandardVarExpert.InitialValueAnalysis.isAlias()
func (a *InlineVariables_InliningBehavior_StandardVarExpert_InitialValueAnalysis) IsAlias() bool {
	return a.GetAliasedVar() != nil
}

// Java: InlineVariables.InliningBehavior.StandardVarExpert.InitialValueAnalysis.getAliasedVar()
func (a *InlineVariables_InliningBehavior_StandardVarExpert_InitialValueAnalysis) GetAliasedVar() *Var {
	if a.aliasedVar.IsKnown() {
		return a.aliasedVar.GetKnownValue()
	} else {
		if name, ok := ast.AsName(a.value); ok {
			aliasedName := name.Ident()
			// Never consider this variable to be an alias of itself.
			if aliasedName == a.outer.v.GetName() {
				return a.aliasedVar.SetKnownValueOnce(nil)
			}
			return a.aliasedVar.SetKnownValueOnce(a.outer.v.GetScope().(*Scope).GetVar(aliasedName))
		}
		return a.aliasedVar.SetKnownValueOnce(nil)
	}
}

// checkScopeVarEligibility handles the fast bailouts of
// Java's InliningBehavior.createVarExpert before initializing a
// StandardVarExpert.
func (b *InlineVariables_InliningBehavior) checkScopeVarEligibility(
	v *Var,
	referenceInfo *ReferenceCollection,
	isLocal bool,
) (isDeclaredOrInferredConstant bool, bailoutState inlineVariables_handledState, bailedOut bool) {
	if referenceInfo == nil {
		// If we couldn't collect any reference info, don't try to inline and assume it's unsafe
		// to inline aliases of this variable, too.
		return false, handledVarUnsafe, true
	}

	if v.IsExtern() {
		// TODO(bradfordcsmith): Extern variables are generally unsafe to inline.
		return false, handledVarUnsafe, true
	}

	isDeclaredOrInferredConstant = inlineVariables_isDeclaredOrInferredConst(v)
	if !isDeclaredOrInferredConstant && b.outer.mode == InlineVariables_Mode_CONSTANTS_ONLY {
		// If we're only inlining constants, then we shouldn't inline an alias.
		return false, handledVarUnsafe, true
	}

	cc := b.codingConvention
	if cc == nil {
		cc = b.outer.compiler.GetCodingConvention()
		b.codingConvention = cc
	}

	if cc.IsExportedWithLocal(v.GetName() /* local= */, isLocal) {
		// If the variable is exported, it might be assigned a new value by code we cannot see,
		// so aliases to it are creating snapshots of its state.
		// We cannot inline this variable or its aliases.
		return isDeclaredOrInferredConstant, handledVarUnsafe, true
	}

	if cc.IsPropertyRenameFunction(v.nameNode) {
		// It's not terribly likely that anything creates an alias of our special property rename
		// function, but its value never changes, so it should be safe to inline aliases to it.
		return isDeclaredOrInferredConstant, handledVarSafe, true
	}

	return isDeclaredOrInferredConstant, handledVarUnseen, false
}

// IsSafeUseOfArguments: true if `argumentsNode` is a use of `arguments` we
// know won't modify any parameter values.
//
// In sloppy mode it is possible to change the value of a function parameter
// by assigning to the corresponding entry in `arguments`.
//
// This method checks for just a few common uses that are known to be safe.
// However, we may soon remove this check entirely because unsafe uses
// aren't worth worrying about. See the comment on
// varsInThisScopeMayBeModifiedUsingArguments.
//
// Java: InlineVariables.InliningBehavior.isSafeUseOfArguments(Node)
func (b *InlineVariables_InliningBehavior) IsSafeUseOfArguments(argumentsNode rhino.Node) bool {
	jutil.CheckArgument(argumentsNode.MatchesName("arguments"))
	// `arguments[i]` that is only read, not assigned
	// or `fn.apply(thisArg, arguments)`
	return b.IsTargetOfPropertyRead(argumentsNode) ||
		b.IsSecondArgumentToDotApplyMethod(argumentsNode)
}

// IsTargetOfPropertyRead: true if `n` is the object in a property access
// expression (`obj[prop]` or `obj.prop` or an optional chain version of one
// of those) and the property is only being read, not modified.
//
// Java: InlineVariables.InliningBehavior.isTargetOfPropertyRead(Node)
func (b *InlineVariables_InliningBehavior) IsTargetOfPropertyRead(n rhino.Node) bool {
	get, ok := ast.AsNormalOrOptChainGet(n.GetParent())
	return ok && get.Object().Is(n) && !NodeUtil_IsLValue(get.Node())
}

// IsSecondArgumentToDotApplyMethod: true if `n` is being used in a call like
// this: `fn.apply(thisArg, n)`.
//
// Java: InlineVariables.InliningBehavior.isSecondArgumentToDotApplyMethod(Node)
func (b *InlineVariables_InliningBehavior) IsSecondArgumentToDotApplyMethod(n rhino.Node) bool {
	call, ok := ast.AsNormalOrOptChainCall(n.GetParent())
	if !ok {
		return false
	}
	callee, ok := call.Callee().AsNormalOrOptChainGetProp()
	if !ok || callee.Property() != "apply" {
		return false
	}
	thisArg, ok := call.FirstArg()
	if !ok {
		return false
	}
	secondArg, ok := thisArg.Next()
	return ok && secondArg.Is(n)
}

// Inline does the actual work of inlining a single declaration into a
// single reference.
//
// Java: InlineVariables.InliningBehavior.inline(Reference, Reference, Reference)
func (b *InlineVariables_InliningBehavior) Inline(decl *Reference, init *Reference, ref *Reference) {
	value := init.GetAssignedValue()
	jutil.CheckState(!value.IsNil())
	// Check for function declarations before the value is moved in the AST.
	isFunctionDeclaration := NodeUtil_IsFunctionDeclaration(value)
	if isFunctionDeclaration {
		// We're inlining this function declaration because there was only one reference to it,
		// probably the rhs of an assignment statement. Since we're eliminating the only
		// reference, we can get rid of it. See b/277538101.
		fn := ast.MustFunction(value)
		fn.RemoveName()

		// In addition to changing the containing scope, inlining function declarations also changes
		// the function name scope from the containing scope to the inner scope.
		b.outer.compiler.ReportChangeToChangeScope(value)
		b.outer.compiler.ReportChangeToEnclosingScope(value.GetParent())
	}
	b.InlineValue(ref.GetNode(), value.Detach())
	if decl != init {
		expressRoot := init.GetGrandparent()
		jutil.CheckState(expressRoot.IsExprResult())
		NodeUtil_RemoveChild(expressRoot.GetParent(), expressRoot)
	}
	// Function declarations have already been removed.
	if !isFunctionDeclaration {
		b.RemoveDeclaration(decl)
	}
}

// InlineWellDefinedVariable inlines an immutable variable into all of its
// references.
//
// Java: InlineVariables.InliningBehavior.inlineWellDefinedVariable(Var, Node, List<Reference>)
func (b *InlineVariables_InliningBehavior) InlineWellDefinedVariable(v *Var, value rhino.Node, refSet []*Reference) {
	for _, r := range refSet {
		if r.GetNode() == v.GetNameNode() {
			b.RemoveDeclaration(r)
		} else if r.IsSimpleAssignmentToName() {
			/*
			 * This is the initialization.
			 *
			 * <p>Replace the entire assignment with just the value, and use the original value node
			 * in case it contains references to variables that still require inlining.
			 */
			b.InlineValue(r.GetParent(), value.Detach())
		} else {
			toRemove := r.GetNode()
			if toRemove.CanCloneLeafInPlace() && value.CanCloneLeafInPlace() {
				b.outer.compiler.ReportChangeToEnclosingScope(toRemove)
				castColor := toRemove.GetColor()
				hasCastColor := castColor != nil && toRemove.IsColorFromTypeCast()
				toRemove.OverwriteLeafFromClone(value)
				if hasCastColor {
					toRemove.SetColor(castColor)
					toRemove.SetColorFromTypeCast()
				}
			} else {
				clonedValue := value.CloneTree()
				NodeUtil_MarkNewScopesChanged(clonedValue, b.outer.compiler)
				b.InlineValue(toRemove, clonedValue)
			}
		}
	}
}

// RemoveDeclaration removes the given VAR declaration.
//
// Java: InlineVariables.InliningBehavior.removeDeclaration(Reference)
func (b *InlineVariables_InliningBehavior) RemoveDeclaration(decl *Reference) {
	varNode, ok := ast.AsNameDeclaration(decl.GetParent())
	jutil.CheckStateMsg1(ok, "%s", decl.GetParent())
	grandparent := decl.GetGrandparent()

	b.outer.compiler.ReportChangeToEnclosingScope(decl.GetNode())
	decl.GetNode().Detach()
	// Remove var node if empty
	if !varNode.HasDeclarators() {
		NodeUtil_RemoveChild(grandparent, varNode.Node())
	}
}

// Java: InlineVariables.InliningBehavior.inlineValue(Node, Node)
func (b *InlineVariables_InliningBehavior) InlineValue(toRemove rhino.Node, toInsert rhino.Node) {
	b.outer.compiler.ReportChangeToEnclosingScope(toRemove)

	// Help type-based optimizations by propagating more specific types from type assertions
	if toRemove.GetColor() != nil && toRemove.IsColorFromTypeCast() {
		toInsert.SetColor(toRemove.GetColor())
		toInsert.SetColorFromTypeCast()
	}
	toRemove.ReplaceWith(toInsert)
	NodeUtil_MarkFunctionsDeleted(toRemove, b.outer.compiler)
}

// CanInline returns true if the provided reference and declaration can be
// safely inlined according to our criteria.
//
// Java: InlineVariables.InliningBehavior.canInline(Reference, Reference, Reference, Node)
func (b *InlineVariables_InliningBehavior) CanInline(
	declaration *Reference, initialization *Reference, reference *Reference, initValue rhino.Node) bool {
	// If the value is read more than once, skip it.
	// VAR declarations and EXPR_RESULT don't need the value, but other
	// ASSIGN expressions parents do.
	if declaration != initialization && !initialization.GetGrandparent().IsExprResult() {
		return false
	}

	// Be very conservative and do not cross control structures or scope boundaries
	if declaration.GetBasicBlock() != initialization.GetBasicBlock() ||
		declaration.GetBasicBlock() != reference.GetBasicBlock() {
		return false
	}

	// Do not inline into a call node. This would change
	// the context in which it was being called. For example,
	//   var a = b.c;
	//   a();
	// should not be inlined, because it calls a in the context of b
	// rather than the context of the window.
	//   var a = b.c;
	//   f(a)
	// is OK.
	jutil.CheckState(!initValue.IsNil())
	call, isCall := ast.AsCall(reference.GetParent())
	if initValue.IsGetProp() && isCall && call.Callee().Is(reference.GetNode()) {
		return false
	}

	if initValue.IsFunction() && isCall {
		convention := b.outer.compiler.GetCodingConvention()
		// Bug 2388531: Don't inline subclass definitions into class defining
		// calls as this confused class removing logic.
		relationship := convention.GetClassesDefinedByCall(call.Node())
		if relationship != nil {
			return false
		}

		// issue 668: Don't inline singleton getter methods
		// calls as this confused class removing logic.
		if convention.GetSingletonGetterClassName(call.Node()) != nil {
			return false
		}
	}

	if initialization.GetScope() != declaration.GetScope() ||
		!initialization.GetScope().Contains(reference.GetScope()) {
		return false
	}

	return b.CanMoveAggressively(initValue) || b.CanMoveModerately(initialization, reference)
}

// CanMoveAggressively: if the value is a literal, we can cross more
// boundaries to inline it.
//
// Java: InlineVariables.InliningBehavior.canMoveAggressively(Node)
func (b *InlineVariables_InliningBehavior) CanMoveAggressively(value rhino.Node) bool {
	// Function expressions and other mutable objects can move within
	// the same basic block.
	return NodeUtil_IsLiteralValue(value, true) || value.IsFunction()
}

// InlineVariables_InliningBehavior_MAX_MODERATE_MOTION_STEPS is the maximum
// number of AST steps LocalVarMotion will traverse before giving up. Prevents
// quadratic traversal in large monolithic scopes with many definitions
// (b/545160233).
//
// 5000 steps is large enough to span massive side-effect-free AST structures
// like protobuf descriptor arrays without regressing code size, while
// bounding worst-case traversal time in huge monolithic scopes to fractions
// of a second.
//
// Java: InlineVariables.InliningBehavior.MAX_MODERATE_MOTION_STEPS
const InlineVariables_InliningBehavior_MAX_MODERATE_MOTION_STEPS = 5000

// CanMoveModerately: if the value of a variable is not constant, then it may
// read or modify state. Therefore it cannot be moved past anything else that
// may modify the value being read or read values that are modified.
//
// Java: InlineVariables.InliningBehavior.canMoveModerately(Reference, Reference)
func (b *InlineVariables_InliningBehavior) CanMoveModerately(initialization *Reference, reference *Reference) bool {
	// Check if declaration can be inlined without passing
	// any side-effect causing nodes.
	it := &b.motionScratch
	if NodeUtil_IsNameDeclaration(initialization.GetParent()) {
		it.initForVar(
			b.outer.compiler,
			initialization.GetNode(),        // NAME
			initialization.GetParent(),      // VAR/LET/CONST
			initialization.GetGrandparent()) // VAR/LET/CONST container
	} else if initialization.GetParent().IsAssign() {
		jutil.CheckState(initialization.GetGrandparent().IsExprResult())
		it.initForAssign(
			b.outer.compiler,
			initialization.GetNode(),                    // NAME
			initialization.GetParent(),                  // ASSIGN
			initialization.GetGrandparent(),             // EXPR_RESULT
			initialization.GetGrandparent().GetParent()) // EXPR container
	} else {
		panic(jutil.NewIllegalStateExceptionMsg(
			"Unexpected initialization parent\n" + initialization.GetParent().ToStringTree()))
	}
	targetName := reference.GetNode()
	steps := 0
	for it.HasNext() {
		curNode := it.Next()
		if curNode == targetName {
			return true
		}
		steps++
		if steps > InlineVariables_InliningBehavior_MAX_MODERATE_MOTION_STEPS {
			return false
		}
	}

	return false
}

// IsValidDeclaration returns true if the reference is a normal VAR or
// FUNCTION declaration.
//
// Java: InlineVariables.InliningBehavior.isValidDeclaration(Reference)
func (b *InlineVariables_InliningBehavior) IsValidDeclaration(declaration *Reference) bool {
	return (NodeUtil_IsNameDeclaration(declaration.GetParent()) &&
		!NodeUtil_IsLoopStructure(declaration.GetGrandparent())) ||
		NodeUtil_IsFunctionDeclaration(declaration.GetParent())
}

// IsValidInitialization returns whether there is a initial value.
//
// Java: InlineVariables.InliningBehavior.isValidInitialization(Reference)
func (b *InlineVariables_InliningBehavior) IsValidInitialization(initialization *Reference) bool {
	if initialization == nil {
		return false
	} else if initialization.IsDeclaration() {
		// The reference is a FUNCTION declaration or normal VAR declaration
		// with a value.
		if !NodeUtil_IsFunctionDeclaration(initialization.GetParent()) &&
			!initialization.GetNode().HasChildren() {
			return false
		}
	} else {
		assign, ok := ast.AsAssign(initialization.GetParent())
		jutil.CheckState(ok && assign.Target().Is(initialization.GetNode()))
	}

	return true
}

// IsValidReference returns true if the reference is a candidate for
// inlining.
//
// Java: InlineVariables.InliningBehavior.isValidReference(Reference)
func (b *InlineVariables_InliningBehavior) IsValidReference(reference *Reference) bool {
	return !reference.IsDeclaration() && !reference.IsLvalue()
}

// Java: InlineVariables.hasNoInlineAnnotation(Var)
func InlineVariables_HasNoInlineAnnotation(var_ *Var) bool {
	n := var_.nameNode
	if n.IsNil() {
		return false
	}
	if info := n.GetJSDocInfo(); info != nil {
		return info.IsNoInline()
	}
	parent := n.GetParent()
	if parent.IsNil() {
		return false
	}
	switch parent.GetToken() {
	case rhino.Token_PARAM_LIST:
		return false
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST, rhino.Token_FUNCTION, rhino.Token_CLASS:
		if parent.GetFirstChild() == n {
			if info := parent.GetJSDocInfo(); info != nil {
				return info.IsNoInline()
			}
			if gp := parent.GetParent(); !gp.IsNil() && gp.IsExport() {
				info := gp.GetJSDocInfo()
				return info != nil && info.IsNoInline()
			}
			return false
		}
	}
	jsDocInfo := NodeUtil_GetBestJSDocInfo(n)
	return jsDocInfo != nil && jsDocInfo.IsNoInline()
}

// InlineVariables_InitiallyUnknown is InlineVariables.InitiallyUnknown<T>.
//
// Go-only: language-adapter — T is comparable because Java compares values
// with ==; T's zero value is Java's null.
type InlineVariables_InitiallyUnknown[T comparable] struct {
	isKnown bool
	value   T
}

// Java: InlineVariables.InitiallyUnknown()
func NewInlineVariables_InitiallyUnknown[T comparable]() *InlineVariables_InitiallyUnknown[T] {
	return &InlineVariables_InitiallyUnknown[T]{isKnown: false}
}

// Java: InlineVariables.InitiallyUnknown.isKnown()
func (u *InlineVariables_InitiallyUnknown[T]) IsKnown() bool {
	return u.isKnown
}

// Java: InlineVariables.InitiallyUnknown.isKnownNotNull()
func (u *InlineVariables_InitiallyUnknown[T]) IsKnownNotNull() bool {
	var null T
	return u.IsKnownNotToBe(null)
}

// Java: InlineVariables.InitiallyUnknown.isKnownToBe(T)
func (u *InlineVariables_InitiallyUnknown[T]) IsKnownToBe(other T) bool {
	return u.isKnown && u.value == other
}

// Java: InlineVariables.InitiallyUnknown.isKnownNotToBe(T)
func (u *InlineVariables_InitiallyUnknown[T]) IsKnownNotToBe(other T) bool {
	return u.isKnown && u.value != other
}

// Java: InlineVariables.InitiallyUnknown.setKnownValueOnce(T)
func (u *InlineVariables_InitiallyUnknown[T]) SetKnownValueOnce(value T) T {
	jutil.CheckStateMsg(!u.isKnown, "already known")
	u.value = value
	u.isKnown = true
	return value
}

// Java: InlineVariables.InitiallyUnknown.getKnownValue()
func (u *InlineVariables_InitiallyUnknown[T]) GetKnownValue() T {
	jutil.CheckStateMsg(u.isKnown, "not yet known")
	return u.value
}
