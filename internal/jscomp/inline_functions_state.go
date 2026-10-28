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
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
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

func newNamedFnTarget(fn rhino.Node) inlineFnTarget {
	name, ok := ast.MustFunction(fn).Name()
	jutil.CheckState(ok)
	return inlineFnTarget{
		kind:     fnKindDecl,
		name:     name.Ident(),
		nameNode: name.Node(),
		fnNode:   fn,
		declNode: fn,
	}
}

func newVarFnTarget(declNode rhino.Node) inlineFnTarget {
	d, _ := ast.MustNameDeclaration(declNode).FirstDeclarator()
	name := d.MustName()
	initVal, _ := name.Init()
	return inlineFnTarget{
		kind:     fnKindVar,
		name:     name.Ident(),
		nameNode: name.Node(),
		fnNode:   initVal.MustFunction().Node(),
		declNode: declNode,
	}
}

func newExprFnTarget(fn rhino.Node, index int) inlineFnTarget {
	return inlineFnTarget{
		kind:   fnKindExpr,
		name:   strconv.Itoa(index),
		fnNode: fn,
	}
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
	target           inlineFnTarget
	safeFnNode       rhino.Node
	inline           bool
	remove           bool
	inlineDirectly   bool
	referencesThis   bool
	hasInnerClosures bool
	summaryValid     bool
	references       *jutil.LinkedHashMap[rhino.Node, *FunctionInjector_Reference]
	refMapVal        jutil.LinkedHashMap[rhino.Node, *FunctionInjector_Reference]
	chunk            *JSChunk
	namesToAlias     *jutil.LinkedHashSet[string]
	summary          *calleeSummary

	hasRequireInliningAnnotation            bool
	hasRequireInliningAnnotationInitialized bool

	hasEncourageInliningAnnotation            bool
	hasEncourageInliningAnnotationInitialized bool
	bodyAstSize                               int32
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
	if s.references == nil {
		return
	}
	for key, value := range s.references.All() {
		if value.mode == FunctionInjector_InliningMode_BLOCK {
			s.references.Remove(key)
		}
	}
}

// Java: InlineFunctions.FunctionState.hasBlockInliningReferences()
func (s *InlineFunctions_FunctionState) HasBlockInliningReferences() bool {
	if s.references == nil {
		return false
	}
	for r := range s.references.Values() {
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
	if s.summaryValid && s.summary != nil {
		if s.summary.fnNode == s.GetSafeFnNode() {
			s.summary.fnNode = safeFnNode
		} else {
			s.summaryValid = false
		}
	}
	s.safeFnNode = safeFnNode
}

func (s *InlineFunctions_FunctionState) readOnlySummary() *calleeSummary {
	if s.summaryValid && s.summary != nil && s.summary.fnNode == s.GetSafeFnNode() {
		return s.summary
	}
	return nil
}

func (s *InlineFunctions_FunctionState) ensureSummary() *calleeSummary {
	fnNode := s.GetSafeFnNode()
	if s.summaryValid && s.summary != nil && s.summary.fnNode == fnNode {
		return s.summary
	}
	if s.summary == nil {
		s.summary = new(calleeSummary)
	}
	if scanCalleeSummary(s.target.name, fnNode, s.summary) {
		s.summaryValid = true
		return s.summary
	}
	s.summaryValid = false
	return nil
}

// Java: InlineFunctions.FunctionState.canInline()
func (s *InlineFunctions_FunctionState) CanInline() bool {
	return s.inline
}

// Java: InlineFunctions.FunctionState.encourageInlining()
func (s *InlineFunctions_FunctionState) EncourageInlining() bool {
	if s.target.kind == fnKindNone {
		return false
	}
	if !s.hasEncourageInliningAnnotationInitialized {
		s.hasEncourageInliningAnnotationInitialized = true
		s.hasEncourageInliningAnnotation = InlineFunctions_HasEncourageInliningAnnotation(s.target.fnNode)
	}
	return s.hasEncourageInliningAnnotation
}

// Java: InlineFunctions.FunctionState.requireInlining()
func (s *InlineFunctions_FunctionState) RequireInlining() bool {
	if s.target.kind == fnKindNone {
		return false
	}
	if !s.hasRequireInliningAnnotationInitialized {
		s.hasRequireInliningAnnotationInitialized = true
		s.hasRequireInliningAnnotation = InlineFunctions_HasRequireInliningAnnotation(s.target.fnNode)
	}
	return s.hasRequireInliningAnnotation
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
	if s.refMapVal.Size() > 0 {
		s.refMapVal.Clear()
	}
	s.references = nil
	s.summary = nil
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
	return s.references != nil && !s.references.IsEmpty()
}

// Java: InlineFunctions.FunctionState.getReferencesInternal()
func (s *InlineFunctions_FunctionState) GetReferencesInternal() *jutil.LinkedHashMap[rhino.Node, *FunctionInjector_Reference] {
	if s.references == nil {
		return jutil.NewLinkedHashMap[rhino.Node, *FunctionInjector_Reference]()
	}
	return s.references
}

// Java: InlineFunctions.FunctionState.addReference(Reference)
func (s *InlineFunctions_FunctionState) AddReference(ref *FunctionInjector_Reference) {
	if s.references == nil {
		s.references = &s.refMapVal
	}
	s.references.PutNew(ref.callNode, ref)
}

var inlineFunctions_emptyStringSet = jutil.NewLinkedHashSet[string]()

// GetReferences returns the references (the values of the references map).
//
// Go-only: language-adapter — Java's live values() view becomes a snapshot
// slice; no caller modifies the map while iterating it.
//
// Java: InlineFunctions.FunctionState.getReferences()
func (s *InlineFunctions_FunctionState) GetReferences() []*FunctionInjector_Reference {
	if s.references == nil || s.references.IsEmpty() {
		return nil
	}
	out := make([]*FunctionInjector_Reference, 0, s.references.Size())
	for ref := range s.references.Values() {
		out = append(out, ref)
	}
	return out
}

// Java: InlineFunctions.FunctionState.getReference(Node)
func (s *InlineFunctions_FunctionState) GetReference(n rhino.Node) *FunctionInjector_Reference {
	if s.references == nil {
		return nil
	}
	return s.references.Get(n)
}

// Java: InlineFunctions.FunctionState.getNamesToAlias()
func (s *InlineFunctions_FunctionState) GetNamesToAlias() *jutil.LinkedHashSet[string] {
	if s.namesToAlias == nil {
		return inlineFunctions_emptyStringSet
	}
	return s.namesToAlias
}

// Java: InlineFunctions.FunctionState.setNamesToAlias(Set<String>)
func (s *InlineFunctions_FunctionState) SetNamesToAlias(names *jutil.LinkedHashSet[string]) {
	s.namesToAlias = names
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

// Java: InlineFunctions.ShouldWarnWhenRequireInliningCannotInline.ordinal()
func (e InlineFunctions_ShouldWarnWhenRequireInliningCannotInline) Ordinal() int { return int(e) }

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

// Java: InlineFunctions.CannotRemoveReason.ordinal()
func (e InlineFunctions_CannotRemoveReason) Ordinal() int { return int(e) }

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

// Java: InlineFunctions.DisallowInliningReason.ordinal()
func (e InlineFunctions_DisallowInliningReason) Ordinal() int { return int(e) }

// InlineFunctions_ShouldRequireInlining is the record
// InlineFunctions.ShouldRequireInlining(boolean shouldValidate, String
// optMessage).
type InlineFunctions_ShouldRequireInlining struct {
	shouldValidate bool
	optMessage     string
}

// Java: InlineFunctions.ShouldRequireInlining(boolean, String)
func NewInlineFunctions_ShouldRequireInlining(shouldValidate bool, optMessage string) *InlineFunctions_ShouldRequireInlining {
	return &InlineFunctions_ShouldRequireInlining{shouldValidate: shouldValidate, optMessage: optMessage}
}

// Java: InlineFunctions.ShouldRequireInlining.shouldValidate()
func (r *InlineFunctions_ShouldRequireInlining) ShouldValidate() bool { return r.shouldValidate }

// Java: InlineFunctions.ShouldRequireInlining.optMessage()
func (r *InlineFunctions_ShouldRequireInlining) OptMessage() string { return r.optMessage }

// Java: InlineFunctions.ShouldRequireInlining.equals(Object)
func (r *InlineFunctions_ShouldRequireInlining) Equals(o any) bool {
	other, ok := o.(*InlineFunctions_ShouldRequireInlining)
	return ok && other != nil && r.shouldValidate == other.shouldValidate && r.optMessage == other.optMessage
}

// Java: InlineFunctions.ShouldRequireInlining.hashCode()
func (r *InlineFunctions_ShouldRequireInlining) HashCode() int32 {
	return 31*jutil.ObjectHashCode(r.shouldValidate) + jutil.StringHashCode(r.optMessage)
}

// Java: InlineFunctions.ShouldRequireInlining.toString()
func (r *InlineFunctions_ShouldRequireInlining) String() string {
	return "ShouldRequireInlining[shouldValidate=" + strconv.FormatBool(r.shouldValidate) + ", optMessage=" + r.optMessage + "]"
}

// Java: InlineFunctions.checkOptimizationLevelSupportsRequireInlining(CompilerOptions)
func InlineFunctions_CheckOptimizationLevelSupportsRequireInlining(
	options *CompilerOptions) *InlineFunctions_ShouldRequireInlining {
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
		return NewInlineFunctions_ShouldRequireInlining(true, "")
	}
	return NewInlineFunctions_ShouldRequireInlining(false, "["+strings.Join(problems, ", ")+"]")
}
