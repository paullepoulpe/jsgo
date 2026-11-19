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
// Port of com/google/javascript/jscomp/base/format/MissingFormatWidthException.java (generated mechanically by j2go).

package baseformat

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.base.format.MissingFormatWidthException
type MissingFormatWidthException struct {
	javalang.RuntimeException
	idHash_ atomic.Int32
}

func (this_ *MissingFormatWidthException) BindSelf_(s any) {
	_ = s
}
func (this_ *MissingFormatWidthException) IsJavaNil_() bool { return this_ == nil }
func (*MissingFormatWidthException) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.base.format.MissingFormatWidthException", "MissingFormatWidthException" }
func (this_ *MissingFormatWidthException) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.base.format.MissingFormatWidthException", "com.google.javascript.jscomp.base.format.MissingFormatWidthException", func(v_ any) bool {
		return j2gort.InstanceOf[*MissingFormatWidthException](v_)
	})
}

func (this_ *MissingFormatWidthException) initFields_MissingFormatWidthException() {
	_ = this_
}

// Java constructor: MissingFormatWidthException(java.lang.String msg)
func NewMissingFormatWidthException(msg string) *MissingFormatWidthException {
	this_ := &MissingFormatWidthException{}
	InitMissingFormatWidthException(this_, msg)
	return this_
}

func InitMissingFormatWidthException(this_ *MissingFormatWidthException, msg string) {
	_ = msg
	j2gort.InitUserThrowable(this_, "com.google.javascript.jscomp.base.format.MissingFormatWidthException", j2gort.BoxString(msg))
	this_.initFields_MissingFormatWidthException()
}

