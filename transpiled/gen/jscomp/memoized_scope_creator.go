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
// Port of com/google/javascript/jscomp/MemoizedScopeCreator.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.MemoizedScopeCreator
type MemoizedScopeCreator struct {
	scopesByScopeRoot j2gort.OrderedMap[pkg_rhino.Node, AbstractScope] `gson:"scopesByScopeRoot"`
	delegate          ScopeCreator                                     `gson:"delegate"`
	idHash_           atomic.Int32
}

func (this_ *MemoizedScopeCreator) BindSelf_(s any) {
	_ = s
}
func (this_ *MemoizedScopeCreator) IsJavaNil_() bool { return this_ == nil }
func (*MemoizedScopeCreator) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.MemoizedScopeCreator", "MemoizedScopeCreator"
}
func (this_ *MemoizedScopeCreator) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*MemoizedScopeCreator) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.MemoizedScopeCreator", "com.google.javascript.jscomp.MemoizedScopeCreator", func(v_ any) bool {
		return j2gort.InstanceOf[*MemoizedScopeCreator](v_)
	})
}
func (*MemoizedScopeCreator) Implements_jscomp_ScopeCreator() {}

func (this_ *MemoizedScopeCreator) initFields_MemoizedScopeCreator() {
	_ = this_
}

// Java constructor: MemoizedScopeCreator(com.google.javascript.jscomp.ScopeCreator delegate)
func NewMemoizedScopeCreator_1(delegate ScopeCreator) *MemoizedScopeCreator {
	this_ := &MemoizedScopeCreator{}
	InitMemoizedScopeCreator_1(this_, delegate)
	return this_
}

func InitMemoizedScopeCreator_1(this_ *MemoizedScopeCreator, delegate ScopeCreator) {
	_ = delegate
	InitMemoizedScopeCreator_2(this_, delegate, int32(16))
}

// Java constructor: MemoizedScopeCreator(com.google.javascript.jscomp.ScopeCreator delegate, int expectedScopes)
func NewMemoizedScopeCreator_2(delegate ScopeCreator, expectedScopes int32) *MemoizedScopeCreator {
	this_ := &MemoizedScopeCreator{}
	InitMemoizedScopeCreator_2(this_, delegate, expectedScopes)
	return this_
}

func InitMemoizedScopeCreator_2(this_ *MemoizedScopeCreator, delegate ScopeCreator, expectedScopes int32) {
	_ = delegate
	_ = expectedScopes
	this_.initFields_MemoizedScopeCreator()
	this_.scopesByScopeRoot = j2gort.NewOrderedMapCap[pkg_rhino.Node, AbstractScope](expectedScopes)
	this_.delegate = delegate
}

// Java method: MemoizedScopeCreator.createScope(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.AbstractScope<?,?> parent)
func (this_ *MemoizedScopeCreator) CreateScope_2_Node_AbstractScope(n pkg_rhino.Node, parent AbstractScope) any {
	_ = n
	_ = parent
	var scope AbstractScope = this_.scopesByScopeRoot.Get(n)
	_ = scope
	if scope == nil || scope.IsJavaNil_() {
		scope = j2gort.AsIface[AbstractScope](this_.delegate.CreateScope_2_Node_AbstractScope(n, parent))
		this_.scopesByScopeRoot.Put(n, scope)
	} else {
		if !(j2gort.IfaceEquals[AbstractScope](parent, j2gort.AsIface[AbstractScope](scope.GetParent()))) {
			j2gort.CheckState(false)
		}
	}
	return scope
}
func (this_ *MemoizedScopeCreator) ReleaseSlabs_() {
	if this_ == nil {
		return
	}
	if this_.delegate != nil {
		if r_, ok_ := any(this_.delegate).(interface{ ReleaseSlabs_() }); ok_ {
			r_.ReleaseSlabs_()
		}
	}
	this_.scopesByScopeRoot.ReleaseToPool()

}
