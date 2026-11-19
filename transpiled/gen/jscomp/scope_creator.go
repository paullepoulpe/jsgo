/*
 * Copyright 2006 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/ScopeCreator.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java interface: com.google.javascript.jscomp.ScopeCreator
type ScopeCreator interface {
	Implements_jscomp_ScopeCreator()
	IsJavaNil_() bool
	CreateScope_2_Node_AbstractScope(pkg_rhino.Node, AbstractScope) any
}

func ToIface_ScopeCreator[T interface {
	ScopeCreator
	comparable
}](v T) ScopeCreator {
	var z T
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ScopeCreator", "com.google.javascript.jscomp.ScopeCreator", func(v_ any) bool {
		return j2gort.InstanceOf[ScopeCreator](v_)
	})
}

var static_ScopeCreator_ASSERT_NO_SCOPES_CREATED ScopeCreator

func StaticGet_ScopeCreator_ASSERT_NO_SCOPES_CREATED() ScopeCreator {
	clinit_ScopeCreator()
	return static_ScopeCreator_ASSERT_NO_SCOPES_CREATED
}
func StaticSet_ScopeCreator_ASSERT_NO_SCOPES_CREATED(v ScopeCreator) {
	clinit_ScopeCreator()
	static_ScopeCreator_ASSERT_NO_SCOPES_CREATED = v
}

var clinitDone_ScopeCreator uint32
var clinitState_ScopeCreator j2gort.ClinitState

func clinit_ScopeCreator() {
	if atomic.LoadUint32(&clinitDone_ScopeCreator) != 0 {
		return
	}
	clinitSlow_ScopeCreator()
}

func clinitSlow_ScopeCreator() {
	if !j2gort.EnterClinit(&clinitState_ScopeCreator, &clinitDone_ScopeCreator) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ScopeCreator, &clinitDone_ScopeCreator)
	static_ScopeCreator_ASSERT_NO_SCOPES_CREATED = NewScopeCreator__Anon_78()
}

// Java class:
type ScopeCreator__Anon_78 struct {
	idHash_ atomic.Int32
}

func (this_ *ScopeCreator__Anon_78) BindSelf_(s any) {
	_ = s
}
func (this_ *ScopeCreator__Anon_78) IsJavaNil_() bool { return this_ == nil }
func (*ScopeCreator__Anon_78) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ScopeCreator$1", ""
}
func (this_ *ScopeCreator__Anon_78) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ScopeCreator__Anon_78) HasIdentityEquals_()             {}
func (*ScopeCreator__Anon_78) Implements_jscomp_ScopeCreator() {}

func (this_ *ScopeCreator__Anon_78) initFields_ScopeCreator__Anon_78() {
	_ = this_
}

func NewScopeCreator__Anon_78() *ScopeCreator__Anon_78 {
	this_ := &ScopeCreator__Anon_78{}
	this_.initFields_ScopeCreator__Anon_78()
	return this_
}

// Java method: .createScope(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.AbstractScope<?,?> parent)
func (this_ *ScopeCreator__Anon_78) CreateScope_2_Node_AbstractScope(n pkg_rhino.Node, parent AbstractScope) any {
	_ = n
	_ = parent
	panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "ScopeCreator.ASSERT_NO_SCOPES_CREATED cannot create child scopes.")))
}
