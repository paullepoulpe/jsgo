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
// Port of com/google/javascript/rhino/jstype/TemplatizedType.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.rhino.jstype.TemplatizedType
type TemplatizedType struct {
	ProxyObjectTypeBase
	templateTypes                []JSType         `gson:"templateTypes"`
	rebuiltResolveResult         *TemplatizedType `gson:"rebuiltResolveResult"`
	isSpecializedOnlyWithUnknown bool             `gson:"isSpecializedOnlyWithUnknown"`
}

func (this_ *TemplatizedType) BindSelf_(s any) {
	_ = s
	this_.ProxyObjectTypeBase.BindSelf_(s)
}
func (this_ *TemplatizedType) IsJavaNil_() bool { return this_ == nil }
func (*TemplatizedType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.TemplatizedType", "TemplatizedType"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.TemplatizedType", "com.google.javascript.rhino.jstype.TemplatizedType", func(v_ any) bool {
		return j2gort.InstanceOf[*TemplatizedType](v_)
	})
}
func (*TemplatizedType) Implements_rhino_ProxyObjectType() {}
func (*TemplatizedType) Implements_rhino_ObjectType()      {}
func (*TemplatizedType) Implements_rhino_JSType()          {}

var static_TemplatizedType_TYPE_CLASS JSTypeClass

func StaticGet_TemplatizedType_TYPE_CLASS() JSTypeClass {
	clinit_TemplatizedType()
	return static_TemplatizedType_TYPE_CLASS
}
func StaticSet_TemplatizedType_TYPE_CLASS(v JSTypeClass) {
	clinit_TemplatizedType()
	static_TemplatizedType_TYPE_CLASS = v
}

var static_TemplatizedType_EMPTY_JSTYPE_LIST j2gort.Ext_com_google_common_collect_ImmutableList

func StaticGet_TemplatizedType_EMPTY_JSTYPE_LIST() j2gort.Ext_com_google_common_collect_ImmutableList {
	clinit_TemplatizedType()
	return static_TemplatizedType_EMPTY_JSTYPE_LIST
}
func StaticSet_TemplatizedType_EMPTY_JSTYPE_LIST(v j2gort.Ext_com_google_common_collect_ImmutableList) {
	clinit_TemplatizedType()
	static_TemplatizedType_EMPTY_JSTYPE_LIST = v
}

var static_TemplatizedType_EMPTY_OBJECT_LIST j2gort.Ext_com_google_common_collect_ImmutableList

func StaticGet_TemplatizedType_EMPTY_OBJECT_LIST() j2gort.Ext_com_google_common_collect_ImmutableList {
	clinit_TemplatizedType()
	return static_TemplatizedType_EMPTY_OBJECT_LIST
}
func StaticSet_TemplatizedType_EMPTY_OBJECT_LIST(v j2gort.Ext_com_google_common_collect_ImmutableList) {
	clinit_TemplatizedType()
	static_TemplatizedType_EMPTY_OBJECT_LIST = v
}

var clinitDone_TemplatizedType uint32
var clinitState_TemplatizedType j2gort.ClinitState

func clinit_TemplatizedType() {
	if atomic.LoadUint32(&clinitDone_TemplatizedType) != 0 {
		return
	}
	clinitSlow_TemplatizedType()
}

func clinitSlow_TemplatizedType() {
	if !j2gort.EnterClinit(&clinitState_TemplatizedType, &clinitDone_TemplatizedType) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_TemplatizedType, &clinitDone_TemplatizedType)
	static_TemplatizedType_TYPE_CLASS = StaticGet_JSTypeClass_TEMPLATIZED()
	static_TemplatizedType_EMPTY_JSTYPE_LIST = j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	static_TemplatizedType_EMPTY_OBJECT_LIST = j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
}

func (this_ *TemplatizedType) initFields_TemplatizedType() {
	_ = this_

}

// Java constructor: TemplatizedType(com.google.javascript.rhino.jstype.JSTypeRegistry registry, com.google.javascript.rhino.jstype.ObjectType objectType, com.google.common.collect.ImmutableList<com.google.javascript.rhino.jstype.JSType> templateTypes)
func NewTemplatizedType(registry *JSTypeRegistry, objectType ObjectType, templateTypes j2gort.Ext_com_google_common_collect_ImmutableList) *TemplatizedType {
	this_ := &TemplatizedType{}
	this_.BindSelf_(this_)
	InitTemplatizedType(this_, registry, objectType, templateTypes)
	return this_
}
func allocPooled_TemplatizedType(slab_ *j2gort.PooledSlab[TemplatizedType], registry *JSTypeRegistry, objectType ObjectType, templateTypes j2gort.Ext_com_google_common_collect_ImmutableList) *TemplatizedType {
	clinit_TemplatizedType()
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitTemplatizedType(this_, registry, objectType, templateTypes)
	return this_
}

func InitTemplatizedType(this_ *TemplatizedType, registry *JSTypeRegistry, objectType ObjectType, templateTypes j2gort.Ext_com_google_common_collect_ImmutableList) {
	_ = registry
	_ = objectType
	_ = templateTypes
	InitProxyObjectType_4(&this_.ProxyObjectTypeBase, registry, Upcast_JSType(objectType), objectType.GetTemplateTypeMap().CopyFilledWithValues(templateTypes), int32(4096))
	this_.initFields_TemplatizedType()
	var paramCount int32 = objectType.GetTemplateParamCount()
	_ = paramCount
	var unknownType JSType = this_.GetNativeType(StaticGet_JSTypeNative_UNKNOWN_TYPE())
	_ = unknownType
	if paramCount == int32(0) {
		this_.templateTypes = j2gort.ImmListToSliceNullable[JSType](StaticGet_TemplatizedType_EMPTY_JSTYPE_LIST())
		this_.isSpecializedOnlyWithUnknown = true
	} else if paramCount == int32(1) {
		var objMap *TemplateTypeMap = objectType.GetTemplateTypeMap()
		_ = objMap
		var key0 *TemplateType = objMap.GetTemplateKeyAt((objMap.Size() - int32(1)))
		_ = key0
		var resolved0 JSType = this_.GetTemplateTypeMap().GetResolvedTemplateType(key0)
		_ = resolved0
		this_.templateTypes = []JSType{j2gort.CheckNotNullElem(resolved0)}
		this_.isSpecializedOnlyWithUnknown = unknownType.GetJSTypeBase().EqualsJSType(resolved0)
	} else if paramCount == int32(2) {
		var objMap_2 *TemplateTypeMap = objectType.GetTemplateTypeMap()
		_ = objMap_2
		var start int32 = (objMap_2.Size() - int32(2))
		_ = start
		var ownMap *TemplateTypeMap = this_.GetTemplateTypeMap()
		_ = ownMap
		var resolved0_2 JSType = ownMap.GetResolvedTemplateType(objMap_2.GetTemplateKeyAt(start))
		_ = resolved0_2
		var resolved1 JSType = ownMap.GetResolvedTemplateType(objMap_2.GetTemplateKeyAt((start + int32(1))))
		_ = resolved1
		this_.templateTypes = []JSType{j2gort.CheckNotNullElem(resolved0_2), j2gort.CheckNotNullElem(resolved1)}
		this_.isSpecializedOnlyWithUnknown = (unknownType.GetJSTypeBase().EqualsJSType(resolved0_2) && unknownType.GetJSTypeBase().EqualsJSType(resolved1))
	} else {
		var objMap_3 *TemplateTypeMap = objectType.GetTemplateTypeMap()
		_ = objMap_3
		var ownMap_2 *TemplateTypeMap = this_.GetTemplateTypeMap()
		_ = ownMap_2
		var end int32 = objMap_3.Size()
		_ = end
		var start_2 int32 = (end - paramCount)
		_ = start_2
		var builder *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = j2gort.InitImmListBuilderCap(&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{}, paramCount)
		_ = builder
		var maybeIsSpecializedOnlyWithUnknown bool = true
		_ = maybeIsSpecializedOnlyWithUnknown
		{
			var i int32 = start_2
			_ = i
			for ; i < end; i = i + 1 {
				var resolvedType JSType = ownMap_2.GetResolvedTemplateType(objMap_3.GetTemplateKeyAt(i))
				_ = resolvedType
				builder.AddVoid_(resolvedType)
				if maybeIsSpecializedOnlyWithUnknown {
					maybeIsSpecializedOnlyWithUnknown = unknownType.GetJSTypeBase().EqualsJSType(resolvedType)
				}
			}
		}
		this_.templateTypes = j2gort.ImmListToSliceNullable[JSType](builder.BuildImmList_())
		this_.isSpecializedOnlyWithUnknown = maybeIsSpecializedOnlyWithUnknown
	}
	registry.GetResolver().ResolveIfClosed(ToIface_JSType(this_), StaticGet_TemplatizedType_TYPE_CLASS())
}

// Java method: TemplatizedType.getTypeClass()
func (this_ *TemplatizedType) GetTypeClass() JSTypeClass {
	return StaticGet_TemplatizedType_TYPE_CLASS()
}

// Java method: TemplatizedType.getCtorImplementedInterfaces()
func (this_ *TemplatizedType) GetCtorImplementedInterfaces() any {
	var resolvedImplementedInterfaces *j2gort.Ext_java_util_LinkedHashSet = nil
	_ = resolvedImplementedInterfaces
	var ownMap *TemplateTypeMap = this_.GetTemplateTypeMap()
	_ = ownMap
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.GetReferencedObjTypeInternal().GetCtorImplementedInterfaces())))
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			obj := j2gort.AsIface[ObjectType](rawElem_2)
			_ = obj
			if resolvedImplementedInterfaces == nil {
				resolvedImplementedInterfaces = j2gort.NewExt_java_util_LinkedHashSet()
			}
			resolvedImplementedInterfaces.Add_1(this_.GetJSTypeBase().registry.ReplaceTemplateTypesPartial(Upcast_JSType(obj), ownMap).GetJSTypeBase().ToObjectType())
		}
	}
	if resolvedImplementedInterfaces != nil {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](resolvedImplementedInterfaces)
	} else {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](StaticGet_TemplatizedType_EMPTY_OBJECT_LIST())
	}
}

// Java method: TemplatizedType.getCtorExtendedInterfaces()
func (this_ *TemplatizedType) GetCtorExtendedInterfaces() any {
	var resolvedExtendedInterfaces *j2gort.Ext_java_util_LinkedHashSet = nil
	_ = resolvedExtendedInterfaces
	var ownMap *TemplateTypeMap = this_.GetTemplateTypeMap()
	_ = ownMap
	{
		cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.GetReferencedObjTypeInternal().GetCtorExtendedInterfaces())))
		for cur_3.HasNext() {
			rawElem_4 := cur_3.Next()
			obj := j2gort.AsIface[ObjectType](rawElem_4)
			_ = obj
			if resolvedExtendedInterfaces == nil {
				resolvedExtendedInterfaces = j2gort.NewExt_java_util_LinkedHashSet()
			}
			resolvedExtendedInterfaces.Add_1(this_.GetJSTypeBase().registry.ReplaceTemplateTypesPartial(Upcast_JSType(obj), ownMap).GetJSTypeBase().ToObjectType())
		}
	}
	if resolvedExtendedInterfaces != nil {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](resolvedExtendedInterfaces)
	} else {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](StaticGet_TemplatizedType_EMPTY_OBJECT_LIST())
	}
}

// Java method: TemplatizedType.appendTo(com.google.javascript.rhino.jstype.TypeStringBuilder sb)
func (this_ *TemplatizedType) AppendTo(sb *TypeStringBuilder) {
	_ = sb
	this_.ProxyObjectTypeBase.AppendTo(sb)
	if len(this_.templateTypes) != 0 {
		sb.Append_1_string("<").AppendAll(j2gort.ImmListFromTypedSlice(this_.templateTypes), ",").Append_1_string(">")
	}
}

// Java method: TemplatizedType.recursionUnsafeHashCode()
func (this_ *TemplatizedType) RecursionUnsafeHashCode() int32 {
	var baseHash int32 = this_.ProxyObjectTypeBase.RecursionUnsafeHashCode()
	_ = baseHash
	if this_.isSpecializedOnlyWithUnknown {
		return baseHash
	}
	return func() int32 {
		var oh_5 int32 = 1
		oh_5 = 31*oh_5 + j2gort.ObjectHash(j2gort.ImmListFromTypedSlice(this_.templateTypes))
		oh_5 = 31*oh_5 + (baseHash)
		return oh_5
	}()
}

// Java method: TemplatizedType.visit(com.google.javascript.rhino.jstype.Visitor<T> visitor)
func (this_ *TemplatizedType) Visit_1(visitor Visitor) any {
	_ = visitor
	return visitor.CaseTemplatizedType(this_)
}

// Java method: TemplatizedType.visit(com.google.javascript.rhino.jstype.RelationshipVisitor<T> visitor, com.google.javascript.rhino.jstype.JSType that)
func (this_ *TemplatizedType) Visit_2(visitor RelationshipVisitor, that JSType) any {
	_ = visitor
	_ = that
	return visitor.CaseTemplatizedType(this_, that)
}

// Java method: TemplatizedType.toMaybeTemplatizedType()
func (this_ *TemplatizedType) ToMaybeTemplatizedType() *TemplatizedType {
	return this_
}

// Java method: TemplatizedType.getTemplateTypes()
func (this_ *TemplatizedType) GetTemplateTypes() j2gort.Ext_com_google_common_collect_ImmutableList {
	return j2gort.ImmListFromTypedSlice(this_.templateTypes)
}

// Java method: TemplatizedType.getPropertyType(com.google.javascript.rhino.jstype.Property.Key propertyName)
func (this_ *TemplatizedType) GetPropertyType_1_Property_Key(propertyName Property_Key) JSType {
	_ = propertyName
	var result JSType = this_.ProxyObjectTypeBase.GetPropertyType_1_Property_Key(propertyName)
	_ = result
	if result == nil || result.IsJavaNil_() {
		return nil
	} else {
		return this_.GetJSTypeBase().registry.ReplaceTemplateTypesPartial(result, this_.GetTemplateTypeMap())
	}
}

// Java method: TemplatizedType.wrapsSameRawType(com.google.javascript.rhino.jstype.JSType that)
func (this_ *TemplatizedType) WrapsSameRawType(that JSType) bool {
	_ = that
	return (that.GetJSTypeBase().IsTemplatizedType() && this_.GetReferencedTypeInternal().GetJSTypeBase().EqualsJSType(that.ToMaybeTemplatizedType().GetReferencedTypeInternal()))
}

// Java method: TemplatizedType.getGreatestSubtypeHelper(com.google.javascript.rhino.jstype.JSType rawThat)
func (this_ *TemplatizedType) GetGreatestSubtypeHelper(rawThat JSType) JSType {
	_ = rawThat
	j2gort.CheckNotNullIface(rawThat)
	if !(this_.WrapsSameRawType(rawThat)) {
		if this_.IsSubtype_1(rawThat) {
			return this_
		} else if rawThat.GetJSTypeBase().IsSubtypeOf_1(this_) {
			return rawThat
		}
		if this_.IsObject() && rawThat.IsObject() {
			return this_.GetNativeType(StaticGet_JSTypeNative_NO_OBJECT_TYPE())
		}
		return this_.GetNativeType(StaticGet_JSTypeNative_NO_TYPE())
	}
	var that *TemplatizedType = rawThat.ToMaybeTemplatizedType()
	_ = that
	j2gort.CheckNotNullPtr(that)
	if this_.EqualsJSType(ToIface_JSType(that)) {
		return this_
	}
	return Upcast_JSType(this_.GetReferencedObjTypeInternal())
}

// Java method: TemplatizedType.getTemplateTypeMap()
func (this_ *TemplatizedType) GetTemplateTypeMap() *TemplateTypeMap {
	return this_.GetJSTypeBase().templateTypeMap
}

// Java method: TemplatizedType.hasAnyTemplateTypesInternal()
func (this_ *TemplatizedType) HasAnyTemplateTypesInternal() bool {
	return this_.GetJSTypeBase().templateTypeMap.HasAnyTemplateTypesInternal()
}

// Java method: TemplatizedType.getReferencedType()
func (this_ *TemplatizedType) GetReferencedType() ObjectType {
	return this_.GetReferencedObjTypeInternal()
}

// Java method: TemplatizedType.getCachedResolveResult()
func (this_ *TemplatizedType) GetCachedResolveResult() JSType {
	if !(this_.IsResolved()) {
		return nil
	}
	if this_.rebuiltResolveResult != nil {
		return ToIface_JSType(this_.rebuiltResolveResult)
	} else {
		return this_
	}
}

// Java method: TemplatizedType.setCachedResolveResult(com.google.javascript.rhino.jstype.JSType result)
func (this_ *TemplatizedType) SetCachedResolveResult(result JSType) {
	_ = result
	j2gort.CheckNotNullIface(result)
	if result == this_ || ((result == nil || result.IsJavaNil_()) && (this_ == nil || this_.IsJavaNil_())) {
		this_.rebuiltResolveResult = nil
	} else {
		this_.rebuiltResolveResult = j2gort.CastPtr[TemplatizedType](result)
	}
	this_.MarkResolved()
}

// Java method: TemplatizedType.clearCachedResolveResult()
func (this_ *TemplatizedType) ClearCachedResolveResult() {
	this_.rebuiltResolveResult = nil
	this_.ProxyObjectTypeBase.ClearCachedResolveResult()
}

// Java method: TemplatizedType.resolveInternal(com.google.javascript.rhino.ErrorReporter reporter)
func (this_ *TemplatizedType) ResolveInternal(reporter ErrorReporter) JSType {
	_ = reporter
	var baseTypeBefore JSType = Upcast_JSType(this_.GetReferencedType())
	_ = baseTypeBefore
	this_.ProxyObjectTypeBase.ResolveInternal(reporter)
	var rebuild bool = !(j2gort.IfaceEquals[JSType](baseTypeBefore, this_.GetReferencedType()))
	_ = rebuild
	var n int32 = int32(len(this_.templateTypes))
	_ = n
	var builder *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = nil
	_ = builder
	if rebuild && (n > int32(0)) {
		builder = j2gort.InitImmListBuilderCap(&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{}, n)
	}
	{
		var i int32 = int32(0)
		_ = i
		for ; i < n; i = i + 1 {
			var type_ JSType = this_.templateTypes[i]
			_ = type_
			var resolved JSType = type_.GetJSTypeBase().Resolve(reporter)
			_ = resolved
			if builder != nil {
				builder.AddVoid_(resolved)
			} else if !(resolved == type_ || ((resolved == nil || resolved.IsJavaNil_()) && (type_ == nil || type_.IsJavaNil_()))) {
				rebuild = true
				builder = j2gort.InitImmListBuilderCap(&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{}, n)
				{
					var j int32 = int32(0)
					_ = j
					for ; j < i; j = j + 1 {
						builder.AddVoid_(this_.templateTypes[j])
					}
				}
				builder.AddVoid_(resolved)
			}
		}
	}
	if rebuild {
		return ToIface_JSType(allocPooled_TemplatizedType(&this_.registry.slab_TemplatizedType, this_.GetJSTypeBase().registry, this_.GetReferencedType(), func() j2gort.Ext_com_google_common_collect_ImmutableList {
			if builder != nil {
				return builder.BuildImmList_()
			}
			return StaticGet_TemplatizedType_EMPTY_JSTYPE_LIST()
		}()))
	} else {
		return this_
	}
}
func (this_ *TemplatizedType) Slice_GetTemplateTypes() []JSType {
	return this_.templateTypes
}
