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
// Port of com/google/javascript/jscomp/base/JSCompStrings.java (generated mechanically by j2go).

package base

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.base.JSCompStrings
type JSCompStrings struct {
	idHash_ atomic.Int32
}

func (this_ *JSCompStrings) BindSelf_(s any) {
	_ = s
}
func (this_ *JSCompStrings) IsJavaNil_() bool { return this_ == nil }
func (*JSCompStrings) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.base.JSCompStrings", "JSCompStrings"
}
func (this_ *JSCompStrings) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JSCompStrings) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.base.JSCompStrings", "com.google.javascript.jscomp.base.JSCompStrings", func(v_ any) bool {
		return j2gort.InstanceOf[*JSCompStrings](v_)
	})
}

func (this_ *JSCompStrings) initFields_JSCompStrings() {
	_ = this_
}

// Java constructor: JSCompStrings()
func NewJSCompStrings() *JSCompStrings {
	this_ := &JSCompStrings{}
	InitJSCompStrings(this_)
	return this_
}

func InitJSCompStrings(this_ *JSCompStrings) {
	this_.initFields_JSCompStrings()
	panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError")))
}

// Java static method: JSCompStrings.lines(java.lang.String[] lines)
func JSCompStrings_Lines_1_slice_string(lines []string) string {
	_ = lines
	return j2gort.Ext_java_lang_String_Join_2_any_slice_any("\n", func() []any {
		arr_1 := lines
		if arr_1 == nil {
			return nil
		}
		out_2 := make([]any, len(arr_1))
		for i_, el_3 := range arr_1 {
			out_2[i_] = j2gort.BoxString(el_3)
		}
		return out_2
	}())
}

// Java static method: JSCompStrings.lines(java.lang.Iterable<java.lang.String> lines)
func JSCompStrings_Lines_1_Ext_java_lang_Iterable(lines j2gort.Ext_java_lang_Iterable) string {
	_ = lines
	return j2gort.Ext_java_lang_String_Join_2_any_Ext_java_lang_Iterable("\n", lines)
}
