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
	"slices"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type templateTypeSlot struct {
	key *rhino.TemplateType
	typ rhino.JSType
}

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
	matchedBuf [8]templateTypeSlot
	matched    []templateTypeSlot
	seenBuf    [4]rhino.JSType
	seenTypes  []rhino.JSType

	registry      *rhino.JSTypeRegistry
	calleeType    rhino.FunctionType
	localThisType rhino.JSType
	invocation    rhino.Node

	unknownType rhino.JSType
}

// Java: InvocationTemplateTypeMatcher(JSTypeRegistry, FunctionType, JSType, Node)
func NewInvocationTemplateTypeMatcher(
	registry *rhino.JSTypeRegistry, calleeType rhino.FunctionType, localThisType rhino.JSType, invocation rhino.Node) *InvocationTemplateTypeMatcher {
	return &InvocationTemplateTypeMatcher{
		registry:      registry,
		calleeType:    calleeType,
		localThisType: localThisType,
		invocation:    invocation,
		unknownType:   registry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE),
	}
}

// Java: InvocationTemplateTypeMatcher.match()
func (m *InvocationTemplateTypeMatcher) match() {
	m.matched = m.matchedBuf[:0]
	m.seenTypes = m.seenBuf[:0]
	if m.calleeType.GetTemplateTypeMap().IsEmpty() {
		return
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
			return
		}

		// The first argument to the tag function is an array of strings (typed as ITemplateArray)
		// but not an actual AST node
		m.matchTemplateTypesRecursive(
			calleeParameters[0].GetJSType(), m.registry.GetNativeType(rhino.JSTypeNative_I_TEMPLATE_ARRAY_TYPE))

		// Resolve the remaining template types from the template literal substitutions.
		declParams := calleeParameters[1:]
		i := 0
		var variadicParam rhino.FunctionType_Parameter
		for callParam := range NewNodeUtil_TemplateArgsIterable(m.invocation.GetLastChild()).All() {
			if !m.matchCallArg(declParams, &i, &variadicParam, callParam) {
				break
			}
		}
	} else if m.invocation.HasMoreThanOneChild() {
		declParams := m.calleeType.GetParameters()
		i := 0
		var variadicParam rhino.FunctionType_Parameter
		for callParam := m.invocation.GetSecondChild(); !callParam.IsNil(); callParam = callParam.GetNext() {
			if !m.matchCallArg(declParams, &i, &variadicParam, callParam) {
				break
			}
		}
	}
}

func (m *InvocationTemplateTypeMatcher) getOrDefault(key *rhino.TemplateType, defaultValue rhino.JSType) rhino.JSType {
	for i := range m.matched {
		if m.matched[i].key == key && m.matched[i].typ != nil {
			return m.matched[i].typ
		}
	}
	return defaultValue
}

func (m *InvocationTemplateTypeMatcher) matchCallArg(
	declParams []rhino.FunctionType_Parameter,
	i *int,
	variadicParam *rhino.FunctionType_Parameter,
	callParam rhino.Node) bool {
	if !variadicParam.IsZero() {
		m.matchTemplateTypesRecursive(variadicParam.GetJSType(), m.getTypeOrUnknown(callParam))
		return true
	}
	if *i >= len(declParams) {
		return false
	}
	declParam := declParams[*i]
	*i++
	m.matchTemplateTypesRecursive(declParam.GetJSType(), m.getTypeOrUnknown(callParam))
	if declParam.IsVariadic() {
		*variadicParam = declParam
	}
	return true
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
		if !slices.Contains(m.seenTypes, paramType) {
			m.seenTypes = append(m.seenTypes, paramType)
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
			m.seenTypes = m.seenTypes[:len(m.seenTypes)-1]
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

// Java: InvocationTemplateTypeMatcher.matchTemplateTypesFromParameters(Iterator<Parameter>, Iterator<Parameter>)
func (m *InvocationTemplateTypeMatcher) matchTemplateTypesFromParameters(
	declParams []rhino.FunctionType_Parameter, callParams []rhino.FunctionType_Parameter) {
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
	for i := range m.matched {
		if m.matched[i].key == template {
			if m.matched[i].typ != nil {
				m.matched[i].typ = m.matched[i].typ.GetLeastSupertype(match)
			} else {
				m.matched[i].typ = match
			}
			return
		}
	}
	m.matched = append(m.matched, templateTypeSlot{key: template, typ: match})
}

// Java: InvocationTemplateTypeMatcher.getTypeOrUnknown(Node)
func (m *InvocationTemplateTypeMatcher) getTypeOrUnknown(n rhino.Node) rhino.JSType {
	type_ := n.GetJSType()
	if type_ == nil {
		return m.unknownType
	}
	return type_
}
