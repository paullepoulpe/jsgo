/*
 * Copyright 2015 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/BlackHoleErrorManager.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.BlackHoleErrorManager
type BlackHoleErrorManager struct {
	SortingErrorManagerBase
}

func (this_ *BlackHoleErrorManager) BindSelf_(s any) {
	_ = s
	this_.SortingErrorManagerBase.BindSelf_(s)
}
func (this_ *BlackHoleErrorManager) IsJavaNil_() bool { return this_ == nil }
func (*BlackHoleErrorManager) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.BlackHoleErrorManager", "BlackHoleErrorManager" }
func (this_ *BlackHoleErrorManager) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.SortingErrorManagerBase.IdentityHash_()
}
func (*BlackHoleErrorManager) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.BlackHoleErrorManager", "com.google.javascript.jscomp.BlackHoleErrorManager", func(v_ any) bool {
		return j2gort.InstanceOf[*BlackHoleErrorManager](v_)
	})
}
func (*BlackHoleErrorManager) Implements_jscomp_SortingErrorManager() {}
func (*BlackHoleErrorManager) Implements_jscomp_ErrorManager() {}
func (*BlackHoleErrorManager) Implements_jscomp_ErrorHandler() {}

func (this_ *BlackHoleErrorManager) initFields_BlackHoleErrorManager() {
	_ = this_
}

// Java constructor: BlackHoleErrorManager()
func NewBlackHoleErrorManager() *BlackHoleErrorManager {
	this_ := &BlackHoleErrorManager{}
	this_.BindSelf_(this_)
	InitBlackHoleErrorManager(this_)
	return this_
}

func InitBlackHoleErrorManager(this_ *BlackHoleErrorManager) {
	InitSortingErrorManager(&this_.SortingErrorManagerBase, j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()))
	this_.initFields_BlackHoleErrorManager()
}

