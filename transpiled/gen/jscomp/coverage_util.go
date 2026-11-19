/*
 * Copyright 2009 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/instrumentation/CoverageUtil.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.instrumentation.CoverageUtil
type CoverageUtil struct {
	idHash_ atomic.Int32
}

func (this_ *CoverageUtil) BindSelf_(s any) {
	_ = s
}
func (this_ *CoverageUtil) IsJavaNil_() bool { return this_ == nil }
func (*CoverageUtil) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.instrumentation.CoverageUtil", "CoverageUtil" }
func (this_ *CoverageUtil) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CoverageUtil) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.instrumentation.CoverageUtil", "com.google.javascript.jscomp.instrumentation.CoverageUtil", func(v_ any) bool {
		return j2gort.InstanceOf[*CoverageUtil](v_)
	})
}

func (this_ *CoverageUtil) initFields_CoverageUtil() {
	_ = this_
}

// Java constructor: CoverageUtil()
func NewCoverageUtil() *CoverageUtil {
	this_ := &CoverageUtil{}
	InitCoverageUtil(this_)
	return this_
}

func InitCoverageUtil(this_ *CoverageUtil) {
	this_.initFields_CoverageUtil()
}

// Java static method: CoverageUtil.createIdentifierFromText(java.lang.String inputText)
func CoverageUtil_CreateIdentifierFromText(inputText string) string {
	_ = inputText
	return j2gort.JStrReplaceAll(inputText, "[^\\p{Alnum}]", "_")
}

