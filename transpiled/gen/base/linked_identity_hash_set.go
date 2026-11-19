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
// Port of com/google/javascript/jscomp/base/LinkedIdentityHashSet.java (generated mechanically by j2go).

package base

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.base.LinkedIdentityHashSet
type LinkedIdentityHashSet[E comparable] struct {
	innerMap *LinkedIdentityHashMap[E, *LinkedIdentityHashSet_Sentinel] `gson:"innerMap"`
	idHash_ atomic.Int32
}

func (this_ *LinkedIdentityHashSet[E]) BindSelf_(s any) {
	_ = s
}
func (this_ *LinkedIdentityHashSet[E]) IsJavaNil_() bool { return this_ == nil }
func (*LinkedIdentityHashSet[E]) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.base.LinkedIdentityHashSet", "LinkedIdentityHashSet" }
func (this_ *LinkedIdentityHashSet[E]) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*LinkedIdentityHashSet[E]) HasIdentityEquals_() {}
func (*LinkedIdentityHashSet[E]) IsGenericInst_base_LinkedIdentityHashSet_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.base.LinkedIdentityHashSet", "com.google.javascript.jscomp.base.LinkedIdentityHashSet", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			IsGenericInst_base_LinkedIdentityHashSet_()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}

var static_LinkedIdentityHashSet_SENTINEL *LinkedIdentityHashSet_Sentinel
func StaticGet_LinkedIdentityHashSet_SENTINEL() *LinkedIdentityHashSet_Sentinel {
	clinit_LinkedIdentityHashSet()
	return static_LinkedIdentityHashSet_SENTINEL
}
func StaticSet_LinkedIdentityHashSet_SENTINEL(v *LinkedIdentityHashSet_Sentinel) {
	clinit_LinkedIdentityHashSet()
	static_LinkedIdentityHashSet_SENTINEL = v
}
var clinitDone_LinkedIdentityHashSet uint32
var clinitState_LinkedIdentityHashSet j2gort.ClinitState
func clinit_LinkedIdentityHashSet() {
	if atomic.LoadUint32(&clinitDone_LinkedIdentityHashSet) != 0 {
		return
	}
	clinitSlow_LinkedIdentityHashSet()
}

func clinitSlow_LinkedIdentityHashSet() {
	if !j2gort.EnterClinit(&clinitState_LinkedIdentityHashSet, &clinitDone_LinkedIdentityHashSet) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_LinkedIdentityHashSet, &clinitDone_LinkedIdentityHashSet)
	static_LinkedIdentityHashSet_SENTINEL = NewLinkedIdentityHashSet_Sentinel()
}

func (this_ *LinkedIdentityHashSet[E]) initFields_LinkedIdentityHashSet() {
	_ = this_
	this_.innerMap = NewLinkedIdentityHashMap[E, *LinkedIdentityHashSet_Sentinel]()
}

// Java constructor: LinkedIdentityHashSet()
func NewLinkedIdentityHashSet[E comparable]() *LinkedIdentityHashSet[E] {
	clinit_LinkedIdentityHashSet()
	this_ := &LinkedIdentityHashSet[E]{}
	InitLinkedIdentityHashSet[E](this_)
	return this_
}

func InitLinkedIdentityHashSet[E comparable](this_ *LinkedIdentityHashSet[E]) {
	this_.initFields_LinkedIdentityHashSet()
}

// Java method: LinkedIdentityHashSet.contains(E element)
func (this_ *LinkedIdentityHashSet[E]) Contains(element E) bool {
	_ = element
	return (this_.innerMap.Get(element) == StaticGet_LinkedIdentityHashSet_SENTINEL())
}

// Java method: LinkedIdentityHashSet.add(E element)
func (this_ *LinkedIdentityHashSet[E]) Add(element E) bool {
	_ = element
	return (this_.innerMap.Put(element, StaticGet_LinkedIdentityHashSet_SENTINEL()) == nil)
}

// Java method: LinkedIdentityHashSet.remove(E element)
func (this_ *LinkedIdentityHashSet[E]) Remove(element E) bool {
	_ = element
	return (this_.innerMap.Put(element, nil) == StaticGet_LinkedIdentityHashSet_SENTINEL())
}

// Java class: com.google.javascript.jscomp.base.LinkedIdentityHashSet.Sentinel
type LinkedIdentityHashSet_Sentinel struct {
	idHash_ atomic.Int32
}

func (this_ *LinkedIdentityHashSet_Sentinel) BindSelf_(s any) {
	_ = s
}
func (this_ *LinkedIdentityHashSet_Sentinel) IsJavaNil_() bool { return this_ == nil }
func (*LinkedIdentityHashSet_Sentinel) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.base.LinkedIdentityHashSet$Sentinel", "Sentinel" }
func (this_ *LinkedIdentityHashSet_Sentinel) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*LinkedIdentityHashSet_Sentinel) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.base.LinkedIdentityHashSet.Sentinel", "com.google.javascript.jscomp.base.LinkedIdentityHashSet$Sentinel", func(v_ any) bool {
		return j2gort.InstanceOf[*LinkedIdentityHashSet_Sentinel](v_)
	})
}

func (this_ *LinkedIdentityHashSet_Sentinel) initFields_LinkedIdentityHashSet_Sentinel() {
	_ = this_
}

// Java constructor: Sentinel()
func NewLinkedIdentityHashSet_Sentinel() *LinkedIdentityHashSet_Sentinel {
	this_ := &LinkedIdentityHashSet_Sentinel{}
	InitLinkedIdentityHashSet_Sentinel(this_)
	return this_
}

func InitLinkedIdentityHashSet_Sentinel(this_ *LinkedIdentityHashSet_Sentinel) {
	this_.initFields_LinkedIdentityHashSet_Sentinel()
}

