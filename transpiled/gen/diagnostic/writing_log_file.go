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
// Port of com/google/javascript/jscomp/diagnostic/WritingLogFile.java (generated mechanically by j2go).

package diagnostic

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.diagnostic.WritingLogFile
type WritingLogFile struct {
	LogFileBase
	writer *j2gort.Ext_java_io_BufferedWriter `gson:"writer"`
}

func (this_ *WritingLogFile) BindSelf_(s any) {
	_ = s
}
func (this_ *WritingLogFile) IsJavaNil_() bool { return this_ == nil }
func (*WritingLogFile) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.diagnostic.WritingLogFile", "WritingLogFile"
}
func (this_ *WritingLogFile) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.LogFileBase.IdentityHash_()
}
func (*WritingLogFile) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.diagnostic.WritingLogFile", "com.google.javascript.jscomp.diagnostic.WritingLogFile", func(v_ any) bool {
		return j2gort.InstanceOf[*WritingLogFile](v_)
	})
}
func (*WritingLogFile) Implements_diagnostic_LogFile()                 {}
func (*WritingLogFile) Implements_j2gort_Ext_java_lang_AutoCloseable() {}

func (this_ *WritingLogFile) initFields_WritingLogFile() {
	_ = this_
}

// Java constructor: WritingLogFile(java.io.BufferedWriter writer)
func NewWritingLogFile(writer *j2gort.Ext_java_io_BufferedWriter) *WritingLogFile {
	this_ := &WritingLogFile{}
	InitWritingLogFile(this_, writer)
	return this_
}

func InitWritingLogFile(this_ *WritingLogFile, writer *j2gort.Ext_java_io_BufferedWriter) {
	_ = writer
	InitLogFile(&this_.LogFileBase)
	this_.initFields_WritingLogFile()
	this_.writer = writer
}

// Java static method: WritingLogFile.create(java.nio.file.Path file)
func WritingLogFile_Create(file j2gort.Ext_java_nio_file_Path) LogFile {
	_ = file
	var didRet_1 bool
	var retVal_2 LogFile
	func() {
		defer func() {
			if rec_3 := recover(); rec_3 != nil {
				jex_4 := javalang.RecoverException(rec_3)
				_ = jex_4
				if j2gort.InstanceOf[javaio.IIOException](jex_4) {
					e := j2gort.AsIface[javaio.IIOException](jex_4)
					_ = e
					panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", e)))
				} else {
					panic(rec_3)
				}
			}
		}()
		var dir j2gort.Ext_java_nio_file_Path = file.GetParent()
		_ = dir
		j2gort.Ext_java_nio_file_Files_CreateDirectories(dir, []j2gort.Ext_java_nio_file_attribute_FileAttribute{})
		var encoder j2gort.Ext_java_nio_charset_CharsetEncoder = j2gort.Ext_java_nio_charset_StandardCharsets_F_UTF_8.NewEncoder()
		_ = encoder
		encoder.OnMalformedInput(j2gort.Ext_java_nio_charset_CodingErrorAction_F_REPLACE).OnUnmappableCharacter(j2gort.Ext_java_nio_charset_CodingErrorAction_F_REPLACE)
		retVal_2 = NewWritingLogFile(j2gort.NewExt_java_io_BufferedWriter_1(j2gort.NewExt_java_io_OutputStreamWriter_2_Ext_java_io_OutputStream_Ext_java_nio_charset_CharsetEncoder(j2gort.Ext_java_nio_file_Files_NewOutputStream(file, []j2gort.Ext_java_nio_file_OpenOption{j2gort.AsIface[j2gort.Ext_java_nio_file_OpenOption](j2gort.Ext_java_nio_file_StandardOpenOption_F_CREATE), j2gort.AsIface[j2gort.Ext_java_nio_file_OpenOption](j2gort.Ext_java_nio_file_StandardOpenOption_F_APPEND), j2gort.AsIface[j2gort.Ext_java_nio_file_OpenOption](j2gort.Ext_java_nio_file_StandardOpenOption_F_WRITE)}), encoder)))
		didRet_1 = true
		return
	}()
	if didRet_1 {
		return retVal_2
	}
	panic("unreachable")
}

// Java method: WritingLogFile.log(java.lang.Object value)
func (this_ *WritingLogFile) Log_1_any(value any) LogFile {
	_ = value
	return this_.priv_WritingLogFile_logInternal_6220(javalang.ToJavaString(j2gort.CheckNotNull(value)))
}

// Java method: WritingLogFile.log(java.lang.String value)
func (this_ *WritingLogFile) Log_1_string(value string) LogFile {
	_ = value
	return this_.priv_WritingLogFile_logInternal_6220(value)
}

// Java method: WritingLogFile.log(java.util.function.Supplier<java.lang.String> value)
func (this_ *WritingLogFile) Log_1_Ext_java_util_function_Supplier(value j2gort.Ext_java_util_function_Supplier) LogFile {
	_ = value
	return this_.priv_WritingLogFile_logInternal_6220(j2gort.CastString(value.Get()))
}

// Java method: WritingLogFile.log(java.lang.String template, java.lang.Object[] values)
func (this_ *WritingLogFile) Log_2(template string, values []any) LogFile {
	_ = template
	_ = values
	return this_.priv_WritingLogFile_logInternal_6220(j2gort.Ext_java_lang_String_Format_2(template, values))
}

// Java method: WritingLogFile.logJson(java.lang.Object value)
func (this_ *WritingLogFile) LogJson_1_any(value any) LogFile {
	_ = value
	return this_.priv_WritingLogFile_logInternal_6220(LogsGson_ToJson(value))
}

// Java method: WritingLogFile.logJson(java.util.function.Supplier<java.lang.Object> value)
func (this_ *WritingLogFile) LogJson_1_Ext_java_util_function_Supplier(value j2gort.Ext_java_util_function_Supplier) LogFile {
	_ = value
	return this_.priv_WritingLogFile_logInternal_6220(LogsGson_ToJson(value.Get()))
}

// Java method: WritingLogFile.logJson(com.google.javascript.jscomp.diagnostic.LogFile.StreamedJsonProducer producer)
func (this_ *WritingLogFile) LogJson_1_LogFile_StreamedJsonProducer(producer LogFile_StreamedJsonProducer) LogFile {
	_ = producer
	func() {
		defer func() {
			if rec_5 := recover(); rec_5 != nil {
				jex_6 := javalang.RecoverException(rec_5)
				_ = jex_6
				if j2gort.InstanceOf[javaio.IIOException](jex_6) {
					ex := j2gort.AsIface[javaio.IIOException](jex_6)
					_ = ex
					panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError", ex)))
				} else {
					panic(rec_5)
				}
			}
		}()
		var writer *j2gort.Ext_com_google_gson_stream_JsonWriter = j2gort.NewExt_com_google_gson_stream_JsonWriter(j2gort.AsIface[j2gort.Ext_java_io_Writer](this_.writer))
		_ = writer
		defer j2gort.CloseResource(writer)
		producer.WriteJson(writer)
	}()
	return this_
}

// Java method: WritingLogFile.logInternal(java.lang.String value)
func (this_ *WritingLogFile) priv_WritingLogFile_logInternal_6220(value string) LogFile {
	_ = value
	func() {
		defer func() {
			if rec_7 := recover(); rec_7 != nil {
				jex_8 := javalang.RecoverException(rec_7)
				_ = jex_8
				if j2gort.InstanceOf[javaio.IIOException](jex_8) {
					e := j2gort.AsIface[javaio.IIOException](jex_8)
					_ = e
					panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", e)))
				} else {
					panic(rec_7)
				}
			}
		}()
		j2gort.AsIface[j2gort.Ext_java_io_Writer](this_.writer.Append_1_any(j2gort.BoxString(value))).Append_1_uint16(uint16(10))
	}()
	return this_
}

// Java method: WritingLogFile.close()
func (this_ *WritingLogFile) Close_0() {
	func() {
		defer func() {
			if rec_9 := recover(); rec_9 != nil {
				jex_10 := javalang.RecoverException(rec_9)
				_ = jex_10
				if j2gort.InstanceOf[javaio.IIOException](jex_10) {
					e := j2gort.AsIface[javaio.IIOException](jex_10)
					_ = e
					panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", e)))
				} else {
					panic(rec_9)
				}
			}
		}()
		this_.writer.Close_0()
	}()
}

// Java method: WritingLogFile.isLogging()
func (this_ *WritingLogFile) IsLogging() bool {
	return true
}
