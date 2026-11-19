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
// Port of com/google/javascript/rhino/jstype/InstanceObjectType.java (generated mechanically by j2go).

package rhino

import (
	"strconv"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.rhino.jstype.InstanceObjectType
type InstanceObjectType struct {
	PrototypeObjectTypeBase
	constructor *FunctionType `gson:"constructor"`
}

func (this_ *InstanceObjectType) BindSelf_(s any) {
	_ = s
	this_.PrototypeObjectTypeBase.BindSelf_(s)
}
func (this_ *InstanceObjectType) IsJavaNil_() bool { return this_ == nil }
func (*InstanceObjectType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.InstanceObjectType", "InstanceObjectType"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.InstanceObjectType", "com.google.javascript.rhino.jstype.InstanceObjectType", func(v_ any) bool {
		return j2gort.InstanceOf[*InstanceObjectType](v_)
	})
}
func (*InstanceObjectType) Implements_rhino_PrototypeObjectType() {}
func (*InstanceObjectType) Implements_rhino_ObjectType()          {}
func (*InstanceObjectType) Implements_rhino_JSType()              {}

var static_InstanceObjectType_TYPE_CLASS JSTypeClass

func StaticGet_InstanceObjectType_TYPE_CLASS() JSTypeClass {
	clinit_InstanceObjectType()
	return static_InstanceObjectType_TYPE_CLASS
}
func StaticSet_InstanceObjectType_TYPE_CLASS(v JSTypeClass) {
	clinit_InstanceObjectType()
	static_InstanceObjectType_TYPE_CLASS = v
}

var clinitDone_InstanceObjectType uint32
var clinitState_InstanceObjectType j2gort.ClinitState

func clinit_InstanceObjectType() {
	if atomic.LoadUint32(&clinitDone_InstanceObjectType) != 0 {
		return
	}
	clinitSlow_InstanceObjectType()
}

func clinitSlow_InstanceObjectType() {
	if !j2gort.EnterClinit(&clinitState_InstanceObjectType, &clinitDone_InstanceObjectType) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_InstanceObjectType, &clinitDone_InstanceObjectType)
	static_InstanceObjectType_TYPE_CLASS = StaticGet_JSTypeClass_INSTANCE_OBJECT()
}

func (this_ *InstanceObjectType) initFields_InstanceObjectType() {
	_ = this_
}

// Java constructor: InstanceObjectType(com.google.javascript.rhino.jstype.InstanceObjectType.Builder builder)
func NewInstanceObjectType(builder *InstanceObjectType_Builder) *InstanceObjectType {
	this_ := &InstanceObjectType{}
	this_.BindSelf_(this_)
	InitInstanceObjectType(this_, builder)
	return this_
}

func InitInstanceObjectType(this_ *InstanceObjectType, builder *InstanceObjectType_Builder) {
	_ = builder
	InitPrototypeObjectType_1(&this_.PrototypeObjectTypeBase, ToIface_PrototypeObjectType_Builder(builder))
	this_.initFields_InstanceObjectType()
	this_.constructor = j2gort.CheckNotNullPtr(builder.constructor)
	this_.GetJSTypeBase().registry.GetResolver().ResolveIfClosed(ToIface_JSType(this_), StaticGet_InstanceObjectType_TYPE_CLASS())
}

// Java static method: InstanceObjectType.builderForCtor(com.google.javascript.rhino.jstype.FunctionType ctor)
func InstanceObjectType_BuilderForCtor(ctor *FunctionType) *InstanceObjectType_Builder {
	_ = ctor
	return j2gort.CastPtr[InstanceObjectType_Builder](j2gort.CastPtr[InstanceObjectType_Builder](j2gort.CastPtr[InstanceObjectType_Builder](j2gort.CastPtr[InstanceObjectType_Builder](j2gort.CastPtr[InstanceObjectType_Builder](j2gort.CastPtr[InstanceObjectType_Builder](NewInstanceObjectType_Builder(ctor.GetJSTypeBase().registry).SetName(ctor.GetReferenceName())).SetImplicitPrototype(nil)).SetNative(ctor.IsNativeObjectType())).SetAnonymous(ctor.IsAnonymous())).SetTemplateTypeMap(ctor.GetTemplateTypeMap())).SetTemplateParamCount(ctor.GetTemplateParamCount())).SetConstructor(ctor)
}

// Java method: InstanceObjectType.getTypeClass()
func (this_ *InstanceObjectType) GetTypeClass() JSTypeClass {
	return StaticGet_InstanceObjectType_TYPE_CLASS()
}

// Java method: InstanceObjectType.getReferenceName()
func (this_ *InstanceObjectType) GetReferenceName() string {
	return this_.GetConstructor().GetReferenceName()
}

// Java method: InstanceObjectType.getImplicitPrototype()
func (this_ *InstanceObjectType) GetImplicitPrototype() ObjectType {
	return this_.GetConstructor().GetPrototype()
}

// Java method: InstanceObjectType.getConstructor()
func (this_ *InstanceObjectType) GetConstructor() *FunctionType {
	return this_.constructor
}

// Java method: InstanceObjectType.defineProperty(java.lang.String name, com.google.javascript.rhino.jstype.JSType type, boolean inferred, com.google.javascript.rhino.Node propertyNode)
func (this_ *InstanceObjectType) DefineProperty_4_string_JSType_bool_Node(name string, type_ JSType, inferred bool, propertyNode Node) bool {
	_ = name
	_ = type_
	_ = inferred
	_ = propertyNode
	var proto ObjectType = this_.GetImplicitPrototype()
	_ = proto
	if (proto != nil && !proto.IsJavaNil_()) && proto.GetObjectTypeBase().HasOwnDeclaredProperty_1_string(name) {
		return false
	}
	return this_.PrototypeObjectTypeBase.DefineProperty_4_string_JSType_bool_Node(name, type_, inferred, propertyNode)
}

// Java method: InstanceObjectType.defineProperty(com.google.javascript.rhino.jstype.Property.Key name, com.google.javascript.rhino.jstype.JSType type, boolean inferred, com.google.javascript.rhino.Node propertyNode)
func (this_ *InstanceObjectType) DefineProperty_4_Property_Key_JSType_bool_Node(name Property_Key, type_ JSType, inferred bool, propertyNode Node) bool {
	_ = name
	_ = type_
	_ = inferred
	_ = propertyNode
	var proto ObjectType = this_.GetImplicitPrototype()
	_ = proto
	if (proto != nil && !proto.IsJavaNil_()) && proto.GetObjectTypeBase().HasOwnDeclaredProperty_1_Property_Key(name) {
		return false
	}
	return this_.PrototypeObjectTypeBase.DefineProperty_4_Property_Key_JSType_bool_Node(name, type_, inferred, propertyNode)
}

// Java method: InstanceObjectType.appendTo(com.google.javascript.rhino.jstype.TypeStringBuilder sb)
func (this_ *InstanceObjectType) AppendTo(sb *TypeStringBuilder) {
	_ = sb
	if !(this_.constructor.HasReferenceName()) {
		this_.PrototypeObjectTypeBase.AppendTo(sb)
		return
	} else if sb.IsForAnnotations() {
		sb.Append_1_string(this_.constructor.GetNormalizedReferenceName())
		return
	}
	var name string = this_.constructor.GetReferenceName()
	_ = name
	if j2gort.CheckNotNullString(name) == "" {
		var n Node = this_.constructor.GetSource()
		_ = n
		sb.Append_1_string("<anonymous@").Append_1_string(func() string {
			if n != 0 {
				return n.GetSourceFileName()
			}
			return "unknown"
		}()).Append_1_string(":").Append_1_string(strconv.FormatInt(int64(func() int32 {
			if n != 0 {
				return n.GetLineno()
			}
			return int32(0)
		}()), 10)).Append_1_string(">")
	}
	sb.Append_1_string(name)
}

// Java method: InstanceObjectType.isTheObjectType()
func (this_ *InstanceObjectType) IsTheObjectType() bool {
	return (this_.GetConstructor().IsNativeObjectType() && func() bool {
		sr_1 := "Object"
		se_2 := this_.GetReferenceName()
		return j2gort.CheckNotNullString(sr_1) == se_2
	}())
}

// Java method: InstanceObjectType.isInstanceType()
func (this_ *InstanceObjectType) IsInstanceType() bool {
	return true
}

// Java method: InstanceObjectType.isArrayType()
func (this_ *InstanceObjectType) IsArrayType() bool {
	return (this_.GetConstructor().IsNativeObjectType() && func() bool {
		sr_3 := "Array"
		se_4 := this_.GetReferenceName()
		return j2gort.CheckNotNullString(sr_3) == se_4
	}())
}

// Java method: InstanceObjectType.isReadonlyArrayType()
func (this_ *InstanceObjectType) IsReadonlyArrayType() bool {
	return (this_.GetConstructor().IsNativeObjectType() && func() bool {
		sr_5 := "ReadonlyArray"
		se_6 := this_.GetReferenceName()
		return j2gort.CheckNotNullString(sr_5) == se_6
	}())
}

// Java method: InstanceObjectType.isBigIntObjectType()
func (this_ *InstanceObjectType) IsBigIntObjectType() bool {
	return (this_.GetConstructor().IsNativeObjectType() && func() bool {
		sr_7 := "BigInt"
		se_8 := this_.GetReferenceName()
		return j2gort.CheckNotNullString(sr_7) == se_8
	}())
}

// Java method: InstanceObjectType.isStringObjectType()
func (this_ *InstanceObjectType) IsStringObjectType() bool {
	return (this_.GetConstructor().IsNativeObjectType() && func() bool {
		sr_9 := "String"
		se_10 := this_.GetReferenceName()
		return j2gort.CheckNotNullString(sr_9) == se_10
	}())
}

// Java method: InstanceObjectType.isSymbolObjectType()
func (this_ *InstanceObjectType) IsSymbolObjectType() bool {
	return (this_.GetConstructor().IsNativeObjectType() && func() bool {
		sr_11 := "Symbol"
		se_12 := this_.GetReferenceName()
		return j2gort.CheckNotNullString(sr_11) == se_12
	}())
}

// Java method: InstanceObjectType.isBooleanObjectType()
func (this_ *InstanceObjectType) IsBooleanObjectType() bool {
	return (this_.GetConstructor().IsNativeObjectType() && func() bool {
		sr_13 := "Boolean"
		se_14 := this_.GetReferenceName()
		return j2gort.CheckNotNullString(sr_13) == se_14
	}())
}

// Java method: InstanceObjectType.isNumberObjectType()
func (this_ *InstanceObjectType) IsNumberObjectType() bool {
	return (this_.GetConstructor().IsNativeObjectType() && func() bool {
		sr_15 := "Number"
		se_16 := this_.GetReferenceName()
		return j2gort.CheckNotNullString(sr_15) == se_16
	}())
}

// Java method: InstanceObjectType.isDateType()
func (this_ *InstanceObjectType) IsDateType() bool {
	return (this_.GetConstructor().IsNativeObjectType() && func() bool {
		sr_17 := "Date"
		se_18 := this_.GetReferenceName()
		return j2gort.CheckNotNullString(sr_17) == se_18
	}())
}

// Java method: InstanceObjectType.isRegexpType()
func (this_ *InstanceObjectType) IsRegexpType() bool {
	return (this_.GetConstructor().IsNativeObjectType() && func() bool {
		sr_19 := "RegExp"
		se_20 := this_.GetReferenceName()
		return j2gort.CheckNotNullString(sr_19) == se_20
	}())
}

// Java method: InstanceObjectType.isNominalType()
func (this_ *InstanceObjectType) IsNominalType() bool {
	return this_.HasReferenceName()
}

// Java method: InstanceObjectType.recursionUnsafeHashCode()
func (this_ *InstanceObjectType) RecursionUnsafeHashCode() int32 {
	if this_.HasReferenceName() {
		return NamedType_NominalHashCode(this_)
	} else {
		return this_.PrototypeObjectTypeBase.HashCode()
	}
}

// Java method: InstanceObjectType.getCtorImplementedInterfaces()
func (this_ *InstanceObjectType) GetCtorImplementedInterfaces() any {
	return this_.GetConstructor().GetImplementedInterfaces()
}

// Java method: InstanceObjectType.getCtorExtendedInterfaces()
func (this_ *InstanceObjectType) GetCtorExtendedInterfaces() any {
	return this_.GetConstructor().GetExtendedInterfaces()
}

// Java method: InstanceObjectType.resolveInternal(com.google.javascript.rhino.ErrorReporter reporter)
func (this_ *InstanceObjectType) ResolveInternal(reporter ErrorReporter) JSType {
	_ = reporter
	var resolved JSType = this_.PrototypeObjectTypeBase.ResolveInternal(reporter)
	_ = resolved
	var unused JSType = JSType_SafeResolve(ToIface_JSType(this_.constructor), reporter)
	_ = unused
	return resolved
}

// Java class: com.google.javascript.rhino.jstype.InstanceObjectType.Builder
type InstanceObjectType_Builder struct {
	PrototypeObjectType_BuilderBase
	constructor *FunctionType `gson:"constructor"`
}

func (this_ *InstanceObjectType_Builder) BindSelf_(s any) {
	_ = s
	this_.PrototypeObjectType_BuilderBase.BindSelf_(s)
}
func (this_ *InstanceObjectType_Builder) IsJavaNil_() bool { return this_ == nil }
func (*InstanceObjectType_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.rhino.jstype.InstanceObjectType$Builder", "Builder"
}
func (this_ *InstanceObjectType_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.PrototypeObjectType_BuilderBase.IdentityHash_()
}
func (*InstanceObjectType_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.InstanceObjectType.Builder", "com.google.javascript.rhino.jstype.InstanceObjectType$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*InstanceObjectType_Builder](v_)
	})
}
func (*InstanceObjectType_Builder) Implements_rhino_PrototypeObjectType_Builder() {}

func (this_ *InstanceObjectType_Builder) initFields_InstanceObjectType_Builder() {
	_ = this_
}

// Java constructor: Builder(com.google.javascript.rhino.jstype.JSTypeRegistry registry)
func NewInstanceObjectType_Builder(registry *JSTypeRegistry) *InstanceObjectType_Builder {
	this_ := &InstanceObjectType_Builder{}
	this_.BindSelf_(this_)
	InitInstanceObjectType_Builder(this_, registry)
	return this_
}

func InitInstanceObjectType_Builder(this_ *InstanceObjectType_Builder, registry *JSTypeRegistry) {
	_ = registry
	InitPrototypeObjectType_Builder(&this_.PrototypeObjectType_BuilderBase, registry)
	this_.initFields_InstanceObjectType_Builder()
}

// Java method: Builder.setConstructor(com.google.javascript.rhino.jstype.FunctionType x)
func (this_ *InstanceObjectType_Builder) SetConstructor(x *FunctionType) *InstanceObjectType_Builder {
	_ = x
	this_.constructor = x
	return this_
}

// Java method: Builder.build()
func (this_ *InstanceObjectType_Builder) Build() any {
	return NewInstanceObjectType(this_)
}
