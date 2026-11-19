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
// Port of com/google/javascript/jscomp/base/format/IllegalFormatWidthException.java (generated mechanically by j2go).

package baseformat

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.base.format.IllegalFormatWidthException
type IllegalFormatWidthException struct {
	javalang.RuntimeException
	idHash_ atomic.Int32
}

func (this_ *IllegalFormatWidthException) BindSelf_(s any) {
	_ = s
}
func (this_ *IllegalFormatWidthException) IsJavaNil_() bool { return this_ == nil }
func (*IllegalFormatWidthException) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.base.format.IllegalFormatWidthException", "IllegalFormatWidthException"
}
func (this_ *IllegalFormatWidthException) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.base.format.IllegalFormatWidthException", "com.google.javascript.jscomp.base.format.IllegalFormatWidthException", func(v_ any) bool {
		return j2gort.InstanceOf[*IllegalFormatWidthException](v_)
	})
}

func (this_ *IllegalFormatWidthException) initFields_IllegalFormatWidthException() {
	_ = this_
}

// Java constructor: IllegalFormatWidthException(int width)
func NewIllegalFormatWidthException(width int32) *IllegalFormatWidthException {
	this_ := &IllegalFormatWidthException{}
	InitIllegalFormatWidthException(this_, width)
	return this_
}

func InitIllegalFormatWidthException(this_ *IllegalFormatWidthException, width int32) {
	_ = width
	j2gort.InitUserThrowable(this_, "com.google.javascript.jscomp.base.format.IllegalFormatWidthException", "Width:"+j2gort.StrValInt32(width))
	this_.initFields_IllegalFormatWidthException()
}
