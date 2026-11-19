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
 *   Nick Santos
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */
// Port of com/google/javascript/rhino/jstype/PropertyMap.java (generated mechanically by j2go).

package rhino

import (
	"fmt"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.rhino.jstype.PropertyMap
type PropertyMap struct {
	parentSource ObjectType             `gson:"parentSource"`
	prop0        *Property              `gson:"prop0"`
	prop1        *Property              `gson:"prop1"`
	store        *PropertyMap_PropStore `gson:"store"`
	extra        *PropertyMap_Extra     `gson:"extra"`
}

func (this_ *PropertyMap) BindSelf_(s any) {
	_ = s
}
func (this_ *PropertyMap) IsJavaNil_() bool { return this_ == nil }
func (*PropertyMap) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.PropertyMap", "PropertyMap"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.PropertyMap", "com.google.javascript.rhino.jstype.PropertyMap", func(v_ any) bool {
		return j2gort.InstanceOf[*PropertyMap](v_)
	})
}

var static_PropertyMap_EMPTY_MAP *PropertyMap

func StaticGet_PropertyMap_EMPTY_MAP() *PropertyMap {
	clinit_PropertyMap()
	return static_PropertyMap_EMPTY_MAP
}
func StaticSet_PropertyMap_EMPTY_MAP(v *PropertyMap) {
	clinit_PropertyMap()
	static_PropertyMap_EMPTY_MAP = v
}

var static_PropertyMap_EMPTY_PROPERTY_NAMES j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_PropertyMap_EMPTY_PROPERTY_NAMES() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_PropertyMap()
	return static_PropertyMap_EMPTY_PROPERTY_NAMES
}
func StaticSet_PropertyMap_EMPTY_PROPERTY_NAMES(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_PropertyMap()
	static_PropertyMap_EMPTY_PROPERTY_NAMES = v
}

var static_PropertyMap_EMPTY_KNOWN_SYMBOLS j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_PropertyMap_EMPTY_KNOWN_SYMBOLS() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_PropertyMap()
	return static_PropertyMap_EMPTY_KNOWN_SYMBOLS
}
func StaticSet_PropertyMap_EMPTY_KNOWN_SYMBOLS(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_PropertyMap()
	static_PropertyMap_EMPTY_KNOWN_SYMBOLS = v
}

var static_PropertyMap_EMPTY_ALL_KEYS *PropertyMap_AllKeys

func StaticGet_PropertyMap_EMPTY_ALL_KEYS() *PropertyMap_AllKeys {
	clinit_PropertyMap()
	return static_PropertyMap_EMPTY_ALL_KEYS
}
func StaticSet_PropertyMap_EMPTY_ALL_KEYS(v *PropertyMap_AllKeys) {
	clinit_PropertyMap()
	static_PropertyMap_EMPTY_ALL_KEYS = v
}

var clinitDone_PropertyMap uint32
var clinitState_PropertyMap j2gort.ClinitState

func clinit_PropertyMap() {
	if atomic.LoadUint32(&clinitDone_PropertyMap) != 0 {
		return
	}
	clinitSlow_PropertyMap()
}

func clinitSlow_PropertyMap() {
	if !j2gort.EnterClinit(&clinitState_PropertyMap, &clinitDone_PropertyMap) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_PropertyMap, &clinitDone_PropertyMap)
	static_PropertyMap_EMPTY_MAP = NewPropertyMap()
	static_PropertyMap_EMPTY_PROPERTY_NAMES = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()
	static_PropertyMap_EMPTY_KNOWN_SYMBOLS = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()
	static_PropertyMap_EMPTY_ALL_KEYS = NewPropertyMap_AllKeys(j2gort.Ext_com_google_common_collect_ImmutableSortedSet_Of_0(), j2gort.AsIface[j2gort.Ext_java_util_Set](StaticGet_PropertyMap_EMPTY_KNOWN_SYMBOLS()))
}

func (this_ *PropertyMap) initFields_PropertyMap() {
	_ = this_

}

// Java constructor: PropertyMap()
func NewPropertyMap() *PropertyMap {
	clinit_PropertyMap()
	this_ := &PropertyMap{}
	InitPropertyMap(this_)
	return this_
}

func InitPropertyMap(this_ *PropertyMap) {
	this_.initFields_PropertyMap()
}

// Java method: PropertyMap.ensureExtra()
func (this_ *PropertyMap) priv_PropertyMap_ensureExtra_9380() *PropertyMap_Extra {
	var e *PropertyMap_Extra = this_.extra
	_ = e
	if e == nil {
		e = NewPropertyMap_Extra()
		this_.extra = e
	}
	return e
}

// Java static method: PropertyMap.immutableEmptyMap()
func PropertyMap_ImmutableEmptyMap() *PropertyMap {
	clinit_PropertyMap()
	return StaticGet_PropertyMap_EMPTY_MAP()
}

// Java method: PropertyMap.setParentSource(com.google.javascript.rhino.jstype.ObjectType ownerType)
func (this_ *PropertyMap) SetParentSource(ownerType ObjectType) {
	_ = ownerType
	if this_ == StaticGet_PropertyMap_EMPTY_MAP() {
		return
	}
	this_.parentSource = ownerType
	this_.priv_PropertyMap_incrementCachedKeySetCounter_9413()
}

// Java method: PropertyMap.setParentForTesting(com.google.javascript.rhino.jstype.PropertyMap parent)
func (this_ *PropertyMap) SetParentForTesting(parent *PropertyMap) {
	_ = parent
	this_.parentSource = parent.parentSource
	this_.priv_PropertyMap_incrementCachedKeySetCounter_9413()
}

// Java method: PropertyMap.getPrimaryParent()
func (this_ *PropertyMap) GetPrimaryParent() *PropertyMap {
	if this_.parentSource == nil || this_.parentSource.IsJavaNil_() {
		return nil
	}
	var iProto ObjectType = this_.parentSource.GetImplicitPrototype()
	_ = iProto
	if iProto == nil || iProto.IsJavaNil_() {
		return nil
	} else {
		return iProto.GetPropertyMap()
	}
}

// Java static method: PropertyMap.getSecondaryParentObjects(com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable ObjectType parentSource)
func PropertyMap_GetSecondaryParentObjects(parentSource ObjectType) j2gort.Ext_java_lang_Iterable {
	clinit_PropertyMap()
	_ = parentSource
	if parentSource == nil || parentSource.IsJavaNil_() {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
	}
	var ctor *FunctionType = parentSource.GetConstructor()
	_ = ctor
	if (ctor != nil) && ctor.IsAbstract() {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](ctor.GetOwnImplementedInterfaces())
	}
	return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](parentSource.GetCtorExtendedInterfaces())
}

// Java method: PropertyMap.getSecondaryParentObjects()
func (this_ *PropertyMap) priv_PropertyMap_getSecondaryParentObjects_9384() j2gort.Ext_java_lang_Iterable {
	return PropertyMap_GetSecondaryParentObjects(this_.parentSource)
}

// Java static method: PropertyMap.findClosest(com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable PrototypeObjectType start, java.lang.String name)
func PropertyMap_FindClosest_2_PrototypeObjectType_string(start PrototypeObjectType, name string) *Property_OwnedProperty {
	clinit_PropertyMap()
	_ = start
	_ = name
	{
		var curr PrototypeObjectType = start
		_ = curr
		for ; curr != nil && !curr.IsJavaNil_(); curr = ObjectType_UnwrapToPrototypeObjectType(curr.GetImplicitPrototype()) {
			var ownMap *PropertyMap = curr.GetPrototypeObjectTypeBase().GetOwnPropertyMapOrNull()
			_ = ownMap
			if ownMap != nil {
				var prop *Property = ownMap.GetOwnProperty_1_string(name)
				_ = prop
				if prop != nil {
					return NewProperty_OwnedProperty(Upcast_ObjectType(curr), prop)
				}
			}
		}
	}
	{
		var curr_2 PrototypeObjectType = start
		_ = curr_2
		for ; curr_2 != nil && !curr_2.IsJavaNil_(); curr_2 = ObjectType_UnwrapToPrototypeObjectType(curr_2.GetImplicitPrototype()) {
			{
				cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(PropertyMap_GetSecondaryParentObjects(Upcast_ObjectType(curr_2))))
				for cur_1.HasNext() {
					rawElem_2 := cur_1.Next()
					o := j2gort.AsIface[ObjectType](rawElem_2)
					_ = o
					var e *Property_OwnedProperty = PropertyMap_FindClosest_2_PrototypeObjectType_string(ObjectType_UnwrapToPrototypeObjectType(o), name)
					_ = e
					if e != nil {
						return e
					}
				}
			}
		}
	}
	return nil
}

// Java method: PropertyMap.findClosest(java.lang.String name)
func (this_ *PropertyMap) FindClosest_1_string(name string) *Property_OwnedProperty {
	_ = name
	if this_.parentSource == nil || this_.parentSource.IsJavaNil_() {
		var prop *Property = this_.GetOwnProperty_1_string(name)
		_ = prop
		if prop != nil {
			return NewProperty_OwnedProperty(nil, prop)
		} else {
			return nil
		}
	}
	return PropertyMap_FindClosest_2_PrototypeObjectType_string(ObjectType_UnwrapToPrototypeObjectType(this_.parentSource), name)
}

// Java static method: PropertyMap.findClosest(com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable PrototypeObjectType start, com.google.javascript.rhino.jstype.Property.Key name)
func PropertyMap_FindClosest_2_PrototypeObjectType_Property_Key(start PrototypeObjectType, name Property_Key) *Property_OwnedProperty {
	clinit_PropertyMap()
	_ = start
	_ = name
	if name.Kind() == StaticGet_Property_KeyKind_STRING() {
		return PropertyMap_FindClosest_2_PrototypeObjectType_string(start, name.M_String())
	}
	var symbol *KnownSymbolType = name.Symbol()
	_ = symbol
	{
		var curr PrototypeObjectType = start
		_ = curr
		for ; curr != nil && !curr.IsJavaNil_(); curr = ObjectType_UnwrapToPrototypeObjectType(curr.GetImplicitPrototype()) {
			var ownMap *PropertyMap = curr.GetPrototypeObjectTypeBase().GetOwnPropertyMapOrNull()
			_ = ownMap
			if ownMap != nil {
				var prop *Property = ownMap.priv_PropertyMap_getOwnSymbolProperty_9389(symbol)
				_ = prop
				if prop != nil {
					return NewProperty_OwnedProperty(Upcast_ObjectType(curr), prop)
				}
			}
		}
	}
	{
		var curr_2 PrototypeObjectType = start
		_ = curr_2
		for ; curr_2 != nil && !curr_2.IsJavaNil_(); curr_2 = ObjectType_UnwrapToPrototypeObjectType(curr_2.GetImplicitPrototype()) {
			{
				cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(PropertyMap_GetSecondaryParentObjects(Upcast_ObjectType(curr_2))))
				for cur_3.HasNext() {
					rawElem_4 := cur_3.Next()
					o := j2gort.AsIface[ObjectType](rawElem_4)
					_ = o
					var e *Property_OwnedProperty = PropertyMap_FindClosest_2_PrototypeObjectType_Property_Key(ObjectType_UnwrapToPrototypeObjectType(o), name)
					_ = e
					if e != nil {
						return e
					}
				}
			}
		}
	}
	return nil
}

// Java method: PropertyMap.findClosest(com.google.javascript.rhino.jstype.Property.Key name)
func (this_ *PropertyMap) FindClosest_1_Property_Key(name Property_Key) *Property_OwnedProperty {
	_ = name
	if this_.parentSource == nil || this_.parentSource.IsJavaNil_() {
		var prop *Property = this_.GetOwnProperty_1_Property_Key(name)
		_ = prop
		if prop != nil {
			return NewProperty_OwnedProperty(nil, prop)
		} else {
			return nil
		}
	}
	return PropertyMap_FindClosest_2_PrototypeObjectType_Property_Key(ObjectType_UnwrapToPrototypeObjectType(this_.parentSource), name)
}

// Java static method: PropertyMap.findClosestProperty(com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable PrototypeObjectType start, java.lang.String name)
func PropertyMap_FindClosestProperty_2_PrototypeObjectType_string(start PrototypeObjectType, name string) *Property {
	clinit_PropertyMap()
	_ = start
	_ = name
	{
		var curr PrototypeObjectType = start
		_ = curr
		for ; curr != nil && !curr.IsJavaNil_(); curr = ObjectType_UnwrapToPrototypeObjectType(curr.GetImplicitPrototype()) {
			var ownMap *PropertyMap = curr.GetPrototypeObjectTypeBase().GetOwnPropertyMapOrNull()
			_ = ownMap
			if ownMap != nil {
				var prop *Property = ownMap.GetOwnProperty_1_string(name)
				_ = prop
				if prop != nil {
					return prop
				}
			}
		}
	}
	{
		var curr_2 PrototypeObjectType = start
		_ = curr_2
		for ; curr_2 != nil && !curr_2.IsJavaNil_(); curr_2 = ObjectType_UnwrapToPrototypeObjectType(curr_2.GetImplicitPrototype()) {
			{
				cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(PropertyMap_GetSecondaryParentObjects(Upcast_ObjectType(curr_2))))
				for cur_5.HasNext() {
					rawElem_6 := cur_5.Next()
					o := j2gort.AsIface[ObjectType](rawElem_6)
					_ = o
					var prop_2 *Property = PropertyMap_FindClosestProperty_2_PrototypeObjectType_string(ObjectType_UnwrapToPrototypeObjectType(o), name)
					_ = prop_2
					if prop_2 != nil {
						return prop_2
					}
				}
			}
		}
	}
	return nil
}

// Java method: PropertyMap.findClosestProperty(java.lang.String name)
func (this_ *PropertyMap) FindClosestProperty_1_string(name string) *Property {
	_ = name
	if this_.parentSource == nil || this_.parentSource.IsJavaNil_() {
		return this_.GetOwnProperty_1_string(name)
	}
	return PropertyMap_FindClosestProperty_2_PrototypeObjectType_string(ObjectType_UnwrapToPrototypeObjectType(this_.parentSource), name)
}

// Java static method: PropertyMap.findClosestProperty(com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable PrototypeObjectType start, com.google.javascript.rhino.jstype.Property.Key name)
func PropertyMap_FindClosestProperty_2_PrototypeObjectType_Property_Key(start PrototypeObjectType, name Property_Key) *Property {
	clinit_PropertyMap()
	_ = start
	_ = name
	if name.Kind() == StaticGet_Property_KeyKind_STRING() {
		return PropertyMap_FindClosestProperty_2_PrototypeObjectType_string(start, name.M_String())
	}
	var symbol *KnownSymbolType = name.Symbol()
	_ = symbol
	{
		var curr PrototypeObjectType = start
		_ = curr
		for ; curr != nil && !curr.IsJavaNil_(); curr = ObjectType_UnwrapToPrototypeObjectType(curr.GetImplicitPrototype()) {
			var ownMap *PropertyMap = curr.GetPrototypeObjectTypeBase().GetOwnPropertyMapOrNull()
			_ = ownMap
			if ownMap != nil {
				var prop *Property = ownMap.priv_PropertyMap_getOwnSymbolProperty_9389(symbol)
				_ = prop
				if prop != nil {
					return prop
				}
			}
		}
	}
	{
		var curr_2 PrototypeObjectType = start
		_ = curr_2
		for ; curr_2 != nil && !curr_2.IsJavaNil_(); curr_2 = ObjectType_UnwrapToPrototypeObjectType(curr_2.GetImplicitPrototype()) {
			{
				cur_7 := j2gort.NewIterCursor(j2gort.CheckNotNull(PropertyMap_GetSecondaryParentObjects(Upcast_ObjectType(curr_2))))
				for cur_7.HasNext() {
					rawElem_8 := cur_7.Next()
					o := j2gort.AsIface[ObjectType](rawElem_8)
					_ = o
					var prop_2 *Property = PropertyMap_FindClosestProperty_2_PrototypeObjectType_Property_Key(ObjectType_UnwrapToPrototypeObjectType(o), name)
					_ = prop_2
					if prop_2 != nil {
						return prop_2
					}
				}
			}
		}
	}
	return nil
}

// Java method: PropertyMap.findClosestProperty(com.google.javascript.rhino.jstype.Property.Key name)
func (this_ *PropertyMap) FindClosestProperty_1_Property_Key(name Property_Key) *Property {
	_ = name
	if this_.parentSource == nil || this_.parentSource.IsJavaNil_() {
		return this_.GetOwnProperty_1_Property_Key(name)
	}
	return PropertyMap_FindClosestProperty_2_PrototypeObjectType_Property_Key(ObjectType_UnwrapToPrototypeObjectType(this_.parentSource), name)
}

// Java method: PropertyMap.getOwnSymbolProperty(com.google.javascript.rhino.jstype.KnownSymbolType symbol)
func (this_ *PropertyMap) priv_PropertyMap_getOwnSymbolProperty_9389(symbol *KnownSymbolType) *Property {
	_ = symbol
	var e *PropertyMap_Extra = this_.extra
	_ = e
	if (e != nil) && (e.knownSymbols != nil) {
		return e.knownSymbols.Get(symbol)
	} else {
		return nil
	}
}

// Java method: PropertyMap.getOwnProperty(java.lang.String propertyName)
func (this_ *PropertyMap) GetOwnProperty_1_string(propertyName string) *Property {
	_ = propertyName
	var s *PropertyMap_PropStore = this_.store
	_ = s
	if s == nil {
		var p0 *Property = this_.prop0
		_ = p0
		if p0 != nil {
			var n0 string = p0.GetName()
			_ = n0
			if (n0 == propertyName) || (j2gort.CheckNotNullString(n0) == propertyName) {
				return p0
			}
			var p1 *Property = this_.prop1
			_ = p1
			if p1 != nil {
				var n1 string = p1.GetName()
				_ = n1
				if (n1 == propertyName) || (j2gort.CheckNotNullString(n1) == propertyName) {
					return p1
				}
			}
		}
		return nil
	}
	if s.map_ == nil {
		var items []*Property = s.items
		_ = items
		{
			var i int32 = int32(0)
			_ = i
			var n int32 = s.size
			_ = n
			for ; i < n; i = i + 1 {
				var p *Property = items[i]
				_ = p
				var pn string = p.GetName()
				_ = pn
				if (pn == propertyName) || (j2gort.CheckNotNullString(pn) == propertyName) {
					return p
				}
			}
		}
		return nil
	}
	return s.map_.Get(propertyName)
}

// Java method: PropertyMap.getOwnProperty(com.google.javascript.rhino.jstype.Property.Key propertyName)
func (this_ *PropertyMap) GetOwnProperty_1_Property_Key(propertyName Property_Key) *Property {
	_ = propertyName
	switch (propertyName.Kind()).Ordinal() {
	case int32(0) /* KeyKind.STRING */ :
		return this_.GetOwnProperty_1_string(propertyName.M_String())
	case int32(1) /* KeyKind.SYMBOL */ :
		return this_.priv_PropertyMap_getOwnSymbolProperty_9389(propertyName.Symbol())
	}
	panic("unreachable switch expression")
}

// Java method: PropertyMap.getPropertiesCount()
func (this_ *PropertyMap) GetPropertiesCount() int32 {
	var primaryParent *PropertyMap = this_.GetPrimaryParent()
	_ = primaryParent
	if primaryParent == nil {
		return this_.GetOwnPropertiesCount()
	}
	return this_.GetAllKeys().StringKeys().Size()
}

// Java method: PropertyMap.hasInlinePropsOnly()
func (this_ *PropertyMap) HasInlinePropsOnly() bool {
	return (this_.store == nil)
}

// Java method: PropertyMap.hasFlatProps()
func (this_ *PropertyMap) HasFlatProps() bool {
	var s *PropertyMap_PropStore = this_.store
	_ = s
	return ((s == nil) || (s.map_ == nil))
}

// Java method: PropertyMap.getFlatPropAt(int index)
func (this_ *PropertyMap) GetFlatPropAt(index int32) *Property {
	_ = index
	var s *PropertyMap_PropStore = this_.store
	_ = s
	if s != nil {
		return s.items[index]
	}
	if index == int32(0) {
		return this_.prop0
	} else {
		return this_.prop1
	}
}

// Java method: PropertyMap.getInlineProp0()
func (this_ *PropertyMap) GetInlineProp0() *Property {
	return this_.prop0
}

// Java method: PropertyMap.getInlineProp1()
func (this_ *PropertyMap) GetInlineProp1() *Property {
	return this_.prop1
}

// Java method: PropertyMap.hasOwnKnownSymbols()
func (this_ *PropertyMap) HasOwnKnownSymbols() bool {
	var e *PropertyMap_Extra = this_.extra
	_ = e
	return (((e != nil) && (e.knownSymbols != nil)) && !(e.knownSymbols.IsEmpty()))
}

// Java method: PropertyMap.getOwnPropertiesCount()
func (this_ *PropertyMap) GetOwnPropertiesCount() int32 {
	var s *PropertyMap_PropStore = this_.store
	_ = s
	if s != nil {
		if s.map_ != nil {
			return s.map_.Size()
		} else {
			return s.size
		}
	}
	if this_.prop0 == nil {
		return int32(0)
	}
	if this_.prop1 == nil {
		return int32(1)
	} else {
		return int32(2)
	}
}

// Java method: PropertyMap.collectOwnPropertyNames(java.lang.String[] out, int idx, boolean skipPrototype)
func (this_ *PropertyMap) CollectOwnPropertyNames(out []string, idx int32, skipPrototype bool) int32 {
	_ = out
	_ = idx
	_ = skipPrototype
	var s *PropertyMap_PropStore = this_.store
	_ = s
	if s != nil {
		if s.map_ != nil {
			{
				cur_9 := s.map_.IterKeys()
				for cur_9.HasNext() {
					key := cur_9.
						Next()
					_ = key
					if !(skipPrototype) || !("prototype" == key) {
						{
							a_10 := out
							i_11 := func() int32 { old_12 := idx; nw_13 := old_12 + 1; idx = nw_13; return old_12 }()
							a_10[i_11] = key
						}
					}
				}
			}
			return idx
		}
		var items []*Property = s.items
		_ = items
		{
			var i int32 = int32(0)
			_ = i
			var n int32 = s.size
			_ = n
			for ; i < n; i = i + 1 {
				var key_2 string = items[i].GetName()
				_ = key_2
				if !(skipPrototype) || !("prototype" == key_2) {
					{
						a_14 := out
						i_15 := func() int32 { old_16 := idx; nw_17 := old_16 + 1; idx = nw_17; return old_16 }()
						a_14[i_15] = key_2
					}
				}
			}
		}
		return idx
	}
	var p0 *Property = this_.prop0
	_ = p0
	if p0 != nil {
		var k0 string = p0.GetName()
		_ = k0
		if !(skipPrototype) || !("prototype" == k0) {
			{
				a_18 := out
				i_19 := func() int32 { old_20 := idx; nw_21 := old_20 + 1; idx = nw_21; return old_20 }()
				a_18[i_19] = k0
			}
		}
		var p1 *Property = this_.prop1
		_ = p1
		if p1 != nil {
			var k1 string = p1.GetName()
			_ = k1
			if !(skipPrototype) || !("prototype" == k1) {
				{
					a_22 := out
					i_23 := func() int32 { old_24 := idx; nw_25 := old_24 + 1; idx = nw_25; return old_24 }()
					a_22[i_23] = k1
				}
			}
		}
	}
	return idx
}

// Java method: PropertyMap.addOwnPropertyNamesTo(java.util.Set<java.lang.String> out)
func (this_ *PropertyMap) AddOwnPropertyNamesTo(out j2gort.Ext_java_util_Set) {
	_ = out
	var s *PropertyMap_PropStore = this_.store
	_ = s
	if s != nil {
		if s.map_ != nil {
			out.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](s.map_.LiveKeySet()))
			return
		}
		var items []*Property = s.items
		_ = items
		{
			var i int32 = int32(0)
			_ = i
			var n int32 = s.size
			_ = n
			for ; i < n; i = i + 1 {
				j2gort.SetAddStr(out, items[i].GetName())
			}
		}
		return
	}
	var p0 *Property = this_.prop0
	_ = p0
	if p0 != nil {
		j2gort.SetAddStr(out, p0.GetName())
		var p1 *Property = this_.prop1
		_ = p1
		if p1 != nil {
			j2gort.SetAddStr(out, p1.GetName())
		}
	}
}

// Java method: PropertyMap.getOwnPropertyNames()
func (this_ *PropertyMap) GetOwnPropertyNames() j2gort.Ext_java_util_Set {
	var s *PropertyMap_PropStore = this_.store
	_ = s
	if s != nil {
		if s.map_ != nil {
			return j2gort.AsIface[j2gort.Ext_java_util_Set](s.map_.LiveKeySet())
		}
		var n int32 = s.size
		_ = n
		var items []*Property = s.items
		_ = items
		var builder j2gort.Ext_com_google_common_collect_ImmutableSet_Builder = j2gort.Ext_com_google_common_collect_ImmutableSet_BuilderWithExpectedSize(n)
		_ = builder
		{
			var i int32 = int32(0)
			_ = i
			for ; i < n; i = i + 1 {
				j2gort.ImmSetBuilderAddStrVoid(builder, items[i].GetName())
			}
		}
		return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](builder.Build())
	}
	var p0 *Property = this_.prop0
	_ = p0
	if p0 == nil {
		return j2gort.AsIface[j2gort.Ext_java_util_Set](StaticGet_PropertyMap_EMPTY_PROPERTY_NAMES())
	}
	var p1 *Property = this_.prop1
	_ = p1
	if p1 == nil {
		return j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.ImmSetOfStr1(p0.GetName()))
	}
	return j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.ImmSetOfStr2(p0.GetName(), p1.GetName()))
}

// Java method: PropertyMap.getOwnKnownSymbols()
func (this_ *PropertyMap) GetOwnKnownSymbols() j2gort.Ext_java_util_Set {
	var e *PropertyMap_Extra = this_.extra
	_ = e
	if (e != nil) && (e.knownSymbols != nil) {
		return j2gort.AsIface[j2gort.Ext_java_util_Set](e.knownSymbols.LiveKeySet())
	} else {
		return j2gort.AsIface[j2gort.Ext_java_util_Set](StaticGet_PropertyMap_EMPTY_KNOWN_SYMBOLS())
	}
}

// Java method: PropertyMap.getAllKeys()
func (this_ *PropertyMap) GetAllKeys() *PropertyMap_AllKeys {
	if (this_ == StaticGet_PropertyMap_EMPTY_MAP()) || ((((this_.parentSource == nil || this_.parentSource.IsJavaNil_()) && (this_.prop0 == nil)) && (this_.store == nil)) && ((this_.extra == nil) || (this_.extra.knownSymbols == nil))) {
		return StaticGet_PropertyMap_EMPTY_ALL_KEYS()
	}
	var ancestors j2gort.OrderedSet[*PropertyMap] = j2gort.NewOrderedSet[*PropertyMap]()
	_ = ancestors
	this_.priv_PropertyMap_collectAllAncestors_9405(&ancestors)
	var maxAncestorCounter int32 = int32(1)
	_ = maxAncestorCounter
	var ancestorHash int32 = int32(0)
	_ = ancestorHash
	var missingAncestorExtra bool = false
	_ = missingAncestorExtra
	{
		cur_26 := ancestors.Iter()
		for cur_26.HasNext() {

			ancestor := cur_26.
				Next()
			_ = ancestor
			if ancestor == StaticGet_PropertyMap_EMPTY_MAP() {
				continue
			}
			ancestorHash = ((int32(31) * ancestorHash) + j2gort.IdentityHash(j2gort.PtrToAny(ancestor)))
			var ae *PropertyMap_Extra = ancestor.extra
			_ = ae
			if ae == nil {
				missingAncestorExtra = true
			} else if ae.cachedKeySetCounter > maxAncestorCounter {
				maxAncestorCounter = ae.cachedKeySetCounter
			}
		}
	}
	var selfExtra *PropertyMap_Extra = this_.extra
	_ = selfExtra
	if ((((selfExtra == nil) || missingAncestorExtra) || (selfExtra.cachedKeySet == nil)) || (maxAncestorCounter != selfExtra.cachedKeySetCounter)) || (ancestorHash != selfExtra.cachedAncestorHash) {
		var keys *j2gort.Ext_java_util_TreeSet = j2gort.NewExt_java_util_TreeSet()
		_ = keys
		var knownSymbolsKeys *j2gort.Ext_java_util_LinkedHashSet = nil
		_ = knownSymbolsKeys
		{
			cur_28 := ancestors.Iter()
			for cur_28.HasNext() {

				ancestor_2 := cur_28.
					Next()
				_ = ancestor_2
				if ancestor_2 == StaticGet_PropertyMap_EMPTY_MAP() {
					continue
				}
				var ae_2 *PropertyMap_Extra = ancestor_2.priv_PropertyMap_ensureExtra_9380()
				_ = ae_2
				{
					r_30 := ae_2
					r_30.cachedKeySetCounter = maxAncestorCounter
				}
				{
					r_31 := ae_2
					r_31.cachedKeySet = nil
				}
				{
					r_32 := ae_2
					r_32.cachedKnownSymbolsKeySet = nil
				}
				ancestor_2.AddOwnPropertyNamesTo(keys)
				if (ae_2.knownSymbols != nil) && !(ae_2.knownSymbols.IsEmpty()) {
					if knownSymbolsKeys == nil {
						knownSymbolsKeys = j2gort.NewExt_java_util_LinkedHashSet()
					}
					knownSymbolsKeys.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](ae_2.knownSymbols.LiveKeySet()))
				}
			}
		}
		selfExtra = this_.priv_PropertyMap_ensureExtra_9380()
		{
			r_33 := selfExtra
			r_33.cachedKeySetCounter = maxAncestorCounter
		}
		{
			r_34 := selfExtra
			r_34.cachedAncestorHash = ancestorHash
		}
		{
			r_35 := selfExtra
			r_35.cachedKeySet = j2gort.Ext_com_google_common_collect_ImmutableSortedSet_CopyOfSorted(keys)
		}
		{
			r_36 := selfExtra
			r_36.cachedKnownSymbolsKeySet = func() j2gort.Ext_com_google_common_collect_ImmutableSet {
				if knownSymbolsKeys != nil {
					return j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](knownSymbolsKeys))
				}
				return StaticGet_PropertyMap_EMPTY_KNOWN_SYMBOLS()
			}()
		}
	}
	return NewPropertyMap_AllKeys(selfExtra.cachedKeySet, j2gort.AsIface[j2gort.Ext_java_util_Set](selfExtra.cachedKnownSymbolsKeySet))
}

// Java method: PropertyMap.collectAllAncestors(java.util.LinkedHashSet<com.google.javascript.rhino.jstype.PropertyMap> ancestors)
func (this_ *PropertyMap) priv_PropertyMap_collectAllAncestors_9405(ancestors *j2gort.OrderedSet[*PropertyMap]) {
	_ = ancestors
	if !(ancestors.Add(this_)) {
		return
	}
	var primaryParent *PropertyMap = this_.GetPrimaryParent()
	_ = primaryParent
	if primaryParent != nil {
		primaryParent.priv_PropertyMap_collectAllAncestors_9405(ancestors)
	}
	{
		cur_37 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.priv_PropertyMap_getSecondaryParentObjects_9384()))
		for cur_37.HasNext() {
			rawElem_38 := cur_37.Next()
			parentType := j2gort.AsIface[ObjectType](rawElem_38)
			_ = parentType
			var parentMap *PropertyMap = parentType.GetPropertyMap()
			_ = parentMap
			if parentMap != nil {
				parentMap.priv_PropertyMap_collectAllAncestors_9405(ancestors)
			}
		}
	}
}

// Java method: PropertyMap.putProperty(java.lang.String name, com.google.javascript.rhino.jstype.Property newProp)
func (this_ *PropertyMap) PutProperty_2_string_Property(name string, newProp *Property) {
	_ = name
	_ = newProp
	var s *PropertyMap_PropStore = this_.store
	_ = s
	if s == nil {
		var p0 *Property = this_.prop0
		_ = p0
		if p0 == nil {
			this_.prop0 = newProp
			this_.priv_PropertyMap_incrementCachedKeySetCounter_9413()
			return
		}
		var n0 string = p0.GetName()
		_ = n0
		if (n0 == name) || (j2gort.CheckNotNullString(n0) == name) {
			newProp.SetJSDocInfo(p0.GetJSDocInfo())
			this_.prop0 = newProp
			return
		}
		var p1 *Property = this_.prop1
		_ = p1
		if p1 == nil {
			this_.prop1 = newProp
			this_.priv_PropertyMap_incrementCachedKeySetCounter_9413()
			return
		}
		var n1 string = p1.GetName()
		_ = n1
		if (n1 == name) || (j2gort.CheckNotNullString(n1) == name) {
			newProp.SetJSDocInfo(p1.GetJSDocInfo())
			this_.prop1 = newProp
			return
		}
		this_.store = NewPropertyMap_PropStore(p0, p1, newProp)
		this_.prop0 = nil
		this_.prop1 = nil
		this_.priv_PropertyMap_incrementCachedKeySetCounter_9413()
		return
	}
	if s.map_ == nil {
		var items []*Property = s.items
		_ = items
		var n int32 = s.size
		_ = n
		{
			var i int32 = int32(0)
			_ = i
			for ; i < n; i = i + 1 {
				var existing *Property = items[i]
				_ = existing
				var en string = existing.GetName()
				_ = en
				if (en == name) || (j2gort.CheckNotNullString(en) == name) {
					newProp.SetJSDocInfo(existing.GetJSDocInfo())
					{
						a_39 := items
						i_40 := i
						a_39[i_40] = newProp
					}
					return
				}
			}
		}
		if n < int32(16) {
			if n == int32(len(items)) {
				items = func() []*Property {
					arrSrc_41 := items
					arrLen_42 := (n * int32(2))
					j2gort.CheckNotNullSlice(arrSrc_41)
					if arrLen_42 < 0 {
						panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
					}
					arrOut_43 := make([]*Property, arrLen_42)
					copy(arrOut_43, arrSrc_41)
					return arrOut_43
				}()
				{
					r_44 := s
					r_44.items = items
				}
			}
			{
				a_45 := items
				i_46 := n
				a_45[i_46] = newProp
			}
			{
				r_47 := s
				r_47.size = (n + int32(1))
			}
			this_.priv_PropertyMap_incrementCachedKeySetCounter_9413()
			return
		}
		{
			r_48 := s
			r_48.map_ = j2gort.NewOrderedMapPtrCap[string, *Property](int32(32))
		}
		{
			var i_2 int32 = int32(0)
			_ = i_2
			for ; i_2 < int32(16); i_2 = i_2 + 1 {
				var existing_2 *Property = items[i_2]
				_ = existing_2
				s.map_.Put(existing_2.GetName(), existing_2)
			}
		}
		s.map_.Put(name, newProp)
		{
			r_49 := s
			r_49.items = nil
		}
		{
			r_50 := s
			r_50.size = int32(0)
		}
		this_.priv_PropertyMap_incrementCachedKeySetCounter_9413()
		return
	}
	var oldProp *Property = s.map_.Get(name)
	_ = oldProp
	if oldProp == nil {
		this_.priv_PropertyMap_incrementCachedKeySetCounter_9413()
	} else {
		newProp.SetJSDocInfo(oldProp.GetJSDocInfo())
	}
	s.map_.Put(name, newProp)
}

// Java method: PropertyMap.removeProperty(java.lang.String name)
func (this_ *PropertyMap) RemoveProperty(name string) bool {
	_ = name
	var s *PropertyMap_PropStore = this_.store
	_ = s
	if s == nil {
		var p0 *Property = this_.prop0
		_ = p0
		if p0 == nil {
			return false
		}
		var n0 string = p0.GetName()
		_ = n0
		if (n0 == name) || (j2gort.CheckNotNullString(n0) == name) {
			this_.prop0 = this_.prop1
			this_.prop1 = nil
			this_.priv_PropertyMap_incrementCachedKeySetCounter_9413()
			return true
		}
		var p1 *Property = this_.prop1
		_ = p1
		if p1 != nil {
			var n1 string = p1.GetName()
			_ = n1
			if (n1 == name) || (j2gort.CheckNotNullString(n1) == name) {
				this_.prop1 = nil
				this_.priv_PropertyMap_incrementCachedKeySetCounter_9413()
				return true
			}
		}
		return false
	}
	if s.map_ == nil {
		var items []*Property = s.items
		_ = items
		var n int32 = s.size
		_ = n
		{
			var i int32 = int32(0)
			_ = i
			for ; i < n; i = i + 1 {
				var iname string = items[i].GetName()
				_ = iname
				if (iname == name) || (j2gort.CheckNotNullString(iname) == name) {
					var last int32 = (n - int32(1))
					_ = last
					{
						var j int32 = i
						_ = j
						for ; j < last; j = j + 1 {
							{
								a_51 := items
								i_52 := j
								a_51[i_52] = items[(j + int32(1))]
							}
						}
					}
					{
						a_53 := items
						i_54 := last
						a_53[i_54] = nil
					}
					{
						r_55 := s
						r_55.size = last
					}
					if last == int32(2) {
						this_.prop0 = items[int32(0)]
						this_.prop1 = items[int32(1)]
						this_.store = nil
					}
					this_.priv_PropertyMap_incrementCachedKeySetCounter_9413()
					return true
				}
			}
		}
		return false
	}
	if s.map_.Remove(name) != nil {
		this_.priv_PropertyMap_incrementCachedKeySetCounter_9413()
		return true
	}
	return false
}

// Java method: PropertyMap.putProperty(com.google.javascript.rhino.jstype.KnownSymbolType symbol, com.google.javascript.rhino.jstype.Property newProp)
func (this_ *PropertyMap) PutProperty_2_KnownSymbolType_Property(symbol *KnownSymbolType, newProp *Property) {
	_ = symbol
	_ = newProp
	var e *PropertyMap_Extra = this_.priv_PropertyMap_ensureExtra_9380()
	_ = e
	var oldProp *Property
	_ = oldProp
	if e.knownSymbols != nil {
		oldProp = e.knownSymbols.Get(symbol)
	} else {
		oldProp = nil
	}
	if oldProp == nil {
		this_.priv_PropertyMap_incrementCachedKeySetCounter_9413()
	}
	if e.knownSymbols == nil {
		{
			r_56 := e
			r_56.knownSymbols = &j2gort.CustomOrderedMap[*KnownSymbolType, *Property]{}
		}
	}
	e.knownSymbols.Put(symbol, newProp)
}

// Java method: PropertyMap.putProperty(com.google.javascript.rhino.jstype.Property.Key name, com.google.javascript.rhino.jstype.Property newProp)
func (this_ *PropertyMap) PutProperty_2_Property_Key_Property(name Property_Key, newProp *Property) {
	_ = name
	_ = newProp
	switch (name.Kind()).Ordinal() {
	case int32(0) /* KeyKind.STRING */ :
		this_.PutProperty_2_string_Property(name.M_String(), newProp)
	case int32(1) /* KeyKind.SYMBOL */ :
		this_.PutProperty_2_KnownSymbolType_Property(name.Symbol(), newProp)
	}
}

// Java method: PropertyMap.resolveOwnPropertyTypes(com.google.javascript.rhino.ErrorReporter reporter)
func (this_ *PropertyMap) ResolveOwnPropertyTypes(reporter ErrorReporter) {
	_ = reporter
	var s *PropertyMap_PropStore = this_.store
	_ = s
	if s != nil {
		if s.map_ != nil {
			{
				cur_57 := s.map_.IterValues()
				for cur_57.HasNext() {

					prop := cur_57.
						Next()
					_ = prop
					prop.SetType(JSType_SafeResolve(prop.GetType(), reporter))
				}
			}
			return
		}
		var items []*Property = s.items
		_ = items
		{
			var i int32 = int32(0)
			_ = i
			var n int32 = s.size
			_ = n
			for ; i < n; i = i + 1 {
				var prop_2 *Property = items[i]
				_ = prop_2
				prop_2.SetType(JSType_SafeResolve(prop_2.GetType(), reporter))
			}
		}
		return
	}
	var p0 *Property = this_.prop0
	_ = p0
	if p0 != nil {
		p0.SetType(JSType_SafeResolve(p0.GetType(), reporter))
		var p1 *Property = this_.prop1
		_ = p1
		if p1 != nil {
			p1.SetType(JSType_SafeResolve(p1.GetType(), reporter))
		}
	}
}

// Java method: PropertyMap.values()
func (this_ *PropertyMap) Values() j2gort.Ext_java_lang_Iterable {
	var s *PropertyMap_PropStore = this_.store
	_ = s
	if s != nil {
		if s.map_ != nil {
			return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](s.map_.LiveValues())
		}
		var n int32 = s.size
		_ = n
		var items []*Property = s.items
		_ = items
		var builder *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = j2gort.InitImmListBuilderCap(&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{}, n)
		_ = builder
		{
			var i int32 = int32(0)
			_ = i
			for ; i < n; i = i + 1 {
				builder.AddVoid_(j2gort.PtrToAny(items[i]))
			}
		}
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](builder.BuildImmList_())
	}
	var p0 *Property = this_.prop0
	_ = p0
	if p0 == nil {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
	}
	var p1 *Property = this_.prop1
	_ = p1
	if p1 == nil {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.ImmListOfPtr(p0))
	}
	return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_ImmutableList_Of_2(j2gort.PtrToAny(p0), j2gort.PtrToAny(p1)))
}

// Java method: PropertyMap.hashCode()
func (this_ *PropertyMap) HashCode() int32 {
	var h int32 = int32(0)
	_ = h
	var s *PropertyMap_PropStore = this_.store
	_ = s
	if s != nil {
		if s.map_ != nil {
			{
				cur_59 := s.map_.IterKeys()
				for cur_59.HasNext() {
					key := cur_59.
						Next()
					_ = key
					h = (h + javalang.StringHashCode(j2gort.CheckNotNullString(key)))
				}
			}
		} else {
			var items []*Property = s.items
			_ = items
			{
				var i int32 = int32(0)
				_ = i
				var n int32 = s.size
				_ = n
				for ; i < n; i = i + 1 {
					h = (h + javalang.StringHashCode(j2gort.CheckNotNullString(items[i].GetName())))
				}
			}
		}
	} else if this_.prop0 != nil {
		h = javalang.StringHashCode(j2gort.CheckNotNullString(this_.prop0.GetName()))
		if this_.prop1 != nil {
			h = (h + javalang.StringHashCode(j2gort.CheckNotNullString(this_.prop1.GetName())))
		}
	}
	var e *PropertyMap_Extra = this_.extra
	_ = e
	if (e == nil) || (e.knownSymbols == nil) {
		return h
	}
	return ((int32(31) * (int32(31) + h)) + j2gort.AsIface[j2gort.Ext_java_util_Set](e.knownSymbols.LiveKeySet()).HashCode())
}

// Java method: PropertyMap.incrementCachedKeySetCounter()
func (this_ *PropertyMap) priv_PropertyMap_incrementCachedKeySetCounter_9413() {
	var e *PropertyMap_Extra = this_.extra
	_ = e
	if e != nil {
		r_60 := e
		r_60.cachedKeySetCounter = r_60.cachedKeySetCounter + 1
		{
			r_61 := e
			r_61.cachedKeySet = nil
		}
		{
			r_62 := e
			r_62.cachedKnownSymbolsKeySet = nil
		}
		if !(e.cachedKeySetCounter >= int32(0)) {
			j2gort.CheckState(false)
		}
	}
}

// Java class: com.google.javascript.rhino.jstype.PropertyMap.Extra
type PropertyMap_Extra struct {
	knownSymbols             *j2gort.CustomOrderedMap[*KnownSymbolType, *Property]   `gson:"knownSymbols"`
	cachedKeySet             j2gort.Ext_com_google_common_collect_ImmutableSortedSet `gson:"cachedKeySet"`
	cachedKnownSymbolsKeySet j2gort.Ext_com_google_common_collect_ImmutableSet       `gson:"cachedKnownSymbolsKeySet"`
	cachedKeySetCounter      int32                                                   `gson:"cachedKeySetCounter"`
	cachedAncestorHash       int32                                                   `gson:"cachedAncestorHash"`
	idHash_                  atomic.Int32
}

func (this_ *PropertyMap_Extra) BindSelf_(s any) {
	_ = s
}
func (this_ *PropertyMap_Extra) IsJavaNil_() bool { return this_ == nil }
func (*PropertyMap_Extra) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.PropertyMap$Extra", "Extra"
}
func (this_ *PropertyMap_Extra) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PropertyMap_Extra) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.PropertyMap.Extra", "com.google.javascript.rhino.jstype.PropertyMap$Extra", func(v_ any) bool {
		return j2gort.InstanceOf[*PropertyMap_Extra](v_)
	})
}

func (this_ *PropertyMap_Extra) initFields_PropertyMap_Extra() {
	_ = this_

}

// Java constructor: Extra()
func NewPropertyMap_Extra() *PropertyMap_Extra {
	this_ := &PropertyMap_Extra{}
	InitPropertyMap_Extra(this_)
	return this_
}

func InitPropertyMap_Extra(this_ *PropertyMap_Extra) {
	this_.initFields_PropertyMap_Extra()
}

// Java class: com.google.javascript.rhino.jstype.PropertyMap.PropStore
type PropertyMap_PropStore struct {
	items   []*Property                           `gson:"items"`
	map_    *j2gort.OrderedMap[string, *Property] `gson:"map"`
	size    int32                                 `gson:"size"`
	idHash_ atomic.Int32
}

func (this_ *PropertyMap_PropStore) BindSelf_(s any) {
	_ = s
}
func (this_ *PropertyMap_PropStore) IsJavaNil_() bool { return this_ == nil }
func (*PropertyMap_PropStore) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.PropertyMap$PropStore", "PropStore"
}
func (this_ *PropertyMap_PropStore) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PropertyMap_PropStore) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.PropertyMap.PropStore", "com.google.javascript.rhino.jstype.PropertyMap$PropStore", func(v_ any) bool {
		return j2gort.InstanceOf[*PropertyMap_PropStore](v_)
	})
}

func (this_ *PropertyMap_PropStore) initFields_PropertyMap_PropStore() {
	_ = this_

}

// Java constructor: PropStore(com.google.javascript.rhino.jstype.Property p0, com.google.javascript.rhino.jstype.Property p1, com.google.javascript.rhino.jstype.Property p2)
func NewPropertyMap_PropStore(p0 *Property, p1 *Property, p2 *Property) *PropertyMap_PropStore {
	this_ := &PropertyMap_PropStore{}
	InitPropertyMap_PropStore(this_, p0, p1, p2)
	return this_
}

func InitPropertyMap_PropStore(this_ *PropertyMap_PropStore, p0 *Property, p1 *Property, p2 *Property) {
	_ = p0
	_ = p1
	_ = p2
	this_.initFields_PropertyMap_PropStore()
	var arr []*Property = make([]*Property, int32(4))
	_ = arr
	{
		a_63 := arr
		i_64 := int32(0)
		a_63[i_64] = p0
	}
	{
		a_65 := arr
		i_66 := int32(1)
		a_65[i_66] = p1
	}
	{
		a_67 := arr
		i_68 := int32(2)
		a_67[i_68] = p2
	}
	this_.items = arr
	this_.size = int32(3)
}

// Java record: com.google.javascript.rhino.jstype.PropertyMap.AllKeys
type PropertyMap_AllKeys struct {
	stringKeys      j2gort.Ext_com_google_common_collect_ImmutableSortedSet `gson:"stringKeys"`
	knownSymbolKeys j2gort.Ext_java_util_Set                                `gson:"knownSymbolKeys"`
}

func (this_ *PropertyMap_AllKeys) BindSelf_(s any) {
	_ = s
}
func (this_ *PropertyMap_AllKeys) IsJavaNil_() bool { return this_ == nil }
func (*PropertyMap_AllKeys) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.PropertyMap$AllKeys", "AllKeys"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.PropertyMap.AllKeys", "com.google.javascript.rhino.jstype.PropertyMap$AllKeys", func(v_ any) bool {
		return j2gort.InstanceOf[*PropertyMap_AllKeys](v_)
	})
}
func (*PropertyMap_AllKeys) Implements_j2gort_Ext_java_lang_Record() {}

func (this_ *PropertyMap_AllKeys) initFields_PropertyMap_AllKeys() {
	_ = this_
}

// Java constructor: AllKeys(com.google.common.collect.ImmutableSortedSet<java.lang.String> stringKeys, java.util.Set<com.google.javascript.rhino.jstype.KnownSymbolType> knownSymbolKeys)
func NewPropertyMap_AllKeys(stringKeys j2gort.Ext_com_google_common_collect_ImmutableSortedSet, knownSymbolKeys j2gort.Ext_java_util_Set) *PropertyMap_AllKeys {
	this_ := &PropertyMap_AllKeys{}
	InitPropertyMap_AllKeys(this_, stringKeys, knownSymbolKeys)
	return this_
}

func InitPropertyMap_AllKeys(this_ *PropertyMap_AllKeys, stringKeys j2gort.Ext_com_google_common_collect_ImmutableSortedSet, knownSymbolKeys j2gort.Ext_java_util_Set) {
	_ = stringKeys
	_ = knownSymbolKeys
	this_.initFields_PropertyMap_AllKeys()
	this_.stringKeys = stringKeys
	this_.knownSymbolKeys = knownSymbolKeys
}

func (this_ *PropertyMap_AllKeys) StringKeys() j2gort.Ext_com_google_common_collect_ImmutableSortedSet {
	return this_.stringKeys
}
func (this_ *PropertyMap_AllKeys) KnownSymbolKeys() j2gort.Ext_java_util_Set {
	return this_.knownSymbolKeys
}
func (this_ *PropertyMap_AllKeys) Equals(other_ any) bool {
	o_, ok_ := other_.(*PropertyMap_AllKeys)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if !javalang.ObjectsEquals(this_.stringKeys, o_.stringKeys) {
		return false
	}
	if !javalang.ObjectsEquals(this_.knownSymbolKeys, o_.knownSymbolKeys) {
		return false
	}
	return true
}
func (this_ *PropertyMap_AllKeys) HashCode() int32 {
	var h_ int32 = 0
	h_ = 31*h_ + j2gort.ObjectHash(this_.stringKeys)
	h_ = 31*h_ + j2gort.ObjectHash(this_.knownSymbolKeys)
	return h_
}
func (this_ *PropertyMap_AllKeys) String() string {
	_ = fmt.Sprint
	return "AllKeys[" + "stringKeys=" + j2gort.StrVal(this_.stringKeys) + ", " + "knownSymbolKeys=" + j2gort.StrVal(this_.knownSymbolKeys) + "]"
}
