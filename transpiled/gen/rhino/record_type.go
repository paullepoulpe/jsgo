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
// Port of com/google/javascript/rhino/jstype/RecordType.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.rhino.jstype.RecordType
type RecordType struct {
	PrototypeObjectTypeBase
}

func (this_ *RecordType) BindSelf_(s any) {
	_ = s
	this_.PrototypeObjectTypeBase.BindSelf_(s)
}
func (this_ *RecordType) IsJavaNil_() bool { return this_ == nil }
func (*RecordType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.RecordType", "RecordType"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.RecordType", "com.google.javascript.rhino.jstype.RecordType", func(v_ any) bool {
		return j2gort.InstanceOf[*RecordType](v_)
	})
}
func (*RecordType) Implements_rhino_PrototypeObjectType() {}
func (*RecordType) Implements_rhino_ObjectType()          {}
func (*RecordType) Implements_rhino_JSType()              {}

var static_RecordType_TYPE_CLASS JSTypeClass

func StaticGet_RecordType_TYPE_CLASS() JSTypeClass {
	clinit_RecordType()
	return static_RecordType_TYPE_CLASS
}
func StaticSet_RecordType_TYPE_CLASS(v JSTypeClass) {
	clinit_RecordType()
	static_RecordType_TYPE_CLASS = v
}

var clinitDone_RecordType uint32
var clinitState_RecordType j2gort.ClinitState

func clinit_RecordType() {
	if atomic.LoadUint32(&clinitDone_RecordType) != 0 {
		return
	}
	clinitSlow_RecordType()
}

func clinitSlow_RecordType() {
	if !j2gort.EnterClinit(&clinitState_RecordType, &clinitDone_RecordType) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_RecordType, &clinitDone_RecordType)
	static_RecordType_TYPE_CLASS = StaticGet_JSTypeClass_RECORD()
}

func (this_ *RecordType) initFields_RecordType() {
	_ = this_
}

// Java constructor: RecordType(com.google.javascript.rhino.jstype.JSTypeRegistry registry, java.util.Map<java.lang.String,com.google.javascript.rhino.jstype.RecordTypeBuilder.RecordProperty> properties, boolean declared)
func NewRecordType(registry *JSTypeRegistry, properties j2gort.Ext_java_util_Map, declared bool) *RecordType {
	this_ := &RecordType{}
	this_.BindSelf_(this_)
	InitRecordType(this_, registry, properties, declared)
	return this_
}

func InitRecordType(this_ *RecordType, registry *JSTypeRegistry, properties j2gort.Ext_java_util_Map, declared bool) {
	_ = registry
	_ = properties
	_ = declared
	InitPrototypeObjectType_2(&this_.PrototypeObjectTypeBase, PrototypeObjectType_S_Builder(registry), int32(2048))
	this_.initFields_RecordType()
	this_.SetPrettyPrint(true)
	if declared {
		this_.GetJSTypeBase().protoFlags = int8((int32(this_.GetJSTypeBase().protoFlags) | int32(int8(32))))
	}
	{
		cur_1 := j2gort.NewMapKeyStrCursor(properties)
		for cur_1.HasNext() {
			property := cur_1.NextStr()
			_ = property
			var prop *RecordTypeBuilder_RecordProperty = j2gort.CastPtr[RecordTypeBuilder_RecordProperty](j2gort.MapGetStr(properties, property))
			_ = prop
			if prop == nil {
				panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "RecordProperty associated with a property should not be null!")))
			}
			if declared {
				this_.DefineDeclaredProperty_3_string_JSType_Node(property, prop.GetType(), prop.GetPropertyNode())
			} else {
				this_.DefineSynthesizedProperty(property, prop.GetType(), prop.GetPropertyNode())
			}
		}
	}
	this_.GetJSTypeBase().protoFlags = int8((int32(this_.GetJSTypeBase().protoFlags) | int32(int8(64))))
	registry.GetResolver().ResolveIfClosed(ToIface_JSType(this_), StaticGet_RecordType_TYPE_CLASS())
}

// Java method: RecordType.getTypeClass()
func (this_ *RecordType) GetTypeClass() JSTypeClass {
	return StaticGet_RecordType_TYPE_CLASS()
}

// Java method: RecordType.isSynthetic()
func (this_ *RecordType) IsSynthetic() bool {
	return ((int32(this_.GetJSTypeBase().protoFlags) & int32(int8(32))) == int32(0))
}

// Java method: RecordType.getImplicitPrototype()
func (this_ *RecordType) GetImplicitPrototype() ObjectType {
	return this_.GetJSTypeBase().registry.GetNativeObjectType(StaticGet_JSTypeNative_OBJECT_TYPE())
}

// Java method: RecordType.defineProperty(java.lang.String propertyName, com.google.javascript.rhino.jstype.JSType type, boolean inferred, com.google.javascript.rhino.Node propertyNode)
func (this_ *RecordType) DefineProperty_4_string_JSType_bool_Node(propertyName string, type_ JSType, inferred bool, propertyNode Node) bool {
	_ = propertyName
	_ = type_
	_ = inferred
	_ = propertyNode
	if (int32(this_.GetJSTypeBase().protoFlags) & int32(int8(64))) != int32(0) {
		return false
	}
	return this_.PrototypeObjectTypeBase.DefineProperty_4_string_JSType_bool_Node(propertyName, type_, inferred, propertyNode)
}

// Java method: RecordType.defineProperty(com.google.javascript.rhino.jstype.Property.Key propertyName, com.google.javascript.rhino.jstype.JSType type, boolean inferred, com.google.javascript.rhino.Node propertyNode)
func (this_ *RecordType) DefineProperty_4_Property_Key_JSType_bool_Node(propertyName Property_Key, type_ JSType, inferred bool, propertyNode Node) bool {
	_ = propertyName
	_ = type_
	_ = inferred
	_ = propertyNode
	if (int32(this_.GetJSTypeBase().protoFlags) & int32(int8(64))) != int32(0) {
		return false
	}
	return this_.PrototypeObjectTypeBase.DefineProperty_4_Property_Key_JSType_bool_Node(propertyName, type_, inferred, propertyNode)
}

// Java method: RecordType.getGreatestSubtypeHelper(com.google.javascript.rhino.jstype.JSType that)
func (this_ *RecordType) GetGreatestSubtypeHelper(that JSType) JSType {
	_ = that
	if that.GetJSTypeBase().IsRecordType() {
		var thatRecord *RecordType = that.ToMaybeRecordType()
		_ = thatRecord
		var val_2 RecordTypeBuilder
		InitRecordTypeBuilder(&val_2, this_.GetJSTypeBase().registry)
		var builder *RecordTypeBuilder = &val_2
		_ = builder
		builder.SetSynthesized(true)
		var noType JSType = Upcast_JSType(this_.GetJSTypeBase().registry.GetNativeObjectType(StaticGet_JSTypeNative_NO_TYPE()))
		_ = noType
		var thisMap *PropertyMap = this_.GetOwnPropertyMapOrNull()
		_ = thisMap
		if thisMap != nil {
			if thisMap.HasFlatProps() {
				{
					var i int32 = int32(0)
					_ = i
					var n int32 = thisMap.GetOwnPropertiesCount()
					_ = n
					for ; i < n; i = i + 1 {
						var thisProp *Property = thisMap.GetFlatPropAt(i)
						_ = thisProp
						var property string = thisProp.GetName()
						_ = property
						var thisPropertyType JSType = this_.GetPropertyType_1_string(property)
						_ = thisPropertyType
						var propType JSType = nil
						_ = propType
						if thatRecord.HasProperty_1_string(property) {
							var thatPropertyType JSType = thatRecord.GetPropertyType_1_string(property)
							_ = thatPropertyType
							propType = thisPropertyType.GetJSTypeBase().GetGreatestSubtype(thatPropertyType)
							if propType.GetJSTypeBase().EqualsJSType(noType) {
								return noType
							}
						} else {
							propType = thisPropertyType
						}
						builder.AddProperty(property, propType, thisProp.GetNode())
					}
				}
			} else {
				{
					cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.GetOwnPropertyNames()))
					for cur_3.HasNext() {
						property_2 := cur_3.NextStr()
						_ = property_2
						var thisPropertyType_2 JSType = this_.GetPropertyType_1_string(property_2)
						_ = thisPropertyType_2
						var propType_2 JSType = nil
						_ = propType_2
						if thatRecord.HasProperty_1_string(property_2) {
							var thatPropertyType_2 JSType = thatRecord.GetPropertyType_1_string(property_2)
							_ = thatPropertyType_2
							propType_2 = thisPropertyType_2.GetJSTypeBase().GetGreatestSubtype(thatPropertyType_2)
							if propType_2.GetJSTypeBase().EqualsJSType(noType) {
								return noType
							}
						} else {
							propType_2 = thisPropertyType_2
						}
						builder.AddProperty(property_2, propType_2, this_.GetPropertyNode_1_string(property_2))
					}
				}
			}
		}
		var thatMap *PropertyMap = thatRecord.GetOwnPropertyMapOrNull()
		_ = thatMap
		if thatMap != nil {
			if thatMap.HasFlatProps() {
				{
					var i_2 int32 = int32(0)
					_ = i_2
					var n_2 int32 = thatMap.GetOwnPropertiesCount()
					_ = n_2
					for ; i_2 < n_2; i_2 = i_2 + 1 {
						var thatProp *Property = thatMap.GetFlatPropAt(i_2)
						_ = thatProp
						var property_3 string = thatProp.GetName()
						_ = property_3
						if !(this_.HasProperty_1_string(property_3)) {
							builder.AddProperty(property_3, thatRecord.GetPropertyType_1_string(property_3), thatProp.GetNode())
						}
					}
				}
			} else {
				{
					cur_4 := j2gort.NewIterCursor(j2gort.CheckNotNull(thatRecord.GetOwnPropertyNames()))
					for cur_4.HasNext() {
						property_4 := cur_4.NextStr()
						_ = property_4
						if !(this_.HasProperty_1_string(property_4)) {
							builder.AddProperty(property_4, thatRecord.GetPropertyType_1_string(property_4), thatRecord.GetPropertyNode_1_string(property_4))
						}
					}
				}
			}
		}
		return builder.Build()
	}
	var greatestSubtype JSType = this_.GetJSTypeBase().registry.GetNativeType(StaticGet_JSTypeNative_NO_OBJECT_TYPE())
	_ = greatestSubtype
	var thatRestrictedToObj JSType = this_.GetJSTypeBase().registry.GetNativeType(StaticGet_JSTypeNative_OBJECT_TYPE()).GetJSTypeBase().GetGreatestSubtype(that)
	_ = thatRestrictedToObj
	if !(thatRestrictedToObj.GetJSTypeBase().IsEmptyType()) {
		var thisMap_2 *PropertyMap = this_.GetOwnPropertyMapOrNull()
		_ = thisMap_2
		if thisMap_2 != nil {
			if thisMap_2.HasFlatProps() {
				{
					var i_3 int32 = int32(0)
					_ = i_3
					var n_3 int32 = thisMap_2.GetOwnPropertiesCount()
					_ = n_3
					for ; i_3 < n_3; i_3 = i_3 + 1 {
						var propName string = thisMap_2.GetFlatPropAt(i_3).GetName()
						_ = propName
						var propType_3 JSType = this_.GetPropertyType_1_string(propName)
						_ = propType_3
						var builder_2 *UnionType_Builder = this_.GetJSTypeBase().registry.AcquireUnionTypeBuilder()
						_ = builder_2
						var propUnion JSType = nil
						_ = propUnion
						func() {
							defer func() {
								this_.GetJSTypeBase().registry.ReleaseUnionTypeBuilder(builder_2)
							}()
							{
								cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.GetJSTypeBase().registry.GetEachReferenceTypeWithProperty(propName)))
								for cur_5.HasNext() {
									rawElem_6 := cur_5.Next()
									alt := j2gort.AsIface[ObjectType](rawElem_6)
									_ = alt
									var altPropType JSType = alt.GetObjectTypeBase().GetPropertyType_1_string(propName)
									_ = altPropType
									if (((altPropType != nil && !altPropType.IsJavaNil_()) && !(alt.GetObjectTypeBase().GetJSTypeBase().EqualsJSType(this_))) && alt.GetObjectTypeBase().IsSubtypeOf_1(that)) && altPropType.GetJSTypeBase().IsSubtypeOf_1(propType_3) {
										builder_2.AddAlternate(Upcast_JSType(alt))
									}
								}
							}
							propUnion = builder_2.Build()
						}()
						greatestSubtype = greatestSubtype.GetLeastSupertype(propUnion)
					}
				}
			} else {
				{
					cur_7 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.GetOwnPropertyNames()))
					for cur_7.HasNext() {
						propName_2 := cur_7.NextStr()
						_ = propName_2
						var propType_4 JSType = this_.GetPropertyType_1_string(propName_2)
						_ = propType_4
						var builder_3 *UnionType_Builder = this_.GetJSTypeBase().registry.AcquireUnionTypeBuilder()
						_ = builder_3
						var propUnion_2 JSType = nil
						_ = propUnion_2
						func() {
							defer func() {
								this_.GetJSTypeBase().registry.ReleaseUnionTypeBuilder(builder_3)
							}()
							{
								cur_8 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.GetJSTypeBase().registry.GetEachReferenceTypeWithProperty(propName_2)))
								for cur_8.HasNext() {
									rawElem_9 := cur_8.Next()
									alt_2 := j2gort.AsIface[ObjectType](rawElem_9)
									_ = alt_2
									var altPropType_2 JSType = alt_2.GetObjectTypeBase().GetPropertyType_1_string(propName_2)
									_ = altPropType_2
									if (((altPropType_2 != nil && !altPropType_2.IsJavaNil_()) && !(alt_2.GetObjectTypeBase().GetJSTypeBase().EqualsJSType(this_))) && alt_2.GetObjectTypeBase().IsSubtypeOf_1(that)) && altPropType_2.GetJSTypeBase().IsSubtypeOf_1(propType_4) {
										builder_3.AddAlternate(Upcast_JSType(alt_2))
									}
								}
							}
							propUnion_2 = builder_3.Build()
						}()
						greatestSubtype = greatestSubtype.GetLeastSupertype(propUnion_2)
					}
				}
			}
		}
	}
	return greatestSubtype
}

// Java method: RecordType.toMaybeRecordType()
func (this_ *RecordType) ToMaybeRecordType() *RecordType {
	return this_
}

// Java method: RecordType.isStructuralType()
func (this_ *RecordType) IsStructuralType() bool {
	return true
}
