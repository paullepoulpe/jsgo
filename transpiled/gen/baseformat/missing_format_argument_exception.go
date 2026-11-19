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
// Port of com/google/javascript/jscomp/base/format/MissingFormatArgumentException.java (generated mechanically by j2go).

package baseformat

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.base.format.MissingFormatArgumentException
type MissingFormatArgumentException struct {
	javalang.RuntimeException
	idHash_ atomic.Int32
}

func (this_ *MissingFormatArgumentException) BindSelf_(s any) {
	_ = s
}
func (this_ *MissingFormatArgumentException) IsJavaNil_() bool { return this_ == nil }
func (*MissingFormatArgumentException) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.base.format.MissingFormatArgumentException", "MissingFormatArgumentException" }
func (this_ *MissingFormatArgumentException) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.base.format.MissingFormatArgumentException", "com.google.javascript.jscomp.base.format.MissingFormatArgumentException", func(v_ any) bool {
		return j2gort.InstanceOf[*MissingFormatArgumentException](v_)
	})
}

func (this_ *MissingFormatArgumentException) initFields_MissingFormatArgumentException() {
	_ = this_
}

// Java constructor: MissingFormatArgumentException(java.lang.String message)
func NewMissingFormatArgumentException(message string) *MissingFormatArgumentException {
	this_ := &MissingFormatArgumentException{}
	InitMissingFormatArgumentException(this_, message)
	return this_
}

func InitMissingFormatArgumentException(this_ *MissingFormatArgumentException, message string) {
	_ = message
	j2gort.InitUserThrowable(this_, "com.google.javascript.jscomp.base.format.MissingFormatArgumentException", j2gort.BoxString(message))
	this_.initFields_MissingFormatArgumentException()
}

