// Port of com/google/javascript/jscomp/FunctionArgumentInjector.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: the parameters of a function and the arguments of a call are
// read through the typed view (internal/ast), and so are the targets of
// assignments and the declared names; internal parameter/local sets and
// argument replacement maps are keyed by 32-bit strpool IDs instead of
// hashing materialized strings (see function_argument_injector_sets.go).

// FunctionArgumentInjector is a nifty set of functions to deal with the
// issues of replacing function parameters with a set of call argument
// expressions.
type FunctionArgumentInjector struct {
	astAnalyzer *AstAnalyzer

	sentinelUndefinedNode rhino.Node
	templateUndefinedNode rhino.Node
	sentinelThisNode      rhino.Node
	paramBindings         []callParamBinding
	scratchLocals         inlineIDSet
	scratchSummary        calleeSummary
	seScratchSummary      calleeSummary
	referencedAfter       FunctionArgumentInjector_ReferencedAfterSideEffect

	// readOnlySummaries is enabled on parallel Phase 2 workers so lazy side-effect
	// summaries are stored in worker-local sideEffectCache instead of mutating
	// shared FunctionState summaries concurrently.
	readOnlySummaries bool
	sideEffectCache   map[rhino.Node]calleeSideEffectSummary
}

// Java: FunctionArgumentInjector(AstAnalyzer)
func NewFunctionArgumentInjector(astAnalyzer *AstAnalyzer) *FunctionArgumentInjector {
	f := &FunctionArgumentInjector{
		astAnalyzer:           astAnalyzer,
		sentinelUndefinedNode: NodeUtil_NewUndefinedNode(rhino.NilNode),
		templateUndefinedNode: NodeUtil_NewUndefinedNode(rhino.NilNode),
		sentinelThisNode:      rhino.IR_ThisNode(),
	}
	f.referencedAfter.outer = f
	return f
}

func (f *FunctionArgumentInjector) getTemplateUndefinedNode(callNode rhino.Node) rhino.Node {
	if !f.templateUndefinedNode.IsNil() && f.templateUndefinedNode.ResetTemplateUndefinedSrcref(callNode) {
		return f.templateUndefinedNode
	}
	return NodeUtil_NewUndefinedNode(callNode)
}

func (f *FunctionArgumentInjector) injectByID(
	compiler AbstractCompiler,
	node rhino.Node,
	parent rhino.Node,
	replacements *inlineIDNodeMap,
	replaceThis bool) rhino.Node {
	if node.IsName() {
		replacementTemplate := replacements.get(node.GetStringID())
		if !replacementTemplate.IsNil() {
			// This should not be replacing declared names.
			jutil.CheckStateMsg1(!(parent.IsFunction() || parent.IsVar() || parent.IsCatch()), "%s", parent)
			if !parent.IsNil() && node.CanCloneLeafInPlace() && replacementTemplate.CanCloneLeafInPlace() {
				node.OverwriteLeafFromClone(replacementTemplate)
				return node
			}
			// The name may need to be replaced more than once,
			// so we need to clone the node.
			replacement := replacementTemplate.CloneTree()
			node.ReplaceWith(replacement)
			return replacement
		}
	} else if replaceThis && node.IsThis() {
		replacementTemplate := replacements.get(calleeThisMarkerID)
		rhino.CheckNotNull(replacementTemplate)
		if !replacementTemplate.IsThis() {
			if !parent.IsNil() && node.CanCloneLeafInPlace() && replacementTemplate.CanCloneLeafInPlace() {
				node.OverwriteLeafFromClone(replacementTemplate)
				if compiler.GetAstAnalyzer().MayHaveSideEffects(replacementTemplate) {
					replacements.remove(calleeThisMarkerID)
				}
				return node
			}
			// The name may need to be replaced more than once,
			// so we need to clone the node.
			replacement := replacementTemplate.CloneTree()
			node.ReplaceWith(replacement)

			// Remove the value.  This isn't required but it ensures that we won't
			// inject side-effects multiple times as it will trigger the null
			// check above if we do.
			if compiler.GetAstAnalyzer().MayHaveSideEffects(replacementTemplate) {
				replacements.remove(calleeThisMarkerID)
			}

			return replacement
		}
	} else if (node.IsFunction() && !node.IsArrowFunction()) ||
		node.IsMemberFieldDef() ||
		NodeUtil_IsClassStaticBlock(node) {
		// Once we enter another non-arrow function, class field initializer, or static block,
		// the "this" value changes. Don't try to replace it within an inner scope.
		replaceThis = false
	}

	if field, ok := ast.AsComputedFieldDef(node); ok {
		// The key is evaluated where the class is, the value later, in the
		// field's own this.
		f.injectByID(compiler, field.Key().Node(), node, replacements, replaceThis)
		if value, ok := field.Value(); ok { // read after the key was replaced
			f.injectByID(compiler, value.Node(), node, replacements /* replaceThis= */, false)
		}
		return node
	}

	for c := node.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		// We have to reassign c in case it was replaced, because the removed c's
		// getNext() would no longer be correct.
		c = f.injectByID(compiler, c, node, replacements, replaceThis)
	}

	return node
}

type callParamWalkMode uint8

const (
	callParamWalkOwned callParamWalkMode = iota
	callParamWalkTemplate
	callParamWalkReadOnly
)

type callParamBinding struct {
	nameID    uint32
	paramNode rhino.Node
	argNode   rhino.Node
}

func (b callParamBinding) nameStr() string {
	return rhino.RhinoStringPool_StringFromID(b.nameID)
}

func hasDuplicateBindingID(bindings []callParamBinding) bool {
	n := len(bindings)
	if n <= 1 {
		return false
	}
	if n <= inlineSmallCap {
		for i := 1; i < n; i++ {
			id := bindings[i].nameID
			for j := 0; j < i; j++ {
				if bindings[j].nameID == id {
					return true
				}
			}
		}
		return false
	}
	seen := make(map[uint32]struct{}, n)
	for _, b := range bindings {
		if _, dup := seen[b.nameID]; dup {
			return true
		}
		seen[b.nameID] = struct{}{}
	}
	return false
}

func (f *FunctionArgumentInjector) resolveUndefinedArg(
	callNode rhino.Node,
	mode callParamWalkMode,
	undefTemplate *rhino.Node,
) rhino.Node {
	switch mode {
	case callParamWalkOwned:
		return NodeUtil_NewUndefinedNode(callNode)
	case callParamWalkReadOnly:
		return f.sentinelUndefinedNode
	default:
		if undefTemplate.IsNil() {
			*undefTemplate = f.getTemplateUndefinedNode(callNode)
		}
		return *undefTemplate
	}
}

func (f *FunctionArgumentInjector) appendCallParamBindings(
	dst []callParamBinding,
	fnNode rhino.Node,
	callNode rhino.Node,
	safeNameIdSupplier func() string,
	mode callParamWalkMode,
) []callParamBinding {
	rhino.CheckNotNull(fnNode)
	var undefTemplate rhino.Node

	// CALL NODE: [ NAME, ARG1, ARG2, ... ]
	cArg, hasArg := ast.MustNormalOrOptChainCall(callNode).FirstArg()
	if hasArg && NodeUtil_IsFunctionObjectCall(callNode) {
		thisArg := cArg.Node()
		if fnNode.IsArrowFunction() {
			if mode == callParamWalkOwned {
				thisArg = rhino.IR_ThisNode()
			} else {
				thisArg = f.sentinelThisNode
			}
		}
		dst = append(dst, callParamBinding{
			nameID:    calleeThisMarkerID,
			paramNode: FunctionArgumentInjector_THIS_MARKER_NODE,
			argNode:   thisArg,
		})
		cArg, hasArg = cArg.Next()
	} else {
		// 'apply' isn't supported yet.
		jutil.CheckStateMsg1(!NodeUtil_IsFunctionObjectApply(callNode), "%s", callNode)
		dst = append(dst, callParamBinding{
			nameID:    calleeThisMarkerID,
			paramNode: FunctionArgumentInjector_THIS_MARKER_NODE,
			argNode:   f.resolveUndefinedArg(callNode, mode, &undefTemplate),
		})
	}

	// (A plain loop: ranging over Params() with a return in the body costs
	// a closure call per parameter in this hot path.)
	for fnParam, more := ast.MustFunction(fnNode).Params().FirstParam(); more; fnParam, more = fnParam.Next() {
		if rest, ok := fnParam.AsIterRest(); ok {
			restName, ok := rest.Target().AsName()
			if hasArg {
				jutil.CheckStateMsg1(ok, "%s", rest.Target().Node())
				// The rest param takes the remaining arguments, as an array.
				array := rhino.IR_Arraylit()
				array.SrcrefTreeIfMissing(cArg.Node())
				for ; hasArg; cArg, hasArg = cArg.Next() {
					array.AddChildToBack(cArg.Node().CloneTree())
				}
				return append(dst, callParamBinding{
					nameID:    restName.Node().GetStringID(),
					paramNode: restName.Node(),
					argNode:   array,
				})
			}
			jutil.CheckStateMsg1(ok, "%s", fnParam.Node())
			dst = append(dst, callParamBinding{
				nameID:    restName.Node().GetStringID(),
				paramNode: restName.Node(),
				argNode:   rhino.IR_Arraylit().Srcref(fnParam.Node()),
			})
			continue
		}
		name, ok := fnParam.AsName()
		jutil.CheckStateMsg1(ok, "%s", fnParam.Node())
		nameNode := name.Node()
		if hasArg {
			dst = append(dst, callParamBinding{
				nameID:    nameNode.GetStringID(),
				paramNode: nameNode,
				argNode:   cArg.Node(),
			})
			cArg, hasArg = cArg.Next()
		} else {
			dst = append(dst, callParamBinding{
				nameID:    nameNode.GetStringID(),
				paramNode: nameNode,
				argNode:   f.resolveUndefinedArg(callNode, mode, &undefTemplate),
			})
		}
	}

	// Add temp names for arguments that don't have named parameters in the
	// called function.
	for ; hasArg; cArg, hasArg = cArg.Next() {
		uniquePlaceholder := FunctionArgumentInjector_GetUniqueAnonymousParameterName(safeNameIdSupplier)
		if mode == callParamWalkReadOnly {
			dst = append(dst, callParamBinding{
				nameID:    rhino.RhinoStringPool_InternID(uniquePlaceholder),
				paramNode: FunctionArgumentInjector_THIS_MARKER_NODE,
				argNode:   cArg.Node(),
			})
		} else {
			paramNode := rhino.IR_Name(uniquePlaceholder)
			dst = append(dst, callParamBinding{
				nameID:    paramNode.GetStringID(),
				paramNode: paramNode,
				argNode:   cArg.Node(),
			})
		}
	}
	return dst
}

func throwDuplicateParamBinding(bindings []callParamBinding) {
	n := len(bindings)
	for i := n - 1; i >= 0; i-- {
		id := bindings[i].nameID
		for j := i + 1; j < n; j++ {
			if bindings[j].nameID == id {
				k := bindings[i].nameStr()
				laterStr := "{" + bindings[j].paramNode.String() + " " + bindings[j].argNode.String() + "}"
				earlierStr := "{" + bindings[i].paramNode.String() + " " + bindings[i].argNode.String() + "}"
				panic(jutil.NewIllegalArgumentExceptionMsg("Multiple entries with same key: " +
					k + "=" + laterStr + " and " +
					k + "=" + earlierStr))
			}
		}
	}
}

func (f *FunctionArgumentInjector) failDuplicateParamBindings(
	fnNode rhino.Node, callNode rhino.Node, safeNameIdSupplier func() string,
) {
	var bindingsBuf [inlineSmallCap]callParamBinding
	bindings := f.appendCallParamBindings(bindingsBuf[:0], fnNode, callNode, safeNameIdSupplier, callParamWalkOwned)
	throwDuplicateParamBinding(bindings)
}

func (f *FunctionArgumentInjector) fillCallParamNodeMap(
	fnNode rhino.Node,
	callNode rhino.Node,
	safeNameIdSupplier func() string,
	dst *inlineIDNodeMap,
) {
	dst.clear()
	f.paramBindings = f.appendCallParamBindings(f.paramBindings[:0], fnNode, callNode, safeNameIdSupplier, callParamWalkTemplate)
	for _, b := range f.paramBindings {
		if !dst.put(b.nameID, b.argNode).IsNil() {
			f.failDuplicateParamBindings(fnNode, callNode, safeNameIdSupplier)
		}
	}
}

func (f *FunctionArgumentInjector) fillCallParamBindings(
	dst []callParamBinding,
	fnNode rhino.Node,
	callNode rhino.Node,
	safeNameIdSupplier func() string,
) []callParamBinding {
	dst = f.appendCallParamBindings(dst[:0], fnNode, callNode, safeNameIdSupplier, callParamWalkTemplate)
	if hasDuplicateBindingID(dst) {
		f.failDuplicateParamBindings(fnNode, callNode, safeNameIdSupplier)
	}
	return dst
}

func (f *FunctionArgumentInjector) getCallParamBindingsScratch(
	fnNode rhino.Node,
	callNode rhino.Node,
	safeNameIdSupplier func() string,
) []callParamBinding {
	f.paramBindings = f.appendCallParamBindings(f.paramBindings[:0], fnNode, callNode, safeNameIdSupplier, callParamWalkReadOnly)
	if hasDuplicateBindingID(f.paramBindings) {
		f.failDuplicateParamBindings(fnNode, callNode, safeNameIdSupplier)
	}
	return f.paramBindings
}

// FunctionArgumentInjector_GetUniqueAnonymousParameterName: parameter names
// will be name unique when at a later time.
//
// Java: FunctionArgumentInjector.getUniqueAnonymousParameterName(Supplier<String>)
func FunctionArgumentInjector_GetUniqueAnonymousParameterName(safeNameIdSupplier func() string) string {
	return "JSCompiler_inline_anon_param_" + safeNameIdSupplier()
}

// FunctionArgumentInjector_CanNameValueChange is similar to NodeUtil.isLValue
// except that object properties and array member modification aren't
// important ("o" in "o.a = 2" is still "o" after assignment, where in as "o
// = x", "o" is now "x").
//
// This also looks for the redefinition of a name. function (x) {var x;}
//
// n is the NAME node in question.
//
// Java: FunctionArgumentInjector.canNameValueChange(Node)
func FunctionArgumentInjector_CanNameValueChange(n rhino.Node) bool {
	if !NodeUtil_IsLValue(n) {
		return false
	}
	if NodeUtil_IsDeclarationLValue(n) {
		declaring := NodeUtil_GetDeclaringParent(n)
		return !declaring.IsConst() && !declaring.IsLet()
	}
	return true
}

func findThisArgInBindings(bindings []callParamBinding) rhino.Node {
	if len(bindings) > 0 && bindings[0].nameID == calleeThisMarkerID {
		return bindings[0].argNode
	}
	for _, b := range bindings {
		if b.nameID == calleeThisMarkerID {
			return b.argNode
		}
	}
	return rhino.NilNode
}

// findRequiresTempsUpToParameterID returns the strpool ID of the last parameter
// whose call argument requires hoisting into a temporary (or 0 if none do).
//
// If an argument requires a temp (e.g. it is side-effectful, affects mutable
// state or due to other reasons), we need to evaluate it into a temporary at
// the callsite before inlining the function block into that callsite.
// However, that argument's early evaluation could also change the values of
// previous arguments. Hence, in this function we decide to hoist all those
// previous args and evaluate them before this side-effectful argument, even
// if they are not actually affected by that side-effectful argument.
func (f *FunctionArgumentInjector) findRequiresTempsUpToParameterID(
	compiler AbstractCompiler,
	fnNode rhino.Node,
	bindings []callParamBinding,
	modifiedParameters *inlineIDSet,
	convention CodingConvention,
	summary *calleeSummary,
	renamedClone bool,
) uint32 {
	if summary == nil || (!renamedClone && summary.fnNode != fnNode) || !summary.canReuseAfterRename(renamedClone) {
		scanCalleeSummary("", fnNode, &f.scratchSummary)
		summary = &f.scratchSummary
		renamedClone = false
	}

	var requiresTempsUpToID uint32

	argCount := len(bindings)
	isTrivialBody := summary.isTrivialBody
	thisArg := findThisArgInBindings(bindings)
	hasMinimalParameters := (NodeUtil_IsUndefined(thisArg) ||
		(fnNode.IsArrowFunction() && thisArg.IsThis())) &&
		argCount <= 2 // this + one parameter

	useSideEffectSummary := summary.canReuseSideEffects(renamedClone)
	var seSummary calleeSideEffectSummary
	var seSummaryReady bool
	var sideEffectsCollected bool

	collectFallbackSideEffects := func() {
		if sideEffectsCollected {
			return
		}
		seTarget := summary
		if seTarget.fnNode != fnNode {
			scanCalleeSummary("", fnNode, &f.seScratchSummary)
			seTarget = &f.seScratchSummary
		}
		f.collectReferencedAfterSideEffect(seTarget)
		sideEffectsCollected = true
	}

	isReferencedAfterSideEffect := func(paramID uint32, paramIdx int) bool {
		if useSideEffectSummary {
			if !seSummaryReady {
				var ok bool
				seSummary, ok = f.ensureCalleeSideEffects(summary)
				if ok {
					seSummaryReady = true
				} else {
					useSideEffectSummary = false
				}
			}
			if seSummaryReady {
				if paramID == calleeThisMarkerID {
					return seSummary.thisAfterSideEffect
				}
				if paramIdx < 0 || paramIdx >= int(summary.numParams) {
					return false
				}
				return (seSummary.afterSideEffectMask & (uint64(1) << uint(paramIdx))) != 0
			}
		}
		collectFallbackSideEffects()
		return f.referencedAfter.parametersReferenced.contains(paramID)
	}

	hasNamesReferencedBefore := func(paramID uint32, paramIdx int) bool {
		if useSideEffectSummary {
			if !seSummaryReady {
				var ok bool
				seSummary, ok = f.ensureCalleeSideEffects(summary)
				if ok {
					seSummaryReady = true
				} else {
					useSideEffectSummary = false
				}
			}
			if seSummaryReady {
				if paramIdx < 0 || paramIdx >= int(summary.numParams) {
					return false
				}
				return (seSummary.namesBeforeParamMask & (uint64(1) << uint(paramIdx))) != 0
			}
		}
		collectFallbackSideEffects()
		return f.referencedAfter.parametersWithNamesReferencedBeforeParameter.contains(paramID)
	}

	astAnalyzer := compiler.GetAstAnalyzer()
	summaryParams := summary.params()

	// Check for arguments that are evaluated more than once.
	for entryIdx, b := range bindings {
		paramID := b.nameID
		paramIdx := -1
		if renamedClone {
			if entryIdx > 0 && entryIdx-1 < len(summaryParams) {
				paramIdx = entryIdx - 1
			}
		} else if paramID != calleeThisMarkerID {
			paramIdx = summary.findParamByID(paramID)
		}

		cArg := b.argNode
		if !modifiedParameters.isEmpty() && modifiedParameters.contains(paramID) {
			requiresTempsUpToID = paramID
			continue
		}
		// Stores whether this arg need to get hoisted using a temporary
		requiresTemporary := false
		var references int
		if paramIdx >= 0 {
			references = int(summaryParams[paramIdx].refCount)
		}

		argSideEffects := astAnalyzer.MayHaveSideEffects(cArg)
		if !argSideEffects && references == 0 {
			requiresTemporary = false
		} else if isTrivialBody &&
			hasMinimalParameters &&
			references == 1 &&
			// The below line is checking: Can this be affected by side-effects in the function body?
			!(NodeUtil_CanBeSideEffected(cArg) && isReferencedAfterSideEffect(paramID, paramIdx)) &&
			// The below line is checking: Can the function body be affected by the side effects of
			// the argument? If the argument has side effects and the function body contains names
			// that are referenced before the argument, then we need to inline using temporaries. If
			// not, then we can enter this condition's else-if block.
			(!argSideEffects || !hasNamesReferencedBefore(paramID, paramIdx)) {
			// For functions with a trivial body, and where the parameter evaluation order
			// can't change, and there aren't any side-effect before the parameter, we can
			// avoid creating a temporary.
			//
			// This is done to help inline common trivial functions
			// TODO: b/407603216 - return `true` when return expressions don't have
			// references that can be side-effected
			requiresTemporary = false
		} else if astAnalyzer.MayEffectMutableState(cArg) && references > 0 {
			// Note: Mutable arguments should be assigned to temps, as the
			// may be within in a loop:
			//   function x(a) {
			//     for(var i=0; i<0; i++) {
			//       foo(a);
			//     }
			//   x( [] );
			// The parameter in the call to foo should not become "[]".
			requiresTemporary = true
		} else if argSideEffects {
			// even if there are no references, we still need to evaluate the
			// expression if it has side-effects.
			requiresTemporary = true
		} else if NodeUtil_CanBeSideEffected(cArg) && isReferencedAfterSideEffect(paramID, paramIdx) {
			requiresTemporary = true
		} else if references > 1 {
			// Safe is a misnomer, this is a check for "large".
			switch cArg.GetToken() {
			case rhino.Token_NAME:
				name := cArg.GetString()
				// Don't worry about whether this is global or local, just check if it is
				// "exported" in either case.
				requiresTemporary =
					(convention.IsExportedWithLocal(name, true) || convention.IsExportedWithLocal(name, false))
			case rhino.Token_THIS:
				requiresTemporary = false
			case rhino.Token_STRINGLIT:
				requiresTemporary = (jutil.StringLength(cArg.GetString()) >= 2)
			default:
				requiresTemporary = !NodeUtil_IsImmutableValue(cArg)
			}
		}

		if requiresTemporary {
			requiresTempsUpToID = paramID
		}
	}

	return requiresTempsUpToID
}

func (f *FunctionArgumentInjector) hasCallArgumentsNeedingTempsForBindings(
	compiler AbstractCompiler,
	fnNode rhino.Node,
	bindings []callParamBinding,
	modifiedParameters *inlineIDSet,
	convention CodingConvention,
	summary *calleeSummary,
) bool {
	jutil.CheckArgumentMsg1(fnNode.IsFunction(), "%s", fnNode)
	if !modifiedParameters.isEmpty() {
		return true
	}
	if len(bindings) == 0 {
		return false
	}
	requiresTempsUpToID := f.findRequiresTempsUpToParameterID(
		compiler, fnNode, bindings, modifiedParameters, convention, summary, false)
	if requiresTempsUpToID == 0 {
		return false
	}
	thisArg := findThisArgInBindings(bindings)
	for _, b := range bindings {
		if b.nameID == calleeThisMarkerID &&
			(NodeUtil_IsUndefined(thisArg) ||
				(fnNode.IsArrowFunction() && thisArg.IsThis())) {
			continue
		}
		return true
	}
	return false
}

func (f *FunctionArgumentInjector) gatherCallArgumentsNeedingTempsByID(
	compiler AbstractCompiler,
	fnNode rhino.Node,
	bindings []callParamBinding,
	modifiedParameters *inlineIDSet,
	convention CodingConvention,
	summary *calleeSummary,
	renamedClone bool,
	dst *inlineIDSet,
) {
	jutil.CheckArgumentMsg1(fnNode.IsFunction(), "%s", fnNode)
	dst.clear()
	dst.addAll(modifiedParameters)
	if len(bindings) == 0 {
		return
	}
	requiresTempsUpToID := f.findRequiresTempsUpToParameterID(
		compiler, fnNode, bindings, modifiedParameters, convention, summary, renamedClone)
	if requiresTempsUpToID == 0 {
		return
	}
	addCallArgumentsNeedingTempsByID(fnNode, bindings, requiresTempsUpToID, dst)
}

func addCallArgumentsNeedingTempsByID(
	fnNode rhino.Node,
	bindings []callParamBinding,
	requiresTempsUpToID uint32,
	dst *inlineIDSet,
) {
	thisArg := findThisArgInBindings(bindings)
	// mark all names up to requiresTempsUpToID as namesNeedingTemps
	for _, b := range bindings {
		if b.nameID == calleeThisMarkerID &&
			(NodeUtil_IsUndefined(thisArg) ||
				(fnNode.IsArrowFunction() && thisArg.IsThis())) {
			// Skip generating a temporary for an implicit "this" arg.
			continue
		}
		dst.add(b.nameID)
		if b.nameID == requiresTempsUpToID {
			break
		}
	}
}

func (f *FunctionArgumentInjector) collectReferencedAfterSideEffect(summary *calleeSummary) {
	root := ast.MustFunction(summary.fnNode).Body().Node()
	f.scratchLocals.clear()
	gatherLocalNameIDs(root, &f.scratchLocals)

	collector := &f.referencedAfter
	collector.reset(f, summary, &f.scratchLocals)
	collector.walk(root)
	collector.summary = nil
}

func (f *FunctionArgumentInjector) ensureCalleeSideEffects(summary *calleeSummary) (calleeSideEffectSummary, bool) {
	if summary == nil || int(summary.numParams) > 64 {
		return calleeSideEffectSummary{overflow: true}, false
	}
	if summary.sideEffects.computed {
		return summary.sideEffects, true
	}
	if cached, ok := f.sideEffectCache[summary.fnNode]; ok {
		if !f.readOnlySummaries {
			summary.sideEffects = cached
		}
		return cached, true
	}

	f.collectReferencedAfterSideEffect(summary)
	collector := &f.referencedAfter

	se := calleeSideEffectSummary{
		computed:            true,
		thisAfterSideEffect: collector.parametersReferenced.contains(calleeThisMarkerID),
	}
	for i, p := range summary.params() {
		bit := uint64(1) << uint(i)
		if collector.parametersReferenced.contains(p.nameID) {
			se.afterSideEffectMask |= bit
		}
		if collector.parametersWithNamesReferencedBeforeParameter.contains(p.nameID) {
			se.namesBeforeParamMask |= bit
		}
	}
	if f.readOnlySummaries {
		if f.sideEffectCache == nil {
			f.sideEffectCache = make(map[rhino.Node]calleeSideEffectSummary, 16)
		}
		f.sideEffectCache[summary.fnNode] = se
	} else {
		summary.sideEffects = se
	}
	return se, true
}

// FunctionArgumentInjector_ReferencedAfterSideEffect collects parameter
// names referenced after a non-local side-effect.
//
// Assumptions:
//
//   - We assume parameters are not modified in the function body (that is
//     checked separately).
//   - There are no inner functions (also checked separately).
//
// As we are trying to replace parameters with there passed in values we are
// interested in anything that may affect those value. So, ignoring changes
// to local variables, we look for things that may affect anything outside
// the local-state. Once such a side-effect is seen any following reference
// to the function parameters are collected. These will need to be assigned
// to temporaries to prevent changes to their value as would have happened
// during the function call.
//
// To properly handle loop structures all references to the function
// parameters are recorded and the decision to keep or throw away those
// references is deferred until exiting the loop structure.
type FunctionArgumentInjector_ReferencedAfterSideEffect struct {
	outer                *FunctionArgumentInjector
	summary              *calleeSummary
	locals               *inlineIDSet
	sideEffectSeen       bool
	parametersReferenced inlineIDSet

	nameNodeHasBeenSeen                          bool
	parametersWithNamesReferencedBeforeParameter inlineIDSet

	loopsEntered int
}

func (r *FunctionArgumentInjector_ReferencedAfterSideEffect) reset(
	outer *FunctionArgumentInjector,
	summary *calleeSummary,
	locals *inlineIDSet,
) {
	r.outer = outer
	r.summary = summary
	r.locals = locals
	r.sideEffectSeen = false
	r.nameNodeHasBeenSeen = false
	r.loopsEntered = 0
	r.parametersReferenced.clear()
	r.parametersWithNamesReferencedBeforeParameter.clear()
}

func (r *FunctionArgumentInjector_ReferencedAfterSideEffect) paramCount() int {
	return int(r.summary.numParams) + 1
}

func (r *FunctionArgumentInjector_ReferencedAfterSideEffect) hasParameterID(nameID uint32) bool {
	return r.summary.findParamByID(nameID) >= 0
}

func (r *FunctionArgumentInjector_ReferencedAfterSideEffect) walk(node rhino.Node) {
	if r.Apply(node) {
		for c := node.GetFirstChild(); !c.IsNil(); {
			next := c.GetNext()
			r.walk(c)
			c = next
		}
	}
	r.Visit(node)
}

// Java: FunctionArgumentInjector.ReferencedAfterSideEffect.apply(Node)
func (r *FunctionArgumentInjector_ReferencedAfterSideEffect) Apply(node rhino.Node) bool {
	// Keep track of any loop structures entered.
	if NodeUtil_IsLoopStructure(node) {
		r.loopsEntered++
	}

	// If we have found all the parameters, don't bother looking
	// at the children.
	return !(r.sideEffectSeen && r.paramCount() == r.parametersReferenced.size())
}

// Java: FunctionArgumentInjector.ReferencedAfterSideEffect.inLoop()
func (r *FunctionArgumentInjector_ReferencedAfterSideEffect) InLoop() bool {
	return r.loopsEntered != 0
}

// Java: FunctionArgumentInjector.ReferencedAfterSideEffect.visit(Node)
func (r *FunctionArgumentInjector_ReferencedAfterSideEffect) Visit(n rhino.Node) {
	if n.IsName() {
		nameID := n.GetStringID()
		if r.hasParameterID(nameID) && r.nameNodeHasBeenSeen {
			// We have seen a name node before this parameter. If this parameter has side effects, it
			// is not safe to inline.
			r.parametersWithNamesReferencedBeforeParameter.add(nameID)
		}
		r.nameNodeHasBeenSeen = true
	}

	// If we are exiting a loop.
	if NodeUtil_IsLoopStructure(n) {
		r.loopsEntered--
		if !r.InLoop() && !r.sideEffectSeen {
			// Now that the loops has been fully traversed and
			// no side-effects have been seen, throw away
			// the references seen in them.
			r.parametersReferenced.clear()
		}
	}

	if !r.sideEffectSeen {
		// Look for side-effects.
		if r.HasNonLocalSideEffect(n) {
			r.sideEffectSeen = true
		}
	}

	// If traversing the nodes of a loop save any references
	// that are seen.
	if r.InLoop() || r.sideEffectSeen {
		// Record references to parameters.
		if n.IsName() {
			nameID := n.GetStringID()
			if r.hasParameterID(nameID) {
				r.parametersReferenced.add(nameID)
			}
		} else if n.IsThis() {
			r.parametersReferenced.add(calleeThisMarkerID)
		}
	}
}

// HasNonLocalSideEffect returns whether the node may have non-local
// side-effects.
//
// Java: FunctionArgumentInjector.ReferencedAfterSideEffect.hasNonLocalSideEffect(Node)
func (r *FunctionArgumentInjector_ReferencedAfterSideEffect) HasNonLocalSideEffect(n rhino.Node) bool {
	sideEffect := false
	type_ := n.GetToken()
	// Note: Only care about changes to non-local names, specifically
	// ignore VAR declaration assignments.
	if assign, ok := ast.AsAssignment(n); ok {
		// Ignore changes to local names.
		if !r.IsLocalName(assign.Target().Node()) {
			sideEffect = true
		}
	} else if incDec, ok := ast.AsIncDec(n); ok {
		if !r.IsLocalName(incDec.Target().Node()) {
			sideEffect = true
		}
	} else if type_ == rhino.Token_CALL {
		sideEffect = r.outer.astAnalyzer.FunctionCallHasSideEffects(n)
	} else if type_ == rhino.Token_NEW {
		sideEffect = r.outer.astAnalyzer.ConstructorCallHasSideEffects(n)
	} else if type_ == rhino.Token_DELPROP {
		sideEffect = true
	}

	return sideEffect
}

// IsLocalName returns whether node is a reference to locally declared name.
//
// Java: FunctionArgumentInjector.ReferencedAfterSideEffect.isLocalName(Node)
func (r *FunctionArgumentInjector_ReferencedAfterSideEffect) IsLocalName(node rhino.Node) bool {
	if node.IsName() {
		nameID := node.GetStringID()
		return r.hasParameterID(nameID) || r.locals.contains(nameID)
	}
	return false
}
