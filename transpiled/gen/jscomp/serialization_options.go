/*
 * Copyright 2020 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/serialization/SerializationOptions.java (generated mechanically by j2go).

package jscomp

import (
	"fmt"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java record: com.google.javascript.jscomp.serialization.SerializationOptions
type SerializationOptions struct {
	runtimeLibraries j2gort.Ext_com_google_common_collect_ImmutableList `gson:"runtimeLibraries"`
	includeDebugInfo bool `gson:"includeDebugInfo"`
	runValidation bool `gson:"runValidation"`
}

func (this_ *SerializationOptions) BindSelf_(s any) {
	_ = s
}
func (this_ *SerializationOptions) IsJavaNil_() bool { return this_ == nil }
func (*SerializationOptions) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.serialization.SerializationOptions", "SerializationOptions" }
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.SerializationOptions", "com.google.javascript.jscomp.serialization.SerializationOptions", func(v_ any) bool {
		return j2gort.InstanceOf[*SerializationOptions](v_)
	})
}
func (*SerializationOptions) Implements_j2gort_Ext_java_lang_Record() {}

func (this_ *SerializationOptions) initFields_SerializationOptions() {
	_ = this_
}

// Java constructor: SerializationOptions(boolean includeDebugInfo, boolean runValidation, com.google.common.collect.ImmutableList<java.lang.String> runtimeLibraries)
func NewSerializationOptions(includeDebugInfo bool, runValidation bool, runtimeLibraries j2gort.Ext_com_google_common_collect_ImmutableList) *SerializationOptions {
	this_ := &SerializationOptions{}
	InitSerializationOptions(this_, includeDebugInfo, runValidation, runtimeLibraries)
	return this_
}

func InitSerializationOptions(this_ *SerializationOptions, includeDebugInfo bool, runValidation bool, runtimeLibraries j2gort.Ext_com_google_common_collect_ImmutableList) {
	_ = includeDebugInfo
	_ = runValidation
	_ = runtimeLibraries
	this_.initFields_SerializationOptions()
	j2gort.CheckNotNull(runtimeLibraries)
	this_.includeDebugInfo = includeDebugInfo
	this_.runValidation = runValidation
	this_.runtimeLibraries = runtimeLibraries
}

// Java static method: SerializationOptions.builder()
func SerializationOptions_S_Builder() SerializationOptions_Builder {
	return NewAutoBuilder_SerializationOptions_Builder().SetRunValidation(false).SetIncludeDebugInfo(false).SetRuntimeLibraries(j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
}

func (this_ *SerializationOptions) IncludeDebugInfo() bool { return this_.includeDebugInfo }
func (this_ *SerializationOptions) RunValidation() bool { return this_.runValidation }
func (this_ *SerializationOptions) RuntimeLibraries() j2gort.Ext_com_google_common_collect_ImmutableList { return this_.runtimeLibraries }
func (this_ *SerializationOptions) Equals(other_ any) bool {
	o_, ok_ := other_.(*SerializationOptions)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if this_.includeDebugInfo != o_.includeDebugInfo { return false }
	if this_.runValidation != o_.runValidation { return false }
	if !javalang.ObjectsEquals(this_.runtimeLibraries, o_.runtimeLibraries) { return false }
	return true
}
func (this_ *SerializationOptions) HashCode() int32 {
	var h_ int32 = 0
	if this_.includeDebugInfo { h_ = 31*h_ + 1231 } else { h_ = 31*h_ + 1237 }
	if this_.runValidation { h_ = 31*h_ + 1231 } else { h_ = 31*h_ + 1237 }
	h_ = 31*h_ + j2gort.ObjectHash(this_.runtimeLibraries)
	return h_
}
func (this_ *SerializationOptions) String() string {
	_ = fmt.Sprint
	return "SerializationOptions[" + "includeDebugInfo=" + j2gort.StrVal(j2gort.ValueOfBoolean(this_.includeDebugInfo)) + ", " + "runValidation=" + j2gort.StrVal(j2gort.ValueOfBoolean(this_.runValidation)) + ", " + "runtimeLibraries=" + j2gort.StrVal(this_.runtimeLibraries) + "]"
}

// Java polymorphic class interface: com.google.javascript.jscomp.serialization.SerializationOptions.Builder
type SerializationOptions_Builder interface {
	Implements_jscomp_SerializationOptions_Builder()
	IsJavaNil_() bool
	GetSerializationOptions_BuilderBase() *SerializationOptions_BuilderBase
	SetIncludeDebugInfo(bool) SerializationOptions_Builder
	SetRunValidation(bool) SerializationOptions_Builder
	SetRuntimeLibraries(j2gort.Ext_com_google_common_collect_ImmutableList) SerializationOptions_Builder
	Build() *SerializationOptions
}
func ToIface_SerializationOptions_Builder[T interface {
	SerializationOptions_Builder
	comparable
}](v T) SerializationOptions_Builder {
	var z T
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.jscomp.serialization.SerializationOptions.Builder
type SerializationOptions_BuilderBase struct {
	idHash_ atomic.Int32
}

func (this_ *SerializationOptions_BuilderBase) BindSelf_(s any) {
	_ = s
}
func (this_ *SerializationOptions_BuilderBase) IsJavaNil_() bool { return this_ == nil }
func (*SerializationOptions_BuilderBase) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.serialization.SerializationOptions$Builder", "Builder" }
func (this_ *SerializationOptions_BuilderBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.SerializationOptions.Builder", "com.google.javascript.jscomp.serialization.SerializationOptions$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[SerializationOptions_Builder](v_)
	})
}
func (b *SerializationOptions_BuilderBase) GetSerializationOptions_BuilderBase() *SerializationOptions_BuilderBase { return b }
func (*SerializationOptions_BuilderBase) Implements_jscomp_SerializationOptions_Builder() {}

func (this_ *SerializationOptions_BuilderBase) initFields_SerializationOptions_Builder() {
	_ = this_
}

func InitSerializationOptions_Builder(this_ *SerializationOptions_BuilderBase) {
	this_.initFields_SerializationOptions_Builder()
}

// Java method: Builder.setIncludeDebugInfo(boolean includeDebugInfo)
func (this_ *SerializationOptions_BuilderBase) SetIncludeDebugInfo(includeDebugInfo bool) SerializationOptions_Builder {
	_ = includeDebugInfo
	panic("abstract or native method invoked")
}

// Java method: Builder.setRunValidation(boolean runValidation)
func (this_ *SerializationOptions_BuilderBase) SetRunValidation(runValidation bool) SerializationOptions_Builder {
	_ = runValidation
	panic("abstract or native method invoked")
}

// Java method: Builder.setRuntimeLibraries(com.google.common.collect.ImmutableList<java.lang.String> runtimeLibraries)
func (this_ *SerializationOptions_BuilderBase) SetRuntimeLibraries(runtimeLibraries j2gort.Ext_com_google_common_collect_ImmutableList) SerializationOptions_Builder {
	_ = runtimeLibraries
	panic("abstract or native method invoked")
}

// Java method: Builder.build()
func (this_ *SerializationOptions_BuilderBase) Build() *SerializationOptions {
	panic("abstract or native method invoked")
}

