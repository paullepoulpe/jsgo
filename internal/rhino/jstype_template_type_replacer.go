// Port of com/google/javascript/rhino/jstype/TemplateTypeReplacer.java @ closure-compiler v20261021

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   John Lenz
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino

import (
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// TemplateTypeReplacer specializes TemplatizedTypes according to provided
// bindings.
//
// Java: implements Visitor<JSType>
type TemplateTypeReplacer struct {
	registry *JSTypeRegistry
	bindings *TemplateTypeMap

	visitProperties bool
	// TODO(nickreid): We should only need `useUnknownForMissingBinding`.
	// Keeping two separate bits was a quick fix for collapsing two different
	// classes.
	useUnknownForMissingKeys   bool
	useUnknownForMissingValues bool

	hasMadeReplacement bool
	keyType            *TemplateType
	// The index in the corresponding bindings pointing to where the submap for
	// the specified key type ends, exclusive; or if keyType is null, defaults
	// to -1.
	ownSubMapBoundary int

	seenTypes          []JSType
	visitedObjectTypes map[JSType]JSType
}

// TemplateTypeReplacer_ForInference creates a replacer for use during
// TypeInference.
//
// Java: TemplateTypeReplacer.forInference(JSTypeRegistry, Map<TemplateType, JSType>)
func TemplateTypeReplacer_ForInference(
	registry *JSTypeRegistry, bindings *jutil.LinkedHashMap[*TemplateType, JSType]) *TemplateTypeReplacer {
	r := TemplateTypeReplacer_MakeForInference(registry, bindings)
	return &r
}

// TemplateTypeReplacer_MakeForInference creates a value replacer for use during
// TypeInference without heap allocation.
func TemplateTypeReplacer_MakeForInference(
	registry *JSTypeRegistry, bindings *jutil.LinkedHashMap[*TemplateType, JSType]) TemplateTypeReplacer {
	n := bindings.Size()
	if n == 0 {
		return MakeTemplateTypeReplacer(registry, registry.GetEmptyTemplateTypeMap(), true, true, true)
	}
	keys := make([]*TemplateType, 0, n)
	values := make([]JSType, 0, n)
	unknown := registry.GetNativeType(JSTypeNative_UNKNOWN_TYPE)
	for key, value := range bindings.All() {
		if key == nil {
			panic(jutil.NewNullPointerExceptionMsg("at index " + jutil.ToJavaString(len(keys))))
		}
		keys = append(keys, key)
		if value != nil {
			values = append(values, value)
		} else {
			values = append(values, unknown)
		}
	}
	m := registry.GetEmptyTemplateTypeMap().CopyWithExtensionOfKeys(keys, values)
	return MakeTemplateTypeReplacer(registry, m, true, true, true)
}

// TemplateTypeReplacer_ForTotalReplacement creates a replacer that will
// always totally eliminate TemplateTypes from the definitions of the types
// it performs replacement on.
//
// If a binding for a TemplateType is required but not provided, `?` will be
// used.
//
// Java: TemplateTypeReplacer.forTotalReplacement(JSTypeRegistry, TemplateTypeMap)
func TemplateTypeReplacer_ForTotalReplacement(registry *JSTypeRegistry, bindings *TemplateTypeMap) *TemplateTypeReplacer {
	return NewTemplateTypeReplacer(registry, bindings, false, false, true)
}

// TemplateTypeReplacer_ForPartialReplacement creates a replacer that may
// not totally eliminate TemplateTypes from the definitions of the types it
// performs replacement on.
//
// If a binding for a TemplateType is required but not provided, uses of that
// type will not be replaced.
//
// Java: TemplateTypeReplacer.forPartialReplacement(JSTypeRegistry, TemplateTypeMap)
func TemplateTypeReplacer_ForPartialReplacement(registry *JSTypeRegistry, bindings *TemplateTypeMap) *TemplateTypeReplacer {
	return NewTemplateTypeReplacer(registry, bindings, false, false, false)
}

// TemplateTypeReplacer_MakeForPartialReplacement creates a value replacer for partial replacement without heap allocation.
func TemplateTypeReplacer_MakeForPartialReplacement(registry *JSTypeRegistry, bindings *TemplateTypeMap) TemplateTypeReplacer {
	return MakeTemplateTypeReplacer(registry, bindings, false, false, false)
}

// MakeTemplateTypeReplacer constructs a TemplateTypeReplacer value.
func MakeTemplateTypeReplacer(
	registry *JSTypeRegistry,
	bindings *TemplateTypeMap,
	visitProperties bool,
	useUnknownForMissingKeys bool,
	useUnknownForMissingValues bool) TemplateTypeReplacer {
	return TemplateTypeReplacer{
		registry:                   registry,
		bindings:                   bindings,
		visitProperties:            visitProperties,
		useUnknownForMissingKeys:   useUnknownForMissingKeys,
		useUnknownForMissingValues: useUnknownForMissingValues,
		ownSubMapBoundary:          -1,
	}
}

// Java: TemplateTypeReplacer(JSTypeRegistry, TemplateTypeMap, boolean, boolean, boolean)
func NewTemplateTypeReplacer(
	registry *JSTypeRegistry,
	bindings *TemplateTypeMap,
	visitProperties bool,
	useUnknownForMissingKeys bool,
	useUnknownForMissingValues bool) *TemplateTypeReplacer {
	r := MakeTemplateTypeReplacer(registry, bindings, visitProperties, useUnknownForMissingKeys, useUnknownForMissingValues)
	return &r
}

// Java: TemplateTypeReplacer.hasMadeReplacement()
func (r *TemplateTypeReplacer) HasMadeReplacement() bool {
	return r.hasMadeReplacement
}

// InitSeenTypes: determinism is unnecessary for seenTypes, because we only
// call contains/add/remove.
//
// Java: TemplateTypeReplacer.initSeenTypes()
func (r *TemplateTypeReplacer) InitSeenTypes() {
}

func (r *TemplateTypeReplacer) seenContains(t JSType) bool {
	return slices.Contains(r.seenTypes, t)
}

func (r *TemplateTypeReplacer) seenAdd(t JSType) bool {
	if slices.Contains(r.seenTypes, t) {
		return false
	}
	r.seenTypes = append(r.seenTypes, t)
	return true
}

func (r *TemplateTypeReplacer) seenRemove(t JSType) {
	for i := len(r.seenTypes) - 1; i >= 0; i-- {
		if r.seenTypes[i] == t {
			last := len(r.seenTypes) - 1
			r.seenTypes[i] = r.seenTypes[last]
			r.seenTypes[last] = nil
			r.seenTypes = r.seenTypes[:last]
			return
		}
	}
}

// VisitJSType visits type_ directly without boxing r into a Visitor interface.
func (r *TemplateTypeReplacer) VisitJSType(type_ JSType) JSType {
	return r.visitJSType(type_)
}

func (r *TemplateTypeReplacer) visitJSType(type_ JSType) JSType {
	switch x := type_.(type) {
	case *TemplateType:
		return r.CaseTemplateType(x)
	case *UnionType:
		return r.CaseUnionType(x)
	case *TemplatizedType:
		return r.CaseTemplatizedType(x)
	case *FunctionTypeBase:
		return r.CaseFunctionType(x)
	case *InstanceObjectType, *PrototypeObjectTypeBase, *EnumType, *NoResolvedType, *EnumElementType, *NamedType, *NoType:
		return type_
	case *RecordType:
		return r.CaseObjectType(x)
	case *UnknownType:
		return r.GetNativeType(JSTypeNative_UNKNOWN_TYPE)
	case *NumberType:
		return r.GetNativeType(JSTypeNative_NUMBER_TYPE)
	case *StringType:
		return r.GetNativeType(JSTypeNative_STRING_TYPE)
	case *BooleanType:
		return r.GetNativeType(JSTypeNative_BOOLEAN_TYPE)
	case *NullType:
		return r.GetNativeType(JSTypeNative_NULL_TYPE)
	case *VoidType:
		return r.GetNativeType(JSTypeNative_VOID_TYPE)
	case *AllType:
		return r.GetNativeType(JSTypeNative_ALL_TYPE)
	case *BigIntType:
		return r.GetNativeType(JSTypeNative_BIGINT_TYPE)
	case *NoObjectTypeBase:
		return r.GetNativeType(JSTypeNative_NO_OBJECT_TYPE)
	case *SymbolTypeBase, *KnownSymbolType:
		return type_
	case *ProxyObjectTypeBase:
		return r.CaseProxyObjectType(x)
	default:
		if fn, ok := type_.(FunctionType); ok {
			if _, isNoObj := type_.(NoObjectType); isNoObj {
				return r.GetNativeType(JSTypeNative_NO_OBJECT_TYPE)
			}
			return r.CaseFunctionType(fn)
		}
		if proxy, ok := type_.(ProxyObjectType); ok {
			return r.CaseProxyObjectType(proxy)
		}
		if obj, ok := type_.(ObjectType); ok {
			return r.CaseObjectType(obj)
		}
		if _, ok := type_.(SymbolType); ok {
			return type_
		}
		panic(jutil.NewUnsupportedOperationException())
	}
}

// Java: TemplateTypeReplacer.caseNoType(NoType)
func (r *TemplateTypeReplacer) CaseNoType(type_ *NoType) JSType {
	return type_
}

// Java: TemplateTypeReplacer.caseEnumElementType(EnumElementType)
func (r *TemplateTypeReplacer) CaseEnumElementType(type_ *EnumElementType) JSType {
	return type_
}

// Java: TemplateTypeReplacer.caseAllType()
func (r *TemplateTypeReplacer) CaseAllType() JSType {
	return r.GetNativeType(JSTypeNative_ALL_TYPE)
}

// Java: TemplateTypeReplacer.caseBooleanType()
func (r *TemplateTypeReplacer) CaseBooleanType() JSType {
	return r.GetNativeType(JSTypeNative_BOOLEAN_TYPE)
}

// Java: TemplateTypeReplacer.caseNoObjectType()
func (r *TemplateTypeReplacer) CaseNoObjectType() JSType {
	return r.GetNativeType(JSTypeNative_NO_OBJECT_TYPE)
}

// Java: TemplateTypeReplacer.caseFunctionType(FunctionType)
func (r *TemplateTypeReplacer) CaseFunctionType(type_ FunctionType) JSType {
	if r.IsNativeFunctionType(type_) {
		return type_
	}
	if !type_.IsOrdinaryFunction() && !type_.IsConstructor() {
		return type_
	}
	if !r.seenAdd(type_) {
		return type_
	}
	res := r.CaseFunctionTypeUnguarded(type_)
	r.seenRemove(type_)
	return res
}

// Java: TemplateTypeReplacer.caseFunctionTypeUnguarded(FunctionType)
func (r *TemplateTypeReplacer) CaseFunctionTypeUnguarded(type_ FunctionType) JSType {
	if r.IsNativeFunctionType(type_) {
		return type_
	}

	if !type_.IsOrdinaryFunction() && !type_.IsConstructor() {
		return type_
	}

	changed := false

	beforeThis := type_.GetTypeOfThis()
	afterThis := r.CoerseToThisType(r.visitJSType(beforeThis))
	if beforeThis != afterThis {
		changed = true
	}

	beforeReturn := type_.GetReturnType()
	afterReturn := r.visitJSType(beforeReturn)
	if beforeReturn != afterReturn {
		changed = true
	}

	paramsChanged := false
	parameters := type_.GetParameters()
	numParams := len(parameters)
	var paramBuilder FunctionParamBuilder
	for i := 0; i < numParams; i++ {
		parameter := parameters[i]
		beforeParamType := parameter.GetJSType()
		afterParamType := r.visitJSType(beforeParamType)

		if beforeParamType != afterParamType {
			if !paramsChanged {
				changed = true
				paramsChanged = true
				paramBuilder.registry = r.registry
				paramBuilder.parameters = make([]*FunctionType_Parameter, i, numParams)
				copy(paramBuilder.parameters, parameters[:i])
			}
			if parameter.IsOptional() {
				paramBuilder.AddOptionalParam(afterParamType)
			} else if parameter.IsVariadic() {
				paramBuilder.AddVarArgs(afterParamType)
			} else {
				paramBuilder.AddRequiredParam(afterParamType)
			}
		} else if paramsChanged {
			paramBuilder.parameters = append(paramBuilder.parameters, parameter)
		}
	}

	if changed {
		var params []*FunctionType_Parameter
		if paramsChanged {
			params = paramBuilder.parameters
		} else {
			params = parameters
		}
		if fb, ok := type_.(*FunctionTypeBase); ok {
			var builder FunctionType_Builder
			fb.initBuilder(&builder)
			return builder.
				WithParametersOfList(params).
				WithReturnType(afterReturn).
				WithTypeOfThis(afterThis).
				WithIsAbstract(false). // TODO(b/187989034): Copy this from the source function.
				BuildAndResolve()
		}
		return type_.ToBuilder().
			WithParametersOfList(params).
			WithReturnType(afterReturn).
			WithTypeOfThis(afterThis).
			WithIsAbstract(false). // TODO(b/187989034): Copy this from the source function.
			Build()
	}

	return type_
}

// Java: TemplateTypeReplacer.coerseToThisType(JSType)
func (r *TemplateTypeReplacer) CoerseToThisType(type_ JSType) JSType {
	if type_ != nil {
		return type_
	}
	return r.registry.GetNativeObjectType(JSTypeNative_UNKNOWN_TYPE)
}

// Java: TemplateTypeReplacer.caseObjectType(ObjectType)
func (r *TemplateTypeReplacer) CaseObjectType(objType ObjectType) JSType {
	_, isProxy := objType.(ProxyObjectType)
	if !r.visitProperties ||
		objType.IsNominalType() ||
		isProxy ||
		!objType.IsRecordType() {
		return objType
	}
	if !r.seenAdd(objType) {
		return objType
	}
	res := r.CaseObjectTypeUnguarded(objType)
	r.seenRemove(objType)
	return res
}

// Java: TemplateTypeReplacer.caseObjectTypeUnguarded(ObjectType)
func (r *TemplateTypeReplacer) CaseObjectTypeUnguarded(objType ObjectType) JSType {
	_, isProxy := objType.(ProxyObjectType)
	if !r.visitProperties ||
		objType.IsNominalType() ||
		isProxy ||
		!objType.IsRecordType() {
		return objType
	}

	// If we've visited this object type, we don't need to walk it again.
	cached := r.GetVisitedObjectTypeOrNull(objType)
	if cached != nil {
		return cached
	}

	ownProps := objType.GetOwnPropertyNames()
	var builder *RecordTypeBuilder
	for i, prop := range ownProps {
		beforeType := objType.GetPropertyType(prop)
		afterType := r.visitJSType(beforeType)
		if builder != nil {
			builder.AddProperty(prop, afterType, objType.GetPropertyNode(prop))
		} else if beforeType != afterType {
			builder = NewRecordTypeBuilder(r.registry)
			for j := 0; j < i; j++ {
				prevProp := ownProps[j]
				builder.AddProperty(prevProp, objType.GetPropertyType(prevProp), objType.GetPropertyNode(prevProp))
			}
			builder.AddProperty(prop, afterType, objType.GetPropertyNode(prop))
		}
	}

	// Use our new type if anything changed, and be sure to update the cache.
	var result JSType = objType
	if builder != nil {
		result = builder.Build()
	}
	r.putVisitedObjectType(objType, result)
	return result
}

// Java: TemplateTypeReplacer.caseTemplatizedType(TemplatizedType)
func (r *TemplateTypeReplacer) CaseTemplatizedType(type_ *TemplatizedType) JSType {
	if !r.seenAdd(type_) {
		return type_
	}
	res := r.CaseTemplatizedTypeUnguarded(type_)
	r.seenRemove(type_)
	return res
}

// Java: TemplateTypeReplacer.caseTemplatizedTypeUnguarded(TemplatizedType)
func (r *TemplateTypeReplacer) CaseTemplatizedTypeUnguarded(type_ *TemplatizedType) JSType {
	changed := false
	beforeBaseType := type_.GetReferencedType()
	afterBaseType := ObjectType_Cast(r.visitJSType(beforeBaseType))
	if beforeBaseType != afterBaseType {
		changed = true
	}

	templateTypes := type_.GetTemplateTypes()
	n := len(templateTypes)
	var builder []JSType
	for i := 0; i < n; i++ {
		beforeTemplateType := templateTypes[i]
		afterTemplateType := r.visitJSType(beforeTemplateType)
		jutil.CheckNotNull(afterTemplateType)
		if builder != nil {
			builder[i] = afterTemplateType
		} else if beforeTemplateType != afterTemplateType {
			changed = true
			builder = make([]JSType, n)
			copy(builder[:i], templateTypes[:i])
			builder[i] = afterTemplateType
		}
	}

	if changed {
		if builder == nil {
			builder = templateTypes
		}
		type_ = r.registry.CreateTemplatizedTypeOfList(afterBaseType, builder)
	}
	return type_
}

// Java: TemplateTypeReplacer.caseUnknownType()
func (r *TemplateTypeReplacer) CaseUnknownType() JSType {
	return r.GetNativeType(JSTypeNative_UNKNOWN_TYPE)
}

// Java: TemplateTypeReplacer.caseNullType()
func (r *TemplateTypeReplacer) CaseNullType() JSType {
	return r.GetNativeType(JSTypeNative_NULL_TYPE)
}

// Java: TemplateTypeReplacer.caseNumberType()
func (r *TemplateTypeReplacer) CaseNumberType() JSType {
	return r.GetNativeType(JSTypeNative_NUMBER_TYPE)
}

// Java: TemplateTypeReplacer.caseBigIntType()
func (r *TemplateTypeReplacer) CaseBigIntType() JSType {
	return r.GetNativeType(JSTypeNative_BIGINT_TYPE)
}

// Java: TemplateTypeReplacer.caseStringType()
func (r *TemplateTypeReplacer) CaseStringType() JSType {
	return r.GetNativeType(JSTypeNative_STRING_TYPE)
}

// Java: TemplateTypeReplacer.caseSymbolType(SymbolType)
func (r *TemplateTypeReplacer) CaseSymbolType(type_ SymbolType) JSType {
	return type_
}

// Java: TemplateTypeReplacer.caseVoidType()
func (r *TemplateTypeReplacer) CaseVoidType() JSType {
	return r.GetNativeType(JSTypeNative_VOID_TYPE)
}

// Java: TemplateTypeReplacer.caseUnionType(UnionType)
func (r *TemplateTypeReplacer) CaseUnionType(type_ *UnionType) JSType {
	if !r.seenAdd(type_) {
		return type_
	}
	res := r.CaseUnionTypeUnguarded(type_)
	r.seenRemove(type_)
	return res
}

// Java: TemplateTypeReplacer.caseUnionTypeUnguarded(UnionType)
func (r *TemplateTypeReplacer) CaseUnionTypeUnguarded(type_ *UnionType) JSType {
	alternates := type_.GetAlternates()
	alternateCount := len(alternates)
	var results []JSType
	for i := 0; i < alternateCount; i++ {
		alternative := alternates[i]
		replacement := r.visitJSType(alternative)
		if results != nil {
			results[i] = replacement
		} else if replacement != alternative {
			results = make([]JSType, alternateCount)
			copy(results[:i], alternates[:i])
			results[i] = replacement
		}
	}

	if results != nil {
		b := r.registry.acquireUnionBuilder()
		res := b.AddAlternates(results).Build() // maybe not a union
		r.registry.releaseUnionBuilder(b)
		return res
	}

	return type_
}

// Java: TemplateTypeReplacer.caseTemplateType(TemplateType)
func (r *TemplateTypeReplacer) CaseTemplateType(type_ *TemplateType) JSType {
	r.hasMadeReplacement = true

	if type_ != r.keyType && !r.bindings.HasTemplateKeyFrom(type_, r.ownSubMapBoundary) {
		if r.useUnknownForMissingKeys {
			return r.GetNativeType(JSTypeNative_UNKNOWN_TYPE)
		}
		return type_
	}

	if r.seenContains(type_) {
		// If we have already encountered this TemplateType during replacement
		// (i.e. there is a reference loop) then return the TemplateType type
		// itself.
		return type_
	} else if !r.bindings.HasTemplateType(type_) {
		// If there is no JSType substitution for the TemplateType, return either
		// the UNKNOWN_TYPE or the TemplateType type itself, depending on
		// configuration.
		if r.useUnknownForMissingValues {
			return r.GetNativeType(JSTypeNative_UNKNOWN_TYPE)
		}
		return type_
	}

	replacement := r.bindings.GetUnresolvedOriginalTemplateType(type_)
	// Recursive templatized types, such as T => Foo<T>. Don't do any
	// replacement - we should preserve Foo<T> as it was before. Note: this
	// isn't perfect - if we had T => Foo<T, U>, this will skip replacing U -
	// but it's better than a stack overflow. Return the 'type' instead of
	// 'replacement' to notify callers that no changes were made. If we didn't
	// check this here, we'd still skip doing any work in
	// replacement.visit(this) because guardAgainstCycles would have already
	// seen 'replacement'. But we'd return 'visitedReplacement' instead of
	// 'type'.
	if replacement == r.keyTypeAsJSType() || r.seenContains(replacement) {
		return type_
	}
	r.seenAdd(type_)
	visitedReplacement := r.visitJSType(replacement)
	r.seenRemove(type_)

	jutil.CheckStateMsg1(
		visitedReplacement != r.keyTypeAsJSType(),
		"Trying to replace key %s with the same value",
		r.keyTypeAsJSType())
	return visitedReplacement
}

// keyTypeAsJSType is keyType as a JSType, null if it is null.
func (r *TemplateTypeReplacer) keyTypeAsJSType() JSType {
	if r.keyType == nil {
		return nil
	}
	return r.keyType
}

// Java: TemplateTypeReplacer.getNativeType(JSTypeNative)
func (r *TemplateTypeReplacer) GetNativeType(nativeType JSTypeNative) JSType {
	return r.registry.GetNativeType(nativeType)
}

// Java: TemplateTypeReplacer.isNativeFunctionType(FunctionType)
func (r *TemplateTypeReplacer) IsNativeFunctionType(type_ FunctionType) bool {
	return type_.IsNativeObjectType()
}

// Java: TemplateTypeReplacer.caseNamedType(NamedType)
func (r *TemplateTypeReplacer) CaseNamedType(type_ *NamedType) JSType {
	// The internals of a named type aren't interesting.
	return type_
}

// Java: TemplateTypeReplacer.caseProxyObjectType(ProxyObjectType)
func (r *TemplateTypeReplacer) CaseProxyObjectType(type_ ProxyObjectType) JSType {
	if !r.seenAdd(type_) {
		return type_
	}
	res := r.CaseProxyObjectTypeUnguarded(type_)
	r.seenRemove(type_)
	return res
}

// Java: TemplateTypeReplacer.caseProxyObjectTypeUnguarded(ProxyObjectType)
func (r *TemplateTypeReplacer) CaseProxyObjectTypeUnguarded(type_ ProxyObjectType) JSType {
	// Be careful not to unwrap a type unless it has changed.
	beforeType := type_.GetReferencedTypeInternal()
	replacement := r.visitJSType(beforeType)
	if replacement != beforeType {
		return replacement
	}
	return type_
}

// Java: TemplateTypeReplacer.setKeyType(TemplateType, int)
func (r *TemplateTypeReplacer) SetKeyType(keyType *TemplateType, ownSubMapBoundary int) {
	r.keyType = keyType
	r.ownSubMapBoundary = ownSubMapBoundary
}

// GetVisitedObjectTypeOrNull: determinism is unnecessary for
// visitedObjectTypes, because we only call get/set
//
// Java: TemplateTypeReplacer.getVisitedObjectTypeOrNull(JSType)
func (r *TemplateTypeReplacer) GetVisitedObjectTypeOrNull(type_ JSType) JSType {
	return r.visitedObjectTypes[type_]
}

func (r *TemplateTypeReplacer) putVisitedObjectType(k, v JSType) {
	if r.visitedObjectTypes == nil {
		r.visitedObjectTypes = make(map[JSType]JSType)
	}
	r.visitedObjectTypes[k] = v
}
