// Port of com/google/javascript/rhino/jstype/SubtypeChecker.java @ closure-compiler v20260415

package rhino

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// SubtypeChecker represents the computation of a single supertype-subtype
// relationship.
//
// Each instance can only be used once. It acts as a scope for the
// computation and may accumulate state.
type SubtypeChecker struct {
	initialSupertype JSType
	initialSubtype   JSType

	registry *JSTypeRegistry
	// Go: nil is Java null.
	isUsingStructuralTyping *bool
	subtypingMode           *JSType_SubtypingMode

	subtypeCache *jutil.ObjectLinkedHashMap[*SubtypeChecker_CacheKey, JSType_MatchStatus]

	hasRun         bool
	recursionDepth int
}

// Java: SubtypeChecker.BIVARIANT_TYPES
var SubtypeChecker_BIVARIANT_TYPES = jutil.NewLinkedHashSet("Object", "IArrayLike", "Array")

// SubtypeChecker_POTENTIALLY_CYCLIC_RECURSION_DEPTH is an arbitrary depth at
// which to start checking for cyclic recursion.
//
// Checking for cycles is expensive. Most types are not cyclic; they have
// have fairly shallow composition. Therefore, we want to avoid checking
// cycles them until we have a strong expectation that it's necessary.
//
// This number is an optimization, not a correctness requirement. The value
// could be anything (short of stack overflow) and we should get the same
// result.
//
// Java: SubtypeChecker.POTENTIALLY_CYCLIC_RECURSION_DEPTH
const SubtypeChecker_POTENTIALLY_CYCLIC_RECURSION_DEPTH = 20

// Java: SubtypeChecker.setSupertype(JSType)
func (c *SubtypeChecker) SetSupertype(value JSType) *SubtypeChecker {
	c.CheckHasNotRun()
	jutil.CheckState(c.initialSupertype == nil)
	c.initialSupertype = jutil.CheckNotNull(value)
	return c
}

// Java: SubtypeChecker.setSubtype(JSType)
func (c *SubtypeChecker) SetSubtype(value JSType) *SubtypeChecker {
	c.CheckHasNotRun()
	jutil.CheckState(c.initialSubtype == nil)
	c.initialSubtype = jutil.CheckNotNull(value)
	return c
}

// Java: SubtypeChecker.setUsingStructuralSubtyping(boolean)
func (c *SubtypeChecker) SetUsingStructuralSubtyping(value bool) *SubtypeChecker {
	c.CheckHasNotRun()
	jutil.CheckState(c.isUsingStructuralTyping == nil)
	c.isUsingStructuralTyping = &value
	return c
}

// Java: SubtypeChecker.setSubtypingMode(SubtypingMode)
func (c *SubtypeChecker) SetSubtypingMode(value JSType_SubtypingMode) *SubtypeChecker {
	c.CheckHasNotRun()
	jutil.CheckState(c.subtypingMode == nil)
	c.subtypingMode = &value
	return c
}

// Java: SubtypeChecker.checkHasNotRun()
func (c *SubtypeChecker) CheckHasNotRun() {
	jutil.CheckState(!c.hasRun)
}

// Java: SubtypeChecker(JSTypeRegistry)
func NewSubtypeChecker(registry *JSTypeRegistry) *SubtypeChecker {
	return &SubtypeChecker{registry: registry}
}

// Java: SubtypeChecker.check()
func (c *SubtypeChecker) Check() bool {
	c.CheckHasNotRun()
	c.hasRun = true
	return c.IsSubtypeCaching(c.initialSubtype, c.initialSupertype)
}

// IsSubtypeCaching is the top-level recursive entrypoint for subtyping
// logic.
//
// Caching is necessary to catch cyclic types. Identity caching is
// insufficient because some types (e.g. TemplatizedType) can generate new
// type instances on the fly.
//
// Java: SubtypeChecker.isSubtypeCaching(JSType, JSType)
func (c *SubtypeChecker) IsSubtypeCaching(subtype JSType, supertype JSType) bool {
	jutil.CheckNotNull(subtype)
	jutil.CheckNotNull(supertype)

	// Wait to instantiate/use the cache until we have some hint that there may
	// be recursion.
	if c.recursionDepth > SubtypeChecker_POTENTIALLY_CYCLIC_RECURSION_DEPTH {
		if c.subtypeCache == nil {
			c.subtypeCache = jutil.NewObjectLinkedHashMap[*SubtypeChecker_CacheKey, JSType_MatchStatus]()
		}
	}

	// Once the cache exists, use it consistently.
	if c.subtypeCache == nil {
		c.recursionDepth++
		defer func() { c.recursionDepth-- }()
		return c.IsSubtypeDispatching(subtype, supertype)
	}

	key := NewSubtypeChecker_CacheKey(subtype, supertype)
	cached, present := c.subtypeCache.Lookup(key)
	if !present {
		c.subtypeCache.Put(key, JSType_MatchStatus_PROCESSING)
		result := c.IsSubtypeDispatching(subtype, supertype)
		c.subtypeCache.Put(key, JSType_MatchStatus_ValueOf(result))
		return result
	}

	if cached == JSType_MatchStatus_PROCESSING {
		c.subtypeCache.Put(key, JSType_MatchStatus_MATCH)
		return true
	}

	return cached.SubtypeValue()
}

// IsSubtypeDispatching is the custom dynamic dispatcher for various
// subtypes.
//
// This method should be temporary. It exists to delegate to subclass
// specifc logic that used to be distributed across all the overrides of
// JSType::isSubtype. In order to consolidate all the behaviour, a limited
// form of dynamic dispatch was created here. Over time, this behaviour
// should be integrated into the main body of subtyping logic.
//
// Java: SubtypeChecker.isSubtypeDispatching(JSType, JSType)
func (c *SubtypeChecker) IsSubtypeDispatching(subtype JSType, supertype JSType) bool {
	switch subtype.GetTypeClass() {
	case JSTypeClass_ARROW:
		return c.IsArrowTypeSubtype(subtype.(*ArrowType), supertype)
	case JSTypeClass_ENUM_ELEMENT:
		return c.IsEnumElementSubtype(subtype.(*EnumElementType), supertype)
	case JSTypeClass_ENUM:
		return c.IsEnumSubtype(subtype.(*EnumType), supertype)
	case JSTypeClass_NO_OBJECT, JSTypeClass_NO:
		return c.IsVariousBottomsSubtype(subtype, supertype)
	case JSTypeClass_FUNCTION:
		return c.IsFunctionSubtype(subtype.(FunctionType), supertype)
	case JSTypeClass_TEMPLATE:
		return c.IsTemplateSubtype(subtype.(*TemplateType), supertype)
	case JSTypeClass_PROXY_OBJECT:
		return c.IsProxyObjectSubtype(subtype.(ProxyObjectType), supertype)
	case JSTypeClass_WELL_KNOWN_SYMBOL:
		return c.IsWellKnownSymbolSubtype(subtype.(*KnownSymbolType), supertype)
	}
	return c.IsSubtypeHelper(subtype, supertype)
}

// Java: SubtypeChecker.isWellKnownSymbolSubtype(KnownSymbolType, JSType)
func (c *SubtypeChecker) IsWellKnownSymbolSubtype(subtype *KnownSymbolType, supertype JSType) bool {
	if supertype.IsKnownSymbolValueType() {
		return subtype.Equals(supertype)
	}
	if supertype.IsSymbolValueType() {
		return true
	}
	return c.IsSubtypeHelper(subtype, supertype)
}

// Java: SubtypeChecker.isSubtypeHelper(JSType, JSType)
func (c *SubtypeChecker) IsSubtypeHelper(subtype JSType, supertype JSType) bool {
	// Axiomatic cases.
	if subtype == supertype ||
		supertype.IsUnknownType() ||
		supertype.IsAllType() ||
		subtype.IsUnknownType() ||
		subtype.IsNoType() {
		return true
	}

	// Special consideration for `null` and `undefined` in J2CL.
	if c.subtypingMode != nil && *c.subtypingMode == JSType_SubtypingMode_IGNORE_NULL_UNDEFINED &&
		(subtype.IsNullType() || subtype.IsVoidType()) {
		return true
	}

	// Reflexive case.
	if jutil.ObjectsEquals(subtype, supertype) {
		return true
	}

	// Unwrap proxy types.
	//
	// Only named types are unwrapped because other subclasses of
	// `ProxyObjectType` should not be considered proxies; they have additional
	// behaviour.
	//
	// We don't want to check the cache here. Since `NamedType`s are generally
	// equal to their referenced types, we'd get a false cache hit.
	if subtype.IsNamedType() {
		return c.IsSubtypeDispatching(subtype.ToMaybeNamedType().GetReferencedType(), supertype)
	} else if supertype.IsNamedType() {
		return c.IsSubtypeDispatching(subtype, supertype.ToMaybeNamedType().GetReferencedType())
	}

	// Union decomposition.
	if subtype.IsUnionType() {
		// All alternates must be subtypes.
		return JSTypeIterations_AllTypesMatchInUnion(
			func(sub JSType) bool { return c.IsSubtypeCaching(sub, supertype) }, subtype.ToMaybeUnionType())
	} else if supertype.IsUnionType() {
		// Some alternate must be a supertype.
		return JSTypeIterations_AnyTypeMatchesInUnion(
			func(sup JSType) bool { return c.IsSubtypeCaching(subtype, sup) }, supertype.ToMaybeUnionType())
	}

	if !subtype.IsObjectType() || !supertype.IsObjectType() {
		return false // All cases for non object types have been covered by this point.
	}

	return c.IsObjectSubtypeHelper(subtype.AssertObjectType(), supertype.AssertObjectType())
}

// Java: SubtypeChecker.isObjectSubtypeHelper(ObjectType, ObjectType)
func (c *SubtypeChecker) IsObjectSubtypeHelper(subtype ObjectType, supertype ObjectType) bool {
	subtypeParams := subtype.GetTemplateTypeMap()
	supertypeParams := supertype.GetTemplateTypeMap()
	subAndSuperAreSameBaseType := false
	if subtype.IsTemplatizedType() && supertype.IsTemplatizedType() {
		templatizedSubType := subtype.ToMaybeTemplatizedType()
		templatizedSuperType := supertype.ToMaybeTemplatizedType()
		superTypeReferencedType := templatizedSuperType.GetReferencedType()
		subTypeReferencedType := templatizedSubType.GetReferencedType()
		if subTypeReferencedType.Equals(superTypeReferencedType) {
			subAndSuperAreSameBaseType = true
		}
	}
	bivarantMatch := false

	// Array and Object are exempt from template type invariance.
	//
	// They also have to be checked first because the `Object` index key acts
	// like an operator; it's not visible as a property but it has typing.
	// There are types that would otherwise match structurally but should not
	// be be considered subtypes because their types for this operator are
	// different.
	if SubtypeChecker_IsBivariantType(supertype) {
		key := subtype.GetJSTypeBase().registry.GetObjectElementKey()
		thisElement := subtypeParams.GetResolvedTemplateType(key)
		thatElement := supertypeParams.GetResolvedTemplateType(key)

		if !c.MeetVarianceConstraint(SubtypeChecker_Variance_BIVARIANT, thisElement, thatElement) {
			return false
		}
		bivarantMatch = true
	}
	var covariantKey *TemplateType
	if !bivarantMatch {
		covariantKey = SubtypeChecker_GetTemplateKeyIfCovariantType(supertype)
	}

	if !subAndSuperAreSameBaseType &&
		*c.isUsingStructuralTyping &&
		supertype.IsStructuralType() {
		// Do this before considering templatization in general.
		//
		// If the super type is a structural type, then we can't safely unwrap
		// any templatized types. The templates might affect the types of the
		// properties.
		if covariantKey == c.registry.GetReadonlyArrayElementKey() &&
			c.IsNativeArrayType(subtype) {
			// patch for covariant subtyping of Array and ReadonlyArray. structural
			// subtyping generally doesn't really support covariance of generics,
			// but Array and ReadonlyArray are a special case we want to support as
			// it has been a problem in practice.
			thisElement := subtypeParams.GetResolvedTemplateType(covariantKey)
			thatElement := supertypeParams.GetResolvedTemplateType(covariantKey)
			return c.MeetVarianceConstraint(SubtypeChecker_Variance_COVARIANT, thisElement, thatElement)
		}
		return c.IsStructuralSubtypeHelper(
			subtype, supertype, SubtypeChecker_PropertyOptionality_VOIDABLE_PROPS_ARE_OPTIONAL)
	}

	// Wait to check template types until after structural checks.
	//
	// It's possible for a subtructural type to satisfy the shape defined by a
	// template specialization, without actually having template params. The
	// structural type may just have had properties with the right types.
	if !bivarantMatch {
		if covariantKey != nil {
			thisElement := subtypeParams.GetResolvedTemplateType(covariantKey)
			thatElement := supertypeParams.GetResolvedTemplateType(covariantKey)
			if !c.MeetVarianceConstraint(SubtypeChecker_Variance_COVARIANT, thisElement, thatElement) {
				return false
			}
		} else if !c.IsTypeMapSubmap(subtype, supertype) {
			return false
		}
	}

	// Interfaces
	// Find all the interfaces implemented by this class and compare each one
	// to the interface instance.
	subtypeCtor := subtype.GetConstructor()
	supertypeCtor := supertype.GetConstructor()
	if subtypeCtor != nil && subtypeCtor.IsInterface() {
		for _, subtypeInterface := range subtype.GetCtorExtendedInterfaces() {
			if c.IsSubtypeCaching(subtypeInterface, supertype) {
				return true
			}
		}
	} else if supertypeCtor != nil && supertypeCtor.IsInterface() {
		for _, subtypeInterface := range subtype.GetCtorImplementedInterfaces() {
			if c.IsSubtypeCaching(subtypeInterface, supertype) {
				return true
			}
		}
	}

	return supertype.IsImplicitPrototypeOf(subtype)
}

// Java: SubtypeChecker.isStructuralSubtypeHelper(ObjectType, ObjectType, PropertyOptionality)
func (c *SubtypeChecker) IsStructuralSubtypeHelper(
	subtype ObjectType, supertype ObjectType, optionality SubtypeChecker_PropertyOptionality) bool {

	// subtype is a subtype of record type supertype iff:
	// 1) subtype has all the non-optional properties declared in supertype.
	// 2) And for each property of supertype, its type must be
	//    a super type of the corresponding property of subtype.

	var keys *PropertyMap_AllKeys
	if !supertype.IsRecordType() {
		keys = supertype.GetPropertyMap().GetAllKeys()
	}
	// NOTE: Inline record literal types always have Object as a supertype. In
	// these cases, we really only care about the properties explicitly
	// declared in the record literal, and not about any properties inherited
	// from Object.prototype. On the other hand, @record types allow
	// inheritance and we need to match against inherited properties as well.
	var props []string
	if supertype.IsRecordType() {
		props = supertype.GetOwnPropertyNames()
	} else {
		props = keys.StringKeys()
	}

	for _, property := range props {
		if !c.CheckPropertyPresence(
			NewProperty_StringKey(property), subtype, supertype, optionality) {
			return false
		}
	}
	var knownSymbols []*KnownSymbolType
	if supertype.IsRecordType() {
		knownSymbols = supertype.GetOwnPropertyKnownSymbols()
	} else {
		knownSymbols = keys.KnownSymbolKeys()
	}

	for _, knownSymbolType := range knownSymbols {
		knownSymbolTypeKey := NewProperty_SymbolKey(knownSymbolType)
		if !c.CheckPropertyPresence(knownSymbolTypeKey, subtype, supertype, optionality) {
			return false
		}
	}

	return true
}

// Java: SubtypeChecker.checkPropertyPresence(Property.Key, ObjectType, ObjectType, PropertyOptionality)
func (c *SubtypeChecker) CheckPropertyPresence(
	propertyKey Property_Key,
	subtype ObjectType,
	supertype ObjectType,
	optionality SubtypeChecker_PropertyOptionality) bool {

	supertypeProp := supertype.GetPropertyTypeByKey(propertyKey)
	if subtype.HasPropertyByKey(propertyKey) {
		subtypeProp := subtype.GetPropertyTypeByKey(propertyKey)
		if !c.IsSubtypeCaching(subtypeProp, supertypeProp) {
			return false
		}
	} else if !optionality.IsOptional(supertypeProp) {
		return false
	}
	return true
}

// Java: SubtypeChecker.isNativeArrayType(ObjectType)
func (c *SubtypeChecker) IsNativeArrayType(type_ ObjectType) bool {
	unwrapped := SubtypeChecker_GetObjectTypeIfNative(type_)
	return unwrapped != nil && unwrapped.Equals(c.registry.GetNativeType(JSTypeNative_ARRAY_TYPE))
}

// Java: SubtypeChecker.isFunctionSubtype(FunctionType, JSType)
func (c *SubtypeChecker) IsFunctionSubtype(subtype FunctionType, nonFunctionSupertype JSType) bool {
	if c.IsSubtypeHelper(subtype, nonFunctionSupertype) {
		return true
	}

	if nonFunctionSupertype.IsFunctionType() {
		supertype := nonFunctionSupertype.ToMaybeFunctionType()
		if supertype.IsInterface() {
			// Any function can be assigned to an interface function.
			return true
		}
		if subtype.IsInterface() {
			// An interface function cannot be assigned to anything.
			return false
		}

		return c.ShouldTreatThisTypesAsCovariant(subtype, supertype) &&
			c.IsSubtypeCaching(
				subtype.GetInternalArrowType(), supertype.GetInternalArrowType())
	}

	return c.IsSubtypeCaching(
		c.registry.GetNativeType(JSTypeNative_FUNCTION_PROTOTYPE), nonFunctionSupertype)
}

// Java: SubtypeChecker.isVariousBottomsSubtype(JSType, JSType)
func (c *SubtypeChecker) IsVariousBottomsSubtype(subtype JSType, supertype JSType) bool {
	if c.IsSubtypeHelper(subtype, supertype) {
		return true
	}

	return supertype.IsObject() && !supertype.IsNoType()
}

// Java: SubtypeChecker.isArrowTypeSubtype(ArrowType, JSType)
func (c *SubtypeChecker) IsArrowTypeSubtype(subtype *ArrowType, nonArrowSupertype JSType) bool {
	supertype, ok := nonArrowSupertype.(*ArrowType)
	if !ok {
		return false
	}

	// This is described in Draft 2 of the ES4 spec, Section 3.4.7: Subtyping
	// Function Types.

	// this.returnType <: that.returnType (covariant)
	if !c.IsSubtypeCaching(subtype.GetReturnType(), supertype.GetReturnType()) {
		return false
	}

	// that.paramType[i] <: this.paramType[i] (contravariant)
	//
	// If this.paramType[i] is required, then that.paramType[i] is required.
	//
	// In theory, the "required-ness" should work in the other direction as
	// well. In other words, if we have
	//
	// function f(number, number) {}
	// function g(number) {}
	//
	// Then f *should* not be a subtype of g, and g *should* not be a subtype
	// of f. But in practice, we do not implement it this way. We want to
	// support the use case where you can pass g where f is expected, and
	// pretend that g ignores the second argument. That way, you can have a
	// single "no-op" function, and you don't have to create a new no-op
	// function for every possible type signature.
	//
	// So, in this case, g < f, but f !< g
	subtypeParameters := subtype.GetParameterList()
	supertypeParameters := supertype.GetParameterList()
	si, pi := 0, 0
	var subtypeParam, supertypeParam *FunctionType_Parameter
	if si < len(subtypeParameters) {
		subtypeParam = subtypeParameters[si]
		si++
	}
	if pi < len(supertypeParameters) {
		supertypeParam = supertypeParameters[pi]
		pi++
	}
	for subtypeParam != nil && supertypeParam != nil {

		subtypeParamType := subtypeParam.GetJSType()
		supertypeParamType := supertypeParam.GetJSType()
		if subtypeParamType != nil {
			if supertypeParamType == nil ||
				!c.IsSubtypeCaching(supertypeParamType, subtypeParamType) {
				return false
			}
		}

		thisIsVarArgs := subtypeParam.IsVariadic()
		thatIsVarArgs := supertypeParam.IsVariadic()
		thisIsOptional := thisIsVarArgs || subtypeParam.IsOptional()
		thatIsOptional := thatIsVarArgs || supertypeParam.IsOptional()

		// "that" can't be a supertype, because it's missing a required argument.
		if !thisIsOptional && thatIsOptional {
			// NOTE(nicksantos): In our type system, we use {function(...?)} and
			// {function(...NoType)} to to indicate that arity should not be
			// checked. Strictly speaking, this is not a correct formulation,
			// because now a sub-function can required arguments that are var_args
			// in the super-function. So we special-case this.
			isTopFunction := thatIsVarArgs &&
				(supertypeParamType == nil ||
					supertypeParamType.IsUnknownType() ||
					supertypeParamType.IsNoType())
			if !isTopFunction {
				return false
			}
		}

		// don't advance if we have variable arguments
		if !thisIsVarArgs {
			subtypeParam = nil
			if si < len(subtypeParameters) {
				subtypeParam = subtypeParameters[si]
				si++
			}
		}
		if !thatIsVarArgs {
			supertypeParam = nil
			if pi < len(supertypeParameters) {
				supertypeParam = supertypeParameters[pi]
				pi++
			}
		}

		// both var_args indicates the end
		if thisIsVarArgs && thatIsVarArgs {
			break
		}
	}

	// "that" can't be a supertype, because it's missing a required argument.
	return subtypeParam == nil ||
		subtypeParam.IsOptional() ||
		subtypeParam.IsVariadic() ||
		supertypeParam != nil
}

// Java: SubtypeChecker.isEnumSubtype(EnumType, JSType)
func (c *SubtypeChecker) IsEnumSubtype(subtype *EnumType, supertype JSType) bool {
	return supertype.Equals(c.registry.GetNativeType(JSTypeNative_OBJECT_TYPE)) ||
		subtype.Equals(c.registry.GetNativeType(JSTypeNative_OBJECT_PROTOTYPE)) ||
		c.IsSubtypeHelper(subtype, supertype)
}

// Java: SubtypeChecker.isEnumElementSubtype(EnumElementType, JSType)
func (c *SubtypeChecker) IsEnumElementSubtype(subtype *EnumElementType, supertype JSType) bool {
	// TODO(b/136298690): stop creating such a 'meet' and remove this logic.
	// ** start hack **
	if supertype.IsEnumElementType() &&
		subtype.GetEnumType() == supertype.ToMaybeEnumElementType().GetEnumType() {
		// This can happen if, e.g., we have the 'meet' of enum elements
		// Foo<number> and Bar<number>, which is Foo<Bar<number>>; and are
		// comparing it to Foo<number>.
		return c.IsSubtypeCaching(
			subtype.GetPrimitiveType(), supertype.ToMaybeEnumElementType().GetPrimitiveType())
	}
	// ** end hack **

	if c.IsSubtypeHelper(subtype, supertype) {
		return true
	} else {
		return c.IsSubtypeCaching(subtype.GetPrimitiveType(), supertype)
	}
}

// Java: SubtypeChecker.isProxyObjectSubtype(ProxyObjectType, JSType)
func (c *SubtypeChecker) IsProxyObjectSubtype(subtype ProxyObjectType, supertype JSType) bool {
	// Don't check the cache here.
	//
	// If we do, we get false positives for recursion because proxy types are
	// equal to their referenced types. Fortunately, we can leverage this
	// guarantee of equlity to just check subtyping directly on the referenced
	// type.
	return c.IsSubtypeDispatching(subtype.GetReferencedTypeInternal(), supertype)
}

// Java: SubtypeChecker.isTemplateSubtype(TemplateType, JSType)
func (c *SubtypeChecker) IsTemplateSubtype(subtype *TemplateType, supertype JSType) bool {
	if !subtype.GetBound().IsUnknownType() &&
		supertype.IsTemplateType() &&
		!supertype.ToMaybeTemplateType().GetBound().IsUnknownType() {
		return subtype.Visit(NewContainsUpperBoundSuperTypeVisitor(supertype)) ==
			ContainsUpperBoundSuperTypeVisitor_Result_PRESENT
	} else {
		return c.IsProxyObjectSubtype(subtype, supertype)
	}
}

// Java: SubtypeChecker.isTypeMapSubmap(ObjectType, ObjectType)
func (c *SubtypeChecker) IsTypeMapSubmap(subtype ObjectType, supertype ObjectType) bool {
	if subtype.IsFunctionPrototypeType() {
		// TODO(b/145609962): Prototype types never have template type maps, so
		// the result is predetermined. We assume true for legacy reasons.
		// TODO(b/145610392): It's unclear what the correct result should be.
		return true
	}

	submap := subtype.GetTemplateTypeMap()
	supermap := supertype.GetTemplateTypeMap()

	// We only need to iterate the keys of the supermap.
	//
	// A submap may have additional entries not present in the supermap, so
	// long as it also has all supermap entries.
	for _, key := range supermap.GetTemplateKeys() {
		keyCount := submap.GetTemplateKeyCountThisShouldAlwaysBeOneOrZeroButIsnt(key)
		if keyCount == 0 {
			if subtype.LoosenTypecheckingDueToForwardReferencedSupertype() {
				// TODO(b/145145406): Delete this case.
				continue
			}
			return false
		} else if keyCount > 1 {
			// TOOD(b/139230800): Delete this case.
			continue
		}

		subvalue := submap.GetResolvedTemplateType(key)
		supervalue := supermap.GetResolvedTemplateType(key)
		if !c.MeetVarianceConstraint(SubtypeChecker_Variance_INVARIANT, subvalue, supervalue) {
			return false
		}
	}

	return true
}

// Java: SubtypeChecker.meetVarianceConstraint(Variance, JSType, JSType)
func (c *SubtypeChecker) MeetVarianceConstraint(variance SubtypeChecker_Variance, override JSType, reference JSType) bool {
	switch variance {
	case SubtypeChecker_Variance_COVARIANT:
		return c.IsSubtypeCaching(override, reference)
	case SubtypeChecker_Variance_CONTRAVARIANT:
		return c.IsSubtypeCaching(reference, override)
	case SubtypeChecker_Variance_BIVARIANT:
		return c.MeetVarianceConstraint(SubtypeChecker_Variance_COVARIANT, reference, override) ||
			c.MeetVarianceConstraint(SubtypeChecker_Variance_CONTRAVARIANT, reference, override)
	}
	// INVARIANT
	return c.MeetVarianceConstraint(SubtypeChecker_Variance_COVARIANT, reference, override) &&
		c.MeetVarianceConstraint(SubtypeChecker_Variance_CONTRAVARIANT, reference, override)
}

// SubtypeChecker_IsBivariantType determines if the supplied type should be
// checked as a bivariant templatized type rather the standard invariant
// templatized type rules.
//
// Java: SubtypeChecker.isBivariantType(JSType)
func SubtypeChecker_IsBivariantType(type_ JSType) bool {
	unwrapped := SubtypeChecker_GetObjectTypeIfNative(type_)
	if unwrapped == nil {
		return false
	}
	name := unwrapped.GetReferenceName()
	return name != nil && SubtypeChecker_BIVARIANT_TYPES.Contains(*name)
}

// SubtypeChecker_GetTemplateKeyIfCovariantType determines if the specified
// type should be checked as covariant rather than the standard invariant
// type. If so, returns the template type to check covariantly.
//
// Java: SubtypeChecker.getTemplateKeyIfCovariantType(JSType)
func SubtypeChecker_GetTemplateKeyIfCovariantType(type_ JSType) *TemplateType {
	if type_.IsTemplatizedType() {
		// Unlike other covariant/bivariant types, even non-native subtypes of
		// IThenable are covariant, so IThenable is special-cased here.
		ttype := type_.ToMaybeTemplatizedType()
		if ttype.GetTemplateTypeMap().HasTemplateKey(ttype.registry.GetIThenableTemplate()) {
			return ttype.registry.GetIThenableTemplate()
		}
	}
	unwrapped := SubtypeChecker_GetObjectTypeIfNative(type_)
	var unwrappedTypeName *string
	if unwrapped != nil {
		unwrappedTypeName = unwrapped.GetReferenceName()
	}
	if unwrappedTypeName == nil {
		return nil
	}
	registry := unwrapped.GetJSTypeBase().registry
	switch *unwrappedTypeName {
	case "ReadonlyArray":
		return registry.GetReadonlyArrayElementKey()
	case "Iterator":
		return registry.GetIteratorValueTemplate()
	case "Generator":
		return registry.GetGeneratorValueTemplate()
	case "AsyncIterator":
		return registry.GetAsyncIteratorValueTemplate()
	case "Iterable":
		return registry.GetIterableValueTemplate()
	case "IteratorIterable":
		return registry.GetIteratorIterableValueTemplate()
	case "IIterableResult":
		return registry.GetIIterableResultValueTemplate()
	case "AsyncIterable":
		return registry.GetAsyncIterableValueTemplate()
	}
	// All other types are either invariant or bivariant
	return nil
}

// Java: SubtypeChecker.shouldTreatThisTypesAsCovariant(FunctionType, FunctionType)
func (c *SubtypeChecker) ShouldTreatThisTypesAsCovariant(subtype FunctionType, supertype FunctionType) bool {
	// If functionA is a subtype of functionB, then their "this" types should
	// be contravariant. However, this causes problems because of the way we
	// enforce overrides. Because function(this:SubFoo) is not a subtype of
	// function(this:Foo), our override check treats this as an error. Let's
	// punt on all this for now.

	// TODO(nicksantos): fix this.
	// An interface 'this'-type is non-restrictive. In practical terms, if C
	// implements I, and I has a method m, then any m doesn't necessarily have
	// to C#m's 'this' type doesn't need to match I.
	if supertype.GetTypeOfThis().ToObjectType() != nil &&
		supertype.GetTypeOfThis().ToObjectType().GetConstructor() != nil &&
		supertype.GetTypeOfThis().ToObjectType().GetConstructor().IsInterface() {
		return true
	}

	// If one of the 'this' types is covariant of the other, then we'll treat
	// them as covariant (see comment above).
	return c.HackTemporarilyChangeSubtypingMode(
		JSType_SubtypingMode_NORMAL,
		func() bool {
			return c.IsSubtypeCaching(supertype.GetTypeOfThis(), subtype.GetTypeOfThis()) ||
				c.IsSubtypeCaching(subtype.GetTypeOfThis(), supertype.GetTypeOfThis())
		})
}

// Java: SubtypeChecker.hackTemporarilyChangeSubtypingMode(SubtypingMode, BooleanSupplier)
func (c *SubtypeChecker) HackTemporarilyChangeSubtypingMode(tempMode JSType_SubtypingMode, callback func() bool) bool {
	originalMode := c.subtypingMode
	defer func() { c.subtypingMode = originalMode }()
	c.subtypingMode = &tempMode
	return callback()
}

// Java: SubtypeChecker.getObjectTypeIfNative(JSType)
func SubtypeChecker_GetObjectTypeIfNative(type_ JSType) ObjectType {
	objType := type_.ToObjectType()
	unwrapped := ObjectType_DeeplyUnwrap(objType)
	if unwrapped != nil && unwrapped.IsNativeObjectType() {
		return unwrapped
	}
	return nil
}

// SubtypeChecker_PropertyOptionality is how to treat explicitly voidable
// properties for structural subtype checking.
type SubtypeChecker_PropertyOptionality int32

const (
	// Explicitly voidable properties are treated as optional.
	SubtypeChecker_PropertyOptionality_VOIDABLE_PROPS_ARE_OPTIONAL SubtypeChecker_PropertyOptionality = iota
	// All properties are always required, even if explicitly voidable.
	SubtypeChecker_PropertyOptionality_ALL_PROPS_ARE_REQUIRED
)

// Java: SubtypeChecker.PropertyOptionality.isOptional(JSType)
func (o SubtypeChecker_PropertyOptionality) IsOptional(propType JSType) bool {
	return o == SubtypeChecker_PropertyOptionality_VOIDABLE_PROPS_ARE_OPTIONAL && propType.IsExplicitlyVoidable()
}

// SubtypeChecker_Variance is SubtypeChecker.Variance.
type SubtypeChecker_Variance int32

const (
	SubtypeChecker_Variance_COVARIANT SubtypeChecker_Variance = iota
	SubtypeChecker_Variance_CONTRAVARIANT
	SubtypeChecker_Variance_BIVARIANT
	SubtypeChecker_Variance_INVARIANT
)

// SubtypeChecker_CacheKey is SubtypeChecker.CacheKey.
type SubtypeChecker_CacheKey struct {
	left     JSType
	right    JSType
	hashCode int32 // Cache this calculation because it is made often.
}

// Java: SubtypeChecker.CacheKey.hashCode()
func (k *SubtypeChecker_CacheKey) HashCode() int32 {
	return k.hashCode
}

// Java: SubtypeChecker.CacheKey.equals(Object)
func (k *SubtypeChecker_CacheKey) Equals(other any) bool {
	// Calling this with `null` or not a `Key` should never happen, so it's
	// fine to crash.
	that := other.(*SubtypeChecker_CacheKey)
	if k.left == that.left && k.right == that.right {
		return true
	}

	// The vast majority of cases will have already returned by now, since
	// equality isn't even checked unless the hash code matches, and in most
	// cases there's only one instance of any equivalent JSType floating
	// around. The remainder only occurs for cyclic (or otherwise complicated)
	// data structures where equivalent types are being synthesized by
	// recursive application of type parameters, or (even more rarely) for hash
	// collisions. Identity is insufficient in the case of recursive
	// parameterized types because new equivalent copies of the type are
	// generated on-the-fly to compute the types of various properties.
	return jutil.ObjectsEquals(k.left, that.left) && jutil.ObjectsEquals(k.right, that.right)
}

// Java: SubtypeChecker.CacheKey(JSType, JSType)
func NewSubtypeChecker_CacheKey(left JSType, right JSType) *SubtypeChecker_CacheKey {
	// NOTE: order matters here, since we're expressing an asymmetric
	// relationship.
	return &SubtypeChecker_CacheKey{left: left, right: right, hashCode: 31*left.HashCode() + right.HashCode()}
}
