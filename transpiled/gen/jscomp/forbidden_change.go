/*
 * Copyright 2013 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/ForbiddenChange.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.ForbiddenChange
type ForbiddenChange struct {
	idHash_ atomic.Int32
}

func (this_ *ForbiddenChange) BindSelf_(s any) {
	_ = s
}
func (this_ *ForbiddenChange) IsJavaNil_() bool { return this_ == nil }
func (*ForbiddenChange) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ForbiddenChange", "ForbiddenChange"
}
func (this_ *ForbiddenChange) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ForbiddenChange) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ForbiddenChange", "com.google.javascript.jscomp.ForbiddenChange", func(v_ any) bool {
		return j2gort.InstanceOf[*ForbiddenChange](v_)
	})
}
func (*ForbiddenChange) Implements_jscomp_CodeChangeHandler() {}

func (this_ *ForbiddenChange) initFields_ForbiddenChange() {
	_ = this_
}

// Java constructor: ForbiddenChange()
func NewForbiddenChange() *ForbiddenChange {
	this_ := &ForbiddenChange{}
	InitForbiddenChange(this_)
	return this_
}

func InitForbiddenChange(this_ *ForbiddenChange) {
	this_.initFields_ForbiddenChange()
}

// Java method: ForbiddenChange.reportChange()
func (this_ *ForbiddenChange) ReportChange() {
	panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Code changes forbidden")))
}
