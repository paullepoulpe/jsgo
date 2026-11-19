/*
 * Copyright 2011 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/parsing/parser/util/ErrorReporter.java (generated mechanically by j2go).

package parser

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java polymorphic class interface: com.google.javascript.jscomp.parsing.parser.util.ErrorReporter
type ErrorReporter interface {
	Implements_parser_ErrorReporter()
	IsJavaNil_() bool
	GetErrorReporterBase() *ErrorReporterBase
	ReportError_3(SourcePosition, string, []any)
	ReportError_2(SourcePosition, string)
	ReportWarning_3(SourcePosition, string, []any)
	ReportWarning_2(SourcePosition, string)
	HadError() bool
}
func ToIface_ErrorReporter[T interface {
	ErrorReporter
	comparable
}](v T) ErrorReporter {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_ErrorReporter[T ErrorReporter](v T) ErrorReporter {
	var vi ErrorReporter = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.parsing.parser.util.ErrorReporter
type ErrorReporterBase struct {
	self_ ErrorReporter
	idHash_ atomic.Int32
	hadError bool `gson:"hadError"`
}

func (this_ *ErrorReporterBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(ErrorReporter)
}
func (this_ *ErrorReporterBase) IsJavaNil_() bool { return this_ == nil }
func (*ErrorReporterBase) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.parsing.parser.util.ErrorReporter", "ErrorReporter" }
func (this_ *ErrorReporterBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.util.ErrorReporter", "com.google.javascript.jscomp.parsing.parser.util.ErrorReporter", func(v_ any) bool {
		return j2gort.InstanceOf[ErrorReporter](v_)
	})
}
func (b *ErrorReporterBase) GetErrorReporterBase() *ErrorReporterBase { return b }
func (*ErrorReporterBase) Implements_parser_ErrorReporter() {}

func (this_ *ErrorReporterBase) initFields_ErrorReporter() {
	_ = this_
}

func InitErrorReporter(this_ *ErrorReporterBase) {
	this_.initFields_ErrorReporter()
}

// Java method: ErrorReporter.reportError(com.google.javascript.jscomp.parsing.parser.util.SourcePosition location, java.lang.String format, java.lang.Object[] arguments)
func (this_ *ErrorReporterBase) ReportError_3(location SourcePosition, format string, arguments []any) {
	_ = location
	_ = format
	_ = arguments
	this_.hadError = true
	var message string = j2gort.Ext_java_lang_String_Format_2(format, arguments)
	_ = message
	this_.self_.ReportError_2(location, message)
}

// Java method: ErrorReporter.reportError(com.google.javascript.jscomp.parsing.parser.util.SourcePosition location, java.lang.String message)
func (this_ *ErrorReporterBase) ReportError_2(location SourcePosition, message string) {
	_ = location
	_ = message
	panic("abstract or native method invoked")
}

// Java method: ErrorReporter.reportWarning(com.google.javascript.jscomp.parsing.parser.util.SourcePosition location, java.lang.String format, java.lang.Object[] arguments)
func (this_ *ErrorReporterBase) ReportWarning_3(location SourcePosition, format string, arguments []any) {
	_ = location
	_ = format
	_ = arguments
	var message string = j2gort.Ext_java_lang_String_Format_2(format, arguments)
	_ = message
	this_.self_.ReportWarning_2(location, message)
}

// Java method: ErrorReporter.reportWarning(com.google.javascript.jscomp.parsing.parser.util.SourcePosition location, java.lang.String message)
func (this_ *ErrorReporterBase) ReportWarning_2(location SourcePosition, message string) {
	_ = location
	_ = message
	panic("abstract or native method invoked")
}

// Java method: ErrorReporter.hadError()
func (this_ *ErrorReporterBase) HadError() bool {
	return this_.hadError
}

