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
// Port of com/google/javascript/jscomp/EmptyMessageBundle.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.EmptyMessageBundle
type EmptyMessageBundle struct {
	idHash_ atomic.Int32
}

func (this_ *EmptyMessageBundle) BindSelf_(s any) {
	_ = s
}
func (this_ *EmptyMessageBundle) IsJavaNil_() bool { return this_ == nil }
func (*EmptyMessageBundle) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.EmptyMessageBundle", "EmptyMessageBundle" }
func (this_ *EmptyMessageBundle) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*EmptyMessageBundle) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.EmptyMessageBundle", "com.google.javascript.jscomp.EmptyMessageBundle", func(v_ any) bool {
		return j2gort.InstanceOf[*EmptyMessageBundle](v_)
	})
}
func (*EmptyMessageBundle) Implements_jscomp_MessageBundle() {}

func (this_ *EmptyMessageBundle) initFields_EmptyMessageBundle() {
	_ = this_
}

// Java constructor: EmptyMessageBundle()
func NewEmptyMessageBundle() *EmptyMessageBundle {
	this_ := &EmptyMessageBundle{}
	InitEmptyMessageBundle(this_)
	return this_
}

func InitEmptyMessageBundle(this_ *EmptyMessageBundle) {
	this_.initFields_EmptyMessageBundle()
}

// Java method: EmptyMessageBundle.idGenerator()
func (this_ *EmptyMessageBundle) IdGenerator() JsMessage_IdGenerator {
	return nil
}

// Java method: EmptyMessageBundle.getMessage(java.lang.String id)
func (this_ *EmptyMessageBundle) GetMessage(id string) *JsMessage {
	_ = id
	return nil
}

// Java method: EmptyMessageBundle.getAllMessages()
func (this_ *EmptyMessageBundle) GetAllMessages() any {
	return j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
}

