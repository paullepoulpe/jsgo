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

// Java: TypeInference.traverseCall(Node, FlowScope)
func (ti *TypeInference) TraverseCall(callNode rhino.Node, originalScope *LinkedFlowScope) *LinkedFlowScope {
	jutil.CheckArgumentMsg1(callNode.IsCall() || callNode.IsOptChainCall(), "%s", callNode)
	scopeAfterChildren := ti.TraverseChildren(callNode, originalScope)
	scopeAfterExecution :=
		ti.SetCallNodeTypeAfterChildrenTraversed(callNode, scopeAfterChildren)
	// e.g. after `goog.assertString(x);` we can infer `x` is a string.
	return ti.TightenTypesAfterAssertions(scopeAfterExecution, callNode)
}

// Java: TypeInference.traverseTaggedTemplateLit(Node, FlowScope)
func (ti *TypeInference) TraverseTaggedTemplateLit(tagTempLitNode rhino.Node, originalScope *LinkedFlowScope) *LinkedFlowScope {
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

// SetCallNodeTypeAfterChildrenTraversed: n is a non-constructor function
// invocation, i.e. CALL or TAGGED_TEMPLATELIT.
//
// Java: TypeInference.setCallNodeTypeAfterChildrenTraversed(Node, FlowScope)
func (ti *TypeInference) SetCallNodeTypeAfterChildrenTraversed(n rhino.Node, scopeAfterChildren *LinkedFlowScope) *LinkedFlowScope {
	inv := ast.MustInvocation(n)
	left := inv.Callee().Node()
	// Resolve goog.{require,requireType,requireDynamic,forwardDeclare,module.get} calls separately,
	// as they are not normal functions.
	if n.IsCall() &&
		NodeUtil_IsExpressionResultUsed(n) && // Don't bother typing calls if the result is unused.
		ModuleImportResolver_IsGoogModuleDependencyCall(n) {
		if left.MatchesQualifiedName(TypeInference_GOOG_REQUIREDYNAMIC_NAME) {
			// Result of `goog.requireDynamic('module.name')` is `IThenable<MODULE_TYPE>`
			// IThenable is a supertype of Promise.
			n.SetJSType(Promises_WrapInIThenable(ti.registry, ti.GetGoogModuleDependencyCallResultType(n)))
		} else {
			n.SetJSType(ti.GetGoogModuleDependencyCallResultType(n))
		}
		return scopeAfterChildren
	}

	functionType := ti.GetJSType(left).RestrictByNotNullOrUndefined()
	if left.IsSuper() {
		// TODO(sdh): This will probably return the super type; might want to return 'this' instead?
		return ti.TraverseInstantiation(n, functionType, scopeAfterChildren)
	} else if functionType.IsFunctionType() {
		fnType := functionType.ToMaybeFunctionType()
		n.SetJSType(fnType.GetReturnType())
		ti.BackwardsInferenceFromCallSite(n, fnType, scopeAfterChildren)
	} else if checkedUnknown := ti.GetNativeType(rhino.JSTypeNative_CHECKED_UNKNOWN_TYPE); functionType.Equals(checkedUnknown) {
		n.SetJSType(checkedUnknown)
	} else if left.GetJSType() != nil && left.GetJSType().IsUnknownType() {
		// TODO(lharker): do we also want to set this to unknown if the left's type is null? We would
		// lose some inference that TypeCheck does when given a null type.
		n.SetJSType(ti.unknownType)
	}
	return scopeAfterChildren
}

// Java: TypeInference.getGoogModuleDependencyCallResultType(Node)
func (ti *TypeInference) GetGoogModuleDependencyCallResultType(callNode rhino.Node) rhino.JSType {
	firstArg, _ := ast.MustCall(callNode).FirstArg()
	moduleId := firstArg.Node().GetString()
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
			otherVar = otherModuleScope.GetVar(name.GetName())
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
func (ti *TypeInference) TightenTypesAfterAssertions(scope *LinkedFlowScope, callNode rhino.Node) *LinkedFlowScope {
	call := ast.MustNormalOrOptChainCall(callNode)
	left := call.Callee().Node()
	firstParamExpr, ok := call.FirstArg()
	if !ok {
		// this may be an assertion call but there are no arguments to assert
		return scope
	}
	firstParam := firstParamExpr.Node()
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
func (ti *TypeInference) NarrowScope(scope *LinkedFlowScope, node rhino.Node, narrowed rhino.JSType) *LinkedFlowScope {
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
func (ti *TypeInference) BackwardsInferenceFromCallSite(n rhino.Node, fnType rhino.FunctionType, scope *LinkedFlowScope) {
	updatedFnType := ti.InferTemplatedTypesForCall(n, fnType, scope)
	if updatedFnType {
		fnType = ast.MustInvocation(n).Callee().Node().GetJSType().ToMaybeFunctionType()
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
	callee := ast.MustInvocation(n).Callee().Node()
	bindType := ti.GetJSType(callee).ToMaybeFunctionType()
	if bindType != nil && callee != target {
		// Update the type of the
		bindType =
			bindType.ToBuilder().WithReturnType(returnType).WithSourceNode(rhino.NilNode).BuildAndResolve()
		callee.SetJSType(bindType)
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

// evaluateTypeTransformations evaluates the type transformations associated
// with the template types and updates values in place.
//
// Java: TypeInference.evaluateTypeTransformations(ImmutableList<TemplateType>, Map<TemplateType, JSType>, FlowScope)
func (ti *TypeInference) evaluateTypeTransformations(
	templateTypes []*rhino.TemplateType,
	values []rhino.JSType,
	scope *LinkedFlowScope) {
	var typeVarBuf [8]ttlTypeVar
	typeVars := typeVarBuf[:0]
	var ttlObj TypeTransformation
	initialized := false

	for idx, type_ := range templateTypes {
		if !type_.IsTypeTransformation() {
			continue
		}
		// Lazy initialization when the first type transformation is found
		if !initialized {
			initialized = true
			ttlObj = TypeTransformation{
				compiler: ti.compiler,
				registry: ti.registry,
				typeEnv:  scope.GetDeclarationScope(),
			}
			for j, key := range templateTypes {
				// Only add the template types that do not have a type transformation
				if !key.IsTypeTransformation() {
					dup := false
					for k := 0; k < j; k++ {
						if templateTypes[k] == key {
							dup = true
							break
						}
					}
					if !dup {
						typeVars = append(typeVars, ttlTypeVar{name: *key.GetReferenceName(), typ: values[j]})
					}
				}
			}
		}
		// Evaluate the type transformation expression using the current
		// known types for the template type variables
		transformedType := ttlObj.evalWithSlices(type_.GetTypeTransformation(), typeVars, nil)
		values[idx] = transformedType
		for j := range templateTypes {
			if j != idx && templateTypes[j] == type_ {
				values[j] = transformedType
			}
		}
		// Add the transformed type to the type variables
		refName := *type_.GetReferenceName()
		updated := false
		for i := range typeVars {
			if typeVars[i].name == refName {
				typeVars[i].typ = transformedType
				updated = true
				break
			}
		}
		if !updated {
			typeVars = append(typeVars, ttlTypeVar{name: refName, typ: transformedType})
		}
	}
}

// InferTemplatedTypesForCall: for functions that use template types,
// specialize the function type for the call target based on the call-site
// specific arguments. Specifically, this enables inference to set the type
// of any function literal parameters based on these inferred types.
//
// Java: TypeInference.inferTemplatedTypesForCall(Node, FunctionType, FlowScope)
func (ti *TypeInference) InferTemplatedTypesForCall(n rhino.Node, fnType rhino.FunctionType, scope *LinkedFlowScope) bool {
	keys := fnType.GetTemplateTypeMap().GetTemplateKeys()
	if len(keys) == 0 {
		return false
	}

	// Try to infer the template types
	matcher := InvocationTemplateTypeMatcher{
		registry:      ti.registry,
		calleeType:    fnType,
		localThisType: scope.GetTypeOfThis(),
		invocation:    n,
		unknownType:   ti.unknownType,
	}
	matcher.match()

	var valuesBuf [8]rhino.JSType
	var values []rhino.JSType
	if len(keys) <= len(valuesBuf) {
		values = valuesBuf[:len(keys)]
	} else {
		values = make([]rhino.JSType, len(keys))
	}
	// If the inferred type doesn't satisfy the template bound, swap to using the bound. This
	// ensures errors will be reported in type-checking.
	for i, key := range keys {
		v := matcher.getOrDefault(key, ti.unknownType)
		if !v.IsSubtypeOf(key.GetBound()) {
			v = key.GetBound()
		}
		values[i] = v
	}

	// Try to infer the template types using the type transformations
	ti.evaluateTypeTransformations(keys, values, scope)

	for i, v := range values {
		if v == nil {
			values[i] = ti.unknownType
		}
	}

	// Replace all template types. If we couldn't find a replacement, we
	// replace it with UNKNOWN.
	replacer := rhino.MakeTemplateTypeReplacer(
		ti.registry,
		ti.registry.GetEmptyTemplateTypeMap().CopyWithExtensionOfKeys(keys, values),
		true,
		true,
		true)
	callTarget := ast.MustInvocation(n).Callee().Node()

	replacementFnType := replacer.VisitJSType(fnType).ToMaybeFunctionType()
	jutil.CheckNotNull(replacementFnType)
	callTarget.SetJSType(replacementFnType)
	n.SetJSType(replacementFnType.GetReturnType())

	return replacer.HasMadeReplacement()
}

// Java: TypeInference.traverseNew(Node, FlowScope)
func (ti *TypeInference) TraverseNew(n rhino.Node, scope *LinkedFlowScope) *LinkedFlowScope {
	scope = ti.TraverseChildren(n, scope)
	constructor := ast.MustNew(n).Callee().Node()
	constructorType := constructor.GetJSType()
	return ti.TraverseInstantiation(n, constructorType, scope)
}

// Java: TypeInference.traverseInstantiation(Node, JSType, FlowScope)
func (ti *TypeInference) TraverseInstantiation(n rhino.Node, ctorType rhino.JSType, scope *LinkedFlowScope) *LinkedFlowScope {
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
		matcher := InvocationTemplateTypeMatcher{
			registry:      ti.registry,
			calleeType:    ctorFnType,
			localThisType: scope.GetTypeOfThis(),
			invocation:    n,
			unknownType:   ti.unknownType,
		}
		matcher.match()
		typeParams := instantiatedType.GetTypeParameters()
		var typeArgs []rhino.JSType
		if len(typeParams) > 0 {
			typeArgs = make([]rhino.JSType, len(typeParams))
			for i, param := range typeParams {
				typeArgs[i] = matcher.getOrDefault(param, ti.unknownType)
			}
		}
		instantiatedType =
			ti.registry.CreateTemplatizedTypeOfList(instantiatedType, typeArgs).ToMaybeObjectType()
	}

	if instantiatedType != nil {
		n.SetJSType(instantiatedType)
	} else {
		n.SetJSType(ti.unknownType)
	}
	return scope
}

// Java: TypeInference.traverseSuper(Node, FlowScope)
func (ti *TypeInference) TraverseSuper(superNode rhino.Node, currentScope *LinkedFlowScope) {
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
			if superCtor := enclosingFunctionType.GetSuperClassConstructor(); superCtor != nil {
				superNodeType = superCtor
			}
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
