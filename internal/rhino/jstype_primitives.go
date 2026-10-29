// Port of com/google/javascript/rhino/jstype/ValueType.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/jstype/BooleanType.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/jstype/NumberType.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/jstype/StringType.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/jstype/BigIntType.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/jstype/NullType.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/jstype/VoidType.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/jstype/SymbolType.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/jstype/KnownSymbolType.java @ closure-compiler v20261021

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
 *   Brock Smickley
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

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
 *   Nick Santos
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

// ValueTypeBase holds ValueType's base method implementations (for null, void,
// number, boolean, string, bigint, symbol); the virtual calls go through
// JSTypeBase.self().
//
// Java: abstract class ValueType extends JSType.
type ValueTypeBase struct {
	JSTypeBase
}

// InitValueTypeBase is the Go form of the ValueType constructor.
//
// Java: ValueType(JSTypeRegistry)
func (v *ValueTypeBase) InitValueTypeBase(self JSType, registry *JSTypeRegistry) {
	v.InitJSTypeBase(self, registry)
	v.EagerlyResolveToSelf()
}

// Java: ValueType.resolveInternal(ErrorReporter)
func (v *ValueTypeBase) ResolveInternal(reporter ErrorReporter) JSType {
	panic(jutil.NewAssertionErrorNoMessage())
}

// Java: ValueType.appendTo(TypeStringBuilder)
func (v *ValueTypeBase) AppendTo(sb *TypeStringBuilder) {
	sb.Append(jutil.ToJavaString(v.self().GetDisplayName()))
}

// Java: ValueType.hasDisplayName()
func (v *ValueTypeBase) HasDisplayName() bool {
	return true
}

// Java: ValueType.getPropertyKind(Property.Key, boolean)
func (v *ValueTypeBase) GetPropertyKindByKeyWithAutobox(propertyName Property_Key, autobox bool) JSType_HasPropertyKind {
	if autobox && v.self().IsBoxableScalar() {
		return v.self().AutoboxesTo().GetPropertyKindByKeyWithAutobox(propertyName, autobox)
	} else {
		return v.JSTypeBase.GetPropertyKindByKeyWithAutobox(propertyName, autobox)
	}
}

// Java: ValueType.recursionUnsafeHashCode()
func (v *ValueTypeBase) RecursionUnsafeHashCode() int32 {
	// Subclasses of this type are unique within a JSTypeRegistry.
	return jutil.IdentityHashCode(v.self())
}

// BooleanType is the boolean type.
type BooleanType struct {
	ValueTypeBase
}

// Java: BooleanType(JSTypeRegistry)
func NewBooleanType(registry *JSTypeRegistry) *BooleanType {
	t := &BooleanType{}
	t.objFlags = jsKindScalarValue
	t.InitValueTypeBase(t, registry)
	return t
}

// Java: BooleanType.getTypeClass()
func (t *BooleanType) GetTypeClass() JSTypeClass {
	return JSTypeClass_BOOLEAN
}

// Java: BooleanType.testForEquality(JSType)
func (t *BooleanType) TestForEquality(that JSType) base.Tri {
	result := t.ValueTypeBase.TestForEquality(that)
	if result != jsType_TRI_NULL {
		return result
	}
	if that.IsUnknownType() ||
		that.IsSubtypeOf(t.GetNativeType(JSTypeNative_NUMBER_STRING_BOOLEAN)) ||
		that.IsSubtypeOf(t.GetNativeType(JSTypeNative_BIGINT_TYPE)) ||
		that.IsObject() {
		return base.Tri_UNKNOWN
	}
	return base.Tri_FALSE
}

// Java: BooleanType.isBooleanValueType()
func (t *BooleanType) IsBooleanValueType() bool {
	return true
}

// Java: BooleanType.matchesNumberContext()
func (t *BooleanType) MatchesNumberContext() bool {
	return true
}

// Java: BooleanType.matchesStringContext()
func (t *BooleanType) MatchesStringContext() bool {
	return true
}

// Java: BooleanType.matchesObjectContext()
func (t *BooleanType) MatchesObjectContext() bool {
	// TODO(user): Revisit this for ES4, which is stricter.
	return true
}

// Java: BooleanType.autoboxesTo()
func (t *BooleanType) AutoboxesTo() JSType {
	return t.GetNativeType(JSTypeNative_BOOLEAN_OBJECT_TYPE)
}

// Java: BooleanType.getDisplayName()
func (t *BooleanType) GetDisplayName() *string {
	return &jsType_displayName_boolean
}

// Java: BooleanType.getPossibleToBooleanOutcomes()
func (t *BooleanType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return BooleanLiteralSet_BOTH
}

var jsType_displayName_boolean = "boolean"

// NumberType is the number type.
type NumberType struct {
	ValueTypeBase
}

// Java: NumberType(JSTypeRegistry)
func NewNumberType(registry *JSTypeRegistry) *NumberType {
	t := &NumberType{}
	t.objFlags = jsKindScalarValue
	t.InitValueTypeBase(t, registry)
	return t
}

// Java: NumberType.getTypeClass()
func (t *NumberType) GetTypeClass() JSTypeClass {
	return JSTypeClass_NUMBER
}

// Java: NumberType.testForEquality(JSType)
func (t *NumberType) TestForEquality(that JSType) base.Tri {
	result := t.ValueTypeBase.TestForEquality(that)
	if result != jsType_TRI_NULL {
		return result
	}
	if that.IsUnknownType() ||
		that.IsSubtypeOf(t.GetNativeType(JSTypeNative_OBJECT_TYPE)) ||
		that.IsSubtypeOf(t.GetNativeType(JSTypeNative_NUMBER_TYPE)) ||
		that.IsSubtypeOf(t.GetNativeType(JSTypeNative_STRING_TYPE)) ||
		that.IsSubtypeOf(t.GetNativeType(JSTypeNative_BOOLEAN_TYPE)) ||
		that.IsSubtypeOf(t.GetNativeType(JSTypeNative_BIGINT_TYPE)) {
		return base.Tri_UNKNOWN
	}
	return base.Tri_FALSE
}

// Java: NumberType.isNumberValueType()
func (t *NumberType) IsNumberValueType() bool {
	return true
}

// Java: NumberType.matchesNumberContext()
func (t *NumberType) MatchesNumberContext() bool {
	return true
}

// Java: NumberType.matchesStringContext()
func (t *NumberType) MatchesStringContext() bool {
	return true
}

// Java: NumberType.matchesObjectContext()
func (t *NumberType) MatchesObjectContext() bool {
	// TODO(user): Revisit this for ES4, which is stricter.
	return true
}

// Java: NumberType.getDisplayName()
func (t *NumberType) GetDisplayName() *string {
	return &jsType_displayName_number
}

// Java: NumberType.getPossibleToBooleanOutcomes()
func (t *NumberType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return BooleanLiteralSet_BOTH
}

// Java: NumberType.autoboxesTo()
func (t *NumberType) AutoboxesTo() JSType {
	return t.GetNativeType(JSTypeNative_NUMBER_OBJECT_TYPE)
}

var jsType_displayName_number = "number"

// StringType is the string type.
type StringType struct {
	ValueTypeBase
}

// Java: StringType(JSTypeRegistry)
func NewStringType(registry *JSTypeRegistry) *StringType {
	t := &StringType{}
	t.objFlags = jsKindScalarValue
	t.InitValueTypeBase(t, registry)
	return t
}

// Java: StringType.getTypeClass()
func (t *StringType) GetTypeClass() JSTypeClass {
	return JSTypeClass_STRING
}

// Java: StringType.testForEquality(JSType)
func (t *StringType) TestForEquality(that JSType) base.Tri {
	result := t.ValueTypeBase.TestForEquality(that)
	if result != jsType_TRI_NULL {
		return result
	}
	if that.IsUnknownType() ||
		that.IsSubtypeOf(t.GetNativeType(JSTypeNative_OBJECT_TYPE)) ||
		that.IsSubtypeOf(t.GetNativeType(JSTypeNative_NUMBER_TYPE)) ||
		that.IsSubtypeOf(t.GetNativeType(JSTypeNative_STRING_TYPE)) ||
		that.IsSubtypeOf(t.GetNativeType(JSTypeNative_BOOLEAN_TYPE)) ||
		that.IsSubtypeOf(t.GetNativeType(JSTypeNative_BIGINT_TYPE)) {
		return base.Tri_UNKNOWN
	}
	return base.Tri_FALSE
}

// Java: StringType.isStringValueType()
func (t *StringType) IsStringValueType() bool {
	return true
}

// Java: StringType.matchesNumberContext()
func (t *StringType) MatchesNumberContext() bool {
	return true
}

// Java: StringType.matchesStringContext()
func (t *StringType) MatchesStringContext() bool {
	return true
}

// Java: StringType.matchesObjectContext()
func (t *StringType) MatchesObjectContext() bool {
	// TODO(user): Revisit this for ES4, which is stricter.
	return true
}

// Java: StringType.getDisplayName()
func (t *StringType) GetDisplayName() *string {
	return &jsType_displayName_string
}

// Java: StringType.autoboxesTo()
func (t *StringType) AutoboxesTo() JSType {
	return t.GetNativeType(JSTypeNative_STRING_OBJECT_TYPE)
}

// Java: StringType.getPossibleToBooleanOutcomes()
func (t *StringType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return BooleanLiteralSet_BOTH
}

var jsType_displayName_string = "string"

// BigIntType is the bigint type.
type BigIntType struct {
	ValueTypeBase
}

// Java: BigIntType(JSTypeRegistry)
func NewBigIntType(registry *JSTypeRegistry) *BigIntType {
	t := &BigIntType{}
	t.objFlags = jsKindScalarValue
	t.InitValueTypeBase(t, registry)
	return t
}

// Java: BigIntType.getTypeClass()
func (t *BigIntType) GetTypeClass() JSTypeClass {
	return JSTypeClass_BIGINT
}

// Java: BigIntType.testForEquality(JSType)
func (t *BigIntType) TestForEquality(that JSType) base.Tri {
	result := t.ValueTypeBase.TestForEquality(that)
	if result != jsType_TRI_NULL {
		return result
	}
	if that.IsUnknownType() ||
		that.IsSubtypeOf(t.GetNativeType(JSTypeNative_OBJECT_TYPE)) ||
		that.IsSubtypeOf(t.GetNativeType(JSTypeNative_NUMBER_TYPE)) ||
		that.IsSubtypeOf(t.GetNativeType(JSTypeNative_STRING_TYPE)) ||
		that.IsSubtypeOf(t.GetNativeType(JSTypeNative_BOOLEAN_TYPE)) ||
		that.IsSubtypeOf(t.GetNativeType(JSTypeNative_BIGINT_TYPE)) {
		return base.Tri_UNKNOWN
	}
	return base.Tri_FALSE
}

// Java: BigIntType.isBigIntValueType()
func (t *BigIntType) IsBigIntValueType() bool {
	return true
}

// Java: BigIntType.matchesNumberContext()
func (t *BigIntType) MatchesNumberContext() bool {
	// BigInt cannot be coerced to a Number
	return false
}

// Java: BigIntType.matchesStringContext()
func (t *BigIntType) MatchesStringContext() bool {
	return true
}

// Java: BigIntType.matchesObjectContext()
func (t *BigIntType) MatchesObjectContext() bool {
	return true
}

// Java: BigIntType.getDisplayName()
func (t *BigIntType) GetDisplayName() *string {
	return &jsType_displayName_bigint
}

// Java: BigIntType.getPossibleToBooleanOutcomes()
func (t *BigIntType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return BooleanLiteralSet_BOTH
}

// Java: BigIntType.autoboxesTo()
func (t *BigIntType) AutoboxesTo() JSType {
	return t.GetNativeType(JSTypeNative_BIGINT_OBJECT_TYPE)
}

var jsType_displayName_bigint = "bigint"

// NullType is the null type.
type NullType struct {
	ValueTypeBase
}

// Java: NullType(JSTypeRegistry)
func NewNullType(registry *JSTypeRegistry) *NullType {
	t := &NullType{}
	t.objFlags = jsKindNull
	t.InitValueTypeBase(t, registry)
	return t
}

// Java: NullType.getTypeClass()
func (t *NullType) GetTypeClass() JSTypeClass {
	return JSTypeClass_NULL
}

// Java: NullType.isNullType()
func (t *NullType) IsNullType() bool {
	return true
}

// Java: NullType.isNullable()
func (t *NullType) IsNullable() bool {
	return true
}

// Java: NullType.matchesNumberContext()
func (t *NullType) MatchesNumberContext() bool {
	return true
}

// Java: NullType.matchesObjectContext()
func (t *NullType) MatchesObjectContext() bool {
	return false
}

// Java: NullType.matchesStringContext()
func (t *NullType) MatchesStringContext() bool {
	return true
}

// Java: NullType.restrictByNotNullOrUndefined()
func (t *NullType) RestrictByNotNullOrUndefined() JSType {
	return t.GetRegistry().GetNativeType(JSTypeNative_NO_TYPE)
}

// Java: NullType.restrictByNotNull()
func (t *NullType) RestrictByNotNull() JSType {
	return t.GetRegistry().GetNativeType(JSTypeNative_NO_TYPE)
}

// Java: NullType.testForEquality(JSType)
func (t *NullType) TestForEquality(that JSType) base.Tri {
	result := t.ValueTypeBase.TestForEquality(that)
	if result != jsType_TRI_NULL {
		return result
	}
	if that.IsNullType() || that.IsVoidType() {
		return base.Tri_TRUE
	}
	if that.IsUnknownType() || that.IsNullable() {
		return base.Tri_UNKNOWN
	}
	return base.Tri_FALSE
}

// Java: NullType.getDisplayName()
func (t *NullType) GetDisplayName() *string {
	return &jsType_displayName_null
}

// Java: NullType.getPossibleToBooleanOutcomes()
func (t *NullType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return BooleanLiteralSet_FALSE
}

var jsType_displayName_null = "null"

// VoidType is the void type.
type VoidType struct {
	ValueTypeBase
}

// Java: VoidType(JSTypeRegistry)
func NewVoidType(registry *JSTypeRegistry) *VoidType {
	t := &VoidType{}
	t.objFlags = jsKindVoid
	t.InitValueTypeBase(t, registry)
	return t
}

// Java: VoidType.getTypeClass()
func (t *VoidType) GetTypeClass() JSTypeClass {
	return JSTypeClass_VOID
}

// Java: VoidType.restrictByNotNullOrUndefined()
func (t *VoidType) RestrictByNotNullOrUndefined() JSType {
	return t.GetRegistry().GetNativeType(JSTypeNative_NO_TYPE)
}

// Java: VoidType.restrictByNotUndefined()
func (t *VoidType) RestrictByNotUndefined() JSType {
	return t.GetRegistry().GetNativeType(JSTypeNative_NO_TYPE)
}

// Java: VoidType.testForEquality(JSType)
func (t *VoidType) TestForEquality(that JSType) base.Tri {
	if base.Tri_UNKNOWN == t.ValueTypeBase.TestForEquality(that) {
		return base.Tri_UNKNOWN
	}
	if that.IsSubtypeOf(t) || that.IsSubtypeOf(t.GetNativeType(JSTypeNative_NULL_TYPE)) {
		return base.Tri_TRUE
	}
	return base.Tri_FALSE
}

// Java: VoidType.matchesNumberContext()
func (t *VoidType) MatchesNumberContext() bool {
	return false
}

// Java: VoidType.matchesObjectContext()
func (t *VoidType) MatchesObjectContext() bool {
	return false
}

// Java: VoidType.matchesStringContext()
func (t *VoidType) MatchesStringContext() bool {
	return true
}

// Java: VoidType.isVoidType()
func (t *VoidType) IsVoidType() bool {
	return true
}

// Java: VoidType.isVoidable()
func (t *VoidType) IsVoidable() bool {
	return true
}

// Java: VoidType.isExplicitlyVoidable()
func (t *VoidType) IsExplicitlyVoidable() bool {
	return true
}

// Java: VoidType.getDisplayName()
func (t *VoidType) GetDisplayName() *string {
	return &jsType_displayName_undefined
}

// Java: VoidType.getPossibleToBooleanOutcomes()
func (t *VoidType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return BooleanLiteralSet_FALSE
}

var jsType_displayName_undefined = "undefined"

// SymbolTypeBase holds SymbolType's methods; plain symbol types are
// *SymbolTypeBase values.
//
// Java: class SymbolType extends ValueType (subclassed by KnownSymbolType).
type SymbolTypeBase struct {
	ValueTypeBase
}

// NewSymbolType constructs the plain symbol type.
//
// Java: SymbolType(JSTypeRegistry)
func NewSymbolType(registry *JSTypeRegistry) *SymbolTypeBase {
	t := &SymbolTypeBase{}
	t.objFlags = jsKindScalarValue
	t.InitSymbolTypeBase(t, registry)
	return t
}

// InitSymbolTypeBase is the Go form of the constructor for subclasses.
//
// Java: SymbolType(JSTypeRegistry)
func (t *SymbolTypeBase) InitSymbolTypeBase(self JSType, registry *JSTypeRegistry) {
	t.InitValueTypeBase(self, registry)
}

// Java: SymbolType.getTypeClass()
func (t *SymbolTypeBase) GetTypeClass() JSTypeClass {
	return JSTypeClass_SYMBOL
}

// Java: SymbolType.testForEquality(JSType)
func (t *SymbolTypeBase) TestForEquality(that JSType) base.Tri {
	result := t.ValueTypeBase.TestForEquality(that)
	if result != jsType_TRI_NULL {
		return result
	}

	if that.CanCastTo(t.GetNativeType(JSTypeNative_SYMBOL_TYPE)) ||
		that.CanCastTo(t.GetNativeType(JSTypeNative_SYMBOL_OBJECT_TYPE)) {
		return base.Tri_UNKNOWN
	}
	return base.Tri_FALSE
}

// Java: SymbolType.isSymbolValueType()
func (t *SymbolTypeBase) IsSymbolValueType() bool {
	return true
}

// Java: SymbolType.matchesNumberContext()
func (t *SymbolTypeBase) MatchesNumberContext() bool {
	return false
}

// Java: SymbolType.matchesStringContext()
func (t *SymbolTypeBase) MatchesStringContext() bool {
	return false
}

// Java: SymbolType.matchesSymbolContext()
func (t *SymbolTypeBase) MatchesSymbolContext() bool {
	return true
}

// Java: SymbolType.matchesObjectContext()
func (t *SymbolTypeBase) MatchesObjectContext() bool {
	return true
}

// Java: SymbolType.getDisplayName()
func (t *SymbolTypeBase) GetDisplayName() *string {
	return &jsType_displayName_symbol
}

var jsType_displayName_symbol = "symbol"

// Java: SymbolType.autoboxesTo()
func (t *SymbolTypeBase) AutoboxesTo() JSType {
	return t.GetNativeType(JSTypeNative_SYMBOL_OBJECT_TYPE)
}

// Java: SymbolType.getPossibleToBooleanOutcomes()
func (t *SymbolTypeBase) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return BooleanLiteralSet_TRUE
}

// KnownSymbolType is the type of a well-known (or unique) symbol, such as
// Symbol.iterator.
type KnownSymbolType struct {
	SymbolTypeBase

	name string
}

// Java: KnownSymbolType.TYPE_CLASS
const KnownSymbolType_TYPE_CLASS = JSTypeClass_WELL_KNOWN_SYMBOL

// Java: KnownSymbolType(JSTypeRegistry, String)
func NewKnownSymbolType(registry *JSTypeRegistry, name string) *KnownSymbolType {
	t := &KnownSymbolType{}
	t.objFlags = jsKindKnownSymbol
	t.InitSymbolTypeBase(t, registry)
	t.name = name
	return t
}

// Java: KnownSymbolType.getTypeClass()
func (t *KnownSymbolType) GetTypeClass() JSTypeClass {
	return KnownSymbolType_TYPE_CLASS
}

// Java: KnownSymbolType.getDisplayName()
func (t *KnownSymbolType) GetDisplayName() *string {
	return &t.name
}

// Java: KnownSymbolType.isKnownSymbolValueType()
func (t *KnownSymbolType) IsKnownSymbolValueType() bool {
	return true
}

// Java: KnownSymbolType.toMaybeKnownSymbolType()
func (t *KnownSymbolType) ToMaybeKnownSymbolType() *KnownSymbolType {
	return t
}
