/*
 * Copyright 2019 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
// Port of com/google/javascript/jscomp/AccessorSummary.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.AccessorSummary
type AccessorSummary struct {
	accessors                   j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"accessors"`
	idHash_                     atomic.Int32
	assumeAlwaysGetterAndSetter bool `gson:"assumeAlwaysGetterAndSetter"`
}

func (this_ *AccessorSummary) BindSelf_(s any) {
	_ = s
}
func (this_ *AccessorSummary) IsJavaNil_() bool { return this_ == nil }
func (*AccessorSummary) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AccessorSummary", "AccessorSummary"
}
func (this_ *AccessorSummary) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*AccessorSummary) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AccessorSummary", "com.google.javascript.jscomp.AccessorSummary", func(v_ any) bool {
		return j2gort.InstanceOf[*AccessorSummary](v_)
	})
}

func (this_ *AccessorSummary) initFields_AccessorSummary() {
	_ = this_
}

// Java constructor: AccessorSummary(com.google.common.collect.ImmutableMap<java.lang.String,com.google.javascript.jscomp.AccessorSummary.PropertyAccessKind> accessors)
func NewAccessorSummary_1_Ext_com_google_common_collect_ImmutableMap(accessors j2gort.Ext_com_google_common_collect_ImmutableMap) *AccessorSummary {
	this_ := &AccessorSummary{}
	InitAccessorSummary_1_Ext_com_google_common_collect_ImmutableMap(this_, accessors)
	return this_
}

func InitAccessorSummary_1_Ext_com_google_common_collect_ImmutableMap(this_ *AccessorSummary, accessors j2gort.Ext_com_google_common_collect_ImmutableMap) {
	_ = accessors
	this_.initFields_AccessorSummary()
	this_.accessors = accessors
	this_.assumeAlwaysGetterAndSetter = false
}

// Java constructor: AccessorSummary(boolean assumeAlwaysGetterAndSetter)
func NewAccessorSummary_1_bool(assumeAlwaysGetterAndSetter bool) *AccessorSummary {
	this_ := &AccessorSummary{}
	InitAccessorSummary_1_bool(this_, assumeAlwaysGetterAndSetter)
	return this_
}

func InitAccessorSummary_1_bool(this_ *AccessorSummary, assumeAlwaysGetterAndSetter bool) {
	_ = assumeAlwaysGetterAndSetter
	this_.initFields_AccessorSummary()
	this_.accessors = j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0()
	this_.assumeAlwaysGetterAndSetter = assumeAlwaysGetterAndSetter
}

// Java static method: AccessorSummary.create(java.util.Map<java.lang.String,com.google.javascript.jscomp.AccessorSummary.PropertyAccessKind> accessors)
func AccessorSummary_Create(accessors j2gort.Ext_java_util_Map) *AccessorSummary {
	_ = accessors
	return NewAccessorSummary_1_Ext_com_google_common_collect_ImmutableMap(j2gort.Ext_com_google_common_collect_ImmutableMap_CopyOf_1_Ext_java_util_Map(accessors))
}

// Java method: AccessorSummary.getAccessors()
func (this_ *AccessorSummary) GetAccessors() j2gort.Ext_com_google_common_collect_ImmutableMap {
	return this_.accessors
}

// Java method: AccessorSummary.getKind(java.lang.String name)
func (this_ *AccessorSummary) GetKind(name string) *AccessorSummary_PropertyAccessKind {
	_ = name
	if this_.assumeAlwaysGetterAndSetter {
		return StaticGet_AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER()
	}
	return j2gort.CastPtr[AccessorSummary_PropertyAccessKind](j2gort.MapGetOrDefaultStr(this_.accessors, name, j2gort.PtrToAny(StaticGet_AccessorSummary_PropertyAccessKind_NORMAL())))
}

// Java static method: AccessorSummary.createAssumingAlwaysGetterAndSetter()
func AccessorSummary_CreateAssumingAlwaysGetterAndSetter() *AccessorSummary {
	return NewAccessorSummary_1_bool(true)
}

// Java method: AccessorSummary.toProto()
func (this_ *AccessorSummary) ToProto() *j2gort.Ext_com_google_javascript_jscomp_AccessorSummaryProto {
	var builder *j2gort.Ext_com_google_javascript_jscomp_AccessorSummaryProto_Builder = j2gort.Ext_com_google_javascript_jscomp_AccessorSummaryProto_NewBuilder_0().SetAssumeAlwaysGetterAndSetter(this_.assumeAlwaysGetterAndSetter)
	_ = builder
	if this_.assumeAlwaysGetterAndSetter {
		if !(this_.accessors.IsEmpty()) {
			j2gort.CheckState(false)
		}
	} else {
		{
			cur_1 := j2gort.NewMapEntryCursor(this_.accessors)
			for cur_1.NextEntry() {
				builder.AddAccessors_1_Ext_com_google_javascript_jscomp_AccessorSummaryEntryProto(j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_AccessorSummaryEntryProto](j2gort.Ext_com_google_javascript_jscomp_AccessorSummaryEntryProto_NewBuilder_0().SetName(cur_1.KeyStr()).SetKind(AccessorSummary_ToProto(j2gort.CastPtr[AccessorSummary_PropertyAccessKind](cur_1.GetValue()))).Build()))
			}
		}
	}
	return j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_AccessorSummaryProto](builder.Build())
}

// Java static method: AccessorSummary.toProto(com.google.javascript.jscomp.AccessorSummary.PropertyAccessKind kind)
func AccessorSummary_ToProto(kind *AccessorSummary_PropertyAccessKind) *j2gort.Ext_com_google_javascript_jscomp_PropertyAccessKindProto {
	_ = kind
	switch j2gort.CheckNotNullPtr(kind).Ordinal() {
	case int32(0) /* PropertyAccessKind.NORMAL */ :
		return j2gort.Ext_com_google_javascript_jscomp_PropertyAccessKindProto_F_KIND_NORMAL
	case int32(1) /* PropertyAccessKind.GETTER_ONLY */ :
		return j2gort.Ext_com_google_javascript_jscomp_PropertyAccessKindProto_F_KIND_GETTER_ONLY
	case int32(2) /* PropertyAccessKind.SETTER_ONLY */ :
		return j2gort.Ext_com_google_javascript_jscomp_PropertyAccessKindProto_F_KIND_SETTER_ONLY
	case int32(3) /* PropertyAccessKind.GETTER_AND_SETTER */ :
		return j2gort.Ext_com_google_javascript_jscomp_PropertyAccessKindProto_F_KIND_GETTER_AND_SETTER
	}
	panic("unreachable switch expression")
}

// Java static method: AccessorSummary.fromProto(com.google.javascript.jscomp.AccessorSummaryProto proto)
func AccessorSummary_FromProto_1_Ext_com_google_javascript_jscomp_AccessorSummaryProto(proto *j2gort.Ext_com_google_javascript_jscomp_AccessorSummaryProto) *AccessorSummary {
	_ = proto
	if proto.GetAssumeAlwaysGetterAndSetter() {
		return AccessorSummary_CreateAssumingAlwaysGetterAndSetter()
	}
	var builder j2gort.Ext_com_google_common_collect_ImmutableMap_Builder = j2gort.Ext_com_google_common_collect_ImmutableMap_S_Builder()
	_ = builder
	{
		cur_2 := j2gort.NewIterCursor(j2gort.CheckNotNull(proto.GetAccessorsList()))
		for cur_2.HasNext() {
			rawElem_3 := cur_2.Next()
			entry := j2gort.CastPtr[j2gort.Ext_com_google_javascript_jscomp_AccessorSummaryEntryProto](rawElem_3)
			_ = entry
			builder.Put_2(j2gort.BoxString(entry.GetName()), j2gort.PtrToAny(AccessorSummary_FromProto_1_Ext_com_google_javascript_jscomp_PropertyAccessKindProto(entry.GetKind())))
		}
	}
	return AccessorSummary_Create(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](builder.BuildOrThrow()))
}

// Java static method: AccessorSummary.fromProto(com.google.javascript.jscomp.PropertyAccessKindProto kind)
func AccessorSummary_FromProto_1_Ext_com_google_javascript_jscomp_PropertyAccessKindProto(kind *j2gort.Ext_com_google_javascript_jscomp_PropertyAccessKindProto) *AccessorSummary_PropertyAccessKind {
	_ = kind
	switch j2gort.CheckNotNullPtr(kind).Ordinal() {
	case int32(1) /* PropertyAccessKindProto.KIND_NORMAL */ :
		return StaticGet_AccessorSummary_PropertyAccessKind_NORMAL()
	case int32(2) /* PropertyAccessKindProto.KIND_GETTER_ONLY */ :
		return StaticGet_AccessorSummary_PropertyAccessKind_GETTER_ONLY()
	case int32(3) /* PropertyAccessKindProto.KIND_SETTER_ONLY */ :
		return StaticGet_AccessorSummary_PropertyAccessKind_SETTER_ONLY()
	case int32(4) /* PropertyAccessKindProto.KIND_GETTER_AND_SETTER */ :
		return StaticGet_AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER()
	default:
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Unknown kind: "+j2gort.StrVal(kind))))
	}
	panic("unreachable switch expression")
}

// Java enum: com.google.javascript.jscomp.AccessorSummary.PropertyAccessKind
type AccessorSummary_PropertyAccessKind struct {
	j2gort.EnumBase
	flags int8
}

func (this_ *AccessorSummary_PropertyAccessKind) GetEnumBase() *j2gort.EnumBase {
	if this_ == nil {
		return nil
	}
	return &this_.EnumBase
}
func (this_ *AccessorSummary_PropertyAccessKind) EnumOrdinal_() int32 {
	if this_ == nil {
		return -1
	}
	return this_.EnumBase.Ordinal_
}
func (*AccessorSummary_PropertyAccessKind) Implements_j2gort_Ext_java_lang_Enum()               {}
func (*AccessorSummary_PropertyAccessKind) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ *AccessorSummary_PropertyAccessKind) initFields_AccessorSummary_PropertyAccessKind() {
	_ = this_
}

func initEnumAccessorSummary_PropertyAccessKind(this_ *AccessorSummary_PropertyAccessKind, enumName_ string, enumOrd_ int32, flags int32) *AccessorSummary_PropertyAccessKind {
	this_.EnumBase.Name_ = enumName_
	this_.EnumBase.Ordinal_ = enumOrd_
	_ = flags
	this_.initFields_AccessorSummary_PropertyAccessKind()
	this_.flags = int8(flags)
	return this_
}

var enum_AccessorSummary_PropertyAccessKind_NORMAL *AccessorSummary_PropertyAccessKind

func StaticGet_AccessorSummary_PropertyAccessKind_NORMAL() *AccessorSummary_PropertyAccessKind {
	clinit_AccessorSummary_PropertyAccessKind()
	return enum_AccessorSummary_PropertyAccessKind_NORMAL
}

var enum_AccessorSummary_PropertyAccessKind_GETTER_ONLY *AccessorSummary_PropertyAccessKind

func StaticGet_AccessorSummary_PropertyAccessKind_GETTER_ONLY() *AccessorSummary_PropertyAccessKind {
	clinit_AccessorSummary_PropertyAccessKind()
	return enum_AccessorSummary_PropertyAccessKind_GETTER_ONLY
}

var enum_AccessorSummary_PropertyAccessKind_SETTER_ONLY *AccessorSummary_PropertyAccessKind

func StaticGet_AccessorSummary_PropertyAccessKind_SETTER_ONLY() *AccessorSummary_PropertyAccessKind {
	clinit_AccessorSummary_PropertyAccessKind()
	return enum_AccessorSummary_PropertyAccessKind_SETTER_ONLY
}

var enum_AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER *AccessorSummary_PropertyAccessKind

func StaticGet_AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER() *AccessorSummary_PropertyAccessKind {
	clinit_AccessorSummary_PropertyAccessKind()
	return enum_AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER
}

var clinitDone_AccessorSummary_PropertyAccessKind uint32
var clinitState_AccessorSummary_PropertyAccessKind j2gort.ClinitState

func clinit_AccessorSummary_PropertyAccessKind() {
	if atomic.LoadUint32(&clinitDone_AccessorSummary_PropertyAccessKind) != 0 {
		return
	}
	clinitSlow_AccessorSummary_PropertyAccessKind()
}

func clinitSlow_AccessorSummary_PropertyAccessKind() {
	if !j2gort.EnterClinit(&clinitState_AccessorSummary_PropertyAccessKind, &clinitDone_AccessorSummary_PropertyAccessKind) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_AccessorSummary_PropertyAccessKind, &clinitDone_AccessorSummary_PropertyAccessKind)
	enum_AccessorSummary_PropertyAccessKind_NORMAL = &AccessorSummary_PropertyAccessKind{}
	initEnumAccessorSummary_PropertyAccessKind(enum_AccessorSummary_PropertyAccessKind_NORMAL, "NORMAL", int32(0), int32(0))
	enum_AccessorSummary_PropertyAccessKind_GETTER_ONLY = &AccessorSummary_PropertyAccessKind{}
	initEnumAccessorSummary_PropertyAccessKind(enum_AccessorSummary_PropertyAccessKind_GETTER_ONLY, "GETTER_ONLY", int32(1), int32(1))
	enum_AccessorSummary_PropertyAccessKind_SETTER_ONLY = &AccessorSummary_PropertyAccessKind{}
	initEnumAccessorSummary_PropertyAccessKind(enum_AccessorSummary_PropertyAccessKind_SETTER_ONLY, "SETTER_ONLY", int32(2), int32(2))
	enum_AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER = &AccessorSummary_PropertyAccessKind{}
	initEnumAccessorSummary_PropertyAccessKind(enum_AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER, "GETTER_AND_SETTER", int32(3), int32(3))
	enumValues_AccessorSummary_PropertyAccessKind = []*AccessorSummary_PropertyAccessKind{enum_AccessorSummary_PropertyAccessKind_NORMAL, enum_AccessorSummary_PropertyAccessKind_GETTER_ONLY, enum_AccessorSummary_PropertyAccessKind_SETTER_ONLY, enum_AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.AccessorSummary.PropertyAccessKind", []any{enum_AccessorSummary_PropertyAccessKind_NORMAL, enum_AccessorSummary_PropertyAccessKind_GETTER_ONLY, enum_AccessorSummary_PropertyAccessKind_SETTER_ONLY, enum_AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER})
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.AccessorSummary$PropertyAccessKind", []any{enum_AccessorSummary_PropertyAccessKind_NORMAL, enum_AccessorSummary_PropertyAccessKind_GETTER_ONLY, enum_AccessorSummary_PropertyAccessKind_SETTER_ONLY, enum_AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER})
}

func init() {
	j2gort.RegisterEnumConstantsLoader("com.google.javascript.jscomp.AccessorSummary.PropertyAccessKind", "com.google.javascript.jscomp.AccessorSummary$PropertyAccessKind", func() {
		clinit_AccessorSummary_PropertyAccessKind()
	})
}

var enumValues_AccessorSummary_PropertyAccessKind []*AccessorSummary_PropertyAccessKind

func AccessorSummary_PropertyAccessKind_Values() []*AccessorSummary_PropertyAccessKind {
	clinit_AccessorSummary_PropertyAccessKind()
	return j2gort.CloneSlice(enumValues_AccessorSummary_PropertyAccessKind)
}

func AccessorSummary_PropertyAccessKind_ValueOf(name string) *AccessorSummary_PropertyAccessKind {
	clinit_AccessorSummary_PropertyAccessKind()
	for _, c := range enumValues_AccessorSummary_PropertyAccessKind {
		if c.Name() == name {
			return c
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.AccessorSummary.PropertyAccessKind." + name))
}

// Java method: PropertyAccessKind.hasGetter()
func (this_ *AccessorSummary_PropertyAccessKind) HasGetter() bool {
	return ((int32(this_.flags) & int32(1)) != int32(0))
}

// Java method: PropertyAccessKind.hasSetter()
func (this_ *AccessorSummary_PropertyAccessKind) HasSetter() bool {
	return ((int32(this_.flags) & int32(2)) != int32(0))
}

// Java method: PropertyAccessKind.hasGetterOrSetter()
func (this_ *AccessorSummary_PropertyAccessKind) HasGetterOrSetter() bool {
	return ((int32(this_.flags) & int32(3)) != int32(0))
}

// Java method: PropertyAccessKind.unionWith(com.google.javascript.jscomp.AccessorSummary.PropertyAccessKind other)
func (this_ *AccessorSummary_PropertyAccessKind) UnionWith(other *AccessorSummary_PropertyAccessKind) *AccessorSummary_PropertyAccessKind {
	_ = other
	var combinedFlags int32 = (int32(this_.flags) | int32(other.flags))
	_ = combinedFlags
	switch combinedFlags {
	case int32(0):
		return StaticGet_AccessorSummary_PropertyAccessKind_NORMAL()
	case int32(1):
		return StaticGet_AccessorSummary_PropertyAccessKind_GETTER_ONLY()
	case int32(2):
		return StaticGet_AccessorSummary_PropertyAccessKind_SETTER_ONLY()
	case int32(3):
		return StaticGet_AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER()
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "unexpected value: "+j2gort.StrValInt32(combinedFlags))))
	}
	panic("unreachable switch expression")
}

func (this_ *AccessorSummary_PropertyAccessKind) IsJavaNil_() bool { return this_ == nil }
func (*AccessorSummary_PropertyAccessKind) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AccessorSummary$PropertyAccessKind", "PropertyAccessKind"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AccessorSummary.PropertyAccessKind", "com.google.javascript.jscomp.AccessorSummary$PropertyAccessKind", func(v_ any) bool {
		return j2gort.InstanceOf[*AccessorSummary_PropertyAccessKind](v_)
	})
}
