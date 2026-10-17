// Port of com/google/javascript/rhino/jstype/EnumElementType.java @ closure-compiler v20260415

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
)

// EnumElementType is the type of individual elements of an enum type (see
// EnumType).
type EnumElementType struct {
	ObjectTypeBase

	// The primitive type this enum element type wraps. For instance, in the
	// following code defining the LOCAL_CODES enum
	//
	//	var LOCAL_CODES = {A: 3, B: 9, C: 8}
	//
	// the primitive type of the the constants is number.
	primitiveType JSType

	// The primitive type, if it is an object.
	primitiveObjectType ObjectType

	name *string

	enumType *EnumType
}

// Java: EnumElementType.TYPE_CLASS
const EnumElementType_TYPE_CLASS = JSTypeClass_ENUM_ELEMENT

// Java: EnumElementType(JSTypeRegistry, JSType, String, EnumType)
func NewEnumElementType(registry *JSTypeRegistry, elementType JSType, name *string, enumType *EnumType) *EnumElementType {
	t := &EnumElementType{}
	t.InitObjectTypeBase(t, registry)
	t.primitiveType = elementType
	t.primitiveObjectType = elementType.ToObjectType()
	t.name = name
	t.enumType = enumType

	registry.GetResolver().ResolveIfClosed(t, EnumElementType_TYPE_CLASS)
	return t
}

// Java: EnumElementType.getTypeClass()
func (t *EnumElementType) GetTypeClass() JSTypeClass {
	return EnumElementType_TYPE_CLASS
}

// Java: EnumElementType.getEnumType()
func (t *EnumElementType) GetEnumType() *EnumType {
	return t.enumType
}

// Java: EnumElementType.getPropertyKind(Property.Key, boolean)
func (t *EnumElementType) GetPropertyKindByKeyWithAutobox(propertyName Property_Key, autobox bool) JSType_HasPropertyKind {
	return t.primitiveType.GetPropertyKindByKeyWithAutobox(propertyName, autobox)
}

// Java: EnumElementType.getPropertyMap()
func (t *EnumElementType) GetPropertyMap() *PropertyMap {
	if t.primitiveObjectType == nil {
		return PropertyMap_ImmutableEmptyMap()
	}
	return t.primitiveObjectType.GetPropertyMap()
}

// Java: EnumElementType.toMaybeEnumElementType()
func (t *EnumElementType) ToMaybeEnumElementType() *EnumElementType {
	return t
}

// Java: EnumElementType.matchesNumberContext()
func (t *EnumElementType) MatchesNumberContext() bool {
	return t.primitiveType.MatchesNumberContext()
}

// Java: EnumElementType.matchesStringContext()
func (t *EnumElementType) MatchesStringContext() bool {
	return t.primitiveType.MatchesStringContext()
}

// Java: EnumElementType.matchesObjectContext()
func (t *EnumElementType) MatchesObjectContext() bool {
	return t.primitiveType.MatchesObjectContext()
}

// Java: EnumElementType.matchesSymbolContext()
func (t *EnumElementType) MatchesSymbolContext() bool {
	return t.primitiveType.MatchesSymbolContext()
}

// Java: EnumElementType.canBeCalled()
func (t *EnumElementType) CanBeCalled() bool {
	return t.primitiveType.CanBeCalled()
}

// Java: EnumElementType.isObject()
func (t *EnumElementType) IsObject() bool {
	return t.primitiveType.IsObject()
}

// Java: EnumElementType.testForEquality(JSType)
func (t *EnumElementType) TestForEquality(that JSType) base.Tri {
	return t.primitiveType.TestForEquality(that)
}

// Java: EnumElementType.isNullable()
func (t *EnumElementType) IsNullable() bool {
	return t.primitiveType.IsNullable()
}

// Java: EnumElementType.isVoidable()
func (t *EnumElementType) IsVoidable() bool {
	return t.primitiveType.IsVoidable()
}

// Java: EnumElementType.isNominalType()
func (t *EnumElementType) IsNominalType() bool {
	return t.HasReferenceName()
}

// Java: EnumElementType.recursionUnsafeHashCode()
func (t *EnumElementType) RecursionUnsafeHashCode() int32 {
	if !t.HasReferenceName() {
		// TODO(nickreid): Apparently this can happen if the l-value the enum is
		// assinged to is not a qname. Fortunatly, this whole thing should become
		// redundant once equality cannot be checked before resolution.
		return 2
	}
	return NamedType_NominalHashCode(t)
}

// Java: EnumElementType.appendTo(TypeStringBuilder)
func (t *EnumElementType) AppendTo(sb *TypeStringBuilder) {
	if sb.IsForAnnotations() {
		// TODO(dimvar): this should use getReferenceName() instead of
		// this.primitiveType
		sb.AppendType(t.primitiveType)
	} else {
		sb.Append(jsTypeNullableToString(t.GetReferenceName())).Append("<").AppendType(t.primitiveType).Append(">")
	}
}

// jsTypeNullableToString is String.valueOf for a nullable String, as
// StringBuilder.append(String) renders it ("null" for null).
func jsTypeNullableToString(s *string) string {
	if s == nil {
		return "null"
	}
	return *s
}

// Java: EnumElementType.getReferenceName()
func (t *EnumElementType) GetReferenceName() *string {
	return t.name
}

// Java: EnumElementType.visit(Visitor<T>)
func (t *EnumElementType) Visit(visitor Visitor) any {
	return visitor.CaseEnumElementType(t)
}

// Java: EnumElementType.visit(RelationshipVisitor<T>, JSType)
func (t *EnumElementType) VisitRelationship(visitor RelationshipVisitor, that JSType) any {
	return visitor.CaseEnumElementType(t, that)
}

// Java: EnumElementType.defineProperty(Property.Key, JSType, boolean, Node)
func (t *EnumElementType) DefinePropertyByKey(propertyName Property_Key, type_ JSType, inferred bool, propertyNode *Node) bool {
	// nothing
	return true
}

// Java: EnumElementType.getImplicitPrototype()
func (t *EnumElementType) GetImplicitPrototype() ObjectType {
	return nil
}

// Java: EnumElementType.findPropertyTypeWithoutConsideringTemplateTypes(Property.Key)
func (t *EnumElementType) FindPropertyTypeWithoutConsideringTemplateTypes(propertyName Property_Key) JSType {
	return t.primitiveType.FindPropertyTypeByKey(propertyName)
}

// Java: EnumElementType.getConstructor()
func (t *EnumElementType) GetConstructor() FunctionType {
	// TODO(b/147236174): This should always return null.
	if t.primitiveObjectType == nil {
		return nil
	}
	return t.primitiveObjectType.GetConstructor()
}

// Java: EnumElementType.autoboxesTo()
func (t *EnumElementType) AutoboxesTo() JSType {
	return t.primitiveType.AutoboxesTo()
}

// GetPrimitiveType gets the primitive type of this enum element.
//
// Java: EnumElementType.getPrimitiveType()
func (t *EnumElementType) GetPrimitiveType() JSType {
	return t.primitiveType
}

// EnumElementType_GetGreatestSubtype returns the infimum of a enum element
// type and another type, or null if the infimum is empty.
//
// This can be a little bit weird. For example, suppose you have an enum of
// {(string|number)}, and you want the greatest subtype of the enum and a
// {number}.
//
// The infimum is non-empty. But at the same time, we don't really have a
// name for this infimum. It's equivalent to "elements of this enum that are
// numbers".
//
// The best we can do is make up a new type. This is similar to what we do
// in UnionType#meet, which kind-of-sort-of makes sense, because an
// EnumElementType is a union of instances of a type.
//
// Java: EnumElementType.getGreatestSubtype(EnumElementType, JSType)
func EnumElementType_GetGreatestSubtype(element *EnumElementType, that JSType) JSType {
	// This method is implemented as a static because we don't want
	// polymorphism. Ideally all the `greatestSubtype` code would be in one
	// place. Until then, using static calls minimizes confusion.

	meetPrimitive := element.primitiveType.GetGreatestSubtype(that)
	if meetPrimitive.IsEmptyType() {
		return nil
	} else {
		return NewEnumElementType(element.registry, meetPrimitive, element.name, element.GetEnumType())
	}
}

// Java: EnumElementType.resolveInternal(ErrorReporter)
func (t *EnumElementType) ResolveInternal(reporter ErrorReporter) JSType {
	t.primitiveType = t.primitiveType.Resolve(reporter)
	t.primitiveObjectType = ObjectType_Cast(t.primitiveType)
	return t
}
