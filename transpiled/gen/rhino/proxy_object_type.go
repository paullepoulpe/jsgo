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
// Port of com/google/javascript/rhino/jstype/ProxyObjectType.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java polymorphic class interface: com.google.javascript.rhino.jstype.ProxyObjectType
type ProxyObjectType interface {
	Implements_rhino_ProxyObjectType()
	IsJavaNil_() bool
	GetProxyObjectTypeBase() *ProxyObjectTypeBase
	ObjectType
	GetReferencedTypeInternal() JSType
	GetReferencedObjTypeInternal() ObjectType
	SetReferencedType(JSType)
	VisitReferenceType(Visitor) any
}

func ToIface_ProxyObjectType[T interface {
	ProxyObjectType
	comparable
}](v T) ProxyObjectType {
	var z T
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.rhino.jstype.ProxyObjectType
type ProxyObjectTypeBase struct {
	ObjectTypeBase
	referencedType    JSType     `gson:"referencedType"`
	referencedObjType ObjectType `gson:"referencedObjType"`
}

func (this_ *ProxyObjectTypeBase) BindSelf_(s any) {
	_ = s
	this_.ObjectTypeBase.BindSelf_(s)
}
func (this_ *ProxyObjectTypeBase) IsJavaNil_() bool { return this_ == nil }
func (*ProxyObjectTypeBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.ProxyObjectType", "ProxyObjectType"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.ProxyObjectType", "com.google.javascript.rhino.jstype.ProxyObjectType", func(v_ any) bool {
		return j2gort.InstanceOf[ProxyObjectType](v_)
	})
}
func (b *ProxyObjectTypeBase) GetProxyObjectTypeBase() *ProxyObjectTypeBase { return b }
func (this_ *ProxyObjectTypeBase) self_ProxyObjectType_() ProxyObjectType {
	if s := this_.JSTypeBase.self_; s != nil {
		return s.(ProxyObjectType)
	}
	return nil
}
func (*ProxyObjectTypeBase) Implements_rhino_ProxyObjectType() {}
func (*ProxyObjectTypeBase) Implements_rhino_ObjectType()      {}
func (*ProxyObjectTypeBase) Implements_rhino_JSType()          {}

var static_ProxyObjectType_TYPE_CLASS JSTypeClass

func StaticGet_ProxyObjectType_TYPE_CLASS() JSTypeClass {
	clinit_ProxyObjectType()
	return static_ProxyObjectType_TYPE_CLASS
}
func StaticSet_ProxyObjectType_TYPE_CLASS(v JSTypeClass) {
	clinit_ProxyObjectType()
	static_ProxyObjectType_TYPE_CLASS = v
}

var clinitDone_ProxyObjectType uint32
var clinitState_ProxyObjectType j2gort.ClinitState

func clinit_ProxyObjectType() {
	if atomic.LoadUint32(&clinitDone_ProxyObjectType) != 0 {
		return
	}
	clinitSlow_ProxyObjectType()
}

func clinitSlow_ProxyObjectType() {
	if !j2gort.EnterClinit(&clinitState_ProxyObjectType, &clinitDone_ProxyObjectType) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ProxyObjectType, &clinitDone_ProxyObjectType)
	static_ProxyObjectType_TYPE_CLASS = StaticGet_JSTypeClass_PROXY_OBJECT()
}

func (this_ *ProxyObjectTypeBase) initFields_ProxyObjectType() {
	_ = this_
}

// Java constructor: ProxyObjectType(com.google.javascript.rhino.jstype.JSTypeRegistry registry, com.google.javascript.rhino.jstype.JSType referencedType)
func NewProxyObjectType_2(registry *JSTypeRegistry, referencedType JSType) ProxyObjectType {
	this_ := &ProxyObjectTypeBase{}
	this_.BindSelf_(this_)
	InitProxyObjectType_2(this_, registry, referencedType)
	return this_
}

func InitProxyObjectType_2(this_ *ProxyObjectTypeBase, registry *JSTypeRegistry, referencedType JSType) {
	_ = registry
	_ = referencedType
	InitProxyObjectType_4(this_, registry, referencedType, nil, int32(0))
}

// Java constructor: ProxyObjectType(com.google.javascript.rhino.jstype.JSTypeRegistry registry, com.google.javascript.rhino.jstype.JSType referencedType, com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable TemplateTypeMap templateTypeMap)
func NewProxyObjectType_3(registry *JSTypeRegistry, referencedType JSType, templateTypeMap *TemplateTypeMap) ProxyObjectType {
	this_ := &ProxyObjectTypeBase{}
	this_.BindSelf_(this_)
	InitProxyObjectType_3(this_, registry, referencedType, templateTypeMap)
	return this_
}

func InitProxyObjectType_3(this_ *ProxyObjectTypeBase, registry *JSTypeRegistry, referencedType JSType, templateTypeMap *TemplateTypeMap) {
	_ = registry
	_ = referencedType
	_ = templateTypeMap
	InitProxyObjectType_4(this_, registry, referencedType, templateTypeMap, int32(0))
}

// Java constructor: ProxyObjectType(com.google.javascript.rhino.jstype.JSTypeRegistry registry, com.google.javascript.rhino.jstype.JSType referencedType, com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable TemplateTypeMap templateTypeMap, int extraKindMask)
func NewProxyObjectType_4(registry *JSTypeRegistry, referencedType JSType, templateTypeMap *TemplateTypeMap, extraKindMask int32) ProxyObjectType {
	this_ := &ProxyObjectTypeBase{}
	this_.BindSelf_(this_)
	InitProxyObjectType_4(this_, registry, referencedType, templateTypeMap, extraKindMask)
	return this_
}

func InitProxyObjectType_4(this_ *ProxyObjectTypeBase, registry *JSTypeRegistry, referencedType JSType, templateTypeMap *TemplateTypeMap, extraKindMask int32) {
	_ = registry
	_ = referencedType
	_ = templateTypeMap
	_ = extraKindMask
	InitObjectType_3(&this_.ObjectTypeBase, registry, templateTypeMap, (int32(1) | extraKindMask))
	this_.initFields_ProxyObjectType()
	this_.SetReferencedType(j2gort.CheckNotNullIface(referencedType))
	registry.GetResolver().ResolveIfClosed(this_.JSTypeBase.self_, StaticGet_ProxyObjectType_TYPE_CLASS())
}

// Java method: ProxyObjectType.hasKindMaskProxy(int mask)
func (this_ *ProxyObjectTypeBase) HasKindMaskProxy(mask int32) bool {
	_ = mask
	return this_.referencedType.HasKindMask(mask)
}

// Java method: ProxyObjectType.getTypeClass()
func (this_ *ProxyObjectTypeBase) GetTypeClass() JSTypeClass {
	return StaticGet_ProxyObjectType_TYPE_CLASS()
}

// Java method: ProxyObjectType.getPropertyKind(com.google.javascript.rhino.jstype.Property.Key propertyName, boolean autobox)
func (this_ *ProxyObjectTypeBase) GetPropertyKind_2_Property_Key_bool(propertyName Property_Key, autobox bool) JSType_HasPropertyKind {
	_ = propertyName
	_ = autobox
	return this_.referencedType.GetPropertyKind_2_Property_Key_bool(propertyName, autobox)
}

// Java method: ProxyObjectType.getPropertyKind(java.lang.String propertyName, boolean autobox)
func (this_ *ProxyObjectTypeBase) GetPropertyKind_2_string_bool(propertyName string, autobox bool) JSType_HasPropertyKind {
	_ = propertyName
	_ = autobox
	return this_.referencedType.GetPropertyKind_2_string_bool(propertyName, autobox)
}

// Java method: ProxyObjectType.getPropertyMap()
func (this_ *ProxyObjectTypeBase) GetPropertyMap() *PropertyMap {
	if this_.referencedObjType == nil || this_.referencedObjType.IsJavaNil_() {
		return PropertyMap_ImmutableEmptyMap()
	} else {
		return this_.referencedObjType.GetPropertyMap()
	}
}

// Java method: ProxyObjectType.getReferencedTypeInternal()
func (this_ *ProxyObjectTypeBase) GetReferencedTypeInternal() JSType {
	return this_.referencedType
}

// Java method: ProxyObjectType.getReferencedObjTypeInternal()
func (this_ *ProxyObjectTypeBase) GetReferencedObjTypeInternal() ObjectType {
	return this_.referencedObjType
}

// Java method: ProxyObjectType.setReferencedType(com.google.javascript.rhino.jstype.JSType referencedType)
func (this_ *ProxyObjectTypeBase) SetReferencedType(referencedType JSType) {
	_ = referencedType
	j2gort.CheckNotNullIface(referencedType)
	this_.referencedType = referencedType
	var objectType ObjectType
	_ = objectType
	var ok_1 bool
	objectType, ok_1 = j2gort.TryCast[ObjectType](referencedType)
	if ok_1 {
		this_.referencedObjType = objectType
	} else {
		this_.referencedObjType = nil
	}
}

// Java method: ProxyObjectType.loosenTypecheckingDueToForwardReferencedSupertype()
func (this_ *ProxyObjectTypeBase) LoosenTypecheckingDueToForwardReferencedSupertype() bool {
	return this_.referencedType.LoosenTypecheckingDueToForwardReferencedSupertype()
}

// Java method: ProxyObjectType.setValidator(com.google.common.base.Predicate<com.google.javascript.rhino.jstype.JSType> validator)
func (this_ *ProxyObjectTypeBase) SetValidator(validator j2gort.Ext_com_google_common_base_Predicate) bool {
	_ = validator
	return this_.referencedType.SetValidator(validator)
}

// Java method: ProxyObjectType.getReferenceName()
func (this_ *ProxyObjectTypeBase) GetReferenceName() string {
	if this_.referencedObjType == nil || this_.referencedObjType.IsJavaNil_() {
		return ""
	} else {
		return this_.referencedObjType.GetReferenceName()
	}
}

// Java method: ProxyObjectType.matchesNumberContext()
func (this_ *ProxyObjectTypeBase) MatchesNumberContext() bool {
	return this_.referencedType.MatchesNumberContext()
}

// Java method: ProxyObjectType.matchesStringContext()
func (this_ *ProxyObjectTypeBase) MatchesStringContext() bool {
	return this_.referencedType.MatchesStringContext()
}

// Java method: ProxyObjectType.matchesSymbolContext()
func (this_ *ProxyObjectTypeBase) MatchesSymbolContext() bool {
	return this_.referencedType.MatchesSymbolContext()
}

// Java method: ProxyObjectType.matchesObjectContext()
func (this_ *ProxyObjectTypeBase) MatchesObjectContext() bool {
	return this_.referencedType.MatchesObjectContext()
}

// Java method: ProxyObjectType.canBeCalled()
func (this_ *ProxyObjectTypeBase) CanBeCalled() bool {
	return this_.referencedType.CanBeCalled()
}

// Java method: ProxyObjectType.isStructuralType()
func (this_ *ProxyObjectTypeBase) IsStructuralType() bool {
	return this_.referencedType.IsStructuralType()
}

// Java method: ProxyObjectType.isUnknownType()
func (this_ *ProxyObjectTypeBase) IsUnknownType() bool {
	return this_.referencedType.IsUnknownType()
}

// Java method: ProxyObjectType.isNullable()
func (this_ *ProxyObjectTypeBase) IsNullable() bool {
	return this_.referencedType.IsNullable()
}

// Java method: ProxyObjectType.isVoidable()
func (this_ *ProxyObjectTypeBase) IsVoidable() bool {
	return this_.referencedType.IsVoidable()
}

// Java method: ProxyObjectType.toMaybeEnumType()
func (this_ *ProxyObjectTypeBase) ToMaybeEnumType() *EnumType {
	return this_.referencedType.ToMaybeEnumType()
}

// Java method: ProxyObjectType.isNominalType()
func (this_ *ProxyObjectTypeBase) IsNominalType() bool {
	return this_.referencedType.IsNominalType()
}

// Java method: ProxyObjectType.isInstanceType()
func (this_ *ProxyObjectTypeBase) IsInstanceType() bool {
	return this_.referencedType.IsInstanceType()
}

// Java method: ProxyObjectType.isStruct()
func (this_ *ProxyObjectTypeBase) IsStruct() bool {
	return this_.referencedType.IsStruct()
}

// Java method: ProxyObjectType.isDict()
func (this_ *ProxyObjectTypeBase) IsDict() bool {
	return this_.referencedType.IsDict()
}

// Java method: ProxyObjectType.toMaybeRecordType()
func (this_ *ProxyObjectTypeBase) ToMaybeRecordType() *RecordType {
	return this_.referencedType.ToMaybeRecordType()
}

// Java method: ProxyObjectType.toMaybeUnionType()
func (this_ *ProxyObjectTypeBase) ToMaybeUnionType() *UnionType {
	return this_.referencedType.ToMaybeUnionType()
}

// Java method: ProxyObjectType.toMaybeFunctionType()
func (this_ *ProxyObjectTypeBase) ToMaybeFunctionType() *FunctionType {
	return this_.referencedType.ToMaybeFunctionType()
}

// Java method: ProxyObjectType.toMaybeEnumElementType()
func (this_ *ProxyObjectTypeBase) ToMaybeEnumElementType() *EnumElementType {
	return this_.referencedType.ToMaybeEnumElementType()
}

// Java method: ProxyObjectType.toMaybeKnownSymbolType()
func (this_ *ProxyObjectTypeBase) ToMaybeKnownSymbolType() *KnownSymbolType {
	return this_.referencedType.ToMaybeKnownSymbolType()
}

// Java method: ProxyObjectType.testForEquality(com.google.javascript.rhino.jstype.JSType that)
func (this_ *ProxyObjectTypeBase) TestForEquality(that JSType) pkg_base.Tri {
	_ = that
	return this_.referencedType.TestForEquality(that)
}

// Java method: ProxyObjectType.getOwnerFunction()
func (this_ *ProxyObjectTypeBase) GetOwnerFunction() *FunctionType {
	if this_.referencedObjType == nil || this_.referencedObjType.IsJavaNil_() {
		return nil
	} else {
		return this_.referencedObjType.GetOwnerFunction()
	}
}

// Java method: ProxyObjectType.getCtorImplementedInterfaces()
func (this_ *ProxyObjectTypeBase) GetCtorImplementedInterfaces() any {
	if this_.referencedObjType == nil || this_.referencedObjType.IsJavaNil_() {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
	} else {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.referencedObjType.GetCtorImplementedInterfaces())
	}
}

// Java method: ProxyObjectType.getCtorExtendedInterfaces()
func (this_ *ProxyObjectTypeBase) GetCtorExtendedInterfaces() any {
	if this_.referencedObjType == nil || this_.referencedObjType.IsJavaNil_() {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_java_util_Collections_EmptyList())
	} else {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.referencedObjType.GetCtorExtendedInterfaces())
	}
}

// Java method: ProxyObjectType.recursionUnsafeHashCode()
func (this_ *ProxyObjectTypeBase) RecursionUnsafeHashCode() int32 {
	return this_.referencedType.HashCode()
}

// Java method: ProxyObjectType.appendTo(com.google.javascript.rhino.jstype.TypeStringBuilder sb)
func (this_ *ProxyObjectTypeBase) AppendTo(sb *TypeStringBuilder) {
	_ = sb
	this_.referencedType.AppendTo(sb)
}

// Java method: ProxyObjectType.getImplicitPrototype()
func (this_ *ProxyObjectTypeBase) GetImplicitPrototype() ObjectType {
	if this_.referencedObjType == nil || this_.referencedObjType.IsJavaNil_() {
		return nil
	} else {
		return this_.referencedObjType.GetImplicitPrototype()
	}
}

// Java method: ProxyObjectType.defineProperty(java.lang.String propertyName, com.google.javascript.rhino.jstype.JSType type, boolean inferred, com.google.javascript.rhino.Node propertyNode)
func (this_ *ProxyObjectTypeBase) DefineProperty_4_string_JSType_bool_Node(propertyName string, type_ JSType, inferred bool, propertyNode Node) bool {
	_ = propertyName
	_ = type_
	_ = inferred
	_ = propertyNode
	return ((this_.referencedObjType == nil || this_.referencedObjType.IsJavaNil_()) || this_.referencedObjType.DefineProperty_4_string_JSType_bool_Node(propertyName, type_, inferred, propertyNode))
}

// Java method: ProxyObjectType.defineProperty(com.google.javascript.rhino.jstype.Property.Key propertyName, com.google.javascript.rhino.jstype.JSType type, boolean inferred, com.google.javascript.rhino.Node propertyNode)
func (this_ *ProxyObjectTypeBase) DefineProperty_4_Property_Key_JSType_bool_Node(propertyName Property_Key, type_ JSType, inferred bool, propertyNode Node) bool {
	_ = propertyName
	_ = type_
	_ = inferred
	_ = propertyNode
	return ((this_.referencedObjType == nil || this_.referencedObjType.IsJavaNil_()) || this_.referencedObjType.DefineProperty_4_Property_Key_JSType_bool_Node(propertyName, type_, inferred, propertyNode))
}

// Java method: ProxyObjectType.findPropertyTypeWithoutConsideringTemplateTypes(com.google.javascript.rhino.jstype.Property.Key propertyName)
func (this_ *ProxyObjectTypeBase) FindPropertyTypeWithoutConsideringTemplateTypes_1_Property_Key(propertyName Property_Key) JSType {
	_ = propertyName
	return this_.referencedType.FindPropertyType_1_Property_Key(propertyName)
}

// Java method: ProxyObjectType.findPropertyTypeWithoutConsideringTemplateTypes(java.lang.String propertyName)
func (this_ *ProxyObjectTypeBase) FindPropertyTypeWithoutConsideringTemplateTypes_1_string(propertyName string) JSType {
	_ = propertyName
	return this_.referencedType.FindPropertyType_1_string(propertyName)
}

// Java method: ProxyObjectType.getJSDocInfo()
func (this_ *ProxyObjectTypeBase) GetJSDocInfo() *JSDocInfo {
	return this_.referencedType.GetJSDocInfo()
}

// Java method: ProxyObjectType.setJSDocInfo(com.google.javascript.rhino.JSDocInfo info)
func (this_ *ProxyObjectTypeBase) SetJSDocInfo(info *JSDocInfo) {
	_ = info
	if this_.referencedObjType != nil && !this_.referencedObjType.IsJavaNil_() {
		this_.referencedObjType.SetJSDocInfo(info)
	}
}

// Java method: ProxyObjectType.setPropertyJSDocInfo(java.lang.String propertyName, com.google.javascript.rhino.JSDocInfo info)
func (this_ *ProxyObjectTypeBase) SetPropertyJSDocInfo_2_string_JSDocInfo(propertyName string, info *JSDocInfo) {
	_ = propertyName
	_ = info
	if this_.referencedObjType != nil && !this_.referencedObjType.IsJavaNil_() {
		this_.referencedObjType.SetPropertyJSDocInfo_2_string_JSDocInfo(propertyName, info)
	} else {
		j2gort.CheckNotNullString(propertyName)
	}
}

// Java method: ProxyObjectType.setPropertyJSDocInfo(com.google.javascript.rhino.jstype.Property.Key propertyName, com.google.javascript.rhino.JSDocInfo info)
func (this_ *ProxyObjectTypeBase) SetPropertyJSDocInfo_2_Property_Key_JSDocInfo(propertyName Property_Key, info *JSDocInfo) {
	_ = propertyName
	_ = info
	if this_.referencedObjType != nil && !this_.referencedObjType.IsJavaNil_() {
		this_.referencedObjType.SetPropertyJSDocInfo_2_Property_Key_JSDocInfo(propertyName, info)
	}
}

// Java method: ProxyObjectType.getConstructor()
func (this_ *ProxyObjectTypeBase) GetConstructor() *FunctionType {
	if this_.referencedObjType == nil || this_.referencedObjType.IsJavaNil_() {
		return nil
	} else {
		return this_.referencedObjType.GetConstructor()
	}
}

// Java method: ProxyObjectType.getTemplateTypes()
func (this_ *ProxyObjectTypeBase) GetTemplateTypes() j2gort.Ext_com_google_common_collect_ImmutableList {
	if this_.referencedObjType == nil || this_.referencedObjType.IsJavaNil_() {
		return nil
	} else {
		return this_.referencedObjType.GetTemplateTypes()
	}
}

// Java method: ProxyObjectType.getTemplateParamCount()
func (this_ *ProxyObjectTypeBase) GetTemplateParamCount() int32 {
	return this_.referencedType.GetTemplateParamCount()
}

// Java method: ProxyObjectType.visitReferenceType(com.google.javascript.rhino.jstype.Visitor<T> visitor)
func (this_ *ProxyObjectTypeBase) VisitReferenceType(visitor Visitor) any {
	_ = visitor
	return this_.referencedType.Visit_1(visitor)
}

// Java method: ProxyObjectType.visit(com.google.javascript.rhino.jstype.Visitor<T> visitor)
func (this_ *ProxyObjectTypeBase) Visit_1(visitor Visitor) any {
	_ = visitor
	return visitor.CaseProxyObjectType(this_.self_ProxyObjectType_())
}

// Java method: ProxyObjectType.visit(com.google.javascript.rhino.jstype.RelationshipVisitor<T> visitor, com.google.javascript.rhino.jstype.JSType that)
func (this_ *ProxyObjectTypeBase) Visit_2(visitor RelationshipVisitor, that JSType) any {
	_ = visitor
	_ = that
	return this_.referencedType.Visit_2(visitor, that)
}

// Java method: ProxyObjectType.resolveInternal(com.google.javascript.rhino.ErrorReporter reporter)
func (this_ *ProxyObjectTypeBase) ResolveInternal(reporter ErrorReporter) JSType {
	_ = reporter
	this_.SetReferencedType(this_.referencedType.Resolve(reporter))
	return this_.JSTypeBase.self_
}

// Java method: ProxyObjectType.getTypeOfThis()
func (this_ *ProxyObjectTypeBase) GetTypeOfThis() JSType {
	if this_.referencedObjType != nil && !this_.referencedObjType.IsJavaNil_() {
		return this_.referencedObjType.GetTypeOfThis()
	}
	return this_.ObjectTypeBase.GetTypeOfThis()
}

// Java method: ProxyObjectType.collapseUnion()
func (this_ *ProxyObjectTypeBase) CollapseUnion() JSType {
	if this_.referencedType.IsUnionType() {
		return this_.referencedType.CollapseUnion()
	}
	return this_.JSTypeBase.self_
}

// Java method: ProxyObjectType.matchConstraint(com.google.javascript.rhino.jstype.JSType constraint)
func (this_ *ProxyObjectTypeBase) MatchConstraint(constraint JSType) {
	_ = constraint
	this_.referencedType.MatchConstraint(constraint)
}

// Java method: ProxyObjectType.toMaybeTemplatizedType()
func (this_ *ProxyObjectTypeBase) ToMaybeTemplatizedType() *TemplatizedType {
	return this_.referencedType.ToMaybeTemplatizedType()
}

// Java method: ProxyObjectType.toMaybeTemplateType()
func (this_ *ProxyObjectTypeBase) ToMaybeTemplateType() *TemplateType {
	return this_.referencedType.ToMaybeTemplateType()
}

// Java method: ProxyObjectType.hasAnyTemplateTypesInternal()
func (this_ *ProxyObjectTypeBase) HasAnyTemplateTypesInternal() bool {
	return this_.referencedType.HasAnyTemplateTypes()
}

// Java method: ProxyObjectType.getTemplateTypeMap()
func (this_ *ProxyObjectTypeBase) GetTemplateTypeMap() *TemplateTypeMap {
	return this_.referencedType.GetTemplateTypeMap()
}
