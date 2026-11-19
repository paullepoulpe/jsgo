/*
 * Copyright 2017 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/InvalidatingTypes.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.InvalidatingTypes
type InvalidatingTypes struct {
	typeToLocation *j2gort.Ext_com_google_common_collect_ImmutableSetMultimap `gson:"typeToLocation"`
	idHash_        atomic.Int32
}

func (this_ *InvalidatingTypes) BindSelf_(s any) {
	_ = s
}
func (this_ *InvalidatingTypes) IsJavaNil_() bool { return this_ == nil }
func (*InvalidatingTypes) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InvalidatingTypes", "InvalidatingTypes"
}
func (this_ *InvalidatingTypes) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*InvalidatingTypes) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InvalidatingTypes", "com.google.javascript.jscomp.InvalidatingTypes", func(v_ any) bool {
		return j2gort.InstanceOf[*InvalidatingTypes](v_)
	})
}

var static_InvalidatingTypes_ALWAYS_INVALIDATING_TYPES j2gort.Ext_com_google_common_collect_ImmutableList

func StaticGet_InvalidatingTypes_ALWAYS_INVALIDATING_TYPES() j2gort.Ext_com_google_common_collect_ImmutableList {
	clinit_InvalidatingTypes()
	return static_InvalidatingTypes_ALWAYS_INVALIDATING_TYPES
}
func StaticSet_InvalidatingTypes_ALWAYS_INVALIDATING_TYPES(v j2gort.Ext_com_google_common_collect_ImmutableList) {
	clinit_InvalidatingTypes()
	static_InvalidatingTypes_ALWAYS_INVALIDATING_TYPES = v
}

var static_InvalidatingTypes_ALWAYS_INVALIDATING_LOCATION pkg_rhino.Node

func StaticGet_InvalidatingTypes_ALWAYS_INVALIDATING_LOCATION() pkg_rhino.Node {
	clinit_InvalidatingTypes()
	return static_InvalidatingTypes_ALWAYS_INVALIDATING_LOCATION
}
func StaticSet_InvalidatingTypes_ALWAYS_INVALIDATING_LOCATION(v pkg_rhino.Node) {
	clinit_InvalidatingTypes()
	static_InvalidatingTypes_ALWAYS_INVALIDATING_LOCATION = v
}

var clinitDone_InvalidatingTypes uint32
var clinitState_InvalidatingTypes j2gort.ClinitState

func clinit_InvalidatingTypes() {
	if atomic.LoadUint32(&clinitDone_InvalidatingTypes) != 0 {
		return
	}
	clinitSlow_InvalidatingTypes()
}

func clinitSlow_InvalidatingTypes() {
	if !j2gort.EnterClinit(&clinitState_InvalidatingTypes, &clinitDone_InvalidatingTypes) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_InvalidatingTypes, &clinitDone_InvalidatingTypes)
	static_InvalidatingTypes_ALWAYS_INVALIDATING_TYPES = j2gort.Ext_com_google_common_collect_ImmutableList_Of_7(j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_FUNCTION_FUNCTION_TYPE()), j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_FUNCTION_TYPE()), j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_FUNCTION_PROTOTYPE()), j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_FUNCTION_INSTANCE_PROTOTYPE()), j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_OBJECT_TYPE()), j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_OBJECT_PROTOTYPE()), j2gort.BoxNullableVal(pkg_rhino.StaticGet_JSTypeNative_OBJECT_FUNCTION_TYPE()))
	static_InvalidatingTypes_ALWAYS_INVALIDATING_LOCATION = pkg_rhino.IR_Name("alwaysInvalidatingLocation").SetStaticSourceFile(pkg_rhino.ToIface_StaticSourceFile(SourceFile_FromCode_2("InvalidatingTypes_alwaysInvalidatingLocation", "")))
}

func (this_ *InvalidatingTypes) initFields_InvalidatingTypes() {
	_ = this_
}

// Java constructor: InvalidatingTypes(com.google.common.collect.ImmutableSetMultimap<com.google.javascript.rhino.jstype.JSType,com.google.javascript.rhino.Node> typeToLocation)
func NewInvalidatingTypes(typeToLocation *j2gort.Ext_com_google_common_collect_ImmutableSetMultimap) *InvalidatingTypes {
	clinit_InvalidatingTypes()
	this_ := &InvalidatingTypes{}
	InitInvalidatingTypes(this_, typeToLocation)
	return this_
}

func InitInvalidatingTypes(this_ *InvalidatingTypes, typeToLocation *j2gort.Ext_com_google_common_collect_ImmutableSetMultimap) {
	_ = typeToLocation
	this_.initFields_InvalidatingTypes()
	this_.typeToLocation = typeToLocation
}

// Java method: InvalidatingTypes.isInvalidating(com.google.javascript.rhino.jstype.JSType type)
func (this_ *InvalidatingTypes) IsInvalidating(type_ pkg_rhino.JSType) bool {
	_ = type_
	if ((type_ == nil || type_.IsJavaNil_()) || type_.IsUnknownType()) || type_.IsEmptyType() {
		return true
	}
	if type_.IsUnionType() {
		type_ = type_.RestrictByNotNullOrUndefined()
		if type_.IsUnionType() {
			{
				cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(type_.GetUnionMembers()))
				for cur_1.HasNext() {
					rawElem_2 := cur_1.Next()
					alt := j2gort.AsIface[pkg_rhino.JSType](rawElem_2)
					_ = alt
					if this_.IsInvalidating(alt) {
						return true
					}
				}
			}
			return false
		}
	}
	var objType pkg_rhino.ObjectType = type_.ToMaybeObjectType()
	_ = objType
	if objType == nil || objType.IsJavaNil_() {
		return false
	}
	if objType.IsTemplatizedType() {
		objType = objType.ToMaybeTemplatizedType().GetReferencedType()
	}
	return (this_.typeToLocation.ContainsKey(objType) || InvalidatingTypes_IsAmbiguousOrStructuralType(objType))
}

// Java method: InvalidatingTypes.getMismatchLocations()
func (this_ *InvalidatingTypes) GetMismatchLocations() *j2gort.Ext_com_google_common_collect_ImmutableSetMultimap {
	return this_.typeToLocation
}

// Java static method: InvalidatingTypes.isAmbiguousOrStructuralType(com.google.javascript.rhino.jstype.ObjectType type)
func InvalidatingTypes_IsAmbiguousOrStructuralType(type_ pkg_rhino.ObjectType) bool {
	clinit_InvalidatingTypes()
	_ = type_
	if type_.IsEnumType() {
		return false
	} else if type_.IsEnumElementType() {
		var primitive pkg_rhino.ObjectType = type_.ToMaybeEnumElementType().GetPrimitiveType().ToMaybeObjectType()
		_ = primitive
		return ((primitive == nil || primitive.IsJavaNil_()) || InvalidatingTypes_IsAmbiguousOrStructuralType(primitive))
	} else if type_.IsFunctionType() {
		return (!(type_.IsNominalConstructorOrInterface()) || type_.ToMaybeFunctionType().IsAmbiguousConstructor())
	} else if type_.IsFunctionPrototypeType() {
		var ownerFunction *pkg_rhino.FunctionType = type_.GetOwnerFunction()
		_ = ownerFunction
		return ((((ownerFunction == nil) || !(ownerFunction.IsNominalConstructorOrInterface())) || ownerFunction.IsAmbiguousConstructor()) || ownerFunction.IsStructuralInterface())
	} else if type_.IsInstanceType() {
		var ctor *pkg_rhino.FunctionType = type_.GetConstructor()
		_ = ctor
		return (((ctor == nil) || ctor.IsAmbiguousConstructor()) || ctor.IsStructuralInterface())
	}
	return true
}

// Java class: com.google.javascript.jscomp.InvalidatingTypes.Builder
type InvalidatingTypes_Builder struct {
	registry       *pkg_rhino.JSTypeRegistry                                          `gson:"registry"`
	typeToLocation *j2gort.Ext_com_google_common_collect_ImmutableSetMultimap_Builder `gson:"typeToLocation"`
	idHash_        atomic.Int32
}

func (this_ *InvalidatingTypes_Builder) BindSelf_(s any) {
	_ = s
}
func (this_ *InvalidatingTypes_Builder) IsJavaNil_() bool { return this_ == nil }
func (*InvalidatingTypes_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InvalidatingTypes$Builder", "Builder"
}
func (this_ *InvalidatingTypes_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*InvalidatingTypes_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InvalidatingTypes.Builder", "com.google.javascript.jscomp.InvalidatingTypes$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*InvalidatingTypes_Builder](v_)
	})
}

func (this_ *InvalidatingTypes_Builder) initFields_InvalidatingTypes_Builder() {
	_ = this_
	this_.typeToLocation = j2gort.Ext_com_google_common_collect_ImmutableSetMultimap_S_Builder()
}

// Java constructor: Builder(com.google.javascript.rhino.jstype.JSTypeRegistry registry)
func NewInvalidatingTypes_Builder(registry *pkg_rhino.JSTypeRegistry) *InvalidatingTypes_Builder {
	this_ := &InvalidatingTypes_Builder{}
	InitInvalidatingTypes_Builder(this_, registry)
	return this_
}

func InitInvalidatingTypes_Builder(this_ *InvalidatingTypes_Builder, registry *pkg_rhino.JSTypeRegistry) {
	_ = registry
	this_.initFields_InvalidatingTypes_Builder()
	this_.registry = registry
}

// Java method: Builder.build()
func (this_ *InvalidatingTypes_Builder) Build() *InvalidatingTypes {
	{
		cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(StaticGet_InvalidatingTypes_ALWAYS_INVALIDATING_TYPES()))
		for cur_3.HasNext() {
			rawElem_4 := cur_3.Next()
			t := j2gort.Cast[pkg_rhino.JSTypeNative](rawElem_4)
			_ = t
			this_.typeToLocation.Put_2(this_.registry.GetNativeType(t), j2gort.BoxNode(StaticGet_InvalidatingTypes_ALWAYS_INVALIDATING_LOCATION()))
		}
	}
	return NewInvalidatingTypes(j2gort.CastPtr[j2gort.Ext_com_google_common_collect_ImmutableSetMultimap](this_.typeToLocation.Build()))
}

// Java method: Builder.addAllTypeMismatches(java.lang.Iterable<com.google.javascript.jscomp.TypeMismatch> mismatches)
func (this_ *InvalidatingTypes_Builder) AddAllTypeMismatches(mismatches j2gort.Ext_java_lang_Iterable) *InvalidatingTypes_Builder {
	_ = mismatches
	{
		cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(mismatches))
		for cur_5.HasNext() {
			rawElem_6 := cur_5.Next()
			mismatch := j2gort.CastPtr[TypeMismatch](rawElem_6)
			_ = mismatch
			this_.priv_Builder_addTypeWithReason_3431(mismatch.Found(), mismatch.Location())
			this_.priv_Builder_addTypeWithReason_3431(mismatch.Required(), mismatch.Location())
		}
	}
	return this_
}

// Java method: Builder.addTypeWithReason(com.google.javascript.rhino.jstype.JSType type, com.google.javascript.rhino.Node location)
func (this_ *InvalidatingTypes_Builder) priv_Builder_addTypeWithReason_3431(type_ pkg_rhino.JSType, location pkg_rhino.Node) {
	_ = type_
	_ = location
	type_ = type_.RestrictByNotNullOrUndefined()
	if type_.IsUnionType() {
		{
			cur_7 := j2gort.NewIterCursor(j2gort.CheckNotNull(type_.GetUnionMembers()))
			for cur_7.HasNext() {
				rawElem_8 := cur_7.Next()
				alt := j2gort.AsIface[pkg_rhino.JSType](rawElem_8)
				_ = alt
				this_.priv_Builder_addTypeWithReason_3431(alt, location)
			}
		}
		return
	}
	if type_.IsEnumElementType() {
		this_.priv_Builder_addTypeWithReason_3431(type_.GetEnumeratedTypeOfEnumElement(), location)
		return
	}
	var objType pkg_rhino.ObjectType = type_.ToMaybeObjectType()
	_ = objType
	if objType == nil || objType.IsJavaNil_() {
		return
	}
	this_.priv_Builder_recordTypeWithReason_3432(objType, location)
	this_.priv_Builder_recordTypeWithReason_3432(objType.GetImplicitPrototype(), location)
	if objType.IsConstructor() {
		var instanceType pkg_rhino.ObjectType = objType.ToMaybeFunctionType().GetInstanceType()
		_ = instanceType
		this_.priv_Builder_recordTypeWithReason_3432(instanceType, location)
	} else if objType.IsInstanceType() {
		this_.priv_Builder_recordTypeWithReason_3432(pkg_rhino.ToIface_ObjectType(objType.GetConstructor()), location)
	}
}

// Java method: Builder.recordTypeWithReason(com.google.javascript.rhino.jstype.ObjectType type, com.google.javascript.rhino.Node location)
func (this_ *InvalidatingTypes_Builder) priv_Builder_recordTypeWithReason_3432(type_ pkg_rhino.ObjectType, location pkg_rhino.Node) {
	_ = type_
	_ = location
	if type_ == nil || type_.IsJavaNil_() {
		return
	}
	if type_.IsTemplatizedType() {
		type_ = type_.ToMaybeTemplatizedType().GetReferencedType()
	}
	if InvalidatingTypes_IsAmbiguousOrStructuralType(type_) {
		return
	}
	this_.typeToLocation.Put_2(type_, j2gort.BoxNode(location))
}
