// Port of com/google/javascript/rhino/jstype/EqualityChecker.java @ closure-compiler v20261021

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
 *   Bob Jervis
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino

import (
	"unsafe"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// EqualityChecker represents the computation of a single equality
// relationship.
//
// Each instance can only be used once. It acts as a scope for the
// computation and may accumulate state.
type EqualityChecker struct {
	// Go: nil is Java null.
	eqMethod *EqualityChecker_EqMethod

	eqCache        *jutil.ObjectLinkedHashMap[*EqualityChecker_CacheKey, JSType_MatchStatus]
	recursionDepth int
	hasRun         bool
}

// EqualityChecker_EqMethod represents different ways for comparing equality
// among types.
type EqualityChecker_EqMethod int32

const (
	// EqualityChecker_EqMethod_IDENTITY indicates that the two types should
	// behave exactly the same under all type operations.
	//
	// Thus, {string} != {?} and {Unresolved} != {?}
	EqualityChecker_EqMethod_IDENTITY EqualityChecker_EqMethod = iota

	// EqualityChecker_EqMethod_DATA_FLOW indicates that the two types are
	// almost exactly the same, and that a data flow analysis algorithm
	// comparing them should consider them equal.
	//
	// In traditional type inference, the types form a finite lattice, and
	// this ensures that type inference will terminate.
	//
	// In our type system, the unknown types do not obey the lattice rules. So
	// if we continue to perform inference over the unknown types, we may never
	// terminate.
	//
	// By treating all unknown types as equivalent for the purposes of data
	// flow analysis, we ensure that the algorithm will terminate.
	//
	// Thus, {string} != {?} and {Unresolved} ~= {?}
	EqualityChecker_EqMethod_DATA_FLOW
)

var equalityChecker_EqMethod_names = [...]string{"IDENTITY", "DATA_FLOW"}

// String is Enum.name().
func (m EqualityChecker_EqMethod) String() string { return equalityChecker_EqMethod_names[m] }

// EqualityChecker_POTENTIALLY_CYCLIC_RECURSION_DEPTH is an arbitrary depth
// at which to start checking for cyclic recursion.
//
// Checking for cycles is expensive. Most types are not cyclic; they have
// have fairly shallow composition. Therefore, we want to avoid checking
// cycles them until we have a strong expectation that it's necessary.
//
// This number is an optimization, not a correctness requirement. The value
// could be anything (short of stack overflow) and we should get the same
// result.
//
// Java: EqualityChecker.POTENTIALLY_CYCLIC_RECURSION_DEPTH
const EqualityChecker_POTENTIALLY_CYCLIC_RECURSION_DEPTH = 20

//go:noinline
func failJSTypeState() {
	panic(jutil.NewIllegalStateException())
}

//go:noinline
func failJSTypeNull() {
	panic(jutil.NewNullPointerException())
}

// Java: EqualityChecker.setEqMethod(EqMethod)
func (c *EqualityChecker) SetEqMethod(x EqualityChecker_EqMethod) *EqualityChecker {
	c.CheckHasNotRun()
	if c.eqMethod != nil {
		failJSTypeState()
	}
	// Go-only: performance — the nullable enum points into a shared table
	// rather than at a fresh copy of x (it is only ever read through). The
	// copy is a variable of the else branch: taking &x would move x to the
	// heap on every call.
	if x >= 0 && int(x) < len(equalityChecker_eqMethods) {
		c.eqMethod = &equalityChecker_eqMethods[x]
	} else {
		other := x
		c.eqMethod = &other
	}
	return c
}

// equalityChecker_eqMethods holds each EqMethod, for SetEqMethod.
//
// Go-only: performance — see SetEqMethod.
var equalityChecker_eqMethods = [...]EqualityChecker_EqMethod{
	EqualityChecker_EqMethod_IDENTITY, EqualityChecker_EqMethod_DATA_FLOW}

// Java: EqualityChecker.checkHasNotRun()
func (c *EqualityChecker) CheckHasNotRun() {
	if c.hasRun {
		failJSTypeState()
	}
}

// Java: EqualityChecker()
func NewEqualityChecker() *EqualityChecker {
	return &EqualityChecker{}
}

// Java: EqualityChecker.check(JSType, JSType)
func (c *EqualityChecker) Check(left JSType, right JSType) bool {
	c.CheckHasNotRun()
	c.hasRun = true
	return c.AreEqualCaching(left, right)
}

func equalityChecker_AreIdentityEqual(left JSType, right JSType) bool {
	lb := jsTypeBase(left)
	rb := jsTypeBase(right)
	if lb == rb {
		return true
	}
	if lb == nil || rb == nil {
		return false
	}
	return equalityChecker_AreIdentityEqualWithBases(left, lb, right, rb)
}

func equalityChecker_AreIdentityEqualWithBases(left JSType, lb *JSTypeBase, right JSType, rb *JSTypeBase) bool {
	lk := lb.objFlags & objKindMask
	rk := rb.objFlags & objKindMask
	if (((uint32(1)<<lk)|(uint32(1)<<rk))&jsKindComplexEqMask) == 0 &&
		!lb.unknown && !rb.unknown &&
		lb.resolveResult != nil && rb.resolveResult != nil {
		if lk != rk {
			return false
		}
		switch lk {
		case objKindFunction:
			lf := (*FunctionTypeBase)(unsafe.Pointer(lb))
			rf := (*FunctionTypeBase)(unsafe.Pointer(rb))
			if lf.call.fnKind != rf.call.fnKind {
				return false
			}
			if lf.call.fnKind == FunctionType_Kind_CONSTRUCTOR || lf.call.fnKind == FunctionType_Kind_INTERFACE {
				return false
			}
			c := EqualityChecker{
				eqMethod: &equalityChecker_eqMethods[EqualityChecker_EqMethod_IDENTITY],
				hasRun:   true,
			}
			return c.AreFunctionEqual(lf, rf)
		case objKindRecord:
			c := EqualityChecker{
				eqMethod: &equalityChecker_eqMethods[EqualityChecker_EqMethod_IDENTITY],
				hasRun:   true,
			}
			if !c.AreTypeMapEqual(lb.templateTypeMap, rb.templateTypeMap) {
				return false
			}
			return c.AreRecordEqual((*RecordType)(unsafe.Pointer(lb)), (*RecordType)(unsafe.Pointer(rb)))
		case objKindNoResolved:
			return (*NoResolvedType)(unsafe.Pointer(lb)).referenceName ==
				(*NoResolvedType)(unsafe.Pointer(rb)).referenceName
		default:
			return false
		}
	}
	c := EqualityChecker{
		eqMethod: &equalityChecker_eqMethods[EqualityChecker_EqMethod_IDENTITY],
		hasRun:   true,
	}
	return c.AreEqualInternal(left, right)
}

// CheckParameters returns whether the parameters (ignoring any other
// aspects) of the two types are equal.
//
// Java: EqualityChecker.checkParameters(ArrowType, ArrowType)
func (c *EqualityChecker) CheckParameters(left *ArrowType, right *ArrowType) bool {
	c.CheckHasNotRun()
	c.hasRun = true
	return left == right || c.AreArrowParameterEqual(left, right)
}

func (c *EqualityChecker) CheckArrows(left *ArrowType, right *ArrowType) bool {
	c.CheckHasNotRun()
	c.hasRun = true
	return left == right || c.AreArrowEqual(left, right)
}

// Java: EqualityChecker.areEqualCaching(JSType, JSType)
func (c *EqualityChecker) AreEqualCaching(left JSType, right JSType) bool {
	// Wait to instantiate/use the cache until we have some hint that there may
	// be recursion.
	if c.recursionDepth > EqualityChecker_POTENTIALLY_CYCLIC_RECURSION_DEPTH {
		if c.eqCache == nil {
			c.eqCache = jutil.NewObjectLinkedHashMap[*EqualityChecker_CacheKey, JSType_MatchStatus]()
		}
	}

	// Once the cache exists, use it consistently.
	if c.eqCache == nil {
		c.recursionDepth++
		res := c.AreEqualInternal(left, right)
		c.recursionDepth--
		return res
	}

	key := NewEqualityChecker_CacheKey(left, right)
	cached, present := c.eqCache.Lookup(key)

	if !present {
		c.eqCache.Put(key, JSType_MatchStatus_PROCESSING)
		result := c.AreEqualInternal(left, right)
		c.eqCache.Put(key, JSType_MatchStatus_ValueOf(result))
		return result
	}

	if cached == JSType_MatchStatus_PROCESSING {
		c.eqCache.Put(key, JSType_MatchStatus_MATCH)
		return true
	}

	return cached.SubtypeValue()
}

// Java: EqualityChecker.areEqualInternal(JSType, JSType)
func (c *EqualityChecker) AreEqualInternal(left JSType, right JSType) bool {
	lb := jsTypeBase(left)
	rb := jsTypeBase(right)
	if lb == rb {
		return true
	} else if lb == nil || rb == nil {
		return false
	}

	lk := lb.objFlags & objKindMask
	rk := rb.objFlags & objKindMask
	if c.eqCache == nil &&
		(((uint32(1)<<lk)|(uint32(1)<<rk))&jsKindComplexEqMask) == 0 &&
		!lb.unknown && !rb.unknown &&
		lb.resolveResult != nil && rb.resolveResult != nil {
		if lk != rk {
			return false
		}
		switch lk {
		case objKindFunction:
			return c.AreFunctionEqual((*FunctionTypeBase)(unsafe.Pointer(lb)), (*FunctionTypeBase)(unsafe.Pointer(rb)))
		case objKindRecord:
			if !c.AreTypeMapEqual(lb.templateTypeMap, rb.templateTypeMap) {
				return false
			}
			return c.AreRecordEqual((*RecordType)(unsafe.Pointer(lb)), (*RecordType)(unsafe.Pointer(rb)))
		case objKindNoResolved:
			return (*NoResolvedType)(unsafe.Pointer(lb)).referenceName ==
				(*NoResolvedType)(unsafe.Pointer(rb)).referenceName
		default:
			return false
		}
	}

	if left.IsNoResolvedType() && right.IsNoResolvedType() {
		return jsTypeNullableStringEquals(
			left.ToObjectType().GetReferenceName(),
			right.ToObjectType().GetReferenceName())
	}

	leftUnknown := left.IsUnknownType()
	rightUnknown := right.IsUnknownType()
	if leftUnknown || rightUnknown {
		if *c.eqMethod == EqualityChecker_EqMethod_DATA_FLOW {
			// If we're checkings data flow, then two types are the same if they're
			// both unknown.
			return leftUnknown && rightUnknown
		} else if leftUnknown && rightUnknown && (left.IsNominalType() != right.IsNominalType()) {
			// If they're both unknown, but one is a nominal type and the other is
			// not, then we should fail out immediately. left ensures right we won't
			// unbox the unknowns further down.
			return false
		}
	}

	if left.IsUnionType() {
		if right.IsUnionType() {
			return c.AreUnionEqual(left.ToMaybeUnionType(), right.ToMaybeUnionType())
		}
		leftAlts := left.ToMaybeUnionType().GetAlternates()
		if len(leftAlts) == 1 {
			return c.AreEqualInternal(leftAlts[0], right)
		}
		if proxyObjectType, ok := right.(ProxyObjectType); ok {
			if _, isTemplate := right.(*TemplateType); !isTemplate {
				return c.AreTypeMapEqual(left.GetTemplateTypeMap(), right.GetTemplateTypeMap()) &&
					c.AreEqualCaching(left, proxyObjectType.GetReferencedTypeInternal())
			}
		}
		return false
	} else if right.IsUnionType() {
		rightAlts := right.ToMaybeUnionType().GetAlternates()
		if len(rightAlts) == 1 {
			return c.AreEqualInternal(left, rightAlts[0])
		}
		if proxyObjectType, ok := left.(ProxyObjectType); ok {
			if _, isTemplate := left.(*TemplateType); !isTemplate {
				return c.AreTypeMapEqual(left.GetTemplateTypeMap(), right.GetTemplateTypeMap()) &&
					c.AreEqualCaching(proxyObjectType.GetReferencedTypeInternal(), right)
			}
		}
		return false
	}

	if left.IsFunctionType() && right.IsFunctionType() {
		return c.AreFunctionEqual(left.ToMaybeFunctionType(), right.ToMaybeFunctionType())
	}

	if left.IsRecordType() && right.IsRecordType() {
		if !c.AreTypeMapEqual(left.GetTemplateTypeMap(), right.GetTemplateTypeMap()) {
			return false
		}
		return c.AreRecordEqual(left.ToMaybeRecordType(), right.ToMaybeRecordType())
	}

	if left.IsNominalType() && right.IsNominalType() {
		if !c.AreTypeMapEqual(left.GetTemplateTypeMap(), right.GetTemplateTypeMap()) {
			return false
		}
		leftUnwrapped := EqualityChecker_UnwrapNominalTypeProxies(left.ToObjectType())
		rightUnwrapped := EqualityChecker_UnwrapNominalTypeProxies(right.ToObjectType())

		jutil.CheckState(leftUnwrapped.IsNominalType() && rightUnwrapped.IsNominalType())
		if left.IsResolved() && right.IsResolved() {
			return leftUnwrapped == rightUnwrapped
		} else {
			// TODO(b/140763807): this is not valid across scopes pre-resolution.
			nameOfleft := jutil.CheckNotNull(leftUnwrapped.GetReferenceName())
			nameOfright := jutil.CheckNotNull(rightUnwrapped.GetReferenceName())
			googModuleIdOfLeft := EqualityChecker_GetGoogModuleId(leftUnwrapped)
			googModuleIdOfRight := EqualityChecker_GetGoogModuleId(rightUnwrapped)
			return jsTypeNullableStringEquals(nameOfleft, nameOfright) &&
				jsTypeNullableStringEquals(googModuleIdOfLeft, googModuleIdOfRight)
		}
	}

	// Unwrap proxies.
	//
	// Remember that `TemplateType` has identity semantics and shouldn't be
	// unwrapped.
	if proxyObjectType, ok := left.(ProxyObjectType); ok {
		if _, isTemplate := left.(*TemplateType); !isTemplate {
			if !c.AreTypeMapEqual(left.GetTemplateTypeMap(), right.GetTemplateTypeMap()) {
				return false
			}
			return c.AreEqualCaching(proxyObjectType.GetReferencedTypeInternal(), right)
		}
	}
	if proxyObjectType, ok := right.(ProxyObjectType); ok {
		if _, isTemplate := right.(*TemplateType); !isTemplate {
			if !c.AreTypeMapEqual(left.GetTemplateTypeMap(), right.GetTemplateTypeMap()) {
				return false
			}
			return c.AreEqualCaching(left, proxyObjectType.GetReferencedTypeInternal())
		}
	}

	// Relies on the fact right for the base JSType, only one instance of each
	// sub-type will ever be created in a given registry, so there is no need
	// to verify members. If the object pointers are not identical, then the
	// type member must be different.
	return false
}

// jsTypeNullableStringEquals is Objects.equals on two nullable Strings.
func jsTypeNullableStringEquals(a, b *string) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return *a == *b
}

// AreUnionEqual: two union types are equal if, after flattening nested
// union types, they have the same number of alternates and all alternates
// are equal.
//
// Java: EqualityChecker.areUnionEqual(UnionType, UnionType)
func (c *EqualityChecker) AreUnionEqual(left *UnionType, right *UnionType) bool {
	leftAlternates := left.GetAlternates()
	rightAlternates := right.GetAlternates()
	if *c.eqMethod == EqualityChecker_EqMethod_IDENTITY && len(leftAlternates) != len(rightAlternates) {
		return false
	}

outer:
	for i := 0; i < len(rightAlternates); i++ {
		rightAlt := rightAlternates[i]
		rightBase := jsTypeBase(rightAlt)
		if i < len(leftAlternates) && jsTypeBase(leftAlternates[i]) == rightBase {
			continue outer
		}
		for k := 0; k < len(leftAlternates); k++ {
			if jsTypeBase(leftAlternates[k]) == rightBase {
				continue outer
			}
		}

		for k := 0; k < len(leftAlternates); k++ {
			leftAlt := leftAlternates[k]
			if c.AreEqualCaching(leftAlt, rightAlt) {
				continue outer
			}
		}

		return false
	}

	return true
}

// Java: EqualityChecker.getGoogModuleId(ObjectType)
func EqualityChecker_GetGoogModuleId(type_ ObjectType) *string {
	if functionType, ok := type_.(FunctionType); ok {
		return functionType.GetGoogModuleId()
	}
	if type_.GetConstructor() != nil {
		return type_.GetConstructor().GetGoogModuleId()
	}
	return nil
}

// AreFunctionEqual: two function types are equal if their signatures match.
// Since they don't have signatures, two interfaces are equal if their names
// match.
//
// Java: EqualityChecker.areFunctionEqual(FunctionType, FunctionType)
func (c *EqualityChecker) AreFunctionEqual(left FunctionType, right FunctionType) bool {
	if left == right {
		// Identity needs to be re-checked in case a proxy was unwrapped when
		// entering this case.
		return true
	}

	if left.GetKind() != right.GetKind() {
		return false
	}

	switch left.GetKind() {
	case FunctionType_Kind_CONSTRUCTOR, FunctionType_Kind_INTERFACE:
		// constructors and interfaces use identity semantics, which we checked
		// for above.
		return false
	case FunctionType_Kind_ORDINARY:
		return c.AreEqualCaching(left.GetTypeOfThis(), right.GetTypeOfThis()) &&
			c.AreArrowEqual(left.GetInternalArrowType(), right.GetInternalArrowType()) &&
			jsTypeClosurePrimitiveEquals(left.GetClosurePrimitive(), right.GetClosurePrimitive())
	}
	panic(jutil.NewAssertionErrorNoMessage())
}

// jsTypeClosurePrimitiveEquals is Objects.equals on two nullable
// ClosurePrimitives.
func jsTypeClosurePrimitiveEquals(a, b *ClosurePrimitive) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return *a == *b
}

// Java: EqualityChecker.areArrowEqual(ArrowType, ArrowType)
func (c *EqualityChecker) AreArrowEqual(left *ArrowType, right *ArrowType) bool {
	return c.AreEqualCaching(left.GetReturnType(), right.GetReturnType()) &&
		c.AreArrowParameterEqual(left, right)
}

// Java: EqualityChecker.areArrowParameterEqual(ArrowType, ArrowType)
func (c *EqualityChecker) AreArrowParameterEqual(left *ArrowType, right *ArrowType) bool {
	if len(left.GetParameterList()) != len(right.GetParameterList()) {
		return false
	}

	for i := 0; i < len(left.GetParameterList()); i++ {
		leftParam := left.GetParameterList()[i]
		rightParam := right.GetParameterList()[i]

		leftParamType := leftParam.GetJSType()
		rightParamType := rightParam.GetJSType()
		if leftParamType != nil {
			// Both parameter lists give a type for this param, it should be equal
			if rightParamType != nil && !c.AreEqualCaching(leftParamType, rightParamType) {
				return false
			}
		} else {
			if rightParamType != nil {
				return false
			}
		}

		// Check var_args/optionality
		if leftParam.IsOptional() != rightParam.IsOptional() {
			return false
		}

		if leftParam.IsVariadic() != rightParam.IsVariadic() {
			return false
		}
	}

	return true
}

// AreRecordEqual checks for equality on inline record types (e.g. `{a:
// number}`).
//
// Only inline record types can be compared for equality based on structure.
// All other structural types, such as @records and object-literal types,
// have some concept of uniquenes. Such structural relationships are
// correctly expressed in terms of subtyping.
//
// Java: EqualityChecker.areRecordEqual(RecordType, RecordType)
func (c *EqualityChecker) AreRecordEqual(left *RecordType, right *RecordType) bool {
	// Don't check inherited properties; checking them is both incorrect and
	// slow.
	//
	// The full definition of a record type is contained in its "own"
	// properties (i.e. `{a: boolean, toString: function(...)}` and `{a:
	// boolean}` are not interchangable). This is in part because all inline
	// record types share the same inheritance.
	//
	// Additionally, code that makes heavy use of inline record types compiles
	// very slowly if the set of inherited properties is recomputed during
	// every equality check.
	leftEntries := left.properties.entries()
	rightEntries := right.properties.entries()
	if len(leftEntries) != len(rightEntries) {
		return false
	}
	left.properties.ensureSorted()
	right.properties.ensureSorted()
	leftEntries = left.properties.entries()
	rightEntries = right.properties.entries()
	for i := range leftEntries {
		if leftEntries[i].nameID != rightEntries[i].nameID {
			return false
		}
	}
	for i := range leftEntries {
		if !c.AreEqualCaching(leftEntries[i].GetType(), rightEntries[i].GetType()) {
			return false
		}
	}

	return true
}

// jsTypeStringSetsEqual is Set.equals for two sets of strings held as
// slices without duplicates.
//
// Go-only: language-adapter — getOwnPropertyNames() returns a slice.
func jsTypeStringSetsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	set := make(map[string]struct{}, len(b))
	for _, s := range b {
		set[s] = struct{}{}
	}
	for _, s := range a {
		if _, ok := set[s]; !ok {
			return false
		}
	}
	return true
}

// Java: EqualityChecker.areTypeMapEqual(TemplateTypeMap, TemplateTypeMap)
func (c *EqualityChecker) AreTypeMapEqual(left *TemplateTypeMap, right *TemplateTypeMap) bool {
	if left == right {
		return true
	}
	leftKeys := left.GetTemplateKeys()
	if len(leftKeys) == 0 {
		return true
	}
	rightKeys := right.GetTemplateKeys()
	if len(rightKeys) == 0 {
		return false
	}
	leftValues := left.GetTemplateValues()
	rightValues := right.GetTemplateValues()

outer:
	for i := 0; i < len(leftKeys); i++ {
		leftKey := leftKeys[i]

		for j := 0; j < len(rightKeys); j++ {
			rightKey := rightKeys[j]

			// Cross-compare every key-value pair in this TemplateTypeMap with
			// those in that TemplateTypeMap. Update the Equivalence match for
			// both key-value pairs involved.
			if leftKey != rightKey {
				continue
			}

			// The arrays are index-aligned but we can have more keys than values.
			var leftValue, rightValue JSType
			if i < len(leftValues) {
				leftValue = leftValues[i]
			} else {
				leftValue = left.DefaultValueType(leftKey)
			}
			if j < len(rightValues) {
				rightValue = rightValues[j]
			} else {
				rightValue = right.DefaultValueType(rightKey)
			}
			if c.AreEqualCaching(leftValue, rightValue) {
				continue outer
			}
		}

		return false
	}

	return true
}

// EqualityChecker_CacheKey is EqualityChecker.CacheKey.
type EqualityChecker_CacheKey struct {
	left     JSType
	right    JSType
	hashCode int32 // Cache this calculation because it is made often.
}

// Java: EqualityChecker.CacheKey.hashCode()
func (k *EqualityChecker_CacheKey) HashCode() int32 {
	return k.hashCode
}

// Java: EqualityChecker.CacheKey.equals(Object)
func (k *EqualityChecker_CacheKey) Equals(other any) bool {
	// Calling left with `null` or not a `Key` should cause a crash.
	right := other.(*EqualityChecker_CacheKey)
	if k == right {
		return true
	}

	// Recall right `Key` implements identity equality on `left` and `right`.
	//
	// Recall right `left` and `right` are not ordered.
	return (k.left == right.left && k.right == right.right) ||
		(k.left == right.right && k.right == right.left)
}

// Java: EqualityChecker.CacheKey(JSType, JSType)
func NewEqualityChecker_CacheKey(left JSType, right JSType) *EqualityChecker_CacheKey {
	// XOR the component hashcodes because:
	//   - It's a symmetric operator, so we don't have to worry about order.
	//   - It's assumed the inputs are already uniformly distributed and
	//     unrelated.
	//     - `left` and `right` should never be identical.
	// Recall right `Key` implements identity equality on `left` and `right`.
	return &EqualityChecker_CacheKey{
		left:     left,
		right:    right,
		hashCode: jutil.IdentityHashCode(left) ^ jutil.IdentityHashCode(right),
	}
}

// Java: EqualityChecker.unwrapNominalTypeProxies(ObjectType)
func EqualityChecker_UnwrapNominalTypeProxies(objType ObjectType) ObjectType {
	if !objType.IsResolved() || (!objType.IsNamedType() && !objType.IsTemplatizedType()) {
		// Don't unwrap TemplateTypes, as they should use identity semantics even
		// if their bounds are compatible. On the other hand, different
		// TemplatizedType instances may be equal if their TemplateTypeMaps are
		// compatible (which was checked for earlier).
		return objType
	}

	var internal ObjectType
	if objType.IsNamedType() {
		internal = objType.ToMaybeNamedType().GetReferencedObjTypeInternal()
	} else {
		internal = objType.ToMaybeTemplatizedType().GetReferencedObjTypeInternal()
	}
	return EqualityChecker_UnwrapNominalTypeProxies(internal)
}
