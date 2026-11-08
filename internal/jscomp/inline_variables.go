// Port of com/google/javascript/jscomp/InlineVariables.java @ closure-compiler v20261104

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
	"github.com/paullepoulpe/jsgo/internal/strpool"
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

	mode             InlineVariables_Mode
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

	scratchExpert inlineVariables_standardVarExpert
	motionScratch NodeIterators_LocalVarMotion
}

// InlineVariables_InliningBehavior builds up information about nodes in each
// scope. When exiting the scope, inspects all variables in that scope, and
// inlines any that we can.
type InlineVariables_InliningBehavior = InlineVariables

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
	if name, ok := ast.AsName(n); ok {
		return name.IsConstantName()
	}
	return n.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME)
}

// Java: InlineVariables.Mode.name()
func (m InlineVariables_Mode) String() string { return inlineVariables_Mode_names[m] }

// Java: InlineVariables(AbstractCompiler, Mode)
func NewInlineVariables(compiler AbstractCompiler, mode InlineVariables_Mode) *InlineVariables {
	return &InlineVariables{compiler: compiler, mode: mode}
}

func (iv *InlineVariables) resetBehavior() {
	if iv.compiler != nil {
		iv.codingConvention = iv.compiler.GetCodingConvention()
	}
	iv.currentScope = nil
	clear(iv.handledScopeVars)
	iv.handledScopeVars = iv.handledScopeVars[:0]
	if len(iv.aliasRetryHandlers) > 0 {
		clear(iv.aliasRetryHandlers)
	}
	iv.unpinnedRetries = 0
	iv.hasGlobalAliasRetries = false
	iv.onGlobalAliasRetry = nil
}

// Java: InlineVariables.process(Node, Node)
func (iv *InlineVariables) Process(externs rhino.Node, root rhino.Node) {
	if iv.processParallel(externs, root) {
		return
	}
	iv.resetBehavior()
	callback := NewReferenceCollectorWithVarFilter(
		iv.compiler,
		iv,
		NewSyntacticScopeCreator(iv.compiler),
		inlineVariables_Mode_varPredicate[iv.mode])
	callback.enableLocalReclamation(func() bool {
		if len(iv.aliasRetryHandlers) == 0 {
			iv.unpinnedRetries = 0
			callback.clearLocalBase()
			return true
		}
		if iv.unpinnedRetries == 0 {
			return true
		}
		iv.unpinnedRetries = 0
		return false
	})
	callback.Process(externs, root)
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

func scopeVarSlot(scope *Scope, v *Var) int {
	if idx := v.GetIndex(); uint(idx) < uint(len(scope.varList)) && scope.varList[idx] == v {
		return idx
	}
	return slices.Index(scope.varList, v)
}

func (iv *InlineVariables) recordHandledVar(v *Var, hv inlineVariables_handledVar) {
	if iv.mode == InlineVariables_Mode_CONSTANTS_ONLY {
		return
	}
	handledPtr := &hv
	if scope := iv.currentScope; v != nil && scope != nil && v.scope == scope {
		if slot := scopeVarSlot(scope, v); slot >= 0 {
			iv.handledScopeVars[slot] = hv
			handledPtr = &iv.handledScopeVars[slot]
		}
	}
	// Retry any aliases we saw before handling this Var
	iv.RetryAliases(v, handledPtr)
}

func (iv *InlineVariables) getHandledVar(scope *Scope, aliasedVar *Var) *inlineVariables_handledVar {
	if aliasedVar == nil || aliasedVar.scope != scope || scope != iv.currentScope {
		return nil
	}
	if slot := scopeVarSlot(scope, aliasedVar); slot >= 0 && iv.handledScopeVars[slot].state != handledVarUnseen {
		return &iv.handledScopeVars[slot]
	}
	return nil
}

func (iv *InlineVariables) initScratchVarExpert(
	v *Var,
	referenceInfo *ReferenceCollection,
	isDeclaredOrInferredConstant bool,
	mayBeAParameterModifiedViaArguments bool,
) *inlineVariables_standardVarExpert {
	iv.scratchExpert = inlineVariables_standardVarExpert{
		outer:                               iv,
		v:                                   v,
		referenceInfo:                       referenceInfo,
		declaration:                         referenceInfo.getRef(0),
		isDeclaredOrInferredConstant:        isDeclaredOrInferredConstant,
		mayBeAParameterModifiedViaArguments: mayBeAParameterModifiedViaArguments,
	}
	return &iv.scratchExpert
}

func (iv *InlineVariables) allocRetryExpert(
	src *inlineVariables_standardVarExpert,
) *inlineVariables_standardVarExpert {
	e := &inlineVariables_standardVarExpert{}
	*e = *src
	e.initValAnalysis.outer = e
	return e
}

func (iv *InlineVariables) allocRetryHandler(
	v *Var,
	expert *inlineVariables_standardVarExpert,
) *InlineVariables_InliningBehavior_AliasInlineRetryHandler {
	return &InlineVariables_InliningBehavior_AliasInlineRetryHandler{
		outer:  iv,
		v:      v,
		expert: expert,
	}
}

func (iv *InlineVariables) putAliasRetryHandler(
	aliasedVar *Var,
	retryHandler *InlineVariables_InliningBehavior_AliasInlineRetryHandler,
) {
	if iv.onGlobalAliasRetry != nil && !aliasedVar.IsLocal() {
		iv.hasGlobalAliasRetries = true
		iv.unpinnedRetries++
		iv.onGlobalAliasRetry(aliasedVar, retryHandler)
		return
	}
	if iv.aliasRetryHandlers == nil {
		iv.aliasRetryHandlers = make(map[*Var][]*InlineVariables_InliningBehavior_AliasInlineRetryHandler, 8)
	}
	iv.aliasRetryHandlers[aliasedVar] = append(iv.aliasRetryHandlers[aliasedVar], retryHandler)
	iv.unpinnedRetries++
}

// Java: InlineVariables.InliningBehavior.afterExitScope(NodeTraversal, ReferenceMap)
func (iv *InlineVariables) AfterExitScope(t *NodeTraversal, referenceMap *ReferenceCollector) {
	iv.DoInlinesForScope(t, referenceMap)
}

// DoInlinesForScope: for all variables in this scope, see if they are only
// used once. If it looks safe to do so, inline them.
//
// Java: InlineVariables.InliningBehavior.doInlinesForScope(NodeTraversal, ReferenceMap)
func (iv *InlineVariables) DoInlinesForScope(t *NodeTraversal, referenceMap *ReferenceCollector) {
	// Any variables we completed in an earlier scope are now out of scope,
	// so we can clear this slice.
	clear(iv.handledScopeVars)
	iv.handledScopeVars = iv.handledScopeVars[:0]
	iv.currentScope = nil
	if len(t.scopes) > 0 && t.scopes[len(t.scopes)-1].scope == nil {
		return
	}
	scope := t.GetScope()
	if scope.GetVarCount() == 0 {
		return
	}
	isLocalScope := scope.IsLocal()
	if iv.mode == InlineVariables_Mode_LOCALS_ONLY && !isLocalScope && len(iv.aliasRetryHandlers) == 0 {
		return
	}
	if iv.mode != InlineVariables_Mode_CONSTANTS_ONLY {
		iv.currentScope = scope
		nVars := len(scope.varList)
		iv.handledScopeVars = slices.Grow(iv.handledScopeVars, nVars)[:nVars]
	}
	mayBeAParameterModifiedViaArguments :=
		iv.VarsInThisScopeMayBeModifiedUsingArguments(scope, referenceMap)
	for i := 0; i < len(scope.varList); i++ {
		v := scope.varList[i]
		if v == nil {
			continue
		}
		if iv.mode == InlineVariables_Mode_CONSTANTS_ONLY && !inlineVariables_isDeclaredOrInferredConst(v) {
			continue
		}
		referenceInfo := referenceMap.GetReferences(v)
		isConst, bailoutState, bailedOut := iv.checkScopeVarEligibility(v, referenceInfo, isLocalScope)
		if bailedOut {
			iv.recordHandledVar(v, inlineVariables_handledVar{state: bailoutState})
			continue
		}
		expert := iv.initScratchVarExpert(v, referenceInfo, isConst, mayBeAParameterModifiedViaArguments)
		expert.Analyze()
		if expert.ShouldWaitForAliasedVar() {
			// We need to wait for the aliased variable to be handled.
			aliasedVar := expert.GetAliasedVar()
			if aliasedHandled := iv.getHandledVar(scope, aliasedVar); aliasedHandled != nil {
				// We've actually already completed the aliased Var
				iv.handleAliasedVarCompletion(v, expert, aliasedVar, aliasedHandled)
			} else {
				// wait for completion of the aliased Var
				if expert == &iv.scratchExpert {
					expert = iv.allocRetryExpert(&iv.scratchExpert)
				}
				retryHandler := iv.allocRetryHandler(v, expert)
				iv.putAliasRetryHandler(aliasedVar, retryHandler)
			}
		} else {
			iv.completeVarExpert(v, expert)
		}
	}
}

func (iv *InlineVariables) completeVarExpert(
	v *Var,
	expert *inlineVariables_standardVarExpert,
) {
	var hv inlineVariables_handledVar
	if expert.ShouldInline() {
		expert.PerformInline()
		hv = inlineVariables_handledVar{state: handledVarSafe}
	} else {
		hv = expert.negativeHandledVar()
	}
	iv.recordHandledVar(v, hv)
}

func (iv *InlineVariables) handleAliasedVarCompletion(
	v *Var,
	expert *inlineVariables_standardVarExpert,
	aliasedVar *Var,
	aliasedHandled *inlineVariables_handledVar,
) {
	expert.ReanalyzeAfterAliasedVar(aliasedVar, aliasedHandled)
	jutil.CheckStateMsg1(
		!expert.ShouldWaitForAliasedVar(), "expert for %s asked to wait a second time", v)
	iv.completeVarExpert(v, expert)
}

var inlineVariablesApplyID = uint32(strpool.Intern("apply"))

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
func (iv *InlineVariables) VarsInThisScopeMayBeModifiedUsingArguments(
	scope *Scope, referenceMap *ReferenceCollector) bool {
	if fn, ok := ast.AsFunction(scope.GetRootNode()); ok && !fn.IsArrow() {
		refs := referenceMap.getReferencesByScopedNameID(fn.Node(), implicitArgumentsID)
		if refs != nil && len(refs.references) > 0 {
			for _, ref := range refs.references {
				if !iv.IsSafeUseOfArguments(ref.GetNode()) {
					return true
				}
			}
		}
	}
	return false
}

// Java: InlineVariables.InliningBehavior.retryAliases(Var, InlineVarAnalysis)
func (iv *InlineVariables) RetryAliases(aliasedVar *Var, aliasedHandled *inlineVariables_handledVar) {
	if len(iv.aliasRetryHandlers) == 0 {
		return
	}
	handlers, ok := iv.aliasRetryHandlers[aliasedVar]
	if !ok {
		return
	}
	delete(iv.aliasRetryHandlers, aliasedVar)
	if len(iv.aliasRetryHandlers) == 0 {
		iv.unpinnedRetries = 0
	} else if iv.unpinnedRetries >= len(handlers) {
		iv.unpinnedRetries -= len(handlers)
	} else {
		iv.unpinnedRetries = 0
	}
	for _, aliasInlineRetryHandler := range handlers {
		aliasInlineRetryHandler.HandleAliasedVarCompletion(aliasedVar, aliasedHandled)
	}
}

// InlineVariables_InliningBehavior_AliasInlineRetryHandler retries inlining
// an alias variable after the original variable has been dealt with.
type InlineVariables_InliningBehavior_AliasInlineRetryHandler struct {
	outer *InlineVariables

	v      *Var
	expert *inlineVariables_standardVarExpert
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

type varExpertBit uint16

const (
	veNeverAssigned varExpertBit = 1 << iota
	veWellDefined
	veAssignedOnceInLifetime
	veReferencedWeakly
	veWellDefinedAssignedOnce
	veInitRef
	veAllRefsValid
)

// inlineVariables_standardVarExpert knows how to analyze a variable to
// determine whether it should be inlined, how to do the inlining, and
// whether it's OK to inline aliases of the variable.
type inlineVariables_standardVarExpert struct {
	outer *InlineVariables

	v                                   *Var
	referenceInfo                       *ReferenceCollection
	declaration                         *Reference
	initRef                             *Reference
	known                               varExpertBit
	val                                 varExpertBit
	isDeclaredOrInferredConstant        bool
	mayBeAParameterModifiedViaArguments bool

	initValAnalysis       inlineVariables_initialValueAnalysis
	analysisKind          standardVarAnalysisKind
	analysisAliasedVar    *Var
	analysisInitValue     rhino.Node
	analysisInitRef       *Reference
	analysisSingleReadRef *Reference
}

func (e *inlineVariables_standardVarExpert) knows(bit varExpertBit) bool {
	return e.known&bit != 0
}

func (e *inlineVariables_standardVarExpert) knownVal(bit varExpertBit) bool {
	return e.val&bit != 0
}

func (e *inlineVariables_standardVarExpert) isKnownTrue(bit varExpertBit) bool {
	return e.known&e.val&bit != 0
}

func (e *inlineVariables_standardVarExpert) isKnownFalse(bit varExpertBit) bool {
	return e.known&^e.val&bit != 0
}

func (e *inlineVariables_standardVarExpert) setOnce(bit varExpertBit, v bool) bool {
	jutil.CheckStateMsg(e.known&bit == 0, "already known")
	e.known |= bit
	if v {
		e.val |= bit
	}
	return v
}

func (e *inlineVariables_standardVarExpert) ShouldInline() bool {
	return e.analysisKind >= analysisKindInlineWellDefined
}

func (e *inlineVariables_standardVarExpert) ShouldWaitForAliasedVar() bool {
	return e.analysisKind == analysisKindWaitForAlias
}

func (e *inlineVariables_standardVarExpert) GetAliasedVar() *Var {
	if e.analysisKind == analysisKindWaitForAlias {
		return e.analysisAliasedVar
	}
	panic(jutil.NewUnsupportedOperationExceptionMsg("no aliased Var"))
}

func (e *inlineVariables_standardVarExpert) PerformInline() {
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
func (e *inlineVariables_standardVarExpert) IsNeverAssigned() bool {
	if e.knows(veNeverAssigned) {
		return e.knownVal(veNeverAssigned)
	}
	if (e.knows(veInitRef) && e.initRef != nil) || e.isKnownTrue(veAssignedOnceInLifetime) {
		return e.setOnce(veNeverAssigned, false)
	}
	return e.setOnce(veNeverAssigned, e.referenceInfo.IsNeverAssigned())
}

// Java: InlineVariables.InliningBehavior.StandardVarExpert.isWellDefined()
func (e *inlineVariables_standardVarExpert) IsWellDefined() bool {
	if e.knows(veWellDefined) {
		return e.knownVal(veWellDefined)
	}
	return e.setOnce(veWellDefined, e.referenceInfo.IsWellDefined())
}

// Java: InlineVariables.InliningBehavior.StandardVarExpert.isAssignedOnceInLifetime()
func (e *inlineVariables_standardVarExpert) IsAssignedOnceInLifetime() bool {
	if e.knows(veAssignedOnceInLifetime) {
		return e.knownVal(veAssignedOnceInLifetime)
	}
	return e.setOnce(veAssignedOnceInLifetime, e.referenceInfo.IsAssignedOnceInLifetime())
}

// Java: InlineVariables.InliningBehavior.StandardVarExpert.isReferencedWeakly()
func (e *inlineVariables_standardVarExpert) IsReferencedWeakly() bool {
	if e.knows(veReferencedWeakly) {
		return e.knownVal(veReferencedWeakly)
	}
	return e.setOnce(veReferencedWeakly, e.referenceInfo.IsReferencedWeakly())
}

// IsWellDefinedAssignedOnce is a more efficient way to ask if the
// variable is both well defined and assigned exactly once in its lifetime.
//
// Java: InlineVariables.InliningBehavior.StandardVarExpert.isWellDefinedAssignedOnce()
func (e *inlineVariables_standardVarExpert) IsWellDefinedAssignedOnce() bool {
	if e.knows(veWellDefinedAssignedOnce) {
		return e.knownVal(veWellDefinedAssignedOnce)
	}
	// Check first to see if either is already known to be false before calculating anything.
	return e.setOnce(veWellDefinedAssignedOnce,
		!e.isKnownFalse(veWellDefined) &&
			!e.isKnownFalse(veAssignedOnceInLifetime) &&
			// isWellDefined is generally less expensive to calculate
			e.IsWellDefined() &&
			e.IsAssignedOnceInLifetime())
}

// Java: InlineVariables.InliningBehavior.StandardVarExpert.getInitialization()
func (e *inlineVariables_standardVarExpert) GetInitialization() *Reference {
	if e.knows(veInitRef) {
		return e.initRef
	}
	var initRef *Reference
	if e.isDeclaredOrInferredConstant {
		initRef = e.referenceInfo.GetInitializingReferenceForConstants()
	} else {
		initRef = e.referenceInfo.GetInitializingReference()
	}
	e.known |= veInitRef
	e.initRef = initRef
	return initRef
}

// Java: InlineVariables.InliningBehavior.StandardVarExpert.hasValidDeclaration()
func (e *inlineVariables_standardVarExpert) HasValidDeclaration() bool {
	return e.outer.IsValidDeclaration(e.declaration)
}

// Java: InlineVariables.InliningBehavior.StandardVarExpert.hasValidInitialization()
func (e *inlineVariables_standardVarExpert) HasValidInitialization() bool {
	return e.outer.IsValidInitialization(e.GetInitialization())
}

// Java: InlineVariables.InliningBehavior.StandardVarExpert.allReferencesAreValid()
func (e *inlineVariables_standardVarExpert) AllReferencesAreValid() bool {
	if e.knows(veAllRefsValid) {
		return e.knownVal(veAllRefsValid)
	}
	if e.HasValidDeclaration() && (e.HasValidInitialization() || e.IsNeverAssigned()) {
		initialization := e.GetInitialization()
		for _, ref := range e.referenceInfo.references {
			if ref != e.declaration && ref != initialization && !e.outer.IsValidReference(ref) {
				return e.setOnce(veAllRefsValid, false)
			}
		}
		return e.setOnce(veAllRefsValid, true)
	}
	return e.setOnce(veAllRefsValid, false)
}

// Java: InlineVariables.InliningBehavior.StandardVarExpert.analyze()
func (e *inlineVariables_standardVarExpert) Analyze() {
	if InlineVariables_HasNoInlineAnnotation(e.v) || e.IsReferencedWeakly() {
		e.markNegativeAnalysis()
		return
	}
	// Go-only: a variable whose only reference is its declaration can only
	// be removed as an alias (var a = b; see ReanalyzeAfterAliasedVar), so
	// the others skip the analysis.
	if len(e.referenceInfo.references) == 1 && !e.isDeclaredOrInferredConstant {
		if e.outer.mode == InlineVariables_Mode_CONSTANTS_ONLY ||
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
	e.initValAnalysis = inlineVariables_initialValueAnalysis{
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

	if e.outer.mode == InlineVariables_Mode_CONSTANTS_ONLY {
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
func (e *inlineVariables_standardVarExpert) AnalyzeWithInitialValue(
	initialization *Reference, initValue rhino.Node, initialValueAnalysis *inlineVariables_initialValueAnalysis) {
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
		_, isThis := ast.AsThis(initValue)
		if e.IsWellDefined() &&
			(initialValueAnalysis.IsImmutableValue() ||
				(isThis && !e.referenceInfo.IsThisRebound())) {
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
func (e *inlineVariables_standardVarExpert) markNegativeAnalysis() {
	// If we already know whether it is safe to inline aliases of this variable,
	// record the resolved state immediately.
	if e.mayBeAParameterModifiedViaArguments ||
		e.outer.mode == InlineVariables_Mode_CONSTANTS_ONLY ||
		len(e.referenceInfo.references) <= 1 ||
		e.isKnownFalse(veWellDefinedAssignedOnce) ||
		e.isKnownFalse(veWellDefined) ||
		e.isKnownFalse(veAssignedOnceInLifetime) {
		e.analysisKind = analysisKindNegativeUnsafe
		return
	}
	if e.isKnownTrue(veWellDefinedAssignedOnce) ||
		(e.isKnownTrue(veWellDefined) && e.isKnownTrue(veAssignedOnceInLifetime)) {
		e.analysisKind = analysisKindNegativeSafe
		return
	}

	// Delay calculating safety until we're actually asked.
	e.analysisKind = analysisKindDelayedNegative
}

func (e *inlineVariables_standardVarExpert) negativeHandledVar() inlineVariables_handledVar {
	switch e.analysisKind {
	case analysisKindNegativeSafe:
		return inlineVariables_handledVar{state: handledVarSafe}
	case analysisKindDelayedNegative:
		if e.isKnownTrue(veWellDefined) {
			return inlineVariables_handledVar{referenceInfo: e.referenceInfo, state: handledVarLazyWellDefined}
		}
		return inlineVariables_handledVar{referenceInfo: e.referenceInfo, state: handledVarLazy}
	default:
		return inlineVariables_handledVar{state: handledVarUnsafe}
	}
}

func (e *inlineVariables_standardVarExpert) inlineWellDefinedVarAnalysis(initValue rhino.Node) {
	e.analysisKind = analysisKindInlineWellDefined
	e.analysisInitValue = initValue
}

// Java: InlineVariables.InliningBehavior.StandardVarExpert.reanalyzeAfterAliasedVar(Var, InlineVarAnalysis)
func (e *inlineVariables_standardVarExpert) ReanalyzeAfterAliasedVar(
	aliasedVar *Var, aliasedHandled *inlineVariables_handledVar) {
	initialization := e.GetInitialization()
	var initValue rhino.Node
	if initialization != nil {
		initValue = initialization.GetAssignedValue()
	}
	e.initValAnalysis = inlineVariables_initialValueAnalysis{
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
