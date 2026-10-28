// Port of com/google/javascript/jscomp/OptimizeParameters.java @ closure-compiler v20261021

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

	"github.com/paullepoulpe/jsgo/internal/ast"
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
	decisionsLog diagnostic.LogFile
	fnNodes      []rhino.Node
	defNodes     []rhino.Node
	ttlList      []rhino.Node
	fixedParams  []optimizeParametersParam
}

type optimizeParametersParam struct {
	arg               rhino.Node
	shouldRemove      bool
	hasSideEffects    bool
	canBeSideEffected bool
	mayBeUndefined    bool
}

// Java: OptimizeParameters(AbstractCompiler)
func NewOptimizeParameters(compiler AbstractCompiler) *OptimizeParameters {
	compiler = jutil.CheckNotNull(compiler)
	return &OptimizeParameters{
		compiler:    compiler,
		astAnalyzer: compiler.GetAstAnalyzer(),
	}
}

// Process is CompilerPass.process.
//
// Go-only: language-adapter — the two Java overloads process(Node, Node)
// (CompilerPass) and process(Node, Node, ReferenceMap)
// (CallGraphCompilerPass) cannot share a Go name; the second one is
// ProcessWithReferenceMap, and AsCallGraphCompilerPass adapts it.
//
// Java: OptimizeParameters.process(Node, Node)
func (o *OptimizeParameters) Process(externs rhino.Node, root rhino.Node) {
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
func (o *OptimizeParameters) ProcessWithReferenceMap(externs rhino.Node, root rhino.Node, refMap *OptimizeCalls_ReferenceMap) {
	decisionsLog := o.compiler.CreateOrReopenIndexedLog("OptimizeParameters", "decisions.log")
	defer decisionsLog.Close()
	defer func() {
		o.decisionsLog = diagnostic.LogFile{}
		clear(o.fnNodes)
		o.fnNodes = o.fnNodes[:0]
		clear(o.defNodes)
		o.defNodes = o.defNodes[:0]
		clear(o.ttlList)
		o.ttlList = o.ttlList[:0]
		clear(o.fixedParams)
		o.fixedParams = o.fixedParams[:0]
	}()
	o.decisionsLog = decisionsLog
	o.globalScope = refMap.GetGlobalScope()

	var toOptimize [][]rhino.Node

	for key, refs := range refMap.GetNameReferences() {
		if o.analyzeCandidate("name", key, refs) {
			toOptimize = append(toOptimize, refs)
		}
	}

	for key, refs := range refMap.GetPropReferences() {
		if o.analyzeCandidate("property", key, refs) {
			toOptimize = append(toOptimize, refs)
		}
	}

	// NOTE: Optimizations must keep ReferenceMap consistent:
	// 1. Optional-arg elimination only moves parameter names/defaults into bodies.
	// 2. Constant-arg elimination only runs on single-definition functions, moving
	//    arguments from call sites into the body.
	// 3. Unused-arg elimination may delete references and therefore runs last.
	for _, refs := range toOptimize {
		o.tryEliminateOptionalArgs(refs)
	}

	for _, refs := range toOptimize {
		o.tryEliminateConstantArgs(refs)
	}

	optimizer := optimizeParametersUnusedOptimizer{outer: o}
	for _, refs := range toOptimize {
		optimizer.tryEliminateUnusedArgs(refs)
	}
	optimizer.applyChanges()
}

type optimizeParametersUnusedOptimizer struct {
	outer             *OptimizeParameters
	toRemove          []rhino.Node
	toReplaceWithZero []rhino.Node
	used              jutil.BitSet
	unused            jutil.BitSet
	unremovable       jutil.BitSet
}

// Java: OptimizeParameters.UnusedParameterOptimizer.tryEliminateUnusedArgs(ArrayList<Node>)
func (u *optimizeParametersUnusedOptimizer) tryEliminateUnusedArgs(refs []rhino.Node) {
	fns := OptimizeCalls_ReferenceMap_AppendFunctionNodes(u.outer.fnNodes[:0], refs)
	u.outer.fnNodes = fns
	jutil.CheckState(len(fns) > 0)

	maxFormalsCount := 0
	lowestUsedRest := math.MaxInt32
	used := &u.used
	used.ClearAll()
	for _, fnNode := range fns {
		fn := ast.MustFunction(fnNode)
		index := -1
		for param := range fn.Params().Params() {
			index++
			if !param.IsUnused() {
				used.Set(index)
				if param.Node().IsRest() {
					lowestUsedRest = min(lowestUsedRest, index)
					if lowestUsedRest == 0 {
						// all parameters are used.
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

	if lowestUsedRest < maxFormalsCount && unused.Cardinality() == 0 {
		return
	}

	unremovable := &u.unremovable
	unremovable.CopyFrom(used)

	lowestSpread := math.MaxInt32
	for _, n := range refs {
		if OptimizeCalls_ReferenceMap_IsNormalOrOptChainCallOrNewTarget(n) {
			if optimizeParameters_hasSpreadInFirstCallParam(n) {
				return
			}
			param := OptimizeCalls_ReferenceMap_GetFirstArgumentForCallOrNewOrDotCall(n)
			paramIndex := 0
			for !param.IsNil() {
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

	if lowestSpread < maxFormalsCount {
		unremovable.SetRange(lowestSpread, maxFormalsCount)
	}

	removeAllAfterIndex := math.MaxInt32
	if lowestUsedRest == math.MaxInt32 {
		removeAllAfterIndex = maxFormalsCount - 1
	}

	for _, n := range refs {
		if OptimizeCalls_ReferenceMap_IsNormalOrOptChainCallOrNewTarget(n) && !optimizeParameters_alreadyRemoved(n) {
			arg := OptimizeCalls_ReferenceMap_GetFirstArgumentForCallOrNewOrDotCall(n)
			u.recordRemovalCallArguments(lowestUsedRest, removeAllAfterIndex, unused, unremovable, arg, 0)
		}
	}

	for _, fnNode := range fns {
		firstParam, _ := ast.MustFunction(fnNode).Params().FirstParam()
		u.removeUnusedFunctionParameters(unremovable, firstParam.Node(), 0)
	}
}

// Java: OptimizeParameters.UnusedParameterOptimizer.recordRemovalCallArguments(int, int, BitSet, BitSet, Node, int)
func (u *optimizeParametersUnusedOptimizer) recordRemovalCallArguments(
	firstRestIndex int,
	removeAllAfterIndex int,
	unused *jutil.BitSet,
	unremovable *jutil.BitSet,
	arg rhino.Node,
	index int,
) bool {
	if arg.IsNil() {
		return true
	}
	if index > removeAllAfterIndex {
		return u.removeArgAndFollowing(arg)
	}
	if arg.IsSpread() {
		return false
	}

	removedAllTrailing := u.recordRemovalCallArguments(
		firstRestIndex, removeAllAfterIndex, unused, unremovable, arg.GetNext(), index+1)
	if index < firstRestIndex {
		isRemovableTrailingUndefined := NodeUtil_IsUndefined(arg) && removedAllTrailing
		if isRemovableTrailingUndefined && !u.outer.astAnalyzer.MayHaveSideEffects(arg) {
			u.toRemove = append(u.toRemove, arg)
			return true
		}

		if unused.Get(index) {
			hasSideEffects := isRemovableTrailingUndefined || u.outer.astAnalyzer.MayHaveSideEffects(arg)
			if !hasSideEffects {
				if unremovable.Get(index) {
					u.recordReplaceWithZero(arg)
				} else {
					u.toRemove = append(u.toRemove, arg)
					return removedAllTrailing
				}
			} else {
				u.removeFromCommaIfNoSideEffects(arg)
			}
		}
	}
	return false
}

// Java: OptimizeParameters.UnusedParameterOptimizer.removeFromCommaIfNoSideEffects(Node)
func (u *optimizeParametersUnusedOptimizer) removeFromCommaIfNoSideEffects(arg rhino.Node) {
	if comma, ok := ast.AsBinaryOf(arg, rhino.Token_COMMA); ok {
		right := comma.Right().Node()
		if !u.outer.astAnalyzer.MayHaveSideEffects(right) {
			u.recordReplaceWithZero(right)
		} else {
			u.removeFromCommaIfNoSideEffects(right)
		}
	}
}

// Java: OptimizeParameters.UnusedParameterOptimizer.recordReplaceWithZero(Node)
func (u *optimizeParametersUnusedOptimizer) recordReplaceWithZero(arg rhino.Node) {
	if num, ok := ast.AsNumber(arg); !ok || num.Value() != 0 {
		u.toReplaceWithZero = append(u.toReplaceWithZero, arg)
	}
}

// Java: OptimizeParameters.UnusedParameterOptimizer.removeArgAndFollowing(Node)
func (u *optimizeParametersUnusedOptimizer) removeArgAndFollowing(arg rhino.Node) bool {
	if !arg.IsNil() {
		removedAll := u.removeArgAndFollowing(arg.GetNext())
		if !u.outer.astAnalyzer.MayHaveSideEffects(arg) {
			u.toRemove = append(u.toRemove, arg)
			return removedAll
		}
		return false
	}
	return true
}

// Java: OptimizeParameters.UnusedParameterOptimizer.removeUnusedFunctionParameters(BitSet, Node, int)
func (u *optimizeParametersUnusedOptimizer) removeUnusedFunctionParameters(unremovable *jutil.BitSet, param rhino.Node, index int) {
	if !param.IsNil() {
		u.removeUnusedFunctionParameters(unremovable, param.GetNext(), index+1)
		if !unremovable.Get(index) {
			jutil.CheckState(param.IsName()) // update for ES6
			u.outer.compiler.ReportChangeToEnclosingScope(param)
			param.Detach()
		}
	}
}

// Java: OptimizeParameters.UnusedParameterOptimizer.applyChanges()
func (u *optimizeParametersUnusedOptimizer) applyChanges() {
	for _, n := range u.toRemove {
		if optimizeParameters_alreadyRemoved(n) {
			continue
		}
		u.outer.compiler.ReportChangeToEnclosingScope(n)
		n.Detach()
		NodeUtil_MarkFunctionsDeleted(n, u.outer.compiler)
	}
	for _, n := range u.toReplaceWithZero {
		num, isNum := ast.AsNumber(n)
		jutil.CheckState(!isNum || num.Value() != 0.0)
		if optimizeParameters_alreadyRemoved(n) {
			continue
		}
		u.outer.compiler.ReportChangeToEnclosingScope(n)
		n.ReplaceWith(rhino.IR_Number(0).Srcref(n))
		NodeUtil_MarkFunctionsDeleted(n, u.outer.compiler)
	}
}

// Java: OptimizeParameters.alreadyRemoved(Node)
func optimizeParameters_alreadyRemoved(n rhino.Node) bool {
	parent := n.GetParent()
	for !parent.IsNil() {
		n = parent
		parent = n.GetParent()
	}
	return !n.IsRoot()
}

// analyzeCandidate reports whether a symbol's references are safe for parameter
// optimization (converting any eligible tagged template literals to normal calls).
//
// Java: OptimizeParameters.analyzeCandidate(String, String, ArrayList<Node>)
func (o *OptimizeParameters) analyzeCandidate(refKind string, key string, refs []rhino.Node) bool {
	if !OptimizeCalls_MayBeOptimizableName(o.compiler, key) {
		if o.decisionsLog.IsLogging() {
			o.decisionsLog.LogFormat("%s\t%s\tnot an optimizable name", refKind, key)
		}
		return false
	}

	definitions := o.defNodes[:0]
	ttls := o.ttlList[:0]
	seenCandidateUse := false
	for _, n := range refs {
		if OptimizeCalls_ReferenceMap_IsNormalOrOptChainCallOrNewTarget(n) {
			seenCandidateUse = true
		} else if parent := n.GetParent(); parent.IsTaggedTemplateLit() {
			ttls = append(ttls, parent)
		} else if o.isCandidateDefinition(n) {
			definitions = append(definitions, n)
		} else if !OptimizeCalls_IsAllowedReference(n) {
			if o.decisionsLog.IsLogging() {
				o.decisionsLog.LogFormat("%s\t%s\tnot an allowed reference: %s", refKind, key, n.GetLocation())
			}
			o.defNodes = definitions
			o.ttlList = ttls
			return false
		}
	}
	o.defNodes = definitions
	o.ttlList = ttls

	if len(definitions) == 0 {
		if o.decisionsLog.IsLogging() {
			o.decisionsLog.LogFormat("%s\t%s\tno definition found", refKind, key)
		}
		return false
	}

	if len(ttls) != 0 {
		functionNodes := OptimizeCalls_ReferenceMap_AppendFunctionNodes(o.fnNodes[:0], definitions)
		o.fnNodes = functionNodes

		for _, functionNode := range functionNodes {
			if firstParam, ok := ast.MustFunction(functionNode).Params().FirstParam(); ok && !firstParam.IsUnused() {
				if o.decisionsLog.IsLogging() {
					o.decisionsLog.LogFormat(
						"%s\t%s\twill not optimize parameters for tagged template literals", refKind, key)
				}
				return false
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
		return false
	}

	if len(ttls) > 0 {
		for _, ttlNode := range ttls {
			o.convertTaggedTemplateLiteralToNormalCall(ast.MustTaggedTemplateLit(ttlNode))
		}
		clear(ttls)
		o.ttlList = ttls[:0]
	}
	return true
}

// convertTaggedTemplateLiteralToNormalCall converts "tagFunction`template ${1} literal ${2}`"
// to `tagFunction([], 1, 2)`.
//
// Java: OptimizeParameters.convertTaggedTemplateLiteralToNormalCall(Node)
func (o *OptimizeParameters) convertTaggedTemplateLiteralToNormalCall(ttl ast.TaggedTemplateLit) {
	callee := ttl.Tag()
	tpl := ttl.Template()
	firstArg := rhino.IR_Arraylit().Srcref(tpl.Node())
	callee.Detach()
	callNode := NodeUtil_NewCallNode(callee.Node(), firstArg).Srcref(ttl.Node())
	for part := range ast.ReadAhead(tpl.FirstPart()) {
		if sub, ok := part.AsTemplateLitSub(); ok {
			arg := sub.Expr().Detach().Node()
			callNode.AddChildToBack(arg)
		}
	}
	ttl.Node().ReplaceWith(callNode)
	o.compiler.ReportChangeToEnclosingScope(callNode)
}

// Java: OptimizeParameters.isCandidateDefinition(Node)
func (o *OptimizeParameters) isCandidateDefinition(n rhino.Node) bool {
	parent := n.GetParent()

	var functionExpr rhino.Node
	if parent.IsFunction() && NodeUtil_IsFunctionDeclaration(parent) {
		functionExpr = parent
	} else if OptimizeCalls_ReferenceMap_IsSimpleAssignmentTarget(n) {
		functionExpr = ast.MustAssign(parent).Value().Node()
	} else if name, ok := ast.AsName(n); ok && n.HasChildren() {
		init, _ := name.Init()
		functionExpr = init.Node()
	} else if m, ok := ast.AsMemberFunctionDef(n); ok && parent.IsClassMembers() {
		functionExpr = m.Function().Node()
	} else if cls, ok := ast.AsClass(parent); ok && cls.NameNode() == n {
		functionExpr = parent
	} else {
		return false
	}

	return o.allDefinitionsAreCandidateFunctions(functionExpr)
}

// Java: OptimizeParameters.classContainsClassFieldWithRHSSideEffects(Node)
func (o *OptimizeParameters) classContainsClassFieldWithRHSSideEffects(cls ast.Class) bool {
	for member := range cls.Body().Members() {
		if fd, ok := member.AsFieldDef(); ok && !fd.IsStatic() {
			if init, hasVal := fd.Value(); hasVal && o.astAnalyzer.MayHaveSideEffects(init.Node()) {
				return true
			}
		}
	}
	return false
}

// Java: OptimizeParameters.allDefinitionsAreCandidateFunctions(Node)
func (o *OptimizeParameters) allDefinitionsAreCandidateFunctions(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_CLASS:
		if NodeUtil_IsNamedClassExpression(n) {
			return false
		}
		cls := ast.MustClass(n)
		ctor, hasCtor := cls.Constructor()
		if !hasCtor || o.classContainsClassFieldWithRHSSideEffects(cls) {
			return false
		}
		fn := ctor.Function()
		return !NodeUtil_DoesFunctionReferenceOwnArgumentsObject(fn.Node()) &&
			!optimizeParameters_mayReferenceParamBeforeBody(fn) &&
			!optimizeParameters_isUsedViaDotConstructor(fn.Node())

	case rhino.Token_FUNCTION:
		fn := ast.MustFunction(n)
		return !NodeUtil_IsNamedFunctionExpression(n) &&
			!NodeUtil_DoesFunctionReferenceOwnArgumentsObject(n) &&
			!optimizeParameters_mayReferenceParamBeforeBody(fn) &&
			!optimizeParameters_isUsedViaDotConstructor(n)

	case rhino.Token_CAST:
		return o.allDefinitionsAreCandidateFunctions(ast.MustCast(n).Expr().Node())
	case rhino.Token_COMMA:
		return o.allDefinitionsAreCandidateFunctions(ast.MustBinary(n).Right().Node())
	case rhino.Token_HOOK:
		h := ast.MustHook(n)
		return o.allDefinitionsAreCandidateFunctions(h.Then().Node()) &&
			o.allDefinitionsAreCandidateFunctions(h.Else().Node())
	case rhino.Token_OR, rhino.Token_AND, rhino.Token_COALESCE:
		b := ast.MustBinary(n)
		return o.allDefinitionsAreCandidateFunctions(b.Left().Node()) &&
			o.allDefinitionsAreCandidateFunctions(b.Right().Node())
	default:
		return false
	}
}

// Java: OptimizeParameters.isUsedViaDotConstructor(Node)
func optimizeParameters_isUsedViaDotConstructor(function rhino.Node) bool {
	docInfo := NodeUtil_GetBestJSDocInfo(function)
	return docInfo != nil && docInfo.IsUsedViaDotConstructor()
}

// Java: OptimizeParameters.mayReferenceParamBeforeBody(Node)
func optimizeParameters_mayReferenceParamBeforeBody(fn ast.Function) bool {
	paramList := fn.Params()
	if !paramList.HasParams() {
		return false
	}

	allSimpleNames := true
	count := 0
	for p := range paramList.Params() {
		if !p.Node().IsName() || p.Node().HasChildren() {
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
			for a := paramList.Node().GetFirstChild(); !a.IsNil(); a = a.GetNext() {
				aID := a.GetStringID()
				for b := a.GetNext(); !b.IsNil(); b = b.GetNext() {
					if aID == b.GetStringID() {
						return true
					}
				}
			}
			return false
		}
	}

	var seen optimizeParametersParamNameTracker
	return seen.scan(paramList.Node())
}

type optimizeParametersParamNameEntry struct {
	id    uint32
	first rhino.Node
	multi bool
}

type optimizeParametersParamNameTracker struct {
	inline  [16]optimizeParametersParamNameEntry
	entries []optimizeParametersParamNameEntry
	byID    map[uint32]int
}

func (t *optimizeParametersParamNameTracker) scan(n rhino.Node) bool {
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if t.scan(c) {
			return true
		}
	}
	if !n.IsName() {
		return false
	}
	id := n.GetStringID()
	if t.entries == nil {
		t.entries = t.inline[:0]
	}
	idx := -1
	if t.byID != nil {
		if existing, ok := t.byID[id]; ok {
			idx = existing
		}
	} else {
		for i := range t.entries {
			if t.entries[i].id == id {
				idx = i
				break
			}
		}
	}
	if idx < 0 {
		newIdx := len(t.entries)
		t.entries = append(t.entries, optimizeParametersParamNameEntry{id: id, first: n})
		if t.byID != nil {
			t.byID[id] = newIdx
		} else if len(t.entries) == len(t.inline) {
			t.byID = make(map[uint32]int, 32)
			for i := range t.entries {
				t.byID[t.entries[i].id] = i
			}
		}
		return false
	}
	e := &t.entries[idx]
	if !e.multi {
		e.multi = true
		if NodeUtil_IsLValue(e.first) {
			return true
		}
	}
	return NodeUtil_IsLValue(n)
}

// Java: OptimizeParameters.tryEliminateOptionalArgs(ArrayList<Node>)
func (o *OptimizeParameters) tryEliminateOptionalArgs(refs []rhino.Node) {
	maxArgs := -1

	for _, n := range refs {
		if OptimizeCalls_ReferenceMap_IsNormalOrOptChainCallOrNewTarget(n) {
			if optimizeParameters_hasSpreadInFirstCallParam(n) {
				return
			}
			numArgs := 0
			firstArg := OptimizeCalls_ReferenceMap_GetFirstArgumentForCallOrNewOrDotCall(n)
			for c := firstArg; !c.IsNil(); c = c.GetNext() {
				numArgs++
				if c.IsSpread() {
					return
				}
			}
			if numArgs > maxArgs {
				maxArgs = numArgs
			}
		}
	}

	fns := OptimizeCalls_ReferenceMap_AppendFunctionNodes(o.fnNodes[:0], refs)
	o.fnNodes = fns
	for _, fn := range fns {
		o.eliminateParamsAfter(ast.MustFunction(fn), maxArgs)
	}
}

// Java: OptimizeParameters.tryEliminateConstantArgs(ArrayList<Node>)
func (o *OptimizeParameters) tryEliminateConstantArgs(refs []rhino.Node) {
	parameters, ok := o.findFixedArguments(refs)
	if !ok {
		return
	}

	fnNode, fnCount := optimizeCalls_getSingleFunctionNode(refs)
	if fnCount > 1 {
		return
	}
	if fnCount == 0 {
		panic(jutil.NewNoSuchElementException())
	}

	fn := ast.MustFunction(fnNode)
	if !optimizeParameters_adjustForConstraints(fn, parameters) {
		return
	}

	for _, n := range refs {
		if !optimizeParameters_alreadyRemoved(n) && OptimizeCalls_ReferenceMap_IsNormalOrOptChainCallOrNewTarget(n) {
			o.optimizeCallSite(parameters, n)
		}
	}

	o.optimizeFunctionDefinition(parameters, fn)
}

// Java: OptimizeParameters.findFixedArguments(ArrayList<Node>)
func (o *OptimizeParameters) findFixedArguments(refs []rhino.Node) ([]optimizeParametersParam, bool) {
	clear(o.fixedParams)
	o.fixedParams = o.fixedParams[:0]
	firstCall := true

	continueLooking := false
	for _, n := range refs {
		if OptimizeCalls_ReferenceMap_IsNormalOrOptChainCallOrNewTarget(n) {
			if optimizeParameters_hasSpreadInFirstCallParam(n) {
				continueLooking = false
				break
			}
			cur := OptimizeCalls_ReferenceMap_GetFirstArgumentForCallOrNewOrDotCall(n)
			if firstCall {
				continueLooking = o.buildInitialParameterList(cur)
				firstCall = false
			} else {
				continueLooking = o.findFixedParameters(cur)
			}
			if !continueLooking {
				break
			}
		}
	}

	if continueLooking {
		return o.fixedParams, true
	}
	return nil, false
}

// Java: OptimizeParameters.hasSpreadInFirstCallParam(Node)
func optimizeParameters_hasSpreadInFirstCallParam(n rhino.Node) bool {
	call := OptimizeCalls_ReferenceMap_GetCallOrNewNodeForTarget(n)
	if NodeUtil_IsFunctionObjectCall(call) {
		firstArg, ok := ast.MustCallOrNew(call).FirstArg()
		return ok && firstArg.Node().IsSpread()
	}
	return false
}

// Java: OptimizeParameters.adjustForConstraints(Node, List<Parameter>)
func optimizeParameters_adjustForConstraints(fn ast.Function, parameters []optimizeParametersParam) bool {
	info := NodeUtil_GetBestJSDocInfo(fn.Node())
	if info != nil && info.IsNoInline() {
		return false
	}

	paramList := fn.Params()
	lastFormal, hasLastFormal := paramList.LastParam()
	restIndex := math.MaxInt32
	lastNonRestFormal := paramList.NumParams() - 1
	formal := lastFormal.Node()
	if hasLastFormal && formal.IsRest() {
		restIndex = lastNonRestFormal
		lastNonRestFormal--
		formal = formal.GetPrevious()
	}

	anyMovable := false
	seenUnmovableSideEffects := false
	seenUnmoveableSideEffected := false
	allRestValueRemovable := true
	for i := len(parameters) - 1; i >= 0; i-- {
		current := &parameters[i]

		if i <= lastNonRestFormal {
			if formal.IsDefaultValue() && current.mayBeUndefined {
				current.shouldRemove = false
			}
			formal = formal.GetPrevious()
		}

		if current.shouldRemove &&
			((seenUnmovableSideEffects && current.canBeSideEffected) ||
				(seenUnmoveableSideEffected && current.hasSideEffects)) {
			current.shouldRemove = false
		}

		if i >= restIndex {
			if allRestValueRemovable {
				if !current.shouldRemove {
					anyMovable = false
					allRestValueRemovable = false
					for j := i + 1; j < len(parameters); j++ {
						p := &parameters[j]
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

// Java: OptimizeParameters.findFixedParameters(List<Parameter>, Node)
func (o *OptimizeParameters) findFixedParameters(cur rhino.Node) bool {
	anyMovable := false
	index := 0
	for !cur.IsNil() {
		if index >= len(o.fixedParams) {
			p := optimizeParametersParam{arg: cur}
			o.setParameterSideEffectInfo(&p, cur)
			o.fixedParams = append(o.fixedParams, p)
		} else {
			p := &o.fixedParams[index]
			if p.shouldRemove {
				if !cur.IsEquivalentTo(p.arg) {
					p.shouldRemove = false
				} else {
					anyMovable = true
				}
			}
			p.hasSideEffects = p.hasSideEffects || o.astAnalyzer.MayHaveSideEffects(cur)
			p.canBeSideEffected = p.canBeSideEffected || NodeUtil_CanBeSideEffected(cur)
		}
		if cur.IsSpread() {
			break
		}
		cur = cur.GetNext()
		index++
	}

	for ; index < len(o.fixedParams); index++ {
		o.fixedParams[index].shouldRemove = false
	}

	return anyMovable
}

// Java: OptimizeParameters.buildInitialParameterList(List<Parameter>, Node)
func (o *OptimizeParameters) buildInitialParameterList(cur rhino.Node) bool {
	anyMovable := false
	for !cur.IsNil() {
		movable := optimizeParameters_isMovableValue(cur, o.globalScope)
		p := optimizeParametersParam{arg: cur, shouldRemove: movable}
		o.setParameterSideEffectInfo(&p, cur)
		o.fixedParams = append(o.fixedParams, p)
		if movable {
			anyMovable = true
		}
		if cur.IsSpread() {
			break
		}
		cur = cur.GetNext()
	}
	return anyMovable
}

// Java: OptimizeParameters.setParameterSideEffectInfo(Parameter, Node)
func (o *OptimizeParameters) setParameterSideEffectInfo(p *optimizeParametersParam, value rhino.Node) {
	p.hasSideEffects = o.astAnalyzer.MayHaveSideEffects(value)
	p.canBeSideEffected = NodeUtil_CanBeSideEffected(value)
	if !value.IsSpread() {
		p.mayBeUndefined = NodeUtil_MayBeUndefined(value)
	}
}

var optimizeParametersArgumentsID = rhino.RhinoStringPool_InternID("arguments")

// Java: OptimizeParameters.isMovableValue(Node, Scope)
func optimizeParameters_isMovableValue(n rhino.Node, globalScope *Scope) bool {
	switch n.GetToken() {
	case rhino.Token_AWAIT, rhino.Token_YIELD, rhino.Token_THIS, rhino.Token_SUPER,
		rhino.Token_ITER_SPREAD, rhino.Token_OBJECT_SPREAD, rhino.Token_NEW_TARGET, rhino.Token_IMPORT_META:
		return false
	case rhino.Token_FUNCTION:
		return false
	case rhino.Token_NAME:
		id := n.GetStringID()
		if id == optimizeParametersArgumentsID || globalScope.GetVarByID(id) == nil {
			return false
		}
	}

	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if !optimizeParameters_isMovableValue(c, globalScope) {
			return false
		}
	}
	return true
}

// Java: OptimizeParameters.optimizeFunctionDefinition(List<Parameter>, Node)
func (o *OptimizeParameters) optimizeFunctionDefinition(parameters []optimizeParametersParam, fn ast.Function) {
	paramList := fn.Params()
	lastParameter := len(parameters) - 1
	if lastParam, hasLast := paramList.LastParam(); hasLast {
		if rest, isRest := lastParam.AsIterRest(); isRest {
			restIndex := paramList.NumParams() - 1
			if len(parameters) < restIndex ||
				(restIndex < len(parameters) && parameters[restIndex].shouldRemove) {
				value := rhino.IR_Arraylit().Srcref(rest.Node())
				for i := restIndex; i < len(parameters); i++ {
					jutil.CheckState(parameters[i].shouldRemove)
					value.AddChildToBack(parameters[i].arg)
				}
				rest.Detach()
				lhs := rest.Target().Detach().Node()
				o.addVariableDeclarationToFunction(fn, lhs, value)
			}
			lastParameter = min(len(parameters)-1, restIndex-1)
		}
	}

	for i := lastParameter; i >= 0; i-- {
		parameter := &parameters[i]
		if !parameter.shouldRemove {
			continue
		}
		var formalLhs rhino.Node
		if formalParam, ok := paramList.ParamAt(i); ok {
			formalParam.Detach()
			if dv, isDefault := formalParam.AsDefaultValue(); isDefault {
				jutil.CheckState(!parameter.mayBeUndefined)
				formalLhs = dv.Target().Detach().Node()
			} else {
				formalLhs = formalParam.Node()
			}
		}

		if !formalLhs.IsNil() {
			o.addVariableDeclarationToFunction(fn, formalLhs, parameter.arg)
		} else {
			o.insertStatements(fn, []rhino.Node{rhino.IR_ExprResult(parameter.arg).Srcref(parameter.arg)})
		}
	}
}

// Java: OptimizeParameters.optimizeCallSite(List<Parameter>, Node)
func (o *OptimizeParameters) optimizeCallSite(parameters []optimizeParametersParam, target rhino.Node) {
	call := OptimizeCalls_ReferenceMap_GetCallOrNewNodeForTarget(target)
	callMayMutateArgs := call.MayMutateArguments()
	callMayMutateGlobalsOrThrow := call.MayMutateGlobalStateOrThrow()

	for index := len(parameters) - 1; index >= 0; index-- {
		p := &parameters[index]
		if !p.shouldRemove {
			continue
		}
		if callArgNode := OptimizeCalls_ReferenceMap_GetArgumentForCallOrNewOrDotCall(target, index); !callArgNode.IsNil() {
			NodeUtil_DeleteNode(callArgNode, o.compiler)
		}

		if !callMayMutateGlobalsOrThrow &&
			(p.hasSideEffects || (callMayMutateArgs && !NodeUtil_IsImmutableValue(p.arg))) {
			callMayMutateGlobalsOrThrow = true
			call.SetSideEffectFlagsObject(
				rhino.NewNode_SideEffectFlagsWithValue(call.GetSideEffectFlags()).SetMutatesGlobalState())
		}
	}
}

// Java: OptimizeParameters.addVariableDeclarationToFunction(Node, Node, Node)
func (o *OptimizeParameters) addVariableDeclarationToFunction(fn ast.Function, lhs rhino.Node, value rhino.Node) {
	jutil.CheckStateMsg(value.GetParent().IsNil(), "Value must be detached")
	jutil.CheckStateMsg(!lhs.IsNil(), "formal parameter being declared must not be null")
	jutil.CheckStateMsg(lhs.GetParent().IsNil(), "formal parameter being declared must be detached")
	var stmts []rhino.Node
	if lhs.IsDestructuringPattern() {
		o.rewriteDestructuringPattern(lhs, value, &stmts)
	} else {
		stmts = append(stmts, NodeUtil_NewVarNodeWithLhs(lhs, value))
	}
	o.insertStatements(fn, stmts)
}

// Java: OptimizeParameters.insertStatements(Node, List<Node>)
func (o *OptimizeParameters) insertStatements(fn ast.Function, stmts []rhino.Node) {
	if len(stmts) == 0 {
		return
	}
	block := fn.Body().Node()
	insertionPoint := NodeUtil_GetInsertionPointAfterAllInnerFunctionDeclarations(block)
	if insertionPoint.IsNil() {
		for _, stmt := range stmts {
			block.AddChildToBack(stmt)
		}
	} else {
		for i := len(stmts) - 1; i >= 0; i-- {
			stmts[i].InsertBefore(insertionPoint)
		}
	}
	o.compiler.ReportChangeToEnclosingScope(stmts[0])
}

// Java: OptimizeParameters.eliminateParamsAfter(Node, int)
func (o *OptimizeParameters) eliminateParamsAfter(fn ast.Function, argIndex int) {
	formal, ok := fn.Params().ParamAt(argIndex)
	if !ok {
		return
	}
	o.eliminateParamsAfterNode(fn, formal.Node())
}

// Java: OptimizeParameters.eliminateParamsAfter(Node, Node)
func (o *OptimizeParameters) eliminateParamsAfterNode(fn ast.Function, formalNode rhino.Node) {
	if formalNode.IsNil() {
		return
	}
	// Keep the args in the same order: process trailing parameters first.
	o.eliminateParamsAfterNode(fn, formalNode.GetNext())

	formal := ast.MustParam(formalNode)
	formal.Detach()
	if rest, ok := formal.AsIterRest(); ok {
		lhs := rest.Target().Detach().Node()
		value := rhino.IR_Arraylit().Srcref(formalNode)
		o.addVariableDeclarationToFunction(fn, lhs, value)
	} else if dv, ok := formal.AsDefaultValue(); ok {
		target, val := dv.Target(), dv.Value()
		lhs := target.Detach().Node()
		value := val.Detach().Node()
		o.addVariableDeclarationToFunction(fn, lhs, value)
	} else if _, ok := formal.AsPattern(); ok {
		var stmts []rhino.Node
		value := NodeUtil_NewUndefinedNode(formalNode)
		o.rewriteDestructuringPattern(formalNode, value, &stmts)
		o.insertStatements(fn, stmts)
	} else {
		o.insertStatements(fn, []rhino.Node{rhino.IR_Var(formalNode).SrcrefIfMissing(formalNode)})
	}
}

// Java: OptimizeParameters.rewriteDestructuringPattern(Node, Node, List<Node>)
func (o *OptimizeParameters) rewriteDestructuringPattern(
	destructuringPattern rhino.Node, value rhino.Node, stmts *[]rhino.Node) {
	jutil.CheckState(destructuringPattern.IsDestructuringPattern())
	jutil.CheckStateMsg(!destructuringPattern.HasParent(), "Formal parameter must be detached for rewriting")

	NodeUtil_VisitLhsNodesInDestructuringPattern(
		destructuringPattern,
		func(name rhino.Node) {
			jutil.CheckStateMsg1(
				name.IsName(),
				"lhs in destructuring declaration should be a simple name. (%s)",
				name)
			newName := rhino.IR_Name(name.GetString()).Srcref(name)
			newVar := rhino.IR_Var(newName).Srcref(name)
			*stmts = append(*stmts, newVar)
		})
	expr := rhino.IR_ExprResult(rhino.IR_Assign(destructuringPattern, value).Srcref(destructuringPattern)).
		Srcref(destructuringPattern)
	*stmts = append(*stmts, expr)
}
