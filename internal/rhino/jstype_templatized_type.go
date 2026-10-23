// Port of com/google/javascript/rhino/jstype/TemplatizedType.java @ closure-compiler v20261015

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
}

// Java: TemplatizedType.TYPE_CLASS
const TemplatizedType_TYPE_CLASS = JSTypeClass_TEMPLATIZED

// Java: TemplatizedType(JSTypeRegistry, ObjectType, ImmutableList<JSType>)
func NewTemplatizedType(registry *JSTypeRegistry, objectType ObjectType, templateTypes []JSType) *TemplatizedType {
	t := &TemplatizedType{}
	t.objFlags = objKindTemplatized
	ttm := objectType.GetTemplateTypeMap().CopyFilledWithValues(templateTypes)
	t.InitProxyObjectTypeBaseWithTemplateTypeMap(t, registry, objectType, ttm)

	typeParams := objectType.GetTypeParameters()
	nParams := len(typeParams)
	if nParams == 0 {
		t.templateTypes = emptyJSTypeSlice
		t.isSpecializedOnlyWithUnknown = true
	} else if nParams == 1 && len(ttm.resolvedTemplateValues) == len(ttm.templateKeys) {
		nKeys := len(ttm.templateKeys)
		resolvedType := jutil.CheckNotNull(ttm.resolvedTemplateValues[nKeys-1])
		t.templateTypes = ttm.resolvedTemplateValues[nKeys-1 : nKeys]
		t.isSpecializedOnlyWithUnknown = jsTypeEquals(t.GetNativeType(JSTypeNative_UNKNOWN_TYPE), resolvedType)
	} else if nParams == 2 && len(ttm.resolvedTemplateValues) == len(ttm.templateKeys) && typeParams[0] != typeParams[1] {
		nKeys := len(ttm.templateKeys)
		r0 := jutil.CheckNotNull(ttm.resolvedTemplateValues[nKeys-2])
		r1 := jutil.CheckNotNull(ttm.resolvedTemplateValues[nKeys-1])
		t.templateTypes = ttm.resolvedTemplateValues[nKeys-2 : nKeys]
		unk := t.GetNativeType(JSTypeNative_UNKNOWN_TYPE)
		t.isSpecializedOnlyWithUnknown = jsTypeEquals(unk, r0) && jsTypeEquals(unk, r1)
	} else {
		builder := make([]JSType, 0, nParams)
		maybeIsSpecializedOnlyWithUnknown := true
		for _, newlyFilledTemplateKey := range typeParams {
			resolvedType := ttm.GetResolvedTemplateType(newlyFilledTemplateKey)

			builder = append(builder, jutil.CheckNotNull(resolvedType))
			if maybeIsSpecializedOnlyWithUnknown {
				maybeIsSpecializedOnlyWithUnknown = jsTypeEquals(t.GetNativeType(JSTypeNative_UNKNOWN_TYPE), resolvedType)
			}
		}
		t.templateTypes = builder
		t.isSpecializedOnlyWithUnknown = maybeIsSpecializedOnlyWithUnknown
	}

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
	replacer := MakeTemplateTypeReplacerForPartialReplacement(t.GetRegistry(), t.templateTypeMap)
	for _, obj := range t.GetReferencedObjTypeInternal().GetCtorImplementedInterfaces() {
		resolvedImplementedInterfaces.Add(replacer.VisitJSType(obj).ToObjectType())
	}
	return resolvedImplementedInterfaces.ToSlice()
}

// Java: TemplatizedType.getCtorExtendedInterfaces()
func (t *TemplatizedType) GetCtorExtendedInterfaces() []ObjectType {
	resolvedExtendedInterfaces := jutil.NewObjectLinkedHashSet[ObjectType]()
	replacer := MakeTemplateTypeReplacerForPartialReplacement(t.GetRegistry(), t.templateTypeMap)
	for _, obj := range t.GetReferencedObjTypeInternal().GetCtorExtendedInterfaces() {
		resolvedExtendedInterfaces.Add(replacer.VisitJSType(obj).ToObjectType())
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

// Java: JSType.isTemplatizedType()
func (t *TemplatizedType) IsTemplatizedType() bool {
	return true
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
	replacer := MakeTemplateTypeReplacerForPartialReplacement(t.GetRegistry(), t.templateTypeMap)
	return replacer.VisitJSType(result)
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
		if t.IsSubtype(rawThat) {
			return t
		} else if rawThat.IsSubtypeOf(t) {
			return rawThat
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
	var builderBuf [4]JSType
	var builder []JSType
	for i, type_ := range t.templateTypes {
		resolved := type_.Resolve(reporter)
		if builder != nil {
			builder[i] = resolved
		} else if resolved != type_ {
			rebuild = true
			if len(t.templateTypes) <= len(builderBuf) {
				builder = builderBuf[:len(t.templateTypes)]
			} else {
				builder = make([]JSType, len(t.templateTypes))
			}
			copy(builder[:i], t.templateTypes[:i])
			builder[i] = resolved
		}
	}

	if rebuild {
		if builder == nil {
			builder = t.templateTypes
		}
		return t.GetRegistry().CreateTemplatizedTypeOfList(t.GetReferencedType(), builder)
	} else {
		return t
	}
}
