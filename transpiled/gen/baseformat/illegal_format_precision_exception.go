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
// Port of com/google/javascript/jscomp/base/format/IllegalFormatPrecisionException.java (generated mechanically by j2go).

package baseformat

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.base.format.IllegalFormatPrecisionException
type IllegalFormatPrecisionException struct {
	javalang.RuntimeException
	idHash_ atomic.Int32
}

func (this_ *IllegalFormatPrecisionException) BindSelf_(s any) {
	_ = s
}
func (this_ *IllegalFormatPrecisionException) IsJavaNil_() bool { return this_ == nil }
func (*IllegalFormatPrecisionException) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.base.format.IllegalFormatPrecisionException", "IllegalFormatPrecisionException"
}
func (this_ *IllegalFormatPrecisionException) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.base.format.IllegalFormatPrecisionException", "com.google.javascript.jscomp.base.format.IllegalFormatPrecisionException", func(v_ any) bool {
		return j2gort.InstanceOf[*IllegalFormatPrecisionException](v_)
	})
}

func (this_ *IllegalFormatPrecisionException) initFields_IllegalFormatPrecisionException() {
	_ = this_
}

// Java constructor: IllegalFormatPrecisionException(int precision)
func NewIllegalFormatPrecisionException(precision int32) *IllegalFormatPrecisionException {
	this_ := &IllegalFormatPrecisionException{}
	InitIllegalFormatPrecisionException(this_, precision)
	return this_
}

func InitIllegalFormatPrecisionException(this_ *IllegalFormatPrecisionException, precision int32) {
	_ = precision
	j2gort.InitUserThrowable(this_, "com.google.javascript.jscomp.base.format.IllegalFormatPrecisionException", "Precision: "+j2gort.StrValInt32(precision))
	this_.initFields_IllegalFormatPrecisionException()
}
