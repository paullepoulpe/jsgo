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
// Port of com/google/javascript/jscomp/RecentChange.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.RecentChange
type RecentChange struct {
	idHash_    atomic.Int32
	hasChanged bool `gson:"hasChanged"`
}

func (this_ *RecentChange) BindSelf_(s any) {
	_ = s
}
func (this_ *RecentChange) IsJavaNil_() bool { return this_ == nil }
func (*RecentChange) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RecentChange", "RecentChange"
}
func (this_ *RecentChange) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RecentChange) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RecentChange", "com.google.javascript.jscomp.RecentChange", func(v_ any) bool {
		return j2gort.InstanceOf[*RecentChange](v_)
	})
}
func (*RecentChange) Implements_jscomp_CodeChangeHandler() {}

func (this_ *RecentChange) initFields_RecentChange() {
	_ = this_

}

// Java constructor: RecentChange()
func NewRecentChange() *RecentChange {
	this_ := &RecentChange{}
	InitRecentChange(this_)
	return this_
}

func InitRecentChange(this_ *RecentChange) {
	this_.initFields_RecentChange()
}

// Java method: RecentChange.reportChange()
func (this_ *RecentChange) ReportChange() {
	this_.hasChanged = true
}

// Java method: RecentChange.hasCodeChanged()
func (this_ *RecentChange) HasCodeChanged() bool {
	return this_.hasChanged
}

// Java method: RecentChange.reset()
func (this_ *RecentChange) Reset() {
	this_.hasChanged = false
}
