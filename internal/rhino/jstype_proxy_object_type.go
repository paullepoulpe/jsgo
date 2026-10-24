// Port of com/google/javascript/rhino/jstype/ProxyObjectType.java @ closure-compiler v20261015

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

// ProxyObjectType is an object type which uses composition to delegate all
// calls.
//
// Java: class ProxyObjectType extends ObjectType (subclassed by NamedType,
// TemplateType and TemplatizedType).
type ProxyObjectType interface {
	ObjectType

	GetReferencedTypeInternal() JSType
	GetReferencedObjTypeInternal() ObjectType
	SetReferencedType(referencedType JSType)
	// Java: ProxyObjectType.visitReferenceType(Visitor<T>)
	VisitReferenceType(visitor Visitor) any

	// GetProxyObjectTypeBase gives access to the base through the interface.
	//
	// Go-only: language-adapter — field access through the interface
	// (PORTING.md §4).
	GetProxyObjectTypeBase() *ProxyObjectTypeBase
}

// ProxyObjectTypeBase holds ProxyObjectType's fields and methods; plain
// proxies are *ProxyObjectTypeBase values.
type ProxyObjectTypeBase struct {
	ObjectTypeBase

	referencedType    JSType
	referencedObjType ObjectType
}

// Java: ProxyObjectType.TYPE_CLASS
const ProxyObjectType_TYPE_CLASS = JSTypeClass_PROXY_OBJECT

// Java: ProxyObjectType(JSTypeRegistry, JSType)
func NewProxyObjectType(registry *JSTypeRegistry, referencedType JSType) *ProxyObjectTypeBase {
	return NewProxyObjectTypeWithTemplateTypeMap(registry, referencedType, nil)
}

// Java: ProxyObjectType(JSTypeRegistry, JSType, TemplateTypeMap)
func NewProxyObjectTypeWithTemplateTypeMap(
	registry *JSTypeRegistry, referencedType JSType, templateTypeMap *TemplateTypeMap) *ProxyObjectTypeBase {
	t := &ProxyObjectTypeBase{}
	t.objFlags = objKindProxy
	t.InitProxyObjectTypeBaseWithTemplateTypeMap(t, registry, referencedType, templateTypeMap)
	return t
}

// InitProxyObjectTypeBase is the Go form of the ProxyObjectType(JSTypeRegistry,
// JSType) constructor for subclasses.
//
// Java: ProxyObjectType(JSTypeRegistry, JSType)
func (p *ProxyObjectTypeBase) InitProxyObjectTypeBase(self ProxyObjectType, registry *JSTypeRegistry, referencedType JSType) {
	p.InitProxyObjectTypeBaseWithTemplateTypeMap(self, registry, referencedType, nil)
}

// InitProxyObjectTypeBaseWithTemplateTypeMap is the Go form of the
// ProxyObjectType(JSTypeRegistry, JSType, TemplateTypeMap) constructor for
// subclasses.
//
// Java: ProxyObjectType(JSTypeRegistry, JSType, TemplateTypeMap)
func (p *ProxyObjectTypeBase) InitProxyObjectTypeBaseWithTemplateTypeMap(
	self ProxyObjectType, registry *JSTypeRegistry, referencedType JSType, templateTypeMap *TemplateTypeMap) {
	p.InitObjectTypeBaseWithTemplateTypeMap(self, registry, templateTypeMap)
	p.SetReferencedType(jutil.CheckNotNull(referencedType))

	registry.GetResolver().ResolveIfClosed(self, ProxyObjectType_TYPE_CLASS)
}

// GetProxyObjectTypeBase gives access to the base through the interface.
//
// Go-only: language-adapter — field access through the interface (PORTING.md §4).
func (p *ProxyObjectTypeBase) GetProxyObjectTypeBase() *ProxyObjectTypeBase { return p }

// Java: ProxyObjectType.getTypeClass()
func (p *ProxyObjectTypeBase) GetTypeClass() JSTypeClass {
	return ProxyObjectType_TYPE_CLASS
}

// Java: ProxyObjectType.getPropertyKind(Property.Key, boolean)
func (p *ProxyObjectTypeBase) GetPropertyKindByKeyWithAutobox(propertyName Property_Key, autobox bool) JSType_HasPropertyKind {
	return p.referencedType.GetPropertyKindByKeyWithAutobox(propertyName, autobox)
}

// Java: ProxyObjectType.getPropertyMap()
func (p *ProxyObjectTypeBase) GetPropertyMap() *PropertyMap {
	if p.referencedObjType == nil {
		return PropertyMap_ImmutableEmptyMap()
	}
	return p.referencedObjType.GetPropertyMap()
}

// Java: ProxyObjectType.getReferencedTypeInternal()
func (p *ProxyObjectTypeBase) GetReferencedTypeInternal() JSType {
	return p.referencedType
}

// Java: ProxyObjectType.getReferencedObjTypeInternal()
func (p *ProxyObjectTypeBase) GetReferencedObjTypeInternal() ObjectType {
	return p.referencedObjType
}

// Java: ProxyObjectType.setReferencedType(JSType)
func (p *ProxyObjectTypeBase) SetReferencedType(referencedType JSType) {
	jutil.CheckNotNull(referencedType)
	p.referencedType = referencedType
	if objectType, ok := referencedType.(ObjectType); ok {
		p.referencedObjType = objectType
	} else {
		p.referencedObjType = nil
	}
}

// Java: ProxyObjectType.loosenTypecheckingDueToForwardReferencedSupertype()
func (p *ProxyObjectTypeBase) LoosenTypecheckingDueToForwardReferencedSupertype() bool {
	return p.referencedType.LoosenTypecheckingDueToForwardReferencedSupertype()
}

// Java: ProxyObjectType.setValidator(Predicate<JSType>)
func (p *ProxyObjectTypeBase) SetValidator(validator func(JSType) bool) bool {
	// The referenced type might have specialized behavior for validation, e.g.
	// NamedType defers validation until after named type resolution.
	return p.referencedType.SetValidator(validator)
}

// Java: ProxyObjectType.getReferenceName()
func (p *ProxyObjectTypeBase) GetReferenceName() *string {
	if p.referencedObjType == nil {
		empty := ""
		return &empty
	}
	return p.referencedObjType.GetReferenceName()
}

// Java: ProxyObjectType.matchesNumberContext()
func (p *ProxyObjectTypeBase) MatchesNumberContext() bool {
	return p.referencedType.MatchesNumberContext()
}

// Java: ProxyObjectType.matchesStringContext()
func (p *ProxyObjectTypeBase) MatchesStringContext() bool {
	return p.referencedType.MatchesStringContext()
}

// Java: ProxyObjectType.matchesSymbolContext()
func (p *ProxyObjectTypeBase) MatchesSymbolContext() bool {
	return p.referencedType.MatchesSymbolContext()
}

// Java: ProxyObjectType.matchesObjectContext()
func (p *ProxyObjectTypeBase) MatchesObjectContext() bool {
	return p.referencedType.MatchesObjectContext()
}

// Java: ProxyObjectType.canBeCalled()
func (p *ProxyObjectTypeBase) CanBeCalled() bool {
	return p.referencedType.CanBeCalled()
}

// Java: ProxyObjectType.isStructuralType()
func (p *ProxyObjectTypeBase) IsStructuralType() bool {
	return p.referencedType.IsStructuralType()
}

// Java: ProxyObjectType.isNoType()
func (p *ProxyObjectTypeBase) IsNoType() bool {
	return p.referencedType.IsNoType()
}

// Java: ProxyObjectType.isNoObjectType()
func (p *ProxyObjectTypeBase) IsNoObjectType() bool {
	return p.referencedType.IsNoObjectType()
}

// Java: ProxyObjectType.isNoResolvedType()
func (p *ProxyObjectTypeBase) IsNoResolvedType() bool {
	return p.referencedType.IsNoResolvedType()
}

// Java: ProxyObjectType.isUnknownType()
func (p *ProxyObjectTypeBase) IsUnknownType() bool {
	return p.referencedType.IsUnknownType()
}

// Java: ProxyObjectType.isCheckedUnknownType()
func (p *ProxyObjectTypeBase) IsCheckedUnknownType() bool {
	return p.referencedType.IsCheckedUnknownType()
}

// Java: ProxyObjectType.isNullable()
func (p *ProxyObjectTypeBase) IsNullable() bool {
	return p.referencedType.IsNullable()
}

// Java: ProxyObjectType.isVoidable()
func (p *ProxyObjectTypeBase) IsVoidable() bool {
	return p.referencedType.IsVoidable()
}

// Java: JSType.isEnumType()
func (p *ProxyObjectTypeBase) IsEnumType() bool {
	return p.referencedType.IsEnumType()
}

// Java: ProxyObjectType.toMaybeEnumType()
func (p *ProxyObjectTypeBase) ToMaybeEnumType() *EnumType {
	return p.referencedType.ToMaybeEnumType()
}

// Java: ProxyObjectType.isConstructor()
func (p *ProxyObjectTypeBase) IsConstructor() bool {
	return p.referencedType.IsConstructor()
}

// Java: ProxyObjectType.isNominalType()
func (p *ProxyObjectTypeBase) IsNominalType() bool {
	return p.referencedType.IsNominalType()
}

// Java: ProxyObjectType.isInstanceType()
func (p *ProxyObjectTypeBase) IsInstanceType() bool {
	return p.referencedType.IsInstanceType()
}

// Java: ProxyObjectType.isInterface()
func (p *ProxyObjectTypeBase) IsInterface() bool {
	return p.referencedType.IsInterface()
}

// Java: ProxyObjectType.isOrdinaryFunction()
func (p *ProxyObjectTypeBase) IsOrdinaryFunction() bool {
	return p.referencedType.IsOrdinaryFunction()
}

// Java: ProxyObjectType.isAllType()
func (p *ProxyObjectTypeBase) IsAllType() bool {
	return p.referencedType.IsAllType()
}

// Java: ProxyObjectType.isStruct()
func (p *ProxyObjectTypeBase) IsStruct() bool {
	return p.referencedType.IsStruct()
}

// Java: ProxyObjectType.isDict()
func (p *ProxyObjectTypeBase) IsDict() bool {
	return p.referencedType.IsDict()
}

// Java: ProxyObjectType.isNativeObjectType()
func (p *ProxyObjectTypeBase) IsNativeObjectType() bool {
	if p.referencedObjType == nil {
		return false
	}
	return p.referencedObjType.IsNativeObjectType()
}

// Java: JSType.isRecordType()
func (p *ProxyObjectTypeBase) IsRecordType() bool {
	return p.referencedType.IsRecordType()
}

// Java: ProxyObjectType.toMaybeRecordType()
func (p *ProxyObjectTypeBase) ToMaybeRecordType() *RecordType {
	return p.referencedType.ToMaybeRecordType()
}

// Java: JSType.isUnionType()
func (p *ProxyObjectTypeBase) IsUnionType() bool {
	return p.referencedType.IsUnionType()
}

// Java: ProxyObjectType.toMaybeUnionType()
func (p *ProxyObjectTypeBase) ToMaybeUnionType() *UnionType {
	return p.referencedType.ToMaybeUnionType()
}

// Java: JSType.isFunctionType()
func (p *ProxyObjectTypeBase) IsFunctionType() bool {
	return p.referencedType.IsFunctionType()
}

// Java: ProxyObjectType.toMaybeFunctionType()
func (p *ProxyObjectTypeBase) ToMaybeFunctionType() FunctionType {
	return p.referencedType.ToMaybeFunctionType()
}

// Java: JSType.isEnumElementType()
func (p *ProxyObjectTypeBase) IsEnumElementType() bool {
	return p.referencedType.IsEnumElementType()
}

// Java: ProxyObjectType.toMaybeEnumElementType()
func (p *ProxyObjectTypeBase) ToMaybeEnumElementType() *EnumElementType {
	return p.referencedType.ToMaybeEnumElementType()
}

// Java: ProxyObjectType.toMaybeKnownSymbolType()
func (p *ProxyObjectTypeBase) ToMaybeKnownSymbolType() *KnownSymbolType {
	return p.referencedType.ToMaybeKnownSymbolType()
}

// Java: ProxyObjectType.testForEquality(JSType)
func (p *ProxyObjectTypeBase) TestForEquality(that JSType) base.Tri {
	return p.referencedType.TestForEquality(that)
}

// Java: ProxyObjectType.getOwnerFunction()
func (p *ProxyObjectTypeBase) GetOwnerFunction() FunctionType {
	if p.referencedObjType == nil {
		return nil
	}
	return p.referencedObjType.GetOwnerFunction()
}

// Java: ProxyObjectType.getCtorImplementedInterfaces()
func (p *ProxyObjectTypeBase) GetCtorImplementedInterfaces() []ObjectType {
	if p.referencedObjType == nil {
		return []ObjectType{}
	}
	return p.referencedObjType.GetCtorImplementedInterfaces()
}

// Java: ProxyObjectType.getCtorExtendedInterfaces()
func (p *ProxyObjectTypeBase) GetCtorExtendedInterfaces() []ObjectType {
	if p.referencedObjType == nil {
		return []ObjectType{}
	}
	return p.referencedObjType.GetCtorExtendedInterfaces()
}

// Java: ProxyObjectType.recursionUnsafeHashCode()
func (p *ProxyObjectTypeBase) RecursionUnsafeHashCode() int32 {
	return p.referencedType.HashCode()
}

// Java: ProxyObjectType.appendTo(TypeStringBuilder)
func (p *ProxyObjectTypeBase) AppendTo(sb *TypeStringBuilder) {
	p.referencedType.AppendTo(sb)
}

// Java: ProxyObjectType.getImplicitPrototype()
func (p *ProxyObjectTypeBase) GetImplicitPrototype() ObjectType {
	if p.referencedObjType == nil {
		return nil
	}
	return p.referencedObjType.GetImplicitPrototype()
}

// Java: ProxyObjectType.defineProperty(Property.Key, JSType, boolean, Node)
func (p *ProxyObjectTypeBase) DefinePropertyByKey(propertyName Property_Key, type_ JSType, inferred bool, propertyNode Node) bool {
	return p.referencedObjType == nil ||
		p.referencedObjType.DefinePropertyByKey(propertyName, type_, inferred, propertyNode)
}

// Java: ProxyObjectType.findPropertyTypeWithoutConsideringTemplateTypes(Property.Key)
func (p *ProxyObjectTypeBase) FindPropertyTypeWithoutConsideringTemplateTypes(propertyName Property_Key) JSType {
	return p.referencedType.FindPropertyTypeByKey(propertyName)
}

// Java: ProxyObjectType.getJSDocInfo()
func (p *ProxyObjectTypeBase) GetJSDocInfo() *JSDocInfo {
	return p.referencedType.GetJSDocInfo()
}

// Java: ProxyObjectType.setJSDocInfo(JSDocInfo)
func (p *ProxyObjectTypeBase) SetJSDocInfo(info *JSDocInfo) {
	if p.referencedObjType != nil {
		p.referencedObjType.SetJSDocInfo(info)
	}
}

// Java: ProxyObjectType.setPropertyJSDocInfo(Property.Key, JSDocInfo)
func (p *ProxyObjectTypeBase) SetPropertyJSDocInfoByKey(propertyName Property_Key, info *JSDocInfo) {
	if p.referencedObjType != nil {
		p.referencedObjType.SetPropertyJSDocInfoByKey(propertyName, info)
	}
}

// Java: ProxyObjectType.getConstructor()
func (p *ProxyObjectTypeBase) GetConstructor() FunctionType {
	if p.referencedObjType == nil {
		return nil
	}
	return p.referencedObjType.GetConstructor()
}

// Java: ProxyObjectType.getTemplateTypes()
func (p *ProxyObjectTypeBase) GetTemplateTypes() []JSType {
	if p.referencedObjType == nil {
		return nil
	}
	return p.referencedObjType.GetTemplateTypes()
}

// Java: ProxyObjectType.getTemplateParamCount()
func (p *ProxyObjectTypeBase) GetTemplateParamCount() int {
	return p.referencedType.GetTemplateParamCount()
}

// Java: ProxyObjectType.visitReferenceType(Visitor<T>)
func (p *ProxyObjectTypeBase) VisitReferenceType(visitor Visitor) any {
	return p.referencedType.Visit(visitor)
}

// Java: ProxyObjectType.visit(Visitor<T>)
func (p *ProxyObjectTypeBase) Visit(visitor Visitor) any {
	return visitor.CaseProxyObjectType(p.self().(ProxyObjectType))
}

// Java: ProxyObjectType.visit(RelationshipVisitor<T>, JSType)
func (p *ProxyObjectTypeBase) VisitRelationship(visitor RelationshipVisitor, that JSType) any {
	return p.referencedType.VisitRelationship(visitor, that)
}

// Java: ProxyObjectType.resolveInternal(ErrorReporter)
func (p *ProxyObjectTypeBase) ResolveInternal(reporter ErrorReporter) JSType {
	p.SetReferencedType(p.referencedType.Resolve(reporter))
	return p.self()
}

// Java: ProxyObjectType.getTypeOfThis()
func (p *ProxyObjectTypeBase) GetTypeOfThis() JSType {
	if p.referencedObjType != nil {
		return p.referencedObjType.GetTypeOfThis()
	}
	return p.ObjectTypeBase.GetTypeOfThis()
}

// Java: ProxyObjectType.collapseUnion()
func (p *ProxyObjectTypeBase) CollapseUnion() JSType {
	if p.referencedType.IsUnionType() {
		return p.referencedType.CollapseUnion()
	}
	return p.self()
}

// Java: ProxyObjectType.matchConstraint(JSType)
func (p *ProxyObjectTypeBase) MatchConstraint(constraint JSType) {
	p.referencedType.MatchConstraint(constraint)
}

// Java: JSType.isTemplatizedType()
func (p *ProxyObjectTypeBase) IsTemplatizedType() bool {
	return p.referencedType.IsTemplatizedType()
}

// Java: ProxyObjectType.toMaybeTemplatizedType()
func (p *ProxyObjectTypeBase) ToMaybeTemplatizedType() *TemplatizedType {
	return p.referencedType.ToMaybeTemplatizedType()
}

// Java: JSType.isTemplateType()
func (p *ProxyObjectTypeBase) IsTemplateType() bool {
	return p.referencedType.IsTemplateType()
}

// Java: ProxyObjectType.toMaybeTemplateType()
func (p *ProxyObjectTypeBase) ToMaybeTemplateType() *TemplateType {
	return p.referencedType.ToMaybeTemplateType()
}

// Java: ProxyObjectType.hasAnyTemplateTypesInternal()
func (p *ProxyObjectTypeBase) HasAnyTemplateTypesInternal() bool {
	return p.referencedType.HasAnyTemplateTypes()
}

// Java: ProxyObjectType.getTemplateTypeMap()
func (p *ProxyObjectTypeBase) GetTemplateTypeMap() *TemplateTypeMap {
	return p.referencedType.GetTemplateTypeMap()
}
