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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

var inlineFnCallPropID = rhino.RhinoStringPool_InternID("call")

// Java: InlineFunctions.targetSizeAfterInlineExceedsLimit(NodeTraversal, FunctionState)
func inlineTargetSizeExceedsLimit(
	t *NodeTraversal,
	functionState *InlineFunctions_FunctionState,
	maxSizeAfterInlining int,
	lastContainingFn *rhino.Node,
	lastContainingFnSize *int,
	cacheBodySize bool,
) bool {
	containingFunction := t.GetEnclosingFunction()
	// Always inline at the top level,
	// unless maybeAddFunction has marked functionState as not inlinable.
	if containingFunction.IsNil() {
		return false
	}
	inlinedFun := functionState.functionNode()
	if InlineFunctions_IsAlwaysInlinable(inlinedFun) {
		return false
	}

	inlinedFunSize := int(functionState.bodyAstSize)
	if inlinedFunSize == 0 {
		inlinedFunSize = NodeUtil_CountAstSizeUpToLimit(ast.MustFunction(inlinedFun).Body().Node(), maxSizeAfterInlining)
		if cacheBodySize {
			functionState.bodyAstSize = int32(inlinedFunSize)
		}
	}
	var targetFunSize int
	if containingFunction == *lastContainingFn {
		targetFunSize = *lastContainingFnSize
	} else {
		targetFunSize = NodeUtil_CountAstSizeUpToLimit(containingFunction, maxSizeAfterInlining)
		*lastContainingFn = containingFunction
		*lastContainingFnSize = targetFunSize
	}
	return inlinedFunSize+targetFunSize > maxSizeAfterInlining
}

var (
	inlineCallAndNameVisitTokens = makeVisitTokenSet(
		rhino.Token_NAME,
		rhino.Token_CALL,
		rhino.Token_OPTCHAIN_CALL,
	)
	inlineCallOnlyVisitTokens = makeVisitTokenSet(
		rhino.Token_CALL,
		rhino.Token_OPTCHAIN_CALL,
	)
)

// inlineCallVisitor visits candidate call sites (and optionally NAME nodes)
// without interface or base-struct wrapper ceremony.
type inlineCallVisitor struct {
	namedFns   map[uint32]*InlineFunctions_FunctionState
	anonFns    map[rhino.Node]*InlineFunctions_FunctionState
	onCallSite func(t *NodeTraversal, callNode rhino.Node, fs *InlineFunctions_FunctionState)
	onName     func(n, parent rhino.Node)
}

var (
	_ nodeTraversal_visitTokens      = (*inlineCallVisitor)(nil)
	_ nodeTraversal_onlyNameLeaves   = (*inlineCallVisitor)(nil)
	_ nodeTraversal_noShouldTraverse = (*inlineCallVisitor)(nil)
)

func (*inlineCallVisitor) ShouldTraverse(*NodeTraversal, rhino.Node, rhino.Node) bool { return true }
func (*inlineCallVisitor) nodeTraversalNoShouldTraverse()                             {}
func (*inlineCallVisitor) onlyNameLeaves()                                            {}

func (v *inlineCallVisitor) visitTokens() *ast.TokenSet {
	if v.onName != nil {
		return &inlineCallAndNameVisitTokens
	}
	return &inlineCallOnlyVisitTokens
}

func (v *inlineCallVisitor) inlineWalk(n rhino.Node, parent rhino.Node) {
	if n.GetToken() == rhino.Token_CLASS_MEMBERS {
		body := ast.MustClassBody(n)
		for m := range ast.ReadAhead(body.FirstMember()) {
			switch m.Token() {
			case rhino.Token_COMPUTED_PROP:
				v.inlineWalk(m.MustComputedProp().Key().Node(), m.Node())
			case rhino.Token_COMPUTED_FIELD_DEF:
				v.inlineWalk(m.MustComputedFieldDef().Key().Node(), m.Node())
			}
		}
		for m := range ast.ReadAhead(body.FirstMember()) {
			switch m.Token() {
			case rhino.Token_COMPUTED_PROP:
				if value, ok := m.MustComputedProp().Value(); ok {
					v.inlineWalk(value.Node(), m.Node())
				}
			case rhino.Token_COMPUTED_FIELD_DEF:
				if value, ok := m.MustComputedFieldDef().Value(); ok {
					v.inlineWalk(value.Node(), m.Node())
				}
			default:
				v.inlineWalk(m.Node(), n)
			}
		}
		return
	}
	for c := n.GetFirstChild(); !c.IsNil(); {
		next := c.GetNext()
		if !c.GetFirstChild().IsNil() {
			v.inlineWalk(c, n)
		}
		c = next
	}
	switch n.GetToken() {
	case rhino.Token_OPTCHAIN_CALL, rhino.Token_CALL:
		v.visitCallNode(nil, n)
	}
}

func (v *inlineCallVisitor) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_NAME:
		if v.onName != nil {
			v.onName(n, parent)
		}
	case rhino.Token_OPTCHAIN_CALL, rhino.Token_CALL:
		v.visitCallNode(t, n)
	}
}

func (v *inlineCallVisitor) visitCallNode(t *NodeTraversal, n rhino.Node) {
	child := ast.MustNormalOrOptChainCall(n).Callee()
	var functionState *InlineFunctions_FunctionState
	// NOTE: The normalization pass ensures that local names do not collide with global names.
	if callee, ok := child.AsName(); ok {
		functionState = v.namedFns[callee.Node().GetStringID()]
	} else if child.Token() == rhino.Token_FUNCTION {
		functionState = v.anonFns[child.Node()]
	} else if NodeUtil_IsFunctionObjectCall(n) {
		fnIdentifyingNode := child.MustNormalOrOptChainGet().Object()
		if fnName, ok := fnIdentifyingNode.AsName(); ok {
			functionState = v.namedFns[fnName.Node().GetStringID()]
		} else if fnIdentifyingNode.Token() == rhino.Token_FUNCTION {
			functionState = v.anonFns[fnIdentifyingNode.Node()]
		}
	}

	if functionState != nil {
		v.onCallSite(t, n, functionState)
	}
}

// InlineFunctions_IsCandidateUsage returns whether the name is used in a way
// that might be a candidate for inlining.
//
// Java: InlineFunctions.isCandidateUsage(Node)
func InlineFunctions_IsCandidateUsage(name rhino.Node) bool {
	jutil.CheckState(name.IsName())
	switch ast.RoleOf(name) {
	case ast.RoleNameDeclarationDeclarators, ast.RoleFunctionName, ast.RoleFunctionBody:
		// This is a declaration.  Duplicate declarations are handle during
		// function candidate gathering.
		// (Java tests that the parent is a function, so the expression body
		// of an arrow function, () => fn, counts too.)
		return true
	case ast.RoleCallCallee, ast.RoleOptChainCallCallee:
		// This is a normal reference to the function.
		return true
	case ast.RoleGetElemObject:
		// Check for a ".call" to the named function:
		//   CALL
		//     GETPROP/GETELEM
		//       NAME
		//       STRING == "call"
		//     This-Value
		//     Function-parameter-1
		//     ...
		index, ok := ast.MustGetElem(name.GetParent()).Index().AsStringLit()
		if !ok || index.Node().GetStringID() != inlineFnCallPropID {
			return false
		}
	case ast.RoleGetPropObject:
		if name.GetParent().GetStringID() != inlineFnCallPropID {
			return false
		}
	default:
		return false
	}
	// Yep, a ".call".
	return ast.RoleOf(name.GetParent()) == ast.RoleCallCallee
}

// inlineFunctions_isConstAlias reports whether the name n (whose parent is
// parent) is the value of a constant: a declared const (const bar = foo)
// or an assignment marked @const (/** @const */ exports.bar = foo).
func inlineFunctions_isConstAlias(n, parent rhino.Node) bool {
	switch ast.RoleOf(n) {
	case ast.RoleAssignValue:
		info := parent.GetJSDocInfo()
		return info != nil && info.IsConstant()
	case ast.RoleNameInit:
		return parent.GetParent().IsConst()
	}
	return false
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

// inlineRefCollector collects Phase-2 function call references and non-call
// name usages for both sequential and parallel traversals.
type inlineRefCollector struct {
	visitor             inlineCallVisitor
	outer               *InlineFunctions
	injector            *FunctionInjector
	lastContainingFnPtr *rhino.Node
	lastContainingSzPtr *int
	localLastFn         rhino.Node
	localLastFnSize     int
	safeIdsCounter      int
	scratchRef          FunctionInjector_Reference
	localDisallowed     map[*InlineFunctions_FunctionState]struct{}
	events              []ifPhase2Event
	recordEvents        bool
}

func (r *inlineRefCollector) initSerial(f *InlineFunctions) {
	r.outer = f
	r.injector = f.injector
	r.lastContainingFnPtr = &f.lastContainingFn
	r.lastContainingSzPtr = &f.lastContainingFnSize
	r.recordEvents = false
	r.visitor = inlineCallVisitor{
		namedFns:   f.namedFns,
		anonFns:    f.anonFns,
		onCallSite: r.visitCallSite,
		onName:     r.checkNameUsage,
	}
}

func (r *inlineRefCollector) initParallel(f *InlineFunctions, injector *FunctionInjector) {
	r.outer = f
	r.injector = injector
	r.lastContainingFnPtr = &r.localLastFn
	r.lastContainingSzPtr = &r.localLastFnSize
	r.recordEvents = true
	r.visitor = inlineCallVisitor{
		namedFns:   f.namedFns,
		anonFns:    f.anonFns,
		onCallSite: r.visitCallSite,
		onName:     r.checkNameUsage,
	}
}

func (r *inlineRefCollector) resetForScript() {
	r.events = nil
	r.localLastFn = rhino.NilNode
	if r.localDisallowed != nil {
		clear(r.localDisallowed)
	}
}

func (r *inlineRefCollector) isDisallowed(fs *InlineFunctions_FunctionState) bool {
	if !fs.CanInline() {
		return true
	}
	if r.localDisallowed != nil {
		_, ok := r.localDisallowed[fs]
		return ok
	}
	return false
}

func (r *inlineRefCollector) visitCallSite(t *NodeTraversal, callNode rhino.Node, functionState *InlineFunctions_FunctionState) {
	if r.isDisallowed(functionState) {
		return
	}

	chunk := t.GetChunk()
	mode := FunctionInjector_InliningMode_BLOCK
	if functionState.CanInlineDirectly() {
		mode = FunctionInjector_InliningMode_DIRECT
	}
	beforeIds := r.safeIdsCounter
	added, scope, reqDecomp := r.tryMode(t, functionState, callNode, chunk, mode)
	if !added && mode == FunctionInjector_InliningMode_DIRECT {
		// This reference can not be directly inlined, see if
		// block replacement inlining is possible.
		mode = FunctionInjector_InliningMode_BLOCK
		added, scope, reqDecomp = r.tryMode(t, functionState, callNode, chunk, mode)
	}

	if !r.recordEvents {
		if added {
			candidate := r.outer.allocReference(callNode, scope, chunk, mode)
			candidate.SetRequiresDecomposition(reqDecomp)
			functionState.AddReference(candidate)
		} else {
			// Don't try to remove a function if we can't inline all
			// the references.
			functionState.SetRemoveWithContext(
				false,
				r.outer.compiler,
				InlineFunctions_CannotRemoveReason_UNINLINABLE_REFERENCE,
				InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_NO,
				/* contextNode= */ callNode)
		}
		return
	}

	consumed := uint32(r.safeIdsCounter - beforeIds)
	if added {
		r.events = append(r.events, ifPhase2Event{
			kind:                  ifEvAddRef,
			fs:                    functionState,
			node:                  callNode,
			scope:                 scope,
			chunk:                 chunk,
			mode:                  mode,
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

func (r *inlineRefCollector) tryMode(
	t *NodeTraversal,
	functionState *InlineFunctions_FunctionState,
	callNode rhino.Node,
	chunk *JSChunk,
	mode FunctionInjector_InliningMode,
) (bool, *Scope, bool) {
	// If many functions are inlined into the same function F in the same
	// inlining round, then the size of F may exceed the max size.
	// This could be avoided if we bail later, during the inlining phase, eg,
	// in Inline#visitCallSite. However, that is not safe, because at that
	// point expression decomposition has already run, and we want to
	// decompose expressions only for the calls that are actually inlined.
	if r.outer.enforceMaxSizeAfterInlining &&
		inlineTargetSizeExceedsLimit(
			t, functionState, r.outer.maxSizeAfterInlining,
			r.lastContainingFnPtr, r.lastContainingSzPtr, !r.recordEvents) {
		return false, nil, false
	}

	var summary *calleeSummary
	if r.recordEvents {
		summary = functionState.readOnlySummary()
	} else {
		summary = functionState.ensureSummary()
	}
	return inlineCheckCallSiteWithMode(
		t, r.injector, &r.scratchRef, functionState, summary, callNode, chunk, mode)
}

func inlineCheckCallSiteWithMode(
	t *NodeTraversal,
	injector *FunctionInjector,
	scratchRef *FunctionInjector_Reference,
	functionState *InlineFunctions_FunctionState,
	summary *calleeSummary,
	callNode rhino.Node,
	chunk *JSChunk,
	mode FunctionInjector_InliningMode,
) (bool, *Scope, bool) {
	scratchRef.init(callNode, nil, chunk, mode, false)
	scratchRef.t = t
	result := injector.canInlineReference(
		scratchRef,
		functionState.functionNode(),
		functionStateAliasCount(functionState),
		functionState.GetReferencesThis(),
		functionState.HasInnerClosures(),
		summary)
	if result == FunctionInjector_CanInlineResult_NO {
		return false, nil, false
	}
	var scope *Scope
	if result == FunctionInjector_CanInlineResult_AFTER_PREPARATION {
		scope = t.GetScope()
	} else if mode != FunctionInjector_InliningMode_DIRECT && !NodeUtil_IsExprCall(callNode.GetParent()) {
		scope = t.resolveLookupScope(len(t.scopes) - 1)
	}
	return true, scope, result == FunctionInjector_CanInlineResult_AFTER_PREPARATION
}

func (r *inlineRefCollector) checkNameUsage(n rhino.Node, parent rhino.Node) {
	jutil.CheckStateMsg1(n.IsName(), "%s", n)

	id := n.GetStringID()
	if id == rhino.RhinoStringPool_EmptyID {
		return
	}

	// Other refs to a function name remove its candidacy for inlining
	functionState := r.outer.namedFns[id]
	if functionState == nil || InlineFunctions_IsCandidateUsage(n) {
		return
	}

	// If the name is being assigned to or modified as an L-value (e.g. `bar = something;`,
	// destructuring `[bar] = ...`, or `for (bar of ...)`), it cannot be inlined.
	if NodeUtil_IsLValue(n) {
		if !r.recordEvents {
			functionState.DisallowInlining(r.outer.compiler, InlineFunctions_DisallowInliningReason_REASSIGNED)
			return
		}
		if r.localDisallowed == nil {
			r.localDisallowed = make(map[*InlineFunctions_FunctionState]struct{}, 4)
		}
		r.localDisallowed[functionState] = struct{}{}
		r.events = append(r.events, ifPhase2Event{
			kind: ifEvDisallowReassigned,
			fs:   functionState,
		})
	} else if inlineFunctions_isConstAlias(n, parent) {
		// e.g. const bar = foo; exports.bar = foo;
		// We can't see through this reference right now, but we will be able
		// to once `bar` is inlined; so we shouldn't remove `foo` right now but
		// this is probably okay with @requireInlining.
		if !r.recordEvents {
			functionState.SetRemove(
				false,
				r.outer.compiler,
				InlineFunctions_CannotRemoveReason_UNREMOVABLE_REFERENCE,
				InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_NO)
			return
		}
		r.events = append(r.events, ifPhase2Event{
			kind: ifEvUnremovableConst,
			fs:   functionState,
		})
	} else {
		// e.g. var fn = bar; <== we can't inline "bar"
		// As this reference can't be inlined mark the function as
		// unremovable.
		if !r.recordEvents {
			functionState.SetRemoveWithContext(
				false,
				r.outer.compiler,
				InlineFunctions_CannotRemoveReason_UNREMOVABLE_REFERENCE,
				InlineFunctions_ShouldWarnWhenRequireInliningCannotInline_YES,
				/* contextNode= */ parent)
			return
		}
		r.events = append(r.events, ifPhase2Event{
			kind: ifEvUnremovableContext,
			fs:   functionState,
			node: parent,
		})
	}
}

func (f *InlineFunctions) replayPhase2Events(events []ifPhase2Event) {
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

func (f *InlineFunctions) inlineCallSite(t *NodeTraversal, callNode rhino.Node, functionState *InlineFunctions_FunctionState) {
	jutil.CheckState(functionState.HasExistingFunctionDefinition())
	if !functionState.CanInline() {
		return
	}
	ref := functionState.GetReference(callNode)
	// There are two cases ref can be null: if the call site was introduced
	// because it was part of a function that was inlined during this pass
	// or if the call site was trimmed from the list of references because
	// the function couldn't be inlined at this location.
	if ref == nil {
		return
	}
	fnName := functionState.functionName()
	if fnName == "" && functionState.target.kind == fnKindExpr && ref.mode == FunctionInjector_InliningMode_BLOCK {
		fnName = f.functionStateName(functionState)
	}
	fnNode := functionState.GetSafeFnNode()
	newExpr := f.injector.inlineWithSummary(ref, fnName, fnNode, functionState.ensureSummary())
	if newExpr != ref.callNode {
		f.invalidateEnclosingSummaries(newExpr)
		if t != nil {
			t.GetCompiler().ReportChangeToEnclosingScope(newExpr)
		} else {
			f.injector.compiler.ReportChangeToEnclosingScope(newExpr)
		}
	}
	// Keep track of references that have been inlined so that
	// we can verify that none have been missed.
	ref.inlined = true
}

func (f *InlineFunctions) invalidateEnclosingSummaries(node rhino.Node) {
	if len(f.fns) == 0 {
		return
	}
	hasAnon := len(f.anonFns) > 0
	for anc := node.GetParent(); !anc.IsNil(); anc = anc.GetParent() {
		if !anc.IsFunction() {
			continue
		}
		if name, ok := ast.MustFunction(anc).Name(); ok {
			if fs := f.namedFns[name.Node().GetStringID()]; fs != nil && fs.GetSafeFnNode() == anc {
				fs.summaryValid = false
			}
		}
		if p := anc.GetParent(); p.IsName() {
			if fs := f.namedFns[p.GetStringID()]; fs != nil && fs.GetSafeFnNode() == anc {
				fs.summaryValid = false
			}
		}
		if hasAnon {
			if fs := f.anonFns[anc]; fs != nil && fs.GetSafeFnNode() == anc {
				fs.summaryValid = false
			}
		}
	}
}
