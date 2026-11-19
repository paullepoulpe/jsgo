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
// Port of com/google/javascript/jscomp/base/format/IllegalFormatConversionException.java (generated mechanically by j2go).

package baseformat

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.base.format.IllegalFormatConversionException
type IllegalFormatConversionException struct {
	javalang.RuntimeException
	idHash_ atomic.Int32
}

func (this_ *IllegalFormatConversionException) BindSelf_(s any) {
	_ = s
}
func (this_ *IllegalFormatConversionException) IsJavaNil_() bool { return this_ == nil }
func (*IllegalFormatConversionException) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.base.format.IllegalFormatConversionException", "IllegalFormatConversionException" }
func (this_ *IllegalFormatConversionException) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.base.format.IllegalFormatConversionException", "com.google.javascript.jscomp.base.format.IllegalFormatConversionException", func(v_ any) bool {
		return j2gort.InstanceOf[*IllegalFormatConversionException](v_)
	})
}

func (this_ *IllegalFormatConversionException) initFields_IllegalFormatConversionException() {
	_ = this_
}

// Java constructor: IllegalFormatConversionException(char conversionType, java.lang.Class<? extends java.lang.Object> class1)
func NewIllegalFormatConversionException(conversionType uint16, class1 *javalang.Class) *IllegalFormatConversionException {
	this_ := &IllegalFormatConversionException{}
	InitIllegalFormatConversionException(this_, conversionType, class1)
	return this_
}

func InitIllegalFormatConversionException(this_ *IllegalFormatConversionException, conversionType uint16, class1 *javalang.Class) {
	_ = conversionType
	_ = class1
	j2gort.InitUserThrowable(this_, "com.google.javascript.jscomp.base.format.IllegalFormatConversionException", j2gort.BoxString(("ConversionType = " + j2gort.StrValChar(conversionType) + " class = " + j2gort.StrVal(class1))))
	this_.initFields_IllegalFormatConversionException()
}

