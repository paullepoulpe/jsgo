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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// TypeInference is type inference within a script node or a function body,
// using the data-flow analysis framework.
type TypeInference struct {
	DataFlowAnalysisBase[rhino.Node, FlowScope]

	compiler           AbstractCompiler
	registry           *rhino.JSTypeRegistry
	reverseInterpreter ReverseAbstractInterpreter
	bottomScope        FlowScope
	// global scope, function scope, or static block scope
	containerScope          *TypedScope
	scopeCreator            *TypedScopeCreator
	assertionFunctionLookup *CodingConvention_AssertionFunctionLookup
	moduleImportResolver    *ModuleImportResolver
	// A record is pushed onto this stack during the traversal of each optional chain.
	// Inference of the type at the end of the optional chain requires scope information
	// from traversing the start of the chain. This stack serves as a communication channel
	// for that information.
	//
	// Go-only: language-adapter — the ArrayDeque is only used as a stack
	// (addFirst/peekFirst/removeFirst), so it is a slice whose last element is
	// the deque's first.
	optChainArrayDeque []*TypeInference_OptChainInfo

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
	flowBrancher             DataFlowAnalysis_FlowBrancher[FlowScope]
	branchSource             rhino.Node
	branchOutput             FlowScope
	branchCondition          rhino.Node
	branchConditionFlowScope FlowScope
	branchConditionOutcomes  *TypeInference_BooleanOutcomePair

	templatedMatched  jutil.LinkedHashMap[*rhino.TemplateType, rhino.JSType]
	templatedInferred jutil.LinkedHashMap[*rhino.TemplateType, rhino.JSType]

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
	cfg ControlFlowGraph[rhino.Node],
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
	ti.releaseLargeStorage()
}

func (ti *TypeInference) initTypeInference(
	compiler AbstractCompiler,
	cfg ControlFlowGraph[rhino.Node],
	reverseInterpreter ReverseAbstractInterpreter,
	syntacticScope *TypedScope,
	scopeCreator *TypedScopeCreator,
	assertionFunctionLookup *CodingConvention_AssertionFunctionLookup) {
	ti.InitDataFlowAnalysisBase(ti, cfg)
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
	ti.optChainArrayDeque = ti.optChainArrayDeque[:0]
	// Constructor body.
	jutil.CheckArgumentMsg2(
		syntacticScope.GetRootNode() == cfg.GetEntry().GetValue(),
		"Expected syntactic scope to be rooted at CFG root node of %s, but instead got syntactic"+
			" scope root node of %s",
		cfg.GetEntry(),
		syntacticScope.GetRootNode())
	ti.compiler = jutil.CheckNotNull(compiler)
	registry := jutil.CheckNotNull(compiler.GetTypeRegistry())
	if ti.registry != registry {
		ti.registry = registry
		ti.unknownType = registry.GetNativeObjectType(rhino.JSTypeNative_UNKNOWN_TYPE)
		ti.numberAdditionSupertype = registry.GetNativeType(rhino.JSTypeNative_NUMBER_ADDITION_SUPERTYPE)
		ti.arrayOfUnknownType = nil
	}
	ti.reverseInterpreter = jutil.CheckNotNull(reverseInterpreter)
	ti.scopeCreator = jutil.CheckNotNull(scopeCreator)
	if scopeCreator.moduleImportResolver != nil {
		ti.moduleImportResolver = scopeCreator.moduleImportResolver
	} else {
		ti.moduleImportResolver =
			NewModuleImportResolver(
				compiler.GetModuleMap(), scopeCreator.GetNodeToScopeMapper(), ti.registry)
	}
	ti.containerScope = syntacticScope
	ti.assertionFunctionLookup = jutil.CheckNotNull(assertionFunctionLookup)

	rootNode := syntacticScope.GetRootNode()
	clear(ti.bottomImplicitVars[:])
	ti.bottomTypedScope = TypedScope{}
	ti.bottomTypedScope.InitAbstractScopeBase(&ti.bottomTypedScope, rootNode)
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

// Java: TypeInference.inferDeclarativelyUnboundVarsWithoutTypes(FlowScope)
func (ti *TypeInference) InferDeclarativelyUnboundVarsWithoutTypes(flow FlowScope) FlowScope {
	scope := flow.GetDeclarationScope().(*TypedScope)
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
	l := flow.(*LinkedFlowScope)
	var overlay *LinkedFlowScope_OverlayScope
	var slots *rhino.HamtPMap[string, *LinkedFlowScope_OverlaySlot]
	var voidType rhino.JSType
	for var_ := range scope.GetVarIterable() {
		if !scope.IsDeclarativelyUnboundVarWithoutType(var_) || ti.IsUnflowable(var_) {
			continue
		}
		if voidType == nil {
			voidType = ti.GetNativeType(rhino.JSTypeNative_VOID_TYPE)
			overlay = l.scopes.Get(scope)
			if overlay != nil {
				slots = overlay.slots
			} else {
				slots = LinkedFlowScope_EMPTY_SLOTS
			}
		}
		name := var_.GetName()
		slot := slots.Get(name)
		if slot == nil || slot.type_ != voidType {
			slots = slots.PlusMapWithArena(l.arena.slotMaps(), name, newLinkedFlowScope_OverlaySlotWithArena(l.arena, name, voidType))
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
	newScope := newLinkedFlowScope_OverlayScopeWithArena(l.arena, scope, slots)
	var newScopes *rhino.HamtPMap[*TypedScope, *LinkedFlowScope_OverlayScope]
	if !newScope.slots.IsEmpty() {
		newScopes = l.scopes.PlusMapWithArena(l.arena.scopeMaps(), scope, newScope)
	} else {
		newScopes = l.scopes.MinusMapWithArena(l.arena.scopeMaps(), scope)
	}
	if newScopes != l.scopes {
		return newLinkedFlowScopeWithArena(l.arena, l.inputProvider, newScopes, l.syntacticScope, l.functionScope)
	}
	return flow
}

// InferParameters infers all of a function's parameters if their types
// aren't declared.
//
// Java: TypeInference.inferParameters(FlowScope)
func (ti *TypeInference) InferParameters(entryFlowScope FlowScope) FlowScope {
	functionNode := ti.containerScope.GetRootNode()
	fn, ok := ast.AsFunction(functionNode)
	if !ok {
		return entryFlowScope // we're in the global scope
	} else if NodeUtil_IsBundledGoogModuleCall(functionNode.GetParent()) {
		// Pretend the function literal in `goog.loadModule(function(exports) {` does not exist.
		return entryFlowScope
	}
	astParameters := fn.Params()
	var iifeArgumentNode rhino.Node

	if NodeUtil_IsInvocationTarget(functionNode) {
		iifeArgumentNode = functionNode.GetNext()
	}

	functionType := rhino.JSType_ToMaybeFunctionType(functionNode.GetJSType())
	parameterTypes := functionType.GetParameters()
	nextParameter := 0
	var parameter *rhino.FunctionType_Parameter
	if nextParameter < len(parameterTypes) {
		parameter = parameterTypes[nextParameter]
		nextParameter++
	}

	// This really iterates over three different things at once:
	//   - the actual AST parameter nodes (which may be REST, DEFAULT_VALUE, etc.)
	//   - the argument nodes in an IIFE
	//   - the parameter type nodes from the FunctionType on the FUNCTION node
	// Always visit every AST parameter once, regardless of how many IIFE arguments or
	// FunctionType param nodes there are.
	for astParamItr := range astParameters.Params() {
		astParam := astParamItr.Node()
		if !iifeArgumentNode.IsNil() && iifeArgumentNode.IsSpread() {
			// block inference on all parameters that might possibly be set by a spread, e.g. `z` in
			// (function f(x, y, z = 1))(...[1, 2], 'foo')
			iifeArgumentNode = rhino.NilNode
		}

		// Running variable for the type of the param within the body of the function. We use the
		// existing type on the param node as the default, and then transform it according to the
		// declaration syntax.
		inferredType := ti.GetJSType(astParam)

		if !iifeArgumentNode.IsNil() {
			if iifeArgumentNode.GetJSType() != nil {
				inferredType = iifeArgumentNode.GetJSType()
			}
		} else if parameter != nil {
			if parameter.GetJSType() != nil {
				inferredType = parameter.GetJSType()
			}
		}

		var defaultValue rhino.Node
		if defVal, ok := astParamItr.AsDefaultValue(); ok {
			defaultValue = defVal.Value().Node()
			// must call `traverse` to correctly type the default value
			entryFlowScope = ti.Traverse(defaultValue, entryFlowScope)
			astParam = defVal.Target().Node()
		} else if rest, ok := astParamItr.AsIterRest(); ok {
			// e.g. `function f(p1, ...restParamName) {}`
			// set astParam = restParamName
			astParam = rest.Target().Node()
			// convert 'number' into 'Array<number>' for rest parameters
			inferredType =
				ti.registry.CreateTemplatizedType(ti.registry.GetNativeObjectType(rhino.JSTypeNative_ARRAY_TYPE), inferredType)
		}

		if !defaultValue.IsNil() {
			// The param could possibly be the default type, and `undefined` args won't propagate in.
			inferredType =
				ti.registry.CreateUnionType(
					inferredType.RestrictByNotUndefined(), ti.GetJSType(defaultValue))
		}

		if astParam.IsDestructuringPattern() {
			// even if the inferredType is null, we still need to type all the nodes inside the
			// destructuring pattern. (e.g. in computed properties or default value expressions)
			entryFlowScope = ti.UpdateDestructuringParameter(astParam, inferredType, entryFlowScope)
		} else {
			// for simple named parameters, we only need to update the scope/AST if we have a new
			// inferred type.
			entryFlowScope =
				ti.UpdateNamedParameter(astParam, !defaultValue.IsNil(), inferredType, entryFlowScope)
		}

		parameter = nil
		if nextParameter < len(parameterTypes) {
			parameter = parameterTypes[nextParameter]
			nextParameter++
		}
		if !iifeArgumentNode.IsNil() {
			iifeArgumentNode = iifeArgumentNode.GetNext()
		}
	}

	return entryFlowScope
}

// UpdateDestructuringParameter sets the types of a un-named/destructuring
// function parameter to an inferred type.
//
// This method is responsible for typing:
//   - The scope slot
//   - The pattern nodes
//
// Java: TypeInference.updateDestructuringParameter(Node, JSType, FlowScope)
func (ti *TypeInference) UpdateDestructuringParameter(
	pattern rhino.Node, inferredType rhino.JSType, entryFlowScope FlowScope) FlowScope {
	// look through all expressions and lvalues in the pattern.
	// given an lvalue, change its type if either a) it's inferred (not declared in
	// TypedScopeCreator) or b) it has a default value
	entryFlowScope =
		ti.TraverseDestructuringPatternHelper(
			pattern,
			entryFlowScope,
			inferredType,
			func(scope FlowScope, lvalue rhino.Node, type_ rhino.JSType) FlowScope {
				var_ := ti.containerScope.GetVar(lvalue.GetString())
				jutil.CheckNotNull(var_)
				// This condition will trigger on cases like
				//   (function f({x}) {})({x: 3})
				// where `x` is of unknown type during the typed scope creation phase, but
				// here we can infer that it is of type `number`
				if var_.IsTypeInferred() {
					var_.SetType(type_)
					lvalue.SetJSType(type_)
				}
				if lvalue.GetParent().IsDefaultValue() {
					// Given
					//   /** @param {{age: (number|undefined)}} data */
					//   function f({age = 99}) {}
					// infer that `age` is now a `number` and not `number|undefined`
					// treat this similarly to if there was an assignment inside the function body
					// TODO(b/117162687): allow people to narrow the declared type to
					// exclude 'undefined' inside the function body.
					scope = ti.UpdateScopeForAssignment(scope, lvalue, type_, TypeInference_AssignmentType_ASSIGN)
				}
				return scope
			})

	return entryFlowScope
}

// UpdateNamedParameter sets the types of a named/non-destructuring function
// parameter to an inferred type.
//
// This method is responsible for typing:
//   - The scope slot
//   - The param node
//
// Java: TypeInference.updateNamedParameter(Node, boolean, JSType, FlowScope)
func (ti *TypeInference) UpdateNamedParameter(
	paramName rhino.Node, hasDefaultValue bool, inferredType rhino.JSType, entryFlowScope FlowScope) FlowScope {
	var_ := ti.containerScope.GetVar(paramName.GetString())
	jutil.CheckNotNullMsg1(var_, "Missing var for parameter %s", paramName)

	paramName.SetJSType(inferredType)

	if var_.IsTypeInferred() {
		var_.SetType(inferredType)
	} else if hasDefaultValue {
		// If this is a declared type with a default value, update the LinkedFlowScope slots but not
		// the actual TypedVar. This is similar to what would happen if the default value was moved
		// into an assignment in the fn body
		entryFlowScope = ti.RedeclareSimpleVar(entryFlowScope, paramName, inferredType)
	}
	return entryFlowScope
}

// TypeInference_TypeDeclaringCallback abstracts logic for declaring an
// lvalue in a particular scope.
//
// It updates the given scope upon seeing an assignment or declaration:
// scope is the scope we are in, lvalue the value being updated (a NAME,
// GETPROP, GETELEM, or CAST), type_ the type we've inferred for the lvalue;
// it returns the updated flow scope.
//
// Java: TypeInference.TypeDeclaringCallback.declareTypeInScope(FlowScope, Node, JSType)
//
// Go-only: language-adapter — a single-method interface implemented only by
// lambdas, so a Go func type.
type TypeInference_TypeDeclaringCallback func(scope FlowScope, lvalue rhino.Node, type_ rhino.JSType) FlowScope

// Java: TypeInference.createInitialEstimateLattice()
func (ti *TypeInference) CreateInitialEstimateLattice() FlowScope {
	return ti.bottomScope
}

// Java: TypeInference.createEntryLattice()
func (ti *TypeInference) CreateEntryLattice() FlowScope {
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
func (ti *TypeInference) FlowThrough(n rhino.Node, input FlowScope) FlowScope {
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
	output := input.WithSyntacticScope(syntacticBlockScope)
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
		exportsVar :=
			jutil.CheckNotNullMsg1(
				syntacticBlockScope.GetVar("exports"),
				"Missing exports var for %s",
				module.Metadata())
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

// Java: TypeInference.isForward()
func (ti *TypeInference) IsForward() bool {
	return true
}

// Java: TypeInference.isBranched()
func (ti *TypeInference) IsBranched() bool {
	return true
}

// Java: TypeInference.createFlowJoiner()
func (ti *TypeInference) CreateFlowJoiner() DataFlowAnalysis_FlowJoiner[FlowScope] {
	ti.flowJoiner.result = nil
	ti.flowJoiner.inputProvider = ti.compiler
	return &ti.flowJoiner
}

// Java: TypeInference.createFlowBrancher(Node, FlowScope)
func (ti *TypeInference) CreateFlowBrancher(source rhino.Node, output FlowScope) DataFlowAnalysis_FlowBrancher[FlowScope] {
	// NOTE(nicksantos): Right now, we just treat ON_EX edges like UNCOND
	// edges. If we wanted to be perfect, we'd actually JOIN all the out
	// lattices of this flow with the in lattice, and then make that the out
	// lattice for the ON_EX edge. But it's probably too expensive to be
	// worthwhile.
	ti.branchSource = source
	ti.branchOutput = output
	ti.branchCondition = rhino.NilNode
	ti.branchConditionFlowScope = nil
	ti.branchConditionOutcomes = nil
	if ti.flowBrancher == nil {
		ti.flowBrancher = ti.branchFlow
	}
	return ti.flowBrancher
}

func (ti *TypeInference) branchFlow(branch ControlFlowGraph_Branch) FlowScope {
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
			if ti.branchConditionOutcomes == nil {
				if ti.branchCondition.IsAnd() {
					ti.branchConditionOutcomes = ti.TraverseAnd(ti.branchCondition, output)
				} else {
					ti.branchConditionOutcomes = ti.TraverseOr(ti.branchCondition, output)
				}
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
func (ti *TypeInference) Traverse(n rhino.Node, scope FlowScope) FlowScope {
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
		scope = ti.TraverseAnd(n, scope).GetJoinedFlowScope()
	case rhino.Token_OR:
		scope = ti.TraverseOr(n, scope).GetJoinedFlowScope()
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
		jutil.CheckNotNullMsg(info, "CAST node should always have JSDocInfo")
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
func (ti *TypeInference) InitializeEnhancedForScope(source rhino.Node, output FlowScope) FlowScope {
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

// Java: TypeInference.traverseUnsignedRightShift(Node, FlowScope)
func (ti *TypeInference) TraverseUnsignedRightShift(n rhino.Node, scope FlowScope) FlowScope {
	scope = ti.TraverseChildren(n, scope)
	var left, right rhino.Node
	if bin, ok := ast.AsBinary(n); ok {
		left = bin.Left().Node()
		right = bin.Right().Node()
	} else {
		aop := ast.MustAssignOp(n)
		left = aop.Target().Node()
		right = aop.Value().Node()
	}
	if TypeInference_GetBigIntPresence(ti.GetJSType(left)) != TypeInference_BigIntPresence_NO_BIGINT ||
		TypeInference_GetBigIntPresence(ti.GetJSType(right)) != TypeInference_BigIntPresence_NO_BIGINT {
		// The spec does not allow for BigInts in an unsigned right shift, so we will report an error
		n.SetJSType(ti.GetNativeType(rhino.JSTypeNative_NO_TYPE))
	} else {
		n.SetJSType(ti.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE))
	}
	return scope
}

// Java: TypeInference.traverseAssignUnsignedRightShift(Node, FlowScope)
func (ti *TypeInference) TraverseAssignUnsignedRightShift(n rhino.Node, scope FlowScope) FlowScope {
	left := ast.MustAssignOp(n).Target().Node()
	scope = ti.TraverseUnsignedRightShift(n, scope)
	return ti.UpdateScopeForAssignmentWithUpdateNode(
		scope, left, ti.GetJSType(n) /* updateNode= */, rhino.NilNode, TypeInference_AssignmentType_ASSIGN)
}

// Java: TypeInference.traverseSuper(Node, FlowScope)
func (ti *TypeInference) TraverseSuper(superNode rhino.Node, currentScope FlowScope) {
	var superNodeType rhino.ObjectType

	switch superNode.GetParent().GetToken() {
	case rhino.Token_CALL:
		// Find the closest non-arrow function (TODO(sdh): this could be an AbstractScope method).
		scope := ti.containerScope
		for scope != nil && !NodeUtil_IsNonArrowFunction(scope.GetRootNode()) {
			scope = scope.GetParent()
		}
		if scope == nil {
			superNode.SetJSType(ti.unknownType)
			return
		}
		enclosingFunctionType :=
			rhino.JSType_ToMaybeFunctionType(scope.GetRootNode().GetJSType())
		// Inside a constructor, `super` may have two different types. Calls to `super()` use the
		// super-ctor type, while property accesses use the super-instance type. `Scopes` are only
		// aware of the latter case.
		if enclosingFunctionType != nil && enclosingFunctionType.IsConstructor() {
			superNodeType = enclosingFunctionType.GetSuperClassConstructor()
		}
	case rhino.Token_GETELEM, rhino.Token_GETPROP:
		currentSyntacticScope := currentScope.GetDeclarationScope()
		// Go-only: language-adapter — StaticTypedScope.getSlot returns a StaticTypedSlot in Java
		// (a covariant return, not a cast), so a missing slot stays null and getType() on it
		// throws the NullPointerException (`super.a` outside a class); a Go type assertion on
		// the nil StaticSlot would throw a ClassCastException instead.
		var superSlot rhino.StaticTypedSlot
		if slot := currentSyntacticScope.GetSlot("super"); slot != nil {
			superSlot = slot.(rhino.StaticTypedSlot)
		}
		superNodeType = rhino.ObjectType_Cast(superSlot.GetType())
	default:
		panic(jutil.NewIllegalStateExceptionMsg(
			"Unexpected parent of SUPER: " + superNode.GetParent().ToStringTree()))
	}

	if superNodeType != nil {
		superNode.SetJSType(superNodeType)
	} else {
		superNode.SetJSType(ti.unknownType)
	}
}

// Java: TypeInference.traverseNewTarget(Node)
func (ti *TypeInference) TraverseNewTarget(newTargetNode rhino.Node) {
	// new.target is (undefined|!Function) within a vanilla function and !Function within an ES6
	// constructor.
	// Find the closest non-arrow function (TODO(sdh): this could be an AbstractScope method).
	scope := ti.containerScope
	for scope != nil && !NodeUtil_IsNonArrowFunction(scope.GetRootNode()) {
		scope = scope.GetParent()
	}
	if scope == nil {
		// `new.target` is not inside any non-arrow function. The loop above skips arrow functions
		// because they use the `new.target` of the code around them, so this also covers arrow
		// functions with no non-arrow function around them, e.g. a top-level `() => new.target`.
		// IRFactory already reports a parse error for all of these, so just use the unknown type.
		newTargetNode.SetJSType(ti.unknownType)
		return
	}
	root := scope.GetRootNode()
	parent := root.GetParent()
	if parent.GetGrandparent().IsClass() {
		// In an ES6 constuctor, new.target may not be undefined.  In any other method, it must be
		// undefined, since methods are not constructable.
		var type_ rhino.JSTypeNative
		if NodeUtil_IsEs6ConstructorMemberFunctionDef(parent) {
			type_ = rhino.JSTypeNative_FUNCTION_TYPE
		} else {
			type_ = rhino.JSTypeNative_VOID_TYPE
		}
		newTargetNode.SetJSType(ti.registry.GetNativeType(type_))
	} else {
		// Other functions also include undefined, in case they are not called with 'new'.
		newTargetNode.SetJSType(
			ti.registry.CreateUnionType(
				ti.registry.GetNativeType(rhino.JSTypeNative_FUNCTION_TYPE),
				ti.registry.GetNativeType(rhino.JSTypeNative_VOID_TYPE)))
	}
}

// TraverseReturn traverses a return value.
//
// Java: TypeInference.traverseReturn(Node, FlowScope)
func (ti *TypeInference) TraverseReturn(n rhino.Node, scope FlowScope) FlowScope {
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

// Java: TypeInference.traverseYieldAll(Node, FlowScope)
func (ti *TypeInference) TraverseYieldAll(n rhino.Node, scope FlowScope) FlowScope {
	// A yield* expression will first yield all the elements of the given iterable, and then
	// evaluate to whatever the iterable returns when done.
	// The yielded type and done return type are not necessarily the same. Here, we look for the
	// done type - TReturn in Iterable<T, TReturn, TNext>.
	scope = ti.TraverseChildren(n, scope)
	val, _ := ast.MustYield(n).Value()
	innerType := ti.GetJSType(val.Node())
	yieldAllResult := JsIterables_GetReturnElementType(innerType, ti.registry)
	n.SetJSType(yieldAllResult)
	return scope
}

// TraverseCatch: any value can be thrown, so it's really impossible to
// determine the type of a CATCH param. Treat it as the UNKNOWN type.
//
// Java: TypeInference.traverseCatch(Node, FlowScope)
func (ti *TypeInference) TraverseCatch(catchNode rhino.Node, scope FlowScope) FlowScope {
	catch := ast.MustCatch(catchNode)
	param, ok := catch.Param()
	if !ok {
		// ES2019 allows `try {} catch {}` with no catch expression
		return scope
	}
	if name, isName := param.AsName(); isName {
		nameNode := name.Node()
		type_ := nameNode.GetJSType()
		// If the catch expression name was declared in the catch in TypedScopeCreator use that type.
		// Otherwise use "unknown".
		if type_ == nil {
			type_ = ti.unknownType
			nameNode.SetJSType(ti.unknownType)
		}
		return ti.RedeclareSimpleVar(scope, nameNode, type_)
	}
	return ti.TraverseDestructuringPattern(param.Node(), scope, ti.unknownType, TypeInference_AssignmentType_DECLARATION)
}

// Java: TypeInference.traverseAssign(Node, FlowScope)
func (ti *TypeInference) TraverseAssign(n rhino.Node, scope FlowScope) FlowScope {
	assign := ast.MustAssign(n)
	target := assign.Target().Node()
	value := assign.Value().Node()
	if target.IsDestructuringPattern() {
		scope = ti.Traverse(value, scope)
		valueType := ti.GetJSType(value)
		n.SetJSType(valueType)
		return ti.TraverseDestructuringPattern(target, scope, valueType, TypeInference_AssignmentType_ASSIGN)
	} else {
		scope = ti.TraverseChildren(n, scope)

		valueType := ti.GetJSType(value)
		n.SetJSType(valueType)

		return ti.UpdateScopeForAssignment(scope, target, valueType, TypeInference_AssignmentType_ASSIGN)
	}
}

// Java: TypeInference.traverseAssignOp(Node, FlowScope)
func (ti *TypeInference) TraverseAssignOp(n rhino.Node, scope FlowScope) FlowScope {
	left := ast.MustAssignOp(n).Target().Node()
	scope = ti.TraverseBigIntCompatibleBinaryOperator(n, scope)

	// The lhs is both an input and an output, so don't update the input type here.
	return ti.UpdateScopeForAssignmentWithUpdateNode(
		scope, left, ti.GetJSType(n) /* updateNode= */, rhino.NilNode, TypeInference_AssignmentType_ASSIGN)
}

// Java: TypeInference.isInExternFile(Node)
func TypeInference_IsInExternFile(n rhino.Node) bool {
	return NodeUtil_GetSourceFile(n).IsExtern()
}

// Java: TypeInference.isPossibleMixinApplication(Node, Node)
func TypeInference_IsPossibleMixinApplication(lvalue rhino.Node, rvalue rhino.Node) bool {
	if TypeInference_IsInExternFile(lvalue) {
		return true
	}

	jsdoc := NodeUtil_GetBestJSDocInfo(lvalue)
	return jsdoc != nil &&
		jsdoc.IsConstructor() &&
		jsdoc.GetImplementedInterfaceCount() > 0 &&
		lvalue.IsQualifiedName() &&
		(!rvalue.IsNil() && rvalue.IsCall())
}

// TypeInference_AddMissingInterfaceProperties: constructor is a constructor
// function defined by a call, which may be a mixin application. The
// constructor implements at least one interface. If the constructor is
// missing some properties of the inherited interfaces, this method declares
// these properties.
//
// Java: TypeInference.addMissingInterfaceProperties(JSType)
func TypeInference_AddMissingInterfaceProperties(constructor rhino.JSType) {
	if constructor == nil || !constructor.IsConstructor() {
		return
	}
	f := constructor.ToMaybeFunctionType()
	proto := f.GetPrototype()
	for _, interf := range f.GetImplementedInterfaces() {
		for _, pname := range interf.GetPropertyNames() {
			if !proto.HasProperty(pname) {
				proto.DefineDeclaredProperty(pname, interf.GetPropertyType(pname), rhino.NilNode)
			}
		}
	}
}

// TypeInference_AssignmentType is either a combined
// declaration/initialization or a regular assignment.
type TypeInference_AssignmentType int32

const (
	TypeInference_AssignmentType_DECLARATION TypeInference_AssignmentType = iota // var x = 3;
	TypeInference_AssignmentType_ASSIGN                                          // `a.b.c = d;` or `x = 4;`
)

var typeInference_AssignmentType_names = [...]string{"DECLARATION", "ASSIGN"}

// Name is Java's Enum.name().
func (t TypeInference_AssignmentType) Name() string { return typeInference_AssignmentType_names[t] }

// String is Java's Enum.toString().
func (t TypeInference_AssignmentType) String() string { return t.Name() }

// Ordinal is Java's Enum.ordinal().
func (t TypeInference_AssignmentType) Ordinal() int { return int(t) }

// Java: TypeInference.updateScopeForAssignment(FlowScope, Node, JSType, AssignmentType)
func (ti *TypeInference) UpdateScopeForAssignment(
	scope FlowScope, target rhino.Node, resultType rhino.JSType, type_ TypeInference_AssignmentType) FlowScope {
	return ti.UpdateScopeForAssignmentWithUpdateNode(scope, target, resultType, target, type_)
}

// UpdateScopeForAssignmentWithUpdateNode updates the scope according to the
// result of an assignment.
//
// Java: TypeInference.updateScopeForAssignment(FlowScope, Node, JSType, Node, AssignmentType)
func (ti *TypeInference) UpdateScopeForAssignmentWithUpdateNode(
	scope FlowScope,
	target rhino.Node,
	resultType rhino.JSType,
	updateNode rhino.Node,
	type_ TypeInference_AssignmentType) FlowScope {
	jutil.CheckNotNull(resultType)
	jutil.CheckState(updateNode.IsNil() || updateNode == target)

	targetType := target.GetJSType() // may be null

	right := NodeUtil_GetRValueOfLValue(target)
	if TypeInference_IsPossibleMixinApplication(target, right) {
		TypeInference_AddMissingInterfaceProperties(targetType)
	}

	switch target.GetToken() {
	case rhino.Token_NAME:
		name := ast.MustName(target)
		varName := name.Ident()
		var_ := TypeInference_GetDeclaredVar(scope, varName)
		var varType rhino.JSType
		if var_ != nil {
			varType = var_.GetType()
		}
		isVarDeclaration :=
			type_ == TypeInference_AssignmentType_DECLARATION &&
				var_ != nil &&
				!var_.IsTypeInferred() &&
				!var_.GetNameNode().IsNil() // implicit vars (like arguments) have no nameNode

		// Whether this variable is declared not because it has JSDoc with a declaration, but
		// because it is const and the right-hand-side is easily inferrable.
		// e.g. these are 'typeless const declarations':
		//   const x = 0;
		//   /** @const */
		//   a.b.c = SomeOtherConstructor;
		// but these are not:
		//    let x = 0;
		//    /** @const @constructor */
		//    a.b.c = someMixin();
		// This is messy, since the definition of 'typeless const' is duplicated in
		// TypedScopeCreator and this code.
		isTypelessConstDecl :=
			isVarDeclaration &&
				var_.GetNameNode().IsName() && // ignore redeclarations of implicit globals
				NodeUtil_IsConstantDeclaration(var_.GetJSDocInfo(), var_.GetNameNode()) &&
				!(var_.GetJSDocInfo() != nil &&
					var_.GetJSDocInfo().ContainsDeclarationExcludingTypelessConst()) &&
				!(var_.GetType() != nil && var_.GetType().IsNoResolvedType())

		// When looking at VAR initializers for declared VARs, we tend
		// to use the declared type over the type it's being
		// initialized to in the global scope.
		//
		// For example,
		// /** @param {number} */ var f = goog.abstractMethod;
		// it's obvious that the programmer wants you to use
		// the declared function signature, not the inferred signature.
		//
		// Or,
		// /** @type {Object.<string>} */ var x = {};
		// the one-time anonymous object on the right side
		// is as narrow as it can possibly be, but we need to make
		// sure we back-infer the <string> element constraint on
		// the left hand side, so we use the left hand side.

		isVarTypeBetter :=
			isVarDeclaration &&
				// Makes it easier to check for NPEs.
				!resultType.IsNullType() &&
				!resultType.IsVoidType() &&
				// Do not use the var type if the declaration looked like
				// /** @const */ var x = 3;
				// because this type was computed from the RHS
				!isTypelessConstDecl

		// TODO(nicksantos): This might be a better check once we have
		// back-inference of object/array constraints.  It will probably
		// introduce more type warnings.  It uses the result type iff it's
		// strictly narrower than the declared var type.
		//
		// boolean isVarTypeBetter = isVarDeclaration &&
		//    (varType.restrictByNotNullOrUndefined().isSubtype(resultType)
		//     || !resultType.isSubtype(varType));

		if isVarTypeBetter {
			scope = ti.redeclareSimpleVar(scope, target, var_, varType)
		} else {
			scope = ti.redeclareSimpleVar(scope, target, var_, resultType)
		}

		if !updateNode.IsNil() {
			updateNode.SetJSType(resultType)
		}

		_, hasInit := name.Init()
		if var_ != nil &&
			var_.IsTypeInferred() &&
			// Don't change the typed scope to include "undefined" upon seeing "let foo;", because
			// this is incompatible with how we currently handle VARs and breaks existing code.
			// TODO(sdh): remove this condition after cleaning up code depending on it.
			!(target.GetParent().IsLet() && !hasInit) {
			oldType := var_.GetType()
			if oldType == nil {
				var_.SetType(resultType)
			} else {
				var_.SetType(oldType.GetLeastSupertype(resultType))
			}
		} else if isTypelessConstDecl {
			// /** @const */ var x = y;
			// should be redeclared, so that the type of y
			// gets propagated to inner scopes.
			var_.SetType(resultType)
		}
	case rhino.Token_GETPROP:
		gp := ast.MustGetProp(target)
		qualifiedName, isQName := target.GetQualifiedNameString()
		if isQName {
			declaredSlotType := false
			rawObjType := gp.Object().Node().GetJSType()
			if rawObjType != nil {
				objType := rhino.ObjectType_Cast(rawObjType.RestrictByNotNullOrUndefined())
				if objType != nil {
					propName := gp.Property()
					declaredSlotType = objType.IsPropertyTypeDeclared(propName)
				}
			}
			var safeLeftType rhino.JSType
			if targetType == nil {
				safeLeftType = ti.unknownType
			} else {
				safeLeftType = targetType
			}
			scope =
				scope.InferQualifiedSlot(
					target, qualifiedName, safeLeftType, resultType, declaredSlotType)
		}

		if !updateNode.IsNil() {
			updateNode.SetJSType(resultType)
		}
		ti.ensurePropertyDefined(target, qualifiedName, isQName, resultType, scope)
	default:
	}
	return scope
}

// EnsurePropertyDefined defines a property if the property has not been
// defined yet.
//
// Java: TypeInference.ensurePropertyDefined(Node, JSType, FlowScope)
func (ti *TypeInference) EnsurePropertyDefined(getprop rhino.Node, rightType rhino.JSType, scope FlowScope) {
	qName, isQName := getprop.GetQualifiedNameString()
	ti.ensurePropertyDefined(getprop, qName, isQName, rightType, scope)
}

func (ti *TypeInference) ensurePropertyDefined(getprop rhino.Node, qName string, isQName bool, rightType rhino.JSType, scope FlowScope) {
	gp := ast.MustGetProp(getprop)
	propName := gp.Property()
	obj := gp.Object().Node()
	nodeType := ti.GetJSType(obj)
	objectType := rhino.ObjectType_Cast(nodeType.RestrictByNotNullOrUndefined())
	propCreationInConstructor :=
		obj.IsThis() && ti.GetJSType(ti.containerScope.GetRootNode()).IsConstructor()

	if objectType == nil {
		ti.registry.RegisterPropertyOnType(propName, nodeType)
	} else {
		if nodeType.IsStruct() && !objectType.HasProperty(propName) {
			// In general, we don't want to define a property on a struct object,
			// b/c TypeCheck will later check for improper property creation on
			// structs. There are two exceptions.
			// 1) If it's a property created inside the constructor, on the newly
			//    created instance, allow it.
			// 2) If it's a prototype property, allow it. For example:
			//    Foo.prototype.bar = baz;
			//    where Foo.prototype is a struct and the assignment happens at the
			//    top level and the constructor Foo is defined in the same file.
			staticPropCreation := false
			maybeAssignStm := getprop.GetGrandparent()
			if ti.containerScope.IsGlobal() && NodeUtil_IsPrototypePropertyDeclaration(maybeAssignStm) {
				propCreationFilename := maybeAssignStm.GetSourceFileName()
				ctor := objectType.GetOwnerFunction().GetSource()
				if !ctor.IsNil() {
					// ctor.getSourceFileName().equals(propCreationFilename)
					ctorFileName := *ctor.GetSourceFileName()
					if propCreationFilename != nil && ctorFileName == *propCreationFilename {
						staticPropCreation = true
					}
				}
			}
			if !propCreationInConstructor && !staticPropCreation {
				return // Early return to avoid creating the property below.
			}
		}

		if ti.ensurePropertyDeclaredHelper(getprop, qName, isQName, objectType, scope) {
			return
		}

		if !objectType.IsPropertyTypeDeclared(propName) {
			// We do not want a "stray" assign to define an inferred property
			// for every object of this type in the program. So we use a heuristic
			// approach to determine whether to infer the property.
			//
			// 1) If the property is already defined, join it with the previously
			//    inferred type.
			// 2) If this isn't an instance object, define it.
			// 3) If the property of an object is being assigned in the constructor,
			//    define it.
			// 4) If this is a stub, define it.
			// 5) Otherwise, do not define the type, but declare it in the registry
			//    so that we can use it for missing property checks.
			if objectType.HasProperty(propName) || !objectType.IsInstanceType() {
				if "prototype" == propName {
					ti.DefineDeclaredProperty(objectType, propName, rightType, getprop)
				} else {
					objectType.DefineInferredProperty(propName, rightType, getprop)
				}
			} else if propCreationInConstructor {
				objectType.DefineInferredProperty(propName, rightType, getprop)
			} else {
				ti.registry.RegisterPropertyOnType(propName, objectType)
			}
		}
	}
}

// Java: TypeInference.defineDeclaredProperty(ObjectType, String, JSType, Node)
func (ti *TypeInference) DefineDeclaredProperty(
	objectType rhino.ObjectType, propName string, rightType rhino.JSType, getprop rhino.Node) bool {
	if propName == "prototype" && !getprop.GetParent().IsExprResult() {
		functionType := objectType.ToMaybeFunctionType()
		if functionType != nil &&
			!functionType.GetSource().IsNil() &&
			functionType.GetSource().IsClass() {
			ti.compiler.Report(JSError_MakeWithNode(getprop.GetParent(), TypeInference_REASSIGN_CLASS_PROTOTYPE))
			return true
		}
	}
	return objectType.DefineDeclaredProperty(propName, rightType, getprop)
}

// EnsurePropertyDeclaredHelper declares a property on its owner, if
// necessary.
//
// Returns true if a property was declared.
//
// Java: TypeInference.ensurePropertyDeclaredHelper(Node, ObjectType, FlowScope)
func (ti *TypeInference) EnsurePropertyDeclaredHelper(
	getprop rhino.Node, objectType rhino.ObjectType, scope FlowScope) bool {
	qName, isQName := getprop.GetQualifiedNameString()
	return ti.ensurePropertyDeclaredHelper(getprop, qName, isQName, objectType, scope)
}

func (ti *TypeInference) ensurePropertyDeclaredHelper(
	getprop rhino.Node, qName string, isQName bool, objectType rhino.ObjectType, scope FlowScope) bool {
	if isQName {
		propName := getprop.GetString()
		var_ := TypeInference_GetDeclaredVar(scope, qName)
		if var_ != nil && !var_.IsTypeInferred() {
			// Handle normal declarations that could not be addressed earlier.
			if propName == "prototype" ||
				// Handle prototype declarations that could not be addressed earlier.
				(!objectType.HasOwnProperty(propName) &&
					(!objectType.IsInstanceType() ||
						(var_.IsExtern() && !objectType.IsNativeObjectType()))) {
				return ti.DefineDeclaredProperty(objectType, propName, var_.GetType(), getprop)
			}
		}
	}
	return false
}

// Java: TypeInference.traverseDeclaration(Node, FlowScope)
func (ti *TypeInference) TraverseDeclaration(n rhino.Node, scope FlowScope) FlowScope {
	for declarationChild := n.GetFirstChild(); !declarationChild.IsNil(); declarationChild = declarationChild.GetNext() {
		scope = ti.TraverseDeclarationChild(declarationChild, scope)
	}

	return scope
}

// Java: TypeInference.traverseDeclarationChild(Node, FlowScope)
func (ti *TypeInference) TraverseDeclarationChild(n rhino.Node, scope FlowScope) FlowScope {
	if n.IsName() {
		return ti.TraverseName(n, scope)
	}

	lhs := ast.MustDestructuringLhs(n)
	init, _ := lhs.Init()
	scope = ti.Traverse(init.Node(), scope)
	return ti.TraverseDestructuringPattern(
		lhs.Pattern().Node(), scope, ti.GetJSType(init.Node()), TypeInference_AssignmentType_DECLARATION)
}

// TraverseDestructuringPattern traverses a destructuring pattern in an
// assignment or declaration.
//
// Java: TypeInference.traverseDestructuringPattern(Node, FlowScope, JSType, AssignmentType)
func (ti *TypeInference) TraverseDestructuringPattern(
	pattern rhino.Node, scope FlowScope, patternType rhino.JSType, assignmentType TypeInference_AssignmentType) FlowScope {
	return ti.TraverseDestructuringPatternHelper(
		pattern,
		scope,
		patternType,
		func(flowScope FlowScope, targetNode rhino.Node, targetType rhino.JSType) FlowScope {
			if targetType == nil {
				targetType = ti.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
			}
			return ti.UpdateScopeForAssignment(flowScope, targetNode, targetType, assignmentType)
		})
}

// TraverseDestructuringPatternHelper traverses a destructuring pattern, and
// calls declarer on each lvalue.
//
// The purpose of the callback is to abstract different logic for declaring
// lvalues in function parameters vs. in regular assignments/declarations.
//
// declarer contains a callback called on every lvalue.
//
// Java: TypeInference.traverseDestructuringPatternHelper(Node, FlowScope, JSType, TypeDeclaringCallback)
func (ti *TypeInference) TraverseDestructuringPatternHelper(
	pattern rhino.Node, scope FlowScope, patternType rhino.JSType, declarer TypeInference_TypeDeclaringCallback) FlowScope {
	jutil.CheckArgumentMsg1(pattern.IsDestructuringPattern(), "%s", pattern)
	jutil.CheckNotNull(patternType)
	for _, target := range DestructuredTarget_CreateAllNonEmptyTargetsInPattern(ti.registry, patternType, pattern) {

		// The computed property is always evaluated first.
		if target.HasComputedProperty() {
			scope = ti.Traverse(target.GetComputedProperty().GetFirstChild(), scope)
		}
		targetNode := target.GetNode()

		if targetNode.IsDestructuringPattern() {
			if target.HasDefaultValue() {
				ti.Traverse(target.GetDefaultValue(), scope)
			}

			// traverse into nested patterns
			targetType := target.InferType()
			if targetType == nil {
				targetType = ti.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
			}
			scope = ti.TraverseDestructuringPatternHelper(targetNode, scope, targetType, declarer)
		} else {
			scope = ti.Traverse(targetNode, scope)

			if target.HasDefaultValue() {
				// TODO(lharker): what do we do with the inferred slots in the scope?
				// throw them away or join them with the previous scope?
				ti.Traverse(target.GetDefaultValue(), scope)
			}

			// declare in the scope
			scope = declarer(scope, targetNode, target.InferType())
		}
	}
	// put the `inferred type` of a pattern on it, to make it easier to do typechecking
	pattern.SetJSType(patternType)
	return scope
}

// Java: TypeInference.traverseName(Node, FlowScope)
func (ti *TypeInference) TraverseName(n rhino.Node, scope FlowScope) FlowScope {
	name := ast.MustName(n)
	varName := name.Ident()
	type_ := n.GetJSType()
	if value, ok := name.Init(); ok {
		// The only case where `value` isn't null is when we are in a name declaration/initialization
		//     var x = 3;
		scope = ti.Traverse(value.Node(), scope)
		return ti.UpdateScopeForAssignment(scope, n, ti.GetJSType(value.Node()), TypeInference_AssignmentType_DECLARATION)
	}

	if NodeUtil_IsNameDeclaration(n.GetParent()) &&
		TypeInference_IsPossibleMixinApplication(n /* rvalue= */, rhino.NilNode) {
		TypeInference_AddMissingInterfaceProperties(type_)
	}

	if n.GetParent().IsLet() {
		// Whenever we see a LET, we're guaranteed it's not yet in the scope, and we don't need to
		// worry about it being from an outer scope.  In this case, it has no child, so the actual
		// type should be undefined, but we make a special allowance for type-annotated variables.
		// In that case, we use the annotated type instead.
		// TODO(sdh): I would have thought that #updateScopeForTypeChange would handle using the
		// declared type correctly, but for some reason it doesn't so we handle it here.
		var resultType rhino.JSType
		if type_ != nil {
			resultType = type_
		} else {
			resultType = ti.GetNativeType(rhino.JSTypeNative_VOID_TYPE)
		}
		scope = ti.UpdateScopeForAssignment(scope, n, resultType, TypeInference_AssignmentType_DECLARATION)
		type_ = resultType
	} else {
		var_, _ := scope.GetSlot(varName).(rhino.StaticTypedSlot)
		if var_ != nil {
			// There are two situations where we don't want to use type information
			// from the scope, even if we have it.

			// 1) The var is escaped and assigned in an inner scope, e.g.,
			// function f() { var x = 3; function g() { x = null } (x); }
			isInferred := var_.IsTypeInferred()
			unflowable := isInferred && ti.IsUnflowable(TypeInference_GetDeclaredVar(scope, varName))

			// 2) We're reading type information from another scope for an
			// inferred variable. That variable is assigned more than once,
			// and we can't know which type we're getting.
			//
			// var t = null; function f() { (t); } doStuff(); t = {};
			//
			// Notice that this heuristic isn't perfect. For example, you might
			// have:
			//
			// function f() { (t); } f(); var t = 3;
			//
			// In this case, we would infer the first reference to t as
			// type {number}, even though it's undefined.
			var maybeOuterVar *TypedVar
			if isInferred && ti.containerScope.IsLocal() {
				maybeOuterVar = ti.containerScope.GetParent().GetVar(varName)
			}
			// Go-only: language-adapter — Java's var.equals(maybeOuterVar) is
			// TypedVar's identity equality; a null maybeOuterVar is never equal.
			nonLocalInferredSlot :=
				maybeOuterVar != nil && var_ == rhino.StaticTypedSlot(maybeOuterVar) &&
					!maybeOuterVar.IsMarkedAssignedExactlyOnce()

			if !unflowable && !nonLocalInferredSlot {
				type_ = var_.GetType()
				if type_ == nil {
					type_ = ti.unknownType
				}
			}
		}
	}
	n.SetJSType(type_)
	return scope
}

// Java: TypeInference.traverseImportSpec(FlowScope, Node)
func (ti *TypeInference) TraverseImportSpec(scope FlowScope, spec rhino.Node) FlowScope {
	impSpec := ast.MustImportSpec(spec)
	exportedName := impSpec.Imported().Node()
	localName := impSpec.Local().Node()

	scope = ti.Traverse(localName, scope)

	exportedName.SetJSType(localName.GetJSType())
	return scope
}

// Java: TypeInference.redeclareSimpleVar(FlowScope, Node, JSType)
func (ti *TypeInference) RedeclareSimpleVar(scope FlowScope, nameNode rhino.Node, varType rhino.JSType) FlowScope {
	jutil.CheckStateMsg1(nameNode.IsName(), "%s", nameNode)
	varName := nameNode.GetString()
	return ti.redeclareSimpleVar(scope, nameNode, TypeInference_GetDeclaredVar(scope, varName), varType)
}

func (ti *TypeInference) redeclareSimpleVar(scope FlowScope, nameNode rhino.Node, declaredVar *TypedVar, varType rhino.JSType) FlowScope {
	jutil.CheckStateMsg1(nameNode.IsName(), "%s", nameNode)
	varName := nameNode.GetString()
	if varType == nil {
		varType = ti.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
	}
	if ti.IsUnflowable(declaredVar) {
		return scope
	}
	return scope.InferSlotType(varName, varType)
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
	} else {
		return jsType
	}
}

// Java: TypeInference.getNativeType(JSTypeNative)
func (ti *TypeInference) GetNativeType(typeId rhino.JSTypeNative) rhino.JSType {
	return ti.registry.GetNativeType(typeId)
}

// Java: TypeInference.getDeclaredVar(FlowScope, String)
func TypeInference_GetDeclaredVar(scope FlowScope, name string) *TypedVar {
	return scope.GetDeclarationScope().(*TypedScope).GetVar(name)
}
