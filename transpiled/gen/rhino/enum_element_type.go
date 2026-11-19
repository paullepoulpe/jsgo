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
// Port of com/google/javascript/rhino/jstype/EnumElementType.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.rhino.jstype.EnumElementType
type EnumElementType struct {
	ObjectTypeBase
	primitiveType       JSType     `gson:"primitiveType"`
	primitiveObjectType ObjectType `gson:"primitiveObjectType"`
	name                string     `gson:"name"`
	enumType            *EnumType  `gson:"enumType"`
}

func (this_ *EnumElementType) BindSelf_(s any) {
	_ = s
	this_.ObjectTypeBase.BindSelf_(s)
}
func (this_ *EnumElementType) IsJavaNil_() bool { return this_ == nil }
func (*EnumElementType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.EnumElementType", "EnumElementType"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.EnumElementType", "com.google.javascript.rhino.jstype.EnumElementType", func(v_ any) bool {
		return j2gort.InstanceOf[*EnumElementType](v_)
	})
}
func (*EnumElementType) Implements_rhino_ObjectType() {}
func (*EnumElementType) Implements_rhino_JSType()     {}

var static_EnumElementType_TYPE_CLASS JSTypeClass

func StaticGet_EnumElementType_TYPE_CLASS() JSTypeClass {
	clinit_EnumElementType()
	return static_EnumElementType_TYPE_CLASS
}
func StaticSet_EnumElementType_TYPE_CLASS(v JSTypeClass) {
	clinit_EnumElementType()
	static_EnumElementType_TYPE_CLASS = v
}

var clinitDone_EnumElementType uint32
var clinitState_EnumElementType j2gort.ClinitState

func clinit_EnumElementType() {
	if atomic.LoadUint32(&clinitDone_EnumElementType) != 0 {
		return
	}
	clinitSlow_EnumElementType()
}

func clinitSlow_EnumElementType() {
	if !j2gort.EnterClinit(&clinitState_EnumElementType, &clinitDone_EnumElementType) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_EnumElementType, &clinitDone_EnumElementType)
	static_EnumElementType_TYPE_CLASS = StaticGet_JSTypeClass_ENUM_ELEMENT()
}

func (this_ *EnumElementType) initFields_EnumElementType() {
	_ = this_
	this_.name = j2gort.NullString
}

// Java constructor: EnumElementType(com.google.javascript.rhino.jstype.JSTypeRegistry registry, com.google.javascript.rhino.jstype.JSType elementType, java.lang.String name, com.google.javascript.rhino.jstype.EnumType enumType)
func NewEnumElementType(registry *JSTypeRegistry, elementType JSType, name string, enumType *EnumType) *EnumElementType {
	this_ := &EnumElementType{}
	this_.BindSelf_(this_)
	InitEnumElementType(this_, registry, elementType, name, enumType)
	return this_
}

func InitEnumElementType(this_ *EnumElementType, registry *JSTypeRegistry, elementType JSType, name string, enumType *EnumType) {
	_ = registry
	_ = elementType
	_ = name
	_ = enumType
	InitObjectType_2_JSTypeRegistry_int32(&this_.ObjectTypeBase, registry, int32(256))
	this_.initFields_EnumElementType()
	this_.GetJSTypeBase().jstypeFlags = int8((int32(this_.GetJSTypeBase().jstypeFlags) | int32(int8(32))))
	this_.primitiveType = elementType
	this_.primitiveObjectType = elementType.GetJSTypeBase().ToObjectType()
	this_.name = name
	this_.enumType = enumType
	if (elementType == nil || elementType.IsJavaNil_()) || (elementType.GetJSTypeBase().GetWellFormedTemplatizedKind() == int32(0)) {
		this_.GetJSTypeBase().protoFlags = int8((int32(this_.GetJSTypeBase().protoFlags) | int32(int8(-128))))
	}
	registry.GetResolver().ResolveIfClosed(ToIface_JSType(this_), StaticGet_EnumElementType_TYPE_CLASS())
}

// Java method: EnumElementType.getTypeClass()
func (this_ *EnumElementType) GetTypeClass() JSTypeClass {
	return StaticGet_EnumElementType_TYPE_CLASS()
}

// Java method: EnumElementType.getEnumType()
func (this_ *EnumElementType) GetEnumType() *EnumType {
	return this_.enumType
}

// Java method: EnumElementType.getPropertyKind(com.google.javascript.rhino.jstype.Property.Key propertyName, boolean autobox)
func (this_ *EnumElementType) GetPropertyKind_2_Property_Key_bool(propertyName Property_Key, autobox bool) JSType_HasPropertyKind {
	_ = propertyName
	_ = autobox
	return this_.primitiveType.GetPropertyKind_2_Property_Key_bool(propertyName, autobox)
}

// Java method: EnumElementType.getPropertyKind(java.lang.String propertyName, boolean autobox)
func (this_ *EnumElementType) GetPropertyKind_2_string_bool(propertyName string, autobox bool) JSType_HasPropertyKind {
	_ = propertyName
	_ = autobox
	return this_.primitiveType.GetPropertyKind_2_string_bool(propertyName, autobox)
}

// Java method: EnumElementType.getPropertyMap()
func (this_ *EnumElementType) GetPropertyMap() *PropertyMap {
	if this_.primitiveObjectType == nil || this_.primitiveObjectType.IsJavaNil_() {
		return PropertyMap_ImmutableEmptyMap()
	} else {
		return this_.primitiveObjectType.GetPropertyMap()
	}
}

// Java method: EnumElementType.toMaybeEnumElementType()
func (this_ *EnumElementType) ToMaybeEnumElementType() *EnumElementType {
	return this_
}

// Java method: EnumElementType.matchesNumberContext()
func (this_ *EnumElementType) MatchesNumberContext() bool {
	return this_.primitiveType.MatchesNumberContext()
}

// Java method: EnumElementType.matchesStringContext()
func (this_ *EnumElementType) MatchesStringContext() bool {
	return this_.primitiveType.MatchesStringContext()
}

// Java method: EnumElementType.matchesObjectContext()
func (this_ *EnumElementType) MatchesObjectContext() bool {
	return this_.primitiveType.MatchesObjectContext()
}

// Java method: EnumElementType.matchesSymbolContext()
func (this_ *EnumElementType) MatchesSymbolContext() bool {
	return this_.primitiveType.MatchesSymbolContext()
}

// Java method: EnumElementType.canBeCalled()
func (this_ *EnumElementType) CanBeCalled() bool {
	return this_.primitiveType.CanBeCalled()
}

// Java method: EnumElementType.isObject()
func (this_ *EnumElementType) IsObject() bool {
	return this_.primitiveType.IsObject()
}

// Java method: EnumElementType.testForEquality(com.google.javascript.rhino.jstype.JSType that)
func (this_ *EnumElementType) TestForEquality(that JSType) pkg_base.Tri {
	_ = that
	return this_.primitiveType.TestForEquality(that)
}

// Java method: EnumElementType.isNullable()
func (this_ *EnumElementType) IsNullable() bool {
	return this_.primitiveType.IsNullable()
}

// Java method: EnumElementType.isVoidable()
func (this_ *EnumElementType) IsVoidable() bool {
	return this_.primitiveType.IsVoidable()
}

// Java method: EnumElementType.isNominalType()
func (this_ *EnumElementType) IsNominalType() bool {
	return this_.HasReferenceName()
}

// Java method: EnumElementType.recursionUnsafeHashCode()
func (this_ *EnumElementType) RecursionUnsafeHashCode() int32 {
	if !(this_.HasReferenceName()) {
		return int32(2)
	}
	return NamedType_NominalHashCode(this_)
}

// Java method: EnumElementType.appendTo(com.google.javascript.rhino.jstype.TypeStringBuilder sb)
func (this_ *EnumElementType) AppendTo(sb *TypeStringBuilder) {
	_ = sb
	if sb.IsForAnnotations() {
		sb.Append_1_JSType(this_.primitiveType)
	} else {
		sb.Append_1_string(this_.GetReferenceName()).Append_1_string("<").Append_1_JSType(this_.primitiveType).Append_1_string(">")
	}
}

// Java method: EnumElementType.getReferenceName()
func (this_ *EnumElementType) GetReferenceName() string {
	return this_.name
}

// Java method: EnumElementType.visit(com.google.javascript.rhino.jstype.Visitor<T> visitor)
func (this_ *EnumElementType) Visit_1(visitor Visitor) any {
	_ = visitor
	return visitor.CaseEnumElementType(this_)
}

// Java method: EnumElementType.visit(com.google.javascript.rhino.jstype.RelationshipVisitor<T> visitor, com.google.javascript.rhino.jstype.JSType that)
func (this_ *EnumElementType) Visit_2(visitor RelationshipVisitor, that JSType) any {
	_ = visitor
	_ = that
	return visitor.CaseEnumElementType(this_, that)
}

// Java method: EnumElementType.defineProperty(java.lang.String propertyName, com.google.javascript.rhino.jstype.JSType type, boolean inferred, com.google.javascript.rhino.Node propertyNode)
func (this_ *EnumElementType) DefineProperty_4_string_JSType_bool_Node(propertyName string, type_ JSType, inferred bool, propertyNode Node) bool {
	_ = propertyName
	_ = type_
	_ = inferred
	_ = propertyNode
	return true
}

// Java method: EnumElementType.defineProperty(com.google.javascript.rhino.jstype.Property.Key propertyName, com.google.javascript.rhino.jstype.JSType type, boolean inferred, com.google.javascript.rhino.Node propertyNode)
func (this_ *EnumElementType) DefineProperty_4_Property_Key_JSType_bool_Node(propertyName Property_Key, type_ JSType, inferred bool, propertyNode Node) bool {
	_ = propertyName
	_ = type_
	_ = inferred
	_ = propertyNode
	return true
}

// Java method: EnumElementType.getImplicitPrototype()
func (this_ *EnumElementType) GetImplicitPrototype() ObjectType {
	return nil
}

// Java method: EnumElementType.findPropertyTypeWithoutConsideringTemplateTypes(com.google.javascript.rhino.jstype.Property.Key propertyName)
func (this_ *EnumElementType) FindPropertyTypeWithoutConsideringTemplateTypes_1_Property_Key(propertyName Property_Key) JSType {
	_ = propertyName
	return this_.primitiveType.GetJSTypeBase().FindPropertyType_1_Property_Key(propertyName)
}

// Java method: EnumElementType.findPropertyTypeWithoutConsideringTemplateTypes(java.lang.String propertyName)
func (this_ *EnumElementType) FindPropertyTypeWithoutConsideringTemplateTypes_1_string(propertyName string) JSType {
	_ = propertyName
	return this_.primitiveType.GetJSTypeBase().FindPropertyType_1_string(propertyName)
}

// Java method: EnumElementType.getConstructor()
func (this_ *EnumElementType) GetConstructor() *FunctionType {
	if this_.primitiveObjectType == nil || this_.primitiveObjectType.IsJavaNil_() {
		return nil
	} else {
		return this_.primitiveObjectType.GetConstructor()
	}
}

// Java method: EnumElementType.autoboxesTo()
func (this_ *EnumElementType) AutoboxesTo() JSType {
	return this_.primitiveType.AutoboxesTo()
}

// Java method: EnumElementType.getPrimitiveType()
func (this_ *EnumElementType) GetPrimitiveType() JSType {
	return this_.primitiveType
}

// Java static method: EnumElementType.getGreatestSubtype(com.google.javascript.rhino.jstype.EnumElementType element, com.google.javascript.rhino.jstype.JSType that)
func EnumElementType_GetGreatestSubtype(element *EnumElementType, that JSType) JSType {
	_ = element
	_ = that
	var meetPrimitive JSType = element.primitiveType.GetJSTypeBase().GetGreatestSubtype(that)
	_ = meetPrimitive
	if meetPrimitive.GetJSTypeBase().IsEmptyType() {
		return nil
	} else {
		return NewEnumElementType(element.GetJSTypeBase().registry, meetPrimitive, element.name, element.GetEnumType())
	}
}

// Java method: EnumElementType.resolveInternal(com.google.javascript.rhino.ErrorReporter reporter)
func (this_ *EnumElementType) ResolveInternal(reporter ErrorReporter) JSType {
	_ = reporter
	this_.primitiveType = this_.primitiveType.GetJSTypeBase().Resolve(reporter)
	this_.primitiveObjectType = ObjectType_Cast(this_.primitiveType)
	if (this_.primitiveType == nil || this_.primitiveType.IsJavaNil_()) || (this_.primitiveType.GetJSTypeBase().GetWellFormedTemplatizedKind() == int32(0)) {
		this_.GetJSTypeBase().protoFlags = int8((int32(this_.GetJSTypeBase().protoFlags) | int32(int8(-128))))
	} else {
		this_.GetJSTypeBase().protoFlags = int8((int32(this_.GetJSTypeBase().protoFlags) & ^(int32(int8(-128)))))
	}
	return this_
}
