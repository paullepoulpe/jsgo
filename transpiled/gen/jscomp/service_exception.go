/*
 * Copyright 2008 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/deps/ServiceException.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.deps.ServiceException
type ServiceException struct {
	javalang.JavaException
	idHash_ atomic.Int32
}

func (this_ *ServiceException) BindSelf_(s any) {
	_ = s
}
func (this_ *ServiceException) IsJavaNil_() bool { return this_ == nil }
func (*ServiceException) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.deps.ServiceException", "ServiceException" }
func (this_ *ServiceException) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.ServiceException", "com.google.javascript.jscomp.deps.ServiceException", func(v_ any) bool {
		return j2gort.InstanceOf[*ServiceException](v_)
	})
}

func (this_ *ServiceException) initFields_ServiceException() {
	_ = this_
}

// Java constructor: ServiceException(java.lang.String message)
func NewServiceException(message string) *ServiceException {
	this_ := &ServiceException{}
	InitServiceException(this_, message)
	return this_
}

func InitServiceException(this_ *ServiceException, message string) {
	_ = message
	j2gort.InitUserThrowable(this_, "com.google.javascript.jscomp.deps.ServiceException", j2gort.BoxString(message))
	this_.initFields_ServiceException()
}

