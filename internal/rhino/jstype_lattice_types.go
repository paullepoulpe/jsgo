// Port of com/google/javascript/rhino/jstype/AllType.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/jstype/UnknownType.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/jstype/NoObjectType.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/jstype/NoType.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/jstype/NoResolvedType.java @ closure-compiler v20261021

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

// AllType is the all type, representing all values.
type AllType struct {
	JSTypeBase
}

// Java: AllType(JSTypeRegistry)
func NewAllType(registry *JSTypeRegistry) *AllType {
	t := &AllType{}
	t.objFlags = jsKindAll
	t.InitJSTypeBase(t, registry)
	t.EagerlyResolveToSelf()
	return t
}

// Java: AllType.getTypeClass()
func (t *AllType) GetTypeClass() JSTypeClass {
	return JSTypeClass_ALL
}

// Java: AllType.isAllType()
func (t *AllType) IsAllType() bool {
	return true
}

// Java: AllType.matchesStringContext()
func (t *AllType) MatchesStringContext() bool {
	// Be lenient.
	return true
}

// Java: AllType.matchesObjectContext()
func (t *AllType) MatchesObjectContext() bool {
	// Be lenient.
	return true
}

// Java: AllType.testForEquality(JSType)
func (t *AllType) TestForEquality(that JSType) base.Tri {
	return base.Tri_UNKNOWN
}

// Java: AllType.appendTo(TypeStringBuilder)
func (t *AllType) AppendTo(sb *TypeStringBuilder) {
	sb.Append("*")
}

// Java: AllType.getDisplayName()
func (t *AllType) GetDisplayName() *string {
	return &allType_displayName
}

var allType_displayName = "<Any Type>"

// Java: AllType.hasDisplayName()
func (t *AllType) HasDisplayName() bool {
	return true
}

// Java: AllType.getPossibleToBooleanOutcomes()
func (t *AllType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return BooleanLiteralSet_BOTH
}

// Java: AllType.resolveInternal(ErrorReporter)
func (t *AllType) ResolveInternal(reporter ErrorReporter) JSType {
	panic(jutil.NewAssertionErrorNoMessage())
}

// Java: AllType.recursionUnsafeHashCode()
func (t *AllType) RecursionUnsafeHashCode() int32 {
	return jutil.IdentityHashCode(t)
}

// Java: AllType.isNullable()
func (t *AllType) IsNullable() bool {
	return true
}

// Java: AllType.isVoidable()
func (t *AllType) IsVoidable() bool {
	return true
}

// UnknownType is the unknown type.
type UnknownType struct {
	ObjectTypeBase

	// See the explanation of checked unknown types in JSTypeNative.
	isChecked bool
}

// Java: UnknownType(JSTypeRegistry, boolean)
func NewUnknownType(registry *JSTypeRegistry, isChecked bool) *UnknownType {
	t := &UnknownType{}
	t.objFlags = objKindUnknown
	t.InitObjectTypeBase(t, registry)
	t.EagerlyResolveToSelf()
	t.isChecked = isChecked
	return t
}

// Java: UnknownType.getTypeClass()
func (t *UnknownType) GetTypeClass() JSTypeClass {
	return JSTypeClass_UNKNOWN
}

// Java: UnknownType.isUnknownType()
func (t *UnknownType) IsUnknownType() bool {
	return true
}

// Java: UnknownType.isCheckedUnknownType()
func (t *UnknownType) IsCheckedUnknownType() bool {
	return t.isChecked
}

// Java: UnknownType.canBeCalled()
func (t *UnknownType) CanBeCalled() bool {
	return true
}

// Java: UnknownType.matchesNumberContext()
func (t *UnknownType) MatchesNumberContext() bool {
	return true
}

// Java: UnknownType.matchesObjectContext()
func (t *UnknownType) MatchesObjectContext() bool {
	return true
}

// Java: UnknownType.matchesStringContext()
func (t *UnknownType) MatchesStringContext() bool {
	return true
}

// Java: UnknownType.matchesSymbolContext()
func (t *UnknownType) MatchesSymbolContext() bool {
	return true
}

// Java: UnknownType.testForEquality(JSType)
func (t *UnknownType) TestForEquality(that JSType) base.Tri {
	return base.Tri_UNKNOWN
}

// Java: UnknownType.isNullable()
func (t *UnknownType) IsNullable() bool {
	return true
}

// Java: UnknownType.isVoidable()
func (t *UnknownType) IsVoidable() bool {
	return true
}

// Java: UnknownType.appendTo(TypeStringBuilder)
func (t *UnknownType) AppendTo(sb *TypeStringBuilder) {
	sb.Append(*t.GetReferenceName())
}

// Java: UnknownType.defineProperty(Property.Key, JSType, boolean, Node)
func (t *UnknownType) DefinePropertyByKey(propertyName Property_Key, type_ JSType, inferred bool, propertyNode Node) bool {
	// nothing to define
	return true
}

// Java: UnknownType.getImplicitPrototype()
func (t *UnknownType) GetImplicitPrototype() ObjectType {
	return nil
}

// Java: UnknownType.getConstructor()
func (t *UnknownType) GetConstructor() FunctionType {
	return nil
}

// Java: UnknownType.getReferenceName()
func (t *UnknownType) GetReferenceName() *string {
	if t.isChecked {
		return &unknownType_checkedName
	}
	return &unknownType_name
}

var (
	unknownType_checkedName = "??"
	unknownType_name        = "?"
	unknownType_displayName = "Unknown"
)

// Java: UnknownType.getDisplayName()
func (t *UnknownType) GetDisplayName() *string {
	return &unknownType_displayName
}

// Java: UnknownType.hasDisplayName()
func (t *UnknownType) HasDisplayName() bool {
	return true
}

// Java: UnknownType.getPossibleToBooleanOutcomes()
func (t *UnknownType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return BooleanLiteralSet_BOTH
}

// Java: UnknownType.resolveInternal(ErrorReporter)
func (t *UnknownType) ResolveInternal(reporter ErrorReporter) JSType {
	panic(jutil.NewAssertionErrorNoMessage())
}

// Java: UnknownType.recursionUnsafeHashCode()
func (t *UnknownType) RecursionUnsafeHashCode() int32 {
	return jutil.IdentityHashCode(t)
}

// NoObjectType is the bottom Object type, representing the subclass of all
// objects.
//
// Although JavaScript programmers can't explicitly denote the bottom Object
// type, it comes up in static analysis. For example, if we have:
//
//	var x = function() {};
//	if (x instanceof Array) {
//	  f(x);
//	}
//
// We need to be able to assign x a type within the f(x) call. It has no
// possible type, but x would not be legal if f expected a string. So we
// assign it the NoObjectType.
//
// Java: class NoObjectType extends FunctionType (subclassed by NoType).
type NoObjectType interface {
	FunctionType
	isNoObjectType()
}

// NoObjectTypeBase holds NoObjectType's methods; the bottom Object type is
// a *NoObjectTypeBase value.
type NoObjectTypeBase struct {
	FunctionTypeBase
}

func (*NoObjectTypeBase) isNoObjectType() {}

// Java: NoObjectType(JSTypeRegistry)
func NewNoObjectType(registry *JSTypeRegistry) *NoObjectTypeBase {
	t := &NoObjectTypeBase{}
	t.objFlags = objKindNoObject
	t.InitNoObjectTypeBase(t, registry)
	return t
}

// InitNoObjectTypeBase is the Go form of the constructor for subclasses.
//
// Java: NoObjectType(JSTypeRegistry)
func (t *NoObjectTypeBase) InitNoObjectTypeBase(self NoObjectType, registry *JSTypeRegistry) {
	b := FunctionType_Builder_(registry).
		WithKind(FunctionType_Kind_NONE).
		WithReturnsOwnInstanceType().
		ForNativeType()
	t.InitFunctionTypeBase(self, b)
	registry.releaseFunctionBuilder(b)
	t.EagerlyResolveToSelf()
}

// Java: NoObjectType.getTypeClass()
func (t *NoObjectTypeBase) GetTypeClass() JSTypeClass {
	return JSTypeClass_NO_OBJECT
}

// Java: JSType.isFunctionType()
func (t *NoObjectTypeBase) IsFunctionType() bool {
	return false
}

// Java: NoObjectType.toMaybeFunctionType()
func (t *NoObjectTypeBase) ToMaybeFunctionType() FunctionType {
	return nil
}

// Java: NoObjectType.isNoObjectType()
func (t *NoObjectTypeBase) IsNoObjectType() bool {
	return true
}

// Java: NoObjectType.getImplicitPrototype()
func (t *NoObjectTypeBase) GetImplicitPrototype() ObjectType {
	return nil
}

// Java: ObjectType.hasReferenceName()
func (t *NoObjectTypeBase) HasReferenceName() bool {
	return false
}

// Java: NoObjectType.getReferenceName()
func (t *NoObjectTypeBase) GetReferenceName() *string {
	return nil
}

// Java: NoObjectType.matchesNumberContext()
func (t *NoObjectTypeBase) MatchesNumberContext() bool {
	return true
}

// Java: NoObjectType.matchesObjectContext()
func (t *NoObjectTypeBase) MatchesObjectContext() bool {
	return true
}

// Java: NoObjectType.matchesStringContext()
func (t *NoObjectTypeBase) MatchesStringContext() bool {
	return true
}

// Java: NoObjectType.matchesSymbolContext()
func (t *NoObjectTypeBase) MatchesSymbolContext() bool {
	return true
}

// Java: NoObjectType.recursionUnsafeHashCode()
func (t *NoObjectTypeBase) RecursionUnsafeHashCode() int32 {
	return jutil.IdentityHashCode(t.ObjectTypeBase.self())
}

// Java: NoObjectType.defineProperty(Property.Key, JSType, boolean, Node)
func (t *NoObjectTypeBase) DefinePropertyByKey(propertyName Property_Key, type_ JSType, inferred bool, propertyNode Node) bool {
	// nothing, all properties are defined
	return true
}

// Java: NoObjectType.setPropertyJSDocInfo(Property.Key, JSDocInfo)
func (t *NoObjectTypeBase) SetPropertyJSDocInfoByKey(propertyName Property_Key, info *JSDocInfo) {
	// Do nothing, specific properties do not have JSDocInfo.
}

// Java: NoObjectType.appendTo(TypeStringBuilder)
func (t *NoObjectTypeBase) AppendTo(sb *TypeStringBuilder) {
	if sb.IsForAnnotations() {
		sb.Append("?")
	} else {
		sb.Append("NoObject")
	}
}

// Java: NoObjectType.getConstructor()
func (t *NoObjectTypeBase) GetConstructor() FunctionType {
	return nil
}

// Java: NoObjectType.resolveInternal(ErrorReporter)
func (t *NoObjectTypeBase) ResolveInternal(reporter ErrorReporter) JSType {
	panic(jutil.NewAssertionErrorNoMessage())
}

// NoType is the bottom type, representing the subclass of any value or
// object.
//
// Although JavaScript programmers can't explicitly denote the bottom type,
// it comes up in static analysis. For example, if we have:
//
//	var x = null;
//	if (x) {
//	  f(x);
//	}
//
// We need to be able to assign x a type within the f(x) call. Since it has
// no possible type, we assign x the NoType, so that f(x) is legal no matter
// what the type of f's first argument is.
type NoType struct {
	NoObjectTypeBase
}

// Java: NoType(JSTypeRegistry)
func NewNoType(registry *JSTypeRegistry) *NoType {
	t := &NoType{}
	t.objFlags = objKindNo
	t.InitNoObjectTypeBase(t, registry)
	return t
}

// Java: NoType.getTypeClass()
func (t *NoType) GetTypeClass() JSTypeClass {
	return JSTypeClass_NO
}

// Java: NoType.isNoObjectType()
func (t *NoType) IsNoObjectType() bool {
	return false
}

// Java: NoType.isNoType()
func (t *NoType) IsNoType() bool {
	return true
}

// Java: NoType.isNullable()
func (t *NoType) IsNullable() bool {
	return true
}

// Java: NoType.isVoidable()
func (t *NoType) IsVoidable() bool {
	return true
}

// Java: NoType.getPossibleToBooleanOutcomes()
func (t *NoType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return BooleanLiteralSet_BOTH
}

// Java: NoType.appendTo(TypeStringBuilder)
func (t *NoType) AppendTo(sb *TypeStringBuilder) {
	if sb.IsForAnnotations() {
		sb.Append("?")
	} else {
		sb.Append("None")
	}
}

// NoResolvedType is an unresolved type that was forward declared. So we
// know it exists, but that it wasn't pulled into this binary.
//
// In most cases, it behaves like a sibling of the top object type.
//
// Why have this class at all, and not just use the top object type or `*`?
//
// Clutz relies on this class to emit .d.ts typings for JS files without
// having to see all their dependencies: Clutz expects unresolved types to be
// available when it runs, so that it can map them back to the corresponding
// goog.required import.
//
// Clutz also sets a special JSCompiler flag so that all unresolved types are
// treated as forward declared. (In most non-Clutz cases, this class is only
// rarely used, and only if there's actually an explicit
// `goog.forwardDeclare` call.)
//
// To address this, we would need to significantly refactor Clutz.
type NoResolvedType struct {
	ObjectTypeBase

	// The name originally used to reference this type.
	referenceName string

	// Any template arguments to this type, or null if none. This field is not
	// used for JSCompiler's type checking; it is only needed by Clutz.
	templateTypes []JSType
}

// Java: NoResolvedType(JSTypeRegistry, String, ImmutableList<JSType>)
func NewNoResolvedType(registry *JSTypeRegistry, referenceName string, templateTypes []JSType) *NoResolvedType {
	t := &NoResolvedType{}
	t.objFlags = objKindNoResolved
	t.InitObjectTypeBase(t, registry)
	t.referenceName = referenceName
	t.templateTypes = templateTypes
	t.EagerlyResolveToSelf()
	return t
}

// Java: NoResolvedType.getTypeClass()
func (t *NoResolvedType) GetTypeClass() JSTypeClass {
	return JSTypeClass_NO_RESOLVED
}

// Java: NoResolvedType.getReferenceName()
func (t *NoResolvedType) GetReferenceName() *string {
	return &t.referenceName
}

// Java: NoResolvedType.getTemplateTypes()
func (t *NoResolvedType) GetTemplateTypes() []JSType {
	return t.templateTypes
}

// Java: NoResolvedType.isNoResolvedType()
func (t *NoResolvedType) IsNoResolvedType() bool {
	return true
}

// Java: NoResolvedType.getPossibleToBooleanOutcomes()
func (t *NoResolvedType) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return BooleanLiteralSet_BOTH
}

// Java: NoResolvedType.testForEquality(JSType)
func (t *NoResolvedType) TestForEquality(that JSType) base.Tri {
	return base.Tri_UNKNOWN
}

// Java: NoResolvedType.resolveInternal(ErrorReporter)
func (t *NoResolvedType) ResolveInternal(reporter ErrorReporter) JSType {
	panic(jutil.NewAssertionErrorNoMessage())
}

// Java: NoResolvedType.getConstructor()
func (t *NoResolvedType) GetConstructor() FunctionType {
	return nil
}

// Java: NoResolvedType.getImplicitPrototype()
func (t *NoResolvedType) GetImplicitPrototype() ObjectType {
	return nil
}

// Java: NoResolvedType.recursionUnsafeHashCode()
func (t *NoResolvedType) RecursionUnsafeHashCode() int32 {
	// NoResolvedType instances are unique within a JSTypeRegistry.
	return jutil.IdentityHashCode(t)
}

// Java: NoResolvedType.defineProperty(Property.Key, JSType, boolean, Node)
func (t *NoResolvedType) DefinePropertyByKey(propertyName Property_Key, type_ JSType, inferred bool, propertyNode Node) bool {
	// nothing, all properties are defined
	return true
}

// Java: NoResolvedType.appendTo(TypeStringBuilder)
func (t *NoResolvedType) AppendTo(sb *TypeStringBuilder) {
	if sb.IsForAnnotations() {
		sb.Append("?")
	} else {
		sb.Append("NoResolvedType<" + t.referenceName + ">")
	}
}
