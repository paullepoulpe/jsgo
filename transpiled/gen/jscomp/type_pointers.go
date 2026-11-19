/*
 * Copyright 2021 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/serialization/TypePointers.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.serialization.TypePointers
type TypePointers struct {
	idHash_ atomic.Int32
}

func (this_ *TypePointers) BindSelf_(s any) {
	_ = s
}
func (this_ *TypePointers) IsJavaNil_() bool { return this_ == nil }
func (*TypePointers) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.serialization.TypePointers", "TypePointers" }
func (this_ *TypePointers) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TypePointers) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.TypePointers", "com.google.javascript.jscomp.serialization.TypePointers", func(v_ any) bool {
		return j2gort.InstanceOf[*TypePointers](v_)
	})
}

var static_TypePointers_OFFSET_TO_AXIOMATIC_COLOR j2gort.Ext_com_google_common_collect_ImmutableList
func StaticGet_TypePointers_OFFSET_TO_AXIOMATIC_COLOR() j2gort.Ext_com_google_common_collect_ImmutableList {
	clinit_TypePointers()
	return static_TypePointers_OFFSET_TO_AXIOMATIC_COLOR
}
func StaticSet_TypePointers_OFFSET_TO_AXIOMATIC_COLOR(v j2gort.Ext_com_google_common_collect_ImmutableList) {
	clinit_TypePointers()
	static_TypePointers_OFFSET_TO_AXIOMATIC_COLOR = v
}
var static_TypePointers_AXIOMATIC_COLOR_COUNT int32
func StaticGet_TypePointers_AXIOMATIC_COLOR_COUNT() int32 {
	clinit_TypePointers()
	return static_TypePointers_AXIOMATIC_COLOR_COUNT
}
func StaticSet_TypePointers_AXIOMATIC_COLOR_COUNT(v int32) {
	clinit_TypePointers()
	static_TypePointers_AXIOMATIC_COLOR_COUNT = v
}
var clinitDone_TypePointers uint32
var clinitState_TypePointers j2gort.ClinitState
func clinit_TypePointers() {
	if atomic.LoadUint32(&clinitDone_TypePointers) != 0 {
		return
	}
	clinitSlow_TypePointers()
}

func clinitSlow_TypePointers() {
	if !j2gort.EnterClinit(&clinitState_TypePointers, &clinitDone_TypePointers) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_TypePointers, &clinitDone_TypePointers)
	static_TypePointers_OFFSET_TO_AXIOMATIC_COLOR = j2gort.Ext_com_google_common_collect_ImmutableList_Of_10(j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_UNKNOWN()), j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_BOOLEAN()), j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_STRING()), j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_NUMBER()), j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_NULL_OR_VOID()), j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_SYMBOL()), j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_BIGINT()), j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_TOP_OBJECT()), j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_TOP_FUNCTION()), j2gort.PtrToAny(pkg_colors.StaticGet_StandardColors_GBIGINT()))
	static_TypePointers_AXIOMATIC_COLOR_COUNT = StaticGet_TypePointers_OFFSET_TO_AXIOMATIC_COLOR().Size()
}

func (this_ *TypePointers) initFields_TypePointers() {
	_ = this_
}

// Java constructor: TypePointers()
func NewTypePointers() *TypePointers {
	clinit_TypePointers()
	this_ := &TypePointers{}
	InitTypePointers(this_)
	return this_
}

func InitTypePointers(this_ *TypePointers) {
	this_.initFields_TypePointers()
	panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError")))
}

// Java static method: TypePointers.trimOffset(int x)
func TypePointers_TrimOffset(x int32) int32 {
	clinit_TypePointers()
	_ = x
	if !((x >= StaticGet_TypePointers_AXIOMATIC_COLOR_COUNT())) { j2gort.CheckStateMsg(false, j2gort.ValueOfInt(x)) }
	return (x - StaticGet_TypePointers_AXIOMATIC_COLOR_COUNT())
}

// Java static method: TypePointers.untrimOffset(int x)
func TypePointers_UntrimOffset(x int32) int32 {
	clinit_TypePointers()
	_ = x
	if !((x >= int32(0))) { j2gort.CheckStateMsg(false, j2gort.ValueOfInt(x)) }
	return (x + StaticGet_TypePointers_AXIOMATIC_COLOR_COUNT())
}

// Java static method: TypePointers.isAxiomatic(int x)
func TypePointers_IsAxiomatic(x int32) bool {
	clinit_TypePointers()
	_ = x
	return (x < StaticGet_TypePointers_AXIOMATIC_COLOR_COUNT())
}

