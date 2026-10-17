// Port of com/google/javascript/rhino/jstype/TemplatizedType.java @ closure-compiler v20260415

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

// TemplatizedType is an object type with declared template types, such as
// Array<string>.
type TemplatizedType struct {
	ProxyObjectTypeBase

	// A cache of the type parameter values for this specialization.
	templateTypes []JSType
	// Whether all type parameter values for this specialization are `?`.
	isSpecializedOnlyWithUnknown bool

	replacer *TemplateTypeReplacer
}

// Java: TemplatizedType.TYPE_CLASS
const TemplatizedType_TYPE_CLASS = JSTypeClass_TEMPLATIZED

// Java: TemplatizedType(JSTypeRegistry, ObjectType, ImmutableList<JSType>)
func NewTemplatizedType(registry *JSTypeRegistry, objectType ObjectType, templateTypes []JSType) *TemplatizedType {
	t := &TemplatizedType{}
	t.InitProxyObjectTypeBaseWithTemplateTypeMap(
		t, registry, objectType, objectType.GetTemplateTypeMap().CopyFilledWithValues(templateTypes))

	builder := []JSType{}
	maybeIsSpecializedOnlyWithUnknown := true
	for _, newlyFilledTemplateKey := range objectType.GetTypeParameters() {
		resolvedType := t.GetTemplateTypeMap().GetResolvedTemplateType(newlyFilledTemplateKey)

		builder = append(builder, jutil.CheckNotNull(resolvedType))
		if maybeIsSpecializedOnlyWithUnknown {
			maybeIsSpecializedOnlyWithUnknown = t.GetNativeType(JSTypeNative_UNKNOWN_TYPE).Equals(resolvedType)
		}
	}
	t.templateTypes = builder
	t.isSpecializedOnlyWithUnknown = maybeIsSpecializedOnlyWithUnknown

	t.replacer = TemplateTypeReplacer_ForPartialReplacement(registry, t.GetTemplateTypeMap())

	registry.GetResolver().ResolveIfClosed(t, TemplatizedType_TYPE_CLASS)
	return t
}

// Java: TemplatizedType.getTypeClass()
func (t *TemplatizedType) GetTypeClass() JSTypeClass {
	return TemplatizedType_TYPE_CLASS
}

// GetCtorImplementedInterfaces: NOTE(dimvar): If getCtorImplementedInterfaces
// is implemented here, this is the correct implementation. The one inherited
// from ProxyObjectType is not correct because it doesn't instantiate the
// generic types. However, our unit tests don't actually ever call this
// method, so it could alternatively just throw an
// UnsupportedOperationException.
//
// Java: TemplatizedType.getCtorImplementedInterfaces()
func (t *TemplatizedType) GetCtorImplementedInterfaces() []ObjectType {
	resolvedImplementedInterfaces := jutil.NewObjectLinkedHashSet[ObjectType]()
	for _, obj := range t.GetReferencedObjTypeInternal().GetCtorImplementedInterfaces() {
		resolvedImplementedInterfaces.Add(JSType_VisitForJSType(obj, t.replacer).ToObjectType())
	}
	return resolvedImplementedInterfaces.ToSlice()
}

// Java: TemplatizedType.getCtorExtendedInterfaces()
func (t *TemplatizedType) GetCtorExtendedInterfaces() []ObjectType {
	resolvedExtendedInterfaces := jutil.NewObjectLinkedHashSet[ObjectType]()
	for _, obj := range t.GetReferencedObjTypeInternal().GetCtorExtendedInterfaces() {
		resolvedExtendedInterfaces.Add(JSType_VisitForJSType(obj, t.replacer).ToObjectType())
	}
	return resolvedExtendedInterfaces.ToSlice()
}

// Java: TemplatizedType.appendTo(TypeStringBuilder)
func (t *TemplatizedType) AppendTo(sb *TypeStringBuilder) {
	t.ProxyObjectTypeBase.AppendTo(sb)
	if len(t.templateTypes) != 0 {
		elements := make([]any, len(t.templateTypes))
		for i, e := range t.templateTypes {
			elements[i] = e
		}
		sb.Append("<").AppendAll(elements, ",").Append(">")
	}
}

// Java: TemplatizedType.recursionUnsafeHashCode()
func (t *TemplatizedType) RecursionUnsafeHashCode() int32 {
	baseHash := t.ProxyObjectTypeBase.RecursionUnsafeHashCode()

	// TODO(b/110224889): This case can probably be removed if `equals()` is
	// updated.
	if t.isSpecializedOnlyWithUnknown {
		return baseHash
	}
	return 31*(31*1+jutil.ListHashCode(t.templateTypes)) + baseHash
}

// Java: TemplatizedType.visit(Visitor<T>)
func (t *TemplatizedType) Visit(visitor Visitor) any {
	return visitor.CaseTemplatizedType(t)
}

// Java: TemplatizedType.visit(RelationshipVisitor<T>, JSType)
func (t *TemplatizedType) VisitRelationship(visitor RelationshipVisitor, that JSType) any {
	return visitor.CaseTemplatizedType(t, that)
}

// Java: TemplatizedType.toMaybeTemplatizedType()
func (t *TemplatizedType) ToMaybeTemplatizedType() *TemplatizedType {
	return t
}

// Java: TemplatizedType.getTemplateTypes()
func (t *TemplatizedType) GetTemplateTypes() []JSType {
	return t.templateTypes
}

// Java: TemplatizedType.getPropertyType(Property.Key)
func (t *TemplatizedType) GetPropertyTypeByKey(propertyName Property_Key) JSType {
	result := t.ProxyObjectTypeBase.GetPropertyTypeByKey(propertyName)
	if result == nil {
		return nil
	}
	return JSType_VisitForJSType(result, t.replacer)
}

// Java: TemplatizedType.wrapsSameRawType(JSType)
func (t *TemplatizedType) WrapsSameRawType(that JSType) bool {
	return that.IsTemplatizedType() && t.GetReferencedTypeInternal().
		Equals(that.ToMaybeTemplatizedType().GetReferencedTypeInternal())
}

// GetGreatestSubtypeHelper computes the greatest subtype of two related
// templatized types.
//
// Java: TemplatizedType.getGreatestSubtypeHelper(JSType)
func (t *TemplatizedType) GetGreatestSubtypeHelper(rawThat JSType) JSType {
	jutil.CheckNotNull(rawThat)

	if !t.WrapsSameRawType(rawThat) {
		if !rawThat.IsTemplatizedType() {
			if t.IsSubtype(rawThat) {
				return t
			} else if rawThat.IsSubtypeOf(t) {
				return rawThat
			}
		}
		if t.IsObject() && rawThat.IsObject() {
			return t.GetNativeType(JSTypeNative_NO_OBJECT_TYPE)
		}
		return t.GetNativeType(JSTypeNative_NO_TYPE)
	}

	that := rawThat.ToMaybeTemplatizedType()
	jutil.CheckNotNull(that)

	if t.Equals(that) {
		return t
	}

	// For types that have the same raw type but different type parameters, we
	// simply create a type has a "unknown" type parameter. This is equivalent
	// to the raw type.
	return t.GetReferencedObjTypeInternal()
}

// Java: TemplatizedType.getTemplateTypeMap()
func (t *TemplatizedType) GetTemplateTypeMap() *TemplateTypeMap {
	return t.templateTypeMap
}

// Java: TemplatizedType.hasAnyTemplateTypesInternal()
func (t *TemplatizedType) HasAnyTemplateTypesInternal() bool {
	return t.templateTypeMap.HasAnyTemplateTypesInternal()
}

// GetReferencedType returns the referenced ObjectType wrapped by this
// TemplatizedType.
//
// Java: TemplatizedType.getReferencedType()
func (t *TemplatizedType) GetReferencedType() ObjectType {
	return t.GetReferencedObjTypeInternal()
}

// Java: TemplatizedType.resolveInternal(ErrorReporter)
func (t *TemplatizedType) ResolveInternal(reporter ErrorReporter) JSType {
	baseTypeBefore := t.GetReferencedType()
	t.ProxyObjectTypeBase.ResolveInternal(reporter)

	rebuild := baseTypeBefore != t.GetReferencedType()
	builder := []JSType{}
	for _, type_ := range t.templateTypes {
		resolved := type_.Resolve(reporter)
		rebuild = rebuild || resolved != type_
		builder = append(builder, resolved)
	}

	if rebuild {
		return NewTemplatizedType(t.registry, t.GetReferencedType(), builder)
	} else {
		return t
	}
}
