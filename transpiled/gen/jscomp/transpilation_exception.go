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
// Port of com/google/javascript/jscomp/bundle/TranspilationException.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.bundle.TranspilationException
type TranspilationException struct {
	javalang.RuntimeException
	errors   j2gort.Ext_com_google_common_collect_ImmutableList `gson:"errors"`
	warnings j2gort.Ext_com_google_common_collect_ImmutableList `gson:"warnings"`
	idHash_  atomic.Int32
}

func (this_ *TranspilationException) BindSelf_(s any) {
	_ = s
}
func (this_ *TranspilationException) IsJavaNil_() bool { return this_ == nil }
func (*TranspilationException) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.bundle.TranspilationException", "TranspilationException"
}
func (this_ *TranspilationException) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.bundle.TranspilationException", "com.google.javascript.jscomp.bundle.TranspilationException", func(v_ any) bool {
		return j2gort.InstanceOf[*TranspilationException](v_)
	})
}

func (this_ *TranspilationException) initFields_TranspilationException() {
	_ = this_
}

// Java constructor: TranspilationException(java.lang.Exception cause)
func NewTranspilationException_1(cause javalang.IException) *TranspilationException {
	this_ := &TranspilationException{}
	InitTranspilationException_1(this_, cause)
	return this_
}

func InitTranspilationException_1(this_ *TranspilationException, cause javalang.IException) {
	_ = cause
	InitTranspilationException_2(this_, TranspilationException_TryCastToTranspilationException(cause.GetCause()), cause)
}

// Java constructor: TranspilationException(com.google.javascript.jscomp.SourceExcerptProvider source, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.JSError> errors, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.JSError> warnings)
func NewTranspilationException_3(source SourceExcerptProvider, errors j2gort.Ext_com_google_common_collect_ImmutableList, warnings j2gort.Ext_com_google_common_collect_ImmutableList) *TranspilationException {
	this_ := &TranspilationException{}
	InitTranspilationException_3(this_, source, errors, warnings)
	return this_
}

func InitTranspilationException_3(this_ *TranspilationException, source SourceExcerptProvider, errors j2gort.Ext_com_google_common_collect_ImmutableList, warnings j2gort.Ext_com_google_common_collect_ImmutableList) {
	_ = source
	_ = errors
	_ = warnings
	InitTranspilationException_4(this_, errors, warnings, TranspilationException_Format(source, errors, warnings), nil)
}

// Java constructor: TranspilationException(com.google.javascript.jscomp.bundle.TranspilationException root, java.lang.Exception cause)
func NewTranspilationException_2(root *TranspilationException, cause javalang.IException) *TranspilationException {
	this_ := &TranspilationException{}
	InitTranspilationException_2(this_, root, cause)
	return this_
}

func InitTranspilationException_2(this_ *TranspilationException, root *TranspilationException, cause javalang.IException) {
	_ = root
	_ = cause
	InitTranspilationException_4(this_, root.errors, root.warnings, j2gort.ThrowableGetMessage(root), cause)
}

// Java constructor: TranspilationException(com.google.common.collect.ImmutableList<com.google.javascript.jscomp.JSError> errors, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.JSError> warnings, java.lang.String formatted, java.lang.@org.jspecify.annotations.Nullable Exception cause)
func NewTranspilationException_4(errors j2gort.Ext_com_google_common_collect_ImmutableList, warnings j2gort.Ext_com_google_common_collect_ImmutableList, formatted string, cause javalang.IException) *TranspilationException {
	this_ := &TranspilationException{}
	InitTranspilationException_4(this_, errors, warnings, formatted, cause)
	return this_
}

func InitTranspilationException_4(this_ *TranspilationException, errors j2gort.Ext_com_google_common_collect_ImmutableList, warnings j2gort.Ext_com_google_common_collect_ImmutableList, formatted string, cause javalang.IException) {
	_ = errors
	_ = warnings
	_ = formatted
	_ = cause
	j2gort.InitUserThrowable(this_, "com.google.javascript.jscomp.bundle.TranspilationException", j2gort.BoxString(formatted), cause)
	this_.initFields_TranspilationException()
	this_.errors = errors
	this_.warnings = warnings
}

// Java method: TranspilationException.errors()
func (this_ *TranspilationException) Errors() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.errors
}

// Java method: TranspilationException.warnings()
func (this_ *TranspilationException) Warnings() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.warnings
}

// Java static method: TranspilationException.format(com.google.javascript.jscomp.SourceExcerptProvider source, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.JSError> errors, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.JSError> warnings)
func TranspilationException_Format(source SourceExcerptProvider, errors j2gort.Ext_com_google_common_collect_ImmutableList, warnings j2gort.Ext_com_google_common_collect_ImmutableList) string {
	_ = source
	_ = errors
	_ = warnings
	var sb j2gort.LocalStringBuilder
	sb.AppendStr("Transpilation failed:\n")

	var formatter MessageFormatter
	_ = formatter
	if source != nil && !source.IsJavaNil_() {
		formatter = StaticGet_ErrorFormat_SINGLELINE().ToFormatter(source, false)
	} else {
		formatter = StaticGet_ErrorFormat_SOURCELESS().ToFormatter(source, false)
	}
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(errors))
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			error_ := j2gort.CastPtr[JSError](rawElem_2)
			_ = error_
			sb.AppendStr("\n").AppendStr(formatter.FormatError(error_))
		}
	}
	{
		cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(warnings))
		for cur_3.HasNext() {
			rawElem_4 := cur_3.Next()
			warning := j2gort.CastPtr[JSError](rawElem_4)
			_ = warning
			sb.AppendStr("\n").AppendStr(formatter.FormatError(warning))
		}
	}
	return sb.String()
}

// Java static method: TranspilationException.tryCastToTranspilationException(java.lang.Throwable t)
func TranspilationException_TryCastToTranspilationException(t javalang.Throwable) *TranspilationException {
	_ = t
	var transpilationException *TranspilationException
	_ = transpilationException
	var ok_5 bool
	transpilationException, ok_5 = j2gort.TryCast[*TranspilationException](t)
	if ok_5 {
		return transpilationException
	}
	return nil
}
