// Port of com/google/javascript/rhino/jstype/SubtypeChecker.java
// @ closure-compiler v20260415
//
// Represents the computation of a single supertype-subtype relationship.
// Each instance can only be used once. It acts as a scope for the computation
// and may accumulate state.
//
// NOTE: Types referenced by the Java original that are not yet ported to Go
// (ArrowType, ProxyObjectType, NamedType, RecordType, KnownSymbolType,
// ContainsUpperBoundSuperTypeVisitor, PropertyMap, Property, TemplateTypeMap
// full API) are stubbed with TODO markers. The public API (SubtypeChecker,
// NewSubtypeChecker, SubtypingMode) is complete and callable.

package jstype

// SubtypingMode mirrors JSType.SubtypingMode in Java — controls how null
// and undefined are treated during subtype checks.
type SubtypingMode int

const (
	// SubtypingModeNormal is the standard mode; null and void are not
	// auto-subtypes of everything.
	SubtypingModeNormal SubtypingMode = iota
	// SubtypingModeIgnoreNullUndefined is used by J2CL; null and void
	// are treated as subtypes of all types.
	SubtypingModeIgnoreNullUndefined
)

// matchStatus mirrors JSType.MatchStatus — used in the cycle-detection cache.
type matchStatus int

const (
	matchStatusProcessing matchStatus = iota
	matchStatusMatch
	matchStatusNoMatch
)

func (ms matchStatus) subtypeValue() bool { return ms == matchStatusMatch || ms == matchStatusProcessing }

func matchStatusValueOf(result bool) matchStatus {
	if result {
		return matchStatusMatch
	}
	return matchStatusNoMatch
}

// subtypeCacheKey is the cache key for SubtypeChecker's cycle-detection map.
// Mirrors SubtypeChecker.CacheKey in Java.
type subtypeCacheKey struct {
	left  JSType
	right JSType
}

// potentiallyCyclicRecursionDepth is an arbitrary depth at which to start
// checking for cyclic recursion. Mirrors SubtypeChecker.POTENTIALLY_CYCLIC_RECURSION_DEPTH.
const potentiallyCyclicRecursionDepth = 20

// bivariantTypes mirrors SubtypeChecker.BIVARIANT_TYPES. Array and Object are
// exempt from template type invariance.
var bivariantTypes = map[string]bool{
	"Object":      true,
	"IArrayLike":  true,
	"Array":       true,
}

// subtypeVariance mirrors SubtypeChecker.Variance.
type subtypeVariance int

const (
	varianceCovariant     subtypeVariance = iota
	varianceContravariant subtypeVariance = iota
	varianceBivariant     subtypeVariance = iota
	varianceInvariant     subtypeVariance = iota
)

// propertyOptionality mirrors SubtypeChecker.PropertyOptionality.
type propertyOptionality int

const (
	// voidablePropsAreOptional: explicitly voidable properties are treated as optional.
	voidablePropsAreOptional propertyOptionality = iota
	// allPropsAreRequired: all properties are always required, even if explicitly voidable.
	allPropsAreRequired propertyOptionality = iota
)

func (po propertyOptionality) isOptional(propType JSType) bool {
	return po == voidablePropsAreOptional && containsVoid(propType)
}

// SubtypeChecker represents the computation of a single supertype-subtype
// relationship. Each instance can only be used once. Mirrors
// SubtypeChecker in Java.
type SubtypeChecker struct {
	initialSupertype JSType
	initialSubtype   JSType

	registry              *Registry
	isUsingStructuralTyping *bool
	subtypingMode           SubtypingMode

	subtypeCache   map[subtypeCacheKey]matchStatus
	hasRun         bool
	recursionDepth int
}

// NewSubtypeChecker allocates a SubtypeChecker bound to the given registry.
// Mirrors SubtypeChecker(JSTypeRegistry) constructor in Java.
func NewSubtypeChecker(registry *Registry) *SubtypeChecker {
	return &SubtypeChecker{registry: registry}
}

// SetSupertype sets the supertype for the check. Must be called before Check.
// Mirrors SubtypeChecker.setSupertype in Java.
func (sc *SubtypeChecker) SetSupertype(value JSType) *SubtypeChecker {
	sc.checkHasNotRun()
	if sc.initialSupertype != nil {
		panic("initialSupertype already set")
	}
	if value == nil {
		panic("supertype must not be nil")
	}
	sc.initialSupertype = value
	return sc
}

// SetSubtype sets the subtype for the check. Must be called before Check.
// Mirrors SubtypeChecker.setSubtype in Java.
func (sc *SubtypeChecker) SetSubtype(value JSType) *SubtypeChecker {
	sc.checkHasNotRun()
	if sc.initialSubtype != nil {
		panic("initialSubtype already set")
	}
	if value == nil {
		panic("subtype must not be nil")
	}
	sc.initialSubtype = value
	return sc
}

// SetUsingStructuralSubtyping configures whether structural subtyping is used.
// Mirrors SubtypeChecker.setUsingStructuralSubtyping in Java.
func (sc *SubtypeChecker) SetUsingStructuralSubtyping(value bool) *SubtypeChecker {
	sc.checkHasNotRun()
	if sc.isUsingStructuralTyping != nil {
		panic("isUsingStructuralTyping already set")
	}
	sc.isUsingStructuralTyping = &value
	return sc
}

// SetSubtypingMode sets the subtyping mode (NORMAL or IGNORE_NULL_UNDEFINED).
// Mirrors SubtypeChecker.setSubtypingMode in Java.
func (sc *SubtypeChecker) SetSubtypingMode(value SubtypingMode) *SubtypeChecker {
	sc.checkHasNotRun()
	sc.subtypingMode = value
	return sc
}

func (sc *SubtypeChecker) checkHasNotRun() {
	if sc.hasRun {
		panic("SubtypeChecker has already run")
	}
}

func (sc *SubtypeChecker) useStructuralTyping() bool {
	return sc.isUsingStructuralTyping != nil && *sc.isUsingStructuralTyping
}

// Check runs the subtype check. Can only be called once per instance.
// Mirrors SubtypeChecker.check in Java.
func (sc *SubtypeChecker) Check() bool {
	sc.checkHasNotRun()
	sc.hasRun = true
	return sc.isSubtypeCaching(sc.initialSubtype, sc.initialSupertype)
}

// isSubtypeCaching is the top-level recursive entrypoint for subtyping logic.
// Caching is necessary to catch cyclic types. Mirrors
// SubtypeChecker.isSubtypeCaching in Java.
func (sc *SubtypeChecker) isSubtypeCaching(subtype, supertype JSType) bool {
	if subtype == nil {
		panic("subtype must not be nil")
	}
	if supertype == nil {
		panic("supertype must not be nil")
	}

	// Wait to instantiate/use the cache until we have some hint of recursion.
	if sc.recursionDepth > potentiallyCyclicRecursionDepth {
		if sc.subtypeCache == nil {
			sc.subtypeCache = make(map[subtypeCacheKey]matchStatus)
		}
	}

	// Once the cache exists, use it consistently.
	if sc.subtypeCache == nil {
		sc.recursionDepth++
		result := sc.isSubtypeDispatching(subtype, supertype)
		sc.recursionDepth--
		return result
	}

	key := subtypeCacheKey{left: subtype, right: supertype}
	cached, exists := sc.subtypeCache[key]

	if !exists {
		sc.subtypeCache[key] = matchStatusProcessing
		result := sc.isSubtypeDispatching(subtype, supertype)
		sc.subtypeCache[key] = matchStatusValueOf(result)
		return result
	}

	if cached == matchStatusProcessing {
		sc.subtypeCache[key] = matchStatusMatch
		return true
	}

	return cached.subtypeValue()
}

// isSubtypeDispatching is the custom dynamic dispatcher for various subtypes.
// Mirrors SubtypeChecker.isSubtypeDispatching in Java. Uses Go type switches
// instead of Java's getTypeClass() enum switch.
func (sc *SubtypeChecker) isSubtypeDispatching(subtype, supertype JSType) bool {
	switch sub := subtype.(type) {
	case *ArrowType:
		return sc.isArrowTypeSubtype(sub, supertype)
	case *EnumElementType:
		return sc.isEnumElementSubtype(sub, supertype)
	case *EnumType:
		return sc.isEnumSubtype(sub, supertype)
	case *noObjectType, *noType:
		return sc.isVariousBottomsSubtype(subtype, supertype)
	case *FunctionType:
		return sc.isFunctionSubtype(sub, supertype)
	case *TemplateType:
		return sc.isTemplateSubtype(sub, supertype)
	case *ProxyObjectType:
		return sc.isProxyObjectSubtype(sub, supertype)
	case *KnownSymbolType:
		return sc.isWellKnownSymbolSubtype(sub, supertype)
	default:
		return sc.isSubtypeHelper(subtype, supertype)
	}
}

func (sc *SubtypeChecker) isWellKnownSymbolSubtype(subtype *KnownSymbolType, supertype JSType) bool {
	if sup, ok := supertype.(*KnownSymbolType); ok {
		return subtype == sup
	}
	if _, ok := supertype.(*symbolType); ok {
		return true
	}
	return sc.isSubtypeHelper(subtype, supertype)
}

// isSubtypeHelper handles the common cases. Mirrors SubtypeChecker.isSubtypeHelper.
func (sc *SubtypeChecker) isSubtypeHelper(subtype, supertype JSType) bool {
	// Axiomatic cases.
	if subtype == supertype ||
		isUnknown(supertype) ||
		isAll(supertype) ||
		isUnknown(subtype) ||
		isNoType(subtype) {
		return true
	}

	// Special consideration for null and undefined in J2CL.
	if sc.subtypingMode == SubtypingModeIgnoreNullUndefined {
		if isNullType(subtype) || isVoidType(subtype) {
			return true
		}
	}

	// Reflexive case (value equality, not pointer equality).
	// In Go, value types already covered above via ==; for structs we rely on
	// pointer identity like Java does for most interned types.

	// Unwrap named types (proxy unwrapping for NamedType only).
	// TODO(SubtypeChecker): implement when NamedType is ported.
	// if nt, ok := subtype.(*NamedType); ok {
	//     return sc.isSubtypeDispatching(nt.getReferencedType(), supertype)
	// }
	// if nt, ok := supertype.(*NamedType); ok {
	//     return sc.isSubtypeDispatching(subtype, nt.getReferencedType())
	// }

	// Union decomposition.
	if sub, ok := subtype.(*UnionType); ok {
		// All alternates must be subtypes.
		for _, m := range sub.members {
			if !sc.isSubtypeCaching(m, supertype) {
				return false
			}
		}
		return true
	}
	if sup, ok := supertype.(*UnionType); ok {
		// Some alternate must be a supertype.
		for _, m := range sup.members {
			if sc.isSubtypeCaching(subtype, m) {
				return true
			}
		}
		return false
	}

	subtypeObj := toObjectType(subtype)
	supertypeObj := toObjectType(supertype)
	if subtypeObj == nil || supertypeObj == nil {
		return false // All cases for non object types covered by this point.
	}

	return sc.isObjectSubtypeHelper(subtypeObj, supertypeObj)
}

// isObjectSubtypeHelper handles object-to-object subtyping.
// Mirrors SubtypeChecker.isObjectSubtypeHelper.
func (sc *SubtypeChecker) isObjectSubtypeHelper(subtype, supertype objectTypeI) bool {
	// TODO(SubtypeChecker): TemplateTypeMap variance checks.
	// Java computes subtypeParams/supertypeParams from getTemplateTypeMap() and
	// checks bivariant/covariant/invariant constraints. Stubbed until
	// TemplateTypeMap full API is ported.

	// Check bivariant types (Object, IArrayLike, Array).
	if sc.isBivariantType(supertype) {
		// TODO(SubtypeChecker): bivariant template element check via getObjectElementKey.
		// For now fall through to structural/nominal checks.
		_ = true
	}

	// Structural subtyping when supertype is a structural type.
	if sc.useStructuralTyping() && isStructuralObjectType(supertype) {
		return sc.isStructuralSubtypeHelper(
			subtype, supertype, voidablePropsAreOptional)
	}

	// TODO(SubtypeChecker): covariant key checks (IThenable, ReadonlyArray, Iterator, etc.)
	// TODO(SubtypeChecker): isTypeMapSubmap

	// Interfaces: find all the interfaces implemented by this class and compare.
	// TODO(SubtypeChecker): interface subtype check via getCtorImplementedInterfaces / getCtorExtendedInterfaces.

	return isImplicitPrototypeOf(supertype, subtype)
}

// isStructuralSubtypeHelper checks structural subtyping.
// Mirrors SubtypeChecker.isStructuralSubtypeHelper.
func (sc *SubtypeChecker) isStructuralSubtypeHelper(
	subtype, supertype objectTypeI, optionality propertyOptionality,
) bool {
	for _, pname := range getPropertyNames(supertype) {
		supertypeProp := getPropertyType(supertype, pname)
		if supertypeProp == nil {
			continue
		}
		if hasProperty(subtype, pname) {
			subtypeProp := getPropertyType(subtype, pname)
			if !sc.isSubtypeCaching(subtypeProp, supertypeProp) {
				return false
			}
		} else if !optionality.isOptional(supertypeProp) {
			return false
		}
	}
	return true
}

func (sc *SubtypeChecker) isFunctionSubtype(subtype *FunctionType, nonFunctionSupertype JSType) bool {
	if sc.isSubtypeHelper(subtype, nonFunctionSupertype) {
		return true
	}

	if sup, ok := nonFunctionSupertype.(*FunctionType); ok {
		// TODO(SubtypeChecker): isInterface check. In Java: if supertype.isInterface() return true.
		// TODO(SubtypeChecker): if subtype.isInterface() return false.
		// TODO(SubtypeChecker): shouldTreatThisTypesAsCovariant + arrow subtype check.
		_ = sup
		return false
	}

	// TODO(SubtypeChecker): fall through to FUNCTION_PROTOTYPE check.
	// return sc.isSubtypeCaching(registry.getNativeType(JSTypeNative.FUNCTION_PROTOTYPE), nonFunctionSupertype)
	return false
}

func (sc *SubtypeChecker) isVariousBottomsSubtype(subtype, supertype JSType) bool {
	if sc.isSubtypeHelper(subtype, supertype) {
		return true
	}
	return isObjectType(supertype) && !isNoType(supertype)
}

func (sc *SubtypeChecker) isArrowTypeSubtype(subtype *ArrowType, nonArrowSupertype JSType) bool {
	supertype, ok := nonArrowSupertype.(*ArrowType)
	if !ok {
		return false
	}

	// Return type: covariant.
	if !sc.isSubtypeCaching(subtype.ReturnType, supertype.ReturnType) {
		return false
	}

	// Parameters: contravariant.
	// Mirrors the Java iterator loop in isArrowTypeSubtype.
	si := 0
	ti := 0
	for si < len(subtype.Params) && ti < len(supertype.Params) {
		subtypeParam := subtype.Params[si]
		supertypeParam := supertype.Params[ti]

		subtypeParamType := subtypeParam.JSType
		supertypeParamType := supertypeParam.JSType
		if subtypeParamType != nil {
			if supertypeParamType == nil || !sc.isSubtypeCaching(supertypeParamType, subtypeParamType) {
				return false
			}
		}

		thisIsVarArgs := subtypeParam.Variadic
		thatIsVarArgs := supertypeParam.Variadic
		thisIsOptional := thisIsVarArgs || subtypeParam.Optional
		thatIsOptional := thatIsVarArgs || supertypeParam.Optional

		if !thisIsOptional && thatIsOptional {
			isTopFunction := thatIsVarArgs &&
				(supertypeParamType == nil || isUnknown(supertypeParamType) || isNoType(supertypeParamType))
			if !isTopFunction {
				return false
			}
		}

		if !thisIsVarArgs {
			si++
		}
		if !thatIsVarArgs {
			ti++
		}

		if thisIsVarArgs && thatIsVarArgs {
			break
		}
	}

	if si < len(subtype.Params) {
		p := subtype.Params[si]
		return p.Optional || p.Variadic || ti < len(supertype.Params)
	}
	return true
}

func (sc *SubtypeChecker) isEnumSubtype(subtype *EnumType, supertype JSType) bool {
	// TODO(SubtypeChecker): check against OBJECT_TYPE native / OBJECT_PROTOTYPE native.
	return sc.isSubtypeHelper(subtype, supertype)
}

func (sc *SubtypeChecker) isEnumElementSubtype(subtype *EnumElementType, supertype JSType) bool {
	// TODO(b/136298690): stop creating such a 'meet' and remove this logic.
	if supElem, ok := supertype.(*EnumElementType); ok {
		if subtype.Enum == supElem.Enum {
			return sc.isSubtypeCaching(subtype.Enum.Underlying, supElem.Enum.Underlying)
		}
	}

	if sc.isSubtypeHelper(subtype, supertype) {
		return true
	}
	return sc.isSubtypeCaching(subtype.Enum.Underlying, supertype)
}

func (sc *SubtypeChecker) isProxyObjectSubtype(subtype *ProxyObjectType, supertype JSType) bool {
	// Don't check the cache here; proxy types are equal to their referenced types.
	return sc.isSubtypeDispatching(subtype.ReferencedType, supertype)
}

func (sc *SubtypeChecker) isTemplateSubtype(subtype *TemplateType, supertype JSType) bool {
	// TODO(SubtypeChecker): bounds-based ContainsUpperBoundSuperTypeVisitor check.
	// For now fall through to proxy subtype (TemplateType embeds ProxyObjectType-like behaviour).
	_ = subtype
	_ = supertype
	// Mirrors isProxyObjectSubtype path for unbound templates.
	// TemplateType in Java is a ProxyObjectType whose referenced type is its bound.
	// Since our TemplateType has no bound, treat as unknown (always true).
	return true
}

// isBivariantType mirrors SubtypeChecker.isBivariantType — returns true for
// Object, IArrayLike, and Array (native object types exempt from invariance).
func (sc *SubtypeChecker) isBivariantType(t objectTypeI) bool {
	name := objectTypeName(t)
	return bivariantTypes[name]
}

// GetTemplateKeyIfCovariantType mirrors SubtypeChecker.getTemplateKeyIfCovariantType.
// TODO(SubtypeChecker): full implementation requires TemplateTypeMap.hasTemplateKey and
// registry accessor methods (getIThenableTemplate, getReadonlyArrayElementKey, etc.).
// Returns nil until those are ported.
func GetTemplateKeyIfCovariantType(t JSType) *TemplateType {
	// TODO(SubtypeChecker): implement covariant key detection.
	return nil
}

// ---- helper interfaces and accessors for ObjectType variants ----

// objectTypeI is a minimal interface satisfied by *ObjectType and
// *TemplatizedType for use inside SubtypeChecker's structural checks.
// TODO(SubtypeChecker): expand when RecordType and other structural types are ported.
type objectTypeI interface {
	JSType
	getMembers() map[string]JSType
	getName() string
}

func (o *ObjectType) getMembers() map[string]JSType { return o.members }
func (o *ObjectType) getName() string               { return o.Name }
func (t *TemplatizedType) getMembers() map[string]JSType {
	if t.Base != nil {
		return t.Base.members
	}
	return nil
}
func (t *TemplatizedType) getName() string {
	if t.Base != nil {
		return t.Base.Name
	}
	return ""
}

func toObjectType(t JSType) objectTypeI {
	switch v := t.(type) {
	case *ObjectType:
		return v
	case *TemplatizedType:
		return v
	}
	return nil
}

func getPropertyNames(t objectTypeI) []string {
	switch v := t.(type) {
	case *ObjectType:
		names := make([]string, 0, len(v.members))
		for k := range v.members {
			names = append(names, k)
		}
		return names
	case *TemplatizedType:
		if v.Base != nil {
			names := make([]string, 0, len(v.Base.members))
			for k := range v.Base.members {
				names = append(names, k)
			}
			return names
		}
	}
	return nil
}

func getPropertyType(t objectTypeI, name string) JSType {
	switch v := t.(type) {
	case *ObjectType:
		return v.GetPropertyType(name)
	case *TemplatizedType:
		if v.Base != nil {
			return v.Base.GetPropertyType(name)
		}
	}
	return nil
}

func hasProperty(t objectTypeI, name string) bool {
	return getPropertyType(t, name) != nil
}

func isStructuralObjectType(t objectTypeI) bool {
	if obj, ok := t.(*ObjectType); ok {
		return obj.IsStructural
	}
	return false
}

func objectTypeName(t objectTypeI) string {
	if obj, ok := t.(*ObjectType); ok {
		return obj.Name
	}
	if tt, ok := t.(*TemplatizedType); ok && tt.Base != nil {
		return tt.Base.Name
	}
	return ""
}

func isImplicitPrototypeOf(supertype, subtype objectTypeI) bool {
	// Mirrors ObjectType.isImplicitPrototypeOf in Java — walks the @extends chain.
	// TODO(SubtypeChecker): full prototype chain walk (requires super field access).
	sub, ok := subtype.(*ObjectType)
	if !ok {
		return false
	}
	sup, ok := supertype.(*ObjectType)
	if !ok {
		return false
	}
	// Walk sub's prototype chain to find sup.
	for cur := sub; cur != nil; cur = cur.super {
		if cur == sup {
			return true
		}
	}
	return false
}

// isObjectType returns true when t is any ObjectType variant.
func isObjectType(t JSType) bool {
	switch t.(type) {
	case *ObjectType, *TemplatizedType, *FunctionType:
		return true
	}
	return false
}

// ---- helper predicates for primitive singletons ----

func isNullType(t JSType) bool { _, ok := t.(nullType); return ok }
func isVoidType(t JSType) bool { _, ok := t.(voidType); return ok }
func isNoType(t JSType) bool {
	switch t.(type) {
	case *noType, *noObjectType:
		return true
	}
	return false
}

// ---- stub types for not-yet-ported Java types ----

// ArrowType is the internal function signature type. Stub: full port deferred.
// Mirrors com/google/javascript/rhino/jstype/ArrowType.java.
// TODO(ArrowType): port ArrowType.java including getParameterList/getReturnType.
type ArrowType struct {
	// ReturnType is the covariant return.
	ReturnType JSType
	// Params is the ordered parameter list.
	Params []arrowParam
}

// arrowParam mirrors FunctionType.Parameter in Java.
type arrowParam struct {
	JSType   JSType
	Optional bool
	Variadic bool
}

func (a *ArrowType) TypeName() string   { return "arrow" }
func (a *ArrowType) IsSubtypeOf(o JSType) bool {
	return NewSubtypeChecker(nil).
		SetSubtype(a).
		SetSupertype(o).
		SetUsingStructuralSubtyping(false).
		SetSubtypingMode(SubtypingModeNormal).
		Check()
}

// ProxyObjectType is a transparent wrapper around another JSType. Stub.
// Mirrors com/google/javascript/rhino/jstype/ProxyObjectType.java.
// TODO(ProxyObjectType): port ProxyObjectType.java.
type ProxyObjectType struct {
	ReferencedType JSType
}

func (p *ProxyObjectType) TypeName() string        { return p.ReferencedType.TypeName() }
func (p *ProxyObjectType) IsSubtypeOf(o JSType) bool {
	return NewSubtypeChecker(nil).
		SetSubtype(p).
		SetSupertype(o).
		SetUsingStructuralSubtyping(false).
		SetSubtypingMode(SubtypingModeNormal).
		Check()
}

// KnownSymbolType is a well-known JavaScript symbol (e.g. Symbol.iterator). Stub.
// Mirrors com/google/javascript/rhino/jstype/KnownSymbolType.java.
// TODO(KnownSymbolType): port KnownSymbolType.java.
type KnownSymbolType struct {
	SymbolName string
}

func (k *KnownSymbolType) TypeName() string       { return "symbol(" + k.SymbolName + ")" }
func (k *KnownSymbolType) IsSubtypeOf(o JSType) bool {
	return NewSubtypeChecker(nil).
		SetSubtype(k).
		SetSupertype(o).
		SetUsingStructuralSubtyping(false).
		SetSubtypingMode(SubtypingModeNormal).
		Check()
}

// symbolType is the general Symbol primitive type (non-well-known). Stub.
// Mirrors com/google/javascript/rhino/jstype/SymbolType.java.
// TODO(SymbolType): port SymbolType.java.
type symbolType struct{}

func (s *symbolType) TypeName() string        { return "symbol" }
func (s *symbolType) IsSubtypeOf(o JSType) bool {
	if _, ok := o.(*symbolType); ok {
		return true
	}
	return isAll(o) || isUnknown(o)
}

// noObjectType is the bottom of the object hierarchy. Stub.
// Mirrors com/google/javascript/rhino/jstype/NoObjectType.java.
// TODO(NoObjectType): port NoObjectType.java.
type noObjectType struct{}

func (n *noObjectType) TypeName() string        { return "NoObject" }
func (n *noObjectType) IsSubtypeOf(_ JSType) bool { return true }

// noType is the absolute bottom type. Stub.
// Mirrors com/google/javascript/rhino/jstype/NoType.java.
// TODO(NoType): port NoType.java.
type noType struct{}

func (n *noType) TypeName() string        { return "NoType" }
func (n *noType) IsSubtypeOf(_ JSType) bool { return true }
