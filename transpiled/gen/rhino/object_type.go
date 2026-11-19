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
// Port of com/google/javascript/rhino/jstype/ObjectType.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java polymorphic class interface: com.google.javascript.rhino.jstype.ObjectType
type ObjectType interface {
	Implements_rhino_ObjectType()
	IsJavaNil_() bool
	GetObjectTypeBase() *ObjectTypeBase
	JSType
	GetPropertyMap() *PropertyMap
	GetSlot_1_Property_Key(Property_Key) *Property
	GetSlot_1_string(string) *Property
	GetOwnSlot_1_string(string) *Property
	GetOwnSlot_1_Property_Key(Property_Key) *Property
	GetTypeOfThis() JSType
	GetTemplateTypes() j2gort.Ext_com_google_common_collect_ImmutableList
	SetJSDocInfo(*JSDocInfo)
	DetectImplicitPrototypeCycle() bool
	DetectInheritanceCycle() bool
	GetReferenceName() string
	HasReferenceName() bool
	GetNormalizedReferenceName() string
	GetRawType() ObjectType
	GetConstructor() *FunctionType
	GetSuperClassConstructor() *FunctionType
	GetClosestDefiningType(string) ObjectType
	FindClosestDefinition_1_string(string) *Property_OwnedProperty
	FindClosestDefinition_1_Property_Key(Property_Key) *Property_OwnedProperty
	GetImplicitPrototype() ObjectType
	GetImplicitPrototypeChain() j2gort.Ext_java_lang_Iterable
	DefineDeclaredProperty_3_string_JSType_Node(string, JSType, Node) bool
	DefineDeclaredProperty_3_Property_Key_JSType_Node(Property_Key, JSType, Node) bool
	DefineSynthesizedProperty(string, JSType, Node) bool
	DefineInferredProperty_3_Property_Key_JSType_Node(Property_Key, JSType, Node) bool
	DefineInferredProperty_3_string_JSType_Node(string, JSType, Node) bool
	DefineProperty_4_string_JSType_bool_Node(string, JSType, bool, Node) bool
	DefineProperty_4_Property_Key_JSType_bool_Node(Property_Key, JSType, bool, Node) bool
	GetPropertyNode_1_string(string) Node
	GetPropertyNode_1_Property_Key(Property_Key) Node
	GetPropertyDefSite_1_string(string) Node
	GetPropertyDefSite_1_Property_Key(Property_Key) Node
	GetPropertyJSDocInfo_1_string(string) *JSDocInfo
	GetPropertyJSDocInfo_1_Property_Key(Property_Key) *JSDocInfo
	GetOwnPropertyJSDocInfo_1_string(string) *JSDocInfo
	GetOwnPropertyJSDocInfo_1_Property_Key(Property_Key) *JSDocInfo
	GetOwnPropertyDefSite_1_string(string) Node
	GetOwnPropertyDefSite_1_Property_Key(Property_Key) Node
	SetPropertyJSDocInfo_2_string_JSDocInfo(string, *JSDocInfo)
	SetPropertyJSDocInfo_2_Property_Key_JSDocInfo(Property_Key, *JSDocInfo)
	SetPropertyNode_2_string_Node(string, Node)
	SetPropertyNode_2_Property_Key_Node(Property_Key, Node)
	GetPropertyType_1_string(string) JSType
	GetPropertyType_1_Property_Key(Property_Key) JSType
	GetOwnPropertyKind_1_string(string) JSType_HasPropertyKind
	GetOwnPropertyKind_1_Property_Key(Property_Key) JSType_HasPropertyKind
	HasOwnProperty_1_string(string) bool
	HasOwnProperty_1_Property_Key(Property_Key) bool
	GetOwnPropertyKnownSymbols() j2gort.Ext_java_util_Set
	GetOwnPropertyNames() j2gort.Ext_java_util_Set
	GetOwnPropertyNamesUpperBound() int32
	CollectOwnPropertyNames([]string) int32
	GetOwnPropertyKeys() j2gort.Ext_com_google_common_collect_ImmutableSet
	IsPropertyTypeInferred_1_string(string) bool
	IsPropertyTypeInferred_1_Property_Key(Property_Key) bool
	IsPropertyTypeDeclared_1_string(string) bool
	IsPropertyTypeDeclared_1_Property_Key(Property_Key) bool
	HasOwnDeclaredProperty_1_string(string) bool
	HasOwnDeclaredProperty_1_Property_Key(Property_Key) bool
	IsPropertyInExterns(string) bool
	GetPropertiesCount() int32
	GetPropertyNames() j2gort.Ext_com_google_common_collect_ImmutableSortedSet
	GetAllKeys() *PropertyMap_AllKeys
	IsImplicitPrototypeOf(ObjectType) bool
	HasCachedValues() bool
	ClearCachedValues()
	GetOwnerFunction() *FunctionType
	SetOwnerFunction(*FunctionType)
	GetCtorImplementedInterfaces() any
	GetCtorExtendedInterfaces() any
	GetPropertyTypeMap() j2gort.Ext_java_util_Map
	GetEnumeratedTypeOfEnumObject() JSType
}

func ToIface_ObjectType[T interface {
	ObjectType
	comparable
}](v T) ObjectType {
	var z T
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.rhino.jstype.ObjectType
type ObjectTypeBase struct {
	JSTypeBase
	docInfo *JSDocInfo `gson:"docInfo"`
}

func (this_ *ObjectTypeBase) BindSelf_(s any) {
	_ = s
	this_.JSTypeBase.BindSelf_(s)
}
func (this_ *ObjectTypeBase) IsJavaNil_() bool { return this_ == nil }
func (*ObjectTypeBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.ObjectType", "ObjectType"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.ObjectType", "com.google.javascript.rhino.jstype.ObjectType", func(v_ any) bool {
		return j2gort.InstanceOf[ObjectType](v_)
	})
}
func (b *ObjectTypeBase) GetObjectTypeBase() *ObjectTypeBase { return b }
func (this_ *ObjectTypeBase) self_ObjectType_() ObjectType {
	if s := this_.JSTypeBase.self_; s != nil {
		return s.(ObjectType)
	}
	return nil
}
func (*ObjectTypeBase) Implements_rhino_ObjectType() {}
func (*ObjectTypeBase) Implements_rhino_JSType()     {}

var static_ObjectType_EMPTY_STRING_SET j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_ObjectType_EMPTY_STRING_SET() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_ObjectType()
	return static_ObjectType_EMPTY_STRING_SET
}
func StaticSet_ObjectType_EMPTY_STRING_SET(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_ObjectType()
	static_ObjectType_EMPTY_STRING_SET = v
}

var static_ObjectType_PROTOTYPE_KEY *Property_StringKey

func StaticGet_ObjectType_PROTOTYPE_KEY() *Property_StringKey {
	clinit_ObjectType()
	return static_ObjectType_PROTOTYPE_KEY
}
func StaticSet_ObjectType_PROTOTYPE_KEY(v *Property_StringKey) {
	clinit_ObjectType()
	static_ObjectType_PROTOTYPE_KEY = v
}

var static_ObjectType_CALL_KEY *Property_StringKey

func StaticGet_ObjectType_CALL_KEY() *Property_StringKey {
	clinit_ObjectType()
	return static_ObjectType_CALL_KEY
}
func StaticSet_ObjectType_CALL_KEY(v *Property_StringKey) {
	clinit_ObjectType()
	static_ObjectType_CALL_KEY = v
}

var static_ObjectType_BIND_KEY *Property_StringKey

func StaticGet_ObjectType_BIND_KEY() *Property_StringKey {
	clinit_ObjectType()
	return static_ObjectType_BIND_KEY
}
func StaticSet_ObjectType_BIND_KEY(v *Property_StringKey) {
	clinit_ObjectType()
	static_ObjectType_BIND_KEY = v
}

var static_ObjectType_APPLY_KEY *Property_StringKey

func StaticGet_ObjectType_APPLY_KEY() *Property_StringKey {
	clinit_ObjectType()
	return static_ObjectType_APPLY_KEY
}
func StaticSet_ObjectType_APPLY_KEY(v *Property_StringKey) {
	clinit_ObjectType()
	static_ObjectType_APPLY_KEY = v
}

var clinitDone_ObjectType uint32
var clinitState_ObjectType j2gort.ClinitState

func clinit_ObjectType() {
	if atomic.LoadUint32(&clinitDone_ObjectType) != 0 {
		return
	}
	clinitSlow_ObjectType()
}

func clinitSlow_ObjectType() {
	if !j2gort.EnterClinit(&clinitState_ObjectType, &clinitDone_ObjectType) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ObjectType, &clinitDone_ObjectType)
	static_ObjectType_EMPTY_STRING_SET = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()
	static_ObjectType_PROTOTYPE_KEY = NewProperty_StringKey("prototype")
	static_ObjectType_CALL_KEY = NewProperty_StringKey("call")
	static_ObjectType_BIND_KEY = NewProperty_StringKey("bind")
	static_ObjectType_APPLY_KEY = NewProperty_StringKey("apply")
}

func (this_ *ObjectTypeBase) initFields_ObjectType() {
	_ = this_
	this_.docInfo = nil
}

func InitObjectType_1(this_ *ObjectTypeBase, registry *JSTypeRegistry) {
	_ = registry
	InitJSType_1(&this_.JSTypeBase, registry)
	this_.initFields_ObjectType()
}

func InitObjectType_2_JSTypeRegistry_int32(this_ *ObjectTypeBase, registry *JSTypeRegistry, typeKindMask int32) {
	_ = registry
	_ = typeKindMask
	InitJSType_3(&this_.JSTypeBase, registry, nil, typeKindMask)
	this_.initFields_ObjectType()
}

func InitObjectType_2_JSTypeRegistry_TemplateTypeMap(this_ *ObjectTypeBase, registry *JSTypeRegistry, templateTypeMap *TemplateTypeMap) {
	_ = registry
	_ = templateTypeMap
	InitJSType_2(&this_.JSTypeBase, registry, templateTypeMap)
	this_.initFields_ObjectType()
}

func InitObjectType_3(this_ *ObjectTypeBase, registry *JSTypeRegistry, templateTypeMap *TemplateTypeMap, typeKindMask int32) {
	_ = registry
	_ = templateTypeMap
	_ = typeKindMask
	InitJSType_3(&this_.JSTypeBase, registry, templateTypeMap, typeKindMask)
	this_.initFields_ObjectType()
}

// Java method: ObjectType.getPropertyMap()
func (this_ *ObjectTypeBase) GetPropertyMap() *PropertyMap {
	return PropertyMap_ImmutableEmptyMap()
}

// Java static method: ObjectType.unwrapToPrototypeObjectType(com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable ObjectType type)
func ObjectType_UnwrapToPrototypeObjectType(type_ ObjectType) PrototypeObjectType {
	clinit_ObjectType()
	_ = type_
	var curr ObjectType = type_
	_ = curr
	for curr != nil && !curr.IsJavaNil_() {
		var pot PrototypeObjectType
		_ = pot
		var ok_1 bool
		pot, ok_1 = j2gort.TryCast[PrototypeObjectType](curr)
		if ok_1 {
			return pot
		}
		var proxy ProxyObjectType
		_ = proxy
		var enumElem *EnumElementType
		_ = enumElem
		var ok_2 bool
		proxy, ok_2 = j2gort.TryCast[ProxyObjectType](curr)
		if ok_2 {
			curr = proxy.GetReferencedObjTypeInternal()
		} else {
			var ok_3 bool
			enumElem, ok_3 = j2gort.TryCast[*EnumElementType](curr)
			if ok_3 {
				curr = enumElem.GetPrimitiveType().ToObjectType()
			} else {
				return nil
			}
		}
	}
	return nil
}

// Java method: ObjectType.getSlot(com.google.javascript.rhino.jstype.Property.Key name)
func (this_ *ObjectTypeBase) GetSlot_1_Property_Key(name Property_Key) *Property {
	_ = name
	return PropertyMap_FindClosestProperty_2_PrototypeObjectType_Property_Key(ObjectType_UnwrapToPrototypeObjectType(this_.self_ObjectType_()), name)
}

// Java method: ObjectType.getSlot(java.lang.String name)
func (this_ *ObjectTypeBase) GetSlot_1_string(name string) *Property {
	_ = name
	if j2gort.InstanceOf[*FunctionType](this_.self_ObjectType_()) && ("prototype" == name) {
		return this_.self_ObjectType_().GetSlot_1_Property_Key(ToIface_Property_Key(StaticGet_ObjectType_PROTOTYPE_KEY()))
	}
	return PropertyMap_FindClosestProperty_2_PrototypeObjectType_string(ObjectType_UnwrapToPrototypeObjectType(this_.self_ObjectType_()), name)
}

// Java method: ObjectType.getOwnSlot(java.lang.String name)
func (this_ *ObjectTypeBase) GetOwnSlot_1_string(name string) *Property {
	_ = name
	var src PrototypeObjectType = ObjectType_UnwrapToPrototypeObjectType(this_.self_ObjectType_())
	_ = src
	if src == nil || src.IsJavaNil_() {
		return nil
	}
	var map_ *PropertyMap = src.GetOwnPropertyMapOrNull()
	_ = map_
	if map_ != nil {
		return map_.GetOwnProperty_1_string(name)
	} else {
		return nil
	}
}

// Java method: ObjectType.getOwnSlot(com.google.javascript.rhino.jstype.Property.Key name)
func (this_ *ObjectTypeBase) GetOwnSlot_1_Property_Key(name Property_Key) *Property {
	_ = name
	var src PrototypeObjectType = ObjectType_UnwrapToPrototypeObjectType(this_.self_ObjectType_())
	_ = src
	if src == nil || src.IsJavaNil_() {
		return nil
	}
	var map_ *PropertyMap = src.GetOwnPropertyMapOrNull()
	_ = map_
	if map_ != nil {
		return map_.GetOwnProperty_1_Property_Key(name)
	} else {
		return nil
	}
}

// Java method: ObjectType.getTypeOfThis()
func (this_ *ObjectTypeBase) GetTypeOfThis() JSType {
	return nil
}

// Java method: ObjectType.getTemplateTypes()
func (this_ *ObjectTypeBase) GetTemplateTypes() j2gort.Ext_com_google_common_collect_ImmutableList {
	return j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
}

// Java method: ObjectType.getJSDocInfo()
func (this_ *ObjectTypeBase) GetJSDocInfo() *JSDocInfo {
	return this_.docInfo
}

// Java method: ObjectType.setJSDocInfo(com.google.javascript.rhino.JSDocInfo info)
func (this_ *ObjectTypeBase) SetJSDocInfo(info *JSDocInfo) {
	_ = info
	this_.docInfo = info
}

// Java method: ObjectType.detectImplicitPrototypeCycle()
func (this_ *ObjectTypeBase) DetectImplicitPrototypeCycle() bool {
	this_.GetJSTypeBase().jstypeFlags = int8((int32(this_.GetJSTypeBase().jstypeFlags) | int32(int8(16))))
	var p ObjectType = this_.self_ObjectType_().GetImplicitPrototype()
	_ = p
	for p != nil && !p.IsJavaNil_() {
		if (int32(p.GetJSTypeBase().jstypeFlags) & int32(int8(16))) != int32(0) {
			return true
		} else {
			r_4 := p
			r_4.GetJSTypeBase().jstypeFlags = int8((int32(r_4.GetJSTypeBase().jstypeFlags) | int32(int8(16))))
		}
		p = p.GetImplicitPrototype()
	}
	p = this_.self_ObjectType_()
	for doFirst_5 := true; doFirst_5 || (p != nil && !p.IsJavaNil_()); doFirst_5 = false {
		r_6 := p
		r_6.GetJSTypeBase().jstypeFlags = int8((int32(r_6.GetJSTypeBase().jstypeFlags) & ^(int32(int8(16)))))
		p = p.GetImplicitPrototype()
	}
	return false
}

// Java method: ObjectType.detectInheritanceCycle()
func (this_ *ObjectTypeBase) DetectInheritanceCycle() bool {
	if this_.DetectImplicitPrototypeCycle() || j2gort.Ext_com_google_common_collect_Iterables_Contains(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.self_ObjectType_().GetCtorImplementedInterfaces()), this_.JSTypeBase.self_) {
		return true
	}
	var fnType *FunctionType = this_.self_ObjectType_().GetConstructor()
	_ = fnType
	return ((fnType != nil) && (fnType.CheckExtendsLoop_0() != nil))
}

// Java method: ObjectType.getReferenceName()
func (this_ *ObjectTypeBase) GetReferenceName() string {
	panic("abstract or native method invoked")
}

// Java method: ObjectType.hasReferenceName()
func (this_ *ObjectTypeBase) HasReferenceName() bool {
	return (this_.self_ObjectType_().GetReferenceName() != j2gort.NullString)
}

// Java method: ObjectType.getNormalizedReferenceName()
func (this_ *ObjectTypeBase) GetNormalizedReferenceName() string {
	var name string = this_.self_ObjectType_().GetReferenceName()
	_ = name
	if name != j2gort.NullString {
		var start int32 = j2gort.JStrIndexOfCp(name, int32(uint16(40)))
		_ = start
		if start != int32(-1) {
			var end int32 = j2gort.JStrLastIndexOfCp(name, int32(uint16(41)))
			_ = end
			var prefix string = j2gort.JStrSub(name, int32(0), start)
			_ = prefix
			if (end + (int32(1) % j2gort.JStrLen(name))) == int32(0) {
				return prefix
			} else {
				return (j2gort.StrValString(prefix) + j2gort.StrValString(j2gort.JStrSubFrom(name, (end+int32(1)))))
			}
		}
	}
	return name
}

// Java method: ObjectType.getDisplayName()
func (this_ *ObjectTypeBase) GetDisplayName() string {
	return this_.GetNormalizedReferenceName()
}

// Java static method: ObjectType.createDelegateSuffix(java.lang.String suffix)
func ObjectType_CreateDelegateSuffix(suffix string) string {
	clinit_ObjectType()
	_ = suffix
	return ("(" + j2gort.StrValString(suffix) + ")")
}

// Java method: ObjectType.getRawType()
func (this_ *ObjectTypeBase) GetRawType() ObjectType {
	var t *TemplatizedType = this_.JSTypeBase.self_.ToMaybeTemplatizedType()
	_ = t
	if t == nil {
		return this_.self_ObjectType_()
	} else {
		return t.GetReferencedType()
	}
}

// Java method: ObjectType.testForEquality(com.google.javascript.rhino.jstype.JSType that)
func (this_ *ObjectTypeBase) TestForEquality(that JSType) pkg_base.Tri {
	_ = that
	var result pkg_base.Tri = this_.JSTypeBase.TestForEquality(that)
	_ = result
	if result != 0 {
		return result
	}
	if (((((that.IsUnknownType() || that.IsSubtypeOf_1(this_.GetNativeType(StaticGet_JSTypeNative_OBJECT_TYPE()))) || that.IsSubtypeOf_1(this_.GetNativeType(StaticGet_JSTypeNative_NUMBER_TYPE()))) || that.IsSubtypeOf_1(this_.GetNativeType(StaticGet_JSTypeNative_STRING_TYPE()))) || that.IsSubtypeOf_1(this_.GetNativeType(StaticGet_JSTypeNative_BOOLEAN_TYPE()))) || that.IsSubtypeOf_1(this_.GetNativeType(StaticGet_JSTypeNative_SYMBOL_TYPE()))) || that.IsSubtypeOf_1(this_.GetNativeType(StaticGet_JSTypeNative_BIGINT_TYPE())) {
		return pkg_base.StaticGet_Tri_UNKNOWN()
	}
	return pkg_base.StaticGet_Tri_FALSE()
}

// Java method: ObjectType.getConstructor()
func (this_ *ObjectTypeBase) GetConstructor() *FunctionType {
	panic("abstract or native method invoked")
}

// Java method: ObjectType.getSuperClassConstructor()
func (this_ *ObjectTypeBase) GetSuperClassConstructor() *FunctionType {
	var iproto ObjectType = this_.self_ObjectType_().GetImplicitPrototype()
	_ = iproto
	if iproto == nil || iproto.IsJavaNil_() {
		return nil
	}
	iproto = iproto.GetImplicitPrototype()
	if iproto == nil || iproto.IsJavaNil_() {
		return nil
	} else {
		return iproto.GetConstructor()
	}
}

// Java method: ObjectType.getClosestDefiningType(java.lang.String propertyName)
func (this_ *ObjectTypeBase) GetClosestDefiningType(propertyName string) ObjectType {
	_ = propertyName
	var property *Property_OwnedProperty = PropertyMap_FindClosest_2_PrototypeObjectType_string(ObjectType_UnwrapToPrototypeObjectType(this_.self_ObjectType_()), propertyName)
	_ = property
	if property == nil {
		return nil
	} else {
		return property.GetOwner()
	}
}

// Java method: ObjectType.findClosestDefinition(java.lang.String propertyName)
func (this_ *ObjectTypeBase) FindClosestDefinition_1_string(propertyName string) *Property_OwnedProperty {
	_ = propertyName
	return PropertyMap_FindClosest_2_PrototypeObjectType_string(ObjectType_UnwrapToPrototypeObjectType(this_.self_ObjectType_()), propertyName)
}

// Java method: ObjectType.findClosestDefinition(com.google.javascript.rhino.jstype.Property.Key propertyName)
func (this_ *ObjectTypeBase) FindClosestDefinition_1_Property_Key(propertyName Property_Key) *Property_OwnedProperty {
	_ = propertyName
	return PropertyMap_FindClosest_2_PrototypeObjectType_Property_Key(ObjectType_UnwrapToPrototypeObjectType(this_.self_ObjectType_()), propertyName)
}

// Java method: ObjectType.getImplicitPrototype()
func (this_ *ObjectTypeBase) GetImplicitPrototype() ObjectType {
	panic("abstract or native method invoked")
}

// Java method: ObjectType.getImplicitPrototypeChain()
func (this_ *ObjectTypeBase) GetImplicitPrototypeChain() j2gort.Ext_java_lang_Iterable {
	var self ObjectType = this_.self_ObjectType_()
	_ = self
	return j2gort.NewLambda_Ext_java_lang_Iterable(func() any {
		return NewObjectType__Anon_159(this_.self_ObjectType_(), self)
	})
}

// Java method: ObjectType.defineDeclaredProperty(java.lang.String propertyName, com.google.javascript.rhino.jstype.JSType type, com.google.javascript.rhino.Node propertyNode)
func (this_ *ObjectTypeBase) DefineDeclaredProperty_3_string_JSType_Node(propertyName string, type_ JSType, propertyNode Node) bool {
	_ = propertyName
	_ = type_
	_ = propertyNode
	j2gort.CheckNotNullString(propertyName)
	var result bool = this_.self_ObjectType_().DefineProperty_4_string_JSType_bool_Node(propertyName, type_, false, propertyNode)
	_ = result
	this_.GetJSTypeBase().registry.RegisterPropertyOnType(propertyName, this_.JSTypeBase.self_)
	return result
}

// Java method: ObjectType.defineDeclaredProperty(com.google.javascript.rhino.jstype.Property.Key propertyName, com.google.javascript.rhino.jstype.JSType type, com.google.javascript.rhino.Node propertyNode)
func (this_ *ObjectTypeBase) DefineDeclaredProperty_3_Property_Key_JSType_Node(propertyName Property_Key, type_ JSType, propertyNode Node) bool {
	_ = propertyName
	_ = type_
	_ = propertyNode
	var result bool = this_.self_ObjectType_().DefineProperty_4_Property_Key_JSType_bool_Node(propertyName, type_, false, propertyNode)
	_ = result
	{
		var iifeCond_1 bool
		{

			re_7 := propertyName.Kind()
			ae_8 := StaticGet_Property_KeyKind_STRING()
			if re_7 == 0 {
				j2gort.ThrowNPE()
			}
			iifeCond_1 = re_7 == ae_8
		}
		if iifeCond_1 {
			this_.GetJSTypeBase().registry.RegisterPropertyOnType(propertyName.M_String(), this_.JSTypeBase.self_)
		}
	}
	return result
}

// Java method: ObjectType.defineSynthesizedProperty(java.lang.String propertyName, com.google.javascript.rhino.jstype.JSType type, com.google.javascript.rhino.Node propertyNode)
func (this_ *ObjectTypeBase) DefineSynthesizedProperty(propertyName string, type_ JSType, propertyNode Node) bool {
	_ = propertyName
	_ = type_
	_ = propertyNode
	return this_.self_ObjectType_().DefineProperty_4_string_JSType_bool_Node(propertyName, type_, false, propertyNode)
}

// Java method: ObjectType.defineInferredProperty(com.google.javascript.rhino.jstype.Property.Key propertyName, com.google.javascript.rhino.jstype.JSType type, com.google.javascript.rhino.Node propertyNode)
func (this_ *ObjectTypeBase) DefineInferredProperty_3_Property_Key_JSType_Node(propertyName Property_Key, type_ JSType, propertyNode Node) bool {
	_ = propertyName
	_ = type_
	_ = propertyNode
	if !(propertyName.M_String() != j2gort.NullString) {
		j2gort.CheckArgumentStr(false, "Symbol-based property keys are not supported for inferred properties.")
	}
	return this_.DefineInferredProperty_3_string_JSType_Node(propertyName.M_String(), type_, propertyNode)
}

// Java method: ObjectType.defineInferredProperty(java.lang.String propertyName, com.google.javascript.rhino.jstype.JSType type, com.google.javascript.rhino.Node propertyNode)
func (this_ *ObjectTypeBase) DefineInferredProperty_3_string_JSType_Node(propertyName string, type_ JSType, propertyNode Node) bool {
	_ = propertyName
	_ = type_
	_ = propertyNode
	if j2gort.InstanceOf[ProxyObjectType](this_.self_ObjectType_()) || j2gort.InstanceOf[*EnumElementType](this_.self_ObjectType_()) {
		if this_.HasProperty_1_string(propertyName) {
			if this_.IsPropertyTypeDeclared_1_string(propertyName) {
				return true
			}
			var originalType JSType = j2gort.CheckNotNullIface(this_.GetPropertyType_1_string(propertyName))
			_ = originalType
			type_ = originalType.GetLeastSupertype(type_)
		}
	} else {
		var slot *Property = this_.GetSlot_1_string(propertyName)
		_ = slot
		if ((slot != nil) || this_.IsEmptyType()) || this_.self_ObjectType_().IsUnknownType() {
			if (slot != nil) && !(slot.IsTypeInferred()) {
				return true
			}
			var originalType_2 JSType = j2gort.CheckNotNullIface(this_.GetPropertyType_1_string(propertyName))
			_ = originalType_2
			type_ = originalType_2.GetLeastSupertype(type_)
		}
	}
	var result bool = this_.self_ObjectType_().DefineProperty_4_string_JSType_bool_Node(propertyName, type_, true, propertyNode)
	_ = result
	this_.GetJSTypeBase().registry.RegisterPropertyOnType(propertyName, this_.JSTypeBase.self_)
	return result
}

// Java method: ObjectType.defineProperty(java.lang.String propertyName, com.google.javascript.rhino.jstype.JSType type, boolean inferred, com.google.javascript.rhino.Node propertyNode)
func (this_ *ObjectTypeBase) DefineProperty_4_string_JSType_bool_Node(propertyName string, type_ JSType, inferred bool, propertyNode Node) bool {
	_ = propertyName
	_ = type_
	_ = inferred
	_ = propertyNode
	panic("abstract or native method invoked")
}

// Java method: ObjectType.defineProperty(com.google.javascript.rhino.jstype.Property.Key propertyName, com.google.javascript.rhino.jstype.JSType type, boolean inferred, com.google.javascript.rhino.Node propertyNode)
func (this_ *ObjectTypeBase) DefineProperty_4_Property_Key_JSType_bool_Node(propertyName Property_Key, type_ JSType, inferred bool, propertyNode Node) bool {
	_ = propertyName
	_ = type_
	_ = inferred
	_ = propertyNode
	panic("abstract or native method invoked")
}

// Java method: ObjectType.getPropertyNode(java.lang.String propertyName)
func (this_ *ObjectTypeBase) GetPropertyNode_1_string(propertyName string) Node {
	_ = propertyName
	var p *Property = this_.GetSlot_1_string(propertyName)
	_ = p
	if p == nil {
		return Node(0)
	} else {
		return p.GetNode()
	}
}

// Java method: ObjectType.getPropertyNode(com.google.javascript.rhino.jstype.Property.Key propertyName)
func (this_ *ObjectTypeBase) GetPropertyNode_1_Property_Key(propertyName Property_Key) Node {
	_ = propertyName
	var p *Property = this_.self_ObjectType_().GetSlot_1_Property_Key(propertyName)
	_ = p
	if p == nil {
		return Node(0)
	} else {
		return p.GetNode()
	}
}

// Java method: ObjectType.getPropertyDefSite(java.lang.String propertyName)
func (this_ *ObjectTypeBase) GetPropertyDefSite_1_string(propertyName string) Node {
	_ = propertyName
	return this_.GetPropertyNode_1_string(propertyName)
}

// Java method: ObjectType.getPropertyDefSite(com.google.javascript.rhino.jstype.Property.Key propertyName)
func (this_ *ObjectTypeBase) GetPropertyDefSite_1_Property_Key(propertyName Property_Key) Node {
	_ = propertyName
	return this_.GetPropertyNode_1_Property_Key(propertyName)
}

// Java method: ObjectType.getPropertyJSDocInfo(java.lang.String propertyName)
func (this_ *ObjectTypeBase) GetPropertyJSDocInfo_1_string(propertyName string) *JSDocInfo {
	_ = propertyName
	var p *Property = this_.GetSlot_1_string(propertyName)
	_ = p
	if p == nil {
		return nil
	} else {
		return p.GetJSDocInfo()
	}
}

// Java method: ObjectType.getPropertyJSDocInfo(com.google.javascript.rhino.jstype.Property.Key propertyName)
func (this_ *ObjectTypeBase) GetPropertyJSDocInfo_1_Property_Key(propertyName Property_Key) *JSDocInfo {
	_ = propertyName
	var p *Property = this_.self_ObjectType_().GetSlot_1_Property_Key(propertyName)
	_ = p
	if p == nil {
		return nil
	} else {
		return p.GetJSDocInfo()
	}
}

// Java method: ObjectType.getOwnPropertyJSDocInfo(java.lang.String propertyName)
func (this_ *ObjectTypeBase) GetOwnPropertyJSDocInfo_1_string(propertyName string) *JSDocInfo {
	_ = propertyName
	var p *Property = this_.GetOwnSlot_1_string(propertyName)
	_ = p
	if p == nil {
		return nil
	} else {
		return p.GetJSDocInfo()
	}
}

// Java method: ObjectType.getOwnPropertyJSDocInfo(com.google.javascript.rhino.jstype.Property.Key propertyName)
func (this_ *ObjectTypeBase) GetOwnPropertyJSDocInfo_1_Property_Key(propertyName Property_Key) *JSDocInfo {
	_ = propertyName
	var p *Property = this_.GetOwnSlot_1_Property_Key(propertyName)
	_ = p
	if p == nil {
		return nil
	} else {
		return p.GetJSDocInfo()
	}
}

// Java method: ObjectType.getOwnPropertyDefSite(java.lang.String propertyName)
func (this_ *ObjectTypeBase) GetOwnPropertyDefSite_1_string(propertyName string) Node {
	_ = propertyName
	var p *Property = this_.GetOwnSlot_1_string(propertyName)
	_ = p
	if p == nil {
		return Node(0)
	} else {
		return p.GetNode()
	}
}

// Java method: ObjectType.getOwnPropertyDefSite(com.google.javascript.rhino.jstype.Property.Key propertyName)
func (this_ *ObjectTypeBase) GetOwnPropertyDefSite_1_Property_Key(propertyName Property_Key) Node {
	_ = propertyName
	var p *Property = this_.GetOwnSlot_1_Property_Key(propertyName)
	_ = p
	if p == nil {
		return Node(0)
	} else {
		return p.GetNode()
	}
}

// Java method: ObjectType.setPropertyJSDocInfo(java.lang.String propertyName, com.google.javascript.rhino.JSDocInfo info)
func (this_ *ObjectTypeBase) SetPropertyJSDocInfo_2_string_JSDocInfo(propertyName string, info *JSDocInfo) {
	_ = propertyName
	_ = info
	j2gort.CheckNotNullString(propertyName)
}

// Java method: ObjectType.setPropertyJSDocInfo(com.google.javascript.rhino.jstype.Property.Key propertyName, com.google.javascript.rhino.JSDocInfo info)
func (this_ *ObjectTypeBase) SetPropertyJSDocInfo_2_Property_Key_JSDocInfo(propertyName Property_Key, info *JSDocInfo) {
	_ = propertyName
	_ = info
}

// Java method: ObjectType.setPropertyNode(java.lang.String propertyName, com.google.javascript.rhino.Node defSite)
func (this_ *ObjectTypeBase) SetPropertyNode_2_string_Node(propertyName string, defSite Node) {
	_ = propertyName
	_ = defSite
	j2gort.CheckNotNullString(propertyName)
}

// Java method: ObjectType.setPropertyNode(com.google.javascript.rhino.jstype.Property.Key propertyName, com.google.javascript.rhino.Node defSite)
func (this_ *ObjectTypeBase) SetPropertyNode_2_Property_Key_Node(propertyName Property_Key, defSite Node) {
	_ = propertyName
	_ = defSite
}

// Java method: ObjectType.findPropertyTypeWithoutConsideringTemplateTypes(com.google.javascript.rhino.jstype.Property.Key propertyName)
func (this_ *ObjectTypeBase) FindPropertyTypeWithoutConsideringTemplateTypes_1_Property_Key(propertyName Property_Key) JSType {
	_ = propertyName
	if this_.HasProperty_1_Property_Key(propertyName) {
		return this_.self_ObjectType_().GetPropertyType_1_Property_Key(propertyName)
	} else {
		return nil
	}
}

// Java method: ObjectType.findPropertyTypeWithoutConsideringTemplateTypes(java.lang.String propertyName)
func (this_ *ObjectTypeBase) FindPropertyTypeWithoutConsideringTemplateTypes_1_string(propertyName string) JSType {
	_ = propertyName
	if this_.HasProperty_1_string(propertyName) {
		return this_.GetPropertyType_1_string(propertyName)
	} else {
		return nil
	}
}

// Java method: ObjectType.getPropertyType(java.lang.String propertyName)
func (this_ *ObjectTypeBase) GetPropertyType_1_string(propertyName string) JSType {
	_ = propertyName
	var templatized *TemplatizedType
	_ = templatized
	var ok_9 bool
	templatized, ok_9 = j2gort.TryCast[*TemplatizedType](this_.self_ObjectType_())
	if ok_9 {
		var result JSType = this_.priv_ObjectType_getPropertyTypeFromSlot_9302(this_.GetSlot_1_string(propertyName))
		_ = result
		if result == nil || result.IsJavaNil_() {
			return nil
		} else {
			return this_.GetJSTypeBase().registry.ReplaceTemplateTypesPartial(result, templatized.GetTemplateTypeMap())
		}
	}
	if j2gort.InstanceOf[*FunctionType](this_.self_ObjectType_()) {
		if "call" == propertyName {
			return this_.self_ObjectType_().GetPropertyType_1_Property_Key(ToIface_Property_Key(StaticGet_ObjectType_CALL_KEY()))
		}
		if "bind" == propertyName {
			return this_.self_ObjectType_().GetPropertyType_1_Property_Key(ToIface_Property_Key(StaticGet_ObjectType_BIND_KEY()))
		}
		if "apply" == propertyName {
			return this_.self_ObjectType_().GetPropertyType_1_Property_Key(ToIface_Property_Key(StaticGet_ObjectType_APPLY_KEY()))
		}
	}
	return this_.priv_ObjectType_getPropertyTypeFromSlot_9302(this_.GetSlot_1_string(propertyName))
}

// Java method: ObjectType.getPropertyType(com.google.javascript.rhino.jstype.Property.Key propertyName)
func (this_ *ObjectTypeBase) GetPropertyType_1_Property_Key(propertyName Property_Key) JSType {
	_ = propertyName
	var slot *Property = this_.self_ObjectType_().GetSlot_1_Property_Key(propertyName)
	_ = slot
	return this_.priv_ObjectType_getPropertyTypeFromSlot_9302(slot)
}

// Java method: ObjectType.getPropertyTypeFromSlot(com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable Property slot)
func (this_ *ObjectTypeBase) priv_ObjectType_getPropertyTypeFromSlot_9302(slot *Property) JSType {
	_ = slot
	if slot == nil {
		if this_.IsNoResolvedType() || this_.IsCheckedUnknownType() {
			return this_.GetNativeType(StaticGet_JSTypeNative_CHECKED_UNKNOWN_TYPE())
		} else if this_.IsEmptyType() {
			return this_.GetNativeType(StaticGet_JSTypeNative_NO_TYPE())
		}
		return this_.GetNativeType(StaticGet_JSTypeNative_UNKNOWN_TYPE())
	}
	return slot.GetType()
}

// Java method: ObjectType.getPropertyKind(com.google.javascript.rhino.jstype.Property.Key propertyName, boolean autobox)
func (this_ *ObjectTypeBase) GetPropertyKind_2_Property_Key_bool(propertyName Property_Key, autobox bool) JSType_HasPropertyKind {
	_ = propertyName
	_ = autobox
	return JSType_HasPropertyKind_Of(((this_.IsEmptyType() || this_.self_ObjectType_().IsUnknownType()) || (this_.self_ObjectType_().GetSlot_1_Property_Key(propertyName) != nil)))
}

// Java method: ObjectType.getPropertyKind(java.lang.String propertyName, boolean autobox)
func (this_ *ObjectTypeBase) GetPropertyKind_2_string_bool(propertyName string, autobox bool) JSType_HasPropertyKind {
	_ = propertyName
	_ = autobox
	j2gort.CheckNotNullString(propertyName)
	return JSType_HasPropertyKind_Of(((this_.IsEmptyType() || this_.self_ObjectType_().IsUnknownType()) || (this_.GetSlot_1_string(propertyName) != nil)))
}

// Java method: ObjectType.getOwnPropertyKind(java.lang.String propertyName)
func (this_ *ObjectTypeBase) GetOwnPropertyKind_1_string(propertyName string) JSType_HasPropertyKind {
	_ = propertyName
	if this_.GetOwnSlot_1_string(propertyName) != nil {
		return StaticGet_JSType_HasPropertyKind_KNOWN_PRESENT()
	} else {
		return StaticGet_JSType_HasPropertyKind_ABSENT()
	}
}

// Java method: ObjectType.getOwnPropertyKind(com.google.javascript.rhino.jstype.Property.Key propertyName)
func (this_ *ObjectTypeBase) GetOwnPropertyKind_1_Property_Key(propertyName Property_Key) JSType_HasPropertyKind {
	_ = propertyName
	if this_.GetOwnSlot_1_Property_Key(propertyName) != nil {
		return StaticGet_JSType_HasPropertyKind_KNOWN_PRESENT()
	} else {
		return StaticGet_JSType_HasPropertyKind_ABSENT()
	}
}

// Java method: ObjectType.hasOwnProperty(java.lang.String propertyName)
func (this_ *ObjectTypeBase) HasOwnProperty_1_string(propertyName string) bool {
	_ = propertyName
	return (this_.GetOwnSlot_1_string(propertyName) != nil)
}

// Java method: ObjectType.hasOwnProperty(com.google.javascript.rhino.jstype.Property.Key propertyName)
func (this_ *ObjectTypeBase) HasOwnProperty_1_Property_Key(propertyName Property_Key) bool {
	_ = propertyName
	return !(func() bool {
		re_10 := this_.GetOwnPropertyKind_1_Property_Key(propertyName)
		ae_11 := StaticGet_JSType_HasPropertyKind_ABSENT()
		if re_10 == 0 {
			j2gort.ThrowNPE()
		}
		return re_10 == ae_11
	}())
}

// Java method: ObjectType.getOwnPropertyKnownSymbols()
func (this_ *ObjectTypeBase) GetOwnPropertyKnownSymbols() j2gort.Ext_java_util_Set {
	var src PrototypeObjectType = ObjectType_UnwrapToPrototypeObjectType(this_.self_ObjectType_())
	_ = src
	if src == nil || src.IsJavaNil_() {
		return j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0())
	}
	var map_ *PropertyMap = src.GetOwnPropertyMapOrNull()
	_ = map_
	if map_ != nil {
		return map_.GetOwnKnownSymbols()
	} else {
		return j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0())
	}
}

// Java method: ObjectType.getOwnPropertyNames()
func (this_ *ObjectTypeBase) GetOwnPropertyNames() j2gort.Ext_java_util_Set {
	var src PrototypeObjectType = ObjectType_UnwrapToPrototypeObjectType(this_.self_ObjectType_())
	_ = src
	if src == nil || src.IsJavaNil_() {
		return j2gort.AsIface[j2gort.Ext_java_util_Set](StaticGet_ObjectType_EMPTY_STRING_SET())
	}
	var map_ *PropertyMap = src.GetOwnPropertyMapOrNull()
	_ = map_
	if map_ != nil {
		return map_.GetOwnPropertyNames()
	} else {
		return j2gort.AsIface[j2gort.Ext_java_util_Set](StaticGet_ObjectType_EMPTY_STRING_SET())
	}
}

// Java method: ObjectType.getOwnPropertyNamesUpperBound()
func (this_ *ObjectTypeBase) GetOwnPropertyNamesUpperBound() int32 {
	var count int32 = int32(0)
	_ = count
	var fn *FunctionType
	_ = fn
	var ok_12 bool
	fn, ok_12 = j2gort.TryCast[*FunctionType](this_.self_ObjectType_())
	if ok_12 && (fn.HasPrototypeSlot()) {
		count = int32(1)
	}
	var src PrototypeObjectType = ObjectType_UnwrapToPrototypeObjectType(this_.self_ObjectType_())
	_ = src
	if src != nil && !src.IsJavaNil_() {
		var map_ *PropertyMap = src.GetOwnPropertyMapOrNull()
		_ = map_
		if map_ != nil {
			count = (count + map_.GetOwnPropertiesCount())
		}
	}
	return count
}

// Java method: ObjectType.collectOwnPropertyNames(java.lang.String[] out)
func (this_ *ObjectTypeBase) CollectOwnPropertyNames(out []string) int32 {
	_ = out
	var idx int32 = int32(0)
	_ = idx
	var fn *FunctionType
	_ = fn
	var hasProto bool = ((j2gort.InstanceOfBind[*FunctionType](this_.self_ObjectType_(), &fn)) && fn.HasPrototypeSlot())
	_ = hasProto
	if hasProto {
		{
			a_13 := out
			i_14 := func() int32 { old_15 := idx; nw_16 := old_15 + 1; idx = nw_16; return old_15 }()
			a_13[i_14] = "prototype"
		}
	}
	var src PrototypeObjectType = ObjectType_UnwrapToPrototypeObjectType(this_.self_ObjectType_())
	_ = src
	if src != nil && !src.IsJavaNil_() {
		var map_ *PropertyMap = src.GetOwnPropertyMapOrNull()
		_ = map_
		if map_ != nil {
			idx = map_.CollectOwnPropertyNames(out, idx, hasProto)
		}
	}
	return idx
}

// Java method: ObjectType.getOwnPropertyKeys()
func (this_ *ObjectTypeBase) GetOwnPropertyKeys() j2gort.Ext_com_google_common_collect_ImmutableSet {
	var ownPropertyNames j2gort.Ext_java_util_Set = this_.self_ObjectType_().GetOwnPropertyNames()
	_ = ownPropertyNames
	var ownSymbols j2gort.Ext_java_util_Set = this_.GetOwnPropertyKnownSymbols()
	_ = ownSymbols
	var ownPropertyKeys j2gort.Ext_com_google_common_collect_ImmutableSet_Builder = j2gort.Ext_com_google_common_collect_ImmutableSet_S_Builder()
	_ = ownPropertyKeys
	{
		cur_17 := j2gort.NewIterCursor(j2gort.CheckNotNull(ownPropertyNames))
		for cur_17.HasNext() {
			propertyName := cur_17.NextStr()
			_ = propertyName
			ownPropertyKeys.Add_1_any(NewProperty_StringKey(propertyName))
		}
	}
	{
		cur_18 := j2gort.NewIterCursor(j2gort.CheckNotNull(ownSymbols))
		for cur_18.HasNext() {
			rawElem_19 := cur_18.Next()
			symbol := j2gort.CastPtr[KnownSymbolType](rawElem_19)
			_ = symbol
			ownPropertyKeys.Add_1_any(NewProperty_SymbolKey(symbol))
		}
	}
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](ownPropertyKeys.Build())
}

// Java method: ObjectType.isPropertyTypeInferred(java.lang.String propertyName)
func (this_ *ObjectTypeBase) IsPropertyTypeInferred_1_string(propertyName string) bool {
	_ = propertyName
	var slot *Property = this_.GetSlot_1_string(propertyName)
	_ = slot
	return ((slot != nil) && slot.IsTypeInferred())
}

// Java method: ObjectType.isPropertyTypeInferred(com.google.javascript.rhino.jstype.Property.Key propertyName)
func (this_ *ObjectTypeBase) IsPropertyTypeInferred_1_Property_Key(propertyName Property_Key) bool {
	_ = propertyName
	var slot *Property = this_.self_ObjectType_().GetSlot_1_Property_Key(propertyName)
	_ = slot
	return ((slot != nil) && slot.IsTypeInferred())
}

// Java method: ObjectType.isPropertyTypeDeclared(java.lang.String propertyName)
func (this_ *ObjectTypeBase) IsPropertyTypeDeclared_1_string(propertyName string) bool {
	_ = propertyName
	var slot *Property = this_.GetSlot_1_string(propertyName)
	_ = slot
	return ((slot != nil) && !(slot.IsTypeInferred()))
}

// Java method: ObjectType.isPropertyTypeDeclared(com.google.javascript.rhino.jstype.Property.Key propertyName)
func (this_ *ObjectTypeBase) IsPropertyTypeDeclared_1_Property_Key(propertyName Property_Key) bool {
	_ = propertyName
	var slot *Property = this_.self_ObjectType_().GetSlot_1_Property_Key(propertyName)
	_ = slot
	return ((slot != nil) && !(slot.IsTypeInferred()))
}

// Java method: ObjectType.isStructuralType()
func (this_ *ObjectTypeBase) IsStructuralType() bool {
	var constructor *FunctionType = this_.self_ObjectType_().GetConstructor()
	_ = constructor
	return ((constructor != nil) && constructor.IsStructuralInterface())
}

// Java method: ObjectType.hasOwnDeclaredProperty(java.lang.String name)
func (this_ *ObjectTypeBase) HasOwnDeclaredProperty_1_string(name string) bool {
	_ = name
	return (this_.HasOwnProperty_1_string(name) && this_.IsPropertyTypeDeclared_1_string(name))
}

// Java method: ObjectType.hasOwnDeclaredProperty(com.google.javascript.rhino.jstype.Property.Key name)
func (this_ *ObjectTypeBase) HasOwnDeclaredProperty_1_Property_Key(name Property_Key) bool {
	_ = name
	return (this_.HasOwnProperty_1_Property_Key(name) && this_.IsPropertyTypeDeclared_1_Property_Key(name))
}

// Java method: ObjectType.isPropertyInExterns(java.lang.String propertyName)
func (this_ *ObjectTypeBase) IsPropertyInExterns(propertyName string) bool {
	_ = propertyName
	var p *Property = this_.GetSlot_1_string(propertyName)
	_ = p
	if p == nil {
		return false
	} else {
		return p.IsFromExterns()
	}
}

// Java method: ObjectType.getPropertiesCount()
func (this_ *ObjectTypeBase) GetPropertiesCount() int32 {
	var src PrototypeObjectType = ObjectType_UnwrapToPrototypeObjectType(this_.self_ObjectType_())
	_ = src
	if src == nil || src.IsJavaNil_() {
		return int32(0)
	}
	if j2gort.IfaceEquals[PrototypeObjectType](ObjectType_UnwrapToPrototypeObjectType(src.GetImplicitPrototype()), nil) {
		var map_ *PropertyMap = src.GetOwnPropertyMapOrNull()
		_ = map_
		if map_ != nil {
			return map_.GetOwnPropertiesCount()
		} else {
			return int32(0)
		}
	}
	return src.GetPropertyMap().GetPropertiesCount()
}

// Java method: ObjectType.getPropertyNames()
func (this_ *ObjectTypeBase) GetPropertyNames() j2gort.Ext_com_google_common_collect_ImmutableSortedSet {
	return this_.self_ObjectType_().GetPropertyMap().GetAllKeys().StringKeys()
}

// Java method: ObjectType.getAllKeys()
func (this_ *ObjectTypeBase) GetAllKeys() *PropertyMap_AllKeys {
	return this_.self_ObjectType_().GetPropertyMap().GetAllKeys()
}

// Java method: ObjectType.visit(com.google.javascript.rhino.jstype.Visitor<T> visitor)
func (this_ *ObjectTypeBase) Visit_1(visitor Visitor) any {
	_ = visitor
	return visitor.CaseObjectType(this_.self_ObjectType_())
}

// Java method: ObjectType.visit(com.google.javascript.rhino.jstype.RelationshipVisitor<T> visitor, com.google.javascript.rhino.jstype.JSType that)
func (this_ *ObjectTypeBase) Visit_2(visitor RelationshipVisitor, that JSType) any {
	_ = visitor
	_ = that
	return visitor.CaseObjectType(this_.self_ObjectType_(), that)
}

// Java method: ObjectType.isImplicitPrototypeOf(com.google.javascript.rhino.jstype.ObjectType other)
func (this_ *ObjectTypeBase) IsImplicitPrototypeOf(other ObjectType) bool {
	_ = other
	var unwrappedThis ObjectType = ObjectType_DeeplyUnwrap(this_.self_ObjectType_())
	_ = unwrappedThis
	{
		other = ObjectType_DeeplyUnwrap(other)
		for ; other != nil && !other.IsJavaNil_(); other = ObjectType_DeeplyUnwrap(other.GetImplicitPrototype()) {
			if unwrappedThis == other || ((unwrappedThis == nil || unwrappedThis.IsJavaNil_()) && (other == nil || other.IsJavaNil_())) {
				return true
			}
		}
	}
	return false
}

// Java static method: ObjectType.deeplyUnwrap(com.google.javascript.rhino.jstype.ObjectType original)
func ObjectType_DeeplyUnwrap(original ObjectType) ObjectType {
	clinit_ObjectType()
	_ = original
	var current ObjectType = original
	_ = current
	for j2gort.InstanceOf[ProxyObjectType](current) {
		if current.IsTemplatizedType() {
			current = current.ToMaybeTemplatizedType().GetReferencedType()
		} else if current.IsNamedType() {
			if !(current.IsSuccessfullyResolved()) {
				break
			}
			current = current.ToMaybeNamedType().GetReferencedObjTypeInternal()
		} else {
			current = (j2gort.AsIface[ProxyObjectType](current)).GetReferencedObjTypeInternal()
		}
	}
	return current
}

// Java method: ObjectType.getPossibleToBooleanOutcomes()
func (this_ *ObjectTypeBase) GetPossibleToBooleanOutcomes() BooleanLiteralSet {
	return StaticGet_BooleanLiteralSet_TRUE()
}

// Java method: ObjectType.isUnknownType()
func (this_ *ObjectTypeBase) IsUnknownType() bool {
	if (int32(this_.GetJSTypeBase().jstypeFlags) & int32(int8(32))) == int32(0) {
		var implicitProto ObjectType = this_.self_ObjectType_().GetImplicitPrototype()
		_ = implicitProto
		if (implicitProto == nil || implicitProto.IsJavaNil_()) || implicitProto.IsNativeObjectType() {
			this_.GetJSTypeBase().jstypeFlags = int8((int32(this_.GetJSTypeBase().jstypeFlags) | int32(int8(32))))
			{
				cur_20 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.self_ObjectType_().GetCtorExtendedInterfaces())))
				for cur_20.HasNext() {
					rawElem_21 := cur_20.Next()
					interfaceType := j2gort.AsIface[ObjectType](rawElem_21)
					_ = interfaceType
					if interfaceType.IsUnknownType() {
						this_.GetJSTypeBase().jstypeFlags = int8((int32(this_.GetJSTypeBase().jstypeFlags) & ^(int32(int8(32)))))
						break
					}
				}
			}
		} else if !(implicitProto.IsUnknownType()) {
			this_.GetJSTypeBase().jstypeFlags = int8((int32(this_.GetJSTypeBase().jstypeFlags) | int32(int8(32))))
		}
	}
	return ((int32(this_.GetJSTypeBase().jstypeFlags) & int32(int8(32))) == int32(0))
}

// Java method: ObjectType.isObject()
func (this_ *ObjectTypeBase) IsObject() bool {
	return true
}

// Java method: ObjectType.hasCachedValues()
func (this_ *ObjectTypeBase) HasCachedValues() bool {
	return ((int32(this_.GetJSTypeBase().jstypeFlags) & int32(int8(32))) != int32(0))
}

// Java method: ObjectType.clearCachedValues()
func (this_ *ObjectTypeBase) ClearCachedValues() {
	this_.GetJSTypeBase().jstypeFlags = int8((int32(this_.GetJSTypeBase().jstypeFlags) & ^(int32(int8(32)))))
}

// Java static method: ObjectType.cast(com.google.javascript.rhino.jstype.JSType type)
func ObjectType_Cast(type_ JSType) ObjectType {
	clinit_ObjectType()
	_ = type_
	if type_ == nil || type_.IsJavaNil_() {
		return nil
	} else {
		return type_.ToObjectType()
	}
}

// Java method: ObjectType.isFunctionPrototypeType()
func (this_ *ObjectTypeBase) IsFunctionPrototypeType() bool {
	return (this_.self_ObjectType_().GetOwnerFunction() != nil)
}

// Java method: ObjectType.getOwnerFunction()
func (this_ *ObjectTypeBase) GetOwnerFunction() *FunctionType {
	return nil
}

// Java method: ObjectType.setOwnerFunction(com.google.javascript.rhino.jstype.FunctionType type)
func (this_ *ObjectTypeBase) SetOwnerFunction(type_ *FunctionType) {
	_ = type_
}

// Java method: ObjectType.getCtorImplementedInterfaces()
func (this_ *ObjectTypeBase) GetCtorImplementedInterfaces() any {
	return j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()
}

// Java method: ObjectType.getCtorExtendedInterfaces()
func (this_ *ObjectTypeBase) GetCtorExtendedInterfaces() any {
	return j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()
}

// Java method: ObjectType.getPropertyTypeMap()
func (this_ *ObjectTypeBase) GetPropertyTypeMap() j2gort.Ext_java_util_Map {
	var propTypeMap j2gort.Ext_com_google_common_collect_ImmutableMap_Builder = j2gort.Ext_com_google_common_collect_ImmutableMap_S_Builder()
	_ = propTypeMap
	{
		cur_22 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.GetPropertyNames()))
		for cur_22.HasNext() {
			name := cur_22.NextStr()
			_ = name
			propTypeMap.Put_2(j2gort.BoxString(name), this_.GetPropertyType_1_string(name))
		}
	}
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](propTypeMap.BuildOrThrow())
}

// Java method: ObjectType.getEnumeratedTypeOfEnumObject()
func (this_ *ObjectTypeBase) GetEnumeratedTypeOfEnumObject() JSType {
	return nil
}

// Java class:
type ObjectType__Anon_159 struct {
	j2gort.Ext_com_google_common_collect_AbstractIteratorBase
	outer_     ObjectType
	cap_0_self ObjectType
	next       ObjectType `gson:"next"`
	idHash_    atomic.Int32
}

func (this_ *ObjectType__Anon_159) BindSelf_(s any) {
	_ = s
	this_.Ext_com_google_common_collect_AbstractIteratorBase.BindExtSelf_(s)
}
func (this_ *ObjectType__Anon_159) IsJavaNil_() bool { return this_ == nil }
func (*ObjectType__Anon_159) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.ObjectType$1", ""
}
func (this_ *ObjectType__Anon_159) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (this_ *ObjectType__Anon_159) SetOuter_(o ObjectType)                                          { this_.outer_ = o }
func (this_ *ObjectType__Anon_159) GetOuter_() ObjectType                                           { return this_.outer_ }
func (*ObjectType__Anon_159) Implements_j2gort_Ext_com_google_common_collect_AbstractIterator()     {}
func (*ObjectType__Anon_159) Implements_j2gort_Ext_com_google_common_collect_UnmodifiableIterator() {}
func (*ObjectType__Anon_159) Implements_j2gort_Ext_java_util_Iterator()                             {}

func (this_ *ObjectType__Anon_159) initFields_ObjectType__Anon_159() {
	_ = this_
	this_.next = this_.cap_0_self
}

func NewObjectType__Anon_159(outer_ ObjectType, cap_p_self ObjectType) *ObjectType__Anon_159 {
	this_ := &ObjectType__Anon_159{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	this_.cap_0_self = cap_p_self
	j2gort.InitExt_com_google_common_collect_AbstractIterator(&this_.Ext_com_google_common_collect_AbstractIteratorBase)
	this_.initFields_ObjectType__Anon_159()
	return this_
}

// Java method: .computeNext()
func (this_ *ObjectType__Anon_159) ComputeNext() any {
	this_.next = this_.next.GetImplicitPrototype()
	if this_.next != nil && !this_.next.IsJavaNil_() {
		return this_.next
	} else {
		return j2gort.AsIface[ObjectType](this_.EndOfData())
	}
}
