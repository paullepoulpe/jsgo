// Port of com/google/javascript/rhino/jstype/NoObjectType.java @ closure-compiler v20261021

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

import "github.com/paullepoulpe/jsgo/internal/jutil"

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

	// GetNoObjectTypeBase gives access to the base through the interface.
	//
	// Go-only: language-adapter — field access through the interface
	// (PORTING.md §4).
	GetNoObjectTypeBase() *NoObjectTypeBase
}

// NoObjectTypeBase holds NoObjectType's methods; the bottom Object type is
// a *NoObjectTypeBase value.
type NoObjectTypeBase struct {
	FunctionTypeBase
}

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

// GetNoObjectTypeBase gives access to the base through the interface.
//
// Go-only: language-adapter — field access through the interface (PORTING.md §4).
func (t *NoObjectTypeBase) GetNoObjectTypeBase() *NoObjectTypeBase { return t }

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
