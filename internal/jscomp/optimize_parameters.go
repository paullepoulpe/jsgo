// Port of com/google/javascript/jscomp/OptimizeParameters.java @ closure-compiler v20261015

/*
 * Copyright 2009 The Closure Compiler Authors.
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
	"reflect"

	"github.com/paullepoulpe/jsgo/internal/diagnostic"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// OptimizeParameters optimizes function calls and function signatures.
//
//   - Removes optional parameters if no caller specifies it as argument.
//   - Removes arguments at call site to function that ignores the parameter.
//   - Inline a parameter if the function is always called with that
//     constant.
//   - Removes trailing `undefined` values if the callee doesn't query
//     `arguments` or have a rest parameter
type OptimizeParameters struct {
	compiler    AbstractCompiler
	astAnalyzer *AstAnalyzer
	globalScope *Scope

	// Allocated & cleaned up by process()
	decisionsLog      diagnostic.LogFile
	fnNodes           []*rhino.Node
	defNodes          []*rhino.Node
	ttlList           jutil.ArrayList[*rhino.Node]
	fixedParams       jutil.ArrayList[*OptimizeParameters_Parameter]
	paramSlab         []OptimizeParameters_Parameter
	candidateBuilder  OptimizeParameters_CandidateAnalysisBuilder
	candidateAnalysis OptimizeParameters_CandidateAnalysis
}

// Java: OptimizeParameters(AbstractCompiler)
func NewOptimizeParameters(compiler AbstractCompiler) *OptimizeParameters {
	o := &OptimizeParameters{}
	o.compiler = jutil.CheckNotNull(compiler)
	o.astAnalyzer = compiler.GetAstAnalyzer()
	return o
}

// Process is CompilerPass.process.
//
// Go-only: language-adapter — the two Java overloads process(Node, Node)
// (CompilerPass) and process(Node, Node, ReferenceMap)
// (CallGraphCompilerPass) cannot share a Go name; the second one is
// ProcessWithReferenceMap, and AsCallGraphCompilerPass adapts it.
//
// Java: OptimizeParameters.process(Node, Node)
func (o *OptimizeParameters) Process(externs *rhino.Node, root *rhino.Node) {
	jutil.CheckState(o.compiler.GetLifeCycleStage() == AbstractCompiler_LifeCycleStage_NORMALIZED)

	OptimizeCalls_Builder_().
		SetCompiler(o.compiler).
		SetConsiderExterns(false).
		AddPass(o.AsCallGraphCompilerPass()).
		Build().
		Process(externs, root)
}

// AsCallGraphCompilerPass returns this pass as an
// OptimizeCalls.CallGraphCompilerPass (what Java does by passing `this`).
//
// Go-only: language-adapter — see Process.
func (o *OptimizeParameters) AsCallGraphCompilerPass() OptimizeCalls_CallGraphCompilerPass {
	return OptimizeCalls_CallGraphCompilerPassFunc(o.ProcessWithReferenceMap)
}

// Java: OptimizeParameters.process(Node, Node, ReferenceMap)
func (o *OptimizeParameters) ProcessWithReferenceMap(externs *rhino.Node, root *rhino.Node, refMap *OptimizeCalls_ReferenceMap) {
	decisionsLog :=
		o.compiler.CreateOrReopenIndexedLog(reflect.TypeFor[OptimizeParameters](), "decisions.log")
	defer decisionsLog.Close()
	defer func() {
		o.decisionsLog = nil
		clear(o.fnNodes)
		o.fnNodes = o.fnNodes[:0]
		clear(o.defNodes)
		o.defNodes = o.defNodes[:0]
		o.ttlList.Clear()
		o.fixedParams.Clear()
		clear(o.paramSlab)
		o.paramSlab = o.paramSlab[:0]
	}()
	// Save the LogFile into a field to avoid bucket-brigade passing it through a bunch of methods
	o.decisionsLog = decisionsLog
	o.globalScope = refMap.GetGlobalScope()

	// Find all function nodes that are possible candidates for parameter removal.
	var toOptimize []*jutil.ArrayList[*rhino.Node]

	for key, refs := range refMap.GetNameReferences() {
		candidateAnalysis := o.AnalyzeCandidateName(key, refs)
		if candidateAnalysis.ShouldConvertTaggedTemplateLiterals() {
			candidateAnalysis.ConvertTaggedTemplateLiterals()
		}
		if candidateAnalysis.IsSafeToOptimize() {
			toOptimize = append(toOptimize, refs)
		}
	}

	for key, refs := range refMap.GetPropReferences() {
		candidateAnalysis := o.AnalyzeCandidateProperty(key, refs)
		if candidateAnalysis.ShouldConvertTaggedTemplateLiterals() {
			candidateAnalysis.ConvertTaggedTemplateLiterals()
		}
		if candidateAnalysis.IsSafeToOptimize() {
			toOptimize = append(toOptimize, refs)
		}
	}

	// NOTE: The optimization that are perform must be careful to keep the
	// ReferenceMap in a consistent state. They should be careful to not
	// remove or add global references without update the reference map.
	// While adding references to the map is O(1), removing references
	// O(n) where (n) is the number of references.
	//
	// So the most transformative pass should be last. So:
	//
	// - Removing parameters not provided by any call-site only
	// moves the name and default values from the parameter list to the
	// functions' bodies, so no updates are needed if the same node
	// are reused.
	//
	// - Moving parameters that are provided the same value by every
	// call-site to the function bodies, is currently limited to cases
	// where there are only one function definition, so no updates
	// are needed if the same nodes are reused.
	//
	// - Removing parameters that are unreferenced from call-sites, may
	// remove references, as it is run last.

	for _, refs := range toOptimize {
		o.TryEliminateOptionalArgs(refs)
	}

	for _, refs := range toOptimize {
		o.TryEliminateConstantArgs(refs)
	}

	// tryEliminateUnusedArgs may mutate
	optimizer := NewOptimizeParameters_UnusedParameterOptimizer(o)
	for _, refs := range toOptimize {
		optimizer.TryEliminateUnusedArgs(refs)
	}
	optimizer.ApplyChanges()
}

// optimizeParameters_values is ImmutableListMultimap.values().
//
// Go-only: language-adapter — collects the multimap's values in order.
func optimizeParameters_values(fns *jutil.ListMultimap[*rhino.Node, *rhino.Node]) []*rhino.Node {
	var values []*rhino.Node
	for _, fn := range fns.Entries() {
		values = append(values, fn)
	}
	return values
}

// OptimizeParameters_UnusedParameterOptimizer is the inner class
// UnusedParameterOptimizer.
type OptimizeParameters_UnusedParameterOptimizer struct {
	outer             *OptimizeParameters
	toRemove          []*rhino.Node
	toReplaceWithZero []*rhino.Node
	used              jutil.BitSet
	unused            jutil.BitSet
	unremovable       jutil.BitSet
}

// Java: OptimizeParameters.UnusedParameterOptimizer()
func NewOptimizeParameters_UnusedParameterOptimizer(outer *OptimizeParameters) *OptimizeParameters_UnusedParameterOptimizer {
	return &OptimizeParameters_UnusedParameterOptimizer{outer: outer}
}

// TryEliminateUnusedArgs attempts to eliminate unused parameters by removing
// them from both the call sites and the function definitions.
//
// An unused first parameter: function foo(a, b) {use(b);} foo(1,2);
// foo(1,3) becomes function foo(b) {use(b);} foo(2); foo(3);
//
// refs is a list of references to the symbol (name or property) as vetted
// by #analyzeCandidate.
//
// Java: OptimizeParameters.UnusedParameterOptimizer.tryEliminateUnusedArgs(ArrayList<Node>)
func (u *OptimizeParameters_UnusedParameterOptimizer) TryEliminateUnusedArgs(refs *jutil.ArrayList[*rhino.Node]) {
	// An argument is unused if its position is greater than the number of declared parameters
	// or if it marked as unused.
	fns := OptimizeCalls_ReferenceMap_AppendFunctionNodes(u.outer.fnNodes[:0], refs.Slice())
	u.outer.fnNodes = fns
	jutil.CheckState(len(fns) > 0)

	// Examine all function definitions that are ever assigned to the symbol to determine:
	// 1. Which formal parameter positions are used by at least one of the definitions?
	// 2. What is the largest number of formal parameters across all of the functions?
	// 3. The lowest formal parameter position that contains a rest parameter that is used.
	// e.g.
	// foo = function(used0, unused1, ...usedRest) {}
	// foo = function(unused0, ...unusedRest) {}
	// In this case maxFormalsCount = 3, lowestUsedRest = 2, and used = { 0, 2 }
	maxFormalsCount := 0
	lowestUsedRest := math.MaxInt32
	used := &u.used
	used.ClearAll()
	for _, fn := range fns {
		paramList := NodeUtil_GetFunctionParameters(fn)
		index := -1
		for c := paramList.GetFirstChild(); c != nil; c = c.GetNext() {
			index++
			if !c.IsUnusedParameter() {
				used.Set(index)
				if c.IsRest() {
					lowestUsedRest = min(lowestUsedRest, index)
					if lowestUsedRest == 0 {
						// don't bother doing anything more, all the parameters are used.
						return
					}
				}
			}
		}

		maxFormalsCount = max(maxFormalsCount, index+1)
	}

	// every argument slot after the earliest rest is used
	if lowestUsedRest < maxFormalsCount {
		used.SetRange(lowestUsedRest, maxFormalsCount)
	}

	unused := &u.unused
	unused.CopyFrom(used)
	unused.FlipRange(0, maxFormalsCount)

	// If was a used "rest" declaration, there are no trailing parameters to remove, so
	// bail out now if there are no unused formals.
	if lowestUsedRest < maxFormalsCount && unused.Cardinality() == 0 {
		return
	}

	// NOTE: RemoveUnusedCode removes any trailing unused formals, so we don't need to
	// do for that case.

	unremovable := &u.unremovable
	unremovable.CopyFrom(used)

	// A parameter is removable from a call-site if the parameter value
	// has no side-effects. If all call-sites can be updated, the
	// parameter can be removed from the function definitions.
	// Regardless, the side-effect free values can still be replaced
	// with a simpler expression.

	// 3 values determine the range of removable parameters:
	// the number of formal parameters, the unused parameters bitset,
	// and the position of a used rest parameter (if any).

	// - If the parameter index is higher >= the "rest", it
	// is not a candidate, otherwise if it is in bitset or
	// the greater than the declared formals.

	// To be removed, the candidate must be side-effect free.

	// If not all candidates can be removed, the removable
	// candidates are still removed if there are no following parameters.
	// If there are following parameters the removable candidates are
	// replaced with a literal zero instead.

	// Build a list of parameters to remove
	lowestSpread := math.MaxInt32
	for _, n := range refs.Slice() {
		if OptimizeCalls_ReferenceMap_IsNormalOrOptChainCallOrNewTarget(n) {
			if OptimizeParameters_HasSpreadInFirstCallParam(n) {
				// Bail: with spread in the first argument of .call, we know nothing about any of the
				// parameters, so bail out now.
				return
			}
			param := OptimizeCalls_ReferenceMap_GetFirstArgumentForCallOrNewOrDotCall(n)
			paramIndex := 0
			for param != nil {
				if paramIndex >= maxFormalsCount {
					break
				}

				if param.IsSpread() {
					lowestSpread = min(lowestSpread, paramIndex)
					break
				}

				if !unremovable.Get(paramIndex) && u.outer.astAnalyzer.MayHaveSideEffects(param) {
					unremovable.Set(paramIndex)
				}

				param = param.GetNext()
				paramIndex++
			}
		}
	}

	// Although, a spread prevents the removal of it and all following used slots, it doesn't
	// prevent the replacement of unused values from other call sites
	if lowestSpread < maxFormalsCount {
		unremovable.SetRange(lowestSpread, maxFormalsCount)
	}

	// Only remove trailing parameters if there isn't a rest arguments in any of the
	// definition sites.
	removeAllAfterIndex := math.MaxInt32
	if lowestUsedRest == math.MaxInt32 {
		removeAllAfterIndex = maxFormalsCount - 1
	}

	for _, n := range refs.Slice() {
		if OptimizeCalls_ReferenceMap_IsNormalOrOptChainCallOrNewTarget(n) && !OptimizeParameters_AlreadyRemoved(n) {
			arg := OptimizeCalls_ReferenceMap_GetFirstArgumentForCallOrNewOrDotCall(n)
			u.RecordRemovalCallArguments(
				lowestUsedRest, removeAllAfterIndex, unused, unremovable, arg, 0)
		}
	}

	for _, fn := range fns {
		paramList := NodeUtil_GetFunctionParameters(fn)
		param := paramList.GetFirstChild()
		u.RemoveUnusedFunctionParameters(unremovable, param, 0)
	}
}

// RecordRemovalCallArguments: either firstRestIndex will be MAX_VALUE or
// removeAllAfterIndex will be MAX_VALUE.
//
// firstRestIndex is the lowest index that might match a rest parameter;
// removeAllAfterIndex the index of the last known parameter, all arguments
// with higher indexes must be unused, but it won't be tracked in unused;
// unused all parameter indexes that are known to be unused by the callee;
// unremovable all parameter indexes that are used by a callee, or have
// defaults with side effects; arg the current argument being considered for
// removal; index the index of arg in the argument list. Returns true if all
// arguments >= index have been removed.
//
// Java: OptimizeParameters.UnusedParameterOptimizer.recordRemovalCallArguments(int, int, BitSet, BitSet, Node, int)
func (u *OptimizeParameters_UnusedParameterOptimizer) RecordRemovalCallArguments(
	firstRestIndex int,
	removeAllAfterIndex int,
	unused *jutil.BitSet,
	unremovable *jutil.BitSet,
	arg *rhino.Node,
	index int) bool {
	if arg == nil {
		// base case
		return true
	}

	if index > removeAllAfterIndex {
		return u.RemoveArgAndFollowing(arg)
	}

	if arg.IsSpread() {
		// There is no meaningful "index" after a spread.
		return false
	}

	removedAllTrailing :=
		u.RecordRemovalCallArguments(
			firstRestIndex, removeAllAfterIndex, unused, unremovable, arg.GetNext(), index+1)
	if index < firstRestIndex {
		// An 'undefined' in trailing position calling a function that doesn't reference
		// `arguments` is removable, since the function will see an `undefined` value there
		// even if we remove it.
		isRemovableTrailingUndefined := NodeUtil_IsUndefined(arg) && removedAllTrailing
		if isRemovableTrailingUndefined && !u.outer.astAnalyzer.MayHaveSideEffects(arg) {
			u.toRemove = append(u.toRemove, arg)
			return true
		}

		if unused.Get(index) {
			// If isRemovableTrailingUndefined is true, then we know the arg has side effects
			// otherwise we would have already returned above.
			hasSideEffects :=
				isRemovableTrailingUndefined || u.outer.astAnalyzer.MayHaveSideEffects(arg)
			if !hasSideEffects {
				if unremovable.Get(index) {
					u.RecordReplaceWithZero(arg)
				} else {
					u.toRemove = append(u.toRemove, arg)
					return removedAllTrailing
				}
			} else {
				// If there is a comma operator with side effects on the first node, the
				// second node can be optimized knowing that the parameter is not used.
				u.RemoveFromCommaIfNoSideEffects(arg)
			}
		}
	}
	return false
}

// RemoveFromCommaIfNoSideEffects replaces a no side-effect value in a comma
// operator recursively.
//
// Java: OptimizeParameters.UnusedParameterOptimizer.removeFromCommaIfNoSideEffects(Node)
func (u *OptimizeParameters_UnusedParameterOptimizer) RemoveFromCommaIfNoSideEffects(arg *rhino.Node) {
	if arg.IsComma() {
		secondChild := arg.GetSecondChild()
		if !u.outer.astAnalyzer.MayHaveSideEffects(secondChild) {
			u.RecordReplaceWithZero(secondChild)
		} else {
			u.RemoveFromCommaIfNoSideEffects(secondChild)
		}
	}
}

// RecordReplaceWithZero records node to be replaced with zero.
//
// Java: OptimizeParameters.UnusedParameterOptimizer.recordReplaceWithZero(Node)
func (u *OptimizeParameters_UnusedParameterOptimizer) RecordReplaceWithZero(arg *rhino.Node) {
	if !arg.IsNumber() || arg.GetDouble() != 0 {
		u.toReplaceWithZero = append(u.toReplaceWithZero, arg)
	}
}

// RemoveArgAndFollowing returns true if all following args were removed.
//
// Java: OptimizeParameters.UnusedParameterOptimizer.removeArgAndFollowing(Node)
func (u *OptimizeParameters_UnusedParameterOptimizer) RemoveArgAndFollowing(arg *rhino.Node) bool {
	if arg != nil {
		removedAll := u.RemoveArgAndFollowing(arg.GetNext())
		if !u.outer.astAnalyzer.MayHaveSideEffects(arg) {
			u.toRemove = append(u.toRemove, arg)
			return removedAll
		}
		return false
	}
	return true
}

// Java: OptimizeParameters.UnusedParameterOptimizer.removeUnusedFunctionParameters(BitSet, Node, int)
func (u *OptimizeParameters_UnusedParameterOptimizer) RemoveUnusedFunctionParameters(unremovable *jutil.BitSet, param *rhino.Node, index int) {
	if param != nil {
		u.RemoveUnusedFunctionParameters(unremovable, param.GetNext(), index+1)
		if !unremovable.Get(index) {
			jutil.CheckState(param.IsName()) // update for ES6
			// params are not otherwise referenceable.
			u.outer.compiler.ReportChangeToEnclosingScope(param)
			param.Detach()
		}
	}
}

// ApplyChanges applies optimizations to all previously marked nodes.
//
// Java: OptimizeParameters.UnusedParameterOptimizer.applyChanges()
func (u *OptimizeParameters_UnusedParameterOptimizer) ApplyChanges() {
	for _, n := range u.toRemove {
		// Don't remove any nodes twice since doing so would violate change reporting constraints.
		if OptimizeParameters_AlreadyRemoved(n) {
			continue
		}
		u.outer.compiler.ReportChangeToEnclosingScope(n)
		n.Detach()
		NodeUtil_MarkFunctionsDeleted(n, u.outer.compiler)
	}
	for _, n := range u.toReplaceWithZero {
		jutil.CheckState(!n.IsNumber() || n.GetDouble() != 0.0)
		// Don't remove any nodes twice since doing so would violate change reporting constraints.
		if OptimizeParameters_AlreadyRemoved(n) {
			continue
		}
		u.outer.compiler.ReportChangeToEnclosingScope(n)
		n.ReplaceWith(rhino.IR_Number(0).Srcref(n))
		NodeUtil_MarkFunctionsDeleted(n, u.outer.compiler)
	}
}

// Java: OptimizeParameters.alreadyRemoved(Node)
func OptimizeParameters_AlreadyRemoved(n *rhino.Node) bool {
	parent := n.GetParent()
	for parent != nil {
		n = parent
		parent = n.GetParent()
	}
	return !n.IsRoot()
}

// Java: OptimizeParameters.analyzeCandidateName(String, ArrayList<Node>)
func (o *OptimizeParameters) AnalyzeCandidateName(key string, refs *jutil.ArrayList[*rhino.Node]) *OptimizeParameters_CandidateAnalysis {
	return o.AnalyzeCandidate("name", key, refs)
}

// Java: OptimizeParameters.analyzeCandidateProperty(String, ArrayList<Node>)
func (o *OptimizeParameters) AnalyzeCandidateProperty(key string, refs *jutil.ArrayList[*rhino.Node]) *OptimizeParameters_CandidateAnalysis {
	return o.AnalyzeCandidate("property", key, refs)
}

// AnalyzeCandidate: a function is a candidate for parameter optimization if:
//
//   - if all call sites are known (no aliasing)
//   - if all definition sites are known (the possible values are known
//     functions)
//   - there is at least one definition
//   - none of the calls are tagged template literals, OR the special first
//     argument is unused
//
// Java: OptimizeParameters.analyzeCandidate(String, String, ArrayList<Node>)
func (o *OptimizeParameters) AnalyzeCandidate(refKind string, key string, refs *jutil.ArrayList[*rhino.Node]) *OptimizeParameters_CandidateAnalysis {
	o.ttlList.Clear()
	analysisBuilder := &o.candidateBuilder
	analysisBuilder.outer = o
	analysisBuilder.taggedTemplateLiterals = &o.ttlList
	analysisBuilder.isSafeToOptimize = false
	if !OptimizeCalls_MayBeOptimizableName(o.compiler, key) {
		if o.decisionsLog.IsLogging() {
			o.decisionsLog.LogFormat("%s\t%s\tnot an optimizable name", refKind, key)
		}
		return analysisBuilder.SetIsSafeToOptimize(false).Build()
	}
	definitions := o.defNodes[:0]
	seenCandidateUse := false
	for _, n := range refs.Slice() {
		// TODO(johnlenz): Determine what to do about ".constructor" references.
		// Currently classes that are super classes or have superclasses aren't optimized
		//
		// if (parent.isCall() && n != parent.getFirstChild() && isClassDefiningCall(parent)) {
		//   continue;
		// } else

		if OptimizeCalls_ReferenceMap_IsNormalOrOptChainCallOrNewTarget(n) {
			// TODO(johnlenz): filter .apply when we support it
			seenCandidateUse = true
		} else if n.GetParent().IsTaggedTemplateLit() {
			analysisBuilder.AddTaggedTemplateLiteral(n.GetParent())
		} else if o.IsCandidateDefinition(n) {
			definitions = append(definitions, n)
		} else if !OptimizeCalls_IsAllowedReference(n) {
			if o.decisionsLog.IsLogging() { // avoid build location string when not logging
				o.decisionsLog.LogFormat("%s\t%s\tnot an allowed reference: %s", refKind, key, n.GetLocation())
			}
			// TODO(johnlenz): allow extends clauses.
			o.defNodes = definitions
			return analysisBuilder.SetIsSafeToOptimize(false).Build()
		}
	}
	o.defNodes = definitions
	if len(definitions) == 0 {
		if o.decisionsLog.IsLogging() {
			o.decisionsLog.LogFormat("%s\t%s\tno definition found", refKind, key)
		}
		return analysisBuilder.SetIsSafeToOptimize(false).Build()
	}
	if analysisBuilder.taggedTemplateLiterals.Size() != 0 {
		functionNodes := OptimizeCalls_ReferenceMap_AppendFunctionNodes(o.fnNodes[:0], definitions)
		o.fnNodes = functionNodes

		for _, functionNode := range functionNodes {
			firstParam := NodeUtil_GetFunctionParameters(functionNode).GetFirstChild()
			if firstParam != nil && !firstParam.IsUnusedParameter() {
				if o.decisionsLog.IsLogging() {
					o.decisionsLog.LogFormat(
						"%s\t%s\twill not optimize parameters for tagged template literals", refKind, key)
				}
				return analysisBuilder.SetIsSafeToOptimize(false).Build()
			}
		}
		if o.decisionsLog.IsLogging() {
			o.decisionsLog.LogFormat(
				"%s\t%s\t%s",
				refKind,
				key,
				"first param is unused, so we will convert tagged template literals to normal calls")
		}
		seenCandidateUse = true
	}

	if !seenCandidateUse {
		if o.decisionsLog.IsLogging() {
			o.decisionsLog.LogFormat("%s\t%s\tno usage found", refKind, key)
		}
		return analysisBuilder.SetIsSafeToOptimize(false).Build()
	}
	return analysisBuilder.SetIsSafeToOptimize(true).Build()
}

// OptimizeParameters_CandidateAnalysis holds the results of analyzing a
// candidate function's definition and references.
type OptimizeParameters_CandidateAnalysis struct {
	outer                  *OptimizeParameters
	taggedTemplateLiterals *jutil.ArrayList[*rhino.Node]
	isSafeToOptimize       bool
}

// Java: OptimizeParameters.CandidateAnalysis(CandidateAnalysisBuilder)
func NewOptimizeParameters_CandidateAnalysis(outer *OptimizeParameters, builder *OptimizeParameters_CandidateAnalysisBuilder) *OptimizeParameters_CandidateAnalysis {
	return &OptimizeParameters_CandidateAnalysis{
		outer:                  outer,
		taggedTemplateLiterals: builder.taggedTemplateLiterals,
		isSafeToOptimize:       builder.isSafeToOptimize,
	}
}

// Java: OptimizeParameters.CandidateAnalysis.shouldConvertTaggedTemplateLiterals()
func (c *OptimizeParameters_CandidateAnalysis) ShouldConvertTaggedTemplateLiterals() bool {
	return !c.taggedTemplateLiterals.IsEmpty() && c.isSafeToOptimize
}

// Java: OptimizeParameters.CandidateAnalysis.convertTaggedTemplateLiterals()
func (c *OptimizeParameters_CandidateAnalysis) ConvertTaggedTemplateLiterals() {
	for _, taggedTemplateLiteral := range c.taggedTemplateLiterals.Slice() {
		c.outer.ConvertTaggedTemplateLiteralToNormalCall(taggedTemplateLiteral)
	}
	c.taggedTemplateLiterals.Clear()
}

// Java: OptimizeParameters.CandidateAnalysis.isSafeToOptimize()
func (c *OptimizeParameters_CandidateAnalysis) IsSafeToOptimize() bool {
	return c.isSafeToOptimize && c.taggedTemplateLiterals.IsEmpty()
}

// OptimizeParameters_CandidateAnalysisBuilder is the builder class for
// CandidateAnalysis.
type OptimizeParameters_CandidateAnalysisBuilder struct {
	outer                  *OptimizeParameters
	taggedTemplateLiterals *jutil.ArrayList[*rhino.Node]
	isSafeToOptimize       bool
}

// Java: OptimizeParameters.CandidateAnalysisBuilder()
func NewOptimizeParameters_CandidateAnalysisBuilder(outer *OptimizeParameters) *OptimizeParameters_CandidateAnalysisBuilder {
	return &OptimizeParameters_CandidateAnalysisBuilder{
		outer:                  outer,
		taggedTemplateLiterals: jutil.NewArrayList[*rhino.Node](),
		isSafeToOptimize:       false,
	}
}

// Java: OptimizeParameters.CandidateAnalysisBuilder.setIsSafeToOptimize(boolean)
func (b *OptimizeParameters_CandidateAnalysisBuilder) SetIsSafeToOptimize(isSafeToOptimize bool) *OptimizeParameters_CandidateAnalysisBuilder {
	b.isSafeToOptimize = isSafeToOptimize
	return b
}

// Java: OptimizeParameters.CandidateAnalysisBuilder.addTaggedTemplateLiteral(Node)
func (b *OptimizeParameters_CandidateAnalysisBuilder) AddTaggedTemplateLiteral(ttlNode *rhino.Node) *OptimizeParameters_CandidateAnalysisBuilder {
	b.taggedTemplateLiterals.Add(ttlNode)
	return b
}

// Java: OptimizeParameters.CandidateAnalysisBuilder.build()
func (b *OptimizeParameters_CandidateAnalysisBuilder) Build() *OptimizeParameters_CandidateAnalysis {
	if b.outer != nil && b == &b.outer.candidateBuilder {
		b.outer.candidateAnalysis = OptimizeParameters_CandidateAnalysis{
			outer:                  b.outer,
			taggedTemplateLiterals: b.taggedTemplateLiterals,
			isSafeToOptimize:       b.isSafeToOptimize,
		}
		return &b.outer.candidateAnalysis
	}
	return NewOptimizeParameters_CandidateAnalysis(b.outer, b)
}

// ConvertTaggedTemplateLiteralToNormalCall converts "tagFunction`template
// ${1} literal ${2}`" to `tagFunction([], 1, 2)`.
//
// Java: OptimizeParameters.convertTaggedTemplateLiteralToNormalCall(Node)
func (o *OptimizeParameters) ConvertTaggedTemplateLiteralToNormalCall(ttlNode *rhino.Node) {
	callee := jutil.CheckNotNull(ttlNode.GetFirstChild())
	templateLiteral := callee.GetNext()
	firstArg := rhino.IR_Arraylit().Srcref(templateLiteral)
	callee.Detach()
	callNode := NodeUtil_NewCallNode(callee, firstArg).Srcref(ttlNode)
	for tlChild := templateLiteral.GetFirstChild(); tlChild != nil; tlChild = tlChild.GetNext() {
		if tlChild.IsTemplateLitSub() {
			arg := tlChild.GetOnlyChild()
			arg.Detach()
			callNode.AddChildToBack(arg)
		}
	}
	ttlNode.ReplaceWith(callNode)
	o.compiler.ReportChangeToEnclosingScope(callNode)
}

// Java: OptimizeParameters.isCandidateDefinition(Node)
func (o *OptimizeParameters) IsCandidateDefinition(n *rhino.Node) bool {
	parent := n.GetParent()

	var functionExpr *rhino.Node
	if parent.IsFunction() && NodeUtil_IsFunctionDeclaration(parent) {
		functionExpr = parent
	} else if OptimizeCalls_ReferenceMap_IsSimpleAssignmentTarget(n) {
		functionExpr = parent.GetLastChild()
	} else if n.IsName() && n.HasChildren() {
		functionExpr = n.GetFirstChild()
	} else if o.IsClassMemberDefinition(n) {
		functionExpr = n.GetFirstChild()
	} else if parent.IsClass() && n.IsFirstChildOf(parent) {
		// allDefinitionsAreCandidateFunctions() understands classes and will check for
		// candidacy correctly.
		functionExpr = parent
	} else {
		return false // Couldn't find a function.
	}

	return o.AllDefinitionsAreCandidateFunctions(functionExpr)
}

// Java: OptimizeParameters.isClassMemberDefinition(Node)
func (o *OptimizeParameters) IsClassMemberDefinition(n *rhino.Node) bool {
	return n.IsMemberFunctionDef() && n.GetParent().IsClassMembers()
}

// ClassContainsClassFieldWithRHSSideEffects: when class fields are set to
// RHS with side effects, the order of execution may be moved into incorrect
// order if optimized in this pass. Example:
//
// class C { field2 = alert(2);", constructor(a) {", use(a); } } var c = new
// C(alert(1));
//
// would optimize to
//
// class C { field2 = alert(2);", constructor(a) {", var a = alert(1)
// use(a); } } var c = new C(alert(1));
//
// which would mean alert(2) gets executed before alert(1) - this behavior is
// incorrect.
//
// So, if there are any RHS side effects, we skip the optimization.
//
// Java: OptimizeParameters.classContainsClassFieldWithRHSSideEffects(Node)
func (o *OptimizeParameters) ClassContainsClassFieldWithRHSSideEffects(classNode *rhino.Node) bool {
	classMembersNode := NodeUtil_GetClassMembers(classNode)
	for member := classMembersNode.GetFirstChild(); member != nil; member = member.GetNext() {
		initializer := member.GetLastChild()
		if initializer != nil &&
			NodeUtil_IsNonStaticClassFieldInitializer(initializer) &&
			o.astAnalyzer.MayHaveSideEffects(initializer) {
			return true
		}
	}
	return false
}

// Java: OptimizeParameters.allDefinitionsAreCandidateFunctions(Node)
func (o *OptimizeParameters) AllDefinitionsAreCandidateFunctions(n *rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_CLASS:
		if NodeUtil_IsNamedClassExpression(n) {
			// name creates an alias, making it hard to be sure we've seen all calls
			return false
		} else {
			// `class NameNode {`
			// find the constructor
			constructorMemberFunctionDef := NodeUtil_GetEs6ClassConstructorMemberFunctionDef(n)
			if constructorMemberFunctionDef == nil {
				// unable to find the constructor
				// TODO(bradfordcsmith): Ideally we should find the parent class constructor.
				return false
			} else if o.ClassContainsClassFieldWithRHSSideEffects(n) {
				return false
			} else {
				functionNode := constructorMemberFunctionDef.GetOnlyChild()
				// "arguments" can refer to all parameters or their count.
				return !NodeUtil_DoesFunctionReferenceOwnArgumentsObject(functionNode) &&
					// In `function f(a, b = a) { ... }` it's very difficult to determine if `a` is
					// movable.
					!OptimizeParameters_MayReferenceParamBeforeBody(functionNode) &&
					// `/** @usedViaDotConstructor */ function f(a) {...}` means back off.
					!OptimizeParameters_IsUsedViaDotConstructor(functionNode)
			}
		}
	case rhino.Token_FUNCTION:
		// Named function expression can refer to themselves,
		return !NodeUtil_IsNamedFunctionExpression(n) &&
			// "arguments" can refer to all parameters or their count.
			!NodeUtil_DoesFunctionReferenceOwnArgumentsObject(n) &&
			// In `function f(a, b = a) { ... }` it's very difficult to determine if `a` is movable.
			!OptimizeParameters_MayReferenceParamBeforeBody(n) &&
			// `/** @usedViaDotConstructor */ function f(a) {...}` means back off.
			!OptimizeParameters_IsUsedViaDotConstructor(n)
	case rhino.Token_CAST, rhino.Token_COMMA:
		return o.AllDefinitionsAreCandidateFunctions(n.GetLastChild())
	case rhino.Token_HOOK:
		return o.AllDefinitionsAreCandidateFunctions(n.GetSecondChild()) &&
			o.AllDefinitionsAreCandidateFunctions(n.GetLastChild())
	case rhino.Token_OR, rhino.Token_AND, rhino.Token_COALESCE:
		return o.AllDefinitionsAreCandidateFunctions(n.GetFirstChild()) &&
			o.AllDefinitionsAreCandidateFunctions(n.GetLastChild())
	default:
		return false
	}
}

// OptimizeParameters_IsUsedViaDotConstructor returns true if the function is
// annotated with @usedViaDotConstructor
//
// In this case we just back off on all parameter optimizations since we
// cannot detect which parameters are used via calls through .constructor.
//
// Java: OptimizeParameters.isUsedViaDotConstructor(Node)
func OptimizeParameters_IsUsedViaDotConstructor(function *rhino.Node) bool {
	docInfo := NodeUtil_GetBestJSDocInfo(function)
	return docInfo != nil && docInfo.IsUsedViaDotConstructor()
}

// OptimizeParameters_MayReferenceParamBeforeBody: does the function use one
// of its parameters in code before the body?
//
// Having that property is risky for inlining. Example `function f(a, b = a)
// { ... }`. We can't trivially inline `a` in this case because the inlined
// var can't precede `b = a`.
//
// This case is very rare so for now we just back-off completely. If it
// becomes more common, we can tighten the detection of problematic cases, or
// back-off only for the dangerous params.
//
// Java: OptimizeParameters.mayReferenceParamBeforeBody(Node)
func OptimizeParameters_MayReferenceParamBeforeBody(function *rhino.Node) bool {
	paramList := function.GetSecondChild()
	if !paramList.HasChildren() {
		return false // Fast path; there can't possibly be back-refs.
	}

	allSimpleNames := true
	count := 0
	for c := paramList.GetFirstChild(); c != nil; c = c.GetNext() {
		if !c.IsName() || c.HasChildren() {
			allSimpleNames = false
			break
		}
		count++
	}
	if allSimpleNames {
		if count <= 1 {
			return false
		}
		if count <= 16 {
			for a := paramList.GetFirstChild(); a != nil; a = a.GetNext() {
				aStr := a.GetString()
				for b := a.GetNext(); b != nil; b = b.GetNext() {
					if aStr == b.GetString() {
						return true
					}
				}
			}
			return false
		}
	}

	// Go-only: language-adapter — an ArrayListMultimap (HashMap of String
	// keys) whose iteration order cannot reach the result (any group
	// qualifies), so insertion order is used.
	namesByNames := jutil.ArrayListMultimap_Create[string, *rhino.Node]()
	NodeUtil_VisitPostOrder(
		paramList,
		NodeUtil_VisitorFunc(func(n *rhino.Node) {
			if n.IsName() {
				namesByNames.Put(n.GetString(), n)
			}
		}))

	for _, names := range namesByNames.AsMap().All() {
		if len(names) == 1 {
			continue // There can't be back-refs if there's only one ref.
		}

		for _, name := range names {
			if NodeUtil_IsLValue(name) {
				return true // One ref is a definition, so the rest of might be back-refs.
			}
		}
	}

	return false
}

// TryEliminateOptionalArgs removes any optional parameters if no callers
// specifies it as an argument.
//
// Java: OptimizeParameters.tryEliminateOptionalArgs(ArrayList<Node>)
func (o *OptimizeParameters) TryEliminateOptionalArgs(refs *jutil.ArrayList[*rhino.Node]) {
	// Count the maximum number of arguments passed into this function all
	// all points of the program.
	maxArgs := -1

	for _, n := range refs.Slice() {
		if OptimizeCalls_ReferenceMap_IsNormalOrOptChainCallOrNewTarget(n) {
			if OptimizeParameters_HasSpreadInFirstCallParam(n) {
				// Bail: with spread we must assume all parameters are used, don't waste
				// any more time.
				return
			}
			numArgs := 0
			firstArg := OptimizeCalls_ReferenceMap_GetFirstArgumentForCallOrNewOrDotCall(n)
			for c := firstArg; c != nil; c = c.GetNext() {
				numArgs++
				if c.IsSpread() {
					// Bail: with spread we must assume all parameters are used, don't waste
					// any more time.
					return
				}
			}

			if numArgs > maxArgs {
				maxArgs = numArgs
			}
		}
	}

	fns := OptimizeCalls_ReferenceMap_AppendFunctionNodes(o.fnNodes[:0], refs.Slice())
	o.fnNodes = fns
	for _, fn := range fns {
		o.EliminateParamsAfter(fn, maxArgs)
	}
}

// TryEliminateConstantArgs eliminates parameters if they are always
// constant.
//
// function foo(a, b) {...} foo(1,2); foo(1,3) becomes function foo(b) { var
// a = 1 ... } foo(2); foo(3);
//
// refs is a list of references to the symbol (name or property) as vetted
// by #analyzeCandidate.
//
// Java: OptimizeParameters.tryEliminateConstantArgs(ArrayList<Node>)
func (o *OptimizeParameters) TryEliminateConstantArgs(refs *jutil.ArrayList[*rhino.Node]) {
	parameters := o.FindFixedArguments(refs)
	if parameters == nil {
		return
	}

	fn, fnCount := optimizeCalls_getSingleFunctionNode(refs.Slice())
	if fnCount > 1 {
		// TODO(johnlenz): support moving simple constants.
		// This requires cloning the tree and avoiding adding additional calls/definitions that will
		// invalidate the reference map
		return
	}

	// Only one definition is currently supported.
	if fnCount == 0 {
		// Iterables.getOnlyElement of an empty iterable (size() <= 1 here).
		panic(jutil.NewNoSuchElementException())
	}

	continueLooking := OptimizeParameters_AdjustForConstraints(fn, parameters.Slice())
	if !continueLooking {
		return
	}

	// Found something to do, move the values from the call sites to the function definitions.
	for _, n := range refs.Slice() {
		if !OptimizeParameters_AlreadyRemoved(n) && OptimizeCalls_ReferenceMap_IsNormalOrOptChainCallOrNewTarget(n) {
			o.OptimizeCallSite(parameters.Slice(), n)
		}
	}

	o.OptimizeFunctionDefinition(parameters.Slice(), fn)
}

// FindFixedArguments: refs is a list of references to the symbol (name or
// property) as vetted by #analyzeCandidate. Returns a list of Parameter
// objects, in the declaration order, which represent potentially movable
// values fixed values from all call sites or null if there are no candidate
// values.
//
func (o *OptimizeParameters) allocParam(arg *rhino.Node, shouldRemove bool) *OptimizeParameters_Parameter {
	if len(o.paramSlab) == cap(o.paramSlab) {
		nextCap := max(cap(o.paramSlab)*2, 16)
		o.paramSlab = make([]OptimizeParameters_Parameter, 0, nextCap)
	}
	idx := len(o.paramSlab)
	o.paramSlab = o.paramSlab[:idx+1]
	p := &o.paramSlab[idx]
	*p = OptimizeParameters_Parameter{shouldRemove: shouldRemove, arg: arg}
	return p
}

// Java: OptimizeParameters.findFixedArguments(ArrayList<Node>)
func (o *OptimizeParameters) FindFixedArguments(refs *jutil.ArrayList[*rhino.Node]) *jutil.ArrayList[*OptimizeParameters_Parameter] {
	parameters := &o.fixedParams
	parameters.Clear()
	clear(o.paramSlab)
	o.paramSlab = o.paramSlab[:0]
	firstCall := true

	// Build a list of parameters to remove
	continueLooking := false
	for _, n := range refs.Slice() {
		if OptimizeCalls_ReferenceMap_IsNormalOrOptChainCallOrNewTarget(n) {
			// Normally, we ignore the first parameter to a .call expression (the 'this' value)
			// but if it is a spread, we know nothing about any of the parameters, so bail out now.
			if OptimizeParameters_HasSpreadInFirstCallParam(n) {
				continueLooking = false
				break
			}
			cur := OptimizeCalls_ReferenceMap_GetFirstArgumentForCallOrNewOrDotCall(n)
			if firstCall {
				// Use the first call to construct a list of parameter values of the
				// function.
				continueLooking = o.BuildInitialParameterList(parameters, cur)
				firstCall = false
			} else {
				// All the rest must match
				continueLooking = o.FindFixedParameters(parameters, cur)
			}
			if !continueLooking {
				break
			}
		}
	}

	if continueLooking {
		return parameters
	}
	return nil
}

// OptimizeParameters_HasSpreadInFirstCallParam returns whether the target
// call is a .call expression whose first argument (the this argument) is a
// spread element (e.g., foo.call(...arr)).
//
// ReferenceMap#getFirstArgumentForCallOrNewOrDotCall skips the first
// argument to .call (treating it as the this value) and callers inspect all
// subsequent arguments, which correctly handles spread arguments in normal
// calls like foo(p1, ...args) and .call invocations with an explicit this
// argument like foo.call(thisArg, p1, ...args). However, if the first
// argument to .call itself is a spread (e.g., foo.call(...arr)), it spreads
// across both the this value and subsequent formal parameter positions, so
// we cannot safely optimize any parameters.
//
// Java: OptimizeParameters.hasSpreadInFirstCallParam(Node)
func OptimizeParameters_HasSpreadInFirstCallParam(n *rhino.Node) bool {
	call := OptimizeCalls_ReferenceMap_GetCallOrNewNodeForTarget(n)
	if NodeUtil_IsFunctionObjectCall(call) {
		firstDotCallParam := call.GetSecondChild()
		return firstDotCallParam != nil && firstDotCallParam.IsSpread()
	}
	return false
}

// OptimizeParameters_AdjustForConstraints adjusts the provided Parameter
// objects value created by #findFixedArguments for "rest" value and
// side-effects which might prevent the motion of the parameters from the
// call sites to the function body.
//
// parameters is a list of Parameter objects summarizing all the call-sites
// and whether any of the parameters are fixed. Returns whether there are any
// movable parameters.
//
// Java: OptimizeParameters.adjustForConstraints(Node, List<Parameter>)
func OptimizeParameters_AdjustForConstraints(fn *rhino.Node, parameters []*OptimizeParameters_Parameter) bool {
	info := NodeUtil_GetBestJSDocInfo(fn)
	if info != nil && info.IsNoInline() {
		return false
	}

	paramList := NodeUtil_GetFunctionParameters(fn)
	lastFormal := paramList.GetLastChild()
	restIndex := math.MaxInt32
	lastNonRestFormal := paramList.GetChildCount() - 1
	formal := lastFormal
	if lastFormal != nil && lastFormal.IsRest() {
		restIndex = lastNonRestFormal
		lastNonRestFormal--
		formal = formal.GetPrevious()
	}

	// A parameter with side-effects can move if there are no following parameters
	// that can be affected.

	// A parameter can be moved if it can't be side-effected (a literal),
	// or there are no following side-effects, that aren't moved.

	anyMovable := false
	seenUnmovableSideEffects := false
	seenUnmoveableSideEffected := false
	allRestValueRemovable := true
	for i := len(parameters) - 1; i >= 0; i-- {
		current := parameters[i]

		// back-off for default values whose default value maybe needed.
		// TODO(johnlenz): handle used default
		if i <= lastNonRestFormal {
			if formal.IsDefaultValue() && current.mayBeUndefined {
				current.shouldRemove = false
			}
			formal = formal.GetPrevious()
		}

		// Preserve side-effect ordering, don't move this parameter if:
		// * the current parameter has side-effects and a following
		// parameters that will not be move can be effected.
		// * the current parameter can be effected and a following
		// parameter that will not be moved has side-effects

		if current.shouldRemove &&
			((seenUnmovableSideEffects && current.CanBeSideEffected()) ||
				(seenUnmoveableSideEffected && current.HasSideEffects())) {
			current.shouldRemove = false
		}

		// If any values that are part of the rest cannot be moved to the function body,
		// then all the rest values must remain at the callsite.
		if i >= restIndex {
			if allRestValueRemovable {
				if !current.shouldRemove {
					anyMovable = false
					allRestValueRemovable = false
					// revisit the trailing params and remark them now that we know they are unremovable.
					for j := i + 1; j < len(parameters); j++ {
						p := parameters[j]
						p.shouldRemove = false
						if p.canBeSideEffected {
							seenUnmoveableSideEffected = true
						}
						if p.hasSideEffects {
							seenUnmovableSideEffects = true
						}
					}
				}
			} else {
				current.shouldRemove = false
			}
		}

		if current.shouldRemove {
			anyMovable = true
		} else {
			if current.canBeSideEffected {
				seenUnmoveableSideEffected = true
			}

			if current.hasSideEffects {
				seenUnmovableSideEffects = true
			}
		}
	}
	return anyMovable
}

// FindFixedParameters determines which parameters use the same expression.
//
// Returns whether any parameter was found that can be updated.
//
// Java: OptimizeParameters.findFixedParameters(List<Parameter>, Node)
func (o *OptimizeParameters) FindFixedParameters(parameters *jutil.ArrayList[*OptimizeParameters_Parameter], cur *rhino.Node) bool {
	anyMovable := false
	index := 0
	for cur != nil {
		var p *OptimizeParameters_Parameter
		if index >= parameters.Size() {
			p = o.allocParam(cur, false)
			parameters.Add(p)
			o.SetParameterSideEffectInfo(p, cur)
		} else {
			p = parameters.Get(index)
			if p.ShouldRemove() {
				value := p.GetArg()
				if !cur.IsEquivalentTo(value) {
					p.SetShouldRemove(false)
				} else {
					anyMovable = true
				}
			}
			p.SetHasSideEffects(p.HasSideEffects() || o.astAnalyzer.MayHaveSideEffects(cur))
			p.SetCanBeSideEffected(p.CanBeSideEffected() || NodeUtil_CanBeSideEffected(cur))
		}
		// Back off optimizing arguments following spread
		if cur.IsSpread() {
			break
		}
		cur = cur.GetNext()
		index++
	}

	for ; index < parameters.Size(); index++ {
		parameters.Get(index).SetShouldRemove(false)
	}

	return anyMovable
}

// BuildInitialParameterList returns whether any parameter was movable.
//
// Java: OptimizeParameters.buildInitialParameterList(List<Parameter>, Node)
func (o *OptimizeParameters) BuildInitialParameterList(parameters *jutil.ArrayList[*OptimizeParameters_Parameter], cur *rhino.Node) bool {
	anyMovable := false
	for cur != nil {
		movable := OptimizeParameters_IsMovableValue(cur, o.globalScope)
		p := o.allocParam(cur, movable)
		o.SetParameterSideEffectInfo(p, cur)
		parameters.Add(p)
		if movable {
			anyMovable = true
		}
		// Back off optimizing arguments following spread
		if cur.IsSpread() {
			break
		}
		cur = cur.GetNext()
	}
	return anyMovable
}

// Java: OptimizeParameters.setParameterSideEffectInfo(Parameter, Node)
func (o *OptimizeParameters) SetParameterSideEffectInfo(p *OptimizeParameters_Parameter, value *rhino.Node) {
	p.SetHasSideEffects(o.astAnalyzer.MayHaveSideEffects(value))
	p.SetCanBeSideEffected(NodeUtil_CanBeSideEffected(value))
	if !value.IsSpread() {
		p.SetMayBeUndefined(NodeUtil_MayBeUndefined(value))
	}
}

// OptimizeParameters_IsMovableValue returns whether the expression can be
// safely moved to another function in another scope.
//
// Java: OptimizeParameters.isMovableValue(Node, Scope)
func OptimizeParameters_IsMovableValue(n *rhino.Node, globalScope *Scope) bool {
	// Things that can change value or are inaccessible can't be moved, these
	// are "this", "arguments", local names, and functions that capture local
	// values.
	switch n.GetToken() {
	case rhino.Token_AWAIT, rhino.Token_YIELD, rhino.Token_THIS, rhino.Token_SUPER, rhino.Token_ITER_SPREAD, rhino.Token_OBJECT_SPREAD, rhino.Token_NEW_TARGET, rhino.Token_IMPORT_META:
		return false
	case rhino.Token_FUNCTION:
		// Don't move function closures.
		// TODO(johnlenz): Closure that only contain global reference can be
		// moved.
		return false
	case rhino.Token_NAME:
		if n.GetString() == "arguments" {
			return false
		} else {
			// If it isn't in global scope, then it is in local scope.  This logic depends on
			// "Normalize" creating distinct local names.
			v := globalScope.GetVar(n.GetString())
			if v == nil {
				return false
			}
		}
	default:
	}

	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		if !OptimizeParameters_IsMovableValue(c, globalScope) {
			return false
		}
	}
	return true
}

// Java: OptimizeParameters.optimizeFunctionDefinition(List<Parameter>, Node)
func (o *OptimizeParameters) OptimizeFunctionDefinition(parameters []*OptimizeParameters_Parameter, fn *rhino.Node) {
	paramList := NodeUtil_GetFunctionParameters(fn)
	maybeRest := paramList.GetLastChild()
	lastParameter := len(parameters) - 1
	if maybeRest != nil && maybeRest.IsRest() {
		restIndex := paramList.GetChildCount() - 1
		// If the rest parameter is removable they all are.
		if len(parameters) < restIndex ||
			(restIndex < len(parameters) && parameters[restIndex].ShouldRemove()) {
			value := rhino.IR_Arraylit().Srcref(maybeRest)
			for i := restIndex; i < len(parameters); i++ {
				parameter := parameters[i]

				jutil.CheckState(parameter.ShouldRemove())
				value.AddChildToBack(parameters[i].GetArg())
			}
			maybeRest.Detach()
			lhs := maybeRest.RemoveFirstChild()
			o.AddRestVariableToFunction(fn, lhs, value)
		}

		// process the rest.
		lastParameter = min(len(parameters)-1, restIndex-1)
	}

	for i := lastParameter; i >= 0; i-- {
		parameter := parameters[i]
		if parameter.ShouldRemove() {
			formalParam := NodeUtil_GetArgumentForFunction(fn, i)
			if formalParam != nil {
				formalParam.Detach()
				if formalParam.IsDefaultValue() {
					// Drop the default value as we should only get here if the default value isn't going
					// to be used.
					jutil.CheckState(!parameter.mayBeUndefined)
					formalParam = formalParam.RemoveFirstChild()
				}
			}

			value := parameters[i].GetArg()

			if formalParam != nil {
				o.AddVariableDeclarationToFunction(fn, formalParam, value)
			} else {
				// no formal param, just add the value to the function body. This is the case when the
				// parameter is unused in the function body; we can skip declaring it and simply adding
				// the argument value to the function body.
				o.AddExpressionToFunction(fn, value)
			}
		}
	}
}

// Java: OptimizeParameters.optimizeCallSite(List<Parameter>, Node)
func (o *OptimizeParameters) OptimizeCallSite(parameters []*OptimizeParameters_Parameter, target *rhino.Node) {
	call := OptimizeCalls_ReferenceMap_GetCallOrNewNodeForTarget(target)
	callMayMutateArgs := call.MayMutateArguments()
	callMayMutateGlobalsOrThrow := call.MayMutateGlobalStateOrThrow()

	for index := len(parameters) - 1; index >= 0; index-- {
		p := parameters[index]
		if p.ShouldRemove() {
			o.EliminateCallTargetArgAt(target, index)

			// Inlining parameters into a function may cause the function call to newly mutate global
			// state if either
			//  1) the parameter itself mutates global state
			//  2) the function may mutate its arguments, and the argument being inlined is mutable.
			//     We're deliberately conservative here - possibly the argument being inlined is
			//     not global state or not mutated in practice - because it's difficult to test all
			//     the edge cases.
			if !callMayMutateGlobalsOrThrow &&
				(p.HasSideEffects() || // case (1)
					(callMayMutateArgs && !NodeUtil_IsImmutableValue(p.GetArg()))) { // case (2)
				callMayMutateGlobalsOrThrow = true
				call.SetSideEffectFlagsObject(
					rhino.NewNode_SideEffectFlagsWithValue(call.GetSideEffectFlags()).SetMutatesGlobalState())
			}
		}
	}
}

// OptimizeParameters_Parameter is a simple container class that keeps
// tracks of a parameter and whether it should be removed.
type OptimizeParameters_Parameter struct {
	arg               *rhino.Node
	shouldRemove      bool
	hasSideEffects    bool
	canBeSideEffected bool
	mayBeUndefined    bool
}

// Java: OptimizeParameters.Parameter(Node, boolean)
func NewOptimizeParameters_Parameter(arg *rhino.Node, shouldRemove bool) *OptimizeParameters_Parameter {
	return &OptimizeParameters_Parameter{shouldRemove: shouldRemove, arg: arg}
}

// Java: OptimizeParameters.Parameter.getArg()
func (p *OptimizeParameters_Parameter) GetArg() *rhino.Node {
	return p.arg
}

// Java: OptimizeParameters.Parameter.shouldRemove()
func (p *OptimizeParameters_Parameter) ShouldRemove() bool {
	return p.shouldRemove
}

// Java: OptimizeParameters.Parameter.setShouldRemove(boolean)
func (p *OptimizeParameters_Parameter) SetShouldRemove(value bool) {
	p.shouldRemove = value
}

// Java: OptimizeParameters.Parameter.setHasSideEffects(boolean)
func (p *OptimizeParameters_Parameter) SetHasSideEffects(hasSideEffects bool) {
	p.hasSideEffects = hasSideEffects
}

// Java: OptimizeParameters.Parameter.hasSideEffects()
func (p *OptimizeParameters_Parameter) HasSideEffects() bool {
	return p.hasSideEffects
}

// Java: OptimizeParameters.Parameter.setCanBeSideEffected(boolean)
func (p *OptimizeParameters_Parameter) SetCanBeSideEffected(canBeSideEffected bool) {
	p.canBeSideEffected = canBeSideEffected
}

// Java: OptimizeParameters.Parameter.canBeSideEffected()
func (p *OptimizeParameters_Parameter) CanBeSideEffected() bool {
	return p.canBeSideEffected
}

// Java: OptimizeParameters.Parameter.setMayBeUndefined(boolean)
func (p *OptimizeParameters_Parameter) SetMayBeUndefined(mayBeUndefined bool) {
	p.mayBeUndefined = mayBeUndefined
}

// Java: OptimizeParameters.addRestVariableToFunction(Node, Node, Node)
func (o *OptimizeParameters) AddRestVariableToFunction(function *rhino.Node, lhs *rhino.Node, value *rhino.Node) {
	jutil.CheckState(lhs.GetParent() == nil)
	o.AddVariableDeclarationToFunction(function, lhs, value)
}

// AddVariableDeclarationToFunction adds a variable to the function block
// after all hoisted functions.
//
// function is a function node; lhs the lhs expression; value the initial
// value of the variable.
//
// Java: OptimizeParameters.addVariableDeclarationToFunction(Node, Node, Node)
func (o *OptimizeParameters) AddVariableDeclarationToFunction(function *rhino.Node, lhs *rhino.Node, value *rhino.Node) {
	jutil.CheckStateMsg(value.GetParent() == nil, "Value must be detached")
	jutil.CheckStateMsg(lhs != nil, "formal parameter being declared must not be null")
	jutil.CheckStateMsg(lhs.GetParent() == nil, "formal parameter being declared must be detached")
	var stmts []*rhino.Node
	if lhs.IsDestructuringPattern() {
		o.RewriteDestructuringPattern(lhs, value, &stmts)
	} else {
		stmts = append(stmts, NodeUtil_NewVarNodeWithLhs(lhs, value))
	}
	o.InsertStatements(function, stmts)
}

// AddExpressionToFunction adds an expression to the function block after
// all hoisted functions.
//
// Java: OptimizeParameters.addExpressionToFunction(Node, Node)
func (o *OptimizeParameters) AddExpressionToFunction(function *rhino.Node, value *rhino.Node) {
	var stmts []*rhino.Node
	stmts = append(stmts, rhino.IR_ExprResult(value).Srcref(value))
	o.InsertStatements(function, stmts)
}

// InsertStatements inserts the statements at the beginning of the function
// body, but after any hoisted function declarations so that the AST stays
// normalized.
//
// Java: OptimizeParameters.insertStatements(Node, List<Node>)
func (o *OptimizeParameters) InsertStatements(function *rhino.Node, stmts []*rhino.Node) {
	if len(stmts) == 0 {
		return
	}
	block := NodeUtil_GetFunctionBody(function)
	insertionPoint := NodeUtil_GetInsertionPointAfterAllInnerFunctionDeclarations(block)
	if insertionPoint == nil {
		o.AddStatementsToBack(block, stmts)
	} else {
		o.AddStatementsBefore(insertionPoint, stmts)
	}
	o.compiler.ReportChangeToEnclosingScope(stmts[0])
}

// Java: OptimizeParameters.addStatementsToBack(Node, List<Node>)
func (o *OptimizeParameters) AddStatementsToBack(block *rhino.Node, stmts []*rhino.Node) {
	for _, stmt := range stmts {
		block.AddChildToBack(stmt)
	}
}

// Java: OptimizeParameters.addStatementsBefore(Node, List<Node>)
func (o *OptimizeParameters) AddStatementsBefore(insertionPoint *rhino.Node, stmts []*rhino.Node) {
	// Lists.reverse(stmts)
	for i := len(stmts) - 1; i >= 0; i-- {
		stmts[i].InsertBefore(insertionPoint)
	}
}

// EliminateParamsAfter removes all formal parameters starting at argIndex.
//
// Java: OptimizeParameters.eliminateParamsAfter(Node, int)
func (o *OptimizeParameters) EliminateParamsAfter(fnNode *rhino.Node, argIndex int) {
	formalArgPtr := NodeUtil_GetFunctionParameters(fnNode).GetFirstChild()
	for argIndex != 0 && formalArgPtr != nil {
		formalArgPtr = formalArgPtr.GetNext()
		argIndex--
	}
	o.EliminateParamsAfterNode(fnNode, formalArgPtr)
}

// Java: OptimizeParameters.eliminateParamsAfter(Node, Node)
func (o *OptimizeParameters) EliminateParamsAfterNode(fnNode *rhino.Node, formal *rhino.Node) {
	if formal != nil {
		// Keep the args in the same order, do the last first.
		o.EliminateParamsAfterNode(fnNode, formal.GetNext())
		var stmts []*rhino.Node
		formal.Detach()
		if formal.IsRest() {
			jutil.CheckState(formal.GetNext() == nil)
			lhs := formal.RemoveFirstChild()
			value := rhino.IR_Arraylit().Srcref(formal)
			o.AddVariableDeclarationToFunction(fnNode, lhs, value)
		} else if formal.IsDefaultValue() {
			lhs := formal.RemoveFirstChild()
			value := formal.GetLastChild().Detach()
			o.AddVariableDeclarationToFunction(fnNode, lhs, value)
		} else if formal.IsDestructuringPattern() {
			// Destructuring declarations must have an rhs.
			// NOTE: assigning undefined will cause an exception at runtime if this code is evaluated,
			// which matches the behavior of the input code. It's also possible this method will never
			// be evaluated at runtime. This pass jointly optimizes all methods with the same name.
			value := NodeUtil_NewUndefinedNode(formal)
			o.RewriteDestructuringPattern(formal, value, &stmts)
			o.InsertStatements(fnNode, stmts)
		} else {
			stmts = append(stmts, rhino.IR_Var(formal).SrcrefIfMissing(formal))
			o.InsertStatements(fnNode, stmts)
		}
	}
}

// RewriteDestructuringPattern appends the statements to *stmts.
//
// Go-only: language-adapter — the List<Node> that Java mutates is passed as
// a pointer to the Go slice.
//
// Java: OptimizeParameters.rewriteDestructuringPattern(Node, Node, List<Node>)
func (o *OptimizeParameters) RewriteDestructuringPattern(
	destructuringPattern *rhino.Node, value *rhino.Node, stmts *[]*rhino.Node) {
	jutil.CheckState(destructuringPattern.IsDestructuringPattern())
	jutil.CheckStateMsg(
		!destructuringPattern.HasParent(), "Formal parameter must be detached for rewriting")

	NodeUtil_VisitLhsNodesInDestructuringPattern(
		destructuringPattern,
		func(name *rhino.Node) {
			// Add a declaration outside the array pattern for the given name.
			jutil.CheckStateMsg1(
				name.IsName(),
				// TODO(rishipal): Is this always true? What about var [{a}] = [{a:4}]; ?
				"lhs in destructuring declaration should be a simple name. (%s)",
				name)
			newName := rhino.IR_Name(name.GetString()).Srcref(name)
			newVar := rhino.IR_Var(newName).Srcref(name)
			*stmts = append(*stmts, newVar)
		})
	expr :=
		rhino.IR_ExprResult(rhino.IR_Assign(destructuringPattern, value).Srcref(destructuringPattern)).
			Srcref(destructuringPattern)
	*stmts = append(*stmts, expr)
}

// EliminateCallTargetArgAt eliminates the parameter from a function call.
//
// Java: OptimizeParameters.eliminateCallTargetArgAt(Node, int)
func (o *OptimizeParameters) EliminateCallTargetArgAt(ref *rhino.Node, argIndex int) {
	callArgNode := OptimizeCalls_ReferenceMap_GetArgumentForCallOrNewOrDotCall(ref, argIndex)
	if callArgNode != nil {
		NodeUtil_DeleteNode(callArgNode, o.compiler)
	}
}
