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
// Port of com/google/javascript/jscomp/base/JSCompObjects.java (generated mechanically by j2go).

package base

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.base.JSCompObjects
type JSCompObjects struct {
	idHash_ atomic.Int32
}

func (this_ *JSCompObjects) BindSelf_(s any) {
	_ = s
}
func (this_ *JSCompObjects) IsJavaNil_() bool { return this_ == nil }
func (*JSCompObjects) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.base.JSCompObjects", "JSCompObjects" }
func (this_ *JSCompObjects) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JSCompObjects) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.base.JSCompObjects", "com.google.javascript.jscomp.base.JSCompObjects", func(v_ any) bool {
		return j2gort.InstanceOf[*JSCompObjects](v_)
	})
}

func (this_ *JSCompObjects) initFields_JSCompObjects() {
	_ = this_
}

// Java constructor: JSCompObjects()
func NewJSCompObjects() *JSCompObjects {
	this_ := &JSCompObjects{}
	InitJSCompObjects(this_)
	return this_
}

func InitJSCompObjects(this_ *JSCompObjects) {
	this_.initFields_JSCompObjects()
	panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError")))
}

// Java static method: JSCompObjects.identical(java.lang.Object a, java.lang.Object b)
func JSCompObjects_Identical(a any, b any) bool {
	_ = a
	_ = b
	return j2gort.RefEquals(a, b)
}

