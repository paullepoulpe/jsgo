// Port of com/google/javascript/rhino/jstype/JSType.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// CanCastTo tests whether values of this type can be safely assigned to
// values of that type.
//
// The default implementation verifies that this is a subtype of that.
//
// Java: JSType.canCastTo(JSType)
func (t *JSTypeBase) CanCastTo(that JSType) bool {
	return canCastTo(t.self(), that)
}

// AutoboxesTo turns a scalar type to the corresponding object type. Returns
// the auto-boxed type or null if this type is not a scalar.
//
// Java: JSType.autoboxesTo()
func (t *JSTypeBase) AutoboxesTo() JSType {
	return nil
}

// Java: JSType.isBoxableScalar()
func (t *JSTypeBase) IsBoxableScalar() bool {
	return t.self().AutoboxesTo() != nil
}

// ToObjectType casts this to an ObjectType, or returns null if this is not
// an ObjectType. If this is a scalar type, it will *not* be converted to an
// object type. If you want to simulate JS autoboxing or dereferencing, you
// should use autoboxesTo() or dereference().
//
// Java: JSType.toObjectType()
func (t *JSTypeBase) ToObjectType() ObjectType {
	if t.objFlags&objKindMask >= jsKindUnion {
		return nil
	}
	if objectType, ok := t.self_.(ObjectType); ok {
		return objectType
	}
	return nil
}

// Autobox dereferences a type for property access.
//
// Filters null/undefined and autoboxes the resulting type. Never returns
// null.
//
// Java: JSType.autobox()
func (t *JSTypeBase) Autobox() JSType {
	restricted := t.self().RestrictByNotNullOrUndefined()
	autobox := restricted.AutoboxesTo()
	if autobox == nil {
		return restricted
	}
	return autobox
}

// Dereference dereferences a type for property access.
//
// Filters null/undefined, autoboxes the resulting type, and returns it iff
// it's an object. If not an object, returns null.
//
// Java: JSType.dereference()
func (t *JSTypeBase) Dereference() ObjectType {
	return t.self().Autobox().ToObjectType()
}

// CanTestForEqualityWith tests whether this and that are meaningfully
// comparable. By meaningfully, we mean compatible types that do not lead to
// step 22 of the definition of the Abstract Equality Comparison Algorithm
// (11.9.3, page 55–56) of the ECMA-262 specification.
//
// Java: JSType.canTestForEqualityWith(JSType)
func (t *JSTypeBase) CanTestForEqualityWith(that JSType) bool {
	result := t.self().TestForEquality(that)
	jsTypeCheckTriNotNull(result)
	return result == base.Tri_UNKNOWN
}

// jsType_TRI_NULL stands for a null Tri reference.
//
// Go-only: language-adapter — testForEquality returns a nullable Tri in
// Java: the base JSType implementation (testForEqualityHelper) returns null
// when it cannot decide, and the overrides test `result != null`. base.Tri is
// a value type, so null is this out-of-range value, which never escapes the
// package's own overrides.
const jsType_TRI_NULL base.Tri = -1

// jsTypeCheckTriNotNull throws NullPointerException where Java dereferences
// a null Tri.
func jsTypeCheckTriNotNull(t base.Tri) {
	if t == jsType_TRI_NULL {
		panic(jutil.NewNullPointerException())
	}
}

// TestForEquality compares this and that.
//
// Returns
//
//   - Tri#TRUE if the comparison of values of this type and that always
//     succeed (such as undefined compared to null)
//   - Tri#FALSE if the comparison of values of this type and that always
//     fails (such as undefined compared to number)
//   - Tri#UNKNOWN if the comparison can succeed or fail depending on the
//     concrete values
//
// Java: JSType.testForEquality(JSType)
func (t *JSTypeBase) TestForEquality(that JSType) base.Tri {
	return t.self().TestForEqualityHelper(t.self(), that)
}

// Java: JSType.testForEqualityHelper(JSType, JSType)
func (t *JSTypeBase) TestForEqualityHelper(aType JSType, bType JSType) base.Tri {
	if bType.IsAllType() || bType.IsUnknownType() ||
		bType.IsNoResolvedType() ||
		aType.IsAllType() || aType.IsUnknownType() ||
		aType.IsNoResolvedType() {
		return base.Tri_UNKNOWN
	}

	aIsEmpty := aType.IsEmptyType()
	bIsEmpty := bType.IsEmptyType()
	if aIsEmpty || bIsEmpty {
		if aIsEmpty && bIsEmpty {
			return base.Tri_TRUE
		} else {
			return base.Tri_UNKNOWN
		}
	}

	if bType.IsUnionType() {
		return bType.TestForEquality(aType)
	}

	if aType.IsFunctionType() || bType.IsFunctionType() {
		otherType := aType
		if aType.IsFunctionType() {
			otherType = bType
		}

		// TODO(johnlenz): tighten function type comparisons in general.
		if otherType.IsSymbol() {
			return base.Tri_FALSE
		}

		// In theory, functions are comparable to anything except null/undefined.
		// For example, on FF3:
		// function() {} == 'function () {\n}'
		// In practice, how a function serializes to a string is
		// implementation-dependent, so it does not really make sense to test
		// for equality with a string.
		greatestSubtype := otherType.GetGreatestSubtype(t.GetNativeType(JSTypeNative_OBJECT_TYPE))
		if greatestSubtype.IsNoType() || greatestSubtype.IsNoObjectType() {
			return base.Tri_FALSE
		} else {
			return base.Tri_UNKNOWN
		}
	}

	if bType.IsEnumElementType() {
		return bType.TestForEquality(aType)
	}

	// If this is a "Symbol" or that is "symbol" or "Symbol"
	if aType.IsSymbol() {
		if bType.CanCastTo(t.GetNativeType(JSTypeNative_SYMBOL_TYPE)) ||
			bType.CanCastTo(t.GetNativeType(JSTypeNative_SYMBOL_OBJECT_TYPE)) {
			return base.Tri_UNKNOWN
		}
		return base.Tri_FALSE
	}

	if bType.IsSymbol() {
		if aType.CanCastTo(t.GetNativeType(JSTypeNative_SYMBOL_TYPE)) ||
			aType.CanCastTo(t.GetNativeType(JSTypeNative_SYMBOL_OBJECT_TYPE)) {
			return base.Tri_UNKNOWN
		}
		return base.Tri_FALSE
	}

	return jsType_TRI_NULL
}

// CanTestForShallowEqualityWith tests whether this and that are
// meaningfully comparable using shallow comparison. By meaningfully, we mean
// compatible types that are not rejected by step 1 of the definition of the
// Strict Equality Comparison Algorithm (11.9.6, page 56–57) of the ECMA-262
// specification.
//
// Java: JSType.canTestForShallowEqualityWith(JSType)
func (t *JSTypeBase) CanTestForShallowEqualityWith(that JSType) bool {
	if t.self().IsNoResolvedType() || that.IsNoResolvedType() {
		return true
	}
	if t.self().IsEmptyType() || that.IsEmptyType() {
		return t.self().IsSubtypeOf(that) || that.IsSubtypeOf(t.self())
	}

	inf := t.self().GetGreatestSubtype(that)
	return !inf.IsEmptyType() ||
		// Our getGreatestSubtype relation on functions is pretty bad. Let's
		// just say it's always ok to compare two functions. Once the TODO in
		// FunctionType is fixed, we should be able to remove this.
		inf == t.GetRegistry().GetNativeType(JSTypeNative_LEAST_FUNCTION_TYPE)
}

// IsNullable tests whether this type is nullable.
//
// Java: JSType.isNullable()
func (t *JSTypeBase) IsNullable() bool {
	return false
}

// IsVoidable tests whether this type is voidable.
//
// Java: JSType.isVoidable()
func (t *JSTypeBase) IsVoidable() bool {
	return false
}

// IsExplicitlyVoidable tests whether this type explicitly allows undefined,
// as opposed to ? or *. This is required for a property to be optional.
//
// Java: JSType.isExplicitlyVoidable()
func (t *JSTypeBase) IsExplicitlyVoidable() bool {
	return false
}

// CollapseUnion gets the least supertype of this that's not a union.
//
// Java: JSType.collapseUnion()
func (t *JSTypeBase) CollapseUnion() JSType {
	return t.self()
}

// GetLeastSupertype gets the least supertype of this and that. The least
// supertype is the join (∨) or supremum of both types in the type lattice.
//
// Examples:
//
//   - number ∨ * = *
//   - number ∨ Object = (number, Object)
//   - Number ∨ Object = Object
//
// Java: JSType.getLeastSupertype(JSType)
func (t *JSTypeBase) GetLeastSupertype(that JSType) JSType {
	if t.self() == that {
		return t.self()
	}

	if that.IsUnionType() {
		// Union types have their own implementation of getLeastSupertype.
		return that.ToMaybeUnionType().GetLeastSupertype(t.self())
	}
	return JSType_GetLeastSupertype(t.self(), that)
}

// JSType_GetLeastSupertype is a generic implementation meant to be used as a
// helper for common getLeastSupertype implementations.
//
// Java: JSType.getLeastSupertype(JSType, JSType)
func JSType_GetLeastSupertype(thisType JSType, thatType JSType) JSType {
	thisTb := jsTypeBase(thisType)
	thatTb := jsTypeBase(thatType)
	if thisTb == thatTb {
		return thisType
	}
	registry := thisTb.GetRegistry()
	if registry != nil && registry.resolver.state == JSTypeResolver_State_CLOSED && thisTb.flags&jsTypeFlagResolvedToSelf != 0 && thatTb != nil {
		if thatTb.flags&jsTypeFlagResolvedToSelf != 0 {
			if cached := registry.lookupLeastSupertype(thisTb, thatTb); cached != nil {
				return cached
			}
			areEquivalent := equalityChecker_AreIdentityEqualWithBases(thisType, thisTb, thatType, thatTb)
			if areEquivalent {
				registry.storeLeastSupertype(thisTb, thatTb, thisType)
				return thisType
			}
			res := registry.CreateUnionType(thisType, thatType)
			resTb := jsTypeBase(res)
			if resTb == thisTb || resTb == thatTb {
				registry.storeLeastSupertype(thisTb, thatTb, res)
			}
			return res
		}
	}
	areEquivalent := jsTypeEquals(thisType, thatType)
	if areEquivalent {
		return thisType
	}
	return registry.CreateUnionType(thisType, thatType)
}

// GetGreatestSubtype gets the greatest subtype of this and that. The
// greatest subtype is the meet (∧) or infimum of both types in the type
// lattice.
//
// Examples
//
//   - Number ∧ Any = Any
//   - number ∧ Object = Any
//   - Number ∧ Object = Number
//
// Java: JSType.getGreatestSubtype(JSType)
func (t *JSTypeBase) GetGreatestSubtype(that JSType) JSType {
	return JSType_GetGreatestSubtype(t.self(), that)
}

// JSType_GetGreatestSubtype is a generic implementation meant to be used as
// a helper for common getGreatestSubtype implementations.
//
// Java: JSType.getGreatestSubtype(JSType, JSType)
func JSType_GetGreatestSubtype(thisType JSType, thatType JSType) JSType {
	if thisType.IsFunctionType() && thatType.IsFunctionType() {
		// The FunctionType sub-lattice is not well-defined. i.e., the
		// proposition
		// A < B => sup(A, B) == B
		// does not hold because of unknown parameters and return types.
		// See the comment in supAndInfHelper for more info on this.
		return thisType.ToMaybeFunctionType().SupAndInfHelper(thatType.ToMaybeFunctionType(), false)
	} else if thisType.Equals(thatType) {
		return thisType
	} else if thisType.IsUnknownType() {
		// Handle unknowns by returning the tightest type possible:
		//  given a known and unknown type, return the known type
		//  given a checked unknown type and some unknown type, return a checked unknown
		//  otherwise return an unchecked unknown
		if !thatType.IsUnknownType() {
			return thatType
		} else if thisType.IsCheckedUnknownType() || thatType.IsCheckedUnknownType() {
			return thisType.GetNativeType(JSTypeNative_CHECKED_UNKNOWN_TYPE)
		} else {
			return thisType.GetNativeType(JSTypeNative_UNKNOWN_TYPE)
		}
	} else if thatType.IsUnknownType() {
		// thisType must be known
		return thisType
	} else if thisType.IsUnionType() {
		return UnionType_GetGreatestSubtype(thisType.ToMaybeUnionType(), thatType)
	} else if thatType.IsUnionType() {
		return UnionType_GetGreatestSubtype(thatType.ToMaybeUnionType(), thisType)
	} else if thisType.IsTemplatizedType() {
		return thisType.ToMaybeTemplatizedType().GetGreatestSubtypeHelper(thatType)
	} else if thatType.IsTemplatizedType() {
		return thatType.ToMaybeTemplatizedType().GetGreatestSubtypeHelper(thisType)
	} else if thisType.IsSubtypeOf(thatType) {
		return thisType
	} else if thatType.IsSubtypeOf(thisType) {
		return thatType
	} else if thisType.IsRecordType() {
		return thisType.ToMaybeRecordType().GetGreatestSubtypeHelper(thatType)
	} else if thatType.IsRecordType() {
		return thatType.ToMaybeRecordType().GetGreatestSubtypeHelper(thisType)
	}

	if thisType.IsEnumElementType() {
		inf := EnumElementType_GetGreatestSubtype(thisType.ToMaybeEnumElementType(), thatType)
		if inf != nil {
			return inf
		}
	} else if thatType.IsEnumElementType() {
		inf := EnumElementType_GetGreatestSubtype(thatType.ToMaybeEnumElementType(), thisType)
		if inf != nil {
			return inf
		}
	}

	if thisType.IsObject() && thatType.IsObject() {
		return thisType.GetNativeType(JSTypeNative_NO_OBJECT_TYPE)
	}
	return thisType.GetNativeType(JSTypeNative_NO_TYPE)
}

// GetRestrictedTypeGivenOutcome computes the restricted type of this type
// knowing that the ToBoolean predicate has a specific value. For more
// information about the ToBoolean predicate, see
// #getPossibleToBooleanOutcomes.
//
// Returns the restricted type, or the Any Type if the underlying type could
// not have yielded this ToBoolean value.
//
// TODO(user): Move this method to the SemanticRAI and use the visit method
// of types to get the restricted type.
//
// Java: JSType.getRestrictedTypeGivenOutcome(Outcome)
func (t *JSTypeBase) GetRestrictedTypeGivenOutcome(outcome Outcome) JSType {
	if outcome.IsTruthy() && t.self() == t.GetNativeType(JSTypeNative_UNKNOWN_TYPE) {
		return t.GetNativeType(JSTypeNative_CHECKED_UNKNOWN_TYPE)
	}

	// Only `NULL_TYPE` and `VOID_TYPE` can really be nullish, so the result
	// for any other non-union type must be NO_TYPE. Note that the UnionType
	// class is responsible for handling the union case.
	if outcome.IsNullish().ToBoolean(false) {
		if t.self() == t.GetNativeType(JSTypeNative_VOID_TYPE) ||
			t.self() == t.GetNativeType(JSTypeNative_NULL_TYPE) {
			return t.self()
		} else {
			return t.GetNativeType(JSTypeNative_NO_TYPE)
		}
	}

	literals := t.self().GetPossibleToBooleanOutcomes()
	if literals.Contains(outcome.IsTruthy()) {
		return t.self()
	} else {
		return t.GetNativeType(JSTypeNative_NO_TYPE)
	}
}

// GetTypesUnderEquality computes the subset of this and that types if
// equality is observed. If a value v1 of type null is equal to a value v2 of
// type (undefined,number), we can infer that the type of v1 is null and the
// type of v2 is undefined.
//
// Returns a pair containing the restricted type of this as the first
// component and the restricted type of that as the second element. The
// returned pair is never null even though its components may be null
//
// Java: JSType.getTypesUnderEquality(JSType)
func (t *JSTypeBase) GetTypesUnderEquality(that JSType) *JSType_TypePair {
	a, b := t.GetTypesUnderEqualityPair(that)
	return NewJSType_TypePair(a, b)
}

func (t *JSTypeBase) GetTypesUnderEqualityPair(that JSType) (JSType, JSType) {
	// unions types
	if that.IsUnionType() {
		a, b := that.ToMaybeUnionType().GetTypesUnderEqualityPair(t.self())
		return b, a
	}

	// other types
	r := t.self().TestForEquality(that)
	jsTypeCheckTriNotNull(r)
	switch r {
	case base.Tri_FALSE:
		return nil, nil
	}
	return t.self(), that // TRUE, UNKNOWN
}

// GetTypesUnderInequality computes the subset of this and that types if
// inequality is observed. If a value v1 of type number is not equal to a
// value v2 of type (undefined,number), we can infer that the type of v1 is
// number and the type of v2 is number as well.
//
// Java: JSType.getTypesUnderInequality(JSType)
func (t *JSTypeBase) GetTypesUnderInequality(that JSType) *JSType_TypePair {
	a, b := t.GetTypesUnderInequalityPair(that)
	return NewJSType_TypePair(a, b)
}

func (t *JSTypeBase) GetTypesUnderInequalityPair(that JSType) (JSType, JSType) {
	// unions types
	if that.IsUnionType() {
		a, b := that.ToMaybeUnionType().GetTypesUnderInequalityPair(t.self())
		return b, a
	}

	// For the remaining types, the only way to restrict anything is if both
	// types are null or void, in which case we're left with nothing.
	if t.IsNullTypeOrVoidType() && jsTypeIsNullTypeOrVoidType(that) {
		noType := t.GetRegistry().GetNativeType(JSTypeNative_NO_TYPE)
		return noType, noType
	}
	return t.self(), that
}

// GetTypesUnderShallowEquality computes the subset of this and that types
// under shallow equality.
//
// Java: JSType.getTypesUnderShallowEquality(JSType)
func (t *JSTypeBase) GetTypesUnderShallowEquality(that JSType) *JSType_TypePair {
	a, b := t.GetTypesUnderShallowEqualityPair(that)
	return NewJSType_TypePair(a, b)
}

func (t *JSTypeBase) GetTypesUnderShallowEqualityPair(that JSType) (JSType, JSType) {
	commonType := t.self().GetGreatestSubtype(that)
	return commonType, commonType
}

// GetTypesUnderShallowInequality computes the subset of this and that types
// under shallow inequality.
//
// Java: JSType.getTypesUnderShallowInequality(JSType)
func (t *JSTypeBase) GetTypesUnderShallowInequality(that JSType) *JSType_TypePair {
	a, b := t.GetTypesUnderShallowInequalityPair(that)
	return NewJSType_TypePair(a, b)
}

func (t *JSTypeBase) GetTypesUnderShallowInequalityPair(that JSType) (JSType, JSType) {
	// union types
	if that.IsUnionType() {
		a, b := that.ToMaybeUnionType().GetTypesUnderShallowInequalityPair(t.self())
		return b, a
	}

	// Other types. There are only two types whose shallow inequality is
	// deterministically true -- null and undefined. We can just enumerate
	// them.
	if t.self() == that && t.IsNullTypeOrVoidType() {
		return nil, nil
	}
	// Since unions have already been removed, if this and that don't have
	// the same unit type then there is no further narrowing that can be done
	// on them.
	return t.self(), that
}

// Java: JSType.IsNullTypeOrVoidType()
func (t *JSTypeBase) IsNullTypeOrVoidType() bool {
	// TODO(sdh): Should this be exposed more publicly?
	return t.self().IsNullType() || t.self().IsVoidType()
}

// jsTypeIsNullTypeOrVoidType is that.IsNullTypeOrVoidType() (a private
// method called on another instance).
func jsTypeIsNullTypeOrVoidType(that JSType) bool {
	return that.GetJSTypeBase().IsNullTypeOrVoidType()
}

// Java: JSType.getUnionMembers()
func (t *JSTypeBase) GetUnionMembers() []JSType {
	if t.self().IsUnionType() {
		return t.self().ToMaybeUnionType().GetAlternates()
	}
	return nil
}

// RestrictByNotNullOrUndefined returns, if this is a union type, a union
// type that does not include the null or undefined type.
//
// Java: JSType.restrictByNotNullOrUndefined()
func (t *JSTypeBase) RestrictByNotNullOrUndefined() JSType {
	return t.self()
}

// RestrictByNotUndefined returns, if this is a union type, a union type
// that does not include the undefined type.
//
// Java: JSType.restrictByNotUndefined()
func (t *JSTypeBase) RestrictByNotUndefined() JSType {
	return t.self()
}

// RestrictByNotNull returns, if this is a union type, a union type that does
// not include the null type.
//
// Java: JSType.restrictByNotNull()
func (t *JSTypeBase) RestrictByNotNull() JSType {
	return t.self()
}

// IsSubtypeWithoutStructuralTyping: the logic of this method is similar to
// isSubtype, except that it does not perform structural interface matching
//
// This function is added for disambiguate properties, and is deprecated for
// the other use cases.
//
// Java: JSType.isSubtypeWithoutStructuralTyping(JSType)
func (t *JSTypeBase) IsSubtypeWithoutStructuralTyping(supertype JSType) bool {
	supBase := jsTypeBase(supertype)
	if supBase == nil {
		failJSTypeNull()
	}
	if t == supBase {
		return true
	}
	c := SubtypeChecker{
		initialSubtype:        t.self(),
		initialSupertype:      supertype,
		registry:              t.GetRegistry(),
		usingStructuralTyping: optBoolFalse,
		subtypingMode:         optSubtypingModeNormal,
		hasRun:                true,
	}
	return c.IsSubtypeDispatching(t.self(), supertype)
}

// JSType_SubtypingMode: in files translated from Java, we typecheck null
// and undefined loosely.
type JSType_SubtypingMode int32

const (
	JSType_SubtypingMode_NORMAL JSType_SubtypingMode = iota
	JSType_SubtypingMode_IGNORE_NULL_UNDEFINED
)

var jsType_SubtypingMode_names = [...]string{"NORMAL", "IGNORE_NULL_UNDEFINED"}

// String is Enum.name().
func (m JSType_SubtypingMode) String() string { return jsType_SubtypingMode_names[m] }

// IsSubtype checks whether this is a subtype of that.
//
// Note this function also returns true if this type structurally matches
// the protocol define by that type (if that type is an interface function
// type)
//
// Subtyping rules:
//
//   - (unknown) — every type is a subtype of the Unknown type.
//   - (no) — the No type is a subtype of every type.
//   - (no-object) — the NoObject type is a subtype of every object type
//     (i.e. subtypes of the Object type).
//   - (ref) — a type is a subtype of itself.
//   - (union-l) — A union type is a subtype of a type U if all the union
//     type's constituents are a subtype of U. Formally (T1, …, Tn) <: U if
//     and only Tk <: U for all k ∈ 1..n.
//   - (union-r) — A type U is a subtype of a union type if it is a subtype
//     of one of the union type's constituents. Formally U <: (T1, …, Tn) if
//     and only if U <: Tk for some index k.
//   - (objects) — an Object O1 is a subtype of an object O2 if it has more
//     properties than O2 and all common properties are pairwise subtypes.
//
// Returns this <: that.
//
// Java: JSType.isSubtype(JSType)
func (t *JSTypeBase) IsSubtype(supertype JSType) bool {
	return t.IsSubtypeOf(supertype)
}

// IsSubtypeWithMode is deprecated: prefer isSubtype(JSType) instead.
//
// Java: JSType.isSubtype(JSType, SubtypingMode)
func (t *JSTypeBase) IsSubtypeWithMode(supertype JSType, mode JSType_SubtypingMode) bool {
	return t.IsSubtypeOfWithMode(supertype, mode)
}

// Java: JSType.isSubtypeOf(JSType)
func (t *JSTypeBase) IsSubtypeOf(supertype JSType) bool {
	supBase := jsTypeBase(supertype)
	if supBase == nil {
		failJSTypeNull()
	}
	if t == supBase {
		return true
	}
	c := SubtypeChecker{
		initialSubtype:        t.self(),
		initialSupertype:      supertype,
		registry:              t.GetRegistry(),
		usingStructuralTyping: optBoolTrue,
		subtypingMode:         optSubtypingModeNormal,
		hasRun:                true,
	}
	return c.IsSubtypeDispatching(t.self(), supertype)
}

// Java: JSType.isSubtypeOf(JSType, SubtypingMode)
func (t *JSTypeBase) IsSubtypeOfWithMode(supertype JSType, mode JSType_SubtypingMode) bool {
	supBase := jsTypeBase(supertype)
	if supBase == nil {
		failJSTypeNull()
	}
	if t == supBase {
		return true
	}
	c := SubtypeChecker{
		initialSubtype:        t.self(),
		initialSupertype:      supertype,
		registry:              t.GetRegistry(),
		usingStructuralTyping: optBoolTrue,
		subtypingMode:         optSubtypingModeOf(mode),
		hasRun:                true,
	}
	return c.IsSubtypeDispatching(t.self(), supertype)
}

// Resolve resolves this type in the given scope.
//
// The returned value must be equal to this, as defined by #equals. It may
// or may not be the same object. This method may modify the internal state
// of this, as long as it does so in a way that preserves Object equality.
//
// For efficiency, we should only resolve a type once per compilation job.
// For incremental compilations, one compilation job may need the artifacts
// from a previous generation, so we will eventually need a generational flag
// instead of a boolean one.
//
// Java: JSType.resolve(ErrorReporter)
func (t *JSTypeBase) Resolve(reporter ErrorReporter) JSType {
	t.GetRegistry().GetResolver().AssertLegalToResolveTypes()
	if t.resolveResult == nil {
		// Prevent infinite recursion in cyclically defined types.
		//
		// If resolve is called a twice on a type, before the first call
		// completes, there is a cycle somewhere. The only thing we can do is
		// return `this`. Doing so is safe because the resolve result *should*
		// be equal to `this`, even though the different calls will return
		// different objects.
		t.resolveResult = t
		res := t.self().ResolveInternal(reporter)
		if res != nil {
			t.resolveResult = res.GetJSTypeBase()
		} else {
			t.resolveResult = nil
		}
		if t.resolveResult == t {
			t.flags |= jsTypeFlagResolvedToSelf
		}
		jutil.CheckState(t.resolveResult != nil)
	}
	return t.resolveResult.self()
}

// EagerlyResolveToSelf bypasses the normal call to resolve(ErrorReporter).
// Useful if you are aware that no resolution work needs to be done.
//
// Java: JSType.eagerlyResolveToSelf()
func (t *JSTypeBase) EagerlyResolveToSelf() {
	jutil.CheckStateMsg1(!t.self().IsResolved(), "%s", t.self())
	t.resolveResult = t
	t.flags |= jsTypeFlagResolvedToSelf
	t.GetRegistry().GetResolver().ResolveIfClosed(t.self(), t.self().GetTypeClass())
}

// IsResolved returns whether the type has undergone resolution.
//
// A value of true does not indicate that resolution was successful, only
// that it was attempted and has finished.
//
// Java: JSType.isResolved()
func (t *JSTypeBase) IsResolved() bool {
	return t.resolveResult != nil
}

// IsSuccessfullyResolved returns whether the type has undergone resolution
// and resolved to a "useful" type.
//
// Java: JSType.isSuccessfullyResolved()
func (t *JSTypeBase) IsSuccessfullyResolved() bool {
	return t.self().IsResolved() && !t.self().IsNoResolvedType()
}

// IsUnsuccessfullyResolved returns whether the type has undergone
// resolution and resolved to a "useless" type.
//
// Java: JSType.isUnsuccessfullyResolved()
func (t *JSTypeBase) IsUnsuccessfullyResolved() bool {
	return t.self().IsResolved() && t.self().IsNoResolvedType()
}

// JSType_SafeResolve is a null-safe resolve.
//
// Java: JSType.safeResolve(JSType, ErrorReporter)
func JSType_SafeResolve(type_ JSType, reporter ErrorReporter) JSType {
	if type_ == nil {
		return nil
	}
	return type_.Resolve(reporter)
}

// SetValidator: certain types have constraints on them at resolution-time.
// For example, a type in an @extends annotation must be an object. Clients
// should inject a validator that emits a warning if the type does not
// validate, and return false.
//
// Java: JSType.setValidator(Predicate<JSType>)
func (t *JSTypeBase) SetValidator(validator func(JSType) bool) bool {
	return validator(t.self())
}

// JSType_TypePair represents a pair of types.
type JSType_TypePair struct {
	TypeA JSType
	TypeB JSType
}

// Java: JSType.TypePair(JSType, JSType)
func NewJSType_TypePair(typeA JSType, typeB JSType) *JSType_TypePair {
	return &JSType_TypePair{TypeA: typeA, TypeB: typeB}
}

// MatchConstraint modifies this type so that it matches the specified type.
//
// This is useful for reverse type-inference, where we want to infer that an
// object literal matches its constraint (much like how the java compiler
// does reverse-inference to figure out generics).
//
// Java: JSType.matchConstraint(JSType)
func (t *JSTypeBase) MatchConstraint(constraint JSType) {}

// JSType_MatchStatus describes the status of checking that a function
// implicitly implements an interface.
//
// It also be used to describe the status of checking that a record type
// structurally matches another record type
//
// A function implicitly implements an interface if the function does not use
// @implements to declare that it implements the interface, but its class
// structure complies with the protocol defined by the interface
type JSType_MatchStatus int32

const (
	// JSType_MatchStatus_MATCH indicates that a function implicitly
	// implements an interface (i.e., the function structurally complies with
	// the protocol defined in interface) or a record type matches another
	// record type
	JSType_MatchStatus_MATCH JSType_MatchStatus = iota
	// JSType_MatchStatus_NOT_MATCH indicates that a function does not
	// implicitly implements an interface (i.e., the function does not
	// structurally comply with the protocol defined in interface) or a record
	// type does not match another record type
	JSType_MatchStatus_NOT_MATCH
	// JSType_MatchStatus_PROCESSING indicates that the interface and function
	// relationship is under processing
	JSType_MatchStatus_PROCESSING
)

var jsType_MatchStatus_names = [...]string{"MATCH", "NOT_MATCH", "PROCESSING"}

// jsType_MatchStatus_isSubtype is the isSubtype field of each constant.
var jsType_MatchStatus_isSubtype = [...]bool{true, false, true}

// String is Enum.name().
func (s JSType_MatchStatus) String() string { return jsType_MatchStatus_names[s] }

// Java: JSType.MatchStatus.subtypeValue()
func (s JSType_MatchStatus) SubtypeValue() bool {
	return jsType_MatchStatus_isSubtype[s]
}

// Java: JSType.MatchStatus.valueOf(boolean)
func JSType_MatchStatus_ValueOf(match bool) JSType_MatchStatus {
	if match {
		return JSType_MatchStatus_MATCH
	}
	return JSType_MatchStatus_NOT_MATCH
}
