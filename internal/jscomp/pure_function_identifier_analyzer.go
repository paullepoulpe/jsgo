// Port of com/google/javascript/jscomp/PureFunctionIdentifier.java @ closure-compiler v20261021

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
	"slices"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// pfiPendingCheck records a local variable passed as receiver (argIndex == -1)
// or argument (argIndex >= 0) to a call inside the current function before the
// function's skiplistedVars / skiplistedNameNodes are finalized at ExitScope.
type pfiPendingCheck struct {
	v         *Var
	edgeStart int32
	edgeEnd   int32
	argIndex  int16
}

// pfiFunctionStackEntry holds transient per-function analysis state, reused
// in-place across functions as scopes are entered and exited.
type pfiFunctionStackEntry struct {
	root                rhino.Node
	isArtificiallyPure  bool
	isArrow             bool
	catchDepth          int
	skiplistedVars      []*Var
	skiplistedNameNodes []rhino.Node
	taintedVars         []*Var
	pendingChecks       []pfiPendingCheck
	activeBuf           [4]pfiSummaryID
}

func (e *pfiFunctionStackEntry) addSkiplistedVar(v *Var) {
	if slices.Contains(e.skiplistedVars, v) {
		return
	}
	e.skiplistedVars = append(e.skiplistedVars, v)
}

func (e *pfiFunctionStackEntry) isSkiplistedVar(v *Var) bool {
	if slices.Contains(e.skiplistedVars, v) {
		return true
	}
	for _, node := range e.skiplistedNameNodes {
		if node == v.GetNameNode() {
			return true
		}
		if node.GetStringID() == v.nameID && node.GetParent().IsVar() && v.scope.IsHoistScope() {
			return true
		}
	}
	return false
}

func (e *pfiFunctionStackEntry) addTaintedVar(v *Var) {
	if slices.Contains(e.taintedVars, v) {
		return
	}
	e.taintedVars = append(e.taintedVars, v)
}

func (e *pfiFunctionStackEntry) isTaintedVar(v *Var) bool {
	return slices.Contains(e.taintedVars, v)
}

func (e *pfiFunctionStackEntry) isLocalVar(v *Var, scope *Scope) bool {
	return PureFunctionIdentifier_FunctionBodyAnalyzer_IsVarDeclaredInSameContainerScope(v, scope) &&
		!e.isSkiplistedVar(v) &&
		!v.IsParam() &&
		!v.IsCatch()
}

// pfiFunctionBodyAnalyzer inspects function bodies for side effects, records
// call sites once, and resolves local variable escaping at each function's ExitScope.
type pfiFunctionBodyAnalyzer struct {
	outer    *PureFunctionIdentifier
	stack    []pfiFunctionStackEntry
	stackTop int
}

func newPFIFunctionBodyAnalyzer(outer *PureFunctionIdentifier) *pfiFunctionBodyAnalyzer {
	b := &pfiFunctionBodyAnalyzer{
		outer:    outer,
		stack:    make([]pfiFunctionStackEntry, 1, 8),
		stackTop: 0,
	}
	return b
}

func (b *pfiFunctionBodyAnalyzer) currentFunction() *pfiFunctionStackEntry {
	if b.stackTop < 0 {
		panic(jutil.NewNoSuchElementException())
	}
	return &b.stack[b.stackTop]
}

func (b *pfiFunctionBodyAnalyzer) activeSummaries(e *pfiFunctionStackEntry) []pfiSummaryID {
	p := b.outer
	single, multi := p.getFunctionSummaries(e.root)
	if single != 0 {
		if p.summaries[single].bitmask&pfiMutatesGlobalState != 0 {
			return nil
		}
		e.activeBuf[0] = single
		return e.activeBuf[:1]
	}
	out := e.activeBuf[:0]
	for _, sid := range multi {
		if p.summaries[sid].bitmask&pfiMutatesGlobalState == 0 {
			out = append(out, sid)
		}
	}
	return out
}

func (b *pfiFunctionBodyAnalyzer) setMutatesGlobalState(active []pfiSummaryID, reason rhino.Node) {
	p := b.outer
	for _, sid := range active {
		p.setSummaryMutatesGlobalState(sid, &p.summaries[sid], reason)
	}
}

func (b *pfiFunctionBodyAnalyzer) setMutatesThis(active []pfiSummaryID, reason rhino.Node) {
	p := b.outer
	for _, sid := range active {
		p.setSummaryMutatesThis(sid, &p.summaries[sid], reason)
	}
}

func (b *pfiFunctionBodyAnalyzer) recordThrowsBasedOnContext(e *pfiFunctionStackEntry, active []pfiSummaryID) {
	if e.catchDepth == 0 {
		p := b.outer
		for _, sid := range active {
			p.setSummaryThrows(sid, &p.summaries[sid], e.root)
		}
	}
}

func (b *pfiFunctionBodyAnalyzer) checkIteratesImpureIterable(node rhino.Node, active []pfiSummaryID) {
	if NodeUtil_IteratesImpureIterable(node) {
		b.setMutatesGlobalState(active, node)
	}
}

// Java: PureFunctionIdentifier.FunctionBodyAnalyzer.shouldTraverse(NodeTraversal, Node, Node)
func (b *pfiFunctionBodyAnalyzer) ShouldTraverse(traversal *NodeTraversal, node rhino.Node, parent rhino.Node) bool {
	b.addToCatchDepthIfTryBlock(node, 1)
	return true
}

// Java: PureFunctionIdentifier.FunctionBodyAnalyzer.visit(NodeTraversal, Node, Node)
func (b *pfiFunctionBodyAnalyzer) Visit(traversal *NodeTraversal, node rhino.Node, parent rhino.Node) {
	b.addToCatchDepthIfTryBlock(node, -1)

	if NodeUtil_IsInvocation(node) {
		site := b.outer.recordCallSite(node)
		e := b.currentFunction()
		if e.root.IsNil() || e.isArtificiallyPure || site.flags&pfiCallSkipAnalyzer != 0 {
			return
		}
		active := b.activeSummaries(e)
		if len(active) == 0 {
			return
		}
		b.visitResolvedCall(e, active, traversal, node, site)
		return
	}

	e := b.currentFunction()
	if e.root.IsNil() || e.isArtificiallyPure {
		return
	}
	active := b.activeSummaries(e)
	if len(active) == 0 {
		return
	}
	b.updateSideEffectsForNonCallNode(e, active, traversal, node)
}

func (b *pfiFunctionBodyAnalyzer) updateSideEffectsForNonCallNode(
	e *pfiFunctionStackEntry, active []pfiSummaryID, traversal *NodeTraversal, node rhino.Node) {
	switch node.GetToken() {
	case rhino.Token_ASSIGN:
		lhs := ast.MustAssign(node).Target()
		if _, isPat := lhs.AsPattern(); isPat {
			NodeUtil_VisitLhsNodesInNode(
				node,
				func(lhsNode rhino.Node) {
					b.visitLhsNode(e, active, traversal, lhsNode, PureFunctionIdentifier_RHS_IS_NEVER_LOCAL)
				})
		} else {
			b.visitLhsNode(e, active, traversal, lhs.Node(), PureFunctionIdentifier_FIND_RHS_AND_CHECK_FOR_LOCAL_VALUE)
		}

	case rhino.Token_INC, rhino.Token_DEC:
		b.visitLhsNode(
			e,
			active,
			traversal,
			ast.MustIncDec(node).Target().Node(),
			PureFunctionIdentifier_RHS_IS_ALWAYS_LOCAL)

	case rhino.Token_DELPROP:
		b.visitLhsNode(
			e,
			active,
			traversal,
			ast.MustUnary(node).Operand().Node(),
			PureFunctionIdentifier_RHS_IS_ALWAYS_LOCAL)

	case rhino.Token_FOR_AWAIT_OF:
		b.recordThrowsBasedOnContext(e, active)
		fallthrough
	case rhino.Token_FOR_OF:
		NodeUtil_VisitLhsNodesInNode(
			node,
			func(lhsNode rhino.Node) {
				b.visitLhsNode(
					e,
					active,
					traversal,
					lhsNode,
					PureFunctionIdentifier_RHS_IS_NEVER_LOCAL)
			})
		b.checkIteratesImpureIterable(node, active)

	case rhino.Token_FOR_IN:
		NodeUtil_VisitLhsNodesInNode(
			node,
			func(lhsNode rhino.Node) {
				b.visitLhsNode(
					e,
					active,
					traversal,
					lhsNode,
					PureFunctionIdentifier_RHS_IS_ALWAYS_LOCAL)
			})

	case rhino.Token_DESTRUCTURING_LHS:
		if NodeUtil_IsAnyFor(node.GetParent()) {
			break
		}
		NodeUtil_VisitLhsNodesInNode(
			node.GetParent(),
			func(lhsNode rhino.Node) {
				b.visitLhsNode(e, active, traversal, lhsNode, PureFunctionIdentifier_RHS_IS_NEVER_LOCAL)
			})

	case rhino.Token_NAME:
		if _, ok := ast.AsNameDeclaration(node.GetParent()); ok {
			if value, hasInit := ast.MustName(node).Init(); hasInit && !NodeUtil_EvaluatesToLocalValue(value.Node()) {
				e.skiplistedNameNodes = append(e.skiplistedNameNodes, node)
			}
		}

	case rhino.Token_THROW:
		b.recordThrowsBasedOnContext(e, active)

	case rhino.Token_YIELD:
		b.checkIteratesImpureIterable(node, active)
		b.recordThrowsBasedOnContext(e, active)

	case rhino.Token_AWAIT:
		b.recordThrowsBasedOnContext(e, active)

	case rhino.Token_OBJECT_REST, rhino.Token_OBJECT_SPREAD:
		if !b.outer.assumeGettersArePure {
			b.setMutatesGlobalState(active, node)
		}

	case rhino.Token_ITER_REST, rhino.Token_ITER_SPREAD:
		b.checkIteratesImpureIterable(node, active)

	case rhino.Token_STRING_KEY:
		if _, ok := ast.AsObjectPattern(node.GetParent()); ok {
			if b.outer.getPropertyKindByID(node.GetStringID()).HasGetter() {
				b.setMutatesGlobalState(active, node)
			}
		}

	case rhino.Token_OPTCHAIN_GETPROP, rhino.Token_GETPROP:
		if b.outer.getPropertyKindByID(node.GetStringID()).HasGetterOrSetter() {
			b.setMutatesGlobalState(active, node)
		}

	case rhino.Token_DYNAMIC_IMPORT:
		b.setMutatesGlobalState(active, node)

	default:
		if NodeUtil_IsCompoundAssignmentOp(node) {
			b.visitLhsNode(
				e,
				active,
				traversal,
				ast.MustAssignOp(node).Target().Node(),
				PureFunctionIdentifier_RHS_IS_ALWAYS_LOCAL)
			break
		}

		if b.outer.astAnalyzer.NodeTypeMayHaveSideEffects(node) {
			panic(jutil.NewIllegalArgumentExceptionMsg("Unhandled side effect node type " + node.String()))
		}
	}
}

// Java: PureFunctionIdentifier.FunctionBodyAnalyzer.enterScope(NodeTraversal)
func (b *pfiFunctionBodyAnalyzer) EnterScope(t *NodeTraversal) {
	function := t.GetScopeRoot()
	if !function.IsFunction() {
		if b.currentFunction().root.IsNil() {
			// Variables declared outside all functions can never share a container
			// scope with any scope inside a function. Pre-populating an empty
			// Scope avoids scanning the global/module AST in SyntacticScopeCreator.
			if len(t.scopes) == 1 {
				t.scopes[0].untypedScope = Scope_CreateGlobalScope(function)
			} else if parentScope := t.scopes[len(t.scopes)-2].untypedScope; parentScope != nil {
				t.scopes[len(t.scopes)-1].untypedScope = Scope_CreateChildScope(parentScope, function)
			}
		}
		return
	}

	jsdoc := NodeUtil_GetBestJSDocInfo(function)
	isArtificiallyPure := jsdoc != nil && jsdoc.IsNoSideEffects()

	b.stackTop++
	if b.stackTop == len(b.stack) {
		b.stack = append(b.stack, pfiFunctionStackEntry{})
	}
	e := &b.stack[b.stackTop]
	e.root = function
	e.isArtificiallyPure = isArtificiallyPure
	e.isArrow = function.IsArrowFunction()
	e.catchDepth = 0
	e.skiplistedVars = e.skiplistedVars[:0]
	e.skiplistedNameNodes = e.skiplistedNameNodes[:0]
	e.taintedVars = e.taintedVars[:0]
	e.pendingChecks = e.pendingChecks[:0]

	single, multi := b.outer.getFunctionSummaries(function)
	if single == 0 && len(multi) == 0 {
		sid := b.outer.allocSummary(0, pfiNameAnonymous)
		b.outer.putSummaryForFunctionNode(function, sid)
	}
}

// Java: PureFunctionIdentifier.FunctionBodyAnalyzer.exitScope(NodeTraversal)
func (b *pfiFunctionBodyAnalyzer) ExitScope(t *NodeTraversal) {
	e := b.currentFunction()
	if e.root.IsNil() {
		return
	}
	exitingFunction := t.GetScopeRoot() == e.root

	if len(e.taintedVars) != 0 {
		if scope := t.scopes[len(t.scopes)-1].untypedScope; scope != nil {
			functionScope := scope.GetClosestContainerScope()
			if functionScope.IsFunctionScope() {
				jutil.CheckStateMsg1(
					functionScope.GetRootNode() == e.root, "%s", functionScope.GetRootNode())
				b.resolveTaintedVarsInScope(e, scope)
			}
		}
	}

	if !exitingFunction {
		return
	}

	// All AST nodes, inner blocks, and nested closures of e.root have been visited:
	// skiplistedVars and skiplistedNameNodes are now 100% finalized. Resolve any
	// pending local-variable checks into compact bitmasks on pfiEdge and drop Scope/Var!
	if len(e.pendingChecks) != 0 && (len(e.skiplistedVars) != 0 || len(e.skiplistedNameNodes) != 0) {
		edges := b.outer.edges
		for i := range e.pendingChecks {
			pc := e.pendingChecks[i]
			if !e.isSkiplistedVar(pc.v) {
				continue
			}
			for edgeIdx := pc.edgeStart; edgeIdx < pc.edgeEnd; edgeIdx++ {
				if pc.argIndex < 0 {
					edges[edgeIdx].flags |= pfiEdgeThisGlobal
				} else if pc.argIndex >= 64 {
					edges[edgeIdx].flags |= pfiEdgeAlwaysNonLocalArgs
				} else {
					edges[edgeIdx].nonLocalArgMask |= uint64(1) << pc.argIndex
				}
			}
		}
	}

	clear(e.skiplistedVars)
	e.skiplistedVars = e.skiplistedVars[:0]
	e.skiplistedNameNodes = e.skiplistedNameNodes[:0]
	clear(e.taintedVars)
	e.taintedVars = e.taintedVars[:0]
	clear(e.pendingChecks)
	e.pendingChecks = e.pendingChecks[:0]
	e.root = rhino.NilNode
	b.stackTop--
}

func (b *pfiFunctionBodyAnalyzer) resolveTaintedVarsInScope(e *pfiFunctionStackEntry, scope *Scope) {
	p := b.outer
	active := b.activeSummaries(e)
	for _, sid := range active {
		s := &p.summaries[sid]
		if s.bitmask&pfiMutatesGlobalState != 0 {
			continue
		}
		for _, v := range scope.varList {
			if v == nil || !e.isTaintedVar(v) {
				continue
			}
			isFromDestructuring := NodeUtil_IsLhsByDestructuring(v.GetNameNode())
			if v.IsParam() && !isFromDestructuring && !e.isSkiplistedVar(v) {
				p.setSummaryMutatesArgumentAt(sid, s, PureFunctionIdentifier_GetParamIndex(v), v.GetNode())
				continue
			}
			if !e.isLocalVar(v, scope) {
				p.setSummaryMutatesGlobalState(sid, s, v.GetNode())
				break
			}
		}
	}
}

// Java: PureFunctionIdentifier.FunctionBodyAnalyzer.isVarDeclaredInSameContainerScope(Var, Scope)
func PureFunctionIdentifier_FunctionBodyAnalyzer_IsVarDeclaredInSameContainerScope(v *Var, scope *Scope) bool {
	if v == nil {
		return false
	}
	varScope, _ := v.GetScope().(*Scope)
	return varScope.HasSameContainerScope(scope)
}

func (b *pfiFunctionBodyAnalyzer) visitLhsNode(
	e *pfiFunctionStackEntry,
	active []pfiSummaryID,
	traversal *NodeTraversal,
	lhs rhino.Node,
	hasLocalRhs func(rhino.Node) bool) {
	if get, ok := ast.AsNormalOrOptChainGet(lhs); ok {
		object := get.Object()
		if _, isThis := object.AsThis(); isThis {
			if e.isArrow {
				b.setMutatesGlobalState(active, lhs)
			} else {
				b.setMutatesThis(active, lhs)
			}
		} else if objName, isName := object.AsName(); isName {
			scope := traversal.GetScope()
			var_ := scope.GetVarByID(objName.Node().GetStringID())
			if PureFunctionIdentifier_FunctionBodyAnalyzer_IsVarDeclaredInSameContainerScope(var_, scope) {
				e.addTaintedVar(var_)
			} else {
				b.setMutatesGlobalState(active, lhs)
			}
		} else {
			b.setMutatesGlobalState(active, lhs)
		}
	} else {
		jutil.CheckStateMsg1(lhs.IsName(), "%s", lhs)
		scope := traversal.GetScope()
		var_ := scope.GetVarByID(lhs.GetStringID())
		if PureFunctionIdentifier_FunctionBodyAnalyzer_IsVarDeclaredInSameContainerScope(var_, scope) {
			if !hasLocalRhs(lhs) {
				e.addSkiplistedVar(var_)
			}
		} else {
			b.setMutatesGlobalState(active, lhs)
		}
	}
}

// classifyLocalExpr checks whether expr passed to a callee holds a value local
// to the caller e. Returns (isDefinitelyLocal, pendingVar):
//   - (true, nil): definitely local
//   - (false, nil): definitely non-local
//   - (true, v): local iff !e.isSkiplistedVar(v) at e's ExitScope
func (b *pfiFunctionBodyAnalyzer) classifyLocalExpr(
	e *pfiFunctionStackEntry, traversal *NodeTraversal, scope **Scope, expr rhino.Node) (bool, *Var) {
	if NodeUtil_EvaluatesToLocalValue(expr) {
		return true, nil
	}
	if !expr.IsName() {
		return false, nil
	}
	s := *scope
	if s == nil {
		s = traversal.GetScope()
		*scope = s
	}
	v := s.GetVarByID(expr.GetStringID())
	if v == nil || !PureFunctionIdentifier_FunctionBodyAnalyzer_IsVarDeclaredInSameContainerScope(v, s) || v.IsParam() || v.IsCatch() {
		return false, nil
	}
	if e.isSkiplistedVar(v) {
		return false, nil
	}
	return true, v
}

func (b *pfiFunctionBodyAnalyzer) visitResolvedCall(
	e *pfiFunctionStackEntry,
	active []pfiSummaryID,
	traversal *NodeTraversal,
	invocation rhino.Node,
	site pfiCallSite) {
	if site.calleeCount == 0 {
		b.setMutatesGlobalState(active, rhino.NilNode)
		return
	}

	var edgeFlags uint8
	if e.catchDepth == 0 {
		edgeFlags |= pfiEdgePropagateThrows
	}

	isViaCallOrApply := site.flags&pfiCallViaCallOrApply != 0
	isCallOrTagged := site.flags&pfiCallIsCallOrTagged != 0
	isCacheCall := site.flags&pfiCallIsCacheCall != 0

	var lazyScope *Scope
	var pendingReceiverVar *Var

	if !invocation.IsNew() {
		if pfiCalleeAndCallerShareThis(invocation, isCallOrTagged, isViaCallOrApply) && !e.isArrow {
			edgeFlags |= pfiEdgeThisCaller
		} else {
			callee := pureFunctionIdentifier_invocationCallee(invocation)
			if !isViaCallOrApply {
				if getProp, isGetProp := callee.AsNormalOrOptChainGetProp(); isGetProp {
					isLocal, pendingVar := b.classifyLocalExpr(e, traversal, &lazyScope, getProp.Object().Node())
					if !isLocal {
						edgeFlags |= pfiEdgeThisGlobal
					} else {
						pendingReceiverVar = pendingVar
					}
				} else {
					edgeFlags |= pfiEdgeThisGlobal
				}
			} else {
				edgeFlags |= pfiEdgeThisGlobal
			}
		}
	}

	var nonLocalArgMask uint64
	prevPendingLen := len(e.pendingChecks)

	argumentsMatchParameters := !invocation.IsTaggedTemplateLit() && !isViaCallOrApply && !isCacheCall
	if !argumentsMatchParameters {
		if !NodeUtil_AllArgsUnescapedLocal(invocation) {
			edgeFlags |= pfiEdgeAlwaysNonLocalArgs
		}
	} else {
		index := 0
		for arg := range ast.MustCallOrNew(invocation).Args() {
			argNode := arg.Node()
			if argNode.IsSpread() {
				edgeFlags |= pfiEdgeAlwaysNonLocalArgs
				e.pendingChecks = e.pendingChecks[:prevPendingLen]
				break
			}
			isLocal, pendingVar := b.classifyLocalExpr(e, traversal, &lazyScope, argNode)
			if !isLocal {
				if index >= 64 {
					edgeFlags |= pfiEdgeAlwaysNonLocalArgs
					e.pendingChecks = e.pendingChecks[:prevPendingLen]
					break
				}
				nonLocalArgMask |= uint64(1) << index
			} else if pendingVar != nil {
				e.pendingChecks = append(e.pendingChecks, pfiPendingCheck{
					v:        pendingVar,
					argIndex: int16(index),
				})
			}
			index++
		}
	}

	p := b.outer
	edgeStart := int32(len(p.edges))
	if site.calleeCount == 1 {
		calleeID := pfiSummaryID(site.firstCallee)
		for _, callerID := range active {
			p.addEdge(calleeID, callerID, e.root, edgeFlags, nonLocalArgMask)
		}
	} else {
		callees := p.multiCallees[site.firstCallee : site.firstCallee+site.calleeCount]
		for _, calleeID := range callees {
			for _, callerID := range active {
				p.addEdge(calleeID, callerID, e.root, edgeFlags, nonLocalArgMask)
			}
		}
	}
	edgeEnd := int32(len(p.edges))

	for i := prevPendingLen; i < len(e.pendingChecks); i++ {
		e.pendingChecks[i].edgeStart = edgeStart
		e.pendingChecks[i].edgeEnd = edgeEnd
	}
	if pendingReceiverVar != nil {
		e.pendingChecks = append(e.pendingChecks, pfiPendingCheck{
			v:         pendingReceiverVar,
			edgeStart: edgeStart,
			edgeEnd:   edgeEnd,
			argIndex:  -1,
		})
	}
}

func pfiCalleeAndCallerShareThis(invocation rhino.Node, isCallOrTagged bool, isViaCallOrApply bool) bool {
	if !isCallOrTagged {
		return false
	}
	callee := pureFunctionIdentifier_invocationCallee(invocation)
	calleeNode := callee.Node()
	if calleeNode.IsSuper() {
		return true
	}

	var thisArg rhino.Node
	if isViaCallOrApply {
		if firstArg, ok := ast.MustCallOrNew(invocation).FirstArg(); ok {
			thisArg = firstArg.Node()
		}
	} else if getProp, ok := callee.AsNormalOrOptChainGetProp(); ok {
		thisArg = getProp.Object().Node()
	}

	if thisArg.IsNil() {
		return false
	}
	return thisArg.IsThis() || thisArg.IsSuper()
}

func (b *pfiFunctionBodyAnalyzer) addToCatchDepthIfTryBlock(n rhino.Node, delta int) {
	parent := n.GetParent()
	if ast.RoleIn(n, parent) != ast.RoleTryBody {
		return
	}
	if _, hasCatch := ast.MustTry(parent).Catch(); !hasCatch {
		return
	}
	b.currentFunction().catchDepth += delta
}

// Java: PureFunctionIdentifier.getParamIndex(Var)
func PureFunctionIdentifier_GetParamIndex(param *Var) int {
	var node rhino.Node
	if param.IsDefaultParam() {
		node = param.GetParentNode()
	} else {
		node = param.GetNameNode()
	}
	_, isParamList := ast.AsParamList(node.GetParent())
	jutil.CheckStateMsg1(isParamList, "%s", node)
	return node.GetParent().GetIndexOfChild(node)
}

// Java: PureFunctionIdentifier.isInvocationViaCallOrApply(Node)
func PureFunctionIdentifier_IsInvocationViaCallOrApply(callSite rhino.Node) bool {
	tok := callSite.GetToken()
	if tok != rhino.Token_CALL && tok != rhino.Token_OPTCHAIN_CALL {
		return false
	}
	calleeNode := callSite.GetFirstChild()
	calleeTok := calleeNode.GetToken()
	var propID uint32
	switch calleeTok {
	case rhino.Token_GETPROP, rhino.Token_OPTCHAIN_GETPROP:
		propID = calleeNode.GetStringID()
	case rhino.Token_GETELEM, rhino.Token_OPTCHAIN_GETELEM:
		key := calleeNode.GetSecondChild()
		if !key.IsStringLit() {
			return false
		}
		propID = key.GetStringID()
	default:
		return false
	}
	if propID != pfiCallStringID && propID != pfiApplyStringID {
		return false
	}
	recvTok := calleeNode.GetFirstChild().GetToken()
	return recvTok == rhino.Token_NAME || recvTok == rhino.Token_GETPROP || recvTok == rhino.Token_OPTCHAIN_GETPROP
}

// Java: PureFunctionIdentifier.isCallOrTaggedTemplateLit(Node)
func PureFunctionIdentifier_IsCallOrTaggedTemplateLit(invocation rhino.Node) bool {
	tok := invocation.GetToken()
	return tok == rhino.Token_CALL || tok == rhino.Token_OPTCHAIN_CALL || tok == rhino.Token_TAGGED_TEMPLATELIT
}

func (p *PureFunctionIdentifier) getPropertyKindByID(id uint32) AccessorSummary_PropertyAccessKind {
	if p.assumeGettersArePure {
		return AccessorSummary_PropertyAccessKind_NORMAL
	}
	return p.accessorSummary.GetKindByID(id)
}

// Java: PureFunctionIdentifier.getPropertyKind(String)
func (p *PureFunctionIdentifier) GetPropertyKind(name string) AccessorSummary_PropertyAccessKind {
	if p.assumeGettersArePure {
		return AccessorSummary_PropertyAccessKind_NORMAL
	}
	return p.compiler.GetAccessorSummary().GetKind(name)
}
