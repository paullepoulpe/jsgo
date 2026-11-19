/*
 * Copyright 2026 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/serialization/FastGzipOutputStream.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.serialization.FastGzipOutputStream
type FastGzipOutputStream struct {
	j2gort.Ext_java_util_zip_GZIPOutputStreamBase
	idHash_ atomic.Int32
}

func (this_ *FastGzipOutputStream) BindSelf_(s any) {
	_ = s
	this_.Ext_java_util_zip_GZIPOutputStreamBase.BindExtSelf_(s)
}
func (this_ *FastGzipOutputStream) IsJavaNil_() bool { return this_ == nil }
func (*FastGzipOutputStream) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.serialization.FastGzipOutputStream", "FastGzipOutputStream" }
func (this_ *FastGzipOutputStream) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.FastGzipOutputStream", "com.google.javascript.jscomp.serialization.FastGzipOutputStream", func(v_ any) bool {
		return j2gort.InstanceOf[*FastGzipOutputStream](v_)
	})
}
func (*FastGzipOutputStream) Implements_j2gort_Ext_java_util_zip_GZIPOutputStream() {}
func (*FastGzipOutputStream) Implements_j2gort_Ext_java_util_zip_DeflaterOutputStream() {}
func (*FastGzipOutputStream) Implements_j2gort_Ext_java_io_FilterOutputStream() {}
func (*FastGzipOutputStream) Implements_j2gort_Ext_java_io_OutputStream() {}
func (*FastGzipOutputStream) Implements_j2gort_Ext_java_io_Closeable() {}
func (*FastGzipOutputStream) Implements_j2gort_Ext_java_io_Flushable() {}
func (*FastGzipOutputStream) Implements_j2gort_Ext_java_lang_AutoCloseable() {}

func StaticGet_FastGzipOutputStream_DEFAULT_BUFFER_SIZE() int32 { return int32(65536) }
func (this_ *FastGzipOutputStream) initFields_FastGzipOutputStream() {
	_ = this_
}

// Java constructor: FastGzipOutputStream(java.io.OutputStream out)
func NewFastGzipOutputStream_1(out j2gort.Ext_java_io_OutputStream) *FastGzipOutputStream {
	this_ := &FastGzipOutputStream{}
	this_.BindSelf_(this_)
	InitFastGzipOutputStream_1(this_, out)
	return this_
}

func InitFastGzipOutputStream_1(this_ *FastGzipOutputStream, out j2gort.Ext_java_io_OutputStream) {
	_ = out
	InitFastGzipOutputStream_2(this_, out, int32(65536))
}

// Java constructor: FastGzipOutputStream(java.io.OutputStream out, int size)
func NewFastGzipOutputStream_2(out j2gort.Ext_java_io_OutputStream, size int32) *FastGzipOutputStream {
	this_ := &FastGzipOutputStream{}
	this_.BindSelf_(this_)
	InitFastGzipOutputStream_2(this_, out, size)
	return this_
}

func InitFastGzipOutputStream_2(this_ *FastGzipOutputStream, out j2gort.Ext_java_io_OutputStream, size int32) {
	_ = out
	_ = size
	j2gort.InitExt_java_util_zip_GZIPOutputStream_2_Ext_java_io_OutputStream_int32(&this_.Ext_java_util_zip_GZIPOutputStreamBase, out, size)
	this_.initFields_FastGzipOutputStream()
	this_.GetExt_java_util_zip_DeflaterOutputStreamBase().F_Def.SetLevel(int32(1))
}

