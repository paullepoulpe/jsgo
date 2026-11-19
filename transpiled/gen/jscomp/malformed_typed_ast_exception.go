/*
 * Copyright 2020 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/serialization/MalformedTypedAstException.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.serialization.MalformedTypedAstException
type MalformedTypedAstException struct {
	javalang.RuntimeException
	idHash_ atomic.Int32
}

func (this_ *MalformedTypedAstException) BindSelf_(s any) {
	_ = s
}
func (this_ *MalformedTypedAstException) IsJavaNil_() bool { return this_ == nil }
func (*MalformedTypedAstException) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.serialization.MalformedTypedAstException", "MalformedTypedAstException"
}
func (this_ *MalformedTypedAstException) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.MalformedTypedAstException", "com.google.javascript.jscomp.serialization.MalformedTypedAstException", func(v_ any) bool {
		return j2gort.InstanceOf[*MalformedTypedAstException](v_)
	})
}

func (this_ *MalformedTypedAstException) initFields_MalformedTypedAstException() {
	_ = this_
}

// Java constructor: MalformedTypedAstException(java.lang.Object msg)
func NewMalformedTypedAstException_1(msg any) *MalformedTypedAstException {
	this_ := &MalformedTypedAstException{}
	InitMalformedTypedAstException_1(this_, msg)
	return this_
}

func InitMalformedTypedAstException_1(this_ *MalformedTypedAstException, msg any) {
	_ = msg
	j2gort.InitUserThrowable(this_, "com.google.javascript.jscomp.serialization.MalformedTypedAstException", javalang.ToJavaString(j2gort.CheckNotNull(msg)))
	this_.initFields_MalformedTypedAstException()
}

// Java constructor: MalformedTypedAstException(java.lang.Object msg, java.lang.Throwable cause)
func NewMalformedTypedAstException_2(msg any, cause javalang.Throwable) *MalformedTypedAstException {
	this_ := &MalformedTypedAstException{}
	InitMalformedTypedAstException_2(this_, msg, cause)
	return this_
}

func InitMalformedTypedAstException_2(this_ *MalformedTypedAstException, msg any, cause javalang.Throwable) {
	_ = msg
	_ = cause
	j2gort.InitUserThrowable(this_, "com.google.javascript.jscomp.serialization.MalformedTypedAstException", javalang.ToJavaString(j2gort.CheckNotNull(msg)), cause)
	this_.initFields_MalformedTypedAstException()
}

// Java static method: MalformedTypedAstException.checkWellFormed(boolean condition, java.lang.String description)
func MalformedTypedAstException_CheckWellFormed_2(condition bool, description string) {
	_ = condition
	_ = description
	if !(condition) {
		panic(NewMalformedTypedAstException_1(j2gort.BoxString(description)))
	}
}

// Java static method: MalformedTypedAstException.checkWellFormed(boolean condition, java.lang.String description, java.lang.@org.jspecify.annotations.Nullable Object param)
func MalformedTypedAstException_CheckWellFormed_3(condition bool, description string, param any) {
	_ = condition
	_ = description
	_ = param
	if !(condition) {
		var message string = (j2gort.StrValString(description) + ": " + j2gort.StrVal(param))
		_ = message
		panic(NewMalformedTypedAstException_1(j2gort.BoxString(message)))
	}
}
