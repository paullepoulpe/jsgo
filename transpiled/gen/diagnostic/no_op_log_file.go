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
// Port of com/google/javascript/jscomp/diagnostic/NoOpLogFile.java (generated mechanically by j2go).

package diagnostic

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.diagnostic.NoOpLogFile
type NoOpLogFile struct {
	LogFileBase
}

func (this_ *NoOpLogFile) BindSelf_(s any) {
	_ = s
}
func (this_ *NoOpLogFile) IsJavaNil_() bool { return this_ == nil }
func (*NoOpLogFile) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.diagnostic.NoOpLogFile", "NoOpLogFile"
}
func (this_ *NoOpLogFile) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.LogFileBase.IdentityHash_()
}
func (*NoOpLogFile) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.diagnostic.NoOpLogFile", "com.google.javascript.jscomp.diagnostic.NoOpLogFile", func(v_ any) bool {
		return j2gort.InstanceOf[*NoOpLogFile](v_)
	})
}
func (*NoOpLogFile) Implements_diagnostic_LogFile()                 {}
func (*NoOpLogFile) Implements_j2gort_Ext_java_lang_AutoCloseable() {}

func (this_ *NoOpLogFile) initFields_NoOpLogFile() {
	_ = this_
}

// Java constructor: NoOpLogFile()
func NewNoOpLogFile() *NoOpLogFile {
	this_ := &NoOpLogFile{}
	InitNoOpLogFile(this_)
	return this_
}

func InitNoOpLogFile(this_ *NoOpLogFile) {
	InitLogFile(&this_.LogFileBase)
	this_.initFields_NoOpLogFile()
}

// Java method: NoOpLogFile.log(java.lang.Object value)
func (this_ *NoOpLogFile) Log_1_any(value any) LogFile {
	_ = value
	return this_
}

// Java method: NoOpLogFile.log(java.lang.String value)
func (this_ *NoOpLogFile) Log_1_string(value string) LogFile {
	_ = value
	return this_
}

// Java method: NoOpLogFile.log(java.util.function.Supplier<java.lang.String> value)
func (this_ *NoOpLogFile) Log_1_Ext_java_util_function_Supplier(value j2gort.Ext_java_util_function_Supplier) LogFile {
	_ = value
	return this_
}

// Java method: NoOpLogFile.log(java.lang.String template, java.lang.Object[] values)
func (this_ *NoOpLogFile) Log_2(template string, values []any) LogFile {
	_ = template
	_ = values
	return this_
}

// Java method: NoOpLogFile.logJson(java.lang.Object value)
func (this_ *NoOpLogFile) LogJson_1_any(value any) LogFile {
	_ = value
	return this_
}

// Java method: NoOpLogFile.logJson(java.util.function.Supplier<java.lang.Object> value)
func (this_ *NoOpLogFile) LogJson_1_Ext_java_util_function_Supplier(value j2gort.Ext_java_util_function_Supplier) LogFile {
	_ = value
	return this_
}

// Java method: NoOpLogFile.logJson(com.google.javascript.jscomp.diagnostic.LogFile.StreamedJsonProducer producer)
func (this_ *NoOpLogFile) LogJson_1_LogFile_StreamedJsonProducer(producer LogFile_StreamedJsonProducer) LogFile {
	_ = producer
	return this_
}

// Java method: NoOpLogFile.close()
func (this_ *NoOpLogFile) Close_0() {
}

// Java method: NoOpLogFile.isLogging()
func (this_ *NoOpLogFile) IsLogging() bool {
	return false
}
