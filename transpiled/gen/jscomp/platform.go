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
// Port of com/google/javascript/jscomp/Platform.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.Platform
type Platform struct {
	idHash_ atomic.Int32
}

func (this_ *Platform) BindSelf_(s any) {
	_ = s
}
func (this_ *Platform) IsJavaNil_() bool { return this_ == nil }
func (*Platform) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.Platform", "Platform" }
func (this_ *Platform) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Platform) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Platform", "com.google.javascript.jscomp.Platform", func(v_ any) bool {
		return j2gort.InstanceOf[*Platform](v_)
	})
}

func (this_ *Platform) initFields_Platform() {
	_ = this_
}

// Java constructor: Platform()
func NewPlatform() *Platform {
	this_ := &Platform{}
	InitPlatform(this_)
	return this_
}

func InitPlatform(this_ *Platform) {
	this_.initFields_Platform()
}

// Java static method: Platform.getFileSeperator()
func Platform_GetFileSeperator() string {
	return j2gort.Ext_java_io_File_F_Separator
}

// Java static method: Platform.isThreadInterrupted()
func Platform_IsThreadInterrupted() bool {
	return j2gort.Ext_java_lang_Thread_Interrupted()
}

// Java static method: Platform.freeMemory()
func Platform_FreeMemory() int64 {
	return j2gort.Ext_java_lang_Runtime_GetRuntime().FreeMemory()
}

// Java static method: Platform.totalMemory()
func Platform_TotalMemory() int64 {
	return j2gort.Ext_java_lang_Runtime_GetRuntime().TotalMemory()
}

// Java static method: Platform.formatMessage(java.lang.String message, java.lang.String[] arguments)
func Platform_FormatMessage(message string, arguments []string) string {
	_ = message
	_ = arguments
	return j2gort.Ext_java_text_MessageFormat_Format(message, func() []any { arr_1 := arguments; if arr_1 == nil { return nil }; out_2 := make([]any, len(arr_1)); for i_, el_3 := range arr_1 { out_2[i_] = j2gort.BoxString(el_3) }; return out_2 }())
}

