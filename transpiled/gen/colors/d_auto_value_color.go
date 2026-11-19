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
// Port of com/google/javascript/jscomp/colors/$AutoValue_Color.java (generated mechanically by j2go).

package colors

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.colors.$AutoValue_Color
type D_AutoValue_Color struct {
	idHash_ atomic.Int32
}

func (this_ *D_AutoValue_Color) BindSelf_(s any) {
	_ = s
}
func (this_ *D_AutoValue_Color) IsJavaNil_() bool { return this_ == nil }
func (*D_AutoValue_Color) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.colors.$AutoValue_Color", "$AutoValue_Color" }
func (this_ *D_AutoValue_Color) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*D_AutoValue_Color) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.colors.$AutoValue_Color", "com.google.javascript.jscomp.colors.$AutoValue_Color", func(v_ any) bool {
		return j2gort.InstanceOf[*D_AutoValue_Color](v_)
	})
}

func (this_ *D_AutoValue_Color) initFields_D_AutoValue_Color() {
	_ = this_
}

// Java constructor: $AutoValue_Color()
func NewD_AutoValue_Color() *D_AutoValue_Color {
	this_ := &D_AutoValue_Color{}
	InitD_AutoValue_Color(this_)
	return this_
}

func InitD_AutoValue_Color(this_ *D_AutoValue_Color) {
	this_.initFields_D_AutoValue_Color()
}

