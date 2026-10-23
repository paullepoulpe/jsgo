// Port of com/google/javascript/jscomp/InlineVariables.java @ closure-compiler v20261015

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
	"sync"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

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
	if n == nil {
		return false
	}
	if n.GetConstantVarFlags()&(rhino.Node_ConstantVarFlags_DECLARED|rhino.Node_ConstantVarFlags_INFERRED) != 0 {
		return true
	}
	return n.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME)
}

// Java: InlineVariables.Mode.values()
func InlineVariables_Mode_values() []InlineVariables_Mode {
	return []InlineVariables_Mode{
		InlineVariables_Mode_CONSTANTS_ONLY,
		InlineVariables_Mode_LOCALS_ONLY,
		InlineVariables_Mode_ALL,
	}
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
func (iv *InlineVariables) Process(externs *rhino.Node, root *rhino.Node) {
	behavior := acquireInliningBehavior(iv)
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
	callback.release()
	behavior.release()
}

// InlineVariables_VarExpert is responsible for analyzing a variable to
// determine if it can be inlined.
type InlineVariables_VarExpert interface {
	// Analyze is called to conduct the initial analysis.
	//
	// Java: InlineVariables.VarExpert.analyze()
	Analyze() InlineVariables_InlineVarAnalysis

	// ReanalyzeAfterAliasedVar is called for a Var that is an alias after the
	// original variable is handled.
	//
	// aliasedVar is the variable that was aliased; aliasedVarAnalysis the
	// analysis results for the aliased variable.
	//
	// Java: InlineVariables.VarExpert.reanalyzeAfterAliasedVar(Var, InlineVarAnalysis)
	ReanalyzeAfterAliasedVar(aliasedVar *Var, aliasedVarAnalysis InlineVariables_InlineVarAnalysis) InlineVariables_InlineVarAnalysis
}

// InlineVariables_VarExpertBase holds VarExpert's base implementations.
type InlineVariables_VarExpertBase struct{}

// Java: InlineVariables.VarExpert.reanalyzeAfterAliasedVar(Var, InlineVarAnalysis)
func (b *InlineVariables_VarExpertBase) ReanalyzeAfterAliasedVar(
	aliasedVar *Var, aliasedVarAnalysis InlineVariables_InlineVarAnalysis) InlineVariables_InlineVarAnalysis {
	panic(jutil.NewUnsupportedOperationExceptionMsg("not waiting for an aliased variable"))
}

// inlineVariables_noInlineSelfOrAliasesExpert is the anonymous VarExpert of
// NO_INLINE_SELF_OR_ALIASES_EXPERT.
type inlineVariables_noInlineSelfOrAliasesExpert struct {
	InlineVariables_VarExpertBase
}

// Java: InlineVariables.NO_INLINE_SELF_OR_ALIASES_EXPERT.analyze()
func (*inlineVariables_noInlineSelfOrAliasesExpert) Analyze() InlineVariables_InlineVarAnalysis {
	return InlineVariables_NO_INLINE_SELF_OR_ALIASES_ANALYSIS
}

// inlineVariables_noInlineSelfAliasesOkExpert is the anonymous VarExpert of
// NO_INLINE_SELF_ALIASES_OK_EXPERT.
type inlineVariables_noInlineSelfAliasesOkExpert struct {
	InlineVariables_VarExpertBase
}

// Java: InlineVariables.NO_INLINE_SELF_ALIASES_OK_EXPERT.analyze()
func (*inlineVariables_noInlineSelfAliasesOkExpert) Analyze() InlineVariables_InlineVarAnalysis {
	return InlineVariables_NO_INLINE_SELF_ALIASES_OK_ANALYSIS
}

// Canonical `VarExpert` objects for common cases.
var (
	// Java: InlineVariables.NO_INLINE_SELF_OR_ALIASES_EXPERT
	InlineVariables_NO_INLINE_SELF_OR_ALIASES_EXPERT InlineVariables_VarExpert = &inlineVariables_noInlineSelfOrAliasesExpert{}
	// Java: InlineVariables.NO_INLINE_SELF_ALIASES_OK_EXPERT
	InlineVariables_NO_INLINE_SELF_ALIASES_OK_EXPERT InlineVariables_VarExpert = &inlineVariables_noInlineSelfAliasesOkExpert{}
)

// InlineVariables_InlineVarAnalysis is the result of analyzing a variable to
// see if it can be inlined.
//
// Non-abstract classes should be created by extending this class and
// overriding the methods that should return `true` or perform an operation.
type InlineVariables_InlineVarAnalysis interface {
	// ShouldInline: should we inline this variable?
	//
	// Mutually exclusive with `shouldWaitForAliasedVar()`.
	//
	// Java: InlineVariables.InlineVarAnalysis.shouldInline()
	ShouldInline() bool

	// ShouldWaitForAliasedVar: true if this variable is an alias.
	//
	// The caller should wait for the aliased variable to be handled, then
	// call the expert's `reanalyzeAfterAliasedVar()` method.
	//
	// Mutually exclusive with `shouldInline()`
	//
	// Java: InlineVariables.InlineVarAnalysis.shouldWaitForAliasedVar()
	ShouldWaitForAliasedVar() bool

	// GetAliasedVar gets the aliased variable, if this is an alias.
	//
	// Returns the `Var` for which this one is an alias; throws
	// `UnsupportedOperationException` if `shouldWaitForAliasedVar()` is
	// `false`.
	//
	// Java: InlineVariables.InlineVarAnalysis.getAliasedVar()
	GetAliasedVar() *Var

	// IsSafeToInlineAliases: true if it is safe for aliases of this variable
	// to inline this variable's name.
	//
	// Java: InlineVariables.InlineVarAnalysis.isSafeToInlineAliases()
	IsSafeToInlineAliases() bool

	// PerformInline performs the inline operation.
	//
	// Java: InlineVariables.InlineVarAnalysis.performInline()
	PerformInline()
}

// InlineVariables_InlineVarAnalysisBase holds InlineVarAnalysis's base
// implementations.
type InlineVariables_InlineVarAnalysisBase struct{}

// Java: InlineVariables.InlineVarAnalysis.shouldInline()
func (b *InlineVariables_InlineVarAnalysisBase) ShouldInline() bool {
	return false
}

// Java: InlineVariables.InlineVarAnalysis.shouldWaitForAliasedVar()
func (b *InlineVariables_InlineVarAnalysisBase) ShouldWaitForAliasedVar() bool {
	return false
}

// Java: InlineVariables.InlineVarAnalysis.getAliasedVar()
func (b *InlineVariables_InlineVarAnalysisBase) GetAliasedVar() *Var {
	panic(jutil.NewUnsupportedOperationExceptionMsg("no aliased Var"))
}

// Java: InlineVariables.InlineVarAnalysis.isSafeToInlineAliases()
func (b *InlineVariables_InlineVarAnalysisBase) IsSafeToInlineAliases() bool {
	return false
}

// Java: InlineVariables.InlineVarAnalysis.performInline()
func (b *InlineVariables_InlineVarAnalysisBase) PerformInline() {
	panic(jutil.NewUnsupportedOperationExceptionMsg("cannot inline"))
}

// inlineVariables_noInlineSelfOrAliasesAnalysis is the anonymous
// InlineVarAnalysis of NO_INLINE_SELF_OR_ALIASES_ANALYSIS.
type inlineVariables_noInlineSelfOrAliasesAnalysis struct {
	InlineVariables_InlineVarAnalysisBase
}

// inlineVariables_noInlineSelfAliasesOkAnalysis is the anonymous
// InlineVarAnalysis of NO_INLINE_SELF_ALIASES_OK_ANALYSIS.
type inlineVariables_noInlineSelfAliasesOkAnalysis struct {
	InlineVariables_InlineVarAnalysisBase
}

// Java: InlineVariables.NO_INLINE_SELF_ALIASES_OK_ANALYSIS.isSafeToInlineAliases()
func (*inlineVariables_noInlineSelfAliasesOkAnalysis) IsSafeToInlineAliases() bool {
	return true
}

// Canonical `InlineVarAnalysis` values.
// We'll use these instead of creating new objects for each analysis.
var (
	// Java: InlineVariables.NO_INLINE_SELF_OR_ALIASES_ANALYSIS
	InlineVariables_NO_INLINE_SELF_OR_ALIASES_ANALYSIS InlineVariables_InlineVarAnalysis = &inlineVariables_noInlineSelfOrAliasesAnalysis{}
	// Java: InlineVariables.NO_INLINE_SELF_ALIASES_OK_ANALYSIS
	InlineVariables_NO_INLINE_SELF_ALIASES_OK_ANALYSIS InlineVariables_InlineVarAnalysis = &inlineVariables_noInlineSelfAliasesOkAnalysis{}
)

type handledVarAnalysisEntry struct {
	v        *Var
	analysis InlineVariables_InlineVarAnalysis
}

// InlineVariables_InliningBehavior builds up information about nodes in each
// scope. When exiting the scope, inspects all variables in that scope, and
// inlines any that we can.
type InlineVariables_InliningBehavior struct {
	outer            *InlineVariables
	codingConvention CodingConvention

	// Records the analyses of variables that have already been handled in
	// the current scope.
	//
	// This is necessary in order for aliases of those variables to determine
	// whether they may be inlined.
	//
	// Go-only: language-adapter — Var overrides equals/hashCode (ScopedName),
	// hence the Object map.
	currentScopeHandledVarAnalysesMap *jutil.ObjectLinkedHashMap[*Var, InlineVariables_InlineVarAnalysis]
	currentScope                      *Scope
	handledScopeVars                  []InlineVariables_InlineVarAnalysis
	handledScopeVarsBuf               [64]InlineVariables_InlineVarAnalysis
	handledFallbackVars               []handledVarAnalysisEntry
	handledFallbackBuf                [8]handledVarAnalysisEntry

	// Records alias variables that are waiting for the original variables to
	// be handled.
	//
	// The value is an object that knows what to do when the original variable
	// is handled.
	varToAliasRetryHandlersMap *jutil.SetMultimap[*Var, *InlineVariables_InliningBehavior_AliasInlineRetryHandler]
	aliasRetryHandlers         map[*Var][]*InlineVariables_InliningBehavior_AliasInlineRetryHandler
	unpinnedRetries            int

	scratchExpert        InlineVariables_InliningBehavior_StandardVarExpert
	delayedNegSlab       []inlineVariables_delayedNegativeAnalysis
	delayedNegChunks     [][]inlineVariables_delayedNegativeAnalysis
	delayedNegChunkIdx   int
	delayedNegBuf        [32]inlineVariables_delayedNegativeAnalysis
	retryExpertSlab      []InlineVariables_InliningBehavior_StandardVarExpert
	retryExpertChunks    [][]InlineVariables_InliningBehavior_StandardVarExpert
	retryExpertChunkIdx  int
	retryExpertBuf       [8]InlineVariables_InliningBehavior_StandardVarExpert
	retryHandlerSlab     []InlineVariables_InliningBehavior_AliasInlineRetryHandler
	retryHandlerChunks   [][]InlineVariables_InliningBehavior_AliasInlineRetryHandler
	retryHandlerChunkIdx int
	retryHandlerBuf      [8]InlineVariables_InliningBehavior_AliasInlineRetryHandler
	motionScratch        NodeIterators_LocalVarMotion
}

var inliningBehaviorPool = sync.Pool{
	New: func() any {
		b := &InlineVariables_InliningBehavior{}
		b.handledScopeVars = b.handledScopeVarsBuf[:0]
		b.handledFallbackVars = b.handledFallbackBuf[:0]
		b.delayedNegSlab = b.delayedNegBuf[:0]
		b.retryExpertSlab = b.retryExpertBuf[:0]
		b.retryHandlerSlab = b.retryHandlerBuf[:0]
		return b
	},
}

func acquireInliningBehavior(outer *InlineVariables) *InlineVariables_InliningBehavior {
	b := inliningBehaviorPool.Get().(*InlineVariables_InliningBehavior)
	b.outer = outer
	if outer != nil && outer.compiler != nil {
		b.codingConvention = outer.compiler.GetCodingConvention()
	} else {
		b.codingConvention = nil
	}
	return b
}

func (b *InlineVariables_InliningBehavior) release() {
	b.outer = nil
	b.codingConvention = nil
	b.currentScope = nil
	if len(b.handledScopeVars) > 0 {
		clear(b.handledScopeVars)
		b.handledScopeVars = b.handledScopeVars[:0]
	}
	if len(b.handledFallbackVars) > 0 {
		clear(b.handledFallbackVars)
		b.handledFallbackVars = b.handledFallbackVars[:0]
	}
	if len(b.aliasRetryHandlers) > 0 {
		clear(b.aliasRetryHandlers)
	}
	b.unpinnedRetries = 0
	b.scratchExpert = InlineVariables_InliningBehavior_StandardVarExpert{}
	clear(b.delayedNegBuf[:])
	for i := range b.delayedNegChunks {
		clear(b.delayedNegChunks[i][:cap(b.delayedNegChunks[i])])
	}
	b.delayedNegChunkIdx = 0
	b.delayedNegSlab = b.delayedNegBuf[:0]
	clear(b.retryExpertBuf[:])
	for i := range b.retryExpertChunks {
		clear(b.retryExpertChunks[i][:cap(b.retryExpertChunks[i])])
	}
	b.retryExpertChunkIdx = 0
	b.retryExpertSlab = b.retryExpertBuf[:0]
	clear(b.retryHandlerBuf[:])
	for i := range b.retryHandlerChunks {
		clear(b.retryHandlerChunks[i][:cap(b.retryHandlerChunks[i])])
	}
	b.retryHandlerChunkIdx = 0
	b.retryHandlerSlab = b.retryHandlerBuf[:0]
	b.motionScratch.compiler = nil
	b.motionScratch.varName = ""
	b.motionScratch.lookAhead = nil
	clear(b.motionScratch.iteratorStorage.ancestors)
	b.motionScratch.iteratorStorage.ancestors = b.motionScratch.iteratorStorage.ancestors[:0]
	inliningBehaviorPool.Put(b)
}

// Java: InlineVariables.InliningBehavior()
func NewInlineVariables_InliningBehavior(outer *InlineVariables) *InlineVariables_InliningBehavior {
	var cc CodingConvention
	if outer != nil && outer.compiler != nil {
		cc = outer.compiler.GetCodingConvention()
	}
	b := &InlineVariables_InliningBehavior{
		outer:            outer,
		codingConvention: cc,
	}
	b.handledScopeVars = b.handledScopeVarsBuf[:0]
	b.handledFallbackVars = b.handledFallbackBuf[:0]
	b.delayedNegSlab = b.delayedNegBuf[:0]
	b.retryExpertSlab = b.retryExpertBuf[:0]
	b.retryHandlerSlab = b.retryHandlerBuf[:0]
	return b
}

func (b *InlineVariables_InliningBehavior) putHandledVarAnalysis(v *Var, analysis InlineVariables_InlineVarAnalysis) {
	scope := b.currentScope
	if v == nil || scope == nil || v.scope != scope {
		return
	}
	if idx := v.GetIndex(); uint(idx) < uint(len(b.handledScopeVars)) && scope.varList[idx] == v {
		b.handledScopeVars[idx] = analysis
		return
	}
	b.handledFallbackVars = append(b.handledFallbackVars, handledVarAnalysisEntry{v: v, analysis: analysis})
}

func (b *InlineVariables_InliningBehavior) getHandledVarAnalysis(scope *Scope, aliasedVar *Var) InlineVariables_InlineVarAnalysis {
	if aliasedVar == nil || aliasedVar.scope != scope {
		return nil
	}
	if idx := aliasedVar.GetIndex(); uint(idx) < uint(len(b.handledScopeVars)) && scope.varList[idx] == aliasedVar {
		return b.handledScopeVars[idx]
	}
	if canon := scope.lookupOwnVarByID(aliasedVar.nameID); canon != nil {
		if idx := canon.GetIndex(); uint(idx) < uint(len(b.handledScopeVars)) && scope.varList[idx] == canon {
			return b.handledScopeVars[idx]
		}
	}
	nameID := aliasedVar.nameID
	for _, entry := range b.handledFallbackVars {
		if entry.v == aliasedVar || entry.v.nameID == nameID {
			return entry.analysis
		}
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
		declaration:                         referenceInfo.references.Get(0),
		isDeclaredOrInferredConstant:        isDeclaredOrInferredConstant,
		mayBeAParameterModifiedViaArguments: mayBeAParameterModifiedViaArguments,
	}
	return &b.scratchExpert
}

func (b *InlineVariables_InliningBehavior) allocDelayedNegativeAnalysis(
	src *InlineVariables_InliningBehavior_StandardVarExpert,
) *inlineVariables_delayedNegativeAnalysis {
	if len(b.delayedNegSlab) == cap(b.delayedNegSlab) {
		if b.delayedNegChunkIdx < len(b.delayedNegChunks) {
			ch := b.delayedNegChunks[b.delayedNegChunkIdx]
			b.delayedNegChunkIdx++
			b.delayedNegSlab = ch[:0:cap(ch)]
		} else {
			newCap := 64
			if cap(b.delayedNegSlab) >= 64 {
				newCap = min(cap(b.delayedNegSlab)*2, 512)
			}
			ch := make([]inlineVariables_delayedNegativeAnalysis, 0, newCap)
			b.delayedNegChunks = append(b.delayedNegChunks, ch)
			b.delayedNegChunkIdx = len(b.delayedNegChunks)
			b.delayedNegSlab = ch
		}
	}
	idx := len(b.delayedNegSlab)
	b.delayedNegSlab = b.delayedNegSlab[:idx+1]
	a := &b.delayedNegSlab[idx]
	var state uint8
	if src.isWellDefined.IsKnownToBe(true) {
		state = 3
	}
	*a = inlineVariables_delayedNegativeAnalysis{
		referenceInfo: src.referenceInfo,
		state:         state,
	}
	return a
}

func (b *InlineVariables_InliningBehavior) allocRetryExpert(
	src *InlineVariables_InliningBehavior_StandardVarExpert,
) *InlineVariables_InliningBehavior_StandardVarExpert {
	if len(b.retryExpertSlab) == cap(b.retryExpertSlab) {
		if b.retryExpertChunkIdx < len(b.retryExpertChunks) {
			ch := b.retryExpertChunks[b.retryExpertChunkIdx]
			b.retryExpertChunkIdx++
			b.retryExpertSlab = ch[:0:cap(ch)]
		} else {
			newCap := 32
			if cap(b.retryExpertSlab) >= 32 {
				newCap = min(cap(b.retryExpertSlab)*2, 512)
			}
			ch := make([]InlineVariables_InliningBehavior_StandardVarExpert, 0, newCap)
			b.retryExpertChunks = append(b.retryExpertChunks, ch)
			b.retryExpertChunkIdx = len(b.retryExpertChunks)
			b.retryExpertSlab = ch
		}
	}
	idx := len(b.retryExpertSlab)
	b.retryExpertSlab = b.retryExpertSlab[:idx+1]
	e := &b.retryExpertSlab[idx]
	*e = *src
	e.initValAnalysis.outer = e
	return e
}

func (b *InlineVariables_InliningBehavior) allocRetryHandler(
	v *Var,
	expert InlineVariables_VarExpert,
) *InlineVariables_InliningBehavior_AliasInlineRetryHandler {
	if len(b.retryHandlerSlab) == cap(b.retryHandlerSlab) {
		if b.retryHandlerChunkIdx < len(b.retryHandlerChunks) {
			ch := b.retryHandlerChunks[b.retryHandlerChunkIdx]
			b.retryHandlerChunkIdx++
			b.retryHandlerSlab = ch[:0:cap(ch)]
		} else {
			newCap := 32
			if cap(b.retryHandlerSlab) >= 32 {
				newCap = min(cap(b.retryHandlerSlab)*2, 512)
			}
			ch := make([]InlineVariables_InliningBehavior_AliasInlineRetryHandler, 0, newCap)
			b.retryHandlerChunks = append(b.retryHandlerChunks, ch)
			b.retryHandlerChunkIdx = len(b.retryHandlerChunks)
			b.retryHandlerSlab = ch
		}
	}
	idx := len(b.retryHandlerSlab)
	b.retryHandlerSlab = b.retryHandlerSlab[:idx+1]
	h := &b.retryHandlerSlab[idx]
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
	// so we can clear this map.
	if len(b.handledScopeVars) > 0 {
		clear(b.handledScopeVars)
		b.handledScopeVars = b.handledScopeVars[:0]
	}
	if len(b.handledFallbackVars) > 0 {
		clear(b.handledFallbackVars)
		b.handledFallbackVars = b.handledFallbackVars[:0]
	}
	b.currentScope = nil
	b.delayedNegChunkIdx = 0
	b.delayedNegSlab = b.delayedNegBuf[:0]
	if len(b.aliasRetryHandlers) == 0 {
		b.retryExpertChunkIdx = 0
		b.retryExpertSlab = b.retryExpertBuf[:0]
		b.retryHandlerChunkIdx = 0
		b.retryHandlerSlab = b.retryHandlerBuf[:0]
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
		if cap(b.handledScopeVars) < nVars {
			newCap := max(len(b.handledScopeVarsBuf), cap(b.handledScopeVars)*2)
			if newCap < nVars {
				newCap = nVars
			}
			b.handledScopeVars = make([]InlineVariables_InlineVarAnalysis, nVars, newCap)
		} else {
			b.handledScopeVars = b.handledScopeVars[:nVars]
		}
	}
	mayBeAParameterModifiedViaArguments :=
		b.VarsInThisScopeMayBeModifiedUsingArguments(scope, referenceMap)
	wrapper, _ := referenceMap.(*ReferenceCollector_ReferenceMapWrapper)
	var collector *ReferenceCollector
	if wrapper != nil {
		collector = wrapper.c
	}
	for i := 0; i < len(scope.varList); i++ {
		v := scope.varList[i]
		if v == nil {
			continue
		}
		if b.outer.mode == InlineVariables_Mode_CONSTANTS_ONLY && !inlineVariables_isDeclaredOrInferredConst(v) {
			continue
		}
		var referenceInfo *ReferenceCollection
		if collector != nil {
			referenceInfo = collector.GetReferences(v)
		} else if wrapper != nil {
			referenceInfo = wrapper.GetReferences(v)
		} else {
			referenceInfo = referenceMap.GetReferences(v)
		}
		expert := b.createVarExpertForScope(v, referenceInfo, mayBeAParameterModifiedViaArguments, isLocalScope)
		analysis := expert.Analyze()
		if analysis.ShouldWaitForAliasedVar() {
			// We need to wait for the aliased variable to be handled.
			aliasedVar := analysis.GetAliasedVar()
			aliasedVarAnalysis := b.getHandledVarAnalysis(scope, aliasedVar)
			if aliasedVarAnalysis != nil {
				// We've actually already completed the aliased Var
				b.handleAliasedVarCompletion(v, expert, aliasedVar, aliasedVarAnalysis)
			} else {
				// wait for completion of the aliased Var
				if expert == &b.scratchExpert {
					expert = b.allocRetryExpert(&b.scratchExpert)
				}
				retryHandler := b.allocRetryHandler(v, expert)
				b.putAliasRetryHandler(aliasedVar, retryHandler)
			}
		} else {
			b.completeVarAnalysis(v, analysis)
		}
	}
}

func (b *InlineVariables_InliningBehavior) completeVarAnalysis(v *Var, analysis InlineVariables_InlineVarAnalysis) {
	if analysis.ShouldInline() {
		analysis.PerformInline()
		if analysis == &b.scratchExpert {
			analysis = InlineVariables_NO_INLINE_SELF_ALIASES_OK_ANALYSIS
		}
	} else if analysis == &b.scratchExpert {
		analysis = b.allocDelayedNegativeAnalysis(&b.scratchExpert)
	}
	if b.outer.mode != InlineVariables_Mode_CONSTANTS_ONLY {
		// Record the results for aliases we may find later
		b.putHandledVarAnalysis(v, analysis)
		// Retry any aliases we saw before handling this Var
		b.RetryAliases(v, analysis)
	}
}

func (b *InlineVariables_InliningBehavior) handleAliasedVarCompletion(
	v *Var, expert InlineVariables_VarExpert, aliasedVar *Var, aliasedVarAnalysis InlineVariables_InlineVarAnalysis) {
	newAnalysis := expert.ReanalyzeAfterAliasedVar(aliasedVar, aliasedVarAnalysis)
	jutil.CheckStateMsg1(
		!newAnalysis.ShouldWaitForAliasedVar(), "expert for %s asked to wait a second time", v)
	b.completeVarAnalysis(v, newAnalysis)
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
		if refs != nil && !refs.references.IsEmpty() {
			for _, ref := range refs.references.Slice() {
				if !b.IsSafeUseOfArguments(ref.GetNode()) {
					return true
				}
			}
		}
	}
	return false
}

// Java: InlineVariables.InliningBehavior.retryAliases(Var, InlineVarAnalysis)
func (b *InlineVariables_InliningBehavior) RetryAliases(aliasedVar *Var, aliasedVarAnalysis InlineVariables_InlineVarAnalysis) {
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
		aliasInlineRetryHandler.HandleAliasedVarCompletion(aliasedVar, aliasedVarAnalysis)
	}
}

// InlineVariables_InliningBehavior_AliasInlineRetryHandler retries inlining
// an alias variable after the original variable has been dealt with.
type InlineVariables_InliningBehavior_AliasInlineRetryHandler struct {
	outer *InlineVariables_InliningBehavior

	v      *Var
	expert InlineVariables_VarExpert
}

// Java: InlineVariables.InliningBehavior.AliasInlineRetryHandler(Var, VarExpert)
func NewInlineVariables_InliningBehavior_AliasInlineRetryHandler(
	outer *InlineVariables_InliningBehavior, v *Var, expert InlineVariables_VarExpert) *InlineVariables_InliningBehavior_AliasInlineRetryHandler {
	return &InlineVariables_InliningBehavior_AliasInlineRetryHandler{outer: outer, v: v, expert: expert}
}

// Java: InlineVariables.InliningBehavior.AliasInlineRetryHandler.handleAliasedVarCompletion(Var, InlineVarAnalysis)
func (h *InlineVariables_InliningBehavior_AliasInlineRetryHandler) HandleAliasedVarCompletion(
	aliasedVar *Var, aliasedVarAnalyis InlineVariables_InlineVarAnalysis) {
	h.outer.handleAliasedVarCompletion(h.v, h.expert, aliasedVar, aliasedVarAnalyis)
}

type standardVarAnalysisKind uint8

const (
	analysisKindNone standardVarAnalysisKind = iota
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
	InlineVariables_VarExpertBase
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
	analysisInitValue     *rhino.Node
	analysisInitRef       *Reference
	analysisSingleReadRef *Reference
}

// Java: InlineVariables.InliningBehavior.StandardVarExpert(VarExpertInitData)
func NewInlineVariables_InliningBehavior_StandardVarExpert(
	outer *InlineVariables_InliningBehavior, initData *InlineVariables_InliningBehavior_VarExpertInitData) *InlineVariables_InliningBehavior_StandardVarExpert {
	return &InlineVariables_InliningBehavior_StandardVarExpert{
		outer:                               outer,
		v:                                   initData.v,
		referenceInfo:                       initData.referenceInfo,
		declaration:                         initData.referenceInfo.references.Get(0),
		isDeclaredOrInferredConstant:        initData.isDeclaredOrInferredConstant,
		mayBeAParameterModifiedViaArguments: initData.mayBeAParameterModifiedViaArguments,
	}
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

func (e *InlineVariables_InliningBehavior_StandardVarExpert) IsSafeToInlineAliases() bool {
	switch e.analysisKind {
	case analysisKindDelayedNegative:
		return e.IsWellDefinedAssignedOnce()
	case analysisKindInlineWellDefined, analysisKindInlineNeverAssigned, analysisKindInlineSingleRef:
		return true
	case analysisKindWaitForAlias:
		panic(jutil.NewUnsupportedOperationExceptionMsg("analysis is incomplete"))
	default:
		return false
	}
}

func (e *InlineVariables_InliningBehavior_StandardVarExpert) PerformInline() {
	switch e.analysisKind {
	case analysisKindInlineWellDefined:
		e.outer.InlineWellDefinedVariable(e.v, e.analysisInitValue, &e.referenceInfo.references)
	case analysisKindInlineNeverAssigned:
		srcLocation := e.declaration.GetNode()
		undefinedNode := NodeUtil_NewUndefinedNode(srcLocation)
		e.outer.InlineWellDefinedVariable(e.v, undefinedNode, &e.referenceInfo.references)
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
			for _, ref := range e.referenceInfo.references.Slice() {
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
func (e *InlineVariables_InliningBehavior_StandardVarExpert) Analyze() InlineVariables_InlineVarAnalysis {
	if InlineVariables_HasNoInlineAnnotation(e.v) || e.IsReferencedWeakly() {
		return e.GetNegativeInlineVarAnalysis()
	}
	if e.referenceInfo.references.Size() == 1 && !e.isDeclaredOrInferredConstant {
		if e.outer.outer.mode == InlineVariables_Mode_CONSTANTS_ONLY {
			return e.GetNegativeInlineVarAnalysis()
		}
		if firstChild := e.declaration.nameNode.GetFirstChild(); firstChild == nil || !firstChild.IsName() {
			return e.GetNegativeInlineVarAnalysis()
		}
	}
	initialization := e.GetInitialization()
	var initValue *rhino.Node
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
		return e.inlineWellDefinedVarAnalysis(initValue)
	}

	if e.outer.outer.mode == InlineVariables_Mode_CONSTANTS_ONLY {
		// If we're in constants-only mode, don't run more aggressive
		// inlining heuristics. See InlineConstantsTest.
		return e.GetNegativeInlineVarAnalysis()
	}

	if initialValueAnalysis.IsAlias() {
		e.analysisKind = analysisKindWaitForAlias
		e.analysisAliasedVar = initialValueAnalysis.GetAliasedVar()
		return e
	}

	return e.AnalyzeWithInitialValue(initialization, initValue, initialValueAnalysis)
}

// Java: InlineVariables.InliningBehavior.StandardVarExpert.analyzeWithInitialValue(Reference, Node, InitialValueAnalysis)
func (e *InlineVariables_InliningBehavior_StandardVarExpert) AnalyzeWithInitialValue(
	initialization *Reference, initValue *rhino.Node, initialValueAnalysis *InlineVariables_InliningBehavior_StandardVarExpert_InitialValueAnalysis) InlineVariables_InlineVarAnalysis {
	refCount := e.referenceInfo.references.Size()

	// TODO(bradfordcsmith): We could remove the `refCount > 1` here, but:
	//  1. That will require some additional logic to handle stuff like named function
	//     expressions and avoiding the removal of side-effects.
	//  2. RemoveUnusedCode will remove those cases anyway.
	//  3. The unit tests will have to be more verbose to make sure stuff we don't want to be
	//  removed is used.
	if refCount > 1 && e.AllReferencesAreValid() {
		if e.IsNeverAssigned() {
			e.analysisKind = analysisKindInlineNeverAssigned
			return e
		}
		if e.IsWellDefined() &&
			(initialValueAnalysis.IsImmutableValue() ||
				(initValue.IsThis() && !e.referenceInfo.IsThisRebound())) {
			// if the variable is referenced more than once, we can only
			// inline it if it's immutable and never defined before referenced.
			return e.inlineWellDefinedVarAnalysis(initValue)
		}

		firstReadRefIndex := 2
		if e.declaration == initialization {
			firstReadRefIndex = 1
		}
		numReadRefs := refCount - firstReadRefIndex
		if numReadRefs == 0 {
			// The only reference is the initialization.
			// Remove the assignment and the variable declaration.
			return e.inlineWellDefinedVarAnalysis(initValue)
		}

		if numReadRefs == 1 {
			// The variable is likely only read once, so we can try some more complex inlining
			// heuristics.
			singleReadReference := e.referenceInfo.references.Get(firstReadRefIndex)
			if e.outer.CanInline(e.declaration, initialization, singleReadReference, initValue) {
				// A custom inline method is needed for this case.
				e.analysisKind = analysisKindInlineSingleRef
				e.analysisInitRef = initialization
				e.analysisSingleReadRef = singleReadReference
				return e
			}
		}
	}

	return e.GetNegativeInlineVarAnalysis()
}

// Java: InlineVariables.InliningBehavior.StandardVarExpert.getNegativeInlineVarAnalysis()
func (e *InlineVariables_InliningBehavior_StandardVarExpert) GetNegativeInlineVarAnalysis() InlineVariables_InlineVarAnalysis {

	// If we already know whether it is safe to inline aliases of this variable,
	// use canonical analysis values to save on memory space.
	if e.mayBeAParameterModifiedViaArguments ||
		e.outer.outer.mode == InlineVariables_Mode_CONSTANTS_ONLY ||
		e.referenceInfo.references.Size() <= 1 ||
		e.isWellDefinedAssignedOnce.IsKnownToBe(false) ||
		e.isWellDefined.IsKnownToBe(false) ||
		e.isAssignedOnceInLifetime.IsKnownToBe(false) {
		return InlineVariables_NO_INLINE_SELF_OR_ALIASES_ANALYSIS
	}
	if e.isWellDefinedAssignedOnce.IsKnownToBe(true) ||
		(e.isWellDefined.IsKnownToBe(true) && e.isAssignedOnceInLifetime.IsKnownToBe(true)) {
		return InlineVariables_NO_INLINE_SELF_ALIASES_OK_ANALYSIS
	}

	// Delay calculating safety until we're actually asked.
	e.analysisKind = analysisKindDelayedNegative
	return e
}

// inlineVariables_delayedNegativeAnalysis is the anonymous InlineVarAnalysis
// of StandardVarExpert.getNegativeInlineVarAnalysis().
type inlineVariables_delayedNegativeAnalysis struct {
	InlineVariables_InlineVarAnalysisBase
	referenceInfo *ReferenceCollection
	state         uint8 // 0: unknown, 1: false, 2: true, 3: wellDefined known true
}

// Java: InlineVariables.InliningBehavior.StandardVarExpert.getNegativeInlineVarAnalysis().new InlineVarAnalysis() {...}.isSafeToInlineAliases()
func (a *inlineVariables_delayedNegativeAnalysis) IsSafeToInlineAliases() bool {
	switch a.state {
	case 0:
		if a.referenceInfo.IsWellDefined() && a.referenceInfo.IsAssignedOnceInLifetime() {
			a.state = 2
			return true
		}
		a.state = 1
		return false
	case 3:
		if a.referenceInfo.IsAssignedOnceInLifetime() {
			a.state = 2
			return true
		}
		a.state = 1
		return false
	default:
		return a.state == 2
	}
}

func (e *InlineVariables_InliningBehavior_StandardVarExpert) inlineWellDefinedVarAnalysis(initValue *rhino.Node) InlineVariables_InlineVarAnalysis {
	e.analysisKind = analysisKindInlineWellDefined
	e.analysisInitValue = initValue
	return e
}

// Java: InlineVariables.InliningBehavior.StandardVarExpert.createInlineWellDefinedVarAnalysis(Node)
func (e *InlineVariables_InliningBehavior_StandardVarExpert) CreateInlineWellDefinedVarAnalysis(initValue *rhino.Node) *InlineVariables_InliningBehavior_PositiveInlineVarAnalysis {
	return NewInlineVariables_InliningBehavior_PositiveInlineVarAnalysis(
		func() { e.outer.InlineWellDefinedVariable(e.v, initValue, &e.referenceInfo.references) })
}

// Java: InlineVariables.InliningBehavior.StandardVarExpert.reanalyzeAfterAliasedVar(Var, InlineVarAnalysis)
func (e *InlineVariables_InliningBehavior_StandardVarExpert) ReanalyzeAfterAliasedVar(
	aliasedVar *Var, aliasedVarAnalysis InlineVariables_InlineVarAnalysis) InlineVariables_InlineVarAnalysis {
	initialization := e.GetInitialization()
	var initValue *rhino.Node
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
		return e.inlineWellDefinedVarAnalysis(initValue)
	}

	if initialValueAnalysis.IsAlias() &&
		aliasedVar == initialValueAnalysis.GetAliasedVar() &&
		aliasedVarAnalysis.IsSafeToInlineAliases() &&
		e.IsWellDefinedAssignedOnce() {
		// The variable we aliased couldn't be inlined itself, or it was an alias for another
		// variable that got inlined in its place.
		// However, it is safe to inline the name assigned to this variable now.
		return e.inlineWellDefinedVarAnalysis(initValue)
	}
	return e.AnalyzeWithInitialValue(initialization, initValue, initialValueAnalysis)
}

// InlineVariables_InliningBehavior_StandardVarExpert_InitialValueAnalysis is
// information about a value used to initialize a variable.
type InlineVariables_InliningBehavior_StandardVarExpert_InitialValueAnalysis struct {
	outer *InlineVariables_InliningBehavior_StandardVarExpert

	value *rhino.Node

	isImmutableValue InlineVariables_InitiallyUnknown[bool]
	aliasedVar       InlineVariables_InitiallyUnknown[*Var]
}

// Java: InlineVariables.InliningBehavior.StandardVarExpert.InitialValueAnalysis(Node)
func NewInlineVariables_InliningBehavior_StandardVarExpert_InitialValueAnalysis(
	outer *InlineVariables_InliningBehavior_StandardVarExpert, value *rhino.Node) *InlineVariables_InliningBehavior_StandardVarExpert_InitialValueAnalysis {
	return &InlineVariables_InliningBehavior_StandardVarExpert_InitialValueAnalysis{
		outer: outer,
		value: value,
	}
}

// Java: InlineVariables.InliningBehavior.StandardVarExpert.InitialValueAnalysis.isImmutableValue()
func (a *InlineVariables_InliningBehavior_StandardVarExpert_InitialValueAnalysis) IsImmutableValue() bool {
	if a.isImmutableValue.IsKnown() {
		return a.isImmutableValue.GetKnownValue()
	} else {
		return a.isImmutableValue.SetKnownValueOnce(
			a.value != nil && NodeUtil_IsImmutableValue(a.value))
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
		if a.value != nil && a.value.IsName() {
			aliasedName := a.value.GetString()
			// Never consider this variable to be an alias of itself.
			if aliasedName == a.outer.v.GetName() {
				return a.aliasedVar.SetKnownValueOnce(nil)
			}
			return a.aliasedVar.SetKnownValueOnce(a.outer.v.GetScope().(*Scope).GetVar(aliasedName))
		}
		return a.aliasedVar.SetKnownValueOnce(nil)
	}
}

// InlineVariables_InliningBehavior_PositiveInlineVarAnalysis indicates that
// the analyzed variable may be inlined.
type InlineVariables_InliningBehavior_PositiveInlineVarAnalysis struct {
	InlineVariables_InlineVarAnalysisBase

	inliner func()
}

// Java: InlineVariables.InliningBehavior.PositiveInlineVarAnalysis(Runnable)
func NewInlineVariables_InliningBehavior_PositiveInlineVarAnalysis(inliner func()) *InlineVariables_InliningBehavior_PositiveInlineVarAnalysis {
	return &InlineVariables_InliningBehavior_PositiveInlineVarAnalysis{inliner: inliner}
}

// Java: InlineVariables.InliningBehavior.PositiveInlineVarAnalysis.shouldInline()
func (a *InlineVariables_InliningBehavior_PositiveInlineVarAnalysis) ShouldInline() bool {
	return true
}

// Java: InlineVariables.InliningBehavior.PositiveInlineVarAnalysis.performInline()
func (a *InlineVariables_InliningBehavior_PositiveInlineVarAnalysis) PerformInline() {
	a.inliner()
}

// Java: InlineVariables.InliningBehavior.PositiveInlineVarAnalysis.isSafeToInlineAliases()
func (a *InlineVariables_InliningBehavior_PositiveInlineVarAnalysis) IsSafeToInlineAliases() bool {
	// If we've inlined this variable, then any aliases of it should definitely try to
	// inline themselves with the new value that replaced this variable.
	// If for some reason the logic that requested this analysis decided not to actually
	// do the inlining, it's still safe to inline the name of this variable.
	// If it weren't, we wouldn't have said it was OK to inline its value.
	return true
}

// InlineVariables_InliningBehavior_VarIsAliasAnalysis indicates that the
// analyzed variable is an alias and the decision about whether to inline it
// must wait until inlining has been done (or not) for the original value it
// aliases.
type InlineVariables_InliningBehavior_VarIsAliasAnalysis struct {
	InlineVariables_InlineVarAnalysisBase

	aliasedVar *Var
}

// Java: InlineVariables.InliningBehavior.VarIsAliasAnalysis(Var)
func NewInlineVariables_InliningBehavior_VarIsAliasAnalysis(aliasedVar *Var) *InlineVariables_InliningBehavior_VarIsAliasAnalysis {
	return &InlineVariables_InliningBehavior_VarIsAliasAnalysis{aliasedVar: aliasedVar}
}

// Java: InlineVariables.InliningBehavior.VarIsAliasAnalysis.shouldWaitForAliasedVar()
func (a *InlineVariables_InliningBehavior_VarIsAliasAnalysis) ShouldWaitForAliasedVar() bool {
	return true
}

// Java: InlineVariables.InliningBehavior.VarIsAliasAnalysis.getAliasedVar()
func (a *InlineVariables_InliningBehavior_VarIsAliasAnalysis) GetAliasedVar() *Var {
	return a.aliasedVar
}

// Java: InlineVariables.InliningBehavior.VarIsAliasAnalysis.isSafeToInlineAliases()
func (a *InlineVariables_InliningBehavior_VarIsAliasAnalysis) IsSafeToInlineAliases() bool {
	panic(jutil.NewUnsupportedOperationExceptionMsg("analysis is incomplete"))
}

// CreateVarExpert creates a VarExpert object appropriate for the given
// variable.
//
// Java: InlineVariables.InliningBehavior.createVarExpert(Var, ReferenceCollection, boolean)
func (b *InlineVariables_InliningBehavior) CreateVarExpert(
	v *Var, referenceInfo *ReferenceCollection, mayBeAParameterModifiedViaArguments bool) InlineVariables_VarExpert {
	return b.createVarExpertForScope(v, referenceInfo, mayBeAParameterModifiedViaArguments, v.IsLocal())
}

func (b *InlineVariables_InliningBehavior) createVarExpertForScope(
	v *Var, referenceInfo *ReferenceCollection, mayBeAParameterModifiedViaArguments bool, isLocal bool) InlineVariables_VarExpert {
	if referenceInfo == nil {
		// If we couldn't collect any reference info, don't try to inline and assume it's unsafe
		// to inline aliases of this variable, too.
		return InlineVariables_NO_INLINE_SELF_OR_ALIASES_EXPERT
	}

	if v.IsExtern() {
		// TODO(bradfordcsmith): Extern variables are generally unsafe to inline.
		return InlineVariables_NO_INLINE_SELF_OR_ALIASES_EXPERT
	}

	isDeclaredOrInferredConstant := inlineVariables_isDeclaredOrInferredConst(v)
	if !isDeclaredOrInferredConstant && b.outer.mode == InlineVariables_Mode_CONSTANTS_ONLY {
		// If we're only inlining constants, then we shouldn't inline an alias.
		return InlineVariables_NO_INLINE_SELF_OR_ALIASES_EXPERT
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
		return InlineVariables_NO_INLINE_SELF_OR_ALIASES_EXPERT
	}

	if cc.IsPropertyRenameFunction(v.nameNode) {
		// It's not terribly likely that anything creates an alias of our special property rename
		// function, but its value never changes, so it should be safe to inline aliases to it.
		return InlineVariables_NO_INLINE_SELF_ALIASES_OK_EXPERT
	}

	return b.initScratchVarExpert(v, referenceInfo, isDeclaredOrInferredConstant, mayBeAParameterModifiedViaArguments)
}

// InlineVariables_InliningBehavior_VarExpertInitData is used to initialize
// fields in a `StandardVarExpert` object.
type InlineVariables_InliningBehavior_VarExpertInitData struct {
	v                                   *Var
	referenceInfo                       *ReferenceCollection
	isDeclaredOrInferredConstant        bool
	mayBeAParameterModifiedViaArguments bool
}

// Java: InlineVariables.InliningBehavior.VarExpertInitData()
func NewInlineVariables_InliningBehavior_VarExpertInitData() *InlineVariables_InliningBehavior_VarExpertInitData {
	return &InlineVariables_InliningBehavior_VarExpertInitData{}
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
func (b *InlineVariables_InliningBehavior) IsSafeUseOfArguments(argumentsNode *rhino.Node) bool {
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
func (b *InlineVariables_InliningBehavior) IsTargetOfPropertyRead(n *rhino.Node) bool {
	getNode := n.GetParent()
	return n.IsFirstChildOf(getNode) &&
		NodeUtil_IsNormalOrOptChainGet(getNode) &&
		!NodeUtil_IsLValue(getNode)
}

// IsSecondArgumentToDotApplyMethod: true if `n` is being used in a call like
// this: `fn.apply(thisArg, n)`.
//
// Java: InlineVariables.InliningBehavior.isSecondArgumentToDotApplyMethod(Node)
func (b *InlineVariables_InliningBehavior) IsSecondArgumentToDotApplyMethod(n *rhino.Node) bool {
	callNode := n.GetParent()
	if NodeUtil_IsNormalOrOptChainCall(callNode) {
		calleeNode := callNode.GetFirstChild()
		if NodeUtil_IsNormalOrOptChainGetProp(calleeNode) {
			if calleeNode.GetString() == "apply" {
				thisArgNode := calleeNode.GetNext()
				if thisArgNode != nil {
					return thisArgNode.GetNext() == n
				}
			}
		}
	}
	return false
}

// Inline does the actual work of inlining a single declaration into a
// single reference.
//
// Java: InlineVariables.InliningBehavior.inline(Reference, Reference, Reference)
func (b *InlineVariables_InliningBehavior) Inline(decl *Reference, init *Reference, ref *Reference) {
	value := init.GetAssignedValue()
	jutil.CheckState(value != nil)
	// Check for function declarations before the value is moved in the AST.
	isFunctionDeclaration := NodeUtil_IsFunctionDeclaration(value)
	if isFunctionDeclaration {
		// We're inlining this function declaration because there was only one reference to it,
		// probably the rhs of an assignment statement. Since we're eliminating the only
		// reference, we can get rid of it. See b/277538101.
		funcName := value.GetFirstChild()
		jutil.CheckNotNull(funcName)
		jutil.CheckState(funcName.IsName())
		funcName.SetString("")

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
func (b *InlineVariables_InliningBehavior) InlineWellDefinedVariable(v *Var, value *rhino.Node, refSet *jutil.ArrayList[*Reference]) {
	for _, r := range refSet.Slice() {
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
			clonedValue := value.CloneTree()
			NodeUtil_MarkNewScopesChanged(clonedValue, b.outer.compiler)
			b.InlineValue(r.GetNode(), clonedValue)
		}
	}
}

// RemoveDeclaration removes the given VAR declaration.
//
// Java: InlineVariables.InliningBehavior.removeDeclaration(Reference)
func (b *InlineVariables_InliningBehavior) RemoveDeclaration(decl *Reference) {
	varNode := decl.GetParent()
	jutil.CheckStateMsg1(NodeUtil_IsNameDeclaration(varNode), "%s", varNode)
	grandparent := decl.GetGrandparent()

	b.outer.compiler.ReportChangeToEnclosingScope(decl.GetNode())
	decl.GetNode().Detach()
	// Remove var node if empty
	if !varNode.HasChildren() {
		NodeUtil_RemoveChild(grandparent, varNode)
	}
}

// Java: InlineVariables.InliningBehavior.inlineValue(Node, Node)
func (b *InlineVariables_InliningBehavior) InlineValue(toRemove *rhino.Node, toInsert *rhino.Node) {
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
	declaration *Reference, initialization *Reference, reference *Reference, initValue *rhino.Node) bool {
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
	jutil.CheckState(initValue != nil)
	if initValue.IsGetProp() &&
		reference.GetParent().IsCall() &&
		reference.GetParent().GetFirstChild() == reference.GetNode() {
		return false
	}

	if initValue.IsFunction() {
		callNode := reference.GetParent()
		if reference.GetParent().IsCall() {
			convention := b.outer.compiler.GetCodingConvention()
			// Bug 2388531: Don't inline subclass definitions into class defining
			// calls as this confused class removing logic.
			relationship := convention.GetClassesDefinedByCall(callNode)
			if relationship != nil {
				return false
			}

			// issue 668: Don't inline singleton getter methods
			// calls as this confused class removing logic.
			if convention.GetSingletonGetterClassName(callNode) != nil {
				return false
			}
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
func (b *InlineVariables_InliningBehavior) CanMoveAggressively(value *rhino.Node) bool {
	// Function expressions and other mutable objects can move within
	// the same basic block.
	return NodeUtil_IsLiteralValue(value, true) || value.IsFunction()
}

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
	for it.HasNext() {
		curNode := it.Next()
		if curNode == targetName {
			return true
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
		parent := initialization.GetParent()
		jutil.CheckState(parent.IsAssign() && parent.GetFirstChild() == initialization.GetNode())
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
	if n == nil {
		return false
	}
	if info := n.GetJSDocInfo(); info != nil {
		return info.IsNoInline()
	}
	parent := n.GetParent()
	if parent == nil {
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
			if gp := parent.GetParent(); gp != nil && gp.IsExport() {
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
