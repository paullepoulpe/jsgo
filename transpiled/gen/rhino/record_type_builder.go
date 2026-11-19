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
// Port of com/google/javascript/rhino/jstype/RecordTypeBuilder.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.rhino.jstype.RecordTypeBuilder
type RecordTypeBuilder struct {
	registry   *JSTypeRegistry                     `gson:"registry"`
	properties *j2gort.Ext_java_util_LinkedHashMap `gson:"properties"`
	idHash_    atomic.Int32
	isEmpty    bool `gson:"isEmpty"`
	isDeclared bool `gson:"isDeclared"`
}

func (this_ *RecordTypeBuilder) BindSelf_(s any) {
	_ = s
}
func (this_ *RecordTypeBuilder) IsJavaNil_() bool { return this_ == nil }
func (*RecordTypeBuilder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.RecordTypeBuilder", "RecordTypeBuilder"
}
func (this_ *RecordTypeBuilder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RecordTypeBuilder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.RecordTypeBuilder", "com.google.javascript.rhino.jstype.RecordTypeBuilder", func(v_ any) bool {
		return j2gort.InstanceOf[*RecordTypeBuilder](v_)
	})
}

func (this_ *RecordTypeBuilder) initFields_RecordTypeBuilder() {
	_ = this_
	this_.isEmpty = true
	this_.isDeclared = true
	this_.properties = j2gort.NewExt_java_util_LinkedHashMap()
}

// Java constructor: RecordTypeBuilder(com.google.javascript.rhino.jstype.JSTypeRegistry registry)
func NewRecordTypeBuilder(registry *JSTypeRegistry) *RecordTypeBuilder {
	this_ := &RecordTypeBuilder{}
	InitRecordTypeBuilder(this_, registry)
	return this_
}

func InitRecordTypeBuilder(this_ *RecordTypeBuilder, registry *JSTypeRegistry) {
	_ = registry
	this_.initFields_RecordTypeBuilder()
	this_.registry = registry
}

// Java method: RecordTypeBuilder.setSynthesized(boolean synthesized)
func (this_ *RecordTypeBuilder) SetSynthesized(synthesized bool) {
	_ = synthesized
	this_.isDeclared = !(synthesized)
}

// Java method: RecordTypeBuilder.addProperty(java.lang.String name, com.google.javascript.rhino.jstype.JSType type, com.google.javascript.rhino.Node propertyNode)
func (this_ *RecordTypeBuilder) AddProperty(name string, type_ JSType, propertyNode Node) *RecordTypeBuilder {
	_ = name
	_ = type_
	_ = propertyNode
	this_.isEmpty = false
	j2gort.MapPutStr(this_.properties, name, NewRecordTypeBuilder_RecordProperty(type_, propertyNode))
	return this_
}

// Java method: RecordTypeBuilder.build()
func (this_ *RecordTypeBuilder) Build() JSType {
	if this_.isEmpty {
		return Upcast_JSType(this_.registry.GetNativeObjectType(StaticGet_JSTypeNative_OBJECT_TYPE()))
	}
	var m *j2gort.Ext_com_google_common_collect_ImmutableSortedMap_Builder = j2gort.Ext_com_google_common_collect_ImmutableSortedMap_NaturalOrder()
	_ = m
	m.PutAll_1_Ext_java_util_Map(j2gort.AsIface[j2gort.Ext_java_util_Map](this_.properties))
	return NewRecordType(this_.registry, j2gort.AsIface[j2gort.Ext_java_util_Map](j2gort.CastPtr[j2gort.Ext_com_google_common_collect_ImmutableSortedMap](m.BuildOrThrow())), this_.isDeclared)
}

// Java class: com.google.javascript.rhino.jstype.RecordTypeBuilder.RecordProperty
type RecordTypeBuilder_RecordProperty struct {
	type_        JSType `gson:"type"`
	propertyNode Node   `gson:"propertyNode"`
	idHash_      atomic.Int32
}

func (this_ *RecordTypeBuilder_RecordProperty) BindSelf_(s any) {
	_ = s
}
func (this_ *RecordTypeBuilder_RecordProperty) IsJavaNil_() bool { return this_ == nil }
func (*RecordTypeBuilder_RecordProperty) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.RecordTypeBuilder$RecordProperty", "RecordProperty"
}
func (this_ *RecordTypeBuilder_RecordProperty) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RecordTypeBuilder_RecordProperty) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.RecordTypeBuilder.RecordProperty", "com.google.javascript.rhino.jstype.RecordTypeBuilder$RecordProperty", func(v_ any) bool {
		return j2gort.InstanceOf[*RecordTypeBuilder_RecordProperty](v_)
	})
}

func (this_ *RecordTypeBuilder_RecordProperty) initFields_RecordTypeBuilder_RecordProperty() {
	_ = this_
}

// Java constructor: RecordProperty(com.google.javascript.rhino.jstype.JSType type, com.google.javascript.rhino.Node propertyNode)
func NewRecordTypeBuilder_RecordProperty(type_ JSType, propertyNode Node) *RecordTypeBuilder_RecordProperty {
	this_ := &RecordTypeBuilder_RecordProperty{}
	InitRecordTypeBuilder_RecordProperty(this_, type_, propertyNode)
	return this_
}

func InitRecordTypeBuilder_RecordProperty(this_ *RecordTypeBuilder_RecordProperty, type_ JSType, propertyNode Node) {
	_ = type_
	_ = propertyNode
	this_.initFields_RecordTypeBuilder_RecordProperty()
	this_.type_ = type_
	this_.propertyNode = propertyNode
}

// Java method: RecordProperty.getType()
func (this_ *RecordTypeBuilder_RecordProperty) GetType() JSType {
	return this_.type_
}

// Java method: RecordProperty.getPropertyNode()
func (this_ *RecordTypeBuilder_RecordProperty) GetPropertyNode() Node {
	return this_.propertyNode
}

// Java method: RecordProperty.toString()
func (this_ *RecordTypeBuilder_RecordProperty) String() string {
	return ("RecordProperty{type: " + j2gort.StrVal(this_.type_) + ", node: " + j2gort.StrVal(this_.propertyNode) + "}")
}
