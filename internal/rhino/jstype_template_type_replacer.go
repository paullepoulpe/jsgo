// Port of com/google/javascript/rhino/jstype/TemplateTypeReplacer.java @ closure-compiler v20260415

package rhino

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// TemplateTypeReplacer specializes TemplatizedTypes according to provided
// bindings.
//
// Java: implements Visitor<JSType>; the case methods return JSTypes boxed in
// `any` (see Visitor).
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

	// Initialize data structures to `null` because these are unused in ~40% of
	// TemplateTypeReplacers.
	seenTypes          *jutil.LinkedHashSet[JSType]
	visitedObjectTypes *jutil.LinkedHashMap[JSType, JSType]
}

// TemplateTypeReplacer_ForInference creates a replacer for use during
// TypeInference.
//
// Java: TemplateTypeReplacer.forInference(JSTypeRegistry, Map<TemplateType, JSType>)
func TemplateTypeReplacer_ForInference(
	registry *JSTypeRegistry, bindings *jutil.LinkedHashMap[*TemplateType, JSType]) *TemplateTypeReplacer {
	keys := jsTypeCopyOfTemplateTypes(bindings.KeySlice())
	values := []JSType{}
	for _, key := range keys {
		value := bindings.Get(key)
		if value != nil {
			values = append(values, value)
		} else {
			values = append(values, registry.GetNativeType(JSTypeNative_UNKNOWN_TYPE))
		}
	}
	m := registry.GetEmptyTemplateTypeMap().CopyWithExtensionOfKeys(keys, values)
	return NewTemplateTypeReplacer(registry, m, true, true, true)
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

// Java: TemplateTypeReplacer(JSTypeRegistry, TemplateTypeMap, boolean, boolean, boolean)
func NewTemplateTypeReplacer(
	registry *JSTypeRegistry,
	bindings *TemplateTypeMap,
	visitProperties bool,
	useUnknownForMissingKeys bool,
	useUnknownForMissingValues bool) *TemplateTypeReplacer {
	return &TemplateTypeReplacer{
		registry:                   registry,
		bindings:                   bindings,
		visitProperties:            visitProperties,
		useUnknownForMissingKeys:   useUnknownForMissingKeys,
		useUnknownForMissingValues: useUnknownForMissingValues,
		ownSubMapBoundary:          -1,
	}
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
	if r.seenTypes == nil {
		// Sets.newIdentityHashSet()
		r.seenTypes = jutil.NewLinkedHashSet[JSType]()
	}
}

// Java: TemplateTypeReplacer.caseNoType(NoType)
func (r *TemplateTypeReplacer) CaseNoType(type_ *NoType) any {
	return JSType(type_)
}

// Java: TemplateTypeReplacer.caseEnumElementType(EnumElementType)
func (r *TemplateTypeReplacer) CaseEnumElementType(type_ *EnumElementType) any {
	return JSType(type_)
}

// Java: TemplateTypeReplacer.caseAllType()
func (r *TemplateTypeReplacer) CaseAllType() any {
	return r.GetNativeType(JSTypeNative_ALL_TYPE)
}

// Java: TemplateTypeReplacer.caseBooleanType()
func (r *TemplateTypeReplacer) CaseBooleanType() any {
	return r.GetNativeType(JSTypeNative_BOOLEAN_TYPE)
}

// Java: TemplateTypeReplacer.caseNoObjectType()
func (r *TemplateTypeReplacer) CaseNoObjectType() any {
	return r.GetNativeType(JSTypeNative_NO_OBJECT_TYPE)
}

// Java: TemplateTypeReplacer.caseFunctionType(FunctionType)
func (r *TemplateTypeReplacer) CaseFunctionType(type_ FunctionType) any {
	return r.GuardAgainstCycles(type_, func(t JSType) JSType { return r.CaseFunctionTypeUnguarded(t.(FunctionType)) })
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
	afterThis := r.CoerseToThisType(JSType_VisitForJSType(beforeThis, r))
	if beforeThis != afterThis {
		changed = true
	}

	beforeReturn := type_.GetReturnType()
	afterReturn := JSType_VisitForJSType(beforeReturn, r)
	if beforeReturn != afterReturn {
		changed = true
	}

	paramsChanged := false
	numParams := len(type_.GetParameters())
	var paramBuilder *FunctionParamBuilder
	if numParams != 0 {
		paramBuilder = NewFunctionParamBuilderWithCapacity(r.registry, numParams)
	}
	for i := 0; i < numParams; i++ {
		parameter := type_.GetParameters()[i]
		beforeParamType := parameter.GetJSType()
		afterParamType := JSType_VisitForJSType(beforeParamType, r)

		if beforeParamType != afterParamType {
			changed = true
			paramsChanged = true
			// TODO(lharker): we could also lazily create the FunctionParamBuilder
			// here, but that would require re-iterating over all previously seen
			// params so unclear it's worth the complexity
			if parameter.IsOptional() {
				paramBuilder.AddOptionalParams(afterParamType)
			} else if parameter.IsVariadic() {
				paramBuilder.AddVarArgs(afterParamType)
			} else {
				paramBuilder.AddRequiredParams(afterParamType)
			}
		} else {
			paramBuilder.NewParameterFrom(parameter)
		}
	}

	if changed {
		var params []*FunctionType_Parameter
		if paramsChanged {
			params = paramBuilder.Build()
		} else {
			params = type_.GetParameters()
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
func (r *TemplateTypeReplacer) CaseObjectType(objType ObjectType) any {
	return r.GuardAgainstCycles(objType, func(t JSType) JSType { return r.CaseObjectTypeUnguarded(t.(ObjectType)) })
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

	changed := false
	builder := NewRecordTypeBuilder(r.registry)
	for _, prop := range objType.GetOwnPropertyNames() {
		propertyNode := objType.GetPropertyNode(prop)
		beforeType := objType.GetPropertyType(prop)
		afterType := JSType_VisitForJSType(beforeType, r)
		if beforeType != afterType {
			changed = true
		}
		builder.AddProperty(prop, afterType, propertyNode)
	}

	// Use our new type if anything changed, and be sure to update the cache.
	var result JSType = objType
	if changed {
		result = builder.Build()
	}
	r.visitedObjectTypes.Put(objType, result)
	return result
}

// Java: TemplateTypeReplacer.caseTemplatizedType(TemplatizedType)
func (r *TemplateTypeReplacer) CaseTemplatizedType(type_ *TemplatizedType) any {
	return r.GuardAgainstCycles(type_, func(t JSType) JSType { return r.CaseTemplatizedTypeUnguarded(t.(*TemplatizedType)) })
}

// Java: TemplateTypeReplacer.caseTemplatizedTypeUnguarded(TemplatizedType)
func (r *TemplateTypeReplacer) CaseTemplatizedTypeUnguarded(type_ *TemplatizedType) JSType {
	changed := false
	beforeBaseType := type_.GetReferencedType()
	afterBaseType := ObjectType_Cast(JSType_VisitForJSType(beforeBaseType, r))
	if beforeBaseType != afterBaseType {
		changed = true
	}

	builder := []JSType{}
	for _, beforeTemplateType := range type_.GetTemplateTypes() {
		afterTemplateType := JSType_VisitForJSType(beforeTemplateType, r)
		if beforeTemplateType != afterTemplateType {
			changed = true
		}
		builder = append(builder, jutil.CheckNotNull(afterTemplateType))
	}

	if changed {
		type_ = r.registry.CreateTemplatizedTypeOfList(afterBaseType, builder)
	}
	return type_
}

// Java: TemplateTypeReplacer.caseUnknownType()
func (r *TemplateTypeReplacer) CaseUnknownType() any {
	return r.GetNativeType(JSTypeNative_UNKNOWN_TYPE)
}

// Java: TemplateTypeReplacer.caseNullType()
func (r *TemplateTypeReplacer) CaseNullType() any {
	return r.GetNativeType(JSTypeNative_NULL_TYPE)
}

// Java: TemplateTypeReplacer.caseNumberType()
func (r *TemplateTypeReplacer) CaseNumberType() any {
	return r.GetNativeType(JSTypeNative_NUMBER_TYPE)
}

// Java: TemplateTypeReplacer.caseBigIntType()
func (r *TemplateTypeReplacer) CaseBigIntType() any {
	return r.GetNativeType(JSTypeNative_BIGINT_TYPE)
}

// Java: TemplateTypeReplacer.caseStringType()
func (r *TemplateTypeReplacer) CaseStringType() any {
	return r.GetNativeType(JSTypeNative_STRING_TYPE)
}

// Java: TemplateTypeReplacer.caseSymbolType(SymbolType)
func (r *TemplateTypeReplacer) CaseSymbolType(type_ SymbolType) any {
	return JSType(type_)
}

// Java: TemplateTypeReplacer.caseVoidType()
func (r *TemplateTypeReplacer) CaseVoidType() any {
	return r.GetNativeType(JSTypeNative_VOID_TYPE)
}

// Java: TemplateTypeReplacer.caseUnionType(UnionType)
func (r *TemplateTypeReplacer) CaseUnionType(type_ *UnionType) any {
	return r.GuardAgainstCycles(type_, func(t JSType) JSType { return r.CaseUnionTypeUnguarded(t.(*UnionType)) })
}

// Java: TemplateTypeReplacer.caseUnionTypeUnguarded(UnionType)
func (r *TemplateTypeReplacer) CaseUnionTypeUnguarded(type_ *UnionType) JSType {
	changed := false
	results := []JSType{}
	alternates := type_.GetAlternates()
	alternateCount := len(alternates)
	for i := 0; i < alternateCount; i++ {
		alternative := alternates[i]
		replacement := JSType_VisitForJSType(alternative, r)
		if replacement != alternative {
			changed = true
		}
		results = append(results, replacement)
	}

	if changed {
		return r.registry.CreateUnionTypeOfList(results) // maybe not a union
	}

	return type_
}

// Java: TemplateTypeReplacer.caseTemplateType(TemplateType)
func (r *TemplateTypeReplacer) CaseTemplateType(type_ *TemplateType) any {
	r.hasMadeReplacement = true

	if type_ != r.keyType && !r.bindings.HasTemplateKeyFrom(type_, r.ownSubMapBoundary) {
		if r.useUnknownForMissingKeys {
			return r.GetNativeType(JSTypeNative_UNKNOWN_TYPE)
		}
		return JSType(type_)
	}

	r.InitSeenTypes()
	if r.seenTypes.Contains(type_) {
		// If we have already encountered this TemplateType during replacement
		// (i.e. there is a reference loop) then return the TemplateType type
		// itself.
		return JSType(type_)
	} else if !r.bindings.HasTemplateType(type_) {
		// If there is no JSType substitution for the TemplateType, return either
		// the UNKNOWN_TYPE or the TemplateType type itself, depending on
		// configuration.
		if r.useUnknownForMissingValues {
			return r.GetNativeType(JSTypeNative_UNKNOWN_TYPE)
		}
		return JSType(type_)
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
	if replacement == r.keyTypeAsJSType() || r.seenTypes.Contains(replacement) {
		return JSType(type_)
	}
	r.seenTypes.Add(type_)
	visitedReplacement := JSType_VisitForJSType(replacement, r)
	r.seenTypes.Remove(type_)

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
func (r *TemplateTypeReplacer) CaseNamedType(type_ *NamedType) any {
	// The internals of a named type aren't interesting.
	return JSType(type_)
}

// Java: TemplateTypeReplacer.caseProxyObjectType(ProxyObjectType)
func (r *TemplateTypeReplacer) CaseProxyObjectType(type_ ProxyObjectType) any {
	return r.GuardAgainstCycles(type_, func(t JSType) JSType { return r.CaseProxyObjectTypeUnguarded(t.(ProxyObjectType)) })
}

// Java: TemplateTypeReplacer.caseProxyObjectTypeUnguarded(ProxyObjectType)
func (r *TemplateTypeReplacer) CaseProxyObjectTypeUnguarded(type_ ProxyObjectType) JSType {
	// Be careful not to unwrap a type unless it has changed.
	beforeType := type_.GetReferencedTypeInternal()
	replacement := JSType_VisitForJSType(beforeType, r)
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

// Java: TemplateTypeReplacer.guardAgainstCycles(T, Function<T, JSType>)
func (r *TemplateTypeReplacer) GuardAgainstCycles(type_ JSType, mapper func(JSType) JSType) any {
	r.InitSeenTypes()
	if !r.seenTypes.Add(type_) {
		return type_
	}
	defer r.seenTypes.Remove(type_)
	return mapper(type_)
}

// GetVisitedObjectTypeOrNull: determinism is unnecessary for
// visitedObjectTypes, because we only call get/set
//
// Java: TemplateTypeReplacer.getVisitedObjectTypeOrNull(JSType)
func (r *TemplateTypeReplacer) GetVisitedObjectTypeOrNull(type_ JSType) JSType {
	// If we've visited this object type, we don't need to walk it again.
	if r.visitedObjectTypes == nil {
		// new IdentityHashMap<>()
		r.visitedObjectTypes = jutil.NewLinkedHashMap[JSType, JSType]()
		return nil
	}
	return r.visitedObjectTypes.Get(type_)
}
