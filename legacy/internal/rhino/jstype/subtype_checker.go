// Port of com/google/javascript/rhino/jstype/SubtypeChecker.java @ closure-compiler v20260415
/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1/GPL 2.0
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
 *   Bob Jervis
 *   Google Inc.
 *
 * Alternatively, the contents of this file may be used under the terms of
 * the GNU General Public License Version 2 or later (the "GPL"), in which
 * case the provisions of the GPL are applicable instead of those above. If
 * you wish to allow use of your version of this file only under the terms of
 * the GPL and not to allow others to use your version of this file under the
 * MPL, indicate your decision by deleting the provisions above and replacing
 * them with the notice and other provisions required by the GPL. If you do
 * not delete the provisions above, a recipient may use your version of this
 * file under either the MPL or the GPL.
 *
 * ***** END LICENSE BLOCK ***** */

// Port of com/google/javascript/rhino/jstype/SubtypeChecker.java
// @ closure-compiler v20260415
//
// Represents the computation of a single supertype-subtype relationship.
//
// Each instance can only be used once. It acts as a scope for the computation
// and may accumulate state.

package jstype

// SubtypingMode controls how null/undefined are treated during subtype checks.
// Port of JSType.SubtypingMode in Java.
type SubtypingMode int

const (
	// SubtypingModeNormal is the standard subtyping mode.
	SubtypingModeNormal SubtypingMode = iota
	// SubtypingModeIgnoreNullUndefined treats null and undefined as subtypes
	// of everything (used by J2CL).
	SubtypingModeIgnoreNullUndefined
)

// subtypeMatchStatus mirrors JSType.MatchStatus in Java.
type subtypeMatchStatus int

const (
	subtypeProcessing subtypeMatchStatus = iota
	subtypeMatch
	subtypeNoMatch
)

func (s subtypeMatchStatus) subtypeValue() bool {
	return s == subtypeMatch
}

func subtypeMatchStatusOf(b bool) subtypeMatchStatus {
	if b {
		return subtypeMatch
	}
	return subtypeNoMatch
}

// subtypeCacheKey is a pair of JSType values used as a cache key.
// Port of SubtypeChecker.CacheKey in Java.
type subtypeCacheKey struct {
	left  JSType
	right JSType
}

// subtypeVariance mirrors SubtypeChecker.Variance in Java.
type subtypeVariance int

const (
	varianceCovariant     subtypeVariance = iota
	varianceContravariant subtypeVariance = iota
	varianceBivariant     subtypeVariance = iota
	varianceInvariant     subtypeVariance = iota
)

// subtypePropertyOptionality mirrors SubtypeChecker.PropertyOptionality in Java.
type subtypePropertyOptionality int

const (
	// voidablePropsAreOptional treats explicitly voidable properties as optional.
	voidablePropsAreOptional subtypePropertyOptionality = iota
	// allPropsAreRequired treats all properties as required even if explicitly voidable.
	allPropsAreRequired subtypePropertyOptionality = iota
)

func (o subtypePropertyOptionality) isOptional(propType JSType) bool {
	return o == voidablePropsAreOptional && containsVoid(propType)
}

// bivariantTypes are exempt from template type invariance.
// Port of SubtypeChecker.BIVARIANT_TYPES.
var bivariantTypeNames = map[string]bool{
	"Object":     true,
	"IArrayLike": true,
	"Array":      true,
}

// potentiallyCyclicRecursionDepth is an arbitrary depth at which to start
// checking for cyclic recursion. Port of POTENTIALLY_CYCLIC_RECURSION_DEPTH.
const potentiallyCyclicRecursionDepth = 20

// SubtypeChecker represents the computation of a single supertype-subtype
// relationship. Each instance can only be used once.
// Port of com/google/javascript/rhino/jstype/SubtypeChecker.java.
type SubtypeChecker struct {
	initialSupertype JSType
	initialSubtype   JSType

	isUsingStructuralTyping *bool
	subtypingMode           *SubtypingMode

	subtypeCache map[subtypeCacheKey]subtypeMatchStatus

	hasRun         bool
	recursionDepth int
}

// NewSubtypeChecker constructs a SubtypeChecker.
// Port of SubtypeChecker(JSTypeRegistry registry) constructor in Java.
func NewSubtypeChecker() *SubtypeChecker {
	return &SubtypeChecker{}
}

// SetSupertype sets the supertype for the check. Must be called before Check().
func (s *SubtypeChecker) SetSupertype(value JSType) *SubtypeChecker {
	s.checkHasNotRun()
	if s.initialSupertype != nil {
		panic("supertype already set")
	}
	if value == nil {
		panic("supertype must not be nil")
	}
	s.initialSupertype = value
	return s
}

// SetSubtype sets the subtype for the check. Must be called before Check().
func (s *SubtypeChecker) SetSubtype(value JSType) *SubtypeChecker {
	s.checkHasNotRun()
	if s.initialSubtype != nil {
		panic("subtype already set")
	}
	if value == nil {
		panic("subtype must not be nil")
	}
	s.initialSubtype = value
	return s
}

// SetUsingStructuralSubtyping configures whether structural subtyping is used.
func (s *SubtypeChecker) SetUsingStructuralSubtyping(value bool) *SubtypeChecker {
	s.checkHasNotRun()
	if s.isUsingStructuralTyping != nil {
		panic("isUsingStructuralTyping already set")
	}
	s.isUsingStructuralTyping = &value
	return s
}

// SetSubtypingMode configures the subtyping mode.
func (s *SubtypeChecker) SetSubtypingMode(value SubtypingMode) *SubtypeChecker {
	s.checkHasNotRun()
	if s.subtypingMode != nil {
		panic("subtypingMode already set")
	}
	s.subtypingMode = &value
	return s
}

func (s *SubtypeChecker) checkHasNotRun() {
	if s.hasRun {
		panic("SubtypeChecker has already run")
	}
}

// Check runs the subtype check. May only be called once.
// Port of SubtypeChecker.check() in Java.
func (s *SubtypeChecker) Check() bool {
	s.checkHasNotRun()
	s.hasRun = true
	return s.isSubtypeCaching(s.initialSubtype, s.initialSupertype)
}

// isSubtypeCaching is the top-level recursive entrypoint for subtyping logic.
// Caching is necessary to catch cyclic types.
// Port of SubtypeChecker.isSubtypeCaching() in Java.
func (s *SubtypeChecker) isSubtypeCaching(subtype, supertype JSType) bool {
	if subtype == nil || supertype == nil {
		return true
	}

	// Wait to instantiate/use the cache until we have some hint that there
	// may be recursion.
	if s.recursionDepth > potentiallyCyclicRecursionDepth {
		if s.subtypeCache == nil {
			s.subtypeCache = make(map[subtypeCacheKey]subtypeMatchStatus)
		}
	}

	// Once the cache exists, use it consistently.
	if s.subtypeCache == nil {
		s.recursionDepth++
		result := s.isSubtypeDispatching(subtype, supertype)
		s.recursionDepth--
		return result
	}

	key := subtypeCacheKey{left: subtype, right: supertype}
	if cached, ok := s.subtypeCache[key]; ok {
		if cached == subtypeProcessing {
			s.subtypeCache[key] = subtypeMatch
			return true
		}
		return cached.subtypeValue()
	}

	s.subtypeCache[key] = subtypeProcessing
	result := s.isSubtypeDispatching(subtype, supertype)
	s.subtypeCache[key] = subtypeMatchStatusOf(result)
	return result
}

// isSubtypeDispatching is a custom dynamic dispatcher for various subtypes.
// Port of SubtypeChecker.isSubtypeDispatching() in Java.
func (s *SubtypeChecker) isSubtypeDispatching(subtype, supertype JSType) bool {
	switch sub := subtype.(type) {
	case *EnumElementType:
		return s.isEnumElementSubtype(sub, supertype)
	case *EnumType:
		return s.isEnumSubtype(sub, supertype)
	case *FunctionType:
		return s.isFunctionSubtype(sub, supertype)
	case *TemplateType:
		return s.isTemplateSubtype(sub, supertype)
	default:
		return s.isSubtypeHelper(subtype, supertype)
	}
}

// isSubtypeHelper handles the common subtype logic.
// Port of SubtypeChecker.isSubtypeHelper() in Java.
func (s *SubtypeChecker) isSubtypeHelper(subtype, supertype JSType) bool {
	// Axiomatic cases.
	if subtype == supertype ||
		isAll(supertype) ||
		isUnknown(supertype) ||
		isUnknown(subtype) {
		return true
	}

	// NoType (bottom) is a subtype of everything.
	if isNoType(subtype) {
		return true
	}

	// Special consideration for null and undefined in J2CL.
	if s.subtypingMode != nil && *s.subtypingMode == SubtypingModeIgnoreNullUndefined {
		if _, ok := subtype.(nullType); ok {
			return true
		}
		if _, ok := subtype.(voidType); ok {
			return true
		}
	}

	// Union decomposition.
	if subUnion, ok := subtype.(*UnionType); ok {
		// All alternates must be subtypes.
		for _, alt := range subUnion.members {
			if !s.isSubtypeCaching(alt, supertype) {
				return false
			}
		}
		return true
	}
	if supUnion, ok := supertype.(*UnionType); ok {
		// Some alternate must be a supertype.
		for _, alt := range supUnion.members {
			if s.isSubtypeCaching(subtype, alt) {
				return true
			}
		}
		return false
	}

	subObj, subIsObj := subtype.(*ObjectType)
	supObj, supIsObj := supertype.(*ObjectType)
	if !subIsObj || !supIsObj {
		return false // All cases for non object types have been covered.
	}

	return s.isObjectSubtypeHelper(subObj, supObj)
}

// isObjectSubtypeHelper handles object-to-object subtyping.
// Port of SubtypeChecker.isObjectSubtypeHelper() in Java.
func (s *SubtypeChecker) isObjectSubtypeHelper(subtype, supertype *ObjectType) bool {
	// --- Template type checks ---
	// Determine whether both sides are templatized to the same base type.
	subTT := subtype.templatizedSelf
	supTT := supertype.templatizedSelf
	subAndSuperAreSameBaseType := false
	if subTT != nil && supTT != nil && subTT.Base == supTT.Base {
		subAndSuperAreSameBaseType = true
	}
	bivariantMatch := false

	// Array, Object, and IArrayLike are exempt from template invariance
	// (bivariant element types). Port of isBivariantType + meetVarianceConstraint.
	if objectTypeIfNativeOf(supertype) != nil && bivariantTypeNames[supertype.GetReferenceName()] {
		if subTT != nil && supTT != nil {
			// Check bivariance for matching template args by position.
			for i := 0; i < len(supTT.Args) && i < len(subTT.Args); i++ {
				if !s.meetVarianceConstraint(varianceBivariant, subTT.Args[i], supTT.Args[i]) {
					return false
				}
			}
		}
		bivariantMatch = true
	}

	// Find a single covariant template key for types like ReadonlyArray,
	// Iterator, Generator, etc. In Go we use positional arg index 0.
	covariantArgIdx := -1
	if !bivariantMatch {
		covariantArgIdx = covariantArgIndexOf(supertype)
	}

	if !subAndSuperAreSameBaseType &&
		s.isUsingStructuralTyping != nil && *s.isUsingStructuralTyping &&
		supertype.IsStructuralType() {
		// Structural subtyping branch — must be done before templatization.
		// Special case: Array subtype of ReadonlyArray is checked covariantly.
		if covariantArgIdx == 0 && isNativeArrayObjectType(subtype) {
			if subTT != nil && supTT != nil && covariantArgIdx < len(supTT.Args) {
				subVal := JSType(nil)
				if covariantArgIdx < len(subTT.Args) {
					subVal = subTT.Args[covariantArgIdx]
				}
				superVal := supTT.Args[covariantArgIdx]
				return s.meetVarianceConstraint(varianceCovariant, subVal, superVal)
			}
		}
		return s.isStructuralSubtypeHelper(subtype, supertype, voidablePropsAreOptional)
	}

	// Template type submap check (invariant or covariant).
	if !bivariantMatch {
		if covariantArgIdx >= 0 {
			// Covariant single-arg check.
			if subTT != nil && supTT != nil {
				if covariantArgIdx < len(supTT.Args) {
					subVal := JSType(nil)
					if covariantArgIdx < len(subTT.Args) {
						subVal = subTT.Args[covariantArgIdx]
					}
					superVal := supTT.Args[covariantArgIdx]
					if !s.meetVarianceConstraint(varianceCovariant, subVal, superVal) {
						return false
					}
				}
			}
		} else if !s.isTypeMapSubmap(subtype, supertype) {
			return false
		}
	}

	// --- Interface chain walking ---
	// Find all the interfaces implemented/extended by subtype and compare
	// each one to the supertype. Port of the "Interfaces" block in Java.
	subtypeCtor := subtype.GetConstructor()
	supertypeCtor := supertype.GetConstructor()
	if subtypeCtor != nil && subtypeCtor.IsInterface() {
		for _, subtypeInterface := range subtype.GetCtorExtendedInterfaces() {
			if s.isSubtypeCaching(subtypeInterface, supertype) {
				return true
			}
		}
	} else if supertypeCtor != nil && supertypeCtor.IsInterface() {
		for _, subtypeInterface := range subtype.GetCtorImplementedInterfaces() {
			if s.isSubtypeCaching(subtypeInterface, supertype) {
				return true
			}
		}
	}

	// --- Nominal prototype-chain check ---
	return supertype.IsImplicitPrototypeOf(subtype)
}

// objectTypeIfNativeOf returns obj if it is a native ObjectType, else nil.
// Port of SubtypeChecker.getObjectTypeIfNative() in Java.
func objectTypeIfNativeOf(obj *ObjectType) *ObjectType {
	if obj != nil && obj.IsNativeObjectType() {
		return obj
	}
	return nil
}

// isNativeArrayObjectType reports whether obj is the native Array type.
// Port of SubtypeChecker.isNativeArrayType() in Java.
func isNativeArrayObjectType(obj *ObjectType) bool {
	return obj.IsNativeObjectType() && obj.GetReferenceName() == "Array"
}

// covariantArgIndexOf returns 0 when supertype uses single-arg covariant
// template semantics (ReadonlyArray, Iterator, Generator, AsyncIterator,
// Iterable, IteratorIterable, IIterableResult, AsyncIterable), or -1
// otherwise. In Go the template args are positional so index 0 is used.
// Port of SubtypeChecker.getTemplateKeyIfCovariantType() in Java.
func covariantArgIndexOf(supertype *ObjectType) int {
	native := objectTypeIfNativeOf(supertype)
	if native == nil {
		// IThenable covariance applies to non-native subtypes too in Java,
		// but requires registry.getIThenableTemplate() which is not yet
		// available in the Go port.  Fall through to invariant check.
		return -1
	}
	switch native.GetReferenceName() {
	case "ReadonlyArray", "Iterator", "Generator", "AsyncIterator",
		"Iterable", "IteratorIterable", "IIterableResult", "AsyncIterable":
		return 0
	}
	return -1
}

// isTypeMapSubmap checks that every template argument in supertype is matched
// by an invariant-compatible argument in subtype.
// Port of SubtypeChecker.isTypeMapSubmap() in Java.
func (s *SubtypeChecker) isTypeMapSubmap(subtype, supertype *ObjectType) bool {
	// Prototype types never have template type maps — assume true.
	// Port of the isFunctionPrototypeType() early-return in Java.
	// In Go, prototype objects have a non-nil ownerFunction (GetOwnerFunction).
	if subtype.GetOwnerFunction() != nil {
		return true
	}

	supTT := supertype.templatizedSelf
	if supTT == nil {
		// Supertype has no template args; nothing to check.
		return true
	}
	subTT := subtype.templatizedSelf

	for i, superVal := range supTT.Args {
		subVal := JSType(nil)
		if subTT != nil && i < len(subTT.Args) {
			subVal = subTT.Args[i]
		} else {
			// Key not present in subtype map — loosenTypechecking is not
			// modelled in Go yet; conservatively fail.
			return false
		}
		if !s.meetVarianceConstraint(varianceInvariant, subVal, superVal) {
			return false
		}
	}
	return true
}

// meetVarianceConstraint checks whether override is variance-compatible
// with reference.  Port of SubtypeChecker.meetVarianceConstraint() in Java.
func (s *SubtypeChecker) meetVarianceConstraint(variance subtypeVariance, override, reference JSType) bool {
	switch variance {
	case varianceCovariant:
		return s.isSubtypeCaching(override, reference)
	case varianceContravariant:
		return s.isSubtypeCaching(reference, override)
	case varianceBivariant:
		return s.meetVarianceConstraint(varianceCovariant, reference, override) ||
			s.meetVarianceConstraint(varianceContravariant, reference, override)
	case varianceInvariant:
		return s.meetVarianceConstraint(varianceCovariant, reference, override) &&
			s.meetVarianceConstraint(varianceContravariant, reference, override)
	default:
		return false
	}
}

// isStructuralSubtypeHelper handles structural subtype checking for @record types.
// Port of SubtypeChecker.isStructuralSubtypeHelper() in Java.
//
// subtype is a subtype of record type supertype iff:
//  1. subtype has all the non-optional properties declared in supertype.
//  2. For each property of supertype, its type must be a supertype of the
//     corresponding property of subtype.
func (s *SubtypeChecker) isStructuralSubtypeHelper(subtype, supertype *ObjectType, optionality subtypePropertyOptionality) bool {
	for prop, supertypeProp := range supertype.members {
		if !s.checkPropertyPresence(prop, subtype, supertypeProp, optionality) {
			return false
		}
	}
	return true
}

// checkPropertyPresence checks whether a required property is present and
// has a compatible type.
// Port of SubtypeChecker.checkPropertyPresence() in Java.
func (s *SubtypeChecker) checkPropertyPresence(property string, subtype *ObjectType, supertypeProp JSType, optionality subtypePropertyOptionality) bool {
	subtypeProp := subtype.GetPropertyType(property)
	if subtypeProp != nil {
		if !s.isSubtypeCaching(subtypeProp, supertypeProp) {
			return false
		}
	} else if !optionality.isOptional(supertypeProp) {
		return false
	}
	return true
}

// isFunctionSubtype handles function-to-anything subtyping.
// Port of SubtypeChecker.isFunctionSubtype() in Java.
func (s *SubtypeChecker) isFunctionSubtype(subtype *FunctionType, nonFunctionSupertype JSType) bool {
	if s.isSubtypeHelper(subtype, nonFunctionSupertype) {
		return true
	}

	if superFn, ok := nonFunctionSupertype.(*FunctionType); ok {
		// Covariant return type.
		if !s.isSubtypeCaching(subtype.Return, superFn.Return) {
			return false
		}

		// Contravariant parameter types.
		// If super has fewer params, that's OK (g ignores extra args).
		// If sub requires more params than super provides, that's a mismatch.
		subLen := len(subtype.Params)
		supLen := len(superFn.Params)
		minLen := subLen
		if supLen < minLen {
			minLen = supLen
		}
		for i := 0; i < minLen; i++ {
			// contravariant: supertype param must be subtype of subtype param
			if !s.isSubtypeCaching(superFn.Params[i], subtype.Params[i]) {
				return false
			}
		}
		return true
	}

	// Any FunctionType is a subtype of the built-in Function instance.
	if obj, ok := nonFunctionSupertype.(*ObjectType); ok && obj.Name == "Function" {
		return true
	}

	return false
}

// isEnumSubtype handles enum-to-anything subtyping.
// Port of SubtypeChecker.isEnumSubtype() in Java.
func (s *SubtypeChecker) isEnumSubtype(subtype *EnumType, supertype JSType) bool {
	if obj, ok := supertype.(*ObjectType); ok && obj.Name == "Object" {
		return true
	}
	return s.isSubtypeHelper(subtype, supertype)
}

// isEnumElementSubtype handles enum element-to-anything subtyping.
// Port of SubtypeChecker.isEnumElementSubtype() in Java.
func (s *SubtypeChecker) isEnumElementSubtype(subtype *EnumElementType, supertype JSType) bool {
	// Hack: if the supertype is an enum element with the same enum pointer,
	// compare the primitive types instead.
	// TODO(b/136298690): stop creating such a 'meet' and remove this logic.
	if supElem, ok := supertype.(*EnumElementType); ok && subtype.Enum == supElem.Enum {
		return s.isSubtypeCaching(subtype.Enum.Underlying, supElem.Enum.Underlying)
	}

	if s.isSubtypeHelper(subtype, supertype) {
		return true
	}
	return s.isSubtypeCaching(subtype.Enum.Underlying, supertype)
}

// isTemplateSubtype handles template type subtyping.
// Port of SubtypeChecker.isTemplateSubtype() in Java.
// For our simplified type system, template types are treated as unknown.
func (s *SubtypeChecker) isTemplateSubtype(subtype *TemplateType, supertype JSType) bool {
	// In our simplified type system, template types act like Unknown.
	return true
}

// isNoType reports whether t is the bottom type (NoType).
func isNoType(t JSType) bool {
	_, ok := t.(*NoType)
	return ok
}
