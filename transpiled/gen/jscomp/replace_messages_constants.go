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
// Port of com/google/javascript/jscomp/ReplaceMessagesConstants.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.ReplaceMessagesConstants
type ReplaceMessagesConstants struct {
	idHash_ atomic.Int32
}

func (this_ *ReplaceMessagesConstants) BindSelf_(s any) {
	_ = s
}
func (this_ *ReplaceMessagesConstants) IsJavaNil_() bool { return this_ == nil }
func (*ReplaceMessagesConstants) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.ReplaceMessagesConstants", "ReplaceMessagesConstants" }
func (this_ *ReplaceMessagesConstants) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ReplaceMessagesConstants) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReplaceMessagesConstants", "com.google.javascript.jscomp.ReplaceMessagesConstants", func(v_ any) bool {
		return j2gort.InstanceOf[*ReplaceMessagesConstants](v_)
	})
}

func StaticGet_ReplaceMessagesConstants_DEFINE_MSG_CALLEE() string { return "__jscomp_define_msg__" }
func StaticGet_ReplaceMessagesConstants_FALLBACK_MSG_CALLEE() string { return "__jscomp_msg_fallback__" }
var static_ReplaceMessagesConstants_PROTECTED_FUNCTION_NAMES j2gort.Ext_com_google_common_collect_ImmutableSet
func StaticGet_ReplaceMessagesConstants_PROTECTED_FUNCTION_NAMES() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_ReplaceMessagesConstants()
	return static_ReplaceMessagesConstants_PROTECTED_FUNCTION_NAMES
}
func StaticSet_ReplaceMessagesConstants_PROTECTED_FUNCTION_NAMES(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_ReplaceMessagesConstants()
	static_ReplaceMessagesConstants_PROTECTED_FUNCTION_NAMES = v
}
var clinitDone_ReplaceMessagesConstants uint32
var clinitState_ReplaceMessagesConstants j2gort.ClinitState
func clinit_ReplaceMessagesConstants() {
	if atomic.LoadUint32(&clinitDone_ReplaceMessagesConstants) != 0 {
		return
	}
	clinitSlow_ReplaceMessagesConstants()
}

func clinitSlow_ReplaceMessagesConstants() {
	if !j2gort.EnterClinit(&clinitState_ReplaceMessagesConstants, &clinitDone_ReplaceMessagesConstants) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ReplaceMessagesConstants, &clinitDone_ReplaceMessagesConstants)
	static_ReplaceMessagesConstants_PROTECTED_FUNCTION_NAMES = j2gort.ImmSetOfStr2("__jscomp_define_msg__", "__jscomp_msg_fallback__")
}

func (this_ *ReplaceMessagesConstants) initFields_ReplaceMessagesConstants() {
	_ = this_
}

// Java constructor: ReplaceMessagesConstants()
func NewReplaceMessagesConstants() *ReplaceMessagesConstants {
	this_ := &ReplaceMessagesConstants{}
	InitReplaceMessagesConstants(this_)
	return this_
}

func InitReplaceMessagesConstants(this_ *ReplaceMessagesConstants) {
	this_.initFields_ReplaceMessagesConstants()
}

// Java static method: ReplaceMessagesConstants.isProtectedMessage(com.google.javascript.rhino.Node n)
func ReplaceMessagesConstants_IsProtectedMessage(n pkg_rhino.Node) bool {
	_ = n
	if (!(n.IsCall())) {
		return false
	}
	var callee pkg_rhino.Node = n.GetFirstChild()
	_ = callee
	if (!(callee.IsName())) {
		return false
	}
	return j2gort.SetContainsStr(StaticGet_ReplaceMessagesConstants_PROTECTED_FUNCTION_NAMES(), callee.GetString())
}

