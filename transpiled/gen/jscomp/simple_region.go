/*
 * Copyright 2007 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/SimpleRegion.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.SimpleRegion
type SimpleRegion struct {
	source string `gson:"source"`
	beginningLineNumber int32 `gson:"beginningLineNumber"`
	endingLineNumber int32 `gson:"endingLineNumber"`
	idHash_ atomic.Int32
}

func (this_ *SimpleRegion) BindSelf_(s any) {
	_ = s
}
func (this_ *SimpleRegion) IsJavaNil_() bool { return this_ == nil }
func (*SimpleRegion) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.SimpleRegion", "SimpleRegion" }
func (this_ *SimpleRegion) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SimpleRegion) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SimpleRegion", "com.google.javascript.jscomp.SimpleRegion", func(v_ any) bool {
		return j2gort.InstanceOf[*SimpleRegion](v_)
	})
}
func (*SimpleRegion) Implements_jscomp_Region() {}

func (this_ *SimpleRegion) initFields_SimpleRegion() {
	_ = this_
	this_.source = j2gort.NullString
}

// Java constructor: SimpleRegion(int beginningLineNumber, int endingLineNumber, java.lang.String source)
func NewSimpleRegion(beginningLineNumber int32, endingLineNumber int32, source string) *SimpleRegion {
	this_ := &SimpleRegion{}
	InitSimpleRegion(this_, beginningLineNumber, endingLineNumber, source)
	return this_
}

func InitSimpleRegion(this_ *SimpleRegion, beginningLineNumber int32, endingLineNumber int32, source string) {
	_ = beginningLineNumber
	_ = endingLineNumber
	_ = source
	this_.initFields_SimpleRegion()
	this_.beginningLineNumber = beginningLineNumber
	this_.endingLineNumber = endingLineNumber
	this_.source = source
}

// Java method: SimpleRegion.getBeginningLineNumber()
func (this_ *SimpleRegion) GetBeginningLineNumber() int32 {
	return this_.beginningLineNumber
}

// Java method: SimpleRegion.getEndingLineNumber()
func (this_ *SimpleRegion) GetEndingLineNumber() int32 {
	return this_.endingLineNumber
}

// Java method: SimpleRegion.getSourceExcerpt()
func (this_ *SimpleRegion) GetSourceExcerpt() string {
	return this_.source
}

