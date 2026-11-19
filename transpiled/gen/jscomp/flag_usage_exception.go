/*
 * Copyright 2018 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/FlagUsageException.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.FlagUsageException
type FlagUsageException struct {
	javalang.RuntimeException
	idHash_ atomic.Int32
}

func (this_ *FlagUsageException) BindSelf_(s any) {
	_ = s
}
func (this_ *FlagUsageException) IsJavaNil_() bool { return this_ == nil }
func (*FlagUsageException) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.FlagUsageException", "FlagUsageException" }
func (this_ *FlagUsageException) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.FlagUsageException", "com.google.javascript.jscomp.FlagUsageException", func(v_ any) bool {
		return j2gort.InstanceOf[*FlagUsageException](v_)
	})
}

func StaticGet_FlagUsageException_SerialVersionUID() int64 { return int64(1) }
func (this_ *FlagUsageException) initFields_FlagUsageException() {
	_ = this_
}

// Java constructor: FlagUsageException(java.lang.String message)
func NewFlagUsageException(message string) *FlagUsageException {
	this_ := &FlagUsageException{}
	InitFlagUsageException(this_, message)
	return this_
}

func InitFlagUsageException(this_ *FlagUsageException, message string) {
	_ = message
	j2gort.InitUserThrowable(this_, "com.google.javascript.jscomp.FlagUsageException", j2gort.BoxString(message))
	this_.initFields_FlagUsageException()
}

