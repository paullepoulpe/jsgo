// Port of com/google/javascript/jscomp/InvocationTemplateTypeMatcher.java @ closure-compiler v20261021

/*
 * Copyright 2007 The Closure Compiler Authors.
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
	"iter"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// InvocationTemplateTypeMatcher determines the types that fill any template
// parameters at a function invocation.
//
// Given an invocation of some function with type `F` templated on `T`, this
// class traverses the arguments of the invocation to determine what type
// `T` should be. Argument types and `F`'s parameter types are
// traversed/recursed in parallel to match template types at any depth
// inside `F`'s signature.
//
// Instances of this class are single use. They provide a "scope" for the
// matching but accumulate state as the matching progresses.
type InvocationTemplateTypeMatcher struct {
	matchedTypes *base.LinkedIdentityHashMap[*rhino.TemplateType, rhino.JSType]
	seenTypes    *base.LinkedIdentityHashSet[rhino.JSType]

	registry      *rhino.JSTypeRegistry
	calleeType    rhino.FunctionType
	localThisType rhino.JSType
	invocation    rhino.Node

	unknownType rhino.JSType
}

// Java: InvocationTemplateTypeMatcher(JSTypeRegistry, FunctionType, JSType, Node)
func NewInvocationTemplateTypeMatcher(
	registry *rhino.JSTypeRegistry, calleeType rhino.FunctionType, localThisType rhino.JSType, invocation rhino.Node) *InvocationTemplateTypeMatcher {
	m := &InvocationTemplateTypeMatcher{
		matchedTypes: base.NewLinkedIdentityHashMap[*rhino.TemplateType, rhino.JSType](),
		seenTypes:    base.NewLinkedIdentityHashSet[rhino.JSType](),
	}
	m.registry = registry
	m.calleeType = calleeType
	m.localThisType = localThisType
	m.invocation = invocation

	m.unknownType = registry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
	return m
}

// Java: InvocationTemplateTypeMatcher.match()
func (m *InvocationTemplateTypeMatcher) Match() *jutil.LinkedHashMap[*rhino.TemplateType, rhino.JSType] {
	if m.calleeType.GetTemplateTypeMap().IsEmpty() {
		return jutil.NewLinkedHashMap[*rhino.TemplateType, rhino.JSType]()
	}

	target := m.invocation.GetFirstChild()
	if NodeUtil_IsNormalGet(target) {
		obj := target.GetFirstChild()
		typeOfThisRequiredByTheFunction := m.calleeType.GetTypeOfThis()
		// The type placed on a SUPER node is the superclass type, which allows us to infer the right
		// property types for the GETPROP or GETELEM nodes built on it.
		// However, the type actually passed as `this` when making calls this way is the `this`
		// of the scope where the `super` appears.
		var typeOfThisProvidedByTheCall rhino.JSType
		if obj.IsSuper() {
			typeOfThisProvidedByTheCall = m.localThisType
		} else {
			typeOfThisProvidedByTheCall = m.getTypeOrUnknown(obj)
		}
		// We're looking at a call made as `obj['method']()` or `obj.method()` (see enclosing if),
		// so if the call is successfully made, then the object through which it is made isn't null
		// or undefined.
		typeOfThisProvidedByTheCall = typeOfThisProvidedByTheCall.RestrictByNotNullOrUndefined()
		m.matchTemplateTypesRecursive(
			typeOfThisRequiredByTheFunction, typeOfThisProvidedByTheCall)
	}

	if m.invocation.IsTaggedTemplateLit() {
		calleeParameters := m.calleeType.GetParameters()
		if len(calleeParameters) == 0 {
			// TypeCheck will warn if there are too few function parameters
			return m.copyMatchedTypes()
		}

		// The first argument to the tag function is an array of strings (typed as ITemplateArray)
		// but not an actual AST node
		m.matchTemplateTypesRecursive(
			calleeParameters[0].GetJSType(), m.registry.GetNativeType(rhino.JSTypeNative_I_TEMPLATE_ARRAY_TYPE))

		// Resolve the remaining template types from the template literal substitutions.
		m.matchTemplateTypesFromNodes(
			m.calleeType.GetParameters()[1:],
			NodeUtil_GetInvocationArgsAsIterable(m.invocation))
	} else if m.invocation.HasMoreThanOneChild() {
		m.matchTemplateTypesFromNodes(
			m.calleeType.GetParameters(), NodeUtil_GetInvocationArgsAsIterable(m.invocation))
	}

	return m.copyMatchedTypes()
}

// Java: InvocationTemplateTypeMatcher.matchTemplateTypesRecursive(JSType, JSType)
func (m *InvocationTemplateTypeMatcher) matchTemplateTypesRecursive(paramType rhino.JSType, argType rhino.JSType) {
	if paramType.IsTemplateType() {
		// Recursive base case.
		// example: @param {T}
		m.recordTemplateMatch(paramType.ToMaybeTemplateType(), argType)
		return
	}

	// Unpack unions.
	if paramType.IsUnionType() {
		// example: @param {Array.<T>|NodeList|Arguments|{length:number}}
		unionType := paramType.ToMaybeUnionType()
		for _, alternate := range unionType.GetAlternates() {
			m.matchTemplateTypesRecursive(alternate, argType)
		}
		return
	} else if argType.IsUnionType() {
		unionType := argType.ToMaybeUnionType()
		for _, alternate := range unionType.GetAlternates() {
			m.matchTemplateTypesRecursive(paramType, alternate)
		}
		return
	}

	if paramType.IsFunctionType() {
		paramFunctionType := paramType.ToMaybeFunctionType()
		argFunctionType :=
			argType.RestrictByNotNullOrUndefined().CollapseUnion().ToMaybeFunctionType()
		if argFunctionType != nil && argFunctionType.IsSubtype(paramType) {
			// infer from return type of the function type
			m.matchTemplateTypesRecursive(
				paramFunctionType.GetTypeOfThis(), argFunctionType.GetTypeOfThis())
			// infer from return type of the function type
			m.matchTemplateTypesRecursive(
				paramFunctionType.GetReturnType(), argFunctionType.GetReturnType())
			// infer from parameter types of the function type
			m.matchTemplateTypesFromParameters(
				paramFunctionType.GetParameters(),
				argFunctionType.GetParameters())
		}
	} else if paramType.IsRecordType() && !paramType.IsNominalType() {
		// example: @param {{foo:T}}
		if m.seenTypes.Add(paramType) {
			paramRecordType := paramType.ToObjectType()
			argObjectType := argType.RestrictByNotNullOrUndefined().ToObjectType()
			if argObjectType != nil &&
				!argObjectType.IsUnknownType() &&
				!argObjectType.IsEmptyType() {
				names := paramRecordType.GetPropertyNames()
				for _, name := range names {
					if paramRecordType.HasOwnProperty(name) && argObjectType.HasProperty(name) {
						m.matchTemplateTypesRecursive(
							paramRecordType.GetPropertyType(name), argObjectType.GetPropertyType(name))
					}
				}
			}
			m.seenTypes.Remove(paramType)
		}
	} else if paramType.IsTemplatizedType() {
		// example: @param {Array<T>}
		templatizedParamType := paramType.ToMaybeTemplatizedType()
		keyCount := len(templatizedParamType.GetTemplateTypes())
		// TODO(johnlenz): determine why we are creating TemplatizedTypes for
		// types with no type arguments.
		if keyCount == 0 {
			return
		}

		argObjectType := argType.RestrictByNotNullOrUndefined().CollapseUnion()
		// Resolve any template types in common between the argument type and parameter type
		paramTypeMap := paramType.GetTemplateTypeMap()

		keys := paramTypeMap.GetTemplateKeys()
		argTypeMap := argObjectType.GetTemplateTypeMap()
		for index := len(keys) - keyCount; index < len(keys); index++ {
			key := keys[index]
			m.matchTemplateTypesRecursive(
				paramTypeMap.GetResolvedTemplateType(key), argTypeMap.GetResolvedTemplateType(key))
		}
	}
}

// matchTemplateTypesFromNodes matches the declared parameters against the
// types of the call's argument nodes.
//
// Go-only: language-adapter — Java's overload on (Iterable<Parameter>,
// Iterable<Node>) only takes the iterators of its arguments and delegates
// to the (Iterator, Iterator) overload; both are this function, which walks
// the parameter slice while ranging over the argument nodes (taking a
// parameter and an argument per step, then every remaining argument once a
// variadic parameter is reached, exactly as the Java loops do).
//
// Java: InvocationTemplateTypeMatcher.matchTemplateTypesFromNodes(Iterable<Parameter>, Iterable<Node>)
// Java: InvocationTemplateTypeMatcher.matchTemplateTypesFromNodes(Iterator<Parameter>, Iterator<Node>)
func (m *InvocationTemplateTypeMatcher) matchTemplateTypesFromNodes(
	declParams []*rhino.FunctionType_Parameter, callParams iter.Seq[rhino.Node]) {
	i := 0
	var variadicParam *rhino.FunctionType_Parameter
	for callParam := range callParams {
		if variadicParam != nil {
			m.matchTemplateTypesRecursive(
				variadicParam.GetJSType(), m.getTypeOrUnknown(callParam))
			continue
		}
		if i >= len(declParams) {
			break
		}
		declParam := declParams[i]
		i++
		m.matchTemplateTypesRecursive(
			declParam.GetJSType(), m.getTypeOrUnknown(callParam))

		if declParam.IsVariadic() {
			variadicParam = declParam
		}
	}
}

// Java: InvocationTemplateTypeMatcher.matchTemplateTypesFromParameters(Iterator<Parameter>, Iterator<Parameter>)
func (m *InvocationTemplateTypeMatcher) matchTemplateTypesFromParameters(
	declParams []*rhino.FunctionType_Parameter, callParams []*rhino.FunctionType_Parameter) {
	i, j := 0, 0
	for i < len(declParams) && j < len(callParams) {
		declParam := declParams[i]
		i++
		m.matchTemplateTypesRecursive(declParam.GetJSType(), callParams[j].GetJSType())
		j++

		if declParam.IsVariadic() {
			for j < len(callParams) {
				m.matchTemplateTypesRecursive(declParam.GetJSType(), callParams[j].GetJSType())
				j++
			}
		}
	}
}

// Java: InvocationTemplateTypeMatcher.recordTemplateMatch(TemplateType, JSType)
func (m *InvocationTemplateTypeMatcher) recordTemplateMatch(template *rhino.TemplateType, match rhino.JSType) {
	if match.IsUnknownType() {
		return
	}

	// Don't worry about checking bounds here. We'll validate them once they're all collected.
	m.matchedTypes.Merge(template, match, func(a rhino.JSType, b rhino.JSType) rhino.JSType {
		return a.GetLeastSupertype(b)
	})
}

// Java: InvocationTemplateTypeMatcher.copyMatchedTypes()
func (m *InvocationTemplateTypeMatcher) copyMatchedTypes() *jutil.LinkedHashMap[*rhino.TemplateType, rhino.JSType] {
	builder := jutil.NewImmutableMapBuilder[*rhino.TemplateType, rhino.JSType]()
	m.matchedTypes.ForEach(func(k *rhino.TemplateType, v rhino.JSType) { builder.Put(k, v) })
	return builder.BuildOrThrow()
}

// Java: InvocationTemplateTypeMatcher.getTypeOrUnknown(Node)
func (m *InvocationTemplateTypeMatcher) getTypeOrUnknown(n rhino.Node) rhino.JSType {
	type_ := n.GetJSType()
	if type_ == nil {
		return m.unknownType
	}
	return type_
}
