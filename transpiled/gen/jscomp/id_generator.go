/*
 * Copyright 2017 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/IdGenerator.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.IdGenerator
type IdGenerator struct {
	currentId int32 `gson:"currentId"`
	idHash_   atomic.Int32
}

func (this_ *IdGenerator) BindSelf_(s any) {
	_ = s
}
func (this_ *IdGenerator) IsJavaNil_() bool { return this_ == nil }
func (*IdGenerator) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.IdGenerator", "IdGenerator"
}
func (this_ *IdGenerator) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*IdGenerator) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.IdGenerator", "com.google.javascript.jscomp.IdGenerator", func(v_ any) bool {
		return j2gort.InstanceOf[*IdGenerator](v_)
	})
}

func (this_ *IdGenerator) initFields_IdGenerator() {
	_ = this_

}

// Java constructor: IdGenerator()
func NewIdGenerator() *IdGenerator {
	this_ := &IdGenerator{}
	InitIdGenerator(this_)
	return this_
}

func InitIdGenerator(this_ *IdGenerator) {
	this_.initFields_IdGenerator()
}

// Java method: IdGenerator.hasGeneratedAnyIds()
func (this_ *IdGenerator) HasGeneratedAnyIds() bool {
	return (this_.currentId != int32(0))
}

// Java method: IdGenerator.newId()
func (this_ *IdGenerator) NewId() int32 {
	return func() int32 { old_1 := this_.currentId; nw_2 := old_1 + 1; this_.currentId = nw_2; return old_1 }()
}

// Java method: IdGenerator.getCurrentId()
func (this_ *IdGenerator) GetCurrentId() int32 {
	return this_.currentId
}

// Java method: IdGenerator.setCurrentId(int id)
func (this_ *IdGenerator) SetCurrentId(id int32) {
	_ = id
	this_.currentId = id
}
