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
// Port of com/google/javascript/rhino/jstype/ValueType.java (generated mechanically by j2go).

package rhino

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java polymorphic class interface: com.google.javascript.rhino.jstype.ValueType
type ValueType interface {
	Implements_rhino_ValueType()
	IsJavaNil_() bool
	GetValueTypeBase() *ValueTypeBase
	JSType
}
func ToIface_ValueType[T interface {
	ValueType
	comparable
}](v T) ValueType {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_ValueType[T ValueType](v T) ValueType {
	var vi ValueType = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.rhino.jstype.ValueType
type ValueTypeBase struct {
	JSTypeBase
}

func (this_ *ValueTypeBase) BindSelf_(s any) {
	_ = s
	this_.JSTypeBase.BindSelf_(s)
}
func (this_ *ValueTypeBase) IsJavaNil_() bool { return this_ == nil }
func (*ValueTypeBase) JavaClassInfo_() (string, string) { return "com.google.javascript.rhino.jstype.ValueType", "ValueType" }
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.ValueType", "com.google.javascript.rhino.jstype.ValueType", func(v_ any) bool {
		return j2gort.InstanceOf[ValueType](v_)
	})
}
func (b *ValueTypeBase) GetValueTypeBase() *ValueTypeBase { return b }
func (this_ *ValueTypeBase) self_ValueType_() ValueType {
	if s := this_.JSTypeBase.self_; s != nil {
		return s.(ValueType)
	}
	return nil
}
func (*ValueTypeBase) Implements_rhino_ValueType() {}
func (*ValueTypeBase) Implements_rhino_JSType() {}

func (this_ *ValueTypeBase) initFields_ValueType() {
	_ = this_
}

func InitValueType(this_ *ValueTypeBase, registry *JSTypeRegistry, typeKindMask int32) {
	_ = registry
	_ = typeKindMask
	InitJSType_3(&this_.JSTypeBase, registry, nil, typeKindMask)
	this_.initFields_ValueType()
	this_.GetJSTypeBase().jstypeFlags = int8((int32(this_.GetJSTypeBase().jstypeFlags) | int32(int8(32))))
	this_.EagerlyResolveToSelf()
}

// Java method: ValueType.resolveInternal(com.google.javascript.rhino.ErrorReporter reporter)
func (this_ *ValueTypeBase) ResolveInternal(reporter ErrorReporter) JSType {
	_ = reporter
	panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError")))
}

// Java method: ValueType.appendTo(com.google.javascript.rhino.jstype.TypeStringBuilder sb)
func (this_ *ValueTypeBase) AppendTo(sb *TypeStringBuilder) {
	_ = sb
	sb.Append_1_string(this_.self_ValueType_().GetDisplayName())
}

// Java method: ValueType.getDisplayName()
func (this_ *ValueTypeBase) GetDisplayName() string {
	panic("abstract or native method invoked")
}

// Java method: ValueType.hasDisplayName()
func (this_ *ValueTypeBase) HasDisplayName() bool {
	return true
}

// Java method: ValueType.visit(com.google.javascript.rhino.jstype.RelationshipVisitor<T> visitor, com.google.javascript.rhino.jstype.JSType that)
func (this_ *ValueTypeBase) Visit_2(visitor RelationshipVisitor, that JSType) any {
	_ = visitor
	_ = that
	return visitor.CaseValueType(this_.self_ValueType_(), that)
}

// Java method: ValueType.getPropertyKind(com.google.javascript.rhino.jstype.Property.Key propertyName, boolean autobox)
func (this_ *ValueTypeBase) GetPropertyKind_2_Property_Key_bool(propertyName Property_Key, autobox bool) JSType_HasPropertyKind {
	_ = propertyName
	_ = autobox
	if ((autobox && this_.IsBoxableScalar())) {
		return this_.JSTypeBase.self_.AutoboxesTo().GetPropertyKind_2_Property_Key_bool(propertyName, autobox)
	} else {
		return this_.JSTypeBase.GetPropertyKind_2_Property_Key_bool(propertyName, autobox)
	}
}

// Java method: ValueType.getPropertyKind(java.lang.String propertyName, boolean autobox)
func (this_ *ValueTypeBase) GetPropertyKind_2_string_bool(propertyName string, autobox bool) JSType_HasPropertyKind {
	_ = propertyName
	_ = autobox
	if ((autobox && this_.IsBoxableScalar())) {
		return this_.JSTypeBase.self_.AutoboxesTo().GetPropertyKind_2_string_bool(propertyName, autobox)
	} else {
		return this_.JSTypeBase.GetPropertyKind_2_string_bool(propertyName, autobox)
	}
}

// Java method: ValueType.recursionUnsafeHashCode()
func (this_ *ValueTypeBase) RecursionUnsafeHashCode() int32 {
	return j2gort.IdentityHash(this_.JSTypeBase.self_)
}

