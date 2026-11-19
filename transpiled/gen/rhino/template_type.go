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

/**
 * For functions with function(this: T, ...) and T as arguments, type inference
 * will set the type of this on a function literal argument to the actual type
 * of T.
 *
 */
// Port of com/google/javascript/rhino/jstype/TemplateType.java (generated mechanically by j2go).

package rhino

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.rhino.jstype.TemplateType
type TemplateType struct {
	ProxyObjectTypeBase
	name               string `gson:"name"`
	bound              JSType `gson:"bound"`
	typeTransformation Node   `gson:"typeTransformation"`
}

func (this_ *TemplateType) BindSelf_(s any) {
	_ = s
	this_.ProxyObjectTypeBase.BindSelf_(s)
}
func (this_ *TemplateType) IsJavaNil_() bool { return this_ == nil }
func (*TemplateType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.TemplateType", "TemplateType"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.TemplateType", "com.google.javascript.rhino.jstype.TemplateType", func(v_ any) bool {
		return j2gort.InstanceOf[*TemplateType](v_)
	})
}
func (*TemplateType) Implements_rhino_ProxyObjectType() {}
func (*TemplateType) Implements_rhino_ObjectType()      {}
func (*TemplateType) Implements_rhino_JSType()          {}

var static_TemplateType_TYPE_CLASS JSTypeClass

func StaticGet_TemplateType_TYPE_CLASS() JSTypeClass {
	clinit_TemplateType()
	return static_TemplateType_TYPE_CLASS
}
func StaticSet_TemplateType_TYPE_CLASS(v JSTypeClass) {
	clinit_TemplateType()
	static_TemplateType_TYPE_CLASS = v
}

var clinitDone_TemplateType uint32
var clinitState_TemplateType j2gort.ClinitState

func clinit_TemplateType() {
	if atomic.LoadUint32(&clinitDone_TemplateType) != 0 {
		return
	}
	clinitSlow_TemplateType()
}

func clinitSlow_TemplateType() {
	if !j2gort.EnterClinit(&clinitState_TemplateType, &clinitDone_TemplateType) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_TemplateType, &clinitDone_TemplateType)
	static_TemplateType_TYPE_CLASS = StaticGet_JSTypeClass_TEMPLATE()
}

func (this_ *TemplateType) initFields_TemplateType() {
	_ = this_
	this_.name = j2gort.NullString
}

// Java constructor: TemplateType(com.google.javascript.rhino.jstype.JSTypeRegistry registry, java.lang.String name)
func NewTemplateType_2(registry *JSTypeRegistry, name string) *TemplateType {
	this_ := &TemplateType{}
	this_.BindSelf_(this_)
	InitTemplateType_2(this_, registry, name)
	return this_
}

func InitTemplateType_2(this_ *TemplateType, registry *JSTypeRegistry, name string) {
	_ = registry
	_ = name
	InitTemplateType_4(this_, registry, name, nil, Node(0))
}

// Java constructor: TemplateType(com.google.javascript.rhino.jstype.JSTypeRegistry registry, java.lang.String name, com.google.javascript.rhino.jstype.JSType bound)
func NewTemplateType_3_JSTypeRegistry_string_JSType(registry *JSTypeRegistry, name string, bound JSType) *TemplateType {
	this_ := &TemplateType{}
	this_.BindSelf_(this_)
	InitTemplateType_3_JSTypeRegistry_string_JSType(this_, registry, name, bound)
	return this_
}

func InitTemplateType_3_JSTypeRegistry_string_JSType(this_ *TemplateType, registry *JSTypeRegistry, name string, bound JSType) {
	_ = registry
	_ = name
	_ = bound
	InitTemplateType_4(this_, registry, name, bound, Node(0))
}

// Java constructor: TemplateType(com.google.javascript.rhino.jstype.JSTypeRegistry registry, java.lang.String name, com.google.javascript.rhino.Node typeTransformation)
func NewTemplateType_3_JSTypeRegistry_string_Node(registry *JSTypeRegistry, name string, typeTransformation Node) *TemplateType {
	this_ := &TemplateType{}
	this_.BindSelf_(this_)
	InitTemplateType_3_JSTypeRegistry_string_Node(this_, registry, name, typeTransformation)
	return this_
}

func InitTemplateType_3_JSTypeRegistry_string_Node(this_ *TemplateType, registry *JSTypeRegistry, name string, typeTransformation Node) {
	_ = registry
	_ = name
	_ = typeTransformation
	InitTemplateType_4(this_, registry, name, nil, typeTransformation)
}

// Java constructor: TemplateType(com.google.javascript.rhino.jstype.JSTypeRegistry registry, java.lang.String name, com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType bound, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node typeTransformation)
func NewTemplateType_4(registry *JSTypeRegistry, name string, bound JSType, typeTransformation Node) *TemplateType {
	this_ := &TemplateType{}
	this_.BindSelf_(this_)
	InitTemplateType_4(this_, registry, name, bound, typeTransformation)
	return this_
}

func InitTemplateType_4(this_ *TemplateType, registry *JSTypeRegistry, name string, bound JSType, typeTransformation Node) {
	_ = registry
	_ = name
	_ = bound
	_ = typeTransformation
	InitProxyObjectType_4(&this_.ProxyObjectTypeBase, registry, func() JSType {
		if bound == nil || bound.IsJavaNil_() {
			return j2gort.AsIface[JSType](registry.GetNativeObjectType(StaticGet_JSTypeNative_UNKNOWN_TYPE()))
		}
		return bound
	}(), nil, int32(8192))
	this_.initFields_TemplateType()
	this_.name = name
	if bound == nil || bound.IsJavaNil_() {
		this_.bound = j2gort.AsIface[JSType](registry.GetNativeObjectType(StaticGet_JSTypeNative_UNKNOWN_TYPE()))
	} else {
		this_.bound = bound
	}
	this_.typeTransformation = typeTransformation
	registry.GetResolver().ResolveIfClosed(ToIface_JSType(this_), StaticGet_TemplateType_TYPE_CLASS())
}

// Java method: TemplateType.getTypeClass()
func (this_ *TemplateType) GetTypeClass() JSTypeClass {
	return StaticGet_TemplateType_TYPE_CLASS()
}

// Java method: TemplateType.getReferenceName()
func (this_ *TemplateType) GetReferenceName() string {
	return this_.name
}

// Java method: TemplateType.appendTo(com.google.javascript.rhino.jstype.TypeStringBuilder sb)
func (this_ *TemplateType) AppendTo(sb *TypeStringBuilder) {
	_ = sb
	if j2gort.IfaceEquals[JSType](this_.bound, this_.GetJSTypeBase().registry.GetNativeObjectType(StaticGet_JSTypeNative_UNKNOWN_TYPE())) {
		sb.Append_1_string(this_.name)
	} else {
		sb.Append_1_string(this_.name).Append_1_string(" extends ").Append_1_JSType(this_.bound)
	}
}

// Java method: TemplateType.toMaybeTemplateType()
func (this_ *TemplateType) ToMaybeTemplateType() *TemplateType {
	return this_
}

// Java method: TemplateType.hasAnyTemplateTypesInternal()
func (this_ *TemplateType) HasAnyTemplateTypesInternal() bool {
	return true
}

// Java method: TemplateType.visit(com.google.javascript.rhino.jstype.Visitor<T> visitor)
func (this_ *TemplateType) Visit_1(visitor Visitor) any {
	_ = visitor
	return visitor.CaseTemplateType(this_)
}

// Java method: TemplateType.visit(com.google.javascript.rhino.jstype.RelationshipVisitor<T> visitor, com.google.javascript.rhino.jstype.JSType that)
func (this_ *TemplateType) Visit_2(visitor RelationshipVisitor, that JSType) any {
	_ = visitor
	_ = that
	return visitor.CaseTemplateType(this_, that)
}

// Java method: TemplateType.isTypeTransformation()
func (this_ *TemplateType) IsTypeTransformation() bool {
	return (this_.typeTransformation != 0)
}

// Java method: TemplateType.getTypeTransformation()
func (this_ *TemplateType) GetTypeTransformation() Node {
	return this_.typeTransformation
}

// Java method: TemplateType.getBound()
func (this_ *TemplateType) GetBound() JSType {
	return this_.bound
}

// Java method: TemplateType.setBound(com.google.javascript.rhino.jstype.JSType bound)
func (this_ *TemplateType) SetBound(bound JSType) {
	_ = bound
	this_.bound = bound
	this_.SetReferencedType(bound)
}

// Java method: TemplateType.setValidator(com.google.common.base.Predicate<com.google.javascript.rhino.jstype.JSType> validator)
func (this_ *TemplateType) SetValidator(validator j2gort.Ext_com_google_common_base_Predicate) bool {
	_ = validator
	return validator.Apply(this_)
}

// Java method: TemplateType.containsCycle()
func (this_ *TemplateType) ContainsCycle() bool {
	var typeVisitor *ContainsUpperBoundSuperTypeVisitor = NewContainsUpperBoundSuperTypeVisitor(nil)
	_ = typeVisitor
	return (j2gort.Cast[ContainsUpperBoundSuperTypeVisitor_Result](this_.Visit_1(typeVisitor)) == StaticGet_ContainsUpperBoundSuperTypeVisitor_Result_CYCLE())
}

// Java method: TemplateType.recursionUnsafeHashCode()
func (this_ *TemplateType) RecursionUnsafeHashCode() int32 {
	return j2gort.IdentityHash(this_)
}
