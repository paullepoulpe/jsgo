/*
 * Copyright 2009 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/deps/SimpleDependencyInfo.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java polymorphic class interface: com.google.javascript.jscomp.deps.SimpleDependencyInfo
type SimpleDependencyInfo interface {
	Implements_jscomp_SimpleDependencyInfo()
	IsJavaNil_() bool
	GetSimpleDependencyInfoBase() *SimpleDependencyInfoBase
	DependencyInfo
}

func ToIface_SimpleDependencyInfo[T interface {
	SimpleDependencyInfo
	comparable
}](v T) SimpleDependencyInfo {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_SimpleDependencyInfo[T SimpleDependencyInfo](v T) SimpleDependencyInfo {
	var vi SimpleDependencyInfo = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.deps.SimpleDependencyInfo
type SimpleDependencyInfoBase struct {
	self_                 SimpleDependencyInfo
	cachedRequiredSymbols j2gort.Ext_com_google_common_collect_ImmutableList `gson:"-"`
	idHash_               atomic.Int32
}

func (this_ *SimpleDependencyInfoBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(SimpleDependencyInfo)
}
func (this_ *SimpleDependencyInfoBase) IsJavaNil_() bool { return this_ == nil }
func (*SimpleDependencyInfoBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.SimpleDependencyInfo", "SimpleDependencyInfo"
}
func (this_ *SimpleDependencyInfoBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.SimpleDependencyInfo", "com.google.javascript.jscomp.deps.SimpleDependencyInfo", func(v_ any) bool {
		return j2gort.InstanceOf[SimpleDependencyInfo](v_)
	})
}
func (b *SimpleDependencyInfoBase) GetSimpleDependencyInfoBase() *SimpleDependencyInfoBase { return b }
func (*SimpleDependencyInfoBase) Implements_jscomp_SimpleDependencyInfo()                  {}
func (*SimpleDependencyInfoBase) Implements_jscomp_DependencyInfo()                        {}

var static_SimpleDependencyInfo_EMPTY SimpleDependencyInfo

func StaticGet_SimpleDependencyInfo_EMPTY() SimpleDependencyInfo {
	clinit_SimpleDependencyInfo()
	return static_SimpleDependencyInfo_EMPTY
}
func StaticSet_SimpleDependencyInfo_EMPTY(v SimpleDependencyInfo) {
	clinit_SimpleDependencyInfo()
	static_SimpleDependencyInfo_EMPTY = v
}

var clinitDone_SimpleDependencyInfo uint32
var clinitState_SimpleDependencyInfo j2gort.ClinitState

func clinit_SimpleDependencyInfo() {
	if atomic.LoadUint32(&clinitDone_SimpleDependencyInfo) != 0 {
		return
	}
	clinitSlow_SimpleDependencyInfo()
}

func clinitSlow_SimpleDependencyInfo() {
	if !j2gort.EnterClinit(&clinitState_SimpleDependencyInfo, &clinitDone_SimpleDependencyInfo) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_SimpleDependencyInfo, &clinitDone_SimpleDependencyInfo)
	static_SimpleDependencyInfo_EMPTY = SimpleDependencyInfo_S_Builder("", "").Build()
}

func (this_ *SimpleDependencyInfoBase) initFields_SimpleDependencyInfo() {
	_ = this_
}

func InitSimpleDependencyInfo(this_ *SimpleDependencyInfoBase) {
	this_.initFields_SimpleDependencyInfo()
}

// Java method: SimpleDependencyInfo.getRequiredSymbols()
func (this_ *SimpleDependencyInfoBase) GetRequiredSymbols() j2gort.Ext_com_google_common_collect_ImmutableList {
	var syms j2gort.Ext_com_google_common_collect_ImmutableList = this_.cachedRequiredSymbols
	_ = syms
	if syms == nil {
		syms = DependencyInfo_Require_AsSymbolList(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.self_.GetRequires()))
		this_.cachedRequiredSymbols = syms
	}
	return syms
}

// Java static method: SimpleDependencyInfo.builder(java.lang.String srcPathRelativeToClosure, java.lang.String pathOfDefiningFile)
func SimpleDependencyInfo_S_Builder(srcPathRelativeToClosure string, pathOfDefiningFile string) SimpleDependencyInfo_Builder {
	clinit_SimpleDependencyInfo()
	_ = srcPathRelativeToClosure
	_ = pathOfDefiningFile
	return NewAutoValue_SimpleDependencyInfo_Builder().SetName(pathOfDefiningFile).SetPathRelativeToClosureBase(srcPathRelativeToClosure).SetProvides_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())).SetRequires_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())).SetTypeRequires_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())).SetLoadFlags(j2gort.AsIface[j2gort.Ext_java_util_Map](j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0())).SetHasExternsAnnotation(false).SetHasNoCompileAnnotation(false)
}

func (this_ *SimpleDependencyInfoBase) IsEs6Module() bool {
	return IfaceDefault_DependencyInfo_IsEs6Module(this_.self_)
}

func (this_ *SimpleDependencyInfoBase) IsGoogModule() bool {
	return IfaceDefault_DependencyInfo_IsGoogModule(this_.self_)
}

// Java polymorphic class interface: com.google.javascript.jscomp.deps.SimpleDependencyInfo.Builder
type SimpleDependencyInfo_Builder interface {
	Implements_jscomp_SimpleDependencyInfo_Builder()
	IsJavaNil_() bool
	GetSimpleDependencyInfo_BuilderBase() *SimpleDependencyInfo_BuilderBase
	SetName(string) SimpleDependencyInfo_Builder
	SetPathRelativeToClosureBase(string) SimpleDependencyInfo_Builder
	SetProvides_1_Ext_java_util_Collection(j2gort.Ext_java_util_Collection) SimpleDependencyInfo_Builder
	SetProvides_1_slice_string([]string) SimpleDependencyInfo_Builder
	SetRequires_1_Ext_java_util_Collection(j2gort.Ext_java_util_Collection) SimpleDependencyInfo_Builder
	SetRequires_1_DependencyInfo_Require([]*DependencyInfo_Require) SimpleDependencyInfo_Builder
	SetTypeRequires_1_Ext_java_util_Collection(j2gort.Ext_java_util_Collection) SimpleDependencyInfo_Builder
	SetTypeRequires_1_slice_string([]string) SimpleDependencyInfo_Builder
	SetLoadFlags(j2gort.Ext_java_util_Map) SimpleDependencyInfo_Builder
	SetHasExternsAnnotation(bool) SimpleDependencyInfo_Builder
	SetHasNoCompileAnnotation(bool) SimpleDependencyInfo_Builder
	SetGoogModule(bool) SimpleDependencyInfo_Builder
	Build() SimpleDependencyInfo
}

func ToIface_SimpleDependencyInfo_Builder[T interface {
	SimpleDependencyInfo_Builder
	comparable
}](v T) SimpleDependencyInfo_Builder {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_SimpleDependencyInfo_Builder[T SimpleDependencyInfo_Builder](v T) SimpleDependencyInfo_Builder {
	var vi SimpleDependencyInfo_Builder = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.deps.SimpleDependencyInfo.Builder
type SimpleDependencyInfo_BuilderBase struct {
	self_   SimpleDependencyInfo_Builder
	idHash_ atomic.Int32
}

func (this_ *SimpleDependencyInfo_BuilderBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(SimpleDependencyInfo_Builder)
}
func (this_ *SimpleDependencyInfo_BuilderBase) IsJavaNil_() bool { return this_ == nil }
func (*SimpleDependencyInfo_BuilderBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.SimpleDependencyInfo$Builder", "Builder"
}
func (this_ *SimpleDependencyInfo_BuilderBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.SimpleDependencyInfo.Builder", "com.google.javascript.jscomp.deps.SimpleDependencyInfo$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[SimpleDependencyInfo_Builder](v_)
	})
}
func (b *SimpleDependencyInfo_BuilderBase) GetSimpleDependencyInfo_BuilderBase() *SimpleDependencyInfo_BuilderBase {
	return b
}
func (*SimpleDependencyInfo_BuilderBase) Implements_jscomp_SimpleDependencyInfo_Builder() {}

var static_SimpleDependencyInfo_Builder_GOOG_MODULE_FLAGS j2gort.Ext_com_google_common_collect_ImmutableMap

func StaticGet_SimpleDependencyInfo_Builder_GOOG_MODULE_FLAGS() j2gort.Ext_com_google_common_collect_ImmutableMap {
	clinit_SimpleDependencyInfo_Builder()
	return static_SimpleDependencyInfo_Builder_GOOG_MODULE_FLAGS
}
func StaticSet_SimpleDependencyInfo_Builder_GOOG_MODULE_FLAGS(v j2gort.Ext_com_google_common_collect_ImmutableMap) {
	clinit_SimpleDependencyInfo_Builder()
	static_SimpleDependencyInfo_Builder_GOOG_MODULE_FLAGS = v
}

var clinitDone_SimpleDependencyInfo_Builder uint32
var clinitState_SimpleDependencyInfo_Builder j2gort.ClinitState

func clinit_SimpleDependencyInfo_Builder() {
	if atomic.LoadUint32(&clinitDone_SimpleDependencyInfo_Builder) != 0 {
		return
	}
	clinitSlow_SimpleDependencyInfo_Builder()
}

func clinitSlow_SimpleDependencyInfo_Builder() {
	if !j2gort.EnterClinit(&clinitState_SimpleDependencyInfo_Builder, &clinitDone_SimpleDependencyInfo_Builder) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_SimpleDependencyInfo_Builder, &clinitDone_SimpleDependencyInfo_Builder)
	static_SimpleDependencyInfo_Builder_GOOG_MODULE_FLAGS = j2gort.Ext_com_google_common_collect_ImmutableMap_Of_2("module", "goog")
}

func (this_ *SimpleDependencyInfo_BuilderBase) initFields_SimpleDependencyInfo_Builder() {
	_ = this_
}

func InitSimpleDependencyInfo_Builder(this_ *SimpleDependencyInfo_BuilderBase) {
	this_.initFields_SimpleDependencyInfo_Builder()
}

// Java static method: Builder.from(com.google.javascript.jscomp.deps.DependencyInfo copy)
func SimpleDependencyInfo_Builder_From(copy_ DependencyInfo) SimpleDependencyInfo_Builder {
	clinit_SimpleDependencyInfo_Builder()
	_ = copy_
	return NewAutoValue_SimpleDependencyInfo_Builder().SetName(copy_.GetName()).SetPathRelativeToClosureBase(copy_.GetPathRelativeToClosureBase()).SetProvides_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](copy_.GetProvides())).SetRequires_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](copy_.GetRequires())).SetTypeRequires_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](copy_.GetTypeRequires())).SetLoadFlags(j2gort.AsIface[j2gort.Ext_java_util_Map](copy_.GetLoadFlags())).SetHasExternsAnnotation(copy_.GetHasExternsAnnotation()).SetHasNoCompileAnnotation(copy_.GetHasNoCompileAnnotation())
}

// Java method: Builder.setName(java.lang.String name)
func (this_ *SimpleDependencyInfo_BuilderBase) SetName(name string) SimpleDependencyInfo_Builder {
	_ = name
	panic("abstract or native method invoked")
}

// Java method: Builder.setPathRelativeToClosureBase(java.lang.String srcPathRelativeToClosure)
func (this_ *SimpleDependencyInfo_BuilderBase) SetPathRelativeToClosureBase(srcPathRelativeToClosure string) SimpleDependencyInfo_Builder {
	_ = srcPathRelativeToClosure
	panic("abstract or native method invoked")
}

// Java method: Builder.setProvides(java.util.Collection<java.lang.String> provides)
func (this_ *SimpleDependencyInfo_BuilderBase) SetProvides_1_Ext_java_util_Collection(provides j2gort.Ext_java_util_Collection) SimpleDependencyInfo_Builder {
	_ = provides
	panic("abstract or native method invoked")
}

// Java method: Builder.setProvides(java.lang.String[] provides)
func (this_ *SimpleDependencyInfo_BuilderBase) SetProvides_1_slice_string(provides []string) SimpleDependencyInfo_Builder {
	_ = provides
	panic("abstract or native method invoked")
}

// Java method: Builder.setRequires(java.util.Collection<com.google.javascript.jscomp.deps.DependencyInfo.Require> requires)
func (this_ *SimpleDependencyInfo_BuilderBase) SetRequires_1_Ext_java_util_Collection(requires j2gort.Ext_java_util_Collection) SimpleDependencyInfo_Builder {
	_ = requires
	panic("abstract or native method invoked")
}

// Java method: Builder.setRequires(com.google.javascript.jscomp.deps.DependencyInfo.Require[] requires)
func (this_ *SimpleDependencyInfo_BuilderBase) SetRequires_1_DependencyInfo_Require(requires []*DependencyInfo_Require) SimpleDependencyInfo_Builder {
	_ = requires
	panic("abstract or native method invoked")
}

// Java method: Builder.setTypeRequires(java.util.Collection<java.lang.String> typeRequires)
func (this_ *SimpleDependencyInfo_BuilderBase) SetTypeRequires_1_Ext_java_util_Collection(typeRequires j2gort.Ext_java_util_Collection) SimpleDependencyInfo_Builder {
	_ = typeRequires
	panic("abstract or native method invoked")
}

// Java method: Builder.setTypeRequires(java.lang.String[] typeRequires)
func (this_ *SimpleDependencyInfo_BuilderBase) SetTypeRequires_1_slice_string(typeRequires []string) SimpleDependencyInfo_Builder {
	_ = typeRequires
	panic("abstract or native method invoked")
}

// Java method: Builder.setLoadFlags(java.util.Map<java.lang.String,java.lang.String> loadFlags)
func (this_ *SimpleDependencyInfo_BuilderBase) SetLoadFlags(loadFlags j2gort.Ext_java_util_Map) SimpleDependencyInfo_Builder {
	_ = loadFlags
	panic("abstract or native method invoked")
}

// Java method: Builder.setHasExternsAnnotation(boolean hasExternsAnnotation)
func (this_ *SimpleDependencyInfo_BuilderBase) SetHasExternsAnnotation(hasExternsAnnotation bool) SimpleDependencyInfo_Builder {
	_ = hasExternsAnnotation
	panic("abstract or native method invoked")
}

// Java method: Builder.setHasNoCompileAnnotation(boolean hasNoCompileAnnotation)
func (this_ *SimpleDependencyInfo_BuilderBase) SetHasNoCompileAnnotation(hasNoCompileAnnotation bool) SimpleDependencyInfo_Builder {
	_ = hasNoCompileAnnotation
	panic("abstract or native method invoked")
}

// Java method: Builder.setGoogModule(boolean isModule)
func (this_ *SimpleDependencyInfo_BuilderBase) SetGoogModule(isModule bool) SimpleDependencyInfo_Builder {
	_ = isModule
	return this_.self_.SetLoadFlags(func() j2gort.Ext_java_util_Map {
		if isModule {
			return j2gort.AsIface[j2gort.Ext_java_util_Map](StaticGet_SimpleDependencyInfo_Builder_GOOG_MODULE_FLAGS())
		}
		return j2gort.AsIface[j2gort.Ext_java_util_Map](j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0())
	}())
}

// Java method: Builder.build()
func (this_ *SimpleDependencyInfo_BuilderBase) Build() SimpleDependencyInfo {
	panic("abstract or native method invoked")
}
