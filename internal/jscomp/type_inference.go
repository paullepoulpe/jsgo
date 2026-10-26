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
	if !functionNode.IsFunction() {
		return entryFlowScope // we're in the global scope
	} else if NodeUtil_IsBundledGoogModuleCall(functionNode.GetParent()) {
		// Pretend the function literal in `goog.loadModule(function(exports) {` does not exist.
		return entryFlowScope
	}
	astParameters := functionNode.GetSecondChild()
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
	for astParamItr := astParameters.GetFirstChild(); !astParamItr.IsNil(); astParamItr = astParamItr.GetNext() {
		astParam := astParamItr
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
		if astParam.IsDefaultValue() {
			defaultValue = astParam.GetSecondChild()
			// must call `traverse` to correctly type the default value
			entryFlowScope = ti.Traverse(defaultValue, entryFlowScope)
			astParam = astParam.GetFirstChild()
		} else if astParam.IsRest() {
			// e.g. `function f(p1, ...restParamName) {}`
			// set astParam = restParamName
			astParam = astParam.GetOnlyChild()
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
		n.SetJSType(ti.GetJSType(n.GetLastChild()))
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
		if n.GetFirstChild().IsGetProp() {
			getprop := n.GetFirstChild()
			ownerType :=
				rhino.ObjectType_Cast(ti.GetJSType(getprop.GetFirstChild()).RestrictByNotNullOrUndefined())
			if ownerType != nil {
				ti.EnsurePropertyDeclaredHelper(getprop, ownerType, scope)
			}
		}
		isTypeable = false
	case rhino.Token_SWITCH:
		scope = ti.Traverse(n.GetFirstChild(), scope)
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
		if n.GetBooleanProp(rhino.Node_EXPORT_DEFAULT) {
			// TypedScopeCreator declared a dummy variable *default* to store this type. Update the
			// variable with the inferred type.
			defaultExport := TypeInference_GetDeclaredVar(scope, Export_DEFAULT_EXPORT_NAME)
			if defaultExport.IsTypeInferred() {
				defaultExport.SetType(ti.GetJSType(n.GetOnlyChild()))
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
	item := source.GetFirstChild()
	obj := item.GetNext()

	informed := ti.Traverse(obj, output)

	var assignmentType TypeInference_AssignmentType
	if NodeUtil_IsNameDeclaration(item) {
		item = item.GetFirstChild()
		assignmentType = TypeInference_AssignmentType_DECLARATION
	} else {
		assignmentType = TypeInference_AssignmentType_ASSIGN
	}
	if item.IsDestructuringLhs() {
		item = item.GetFirstChild()
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
	if TypeInference_GetBigIntPresence(ti.GetJSType(n.GetFirstChild())) != TypeInference_BigIntPresence_NO_BIGINT ||
		TypeInference_GetBigIntPresence(ti.GetJSType(n.GetLastChild())) != TypeInference_BigIntPresence_NO_BIGINT {
		// The spec does not allow for BigInts in an unsigned right shift, so we will report an error
		n.SetJSType(ti.GetNativeType(rhino.JSTypeNative_NO_TYPE))
	} else {
		n.SetJSType(ti.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE))
	}
	return scope
}

// Java: TypeInference.traverseAssignUnsignedRightShift(Node, FlowScope)
func (ti *TypeInference) TraverseAssignUnsignedRightShift(n rhino.Node, scope FlowScope) FlowScope {
	left := n.GetFirstChild()
	scope = ti.TraverseUnsignedRightShift(n, scope)
	return ti.UpdateScopeForAssignmentWithUpdateNode(
		scope, left, ti.GetJSType(n) /* updateNode= */, rhino.NilNode, TypeInference_AssignmentType_ASSIGN)
}

// Java: TypeInference.traverseCall(Node, FlowScope)
func (ti *TypeInference) TraverseCall(callNode rhino.Node, originalScope FlowScope) FlowScope {
	jutil.CheckArgumentMsg1(callNode.IsCall() || callNode.IsOptChainCall(), "%s", callNode)
	scopeAfterChildren := ti.TraverseChildren(callNode, originalScope)
	scopeAfterExecution :=
		ti.SetCallNodeTypeAfterChildrenTraversed(callNode, scopeAfterChildren)
	// e.g. after `goog.assertString(x);` we can infer `x` is a string.
	return ti.TightenTypesAfterAssertions(scopeAfterExecution, callNode)
}

// Java: TypeInference.traverseTaggedTemplateLit(Node, FlowScope)
func (ti *TypeInference) TraverseTaggedTemplateLit(tagTempLitNode rhino.Node, originalScope FlowScope) FlowScope {
	jutil.CheckArgumentMsg1(tagTempLitNode.IsTaggedTemplateLit(), "%s", tagTempLitNode)
	scopeAfterChildren := ti.TraverseChildren(tagTempLitNode, originalScope)
	// A tagged template literal is really a special kind of function call.
	scopeAfterExecution :=
		ti.SetCallNodeTypeAfterChildrenTraversed(tagTempLitNode, scopeAfterChildren)
	if tagTempLitNode.GetJSType() == nil {
		tagTempLitNode.SetJSType(ti.unknownType)
	}
	return scopeAfterExecution
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

	retValue := n.GetFirstChild()
	if !retValue.IsNil() {
		type_ := ti.containerScope.GetRootNode().GetJSType()
		if type_ != nil {
			fnType := type_.ToMaybeFunctionType()
			if fnType != nil {
				TypeInference_InferPropertyTypesToMatchConstraint(retValue.GetJSType(), fnType.GetReturnType())
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
	innerType := ti.GetJSType(n.GetFirstChild())
	yieldAllResult := JsIterables_GetReturnElementType(innerType, ti.registry)
	n.SetJSType(yieldAllResult)
	return scope
}

// TraverseCatch: any value can be thrown, so it's really impossible to
// determine the type of a CATCH param. Treat it as the UNKNOWN type.
//
// Java: TypeInference.traverseCatch(Node, FlowScope)
func (ti *TypeInference) TraverseCatch(catchNode rhino.Node, scope FlowScope) FlowScope {
	catchTarget := catchNode.GetFirstChild()
	if catchTarget.IsName() {
		name := catchNode.GetFirstChild()
		type_ := name.GetJSType()
		// If the catch expression name was declared in the catch in TypedScopeCreator use that type.
		// Otherwise use "unknown".
		if type_ == nil {
			type_ = ti.unknownType
			name.SetJSType(ti.unknownType)
		}
		return ti.RedeclareSimpleVar(scope, name, type_)
	} else if catchTarget.IsDestructuringPattern() {
		pattern := catchNode.GetFirstChild()
		return ti.TraverseDestructuringPattern(pattern, scope, ti.unknownType, TypeInference_AssignmentType_DECLARATION)
	} else {
		jutil.CheckStateMsg1(catchTarget.IsEmpty(), "%s", catchTarget)
		// ES2019 allows `try {} catch {}` with no catch expression
		return scope
	}
}

// Java: TypeInference.traverseAssign(Node, FlowScope)
func (ti *TypeInference) TraverseAssign(n rhino.Node, scope FlowScope) FlowScope {
	target := n.GetFirstChild()
	value := n.GetLastChild()
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
	left := n.GetFirstChild()
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
		varName := target.GetString()
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

		if var_ != nil &&
			var_.IsTypeInferred() &&
			// Don't change the typed scope to include "undefined" upon seeing "let foo;", because
			// this is incompatible with how we currently handle VARs and breaks existing code.
			// TODO(sdh): remove this condition after cleaning up code depending on it.
			!(target.GetParent().IsLet() && !target.HasChildren()) {
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
		qualifiedName, isQName := target.GetQualifiedNameString()
		if isQName {
			declaredSlotType := false
			rawObjType := target.GetFirstChild().GetJSType()
			if rawObjType != nil {
				objType := rhino.ObjectType_Cast(rawObjType.RestrictByNotNullOrUndefined())
				if objType != nil {
					propName := target.GetString()
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
	propName := getprop.GetString()
	obj := getprop.GetFirstChild()
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

	jutil.CheckStateMsg1(n.IsDestructuringLhs(), "%s", n)
	scope = ti.Traverse(n.GetSecondChild(), scope)
	return ti.TraverseDestructuringPattern(
		n.GetFirstChild(), scope, ti.GetJSType(n.GetSecondChild()), TypeInference_AssignmentType_DECLARATION)
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
	varName := n.GetString()
	value := n.GetFirstChild()
	type_ := n.GetJSType()
	if !value.IsNil() {
		// The only case where `value` isn't null is when we are in a name declaration/initialization
		//     var x = 3;
		scope = ti.Traverse(value, scope)
		return ti.UpdateScopeForAssignment(scope, n, ti.GetJSType(value), TypeInference_AssignmentType_DECLARATION)
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
	exportedName := spec.GetFirstChild()
	localName := spec.GetSecondChild()

	scope = ti.Traverse(localName, scope)

	exportedName.SetJSType(localName.GetJSType())
	return scope
}

// TypeInference_BigIntPresence: now that BigInt is a factor, we need a
// better understanding of which types are involved in any given operation.
// In most cases, simply knowing that BigInt is involved is not enough
// information to determine that we should report an error, for example. So,
// we have designed this enumeration (and a respective utility function
// below) to aid in providing more context for cases that are more complex.
// Currently, there are four possibilities that concern us. This is subject to
// change as we add more support for type inference with BigInts.
type TypeInference_BigIntPresence int32

const (
	TypeInference_BigIntPresence_NO_BIGINT TypeInference_BigIntPresence = iota
	TypeInference_BigIntPresence_ALL_BIGINT
	TypeInference_BigIntPresence_BIGINT_OR_NUMBER
	TypeInference_BigIntPresence_BIGINT_OR_OTHER
)

var typeInference_BigIntPresence_names = [...]string{
	"NO_BIGINT", "ALL_BIGINT", "BIGINT_OR_NUMBER", "BIGINT_OR_OTHER",
}

// Name is Java's Enum.name().
func (p TypeInference_BigIntPresence) Name() string { return typeInference_BigIntPresence_names[p] }

// String is Java's Enum.toString().
func (p TypeInference_BigIntPresence) String() string { return p.Name() }

// Ordinal is Java's Enum.ordinal().
func (p TypeInference_BigIntPresence) Ordinal() int { return int(p) }

// TypeInference_GetBigIntPresence is a utility function for determining the
// BigIntPresence for any type.
//
// Java: TypeInference.getBigIntPresence(JSType)
func TypeInference_GetBigIntPresence(type_ rhino.JSType) TypeInference_BigIntPresence {
	// Base case
	if type_.IsOnlyBigInt() {
		return TypeInference_BigIntPresence_ALL_BIGINT
	}

	// Checking enum case
	typeAsEnumElement := type_.ToMaybeEnumElementType()
	if typeAsEnumElement != nil {
		// No matter what type the enum element is, this function can resolve it to a BigIntPresence
		return TypeInference_GetBigIntPresence(typeAsEnumElement.GetPrimitiveType())
	}

	// Union case
	typeAsUnion := type_.ToMaybeUnionType()
	if typeAsUnion != nil {
		containsBigInt := false
		containsNumber := false
		containsOther := false
		for _, alternate := range typeAsUnion.GetAlternates() {
			if TypeInference_GetBigIntPresence(alternate) != TypeInference_BigIntPresence_NO_BIGINT {
				containsBigInt = true
			} else if alternate.IsNumber() && !alternate.IsUnknownType() {
				containsNumber = true
			} else {
				containsOther = true
			}
		}
		if containsBigInt {
			if containsOther {
				return TypeInference_BigIntPresence_BIGINT_OR_OTHER
			} else if containsNumber {
				return TypeInference_BigIntPresence_BIGINT_OR_NUMBER
			} else {
				return TypeInference_BigIntPresence_ALL_BIGINT
			}
		}
	}

	// If it’s not a bigint object, bigint value, enum containing either, or a union, then we can
	// safely assume that it’s not a bigint in anyway
	return TypeInference_BigIntPresence_NO_BIGINT
}

// TraverseUnaryPlus checks for BigInt with a unary plus.
//
// Java: TypeInference.traverseUnaryPlus(Node, FlowScope)
func (ti *TypeInference) TraverseUnaryPlus(n rhino.Node, scope FlowScope) FlowScope {
	scope = ti.TraverseChildren(n, scope) // Find types.
	bigintPresenceInOperand := TypeInference_GetBigIntPresence(ti.GetJSType(n.GetFirstChild()))
	if bigintPresenceInOperand != TypeInference_BigIntPresence_NO_BIGINT {
		// Unary plus throws an exception when applied to a bigint
		n.SetJSType(ti.GetNativeType(rhino.JSTypeNative_NO_TYPE))
	} else {
		n.SetJSType(ti.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE))
	}
	return scope
}

// TraverseBigIntCompatibleUnaryOperator traverses unary minus, bitwise NOT,
// increment, and decrement.
//
// Java: TypeInference.traverseBigIntCompatibleUnaryOperator(Node, FlowScope)
func (ti *TypeInference) TraverseBigIntCompatibleUnaryOperator(n rhino.Node, scope FlowScope) FlowScope {
	scope = ti.TraverseChildren(n, scope) // Find types.
	switch TypeInference_GetBigIntPresence(ti.GetJSType(n.GetFirstChild())) {
	case TypeInference_BigIntPresence_ALL_BIGINT:
		// BigIntPresence in operand
		n.SetJSType(ti.GetNativeType(rhino.JSTypeNative_BIGINT_TYPE))
	case TypeInference_BigIntPresence_NO_BIGINT:
		n.SetJSType(ti.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE))
	case TypeInference_BigIntPresence_BIGINT_OR_NUMBER, TypeInference_BigIntPresence_BIGINT_OR_OTHER:
		n.SetJSType(ti.GetNativeType(rhino.JSTypeNative_BIGINT_NUMBER))
	}
	return scope
}

// TraverseBigIntCompatibleBinaryOperator traverses all binary numeric
// operations (except for URSH and string concatenation with ADD).
//
// Java: TypeInference.traverseBigIntCompatibleBinaryOperator(Node, FlowScope)
func (ti *TypeInference) TraverseBigIntCompatibleBinaryOperator(n rhino.Node, scope FlowScope) FlowScope {
	scope = ti.TraverseChildren(n, scope) // Find types.
	leftBigIntPresence := TypeInference_GetBigIntPresence(ti.GetJSType(n.GetFirstChild()))
	rightBigIntPresence := TypeInference_GetBigIntPresence(ti.GetJSType(n.GetLastChild()))
	if leftBigIntPresence != rightBigIntPresence {
		// disallowed mixture of bigint and non-bigint operands
		n.SetJSType(ti.GetNativeType(rhino.JSTypeNative_NO_TYPE))
	} else {
		switch leftBigIntPresence {
		case TypeInference_BigIntPresence_NO_BIGINT:
			n.SetJSType(ti.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE))
		case TypeInference_BigIntPresence_ALL_BIGINT:
			n.SetJSType(ti.GetNativeType(rhino.JSTypeNative_BIGINT_TYPE))
		case TypeInference_BigIntPresence_BIGINT_OR_NUMBER:
			n.SetJSType(ti.GetNativeType(rhino.JSTypeNative_BIGINT_NUMBER))
		case TypeInference_BigIntPresence_BIGINT_OR_OTHER:
			// In the case of arithmetic operations, BigInts are only compatible with other BigInts.
			// So if bigint is in a union with anything but number (and even then they both have to
			// be
			// {bigint|number}), then an error is reported.
			n.SetJSType(ti.GetNativeType(rhino.JSTypeNative_NO_TYPE))
		}
	}
	return scope
}

// Java: TypeInference.traverseClass(Node, FlowScope)
func (ti *TypeInference) TraverseClass(n rhino.Node, scope FlowScope) FlowScope {
	// The name already has a type applied (from TypedScopeCreator) if it's non-empty, and the
	// members are traversed in the class scope (and in their own function scopes).  But the extends
	// clause and computed property keys are in the outer scope and must be traversed here.
	scope = ti.Traverse(n.GetSecondChild(), scope)
	classMembers := NodeUtil_GetClassMembers(n)

	for member := classMembers.GetFirstChild(); !member.IsNil(); member = member.GetNext() {
		// Computed properties LHS need to happen before any RHS values
		if member.IsComputedProp() || member.IsComputedFieldDef() {
			scope = ti.Traverse(member.GetFirstChild(), scope)
		}
	}
	for member := classMembers.GetFirstChild(); !member.IsNil(); member = member.GetNext() {
		scope = ti.TraverseClassMemberRhs(member, scope)
	}
	return scope
}

// Java: TypeInference.traverseClassMemberRhs(Node, FlowScope)
func (ti *TypeInference) TraverseClassMemberRhs(member rhino.Node, scope FlowScope) FlowScope {
	switch member.GetToken() {
	case rhino.Token_MEMBER_FIELD_DEF, rhino.Token_COMPUTED_FIELD_DEF:
		rhs := TypeInference_GetRhsOfField(member)
		if !rhs.IsNil() {
			computedFieldDefTypedScope := ti.scopeCreator.CreateScopeOfNode(member)
			computedFieldDefFlowScope :=
				scope.WithSyntacticScope(computedFieldDefTypedScope)
			rhsScope :=
				ti.Traverse(rhs, computedFieldDefFlowScope).
					WithSyntacticScope(scope.GetDeclarationScope())
			if member.IsStaticMember() {
				return rhsScope
			}
		}
		return scope
	case rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_COMPUTED_PROP, rhino.Token_BLOCK,
		rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF:
		return scope
	default:
		panic(jutil.NewAssertionErrorNoMessage())
	}
}

// Java: TypeInference.getRhsOfField(Node)
func TypeInference_GetRhsOfField(fieldNode rhino.Node) rhino.Node {
	switch fieldNode.GetToken() {
	case rhino.Token_MEMBER_FIELD_DEF:
		if fieldNode.HasOneChild() {
			return fieldNode.GetFirstChild()
		}
		return rhino.NilNode
	case rhino.Token_COMPUTED_FIELD_DEF:
		if fieldNode.HasTwoChildren() {
			return fieldNode.GetSecondChild()
		}
		return rhino.NilNode
	default:
		panic(jutil.NewAssertionErrorNoMessage())
	}
}

// TraverseArrayLiteral traverses each element of the array.
//
// Java: TypeInference.traverseArrayLiteral(Node, FlowScope)
func (ti *TypeInference) TraverseArrayLiteral(n rhino.Node, scope FlowScope) FlowScope {
	scope = ti.TraverseChildren(n, scope)
	if ti.arrayOfUnknownType == nil {
		ti.arrayOfUnknownType = ti.registry.CreateTemplatizedType(
			ti.registry.GetNativeObjectType(rhino.JSTypeNative_ARRAY_TYPE), ti.unknownType)
	}
	n.SetJSType(ti.arrayOfUnknownType)
	return scope
}

// Java: TypeInference.traverseObjectLiteral(Node, FlowScope)
func (ti *TypeInference) TraverseObjectLiteral(n rhino.Node, scope FlowScope) FlowScope {
	type_ := n.GetJSType()
	jutil.CheckNotNull(type_)

	spreadOperatorSeen := false
	for name := n.GetFirstChild(); !name.IsNil(); name = name.GetNext() {
		spreadOperatorSeen = spreadOperatorSeen || name.IsSpread()
		scope = ti.TraverseChildren(name, scope)
	}

	// Object literals can be reflected on other types.
	// See CodingConvention#getObjectLiteralCast and goog.reflect.object
	// Ignore these types of literals.
	objectType := rhino.ObjectType_Cast(type_)
	if objectType == nil || n.GetBooleanProp(rhino.Node_REFLECTED_OBJECT) || objectType.IsEnumType() {
		return scope
	}

	var qObjName string
	var hasQObjName bool
	if !n.GetFirstChild().IsNil() {
		qObjName, hasQObjName = ti.getBestLValueNameString(NodeUtil_GetBestLValue(n))
	}
	for key := n.GetFirstChild(); !key.IsNil(); key = key.GetNext() {
		if key.IsComputedProp() {
			// Don't define computed properties as inferred properties on the object
			continue
		}

		if key.IsSpread() {
			// TODO(b/128355893): Do smarter inferrence. There are a lot of potential issues with
			// inference on object-spread, so for now we just give up and say `Object`.
			n.SetJSType(ti.registry.GetNativeType(rhino.JSTypeNative_OBJECT_TYPE))
			break
		}

		// Java tests `memberName != null` here and otherwise sets the unknown type
		// on n; NodeUtil.getObjectOrClassLitKeyName never returns null (it throws).
		memberName := NodeUtil_GetObjectOrClassLitKeyName(key)
		rawValueType := key.GetFirstChild().GetJSType()
		valueType := TypeCheck_GetObjectLitKeyTypeFromValueType(key, rawValueType)
		if valueType == nil {
			valueType = ti.unknownType
		}

		// See b/260837012. traverseObjectLiteral() can be invoked multiple times for the same
		// literal. When a spread operator is encountered, the type of the literal is changed to
		// OBJECT_TYPE. A second invocation of this code may then erroneously define new properties
		// on the global OBJECT_TYPE. So, if there are spread operators in the literal and the type
		// has already been propagated to OBJECT_TYPE, do not add newly inferred properties to the
		// type.
		if !spreadOperatorSeen || !n.GetJSType().IsNativeObjectType() {
			objectType.DefineInferredProperty(memberName, valueType, key)
		}

		// Do normal flow inference if this is a direct property assignment.
		if hasQObjName && key.IsStringKey() {
			qKeyName := qObjName + "." + memberName
			var_ := TypeInference_GetDeclaredVar(scope, qKeyName)
			var oldType rhino.JSType
			if var_ != nil {
				oldType = var_.GetType()
			}
			if var_ != nil && var_.IsTypeInferred() {
				if oldType == nil {
					var_.SetType(valueType)
				} else {
					var_.SetType(oldType.GetLeastSupertype(oldType))
				}
			}

			var bottomType rhino.JSType
			if oldType == nil {
				bottomType = ti.unknownType
			} else {
				bottomType = oldType
			}
			scope =
				scope.InferQualifiedSlot(
					key, qKeyName, bottomType, valueType, false)
		}
	}
	return scope
}

// Java: TypeInference.traverseAdd(Node, FlowScope)
func (ti *TypeInference) TraverseAdd(n rhino.Node, scope FlowScope) FlowScope {
	left := n.GetFirstChild()
	right := left.GetNext()
	scope = ti.TraverseChildren(n, scope)

	leftType := left.GetJSType()
	rightType := right.GetJSType()

	var type_ rhino.JSType = ti.unknownType
	if leftType != nil && rightType != nil {
		leftIsUnknown := leftType.IsUnknownType()
		rightIsUnknown := rightType.IsUnknownType()
		if (!leftIsUnknown && leftType.IsString()) || (!rightIsUnknown && rightType.IsString()) {
			type_ = ti.GetNativeType(rhino.JSTypeNative_STRING_TYPE)
		} else if TypeInference_GetBigIntPresence(leftType) != TypeInference_BigIntPresence_NO_BIGINT ||
			TypeInference_GetBigIntPresence(rightType) != TypeInference_BigIntPresence_NO_BIGINT {
			return ti.TraverseBigIntCompatibleBinaryOperator(n, scope)
		} else if leftIsUnknown || rightIsUnknown {
			type_ = ti.unknownType
		} else if ti.IsAddedAsNumber(leftType) && ti.IsAddedAsNumber(rightType) {
			type_ = ti.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)
		} else {
			type_ = ti.GetNativeType(rhino.JSTypeNative_NUMBER_STRING)
		}
	}
	n.SetJSType(type_)

	if n.IsAssignAdd() {
		scope =
			ti.UpdateScopeForAssignmentWithUpdateNode(
				scope, left, type_, rhino.NilNode /* updateNode= */, TypeInference_AssignmentType_ASSIGN)
	}

	return scope
}

// Java: TypeInference.isAddedAsNumber(JSType)
func (ti *TypeInference) IsAddedAsNumber(type_ rhino.JSType) bool {
	return type_.IsSubtypeOf(ti.numberAdditionSupertype)
}

// Java: TypeInference.traverseHook(Node, FlowScope)
func (ti *TypeInference) TraverseHook(n rhino.Node, scope FlowScope) FlowScope {
	condition := n.GetFirstChild()
	trueNode := condition.GetNext()
	falseNode := n.GetLastChild()

	// verify the condition
	scope = ti.Traverse(condition, scope)

	// reverse abstract interpret the condition to produce two new scopes
	trueScope :=
		ti.reverseInterpreter.GetPreciserScopeKnowingConditionOutcome(condition, scope, rhino.Outcome_TRUE)
	falseScope :=
		ti.reverseInterpreter.GetPreciserScopeKnowingConditionOutcome(condition, scope, rhino.Outcome_FALSE)

	// traverse the true node with the trueScope
	ti.Traverse(trueNode, trueScope)

	// traverse the false node with the falseScope
	ti.Traverse(falseNode, falseScope)

	// meet true and false nodes' types and assign
	trueType := trueNode.GetJSType()
	falseType := falseNode.GetJSType()
	if trueType != nil && falseType != nil {
		n.SetJSType(trueType.GetLeastSupertype(falseType))
	} else {
		n.SetJSType(ti.unknownType)
	}

	return scope
}

// Java: TypeInference.traverseNullishCoalesce(Node, FlowScope)
func (ti *TypeInference) TraverseNullishCoalesce(n rhino.Node, scope FlowScope) FlowScope {
	jutil.CheckArgument(n.IsNullishCoalesce() || n.IsAssignNullishCoalesce())
	left := n.GetFirstChild()
	right := n.GetLastChild()

	scope = ti.Traverse(left, scope)

	rightScope :=
		ti.reverseInterpreter.GetPreciserScopeKnowingConditionOutcome(left, scope, rhino.Outcome_NULLISH)

	scopeAfterTraverseRight := ti.Traverse(right, rightScope)

	leftType := left.GetJSType()
	rightType := right.GetJSType()

	var type_ rhino.JSType = ti.unknownType
	if leftType != nil {
		if !leftType.IsNullable() && !leftType.IsVoidable() {
			type_ = leftType
		} else if rightType != nil {
			type_ = ti.registry.CreateUnionType(leftType.RestrictByNotNullOrUndefined(), rightType)
			scope = ti.Join(scope, scopeAfterTraverseRight)
			// Assignment occurs if lhs is null
			if n.IsAssignNullishCoalesce() {
				scope = ti.UpdateScopeForAssignment(scope, left, type_, TypeInference_AssignmentType_ASSIGN)
			}
		}
	}
	n.SetJSType(type_)
	return scope
}

// SetCallNodeTypeAfterChildrenTraversed: n is a non-constructor function
// invocation, i.e. CALL or TAGGED_TEMPLATELIT.
//
// Java: TypeInference.setCallNodeTypeAfterChildrenTraversed(Node, FlowScope)
func (ti *TypeInference) SetCallNodeTypeAfterChildrenTraversed(n rhino.Node, scopeAfterChildren FlowScope) FlowScope {
	// Resolve goog.{require,requireType,requireDynamic,forwardDeclare,module.get} calls separately,
	// as they are not normal functions.
	if n.IsCall() &&
		NodeUtil_IsExpressionResultUsed(n) && // Don't bother typing calls if the result is unused.
		ModuleImportResolver_IsGoogModuleDependencyCall(n) {
		calleeNode := n.GetFirstChild()
		if calleeNode.MatchesQualifiedName(TypeInference_GOOG_REQUIREDYNAMIC_NAME) {
			// Result of `goog.requireDynamic('module.name')` is `IThenable<MODULE_TYPE>`
			// IThenable is a supertype of Promise.
			n.SetJSType(Promises_WrapInIThenable(ti.registry, ti.GetGoogModuleDependencyCallResultType(n)))
		} else {
			n.SetJSType(ti.GetGoogModuleDependencyCallResultType(n))
		}
		return scopeAfterChildren
	}

	left := n.GetFirstChild()
	functionType := ti.GetJSType(left).RestrictByNotNullOrUndefined()
	if left.IsSuper() {
		// TODO(sdh): This will probably return the super type; might want to return 'this' instead?
		return ti.TraverseInstantiation(n, functionType, scopeAfterChildren)
	} else if functionType.IsFunctionType() {
		fnType := functionType.ToMaybeFunctionType()
		n.SetJSType(fnType.GetReturnType())
		ti.BackwardsInferenceFromCallSite(n, fnType, scopeAfterChildren)
	} else if functionType.Equals(ti.GetNativeType(rhino.JSTypeNative_CHECKED_UNKNOWN_TYPE)) {
		n.SetJSType(ti.GetNativeType(rhino.JSTypeNative_CHECKED_UNKNOWN_TYPE))
	} else if left.GetJSType() != nil && left.GetJSType().IsUnknownType() {
		// TODO(lharker): do we also want to set this to unknown if the left's type is null? We would
		// lose some inference that TypeCheck does when given a null type.
		n.SetJSType(ti.unknownType)
	}
	return scopeAfterChildren
}

// Java: TypeInference.getGoogModuleDependencyCallResultType(Node)
func (ti *TypeInference) GetGoogModuleDependencyCallResultType(callNode rhino.Node) rhino.JSType {
	moduleId := callNode.GetSecondChild().GetString()
	module := ti.compiler.GetModuleMap().GetClosureModule(moduleId)
	// Fall back to the `?` type if the module is unknown to the compiler
	if module == nil {
		return ti.unknownType
	}

	name := ti.moduleImportResolver.GetClosureNamespaceTypeFromCall(callNode)
	if name != nil {
		otherModuleScope := ti.scopeCreator.GetNodeToScopeMapper()(name.GetScopeRoot())
		var otherVar *TypedVar
		if otherModuleScope != nil {
			otherVar, _ = otherModuleScope.GetSlot(name.GetName()).(*TypedVar)
		}
		if otherVar != nil {
			if otherVar.GetType() != nil {
				return otherVar.GetType()
			}
			return ti.unknownType
		}

		if otherModuleScope != nil && module.Metadata().ModuleType() == ModuleMetadataMap_ModuleType_GOOG_PROVIDE {
			// A "provideAlreadyProvided" is type through the properties of the containing namespace
			// since it is a provided name it is rooted in the global scope.

			// We validated that this is a valid "provided name" above so it is ok to look it up by
			// properties.
			return otherModuleScope.GetTypeThroughNamespace(moduleId)
		}
	}

	return ti.unknownType
}

// Java: TypeInference.tightenTypesAfterAssertions(FlowScope, Node)
func (ti *TypeInference) TightenTypesAfterAssertions(scope FlowScope, callNode rhino.Node) FlowScope {
	left := callNode.GetFirstChild()
	firstParam := left.GetNext()
	if firstParam.IsNil() {
		// this may be an assertion call but there are no arguments to assert
		return scope
	}
	assertionFunctionSpec := ti.assertionFunctionLookup.LookupByCallee(left)
	if assertionFunctionSpec == nil {
		// this is not a recognized assertion function
		return scope
	}
	assertedNode := assertionFunctionSpec.GetAssertedArg(firstParam)
	if assertedNode.IsNil() {
		return scope
	}

	// Handle assertions that enforce expressions evaluate to true.
	switch assertionFunctionSpec.GetAssertionKind() {
	case CodingConvention_AssertionFunctionSpec_AssertionKind_TRUTHY:
		// Handle arbitrary expressions within the assert.
		// e.g. given `assert(typeof x === 'string')`, the resulting scope will infer x to be a
		// string.
		scope =
			ti.reverseInterpreter.GetPreciserScopeKnowingConditionOutcome(
				assertedNode, scope, rhino.Outcome_TRUE)
		// Build the result of the assertExpression
		truthyType := ti.GetJSType(assertedNode).RestrictByNotNullOrUndefined()
		callNode.SetJSType(truthyType)
	case CodingConvention_AssertionFunctionSpec_AssertionKind_MATCHES_RETURN_TYPE:
		// Handle assertions that enforce expressions match the return type of the function
		callType := rhino.JSType_ToMaybeFunctionType(left.GetJSType())
		var assertedType rhino.JSType
		if callType != nil {
			assertedType = callType.GetReturnType()
		} else {
			assertedType = ti.unknownType
		}
		type_ := ti.GetJSType(assertedNode)
		var narrowed rhino.JSType
		if assertedType.IsUnknownType() || type_.IsUnknownType() {
			narrowed = assertedType
		} else {
			narrowed = type_.GetGreatestSubtype(assertedType)
		}
		callNode.SetJSType(narrowed)
		if assertedNode.IsQualifiedName() && type_.DiffersFrom(narrowed) {
			scope = ti.NarrowScope(scope, assertedNode, narrowed)
		}
	}

	return scope
}

// Java: TypeInference.narrowScope(FlowScope, Node, JSType)
func (ti *TypeInference) NarrowScope(scope FlowScope, node rhino.Node, narrowed rhino.JSType) FlowScope {
	if node.IsThis() {
		// "this" references don't need to be modeled in the control flow graph.
		return scope
	}

	if node.IsGetProp() {
		qName, _ := node.GetQualifiedNameString()
		return scope.InferQualifiedSlot(
			node, qName, ti.GetJSType(node), narrowed, false)
	}
	return ti.RedeclareSimpleVar(scope, node, narrowed)
}

// BackwardsInferenceFromCallSite: we only do forward type inference. We do
// not do full backwards type inference.
//
// In other words, if we have,
//
//	var x = f();
//	g(x);
//
// a forward type-inference engine would try to figure out the type of "x"
// from the return type of "f". A backwards type-inference engine would try to
// figure out the type of "x" from the parameter type of "g".
//
// However, there are a few special syntactic forms where we do some some
// half-assed backwards type-inference, because programmers expect it in this
// day and age. To take an example from Java,
//
//	List<String> x = Lists.newArrayList();
//
// The Java compiler will be able to infer the generic type of the List
// returned by newArrayList().
//
// In much the same way, we do some special-case backwards inference for JS.
// Those cases are enumerated here.
//
// Java: TypeInference.backwardsInferenceFromCallSite(Node, FunctionType, FlowScope)
func (ti *TypeInference) BackwardsInferenceFromCallSite(n rhino.Node, fnType rhino.FunctionType, scope FlowScope) {
	updatedFnType := ti.InferTemplatedTypesForCall(n, fnType, scope)
	if updatedFnType {
		fnType = n.GetFirstChild().GetJSType().ToMaybeFunctionType()
	}
	ti.UpdateTypeOfArguments(n, fnType)
	ti.UpdateBind(n)
}

// UpdateBind: when "bind" is called on a function, we infer the type of the
// returned "bound" function by looking at the number of parameters in the
// call site. We also infer the "this" type of the target, if it's a function
// expression.
//
// Java: TypeInference.updateBind(Node)
func (ti *TypeInference) UpdateBind(n rhino.Node) {
	bind := ti.compiler.GetCodingConvention().DescribeFunctionBind(n, true)
	if bind == nil {
		return
	}

	target := bind.target
	callTargetFn :=
		ti.GetJSType(target).RestrictByNotNullOrUndefined().ToMaybeFunctionType()
	if callTargetFn == nil {
		return
	}

	if !bind.thisValue.IsNil() && target.IsFunction() {
		thisType := ti.GetJSType(bind.thisValue)
		if thisType.ToObjectType() != nil &&
			!thisType.IsUnknownType() &&
			callTargetFn.GetTypeOfThis().IsUnknownType() {
			callTargetFn =
				callTargetFn.ToBuilder().
					WithTypeOfThis(thisType.ToObjectType()).
					WithSourceNode(rhino.NilNode).
					BuildAndResolve()
			target.SetJSType(callTargetFn)
		}
	}

	returnType :=
		callTargetFn.GetBindReturnType(
			// getBindReturnType expects the 'this' argument to be included.
			bind.GetBoundParameterCount() + 1)
	bindType := ti.GetJSType(n.GetFirstChild()).ToMaybeFunctionType()
	if bindType != nil && n.GetFirstChild() != target {
		// Update the type of the
		bindType =
			bindType.ToBuilder().WithReturnType(returnType).WithSourceNode(rhino.NilNode).BuildAndResolve()
		n.GetFirstChild().SetJSType(bindType)
	}

	n.SetJSType(returnType)
}

// UpdateTypeOfArguments performs a limited back-inference on function
// arguments based on the expected parameter types.
//
// Currently this only does back-inference in two cases: it infers the type of
// function literal arguments and adds inferred properties to inferred
// object-typed arguments.
//
// For example: if someone calls `Promise<string>.prototype.then` with
// `(result) => ...` then we infer that the type of the arrow function is
// `function(string): ?`, and inside the arrow function body we know that
// `result` is a string.
//
// Java: TypeInference.updateTypeOfArguments(Node, FunctionType)
func (ti *TypeInference) UpdateTypeOfArguments(n rhino.Node, fnType rhino.FunctionType) {
	jutil.CheckStateMsg1(NodeUtil_IsInvocation(n), "%s", n)
	parameters := fnType.GetParameters()
	nextParameter := 0
	var curArg rhino.Node
	isTaggedTemplateLit := n.IsTaggedTemplateLit()
	if isTaggedTemplateLit {
		// Skip the first parameter because it corresponds to a constructed array of the template lit
		// subs, not an actual AST node, so there's nothing to update.
		if nextParameter >= len(parameters) {
			// TypeCheck will warn if there is no first parameter. Just bail out here.
			return
		}
		nextParameter++
		for curArg = n.GetLastChild().GetFirstChild(); !curArg.IsNil() && !curArg.IsTemplateLitSub(); curArg = curArg.GetNext() {
		}
	} else {
		curArg = n.GetSecondChild()
	}

	// Note: if there are too many or too few arguments, TypeCheck will warn.
	for nextParameter < len(parameters) && !curArg.IsNil() {
		iArgument := curArg
		curArg = curArg.GetNext()
		if isTaggedTemplateLit {
			iArgument = iArgument.GetFirstChild()
			for !curArg.IsNil() && !curArg.IsTemplateLitSub() {
				curArg = curArg.GetNext()
			}
		}
		iArgumentType := ti.GetJSType(iArgument)

		iParameter := parameters[nextParameter]
		nextParameter++
		var iParameterType rhino.JSType
		if iParameter.GetJSType() != nil {
			iParameterType = iParameter.GetJSType()
		} else {
			iParameterType = ti.unknownType
		}

		TypeInference_InferPropertyTypesToMatchConstraint(iArgumentType, iParameterType)

		// If the parameter to the call is a function expression, propagate the
		// function signature from the call site to the function node.

		// Filter out non-function types (such as null and undefined) as
		// we only care about FUNCTION subtypes here.
		var restrictedParameter rhino.FunctionType
		if iParameterType.IsUnionType() {
			union := iParameterType.ToMaybeUnionType()
			for _, alternative := range union.GetAlternates() {
				if alternative.IsFunctionType() {
					// There is only one function type per union.
					restrictedParameter = alternative.ToMaybeFunctionType()
					break
				}
			}
		} else {
			restrictedParameter = iParameterType.ToMaybeFunctionType()
		}

		if restrictedParameter != nil && iArgument.IsFunction() && iArgumentType.IsFunctionType() {
			argFnType := iArgumentType.ToMaybeFunctionType()
			argJsdoc := iArgument.GetJSDocInfo()
			// Treat the parameter & return types of the function as 'declared' if the function has
			// JSDoc with type annotations, or a parameter has inline JSDoc.
			// Note that this does not distinguish between cases where all parameters have JSDoc vs
			// only one parameter has JSDoc.
			declared :=
				(argJsdoc != nil && argJsdoc.ContainsDeclaration()) ||
					NodeUtil_FunctionHasInlineJsdocs(iArgument)
			iArgument.SetJSType(ti.MatchFunction(restrictedParameter, argFnType, declared))
		}
	}
}

// MatchFunction takes the current function type, and tries to match the
// expected function type. This is a form of backwards-inference, like
// record-type constraint matching.
//
// declared is whether the given function type is user-provided as opposed to
// inferred.
//
// Java: TypeInference.matchFunction(FunctionType, FunctionType, boolean)
func (ti *TypeInference) MatchFunction(
	expectedType rhino.FunctionType, currentType rhino.FunctionType, declared bool) rhino.FunctionType {
	if declared {
		// If the function was declared but it doesn't have a known "this"
		// but the expected type does, back fill it.
		if currentType.GetTypeOfThis().IsUnknownType() &&
			!expectedType.GetTypeOfThis().IsUnknownType() {
			return currentType.ToBuilder().
				WithTypeOfThis(expectedType.GetTypeOfThis()).
				BuildAndResolve()
		}
	} else {
		// For now, we just make sure the current type has enough
		// arguments to match the expected type, and return the
		// expected type if it does.
		if currentType.GetMaxArity() <= expectedType.GetMaxArity() {
			return expectedType
		}
	}
	return currentType
}

// BuildTypeVariables builds the type environment where type transformations
// will be evaluated. It only considers the template type variables that do
// not have a type transformation.
//
// Java: TypeInference.buildTypeVariables(Map<TemplateType, JSType>)
func (ti *TypeInference) BuildTypeVariables(
	inferredTypes *jutil.LinkedHashMap[*rhino.TemplateType, rhino.JSType]) *jutil.LinkedHashMap[string, rhino.JSType] {
	typeVars := jutil.NewLinkedHashMap[string, rhino.JSType]()
	for key, value := range inferredTypes.All() {
		// Only add the template type that do not have a type transformation
		if !key.IsTypeTransformation() {
			typeVars.Put(*key.GetReferenceName(), value)
		}
	}
	return typeVars
}

// EvaluateTypeTransformations evaluates the type transformations associated
// to the template types.
//
// Java: TypeInference.evaluateTypeTransformations(ImmutableList<TemplateType>, Map<TemplateType, JSType>, FlowScope)
func (ti *TypeInference) EvaluateTypeTransformations(
	templateTypes []*rhino.TemplateType,
	inferredTypes *jutil.LinkedHashMap[*rhino.TemplateType, rhino.JSType],
	scope FlowScope) *jutil.LinkedHashMap[*rhino.TemplateType, rhino.JSType] {

	var typeVars *jutil.LinkedHashMap[string, rhino.JSType]
	var result *jutil.LinkedHashMap[*rhino.TemplateType, rhino.JSType]
	var ttlObj *TypeTransformation

	for _, type_ := range templateTypes {
		if type_.IsTypeTransformation() {
			// Lazy initialization when the first type transformation is found
			if ttlObj == nil {
				ttlObj = NewTypeTransformation(ti.compiler, scope.GetDeclarationScope())
				typeVars = ti.BuildTypeVariables(inferredTypes)
				result = jutil.NewLinkedHashMap[*rhino.TemplateType, rhino.JSType]()
			}
			// Evaluate the type transformation expression using the current
			// known types for the template type variables
			transformedType :=
				ttlObj.Eval(type_.GetTypeTransformation(), jutil.NewLinkedHashMapFrom(typeVars))
			result.Put(type_, transformedType)
			// Add the transformed type to the type variables
			typeVars.Put(*type_.GetReferenceName(), transformedType)
		}
	}
	return result
}

func (ti *TypeInference) matchInvocationTemplateTypes(
	calleeType rhino.FunctionType, localThisType rhino.JSType, invocation rhino.Node) *jutil.LinkedHashMap[*rhino.TemplateType, rhino.JSType] {
	ti.templatedMatched.Clear()
	if calleeType.GetTemplateTypeMap().IsEmpty() {
		return &ti.templatedMatched
	}
	matcher := InvocationTemplateTypeMatcher{
		matchedTypes:  &ti.templatedMatched,
		registry:      ti.registry,
		calleeType:    calleeType,
		localThisType: localThisType,
		invocation:    invocation,
		unknownType:   ti.unknownType,
	}
	return matcher.Match()
}

// InferTemplatedTypesForCall: for functions that use template types,
// specialize the function type for the call target based on the call-site
// specific arguments. Specifically, this enables inference to set the type
// of any function literal parameters based on these inferred types.
//
// Java: TypeInference.inferTemplatedTypesForCall(Node, FunctionType, FlowScope)
func (ti *TypeInference) InferTemplatedTypesForCall(n rhino.Node, fnType rhino.FunctionType, scope FlowScope) bool {
	keys := fnType.GetTemplateTypeMap().GetTemplateKeys()
	if len(keys) == 0 {
		return false
	}

	// Try to infer the template types
	bindings := ti.matchInvocationTemplateTypes(fnType, scope.GetTypeOfThis(), n)
	inferred := &ti.templatedInferred
	inferred.Clear()
	// If the inferred type doesn't satisfy the template bound, swap to using the bound. This
	// ensures errors will be reported in type-checking.
	for _, key := range keys {
		v := bindings.GetOrDefault(key, ti.unknownType)
		if !v.IsSubtypeOf(key.GetBound()) {
			v = key.GetBound()
		}
		inferred.Put(key, v)
	}

	// Try to infer the template types using the type transformations
	typeTransformations :=
		ti.EvaluateTypeTransformations(keys, inferred, scope)
	if typeTransformations != nil {
		for k, v := range typeTransformations.All() {
			inferred.Put(k, v)
		}
	}

	// Replace all template types. If we couldn't find a replacement, we
	// replace it with UNKNOWN.
	replacer := rhino.TemplateTypeReplacer_MakeForInference(ti.registry, inferred)
	callTarget := n.GetFirstChild()

	replacementFnType := replacer.VisitJSType(fnType).ToMaybeFunctionType()
	jutil.CheckNotNull(replacementFnType)
	callTarget.SetJSType(replacementFnType)
	n.SetJSType(replacementFnType.GetReturnType())

	return replacer.HasMadeReplacement()
}

// Java: TypeInference.traverseNew(Node, FlowScope)
func (ti *TypeInference) TraverseNew(n rhino.Node, scope FlowScope) FlowScope {
	scope = ti.TraverseChildren(n, scope)
	constructor := n.GetFirstChild()
	constructorType := constructor.GetJSType()
	return ti.TraverseInstantiation(n, constructorType, scope)
}

// Java: TypeInference.traverseInstantiation(Node, JSType, FlowScope)
func (ti *TypeInference) TraverseInstantiation(n rhino.Node, ctorType rhino.JSType, scope FlowScope) FlowScope {
	if ctorType == nil || ctorType.IsUnknownType() {
		n.SetJSType(ti.unknownType)
		return scope
	}

	ctorType = ctorType.RestrictByNotNullOrUndefined()

	ctorFnType := ctorType.ToMaybeFunctionType()
	if functionType, ok := ctorType.(rhino.FunctionType); ctorFnType == nil && ok {
		// If ctorType is a NoObjectType, then toMaybeFunctionType will
		// return null. But NoObjectType implements the FunctionType
		// interface, precisely because it can validly construct objects.
		ctorFnType = functionType
	}

	if ctorFnType == nil || !ctorFnType.IsConstructor() {
		n.SetJSType(ti.unknownType)
		return scope
	}

	// TODO(nickreid): This probably isn't the right thing to do based on the differences between
	// the `this` parameter between CALL and NEW.
	ti.BackwardsInferenceFromCallSite(n, ctorFnType, scope)

	instantiatedType := ctorFnType.GetInstanceType()
	if rhino.JSType(ctorFnType) == ti.registry.GetNativeType(rhino.JSTypeNative_OBJECT_FUNCTION_TYPE) {
		// TODO(b/138617950): Delete this case when `Object` and `Object<?, ?> are sparate.
	} else if ctorFnType.HasAnyTemplateTypes() {
		if instantiatedType.IsTemplatizedType() {
			instantiatedType = instantiatedType.ToMaybeTemplatizedType().GetRawType()
		}
		// If necessary, templatized the instance type based on the the constructor parameters.
		inferredTypes := ti.matchInvocationTemplateTypes(ctorFnType, scope.GetTypeOfThis(), n)
		instantiatedType =
			ti.registry.CreateTemplatizedTypeOfMap(instantiatedType, inferredTypes).ToMaybeObjectType()
	}

	if instantiatedType != nil {
		n.SetJSType(instantiatedType)
	} else {
		n.SetJSType(ti.unknownType)
	}
	return scope
}

// Java: TypeInference.traverseAnd(Node, FlowScope)
func (ti *TypeInference) TraverseAnd(n rhino.Node, scope FlowScope) *TypeInference_BooleanOutcomePair {
	left := n.GetFirstChild()
	leftOutcome := ti.TraverseWithinShortCircuitingBinOp(left, scope)
	return ti.TraverseShortCircuitingBinOp(n, left, leftOutcome)
}

// Java: TypeInference.traverseChildren(Node, FlowScope)
func (ti *TypeInference) TraverseChildren(n rhino.Node, scope FlowScope) FlowScope {
	for el := n.GetFirstChild(); !el.IsNil(); el = el.GetNext() {
		scope = ti.Traverse(el, scope)
	}
	return scope
}

// Java: TypeInference.traverseGetElem(Node, FlowScope)
func (ti *TypeInference) TraverseGetElem(n rhino.Node, originalScope FlowScope) FlowScope {
	executedScope := ti.TraverseChildren(n, originalScope)
	return ti.SetGetElemNodeTypeAfterChildrenTraversed(n, executedScope)
}

// SetGetElemNodeTypeAfterChildrenTraversed sets the appropriate type on the
// GETELEM node `n` after its children have been traversed.
//
// n is the node that we want to finishTraversing; scopeAfterChildren the
// scope after children are traversed.
//
// Java: TypeInference.setGetElemNodeTypeAfterChildrenTraversed(Node, FlowScope)
func (ti *TypeInference) SetGetElemNodeTypeAfterChildrenTraversed(n rhino.Node, scopeAfterChildren FlowScope) FlowScope {
	jutil.CheckArgument(n.IsGetElem() || n.IsOptChainGetElem())
	ti.InferGetElemType(n)
	scopeAfterChildren = ti.TightenTypeAfterDereference(n.GetFirstChild(), scopeAfterChildren)
	return scopeAfterChildren
}

// Java: TypeInference.inferGetElemType(Node)
func (ti *TypeInference) InferGetElemType(n rhino.Node) {
	objType := ti.GetJSType(n.GetFirstChild()).RestrictByNotNullOrUndefined()
	indexKey := n.GetLastChild()
	indexType := ti.GetJSType(indexKey)

	var inferredType rhino.JSType
	if indexType.IsKnownSymbolValueType() {
		inferredType = ti.DereferenceKnownSymbolProp(objType, indexType.ToMaybeKnownSymbolType())
	} else if indexType.IsSymbolValueType() {
		// For now, allow symbols definitions/access on any type. In the future only allow them
		// on the subtypes for which they are defined.
		// TODO(b/77474174): be stricter about accesses for non-well-known symbols
		inferredType = ti.unknownType
	} else {
		inferredType = ti.DereferenceIndexSignature(objType)
	}
	if inferredType != nil {
		n.SetJSType(inferredType)
	} else {
		n.SetJSType(ti.unknownType)
	}
}

// Java: TypeInference.dereferenceKnownSymbolProp(JSType, KnownSymbolType)
func (ti *TypeInference) DereferenceKnownSymbolProp(obj rhino.JSType, indexSymbol *rhino.KnownSymbolType) rhino.JSType {
	if !obj.IsUnionType() {
		objType := obj.ToMaybeObjectType()
		if objType == nil {
			return nil
		}
		propertyType := objType.FindPropertyTypeByKey(rhino.NewProperty_SymbolKey(indexSymbol))
		if propertyType == nil {
			return nil
		}
		if !propertyType.IsUnionType() && !propertyType.IsNoType() {
			return propertyType
		}
		argumentTypes := rhino.UnionType_Builder_(ti.registry)
		argumentTypes.AddAlternate(propertyType)
		return argumentTypes.Build()
	}

	// If this is a union type, then we must extract type arguments from each option.
	argumentTypes := rhino.UnionType_Builder_(ti.registry)
	alternates := obj.ToMaybeUnionType().GetAlternates()
	key := rhino.NewProperty_SymbolKey(indexSymbol)
	for _, option := range alternates {
		if option.ToMaybeObjectType() == nil {
			// This isn't an array or object, drop out.
			argumentTypes = nil
			break
		}

		// Extract the element type and add all options to our set of alternates.
		propertyType := option.ToMaybeObjectType().FindPropertyTypeByKey(key)
		if propertyType == nil {
			// this union member doesn't have the property. just make it unknown.
			argumentTypes = nil
			break
		}
		argumentTypes.AddAlternate(propertyType)
	}

	// Unwrap the union if possible, and fail if we had no alternates.
	if argumentTypes == nil {
		return nil
	}
	return argumentTypes.Build()
}

// Java: TypeInference.dereferenceIndexSignature(JSType)
func (ti *TypeInference) DereferenceIndexSignature(obj rhino.JSType) rhino.JSType {
	if !obj.IsUnionType() {
		typeMap := obj.GetTemplateTypeMap()
		if !typeMap.HasTemplateType(ti.registry.GetObjectElementKey()) {
			return nil
		}
		resolved := typeMap.GetResolvedTemplateType(ti.registry.GetObjectElementKey())
		if !resolved.IsUnionType() && !resolved.IsNoType() {
			return resolved
		}
		argumentTypes := rhino.UnionType_Builder_(ti.registry)
		argumentTypes.AddAlternate(resolved)
		return argumentTypes.Build()
	}

	// If this is a union type, then we must extract type arguments from each option.
	argumentTypes := rhino.UnionType_Builder_(ti.registry)
	alternates := obj.ToMaybeUnionType().GetAlternates()
	for _, option := range alternates {
		typeMap := option.GetTemplateTypeMap()
		if !typeMap.HasTemplateType(ti.registry.GetObjectElementKey()) {
			// This isn't an array or object, drop out.
			argumentTypes = nil
			break
		}

		// Extract the element type and add all options to our set of alternates.
		argumentTypes.AddAlternate(typeMap.GetResolvedTemplateType(ti.registry.GetObjectElementKey()))
	}

	// Unwrap the union if possible, and fail if we had no alternates.
	if argumentTypes == nil {
		return nil
	}
	return argumentTypes.Build()
}

// Java: TypeInference.traverseGetProp(Node, FlowScope)
func (ti *TypeInference) TraverseGetProp(n rhino.Node, scope FlowScope) FlowScope {
	scope = ti.TraverseChildren(n, scope)
	return ti.SetGetPropNodeTypeAfterChildrenTraversed(n, scope)
}

// SetGetPropNodeTypeAfterChildrenTraversed sets the appropriate type on the
// GETPROP node `n` after its children have been traversed.
//
// Java: TypeInference.setGetPropNodeTypeAfterChildrenTraversed(Node, FlowScope)
func (ti *TypeInference) SetGetPropNodeTypeAfterChildrenTraversed(n rhino.Node, scopeAfterChildren FlowScope) FlowScope {
	jutil.CheckArgument(n.IsGetProp() || n.IsOptChainGetProp())
	objNode := n.GetFirstChild()
	n.SetJSType(ti.GetPropertyType(objNode.GetJSType(), n.GetString(), n, scopeAfterChildren))
	return ti.TightenTypeAfterDereference(n.GetFirstChild(), scopeAfterChildren)
}

// SetOptChainNodeTypeAfterChildrenTraversed sets the appropriate type on the
// OptChain node `n` after its children have been traversed.
//
// Java: TypeInference.setOptChainNodeTypeAfterChildrenTraversed(Node, FlowScope)
func (ti *TypeInference) SetOptChainNodeTypeAfterChildrenTraversed(
	n rhino.Node, scopeAfterChildren FlowScope) FlowScope {
	switch n.GetToken() {
	case rhino.Token_OPTCHAIN_GETPROP:
		return ti.SetGetPropNodeTypeAfterChildrenTraversed(n, scopeAfterChildren)
	case rhino.Token_OPTCHAIN_GETELEM:
		return ti.SetGetElemNodeTypeAfterChildrenTraversed(n, scopeAfterChildren)
	case rhino.Token_OPTCHAIN_CALL:
		return ti.SetCallNodeTypeAfterChildrenTraversed(n, scopeAfterChildren)
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Illegal token inside finishTraversingOptChain"))
	}
}

// TypeInference_OptChainInfo holds flow scopes for conditional and
// unconditional parts during traversal of an optional chain.
type TypeInference_OptChainInfo struct {
	endOfChain         rhino.Node
	startOfChain       rhino.Node
	unconditionalScope FlowScope
}

// Java: TypeInference.OptChainInfo(Node, Node)
func NewTypeInference_OptChainInfo(endOfChain rhino.Node, startOfChain rhino.Node) *TypeInference_OptChainInfo {
	o := &TypeInference_OptChainInfo{}
	o.endOfChain = endOfChain
	o.startOfChain = startOfChain
	o.unconditionalScope = nil
	return o
}

// TraverseOptChain: traversal requires holding scopes from the
// unconditional start (lhs of the `?.`) and the conditional part (rhs of the
// `?.`), and using the startType to determine the resulting scope.
//
// Java: TypeInference.traverseOptChain(Node, FlowScope)
func (ti *TypeInference) TraverseOptChain(n rhino.Node, scope FlowScope) FlowScope {
	jutil.CheckArgument(NodeUtil_IsOptChainNode(n))

	if NodeUtil_IsEndOfOptChainSegment(n) {
		// Create new optional chain tracking object and push it onto the stack.
		startOfChain := NodeUtil_GetStartOfOptChainSegment(n)
		optChainInfo := NewTypeInference_OptChainInfo(n, startOfChain)
		// optChainArrayDeque.addFirst(optChainInfo)
		ti.optChainArrayDeque = append(ti.optChainArrayDeque, optChainInfo)
	}

	lhsScope := ti.Traverse(n.GetFirstChild(), scope)

	if n.IsOptionalChainStart() {
		// Store lhsScope into top-of-stack as unexecuted (unconditional) scope.
		ti.optChainArrayDequePeekFirst().unconditionalScope = lhsScope
	}

	// Traverse the remaining children and capture their changes into a new FlowScope var
	// `aboutToExecuteScope`. This FlowScope must be constructed on top of the `lhsScope` and not
	// the original scope `scope`, otherwise changes to outer variable that are preserved in the
	// lhsScope would not be captured in the `aboutToExecuteScope`.
	aboutToExecuteScope := lhsScope
	nextChild := n.GetSecondChild()
	for !nextChild.IsNil() {
		aboutToExecuteScope = ti.Traverse(nextChild, aboutToExecuteScope)
		nextChild = nextChild.GetNext()
	}

	// Assigns the type to `n` assuming the entire chain executes and returns the the executed
	// scope.
	executedScope := ti.SetOptChainNodeTypeAfterChildrenTraversed(n, aboutToExecuteScope)

	// Unlike CALL, the OPTCHAIN_CALL nodes must not remain untyped when left child is untyped.
	if n.GetJSType() == nil {
		n.SetJSType(ti.unknownType)
	}

	if NodeUtil_IsEndOfOptChainSegment(n) {
		// Use the startNode's type to selectively join the executed scope with the unexecuted scope,
		// and update the type assigned to `n` in `setXAfterChildrenTraversed()`
		startOfChain := NodeUtil_GetStartOfOptChainSegment(n)

		// Pop the stack to obtain the current chain.
		currentChain := ti.optChainArrayDequeRemoveFirst()

		// Sanity check that the popped chain correctly corresponds to the current optional chain
		jutil.CheckState(currentChain.endOfChain == n)
		jutil.CheckState(currentChain.startOfChain == startOfChain)

		startNode := rhino.CheckNotNull(startOfChain.GetFirstChild())
		return ti.UpdateTypeWhenEndOfOptChain(n, startNode, currentChain, executedScope)
	} else {
		// `n` is just an inner node (i.e. not the end of the current chain). Simply return the
		// executedScope.
		return executedScope
	}
}

// optChainArrayDequePeekFirst is ArrayDeque.peekFirst(): null when empty.
//
// Go-only: language-adapter — see optChainArrayDeque.
func (ti *TypeInference) optChainArrayDequePeekFirst() *TypeInference_OptChainInfo {
	if len(ti.optChainArrayDeque) == 0 {
		return nil
	}
	return ti.optChainArrayDeque[len(ti.optChainArrayDeque)-1]
}

// optChainArrayDequeRemoveFirst is ArrayDeque.removeFirst(): throws
// NoSuchElementException when empty.
//
// Go-only: language-adapter — see optChainArrayDeque.
func (ti *TypeInference) optChainArrayDequeRemoveFirst() *TypeInference_OptChainInfo {
	if len(ti.optChainArrayDeque) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	first := ti.optChainArrayDeque[len(ti.optChainArrayDeque)-1]
	ti.optChainArrayDeque[len(ti.optChainArrayDeque)-1] = nil
	ti.optChainArrayDeque = ti.optChainArrayDeque[:len(ti.optChainArrayDeque)-1]
	return first
}

// UpdateTypeWhenEndOfOptChain uses the startNode's type to selectively join
// the executed scope with the unexecuted scope, and updates the type assigned
// to `optChain` in `setXAfterChildrenTraversed()`.
//
// Java: TypeInference.updateTypeWhenEndOfOptChain(Node, Node, OptChainInfo, FlowScope)
func (ti *TypeInference) UpdateTypeWhenEndOfOptChain(
	optChain rhino.Node, startNode rhino.Node, currentChain *TypeInference_OptChainInfo, executedScope FlowScope) FlowScope {

	startType := ti.GetJSType(startNode)
	if startType.IsUnknownType() {
		// Unknown startType: Conditional part may execute. Result type must be unknown.
		optChain.SetJSType(ti.unknownType)
		return ti.Join(executedScope, currentChain.unconditionalScope)
	} else if startType.IsNullType() || startType.IsVoidType() {
		// Conditional part will not execute.
		optChain.SetJSType(ti.registry.GetNativeType(rhino.JSTypeNative_VOID_TYPE))
		return currentChain.unconditionalScope
	} else if !startType.IsNullable() {
		// Conditional part will execute; `optChain` was assigned the right type in
		// `setXAfterChildrenTraversed()`.
		return executedScope
	} else {
		// Nullable start: Conditional part may execute. Need to add a VOID_TYPE to the type
		// assigned to `optChain` within `setXAfterChildrenTraversed()`.
		optChain.SetJSType(
			ti.registry.CreateUnionType(ti.registry.GetNativeType(rhino.JSTypeNative_VOID_TYPE), optChain.GetJSType()))
		unexecutedScope :=
			ti.reverseInterpreter.GetPreciserScopeKnowingConditionOutcome(
				startNode, currentChain.unconditionalScope, rhino.Outcome_NULLISH)
		return ti.Join(executedScope, unexecutedScope)
	}
}

// TypeInference_InferPropertyTypesToMatchConstraint: suppose X is an object
// with inferred properties. Suppose also that X is used in a way where it
// would only type-check correctly if some of those properties are widened.
// Then we should be polite and automatically widen X's properties.
//
// For a concrete example, consider: param x {{prop: (number|undefined)}}
// function f(x) {} f({});
//
// If we give the anonymous object an inferred property of
// (number|undefined), then this code will type-check appropriately.
//
// Java: TypeInference.inferPropertyTypesToMatchConstraint(JSType, JSType)
func TypeInference_InferPropertyTypesToMatchConstraint(type_ rhino.JSType, constraint rhino.JSType) {
	if type_ == nil || constraint == nil {
		return
	}

	type_.MatchConstraint(constraint)
}

// TightenTypeAfterDereference: if we access a property of a symbol, then that
// symbol is not null or undefined.
//
// Java: TypeInference.tightenTypeAfterDereference(Node, FlowScope)
func (ti *TypeInference) TightenTypeAfterDereference(n rhino.Node, scope FlowScope) FlowScope {
	switch n.GetToken() {
	case rhino.Token_NAME:
		if n.GetString() == "" {
			return scope
		}
		type_ := ti.GetJSType(n)
		narrowed := type_.RestrictByNotNullOrUndefined()
		if type_ != narrowed && !type_.Equals(narrowed) {
			scope = ti.RedeclareSimpleVar(scope, n, narrowed)
		}
	case rhino.Token_GETPROP:
		type_ := ti.GetJSType(n)
		narrowed := type_.RestrictByNotNullOrUndefined()
		if type_ != narrowed && !type_.Equals(narrowed) {
			if qName, ok := n.GetQualifiedNameString(); ok {
				scope = scope.InferQualifiedSlot(n, qName, type_, narrowed, false)
			}
		}
	case rhino.Token_SUPER:
		type_ := ti.GetJSType(n)
		narrowed := type_.RestrictByNotNullOrUndefined()
		if type_ != narrowed && !type_.Equals(narrowed) {
			scope = ti.RedeclareSimpleVar(scope, n, narrowed)
		}
	}
	return scope
}

func (ti *TypeInference) getBestLValueNameString(lValue rhino.Node) (string, bool) {
	if lValue.IsNil() || lValue.GetParent().IsObjectLit() {
		return "", false
	}
	if NodeUtil_IsGetOrSetKey(lValue) || lValue.IsComputedProp() {
		return "", false
	}
	if lValue.IsStringKey() || lValue.IsMemberFunctionDef() || lValue.IsMemberFieldDef() {
		objectLit := lValue.GetParent()
		owner, ok := ti.getBestLValueNameString(NodeUtil_GetBestLValue(objectLit))
		if ok {
			return owner + "." + NodeUtil_GetObjectOrClassLitKeyName(lValue), true
		}
		return "", false
	}
	return lValue.GetQualifiedNameString()
}

// Java: TypeInference.getPropertyType(JSType, String, Node, FlowScope)
func (ti *TypeInference) GetPropertyType(objType rhino.JSType, propName string, n rhino.Node, scope FlowScope) rhino.JSType {
	// We often have a couple of different types to choose from for the
	// property. Ordered by accuracy, we have
	// 1) A locally inferred qualified name (which is in the FlowScope)
	// 2) A globally declared qualified name (which is in the FlowScope)
	// 3) A property on the owner type (which is on objType)
	// 4) A name in the type registry (as a last resort)
	var propertyType rhino.JSType
	isLocallyInferred := false

	// Scopes sometimes contain inferred type info about qualified names.
	qualifiedName, hasQName := n.GetQualifiedNameString()
	var var_ rhino.StaticTypedSlot
	if hasQName {
		var_, _ = scope.GetSlot(qualifiedName).(rhino.StaticTypedSlot)
	}
	if var_ != nil {
		varType := var_.GetType()
		if varType != nil {
			isDeclared := !var_.IsTypeInferred()
			if !isDeclared || varType.Equals(ti.unknownType) {
				declaredVar := TypeInference_GetDeclaredVar(scope, qualifiedName)
				// Go-only: language-adapter — Java's `var != getDeclaredVar(...)` is an
				// identity comparison; a null declared var differs from var.
				isLocallyInferred = declaredVar == nil || var_ != rhino.StaticTypedSlot(declaredVar)
			}
			if isDeclared || isLocallyInferred {
				propertyType = varType
			}
		}
	}

	if propertyType == nil && objType != nil {
		foundType := objType.FindPropertyType(propName)
		if foundType != nil {
			propertyType = foundType
		}
	}

	if (propertyType == nil || propertyType.IsUnknownType()) && hasQName {
		// If we find this node in the registry, then we can infer its type.
		regType :=
			rhino.ObjectType_Cast(ti.registry.GetType(scope.GetDeclarationScope(), qualifiedName))
		if regType != nil {
			propertyType = regType.GetConstructor()
		}
	}

	if propertyType == nil {
		return ti.unknownType
	} else if propertyType.Equals(ti.unknownType) && isLocallyInferred {
		// If the type has been checked in this scope,
		// then use CHECKED_UNKNOWN_TYPE instead to indicate that.
		return ti.GetNativeType(rhino.JSTypeNative_CHECKED_UNKNOWN_TYPE)
	} else {
		return propertyType
	}
}

// Java: TypeInference.traverseOr(Node, FlowScope)
func (ti *TypeInference) TraverseOr(n rhino.Node, scope FlowScope) *TypeInference_BooleanOutcomePair {
	left := n.GetFirstChild()
	leftOutcome := ti.TraverseWithinShortCircuitingBinOp(left, scope)
	return ti.TraverseShortCircuitingBinOp(n, left, leftOutcome)
}

// Java: TypeInference.traverseShortCircuitingBinOpAssignment(Node, FlowScope)
func (ti *TypeInference) TraverseShortCircuitingBinOpAssignment(n rhino.Node, scope FlowScope) FlowScope {
	left := n.GetFirstChild()
	nIsAnd := n.IsAssignAnd()
	leftOutcome := ti.TraverseWithinShortCircuitingBinOp(left, scope)
	outcome := ti.TraverseShortCircuitingBinOp(n, left, leftOutcome)

	outcomeJoinedFlowScope := outcome.GetJoinedFlowScope()

	if leftOutcome.toBooleanOutcomes == rhino.BooleanLiteralSet_Get(!nIsAnd) {
		// Either n is && and lhs has a toBooleanOutcome of false,
		// or n is || and lhs has a toBooleanOutcome of true, so assignment does not occur
		// The scope is otherwise updated according to the result of an assignment.
		return outcomeJoinedFlowScope
	}
	return ti.UpdateScopeForAssignment(
		outcomeJoinedFlowScope, n.GetFirstChild(), n.GetJSType(), TypeInference_AssignmentType_ASSIGN)
}

// Java: TypeInference.traverseShortCircuitingBinOp(Node, Node, BooleanOutcomePair)
func (ti *TypeInference) TraverseShortCircuitingBinOp(
	n rhino.Node, left rhino.Node, leftOutcome *TypeInference_BooleanOutcomePair) *TypeInference_BooleanOutcomePair {
	jutil.CheckArgument(n.IsAnd() || n.IsOr() || n.IsAssignAnd() || n.IsAssignOr())
	nIsAnd := n.IsAnd() || n.IsAssignAnd()
	right := n.GetLastChild()

	// type the left node
	leftType := left.GetJSType()

	// reverse abstract interpret the left node to produce the correct
	// scope in which to verify the right node
	rightScope :=
		ti.reverseInterpreter.GetPreciserScopeKnowingConditionOutcome(
			left,
			leftOutcome.GetOutcomeFlowScope(left.GetToken(), nIsAnd),
			rhino.Outcome_ForBoolean(nIsAnd))

	// type the right node
	rightOutcome := ti.TraverseWithinShortCircuitingBinOp(right, rightScope)
	rightType := right.GetJSType()

	var type_ rhino.JSType
	var outcome *TypeInference_BooleanOutcomePair
	if leftType != nil && rightType != nil {
		leftType = leftType.GetRestrictedTypeGivenOutcome(rhino.Outcome_ForBoolean(!nIsAnd))
		if leftOutcome.toBooleanOutcomes == rhino.BooleanLiteralSet_Get(!nIsAnd) {
			// Either n is && and lhs is false, or n is || and lhs is true.
			// Use the restricted left type; the right side never gets evaluated.
			type_ = leftType
			outcome = leftOutcome
		} else {
			// Use the join of the restricted left type knowing the outcome of the
			// ToBoolean predicate and of the right type.
			type_ = leftType.GetLeastSupertype(rightType)
			outcome =
				NewTypeInference_BooleanOutcomePair(
					ti,
					TypeInference_JoinBooleanOutcomes(
						nIsAnd, leftOutcome.toBooleanOutcomes, rightOutcome.toBooleanOutcomes),
					TypeInference_JoinBooleanOutcomes(nIsAnd, leftOutcome.booleanValues, rightOutcome.booleanValues),
					leftOutcome.GetJoinedFlowScope(),
					rightOutcome.GetJoinedFlowScope())
		}
		// Exclude the boolean type if the literal set is empty because a boolean
		// can never actually be returned.
		if outcome.booleanValues == rhino.BooleanLiteralSet_EMPTY &&
			ti.GetNativeType(rhino.JSTypeNative_BOOLEAN_TYPE).IsSubtypeOf(type_) {
			// Exclusion only makes sense for a union type.
			if type_.IsUnionType() {
				type_ = type_.ToMaybeUnionType().GetRestrictedUnion(ti.GetNativeType(rhino.JSTypeNative_BOOLEAN_TYPE))
			}
		}
	} else {
		type_ = ti.unknownType
		outcome =
			NewTypeInference_BooleanOutcomePair(
				ti,
				rhino.BooleanLiteralSet_BOTH,
				rhino.BooleanLiteralSet_BOTH,
				leftOutcome.GetJoinedFlowScope(),
				rightOutcome.GetJoinedFlowScope())
	}
	n.SetJSType(type_)
	return outcome
}

// Java: TypeInference.traverseWithinShortCircuitingBinOp(Node, FlowScope)
func (ti *TypeInference) TraverseWithinShortCircuitingBinOp(n rhino.Node, scope FlowScope) *TypeInference_BooleanOutcomePair {
	switch n.GetToken() {
	case rhino.Token_AND:
		return ti.TraverseAnd(n, scope)
	case rhino.Token_OR:
		return ti.TraverseOr(n, scope)
	default:
		scope = ti.Traverse(n, scope)
		return ti.NewBooleanOutcomePair(n.GetJSType(), scope)
	}
}

// Java: TypeInference.traverseAwait(Node, FlowScope)
func (ti *TypeInference) TraverseAwait(await rhino.Node, scope FlowScope) FlowScope {
	scope = ti.TraverseChildren(await, scope)

	expr := await.GetFirstChild()
	exprType := ti.GetJSType(expr)
	await.SetJSType(Promises_GetResolvedType(ti.registry, exprType))

	return scope
}

// Java: TypeInference.traverseDynamicImport(Node, FlowScope)
func (ti *TypeInference) TraverseDynamicImport(dynamicImport rhino.Node, scope FlowScope) FlowScope {
	templateType := ti.registry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)

	// If the module specifier is a string, attempt to resolve the module
	moduleMap := ti.compiler.GetModuleMap()
	importSpecifier := dynamicImport.GetFirstChild()
	input := ti.compiler.GetInput(NodeUtil_GetInputId(dynamicImport))
	if importSpecifier.IsStringLit() && moduleMap != nil && input != nil {
		targetPath :=
			input.
				GetPath().
				ResolveJsModule(
					importSpecifier.GetString(),
					*importSpecifier.GetSourceFileName(),
					importSpecifier.GetLineno(),
					importSpecifier.GetCharno())
		var targetModule *Module
		if targetPath != nil {
			targetModule = moduleMap.GetModuleByPath(targetPath)
		}
		if targetModule != nil {
			// TypedScopeCreator ensures that the MODULE_BODY type is the export namespace type
			scriptNode := targetModule.Metadata().RootNode()
			if scriptNode.HasOneChild() && scriptNode.GetFirstChild().IsModuleBody() {
				var exportNamespaceType rhino.JSType
				if !scriptNode.IsNil() {
					exportNamespaceType = scriptNode.GetOnlyChild().GetJSType()
				}
				if exportNamespaceType != nil {
					templateType = exportNamespaceType
				}
			} else {
				// Module transpilation has occurred before type inference so a MODULE_BODY node
				// no longer exists. Traverse the script to locate the module namespace variable
				// and retrieve the type from it.
				moduleName :=
					NodeUtil_FindPreorder(
						scriptNode,
						func(node rhino.Node) bool { return node.MatchesQualifiedName(targetPath.ToModuleName()) },
						func(rhino.Node) bool { return true })
				if !moduleName.IsNil() && moduleName.GetJSType() != nil {
					templateType = moduleName.GetJSType()
				}
			}
		}
	}

	dynamicImport.SetJSType(
		ti.registry.CreateTemplatizedType(
			ti.registry.GetNativeObjectType(rhino.JSTypeNative_PROMISE_TYPE), templateType))

	return ti.TraverseChildren(dynamicImport, scope)
}

// Java: TypeInference.joinBooleanOutcomes(boolean, BooleanLiteralSet, BooleanLiteralSet)
func TypeInference_JoinBooleanOutcomes(
	isAnd bool, left rhino.BooleanLiteralSet, right rhino.BooleanLiteralSet) rhino.BooleanLiteralSet {
	// A truthy value on the lhs of an {@code &&} can never make it to the
	// result. Same for a falsy value on the lhs of an {@code ||}.
	// Hence the intersection.
	return right.Union(left.Intersection(rhino.BooleanLiteralSet_Get(!isAnd)))
}

// TypeInference_BooleanOutcomePair: when traversing short-circuiting binary
// operations, we need to keep track of two sets of boolean literals: 1.
// toBooleanOutcomes: boolean literals as converted from any types, 2.
// booleanValues: boolean literals from just boolean types.
type TypeInference_BooleanOutcomePair struct {
	outer *TypeInference

	toBooleanOutcomes rhino.BooleanLiteralSet
	booleanValues     rhino.BooleanLiteralSet

	// The scope if only half of the expression executed, when applicable.
	leftScope FlowScope

	// The scope when the whole expression executed.
	rightScope FlowScope

	// The scope when we don't know how much of the expression is executed.
	joinedScope FlowScope
}

// Java: TypeInference.BooleanOutcomePair(BooleanLiteralSet, BooleanLiteralSet, FlowScope, FlowScope)
func NewTypeInference_BooleanOutcomePair(
	outer *TypeInference,
	toBooleanOutcomes rhino.BooleanLiteralSet,
	booleanValues rhino.BooleanLiteralSet,
	leftScope FlowScope,
	rightScope FlowScope) *TypeInference_BooleanOutcomePair {
	return &TypeInference_BooleanOutcomePair{
		outer:             outer,
		toBooleanOutcomes: toBooleanOutcomes,
		booleanValues:     booleanValues,
		leftScope:         leftScope,
		rightScope:        rightScope,
	}
}

// GetJoinedFlowScope gets the safe estimated scope without knowing if all of
// the subexpressions will be evaluated.
//
// Java: TypeInference.BooleanOutcomePair.getJoinedFlowScope()
func (p *TypeInference_BooleanOutcomePair) GetJoinedFlowScope() FlowScope {
	if p.joinedScope == nil {
		if p.leftScope == p.rightScope {
			p.joinedScope = p.rightScope
		} else {
			p.joinedScope = p.outer.Join(p.leftScope, p.rightScope)
		}
	}
	return p.joinedScope
}

// GetOutcomeFlowScope gets the outcome scope if we do know the outcome of
// the entire expression.
//
// Java: TypeInference.BooleanOutcomePair.getOutcomeFlowScope(Token, boolean)
func (p *TypeInference_BooleanOutcomePair) GetOutcomeFlowScope(nodeType rhino.Token, outcome bool) FlowScope {
	if (nodeType == rhino.Token_AND && outcome) || (nodeType == rhino.Token_OR && !outcome) {
		// We know that the whole expression must have executed.
		return p.rightScope
	} else {
		return p.GetJoinedFlowScope()
	}
}

// Java: TypeInference.newBooleanOutcomePair(JSType, FlowScope)
func (ti *TypeInference) NewBooleanOutcomePair(jsType rhino.JSType, flowScope FlowScope) *TypeInference_BooleanOutcomePair {
	if jsType == nil {
		return NewTypeInference_BooleanOutcomePair(
			ti, rhino.BooleanLiteralSet_BOTH, rhino.BooleanLiteralSet_BOTH, flowScope, flowScope)
	}
	var booleanValues rhino.BooleanLiteralSet
	if ti.registry.GetNativeType(rhino.JSTypeNative_BOOLEAN_TYPE).IsSubtypeOf(jsType) {
		booleanValues = rhino.BooleanLiteralSet_BOTH
	} else {
		booleanValues = rhino.BooleanLiteralSet_EMPTY
	}
	return NewTypeInference_BooleanOutcomePair(
		ti,
		jsType.GetPossibleToBooleanOutcomes(),
		booleanValues,
		flowScope,
		flowScope)
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
