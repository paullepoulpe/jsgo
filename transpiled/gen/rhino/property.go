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
// Port of com/google/javascript/rhino/jstype/Property.java (generated mechanically by j2go).

package rhino

import (
	"fmt"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.rhino.jstype.Property
type Property struct {
	name         string     `gson:"name"`
	type_        JSType     `gson:"type"`
	docInfo      *JSDocInfo `gson:"docInfo"`
	propertyNode Node       `gson:"propertyNode"`
	inferred     bool       `gson:"inferred"`
}

func (this_ *Property) BindSelf_(s any) {
	_ = s
}
func (this_ *Property) IsJavaNil_() bool { return this_ == nil }
func (*Property) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.Property", "Property"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.Property", "com.google.javascript.rhino.jstype.Property", func(v_ any) bool {
		return j2gort.InstanceOf[*Property](v_)
	})
}
func (*Property) Implements_rhino_StaticTypedSlot() {}
func (*Property) Implements_rhino_StaticTypedRef()  {}
func (*Property) Implements_rhino_StaticSlot()      {}
func (*Property) Implements_rhino_StaticRef()       {}

func StaticGet_Property_SerialVersionUID() int64 { return int64(1) }
func (this_ *Property) initFields_Property() {
	_ = this_
	this_.name = j2gort.NullString

}

// Java constructor: Property(java.lang.String name, com.google.javascript.rhino.jstype.JSType type, boolean inferred, com.google.javascript.rhino.Node propertyNode)
func NewProperty_4_string_JSType_bool_Node(name string, type_ JSType, inferred bool, propertyNode Node) *Property {
	this_ := &Property{}
	InitProperty_4_string_JSType_bool_Node(this_, name, type_, inferred, propertyNode)
	return this_
}

func InitProperty_4_string_JSType_bool_Node(this_ *Property, name string, type_ JSType, inferred bool, propertyNode Node) {
	_ = name
	_ = type_
	_ = inferred
	_ = propertyNode
	this_.initFields_Property()
	this_.name = j2gort.CheckNotNullString(name)
	this_.type_ = j2gort.CheckNotNullIface(type_)
	this_.inferred = inferred
	this_.propertyNode = propertyNode
}

// Java constructor: Property(com.google.javascript.rhino.jstype.Property.Key name, com.google.javascript.rhino.jstype.JSType type, boolean inferred, com.google.javascript.rhino.Node propertyNode)
func NewProperty_4_Property_Key_JSType_bool_Node(name Property_Key, type_ JSType, inferred bool, propertyNode Node) *Property {
	this_ := &Property{}
	InitProperty_4_Property_Key_JSType_bool_Node(this_, name, type_, inferred, propertyNode)
	return this_
}

func InitProperty_4_Property_Key_JSType_bool_Node(this_ *Property, name Property_Key, type_ JSType, inferred bool, propertyNode Node) {
	_ = name
	_ = type_
	_ = inferred
	_ = propertyNode
	this_.initFields_Property()
	j2gort.CheckNotNullIface(name)
	if name.Kind() == StaticGet_Property_KeyKind_STRING() {
		this_.name = j2gort.CheckNotNullString(name.M_String())
	} else {
		this_.name = j2gort.CheckNotNullPtr(name.Symbol()).GetDisplayName()
	}
	this_.type_ = j2gort.CheckNotNullIface(type_)
	this_.inferred = inferred
	this_.propertyNode = propertyNode
}

// Java method: Property.getName()
func (this_ *Property) GetName() string {
	return this_.name
}

// Java method: Property.getNode()
func (this_ *Property) GetNode() Node {
	return this_.propertyNode
}

// Java method: Property.getSourceFile()
func (this_ *Property) GetSourceFile() StaticSourceFile {
	if this_.propertyNode == 0 {
		return nil
	} else {
		return this_.propertyNode.GetStaticSourceFile()
	}
}

// Java method: Property.getSymbol()
func (this_ *Property) GetSymbol() any {
	return this_
}

// Java method: Property.getDeclaration()
func (this_ *Property) GetDeclaration() any {
	if this_.propertyNode == 0 {
		return nil
	} else {
		return this_
	}
}

// Java method: Property.getType()
func (this_ *Property) GetType() JSType {
	return this_.type_
}

// Java method: Property.isTypeInferred()
func (this_ *Property) IsTypeInferred() bool {
	return this_.inferred
}

// Java method: Property.isFromExterns()
func (this_ *Property) IsFromExterns() bool {
	if this_.propertyNode == 0 {
		return false
	} else {
		return this_.propertyNode.IsFromExterns()
	}
}

// Java method: Property.setType(com.google.javascript.rhino.jstype.JSType type)
func (this_ *Property) SetType(type_ JSType) {
	_ = type_
	this_.type_ = j2gort.CheckNotNullIface(type_)
}

// Java method: Property.getJSDocInfo()
func (this_ *Property) GetJSDocInfo() *JSDocInfo {
	return this_.docInfo
}

// Java method: Property.setJSDocInfo(com.google.javascript.rhino.JSDocInfo info)
func (this_ *Property) SetJSDocInfo(info *JSDocInfo) {
	_ = info
	this_.docInfo = info
}

// Java method: Property.setNode(com.google.javascript.rhino.Node n)
func (this_ *Property) SetNode(n Node) {
	_ = n
	this_.propertyNode = n
}

// Java method: Property.toString()
func (this_ *Property) String() string {
	return ("Property {  name: " + j2gort.StrValString(this_.name) + ", type:" + j2gort.StrVal(this_.type_) + ", inferred: " + j2gort.StrValBool(this_.inferred) + "}")
}

// Java method: Property.getScope()
func (this_ *Property) GetScope() any {
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException")))
}

// Java method: Property.hashCode()
func (this_ *Property) HashCode() int32 {
	return ((int32(31) * (int32(31) + javalang.StringHashCode(j2gort.CheckNotNullString(this_.name)))) + func() int32 {
		hi_1 := this_.type_
		if hi_1 == nil || hi_1.IsJavaNil_() {
			return 0
		}
		return hi_1.GetJSTypeBase().HashCode()
	}())
}

// Java enum: com.google.javascript.rhino.jstype.Property.KeyKind
type Property_KeyKind uint8

const (
	enum_Property_KeyKind_STRING Property_KeyKind = 1
	enum_Property_KeyKind_SYMBOL Property_KeyKind = 2
)

var enumBaseTable_Property_KeyKind = [...]j2gort.EnumBase{
	{Name_: "STRING", Ordinal_: 0},
	{Name_: "SYMBOL", Ordinal_: 1},
}
var enumValues_Property_KeyKind = []Property_KeyKind{1, 2}

func StaticGet_Property_KeyKind_STRING() Property_KeyKind {
	return enum_Property_KeyKind_STRING
}
func StaticGet_Property_KeyKind_SYMBOL() Property_KeyKind {
	return enum_Property_KeyKind_SYMBOL
}

func Property_KeyKind_Values() []Property_KeyKind {
	return j2gort.CloneSlice(enumValues_Property_KeyKind)
}

func Property_KeyKind_ValueOf(name string) Property_KeyKind {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_Property_KeyKind {
		if enumBaseTable_Property_KeyKind[i_].Name_ == name {
			return Property_KeyKind(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.rhino.jstype.Property.KeyKind." + name))
}

func (this_ Property_KeyKind) IsJavaNil_() bool                              { return this_ == 0 }
func (Property_KeyKind) HasIdentityEquals_()                                 {}
func (Property_KeyKind) Implements_j2gort_Ext_java_lang_Enum()               {}
func (Property_KeyKind) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ Property_KeyKind) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (Property_KeyKind) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.Property$KeyKind", "KeyKind"
}
func (this_ Property_KeyKind) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_Property_KeyKind[this_-1]
}
func (this_ Property_KeyKind) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ Property_KeyKind) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ Property_KeyKind) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Property_KeyKind[this_-1].Name_
}
func (this_ Property_KeyKind) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_Property_KeyKind[this_-1].Name_
}
func (this_ Property_KeyKind) Equals(other_ any) bool {
	o_, ok_ := other_.(Property_KeyKind)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ Property_KeyKind) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_Property_KeyKind[this_-1].HashCode()
}
func (this_ Property_KeyKind) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_Property_KeyKind[this_-1].CompareTo(other_)
}
func (this_ Property_KeyKind) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_Property_KeyKind[this_-1].Name_)
}
func init() {
	vals_ := []any{Property_KeyKind(1), Property_KeyKind(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.jstype.Property.KeyKind", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.rhino.jstype.Property$KeyKind", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.Property.KeyKind", "com.google.javascript.rhino.jstype.Property$KeyKind", func(v_ any) bool {
		return j2gort.InstanceOf[Property_KeyKind](v_)
	})
}

// Java interface: com.google.javascript.rhino.jstype.Property.Key
type Property_Key interface {
	Implements_rhino_Property_Key()
	IsJavaNil_() bool
	Kind() Property_KeyKind
	M_String() string
	Symbol() *KnownSymbolType
	HumanReadableName() string
	Matches(string) bool
}

func ToIface_Property_Key[T interface {
	Property_Key
	comparable
}](v T) Property_Key {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_Property_Key[T Property_Key](v T) Property_Key {
	var vi Property_Key = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.Property.Key", "com.google.javascript.rhino.jstype.Property$Key", func(v_ any) bool {
		return j2gort.InstanceOf[Property_Key](v_)
	})
}

// Java default interface method: Key.matches(java.lang.String stringKey)
func IfaceDefault_Property_Key_Matches(this_ Property_Key, stringKey string) bool {
	_ = this_
	_ = stringKey
	return (func() bool {
		re_2 := this_.Kind()
		ae_3 := StaticGet_Property_KeyKind_STRING()
		if re_2 == 0 {
			j2gort.ThrowNPE()
		}
		return re_2 == ae_3
	}() && (j2gort.CheckNotNullString(this_.M_String()) == stringKey))
}

// Java record: com.google.javascript.rhino.jstype.Property.StringKey
type Property_StringKey struct {
	string_ string `gson:"string"`
}

func (this_ *Property_StringKey) BindSelf_(s any) {
	_ = s
}
func (this_ *Property_StringKey) IsJavaNil_() bool { return this_ == nil }
func (*Property_StringKey) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.Property$StringKey", "StringKey"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.Property.StringKey", "com.google.javascript.rhino.jstype.Property$StringKey", func(v_ any) bool {
		return j2gort.InstanceOf[*Property_StringKey](v_)
	})
}
func (*Property_StringKey) Implements_j2gort_Ext_java_lang_Record() {}
func (*Property_StringKey) Implements_rhino_Property_Key()          {}

func (this_ *Property_StringKey) initFields_Property_StringKey() {
	_ = this_
}

// Java constructor: StringKey(java.lang.String string)
func NewProperty_StringKey(string_ string) *Property_StringKey {
	this_ := &Property_StringKey{}
	InitProperty_StringKey(this_, string_)
	return this_
}

func InitProperty_StringKey(this_ *Property_StringKey, string_ string) {
	_ = string_
	this_.initFields_Property_StringKey()
	j2gort.CheckNotNullString(string_)
	this_.string_ = string_
}

// Java method: StringKey.symbol()
func (this_ *Property_StringKey) Symbol() *KnownSymbolType {
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException")))
}

// Java method: StringKey.kind()
func (this_ *Property_StringKey) Kind() Property_KeyKind {
	return StaticGet_Property_KeyKind_STRING()
}

// Java method: StringKey.humanReadableName()
func (this_ *Property_StringKey) HumanReadableName() string {
	return this_.M_String()
}

func (this_ *Property_StringKey) M_String() string { return this_.string_ }
func (this_ *Property_StringKey) Equals(other_ any) bool {
	o_, ok_ := other_.(*Property_StringKey)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if this_.string_ != o_.string_ {
		return false
	}
	return true
}
func (this_ *Property_StringKey) HashCode() int32 {
	var h_ int32 = 0
	{
		hs_4 := this_.string_
		h_ = 31 * h_
		if hs_4 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_4)
		}
	}
	return h_
}
func (this_ *Property_StringKey) String() string {
	_ = fmt.Sprint
	return "StringKey[" + "string=" + j2gort.StrVal(j2gort.BoxString(this_.string_)) + "]"
}

func (this_ *Property_StringKey) Matches(p0_ string) bool {
	return IfaceDefault_Property_Key_Matches(this_, p0_)
}

// Java record: com.google.javascript.rhino.jstype.Property.SymbolKey
type Property_SymbolKey struct {
	symbol *KnownSymbolType `gson:"symbol"`
}

func (this_ *Property_SymbolKey) BindSelf_(s any) {
	_ = s
}
func (this_ *Property_SymbolKey) IsJavaNil_() bool { return this_ == nil }
func (*Property_SymbolKey) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.Property$SymbolKey", "SymbolKey"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.Property.SymbolKey", "com.google.javascript.rhino.jstype.Property$SymbolKey", func(v_ any) bool {
		return j2gort.InstanceOf[*Property_SymbolKey](v_)
	})
}
func (*Property_SymbolKey) Implements_j2gort_Ext_java_lang_Record() {}
func (*Property_SymbolKey) Implements_rhino_Property_Key()          {}

func (this_ *Property_SymbolKey) initFields_Property_SymbolKey() {
	_ = this_
}

// Java constructor: SymbolKey(com.google.javascript.rhino.jstype.KnownSymbolType symbol)
func NewProperty_SymbolKey(symbol *KnownSymbolType) *Property_SymbolKey {
	this_ := &Property_SymbolKey{}
	InitProperty_SymbolKey(this_, symbol)
	return this_
}

func InitProperty_SymbolKey(this_ *Property_SymbolKey, symbol *KnownSymbolType) {
	_ = symbol
	this_.initFields_Property_SymbolKey()
	j2gort.CheckNotNullPtr(symbol)
	this_.symbol = symbol
}

// Java method: SymbolKey.string()
func (this_ *Property_SymbolKey) M_String() string {
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException")))
}

// Java method: SymbolKey.kind()
func (this_ *Property_SymbolKey) Kind() Property_KeyKind {
	return StaticGet_Property_KeyKind_SYMBOL()
}

// Java method: SymbolKey.humanReadableName()
func (this_ *Property_SymbolKey) HumanReadableName() string {
	return this_.Symbol().GetDisplayName()
}

func (this_ *Property_SymbolKey) Symbol() *KnownSymbolType { return this_.symbol }
func (this_ *Property_SymbolKey) Equals(other_ any) bool {
	o_, ok_ := other_.(*Property_SymbolKey)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if !javalang.ObjectsEquals(j2gort.PtrToAny(this_.symbol), j2gort.PtrToAny(o_.symbol)) {
		return false
	}
	return true
}
func (this_ *Property_SymbolKey) HashCode() int32 {
	var h_ int32 = 0
	{
		hp_4 := this_.symbol
		h_ = 31 * h_
		if hp_4 != nil {
			h_ += hp_4.HashCode()
		}
	}
	return h_
}
func (this_ *Property_SymbolKey) String() string {
	_ = fmt.Sprint
	return "SymbolKey[" + "symbol=" + j2gort.StrVal(j2gort.PtrToAny(this_.symbol)) + "]"
}

func (this_ *Property_SymbolKey) Matches(p0_ string) bool {
	return IfaceDefault_Property_Key_Matches(this_, p0_)
}

// Java class: com.google.javascript.rhino.jstype.Property.OwnedProperty
type Property_OwnedProperty struct {
	owner   ObjectType `gson:"owner"`
	value   *Property  `gson:"value"`
	idHash_ atomic.Int32
}

func (this_ *Property_OwnedProperty) BindSelf_(s any) {
	_ = s
}
func (this_ *Property_OwnedProperty) IsJavaNil_() bool { return this_ == nil }
func (*Property_OwnedProperty) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.Property$OwnedProperty", "OwnedProperty"
}
func (this_ *Property_OwnedProperty) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Property_OwnedProperty) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.Property.OwnedProperty", "com.google.javascript.rhino.jstype.Property$OwnedProperty", func(v_ any) bool {
		return j2gort.InstanceOf[*Property_OwnedProperty](v_)
	})
}

func (this_ *Property_OwnedProperty) initFields_Property_OwnedProperty() {
	_ = this_
}

// Java constructor: OwnedProperty(com.google.javascript.rhino.jstype.ObjectType owner, com.google.javascript.rhino.jstype.Property value)
func NewProperty_OwnedProperty(owner ObjectType, value *Property) *Property_OwnedProperty {
	this_ := &Property_OwnedProperty{}
	InitProperty_OwnedProperty(this_, owner, value)
	return this_
}

func InitProperty_OwnedProperty(this_ *Property_OwnedProperty, owner ObjectType, value *Property) {
	_ = owner
	_ = value
	this_.initFields_Property_OwnedProperty()
	this_.owner = owner
	this_.value = value
}

// Java method: OwnedProperty.getOwner()
func (this_ *Property_OwnedProperty) GetOwner() ObjectType {
	return this_.owner
}

// Java method: OwnedProperty.getValue()
func (this_ *Property_OwnedProperty) GetValue() *Property {
	return this_.value
}

// Java method: OwnedProperty.getOwnerInstanceType()
func (this_ *Property_OwnedProperty) GetOwnerInstanceType() ObjectType {
	if this_.owner.GetObjectTypeBase().IsFunctionPrototypeType() {
		return this_.owner.GetOwnerFunction().GetInstanceType()
	} else {
		return this_.owner
	}
}

// Java method: OwnedProperty.isOwnedByInterface()
func (this_ *Property_OwnedProperty) IsOwnedByInterface() bool {
	if this_.owner.GetObjectTypeBase().IsFunctionPrototypeType() {
		return this_.owner.GetOwnerFunction().IsInterface()
	} else {
		return this_.owner.GetObjectTypeBase().IsInterface()
	}
}
