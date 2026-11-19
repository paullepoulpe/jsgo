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
// Port of com/google/javascript/rhino/jstype/PrototypeObjectType.java (generated mechanically by j2go).

package rhino

import (
	"fmt"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_prototype_object_type_1 j2gort.Ext_java_util_Comparator

func init() {
	lambdaSingleton_prototype_object_type_1 = j2gort.NewLambda_Ext_java_util_Comparator(func(lp_1 any, lp_2 any) int32 {
		thisProp := j2gort.CastPtr[PrototypeObjectType_PropertyForPrettyPrinting](lp_1)
		_ = thisProp
		otherProp := j2gort.CastPtr[PrototypeObjectType_PropertyForPrettyPrinting](lp_2)
		_ = otherProp
		{
			var iifeCond_1 bool
			{

				sr_3 := thisProp.Key().HumanReadableName()
				se_4 := otherProp.Key().HumanReadableName()
				iifeCond_1 = j2gort.CheckNotNullString(sr_3) == se_4
			}
			if !(iifeCond_1) {
				return j2gort.JStrCompareTo(thisProp.Key().HumanReadableName(), otherProp.Key().HumanReadableName())
			}
		}
		switch (thisProp.Key().Kind()).Ordinal() {
		case int32(0) /* KeyKind.STRING */ :
			return func() int32 {
				switch (otherProp.Key().Kind()).Ordinal() {
				case int32(0) /* KeyKind.STRING */ :
					return j2gort.JStrCompareTo(thisProp.Key().M_String(), otherProp.Key().M_String())
				case int32(1) /* KeyKind.SYMBOL */ :
					return int32(-1)
				}
				panic("unreachable switch expression")
			}()
		case int32(1) /* KeyKind.SYMBOL */ :
			return func() int32 {
				switch (otherProp.Key().Kind()).Ordinal() {
				case int32(0) /* KeyKind.STRING */ :
					return int32(1)
				case int32(1) /* KeyKind.SYMBOL */ :
					return j2gort.JStrCompareTo(thisProp.Type().GetJSTypeBase().String(), otherProp.Type().GetJSTypeBase().String())
				}
				panic("unreachable switch expression")
			}()
		}
		panic("unreachable switch expression")
	})
}

// Java polymorphic class interface: com.google.javascript.rhino.jstype.PrototypeObjectType
type PrototypeObjectType interface {
	Implements_rhino_PrototypeObjectType()
	IsJavaNil_() bool
	GetPrototypeObjectTypeBase() *PrototypeObjectTypeBase
	ObjectType
	EnsureCold() *PrototypeObjectType_ColdData
	GetOwnPropertyMapOrNull() *PropertyMap
	EnsurePropertyMap() *PropertyMap
	SetPrettyPrint(bool)
	IsPrettyPrint() bool
	SetImplicitPrototype(ObjectType)
	IsAnonymous() bool
	MatchRecordTypeConstraint(ObjectType)
}

func ToIface_PrototypeObjectType[T interface {
	PrototypeObjectType
	comparable
}](v T) PrototypeObjectType {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_PrototypeObjectType[T PrototypeObjectType](v T) PrototypeObjectType {
	var vi PrototypeObjectType = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.rhino.jstype.PrototypeObjectType
type PrototypeObjectTypeBase struct {
	ObjectTypeBase
	className  string                        `gson:"className"`
	properties *PropertyMap                  `gson:"properties"`
	cold       *PrototypeObjectType_ColdData `gson:"cold"`
}

func (this_ *PrototypeObjectTypeBase) BindSelf_(s any) {
	_ = s
	this_.ObjectTypeBase.BindSelf_(s)
}
func (this_ *PrototypeObjectTypeBase) IsJavaNil_() bool { return this_ == nil }
func (*PrototypeObjectTypeBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.PrototypeObjectType", "PrototypeObjectType"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.PrototypeObjectType", "com.google.javascript.rhino.jstype.PrototypeObjectType", func(v_ any) bool {
		return j2gort.InstanceOf[PrototypeObjectType](v_)
	})
}
func (b *PrototypeObjectTypeBase) GetPrototypeObjectTypeBase() *PrototypeObjectTypeBase { return b }
func (this_ *PrototypeObjectTypeBase) self_PrototypeObjectType_() PrototypeObjectType {
	if s := this_.JSTypeBase.self_; s != nil {
		return s.(PrototypeObjectType)
	}
	return nil
}
func (*PrototypeObjectTypeBase) Implements_rhino_PrototypeObjectType() {}
func (*PrototypeObjectTypeBase) Implements_rhino_ObjectType()          {}
func (*PrototypeObjectTypeBase) Implements_rhino_JSType()              {}

var static_PrototypeObjectType_TYPE_CLASS JSTypeClass

func StaticGet_PrototypeObjectType_TYPE_CLASS() JSTypeClass {
	clinit_PrototypeObjectType()
	return static_PrototypeObjectType_TYPE_CLASS
}
func StaticSet_PrototypeObjectType_TYPE_CLASS(v JSTypeClass) {
	clinit_PrototypeObjectType()
	static_PrototypeObjectType_TYPE_CLASS = v
}

var static_PrototypeObjectType_EMPTY_OBJECT_TYPE_LIST j2gort.Ext_com_google_common_collect_ImmutableList

func StaticGet_PrototypeObjectType_EMPTY_OBJECT_TYPE_LIST() j2gort.Ext_com_google_common_collect_ImmutableList {
	clinit_PrototypeObjectType()
	return static_PrototypeObjectType_EMPTY_OBJECT_TYPE_LIST
}
func StaticSet_PrototypeObjectType_EMPTY_OBJECT_TYPE_LIST(v j2gort.Ext_com_google_common_collect_ImmutableList) {
	clinit_PrototypeObjectType()
	static_PrototypeObjectType_EMPTY_OBJECT_TYPE_LIST = v
}
func StaticGet_PrototypeObjectType_PROTO_FLAG_NATIVE() int8              { return int8(1) }
func StaticGet_PrototypeObjectType_PROTO_FLAG_ANONYMOUS() int8           { return int8(2) }
func StaticGet_PrototypeObjectType_PROTO_FLAG_PRETTY_PRINT() int8        { return int8(4) }
func StaticGet_PrototypeObjectType_PROTO_IMPLICIT_KIND_SHIFT() int32     { return int32(3) }
func StaticGet_PrototypeObjectType_PROTO_NULL() int32                    { return int32(0) }
func StaticGet_PrototypeObjectType_PROTO_OBJECT() int32                  { return int32(1) }
func StaticGet_PrototypeObjectType_PROTO_FUNCTION() int32                { return int32(2) }
func StaticGet_PrototypeObjectType_PROTO_CUSTOM() int32                  { return int32(3) }
func StaticGet_PrototypeObjectType_PROTO_FLAG_RECORD_DECLARED() int8     { return int8(32) }
func StaticGet_PrototypeObjectType_PROTO_FLAG_RECORD_FROZEN() int8       { return int8(64) }
func StaticGet_PrototypeObjectType_MAX_PRETTY_PRINTED_PROPERTIES() int32 { return int32(10) }

var static_PrototypeObjectType_PROP_COMPARATOR j2gort.Ext_java_util_Comparator

func StaticGet_PrototypeObjectType_PROP_COMPARATOR() j2gort.Ext_java_util_Comparator {
	clinit_PrototypeObjectType()
	return static_PrototypeObjectType_PROP_COMPARATOR
}
func StaticSet_PrototypeObjectType_PROP_COMPARATOR(v j2gort.Ext_java_util_Comparator) {
	clinit_PrototypeObjectType()
	static_PrototypeObjectType_PROP_COMPARATOR = v
}

var clinitDone_PrototypeObjectType uint32
var clinitState_PrototypeObjectType j2gort.ClinitState

func clinit_PrototypeObjectType() {
	if atomic.LoadUint32(&clinitDone_PrototypeObjectType) != 0 {
		return
	}
	clinitSlow_PrototypeObjectType()
}

func clinitSlow_PrototypeObjectType() {
	if !j2gort.EnterClinit(&clinitState_PrototypeObjectType, &clinitDone_PrototypeObjectType) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_PrototypeObjectType, &clinitDone_PrototypeObjectType)
	static_PrototypeObjectType_TYPE_CLASS = StaticGet_JSTypeClass_PROTOTYPE_OBJECT()
	static_PrototypeObjectType_EMPTY_OBJECT_TYPE_LIST = j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	static_PrototypeObjectType_PROP_COMPARATOR = lambdaSingleton_prototype_object_type_1
}

func (this_ *PrototypeObjectTypeBase) initFields_PrototypeObjectType() {
	_ = this_
	this_.className = j2gort.NullString
	this_.properties = nil
	this_.cold = nil
}

// Java constructor: PrototypeObjectType(com.google.javascript.rhino.jstype.PrototypeObjectType.Builder<?> builder)
func NewPrototypeObjectType_1(builder PrototypeObjectType_Builder) PrototypeObjectType {
	this_ := &PrototypeObjectTypeBase{}
	this_.BindSelf_(this_)
	InitPrototypeObjectType_1(this_, builder)
	return this_
}

func InitPrototypeObjectType_1(this_ *PrototypeObjectTypeBase, builder PrototypeObjectType_Builder) {
	_ = builder
	InitPrototypeObjectType_2(this_, builder, int32(0))
}

// Java constructor: PrototypeObjectType(com.google.javascript.rhino.jstype.PrototypeObjectType.Builder<?> builder, int extraKindMask)
func NewPrototypeObjectType_2(builder PrototypeObjectType_Builder, extraKindMask int32) PrototypeObjectType {
	this_ := &PrototypeObjectTypeBase{}
	this_.BindSelf_(this_)
	InitPrototypeObjectType_2(this_, builder, extraKindMask)
	return this_
}

func InitPrototypeObjectType_2(this_ *PrototypeObjectTypeBase, builder PrototypeObjectType_Builder, extraKindMask int32) {
	_ = builder
	_ = extraKindMask
	initPrototypeObjectTypeFromBase_(this_, builder.GetPrototypeObjectType_BuilderBase(), extraKindMask)

}
func initPrototypeObjectTypeFromBase_(this_ *PrototypeObjectTypeBase, builder *PrototypeObjectType_BuilderBase, extraKindMask int32) {
	_ = builder
	_ = extraKindMask
	InitObjectType_3(&this_.ObjectTypeBase, builder.registry, builder.templateTypeMap, ((func() int32 {
		if builder.nativeType {
			return int32(33554432)
		}
		return int32(0)
	}()) | extraKindMask))
	this_.initFields_PrototypeObjectType()
	this_.className = builder.className
	var paramCount int32 = builder.templateParamCount
	_ = paramCount
	if !(paramCount >= int32(0)) {
		j2gort.CheckState(false)
	}
	if paramCount < int32(255) {
		this_.GetJSTypeBase().protoTemplateParamCountByte = int8(paramCount)
	} else {
		this_.GetJSTypeBase().protoTemplateParamCountByte = int8(-1)
		{
			r_5 := this_.EnsureCold()
			r_5.templateParamCountOverflow = paramCount
		}
	}
	var pFlags int8 = int8(int32(0))
	_ = pFlags
	if builder.nativeType {
		pFlags = int8((int32(pFlags) | int32(int8(1))))
	}
	if builder.anonymousType {
		pFlags = int8((int32(pFlags) | int32(int8(2))))
	}
	this_.GetJSTypeBase().protoFlags = pFlags
	if builder.nativeType || builder.implicitPrototype != nil && !builder.implicitPrototype.IsJavaNil_() {
		this_.SetImplicitPrototype(builder.implicitPrototype)
	} else {
		this_.SetImplicitPrototype(this_.GetJSTypeBase().registry.GetNativeObjectType(StaticGet_JSTypeNative_OBJECT_TYPE()))
	}
	if builder.anonymousType {
		if !(this_.className == j2gort.NullString) {
			j2gort.CheckState(false)
		}
	}
	j2gort.CheckNotNullPtr(this_.GetJSTypeBase().templateTypeMap)
	if !(this_.GetJSTypeBase().templateTypeMap.Size() >= paramCount) {
		j2gort.CheckState(false)
	}
	this_.GetJSTypeBase().registry.GetResolver().ResolveIfClosed(this_.JSTypeBase.self_, StaticGet_PrototypeObjectType_TYPE_CLASS())
}

// Java method: PrototypeObjectType.ensureCold()
func (this_ *PrototypeObjectTypeBase) EnsureCold() *PrototypeObjectType_ColdData {
	var c *PrototypeObjectType_ColdData = this_.cold
	_ = c
	if c == nil {
		c = NewPrototypeObjectType_ColdData()
		this_.cold = c
	}
	return c
}

// Java method: PrototypeObjectType.getTypeClass()
func (this_ *PrototypeObjectTypeBase) GetTypeClass() JSTypeClass {
	return StaticGet_PrototypeObjectType_TYPE_CLASS()
}

// Java static method: PrototypeObjectType.builder(com.google.javascript.rhino.jstype.JSTypeRegistry registry)
func PrototypeObjectType_S_Builder(registry *JSTypeRegistry) PrototypeObjectType_Builder {
	_ = registry
	return NewPrototypeObjectType_Builder(registry)
}

// Java method: PrototypeObjectType.getOwnPropertyMapOrNull()
func (this_ *PrototypeObjectTypeBase) GetOwnPropertyMapOrNull() *PropertyMap {
	return this_.properties
}

// Java method: PrototypeObjectType.ensurePropertyMap()
func (this_ *PrototypeObjectTypeBase) EnsurePropertyMap() *PropertyMap {
	var map_ *PropertyMap = this_.properties
	_ = map_
	if map_ == nil {
		map_ = NewPropertyMap()
		map_.SetParentSource(Upcast_ObjectType(this_.self_PrototypeObjectType_()))
		this_.properties = map_
	}
	return map_
}

// Java method: PrototypeObjectType.getPropertyMap()
func (this_ *PrototypeObjectTypeBase) GetPropertyMap() *PropertyMap {
	return this_.EnsurePropertyMap()
}

// Java method: PrototypeObjectType.defineProperty(java.lang.String name, com.google.javascript.rhino.jstype.JSType type, boolean inferred, com.google.javascript.rhino.Node propertyNode)
func (this_ *PrototypeObjectTypeBase) DefineProperty_4_string_JSType_bool_Node(name string, type_ JSType, inferred bool, propertyNode Node) bool {
	_ = name
	_ = type_
	_ = inferred
	_ = propertyNode
	if this_.HasOwnDeclaredProperty_1_string(name) {
		return false
	}
	var val_6 Property
	InitProperty_4_string_JSType_bool_Node(&val_6, name, type_, inferred, propertyNode)
	var newProp *Property = &val_6
	_ = newProp
	this_.EnsurePropertyMap().PutProperty_2_string_Property(name, newProp)
	return true
}

// Java method: PrototypeObjectType.defineProperty(com.google.javascript.rhino.jstype.Property.Key name, com.google.javascript.rhino.jstype.JSType type, boolean inferred, com.google.javascript.rhino.Node propertyNode)
func (this_ *PrototypeObjectTypeBase) DefineProperty_4_Property_Key_JSType_bool_Node(name Property_Key, type_ JSType, inferred bool, propertyNode Node) bool {
	_ = name
	_ = type_
	_ = inferred
	_ = propertyNode
	if this_.HasOwnDeclaredProperty_1_Property_Key(name) {
		return false
	}
	var val_7 Property
	InitProperty_4_Property_Key_JSType_bool_Node(&val_7, name, type_, inferred, propertyNode)
	var newProp *Property = &val_7
	_ = newProp
	this_.EnsurePropertyMap().PutProperty_2_Property_Key_Property(name, newProp)
	return true
}

// Java method: PrototypeObjectType.setPropertyJSDocInfo(java.lang.String propertyName, com.google.javascript.rhino.JSDocInfo info)
func (this_ *PrototypeObjectTypeBase) SetPropertyJSDocInfo_2_string_JSDocInfo(propertyName string, info *JSDocInfo) {
	_ = propertyName
	_ = info
	j2gort.CheckNotNullString(propertyName)
	if info != nil {
		if (this_.properties == nil) || (this_.properties.GetOwnProperty_1_string(propertyName) == nil) {
			this_.DefineInferredProperty_3_string_JSType_Node(propertyName, this_.GetPropertyType_1_string(propertyName), Node(0))
		}
		var property *Property
		_ = property
		if this_.properties != nil {
			property = this_.properties.GetOwnProperty_1_string(propertyName)
		} else {
			property = nil
		}
		if property != nil {
			property.SetJSDocInfo(info)
		}
	}
}

// Java method: PrototypeObjectType.setPropertyJSDocInfo(com.google.javascript.rhino.jstype.Property.Key propertyName, com.google.javascript.rhino.JSDocInfo info)
func (this_ *PrototypeObjectTypeBase) SetPropertyJSDocInfo_2_Property_Key_JSDocInfo(propertyName Property_Key, info *JSDocInfo) {
	_ = propertyName
	_ = info
	if info != nil {
		if (this_.properties == nil) || (this_.properties.GetOwnProperty_1_Property_Key(propertyName) == nil) {
			this_.DefineInferredProperty_3_Property_Key_JSType_Node(propertyName, this_.self_PrototypeObjectType_().GetPropertyType_1_Property_Key(propertyName), Node(0))
		}
		var property *Property
		_ = property
		if this_.properties != nil {
			property = this_.properties.GetOwnProperty_1_Property_Key(propertyName)
		} else {
			property = nil
		}
		if property != nil {
			property.SetJSDocInfo(info)
		}
	}
}

// Java method: PrototypeObjectType.setPropertyNode(java.lang.String propertyName, com.google.javascript.rhino.Node defSite)
func (this_ *PrototypeObjectTypeBase) SetPropertyNode_2_string_Node(propertyName string, defSite Node) {
	_ = propertyName
	_ = defSite
	j2gort.CheckNotNullString(propertyName)
	if this_.properties != nil {
		var property *Property = this_.properties.GetOwnProperty_1_string(propertyName)
		_ = property
		if property != nil {
			property.SetNode(defSite)
		}
	}
}

// Java method: PrototypeObjectType.setPropertyNode(com.google.javascript.rhino.jstype.Property.Key propertyName, com.google.javascript.rhino.Node defSite)
func (this_ *PrototypeObjectTypeBase) SetPropertyNode_2_Property_Key_Node(propertyName Property_Key, defSite Node) {
	_ = propertyName
	_ = defSite
	if this_.properties != nil {
		var property *Property = this_.properties.GetOwnProperty_1_Property_Key(propertyName)
		_ = property
		if property != nil {
			property.SetNode(defSite)
		}
	}
}

// Java method: PrototypeObjectType.matchesNumberContext()
func (this_ *PrototypeObjectTypeBase) MatchesNumberContext() bool {
	return ((((this_.JSTypeBase.self_.IsNumberObjectType() || this_.JSTypeBase.self_.IsDateType()) || this_.JSTypeBase.self_.IsBooleanObjectType()) || this_.JSTypeBase.self_.IsStringObjectType()) || this_.priv_PrototypeObjectType_hasOverriddenNativeProperty_9412("valueOf"))
}

// Java method: PrototypeObjectType.matchesStringContext()
func (this_ *PrototypeObjectTypeBase) MatchesStringContext() bool {
	return ((((((((this_.JSTypeBase.self_.IsTheObjectType() || this_.JSTypeBase.self_.IsStringObjectType()) || this_.JSTypeBase.self_.IsDateType()) || this_.JSTypeBase.self_.IsRegexpType()) || this_.JSTypeBase.self_.IsArrayType()) || this_.JSTypeBase.self_.IsNumberObjectType()) || this_.JSTypeBase.self_.IsBigIntObjectType()) || this_.JSTypeBase.self_.IsBooleanObjectType()) || this_.priv_PrototypeObjectType_hasOverriddenNativeProperty_9412("toString"))
}

// Java method: PrototypeObjectType.matchesSymbolContext()
func (this_ *PrototypeObjectTypeBase) MatchesSymbolContext() bool {
	return this_.JSTypeBase.self_.IsSymbolObjectType()
}

// Java method: PrototypeObjectType.hasOverriddenNativeProperty(java.lang.String propertyName)
func (this_ *PrototypeObjectTypeBase) priv_PrototypeObjectType_hasOverriddenNativeProperty_9412(propertyName string) bool {
	_ = propertyName
	if this_.IsNativeObjectType() {
		return false
	}
	var propertyType JSType = this_.GetPropertyType_1_string(propertyName)
	_ = propertyType
	var nativeType ObjectType
	_ = nativeType
	if this_.IsFunctionType() {
		nativeType = this_.GetJSTypeBase().registry.GetNativeObjectType(StaticGet_JSTypeNative_FUNCTION_PROTOTYPE())
	} else {
		nativeType = this_.GetJSTypeBase().registry.GetNativeObjectType(StaticGet_JSTypeNative_OBJECT_PROTOTYPE())
	}
	var nativePropertyType JSType = nativeType.GetObjectTypeBase().GetPropertyType_1_string(propertyName)
	_ = nativePropertyType
	return !(propertyType == nativePropertyType || ((propertyType == nil || propertyType.

		// Java method: PrototypeObjectType.matchesObjectContext()
		IsJavaNil_()) && (nativePropertyType == nil || nativePropertyType.IsJavaNil_())))
}

func (this_ *PrototypeObjectTypeBase) MatchesObjectContext() bool {
	return true
}

// Java method: PrototypeObjectType.appendTo(com.google.javascript.rhino.jstype.TypeStringBuilder sb)
func (this_ *PrototypeObjectTypeBase) AppendTo(sb *TypeStringBuilder) {
	_ = sb
	if this_.HasReferenceName() {
		sb.Append_1_string(func() string {
			if sb.IsForAnnotations() {
				return this_.GetNormalizedReferenceName()
			}
			return this_.self_PrototypeObjectType_().GetReferenceName()
		}())
		return
	}
	if (int32(this_.GetJSTypeBase().protoFlags) & int32(int8(4))) == int32(0) {
		sb.Append_1_string(func() string {
			if sb.IsForAnnotations() {
				return "?"
			}
			return "{...}"
		}())
		return
	}
	this_.GetJSTypeBase().protoFlags = int8((int32(this_.GetJSTypeBase().protoFlags) & ^(int32(int8(4)))))
	var properties j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_TreeSet_1_Ext_java_util_Comparator(StaticGet_PrototypeObjectType_PROP_COMPARATOR())
	_ = properties
	{
		var current ObjectType = Upcast_ObjectType(this_.self_PrototypeObjectType_())
		_ = current
		for ; current != nil && !current.IsJavaNil_(); current = current.GetImplicitPrototype() {
			if current.GetObjectTypeBase().IsNativeObjectType() || (properties.Size() > int32(10)) {
				break
			}
			{
				cur_8 := j2gort.NewIterCursor(j2gort.CheckNotNull(current.GetOwnPropertyNames()))
				for cur_8.HasNext() {
					stringKey := cur_8.NextStr()
					_ = stringKey
					var key *Property_StringKey = NewProperty_StringKey(stringKey)
					_ = key
					properties.Add_1(NewPrototypeObjectType_PropertyForPrettyPrinting(key, current.GetPropertyType_1_Property_Key(key)))
				}
			}
			{
				cur_9 := j2gort.NewIterCursor(j2gort.CheckNotNull(current.GetObjectTypeBase().GetOwnPropertyKnownSymbols()))
				for cur_9.HasNext() {
					rawElem_10 := cur_9.Next()
					symbolKey := j2gort.CastPtr[KnownSymbolType](rawElem_10)
					_ = symbolKey
					var key_2 *Property_SymbolKey = NewProperty_SymbolKey(symbolKey)
					_ = key_2
					properties.Add_1(NewPrototypeObjectType_PropertyForPrettyPrinting(key_2, current.GetPropertyType_1_Property_Key(key_2)))
				}
			}
		}
	}
	var multiline bool = (!(sb.IsForAnnotations()) && (properties.Size() > int32(1)))
	_ = multiline
	sb.Append_1_string("{").Indent(j2gort.NewLambda_Ext_java_lang_Runnable(func() {
		if multiline {
			sb.BreakLineAndIndent()
		}
		var i int32 = int32(0)
		_ = i
		{
			cur_11 := j2gort.NewIterCursor(properties)
			for cur_11.HasNext() {
				rawElem_12 := cur_11.Next()
				property := j2gort.CastPtr[PrototypeObjectType_PropertyForPrettyPrinting](rawElem_12)
				_ = property
				i = i + 1
				if !(sb.IsForAnnotations()) && (i > int32(10)) {
					sb.Append_1_string("...")
					break
				}
				property.priv_PropertyForPrettyPrinting_appendTo_9427(sb)
				if i < properties.Size() {
					sb.Append_1_string(",")
					if multiline {
						sb.BreakLineAndIndent()
					} else {
						sb.Append_1_string(" ")
					}
				}
			}
		}
	}))
	if multiline {
		sb.BreakLineAndIndent()
	}
	sb.Append_1_string("}")
	this_.GetJSTypeBase().protoFlags = int8((int32(this_.GetJSTypeBase().protoFlags) | int32(int8(4))))
}

// Java method: PrototypeObjectType.setPrettyPrint(boolean prettyPrint)
func (this_ *PrototypeObjectTypeBase) SetPrettyPrint(prettyPrint bool) {
	_ = prettyPrint
	if prettyPrint {
		this_.GetJSTypeBase().protoFlags = int8((int32(this_.GetJSTypeBase().protoFlags) | int32(int8(4))))
	} else {
		this_.GetJSTypeBase().protoFlags = int8((int32(this_.GetJSTypeBase().protoFlags) & ^(int32(int8(4)))))
	}
}

// Java method: PrototypeObjectType.isPrettyPrint()
func (this_ *PrototypeObjectTypeBase) IsPrettyPrint() bool {
	return ((int32(this_.GetJSTypeBase().protoFlags) & int32(int8(4))) != int32(0))
}

// Java method: PrototypeObjectType.getConstructor()
func (this_ *PrototypeObjectTypeBase) GetConstructor() *FunctionType {
	return nil
}

// Java method: PrototypeObjectType.getImplicitPrototypeFallback()
func (this_ *PrototypeObjectTypeBase) priv_PrototypeObjectType_getImplicitPrototypeFallback_9415() ObjectType {
	var kind int32 = ((javalang.Ushr32(int32(this_.GetJSTypeBase().protoFlags), int(int32(3)))) & int32(3))
	_ = kind
	switch kind {
	case int32(0):
		return nil
	case int32(1):
		return this_.GetJSTypeBase().registry.GetNativeObjectType(StaticGet_JSTypeNative_OBJECT_TYPE())
	case int32(2):
		return this_.GetJSTypeBase().registry.GetNativeObjectType(StaticGet_JSTypeNative_FUNCTION_PROTOTYPE())
	default:
		if this_.cold != nil {
			return this_.cold.customImplicitPrototype
		} else {
			return nil
		}
	}
	panic("unreachable")
}

// Java method: PrototypeObjectType.setImplicitPrototypeFallback(com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable ObjectType implicitPrototype)
func (this_ *PrototypeObjectTypeBase) priv_PrototypeObjectType_setImplicitPrototypeFallback_9416(implicitPrototype ObjectType) {
	_ = implicitPrototype
	var kind int32 = int32(0)
	_ = kind
	if implicitPrototype == nil || implicitPrototype.IsJavaNil_() {
		kind = int32(0)
		if this_.cold != nil {
			{
				r_13 := this_.cold
				r_13.customImplicitPrototype = nil
			}
		}
	} else if j2gort.IfaceEquals[ObjectType](implicitPrototype, this_.GetJSTypeBase().registry.GetNativeObjectType(StaticGet_JSTypeNative_OBJECT_TYPE())) {
		kind = int32(1)
		if this_.cold != nil {
			{
				r_14 := this_.cold
				r_14.customImplicitPrototype = nil
			}
		}
	} else if j2gort.IfaceEquals[ObjectType](implicitPrototype, this_.GetJSTypeBase().registry.GetNativeObjectType(StaticGet_JSTypeNative_FUNCTION_PROTOTYPE())) {
		kind = int32(2)
		if this_.cold != nil {
			{
				r_15 := this_.cold
				r_15.customImplicitPrototype = nil
			}
		}
	} else {
		kind = int32(3)
		{
			r_16 := this_.EnsureCold()
			r_16.customImplicitPrototype = implicitPrototype
		}
	}
	this_.GetJSTypeBase().protoFlags = int8(((int32(this_.GetJSTypeBase().protoFlags) & ^(javalang.Shl32(int32(3), int(int32(3))))) | (javalang.Shl32(kind, int(int32(3))))))
}

// Java method: PrototypeObjectType.getImplicitPrototype()
func (this_ *PrototypeObjectTypeBase) GetImplicitPrototype() ObjectType {
	return this_.priv_PrototypeObjectType_getImplicitPrototypeFallback_9415()
}

// Java method: PrototypeObjectType.setImplicitPrototype(com.google.javascript.rhino.jstype.ObjectType implicitPrototype)
func (this_ *PrototypeObjectTypeBase) SetImplicitPrototype(implicitPrototype ObjectType) {
	_ = implicitPrototype
	if !(!(this_.self_PrototypeObjectType_().HasCachedValues())) {
		j2gort.CheckState(false)
	}
	var hadPreviousPrototype bool = !(j2gort.IfaceEquals[ObjectType](this_.priv_PrototypeObjectType_getImplicitPrototypeFallback_9415(), nil))
	_ = hadPreviousPrototype
	this_.priv_PrototypeObjectType_setImplicitPrototypeFallback_9416(implicitPrototype)
	if hadPreviousPrototype {
		this_.EnsurePropertyMap().SetParentSource(Upcast_ObjectType(this_.self_PrototypeObjectType_()))
	}
	if implicitPrototype != nil && !implicitPrototype.IsJavaNil_() {
		this_.MaybeLoosenTypecheckingDueToForwardReferencedSupertype(Upcast_JSType(implicitPrototype))
	}
}

// Java method: PrototypeObjectType.getTemplateParamCount()
func (this_ *PrototypeObjectTypeBase) GetTemplateParamCount() int32 {
	var b int32 = (int32(this_.GetJSTypeBase().protoTemplateParamCountByte) & int32(255))
	_ = b
	if b < int32(255) {
		return b
	} else {
		return this_.cold.templateParamCountOverflow
	}
}

// Java method: PrototypeObjectType.getReferenceName()
func (this_ *PrototypeObjectTypeBase) GetReferenceName() string {
	if this_.className != j2gort.NullString {
		return this_.className
	} else if (this_.cold != nil) && (this_.cold.ownerFunction != nil) {
		return (j2gort.StrValString(this_.cold.ownerFunction.GetReferenceName()) + ".prototype")
	} else {
		return j2gort.NullString
	}
}

// Java method: PrototypeObjectType.isAnonymous()
func (this_ *PrototypeObjectTypeBase) IsAnonymous() bool {
	return ((int32(this_.GetJSTypeBase().protoFlags) & int32(int8(2))) != int32(0))
}

// Java method: PrototypeObjectType.setOwnerFunction(com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable FunctionType type)
func (this_ *PrototypeObjectTypeBase) SetOwnerFunction(type_ *FunctionType) {
	_ = type_
	var currentOwner *FunctionType
	_ = currentOwner
	if this_.cold != nil {
		currentOwner = this_.cold.ownerFunction
	} else {
		currentOwner = nil
	}
	if !((currentOwner == nil) || (type_ == nil)) {
		j2gort.CheckState(false)
	}
	if type_ != nil {
		{
			r_17 := this_.EnsureCold()
			r_17.ownerFunction = type_
		}
	} else if this_.cold != nil {
		{
			r_18 := this_.cold
			r_18.ownerFunction = nil
		}
	}
}

// Java method: PrototypeObjectType.getOwnerFunction()
func (this_ *PrototypeObjectTypeBase) GetOwnerFunction() *FunctionType {
	if this_.cold != nil {
		return this_.cold.ownerFunction
	} else {
		return nil
	}
}

// Java method: PrototypeObjectType.getCtorImplementedInterfaces()
func (this_ *PrototypeObjectTypeBase) GetCtorImplementedInterfaces() any {
	if this_.IsFunctionPrototypeType() {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.GetOwnerFunction().GetImplementedInterfaces())
	} else {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](StaticGet_PrototypeObjectType_EMPTY_OBJECT_TYPE_LIST())
	}
}

// Java method: PrototypeObjectType.getCtorExtendedInterfaces()
func (this_ *PrototypeObjectTypeBase) GetCtorExtendedInterfaces() any {
	if this_.IsFunctionPrototypeType() {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.GetOwnerFunction().GetExtendedInterfaces())
	} else {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](StaticGet_PrototypeObjectType_EMPTY_OBJECT_TYPE_LIST())
	}
}

// Java method: PrototypeObjectType.resolveInternal(com.google.javascript.rhino.ErrorReporter reporter)
func (this_ *PrototypeObjectTypeBase) ResolveInternal(reporter ErrorReporter) JSType {
	_ = reporter
	var implicitPrototype ObjectType = this_.self_PrototypeObjectType_().GetImplicitPrototype()
	_ = implicitPrototype
	if implicitPrototype != nil && !implicitPrototype.IsJavaNil_() {
		var resolvedProto ObjectType = j2gort.AsIface[ObjectType](implicitPrototype.GetObjectTypeBase().Resolve(reporter))
		_ = resolvedProto
		if !(j2gort.InstanceOf[*InstanceObjectType](this_.self_PrototypeObjectType_())) {
			this_.priv_PrototypeObjectType_setImplicitPrototypeFallback_9416(resolvedProto)
		}
	}
	if this_.properties != nil {
		this_.properties.ResolveOwnPropertyTypes(reporter)
	}
	return this_.JSTypeBase.self_
}

// Java method: PrototypeObjectType.matchConstraint(com.google.javascript.rhino.jstype.JSType constraint)
func (this_ *PrototypeObjectTypeBase) MatchConstraint(constraint JSType) {
	_ = constraint
	if this_.HasReferenceName() {
		return
	}
	base_constraint_ := constraint.GetJSTypeBase()
	if base_constraint_.IsRecordType() {
		this_.MatchRecordTypeConstraint(base_constraint_.ToObjectType())
	} else if base_constraint_.IsUnionType() {
		var union *UnionType = constraint.ToMaybeUnionType()
		_ = union
		{
			var i int32 = int32(0)
			_ = i
			var n int32 = union.GetAlternatesCount()
			_ = n
			for ; i < n; i = i + 1 {
				var alt JSType = union.GetAlternateAt(i)
				_ = alt
				if alt.GetJSTypeBase().IsRecordType() {
					this_.MatchRecordTypeConstraint(alt.GetJSTypeBase().ToObjectType())
				}
			}
		}
	}
}

// Java method: PrototypeObjectType.matchRecordTypeConstraint(com.google.javascript.rhino.jstype.ObjectType constraintObj)
func (this_ *PrototypeObjectTypeBase) MatchRecordTypeConstraint(constraintObj ObjectType) {
	_ = constraintObj
	if !(j2gort.IfaceEquals[ObjectType](constraintObj, this_.self_PrototypeObjectType_())) && !(j2gort.InstanceOf[*FunctionType](constraintObj)) {
		var src PrototypeObjectType = ObjectType_UnwrapToPrototypeObjectType(constraintObj)
		_ = src
		if src != nil && !src.IsJavaNil_() {
			var ownMap *PropertyMap = src.GetPrototypeObjectTypeBase().GetOwnPropertyMapOrNull()
			_ = ownMap
			if ownMap == nil {
				return
			}
			if ownMap.HasFlatProps() {
				{
					var i int32 = int32(0)
					_ = i
					var n int32 = ownMap.GetOwnPropertiesCount()
					_ = n
					for ; i < n; i = i + 1 {
						var prop string = ownMap.GetFlatPropAt(i).GetName()
						_ = prop
						var propType JSType = constraintObj.GetObjectTypeBase().GetPropertyType_1_string(prop)
						_ = propType
						var slot *Property = this_.GetSlot_1_string(prop)
						_ = slot
						if (slot == nil) || slot.IsTypeInferred() {
							var typeToInfer JSType = propType
							_ = typeToInfer
							if ((slot == nil) && !(this_.IsEmptyType())) && !(this_.IsUnknownType()) {
								typeToInfer = this_.GetNativeType(StaticGet_JSTypeNative_VOID_TYPE()).GetLeastSupertype(propType)
							}
							this_.DefineInferredProperty_3_string_JSType_Node(prop, typeToInfer, Node(0))
						}
					}
				}
				return
			}
		}
	}
	{
		cur_19 := j2gort.NewIterCursor(j2gort.CheckNotNull(constraintObj.GetOwnPropertyNames()))
		for cur_19.HasNext() {
			prop_2 := cur_19.NextStr()
			_ = prop_2
			var propType_2 JSType = constraintObj.GetObjectTypeBase().GetPropertyType_1_string(prop_2)
			_ = propType_2
			var slot_2 *Property = this_.GetSlot_1_string(prop_2)
			_ = slot_2
			if (slot_2 == nil) || slot_2.IsTypeInferred() {
				var typeToInfer_2 JSType = propType_2
				_ = typeToInfer_2
				if ((slot_2 == nil) && !(this_.IsEmptyType())) && !(this_.IsUnknownType()) {
					typeToInfer_2 = this_.GetNativeType(StaticGet_JSTypeNative_VOID_TYPE()).GetLeastSupertype(propType_2)
				}
				this_.DefineInferredProperty_3_string_JSType_Node(prop_2, typeToInfer_2, Node(0))
			}
		}
	}
}

// Java method: PrototypeObjectType.recursionUnsafeHashCode()
func (this_ *PrototypeObjectTypeBase) RecursionUnsafeHashCode() int32 {
	if this_.self_PrototypeObjectType_().IsStructuralType() {
		var propsHash int32
		_ = propsHash
		if this_.properties != nil {
			propsHash = this_.properties.HashCode()
		} else {
			propsHash = int32(0)
		}
		return ((int32(31) * (int32(31) + func() int32 {
			hs_20 := this_.className
			if hs_20 == j2gort.NullString {
				return 0
			}
			return javalang.StringHashCode(hs_20)
		}())) + propsHash)
	} else {
		return j2gort.IdentityHash(this_.JSTypeBase.self_)
	}
}

// Java class: com.google.javascript.rhino.jstype.PrototypeObjectType.ColdData
type PrototypeObjectType_ColdData struct {
	customImplicitPrototype    ObjectType                                         `gson:"customImplicitPrototype"`
	ownerFunction              *FunctionType                                      `gson:"ownerFunction"`
	prototypeSlot              *Property                                          `gson:"prototypeSlot"`
	googModuleId               string                                             `gson:"googModuleId"`
	implementedInterfaces      j2gort.Ext_com_google_common_collect_ImmutableList `gson:"implementedInterfaces"`
	extendedInterfaces         j2gort.Ext_com_google_common_collect_ImmutableList `gson:"extendedInterfaces"`
	canonicalRepresentation    *FunctionType                                      `gson:"canonicalRepresentation"`
	templateParamCountOverflow int32                                              `gson:"templateParamCountOverflow"`
	idHash_                    atomic.Int32
	constructorAmbiguity       int8             `gson:"constructorAmbiguity"`
	propAccess                 int8             `gson:"propAccess"`
	closurePrimitive           ClosurePrimitive `gson:"closurePrimitive"`
}

func (this_ *PrototypeObjectType_ColdData) BindSelf_(s any) {
	_ = s
}
func (this_ *PrototypeObjectType_ColdData) IsJavaNil_() bool { return this_ == nil }
func (*PrototypeObjectType_ColdData) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.PrototypeObjectType$ColdData", "ColdData"
}
func (this_ *PrototypeObjectType_ColdData) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PrototypeObjectType_ColdData) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.PrototypeObjectType.ColdData", "com.google.javascript.rhino.jstype.PrototypeObjectType$ColdData", func(v_ any) bool {
		return j2gort.InstanceOf[*PrototypeObjectType_ColdData](v_)
	})
}

func (this_ *PrototypeObjectType_ColdData) initFields_PrototypeObjectType_ColdData() {
	_ = this_
	this_.googModuleId = j2gort.NullString
	this_.implementedInterfaces = StaticGet_PrototypeObjectType_EMPTY_OBJECT_TYPE_LIST()
	this_.extendedInterfaces = StaticGet_PrototypeObjectType_EMPTY_OBJECT_TYPE_LIST()
}

// Java constructor: ColdData()
func NewPrototypeObjectType_ColdData() *PrototypeObjectType_ColdData {
	this_ := &PrototypeObjectType_ColdData{}
	InitPrototypeObjectType_ColdData(this_)
	return this_
}

func InitPrototypeObjectType_ColdData(this_ *PrototypeObjectType_ColdData) {
	this_.initFields_PrototypeObjectType_ColdData()
}

// Java polymorphic class interface: com.google.javascript.rhino.jstype.PrototypeObjectType.Builder
type PrototypeObjectType_Builder interface {
	Implements_rhino_PrototypeObjectType_Builder()
	IsJavaNil_() bool
	GetPrototypeObjectType_BuilderBase() *PrototypeObjectType_BuilderBase
	SetName(string) any
	SetImplicitPrototype(ObjectType) PrototypeObjectType_Builder
	SetNative(bool) PrototypeObjectType_Builder
	SetAnonymous(bool) PrototypeObjectType_Builder
	SetTemplateTypeMap(*TemplateTypeMap) PrototypeObjectType_Builder
	SetTemplateParamCount(int32) PrototypeObjectType_Builder
	CastThis() PrototypeObjectType_Builder
	Build() any
	ResetBase()
}

func ToIface_PrototypeObjectType_Builder[T interface {
	PrototypeObjectType_Builder
	comparable
}](v T) PrototypeObjectType_Builder {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_PrototypeObjectType_Builder[T PrototypeObjectType_Builder](v T) PrototypeObjectType_Builder {
	var vi PrototypeObjectType_Builder = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.rhino.jstype.PrototypeObjectType.Builder
type PrototypeObjectType_BuilderBase struct {
	self_              PrototypeObjectType_Builder
	registry           *JSTypeRegistry  `gson:"registry"`
	className          string           `gson:"className"`
	implicitPrototype  ObjectType       `gson:"implicitPrototype"`
	templateTypeMap    *TemplateTypeMap `gson:"templateTypeMap"`
	templateParamCount int32            `gson:"templateParamCount"`
	idHash_            atomic.Int32
	nativeType         bool `gson:"nativeType"`
	anonymousType      bool `gson:"anonymousType"`
}

func (this_ *PrototypeObjectType_BuilderBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(PrototypeObjectType_Builder)
}
func (this_ *PrototypeObjectType_BuilderBase) IsJavaNil_() bool { return this_ == nil }
func (*PrototypeObjectType_BuilderBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.PrototypeObjectType$Builder", "Builder"
}
func (this_ *PrototypeObjectType_BuilderBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.PrototypeObjectType.Builder", "com.google.javascript.rhino.jstype.PrototypeObjectType$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[PrototypeObjectType_Builder](v_)
	})
}
func (b *PrototypeObjectType_BuilderBase) GetPrototypeObjectType_BuilderBase() *PrototypeObjectType_BuilderBase {
	return b
}
func (*PrototypeObjectType_BuilderBase) Implements_rhino_PrototypeObjectType_Builder() {}

func (this_ *PrototypeObjectType_BuilderBase) initFields_PrototypeObjectType_Builder() {
	_ = this_
	this_.className = j2gort.NullString
}

// Java constructor: Builder(com.google.javascript.rhino.jstype.JSTypeRegistry registry)
func NewPrototypeObjectType_Builder(registry *JSTypeRegistry) PrototypeObjectType_Builder {
	this_ := &PrototypeObjectType_BuilderBase{}
	this_.BindSelf_(this_)
	InitPrototypeObjectType_Builder(this_, registry)
	return this_
}

func InitPrototypeObjectType_Builder(this_ *PrototypeObjectType_BuilderBase, registry *JSTypeRegistry) {
	_ = registry
	this_.initFields_PrototypeObjectType_Builder()
	this_.registry = registry
	this_.templateTypeMap = registry.GetEmptyTemplateTypeMap()
}

// Java method: Builder.setName(java.lang.String x)
func (this_ *PrototypeObjectType_BuilderBase) SetName(x string) any {
	_ = x
	this_.className = x
	return this_.CastThis()
}

// Java method: Builder.setImplicitPrototype(com.google.javascript.rhino.jstype.ObjectType x)
func (this_ *PrototypeObjectType_BuilderBase) SetImplicitPrototype(x ObjectType) PrototypeObjectType_Builder {
	_ = x
	this_.implicitPrototype = x
	return this_.CastThis()
}

// Java method: Builder.setNative(boolean x)
func (this_ *PrototypeObjectType_BuilderBase) SetNative(x bool) PrototypeObjectType_Builder {
	_ = x
	this_.nativeType = x
	return this_.CastThis()
}

// Java method: Builder.setAnonymous(boolean x)
func (this_ *PrototypeObjectType_BuilderBase) SetAnonymous(x bool) PrototypeObjectType_Builder {
	_ = x
	this_.anonymousType = x
	return this_.CastThis()
}

// Java method: Builder.setTemplateTypeMap(com.google.javascript.rhino.jstype.TemplateTypeMap x)
func (this_ *PrototypeObjectType_BuilderBase) SetTemplateTypeMap(x *TemplateTypeMap) PrototypeObjectType_Builder {
	_ = x
	this_.templateTypeMap = x
	return this_.CastThis()
}

// Java method: Builder.setTemplateParamCount(int x)
func (this_ *PrototypeObjectType_BuilderBase) SetTemplateParamCount(x int32) PrototypeObjectType_Builder {
	_ = x
	this_.templateParamCount = x
	return this_.CastThis()
}

// Java method: Builder.castThis()
func (this_ *PrototypeObjectType_BuilderBase) CastThis() PrototypeObjectType_Builder {
	return this_.self_
}

// Java method: Builder.build()
func (this_ *PrototypeObjectType_BuilderBase) Build() any {
	return NewPrototypeObjectType_1(this_.self_)
}

// Java method: Builder.resetBase()
func (this_ *PrototypeObjectType_BuilderBase) ResetBase() {
	this_.className = j2gort.NullString
	this_.implicitPrototype = nil
	this_.nativeType = false
	this_.anonymousType = false
	this_.templateTypeMap = this_.registry.GetEmptyTemplateTypeMap()
	this_.templateParamCount = int32(0)
}

// Java record: com.google.javascript.rhino.jstype.PrototypeObjectType.PropertyForPrettyPrinting
type PrototypeObjectType_PropertyForPrettyPrinting struct {
	key   Property_Key `gson:"key"`
	type_ JSType       `gson:"type"`
}

func (this_ *PrototypeObjectType_PropertyForPrettyPrinting) BindSelf_(s any) {
	_ = s
}
func (this_ *PrototypeObjectType_PropertyForPrettyPrinting) IsJavaNil_() bool { return this_ == nil }
func (*PrototypeObjectType_PropertyForPrettyPrinting) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.PrototypeObjectType$PropertyForPrettyPrinting", "PropertyForPrettyPrinting"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.PrototypeObjectType.PropertyForPrettyPrinting", "com.google.javascript.rhino.jstype.PrototypeObjectType$PropertyForPrettyPrinting", func(v_ any) bool {
		return j2gort.InstanceOf[*PrototypeObjectType_PropertyForPrettyPrinting](v_)
	})
}
func (*PrototypeObjectType_PropertyForPrettyPrinting) Implements_j2gort_Ext_java_lang_Record() {}

func (this_ *PrototypeObjectType_PropertyForPrettyPrinting) initFields_PrototypeObjectType_PropertyForPrettyPrinting() {
	_ = this_
}

// Java constructor: PropertyForPrettyPrinting(com.google.javascript.rhino.jstype.Property.Key key, com.google.javascript.rhino.jstype.JSType type)
func NewPrototypeObjectType_PropertyForPrettyPrinting(key Property_Key, type_ JSType) *PrototypeObjectType_PropertyForPrettyPrinting {
	this_ := &PrototypeObjectType_PropertyForPrettyPrinting{}
	InitPrototypeObjectType_PropertyForPrettyPrinting(this_, key, type_)
	return this_
}

func InitPrototypeObjectType_PropertyForPrettyPrinting(this_ *PrototypeObjectType_PropertyForPrettyPrinting, key Property_Key, type_ JSType) {
	_ = key
	_ = type_
	this_.initFields_PrototypeObjectType_PropertyForPrettyPrinting()
	this_.key = key
	this_.type_ = type_
}

// Java method: PropertyForPrettyPrinting.appendTo(com.google.javascript.rhino.jstype.TypeStringBuilder sb)
func (this_ *PrototypeObjectType_PropertyForPrettyPrinting) priv_PropertyForPrettyPrinting_appendTo_9427(sb *TypeStringBuilder) {
	_ = sb
	switch (this_.key.Kind()).Ordinal() {
	case int32(0) /* KeyKind.STRING */ :
		sb.Append_1_string(this_.key.HumanReadableName())
	case int32(1) /* KeyKind.SYMBOL */ :
		sb.Append_1_string("[")
		sb.Append_1_string(this_.key.HumanReadableName())
		sb.Append_1_string("]")
	}
	sb.Append_1_string(": ").AppendNonNull(this_.type_)
}

func (this_ *PrototypeObjectType_PropertyForPrettyPrinting) Key() Property_Key { return this_.key }
func (this_ *PrototypeObjectType_PropertyForPrettyPrinting) Type() JSType      { return this_.type_ }
func (this_ *PrototypeObjectType_PropertyForPrettyPrinting) Equals(other_ any) bool {
	o_, ok_ := other_.(*PrototypeObjectType_PropertyForPrettyPrinting)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if !javalang.ObjectsEquals(this_.key, o_.key) {
		return false
	}
	if !javalang.ObjectsEquals(this_.type_, o_.type_) {
		return false
	}
	return true
}
func (this_ *PrototypeObjectType_PropertyForPrettyPrinting) HashCode() int32 {
	var h_ int32 = 0
	h_ = 31*h_ + j2gort.ObjectHash(this_.key)
	{
		hi_21 := this_.type_
		h_ = 31 * h_
		if hi_21 != nil && !hi_21.IsJavaNil_() {
			h_ += hi_21.GetJSTypeBase().HashCode()
		}
	}
	return h_
}
func (this_ *PrototypeObjectType_PropertyForPrettyPrinting) String() string {
	_ = fmt.Sprint
	return "PropertyForPrettyPrinting[" + "key=" + j2gort.StrVal(this_.key) + ", " + "type=" + j2gort.StrVal(this_.type_) + "]"
}
