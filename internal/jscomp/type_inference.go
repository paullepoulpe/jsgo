// Port of com/google/javascript/jscomp/TypeInference.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/graph"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type tiNodeState struct {
	in        *LinkedFlowScope
	out       *LinkedFlowScope
	stepCount int
}

// TypeInference is type inference within a script node or a function body,
// using a monomorphic forward branched data-flow fixed-point worklist.
type TypeInference struct {
	cfg        *ControlFlowGraph[rhino.Node]
	worklist   forwardCfgWorklist
	nodeStates []tiNodeState
	edgeStates []*LinkedFlowScope

	compiler           AbstractCompiler
	registry           *rhino.JSTypeRegistry
	reverseInterpreter ReverseAbstractInterpreter
	bottomScope        *LinkedFlowScope
	// global scope, function scope, or static block scope
	containerScope          *TypedScope
	scopeCreator            *TypedScopeCreator
	assertionFunctionLookup *CodingConvention_AssertionFunctionLookup
	moduleImportResolver    *ModuleImportResolver
	// A record is pushed onto this stack during the traversal of each optional chain.
	// Inference of the type at the end of the optional chain requires scope information
	// from traversing the start of the chain. This stack serves as a communication channel
	// for that information.
	optChainStack []TypeInference_OptChainInfo

	// Scopes that have had their unbound untyped vars inferred as undefined.
	inferredUnboundVarScope *TypedScope
	inferredUnboundVarBuf   [8]*TypedScope
	inferredUnboundVarList  []*TypedScope
	inferredUnboundVars     map[*TypedScope]struct{}

	// For convenience
	unknownType             rhino.ObjectType
	numberAdditionSupertype rhino.JSType
	arrayOfUnknownType      rhino.JSType

	// Reusable per-analysis state to avoid allocations on every scope / CFG node.
	flowScopeArena           linkedFlowScopeArena
	bottomTypedScope         TypedScope
	bottomImplicitVars       [4]*TypedVar
	bottomFlowScope          LinkedFlowScope
	flowJoiner               LinkedFlowScope_FlowScopeJoinOp
	branchSource             rhino.Node
	branchOutput             *LinkedFlowScope
	branchCondition          rhino.Node
	branchConditionFlowScope *LinkedFlowScope
	branchConditionOutcomes  TypeInference_BooleanOutcomePair
	hasBranchOutcomes        bool

	lastScopeRoot           rhino.Node
	lastModule              *Module
	lastSyntacticBlockScope *TypedScope
}

var TypeInference_REASSIGN_CLASS_PROTOTYPE = DiagnosticType_Error(
	"JSC_REASSIGN_CLASS_PROTOTYPE", "Reassigning a class prototype is not allowed")

const TypeInference_GOOG_REQUIREDYNAMIC_NAME = "goog.requireDynamic"

// Java: TypeInference(AbstractCompiler, ControlFlowGraph<Node>, ReverseAbstractInterpreter, TypedScope, TypedScopeCreator, AssertionFunctionLookup)
func NewTypeInference(
	compiler AbstractCompiler,
	cfg *ControlFlowGraph[rhino.Node],
	reverseInterpreter ReverseAbstractInterpreter,
	syntacticScope *TypedScope,
	scopeCreator *TypedScopeCreator,
	assertionFunctionLookup *CodingConvention_AssertionFunctionLookup) *TypeInference {
	ti := &TypeInference{}
	ti.initTypeInference(compiler, cfg, reverseInterpreter, syntacticScope, scopeCreator, assertionFunctionLookup)
	return ti
}

func (ti *TypeInference) ReleaseLargeBuffers() {
	ti.flowScopeArena.trimLargeChunks()
	ti.worklist.releaseLarge()
	if cap(ti.nodeStates) > dfaReuseMaxNodes {
		ti.nodeStates = nil
	}
	if cap(ti.edgeStates) > dfaReuseMaxNodes*2 {
		ti.edgeStates = nil
	}
}

func (ti *TypeInference) initTypeInference(
	compiler AbstractCompiler,
	cfg *ControlFlowGraph[rhino.Node],
	reverseInterpreter ReverseAbstractInterpreter,
	syntacticScope *TypedScope,
	scopeCreator *TypedScopeCreator,
	assertionFunctionLookup *CodingConvention_AssertionFunctionLookup) {
	ti.cfg = cfg
	ti.flowScopeArena.reset()
	// Field initializers.
	ti.lastScopeRoot = rhino.NilNode
	ti.lastModule = nil
	ti.lastSyntacticBlockScope = nil
	ti.inferredUnboundVarScope = nil
	ti.inferredUnboundVarList = ti.inferredUnboundVarBuf[:0]
	if n := len(ti.inferredUnboundVars); n > 0 {
		if n <= 64 {
			// order-independent: clears the map
			for k := range ti.inferredUnboundVars {
				delete(ti.inferredUnboundVars, k)
			}
		} else {
			clear(ti.inferredUnboundVars)
		}
	}
	ti.optChainStack = ti.optChainStack[:0]
	ti.hasBranchOutcomes = false
	ti.branchConditionOutcomes = TypeInference_BooleanOutcomePair{}
	// Constructor body.
	jutil.CheckArgumentMsg2(
		syntacticScope.GetRootNode() == cfg.GetEntry().GetValue(),
		"Expected syntactic scope to be rooted at CFG root node of %s, but instead got syntactic"+
			" scope root node of %s",
		cfg.GetEntry(),
		syntacticScope.GetRootNode())
	if compiler == nil {
		panic(jutil.NewNullPointerException())
	}
	ti.compiler = compiler
	registry := jutil.CheckNotNullPtr(compiler.GetTypeRegistry())
	if ti.registry != registry {
		ti.registry = registry
		ti.unknownType = registry.GetNativeObjectType(rhino.JSTypeNative_UNKNOWN_TYPE)
		ti.numberAdditionSupertype = registry.GetNativeType(rhino.JSTypeNative_NUMBER_ADDITION_SUPERTYPE)
		ti.arrayOfUnknownType = nil
	}
	ti.reverseInterpreter = jutil.CheckNotNullPtr(reverseInterpreter)
	ti.scopeCreator = jutil.CheckNotNullPtr(scopeCreator)
	if scopeCreator.moduleImportResolver != nil {
		ti.moduleImportResolver = scopeCreator.moduleImportResolver
	} else {
		ti.moduleImportResolver =
			NewModuleImportResolver(
				compiler.GetModuleMap(), scopeCreator.GetNodeToScopeMapper(), ti.registry)
	}
	ti.containerScope = syntacticScope
	ti.assertionFunctionLookup = jutil.CheckNotNullPtr(assertionFunctionLookup)

	rootNode := syntacticScope.GetRootNode()
	clear(ti.bottomImplicitVars[:])
	bottomID := ti.bottomTypedScope.scopeID
	if bottomID == 0 {
		bottomID = allocTypedScopeID()
	}
	ti.bottomTypedScope = TypedScope{}
	ti.bottomTypedScope.initScopeCommon(rootNode)
	ti.bottomTypedScope.scopeID = bottomID
	ti.bottomTypedScope.implicitVars = &ti.bottomImplicitVars
	ti.bottomTypedScope.CheckRootScope()
	ti.bottomTypedScope.isBottom = true
	ti.bottomFlowScope = LinkedFlowScope{
		arena:          &ti.flowScopeArena,
		inputProvider:  compiler,
		scopes:         LinkedFlowScope_EMPTY_SCOPES,
		syntacticScope: &ti.bottomTypedScope,
		functionScope:  &ti.bottomTypedScope,
	}
	ti.bottomScope = &ti.bottomFlowScope
}

// GetCfg returns the control flow graph analyzed by TypeInference.
//
// Java: DataFlowAnalysis.getCfg()
func (ti *TypeInference) GetCfg() *ControlFlowGraph[rhino.Node] {
	return ti.cfg
}

// GetImplicitReturnScope returns the joined input flow scope at the CFG's
// implicit return node after Analyze completes.
func (ti *TypeInference) GetImplicitReturnScope() *LinkedFlowScope {
	return ti.nodeStates[ti.cfg.implicitReturn.Index].in
}

func (ti *TypeInference) stepCountAt(index uint32) int {
	return ti.nodeStates[index].stepCount
}

// Join merges two FlowScopes directly without interface dispatch.
//
// Java: DataFlowAnalysis.join(FlowScope, FlowScope)
func (ti *TypeInference) Join(latticeA, latticeB *LinkedFlowScope) *LinkedFlowScope {
	ti.flowJoiner.result = nil
	ti.flowJoiner.inputProvider = ti.compiler
	ti.flowJoiner.JoinFlow(latticeA)
	ti.flowJoiner.JoinFlow(latticeB)
	return ti.flowJoiner.Finish()
}

// Analyze computes a forward branched fixed-point flow solution over ti.cfg.
//
// Java: DataFlowAnalysis.analyze()
func (ti *TypeInference) Analyze() {
	cfg := ti.cfg
	entry := cfg.entry
	implicitReturn := cfg.implicitReturn
	cfgNodes := &cfg.Nodes

	ti.worklist.reset(cfg, cfgNodes.Size())
	nodeCap := cfg.NodeIndexCap()
	if cap(ti.nodeStates) < nodeCap {
		ti.nodeStates = make([]tiNodeState, nodeCap)
	} else {
		clear(ti.nodeStates[:nodeCap])
		ti.nodeStates = ti.nodeStates[:nodeCap]
	}

	bottom := jutil.CheckNotNullPtr(ti.bottomScope)
	for node := range cfgNodes.Values() {
		node.DfaInQueue = false
		ti.nodeStates[node.Index] = tiNodeState{in: bottom, out: bottom}
		if node != implicitReturn {
			ti.worklist.enqueue(node)
		}
	}

	edgeCap := cfg.EdgeIndexCap()
	if cap(ti.edgeStates) < edgeCap {
		ti.edgeStates = make([]*LinkedFlowScope, edgeCap)
	} else {
		clear(ti.edgeStates[:edgeCap])
		ti.edgeStates = ti.edgeStates[:edgeCap]
	}
	for node := range cfgNodes.Values() {
		for _, edge := range node.OutEdgeList.Slice() {
			ti.edgeStates[edge.Index] = bottom
		}
	}

	for ti.worklist.notEmpty() {
		curNode := ti.worklist.dequeue()
		curState := &ti.nodeStates[curNode.Index]
		stepCount := curState.stepCount
		curState.stepCount++
		if stepCount > DataFlowAnalysis_MAX_STEPS_PER_NODE {
			panic(jutil.NewIllegalStateExceptionMsg("Dataflow analysis appears to diverge around: " + jutil.ToJavaString(curNode)))
		}

		ti.joinNodeInputs(curNode, curState, entry)
		if ti.flowNode(curNode, curState) {
			for _, edge := range curNode.OutEdgeList.Slice() {
				nextNode := edge.DestNode
				if !nextNode.DfaInQueue && nextNode != implicitReturn {
					ti.worklist.enqueue(nextNode)
				}
			}
		}
	}

	if implicitReturn != nil {
		ti.joinNodeInputs(implicitReturn, &ti.nodeStates[implicitReturn.Index], entry)
	}
}

func (ti *TypeInference) joinNodeInputs(
	node *graph.LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch],
	state *tiNodeState,
	entry *graph.LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch],
) {
	if node == entry {
		state.in = jutil.CheckNotNullPtr(ti.CreateEntryLattice())
		return
	}
	inEdges := node.InEdgeList.Slice()
	switch len(inEdges) {
	case 0:
		return
	case 1:
		state.in = jutil.CheckNotNullPtr(ti.edgeStates[inEdges[0].Index])
	default:
		ti.flowJoiner.result = nil
		ti.flowJoiner.inputProvider = ti.compiler
		for _, inEdge := range inEdges {
			ti.flowJoiner.JoinFlow(ti.edgeStates[inEdge.Index])
		}
		state.in = jutil.CheckNotNullPtr(ti.flowJoiner.Finish())
	}
}

func (ti *TypeInference) flowNode(
	node *graph.LinkedDiGraphNode[rhino.Node, ControlFlowGraph_Branch],
	state *tiNodeState,
) bool {
	outBefore := state.out
	outAfter := jutil.CheckNotNullPtr(ti.FlowThrough(node.Value, state.in))
	state.out = outAfter
	changed := !outBefore.EqualScope(outAfter)

	// NOTE(nicksantos): Right now, we just treat ON_EX edges like UNCOND
	// edges. If we wanted to be perfect, we'd actually JOIN all the out
	// lattices of this flow with the in lattice, and then make that the out
	// lattice for the ON_EX edge. But it's probably too expensive to be
	// worthwhile.
	ti.branchSource = node.Value
	ti.branchOutput = outAfter
	ti.branchCondition = rhino.NilNode
	ti.branchConditionFlowScope = nil
	ti.branchConditionOutcomes = TypeInference_BooleanOutcomePair{}
	ti.hasBranchOutcomes = false
	for _, le := range node.OutEdgeList.Slice() {
		outBranchBefore := ti.edgeStates[le.Index]
		outBranchAfter := ti.branchFlow(le.Value)
		ti.edgeStates[le.Index] = outBranchAfter
		if !changed {
			changed = !outBranchBefore.EqualScope(outBranchAfter)
		}
	}
	return changed
}

// Java: TypeInference.inferDeclarativelyUnboundVarsWithoutTypes(FlowScope)
func (ti *TypeInference) InferDeclarativelyUnboundVarsWithoutTypes(flow *LinkedFlowScope) *LinkedFlowScope {
	scope := flow.syntacticScope
	if len(scope.vars) == 0 || ti.inferredUnboundVarScope == scope {
		return flow
	}
	if ti.inferredUnboundVarScope == nil {
		ti.inferredUnboundVarScope = scope
	} else {
		for _, s := range ti.inferredUnboundVarList {
			if s == scope {
				return flow
			}
		}
		if len(ti.inferredUnboundVarList) < len(ti.inferredUnboundVarBuf) {
			ti.inferredUnboundVarList = append(ti.inferredUnboundVarList, scope)
		} else {
			if ti.inferredUnboundVars == nil {
				ti.inferredUnboundVars = make(map[*TypedScope]struct{}, 16)
			} else if _, exists := ti.inferredUnboundVars[scope]; exists {
				return flow
			}
			ti.inferredUnboundVars[scope] = struct{}{}
		}
	}
	// For each local variable declared with the VAR keyword, the entry
	// type is VOID.
	var overlay *LinkedFlowScope_OverlayScope
	var slots *rhino.HamtPMap[string, rhino.JSType]
	var voidType rhino.JSType
	for var_ := range scope.GetVarIterable() {
		if !scope.IsDeclarativelyUnboundVarWithoutType(var_) || ti.IsUnflowable(var_) {
			continue
		}
		if voidType == nil {
			voidType = ti.GetNativeType(rhino.JSTypeNative_VOID_TYPE)
			overlay = flow.scopes.Get(scope)
			if overlay != nil {
				slots = overlay.slots
			} else {
				slots = LinkedFlowScope_EMPTY_SLOTS
			}
		}
		name := var_.GetName()
		if slots.Get(name) != voidType {
			slots = slots.PlusMapWithArena(flow.arena.slotMaps(), name, voidType)
		}
	}
	if voidType == nil {
		return flow
	}
	if overlay != nil && slots == overlay.slots {
		return flow
	}
	if overlay == nil && slots.IsEmpty() {
		return flow
	}
	newScope := newLinkedFlowScope_OverlayScopeWithArena(flow.arena, scope, slots)
	var newScopes *rhino.HamtPMap[*TypedScope, *LinkedFlowScope_OverlayScope]
	if !newScope.slots.IsEmpty() {
		newScopes = flow.scopes.PlusMapWithArena(flow.arena.scopeMaps(), scope, newScope)
	} else {
		newScopes = flow.scopes.MinusMapWithArena(flow.arena.scopeMaps(), scope)
	}
	if newScopes != flow.scopes {
		return newLinkedFlowScopeWithArena(flow.arena, flow.inputProvider, newScopes, flow.syntacticScope, flow.functionScope)
	}
	return flow
}

// Java: TypeInference.createInitialEstimateLattice()
func (ti *TypeInference) CreateInitialEstimateLattice() *LinkedFlowScope {
	return ti.bottomScope
}

// Java: TypeInference.createEntryLattice()
func (ti *TypeInference) CreateEntryLattice() *LinkedFlowScope {
	// only ever called once so we don't need to cache this computation
	entryScope :=
		ti.InferDeclarativelyUnboundVarsWithoutTypes(
			newLinkedFlowScopeWithArena(&ti.flowScopeArena, ti.compiler, LinkedFlowScope_EMPTY_SCOPES, ti.containerScope, ti.containerScope))

	return ti.InferParameters(entryScope)
}

func typeInference_getEnclosingScopeRoot(n rhino.Node) rhino.Node {
	for ; !n.IsNil(); n = n.GetParent() {
		if NodeUtil_CreatesScope(n) {
			return n
		}
	}
	return rhino.NilNode
}

// Java: TypeInference.flowThrough(Node, FlowScope)
func (ti *TypeInference) FlowThrough(n rhino.Node, input *LinkedFlowScope) *LinkedFlowScope {
	// If we have not walked a path from <entry> to <n>, then we don't
	// want to infer anything about this scope.
	if input == ti.bottomScope {
		return input
	}

	// This method also does some logic for ES modules right before and after entering/exiting the
	// scope rooted at the module. The reasoning for separating out this logic is that we can just
	// ignore the actual AST nodes for IMPORT/EXPORT, in most cases, because we have already
	// created an abstraction of imports and exports.
	root := typeInference_getEnclosingScopeRoot(n)
	// Inferred types of ES module imports/exports aren't knowable until after TypeInference runs.
	// First update the type of all imports in the scope, then do flow-sensitive inference, then
	// update the implicit '*exports*' object.
	var module *Module
	var syntacticBlockScope *TypedScope
	if root == ti.lastScopeRoot && !root.IsNil() {
		module = ti.lastModule
		syntacticBlockScope = ti.lastSyntacticBlockScope
	} else {
		module = ModuleImportResolver_GetModuleFromScopeRoot(ti.compiler.GetModuleMap(), ti.compiler, root)
		syntacticBlockScope = ti.scopeCreator.CreateScopeOfNode(root)
		ti.lastScopeRoot = root
		ti.lastModule = module
		ti.lastSyntacticBlockScope = syntacticBlockScope
	}
	if module != nil && module.Metadata().IsEs6Module() {
		ti.moduleImportResolver.DeclareEsModuleImports(
			module, syntacticBlockScope, ti.compiler.GetInput(NodeUtil_GetInputId(n)))
	}

	// This logic is not specific to ES modules.
	output := input.withTypedScope(syntacticBlockScope)
	output = ti.InferDeclarativelyUnboundVarsWithoutTypes(output)
	output = ti.Traverse(n, output)

	ti.UpdateModuleScope(module, syntacticBlockScope)
	return output
}

// UpdateModuleScope updates the given scope after running inference over a
// goog.module or ES module.
//
// Java: TypeInference.updateModuleScope(Module, TypedScope)
func (ti *TypeInference) UpdateModuleScope(module *Module, syntacticBlockScope *TypedScope) {
	if module == nil {
		return
	}
	switch module.Metadata().ModuleType() {
	case ModuleMetadataMap_ModuleType_ES6_MODULE:
		// This call only affects exports with an inferred, not declared, type. Declared exports
		// were already added to the namespace object type in TypedScopeCreator.
		ti.moduleImportResolver.UpdateEsModuleNamespaceType(
			syntacticBlockScope.GetVar(Export_NAMESPACE).GetType().ToObjectType(),
			module,
			syntacticBlockScope)
		return
	case ModuleMetadataMap_ModuleType_GOOG_MODULE, ModuleMetadataMap_ModuleType_LEGACY_GOOG_MODULE:
		exportsVar := syntacticBlockScope.GetVarByID(implicitExportsID)
		if exportsVar == nil {
			panic(jutil.NewNullPointerExceptionMsg(jutil.StringFormat("Missing exports var for %s", module.Metadata())))
		}
		var exportsType rhino.JSType
		if exportsVar.GetType() != nil {
			exportsType = exportsVar.GetType()
		} else {
			exportsType = ti.unknownType
		}
		// Store the type of the namespace on the AST for the convenience of later passes that want
		// to access it.
		rootNode := syntacticBlockScope.GetRootNode()
		if rootNode.IsModuleBody() {
			rootNode.SetJSType(exportsType)
		} else {
			// For goog.loadModule, give the `exports` parameter the correct type.
			jutil.CheckStateMsg1(rootNode.IsBlock(), "%s", rootNode)
			paramList := NodeUtil_GetFunctionParameters(rootNode.GetParent())
			paramList.GetOnlyChild().SetJSType(exportsType)
		}
		if !module.Metadata().IsLegacyGoogModule() || exportsVar.GetType() == nil {
			break
		}

		// Update the global scope for the implicit assignment "legacy.module.id = exports;" created
		// by `goog.module.declareLegacyNamespace();`
		moduleId := *module.ClosureNamespace()
		globalScope := syntacticBlockScope.GetGlobalScope()
		globalVar := globalScope.GetVar(moduleId)
		if globalVar.IsTypeInferred() {
			globalVar.SetType(exportsType)
		}
		// Update the property slot on the parent namespace.
		moduleQname := rhino.QualifiedName_Of(moduleId)
		if !moduleQname.IsSimple() {
			parentType := globalScope.LookupQualifiedName(moduleQname.GetOwner())
			var parentObjectType rhino.ObjectType
			if parentType != nil {
				parentObjectType = parentType.ToMaybeObjectType()
			}
			if parentObjectType != nil &&
				!parentObjectType.IsPropertyTypeDeclared(moduleQname.GetComponent()) {
				parentObjectType.DefineInferredProperty(
					moduleQname.GetComponent(), exportsType, exportsVar.GetNode())
			}
		}
		return
	default:
	}
}

func (ti *TypeInference) branchFlow(branch ControlFlowGraph_Branch) *LinkedFlowScope {
	source := ti.branchSource
	output := ti.branchOutput
	switch branch {
	case ControlFlowGraph_Branch_ON_TRUE:
		if NodeUtil_IsEnhancedFor(source) {
			return ti.InitializeEnhancedForScope(source, output)
		}
		// FALL THROUGH
		fallthrough

	case ControlFlowGraph_Branch_ON_FALSE:
		if ti.branchCondition.IsNil() {
			if source.IsCase() {
				ti.branchCondition = source
				ti.branchConditionFlowScope = ti.Traverse(ti.branchCondition.GetFirstChild(), output)
			} else {
				ti.branchCondition = NodeUtil_GetConditionExpression(source)
				if ti.branchCondition.IsNil() {
					return output
				}
			}
		}

		if ti.branchCondition.IsAnd() || ti.branchCondition.IsOr() {
			// When handling the short-circuiting binary operators,
			// the outcome scope on true can be different than the outcome
			// scope on false.
			//
			// TODO(nicksantos): The "right" way to do this is to
			// carry the known outcome all the way through the
			// recursive traversal, so that we can construct a
			// different flow scope based on the outcome. However,
			// this would require a bunch of code and a bunch of
			// extra computation for an edge case. This seems to be
			// a "good enough" approximation.

			// conditionOutcomes is cached from previous calls to the brancher
			if !ti.hasBranchOutcomes {
				if ti.branchCondition.IsAnd() {
					ti.branchConditionOutcomes = ti.TraverseAnd(ti.branchCondition, output)
				} else {
					ti.branchConditionOutcomes = ti.TraverseOr(ti.branchCondition, output)
				}
				ti.hasBranchOutcomes = true
			}
			return ti.reverseInterpreter.GetPreciserScopeKnowingConditionOutcome(
				ti.branchCondition,
				ti.branchConditionOutcomes.GetOutcomeFlowScope(
					ti.branchCondition.GetToken(), branch == ControlFlowGraph_Branch_ON_TRUE),
				rhino.Outcome_ForBoolean(branch == ControlFlowGraph_Branch_ON_TRUE))
		}

		// conditionFlowScope is cached from previous calls to the brancher
		if ti.branchConditionFlowScope == nil {
			ti.branchConditionFlowScope = ti.Traverse(ti.branchCondition, output)
		}
		return ti.reverseInterpreter.GetPreciserScopeKnowingConditionOutcome(
			ti.branchCondition, ti.branchConditionFlowScope, rhino.Outcome_ForBoolean(branch == ControlFlowGraph_Branch_ON_TRUE))

	default:
		return output
	}
}

// Java: TypeInference.traverse(Node, FlowScope)
func (ti *TypeInference) Traverse(n rhino.Node, scope *LinkedFlowScope) *LinkedFlowScope {
	isTypeable := true
	switch n.GetToken() {
	case rhino.Token_ASSIGN:
		scope = ti.TraverseAssign(n, scope)
	case rhino.Token_NAME:
		scope = ti.TraverseName(n, scope)
	case rhino.Token_OPTCHAIN_GETPROP, rhino.Token_OPTCHAIN_CALL, rhino.Token_OPTCHAIN_GETELEM:
		scope = ti.TraverseOptChain(n, scope)
	case rhino.Token_GETPROP:
		scope = ti.TraverseGetProp(n, scope)
	case rhino.Token_CLASS:
		scope = ti.TraverseClass(n, scope)
	case rhino.Token_ASSIGN_AND, rhino.Token_ASSIGN_OR:
		scope = ti.TraverseShortCircuitingBinOpAssignment(n, scope)
	case rhino.Token_AND:
		pair := ti.TraverseAnd(n, scope)
		scope = pair.GetJoinedFlowScope()
	case rhino.Token_OR:
		pair := ti.TraverseOr(n, scope)
		scope = pair.GetJoinedFlowScope()
	case rhino.Token_ASSIGN_COALESCE, rhino.Token_COALESCE:
		scope = ti.TraverseNullishCoalesce(n, scope)
	case rhino.Token_HOOK:
		scope = ti.TraverseHook(n, scope)
	case rhino.Token_OBJECTLIT:
		scope = ti.TraverseObjectLiteral(n, scope)
	case rhino.Token_CALL:
		scope = ti.TraverseCall(n, scope)
	case rhino.Token_NEW:
		scope = ti.TraverseNew(n, scope)
	case rhino.Token_NEW_TARGET:
		ti.TraverseNewTarget(n)
	case rhino.Token_ASSIGN_ADD, rhino.Token_ADD:
		scope = ti.TraverseAdd(n, scope)
	case rhino.Token_POS:
		scope = ti.TraverseUnaryPlus(n, scope)
	case rhino.Token_NEG, rhino.Token_BITNOT, rhino.Token_DEC, rhino.Token_INC:
		scope = ti.TraverseBigIntCompatibleUnaryOperator(n, scope)
	case rhino.Token_ARRAYLIT:
		scope = ti.TraverseArrayLiteral(n, scope)
	case rhino.Token_THIS:
		n.SetJSType(scope.GetTypeOfThis())
	case rhino.Token_ASSIGN_LSH,
		rhino.Token_ASSIGN_RSH,
		rhino.Token_ASSIGN_DIV,
		rhino.Token_ASSIGN_MOD,
		rhino.Token_ASSIGN_BITAND,
		rhino.Token_ASSIGN_BITXOR,
		rhino.Token_ASSIGN_BITOR,
		rhino.Token_ASSIGN_MUL,
		rhino.Token_ASSIGN_SUB,
		rhino.Token_ASSIGN_EXPONENT:
		scope = ti.TraverseAssignOp(n, scope)
	case rhino.Token_ASSIGN_URSH:
		// >>> is not compatible with BigInt
		scope = ti.TraverseAssignUnsignedRightShift(n, scope)
	case rhino.Token_BITAND, rhino.Token_BITXOR, rhino.Token_BITOR, rhino.Token_LSH, rhino.Token_RSH,
		rhino.Token_SUB, rhino.Token_MUL, rhino.Token_DIV, rhino.Token_MOD, rhino.Token_EXPONENT:
		scope = ti.TraverseBigIntCompatibleBinaryOperator(n, scope)
	case rhino.Token_URSH:
		// >>> is not compatible with BigInt
		scope = ti.TraverseUnsignedRightShift(n, scope)
	case rhino.Token_COMMA:
		scope = ti.TraverseChildren(n, scope)
		n.SetJSType(ti.GetJSType(ast.MustBinary(n).Right().Node()))
	case rhino.Token_TEMPLATELIT, rhino.Token_TYPEOF:
		scope = ti.TraverseChildren(n, scope)
		n.SetJSType(ti.GetNativeType(rhino.JSTypeNative_STRING_TYPE))
	case rhino.Token_TEMPLATELIT_SUB,
		rhino.Token_THROW,
		rhino.Token_ITER_SPREAD,
		rhino.Token_OBJECT_SPREAD,
		rhino.Token_IMPORT,
		rhino.Token_IMPORT_SPECS,
		rhino.Token_IMPORT_STAR:
		// these nodes are untyped but have children that may affect the flow scope and need to
		// be typed.
		scope = ti.TraverseChildren(n, scope)
		isTypeable = false
	case rhino.Token_IMPORT_SPEC:
		// these nodes are untyped but have children that need to be typed.
		ti.TraverseImportSpec(scope, n)
		isTypeable = false
	case rhino.Token_TAGGED_TEMPLATELIT:
		scope = ti.TraverseTaggedTemplateLit(n, scope)
	case rhino.Token_DELPROP, rhino.Token_LT, rhino.Token_LE, rhino.Token_GT, rhino.Token_GE, rhino.Token_NOT,
		rhino.Token_EQ, rhino.Token_NE, rhino.Token_SHEQ, rhino.Token_SHNE, rhino.Token_INSTANCEOF, rhino.Token_IN:
		scope = ti.TraverseChildren(n, scope)
		n.SetJSType(ti.GetNativeType(rhino.JSTypeNative_BOOLEAN_TYPE))
	case rhino.Token_GETELEM:
		scope = ti.TraverseGetElem(n, scope)
	case rhino.Token_EXPR_RESULT:
		scope = ti.TraverseChildren(n, scope)
		if getprop, ok := ast.MustExprStmt(n).Expr().AsGetProp(); ok {
			ownerType :=
				rhino.ObjectType_Cast(ti.GetJSType(getprop.Object().Node()).RestrictByNotNullOrUndefined())
			if ownerType != nil {
				ti.EnsurePropertyDeclaredHelper(getprop.Node(), ownerType, scope)
			}
		}
		isTypeable = false
	case rhino.Token_SWITCH:
		scope = ti.Traverse(ast.MustSwitch(n).Discriminant().Node(), scope)
		isTypeable = false
	case rhino.Token_RETURN:
		scope = ti.TraverseReturn(n, scope)
		isTypeable = false
	case rhino.Token_YIELD:
		if n.IsYieldAll() {
			scope = ti.TraverseYieldAll(n, scope)
		} else {
			scope = ti.TraverseChildren(n, scope)
			n.SetJSType(ti.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE))
		}
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		scope = ti.TraverseDeclaration(n, scope)
		isTypeable = false
	case rhino.Token_CATCH:
		scope = ti.TraverseCatch(n, scope)
		isTypeable = false
	case rhino.Token_CAST:
		scope = ti.TraverseChildren(n, scope)
		info := n.GetJSDocInfo()
		// TODO(b/123955687): also check that info.hasType() is true
		if info == nil {
			panic(jutil.NewNullPointerExceptionMsg("CAST node should always have JSDocInfo"))
		}
		if info.HasType() {
			// NOTE(lharker) - I tried moving CAST type evaluation into the typed scope creation
			// phase.
			// Since it caused a few new, seemingly spurious, 'Bad type annotation' and
			// 'unknown property type' warnings, and having it in TypeInference seems to work, we just
			// do the lookup + resolution here.
			n.SetJSType(
				info.GetType().
					Evaluate(scope.GetDeclarationScope(), ti.registry).
					Resolve(ti.registry.GetErrorReporter()))
		} else {
			n.SetJSType(ti.unknownType)
		}
	case rhino.Token_SUPER:
		ti.TraverseSuper(n, scope)
	case rhino.Token_AWAIT:
		scope = ti.TraverseAwait(n, scope)
	case rhino.Token_VOID:
		n.SetJSType(ti.GetNativeType(rhino.JSTypeNative_VOID_TYPE))
		scope = ti.TraverseChildren(n, scope)
	case rhino.Token_EXPORT:
		scope = ti.TraverseChildren(n, scope)
		exp := ast.MustExport(n)
		if exp.IsDefault() {
			// TypedScopeCreator declared a dummy variable *default* to store this type. Update the
			// variable with the inferred type.
			defaultExport := TypeInference_GetDeclaredVar(scope, Export_DEFAULT_EXPORT_NAME)
			if defaultExport.IsTypeInferred() {
				defaultExport.SetType(ti.GetJSType(exp.Body().Node()))
			}
		}
		isTypeable = false
	case rhino.Token_IMPORT_META:
		// TODO(b/137797083): Set an appropriate type.
		n.SetJSType(ti.unknownType)
	case rhino.Token_ROOT,
		rhino.Token_SCRIPT,
		rhino.Token_MODULE_BODY,
		rhino.Token_SWITCH_BODY,
		rhino.Token_FUNCTION,
		rhino.Token_PARAM_LIST,
		rhino.Token_BLOCK,
		rhino.Token_EMPTY,
		rhino.Token_IF,
		rhino.Token_WHILE,
		rhino.Token_DO,
		rhino.Token_FOR,
		rhino.Token_FOR_IN,
		rhino.Token_FOR_OF,
		rhino.Token_FOR_AWAIT_OF,
		rhino.Token_BREAK,
		rhino.Token_CONTINUE,
		rhino.Token_TRY,
		rhino.Token_CASE,
		rhino.Token_DEFAULT_CASE,
		rhino.Token_WITH,
		rhino.Token_DEBUGGER,
		rhino.Token_EXPORT_SPECS,
		rhino.Token_LABEL:
		// These don't need to be typed here, since they only affect control flow.
		isTypeable = false
	case rhino.Token_DYNAMIC_IMPORT:
		ti.TraverseDynamicImport(n, scope)
	case rhino.Token_TRUE, rhino.Token_FALSE, rhino.Token_STRINGLIT, rhino.Token_NUMBER, rhino.Token_BIGINT,
		rhino.Token_NULL, rhino.Token_REGEXP, rhino.Token_TEMPLATELIT_STRING:
		// Primitives are typed in TypedScopeCreator.AbstractScopeBuilder#attachLiteralTypes
	default:
		panic(jutil.NewIllegalStateExceptionMsg(
			"Type inference doesn't know to handle token " + n.GetToken().String()))
	}

	if isTypeable && n.GetJSType() == nil && !typeInference_allowsNullType(n.GetToken()) {
		panic(jutil.NewIllegalStateExceptionMsg("Failed to infer JSType for " + n.String()))
	}
	return scope
}

// TODO(b/154044898): delete these exceptions. Names are given null types to enable inference of
// undeclared names assigned in multiple local scopes. The compiler also infers call and new
// types when the invocation target is such a name.
func typeInference_allowsNullType(t rhino.Token) bool {
	return t == rhino.Token_NAME || t == rhino.Token_CALL || t == rhino.Token_NEW
}

// Java: TypeInference.initializeEnhancedForScope(Node, FlowScope)
func (ti *TypeInference) InitializeEnhancedForScope(source rhino.Node, output *LinkedFlowScope) *LinkedFlowScope {
	enhFor := ast.MustEnhancedFor(source)
	target := enhFor.Target()
	item := target.Node()
	obj := enhFor.Collection().Node()

	informed := ti.Traverse(obj, output)

	var assignmentType TypeInference_AssignmentType
	if decl, ok := target.AsNameDeclaration(); ok {
		item = decl.MustFirstDeclarator().Node()
		assignmentType = TypeInference_AssignmentType_DECLARATION
	} else {
		assignmentType = TypeInference_AssignmentType_ASSIGN
	}
	if lhs, ok := ast.AsDestructuringLhs(item); ok {
		item = lhs.Pattern().Node()
	}

	var newType rhino.JSType
	switch source.GetToken() {
	case rhino.Token_FOR_IN:
		// item is assigned a property name, so its type should be string
		iterKeyType := ti.GetNativeType(rhino.JSTypeNative_STRING_TYPE)
		objType := ti.GetJSType(obj).Autobox()
		objIndexType :=
			objType.GetTemplateTypeMap().GetResolvedTemplateType(ti.registry.GetObjectIndexKey())
		if objIndexType != nil && !objIndexType.IsUnknownType() {
			narrowedKeyType := iterKeyType.GetGreatestSubtype(objIndexType)
			if !narrowedKeyType.IsEmptyType() {
				iterKeyType = narrowedKeyType
			}
		}

		newType = iterKeyType
	case rhino.Token_FOR_OF:
		// for/of. The type of `item` is the type parameter of the Iterable type.
		newType =
			JsIterables_MaybeBoxIterableOrAsyncIterable(ti.GetJSType(obj), ti.registry).
				OrElse(ti.unknownType)
	case rhino.Token_FOR_AWAIT_OF:
		// for/await/of. the iterated object is either of the Iterable or AsyncIterable type.
		// the type of `item` is the Promise.resolve() type of the object's type parameter.
		iterableType :=
			JsIterables_MaybeBoxIterableOrAsyncIterable(ti.GetJSType(obj), ti.registry).
				OrElse(ti.unknownType)

		newType = Promises_GetResolvedType(ti.registry, iterableType)
	default:
		panic(jutil.NewIllegalArgumentExceptionMsg("Unexpected source node " + source.String()))
	}

	// Note that `item` can be an arbitrary LHS expression we need to check.
	if item.IsDestructuringPattern() {
		// for (const {x, y} of data) {
		informed = ti.TraverseDestructuringPattern(item, informed, newType, assignmentType)
	} else {
		informed = ti.Traverse(item, informed)
		informed = ti.UpdateScopeForAssignment(informed, item, newType, assignmentType)
	}
	return informed
}

// TraverseReturn traverses a return value.
//
// Java: TypeInference.traverseReturn(Node, FlowScope)
func (ti *TypeInference) TraverseReturn(n rhino.Node, scope *LinkedFlowScope) *LinkedFlowScope {
	scope = ti.TraverseChildren(n, scope)

	if retValue, ok := ast.MustReturn(n).Value(); ok {
		type_ := ti.containerScope.GetRootNode().GetJSType()
		if type_ != nil {
			fnType := type_.ToMaybeFunctionType()
			if fnType != nil {
				TypeInference_InferPropertyTypesToMatchConstraint(retValue.Node().GetJSType(), fnType.GetReturnType())
			}
		}
	}
	return scope
}

// Java: TypeInference.isUnflowable(TypedVar)
func (ti *TypeInference) IsUnflowable(v *TypedVar) bool {
	return v != nil &&
		v.IsLocal() &&
		v.IsMarkedEscaped() &&
		// It's OK to flow a variable in the scope where it's escaped.
		v.GetScope().(*TypedScope).GetClosestContainerScope() == ti.containerScope
}

// GetJSType gets the JSType from the Node argument and verifies that it is
// present.
//
// Java: TypeInference.getJSType(Node)
func (ti *TypeInference) GetJSType(n rhino.Node) rhino.JSType {
	jsType := n.GetJSType()
	if jsType == nil {
		// TODO(nicksantos): This branch indicates a compiler bug, not worthy of
		// halting the compilation but we should log this and analyze to track
		// down why it happens. This is not critical and will be resolved over
		// time as the type checker is extended.
		return ti.unknownType
	}
	return jsType
}

// Java: TypeInference.getNativeType(JSTypeNative)
func (ti *TypeInference) GetNativeType(typeId rhino.JSTypeNative) rhino.JSType {
	return ti.registry.GetNativeType(typeId)
}
