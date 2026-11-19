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
// Port of com/google/javascript/jscomp/diagnostic/LogsGson.java (generated mechanically by j2go).

package diagnostic

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.diagnostic.LogsGson
type LogsGson struct {
	idHash_ atomic.Int32
}

func (this_ *LogsGson) BindSelf_(s any) {
	_ = s
}
func (this_ *LogsGson) IsJavaNil_() bool { return this_ == nil }
func (*LogsGson) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.diagnostic.LogsGson", "LogsGson"
}
func (this_ *LogsGson) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*LogsGson) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.diagnostic.LogsGson", "com.google.javascript.jscomp.diagnostic.LogsGson", func(v_ any) bool {
		return j2gort.InstanceOf[*LogsGson](v_)
	})
}

var static_LogsGson_GSON *j2gort.Ext_com_google_gson_Gson

func StaticGet_LogsGson_GSON() *j2gort.Ext_com_google_gson_Gson {
	clinit_LogsGson()
	return static_LogsGson_GSON
}
func StaticSet_LogsGson_GSON(v *j2gort.Ext_com_google_gson_Gson) {
	clinit_LogsGson()
	static_LogsGson_GSON = v
}

var clinitDone_LogsGson uint32
var clinitState_LogsGson j2gort.ClinitState

func clinit_LogsGson() {
	if atomic.LoadUint32(&clinitDone_LogsGson) != 0 {
		return
	}
	clinitSlow_LogsGson()
}

func clinitSlow_LogsGson() {
	if !j2gort.EnterClinit(&clinitState_LogsGson, &clinitDone_LogsGson) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_LogsGson, &clinitDone_LogsGson)
	static_LogsGson_GSON = j2gort.NewExt_com_google_gson_GsonBuilder().RegisterTypeHierarchyAdapter(j2gort.ClassLiteral("com.google.javascript.jscomp.diagnostic.LogsGson.Able", "Able"), NewLogsGson__Anon_111()).RegisterTypeHierarchyAdapter(j2gort.ClassLiteral("com.google.common.collect.Multimap", "Multimap"), NewLogsGson__Anon_112()).Create()
}

func (this_ *LogsGson) initFields_LogsGson() {
	_ = this_
}

// Java constructor: LogsGson()
func NewLogsGson() *LogsGson {
	clinit_LogsGson()
	this_ := &LogsGson{}
	InitLogsGson(this_)
	return this_
}

func InitLogsGson(this_ *LogsGson) {
	this_.initFields_LogsGson()
	panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError")))
}

// Java static method: LogsGson.toJson(java.lang.Object o)
func LogsGson_ToJson(o any) string {
	clinit_LogsGson()
	_ = o
	return StaticGet_LogsGson_GSON().ToJson_1_any(o)
}

// Java interface: com.google.javascript.jscomp.diagnostic.LogsGson.Able
type LogsGson_Able interface {
	Implements_diagnostic_LogsGson_Able()
	IsJavaNil_() bool
	ToLogsGson() any
}

func ToIface_LogsGson_Able[T interface {
	LogsGson_Able
	comparable
}](v T) LogsGson_Able {
	var z T
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.diagnostic.LogsGson.Able", "com.google.javascript.jscomp.diagnostic.LogsGson$Able", func(v_ any) bool {
		return j2gort.InstanceOf[LogsGson_Able](v_)
	})
}

// Java class:
type LogsGson__Anon_112 struct {
	j2gort.Ext_com_google_gson_TypeAdapterBase
	idHash_ atomic.Int32
}

func (this_ *LogsGson__Anon_112) BindSelf_(s any) {
	_ = s
	this_.Ext_com_google_gson_TypeAdapterBase.BindExtSelf_(s)
}
func (this_ *LogsGson__Anon_112) IsJavaNil_() bool { return this_ == nil }
func (*LogsGson__Anon_112) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.diagnostic.LogsGson$1", ""
}
func (this_ *LogsGson__Anon_112) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*LogsGson__Anon_112) Implements_j2gort_Ext_com_google_gson_TypeAdapter() {}

func (this_ *LogsGson__Anon_112) initFields_LogsGson__Anon_112() {
	_ = this_
}

func NewLogsGson__Anon_112() *LogsGson__Anon_112 {
	this_ := &LogsGson__Anon_112{}
	this_.BindSelf_(this_)
	j2gort.InitExt_com_google_gson_TypeAdapter(&this_.Ext_com_google_gson_TypeAdapterBase)
	this_.initFields_LogsGson__Anon_112()
	return this_
}

// Java method: .read(com.google.gson.stream.JsonReader in)
func (this_ *LogsGson__Anon_112) Read(in *j2gort.Ext_com_google_gson_stream_JsonReader) any {
	_ = in
	panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError")))
}

// Java method: .write(com.google.gson.stream.JsonWriter out, com.google.common.collect.Multimap<?,?> value)
func (this_ *LogsGson__Anon_112) Write(out *j2gort.Ext_com_google_gson_stream_JsonWriter, arg_1 any) {
	_ = out
	value := j2gort.AsIface[j2gort.Ext_com_google_common_collect_Multimap](arg_1)
	_ = value
	out.JsonValue(StaticGet_LogsGson_GSON().ToJson_1_any(j2gort.AsIface[j2gort.Ext_java_util_Map](value.AsMap())))
}

// Java class:
type LogsGson__Anon_111 struct {
	j2gort.Ext_com_google_gson_TypeAdapterBase
	idHash_ atomic.Int32
}

func (this_ *LogsGson__Anon_111) BindSelf_(s any) {
	_ = s
	this_.Ext_com_google_gson_TypeAdapterBase.BindExtSelf_(s)
}
func (this_ *LogsGson__Anon_111) IsJavaNil_() bool { return this_ == nil }
func (*LogsGson__Anon_111) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.diagnostic.LogsGson$2", ""
}
func (this_ *LogsGson__Anon_111) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*LogsGson__Anon_111) Implements_j2gort_Ext_com_google_gson_TypeAdapter() {}

func (this_ *LogsGson__Anon_111) initFields_LogsGson__Anon_111() {
	_ = this_
}

func NewLogsGson__Anon_111() *LogsGson__Anon_111 {
	this_ := &LogsGson__Anon_111{}
	this_.BindSelf_(this_)
	j2gort.InitExt_com_google_gson_TypeAdapter(&this_.Ext_com_google_gson_TypeAdapterBase)
	this_.initFields_LogsGson__Anon_111()
	return this_
}

// Java method: .read(com.google.gson.stream.JsonReader in)
func (this_ *LogsGson__Anon_111) Read(in *j2gort.Ext_com_google_gson_stream_JsonReader) any {
	_ = in
	panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError")))
}

// Java method: .write(com.google.gson.stream.JsonWriter out, com.google.javascript.jscomp.diagnostic.LogsGson.Able value)
func (this_ *LogsGson__Anon_111) Write(out *j2gort.Ext_com_google_gson_stream_JsonWriter, arg_2 any) {
	_ = out
	value := j2gort.AsIface[LogsGson_Able](arg_2)
	_ = value
	out.JsonValue(StaticGet_LogsGson_GSON().ToJson_1_any(value.ToLogsGson()))
}
