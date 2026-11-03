// Port of com/google/javascript/jscomp/InlineFunctions.java @ closure-compiler v20261028
// Port of com/google/javascript/jscomp/InlineCostEstimator.java @ closure-compiler v20261028

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
	"math"
	"os"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

type fnDeclKind uint8

const (
	fnKindNone fnDeclKind = iota
	fnKindDecl            // function foo() {}
	fnKindVar             // var/let/const foo = function() {}
	fnKindExpr            // (function() {})()
)

// inlineFnTarget describes a candidate function definition without interface
// boxing or per-function heap wrapper allocations.
type inlineFnTarget struct {
	kind     fnDeclKind
	name     string
	nameNode rhino.Node
	fnNode   rhino.Node
	declNode rhino.Node
}

func (t inlineFnTarget) declaringBlock() rhino.Node {
	if t.declNode.IsNil() {
		return rhino.NilNode
	}
	return t.declNode.GetParent()
}

func (t inlineFnTarget) remove(compiler AbstractCompiler) {
	if t.declNode.IsNil() {
		return
	}
	compiler.ReportChangeToEnclosingScope(t.declNode)
	NodeUtil_RemoveChild(t.declNode.GetParent(), t.declNode)
	NodeUtil_MarkFunctionsDeleted(t.declNode, compiler)
}

// InlineFunctions_FunctionState is used to track the decisions that have
// been made about a function.
type InlineFunctions_FunctionState struct {
	target                 inlineFnTarget
	safeFnNode             rhino.Node
	inline                 bool
	remove                 bool
	inlineDirectly         bool
	referencesThis         bool
	hasInnerClosures       bool
	summaryValid           bool
	requireInliningState   uint8 // 0=uninitialized, 1=false, 2=true
	encourageInliningState uint8 // 0=uninitialized, 1=false, 2=true
	refs                   smallMap[rhino.Node, *FunctionInjector_Reference]
	chunk                  *JSChunk
	summary                calleeSummary
	bodyAstSize            int32
	aliasCount             int32
}

// Java: InlineFunctions.FunctionState()
func NewInlineFunctions_FunctionState() *InlineFunctions_FunctionState {
	return &InlineFunctions_FunctionState{
		safeFnNode:     rhino.NilNode,
		inline:         true,
		remove:         true,
		inlineDirectly: false,
		referencesThis: false,
	}
}

// Java: InlineFunctions.FunctionState.hasExistingFunctionDefinition()
func (s *InlineFunctions_FunctionState) HasExistingFunctionDefinition() bool {
	return s.target.kind != fnKindNone
}

// Java: InlineFunctions.FunctionState.setReferencesThis(boolean)
func (s *InlineFunctions_FunctionState) SetReferencesThis(referencesThis bool) {
	s.referencesThis = referencesThis
}

// Java: InlineFunctions.FunctionState.getReferencesThis()
func (s *InlineFunctions_FunctionState) GetReferencesThis() bool {
	return s.referencesThis
}

// SetHasInnerClosures: whether the function contains inner functions or
// non-static class field initializers.
//
// Java: InlineFunctions.FunctionState.setHasInnerClosures(boolean)
func (s *InlineFunctions_FunctionState) SetHasInnerClosures(hasInnerClosures bool) {
	s.hasInnerClosures = hasInnerClosures
}

// Java: InlineFunctions.FunctionState.hasInnerClosures()
func (s *InlineFunctions_FunctionState) HasInnerClosures() bool {
	return s.hasInnerClosures
}

// Java: InlineFunctions.FunctionState.removeBlockInliningReferences()
func (s *InlineFunctions_FunctionState) RemoveBlockInliningReferences() {
	s.refs.retainFunc(func(_ rhino.Node, r *FunctionInjector_Reference) bool {
		return r.mode != FunctionInjector_InliningMode_BLOCK
	})
}

// Java: InlineFunctions.FunctionState.hasBlockInliningReferences()
func (s *InlineFunctions_FunctionState) HasBlockInliningReferences() bool {
	for _, r := range s.refs.values() {
		if r.mode == FunctionInjector_InliningMode_BLOCK {
			return true
		}
	}
	return false
}

func (s *InlineFunctions_FunctionState) setTarget(target inlineFnTarget) {
	jutil.CheckState(s.target.kind == fnKindNone)
	s.target = target
}

func (s *InlineFunctions_FunctionState) functionNode() rhino.Node {
	return s.target.fnNode
}

func (s *InlineFunctions_FunctionState) functionName() string {
	return s.target.name
}

// Java: InlineFunctions.FunctionState.getSafeFnNode()
func (s *InlineFunctions_FunctionState) GetSafeFnNode() rhino.Node {
	if !s.safeFnNode.IsNil() {
		return s.safeFnNode
	}
	return s.target.fnNode
}

// Java: InlineFunctions.FunctionState.setSafeFnNode(Node)
func (s *InlineFunctions_FunctionState) SetSafeFnNode(safeFnNode rhino.Node) {
	if s.summaryValid {
		if s.summary.fnNode == s.GetSafeFnNode() {
			s.summary.fnNode = safeFnNode
		} else {
			s.summaryValid = false
		}
	}
	s.safeFnNode = safeFnNode
}

func (s *InlineFunctions_FunctionState) readOnlySummary() *calleeSummary {
	if s.summaryValid && s.summary.fnNode == s.GetSafeFnNode() {
		return &s.summary
	}
	return nil
}

func (s *InlineFunctions_FunctionState) ensureSummary() *calleeSummary {
	fnNode := s.GetSafeFnNode()
	if s.summaryValid && s.summary.fnNode == fnNode {
		return &s.summary
	}
	if scanCalleeSummary(s.target.name, fnNode, &s.summary) {
		s.summaryValid = true
		return &s.summary
	}
	s.summaryValid = false
	return nil
}

// Java: InlineFunctions.FunctionState.canInline()
func (s *InlineFunctions_FunctionState) CanInline() bool {
	return s.inline
}

func (s *InlineFunctions_FunctionState) lazyAnnotationFlag(state *uint8, check func(rhino.Node) bool) bool {
	if s.target.kind == fnKindNone {
		return false
	}
	if *state == 0 {
		if check(s.target.fnNode) {
			*state = 2
		} else {
			*state = 1
		}
	}
	return *state == 2
}

// Java: InlineFunctions.FunctionState.encourageInlining()
func (s *InlineFunctions_FunctionState) EncourageInlining() bool {
	return s.lazyAnnotationFlag(&s.encourageInliningState, InlineFunctions_HasEncourageInliningAnnotation)
}

// Java: InlineFunctions.FunctionState.requireInlining()
func (s *InlineFunctions_FunctionState) RequireInlining() bool {
	return s.lazyAnnotationFlag(&s.requireInliningState, InlineFunctions_HasRequireInliningAnnotation)
}

// Java: InlineFunctions.FunctionState.disallowInlining(AbstractCompiler, DisallowInliningReason)
func (s *InlineFunctions_FunctionState) DisallowInlining(compiler AbstractCompiler, reason InlineFunctions_DisallowInliningReason) {
	s.DisallowInliningWithWarning(compiler, reason, InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_YES)
}

// Java: InlineFunctions.FunctionState.disallowInlining(AbstractCompiler, DisallowInliningReason, ShouldWarnWhenRequireInliningCannotInline)
func (s *InlineFunctions_FunctionState) DisallowInliningWithWarning(
	compiler AbstractCompiler,
	reason InlineFunctions_DisallowInliningReason,
	shouldWarnWhenRequireInliningCannotInline InlineFunctions_ShouldWarnWhenRequireInliningCannotInline) {
	if s.RequireInlining() &&
		compiler.IsDebugLoggingEnabled() &&
		shouldWarnWhenRequireInliningCannotInline ==
			InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_YES {
		compiler.Report(
			JSError_MakeWithNode(s.target.fnNode, InlineFunctions_MISSED_REQUIRED_INLINING, reason.String()))
	}
	s.inline = false

	// No need to keep references to function that can't be inlined.
	s.refs.release()
	s.summary = calleeSummary{}
	s.summaryValid = false
	// Don't remove functions that we aren't inlining.
	s.remove = false
}

// Java: InlineFunctions.FunctionState.canRemove()
func (s *InlineFunctions_FunctionState) CanRemove() bool {
	return s.remove
}

// Java: InlineFunctions.FunctionState.setRemove(boolean, AbstractCompiler, CannotRemoveReason, ShouldWarnWhenRequireInliningCannotInline)
func (s *InlineFunctions_FunctionState) SetRemove(
	remove bool,
	compiler AbstractCompiler,
	reason InlineFunctions_CannotRemoveReason,
	shouldWarnIfRequireInlining InlineFunctions_ShouldWarnWhenRequireInliningCannotInline) {
	s.SetRemoveWithContext(remove, compiler, reason, shouldWarnIfRequireInlining /* contextNode= */, rhino.NilNode)
}

// Java: InlineFunctions.FunctionState.setRemove(boolean, AbstractCompiler, CannotRemoveReason, ShouldWarnWhenRequireInliningCannotInline, Node)
func (s *InlineFunctions_FunctionState) SetRemoveWithContext(
	remove bool,
	compiler AbstractCompiler,
	reason InlineFunctions_CannotRemoveReason,
	shouldWarnIfRequireInlining InlineFunctions_ShouldWarnWhenRequireInliningCannotInline,
	contextNode rhino.Node) {
	if !remove &&
		(shouldWarnIfRequireInlining != InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_NO) &&
		s.RequireInlining() &&
		compiler.IsDebugLoggingEnabled() {
		contextString := "unknown"
		if !contextNode.IsNil() {
			contextString = contextNode.String()
		}
		compiler.Report(
			JSError_MakeWithNode(
				s.target.fnNode,
				InlineFunctions_MISSED_REQUIRED_INLINING,
				reason.String(),
				contextString))
	}
	s.remove = remove
}

// Java: InlineFunctions.FunctionState.canInlineDirectly()
func (s *InlineFunctions_FunctionState) CanInlineDirectly() bool {
	return s.inlineDirectly
}

// Java: InlineFunctions.FunctionState.inlineDirectly(boolean)
func (s *InlineFunctions_FunctionState) InlineDirectly(directReplacement bool) {
	s.inlineDirectly = directReplacement
}

// Java: InlineFunctions.FunctionState.hasReferences()
func (s *InlineFunctions_FunctionState) HasReferences() bool {
	return !s.refs.isEmpty()
}

// Java: InlineFunctions.FunctionState.addReference(Reference)
func (s *InlineFunctions_FunctionState) AddReference(ref *FunctionInjector_Reference) {
	s.refs.put(ref.callNode, ref)
}

// GetReferences returns the references (the values of the references map).
//
// Java: InlineFunctions.FunctionState.getReferences()
func (s *InlineFunctions_FunctionState) GetReferences() []*FunctionInjector_Reference {
	return s.refs.values()
}

// Java: InlineFunctions.FunctionState.getReference(Node)
func (s *InlineFunctions_FunctionState) GetReference(n rhino.Node) *FunctionInjector_Reference {
	return s.refs.getOrZero(n)
}

// Java: InlineFunctions.FunctionState.setChunk(JSChunk)
func (s *InlineFunctions_FunctionState) SetChunk(chunk *JSChunk) {
	s.chunk = chunk
}

// Java: InlineFunctions.FunctionState.getChunk()
func (s *InlineFunctions_FunctionState) GetChunk() *JSChunk {
	return s.chunk
}

// InlineFunctions_ShouldWarnWhenRequireInliningCannotInline is the enum
// InlineFunctions.ShouldWarnWhenRequireInliningCannotInline.
type InlineFunctions_ShouldWarnWhenRequireInliningCannotInline int32

const (
	InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_YES InlineFunctions_ShouldWarnWhenRequireInliningCannotInline = iota
	InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_NO
)

var inlineFunctions_ShouldWarnWhenRequireInliningCannotInline_names = [...]string{"YES", "NO"}

// Java: InlineFunctions.ShouldWarnWhenRequireInliningCannotInline.name()
func (e InlineFunctions_ShouldWarnWhenRequireInliningCannotInline) String() string {
	return inlineFunctions_ShouldWarnWhenRequireInliningCannotInline_names[e]
}

// InlineFunctions_CannotRemoveReason is the enum
// InlineFunctions.CannotRemoveReason.
type InlineFunctions_CannotRemoveReason int32

const (
	InlineFunctions_CannotRemoveReason_BLOCK_INLINING_REFERENCE InlineFunctions_CannotRemoveReason = iota
	InlineFunctions_CannotRemoveReason_CALLED_BY_INLINED_FUNCTION
	InlineFunctions_CannotRemoveReason_UNINLINABLE_REFERENCE
	InlineFunctions_CannotRemoveReason_UNREMOVABLE_REFERENCE
)

var inlineFunctions_CannotRemoveReason_names = [...]string{
	"BLOCK_INLINING_REFERENCE",
	"CALLED_BY_INLINED_FUNCTION",
	"UNINLINABLE_REFERENCE",
	"UNREMOVABLE_REFERENCE",
}

// Java: InlineFunctions.CannotRemoveReason.name()
func (e InlineFunctions_CannotRemoveReason) String() string {
	return inlineFunctions_CannotRemoveReason_names[e]
}

// InlineFunctions_DisallowInliningReason is the enum
// InlineFunctions.DisallowInliningReason.
type InlineFunctions_DisallowInliningReason int32

const (
	InlineFunctions_DisallowInliningReason_ASYNC_FUNCTION InlineFunctions_DisallowInliningReason = iota
	InlineFunctions_DisallowInliningReason_CAPTURED_VARIABLE_NOT_IN_SCOPE
	InlineFunctions_DisallowInliningReason_COST
	InlineFunctions_DisallowInliningReason_GENERATOR_FUNCTION
	InlineFunctions_DisallowInliningReason_HAS_LOCAL_NAMES
	InlineFunctions_DisallowInliningReason_MULTIPLE_FUNCTION_DEFINITIONS
	InlineFunctions_DisallowInliningReason_NON_INLINABLE_PARAM
	InlineFunctions_DisallowInliningReason_NOT_CANDIDATE_FUNCTION
	InlineFunctions_DisallowInliningReason_NO_INLINE_ANNOTATION
	InlineFunctions_DisallowInliningReason_REASSIGNED
	InlineFunctions_DisallowInliningReason_TARGET_SIZE_EXCEEDS_LIMIT
)

var inlineFunctions_DisallowInliningReason_names = [...]string{
	"ASYNC_FUNCTION",
	"CAPTURED_VARIABLE_NOT_IN_SCOPE",
	"COST",
	"GENERATOR_FUNCTION",
	"HAS_LOCAL_NAMES",
	"MULTIPLE_FUNCTION_DEFINITIONS",
	"NON_INLINABLE_PARAM",
	"NOT_CANDIDATE_FUNCTION",
	"NO_INLINE_ANNOTATION",
	"REASSIGNED",
	"TARGET_SIZE_EXCEEDS_LIMIT",
}

// Java: InlineFunctions.DisallowInliningReason.name()
func (e InlineFunctions_DisallowInliningReason) String() string {
	return inlineFunctions_DisallowInliningReason_names[e]
}

// InlineFunctions_ShouldRequireInlining is the record
// InlineFunctions.ShouldRequireInlining(boolean shouldValidate, String
// optMessage).
type InlineFunctions_ShouldRequireInlining struct {
	shouldValidate bool
	optMessage     string
}

// Java: InlineFunctions.ShouldRequireInlining.shouldValidate()
func (r InlineFunctions_ShouldRequireInlining) ShouldValidate() bool { return r.shouldValidate }

// Java: InlineFunctions.ShouldRequireInlining.optMessage()
func (r InlineFunctions_ShouldRequireInlining) OptMessage() string { return r.optMessage }

// Java: InlineFunctions.checkOptimizationLevelSupportsRequireInlining(CompilerOptions)
func InlineFunctions_CheckOptimizationLevelSupportsRequireInlining(
	options *CompilerOptions) InlineFunctions_ShouldRequireInlining {
	// Non-optimized code might not perform all required inlinings.
	var problems []string
	if !options.ShouldOptimize() {
		problems = append(problems, "non-optimized code")
	}
	// Required inlinings don't make sense without inlining. If we don't inline
	// properties we might not be able to inline functions defined on namespaces.
	if !options.ShouldInlineProperties() {
		problems = append(problems, "missing shouldInlineProperties")
	}
	if !options.ShouldInlineConstantVars() {
		problems = append(problems, "missing inlineConstantVars")
	}
	if !options.GetSmartNameRemoval() {
		problems = append(problems, "missing smartNameRemoval")
	}
	if !options.GetInlineFunctionsLevel().IsOn() {
		problems = append(problems, "missing inlineFunctions")
	}
	if !options.GetInlineFunctionsLevel().IncludesGlobals() {
		problems = append(problems, "inlineFunctions does not include globals")
	}
	// Coverage runs may break some inlining.
	if options.GetInstrumentForCoverageOption() != CompilerOptions_InstrumentOption_NONE {
		problems = append(problems, "instrumentForCoverageOption enabled")
	}
	// We need to collapse properties or functions might be retained due to
	// being written on an exported namespace.
	if options.GetPropertyCollapseLevel() != CompilerOptions_PropertyCollapseLevel_ALL {
		problems = append(problems, "missing property collapsing")
	}
	// Avoid inlining failures due to asserts.
	if !options.ShouldRemoveClosureAsserts() {
		problems = append(problems, "missing removeClosureAsserts")
	}
	// Avoid inlining failures due to being unable to devirtualize a method.
	if !options.ShouldDevirtualizeMethods() {
		problems = append(problems, "missing devirtualizeMethods")
	}
	// We need more than one inlining pass for most required inlinings to work. This
	// eliminates 'fast' mode.
	if options.GetMaxOptimizationLoopIterations() == 1 {
		problems = append(problems,
			"optimizationLoopMaxIterations should be >= 2, but is: "+
				strconv.Itoa(options.GetMaxOptimizationLoopIterations()))
	}
	if len(problems) == 0 {
		return InlineFunctions_ShouldRequireInlining{shouldValidate: true}
	}
	return InlineFunctions_ShouldRequireInlining{optMessage: "[" + strings.Join(problems, ", ") + "]"}
}

// -------------------------------------------------------------------------
// InlineCostEstimator
// -------------------------------------------------------------------------

// For now simply assume identifiers are 2 characters.
const InlineCostEstimator_ESTIMATED_IDENTIFIER = "ab"

// Java: InlineCostEstimator.ESTIMATED_IDENTIFIER_COST
const InlineCostEstimator_ESTIMATED_IDENTIFIER_COST = len(InlineCostEstimator_ESTIMATED_IDENTIFIER)

// InlineCostEstimator_GetCost determines the size of the JS code.
//
// Java: InlineCostEstimator.getCost(Node)
func InlineCostEstimator_GetCost(root rhino.Node) int {
	return InlineCostEstimator_GetCostWithThreshold(root, math.MaxInt32)
}

// InlineCostEstimator_GetCostWithThreshold determines the estimated size of
// the JS snippet represented by the node.
//
// Java: InlineCostEstimator.getCost(Node, int)
func InlineCostEstimator_GetCostWithThreshold(root rhino.Node, costThreshold int) int {
	var cg CodeGenerator
	cg.initForCostEstimation(costThreshold)
	cg.AddNode(root)
	return cg.cc.cost
}

// InlineCostEstimator_CompiledSizeEstimator estimates compiled size by
// assuming names are shortened and all whitespace is stripped.
type InlineCostEstimator_CompiledSizeEstimator struct {
	cg CodeGenerator
}

// Java: InlineCostEstimator.CompiledSizeEstimator(int)
func NewInlineCostEstimator_CompiledSizeEstimator(costThreshold int) *InlineCostEstimator_CompiledSizeEstimator {
	e := &InlineCostEstimator_CompiledSizeEstimator{}
	e.cg.initForCostEstimation(costThreshold)
	return e
}

// AddNode is the Java add(Node) overload.
//
// Java: InlineCostEstimator.CompiledSizeEstimator.add(Node)
func (e *InlineCostEstimator_CompiledSizeEstimator) AddNode(root rhino.Node) {
	e.cg.AddNode(root)
}

// Java: InlineCostEstimator.CompiledSizeEstimator.add(String)
func (e *InlineCostEstimator_CompiledSizeEstimator) Add(newCode string) {
	e.cg.cc.Add(newCode)
}

// Java: InlineCostEstimator.CompiledSizeEstimator.getCost()
func (e *InlineCostEstimator_CompiledSizeEstimator) GetCost() int {
	return e.cg.cc.cost
}

// Java: InlineCostEstimator.CompiledSizeEstimator.continueProcessing()
func (e *InlineCostEstimator_CompiledSizeEstimator) ContinueProcessing() bool {
	return e.cg.cc.ContinueProcessing()
}

// Java: InlineCostEstimator.CompiledSizeEstimator.getLastChar()
func (e *InlineCostEstimator_CompiledSizeEstimator) GetLastChar() uint16 {
	return e.cg.cc.GetLastChar()
}

// -------------------------------------------------------------------------
// Go-only: Compositional Per-FUNCTION Subtree Change-Stamp Cache
// -------------------------------------------------------------------------

var (
	inlineFnCacheEnv    = os.Getenv("JSGO_INLINE_FN_CACHE")
	inlineFnCacheVerify = os.Getenv("JSGO_INLINE_FN_CACHE_VERIFY") == "1"
	// inlineFnCacheForce enables the cache outside PhaseOptimizer loops (tests).
	inlineFnCacheForce = false

	inlineFnCacheCandidateHits   atomic.Int64
	inlineFnCacheCandidateMisses atomic.Int64
	inlineFnCacheLeafWalkSkips   atomic.Int64
	inlineFnCacheSideEffectHits  atomic.Int64
	inlineFnCacheFnInfoSeeds     atomic.Int64
	inlineFnCacheCostHits        atomic.Int64
)

const (
	fnCostUnknown uint8 = iota
	fnCostLowerBound
	fnCostExact
)

func inlineFnCacheEnabled() bool {
	return inlineFnCacheEnv != "0"
}

type inlineSideEffectConfig struct {
	typeRegistry                 *rhino.JSTypeRegistry
	accessorSummary              *AccessorSummary
	useTypesForLocalOptimization bool
	assumeGettersArePure         bool
	hasRegexpGlobalReferences    bool
	assumeKnownBuiltinsArePure   bool
}

func currentInlineSideEffectConfig(compiler AbstractCompiler) inlineSideEffectConfig {
	a := compiler.GetAstAnalyzer()
	return inlineSideEffectConfig{
		typeRegistry:                 a.typeRegistry,
		accessorSummary:              a.accessorSummary,
		useTypesForLocalOptimization: a.useTypesForLocalOptimization,
		assumeGettersArePure:         a.assumeGettersArePure,
		hasRegexpGlobalReferences:    a.hasRegexpGlobalReferences,
		assumeKnownBuiltinsArePure:   a.assumeKnownBuiltinsArePure,
	}
}

const (
	fnCacheShallowReady uint8 = 1 << iota
	fnCacheHasDirectInnerFns
	fnCacheAlwaysInlinable
	fnCacheDirectReplace
	fnCacheIsGenerator
	fnCacheIsAsync
	fnCacheHasNonInlinableParam
	fnCacheShallowLocalClosureDisallowed
)

// inlineFnCacheEntry is strictly pointer-free (noscan, 44 bytes) so c.arena
// adds zero GC mark/scan work regardless of function count. Candidate
// calleeSummary records live in c.summaries (indexed by 1-based summaryIdx).
type inlineFnCacheEntry struct {
	ownChangeTime    int32
	lastPassID       uint32
	seedPassID       uint32
	innerOff         uint32
	innerLen         uint32
	summaryIdx       uint32
	fullAstSize      int32
	fullAstSizeLimit int32
	costVal          int32
	fnInnerNode      rhino.Node
	flags            uint8
	evalState        uint8
	fnInnerState     uint8
	costState        uint8
}

func (e inlineFnCacheEntry) hasFlag(f uint8) bool { return e.flags&f != 0 }

type inlineFnCandSummary struct {
	summaryFnNameID uint32
	seCfgGen        uint32
	summary         calleeSummary
	sideEffects     calleeSideEffectSummary
}

type inlineFunctionsCache struct {
	root      rhino.Node
	compacts  int
	phaseOpt  *PhaseOptimizer
	passID    uint32
	seCfg     inlineSideEffectConfig
	seCfgGen  uint32
	index     map[rhino.Node]uint32 // 1-based index into arena (value size 4B: zero map-value heap allocs)
	arena     []inlineFnCacheEntry  // pointer-free noscan slab
	inners    []rhino.Node          // pointer-free noscan slab
	summaries []inlineFnCandSummary // 1-based index from entry.summaryIdx (candidates only)

	// Reusable scratch slices across passes within a loop.
	serialScratch   inlineCandidateScratch
	workerScratches []inlineCandidateScratch
	allCandsBuf     []*InlineFunctions_FunctionState
}

func inlineFnCacheFor(compiler AbstractCompiler, root rhino.Node) *inlineFunctionsCache {
	if !inlineFnCacheEnabled() || root.IsNil() {
		return nil
	}
	c := compiler
	if c == nil || c.parallelRecorder != nil || c.changeTracker == nil {
		return nil
	}
	if !c.GetLifeCycleStage().IsNormalized() {
		return nil
	}
	phaseOpt := c.phaseOptimizer
	if phaseOpt == nil && !inlineFnCacheForce {
		return nil
	}
	cache := c.inlineFnCache
	if cache == nil {
		cache = &inlineFunctionsCache{
			root:      root,
			compacts:  c.astCompact.count,
			phaseOpt:  phaseOpt,
			index:     make(map[rhino.Node]uint32, 256),
			arena:     make([]inlineFnCacheEntry, 1, 256),
			summaries: make([]inlineFnCandSummary, 1, 64),
		}
		c.inlineFnCache = cache
		return cache
	}
	if cache.root != root ||
		cache.compacts != c.astCompact.count ||
		len(cache.arena) == 0 ||
		(!inlineFnCacheForce && cache.phaseOpt != phaseOpt) {
		cache.root = root
		cache.compacts = c.astCompact.count
		cache.phaseOpt = phaseOpt
		clear(cache.index)
		cache.arena = cache.arena[:1]
		cache.inners = cache.inners[:0]
		clear(cache.summaries)
		if cap(cache.summaries) == 0 {
			cache.summaries = make([]inlineFnCandSummary, 1, 64)
		} else {
			cache.summaries = cache.summaries[:1]
		}
	}
	return cache
}

func (c *inlineFunctionsCache) prepareWorkerScratches(numWorkers int, seCfg inlineSideEffectConfig) []inlineCandidateScratch {
	if cap(c.workerScratches) < numWorkers {
		next := make([]inlineCandidateScratch, numWorkers)
		copy(next, c.workerScratches)
		c.workerScratches = next
	} else {
		c.workerScratches = c.workerScratches[:numWorkers]
	}
	for i := range c.workerScratches {
		c.workerScratches[i].resetForWalk(c, seCfg)
	}
	return c.workerScratches
}

func (c *inlineFunctionsCache) lookup(fnNode rhino.Node) (inlineFnCacheEntry, bool) {
	idx := c.index[fnNode]
	if idx == 0 {
		return inlineFnCacheEntry{}, false
	}
	return c.arena[idx], true
}

func (c *inlineFunctionsCache) summaryOf(e inlineFnCacheEntry) *inlineFnCandSummary {
	if e.summaryIdx == 0 {
		return nil
	}
	return &c.summaries[e.summaryIdx]
}

func (c *inlineFunctionsCache) entryMut(fnNode rhino.Node) *inlineFnCacheEntry {
	if c == nil {
		return nil
	}
	idx := c.index[fnNode]
	if idx == 0 {
		return nil
	}
	return &c.arena[idx]
}

func (c *inlineFunctionsCache) getOrAlloc(fnNode rhino.Node) *inlineFnCacheEntry {
	if idx := c.index[fnNode]; idx != 0 {
		return &c.arena[idx]
	}
	idx := uint32(len(c.arena))
	c.arena = append(c.arena, inlineFnCacheEntry{})
	c.index[fnNode] = idx
	return &c.arena[idx]
}

func (c *inlineFunctionsCache) allocOrUpdateSummary(entry *inlineFnCacheEntry, fnNameID uint32, sum calleeSummary) {
	sum.sideEffects = calleeSideEffectSummary{}
	if entry.summaryIdx != 0 {
		c.summaries[entry.summaryIdx] = inlineFnCandSummary{
			summaryFnNameID: fnNameID,
			summary:         sum,
		}
		return
	}
	if len(c.summaries) == 0 {
		c.summaries = make([]inlineFnCandSummary, 1, 64)
	}
	idx := uint32(len(c.summaries))
	c.summaries = append(c.summaries, inlineFnCandSummary{
		summaryFnNameID: fnNameID,
		summary:         sum,
	})
	entry.summaryIdx = idx
}

func (c *inlineFunctionsCache) directInnersOf(e inlineFnCacheEntry) []rhino.Node {
	if e.innerLen == 0 {
		return nil
	}
	return c.inners[e.innerOff : e.innerOff+e.innerLen]
}

func (c *inlineFunctionsCache) appendInners(src []rhino.Node) (uint32, uint32) {
	if len(src) == 0 {
		return 0, 0
	}
	off := uint32(len(c.inners))
	c.inners = append(c.inners, src...)
	return off, uint32(len(src))
}

func (c *inlineFunctionsCache) beginPassMerge(totalLive int, seCfg inlineSideEffectConfig) {
	c.passID++
	if c.passID == 0 {
		c.passID = 1
	}
	if c.seCfgGen == 0 || c.seCfg != seCfg {
		c.seCfg = seCfg
		c.seCfgGen++
		if c.seCfgGen == 0 {
			c.seCfgGen = 1
		}
	}
	if len(c.arena) == 1 && totalLive > 0 {
		if cap(c.arena) < totalLive+64 {
			c.arena = make([]inlineFnCacheEntry, 1, totalLive+64)
		}
		if len(c.index) == 0 {
			c.index = make(map[rhino.Node]uint32, totalLive+totalLive/4+32)
		}
	}
}

func effectiveCalleeFnNameID(target inlineFnTarget) uint32 {
	if target.name == "" {
		return 0
	}
	if !target.nameNode.IsNil() {
		return target.nameNode.GetStringID()
	}
	return uint32(strpool.Lookup(target.name))
}

type inlineFnWalkFrame struct {
	fnNode                    rhino.Node
	innerStart                int
	allInnerSubtreesUnchanged bool
}

type inlineFnStampUpdate struct {
	fnNode        rhino.Node
	ownChangeTime int32
	innerOff      uint32
	innerLen      uint32
	ownUnchanged  bool
	sameInners    bool
}

type inlineFnCandUpdate struct {
	fnNode           rhino.Node
	fullAstSize      int32
	fullAstSizeLimit int32
	summaryIdx       uint32 // 1-based index into scratch.candSummaries (0 = none)
	summaryFnNameID  uint32
	flags            uint8
}

func (c *inlineFunctionsCache) applyScratch(s *inlineCandidateScratch) int {
	passID := c.passID
	for _, fn := range s.unchangedFns {
		if idx := c.index[fn]; idx != 0 {
			c.arena[idx].lastPassID = passID
			c.arena[idx].seedPassID = passID
		}
	}

	for i := range s.stamps {
		st := &s.stamps[i]
		entry := c.getOrAlloc(st.fnNode)
		prevSummaryIdx := entry.summaryIdx
		var innerOff, innerLen uint32
		if st.sameInners {
			innerOff = entry.innerOff
			innerLen = entry.innerLen
		} else if st.innerLen > 0 {
			src := s.stampInners[st.innerOff : st.innerOff+st.innerLen]
			innerOff, innerLen = c.appendInners(src)
		}
		if st.ownUnchanged {
			next := inlineFnCacheEntry{
				ownChangeTime: st.ownChangeTime,
				lastPassID:    passID,
				seedPassID:    passID,
				innerOff:      innerOff,
				innerLen:      innerLen,
				flags:         (entry.flags &^ fnCacheHasDirectInnerFns) | fnCacheHasDirectInnerFns,
				evalState:     entry.evalState,
			}
			if innerLen >= 2 && entry.fnInnerState == fnInnerMultiple {
				next.fnInnerState = fnInnerMultiple
			}
			*entry = next
			continue
		}
		*entry = inlineFnCacheEntry{
			ownChangeTime: st.ownChangeTime,
			lastPassID:    passID,
			innerOff:      innerOff,
			innerLen:      innerLen,
			summaryIdx:    prevSummaryIdx, // keep slot reusable if candUpdates populates it below, else zeroed
			flags:         fnCacheHasDirectInnerFns,
		}
		if prevSummaryIdx != 0 {
			c.summaries[prevSummaryIdx] = inlineFnCandSummary{}
			entry.summaryIdx = 0
		}
	}

	for i := range s.candUpdates {
		u := &s.candUpdates[i]
		entry := c.getOrAlloc(u.fnNode)
		if entry.ownChangeTime == 0 && !entry.hasFlag(fnCacheHasDirectInnerFns) {
			entry.ownChangeTime = u.fnNode.GetChangeTime()
		}
		entry.lastPassID = passID
		entry.flags = (entry.flags & fnCacheHasDirectInnerFns) | u.flags | fnCacheShallowReady
		if u.fullAstSizeLimit > 0 {
			entry.fullAstSize = u.fullAstSize
			entry.fullAstSizeLimit = u.fullAstSizeLimit
		}
		if u.summaryIdx != 0 {
			c.allocOrUpdateSummary(entry, u.summaryFnNameID, s.candSummaries[u.summaryIdx-1])
		}
	}

	if s.candidateHits != 0 {
		inlineFnCacheCandidateHits.Add(s.candidateHits)
	}
	if s.candidateMisses != 0 {
		inlineFnCacheCandidateMisses.Add(s.candidateMisses)
	}
	if s.leafWalkSkips != 0 {
		inlineFnCacheLeafWalkSkips.Add(s.leafWalkSkips)
	}
	if s.sideEffectHits != 0 {
		inlineFnCacheSideEffectHits.Add(s.sideEffectHits)
	}
	return len(s.liveFns)
}

func (c *inlineFunctionsCache) maybePruneDead(totalLive int, scratches []inlineCandidateScratch) {
	if len(c.index) <= totalLive+max(128, totalLive/2) {
		return
	}
	newIndex := make(map[rhino.Node]uint32, totalLive+totalLive/4+32)
	newArena := make([]inlineFnCacheEntry, 1, totalLive+1)
	var newInners []rhino.Node
	if len(c.inners) > 0 {
		newInners = make([]rhino.Node, 0, len(c.inners))
	}
	newSummaries := make([]inlineFnCandSummary, 1, max(16, len(c.summaries)/2))
	for i := range scratches {
		for _, fn := range scratches[i].liveFns {
			oldIdx := c.index[fn]
			if oldIdx == 0 {
				continue
			}
			if _, dup := newIndex[fn]; dup {
				continue
			}
			e := c.arena[oldIdx]
			if e.innerLen > 0 {
				inners := c.inners[e.innerOff : e.innerOff+e.innerLen]
				e.innerOff = uint32(len(newInners))
				newInners = append(newInners, inners...)
			} else {
				e.innerOff = 0
			}
			if e.summaryIdx != 0 {
				sumEntry := c.summaries[e.summaryIdx]
				e.summaryIdx = uint32(len(newSummaries))
				newSummaries = append(newSummaries, sumEntry)
			}
			newIdx := uint32(len(newArena))
			newArena = append(newArena, e)
			newIndex[fn] = newIdx
		}
	}
	c.index = newIndex
	c.arena = newArena
	c.inners = newInners
	c.summaries = newSummaries
}

func (c *inlineFunctionsCache) lookupSeededEval(fnNode rhino.Node) (uint8, bool) {
	if c == nil {
		return 0, false
	}
	idx := c.index[fnNode]
	if idx == 0 {
		return 0, false
	}
	e := c.arena[idx]
	if e.seedPassID != c.passID || e.evalState == 0 {
		return 0, false
	}
	if inlineFnCacheVerify {
		verifySeededFnInfo(fnNode, functionInjectorFnInfo{evalState: e.evalState})
	}
	inlineFnCacheFnInfoSeeds.Add(1)
	return e.evalState, true
}

func (c *inlineFunctionsCache) lookupSeededInner(fnNode rhino.Node) (uint8, rhino.Node, bool) {
	if c == nil {
		return fnInnerUnknown, rhino.NilNode, false
	}
	idx := c.index[fnNode]
	if idx == 0 {
		return fnInnerUnknown, rhino.NilNode, false
	}
	e := c.arena[idx]
	if e.seedPassID != c.passID || e.fnInnerState == fnInnerUnknown {
		return fnInnerUnknown, rhino.NilNode, false
	}
	if inlineFnCacheVerify {
		verifySeededFnInfo(fnNode, functionInjectorFnInfo{innerState: e.fnInnerState, inner: e.fnInnerNode})
	}
	inlineFnCacheFnInfoSeeds.Add(1)
	return e.fnInnerState, e.fnInnerNode, true
}

func (f *InlineFunctions) persistPhase2ToCache(allCands []*InlineFunctions_FunctionState) {
	f.injector.seedCache = nil
	cache := f.activeCache
	if cache == nil {
		return
	}
	seCfgGen := cache.seCfgGen

	for _, fs := range allCands {
		if !fs.summaryValid {
			continue
		}
		fnNode := fs.target.fnNode
		if fs.summary.fnNode != fnNode || !fs.summary.sideEffects.computed {
			continue
		}
		if entry := cache.entryMut(fnNode); entry != nil && entry.summaryIdx != 0 && !fnNode.IsDeleted() && entry.ownChangeTime == fnNode.GetChangeTime() {
			sumEntry := &cache.summaries[entry.summaryIdx]
			if !sumEntry.sideEffects.computed || sumEntry.seCfgGen != seCfgGen {
				sumEntry.sideEffects = fs.summary.sideEffects
				sumEntry.seCfgGen = seCfgGen
			}
		}
	}
	clear(allCands)

	if seMap := f.functionArgumentInjector.sideEffectCache; len(seMap) > 0 {
		// order-independent: updates disjoint per-fnNode side-effect summaries in cache.summaries
		for fnNode, se := range seMap {
			if !se.computed {
				continue
			}
			if entry := cache.entryMut(fnNode); entry != nil && entry.summaryIdx != 0 && !fnNode.IsDeleted() && entry.ownChangeTime == fnNode.GetChangeTime() {
				sumEntry := &cache.summaries[entry.summaryIdx]
				sumEntry.sideEffects = se
				sumEntry.seCfgGen = seCfgGen
			}
		}
	}

	if fnMap := f.injector.fnInfoCache; len(fnMap) > 0 {
		// order-independent: updates disjoint per-fnNode caller fnInfo entries in cache.arena
		for fnNode, info := range fnMap {
			entry := cache.entryMut(fnNode)
			if entry == nil || fnNode.IsDeleted() || entry.ownChangeTime != fnNode.GetChangeTime() {
				continue
			}
			if info.evalState != 0 && entry.evalState != info.evalState {
				entry.evalState = info.evalState
			}
			if info.innerState != fnInnerUnknown &&
				(entry.fnInnerState != info.innerState || entry.fnInnerNode != info.inner) {
				entry.fnInnerState = info.innerState
				entry.fnInnerNode = info.inner
			}
		}
	}
}
