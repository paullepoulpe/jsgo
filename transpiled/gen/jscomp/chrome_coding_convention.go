/*
 * Copyright 2014 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/ChromeCodingConvention.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.ChromeCodingConvention
type ChromeCodingConvention struct {
	CodingConventions_ProxyBase
}

func (this_ *ChromeCodingConvention) BindSelf_(s any) {
	_ = s
	this_.CodingConventions_ProxyBase.BindSelf_(s)
}
func (this_ *ChromeCodingConvention) IsJavaNil_() bool { return this_ == nil }
func (*ChromeCodingConvention) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.ChromeCodingConvention", "ChromeCodingConvention" }
func (this_ *ChromeCodingConvention) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.CodingConventions_ProxyBase.IdentityHash_()
}
func (*ChromeCodingConvention) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ChromeCodingConvention", "com.google.javascript.jscomp.ChromeCodingConvention", func(v_ any) bool {
		return j2gort.InstanceOf[*ChromeCodingConvention](v_)
	})
}
func (*ChromeCodingConvention) Implements_jscomp_CodingConventions_Proxy() {}
func (*ChromeCodingConvention) Implements_jscomp_CodingConvention() {}

var static_ChromeCodingConvention_CR_ADD_SINGLETON_GETTER pkg_rhino.QualifiedName
func StaticGet_ChromeCodingConvention_CR_ADD_SINGLETON_GETTER() pkg_rhino.QualifiedName {
	clinit_ChromeCodingConvention()
	return static_ChromeCodingConvention_CR_ADD_SINGLETON_GETTER
}
func StaticSet_ChromeCodingConvention_CR_ADD_SINGLETON_GETTER(v pkg_rhino.QualifiedName) {
	clinit_ChromeCodingConvention()
	static_ChromeCodingConvention_CR_ADD_SINGLETON_GETTER = v
}
var clinitDone_ChromeCodingConvention uint32
var clinitState_ChromeCodingConvention j2gort.ClinitState
func clinit_ChromeCodingConvention() {
	if atomic.LoadUint32(&clinitDone_ChromeCodingConvention) != 0 {
		return
	}
	clinitSlow_ChromeCodingConvention()
}

func clinitSlow_ChromeCodingConvention() {
	if !j2gort.EnterClinit(&clinitState_ChromeCodingConvention, &clinitDone_ChromeCodingConvention) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ChromeCodingConvention, &clinitDone_ChromeCodingConvention)
	static_ChromeCodingConvention_CR_ADD_SINGLETON_GETTER = pkg_rhino.QualifiedName_Of("cr.addSingletonGetter")
}

func (this_ *ChromeCodingConvention) initFields_ChromeCodingConvention() {
	_ = this_
}

// Java constructor: ChromeCodingConvention()
func NewChromeCodingConvention() *ChromeCodingConvention {
	this_ := &ChromeCodingConvention{}
	this_.BindSelf_(this_)
	InitChromeCodingConvention(this_)
	return this_
}

func InitChromeCodingConvention(this_ *ChromeCodingConvention) {
	InitChromeCodingConvention_1(this_, CodingConventions_GetDefault())
}

// Java constructor: ChromeCodingConvention(com.google.javascript.jscomp.CodingConvention wrapped)
func NewChromeCodingConvention_1(wrapped CodingConvention) *ChromeCodingConvention {
	this_ := &ChromeCodingConvention{}
	this_.BindSelf_(this_)
	InitChromeCodingConvention_1(this_, wrapped)
	return this_
}

func InitChromeCodingConvention_1(this_ *ChromeCodingConvention, wrapped CodingConvention) {
	_ = wrapped
	InitCodingConventions_Proxy(&this_.CodingConventions_ProxyBase, wrapped)
	this_.initFields_ChromeCodingConvention()
}

// Java method: ChromeCodingConvention.getSingletonGetterClassName(com.google.javascript.rhino.Node callNode)
func (this_ *ChromeCodingConvention) GetSingletonGetterClassName(callNode pkg_rhino.Node) string {
	_ = callNode
	var callArg pkg_rhino.Node = callNode.GetFirstChild()
	_ = callArg
	if ((!(StaticGet_ChromeCodingConvention_CR_ADD_SINGLETON_GETTER().Matches(callArg)) || !(callNode.HasTwoChildren()))) {
		return this_.CodingConventions_ProxyBase.GetSingletonGetterClassName(callNode)
	}
	return callArg.GetNext().GetQualifiedName()
}

// Java method: ChromeCodingConvention.applySingletonGetter(com.google.javascript.rhino.NominalTypeBuilder classType, com.google.javascript.rhino.jstype.FunctionType getterType)
func (this_ *ChromeCodingConvention) ApplySingletonGetter(classType *pkg_rhino.NominalTypeBuilder, getterType *pkg_rhino.FunctionType) {
	_ = classType
	_ = getterType
	var defSite pkg_rhino.Node = classType.Constructor().GetSource()
	_ = defSite
	classType.DeclareConstructorProperty("getInstance", pkg_rhino.ToIface_JSType(getterType), defSite)
	classType.DeclareConstructorProperty("instance_", j2gort.AsIface[pkg_rhino.JSType](classType.Instance()), defSite)
}

