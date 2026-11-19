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
// Port of com/google/javascript/jscomp/diagnostic/LogFile.java (generated mechanically by j2go).

package diagnostic

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java polymorphic class interface: com.google.javascript.jscomp.diagnostic.LogFile
type LogFile interface {
	Implements_diagnostic_LogFile()
	IsJavaNil_() bool
	GetLogFileBase() *LogFileBase
	j2gort.Ext_java_lang_AutoCloseable
	Log_1_any(any) LogFile
	Log_1_string(string) LogFile
	Log_1_Ext_java_util_function_Supplier(j2gort.Ext_java_util_function_Supplier) LogFile
	Log_2(string, []any) LogFile
	LogJson_1_any(any) LogFile
	LogJson_1_Ext_java_util_function_Supplier(j2gort.Ext_java_util_function_Supplier) LogFile
	LogJson_1_LogFile_StreamedJsonProducer(LogFile_StreamedJsonProducer) LogFile
	IsLogging() bool
}

func ToIface_LogFile[T interface {
	LogFile
	comparable
}](v T) LogFile {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_LogFile[T LogFile](v T) LogFile {
	var vi LogFile = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.diagnostic.LogFile
type LogFileBase struct {
	idHash_ atomic.Int32
}

func (this_ *LogFileBase) BindSelf_(s any) {
	_ = s
}
func (this_ *LogFileBase) IsJavaNil_() bool { return this_ == nil }
func (*LogFileBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.diagnostic.LogFile", "LogFile"
}
func (this_ *LogFileBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.diagnostic.LogFile", "com.google.javascript.jscomp.diagnostic.LogFile", func(v_ any) bool {
		return j2gort.InstanceOf[LogFile](v_)
	})
}
func (b *LogFileBase) GetLogFileBase() *LogFileBase                 { return b }
func (*LogFileBase) Implements_diagnostic_LogFile()                 {}
func (*LogFileBase) Implements_j2gort_Ext_java_lang_AutoCloseable() {}

func (this_ *LogFileBase) initFields_LogFile() {
	_ = this_
}

func InitLogFile(this_ *LogFileBase) {
	this_.initFields_LogFile()
}

// Java static method: LogFile.createOrReopen(java.nio.file.Path file)
func LogFile_CreateOrReopen(file j2gort.Ext_java_nio_file_Path) LogFile {
	_ = file
	return WritingLogFile_Create(file)
}

// Java static method: LogFile.createNoOp()
func LogFile_CreateNoOp() LogFile {
	return NewNoOpLogFile()
}

// Java method: LogFile.log(java.lang.Object value)
func (this_ *LogFileBase) Log_1_any(value any) LogFile {
	_ = value
	panic("abstract or native method invoked")
}

// Java method: LogFile.log(java.lang.String value)
func (this_ *LogFileBase) Log_1_string(value string) LogFile {
	_ = value
	panic("abstract or native method invoked")
}

// Java method: LogFile.log(java.util.function.Supplier<java.lang.String> value)
func (this_ *LogFileBase) Log_1_Ext_java_util_function_Supplier(value j2gort.Ext_java_util_function_Supplier) LogFile {
	_ = value
	panic("abstract or native method invoked")
}

// Java method: LogFile.log(java.lang.String template, java.lang.Object[] values)
func (this_ *LogFileBase) Log_2(template string, values []any) LogFile {
	_ = template
	_ = values
	panic("abstract or native method invoked")
}

// Java method: LogFile.logJson(java.lang.Object value)
func (this_ *LogFileBase) LogJson_1_any(value any) LogFile {
	_ = value
	panic("abstract or native method invoked")
}

// Java method: LogFile.logJson(java.util.function.Supplier<java.lang.Object> value)
func (this_ *LogFileBase) LogJson_1_Ext_java_util_function_Supplier(value j2gort.Ext_java_util_function_Supplier) LogFile {
	_ = value
	panic("abstract or native method invoked")
}

// Java method: LogFile.logJson(com.google.javascript.jscomp.diagnostic.LogFile.StreamedJsonProducer producer)
func (this_ *LogFileBase) LogJson_1_LogFile_StreamedJsonProducer(producer LogFile_StreamedJsonProducer) LogFile {
	_ = producer
	panic("abstract or native method invoked")
}

// Java method: LogFile.isLogging()
func (this_ *LogFileBase) IsLogging() bool {
	panic("abstract or native method invoked")
}

// Java method: LogFile.close()
func (this_ *LogFileBase) Close_0() {
	panic("abstract or native method invoked")
}

// Java interface: com.google.javascript.jscomp.diagnostic.LogFile.StreamedJsonProducer
type LogFile_StreamedJsonProducer interface {
	Implements_diagnostic_LogFile_StreamedJsonProducer()
	IsJavaNil_() bool
	WriteJson(*j2gort.Ext_com_google_gson_stream_JsonWriter)
}

func ToIface_LogFile_StreamedJsonProducer[T interface {
	LogFile_StreamedJsonProducer
	comparable
}](v T) LogFile_StreamedJsonProducer {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_LogFile_StreamedJsonProducer[T LogFile_StreamedJsonProducer](v T) LogFile_StreamedJsonProducer {
	var vi LogFile_StreamedJsonProducer = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.diagnostic.LogFile.StreamedJsonProducer", "com.google.javascript.jscomp.diagnostic.LogFile$StreamedJsonProducer", func(v_ any) bool {
		return j2gort.InstanceOf[LogFile_StreamedJsonProducer](v_)
	})
}
